package networkdependencies_test

import (
	"testing"

	"github.com/OperationsPAI/chaos-experiment/internal/networkdependencies"
	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
	"github.com/OperationsPAI/chaos-experiment/internal/testdata"
)

type storeBackedMetadataStore struct {
	pairs map[string][]systemconfig.NetworkPairData
}

func (m *storeBackedMetadataStore) GetServiceEndpoints(system string, serviceName string) ([]systemconfig.ServiceEndpointData, error) {
	return nil, nil
}

func (m *storeBackedMetadataStore) GetAllServiceNames(system string) ([]string, error) {
	return nil, nil
}

func (m *storeBackedMetadataStore) GetJavaClassMethods(system string, serviceName string) ([]systemconfig.JavaClassMethodData, error) {
	return nil, nil
}

func (m *storeBackedMetadataStore) GetDatabaseOperations(system string, serviceName string) ([]systemconfig.DatabaseOperationData, error) {
	return nil, nil
}

func (m *storeBackedMetadataStore) GetGRPCOperations(system string, serviceName string) ([]systemconfig.GRPCOperationData, error) {
	return nil, nil
}

func (m *storeBackedMetadataStore) GetNetworkPairs(system string) ([]systemconfig.NetworkPairData, error) {
	return m.pairs[system], nil
}

func (m *storeBackedMetadataStore) GetRuntimeMutatorTargets(system string) ([]systemconfig.RuntimeMutatorTargetData, error) {
	return nil, nil
}

func TestSelectNetworkTargetForService(t *testing.T) {
	// Setup mocks for network dependencies
	cleanup := testdata.SetupNetworkDependenciesMock()
	defer cleanup()

	tests := []struct {
		name           string
		sourceName     string
		targetIndex    int
		wantTargetName string
		wantOk         bool
	}{
		{
			name:           "Valid source and target index",
			sourceName:     "ts-auth-service",
			targetIndex:    0,
			wantTargetName: "ts-verification-code-service", // First dependency from mock
			wantOk:         true,
		},
		{
			name:           "Valid source but negative target index",
			sourceName:     "ts-auth-service",
			targetIndex:    -1,
			wantTargetName: "",
			wantOk:         false,
		},
		{
			name:           "Valid source but out of bounds target index",
			sourceName:     "ts-auth-service",
			targetIndex:    100,
			wantTargetName: "",
			wantOk:         false,
		},
		{
			name:           "Non-existent source service",
			sourceName:     "non-existent-service",
			targetIndex:    0,
			wantTargetName: "",
			wantOk:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetName, ok := networkdependencies.GetServicePairByServiceAndIndex(tt.sourceName, tt.targetIndex)

			if ok != tt.wantOk {
				t.Errorf("GetServicePairByServiceAndIndex() ok = %v, wantOk %v", ok, tt.wantOk)
				return
			}

			if tt.wantOk && targetName != tt.wantTargetName {
				t.Errorf("GetServicePairByServiceAndIndex() targetName = %v, want %v", targetName, tt.wantTargetName)
			}
		})
	}
}

func TestGetAllServiceNames(t *testing.T) {
	// Setup mocks
	cleanup := testdata.SetupNetworkDependenciesMock()
	defer cleanup()

	serviceNames := networkdependencies.ListAllServiceNames()

	if len(serviceNames) == 0 {
		t.Errorf("ListAllServiceNames() returned empty list, expected service names")
	}

	// Check that the list contains the expected service names from mocks
	expectedServices := []string{
		"ts-auth-service",
		"ts-order-service",
		"ts-travel-service",
	}

	for _, expected := range expectedServices {
		found := false
		for _, actual := range serviceNames {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ListAllServiceNames() missing expected service: %s", expected)
		}
	}
}

func TestGetDependenciesForService(t *testing.T) {
	// Setup mocks
	cleanup := testdata.SetupNetworkDependenciesMock()
	defer cleanup()

	tests := []struct {
		name        string
		serviceName string
		wantEmpty   bool
	}{
		{
			name:        "Existing service with dependencies",
			serviceName: "ts-auth-service",
			wantEmpty:   false,
		},
		{
			name:        "Non-existent service",
			serviceName: "non-existent-service",
			wantEmpty:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dependencies := networkdependencies.GetDependenciesForService(tt.serviceName)

			if tt.wantEmpty && len(dependencies) > 0 {
				t.Errorf("GetDependenciesForService() returned %d dependencies, expected empty list", len(dependencies))
			}

			if !tt.wantEmpty && len(dependencies) == 0 {
				t.Errorf("GetDependenciesForService() returned empty list, expected dependencies")
			}
		})
	}
}

