package serviceendpoints

import (
	"testing"

	"github.com/LGU-SE-Internal/chaos-experiment/internal/systemconfig"
)

type mockServiceEndpointProvider struct {
	services  []string
	endpoints map[string][]systemconfig.ServiceEndpointData
}

func (m *mockServiceEndpointProvider) GetServiceNames() []string {
	return m.services
}

func (m *mockServiceEndpointProvider) GetEndpointsByService(serviceName string) []systemconfig.ServiceEndpointData {
	return m.endpoints[serviceName]
}

func TestDynamicServiceEndpointProvider(t *testing.T) {
	const testSystem = systemconfig.SystemType("service-endpoints-dynamic")

	t.Cleanup(func() {
		_ = systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)
		_ = systemconfig.UnregisterSystem(testSystem)
	})

	if err := systemconfig.RegisterSystem(systemconfig.SystemRegistration{
		Name:        testSystem,
		NsPattern:   "^service-endpoints-dynamic\\d+$",
		DisplayName: "ServiceEndpointsDynamic",
	}); err != nil {
		t.Fatalf("RegisterSystem() error = %v", err)
	}

	systemconfig.GetRegistry().RegisterServiceEndpointProvider(testSystem, &mockServiceEndpointProvider{
		services: []string{"dynamic-api"},
		endpoints: map[string][]systemconfig.ServiceEndpointData{
			"dynamic-api": {
				{
					ServiceName:    "dynamic-api",
					RequestMethod:  "GET",
					ResponseStatus: "200",
					Route:          "/api/dynamic",
					ServerAddress:  "dynamic-backend",
					ServerPort:     "8080",
					SpanName:       "dynamic.api",
				},
			},
		},
	})

	if err := systemconfig.SetCurrentSystem(testSystem); err != nil {
		t.Fatalf("SetCurrentSystem() error = %v", err)
	}

	services := GetAllServices()
	if len(services) != 1 || services[0] != "dynamic-api" {
		t.Fatalf("GetAllServices() = %#v, want [dynamic-api]", services)
	}

	endpoints := GetEndpointsByService("dynamic-api")
	if len(endpoints) != 1 {
		t.Fatalf("GetEndpointsByService() returned %d endpoints, want 1", len(endpoints))
	}
	if endpoints[0].SpanName != "dynamic.api" {
		t.Fatalf("GetEndpointsByService()[0].SpanName = %q, want %q", endpoints[0].SpanName, "dynamic.api")
	}
}
