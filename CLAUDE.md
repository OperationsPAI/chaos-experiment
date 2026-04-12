# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a chaos engineering framework for Kubernetes microservices that integrates with Chaos Mesh. It provides programmatic APIs to inject various types of chaos (network, pod, HTTP, JVM, DNS, stress, etc.) into microservice systems for resilience testing.

## Supported Systems

The codebase supports multiple microservice benchmark systems via `internal/systemconfig`:
- **TrainTicket** (`ts`) - Default system
- **OpenTelemetry Demo** (`otel-demo`)
- **Media Microservices** (`media`)
- **Hotel Reservation** (`hs`)
- **Social Network** (`sn`)
- **Online Boutique** (`ob`)

Each system has its own metadata in `internal/{system}/` directories containing service endpoints, database operations, gRPC operations, and Java class methods.

## Build Commands

```bash
# Build the main chaos injection program
go build -o bin/chaos-experiment cmd/main.go

# Build the Java method analyzer
go build -o bin/generate-java-methods cmd/javaanalyzer/main.go

# Build the ClickHouse analyzer for service endpoints
go build -o bin/clickhouse-analyzer cmd/clickhouseanalyzer/main.go

# Build the fault points generator
go build -o bin/faultpoints cmd/faultpoints/main.go

# Build the internal data generator
go build -o bin/internaldata cmd/internaldata/main.go
```

## Code Quality

We use pre-commit hooks to maintain code quality. The configuration includes:
- **golangci-lint**: Comprehensive Go linting with security and style checks
- **gofmt/goimports**: Automatic code formatting and import organization
- **govet**: Static analysis for common Go mistakes
- **Standard hooks**: Trailing whitespace, EOF fixes, file validation

To set up pre-commit:
```bash
# Install pre-commit
pip install pre-commit

# Install git hooks
pre-commit install

# Run on all files
pre-commit run --all-files

# Run specific checks
pre-commit run --all-files --show-diff-on-failure --color=always go-fmt
pre-commit run --all-files --show-diff-on-failure --color=always golangci-lint
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./handler
go test ./internal/resourcelookup
go test ./client

# Run tests with verbose output
go test -v ./...
```

## Code Architecture

### Three-Layer Architecture

