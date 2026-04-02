package handler

import (
	"context"
)

// ChaosResourceMapping maps a chaos type to its available resource targets.
type ChaosResourceMapping struct {
	Services   []string `json:"services,omitempty"`
	Containers []string `json:"containers,omitempty"`
	Endpoints  []string `json:"endpoints,omitempty"`
}

// SystemResource holds resource information for a specific system.
type SystemResource struct {
	Services   []string `json:"services,omitempty"`
	Pods       []string `json:"pods,omitempty"`
	Containers []string `json:"containers,omitempty"`
	Endpoints  []string `json:"endpoints,omitempty"`
}

// GetChaosTypeResourceMappings returns the resource mappings for all chaos types.
// This is a stub implementation that returns an empty map.
func GetChaosTypeResourceMappings() (map[string]ChaosResourceMapping, error) {
	return make(map[string]ChaosResourceMapping), nil
}

// GetSystemResourceMap returns system resources for all registered systems.
// This is a stub implementation that returns an empty map.
func GetSystemResourceMap(ctx context.Context) (map[SystemType]SystemResource, error) {
	return make(map[SystemType]SystemResource), nil
}
