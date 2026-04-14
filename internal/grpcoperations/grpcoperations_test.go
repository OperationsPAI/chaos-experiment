package grpcoperations

import (
	"testing"

	"github.com/LGU-SE-Internal/chaos-experiment/internal/systemconfig"
)

type mockGRPCProvider struct {
	services   []string
	operations map[string][]systemconfig.GRPCOperationData
}

func (m *mockGRPCProvider) GetServiceNames() []string {
	return m.services
}

func (m *mockGRPCProvider) GetOperationsByService(serviceName string) []systemconfig.GRPCOperationData {
	return m.operations[serviceName]
}

func TestDynamicGRPCProvider(t *testing.T) {
	const testSystem = systemconfig.SystemType("grpc-operations-dynamic")

	t.Cleanup(func() {
		_ = systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)
		_ = systemconfig.UnregisterSystem(testSystem)
	})

	if err := systemconfig.RegisterSystem(systemconfig.SystemRegistration{
		Name:        testSystem,
		NsPattern:   "^grpc-operations-dynamic\\d+$",
		DisplayName: "GRPCOperationsDynamic",
	}); err != nil {
		t.Fatalf("RegisterSystem() error = %v", err)
	}

	systemconfig.GetRegistry().RegisterGRPCOperationProvider(testSystem, &mockGRPCProvider{
		services: []string{"dynamic-api"},
		operations: map[string][]systemconfig.GRPCOperationData{
			"dynamic-api": {
				{
					ServiceName:    "dynamic-api",
					RPCSystem:      "grpc",
					RPCService:     "dynamic.API",
					RPCMethod:      "List",
					GRPCStatusCode: "OK",
					SpanKind:       "client",
				},
				{
					ServiceName:    "dynamic-api",
					RPCSystem:      "grpc",
					RPCService:     "dynamic.API",
					RPCMethod:      "Create",
					GRPCStatusCode: "OK",
					SpanKind:       "server",
				},
			},
		},
	})

	if err := systemconfig.SetCurrentSystem(testSystem); err != nil {
		t.Fatalf("SetCurrentSystem() error = %v", err)
	}

	if services := GetAllGRPCServices(); len(services) != 1 || services[0] != "dynamic-api" {
		t.Fatalf("GetAllGRPCServices() = %#v, want [dynamic-api]", services)
	}

	if ops := GetOperationsByService("dynamic-api"); len(ops) != 2 {
		t.Fatalf("GetOperationsByService() returned %d ops, want 2", len(ops))
	}

	if ops := GetClientOperations(); len(ops) != 1 || ops[0].RPCMethod != "List" {
		t.Fatalf("GetClientOperations() = %#v", ops)
	}

	if ops := GetServerOperations(); len(ops) != 1 || ops[0].RPCMethod != "Create" {
		t.Fatalf("GetServerOperations() = %#v", ops)
	}

	if ops := GetOperationsByRPCService("dynamic.API"); len(ops) != 2 {
		t.Fatalf("GetOperationsByRPCService() returned %d ops, want 2", len(ops))
	}
}
