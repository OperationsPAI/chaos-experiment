package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/LGU-SE-Internal/chaos-experiment/chaos"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/rand"
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
	opts ...chaos.OptJVMChaos,
) (string, error) {

	spec := chaos.GenerateJVMChaosSpec(namespace, appName, duration, append([]chaos.OptJVMChaos{
		chaos.WithJVMClass(className),
		chaos.WithJVMMethod(methodName),
	}, opts...)...)

	name := strings.ToLower(fmt.Sprintf("%s-%s-mutator-%s-%s", namespace, appName, mutationType, rand.String(6)))

	jvmChaos, err := chaos.NewJvmChaos(
		chaos.WithAnnotations(annotations),
		chaos.WithLabels(labels),
		chaos.WithName(name),
		chaos.WithNamespace(namespace),
		chaos.WithJVMChaosSpec(spec),
	)

	if err != nil {
		logrus.Errorf("Failed to create chaos: %v", err)
		return "", err
	}

	create, err := jvmChaos.ValidateCreate()
	if err != nil {
		logrus.Errorf("Failed to validate create chaos: %v", err)
		return "", err
	}
	logrus.Infof("create warning: %v", create)

	err = cli.Create(ctx, jvmChaos)
	if err != nil {
		logrus.Errorf("Failed to create chaos: %v", err)
		return "", err
	}

	return name, nil
}
