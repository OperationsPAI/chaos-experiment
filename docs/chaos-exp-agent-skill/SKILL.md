---
name: chaos-exp-agent-skill
description: Use when an automation agent needs to drive the guided `chaos-exp` CLI, continue saved sessions, choose the next selector, fill grouped chaos parameters, and decide when a config is ready to apply.
---

# chaos-exp Agent Skill

Use this skill when you need to operate the guided `chaos-exp` CLI as a step-by-step terminal form.

## Quick start

Start with machine-readable output.

```bash
chaos-exp -output json
```

For a new task, prefer a clean session:

```bash
chaos-exp -output json --reset-config
```

Read these fields first from every response:

- `stage`
- `config`
- `resolved`
- `next`
- `can_apply`
- `errors`

Treat `config` as the canonical reusable session snapshot.

## Default operating loop

1. Run `chaos-exp -output json ...`.
2. Stop if `errors` is non-empty.
3. If `can_apply` is `true`, review `preview` or `apply_payload` before applying.
4. Otherwise inspect `next`.
5. If `next` contains one required selector, prefer `--next <value>`.
6. If the next input is a grouped parameter set, send explicit flags in one call.
7. Repeat until `stage` becomes `ready_to_apply`.

## Continuation rules

Use the shortest valid follow-up call.

### Single next-step selection

If there is one required selector, use:

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

If you already know the field, set it directly and let the saved session fill earlier stages.

```bash
chaos-exp -output json --app ts-auth-service
chaos-exp -output json --chaos-type NetworkDelay
chaos-exp -output json --target-service ts-station-service
```

## Grouped parameter stages

If `next[*].kind` is `group`, do not use `--next`.

Send the required parameter flags together in one call.

Examples:

```bash
chaos-exp -output json --cpu-load 80 --cpu-worker 1
chaos-exp -output json --memory-size 256 --mem-worker 1
chaos-exp -output json --latency 120 --correlation 50 --jitter 10 --direction both
chaos-exp -output json --rate 1024 --limit 20971520 --buffer 10000 --direction both
```

If you omit `--duration`, the resolver defaults it to `5` minutes.

## Session behavior

Guided responses are auto-saved in `~/.chaos-exp/config.yaml` unless disabled.

Practical consequences:

- You do not need to repeat earlier selections on every call.
- Changing `--system`, `--namespace`, `--app`, or `--chaos-type` automatically clears stale downstream state.
- Use `--reset-config` when switching to a different task.
- Use `--config <path>` to isolate a run.

## Object-ref formats

When `next.kind` is `object_ref`, use the exact returned `value` when possible.

Accepted shorthand formats:

- HTTP endpoint: `<METHOD> <ROUTE>`
- JVM method: `<Class>#<method>`
- database operation: `<database>/<table>/<operation>`

## Apply gate

Apply only after the response shows both:

- `stage: ready_to_apply`
- `can_apply: true`

Then run:

```bash
chaos-exp -output json --apply
```

Because the session is already saved, you usually do not need to repeat all flags.

## Do not guess values

- Prefer exact option values returned by `next`.
- Do not invent app names, method refs, mutator configs, or endpoint strings.
- If a grouped stage returns ranges, choose within those ranges.

## Read more only when needed

Open these docs only when the task needs more detail:

- `../chaos-exp-cli.md` for the full guided schema, stage model, and supported chaos types
- `../chaos-exp-cli-examples.md` for complete interaction examples
- `../guided-live-validation-2026-04-17.md` for recorded live validation traces
