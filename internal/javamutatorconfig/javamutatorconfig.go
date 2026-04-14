package javamutatorconfig

import (
	"sort"

	obmutator "github.com/OperationsPAI/chaos-experiment/internal/ob/mutatorconfig"
	oteldemomutator "github.com/OperationsPAI/chaos-experiment/internal/oteldemo/mutatorconfig"
	"github.com/OperationsPAI/chaos-experiment/internal/resourcetypes"
	sockshopmutator "github.com/OperationsPAI/chaos-experiment/internal/sockshop/mutatorconfig"
	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
	teastoremutator "github.com/OperationsPAI/chaos-experiment/internal/teastore/mutatorconfig"
	tsmutator "github.com/OperationsPAI/chaos-experiment/internal/ts/mutatorconfig"
)

// MutationSpec is an alias to the shared runtime mutator type.
type MutationSpec = resourcetypes.RuntimeMutatorMutationSpec

// ValidInjection is an alias to the shared flattened runtime mutator target type.
type ValidInjection = resourcetypes.RuntimeMutatorTarget

type staticRuntimeMutatorProvider struct {
	targets map[string][]ValidInjection
}

type methodMutationEntry struct {
	ClassName  string
	MethodName string
	Mutations  []MutationSpec
}

func init() {
	registry := systemconfig.GetRegistry()
	registry.RegisterRuntimeMutatorProvider(systemconfig.SystemTrainTicket, newStaticRuntimeMutatorProvider(convertTSTargetMap()))
	registry.RegisterRuntimeMutatorProvider(systemconfig.SystemOtelDemo, newStaticRuntimeMutatorProvider(convertOtelDemoTargetMap()))
	registry.RegisterRuntimeMutatorProvider(systemconfig.SystemOnlineBoutique, newStaticRuntimeMutatorProvider(convertOBTargetMap()))
	registry.RegisterRuntimeMutatorProvider(systemconfig.SystemSockShop, newStaticRuntimeMutatorProvider(convertSockShopTargetMap()))
	registry.RegisterRuntimeMutatorProvider(systemconfig.SystemTeaStore, newStaticRuntimeMutatorProvider(convertTeaStoreTargetMap()))
}

func newStaticRuntimeMutatorProvider(targets map[string][]ValidInjection) systemconfig.RuntimeMutatorProvider {
	return &staticRuntimeMutatorProvider{targets: targets}
}

func (p *staticRuntimeMutatorProvider) GetServiceNames() []string {
	services := make([]string, 0, len(p.targets))
	for service := range p.targets {
		services = append(services, service)
	}
	sort.Strings(services)
	return services
}

func (p *staticRuntimeMutatorProvider) GetTargetsByService(serviceName string) []systemconfig.RuntimeMutatorTargetData {
	targets := p.targets[serviceName]
	result := make([]systemconfig.RuntimeMutatorTargetData, len(targets))
	for i, target := range targets {
		result[i] = systemconfig.RuntimeMutatorTargetData{
			AppName:          target.AppName,
			ClassName:        target.ClassName,
			MethodName:       target.MethodName,
			MutationType:     target.Mutation.Type,
			MutationTypeName: target.Mutation.TypeName,
			MutationFrom:     target.Mutation.From,
			MutationTo:       target.Mutation.To,
			MutationStrategy: target.Mutation.Strategy,
			Description:      target.Mutation.Description,
		}
	}
	return result
}

// ListAllValidInjections returns all valid runtime mutator targets for current system.
func ListAllValidInjections() []ValidInjection {
	data, err := systemconfig.GetMetadataStore().GetRuntimeMutatorTargets(string(systemconfig.GetCurrentSystem()))
	if err == nil && len(data) > 0 {
		result := make([]ValidInjection, len(data))
		for i, target := range data {
			result[i] = ValidInjection{
				AppName:    target.AppName,
				ClassName:  target.ClassName,
				MethodName: target.MethodName,
				Mutation: MutationSpec{
					Type:        target.MutationType,
					TypeName:    target.MutationTypeName,
					From:        target.MutationFrom,
					To:          target.MutationTo,
					Strategy:    target.MutationStrategy,
					Description: target.Description,
				},
			}
		}
		return result
	}
	return []ValidInjection{}
}

func convertTSTargetMap() map[string][]ValidInjection {
	return buildTargetMap(
		tsmutator.GetAllServices,
		func(service string) []ValidInjection {
			return convertTSMethodMutationEntries(service, tsmutator.GetMutatorConfigByService(service))
		},
	)
}

func convertOtelDemoTargetMap() map[string][]ValidInjection {
	return buildTargetMap(
		oteldemomutator.GetAllServices,
		func(service string) []ValidInjection {
			return convertOtelDemoMethodMutationEntries(service, oteldemomutator.GetMutatorConfigByService(service))
		},
	)
}

func convertOBTargetMap() map[string][]ValidInjection {
	return buildTargetMap(
		obmutator.GetAllServices,
		func(service string) []ValidInjection {
			return convertOBMethodMutationEntries(service, obmutator.GetMutatorConfigByService(service))
		},
	)
}

func convertSockShopTargetMap() map[string][]ValidInjection {
	return buildTargetMap(
		sockshopmutator.GetAllServices,
		func(service string) []ValidInjection {
			return convertSockShopMethodMutationEntries(service, sockshopmutator.GetMutatorConfigByService(service))
		},
	)
}

