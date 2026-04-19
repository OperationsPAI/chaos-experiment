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
	appIdx, systemIdx, duration, err := resolveAppLevel(ctx, cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.PodKillSpec{Duration: duration, System: systemIdx, AppIdx: appIdx}
	return handler.InjectionConf{PodKill: spec}, payload("PodKill", map[string]any{
		"Duration": duration,
		"System":   systemIdx,
		"AppIdx":   appIdx,
	}), nil
}

func buildPodFailure(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	appIdx, systemIdx, duration, err := resolveAppLevel(ctx, cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.PodFailureSpec{Duration: duration, System: systemIdx, AppIdx: appIdx}
	return handler.InjectionConf{PodFailure: spec}, payload("PodFailure", map[string]any{
		"Duration": duration,
		"System":   systemIdx,
		"AppIdx":   appIdx,
	}), nil
}

func buildJVMGarbageCollector(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	appIdx, systemIdx, duration, err := resolveAppLevel(ctx, cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.JVMGCSpec{Duration: duration, System: systemIdx, AppIdx: appIdx}
	return handler.InjectionConf{JVMGarbageCollector: spec}, payload("JVMGarbageCollector", map[string]any{
		"Duration": duration,
		"System":   systemIdx,
		"AppIdx":   appIdx,
	}), nil
}

func buildContainerKill(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	containerIdx, systemIdx, duration, err := resolveContainerLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.ContainerKillSpec{Duration: duration, System: systemIdx, ContainerIdx: containerIdx}
	return handler.InjectionConf{ContainerKill: spec}, payload("ContainerKill", map[string]any{
		"Duration":     duration,
		"System":       systemIdx,
		"ContainerIdx": containerIdx,
	}), nil
}

func buildCPUStress(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.CPULoad == nil || cfg.CPUWorker == nil {
		return handler.InjectionConf{}, nil, fmt.Errorf("cpu_load and cpu_worker are required")
	}

	containerIdx, systemIdx, duration, err := resolveContainerLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.CPUStressChaosSpec{
		Duration:     duration,
		System:       systemIdx,
		ContainerIdx: containerIdx,
		CPULoad:      *cfg.CPULoad,
		CPUWorker:    *cfg.CPUWorker,
	}
	return handler.InjectionConf{CPUStress: spec}, payload("CPUStress", map[string]any{
		"Duration":     duration,
		"System":       systemIdx,
		"ContainerIdx": containerIdx,
		"CPULoad":      *cfg.CPULoad,
		"CPUWorker":    *cfg.CPUWorker,
	}), nil
}

func buildMemoryStress(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.MemorySize == nil || cfg.MemWorker == nil {
		return handler.InjectionConf{}, nil, fmt.Errorf("memory_size and mem_worker are required")
	}

	containerIdx, systemIdx, duration, err := resolveContainerLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.MemoryStressChaosSpec{
		Duration:     duration,
		System:       systemIdx,
		ContainerIdx: containerIdx,
		MemorySize:   *cfg.MemorySize,
		MemWorker:    *cfg.MemWorker,
	}
	return handler.InjectionConf{MemoryStress: spec}, payload("MemoryStress", map[string]any{
		"Duration":     duration,
		"System":       systemIdx,
		"ContainerIdx": containerIdx,
		"MemorySize":   *cfg.MemorySize,
		"MemWorker":    *cfg.MemWorker,
	}), nil
}

func buildTimeSkew(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.TimeOffset == nil {
		return handler.InjectionConf{}, nil, fmt.Errorf("time_offset is required")
	}

	containerIdx, systemIdx, duration, err := resolveContainerLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.TimeSkewSpec{
		Duration:     duration,
		System:       systemIdx,
		ContainerIdx: containerIdx,
		TimeOffset:   *cfg.TimeOffset,
	}
	return handler.InjectionConf{TimeSkew: spec}, payload("TimeSkew", map[string]any{
		"Duration":     duration,
		"System":       systemIdx,
		"ContainerIdx": containerIdx,
		"TimeOffset":   *cfg.TimeOffset,
	}), nil
}

func buildHTTPRequestAbort(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	endpointIdx, systemIdx, duration, err := resolveEndpointLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.HTTPRequestAbortSpec{Duration: duration, System: systemIdx, EndpointIdx: endpointIdx}
	return handler.InjectionConf{HTTPRequestAbort: spec}, payload("HTTPRequestAbort", map[string]any{
		"Duration":    duration,
		"System":      systemIdx,
		"EndpointIdx": endpointIdx,
	}), nil
}

func buildHTTPResponseAbort(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	endpointIdx, systemIdx, duration, err := resolveEndpointLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.HTTPResponseAbortSpec{Duration: duration, System: systemIdx, EndpointIdx: endpointIdx}
	return handler.InjectionConf{HTTPResponseAbort: spec}, payload("HTTPResponseAbort", map[string]any{
		"Duration":    duration,
		"System":      systemIdx,
		"EndpointIdx": endpointIdx,
	}), nil
}

func buildHTTPRequestDelay(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.DelayDuration == nil {
		return handler.InjectionConf{}, nil, fmt.Errorf("delay_duration is required")
	}

	endpointIdx, systemIdx, duration, err := resolveEndpointLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.HTTPRequestDelaySpec{
		Duration:      duration,
		System:        systemIdx,
		EndpointIdx:   endpointIdx,
		DelayDuration: *cfg.DelayDuration,
	}
	return handler.InjectionConf{HTTPRequestDelay: spec}, payload("HTTPRequestDelay", map[string]any{
		"Duration":      duration,
		"System":        systemIdx,
		"EndpointIdx":   endpointIdx,
		"DelayDuration": *cfg.DelayDuration,
	}), nil
}

func buildHTTPResponseDelay(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.DelayDuration == nil {
		return handler.InjectionConf{}, nil, fmt.Errorf("delay_duration is required")
	}

	endpointIdx, systemIdx, duration, err := resolveEndpointLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.HTTPResponseDelaySpec{
		Duration:      duration,
		System:        systemIdx,
		EndpointIdx:   endpointIdx,
		DelayDuration: *cfg.DelayDuration,
	}
	return handler.InjectionConf{HTTPResponseDelay: spec}, payload("HTTPResponseDelay", map[string]any{
		"Duration":      duration,
		"System":        systemIdx,
		"EndpointIdx":   endpointIdx,
		"DelayDuration": *cfg.DelayDuration,
	}), nil
}
func buildHTTPResponseReplaceBody(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.BodyType == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("body_type is required")
	}

	endpointIdx, systemIdx, duration, err := resolveEndpointLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	bodyType, err := bodyTypeCode(cfg.BodyType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.HTTPResponseReplaceBodySpec{
		Duration:    duration,
		System:      systemIdx,
		EndpointIdx: endpointIdx,
		BodyType:    handler.ReplaceBodyType(bodyType),
	}
	return handler.InjectionConf{HTTPResponseReplaceBody: spec}, payload("HTTPResponseReplaceBody", map[string]any{
		"Duration":    duration,
		"System":      systemIdx,
		"EndpointIdx": endpointIdx,
		"BodyType":    bodyType,
	}), nil
}

func buildHTTPResponsePatchBody(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	endpointIdx, systemIdx, duration, err := resolveEndpointLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.HTTPResponsePatchBodySpec{Duration: duration, System: systemIdx, EndpointIdx: endpointIdx}
	return handler.InjectionConf{HTTPResponsePatchBody: spec}, payload("HTTPResponsePatchBody", map[string]any{
		"Duration":    duration,
		"System":      systemIdx,
		"EndpointIdx": endpointIdx,
	}), nil
}

func buildHTTPRequestReplacePath(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	endpointIdx, systemIdx, duration, err := resolveEndpointLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.HTTPRequestReplacePathSpec{Duration: duration, System: systemIdx, EndpointIdx: endpointIdx}
	return handler.InjectionConf{HTTPRequestReplacePath: spec}, payload("HTTPRequestReplacePath", map[string]any{
		"Duration":    duration,
		"System":      systemIdx,
		"EndpointIdx": endpointIdx,
	}), nil
}

func buildHTTPRequestReplaceMethod(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.ReplaceMethod == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("replace_method is required")
	}

	endpointIdx, systemIdx, duration, err := resolveEndpointLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	methodCode, err := replaceMethodCode(systemType, cfg)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.HTTPRequestReplaceMethodSpec{
		Duration:      duration,
		System:        systemIdx,
		EndpointIdx:   endpointIdx,
		ReplaceMethod: methodCode,
	}
	return handler.InjectionConf{HTTPRequestReplaceMethod: spec}, payload("HTTPRequestReplaceMethod", map[string]any{
		"Duration":      duration,
		"System":        systemIdx,
		"EndpointIdx":   endpointIdx,
		"ReplaceMethod": methodCode,
	}), nil
}

func buildHTTPResponseReplaceCode(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.StatusCode == nil {
		return handler.InjectionConf{}, nil, fmt.Errorf("status_code is required")
	}

	endpointIdx, systemIdx, duration, err := resolveEndpointLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	statusCode, err := statusCodeCode(*cfg.StatusCode)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.HTTPResponseReplaceCodeSpec{
		Duration:    duration,
		System:      systemIdx,
		EndpointIdx: endpointIdx,
		StatusCode:  handler.HTTPStatusCode(statusCode),
	}
	return handler.InjectionConf{HTTPResponseReplaceCode: spec}, payload("HTTPResponseReplaceCode", map[string]any{
		"Duration":    duration,
		"System":      systemIdx,
		"EndpointIdx": endpointIdx,
		"StatusCode":  statusCode,
	}), nil
}

func buildDNSError(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	dnsIdx, systemIdx, duration, err := resolveDNSLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.DNSErrorSpec{Duration: duration, System: systemIdx, DNSEndpointIdx: dnsIdx}
	return handler.InjectionConf{DNSError: spec}, payload("DNSError", map[string]any{
		"Duration":       duration,
		"System":         systemIdx,
		"DNSEndpointIdx": dnsIdx,
	}), nil
}

func buildDNSRandom(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	dnsIdx, systemIdx, duration, err := resolveDNSLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.DNSRandomSpec{Duration: duration, System: systemIdx, DNSEndpointIdx: dnsIdx}
	return handler.InjectionConf{DNSRandom: spec}, payload("DNSRandom", map[string]any{
		"Duration":       duration,
		"System":         systemIdx,
		"DNSEndpointIdx": dnsIdx,
	}), nil
}

func buildNetworkPartition(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.Direction == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("direction is required")
	}

	pairIdx, systemIdx, duration, direction, err := resolveNetworkLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.NetworkPartitionSpec{
		Duration:       duration,
		System:         systemIdx,
		NetworkPairIdx: pairIdx,
		Direction:      direction,
	}
	return handler.InjectionConf{NetworkPartition: spec}, payload("NetworkPartition", map[string]any{
		"Duration":       duration,
		"System":         systemIdx,
		"NetworkPairIdx": pairIdx,
		"Direction":      direction,
	}), nil
}

