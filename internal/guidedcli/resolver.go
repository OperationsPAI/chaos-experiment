package guidedcli

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/OperationsPAI/chaos-experiment/handler"
	"github.com/OperationsPAI/chaos-experiment/internal/resourcelookup"
	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

const defaultDurationMinutes = 5

func Resolve(ctx context.Context, cfg GuidedConfig) (*GuidedResponse, error) {
	if err := normalizeSystemSelection(&cfg); err != nil {
		return errorResponse(cfg, "invalid_system", err), nil
	}

	if cfg.System == "" {
		return resolveSystems(cfg), nil
	}

	systemType, err := systemconfig.ParseSystemType(cfg.SystemType)
	if err != nil {
		return errorResponse(cfg, "invalid_system_type", err), nil
	}
	if err := systemconfig.SetCurrentSystem(systemType); err != nil {
		return errorResponse(cfg, "invalid_system_type", err), nil
	}

	if cfg.App == "" {
		response, err := resolveApps(ctx, cfg, systemType)
		if err != nil {
			return resolutionErrorResponse(cfg, err), nil
		}
		return response, nil
	}

	if cfg.ChaosType == "" {
		response, err := resolveChaosTypes(ctx, cfg, systemType)
		if err != nil {
			return resolutionErrorResponse(cfg, err), nil
		}
		return response, nil
	}

	var response *GuidedResponse
	switch cfg.ChaosType {
	case "PodKill":
		response, err = resolvePodKill(ctx, cfg, systemType)
	case "CPUStress":
		response, err = resolveCPUStress(ctx, cfg, systemType)
	case "NetworkDelay":
		response, err = resolveNetworkDelay(ctx, cfg, systemType)
	case "HTTPRequestDelay":
		response, err = resolveHTTPRequestDelay(ctx, cfg, systemType)
	case "JVMLatency":
		response, err = resolveJVMLatency(ctx, cfg, systemType)
	case "JVMRuntimeMutator":
		response, err = resolveJVMRuntimeMutator(ctx, cfg, systemType)
	default:
		return &GuidedResponse{
			Mode:     "guided",
			Stage:    "unsupported_chaos_type",
			Config:   cfg,
			Resolved: resolvedMap(cfg),
			Errors: []string{
				fmt.Sprintf("chaos type %q is not implemented in the current guided CLI", cfg.ChaosType),
			},
		}, nil
	}
	if err != nil {
		return resolutionErrorResponse(cfg, err), nil
	}
	return response, nil
}

func resolveSystems(cfg GuidedConfig) *GuidedResponse {
	instances, warnings := discoverSystemInstances()
	options := make([]FieldOption, 0, len(instances))
	resourceItems := make([]map[string]any, 0, len(instances))
	for _, instance := range instances {
		options = append(options, FieldOption{
			Value: instance.Name,
			Label: instance.Name,
			Metadata: map[string]any{
				"system_type":  instance.SystemType.String(),
				"namespace":    instance.Namespace,
				"display_name": instance.DisplayName,
			},
		})
		resourceItems = append(resourceItems, map[string]any{
			"name":         instance.Name,
			"system_type":  instance.SystemType.String(),
			"namespace":    instance.Namespace,
			"display_name": instance.DisplayName,
		})
	}

	return &GuidedResponse{
		Mode:     "guided",
		Stage:    "select_system",
		Config:   cfg,
		Resolved: resolvedMap(cfg),
		Next: []FieldSpec{{
			Name:        "system",
			Kind:        "enum",
			Required:    true,
			Description: "Select a system namespace instance",
			Options:     options,
		}},
		Resources: map[string]any{"systems": resourceItems},
		Warnings:  warnings,
	}
}