1. **chaos/** - Low-level chaos specification builders
   - Creates Chaos Mesh CRD specs with functional options pattern
   - Each file corresponds to a chaos type (network_chaos.go, http_chaos.go, jvm_chaos.go, etc.)
   - Uses `OptChaos` and type-specific option functions (e.g., `OptNetworkChaos`, `OptHTTPChaos`)

2. **controllers/** - Mid-level chaos orchestration
   - Provides simplified APIs for creating and scheduling chaos experiments
   - Handles workflow creation with serial/parallel execution
   - Functions like `CreateNetworkDelayChaos()`, `ScheduleHTTPChaos()`, `AddStressChaosWorkflowNodes()`

3. **handler/** - High-level intelligent chaos generation
   - Provides system-aware chaos generation using metadata
   - Automatically selects appropriate targets based on service endpoints, network dependencies, JVM methods, etc.
   - Main entry point: `handler.go` with functions like `GenerateChaos()`, `GenerateAllChaos()`

### Key Packages

- **internal/systemconfig/** - Global system type configuration (singleton pattern)
  - Use `systemconfig.SetCurrentSystem()` to switch between systems
  - Use `systemconfig.GetCurrentSystem()` to get active system

- **internal/resourcelookup/** - Cached resource lookup with lazy loading
  - Singleton cache manager per system type
  - Functions: `GetAllJVMMethods()`, `GetAllHTTPEndpoints()`, `GetAllNetworkPairs()`, `GetAllDNSEndpoints()`, `GetAllDatabaseOperations()`
  - Call `InvalidateCache()` to clear cached data

- **internal/metadataaccess/** - Unified metadata access layer
  - Automatically routes to correct system-specific metadata based on `systemconfig`
  - Provides functions like `GetEndpointsByService()`, `GetClassMethodsByService()`, `GetOperationsByService()`

- **client/** - Kubernetes client wrapper
  - `GetK8sClient()` returns controller-runtime client
  - `GetLabels()`, `GetContainersWithAppLabel()` for cluster introspection

### Data Generation Tools

The framework includes tools to generate internal metadata from live systems:

1. **Java Method Analyzer** (`cmd/javaanalyzer/`)
   - Analyzes Java source code to extract class methods
   - Generates `internal/{system}/javaclassmethods/javaclassmethods.go`
   - Run: `./bin/generate-java-methods --services /path/to/java/services`

2. **ClickHouse Analyzer** (`cmd/clickhouseanalyzer/`)
   - Analyzes OpenTelemetry traces from ClickHouse
   - Generates service endpoints and database operations
   - Run: `go run cmd/clickhouseanalyzer/main.go --host=HOST --username=USER --password=PASS`

3. **Fault Points Generator** (`cmd/faultpoints/`)
   - Generates comprehensive fault injection points for a system
   - Outputs JSON with all possible chaos targets

## Common Patterns

### Creating Single Chaos Experiments

```go
k8sClient := client.GetK8sClient()
namespace := "ts"
appName := "ts-user-service"

// Network delay
controllers.CreateNetworkDelayChaos(k8sClient, namespace, appName,
    "100ms", "25", "10ms", pointer.String("2m"))

// HTTP abort
abort := true
controllers.CreateHTTPChaos(k8sClient, namespace, appName, "request-abort",
    chaos.WithTarget(chaosmeshv1alpha1.PodHttpRequest),
    chaos.WithPort(8080),
    chaos.WithAbort(&abort))

// JVM latency
controllers.CreateJVMChaos(k8sClient, namespace, appName,
    chaosmeshv1alpha1.JVMLatencyAction, pointer.String("2m"),
    chaos.WithJVMClass("com.example.UserService"),
    chaos.WithJVMMethod("getUserById"),
    chaos.WithJVMLatencyDuration(1000))
```

### Creating Workflow (Sequential Chaos)

```go
workflowSpec := controllers.NewWorkflowSpec(namespace)
sleepTime := pointer.String("15m")
injectTime := pointer.String("5m")

// Add CPU stress
stressors := controllers.MakeCPUStressors(100, 5)
controllers.AddStressChaosWorkflowNodes(workflowSpec, namespace, appList,
    stressors, "cpu", injectTime, sleepTime)

// Add pod failure
controllers.AddPodChaosWorkflowNodes(workflowSpec, namespace, appList,
    chaosmeshv1alpha1.PodFailureAction, injectTime, sleepTime)

// Create workflow
controllers.CreateWorkflow(k8sClient, workflowSpec, namespace)
```

### Using Handler for Intelligent Chaos

```go
// Set the system type first
systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)

// Generate chaos for specific type and target
chaos := handler.GenerateChaos(ctx, handler.NetworkDelay,
    handler.NetworkTarget{
        SourceService: "ts-user-service",
        TargetService: "ts-order-service",
    })

// Generate all possible chaos for a type
allChaos := handler.GenerateAllChaos(ctx, handler.HTTPRequestAbort)
```

## Important Notes

- Always call `systemconfig.SetCurrentSystem()` before using handler or metadata access functions
- The `resourcelookup` package caches data - call `InvalidateCache()` if cluster state changes
- Chaos experiments require Chaos Mesh CRDs to be installed in the cluster
- Use `pointer.String()` from `k8s.io/utils/pointer` for optional duration fields
- Network chaos between services requires knowing the network dependencies (use `GetAllNetworkPairs()`)
- DNS chaos does NOT work for gRPC-only connections (automatically filtered)
- JVM chaos only works with Java services that have the Chaos Mesh agent injected
- Database chaos only supports MySQL operations
