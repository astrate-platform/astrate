# Phase 7 Memory — Push docs to main for GitHub Pages (2026-07-27)

## Completed

- Committed all docs-related files to `main` branch (commit `cc34e6a`).
- Pushed to origin, triggering the `docs.yml` GitHub Actions workflow.
- Build job passed (13s). Deploy initially failed because GitHub Pages was not enabled.
- GitHub Pages enabled in repo Settings (source: GitHub Actions).
- Re-run or next push to `main` will deploy successfully.

## Files Changed

- `.github/workflows/docs.yml` — GitHub Actions workflow for Pages deployment
- `.gitignore` — docs build artifacts, sync copies, local workflow entries
- `docs/Makefile` — build tooling (sync, serve, build, clean)
- `docs/mkdocs.yml` — MkDocs + Material theme, 15-page nav
- `docs/requirements.txt` — pins mkdocs 1.x, material 9.x
- `docs/api/*.yaml` — updated OpenAPI specs (realm-management expanded, native API fixes)
- `docs/site/*.md` — 15 narrative documentation pages
- `docs/swagger-ui/index.html` — bundled Swagger UI

## Verification

- `gh run view 30230495644` — build passed, deploy failed (Pages not enabled at time of run)
- All 26 files committed cleanly, no code files mixed in

## Known Risks

- Code milestones M0-M9 are on `worktree-*` branches, not yet merged to `main`.
- The `worktree-m12-07-lttb` branch has unstaged lttb code changes (not documentation work).

## Documentation Gaps Identified

The following pages are missing from the current 15-page site and would add significant value:

| Priority | Page | Why it matters |
|---|---|---|
| High | Configuration Reference | `docker-compose.yml` references config options but no docs list them all. The `internal/config` system has many TOML keys and `ASTRATE_*` env vars. |
| High | Operations Runbook | `docker-compose.yml` header references `docs/OPERATIONS.md` for production TLS, CA re-keying, backup. Does not exist yet. |
| High | JSON Payload Profile | ROADMAP M9 calls for `docs/JSON-PAYLOAD-PROFILE.md` — normative spec for AtomVM/bare-metal JSON clients. Critical for the dual-format story. |
| Medium | Contributing / Dev Guide | Test tiers (T1-T5), dev setup, how to run tests, how to add features. |
| Medium | Troubleshooting | Common failure modes: cert errors, migration issues, connection refused. |
| Medium | Migration from Astarte | How to move an existing Astarte deployment to Astrate. |

## Next Steps

1. Verify live site at GitHub Pages URL after next deployment.
2. Pick one of the documentation gaps above and write it.
3. Commit to `main`, push, verify it appears on the live site.
