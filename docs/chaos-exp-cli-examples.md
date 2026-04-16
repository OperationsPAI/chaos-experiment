# chaos-exp Guided CLI Examples

This document records full raw JSON responses for two guided sessions using the new shorthand flow.

Notes:

- output format: `json`
- state persistence: enabled by default
- isolated session files were used so the examples do not interfere with each other
- the responses below were captured live from the `ts` namespace on 2026-04-17

## Example 1: Simple Flow (`PodKill`)

Session file used in the capture:

```text
%TEMP%\\chaos-exp-examples-podkill.yaml
```

### Call 1

```bash
./chaos-exp.exe -output json --config %TEMP%\chaos-exp-examples-podkill.yaml --reset-config --namespace ts
```

```json
{
  "mode": "guided",
  "stage": "select_app",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts"
  },
  "resolved": {
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "app",
      "kind": "enum",
      "required": true,
      "description": "Select an app label in the namespace",
      "options": [
        {
          "value": "ts-admin-basic-info-service",
          "label": "ts-admin-basic-info-service"
        },
        {
          "value": "ts-admin-order-service",
          "label": "ts-admin-order-service"
        },
        {
          "value": "ts-admin-route-service",
          "label": "ts-admin-route-service"
        },
        {
          "value": "ts-admin-travel-service",
          "label": "ts-admin-travel-service"
        },
        {
          "value": "ts-admin-user-service",
          "label": "ts-admin-user-service"
        },
        {
          "value": "ts-assurance-service",
          "label": "ts-assurance-service"
        },
        {
          "value": "ts-auth-service",
          "label": "ts-auth-service"
        },
        {
          "value": "ts-basic-service",
          "label": "ts-basic-service"
        },
        {
          "value": "ts-cancel-service",
          "label": "ts-cancel-service"
        },
        {
          "value": "ts-config-service",
          "label": "ts-config-service"
        },
        {
          "value": "ts-consign-price-service",
          "label": "ts-consign-price-service"
        },
        {
          "value": "ts-consign-service",
          "label": "ts-consign-service"
        },
        {
          "value": "ts-contacts-service",
          "label": "ts-contacts-service"
        },
        {
          "value": "ts-delivery-service",
          "label": "ts-delivery-service"
        },
        {
          "value": "ts-execute-service",
          "label": "ts-execute-service"
        },
        {
          "value": "ts-food-delivery-service",
          "label": "ts-food-delivery-service"
        },
        {
          "value": "ts-food-service",
          "label": "ts-food-service"
        },
        {
          "value": "ts-inside-payment-service",
          "label": "ts-inside-payment-service"
        },
        {
          "value": "ts-notification-service",
          "label": "ts-notification-service"
        },
        {
          "value": "ts-order-other-service",
          "label": "ts-order-other-service"
        },
        {
          "value": "ts-order-service",
          "label": "ts-order-service"
        },
        {
          "value": "ts-payment-service",
          "label": "ts-payment-service"
        },
        {
          "value": "ts-preserve-service",
          "label": "ts-preserve-service"
        },
        {
          "value": "ts-price-service",
          "label": "ts-price-service"
        },
        {
          "value": "ts-rebook-service",
          "label": "ts-rebook-service"
        },
        {
          "value": "ts-route-plan-service",
          "label": "ts-route-plan-service"
        },
        {
          "value": "ts-route-service",
          "label": "ts-route-service"
        },
        {
          "value": "ts-seat-service",
          "label": "ts-seat-service"
        },
        {
          "value": "ts-security-service",
          "label": "ts-security-service"
        },
        {
          "value": "ts-station-food-service",
          "label": "ts-station-food-service"
        },
        {
          "value": "ts-station-service",
          "label": "ts-station-service"
        },
        {
          "value": "ts-train-food-service",
          "label": "ts-train-food-service"
        },
        {
          "value": "ts-train-service",
          "label": "ts-train-service"
        },
        {
          "value": "ts-travel-plan-service",
          "label": "ts-travel-plan-service"
        },
        {
          "value": "ts-travel-service",
          "label": "ts-travel-service"
        },
        {
          "value": "ts-travel2-service",
          "label": "ts-travel2-service"
        },
        {
          "value": "ts-ui-dashboard",
          "label": "ts-ui-dashboard"
        },
        {
          "value": "ts-user-service",
          "label": "ts-user-service"
        },
        {
          "value": "ts-wait-order-service",
          "label": "ts-wait-order-service"
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
./chaos-exp.exe -output json --config %TEMP%\chaos-exp-examples-podkill.yaml --next ts-auth-service
```

