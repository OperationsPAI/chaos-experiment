package chaoscli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	chaospkg "github.com/OperationsPAI/chaos-experiment/chaos"
	chaosclient "github.com/OperationsPAI/chaos-experiment/client"
	"github.com/OperationsPAI/chaos-experiment/controllers"
	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type DirectSubmitter struct {
	client ctrlclient.Client
	writer io.Writer
}

func NewDirectSubmitter(k8sClient ctrlclient.Client, w io.Writer) *DirectSubmitter {
	if w == nil {
		w = os.Stdout
	}
	return &DirectSubmitter{client: k8sClient, writer: w}
}

func (d *DirectSubmitter) Submit(ctx context.Context, spec Spec) error {
	if d.client == nil {
		d.client = chaosclient.GetK8sClient()
	}
	switch spec.Type {
	case "NetworkDelay":
		targetService, _ := stringParam(spec.Params, "target_service")
		if targetService == "" {
			return fmt.Errorf("--target-service is required")
		}
		direction, err := networkDirection(spec.Params["direction"])
		if err != nil {
			return err
		}
		_, err = controllers.CreateNetworkDelayChaos(
			d.client,
			ctx,
			spec.Namespace,
			spec.Target,
			fmt.Sprintf("%dms", intParam(spec.Params, "latency", 100)),
			strconv.Itoa(intParam(spec.Params, "correlation", 0)),
			fmt.Sprintf("%dms", intParam(spec.Params, "jitter", 0)),
			&spec.Duration,
			nil,
			nil,
			chaospkg.WithNetworkTargetAndDirection(spec.Namespace, targetService, direction),
		)
		return err
	case "NetworkLoss":
		targetService, _ := stringParam(spec.Params, "target_service")
		if targetService == "" {
			return fmt.Errorf("--target-service is required")
		}
		direction, err := networkDirection(spec.Params["direction"])
		if err != nil {
			return err
		}
		_, err = controllers.CreateNetworkLossChaos(
			d.client,
			ctx,
			spec.Namespace,
			spec.Target,
			strconv.Itoa(intParam(spec.Params, "loss", 10)),
			strconv.Itoa(intParam(spec.Params, "correlation", 0)),
			&spec.Duration,
			nil,
			nil,
			chaospkg.WithNetworkTargetAndDirection(spec.Namespace, targetService, direction),
		)
		return err
	case "HTTPRequestAbort":
		route, _ := stringParam(spec.Params, "route")
		method, _ := stringParam(spec.Params, "http_method")
		port := intParam(spec.Params, "port", 80)
		abort := true
		_, err := controllers.CreateHTTPChaos(
			d.client,
			ctx,
			spec.Namespace,
			spec.Target,
			"request-abort",
			&spec.Duration,
			nil,
			nil,
			chaospkg.WithTarget(chaosmeshv1alpha1.PodHttpRequest),
			chaospkg.WithAbort(&abort),
			chaospkg.WithPort(int32(port)),
			chaospkg.WithPath(&route),
			chaospkg.WithMethod(&method),
		)
		return err
	case "JVMLatency":
		className, _ := stringParam(spec.Params, "class")
		methodName, _ := stringParam(spec.Params, "method")
		_, err := controllers.CreateJVMChaos(
			d.client,
			ctx,
			spec.Namespace,
			spec.Target,
			chaosmeshv1alpha1.JVMLatencyAction,
			&spec.Duration,
			nil,
			nil,
			chaospkg.WithJVMClass(className),
			chaospkg.WithJVMMethod(methodName),
			chaospkg.WithJVMLatencyDuration(intParam(spec.Params, "latency_duration", 1000)),
		)
		return err
	case "CPUStress":
		containerName, _ := stringParam(spec.Params, "container")
		stressors := controllers.MakeCPUStressors(
			intParam(spec.Params, "cpu_load", 80),
			intParam(spec.Params, "cpu_worker", 1),
		)
		if containerName == "" {
			_, err := controllers.CreateStressChaos(d.client, spec.Namespace, spec.Target, stressors, "cpu-exhaustion", &spec.Duration)
			return err
		}
		_, err := controllers.CreateStressChaosWithContainer(
			d.client,
			ctx,
			spec.Namespace,
			spec.Target,
			stressors,
			"cpu-exhaustion",
			&spec.Duration,
			nil,
			nil,
			[]string{containerName},
		)
		return err
	case "PodFailure":
		_, err := controllers.CreatePodChaos(
			d.client,
			ctx,
			spec.Namespace,
			spec.Target,
			chaosmeshv1alpha1.PodFailureAction,
			&spec.Duration,
			nil,
			nil,
		)
		return err
	default:
		return fmt.Errorf("unsupported direct fault type %q", spec.Type)
	}
}

func (d *DirectSubmitter) DryRun(_ context.Context, spec Spec, w io.Writer) error {
	return NewPrintSubmitter(w).Submit(context.Background(), spec)
}

func intParam(params map[string]any, key string, fallback int) int {
	if params == nil {
		return fallback
	}
	switch v := params[key].(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return fallback
}

func stringParam(params map[string]any, key string) (string, bool) {
	if params == nil {
		return "", false
	}
	v, ok := params[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func networkDirection(v any) (chaosmeshv1alpha1.Direction, error) {
	switch s, _ := v.(string); s {
	case "", "to":
		return chaosmeshv1alpha1.To, nil
	case "from":
		return chaosmeshv1alpha1.From, nil
	case "both":
		return chaosmeshv1alpha1.Both, nil
	default:
		return "", fmt.Errorf("unsupported direction %q", s)
	}
}
