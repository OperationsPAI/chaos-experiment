package guidedcli

import (
	"context"
	"fmt"
	"sort"

	"github.com/OperationsPAI/chaos-experiment/internal/networkdependencies"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/OperationsPAI/chaos-experiment/internal/resourcelookup"
	"github.com/OperationsPAI/chaos-experiment/internal/serviceendpoints"
	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

func safeAppLabels(namespace string, systemType systemconfig.SystemType) ([]string, error) {
	pods, err := listPodsSafe(namespace)
	if err != nil {
		fallback := filteredServiceLabels(serviceendpoints.GetAllServices())
		sort.Strings(fallback)
		if len(fallback) > 0 {
			return fallback, nil
		}
		return nil, err
	}

	labelKey := systemconfig.GetAppLabelKey(systemType)
	seen := map[string]bool{}
	labels := make([]string, 0)
	for _, pod := range pods {
		if value := pod.Labels[labelKey]; value != "" && !seen[value] {
			seen[value] = true
			labels = append(labels, value)
		}
	}
	sort.Strings(labels)
	if len(labels) == 0 {
		fallback := filteredServiceLabels(serviceendpoints.GetAllServices())
		sort.Strings(fallback)
		return fallback, nil
	}
	filtered := filteredServiceLabels(labels)
	if len(filtered) == 0 {
		return labels, nil
	}
	return filtered, nil
}

func safeContainers(namespace string) ([]resourcelookup.ContainerInfo, error) {
	pods, err := listPodsSafe(namespace)
	if err != nil {
		return nil, err
	}
	labelKey := systemconfig.GetCurrentAppLabelKey()
	result := make([]resourcelookup.ContainerInfo, 0)
	for _, pod := range pods {
		appLabel := pod.Labels[labelKey]
		for _, container := range pod.Spec.Containers {
			result = append(result, resourcelookup.ContainerInfo{
				PodName:       pod.Name,
				AppLabel:      appLabel,
				ContainerName: container.Name,
			})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].AppLabel != result[j].AppLabel {
			return result[i].AppLabel < result[j].AppLabel
		}
		return result[i].ContainerName < result[j].ContainerName
	})
	return result, nil
}

func listPodsSafe(namespace string) ([]corev1.Pod, error) {
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
	list, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list pods in namespace %s: %w", namespace, err)
	}
	return list.Items, nil
}

func filteredServiceLabels(labels []string) []string {
	allowed := networkSourceServices()
	if len(allowed) == 0 {
		return labels
	}

	result := make([]string, 0, len(labels))
	seen := make(map[string]bool, len(labels))
	for _, label := range labels {
		if allowed[label] && !seen[label] {
			seen[label] = true
			result = append(result, label)
		}
	}
	return result
}

func networkSourceServices() map[string]bool {
	pairs := networkdependencies.GetAllServicePairs()
	if len(pairs) == 0 {
		return nil
	}

	allowed := make(map[string]bool, len(pairs))
	for _, pair := range pairs {
		if pair.SourceService != "" {
			allowed[pair.SourceService] = true
		}
	}
	return allowed
}
