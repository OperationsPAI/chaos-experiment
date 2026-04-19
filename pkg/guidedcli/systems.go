package guidedcli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

type systemInstance struct {
	Name        string
	SystemType  systemconfig.SystemType
	Namespace   string
	DisplayName string
}

func discoverSystemInstances() ([]systemInstance, []string) {
	namespaces, err := listNamespacesSafe()
	if err != nil {
		return fallbackSystemInstances(), []string{
			fmt.Sprintf("namespace discovery failed, fallback to default namespace instances: %v", err),
		}
	}

	instances := make([]systemInstance, 0)
	for _, namespace := range namespaces {
		if instance, ok := matchSystemInstance(namespace); ok {
			instances = append(instances, instance)
		}
	}

	sort.Slice(instances, func(i, j int) bool {
		return instances[i].Name < instances[j].Name
	})

	if len(instances) == 0 {
		return fallbackSystemInstances(), []string{"no matching namespaces found, fallback to default namespace instances"}
	}
	return instances, nil
}

func listNamespacesSafe() ([]string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = buildKubeconfigSafe()
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("create kubernetes clientset: %w", err)
	}

	list, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list namespaces: %w", err)
	}

	result := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, item.Name)
	}
	return result, nil
}

func buildKubeconfigSafe() (*rest.Config, error) {
	paths := make([]string, 0, 2)
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		paths = append(paths, filepath.Join(home, ".kube", "config"))
	}
	if userProfile := os.Getenv("USERPROFILE"); userProfile != "" {
		candidate := filepath.Join(userProfile, ".kube", "config")
		if len(paths) == 0 || paths[0] != candidate {
			paths = append(paths, candidate)
		}
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			config, err := clientcmd.BuildConfigFromFlags("", path)
			if err == nil {
				return config, nil
			}
		}
	}
	return nil, fmt.Errorf("kubeconfig not found in default locations")
}

