package chaos

import (
	"errors"

	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
)

// NewJVMRuntimeMutatorChaos creates a new JVM runtime mutator chaos spec
func NewJVMRuntimeMutatorChaos(opts ...OptChaos) (*chaosmeshv1alpha1.JVMChaos, error) {
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
	if config.JVMChaos == nil {
		return nil, errors.New("jvmChaos is required")
	}

	jvmChaos := chaosmeshv1alpha1.JVMChaos{}
	jvmChaos.Name = config.Name
	jvmChaos.Namespace = config.Namespace
	config.JVMChaos.DeepCopyInto(&jvmChaos.Spec)

	if config.Labels != nil {
		jvmChaos.Labels = config.Labels
	}
	if config.Annotations != nil {
		jvmChaos.Annotations = config.Annotations
	}

	return &jvmChaos, nil
}

// OptJVMRuntimeMutator options for JVM runtime mutator chaos
type OptJVMRuntimeMutator func(opt *chaosmeshv1alpha1.JVMChaosSpec)

// WithJVMRuntimeMutatorType sets the mutation type (constant, operator, string)
func WithJVMRuntimeMutatorType(mutationType string) OptJVMChaos {
	return func(opt *chaosmeshv1alpha1.JVMChaosSpec) {
		// Store mutation type in JVMParameter for runtime mutator
		if opt.JVMParameter.JVMClassMethodSpec.Class == "" {
			opt.JVMParameter.JVMClassMethodSpec.Class = mutationType
		}
	}
}

// WithJVMMutationConfig sets the mutation configuration
func WithJVMMutationConfig(from, to string) OptJVMChaos {
	return func(opt *chaosmeshv1alpha1.JVMChaosSpec) {
		// Use ReturnValue field to store mutation config as JSON
		opt.JVMParameter.ReturnValue = from + ":" + to
	}
}

// WithJVMMutationStrategy sets the mutation strategy
func WithJVMMutationStrategy(strategy string) OptJVMChaos {
	return func(opt *chaosmeshv1alpha1.JVMChaosSpec) {
		opt.JVMParameter.ReturnValue = strategy
	}
}
