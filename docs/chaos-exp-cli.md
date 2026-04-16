# chaos-exp Guided CLI

## Overview

`chaos-exp` exposes a guided terminal workflow for building a chaos experiment through repeated, machine-readable calls.

Instead of requiring one fully flattened command up front, each call returns:

- the current reusable `config` snapshot
- the normalized `resolved` context
- the next selectable field or parameter group in `next`
- a preview and apply payload once the request reaches `ready_to_apply`

The intended client is an automation agent or terminal workflow that behaves like a cascading form.

## Core Interaction Model

The guided resolver is a state machine backed by the config file.

A normal session progresses like this:

1. select a system / namespace
2. select an app
3. select a chaos type
4. select any type-specific resource
5. fill any remaining parameter group
6. reach `ready_to_apply`
7. optionally re-run with `--apply`

Every response includes the full `config` snapshot for the current state, so the caller can always continue from the last response alone.

## Session Persistence

Guided sessions now persist by default.

Default config path:

```text
~/.chaos-exp/config.yaml
```

Unless disabled, every guided call saves the returned `config` snapshot into `guided-session.config`.

Relevant flags:

- `--config <path>`: use a specific session file
- `--save-config` : explicitly keep the default auto-save behavior
- `--no-save-config`: disable session persistence for the current call
- `--reset-config`: clear the saved guided session before resolving

This means later calls do not need to repeat already selected earlier stages.

Example:

```bash
chaos-exp -output json --namespace ts
chaos-exp -output json --app ts-auth-service
chaos-exp -output json --chaos-type PodKill
```

The second call reuses the saved namespace. The third call reuses both namespace and app.

## `--next` Shortcuts

`--next` applies a single next-step selection against the current saved or merged session state.

Typical usage:

```bash
chaos-exp -output json --namespace ts
chaos-exp -output json --next ts-auth-service
chaos-exp -output json --next PodKill
```

How it works:

- the CLI first resolves the current state
- it reads the current `next` contract
- it applies the supplied value to the one selectable field for that stage
- it resolves again and returns the next response

`--next` works best for single-choice stages such as:

- `system`
- `app`
- `chaos_type`
- `container`
- `target_service`
- `domain`
- JVM method selection (`Class#method`)
- HTTP endpoint selection (`METHOD /route`)
- database operation selection (`db/table/operation`)
- runtime mutator config
- standalone numeric fields such as `duration`

`--next` is intentionally not supported for grouped parameter stages. For those, use explicit flags such as:

```bash
chaos-exp -output json --latency 120 --correlation 50 --jitter 10 --direction both
```

## Direct-Flag Shorthand

Because the guided session is auto-saved, a later-stage flag can be supplied directly and earlier selections will be filled from the saved session.

Examples:

```bash
chaos-exp -output json --namespace ts
chaos-exp -output json --app ts-order-service
chaos-exp -output json --chaos-type NetworkDelay
chaos-exp -output json --target-service ts-station-service
chaos-exp -output json --latency 120 --correlation 50 --jitter 10 --direction both
```

This produces the same guided cascade as repeating all earlier flags on every call, but with much shorter commands.

## Automatic Downstream Clearing

When the caller changes an earlier-stage selection, the CLI automatically drops stale downstream state before continuing.

Current clearing rules:

- changing `system`, `system_type`, or `namespace` clears the previous guided branch and re-roots the session
- changing `app` clears `chaos_type` and all deeper type-specific selections and parameters
- changing `chaos_type` clears all post-type selections and type-specific parameters
- changing a JVM method selection clears `mutator_config`
- changing an HTTP endpoint selection clears `replace_method`
- changing a database tuple clears the saved database tuple fields before applying the new one

This prevents invalid combinations such as a mutator config being silently reused for a different method.

## System, Namespace, and App Resolution

The CLI accepts either `--system` or `--namespace`.

- `system` is the namespace instance name exposed to the user, for example `ts` or `ts0`
- `system_type` is the internal registered system family, for example `ts`
- `namespace` is the concrete Kubernetes namespace used for lookup and apply

Normalization examples:

- `--system ts0` resolves to `system=ts0`, `system_type=ts`, `namespace=ts0`
- `--namespace ts` resolves to `system=ts`, `system_type=ts`, `namespace=ts`
- if only `system_type` is present, the default namespace instance for that system is inferred

### App Label Filtering

The app picker is intentionally filtered.

Guided app candidates are limited to service names that appear as source services in network dependency metadata. This keeps the top-level app list aligned with dependency-aware chaos workflows and hides infrastructure-only labels that do not act as source applications.

## Response Contract

Every guided response follows the same structure.

```yaml
mode: guided
stage: string
config: object
resolved: object
next: []
preview: object?
apply_payload: object?
result: object?
can_apply: boolean
warnings: []
errors: []
resources: object?
meta: object?
```

Important fields:

- `mode`: always `guided`
- `stage`: the current point in the cascade, such as `select_app`, `select_http_endpoint`, or `ready_to_apply`
- `config`: the reusable guided session snapshot for the next call
- `resolved`: normalized values actually used by the resolver
- `next`: one or more field specifications describing the next selection or parameter input
- `preview`: a human-readable preview of the selected target and parameters
- `apply_payload`: the internal payload once the config is complete
- `can_apply`: whether the request is complete enough to create the chaos resource

### Field Kinds

The guided resolver currently emits these field kinds:

- `enum`: choose one value from a list
- `number_range`: choose a numeric value within a documented range
- `group`: fill multiple related parameters in one step
- `object_ref`: choose a structured resource such as an endpoint, JVM method, or database operation

