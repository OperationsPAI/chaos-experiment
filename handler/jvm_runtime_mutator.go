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

// MutationType constants for JVM runtime mutator
const (
	MutationTypeConstant = 0 // constant mutation
	MutationTypeOperator = 1 // operator mutation
	MutationTypeString   = 2 // string mutation
)

// MutationTypeNames maps mutation type index to name
var MutationTypeNames = []string{"constant", "operator", "string"}

// ConstantMutationConfigs defines predefined constant mutation configurations
// Each entry is a pair [from, to]
var ConstantMutationConfigs = [][2]string{
	{"true", "false"},              // 0: boolean true to false
	{"false", "true"},              // 1: boolean false to true
	{"0", "1"},                     // 2: zero to one
	{"1", "0"},                     // 3: one to zero
	{"100", "0"},                   // 4: hundred to zero
	{"\"success\"", "\"failure\""}, // 5: success to failure string
	{"\"ok\"", "\"error\""},        // 6: ok to error string
	{"-1", "0"},                    // 7: negative one to zero
	{"0", "-1"},                    // 8: zero to negative one
	{"1000", "1"},                  // 9: large value to small
	{"60", "1"},                    // 10: timeout value mutation
}

// JVMRuntimeMutatorSpec defines the JVM runtime mutator chaos injection parameters
type JVMRuntimeMutatorSpec struct {
	Duration     int `range:"1-60" description:"Time Unit Minute"`
	System       int `range:"0-0" dynamic:"true" description:"System Index"`
	MethodIdx    int `range:"0-0" dynamic:"true" description:"Flattened app+method index"`
	MutationType int `range:"0-2" description:"Mutation Type: 0=constant, 1=operator, 2=string"`
	MutationOpt  int `range:"0-10" description:"Mutation option: for constant(0-10 predefined pairs), operator(0-3 strategies), string(0-5 strategies)"`
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

	var optss []chaos.OptChaos

	// Get mutation type name from index
	if s.MutationType < 0 || s.MutationType >= len(MutationTypeNames) {
		return "", fmt.Errorf("mutation type index out of range: %d (max: %d)", s.MutationType, len(MutationTypeNames)-1)
	}
	mutationTypeName := MutationTypeNames[s.MutationType]

	// Configure mutation based on type
	switch s.MutationType {
	case MutationTypeConstant:
		from, to := s.getConstantConfig()
		optss = append(optss,
			chaos.WithRuntimeMutatorAction("constant"),
			chaos.WithRuntimeMutatorConfig(from, to),
		)
	case MutationTypeOperator:
		strategy := s.getMutationStrategy()
		optss = append(optss,
			chaos.WithRuntimeMutatorAction("operator"),
			chaos.WithRuntimeMutatorStrategy(strategy),
		)
	case MutationTypeString:
		strategy := s.getMutationStrategy()
		optss = append(optss,
			chaos.WithRuntimeMutatorAction("string"),
			chaos.WithRuntimeMutatorStrategy(strategy),
		)
	default:
		return "", fmt.Errorf("unsupported mutation type: %d", s.MutationType)
	}

	return controllers.CreateJVMRuntimeMutatorChaos(cli, ctx, ns, appName,
		className, methodName, mutationTypeName, duration, annotations, labels, optss...)
}

func (s *JVMRuntimeMutatorSpec) getMutationStrategy() string {
	strategies := map[int][]string{
		MutationTypeOperator: {
			"add_to_sub", "sub_to_add", "mul_to_div", "div_to_mul",
		},
		MutationTypeString: {
			"empty", "null", "reverse", "uppercase", "lowercase", "random",
		},
	}

	if strats, ok := strategies[s.MutationType]; ok {
		if s.MutationOpt >= 0 && s.MutationOpt < len(strats) {
			return strats[s.MutationOpt]
		}
	}

	// Default strategies
	if s.MutationType == MutationTypeOperator {
		return "add_to_sub"
	}
	return "empty"
}

// getConstantConfig returns the from/to pair for constant mutation based on MutationOpt
func (s *JVMRuntimeMutatorSpec) getConstantConfig() (string, string) {
	if s.MutationOpt >= 0 && s.MutationOpt < len(ConstantMutationConfigs) {
		config := ConstantMutationConfigs[s.MutationOpt]
		return config[0], config[1]
	}
	// Default: true to false
	return "true", "false"
}
