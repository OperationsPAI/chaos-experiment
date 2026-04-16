package guidedcli

import (
	"context"
	"fmt"
	"strings"

	"github.com/OperationsPAI/chaos-experiment/handler"
	"github.com/OperationsPAI/chaos-experiment/internal/resourcelookup"
	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

func buildPodKill(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	apps, err := safeAppLabels(cfg.Namespace, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	appIdx := indexOf(apps, cfg.App)
	if appIdx < 0 {
		return handler.InjectionConf{}, nil, fmt.Errorf("app %q not found", cfg.App)
	}
	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	duration := normalizedDuration(cfg)
	spec := &handler.PodKillSpec{Duration: duration, System: systemIdx, AppIdx: appIdx}
	return handler.InjectionConf{PodKill: spec}, map[string]any{
		"PodKill": map[string]any{"Duration": duration, "System": systemIdx, "AppIdx": appIdx},
	}, nil
}

func buildCPUStress(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.CPULoad == nil || cfg.CPUWorker == nil {
		return handler.InjectionConf{}, nil, fmt.Errorf("cpu_load and cpu_worker are required")
	}
	containers, err := safeContainers(cfg.Namespace)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	containerIdx := -1
	for idx, container := range containers {
		if container.AppLabel == cfg.App && container.ContainerName == cfg.Container {
			containerIdx = idx
			break
		}
	}
	if containerIdx < 0 {
		return handler.InjectionConf{}, nil, fmt.Errorf("container %q not found under app %q", cfg.Container, cfg.App)
	}
	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	duration := normalizedDuration(cfg)
	spec := &handler.CPUStressChaosSpec{
		Duration:     duration,
		System:       systemIdx,
		ContainerIdx: containerIdx,
		CPULoad:      *cfg.CPULoad,
		CPUWorker:    *cfg.CPUWorker,
	}
	return handler.InjectionConf{CPUStress: spec}, map[string]any{
		"CPUStress": map[string]any{
			"Duration": duration, "System": systemIdx, "ContainerIdx": containerIdx, "CPULoad": *cfg.CPULoad, "CPUWorker": *cfg.CPUWorker,
		},
	}, nil
}

func buildNetworkDelay(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.Latency == nil || cfg.Correlation == nil || cfg.Jitter == nil || cfg.Direction == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("latency, correlation, jitter and direction are required")
	}
	pairs, err := resourcelookup.GetSystemCache(systemType).GetAllNetworkPairs()
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	pairIdx := -1
	for idx, pair := range pairs {
		if pair.SourceService == cfg.App && pair.TargetService == cfg.TargetService {
			pairIdx = idx
			break
		}
	}
	if pairIdx < 0 {
		return handler.InjectionConf{}, nil, fmt.Errorf("network pair %q -> %q not found", cfg.App, cfg.TargetService)
	}
	direction, err := directionCode(cfg.Direction)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	duration := normalizedDuration(cfg)
	spec := &handler.NetworkDelaySpec{
		Duration:       duration,
		System:         systemIdx,
		NetworkPairIdx: pairIdx,
		Latency:        *cfg.Latency,
		Correlation:    *cfg.Correlation,
		Jitter:         *cfg.Jitter,
		Direction:      direction,
	}
	return handler.InjectionConf{NetworkDelay: spec}, map[string]any{
		"NetworkDelay": map[string]any{
			"Duration": duration, "System": systemIdx, "NetworkPairIdx": pairIdx, "Latency": *cfg.Latency, "Correlation": *cfg.Correlation, "Jitter": *cfg.Jitter, "Direction": direction,
		},
	}, nil
}

func buildHTTPRequestDelay(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.DelayDuration == nil {
		return handler.InjectionConf{}, nil, fmt.Errorf("delay_duration is required")
	}
	endpoints, err := resourcelookup.GetSystemCache(systemType).GetAllHTTPEndpoints()
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	endpointIdx := -1
	for idx, endpoint := range endpoints {
		if endpoint.AppName == cfg.App && endpoint.Route == cfg.Route && endpoint.Method == cfg.HTTPMethod {
			endpointIdx = idx
			break
		}
	}
	if endpointIdx < 0 {
		return handler.InjectionConf{}, nil, fmt.Errorf("http endpoint %s %s not found under app %q", cfg.HTTPMethod, cfg.Route, cfg.App)
	}
	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	duration := normalizedDuration(cfg)
	spec := &handler.HTTPRequestDelaySpec{Duration: duration, System: systemIdx, EndpointIdx: endpointIdx, DelayDuration: *cfg.DelayDuration}
	return handler.InjectionConf{HTTPRequestDelay: spec}, map[string]any{
		"HTTPRequestDelay": map[string]any{"Duration": duration, "System": systemIdx, "EndpointIdx": endpointIdx, "DelayDuration": *cfg.DelayDuration},
	}, nil
}

func buildJVMLatency(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.LatencyDuration == nil {
		return handler.InjectionConf{}, nil, fmt.Errorf("latency_duration is required")
	}
	methods, err := resourcelookup.GetSystemCache(systemType).GetAllJVMMethods()
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	methodIdx := -1
	for idx, method := range methods {
		if method.AppName == cfg.App && method.ClassName == cfg.Class && method.MethodName == cfg.Method {
			methodIdx = idx
			break
		}
	}
	if methodIdx < 0 {
		return handler.InjectionConf{}, nil, fmt.Errorf("jvm method %s#%s not found under app %q", cfg.Class, cfg.Method, cfg.App)
	}
	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	duration := normalizedDuration(cfg)
	spec := &handler.JVMLatencySpec{Duration: duration, System: systemIdx, MethodIdx: methodIdx, LatencyDuration: *cfg.LatencyDuration}
	return handler.InjectionConf{JVMLatency: spec}, map[string]any{
		"JVMLatency": map[string]any{"Duration": duration, "System": systemIdx, "MethodIdx": methodIdx, "LatencyDuration": *cfg.LatencyDuration},
	}, nil
}

func buildJVMRuntimeMutator(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.MutatorConfig == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("mutator_config is required")
	}
	targets, err := resourcelookup.GetSystemCache(systemType).GetAllJVMRuntimeMutatorTargets()
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	targetIdx := -1
	for idx, target := range targets {
		if target.AppName == cfg.App && target.ClassName == cfg.Class && target.MethodName == cfg.Method && runtimeMutatorKey(target) == cfg.MutatorConfig {
			targetIdx = idx
			break
		}
	}
	if targetIdx < 0 {
		return handler.InjectionConf{}, nil, fmt.Errorf("runtime mutator config %q not found for %s#%s", cfg.MutatorConfig, cfg.Class, cfg.Method)
	}
	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	duration := normalizedDuration(cfg)
	spec := &handler.JVMRuntimeMutatorSpec{Duration: duration, System: systemIdx, MutatorTargetIdx: targetIdx}
	return handler.InjectionConf{JVMRuntimeMutator: spec}, map[string]any{
		"JVMRuntimeMutator": map[string]any{"Duration": duration, "System": systemIdx, "MutatorTargetIdx": targetIdx},
	}, nil
}

func normalizedDuration(cfg GuidedConfig) int {
	if cfg.Duration == nil {
		return defaultDurationMinutes
	}
	return *cfg.Duration
}

func directionCode(direction string) (int, error) {
	switch strings.ToLower(direction) {
	case "to":
		return 1, nil
	case "from":
		return 2, nil
	case "both":
		return 3, nil
	default:
		return 0, fmt.Errorf("invalid direction %q", direction)
	}
}

func indexOf(values []string, target string) int {
	for idx, value := range values {
		if value == target {
			return idx
		}
	}
	return -1
}