```json
{
  "mode": "guided",
  "stage": "select_chaos_type",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service"
  },
  "resolved": {
    "app": "ts-auth-service",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "chaos_type",
      "kind": "enum",
      "required": true,
      "description": "Select a chaos type supported for the current app",
      "options": [
        {
          "value": "CPUStress",
          "label": "CPUStress",
          "description": "Stress a container with CPU load"
        },
        {
          "value": "ContainerKill",
          "label": "ContainerKill",
          "description": "Kill a specific container under the app"
        },
        {
          "value": "DNSError",
          "label": "DNSError",
          "description": "Return DNS errors for a selected domain"
        },
        {
          "value": "DNSRandom",
          "label": "DNSRandom",
          "description": "Return random DNS results for a selected domain"
        },
        {
          "value": "HTTPRequestAbort",
          "label": "HTTPRequestAbort",
          "description": "Abort HTTP requests for a selected endpoint"
        },
        {
          "value": "HTTPRequestDelay",
          "label": "HTTPRequestDelay",
          "description": "Delay HTTP requests for a selected endpoint"
        },
        {
          "value": "HTTPRequestReplaceMethod",
          "label": "HTTPRequestReplaceMethod",
          "description": "Replace request methods"
        },
        {
          "value": "HTTPRequestReplacePath",
          "label": "HTTPRequestReplacePath",
          "description": "Replace request paths"
        },
        {
          "value": "HTTPResponseAbort",
          "label": "HTTPResponseAbort",
          "description": "Abort HTTP responses for a selected endpoint"
        },
        {
          "value": "HTTPResponseDelay",
          "label": "HTTPResponseDelay",
          "description": "Delay HTTP responses for a selected endpoint"
        },
        {
          "value": "HTTPResponsePatchBody",
          "label": "HTTPResponsePatchBody",
          "description": "Patch HTTP response bodies"
        },
        {
          "value": "HTTPResponseReplaceBody",
          "label": "HTTPResponseReplaceBody",
          "description": "Replace HTTP response bodies"
        },
        {
          "value": "HTTPResponseReplaceCode",
          "label": "HTTPResponseReplaceCode",
          "description": "Replace response status codes"
        },
        {
          "value": "JVMCPUStress",
          "label": "JVMCPUStress",
          "description": "Stress CPU inside a JVM method"
        },
        {
          "value": "JVMException",
          "label": "JVMException",
          "description": "Throw a JVM exception"
        },
        {
          "value": "JVMGarbageCollector",
          "label": "JVMGarbageCollector",
          "description": "Trigger JVM garbage collection for the app"
        },
        {
          "value": "JVMLatency",
          "label": "JVMLatency",
          "description": "Inject latency into a JVM method"
        },
        {
          "value": "JVMMemoryStress",
          "label": "JVMMemoryStress",
          "description": "Stress memory inside a JVM method"
        },
        {
          "value": "JVMMySQLException",
          "label": "JVMMySQLException",
          "description": "Inject SQL exceptions into a MySQL operation"
        },
        {
          "value": "JVMMySQLLatency",
          "label": "JVMMySQLLatency",
          "description": "Inject latency into a MySQL operation"
        },
        {
          "value": "JVMReturn",
          "label": "JVMReturn",
          "description": "Override a JVM return value"
        },
        {
          "value": "JVMRuntimeMutator",
          "label": "JVMRuntimeMutator",
          "description": "Apply a runtime mutator strategy to a JVM method"
        },
        {
          "value": "MemoryStress",
          "label": "MemoryStress",
          "description": "Stress a container with memory pressure"
        },
        {
          "value": "NetworkBandwidth",
          "label": "NetworkBandwidth",
          "description": "Limit bandwidth to a downstream service"
        },
        {
          "value": "NetworkCorrupt",
          "label": "NetworkCorrupt",
          "description": "Corrupt traffic to a downstream service"
        },
        {
          "value": "NetworkDelay",
          "label": "NetworkDelay",
          "description": "Delay traffic to a downstream service"
        },
        {
          "value": "NetworkDuplicate",
          "label": "NetworkDuplicate",
          "description": "Duplicate traffic to a downstream service"
        },
        {
          "value": "NetworkLoss",
          "label": "NetworkLoss",
          "description": "Drop traffic to a downstream service"
        },
        {
          "value": "NetworkPartition",
          "label": "NetworkPartition",
          "description": "Partition traffic to a downstream service"
        },
        {
          "value": "PodFailure",
          "label": "PodFailure",
          "description": "Fail pods for the selected app"
        },
        {
          "value": "PodKill",
          "label": "PodKill",
          "description": "Kill pods for the selected app"
        },
        {
          "value": "TimeSkew",
          "label": "TimeSkew",
          "description": "Shift time in a specific container"
        }
      ]
    }
  ],
  "preview": {
    "resource_summary": {
      "containers": 1,
      "database_operations": 10,
      "dns_domains": 2,
      "http_endpoints": 1,
      "jvm_methods": 29,
      "network_targets": 2,
      "runtime_mutator_methods": 18
    }
  },
  "can_apply": false
}
```

