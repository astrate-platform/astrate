# Compatibility

Astrate is wire-compatible with Astarte. This page inventories deliberate deviations and lists tested client versions.

## Tested clients

| Client | Checkpoint |
|---|---|
| `astarte-device-sdk-go` | CP-B (M6), CP-D (M9) |
| `astarte-go` (pairing/agent) | CP-A, CP-B |
| `astartectl` (release binary) | CP-A, CP-C, CP-D |
| `astarte-device-sdk-python` | CP-D |
| AtomVM JSON profile simulator | CP-D |
| Astarte Dashboard v1.2.2 | Runs unmodified |

## Wire-identical surfaces

All guarded by the conformance suite against official clients:

- MQTT topics, Astarte MQTT v1 connection contract (mTLS, CN identity, session handling, ACL model)
- BSON `{v, t}` data documents and zlib + size-prefixed control payloads
- Introspection format (`;`-separated `name:major:minor` triples)
- Pairing REST bodies and status codes
- Certificate `Subject CN = <realm>/<device_id>` with 128-bit serials
- JWT claim model (`a_aea`, `a_rma`, `a_pa`, `a_ha`, `a_ch`)
- Realm management and AppEngine API shapes, including pagination
- Astarte Dashboard v1.2.2 Device Live Events (since M11 Channels socket)

## Deliberate deviations

All additive or strictly-safer; none affect unmodified device SDKs.

### 1. Astarte Channels: two sockets, one bus

The upstream Phoenix socket is served at `/appengine/v1/socket/websocket` (phoenix.js V2 wire format) for Dashboard compatibility. Astrate keeps its own plain WebSocket/SSE endpoint at `/astrate/v1/<realm>/socket`.

- Authentication from query parameters (phoenix.js cannot set `Authorization` header).
- WATCH authorization paths reconstructed from upstream documentation.
- Transient triggers are matcher-only (slow viewers drop frames, not backpressure).

### 2. JSON payload profile

Astrate accepts plain JSON alongside BSON on the same topics. Pure superset; see [Payload Formats](payload-formats.md).

### 3. MQTT 5.0 accepted

Astarte uses MQTT 3.1.1; Astrate also accepts 5.0. Superset, no SDK impact.

### 4. Astrate-native endpoints

Health, readiness, metrics, and the live-stream socket under `/astrate/v1/...` -- a namespace that cannot collide with upstream.

### 5. Uniform 401 vs 403

Astrate returns `401` for every failure to establish identity (wrong secret, unknown device, bad token) and `403` only for authorization refusal after identity is established. Stricter/safer -- eliminates enumeration oracles.

### 6. Housekeeping realm body

Cassandra-specific replication fields accepted but ignored (Astrate is PostgreSQL).

### 7. Latest-serial enforcement

When `pairing.enforce_latest_cert` is enabled, the broker rejects certificates whose serial differs from the device's latest issuance. Stricter than upstream's CRL-less default. Off by default.

### 8. Trigger delivery policies

Astrate enforces trigger delivery policies with its own reading of transport failures. A request that never produced a response (connection refused, DNS failure, timeout) is treated as a server error.

### 9. Synchronous device deletion

`DELETE /realmmanagement/v1/<realm>/devices/<id>` removes the device and all data in one transaction. Upstream runs asynchronous multi-service wipe.

### 10. Emulated API version

`GET /v1/<realm>/version` reports the emulated upstream API level (`1.2.2`), not Astrate's own version.

## Infrastructure differences (by design)

- Single Go binary instead of Elixir microservices; no Kubernetes.
- PostgreSQL + TimescaleDB instead of Cassandra/ScyllaDB.
- Embedded mochi-mqtt broker and CA instead of VerneMQ + CFSSL.
- In-process sharded pipeline instead of RabbitMQ.