func fallbackSystemInstances() []systemInstance {
	result := make([]systemInstance, 0, len(systemconfig.GetAllSystemTypes()))
	for _, system := range systemconfig.GetAllSystemTypes() {
		namespace, err := systemconfig.GetNamespaceByIndex(system, 0)
		if err != nil {
			continue
		}
		reg := systemconfig.GetRegistration(system)
		displayName := system.String()
		if reg != nil && reg.DisplayName != "" {
			displayName = reg.DisplayName
		}
		result = append(result, systemInstance{
			Name:        namespace,
			SystemType:  system,
			Namespace:   namespace,
			DisplayName: displayName,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func matchSystemInstance(namespace string) (systemInstance, bool) {
	for _, system := range systemconfig.GetAllSystemTypes() {
		reg := systemconfig.GetRegistration(system)
		if reg == nil {
			continue
		}
		pattern, err := regexp.Compile(reg.NsPattern)
		if err != nil {
			continue
		}
		if pattern.MatchString(namespace) {
			displayName := system.String()
			if reg.DisplayName != "" {
				displayName = reg.DisplayName
			}
			return systemInstance{
				Name:        namespace,
				SystemType:  system,
				Namespace:   namespace,
				DisplayName: displayName,
			}, true
		}
	}
	return systemInstance{}, false
}

func normalizeSystemSelection(cfg *GuidedConfig) error {
	if cfg.System != "" {
		instance, ok := matchSystemInstance(cfg.System)
		if ok {
			cfg.SystemType = instance.SystemType.String()
			if cfg.Namespace == "" {
				cfg.Namespace = instance.Namespace
			}
			return nil
		}

		system, ok := inferSystemType(cfg.System)
		if !ok {
			return fmt.Errorf("system %q does not match any registered namespace pattern or system name", cfg.System)
		}
		cfg.SystemType = system.String()
		if cfg.Namespace == "" {
			cfg.Namespace = cfg.System
		}
		return nil
	}

	if cfg.Namespace != "" {
		system, ok := inferSystemType(cfg.Namespace)
		if ok {
			cfg.System = cfg.Namespace
			cfg.SystemType = system.String()
			return nil
		}
	}

	if cfg.SystemType != "" && cfg.Namespace == "" {
		system, err := systemconfig.ParseSystemType(cfg.SystemType)
		if err != nil {
			return err
		}
		namespace, err := systemconfig.GetNamespaceByIndex(system, 0)
		if err != nil {
			return err
		}
		cfg.System = namespace
		cfg.Namespace = namespace
		return nil
	}
	return nil
}

func inferSystemType(value string) (systemconfig.SystemType, bool) {
	if instance, ok := matchSystemInstance(value); ok {
		return instance.SystemType, true
	}

	for _, system := range systemconfig.GetAllSystemTypes() {
		if value == system.String() {
			return system, true
		}

		reg := systemconfig.GetRegistration(system)
		if reg == nil {
			continue
		}
		if value == namespaceStem(reg.NsPattern) {
			return system, true
		}
	}
	return "", false
}

func namespaceStem(pattern string) string {
	stem := strings.TrimPrefix(pattern, "^")
	stem = strings.TrimSuffix(stem, "$")
	stem = strings.ReplaceAll(stem, `\d+`, "")
	stem = strings.ReplaceAll(stem, `\d*`, "")
	return stem
}

func systemTypeIndex(name string) (int, error) {
	for idx, system := range systemconfig.GetAllSystemTypes() {
		if system.String() == name {
			return idx, nil
		}
	}
	return 0, fmt.Errorf("system type %q not found", name)
}

func resolvedMap(cfg GuidedConfig) map[string]any {
	result := map[string]any{}
	addString := func(key, value string) {
		if strings.TrimSpace(value) != "" {
			result[key] = value
		}
	}

	addString("system", cfg.System)
	addString("system_type", cfg.SystemType)
	addString("namespace", cfg.Namespace)
	addString("app", cfg.App)
	addString("chaos_type", cfg.ChaosType)
	addString("container", cfg.Container)
	addString("target_service", cfg.TargetService)
	addString("domain", cfg.Domain)
	addString("class", cfg.Class)
	addString("method", cfg.Method)
	addString("mutator_config", cfg.MutatorConfig)
	addString("route", cfg.Route)
	addString("http_method", cfg.HTTPMethod)
	addString("database", cfg.Database)
	addString("table", cfg.Table)
	addString("operation", cfg.Operation)
	if cfg.Duration != nil {
		result["duration"] = *cfg.Duration
	}
	if cfg.MemorySize != nil {
		result["memory_size"] = *cfg.MemorySize
	}
	if cfg.MemWorker != nil {
		result["mem_worker"] = *cfg.MemWorker
	}
	if cfg.TimeOffset != nil {
		result["time_offset"] = *cfg.TimeOffset
	}
	if cfg.CPULoad != nil {
		result["cpu_load"] = *cfg.CPULoad
	}
	if cfg.CPUWorker != nil {
		result["cpu_worker"] = *cfg.CPUWorker
	}
	if cfg.Latency != nil {
		result["latency"] = *cfg.Latency
	}
	if cfg.Correlation != nil {
		result["correlation"] = *cfg.Correlation
	}
	if cfg.Jitter != nil {
		result["jitter"] = *cfg.Jitter
	}
	if cfg.Loss != nil {
		result["loss"] = *cfg.Loss
	}
	if cfg.Duplicate != nil {
		result["duplicate"] = *cfg.Duplicate
	}
	if cfg.Corrupt != nil {
		result["corrupt"] = *cfg.Corrupt
	}
	if cfg.Rate != nil {
		result["rate"] = *cfg.Rate
	}
	if cfg.Limit != nil {
		result["limit"] = *cfg.Limit
	}
	if cfg.Buffer != nil {
		result["buffer"] = *cfg.Buffer
	}
	addString("direction", cfg.Direction)
	if cfg.DelayDuration != nil {
		result["delay_duration"] = *cfg.DelayDuration
	}
	if cfg.LatencyDuration != nil {
		result["latency_duration"] = *cfg.LatencyDuration
	}
	if cfg.LatencyMs != nil {
		result["latency_ms"] = *cfg.LatencyMs
	}
	if cfg.CPUCount != nil {
		result["cpu_count"] = *cfg.CPUCount
	}
	addString("return_type", cfg.ReturnType)
	addString("return_value_opt", cfg.ReturnValueOpt)
	addString("exception_opt", cfg.ExceptionOpt)
	addString("mem_type", cfg.MemType)
	addString("body_type", cfg.BodyType)
	addString("replace_method", cfg.ReplaceMethod)
	if cfg.StatusCode != nil {
		result["status_code"] = *cfg.StatusCode
	}
	return result
}
