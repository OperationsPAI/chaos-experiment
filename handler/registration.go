package handler

import (
	"fmt"
	"sort"

	"github.com/LGU-SE-Internal/chaos-experiment/internal/model"
	"github.com/LGU-SE-Internal/chaos-experiment/internal/registry"
	"github.com/LGU-SE-Internal/chaos-experiment/internal/resourcelookup"
	"github.com/LGU-SE-Internal/chaos-experiment/internal/resourcetypes"
	"github.com/LGU-SE-Internal/chaos-experiment/internal/systemconfig"
)

// Re-export the system-data types so external callers do not need internal imports.
type SystemData = model.SystemData
type SystemHTTPEndpoint = resourcetypes.HTTPEndpoint
type SystemDatabaseOperation = resourcetypes.DatabaseOperation
type SystemRPCOperation = resourcetypes.RPCOperation

// Re-export provider interfaces and payload types for runtime registration.
type ServiceEndpointProvider = systemconfig.ServiceEndpointProvider
type ServiceEndpointData = systemconfig.ServiceEndpointData
type DatabaseOperationProvider = systemconfig.DatabaseOperationProvider
type DatabaseOperationData = systemconfig.DatabaseOperationData
type GRPCOperationProvider = systemconfig.GRPCOperationProvider
type GRPCOperationData = systemconfig.GRPCOperationData
type JavaClassMethodProvider = systemconfig.JavaClassMethodProvider
type JavaClassMethodData = systemconfig.JavaClassMethodData

// SystemConfig is the public-facing configuration for registering a system.
type SystemConfig struct {
	Name        string
	NsPattern   string
	DisplayName string
}

// RegisterSystem registers a new system type with the given configuration.
func RegisterSystem(cfg SystemConfig) error {
	return systemconfig.RegisterSystem(systemconfig.SystemRegistration{
		Name:        systemconfig.SystemType(cfg.Name),
		NsPattern:   cfg.NsPattern,
		DisplayName: cfg.DisplayName,
	})
}

// RegisterSystemData registers a system-data bundle and wires it into both the
// injectionv2 registry layer and the dynamic metadata provider layer.
func RegisterSystemData(name string, data *SystemData) error {
	system := systemconfig.SystemType(name)
	if err := validateRegisteredSystem(system); err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("system data must not be nil")
	}

	normalized := normalizeSystemData(data)
	if err := registry.Register(system, normalized); err != nil {
		return err
	}

	metaRegistry := systemconfig.GetRegistry()
	metaRegistry.RegisterServiceEndpointProvider(system, &systemDataServiceEndpointProvider{data: normalized})
	metaRegistry.RegisterDatabaseOperationProvider(system, &systemDataDatabaseOperationProvider{data: normalized})
	metaRegistry.RegisterGRPCOperationProvider(system, &systemDataGRPCOperationProvider{data: normalized})
	resourcelookup.ResetSystemCache(system)
	return nil
}

// RegisterServiceEndpointProvider registers a runtime service-endpoint provider for a system.
func RegisterServiceEndpointProvider(name string, provider ServiceEndpointProvider) error {
	system := systemconfig.SystemType(name)
	if err := validateRegisteredSystem(system); err != nil {
		return err
	}
	if provider == nil {
		return fmt.Errorf("service endpoint provider must not be nil")
	}

	systemconfig.GetRegistry().RegisterServiceEndpointProvider(system, provider)
	resourcelookup.ResetSystemCache(system)
	return nil
}

// RegisterDatabaseOperationProvider registers a runtime database provider for a system.
func RegisterDatabaseOperationProvider(name string, provider DatabaseOperationProvider) error {
	system := systemconfig.SystemType(name)
	if err := validateRegisteredSystem(system); err != nil {
		return err
	}
	if provider == nil {
		return fmt.Errorf("database operation provider must not be nil")
	}

	systemconfig.GetRegistry().RegisterDatabaseOperationProvider(system, provider)
	resourcelookup.ResetSystemCache(system)
	return nil
}

// RegisterGRPCOperationProvider registers a runtime gRPC provider for a system.
func RegisterGRPCOperationProvider(name string, provider GRPCOperationProvider) error {
	system := systemconfig.SystemType(name)
	if err := validateRegisteredSystem(system); err != nil {
		return err
	}
	if provider == nil {
		return fmt.Errorf("gRPC operation provider must not be nil")
	}

	systemconfig.GetRegistry().RegisterGRPCOperationProvider(system, provider)
	resourcelookup.ResetSystemCache(system)
	return nil
}

// RegisterJavaClassMethodProvider registers a runtime JVM method provider for a system.
func RegisterJavaClassMethodProvider(name string, provider JavaClassMethodProvider) error {
	system := systemconfig.SystemType(name)
	if err := validateRegisteredSystem(system); err != nil {
		return err
	}
	if provider == nil {
		return fmt.Errorf("Java class method provider must not be nil")
	}

	systemconfig.GetRegistry().RegisterJavaClassMethodProvider(system, provider)
	resourcelookup.ResetSystemCache(system)
	return nil
}

// UnregisterSystem removes a previously registered system type and any attached runtime data/providers.
func UnregisterSystem(name string) error {
	system := systemconfig.SystemType(name)
	if err := validateRegisteredSystem(system); err != nil {
		return err
	}

	systemconfig.GetRegistry().UnregisterSystem(system)
	registry.Unregister(system)
	resourcelookup.ResetSystemCache(system)
	return systemconfig.UnregisterSystem(system)
}

