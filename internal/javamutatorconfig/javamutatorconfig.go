// Package javamutatorconfig provides a system-aware routing layer for JVM runtime mutator configs.
package javamutatorconfig

import (
	"sort"

	obmutator "github.com/OperationsPAI/chaos-experiment/internal/ob/mutatorconfig"
	"github.com/OperationsPAI/chaos-experiment/internal/resourcetypes"
	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"

	oteldemomutator "github.com/OperationsPAI/chaos-experiment/internal/oteldemo/mutatorconfig"
	sockshopmutator "github.com/OperationsPAI/chaos-experiment/internal/sockshop/mutatorconfig"
	teastoremutator "github.com/OperationsPAI/chaos-experiment/internal/teastore/mutatorconfig"
	tsmutator "github.com/OperationsPAI/chaos-experiment/internal/ts/mutatorconfig"
)

// MutationSpec is an alias to the shared runtime mutator type.
type MutationSpec = resourcetypes.RuntimeMutatorMutationSpec

// ValidInjection is an alias to the shared flattened runtime mutator target type.
type ValidInjection = resourcetypes.RuntimeMutatorTarget

// ListAllValidInjections returns all valid runtime mutator targets for current system.
func ListAllValidInjections() []ValidInjection {
	services := getAllServicesBySystem()
	result := make([]ValidInjection, 0)

	for _, service := range services {
		entries := getMutatorConfigByService(service)
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
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].AppName != result[j].AppName {
			return result[i].AppName < result[j].AppName
		}
		if result[i].ClassName != result[j].ClassName {
			return result[i].ClassName < result[j].ClassName
		}
		if result[i].MethodName != result[j].MethodName {
			return result[i].MethodName < result[j].MethodName
		}
		if result[i].Mutation.Type != result[j].Mutation.Type {
			return result[i].Mutation.Type < result[j].Mutation.Type
		}
		if result[i].Mutation.Strategy != result[j].Mutation.Strategy {
			return result[i].Mutation.Strategy < result[j].Mutation.Strategy
		}
		if result[i].Mutation.From != result[j].Mutation.From {
			return result[i].Mutation.From < result[j].Mutation.From
		}
		return result[i].Mutation.To < result[j].Mutation.To
	})

	return result
}

func getAllServicesBySystem() []string {
	system := systemconfig.GetCurrentSystem()
	switch system {
	case systemconfig.SystemTrainTicket:
		return tsmutator.GetAllServices()
	case systemconfig.SystemOtelDemo:
		return oteldemomutator.GetAllServices()
	case systemconfig.SystemOnlineBoutique:
		return obmutator.GetAllServices()
	case systemconfig.SystemSockShop:
		return sockshopmutator.GetAllServices()
	case systemconfig.SystemTeaStore:
		return teastoremutator.GetAllServices()
	default:
		return []string{}
	}
}

type methodMutationEntry struct {
	ClassName  string
	MethodName string
	Mutations  []MutationSpec
}

func getMutatorConfigByService(service string) []methodMutationEntry {
	system := systemconfig.GetCurrentSystem()
	switch system {
	case systemconfig.SystemTrainTicket:
		entries := tsmutator.GetMutatorConfigByService(service)
		result := make([]methodMutationEntry, len(entries))
		for i, e := range entries {
			result[i] = methodMutationEntry{
				ClassName:  e.ClassName,
				MethodName: e.MethodName,
				Mutations:  convertTSMutations(e.Mutations),
			}
		}
		return result
	case systemconfig.SystemOtelDemo:
		entries := oteldemomutator.GetMutatorConfigByService(service)
		result := make([]methodMutationEntry, len(entries))
		for i, e := range entries {
			result[i] = methodMutationEntry{
				ClassName:  e.ClassName,
				MethodName: e.MethodName,
				Mutations:  convertOtelMutations(e.Mutations),
			}
		}
		return result
	case systemconfig.SystemOnlineBoutique:
		entries := obmutator.GetMutatorConfigByService(service)
		result := make([]methodMutationEntry, len(entries))
		for i, e := range entries {
			result[i] = methodMutationEntry{
				ClassName:  e.ClassName,
				MethodName: e.MethodName,
				Mutations:  convertOBMutations(e.Mutations),
			}
		}
		return result
	case systemconfig.SystemSockShop:
		entries := sockshopmutator.GetMutatorConfigByService(service)
		result := make([]methodMutationEntry, len(entries))
		for i, e := range entries {
			result[i] = methodMutationEntry{
				ClassName:  e.ClassName,
				MethodName: e.MethodName,
				Mutations:  convertSockShopMutations(e.Mutations),
			}
		}
		return result
	case systemconfig.SystemTeaStore:
		entries := teastoremutator.GetMutatorConfigByService(service)
		result := make([]methodMutationEntry, len(entries))
		for i, e := range entries {
			result[i] = methodMutationEntry{
				ClassName:  e.ClassName,
				MethodName: e.MethodName,
				Mutations:  convertTeaStoreMutations(e.Mutations),
			}
		}
		return result
	default:
		return []methodMutationEntry{}
	}
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