func buildNetworkDelay(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.Latency == nil || cfg.Correlation == nil || cfg.Jitter == nil || cfg.Direction == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("latency, correlation, jitter and direction are required")
	}

	pairIdx, systemIdx, duration, direction, err := resolveNetworkLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.NetworkDelaySpec{
		Duration:       duration,
		System:         systemIdx,
		NetworkPairIdx: pairIdx,
		Latency:        *cfg.Latency,
		Correlation:    *cfg.Correlation,
		Jitter:         *cfg.Jitter,
		Direction:      direction,
	}
	return handler.InjectionConf{NetworkDelay: spec}, payload("NetworkDelay", map[string]any{
		"Duration":       duration,
		"System":         systemIdx,
		"NetworkPairIdx": pairIdx,
		"Latency":        *cfg.Latency,
		"Correlation":    *cfg.Correlation,
		"Jitter":         *cfg.Jitter,
		"Direction":      direction,
	}), nil
}
func buildNetworkLoss(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.Loss == nil || cfg.Correlation == nil || cfg.Direction == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("loss, correlation and direction are required")
	}

	pairIdx, systemIdx, duration, direction, err := resolveNetworkLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.NetworkLossSpec{
		Duration:       duration,
		System:         systemIdx,
		NetworkPairIdx: pairIdx,
		Loss:           *cfg.Loss,
		Correlation:    *cfg.Correlation,
		Direction:      direction,
	}
	return handler.InjectionConf{NetworkLoss: spec}, payload("NetworkLoss", map[string]any{
		"Duration":       duration,
		"System":         systemIdx,
		"NetworkPairIdx": pairIdx,
		"Loss":           *cfg.Loss,
		"Correlation":    *cfg.Correlation,
		"Direction":      direction,
	}), nil
}

