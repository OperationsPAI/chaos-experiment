// Package javaclassmethods provides a system-aware routing layer for Java class method data.
// This package delegates to registered providers instead of hard-coded switch statements.
package javaclassmethods

import (
	"sort"

	"github.com/OperationsPAI/chaos-experiment/internal/serviceendpoints"
	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"

	objvm "github.com/OperationsPAI/chaos-experiment/internal/ob/javaclassmethods"
	oteldemojvm "github.com/OperationsPAI/chaos-experiment/internal/oteldemo/javaclassmethods"
	sockshopjvm "github.com/OperationsPAI/chaos-experiment/internal/sockshop/javaclassmethods"
	teastorejvm "github.com/OperationsPAI/chaos-experiment/internal/teastore/javaclassmethods"
	tsjvm "github.com/OperationsPAI/chaos-experiment/internal/ts/javaclassmethods"
)

// ClassMethodEntry represents a class-method pair from Java analysis.
type ClassMethodEntry struct {
	ClassName  string
	MethodName string
}

type staticJavaClassMethodProvider struct {
	methods map[string][]ClassMethodEntry
}

func init() {
	registry := systemconfig.GetRegistry()
	registry.RegisterJavaClassMethodProvider(systemconfig.SystemTrainTicket, newStaticJavaClassMethodProvider(convertTSMethodMap()))
	registry.RegisterJavaClassMethodProvider(systemconfig.SystemOtelDemo, newStaticJavaClassMethodProvider(convertOtelDemoMethodMap()))
	registry.RegisterJavaClassMethodProvider(systemconfig.SystemOnlineBoutique, newStaticJavaClassMethodProvider(convertOBMethodMap()))
	registry.RegisterJavaClassMethodProvider(systemconfig.SystemSockShop, newStaticJavaClassMethodProvider(convertSockShopMethodMap()))
	registry.RegisterJavaClassMethodProvider(systemconfig.SystemTeaStore, newStaticJavaClassMethodProvider(convertTeaStoreMethodMap()))
}

func newStaticJavaClassMethodProvider(methods map[string][]ClassMethodEntry) systemconfig.JavaClassMethodProvider {
	return &staticJavaClassMethodProvider{methods: methods}
}

func (p *staticJavaClassMethodProvider) GetServiceNames() []string {
	services := make([]string, 0, len(p.methods))
	for service := range p.methods {
		services = append(services, service)
	}
	sort.Strings(services)
	return services
}

func (p *staticJavaClassMethodProvider) GetClassMethodsByService(serviceName string) []systemconfig.JavaClassMethodData {
	methods := p.methods[serviceName]
	result := make([]systemconfig.JavaClassMethodData, len(methods))
	for i, method := range methods {
		result[i] = systemconfig.JavaClassMethodData{
			ClassName:  method.ClassName,
			MethodName: method.MethodName,
		}
	}
	return result
}

// GetClassMethodsByService returns all class-method pairs for a service based on the current system.
func GetClassMethodsByService(serviceName string) []ClassMethodEntry {
	if !isNetworkServiceName(serviceName) {
		return []ClassMethodEntry{}
	}

	data, err := systemconfig.GetMetadataStore().GetJavaClassMethods(string(systemconfig.GetCurrentSystem()), serviceName)
	if err == nil && len(data) > 0 {
		result := make([]ClassMethodEntry, len(data))
		for i, method := range data {
			result[i] = ClassMethodEntry{
				ClassName:  method.ClassName,
				MethodName: method.MethodName,
			}
		}
		return result
	}
	return []ClassMethodEntry{}
}

// GetAllServices returns a list of all available service names based on the current system.
func GetAllServices() []string {
	networkServices := serviceendpoints.GetAllServices()
	if len(networkServices) == 0 {
		return []string{}
	}

	networkSet := make(map[string]struct{}, len(networkServices))
	for _, service := range networkServices {
		networkSet[service] = struct{}{}
	}

	var candidateServices []string
	names, err := systemconfig.GetMetadataStore().GetAllServiceNames(string(systemconfig.GetCurrentSystem()))
	if err == nil && len(names) > 0 {
		for _, service := range names {
			if len(GetClassMethodsByService(service)) > 0 {
				candidateServices = append(candidateServices, service)
			}
		}
	}

	var filtered []string
	for _, service := range candidateServices {
		if _, ok := networkSet[service]; ok {
			filtered = append(filtered, service)
		}
	}
	sort.Strings(filtered)
	return filtered
}

