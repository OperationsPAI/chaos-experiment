package handler

import (
	"fmt"
	"sort"

	_ "github.com/LGU-SE-Internal/chaos-experiment/internal/adapter" // ensure system data is registered
	"github.com/LGU-SE-Internal/chaos-experiment/internal/endpoint"
	"github.com/LGU-SE-Internal/chaos-experiment/internal/javaclassmethods"
	"github.com/LGU-SE-Internal/chaos-experiment/internal/registry"
	"github.com/LGU-SE-Internal/chaos-experiment/internal/resourcelookup"
	"github.com/LGU-SE-Internal/chaos-experiment/internal/systemconfig"
)

type JVMMethodInfo struct {
	ServiceName string
	ClassName   string
	MethodName  string
}

func getAllAppLabels(namespace string) ([]string, error) {
	labels, err := resourcelookup.GetAllAppLabels(namespace, TargetLabelKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get app labels: %w", err)
	}
	return labels, nil
}

func getAppLabelByIndex(namespace string, appIdx int) (string, error) {
	labels, err := getAllAppLabels(namespace)
	if err != nil {
		return "", err
	}

	if appIdx < 0 || appIdx >= len(labels) {
		return "", fmt.Errorf("app index out of range: %d (max: %d)", appIdx, len(labels)-1)
	}

	return labels[appIdx], nil
}

func getAllContainerInfos(namespace string) ([]resourcelookup.ContainerInfo, error) {
	containers, err := resourcelookup.GetAllContainers(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get containers: %w", err)
	}
	return containers, nil
}

func getContainerInfoByIndex(namespace string, containerIdx int) (*resourcelookup.ContainerInfo, error) {
	containers, err := getAllContainerInfos(namespace)
	if err != nil {
		return nil, err
	}

	if containerIdx < 0 || containerIdx >= len(containers) {
		return nil, fmt.Errorf("container index out of range: %d (max: %d)", containerIdx, len(containers)-1)
	}

	return &containers[containerIdx], nil
}

func getAllJVMMethodInfos() ([]JVMMethodInfo, error) {
	sysData := registry.GetCurrent()
	if sysData == nil {
		return nil, fmt.Errorf("current system %s is not registered", systemconfig.GetCurrentSystem())
	}

	result := make([]JVMMethodInfo, 0)
	for _, service := range sysData.GetAllServices() {
		methods := javaclassmethods.GetClassMethodsByService(service)
		for _, method := range methods {
			result = append(result, JVMMethodInfo{
				ServiceName: service,
				ClassName:   method.ClassName,
				MethodName:  method.MethodName,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].ServiceName != result[j].ServiceName {
			return result[i].ServiceName < result[j].ServiceName
		}
		if result[i].ClassName != result[j].ClassName {
			return result[i].ClassName < result[j].ClassName
		}
		return result[i].MethodName < result[j].MethodName
	})

	return result, nil
}

func getJVMMethodInfoByIndex(methodIdx int) (*JVMMethodInfo, error) {
	methods, err := getAllJVMMethodInfos()
	if err != nil {
		return nil, fmt.Errorf("failed to get JVM methods: %w", err)
	}

	if methodIdx < 0 || methodIdx >= len(methods) {
		return nil, fmt.Errorf("method index out of range: %d (max: %d)", methodIdx, len(methods)-1)
	}

	return &methods[methodIdx], nil
}

func getAllHTTPEndpointInfos() ([]endpoint.HTTPEndpointInfo, error) {
	sysData := registry.GetCurrent()
	if sysData == nil {
		return nil, fmt.Errorf("current system %s is not registered", systemconfig.GetCurrentSystem())
	}

	result := make([]endpoint.HTTPEndpointInfo, 0)
	for _, service := range sysData.GetAllServices() {
		for _, ep := range sysData.GetHTTPEndpointsByService(service) {
			if ep.ServerAddress == "ts-rabbitmq" || ep.Route == "" {
				continue
			}
			result = append(result, endpoint.ToHTTPEndpointInfo(ep))
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].ServiceName != result[j].ServiceName {
			return result[i].ServiceName < result[j].ServiceName
		}
		return result[i].Route < result[j].Route
	})

	return result, nil
}