func buildNetworkDuplicate(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.Duplicate == nil || cfg.Correlation == nil || cfg.Direction == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("duplicate, correlation and direction are required")
	}

	pairIdx, systemIdx, duration, direction, err := resolveNetworkLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.NetworkDuplicateSpec{
		Duration:       duration,
		System:         systemIdx,
		NetworkPairIdx: pairIdx,
		Duplicate:      *cfg.Duplicate,
		Correlation:    *cfg.Correlation,
		Direction:      direction,
	}
	return handler.InjectionConf{NetworkDuplicate: spec}, payload("NetworkDuplicate", map[string]any{
		"Duration":       duration,
		"System":         systemIdx,
		"NetworkPairIdx": pairIdx,
		"Duplicate":      *cfg.Duplicate,
		"Correlation":    *cfg.Correlation,
		"Direction":      direction,
	}), nil
}

func buildNetworkCorrupt(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.Corrupt == nil || cfg.Correlation == nil || cfg.Direction == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("corrupt, correlation and direction are required")
	}

	pairIdx, systemIdx, duration, direction, err := resolveNetworkLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.NetworkCorruptSpec{
		Duration:       duration,
		System:         systemIdx,
		NetworkPairIdx: pairIdx,
		Corrupt:        *cfg.Corrupt,
		Correlation:    *cfg.Correlation,
		Direction:      direction,
	}
	return handler.InjectionConf{NetworkCorrupt: spec}, payload("NetworkCorrupt", map[string]any{
		"Duration":       duration,
		"System":         systemIdx,
		"NetworkPairIdx": pairIdx,
		"Corrupt":        *cfg.Corrupt,
		"Correlation":    *cfg.Correlation,
		"Direction":      direction,
	}), nil
}

