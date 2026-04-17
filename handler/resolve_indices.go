package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

func resolveSystemType(systemName string) (systemconfig.SystemType, error) {
	if systemName == "" {
		return systemconfig.SystemTrainTicket, nil
	}

	if system, err := systemconfig.ParseSystemType(systemName); err == nil {
		return system, nil
	}

	for _, system := range systemconfig.GetAllSystemTypes() {
		name := system.String()
		if strings.HasPrefix(name, systemName) || strings.HasPrefix(systemName, name) {
			return system, nil
		}
	}

	return "", fmt.Errorf("unknown system %q", systemName)
}

func ResolveSystemIndex(systemName string) (int, error) {
	system, err := resolveSystemType(systemName)
	if err != nil {
		return 0, err
	}

	for idx, candidate := range systemconfig.GetAllSystemTypes() {
		if candidate == system {
			return idx, nil
		}
	}

	return 0, fmt.Errorf("system %q is not registered", systemName)
}

func ResolveAppIndex(ctx context.Context, systemName, namespace, app string) (int, error) {
	system, err := resolveSystemType(systemName)
	if err != nil {
		return 0, err
	}

	labels, err := getAllAppLabels(ctx, system, namespace)
	if err != nil {
		return 0, err
	}
	for idx, label := range labels {
		if label == app {
			return idx, nil
		}
	}

	return 0, fmt.Errorf("app %q not found in namespace %q", app, namespace)
}

func ResolveContainerIndex(ctx context.Context, systemName, namespace, app, container string) (int, error) {
	system, err := resolveSystemType(systemName)
	if err != nil {
		return 0, err
	}

	containers, err := getAllContainerInfos(ctx, system, namespace)
	if err != nil {
		return 0, err
	}
	for idx, info := range containers {
		if info.AppLabel != app {
			continue
		}
		if container == "" || info.ContainerName == container {
			return idx, nil
		}
	}

	if container == "" {
		return 0, fmt.Errorf("no container found for app %q in namespace %q", app, namespace)
	}
	return 0, fmt.Errorf("container %q for app %q not found in namespace %q", container, app, namespace)
}

func ResolveNetworkPairIndex(systemName, sourceService, targetService string) (int, error) {
	system, err := resolveSystemType(systemName)
	if err != nil {
		return 0, err
	}

	pairs, err := getAllNetworkPairs(system)
	if err != nil {
		return 0, err
	}
	for idx, pair := range pairs {
		if pair.SourceService == sourceService && pair.TargetService == targetService {
			return idx, nil
		}
	}

	return 0, fmt.Errorf("network pair %q -> %q not found", sourceService, targetService)
}

func ResolveHTTPEndpointIndex(systemName, app, route, method string, port int) (int, error) {
	system, err := resolveSystemType(systemName)
	if err != nil {
		return 0, err
	}

	endpoints, err := getAllHTTPEndpointInfos(system)
	if err != nil {
		return 0, err
	}

	wantMethod := strings.ToUpper(method)
	wantPort := strconv.Itoa(port)
	for idx, endpoint := range endpoints {
		if endpoint.AppName != app {
			continue
		}
		if route != "" && endpoint.Route != route {
			continue
		}
		if wantMethod != "" && strings.ToUpper(endpoint.Method) != wantMethod {
			continue
		}
		if port > 0 && endpoint.ServerPort != wantPort {
			continue
		}
		return idx, nil
	}

	return 0, fmt.Errorf("http endpoint for app %q route %q method %q port %d not found", app, route, method, port)
}

func ResolveJVMMethodIndex(systemName, app, className, methodName string) (int, error) {
	system, err := resolveSystemType(systemName)
	if err != nil {
		return 0, err
	}

	methods, err := getAllJVMMethods(system)
	if err != nil {
		return 0, err
	}
	for idx, method := range methods {
		if method.AppName != app {
			continue
		}
		if className != "" && method.ClassName != className {
			continue
		}
		if method.MethodName != methodName {
			continue
		}
		return idx, nil
	}

	return 0, fmt.Errorf("jvm method %q.%q for app %q not found", className, methodName, app)
}
