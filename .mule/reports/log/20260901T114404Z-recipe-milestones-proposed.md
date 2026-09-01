slug: recipe-milestones
verdict: proposed
at:  3a70ddb
ran: 2026-09-01T11:44:04Z on DietPi in 607s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
I've read both files. Let me read the milestones file and check the recipe requirements.
→ Read .mule/milestones.md
→ Read .mule/recipes/astarte-upstream.md
Now I understand the recipe. Let me analyze the milestones file:

1. **v2.0 — astarte-flow feature parity** — `Status: not started.` This is the first section not marked `DONE`. So I work on v2.0 this run.
2. **v3.0** — TBD scope, but I shouldn't get there since v2.0 is the first un-DONE section.
3. **v4.0+** — placeholder.

The v2.0 scope is **not** marked TBD — it has a defined scope (astarte-flow feature parity). So I proceed with step 3 (investigate the gap).

Let me check for existing milestone-2.0 issues and look at Astrate's current state.

$ gh issue list --label milestone-2.0 --state open --limit 50; echo "---exit $?---"
---exit 0---

$ gh issue list --label milestone-2.0 --state all --limit 50
43	CLOSED	flow-container-block: PoC → MVP (custom processing in containers)	enhancement, milestone-2.0	2026-07-29T18:54:52Z
42	CLOSED	flow-rehydrate-edge-cases: follow-ups after durable auto_restart MVP	enhancement, milestone-2.0	2026-07-29T20:39:19Z
41	CLOSED	flow-durable: persist flows + auto_restart rehydrate on process boot	enhancement, milestone-2.0	2026-07-29T18:54:49Z
40	CLOSED	flow-multi-instance: named flows + pipeline config (decided — design then implement)	enhancement, question, milestone-2.0	2026-07-29T18:54:47Z
39	CLOSED	flow-blocks-discovery: GET /flow/v1/{realm}/blocks for registered catalog types	enhancement, milestone-2.0	2026-07-29T18:54:44Z
37	CLOSED	flow-source-pump: poll AstarteSource and Stop on flow teardown	milestone-2.0	2026-07-29T14:13:12Z
27	CLOSED	flow-astarte-source: Astarte device events source block	mule-review, milestone-2.0	2026-07-29T07:45:22Z
26	CLOSED	flow-router: Stream-based message routing through the block graph	milestone-2.0	2026-07-29T07:33:30Z
25	CLOSED	flow-lifecycle: Flow start/stop lifecycle and runtime management	milestone-2.0	2026-07-29T07:33:28Z
24	CLOSED	flow-pipeline-store: Pipeline CRUD storage and graph validation	mule-review, milestone-2.0	2026-07-29T07:45:19Z
23	CLOSED	flow-models: FlowMessage wire format and Block/Pipeline/Flow data models	milestone-2.0	2026-07-29T07:33:26Z

All `milestone-2.0` issues are **CLOSED**, and there are **no open** milestone-2.0 issues. The milestone says "not started" but there's a history of closed milestone-2.0 issues. Let me look at the actual reference (astarte_flow) and Astrate's current code to find any remaining gaps.

Let me investigate the astarte_flow reference and Astrate's current implementation.
$ git log --oneline -20
3a70ddb mule: refill the queue
22875f1 mule: recipe hygiene timed out
5006f6e mule: blocked docs-sync-pairing-status-enum
49fcd62 mule: refill the queue
45171b9 mule: recipe docs-sync ran (203s)
19ad0fe mule: blocked probe-properties-on-connect-encoding
8a983a3 mule: blocked probe-binaryblob-validation
25560b4 mule: blocked probe-mqtt-capabilities-declaration
21f74d8 mule: blocked compat-note-v133
2f072d3 mule: blocked probe-emptycache-resend-device-error
4c1228c mule: blocked compat-note-v140-rc3
acfcb26 mule: blocked probe-value-type-validation
db8048e mule: blocked probe-interface-default-values
eeb69d4 mule: blocked control-producer-properties-compression
d044364 mule: blocked issue-68
8e976a6 mule: blocked issue-91
7b03846 mule: recipe astarte-upstream ran (224s)
6e30a1c queue: drop three lines already done or duplicated, approve the plaintext-purge line, take the issues-first ordering from main
9233c4c mule: pin the model (from main)
6d1055b mule: take the cron survival fix from main (PATH + cwd)

$ ls internal/ | head -50
appengine
auth
broker
config
engine
flow
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil

$ cat docs/COMPATIBILITY.md 2>/dev/null | head -60
# Astrate Compatibility

Astrate is a wire-compatible re-implementation of Astarte: unmodified official
Astarte device SDKs and `astartectl` work against it. This document inventories
the **deliberate** deviations (`docs/DESIGN.md` §3.6) and lists the client
versions the conformance suite (`test/conformance`) pins and exercises.

## Supported / tested clients

These are pinned in `test/conformance/go.mod`, `test/conformance/pysdk/requirements.txt`,
and the astartectl helper; they are upgraded deliberately, never by drift
(`docs/ROADMAP.md` §0.3).

| Client | Pinned version | Checkpoint(s) |
|---|---|---|
| `astarte-device-sdk-go` | v0.90.2 | CP-B (M6), CP-D (M9 `gosdk`) |
| `astarte-go` (pairing/agent client) | v0.90.4 | CP-A, CP-B |
| `astartectl` (release binary) | v26.5.0 | CP-A, CP-C, CP-D regression |
| `astarte-device-sdk-python` | pinned in `pysdk/requirements.txt` | CP-D (`pysdk`) |
| AtomVM JSON profile (Astrate simulator) | n/a (this repo) | CP-D (`atomvm`) |

The conformance checkpoints gate the build at the earliest milestone that could
introduce protocol drift: **CP-A** (pairing, M4), **CP-B** (full device loop,
M6), **CP-C** (`astartectl` operator flow, M7), **CP-D** (full matrix, M9).

## Wire-identical surfaces

Guarded by the conformance suite against the official clients:

- MQTT topics, Astarte MQTT v1 connection contract (mTLS, identity from the
  certificate CN with the wire client ID free-form and remapped to the CN —
  the VerneMQ subscriber-id remap; session handling), and ACL model. The
  official Python SDK connects with a random paho client ID and relies on
  this.
- BSON `{v, t}` data documents and the zlib + size-prefixed control payloads
  (`emptyCache`, `producer/properties`, `consumer/properties`).
- Introspection format (`;`-separated `name:major:minor` triples).
- Pairing REST bodies and status codes (registration, credentials, info,
  verify).
- Certificate `Subject CN = <realm>/<device_id>`; 128-bit serials; clientAuth.
- JWT claim model (`a_aea`, `a_rma`, `a_pa`, `a_ha`, `a_ch`) with implicit
  anchoring and the `"<verb-regex>::<path-regex>"` authorization strings.
- Realm-management interface/trigger install/update/delete semantics and
  AppEngine device/data/query shapes, including body-`links` device-list
  pagination (`?details=true&limit&from_token` with the cursor in
  `links.next`) and `/stats/devices`.
- The Astarte Dashboard v1.2.2 runs unmodified against Astrate (compose
  `full` profile, `http://localhost:4040`), Device Live Events included since
  the M11 Channels socket (deviation 1).

## Deliberate deviations

All additive or strictly-safer; none affect unmodified device SDKs.

1. **Astarte Channels: two sockets, one bus** — the upstream Phoenix socket is
   served at `/appengine/v1/socket/websocket` (phoenix.js V2 wire format,
   `?vsn=2.0.0&realm=&token=`), which is what the Dashboard's Device Live
   Events card speaks; `phx_join`, `watch`, `phx_leave` and the `phoenix`
   heartbeat are answered, and matching events are pushed as `new_event`.
   Alongside it Astrate keeps its own plain WebSocket/SSE endpoint at

The milestone-2.0 issues are all closed, and there's a substantial `internal/flow/` directory. Let me investigate the current state of the flow implementation relative to the astarte_flow reference. Let me look at what's already done in internal/flow and what the reference actually provides.

Let me check the flow package and understand what capabilities exist, and look at the reference astarte_flow repo.
$ ls -la internal/flow/
total 100
drwxr-xr-x  2 root root  4096 Jul 28 04:11 .
drwxr-xr-x 16 root root  4096 Jul 28 02:09 ..
-rw-r--r--  1 root root  2868 Jul 28 04:11 block.go
-rw-r--r--  1 root root  6161 Jul 28 02:53 flow.go
-rw-r--r--  1 root root  4382 Jul 28 04:10 flow_test.go
-rw-r--r--  1 root root  1272 Jul 28 02:18 graph.go
-rw-r--r--  1 root root  7999 Jul 28 02:52 manager_test.go
-rw-r--r--  1 root root  9357 Jul 28 04:11 message.go
-rw-r--r--  1 root root  8287 Jul 28 04:09 message_test.go
-rw-r--r--  1 root root  4709 Jul 28 04:11 pipeline.go
-rw-r--r--  1 root root  6131 Jul 28 04:10 pipeline_test.go
-rw-r--r--  1 root root  6684 Jul 28 02:18 router.go
-rw-r--r--  1 root root 10590 Jul 28 02:17 router_test.go

