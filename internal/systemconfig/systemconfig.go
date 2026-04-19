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

// SystemType represents the type of system being analyzed/targeted.
type SystemType string

const (
	// SystemTrainTicket represents the TrainTicket microservice system.
	SystemTrainTicket SystemType = "ts"
	// SystemOtelDemo represents the OpenTelemetry Demo system.
	SystemOtelDemo SystemType = "otel-demo"
	// SystemMediaMicroservices represents the Media Microservices system.
	SystemMediaMicroservices SystemType = "media"
	// SystemHotelReservation represents the Hotel Reservation system.
	SystemHotelReservation SystemType = "hs"
	// SystemSocialNetwork represents the Social Network system.
	SystemSocialNetwork SystemType = "sn"
	// SystemOnlineBoutique represents the Online Boutique system.
	SystemOnlineBoutique SystemType = "ob"
	// SystemSockShop represents the Sock Shop system.
	SystemSockShop SystemType = "sockshop"
	// SystemTeaStore represents the Tea Store system.
	SystemTeaStore SystemType = "teastore"
)

// SystemRegistration holds registration data for a system.
type SystemRegistration struct {
	Name        SystemType
	NsPattern   string
	DisplayName string
	// AppLabelKey is the pod label key used to select application workloads
	// (e.g. "app", "app.kubernetes.io/name"). An empty value defaults to "app".
	AppLabelKey string
}

var (
	currentSystem   = SystemTrainTicket
	currentSystemMu sync.RWMutex

	registeredSystems   map[SystemType]SystemRegistration
	registeredSystemsMu sync.RWMutex
)

func init() {
	registeredSystems = make(map[SystemType]SystemRegistration)
	for _, reg := range builtinSystemRegistrations() {
		registeredSystems[reg.Name] = reg
	}
}

func builtinSystemRegistrations() []SystemRegistration {
	return []SystemRegistration{
		{Name: SystemTrainTicket, NsPattern: "^ts\\d+$", DisplayName: "TrainTicket", AppLabelKey: "app"},
		{Name: SystemOtelDemo, NsPattern: "^otel-demo\\d+$", DisplayName: "OtelDemo", AppLabelKey: "app.kubernetes.io/name"},
		{Name: SystemMediaMicroservices, NsPattern: "^media\\d+$", DisplayName: "MediaMicroservices", AppLabelKey: "app"},
		{Name: SystemHotelReservation, NsPattern: "^hs\\d+$", DisplayName: "HotelReservation", AppLabelKey: "app"},
		{Name: SystemSocialNetwork, NsPattern: "^sn\\d+$", DisplayName: "SocialNetwork", AppLabelKey: "app"},
		{Name: SystemOnlineBoutique, NsPattern: "^ob\\d+$", DisplayName: "OnlineBoutique", AppLabelKey: "app"},
		{Name: SystemSockShop, NsPattern: "^sockshop\\d+$", DisplayName: "SockShop", AppLabelKey: "app"},
		{Name: SystemTeaStore, NsPattern: "^teastore\\d+$", DisplayName: "TeaStore", AppLabelKey: "app"},
	}
}

// GetAppLabelKey returns the pod selector label key for the given system.
// Defaults to "app" if the system is unregistered or the registration has an empty AppLabelKey.
func GetAppLabelKey(system SystemType) string {
	reg := GetRegistration(system)
	if reg == nil || reg.AppLabelKey == "" {
		return "app"
	}
	return reg.AppLabelKey
}

// GetCurrentAppLabelKey returns the pod selector label key for the current system.
func GetCurrentAppLabelKey() string {
	return GetAppLabelKey(GetCurrentSystem())
}

// RegisterSystem registers a new system type.
func RegisterSystem(reg SystemRegistration) error {
	if reg.Name == "" {
		return fmt.Errorf("system name must not be empty")
	}

	registeredSystemsMu.Lock()
	defer registeredSystemsMu.Unlock()

	if _, exists := registeredSystems[reg.Name]; exists {
		return fmt.Errorf("system %s is already registered", reg.Name)
	}

	registeredSystems[reg.Name] = reg
	return nil
}

// UnregisterSystem removes a previously registered system.
func UnregisterSystem(name SystemType) error {
	registeredSystemsMu.Lock()
	defer registeredSystemsMu.Unlock()

	if _, exists := registeredSystems[name]; !exists {
		return fmt.Errorf("system %s is not registered", name)
	}

	delete(registeredSystems, name)

	currentSystemMu.Lock()
	defer currentSystemMu.Unlock()
	if currentSystem != name {
		return nil
	}

	if _, exists := registeredSystems[SystemTrainTicket]; exists {
		currentSystem = SystemTrainTicket
		return nil
	}

	currentSystem = ""
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

	copied := reg
	return &copied
}

// GetAllRegisteredSystems returns all registered system type names in a stable order.
func GetAllRegisteredSystems() []SystemType {
	registeredSystemsMu.RLock()
	defer registeredSystemsMu.RUnlock()

	systems := make([]SystemType, 0, len(registeredSystems))
	for system := range registeredSystems {
		systems = append(systems, system)
	}

	sort.Slice(systems, func(i, j int) bool {
		return string(systems[i]) < string(systems[j])
	})

	return systems
}

// SetCurrentSystem sets the global system type for the current process.
// This should be called at initialization time before any metadata is accessed.
func SetCurrentSystem(system SystemType) error {
	if !IsRegistered(system) {
		return fmt.Errorf("invalid system type: %s, valid types are: %s", system, strings.Join(registeredSystemNames(), ", "))
	}

	currentSystemMu.Lock()
	defer currentSystemMu.Unlock()

	currentSystem = system
	return nil
}

// GetCurrentSystem returns the current system type.
func GetCurrentSystem() SystemType {
	currentSystemMu.RLock()
	defer currentSystemMu.RUnlock()
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

// IsSockShop returns true if the current system is Sock Shop.
func IsSockShop() bool {
	return GetCurrentSystem() == SystemSockShop
}

// IsTeaStore returns true if the current system is Tea Store.
func IsTeaStore() bool {
	return GetCurrentSystem() == SystemTeaStore
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
	system := SystemType(s)
	if !IsRegistered(system) {
		return "", fmt.Errorf("invalid system type: %s, valid types are: %s", s, strings.Join(registeredSystemNames(), ", "))
	}
	return system, nil
}

func registeredSystemNames() []string {
	systems := GetAllRegisteredSystems()
	names := make([]string, len(systems))
	for i, system := range systems {
		names[i] = system.String()
	}
	return names
}
