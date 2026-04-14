package systemconfig

import (
	"sort"
	"sync"
)

// NetworkPairData represents a network communication pair for the MetadataStore interface.
type NetworkPairData struct {
	Source string
	Target string
}

// MetadataStore is the primary metadata access contract for dynamic system registration.
// Runtime consumers should read through this interface, while provider/system-data
// registrations are adapted into the store internally.
type MetadataStore interface {
	GetServiceEndpoints(system string, serviceName string) ([]ServiceEndpointData, error)
	GetAllServiceNames(system string) ([]string, error)
	GetJavaClassMethods(system string, serviceName string) ([]JavaClassMethodData, error)
	GetDatabaseOperations(system string, serviceName string) ([]DatabaseOperationData, error)
	GetGRPCOperations(system string, serviceName string) ([]GRPCOperationData, error)
	GetNetworkPairs(system string) ([]NetworkPairData, error)
}

type metadataStoreRouter struct {
	mu sync.RWMutex

	external           MetadataStore
	serviceEndpoints   map[SystemType]ServiceEndpointProvider
	databaseOperations map[SystemType]DatabaseOperationProvider
	grpcOperations     map[SystemType]GRPCOperationProvider
	javaClassMethods   map[SystemType]JavaClassMethodProvider
}

var (
	globalMetadataStore     MetadataStore
	globalMetadataStoreOnce sync.Once
)

func getMetadataStoreRouter() *metadataStoreRouter {
	store, ok := GetMetadataStore().(*metadataStoreRouter)
	if ok {
		return store
	}
	return nil
}

func newMetadataStoreRouter() *metadataStoreRouter {
	return &metadataStoreRouter{
		serviceEndpoints:   make(map[SystemType]ServiceEndpointProvider),
		databaseOperations: make(map[SystemType]DatabaseOperationProvider),
		grpcOperations:     make(map[SystemType]GRPCOperationProvider),
		javaClassMethods:   make(map[SystemType]JavaClassMethodProvider),
	}
}

// SetMetadataStore sets the primary external MetadataStore implementation.
// Registered provider/system-data adapters continue to act as fallback sources.
func SetMetadataStore(store MetadataStore) {
	router := getMetadataStoreRouter()
	if router == nil {
		return
	}

	router.mu.Lock()
	defer router.mu.Unlock()
	router.external = store
}

// GetMetadataStore returns the process-wide MetadataStore router.
func GetMetadataStore() MetadataStore {
	globalMetadataStoreOnce.Do(func() {
		globalMetadataStore = newMetadataStoreRouter()
	})
	return globalMetadataStore
}

func (m *metadataStoreRouter) RegisterServiceEndpointProvider(system SystemType, provider ServiceEndpointProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.serviceEndpoints[system] = provider
}

func (m *metadataStoreRouter) RegisterDatabaseOperationProvider(system SystemType, provider DatabaseOperationProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.databaseOperations[system] = provider
}

func (m *metadataStoreRouter) RegisterGRPCOperationProvider(system SystemType, provider GRPCOperationProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.grpcOperations[system] = provider
}

func (m *metadataStoreRouter) RegisterJavaClassMethodProvider(system SystemType, provider JavaClassMethodProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.javaClassMethods[system] = provider
}

func (m *metadataStoreRouter) ClearInternal() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.serviceEndpoints = make(map[SystemType]ServiceEndpointProvider)
	m.databaseOperations = make(map[SystemType]DatabaseOperationProvider)
	m.grpcOperations = make(map[SystemType]GRPCOperationProvider)
	m.javaClassMethods = make(map[SystemType]JavaClassMethodProvider)
}

func (m *metadataStoreRouter) UnregisterSystem(system SystemType) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.serviceEndpoints, system)
	delete(m.databaseOperations, system)
	delete(m.grpcOperations, system)
	delete(m.javaClassMethods, system)
}

func (m *metadataStoreRouter) GetServiceEndpoints(system string, serviceName string) ([]ServiceEndpointData, error) {
	if external := m.getExternal(); external != nil {
		data, err := external.GetServiceEndpoints(system, serviceName)
		if err == nil && len(data) > 0 {
			return data, nil
		}
	}

	provider := m.getServiceEndpointProvider(SystemType(system))
	if provider == nil {
		return nil, nil
	}
	return provider.GetEndpointsByService(serviceName), nil
}

func (m *metadataStoreRouter) GetAllServiceNames(system string) ([]string, error) {
	if external := m.getExternal(); external != nil {
		names, err := external.GetAllServiceNames(system)
		if err == nil && len(names) > 0 {
			return append([]string(nil), names...), nil
		}
	}

	return m.listInternalServices(SystemType(system)), nil
}

