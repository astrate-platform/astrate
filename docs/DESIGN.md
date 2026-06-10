# Astrate — Architectural Design Document (Phase 1)

**Project:** Astrate — a lean, single-binary, Astarte-wire-compatible IoT platform in Go
**Status:** Phase 1 deliverable — awaiting approval before Phase 2 (roadmap) and Phase 3 (code)
**Version:** 0.1 (2026-06-10)

---

## 0. Executive Summary

Astrate is a spiritual fork of the [Astarte IoT Platform](https://github.com/astarte-platform/astarte)
that preserves Astarte's *external contracts* — the Astarte MQTT v1 protocol, the Pairing API, the
AppEngine/Realm Management REST APIs, and the Interface JSON specification — while replacing its
*internals* entirely:

| Concern | Astarte (upstream) | Astrate |
|---|---|---|
| Runtime | ~8 Elixir/OTP microservices | One statically linked Go binary (modular monolith) |
| Orchestration | Kubernetes + astarte-operator | `docker-compose.yml` or bare binary + systemd |
| Time-series store | Cassandra / ScyllaDB (multi-GB RAM floor) | PostgreSQL 16 + TimescaleDB |
| Metadata store | Cassandra | Same PostgreSQL instance |
| Message bus | RabbitMQ (inter-service AMQP) | In-process Go channels (sharded worker pools) |
| MQTT broker | VerneMQ + `astarte_vmq_plugin` | Embedded `mochi-mqtt/server` v2 (in-binary) |
| Certificate authority | CFSSL sidecar | Embedded CA via Go `crypto/x509` |
| Payloads | BSON only | BSON **and** plain JSON (AtomVM / bare-metal friendly) |
| Target footprint | 4–16 GB RAM cluster | ≤ 1–2 GB RAM single VPS / edge node |

**Hard compatibility goal:** an unmodified official Astarte device SDK (C/ESP32, Python, Go, Rust,
Java/Android, Elixir) pointed at Astrate's Pairing API URL must register, obtain credentials,
connect over mutual TLS, exchange introspection/properties/datastreams, and receive server-owned
data — without a single line of SDK code changed.

**Non-goals (v1):** Astarte Flow, the Kubernetes operator, Cassandra migration tooling, the
Dashboard UI (the API is compatible, so the upstream dashboard *may* work later, but it is not a
v1 acceptance criterion), and Astarte Channels' full Phoenix-socket protocol (we provide a simpler
WebSocket/SSE live stream; see §1.4).

---

## 1. Service Mapping: Astarte Components → Astrate Go Modules

### 1.1 Upstream component inventory and their Astrate destination

Astarte decomposes into independently deployed Elixir apps communicating over RabbitMQ and
querying Cassandra. Astrate collapses them into Go packages inside one process, communicating
through typed in-process interfaces and channels. The decoupling is preserved *at the package
boundary* (each domain exposes a narrow Go interface, owns its tables, and is independently
testable), not at the network boundary.

| Astarte component | Responsibility (upstream) | Astrate package | Notes on the mapping |
|---|---|---|---|
| **Pairing API** + CFSSL | Device registration, credentials secret issuance, X.509 CSR signing, broker info | `internal/pairing` | CA is embedded (`crypto/x509`); per-realm CA key in DB (encrypted) or on-disk PEM. Same REST surface (§4). |
| **VerneMQ** + `astarte_vmq_plugin` | MQTT broker, mTLS termination, ACLs, bridging publishes onto AMQP | `internal/broker` | Embedded `mochi-mqtt/server` v2 with auth/ACL hooks; publishes flow into the engine through a Go channel, not AMQP. |
| **Data Updater Plant (DUP)** | Consumes device messages from AMQP, validates against interfaces, writes to Cassandra, detects trigger conditions | `internal/engine` | The heart of Astrate: sharded per-device ordered pipeline → validation → persistence → trigger evaluation → live fan-out. |
| **AppEngine API** | REST API to read device data, publish server-owned data, manage groups/aliases; Channels (WebSockets) | `internal/appengine` | Same `/v1/<realm>/devices/...` resource model; server-owned publishes are handed to `internal/broker` via the engine. |
| **Realm Management API** | Interface (schema) and trigger CRUD per realm | `internal/realm` | Interfaces stored as JSONB + compiled in-memory (§2.6). Trigger CRUD included; trigger *execution* lives in `internal/engine/triggers`. |
| **Housekeeping API** | Realm lifecycle (create/delete realms, instance admin) | `internal/housekeeping` | Realms become rows + `search_path`-free schema-qualified tables (§2.1). Creating a realm is a transaction, not a Cassandra keyspace provision. |
| **Trigger Engine** | Executes trigger actions (HTTP webhooks, AMQP messages) | `internal/engine/triggers` | HTTP webhook actions with retry/backoff; AMQP action replaced by optional NATS/HTTP forwarding (extension point). |
| **Astarte Channels** | Phoenix WebSocket rooms for live data | `internal/appengine/stream` | Simplified WebSocket + SSE endpoint fed by the engine's fan-out bus. Not wire-compatible with Phoenix sockets in v1 (documented deviation). |
| **Astarte Flow** | Dataflow processing framework | — | Out of scope. Triggers + the live stream socket cover the common cases. |
| **astartectl / Dashboard** | Tooling/UI | — | API-compatible by construction; `astartectl` should largely work against Astrate. Verified in Phase 2 test plan. |

### 1.2 Process architecture