func getHTTPChaosEndpointByIndex(endpointIdx int) (*HTTPEndpoint, error) {
	endpoints, err := getAllHTTPEndpointInfos()
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP endpoints: %w", err)
	}

	if endpointIdx < 0 || endpointIdx >= len(endpoints) {
		return nil, fmt.Errorf("endpoint index out of range: %d (max: %d)", endpointIdx, len(endpoints)-1)
	}

	ep := endpoints[endpointIdx]
	return &HTTPEndpoint{
		ServiceName:   ep.ServiceName,
		Route:         ep.Route,
		Method:        ep.Method,
		TargetService: ep.ServerAddress,
		Port:          ep.ServerPort,
	}, nil
}

func getAllNetworkPairs() ([]endpoint.CallPair, error) {
	sysData := registry.GetCurrent()
	if sysData == nil {
		return nil, fmt.Errorf("current system %s is not registered", systemconfig.GetCurrentSystem())
	}

	pairMap := make(map[string]*endpoint.CallPair)

	for _, service := range sysData.GetAllServices() {
		for _, ep := range sysData.GetHTTPEndpointsByService(service) {
			if ep.ServerAddress == "" || ep.ServerAddress == service {
				continue
			}
			key := ep.ServiceName + "->" + ep.ServerAddress
			if pairMap[key] == nil {
				pairMap[key] = &endpoint.CallPair{SourceService: ep.ServiceName, TargetService: ep.ServerAddress, SpanNames: []string{}, OperationTypes: []string{}}
			}
			if ep.SpanName != "" {
				pairMap[key].SpanNames = append(pairMap[key].SpanNames, ep.SpanName)
			}
			if !containsString(pairMap[key].OperationTypes, "http") {
				pairMap[key].OperationTypes = append(pairMap[key].OperationTypes, "http")
			}
		}
	}

	for _, service := range sysData.GetAllRPCServices() {
		for _, op := range sysData.GetRPCOperationsByService(service) {
			if op.ServerAddress == "" || op.ServerAddress == service {
				continue
			}
			key := op.ServiceName + "->" + op.ServerAddress
			if pairMap[key] == nil {
				pairMap[key] = &endpoint.CallPair{SourceService: op.ServiceName, TargetService: op.ServerAddress, SpanNames: []string{}, OperationTypes: []string{}}
			}
			if op.SpanName != "" {
				pairMap[key].SpanNames = append(pairMap[key].SpanNames, op.SpanName)
			}
			if !containsString(pairMap[key].OperationTypes, "rpc") {
				pairMap[key].OperationTypes = append(pairMap[key].OperationTypes, "rpc")
			}
		}
	}

	for _, service := range sysData.GetAllDatabaseServices() {
		for _, op := range sysData.GetDatabaseOperationsByService(service) {
			if op.ServerAddress == "" || op.ServerAddress == service {
				continue
			}
			key := op.ServiceName + "->" + op.ServerAddress
			if pairMap[key] == nil {
				pairMap[key] = &endpoint.CallPair{SourceService: op.ServiceName, TargetService: op.ServerAddress, SpanNames: []string{}, OperationTypes: []string{}}
			}
			if op.SpanName != "" {
				pairMap[key].SpanNames = append(pairMap[key].SpanNames, op.SpanName)
			}
			if !containsString(pairMap[key].OperationTypes, "db") {
				pairMap[key].OperationTypes = append(pairMap[key].OperationTypes, "db")
			}
		}
	}

	result := make([]endpoint.CallPair, 0, len(pairMap))
	for _, pair := range pairMap {
		pair.SpanNames = uniqueSorted(pair.SpanNames)
		sort.Strings(pair.OperationTypes)
		result = append(result, *pair)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].SourceService != result[j].SourceService {
			return result[i].SourceService < result[j].SourceService
		}
		return result[i].TargetService < result[j].TargetService
	})

	return result, nil
}

func getNetworkPairByIndex(networkPairIdx int) (*endpoint.CallPair, error) {
	networkPairs, err := getAllNetworkPairs()
	if err != nil {
		return nil, fmt.Errorf("failed to get network pairs: %w", err)
	}

	if networkPairIdx < 0 || networkPairIdx >= len(networkPairs) {
		return nil, fmt.Errorf("network pair index out of range: %d (max: %d)", networkPairIdx, len(networkPairs)-1)
	}

	return &networkPairs[networkPairIdx], nil
}

