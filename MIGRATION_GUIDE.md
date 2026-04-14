# Migration Guide: Using the New Architecture

## Overview

The new architecture provides complete decoupling between:
1. **Internal Data Storage** (`internal/resourcetypes`) - Raw data types
2. **Handler Endpoints** (`internal/endpoint`) - Unified conversion layer
3. **System Data Model** (`internal/model`) - Container for all resources
4. **Registry** (`internal/registry`) - Auto-registration pattern
5. **Adapter** (`internal/adapter`) - Compatibility layer for existing generated files

## Current Status

### ✅ Completed
- `internal/resourcetypes/types.go` - Decoupled raw data types (HTTPEndpoint, RPCOperation, DatabaseOperation)
- `internal/endpoint/types.go` - Handler endpoint types (Endpoint, CallPair, HTTPEndpointInfo, DNSEndpointInfo, DatabaseInfo)
- `internal/endpoint/converter.go` - Conversion functions from internal data to endpoints
- `internal/model/systemdata.go` - Unified SystemData container
- `internal/registry/registry.go` - Registry pattern implementation
- `internal/adapter/adapter.go` - **NEW**: Bridges old generated code with new registry
- `tools/clickhouseanalyzer` - NO LONGER imports internal packages (fully decoupled)
- **`cmd/faultpoints/main.go`** - **MIGRATED**: Now uses registry pattern!

### ✅ FULLY FUNCTIONAL
The new architecture is **actively in use**:
- cmd/faultpoints eliminates 18+ system imports, uses registry pattern
- Tested and working with 8 registry-registered systems (ts, otel-demo, media, hs, sn, ob, sockshop, teastore)
- Adapter layer provides backward compatibility with existing generated files

### 🔄 In Progress
- Handlers are being migrated behind `handler/endpoint_provider.go` (HTTP/network/DNS/database + app/container/JVM helpers)
- `jvm_chaos.go` now uses provider-based JVM method selection (from `internal/javaclassmethods` + registry services)
- Generate single `data.go` per system with auto-registration

### ❌ Not Started
- Final cleanup: remove `resourcelookup` dependency from handler initialization/cache preloading path

### ✅ Recently Migrated (Important)
- `internal/resourcelookup/lookup.go` is already migrated to registry access (`registry.MustGetCurrent()`)
- `cmd/faultpoints/main.go` is already migrated and uses adapter + registry
- SockShop/TeaStore are now auto-registered in `internal/adapter/adapter.go`

## How to Use the New Architecture

### Real Example: cmd/faultpoints (MIGRATED)

The cmd/faultpoints has been fully migrated and demonstrates the new pattern:

```go
package main

import (
	_ "github.com/OperationsPAI/chaos-experiment/internal/adapter" // Auto-registers all systems
	"github.com/OperationsPAI/chaos-experiment/internal/endpoint"
	"github.com/OperationsPAI/chaos-experiment/internal/registry"
	"github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

func main() {
	// Set system (done once at startup)
	systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)
	
	// Verify registration
	if !registry.IsRegistered(systemconfig.GetCurrentSystem()) {
		panic("System not registered")
	}
	
	// Use the system data
	listHTTPEndpoints()
}

func listHTTPEndpoints() {
	sysData := registry.MustGetCurrent()  // No switch-case needed!
	
	var endpoints []endpoint.HTTPEndpointInfo
	for _, service := range sysData.GetAllServices() {
		for _, ep := range sysData.GetHTTPEndpointsByService(service) {
			if ep.Route != "" {
				endpoints = append(endpoints, endpoint.ToHTTPEndpointInfo(ep))
			}
		}
	}
	// Use endpoints...
}
```

**What this eliminates:**
- ❌ No more 18+ system-specific imports
- ❌ No more switch-case on system type
- ❌ No more manual system initialization
- ✅ Just import adapter and use registry!

### For New Code (Recommended Pattern)

#### 1. Accessing System Data via Registry

