package databaseoperations

import (
	"testing"

	"github.com/LGU-SE-Internal/chaos-experiment/internal/systemconfig"
)

type mockDatabaseOperationProvider struct {
	services   []string
	operations map[string][]systemconfig.DatabaseOperationData
}

func (m *mockDatabaseOperationProvider) GetServiceNames() []string {
	return m.services
}

func (m *mockDatabaseOperationProvider) GetOperationsByService(serviceName string) []systemconfig.DatabaseOperationData {
	return m.operations[serviceName]
}

func TestDynamicDatabaseOperationProvider(t *testing.T) {
	const testSystem = systemconfig.SystemType("database-operations-dynamic")

	t.Cleanup(func() {
		_ = systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)
		_ = systemconfig.UnregisterSystem(testSystem)
	})

	if err := systemconfig.RegisterSystem(systemconfig.SystemRegistration{
		Name:        testSystem,
		NsPattern:   "^database-operations-dynamic\\d+$",
		DisplayName: "DatabaseOperationsDynamic",
	}); err != nil {
		t.Fatalf("RegisterSystem() error = %v", err)
	}

	systemconfig.GetRegistry().RegisterDatabaseOperationProvider(testSystem, &mockDatabaseOperationProvider{
		services: []string{"dynamic-api"},
		operations: map[string][]systemconfig.DatabaseOperationData{
			"dynamic-api": {
				{
					ServiceName: "dynamic-api",
					DBName:      "dynamic",
					DBTable:     "records",
					Operation:   "SELECT",
				},
			},
		},
	})

	if err := systemconfig.SetCurrentSystem(testSystem); err != nil {
		t.Fatalf("SetCurrentSystem() error = %v", err)
	}

	services := GetAllDatabaseServices()
	if len(services) != 1 || services[0] != "dynamic-api" {
		t.Fatalf("GetAllDatabaseServices() = %#v, want [dynamic-api]", services)
	}

	ops := GetOperationsByService("dynamic-api")
	if len(ops) != 1 || ops[0].DBTable != "records" {
		t.Fatalf("GetOperationsByService() = %#v", ops)
	}

	if byDB := GetOperationsByDatabase("dynamic"); len(byDB) != 1 {
		t.Fatalf("GetOperationsByDatabase() returned %d results, want 1", len(byDB))
	}

	if byTable := GetOperationsByTable("records"); len(byTable) != 1 {
		t.Fatalf("GetOperationsByTable() returned %d results, want 1", len(byTable))
	}
}
