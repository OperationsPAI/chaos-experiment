package handler

import (
	"context"
	"fmt"

	"github.com/LGU-SE-Internal/chaos-experiment/internal/resourcelookup"
	"github.com/LGU-SE-Internal/chaos-experiment/internal/systemconfig"
)

func getAllAppLabels(ctx context.Context, system systemconfig.SystemType, namespace string) ([]string, error) {
	labels, err := resourcelookup.GetSystemCache(system).GetAllAppLabels(ctx, namespace, defaultAppLabel)
	if err != nil {
		return nil, fmt.Errorf("failed to get app labels: %w", err)
	}
	return labels, nil
}

func getAppLabelByIndex(ctx context.Context, system systemconfig.SystemType, namespace string, appIdx int) (string, error) {
	labels, err := getAllAppLabels(ctx, system, namespace)
	if err != nil {
		return "", err
	}

	if appIdx < 0 || appIdx >= len(labels) {
		return "", fmt.Errorf("app index out of range: %d (max: %d)", appIdx, len(labels)-1)
	}

	return labels[appIdx], nil
}

func getAllContainerInfos(ctx context.Context, system systemconfig.SystemType, namespace string) ([]resourcelookup.ContainerInfo, error) {
	containers, err := resourcelookup.GetSystemCache(system).GetAllContainers(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get containers: %w", err)
	}
	return containers, nil
}

func getContainerInfoByIndex(ctx context.Context, system systemconfig.SystemType, namespace string, containerIdx int) (*resourcelookup.ContainerInfo, error) {
	containers, err := getAllContainerInfos(ctx, system, namespace)
	if err != nil {
		return nil, err
	}

	if containerIdx < 0 || containerIdx >= len(containers) {
		return nil, fmt.Errorf("container index out of range: %d (max: %d)", containerIdx, len(containers)-1)
	}

	return &containers[containerIdx], nil
}

func getAllHTTPEndpointInfos(system systemconfig.SystemType) ([]resourcelookup.AppEndpointPair, error) {
	endpoints, err := resourcelookup.GetSystemCache(system).GetAllHTTPEndpoints()
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP endpoints: %w", err)
	}
	return endpoints, nil
}

func getHTTPEndpointByIndex(system systemconfig.SystemType, endpointIdx int) (*resourcelookup.AppEndpointPair, error) {
	endpoints, err := getAllHTTPEndpointInfos(system)
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP endpoints: %w", err)
	}

	if endpointIdx < 0 || endpointIdx >= len(endpoints) {
		return nil, fmt.Errorf("endpoint index out of range: %d (max: %d)", endpointIdx, len(endpoints)-1)
	}

	return &endpoints[endpointIdx], nil
}

func getAllNetworkPairs(system systemconfig.SystemType) ([]resourcelookup.AppNetworkPair, error) {
	networkPairs, err := resourcelookup.GetSystemCache(system).GetAllNetworkPairs()
	if err != nil {
		return nil, fmt.Errorf("failed to get network pairs: %w", err)
	}
	return networkPairs, nil
}

func getNetworkPairByIndex(system systemconfig.SystemType, networkPairIdx int) (*resourcelookup.AppNetworkPair, error) {
	networkPairs, err := getAllNetworkPairs(system)
	if err != nil {
		return nil, fmt.Errorf("failed to get network pairs: %w", err)
	}

	if networkPairIdx < 0 || networkPairIdx >= len(networkPairs) {
		return nil, fmt.Errorf("network pair index out of range: %d (max: %d)", networkPairIdx, len(networkPairs)-1)
	}

	return &networkPairs[networkPairIdx], nil
}

func getAllDNSEndpoints(system systemconfig.SystemType) ([]resourcelookup.AppDNSPair, error) {
	endpoints, err := resourcelookup.GetSystemCache(system).GetAllDNSEndpoints()
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS endpoints: %w", err)
	}
	return endpoints, nil
}

func getDNSEndpointByIndex(system systemconfig.SystemType, dnsEndpointIdx int) (*resourcelookup.AppDNSPair, error) {
	endpoints, err := getAllDNSEndpoints(system)
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS endpoints: %w", err)
	}

	if dnsEndpointIdx < 0 || dnsEndpointIdx >= len(endpoints) {
		return nil, fmt.Errorf("endpoint index out of range: %d (max: %d)", dnsEndpointIdx, len(endpoints)-1)
	}

	return &endpoints[dnsEndpointIdx], nil
}

func getAllDatabaseOperations(system systemconfig.SystemType) ([]resourcelookup.AppDatabasePair, error) {
	operations, err := resourcelookup.GetSystemCache(system).GetAllDatabaseOperations()
	if err != nil {
		return nil, fmt.Errorf("failed to get database operations: %w", err)
	}
	return operations, nil
}

func getDatabaseOperationByIndex(system systemconfig.SystemType, databaseIdx int) (*resourcelookup.AppDatabasePair, error) {
	dbOps, err := getAllDatabaseOperations(system)
	if err != nil {
		return nil, fmt.Errorf("failed to get database operations: %w", err)
	}

	if databaseIdx < 0 || databaseIdx >= len(dbOps) {
		return nil, fmt.Errorf("database operation index out of range: %d (max: %d)", databaseIdx, len(dbOps)-1)
	}

	return &dbOps[databaseIdx], nil
}

func getAllJVMMethods(system systemconfig.SystemType) ([]resourcelookup.AppMethodPair, error) {
	methods, err := resourcelookup.GetSystemCache(system).GetAllJVMMethods()
	if err != nil {
		return nil, fmt.Errorf("failed to get JVM methods: %w", err)
	}
	return methods, nil
}

func getJVMMethodByIndex(system systemconfig.SystemType, methodIdx int) (*resourcelookup.AppMethodPair, error) {
	methods, err := getAllJVMMethods(system)
	if err != nil {
		return nil, err
	}

	if methodIdx < 0 || methodIdx >= len(methods) {
		return nil, fmt.Errorf("method index out of range: %d (max: %d)", methodIdx, len(methods)-1)
	}

	return &methods[methodIdx], nil
}

func getAllJVMRuntimeMutatorTargets(system systemconfig.SystemType) ([]resourcelookup.AppRuntimeMutatorTarget, error) {
	targets, err := resourcelookup.GetSystemCache(system).GetAllJVMRuntimeMutatorTargets()
	if err != nil {
		return nil, fmt.Errorf("failed to get JVM runtime mutator targets: %w", err)
	}
	return targets, nil
}

func getJVMRuntimeMutatorTargetByIndex(system systemconfig.SystemType, targetIdx int) (*resourcelookup.AppRuntimeMutatorTarget, error) {
	targets, err := getAllJVMRuntimeMutatorTargets(system)
	if err != nil {
		return nil, err
	}

	if targetIdx < 0 || targetIdx >= len(targets) {
		return nil, fmt.Errorf("target index out of range: %d (max: %d)", targetIdx, len(targets)-1)
	}

	return &targets[targetIdx], nil
}
