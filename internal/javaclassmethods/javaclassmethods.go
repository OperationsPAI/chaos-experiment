// Package javaclassmethods provides a system-aware routing layer for Java class method data.
// This package delegates to the appropriate system-specific package based on the current system configuration.
package javaclassmethods

import (
	"sort"

	"github.com/LGU-SE-Internal/chaos-experiment/internal/serviceendpoints"
	"github.com/LGU-SE-Internal/chaos-experiment/internal/systemconfig"

	objvm "github.com/LGU-SE-Internal/chaos-experiment/internal/ob/javaclassmethods"
	oteldemojvm "github.com/LGU-SE-Internal/chaos-experiment/internal/oteldemo/javaclassmethods"
	sockshopjvm "github.com/LGU-SE-Internal/chaos-experiment/internal/sockshop/javaclassmethods"
	teastorejvm "github.com/LGU-SE-Internal/chaos-experiment/internal/teastore/javaclassmethods"
	tsjvm "github.com/LGU-SE-Internal/chaos-experiment/internal/ts/javaclassmethods"
)

// ClassMethodEntry represents a class-method pair from Java analysis
type ClassMethodEntry struct {
	ClassName  string
	MethodName string
}

// GetClassMethodsByService returns all class-method pairs for a service based on current system
func GetClassMethodsByService(serviceName string) []ClassMethodEntry {
	if !isNetworkServiceName(serviceName) {
		return []ClassMethodEntry{}
	}

	system := systemconfig.GetCurrentSystem()
	switch system {
	case systemconfig.SystemTrainTicket:
		tsMethods := tsjvm.GetClassMethodsByService(serviceName)
		return convertTSMethods(tsMethods)
	case systemconfig.SystemOtelDemo:
		otelMethods := oteldemojvm.GetClassMethodsByService(serviceName)
		return convertOtelDemoMethods(otelMethods)
	case systemconfig.SystemOnlineBoutique:
		obMethods := objvm.GetClassMethodsByService(serviceName)
		return convertOBMethods(obMethods)
	case systemconfig.SystemSockShop:
		sockshopMethods := sockshopjvm.GetClassMethodsByService(serviceName)
		return convertSockShopMethods(sockshopMethods)
	case systemconfig.SystemTeaStore:
		teastoreMethods := teastorejvm.GetClassMethodsByService(serviceName)
		return convertTeaStoreMethods(teastoreMethods)
	default:
		return []ClassMethodEntry{}
	}
}

// GetAllServices returns a list of all available service names based on current system
func GetAllServices() []string {
	networkServices := serviceendpoints.GetAllServices()
	if len(networkServices) == 0 {
		return []string{}
	}

	networkSet := make(map[string]struct{}, len(networkServices))
	for _, name := range networkServices {
		networkSet[name] = struct{}{}
	}

	system := systemconfig.GetCurrentSystem()
	var jvmServices []string
	switch system {
	case systemconfig.SystemTrainTicket:
		jvmServices = tsjvm.GetAllServices()
	case systemconfig.SystemOtelDemo:
		jvmServices = oteldemojvm.GetAllServices()
	case systemconfig.SystemOnlineBoutique:
		jvmServices = objvm.GetAllServices()
	case systemconfig.SystemSockShop:
		jvmServices = sockshopjvm.GetAllServices()
	case systemconfig.SystemTeaStore:
		jvmServices = teastorejvm.GetAllServices()
	default:
		return []string{}
	}

	filtered := make([]string, 0, len(jvmServices))
	for _, service := range jvmServices {
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

// convertTSMethods converts ts-specific methods to the common type
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

// convertOtelDemoMethods converts otel-demo-specific methods to the common type
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

// convertOBMethods converts ob-specific methods to the common type
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

// convertSockShopMethods converts sockshop-specific methods to the common type
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

// convertTeaStoreMethods converts teastore-specific methods to the common type
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
