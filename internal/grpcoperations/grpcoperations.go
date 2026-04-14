// Package grpcoperations provides a system-aware routing layer for gRPC operation data.
// This package delegates to registered providers instead of hard-coded switch statements.
package grpcoperations

import (
	"sort"
	"strings"

	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"

	hsgrpc "github.com/OperationsPAI/chaos-experiment/internal/hs/grpcoperations"
	mediagrpc "github.com/OperationsPAI/chaos-experiment/internal/media/grpcoperations"
	obgrpc "github.com/OperationsPAI/chaos-experiment/internal/ob/grpcoperations"
	oteldemogrpc "github.com/OperationsPAI/chaos-experiment/internal/oteldemo/grpcoperations"
	sngrpc "github.com/OperationsPAI/chaos-experiment/internal/sn/grpcoperations"
	sockshopgrpc "github.com/OperationsPAI/chaos-experiment/internal/sockshop/grpcoperations"
	teastoregrpc "github.com/OperationsPAI/chaos-experiment/internal/teastore/grpcoperations"
)

// GRPCOperation represents a gRPC operation from ClickHouse analysis.
type GRPCOperation struct {
	ServiceName   string
	RPCSystem     string
	RPCService    string
	RPCMethod     string
	StatusCode    string
	ServerAddress string
	ServerPort    string
	SpanKind      string
}

type staticGRPCOperationProvider struct {
	operations map[string][]GRPCOperation
}

func init() {
	registry := systemconfig.GetRegistry()
	registry.RegisterGRPCOperationProvider(systemconfig.SystemOtelDemo, newStaticGRPCOperationProvider(convertOtelDemoOperationMap(oteldemogrpc.GRPCOperations)))
	registry.RegisterGRPCOperationProvider(systemconfig.SystemMediaMicroservices, newStaticGRPCOperationProvider(convertMediaOperationMap(mediagrpc.GRPCOperations)))
	registry.RegisterGRPCOperationProvider(systemconfig.SystemHotelReservation, newStaticGRPCOperationProvider(convertHSOperationMap(hsgrpc.GRPCOperations)))
	registry.RegisterGRPCOperationProvider(systemconfig.SystemSocialNetwork, newStaticGRPCOperationProvider(convertSNOperationMap(sngrpc.GRPCOperations)))
	registry.RegisterGRPCOperationProvider(systemconfig.SystemOnlineBoutique, newStaticGRPCOperationProvider(convertOBOperationMap(obgrpc.GRPCOperations)))
	registry.RegisterGRPCOperationProvider(systemconfig.SystemSockShop, newStaticGRPCOperationProvider(convertSockShopOperationMap(sockshopgrpc.GRPCOperations)))
	registry.RegisterGRPCOperationProvider(systemconfig.SystemTeaStore, newStaticGRPCOperationProvider(convertTeaStoreOperationMap(teastoregrpc.GRPCOperations)))
}

func newStaticGRPCOperationProvider(operations map[string][]GRPCOperation) systemconfig.GRPCOperationProvider {
	return &staticGRPCOperationProvider{operations: operations}
}

func (p *staticGRPCOperationProvider) GetServiceNames() []string {
	services := make([]string, 0, len(p.operations))
	for service := range p.operations {
		services = append(services, service)
	}
	sort.Strings(services)
	return services
}

func (p *staticGRPCOperationProvider) GetOperationsByService(serviceName string) []systemconfig.GRPCOperationData {
	operations := p.operations[serviceName]
	result := make([]systemconfig.GRPCOperationData, len(operations))
	for i, operation := range operations {
		result[i] = systemconfig.GRPCOperationData{
			ServiceName:    operation.ServiceName,
			RPCSystem:      operation.RPCSystem,
			RPCService:     operation.RPCService,
			RPCMethod:      operation.RPCMethod,
			GRPCStatusCode: operation.StatusCode,
			ServerAddress:  operation.ServerAddress,
			ServerPort:     operation.ServerPort,
			SpanKind:       operation.SpanKind,
		}
	}
	return result
}

