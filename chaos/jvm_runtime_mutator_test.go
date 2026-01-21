package chaos

import (
	"testing"

	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/pointer"
)

func TestNewRuntimeMutatorChaos(t *testing.T) {
	tests := []struct {
		name    string
		opts    []OptChaos
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid constant mutation",
			opts: []OptChaos{
				WithName("test-mutation"),
				WithNamespace("test-ns"),
				WithRuntimeMutatorAction("constant"),
				WithRuntimeMutatorClass("com.example.TestClass"),
				WithRuntimeMutatorMethod("testMethod"),
				WithRuntimeMutatorConfig("100", "0"),
				WithRuntimeMutatorPort(9090),
			},
			wantErr: false,
		},
		{
			name: "valid operator mutation",
			opts: []OptChaos{
				WithName("test-operator"),
				WithNamespace("test-ns"),
				WithRuntimeMutatorAction("operator"),
				WithRuntimeMutatorClass("com.example.MathClass"),
				WithRuntimeMutatorMethod("calculate"),
				WithRuntimeMutatorStrategy("add-to-sub"),
				WithRuntimeMutatorPort(9091),
			},
			wantErr: false,
		},
		{
			name: "valid string mutation",
			opts: []OptChaos{
				WithName("test-string"),
				WithNamespace("test-ns"),
				WithRuntimeMutatorAction("string"),
				WithRuntimeMutatorClass("com.example.StringClass"),
				WithRuntimeMutatorMethod("getMessage"),
				WithRuntimeMutatorStrategy("empty-string"),
				WithRuntimeMutatorPort(9092),
			},
			wantErr: false,
		},
		{
			name: "missing name",
			opts: []OptChaos{
				WithNamespace("test-ns"),
				WithRuntimeMutatorAction("constant"),
				WithRuntimeMutatorClass("com.example.TestClass"),
				WithRuntimeMutatorMethod("testMethod"),
				WithRuntimeMutatorConfig("100", "0"),
			},
			wantErr: true,
			errMsg:  "the resource name is required",
		},
		{
			name: "missing namespace",
			opts: []OptChaos{
				WithName("test-mutation"),
				WithRuntimeMutatorAction("constant"),
				WithRuntimeMutatorClass("com.example.TestClass"),
				WithRuntimeMutatorMethod("testMethod"),
				WithRuntimeMutatorConfig("100", "0"),
			},
			wantErr: true,
			errMsg:  "the namespace is required",
		},
		{
			name: "missing RuntimeMutatorChaos spec",
			opts: []OptChaos{
				WithName("test-mutation"),
				WithNamespace("test-ns"),
			},
			wantErr: true,
			errMsg:  "runtimeMutatorChaos is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chaos, err := NewRuntimeMutatorChaos(tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, chaos)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, chaos)
				assert.Equal(t, "test-ns", chaos.Namespace)
			}
		})
	}
}

func TestRuntimeMutatorOptions(t *testing.T) {
	t.Run("WithRuntimeMutatorAction", func(t *testing.T) {
		config := &ConfigChaos{RuntimeMutatorChaos: &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{}}
		opt := WithRuntimeMutatorAction(chaosmeshv1alpha1.RuntimeMutatorOperatorAction)
		opt(config)
		assert.Equal(t, chaosmeshv1alpha1.RuntimeMutatorOperatorAction, config.RuntimeMutatorChaos.Action)
	})

	t.Run("WithRuntimeMutatorClass", func(t *testing.T) {
		config := &ConfigChaos{RuntimeMutatorChaos: &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{}}
		opt := WithRuntimeMutatorClass("com.example.TestClass")
		opt(config)
		assert.Equal(t, "com.example.TestClass", config.RuntimeMutatorChaos.Class)
	})

	t.Run("WithRuntimeMutatorMethod", func(t *testing.T) {
		config := &ConfigChaos{RuntimeMutatorChaos: &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{}}
		opt := WithRuntimeMutatorMethod("testMethod")
		opt(config)
		assert.Equal(t, "testMethod", config.RuntimeMutatorChaos.Method)
	})

	t.Run("WithRuntimeMutatorConfig", func(t *testing.T) {
		config := &ConfigChaos{RuntimeMutatorChaos: &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{}}
		opt := WithRuntimeMutatorConfig("100", "0")
		opt(config)
		assert.Equal(t, pointer.String("100"), config.RuntimeMutatorChaos.From)
		assert.Equal(t, pointer.String("0"), config.RuntimeMutatorChaos.To)
	})

	t.Run("WithRuntimeMutatorStrategy", func(t *testing.T) {
		config := &ConfigChaos{RuntimeMutatorChaos: &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{}}
		opt := WithRuntimeMutatorStrategy("add-to-sub")
		opt(config)
		assert.Equal(t, pointer.String("add-to-sub"), config.RuntimeMutatorChaos.Strategy)
	})

	t.Run("WithRuntimeMutatorPort", func(t *testing.T) {
		config := &ConfigChaos{RuntimeMutatorChaos: &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{}}
		opt := WithRuntimeMutatorPort(9090)
		opt(config)
		assert.Equal(t, int32(9090), config.RuntimeMutatorChaos.Port)
	})

	t.Run("WithRuntimeMutatorChaosSpec", func(t *testing.T) {
		spec := &chaosmeshv1alpha1.RuntimeMutatorChaosSpec{
			Action: chaosmeshv1alpha1.RuntimeMutatorStringAction,
			RuntimeMutatorParameter: chaosmeshv1alpha1.RuntimeMutatorParameter{
				Class:    "com.example.StringClass",
				Method:   "getString",
				Strategy: pointer.String("empty-string"),
			},
		}
		config := &ConfigChaos{}
		opt := WithRuntimeMutatorChaosSpec(spec)
		opt(config)
		assert.Equal(t, spec, config.RuntimeMutatorChaos)
	})
}

