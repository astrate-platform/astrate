# Phase 4 Memory — Static site generator (2026-07-27)

## Completed

- Wired 15 narrative Markdown pages into MkDocs with Material theme.
- Created `docs/mkdocs.yml`, `docs/requirements.txt`, `docs/Makefile`.
- Copied DESIGN.md, ROADMAP.md, api/, swagger-ui/ into `docs/site/` so MkDocs can serve them.
- Fixed 4 relative links in `docs/site/index.md` (removed `../` prefix).
- Verified clean build: `mkdocs build` succeeds, all 15 pages + DESIGN + ROADMAP + api + swagger-ui present in output.

## Files created

- `docs/mkdocs.yml` — MkDocs config: Material theme, nav for 15 pages, markdown extensions
- `docs/requirements.txt` — mkdocs + mkdocs-material pinned
- `docs/Makefile` — targets: install, sync, serve, build, clean

## Files changed

- `docs/site/index.md` — fixed 4 links: `../DESIGN.md` → `DESIGN.md`, etc.
- `docs/site/DESIGN.md` — copied from `docs/DESIGN.md` (was symlink, now real file)
- `docs/site/ROADMAP.md` — copied from `docs/ROADMAP.md`
- `docs/site/api/` — copied from `docs/api/`
- `docs/site/swagger-ui/` — copied from `docs/swagger-ui/`

## Verification

- `cd docs && mkdocs build -f mkdocs.yml` — exit code 0, 15 nav pages + 2 extra pages built
- All nav links, DESIGN, ROADMAP, api/, swagger-ui/ links resolve in built output
- `site-dist/` contains: index.html, 15 page dirs, DESIGN/, ROADMAP/, api/, swagger-ui/, assets/, search/

## Known Risks

- DESIGN.md and ROADMAP.md are copied into `docs/site/` and will go stale when the originals change. The `make sync` target handles this but must be run before build.
- `site-dist/` is a build artifact — should be added to `.gitignore` if not already.
- The MkDocs 2.0 upstream warning is informational only; current MkDocs 1.x is stable.
- No GitHub Pages or CI deployment configured yet — local `mkdocs serve` only.

## Next phase

- Phase 5: Add `site-dist/` to `.gitignore`, optionally set up GitHub Pages deployment, or begin code generation milestones (M0-M9 per ROADMAP.md).