func resolveApps(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (*GuidedResponse, error) {
	apps, err := safeAppLabels(cfg.Namespace, systemType)
	if err != nil {
		return nil, err
	}
	options := make([]FieldOption, 0, len(apps))
	for _, app := range apps {
		options = append(options, FieldOption{Value: app, Label: app})
	}

	return &GuidedResponse{
		Mode:     "guided",
		Stage:    "select_app",
		Config:   cfg,
		Resolved: resolvedMap(cfg),
		Next: []FieldSpec{{
			Name:        "app",
			Kind:        "enum",
			Required:    true,
			Description: "Select an app label in the namespace",
			Options:     options,
		}},
	}, nil
}

func resolveChaosTypes(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (*GuidedResponse, error) {
	options, summary, err := availableChaosTypeOptions(ctx, cfg, systemType)
	if err != nil {
		return nil, err
	}

	return &GuidedResponse{
		Mode:     "guided",
		Stage:    "select_chaos_type",
		Config:   cfg,
		Resolved: resolvedMap(cfg),
		Next: []FieldSpec{{
			Name:        "chaos_type",
			Kind:        "enum",
			Required:    true,
			Description: "Select a chaos type supported for the current app",
			Options:     options,
		}},
		Preview: &Preview{ResourceSummary: summary},
	}, nil
}

func resolvePodKill(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (*GuidedResponse, error) {
	return finalizeOrRequest(ctx, cfg, systemType, []FieldSpec{optionalDurationField()}, buildPodKill)
}

func resolveCPUStress(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (*GuidedResponse, error) {
	containers, err := containersByApp(ctx, systemType, cfg.Namespace, cfg.App)
	if err != nil {
		return nil, err
	}
	if cfg.Container == "" {
		options := make([]FieldOption, 0, len(containers))
		for _, container := range containers {
			options = append(options, FieldOption{Value: container, Label: container})
		}
		return &GuidedResponse{
			Mode:     "guided",
			Stage:    "select_container",
			Config:   cfg,
			Resolved: resolvedMap(cfg),
			Next: []FieldSpec{{
				Name:        "container",
				Kind:        "enum",
				Required:    true,
				Description: "Select a container under the app",
				Options:     options,
			}},
		}, nil
	}

	return finalizeOrRequest(ctx, cfg, systemType, []FieldSpec{{
		Name:        "params",
		Kind:        "group",
		Required:    true,
		Description: "Fill CPU stress parameters",
		Fields: []FieldSpec{
			optionalDurationField(),
			requiredNumberField("cpu_load", "CPU load percentage", 1, 100, 1, "%"),
			requiredNumberField("cpu_worker", "CPU stress worker count", 1, 3, 1, ""),
		},
	}}, buildCPUStress)
}

func resolveNetworkDelay(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (*GuidedResponse, error) {
	targets, err := networkTargetsByApp(systemType, cfg.App)
	if err != nil {
		return nil, err
	}
	if cfg.TargetService == "" {
		options := make([]FieldOption, 0, len(targets))
		for _, target := range targets {
			options = append(options, FieldOption{Value: target, Label: target})
		}
		return &GuidedResponse{
			Mode:     "guided",
			Stage:    "select_network_target",
			Config:   cfg,
			Resolved: resolvedMap(cfg),
			Next: []FieldSpec{{
				Name:        "target_service",
				Kind:        "enum",
				Required:    true,
				Description: "Select the target service for the network pair",
				Options:     options,
			}},
		}, nil
	}

	return finalizeOrRequest(ctx, cfg, systemType, []FieldSpec{{
		Name:        "params",
		Kind:        "group",
		Required:    true,
		Description: "Fill network delay parameters",
		Fields: []FieldSpec{
			optionalDurationField(),
			requiredNumberField("latency", "Network latency", 1, 2000, 1, "ms"),
			requiredNumberField("correlation", "Correlation percentage", 0, 100, 1, "%"),
			requiredNumberField("jitter", "Jitter", 0, 1000, 1, "ms"),
			{
				Name:        "direction",
				Kind:        "enum",
				Required:    true,
				Description: "Traffic direction",
				Options:     []FieldOption{{Value: "to", Label: "to"}, {Value: "from", Label: "from"}, {Value: "both", Label: "both"}},
			},
		},
	}}, buildNetworkDelay)
}

func resolveHTTPRequestDelay(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (*GuidedResponse, error) {
	endpoints, err := httpEndpointsByApp(systemType, cfg.App)
	if err != nil {
		return nil, err
	}
	if cfg.Route == "" || cfg.HTTPMethod == "" {
		options := make([]FieldOption, 0, len(endpoints))
		for _, endpoint := range endpoints {
			options = append(options, FieldOption{
				Value: endpoint.Method + " " + endpoint.Route,
				Label: endpoint.Method + " " + endpoint.Route,
				Metadata: map[string]any{
					"http_method":    endpoint.Method,
					"route":          endpoint.Route,
					"target_service": endpoint.ServerAddress,
					"span_name":      endpoint.SpanName,
				},
			})
		}
		return &GuidedResponse{
			Mode:     "guided",
			Stage:    "select_http_endpoint",
			Config:   cfg,
			Resolved: resolvedMap(cfg),
			Next: []FieldSpec{{
				Name:        "endpoint",
				Kind:        "object_ref",
				Required:    true,
				Description: "Select the HTTP endpoint for request delay",
				KeyFields:   []string{"http_method", "route"},
				Options:     options,
			}},
		}, nil
	}

	return finalizeOrRequest(ctx, cfg, systemType, []FieldSpec{{
		Name:        "params",
		Kind:        "group",
		Required:    true,
		Description: "Fill HTTP request delay parameters",
		Fields: []FieldSpec{
			optionalDurationField(),
			requiredNumberField("delay_duration", "Request delay duration", 10, 5000, 1, "ms"),
		},
	}}, buildHTTPRequestDelay)
}

func resolveJVMLatency(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (*GuidedResponse, error) {
	methods, err := jvmMethodsByApp(systemType, cfg.App)
	if err != nil {
		return nil, err
	}
	if cfg.Class == "" || cfg.Method == "" {
		options := make([]FieldOption, 0, len(methods))
		for _, method := range methods {
			options = append(options, FieldOption{
				Value:    method.ClassName + "#" + method.MethodName,
				Label:    method.ClassName + "#" + method.MethodName,
				Metadata: map[string]any{"class": method.ClassName, "method": method.MethodName},
			})
		}
		return &GuidedResponse{
			Mode:     "guided",
			Stage:    "select_jvm_method",
			Config:   cfg,
			Resolved: resolvedMap(cfg),
			Next: []FieldSpec{{
				Name:        "method_ref",
				Kind:        "object_ref",
				Required:    true,
				Description: "Select the JVM method for latency injection",
				KeyFields:   []string{"class", "method"},
				Options:     options,
			}},
		}, nil
	}

	return finalizeOrRequest(ctx, cfg, systemType, []FieldSpec{{
		Name:        "params",
		Kind:        "group",
		Required:    true,
		Description: "Fill JVM latency parameters",
		Fields: []FieldSpec{
			optionalDurationField(),
			requiredNumberField("latency_duration", "JVM latency duration", 1, 5000, 1, "ms"),
		},
	}}, buildJVMLatency)
}

func resolveJVMRuntimeMutator(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (*GuidedResponse, error) {
	methods, err := runtimeMutatorMethodsByApp(systemType, cfg.App)
	if err != nil {
		return nil, err
	}
	if cfg.Class == "" || cfg.Method == "" {
		options := make([]FieldOption, 0, len(methods))
		for _, method := range methods {
			options = append(options, FieldOption{
				Value:    method.ClassName + "#" + method.MethodName,
				Label:    method.ClassName + "#" + method.MethodName,
				Metadata: map[string]any{"class": method.ClassName, "method": method.MethodName},
			})
		}
		return &GuidedResponse{
			Mode:     "guided",
			Stage:    "select_runtime_mutator_method",
			Config:   cfg,
			Resolved: resolvedMap(cfg),
			Next: []FieldSpec{{
				Name:        "method_ref",
				Kind:        "object_ref",
				Required:    true,
				Description: "Select the method for runtime mutator injection",
				KeyFields:   []string{"class", "method"},
				Options:     options,
			}},
		}, nil
	}

	mutators, err := runtimeMutatorsByMethod(systemType, cfg.App, cfg.Class, cfg.Method)
	if err != nil {
		return nil, err
	}
	if cfg.MutatorConfig == "" {
		options := make([]FieldOption, 0, len(mutators))
		for _, mutator := range mutators {
			options = append(options, FieldOption{
				Value: runtimeMutatorKey(mutator),
				Label: runtimeMutatorLabel(mutator),
				Metadata: map[string]any{
					"mutation_type_name": mutator.MutationTypeName,
					"mutation_strategy":  mutator.MutationStrategy,
					"mutation_from":      mutator.MutationFrom,
					"mutation_to":        mutator.MutationTo,
					"description":        mutator.Description,
				},
			})
		}
		return &GuidedResponse{
			Mode:     "guided",
			Stage:    "select_runtime_mutator_config",
			Config:   cfg,
			Resolved: resolvedMap(cfg),
			Next: []FieldSpec{
				{
					Name:        "mutator_config",
					Kind:        "enum",
					Required:    true,
					Description: "Select the runtime mutator config",
					Options:     options,
				},
				optionalDurationField(),
			},
		}, nil
	}

	return finalizeOrRequest(ctx, cfg, systemType, []FieldSpec{optionalDurationField()}, buildJVMRuntimeMutator)
}

type buildFunc func(context.Context, GuidedConfig, systemconfig.SystemType) (handler.InjectionConf, map[string]any, error)

func finalizeOrRequest(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType, next []FieldSpec, builder buildFunc) (*GuidedResponse, error) {
	conf, payload, err := builder(ctx, cfg, systemType)
	if err != nil {
		return &GuidedResponse{
			Mode:     "guided",
			Stage:    "fill_required_fields",
			Config:   cfg,
			Resolved: resolvedMap(cfg),
			Next:     next,
			Errors:   []string{err.Error()},
		}, nil
	}

	normalized := normalizeDuration(cfg)
	response := &GuidedResponse{
		Mode:         "guided",
		Stage:        "ready_to_apply",
		Config:       normalized,
		Resolved:     resolvedMap(normalized),
		Next:         next,
		Preview:      buildPreview(normalized, payload, systemType),
		ApplyPayload: payload,
		CanApply:     true,
	}

	if cfg.Apply {
		names, err := handler.BatchCreate(ctx, []handler.InjectionConf{conf}, handler.SystemType(systemType), cfg.Namespace, map[string]string{}, map[string]string{})
		if err != nil {
			response.Errors = []string{err.Error()}
			return response, nil
		}
		response.Stage = "applied"
		response.Result = map[string]any{
			"created":   names,
			"count":     len(names),
			"namespace": cfg.Namespace,
			"system":    cfg.System,
		}
	}

	return response, nil
}

func availableChaosTypeOptions(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) ([]FieldOption, map[string]any, error) {
	options := make([]FieldOption, 0)
	summary := map[string]any{}

	containers, _ := containersByApp(ctx, systemType, cfg.Namespace, cfg.App)
	if len(containers) > 0 {
		options = append(options,
			FieldOption{Value: "PodKill", Label: "PodKill", Description: "Kill pods for the selected app"},
			FieldOption{Value: "CPUStress", Label: "CPUStress", Description: "Stress a container with CPU load"},
		)
		summary["containers"] = len(containers)
	}

	endpoints, _ := httpEndpointsByApp(systemType, cfg.App)
	if len(endpoints) > 0 {
		options = append(options, FieldOption{Value: "HTTPRequestDelay", Label: "HTTPRequestDelay", Description: "Delay HTTP requests for a selected endpoint"})
		summary["http_endpoints"] = len(endpoints)
	}

	networkTargets, _ := networkTargetsByApp(systemType, cfg.App)
	if len(networkTargets) > 0 {
		options = append(options, FieldOption{Value: "NetworkDelay", Label: "NetworkDelay", Description: "Delay traffic to a downstream service"})
		summary["network_targets"] = len(networkTargets)
	}

	methods, _ := jvmMethodsByApp(systemType, cfg.App)
	if len(methods) > 0 {
		options = append(options, FieldOption{Value: "JVMLatency", Label: "JVMLatency", Description: "Inject latency into a JVM method"})
		summary["jvm_methods"] = len(methods)
	}

	mutatorMethods, _ := runtimeMutatorMethodsByApp(systemType, cfg.App)
	if len(mutatorMethods) > 0 {
		options = append(options, FieldOption{Value: "JVMRuntimeMutator", Label: "JVMRuntimeMutator", Description: "Apply a runtime mutator strategy to a JVM method"})
		summary["runtime_mutator_methods"] = len(mutatorMethods)
	}

	sort.Slice(options, func(i, j int) bool { return options[i].Value < options[j].Value })
	return options, summary, nil
}

func containersByApp(ctx context.Context, systemType systemconfig.SystemType, namespace, app string) ([]string, error) {
	allContainers, err := safeContainers(namespace)
	if err != nil {
		return nil, err
	}
	containers := make([]string, 0)
	seen := map[string]bool{}
	for _, container := range allContainers {
		if container.AppLabel == app && !seen[container.ContainerName] {
			seen[container.ContainerName] = true
			containers = append(containers, container.ContainerName)
		}
	}
	sort.Strings(containers)
	return containers, nil
}

func networkTargetsByApp(systemType systemconfig.SystemType, app string) ([]string, error) {
	pairs, err := resourcelookup.GetSystemCache(systemType).GetAllNetworkPairs()
	if err != nil {
		return nil, err
	}
	targets := make([]string, 0)
	seen := map[string]bool{}
	for _, pair := range pairs {
		if pair.SourceService == app && !seen[pair.TargetService] {
			seen[pair.TargetService] = true
			targets = append(targets, pair.TargetService)
		}
	}
	sort.Strings(targets)
	return targets, nil
}

func httpEndpointsByApp(systemType systemconfig.SystemType, app string) ([]resourcelookup.AppEndpointPair, error) {
	endpoints, err := resourcelookup.GetSystemCache(systemType).GetAllHTTPEndpoints()
	if err != nil {
		return nil, err
	}
	result := make([]resourcelookup.AppEndpointPair, 0)
	for _, endpoint := range endpoints {
		if endpoint.AppName == app {
			result = append(result, endpoint)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Method != result[j].Method {
			return result[i].Method < result[j].Method
		}
		return result[i].Route < result[j].Route
	})
	return result, nil
}

func jvmMethodsByApp(systemType systemconfig.SystemType, app string) ([]resourcelookup.AppMethodPair, error) {
	methods, err := resourcelookup.GetSystemCache(systemType).GetAllJVMMethods()
	if err != nil {
		return nil, err
	}
	result := make([]resourcelookup.AppMethodPair, 0)
	for _, method := range methods {
		if method.AppName == app {
			result = append(result, method)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].ClassName != result[j].ClassName {
			return result[i].ClassName < result[j].ClassName
		}
		return result[i].MethodName < result[j].MethodName
	})
	return result, nil
}

type runtimeMutatorTarget = resourcelookup.AppRuntimeMutatorTarget

func runtimeMutatorMethodsByApp(systemType systemconfig.SystemType, app string) ([]runtimeMutatorTarget, error) {
	targets, err := resourcelookup.GetSystemCache(systemType).GetAllJVMRuntimeMutatorTargets()
	if err != nil {
		return nil, err
	}
	result := make([]runtimeMutatorTarget, 0)
	seen := map[string]bool{}
	for _, target := range targets {
		if target.AppName != app {
			continue
		}
		key := target.AppName + "|" + target.ClassName + "|" + target.MethodName
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, target)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].ClassName != result[j].ClassName {
			return result[i].ClassName < result[j].ClassName
		}
		return result[i].MethodName < result[j].MethodName
	})
	return result, nil
}

func runtimeMutatorsByMethod(systemType systemconfig.SystemType, app, className, methodName string) ([]runtimeMutatorTarget, error) {
	targets, err := resourcelookup.GetSystemCache(systemType).GetAllJVMRuntimeMutatorTargets()
	if err != nil {
		return nil, err
	}
	result := make([]runtimeMutatorTarget, 0)
	for _, target := range targets {
		if target.AppName == app && target.ClassName == className && target.MethodName == methodName {
			result = append(result, target)
		}
	}
	sort.Slice(result, func(i, j int) bool { return runtimeMutatorKey(result[i]) < runtimeMutatorKey(result[j]) })
	return result, nil
}

func runtimeMutatorKey(target runtimeMutatorTarget) string {
	switch target.MutationTypeName {
	case "constant":
		return strings.Join([]string{"constant", target.MutationFrom, target.MutationTo}, ":")
	case "operator", "string":
		return strings.Join([]string{target.MutationTypeName, target.MutationStrategy}, ":")
	default:
		return target.MutationTypeName
	}
}

func runtimeMutatorLabel(target runtimeMutatorTarget) string {
	if target.Description != "" {
		return target.Description
	}
	return runtimeMutatorKey(target)
}

func optionalDurationField() FieldSpec {
	return FieldSpec{
		Name:        "duration",
		Kind:        "number_range",
		Required:    false,
		Description: "Fault duration in minutes",
		Min:         intPtr(1),
		Max:         intPtr(60),
		Step:        intPtr(1),
		Default:     intPtr(defaultDurationMinutes),
		Unit:        "minute",
	}
}

func requiredNumberField(name, description string, min, max, step int, unit string) FieldSpec {
	return FieldSpec{
		Name:        name,
		Kind:        "number_range",
		Required:    true,
		Description: description,
		Min:         intPtr(min),
		Max:         intPtr(max),
		Step:        intPtr(step),
		Unit:        unit,
	}
}

func normalizeDuration(cfg GuidedConfig) GuidedConfig {
	if cfg.Duration == nil {
		cfg.Duration = intPtr(defaultDurationMinutes)
	}
	return cfg
}

func errorResponse(cfg GuidedConfig, stage string, err error) *GuidedResponse {
	return &GuidedResponse{Mode: "guided", Stage: stage, Config: cfg, Resolved: resolvedMap(cfg), Errors: []string{err.Error()}}
}

func resolutionErrorResponse(cfg GuidedConfig, err error) *GuidedResponse {
	return &GuidedResponse{
		Mode:     "guided",
		Stage:    stageForConfig(cfg),
		Config:   cfg,
		Resolved: resolvedMap(cfg),
		Errors:   []string{err.Error()},
	}
}

func stageForConfig(cfg GuidedConfig) string {
	if cfg.System == "" {
		return "select_system"
	}
	if cfg.App == "" {
		return "select_app"
	}
	if cfg.ChaosType == "" {
		return "select_chaos_type"
	}

	switch cfg.ChaosType {
	case "PodKill":
		return "fill_required_fields"
	case "CPUStress":
		if cfg.Container == "" {
			return "select_container"
		}
		return "fill_required_fields"
	case "NetworkDelay":
		if cfg.TargetService == "" {
			return "select_network_target"
		}
		return "fill_required_fields"
	case "HTTPRequestDelay":
		if cfg.Route == "" || cfg.HTTPMethod == "" {
			return "select_http_endpoint"
		}
		return "fill_required_fields"
	case "JVMLatency":
		if cfg.Class == "" || cfg.Method == "" {
			return "select_jvm_method"
		}
		return "fill_required_fields"
	case "JVMRuntimeMutator":
		if cfg.Class == "" || cfg.Method == "" {
			return "select_runtime_mutator_method"
		}
		if cfg.MutatorConfig == "" {
			return "select_runtime_mutator_config"
		}
		return "fill_required_fields"
	default:
		return "fill_required_fields"
	}
}