func isNetworkServiceName(serviceName string) bool {
	if serviceName == "" {
		return false
	}
	for _, networkService := range serviceendpoints.GetAllServices() {
		if networkService == serviceName {
			return true
		}
	}
	return false
}

func convertTSMethodMap() map[string][]ClassMethodEntry {
	return buildMethodMap(
		tsjvm.GetAllServices,
		func(service string) []ClassMethodEntry {
			return convertTSMethods(tsjvm.GetClassMethodsByService(service))
		},
	)
}

func convertOtelDemoMethodMap() map[string][]ClassMethodEntry {
	return buildMethodMap(
		oteldemojvm.GetAllServices,
		func(service string) []ClassMethodEntry {
			return convertOtelDemoMethods(oteldemojvm.GetClassMethodsByService(service))
		},
	)
}

func convertOBMethodMap() map[string][]ClassMethodEntry {
	return buildMethodMap(
		objvm.GetAllServices,
		func(service string) []ClassMethodEntry {
			return convertOBMethods(objvm.GetClassMethodsByService(service))
		},
	)
}

func convertSockShopMethodMap() map[string][]ClassMethodEntry {
	return buildMethodMap(
		sockshopjvm.GetAllServices,
		func(service string) []ClassMethodEntry {
			return convertSockShopMethods(sockshopjvm.GetClassMethodsByService(service))
		},
	)
}

func convertTeaStoreMethodMap() map[string][]ClassMethodEntry {
	return buildMethodMap(
		teastorejvm.GetAllServices,
		func(service string) []ClassMethodEntry {
			return convertTeaStoreMethods(teastorejvm.GetClassMethodsByService(service))
		},
	)
}

func buildMethodMap(services func() []string, loader func(string) []ClassMethodEntry) map[string][]ClassMethodEntry {
	allServices := services()
	result := make(map[string][]ClassMethodEntry, len(allServices))
	for _, service := range allServices {
		result[service] = loader(service)
	}
	return result
}

// convertTSMethods converts ts-specific methods to the common type.
func convertTSMethods(tsMethods []tsjvm.ClassMethodEntry) []ClassMethodEntry {
	result := make([]ClassMethodEntry, len(tsMethods))
	for i, m := range tsMethods {
		result[i] = ClassMethodEntry{
			ClassName:  m.ClassName,
			MethodName: m.MethodName,
		}
	}
	return result
}

// convertOtelDemoMethods converts otel-demo-specific methods to the common type.
func convertOtelDemoMethods(otelMethods []oteldemojvm.ClassMethodEntry) []ClassMethodEntry {
	result := make([]ClassMethodEntry, len(otelMethods))
	for i, m := range otelMethods {
		result[i] = ClassMethodEntry{
			ClassName:  m.ClassName,
			MethodName: m.MethodName,
		}
	}
	return result
}

// convertOBMethods converts ob-specific methods to the common type.
func convertOBMethods(obMethods []objvm.ClassMethodEntry) []ClassMethodEntry {
	result := make([]ClassMethodEntry, len(obMethods))
	for i, m := range obMethods {
		result[i] = ClassMethodEntry{
			ClassName:  m.ClassName,
			MethodName: m.MethodName,
		}
	}
	return result
}

// convertSockShopMethods converts sockshop-specific methods to the common type.
func convertSockShopMethods(sockshopMethods []sockshopjvm.ClassMethodEntry) []ClassMethodEntry {
	result := make([]ClassMethodEntry, len(sockshopMethods))
	for i, m := range sockshopMethods {
		result[i] = ClassMethodEntry{
			ClassName:  m.ClassName,
			MethodName: m.MethodName,
		}
	}
	return result
}

// convertTeaStoreMethods converts teastore-specific methods to the common type.
func convertTeaStoreMethods(teastoreMethods []teastorejvm.ClassMethodEntry) []ClassMethodEntry {
	result := make([]ClassMethodEntry, len(teastoreMethods))
	for i, m := range teastoreMethods {
		result[i] = ClassMethodEntry{
			ClassName:  m.ClassName,
			MethodName: m.MethodName,
		}
	}
	return result
}
