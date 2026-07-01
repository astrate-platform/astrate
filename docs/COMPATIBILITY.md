# Astrate Compatibility

Astrate is a wire-compatible re-implementation of Astarte: unmodified official
Astarte device SDKs and `astartectl` work against it. This document inventories
the **deliberate** deviations (`docs/DESIGN.md` §3.6) and lists the client
versions the conformance suite (`test/conformance`) pins and exercises.

## Supported / tested clients

These are pinned in `test/conformance/go.mod`, `test/conformance/pysdk/requirements.txt`,
and the astartectl helper; they are upgraded deliberately, never by drift
(`docs/ROADMAP.md` §0.3).

| Client | Pinned version | Checkpoint(s) |
|---|---|---|
| `astarte-device-sdk-go` | v0.90.2 | CP-B (M6), CP-D (M9 `gosdk`) |
| `astarte-go` (pairing/agent client) | v0.90.4 | CP-A, CP-B |
| `astartectl` (release binary) | v26.5.0 | CP-A, CP-C, CP-D regression |
| `astarte-device-sdk-python` | pinned in `pysdk/requirements.txt` | CP-D (`pysdk`) |
| AtomVM JSON profile (Astrate simulator) | n/a (this repo) | CP-D (`atomvm`) |

The conformance checkpoints gate the build at the earliest milestone that could
introduce protocol drift: **CP-A** (pairing, M4), **CP-B** (full device loop,
M6), **CP-C** (`astartectl` operator flow, M7), **CP-D** (full matrix, M9).

## Wire-identical surfaces

Guarded by the conformance suite against the official clients:

- MQTT topics, Astarte MQTT v1 connection contract (mTLS, identity from the
  certificate CN with the wire client ID free-form and remapped to the CN —
  the VerneMQ subscriber-id remap; session handling), and ACL model. The
  official Python SDK connects with a random paho client ID and relies on
  this.
- BSON `{v, t}` data documents and the zlib + size-prefixed control payloads
  (`emptyCache`, `producer/properties`, `consumer/properties`).
- Introspection format (`;`-separated `name:major:minor` triples).
- Pairing REST bodies and status codes (registration, credentials, info,
  verify).
- Certificate `Subject CN = <realm>/<device_id>`; 128-bit serials; clientAuth.
- JWT claim model (`a_aea`, `a_rma`, `a_pa`, `a_ha`, `a_ch`) with implicit
  anchoring and the `"<verb-regex>::<path-regex>"` authorization strings.
- Realm-management interface/trigger install/update/delete semantics and
  AppEngine device/data/query shapes.

## Deliberate deviations

All additive or strictly-safer; none affect unmodified device SDKs.

1. **Astarte Channels** — the Phoenix-socket protocol is replaced by a plain
   WebSocket/SSE endpoint at `/astrate/v1/<realm>/socket`. Different contract,
   additive endpoint, honouring `a_ch` claims as room filters. Device SDKs are
   unaffected (Channels is a consumer-side API).

2. **JSON payload profile + `initial_payload_format`** — Astrate accepts a
   documented plain-JSON data encoding alongside BSON on the same topics, and an
   additive registration field selects the device's server→device format. Pure
   superset; see [JSON-PAYLOAD-PROFILE.md](JSON-PAYLOAD-PROFILE.md). Upstream SDK
   behaviour is byte-identical.

3. **MQTT 5.0 accepted** — Astarte uses MQTT 3.1.1; Astrate's broker also
   accepts 5.0 clients (mochi default). A superset with no SDK impact.

4. **Astrate-native endpoints under `/astrate/v1/...`** — health, readiness,
   metrics, and the live-stream socket live in a namespace that cannot collide
   with the upstream API surface.

5. **Uniform pairing-credentials error** — a wrong credentials secret and an
   unknown device both return a uniform `401` (upstream returns `403` via its
   RPC "forbidden" mapping). Stricter/safer: it avoids a device-enumeration
   oracle and matches the per-IP/per-device rate-limited timing.

6. **Housekeeping realm body** — `{realm_name, jwt_public_key_pem,
   device_registration_limit}`; the Cassandra-specific replication fields
   (`replication_class`, `replication_factor`, …) upstream carries are accepted
   but ignored (Astrate is PostgreSQL/TimescaleDB, not Cassandra).

7. **Latest-serial enforcement** (`pairing.enforce_latest_cert`) — when enabled,
   the broker rejects a certificate whose serial differs from the device's
   latest issuance (an always-online CRL equivalent), stricter than upstream's
   CRL-less default. Off by default for fleets that rotate while devices hold
   older still-valid certs.

8. **AppEngine device-list pagination** — the next-page cursor is returned in
   the `X-Astrate-Next-Token` response header (with `data` = the id array),
   rather than an upstream `links.next` body field. Single-page listings are
   wire-identical; pagers that follow the header continue seamlessly.

## Infrastructure differences (by design)

Not protocol deviations — these are the point of the project (`docs/DESIGN.md`):

- **Single Go binary** (modular monolith) instead of fragmented Elixir
  microservices; no Kubernetes.
- **PostgreSQL + TimescaleDB** instead of Cassandra/ScyllaDB; a shared
  hypertable with typed columns instead of per-interface tables.
- **Embedded mochi-mqtt broker** and an **embedded per-realm CA** instead of
  VerneMQ + CFSSL.
- **In-process sharded pipeline** instead of RabbitMQ between the broker and the
  persistence layer.
