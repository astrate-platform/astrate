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

### 11. AppEngine data query: `sort=ascending`

Additive. Upstream 1.2.2 serves time series newest-first and has no `sort` parameter; Astrate also accepts `sort=ascending`. Standard clients never send it.

### 12. Group-listing pagination token

Same wire shape (`?limit&from_token`, cursor in `links.next`), but the opaque token encodes a row offset instead of upstream's `insertion_uuid` keyset. Tokens are server-generated and never portable across servers, so there is no client-visible difference.

### 13. AMQP trigger actions are rejected at creation

Upstream 1.2.2 accepts trigger actions with `amqp_exchange`/`amqp_routing_key` and forwards events to RabbitMQ. Astrate has no AMQP bus, so it fails trigger installation with a per-field error instead of silently dropping events later; stored legacy AMQP triggers fail loudly at reload.

### 14. AppEngine server-write taxonomy

The write-path statuses are measured against upstream 1.2.0 and matched, with two deliberate divergences where upstream returns an unhandled `500`:

- A wrong scalar type (and a malformed object aggregate) keeps upstream's `500` bug-for-bug; only the envelope wording differs.
- `DELETE` on a server-owned *individual* datastream is `500` upstream while its object-aggregated sibling is `405 Cannot write to read-only resource`. Astrate answers the `405` form for both.

### 15. Canonical detail casing follows measurement

Upstream's Phoenix renders `Bad request` and `Internal server error`; Astrate uses the measured lowercase forms rather than the reconstructed title-case ones.

### 16. `virtual_device_pool` publishes without an MQTT session

Upstream's dynamic pool registers each first-seen id through Pairing and then spawns a real MQTT device, keeping the credentials secret in a local store -- losing that store permanently bricks every id whose certificate was issued, with no recovery path. Astrate keeps the observable contract (key grammar, first-seen registration through the pairing door, rows queryable like any device-owned datastream) but lands values through the engine ingest path and keeps secrets server-side, so that failure mode does not exist.

### 17. Always synchronous where upstream 1.4 defaults to asynchronous

Upstream 1.4 runs realm create/delete, interface install/update/delete and delivery-policy delete in the background, letting the caller opt into synchronous execution with `?async_operation=false`. Astrate performs all of them synchronously and answers only once the work is done -- a strictly stronger guarantee. The parameter is accepted and ignored on either value, so upstream clients that send it keep working unchanged.

### 18. Per-service health endpoints answer 503 when a dependency is down

Upstream's unauthenticated `GET /{appengine,realmmanagement,pairing}/health`, which the Astarte Dashboard polls for its per-service status indicators, returns a static `200` -- the indicator cannot go red whatever state the instance is in. Astrate serves the same route and the same `200` envelope, but runs its database probe first and answers `503` when it fails. Astrate additionally keeps the realm-scoped `GET /pairing/v1/<realm>/health`, which upstream 404s: it resolves the realm too, and devices can probe it before they hold credentials.

## Infrastructure differences (by design)

- Single Go binary instead of Elixir microservices; no Kubernetes.
- PostgreSQL + TimescaleDB instead of Cassandra/ScyllaDB.
- Embedded mochi-mqtt broker and CA instead of VerneMQ + CFSSL.
- In-process sharded pipeline instead of RabbitMQ.
