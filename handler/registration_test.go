package handler

import (
	"testing"

	"github.com/OperationsPAI/chaos-experiment/internal/databaseoperations"
	"github.com/OperationsPAI/chaos-experiment/internal/grpcoperations"
	"github.com/OperationsPAI/chaos-experiment/internal/registry"
	"github.com/OperationsPAI/chaos-experiment/internal/serviceendpoints"
	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

func TestRegisterSystemData(t *testing.T) {
	const systemName = "dynamic-system-data"

	t.Cleanup(func() {
		_ = systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)
		if IsSystemRegistered(systemName) {
			_ = UnregisterSystem(systemName)
		}
	})

	if err := RegisterSystem(SystemConfig{
		Name:        systemName,
		NsPattern:   "^dynamic-system-data\\d+$",
		DisplayName: "DynamicSystemData",
	}); err != nil {
		t.Fatalf("RegisterSystem() error = %v", err)
	}

	err := RegisterSystemData(systemName, &SystemData{
		SystemName: "dynamic-system-data",
		HTTPEndpoints: map[string][]SystemHTTPEndpoint{
			"checkout": {
				{
					ServiceName:    "checkout",
					RequestMethod:  "POST",
					ResponseStatus: "200",
					Route:          "/api/checkout",
					ServerAddress:  "payment",
					ServerPort:     "8080",
					SpanName:       "checkout.request",
				},
			},
		},
		DatabaseOperations: map[string][]SystemDatabaseOperation{
			"checkout": {
				{
					ServiceName: "checkout",
					DBName:      "orders",
					DBTable:     "checkout",
					Operation:   "SELECT",
					DBSystem:    "mysql",
				},
			},
		},
		RPCOperations: map[string][]SystemRPCOperation{
			"checkout": {
				{
					ServiceName:   "checkout",
					RPCSystem:     "grpc",
					RPCService:    "payment.Service",
					RPCMethod:     "Charge",
					StatusCode:    "OK",
					ServerAddress: "payment",
					ServerPort:    "9090",
					SpanKind:      "client",
				},
			},
		},
		AllServices: []string{"checkout", "payment"},
	})
	if err != nil {
		t.Fatalf("RegisterSystemData() error = %v", err)
	}

	if !IsSystemDataRegistered(systemName) {
		t.Fatal("IsSystemDataRegistered() = false, want true")
	}

	if err := systemconfig.SetCurrentSystem(systemconfig.SystemType(systemName)); err != nil {
		t.Fatalf("SetCurrentSystem() error = %v", err)
	}

	services := serviceendpoints.GetAllServices()
	if len(services) != 2 {
		t.Fatalf("GetAllServices() returned %d services, want 2", len(services))
	}

	httpEndpoints := serviceendpoints.GetEndpointsByService("checkout")
	if len(httpEndpoints) != 1 || httpEndpoints[0].Route != "/api/checkout" {
		t.Fatalf("GetEndpointsByService() = %#v", httpEndpoints)
	}

	dbOps := databaseoperations.GetOperationsByService("checkout")
	if len(dbOps) != 1 || dbOps[0].DBTable != "checkout" {
		t.Fatalf("GetOperationsByService() = %#v", dbOps)
	}

	rpcOps := grpcoperations.GetOperationsByService("checkout")
	if len(rpcOps) != 1 || rpcOps[0].RPCMethod != "Charge" {
		t.Fatalf("GetOperationsByService() = %#v", rpcOps)
	}

	if !registry.IsRegistered(systemconfig.SystemType(systemName)) {
		t.Fatal("registry.IsRegistered() = false, want true")
	}
}

func TestUnregisterSystemRemovesRuntimeData(t *testing.T) {
	const systemName = "dynamic-system-cleanup"
	system := systemconfig.SystemType(systemName)

	t.Cleanup(func() {
		_ = systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)
		if IsSystemRegistered(systemName) {
			_ = UnregisterSystem(systemName)
		}
	})

	if err := RegisterSystem(SystemConfig{
		Name:        systemName,
		NsPattern:   "^dynamic-system-cleanup\\d+$",
		DisplayName: "DynamicSystemCleanup",
	}); err != nil {
		t.Fatalf("RegisterSystem() error = %v", err)
	}

	if err := RegisterServiceEndpointProvider(systemName, &testServiceEndpointProvider{
		services: []string{"api"},
		endpoints: map[string][]ServiceEndpointData{
			"api": {{ServiceName: "api", Route: "/health"}},
		},
	}); err != nil {
		t.Fatalf("RegisterServiceEndpointProvider() error = %v", err)
	}

	if err := UnregisterSystem(systemName); err != nil {
		t.Fatalf("UnregisterSystem() error = %v", err)
	}

	if IsSystemRegistered(systemName) {
		t.Fatal("IsSystemRegistered() = true after UnregisterSystem()")
	}
	if registry.IsRegistered(system) {
		t.Fatal("registry.IsRegistered() = true after UnregisterSystem()")
	}
	if systemconfig.GetRegistration(system) != nil {
		t.Fatal("GetRegistration() should return nil after unregister")
	}
}

type testServiceEndpointProvider struct {
	services  []string
	endpoints map[string][]ServiceEndpointData
}

func (p *testServiceEndpointProvider) GetServiceNames() []string {
	return p.services
}

func (p *testServiceEndpointProvider) GetEndpointsByService(serviceName string) []ServiceEndpointData {
	return p.endpoints[serviceName]
}
