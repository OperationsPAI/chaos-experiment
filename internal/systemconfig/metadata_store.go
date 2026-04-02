package systemconfig

import (
	"sync"
)

// ServiceEndpointData represents a service endpoint for the MetadataStore interface.
type ServiceEndpointData struct {
	Method string
	Path   string
}

// JavaClassMethodData represents a Java class-method pair for the MetadataStore interface.
type JavaClassMethodData struct {
	ClassName  string
	MethodName string
}

// DatabaseOperationData represents a database operation for the MetadataStore interface.
type DatabaseOperationData struct {
	Operation string
	Table     string
}

// GRPCOperationData represents a gRPC operation for the MetadataStore interface.
type GRPCOperationData struct {
	Service string
	Method  string
}

// NetworkPairData represents a network communication pair for the MetadataStore interface.
type NetworkPairData struct {
	Source string
	Target string
}

// MetadataStore is an interface that external callers (e.g. AegisLab) can implement
// to provide dynamic metadata for dynamically registered systems.
type MetadataStore interface {
	GetServiceEndpoints(system string, serviceName string) ([]ServiceEndpointData, error)
	GetAllServiceNames(system string) ([]string, error)
	GetJavaClassMethods(system string, serviceName string) ([]JavaClassMethodData, error)
	GetDatabaseOperations(system string, serviceName string) ([]DatabaseOperationData, error)
	GetGRPCOperations(system string, serviceName string) ([]GRPCOperationData, error)
	GetNetworkPairs(system string) ([]NetworkPairData, error)
}

var (
	globalMetadataStore MetadataStore
	metadataStoreMu     sync.RWMutex
)

// SetMetadataStore sets the global MetadataStore implementation.
func SetMetadataStore(store MetadataStore) {
	metadataStoreMu.Lock()
	defer metadataStoreMu.Unlock()
	globalMetadataStore = store
}

// GetMetadataStore returns the global MetadataStore, or nil if not set.
func GetMetadataStore() MetadataStore {
	metadataStoreMu.RLock()
	defer metadataStoreMu.RUnlock()
	return globalMetadataStore
}
