# Astrate — Implementation Roadmap & Milestones (Phase 2)

**Project:** Astrate — lean, single-binary, Astarte-wire-compatible IoT platform in Go
**Status:** Phase 2 deliverable — awaiting approval before Phase 3 (code generation)
**Version:** 0.1 (2026-06-10)
**Normative reference:** `docs/DESIGN.md` (Phase 1, approved). Section references (§) below point there.

---

## 0. How to Read This Roadmap

### 0.1 Sequential-generation rules (binding for Phase 3)

1. **Strict dependency order.** Files are listed in generation order; no file may import or
   reference a file that appears later in the roadmap. Every milestone ends with a tree that
   compiles (`go build ./...`) and a green test suite for everything generated so far.
2. **Tests ride with their subject.** Each `foo.go` is immediately followed by `foo_test.go`
   (plus `testdata/` fixtures) in the same generation step — never deferred to a later milestone.
3. **No placeholders.** `// TODO` is forbidden except at the three named extension points:
   external-bus intake (§1.4), NATS/HTTP trigger forwarding (§1.1), and Timescale toolkit `lttb`
   probing (§2.5) — each must still compile and have a working default path.
4. **Interfaces before implementations.** Where package A consumes package B through an
   interface (hexagonal-lite, §1.3), the interface is defined in A's milestone step, with B
   already generated — so wiring never requires editing earlier files.
5. **One milestone per Phase 3 generation session** (M6 and M7 split into two sessions each,
   see §11). Each session's exit criterion is the milestone gate, not "code emitted".

### 0.2 Verification tiers

| Tier | Name | Mechanics | When it runs |
|---|---|---|---|
| **T1** | Unit | Pure Go, `go test ./...`, no Docker, no network | Every commit |
| **T2** | Integration | `//go:build integration`; testcontainers-go boots `timescale/timescaledb:latest-pg16` | Every commit (CI), on-demand locally |
| **T3** | Component / E2E | Full wired binary (in-process `run()` or compose); real broker + real DB + test MQTT/HTTP clients | Every commit (CI) |
| **T4** | Conformance | **Official Astarte SDKs** (Go, Python), `astartectl`, and the AtomVM-profile JSON simulator driven against a composed Astrate instance | PRs to `main` + nightly |
| **T5** | Non-functional | Load/footprint budgets, security probes | Nightly |

### 0.3 Conformance checkpoints (the anti-protocol-drift spine, §6 risk 1)

Conformance is not a final-milestone afterthought; it gates progress at four checkpoints:

| Checkpoint | After | What must pass against unmodified official clients |
|---|---|---|
| **CP-A** | M4 | Pairing REST: device registration + CSR/credentials handshake using the official **Go SDK pairing client** and `astartectl pairing agent register` |
| **CP-B** | M6 | Full device loop with the official **Go SDK**: register → mTLS connect → introspection → datastreams (individual + object, BSON) → properties (set/unset) → server-owned data delivery → `emptyCache` |
| **CP-C** | M7 | `astartectl` smoke: housekeeping realm create → interface install → device register → appengine data queries → trigger CRUD |
| **CP-D** | M9 | Full matrix: Go SDK + **Python SDK** + AtomVM JSON simulator + reconnect/session_present scenarios + `astartectl` regression |

Pinned conformance clients (upgraded deliberately, never floating):
`astarte-device-sdk-go` (latest tagged release at M4 start), `astarte-device-sdk-python`
(PyPI, pinned), `astartectl` (pinned release binary).

### 0.4 Milestone overview

| Milestone | Deliverable | Gate |
|---|---|---|
| **M0** | Repo scaffolding, dependency pinning, test harness, CI skeleton | CI green on empty-ish tree; Timescale container helper boots |
| **M1** | Foundation libraries: `pkg/deviceid`, `pkg/astarteapi`, `pkg/interfaceschema`, `pkg/payload` | T1 green incl. fuzz + zero-alloc benchmarks; golden SDK payload vectors decode |
| **M2** | SQL migrations + `internal/store` repositories | T2 green; hypertables/compression/TTL job verified against live TimescaleDB |
| **M3** | `internal/auth` — JWT + Astarte authz claims | T1 green; claim-matching parity table passes |
| **M4** | `internal/pairing` + embedded CA + Pairing API | T2 green; **CP-A** |
| **M5** | `internal/broker` — embedded mochi-mqtt, mTLS, ACL, sessions | T3 green: auth/ACL matrix, session persistence across restart |
| **M6** | `internal/engine` — shards, validation, persistence, triggers, stream bus | T3 green: ordered ingest E2E both payload formats; **CP-B** |
| **M7** | `internal/realm`, `internal/housekeeping`, `internal/appengine` (+ stream socket) | T2/T3 green; **CP-C** |
| **M8** | `internal/config`, `cmd/astrate`, observability, Dockerfile, compose | Boot/shutdown E2E; compose smoke in CI |
| **M9** | Conformance suite, load/footprint, security hardening, ops docs | **CP-D**; T5 budgets met → v0.1.0 readiness |

Dependency graph (arrows = "required by"):

```
M0 ─▶ M1 ─▶ M2 ─▶ M3 ─▶ M4 ─▶ M5 ─▶ M6 ─▶ M7 ─▶ M8 ─▶ M9
       │     │           │      ▲      ▲     ▲
       │     └───────────┴──────┘      │     │
       └── (pkg/* feeds every later milestone) ┘
```

---

## 1. M0 — Scaffolding, Dependencies, Test Harness, CI Skeleton

**Goal:** a repo where every later milestone only *adds* files. All third-party dependency
choices are pinned here so Phase 3 never has to re-litigate them.

### 1.1 Pinned dependency decisions

