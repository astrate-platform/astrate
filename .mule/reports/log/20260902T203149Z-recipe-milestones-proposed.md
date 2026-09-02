slug: recipe-milestones
verdict: proposed
at:  262f6f6
ran: 2026-09-02T20:31:49Z on DietPi in 336s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
Now I need to read `.mule/milestones.md` to find the first incomplete milestone.
→ Read .mule/milestones.md
→ Read .mule/todo.md
→ Read .mule/for-giulio.md
The milestone file I read still shows v2.0 as "not started", but `for-giulio.md` says this branch is stale. Let me verify the actual state on `origin/main` and check existing issues.

$ git -C /root/astrate-mule show origin/main:.mule/milestones.md 2>/dev/null || echo "FAILED"
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

$ gh issue list --label milestone-2.0 --state all --limit 50 2>/dev/null || echo "GH_FAILED"
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

$ gh issue list --label milestone-3.0 --state all --limit 50 2>/dev/null || echo "GH_FAILED"
(no output)

$ ls -la /root/astrate-mule/.mule/research/ 2>/dev/null || echo "NO_RESEARCH_DIR"
NO_RESEARCH_DIR

I have a stale branch situation. The `git show origin/main:.mule/milestones.md` reveals the real state:

- **v2.0 is DONE** (2026-07-29) — every `milestone-2.0` issue is CLOSED, no open gaps.
- **v3.0** is the real current target: "upstream Astarte feature parity (1.2.2 → 1.3/1.4)". Status: not started. First recipe job: **triage #47–#89 into an ordered plan**, file sub-issues where work splits, escalate the "decide" set.

`gh issue list --label milestone-3.0` returned no output (no milestone-3.0 issues exist).

The for-giulio note from a prior run already flagged exactly this: the branch is stale and targets v3.0 on main. But I need to be careful — the recipe says to read `.mule/milestones.md` (the local file). However, the local file is stale. My job is to work the current milestone correctly.

Per the recipe step 3, I need to investigate the gap. The reference for v3.0 is `.mule/research/upstream-parity-2026-08-22.md` and the backlog issues #47–#89. The first recipe job is to **triage #47–#89 into an ordered plan**: which are audits vs features vs decisions, file sub-issues where work splits, escalate the "decide" set.

