package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/OperationsPAI/chaos-experiment/internal/guidedcli"
	"gopkg.in/yaml.v3"
)

func main() {
	cfgPath := flag.String("config", "", "Path to the guided CLI config file")
	saveConfig := flag.Bool("save-config", false, "Persist the returned config snapshot to the config file")
	resetConfig := flag.Bool("reset-config", false, "Reset the saved guided session before resolving")
	output := flag.String("output", "json", "Output format: json|yaml")

	system := flag.String("system", "", "System namespace instance, for example ts0")
	namespace := flag.String("namespace", "", "Namespace override")
	app := flag.String("app", "", "App label")
	chaosType := flag.String("chaos-type", "", "Chaos type")
	container := flag.String("container", "", "Container name")
	targetService := flag.String("target-service", "", "Target service for network chaos")
	domain := flag.String("domain", "", "Domain for DNS chaos")
	class := flag.String("class", "", "JVM class name")
	method := flag.String("method", "", "JVM method name")
	mutatorConfig := flag.String("mutator-config", "", "Runtime mutator config key")
	route := flag.String("route", "", "HTTP route")
	httpMethod := flag.String("http-method", "", "HTTP method")
	database := flag.String("database", "", "Database name")
	table := flag.String("table", "", "Database table")
	operation := flag.String("operation", "", "Database operation")
	direction := flag.String("direction", "", "Direction: to|from|both")
	apply := flag.Bool("apply", false, "Apply the resolved chaos configuration")

	duration := flag.Int("duration", 0, "Duration in minutes, default is 5")
	cpuLoad := flag.Int("cpu-load", 0, "CPU load percentage")
	cpuWorker := flag.Int("cpu-worker", 0, "CPU worker count")
	latency := flag.Int("latency", 0, "Network latency in milliseconds")
	correlation := flag.Int("correlation", -1, "Correlation percentage")
	jitter := flag.Int("jitter", -1, "Jitter in milliseconds")
	delayDuration := flag.Int("delay-duration", 0, "HTTP delay duration in milliseconds")
	latencyDuration := flag.Int("latency-duration", 0, "JVM latency duration in milliseconds")

	flag.Parse()

	fileCfg, err := guidedcli.LoadConfig(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}
	if *resetConfig {
		fileCfg.GuidedSession = guidedcli.GuidedSession{}
	}

	cliCfg := guidedcli.GuidedConfig{
		System:        *system,
		Namespace:     *namespace,
		App:           *app,
		ChaosType:     *chaosType,
		Container:     *container,
		TargetService: *targetService,
		Domain:        *domain,
		Class:         *class,
		Method:        *method,
		MutatorConfig: *mutatorConfig,
		Route:         *route,
		HTTPMethod:    *httpMethod,
		Database:      *database,
		Table:         *table,
		Operation:     *operation,
		Direction:     *direction,
		Apply:         *apply,
		SaveConfig:    *saveConfig,
		ResetConfig:   *resetConfig,
	}

	if *duration > 0 {
		cliCfg.Duration = intPtr(*duration)
	}
	if *cpuLoad > 0 {
		cliCfg.CPULoad = intPtr(*cpuLoad)
	}
	if *cpuWorker > 0 {
		cliCfg.CPUWorker = intPtr(*cpuWorker)
	}
	if *latency > 0 {
		cliCfg.Latency = intPtr(*latency)
	}
	if *correlation >= 0 {
		cliCfg.Correlation = intPtr(*correlation)
	}
	if *jitter >= 0 {
		cliCfg.Jitter = intPtr(*jitter)
	}
	if *delayDuration > 0 {
		cliCfg.DelayDuration = intPtr(*delayDuration)
	}
	if *latencyDuration > 0 {
		cliCfg.LatencyDuration = intPtr(*latencyDuration)
	}

	merged := guidedcli.MergeConfig(fileCfg, cliCfg)
	response, err := guidedcli.Resolve(context.Background(), merged)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve guided response: %v\n", err)
		os.Exit(1)
	}

	if *saveConfig {
		if err := guidedcli.SaveConfig(*cfgPath, fileCfg, response.Config); err != nil {
			fmt.Fprintf(os.Stderr, "failed to save config: %v\n", err)
			os.Exit(1)
		}
	}

	switch *output {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(response); err != nil {
			fmt.Fprintf(os.Stderr, "failed to encode json response: %v\n", err)
			os.Exit(1)
		}
	case "yaml":
		data, err := yaml.Marshal(response)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to encode yaml response: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprint(os.Stdout, string(data))
	default:
		fmt.Fprintf(os.Stderr, "unsupported output format %q\n", *output)
		os.Exit(1)
	}
}

func intPtr(v int) *int {
	return &v
}
