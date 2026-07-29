# Handoff — after #37 source pump + milestones refresh (2026-07-29)

Use this prompt to continue general Astrate work (not a docs-only phase).

````text
I'm working on the Astrate project in ~/astrate.

Before changing files, read:
- ~/astrate/docs/handoff/session-2026-07-29-handoff.md  (this file)
- ~/astrate/docs/handoff/README.md
- ~/astrate/.mule/milestones.md
- ~/astrate/.mule/for-giulio.md  (open decisions only)

Known state (2026-07-29, later same day):
- main base was d3dc4f3 (PR #36 merged): blocked mule issues #14 #21 #22 #24 #27 closed.
- Local uncommitted work (this session, not committed unless asked):
  - **#37 done in tree:** Flow source pump + Stop on teardown.
    - `flow.Source` (`Emit`) + `flow.Stopper` (`Stop`)
    - Manager pumps Sources → router; `BlockGraph.Run` skips Source stages
    - `StopFlow`: cancel pump → drain router → `Stop()` on Stoppers
    - AstarteSource implements `Emit` (blocking) + existing `Stop`
    - Tests: `internal/flow/pump_test.go`, Emit tests on astartesource
  - **`.mule/milestones.md` refreshed:** v2.0 status "in progress" with landed table
    + remaining gaps (block factory, cmd wiring, HTTP API, block catalog, flows table
    decision, parity audit).
- Flow still not imported from `cmd/astrate` — runtime exists, process bootstrap does not.
- **Scope policy:** Astrate reimplements free/open-source software that is worth recreating
  for the project — a wire-compatible Astarte-platform reimplementation in Go with lighter
  components, plus extras (e.g. AtomVM compatibility). Prefer integrating existing FOSS
  when it already fits (e.g. the original Astarte Dashboard; upstream Edgehog may be fine
  as-is — investigate via #28) rather than rewriting everything.
- Active for-giulio decisions: dormant trigger types (#20 bench on Legion);
  group-scoped triggers (#17); #10 external-bus intake design question.
- Mule: Pi uses planner+cron (not flat mule.timer).
- Lenovo Legion Go is awake (`ssh legion`); [legion] mule / race-check / #20 can run.
- feat/pokemon-agent branch exists separately; main is the product line of work.

## Practical “what should I do?” menu

| Priority | Action |
|---|---|
| **1** | Next Flow gap: **block factory + pipeline instantiate** (stored Pipeline → StartFlow) |
| **2** | Wire **flow.Manager + stream.Bus** into `cmd/astrate` / engine bootstrap |
| **3** | Operator **HTTP API** for pipelines/flows (small, test-backed) |
| **4** | **#20** bench on Legion → decide dormant triggers |
| **5** | Let mule chew **#13**, auto v1.3.0 checks, race-check; review when it lands |
| **6** | **#17–19** when you want more wire parity |

Short answer: Flow runtime core (including source pump) is in tree. Next is make it
operator-usable: factory from stored pipelines, process wiring, then HTTP. Close or
update GitHub **#37** when this work is committed/merged.

Rules:
- Prefer small, test-backed changes; run relevant go test / integration when practical.
- Do not commit unless asked; do not push without confirmation for shared remotes.
- At the end of the session, update this handoff (or add a dated successor) and
  tell the user which handoff file to read next.
````

## Menu only (quick scan)

| Priority | Action |
|---|---|
| **1** | Block factory + pipeline → `StartFlow` |
| **2** | Wire Manager + Bus into process bootstrap |
| **3** | Pipelines/flows HTTP API |
| **4** | **#20** bench on Legion → dormant triggers |
| **5** | Mule: **#13**, auto v1.3.0 checks, race-check |
| **6** | **#17–19** wire/authz parity when wanted |

## Also open (context)

| Issue | Notes |
|---|---|
| #37 | source pump — **implemented in working tree**; close when merged |
| #28 | Edgehog compatibility (readonly, mule) |
| #20 | previous-value lookup bench (readonly; Legion is up) |
| #13 | dashboard error_names audit (readonly) |
| #17–19 | enhancements |
| #10 | design question |
| #1 | wontfix — ignore |
