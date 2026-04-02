// Package systemconfig provides a global configuration for the target system type.
// This package allows different systems (TrainTicket, OtelDemo, etc.) to coexist
// with their own metadata and configurations.
package systemconfig

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// SystemType represents the type of system being analyzed/targeted
type SystemType string

const (
	// SystemTrainTicket represents the TrainTicket microservice system
	SystemTrainTicket SystemType = "ts"
	// SystemOtelDemo represents the OpenTelemetry Demo system
	SystemOtelDemo SystemType = "otel-demo"
	// SystemMediaMicroservices represents the Media Microservices system
	SystemMediaMicroservices SystemType = "media"
	// SystemHotelReservation represents the Hotel Reservation system
	SystemHotelReservation SystemType = "hs"
	// SystemSocialNetwork represents the Social Network system
	SystemSocialNetwork SystemType = "sn"
	// SystemOnlineBoutique represents the Online Boutique system
	SystemOnlineBoutique SystemType = "ob"
)

// SystemRegistration holds registration data for a system.
type SystemRegistration struct {
	Name        SystemType
	NsPattern   string
	DisplayName string
}

var (
	// currentSystem holds the current system type
	currentSystem SystemType = SystemTrainTicket

	// mu protects access to currentSystem
	mu sync.RWMutex

	// registeredSystems holds all registered system types
	registeredSystems   = make(map[SystemType]*SystemRegistration)
	registeredSystemsMu sync.RWMutex
)

func init() {
	// Pre-register the 6 built-in systems
	builtins := []SystemRegistration{
		{Name: SystemTrainTicket, NsPattern: "^ts\\d+$", DisplayName: "TrainTicket"},
		{Name: SystemOtelDemo, NsPattern: "^otel-demo\\d+$", DisplayName: "OtelDemo"},
		{Name: SystemMediaMicroservices, NsPattern: "^media\\d+$", DisplayName: "MediaMicroservices"},
		{Name: SystemHotelReservation, NsPattern: "^hs\\d+$", DisplayName: "HotelReservation"},
		{Name: SystemSocialNetwork, NsPattern: "^sn\\d+$", DisplayName: "SocialNetwork"},
		{Name: SystemOnlineBoutique, NsPattern: "^ob\\d+$", DisplayName: "OnlineBoutique"},
	}
	for i := range builtins {
		registeredSystems[builtins[i].Name] = &builtins[i]
	}
}

// RegisterSystem registers a new system type. Returns an error if the name is
// empty or already registered.
func RegisterSystem(reg SystemRegistration) error {
	if reg.Name == "" {
		return fmt.Errorf("system name must not be empty")
	}
	registeredSystemsMu.Lock()
	defer registeredSystemsMu.Unlock()
	if _, exists := registeredSystems[reg.Name]; exists {
		return fmt.Errorf("system %s is already registered", reg.Name)
	}
	copied := reg
	registeredSystems[reg.Name] = &copied
	return nil
}

// UnregisterSystem removes a previously registered system. Returns an error if
// the system is not registered.
func UnregisterSystem(name SystemType) error {
	registeredSystemsMu.Lock()
	defer registeredSystemsMu.Unlock()
	if _, exists := registeredSystems[name]; !exists {
		return fmt.Errorf("system %s is not registered", name)
	}
	delete(registeredSystems, name)
	return nil
}

// IsRegistered returns true if the given system type is registered.
func IsRegistered(system SystemType) bool {
	registeredSystemsMu.RLock()
	defer registeredSystemsMu.RUnlock()
	_, exists := registeredSystems[system]
	return exists
}

// GetRegistration returns the registration for a system, or nil if not found.
func GetRegistration(system SystemType) *SystemRegistration {
	registeredSystemsMu.RLock()
	defer registeredSystemsMu.RUnlock()
	reg, exists := registeredSystems[system]
	if !exists {
		return nil
	}
	// Return a copy to avoid data races
	copied := *reg
	return &copied
}

// GetAllRegisteredSystems returns all registered system type names, sorted.
func GetAllRegisteredSystems() []SystemType {
	registeredSystemsMu.RLock()
	defer registeredSystemsMu.RUnlock()
	result := make([]SystemType, 0, len(registeredSystems))
	for k := range registeredSystems {
		result = append(result, k)
	}
	sort.Slice(result, func(i, j int) bool {
		return string(result[i]) < string(result[j])
	})
	return result
}

// SetCurrentSystem sets the global system type for the current process.
// This should be called at initialization time before any metadata is accessed.
func SetCurrentSystem(system SystemType) error {
	if !IsRegistered(system) {
		return fmt.Errorf("invalid system type: %s, system is not registered", system)
	}

	mu.Lock()
	defer mu.Unlock()
	currentSystem = system
	return nil
}

// GetCurrentSystem returns the current system type.
func GetCurrentSystem() SystemType {
	mu.RLock()
	defer mu.RUnlock()
	return currentSystem
}

// IsTrainTicket returns true if the current system is TrainTicket.
func IsTrainTicket() bool {
	return GetCurrentSystem() == SystemTrainTicket
}

// IsOtelDemo returns true if the current system is OpenTelemetry Demo.
func IsOtelDemo() bool {
	return GetCurrentSystem() == SystemOtelDemo
}

// IsMediaMicroservices returns true if the current system is Media Microservices.
func IsMediaMicroservices() bool {
	return GetCurrentSystem() == SystemMediaMicroservices
}

// IsHotelReservation returns true if the current system is Hotel Reservation.
func IsHotelReservation() bool {
	return GetCurrentSystem() == SystemHotelReservation
}

// IsSocialNetwork returns true if the current system is Social Network.
func IsSocialNetwork() bool {
	return GetCurrentSystem() == SystemSocialNetwork
}

// IsOnlineBoutique returns true if the current system is Online Boutique.
func IsOnlineBoutique() bool {
	return GetCurrentSystem() == SystemOnlineBoutique
}

// String returns the string representation of the SystemType.
func (s SystemType) String() string {
	return string(s)
}

// IsValid returns true if this SystemType is registered.
func (s SystemType) IsValid() bool {
	return IsRegistered(s)
}

// GetAllSystemTypes returns all registered system types.
func GetAllSystemTypes() []SystemType {
	return GetAllRegisteredSystems()
}

// GetNamespaceByIndex generates a namespace name based on the system type and index.
func GetNamespaceByIndex(system SystemType, index int) (string, error) {
	reg := GetRegistration(system)
	if reg == nil {
		return "", fmt.Errorf("system type not found: %s", system)
	}

	name := strings.TrimPrefix(reg.NsPattern, "^")
	name = strings.TrimSuffix(name, "$")
	name = strings.Replace(name, "\\d+", fmt.Sprintf("%d", index), 1)

	return name, nil
}

// ParseSystemType parses a string into a SystemType.
func ParseSystemType(s string) (SystemType, error) {
	st := SystemType(s)
	if !IsRegistered(st) {
		return "", fmt.Errorf("invalid system type: %s, system is not registered", s)
	}
	return st, nil
}
