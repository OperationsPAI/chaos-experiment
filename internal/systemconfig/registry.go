package systemconfig

import (
	"fmt"
	"sync"
)

// MetadataType represents the type of metadata.
type MetadataType string

const (
	// MetadataServiceEndpoints represents service endpoint metadata.
	MetadataServiceEndpoints MetadataType = "service_endpoints"
	// MetadataDatabaseOperations represents database operation metadata.
	MetadataDatabaseOperations MetadataType = "database_operations"
	// MetadataJavaClassMethods represents Java class method metadata.
	MetadataJavaClassMethods MetadataType = "java_class_methods"
	// MetadataNetworkDependencies represents network dependency metadata.
	MetadataNetworkDependencies MetadataType = "network_dependencies"
	// MetadataGRPCOperations represents gRPC operation metadata.
	MetadataGRPCOperations MetadataType = "grpc_operations"
)

// MetadataProvider is a shared capability for metadata providers.
type MetadataProvider interface {
	// GetServiceNames returns a list of all service names covered by the provider.
	GetServiceNames() []string
}

// ServiceEndpointProvider provides service endpoint data.
type ServiceEndpointProvider interface {
	MetadataProvider
	// GetEndpointsByService returns endpoints for a specific service.
	GetEndpointsByService(serviceName string) []ServiceEndpointData
}

// ServiceEndpointData represents a service endpoint.
type ServiceEndpointData struct {
	ServiceName    string
	RequestMethod  string
	ResponseStatus string
	Route          string
	ServerAddress  string
	ServerPort     string
	SpanName       string
}

// DatabaseOperationProvider provides database operation data.
type DatabaseOperationProvider interface {
	MetadataProvider
	// GetOperationsByService returns database operations for a specific service.
	GetOperationsByService(serviceName string) []DatabaseOperationData
}

// DatabaseOperationData represents a database operation.
type DatabaseOperationData struct {
	ServiceName   string
	DBName        string
	DBTable       string
	Operation     string
	DBSystem      string
	ServerAddress string
	ServerPort    string
}

// GRPCOperationProvider provides gRPC operation data.
type GRPCOperationProvider interface {
	MetadataProvider
	// GetOperationsByService returns gRPC operations for a specific service.
	GetOperationsByService(serviceName string) []GRPCOperationData
}

// GRPCOperationData represents a gRPC operation.
type GRPCOperationData struct {
	ServiceName    string
	RPCSystem      string
	RPCService     string
	RPCMethod      string
	GRPCStatusCode string
	ServerAddress  string
	ServerPort     string
	SpanKind       string
}

// JavaClassMethodProvider provides Java class method data.
type JavaClassMethodProvider interface {
	MetadataProvider
	// GetClassMethodsByService returns Java class methods for a specific service.
	GetClassMethodsByService(serviceName string) []JavaClassMethodData
}

// JavaClassMethodData represents a Java class method.
type JavaClassMethodData struct {
	ClassName  string
	MethodName string
}

// MetadataRegistry holds registered metadata providers for each system.
type MetadataRegistry struct {
	mu                 sync.RWMutex
	serviceEndpoints   map[SystemType]ServiceEndpointProvider
	databaseOperations map[SystemType]DatabaseOperationProvider
	grpcOperations     map[SystemType]GRPCOperationProvider
	javaClassMethods   map[SystemType]JavaClassMethodProvider
}

var (
	globalRegistry *MetadataRegistry
	registryOnce   sync.Once
)

// GetRegistry returns the global metadata registry.
func GetRegistry() *MetadataRegistry {
	registryOnce.Do(func() {
		globalRegistry = &MetadataRegistry{
			serviceEndpoints:   make(map[SystemType]ServiceEndpointProvider),
			databaseOperations: make(map[SystemType]DatabaseOperationProvider),
			grpcOperations:     make(map[SystemType]GRPCOperationProvider),
			javaClassMethods:   make(map[SystemType]JavaClassMethodProvider),
		}
	})
	return globalRegistry
}

// RegisterServiceEndpointProvider registers a service endpoint provider for a system.
func (r *MetadataRegistry) RegisterServiceEndpointProvider(system SystemType, provider ServiceEndpointProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.serviceEndpoints[system] = provider
	if router := getMetadataStoreRouter(); router != nil {
		router.RegisterServiceEndpointProvider(system, provider)
	}
}