| Concern | Choice | Rationale (lean-first) |
|---|---|---|
| HTTP routing | **stdlib `net/http`** (Go ≥ 1.22 `ServeMux` method+wildcard patterns) | Zero deps; upstream-compatible paths need no framework |
| MQTT broker | `github.com/mochi-mqtt/server/v2` | Frozen in §1.1 |
| Postgres | `github.com/jackc/pgx/v5` (+ `pgxpool`) | Frozen in §1.3 |
| Migrations | `github.com/golang-migrate/migrate/v4` with `source/iofs` + `go:embed` | Frozen in §1.3 |
| BSON | `go.mongodb.org/mongo-driver/v2/bson` (raw-document API only) | Frozen in §3.5.5 |
| JWT | `github.com/golang-jwt/jwt/v5` | Maintained, alg-allowlist friendly |
| Password hashing | `golang.org/x/crypto/bcrypt` (cost 10) | Frozen in §4.1 |
| Session store | `go.etcd.io/bbolt` | §3.1 said bbolt/pebble — **pinned: bbolt** (simpler, single file) |
| LRU (token cache) | `github.com/hashicorp/golang-lru/v2` | §4.2 |
| Metrics | `github.com/prometheus/client_golang` | §5.2 |
| WebSocket | `github.com/coder/websocket` | Minimal, context-native |
| TOML config | `github.com/BurntSushi/toml` | §5.1 |
| Test containers | `github.com/testcontainers/testcontainers-go` (+ postgres module), image `timescale/timescaledb:latest-pg16` | §5.4 parity with production image |
| Test MQTT client | `github.com/eclipse/paho.mqtt.golang` (test-only) | Same client family the official Go SDK uses |

### 1.2 Files

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 0.1 | `go.mod` | Module path, Go ≥ 1.22, pinned deps above | — |
| 0.2 | `.gitignore` | Go defaults + `*.db` (bbolt), `dist/` | — |
| 0.3 | `Makefile` | `build` (CGO_ENABLED=0, `-trimpath`, static), `lint`, `test`, `test-integration`, `test-e2e`, `test-conformance`, `up`/`down` (compose) | 60 |
| 0.4 | `.golangci.yml` | govet, staticcheck, errcheck, gosec, revive; test exclusions | 40 |
| 0.5 | `internal/testutil/pg.go` | `StartTimescale(t) *pgxpool.Pool`: container (or `ASTRATE_TEST_DSN` reuse for speed), waits ready, returns pool + teardown | 90 |
| 0.6 | `internal/testutil/golden.go` | Golden-file compare/update helper (`-update` flag) | 50 |
| 0.7 | `internal/testutil/pg_test.go` | T2: container boots; `SELECT extversion FROM pg_extension WHERE extname='timescaledb'` non-empty | 30 |
| 0.8 | `docker-compose.yml` | `timescaledb` service only (tuned per §5.4: `shared_buffers=256MB`, `max_connections=50`); `astrate` service added in M8 | 30 |
| 0.9 | `.github/workflows/ci.yml` | Jobs: `lint`, `unit` (T1), `integration` (T2), placeholders for `e2e`/`conformance`/`nightly` enabled in M6/M9 | 80 |
| 0.10 | `README.md` | One-paragraph statement + pointer to docs/; expanded in M8 | 30 |

### 1.3 Verification / gate

- `make lint test` green; T2 job boots the Timescale container in CI.
- **Gate:** CI pipeline green end-to-end on the scaffold.

---

## 2. M1 — Foundation Libraries (`pkg/`)

**Goal:** the four dependency-free libraries every domain consumes. Pure logic, exhaustively
unit-tested — this is where wire-compat correctness is cheapest to lock in.
**Dependency rule check:** `pkg/*` imports stdlib + bson only; zero `internal/*` imports (§1.3).

### 2.1 Step M1.1 — `pkg/deviceid`

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 1.1 | `pkg/deviceid/deviceid.go` | `type ID [16]byte`; `Parse(s string) (ID, error)` (exactly 22-char unpadded base64url → 16 bytes), `String()`, `FromUUID`/`UUID()`, `Random()`, `FromNamespace(ns uuid, payload string)` (UUIDv5 deterministic derivation, astartectl parity) | 140 |
| 1.2 | `pkg/deviceid/deviceid_test.go` | See below | 160 |

**Verification (T1):**
- Round-trip property: `Parse(id.String()) == id` for random IDs (1000 iterations).
- Known vectors: at least 3 `(UUID, base64url)` pairs cross-checked against `astartectl utils device-id` output (vendored as constants with provenance comment).
- Rejections: 21/23-char strings, padded base64, standard-alphabet chars (`+`, `/`), non-ASCII, empty.
- `FuzzParse`: never panics; accepts ⇒ round-trips.

### 2.2 Step M1.2 — `pkg/astarteapi`

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 1.3 | `pkg/astarteapi/envelope.go` | `WriteData(w, status, v)` → `{"data": v}`; `WriteError(w, status, detail)` → `{"errors":{"detail":"..."}}`; canonical upstream error constructors (404 `"Device not found"`, 401, 403, 422 shapes); request body `{"data": ...}` unwrapper with size cap | 120 |
| 1.4 | `pkg/astarteapi/envelope_test.go` | Golden JSON for every constructor; unwrap rejects missing `data`, oversized bodies | 90 |

**Verification (T1):** golden envelopes byte-compared (these exact bytes are what `astartectl`
and SDK error paths parse — treated as frozen fixtures from here on).

