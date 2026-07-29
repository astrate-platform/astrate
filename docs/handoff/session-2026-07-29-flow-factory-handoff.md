# Handoff — Flow after factory + transforms (2026-07-29)

Use this prompt to continue general Astrate work (not a docs-only phase).

````text
I'm working on the Astrate project in ~/astrate on branch main.

Before changing files, read:
- ~/astrate/docs/handoff/session-2026-07-29-flow-factory-handoff.md  (this file)
- ~/astrate/docs/handoff/README.md
- ~/astrate/.mule/milestones.md
- ~/astrate/.mule/for-giulio.md  (open decisions only)

## Split status: committed vs remaining

**On `main` (committed this session / earlier same day):**
- Flow runtime core, source pump (#37 lineage), AstarteSource, pipeline store.
- **Block factory** — `internal/flow/factory.go` (`Registry`, `Instantiate`, `ParseDefinition`).
- **Built-in catalog** — `astarte_source`, `filter`, `map`, `null_sink`, `log_sink`
  (`catalog.go`, `transform.go`).
- **Process wiring** — `cmd/astrate`: Manager + DefaultRegistry + bus; flowapi mounted;
  shutdown HTTP → broker → flows → engine.
- **Operator HTTP API** — `internal/flowapi` `/flow/v1/{realm}/pipelines` + `/flows` (a_rma).
- Pokémon agent example merged (PR #38) — not the default next task.

**Commits to know:** factory/API/wiring landed first; filter/map transforms next
(check `git log --oneline -5` for SHAs). Branch may be **ahead of origin** — push only
with confirmation.

Known context:
- #37 still open on GitHub until closed if still open (source pump already on main).
- Scope policy: reimplement FOSS worth recreating; prefer integrate when it fits.
- Active for-giulio: **flows table decision (new)**; dormant triggers (#20 Legion);
  group-scoped triggers (#17); #10.
- Mule: Pi planner+cron; Legion for race-check / #20.

## Practical “what should I do?” menu

| Priority | Action |
|---|---|
| **0** | **Push** local Flow commits if not yet on origin (confirm with user first) |
| **1** | **Flows table decision** — wait for Giulio / implement after for-giulio clears |
| **2** | **Parity audit** vs astarte_flow operator concepts; file residual gaps |
| **3** | **#20** bench on Legion → decide dormant triggers |
| **4** | Mule: **#13**, auto v1.3.0 checks, race-check; review when it lands |
| **5** | **#17–19** wire/authz parity when wanted |
| **6** | Close #37 on GitHub if still open (source pump already merged) |
| **7** | Optional: more catalog blocks (HTTP sink, script/container) only if a client needs them |

Short answer: usable source→filter/map→sink pipelines work end-to-end via factory + API.
Next product decisions are **persistence of running flows** and a **parity audit**.
Pokemon is finished on main — ignore unless asked.

### filter / map config (operator reference)

`filter` (AND of set conditions; drop = zero outputs):
- `key_prefix`, `key_contains` (strings)
- `type`: integer|real|boolean|datetime|binary|string|map
- `metadata`: object string→string (all must equal)

`map` (payload unchanged; shallow-copies Metadata):
- `key`: template with `{key}` and `{metadata.<name>}`
- `set_metadata`: object string→string (merge)
- `delete_metadata`: string array

Rules:
- Prefer small, test-backed changes; run relevant go test / integration when practical.
- Do not commit unless asked; do not push without confirmation for shared remotes.
- At the end of the session, update this handoff (or add a dated successor) and
  tell the user which handoff file to read next.
````

## Menu only (quick scan)

| Priority | Action |
|---|---|
| **0** | Push Flow commits (if local-only) |
| **1** | Flows table decision (for-giulio) |
| **2** | Parity audit |
| **3** | **#20** Legion bench → dormant triggers |
| **4** | Mule #13 / v1.3.0 / race-check |
| **5** | **#17–19** when wanted |
| **6** | Close #37 if open |
| **7** | More catalog blocks only if needed |

## Landed

| Piece | Path |
|---|---|
| Factory | `internal/flow/factory.go` |
| Catalog + transforms | `internal/flow/blocks/{catalog,transform}.go` |
| API | `internal/flowapi/` |
| Bootstrap | `cmd/astrate/main.go` |
| Milestones | `.mule/milestones.md` (transform gap closed) |
| Escalation | `.mule/for-giulio.md` (flows table) |

## Pokémon (closed)

- PR **#38** merged 2026-07-29 → `examples/pokemon-agent/` on main.

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
