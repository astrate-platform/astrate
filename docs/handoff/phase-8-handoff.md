# Handoff Prompt — Start Phase 8

Use this prompt to start the next documentation session.

````text
I'm working on the astrate project documentation in ~/astrate.

Before changing files, read:
- ~/astrate/docs/AI-DOCUMENTATION-WORKFLOW.md
- ~/astrate/docs/handoff/README.md
- ~/astrate/docs/handoff/phase-7-memory.md
- ~/astrate/docs/handoff/phase-8-memory.md

Known state:
- Phases 1-7 are complete and committed to main (cc34e6a).
- Live site: https://astrate-platform.github.io/astrate/
- OPERATIONS.md and JSON-PAYLOAD-PROFILE.md already exist and are substantive.
- 4 new pages needed: Configuration Reference, Contributing, Troubleshooting, Migration from Astarte.
- MkDocs nav lives in docs/mkdocs.yml on main (16 entries currently).
- Build: cd docs && make sync && make build
- Important: the mkdocs.yml, Makefile, and narrative pages exist on main but not on the worktree-m12-07-lttb branch. Check out main or work against it.

Task:
- Implement Phase 8 per docs/handoff/phase-8-memory.md (4 new pages + updates + verification).

Sub-phases in order:
1. 8A: Create Configuration Reference (docs/CONFIGURATION-REFERENCE.md + docs/site/configuration-reference.md)
   - Read internal/config/config.example.toml, docs/DESIGN.md §5.1, docs/OPERATIONS.md §Configuration
   - Enumerate every TOML key: section, key name, type, default value, env override, required/optional, description
   - Include precedence rule: defaults < TOML < ASTRATE_* env
   - Add to mkdocs.yml nav

2. 8B: Create Contributing / Dev Guide (docs/CONTRIBUTING.md + docs/site/contributing.md)
   - Read ROADMAP.md §0.2 (test tiers T1-T5), root Makefile, .golangci.yml
   - Cover: prerequisites, dev setup (clone, make tools, make up), test tiers with commands, code style, how to contribute
   - Put the authoritative copy at repo root (CONTRIBUTING.md), sync copy to docs/site/

3. 8C: Create Troubleshooting (docs/site/troubleshooting.md)
   - Read DESIGN.md §6, OPERATIONS.md, COMPATIBILITY.md
   - Cover: cert errors (wrong CA, expired, serial mismatch), DB connection refused, migration failures, MQTT connection refused, master key lost, Dashboard 401/403, rate limit 429, payload rejected, device not connecting

4. 8D: Create Migration from Astarte (docs/site/migration-from-astarte.md)
   - Read DESIGN.md §0 (comparison), §1.1 (service mapping), COMPATIBILITY.md
   - Cover: pre-migration checklist, data export, schema mapping, device behavior during cutover, rollback plan
   - Mark as preliminary (no real-world migration yet)

5. 8E: Update existing files
   - Edit docs/mkdocs.yml: add nav entries for all 4 new pages
   - Edit docs/OPERATIONS.md: add cross-refs to config reference and troubleshooting at top
   - Edit docs/site/deployment.md (on main): add cross-refs if appropriate

6. 8F: Verify and handoff
   - cd docs && make sync && make build — must pass
   - Update docs/handoff/phase-8-memory.md and phase-8-handoff.md with what was actually done
   - Update docs/handoff/README.md

Rules:
- Read source/tests before documenting behavior (especially config.example.toml for config ref).
- Keep each page focused and scannable (tables, code blocks, short paragraphs).
- Mark uncertain or preliminary content explicitly.
- Do not edit DESIGN.md, ROADMAP.md, or existing narrative pages unless adding cross-refs.
- Verify build before declaring complete.
- At the end, provide: files changed, verification results, risks, and tell the user which handoff file to read for the next session.
````
