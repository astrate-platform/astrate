# Architecture

Astrate collapses Astarte's ~8 Elixir/OTP microservices into Go packages inside one process. The decoupling is preserved at the package boundary, not at the network boundary.

## Process model

```
                       ┌────────────────────────────────────────────────────────────┐
                       │                     astrate (single binary)                │
                       │                                                            │
  Device SDKs ──mTLS──▶│  internal/broker          internal/engine                  │
  (MQTT 3.1.1, :8883)  │  ┌─────────────────┐      ┌──────────────────────────────┐ │
                       │  │ mochi-mqtt v2   │      │  shard router (FNV(device))  │ │
                       │  │  - TLS+mTLS     │─chan▶│  ┌────────┐  ┌────────┐      │ │
                       │  │  - AuthHook     │      │  │shard 0 │...│shard N │      │ │
                       │  │  - ACLHook      │◀chan─│  └───┬────┘  └───┬────┘      │ │
                       │  │  - inline client│      │      ▼           ▼           │ │
                       │  └─────────────────┘      │  validate → persist → trig   │ │
                       │                           └───────┬──────────┬───────────┘ │
  Devices ──HTTPS────▶ │  internal/pairing                 │          │             │
  (register/CSR :8080) │  internal/appengine  ◀────────────┘          ▼             │
  Operators ─HTTPS───▶ │  internal/realm                       live fan-out bus     │
  (JWT auth)           │  internal/housekeeping                (WebSocket/SSE)      │
                       │                                                            │
                       │  internal/store  (pgxpool) ── one connection pool ─────────┼──▶ PostgreSQL 16
                       └────────────────────────────────────────────────────────────┘    + TimescaleDB
```

## Service mapping

| Astarte component | Astrate package | Responsibility |
|---|---|---|
| Pairing API + CFSSL | `internal/pairing` | Device registration, credentials, CSR signing, broker info |
| VerneMQ + plugin | `internal/broker` | Embedded MQTT broker, mTLS, ACLs, sessions |
| Data Updater Plant | `internal/engine` | Sharded pipeline: validation, persistence, triggers, fan-out |
| AppEngine API | `internal/appengine` | REST API for device data, server-owned publishes, groups |
| Realm Management API | `internal/realm` | Interface and trigger CRUD per realm |
| Housekeeping API | `internal/housekeeping` | Realm lifecycle (create/delete) |
| Trigger Engine | `internal/engine/triggers` | HTTP webhook actions with retry/backoff |
| Astarte Channels | `internal/appengine/stream` | Simplified WebSocket/SSE live stream |

## Package layout

```
astrate/
├── cmd/astrate/              # main: config, wiring, lifecycle, graceful shutdown
├── internal/
│   ├── broker/               # embedded MQTT broker, auth/ACL hooks, session store
│   ├── engine/               # ingestion pipeline: shards, validation, persistence
│   │   ├── triggers/         # trigger matching + action execution (HTTP webhooks)
│   │   └── stream/           # in-process pub/sub bus for live consumers
│   ├── pairing/              # registration, credentials, CSR signing
│   │   └── ca/               # embedded per-realm certificate authority
│   ├── appengine/            # REST API + WebSocket/SSE live stream
│   ├── realm/                # interfaces, triggers CRUD
│   ├── housekeeping/         # realm lifecycle
│   ├── auth/                 # JWT validation, Astarte authz claims
│   ├── store/                # pgx repositories, migrations
│   └── config/               # TOML/env config
├── pkg/
│   ├── interfaceschema/      # Interface JSON parsing, endpoint trie compiler
│   ├── payload/              # BSON/JSON dual codec
│   ├── deviceid/             # 128-bit device ID encoding
│   └── astarteapi/           # shared API envelope types
├── migrations/               # SQL migration files (go:embed)
└── docker-compose.yml
```

**Dependency rule:** `pkg/*` has no `internal/*` imports. `internal/store` is imported by domain packages but never imports them. HTTP API packages depend on domain services through interfaces (hexagonal-lite).

## Concurrency model

Astarte's DUP relies on RabbitMQ queue-per-shard semantics for per-device message ordering. Astrate reproduces this with an in-process shard router:

1. The broker hook delivers every inbound PUBLISH as an `InboundMessage` onto the engine's intake.
2. The router computes `shard = FNV1a(device_id) % N` (default 16) and appends to that shard's bounded channel.
3. One goroutine per shard processes messages strictly in order -- same guarantee as DUP, zero broker dependency.
4. Backpressure: bounded channels (default 4096/shard). When saturated, the broker hook blocks QoS >= 1 PUBACKs (pushing backpressure to the device). QoS 0 messages are dropped with a metric increment.
5. Persistence uses per-shard micro-batching (flush at 64 rows or 50 ms) through `pgx.Batch`/`COPY`.

## Multi-tenancy (realms)

Realms survive from Astarte because they are part of every wire contract:

- MQTT topics: `<realm>/<device_id>/...`
- REST URLs: `/v1/<realm>/...`
- Certificate CNs: `<realm>/<device_id>`

Realms become cheap: a `realms` row + per-realm CA + per-realm JWT public keys. A single-realm install is just a realm named e.g. `home` created at first boot.

## Lifecycle and resilience

- **Graceful shutdown:** broker stops accepting, shards drain (bounded by timeout), batches flush.
- **Crash safety:** QoS >= 1 messages are PUBACK'd only after the persistence batch commits (at-least-once). Datastream inserts are idempotent; properties are upserts.
- **DB outage:** shards park with exponential backoff; broker applies backpressure; QoS 0 data degrades first.