func TestGenerateRuntimeMutatorChaosSpec(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		appName   string
		duration  *string
		opts      []OptChaos
		validate  func(t *testing.T, spec *chaosmeshv1alpha1.RuntimeMutatorChaosSpec)
	}{
		{
			name:      "basic spec with duration",
			namespace: "test-ns",
			appName:   "test-app",
			duration:  pointer.String("5m"),
			opts:      []OptChaos{},
			validate: func(t *testing.T, spec *chaosmeshv1alpha1.RuntimeMutatorChaosSpec) {
				assert.Equal(t, "test-ns", spec.Selector.Namespaces[0])
				assert.Equal(t, "test-app", spec.Selector.LabelSelectors["app"])
				assert.Equal(t, pointer.String("5m"), spec.Duration)
			},
		},
		{
			name:      "spec with mutation configuration",
			namespace: "ts",
			appName:   "ts-user-service",
			duration:  pointer.String("2m"),
			opts: []OptChaos{
				WithRuntimeMutatorAction("constant"),
				WithRuntimeMutatorClass("com.example.UserService"),
				WithRuntimeMutatorMethod("getUserById"),
				WithRuntimeMutatorConfig("100", "0"),
				WithRuntimeMutatorPort(9090),
			},
			validate: func(t *testing.T, spec *chaosmeshv1alpha1.RuntimeMutatorChaosSpec) {
				assert.Equal(t, chaosmeshv1alpha1.RuntimeMutatorConstantAction, spec.Action)
				assert.Equal(t, "com.example.UserService", spec.RuntimeMutatorParameter.Class)
				assert.Equal(t, "getUserById", spec.RuntimeMutatorParameter.Method)
				assert.Equal(t, pointer.String("100"), spec.From)
				assert.Equal(t, pointer.String("0"), spec.To)
				assert.Equal(t, int32(9090), spec.Port)
			},
		},
		{
			name:      "spec with operator mutation",
			namespace: "ts",
			appName:   "ts-order-service",
			duration:  pointer.String("3m"),
			opts: []OptChaos{
				WithRuntimeMutatorAction("operator"),
				WithRuntimeMutatorClass("com.example.OrderCalculator"),
				WithRuntimeMutatorMethod("calculateTotal"),
				WithRuntimeMutatorStrategy("add-to-sub"),
				WithRuntimeMutatorPort(9091),
			},
			validate: func(t *testing.T, spec *chaosmeshv1alpha1.RuntimeMutatorChaosSpec) {
				assert.Equal(t, chaosmeshv1alpha1.RuntimeMutatorOperatorAction, spec.Action)
				assert.Equal(t, "com.example.OrderCalculator", spec.RuntimeMutatorParameter.Class)
				assert.Equal(t, "calculateTotal", spec.RuntimeMutatorParameter.Method)
				assert.Equal(t, pointer.String("add-to-sub"), spec.Strategy)
				assert.Equal(t, int32(9091), spec.Port)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := GenerateRuntimeMutatorChaosSpec(tt.namespace, tt.appName, tt.duration, tt.opts...)
			assert.NotNil(t, spec)
			if tt.validate != nil {
				tt.validate(t, spec)
			}
		})
	}
}

func TestRuntimeMutatorChaosWithLabelsAndAnnotations(t *testing.T) {
	labels := map[string]string{
		"env":     "test",
		"version": "v1.0",
	}
	annotations := map[string]string{
		"description": "Test runtime mutation",
		"owner":       "test-team",
	}

	chaos, err := NewRuntimeMutatorChaos(
		WithName("labeled-mutation"),
		WithNamespace("test-ns"),
		WithRuntimeMutatorAction(chaosmeshv1alpha1.RuntimeMutatorConstantAction),
		WithRuntimeMutatorClass("com.example.TestClass"),
		WithRuntimeMutatorMethod("testMethod"),
		WithRuntimeMutatorConfig("true", "false"),
		WithLabels(labels),
		WithAnnotations(annotations),
	)

	assert.NoError(t, err)
	assert.NotNil(t, chaos)
	assert.Equal(t, labels, chaos.Labels)
	assert.Equal(t, annotations, chaos.Annotations)
}
