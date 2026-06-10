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
- Phase 2 (roadmap): DONE 2026-06-10 → `docs/ROADMAP.md`. Awaiting user approval. 10 milestones M0-M9, 12 codegen sessions (S1-S12), 4 SDK conformance checkpoints CP-A (after M4 pairing), CP-B (after M6 engine), CP-C (after M7 APIs), CP-D (M9 full matrix). Test tiers T1 unit / T2 testcontainers-TimescaleDB / T3 E2E / T4 official-SDK conformance / T5 load+security.
- Phase 3: full source code generation, minimal TODOs/placeholders. ROADMAP.md §0.1 fixes binding generation rules (strict dependency order, tests ride with subject, no TODOs except 3 named extension points) and §11 maps milestones to sessions — follow it literally.

**Dependency pins decided in ROADMAP.md M0 (frozen unless user amends):** stdlib net/http (Go >= 1.22 ServeMux), mochi-mqtt/server/v2, pgx/v5, golang-migrate/v4 + go:embed, mongo-driver/v2 bson (raw API), golang-jwt/v5, x/crypto/bcrypt, bbolt (chosen over pebble), hashicorp/golang-lru/v2, prometheus client_golang, coder/websocket, BurntSushi/toml, testcontainers-go with timescale/timescaledb:latest-pg16, paho.mqtt.golang (test-only).

**How to apply:** Before starting Phase 3, confirm the user approved Phase 2 (docs/ROADMAP.md). All design decisions (shared hypertable with typed columns, in-process sharded pipeline replacing RabbitMQ, embedded CA replacing CFSSL, payload format sniffing) are frozen in docs/DESIGN.md; the file-by-file order, tests, and dependency pins are frozen in docs/ROADMAP.md — read both before generating code, treat them as the spec unless the user amends them. Schema in DESIGN.md §2.2-2.4 must be transcribed verbatim into migrations; deviations are design-change requests, not ad-hoc fixes.