### Call 3

```bash
./chaos-exp.exe -output json --config %TEMP%\chaos-exp-examples-podkill.yaml --next PodKill
```

```json
{
  "mode": "guided",
  "stage": "ready_to_apply",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "PodKill",
    "duration": 5
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "PodKill",
    "duration": 5,
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "duration",
      "kind": "number_range",
      "required": false,
      "description": "Fault duration in minutes",
      "min": 1,
      "max": 60,
      "step": 1,
      "default": 5,
      "unit": "minute"
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "PodKill": {
          "AppIdx": 6,
          "Duration": 5,
          "System": 7
        }
      },
      "chaos_type": "PodKill",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "service": [
        "ts-auth-service"
      ]
    }
  },
  "apply_payload": {
    "PodKill": {
      "AppIdx": 6,
      "Duration": 5,
      "System": 7
    }
  },
  "can_apply": true
}
```

## Example 2: Complex Flow (`JVMRuntimeMutator`)

Session file used in the capture:

```text
%TEMP%\\chaos-exp-examples-mutator.yaml
```

### Call 1

```bash
./chaos-exp.exe -output json --config %TEMP%\chaos-exp-examples-mutator.yaml --reset-config --namespace ts
```

```json
{
  "mode": "guided",
  "stage": "select_app",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts"
  },
  "resolved": {
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "app",
      "kind": "enum",
      "required": true,
      "description": "Select an app label in the namespace",
      "options": [
        {
          "value": "ts-admin-basic-info-service",
          "label": "ts-admin-basic-info-service"
        },
        {
          "value": "ts-admin-order-service",
          "label": "ts-admin-order-service"
        },
        {
          "value": "ts-admin-route-service",
          "label": "ts-admin-route-service"
        },
        {
          "value": "ts-admin-travel-service",
          "label": "ts-admin-travel-service"
        },
        {
          "value": "ts-admin-user-service",
          "label": "ts-admin-user-service"
        },
        {
          "value": "ts-assurance-service",
          "label": "ts-assurance-service"
        },
        {
          "value": "ts-auth-service",
          "label": "ts-auth-service"
        },
        {
          "value": "ts-basic-service",
          "label": "ts-basic-service"
        },
        {
          "value": "ts-cancel-service",
          "label": "ts-cancel-service"
        },
        {
          "value": "ts-config-service",
          "label": "ts-config-service"
        },
        {
          "value": "ts-consign-price-service",
          "label": "ts-consign-price-service"
        },
        {
          "value": "ts-consign-service",
          "label": "ts-consign-service"
        },
        {
          "value": "ts-contacts-service",
          "label": "ts-contacts-service"
        },
        {
          "value": "ts-delivery-service",
          "label": "ts-delivery-service"
        },
        {
          "value": "ts-execute-service",
          "label": "ts-execute-service"
        },
        {
          "value": "ts-food-delivery-service",
          "label": "ts-food-delivery-service"
        },
        {
          "value": "ts-food-service",
          "label": "ts-food-service"
        },
        {
          "value": "ts-inside-payment-service",
          "label": "ts-inside-payment-service"
        },
        {
          "value": "ts-notification-service",
          "label": "ts-notification-service"
        },
        {
          "value": "ts-order-other-service",
          "label": "ts-order-other-service"
        },
        {
          "value": "ts-order-service",
          "label": "ts-order-service"
        },
        {
          "value": "ts-payment-service",
          "label": "ts-payment-service"
        },
        {
          "value": "ts-preserve-service",
          "label": "ts-preserve-service"
        },
        {
          "value": "ts-price-service",
          "label": "ts-price-service"
        },
        {
          "value": "ts-rebook-service",
          "label": "ts-rebook-service"
        },
        {
          "value": "ts-route-plan-service",
          "label": "ts-route-plan-service"
        },
        {
          "value": "ts-route-service",
          "label": "ts-route-service"
        },
        {
          "value": "ts-seat-service",
          "label": "ts-seat-service"
        },
        {
          "value": "ts-security-service",
          "label": "ts-security-service"
        },
        {
          "value": "ts-station-food-service",
          "label": "ts-station-food-service"
        },
        {
          "value": "ts-station-service",
          "label": "ts-station-service"
        },
        {
          "value": "ts-train-food-service",
          "label": "ts-train-food-service"
        },
        {
          "value": "ts-train-service",
          "label": "ts-train-service"
        },
        {
          "value": "ts-travel-plan-service",
          "label": "ts-travel-plan-service"
        },
        {
          "value": "ts-travel-service",
          "label": "ts-travel-service"
        },
        {
          "value": "ts-travel2-service",
          "label": "ts-travel2-service"
        },
        {
          "value": "ts-ui-dashboard",
          "label": "ts-ui-dashboard"
        },
        {
          "value": "ts-user-service",
          "label": "ts-user-service"
        },
        {
          "value": "ts-wait-order-service",
          "label": "ts-wait-order-service"
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
./chaos-exp.exe -output json --config %TEMP%\chaos-exp-examples-mutator.yaml --next ts-delivery-service
```

