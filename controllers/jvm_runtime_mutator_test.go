package controllers

import (
	"context"
	"testing"

	"github.com/OperationsPAI/chaos-experiment/chaos"
	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// MockClient is a mock implementation of client.Client
type MockClient struct {
	mock.Mock
	client.Client
}

func (m *MockClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	args := m.Called(ctx, obj, opts)
	return args.Error(0)
}

func TestCreateJVMRuntimeMutatorChaos(t *testing.T) {
	tests := []struct {
		name         string
		mutationType string
		className    string
		methodName   string
		opts         []chaos.OptChaos
		wantErr      bool
		errorMsg     string
	}{
		{
			name:         "valid constant mutation",
			mutationType: "constant",
			className:    "com.example.TestClass",
			methodName:   "testMethod",
			opts: []chaos.OptChaos{
				chaos.WithRuntimeMutatorConfig("100", "200"),
			},
			wantErr: false,
		},
		{
			name:         "valid operator mutation",
			mutationType: "operator",
			className:    "com.example.TestClass",
			methodName:   "testMethod",
			opts: []chaos.OptChaos{
				chaos.WithRuntimeMutatorStrategy("add_to_sub"),
			},
			wantErr: false,
		},
		{
			name:         "valid string mutation",
			mutationType: "string",
			className:    "com.example.TestClass",
			methodName:   "testMethod",
			opts: []chaos.OptChaos{
				chaos.WithRuntimeMutatorStrategy("empty"),
			},
			wantErr: false,
		},
		{
			name:         "invalid mutation type",
			mutationType: "invalid",
			className:    "com.example.TestClass",
			methodName:   "testMethod",
			wantErr:      true,
			errorMsg:     "invalid mutation type: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fake client
			scheme := runtime.NewScheme()
			_ = v1alpha1.AddToScheme(scheme)
			fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

			ctx := context.Background()
			namespace := "test-namespace"
			appName := "test-app"
			duration := pointer.String("5m")
			annotations := map[string]string{"test": "annotation"}
			labels := map[string]string{"test": "label"}

			name, err := CreateJVMRuntimeMutatorChaos(
				fakeClient,
				ctx,
				namespace,
				appName,
				tt.className,
				tt.methodName,
				tt.mutationType,
				duration,
				annotations,
				labels,
				tt.opts...,
			)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, name)
				assert.Contains(t, name, "test-namespace-test-app-mutator")
			}
		})
	}
}

func TestAddJVMRuntimeMutatorWorkflowNodes(t *testing.T) {
	workflowSpec := &v1alpha1.WorkflowSpec{
		Templates: []v1alpha1.Template{},
	}

	namespace := "test-namespace"
	appList := []string{"app1", "app2"}
	mutationType := "constant"
	className := "com.example.TestClass"
	methodName := "testMethod"
	injectTime := pointer.String("5m")
	sleepTime := pointer.String("10m")

	result := AddJVMRuntimeMutatorWorkflowNodes(
		workflowSpec,
		namespace,
		appList,
		mutationType,
		className,
		methodName,
		injectTime,
		sleepTime,
		chaos.WithRuntimeMutatorConfig("100", "200"),
	)

	// Should add 4 templates (2 apps * (1 chaos + 1 sleep))
	assert.Len(t, result.Templates, 4)

	// Check first chaos template
	assert.Equal(t, v1alpha1.TypeRuntimeMutatorChaos, result.Templates[0].Type)
	assert.NotNil(t, result.Templates[0].EmbedChaos.RuntimeMutatorChaos)
	assert.Equal(t, injectTime, result.Templates[0].Deadline)

	// Check sleep template
	assert.Equal(t, v1alpha1.TypeSuspend, result.Templates[1].Type)
	assert.Equal(t, sleepTime, result.Templates[1].Deadline)
}

func TestAddJVMRuntimeMutatorWorkflowNodes_InvalidType(t *testing.T) {
	workflowSpec := &v1alpha1.WorkflowSpec{
		Templates: []v1alpha1.Template{},
	}

	namespace := "test-namespace"
	appList := []string{"app1"}
	mutationType := "invalid"
	className := "com.example.TestClass"
	methodName := "testMethod"
	injectTime := pointer.String("5m")
	sleepTime := pointer.String("10m")

	result := AddJVMRuntimeMutatorWorkflowNodes(
		workflowSpec,
		namespace,
		appList,
		mutationType,
		className,
		methodName,
		injectTime,
		sleepTime,
	)

	// Should return original spec unchanged for invalid type
	assert.Len(t, result.Templates, 0)
}

func TestScheduleJVMRuntimeMutator(t *testing.T) {
	// Create a mock client
	mockClient := &MockClient{}

	// Set up expectations
	mockClient.On("Create", mock.Anything, mock.AnythingOfType("*v1alpha1.Workflow"), mock.Anything).Return(nil)

	namespace := "test-namespace"
	appList := []string{"app1", "app2"}
	mutationType := "operator"
	className := "com.example.TestClass"
	methodName := "testMethod"

	// This should not panic or error with the mock
	ScheduleJVMRuntimeMutator(
		mockClient,
		namespace,
		appList,
		mutationType,
		className,
		methodName,
		chaos.WithRuntimeMutatorStrategy("add_to_sub"),
	)

	// Verify that Create was called
	mockClient.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("*v1alpha1.Workflow"), mock.Anything)
}

func TestScheduleJVMRuntimeMutator_InvalidType(t *testing.T) {
	// Create a mock client
	mockClient := &MockClient{}

	// No expectations set - Create should not be called for invalid type

	namespace := "test-namespace"
	appList := []string{"app1"}
	mutationType := "invalid"
	className := "com.example.TestClass"
	methodName := "testMethod"

	// This should return early without calling Create
	ScheduleJVMRuntimeMutator(
		mockClient,
		namespace,
		appList,
		mutationType,
		className,
		methodName,
	)

	// Verify that Create was NOT called
	mockClient.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}
