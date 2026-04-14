package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/OperationsPAI/chaos-experiment/handler"
	"github.com/OperationsPAI/chaos-experiment/internal/resourcelookup"
	"github.com/OperationsPAI/chaos-experiment/internal/serviceendpoints"
	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

var rangePattern = regexp.MustCompile(`^(-?\d+)-(-?\d+)$`)

type FieldProbeResult struct {
	Name       string `json:"name"`
	RangeMin   int    `json:"range_min"`
	RangeMax   int    `json:"range_max"`
	Selected   int    `json:"selected"`
	Source     string `json:"source"`
	Candidates int    `json:"candidates"`
}

type ChaosProbeResult struct {
	ChaosType     string                 `json:"chaos_type"`
	Status        string                 `json:"status"`
	Error         string                 `json:"error,omitempty"`
	Fields        []FieldProbeResult     `json:"fields"`
	DisplayConfig map[string]interface{} `json:"display_config,omitempty"`
}

type SystemProbeResult struct {
	System          string             `json:"system"`
	SelectedService string             `json:"selected_service"`
	ServiceCount    int                `json:"service_count"`
	Results         []ChaosProbeResult `json:"results"`
}

type ProbeReport struct {
	Seed      int64               `json:"seed"`
	CreatedAt string              `json:"created_at"`
	Systems   []SystemProbeResult `json:"systems"`
}

func main() {
	var outputPath string
	var seed int64
	flag.StringVar(&outputPath, "output", "artifacts/handler_config_probe_report.json", "Output report file path")
	flag.Int64Var(&seed, "seed", time.Now().UnixNano(), "Random seed")
	flag.Parse()

	rng := rand.New(rand.NewSource(seed))
	ctx := context.Background()

	report := ProbeReport{
		Seed:      seed,
		CreatedAt: time.Now().Format(time.RFC3339),
		Systems:   make([]SystemProbeResult, 0),
	}

	systems := systemconfig.GetAllSystemTypes()
	for _, system := range systems {
		if err := systemconfig.SetCurrentSystem(system); err != nil {
			fmt.Fprintf(os.Stderr, "skip system %s: %v\n", system, err)
			continue
		}

		services := serviceendpoints.GetAllServices()
		sort.Strings(services)
		if len(services) == 0 {
			report.Systems = append(report.Systems, SystemProbeResult{
				System:          system.String(),
				SelectedService: "",
				ServiceCount:    0,
				Results:         []ChaosProbeResult{},
			})
			continue
		}

		selectedService := services[rng.Intn(len(services))]
		namespace, nsErr := systemconfig.GetNamespaceByIndex(system, 0)
		if nsErr != nil {
			namespace = ""
		}

		result := SystemProbeResult{
			System:          system.String(),
			SelectedService: selectedService,
			ServiceCount:    len(services),
			Results:         make([]ChaosProbeResult, 0),
		}

		for chaosType := handler.PodKill; chaosType <= handler.JVMRuntimeMutator; chaosType++ {
			chaosName := handler.GetChaosTypeName(chaosType)
			specAny, ok := handler.SpecMap[chaosType]
			if !ok {
				result.Results = append(result.Results, ChaosProbeResult{ChaosType: chaosName, Status: "skipped", Error: "spec missing"})
				continue
			}

			specType := reflect.TypeOf(specAny)
			specVal := reflect.New(specType).Elem()
			fieldResults := make([]FieldProbeResult, 0)

			for i := 0; i < specType.NumField(); i++ {
				f := specType.Field(i)
				rangeTag := f.Tag.Get("range")
				if rangeTag == "" {
					continue
				}

				minV, maxV, err := parseRangeTag(rangeTag)
				if err != nil {
					continue
				}

				selected := minV
				source := "static"
				candidateCount := max(0, maxV-minV+1)

				if f.Tag.Get("dynamic") == "true" {
					dynMin, dynMax, candidates, dynSource := resolveDynamicField(ctx, namespace, system, f.Name, selectedService)
					minV, maxV = dynMin, dynMax
					candidateCount = len(candidates)
					source = dynSource

					if len(candidates) > 0 {
						selected = candidates[rng.Intn(len(candidates))]
					} else if maxV >= minV {
						selected = minV + rng.Intn(maxV-minV+1)
						source = source + ":fallback-global"
					} else {
						selected = minV
						source = source + ":fallback-zero"
					}
				} else {
					if maxV >= minV {
						selected = minV + rng.Intn(maxV-minV+1)
					}
				}

				specField := specVal.Field(i)
				if specField.CanSet() {
					specField.SetInt(int64(selected))
				}

				fieldResults = append(fieldResults, FieldProbeResult{
					Name:       f.Name,
					RangeMin:   minV,
					RangeMax:   maxV,
					Selected:   selected,
					Source:     source,
					Candidates: candidateCount,
				})
			}

			conf := handler.InjectionConf{}
			confVal := reflect.ValueOf(&conf).Elem()
			confField := confVal.FieldByName(chaosName)
			if !confField.IsValid() {
				result.Results = append(result.Results, ChaosProbeResult{
					ChaosType: chaosName,
					Status:    "error",
					Error:     "injection config field not found",
					Fields:    fieldResults,
				})
				continue
			}

			specPtr := reflect.New(specType)
			specPtr.Elem().Set(specVal)
			confField.Set(specPtr)

			display, err := conf.GetDisplayConfig(ctx)
			if err != nil {
				result.Results = append(result.Results, ChaosProbeResult{
					ChaosType: chaosName,
					Status:    "error",
					Error:     err.Error(),
					Fields:    fieldResults,
				})
				continue
			}

			result.Results = append(result.Results, ChaosProbeResult{
				ChaosType:     chaosName,
				Status:        "ok",
				Fields:        fieldResults,
				DisplayConfig: display,
			})
		}

		report.Systems = append(report.Systems, result)
	}

	if err := writeReport(outputPath, report); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Handler config probe report written to %s\n", outputPath)
}