```json
{
  "mode": "guided",
  "stage": "select_chaos_type",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-delivery-service"
  },
  "resolved": {
    "app": "ts-delivery-service",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "chaos_type",
      "kind": "enum",
      "required": true,
      "description": "Select a chaos type supported for the current app",
      "options": [
        {
          "value": "CPUStress",
          "label": "CPUStress",
          "description": "Stress a container with CPU load"
        },
        {
          "value": "ContainerKill",
          "label": "ContainerKill",
          "description": "Kill a specific container under the app"
        },
        {
          "value": "DNSError",
          "label": "DNSError",
          "description": "Return DNS errors for a selected domain"
        },
        {
          "value": "DNSRandom",
          "label": "DNSRandom",
          "description": "Return random DNS results for a selected domain"
        },
        {
          "value": "JVMCPUStress",
          "label": "JVMCPUStress",
          "description": "Stress CPU inside a JVM method"
        },
        {
          "value": "JVMException",
          "label": "JVMException",
          "description": "Throw a JVM exception"
        },
        {
          "value": "JVMGarbageCollector",
          "label": "JVMGarbageCollector",
          "description": "Trigger JVM garbage collection for the app"
        },
        {
          "value": "JVMLatency",
          "label": "JVMLatency",
          "description": "Inject latency into a JVM method"
        },
        {
          "value": "JVMMemoryStress",
          "label": "JVMMemoryStress",
          "description": "Stress memory inside a JVM method"
        },
        {
          "value": "JVMMySQLException",
          "label": "JVMMySQLException",
          "description": "Inject SQL exceptions into a MySQL operation"
        },
        {
          "value": "JVMMySQLLatency",
          "label": "JVMMySQLLatency",
          "description": "Inject latency into a MySQL operation"
        },
        {
          "value": "JVMReturn",
          "label": "JVMReturn",
          "description": "Override a JVM return value"
        },
        {
          "value": "JVMRuntimeMutator",
          "label": "JVMRuntimeMutator",
          "description": "Apply a runtime mutator strategy to a JVM method"
        },
        {
          "value": "MemoryStress",
          "label": "MemoryStress",
          "description": "Stress a container with memory pressure"
        },
        {
          "value": "NetworkBandwidth",
          "label": "NetworkBandwidth",
          "description": "Limit bandwidth to a downstream service"
        },
        {
          "value": "NetworkCorrupt",
          "label": "NetworkCorrupt",
          "description": "Corrupt traffic to a downstream service"
        },
        {
          "value": "NetworkDelay",
          "label": "NetworkDelay",
          "description": "Delay traffic to a downstream service"
        },
        {
          "value": "NetworkDuplicate",
          "label": "NetworkDuplicate",
          "description": "Duplicate traffic to a downstream service"
        },
        {
          "value": "NetworkLoss",
          "label": "NetworkLoss",
          "description": "Drop traffic to a downstream service"
        },
        {
          "value": "NetworkPartition",
          "label": "NetworkPartition",
          "description": "Partition traffic to a downstream service"
        },
        {
          "value": "PodFailure",
          "label": "PodFailure",
          "description": "Fail pods for the selected app"
        },
        {
          "value": "PodKill",
          "label": "PodKill",
          "description": "Kill pods for the selected app"
        },
        {
          "value": "TimeSkew",
          "label": "TimeSkew",
          "description": "Shift time in a specific container"
        }
      ]
    }
  ],
  "preview": {
    "resource_summary": {
      "containers": 1,
      "database_operations": 5,
      "dns_domains": 2,
      "jvm_methods": 4,
      "network_targets": 2,
      "runtime_mutator_methods": 1
    }
  },
  "can_apply": false
}
```