### 2.3 Step M1.3 — `pkg/interfaceschema`

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 1.5 | `pkg/interfaceschema/types.go` | Enums + parsers: `InterfaceType`, `Ownership`, `Aggregation`, `Reliability` (→ QoS byte), `Retention`, `DatabaseRetentionPolicy`, `ValueType` (all 14: 7 scalars + 7 arrays) | 150 |
| 1.6 | `pkg/interfaceschema/parse.go` | `ParseInterface([]byte) (*Interface, error)`: strict decode of the Astarte Interface JSON; validation — name regex (RFC-compliant reverse-domain, ≤ 128 chars), `version_major/minor` ≥ 0 not both 0, endpoint syntax (`/`-rooted, `%{param}` whole-segment-only, no duplicate/conflicting endpoints), object-aggregation constraints (datastream only, same depth, uniform `explicit_timestamp`/`reliability`/`retention`/`expiry`, last level distinct), properties constraints (no aggregation=object, no explicit_timestamp; `allow_unset` legal), per-field allowed-values | 320 |
| 1.7 | `pkg/interfaceschema/parse_test.go` | Fixture-driven; see below | 200 |
| 1.8 | `pkg/interfaceschema/testdata/valid/*.json` | ≥ 12 fixtures: vendored upstream `astarte-platform/standard-interfaces` (genericsensors, device-info…) + object-aggregated + parametric + properties w/ allow_unset + every value type | — |
| 1.9 | `pkg/interfaceschema/testdata/invalid/*.json` | ≥ 15 fixtures, one rule violation each, expected error substring in a manifest | — |
| 1.10 | `pkg/interfaceschema/trie.go` | `EndpointTrie`: build from endpoints; `Match(path string) (*CompiledMapping, bool)` — segment-wise, exact child before parametric child, zero-alloc (sub-slicing, no `strings.Split`); placeholder hygiene checks (§2.6: non-empty, ≤ 256 B, no `/`, `+`, `#`) | 180 |
| 1.11 | `pkg/interfaceschema/trie_test.go` | Match table incl. parametric/mixed/miss/overlong; `BenchmarkMatch` + `TestMatchZeroAllocs` (`testing.AllocsPerRun == 0`) | 150 |
| 1.12 | `pkg/interfaceschema/compile.go` | `Compile(*Interface, ids EndpointIDResolver) (*CompiledInterface, error)` → §2.6 structs verbatim (`CompiledInterface`, `CompiledMapping`, `ObjectLeaves`) | 140 |
| 1.13 | `pkg/interfaceschema/compat.go` | `CheckMinorUpgrade(old, new *Interface) error` — additive-mappings-only, no mutation of existing mapping attributes, same type/ownership/aggregation (§2.6 versioning parity) | 100 |
| 1.14 | `pkg/interfaceschema/compile_test.go`, `compat_test.go` | Compile fixtures → expected tries/leaves; upgrade accept/reject table | 160 |

**Verification (T1):** all fixtures parse/reject as expected; zero-alloc benchmark gate is a CI
assertion, not advisory; `FuzzParseInterface` never panics.

### 2.4 Step M1.4 — `pkg/payload`

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 1.15 | `pkg/payload/value.go` | `Value` (decoded scalar/array/object), `DecodedPayload{Value, Timestamp *time.Time}`, coercion table per §2.6-step-5 (double widening, integer fit, longinteger string form, UTF-8 + 64 KiB string cap, ≤ 1024 array elements, homogeneity) | 200 |
| 1.16 | `pkg/payload/sniff.go` | `DetectFormat([]byte) Format` — exact §3.5.2 algorithm (`empty` / `bson` / `json` / `invalid`) | 50 |
| 1.17 | `pkg/payload/bson.go` | Decode via `bson.Raw` lookups of `v`/`t` (no maps/reflection); encode `{v,t}` for outbound | 220 |
| 1.18 | `pkg/payload/json.go` | Strict §3.5.3 profile: envelope mandatory, `t` RFC 3339 or epoch-ms, base64 binaryblob, longinteger decimal-string, object aggregation; encoder for outbound | 200 |
| 1.19 | `pkg/payload/payload.go` | Facade: `Decode(p []byte, m *CompiledMapping / leaves) (DecodedPayload, error)`; `Encode(v, t, Format) ([]byte, error)`; size caps; typed reject reasons (enum shared with engine metrics) | 120 |
| 1.20 | `pkg/payload/testdata/bson/*.hex` | **Golden vectors captured from the official Go SDK encoder** (one per value type, plus object aggregation, plus with/without `t`); provenance documented in a README | — |
| 1.21 | `pkg/payload/payload_test.go` (+ `sniff_test.go`, `fuzz_test.go`) | See below | 350 |

**Verification (T1):**
- Every golden BSON vector decodes to the expected `Value` for its declared type; re-encoding produces semantically-equal BSON (field order `v`,`t` matched to SDK output).
- JSON profile table: every type × valid/invalid (incl. `22.5` bare-shorthand **rejected**, `t` both forms, > 2^53 longinteger as string, unpadded base64 rejected).
- Coercion edge cases: int64→double lossless only, double→integer exact-fit only, NaN/±Inf rejected, datetime range bounds.
- `TestSniffNoCollision`: generated corpus of valid JSON docs never classifies as BSON and vice versa; `FuzzDecode` (both formats) never panics, never misclassifies the format of valid inputs.
- 64 KiB cap enforced for both formats.

**M1 gate:** T1 green including fuzz seeds and alloc benchmarks; golden vectors committed.

---

## 3. M2 — Migrations + `internal/store`

**Goal:** the single Postgres access layer (§1.3: imported by domains, imports none of them).
Schema is transcribed **verbatim from §2.2–2.4** — any deviation found while implementing is a
design-change request, not an ad-hoc fix.