func buildNetworkBandwidth(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.Rate == nil || cfg.Limit == nil || cfg.Buffer == nil || cfg.Direction == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("rate, limit, buffer and direction are required")
	}

	pairIdx, systemIdx, duration, direction, err := resolveNetworkLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.NetworkBandwidthSpec{
		Duration:       duration,
		System:         systemIdx,
		NetworkPairIdx: pairIdx,
		Rate:           *cfg.Rate,
		Limit:          *cfg.Limit,
		Buffer:         *cfg.Buffer,
		Direction:      direction,
	}
	return handler.InjectionConf{NetworkBandwidth: spec}, payload("NetworkBandwidth", map[string]any{
		"Duration":       duration,
		"System":         systemIdx,
		"NetworkPairIdx": pairIdx,
		"Rate":           *cfg.Rate,
		"Limit":          *cfg.Limit,
		"Buffer":         *cfg.Buffer,
		"Direction":      direction,
	}), nil
}

func buildJVMLatency(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.LatencyDuration == nil {
		return handler.InjectionConf{}, nil, fmt.Errorf("latency_duration is required")
	}

	methodIdx, systemIdx, duration, err := resolveMethodLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.JVMLatencySpec{
		Duration:        duration,
		System:          systemIdx,
		MethodIdx:       methodIdx,
		LatencyDuration: *cfg.LatencyDuration,
	}
	return handler.InjectionConf{JVMLatency: spec}, payload("JVMLatency", map[string]any{
		"Duration":        duration,
		"System":          systemIdx,
		"MethodIdx":       methodIdx,
		"LatencyDuration": *cfg.LatencyDuration,
	}), nil
}

func buildJVMReturn(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.ReturnType == "" || cfg.ReturnValueOpt == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("return_type and return_value_opt are required")
	}

	methodIdx, systemIdx, duration, err := resolveMethodLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	returnType, err := returnTypeCode(cfg.ReturnType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}
	returnValueOpt, err := returnValueOptCode(cfg.ReturnValueOpt)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.JVMReturnSpec{
		Duration:       duration,
		System:         systemIdx,
		MethodIdx:      methodIdx,
		ReturnType:     handler.JVMReturnType(returnType),
		ReturnValueOpt: returnValueOpt,
	}
	return handler.InjectionConf{JVMReturn: spec}, payload("JVMReturn", map[string]any{
		"Duration":       duration,
		"System":         systemIdx,
		"MethodIdx":      methodIdx,
		"ReturnType":     returnType,
		"ReturnValueOpt": returnValueOpt,
	}), nil
}