func getAllDNSEndpoints() ([]endpoint.DNSEndpointInfo, error) {
	sysData := registry.GetCurrent()
	if sysData == nil {
		return nil, fmt.Errorf("current system %s is not registered", systemconfig.GetCurrentSystem())
	}

	domainMap := make(map[string]*endpoint.DNSEndpointInfo)

	for _, service := range sysData.GetAllServices() {
		for _, ep := range sysData.GetHTTPEndpointsByService(service) {
			if ep.ServerAddress == "" || ep.ServerAddress == service {
				continue
			}
			key := ep.ServiceName + "->" + ep.ServerAddress
			if domainMap[key] == nil {
				domainMap[key] = &endpoint.DNSEndpointInfo{ServiceName: ep.ServiceName, Domain: ep.ServerAddress, SpanNames: []string{}}
			}
			domainMap[key].HasHTTP = true
			if ep.SpanName != "" {
				domainMap[key].SpanNames = append(domainMap[key].SpanNames, ep.SpanName)
			}
		}
	}

	for _, service := range sysData.GetAllDatabaseServices() {
		for _, op := range sysData.GetDatabaseOperationsByService(service) {
			if op.ServerAddress == "" || op.ServerAddress == service {
				continue
			}
			key := op.ServiceName + "->" + op.ServerAddress
			if domainMap[key] == nil {
				domainMap[key] = &endpoint.DNSEndpointInfo{ServiceName: op.ServiceName, Domain: op.ServerAddress, SpanNames: []string{}}
			}
			domainMap[key].HasDB = true
			if op.SpanName != "" {
				domainMap[key].SpanNames = append(domainMap[key].SpanNames, op.SpanName)
			}
		}
	}

	result := make([]endpoint.DNSEndpointInfo, 0, len(domainMap))
	for _, info := range domainMap {
		info.SpanNames = uniqueSorted(info.SpanNames)
		result = append(result, *info)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].ServiceName != result[j].ServiceName {
			return result[i].ServiceName < result[j].ServiceName
		}
		return result[i].Domain < result[j].Domain
	})

	return result, nil
}

func getDNSEndpointByIndex(dnsEndpointIdx int) (*endpoint.DNSEndpointInfo, error) {
	endpoints, err := getAllDNSEndpoints()
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS endpoints: %w", err)
	}

	if dnsEndpointIdx < 0 || dnsEndpointIdx >= len(endpoints) {
		return nil, fmt.Errorf("endpoint index out of range: %d (max: %d)", dnsEndpointIdx, len(endpoints)-1)
	}

	return &endpoints[dnsEndpointIdx], nil
}

func getAllDatabaseInfos() ([]endpoint.DatabaseInfo, error) {
	sysData := registry.GetCurrent()
	if sysData == nil {
		return nil, fmt.Errorf("current system %s is not registered", systemconfig.GetCurrentSystem())
	}

	result := make([]endpoint.DatabaseInfo, 0)
	for _, service := range sysData.GetAllDatabaseServices() {
		for _, op := range sysData.GetDatabaseOperationsByService(service) {
			result = append(result, endpoint.ToDatabaseInfo(op))
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].ServiceName != result[j].ServiceName {
			return result[i].ServiceName < result[j].ServiceName
		}
		if result[i].DBName != result[j].DBName {
			return result[i].DBName < result[j].DBName
		}
		return result[i].TableName < result[j].TableName
	})

	return result, nil
}

func getDatabaseInfoByIndex(databaseIdx int) (*endpoint.DatabaseInfo, error) {
	dbOps, err := getAllDatabaseInfos()
	if err != nil {
		return nil, fmt.Errorf("failed to get database operations: %w", err)
	}

	if databaseIdx < 0 || databaseIdx >= len(dbOps) {
		return nil, fmt.Errorf("database operation index out of range: %d (max: %d)", databaseIdx, len(dbOps)-1)
	}

	return &dbOps[databaseIdx], nil
}

func containsString(arr []string, target string) bool {
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}

func uniqueSorted(input []string) []string {
	if len(input) == 0 {
		return input
	}
	m := make(map[string]struct{}, len(input))
	for _, s := range input {
		m[s] = struct{}{}
	}
	result := make([]string, 0, len(m))
	for s := range m {
		result = append(result, s)
	}
	sort.Strings(result)
	return result
}