func TestGetAllServicePairs(t *testing.T) {
	// Setup mocks
	cleanup := testdata.SetupNetworkDependenciesMock()
	defer cleanup()

	pairs := networkdependencies.GetAllServicePairs()

	if len(pairs) == 0 {
		t.Errorf("GetAllServicePairs() returned empty list, expected service pairs")
	}

	// Verify the structure of the pairs
	for _, pair := range pairs {
		if pair.SourceService == "" {
			t.Errorf("GetAllServicePairs() returned pair with empty source service")
		}

		if pair.TargetService == "" {
			t.Errorf("GetAllServicePairs() returned pair with empty target service")
		}

		if pair.ConnectionDetails == "" {
			t.Errorf("GetAllServicePairs() returned pair with empty connection details")
		}
	}

	// Verify a specific pair exists from our mock data
	foundPair := false
	for _, pair := range pairs {
		if pair.SourceService == "ts-auth-service" &&
			pair.TargetService == "ts-verification-code-service" {
			foundPair = true
			break
		}
	}

	if !foundPair {
		t.Errorf("Expected to find pair ts-auth-service -> ts-verification-code-service")
	}
}

func TestNetworkHelpersIntegration(t *testing.T) {
	// Setup mocks
	cleanup := testdata.SetupNetworkDependenciesMock()
	defer cleanup()

	// Test that our helper functions work well together
	serviceNames := networkdependencies.ListAllServiceNames()
	if len(serviceNames) == 0 {
		t.Fatal("No service names returned")
	}

	sourceName := serviceNames[0]
	dependencies := networkdependencies.GetDependenciesForService(sourceName)

	if len(dependencies) == 0 {
		// Try another service if this one has no dependencies
		if len(serviceNames) > 1 {
			sourceName = serviceNames[1]
			dependencies = networkdependencies.GetDependenciesForService(sourceName)
		}
	}

	if len(dependencies) == 0 {
		t.Skip("No service with dependencies found, skipping integration test")
	}

	// Test that GetServicePairByServiceAndIndex works with the dependencies
	targetName, ok := networkdependencies.GetServicePairByServiceAndIndex(sourceName, 0)
	if !ok {
		t.Fatalf("GetServicePairByServiceAndIndex() failed for %s", sourceName)
	}

	if targetName != dependencies[0] {
		t.Errorf("GetServicePairByServiceAndIndex() targetName = %v, want %v", targetName, dependencies[0])
	}

	// Verify that all pairs contain our source service
	pairs := networkdependencies.GetAllServicePairs()
	foundPair := false

	for _, pair := range pairs {
		if pair.SourceService == sourceName && pair.TargetService == targetName {
			foundPair = true
			break
		}
	}

	if !foundPair {
		t.Errorf("GetAllServicePairs() does not contain expected pair: %s -> %s", sourceName, targetName)
	}
}

func TestMetadataStoreNetworkPairs(t *testing.T) {
	const testSystem = systemconfig.SystemType("network-store-system")

	t.Cleanup(func() {
		systemconfig.SetMetadataStore(nil)
		_ = systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)
		_ = systemconfig.UnregisterSystem(testSystem)
	})

	if err := systemconfig.RegisterSystem(systemconfig.SystemRegistration{
		Name:        testSystem,
		NsPattern:   "^network-store-system\\d+$",
		DisplayName: "NetworkStoreSystem",
	}); err != nil {
		t.Fatalf("RegisterSystem() error = %v", err)
	}

	systemconfig.SetMetadataStore(&storeBackedMetadataStore{
		pairs: map[string][]systemconfig.NetworkPairData{
			string(testSystem): {
				{Source: "frontend", Target: "checkout"},
			},
		},
	})

	if err := systemconfig.SetCurrentSystem(testSystem); err != nil {
		t.Fatalf("SetCurrentSystem() error = %v", err)
	}

	pairs := networkdependencies.GetAllServicePairs()
	if len(pairs) != 1 {
		t.Fatalf("GetAllServicePairs() returned %d pairs, want 1", len(pairs))
	}
	if pairs[0].SourceService != "frontend" || pairs[0].TargetService != "checkout" {
		t.Fatalf("GetAllServicePairs() = %#v", pairs)
	}
}