// GetOperationsByService returns all gRPC operations for a service based on the current system.
func GetOperationsByService(serviceName string) []GRPCOperation {
	data, err := systemconfig.GetMetadataStore().GetGRPCOperations(string(systemconfig.GetCurrentSystem()), serviceName)
	if err == nil && len(data) > 0 {
		result := make([]GRPCOperation, len(data))
		for i, operation := range data {
			result[i] = GRPCOperation{
				ServiceName:   operation.ServiceName,
				RPCSystem:     operation.RPCSystem,
				RPCService:    operation.RPCService,
				RPCMethod:     operation.RPCMethod,
				StatusCode:    operation.GRPCStatusCode,
				ServerAddress: operation.ServerAddress,
				ServerPort:    operation.ServerPort,
				SpanKind:      operation.SpanKind,
			}
		}
		return result
	}
	return []GRPCOperation{}
}

// GetAllGRPCServices returns a list of all services that perform gRPC operations.
func GetAllGRPCServices() []string {
	names, err := systemconfig.GetMetadataStore().GetAllServiceNames(string(systemconfig.GetCurrentSystem()))
	if err == nil && len(names) > 0 {
		var services []string
		for _, service := range names {
			if len(GetOperationsByService(service)) > 0 {
				services = append(services, service)
			}
		}
		if len(services) > 0 {
			return services
		}
	}
	return []string{}
}

// GetClientOperations returns all client-side gRPC operations.
func GetClientOperations() []GRPCOperation {
	return filterOperations(func(operation GRPCOperation) bool {
		return strings.EqualFold(operation.SpanKind, "client")
	})
}

// GetServerOperations returns all server-side gRPC operations.
func GetServerOperations() []GRPCOperation {
	return filterOperations(func(operation GRPCOperation) bool {
		return strings.EqualFold(operation.SpanKind, "server")
	})
}

// GetOperationsByRPCService returns all operations for a specific RPC service.
func GetOperationsByRPCService(rpcService string) []GRPCOperation {
	return filterOperations(func(operation GRPCOperation) bool {
		return operation.RPCService == rpcService
	})
}

func filterOperations(match func(GRPCOperation) bool) []GRPCOperation {
	var results []GRPCOperation
	for _, service := range GetAllGRPCServices() {
		for _, operation := range GetOperationsByService(service) {
			if match(operation) {
				results = append(results, operation)
			}
		}
	}
	return results
}

func convertOtelDemoOperationMap(otelOps map[string][]oteldemogrpc.GRPCOperation) map[string][]GRPCOperation {
	result := make(map[string][]GRPCOperation, len(otelOps))
	for service, operations := range otelOps {
		result[service] = convertOtelDemoOperations(operations)
	}
	return result
}

func convertMediaOperationMap(mediaOps map[string][]mediagrpc.GRPCOperation) map[string][]GRPCOperation {
	result := make(map[string][]GRPCOperation, len(mediaOps))
	for service, operations := range mediaOps {
		result[service] = convertMediaOperations(operations)
	}
	return result
}

func convertHSOperationMap(hsOps map[string][]hsgrpc.GRPCOperation) map[string][]GRPCOperation {
	result := make(map[string][]GRPCOperation, len(hsOps))
	for service, operations := range hsOps {
		result[service] = convertHSOperations(operations)
	}
	return result
}

func convertSNOperationMap(snOps map[string][]sngrpc.GRPCOperation) map[string][]GRPCOperation {
	result := make(map[string][]GRPCOperation, len(snOps))
	for service, operations := range snOps {
		result[service] = convertSNOperations(operations)
	}
	return result
}

func convertOBOperationMap(obOps map[string][]obgrpc.GRPCOperation) map[string][]GRPCOperation {
	result := make(map[string][]GRPCOperation, len(obOps))
	for service, operations := range obOps {
		result[service] = convertOBOperations(operations)
	}
	return result
}

func convertSockShopOperationMap(sockshopOps map[string][]sockshopgrpc.GRPCOperation) map[string][]GRPCOperation {
	result := make(map[string][]GRPCOperation, len(sockshopOps))
	for service, operations := range sockshopOps {
		result[service] = convertSockShopOperations(operations)
	}
	return result
}

func convertTeaStoreOperationMap(teastoreOps map[string][]teastoregrpc.GRPCOperation) map[string][]GRPCOperation {
	result := make(map[string][]GRPCOperation, len(teastoreOps))
	for service, operations := range teastoreOps {
		result[service] = convertTeaStoreOperations(operations)
	}
	return result
}