### Call 3

```bash
./chaos-exp.exe -output json --config %TEMP%\chaos-exp-examples-mutator.yaml --chaos-type JVMRuntimeMutator
```

```json
{
  "mode": "guided",
  "stage": "select_runtime_mutator_method",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-delivery-service",
    "chaos_type": "JVMRuntimeMutator"
  },
  "resolved": {
    "app": "ts-delivery-service",
    "chaos_type": "JVMRuntimeMutator",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "method_ref",
      "kind": "object_ref",
      "required": true,
      "description": "Select the method for runtime mutator injection",
      "options": [
        {
          "value": "delivery.mq.RabbitReceive#process",
          "label": "delivery.mq.RabbitReceive#process",
          "metadata": {
            "class": "delivery.mq.RabbitReceive",
            "method": "process"
          }
        }
      ],
      "key_fields": [
        "class",
        "method"
      ]
    }
  ],
  "can_apply": false
}
```

### Call 4

```bash
./chaos-exp.exe -output json --config %TEMP%\chaos-exp-examples-mutator.yaml --next delivery.mq.RabbitReceive#process
```

```json
{
  "mode": "guided",
  "stage": "select_runtime_mutator_config",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-delivery-service",
    "chaos_type": "JVMRuntimeMutator",
    "class": "delivery.mq.RabbitReceive",
    "method": "process"
  },
  "resolved": {
    "app": "ts-delivery-service",
    "chaos_type": "JVMRuntimeMutator",
    "class": "delivery.mq.RabbitReceive",
    "method": "process",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "mutator_config",
      "kind": "enum",
      "required": true,
      "description": "Select the runtime mutator config",
      "options": [
        {
          "value": "constant:\"[process][Receive delivery object][delivery object: {}]\":\"mutated_[process][Receive delivery object][delivery object: {}]\"",
          "label": "Mutate string constant to an error-oriented value",
          "metadata": {
            "description": "Mutate string constant to an error-oriented value",
            "mutation_from": "\"[process][Receive delivery object][delivery object: {}]\"",
            "mutation_strategy": "",
            "mutation_to": "\"mutated_[process][Receive delivery object][delivery object: {}]\"",
            "mutation_type_name": "constant"
          }
        },
        {
          "value": "operator:add_to_sub",
          "label": "Mutate + to -",
          "metadata": {
            "description": "Mutate + to -",
            "mutation_from": "",
            "mutation_strategy": "add_to_sub",
            "mutation_to": "",
            "mutation_type_name": "operator"
          }
        },
        {
          "value": "string:empty",
          "label": "Replace string result with empty string",
          "metadata": {
            "description": "Replace string result with empty string",
            "mutation_from": "",
            "mutation_strategy": "empty",
            "mutation_to": "",
            "mutation_type_name": "string"
          }
        },
        {
          "value": "string:lowercase",
          "label": "Convert string to lowercase",
          "metadata": {
            "description": "Convert string to lowercase",
            "mutation_from": "",
            "mutation_strategy": "lowercase",
            "mutation_to": "",
            "mutation_type_name": "string"
          }
        },
        {
          "value": "string:null",
          "label": "Replace string result with null",
          "metadata": {
            "description": "Replace string result with null",
            "mutation_from": "",
            "mutation_strategy": "null",
            "mutation_to": "",
            "mutation_type_name": "string"
          }
        },
        {
          "value": "string:random",
          "label": "Replace string with random content",
          "metadata": {
            "description": "Replace string with random content",
            "mutation_from": "",
            "mutation_strategy": "random",
            "mutation_to": "",
            "mutation_type_name": "string"
          }
        },
        {
          "value": "string:reverse",
          "label": "Reverse string content",
          "metadata": {
            "description": "Reverse string content",
            "mutation_from": "",
            "mutation_strategy": "reverse",
            "mutation_to": "",
            "mutation_type_name": "string"
          }
        },
        {
          "value": "string:uppercase",
          "label": "Convert string to uppercase",
          "metadata": {
            "description": "Convert string to uppercase",
            "mutation_from": "",
            "mutation_strategy": "uppercase",
            "mutation_to": "",
            "mutation_type_name": "string"
          }
        }
      ]
    },
    {
      "name": "duration",
      "kind": "number_range",
      "required": false,
      "description": "Fault duration in minutes",
      "min": 1,
      "max": 60,
      "step": 1,
      "default": 5,
      "unit": "minute"
    }
  ],
  "can_apply": false
}
```

