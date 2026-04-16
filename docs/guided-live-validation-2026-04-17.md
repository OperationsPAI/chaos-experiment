# Guided CLI Live Validation Record

- Date: 2026-04-17
- Namespace under validation: `ts`
- Validation mode: guided CLI only, stopped at `ready_to_apply`, no `--apply` executed
- Representative app for most records: `ts-auth-service`
- Representative HTTP endpoint: `GET /api/v1/verifycode/verify/*`
- Representative JVM method: `auth.security.jwt.JWTProvider#createToken`
- Representative downstream target / domain / database tuple: `mysql`, `mysql`, `ts/auth_user/SELECT`
- For parameterized types, the record now keeps the full chain: resource selection -> parameter-range response -> ready_to_apply

## PodKill

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type PodKill
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
          "AppIdx": 7,
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
      "AppIdx": 7,
      "Duration": 5,
      "System": 7
    }
  },
  "can_apply": true
}
```

## PodFailure

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type PodFailure
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
    "chaos_type": "PodFailure",
    "duration": 5
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "PodFailure",
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
        "PodFailure": {
          "AppIdx": 7,
          "Duration": 5,
          "System": 7
        }
      },
      "chaos_type": "PodFailure",
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
    "PodFailure": {
      "AppIdx": 7,
      "Duration": 5,
      "System": 7
    }
  },
  "can_apply": true
}
```

## ContainerKill

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type ContainerKill
```
```json
{
  "mode": "guided",
  "stage": "select_container",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "ContainerKill"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "ContainerKill",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "container",
      "kind": "enum",
      "required": true,
      "description": "Select a container under the app",
      "options": [
        {
          "value": "ts-auth-service",
          "label": "ts-auth-service"
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type ContainerKill --container ts-auth-service
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
    "chaos_type": "ContainerKill",
    "container": "ts-auth-service",
    "duration": 5
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "ContainerKill",
    "container": "ts-auth-service",
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
        "ContainerKill": {
          "ContainerIdx": 8,
          "Duration": 5,
          "System": 7
        }
      },
      "chaos_type": "ContainerKill",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "container_name": "ts-auth-service",
        "namespace": "ts"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "container": [
        "ts-auth-service"
      ],
      "service": [
        "ts-auth-service"
      ]
    }
  },
  "apply_payload": {
    "ContainerKill": {
      "ContainerIdx": 8,
      "Duration": 5,
      "System": 7
    }
  },
  "can_apply": true
}
```

## CPUStress

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type CPUStress
```
```json
{
  "mode": "guided",
  "stage": "select_container",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "CPUStress"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "CPUStress",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "container",
      "kind": "enum",
      "required": true,
      "description": "Select a container under the app",
      "options": [
        {
          "value": "ts-auth-service",
          "label": "ts-auth-service"
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type CPUStress --container ts-auth-service
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "CPUStress",
    "container": "ts-auth-service"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "CPUStress",
    "container": "ts-auth-service",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill CPU stress parameters",
      "fields": [
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
        },
        {
          "name": "cpu_load",
          "kind": "number_range",
          "required": true,
          "description": "CPU load percentage",
          "min": 1,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "cpu_worker",
          "kind": "number_range",
          "required": true,
          "description": "CPU stress worker count",
          "min": 1,
          "max": 3,
          "step": 1
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "cpu_load and cpu_worker are required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type CPUStress --container ts-auth-service --cpu-load 80 --cpu-worker 1
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
    "chaos_type": "CPUStress",
    "container": "ts-auth-service",
    "duration": 5,
    "cpu_load": 80,
    "cpu_worker": 1
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "CPUStress",
    "container": "ts-auth-service",
    "cpu_load": 80,
    "cpu_worker": 1,
    "duration": 5,
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill CPU stress parameters",
      "fields": [
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
        },
        {
          "name": "cpu_load",
          "kind": "number_range",
          "required": true,
          "description": "CPU load percentage",
          "min": 1,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "cpu_worker",
          "kind": "number_range",
          "required": true,
          "description": "CPU stress worker count",
          "min": 1,
          "max": 3,
          "step": 1
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "CPUStress": {
          "CPULoad": 80,
          "CPUWorker": 1,
          "ContainerIdx": 8,
          "Duration": 5,
          "System": 7
        }
      },
      "chaos_type": "CPUStress",
      "cpu_load": 80,
      "cpu_worker": 1,
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "container_name": "ts-auth-service",
        "namespace": "ts"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "container": [
        "ts-auth-service"
      ],
      "metric": [
        "cpu"
      ],
      "service": [
        "ts-auth-service"
      ]
    }
  },
  "apply_payload": {
    "CPUStress": {
      "CPULoad": 80,
      "CPUWorker": 1,
      "ContainerIdx": 8,
      "Duration": 5,
      "System": 7
    }
  },
  "can_apply": true
}
```

## MemoryStress

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type MemoryStress
```
```json
{
  "mode": "guided",
  "stage": "select_container",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "MemoryStress"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "MemoryStress",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "container",
      "kind": "enum",
      "required": true,
      "description": "Select a container under the app",
      "options": [
        {
          "value": "ts-auth-service",
          "label": "ts-auth-service"
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type MemoryStress --container ts-auth-service
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "MemoryStress",
    "container": "ts-auth-service"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "MemoryStress",
    "container": "ts-auth-service",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill memory stress parameters",
      "fields": [
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
        },
        {
          "name": "memory_size",
          "kind": "number_range",
          "required": true,
          "description": "Memory size",
          "min": 1,
          "max": 1024,
          "step": 1,
          "unit": "MiB"
        },
        {
          "name": "mem_worker",
          "kind": "number_range",
          "required": true,
          "description": "Memory stress worker count",
          "min": 1,
          "max": 4,
          "step": 1
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "memory_size and mem_worker are required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type MemoryStress --container ts-auth-service --memory-size 128 --mem-worker 1
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
    "chaos_type": "MemoryStress",
    "container": "ts-auth-service",
    "duration": 5,
    "memory_size": 128,
    "mem_worker": 1
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "MemoryStress",
    "container": "ts-auth-service",
    "duration": 5,
    "mem_worker": 1,
    "memory_size": 128,
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill memory stress parameters",
      "fields": [
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
        },
        {
          "name": "memory_size",
          "kind": "number_range",
          "required": true,
          "description": "Memory size",
          "min": 1,
          "max": 1024,
          "step": 1,
          "unit": "MiB"
        },
        {
          "name": "mem_worker",
          "kind": "number_range",
          "required": true,
          "description": "Memory stress worker count",
          "min": 1,
          "max": 4,
          "step": 1
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "MemoryStress": {
          "ContainerIdx": 8,
          "Duration": 5,
          "MemWorker": 1,
          "MemorySize": 128,
          "System": 7
        }
      },
      "chaos_type": "MemoryStress",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "container_name": "ts-auth-service",
        "namespace": "ts"
      },
      "mem_worker": 1,
      "memory_size": 128,
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "container": [
        "ts-auth-service"
      ],
      "metric": [
        "memory"
      ],
      "service": [
        "ts-auth-service"
      ]
    }
  },
  "apply_payload": {
    "MemoryStress": {
      "ContainerIdx": 8,
      "Duration": 5,
      "MemWorker": 1,
      "MemorySize": 128,
      "System": 7
    }
  },
  "can_apply": true
}
```

## TimeSkew

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type TimeSkew
```
```json
{
  "mode": "guided",
  "stage": "select_container",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "TimeSkew"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "TimeSkew",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "container",
      "kind": "enum",
      "required": true,
      "description": "Select a container under the app",
      "options": [
        {
          "value": "ts-auth-service",
          "label": "ts-auth-service"
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type TimeSkew --container ts-auth-service
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "TimeSkew",
    "container": "ts-auth-service"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "TimeSkew",
    "container": "ts-auth-service",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill time skew parameters",
      "fields": [
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
        },
        {
          "name": "time_offset",
          "kind": "number_range",
          "required": true,
          "description": "Time offset in seconds",
          "min": -600,
          "max": 600,
          "step": 1,
          "unit": "second"
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "time_offset is required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type TimeSkew --container ts-auth-service --time-offset 60
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
    "chaos_type": "TimeSkew",
    "container": "ts-auth-service",
    "duration": 5,
    "time_offset": 60
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "TimeSkew",
    "container": "ts-auth-service",
    "duration": 5,
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts",
    "time_offset": 60
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill time skew parameters",
      "fields": [
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
        },
        {
          "name": "time_offset",
          "kind": "number_range",
          "required": true,
          "description": "Time offset in seconds",
          "min": -600,
          "max": 600,
          "step": 1,
          "unit": "second"
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "TimeSkew": {
          "ContainerIdx": 8,
          "Duration": 5,
          "System": 7,
          "TimeOffset": 60
        }
      },
      "chaos_type": "TimeSkew",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "container_name": "ts-auth-service",
        "namespace": "ts"
      },
      "namespace": "ts",
      "system": "ts",
      "time_offset": 60
    },
    "groundtruth": {
      "container": [
        "ts-auth-service"
      ],
      "service": [
        "ts-auth-service"
      ]
    }
  },
  "apply_payload": {
    "TimeSkew": {
      "ContainerIdx": 8,
      "Duration": 5,
      "System": 7,
      "TimeOffset": 60
    }
  },
  "can_apply": true
}
```

## HTTPRequestAbort

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPRequestAbort
```
```json
{
  "mode": "guided",
  "stage": "select_http_endpoint",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestAbort"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestAbort",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "endpoint",
      "kind": "object_ref",
      "required": true,
      "description": "Select the HTTP endpoint for request abort",
      "options": [
        {
          "value": "GET /api/v1/verifycode/verify/*",
          "label": "GET /api/v1/verifycode/verify/*",
          "metadata": {
            "http_method": "GET",
            "route": "/api/v1/verifycode/verify/*",
            "span_name": "GET",
            "target_service": "ts-verification-code-service"
          }
        }
      ],
      "key_fields": [
        "http_method",
        "route"
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPRequestAbort --route '/api/v1/verifycode/verify/*' --http-method GET
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
    "chaos_type": "HTTPRequestAbort",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET",
    "duration": 5
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestAbort",
    "duration": 5,
    "http_method": "GET",
    "namespace": "ts",
    "route": "/api/v1/verifycode/verify/*",
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
        "HTTPRequestAbort": {
          "Duration": 5,
          "EndpointIdx": 940,
          "System": 7
        }
      },
      "chaos_type": "HTTPRequestAbort",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "http_method": "GET",
        "route": "/api/v1/verifycode/verify/*",
        "server_address": "ts-verification-code-service",
        "server_port": "8080",
        "span_name": "GET"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "http_latency"
      ],
      "service": [
        "ts-auth-service",
        "ts-verification-code-service"
      ],
      "span": [
        "GET"
      ]
    }
  },
  "apply_payload": {
    "HTTPRequestAbort": {
      "Duration": 5,
      "EndpointIdx": 940,
      "System": 7
    }
  },
  "can_apply": true
}
```

## HTTPResponseAbort

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPResponseAbort
```
```json
{
  "mode": "guided",
  "stage": "select_http_endpoint",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseAbort"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseAbort",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "endpoint",
      "kind": "object_ref",
      "required": true,
      "description": "Select the HTTP endpoint for response abort",
      "options": [
        {
          "value": "GET /api/v1/verifycode/verify/*",
          "label": "GET /api/v1/verifycode/verify/*",
          "metadata": {
            "http_method": "GET",
            "route": "/api/v1/verifycode/verify/*",
            "span_name": "GET",
            "target_service": "ts-verification-code-service"
          }
        }
      ],
      "key_fields": [
        "http_method",
        "route"
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPResponseAbort --route '/api/v1/verifycode/verify/*' --http-method GET
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
    "chaos_type": "HTTPResponseAbort",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET",
    "duration": 5
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseAbort",
    "duration": 5,
    "http_method": "GET",
    "namespace": "ts",
    "route": "/api/v1/verifycode/verify/*",
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
        "HTTPResponseAbort": {
          "Duration": 5,
          "EndpointIdx": 940,
          "System": 7
        }
      },
      "chaos_type": "HTTPResponseAbort",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "http_method": "GET",
        "route": "/api/v1/verifycode/verify/*",
        "server_address": "ts-verification-code-service",
        "server_port": "8080",
        "span_name": "GET"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "http_latency"
      ],
      "service": [
        "ts-auth-service",
        "ts-verification-code-service"
      ],
      "span": [
        "GET"
      ]
    }
  },
  "apply_payload": {
    "HTTPResponseAbort": {
      "Duration": 5,
      "EndpointIdx": 940,
      "System": 7
    }
  },
  "can_apply": true
}
```

