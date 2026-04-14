// Package serviceendpoints provides a system-aware routing layer for service endpoint data.
// This package delegates to registered providers instead of hard-coded switch statements.
package serviceendpoints

import (
	"sort"

	"github.com/LGU-SE-Internal/chaos-experiment/internal/systemconfig"

	hsendpoints "github.com/LGU-SE-Internal/chaos-experiment/internal/hs/serviceendpoints"
	mediaendpoints "github.com/LGU-SE-Internal/chaos-experiment/internal/media/serviceendpoints"
	obendpoints "github.com/LGU-SE-Internal/chaos-experiment/internal/ob/serviceendpoints"
	oteldemoendpoints "github.com/LGU-SE-Internal/chaos-experiment/internal/oteldemo/serviceendpoints"
	snendpoints "github.com/LGU-SE-Internal/chaos-experiment/internal/sn/serviceendpoints"
	sockshopendpoints "github.com/LGU-SE-Internal/chaos-experiment/internal/sockshop/serviceendpoints"
	teastoreendpoints "github.com/LGU-SE-Internal/chaos-experiment/internal/teastore/serviceendpoints"
	tsendpoints "github.com/LGU-SE-Internal/chaos-experiment/internal/ts/serviceendpoints"
)

// ServiceEndpoint represents a service endpoint from ClickHouse analysis.
type ServiceEndpoint struct {
	ServiceName    string
	RequestMethod  string
	ResponseStatus string
	Route          string
	ServerAddress  string
	ServerPort     string
	SpanName       string
}

type staticServiceEndpointProvider struct {
	endpoints map[string][]ServiceEndpoint
}

func init() {
	registry := systemconfig.GetRegistry()
	registry.RegisterServiceEndpointProvider(systemconfig.SystemTrainTicket, newStaticServiceEndpointProvider(convertTSEndpointMap(tsendpoints.ServiceEndpoints)))
	registry.RegisterServiceEndpointProvider(systemconfig.SystemOtelDemo, newStaticServiceEndpointProvider(convertOtelDemoEndpointMap(oteldemoendpoints.ServiceEndpoints)))
	registry.RegisterServiceEndpointProvider(systemconfig.SystemMediaMicroservices, newStaticServiceEndpointProvider(convertMediaEndpointMap(mediaendpoints.ServiceEndpoints)))
	registry.RegisterServiceEndpointProvider(systemconfig.SystemHotelReservation, newStaticServiceEndpointProvider(convertHSEndpointMap(hsendpoints.ServiceEndpoints)))
	registry.RegisterServiceEndpointProvider(systemconfig.SystemSocialNetwork, newStaticServiceEndpointProvider(convertSNEndpointMap(snendpoints.ServiceEndpoints)))
	registry.RegisterServiceEndpointProvider(systemconfig.SystemOnlineBoutique, newStaticServiceEndpointProvider(convertOBEndpointMap(obendpoints.ServiceEndpoints)))
	registry.RegisterServiceEndpointProvider(systemconfig.SystemSockShop, newStaticServiceEndpointProvider(convertSockShopEndpointMap(sockshopendpoints.ServiceEndpoints)))
	registry.RegisterServiceEndpointProvider(systemconfig.SystemTeaStore, newStaticServiceEndpointProvider(convertTeaStoreEndpointMap(teastoreendpoints.ServiceEndpoints)))
}

func newStaticServiceEndpointProvider(endpoints map[string][]ServiceEndpoint) systemconfig.ServiceEndpointProvider {
	return &staticServiceEndpointProvider{endpoints: endpoints}
}

func (p *staticServiceEndpointProvider) GetServiceNames() []string {
	services := make([]string, 0, len(p.endpoints))
	for service := range p.endpoints {
		services = append(services, service)
	}
	sort.Strings(services)
	return services
}

func (p *staticServiceEndpointProvider) GetEndpointsByService(serviceName string) []systemconfig.ServiceEndpointData {
	endpoints := p.endpoints[serviceName]
	result := make([]systemconfig.ServiceEndpointData, len(endpoints))
	for i, endpoint := range endpoints {
		result[i] = systemconfig.ServiceEndpointData{
			ServiceName:    endpoint.ServiceName,
			RequestMethod:  endpoint.RequestMethod,
			ResponseStatus: endpoint.ResponseStatus,
			Route:          endpoint.Route,
			ServerAddress:  endpoint.ServerAddress,
			ServerPort:     endpoint.ServerPort,
			SpanName:       endpoint.SpanName,
		}
	}
	return result
}

