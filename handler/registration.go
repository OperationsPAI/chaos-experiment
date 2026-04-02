package handler

import (
	"github.com/LGU-SE-Internal/chaos-experiment/internal/systemconfig"
)

// Re-export types so external callers can implement/use them without importing internal packages.
type MetadataStore = systemconfig.MetadataStore
type ServiceEndpointData = systemconfig.ServiceEndpointData
type JavaClassMethodData = systemconfig.JavaClassMethodData
type DatabaseOperationData = systemconfig.DatabaseOperationData
type GRPCOperationData = systemconfig.GRPCOperationData
type NetworkPairData = systemconfig.NetworkPairData

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

// UnregisterSystem removes a previously registered system type.
func UnregisterSystem(name string) error {
	return systemconfig.UnregisterSystem(systemconfig.SystemType(name))
}

// IsSystemRegistered returns true if the named system type is registered.
func IsSystemRegistered(name string) bool {
	return systemconfig.IsRegistered(systemconfig.SystemType(name))
}

// SetMetadataStore sets the global MetadataStore implementation.
// External callers should implement the MetadataStore interface and call this
// function to provide dynamic metadata for dynamically registered systems.
func SetMetadataStore(store MetadataStore) {
	systemconfig.SetMetadataStore(store)
}