## HTTPRequestDelay

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPRequestDelay
```
```json
{
  "mode": "guided",
  "stage": "select_http_endpoint",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestDelay"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestDelay",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "endpoint",
      "kind": "object_ref",
      "required": true,
      "description": "Select the HTTP endpoint for request delay",
      "options": [
        {
          "value": "GET /api/v1/verifycode/verify/*",
          "label": "GET /api/v1/verifycode/verify/*",
          "metadata": {
            "http_method": "GET",
            "route": "/api/v1/verifycode/verify/*",
            "span_name": "GET",
            "target_service": "ts-verification-code-service"
          }
        }
      ],
      "key_fields": [
        "http_method",
        "route"
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPRequestDelay --route '/api/v1/verifycode/verify/*' --http-method GET
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestDelay",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestDelay",
    "http_method": "GET",
    "namespace": "ts",
    "route": "/api/v1/verifycode/verify/*",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill HTTP request delay parameters",
      "fields": [
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
        },
        {
          "name": "delay_duration",
          "kind": "number_range",
          "required": true,
          "description": "Request delay duration",
          "min": 10,
          "max": 5000,
          "step": 1,
          "unit": "ms"
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "delay_duration is required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPRequestDelay --route '/api/v1/verifycode/verify/*' --http-method GET --delay-duration 200
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
    "chaos_type": "HTTPRequestDelay",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET",
    "duration": 5,
    "delay_duration": 200
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestDelay",
    "delay_duration": 200,
    "duration": 5,
    "http_method": "GET",
    "namespace": "ts",
    "route": "/api/v1/verifycode/verify/*",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill HTTP request delay parameters",
      "fields": [
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
        },
        {
          "name": "delay_duration",
          "kind": "number_range",
          "required": true,
          "description": "Request delay duration",
          "min": 10,
          "max": 5000,
          "step": 1,
          "unit": "ms"
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "HTTPRequestDelay": {
          "DelayDuration": 200,
          "Duration": 5,
          "EndpointIdx": 940,
          "System": 7
        }
      },
      "chaos_type": "HTTPRequestDelay",
      "delay_duration": 200,
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "http_method": "GET",
        "route": "/api/v1/verifycode/verify/*",
        "server_address": "ts-verification-code-service",
        "server_port": "8080",
        "span_name": "GET"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "http_latency"
      ],
      "service": [
        "ts-auth-service",
        "ts-verification-code-service"
      ],
      "span": [
        "GET"
      ]
    }
  },
  "apply_payload": {
    "HTTPRequestDelay": {
      "DelayDuration": 200,
      "Duration": 5,
      "EndpointIdx": 940,
      "System": 7
    }
  },
  "can_apply": true
}
```

## HTTPResponseDelay

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPResponseDelay
```
```json
{
  "mode": "guided",
  "stage": "select_http_endpoint",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseDelay"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseDelay",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "endpoint",
      "kind": "object_ref",
      "required": true,
      "description": "Select the HTTP endpoint for response delay",
      "options": [
        {
          "value": "GET /api/v1/verifycode/verify/*",
          "label": "GET /api/v1/verifycode/verify/*",
          "metadata": {
            "http_method": "GET",
            "route": "/api/v1/verifycode/verify/*",
            "span_name": "GET",
            "target_service": "ts-verification-code-service"
          }
        }
      ],
      "key_fields": [
        "http_method",
        "route"
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPResponseDelay --route '/api/v1/verifycode/verify/*' --http-method GET
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseDelay",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseDelay",
    "http_method": "GET",
    "namespace": "ts",
    "route": "/api/v1/verifycode/verify/*",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill HTTP response delay parameters",
      "fields": [
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
        },
        {
          "name": "delay_duration",
          "kind": "number_range",
          "required": true,
          "description": "Response delay duration",
          "min": 10,
          "max": 5000,
          "step": 1,
          "unit": "ms"
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "delay_duration is required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPResponseDelay --route '/api/v1/verifycode/verify/*' --http-method GET --delay-duration 200
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
    "chaos_type": "HTTPResponseDelay",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET",
    "duration": 5,
    "delay_duration": 200
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseDelay",
    "delay_duration": 200,
    "duration": 5,
    "http_method": "GET",
    "namespace": "ts",
    "route": "/api/v1/verifycode/verify/*",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill HTTP response delay parameters",
      "fields": [
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
        },
        {
          "name": "delay_duration",
          "kind": "number_range",
          "required": true,
          "description": "Response delay duration",
          "min": 10,
          "max": 5000,
          "step": 1,
          "unit": "ms"
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "HTTPResponseDelay": {
          "DelayDuration": 200,
          "Duration": 5,
          "EndpointIdx": 940,
          "System": 7
        }
      },
      "chaos_type": "HTTPResponseDelay",
      "delay_duration": 200,
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "http_method": "GET",
        "route": "/api/v1/verifycode/verify/*",
        "server_address": "ts-verification-code-service",
        "server_port": "8080",
        "span_name": "GET"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "http_latency"
      ],
      "service": [
        "ts-auth-service",
        "ts-verification-code-service"
      ],
      "span": [
        "GET"
      ]
    }
  },
  "apply_payload": {
    "HTTPResponseDelay": {
      "DelayDuration": 200,
      "Duration": 5,
      "EndpointIdx": 940,
      "System": 7
    }
  },
  "can_apply": true
}
```

## HTTPResponseReplaceBody

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPResponseReplaceBody
```
```json
{
  "mode": "guided",
  "stage": "select_http_endpoint",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseReplaceBody"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseReplaceBody",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "endpoint",
      "kind": "object_ref",
      "required": true,
      "description": "Select the HTTP endpoint for response body replacement",
      "options": [
        {
          "value": "GET /api/v1/verifycode/verify/*",
          "label": "GET /api/v1/verifycode/verify/*",
          "metadata": {
            "http_method": "GET",
            "route": "/api/v1/verifycode/verify/*",
            "span_name": "GET",
            "target_service": "ts-verification-code-service"
          }
        }
      ],
      "key_fields": [
        "http_method",
        "route"
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPResponseReplaceBody --route '/api/v1/verifycode/verify/*' --http-method GET
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseReplaceBody",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseReplaceBody",
    "http_method": "GET",
    "namespace": "ts",
    "route": "/api/v1/verifycode/verify/*",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill HTTP response replacement parameters",
      "fields": [
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
        },
        {
          "name": "body_type",
          "kind": "enum",
          "required": true,
          "description": "Replacement body type",
          "options": [
            {
              "value": "empty",
              "label": "empty"
            },
            {
              "value": "random",
              "label": "random"
            }
          ]
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "body_type is required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPResponseReplaceBody --route '/api/v1/verifycode/verify/*' --http-method GET --body-type empty
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
    "chaos_type": "HTTPResponseReplaceBody",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET",
    "duration": 5,
    "body_type": "empty"
  },
  "resolved": {
    "app": "ts-auth-service",
    "body_type": "empty",
    "chaos_type": "HTTPResponseReplaceBody",
    "duration": 5,
    "http_method": "GET",
    "namespace": "ts",
    "route": "/api/v1/verifycode/verify/*",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill HTTP response replacement parameters",
      "fields": [
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
        },
        {
          "name": "body_type",
          "kind": "enum",
          "required": true,
          "description": "Replacement body type",
          "options": [
            {
              "value": "empty",
              "label": "empty"
            },
            {
              "value": "random",
              "label": "random"
            }
          ]
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "HTTPResponseReplaceBody": {
          "BodyType": 0,
          "Duration": 5,
          "EndpointIdx": 940,
          "System": 7
        }
      },
      "body_type": "empty",
      "chaos_type": "HTTPResponseReplaceBody",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "http_method": "GET",
        "route": "/api/v1/verifycode/verify/*",
        "server_address": "ts-verification-code-service",
        "server_port": "8080",
        "span_name": "GET"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "http_latency"
      ],
      "service": [
        "ts-auth-service",
        "ts-verification-code-service"
      ],
      "span": [
        "GET"
      ]
    }
  },
  "apply_payload": {
    "HTTPResponseReplaceBody": {
      "BodyType": 0,
      "Duration": 5,
      "EndpointIdx": 940,
      "System": 7
    }
  },
  "can_apply": true
}
```

## HTTPResponsePatchBody

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPResponsePatchBody
```
```json
{
  "mode": "guided",
  "stage": "select_http_endpoint",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponsePatchBody"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponsePatchBody",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "endpoint",
      "kind": "object_ref",
      "required": true,
      "description": "Select the HTTP endpoint for response body patching",
      "options": [
        {
          "value": "GET /api/v1/verifycode/verify/*",
          "label": "GET /api/v1/verifycode/verify/*",
          "metadata": {
            "http_method": "GET",
            "route": "/api/v1/verifycode/verify/*",
            "span_name": "GET",
            "target_service": "ts-verification-code-service"
          }
        }
      ],
      "key_fields": [
        "http_method",
        "route"
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPResponsePatchBody --route '/api/v1/verifycode/verify/*' --http-method GET
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
    "chaos_type": "HTTPResponsePatchBody",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET",
    "duration": 5
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponsePatchBody",
    "duration": 5,
    "http_method": "GET",
    "namespace": "ts",
    "route": "/api/v1/verifycode/verify/*",
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
        "HTTPResponsePatchBody": {
          "Duration": 5,
          "EndpointIdx": 940,
          "System": 7
        }
      },
      "chaos_type": "HTTPResponsePatchBody",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "http_method": "GET",
        "route": "/api/v1/verifycode/verify/*",
        "server_address": "ts-verification-code-service",
        "server_port": "8080",
        "span_name": "GET"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "http_latency"
      ],
      "service": [
        "ts-auth-service",
        "ts-verification-code-service"
      ],
      "span": [
        "GET"
      ]
    }
  },
  "apply_payload": {
    "HTTPResponsePatchBody": {
      "Duration": 5,
      "EndpointIdx": 940,
      "System": 7
    }
  },
  "can_apply": true
}
```

## HTTPRequestReplacePath

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPRequestReplacePath
```
```json
{
  "mode": "guided",
  "stage": "select_http_endpoint",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestReplacePath"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestReplacePath",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "endpoint",
      "kind": "object_ref",
      "required": true,
      "description": "Select the HTTP endpoint for request path replacement",
      "options": [
        {
          "value": "GET /api/v1/verifycode/verify/*",
          "label": "GET /api/v1/verifycode/verify/*",
          "metadata": {
            "http_method": "GET",
            "route": "/api/v1/verifycode/verify/*",
            "span_name": "GET",
            "target_service": "ts-verification-code-service"
          }
        }
      ],
      "key_fields": [
        "http_method",
        "route"
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPRequestReplacePath --route '/api/v1/verifycode/verify/*' --http-method GET
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
    "chaos_type": "HTTPRequestReplacePath",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET",
    "duration": 5
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestReplacePath",
    "duration": 5,
    "http_method": "GET",
    "namespace": "ts",
    "route": "/api/v1/verifycode/verify/*",
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
        "HTTPRequestReplacePath": {
          "Duration": 5,
          "EndpointIdx": 940,
          "System": 7
        }
      },
      "chaos_type": "HTTPRequestReplacePath",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "http_method": "GET",
        "route": "/api/v1/verifycode/verify/*",
        "server_address": "ts-verification-code-service",
        "server_port": "8080",
        "span_name": "GET"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "http_latency"
      ],
      "service": [
        "ts-auth-service",
        "ts-verification-code-service"
      ],
      "span": [
        "GET"
      ]
    }
  },
  "apply_payload": {
    "HTTPRequestReplacePath": {
      "Duration": 5,
      "EndpointIdx": 940,
      "System": 7
    }
  },
  "can_apply": true
}
```

## HTTPRequestReplaceMethod

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPRequestReplaceMethod
```
```json
{
  "mode": "guided",
  "stage": "select_http_endpoint",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestReplaceMethod"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestReplaceMethod",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "endpoint",
      "kind": "object_ref",
      "required": true,
      "description": "Select the HTTP endpoint for request method replacement",
      "options": [
        {
          "value": "GET /api/v1/verifycode/verify/*",
          "label": "GET /api/v1/verifycode/verify/*",
          "metadata": {
            "http_method": "GET",
            "route": "/api/v1/verifycode/verify/*",
            "span_name": "GET",
            "target_service": "ts-verification-code-service"
          }
        }
      ],
      "key_fields": [
        "http_method",
        "route"
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPRequestReplaceMethod --route '/api/v1/verifycode/verify/*' --http-method GET
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestReplaceMethod",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestReplaceMethod",
    "http_method": "GET",
    "namespace": "ts",
    "route": "/api/v1/verifycode/verify/*",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill HTTP request replacement parameters",
      "fields": [
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
        },
        {
          "name": "replace_method",
          "kind": "enum",
          "required": true,
          "description": "Replacement HTTP method",
          "options": [
            {
              "value": "POST",
              "label": "POST"
            },
            {
              "value": "PUT",
              "label": "PUT"
            },
            {
              "value": "DELETE",
              "label": "DELETE"
            },
            {
              "value": "HEAD",
              "label": "HEAD"
            },
            {
              "value": "OPTIONS",
              "label": "OPTIONS"
            },
            {
              "value": "PATCH",
              "label": "PATCH"
            }
          ]
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "replace_method is required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPRequestReplaceMethod --route '/api/v1/verifycode/verify/*' --http-method GET --replace-method POST
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
    "chaos_type": "HTTPRequestReplaceMethod",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET",
    "duration": 5,
    "replace_method": "POST"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPRequestReplaceMethod",
    "duration": 5,
    "http_method": "GET",
    "namespace": "ts",
    "replace_method": "POST",
    "route": "/api/v1/verifycode/verify/*",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill HTTP request replacement parameters",
      "fields": [
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
        },
        {
          "name": "replace_method",
          "kind": "enum",
          "required": true,
          "description": "Replacement HTTP method",
          "options": [
            {
              "value": "POST",
              "label": "POST"
            },
            {
              "value": "PUT",
              "label": "PUT"
            },
            {
              "value": "DELETE",
              "label": "DELETE"
            },
            {
              "value": "HEAD",
              "label": "HEAD"
            },
            {
              "value": "OPTIONS",
              "label": "OPTIONS"
            },
            {
              "value": "PATCH",
              "label": "PATCH"
            }
          ]
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "HTTPRequestReplaceMethod": {
          "Duration": 5,
          "EndpointIdx": 940,
          "ReplaceMethod": 0,
          "System": 7
        }
      },
      "chaos_type": "HTTPRequestReplaceMethod",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "http_method": "GET",
        "route": "/api/v1/verifycode/verify/*",
        "server_address": "ts-verification-code-service",
        "server_port": "8080",
        "span_name": "GET"
      },
      "namespace": "ts",
      "replace_method": "POST",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "http_latency"
      ],
      "service": [
        "ts-auth-service",
        "ts-verification-code-service"
      ],
      "span": [
        "GET"
      ]
    }
  },
  "apply_payload": {
    "HTTPRequestReplaceMethod": {
      "Duration": 5,
      "EndpointIdx": 940,
      "ReplaceMethod": 0,
      "System": 7
    }
  },
  "can_apply": true
}
```

