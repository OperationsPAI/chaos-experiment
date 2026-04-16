// Package metadataaccess provides a unified interface for accessing system-specific metadata.
// It uses the systemconfig package to determine the current system and routes calls to
// the appropriate registered metadata accessors.
package metadataaccess

import (
	"sync"

	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

// ServiceEndpoint represents a service endpoint from trace analysis.
type ServiceEndpoint struct {
	ServiceName    string
	RequestMethod  string
	ResponseStatus string
	Route          string
	ServerAddress  string
	ServerPort     string
}

// DatabaseOperation represents a database operation from trace analysis.
type DatabaseOperation struct {
	ServiceName   string
	DBName        string
	DBTable       string
	Operation     string
	DBSystem      string
	ServerAddress string
	ServerPort    string
}

// GRPCOperation represents a gRPC operation from trace analysis.
type GRPCOperation struct {
	ServiceName    string
	RPCSystem      string
	RPCService     string
	RPCMethod      string
	GRPCStatusCode string
	ServerAddress  string
	ServerPort     string
	SpanKind       string
}

// JavaClassMethod represents a Java class method for JVM chaos.
type JavaClassMethod struct {
	ClassName  string
	MethodName string
}

// ServiceEndpointAccessor defines the interface for service endpoint access.
type ServiceEndpointAccessor interface {
	GetEndpointsByService(serviceName string) []ServiceEndpoint
	GetAllServices() []string
}

// DatabaseOperationAccessor defines the interface for database operation access.
type DatabaseOperationAccessor interface {
	GetOperationsByService(serviceName string) []DatabaseOperation
	GetAllDatabaseServices() []string
}

// GRPCOperationAccessor defines the interface for gRPC operation access.
type GRPCOperationAccessor interface {
	GetOperationsByService(serviceName string) []GRPCOperation
	GetAllGRPCServices() []string
}

// JavaMethodAccessor defines the interface for Java method access.
type JavaMethodAccessor interface {
	GetClassMethodsByService(serviceName string) []JavaClassMethod
	GetAllServices() []string
}

// MetadataAccessor provides unified access to system-specific metadata.
type MetadataAccessor struct {
	mu sync.RWMutex

	serviceEndpoints   map[systemconfig.SystemType]ServiceEndpointAccessor
	databaseOperations map[systemconfig.SystemType]DatabaseOperationAccessor
	grpcOperations     map[systemconfig.SystemType]GRPCOperationAccessor
	javaMethods        map[systemconfig.SystemType]JavaMethodAccessor
}

var (
	accessor     *MetadataAccessor
	accessorOnce sync.Once
)

// GetAccessor returns the singleton metadata accessor.
func GetAccessor() *MetadataAccessor {
	accessorOnce.Do(func() {
		accessor = &MetadataAccessor{
			serviceEndpoints:   make(map[systemconfig.SystemType]ServiceEndpointAccessor),
			databaseOperations: make(map[systemconfig.SystemType]DatabaseOperationAccessor),
			grpcOperations:     make(map[systemconfig.SystemType]GRPCOperationAccessor),
			javaMethods:        make(map[systemconfig.SystemType]JavaMethodAccessor),
		}
	})
	return accessor
}

// RegisterServiceEndpoints registers a service endpoint accessor for a system.
func (m *MetadataAccessor) RegisterServiceEndpoints(system systemconfig.SystemType, accessor ServiceEndpointAccessor) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.serviceEndpoints[system] = accessor
}

// RegisterDatabaseOperations registers a database operation accessor for a system.
func (m *MetadataAccessor) RegisterDatabaseOperations(system systemconfig.SystemType, accessor DatabaseOperationAccessor) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.databaseOperations[system] = accessor
}

// RegisterGRPCOperations registers a gRPC operation accessor for a system.
func (m *MetadataAccessor) RegisterGRPCOperations(system systemconfig.SystemType, accessor GRPCOperationAccessor) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.grpcOperations[system] = accessor
}

// RegisterJavaMethods registers a Java method accessor for a system.
func (m *MetadataAccessor) RegisterJavaMethods(system systemconfig.SystemType, accessor JavaMethodAccessor) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.javaMethods[system] = accessor
}

// RegisterTrainTicketServiceEndpoints registers the TrainTicket service endpoint accessor.
func (m *MetadataAccessor) RegisterTrainTicketServiceEndpoints(accessor ServiceEndpointAccessor) {
	m.RegisterServiceEndpoints(systemconfig.SystemTrainTicket, accessor)
}