// RegisterDatabaseOperationProvider registers a database operation provider for a system.
func (r *MetadataRegistry) RegisterDatabaseOperationProvider(system SystemType, provider DatabaseOperationProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.databaseOperations[system] = provider
	if router := getMetadataStoreRouter(); router != nil {
		router.RegisterDatabaseOperationProvider(system, provider)
	}
}

// RegisterGRPCOperationProvider registers a gRPC operation provider for a system.
func (r *MetadataRegistry) RegisterGRPCOperationProvider(system SystemType, provider GRPCOperationProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.grpcOperations[system] = provider
	if router := getMetadataStoreRouter(); router != nil {
		router.RegisterGRPCOperationProvider(system, provider)
	}
}

// RegisterJavaClassMethodProvider registers a Java class method provider for a system.
func (r *MetadataRegistry) RegisterJavaClassMethodProvider(system SystemType, provider JavaClassMethodProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.javaClassMethods[system] = provider
	if router := getMetadataStoreRouter(); router != nil {
		router.RegisterJavaClassMethodProvider(system, provider)
	}
}

// GetServiceEndpointProvider returns the service endpoint provider for the current system.
func (r *MetadataRegistry) GetServiceEndpointProvider() (ServiceEndpointProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.serviceEndpoints[GetCurrentSystem()]
	if !exists {
		return nil, fmt.Errorf("no service endpoint provider registered for system: %s", GetCurrentSystem())
	}
	return provider, nil
}

// GetDatabaseOperationProvider returns the database operation provider for the current system.
func (r *MetadataRegistry) GetDatabaseOperationProvider() (DatabaseOperationProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.databaseOperations[GetCurrentSystem()]
	if !exists {
		return nil, fmt.Errorf("no database operation provider registered for system: %s", GetCurrentSystem())
	}
	return provider, nil
}

// GetGRPCOperationProvider returns the gRPC operation provider for the current system.
func (r *MetadataRegistry) GetGRPCOperationProvider() (GRPCOperationProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.grpcOperations[GetCurrentSystem()]
	if !exists {
		return nil, fmt.Errorf("no gRPC operation provider registered for system: %s", GetCurrentSystem())
	}
	return provider, nil
}

// GetJavaClassMethodProvider returns the Java class method provider for the current system.
func (r *MetadataRegistry) GetJavaClassMethodProvider() (JavaClassMethodProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.javaClassMethods[GetCurrentSystem()]
	if !exists {
		return nil, fmt.Errorf("no Java class method provider registered for system: %s", GetCurrentSystem())
	}
	return provider, nil
}

// HasServiceEndpointProvider checks if a service endpoint provider is registered for the current system.
func (r *MetadataRegistry) HasServiceEndpointProvider() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.serviceEndpoints[GetCurrentSystem()]
	return exists
}

// HasDatabaseOperationProvider checks if a database operation provider is registered for the current system.
func (r *MetadataRegistry) HasDatabaseOperationProvider() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.databaseOperations[GetCurrentSystem()]
	return exists
}

// HasGRPCOperationProvider checks if a gRPC operation provider is registered for the current system.
func (r *MetadataRegistry) HasGRPCOperationProvider() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.grpcOperations[GetCurrentSystem()]
	return exists
}

// HasJavaClassMethodProvider checks if a Java class method provider is registered for the current system.
func (r *MetadataRegistry) HasJavaClassMethodProvider() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.javaClassMethods[GetCurrentSystem()]
	return exists
}

// Clear removes all registered providers. Useful for testing.
func (r *MetadataRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.serviceEndpoints = make(map[SystemType]ServiceEndpointProvider)
	r.databaseOperations = make(map[SystemType]DatabaseOperationProvider)
	r.grpcOperations = make(map[SystemType]GRPCOperationProvider)
	r.javaClassMethods = make(map[SystemType]JavaClassMethodProvider)
	if router := getMetadataStoreRouter(); router != nil {
		router.ClearInternal()
	}
}

// UnregisterSystem removes all providers associated with a system.
func (r *MetadataRegistry) UnregisterSystem(system SystemType) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.serviceEndpoints, system)
	delete(r.databaseOperations, system)
	delete(r.grpcOperations, system)
	delete(r.javaClassMethods, system)
	if router := getMetadataStoreRouter(); router != nil {
		router.UnregisterSystem(system)
	}
}
