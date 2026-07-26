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
  AppEngine device/data/query shapes, including body-`links` device-list
  pagination (`?details=true&limit&from_token` with the cursor in
  `links.next`) and `/stats/devices`.
- The Astarte Dashboard v1.2.2 runs unmodified against Astrate (compose
  `full` profile, `http://localhost:4040`), Device Live Events included since
  the M11 Channels socket (deviation 1).

## Deliberate deviations

All additive or strictly-safer; none affect unmodified device SDKs.

1. **Astarte Channels: two sockets, one bus** — the upstream Phoenix socket is
   served at `/appengine/v1/socket/websocket` (phoenix.js V2 wire format,
   `?vsn=2.0.0&realm=&token=`), which is what the Dashboard's Device Live
   Events card speaks; `phx_join`, `watch`, `phx_leave` and the `phoenix`
   heartbeat are answered, and matching events are pushed as `new_event`.
   Alongside it Astrate keeps its own plain WebSocket/SSE endpoint at
   `/astrate/v1/<realm>/socket` (additive, honouring `a_ch` claims as room
   filters). Device SDKs are unaffected either way — Channels is a
   consumer-side API. Three points where the implementation is a reading
   rather than verified parity:

   - **Authentication happens before the upgrade, from query parameters.**
     phoenix.js cannot set an `Authorization` header, so `realm` and `token`
     come from the query string and are verified on the pre-upgrade response.
     Every authentication failure — missing or unknown realm, missing or
     unverifiable token — is a uniform `401` with no existence oracle;
     authorization refusals are `error` reply frames after the upgrade, not
     HTTP status codes.
   - **The WATCH authorization paths are reconstructed from upstream
     documentation, not from upstream source**: `groups/<name>` for a group
     trigger, `<device_id>/<interface>` for a data trigger, and `<device_id>`
     otherwise, each checked as `Authorizes(a_ch, "WATCH", …)`. astartectl-style
     `.*::.*` grants behave identically under any reading; narrowly-scoped
     `a_ch` claims are where a divergence would show.
   - **Transient triggers are matcher-only.** A `watch` payload's
     `simple_trigger` is compiled to a matcher with no action, never stored and
     never handed to the trigger executor. A slow viewer drops frames rather
     than backpressuring ingestion, exactly as the bus itself does.

   One further point, where the constraint is upstream's and Astrate conforms to
   it rather than deviating: **`device_error.error_name` is drawn from a closed
   enum.** Consumers validate it — the Dashboard's card discards the entire
   event and logs "Unrecognised event received" for any name outside the set —
   so Astrate's own reject-reason labels (§2.6), which are a metrics and log
   surface with their own stability requirements, are translated to the upstream
   enum by `triggers.UpstreamErrorName` when an event body is built. The
   original label rides along under `metadata["astrate_reason"]`, so nothing
   diagnostic is lost, and any reason without a specific counterpart becomes
   `interface_loading_failed` rather than an unaccepted name.

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

5. **Uniform `401` where upstream distinguishes with `403`** — Astrate answers
   `401` for every failure to establish *who* is calling, and reserves `403`
   for a caller it did establish who is not permitted the path. Two places
   where that costs a status-code difference:

   - **Pairing credentials** — a wrong credentials secret and an unknown
     device both return `401`; upstream returns `403` via its RPC "forbidden"
     mapping.
   - **Malformed or unverifiable bearer tokens, on every authenticated
     surface** — a token that is not a JWT, one whose signature does not
     verify, and one whose realm key is unknown all return `401` from
     `auth.Middleware`. Upstream answers `403` here too. Observed against
     `api.eu1.astarte.cloud` on 2026-07-26 with no valid credentials of any
     kind, so it is a recorded upstream response rather than a reading of its
     source; the sample is AppEngine, Realm Management and Pairing, and no
     Astarte version other than the one that instance was running.

   Stricter/safer, and deliberately uniform: distinguishing "no such device"
   or "no such realm key" from "wrong secret" hands an unauthenticated caller
   an enumeration oracle, and the uniform answer matches the per-IP and
   per-device rate-limited timing so the status code does not leak what the
   delay is hiding. `403` still appears exactly where upstream puts it once a
   token *has* verified — an authorization path the token does not cover — so
   a client distinguishing the two codes reads Astrate's `401` as "present
   your credentials again" and its `403` as "these credentials will never do",
   which is the useful distinction and the one upstream blurs.

6. **Housekeeping realm body** — `{realm_name, jwt_public_key_pem,
   device_registration_limit}`; the Cassandra-specific replication fields
   (`replication_class`, `replication_factor`, …) upstream carries are accepted
   but ignored (Astrate is PostgreSQL/TimescaleDB, not Cassandra).

7. **Latest-serial enforcement** (`pairing.enforce_latest_cert`) — when enabled,
   the broker rejects a certificate whose serial differs from the device's
   latest issuance (an always-online CRL equivalent), stricter than upstream's
   CRL-less default. Off by default for fleets that rotate while devices hold
   older still-valid certs.

8. **Trigger delivery policies: enforced, with a stated reading of transport
   failures** — a trigger naming a policy gets that policy's `error_handlers`
   verdict per response status, its `retry_times` as the attempt cap, its
   `event_ttl` applied at dequeue, and its `maximum_capacity` as a per-policy
   in-flight bound. Realm Management rejects a trigger naming an unknown
   policy and a policy delete while triggers still reference it. Four points
   where an operator's upstream mental model may not match:

   - **A request that never produced a response** — connection refused, DNS
     failure, timeout — is treated as a *server error*: it matches `any_error`
     and `server_error` handlers and never an explicit status-code list.
     Upstream's schema speaks only in status codes and says nothing about
     transport failures, so this is Astrate's reading rather than a
     documented upstream behaviour. Every retry and discard logs the rule
     that decided it (`reason`, alongside `policy`), so the applied
     interpretation is visible per delivery rather than inferred.
   - **A status no handler claims is discarded.** Handlers are read as the
     whole specification of what to retry.
   - **A trigger naming no policy** keeps Astrate's fixed default — 4xx
     permanent, 5xx and transport failures retried with backoff, capped by
     `engine.triggers` config — rather than resolving upstream's named
     `@default` policy. The two coincide in practice; the `@default` document
     itself is not installable or inspectable.
   - **`maximum_capacity` is realized as a per-policy in-flight counter**, not
     a dedicated queue per policy. The bound is the same; what differs is
     that deliveries over the cap are dropped at enqueue with the existing
     `dropped` metric rather than queued behind a policy-private buffer.

   Not implemented from the policy schema: handler-overlap rejection (two
   handlers claiming the same status code — the first in document order
   wins rather than the install being refused).

9. **Synchronous device deletion** — `DELETE
   /realmmanagement/v1/<realm>/devices/<id>` kicks the live MQTT session and
   removes the device row plus all its data in one transaction. Upstream runs
   an asynchronous multi-service wipe and exposes a transient
   `deletion_in_progress` device-status flag; Astrate's status never carries
   it because the state cannot be observed. Same endpoint, verbs, status
   codes, and end state.

10. **Emulated Realm Management API version** — `GET /v1/<realm>/version`
    reports `realm.APICompatVersion` (currently `1.2.2`), a compatibility
    declaration of the emulated upstream API level, not Astrate's own release
    version. The Dashboard feature-gates its Policies UI on this value
    (>= 1.1.1).

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