Let me check the local research doc (if it exists) and the issue list. Note: my branch is stale, so the local research file may be missing too. Let me check.
$ gh issue list --label upstream-parity --state all --limit 100 2>/dev/null || echo "GH_FAILED"
91	OPEN	Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health	mule-blocked, upstream-parity	2026-08-31T11:31:15Z
89	CLOSED	Dashboard flow-block schema mismatch (split_map/virtual pools hardcoded; null_sink/log_sink unknown)	enhancement, upstream-parity	2026-08-23T15:08:58Z
88	CLOSED	Flow auth: support a_f JWT claim	enhancement, upstream-parity	2026-08-23T15:08:59Z
87	OPEN	Flow block: lua_map — needs embedded Lua runtime (parked)	enhancement, upstream-parity	2026-08-22T02:24:58Z
86	CLOSED	Flow: pipeline source DSL — keep DAG-JSON as documented deviation?	enhancement, upstream-parity	2026-08-23T11:57:02Z
85	CLOSED	Flow API: user-defined composite blocks	enhancement, upstream-parity	2026-08-23T21:47:32Z
84	CLOSED	Flow blocks: virtual_device_pool / dynamic_virtual_device_pool	enhancement, upstream-parity	2026-08-25T14:56:53Z
83	CLOSED	Flow blocks: mqtt_source/mqtt_sink (+modbus_tcp_source?) demand-driven	enhancement, upstream-parity	2026-08-23T15:09:03Z
82	CLOSED	Flow blocks: http_source/http_sink (demand-driven)	enhancement, upstream-parity	2026-08-23T15:09:04Z
81	CLOSED	Flow block: json_path_map	enhancement, upstream-parity	2026-08-23T15:09:06Z
80	CLOSED	Flow blocks: pure-transform set (to_json, update_metadata, split_map, random_source, sort)	enhancement, upstream-parity	2026-08-23T15:09:07Z
79	CLOSED	Verify registration-limit-reached HTTP status vs upstream	enhancement, upstream-parity	2026-08-24T13:21:57Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
77	CLOSED	Verify per-service version endpoints (GET /version, GET /v1/{realm}/version) served everywhere	enhancement, upstream-parity	2026-08-24T13:23:44Z
76	CLOSED	Housekeeping: GET /v1/realm-defaults/replication — decide reject/deviate (Cassandra-shaped)	enhancement, upstream-parity	2026-08-22T19:12:29Z
75	CLOSED	Housekeeping: decide realm-deletion gating/preconditions vs always-sync deviation	enhancement, upstream-parity	2026-08-23T15:09:10Z
74	CLOSED	Housekeeping: PATCH /v1/realms/{realm} (jwt key, registration limit, retention; null=unset)	enhancement, upstream-parity	2026-08-23T15:09:12Z
73	CLOSED	Housekeeping: default datastream retention injection env var (upstream 1.4)	enhancement, upstream-parity	2026-08-23T15:09:13Z
72	CLOSED	Realms: datastream_maximum_storage_retention ceiling (create/patch/enforce)	enhancement, upstream-parity	2026-08-23T15:09:15Z
71	CLOSED	Pairing: realm-scoped health check GET /v1/{realm}/health (upstream 1.3)	enhancement, upstream-parity	2026-08-23T15:09:16Z
70	CLOSED	Triggers: audit wildcard semantics (interface_name '*', match_path '/*' forcing rules)	enhancement, upstream-parity	2026-08-23T15:09:18Z
69	CLOSED	Verify unknown-realm HTTP status on RM endpoints against upstream	enhancement, upstream-parity	2026-08-24T13:23:42Z
68	OPEN	Decide async_operation=false params vs documented always-sync deviation	enhancement, mule-blocked, upstream-parity, upstream-experimental	2026-08-31T11:47:36Z
67	OPEN	Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4)	enhancement, upstream-parity, upstream-experimental	2026-08-30T23:33:48Z
66	CLOSED	Realm Management: detailed=true interface listing with full mappings (upstream 1.4)	enhancement, upstream-parity	2026-08-23T15:09:20Z
65	CLOSED	Policies: handler-overlap rejection + retry_times coupling + prefetch_count	enhancement, upstream-parity	2026-08-23T15:09:22Z
64	CLOSED	Triggers: decide AMQP action behavior (validate-reject vs NATS-forward deviation)	enhancement, upstream-parity	2026-08-23T15:09:23Z
63	CLOSED	Triggers: HTTP action validation limits (URL/method/header blocklist/template size)	enhancement, upstream-parity	2026-08-23T15:09:25Z
62	CLOSED	Realm Management: audit install/update/delete error codes and statuses	enhancement, upstream-parity	2026-08-23T15:09:26Z
61	CLOSED	Realm Management: audit interface/mapping validation matrix against astarte_core	enhancement, upstream-parity	2026-08-23T15:09:28Z
60	CLOSED	Realm Management: GET config/datastream_maximum_storage_retention (since upstream 1.2.0)	enhancement, upstream-parity	2026-08-22T03:23:47Z
59	CLOSED	AppEngine: group create-body validation + UUID-v1 from_token for group device listing	enhancement, upstream-parity	2026-08-22T10:35:18Z
58	CLOSED	AppEngine: PATCH requires Content-Type application/merge-patch+json	enhancement, upstream-parity	2026-08-22T03:15:59Z
57	CLOSED	AppEngine: audit server-write error taxonomy against upstream	enhancement, upstream-parity	2026-08-24T13:24:06Z
56	CLOSED	AppEngine: object-aggregate response formats table and disjoint_tables	enhancement, upstream-parity	2026-08-23T15:09:31Z
55	CLOSED	AppEngine: complete data-query parameter set	enhancement, upstream-parity	2026-08-23T15:09:33Z
54	CLOSED	AppEngine: GET /groups/{g} show-group endpoint	enhancement, upstream-parity	2026-08-22T10:35:16Z
53	CLOSED	AppEngine: group-scoped device PATCH + interface data access under groups/{g}/devices/{d}	enhancement, upstream-parity	2026-08-22T09:18:23Z
52	CLOSED	AppEngine: devices-by-alias full mirror (PATCH via alias + interface data access)	enhancement, upstream-parity	2026-08-22T09:18:20Z
51	CLOSED	keyAgreement control topic + encrypted mappings (upstream 1.4 experimental) — investigate/park	enhancement, upstream-parity	2026-08-22T08:48:28Z
50	CLOSED	Broker/engine hardening: discard malformed introspection/purge/capabilities instead of crashing	enhancement, upstream-parity	2026-08-22T04:04:24Z
49	CLOSED	MQTT v1: binaryblob properties re-send as BSON subtype 0 in consumer/properties	enhancement, upstream-parity	2026-08-22T04:04:23Z
48	CLOSED	MQTT v1: tolerate devices connecting with an empty introspection	enhancement, upstream-parity	2026-08-22T04:04:11Z
47	CLOSED	MQTT v1: accept device capabilities incl. purge_properties_compression_format (plaintext purge)	enhancement, upstream-parity	2026-08-22T04:04:09Z

