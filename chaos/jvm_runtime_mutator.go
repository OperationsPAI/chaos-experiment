package chaos

import (
	"errors"

	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"

	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

// NewRuntimeMutatorChaos creates a new runtime mutator chaos spec
func NewRuntimeMutatorChaos(opts ...OptChaos) (*chaosmeshv1alpha1.RuntimeMutatorChaos, error) {
	config := ConfigChaos{}
	for _, opt := range opts {
		if opt != nil {
			opt(&config)
		}
	}

	if config.Name == "" {
		return nil, errors.New("the resource name is required")
	}
	if config.Namespace == "" {
		return nil, errors.New("the namespace is required")
	}
	if config.RuntimeMutatorChaos == nil {
		return nil, errors.New("runtimeMutatorChaos is required")
	}

	runtimeMutatorChaos := chaosmeshv1alpha1.RuntimeMutatorChaos{}
	runtimeMutatorChaos.Name = config.Name
	runtimeMutatorChaos.Namespace = config.Namespace
	config.RuntimeMutatorChaos.DeepCopyInto(&runtimeMutatorChaos.Spec)

	if config.Labels != nil {
		runtimeMutatorChaos.Labels = config.Labels
	}
	if config.Annotations != nil {
		runtimeMutatorChaos.Annotations = config.Annotations
	}

	return &runtimeMutatorChaos, nil
}

// OptRuntimeMutator options for runtime mutator chaos
type OptRuntimeMutator func(opt *chaosmeshv1alpha1.RuntimeMutatorChaosSpec)

// WithRuntimeMutatorAction sets the mutation action type (constant, operator, string)
func WithRuntimeMutatorAction(action chaosmeshv1alpha1.RuntimeMutatorChaosAction) OptChaos {
	return func(opt *ConfigChaos) {
		if opt.RuntimeMutatorChaos == nil {
			opt.RuntimeMutatorChaos = &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{}
		}
		opt.RuntimeMutatorChaos.Action = action
	}
}

// WithRuntimeMutatorClass sets the Java class to target
func WithRuntimeMutatorClass(class string) OptChaos {
	return func(opt *ConfigChaos) {
		if opt.RuntimeMutatorChaos == nil {
			opt.RuntimeMutatorChaos = &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{}
		}
		opt.RuntimeMutatorChaos.Class = class
	}
}

// WithRuntimeMutatorMethod sets the method to target
func WithRuntimeMutatorMethod(method string) OptChaos {
	return func(opt *ConfigChaos) {
		if opt.RuntimeMutatorChaos == nil {
			opt.RuntimeMutatorChaos = &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{}
		}
		opt.RuntimeMutatorChaos.Method = method
	}
}

// WithRuntimeMutatorConfig sets the mutation configuration (from/to for constant mutation)
func WithRuntimeMutatorConfig(from, to string) OptChaos {
	return func(opt *ConfigChaos) {
		if opt.RuntimeMutatorChaos == nil {
			opt.RuntimeMutatorChaos = &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{}
		}
		opt.RuntimeMutatorChaos.From = &from
		opt.RuntimeMutatorChaos.To = &to
	}
}

// WithRuntimeMutatorStrategy sets the mutation strategy (for operator/string mutations)
func WithRuntimeMutatorStrategy(strategy string) OptChaos {
	return func(opt *ConfigChaos) {
		if opt.RuntimeMutatorChaos == nil {
			opt.RuntimeMutatorChaos = &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{}
		}
		opt.RuntimeMutatorChaos.Strategy = &strategy
	}
}

// WithRuntimeMutatorPort sets the agent server port
func WithRuntimeMutatorPort(port int32) OptChaos {
	return func(opt *ConfigChaos) {
		if opt.RuntimeMutatorChaos == nil {
			opt.RuntimeMutatorChaos = &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{}
		}
		opt.RuntimeMutatorChaos.Port = port
	}
}

// WithRuntimeMutatorChaosSpec sets the RuntimeMutatorChaos spec
func WithRuntimeMutatorChaosSpec(spec *chaosmeshv1alpha1.RuntimeMutatorChaosSpec) OptChaos {
	return func(opt *ConfigChaos) {
		opt.RuntimeMutatorChaos = spec
	}
}

// GenerateRuntimeMutatorChaosSpec generates a RuntimeMutatorChaos spec with the given options
func GenerateRuntimeMutatorChaosSpec(namespace string, appName string, duration *string, opts ...OptChaos) *chaosmeshv1alpha1.RuntimeMutatorChaosSpec {
	spec := &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{
		ContainerSelector: chaosmeshv1alpha1.ContainerSelector{
			PodSelector: chaosmeshv1alpha1.PodSelector{
				Selector: chaosmeshv1alpha1.PodSelectorSpec{
					GenericSelectorSpec: chaosmeshv1alpha1.GenericSelectorSpec{
						Namespaces: []string{namespace},
						LabelSelectors: map[string]string{
							systemconfig.GetCurrentAppLabelKey(): appName,
						},
					},
				},
				Mode: chaosmeshv1alpha1.OneMode,
			},
		},
	}

	if duration != nil {
		spec.Duration = duration
	}

	// Apply options
	config := ConfigChaos{
		RuntimeMutatorChaos: spec,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&config)
		}
	}

	return spec
}