// convertOtelDemoOperations converts otel-demo-specific operations to the common type.
func convertOtelDemoOperations(otelOps []oteldemogrpc.GRPCOperation) []GRPCOperation {
	result := make([]GRPCOperation, len(otelOps))
	for i, op := range otelOps {
		result[i] = GRPCOperation{
			ServiceName:   op.ServiceName,
			RPCSystem:     op.RPCSystem,
			RPCService:    op.RPCService,
			RPCMethod:     op.RPCMethod,
			StatusCode:    op.StatusCode,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
			SpanKind:      op.SpanKind,
		}
	}
	return result
}

// convertMediaOperations converts media-specific operations to the common type.
func convertMediaOperations(mediaOps []mediagrpc.GRPCOperation) []GRPCOperation {
	result := make([]GRPCOperation, len(mediaOps))
	for i, op := range mediaOps {
		result[i] = GRPCOperation{
			ServiceName:   op.ServiceName,
			RPCSystem:     op.RPCSystem,
			RPCService:    op.RPCService,
			RPCMethod:     op.RPCMethod,
			StatusCode:    op.StatusCode,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
			SpanKind:      op.SpanKind,
		}
	}
	return result
}

// convertHSOperations converts hs-specific operations to the common type.
func convertHSOperations(hsOps []hsgrpc.GRPCOperation) []GRPCOperation {
	result := make([]GRPCOperation, len(hsOps))
	for i, op := range hsOps {
		result[i] = GRPCOperation{
			ServiceName:   op.ServiceName,
			RPCSystem:     op.RPCSystem,
			RPCService:    op.RPCService,
			RPCMethod:     op.RPCMethod,
			StatusCode:    op.StatusCode,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
			SpanKind:      op.SpanKind,
		}
	}
	return result
}

// convertSNOperations converts sn-specific operations to the common type.
func convertSNOperations(snOps []sngrpc.GRPCOperation) []GRPCOperation {
	result := make([]GRPCOperation, len(snOps))
	for i, op := range snOps {
		result[i] = GRPCOperation{
			ServiceName:   op.ServiceName,
			RPCSystem:     op.RPCSystem,
			RPCService:    op.RPCService,
			RPCMethod:     op.RPCMethod,
			StatusCode:    op.StatusCode,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
			SpanKind:      op.SpanKind,
		}
	}
	return result
}

// convertOBOperations converts ob-specific operations to the common type.
func convertOBOperations(obOps []obgrpc.GRPCOperation) []GRPCOperation {
	result := make([]GRPCOperation, len(obOps))
	for i, op := range obOps {
		result[i] = GRPCOperation{
			ServiceName:   op.ServiceName,
			RPCSystem:     op.RPCSystem,
			RPCService:    op.RPCService,
			RPCMethod:     op.RPCMethod,
			StatusCode:    op.StatusCode,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
			SpanKind:      op.SpanKind,
		}
	}
	return result
}

// convertSockShopOperations converts sockshop-specific operations to the common type.
func convertSockShopOperations(sockshopOps []sockshopgrpc.GRPCOperation) []GRPCOperation {
	result := make([]GRPCOperation, len(sockshopOps))
	for i, op := range sockshopOps {
		result[i] = GRPCOperation{
			ServiceName:   op.ServiceName,
			RPCSystem:     op.RPCSystem,
			RPCService:    op.RPCService,
			RPCMethod:     op.RPCMethod,
			StatusCode:    op.StatusCode,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
			SpanKind:      op.SpanKind,
		}
	}
	return result
}

// convertTeaStoreOperations converts teastore-specific operations to the common type.
func convertTeaStoreOperations(teastoreOps []teastoregrpc.GRPCOperation) []GRPCOperation {
	result := make([]GRPCOperation, len(teastoreOps))
	for i, op := range teastoreOps {
		result[i] = GRPCOperation{
			ServiceName:   op.ServiceName,
			RPCSystem:     op.RPCSystem,
			RPCService:    op.RPCService,
			RPCMethod:     op.RPCMethod,
			StatusCode:    op.StatusCode,
			ServerAddress: op.ServerAddress,
			ServerPort:    op.ServerPort,
			SpanKind:      op.SpanKind,
		}
	}
	return result
}

// IsGRPCRoutePattern checks if a route looks like a gRPC route pattern.
func IsGRPCRoutePattern(route string) bool {
	if route == "" || len(route) < 3 {
		return false
	}
	if route[0] != '/' {
		return false
	}

	hasDot := false
	for i := 1; i < len(route); i++ {
		if route[i] == '/' {
			break
		}
		if route[i] == '.' {
			hasDot = true
		}
	}
	return hasDot
}