### Call 5

```bash
./chaos-exp.exe -output json --config %TEMP%\chaos-exp-examples-mutator.yaml --next operator:add_to_sub
```

```json
{
  "mode": "guided",
  "stage": "ready_to_apply",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-delivery-service",
    "chaos_type": "JVMRuntimeMutator",
    "class": "delivery.mq.RabbitReceive",
    "method": "process",
    "mutator_config": "operator:add_to_sub",
    "duration": 5
  },
  "resolved": {
    "app": "ts-delivery-service",
    "chaos_type": "JVMRuntimeMutator",
    "class": "delivery.mq.RabbitReceive",
    "duration": 5,
    "method": "process",
    "mutator_config": "operator:add_to_sub",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "duration",
      "kind": "number_range",
      "required": false,
      "description": "Fault duration in minutes",
      "min": 1,
      "max": 60,
      "step": 1,
      "default": 5,
      "unit": "minute"
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "JVMRuntimeMutator": {
          "Duration": 5,
          "MutatorTargetIdx": 1824,
          "System": 7
        }
      },
      "chaos_type": "JVMRuntimeMutator",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-delivery-service",
        "class_name": "delivery.mq.RabbitReceive",
        "description": "Mutate + to -",
        "method_name": "process",
        "mutation_from": "",
        "mutation_strategy": "add_to_sub",
        "mutation_to": "",
        "mutation_type": 1,
        "mutation_type_name": "operator"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "function": [
        "delivery.mq.RabbitReceive.process"
      ],
      "service": [
        "ts-delivery-service"
      ]
    }
  },
  "apply_payload": {
    "JVMRuntimeMutator": {
      "Duration": 5,
      "MutatorTargetIdx": 1824,
      "System": 7
    }
  },
  "can_apply": true
}
```