// IsSystemRegistered returns true if the named system type is registered.
func IsSystemRegistered(name string) bool {
	return systemconfig.IsRegistered(systemconfig.SystemType(name))
}

// IsSystemDataRegistered returns true if system data has been registered for the named system.
func IsSystemDataRegistered(name string) bool {
	return registry.IsRegistered(systemconfig.SystemType(name))
}

func validateRegisteredSystem(system systemconfig.SystemType) error {
	if !systemconfig.IsRegistered(system) {
		return fmt.Errorf("system %s is not registered; call RegisterSystem first", system)
	}
	return nil
}

func normalizeSystemData(data *SystemData) *SystemData {
	normalized := &model.SystemData{
		SystemName:         data.SystemName,
		HTTPEndpoints:      cloneHTTPEndpoints(data.HTTPEndpoints),
		DatabaseOperations: cloneDatabaseOperations(data.DatabaseOperations),
		RPCOperations:      cloneRPCOperations(data.RPCOperations),
		AllServices:        append([]string(nil), data.AllServices...),
	}

	if len(normalized.AllServices) == 0 {
		normalized.AllServices = collectAllServices(normalized)
	}

	return normalized
}

func collectAllServices(data *SystemData) []string {
	serviceSet := make(map[string]struct{})
	for service := range data.HTTPEndpoints {
		serviceSet[service] = struct{}{}
	}
	for service := range data.DatabaseOperations {
		serviceSet[service] = struct{}{}
	}
	for service := range data.RPCOperations {
		serviceSet[service] = struct{}{}
	}

	services := make([]string, 0, len(serviceSet))
	for service := range serviceSet {
		services = append(services, service)
	}
	sort.Strings(services)
	return services
}

func cloneHTTPEndpoints(src map[string][]resourcetypes.HTTPEndpoint) map[string][]resourcetypes.HTTPEndpoint {
	if len(src) == 0 {
		return make(map[string][]resourcetypes.HTTPEndpoint)
	}

	dst := make(map[string][]resourcetypes.HTTPEndpoint, len(src))
	for service, endpoints := range src {
		dst[service] = append([]resourcetypes.HTTPEndpoint(nil), endpoints...)
	}
	return dst
}

func cloneDatabaseOperations(src map[string][]resourcetypes.DatabaseOperation) map[string][]resourcetypes.DatabaseOperation {
	if len(src) == 0 {
		return make(map[string][]resourcetypes.DatabaseOperation)
	}

	dst := make(map[string][]resourcetypes.DatabaseOperation, len(src))
	for service, operations := range src {
		dst[service] = append([]resourcetypes.DatabaseOperation(nil), operations...)
	}
	return dst
}

func cloneRPCOperations(src map[string][]resourcetypes.RPCOperation) map[string][]resourcetypes.RPCOperation {
	if len(src) == 0 {
		return make(map[string][]resourcetypes.RPCOperation)
	}

	dst := make(map[string][]resourcetypes.RPCOperation, len(src))
	for service, operations := range src {
		dst[service] = append([]resourcetypes.RPCOperation(nil), operations...)
	}
	return dst
}

type systemDataServiceEndpointProvider struct {
	data *model.SystemData
}

func (p *systemDataServiceEndpointProvider) GetServiceNames() []string {
	return append([]string(nil), p.data.GetAllServices()...)
}

func (p *systemDataServiceEndpointProvider) GetEndpointsByService(serviceName string) []systemconfig.ServiceEndpointData {
	endpoints := p.data.GetHTTPEndpointsByService(serviceName)
	result := make([]systemconfig.ServiceEndpointData, len(endpoints))
	for i, endpoint := range endpoints {
		result[i] = systemconfig.ServiceEndpointData{
			ServiceName:    endpoint.ServiceName,
			RequestMethod:  endpoint.RequestMethod,
			ResponseStatus: endpoint.ResponseStatus,
			Route:          endpoint.Route,
			ServerAddress:  endpoint.ServerAddress,
			ServerPort:     endpoint.ServerPort,
			SpanName:       endpoint.SpanName,
		}
	}
	return result
}

type systemDataDatabaseOperationProvider struct {
	data *model.SystemData
}

func (p *systemDataDatabaseOperationProvider) GetServiceNames() []string {
	services := make([]string, 0, len(p.data.DatabaseOperations))
	for service := range p.data.DatabaseOperations {
		services = append(services, service)
	}
	sort.Strings(services)
	return services
}

func (p *systemDataDatabaseOperationProvider) GetOperationsByService(serviceName string) []systemconfig.DatabaseOperationData {
	operations := p.data.GetDatabaseOperationsByService(serviceName)
	result := make([]systemconfig.DatabaseOperationData, len(operations))
	for i, operation := range operations {
		result[i] = systemconfig.DatabaseOperationData{
			ServiceName:   operation.ServiceName,
			DBName:        operation.DBName,
			DBTable:       operation.DBTable,
			Operation:     operation.Operation,
			DBSystem:      operation.DBSystem,
			ServerAddress: operation.ServerAddress,
			ServerPort:    operation.ServerPort,
		}
	}
	return result
}

type systemDataGRPCOperationProvider struct {
	data *model.SystemData
}

func (p *systemDataGRPCOperationProvider) GetServiceNames() []string {
	services := make([]string, 0, len(p.data.RPCOperations))
	for service := range p.data.RPCOperations {
		services = append(services, service)
	}
	sort.Strings(services)
	return services
}

func (p *systemDataGRPCOperationProvider) GetOperationsByService(serviceName string) []systemconfig.GRPCOperationData {
	operations := p.data.GetRPCOperationsByService(serviceName)
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