// GetEndpointsByService returns all endpoints for a service based on the current system.
func GetEndpointsByService(serviceName string) []ServiceEndpoint {
	provider, err := systemconfig.GetRegistry().GetServiceEndpointProvider()
	if err != nil {
		return []ServiceEndpoint{}
	}

	data := provider.GetEndpointsByService(serviceName)
	result := make([]ServiceEndpoint, len(data))
	for i, endpoint := range data {
		result[i] = ServiceEndpoint{
			ServiceName:    endpoint.ServiceName,
			RequestMethod:  endpoint.RequestMethod,
			ResponseStatus: endpoint.ResponseStatus,
			Route:          endpoint.Route,
			ServerAddress:  endpoint.ServerAddress,
			ServerPort:     endpoint.ServerPort,
			SpanName:       endpoint.SpanName,
		}
	}
	return result
}

// GetAllServices returns a list of all available service names based on the current system.
func GetAllServices() []string {
	provider, err := systemconfig.GetRegistry().GetServiceEndpointProvider()
	if err != nil {
		return []string{}
	}
	return provider.GetServiceNames()
}

func convertTSEndpointMap(tsEps map[string][]tsendpoints.ServiceEndpoint) map[string][]ServiceEndpoint {
	result := make(map[string][]ServiceEndpoint, len(tsEps))
	for service, endpoints := range tsEps {
		result[service] = convertTSEndpoints(endpoints)
	}
	return result
}

func convertOtelDemoEndpointMap(otelEps map[string][]oteldemoendpoints.ServiceEndpoint) map[string][]ServiceEndpoint {
	result := make(map[string][]ServiceEndpoint, len(otelEps))
	for service, endpoints := range otelEps {
		result[service] = convertOtelDemoEndpoints(endpoints)
	}
	return result
}

func convertMediaEndpointMap(mediaEps map[string][]mediaendpoints.ServiceEndpoint) map[string][]ServiceEndpoint {
	result := make(map[string][]ServiceEndpoint, len(mediaEps))
	for service, endpoints := range mediaEps {
		result[service] = convertMediaEndpoints(endpoints)
	}
	return result
}

func convertHSEndpointMap(hsEps map[string][]hsendpoints.ServiceEndpoint) map[string][]ServiceEndpoint {
	result := make(map[string][]ServiceEndpoint, len(hsEps))
	for service, endpoints := range hsEps {
		result[service] = convertHSEndpoints(endpoints)
	}
	return result
}

func convertSNEndpointMap(snEps map[string][]snendpoints.ServiceEndpoint) map[string][]ServiceEndpoint {
	result := make(map[string][]ServiceEndpoint, len(snEps))
	for service, endpoints := range snEps {
		result[service] = convertSNEndpoints(endpoints)
	}
	return result
}

func convertOBEndpointMap(obEps map[string][]obendpoints.ServiceEndpoint) map[string][]ServiceEndpoint {
	result := make(map[string][]ServiceEndpoint, len(obEps))
	for service, endpoints := range obEps {
		result[service] = convertOBEndpoints(endpoints)
	}
	return result
}

func convertSockShopEndpointMap(sockshopEps map[string][]sockshopendpoints.ServiceEndpoint) map[string][]ServiceEndpoint {
	result := make(map[string][]ServiceEndpoint, len(sockshopEps))
	for service, endpoints := range sockshopEps {
		result[service] = convertSockShopEndpoints(endpoints)
	}
	return result
}

func convertTeaStoreEndpointMap(teastoreEps map[string][]teastoreendpoints.ServiceEndpoint) map[string][]ServiceEndpoint {
	result := make(map[string][]ServiceEndpoint, len(teastoreEps))
	for service, endpoints := range teastoreEps {
		result[service] = convertTeaStoreEndpoints(endpoints)
	}
	return result
}

// convertTSEndpoints converts ts-specific endpoints to the common type.
func convertTSEndpoints(tsEps []tsendpoints.ServiceEndpoint) []ServiceEndpoint {
	result := make([]ServiceEndpoint, len(tsEps))
	for i, ep := range tsEps {
		result[i] = ServiceEndpoint{
			ServiceName:    ep.ServiceName,
			RequestMethod:  ep.RequestMethod,
			ResponseStatus: ep.ResponseStatus,
			Route:          ep.Route,
			ServerAddress:  ep.ServerAddress,
			ServerPort:     ep.ServerPort,
			SpanName:       ep.SpanName,
		}
	}
	return result
}