func (m *metadataStoreRouter) GetJavaClassMethods(system string, serviceName string) ([]JavaClassMethodData, error) {
	if external := m.getExternal(); external != nil {
		data, err := external.GetJavaClassMethods(system, serviceName)
		if err == nil && len(data) > 0 {
			return data, nil
		}
	}

	provider := m.getJavaClassMethodProvider(SystemType(system))
	if provider == nil {
		return nil, nil
	}
	return provider.GetClassMethodsByService(serviceName), nil
}

func (m *metadataStoreRouter) GetDatabaseOperations(system string, serviceName string) ([]DatabaseOperationData, error) {
	if external := m.getExternal(); external != nil {
		data, err := external.GetDatabaseOperations(system, serviceName)
		if err == nil && len(data) > 0 {
			return data, nil
		}
	}

	provider := m.getDatabaseOperationProvider(SystemType(system))
	if provider == nil {
		return nil, nil
	}
	return provider.GetOperationsByService(serviceName), nil
}

func (m *metadataStoreRouter) GetGRPCOperations(system string, serviceName string) ([]GRPCOperationData, error) {
	if external := m.getExternal(); external != nil {
		data, err := external.GetGRPCOperations(system, serviceName)
		if err == nil && len(data) > 0 {
			return data, nil
		}
	}

	provider := m.getGRPCOperationProvider(SystemType(system))
	if provider == nil {
		return nil, nil
	}
	return provider.GetOperationsByService(serviceName), nil
}

func (m *metadataStoreRouter) GetNetworkPairs(system string) ([]NetworkPairData, error) {
	if external := m.getExternal(); external != nil {
		data, err := external.GetNetworkPairs(system)
		if err == nil && len(data) > 0 {
			return data, nil
		}
	}

	systemType := SystemType(system)
	pairs := make(map[string]NetworkPairData)

	serviceProvider := m.getServiceEndpointProvider(systemType)
	if serviceProvider != nil {
		for _, service := range serviceProvider.GetServiceNames() {
			for _, endpoint := range serviceProvider.GetEndpointsByService(service) {
				if endpoint.ServerAddress == "" || endpoint.ServerAddress == service {
					continue
				}
				key := service + "->" + endpoint.ServerAddress
				pairs[key] = NetworkPairData{Source: service, Target: endpoint.ServerAddress}
			}
		}
	}

	grpcProvider := m.getGRPCOperationProvider(systemType)
	if grpcProvider != nil {
		for _, service := range grpcProvider.GetServiceNames() {
			for _, operation := range grpcProvider.GetOperationsByService(service) {
				if operation.ServerAddress == "" || operation.ServerAddress == service {
					continue
				}
				key := service + "->" + operation.ServerAddress
				pairs[key] = NetworkPairData{Source: service, Target: operation.ServerAddress}
			}
		}
	}

	result := make([]NetworkPairData, 0, len(pairs))
	for _, pair := range pairs {
		result = append(result, pair)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Source != result[j].Source {
			return result[i].Source < result[j].Source
		}
		return result[i].Target < result[j].Target
	})
	return result, nil
}

func (m *metadataStoreRouter) getExternal() MetadataStore {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.external
}

func (m *metadataStoreRouter) getServiceEndpointProvider(system SystemType) ServiceEndpointProvider {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.serviceEndpoints[system]
}

func (m *metadataStoreRouter) getDatabaseOperationProvider(system SystemType) DatabaseOperationProvider {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.databaseOperations[system]
}

func (m *metadataStoreRouter) getGRPCOperationProvider(system SystemType) GRPCOperationProvider {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.grpcOperations[system]
}

func (m *metadataStoreRouter) getJavaClassMethodProvider(system SystemType) JavaClassMethodProvider {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.javaClassMethods[system]
}

func (m *metadataStoreRouter) listInternalServices(system SystemType) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if provider := m.serviceEndpoints[system]; provider != nil {
		names := append([]string(nil), provider.GetServiceNames()...)
		sort.Strings(names)
		return names
	}

	serviceSet := make(map[string]struct{})
	addNames := func(provider MetadataProvider) {
		if provider == nil {
			return
		}
		for _, service := range provider.GetServiceNames() {
			serviceSet[service] = struct{}{}
		}
	}

	addNames(m.serviceEndpoints[system])
	addNames(m.databaseOperations[system])
	addNames(m.grpcOperations[system])
	addNames(m.javaClassMethods[system])

	services := make([]string, 0, len(serviceSet))
	for service := range serviceSet {
		services = append(services, service)
	}
	sort.Strings(services)
	return services
}