## HTTPResponseReplaceCode

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPResponseReplaceCode
```
```json
{
  "mode": "guided",
  "stage": "select_http_endpoint",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseReplaceCode"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseReplaceCode",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "endpoint",
      "kind": "object_ref",
      "required": true,
      "description": "Select the HTTP endpoint for response code replacement",
      "options": [
        {
          "value": "GET /api/v1/verifycode/verify/*",
          "label": "GET /api/v1/verifycode/verify/*",
          "metadata": {
            "http_method": "GET",
            "route": "/api/v1/verifycode/verify/*",
            "span_name": "GET",
            "target_service": "ts-verification-code-service"
          }
        }
      ],
      "key_fields": [
        "http_method",
        "route"
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPResponseReplaceCode --route '/api/v1/verifycode/verify/*' --http-method GET
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseReplaceCode",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseReplaceCode",
    "http_method": "GET",
    "namespace": "ts",
    "route": "/api/v1/verifycode/verify/*",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill HTTP response code replacement parameters",
      "fields": [
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
        },
        {
          "name": "status_code",
          "kind": "enum",
          "required": true,
          "description": "Replacement HTTP status code",
          "options": [
            {
              "value": "400",
              "label": "400"
            },
            {
              "value": "401",
              "label": "401"
            },
            {
              "value": "403",
              "label": "403"
            },
            {
              "value": "404",
              "label": "404"
            },
            {
              "value": "405",
              "label": "405"
            },
            {
              "value": "408",
              "label": "408"
            },
            {
              "value": "500",
              "label": "500"
            },
            {
              "value": "502",
              "label": "502"
            },
            {
              "value": "503",
              "label": "503"
            },
            {
              "value": "504",
              "label": "504"
            }
          ]
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "status_code is required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type HTTPResponseReplaceCode --route '/api/v1/verifycode/verify/*' --http-method GET --status-code 503
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
    "chaos_type": "HTTPResponseReplaceCode",
    "route": "/api/v1/verifycode/verify/*",
    "http_method": "GET",
    "duration": 5,
    "status_code": 503
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "HTTPResponseReplaceCode",
    "duration": 5,
    "http_method": "GET",
    "namespace": "ts",
    "route": "/api/v1/verifycode/verify/*",
    "status_code": 503,
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill HTTP response code replacement parameters",
      "fields": [
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
        },
        {
          "name": "status_code",
          "kind": "enum",
          "required": true,
          "description": "Replacement HTTP status code",
          "options": [
            {
              "value": "400",
              "label": "400"
            },
            {
              "value": "401",
              "label": "401"
            },
            {
              "value": "403",
              "label": "403"
            },
            {
              "value": "404",
              "label": "404"
            },
            {
              "value": "405",
              "label": "405"
            },
            {
              "value": "408",
              "label": "408"
            },
            {
              "value": "500",
              "label": "500"
            },
            {
              "value": "502",
              "label": "502"
            },
            {
              "value": "503",
              "label": "503"
            },
            {
              "value": "504",
              "label": "504"
            }
          ]
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "HTTPResponseReplaceCode": {
          "Duration": 5,
          "EndpointIdx": 940,
          "StatusCode": 8,
          "System": 7
        }
      },
      "chaos_type": "HTTPResponseReplaceCode",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "http_method": "GET",
        "route": "/api/v1/verifycode/verify/*",
        "server_address": "ts-verification-code-service",
        "server_port": "8080",
        "span_name": "GET"
      },
      "namespace": "ts",
      "status_code": 503,
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "http_latency"
      ],
      "service": [
        "ts-auth-service",
        "ts-verification-code-service"
      ],
      "span": [
        "GET"
      ]
    }
  },
  "apply_payload": {
    "HTTPResponseReplaceCode": {
      "Duration": 5,
      "EndpointIdx": 940,
      "StatusCode": 8,
      "System": 7
    }
  },
  "can_apply": true
}
```

## DNSError

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type DNSError
```
```json
{
  "mode": "guided",
  "stage": "select_dns_domain",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "DNSError"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "DNSError",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "domain",
      "kind": "enum",
      "required": true,
      "description": "Select the domain for DNS chaos",
      "options": [
        {
          "value": "mysql",
          "label": "mysql",
          "metadata": {
            "domain": "mysql",
            "span_names": [
              "ALTER table ts",
              "CREATE TABLE `ts`.`user_roles`",
              "CREATE table ts",
              "DELETE ts.auth_user",
              "DELETE ts.user_roles",
              "INSERT ts.auth_user",
              "INSERT ts.user_roles",
              "SELECT `ts`.`auth_user`",
              "SELECT ts",
              "SELECT ts.auth_user",
              "SELECT ts.ts",
              "SELECT ts.user_roles"
            ]
          }
        },
        {
          "value": "ts-verification-code-service",
          "label": "ts-verification-code-service",
          "metadata": {
            "domain": "ts-verification-code-service",
            "span_names": [
              "GET"
            ]
          }
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type DNSError --domain mysql
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
    "chaos_type": "DNSError",
    "domain": "mysql",
    "duration": 5
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "DNSError",
    "domain": "mysql",
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
        "DNSError": {
          "DNSEndpointIdx": 17,
          "Duration": 5,
          "System": 7
        }
      },
      "chaos_type": "DNSError",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "domain": "mysql",
        "span_names": [
          "ALTER table ts",
          "CREATE TABLE `ts`.`user_roles`",
          "CREATE table ts",
          "DELETE ts.auth_user",
          "DELETE ts.user_roles",
          "INSERT ts.auth_user",
          "INSERT ts.user_roles",
          "SELECT `ts`.`auth_user`",
          "SELECT ts",
          "SELECT ts.auth_user",
          "SELECT ts.ts",
          "SELECT ts.user_roles"
        ]
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "service": [
        "ts-auth-service",
        "mysql"
      ],
      "span": [
        "ALTER table ts",
        "CREATE TABLE `ts`.`user_roles`",
        "CREATE table ts",
        "DELETE ts.auth_user",
        "DELETE ts.user_roles",
        "INSERT ts.auth_user",
        "INSERT ts.user_roles",
        "SELECT `ts`.`auth_user`",
        "SELECT ts",
        "SELECT ts.auth_user",
        "SELECT ts.ts",
        "SELECT ts.user_roles"
      ]
    }
  },
  "apply_payload": {
    "DNSError": {
      "DNSEndpointIdx": 17,
      "Duration": 5,
      "System": 7
    }
  },
  "can_apply": true
}
```

## DNSRandom

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type DNSRandom
```
```json
{
  "mode": "guided",
  "stage": "select_dns_domain",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "DNSRandom"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "DNSRandom",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "domain",
      "kind": "enum",
      "required": true,
      "description": "Select the domain for DNS chaos",
      "options": [
        {
          "value": "mysql",
          "label": "mysql",
          "metadata": {
            "domain": "mysql",
            "span_names": [
              "ALTER table ts",
              "CREATE TABLE `ts`.`user_roles`",
              "CREATE table ts",
              "DELETE ts.auth_user",
              "DELETE ts.user_roles",
              "INSERT ts.auth_user",
              "INSERT ts.user_roles",
              "SELECT `ts`.`auth_user`",
              "SELECT ts",
              "SELECT ts.auth_user",
              "SELECT ts.ts",
              "SELECT ts.user_roles"
            ]
          }
        },
        {
          "value": "ts-verification-code-service",
          "label": "ts-verification-code-service",
          "metadata": {
            "domain": "ts-verification-code-service",
            "span_names": [
              "GET"
            ]
          }
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type DNSRandom --domain mysql
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
    "chaos_type": "DNSRandom",
    "domain": "mysql",
    "duration": 5
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "DNSRandom",
    "domain": "mysql",
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
        "DNSRandom": {
          "DNSEndpointIdx": 17,
          "Duration": 5,
          "System": 7
        }
      },
      "chaos_type": "DNSRandom",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "domain": "mysql",
        "span_names": [
          "ALTER table ts",
          "CREATE TABLE `ts`.`user_roles`",
          "CREATE table ts",
          "DELETE ts.auth_user",
          "DELETE ts.user_roles",
          "INSERT ts.auth_user",
          "INSERT ts.user_roles",
          "SELECT `ts`.`auth_user`",
          "SELECT ts",
          "SELECT ts.auth_user",
          "SELECT ts.ts",
          "SELECT ts.user_roles"
        ]
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "service": [
        "ts-auth-service",
        "mysql"
      ],
      "span": [
        "ALTER table ts",
        "CREATE TABLE `ts`.`user_roles`",
        "CREATE table ts",
        "DELETE ts.auth_user",
        "DELETE ts.user_roles",
        "INSERT ts.auth_user",
        "INSERT ts.user_roles",
        "SELECT `ts`.`auth_user`",
        "SELECT ts",
        "SELECT ts.auth_user",
        "SELECT ts.ts",
        "SELECT ts.user_roles"
      ]
    }
  },
  "apply_payload": {
    "DNSRandom": {
      "DNSEndpointIdx": 17,
      "Duration": 5,
      "System": 7
    }
  },
  "can_apply": true
}
```

## NetworkPartition

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkPartition
```
```json
{
  "mode": "guided",
  "stage": "select_network_target",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "NetworkPartition"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkPartition",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "target_service",
      "kind": "enum",
      "required": true,
      "description": "Select the target service for the network pair",
      "options": [
        {
          "value": "mysql",
          "label": "mysql"
        },
        {
          "value": "ts-verification-code-service",
          "label": "ts-verification-code-service"
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkPartition --target-service mysql
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "NetworkPartition",
    "target_service": "mysql"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkPartition",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts",
    "target_service": "mysql"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill network partition parameters",
      "fields": [
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
        },
        {
          "name": "direction",
          "kind": "enum",
          "required": true,
          "description": "Traffic direction",
          "options": [
            {
              "value": "to",
              "label": "to"
            },
            {
              "value": "from",
              "label": "from"
            },
            {
              "value": "both",
              "label": "both"
            }
          ]
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "direction is required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkPartition --target-service mysql --direction both
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
    "chaos_type": "NetworkPartition",
    "target_service": "mysql",
    "duration": 5,
    "direction": "both"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkPartition",
    "direction": "both",
    "duration": 5,
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts",
    "target_service": "mysql"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill network partition parameters",
      "fields": [
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
        },
        {
          "name": "direction",
          "kind": "enum",
          "required": true,
          "description": "Traffic direction",
          "options": [
            {
              "value": "to",
              "label": "to"
            },
            {
              "value": "from",
              "label": "from"
            },
            {
              "value": "both",
              "label": "both"
            }
          ]
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "NetworkPartition": {
          "Direction": 3,
          "Duration": 5,
          "NetworkPairIdx": 17,
          "System": 7
        }
      },
      "chaos_type": "NetworkPartition",
      "direction": "both",
      "duration": 5,
      "injection_point": {
        "source_service": "ts-auth-service",
        "span_names": [
          "ALTER table ts",
          "CREATE TABLE `ts`.`user_roles`",
          "CREATE table ts",
          "DELETE ts.auth_user",
          "DELETE ts.user_roles",
          "INSERT ts.auth_user",
          "INSERT ts.user_roles",
          "SELECT `ts`.`auth_user`",
          "SELECT ts",
          "SELECT ts.auth_user",
          "SELECT ts.ts",
          "SELECT ts.user_roles"
        ],
        "target_service": "mysql"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "network_latency"
      ],
      "service": [
        "ts-auth-service",
        "mysql"
      ],
      "span": [
        "ALTER table ts",
        "CREATE TABLE `ts`.`user_roles`",
        "CREATE table ts",
        "DELETE ts.auth_user",
        "DELETE ts.user_roles",
        "INSERT ts.auth_user",
        "INSERT ts.user_roles",
        "SELECT `ts`.`auth_user`",
        "SELECT ts",
        "SELECT ts.auth_user",
        "SELECT ts.ts",
        "SELECT ts.user_roles"
      ]
    }
  },
  "apply_payload": {
    "NetworkPartition": {
      "Direction": 3,
      "Duration": 5,
      "NetworkPairIdx": 17,
      "System": 7
    }
  },
  "can_apply": true
}
```

## NetworkDelay

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkDelay
```
```json
{
  "mode": "guided",
  "stage": "select_network_target",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "NetworkDelay"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkDelay",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "target_service",
      "kind": "enum",
      "required": true,
      "description": "Select the target service for the network pair",
      "options": [
        {
          "value": "mysql",
          "label": "mysql"
        },
        {
          "value": "ts-verification-code-service",
          "label": "ts-verification-code-service"
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkDelay --target-service mysql
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "NetworkDelay",
    "target_service": "mysql"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkDelay",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts",
    "target_service": "mysql"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill network delay parameters",
      "fields": [
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
        },
        {
          "name": "latency",
          "kind": "number_range",
          "required": true,
          "description": "Network latency",
          "min": 1,
          "max": 2000,
          "step": 1,
          "unit": "ms"
        },
        {
          "name": "correlation",
          "kind": "number_range",
          "required": true,
          "description": "Correlation percentage",
          "min": 0,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "jitter",
          "kind": "number_range",
          "required": true,
          "description": "Jitter",
          "min": 0,
          "max": 1000,
          "step": 1,
          "unit": "ms"
        },
        {
          "name": "direction",
          "kind": "enum",
          "required": true,
          "description": "Traffic direction",
          "options": [
            {
              "value": "to",
              "label": "to"
            },
            {
              "value": "from",
              "label": "from"
            },
            {
              "value": "both",
              "label": "both"
            }
          ]
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "latency, correlation, jitter and direction are required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkDelay --target-service mysql --latency 120 --correlation 50 --jitter 10 --direction both
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
    "chaos_type": "NetworkDelay",
    "target_service": "mysql",
    "duration": 5,
    "latency": 120,
    "correlation": 50,
    "jitter": 10,
    "direction": "both"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkDelay",
    "correlation": 50,
    "direction": "both",
    "duration": 5,
    "jitter": 10,
    "latency": 120,
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts",
    "target_service": "mysql"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill network delay parameters",
      "fields": [
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
        },
        {
          "name": "latency",
          "kind": "number_range",
          "required": true,
          "description": "Network latency",
          "min": 1,
          "max": 2000,
          "step": 1,
          "unit": "ms"
        },
        {
          "name": "correlation",
          "kind": "number_range",
          "required": true,
          "description": "Correlation percentage",
          "min": 0,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "jitter",
          "kind": "number_range",
          "required": true,
          "description": "Jitter",
          "min": 0,
          "max": 1000,
          "step": 1,
          "unit": "ms"
        },
        {
          "name": "direction",
          "kind": "enum",
          "required": true,
          "description": "Traffic direction",
          "options": [
            {
              "value": "to",
              "label": "to"
            },
            {
              "value": "from",
              "label": "from"
            },
            {
              "value": "both",
              "label": "both"
            }
          ]
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "NetworkDelay": {
          "Correlation": 50,
          "Direction": 3,
          "Duration": 5,
          "Jitter": 10,
          "Latency": 120,
          "NetworkPairIdx": 17,
          "System": 7
        }
      },
      "chaos_type": "NetworkDelay",
      "correlation": 50,
      "direction": "both",
      "duration": 5,
      "injection_point": {
        "source_service": "ts-auth-service",
        "span_names": [
          "ALTER table ts",
          "CREATE TABLE `ts`.`user_roles`",
          "CREATE table ts",
          "DELETE ts.auth_user",
          "DELETE ts.user_roles",
          "INSERT ts.auth_user",
          "INSERT ts.user_roles",
          "SELECT `ts`.`auth_user`",
          "SELECT ts",
          "SELECT ts.auth_user",
          "SELECT ts.ts",
          "SELECT ts.user_roles"
        ],
        "target_service": "mysql"
      },
      "jitter": 10,
      "latency": 120,
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "network_latency"
      ],
      "service": [
        "ts-auth-service",
        "mysql"
      ],
      "span": [
        "ALTER table ts",
        "CREATE TABLE `ts`.`user_roles`",
        "CREATE table ts",
        "DELETE ts.auth_user",
        "DELETE ts.user_roles",
        "INSERT ts.auth_user",
        "INSERT ts.user_roles",
        "SELECT `ts`.`auth_user`",
        "SELECT ts",
        "SELECT ts.auth_user",
        "SELECT ts.ts",
        "SELECT ts.user_roles"
      ]
    }
  },
  "apply_payload": {
    "NetworkDelay": {
      "Correlation": 50,
      "Direction": 3,
      "Duration": 5,
      "Jitter": 10,
      "Latency": 120,
      "NetworkPairIdx": 17,
      "System": 7
    }
  },
  "can_apply": true
}
```

## NetworkLoss

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkLoss
```
```json
{
  "mode": "guided",
  "stage": "select_network_target",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "NetworkLoss"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkLoss",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "target_service",
      "kind": "enum",
      "required": true,
      "description": "Select the target service for the network pair",
      "options": [
        {
          "value": "mysql",
          "label": "mysql"
        },
        {
          "value": "ts-verification-code-service",
          "label": "ts-verification-code-service"
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkLoss --target-service mysql
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "NetworkLoss",
    "target_service": "mysql"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkLoss",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts",
    "target_service": "mysql"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill network loss parameters",
      "fields": [
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
        },
        {
          "name": "loss",
          "kind": "number_range",
          "required": true,
          "description": "Packet loss percentage",
          "min": 1,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "correlation",
          "kind": "number_range",
          "required": true,
          "description": "Correlation percentage",
          "min": 0,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "direction",
          "kind": "enum",
          "required": true,
          "description": "Traffic direction",
          "options": [
            {
              "value": "to",
              "label": "to"
            },
            {
              "value": "from",
              "label": "from"
            },
            {
              "value": "both",
              "label": "both"
            }
          ]
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "loss, correlation and direction are required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkLoss --target-service mysql --loss 20 --correlation 50 --direction both
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
    "chaos_type": "NetworkLoss",
    "target_service": "mysql",
    "duration": 5,
    "correlation": 50,
    "loss": 20,
    "direction": "both"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkLoss",
    "correlation": 50,
    "direction": "both",
    "duration": 5,
    "loss": 20,
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts",
    "target_service": "mysql"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill network loss parameters",
      "fields": [
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
        },
        {
          "name": "loss",
          "kind": "number_range",
          "required": true,
          "description": "Packet loss percentage",
          "min": 1,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "correlation",
          "kind": "number_range",
          "required": true,
          "description": "Correlation percentage",
          "min": 0,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "direction",
          "kind": "enum",
          "required": true,
          "description": "Traffic direction",
          "options": [
            {
              "value": "to",
              "label": "to"
            },
            {
              "value": "from",
              "label": "from"
            },
            {
              "value": "both",
              "label": "both"
            }
          ]
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "NetworkLoss": {
          "Correlation": 50,
          "Direction": 3,
          "Duration": 5,
          "Loss": 20,
          "NetworkPairIdx": 17,
          "System": 7
        }
      },
      "chaos_type": "NetworkLoss",
      "correlation": 50,
      "direction": "both",
      "duration": 5,
      "injection_point": {
        "source_service": "ts-auth-service",
        "span_names": [
          "ALTER table ts",
          "CREATE TABLE `ts`.`user_roles`",
          "CREATE table ts",
          "DELETE ts.auth_user",
          "DELETE ts.user_roles",
          "INSERT ts.auth_user",
          "INSERT ts.user_roles",
          "SELECT `ts`.`auth_user`",
          "SELECT ts",
          "SELECT ts.auth_user",
          "SELECT ts.ts",
          "SELECT ts.user_roles"
        ],
        "target_service": "mysql"
      },
      "loss": 20,
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "network_latency"
      ],
      "service": [
        "ts-auth-service",
        "mysql"
      ],
      "span": [
        "ALTER table ts",
        "CREATE TABLE `ts`.`user_roles`",
        "CREATE table ts",
        "DELETE ts.auth_user",
        "DELETE ts.user_roles",
        "INSERT ts.auth_user",
        "INSERT ts.user_roles",
        "SELECT `ts`.`auth_user`",
        "SELECT ts",
        "SELECT ts.auth_user",
        "SELECT ts.ts",
        "SELECT ts.user_roles"
      ]
    }
  },
  "apply_payload": {
    "NetworkLoss": {
      "Correlation": 50,
      "Direction": 3,
      "Duration": 5,
      "Loss": 20,
      "NetworkPairIdx": 17,
      "System": 7
    }
  },
  "can_apply": true
}
```