```go
import (
    "github.com/OperationsPAI/chaos-experiment/internal/registry"
    "github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

// Set current system (usually done once at startup)
systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)

// Get system data from registry
sysData := registry.MustGetCurrent()

// Access HTTP endpoints
httpEndpoints := sysData.GetHTTPEndpointsByService("ts-order-service")

// Access RPC operations
rpcOps := sysData.GetRPCOperationsByService("frontend")

// Access database operations
dbOps := sysData.GetDatabaseOperationsByService("account service")
```

#### 2. Converting to Handler Endpoints

```go
import (
    "github.com/OperationsPAI/chaos-experiment/internal/endpoint"
    "github.com/OperationsPAI/chaos-experiment/internal/registry"
)

sysData := registry.MustGetCurrent()

// For HTTP Chaos - get HTTP-specific info
httpEndpoints := sysData.GetHTTPEndpointsByService("myservice")
for _, ep := range httpEndpoints {
    httpInfo := endpoint.ToHTTPEndpointInfo(ep)
    // Use httpInfo for HTTP chaos
}

// For Network Chaos - combine all operation types
// (Build CallPair from HTTP + RPC + DB endpoints)
var networkPairs []endpoint.CallPair
// ... collect from HTTP, RPC, and DB ...

// For DNS Chaos - HTTP + DB only (exclude RPC)
// (Build DNSEndpointInfo from HTTP + DB endpoints)

// For MySQL Chaos - filter by DB system
mysqlOps := sysData.GetDatabaseOperationsByDBSystem("mysql")
for _, op := range mysqlOps {
    dbInfo := endpoint.ToDatabaseInfo(op)
    // Use dbInfo for MySQL chaos
}
```

#### 3. Example: Network Chaos Handler

```go
func GetNetworkPairs() []endpoint.CallPair {
    sysData := registry.MustGetCurrent()
    pairMap := make(map[string]*endpoint.CallPair)
    
    // Collect HTTP pairs
    for _, service := range sysData.GetAllServices() {
        for _, ep := range sysData.GetHTTPEndpointsByService(service) {
            if ep.ServerAddress != "" {
                key := ep.ServiceName + "->" + ep.ServerAddress
                // Add to pairMap with operation type "http"
            }
        }
    }
    
    // Collect RPC pairs
    for _, service := range sysData.GetAllRPCServices() {
        for _, op := range sysData.GetRPCOperationsByService(service) {
            // Add to pairMap with operation type "rpc"
        }
    }
    
    // Collect DB pairs
    for _, service := range sysData.GetAllDatabaseServices() {
        for _, op := range sysData.GetDatabaseOperationsByService(service) {
            // Add to pairMap with operation type "db"
        }
    }
    
    return convertMapToSlice(pairMap)
}
```

## Benefits of New Architecture

### 1. Complete Decoupling
- **Tool Layer** (clickhouseanalyzer): Defines its OWN types, doesn't import internal packages
- **Storage Layer** (resourcetypes): Raw data storage, type-specific fields only
- **Model Layer** (model): Container for all resource types
- **Conversion Layer** (endpoint): Transforms storage to handler needs
- **Registry Layer** (registry): Auto-registration, no switch-case

### 2. No More Switch-Case Statements
OLD:
```go
switch systemconfig.GetCurrentSystem() {
case systemconfig.SystemTrainTicket:
    return tsendpoints.GetEndpoints()
case systemconfig.SystemOtelDemo:
    return oteldemoendpoints.GetEndpoints()
// ... 6 systems ...
}
```

NEW:
```go
sysData := registry.MustGetCurrent()
return sysData.GetHTTPEndpointsByService(service)
```

### 3. No System-Specific Imports
OLD (cmd/faultpoints):
```go
import (
    hsendpoints "github.com/OperationsPAI/chaos-experiment/internal/hs/serviceendpoints"
    mediaendpoints "github.com/OperationsPAI/chaos-experiment/internal/media/serviceendpoints"
    // ... 18+ imports ...
)
```

NEW:
```go
import (
    "github.com/OperationsPAI/chaos-experiment/internal/registry"
    "github.com/OperationsPAI/chaos-experiment/internal/endpoint"
)
```

