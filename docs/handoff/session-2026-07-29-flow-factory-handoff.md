# Handoff — Flow factory + process wiring + HTTP API (2026-07-29)

Use this prompt to continue general Astrate work (not a docs-only phase).

````text
I'm working on the Astrate project in ~/astrate on branch main.

Before changing files, read:
- ~/astrate/docs/handoff/session-2026-07-29-flow-factory-handoff.md  (this file)
- ~/astrate/docs/handoff/README.md
- ~/astrate/.mule/milestones.md
- ~/astrate/.mule/for-giulio.md  (open decisions only)

## Split status: committed vs worktree (important)

**On `main` (committed / merged):**
- HEAD includes PR #38 (pokemon example) and earlier Flow work through source pump:
  runtime core, router, Source/Stopper pump (#37), AstarteSource, pipeline store +
  migration 000008, milestones v2.0 “in progress” table for those pieces.
- Pokémon agent is **done as a project line**: merged, not the default next task.

**In the working tree (uncommitted — restored via stash pop 2026-07-29):**
Flow v2.0 gaps 1–3 from the earlier same-day session, still local:
  1. **Block factory** — `internal/flow/factory.go`: `Registry`, `Instantiate` (topo
     order), `ParseDefinition`, `FlowPipelineID(realm, name)`.
  2. **Built-in catalog** — `internal/flow/blocks/catalog.go`: `astarte_source`,
     `null_sink`, `log_sink` (named wrapper preserves Source/Stopper).
  3. **Process wiring** — `cmd/astrate`: `flow.NewManager()` + `blocks.DefaultRegistry()`
     + `e.Bus()`; `flowapi` mounted; shutdown = HTTP → broker → **flows** → engine.
  4. **Operator HTTP API** — `internal/flowapi` (`service.go`, `http.go`):
     - `GET/POST /flow/v1/{realm}/pipelines`, `GET/PUT/DELETE .../pipelines/{name}`
     - `GET/POST /flow/v1/{realm}/flows`, `GET/DELETE .../flows/{name}`
     - auth: realm JWT `a_rma`; envelope via `pkg/astarteapi`
     - start body: `{"pipeline":"<name>"}`; create body: `{"name","definition"}`
  - Tests: `factory_test.go`, `blocks/catalog_test.go`, `flowapi/*_test.go` (unit).
  - `go test ./internal/flow/... ./internal/flowapi/... ./cmd/astrate` green after pop.
  - `.mule/milestones.md` in tree marks factory/catalog/wiring/API as landed; remaining
    gaps are transforms, flows-table decision, parity audit.

Do not discard this worktree state. Prefer commit + PR when ready (after review).

Known context:
- Stash was `wip-flow-factory-20260729` (from main@c5fa06b); popped cleanly onto
  main@bb67a95 after pokemon merge.
- #37 still open on GitHub until closed if still open (source pump already on main).
- Scope policy: reimplement FOSS worth recreating; prefer integrate when it fits.
- Active for-giulio: dormant triggers (#20 Legion bench); group-scoped triggers (#17); #10.
- Mule: Pi planner+cron; Legion for race-check / #20.

## Practical “what should I do?” menu

| Priority | Action |
|---|---|
| **0** | **Commit/push** the uncommitted factory + catalog + flowapi + cmd wiring (this WIP) |
| **1** | **Native transform block(s)** for the catalog (filter/map) — last big “usable pipeline” gap |
| **2** | Decide **flows table** vs in-memory only (milestone gap; escalate if needed) |
| **3** | **Parity audit** vs astarte_flow operator concepts; file residual gaps |
| **4** | **#20** bench on Legion → decide dormant triggers |
| **5** | Mule: **#13**, auto v1.3.0 checks, race-check; review when it lands |
| **6** | **#17–19** wire/authz parity when wanted |
| **7** | Close #37 on GitHub if still open (source pump already merged) |

Short answer: factory/API/wiring are ready in the tree but not committed. Land that,
then transform + decision/audit items. Pokemon is finished on main — ignore unless asked.

Rules:
- Prefer small, test-backed changes; run relevant go test / integration when practical.
- Do not commit unless asked; do not push without confirmation for shared remotes.
- At the end of the session, update this handoff (or add a dated successor) and
  tell the user which handoff file to read next.
````

## Menu only (quick scan)

| Priority | Action |
|---|---|
| **0** | Commit uncommitted Flow factory/API/wiring WIP |
| **1** | Transform block(s) in catalog |
| **2** | flows table decision |
| **3** | Parity audit |
| **4** | **#20** Legion bench → dormant triggers |
| **5** | Mule #13 / v1.3.0 / race-check |
| **6** | **#17–19** when wanted |
| **7** | Close #37 if open |

## Landed — split

### Committed on `main`
Source pump, Stop on teardown, pipeline store, AstarteSource, Flow runtime core, pokemon example (#38).

### Uncommitted in worktree (paths)
| Piece | Path |
|---|---|
| Factory | `internal/flow/factory.go`, `factory_test.go` |
| Catalog | `internal/flow/blocks/catalog.go`, `catalog_test.go` |
| API | `internal/flowapi/{service,http}.go` + tests |
| Bootstrap | `cmd/astrate/main.go` |
| Milestones | `.mule/milestones.md` (gaps 1–3 marked landed in WIP) |
| Handoff | this file + `docs/handoff/README.md` |

## Pokémon (closed)

- PR **#38** merged 2026-07-29 → `examples/pokemon-agent/` on main.
- Handoffs: `pokemon-agent-handoff.md` / `pokemon-agent-memory.md` (optional follow-ups only).

## Also open (context)

| Issue | Notes |
|---|---|
| #37 | source pump — on main; close on GitHub if still open |
| #28 | Edgehog compatibility (readonly, mule) |
| #20 | previous-value lookup bench (readonly; Legion) |
| #13 | dashboard error_names audit (readonly) |
| #17–19 | enhancements |
| #10 | design question |
| #1 | wontfix — ignore |