Let me look at the milestones file more carefully. The real v3.0 section on `origin/main` says:

> Status: not started. First recipe job: triage #47–#89 into an ordered plan (which are audits vs features vs decisions), file sub-issues where work splits, escalate the "decide" set in one batch.

But looking at the actual issue list, **most of #47–#89 are already CLOSED**. Only a handful are OPEN:

- **#91** Pairing health endpoint divergence — `mule-blocked, upstream-parity` (OPEN)
- **#87** lua_map — parked (OPEN)
- **#78** FDO device onboarding — milestone-4.0 (OPEN)
- **#68** async_operation=false decision — `mule-blocked, upstream-parity, upstream-experimental` (OPEN)
- **#67** Interfaces: decide required/encrypted mapping fields — upstream-experimental (OPEN)

Wait — most are labeled `mule-blocked`. That's the blocked label on the mule/queue side, meaning the mule tried them and gates failed. But the milestone status says "not started". The backlog #47-#89 has largely been worked through (closed via mule tasks).

So the situation: the v3.0 milestone's "first recipe job — triage #47–#89" has essentially been completed by previous runs (most closed). What remains open:

Let me think about what the recipe wants. The recipe v3.0's stated job is to triage #47–#89. Most are closed now. The open ones are:
- #68, #91, #67 are `upstream-experimental` and/or `mule-blocked` — decisions.
- #87, #78 are deliberately parked.

But I should be careful about the staleness. My local `.mule/milestones.md` (which I read, per the recipe step 1) shows v2.0 as "not started" — which is stale and wrong. The author's instruction was to read `.mule/MULE.md` then `.mule/recipes/milestones.md` and carry it out. 

The recipe step 1 says: "Read `.mule/milestones.md`. Find the **first** section not marked `DONE`." The local file's first non-DONE section is v2.0. But that's stale — on main v2.0 is DONE. 

The prior run (2026-09-01) already flagged exactly this and filed nothing, correctly pointing out the branch is stale and needs merge/rebuild against origin/main. My for-giulio entry notes the current real target is v3.0.

I need to decide how to proceed. The task explicitly says this is a *proposal* job producing new task lines appended to `.mule/todo.md`, plus evidence files. I should not start implementing. The recipe targets the current milestone.

Given the branch is stale, the honest and correct action per the prior run's precedent and MULE.md's "say so and stop instead of guessing" is to work against the *actual* current milestone (v3.0), but be careful.

Let me reconsider. The instruction to me says to read the local files and carry out the recipe. But the local milestones.md is stale. The recipe's step 1 uses the local file. However, I've discovered the real state via `origin/main` and `gh`.

The cleanest approach consistent with the established pattern (and the prior run's note): the current milestone is v3.0. For v3.0, the recipe job (per origin/main's milestones.md) is to **triage #47–#89 into an ordered plan** — but most are already closed. So effectively the triage is largely done.

Let me check what's still open among the v3.0 pieces and whether there's anything new to file. The genuinely open upstream-parity issues:

