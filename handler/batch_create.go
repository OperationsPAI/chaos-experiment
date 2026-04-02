package handler

import (
	"context"
	"fmt"

	"github.com/LGU-SE-Internal/chaos-experiment/internal/systemconfig"
)

// BatchCreate creates multiple fault injections and returns their CRD names.
// The system and namespace parameters identify where to inject; namespaceTargetIndex
// defaults to 0 if it cannot be derived from the namespace.
func BatchCreate(ctx context.Context, confs []InjectionConf, system systemconfig.SystemType, namespace string, annotations map[string]string, labels map[string]string) ([]string, error) {
	if len(confs) == 0 {
		return nil, fmt.Errorf("no injection configurations provided")
	}

	names := make([]string, 0, len(confs))
	for _, conf := range confs {
		name, err := conf.Create(ctx, 0, annotations, labels)
		if err != nil {
			return names, fmt.Errorf("failed to create injection: %w", err)
		}
		names = append(names, name)
	}
	return names, nil
}
