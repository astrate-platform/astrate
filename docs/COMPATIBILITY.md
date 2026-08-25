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
- AppEngine device PATCH contract: honored only with `Content-Type:
  application/merge-patch+json` (anything else reproduces upstream's unmapped
  `:patch_mimetype_not_supported` fallback → 500), plus the alias/attribute
  error taxonomy (`Invalid alias`, `Alias already in use`, `Alias tag not
  found`, `Invalid attributes`, `Attribute key not found`).
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
   - **The `a_ch` grammar is measured, and it is not the REST grammar.**
     Recorded against upstream v1.2.0 in
     `test/conformance/upstream/channels.json`. Upstream partitions the `a_ch`
     list by an **exact** match on the verb field, keeping only entries whose
     first field is literally `JOIN` or `WATCH` and discarding every other
     entry; the path regex is then matched within the chosen bucket. So the
     astartectl-style blanket `.*::.*`, which authorizes every REST surface,
     authorizes **nothing** on Channels — `.*` is not the string `JOIN`. A join
     is checked against the room name with `rooms:<realm>:` stripped off, and
     `Token.AuthorizesChannel` implements this separately from `Authorizes`,
     which stays a verb-regex match because that is what upstream's REST plug
     does. Astrate previously read `a_ch` with the REST rule and did not check
     `JOIN` on a join at all; both were corrected in M12 phase 06.
   - **The WATCH authorization paths are measured**, with narrow claims that
     make exactly one of each pair acceptable
     (`test/conformance/upstream/channels.json`). A data trigger is authorized
     against `<device_id>/<interface_name><match_path>` — the match path is
     appended, carrying its own leading slash — and a device trigger against
     the bare `<device_id>`. Astrate built `<device_id>/<interface_name>` for
     data triggers until M12 phase 06b, which was wrong in *both* directions:
     it refused claims upstream accepts and accepted claims upstream refuses.
     The group shapes (`groups/<name>/<interface><match_path>` and
     `groups/<name>`) are measured too (2026-08-22 recording, realm with a
     group): both are accepted exactly as Astrate builds them. The recording
     also settled where `group_name` lives — at the watch payload's **top
     level**; a `group_name` nested inside `simple_trigger` is refused by
     upstream's changeset before authorization is consulted — and that a group
     device_trigger must carry `device_id: "*"` inside `simple_trigger`
     (refusal reason: `device_id must be * for group triggers`, the mirror
     image of plain device triggers, where `"*"` is refused). Astrate read
     `group_name` only from `simple_trigger` until the 2026-08-22
     reconciliation, which means an upstream-shaped group watch silently
     degraded into a device-shaped path check; fixed.
   - **A device trigger's `device_id` must sit inside `simple_trigger` and
     equal the request's own.** Upstream refuses a watch whose `device_id` is
     only at the payload's top level — where the AppEngine REST API puts it —
     and refuses a wildcard `device_id: "*"`, both with the misleading reason
     `unauthorized`. Astrate used to fall back to the top-level `device_id` and
     accept both; since M12 phase 06b it refuses them with the same reason.
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
   `interface_loading_failed` rather than an unaccepted name. The mapping table
   is **measured**, not reconstructed: it was recorded against upstream Astarte
   v1.2.0 in `test/conformance/upstream/channels.json`, and two rows were
   corrected as a result. The last two rows were measured on 2026-08-22:
   `write_on_server_owned_interface` is exact — upstream emits precisely that
   name for a device publishing on a server-owned interface it has declared.
   `value_size_exceeded` is real in upstream's value validator (astarte_core
   rejects strings over 65536 bytes with that name) but **unreachable over
   MQTT on the self-hosted stack**: the transport discards any publish whose
   packet exceeds 65536 bytes before Data Updater Plant ever sees it — broker
   ACKs, nothing stored, no event, no session change — and closes the
   connection outright above ~3 MB (measured: 65468-byte BSON payload
   delivered end-to-end; 65473 bytes silently dropped). Astrate instead
   enforces its own caps (`MaxPayloadBytes` 64 KiB default,
   `MaxStringLen` 65536) and emits a `device_error` naming `value_size_exceeded`
   where upstream says nothing at all — kept as a deliberate divergence, since
   telling the device is strictly more useful than a silent discard, and the
   name itself is upstream's own enum value.

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
   - **Unknown realms** — a request to a realm that does not exist is the
     same case as an unverifiable token (an unknown realm has no key to
     verify against), so Astrate answers `401` uniformly. Measured against
     upstream 1.2.0 on 2026-08-24 (`test/conformance/upstream/verify-versions.json`,
     issue #69): upstream answers `403 Forbidden` for any non-absent
     credential on Realm Management and AppEngine unknown-realm paths — the
     same mapping that produces the bullet above — while a missing token
     still gets `401`. Astrate's uniform `401` is this same deviation applied
     one step further: it does not distinguish "no credentials" from
     "credentials we cannot verify", and an unknown realm is never
     distinguishable from a wrong key.

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

11. **AppEngine data query: `sort=ascending`** — additive extension. Upstream
    1.2.2 serves time series newest-first and has no `sort` parameter; Astrate
    additionally accepts `sort=ascending` (any other value keeps the default
    descending order). Standard clients never send it and are unaffected.

12. **Group-listing pagination token** — same wire shape
    (`?limit&from_token`, next-page cursor in `links.next`), but the opaque
    token encodes a row offset instead of upstream's `insertion_uuid` keyset.
    Tokens are server-generated and never valid across servers, so no
    client-visible difference; adopting the keyset needs an insertion-uuid
    column (migration decision, deferred).

13. **AMQP trigger actions are rejected at creation** — upstream 1.2.2 accepts
    trigger actions with `amqp_exchange`/`amqp_routing_key` and forwards
    events to RabbitMQ; Astrate has no AMQP bus, so such actions fail trigger
    installation with a clear per-field error instead of silently dropping
    events later. Legacy stored AMQP triggers fail loudly at reload rather
    than forwarding nothing.

14. **AppEngine server-write taxonomy, measured with two clean-rejection
    divergences** — the write-path statuses are measured against upstream
    1.2.0 (`test/conformance/upstream/verify-server-writes.json`, 2026-08-24,
    issue #57) and matched: a server write to a device-owned interface is
    `405 Cannot write to device owned resource`; an unknown interface is
    `404 Interface not found in device introspection` (the device's
    introspection is consulted first — an installed-but-undeclared interface
    still answers there); an unknown endpoint is `400 Endpoint not found`;
    a value over 64 KiB is `422 Value size exceeds size limits`; a DELETE on
    a datastream is refused, and an unset of an absent property succeeds
    offline with `204`. Two upstream rows are internal errors Astrate does
    not reproduce:

    - **A value of the wrong scalar type is `500` upstream** (unmapped
      clause in its write path; same for a malformed object aggregate).
      Astrate keeps the `500` status bug-for-bug here — upstream wins on
      wire — so a client sees the same status it would see from upstream;
      only the envelope wording differs (deviation 15).
    - **DELETE on a server-owned *individual* datastream is `500` upstream**
      while the object-aggregated sibling gets `405 Cannot write to
      read-only resource` — clearly the same clause one branch above it.
      Astrate answers the `405 read-only resource` form for both.

    One more measured fact worth keeping: **the REST write path has no
    offline push failure.** A server-owned datastream write for an offline
    device is stored and answered `200` — the issue-guessed
    `503 cannot_push_to_device` is not observable via AppEngine REST on
    1.2.0, and Astrate's persist-then-publish behaves the same way.

15. **Canonical detail casing follows measurement, not memory** — upstream's
    Phoenix renders `Bad request` and `Internal server error`; the
    reconstructed `Bad Request` / `Internal Server Error` constants were
    corrected to the measured forms on 2026-08-24. Every probed row used the
    lowercase variants.

16. **`virtual_device_pool` publishes without an MQTT session** (measured
    against upstream Astarte Flow master @ ad0ae81, 2026-08-24). Upstream's
    dynamic pool registers the first-seen id through Pairing, then spawns a
    *real* MQTT device: client certificate, broker connection, introspection
    from config, ordinary datastream publishes; its local DETS holds the
    credentials secret, and losing `FLOW_PERSISTENCY_DIR` permanently bricks
    every id whose certificate was already issued (`422 already_registered`
    forever — there is no recovery path). Astrate keeps the observable
    contract — key grammar `<device_id>/<interface></path…>` (upstream's
    leading realm segment is dropped: flows are per-realm), first-seen
    registration through the pairing door, rows queryable like any
    device-owned datastream, non-201 registrations logged and dropped with
    the message — but lands values through the engine ingest path instead of
    a per-device MQTT session, skips the introspection gate (virtual devices
    never connect, so there is no introspection to match), and keeps secrets
    server-side, so the brick-on-store-loss failure mode does not exist here.
    Re-registration parity is preserved at the semantic level: upstream
    answers `201` with a fresh secret until first credential issuance and
    `422` afterwards; Astrate's pairing door accepts re-registration of a
    not-yet-confirmed device (the minted secret is discarded) and refuses an
    already-confirmed one (`ErrAlreadyRegistered`, mapped by the block to a
    log-and-drop).

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