$ cat .mule/todo.md | head -60
# Mule queue

One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).
The mule runs the topmost open line. Reorder freely; the order is the priority.

If a line needs more detail than fits, write `.mule/tasks/<slug>.md` and keep the line short.

    tools/mule.sh add "<slug>: <outcome>"
    tools/mule.sh loop

A line tagged `[legion]` needs the Legion Go and is skipped automatically while it is asleep;
the queue moves on to the next runnable line rather than stalling behind it.

A line tagged `[readonly]` is a **standing check**, not a piece of work. It verifies something
and is expected to produce no diff at all — so an empty diff is success, not the `wrote
nothing` failure it would be for any other line. It never gets ticked off, because a gate is
never done. Its output lands in `.mule/reports/<slug>.md` with the sha it ran against, and it
is skipped while that sha is still HEAD: re-checking code that has not moved buys nothing and
spends a call on a free provider.

**`mule/queue` is authoritative for this file.** The mule ticks on the Pi and writes its
progress there, so editing the queue on `main` resurrects completed tasks the moment the two
are merged — which has happened once. Add and reorder tasks on `mule/queue`; let them reach
`main` only when that branch is merged.

A line marked `- [~]` is **parked**: real work, but not tick-sized. The mule only ever
picks up `- [ ]`. Benchmark runs live here because a single tier is 5-20 minutes of ingest
alone -- past the per-task budget -- and because they want someone watching. Run one by hand:

    tools/mule.sh legion bench-push
    ssh legion 'cd ~/astrate/bench && ./scripts/run-tier.sh small astrate -base-url ... -housekeeping-key ...'

## Where tasks come from

**This file is not the whole queue, and for real work it is not even the main part of it.**

The queue is: the standing lines below, plus **every open GitHub issue labelled `mule`**.
Issues are read live on each tick and are never copied into this file — a copy would be a
second place the same fact lives, on a branch the mule commits to and you edit on `main`,
and that produced three merge conflicts in one afternoon.

**Open issues labelled `mule` are worked BEFORE the standing lines below** (changed
2026-08-31). An issue is work a human approved; a standing line is either a gate or something
the mule proposed to itself, and self-proposed work never runs ahead of approved work.

**To give the mule work, file an issue and label it `mule`.** From anywhere, by anyone,
including another model with repo access. No SSH, no editing this file:

    gh issue create --label mule --title "<slug>: <outcome>" --body "<the detail>"

Labels on the issue are the tags: `legion` and `readonly` mean what `[legion]` and
`[readonly]` mean here. State lives on the issue, as labels, because there is exactly one
copy of it there:

| label          | meaning                                                          |
|----------------|------------------------------------------------------------------|
| `mule`         | queued                                                            |
| `mule-review`  | the mule pushed something; **it is not merged and not reviewed**  |
| `mule-blocked` | it tried and could not; re-label `mule` to try again              |

Interesting — there's already a substantial flow package. Let me look at what the flow package covers, and examine the reference astarte_flow capabilities to find gaps. Let me look at the current flow implementation and cross-reference with upstream astarte_flow capabilities.

Let me look at the existing flow package and check the reference repo's capabilities.
→ Read internal/flow/block.go
→ Read internal/flow/flow.go
This is a substantial existing implementation. The milestone-2.0 closed issues indicate the flow feature was built (PoC/MVP). Let me check the reference upstream astarte_flow to understand its full capability set, and compare with what Astrate has. Let me look at what the closed issues covered and figure out the remaining gaps.

Let me check the reference repo's capabilities and how Astrate currently exposes flow (HTTP endpoints, etc.).

$ grep -rl "flow" internal/appengine internal/broker 2>/dev/null | head; echo "---"; grep -rn "flow/v1\|api/flow\|/flow" internal/ --include="*.go" -l 2>/dev/null | head
internal/appengine/downsample_test.go
internal/broker/broker.go
---

$ gh issue view 43 --json title,body,state,labels 2>/dev/null | head -80
{"body":"## Decision (2026-07-29)\n\nContainer blocks are **in scope for Flow v2.0** (phased). Product goal: Flow must be able\nto run operator-supplied containers on the message path.\n\n**Not** a v2.0 gate: native Lua blocks, native MQTT/Modbus sources — containers can host\ncustom logic instead. See `docs/handoff/flow-v2-decisions-2026-07-29.md`.\n\n## Phases\n\n### 0 — Design (required before code)\n\n- Operator concept: a pipeline step whose implementation is a container image + config.\n- Message in/out transport: upstream Astarte Flow uses AMQP; Astrate has no AMQP core.\n  Design may choose another bridge (stdio, gRPC, HTTP, etc.) as long as the *idea* matches.\n- Lifecycle: start/stop with the parent flow; fail loudly if image missing / container exits.\n- Security/ops assumptions for PoC (local Docker only is fine).\n\n### 1 — PoC\n\n- One `container` (or equivalent) block type.\n- Manual/local proof: message enters block → container sees it → output returns (or is logged).\n- Document limitations; no requirement for full operator polish.\n\n### 2 — MVP\n\n- Usable inside stored pipelines and named durable flows.\n- Start/stop with flow; status/errors visible via Flow API.\n- Enough for a power user to run a real custom step.\n\n### 3 — Later (split out)\n\nResource limits, pull policies, multi-arch, richer language SDKs, scaling, etc.\n\n## Related\n\n- Named multi-instance + durable flows (#40 + durable issue) — containers run *inside* those.\n- Parity audit (historical “out of v2.0”) **superseded** for containers by the decisions doc.","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACtkBdug","name":"milestone-2.0","description":"astarte-flow feature parity milestone","color":"0E8A16"}],"state":"CLOSED","title":"flow-container-block: PoC → MVP (custom processing in containers)"}

There's a handoff decisions doc referenced. Let me look at it and understand the full state of flow. Let me also check the flow README/docs and how much of the pipeline/container work has been done.
$ ls docs/handoff/ 2>/dev/null; echo "---"; cat docs/handoff/flow-v2-decisions-2026-07-29.md 2>/dev/null | head -100
README.md
phase-1-handoff.md
phase-1-memory.md
phase-10-handoff.md
phase-10-memory.md
phase-11-handoff.md
phase-2-handoff.md
phase-2-memory.md
phase-3-handoff.md
phase-3-memory.md
phase-4-handoff.md
phase-4-memory.md
phase-5-handoff.md
phase-5-memory.md
phase-6-handoff.md
phase-6-memory.md
phase-7-memory.md
phase-8-handoff.md
phase-8-memory.md
phase-9-handoff.md
phase-9-memory.md
---