```
                       ┌────────────────────────────────────────────────────────────┐
                       │                     astrate (single binary)                │
                       │                                                            │
  Device SDKs ──mTLS──▶│  internal/broker          internal/engine                  │
  (MQTT 3.1.1, :8883)  │  ┌─────────────────┐      ┌──────────────────────────────┐ │
                       │  │ mochi-mqtt v2   │      │  shard router (FNV(device))  │ │
                       │  │  - TLS+mTLS     │─chan▶│  ┌────────┐  ┌────────┐      │ │
                       │  │  - AuthHook     │      │  │shard 0 │… │shard N │      │ │
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

### 1.3 Package layout (top-level skeleton)

```
astrate/
├── cmd/astrate/              # main: config load, wiring, lifecycle, graceful shutdown
├── internal/
│   ├── broker/               # embedded MQTT broker, auth/ACL hooks, session store
│   ├── engine/               # ingestion pipeline: shards, validation, persistence
│   │   ├── triggers/         # trigger matching + action execution (HTTP webhooks)
│   │   └── stream/           # in-process pub/sub bus for live consumers
│   ├── pairing/              # registration, credentials, CSR signing, /pairing/v1 API
│   │   └── ca/               # embedded per-realm certificate authority
│   ├── appengine/            # /appengine/v1 REST API + WebSocket/SSE live stream
│   ├── realm/                # /realmmanagement/v1 API: interfaces, triggers CRUD
│   ├── housekeeping/         # /housekeeping/v1 API: realm lifecycle
│   ├── auth/                 # JWT validation, Astarte authz claim regex matching
│   ├── store/                # pgx repositories, migrations (golang-migrate embedded)
│   └── config/               # TOML/env config
├── pkg/
│   ├── interfaceschema/      # Interface JSON parsing, validation, endpoint trie compiler
│   ├── payload/              # BSON/JSON dual codec for Astarte data payloads (§3.5)
│   ├── deviceid/             # 128-bit device ID <-> base64url (22 chars) <-> UUID
│   └── astarteapi/           # shared API envelope types ({"data": ...}, error format)
├── migrations/               # SQL migration files (embedded via go:embed)
├── docs/
└── docker-compose.yml        # astrate + timescale/timescaledb-ha image
```

Dependency rule: `pkg/*` has no `internal/*` imports; `internal/store` is imported by domain
packages but never imports them; `internal/engine` is the only writer of telemetry; HTTP API
packages depend on domain services through interfaces defined where they are *consumed*
(hexagonal-lite), keeping each module mockable.

### 1.4 Concurrency model (replacing RabbitMQ)

Astarte's DUP relies on RabbitMQ queue-per-shard semantics to guarantee **per-device message
ordering**. Astrate reproduces this with an in-process shard router:

- The broker hook delivers every inbound device PUBLISH as an `InboundMessage` (topic, payload,
  QoS, reception timestamp) onto the engine's intake.
- The router computes `shard = FNV1a(device_id) % N` (N configurable, default 16) and appends to
  that shard's bounded channel. One goroutine per shard processes messages strictly in order —
  same guarantee as DUP, zero broker dependency.
- Backpressure: bounded channels (default 4096/shard). When a shard saturates, the broker hook
  blocks that client's inflight acknowledgment (QoS ≥ 1 PUBACK is deferred), pushing backpressure
  to the device exactly as a slow AMQP consumer would in upstream Astarte. QoS 0 messages are
  dropped with a metric increment when the shard is full.
- Persistence uses per-shard micro-batching (flush at 64 rows or 50 ms, whichever first) through
  `pgx.Batch`/`COPY`, which is where TimescaleDB ingestion throughput comes from.

Chosen explicitly over an external NATS/Mosquitto: one less moving part, and the ordering +
backpressure semantics are easier to make airtight in-process. The `engine` intake is defined as
a Go interface, so an external bus can be reintroduced later without touching the broker or
persistence layers.

### 1.5 Multi-tenancy (realms)

Astarte realms survive in Astrate, because they are part of every wire contract (topics are
`<realm>/<device_id>/...`, URLs are `/v1/<realm>/...`, certificate CNs are
`<realm>/<device_id>`). They become cheap: a `realms` row + per-realm CA + per-realm JWT public
keys. A single-realm install is just a realm named e.g. `home` created at first boot via
`housekeeping` (or auto-provisioned from config for the zero-ceremony path).

---

## 2. Data Modeling: PostgreSQL + TimescaleDB

### 2.1 Tenancy layout decision

Two candidate layouts were considered:

1. **Schema-per-realm** (Postgres schemas as Cassandra-keyspace analogue) — clean isolation, but
   dynamic DDL on realm creation, painful migrations across N schemas, and TimescaleDB jobs
   multiply per realm.
2. **Shared tables with a `realm_id` column** — one migration path, one set of hypertables and
   compression/retention policies, trivially indexable.

**Decision: shared tables + `realm_id`.** Astrate targets small installs (1–5 realms typical);
row-level tenancy with composite keys is simpler and faster at this scale. Realm deletion is a
transactional cascade.

### 2.2 Relational metadata schema

```sql
-- Realms (Housekeeping domain)
CREATE TABLE realms (
    id               smallint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name             text NOT NULL UNIQUE CHECK (name ~ '^[a-z][a-z0-9]*$'),
    jwt_public_keys  jsonb NOT NULL DEFAULT '[]',   -- array of PEM strings (RSA/EC)
    ca_certificate   text NOT NULL,                  -- realm CA cert, PEM
    ca_private_key   bytea NOT NULL,                 -- encrypted at rest (AES-256-GCM,
                                                     -- key from config/KMS env var)
    device_registration_limit integer,
    created_at       timestamptz NOT NULL DEFAULT now()
);

-- Interfaces (Realm Management domain). The raw JSON is the source of truth;
-- generated columns lift the routing-critical fields out for indexing.
CREATE TABLE interfaces (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    realm_id      smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    definition    jsonb NOT NULL,
    name          text     GENERATED ALWAYS AS (definition->>'interface_name') STORED,
    major_version integer  GENERATED ALWAYS AS ((definition->>'version_major')::int) STORED,
    minor_version integer  GENERATED ALWAYS AS ((definition->>'version_minor')::int) STORED,
    type          text     GENERATED ALWAYS AS (definition->>'type') STORED,          -- datastream|properties
    ownership     text     GENERATED ALWAYS AS (definition->>'ownership') STORED,     -- device|server
    aggregation   text     GENERATED ALWAYS AS (coalesce(definition->>'aggregation','individual')) STORED,
    created_at    timestamptz NOT NULL DEFAULT now(),
    UNIQUE (realm_id, name, major_version)
);

-- Mappings, normalized for endpoint-id stability (mirrors Astarte's endpoint UUIDs).
CREATE TABLE endpoints (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    interface_id  bigint NOT NULL REFERENCES interfaces(id) ON DELETE CASCADE,
    endpoint      text NOT NULL,            -- e.g. '/%{sensor_id}/value'
    value_type    text NOT NULL,            -- double|integer|boolean|longinteger|string|
                                            -- binaryblob|datetime|<type>array
    reliability   text NOT NULL DEFAULT 'unreliable',  -- → QoS 0|1|2
    retention     text NOT NULL DEFAULT 'discard',
    expiry        integer NOT NULL DEFAULT 0,
    database_retention_policy text NOT NULL DEFAULT 'no_ttl',
    database_retention_ttl    integer,
    explicit_timestamp boolean NOT NULL DEFAULT false,
    allow_unset   boolean NOT NULL DEFAULT false,
    UNIQUE (interface_id, endpoint)
);

-- Devices (Pairing + AppEngine domains)
CREATE TABLE devices (
    id                  uuid NOT NULL,          -- the 128-bit Astarte device ID
    realm_id            smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    credentials_secret_hash text NOT NULL,      -- bcrypt
    status              text NOT NULL DEFAULT 'registered',  -- registered|confirmed|inhibited
    introspection       jsonb NOT NULL DEFAULT '{}',  -- {"iface.Name": {"major":1,"minor":2}, ...}
    old_introspection   jsonb NOT NULL DEFAULT '{}',
    aliases             jsonb NOT NULL DEFAULT '{}',
    attributes          jsonb NOT NULL DEFAULT '{}',
    cert_serial         text,                   -- serial of last issued client cert
    cert_aki            text,                   -- authority key identifier
    first_registration  timestamptz NOT NULL DEFAULT now(),
    first_credentials_request timestamptz,
    last_credentials_request_ip inet,
    last_connection     timestamptz,
    last_disconnection  timestamptz,
    last_seen_ip        inet,
    connected           boolean NOT NULL DEFAULT false,
    total_received_msgs  bigint NOT NULL DEFAULT 0,
    total_received_bytes bigint NOT NULL DEFAULT 0,
    payload_format_hint  text NOT NULL DEFAULT 'bson',  -- bson|json, see §3.5.4
    PRIMARY KEY (realm_id, id)
);
CREATE INDEX devices_aliases_gin ON devices USING gin (aliases);

-- Device groups (AppEngine)
CREATE TABLE groups (
    id        bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    realm_id  smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    name      text NOT NULL,
    UNIQUE (realm_id, name)
);
CREATE TABLE group_devices (
    group_id  bigint NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    realm_id  smallint NOT NULL,
    device_id uuid NOT NULL,
    PRIMARY KEY (group_id, device_id),
    FOREIGN KEY (realm_id, device_id) REFERENCES devices(realm_id, id) ON DELETE CASCADE
);

-- Triggers (Realm Management domain; executed by the engine)
CREATE TABLE triggers (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    name       text NOT NULL,
    definition jsonb NOT NULL,        -- Astarte trigger JSON (simple_triggers + action)
    UNIQUE (realm_id, name)
);
```

### 2.3 Properties storage

Properties are last-value-wins key/value state — a perfect fit for a plain relational table with
an upsert. Values are stored twice-typed: a `jsonb` rendering for cheap API reads, plus the
endpoint's declared type retained on the row so the API layer can re-encode precisely
(longinteger as number-as-string where required, binaryblob as base64, datetime as RFC3339).

```sql
CREATE TABLE properties (
    realm_id     smallint NOT NULL,
    device_id    uuid NOT NULL,
    interface_id bigint NOT NULL REFERENCES interfaces(id) ON DELETE CASCADE,
    endpoint_id  bigint NOT NULL REFERENCES endpoints(id),
    path         text NOT NULL,               -- concrete path, e.g. '/lcdCmd' or '/4/enable'
    value        jsonb NOT NULL,
    value_type   text NOT NULL,
    set_at       timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (realm_id, device_id, interface_id, path),
    FOREIGN KEY (realm_id, device_id) REFERENCES devices(realm_id, id) ON DELETE CASCADE
);
```

Unset (empty MQTT payload on a property path, `allow_unset: true`) ⇒ `DELETE` of the row.
Server-owned property writes via AppEngine go through the same table *and* are published to the
device (retained delivery semantics via the broker session queue, plus the
`/control/consumer/properties` purge mechanism on session resume — §3.4).

### 2.4 Datastream storage (TimescaleDB hypertables)

Three candidate designs were evaluated:

| Option | Pros | Cons |
|---|---|---|
| (a) Table-per-interface, dynamically created | Perfect typing, per-interface `drop_chunks` retention, best compression | Runtime DDL, migration complexity, thousands of tables possible, planner bloat |
| (b) One JSONB-value hypertable | Simplest | Loses numeric compression + aggregation pushdown; `downsample_to` queries become slow |
| (c) One wide hypertable with sparse typed columns (mirrors Astarte's Cassandra `individual_datastreams`) | No runtime DDL; typed columns compress well (Timescale columnar compression handles NULL-sparse columns cheaply); a single set of policies | Per-interface TTL needs background `DELETE`s instead of `drop_chunks` |

**Decision: (c) for individual datastreams, plus a JSONB hypertable for object aggregations** —
the same shape Astarte itself uses on Cassandra, which de-risks semantics.

```sql
CREATE TABLE individual_datastreams (
    realm_id      smallint NOT NULL,
    device_id     uuid NOT NULL,
    interface_id  bigint NOT NULL,
    endpoint_id   bigint NOT NULL,
    path          text NOT NULL,
    ts            timestamptz NOT NULL,     -- value timestamp (explicit_timestamp or reception)
    reception_ts  timestamptz NOT NULL,
    -- exactly one of the following is non-NULL, per the endpoint's declared type:
    value_double       double precision,
    value_integer      integer,
    value_longinteger  bigint,
    value_boolean      boolean,
    value_string       text,
    value_binaryblob   bytea,
    value_datetime     timestamptz,
    value_array        jsonb               -- all *array types (doublearray, stringarray, ...)
);
SELECT create_hypertable('individual_datastreams', by_range('ts', INTERVAL '7 days'));
CREATE INDEX ids_series_idx ON individual_datastreams
    (realm_id, device_id, interface_id, path, ts DESC);

ALTER TABLE individual_datastreams SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'realm_id, device_id, interface_id, path',
    timescaledb.compress_orderby   = 'ts DESC'
);
SELECT add_compression_policy('individual_datastreams', INTERVAL '7 days');

CREATE TABLE object_datastreams (
    realm_id      smallint NOT NULL,
    device_id     uuid NOT NULL,
    interface_id  bigint NOT NULL,
    path          text NOT NULL,             -- the parametric prefix, e.g. '/12'
    ts            timestamptz NOT NULL,
    reception_ts  timestamptz NOT NULL,
    value         jsonb NOT NULL             -- {"temp": 22.1, "hum": 41.0}
);
SELECT create_hypertable('object_datastreams', by_range('ts', INTERVAL '7 days'));
CREATE INDEX ods_series_idx ON object_datastreams
    (realm_id, device_id, interface_id, path, ts DESC);
-- same compression policy shape as above
```

Sizing rationale for the 1–2 GB VPS target: 7-day chunks with segment-by compression typically
yield 90–95 % reduction on telemetry; the active chunk's indexes for a few hundred devices at
1 msg/s stay well under 200 MB. Chunk interval is configurable for heavier installs.

### 2.5 Retention, expiry, and downsampling

- **`database_retention_policy: use_ttl`** (per endpoint): a single Timescale **user-defined
  action** (background job, e.g. hourly) runs
  `DELETE FROM individual_datastreams WHERE (interface_id, endpoint_id) IN (ttl'd set) AND ts < now() - ttl`,
  batched by chunk to bound transaction size. A global hard-cap retention
  (`retention.max_days`, optional) additionally uses `add_retention_policy` → `drop_chunks` for
  cheap whole-chunk eviction.
- **Datastream `expiry`** (message validity for offline devices) is honoured on the
  *server→device* path by setting MQTT message expiry in the broker's offline queue.
- **AppEngine `downsample_to` queries** map onto Timescale's `time_bucket()` + (optionally)
  `lttb()` from the toolkit when available, with a plain `time_bucket` + `avg/first/last`
  fallback — decided at startup by probing for the toolkit extension.

### 2.6 Dynamic interface validation in Go

Interfaces are dynamic schemas uploaded at runtime; the hot path cannot parse JSONB per message.

**Compilation.** On load (startup, interface install/update, or cache miss), `pkg/interfaceschema`
compiles each interface definition into a `CompiledInterface`:

```go
type CompiledInterface struct {
    ID           int64
    Name         string
    Major, Minor int
    Type         InterfaceType   // Datastream | Properties
    Ownership    Ownership       // Device | Server
    Aggregation  Aggregation     // Individual | Object
    Trie         *EndpointTrie   // segment-wise matcher
    ObjectLeaves map[string]*CompiledMapping // for object aggregation: last-level name → mapping
}

type CompiledMapping struct {
    EndpointID        int64
    ValueType         ValueType   // drives both validation and BSON/JSON decoding
    Reliability       byte        // → MQTT QoS
    Retention         Retention
    Expiry            time.Duration
    ExplicitTimestamp bool
    AllowUnset        bool
    DBRetentionTTL    time.Duration // 0 = no_ttl
}
```

**Endpoint trie.** Astarte endpoints are `/`-separated with `%{param}` placeholder segments
(e.g. `/%{sensor_id}/value`). The trie matches a concrete inbound path (`/4/value`) segment by
segment: exact-match children first, then the (single) parametric child. This is O(depth), with
zero allocations on the hot path (segments are sub-slices of the topic string). Placeholder
values are *not* semantically interpreted (matching upstream behaviour) but are length- and
charset-checked (`no '/'`, `no '#'/'+'`, ≤ 256 bytes, non-empty).

**Validation pipeline** (in `internal/engine`, per message):

1. **Introspection gate** — the sending device's cached introspection must declare the interface
   `name:major`; otherwise reject (and raise the upstream-equivalent
   `device_error`/`unexpected_value` trigger event + metric).
2. **Trie match** of the path → `CompiledMapping`, or reject (`unexpected_path`).
3. **Ownership check** — device-published data only on `ownership: device` interfaces; AppEngine
   publishes only on `ownership: server`.
4. **Payload decode** via `pkg/payload` (§3.5) → `(value any, explicitTS *time.Time)`.
5. **Type check + coercion** against `ValueType` with Astarte's rules:
   - `double` ← BSON double; also int32/int64 (lossless widening). JSON: any number.
   - `integer` ← int32 (or int64/double that fits in int32 exactly).
   - `longinteger` ← int64/int32; JSON: number **or decimal string** (JS 2^53 safety).
   - `boolean`, `string` (must be valid UTF-8, ≤ 64 KiB), `binaryblob` (BSON binary / JSON
     base64 string), `datetime` (BSON UTC datetime / JSON RFC 3339 string or integer epoch ms).
   - Arrays: homogeneous element type, each element checked as above; ≤ 1024 elements.
   - `explicit_timestamp: true` ⇒ `t` field required (datastreams); forbidden otherwise—tolerated
     and ignored, matching upstream leniency.
   - Properties: empty payload ⇒ unset, allowed only if `allow_unset`.
6. **Aggregation shape** — object-aggregated interfaces must arrive as one document of
   last-level keys on the path *prefix*; every key must resolve in `ObjectLeaves`.

Failures are never silent: each rejection increments a per-reason Prometheus counter, is logged
with device/interface/path, and feeds `device_error` triggers (parity with Astarte's DUP).

**Cache & invalidation.** `internal/engine` holds the compiled cache:
`map[realmID]map[name]map[major]*CompiledInterface` behind an `atomic.Pointer` snapshot (copy-on-
write; readers never lock). `internal/realm` bumps a `schema_epoch` and issues a Postgres
`NOTIFY astrate_interfaces` on CRUD; the engine LISTENs and rebuilds the affected realm's
snapshot. NOTIFY is not needed for correctness in a single process (a direct in-process callback
also fires); it exists so an optional second hot-standby instance stays coherent.

**Versioning semantics (parity):** new majors coexist; minor bumps must be non-breaking
(additive mappings only — enforced on PUT exactly as Realm Management does); a device's
introspection pins which `name:major` its messages validate against, and `minor` advertised by
the device may be ≤ the installed minor.

---

## 3. MQTT Topic & Protocol Wire-Compatibility (Astarte MQTT v1)

Reference: the *Astarte MQTT v1 Protocol* specification in the upstream docs. Note that the
`devices/<device_id>/interfaces/<interface_name>/...` shape is Astarte's **REST resource path**
(AppEngine API), which Astrate also implements (§3.7); the **MQTT wire topics** are the
`<realm>/<device_id>/...` scheme below. Both contracts are reproduced exactly.

### 3.1 Broker and connection contract

- **Broker:** `mochi-mqtt/server` v2 embedded in-process. MQTT 3.1.1 (what all Astarte SDKs
  speak; mochi also accepts 5.0). QoS 0/1/2, retained messages, persistent sessions.
- **Listener:** TLS on `:8883` with `ClientAuth: tls.RequireAndVerifyClientCert`, client CAs =
  the per-realm CA pool (all realm CAs loaded; the matching realm is derived from the client
  cert chain). Optional plaintext `:1883` listener exists *only* behind an explicit
  `insecure_dev_mode` flag for local development.
- **Identity:** the client certificate **CN is `<realm>/<device_id>`** (exactly what the Pairing
  CA issued, §4). The auth hook parses CN, checks the cert chains to that realm's CA, checks the
  device row exists and is not `inhibited`, and rejects otherwise. The MQTT *client ID* must
  equal the CN (upstream behaviour); a mismatch is rejected.
- **Sessions:** devices connect with `clean_session=false`. The broker persists session state
  (subscriptions + QoS ≥ 1 offline queue) in a bbolt/pebble file via a mochi storage hook, so
  `session_present` survives Astrate restarts — this matters because SDKs use
  `session_present=0` as the signal to replay introspection, subscriptions, and the empty-cache
  handshake.
- **Connection lifecycle events** (connect/disconnect, with client IP) update
  `devices.connected/last_connection/last_disconnection/last_seen_ip` and fire
  `device_connected`/`device_disconnected` triggers — the work `astarte_vmq_plugin` does
  upstream.

### 3.2 ACL model (enforced in the broker hook)

For an authenticated identity `<realm>/<device_id>` with base topic `B = <realm>/<device_id>`:

| Action | Allowed topics |
|---|---|
| PUBLISH | `B` (introspection), `B/control/emptyCache`, `B/control/producer/properties`, `B/<interface_name><path>` for interfaces in its introspection with `ownership: device` |
| SUBSCRIBE | `B/control/consumer/properties`, `B/<interface_name>/#` for `ownership: server` interfaces in its introspection (the hook permits the `B/#` superset some SDKs request, then delivery is naturally scoped by what the engine publishes) |

Everything else is denied and logged. Server-side (engine/AppEngine) publishes use mochi's
inline client and bypass ACLs.

### 3.3 Device → Astrate message taxonomy

| Topic | Payload | QoS | Astrate handling |
|---|---|---|---|
| `<realm>/<device_id>` | Introspection string: `;`-separated `interface_name:major:minor` triples (UTF-8 plain text, e.g. `com.ex.Sensors:1:0;com.ex.Geo:0:1`) | 2 | Parse; diff against stored introspection; update `devices.introspection` (+`old_introspection` for removed pairs); fire `incoming_introspection` / interface added/removed triggers; recompute the device's server-owned subscription expectations. |
| `<realm>/<device_id>/control/emptyCache` | `1` | 2 | Device lost its local cache: Astrate re-sends the device's server-owned **properties** (each on its data topic, QoS 2) and then publishes the consumer-properties purge message (§3.4). |
| `<realm>/<device_id>/control/producer/properties` | 4-byte **big-endian** uint32 (uncompressed size) + **zlib-deflated** `;`-separated list of `interface_name/path` entries for every device-owned property currently set | 2 | Decompress, parse, and **purge**: delete from `properties` any device-owned row for this device not present in the list (the device is the source of truth for its own properties). |
| `<realm>/<device_id>/<interface_name><path>` (path always starts with `/`) | BSON **or JSON** document `{ "v": <value>, "t": <timestamp, optional> }`; empty payload = property unset | per-mapping reliability | Full pipeline §2.6: validate → persist (upsert property / insert datastream row) → triggers → live fan-out. |

Parsing note: the topic is split as `realm / device_id / rest`; `rest` is matched against the
device's introspected interface names by longest-prefix (interface names contain dots, never
`/`), and the remainder is the path. An interface match failure ⇒ rejection metric + optional
`device_error` trigger, never a crash.

### 3.4 Astrate → device messages

| Topic | Payload | When |
|---|---|---|
| `<realm>/<device_id>/<interface_name><path>` | Same `{v, t}` document; **format per device hint** (§3.5.4); QoS from the mapping's reliability | AppEngine publish on a `ownership: server` interface; re-send of server-owned properties after `emptyCache`. |
| `<realm>/<device_id>/control/consumer/properties` | Same 4-byte BE size + zlib format as producer/properties, listing all currently-set **server-owned** property paths | After `emptyCache`, after session-present=0 reconnect, and after server-owned property unset, so the device can purge stale local state. |

Server-owned **property** values are also retained in the broker per-topic (retain flag) so a
freshly subscribing device converges immediately; the purge message handles deletion races.
Server-owned **datastreams** with `retention: stored/volatile` for offline devices ride on the
broker's persistent-session offline queue with per-message expiry (§2.5).

### 3.5 Dual-payload codec: BSON + JSON (`pkg/payload`)

#### 3.5.1 Why

Astarte mandates BSON (`{v, t}` documents). On ultra-constrained targets — in particular
**AtomVM** (BEAM bytecode on ESP32/RP2040-class MCUs) — a full BSON encoder/decoder is dead
weight, while JSON encoding is a tiny pure-Erlang/Elixir library away. Astrate therefore accepts
**both** encodings *on the same topics with the same semantics*, so an AtomVM device behaves as a
first-class Astarte device minus the BSON dependency.

#### 3.5.2 Format detection (sniffing)

Detection is structural, cheap, and unambiguous:

```
if len(p) == 0                      → control semantics (property unset)
else if len(p) >= 5
     && int32LE(p[0:4]) == len(p)   → BSON  (length prefix self-describes)
     && p[len(p)-1] == 0x00
else if first non-WS byte == '{'    → JSON
else                                → reject (metrics + device_error trigger)
```

A JSON document would have to start with `{` and happen to have a first-4-bytes little-endian
value equal to its own length **and** end with a NUL byte to collide — impossible for valid JSON
text (no NUL allowed). The sniff is therefore safe without any per-device configuration.

#### 3.5.3 JSON payload profile

A strict, documented profile (so AtomVM client authors have a spec):

```json
{ "v": <value>, "t": "2026-06-10T12:34:56.789Z" }
```

- `t` optional; RFC 3339 string **or** integer milliseconds since epoch.
- Scalar mapping by *declared interface type* (the mapping's `ValueType` disambiguates JSON's
  single number type): `double`/`integer` ← JSON number; `longinteger` ← JSON number or decimal
  string (for > 2^53 values); `boolean` ← JSON bool; `string` ← JSON string; `binaryblob` ←
  base64 (standard alphabet, padded) JSON string; `datetime` ← RFC 3339 string or epoch-ms
  number; arrays ← JSON arrays; object aggregation ← `v` is a JSON object of last-level keys.
- A bare-JSON shorthand (`22.5` instead of `{"v":22.5}`) is **not** accepted — keeping the
  envelope mandatory preserves symmetry with BSON and keeps `t` unambiguous.
- Maximum accepted payload size (both formats): 64 KiB default, configurable.

#### 3.5.4 Outbound format selection

Server→device messages must also be decodable by the device. Astrate tracks
`devices.payload_format_hint`:

- Default `bson` (official SDK assumption).
- Flipped to `json` the first time a device publishes a JSON data payload (sticky; reset on
  `emptyCache` only if the next data payload is BSON), and settable explicitly via an Astrate
  extension field at registration (`POST /agent/devices` body `{"data": {"hw_id": ...,
  "initial_payload_format": "json"}}` — additive, ignored-by-upstream-shaped).
- Control payloads (`consumer/properties`) keep the zlib+size format for both kinds of device:
  zlib inflate is available on AtomVM via its standard library, and changing control framing
  would fork the protocol. This keeps the deviation surface limited to the data-document
  encoding only.

#### 3.5.5 BSON specifics

Implemented with `go.mongodb.org/mongo-driver/bson` raw-document API (`bson.Raw` lookups; no
reflection, no intermediate maps) to keep per-message allocations near zero. `t` is BSON UTC
datetime; `v` element type must agree with the mapping per §2.6 rules.

### 3.6 Compatibility deviations (explicit, documented)

1. **Astarte Channels** Phoenix-socket protocol is replaced by plain WebSocket/SSE (different
   contract, additive endpoint). Device SDKs are unaffected (Channels is a consumer-side API).
2. **JSON payloads and `initial_payload_format`** are Astrate extensions — pure supersets;
   upstream SDK behaviour is byte-identical.
3. **MQTT 5.0** clients are accepted (mochi default) though Astarte uses 3.1.1; this is a
   superset with no SDK impact.
4. Everything else — topics, BSON documents, zlib control payloads, introspection format,
   pairing REST bodies, certificate CN, JWT claim model — is wire-identical by design and
   guarded by conformance tests (Phase 2) that run the *official* Astarte Go and Python SDKs
   against Astrate in CI.

### 3.7 REST API surface (compatibility inventory)

Mounted on one HTTP listener (`:8080`) with upstream-compatible base paths (configurable, so
both Astarte-style per-service hostnames behind a reverse proxy and a single host work):

- `/pairing/v1/<realm>/...` — §4.
- `/appengine/v1/<realm>/devices`, `/devices/<id>`, `/devices/<id>/interfaces`,
  `/devices/<id>/interfaces/<iface>[/<path>]` (GET datastream queries with
  `since/since_after/to/limit/downsample_to`, PUT/POST server-owned publishes, DELETE property
  unset), `devices-by-alias/...`, `groups...`. Envelope: `{"data": ...}`; errors:
  `{"errors": {"detail": "..."}}` with upstream status codes.
- `/realmmanagement/v1/<realm>/interfaces...`, `/triggers...`.
- `/housekeeping/v1/realms...`.
- Astrate-native additions under `/astrate/v1/...` (health, metrics, live stream socket) to
  avoid colliding with upstream's namespace.

---

## 4. Security & Pairing

### 4.1 Trust model overview

Three credential planes, identical to Astarte's:

1. **Realm JWTs** — humans/services calling REST APIs. Asymmetric: each realm holds N public
   keys (PEM, RSA-2048+/ECDSA P-256+, alg allowlist `RS256/RS384/RS512/ES256/ES384/ES512`,
   `none` and HMAC rejected); operators keep private keys. Housekeeping has its own key pair
   (instance-level admin).
2. **Credentials secret** — a per-device long-lived bearer secret obtained at registration; used
   *only* against the Pairing API. Stored bcrypt-hashed (cost 10); shown exactly once.
3. **mTLS client certificates** — short-lived X.509 issued by the per-realm CA against a
   device-generated CSR; the only credential the broker accepts.

### 4.2 JWT validation & Astarte authorization claims (`internal/auth`)

Astarte's claim model is reproduced exactly so existing tokens/tooling (`astartectl`) work:

- Claims: `a_aea` (AppEngine), `a_rma` (Realm Management), `a_pa` (Pairing), `a_ha`
  (Housekeeping), `a_ch` (Channels → honoured by Astrate's stream socket). Each is a list of
  authorization strings `"<HTTP-verb-regex>::<path-regex>"` (e.g. `"^POST$::^devices/.*$"`,
  or the catch-all `".*::.*"`), matched against the method and the path *relative to the realm
  base* (e.g. `devices/h4-Dx_RYTU-RbpDOTabhRg/interfaces/...`).
- Regexes are compiled once per token (LRU cache keyed by token hash) and are implicitly
  anchored as upstream does; `exp` honoured if present; `iat` not required (parity).
- Multiple realm public keys allow zero-downtime key rotation (`PUT
  /realmmanagement/v1/<realm>/config/auth` parity endpoint).

### 4.3 Embedded per-realm CA (`internal/pairing/ca`)

Replaces CFSSL:

- On realm creation, generate an ECDSA P-256 CA key + self-signed CA cert (configurable
  lifetime, default 10 years), or import an operator-provided pair. Private key encrypted at
  rest with AES-256-GCM under a master key supplied via env/file (documented: losing the master
  key ⇒ re-issue realm CAs; devices re-pair automatically at next credential rotation since
  their credentials secret still works).
- Issues client certs with: `Subject CN = <realm>/<device_id>`, `serial` = 128-bit random,
  validity = `pairing.cert_ttl` (default 30 days, matching common Astarte deployments;
  configurable down to hours for high-rotation setups), `KeyUsage = digitalSignature`,
  `ExtKeyUsage = clientAuth`. Public key taken from the CSR; **all other CSR-requested
  attributes (including its subject) are ignored and overridden** — the CSR is a proof of key
  possession, exactly as upstream treats it.
- Revocation: issuing new credentials records the new serial on the device row; the broker auth
  hook rejects any presented cert whose serial differs from the recorded latest (an
  always-online CRL equivalent — stricter than upstream's CRL-less default, toggleable via
  `pairing.enforce_latest_cert` for fleets that rotate while devices hold older still-valid
  certs). Inhibiting a device (`credentials_inhibited`, AppEngine PATCH parity) blocks both new
  credentials and new connections.

### 4.4 Pairing flows (wire-compatible)

**Flow A — Registration (operator/agent → Astrate):**

```
POST /pairing/v1/<realm>/agent/devices
Authorization: Bearer <realm JWT with a_pa>
{ "data": { "hw_id": "<22-char base64url 128-bit device ID>" } }
→ 201 { "data": { "credentials_secret": "<random 44-char base64>" } }
```

- `hw_id` validated as exactly 16 bytes after base64url-unpadded decode (UUID-shaped; both
  random and UUIDv5-derived deterministic IDs accepted — `pkg/deviceid`).
- Re-registering an existing device that has **not yet requested credentials** rotates the
  secret (parity); after first credentials request it conflicts (`422`, upstream-shaped error).
- Optional Astrate extension field `initial_payload_format` (§3.5.4).
- Unregister: `DELETE /pairing/v1/<realm>/agent/devices/<device_id>` (makes the device
  registrable again without losing its data — parity).

**Flow B — Credentials (device → Astrate), the SDK hot path:**

```
POST /pairing/v1/<realm>/devices/<device_id>/protocols/astarte_mqtt_v1/credentials
Authorization: Bearer <credentials_secret>
{ "data": { "csr": "-----BEGIN CERTIFICATE REQUEST-----..." } }
→ 201 { "data": { "client_crt": "-----BEGIN CERTIFICATE-----..." } }
```

Constant-time-ish secret verification (bcrypt compare against the stored hash; uniform error +
per-IP/device token-bucket rate limiting to blunt brute force). First successful call stamps
`first_credentials_request` and flips status `registered → confirmed`.

**Flow C — Broker discovery + cert health (device → Astrate):**

```
GET /pairing/v1/<realm>/devices/<device_id>
Authorization: Bearer <credentials_secret>
→ 200 { "data": {
    "status": "confirmed",
    "version": "<astrate version>",
    "protocols": { "astarte_mqtt_v1": {
        "broker_url": "mqtts://<host>:8883",
        "ca_crt": "<realm CA PEM>"
    } } } }

POST .../protocols/astarte_mqtt_v1/credentials/verify
{ "data": { "client_crt": "..." } }
→ 200 { "data": { "valid": true,  "timestamp": ..., "until": "..." } }
   |    { "data": { "valid": false, "cause": "EXPIRED" | "INVALID" | ... } }
```

These two endpoints are what lets stock SDK reconnect/renewal logic run unmodified: SDKs call
`verify` on boot and re-CSR when invalid/near expiry.

**Flow D — MQTT connection:** §3.1 (mTLS, CN check, client-ID = CN, session handling).

```
Agent                Astrate(pairing)          Device                 Astrate(broker)
  │ POST agent/devices    │                       │                        │
  │──────JWT(a_pa)───────▶│                       │                        │
  │◀──credentials_secret──│   (secret delivered   │                        │
  │ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┼── out of band ──────▶│                        │
  │                       │◀──POST credentials────│  keygen + CSR          │
  │                       │───client_crt─────────▶│                        │
  │                       │◀──GET device info─────│                        │
  │                       │───broker_url, ca_crt─▶│                        │
  │                       │                       │──CONNECT (mTLS, CN)───▶│
  │                       │                       │◀─CONNACK(session_present)│
  │                       │                       │── introspection, subs, │
  │                       │                       │   emptyCache, data ───▶│
```

### 4.5 Platform hardening checklist (v1 scope)

- TLS everywhere by default; HTTP listener TLS-terminated in-binary or behind a reverse proxy
  (documented compose profiles for both). HSTS on the API.
- Secrets handling: credentials secrets bcrypt-hashed; CA keys AES-GCM-encrypted; JWT public
  keys only (no symmetric verification accepted).
- Rate limits: pairing endpoints (per-IP and per-device), MQTT CONNECT storm damping
  (per-IP), AppEngine write endpoints (per-token).
- Input bounds everywhere: topic length ≤ 512 B, payload ≤ 64 KiB, introspection ≤ 64 KiB,
  zlib inflate hard-capped at the declared size and an absolute ceiling (zip-bomb guard on
  `producer/properties`).
- Single non-root distroless container; read-only FS except the session-store volume.

---

## 5. Cross-Cutting Concerns

### 5.1 Configuration

One TOML file + `ASTRATE_*` env overrides: listeners, Postgres DSN, shard count, chunk/compress/
retention knobs, cert TTL, master encryption key ref, dev-mode flags, auto-provision realm.
Sane zero-config defaults for the single-VPS case.

### 5.2 Observability

- Prometheus metrics (`/astrate/v1/metrics`): ingest rate, per-reason rejects, shard depth,
  batch flush latency, broker sessions, DB pool stats, trigger delivery outcomes.
- Structured logging (`log/slog`, JSON), per-domain levels.
- `/astrate/v1/health` (liveness) and `/readiness` (DB + broker checks).

### 5.3 Lifecycle & resilience

- Graceful shutdown: broker stops accepting, shards drain (bounded by timeout), batches flush.
- Crash safety: QoS ≥ 1 messages are PUBACK'd **only after** the persistence batch commits —
  at-least-once into Postgres; datastream inserts are idempotence-tolerant (duplicate
  (series, ts) rows are acceptable per Astarte semantics; properties are upserts).
- DB outage: shards park with exponential backoff; broker applies backpressure (§1.4); QoS 0
  data degrades first, by design.

### 5.4 Deployment

`docker-compose.yml`: `timescale/timescaledb:latest-pg16` (tuned: `shared_buffers=256MB`,
`max_connections=50`) + `astrate` (distroless, ~20 MB image). Volumes: pgdata + session store.
Or: one binary + one `apt install postgresql-16 timescaledb` on a VPS. Target steady-state RSS:
Astrate ≤ 150 MB at 1k devices; Postgres tuned to ≤ 768 MB.

---

## 6. Key Risks & Mitigations

| Risk | Mitigation |
|---|---|
| Subtle protocol drift vs. real SDKs | Phase 2 defines CI conformance suites driving the **official** Astarte Go + Python SDKs (and an AtomVM JSON simulator) against Astrate end-to-end. |
| Per-interface TTL via DELETE is heavier than `drop_chunks` | Batched chunk-aware deletes off-peak; global chunk-drop cap as the backstop; revisit table-per-retention-class if profiling demands. |
| mochi-mqtt session persistence semantics under restart | Dedicated storage-hook tests for `session_present`, offline queue replay order, and QoS 2 exactly-once handshake. |
| Embedded CA key compromise | Encrypted at rest, short cert TTLs, latest-serial enforcement, documented re-keying runbook. |
| Single-process blast radius (vs. microservices) | Domain isolation via package contracts, panics confined per shard/connection with recover-and-log, optional hot-standby via shared Postgres + LISTEN/NOTIFY. |

---

## 7. Phase 1 Exit Criteria → Phase 2 Preview

This document freezes: the module decomposition (§1.3), the storage model (§2), the wire
contracts and the dual-payload profile (§3), and the pairing/security design (§4).

Upon approval, **Phase 2** will produce the file-by-file implementation roadmap in dependency
order — `pkg/deviceid` → `pkg/interfaceschema` → `pkg/payload` → `internal/store` + migrations →
`internal/auth` → `internal/pairing(+ca)` → `internal/broker` → `internal/engine` →
`internal/realm`/`appengine`/`housekeeping` → `cmd/astrate` — each step with its unit,
integration (Testcontainers: TimescaleDB), and SDK-conformance verification tests.