func convertTeaStoreTargetMap() map[string][]ValidInjection {
	return buildTargetMap(
		teastoremutator.GetAllServices,
		func(service string) []ValidInjection {
			return convertTeaStoreMethodMutationEntries(service, teastoremutator.GetMutatorConfigByService(service))
		},
	)
}

func buildTargetMap(services func() []string, loader func(string) []ValidInjection) map[string][]ValidInjection {
	allServices := services()
	result := make(map[string][]ValidInjection, len(allServices))
	for _, service := range allServices {
		result[service] = loader(service)
	}
	return result
}

func buildValidInjections(service string, entries []methodMutationEntry) []ValidInjection {
	result := make([]ValidInjection, 0)
	for _, entry := range entries {
		for _, mutation := range entry.Mutations {
			result = append(result, ValidInjection{
				AppName:    service,
				ClassName:  entry.ClassName,
				MethodName: entry.MethodName,
				Mutation:   mutation,
			})
		}
	}
	return result
}

func convertTSMethodMutationEntries(service string, entries []tsmutator.MethodMutationEntry) []ValidInjection {
	result := make([]methodMutationEntry, len(entries))
	for i, entry := range entries {
		result[i] = methodMutationEntry{
			ClassName:  entry.ClassName,
			MethodName: entry.MethodName,
			Mutations:  convertTSMutations(entry.Mutations),
		}
	}
	return buildValidInjections(service, result)
}

func convertOtelDemoMethodMutationEntries(service string, entries []oteldemomutator.MethodMutationEntry) []ValidInjection {
	result := make([]methodMutationEntry, len(entries))
	for i, entry := range entries {
		result[i] = methodMutationEntry{
			ClassName:  entry.ClassName,
			MethodName: entry.MethodName,
			Mutations:  convertOtelMutations(entry.Mutations),
		}
	}
	return buildValidInjections(service, result)
}

func convertOBMethodMutationEntries(service string, entries []obmutator.MethodMutationEntry) []ValidInjection {
	result := make([]methodMutationEntry, len(entries))
	for i, entry := range entries {
		result[i] = methodMutationEntry{
			ClassName:  entry.ClassName,
			MethodName: entry.MethodName,
			Mutations:  convertOBMutations(entry.Mutations),
		}
	}
	return buildValidInjections(service, result)
}

func convertSockShopMethodMutationEntries(service string, entries []sockshopmutator.MethodMutationEntry) []ValidInjection {
	result := make([]methodMutationEntry, len(entries))
	for i, entry := range entries {
		result[i] = methodMutationEntry{
			ClassName:  entry.ClassName,
			MethodName: entry.MethodName,
			Mutations:  convertSockShopMutations(entry.Mutations),
		}
	}
	return buildValidInjections(service, result)
}

func convertTeaStoreMethodMutationEntries(service string, entries []teastoremutator.MethodMutationEntry) []ValidInjection {
	result := make([]methodMutationEntry, len(entries))
	for i, entry := range entries {
		result[i] = methodMutationEntry{
			ClassName:  entry.ClassName,
			MethodName: entry.MethodName,
			Mutations:  convertTeaStoreMutations(entry.Mutations),
		}
	}
	return buildValidInjections(service, result)
}

func convertTSMutations(in []tsmutator.MutationSpec) []MutationSpec {
	out := make([]MutationSpec, len(in))
	for i, m := range in {
		out[i] = MutationSpec{
			Type:        m.Type,
			TypeName:    m.TypeName,
			From:        m.From,
			To:          m.To,
			Strategy:    m.Strategy,
			Description: m.Description,
		}
	}
	return out
}

func convertOtelMutations(in []oteldemomutator.MutationSpec) []MutationSpec {
	out := make([]MutationSpec, len(in))
	for i, m := range in {
		out[i] = MutationSpec{
			Type:        m.Type,
			TypeName:    m.TypeName,
			From:        m.From,
			To:          m.To,
			Strategy:    m.Strategy,
			Description: m.Description,
		}
	}
	return out
}

func convertOBMutations(in []obmutator.MutationSpec) []MutationSpec {
	out := make([]MutationSpec, len(in))
	for i, m := range in {
		out[i] = MutationSpec{
			Type:        m.Type,
			TypeName:    m.TypeName,
			From:        m.From,
			To:          m.To,
			Strategy:    m.Strategy,
			Description: m.Description,
		}
	}
	return out
}

func convertSockShopMutations(in []sockshopmutator.MutationSpec) []MutationSpec {
	out := make([]MutationSpec, len(in))
	for i, m := range in {
		out[i] = MutationSpec{
			Type:        m.Type,
			TypeName:    m.TypeName,
			From:        m.From,
			To:          m.To,
			Strategy:    m.Strategy,
			Description: m.Description,
		}
	}
	return out
}

func convertTeaStoreMutations(in []teastoremutator.MutationSpec) []MutationSpec {
	out := make([]MutationSpec, len(in))
	for i, m := range in {
		out[i] = MutationSpec{
			Type:        m.Type,
			TypeName:    m.TypeName,
			From:        m.From,
			To:          m.To,
			Strategy:    m.Strategy,
			Description: m.Description,
		}
	}
	return out
}