// RegisterTrainTicketDatabaseOperations registers the TrainTicket database operation accessor.
func (m *MetadataAccessor) RegisterTrainTicketDatabaseOperations(accessor DatabaseOperationAccessor) {
	m.RegisterDatabaseOperations(systemconfig.SystemTrainTicket, accessor)
}

// RegisterTrainTicketGRPCOperations registers the TrainTicket gRPC operation accessor.
func (m *MetadataAccessor) RegisterTrainTicketGRPCOperations(accessor GRPCOperationAccessor) {
	m.RegisterGRPCOperations(systemconfig.SystemTrainTicket, accessor)
}

// RegisterTrainTicketJavaMethods registers the TrainTicket Java method accessor.
func (m *MetadataAccessor) RegisterTrainTicketJavaMethods(accessor JavaMethodAccessor) {
	m.RegisterJavaMethods(systemconfig.SystemTrainTicket, accessor)
}

// RegisterOtelDemoServiceEndpoints registers the OtelDemo service endpoint accessor.
func (m *MetadataAccessor) RegisterOtelDemoServiceEndpoints(accessor ServiceEndpointAccessor) {
	m.RegisterServiceEndpoints(systemconfig.SystemOtelDemo, accessor)
}

// RegisterOtelDemoDatabaseOperations registers the OtelDemo database operation accessor.
func (m *MetadataAccessor) RegisterOtelDemoDatabaseOperations(accessor DatabaseOperationAccessor) {
	m.RegisterDatabaseOperations(systemconfig.SystemOtelDemo, accessor)
}

// RegisterOtelDemoGRPCOperations registers the OtelDemo gRPC operation accessor.
func (m *MetadataAccessor) RegisterOtelDemoGRPCOperations(accessor GRPCOperationAccessor) {
	m.RegisterGRPCOperations(systemconfig.SystemOtelDemo, accessor)
}

// RegisterOtelDemoJavaMethods registers the OtelDemo Java method accessor.
func (m *MetadataAccessor) RegisterOtelDemoJavaMethods(accessor JavaMethodAccessor) {
	m.RegisterJavaMethods(systemconfig.SystemOtelDemo, accessor)
}

// RegisterMediaMicroservicesServiceEndpoints registers the MediaMicroservices service endpoint accessor.
func (m *MetadataAccessor) RegisterMediaMicroservicesServiceEndpoints(accessor ServiceEndpointAccessor) {
	m.RegisterServiceEndpoints(systemconfig.SystemMediaMicroservices, accessor)
}

// RegisterMediaMicroservicesDatabaseOperations registers the MediaMicroservices database operation accessor.
func (m *MetadataAccessor) RegisterMediaMicroservicesDatabaseOperations(accessor DatabaseOperationAccessor) {
	m.RegisterDatabaseOperations(systemconfig.SystemMediaMicroservices, accessor)
}

// RegisterMediaMicroservicesGRPCOperations registers the MediaMicroservices gRPC operation accessor.
func (m *MetadataAccessor) RegisterMediaMicroservicesGRPCOperations(accessor GRPCOperationAccessor) {
	m.RegisterGRPCOperations(systemconfig.SystemMediaMicroservices, accessor)
}

// RegisterMediaMicroservicesJavaMethods registers the MediaMicroservices Java method accessor.
func (m *MetadataAccessor) RegisterMediaMicroservicesJavaMethods(accessor JavaMethodAccessor) {
	m.RegisterJavaMethods(systemconfig.SystemMediaMicroservices, accessor)
}

// RegisterHotelReservationServiceEndpoints registers the HotelReservation service endpoint accessor.
func (m *MetadataAccessor) RegisterHotelReservationServiceEndpoints(accessor ServiceEndpointAccessor) {
	m.RegisterServiceEndpoints(systemconfig.SystemHotelReservation, accessor)
}

// RegisterHotelReservationDatabaseOperations registers the HotelReservation database operation accessor.
func (m *MetadataAccessor) RegisterHotelReservationDatabaseOperations(accessor DatabaseOperationAccessor) {
	m.RegisterDatabaseOperations(systemconfig.SystemHotelReservation, accessor)
}

// RegisterHotelReservationGRPCOperations registers the HotelReservation gRPC operation accessor.
func (m *MetadataAccessor) RegisterHotelReservationGRPCOperations(accessor GRPCOperationAccessor) {
	m.RegisterGRPCOperations(systemconfig.SystemHotelReservation, accessor)
}