func buildJVMException(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.ExceptionOpt == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("exception_opt is required")
	}

	methodIdx, systemIdx, duration, err := resolveMethodLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	exceptionOpt, err := exceptionOptCode(cfg.ExceptionOpt)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.JVMExceptionSpec{
		Duration:     duration,
		System:       systemIdx,
		MethodIdx:    methodIdx,
		ExceptionOpt: exceptionOpt,
	}
	return handler.InjectionConf{JVMException: spec}, payload("JVMException", map[string]any{
		"Duration":     duration,
		"System":       systemIdx,
		"MethodIdx":    methodIdx,
		"ExceptionOpt": exceptionOpt,
	}), nil
}
func buildJVMCPUStress(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.CPUCount == nil {
		return handler.InjectionConf{}, nil, fmt.Errorf("cpu_count is required")
	}

	methodIdx, systemIdx, duration, err := resolveMethodLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.JVMCPUStressSpec{
		Duration:  duration,
		System:    systemIdx,
		MethodIdx: methodIdx,
		CPUCount:  *cfg.CPUCount,
	}
	return handler.InjectionConf{JVMCPUStress: spec}, payload("JVMCPUStress", map[string]any{
		"Duration":  duration,
		"System":    systemIdx,
		"MethodIdx": methodIdx,
		"CPUCount":  *cfg.CPUCount,
	}), nil
}

func buildJVMMemoryStress(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.MemType == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("mem_type is required")
	}

	methodIdx, systemIdx, duration, err := resolveMethodLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	memType, err := memTypeCode(cfg.MemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.JVMMemoryStressSpec{
		Duration:  duration,
		System:    systemIdx,
		MethodIdx: methodIdx,
		MemType:   handler.JVMMemoryType(memType),
	}
	return handler.InjectionConf{JVMMemoryStress: spec}, payload("JVMMemoryStress", map[string]any{
		"Duration":  duration,
		"System":    systemIdx,
		"MethodIdx": methodIdx,
		"MemType":   memType,
	}), nil
}

func buildJVMMySQLLatency(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.LatencyMs == nil {
		return handler.InjectionConf{}, nil, fmt.Errorf("latency_ms is required")
	}

	databaseIdx, systemIdx, duration, err := resolveDatabaseLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.JVMMySQLLatencySpec{
		Duration:    duration,
		System:      systemIdx,
		DatabaseIdx: databaseIdx,
		LatencyMs:   *cfg.LatencyMs,
	}
	return handler.InjectionConf{JVMMySQLLatency: spec}, payload("JVMMySQLLatency", map[string]any{
		"Duration":    duration,
		"System":      systemIdx,
		"DatabaseIdx": databaseIdx,
		"LatencyMs":   *cfg.LatencyMs,
	}), nil
}

func buildJVMMySQLException(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	databaseIdx, systemIdx, duration, err := resolveDatabaseLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.JVMMySQLExceptionSpec{
		Duration:    duration,
		System:      systemIdx,
		DatabaseIdx: databaseIdx,
	}
	return handler.InjectionConf{JVMMySQLException: spec}, payload("JVMMySQLException", map[string]any{
		"Duration":    duration,
		"System":      systemIdx,
		"DatabaseIdx": databaseIdx,
	}), nil
}

func buildJVMRuntimeMutator(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (handler.InjectionConf, map[string]any, error) {
	if cfg.MutatorConfig == "" {
		return handler.InjectionConf{}, nil, fmt.Errorf("mutator_config is required")
	}

	targetIdx, systemIdx, duration, err := resolveMutatorLevel(cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, nil, err
	}

	spec := &handler.JVMRuntimeMutatorSpec{Duration: duration, System: systemIdx, MutatorTargetIdx: targetIdx}
	return handler.InjectionConf{JVMRuntimeMutator: spec}, payload("JVMRuntimeMutator", map[string]any{
		"Duration":         duration,
		"System":           systemIdx,
		"MutatorTargetIdx": targetIdx,
	}), nil
}

func resolveAppLevel(ctx context.Context, cfg GuidedConfig, systemType systemconfig.SystemType) (int, int, int, error) {
	apps, err := safeAppLabels(cfg.Namespace, systemType)
	if err != nil {
		return 0, 0, 0, err
	}
	appIdx := indexOf(apps, cfg.App)
	if appIdx < 0 {
		return 0, 0, 0, fmt.Errorf("app %q not found", cfg.App)
	}
	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return 0, 0, 0, err
	}
	return appIdx, systemIdx, normalizedDuration(cfg), nil
}

