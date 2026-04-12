package handler

import (
	"context"
	"fmt"
	"strconv"

	chaos "github.com/LGU-SE-Internal/chaos-experiment/chaos"
	controllers "github.com/LGU-SE-Internal/chaos-experiment/controllers"
	"k8s.io/utils/pointer"
	cli "sigs.k8s.io/controller-runtime/pkg/client"
)

// JVMRuntimeMutatorSpec defines the JVM runtime mutator chaos injection parameters
type JVMRuntimeMutatorSpec struct {
	Duration         int `range:"1-60" description:"Time Unit Minute"`
	System           int `range:"0-0" dynamic:"true" description:"System Index"`
	MutatorTargetIdx int `range:"0-0" dynamic:"true" description:"Flattened valid runtime mutator injection index"`
}

func (s *JVMRuntimeMutatorSpec) Create(cli cli.Client, opts ...Option) (string, error) {
	conf := Conf{}
	for _, opt := range opts {
		opt(&conf)
	}

	annotations := make(map[string]string)
	if conf.Annotations != nil {
		annotations = conf.Annotations
	}

	ctx := context.Background()
	if conf.Context != nil {
		ctx = conf.Context
	}

	labels := make(map[string]string)
	if conf.Labels != nil {
		labels = conf.Labels
	}

	ns := conf.Namespace
	system := conf.System

	target, err := getJVMRuntimeMutatorTargetByIndex(system, s.MutatorTargetIdx)
	if err != nil {
		return "", err
	}
	appName := target.AppName
	className := target.ClassName
	methodName := target.MethodName

	duration := pointer.String(strconv.Itoa(s.Duration) + "m")

	var optss []chaos.OptChaos

	switch target.MutationTypeName {
	case "constant":
		optss = append(optss,
			chaos.WithRuntimeMutatorAction("constant"),
			chaos.WithRuntimeMutatorConfig(target.MutationFrom, target.MutationTo),
		)
	case "operator":
		optss = append(optss,
			chaos.WithRuntimeMutatorAction("operator"),
			chaos.WithRuntimeMutatorStrategy(target.MutationStrategy),
		)
	case "string":
		optss = append(optss,
			chaos.WithRuntimeMutatorAction("string"),
			chaos.WithRuntimeMutatorStrategy(target.MutationStrategy),
		)
	default:
		return "", fmt.Errorf("unsupported mutation type: %s", target.MutationTypeName)
	}

	return controllers.CreateJVMRuntimeMutatorChaos(cli, ctx, ns, appName,
		className, methodName, target.MutationTypeName, duration, annotations, labels, optss...)
}
