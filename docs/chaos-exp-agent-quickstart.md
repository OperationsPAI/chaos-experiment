# Agent Quickstart

This quickstart is for automation agents driving the guided `chaos-exp` CLI.

## 1. Always Ask for Machine-Readable Output

Use JSON unless a human specifically wants YAML.

```bash
chaos-exp -output json
```

Read these response fields first:

- `stage`
- `config`
- `resolved`
- `next`
- `can_apply`
- `errors`

## 2. Treat `config` as the Canonical Session Snapshot

Every response returns a reusable `config` object.

Two supported continuation styles now exist:

- rely on auto-saved guided session state and send only the new flag(s)
- or repeat fields manually if you want the call to be fully explicit

By default, guided responses are saved into `~/.chaos-exp/config.yaml`.

## 3. Prefer the Shortest Valid Follow-Up Call

Use one of these patterns:

### Single next-step selection

If `next` contains one required selector, use:

```bash
chaos-exp -output json --next <value>
```

Examples:

```bash
chaos-exp -output json --next ts-auth-service
chaos-exp -output json --next PodKill
chaos-exp -output json --next delivery.mq.RabbitReceive#process
chaos-exp -output json --next POST /api/v1/orders
chaos-exp -output json --next operator:add_to_sub
```

### Direct stage flag

If you already know the intended field, you can skip `--next` and set the field directly.

```bash
chaos-exp -output json --app ts-auth-service
chaos-exp -output json --chaos-type NetworkDelay
chaos-exp -output json --target-service ts-station-service
```

The current saved session fills the earlier stages automatically.

## 4. Use Explicit Flags for Grouped Parameter Stages

If `next[0].kind == "group"`, do not use `--next`.

Send the required parameter flags directly in one call.

Examples:

```bash
chaos-exp -output json --cpu-load 80 --cpu-worker 1
chaos-exp -output json --memory-size 256 --mem-worker 1
chaos-exp -output json --latency 120 --correlation 50 --jitter 10 --direction both
chaos-exp -output json --rate 1024 --limit 20971520 --buffer 10000 --direction both
```

Optional `duration` can be omitted. The guided resolver defaults it to `5` minutes.

## 5. Expect Automatic Downstream Clearing

When you change an earlier-stage field, the CLI clears stale deeper state automatically.

Practical consequences:

- changing `--namespace` or `--system` starts a new branch
- changing `--app` clears the previous `chaos_type` and deeper selections
- changing `--chaos-type` clears all post-type selections and parameters
- changing a JVM method clears `mutator_config`
- changing an HTTP endpoint clears `replace_method`

This protects against invalid state reuse.

## 6. Know the Common Object-Ref Formats

When `next.kind` is `object_ref`, the accepted shorthand formats are:

- HTTP endpoint: `<METHOD> <ROUTE>`
- JVM method: `<Class>#<method>`
- database operation: `<database>/<table>/<operation>`

Prefer using the exact `value` string from the returned option list.

## 7. Apply Only When `can_apply == true`

Once the response reaches:

- `stage: ready_to_apply`
- `can_apply: true`

review `preview` or `apply_payload`, then apply with:

```bash
chaos-exp -output json --apply
```

Because the session is already saved, no other flags are required unless you want to override something.

## 8. Reset When Switching Tasks

To discard the saved guided session:

```bash
chaos-exp -output json --reset-config
```

To isolate runs, use a dedicated config file:

```bash
chaos-exp -output json --config ./agent-session.yaml --namespace ts
```

## 9. Minimal Agent Loop

1. run `chaos-exp -output json ...`
2. parse `errors`; stop on non-empty errors
3. if `can_apply == true`, decide whether to call `--apply`
4. otherwise inspect `next`
5. if there is one required selector, prefer `--next <value>`
6. if the next input is a grouped parameter set, send explicit flags in one call
7. repeat until `ready_to_apply`
