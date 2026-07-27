# Documentation Handoff Index

This directory contains local handoff files for AI-assisted documentation work.

Start every documentation session by reading:

1. `docs/AI-DOCUMENTATION-WORKFLOW.md`
2. This file
3. The latest `phase-N-memory.md`
4. The latest `phase-N-handoff.md`
5. `g-mind/08 - Progetti/Piano documentazione astrate.md`, if available

## Current State

Latest completed phase: Phase 10 - Fix Observability Metrics + Add Cross-References.

Latest files:

- `docs/handoff/phase-10-memory.md` (Phase 10 completed)
- `docs/handoff/phase-11-handoff.md` (contains the next-session prompt)
- `docs/handoff/phase-9-memory.md`
- `docs/handoff/phase-9-handoff.md`

Current status:

- Docs on `main`: 22 narrative pages, Swagger UI, MkDocs config, GitHub Actions workflow, OpenAPI specs.
- GitHub Pages enabled in repo Settings (source: GitHub Actions).
- Live site: https://astrate-platform.github.io/astrate/
- Code milestones M0-M9 are on `worktree-*` branches, not yet merged to `main`.

Phase 11 plan (next):

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
