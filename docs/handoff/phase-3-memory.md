# Phase 3 Memory — Narrative site pages (2026-07-27)

## Completed

- Created 15 narrative Markdown pages in `docs/site/`.

## Files created

```
docs/site/
├── index.md                   # Overview, comparison table, page index
├── architecture.md            # Process model, service mapping, package layout, concurrency
├── data-modeling.md           # PostgreSQL + TimescaleDB schema (realms, interfaces, devices, properties, datastreams)
├── mqtt-protocol.md           # Astarte MQTT v1 wire format, ACL, message taxonomy, validation pipeline
├── pairing-and-security.md    # Three credential planes, embedded CA, pairing flows A-D
├── payload-formats.md         # BSON + JSON dual codec, sniffing, type mapping, outbound format
├── appengine-api.md           # Device management, data queries, groups, live stream
├── realm-management-api.md    # Interface CRUD, trigger CRUD, auth config
├── housekeeping-api.md        # Realm lifecycle, auto-provisioning
├── interface-schema.md        # Astarte interface JSON spec, types, mappings, compilation
├── triggers.md                # Data/device triggers, HTTP webhooks, delivery policies
├── deployment.md              # Docker Compose, bare VPS, configuration, backups
├── observability.md           # Health/readiness, Prometheus metrics, structured logging
├── compatibility.md           # Wire-identical surfaces, deliberate deviations, tested clients
└── quickstart.md              # 5-minute Docker Compose guide, bare VPS deployment
```

## Evidence read

- `docs/DESIGN.md` — primary source for architecture, data model, protocol, security, payload, and observability details.
- `docs/ROADMAP.md` — milestone structure and verification tiers referenced in architecture.
- `docs/COMPATIBILITY.md` — deviation inventory and tested client matrix used in compatibility page.
- `docs/OPERATIONS.md` — deployment profiles, configuration, backups, CA re-keying used in deployment page.
- `docs/JSON-PAYLOAD-PROFILE.md` — normative JSON profile spec used in payload-formats page.
- `docs/api/*.yaml` — API surface referenced in index and API pages.

## Verification

- `ls docs/site/` — 15 files present.
- All pages cross-reference each other and existing docs correctly.
- No build step required; pages are static Markdown.

## Delegation Log

- Architect (session): Read all handoff/context files, reviewed existing docs, defined 15-page scope.
- Mule (Big Pickle/opencode): Created all 15 Markdown files in parallel batches.

## Known Risks

- Pages are written from source documents (DESIGN.md, ROADMAP.md, etc.) which are design-phase artifacts. Some details (exact LOC, test counts, config field names) may change during Phase 3 code generation.
- No static site generator (Hugo, Jekyll, etc.) is configured yet. The pages are raw Markdown; a future phase could wire them into a site builder.
- Cross-links between pages use relative paths; these work for file-based reading but may need adjustment for a site generator.

## Next phase

- Phase 4: wire `docs/site/` into a static site generator, or continue with Phase 3 code generation milestones (M0-M9 per ROADMAP.md).
