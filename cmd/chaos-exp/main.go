package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/OperationsPAI/chaos-experiment/pkg/guidedcli"
	"gopkg.in/yaml.v3"
)

func main() {
	cfgPath := flag.String("config", "", "Path to the guided CLI config file")
	saveConfig := flag.Bool("save-config", true, "Persist the returned config snapshot to the config file")
	noSaveConfig := flag.Bool("no-save-config", false, "Do not persist the guided session after this call")
	resetConfig := flag.Bool("reset-config", false, "Reset the saved guided session before resolving")
	nextValue := flag.String("next", "", "Apply a single next-step selection using the current guided session state")
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
	returnType := flag.String("return-type", "", "JVM return type: string|int")
	returnValueOpt := flag.String("return-value-opt", "", "JVM return value option: default|random")
	exceptionOpt := flag.String("exception-opt", "", "JVM exception option: default|random")
	memType := flag.String("mem-type", "", "Memory type: heap|stack")
	bodyType := flag.String("body-type", "", "HTTP body type: empty|random")
	replaceMethod := flag.String("replace-method", "", "Replacement HTTP method")
	apply := flag.Bool("apply", false, "Apply the resolved chaos configuration")

	duration := flag.Int("duration", 0, "Duration in minutes, default is 5")
	memorySize := flag.Int("memory-size", 0, "Memory size in MiB")
	memWorker := flag.Int("mem-worker", 0, "Memory stress worker count")
	timeOffset := flag.Int("time-offset", 0, "Time offset in seconds")
	cpuLoad := flag.Int("cpu-load", 0, "CPU load percentage")
	cpuWorker := flag.Int("cpu-worker", 0, "CPU worker count")
	latency := flag.Int("latency", 0, "Network latency in milliseconds")
	correlation := flag.Int("correlation", -1, "Correlation percentage")
	jitter := flag.Int("jitter", -1, "Jitter in milliseconds")
	loss := flag.Int("loss", 0, "Packet loss percentage")
	duplicate := flag.Int("duplicate", 0, "Packet duplication percentage")
	corrupt := flag.Int("corrupt", 0, "Packet corruption percentage")
	rate := flag.Int("rate", 0, "Bandwidth rate in kbps")
	limit := flag.Int("limit", 0, "Bandwidth limit bytes")
	buffer := flag.Int("buffer", 0, "Bandwidth buffer bytes")
	delayDuration := flag.Int("delay-duration", 0, "HTTP delay duration in milliseconds")
	latencyDuration := flag.Int("latency-duration", 0, "JVM latency duration in milliseconds")
	latencyMs := flag.Int("latency-ms", 0, "Database latency in milliseconds")
	cpuCount := flag.Int("cpu-count", 0, "JVM CPU core count")
	statusCode := flag.Int("status-code", 0, "HTTP status code")

	flag.Parse()

	fileCfg, err := guidedcli.LoadConfig(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}
	if *resetConfig {
		fileCfg.GuidedSession = guidedcli.GuidedSession{}
	}
	effectiveSaveConfig := *saveConfig && !*noSaveConfig

	cliCfg := guidedcli.GuidedConfig{
		System:         *system,
		Namespace:      *namespace,
		App:            *app,
		ChaosType:      *chaosType,
		Container:      *container,
		TargetService:  *targetService,
		Domain:         *domain,
		Class:          *class,
		Method:         *method,
		MutatorConfig:  *mutatorConfig,
		Route:          *route,
		HTTPMethod:     *httpMethod,
		Database:       *database,
		Table:          *table,
		Operation:      *operation,
		Direction:      *direction,
		ReturnType:     *returnType,
		ReturnValueOpt: *returnValueOpt,
		ExceptionOpt:   *exceptionOpt,
		MemType:        *memType,
		BodyType:       *bodyType,
		ReplaceMethod:  *replaceMethod,
		Apply:          *apply,
		SaveConfig:     effectiveSaveConfig,
		ResetConfig:    *resetConfig,
	}

	if *duration > 0 {
		cliCfg.Duration = intPtr(*duration)
	}
	if *memorySize > 0 {
		cliCfg.MemorySize = intPtr(*memorySize)
	}
	if *memWorker > 0 {
		cliCfg.MemWorker = intPtr(*memWorker)
	}
	if *timeOffset != 0 {
		cliCfg.TimeOffset = intPtr(*timeOffset)
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
	if *loss > 0 {
		cliCfg.Loss = intPtr(*loss)
	}
	if *duplicate > 0 {
		cliCfg.Duplicate = intPtr(*duplicate)
	}
	if *corrupt > 0 {
		cliCfg.Corrupt = intPtr(*corrupt)
	}
	if *rate > 0 {
		cliCfg.Rate = intPtr(*rate)
	}
	if *limit > 0 {
		cliCfg.Limit = intPtr(*limit)
	}
	if *buffer > 0 {
		cliCfg.Buffer = intPtr(*buffer)
	}
	if *delayDuration > 0 {
		cliCfg.DelayDuration = intPtr(*delayDuration)
	}
	if *latencyDuration > 0 {
		cliCfg.LatencyDuration = intPtr(*latencyDuration)
	}
	if *latencyMs > 0 {
		cliCfg.LatencyMs = intPtr(*latencyMs)
	}
	if *cpuCount > 0 {
		cliCfg.CPUCount = intPtr(*cpuCount)
	}
	if *statusCode > 0 {
		cliCfg.StatusCode = intPtr(*statusCode)
	}

	merged := guidedcli.MergeConfig(fileCfg, cliCfg)
	if *nextValue != "" {
		current, err := guidedcli.Resolve(context.Background(), merged)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to resolve current guided response: %v\n", err)
			os.Exit(1)
		}
		merged, err = guidedcli.ApplyNextSelection(current, *nextValue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to apply --next: %v\n", err)
			os.Exit(1)
		}
		merged.SaveConfig = effectiveSaveConfig
		merged.ResetConfig = *resetConfig
		merged.Apply = *apply
	}

	response, err := guidedcli.Resolve(context.Background(), merged)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve guided response: %v\n", err)
		os.Exit(1)
	}

	if effectiveSaveConfig {
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
