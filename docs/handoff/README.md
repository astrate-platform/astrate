# Documentation Handoff Index

This directory contains local handoff files for AI-assisted documentation work.

Start every documentation session by reading:

1. `docs/AI-DOCUMENTATION-WORKFLOW.md`
2. This file
3. The latest `phase-N-memory.md`
4. The latest `phase-N-handoff.md`
5. `g-mind/08 - Progetti/Piano documentazione astrate.md`, if available

## Current State

### General project work (preferred start for non-docs sessions)

- **`docs/handoff/session-2026-07-29-flow-factory-handoff.md`** — **primary next-session
  prompt for product work.** Flow v2.0 is **mostly wired**:
  - **On `main`:** runtime, source pump, AstarteSource, pipeline store, **block factory**,
    catalog (`astarte_source` / `filter` / `map` / sinks), process wiring, `/flow/v1` API,
    pokemon example (#38).
  - **Still open:** flows-table decision (for-giulio), parity audit; Legion/#20, mule items.
  - Local commits may be **ahead of origin** — push only with confirmation.

- `docs/handoff/session-2026-07-29-handoff.md` — older (source pump only); superseded
  for “what next” by the flow-factory handoff.

### Pokémon agent example — **completed / merged**

- Branch `feat/pokemon-agent` merged to `main` via **PR #38** (`bb67a95`, 2026-07-29).
- Live path proven: intro skip → overworld move → leave Red's House 2F → 1F → Pallet
  Town (`POKEMON_GUIDANCE=light` + Astrate ControlCommand bus).
- Code: `examples/pokemon-agent/` (+ `docs/site/pokemon-agent.md`, handoffs
  `pokemon-agent-*.md`).
- Optional follow-ups only (not blocking product Flow): pure-LLM leave-house,
  mid-dialog flags, endurance demo — see `pokemon-agent-handoff.md` if revisiting.
- **Do not treat pokemon as the default “what next”** for Astrate platform work.

### Documentation phases

Latest completed phase: Phase 10 - Fix Observability Metrics + Add Cross-References.

Latest files:

- `docs/handoff/phase-10-memory.md` (Phase 10 completed)
- `docs/handoff/phase-11-handoff.md` (docs-only next-session prompt)
- `docs/handoff/phase-9-memory.md`
- `docs/handoff/phase-9-handoff.md`

Current status:

- Docs on `main`: 22 narrative pages, Swagger UI, MkDocs config, GitHub Actions workflow, OpenAPI specs.
- GitHub Pages enabled in repo Settings (source: GitHub Actions).
- Live site: https://astrate-platform.github.io/astrate/
- Platform code (M0–M9 lineage + Flow v2.0 through factory/catalog/API/transforms) is on
  `main` (may need push; see flow-factory handoff).

Phase 11 plan (docs, when chosen):

- Review remaining narrative pages for accuracy against current source.
- Create remaining pages from ROADMAP milestones if needed.
- Consider adding ROADMAP.md to MkDocs nav.

Phase 2 local preview (Swagger UI):

- Serve from `~/astrate/docs`, then open `http://localhost:9090/swagger-ui/`.
- Do not serve from `~/astrate/docs/swagger-ui`; the `../api/` YAML paths will not resolve.
- Avoid port `8080` while the Lenovo Legion Go SSH tunnel is active.

Phase 4 local preview (MkDocs site):

- `cd docs && make serve`, then open `http://localhost:8000/`.
- Run `make sync` before build to refresh DESIGN.md, ROADMAP.md, api/, swagger-ui/ copies.

## Operating Rule

Each phase must end with two files:

- `phase-N-memory.md` - factual state, files changed, evidence, verification, risks, and next phase.
- `phase-N-handoff.md` - contains the prompt for the next session. At the end of each phase, the AI tells the user: "read `docs/handoff/phase-(N+1)-handoff.md` to start the next session."

Keep handoff files in English, even if the user speaks Italian in chat.
