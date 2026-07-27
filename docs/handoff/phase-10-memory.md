# Phase 10 Memory — Fix Observability Metrics + Add Cross-References (2026-07-27)

## Completed

- Fixed observability.md: corrected all metric names to match source code, removed fabricated `[metrics]` config block, added missing metrics.
- Added "See also" cross-reference sections to operations.md, deployment.md, and troubleshooting.md.

## Files Changed

- `docs/site/observability.md` — rewritten: correct metric names from source, removed non-existent `[metrics]` config section, added all 15 metrics with accurate types/labels, added "See also" section
- `docs/site/operations.md` — added "See also" section linking to deployment, configuration-reference, observability, pairing-and-security, troubleshooting
- `docs/site/deployment.md` — expanded "See also" section to include operations, compatibility, observability
- `docs/site/troubleshooting.md` — expanded "See also" section to include configuration-reference, operations, observability, pairing-and-security

## Evidence Read

- `internal/engine/router.go:385-486` — engine metric names (`astrate_engine_messages_total`, `astrate_engine_rejects_total`, `astrate_engine_shard_depth`, etc.)
- `internal/engine/triggers/actions.go:137-249` — trigger metric names (`astrate_engine_trigger_deliveries_total`, `astrate_engine_trigger_retries_total`)
- `internal/engine/stream/bus.go:107-119` — stream metric (`astrate_engine_stream_dropped_total`)
- `internal/observability/metrics.go:56-79` — broker and DB pool metrics (`astrate_broker_sessions`, `astrate_db_pool_*_conns`)
- `internal/config/config.example.toml` — confirmed no `[metrics]` section exists

## Key Findings

- `observability.md` had 10+ incorrect metric names (e.g. `astrate_ingest_total` instead of `astrate_engine_messages_total`)
- `observability.md` listed `astrate_broker_connect_total` and `astrate_broker_disconnect_total` which don't exist in source
- `observability.md` listed `astrate_batch_rows` which doesn't exist
- `observability.md` had a fabricated `[metrics] addr` config block not present in `config.example.toml`
- DB pool metric names were incomplete (missing `_conns` suffix and 2 of 4 gauges)

## Verification

- `cd docs && make sync && make build` — exit 0, no project-specific warnings
- All 22 pages present in `site-dist/`

## Known Risks

- The `ROADMAP.md` link in operations.md is a relative link to a file that exists in `docs/site/` but is not in the MkDocs nav — this is pre-existing and acceptable since it links to the raw markdown.
- Pre-existing warning: `index.md` contains an unrecognized relative link `api/` — not introduced by this phase.

## Next Steps

1. Commit to `main`, push, verify on live site
2. Phase 11 candidates: review remaining narrative pages for accuracy (data-modeling, mqtt-protocol, interface-schema, etc.) or create any missing ROADMAP-referenced pages