1. **#91** Pairing health divergence — this is an `upstream-parity` issue that's OPEN and `mule-blocked` (mule tried, gates failed). This is a real gap: upstream serves `GET /pairing/health` (no realm), Astrate serves `GET /pairing/v1/{realm}/health`. This is machine-checkable. But it's already filed as #91.
2. **#67** Interfaces decide required/encrypted mapping fields — a decision (`upstream-experimental`).
3. **#68** async_operation decision — a decision.

So the open v3.0 issues are all either decisions (should go to for-giulio, and indeed #67/#68 relate to decisions) or already-filed machine-checkable gaps (#91).

Let me not over-engineer. The recipe says: if `gh issue list --label milestone-<tag> --state open` is empty and no new gaps → propose "milestone <tag> looks complete, verify and cut the tag". But there are no `milestone-3.0` issues at all (the tag used is `upstream-parity`, not `milestone-3.0`). 

Let me reconsider what milestone tag to use. The v3.0 milestone's issues are labeled `upstream-parity`. The recipe would use `milestone-<tag>` where tag for v3.0 would be `3.0`. But the existing convention uses `upstream-parity` label. Hmm.

Actually, let me re-read the important context. The prior run (2026-09-01) already did this exact analysis and concluded: branch is stale, action needed is "merge/rebuild mule/queue on top of origin/main, then the next milestones run targets v3.0 with the correct research doc." That note is still in for-giulio.md.

Given I'm running on a stale branch, and the previous run's recommendation was to not start v3.0 triage from the stale branch, I should follow that precedent. The correct action: I should not fabricate work. 

But wait — let me reconsider whether the branch is actually still stale. Let me check the git state.
$ cd /root/astrate-mule && git fetch -q origin 2>&1; git status -sb 2>&1 | head; echo "---"; git log --oneline -3 2>&1; echo "---LOCAL vs MAIN---"; git rev-list --count HEAD..origin/main 2>&1
## mule/queue
 M .mule/cron.log
---
262f6f6 mule: blocked lint-flow-stutter
2cbc7d3 mule: blocked lint-pairing-test-param
eab8e24 mule: blocked lint-flow-findings
---LOCAL vs MAIN---
120

The branch is indeed 120 commits behind `origin/main`, exactly as the prior run's for-giulio note said. So my local `.mule/milestones.md` is stale.

Now, the important question: how do I carry out the recipe given this staleness?

The user's instruction says to read `.mule/MULE.md`, then `.mule/recipes/milestones.md`, and carry it out. The recipe's step 1 says read `.mule/milestones.md` (local). But the local file is stale. Following the prior run's established precedent and the recipe's spirit ("investigate the gap between the reference and Astrate's current code"), the real current milestone is v3.0.

The prior run's for-giulio note says: "Action needed: merge/rebuild mule/queue on top of origin/main, then the next milestones run targets v3.0 with the correct research doc." I cannot touch git (no commit/branch/checkout), so I cannot merge/rebuild the branch myself. But I can still read origin/main's milestones.md and the current GitHub issue state to do the v3.0 triage honestly, using the live/remote truth.

For v3.0, the recipe job is to triage #47–#89 into an ordered plan. Looking at the actual issue list, **most of #47–#89 are already CLOSED** (they've been worked through by the mule on the real queue). The open `upstream-parity` ones are:

- **#91** Pairing health endpoint divergence — OPEN, `mule-blocked`. Machine-checkable gap. Already filed.
- **#87** lua_map — OPEN, deliberately parked.
- **#78** FDO — OPEN, milestone-4.0, deliberately parked.
- **#68** async_operation decision — OPEN, `mule-blocked`, decision.
- **#67** Interfaces required/encrypted mapping fields decision — OPEN, decision.

So for v3.0, there's essentially nothing new to file. Everything machine-checkable that the triage would produce has been filed (and mostly closed, i.e., implemented). What remains open is either decisions or deliberately-parked items.

Per the recipe step 4/5: the "decide" set should be escalated to for-giulio. #67 and #68 are decisions. Let me check whether they've already been escalated to for-giulio — and the last run's note already flagged the staleness.

Let me reconsider the whole situation. This is a tricky edge case. The right thing per MULE.md is: **say so and stop instead of guessing**. The branch is 120 commits stale; the prior run already flagged this and the resolution requires a git action (merge/rebuild) I'm not allowed to do.

But I shouldn't just repeat the same note verbatim. Let me consider what NEW value I can add this run. The prior note (2026-09-01) recommended merging the branch. If it's still stale (120 behind), then the blocker is still there. But is there anything actionable I can do as a proposal job?

Let me re-read the recipe steps for v3.0 since v2.0 is DONE on main:
- The v3.0 milestones section says "Status: not started. First recipe job: triage #47–#89 into an ordered plan".

Given most #47-#89 are closed, the triage is effectively done. The remaining open items are decisions/parked.

Actually, let me reconsider whether I should be working against origin/main's milestone file at all, or whether the recipe strictly wants me to use the local `.mule/milestones.md`. The user's instruction references the local path. But the local path is stale and would point me at v2.0 (a DONE milestone on main) — investigating a DONE milestone would duplicate closed work, which the recipe explicitly forbids ("Never file the same gap twice", "trust gh issue list").

The correct, honest behavior: use the live issue state and origin/main's milestones.md as the source of truth (the prior run did exactly this). The current milestone is v3.0. 

For v3.0, checking `gh issue list --label upstream-parity --state open` — open ones are #91, #87, #78, #68, #67. Let me verify none are milestone-3.0-filed (they're upstream-parity labeled). The recipe's step 5 says: "If, after step 3, `gh issue list --label milestone-<tag> --state open` is empty and you found no new gaps: propose one for-giulio line — 'milestone <tag> looks complete, verify and cut the tag' — and stop."

Hmm, but there are no issues labeled `milestone-3.0` at all, and there ARE open upstream-parity issues. The v3.0 milestone work is tracked under `upstream-parity` label, not `milestone-3.0`. So the milestone isn't "complete" in the cut-the-tag sense — several parity items are still open (though mostly decisions or blocked-by-gates).

Given all this complexity and staleness, and following MULE.md's "say so and stop instead of guessing", the right action this run is:

1. Confirm the branch is stale (done: 120 behind).
2. Report that I'm not filing anything new because (a) the branch is stale, (b) the real v3.0 work is already filed, and (c) the remaining open items are decisions/parked that need Giulio, and escalate the current state.

But let me not duplicate the existing note. Let me identify what's CHANGED or actionable new since the last run. 

Actually, let me reconsider. Let me look at whether the remaining open upstream-parity items represent *decisions* that should go to for-giulio. #67 (interfaces required/encrypted) and #68 (async_operation) are marked as decisions. Let me check `gh issue view` on the key open ones to understand their status and whether they're properly escalated.

Let me also double check — is there actually a milestone-3.0 label convention? Let me look.

Actually, I think I'm overcomplicating. Let me step back and determine the cleanest, correct deliverable.

Given:
- The branch is stale (120 behind main). This is the dominant fact — the prior run already recommended fixing it via merge/rebuild.
- The local milestones.md is stale (shows v2.0 not started).
- On main, the current milestone is v3.0, whose triage job has essentially been completed (most #47-#89 closed).
- The remaining open items are either decisions (#67, #68, #91-adjacent) or deliberately parked (#87, #78).

The most honest and useful outcome: 
1. Don't file duplicate issues (nothing new to file that isn't already filed).
2. Escalate the genuine remaining decisions to for-giulio (if not already there).
3. Flag that this run's local milestone state is stale and the recipe can't be cleanly executed until the branch is rebuilt — but with a sharper, updated recommendation.

Let me check the open upstream-parity issues' bodies to determine the "decide" items and whether they're in for-giulio. Let me view a few.
$ gh issue view 67 2>&1 | head -60; echo "======91======"; gh issue view 91 2>&1 | head -60
title:	Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4)
state:	OPEN
author:	atsetilam (Giulio Malitesta)
labels:	enhancement, upstream-experimental, upstream-parity
comments:	1
assignees:	
projects:	
milestone:	
issue-type:	
parent:	
sub-issues:	
sub-issues-completed:	
blocked-by:	
blocking:	
number:	67
--
Neither field exists at 1.2.2/1.3.0; both appear on master. `required` applies to object-aggregated mappings (missing key at runtime → `missing_required_mapping`); `encrypted` ties into the keyAgreement work. Decide now: reject (strict 1.2.2), ignore-and-accept, or implement `required` early.
Size S decision (required impl M).\n\n---\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).
======91======
title:	Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health
state:	OPEN
author:	atsetilam (Giulio Malitesta)
labels:	mule-blocked, upstream-parity
comments:	2
assignees:	
projects:	
milestone:	
issue-type:	
parent:	
sub-issues:	
sub-issues-completed:	
blocked-by:	
blocking:	
number:	91
--
Found during the 2026-08-24 verify batch (#57/#69/#77/#79 probes): upstream Astarte 1.2.0 answers

- `GET /pairing/health` → 200
- `GET /pairing/v1/{realm}/health` → 404 Page not found (bench and scratch realms both)

Astrate implements the unauthenticated health route at `GET /pairing/v1/{realm}/health` (landed in realmcfg-02, issue #71, commit 5111e3c) — the path shape does not exist upstream. The dashboard polls per-service health at the service root, so the current route may be invisible to it.

Decide: move/duplicate the handler at `GET /pairing/health` (keeping or dropping the v1/{realm} variant), and re-check whatever probe led realmcfg-02 to the v1 shape — that session's stack had a port override in play, so the original measurement may have hit a different router.

These are both decision-type items already filed as issues. #67 and #91 both ask for a decision. They're upstream-parity-labeled items (part of the v3.0 body of work) that need Giulio's call — they should be escalated to `for-giulio.md` per the recipe ("Needs a design decision... write it to `.mule/for-giulio.md`, one line").

Let me check if these decision items are already in for-giulio.md. Looking at the for-giulio content I read — the top entries are govulncheck, appengine-api doc divergence, Router.Submit TOCTOU, then #87/#78/#1, then the milestones/nothing-to-file note from 2026-09-01. Neither #67's required/encrypted decision nor #91's pairing health decision is currently in for-giulio.md. #68 async_operation is also a decision.

Let me verify #68 and #67 aren't escalated, and check the current `upstream-parity` open set comprehensively.

The prior run's for-giulio note (2026-09-01) said the branch should be rebuilt and then v3.0 triage happens. Since the branch is still stale, and most #47-#89 are closed, the remaining open upstream-parity items are the "decide" set that the milestone's shape specifically says to escalate: "escalating every issue marked 'decide' to `.mule/for-giulio.md` instead of choosing unilaterally."

So the correct, non-duplicative, useful deliverables this run:
1. Escalate the remaining open "decide" items (#67, #68, #91) to `.mule/for-giulio.md` — these are genuine design decisions the recipe says to escalate, and they're not yet there. This is new, non-duplicative value.
2. Note the staleness blocker (already noted, but I should keep it from going stale / confirm it still applies).

But wait — I must be careful. The recipe says decisions get written to for-giulio "one line each, exactly as github-issues.md already does". But are these already being tracked? #68 is `mule-blocked` and there IS a task line in todo.md for issue-68 (`issue-68-async-operation-accepted`). That task predates the decision — actually no, that's a task to implement accepting async_operation=false. Hmm, but #68 is labeled as a "Decide" issue. Let me not conflate.

Let me reconsider scope. I should be careful not to overload for-giulio with things already known or already handled elsewhere. Let me focus.

The cleanest analysis of the current state:
- v3.0's backlog (#47-#89) has been largely implemented/closed. The milestone per origin/main is "not started" but that's a status-tracking lag; the actual triage+implementation has progressed via the upstream-parity issues.
- Remaining open upstream-parity issues: #91 (health endpoint decision), #87 (lua_map, parked), #78 (FDO, milestone-4.0), #68 (async_operation decision), #67 (required/encrypted decision).

The genuinely-actionable proposal-jobs remaining for v3.0 that are NOT already done:
- The "decide" items should be escalated to for-giulio so Giulio can decide: #67, #68, #91.

These are not currently in for-giulio.md. Filing them is the recipe's specified behavior for design-decision gaps. And I'm not filing GitHub issues / not duplicating (they're already gh issues; escalation to for-giulio is complementary and explicitly required by the recipe for the decide set).

Now, before I write, let me double-check there isn't an existing milestone-3.0 for-giulio escalation I'm missing, and confirm these three decisions aren't already in for-giulio. From what I read, they're not. Let me also check #68's body to escalate accurately.
$ gh issue view 68 2>&1 | head -40
title:	Decide async_operation=false params vs documented always-sync deviation
state:	OPEN
author:	atsetilam (Giulio Malitesta)
labels:	enhancement, mule-blocked, upstream-experimental, upstream-parity
comments:	2
assignees:	
projects:	
milestone:	
issue-type:	
parent:	
sub-issues:	
sub-issues-completed:	
blocked-by:	
blocking:	
number:	68
--
Upstream defaults to async execution with `?async_operation=false` opt-in for: realm create/delete (housekeeping), interface install/update/delete (RM), policy delete (master). Astrate is always-sync everywhere — a superset clients tolerate, but a behavioral difference worth either matching or recording as deviation #11.
Size S decision.\n\n---\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).

Wait — #68 says Astrate is "always-sync everywhere". But todo.md has a `issue-68-async-operation-accepted` task line: "accept and ignore `?async_operation=false`... tests cover the flag parsing (closes #68)". That indicates decision was already made — #68 is effectively decided (implement the accept-and-ignore) and the task is blocked on gates, not on a decision. So #68 is NOT pending a decision anymore; it has a task line. I should not escalate #68 as an open decision.

Let me be careful. Let me reconsider each remaining open item:

- **#91** (pairing health endpoint location): genuinely undecided — "Decide: move/duplicate the handler". Open decision. Not in for-giulio. → escalate.
- **#67** (required/encrypted mapping fields): genuinely undecided — "Decide now: reject / ignore / implement required". Open decision. Not in for-giulio. → escalate.
- **#68** (async_operation): has an implementation task line already (`issue-68-async-operation-accepted`), so decision appears made (accept-and-ignore always-sync). Not pending for-giulio. → do not escalate (would duplicate).
- **#87** (lua_map): already parked in for-giulio (line exists). → skip.
- **#78** (FDO): already parked in for-giulio. → skip.

So the new, non-duplicative escalations this run are **#67** and **#91** — two open design decisions in the v3.0 body of work, currently not in for-giulio.md.

But wait — I must double check these aren't newly-relevant only because the branch is stale. #67 and #91 are on GitHub (visible regardless of branch), labeled upstream-parity (the v3.0 working set). They're genuinely open design decisions that the recipe for v3.0 says to escalate to for-giulio. This is legitimate, in-scope work for the milestones recipe.

Now let me also reconsider the staleness blocker. The prior run's note (2026-09-01) already flagged the branch is stale. It's still stale (120 behind, unchanged). I should update/reiterate but not bloat. Actually, I shouldn't re-file the same stale note. But I should make sure my report reflects that I executed the recipe against the live/remote truth despite the local file being stale.

Let me now also make sure I'm not missing that the recipe wants me to possibly propose a for-giulio "milestone v3.0 looks complete" line. But no — there are still open upstream-parity items (decisions), so I can't propose "complete, cut tag". The recipe's step 5 condition (empty open issue list for the milestone) doesn't hold, AND we have escalation candidates.

Let me write the two escalations to for-giulio.md. I'll add them at the top (newest first), matching the format.

Let me check the exact top format of for-giulio (lines 11-13):
```
---

- **`govulncheck` is not installed on the Pi**, ...
```
The newest entries are at top after the first `---`. I'll insert after line 13 (or before the first existing `- **` item). Actually looking again, line 13 is the last line before line 15's `---`. So there are two `---` then the next item at line 17. The structure: header, `---`, newest item, `---`, next item... So each item is separated by `---`. I'll insert my newest items above the existing newest item (the govulncheck one), keeping them grouped with `---` separators.

Let me write them. Format per recipe: "one line each" but for-giulio entries can be richer. The recipe says for-giulio one line each exactly for these. I'll keep them concise but with the evidence (issue number) inline.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,16 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **v3.0 decide #67** (`upstream-parity`): `required` and `encrypted` mapping fields exist on upstream master but not at 1.2.2/1.3.0 (`required` = missing key at runtime → `missing_required_mapping` on object-aggregated mappings; `encrypted` ties into keyAgreement #51). Decide now: reject (strict 1.2.2), ignore-and-accept, or implement `required` early. Issue is OPEN, no decision on record, not yet escalated here.
+
+---
+
+- **v3.0 decide #91** (`upstream-parity`, `mule-blocked`): upstream serves the unauthenticated health check at `GET /pairing/health` and returns 404 for `GET /pairing/v1/{realm}/health`, but Astrate only implemented the `v1/{realm}` shape (issue #71, commit 5111e3c) — the dashboard polls per-service health at the service root, so the current route may be invisible to it. Decide: move/duplicate the handler at `GET /pairing/health` (keep or drop the v1/{realm} variant), and re-check the original realmcfg-02 measurement (a port override may have skewed it).
+
+---
+
 - **`govulncheck` is not installed on the Pi**, so the hygiene recipe's highest-priority check (reachable dependency vulnerabilities) could not run here. `go install golang.org/x/vuln/cmd/govulncheck@latest` would close it — I did not pick a version because the pin is a decision. The `go list -m -u` sweep shows only version-skew on transitive deps (cloud/azure SDKs, cel/expr, clickhouse); none corresponds to an identified fix this repo needs, so I proposed no bumps. (Hygiene recipe run, 2026-09-02.)
 
 ---
 

Now let me verify the resulting for-giulio.md is well-formed.
→ Read .mule/for-giulio.md [limit=30]
The escalations are in place. Now, per the recipe, I should also consider whether to append task lines to `.mule/todo.md`. But I filed no GitHub issues this run (nothing new to file — all machine-checkable gaps already exist as issues). The two escalations are decisions, which the recipe says go to for-giulio, not to issues or todo task lines.

Regarding the staleness blocker: the prior run's note (2026-09-01) recommended merging/rebuilding the branch against main. The branch is still 120 commits behind. I won't add a new duplicate for-giulio note for that since it's already there, but my report will surface it.

Let me do a final check — confirm I haven't missed any open milestone-3.0-relevant issue that represents a *new machine-checkable* gap worth filing. The open upstream-parity set is: #91 (decision, now escalated), #87 (parked), #78 (milestone-4.0), #68 (has impl task), #67 (decision, now escalated). Nothing new to file. Good.

I did not touch git, did not run `gh issue create` (recipe didn't explicitly require any specific create command this run — I only escalated decisions to for-giulio, which is the correct recipe behavior for the decide set). No new todo.md task lines.

Let me write my report in MULE.md format.

$ cd /root/astrate-mule && git diff --stat .mule/for-giulio.md
.mule/for-giulio.md | 8 ++++++++
 1 file changed, 8 insertions(+)

Confirmed only `.mule/for-giulio.md` changed (a permitted escalation channel, not a frozen file). I did not touch git, did not run `gh issue create` (nothing new to file), and added no task lines.

```
Done: Milestones recipe run — local branch is 120 commits stale (v2.0 still shows "not started"), so worked against origin/main/.mule/milestones.md + live issue state: current milestone is v3.0, whose #47–#89 backlog is largely closed; filed nothing new and instead escalated the two remaining open v3.0 design decisions (#67 required/encrypted mapping fields, #91 pairing health endpoint location) to .mule/for-giulio.md. Skipped #68 (already has an implementation task, decision not pending), #87/#78 (already parked in for-giulio).
Files: .mule/for-giulio.md
Verified: git show origin/main:.mule/milestones.md + gh issue list --label upstream-parity --state all -> no new machine-checkable gap to file (all such gaps already exist as issues); read current for-giulio.md (pass)
Unsure: whether the milestone recipe should have used the stale local .mule/milestones.md verbatim (v2.0, "not started" there but DONE on main); I followed the recipe's own anti-duplication rule and the prior run's precedent and treated v3.0 as the target. The branch is 120 commits behind origin/main and still needs the merge/rebuild the 2026-09-01 note already flags — I can't do that (no git), so the next run should re-check after it lands.
Follow-ups: none (no issue filed, no task lines appended — the two escalations are decisions, which go to for-giulio, not the queue)
```

One note: I filed no GitHub issues this run because every machine-checkable v3.0 gap is already filed; per your instruction, filing an issue was not required by the recipe here. I escalated only the genuinely-undecided items to `for-giulio.md` (#67, #91), which is the recipe's prescribed action for the "decide" set.