## NetworkDuplicate

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkDuplicate
```
```json
{
  "mode": "guided",
  "stage": "select_network_target",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "NetworkDuplicate"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkDuplicate",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "target_service",
      "kind": "enum",
      "required": true,
      "description": "Select the target service for the network pair",
      "options": [
        {
          "value": "mysql",
          "label": "mysql"
        },
        {
          "value": "ts-verification-code-service",
          "label": "ts-verification-code-service"
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkDuplicate --target-service mysql
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "NetworkDuplicate",
    "target_service": "mysql"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkDuplicate",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts",
    "target_service": "mysql"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill network duplicate parameters",
      "fields": [
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
        },
        {
          "name": "duplicate",
          "kind": "number_range",
          "required": true,
          "description": "Packet duplication percentage",
          "min": 1,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "correlation",
          "kind": "number_range",
          "required": true,
          "description": "Correlation percentage",
          "min": 0,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "direction",
          "kind": "enum",
          "required": true,
          "description": "Traffic direction",
          "options": [
            {
              "value": "to",
              "label": "to"
            },
            {
              "value": "from",
              "label": "from"
            },
            {
              "value": "both",
              "label": "both"
            }
          ]
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "duplicate, correlation and direction are required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkDuplicate --target-service mysql --duplicate 10 --correlation 30 --direction both
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
    "chaos_type": "NetworkDuplicate",
    "target_service": "mysql",
    "duration": 5,
    "correlation": 30,
    "duplicate": 10,
    "direction": "both"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkDuplicate",
    "correlation": 30,
    "direction": "both",
    "duplicate": 10,
    "duration": 5,
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts",
    "target_service": "mysql"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill network duplicate parameters",
      "fields": [
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
        },
        {
          "name": "duplicate",
          "kind": "number_range",
          "required": true,
          "description": "Packet duplication percentage",
          "min": 1,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "correlation",
          "kind": "number_range",
          "required": true,
          "description": "Correlation percentage",
          "min": 0,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "direction",
          "kind": "enum",
          "required": true,
          "description": "Traffic direction",
          "options": [
            {
              "value": "to",
              "label": "to"
            },
            {
              "value": "from",
              "label": "from"
            },
            {
              "value": "both",
              "label": "both"
            }
          ]
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "NetworkDuplicate": {
          "Correlation": 30,
          "Direction": 3,
          "Duplicate": 10,
          "Duration": 5,
          "NetworkPairIdx": 17,
          "System": 7
        }
      },
      "chaos_type": "NetworkDuplicate",
      "correlation": 30,
      "direction": "both",
      "duplicate": 10,
      "duration": 5,
      "injection_point": {
        "source_service": "ts-auth-service",
        "span_names": [
          "ALTER table ts",
          "CREATE TABLE `ts`.`user_roles`",
          "CREATE table ts",
          "DELETE ts.auth_user",
          "DELETE ts.user_roles",
          "INSERT ts.auth_user",
          "INSERT ts.user_roles",
          "SELECT `ts`.`auth_user`",
          "SELECT ts",
          "SELECT ts.auth_user",
          "SELECT ts.ts",
          "SELECT ts.user_roles"
        ],
        "target_service": "mysql"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "network_latency"
      ],
      "service": [
        "ts-auth-service",
        "mysql"
      ],
      "span": [
        "ALTER table ts",
        "CREATE TABLE `ts`.`user_roles`",
        "CREATE table ts",
        "DELETE ts.auth_user",
        "DELETE ts.user_roles",
        "INSERT ts.auth_user",
        "INSERT ts.user_roles",
        "SELECT `ts`.`auth_user`",
        "SELECT ts",
        "SELECT ts.auth_user",
        "SELECT ts.ts",
        "SELECT ts.user_roles"
      ]
    }
  },
  "apply_payload": {
    "NetworkDuplicate": {
      "Correlation": 30,
      "Direction": 3,
      "Duplicate": 10,
      "Duration": 5,
      "NetworkPairIdx": 17,
      "System": 7
    }
  },
  "can_apply": true
}
```

## NetworkCorrupt

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkCorrupt
```
```json
{
  "mode": "guided",
  "stage": "select_network_target",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "NetworkCorrupt"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkCorrupt",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "target_service",
      "kind": "enum",
      "required": true,
      "description": "Select the target service for the network pair",
      "options": [
        {
          "value": "mysql",
          "label": "mysql"
        },
        {
          "value": "ts-verification-code-service",
          "label": "ts-verification-code-service"
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkCorrupt --target-service mysql
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "NetworkCorrupt",
    "target_service": "mysql"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkCorrupt",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts",
    "target_service": "mysql"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill network corruption parameters",
      "fields": [
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
        },
        {
          "name": "corrupt",
          "kind": "number_range",
          "required": true,
          "description": "Packet corruption percentage",
          "min": 1,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "correlation",
          "kind": "number_range",
          "required": true,
          "description": "Correlation percentage",
          "min": 0,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "direction",
          "kind": "enum",
          "required": true,
          "description": "Traffic direction",
          "options": [
            {
              "value": "to",
              "label": "to"
            },
            {
              "value": "from",
              "label": "from"
            },
            {
              "value": "both",
              "label": "both"
            }
          ]
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "corrupt, correlation and direction are required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkCorrupt --target-service mysql --corrupt 5 --correlation 20 --direction both
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
    "chaos_type": "NetworkCorrupt",
    "target_service": "mysql",
    "duration": 5,
    "correlation": 20,
    "corrupt": 5,
    "direction": "both"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkCorrupt",
    "correlation": 20,
    "corrupt": 5,
    "direction": "both",
    "duration": 5,
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts",
    "target_service": "mysql"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill network corruption parameters",
      "fields": [
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
        },
        {
          "name": "corrupt",
          "kind": "number_range",
          "required": true,
          "description": "Packet corruption percentage",
          "min": 1,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "correlation",
          "kind": "number_range",
          "required": true,
          "description": "Correlation percentage",
          "min": 0,
          "max": 100,
          "step": 1,
          "unit": "%"
        },
        {
          "name": "direction",
          "kind": "enum",
          "required": true,
          "description": "Traffic direction",
          "options": [
            {
              "value": "to",
              "label": "to"
            },
            {
              "value": "from",
              "label": "from"
            },
            {
              "value": "both",
              "label": "both"
            }
          ]
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "NetworkCorrupt": {
          "Correlation": 20,
          "Corrupt": 5,
          "Direction": 3,
          "Duration": 5,
          "NetworkPairIdx": 17,
          "System": 7
        }
      },
      "chaos_type": "NetworkCorrupt",
      "correlation": 20,
      "corrupt": 5,
      "direction": "both",
      "duration": 5,
      "injection_point": {
        "source_service": "ts-auth-service",
        "span_names": [
          "ALTER table ts",
          "CREATE TABLE `ts`.`user_roles`",
          "CREATE table ts",
          "DELETE ts.auth_user",
          "DELETE ts.user_roles",
          "INSERT ts.auth_user",
          "INSERT ts.user_roles",
          "SELECT `ts`.`auth_user`",
          "SELECT ts",
          "SELECT ts.auth_user",
          "SELECT ts.ts",
          "SELECT ts.user_roles"
        ],
        "target_service": "mysql"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "network_latency"
      ],
      "service": [
        "ts-auth-service",
        "mysql"
      ],
      "span": [
        "ALTER table ts",
        "CREATE TABLE `ts`.`user_roles`",
        "CREATE table ts",
        "DELETE ts.auth_user",
        "DELETE ts.user_roles",
        "INSERT ts.auth_user",
        "INSERT ts.user_roles",
        "SELECT `ts`.`auth_user`",
        "SELECT ts",
        "SELECT ts.auth_user",
        "SELECT ts.ts",
        "SELECT ts.user_roles"
      ]
    }
  },
  "apply_payload": {
    "NetworkCorrupt": {
      "Correlation": 20,
      "Corrupt": 5,
      "Direction": 3,
      "Duration": 5,
      "NetworkPairIdx": 17,
      "System": 7
    }
  },
  "can_apply": true
}
```

## NetworkBandwidth

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkBandwidth
```
```json
{
  "mode": "guided",
  "stage": "select_network_target",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "NetworkBandwidth"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkBandwidth",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "target_service",
      "kind": "enum",
      "required": true,
      "description": "Select the target service for the network pair",
      "options": [
        {
          "value": "mysql",
          "label": "mysql"
        },
        {
          "value": "ts-verification-code-service",
          "label": "ts-verification-code-service"
        }
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkBandwidth --target-service mysql
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "NetworkBandwidth",
    "target_service": "mysql"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "NetworkBandwidth",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts",
    "target_service": "mysql"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill network bandwidth parameters",
      "fields": [
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
        },
        {
          "name": "rate",
          "kind": "number_range",
          "required": true,
          "description": "Bandwidth rate",
          "min": 1,
          "max": 1000000,
          "step": 1,
          "unit": "kbps"
        },
        {
          "name": "limit",
          "kind": "number_range",
          "required": true,
          "description": "Bandwidth queue limit",
          "min": 1,
          "max": 10000,
          "step": 1,
          "unit": "byte"
        },
        {
          "name": "buffer",
          "kind": "number_range",
          "required": true,
          "description": "Bandwidth buffer",
          "min": 1,
          "max": 10000,
          "step": 1,
          "unit": "byte"
        },
        {
          "name": "direction",
          "kind": "enum",
          "required": true,
          "description": "Traffic direction",
          "options": [
            {
              "value": "to",
              "label": "to"
            },
            {
              "value": "from",
              "label": "from"
            },
            {
              "value": "both",
              "label": "both"
            }
          ]
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "rate, limit, buffer and direction are required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type NetworkBandwidth --target-service mysql --rate 1024 --limit 2048 --buffer 4096 --direction both
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
    "chaos_type": "NetworkBandwidth",
    "target_service": "mysql",
    "duration": 5,
    "rate": 1024,
    "limit": 2048,
    "buffer": 4096,
    "direction": "both"
  },
  "resolved": {
    "app": "ts-auth-service",
    "buffer": 4096,
    "chaos_type": "NetworkBandwidth",
    "direction": "both",
    "duration": 5,
    "limit": 2048,
    "namespace": "ts",
    "rate": 1024,
    "system": "ts",
    "system_type": "ts",
    "target_service": "mysql"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill network bandwidth parameters",
      "fields": [
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
        },
        {
          "name": "rate",
          "kind": "number_range",
          "required": true,
          "description": "Bandwidth rate",
          "min": 1,
          "max": 1000000,
          "step": 1,
          "unit": "kbps"
        },
        {
          "name": "limit",
          "kind": "number_range",
          "required": true,
          "description": "Bandwidth queue limit",
          "min": 1,
          "max": 10000,
          "step": 1,
          "unit": "byte"
        },
        {
          "name": "buffer",
          "kind": "number_range",
          "required": true,
          "description": "Bandwidth buffer",
          "min": 1,
          "max": 10000,
          "step": 1,
          "unit": "byte"
        },
        {
          "name": "direction",
          "kind": "enum",
          "required": true,
          "description": "Traffic direction",
          "options": [
            {
              "value": "to",
              "label": "to"
            },
            {
              "value": "from",
              "label": "from"
            },
            {
              "value": "both",
              "label": "both"
            }
          ]
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "NetworkBandwidth": {
          "Buffer": 4096,
          "Direction": 3,
          "Duration": 5,
          "Limit": 2048,
          "NetworkPairIdx": 17,
          "Rate": 1024,
          "System": 7
        }
      },
      "buffer": 4096,
      "chaos_type": "NetworkBandwidth",
      "direction": "both",
      "duration": 5,
      "injection_point": {
        "source_service": "ts-auth-service",
        "span_names": [
          "ALTER table ts",
          "CREATE TABLE `ts`.`user_roles`",
          "CREATE table ts",
          "DELETE ts.auth_user",
          "DELETE ts.user_roles",
          "INSERT ts.auth_user",
          "INSERT ts.user_roles",
          "SELECT `ts`.`auth_user`",
          "SELECT ts",
          "SELECT ts.auth_user",
          "SELECT ts.ts",
          "SELECT ts.user_roles"
        ],
        "target_service": "mysql"
      },
      "limit": 2048,
      "namespace": "ts",
      "rate": 1024,
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "network_latency"
      ],
      "service": [
        "ts-auth-service",
        "mysql"
      ],
      "span": [
        "ALTER table ts",
        "CREATE TABLE `ts`.`user_roles`",
        "CREATE table ts",
        "DELETE ts.auth_user",
        "DELETE ts.user_roles",
        "INSERT ts.auth_user",
        "INSERT ts.user_roles",
        "SELECT `ts`.`auth_user`",
        "SELECT ts",
        "SELECT ts.auth_user",
        "SELECT ts.ts",
        "SELECT ts.user_roles"
      ]
    }
  },
  "apply_payload": {
    "NetworkBandwidth": {
      "Buffer": 4096,
      "Direction": 3,
      "Duration": 5,
      "Limit": 2048,
      "NetworkPairIdx": 17,
      "Rate": 1024,
      "System": 7
    }
  },
  "can_apply": true
}
```

## JVMLatency

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMLatency
```
```json
{
  "mode": "guided",
  "stage": "select_jvm_method",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMLatency"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMLatency",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "method_ref",
      "kind": "object_ref",
      "required": true,
      "description": "Select the JVM method for latency injection",
      "options": [
        {
          "value": "auth.AuthApplication#main",
          "label": "auth.AuthApplication#main",
          "metadata": {
            "class": "auth.AuthApplication",
            "method": "main"
          }
        },
        {
          "value": "auth.constant.AuthConstant#AuthConstant",
          "label": "auth.constant.AuthConstant#AuthConstant",
          "metadata": {
            "class": "auth.constant.AuthConstant",
            "method": "AuthConstant"
          }
        },
        {
          "value": "auth.constant.InfoConstant#InfoConstant",
          "label": "auth.constant.InfoConstant#InfoConstant",
          "metadata": {
            "class": "auth.constant.InfoConstant",
            "method": "InfoConstant"
          }
        },
        {
          "value": "auth.controller.AuthController#createDefaultUser",
          "label": "auth.controller.AuthController#createDefaultUser",
          "metadata": {
            "class": "auth.controller.AuthController",
            "method": "createDefaultUser"
          }
        },
        {
          "value": "auth.controller.AuthController#getHello",
          "label": "auth.controller.AuthController#getHello",
          "metadata": {
            "class": "auth.controller.AuthController",
            "method": "getHello"
          }
        },
        {
          "value": "auth.controller.UserController#deleteUserById",
          "label": "auth.controller.UserController#deleteUserById",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "deleteUserById"
          }
        },
        {
          "value": "auth.controller.UserController#getAllUser",
          "label": "auth.controller.UserController#getAllUser",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getAllUser"
          }
        },
        {
          "value": "auth.controller.UserController#getHello",
          "label": "auth.controller.UserController#getHello",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getHello"
          }
        },
        {
          "value": "auth.controller.UserController#getToken",
          "label": "auth.controller.UserController#getToken",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getToken"
          }
        },
        {
          "value": "auth.entity.User#getAuthorities",
          "label": "auth.entity.User#getAuthorities",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getAuthorities"
          }
        },
        {
          "value": "auth.entity.User#getPassword",
          "label": "auth.entity.User#getPassword",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getPassword"
          }
        },
        {
          "value": "auth.entity.User#getUsername",
          "label": "auth.entity.User#getUsername",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getUsername"
          }
        },
        {
          "value": "auth.entity.User#isAccountNonExpired",
          "label": "auth.entity.User#isAccountNonExpired",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isAccountNonExpired"
          }
        },
        {
          "value": "auth.entity.User#isAccountNonLocked",
          "label": "auth.entity.User#isAccountNonLocked",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isAccountNonLocked"
          }
        },
        {
          "value": "auth.entity.User#isCredentialsNonExpired",
          "label": "auth.entity.User#isCredentialsNonExpired",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isCredentialsNonExpired"
          }
        },
        {
          "value": "auth.entity.User#isEnabled",
          "label": "auth.entity.User#isEnabled",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isEnabled"
          }
        },
        {
          "value": "auth.exception.UserOperationException#UserOperationException",
          "label": "auth.exception.UserOperationException#UserOperationException",
          "metadata": {
            "class": "auth.exception.UserOperationException",
            "method": "UserOperationException"
          }
        },
        {
          "value": "auth.exception.handler.GlobalExceptionHandler#handleUserNotFoundException",
          "label": "auth.exception.handler.GlobalExceptionHandler#handleUserNotFoundException",
          "metadata": {
            "class": "auth.exception.handler.GlobalExceptionHandler",
            "method": "handleUserNotFoundException"
          }
        },
        {
          "value": "auth.init.InitUser#run",
          "label": "auth.init.InitUser#run",
          "metadata": {
            "class": "auth.init.InitUser",
            "method": "run"
          }
        },
        {
          "value": "auth.security.UserDetailsServiceImpl#loadUserByUsername",
          "label": "auth.security.UserDetailsServiceImpl#loadUserByUsername",
          "metadata": {
            "class": "auth.security.UserDetailsServiceImpl",
            "method": "loadUserByUsername"
          }
        },
        {
          "value": "auth.security.jwt.JWTProvider#createToken",
          "label": "auth.security.jwt.JWTProvider#createToken",
          "metadata": {
            "class": "auth.security.jwt.JWTProvider",
            "method": "createToken"
          }
        },
        {
          "value": "auth.security.jwt.JWTProvider#init",
          "label": "auth.security.jwt.JWTProvider#init",
          "metadata": {
            "class": "auth.security.jwt.JWTProvider",
            "method": "init"
          }
        },
        {
          "value": "auth.service.impl.TokenServiceImpl#getServiceUrl",
          "label": "auth.service.impl.TokenServiceImpl#getServiceUrl",
          "metadata": {
            "class": "auth.service.impl.TokenServiceImpl",
            "method": "getServiceUrl"
          }
        },
        {
          "value": "auth.service.impl.TokenServiceImpl#getToken",
          "label": "auth.service.impl.TokenServiceImpl#getToken",
          "metadata": {
            "class": "auth.service.impl.TokenServiceImpl",
            "method": "getToken"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#checkUserCreateInfo",
          "label": "auth.service.impl.UserServiceImpl#checkUserCreateInfo",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "checkUserCreateInfo"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#createDefaultAuthUser",
          "label": "auth.service.impl.UserServiceImpl#createDefaultAuthUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "createDefaultAuthUser"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#deleteByUserId",
          "label": "auth.service.impl.UserServiceImpl#deleteByUserId",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "deleteByUserId"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#getAllUser",
          "label": "auth.service.impl.UserServiceImpl#getAllUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "getAllUser"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#saveUser",
          "label": "auth.service.impl.UserServiceImpl#saveUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "saveUser"
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

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMLatency --class auth.security.jwt.JWTProvider --method createToken
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMLatency",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMLatency",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill JVM latency parameters",
      "fields": [
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
        },
        {
          "name": "latency_duration",
          "kind": "number_range",
          "required": true,
          "description": "JVM latency duration",
          "min": 1,
          "max": 5000,
          "step": 1,
          "unit": "ms"
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "latency_duration is required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMLatency --class auth.security.jwt.JWTProvider --method createToken --latency-duration 100
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
    "chaos_type": "JVMLatency",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken",
    "duration": 5,
    "latency_duration": 100
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMLatency",
    "class": "auth.security.jwt.JWTProvider",
    "duration": 5,
    "latency_duration": 100,
    "method": "createToken",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill JVM latency parameters",
      "fields": [
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
        },
        {
          "name": "latency_duration",
          "kind": "number_range",
          "required": true,
          "description": "JVM latency duration",
          "min": 1,
          "max": 5000,
          "step": 1,
          "unit": "ms"
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "JVMLatency": {
          "Duration": 5,
          "LatencyDuration": 100,
          "MethodIdx": 132,
          "System": 7
        }
      },
      "chaos_type": "JVMLatency",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "class_name": "auth.security.jwt.JWTProvider",
        "method_name": "createToken"
      },
      "latency_duration": 100,
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "function": [
        "auth.security.jwt.JWTProvider.createToken"
      ],
      "metric": [
        "network_latency"
      ],
      "service": [
        "ts-auth-service"
      ]
    }
  },
  "apply_payload": {
    "JVMLatency": {
      "Duration": 5,
      "LatencyDuration": 100,
      "MethodIdx": 132,
      "System": 7
    }
  },
  "can_apply": true
}
```