// RegisterHotelReservationJavaMethods registers the HotelReservation Java method accessor.
func (m *MetadataAccessor) RegisterHotelReservationJavaMethods(accessor JavaMethodAccessor) {
	m.RegisterJavaMethods(systemconfig.SystemHotelReservation, accessor)
}

// RegisterSocialNetworkServiceEndpoints registers the SocialNetwork service endpoint accessor.
func (m *MetadataAccessor) RegisterSocialNetworkServiceEndpoints(accessor ServiceEndpointAccessor) {
	m.RegisterServiceEndpoints(systemconfig.SystemSocialNetwork, accessor)
}

// RegisterSocialNetworkDatabaseOperations registers the SocialNetwork database operation accessor.
func (m *MetadataAccessor) RegisterSocialNetworkDatabaseOperations(accessor DatabaseOperationAccessor) {
	m.RegisterDatabaseOperations(systemconfig.SystemSocialNetwork, accessor)
}

// RegisterSocialNetworkGRPCOperations registers the SocialNetwork gRPC operation accessor.
func (m *MetadataAccessor) RegisterSocialNetworkGRPCOperations(accessor GRPCOperationAccessor) {
	m.RegisterGRPCOperations(systemconfig.SystemSocialNetwork, accessor)
}

// RegisterSocialNetworkJavaMethods registers the SocialNetwork Java method accessor.
func (m *MetadataAccessor) RegisterSocialNetworkJavaMethods(accessor JavaMethodAccessor) {
	m.RegisterJavaMethods(systemconfig.SystemSocialNetwork, accessor)
}

// RegisterOnlineBoutiqueServiceEndpoints registers the OnlineBoutique service endpoint accessor.
func (m *MetadataAccessor) RegisterOnlineBoutiqueServiceEndpoints(accessor ServiceEndpointAccessor) {
	m.RegisterServiceEndpoints(systemconfig.SystemOnlineBoutique, accessor)
}

// RegisterOnlineBoutiqueDatabaseOperations registers the OnlineBoutique database operation accessor.
func (m *MetadataAccessor) RegisterOnlineBoutiqueDatabaseOperations(accessor DatabaseOperationAccessor) {
	m.RegisterDatabaseOperations(systemconfig.SystemOnlineBoutique, accessor)
}

// RegisterOnlineBoutiqueGRPCOperations registers the OnlineBoutique gRPC operation accessor.
func (m *MetadataAccessor) RegisterOnlineBoutiqueGRPCOperations(accessor GRPCOperationAccessor) {
	m.RegisterGRPCOperations(systemconfig.SystemOnlineBoutique, accessor)
}

// RegisterOnlineBoutiqueJavaMethods registers the OnlineBoutique Java method accessor.
func (m *MetadataAccessor) RegisterOnlineBoutiqueJavaMethods(accessor JavaMethodAccessor) {
	m.RegisterJavaMethods(systemconfig.SystemOnlineBoutique, accessor)
}

// RegisterSockShopServiceEndpoints registers the SockShop service endpoint accessor.
func (m *MetadataAccessor) RegisterSockShopServiceEndpoints(accessor ServiceEndpointAccessor) {
	m.RegisterServiceEndpoints(systemconfig.SystemSockShop, accessor)
}

// RegisterSockShopDatabaseOperations registers the SockShop database operation accessor.
func (m *MetadataAccessor) RegisterSockShopDatabaseOperations(accessor DatabaseOperationAccessor) {
	m.RegisterDatabaseOperations(systemconfig.SystemSockShop, accessor)
}

// RegisterSockShopGRPCOperations registers the SockShop gRPC operation accessor.
func (m *MetadataAccessor) RegisterSockShopGRPCOperations(accessor GRPCOperationAccessor) {
	m.RegisterGRPCOperations(systemconfig.SystemSockShop, accessor)
}

// RegisterSockShopJavaMethods registers the SockShop Java method accessor.
func (m *MetadataAccessor) RegisterSockShopJavaMethods(accessor JavaMethodAccessor) {
	m.RegisterJavaMethods(systemconfig.SystemSockShop, accessor)
}

// RegisterTeaStoreServiceEndpoints registers the TeaStore service endpoint accessor.
func (m *MetadataAccessor) RegisterTeaStoreServiceEndpoints(accessor ServiceEndpointAccessor) {
	m.RegisterServiceEndpoints(systemconfig.SystemTeaStore, accessor)
}

// RegisterTeaStoreDatabaseOperations registers the TeaStore database operation accessor.
func (m *MetadataAccessor) RegisterTeaStoreDatabaseOperations(accessor DatabaseOperationAccessor) {
	m.RegisterDatabaseOperations(systemconfig.SystemTeaStore, accessor)
}

