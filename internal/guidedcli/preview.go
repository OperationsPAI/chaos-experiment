package guidedcli

import (
	"github.com/OperationsPAI/chaos-experiment/handler"
	"github.com/OperationsPAI/chaos-experiment/internal/resourcelookup"
	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

func buildPreview(cfg GuidedConfig, payload map[string]any, systemType systemconfig.SystemType) *Preview {
	displayConfig := map[string]any{
		"chaos_type": cfg.ChaosType,
		"system":     cfg.System,
		"namespace":  cfg.Namespace,
		"duration":   normalizedDuration(cfg),
	}
	groundtruth := map[string]any{}

	switch cfg.ChaosType {
	case "PodKill":
		displayConfig["injection_point"] = map[string]any{"app_name": cfg.App}
		groundtruth["service"] = []string{cfg.App}
	case "CPUStress":
		displayConfig["injection_point"] = map[string]any{
			"app_name":        cfg.App,
			"container_name":  cfg.Container,
			"namespace":       cfg.Namespace,
			"container_guess": true,
		}
		displayConfig["cpu_load"] = ptrValue(cfg.CPULoad)
		displayConfig["cpu_worker"] = ptrValue(cfg.CPUWorker)
		groundtruth["service"] = []string{cfg.App}
		if cfg.Container != "" {
			groundtruth["container"] = []string{cfg.Container}
		}
		groundtruth["metric"] = []string{string(handler.MetricCPU)}
	case "NetworkDelay":
		displayConfig["injection_point"] = lookupNetworkPair(cfg, systemType)
		displayConfig["latency"] = ptrValue(cfg.Latency)
		displayConfig["correlation"] = ptrValue(cfg.Correlation)
		displayConfig["jitter"] = ptrValue(cfg.Jitter)
		displayConfig["direction"] = cfg.Direction
		groundtruth["service"] = []string{cfg.App, cfg.TargetService}
		if pair := lookupNetworkPair(cfg, systemType); len(pair) > 0 {
			if spanNames, ok := pair["span_names"]; ok {
				groundtruth["span"] = spanNames
			}
		}
		groundtruth["metric"] = []string{string(handler.MetricNetworkLatency)}
	case "HTTPRequestDelay":
		displayConfig["injection_point"] = lookupHTTPEndpoint(cfg, systemType)
		displayConfig["delay_duration"] = ptrValue(cfg.DelayDuration)
		endpoint := lookupHTTPEndpoint(cfg, systemType)
		if len(endpoint) > 0 {
			groundtruth["service"] = []string{cfg.App, stringValue(endpoint["server_address"])}
			if spanName := stringValue(endpoint["span_name"]); spanName != "" {
				groundtruth["span"] = []string{spanName}
			}
		} else {
			groundtruth["service"] = []string{cfg.App}
		}
		groundtruth["metric"] = []string{string(handler.MetricHTTPLatency)}
	case "JVMLatency":
		displayConfig["injection_point"] = map[string]any{
			"app_name":    cfg.App,
			"class_name":  cfg.Class,
			"method_name": cfg.Method,
		}
		displayConfig["latency_duration"] = ptrValue(cfg.LatencyDuration)
		groundtruth["service"] = []string{cfg.App}
		groundtruth["function"] = []string{cfg.Class + "." + cfg.Method}
		groundtruth["metric"] = []string{string(handler.MetricNetworkLatency)}
	case "JVMRuntimeMutator":
		displayConfig["injection_point"] = map[string]any{
			"app_name":       cfg.App,
			"class_name":     cfg.Class,
			"method_name":    cfg.Method,
			"mutator_config": cfg.MutatorConfig,
		}
		if target := lookupRuntimeMutator(cfg, systemType); len(target) > 0 {
			displayConfig["injection_point"] = target
		}
		groundtruth["service"] = []string{cfg.App}
		groundtruth["function"] = []string{cfg.Class + "." + cfg.Method}
	}

	if len(payload) > 0 {
		displayConfig["apply_payload"] = payload
	}

	return &Preview{
		DisplayConfig: displayConfig,
		Groundtruth:   groundtruth,
	}
}

func lookupNetworkPair(cfg GuidedConfig, systemType systemconfig.SystemType) map[string]any {
	pairs, err := resourcelookup.GetSystemCache(systemType).GetAllNetworkPairs()
	if err != nil {
		return map[string]any{
			"source_service": cfg.App,
			"target_service": cfg.TargetService,
		}
	}
	for _, pair := range pairs {
		if pair.SourceService == cfg.App && pair.TargetService == cfg.TargetService {
			return map[string]any{
				"source_service": pair.SourceService,
				"target_service": pair.TargetService,
				"span_names":     pair.SpanNames,
			}
		}
	}
	return map[string]any{
		"source_service": cfg.App,
		"target_service": cfg.TargetService,
	}
}

func lookupHTTPEndpoint(cfg GuidedConfig, systemType systemconfig.SystemType) map[string]any {
	endpoints, err := resourcelookup.GetSystemCache(systemType).GetAllHTTPEndpoints()
	if err != nil {
		return map[string]any{
			"app_name":    cfg.App,
			"route":       cfg.Route,
			"http_method": cfg.HTTPMethod,
		}
	}
	for _, endpoint := range endpoints {
		if endpoint.AppName == cfg.App && endpoint.Route == cfg.Route && endpoint.Method == cfg.HTTPMethod {
			return map[string]any{
				"app_name":       endpoint.AppName,
				"route":          endpoint.Route,
				"http_method":    endpoint.Method,
				"server_address": endpoint.ServerAddress,
				"server_port":    endpoint.ServerPort,
				"span_name":      endpoint.SpanName,
			}
		}
	}
	return map[string]any{
		"app_name":    cfg.App,
		"route":       cfg.Route,
		"http_method": cfg.HTTPMethod,
	}
}

func lookupRuntimeMutator(cfg GuidedConfig, systemType systemconfig.SystemType) map[string]any {
	targets, err := resourcelookup.GetSystemCache(systemType).GetAllJVMRuntimeMutatorTargets()
	if err != nil {
		return map[string]any{
			"app_name":       cfg.App,
			"class_name":     cfg.Class,
			"method_name":    cfg.Method,
			"mutator_config": cfg.MutatorConfig,
		}
	}
	for _, target := range targets {
		if target.AppName == cfg.App && target.ClassName == cfg.Class && target.MethodName == cfg.Method && runtimeMutatorKey(target) == cfg.MutatorConfig {
			return map[string]any{
				"app_name":           target.AppName,
				"class_name":         target.ClassName,
				"method_name":        target.MethodName,
				"mutation_type":      target.MutationType,
				"mutation_type_name": target.MutationTypeName,
				"mutation_from":      target.MutationFrom,
				"mutation_to":        target.MutationTo,
				"mutation_strategy":  target.MutationStrategy,
				"description":        target.Description,
			}
		}
	}
	return map[string]any{
		"app_name":       cfg.App,
		"class_name":     cfg.Class,
		"method_name":    cfg.Method,
		"mutator_config": cfg.MutatorConfig,
	}
}

func ptrValue(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}

func stringValue(value any) string {
	v, _ := value.(string)
	return v
}