$ grep -rn "pipeline\|Pipeline" internal/flow/*.go | grep -i "store\|http\|api\|crud" | head; echo "---BLOCKS---"; grep -rln "container\|Container\|Lua\|lua\|mqtt\|native" internal/flow/ 2>/dev/null
---BLOCKS---
internal/flow/pipeline_test.go

• Explore flow feature state Explore Agent
• Research upstream astarte_flow Explore Agent
The `flow-v2-decisions` doc doesn't exist in handoff. Let me check the actual state more carefully. The closed issues were about container blocks, durable flows, etc. Let me verify what's actually implemented by checking the flow package more and how flow is exposed via HTTP/API. Let me check the swagger and any flow API endpoints.

Let me delegate a thorough exploration of the current flow state to a subagent while I check the upstream reference.
! permission requested: external_directory (/dev/*); auto-rejecting
✗ Research upstream astarte_flow failed Explore Agent
Error: Subagent failed (task_id: ses_fa33f5f6cffejKVZeVsRum5YWo): The user rejected permission to use this specific tool call.
✓ Explore flow feature state Explore Agent
The upstream research subagent was rejected (likely because `gh api`/network wasn't permitted). Let me do the upstream reference research myself via web. But first note the important finding: the first subagent confirms that while there's a flow kernel, the operator-facing surface (block catalog, container blocks, pipeline store, REST API, wiring into the binary) is all greenfield.

Given the milestone-2.0 issues are all CLOSED but the milestone status still says "not started", and the huge amount of missing operator-facing surface, I need to understand what the closed issues actually accomplished. Let me check the closed issues' bodies to see if they represent completed work or just proposals that were superseded.

Let me look at several closed issues to understand their outcome.

$ for n in 23 24 25 26 27 37 39 40 42; do echo "===== ISSUE $n ====="; gh issue view $n --json title,state,labels,body 2>/dev/null | python3 -c "import sys,json; d=json.load(sys.stdin); print('TITLE:',d['title']); print('STATE:',d['state']); print('LABELS:',[l['name'] for l in d['labels']]); print('BODY:'); print(d['body'][:1500])"; done
===== ISSUE 23 =====
TITLE: flow-models: FlowMessage wire format and Block/Pipeline/Flow data models
STATE: CLOSED
LABELS: ['milestone-2.0']
BODY:
Astarte Flow's core abstractions — the FlowMessage wire format, Block interface (source/sink/transform), Pipeline graph, and Flow runtime instance — have no equivalent in Astrate.

Implement `internal/flow/` with:
- FlowMessage struct matching astarte_flow's `astarte_flow/message/v0.1` JSON schema (key, type, subtype, data, metadata, timestamp_us)
- Block interface: Source, Sink, Transform function signatures
- Pipeline struct: acyclic graph of named blocks with typed connections, JSON-serializable for storage
- Flow struct: running instance of a pipeline with resource lifecycle

Reference: astarte_flow guides/core_concepts/0002-flow-messages.md, 0003-blocks.md, 0004-pipelines.md, 0005-flows.md (https://github.com/astarte-platform/astarte_flow)

This is the foundation — all subsequent flow issues depend on these types. Port the capability, restated for Go; do not copy the Elixir structure.

Files: `internal/flow/`
Task spec: `.mule/tasks/issue-23.md`

Acceptance: Table-driven unit tests in `internal/flow/message_test.go`, `internal/flow/pipeline_test.go`, and `internal/flow/flow_test.go` asserting FlowMessage JSON round-trip against upstream `astarte_flow/message/v0.1` schema, Pipeline graph serialization, and FlowStatus transitions.
===== ISSUE 24 =====
TITLE: flow-pipeline-store: Pipeline CRUD storage and graph validation
STATE: CLOSED
LABELS: ['mule-review', 'milestone-2.0']
BODY:
Pipelines in astarte_flow are stored as JSON documents and validated on create/update. Astrate has no pipeline storage.

Add pipeline CRUD to `internal/store/`:
- New migration: `pipelines` table (id, realm_id, name, definition JSONB, created_at, updated_at)
- Store methods: CreatePipeline, GetPipeline, UpdatePipeline, DeletePipeline, ListPipelines
- Graph validation on create/update: acyclicity check, all block names resolve, connection types are compatible
- Pipeline definition is a JSON-serializable graph matching the data model from issue #23

Reference: astarte_flow `lib/astarte_flow/pipelines.ex`, `lib/astarte_flow/pipeline_builder.ex` (https://github.com/astarte-platform/astarte_flow)

Files: `internal/store/`, `migrations/`
Task spec: `.mule/tasks/issue-24.md`

Acceptance: A unit test suite in `internal/store/pipelines_test.go` testing graph validation (TestValidatePipelineGraph_Acyclic, TestValidatePipelineGraph_Cyclic returning ErrPipelineCyclic, TestValidatePipelineGraph_InvalidBlockRef) and PostgreSQL CRUD storage methods.
===== ISSUE 25 =====
TITLE: flow-lifecycle: Flow start/stop lifecycle and runtime management
STATE: CLOSED
LABELS: ['milestone-2.0']
BODY:
Flows in astarte_flow are running instances of pipelines — they consume CPU/memory and process messages. Astrate has no flow lifecycle.

Implement flow lifecycle in `internal/flow/`:
- StartFlow(pipelineID, config): instantiate a pipeline into a running Flow, allocate resources
- StopFlow(flowID): graceful shutdown — drain in-flight messages, release resources
- ListFlows / GetFlowStatus: inspect running flows (throughput, error counts)
- Flow status enum: creating, running, stopped, failed
- Integration with the store layer (issue #24) for persistence of flow state

Reference: astarte_flow `lib/astarte_flow/flows.ex`, `lib/astarte_flow/restore_flows_task.ex` (https://github.com/astarte-platform/astarte_flow)

Files: `internal/flow/`, `internal/store/`
Task spec: `.mule/tasks/issue-25.md`

Acceptance: A unit test suite in `internal/flow/manager_test.go` verifying StartFlow transitions state to running, StopFlow drains in-flight messages and releases resources to transition state to stopped, and failed initializations set status to failed.
===== ISSUE 26 =====
TITLE: flow-router: Stream-based message routing through the block graph
STATE: CLOSED
LABELS: ['milestone-2.0']
BODY:
When a Flow is running, messages must be routed through the block graph with key-based stream ordering. Astrate has no message routing for flows.

Implement message routing in `internal/flow/`:
- Router: accepts FlowMessage, routes through the block graph (source → transforms → sink)
- Stream assignment: messages with the same key go to the same processing lane (consistent hashing)
- In-order processing within a stream; interleaving across streams is allowed
- Backpressure: bounded channels per lane, configurable overflow behaviour (drop QoS 0 / block QoS 1+)
- Error handling: block failures are caught per-message, logged, and counted; do not crash the flow

Reference: astarte_flow guides/core_concepts/0006-streams.md, `lib/astarte_flow/flows.ex` (https://github.com/astarte-platform/astarte_flow)

Files: `internal/flow/`
Task spec: `.mule/tasks/issue-26.md`

Acceptance: A unit test suite in `internal/flow/router_test.go` verifying in-order delivery of 100 messages for a single key (TestRouter_InOrderDeliveryPerKey), parallel non-blocking delivery across distinct keys, overflow drop policy under queue saturation, and per-block error isolation.
===== ISSUE 27 =====
TITLE: flow-astarte-source: Astarte device events source block
STATE: CLOSED
LABELS: ['mule-review', 'milestone-2.0']
BODY:
The most critical built-in block for Astrate's flow feature parity: a Source block that ingests Astarte device events and converts them to FlowMessages.

Implement `internal/flow/blocks/astartesource/`:
- AstarteSource block: subscribes to Astrate's existing engine stream bus (`internal/engine/stream`) and emits FlowMessages
- Configurable: realm filter, interface filter, path filter
- Maps incoming device events (datastream, properties, control) to FlowMessage type/subtype/data
- Implements the Block Source interface from issue #23

This connects Astrate's existing ingestion pipeline to the new flow system, enabling operator-defined processing of device data.

Reference: astarte_flow `lib/astarte_flow/blocks/device_events_producer.ex` (https://github.com/astarte-platform/astarte_flow)

Depends on: #23 (flow-models), #26 (flow-router for stream integration)

Files: `internal/flow/blocks/astartesource/`, `internal/engine/stream/`
Task spec: `.mule/tasks/issue-27.md`

Acceptance: Unit tests in `internal/flow/blocks/astartesource/source_test.go` asserting stream event ingestion & conversion into FlowMessage with key `<realm>/<device_id>`, realm/interface filtering, and clean unsubscription on Stop().
===== ISSUE 37 =====
TITLE: flow-source-pump: poll AstarteSource and Stop on flow teardown
STATE: CLOSED
LABELS: ['milestone-2.0']
BODY:
## Context

#27 landed `internal/flow/blocks/astartesource` (merged in #36): a Source block that subscribes to `stream.Bus` and converts events to FlowMessages. Unit tests cover ingest, filters, and clean `Stop()`.

The flow runtime is still **push-based** (`Router.Submit` → `BlockGraph.Run`). Nothing polls `Process(nil)` on sources, and `Manager.StopFlow` does not call `astartesource.Stop()`.

## Work needed

1. **Source pump** — when a flow’s graph starts with (or contains) a Source that produces via `Process(nil)`, run a bounded pump that drains it and feeds the router/graph so bus events actually flow in production.
2. **Teardown** — on `StopFlow` / `Manager.Shutdown`, call `Stop()` on blocks that expose it (at least `astartesource.Source`) so bus subscriptions are released.
3. **Tests** — integration-style test: publish on `stream.Bus` → flow with AstarteSource → sink receives FlowMessage; stop flow → subscriber count returns to 0.

## Out of scope

Re-implementing the block itself (#27 done). This is wiring only.

Depends on: #27 (done), #25/#26 (done).
===== ISSUE 39 =====
TITLE: flow-blocks-discovery: GET /flow/v1/{realm}/blocks for registered catalog types
STATE: CLOSED
LABELS: ['enhancement', 'milestone-2.0']
BODY:
## Context

Parity audit 2026-07-29 (`docs/handoff/flow-parity-audit-2026-07-29.md`): upstream exposes `/{realm}/blocks` so operators can list default + user-defined blocks. Astrate has a code registry (`flow.Registry` + `blocks.DefaultRegistry`) but no HTTP surface.

## In scope

1. `GET /flow/v1/{realm}/blocks` — list registered `block_type` strings (and optionally a short human description / role: source|transform|sink).
2. Optional: `GET /flow/v1/{realm}/blocks/{type}` — return static config notes for built-ins (`filter` keys, `map` keys, `astarte_source` filters).
3. Auth: same as other Flow routes (`a_rma`) until a client requires `a_f`.
4. Tests in `internal/flowapi`.

## Out of scope

- User-defined blocks (DSL fragments).
- JSON Schema validation documents matching upstream `ConfigSchema`.
- Container blocks.
- Creating/deleting blocks via API.

## Acceptance

- Unauthenticated request → 401/403 as other Flow routes.
- With valid JWT, list includes at least: `astarte_source`, `filter`, `map`, `null_sink`, `log_sink`.
===== ISSUE 40 =====
TITLE: flow-multi-instance: named flows + pipeline config (decided — design then implement)
STATE: CLOSED
LABELS: ['enhancement', 'question', 'milestone-2.0']
BODY:
## Context

Parity audit 2026-07-29 (`docs/handoff/flow-parity-audit-2026-07-29.md`).

Upstream Astarte Flow treats a **Flow** as a *named* specialization of a **Pipeline**:
- Create body: `{ name, pipeline, config? }`
- Many flows can run the same pipeline with different configs
- Pipeline `source` may contain `${$.config.*}` parameters validated by an optional JSON Schema

Astrate today:
- `POST /flow/v1/{realm}/flows` with `{ "pipeline": "<name>" }` only
- At most one running instance per `realm/pipeline` (`flow.FlowPipelineID`)
- No pipeline parameter substitution; DAG config is fully concrete at store time

## Decision needed (also escalated in `.mule/for-giulio.md`)

For v2.0, choose:

**(a)** Keep 1:1 realm+pipeline; document as Astrate simplification. Close this issue as wontfix for v2.0.

**(b)** Implement named multi-instance flows:
- `name` on start; manager key = realm+flow name
- optional `config` map
- pipeline definitions may reference config (design TBD — avoid full DSL if possible)

Related: durable `flows` table decision (for-giulio) — multi-instance makes durability more useful.

## Do not implement until (a) or (b) is chosen.
===== ISSUE 42 =====
TITLE: flow-rehydrate-edge-cases: follow-ups after durable auto_restart MVP
STATE: CLOSED
LABELS: ['enhancement', 'milestone-2.0']
BODY:
## Context

Parent product decision: durable flows with default `auto_restart` and loud failure
(`docs/handoff/flow-v2-decisions-2026-07-29.md`). First implementation should stay simple;
this issue **tracks** harder cases so they are not forgotten and do not block MVP.

## Candidates (not committed scope — refine during design)

- [ ] Pipeline definition edited while a flow is running (hot-reload vs require stop/start)
- [ ] Pipeline deleted while flows still reference it
- [ ] Config-only change: restart required?
- [ ] Start fails for one flow among many at boot (others should still start)
- [ ] Crash mid-start / partial graph
- [ ] Operator wants “desired running” vs “last known failed” semantics in the API
- [ ] Multi-process / HA managers (likely out of scope for single-binary Astrate)

## Process

Split into real implementation issues only when prioritising after durable+named MVP ships.

$ ls docs/handoff/*.md; echo "---FLOW AUDIT---"; cat docs/handoff/flow-parity-audit-2026-07-29.md 2>/dev/null | head -80
docs/handoff/README.md
docs/handoff/phase-1-handoff.md
docs/handoff/phase-1-memory.md
docs/handoff/phase-10-handoff.md
docs/handoff/phase-10-memory.md
docs/handoff/phase-11-handoff.md
docs/handoff/phase-2-handoff.md
docs/handoff/phase-2-memory.md
docs/handoff/phase-3-handoff.md
docs/handoff/phase-3-memory.md
docs/handoff/phase-4-handoff.md
docs/handoff/phase-4-memory.md
docs/handoff/phase-5-handoff.md
docs/handoff/phase-5-memory.md
docs/handoff/phase-6-handoff.md
docs/handoff/phase-6-memory.md
docs/handoff/phase-7-memory.md
docs/handoff/phase-8-handoff.md
docs/handoff/phase-8-memory.md
docs/handoff/phase-9-handoff.md
docs/handoff/phase-9-memory.md
---FLOW AUDIT---

Now I have a much clearer picture. These closed issues reference docs and an audit (`docs/handoff/flow-parity-audit-2026-07-29.md`) and a code registry (`flow.Registry` + `blocks.DefaultRegistry`) that the explore agent said **does not exist**. This is confusing. The explore agent reported that there is no registry, no flowapi, no HTTP surface. But the issues (37, 39) say things like "#27 landed `internal/flow/blocks/astartesource`" and reference `flow.Registry` + `blocks.DefaultRegistry`.

There's a contradiction. Let me check: did these closed issues' work actually land, or were they closed without the work being done? Let me look at the referenced handoff docs and check the actual current state more carefully.

The explore agent reported:
- No `internal/flow/blocks/astartesource` exists
- No registry
- No flowapi
- No HTTP surface
- No store for pipelines
- No container blocks

This strongly suggests **the closed issues were closed but the work was NOT actually landed in the current codebase** — OR the work landed on a branch that isn't merged, OR these were proposals that got closed/superseded.

Let me verify by checking the referenced audit doc and git history.
$ git log --all --oneline -- internal/flow/ 2>/dev/null | head -30; echo "===BRANCHES==="; git branch -a 2>/dev/null | head -30
84ae2b6 flow: dynamic virtual_device_pool — first-seen auto-registration (#84)
78e293e feat(flow): ExpandComposites — inline expansion of user blocks with cycle guard (issue #85, phase 85a)
6d60193 feat(flowapi): plumb the engine ingest path into flow Deps (issue #84, phase 84c)
ad55cab feat(flow/blocks): virtual_device_pool block over the engine ingest path (issue #84, phase 84b)
970ce71 feat(flow/blocks): register virtual_device_pool type (stub) and Ingest dep (issue #84, prep 84b)
abd1556 feat(flow/blocks): mqtt_source and mqtt_sink blocks (issue #83)
671982e feat(flow/blocks): serve config schemas from the blocks-discovery API (issue #89)
4d456bd feat(flow/blocks): http_source and http_sink blocks (issue #82)
7f91d26 feat(flow/blocks): json_path_map block (issue #81)
4fa0e3d feat(flow/blocks): pure-transform block set — to_json, update_metadata, split_map, random_source, sort (issue #80)
fb3f335 feat(flow): block-death detection + auto-restart with backoff (#45 phase 1)
4c881d0 fix(lint): resolve golangci-lint v2.12.2 failures across flow container work
801fc48 fix(flow): stop already-constructed blocks when Instantiate fails partway
763be9d fix(flow): synchronize test-local state in TestRegistry_InstantiateLinear
89145e6 Flow v2.0: durable named multi-instance flows + container block PoC
28404cf feat(flow): native filter and map catalog transforms
43c2a0e feat(flow): block factory, built-in catalog, process wiring, operator API
c5fa06b feat(flow): source pump + Stop on teardown (#37); refresh v2.0 milestones
aba2d18 fix: land blocked mule issues #14 #21 #22 #24 #27
acf55be mule: issue-23+25+26: flow-models, flow-lifecycle, flow-router (squashed)
a9481d0 mule: issue-23: flow-models: FlowMessage wire format and Block/Pipeline/Flow data models
ee59732 mule: issue-25: flow-lifecycle: Flow start/stop lifecycle and runtime management
e250381 mule: issue-26: flow-router: Stream-based message routing through the block graph
===BRANCHES===
  main
* mule/queue
  remotes/origin/HEAD -> origin/main
  remotes/origin/feat/pokemon-agent
  remotes/origin/fix/blocked-mule-issues
  remotes/origin/m12-06b-watch-path-onto-main
  remotes/origin/main
  remotes/origin/mule/queue
  remotes/origin/mule/research
  remotes/origin/wip/DietPi-20260727T171543Z
  remotes/origin/worktree-m12-05-match-upstream
  remotes/origin/worktree-m12-06-ach-join
  remotes/origin/worktree-m12-06b-watch-path

$ git branch --contains 89145e6 2>/dev/null; echo "===which branches contain flow v2.0 commit==="; git branch -a --contains 89145e6 2>/dev/null; echo "===contains 43c2a0e (block factory)==="; git branch -a --contains 43c2a0e 2>/dev/null
===which branches contain flow v2.0 commit===
  remotes/origin/HEAD -> origin/main
  remotes/origin/main
===contains 43c2a0e (block factory)===
  remotes/origin/HEAD -> origin/main
  remotes/origin/main

The git history is very revealing. There are many commits with flow work: `feat(flow): ...`, `feat(flow/blocks): ...`, "Flow v2.0: durable named multi-instance flows + container block PoC", "feat(flow): block factory, built-in catalog, process wiring, operator API". 

But the explore agent reported these things **don't exist** in the current working tree. This is a strong contradiction. The explore agent may have been working against `main` (the branch's state), but I'm on `mule/queue`. Wait — the git history shows these commits across `--all`. Let me check whether these commits are on the current branch or another.

Let me check what branch these flow commits are on and whether they've been merged into my current working tree.
$ echo "current HEAD:"; git rev-parse --abbrev-ref HEAD; git log --oneline -1; echo "=== is 89145e6 ancestor of HEAD? ==="; git merge-base --is-ancestor 89145e6 HEAD && echo "YES ancestor" || echo "NO not ancestor"
current HEAD:
mule/queue
3a70ddb mule: refill the queue
=== is 89145e6 ancestor of HEAD? ===
NO not ancestor

Critical finding. The flow v2.0 work (block factory, catalog, container blocks, durable flows, operator API, astarte_source, mqtt/http/json_path blocks, etc.) exists on **`origin/main`** but is **NOT** on my current branch `mule/queue` (HEAD 3a70ddb). My branch is behind main.

The recipe says I'm working toward milestone v2.0. The closed milestone-2.0 issues and the enormous amount of flow work on `main` strongly indicate v2.0 has been **implemented on main**. The milestone status in `.mule/milestones.md` says "not started" but that file is stale — it hasn't been updated.

This is a crucial situation. Let me understand the current state of `main` vs my working branch, and check the actual `.mule/milestones.md` — wait, that file is on my branch. Let me check what main's milestones.md says, and whether the flow v2.0 work on main actually completed the milestone.

Let me look at what main has for flow and whether the milestone is actually complete there.
$ git show origin/main:.mule/milestones.md 2>/dev/null | head -80
# Milestones

Release-tag-gated goals. Each section names a body of work that must be true by the time
that tag is cut. `.mule/recipes/milestones.md` reads this file to find the **current
target** — the first section below not marked `DONE` — and works toward it: investigating,
filing GitHub issues (`milestone-<tag>` label, alongside `mule`) with sub-issues where the
work splits, and escalating anything that needs a design decision to `.mule/for-giulio.md`.

**This file is Giulio's.** The recipe reads it and may propose edits via
`.mule/for-giulio.md`, but never edits it directly — same rule as `docs/COMPATIBILITY.md`.
Mark a milestone `DONE` yourself once the tag is actually cut.

Order matters: milestones are worked **in order**, lowest tag first. Do not start v3.0
investigation while v2.0 has open, un-escalated gaps — say so and stop instead.

---

## v2.0 — astarte-flow feature parity

Reference: astarte_flow (https://github.com/astarte-platform/astarte_flow) — upstream's
Elixir "Flow" component. It lets ingested data get piped through a graph of processing
blocks (native blocks and containerised ones) before it lands in storage, with pipelines
described as reusable, parametrised graphs.

Scope for this milestone: whatever set of Flow's *capabilities* Astrate needs to expose the
same wire-visible behaviour and operator-facing concepts (pipelines, blocks, native vs.
containerised blocks) — not a port of the Elixir implementation. See
`.mule/recipes/astarte-upstream.md`'s rule: port the idea, restated in Go, never the code.

Status: **DONE** (2026-07-29). Runtime, factory, catalog (incl. filter/map), process
wiring, and `/flow/v1` API are on `main`. **Parity audit** + **product decisions**
recorded 2026-07-29 (`docs/handoff/flow-parity-audit-2026-07-29.md`,
`docs/handoff/flow-v2-decisions-2026-07-29.md`).

**Design A + B landed on `main`** (commit `89145e6`, 2026-07-29): durable flows +
auto_restart (**#41 closed**); named multi-instance + config (**#40 closed**) —
migration `000009`, store, `${config.*}`, API, boot rehydrate. Container block
PoC→MVP (**#43 closed**) — registered in the catalog, usable inside stored
pipelines/named flows; see `docs/handoff/flow-design-b-container-block-2026-07-29.md`.
Blocks discovery (**#39 closed**).
**Not a v2.0 gate:** native Lua / MQTT blocks.

**#42 closed 2026-07-29** (rehydrate edge cases). Triaged all seven candidates against a
live e2e Docker smoke test of #43: one real bug found and fixed on `main`
(`801fc48` — `Registry.Instantiate` leaked an already-started Docker container when a
later block's constructor failed mid-pipeline-build; mutation-tested regression added);
two candidates ("one flow fails at boot, others still start"; "pipeline deleted while
flows reference it") were already correct by design, no code change needed; three
("hot-reload a running pipeline" #44, "partial restart of failed blocks" #45, "update a
running flow's config" #46) split into their own demand-driven backlog issues per the
explicit v2.0 decisions doc; multi-process/HA managers stay out of scope indefinitely
(single-process design throughout).

### Landed (on main as of 2026-07-29)

| Piece | Where | Notes |
|---|---|---|
| FlowMessage wire format | `internal/flow/message.go` | `astarte_flow/message/v0.1` |
| Block / Source / Stopper / graph | `internal/flow/block.go`, `graph.go` | linear chain; Sources skipped in `Run` |
| Pipeline DAG model + validate | `internal/flow/pipeline.go` | serialisable description |
| Flow lifecycle manager | `internal/flow/flow.go` | Start / Stop / List / Shutdown |
| Stream router (lanes, QoS, metrics) | `internal/flow/router.go` | in-order per key |
| Source pump + Stop on teardown | `internal/flow/flow.go` (#37) | pumps `Source.Emit` → router; `Stopper.Stop` after drain |
| AstarteSource block | `internal/flow/blocks/astartesource` | bus → FlowMessage (#27) |
| Pipeline store + migration | `internal/store/pipelines.go`, `000008_pipelines` | realm-scoped DAG JSON (#24) |
| Durable flows store + migration | `internal/store/flows.go`, `000009_flows` | named multi-instance + auto_restart (#40+#41) |
| Config substitution | `internal/flow/substitute.go` | `${config.key}` in block config strings |
| Block factory + instantiate | `internal/flow/factory.go` | `Registry`, `ParseDefinition`, topo order → `[]Block` |
| Built-in catalog | `internal/flow/blocks/catalog.go` + `transform.go` | `astarte_source`, `filter`, `map`, `null_sink`, `log_sink` |
| Process wiring | `cmd/astrate/main.go` | bus + Manager; **boot rehydrate** before listen; shutdown marks stopped |
| Operator HTTP API | `internal/flowapi` | pipelines + named durable `/flows` (a_rma) |

### Closed gaps (all of them — v2.0 is DONE)

1. **~~Design A: durable + named multi-instance~~** — **implemented, #40+#41 closed.**
2. **~~Design B: container block~~** — **implemented (PoC→MVP), #43 closed.** Doc:
   `docs/handoff/flow-design-b-container-block-2026-07-29.md`. PoC transport: HTTP
   (not AMQP).
3. **~~Blocks discovery API~~** — **implemented, #39 closed.**
4. **~~Parity audit + product decisions~~** — **done** 2026-07-29.

$ git log origin/main --oneline -40 | grep -i "flow\|milestone\|v2.0" | head -40
84ae2b6 flow: dynamic virtual_device_pool — first-seen auto-registration (#84)
14708db feat(flowapi): composites accepted and expanded in the flow lifecycle (issue #85, phase 85d)
4b6dfc7 feat(flowapi): user-block CRUD endpoints with schema validation (issue #85, phase 85c)
1cee640 mule: record FDO milestone-4.0 decision in for-giulio
78e293e feat(flow): ExpandComposites — inline expansion of user blocks with cycle guard (issue #85, phase 85a)
6d60193 feat(flowapi): plumb the engine ingest path into flow Deps (issue #84, phase 84c)
ad55cab feat(flow/blocks): virtual_device_pool block over the engine ingest path (issue #84, phase 84b)
970ce71 feat(flow/blocks): register virtual_device_pool type (stub) and Ingest dep (issue #84, prep 84b)
abd1556 feat(flow/blocks): mqtt_source and mqtt_sink blocks (issue #83)
671982e feat(flow/blocks): serve config schemas from the blocks-discovery API (issue #89)
55bb717 feat(auth,flowapi): a_f JWT claim for the Flow API (issue #88)
4d456bd feat(flow/blocks): http_source and http_sink blocks (issue #82)
7f91d26 feat(flow/blocks): json_path_map block (issue #81)
4fa0e3d feat(flow/blocks): pure-transform block set — to_json, update_metadata, split_map, random_source, sort (issue #80)

**Critical finding.** On `origin/main`, the v2.0 milestone is **marked DONE** (status `DONE (2026-07-29)`). The flow v2.0 work is fully implemented on main: runtime, factory, catalog, process wiring, `/flow/v1` API, durable named multi-instance flows, container block, blocks discovery, and all closed gaps.

But my current branch `mule/queue` has a **stale** `.mule/milestones.md` that still says v2.0 is "not started" — it's behind main.

So the recipe on my branch sees v2.0 as "not started" and would re-investigate and re-file issues. But that would be **duplicating work that's already done on main**, which the recipe explicitly forbids ("Never file the same gap twice").

Let me verify whether main's milestones.md is actually the more current one and whether v2.0 truly is done. Let me check the current state — is `mule/queue` merely behind main, and does main have the authoritative milestones state?

Let me look at the state more carefully — what does main's milestones.md say for v3.0 and what's the actual "first section not marked DONE"?
$ git show origin/main:.mule/milestones.md 2>/dev/null | sed -n '80,200p'
4. **~~Parity audit + product decisions~~** — **done** 2026-07-29.
   Docs: `flow-parity-audit-2026-07-29.md`, `flow-v2-decisions-2026-07-29.md`.
5. **~~Rehydrate edge cases~~** — **#42 closed 2026-07-29.** One real bug fixed
   (`Registry.Instantiate` container-leak on partial pipeline build, `801fc48`); rest
   already correct or split into #44/#45/#46 (demand-driven, no milestone).
6. **Still out of v2.0 gate (demand-driven):** native Lua/JSONPath blocks; native
   MQTT/Modbus/HTTP poll I/O; full pipeline DSL; `a_f` path wire-compat without a
   client; `http_sink` unless a client needs it.

### Explicitly out of this milestone (tracked elsewhere)

- Dormant trigger types / previous-value cost → #20 (Legion bench) + for-giulio decision
- Group-scoped triggers → #17 (and related #18–19 wire/authz)
- External-bus intake design → #10
- Edgehog client compatibility → #28 (`readonly` / mule) — may integrate upstream FOSS as-is

### Project scope (standing)

Astrate reimplements **free and open-source** software worth recreating for its purpose: a
wire-compatible Astarte-platform reimplementation in Go (lighter components) plus extras
such as AtomVM compatibility. Integrating existing FOSS when it already fits is preferred
over a full rewrite (e.g. original Astarte Dashboard; Edgehog under investigation in #28).

---

## v3.0 — upstream Astarte feature parity (1.2.2 → 1.3/1.4)

Scope decided by Giulio on 2026-08-22. The original v3.0 candidate (building the
remaining SECO suite pieces) was dropped for legal reasons, and Edgehog is already
fully compatible with Astrate (verified in #28) — so this milestone takes the body of
work produced by the 2026-08-22 upstream-parity investigation instead:

- Source document: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 /
  v1.4.0-rc.5 diffed service-by-service against our emulated 1.2.2 level).
- Backlog: GitHub issues **#47–#89** (`upstream-parity` label), plus the parked
  **#90** APICompatVersion decision (milestone-4.0 label — pull into this milestone's
  final phase, after the 1.3-surface audit it requires).

Shape of the milestone: work through #47–#89 in the research doc's priority order
(retention ceiling → alias/group surfaces → query formats → wire capabilities →
validation/error-code fidelity), escalating every issue marked "decide" to
`.mule/for-giulio.md` instead of choosing unilaterally. Deliberately parked items
(FDO #78, lua_map #87, pipeline DSL #86, keyAgreement #51) stay out unless Giulio
pulls them in.

Status: not started. First recipe job: triage #47–#89 into an ordered plan (which are
audits vs features vs decisions), file sub-issues where work splits, escalate the
"decide" set in one batch.

---

## v4.0+ — next

Placeholder. Scope undecided: candidates once v3.0 is DONE include the remaining
`upstream-parity` backlog items parked today, or whatever emerges from upstream
releases after v1.4 ships. The recipe should propose the next section's draft
(name + one-line scope) as a `for-giulio.md` entry once v3.0 is `DONE`, rather than
leaving this placeholder to rot.

Status: not started, scope undecided.

$ echo "=== diverged from main? ==="; git rev-list --left-right --count mule/queue...origin/main 2>/dev/null; echo "(left=unique to queue, right=unique to main)"; echo "=== recent main commits ==="; git log origin/main --oneline -15
=== diverged from main? ===
107	120
(left=unique to queue, right=unique to main)
=== recent main commits ===
60a8ec6 mule: approved issues outrank the mule's own proposals
271c055 survey: actually commit what it wrote
4e0ac73 mule: pin the model explicitly — the agent's default no longer resolves
306adc0 mule: survive cron — run from the repo, find opencode under $HOME
c522db1 docs(agents): Legion Go is off by default — probe once, then ask Giulio
72b025e chore(lint): clear golangci-lint debt to zero
7a244b2 feat(payload): reject object documents missing a required key (#67 fase 4b)
fd562f7 payload: pre-wire missing_required reject reason label fixture and error-name translation (fase 4b prep)
9c4d411 feat(interfaceschema): required+encrypted mapping fields parse and gate (#67 fase 4a)
84ae2b6 flow: dynamic virtual_device_pool — first-seen auto-registration (#84)
0376b83 docs: resolve dormant-trigger decision line — measured in #20, activated by 6bd14a7
1e4d4d5 docs: probe-craft lessons from the 2026-08-24 verify batch; clarify the Mac test rule
ab31319 docs(compat): measured write taxonomy + unknown-realm deviation; fix Internal server error casing to measured upstream form
ad4eb45 fix(appengine): server-write error taxonomy matches measured upstream (#57)
cda4010 feat(http): per-service version endpoints, upstream-parity (#77)

This is decisive. On `origin/main`, the milestones file has been **rewritten**:
- **v2.0 is `DONE`** (2026-07-29)
- **v3.0 is now scoped** — "upstream Astarte feature parity (1.2.2 → 1.3/1.4)" with a decided scope (backlog issues #47–#89, upstream-parity label, researched in `.mule/research/upstream-parity-2026-08-22.md`)
- The first section **not marked DONE** on main is **v3.0**, whose first recipe job is: "triage #47–#89 into an ordered plan... file sub-issues where work splits, escalate the 'decide' set in one batch."

My branch `mule/queue` is stale — it has the old milestones.md showing v2.0 "not started" and v3.0 as TBD (CLEA). The **authoritative state is on main**.

But the recipe told me to work on the milestones that my branch's `.mule/milestones.md` presents. There's a conflict: my branch says "work v2.0, not started," but main shows v2.0 done and v3.0 is the actual current target.

The key question: what should I actually do? The recipe's job is to file issues for milestone v2.0 gaps. But there are **no v2.0 gaps** — v2.0 is DONE on main. If I blindly follow my stale branch's milestones.md, I'd investigate v2.0 and find it's already fully implemented (or worse, re-file issues that duplicate closed work).

The recipe's rule: "Never file the same gap twice" and "If, after step 3, `gh issue list --label milestone-2.0 --state open` is empty and you found no new gaps: propose one `.mule/for-giulio.md` line... and stop."

Given that v2.0 is DONE on main (the issues are closed, the code is on main), the honest and correct outcome for this run is:
- The authoritative milestones state is on `main`, and it already reflects v2.0 = DONE, v3.0 scoped as the current target.
- My branch is stale. The recipe on my branch would wrongly target v2.0.

Let me verify the current git situation to understand whether this is a state conflict I should report, and check what `gh issue list` shows for the actual current milestone (v3.0 / upstream-parity). Let me check whether the mule/queue branch is normally expected to have merged main's milestones.

Let me check the divergence between my branch and main, and verify `.mule/milestones.md` and `.mule/research/` on main.
$ echo "=== upstream-parity issues on main (current milestone) ==="; gh issue list --search 'label:upstream-parity' --state all --limit 60 --json number,title,state,labels --jq '.[] | "\(.number)\t\(.state)\t\(.title)"' 2>/dev/null | head -60
=== upstream-parity issues on main (current milestone) ===
91	OPEN	Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health
89	CLOSED	Dashboard flow-block schema mismatch (split_map/virtual pools hardcoded; null_sink/log_sink unknown)
88	CLOSED	Flow auth: support a_f JWT claim
87	OPEN	Flow block: lua_map — needs embedded Lua runtime (parked)
86	CLOSED	Flow: pipeline source DSL — keep DAG-JSON as documented deviation?
85	CLOSED	Flow API: user-defined composite blocks
84	CLOSED	Flow blocks: virtual_device_pool / dynamic_virtual_device_pool
83	CLOSED	Flow blocks: mqtt_source/mqtt_sink (+modbus_tcp_source?) demand-driven
82	CLOSED	Flow blocks: http_source/http_sink (demand-driven)
81	CLOSED	Flow block: json_path_map
80	CLOSED	Flow blocks: pure-transform set (to_json, update_metadata, split_map, random_source, sort)
79	CLOSED	Verify registration-limit-reached HTTP status vs upstream
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
77	CLOSED	Verify per-service version endpoints (GET /version, GET /v1/{realm}/version) served everywhere
76	CLOSED	Housekeeping: GET /v1/realm-defaults/replication — decide reject/deviate (Cassandra-shaped)
75	CLOSED	Housekeeping: decide realm-deletion gating/preconditions vs always-sync deviation
74	CLOSED	Housekeeping: PATCH /v1/realms/{realm} (jwt key, registration limit, retention; null=unset)
73	CLOSED	Housekeeping: default datastream retention injection env var (upstream 1.4)
72	CLOSED	Realms: datastream_maximum_storage_retention ceiling (create/patch/enforce)
71	CLOSED	Pairing: realm-scoped health check GET /v1/{realm}/health (upstream 1.3)
70	CLOSED	Triggers: audit wildcard semantics (interface_name '*', match_path '/*' forcing rules)
69	CLOSED	Verify unknown-realm HTTP status on RM endpoints against upstream
68	OPEN	Decide async_operation=false params vs documented always-sync deviation
67	OPEN	Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4)
66	CLOSED	Realm Management: detailed=true interface listing with full mappings (upstream 1.4)
65	CLOSED	Policies: handler-overlap rejection + retry_times coupling + prefetch_count
64	CLOSED	Triggers: decide AMQP action behavior (validate-reject vs NATS-forward deviation)
63	CLOSED	Triggers: HTTP action validation limits (URL/method/header blocklist/template size)
62	CLOSED	Realm Management: audit install/update/delete error codes and statuses
61	CLOSED	Realm Management: audit interface/mapping validation matrix against astarte_core
60	CLOSED	Realm Management: GET config/datastream_maximum_storage_retention (since upstream 1.2.0)
59	CLOSED	AppEngine: group create-body validation + UUID-v1 from_token for group device listing
58	CLOSED	AppEngine: PATCH requires Content-Type application/merge-patch+json
57	CLOSED	AppEngine: audit server-write error taxonomy against upstream
56	CLOSED	AppEngine: object-aggregate response formats table and disjoint_tables
55	CLOSED	AppEngine: complete data-query parameter set
54	CLOSED	AppEngine: GET /groups/{g} show-group endpoint
53	CLOSED	AppEngine: group-scoped device PATCH + interface data access under groups/{g}/devices/{d}
52	CLOSED	AppEngine: devices-by-alias full mirror (PATCH via alias + interface data access)
51	CLOSED	keyAgreement control topic + encrypted mappings (upstream 1.4 experimental) — investigate/park
50	CLOSED	Broker/engine hardening: discard malformed introspection/purge/capabilities instead of crashing
49	CLOSED	MQTT v1: binaryblob properties re-send as BSON subtype 0 in consumer/properties
48	CLOSED	MQTT v1: tolerate devices connecting with an empty introspection
47	CLOSED	MQTT v1: accept device capabilities incl. purge_properties_compression_format (plaintext purge)

$ echo "=== current branch for-giulio.md ==="; cat .mule/for-giulio.md 2>/dev/null | head -80; echo "=== EXIT $? ==="
=== current branch for-giulio.md ===
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- ~~`device_deletion_started`/`device_deletion_finished` trigger events are not emitted~~ —
  **decided 2026-07-27: emit both, back-to-back, around the synchronous delete.** Filed as
  issue #21 (`mule`). (Cross-project survey, 2026-07-27,
  `.mule/research/survey-2026-07-27.md` source 4.)
- ~~Mustache trigger-action templates are accepted but not rendered~~ — **decided
  2026-07-27: implement it.** Guiding principle clarified: Astarte compatibility means
  SDK/wire compatibility, not minimum dependency count — Astrate is allowed to be a
  compatible *superset*. Library picked: `github.com/cbroglie/mustache`. Filed as issue #22
  (`mule`). (Same survey, source 4.)
- **`value_change`/`value_change_applied`/`path_created`/`path_removed`/`value_stored` trigger
  types compile but never fire** (`internal/engine/triggers/match.go:30-42`). Decision
  deferred pending data: issue #20 (`mule`, `readonly`) asks Big Pickle to benchmark the
  ingest-path cost of a previous-value lookup on the Legion Go before this gets decided one
  way or the other. (Same survey, source 4.)
- **Group-scoped triggers (`group_name` on device/data triggers) compile but never match**
  (`internal/engine/triggers/match.go:11-12`). Decision deferred, tied to issue #17
  (group-WATCH-path reconciliation, trickle work, not mule): whatever group-membership
  mechanism comes out of that phase should also report the perf cost for this decision —
  noted in a comment on #17 so it isn't benchmarked twice. (Same survey, source 4.)

---

- ~~The Pi cannot run the race detector~~ — **resolved 2026-07-27** by installing Go 1.26.5
  as a userland toolchain on the Legion Go (`~/.local/go`, no root, `rm -rf` to undo). The
  Pi still cannot run `-race` (39-bit VMA kernel vs the 48 ThreadSanitizer needs), so its
  gate remains `go vet ./... && go test ./...` — but race coverage now exists on the Legion
  Go, where the full suite runs clean in ~40s on 16 cores. The standing `race-check` task is
  the concurrency gate. Concurrency work is queueable again, provided the race-check runs
  after it.
- **golangci-lint is not installed on the Pi**, so the mule's second gate is silently absent
  there — `gofmt` still runs, the linter does not. `go install
  github.com/golangci/golangci-lint/cmd/golangci-lint@<the pinned version>` on the Pi would
  close it; I did not pick a version for you, since the pin is a decision.
- ~~`/root/astrate` on the Pi has uncommitted work~~ — **resolved 2026-07-27** with the new
  `tools/reconcile.sh`: rescued onto `origin/wip/DietPi-20260727T171543Z` (pushed, not
  reviewed — read the diff before merging anything from it) and `/root/astrate` is now clean
  on `main`. Also had to set `commit.gpgsign false` locally on that clone first — it had
  signing on with no working gpg-agent for a non-interactive session, same class of problem
  `/root/astrate-mule` already solved the same way (your call, confirmed 2026-07-27: align it
  with the existing unattended-Pi-clone convention rather than fix gpg-agent or bypass signing
  silently inside the script).
- **The mule has been idle 15h.** Filed by the dead-man's switch; see journalctl on the Pi.
=== EXIT 0 ===

Now the picture is fully clear. The situation is:

**The milestone work has already been done on `main`.** My current branch `mule/queue` is **120 commits behind `origin/main`** and has a **stale `.mule/milestones.md`** that still shows:
- v2.0 "not started" (but on main it's DONE since 2026-07-29)
- v3.0 as CLEA "TBD" (but on main it's been rescoped to "upstream Astarte feature parity 1.3/1.4")

The **authoritative current target is v3.0** (upstream-parity), and it has a well-defined first recipe job listed right in main's milestones.md: "triage #47–#89 into an ordered plan... file sub-issues where work splits, escalate the 'decide' set in one batch."

My branch is simply stale/behind. The recipe's core rules forbid:
- "Never file the same gap twice"
- Duplicating closed work

If I blindly execute the recipe against my branch's stale milestones.md, I'd target **v2.0** — which is fully DONE on main with all issues CLOSED. Re-investigating and filing v2.0 issues would duplicate completed work. That is exactly what the recipe warns against.

The correct, honest outcome: v2.0 is done (nothing to file there), and my branch's milestones file is stale relative to the authoritative main. I should not file spurious v2.0 issues. The right move is to report this — that v2.0 is complete on main and the real target is v3.0 (already scoped on main), and that my branch needs to sync with main before the milestones recipe can correctly target v3.0.

Let me double-check whether the current branch state is what the mule normally runs against, and confirm the milestone-2.0 open-issue list is empty (it is — all CLOSED). Let me also verify the v3.0 recipe first-job isn't something I should start.

Actually, re-reading the recipe's step 5: "If, after step 3, `gh issue list --label milestone-<tag> --state open` is empty and you found no new gaps: propose one `.mule/for-giulio.md` line — 'milestone <tag> looks complete, verify and cut the tag' — and stop."

For v2.0 on my branch: `gh issue list --label milestone-2.0 --state open` is **empty** (all closed), and I found **no new v2.0 gaps** (the work is done on main). So per the recipe, the correct outcome for the v2.0 milestone is to **propose that v2.0 looks complete** — not to invent issues.

But I need to be careful. The recipe targets the "first section not marked DONE" in the milestones.md **on the branch I'm running on** (mule/queue), which is stale. Let me reconcile this properly.

Let me check whether the milestone-2.0 work is actually present/checked on the current branch or only on main — did the mule already process this? Let me verify whether v2.0 issues being closed corresponds to the work being merged to main, and whether `mule/queue` is expected to have this already.

The key realization: **the mule branch is 120 commits behind main**, which means either (a) the queue branch is mid-rebase/merge-pending, or (b) something's off. But regardless, my job is clear: the v2.0 milestone as presented to me is complete on main, has zero open issues and zero new gaps, so per recipe step 5 I should propose a `.mule/for-giulio.md` line saying v2.0 looks complete/verify-and-cut, and stop — **without filing any issues** (which would duplicate closed work).

However, filing a "v2.0 looks complete, verify and cut the tag" line would be misleading given that this is *already known* on main (main's milestones already marks v2.0 DONE). The real issue is my branch is stale.

Given the recipe's intent and rules, and to be genuinely useful rather than duplicative, the most accurate outcome is:

1. Report that v2.0 (the target my stale milestones.md points at) is **already complete on main** — no gaps, no issues to file (step 5 condition met, and even stronger: it's marked DONE on the authoritative branch).
2. Flag the staleness: `mule/queue` is 120 commits behind `origin/main`; its `.mule/milestones.md` is outdated. The real current milestone is v3.0 (upstream-parity), which is already scoped on main.
3. **File no GitHub issues** (v2.0 gaps would all be duplicates of closed work).

Let me write the `.mule/for-giulio.md` entry reflecting the true state so the next run (after a sync) targets v3.0 correctly. Let me check the current for-giulio.md first.
$ echo "=== milestone-3.0 or v3.0 related tags ==="; gh issue list --label milestone-3.0 --state all --limit 20 --json number,title,state --jq '.[]|"\(.number)\t\(.state)\t\(.title)"' 2>/dev/null; echo "---"; echo "=== is master milestone file on mule/research showing same? ==="; git show origin/mule/research:.mule/milestones.md 2>/dev/null | sed -n '18,30p'
=== milestone-3.0 or v3.0 related tags ===
---
=== is master milestone file on mule/research showing same? ===

Now I've confirmed:
- No `milestone-3.0` label issues exist yet.
- The current branch's for-giulio.md is stale (it's the queue branch's version).

Let me confirm the state of the research doc on main referenced by v3.0 and check the mule/research branch's milestones for comparison. Let me also verify my understanding of the current git branch situation — am I really working on the stale `mule/queue`? Yes, HEAD is `3a70ddb mule: refill the queue` on branch `mule/queue`.

Let me finalize my analysis. The situation is:

**My branch (`mule/queue`) has a stale `.mule/milestones.md`** that presents v2.0 as "not started" — but v2.0 is **already DONE on `origin/main`** (120 commits ahead), with all milestone-2.0 issues CLOSED and the work (factory, catalog, container block, durable named flows, /flow/v1 API) fully merged.

The recipe's intent (work the current target, never duplicate) and its rules (never file the same gap twice; step 5: propose "verify and cut the tag" when nothing's open) both point the same way for the milestone my stale file targets: **v2.0 has no open issues and no new gaps — do not file anything that would duplicate closed work.**

The genuinely useful, honest deliverable for this run is a `.mule/for-giulio.md` entry that:
1. States v2.0 (as my branch's milestones.md targets it) is already complete on main — nothing to file.
2. Flags the real problem: `mule/queue` is 120 commits behind `origin/main`, so its milestones.md is stale; the authoritative current target is **v3.0** (upstream-parity #47–#89), whose first recipe job (triage the backlog) is already spelled out on main.

I'll **not** file any GitHub issues (they'd all be duplicates). Let me write the for-giulio.md entry.

Let me look at the current for-giulio.md's top structure to prepend correctly.
→ Read .mule/for-giulio.md [limit=15]
The current branch's for-giulio.md has content but is stale. I'll prepend my entry at the top (newest at the top). Let me add the entry about the milestones state.

Let me insert the entry right after the top `---` separator (line 11).

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,25 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestones recipe has nothing to file: v2.0 is already DONE, and `mule/queue` is 120
+  commits behind `origin/main`, so this branch's `.mule/milestones.md` is stale.** My branch
+  still shows v2.0 "not started" (and v3.0 as CLEA TBD), so the recipe pointed me at v2.0 —
+  but on `main` v2.0 is marked **DONE (2026-07-29)**: every `milestone-2.0` issue is CLOSED
+  (#23–#27, #37, #39, #40, #41, #42, #43), and the factory/catalog/container-block/durable
+  named flows/`/flow/v1` API all live on `main` (`git show origin/main:.mule/milestones.md`).
+  `gh issue list --label milestone-2.0 --state open` is empty and there are no new v2.0 gaps
+  to file, so per the recipe's step 5 I filed **nothing** rather than duplicate closed work.
+  **Real current target is v3.0** ("upstream Astarte feature parity 1.3/1.4", scoped on main
+  2026-08-22): the first recipe job there is to triage `upstream-parity` backlog **#47–#89**
+  (source `.mule/research/upstream-parity-2026-08-22.md`) into a plan. But I should not start
+  that triage from this stale branch. Action needed: **merge/rebuild `mule/queue` on top of
+  `origin/main`**, then the next milestones run targets v3.0 with the correct research doc.
+  (Milestones recipe run, 2026-09-01.)
+
+---
+
 - ~~`device_deletion_started`/`device_deletion_finished` trigger events are not emitted~~ —
   **decided 2026-07-27: emit both, back-to-back, around the synchronous delete.** Filed as
   issue #21 (`mule`). (Cross-project survey, 2026-07-27,
   `.mule/research/survey-2026-07-27.md` source 4.)

Now let me verify the file edit and confirm I should NOT touch the milestones.md (which I haven't done — the recipe forbids it). Let me confirm no git operations were done and the working tree only has the for-giulio.md addition and no gh issue creation.

Let me verify the final state of my edits.
$ sed -n '11,22p' .mule/for-giulio.md
---

- **Milestones recipe has nothing to file: v2.0 is already DONE, and `mule/queue` is 120
  commits behind `origin/main`, so this branch's `.mule/milestones.md` is stale.** My branch
  still shows v2.0 "not started" (and v3.0 as CLEA TBD), so the recipe pointed me at v2.0 —
  but on `main` v2.0 is marked **DONE (2026-07-29)**: every `milestone-2.0` issue is CLOSED
  (#23–#27, #37, #39, #40, #41, #42, #43), and the factory/catalog/container-block/durable
  named flows/`/flow/v1` API all live on `main` (`git show origin/main:.mule/milestones.md`).
  `gh issue list --label milestone-2.0 --state open` is empty and there are no new v2.0 gaps
  to file, so per the recipe's step 5 I filed **nothing** rather than duplicate closed work.
  **Real current target is v3.0** ("upstream Astarte feature parity 1.3/1.4", scoped on main
  2026-08-22): the first recipe job there is to triage `upstream-parity` backlog **#47–#89**

The for-giulio.md entry is correctly in place. Let me confirm I made no other changes (only for-giulio.md, no git mutations, no gh issue creations).
$ git status --short
M .mule/cron.log
 M .mule/for-giulio.md

The only change I made is `.mule/for-giulio.md` (the `.mule/cron.log` change is pre-existing/environmental, not mine). I did not touch git, did not touch `.mule/milestones.md`, and filed no GitHub issues.

Let me compose the final report.

**Summary of what I did:**
- Read MULE.md and the milestones recipe.
- Read `.mule/milestones.md` → the first section not marked DONE on my branch is **v2.0 — astarte-flow feature parity** (status "not started").
- Investigated: `gh issue list --label milestone-2.0 --state all` → all **CLOSED** (14 issues); `--state open` → **empty**.
- Discovered a critical divergence: my branch `mule/queue` is **120 commits behind `origin/main`**, and on `main` the v2.0 milestone is marked **DONE (2026-07-29)** with the full implementation shipped (factory, block catalog incl. filter/map, container block PoC→MVP, durable named multi-instance flows, `/flow/v1` API), and v3.0 has been **rescoped** to "upstream Astarte feature parity 1.3/1.4" (#47–#89, upstream-parity label).
- Per the recipe's rules (never file the same gap twice; step 5: stop when the issue list is empty and there are no new gaps), I filed **no issues** — anything I filed for v2.0 would duplicate already-closed, already-shipped work.
- Instead, appended a `.mule/for-giulio.md` entry explaining v2.0 is done, that the branch's milestones file is stale (120 commits behind main), and that the real target going forward is v3.0 after `mule/queue` is synced with `origin/main`.
$ echo "audit: no issues created this run (milestone-2.0 all closed, nothing filed)"; gh issue list --label milestone-2.0 --state open --limit 5 2>/dev/null | wc -l
audit: no issues created this run (milestone-2.0 all closed, nothing filed)
0

Done: v2.0 milestone is already complete on `origin/main` — all milestone-2.0 issues CLOSED, no open gaps, so I filed no issues; the real issue is this branch's `.mule/milestones.md` is stale (120 commits behind main, still shows v2.0 "not started" and v3.0 as CLEA TBD, while main marks v2.0 DONE and rescopes v3.0 to upstream-parity #47–#89).

Files: .mule/for-giulio.md
Verified: gh issue list --label milestone-2.0 --state all/open; git log/merge-base origin/main vs mule/queue; git show origin/main:.mule/milestones.md -> pass (no v2.0 issues to file; branch divergence confirmed 120 commits)
Unsure: nothing — I made no code change, no test applies (proposal job, step 5 stop condition met)
Follow-ups: none appended to .mule/todo.md (nothing to propose; I wrote one .mule/for-giulio.md entry instead, per recipe step 5 / "never file the same gap twice")

I deliberately did **not** create any GitHub issues: every v2.0 gap already exists as a closed, shipped issue on `main`. The deliverable is the `for-giulio.md` entry flagging that `mule/queue` needs to be synced with `origin/main` before the next milestones run can correctly target the real current milestone (v3.0, upstream-parity). I also did not touch git or `.mule/milestones.md` (both forbidden).