func resolveDynamicField(
	ctx context.Context,
	namespace string,
	system systemconfig.SystemType,
	fieldName string,
	service string,
) (int, int, []int, string) {
	cache := resourcelookup.GetSystemCache(system)
	all := []int{}
	matched := []int{}

	systems := systemconfig.GetAllSystemTypes()
	if fieldName == "System" {
		all = make([]int, len(systems))
		for i := range systems {
			all[i] = i
		}
		target := 0
		for i, s := range systems {
			if s == system {
				target = i
				break
			}
		}
		return 0, max(0, len(systems)-1), []int{target}, "dynamic:system"
	}

	switch fieldName {
	case "AppIdx":
		labels, err := cache.GetAllAppLabels(ctx, namespace, "app")
		if err != nil {
			return 0, 0, nil, "dynamic:app:error"
		}
		for i, v := range labels {
			all = append(all, i)
			if v == service {
				matched = append(matched, i)
			}
		}
		return rangeFromIndices(all), rangeToIndices(all), matched, "dynamic:app"
	case "MethodIdx":
		methods, err := cache.GetAllJVMMethods()
		if err != nil {
			return 0, 0, nil, "dynamic:method:error"
		}
		for i, v := range methods {
			all = append(all, i)
			if v.AppName == service {
				matched = append(matched, i)
			}
		}
		return rangeFromIndices(all), rangeToIndices(all), matched, "dynamic:method"
	case "MutatorTargetIdx":
		targets, err := cache.GetAllJVMRuntimeMutatorTargets()
		if err != nil {
			return 0, 0, nil, "dynamic:target:error"
		}
		for i, v := range targets {
			all = append(all, i)
			if v.AppName == service {
				matched = append(matched, i)
			}
		}
		return rangeFromIndices(all), rangeToIndices(all), matched, "dynamic:target"
	case "EndpointIdx":
		endpoints, err := cache.GetAllHTTPEndpoints()
		if err != nil {
			return 0, 0, nil, "dynamic:endpoint:error"
		}
		for i, v := range endpoints {
			all = append(all, i)
			if v.AppName == service {
				matched = append(matched, i)
			}
		}
		return rangeFromIndices(all), rangeToIndices(all), matched, "dynamic:endpoint"
	case "NetworkPairIdx":
		pairs, err := cache.GetAllNetworkPairs()
		if err != nil {
			return 0, 0, nil, "dynamic:network:error"
		}
		for i, v := range pairs {
			all = append(all, i)
			if v.SourceService == service || v.TargetService == service {
				matched = append(matched, i)
			}
		}
		return rangeFromIndices(all), rangeToIndices(all), matched, "dynamic:network"
	case "ContainerIdx":
		containers, err := cache.GetAllContainers(ctx, namespace)
		if err != nil {
			return 0, 0, nil, "dynamic:container:error"
		}
		for i, v := range containers {
			all = append(all, i)
			if v.AppLabel == service {
				matched = append(matched, i)
			}
		}
		return rangeFromIndices(all), rangeToIndices(all), matched, "dynamic:container"
	case "DNSEndpointIdx":
		dnsEndpoints, err := cache.GetAllDNSEndpoints()
		if err != nil {
			return 0, 0, nil, "dynamic:dns:error"
		}
		for i, v := range dnsEndpoints {
			all = append(all, i)
			if v.AppName == service {
				matched = append(matched, i)
			}
		}
		return rangeFromIndices(all), rangeToIndices(all), matched, "dynamic:dns"
	case "DatabaseIdx":
		dbOps, err := cache.GetAllDatabaseOperations()
		if err != nil {
			return 0, 0, nil, "dynamic:database:error"
		}
		for i, v := range dbOps {
			all = append(all, i)
			if v.AppName == service {
				matched = append(matched, i)
			}
		}
		return rangeFromIndices(all), rangeToIndices(all), matched, "dynamic:database"
	default:
		return 0, 0, nil, "dynamic:unknown"
	}
}

func parseRangeTag(tag string) (int, int, error) {
	m := rangePattern.FindStringSubmatch(strings.TrimSpace(tag))
	if len(m) != 3 {
		return 0, 0, fmt.Errorf("invalid range tag: %s", tag)
	}

	var start, end int
	_, err := fmt.Sscanf(m[1], "%d", &start)
	if err != nil {
		return 0, 0, err
	}
	_, err = fmt.Sscanf(m[2], "%d", &end)
	if err != nil {
		return 0, 0, err
	}
	return start, end, nil
}

func rangeFromIndices(indices []int) int {
	if len(indices) == 0 {
		return 0
	}
	return indices[0]
}

func rangeToIndices(indices []int) int {
	if len(indices) == 0 {
		return 0
	}
	return indices[len(indices)-1]
}

func writeReport(outputPath string, report ProbeReport) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, b, 0644)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