## JVMReturn

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMReturn
```
```json
{
  "mode": "guided",
  "stage": "select_jvm_method",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMReturn"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMReturn",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "method_ref",
      "kind": "object_ref",
      "required": true,
      "description": "Select the JVM method for return value injection",
      "options": [
        {
          "value": "auth.AuthApplication#main",
          "label": "auth.AuthApplication#main",
          "metadata": {
            "class": "auth.AuthApplication",
            "method": "main"
          }
        },
        {
          "value": "auth.constant.AuthConstant#AuthConstant",
          "label": "auth.constant.AuthConstant#AuthConstant",
          "metadata": {
            "class": "auth.constant.AuthConstant",
            "method": "AuthConstant"
          }
        },
        {
          "value": "auth.constant.InfoConstant#InfoConstant",
          "label": "auth.constant.InfoConstant#InfoConstant",
          "metadata": {
            "class": "auth.constant.InfoConstant",
            "method": "InfoConstant"
          }
        },
        {
          "value": "auth.controller.AuthController#createDefaultUser",
          "label": "auth.controller.AuthController#createDefaultUser",
          "metadata": {
            "class": "auth.controller.AuthController",
            "method": "createDefaultUser"
          }
        },
        {
          "value": "auth.controller.AuthController#getHello",
          "label": "auth.controller.AuthController#getHello",
          "metadata": {
            "class": "auth.controller.AuthController",
            "method": "getHello"
          }
        },
        {
          "value": "auth.controller.UserController#deleteUserById",
          "label": "auth.controller.UserController#deleteUserById",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "deleteUserById"
          }
        },
        {
          "value": "auth.controller.UserController#getAllUser",
          "label": "auth.controller.UserController#getAllUser",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getAllUser"
          }
        },
        {
          "value": "auth.controller.UserController#getHello",
          "label": "auth.controller.UserController#getHello",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getHello"
          }
        },
        {
          "value": "auth.controller.UserController#getToken",
          "label": "auth.controller.UserController#getToken",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getToken"
          }
        },
        {
          "value": "auth.entity.User#getAuthorities",
          "label": "auth.entity.User#getAuthorities",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getAuthorities"
          }
        },
        {
          "value": "auth.entity.User#getPassword",
          "label": "auth.entity.User#getPassword",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getPassword"
          }
        },
        {
          "value": "auth.entity.User#getUsername",
          "label": "auth.entity.User#getUsername",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getUsername"
          }
        },
        {
          "value": "auth.entity.User#isAccountNonExpired",
          "label": "auth.entity.User#isAccountNonExpired",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isAccountNonExpired"
          }
        },
        {
          "value": "auth.entity.User#isAccountNonLocked",
          "label": "auth.entity.User#isAccountNonLocked",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isAccountNonLocked"
          }
        },
        {
          "value": "auth.entity.User#isCredentialsNonExpired",
          "label": "auth.entity.User#isCredentialsNonExpired",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isCredentialsNonExpired"
          }
        },
        {
          "value": "auth.entity.User#isEnabled",
          "label": "auth.entity.User#isEnabled",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isEnabled"
          }
        },
        {
          "value": "auth.exception.UserOperationException#UserOperationException",
          "label": "auth.exception.UserOperationException#UserOperationException",
          "metadata": {
            "class": "auth.exception.UserOperationException",
            "method": "UserOperationException"
          }
        },
        {
          "value": "auth.exception.handler.GlobalExceptionHandler#handleUserNotFoundException",
          "label": "auth.exception.handler.GlobalExceptionHandler#handleUserNotFoundException",
          "metadata": {
            "class": "auth.exception.handler.GlobalExceptionHandler",
            "method": "handleUserNotFoundException"
          }
        },
        {
          "value": "auth.init.InitUser#run",
          "label": "auth.init.InitUser#run",
          "metadata": {
            "class": "auth.init.InitUser",
            "method": "run"
          }
        },
        {
          "value": "auth.security.UserDetailsServiceImpl#loadUserByUsername",
          "label": "auth.security.UserDetailsServiceImpl#loadUserByUsername",
          "metadata": {
            "class": "auth.security.UserDetailsServiceImpl",
            "method": "loadUserByUsername"
          }
        },
        {
          "value": "auth.security.jwt.JWTProvider#createToken",
          "label": "auth.security.jwt.JWTProvider#createToken",
          "metadata": {
            "class": "auth.security.jwt.JWTProvider",
            "method": "createToken"
          }
        },
        {
          "value": "auth.security.jwt.JWTProvider#init",
          "label": "auth.security.jwt.JWTProvider#init",
          "metadata": {
            "class": "auth.security.jwt.JWTProvider",
            "method": "init"
          }
        },
        {
          "value": "auth.service.impl.TokenServiceImpl#getServiceUrl",
          "label": "auth.service.impl.TokenServiceImpl#getServiceUrl",
          "metadata": {
            "class": "auth.service.impl.TokenServiceImpl",
            "method": "getServiceUrl"
          }
        },
        {
          "value": "auth.service.impl.TokenServiceImpl#getToken",
          "label": "auth.service.impl.TokenServiceImpl#getToken",
          "metadata": {
            "class": "auth.service.impl.TokenServiceImpl",
            "method": "getToken"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#checkUserCreateInfo",
          "label": "auth.service.impl.UserServiceImpl#checkUserCreateInfo",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "checkUserCreateInfo"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#createDefaultAuthUser",
          "label": "auth.service.impl.UserServiceImpl#createDefaultAuthUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "createDefaultAuthUser"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#deleteByUserId",
          "label": "auth.service.impl.UserServiceImpl#deleteByUserId",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "deleteByUserId"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#getAllUser",
          "label": "auth.service.impl.UserServiceImpl#getAllUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "getAllUser"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#saveUser",
          "label": "auth.service.impl.UserServiceImpl#saveUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "saveUser"
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

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMReturn --class auth.security.jwt.JWTProvider --method createToken
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMReturn",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMReturn",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill JVM return parameters",
      "fields": [
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
        },
        {
          "name": "return_type",
          "kind": "enum",
          "required": true,
          "description": "Return type",
          "options": [
            {
              "value": "string",
              "label": "string"
            },
            {
              "value": "int",
              "label": "int"
            }
          ]
        },
        {
          "name": "return_value_opt",
          "kind": "enum",
          "required": true,
          "description": "Return value strategy",
          "options": [
            {
              "value": "default",
              "label": "default"
            },
            {
              "value": "random",
              "label": "random"
            }
          ]
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "return_type and return_value_opt are required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMReturn --class auth.security.jwt.JWTProvider --method createToken --return-type string --return-value-opt default
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
    "chaos_type": "JVMReturn",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken",
    "duration": 5,
    "return_type": "string",
    "return_value_opt": "default"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMReturn",
    "class": "auth.security.jwt.JWTProvider",
    "duration": 5,
    "method": "createToken",
    "namespace": "ts",
    "return_type": "string",
    "return_value_opt": "default",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill JVM return parameters",
      "fields": [
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
        },
        {
          "name": "return_type",
          "kind": "enum",
          "required": true,
          "description": "Return type",
          "options": [
            {
              "value": "string",
              "label": "string"
            },
            {
              "value": "int",
              "label": "int"
            }
          ]
        },
        {
          "name": "return_value_opt",
          "kind": "enum",
          "required": true,
          "description": "Return value strategy",
          "options": [
            {
              "value": "default",
              "label": "default"
            },
            {
              "value": "random",
              "label": "random"
            }
          ]
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "JVMReturn": {
          "Duration": 5,
          "MethodIdx": 132,
          "ReturnType": 1,
          "ReturnValueOpt": 0,
          "System": 7
        }
      },
      "chaos_type": "JVMReturn",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "class_name": "auth.security.jwt.JWTProvider",
        "method_name": "createToken"
      },
      "namespace": "ts",
      "return_type": "string",
      "return_value_opt": "default",
      "system": "ts"
    },
    "groundtruth": {
      "function": [
        "auth.security.jwt.JWTProvider.createToken"
      ],
      "service": [
        "ts-auth-service"
      ]
    }
  },
  "apply_payload": {
    "JVMReturn": {
      "Duration": 5,
      "MethodIdx": 132,
      "ReturnType": 1,
      "ReturnValueOpt": 0,
      "System": 7
    }
  },
  "can_apply": true
}
```

## JVMException

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMException
```
```json
{
  "mode": "guided",
  "stage": "select_jvm_method",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMException"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMException",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "method_ref",
      "kind": "object_ref",
      "required": true,
      "description": "Select the JVM method for exception injection",
      "options": [
        {
          "value": "auth.AuthApplication#main",
          "label": "auth.AuthApplication#main",
          "metadata": {
            "class": "auth.AuthApplication",
            "method": "main"
          }
        },
        {
          "value": "auth.constant.AuthConstant#AuthConstant",
          "label": "auth.constant.AuthConstant#AuthConstant",
          "metadata": {
            "class": "auth.constant.AuthConstant",
            "method": "AuthConstant"
          }
        },
        {
          "value": "auth.constant.InfoConstant#InfoConstant",
          "label": "auth.constant.InfoConstant#InfoConstant",
          "metadata": {
            "class": "auth.constant.InfoConstant",
            "method": "InfoConstant"
          }
        },
        {
          "value": "auth.controller.AuthController#createDefaultUser",
          "label": "auth.controller.AuthController#createDefaultUser",
          "metadata": {
            "class": "auth.controller.AuthController",
            "method": "createDefaultUser"
          }
        },
        {
          "value": "auth.controller.AuthController#getHello",
          "label": "auth.controller.AuthController#getHello",
          "metadata": {
            "class": "auth.controller.AuthController",
            "method": "getHello"
          }
        },
        {
          "value": "auth.controller.UserController#deleteUserById",
          "label": "auth.controller.UserController#deleteUserById",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "deleteUserById"
          }
        },
        {
          "value": "auth.controller.UserController#getAllUser",
          "label": "auth.controller.UserController#getAllUser",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getAllUser"
          }
        },
        {
          "value": "auth.controller.UserController#getHello",
          "label": "auth.controller.UserController#getHello",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getHello"
          }
        },
        {
          "value": "auth.controller.UserController#getToken",
          "label": "auth.controller.UserController#getToken",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getToken"
          }
        },
        {
          "value": "auth.entity.User#getAuthorities",
          "label": "auth.entity.User#getAuthorities",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getAuthorities"
          }
        },
        {
          "value": "auth.entity.User#getPassword",
          "label": "auth.entity.User#getPassword",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getPassword"
          }
        },
        {
          "value": "auth.entity.User#getUsername",
          "label": "auth.entity.User#getUsername",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getUsername"
          }
        },
        {
          "value": "auth.entity.User#isAccountNonExpired",
          "label": "auth.entity.User#isAccountNonExpired",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isAccountNonExpired"
          }
        },
        {
          "value": "auth.entity.User#isAccountNonLocked",
          "label": "auth.entity.User#isAccountNonLocked",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isAccountNonLocked"
          }
        },
        {
          "value": "auth.entity.User#isCredentialsNonExpired",
          "label": "auth.entity.User#isCredentialsNonExpired",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isCredentialsNonExpired"
          }
        },
        {
          "value": "auth.entity.User#isEnabled",
          "label": "auth.entity.User#isEnabled",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isEnabled"
          }
        },
        {
          "value": "auth.exception.UserOperationException#UserOperationException",
          "label": "auth.exception.UserOperationException#UserOperationException",
          "metadata": {
            "class": "auth.exception.UserOperationException",
            "method": "UserOperationException"
          }
        },
        {
          "value": "auth.exception.handler.GlobalExceptionHandler#handleUserNotFoundException",
          "label": "auth.exception.handler.GlobalExceptionHandler#handleUserNotFoundException",
          "metadata": {
            "class": "auth.exception.handler.GlobalExceptionHandler",
            "method": "handleUserNotFoundException"
          }
        },
        {
          "value": "auth.init.InitUser#run",
          "label": "auth.init.InitUser#run",
          "metadata": {
            "class": "auth.init.InitUser",
            "method": "run"
          }
        },
        {
          "value": "auth.security.UserDetailsServiceImpl#loadUserByUsername",
          "label": "auth.security.UserDetailsServiceImpl#loadUserByUsername",
          "metadata": {
            "class": "auth.security.UserDetailsServiceImpl",
            "method": "loadUserByUsername"
          }
        },
        {
          "value": "auth.security.jwt.JWTProvider#createToken",
          "label": "auth.security.jwt.JWTProvider#createToken",
          "metadata": {
            "class": "auth.security.jwt.JWTProvider",
            "method": "createToken"
          }
        },
        {
          "value": "auth.security.jwt.JWTProvider#init",
          "label": "auth.security.jwt.JWTProvider#init",
          "metadata": {
            "class": "auth.security.jwt.JWTProvider",
            "method": "init"
          }
        },
        {
          "value": "auth.service.impl.TokenServiceImpl#getServiceUrl",
          "label": "auth.service.impl.TokenServiceImpl#getServiceUrl",
          "metadata": {
            "class": "auth.service.impl.TokenServiceImpl",
            "method": "getServiceUrl"
          }
        },
        {
          "value": "auth.service.impl.TokenServiceImpl#getToken",
          "label": "auth.service.impl.TokenServiceImpl#getToken",
          "metadata": {
            "class": "auth.service.impl.TokenServiceImpl",
            "method": "getToken"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#checkUserCreateInfo",
          "label": "auth.service.impl.UserServiceImpl#checkUserCreateInfo",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "checkUserCreateInfo"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#createDefaultAuthUser",
          "label": "auth.service.impl.UserServiceImpl#createDefaultAuthUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "createDefaultAuthUser"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#deleteByUserId",
          "label": "auth.service.impl.UserServiceImpl#deleteByUserId",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "deleteByUserId"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#getAllUser",
          "label": "auth.service.impl.UserServiceImpl#getAllUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "getAllUser"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#saveUser",
          "label": "auth.service.impl.UserServiceImpl#saveUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "saveUser"
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

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMException --class auth.security.jwt.JWTProvider --method createToken
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMException",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMException",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill JVM exception parameters",
      "fields": [
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
        },
        {
          "name": "exception_opt",
          "kind": "enum",
          "required": true,
          "description": "Exception strategy",
          "options": [
            {
              "value": "default",
              "label": "default"
            },
            {
              "value": "random",
              "label": "random"
            }
          ]
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "exception_opt is required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMException --class auth.security.jwt.JWTProvider --method createToken --exception-opt default
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
    "chaos_type": "JVMException",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken",
    "duration": 5,
    "exception_opt": "default"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMException",
    "class": "auth.security.jwt.JWTProvider",
    "duration": 5,
    "exception_opt": "default",
    "method": "createToken",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill JVM exception parameters",
      "fields": [
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
        },
        {
          "name": "exception_opt",
          "kind": "enum",
          "required": true,
          "description": "Exception strategy",
          "options": [
            {
              "value": "default",
              "label": "default"
            },
            {
              "value": "random",
              "label": "random"
            }
          ]
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "JVMException": {
          "Duration": 5,
          "ExceptionOpt": 0,
          "MethodIdx": 132,
          "System": 7
        }
      },
      "chaos_type": "JVMException",
      "duration": 5,
      "exception_opt": "default",
      "injection_point": {
        "app_name": "ts-auth-service",
        "class_name": "auth.security.jwt.JWTProvider",
        "method_name": "createToken"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "function": [
        "auth.security.jwt.JWTProvider.createToken"
      ],
      "service": [
        "ts-auth-service"
      ]
    }
  },
  "apply_payload": {
    "JVMException": {
      "Duration": 5,
      "ExceptionOpt": 0,
      "MethodIdx": 132,
      "System": 7
    }
  },
  "can_apply": true
}
```

