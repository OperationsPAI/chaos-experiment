package guidedcli

import (
	"context"
	"fmt"

	"github.com/OperationsPAI/chaos-experiment/handler"
	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

// BuildInjection resolves a finalized GuidedConfig into a handler.InjectionConf
// suitable for passing to handler.BatchCreate. The config must be in a
// ready-to-apply state (all required fields filled, chaos type selected).
// If the config is incomplete, returns an error describing what's missing.
//
// Unlike Resolve with Apply=true, this does NOT call BatchCreate — the caller
// is responsible for the actual CRD creation. This is the consumer-side API
// used when the wire carries a GuidedConfig but the inject loop runs elsewhere.
func BuildInjection(ctx context.Context, cfg GuidedConfig) (handler.InjectionConf, handler.SystemType, error) {
	if err := normalizeSystemSelection(&cfg); err != nil {
		return handler.InjectionConf{}, handler.SystemType(""), fmt.Errorf("guided config not ready to apply: %w", err)
	}
	if cfg.System == "" {
		return handler.InjectionConf{}, handler.SystemType(""), fmt.Errorf("guided config not ready to apply: system is required")
	}

	systemType, err := systemconfig.ParseSystemType(cfg.SystemType)
	if err != nil {
		return handler.InjectionConf{}, handler.SystemType(""), fmt.Errorf("guided config not ready to apply: %w", err)
	}
	if err := systemconfig.SetCurrentSystem(systemType); err != nil {
		return handler.InjectionConf{}, handler.SystemType(""), fmt.Errorf("guided config not ready to apply: %w", err)
	}

	if cfg.App == "" {
		return handler.InjectionConf{}, handler.SystemType(""), fmt.Errorf("guided config not ready to apply: app is required")
	}
	if cfg.ChaosType == "" {
		return handler.InjectionConf{}, handler.SystemType(""), fmt.Errorf("guided config not ready to apply: chaos_type is required")
	}

	builder, ok := lookupBuilder(cfg.ChaosType)
	if !ok {
		return handler.InjectionConf{}, handler.SystemType(""), fmt.Errorf("guided config not ready to apply: chaos type %q is not supported", cfg.ChaosType)
	}

	conf, _, err := builder(ctx, cfg, systemType)
	if err != nil {
		return handler.InjectionConf{}, handler.SystemType(""), fmt.Errorf("build %s injection: %w", cfg.ChaosType, err)
	}
	return conf, handler.SystemType(systemType), nil
}

// lookupBuilder returns the builder registered for the given chaos type.
// Mirrors the switch in Resolve; kept independent so BuildInjection can skip
// the dispatch helpers that wrap builders for interactive field requests.
func lookupBuilder(chaosType string) (buildFunc, bool) {
	switch chaosType {
	case "PodKill":
		return buildPodKill, true
	case "PodFailure":
		return buildPodFailure, true
	case "ContainerKill":
		return buildContainerKill, true
	case "CPUStress":
		return buildCPUStress, true
	case "MemoryStress":
		return buildMemoryStress, true
	case "TimeSkew":
		return buildTimeSkew, true
	case "HTTPRequestAbort":
		return buildHTTPRequestAbort, true
	case "HTTPResponseAbort":
		return buildHTTPResponseAbort, true
	case "NetworkDelay":
		return buildNetworkDelay, true
	case "NetworkPartition":
		return buildNetworkPartition, true
	case "NetworkLoss":
		return buildNetworkLoss, true
	case "NetworkDuplicate":
		return buildNetworkDuplicate, true
	case "NetworkCorrupt":
		return buildNetworkCorrupt, true
	case "NetworkBandwidth":
		return buildNetworkBandwidth, true
	case "HTTPRequestDelay":
		return buildHTTPRequestDelay, true
	case "HTTPResponseDelay":
		return buildHTTPResponseDelay, true
	case "HTTPResponseReplaceBody":
		return buildHTTPResponseReplaceBody, true
	case "HTTPResponsePatchBody":
		return buildHTTPResponsePatchBody, true
	case "HTTPRequestReplacePath":
		return buildHTTPRequestReplacePath, true
	case "HTTPRequestReplaceMethod":
		return buildHTTPRequestReplaceMethod, true
	case "HTTPResponseReplaceCode":
		return buildHTTPResponseReplaceCode, true
	case "DNSError":
		return buildDNSError, true
	case "DNSRandom":
		return buildDNSRandom, true
	case "JVMLatency":
		return buildJVMLatency, true
	case "JVMReturn":
		return buildJVMReturn, true
	case "JVMException":
		return buildJVMException, true
	case "JVMGarbageCollector":
		return buildJVMGarbageCollector, true
	case "JVMCPUStress":
		return buildJVMCPUStress, true
	case "JVMMemoryStress":
		return buildJVMMemoryStress, true
	case "JVMMySQLLatency":
		return buildJVMMySQLLatency, true
	case "JVMMySQLException":
		return buildJVMMySQLException, true
	case "JVMRuntimeMutator":
		return buildJVMRuntimeMutator, true
	default:
		return nil, false
	}
}
