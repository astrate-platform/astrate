slug: recipe-milestones
verdict: proposed
at:  bf76d88
ran: 2026-09-04T17:27:55Z on DietPi in 102s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
Now let me read the milestones file and check the current state.
→ Read .mule/milestones.md

$ gh issue list --label milestone- --state all --limit 50 2>&1 | head -50
(no output)

The first section not marked DONE is **v2.0 — astarte-flow feature parity**. Status: not started, no milestone-2.0 issues filed. Let me check the milestone-2.0 issues specifically and understand the reference.

The `gh issue list --label milestone-` returned nothing. Let me check for open milestone issues and read the reference recipe.
→ Read .mule/recipes/astarte-upstream.md

$ gh issue list --label milestone-2.0 --state all --limit 50 2>&1; echo "---open mule---"; gh issue list --label mule --state open --limit 50 2>&1
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
---open mule---

All milestone-2.0 issues are CLOSED. Let me check the recipe step 1: "Find the first section not marked DONE". v2.0 is not marked DONE (it says "Status: not started. No milestone-2.0 issues filed yet" — but actually 11 issues were filed and closed). 

The instructions say: if an earlier milestone still has open, un-escalated gaps, stop. Let me verify v2.0 issue state. All listed milestone-2.0 issues are CLOSED. So no open gaps there.

Wait — the milestones.md says "Status: not started. No `milestone-2.0` issues filed yet." But the actual GitHub state shows 11 issues filed and all closed. The milestone file is stale. But I never touch milestones.md.

Per recipe step 1: first section not marked DONE is v2.0. All its gaps are closed (no open issues). The sweet spot: step 3 says don't duplicate closed issues, check what exists.

Since the milestone-2.0 issues are all closed, that milestone's scope appears complete from the issue list perspective — but the milestone file says "not started". This is a discrepancy. Let me think about what step suggests.

Actually, step 5 applies when after step 3, `gh issue list --label milestone-<tag> --state open` is empty AND no new gaps found → propose "milestone 2.0 looks complete, verify and cut the tag" in for-giulio.md.

Let me follow the recipe for v2.0 properly. Let me investigate the reference and existing coverage. Let me look at Astrate's internal/ and docs for flow coverage, and check the upstream astarte_flow reference.

