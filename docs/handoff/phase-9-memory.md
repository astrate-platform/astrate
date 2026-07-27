# Phase 9 Memory — Add Standalone Docs to MkDocs Nav (2026-07-27)

## Completed

- Copied OPERATIONS.md and JSON-PAYLOAD-PROFILE.md into docs/site/ with fixed internal links.
- Added 2 nav entries to mkdocs.yml (22 entries, up from 20).
- Build verified: `cd docs && make sync && make build` passes with no project-specific warnings.

## Files Changed

- `docs/site/operations.md` — new: OPERATIONS.md adapted for MkDocs (fixed links to architecture.md, compatibility.md, ROADMAP.md, configuration-reference.md, GitHub source)
- `docs/site/json-payload-profile.md` — new: JSON-PAYLOAD-PROFILE.md adapted for MkDocs (fixed links to architecture.md, GitHub source)
- `docs/mkdocs.yml` — added 2 nav entries (operations, json-payload-profile)

## Source Evidence Read

- `docs/OPERATIONS.md` — standalone operations guide (143 lines)
- `docs/JSON-PAYLOAD-PROFILE.md` — normative JSON payload spec (146 lines)
- `docs/site/payload-formats.md` — confirmed it's a summary; JSON-PAYLOAD-PROFILE.md is the detailed spec

## Verification

- `cd docs && make sync && make build` — exit 0, no project-specific warnings
- Both new pages appear in `site-dist/` as `index.html` (operations/, json-payload-profile/)
- MkDocs nav has 22 entries (up from 20)

## Known Risks

- Link from operations.md to `internal/config/config.example.toml` points to GitHub (not a local page) — acceptable since the example TOML isn't in the docs site.
- Link from json-payload-profile.md to `pkg/payload` points to GitHub — same rationale.
- OPERATIONS.md source file in docs/ root still exists and is not in nav; it's the "source of truth" that the site version was derived from.

## Next Steps

1. Commit to `main`, push, verify on live site at https://astrate-platform.github.io/astrate/
2. Phase 10 candidates: review existing narrative pages for accuracy against current source, or create remaining pages from ROADMAP milestones