// RegisterTeaStoreGRPCOperations registers the TeaStore gRPC operation accessor.
func (m *MetadataAccessor) RegisterTeaStoreGRPCOperations(accessor GRPCOperationAccessor) {
	m.RegisterGRPCOperations(systemconfig.SystemTeaStore, accessor)
}

// RegisterTeaStoreJavaMethods registers the TeaStore Java method accessor.
func (m *MetadataAccessor) RegisterTeaStoreJavaMethods(accessor JavaMethodAccessor) {
	m.RegisterJavaMethods(systemconfig.SystemTeaStore, accessor)
}

// GetEndpointsByService returns endpoints for a service based on the current system.
func (m *MetadataAccessor) GetEndpointsByService(serviceName string) []ServiceEndpoint {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if accessor := m.serviceEndpoints[systemconfig.GetCurrentSystem()]; accessor != nil {
		return accessor.GetEndpointsByService(serviceName)
	}
	return nil
}

// GetAllServices returns all services based on the current system.
func (m *MetadataAccessor) GetAllServices() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if accessor := m.serviceEndpoints[systemconfig.GetCurrentSystem()]; accessor != nil {
		return accessor.GetAllServices()
	}
	return nil
}

// GetDatabaseOperationsByService returns database operations for a service based on the current system.
func (m *MetadataAccessor) GetDatabaseOperationsByService(serviceName string) []DatabaseOperation {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if accessor := m.databaseOperations[systemconfig.GetCurrentSystem()]; accessor != nil {
		return accessor.GetOperationsByService(serviceName)
	}
	return nil
}

// GetAllDatabaseServices returns all database services based on the current system.
func (m *MetadataAccessor) GetAllDatabaseServices() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if accessor := m.databaseOperations[systemconfig.GetCurrentSystem()]; accessor != nil {
		return accessor.GetAllDatabaseServices()
	}
	return nil
}

// GetGRPCOperationsByService returns gRPC operations for a service based on the current system.
func (m *MetadataAccessor) GetGRPCOperationsByService(serviceName string) []GRPCOperation {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if accessor := m.grpcOperations[systemconfig.GetCurrentSystem()]; accessor != nil {
		return accessor.GetOperationsByService(serviceName)
	}
	return nil
}

// GetAllGRPCServices returns all gRPC services based on the current system.
func (m *MetadataAccessor) GetAllGRPCServices() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if accessor := m.grpcOperations[systemconfig.GetCurrentSystem()]; accessor != nil {
		return accessor.GetAllGRPCServices()
	}
	return nil
}

// GetJavaMethodsByService returns Java methods for a service based on the current system.
func (m *MetadataAccessor) GetJavaMethodsByService(serviceName string) []JavaClassMethod {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if accessor := m.javaMethods[systemconfig.GetCurrentSystem()]; accessor != nil {
		return accessor.GetClassMethodsByService(serviceName)
	}
	return nil
}

// GetAllJavaServices returns all Java services based on the current system.
func (m *MetadataAccessor) GetAllJavaServices() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if accessor := m.javaMethods[systemconfig.GetCurrentSystem()]; accessor != nil {
		return accessor.GetAllServices()
	}
	return nil
}

// HasServiceEndpoints returns true if service endpoints are available for the current system.
func (m *MetadataAccessor) HasServiceEndpoints() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.serviceEndpoints[systemconfig.GetCurrentSystem()] != nil
}

// HasDatabaseOperations returns true if database operations are available for the current system.
func (m *MetadataAccessor) HasDatabaseOperations() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.databaseOperations[systemconfig.GetCurrentSystem()] != nil
}

// HasGRPCOperations returns true if gRPC operations are available for the current system.
func (m *MetadataAccessor) HasGRPCOperations() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.grpcOperations[systemconfig.GetCurrentSystem()] != nil
}

// HasJavaMethods returns true if Java methods are available for the current system.
func (m *MetadataAccessor) HasJavaMethods() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.javaMethods[systemconfig.GetCurrentSystem()] != nil
}

// Clear removes all registered accessors. Useful for testing.
func (m *MetadataAccessor) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.serviceEndpoints = make(map[systemconfig.SystemType]ServiceEndpointAccessor)
	m.databaseOperations = make(map[systemconfig.SystemType]DatabaseOperationAccessor)
	m.grpcOperations = make(map[systemconfig.SystemType]GRPCOperationAccessor)
	m.javaMethods = make(map[systemconfig.SystemType]JavaMethodAccessor)
}
