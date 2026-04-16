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
	case "PodKill", "PodFailure", "JVMGarbageCollector":
		displayConfig["injection_point"] = map[string]any{"app_name": cfg.App}
		groundtruth["service"] = []string{cfg.App}
	case "ContainerKill", "CPUStress", "MemoryStress", "TimeSkew":
		displayConfig["injection_point"] = map[string]any{
			"app_name":       cfg.App,
			"container_name": cfg.Container,
			"namespace":      cfg.Namespace,
		}
		switch cfg.ChaosType {
		case "CPUStress":
			displayConfig["cpu_load"] = ptrValue(cfg.CPULoad)
			displayConfig["cpu_worker"] = ptrValue(cfg.CPUWorker)
			groundtruth["metric"] = []string{string(handler.MetricCPU)}
		case "MemoryStress":
			displayConfig["memory_size"] = ptrValue(cfg.MemorySize)
			displayConfig["mem_worker"] = ptrValue(cfg.MemWorker)
			groundtruth["metric"] = []string{string(handler.MetricMemory)}
		case "TimeSkew":
			displayConfig["time_offset"] = ptrValue(cfg.TimeOffset)
		}
		groundtruth["service"] = []string{cfg.App}
		if cfg.Container != "" {
			groundtruth["container"] = []string{cfg.Container}
		}
	case "NetworkDelay", "NetworkPartition", "NetworkLoss", "NetworkDuplicate", "NetworkCorrupt", "NetworkBandwidth":
		displayConfig["injection_point"] = lookupNetworkPair(cfg, systemType)
		displayConfig["direction"] = cfg.Direction
		switch cfg.ChaosType {
		case "NetworkDelay":
			displayConfig["latency"] = ptrValue(cfg.Latency)
			displayConfig["correlation"] = ptrValue(cfg.Correlation)
			displayConfig["jitter"] = ptrValue(cfg.Jitter)
		case "NetworkLoss":
			displayConfig["loss"] = ptrValue(cfg.Loss)
			displayConfig["correlation"] = ptrValue(cfg.Correlation)
		case "NetworkDuplicate":
			displayConfig["duplicate"] = ptrValue(cfg.Duplicate)
			displayConfig["correlation"] = ptrValue(cfg.Correlation)
		case "NetworkCorrupt":
			displayConfig["corrupt"] = ptrValue(cfg.Corrupt)
			displayConfig["correlation"] = ptrValue(cfg.Correlation)
		case "NetworkBandwidth":
			displayConfig["rate"] = ptrValue(cfg.Rate)
			displayConfig["limit"] = ptrValue(cfg.Limit)
			displayConfig["buffer"] = ptrValue(cfg.Buffer)
		}
		groundtruth["service"] = []string{cfg.App, cfg.TargetService}
		if pair := lookupNetworkPair(cfg, systemType); len(pair) > 0 {
			if spanNames, ok := pair["span_names"]; ok {
				groundtruth["span"] = spanNames
			}
		}
		groundtruth["metric"] = []string{string(handler.MetricNetworkLatency)}
	case "HTTPRequestAbort", "HTTPResponseAbort", "HTTPRequestDelay", "HTTPResponseDelay", "HTTPResponseReplaceBody", "HTTPResponsePatchBody", "HTTPRequestReplacePath", "HTTPRequestReplaceMethod", "HTTPResponseReplaceCode":
		displayConfig["injection_point"] = lookupHTTPEndpoint(cfg, systemType)
		switch cfg.ChaosType {
		case "HTTPRequestDelay", "HTTPResponseDelay":
			displayConfig["delay_duration"] = ptrValue(cfg.DelayDuration)
		case "HTTPResponseReplaceBody":
			displayConfig["body_type"] = cfg.BodyType
		case "HTTPRequestReplaceMethod":
			displayConfig["replace_method"] = cfg.ReplaceMethod
		case "HTTPResponseReplaceCode":
			displayConfig["status_code"] = ptrValue(cfg.StatusCode)
		}
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
	case "DNSError", "DNSRandom":
		displayConfig["injection_point"] = lookupDNSEndpoint(cfg, systemType)
		groundtruth["service"] = []string{cfg.App, cfg.Domain}
		if endpoint := lookupDNSEndpoint(cfg, systemType); len(endpoint) > 0 {
			if spanNames, ok := endpoint["span_names"]; ok {
				groundtruth["span"] = spanNames
			}
		}
	case "JVMLatency", "JVMReturn", "JVMException", "JVMCPUStress", "JVMMemoryStress":
		displayConfig["injection_point"] = map[string]any{
			"app_name":    cfg.App,
			"class_name":  cfg.Class,
			"method_name": cfg.Method,
		}
		switch cfg.ChaosType {
		case "JVMLatency":
			displayConfig["latency_duration"] = ptrValue(cfg.LatencyDuration)
			groundtruth["metric"] = []string{string(handler.MetricNetworkLatency)}
		case "JVMReturn":
			displayConfig["return_type"] = cfg.ReturnType
			displayConfig["return_value_opt"] = cfg.ReturnValueOpt
		case "JVMException":
			displayConfig["exception_opt"] = cfg.ExceptionOpt
		case "JVMCPUStress":
			displayConfig["cpu_count"] = ptrValue(cfg.CPUCount)
			groundtruth["metric"] = []string{string(handler.MetricCPU)}
		case "JVMMemoryStress":
			displayConfig["mem_type"] = cfg.MemType
			groundtruth["metric"] = []string{string(handler.MetricMemory)}
		}
		groundtruth["service"] = []string{cfg.App}
		groundtruth["function"] = []string{cfg.Class + "." + cfg.Method}
	case "JVMMySQLLatency", "JVMMySQLException":
		displayConfig["injection_point"] = lookupDatabaseOperation(cfg, systemType)
		if cfg.ChaosType == "JVMMySQLLatency" {
			displayConfig["latency_ms"] = ptrValue(cfg.LatencyMs)
			groundtruth["metric"] = []string{string(handler.MetricSQLLatency)}
		}
		groundtruth["service"] = []string{cfg.App}
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
		return map[string]any{"source_service": cfg.App, "target_service": cfg.TargetService}
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
	return map[string]any{"source_service": cfg.App, "target_service": cfg.TargetService}
}

