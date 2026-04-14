package javamutatorconfig_test

import (
	"testing"

	"github.com/OperationsPAI/chaos-experiment/handler"
	"github.com/OperationsPAI/chaos-experiment/internal/javamutatorconfig"
	"github.com/OperationsPAI/chaos-experiment/internal/resourcelookup"
	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

type testRuntimeMutatorProvider struct {
	services map[string][]systemconfig.RuntimeMutatorTargetData
}

func (p *testRuntimeMutatorProvider) GetServiceNames() []string {
	result := make([]string, 0, len(p.services))
	for service := range p.services {
		result = append(result, service)
	}
	return result
}

func (p *testRuntimeMutatorProvider) GetTargetsByService(serviceName string) []systemconfig.RuntimeMutatorTargetData {
	return p.services[serviceName]
}

type runtimeMutatorMetadataStore struct {
	targets map[string][]systemconfig.RuntimeMutatorTargetData
}

func (m *runtimeMutatorMetadataStore) GetServiceEndpoints(system string, serviceName string) ([]systemconfig.ServiceEndpointData, error) {
	return nil, nil
}

func (m *runtimeMutatorMetadataStore) GetAllServiceNames(system string) ([]string, error) {
	return nil, nil
}

func (m *runtimeMutatorMetadataStore) GetJavaClassMethods(system string, serviceName string) ([]systemconfig.JavaClassMethodData, error) {
	return nil, nil
}

func (m *runtimeMutatorMetadataStore) GetDatabaseOperations(system string, serviceName string) ([]systemconfig.DatabaseOperationData, error) {
	return nil, nil
}

func (m *runtimeMutatorMetadataStore) GetGRPCOperations(system string, serviceName string) ([]systemconfig.GRPCOperationData, error) {
	return nil, nil
}

func (m *runtimeMutatorMetadataStore) GetNetworkPairs(system string) ([]systemconfig.NetworkPairData, error) {
	return nil, nil
}

func (m *runtimeMutatorMetadataStore) GetRuntimeMutatorTargets(system string) ([]systemconfig.RuntimeMutatorTargetData, error) {
	return m.targets[system], nil
}

func TestRegisterRuntimeMutatorProvider(t *testing.T) {
	const systemName = "runtime-mutator-dynamic"
	system := systemconfig.SystemType(systemName)

	t.Cleanup(func() {
		systemconfig.SetMetadataStore(nil)
		_ = systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)
		if handler.IsSystemRegistered(systemName) {
			_ = handler.UnregisterSystem(systemName)
		}
	})

	if err := handler.RegisterSystem(handler.SystemConfig{
		Name:        systemName,
		NsPattern:   "^runtime-mutator-dynamic\\d+$",
		DisplayName: "RuntimeMutatorDynamic",
	}); err != nil {
		t.Fatalf("RegisterSystem() error = %v", err)
	}

	if err := handler.RegisterRuntimeMutatorProvider(systemName, &testRuntimeMutatorProvider{
		services: map[string][]systemconfig.RuntimeMutatorTargetData{
			"checkout": {
				{
					AppName:          "checkout",
					ClassName:        "CheckoutService",
					MethodName:       "PlaceOrder",
					MutationType:     2,
					MutationTypeName: "string",
					MutationStrategy: "uppercase",
					Description:      "Uppercase the return value",
				},
			},
		},
	}); err != nil {
		t.Fatalf("RegisterRuntimeMutatorProvider() error = %v", err)
	}

	if err := systemconfig.SetCurrentSystem(system); err != nil {
		t.Fatalf("SetCurrentSystem() error = %v", err)
	}

	injections := javamutatorconfig.ListAllValidInjections()
	if len(injections) != 1 {
		t.Fatalf("ListAllValidInjections() returned %d targets, want 1", len(injections))
	}
	if injections[0].AppName != "checkout" || injections[0].Mutation.Strategy != "uppercase" {
		t.Fatalf("ListAllValidInjections() = %#v", injections)
	}

	targets, err := resourcelookup.GetSystemCache(system).GetAllJVMRuntimeMutatorTargets()
	if err != nil {
		t.Fatalf("GetAllJVMRuntimeMutatorTargets() error = %v", err)
	}
	if len(targets) != 1 || targets[0].MethodName != "PlaceOrder" {
		t.Fatalf("GetAllJVMRuntimeMutatorTargets() = %#v", targets)
	}
}

func TestMetadataStoreOverridesRuntimeMutatorProvider(t *testing.T) {
	const systemName = "runtime-mutator-store"
	system := systemconfig.SystemType(systemName)

	t.Cleanup(func() {
		systemconfig.SetMetadataStore(nil)
		_ = systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)
		if handler.IsSystemRegistered(systemName) {
			_ = handler.UnregisterSystem(systemName)
		}
	})

	if err := handler.RegisterSystem(handler.SystemConfig{
		Name:        systemName,
		NsPattern:   "^runtime-mutator-store\\d+$",
		DisplayName: "RuntimeMutatorStore",
	}); err != nil {
		t.Fatalf("RegisterSystem() error = %v", err)
	}

	if err := handler.RegisterRuntimeMutatorProvider(systemName, &testRuntimeMutatorProvider{
		services: map[string][]systemconfig.RuntimeMutatorTargetData{
			"provider": {
				{
					AppName:          "provider",
					ClassName:        "ProviderClass",
					MethodName:       "ProviderMethod",
					MutationType:     0,
					MutationTypeName: "constant",
					MutationFrom:     "1",
					MutationTo:       "2",
				},
			},
		},
	}); err != nil {
		t.Fatalf("RegisterRuntimeMutatorProvider() error = %v", err)
	}

	systemconfig.SetMetadataStore(&runtimeMutatorMetadataStore{
		targets: map[string][]systemconfig.RuntimeMutatorTargetData{
			systemName: {
				{
					AppName:          "store",
					ClassName:        "StoreClass",
					MethodName:       "StoreMethod",
					MutationType:     1,
					MutationTypeName: "operator",
					MutationStrategy: "add_to_sub",
				},
			},
		},
	})

	if err := systemconfig.SetCurrentSystem(system); err != nil {
		t.Fatalf("SetCurrentSystem() error = %v", err)
	}

	injections := javamutatorconfig.ListAllValidInjections()
	if len(injections) != 1 {
		t.Fatalf("ListAllValidInjections() returned %d targets, want 1", len(injections))
	}
	if injections[0].AppName != "store" || injections[0].MethodName != "StoreMethod" {
		t.Fatalf("ListAllValidInjections() = %#v, want metadata-store target", injections)
	}
}