### 3.1 Files

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 2.1 | `migrations/000001_extensions.up.sql` / `.down.sql` | `CREATE EXTENSION IF NOT EXISTS timescaledb` | 10 |
| 2.2 | `migrations/000002_metadata.up.sql` / `.down.sql` | `realms`, `interfaces` (+generated columns), `endpoints`, `devices` (+aliases GIN), `groups`, `group_devices`, `triggers` — §2.2 verbatim | 130 |
| 2.3 | `migrations/000003_properties.up.sql` / `.down.sql` | `properties` — §2.3 verbatim | 30 |
| 2.4 | `migrations/000004_datastreams.up.sql` / `.down.sql` | `individual_datastreams` + `object_datastreams`: hypertables (7-day chunks), series indexes, compression settings + policies — §2.4 verbatim | 80 |
| 2.5 | `migrations/000005_jobs.up.sql` / `.down.sql` | TTL user-defined action (procedure `astrate_apply_endpoint_ttl()`, chunk-batched DELETE) + `add_job(..., '1 hour')`; optional global `add_retention_policy` applied at runtime from config | 60 |
| 2.6 | `internal/store/store.go` | `New(ctx, dsn) (*Store, error)`: pgxpool (pool sized for §5.4 `max_connections=50` budget), embedded migration runner (`go:embed migrations`), `Health(ctx)`, toolkit/`lttb` capability probe (§2.5) | 160 |
| 2.7 | `internal/store/crypto.go` | AES-256-GCM seal/open for CA private keys; master key from env/file ref (§2.2, §4.3); random nonce per seal | 80 |
| 2.8 | `internal/store/realms.go` | CRUD + `GetByName`; transactional create (realm row + CA material); delete = cascade in one tx (§2.1) | 130 |
| 2.9 | `internal/store/interfaces.go` | Install (interface + endpoints rows in tx), update (minor bump), delete (only if major draining rules satisfied), `LoadRealm(realmID) []StoredInterface` (definition + endpoint IDs for compiler), `EndpointIDResolver` impl | 180 |
| 2.10 | `internal/store/devices.go` | Register/rotate-secret, get/list (paged), introspection update (with `old_introspection` diff), status/inhibit, cert serial+AKI stamp, connection lifecycle updates (`connected`, IPs, timestamps, counters), aliases/attributes patch, by-alias lookup, `payload_format_hint` flip | 220 |
| 2.11 | `internal/store/properties.go` | Upsert, unset (DELETE), get-by-device/interface, device-owned purge-not-in-list (for `producer/properties`, §3.3), server-owned list (for `consumer/properties`) | 120 |
| 2.12 | `internal/store/datastreams.go` | `AppendBatch([]Row)` via `pgx.Batch`/`COPY` into the typed sparse columns / object JSONB; query API: `Series(realm, device, iface, path, opts)` with `since/since_after/to/limit/asc-desc`, `Downsample(..., bucket)` via `time_bucket` (+`lttb` when probed) | 240 |
| 2.13 | `internal/store/groups.go` | Group CRUD + membership (composite FK semantics) | 90 |
| 2.14 | `internal/store/triggers.go` | Trigger CRUD, list-by-realm for engine cache | 70 |
| 2.15 | `internal/store/notify.go` | `NOTIFY astrate_interfaces` emit helper + `Listen(ctx, channel) <-chan Notification` (dedicated conn, auto-reconnect w/ backoff) (§2.6) | 110 |
| 2.16 | `internal/store/*_test.go` (one per repo file) | T2 suites, shared container via `testutil.StartTimescale` | 700 |

### 3.2 Verification (T2 unless noted)

- **Migrations:** fresh-DB up succeeds; `schema_migrations` at head; **catalog assertions** — both hypertables exist (`timescaledb_information.hypertables`), compression enabled with exact `segmentby`/`orderby` (§2.4), compression policy + TTL job registered (`timescaledb_information.jobs`).
- **Realms:** create→get round-trip; CA key seal/open round-trip (T1 for crypto.go: wrong key fails, nonce uniqueness); cascade delete removes interfaces/devices/properties/datastream rows.
- **Interfaces:** install populates generated columns correctly; endpoint IDs stable across reloads; unique `(realm, name, major)` enforced.
- **Devices:** introspection diff moves removed pairs to `old_introspection`; inhibit flag round-trip; alias GIN lookup; pagination cursor stability.
- **Properties:** upsert idempotent; unset deletes; purge-not-in-list deletes exactly the complement (device-owned rows only).
- **Datastreams:** batch of 10k rows lands < 2 s in CI (smoke, non-binding); exactly-one-value-column invariant checked per type; `Series` golden results incl. boundary semantics of `since` vs `since_after`; `Downsample` results match hand-computed `time_bucket` expectations; duplicate `(series, ts)` insert tolerated (§5.3).
- **TTL job:** insert aged rows for a `use_ttl` endpoint, `CALL astrate_apply_endpoint_ttl()`, verify only aged+TTL'd rows gone.
- **Notify:** emit→receive round-trip; listener survives connection kill (`pg_terminate_backend`) and resumes.

**M2 gate:** T2 green in CI; schema matches §2.2–2.4 byte-for-byte (reviewed diff).

---

## 4. M3 — `internal/auth`

**Goal:** the JWT + Astarte-claims layer used by every REST surface (§4.2). Pure T1 — no DB
(key material is injected via a `KeySource` interface; `internal/store` already satisfies it).

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 3.1 | `internal/auth/jwt.go` | Parse+verify against realm key set (PEM multi-key); alg allowlist `RS256/384/512, ES256/384/512`; `none`/HMAC hard-rejected; `exp` honoured, `iat` optional (§4.2) | 140 |
| 3.2 | `internal/auth/claims.go` | `a_aea/a_rma/a_pa/a_ha/a_ch` extraction; `"<verb-regex>::<path-regex>"` matching with **implicit anchoring** (upstream parity), evaluated against method + path relative to realm base | 130 |
| 3.3 | `internal/auth/cache.go` | LRU (size 1024) keyed by SHA-256(token): verified claims + compiled regexes (§4.2) | 70 |
| 3.4 | `internal/auth/middleware.go` | `Require(api Claim, keys KeySource) func(http.Handler) http.Handler` — extracts realm from path, 401/403 via `astarteapi` envelopes; housekeeping variant (instance-level keys) | 110 |
| 3.5 | `internal/auth/*_test.go` | See below | 320 |

**Verification (T1):**
- Parity table for claim matching (the cases `astartectl` emits): `".*::.*"` catch-all; `"^POST$::^devices/.*$"` allows `POST devices/x/...`, denies `GET`; **anchoring**: pattern `devices` must NOT match `devices/abc` unless `.*`-suffixed; multiple authz strings = OR.
- RSA + EC test keys (generated in-test): valid/expired/not-yet-valid/garbage tokens; HMAC token signed with the public key bytes rejected (classic confusion attack); `alg: none` rejected.
- Key rotation: token verifies if *any* of N realm keys matches.
- Cache: hit avoids re-verification (call-count spy); distinct tokens don't collide.
- Middleware: 401 (no/bad token) vs 403 (valid token, claim mismatch) with golden envelopes.

**M3 gate:** T1 green; parity table reviewed against upstream `astarte_rpc`/dashboard token semantics.

---

## 5. M4 — `internal/pairing` + Embedded CA