func lookupHTTPEndpoint(cfg GuidedConfig, systemType systemconfig.SystemType) map[string]any {
	endpoints, err := resourcelookup.GetSystemCache(systemType).GetAllHTTPEndpoints()
	if err != nil {
		return map[string]any{"app_name": cfg.App, "route": cfg.Route, "http_method": cfg.HTTPMethod}
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
	return map[string]any{"app_name": cfg.App, "route": cfg.Route, "http_method": cfg.HTTPMethod}
}

func lookupDNSEndpoint(cfg GuidedConfig, systemType systemconfig.SystemType) map[string]any {
	endpoints, err := resourcelookup.GetSystemCache(systemType).GetAllDNSEndpoints()
	if err != nil {
		return map[string]any{"app_name": cfg.App, "domain": cfg.Domain}
	}
	for _, endpoint := range endpoints {
		if endpoint.AppName == cfg.App && endpoint.Domain == cfg.Domain {
			return map[string]any{
				"app_name":   endpoint.AppName,
				"domain":     endpoint.Domain,
				"span_names": endpoint.SpanNames,
			}
		}
	}
	return map[string]any{"app_name": cfg.App, "domain": cfg.Domain}
}

func lookupDatabaseOperation(cfg GuidedConfig, systemType systemconfig.SystemType) map[string]any {
	ops, err := resourcelookup.GetSystemCache(systemType).GetAllDatabaseOperations()
	if err != nil {
		return map[string]any{"app_name": cfg.App, "database": cfg.Database, "table": cfg.Table, "operation": cfg.Operation}
	}
	for _, op := range ops {
		if op.AppName == cfg.App && op.DBName == cfg.Database && op.TableName == cfg.Table && op.OperationType == cfg.Operation {
			return map[string]any{
				"app_name":  op.AppName,
				"database":  op.DBName,
				"table":     op.TableName,
				"operation": op.OperationType,
			}
		}
	}
	return map[string]any{"app_name": cfg.App, "database": cfg.Database, "table": cfg.Table, "operation": cfg.Operation}
}

func lookupRuntimeMutator(cfg GuidedConfig, systemType systemconfig.SystemType) map[string]any {
	targets, err := resourcelookup.GetSystemCache(systemType).GetAllJVMRuntimeMutatorTargets()
	if err != nil {
		return map[string]any{"app_name": cfg.App, "class_name": cfg.Class, "method_name": cfg.Method, "mutator_config": cfg.MutatorConfig}
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
	return map[string]any{"app_name": cfg.App, "class_name": cfg.Class, "method_name": cfg.Method, "mutator_config": cfg.MutatorConfig}
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