func resolveContainerLevel(cfg GuidedConfig, systemType systemconfig.SystemType) (int, int, int, error) {
	containers, err := safeContainers(cfg.Namespace)
	if err != nil {
		return 0, 0, 0, err
	}

	containerIdx := -1
	for idx, container := range containers {
		if container.AppLabel == cfg.App && container.ContainerName == cfg.Container {
			containerIdx = idx
			break
		}
	}
	if containerIdx < 0 {
		return 0, 0, 0, fmt.Errorf("container %q not found under app %q", cfg.Container, cfg.App)
	}

	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return 0, 0, 0, err
	}
	return containerIdx, systemIdx, normalizedDuration(cfg), nil
}

func resolveEndpointLevel(cfg GuidedConfig, systemType systemconfig.SystemType) (int, int, int, error) {
	endpoints, err := resourcelookup.GetSystemCache(systemType).GetAllHTTPEndpoints()
	if err != nil {
		return 0, 0, 0, err
	}

	endpointIdx := -1
	for idx, endpoint := range endpoints {
		if endpoint.AppName == cfg.App && endpoint.Route == cfg.Route && endpoint.Method == cfg.HTTPMethod {
			endpointIdx = idx
			break
		}
	}
	if endpointIdx < 0 {
		return 0, 0, 0, fmt.Errorf("http endpoint %s %s not found under app %q", cfg.HTTPMethod, cfg.Route, cfg.App)
	}

	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return 0, 0, 0, err
	}
	return endpointIdx, systemIdx, normalizedDuration(cfg), nil
}

func resolveDNSLevel(cfg GuidedConfig, systemType systemconfig.SystemType) (int, int, int, error) {
	endpoints, err := resourcelookup.GetSystemCache(systemType).GetAllDNSEndpoints()
	if err != nil {
		return 0, 0, 0, err
	}

	dnsIdx := -1
	for idx, endpoint := range endpoints {
		if endpoint.AppName == cfg.App && endpoint.Domain == cfg.Domain {
			dnsIdx = idx
			break
		}
	}
	if dnsIdx < 0 {
		return 0, 0, 0, fmt.Errorf("dns domain %q not found under app %q", cfg.Domain, cfg.App)
	}

	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return 0, 0, 0, err
	}
	return dnsIdx, systemIdx, normalizedDuration(cfg), nil
}

func resolveNetworkLevel(cfg GuidedConfig, systemType systemconfig.SystemType) (int, int, int, int, error) {
	pairs, err := resourcelookup.GetSystemCache(systemType).GetAllNetworkPairs()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	pairIdx := -1
	for idx, pair := range pairs {
		if pair.SourceService == cfg.App && pair.TargetService == cfg.TargetService {
			pairIdx = idx
			break
		}
	}
	if pairIdx < 0 {
		return 0, 0, 0, 0, fmt.Errorf("network pair %q -> %q not found", cfg.App, cfg.TargetService)
	}

	direction, err := directionCode(cfg.Direction)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return pairIdx, systemIdx, normalizedDuration(cfg), direction, nil
}

func resolveMethodLevel(cfg GuidedConfig, systemType systemconfig.SystemType) (int, int, int, error) {
	methods, err := resourcelookup.GetSystemCache(systemType).GetAllJVMMethods()
	if err != nil {
		return 0, 0, 0, err
	}

	methodIdx := -1
	for idx, method := range methods {
		if method.AppName == cfg.App && method.ClassName == cfg.Class && method.MethodName == cfg.Method {
			methodIdx = idx
			break
		}
	}
	if methodIdx < 0 {
		return 0, 0, 0, fmt.Errorf("jvm method %s#%s not found under app %q", cfg.Class, cfg.Method, cfg.App)
	}

	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return 0, 0, 0, err
	}
	return methodIdx, systemIdx, normalizedDuration(cfg), nil
}
func resolveDatabaseLevel(cfg GuidedConfig, systemType systemconfig.SystemType) (int, int, int, error) {
	operations, err := resourcelookup.GetSystemCache(systemType).GetAllDatabaseOperations()
	if err != nil {
		return 0, 0, 0, err
	}

	databaseIdx := -1
	for idx, op := range operations {
		if op.AppName == cfg.App && op.DBName == cfg.Database && op.TableName == cfg.Table && strings.EqualFold(op.OperationType, cfg.Operation) {
			databaseIdx = idx
			break
		}
	}
	if databaseIdx < 0 {
		return 0, 0, 0, fmt.Errorf("database operation %s/%s/%s not found under app %q", cfg.Database, cfg.Table, cfg.Operation, cfg.App)
	}

	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return 0, 0, 0, err
	}
	return databaseIdx, systemIdx, normalizedDuration(cfg), nil
}