### 4. Type-Specific Data Storage
Each resource type only stores relevant fields:
- **HTTPEndpoint**: Method, Route, Status (HTTP-specific)
- **RPCOperation**: RPCSystem, RPCService, RPCMethod (RPC-specific)
- **DatabaseOperation**: DBName, DBTable, Operation, DBSystem (DB-specific)

### 5. Fault-Type-Specific Endpoints
Different handlers get different endpoint types based on their needs:
- **Network Chaos**: `CallPair` (all operation types)
- **HTTP Chaos**: `HTTPEndpointInfo` (HTTP only)
- **DNS Chaos**: `DNSEndpointInfo` (HTTP + DB, excludes RPC)
- **MySQL Chaos**: `DatabaseInfo` (MySQL DB only)

## Next Steps

### Phase 1: Update Code Generator (HIGH PRIORITY)
Modify `tools/clickhouseanalyzer/datagenerator.go` to:
1. Generate single `data.go` per system instead of 3 files
2. Add `init()` function that registers data with registry
3. Keep types as aliases to `resourcetypes`

Example generated file structure:
```go
// internal/ts/data/data.go
package data

import (
    "github.com/OperationsPAI/chaos-experiment/internal/model"
    "github.com/OperationsPAI/chaos-experiment/internal/registry"
    "github.com/OperationsPAI/chaos-experiment/internal/resourcetypes"
    "github.com/OperationsPAI/chaos-experiment/internal/systemconfig"
)

func init() {
    // Auto-register on import
    registry.Register(systemconfig.SystemTrainTicket, &model.SystemData{
        SystemName: "ts",
        HTTPEndpoints: map[string][]resourcetypes.HTTPEndpoint{ /* ... */ },
        RPCOperations: map[string][]resourcetypes.RPCOperation{ /* ... */ },
        DatabaseOperations: map[string][]resourcetypes.DatabaseOperation{ /* ... */ },
        AllServices: []string{ /* ... */ },
    })
}
```

### Phase 2: Align with injectionv2 branch (HIGH PRIORITY)
Large overlap exists with `injectionv2` in these areas:
- `internal/resourcetypes`, `internal/model`, `internal/registry`, `internal/resourcelookup`
- `tools/clickhouseanalyzer/datagenerator.go`
- `cmd/faultpoints/main.go`

Recommended merge strategy:
1. Merge core data layers first (`resourcetypes` → `model` → `registry`)
2. Merge generated data pipelines (`tools/clickhouseanalyzer/*`, generated `internal/*/{serviceendpoints,grpcoperations,databaseoperations}`)
3. Merge runtime consumers last (`resourcelookup`, `cmd/faultpoints`, handlers)
4. Run full validation after each layer to prevent conflict cascades

### Phase 3: Update Handlers (LOW PRIORITY)
- ✅ Added provider layer (`handler/endpoint_provider.go`) and migrated core chaos paths (HTTP, network, DNS, JVM MySQL)
- ✅ Migrated app/container/JVM method index resolution in handlers to provider helpers
- Remaining: deprecate direct `resourcelookup` usage in handler initialization (`InitCaches`/`PreloadCaches`) once replacement cache strategy is finalized

## Testing the New Code

Currently, the new architecture can be tested but requires:
1. Manually populating registry with test data
2. Using the conversion functions from `internal/endpoint`

Example test:
```go
func TestRegistryPattern(t *testing.T) {
    // Setup
    testData := &model.SystemData{
        HTTPEndpoints: map[string][]resourcetypes.HTTPEndpoint{
            "service1": {{ServiceName: "service1", Route: "/api/test"}},
        },
    }
    registry.Register(systemconfig.SystemTrainTicket, testData)
    systemconfig.SetCurrentSystem(systemconfig.SystemTrainTicket)
    
    // Use
    sysData := registry.MustGetCurrent()
    eps := sysData.GetHTTPEndpointsByService("service1")
    assert.Equal(t, 1, len(eps))
}
```
