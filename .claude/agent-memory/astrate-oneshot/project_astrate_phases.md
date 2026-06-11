---
name: project-astrate-phases
description: Astrate project goal, hard constraints, and the strict 3-phase workflow (Phases 1-2 done; Phase 3 codegen needs explicit user approval)
metadata:
  type: project
---

Astrate is a lean Go re-implementation (spiritual fork) of the Astarte IoT platform, developed in `/Users/atsetilam/astrate`.

**Why:** Strip Astarte's operational complexity (Kubernetes, Cassandra/ScyllaDB RAM footprint, fragmented Elixir microservices) so it runs on a 1-2 GB VPS or at the edge, while staying wire-compatible with official Astarte device SDKs.

**Hard constraints (from the user's master prompt — do not relax):**
- Go only, single statically linked binary, modular monolith (no microservices)
- PostgreSQL + TimescaleDB only (no Cassandra); embedded MQTT broker (mochi-mqtt) preferred
- docker-compose or bare binary deployment, zero Kubernetes
- SDK compatibility is CRITICAL: exact MQTT topics (`<realm>/<device_id>/...`), Astarte MQTT v1 protocol, pairing REST API, mTLS cert CN `<realm>/<device_id>`, JWT claim model
- Dual payloads: BSON (official SDKs) AND plain JSON (for AtomVM / constrained devices)

**Workflow state:** Strict phase gating — never skip ahead without explicit user approval.
- Phase 1 (architecture design doc): DONE 2026-06-10 → `docs/DESIGN.md`. APPROVED as-is by user, no amendments.
- Phase 2 (roadmap): DONE 2026-06-10 → `docs/ROADMAP.md`. APPROVED by user 2026-06-10. 10 milestones M0-M9, 12 codegen sessions (S1-S12), 4 SDK conformance checkpoints CP-A (after M4 pairing), CP-B (after M6 engine), CP-C (after M7 APIs), CP-D (M9 full matrix). Test tiers T1 unit / T2 testcontainers-TimescaleDB / T3 E2E / T4 official-SDK conformance / T5 load+security.
- Phase 3 (code generation): IN PROGRESS, one session per milestone, each ends with its gate green + a commit on `main`. ROADMAP.md §0.1 fixes binding generation rules (strict dependency order, tests ride with subject, no TODOs except 3 named extension points) and §11 maps milestones to sessions — follow it literally.
  - S1/M0 DONE 2026-06-11, commit 4bbf698: scaffolding, dep pinning, testutil harness, CI skeleton. Gate green locally (lint 0 issues, T1 green, T2 container test boots latest-pg16, timescaledb ext 2.27.2).
  - S2/M1.1-M1.3 DONE 2026-06-11 (committed on main right after 4bbf698): `pkg/deviceid` (astartectl-parity vectors, FuzzParse), `pkg/astarteapi` (frozen golden envelopes incl. Phoenix-default detail strings "Not Found"/"Bad Request" + Astarte-specific "Device not found"; DecodeData unwrapper), `pkg/interfaceschema` (strict parser, zero-alloc backtracking EndpointTrie, Compile→§2.6 structs, CheckMinorUpgrade). Gate green: lint 0 issues, T1+race, both fuzz targets 10s, TestMatchZeroAllocs==0.
  - Next: S3 = M1.4 `pkg/payload` (ROADMAP §2.4): BSON/JSON dual codec, golden BSON vectors captured from the official Go SDK encoder, sniffing per DESIGN §3.5.2.
  - S2 implementation notes S3+ should know: upstream standard-interfaces live in astarte-platform/astarte repo under `standard-interfaces/` (NOT a separate repo); vendored at commit d6a1f5e578 into pkg/interfaceschema/testdata/valid/ with provenance in testdata/README.md. CompiledMapping lives in trie.go (not compile.go) so trie tests compile standalone. EndpointIDResolver is a 2-method interface (ResolveInterface/ResolveEndpoint), nil → zero IDs. gosec #nosec needs to sit on the flagged line itself, not above the import block.

**Phase 3 implementation facts S2+ must know (not derivable at a glance):**
- Module path `github.com/astrate-platform/astrate` (taken from the git remote).
- Do NOT run `go mod tidy` until all pinned deps are actually imported — go.mod hand-pins all §1.1 deps in a direct require block before code imports them; tidy would silently drop the unused ones. If a build hits "missing go.sum entry", `go get <pkg>@pinned-version` for that subpackage instead (this is how jackc/puddle was added for pgxpool).
- golangci-lint pinned at v2.12.2 (v2 config schema), installed via `make tools`; lint runs with `integration`+`e2e` build tags so tag-gated suites are linted. revive's exported/package-comment rules are active — every package needs a package comment and exported symbols need `Name ...`-style doc comments.
- testutil conventions: `testutil.StartTimescale(t)` (T2, honors ASTRATE_TEST_DSN reuse), `testutil.Golden(t, name, got)` with the `-update` flag registered in package testutil — other packages must not re-register an `-update` flag.
- Build tags: T2 = `//go:build integration`, T3 = `e2e` (Makefile runs e2e with both tags).

**Dependency pins decided in ROADMAP.md M0 (frozen unless user amends):** stdlib net/http (Go >= 1.22 ServeMux), mochi-mqtt/server/v2, pgx/v5, golang-migrate/v4 + go:embed, mongo-driver/v2 bson (raw API), golang-jwt/v5, x/crypto/bcrypt, bbolt (chosen over pebble), hashicorp/golang-lru/v2, prometheus client_golang, coder/websocket, BurntSushi/toml, testcontainers-go with timescale/timescaledb:latest-pg16, paho.mqtt.golang (test-only).

**How to apply:** Before starting Phase 3, confirm the user approved Phase 2 (docs/ROADMAP.md). All design decisions (shared hypertable with typed columns, in-process sharded pipeline replacing RabbitMQ, embedded CA replacing CFSSL, payload format sniffing) are frozen in docs/DESIGN.md; the file-by-file order, tests, and dependency pins are frozen in docs/ROADMAP.md — read both before generating code, treat them as the spec unless the user amends them. Schema in DESIGN.md §2.2-2.4 must be transcribed verbatim into migrations; deviations are design-change requests, not ad-hoc fixes.
