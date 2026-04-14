package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/OperationsPAI/chaos-experiment/chaos"
	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateJVMRuntimeMutatorChaos creates a JVM runtime mutator chaos experiment
func CreateJVMRuntimeMutatorChaos(
	cli client.Client,
	ctx context.Context,
	namespace string,
	appName string,
	className string,
	methodName string,
	mutationType string,
	duration *string,
	annotations map[string]string,
	labels map[string]string,
	opts ...chaos.OptChaos,
) (string, error) {

	// Parse mutation type
	var action chaosmeshv1alpha1.RuntimeMutatorChaosAction
	switch mutationType {
	case "constant":
		action = chaosmeshv1alpha1.RuntimeMutatorConstantAction
	case "operator":
		action = chaosmeshv1alpha1.RuntimeMutatorOperatorAction
	case "string":
		action = chaosmeshv1alpha1.RuntimeMutatorStringAction
	default:
		return "", fmt.Errorf("invalid mutation type: %s", mutationType)
	}

	spec := chaos.GenerateRuntimeMutatorChaosSpec(namespace, appName, duration, append([]chaos.OptChaos{
		chaos.WithRuntimeMutatorAction(action),
		chaos.WithRuntimeMutatorClass(className),
		chaos.WithRuntimeMutatorMethod(methodName),
	}, opts...)...)

	name := strings.ToLower(fmt.Sprintf("%s-%s-mutator-%s-%s", namespace, appName, mutationType, rand.String(6)))

	runtimeMutatorChaos, err := chaos.NewRuntimeMutatorChaos(
		chaos.WithAnnotations(annotations),
		chaos.WithLabels(labels),
		chaos.WithName(name),
		chaos.WithNamespace(namespace),
		chaos.WithRuntimeMutatorChaosSpec(spec),
	)

	if err != nil {
		logrus.Errorf("Failed to create chaos: %v", err)
		return "", err
	}

	create, err := runtimeMutatorChaos.ValidateCreate(ctx, runtimeMutatorChaos)
	if err != nil {
		logrus.Errorf("Failed to validate create chaos: %v", err)
		return "", err
	}
	logrus.Infof("create warning: %v", create)

	err = cli.Create(ctx, runtimeMutatorChaos)
	if err != nil {
		logrus.Errorf("Failed to create chaos: %v", err)
		return "", err
	}

	return name, nil
}

// AddJVMRuntimeMutatorWorkflowNodes adds JVM runtime mutator chaos nodes to a workflow
func AddJVMRuntimeMutatorWorkflowNodes(
	workflowSpec *chaosmeshv1alpha1.WorkflowSpec,
	namespace string,
	appList []string,
	mutationType string,
	className string,
	methodName string,
	injectTime *string,
	sleepTime *string,
	opts ...chaos.OptChaos,
) *chaosmeshv1alpha1.WorkflowSpec {
	// Parse mutation type
	var action chaosmeshv1alpha1.RuntimeMutatorChaosAction
	switch mutationType {
	case "constant":
		action = chaosmeshv1alpha1.RuntimeMutatorConstantAction
	case "operator":
		action = chaosmeshv1alpha1.RuntimeMutatorOperatorAction
	case "string":
		action = chaosmeshv1alpha1.RuntimeMutatorStringAction
	default:
		logrus.Errorf("invalid mutation type: %s", mutationType)
		return workflowSpec
	}

	for _, appName := range appList {
		spec := chaos.GenerateRuntimeMutatorChaosSpec(namespace, appName, nil, append([]chaos.OptChaos{
			chaos.WithRuntimeMutatorAction(action),
			chaos.WithRuntimeMutatorClass(className),
			chaos.WithRuntimeMutatorMethod(methodName),
		}, opts...)...)

		workflowSpec.Templates = append(workflowSpec.Templates, chaosmeshv1alpha1.Template{
			Name: strings.ToLower(fmt.Sprintf("%s-%s-mutator-%s-%s", namespace, appName, mutationType, rand.String(6))),
			Type: chaosmeshv1alpha1.TypeRuntimeMutatorChaos,
			EmbedChaos: &chaosmeshv1alpha1.EmbedChaos{
				RuntimeMutatorChaos: spec,
			},
			Deadline: injectTime,
		})

		workflowSpec.Templates = append(workflowSpec.Templates, chaosmeshv1alpha1.Template{
			Name:     fmt.Sprintf("%s-%s", "sleep", rand.String(6)),
			Type:     chaosmeshv1alpha1.TypeSuspend,
			Deadline: sleepTime,
		})
	}

	return workflowSpec
}