**Goal:** wire-identical Pairing API (§4.4 flows A–C) + the CFSSL-replacing CA (§4.3).

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 4.1 | `internal/pairing/ca/ca.go` | `NewRealmCA(lifetime)` (ECDSA P-256, self-signed, default 10 y) or import; `SignCSR(csr, realm, deviceID, ttl) (certPEM, serial, aki, err)` — CN forced to `<realm>/<device_id>`, 128-bit random serial, `KeyUsage=digitalSignature`, `EKU=clientAuth`, **all CSR subject/attrs ignored** (§4.3); `VerifyCert(certPEM)` for the verify endpoint | 200 |
| 4.2 | `internal/pairing/ca/ca_test.go` | T1: field assertions via `x509.ParseCertificate`; CSR-with-hostile-subject override; chain verifies to CA; expired-CA rejection; serial uniqueness (10k draw) | 180 |
| 4.3 | `internal/pairing/service.go` | Flows A–C business logic over `store`: register (secret gen 32 B → 44-char base64, bcrypt store, show-once; re-register rotates iff no credentials request yet else 422; `initial_payload_format` extension §3.5.4), unregister (re-registrable, data kept), credentials (bcrypt compare, uniform error, status flip `registered→confirmed`, stamp first-request + IP, record serial/AKI), info (broker_url + realm `ca_crt`), verify (valid/until vs cause `EXPIRED|INVALID|REVOKED`) | 260 |
| 4.4 | `internal/pairing/ratelimit.go` | Token bucket per-IP and per-device for `/credentials` + register (§4.5); 429 envelope | 80 |
| 4.5 | `internal/pairing/http.go` | Routes: `POST/DELETE /pairing/v1/{realm}/agent/devices[/{id}]` (JWT `a_pa` via M3), `POST .../devices/{id}/protocols/astarte_mqtt_v1/credentials`, `GET .../devices/{id}`, `POST .../credentials/verify` (Bearer credentials-secret auth) — bodies/status codes per §4.4 verbatim | 180 |
| 4.6 | `internal/pairing/service_test.go`, `http_test.go` | T1 (service w/ store iface fake where cheap) + T2 (full HTTP over real store) | 420 |

**Verification:**
- **T2 HTTP flow tests (golden bodies):** Flow A 201 + 44-char secret; second register pre-credentials returns a *different* secret; post-credentials register → 422 upstream-shaped; bad `hw_id` (21 chars, padded, non-url) → 422. Flow B: 201 `client_crt` that (a) parses, (b) has CN `<realm>/<device>`, (c) chains to the realm CA returned by Flow C — verified with `crypto/x509` *and* an `openssl verify` exec smoke; wrong secret → 401 with the same body/timing class as unknown device (uniform error assertion); rate limit → 429. Flow C: info golden body (§4.4 shape exactly); verify on fresh cert → `valid:true` + `until`; on expired (issue with 1 s TTL, sleep) → `EXPIRED`; on cert from a different CA → `INVALID`.
- Inhibited device: credentials → 403 parity.
- **CP-A (T4):** CI job runs the pinned official Go SDK's registration + credentials code (and `astartectl pairing agent register --realm test ...`) against a composed Astrate-pairing instance; asserts the SDK obtains a certificate it considers valid.

**M4 gate:** T2 green + CP-A green.

---

## 6. M5 — `internal/broker`

**Goal:** the embedded mochi-mqtt broker with Astarte semantics (§3.1–3.2): mTLS identity,
ACLs, persistent sessions, lifecycle events, inline publishing. The engine doesn't exist yet, so
the broker emits into an `Intake` interface defined here (M6 implements it; tests use a recorder).

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 5.1 | `internal/broker/intake.go` | `InboundMessage{Realm, DeviceID, Topic, Payload, QoS, ReceivedAt, Ack func()}`; `type Intake interface { Submit(InboundMessage) }`; `LifecycleSink` (connect/disconnect events) | 60 |
| 5.2 | `internal/broker/identity.go` | Parse CN `<realm>/<device_id>` → identity; topic splitter `realm/device/rest` (§3.3 parsing note) | 80 |
| 5.3 | `internal/broker/authhook.go` | `OnConnectAuthenticate`: peer-cert CN parse, chain to *that realm's* CA (per-realm pool, hot-reloadable on realm CRUD), device exists + not inhibited, **client-ID == CN**, latest-serial enforcement behind `pairing.enforce_latest_cert` (§3.1, §4.3) | 160 |
| 5.4 | `internal/broker/aclhook.go` | `OnACLCheck` per §3.2 matrix exactly (publish: base, `control/emptyCache`, `control/producer/properties`, device-owned interface topics from cached introspection; subscribe: `control/consumer/properties`, server-owned `iface/#`, tolerated `B/#` superset); deny+log otherwise | 150 |
| 5.5 | `internal/broker/sessionstore.go` | bbolt-backed mochi storage hook (subscriptions, inflight, offline queue) at configured path; per-message expiry honoured on offline queue (§2.5, §3.4) | 140 |
| 5.6 | `internal/broker/lifecycle.go` | Connect/disconnect → `devices` row updates (connected, IPs, timestamps) + `LifecycleSink` events (`device_connected`/`device_disconnected` trigger feed, §3.1) | 90 |
| 5.7 | `internal/broker/publisher.go` | Inline-client facade: `Publish(topic, payload, qos, retain bool, expiry time.Duration) error` — bypasses ACL (§3.2); used by engine/appengine | 70 |
| 5.8 | `internal/broker/broker.go` | `New(cfg, store, intake, sink)`: mochi server, TLS listener `:8883` (`RequireAndVerifyClientCert`, client-CA pool), optional `:1883` behind `insecure_dev_mode`, MQTT 3.1.1/5.0 accept, hook registration, `OnPublish`→`Intake.Submit` with **deferred-ack wiring** (QoS ≥ 1 PUBACK held until `Ack()`, §1.4/§5.3), graceful stop | 220 |
| 5.9 | `internal/broker/*_test.go` + `internal/testutil/mqttclient.go` | T1 for identity/ACL pure logic; T3 suite with real TLS broker, M4-CA-issued client certs, paho clients, recorder Intake | 600 |

