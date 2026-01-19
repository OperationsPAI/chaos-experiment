package handler

import (
	"context"
	"fmt"
	"strconv"

	chaos "github.com/LGU-SE-Internal/chaos-experiment/chaos"
	controllers "github.com/LGU-SE-Internal/chaos-experiment/controllers"
	"github.com/LGU-SE-Internal/chaos-experiment/internal/resourcelookup"
	"k8s.io/utils/pointer"
	cli "sigs.k8s.io/controller-runtime/pkg/client"
)

// JVMRuntimeMutatorSpec defines the JVM runtime mutator chaos injection parameters
type JVMRuntimeMutatorSpec struct {
	Duration       int    `range:"1-60" description:"Time Unit Minute"`
	System         int    `range:"0-0" dynamic:"true" description:"System Index"`
	MethodIdx      int    `range:"0-0" dynamic:"true" description:"Flattened app+method index"`
	MutationType   string `range:"constant,operator,string" description:"Mutation Type"`
	MutationOpt    int    `range:"0-10" description:"Mutation strategy option"`
	MutationFrom   string `description:"Mutation from value (for constant mutations)"`
	MutationTo     string `description:"Mutation to value (for constant mutations)"`
	MutationStrategy string `description:"Mutation strategy (for operator/string mutations)"`
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

	methods, err := resourcelookup.GetSystemCache(system).GetAllJVMMethods()
	if err != nil {
		return "", fmt.Errorf("failed to get JVM methods: %w", err)
	}

	if s.MethodIdx < 0 || s.MethodIdx >= len(methods) {
		return "", fmt.Errorf("method index out of range: %d (max: %d)", s.MethodIdx, len(methods)-1)
	}

	methodPair := methods[s.MethodIdx]
	appName := methodPair.AppName
	className := methodPair.ClassName
	methodName := methodPair.MethodName

	duration := pointer.String(strconv.Itoa(s.Duration) + "m")

	var optss []chaos.OptJVMChaos

	// Configure mutation based on type
	switch s.MutationType {
	case "constant":
		optss = append(optss,
			chaos.WithJVMRuntimeMutatorType("constant"),
			chaos.WithJVMMutationConfig(s.MutationFrom, s.MutationTo),
		)
	case "operator":
		strategy := s.getMutationStrategy()
		optss = append(optss,
			chaos.WithJVMRuntimeMutatorType("operator"),
			chaos.WithJVMMutationStrategy(strategy),
		)
	case "string":
		strategy := s.getMutationStrategy()
		optss = append(optss,
			chaos.WithJVMRuntimeMutatorType("string"),
			chaos.WithJVMMutationStrategy(strategy),
		)
	default:
		return "", fmt.Errorf("unsupported mutation type: %s", s.MutationType)
	}

	return controllers.CreateJVMRuntimeMutatorChaos(cli, ctx, ns, appName,
		className, methodName, s.MutationType, duration, annotations, labels, optss...)
}

func (s *JVMRuntimeMutatorSpec) getMutationStrategy() string {
	strategies := map[string][]string{
		"operator": {
			"add_to_sub", "sub_to_add", "mul_to_div", "div_to_mul",
		},
		"string": {
			"empty", "null", "reverse", "uppercase", "lowercase", "random",
		},
	}

	if strats, ok := strategies[s.MutationType]; ok {
		if s.MutationOpt >= 0 && s.MutationOpt < len(strats) {
			return strats[s.MutationOpt]
		}
	}

	// Default strategies
	if s.MutationType == "operator" {
		return "add_to_sub"
	}
	return "empty"
}
