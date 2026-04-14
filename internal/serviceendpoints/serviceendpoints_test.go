package serviceendpoints

import (
	"testing"

	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

type mockServiceEndpointProvider struct {
	services  []string
	endpoints map[string][]systemconfig.ServiceEndpointData
}

type mockMetadataStore struct {
	serviceNames map[string][]string
	endpoints    map[string]map[string][]systemconfig.ServiceEndpointData
}

func (m *mockServiceEndpointProvider) GetServiceNames() []string {
	return m.services
}

func (m *mockServiceEndpointProvider) GetEndpointsByService(serviceName string) []systemconfig.ServiceEndpointData {
	return m.endpoints[serviceName]
}

func (m *mockMetadataStore) GetServiceEndpoints(system string, serviceName string) ([]systemconfig.ServiceEndpointData, error) {
	return m.endpoints[system][serviceName], nil
}

func (m *mockMetadataStore) GetAllServiceNames(system string) ([]string, error) {
	return m.serviceNames[system], nil
}

func (m *mockMetadataStore) GetJavaClassMethods(system string, serviceName string) ([]systemconfig.JavaClassMethodData, error) {
	return nil, nil
}

func (m *mockMetadataStore) GetDatabaseOperations(system string, serviceName string) ([]systemconfig.DatabaseOperationData, error) {
	return nil, nil
}

func (m *mockMetadataStore) GetGRPCOperations(system string, serviceName string) ([]systemconfig.GRPCOperationData, error) {
	return nil, nil
}

func (m *mockMetadataStore) GetNetworkPairs(system string) ([]systemconfig.NetworkPairData, error) {
	return nil, nil
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

func TestMetadataStoreOverridesProvider(t *testing.T) {
	const testSystem = systemconfig.SystemType("service-endpoints-metadata-store")

	t.Cleanup(func() {
		systemconfig.SetMetadataStore(nil)
		_ = systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)
		_ = systemconfig.UnregisterSystem(testSystem)
	})

	if err := systemconfig.RegisterSystem(systemconfig.SystemRegistration{
		Name:        testSystem,
		NsPattern:   "^service-endpoints-metadata-store\\d+$",
		DisplayName: "ServiceEndpointsMetadataStore",
	}); err != nil {
		t.Fatalf("RegisterSystem() error = %v", err)
	}

	systemconfig.GetRegistry().RegisterServiceEndpointProvider(testSystem, &mockServiceEndpointProvider{
		services: []string{"provider-api"},
		endpoints: map[string][]systemconfig.ServiceEndpointData{
			"provider-api": {{ServiceName: "provider-api", Route: "/provider"}},
		},
	})

	systemconfig.SetMetadataStore(&mockMetadataStore{
		serviceNames: map[string][]string{
			string(testSystem): {"store-api"},
		},
		endpoints: map[string]map[string][]systemconfig.ServiceEndpointData{
			string(testSystem): {
				"store-api": {{ServiceName: "store-api", Route: "/store"}},
			},
		},
	})

	if err := systemconfig.SetCurrentSystem(testSystem); err != nil {
		t.Fatalf("SetCurrentSystem() error = %v", err)
	}

	services := GetAllServices()
	if len(services) != 1 || services[0] != "store-api" {
		t.Fatalf("GetAllServices() = %#v, want [store-api]", services)
	}

	endpoints := GetEndpointsByService("store-api")
	if len(endpoints) != 1 || endpoints[0].Route != "/store" {
		t.Fatalf("GetEndpointsByService() = %#v, want metadata-store endpoint", endpoints)
	}
}