## JVMGarbageCollector

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMGarbageCollector
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
    "chaos_type": "JVMGarbageCollector",
    "duration": 5
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMGarbageCollector",
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
        "JVMGarbageCollector": {
          "AppIdx": 7,
          "Duration": 5,
          "System": 7
        }
      },
      "chaos_type": "JVMGarbageCollector",
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
    "JVMGarbageCollector": {
      "AppIdx": 7,
      "Duration": 5,
      "System": 7
    }
  },
  "can_apply": true
}
```

## JVMCPUStress

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMCPUStress
```
```json
{
  "mode": "guided",
  "stage": "select_jvm_method",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMCPUStress"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMCPUStress",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "method_ref",
      "kind": "object_ref",
      "required": true,
      "description": "Select the JVM method for CPU stress injection",
      "options": [
        {
          "value": "auth.AuthApplication#main",
          "label": "auth.AuthApplication#main",
          "metadata": {
            "class": "auth.AuthApplication",
            "method": "main"
          }
        },
        {
          "value": "auth.constant.AuthConstant#AuthConstant",
          "label": "auth.constant.AuthConstant#AuthConstant",
          "metadata": {
            "class": "auth.constant.AuthConstant",
            "method": "AuthConstant"
          }
        },
        {
          "value": "auth.constant.InfoConstant#InfoConstant",
          "label": "auth.constant.InfoConstant#InfoConstant",
          "metadata": {
            "class": "auth.constant.InfoConstant",
            "method": "InfoConstant"
          }
        },
        {
          "value": "auth.controller.AuthController#createDefaultUser",
          "label": "auth.controller.AuthController#createDefaultUser",
          "metadata": {
            "class": "auth.controller.AuthController",
            "method": "createDefaultUser"
          }
        },
        {
          "value": "auth.controller.AuthController#getHello",
          "label": "auth.controller.AuthController#getHello",
          "metadata": {
            "class": "auth.controller.AuthController",
            "method": "getHello"
          }
        },
        {
          "value": "auth.controller.UserController#deleteUserById",
          "label": "auth.controller.UserController#deleteUserById",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "deleteUserById"
          }
        },
        {
          "value": "auth.controller.UserController#getAllUser",
          "label": "auth.controller.UserController#getAllUser",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getAllUser"
          }
        },
        {
          "value": "auth.controller.UserController#getHello",
          "label": "auth.controller.UserController#getHello",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getHello"
          }
        },
        {
          "value": "auth.controller.UserController#getToken",
          "label": "auth.controller.UserController#getToken",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getToken"
          }
        },
        {
          "value": "auth.entity.User#getAuthorities",
          "label": "auth.entity.User#getAuthorities",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getAuthorities"
          }
        },
        {
          "value": "auth.entity.User#getPassword",
          "label": "auth.entity.User#getPassword",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getPassword"
          }
        },
        {
          "value": "auth.entity.User#getUsername",
          "label": "auth.entity.User#getUsername",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getUsername"
          }
        },
        {
          "value": "auth.entity.User#isAccountNonExpired",
          "label": "auth.entity.User#isAccountNonExpired",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isAccountNonExpired"
          }
        },
        {
          "value": "auth.entity.User#isAccountNonLocked",
          "label": "auth.entity.User#isAccountNonLocked",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isAccountNonLocked"
          }
        },
        {
          "value": "auth.entity.User#isCredentialsNonExpired",
          "label": "auth.entity.User#isCredentialsNonExpired",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isCredentialsNonExpired"
          }
        },
        {
          "value": "auth.entity.User#isEnabled",
          "label": "auth.entity.User#isEnabled",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isEnabled"
          }
        },
        {
          "value": "auth.exception.UserOperationException#UserOperationException",
          "label": "auth.exception.UserOperationException#UserOperationException",
          "metadata": {
            "class": "auth.exception.UserOperationException",
            "method": "UserOperationException"
          }
        },
        {
          "value": "auth.exception.handler.GlobalExceptionHandler#handleUserNotFoundException",
          "label": "auth.exception.handler.GlobalExceptionHandler#handleUserNotFoundException",
          "metadata": {
            "class": "auth.exception.handler.GlobalExceptionHandler",
            "method": "handleUserNotFoundException"
          }
        },
        {
          "value": "auth.init.InitUser#run",
          "label": "auth.init.InitUser#run",
          "metadata": {
            "class": "auth.init.InitUser",
            "method": "run"
          }
        },
        {
          "value": "auth.security.UserDetailsServiceImpl#loadUserByUsername",
          "label": "auth.security.UserDetailsServiceImpl#loadUserByUsername",
          "metadata": {
            "class": "auth.security.UserDetailsServiceImpl",
            "method": "loadUserByUsername"
          }
        },
        {
          "value": "auth.security.jwt.JWTProvider#createToken",
          "label": "auth.security.jwt.JWTProvider#createToken",
          "metadata": {
            "class": "auth.security.jwt.JWTProvider",
            "method": "createToken"
          }
        },
        {
          "value": "auth.security.jwt.JWTProvider#init",
          "label": "auth.security.jwt.JWTProvider#init",
          "metadata": {
            "class": "auth.security.jwt.JWTProvider",
            "method": "init"
          }
        },
        {
          "value": "auth.service.impl.TokenServiceImpl#getServiceUrl",
          "label": "auth.service.impl.TokenServiceImpl#getServiceUrl",
          "metadata": {
            "class": "auth.service.impl.TokenServiceImpl",
            "method": "getServiceUrl"
          }
        },
        {
          "value": "auth.service.impl.TokenServiceImpl#getToken",
          "label": "auth.service.impl.TokenServiceImpl#getToken",
          "metadata": {
            "class": "auth.service.impl.TokenServiceImpl",
            "method": "getToken"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#checkUserCreateInfo",
          "label": "auth.service.impl.UserServiceImpl#checkUserCreateInfo",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "checkUserCreateInfo"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#createDefaultAuthUser",
          "label": "auth.service.impl.UserServiceImpl#createDefaultAuthUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "createDefaultAuthUser"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#deleteByUserId",
          "label": "auth.service.impl.UserServiceImpl#deleteByUserId",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "deleteByUserId"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#getAllUser",
          "label": "auth.service.impl.UserServiceImpl#getAllUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "getAllUser"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#saveUser",
          "label": "auth.service.impl.UserServiceImpl#saveUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "saveUser"
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

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMCPUStress --class auth.security.jwt.JWTProvider --method createToken
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMCPUStress",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMCPUStress",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill JVM CPU stress parameters",
      "fields": [
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
        },
        {
          "name": "cpu_count",
          "kind": "number_range",
          "required": true,
          "description": "CPU core count",
          "min": 1,
          "max": 8,
          "step": 1
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "cpu_count is required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMCPUStress --class auth.security.jwt.JWTProvider --method createToken --cpu-count 1
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
    "chaos_type": "JVMCPUStress",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken",
    "duration": 5,
    "cpu_count": 1
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMCPUStress",
    "class": "auth.security.jwt.JWTProvider",
    "cpu_count": 1,
    "duration": 5,
    "method": "createToken",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill JVM CPU stress parameters",
      "fields": [
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
        },
        {
          "name": "cpu_count",
          "kind": "number_range",
          "required": true,
          "description": "CPU core count",
          "min": 1,
          "max": 8,
          "step": 1
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "JVMCPUStress": {
          "CPUCount": 1,
          "Duration": 5,
          "MethodIdx": 132,
          "System": 7
        }
      },
      "chaos_type": "JVMCPUStress",
      "cpu_count": 1,
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "class_name": "auth.security.jwt.JWTProvider",
        "method_name": "createToken"
      },
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "function": [
        "auth.security.jwt.JWTProvider.createToken"
      ],
      "metric": [
        "cpu"
      ],
      "service": [
        "ts-auth-service"
      ]
    }
  },
  "apply_payload": {
    "JVMCPUStress": {
      "CPUCount": 1,
      "Duration": 5,
      "MethodIdx": 132,
      "System": 7
    }
  },
  "can_apply": true
}
```