## Supported Chaos Types

The guided CLI currently covers all chaos types exposed by the existing handler layer.

### App-Level

| Chaos Type | Selection Path | Required Parameters |
| --- | --- | --- |
| `PodKill` | `namespace -> app -> chaos_type` | optional `duration` |
| `PodFailure` | `namespace -> app -> chaos_type` | optional `duration` |
| `JVMGarbageCollector` | `namespace -> app -> chaos_type` | optional `duration` |

### Container-Level

| Chaos Type | Selection Path | Required Parameters |
| --- | --- | --- |
| `ContainerKill` | `namespace -> app -> container` | optional `duration` |
| `CPUStress` | `namespace -> app -> container` | `cpu_load`, `cpu_worker`, optional `duration` |
| `MemoryStress` | `namespace -> app -> container` | `memory_size`, `mem_worker`, optional `duration` |
| `TimeSkew` | `namespace -> app -> container` | `time_offset`, optional `duration` |

### HTTP

| Chaos Type | Selection Path | Required Parameters |
| --- | --- | --- |
| `HTTPRequestAbort` | `namespace -> app -> endpoint` | optional `duration` |
| `HTTPResponseAbort` | `namespace -> app -> endpoint` | optional `duration` |
| `HTTPRequestDelay` | `namespace -> app -> endpoint` | `delay_duration`, optional `duration` |
| `HTTPResponseDelay` | `namespace -> app -> endpoint` | `delay_duration`, optional `duration` |
| `HTTPResponseReplaceBody` | `namespace -> app -> endpoint` | `body_type`, optional `duration` |
| `HTTPResponsePatchBody` | `namespace -> app -> endpoint` | optional `duration` |
| `HTTPRequestReplacePath` | `namespace -> app -> endpoint` | optional `duration` |
| `HTTPRequestReplaceMethod` | `namespace -> app -> endpoint` | `replace_method`, optional `duration` |
| `HTTPResponseReplaceCode` | `namespace -> app -> endpoint` | `status_code`, optional `duration` |

### DNS

| Chaos Type | Selection Path | Required Parameters |
| --- | --- | --- |
| `DNSError` | `namespace -> app -> domain` | optional `duration` |
| `DNSRandom` | `namespace -> app -> domain` | optional `duration` |

### Network

| Chaos Type | Selection Path | Required Parameters |
| --- | --- | --- |
| `NetworkPartition` | `namespace -> app -> target_service` | `direction`, optional `duration` |
| `NetworkDelay` | `namespace -> app -> target_service` | `latency`, `correlation`, `jitter`, `direction`, optional `duration` |
| `NetworkLoss` | `namespace -> app -> target_service` | `loss`, `correlation`, `direction`, optional `duration` |
| `NetworkDuplicate` | `namespace -> app -> target_service` | `duplicate`, `correlation`, `direction`, optional `duration` |
| `NetworkCorrupt` | `namespace -> app -> target_service` | `corrupt`, `correlation`, `direction`, optional `duration` |
| `NetworkBandwidth` | `namespace -> app -> target_service` | `rate`, `limit`, `buffer`, `direction`, optional `duration` |

### JVM Method-Level

| Chaos Type | Selection Path | Required Parameters |
| --- | --- | --- |
| `JVMLatency` | `namespace -> app -> class+method` | `latency_duration`, optional `duration` |
| `JVMReturn` | `namespace -> app -> class+method` | `return_type`, `return_value_opt`, optional `duration` |
| `JVMException` | `namespace -> app -> class+method` | `exception_opt`, optional `duration` |
| `JVMCPUStress` | `namespace -> app -> class+method` | `cpu_count`, optional `duration` |
| `JVMMemoryStress` | `namespace -> app -> class+method` | `mem_type`, optional `duration` |

### Database-Level

| Chaos Type | Selection Path | Required Parameters |
| --- | --- | --- |
| `JVMMySQLLatency` | `namespace -> app -> database+table+operation` | `latency_ms`, optional `duration` |
| `JVMMySQLException` | `namespace -> app -> database+table+operation` | optional `duration` |

### Runtime Mutator

| Chaos Type | Selection Path | Required Parameters |
| --- | --- | --- |
| `JVMRuntimeMutator` | `namespace -> app -> class+method -> mutator_config` | optional `duration` |

## Duration Default

If no duration is provided, the guided resolver normalizes the request to a default `duration=5` minutes when building the final apply payload.

The caller may still override it with `--duration <minutes>` or by selecting the `duration` field at a later stage.

## Preview and Apply

Once all required selections and parameters are present, the response includes:

- `stage: ready_to_apply`
- `can_apply: true`
- `apply_payload`
- `preview.display_config`
- `preview.groundtruth` when available

Then the same command can be re-run with `--apply`.

Example:

```bash
chaos-exp -output json \
  --namespace ts \
  --app ts-auth-service \
  --chaos-type NetworkDelay \
  --target-service ts-station-service \
  --latency 120 \
  --correlation 50 \
  --jitter 10 \
  --direction both \
  --apply
```

## Current Limitations

A few guided types still inherit hard-coded values from the existing handler layer.

- `HTTPResponsePatchBody` uses the current built-in patch body payload
- `HTTPRequestReplacePath` uses the current built-in replacement path
- `JVMMySQLException` uses the current built-in exception content

These are handler limitations rather than guided resolver gaps.

## Related Documents

- `docs/chaos-exp-cli-examples.md`: example guided sessions using the new shorthand flow
- `docs/chaos-exp-agent-quickstart.md`: concise operating guide for agents
- `docs/guided-live-validation-2026-04-17.md`: raw live validation notes