func resolveMutatorLevel(cfg GuidedConfig, systemType systemconfig.SystemType) (int, int, int, error) {
	targets, err := resourcelookup.GetSystemCache(systemType).GetAllJVMRuntimeMutatorTargets()
	if err != nil {
		return 0, 0, 0, err
	}

	targetIdx := -1
	for idx, target := range targets {
		if target.AppName == cfg.App && target.ClassName == cfg.Class && target.MethodName == cfg.Method && runtimeMutatorKey(target) == cfg.MutatorConfig {
			targetIdx = idx
			break
		}
	}
	if targetIdx < 0 {
		return 0, 0, 0, fmt.Errorf("runtime mutator config %q not found for %s#%s", cfg.MutatorConfig, cfg.Class, cfg.Method)
	}

	systemIdx, err := systemTypeIndex(systemType.String())
	if err != nil {
		return 0, 0, 0, err
	}
	return targetIdx, systemIdx, normalizedDuration(cfg), nil
}

func payload(name string, fields map[string]any) map[string]any {
	return map[string]any{name: fields}
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

func returnTypeCode(value string) (int, error) {
	switch strings.ToLower(value) {
	case "string":
		return 1, nil
	case "int":
		return 2, nil
	default:
		return 0, fmt.Errorf("invalid return_type %q", value)
	}
}

func returnValueOptCode(value string) (int, error) {
	switch strings.ToLower(value) {
	case "default":
		return 0, nil
	case "random":
		return 1, nil
	default:
		return 0, fmt.Errorf("invalid return_value_opt %q", value)
	}
}

func exceptionOptCode(value string) (int, error) {
	switch strings.ToLower(value) {
	case "default":
		return 0, nil
	case "random":
		return 1, nil
	default:
		return 0, fmt.Errorf("invalid exception_opt %q", value)
	}
}

func memTypeCode(value string) (int, error) {
	switch strings.ToLower(value) {
	case "heap":
		return 1, nil
	case "stack":
		return 2, nil
	default:
		return 0, fmt.Errorf("invalid mem_type %q", value)
	}
}

func bodyTypeCode(value string) (int, error) {
	switch strings.ToLower(value) {
	case "empty":
		return 0, nil
	case "random":
		return 1, nil
	default:
		return 0, fmt.Errorf("invalid body_type %q", value)
	}
}

func statusCodeCode(value int) (int, error) {
	switch value {
	case 400:
		return 0, nil
	case 401:
		return 1, nil
	case 403:
		return 2, nil
	case 404:
		return 3, nil
	case 405:
		return 4, nil
	case 408:
		return 5, nil
	case 500:
		return 6, nil
	case 502:
		return 7, nil
	case 503:
		return 8, nil
	case 504:
		return 9, nil
	default:
		return 0, fmt.Errorf("invalid status_code %d", value)
	}
}

func replaceMethodCode(systemType systemconfig.SystemType, cfg GuidedConfig) (int, error) {
	endpoints, err := resourcelookup.GetSystemCache(systemType).GetAllHTTPEndpoints()
	if err != nil {
		return 0, err
	}

	for _, endpoint := range endpoints {
		if endpoint.AppName == cfg.App && endpoint.Route == cfg.Route && endpoint.Method == cfg.HTTPMethod {
			filtered := handler.GetFilteredHTTPMethods(endpoint.Method)
			targetMethod := strings.ToUpper(cfg.ReplaceMethod)
			for idx, method := range filtered {
				if handler.GetHTTPMethodName(method) == targetMethod {
					return idx, nil
				}
			}
			return 0, fmt.Errorf("invalid replace_method %q for endpoint method %q", cfg.ReplaceMethod, endpoint.Method)
		}
	}

	return 0, fmt.Errorf("http endpoint %s %s not found under app %q", cfg.HTTPMethod, cfg.Route, cfg.App)
}

func indexOf(values []string, target string) int {
	for idx, value := range values {
		if value == target {
			return idx
		}
	}
	return -1
}