## JVMMemoryStress

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMMemoryStress
```
```json
{
  "mode": "guided",
  "stage": "select_jvm_method",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMMemoryStress"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMMemoryStress",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "method_ref",
      "kind": "object_ref",
      "required": true,
      "description": "Select the JVM method for memory stress injection",
      "options": [
        {
          "value": "auth.AuthApplication#main",
          "label": "auth.AuthApplication#main",
          "metadata": {
            "class": "auth.AuthApplication",
            "method": "main"
          }
        },
        {
          "value": "auth.constant.AuthConstant#AuthConstant",
          "label": "auth.constant.AuthConstant#AuthConstant",
          "metadata": {
            "class": "auth.constant.AuthConstant",
            "method": "AuthConstant"
          }
        },
        {
          "value": "auth.constant.InfoConstant#InfoConstant",
          "label": "auth.constant.InfoConstant#InfoConstant",
          "metadata": {
            "class": "auth.constant.InfoConstant",
            "method": "InfoConstant"
          }
        },
        {
          "value": "auth.controller.AuthController#createDefaultUser",
          "label": "auth.controller.AuthController#createDefaultUser",
          "metadata": {
            "class": "auth.controller.AuthController",
            "method": "createDefaultUser"
          }
        },
        {
          "value": "auth.controller.AuthController#getHello",
          "label": "auth.controller.AuthController#getHello",
          "metadata": {
            "class": "auth.controller.AuthController",
            "method": "getHello"
          }
        },
        {
          "value": "auth.controller.UserController#deleteUserById",
          "label": "auth.controller.UserController#deleteUserById",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "deleteUserById"
          }
        },
        {
          "value": "auth.controller.UserController#getAllUser",
          "label": "auth.controller.UserController#getAllUser",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getAllUser"
          }
        },
        {
          "value": "auth.controller.UserController#getHello",
          "label": "auth.controller.UserController#getHello",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getHello"
          }
        },
        {
          "value": "auth.controller.UserController#getToken",
          "label": "auth.controller.UserController#getToken",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getToken"
          }
        },
        {
          "value": "auth.entity.User#getAuthorities",
          "label": "auth.entity.User#getAuthorities",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getAuthorities"
          }
        },
        {
          "value": "auth.entity.User#getPassword",
          "label": "auth.entity.User#getPassword",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getPassword"
          }
        },
        {
          "value": "auth.entity.User#getUsername",
          "label": "auth.entity.User#getUsername",
          "metadata": {
            "class": "auth.entity.User",
            "method": "getUsername"
          }
        },
        {
          "value": "auth.entity.User#isAccountNonExpired",
          "label": "auth.entity.User#isAccountNonExpired",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isAccountNonExpired"
          }
        },
        {
          "value": "auth.entity.User#isAccountNonLocked",
          "label": "auth.entity.User#isAccountNonLocked",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isAccountNonLocked"
          }
        },
        {
          "value": "auth.entity.User#isCredentialsNonExpired",
          "label": "auth.entity.User#isCredentialsNonExpired",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isCredentialsNonExpired"
          }
        },
        {
          "value": "auth.entity.User#isEnabled",
          "label": "auth.entity.User#isEnabled",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isEnabled"
          }
        },
        {
          "value": "auth.exception.UserOperationException#UserOperationException",
          "label": "auth.exception.UserOperationException#UserOperationException",
          "metadata": {
            "class": "auth.exception.UserOperationException",
            "method": "UserOperationException"
          }
        },
        {
          "value": "auth.exception.handler.GlobalExceptionHandler#handleUserNotFoundException",
          "label": "auth.exception.handler.GlobalExceptionHandler#handleUserNotFoundException",
          "metadata": {
            "class": "auth.exception.handler.GlobalExceptionHandler",
            "method": "handleUserNotFoundException"
          }
        },
        {
          "value": "auth.init.InitUser#run",
          "label": "auth.init.InitUser#run",
          "metadata": {
            "class": "auth.init.InitUser",
            "method": "run"
          }
        },
        {
          "value": "auth.security.UserDetailsServiceImpl#loadUserByUsername",
          "label": "auth.security.UserDetailsServiceImpl#loadUserByUsername",
          "metadata": {
            "class": "auth.security.UserDetailsServiceImpl",
            "method": "loadUserByUsername"
          }
        },
        {
          "value": "auth.security.jwt.JWTProvider#createToken",
          "label": "auth.security.jwt.JWTProvider#createToken",
          "metadata": {
            "class": "auth.security.jwt.JWTProvider",
            "method": "createToken"
          }
        },
        {
          "value": "auth.security.jwt.JWTProvider#init",
          "label": "auth.security.jwt.JWTProvider#init",
          "metadata": {
            "class": "auth.security.jwt.JWTProvider",
            "method": "init"
          }
        },
        {
          "value": "auth.service.impl.TokenServiceImpl#getServiceUrl",
          "label": "auth.service.impl.TokenServiceImpl#getServiceUrl",
          "metadata": {
            "class": "auth.service.impl.TokenServiceImpl",
            "method": "getServiceUrl"
          }
        },
        {
          "value": "auth.service.impl.TokenServiceImpl#getToken",
          "label": "auth.service.impl.TokenServiceImpl#getToken",
          "metadata": {
            "class": "auth.service.impl.TokenServiceImpl",
            "method": "getToken"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#checkUserCreateInfo",
          "label": "auth.service.impl.UserServiceImpl#checkUserCreateInfo",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "checkUserCreateInfo"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#createDefaultAuthUser",
          "label": "auth.service.impl.UserServiceImpl#createDefaultAuthUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "createDefaultAuthUser"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#deleteByUserId",
          "label": "auth.service.impl.UserServiceImpl#deleteByUserId",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "deleteByUserId"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#getAllUser",
          "label": "auth.service.impl.UserServiceImpl#getAllUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "getAllUser"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#saveUser",
          "label": "auth.service.impl.UserServiceImpl#saveUser",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "saveUser"
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

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMMemoryStress --class auth.security.jwt.JWTProvider --method createToken
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMMemoryStress",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMMemoryStress",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill JVM memory stress parameters",
      "fields": [
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
        },
        {
          "name": "mem_type",
          "kind": "enum",
          "required": true,
          "description": "Memory type",
          "options": [
            {
              "value": "heap",
              "label": "heap"
            },
            {
              "value": "stack",
              "label": "stack"
            }
          ]
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "mem_type is required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMMemoryStress --class auth.security.jwt.JWTProvider --method createToken --mem-type heap
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
    "chaos_type": "JVMMemoryStress",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken",
    "duration": 5,
    "mem_type": "heap"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMMemoryStress",
    "class": "auth.security.jwt.JWTProvider",
    "duration": 5,
    "mem_type": "heap",
    "method": "createToken",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill JVM memory stress parameters",
      "fields": [
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
        },
        {
          "name": "mem_type",
          "kind": "enum",
          "required": true,
          "description": "Memory type",
          "options": [
            {
              "value": "heap",
              "label": "heap"
            },
            {
              "value": "stack",
              "label": "stack"
            }
          ]
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "JVMMemoryStress": {
          "Duration": 5,
          "MemType": 1,
          "MethodIdx": 132,
          "System": 7
        }
      },
      "chaos_type": "JVMMemoryStress",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "class_name": "auth.security.jwt.JWTProvider",
        "method_name": "createToken"
      },
      "mem_type": "heap",
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "function": [
        "auth.security.jwt.JWTProvider.createToken"
      ],
      "metric": [
        "memory"
      ],
      "service": [
        "ts-auth-service"
      ]
    }
  },
  "apply_payload": {
    "JVMMemoryStress": {
      "Duration": 5,
      "MemType": 1,
      "MethodIdx": 132,
      "System": 7
    }
  },
  "can_apply": true
}
```

## JVMMySQLLatency

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMMySQLLatency
```
```json
{
  "mode": "guided",
  "stage": "select_database_operation",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMMySQLLatency"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMMySQLLatency",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "database_ref",
      "kind": "object_ref",
      "required": true,
      "description": "Select the database/table/operation tuple for MySQL chaos",
      "options": [
        {
          "value": "ts//ALTER table",
          "label": "ts /  / ALTER table",
          "metadata": {
            "database": "ts",
            "operation": "ALTER table",
            "table": ""
          }
        },
        {
          "value": "ts//CREATE table",
          "label": "ts /  / CREATE table",
          "metadata": {
            "database": "ts",
            "operation": "CREATE table",
            "table": ""
          }
        },
        {
          "value": "ts//SELECT",
          "label": "ts /  / SELECT",
          "metadata": {
            "database": "ts",
            "operation": "SELECT",
            "table": ""
          }
        },
        {
          "value": "ts/`ts`.`auth_user`/SELECT",
          "label": "ts / `ts`.`auth_user` / SELECT",
          "metadata": {
            "database": "ts",
            "operation": "SELECT",
            "table": "`ts`.`auth_user`"
          }
        },
        {
          "value": "ts/`ts`.`user_roles`/CREATE TABLE",
          "label": "ts / `ts`.`user_roles` / CREATE TABLE",
          "metadata": {
            "database": "ts",
            "operation": "CREATE TABLE",
            "table": "`ts`.`user_roles`"
          }
        },
        {
          "value": "ts/auth_user/INSERT",
          "label": "ts / auth_user / INSERT",
          "metadata": {
            "database": "ts",
            "operation": "INSERT",
            "table": "auth_user"
          }
        },
        {
          "value": "ts/auth_user/SELECT",
          "label": "ts / auth_user / SELECT",
          "metadata": {
            "database": "ts",
            "operation": "SELECT",
            "table": "auth_user"
          }
        },
        {
          "value": "ts/ts/SELECT",
          "label": "ts / ts / SELECT",
          "metadata": {
            "database": "ts",
            "operation": "SELECT",
            "table": "ts"
          }
        },
        {
          "value": "ts/user_roles/INSERT",
          "label": "ts / user_roles / INSERT",
          "metadata": {
            "database": "ts",
            "operation": "INSERT",
            "table": "user_roles"
          }
        },
        {
          "value": "ts/user_roles/SELECT",
          "label": "ts / user_roles / SELECT",
          "metadata": {
            "database": "ts",
            "operation": "SELECT",
            "table": "user_roles"
          }
        }
      ],
      "key_fields": [
        "database",
        "table",
        "operation"
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMMySQLLatency --database ts --table auth_user --operation SELECT
```
```json
{
  "mode": "guided",
  "stage": "fill_required_fields",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMMySQLLatency",
    "database": "ts",
    "table": "auth_user",
    "operation": "SELECT"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMMySQLLatency",
    "database": "ts",
    "namespace": "ts",
    "operation": "SELECT",
    "system": "ts",
    "system_type": "ts",
    "table": "auth_user"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill JVM MySQL latency parameters",
      "fields": [
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
        },
        {
          "name": "latency_ms",
          "kind": "number_range",
          "required": true,
          "description": "MySQL latency",
          "min": 10,
          "max": 5000,
          "step": 1,
          "unit": "ms"
        }
      ]
    }
  ],
  "can_apply": false,
  "errors": [
    "latency_ms is required"
  ]
}
```

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMMySQLLatency --database ts --table auth_user --operation SELECT --latency-ms 200
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
    "chaos_type": "JVMMySQLLatency",
    "database": "ts",
    "table": "auth_user",
    "operation": "SELECT",
    "duration": 5,
    "latency_ms": 200
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMMySQLLatency",
    "database": "ts",
    "duration": 5,
    "latency_ms": 200,
    "namespace": "ts",
    "operation": "SELECT",
    "system": "ts",
    "system_type": "ts",
    "table": "auth_user"
  },
  "next": [
    {
      "name": "params",
      "kind": "group",
      "required": true,
      "description": "Fill JVM MySQL latency parameters",
      "fields": [
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
        },
        {
          "name": "latency_ms",
          "kind": "number_range",
          "required": true,
          "description": "MySQL latency",
          "min": 10,
          "max": 5000,
          "step": 1,
          "unit": "ms"
        }
      ]
    }
  ],
  "preview": {
    "display_config": {
      "apply_payload": {
        "JVMMySQLLatency": {
          "DatabaseIdx": 15,
          "Duration": 5,
          "LatencyMs": 200,
          "System": 7
        }
      },
      "chaos_type": "JVMMySQLLatency",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "database": "ts",
        "operation": "SELECT",
        "table": "auth_user"
      },
      "latency_ms": 200,
      "namespace": "ts",
      "system": "ts"
    },
    "groundtruth": {
      "metric": [
        "sql_latency"
      ],
      "service": [
        "ts-auth-service"
      ]
    }
  },
  "apply_payload": {
    "JVMMySQLLatency": {
      "DatabaseIdx": 15,
      "Duration": 5,
      "LatencyMs": 200,
      "System": 7
    }
  },
  "can_apply": true
}
```

## JVMMySQLException

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMMySQLException
```
```json
{
  "mode": "guided",
  "stage": "select_database_operation",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMMySQLException"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMMySQLException",
    "namespace": "ts",
    "system": "ts",
    "system_type": "ts"
  },
  "next": [
    {
      "name": "database_ref",
      "kind": "object_ref",
      "required": true,
      "description": "Select the database/table/operation tuple for MySQL chaos",
      "options": [
        {
          "value": "ts//ALTER table",
          "label": "ts /  / ALTER table",
          "metadata": {
            "database": "ts",
            "operation": "ALTER table",
            "table": ""
          }
        },
        {
          "value": "ts//CREATE table",
          "label": "ts /  / CREATE table",
          "metadata": {
            "database": "ts",
            "operation": "CREATE table",
            "table": ""
          }
        },
        {
          "value": "ts//SELECT",
          "label": "ts /  / SELECT",
          "metadata": {
            "database": "ts",
            "operation": "SELECT",
            "table": ""
          }
        },
        {
          "value": "ts/`ts`.`auth_user`/SELECT",
          "label": "ts / `ts`.`auth_user` / SELECT",
          "metadata": {
            "database": "ts",
            "operation": "SELECT",
            "table": "`ts`.`auth_user`"
          }
        },
        {
          "value": "ts/`ts`.`user_roles`/CREATE TABLE",
          "label": "ts / `ts`.`user_roles` / CREATE TABLE",
          "metadata": {
            "database": "ts",
            "operation": "CREATE TABLE",
            "table": "`ts`.`user_roles`"
          }
        },
        {
          "value": "ts/auth_user/INSERT",
          "label": "ts / auth_user / INSERT",
          "metadata": {
            "database": "ts",
            "operation": "INSERT",
            "table": "auth_user"
          }
        },
        {
          "value": "ts/auth_user/SELECT",
          "label": "ts / auth_user / SELECT",
          "metadata": {
            "database": "ts",
            "operation": "SELECT",
            "table": "auth_user"
          }
        },
        {
          "value": "ts/ts/SELECT",
          "label": "ts / ts / SELECT",
          "metadata": {
            "database": "ts",
            "operation": "SELECT",
            "table": "ts"
          }
        },
        {
          "value": "ts/user_roles/INSERT",
          "label": "ts / user_roles / INSERT",
          "metadata": {
            "database": "ts",
            "operation": "INSERT",
            "table": "user_roles"
          }
        },
        {
          "value": "ts/user_roles/SELECT",
          "label": "ts / user_roles / SELECT",
          "metadata": {
            "database": "ts",
            "operation": "SELECT",
            "table": "user_roles"
          }
        }
      ],
      "key_fields": [
        "database",
        "table",
        "operation"
      ]
    }
  ],
  "can_apply": false
}
```

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMMySQLException --database ts --table auth_user --operation SELECT
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
    "chaos_type": "JVMMySQLException",
    "database": "ts",
    "table": "auth_user",
    "operation": "SELECT",
    "duration": 5
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMMySQLException",
    "database": "ts",
    "duration": 5,
    "namespace": "ts",
    "operation": "SELECT",
    "system": "ts",
    "system_type": "ts",
    "table": "auth_user"
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
        "JVMMySQLException": {
          "DatabaseIdx": 15,
          "Duration": 5,
          "System": 7
        }
      },
      "chaos_type": "JVMMySQLException",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "database": "ts",
        "operation": "SELECT",
        "table": "auth_user"
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
    "JVMMySQLException": {
      "DatabaseIdx": 15,
      "Duration": 5,
      "System": 7
    }
  },
  "can_apply": true
}
```

## JVMRuntimeMutator

### Call 1

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMRuntimeMutator
```
```json
{
  "mode": "guided",
  "stage": "select_runtime_mutator_method",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMRuntimeMutator"
  },
  "resolved": {
    "app": "ts-auth-service",
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
          "value": "auth.constant.AuthConstant#AuthConstant",
          "label": "auth.constant.AuthConstant#AuthConstant",
          "metadata": {
            "class": "auth.constant.AuthConstant",
            "method": "AuthConstant"
          }
        },
        {
          "value": "auth.constant.InfoConstant#InfoConstant",
          "label": "auth.constant.InfoConstant#InfoConstant",
          "metadata": {
            "class": "auth.constant.InfoConstant",
            "method": "InfoConstant"
          }
        },
        {
          "value": "auth.controller.AuthController#createDefaultUser",
          "label": "auth.controller.AuthController#createDefaultUser",
          "metadata": {
            "class": "auth.controller.AuthController",
            "method": "createDefaultUser"
          }
        },
        {
          "value": "auth.controller.AuthController#getHello",
          "label": "auth.controller.AuthController#getHello",
          "metadata": {
            "class": "auth.controller.AuthController",
            "method": "getHello"
          }
        },
        {
          "value": "auth.controller.UserController#deleteUserById",
          "label": "auth.controller.UserController#deleteUserById",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "deleteUserById"
          }
        },
        {
          "value": "auth.controller.UserController#getAllUser",
          "label": "auth.controller.UserController#getAllUser",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getAllUser"
          }
        },
        {
          "value": "auth.controller.UserController#getHello",
          "label": "auth.controller.UserController#getHello",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getHello"
          }
        },
        {
          "value": "auth.controller.UserController#getToken",
          "label": "auth.controller.UserController#getToken",
          "metadata": {
            "class": "auth.controller.UserController",
            "method": "getToken"
          }
        },
        {
          "value": "auth.entity.User#isAccountNonExpired",
          "label": "auth.entity.User#isAccountNonExpired",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isAccountNonExpired"
          }
        },
        {
          "value": "auth.entity.User#isAccountNonLocked",
          "label": "auth.entity.User#isAccountNonLocked",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isAccountNonLocked"
          }
        },
        {
          "value": "auth.entity.User#isCredentialsNonExpired",
          "label": "auth.entity.User#isCredentialsNonExpired",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isCredentialsNonExpired"
          }
        },
        {
          "value": "auth.entity.User#isEnabled",
          "label": "auth.entity.User#isEnabled",
          "metadata": {
            "class": "auth.entity.User",
            "method": "isEnabled"
          }
        },
        {
          "value": "auth.init.InitUser#run",
          "label": "auth.init.InitUser#run",
          "metadata": {
            "class": "auth.init.InitUser",
            "method": "run"
          }
        },
        {
          "value": "auth.security.jwt.JWTProvider#createToken",
          "label": "auth.security.jwt.JWTProvider#createToken",
          "metadata": {
            "class": "auth.security.jwt.JWTProvider",
            "method": "createToken"
          }
        },
        {
          "value": "auth.service.impl.TokenServiceImpl#getServiceUrl",
          "label": "auth.service.impl.TokenServiceImpl#getServiceUrl",
          "metadata": {
            "class": "auth.service.impl.TokenServiceImpl",
            "method": "getServiceUrl"
          }
        },
        {
          "value": "auth.service.impl.TokenServiceImpl#getToken",
          "label": "auth.service.impl.TokenServiceImpl#getToken",
          "metadata": {
            "class": "auth.service.impl.TokenServiceImpl",
            "method": "getToken"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#checkUserCreateInfo",
          "label": "auth.service.impl.UserServiceImpl#checkUserCreateInfo",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "checkUserCreateInfo"
          }
        },
        {
          "value": "auth.service.impl.UserServiceImpl#deleteByUserId",
          "label": "auth.service.impl.UserServiceImpl#deleteByUserId",
          "metadata": {
            "class": "auth.service.impl.UserServiceImpl",
            "method": "deleteByUserId"
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

### Call 2

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMRuntimeMutator --class auth.security.jwt.JWTProvider --method createToken
```
```json
{
  "mode": "guided",
  "stage": "select_runtime_mutator_config",
  "config": {
    "system": "ts",
    "system_type": "ts",
    "namespace": "ts",
    "app": "ts-auth-service",
    "chaos_type": "JVMRuntimeMutator",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken"
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMRuntimeMutator",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken",
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
          "value": "operator:add_to_sub",
          "label": "Mutate + to -",
          "metadata": {
            "description": "Mutate + to -",
            "mutation_from": "",
            "mutation_strategy": "add_to_sub",
            "mutation_to": "",
            "mutation_type_name": "operator"
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

### Call 3

```bash
go run ./cmd/chaos-exp -output json --namespace ts --app ts-auth-service --chaos-type JVMRuntimeMutator --class auth.security.jwt.JWTProvider --method createToken --mutator-config 'operator:add_to_sub'
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
    "chaos_type": "JVMRuntimeMutator",
    "class": "auth.security.jwt.JWTProvider",
    "method": "createToken",
    "mutator_config": "operator:add_to_sub",
    "duration": 5
  },
  "resolved": {
    "app": "ts-auth-service",
    "chaos_type": "JVMRuntimeMutator",
    "class": "auth.security.jwt.JWTProvider",
    "duration": 5,
    "method": "createToken",
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
          "MutatorTargetIdx": 1006,
          "System": 7
        }
      },
      "chaos_type": "JVMRuntimeMutator",
      "duration": 5,
      "injection_point": {
        "app_name": "ts-auth-service",
        "class_name": "auth.security.jwt.JWTProvider",
        "description": "Mutate + to -",
        "method_name": "createToken",
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
        "auth.security.jwt.JWTProvider.createToken"
      ],
      "service": [
        "ts-auth-service"
      ]
    }
  },
  "apply_payload": {
    "JVMRuntimeMutator": {
      "Duration": 5,
      "MutatorTargetIdx": 1006,
      "System": 7
    }
  },
  "can_apply": true
}
```


## Session Shorthand Validation (2026-04-17)

This section captures the live validation run for the new guided-session shorthand behavior using a temporary session file in `%TEMP%`.

### Environment

- binary: `./chaos-exp.exe`
- namespace used for validation: `ts`
- session file: `%TEMP%\\chaos-exp-guided-session.yaml`

### 1. Start from Namespace

```bash
./chaos-exp.exe -output json --config %TEMP%\\chaos-exp-guided-session.yaml --reset-config --namespace ts
```

Observed result:

- `stage`: `select_app`
- `config`: `{"system":"ts","system_type":"ts","namespace":"ts"}`
- `next`: `app`

### 2. Use `--next` for App Selection

```bash
./chaos-exp.exe -output json --config %TEMP%\\chaos-exp-guided-session.yaml --next ts-auth-service
```

Observed result:

- `stage`: `select_chaos_type`
- `config`: `{"system":"ts","system_type":"ts","namespace":"ts","app":"ts-auth-service"}`
- `next`: `chaos_type`

### 3. Use `--next` for Chaos Type Selection

```bash
./chaos-exp.exe -output json --config %TEMP%\\chaos-exp-guided-session.yaml --next PodKill
```

Observed result:

- `stage`: `ready_to_apply`
- `config`: `{"system":"ts","system_type":"ts","namespace":"ts","app":"ts-auth-service","chaos_type":"PodKill","duration":5}`
- `next`: `duration`
- `can_apply`: `true`

### 4. Change App Directly and Verify Downstream Clearing

```bash
./chaos-exp.exe -output json --config %TEMP%\\chaos-exp-guided-session.yaml --app ts-delivery-service
```

Observed result:

- `stage`: `select_chaos_type`
- `config`: `{"system":"ts","system_type":"ts","namespace":"ts","app":"ts-delivery-service","duration":5}`
- previous `chaos_type=PodKill` was cleared automatically

### 5. Continue with Direct Stage Flag

```bash
./chaos-exp.exe -output json --config %TEMP%\\chaos-exp-guided-session.yaml --chaos-type JVMRuntimeMutator
```

Observed result:

- `stage`: `select_runtime_mutator_method`
- `config`: `{"system":"ts","system_type":"ts","namespace":"ts","app":"ts-delivery-service","chaos_type":"JVMRuntimeMutator","duration":5}`
- `next`: `method_ref`

### 6. Use `--next` for Object Ref Selection

```bash
./chaos-exp.exe -output json --config %TEMP%\\chaos-exp-guided-session.yaml --next delivery.mq.RabbitReceive#process
```

Observed result:

- `stage`: `select_runtime_mutator_config`
- `config`: `{"system":"ts","system_type":"ts","namespace":"ts","app":"ts-delivery-service","chaos_type":"JVMRuntimeMutator","class":"delivery.mq.RabbitReceive","method":"process","duration":5}`
- `next`: `mutator_config`, `duration`

### 7. Use `--next` When the Stage Has One Required Selector Plus Optional Duration

```bash
./chaos-exp.exe -output json --config %TEMP%\\chaos-exp-guided-session.yaml --next operator:add_to_sub
```

Observed result:

- `stage`: `ready_to_apply`
- `config`: `{"system":"ts","system_type":"ts","namespace":"ts","app":"ts-delivery-service","chaos_type":"JVMRuntimeMutator","class":"delivery.mq.RabbitReceive","method":"process","mutator_config":"operator:add_to_sub","duration":5}`
- `can_apply`: `true`

### Conclusion

The live run confirms the new behavior works as intended:

- saved guided session state auto-fills earlier stages
- `--next` advances single-selector stages, including object refs
- direct stage flags such as `--app` and `--chaos-type` work against the saved session
- changing an earlier stage clears stale downstream selections

### 8. Validate Grouped Parameter Input with Direct Flags

A grouped stage was also validated with `CPUStress`.

#### 8.1 Enter the grouped stage

```bash
./chaos-exp.exe -output json --config %TEMP%\\chaos-exp-guided-session-cpu.yaml --reset-config --namespace ts
./chaos-exp.exe -output json --config %TEMP%\\chaos-exp-guided-session-cpu.yaml --next ts-auth-service
./chaos-exp.exe -output json --config %TEMP%\\chaos-exp-guided-session-cpu.yaml --chaos-type CPUStress
./chaos-exp.exe -output json --config %TEMP%\\chaos-exp-guided-session-cpu.yaml --next ts-auth-service
```

Observed result after the last call:

- `stage`: `fill_required_fields`
- `config`: `{"system":"ts","system_type":"ts","namespace":"ts","app":"ts-auth-service","chaos_type":"CPUStress","container":"ts-auth-service"}`
- `next.kind`: `group`
- required fields: `cpu_load`, `cpu_worker`

#### 8.2 Fill the group with explicit flags

```bash
./chaos-exp.exe -output json --config %TEMP%\\chaos-exp-guided-session-cpu.yaml --cpu-load 80 --cpu-worker 1
```

Observed result:

- `stage`: `ready_to_apply`
- `config`: `{"system":"ts","system_type":"ts","namespace":"ts","app":"ts-auth-service","chaos_type":"CPUStress","container":"ts-auth-service","duration":5,"cpu_load":80,"cpu_worker":1}`
- `can_apply`: `true`

This confirms the intended split:

- `--next` for selector stages
- explicit flags for grouped parameter stages