// ScheduleJVMRuntimeMutator creates and schedules a JVM runtime mutator chaos workflow
func ScheduleJVMRuntimeMutator(
	cli client.Client,
	namespace string,
	appList []string,
	mutationType string,
	className string,
	methodName string,
	opts ...chaos.OptChaos,
) {
	// Parse mutation type
	var action chaosmeshv1alpha1.RuntimeMutatorChaosAction
	switch mutationType {
	case "constant":
		action = chaosmeshv1alpha1.RuntimeMutatorConstantAction
	case "operator":
		action = chaosmeshv1alpha1.RuntimeMutatorOperatorAction
	case "string":
		action = chaosmeshv1alpha1.RuntimeMutatorStringAction
	default:
		logrus.Errorf("invalid mutation type: %s", mutationType)
		return
	}

	workflowName := strings.ToLower(fmt.Sprintf("%s-%s-%s", namespace, mutationType, rand.String(6)))
	workflowSpec := chaosmeshv1alpha1.WorkflowSpec{
		Entry: workflowName,
		Templates: []chaosmeshv1alpha1.Template{
			{
				Name:     workflowName,
				Type:     chaosmeshv1alpha1.TypeSerial,
				Children: nil,
			},
		},
	}

	for idx, appName := range appList {
		spec := chaos.GenerateRuntimeMutatorChaosSpec(namespace, appName, nil, append([]chaos.OptChaos{
			chaos.WithRuntimeMutatorAction(action),
			chaos.WithRuntimeMutatorClass(className),
			chaos.WithRuntimeMutatorMethod(methodName),
		}, opts...)...)

		workflowSpec.Templates = append(workflowSpec.Templates, chaosmeshv1alpha1.Template{
			Name: strings.ToLower(fmt.Sprintf("%s-%s-%s", namespace, appName, mutationType)),
			Type: chaosmeshv1alpha1.TypeRuntimeMutatorChaos,
			EmbedChaos: &chaosmeshv1alpha1.EmbedChaos{
				RuntimeMutatorChaos: spec,
			},
			Deadline: pointer.String("5m"),
		})

		if idx < len(appList)-1 {
			workflowSpec.Templates = append(workflowSpec.Templates, chaosmeshv1alpha1.Template{
				Name:     fmt.Sprintf("%s-%d", "sleep", idx),
				Type:     chaosmeshv1alpha1.TypeSuspend,
				Deadline: pointer.String("10m"),
			})
		}
	}

	for i, template := range workflowSpec.Templates {
		if i == 0 {
			continue
		}
		workflowSpec.Templates[0].Children = append(workflowSpec.Templates[0].Children, template.Name)
	}

	workflowChaos, err := chaos.NewWorkflowChaos(chaos.WithName(workflowName), chaos.WithNamespace(namespace), chaos.WithWorkflowSpec(&workflowSpec))
	if err != nil {
		logrus.Errorf("Failed to create chaos workflow: %v", err)
		return
	}

	create, err := workflowChaos.ValidateCreate(context.Background(), workflowChaos)
	if err != nil {
		logrus.Errorf("Failed to validate create chaos: %v", err)
		return
	}
	logrus.Infof("create warning: %v", create)

	err = cli.Create(context.Background(), workflowChaos)
	if err != nil {
		logrus.Errorf("Failed to create chaos: %v", err)
	}
}
