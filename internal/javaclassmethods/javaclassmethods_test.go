package javaclassmethods

import (
	"testing"

	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

type mockJavaProvider struct {
	services []string
	methods  map[string][]systemconfig.JavaClassMethodData
}

func (m *mockJavaProvider) GetServiceNames() []string {
	return m.services
}

func (m *mockJavaProvider) GetClassMethodsByService(serviceName string) []systemconfig.JavaClassMethodData {
	return m.methods[serviceName]
}

type mockServiceProvider struct {
	services  []string
	endpoints map[string][]systemconfig.ServiceEndpointData
}

func (m *mockServiceProvider) GetServiceNames() []string {
	return m.services
}

func (m *mockServiceProvider) GetEndpointsByService(serviceName string) []systemconfig.ServiceEndpointData {
	return m.endpoints[serviceName]
}

func TestDynamicJavaProvider(t *testing.T) {
	const testSystem = systemconfig.SystemType("java-methods-dynamic")

	t.Cleanup(func() {
		_ = systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)
		_ = systemconfig.UnregisterSystem(testSystem)
	})

	if err := systemconfig.RegisterSystem(systemconfig.SystemRegistration{
		Name:        testSystem,
		NsPattern:   "^java-methods-dynamic\\d+$",
		DisplayName: "JavaMethodsDynamic",
	}); err != nil {
		t.Fatalf("RegisterSystem() error = %v", err)
	}

	systemconfig.GetRegistry().RegisterServiceEndpointProvider(testSystem, &mockServiceProvider{
		services: []string{"dynamic-api"},
		endpoints: map[string][]systemconfig.ServiceEndpointData{
			"dynamic-api": {
				{ServiceName: "dynamic-api", Route: "/api/dynamic"},
			},
		},
	})

	systemconfig.GetRegistry().RegisterJavaClassMethodProvider(testSystem, &mockJavaProvider{
		services: []string{"dynamic-api", "worker-only"},
		methods: map[string][]systemconfig.JavaClassMethodData{
			"dynamic-api": {
				{ClassName: "dynamic.ApiService", MethodName: "List"},
			},
			"worker-only": {
				{ClassName: "dynamic.Worker", MethodName: "Run"},
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

	methods := GetClassMethodsByService("dynamic-api")
	if len(methods) != 1 || methods[0].MethodName != "List" {
		t.Fatalf("GetClassMethodsByService() = %#v", methods)
	}

	if methods := GetClassMethodsByService("worker-only"); len(methods) != 0 {
		t.Fatalf("worker-only should be filtered by network services, got %#v", methods)
	}
}