**Verification:**
- **T1:** identity parse table (bad CN shapes); ACL matrix as an exhaustive table test (every row of §3.2 allowed; ~20 adversarial topics — other device's base, other realm, `control/consumer` publish, server-owned publish, wildcard abuse — denied).
- **T3 (real broker, real certs):** connect happy path; reject: wrong-realm CA cert, inhibited device, stale serial (after issuing a newer one, flag on), client-ID ≠ CN, plaintext attempt with dev-mode off. Session persistence: connect `clean_session=false` → disconnect → broker process restart (new `New()` on same bbolt file) → reconnect sees `session_present=1`; QoS 1 messages published to a subscribed offline topic replay **in order** on reconnect; QoS 2 exactly-once handshake completes (dedicated test per §6 risk 3); offline message with 1 s expiry not delivered after 2 s. Deferred ack: recorder Intake that never calls `Ack()` ⇒ client's QoS 1 publish never PUBACKed (backpressure proof); on `Ack()` ⇒ PUBACK arrives. Lifecycle: device row flips connected/timestamps/IP on connect+disconnect; events arrive at sink.

**M5 gate:** T3 broker suite green (including the restart-persistence tests).

---

## 7. M6 — `internal/engine` (heart of the system)

**Goal:** sharded ordered pipeline (§1.4) + validation (§2.6) + persistence + control channel
semantics (§3.3–3.4) + triggers + live fan-out. Split into two Phase 3 sessions: **M6a**
(pipeline + data path) and **M6b** (control, server-data, triggers, stream).

### 7.1 M6a — pipeline + data path

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 6.1 | `internal/engine/cache.go` | Compiled-interface snapshot `map[realmID]map[name]map[major]*CompiledInterface` behind `atomic.Pointer`, copy-on-write rebuild from `store.LoadRealm`; invalidation via in-process callback **and** `store.Listen("astrate_interfaces")` (§2.6); device-introspection cache (per-connected-device, evicted on disconnect) | 200 |
| 6.2 | `internal/engine/router.go` | `Engine.Submit(InboundMessage)` (implements `broker.Intake`): `shard = FNV1a(deviceID) % N` (default 16), bounded chans (default 4096); QoS ≥ 1 full ⇒ blocking submit (ack deferral = backpressure); QoS 0 full ⇒ drop + metric (§1.4); per-shard goroutine with panic recover-and-log (§6) | 160 |
| 6.3 | `internal/engine/topics.go` | Classify `rest` → introspection / control / data(interface, path) via longest-prefix match against the device's introspected interface names (§3.3) | 90 |
| 6.4 | `internal/engine/data.go` | §2.6 pipeline verbatim: introspection gate → trie match → ownership → `payload.Decode` → type check → aggregation shape → emit `PersistOp`; per-reason reject enum → metrics + `device_error` trigger events; property unset path (`allow_unset`) | 260 |
| 6.5 | `internal/engine/batch.go` | Per-shard micro-batch: flush at 64 rows or 50 ms (§1.4) through `store.AppendBatch` / property upserts in one tx; **`Ack()` called only after commit** (§5.3); DB-outage parking with exponential backoff | 200 |
| 6.6 | `internal/engine/*_test.go` (M6a) | See below | 550 |

### 7.2 M6b — control, server data, triggers, stream

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 6.7 | `internal/engine/introspection.go` | Parse `;`-separated triples; diff vs stored; persist (+`old_introspection`); fire `incoming_introspection`/added/removed trigger events; refresh ACL-relevant introspection cache (§3.3) | 140 |
| 6.8 | `internal/engine/control.go` | `emptyCache`: re-send server-owned properties on data topics (QoS 2, device-hint format) then publish `consumer/properties` purge; `producer/properties`: 4-byte BE size + zlib inflate (declared-size cap + absolute ceiling, §4.5) → purge device-owned properties not listed; `consumer/properties` builder (deflate, size prefix) — sent after emptyCache, after `session_present=0` resume, after server-owned unset (§3.4) | 220 |
| 6.9 | `internal/engine/serverdata.go` | Server-owned publish path for AppEngine: validate against `ownership: server` interface → persist (property upsert/unset or datastream insert) → `broker.Publish` with mapping QoS, retain-for-properties, expiry-for-datastreams, **format per `payload_format_hint`** (+ sticky hint flip on first JSON device payload, §3.5.4) | 160 |
| 6.10 | `internal/engine/triggers/match.go` | Compile stored trigger JSON (§2.2) → matchers: data triggers (interface/path, value conditions), device triggers (connected/disconnected/error, introspection); evaluated post-persist in-shard | 200 |
| 6.11 | `internal/engine/triggers/events.go` | Astarte trigger event JSON payload shapes (SimpleEvent envelope parity: realm, device_id, event type fields) | 120 |
| 6.12 | `internal/engine/triggers/actions.go` | HTTP webhook action: template URL/headers, POST event JSON, retry w/ exponential backoff + cap, delivery-outcome metrics; extension point stub for NATS/HTTP forwarding (§1.1) | 160 |
| 6.13 | `internal/engine/stream/bus.go` | In-process fan-out: `Subscribe(realm, filter) (<-chan Event, cancel)`; non-blocking sends, slow-consumer drop + metric (§1.1/§1.4) | 100 |
| 6.14 | `internal/engine/engine.go` | `New(store, brokerPub, cfg)`: wires cache+router+shards+triggers+bus; `Start/Drain(ctx)` (graceful: stop intake, drain shards bounded by timeout, flush batches, §5.3) | 140 |
| 6.15 | `internal/engine/*_test.go` (M6b) + `internal/testutil/astartedevice.go` | Test device: paho client speaking Astarte MQTT v1 (BSON via mongo-driver, JSON profile, introspection, control payloads) — reused by M7/M9 | 700 |

### 7.3 Verification

- **T1:** ordering property test — N messages for one device through a multi-shard engine with a synthetic slow store arrive persisted in publish order; cross-device parallelism observed. Backpressure: filled shard blocks QoS 1 submit / drops QoS 0 + metric. Reject-reason table covering every §2.6 failure (uses M1 fixtures). zlib: golden `producer/properties` blob round-trip; zip-bomb (1 GiB declared / absolute ceiling) rejected. Trigger matching: fixture upstream trigger JSONs → match/no-match table; webhook retry on 500-then-200 (httptest); event JSON golden vs upstream SimpleEvent shape.
- **T2:** pipeline → real Timescale: datastream rows land in correct typed column; property upsert/unset; `Ack` ordering vs commit asserted with an instrumented store (no ack before commit; injected commit failure ⇒ no ack, retry path).
- **T3 (broker+engine+DB):** test device connects (M4 certs) → introspection → publishes individual BSON, individual JSON, object-aggregated BSON → rows correct; property set/unset; `emptyCache` ⇒ device receives server-owned properties + purge message (decode/inflate asserted); `producer/properties` purge removes the right rows; LISTEN/NOTIFY: interface installed mid-run becomes publishable without restart; cache snapshot swap is race-free under `-race`.
- **CP-B (T4):** official Go SDK end-to-end loop (§0.3) against composed Astrate; asserts SDK-side callbacks fire for server-owned data and Astrate-side rows match.

**M6 gate:** T1–T3 green under `-race`; CP-B green.

---

## 8. M7 — REST Surfaces: `internal/realm`, `internal/housekeeping`, `internal/appengine`

**Goal:** the operator-facing APIs (§3.7), wire-shaped to upstream. Two Phase 3 sessions:
**M7a** (realm + housekeeping) and **M7b** (appengine + stream socket).

### 8.1 M7a — realm management + housekeeping

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 7.1 | `internal/realm/service.go` | Interfaces: install (parse+validate via `pkg/interfaceschema`), update (CheckMinorUpgrade enforcement), delete (upstream draining rules: only major 0 / not in any introspection — parity), list/get; triggers CRUD (validate action shape); `config/auth` JWT key rotation (§4.2); every mutation → store NOTIFY + in-process engine callback | 220 |
| 7.2 | `internal/realm/http.go` | `/realmmanagement/v1/{realm}/interfaces[...]`, `/triggers[...]`, `/config/auth` — `a_rma` | 160 |
| 7.3 | `internal/housekeeping/service.go` | Realm create (tx: row + `ca.NewRealmCA` + sealed key + JWT keys + registration limit), delete (cascade), list/get; instance-admin key auth (`a_ha`) | 140 |
| 7.4 | `internal/housekeeping/http.go` | `/housekeeping/v1/realms[...]` | 100 |
| 7.5 | `internal/realm/*_test.go`, `internal/housekeeping/*_test.go` | T2 HTTP suites | 400 |

### 8.2 M7b — appengine + live stream

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 7.6 | `internal/appengine/service.go` | Device list (paged)/get (upstream status body: introspection, connected, stats…), PATCH (aliases, attributes, `credentials_inhibited`), devices-by-alias; groups CRUD + membership | 240 |
| 7.7 | `internal/appengine/data.go` | GET interface data: properties snapshot tree; datastream queries (`since/since_after/to/limit/downsample_to` → store Series/Downsample; value re-encoding per stored `value_type` — longinteger as string, binaryblob base64, datetime RFC 3339, §2.3); PUT/POST server-owned (→ `engine.ServerData`); DELETE property unset | 260 |
| 7.8 | `internal/appengine/http.go` | `/appengine/v1/{realm}/devices...` routes per §3.7 — `a_aea` | 180 |
| 7.9 | `internal/appengine/stream/ws.go` | `/astrate/v1/{realm}/socket` WebSocket + SSE fallback fed by `engine/stream.Bus`; `a_ch` claims honoured as room filters (§1.1 deviation: not Phoenix-wire) | 180 |
| 7.10 | `internal/appengine/*_test.go` | T2/T3 suites | 550 |

### 8.3 Verification

- **T2 golden suites per service:** envelopes + status codes byte-compared (404 device, 409/422 interface conflicts, 401/403 split). Interfaces: minor-bump additive accepted, mapping-mutation rejected (CheckMinorUpgrade), major coexistence, delete-while-introspected rejected. Housekeeping: created realm immediately serves pairing (cross-domain test with M4). AppEngine: datastream query boundaries (`since` inclusive vs `since_after` exclusive), `limit`, descending default ordering parity, `downsample_to` bucket-count correctness; longinteger > 2^53 round-trips as string through publish→query; property tree shape golden.
- **T3 cross-domain:** PATCH `credentials_inhibited=true` ⇒ broker rejects next CONNECT and pairing rejects credentials. Server-owned PUT ⇒ connected test device receives it (correct QoS + format); offline device receives on reconnect (retention path). Stream socket: test device publishes ⇒ WebSocket subscriber receives the event JSON.
- **CP-C (T4):** scripted `astartectl` run (pinned binary): housekeeping create realm → install `org.astarte-platform.genericsensors.Values` → register device → test device publishes → `astartectl appengine devices data ...` returns the values; trigger CRUD round-trip.

**M7 gate:** T2/T3 green; CP-C green.

---

## 9. M8 — Config, Binary Wiring, Observability, Packaging

**Goal:** one process, one config, ops-ready (§5).

| # | File | Contents | Est. LOC |
|---|---|---|---|
| 8.1 | `internal/config/config.go` | TOML + `ASTRATE_*` env overrides (§5.1): listeners, DSN, shard count/depth, batch knobs, chunk/retention knobs, cert TTL, master-key ref, `insecure_dev_mode`, `enforce_latest_cert`, auto-provision realm block; validation + zero-config defaults; `config.example.toml` generated alongside | 220 |
| 8.2 | `internal/observability/metrics.go` | Prometheus registry + the §5.2 metric set (ingest rate, per-reason rejects, shard depth, flush latency, sessions, pool stats, trigger outcomes); interfaces consumed by engine/broker defined here, satisfied here | 160 |
| 8.3 | `internal/observability/health.go` | `/astrate/v1/health` (liveness), `/astrate/v1/readiness` (DB ping + broker listener check), `/astrate/v1/metrics` | 80 |
| 8.4 | `cmd/astrate/main.go` | `run(ctx, cfg) error`: store→migrations→auth→pairing→engine→broker→HTTP mux (all base paths §3.7)→auto-provision realm if configured; signal-driven graceful shutdown **ordering**: HTTP drain → broker stop-accepting → engine drain/flush (§5.3); `slog` JSON setup, per-domain levels | 240 |
| 8.5 | `cmd/astrate/main_test.go` | T3 boot suite (drives `run()` in-process) | 250 |
| 8.6 | `Dockerfile` | Two-stage: build (CGO_ENABLED=0, `-trimpath -ldflags=-s -w`) → distroless static, non-root, read-only FS + session-store volume (§4.5/§5.4) | 40 |
| 8.7 | `docker-compose.yml` (final) | `astrate` + `timescaledb`, volumes (pgdata, sessions), healthchecks, both TLS-in-binary and reverse-proxy profiles documented (§4.5) | 60 |
| 8.8 | `README.md` (final) + `docs/OPERATIONS.md` | Quickstart (compose + bare-VPS), config reference pointer, backup notes, CA re-keying runbook (§6) | — |

**Verification:**
- **T1:** config precedence (default < TOML < env), validation rejections, example file parses.
- **T3:** full boot on a fresh container DB → readiness 200, auto-provisioned realm serves a pairing registration; SIGTERM mid-traffic: test device publishing QoS 1 during shutdown ⇒ zero acked-but-unpersisted messages (count rows vs PUBACKs); shutdown completes < drain timeout; bbolt sessions survive restart at the binary level (re-runs M5's restart test through `run()`).
- **CI compose smoke:** `docker compose up -d` → health 200 → one full device loop (testutil device) → `docker compose down`; image size budget asserted ≤ 30 MB.
- Metrics endpoint exposes every §5.2 series (presence test).

**M8 gate:** compose smoke green in CI; binary is the single deployable artifact.

---

## 10. M9 — Conformance Matrix, Load/Footprint, Hardening, Release Readiness

**Goal:** prove the §0 hard-compatibility goal with unmodified official clients, enforce the
resource budget, close the §4.5 checklist. All harness code lives under `test/` (own go.mod to
keep SDK deps out of the main module).

| # | File | Contents |
|---|---|---|
| 9.1 | `test/conformance/docker-compose.conformance.yml` | Astrate + Timescale + runner containers |
| 9.2 | `test/conformance/gosdk/main.go` | Official Go SDK runner: register → connect → individual+object datastreams → properties set/unset → receive server-owned property+datastream → `emptyCache` recovery → disconnect/reconnect `session_present` → cert renewal near-expiry (short TTL config) |
| 9.3 | `test/conformance/pysdk/runner.py` (+ pinned `requirements.txt`) | Same scenario matrix on the official Python SDK |
| 9.4 | `test/conformance/atomvm/main.go` | **AtomVM-profile JSON simulator** (§3.5, §6): MQTT 3.1.1 client constrained to the documented JSON profile + zlib control payloads — registration with `initial_payload_format: "json"`, JSON datastreams/properties, asserts server-owned data arrives JSON-encoded |
| 9.5 | `test/conformance/astartectl/run.sh` | Pinned-binary CLI regression (CP-C superset incl. groups, aliases) |
| 9.6 | `test/conformance/verify/main.go` | Cross-checks: DB rows vs published values for every runner; protocol byte traces (mochi hook tap) archived as CI artifacts on failure |
| 9.7 | `test/load/main.go` | 1000 simulated devices @ 1 msg/s for 15 min against the compose stack |
| 9.8 | `test/security/*_test.go` | Zip-bomb on `producer/properties`; pairing brute-force rate-limit; authz sweep (every mux route 401s without token — route-table-driven); TLS config assertions (min version, no client-auth bypass on :8883); oversize topic/payload/introspection bounds (§4.5) |
| 9.9 | `docs/JSON-PAYLOAD-PROFILE.md` | The normative §3.5.3 profile spec for AtomVM client authors |
| 9.10 | `docs/COMPATIBILITY.md` | Deviation inventory (§3.6) + supported-SDK matrix with pinned versions tested |
| 9.11 | `.github/workflows/ci.yml` (final) | `conformance` job on PRs to main; `nightly`: conformance + load + security |

**Verification / gate (CP-D + T5):**
- All four conformance runners green; failure artifacts (byte traces) wired.
- Load budgets **asserted, not observed**: Astrate RSS ≤ 150 MB and Postgres ≤ 768 MB at steady state (§5.4); ingest p99 end-to-end (publish→committed) < 250 ms; shard-depth metric never saturates; zero validation rejects for well-formed traffic.
- Compression verified live: after forcing `compress_chunk` on aged data, `before/after` size ratio recorded (informational) and queries still correct.
- Security suite green; `gosec` high findings zero.
- **Release readiness:** tag `v0.1.0` criteria = CP-D + T5 + docs 9.9/9.10 complete.

---

## 11. Phase 3 Session Plan

Twelve generation sessions, each ending compile-clean with its gate green:

| Session | Scope | Est. new LOC (code + tests) |
|---|---|---|
| S1 | M0 | ~400 |
| S2 | M1.1–M1.2 + M1.3 (`pkg/deviceid`, `astarteapi`, `interfaceschema`) | ~1,800 |
| S3 | M1.4 (`pkg/payload`) | ~1,150 |
| S4 | M2 (migrations + store) | ~2,400 |
| S5 | M3 + M4 (`auth`, `pairing`+CA) → CP-A | ~2,250 |
| S6 | M5 (`broker`) | ~1,600 |
| S7 | M6a (pipeline + data path) | ~1,500 |
| S8 | M6b (control/serverdata/triggers/stream) → CP-B | ~1,950 |
| S9 | M7a (`realm` + `housekeeping`) | ~1,050 |
| S10 | M7b (`appengine` + stream) → CP-C | ~1,400 |
| S11 | M8 (config/binary/packaging) | ~1,050 |
| S12 | M9 (conformance/load/security/docs) → CP-D | ~1,800 |

Total ≈ 18.4 k LOC (≈ 55 % tests/harness). Estimates guide session chunking only — gates, not
line counts, define done.

---

## 12. Standing Risks Tracked Through the Roadmap

| Risk (§6) | Where it is retired |
|---|---|
| Protocol drift vs. real SDKs | CP-A (M4), CP-B (M6), CP-C (M7), CP-D (M9) — drift caught at the earliest milestone that could introduce it |
| mochi session semantics under restart | M5 dedicated restart/replay/QoS 2 tests; re-run at binary level in M8 |
| TTL-by-DELETE cost | M2 TTL-job test + M9 load run measures job duration under telemetry volume; revisit clause stays open until T5 data exists |
| CA key compromise surface | M2 crypto tests, M4 CA tests, M9 security suite + re-key runbook (docs/OPERATIONS.md) |
| Single-process blast radius | M6 per-shard panic-recovery tests; `-race` enforced on T1–T3 from M6 onward |