// convertOtelDemoEndpoints converts otel-demo-specific endpoints to the common type.
func convertOtelDemoEndpoints(otelEps []oteldemoendpoints.ServiceEndpoint) []ServiceEndpoint {
	result := make([]ServiceEndpoint, len(otelEps))
	for i, ep := range otelEps {
		result[i] = ServiceEndpoint{
			ServiceName:    ep.ServiceName,
			RequestMethod:  ep.RequestMethod,
			ResponseStatus: ep.ResponseStatus,
			Route:          ep.Route,
			ServerAddress:  ep.ServerAddress,
			ServerPort:     ep.ServerPort,
			SpanName:       ep.SpanName,
		}
	}
	return result
}

// convertMediaEndpoints converts media-specific endpoints to the common type.
func convertMediaEndpoints(mediaEps []mediaendpoints.ServiceEndpoint) []ServiceEndpoint {
	result := make([]ServiceEndpoint, len(mediaEps))
	for i, ep := range mediaEps {
		result[i] = ServiceEndpoint{
			ServiceName:    ep.ServiceName,
			RequestMethod:  ep.RequestMethod,
			ResponseStatus: ep.ResponseStatus,
			Route:          ep.Route,
			ServerAddress:  ep.ServerAddress,
			ServerPort:     ep.ServerPort,
			SpanName:       ep.SpanName,
		}
	}
	return result
}

// convertHSEndpoints converts hs-specific endpoints to the common type.
func convertHSEndpoints(hsEps []hsendpoints.ServiceEndpoint) []ServiceEndpoint {
	result := make([]ServiceEndpoint, len(hsEps))
	for i, ep := range hsEps {
		result[i] = ServiceEndpoint{
			ServiceName:    ep.ServiceName,
			RequestMethod:  ep.RequestMethod,
			ResponseStatus: ep.ResponseStatus,
			Route:          ep.Route,
			ServerAddress:  ep.ServerAddress,
			ServerPort:     ep.ServerPort,
			SpanName:       ep.SpanName,
		}
	}
	return result
}

// convertSNEndpoints converts sn-specific endpoints to the common type.
func convertSNEndpoints(snEps []snendpoints.ServiceEndpoint) []ServiceEndpoint {
	result := make([]ServiceEndpoint, len(snEps))
	for i, ep := range snEps {
		result[i] = ServiceEndpoint{
			ServiceName:    ep.ServiceName,
			RequestMethod:  ep.RequestMethod,
			ResponseStatus: ep.ResponseStatus,
			Route:          ep.Route,
			ServerAddress:  ep.ServerAddress,
			ServerPort:     ep.ServerPort,
			SpanName:       ep.SpanName,
		}
	}
	return result
}

// convertOBEndpoints converts ob-specific endpoints to the common type.
func convertOBEndpoints(obEps []obendpoints.ServiceEndpoint) []ServiceEndpoint {
	result := make([]ServiceEndpoint, len(obEps))
	for i, ep := range obEps {
		result[i] = ServiceEndpoint{
			ServiceName:    ep.ServiceName,
			RequestMethod:  ep.RequestMethod,
			ResponseStatus: ep.ResponseStatus,
			Route:          ep.Route,
			ServerAddress:  ep.ServerAddress,
			ServerPort:     ep.ServerPort,
			SpanName:       ep.SpanName,
		}
	}
	return result
}

// convertSockShopEndpoints converts sockshop-specific endpoints to the common type.
func convertSockShopEndpoints(sockshopEps []sockshopendpoints.ServiceEndpoint) []ServiceEndpoint {
	result := make([]ServiceEndpoint, len(sockshopEps))
	for i, ep := range sockshopEps {
		result[i] = ServiceEndpoint{
			ServiceName:    ep.ServiceName,
			RequestMethod:  ep.RequestMethod,
			ResponseStatus: ep.ResponseStatus,
			Route:          ep.Route,
			ServerAddress:  ep.ServerAddress,
			ServerPort:     ep.ServerPort,
			SpanName:       ep.SpanName,
		}
	}
	return result
}

// convertTeaStoreEndpoints converts teastore-specific endpoints to the common type.
func convertTeaStoreEndpoints(teastoreEps []teastoreendpoints.ServiceEndpoint) []ServiceEndpoint {
	result := make([]ServiceEndpoint, len(teastoreEps))
	for i, ep := range teastoreEps {
		result[i] = ServiceEndpoint{
			ServiceName:    ep.ServiceName,
			RequestMethod:  ep.RequestMethod,
			ResponseStatus: ep.ResponseStatus,
			Route:          ep.Route,
			ServerAddress:  ep.ServerAddress,
			ServerPort:     ep.ServerPort,
			SpanName:       ep.SpanName,
		}
	}
	return result
}