Let me check what's currently implemented in the repo for flow.
$ ls internal/ 2>&1; echo "=== flow dirs ==="; ls internal/flow 2>/dev/null || ls -d internal/*flow* 2>/dev/null || echo "no flow dir"
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
=== flow dirs ===
block.go
flow.go
flow_test.go
graph.go
manager_test.go
message.go
message_test.go
pipeline.go
pipeline_test.go
router.go
router_test.go

$ gh issue list --label milestone-2.0 --state open --limit 50 2>&1; echo "=== for-giulio current ==="; cat .mule/for-giulio.md 2>/dev/null || echo "no file"
=== for-giulio current ===
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **Dependency sweep corrected: direct (pinned) deps DO have newer versions** — the 2026-09-02 note said the `go list -m -u` sweep showed "only version-skew on transitive deps", but that run hit the recipe's `head -20` cutoff (all cloud/azure/transitive) and never reached the directly-required modules. Full sweep, 2026-09-04. None of these is a fix this repo *needs*, so no bump is proposed — recorded for the decision. Per module (current → available; breaking change; repo use):
  - `github.com/coder/websocket` v1.8.14 → v1.8.15 — no breaking (patch); used in `internal/appengine/stream/ws.go`, `channels/ws.go`; worth it only for the "transmit in single frame when compression enabled" fix + read-path alloc reduction.
  - `go.etcd.io/bbolt` v1.4.3 → v1.5.0 — bbolt's semver promises no API change between patch/minor, so additive-only; used in `internal/broker/sessionstore.go`; v1.5 adds a data-file size limit and panic-recovery hardening, nothing Astrate needs.
  - `go.mongodb.org/mongo-driver/v2` v2.6.0 → v2.8.2 — the 2.8.0 breaking changes are confined to Queryable Encryption string-query options (`options.Text()`→`String()`); Astrate uses only the raw BSON API (`pkg/payload/bson.go`, `internal/engine/capabilities.go`, `bench/`) and is unaffected.
  - `github.com/nats-io/nats.go` v1.52.0 → v1.53.1 — no breaking; the headline fixes (JetStream `resetOrderedConsumer` race, KV dot-rejection) are paths Astrate does not use — `internal/engine/forward/nats.go` is core NATS publish only.
  - `github.com/prometheus/client_golang` v1.23.2 → v1.24.1 — requires Go ≥1.25 (fine, repo is 1.26.1); the breaking `LabelNames`/remote-api renames don't touch repo usage (`prometheus`/`collectors`/`promhttp` in `internal/observability/metrics.go`, flow/engine metrics); would buy `Gather()` panic-recovery and opt-in `CoalesceGather` scrape-pile-up protection.
  - `github.com/testcontainers/testcontainers-go` v0.43.0 → v0.44.0 (modules/postgres v0.42.0, modules/nats v0.43.0) — breaking in `wait.ForSQL` (callback now takes `network.Port`) and `ImageProvider` (new `PullImageWithPlatform`); Astrate's `internal/testutil/pg.go` looks unaffected but it is test-only anyway.
  - `golang.org/x/crypto` v0.53.0 → v0.56.0 — x/crypto keeps API compatibility; used only for bcrypt in `internal/auth`.
  Note: `govulncheck` and `golangci-lint` are still not installed on the Pi (existing entries above), so the vulnerability and lint checks remain unavailable.

---

- **milestone 2.0 looks complete, verify and cut the tag** — all 11 `milestone-2.0` issues
  CLOSED (#23–#27, #37, #39–#43), no open issues, no new gaps after re-checking upstream
  astarte_flow block catalog against `internal/flow/` + git log (MQTT/HTTP source/sink,
  json_path_map, pure-transform set, virtual_device_pool, container block MVP, flow API,
  durable named flows all landed). `mule/queue` branch is stale — see entry below.
  (Milestones recipe run, 2026-09-03.)

- **v3.0 decide #68** (`upstream-parity`, `mule-blocked`): `async_operation=false` acceptance is already decided (Astrate stays always-sync; accept and ignore on mutating endpoints: housekeeping realm create/delete, realm-management interface install/update/delete, trigger/policy delete in `internal/realm/http.go`). The mule's code+test task (`issue-68-async-operation-accepted`) failed gates and is parked. Needs re-queue or gate fix.

---

- **v3.0 decide #67** (`upstream-parity`): `required` and `encrypted` mapping fields exist on upstream master but not at 1.2.2/1.3.0 (`required` = missing key at runtime → `missing_required_mapping` on object-aggregated mappings; `encrypted` ties into keyAgreement #51). Decide now: reject (strict 1.2.2), ignore-and-accept, or implement `required` early. Issue is OPEN, no decision on record, not yet escalated here.

---

- **v3.0 decide #91** (`upstream-parity`, `mule-blocked`): upstream serves the unauthenticated health check at `GET /pairing/health` and returns 404 for `GET /pairing/v1/{realm}/health`, but Astrate only implemented the `v1/{realm}` shape (issue #71, commit 5111e3c) — the dashboard polls per-service health at the service root, so the current route may be invisible to it. Decide: move/duplicate the handler at `GET /pairing/health` (keep or drop the v1/{realm} variant), and re-check the original realmcfg-02 measurement (a port override may have skewed it).

---

- **`govulncheck` is not installed on the Pi**, so the hygiene recipe's highest-priority check (reachable dependency vulnerabilities) could not run here. `go install golang.org/x/vuln/cmd/govulncheck@latest` would close it — I did not pick a version because the pin is a decision. The `go list -m -u` sweep shows only version-skew on transitive deps (cloud/azure SDKs, cel/expr, clickhouse); none corresponds to an identified fix this repo needs, so I proposed no bumps. (Hygiene recipe run, 2026-09-02.)

---

- **`docs/site/appengine-api.md:87-88` documents `GET` and `DELETE /appengine/v1/<realm>/groups/<name>`**, but no such routes exist in `internal/appengine/http.go` (only `/groups/{group}/devices` and the nested device ops — there is no `GET` or `DELETE` on `/groups/{group}`). The endpoint the page describes is absent from the code. Decide: drop the two lines, mark them not-yet-implemented, or have the code grow them. (Docs-sync recipe run, 2026-09-02.)

---

- **`Router.Submit` TOCTOU on `closed` flag** (`internal/flow/router.go:113-120`): Submit
  reads `r.closed` under the mutex, drops the lock, then sends on the channel. A concurrent
  `Drain` could close the channel between the unlock and the send, causing a send-on-closed
  channel panic in the caller (not recovered by `processOne`'s defer). Fix options: (a) hold
  the lock during channel send (hot-path perf hit), (b) use a select with `r.quit` instead of
  sending on a closed channel, (c) accept the risk and document that Submit must not be called
  after Drain. Needs a design call — not proposing a task. (Code review 2026-09-02.)

---

- **#87 `lua_map` — needs embedded Lua runtime, parked.** Design/implementation decision: embedding a Lua VM in Astrate is not machine-checkable by the mule. Consider closing if Lua flow support is not on any active roadmap.
- **#78 FDO device onboarding — milestone-4.0, investigation phase.** Too large for a single mule task; the investigation work (reading upstream's TO2 handling, inventorying endpoints/schema/keys) is a multi-session project. Parking for now until the v3.0 queue clears and this becomes the next milestone target.
- **#1 stale — "Provide an Open Source IoT Platform unironically" (wontfix).** Has been open since the repo's founding with no activity. Consider closing.

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
  (source `.mule/research/upstream-parity-2026-08-22.md`) into a plan. But I should not start
  that triage from this stale branch. Action needed: **merge/rebuild `mule/queue` on top of
  `origin/main`**, then the next milestones run targets v3.0 with the correct research doc.
  (Milestones recipe run, 2026-09-01.)

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

This is very revealing. The `for-giulio.md` file contains crucial context. Two milestones recipe runs have already occurred:

**2026-09-01 run:** Discovered that the current branch's `.mule/milestones.md` is stale — v2.0 is already DONE on `main`, and the real current target is v3.0. It pointed out `mule/queue` is 120 commits behind `origin/main`.

**2026-09-03 run:** Filed nothing for v2.0 (all issues closed), and noted v3.0 scope is "upstream Astarte feature parity 1.3/1.4".

The 2026-09-03 entry says "milestone 2.0 looks complete, verify and cut the tag" was ALREADY proposed.

Let me check the state of the current repo — is this run's branch stale too? The key question: is `.mule/milestones.md` here stale (showing v2.0 as "not started") like it was on 2026-09-01?

Looking at the current milestones.md I read: it shows v2.0 "Status: not started", v3.0 "CLEA architecture, piece 1 — not yet decided". But the for-giulio.md entry from 2026-09-01 says on `main`, v2.0 is DONE and v3.0 is "upstream Astarte feature parity 1.3/1.4".

So the current checked-out milestone file IS stale — it matches the stale branch state described on 2026-09-01, NOT main. The for-giulio entries say the real target is v3.0 (upstream parity), but the file here says v3.0 is CLEA TBD.

Let me check the git state and what's on origin/main.
$ git branch --show-current 2>&1; echo "=== diff count vs main ==="; git rev-list --count HEAD..origin/main 2>&1; echo "=== milestones on main ==="; git show origin/main:.mule/milestones.md 2>&1 | head -80
mule/queue
=== diff count vs main ===
120
=== milestones on main ===
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

This confirms the situation precisely. The current branch `mule/queue` is **120 commits behind `origin/main`**, and its `.mule/milestones.md` is stale — it still shows v2.0 as "not started" and v3.0 as "CLEA architecture, piece 1 — TBD". But on `main`:

- **v2.0 is DONE** (2026-07-29), all 11 milestone-2.0 issues closed.
- **v3.0 is "upstream Astarte feature parity 1.3/1.4"**, scoped on main 2026-08-22.

So the recipe as it will be executed on this stale branch points at v2.0. But this exact situation was already encountered and recorded on 2026-09-01 and 2026-09-03. The for-giulio entries explicitly note the current branch's `.mule/milestones.md` is stale and the recipe fielded nothing for v2.0.

Let me confirm what main's v3.0 section says, so I understand the real target, but the key point: the recipe read the stale file and pointed me at v2.0. Per protocol, I should not start v3.0 work from this stale branch, and v2.0 has no gaps (all closed). This mirrors the existing recorded situation.

Let me check main's v3.0 section to be thorough, and confirm the git state of milestones file.
$ git show origin/main:.mule/milestones.md 2>&1 | sed -n '80,140p'
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

$ gh issue list --label milestone-2.0 --state open --limit 50 2>&1; echo "open-count-exit:$?"
open-count-exit:0

Now the picture is completely clear and matches exactly the situation already recorded in `for-giulio.md` on 2026-09-01 and 2026-09-03.

**The situation:**
- Current branch `mule/queue` is **120 commits behind `origin/main`** and its `.mule/milestones.md` is stale — it shows v2.0 as "not started" and v3.0 as CLEA TBD.
- On `main`, **v2.0 is DONE** (2026-07-29), all 11 `milestone-2.0` issues closed, no open gaps.
- The **real current target is v3.0** ("upstream Astarte feature parity 1.3/1.4"), whose first recipe job is to triage `upstream-parity` backlog #47–#89 using `.mule/research/upstream-parity-2026-08-22.md`.

**Per the recipe on this stale branch:**
- Recipe step 1 → targets v2.0 (first non-DONE section).
- Step 3 → `gh issue list --label milestone-2.0 --state all` shows all closed; no new gaps to file for v2.0.
- Step 5 applies: `milestone-2.0` open is empty, no new gaps found. But the "looks complete, verify and cut the tag" line was **already proposed** in `for-giulio.md` on 2026-09-03.

So there is **nothing new to file or propose this run**. The correct, honest action is to report that this branch is stale (the recipe keeps pointing me at already-DONE v2.0), that v2.0 has no gaps, and that the real v3.0 work cannot be done from this stale branch — exactly as the existing for-giulio entries record.

I should not re-add the "2.0 looks complete" line (already present). I should not duplicate the "merge/rebuild mule/queue" action (already present, 2026-09-01). Since these existing entries already capture everything, and my job is a proposal job, there is genuinely nothing new to append. Per MULE.md and the recipe, finding nothing wrong / nothing new is a good outcome.

Let me verify the milestones recipe has nothing more and double-check there's no `milestone-2.0` sub-task line in todo that needs attention.
$ rg -n "milestone" .mule/todo.md 2>/dev/null || echo "no milestone lines in todo.md"; echo "=== todo tail ==="; tail -8 .mule/todo.md 2>/dev/null; echo "=== research doc present? ==="; ls .mule/research/upstream-parity-2026-08-22.md 2>&1
no milestone lines in todo.md
=== todo tail ===
- [!] lint-flow-stutter: in `internal/flow`, rename exported `FlowStatus`/`FlowConfig`/`FlowMessage` to `Status`/`Config`/`Message` to satisfy revive's stutter rule; touches flow.go, message.go, block.go, graph.go, router.go and their tests — a public API rename, so no wire/JSON change (the Type names are unchanged), it is purely Go identifier renaming. [auto] — BLOCKED: gates failed
- [!] serverdata-aftercommit-test: in `internal/engine/serverdata_test.go`, add a test that wires an `afterCommit` observer on a `newWiredRig` engine, calls `PublishServerValue` for a server-owned property, and asserts whether the observer fires. The batch path fires `afterCommit` for every committed op (batch.go:168), but `PublishServerValue` bypasses the batch entirely (serverdata.go:99-162) — the current tests never check this. The test documents the intended behaviour (exclusion or parity with the data path) and fails if the contract changes. [auto] — BLOCKED: gates failed
- [!] consumer-properties-plaintext-fallback: in `internal/engine/control.go` `sendConsumerProperties` (line 204), when `devices.peek` returns nil (device not cached after eviction), the code defaults to zlib compression — but a device that declared `purge_properties_compression_format: plaintext` via the capabilities topic cannot decompress zlib. Fix by falling back to a store read of the device's capability state (or persisting the capability), and add a test in `internal/engine/capabilities_test.go` that evicts the cache entry for a plaintext device and asserts the purge payload is uncompressed. [auto] — BLOCKED: gates failed
- [!] object-path-ok-comment: in `internal/engine/data.go` `objectPathOK` (line 330), the loop over `ci.ObjectLeaves` returns after checking the first entry. Add a comment explaining why one leaf suffices (upstream validation guarantees uniform prefix depth), so future readers do not mistake the early return for a bug. [auto] — BLOCKED: gates failed
- [!] housekeeping-unit-tests: in `internal/housekeeping`, add a unit test file (no integration build tag) that exercises `CreateRealm` validation paths — blank name returns ErrValidation, blank JWT key returns ErrValidation, negative registration limit returns ErrValidation — using a mock store and sealer. The package currently has only integration tests gated behind `//go:build integration`, so these paths are never tested on the Pi. [auto] — BLOCKED: gates failed
- [!] store-downsample-lttb: in `internal/store`, implement the TODO at store.go:135 — when `s.hasToolkit` is true, switch `Downsample` from `time_bucket+avg` to toolkit `lttb()` downsampling. The probe already records availability; the time_bucket path in datastreams.go is the always-working default. Add a test that exercises both code paths (mocking `hasToolkit`). [auto] — BLOCKED: gates failed
- [!] pairing-info-verify-ratelimit: in `internal/pairing/http.go`, add per-IP (and per-device) token-bucket rate limits to `handleInfo` (line 197) and `handleVerify` (line 238), which currently run a full bcrypt `CompareHashAndPassword` (via `authenticateDevice`, service.go:394) on any bearer string with no throttle — the sibling `handleCredentials` that does identical bcrypt work is limited at burst-5, and DESIGN §4.5 says "pairing endpoints" should be rate-limited. Reuse the credentials limiter defaults (1 req/s, burst 5) and extend `TestPairingHTTPRateLimit` to assert a burst of fake-secret requests to `GET /devices/{id}` and `POST .../verify` returns 429. [auto] — BLOCKED: gates failed
- [!] docs-sync-realm-policies-422: in docs/api/astarte_realm_management_api.yaml, add the undocumented 422 Unprocessable Entity response to DELETE /policies/{name} — the code (internal/realm/service.go:398) returns ErrValidation when a policy is still referenced by a trigger, which maps to 422 via writeError; the YAML only lists 204, 401, 404, 500. [auto] — BLOCKED: gates failed
=== research doc present? ===
ls: cannot access '.mule/research/upstream-parity-2026-08-22.md': No such file or directory

This fully confirms it. The `upstream-parity-2026-08-22.md` research doc that the real v3.0 work needs is **not present on this branch** (it lives on `main`). This branch's todo.md has no milestone lines.

**Conclusion:** This is the identical situation to the 2026-09-01 and 2026-09-03 runs, already recorded verbatim in `for-giulio.md`. The stale `mule/queue` branch (120 commits behind `main`) keeps pointing the milestones recipe at a v2.0 that is already DONE on `main`, with all issues closed and no new gaps. Filing anything for v2.0 would duplicate closed work. Starting v3.0 triage is impossible from this branch because the source research doc doesn't exist here.

Per the recipe and MULE.md principles, the honest, complete deliverable is: **file nothing, propose nothing new**. The two earlier entries in `for-giulio.md` already capture the action needed (rebuild `mule/queue` on `origin/main`). Re-adding them would be noise/duplication.

Let me verify that the two key existing entries are still accurate and don't need updating — the 2026-09-01 entry (rebuild mule/queue) and 2026-09-03 entry (2.0 complete). Both are still accurate given the branch is still 120 behind. Since the recommendation hasn't been acted on yet, I should not re-add a duplicate. But per the recipe step 5, the milestone "looks complete" proposal was already made. Nothing new.

I'll write my report. Since there's nothing to append to todo.md or for-giulio.md (all captured), I will not modify any files. The deliverable for this run is the report itself, mirroring [readonly] behavior.

The recipe says this is a proposal job, but the correct result this run is "propose nothing" because the milestone's state is already fully covered and the branch is stale. My report should state this clearly, following the MULE.md report format.

I made no file changes and no `gh` calls (the recipe's one permitted mutation is filing new issues, and there are none to file). Let me produce the report.

Done: Milestones recipe on stale `mule/queue` branch — v2.0 is already DONE on `main` (all 11 milestone-2.0 issues closed, no open gaps), and this branch (120 commits behind `origin/main`) still points the recipe at v2.0 while its v3.0 source doc (`.mule/research/upstream-parity-2026-08-22.md`) doesn't exist here, so there is nothing new to file or propose — filed nothing.
Files: (none — no file or gh changes this run)
Verified: `gh issue list --label milestone-2.0 --state all` (all 11 closed, none open) -> pass; `git rev-list --count HEAD..origin/main` = 120; `git show origin/main:.mule/milestones.md` (v2.0 DONE, v3.0=upstream parity) -> pass
Unsure: nothing — this mirrors the already-recorded 2026-09-01/09-03 findings in `.mule/for-giulio.md` (rebuild `mule/queue` on `origin/main`; then v3.0 triage of #47–#89 becomes the real next target)
Follow-ups: none — the needed action (merge/rebuild `mule/queue` on top of `origin/main`) is already recorded in `.mule/for-giulio.md`; re-adding it would duplicate. No new todo.md lines, no duplicate issue-filing, no second "2.0 looks complete" line (already there).

**What I did:** Read `.mule/MULE.md`, `.mule/recipes/milestones.md`, `.mule/milestones.md`, and `.mule/for-giulio.md`, then checked live GitHub state and the git branch.

**Result: filed nothing — and that's correct.**

The recipe walks the first non-DONE milestone, which on this `mule/queue` branch is **v2.0** — but this branch is **120 commits behind `origin/main`** and its `.mule/milestones.md` is stale. On `main`, v2.0 is marked **DONE (2026-07-29)** and every one of its 11 `milestone-2.0` issues is CLOSED with no open gaps. So per recipe step 5 there's nothing to file for v2.0 without duplicating closed work.

The real next target on `main` is **v3.0 (upstream Astarte parity)**, but its source doc (`.mule/research/upstream-parity-2026-08-22.md`) doesn't exist on this branch, so triage can't start from here.

Both of these facts were already recorded verbatim in `.mule/for-giulio.md` on 2026-09-01 and 2026-09-03, including the required action (merge/rebuild `mule/queue` off `origin/main`, then the next run targets v3.0). Re-adding them would be duplication, so I appended nothing to `todo.md` or `for-giulio.md`, changed no files, and ran no `gh` mutations.
