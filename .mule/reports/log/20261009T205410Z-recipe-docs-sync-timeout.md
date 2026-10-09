slug: recipe-docs-sync
verdict: timeout
at:  1114f61
ran: 2026-10-09T20:54:10Z on DietPi in 1200s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls docs/api/ && echo "---" && ls internal/ && echo "---MKDOCS---" && ls docs/
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
---
appengine
auth
broker
config
engine
flow
flowapi
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil
---MKDOCS---
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
UPSTREAM-EXPERIMENTAL.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui

$ tail -40 .mule/todo.md && echo "===FORGIULIO===" && ls -la .mule/for-giulio.md 2>/dev/null && tail -30 .mule/for-giulio.md 2>/dev/null
- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362) — the second half docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail: Bad Request`, the reconstructed title-case form that `TestHousekeepingErrorDetailExamples` (internal/swagger/swagger_test.go:687) now pins against everywhere else on this spec, while the wire emits the frozen `astarteapi.DetailBadRequest` = `"Bad request"` (pkg/astarteapi/envelope.go:37, via `WriteBadRequest` on every 400 — internal/housekeeping/http.go:65, 166, 171, 183). Because the example sits on the generic `ErrorDetail` schema, which all eight error responses `$ref` (yaml:377-477), it is the example a generator falls back to for any response whose own example it cannot resolve — so it teaches `"Bad Request"` for a 404 and a 500 too. Change it to `Bad request`, or drop it as the realm-management and pairing `ErrorDetail` schemas do (neither carries one). Then extend `TestHousekeepingErrorDetailExamples`, whose comment currently says the schema-level example is "deliberately not pinned" (swagger_test.go:684-686), to pin it like the responses. One thing not to "fix": the example is nested one level in — it is the `errors` sub-schema's example, not a bare envelope — so do not re-wrap it in another `errors:`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-patch-422-field-error-examples [auto]: add the three missing PATCH 422 bodies to the `ValidationErrors` response component in docs/api/astarte_housekeeping_api.yaml (yaml:427-436, referenced once at yaml:190) — it carries a single example, `error_name: [invalid_update_parameters]` (internal/housekeeping/http.go:176-177), but `patchRealm` emits three more field-error envelopes: `jwt_public_key_pem: ["can't be blank"]` (http.go:188-189), `device_registration_limit: ["is invalid"]` (http.go:193-194) and `datastream_maximum_storage_retention: ["is invalid"]` (http.go:198-199). Add them as further `examples` on the media type (OpenAPI 3.0.3, yaml:1) or split the component per case. Also record that the *messages* differ from POST for the same conditions: create rejects a negative limit/retention with the flat `ErrValidation` detail `device_registration_limit must be non-negative` / `datastream_maximum_storage_retention must be non-negative` (internal/housekeeping/service.go:143-148, 422 `ErrorDetail` shape), while PATCH answers `is invalid` in the FieldErrors shape — and that the service's own three `ErrValidation` branches (service.go:205-213) are unreachable from REST, because `patchRealm` pre-checks the identical conditions and answers first. Say in the `RealmPatch` field descriptions that a negative value is a 422 rather than silently ignored. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [!] docs-sync-hk-wrong-type-field-400 [auto]: widen the `BadRequest` component description in docs/api/astarte_housekeeping_api.yaml (yaml:378, "malformed or missing data envelope") — measured today, a perfectly well-formed `data` envelope carrying a wrong-typed field answers **400**, not the documented 422. `astarteapi.DecodeData` into `map[string]json.RawMessage` (internal/housekeeping/http.go:170) cannot fail on a type mismatch, so the failure surfaces on the second decode into the typed body (http.go:182-184) and is answered `WriteBadRequest`. Probe output: `{"data":{"device_registration_limit":"100"}}` → `json: cannot unmarshal string into Go struct field patchBody.device_registration_limit of type int32`; `{"data":{"jwt_public_key_pem":123}}` → the string twin; `{"data":{"device_registration_limit":1.5}}` → `cannot unmarshal number 1.5`; `{"data":[]}` → `cannot unmarshal array into Go value of type patchBody`. Same on POST, which decodes straight into `realmBody` (http.go:63-67). This is unreachable from a spec-conformant client — `RealmCreate`/`RealmPatch` already type both fields — so document it as a one-line extension of the 400 description on both operations rather than inventing a 422, and do NOT re-shape the `BadRequest` example, whose `detail: Bad request` is the frozen constant `TestHousekeepingErrorDetailExamples` pins. One distinction worth stating in the same description: a wrong-typed field is *not* the `invalid_update_parameters` case — unknown keys are caught earlier by the `patchAllowedFields` sweep (http.go:174-180, 422 FieldErrors) on PATCH only, and silently dropped on POST. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `components.parameters.RealmName` and `RealmCreate.realm_name` both carry `pattern: '^[a-z][a-z0-9]*$'` and `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go) holds both against the realm name CHECK extracted from migrations/000002_metadata.up.sql — but the read side is still bare `type: string`: `Realm.realm_name` (components.schemas.Realm) and `DataEnvelopeRealmNames.items`. Decide whether a response field earns the same pattern (a generated client learns nothing new from it, but a reader diffing create against get sees one constrained and one not) and, if it does, add it to both and extend that test to hold the response-side pattern against the same CHECK. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] store-register-inhibit-preserve: RegisterDevice must not overwrite status to 'registered' when a device was inhibited before first credentials request (preserve 'inhibited' state). File: internal/store/devices.go — BLOCKED: gates failed
- [x] store-aliasvalues-self-exclusion-test: Add test for AliasValuesTaken to assert self-exclusion (same device's own alias values do not count as taken). File: internal/store/devices_test.go
- [!] store-latestindividual-empty-errnotfound: Add store-level test asserting LatestIndividual returns ErrNotFound for empty series and returns newest on multiple samples. File: internal/store/datastreams_test.go — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [!] pairing-burn-bcrypt-cost [auto]: make the dummy-secret timing burn use the configured BcryptCost — burnBcrypt compares against a package-level hash generated once at bcrypt.DefaultCost (internal/pairing/service.go:94-100, 428-430) while real auth hashes and compares at cfg.BcryptCost (service.go:204, 407), and bcrypt_cost is a documented, plumbed lever (internal/config/config.go:106, cmd/astrate/main.go:152, docs/site/configuration-reference.md:76, default 10), so with bcrypt_cost < 10 the unknown-device/empty-hash/malformed-ID failure paths each burn cost-10 work against cost-4 real compares, making failure responses ~an order of magnitude slower and a latency oracle for device existence — the inverse of the DESIGN §4.4 uniform-auth-time property. Keep one dummy hash per configured cost (generate lazily per Service at s.cfg.BcryptCost) and burn against it; add a deterministic test asserting bcrypt.Cost of the burn hash equals cfg.BcryptCost (e.g. a Service built with BcryptCost:4 must burn at cost 4), which fails today. — BLOCKED: lint failed: internal/pairing/service.go:95:7: G101: Potential hardcoded credentials (gosec)
- [x] pairing-onregistered-clock: route the registration event timestamp through the injectable clock — Register calls s.OnRegistered(realmName, hwID, time.Now()) with literal wall-clock time (internal/pairing/service.go:222) while every other timestamp in the package funnels through the s.now seam (VerifyCredentials, service.go:333), so a time-travel test cannot pin the DeviceRegisteredEvent timestamp fired via e.HandleDeviceRegistered (cmd/astrate/main.go:154). Change to s.now(), keep TestRegisterEmitsEvent.
- [x] pairing-unregister-bad-id-test: cover the malformed-device-IDs branch of Unregister in internal/pairing/service_test.go — Unregister wraps deviceid.Parse failure as store.ErrNotFound (internal/pairing/service.go:231-234) but TestUnregister (service_test.go:340) only exercises a well-formed unknown ID; assert both bad-ID shapes and the well-formed unknown-ID 404 in one table.
- [x] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.go
- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
- [x] compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (newest stable, 2026-10-05, maintenance-only: mix.exs/mix.lock version bumps plus a single CHANGELOG "astarte_data_updater_plant RPC server reliability" entry, no wire/API surface change; v1.4.0 is still rc.6-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2/v1.3.4 wording proposal there already stands and its version-reference delta was recorded on 2026-10-05 as reading "newest stable v1.3.5", so fold this in as that one-word delta rather than re-deriving the proposal.
- [!] docs-native-socket-handshake-errors [auto]: add the two WebSocket-handshake failures to both socket operations in docs/api/astrate_native_api.yaml — `/astrate/v1/{realm}/socket` responses (yaml:335-345) and `/appengine/v1/socket/websocket` (yaml:398-404) document only 101/401/403/500 and 101/401/500 — but measured today with the repo's own dependency (probe in /tmp, coder/websocket accept.go `verifyClientRequest`): a plain GET with no `Connection: Upgrade` answers **426 Upgrade Required**, `Content-Type: text/plain; charset=utf-8`, body `WebSocket protocol violation: Connection header "" does not contain Upgrade\n`, and `Sec-WebSocket-Version: 14` answers **400** `unsupported WebSocket protocol version (only 13 is supported): "14"\n`. Add both to each operation and say two things explicitly: the bodies are plain text from `http.Error`, NOT the `{"errors":{"detail":...}}` envelope every other response on this spec uses (a client that JSON-parses them breaks), and they are only reached after authentication — native: the route's guard `mw.RequireRealm(auth.ClaimChannels)` wraps `handle` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112), Phoenix: the realm/token checks at internal/appengine/channels/ws.go:51-68 run before `websocket.Accept` at :71 — and only on the WebSocket branch (`wantsSSE` is tested first at ws.go:78), so 401/403 and the 200 SSE path win. Do NOT add 405: ServeMux's method-mismatch applies equally to every GET route in every spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-native-socket-sse-exact-accept [auto]: fix the SSE negotiation prose on `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml — the description (yaml:289-291), the `transport` parameter (yaml:328-334) and the `200` response (yaml:338-339) read as if any client asking for `text/event-stream` gets SSE, but the predicate `wantsSSE` is exact string equality (internal/appengine/stream/ws.go:151-153): measured, `Accept: text/event-stream, text/plain` and `Accept: text/event-stream;q=1.0` both return false and fall through to `websocket.Accept`, where a non-upgrade client gets the 426 above; `transport` selects SSE only for the literal `sse`, and `transport=websocket` / `transport=` (empty) both mean WebSocket. Say "exactly" for the header, say any other `Accept` falls through to the upgrade attempt, and do NOT claim the server rejects out-of-enum `transport` values — the enum is a client-side contract the handler never enforces. While there, give the `200` its real media type: the handler sets `Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive` and writes one `data: {...}` frame per event followed by a blank line (ws.go:123-126, 156-162). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-socket-security-scheme [auto]: `/astrate/v1/{realm}/socket` is documented as if public in docs/api/astrate_native_api.yaml — the spec has no `components.securitySchemes` at all and root `security: []` (yaml:21), and the operation declares no `security`, so a client generated from this spec sends no `Authorization` header and is answered 401 `{"errors":{"detail":"Unauthorized"}}` by `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112, bearer token read only from the Authorization header at middleware.go:131, 178-196). Add a scheme mirroring `a_aea` (docs/api/astarte_appengine_api.yaml:1436-1443: apiKey in header `Authorization`, described as a realm JWT carrying the `a_ch` claim) and set `security` on this operation only — root `security: []` stays so health/readiness/metrics/version remain documented as unauthenticated. If you model the Phoenix twin too, its credential is `?token=` in the query (internal/appengine/channels/ws.go:53-54), i.e. apiKey `in: query`, not header. Description notes worth carrying: the scheme is case-insensitive `Bearer` with an optional colon, and an unknown realm is 401 not 404 (middleware.go:68-72). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-native-socket-event-schema [auto]: document the native socket's event payload in docs/api/astrate_native_api.yaml — currently only a fenced `incoming_data` example inside the operation description (yaml:300-310), with no `components.schemas` entry anywhere (yaml:425-478), and that example is not the general shape: `wireEvent` marks `interface`, `path` and `value` `omitempty` (internal/appengine/stream/ws.go:51-58), so measured on a verbatim copy a lifecycle event marshals to `{"event":"device_connected","realm":"test","device_id":"dev1","timestamp":"2024-05-30T13:49:57.045Z"}` — no interface/path/value at all (and `value: false`/`value: 0` ARE emitted, because omitempty on an `any` drops only nil). Add a `NativeEvent` schema (event, realm, device_id, optional interface/path/value, timestamp) referenced from the operation, state that `event` is one of the internal/engine/stream/bus.go Kind* constants (bus.go:20-26), and record that this socket drops `InterfaceMajor`, `IP`, `ErrorName` and `ErrorMetadata`, which `stream.Event` carries (bus.go:39-65) — so a `device_error` arrives here with no reason, unlike the Phoenix twin, which marshals the whole struct (internal/appengine/channels/ws.go:406-416). That last sentence documents a code gap, not a spec one: do not add the fields; if they are wanted on the wire, that is a separate code task. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1
- [x] testutil-wait-since-cursor [auto]: give AstarteDevice's message waits a cursor so a second wait for the same topic cannot re-match an older message — `WaitForMessage`/`WaitForTopic` (internal/testutil/astartedevice.go:149-170) scan the capture buffer from index 0 on every poll, so a repeated wait returns the first match again, and in-tree that makes testE2EEmptyCache's resync wait vacuous (see the companion line). Move the paho capture handler out of the `ConnectAstarteDevice` closure (astartedevice.go:61-67) into an injectable method, add `Mark() int` and a `From(start int, …)` wait variant, and add a container-free internal/testutil/astartedevice_test.go that drives two captures on one topic and asserts a wait started after the second returns the second message — it cannot even compile today, which is the failure proof. Keep every new read of `received` under the existing `d.mu`, and say in the report that the change is unverified for races (the Pi gate runs no `-race`). The consumer of this API is engine-e2e-emptycache-resync-assert.
- [ ] engine-e2e-emptycache-resync-assert [legion] [auto]: make the second wait in `testE2EEmptyCache` (internal/engine/engine_e2e_test.go:367) wait for the resync re-publish instead of re-matching the retained set captured at :362 — mark the cursor before `dev.EmptyCache(t)` (:364) and take only messages after the mark, so the "resent property value = int32(42)" assertion at :368-369 fails when the resend never happens; today it cannot fail, because `WaitForTopic` returns the pre-EmptyCache message, which carries the same value 42. Needs `-tags e2e` plus Docker, so this half is Legion-only — the container-free helper and its test land in testutil-wait-since-cursor.
- [!] testutil-control-frame-test [auto]: pin the producer/properties control frame that `SendProducerProperties` puts on the wire (internal/testutil/astartedevice.go:181-216) — `DeflateControlList`/`InflateControlList` have zero tests and every caller of the helpers sits behind an `e2e`/`integration`/conformance build tag, so nothing on the default gate asserts the 4-byte big-endian uncompressed-length prefix plus zlib framing, while the parser side is pinned in internal/engine/control_test.go:89-130 against a CPython golden frame. Add a container-free internal/testutil/astartedevice_test.go with exact-bytes, round-trip, empty-list→nil, sub-4-byte header and declared≠inflated cases; bound the inflate with `io.LimitReader(zr, int64(declared)+1)` the way internal/engine/control.go:301-315 does (astartedevice.go:205 reads with unbounded `io.ReadAll` today, and a lying header from the code under test makes the harness allocate it before failing) and assert the truncated byte count in the mismatch message so the bound itself is pinned. One extra line in internal/engine/control_test.go feeding `testutil.DeflateControlList` into `inflateProperties` is the cross-check that closes the loop. — BLOCKED: wrote nothing
- [x] engine-introspection-producer-roundtrip [auto]: pin the producer half of the docs/DESIGN.md §3.3 introspection rule on the default gate — `testutil.Introspection` (internal/testutil/astartedevice.go:93-105) is called from six places, all behind `e2e`, `integration && e2e` or test/conformance build tags, while the consumer `parseIntrospection` (internal/engine/introspection.go:99-127) is tested only against hand-written literals in internal/engine/introspection_test.go:17-71, so nothing container-free checks that the string the harness renders is the one the engine accepts or that the documented "deterministic (sorted) ordering" holds. Add a case to internal/engine/introspection_test.go that builds the string from an unsorted multi-entry map via `testutil.Introspection`, asserts `parseIntrospection` returns the same map, and pins the exact output (sorted, `name:major:minor;…`, no trailing `;`, empty map → empty string).
- [x] docs-sync-pairing-register-404-unreachable [auto]: drop the `404` `DeviceNotFound` response from `registerDevice` in docs/api/astarte_pairing_api.yaml (yaml:140-141) — it is documented but absent: the route is wrapped by `requireAgent` (internal/pairing/http.go:82-83), i.e. `mw.RequireRealm(auth.ClaimPairing)`, which resolves the realm *before* the handler and answers `401 {"errors":{"detail":"Unauthorized"}}` for an unknown one (internal/auth/middleware.go:66-72, "no existence oracle on auth failures"), so `Service.Register`'s own `GetRealmByName` (internal/pairing/service.go:183-186) can return `store.ErrNotFound` only if the realm is deleted between the middleware's lookup and the handler's — a race, not a contract. Nothing else in Register returns `store.ErrNotFound`: `RegisterDevice`'s `ErrDeviceAlreadyConfirmed` becomes 422 (service.go:209-213) and the CountDevices/SetPayloadFormatHint failures are DB errors → 500. No test covers it (internal/pairing/http_test.go asserts 404 only for the agent DELETE at :527 and for health at :595-598), and the `DeviceNotFound` component's own description — "unknown device ID" — does not match what that path would be even if it fired. The premise came from docs-sync-pairing-register-404, which read service.go without the middleware in front of it. Keep the `404` on `unregisterDevice` (reachable: unknown device, and a malformed `deviceID` via service.go:231-234) and on `getHealth`. State the reachability argument in the report; if you find a reachable path instead, keep the response and fix its description rather than deleting it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
- [!] docs-sync-pairing-body-cap-400 [auto]: document the 64 KiB request-body cap on the three body-taking pairing ops — `registerDevice`, `requestCredentials` and `verifyCredentials` all decode through `astarteapi.DecodeData(r.Body, maxBodyBytes, …)` with `maxBodyBytes = 64 << 10` (internal/pairing/http.go:17, 116, 169, 249), and any body over 65536 bytes fails with an error wrapping `ErrBodyTooLarge` (pkg/astarteapi/envelope.go:233-235) that each handler turns into a plain 400 `{"errors":{"detail":"Bad request"}}` (http.go:117-118, 170-171, 250-251). The shared `BadRequest` component (yaml:578-586) says only "malformed or missing data envelope" and no request-body schema carries a size hint, so a client is told nothing about the cap. Widen that one shared description to name the 64 KiB cap and the failure shape (a single edit covers all three operations) and say explicitly that an over-size body is 400, not 413 — do not invent a 413. One clause on `csr` (yaml:414) and `client_crt` (yaml:422) too: the cap is on the whole envelope, not on the PEM. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML). — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-pairing-realm-name-pattern [auto]: constrain `components.parameters.RealmName` in docs/api/astarte_pairing_api.yaml (yaml:359-365, a bare `type: string` described only as "The realm name.") with `pattern: '^[a-z][a-z0-9]*$'` — the realm name is CHECK-constrained to exactly that (migrations/000002_metadata.up.sql:7, read-only, never edit) and the two surfaces whose tasks already ran carry it on their twins (docs/api/astarte_housekeeping_api.yaml:246, docs/api/astrate_native_api.yaml:420 and 467). Because this one parameter is shared by all seven operations, the description must record that an off-pattern name gets three different answers, not one: the agent and device routes answer **401** (the realm is resolved by the middleware, which treats a miss as unauthenticated rather than 404 — internal/auth/middleware.go:66-72); `getHealth` answers **404** (service.Health returns store.ErrNotFound → WriteNotFound, internal/pairing/service.go:369-371, http.go:289-290); and `getVersion` answers **200** `{"data":"1.2.2"}` for any realm string whatsoever, because `GET /pairing/v1/{realm}/version` is registered with no middleware and no realm lookup (cmd/astrate/main.go:438 → observability.VersionHandler, internal/observability/compat.go:38-47) — say that in the version operation so nobody reads the shared parameter as a guarantee the handler enforces. Mirror `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go:1247) with a pairing twin only if it stays a few lines; the YAML edit is the deliverable. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
- [x] docs-sync-pairing-deviceid-path-param [auto]: give `components.parameters.DeviceID` the constraints `RegisterRequest.hw_id` already carries in docs/api/astarte_pairing_api.yaml — the parameter (yaml:367-377) is prose-only, while `hw_id` (yaml:386-398) got `pattern: '^[A-Za-z0-9_-]{22}$'` plus `minLength`/`maxLength` from docs-sync-pairing-deviceid-base64url, and `TestPairingDeviceIDEncodingDocumented` (internal/swagger/swagger_test.go:1086) pins only the two descriptions and hw_id's bounds, never the parameter's schema; the realm-management twin was constrained later by docs-sync-rm-deviceid-param and its test (swagger_test.go:1168) does hold the parameter's pattern against `deviceid.Parse`. Add the same three schema lines to the parameter, and record in its description the asymmetry the prose omits: a malformed `deviceID` is **401** on all three device-authenticated operations — `deviceid.Parse` failing inside `authenticateDevice` burns bcrypt and returns `ErrUnauthorized` (internal/pairing/service.go:381-385) → `WriteUnauthorized` (http.go:322-323), exactly as an unknown device or a wrong secret is (service.go:395-397), so the client is told its credentials are bad rather than its ID malformed — but **404** on the agent `unregisterDevice` (service.go:231-234 → http.go:326-327). Extend TestPairingDeviceIDEncodingDocumented to compile the new pattern and hold it against `deviceid.Parse` the way the realm-management test does. Do not encode the strict trailing-bits rule in the pattern; that stays open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
- [x] docs-sync-ae-data-envelope-metadata [auto]: add the `metadata` property to `DataEnvelopeObject` in docs/api/astarte_appengine_api.yaml (yaml:1711-1717) — the schema requires only `data` and types it `object`, but every `format=table` data response carries a sibling `metadata` object. Code: `astarteapi.WriteDataWithMetadata` renders `{"data": v, "metadata": m}` through `metadataEnvelope` (`Metadata map[string]any json:"metadata,omitempty"`, pkg/astarteapi/envelope.go:149-158), the `Tabular` carrier sets it (internal/appengine/data.go:59-64), and both renderers populate it — `renderIndividual` (data.go:256-268, `columns` keyed on the path's last segment + `table_header`) and `renderObject` (data.go:293-321, `columns` + sorted `table_header`). Add `metadata` as an optional object property (`columns: map<string,int>`, `table_header: array of string`) described as present only when `format=table`, on all six data GET 200s that `$ref` this schema (yaml:416, 468, 685, 730, 1244, 1290). While there, correct `data`'s type/description: it is not always an object — `format=structured` individual is `[{value,timestamp}]` (data.go:275-281) and `format=disjoint_tables` is an object keyed by path (data.go:269-274, 322-331). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-data-write-object-errors [auto]: document the two object-aggregated write rejections on the six data PUT/POST ops (putData/publishData, by-alias, in-group) in docs/api/astarte_appengine_api.yaml — each refs `BadRequest` for 400 (yaml:518, 570, 775, 822, 1336, 1384) and `ValueTooLarge` for 422 (yaml:528, 580, 785, 832, 1346, 1394), but three bodies are missing. (a) 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` from `writeBadObjectError` (internal/appengine/http.go:722-748, const `detailUnexpectedObjectKey` :711), reached by the `payload.ReasonOf(err) == payload.ReasonBadObject` case (http.go:684-691) — so the shared `ErrorDetail` schema (yaml:1741-1750) needs an optional `unexpected_keys` array, or a dedicated response. (b) 422 `{"errors":{"detail":"Missing required mapping key"}}` (const `detailMissingRequiredMapping`, http.go:713-720) from the `ReasonMissingRequired` case (http.go:692-702) — add an example beside the value-size one already on those ops. Both were implemented by appengine-unexpected-object-key/appengine-missing-required-422 (upstream master b6d46ad4, #2237) with no spec update. Do NOT invent other `payload.ReasonOf` statuses — that family is the separate blocked line appengine-payload-reason-status-map. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-group-create-devices-required [auto]: in docs/api/astarte_appengine_api.yaml, `GroupCreate` (yaml:1652-1663) marks only `group_name` required, but `createGroup` treats `devices` as required — `groupBody.Devices *[]string` (internal/appengine/http.go:391-397) makes a MISSING key a 422 `{"errors":{"devices":["can't be blank"]}}` and an EMPTY array a 422 `{"errors":{"devices":["should have at least 1 item(s)"]}}` (http.go:405-418), and only the MISSING-key case is unreachable from the schema as written. The shared `ValidationErrors` example (yaml:1841-1844) that POST /groups refs (yaml:948) shows `group_name: can't be blank` alone. Add `devices` to `GroupCreate`'s `required` list (the 201 echo always carries it, http.go:424) and add the two device messages to the POST 422 story. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-ae-notfound-component [auto]: the `NotFound` response component in docs/api/astarte_appengine_api.yaml (yaml:1795-1803, example `errors.detail: Not Found`) is `$ref`d by all twelve group-route 404s (yaml:978, 1026, 1070, 1104, 1152, 1179, 1209, 1250, 1296, 1342, 1390, 1429), yet `astarteapi.WriteNotFound`/`DetailNotFound` ("Not Found", pkg/astarteapi/envelope.go:42-43, 205-207) is emitted nowhere in internal/appengine (rg finds only `WriteDeviceNotFound`, http.go:683). A missing group answers 404 `{"errors":{"detail":"Group not found"}}` (`ErrGroupNotFound`, internal/appengine/service.go:41-42, 508-516; writeError http.go:667-668) and a missing member device answers `"Device not found"` (writeError store.ErrNotFound → WriteDeviceNotFound, http.go:682-683). Repoint/rename the component so the group routes show `Group not found` (add a `GroupNotFound` component, keep `DeviceNotFound` on the routes that can 404 for a missing member) and state which routes can emit which. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] realm-error-mapping-test [auto]: add a container-free test file `internal/realm/errors_test.go` pinning the RM error-mapping layer, which is asserted today only by the integration-tagged `internal/realm/http_test.go` (TimescaleDB) so the Pi cannot catch a regression: (1) `writeInterfaceError` (http.go:352-399) — `ErrMaximumDatabaseRetentionExceeded`→422 `{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, a `*interfaceschema.ViolationsError`→422 nested envelope, `errMajorNotFound`→404 "Interface major not found", `ErrNameCollision`/`ErrNameMismatch`/`ErrMajorMismatch`/`interfaceschema.ErrMinorNotIncreased`/`ErrDowngradeNotAllowed`/`ErrMissingEndpoints`/`ErrIncompatibleEndpointChange`/`store.ErrAlreadyExists`→409, `store.ErrInterfaceMajorNotZero`/`store.ErrInterfaceInUse`→403, plain `ErrValidation`→422 with the `realm: validation failed: ` prefix stripped, `store.ErrNotFound`→404 "Interface not found", default→500; (2) `writeError` (http.go:402-425) — its 422/409/404/default arms including that it maps the two interface sentinels to 422, not 403; (3) `renderViolationsBody` (http.go:443-533) byte-for-byte for a mixed input (interface-level field merged first-seen + the full-length index-aligned `mappings` array with `{}` for clean entries); (4) `validationDetail` — prefix stripped, non-prefix passthrough, exact-prefix NOT stripped. `*API{}` needs no Service (these methods never deref `a.svc`), and `interfaceschema.ViolationsError` is built from its exported `Violations`/`MappingCount` fields, so no store or Docker.
- [x] realm-detailed-properties-branch-test [auto]: extend `internal/realm/detailed_required_encrypted_test.go` to cover the properties branch of `writeDetailedMapping` (internal/realm/detailed.go:61-63) — currently only the datastream branch is pinned container-free (the properties path runs only in the integration-tagged `internal/realm/detailed_test.go`): assert a properties mapping renders exactly `{"endpoint":..,"type":..,"allow_unset":<bool>}` with no reliability/retention/expiry/explicit_timestamp/required/encrypted/database_retention_policy keys, for `allow_unset` both true and false. No Docker.
===FORGIULIO===
-rw-r--r-- 1 root root 163554 Oct  9 21:05 .mule/for-giulio.md

**3. Two `[housekeeping]` config keys exist; `configuration-reference.md` does not list them.**
The table at `docs/site/configuration-reference.md:79-84` has exactly two rows and the section
heading calls the block "Instance-admin keys" — but `HousekeepingConfig`
(`internal/config/config.go:109-125`) has four fields:

- `default_datastream_maximum_storage_retention` (`*int64`, env
  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`, validated as a
  non-negative integer at `config.go:286-291`) — injected into a realm at creation when the
  caller omits the field (`internal/housekeeping/service.go:149-150`). An operator who wants it
  has no way to learn it exists.
- `realm_deletion_disabled` (bool, env `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`, a
  fail-loud boolean accepting only `1|true|TRUE|True|0|false|""` and refusing to load on anything
  else, `config.go:296-306`) — this is the switch behind the spec's documented
  `405 "Realm deletion disabled"` on `DELETE /housekeeping/v1/realms/{realm}`, so it changes a
  status code an API client sees.

Both belong in that table next to the two key lists, and the heading wants widening.

Checked and *not* escalated, for the record. The config-key sweep
(`rg -o '\bASTRATE_[A-Z_]+' -N internal/ | sort -u`) found nothing else stale: the two
`jwt_public_key*` keys the page quotes at lines 69-71 are real (`config.go:112-113`), and the
`master_key_file` row's `ASTRATE_MASTER_KEY`/`ASTRATE_MASTER_KEY_FILE` fallbacks are real too
(`internal/store/crypto.go:24-27`). And the one thing that looked like a live bug is not one:
`view()` returns only `JWTPublicKeysPEM[0]` (`internal/housekeeping/service.go:277-281`) while
the realm-management surface joins *all* keys (`internal/realm/service.go:627`) — but every
writer stores a single-element array (create at `internal/store/realms.go:54-61`, PATCH at
`internal/store/realms.go:165`, `putAuthConfig` at `internal/realm/service.go:636`), so a
two-key realm is unreachable and there is nothing to document. Worth noting because it reads
like a bug on first inspection.

$ ls internal/appengine internal/housekeeping internal/pairing internal/realm
internal/appengine:
channels
dashboard_compat_test.go
data.go
downsample.go
downsample_test.go
e2e_test.go
formats_parity_test.go
groups_parity_test.go
http.go
http_test.go
links_test.go
mirror_test.go
patch_test.go
query_opts_test.go
service.go
stream
writeerror_test.go

internal/housekeeping:
async_operation_test.go
http.go
http_test.go
service.go
service_test.go

internal/pairing:
bearersecret_test.go
ca
http.go
http_test.go
ratelimit.go
ratelimit_test.go
remoteip_test.go
service.go
service_test.go
testdata

internal/realm:
async_operation_test.go
dashboard_compat_test.go
detailed.go
detailed_required_encrypted_test.go
detailed_test.go
errors_test.go
helpers_test.go
http.go
http_test.go
policy.go
policy_test.go
service.go

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc|Mount\(' internal/ --glob '*http.go' -g '!*_test.go' | head -100
internal/flowapi/http.go:31:func (a *API) Mount(mux *http.ServeMux) {
internal/appengine/http.go:41:func (a *API) Mount(mux *http.ServeMux) {
internal/appengine/http.go:210:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:231:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:522:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/housekeeping/http.go:35:func (a *API) Mount(mux *http.ServeMux) {
internal/realm/http.go:36:func (a *API) Mount(mux *http.ServeMux) {
internal/pairing/http.go:77:func (a *API) Mount(mux *http.ServeMux) {
internal/pairing/http.go:81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/pairing/http.go:86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
internal/pairing/http.go:88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
internal/pairing/http.go:90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
internal/pairing/http.go:344:	header := r.Header.Get("Authorization")

→ Read internal/realm/http.go [limit=70]
→ Read internal/appengine/http.go [limit=120]

$ rg -n '^\s{2}/' docs/api/astarte_realm_management_api.yaml
27:  /realmmanagement/v1/{realm}/interfaces:
172:  /realmmanagement/v1/{realm}/interfaces/{name}:
203:  /realmmanagement/v1/{realm}/interfaces/{name}/{major}:
352:  /realmmanagement/v1/{realm}/triggers:
449:  /realmmanagement/v1/{realm}/triggers/{name}:
505:  /realmmanagement/v1/{realm}/config/auth:
589:  /realmmanagement/v1/{realm}/config/device_registration_limit:
622:  /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention:
659:  /realmmanagement/v1/{realm}/version:
694:  /realmmanagement/v1/{realm}/devices/{device}:
723:  /realmmanagement/v1/{realm}/policies:
817:  /realmmanagement/v1/{realm}/policies/{name}:

→ Read docs/api/astarte_realm_management_api.yaml
→ Read internal/realm/http.go [offset=70, limit=330]
→ Read internal/realm/http.go [limit=170, offset=400]

$ rg -n 'docs-sync-rm|Interface not found|realm-management_api|Interface major not found' .mule/todo.md .mule/for-giulio.md
.mule/todo.md:80:- [x] docs-sync-rm-policies-delete-422 [auto]: add the missing `422` response to `DELETE /realmmanagement/v1/{realm}/policies/{name}` in docs/api/astarte_realm_management_api.yaml — `Service.DeletePolicy` refuses to delete a policy still referenced by a trigger (`"policy %q is still used by trigger %q"`, internal/realm/service.go:573-574, wrapped as ErrValidation) and `writeError` answers 422 for it (internal/realm/http.go:410-411, detail stripped by validationDetail), but the spec lists only 204/401/404/500 for the delete op (yaml:695-713). Add the 422 ValidationError response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:81:- [x] docs-sync-rm-triggers-422-nested-envelope [auto]: document the nested changeset envelope on `POST /realmmanagement/v1/{realm}/triggers` in docs/api/astarte_realm_management_api.yaml — when trigger compile fails, `createTrigger` answers 422 with the upstream nested body `{"errors":{"action":{...},"simple_triggers":[...]}}` (internal/realm/http.go:252-258, triggerErrorBody at http.go:268-285, issues #63/#70), but the spec's createTrigger 422 (yaml:344-345) refs only the flat `ValidationError` component (`{"errors":{"detail":...}}`) — which the handler also emits for plain ErrValidation (http.go:259-260, e.g. "trigger requires a name", service.go:429). Model the response as oneOf across the flat ErrorDetail and the nested envelope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:123:- [x] docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` route to docs/api/astarte_realm_management_api.yaml — it exists in code (internal/realm/http.go:51, handler getDatastreamMaximumStorageRetention at http.go:135, returns 200 with a data-envelope retention value via Service.GetDatastreamMaximumStorageRetention) but the spec jumps straight from config/device_registration_limit (yaml:428) to /version. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:124:- [x] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:158:- [x] docs-sync-rm-delete-interface-status [auto]: fix the `deleteInterface` response in docs/api/astarte_realm_management_api.yaml — the spec documents `422` with examples "Interface major version is not 0, can't be deleted" / "Cannot delete an interface that is used by a device introspection" (yaml:255-272), but the handler returns `403` with "Interface can't be deleted" / "Interface can't be deleted since it's currently used" (internal/realm/http.go:387-391, `writeInterfaceError`; both are 403 `StatusForbidden`). The spec's detail strings were lifted from the `writeError` branch (http.go:414-419) which DELETE /interfaces does not use. Fix status to `403` and update the two examples to match the actual wire messages. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:159:- [x] docs-sync-rm-mapping-required-encrypted [auto]: add `required` (boolean, datastream-only) and `encrypted` (boolean, datastream-only) to the `InterfaceMapping` schema in docs/api/astarte_realm_management_api.yaml — the parser accepts and persists both on datastream mappings (internal/interfaceschema/parse.go:410-414), rejects them on properties mappings (parse.go:395-398), the detailed listing emits them (`detailed.go:69`), and the compatibility check compares them (interfaceschema/compat.go:121-123); the schema (yaml:820-862) currently has `additionalProperties: false` and neither field, so the OpenAPI would mark a valid datastream interface install body as invalid. Also add both fields to the detailed-listing example (yaml:70-85). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:160:- [x] docs-sync-rm-put-auth-422 [auto]: add the missing `422` (ValidationError) response to `PUT /realmmanagement/v1/{realm}/config/auth` in docs/api/astarte_realm_management_api.yaml — `SetAuthKey` returns `ErrValidation` for a blank `jwt_public_key_pem` (internal/realm/service.go:632-640), which `writeError` maps to 422; the spec (yaml:431-464) documents only 204/400/401/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:161:- [!] docs-sync-rm-version-example [auto]: fix the `getVersion` example in docs/api/astarte_realm_management_api.yaml from `data: "1.1.0"` (yaml:557) to `data: "1.2.2"` — `APICompatVersion` is `"1.2.2"` (internal/realm/service.go:588); the native spec's four compat-version examples already use "1.2.2". Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1
.mule/todo.md:216:- [x] docs-sync-rm-put-interface-409 [auto]: add the missing `409` Conflict response to `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` in docs/api/astarte_realm_management_api.yaml (yaml:222-234 lists only 204/400/401/404/422/500) — `updateInterface` answers 409 for every url/body disagreement and minor-upgrade incompatibility: `ErrNameMismatch` (service.go:204 → "Interface name doesn't match the one in the interface json", http.go:367-369), `ErrMajorMismatch` (service.go:207 → http.go:370-372), and the `CheckMinorUpgrade` sentinels `ErrMinorNotIncreased`/`ErrDowngradeNotAllowed`/`ErrMissingEndpoints`/`ErrIncompatibleEndpointChange` (service.go:223-226 → http.go:373-384) plus `store.ErrAlreadyExists` (http.go:385-386) — the exact set `TestRealmManagementErrorCodes` pins as StatusConflict (http_test.go:417-437). Reuse the `Conflict` response component (already used on POST, yaml:131-132) and add name-mismatch/major-mismatch examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:217:- [x] docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` (yaml:231-232) in docs/api/astarte_realm_management_api.yaml as a oneOf — both handlers go through `writeInterfaceError`, which answers three distinct 422 bodies: the flat ErrorDetail (`ErrValidation` via `validationDetail`, http.go:392-393), the nested violations changeset envelope (`writeViolations`, http.go:358-359 + 434-533, e.g. `{"errors":{"description":["should be at most 1000 character(s)"]}}` and the aligned full-length `mappings` array, http_test.go:480-501), and the named FieldErrors envelope for `ErrMaximumDatabaseRetentionExceeded` (`{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, http.go:355-357, http_test.go:282-320); today only the flat ValidationError is referenced. Mirror the createTrigger oneOf pattern (yaml:344-368). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:218:- [x] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "realm: validation failed: interface definition is invalid"`, but `validationDetail` strips the `realm: validation failed: ` prefix on the wire (http.go:562-568) and the parser's real messages are `invalid interface: ...` (pkg/interfaceschema/parse.go:38, wrapped at service.go:163); the example should show the stripped, real message. Same class as the already-fixed docs-sync-hk-validation-example (housekeeping http.go:251-258). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:219:- [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:228:- [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:229:- [x] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion Go — the device-delete leg added by docs-sync-rm-delete-device-async-operation-param (register a device, `DELETE /devices/{id}?async_operation=…` answers 204 and the row is gone on both values) compiles clean under `go vet -tags integration ./internal/realm/` but has never been executed, because the Pi has no database; until it runs, TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go is the only assertion of that leg that has actually been verified.
.mule/todo.md:230:- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
.mule/todo.md:236:- [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml — the string `"403"` appears nowhere in the file and `components` defines no `Forbidden` response, yet every route is wrapped by `a.require = mw.RequireRealm(auth.ClaimAppEngine)` (internal/appengine/http.go:35), which answers `403 Forbidden` for a verified a_aea JWT whose grants do not authorize the method+path (`astarteapi.WriteForbidden`, internal/auth/middleware.go:90; `DetailForbidden`, pkg/astarteapi/envelope.go:41; pinned today by internal/appengine/http_test.go:184-185). Add a `Forbidden` response component copied from the realm-management one (docs/api/astarte_realm_management_api.yaml:1277-1285 — `ErrorDetail` schema, example `errors.detail: Forbidden`) and `$ref` it from all 32 operations, which is exactly the parity docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go:123-170) pins; add the same docs-half test for the AppEngine spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:254:- [x] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, housekeeping and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
.mule/todo.md:255:- [!] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy alias fields to `InterfaceDefinition` and `InterfaceMapping` in docs/api/astarte_realm_management_api.yaml — both declare `additionalProperties: false` (yaml:935, 972) yet omit interface-level `quality` (alias of `ownership`) and `aggregate` (alias of `aggregation`) and mapping-level `path` (alias of `endpoint`), all three of which the parser accepts and canonicalises: decoded into dedicated fields (pkg/interfaceschema/parse.go:116-122, 136-138), normalised before any rule runs (parse.go:212-236), and re-encoded canonically by `ParseInterfaceCanonical` so the store never sees them (parse.go:165-170, 307-311). A spec-conformant client — or any generated client with this schema baked in — therefore rejects a valid upstream-shaped install/update with a spurious 422 while the service would have accepted and stored it as `ownership`/`aggregation`/`endpoint`. Add the three as documented deprecated aliases (enums mirroring ownership/aggregation, `path` a plain string) or relax `additionalProperties`, and record the three rules the code enforces: `ownership` + `quality` is a violation (parse.go:218) and `aggregation` + `aggregate` is a violation (parse.go:227), while `endpoint` + `path` silently keeps `endpoint` (parse.go:307-311). None of the three may appear in a `required` list. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
.mule/todo.md:256:- [x] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:257:- [x] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstrained string described only as "The device hardware ID." — `Service.DeleteDevice` runs the path value through `deviceid.Parse` (internal/realm/service.go:118-121), which accepts exactly 22 characters of unpadded base64url in strict mode and rejects padding, the standard alphabet's `+` and `/`, and non-canonical trailing bits (pkg/deviceid/deviceid.go:24-26, 33-37, 42-60); every rejection is folded into `store.ErrNotFound` (service.go:119-121) and answered `404 {"errors":{"detail":"Device not found"}}` (internal/realm/http.go:107-109), so a client that sends a well-formed standard-base64 or 23-character padded id is told the device does not exist rather than that its id is malformed. State the encoding and the 22-character length and add `pattern`/`minLength`/`maxLength`, exactly as docs-sync-pairing-deviceid-base64url did for the pairing twin. Do NOT encode the strict trailing-bits rule in the pattern: whether it is a deviation worth recording is still open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:258:- [!] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_realm_management_api.yaml — `data` is declared as a bare `type: object` with no properties (yaml:250-261) although the handler parses it with the same `interfaceschema.ParseInterfaceCanonical` as `installInterface` (internal/realm/service.go:196 vs 161) and then checks the body's `interface_name`/`version_major` against the URL's `{name}`/`{major}` (service.go:203-208), so the update body — the one carrying every interface rule — is the only interface body in the spec with nothing to validate against. `$ref` `#/components/schemas/InterfaceDefinition` as `installInterface` does (yaml:115-117), and say in the operation description that the body must name the same interface and major as the URL (disagreement is 409, already documented) and that `version_minor` must increase (409). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1201s — task too big, split it
.mule/todo.md:268:- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:298:- [x] docs-sync-pairing-deviceid-path-param [auto]: give `components.parameters.DeviceID` the constraints `RegisterRequest.hw_id` already carries in docs/api/astarte_pairing_api.yaml — the parameter (yaml:367-377) is prose-only, while `hw_id` (yaml:386-398) got `pattern: '^[A-Za-z0-9_-]{22}$'` plus `minLength`/`maxLength` from docs-sync-pairing-deviceid-base64url, and `TestPairingDeviceIDEncodingDocumented` (internal/swagger/swagger_test.go:1086) pins only the two descriptions and hw_id's bounds, never the parameter's schema; the realm-management twin was constrained later by docs-sync-rm-deviceid-param and its test (swagger_test.go:1168) does hold the parameter's pattern against `deviceid.Parse`. Add the same three schema lines to the parameter, and record in its description the asymmetry the prose omits: a malformed `deviceID` is **401** on all three device-authenticated operations — `deviceid.Parse` failing inside `authenticateDevice` burns bcrypt and returns `ErrUnauthorized` (internal/pairing/service.go:381-385) → `WriteUnauthorized` (http.go:322-323), exactly as an unknown device or a wrong secret is (service.go:395-397), so the client is told its credentials are bad rather than its ID malformed — but **404** on the agent `unregisterDevice` (service.go:231-234 → http.go:326-327). Extend TestPairingDeviceIDEncodingDocumented to compile the new pattern and hold it against `deviceid.Parse` the way the realm-management test does. Do not encode the strict trailing-bits rule in the pattern; that stays open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
.mule/todo.md:303:- [x] realm-error-mapping-test [auto]: add a container-free test file `internal/realm/errors_test.go` pinning the RM error-mapping layer, which is asserted today only by the integration-tagged `internal/realm/http_test.go` (TimescaleDB) so the Pi cannot catch a regression: (1) `writeInterfaceError` (http.go:352-399) — `ErrMaximumDatabaseRetentionExceeded`→422 `{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, a `*interfaceschema.ViolationsError`→422 nested envelope, `errMajorNotFound`→404 "Interface major not found", `ErrNameCollision`/`ErrNameMismatch`/`ErrMajorMismatch`/`interfaceschema.ErrMinorNotIncreased`/`ErrDowngradeNotAllowed`/`ErrMissingEndpoints`/`ErrIncompatibleEndpointChange`/`store.ErrAlreadyExists`→409, `store.ErrInterfaceMajorNotZero`/`store.ErrInterfaceInUse`→403, plain `ErrValidation`→422 with the `realm: validation failed: ` prefix stripped, `store.ErrNotFound`→404 "Interface not found", default→500; (2) `writeError` (http.go:402-425) — its 422/409/404/default arms including that it maps the two interface sentinels to 422, not 403; (3) `renderViolationsBody` (http.go:443-533) byte-for-byte for a mixed input (interface-level field merged first-seen + the full-length index-aligned `mappings` array with `{}` for clean entries); (4) `validationDetail` — prefix stripped, non-prefix passthrough, exact-prefix NOT stripped. `*API{}` needs no Service (these methods never deref `a.svc`), and `interfaceschema.ViolationsError` is built from its exported `Violations`/`MappingCount` fields, so no store or Docker.
.mule/for-giulio.md:94:  **The knock-on: `docs/COMPATIBILITY.md` deviation 17 repeats the same error, and its wording about one endpoint is now measurably wrong.** It opens "**Always synchronous where upstream 1.4 defaults to asynchronous** — upstream 1.4 runs realm create/delete … and lets the caller opt into synchronous execution with `?async_operation=false`" (docs/COMPATIBILITY.md:335-347). Per the CHANGELOG the async default and the opt-out both date to 1.0.2, not 1.4 — the behaviour is not a moving target, which is the opposite of how that sentence reads. Separately, the sentence lists **device deletion** among the endpoints where the caller can opt into synchronous execution, and upstream has **no such opt-out there**: in `apps/astarte_realm_management/lib/astarte_realm_management_web/controllers/device_controller.ex`, the `operation :delete` doc comment says "Device deletion happens asynchronously, and receiving a 204 response" (:31-38) and the handler is bare — `with :ok <- Devices.delete_device(realm_name, device_id)` (:52-53), with `async_operation` never read. The policy twin does honour it, read as a bare string compare with no parsing (`trigger_policy_controller.ex:155-167`: `if Map.get(params, "async_operation") == "false"`), so upstream accepts only the literal `"false"` and treats absent/`0`/`no` as async. Astrate accepting-and-ignoring the parameter on device delete is therefore an **Astrate extension, not parity** — harmless and stronger, but the record calls it parity. For whoever picks up the queued `docs-sync-rm-delete-device-async-operation-param` line: documenting that parameter on the device-delete operation is still right for an Astrate client, it just has to be recorded as an extension, because upstream's spec has no such parameter on that operation.

$ rg -n 'NotFound|errMajorNotFound|Interface not found' internal/appengine/http.go internal/housekeeping/http.go internal/pairing/http.go | head -40
internal/housekeeping/http.go:244:	case errors.Is(err, store.ErrNotFound):
internal/housekeeping/http.go:245:		_ = astarteapi.WriteNotFound(w)
internal/appengine/http.go:653:	case errors.Is(err, ErrAliasTagNotFound):
internal/appengine/http.go:657:	case errors.Is(err, ErrAttributeKeyNotFound):
internal/appengine/http.go:667:	case errors.Is(err, ErrGroupNotFound):
internal/appengine/http.go:668:		_ = astarteapi.WriteError(w, http.StatusNotFound, "Group not found")
internal/appengine/http.go:669:	case errors.Is(err, ErrPathNotFound):
internal/appengine/http.go:670:		_ = astarteapi.WriteError(w, http.StatusNotFound, "Path not found")
internal/appengine/http.go:671:	// Engine write sentinels must precede store.ErrNotFound: wrapped
internal/appengine/http.go:672:	// ErrInterfaceNotFound errors are distinct from store's, and Go's switch
internal/appengine/http.go:678:	case errors.Is(err, engine.ErrInterfaceNotFound):
internal/appengine/http.go:679:		_ = astarteapi.WriteError(w, http.StatusNotFound, "Interface not found in device introspection")
internal/appengine/http.go:680:	case errors.Is(err, engine.ErrPathNotFound):
internal/appengine/http.go:682:	case errors.Is(err, store.ErrNotFound):
internal/appengine/http.go:683:		_ = astarteapi.WriteDeviceNotFound(w)
internal/pairing/http.go:289:	case errors.Is(err, store.ErrNotFound):
internal/pairing/http.go:290:		_ = astarteapi.WriteNotFound(w)
internal/pairing/http.go:326:	case errors.Is(err, store.ErrNotFound):
internal/pairing/http.go:327:		_ = astarteapi.WriteDeviceNotFound(w)

$ rg -n 'APICompatVersion' internal/realm/service.go
580:// APICompatVersion is the upstream Realm Management API level Astrate
588:const APICompatVersion = "1.2.2"

$ rg -n 'APICompatVersion|1\.2\.2' internal/ cmd/ | head -20
cmd/astrate/main.go:437:		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(realm.APICompatVersion)))
cmd/astrate/main.go:438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
cmd/astrate/version_test.go:25:	want := `{"data":"` + realm.APICompatVersion + `"}`
internal/realm/http.go:119:// feature UI on it — see APICompatVersion).
internal/realm/http.go:121:	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
internal/realm/dashboard_compat_test.go:26:// Dashboard v1.2.2 requires: version, device_registration_limit, delivery
internal/realm/dashboard_compat_test.go:35:		if v != APICompatVersion {
internal/realm/dashboard_compat_test.go:36:			t.Errorf("version = %q, want %q", v, APICompatVersion)
internal/realm/service.go:580:// APICompatVersion is the upstream Realm Management API level Astrate
internal/realm/service.go:588:const APICompatVersion = "1.2.2"
internal/engine/triggers/errorname.go:3:// upstreamErrorNames is the closed set astarte-dashboard 1.2.2 validates
internal/engine/triggers/errorname.go:63:	// 1.2.2's closed set; translated to the closest 1.2.2 member.
internal/broker/aclhook.go:73:		// out of scope for Astrate's current 1.2.2 target, and upstream's wire
internal/appengine/dashboard_compat_test.go:45:// TestDashboardCompat covers the M10 surface the Astarte Dashboard v1.2.2
internal/appengine/service.go:30:// Device-PATCH error taxonomy, frozen to upstream astarte_appengine_api 1.2.2
internal/appengine/service.go:34:	ErrInvalidAlias         = errors.New("Invalid alias")           //nolint:staticcheck // ST1005: detail strings are upstream astarte_appengine_api 1.2.2 wire text

→ Read internal/realm/http_test.go [offset=380, limit=120]

$ rg -n 'bodyNotFound|bodyMajorNotFound|bodyDup|bodyCantDelete|bodyInUse' internal/realm/http_test.go | head -20
338:		bodyDup           = `{"errors":{"detail":"Interface already exists"}}`
346:		bodyMajorNotFound = `{"errors":{"detail":"Interface major not found"}}`
347:		bodyNotFound      = `{"errors":{"detail":"Interface not found"}}`
348:		bodyCantDelete    = `{"errors":{"detail":"Interface can't be deleted"}}`
349:		bodyInUse         = `{"errors":{"detail":"Interface can't be deleted since it's currently used"}}`
396:			nil, http.StatusNotFound, bodyNotFound},
398:			nil, http.StatusNotFound, bodyNotFound},
403:			nil, http.StatusConflict, bodyDup},
442:			nil, http.StatusNotFound, bodyMajorNotFound},
444:			nil, http.StatusNotFound, bodyMajorNotFound},
446:			nil, http.StatusForbidden, bodyCantDelete},
448:			setupDraftInUse, http.StatusForbidden, bodyInUse},

$ rg -n 'json:"' pkg/interfaceschema/*.go | rg -v '_test' | head -60
pkg/interfaceschema/parse.go:111:	InterfaceName *string        `json:"interface_name"`
pkg/interfaceschema/parse.go:112:	VersionMajor  *int           `json:"version_major"`
pkg/interfaceschema/parse.go:113:	VersionMinor  *int           `json:"version_minor"`
pkg/interfaceschema/parse.go:114:	Type          *InterfaceType `json:"type"`
pkg/interfaceschema/parse.go:115:	Ownership     *Ownership     `json:"ownership,omitempty"`
pkg/interfaceschema/parse.go:119:	Quality *Ownership `json:"quality,omitempty"`
pkg/interfaceschema/parse.go:121:	Aggregate   *bool         `json:"aggregate,omitempty"`
pkg/interfaceschema/parse.go:122:	Aggregation *Aggregation  `json:"aggregation,omitempty"`
pkg/interfaceschema/parse.go:123:	Description string        `json:"description,omitempty"`
pkg/interfaceschema/parse.go:124:	Doc         string        `json:"doc,omitempty"`
pkg/interfaceschema/parse.go:125:	Mappings    []mappingJSON `json:"mappings"`
pkg/interfaceschema/parse.go:130:	aliases bool `json:"-"`
pkg/interfaceschema/parse.go:135:	Endpoint *string `json:"endpoint,omitempty"`
pkg/interfaceschema/parse.go:138:	Path                    *string                  `json:"path,omitempty"`
pkg/interfaceschema/parse.go:139:	Type                    *ValueType               `json:"type"`
pkg/interfaceschema/parse.go:140:	Reliability             *Reliability             `json:"reliability,omitempty"`
pkg/interfaceschema/parse.go:141:	Retention               *Retention               `json:"retention,omitempty"`
pkg/interfaceschema/parse.go:142:	Expiry                  *int64                   `json:"expiry,omitempty"`
pkg/interfaceschema/parse.go:143:	DatabaseRetentionPolicy *DatabaseRetentionPolicy `json:"database_retention_policy,omitempty"`
pkg/interfaceschema/parse.go:144:	DatabaseRetentionTTL    *int64                   `json:"database_retention_ttl,omitempty"`
pkg/interfaceschema/parse.go:145:	AllowUnset              *bool                    `json:"allow_unset,omitempty"`
pkg/interfaceschema/parse.go:146:	ExplicitTimestamp       *bool                    `json:"explicit_timestamp,omitempty"`
pkg/interfaceschema/parse.go:147:	Required                *bool                    `json:"required,omitempty"`
pkg/interfaceschema/parse.go:148:	Encrypted               *bool                    `json:"encrypted,omitempty"`
pkg/interfaceschema/parse.go:149:	Description             string                   `json:"description,omitempty"`
pkg/interfaceschema/parse.go:150:	Doc                     string                   `json:"doc,omitempty"`

$ rg -n 'policy must have|policy name must|maximum_capacity|event_ttl|retry_times' internal/engine/triggers/policy.go | head -30
104:		MaximumCapacity int  `json:"maximum_capacity"`
106:		RetryTimes      *int `json:"retry_times"`
107:		EventTTL        *int `json:"event_ttl"`
113:		return nil, fmt.Errorf("triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]")
116:		return nil, fmt.Errorf("triggers: policy must have at least one error handler")
149:			return nil, fmt.Errorf("triggers: retry_times must be 1-100 when any handler retries")
152:		return nil, fmt.Errorf("triggers: retry_times requires a retry handler")
155:		return nil, fmt.Errorf("triggers: maximum_capacity must be a positive integer")
161:		return nil, fmt.Errorf("triggers: event_ttl must be non-negative")

$ rg -o '\bASTRATE_[A-Z_]+' -N internal/ | sort -u
internal/config/config.example.toml:ASTRATE_DATABASE_DSN
internal/config/config.example.toml:ASTRATE_ENGINE_SHARDS
internal/config/config.example.toml:ASTRATE_HTTP_ADDR
internal/config/config.example.toml:ASTRATE_LOG_LEVEL
internal/config/config.example.toml:ASTRATE_MASTER_KEY
internal/config/config.example.toml:ASTRATE_MASTER_KEY_FILE
internal/config/config.example.toml:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/config/config.go:ASTRATE_DATABASE_DSN
internal/config/config.go:ASTRATE_ENGINE_SHARDS
internal/config/config.go:ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
internal/config/config.go:ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
internal/config/config.go:ASTRATE_HTTP_ADDR
internal/config/config.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
internal/config/config.go:ASTRATE_HTTP_TLS_CERT_FILE
internal/config/config.go:ASTRATE_HTTP_TLS_KEY_FILE
internal/config/config.go:ASTRATE_LOG_FORMAT
internal/config/config.go:ASTRATE_LOG_LEVEL
internal/config/config.go:ASTRATE_MASTER_KEY
internal/config/config.go:ASTRATE_MASTER_KEY_FILE
internal/config/config.go:ASTRATE_MQTT_ADDR
internal/config/config.go:ASTRATE_MQTT_ADVERTISED_URL
internal/config/config.go:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/config/config.go:ASTRATE_MQTT_SESSION_STORE_PATH
internal/config/config.go:ASTRATE_MQTT_TLS_CERT_FILE
internal/config/config.go:ASTRATE_MQTT_TLS_KEY_FILE
internal/config/config.go:ASTRATE_REALM_JWT_PUBLIC_KEY
internal/config/config.go:ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
internal/config/config.go:ASTRATE_REALM_NAME
internal/config/config.go:ASTRATE_SECURITY_MASTER_KEY_FILE
internal/config/config_test.go:ASTRATE_DATABASE_DSN
internal/config/config_test.go:ASTRATE_ENGINE_SHARDS
internal/config/config_test.go:ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
internal/config/config_test.go:ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
internal/config/config_test.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
internal/config/config_test.go:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/engine/engine_integration_test.go:ASTRATE_TEST_DSN
internal/flow/blocks/container/block.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/container/block_test.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/container/docker.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/info.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/schema.go:ASTRATE_FLOW_CONFIG
internal/store/crypto.go:ASTRATE_MASTER_KEY
internal/store/crypto.go:ASTRATE_MASTER_KEY_FILE
internal/store/store_test.go:ASTRATE_TEST_DSN
internal/testutil/pg.go:ASTRATE_TEST_DSN

$ rg -n 'ASTRATE_|Astrate Env|env' docs/site/configuration-reference.md | head -60
3:Astrate reads one TOML file (`astrate -config path.toml`) with `ASTRATE_*` environment overrides. Precedence: **built-in defaults < TOML file < environment variables**.
11:- `database.dsn` (or `ASTRATE_DATABASE_DSN`) — PostgreSQL/TimescaleDB connection string.
19:built-in defaults  <  TOML file  <  ASTRATE_* environment
22:Every value below is the built-in default unless noted. Environment overrides exist for operationally critical fields, named `ASTRATE_<SECTION>_<FIELD>` (e.g. `ASTRATE_DATABASE_DSN`, `ASTRATE_HTTP_ADDR`).
30:| `dsn` | string | `""` | `ASTRATE_DATABASE_DSN` | **Yes** | PostgreSQL + TimescaleDB connection string. |
36:| `addr` | string | `":8080"` | `ASTRATE_HTTP_ADDR` | No | Bind address for the single REST listener (pairing, realm management, housekeeping, appengine, native endpoints). |
37:| `tls_cert_file` | string | `""` | `ASTRATE_HTTP_TLS_CERT_FILE` | No | Optional in-binary TLS cert. Leave empty to terminate TLS at a reverse proxy. |
38:| `tls_key_file` | string | `""` | `ASTRATE_HTTP_TLS_KEY_FILE` | No | Optional in-binary TLS key. Must be set together with `tls_cert_file`. |
39:| `cors_allowed_origins` | list[string] | `[]` | `ASTRATE_HTTP_CORS_ALLOWED_ORIGINS` (comma-separated) | No | Browser origins allowed for CORS. `"*"` allows any; empty disables CORS. The Astarte Dashboard SPA needs its origin here (e.g. `"http://localhost:4040"`). |
45:| `addr` | string | `":8883"` | `ASTRATE_MQTT_ADDR` | No | Bind address for the mTLS MQTT listener. |
46:| `tls_cert_file` | string | `""` | `ASTRATE_MQTT_TLS_CERT_FILE` | **Yes*** | Broker server cert. Required unless `insecure_dev_mode` is true. |
47:| `tls_key_file` | string | `""` | `ASTRATE_MQTT_TLS_KEY_FILE` | **Yes*** | Broker server key. Required unless `insecure_dev_mode` is true. |
48:| `insecure_dev_mode` | bool | `true` | `ASTRATE_MQTT_INSECURE_DEV_MODE` | No | Bind a plaintext listener on `dev_addr` that trusts the claimed client ID. **Never enable in production.** |
50:| `session_store_path` | string | `"sessions.db"` | `ASTRATE_MQTT_SESSION_STORE_PATH` | No | bbolt/pebble file for MQTT session persistence. On a read-only container FS, must point at the writable session volume. |
51:| `advertised_url` | string | `""` | `ASTRATE_MQTT_ADVERTISED_URL` | No | Broker URL handed to devices by the pairing info endpoint. Empty derives `"mqtts://<addr>"`; set when devices reach the broker by another host. |
60:| `shards` | int | `16` | `ASTRATE_ENGINE_SHARDS` | No | Number of ordered ingestion shards. Each shard processes messages strictly in order. |
95:| `master_key_file` | string | `""` | `ASTRATE_SECURITY_MASTER_KEY_FILE` | No | File holding the AES-256 master key that seals realm CA private keys. When empty, falls back to `ASTRATE_MASTER_KEY` (64 hex chars) or `ASTRATE_MASTER_KEY_FILE`. |
104:| `name` | string | `""` | `ASTRATE_REALM_NAME` | No | Realm name. Empty disables auto-provisioning. |
105:| `jwt_public_key` | string | `""` | `ASTRATE_REALM_JWT_PUBLIC_KEY` | Conditional | PEM public key for the auto-provisioned realm. Required when `name` is set. |
106:| `jwt_public_key_file` | string | `""` | `ASTRATE_REALM_JWT_PUBLIC_KEY_FILE` | Conditional | File path to the PEM public key. Alternative to inline `jwt_public_key`. |
113:| `level` | string | `"info"` | `ASTRATE_LOG_LEVEL` | No | Log level: `debug`, `info`, `warn`, `error`. |
114:| `format` | string | `"json"` | `ASTRATE_LOG_FORMAT` | No | Log format: `json` (structured) or `text`. |
149:The included `docker-compose.yml` uses environment variables:
155:This sets `ASTRATE_MQTT_INSECURE_DEV_MODE=true`, `ASTRATE_MASTER_KEY` (throwaway), and auto-provisions a `test` realm. See [Deployment](deployment.md) for details.

$ rg -n 'func \(s \*Service\) (ListInterfaceMajors|GetInterface|UpdateInterface|DeleteInterface)|errMajorNotFound' internal/realm/service.go
48:// errMajorNotFound marks a PUT/DELETE lookup miss on the (name, major) pair:
52:var errMajorNotFound = errors.New("realm: interface major not found")
191:func (s *Service) UpdateInterface(ctx context.Context, realm, urlName string, urlMajor int, def []byte) (*store.StoredInterface, error) {
215:			return nil, fmt.Errorf("%w: %w", errMajorNotFound, err)
254:// a lookup miss is marked errMajorNotFound for the HTTP layer (#62).
255:func (s *Service) DeleteInterface(ctx context.Context, realm, name string, major int) error {
262:			return fmt.Errorf("%w: %w", errMajorNotFound, err)
364:func (s *Service) ListInterfaceMajors(ctx context.Context, realm, name string) ([]int, error) {
388:func (s *Service) GetInterface(ctx context.Context, realm, name string, major int) (json.RawMessage, error) {

→ Read internal/realm/service.go [limit=50, offset=364]

$ rg -n 'json:"' internal/engine/triggers/*.go | rg -v '_test' | rg -i 'trigger|on|type|policy|action|simple|match|value|group|device|interface|url|method|header|ssl|template' | head -60
internal/engine/triggers/actions.go:49:	HTTPURL           string            `json:"http_url"`
internal/engine/triggers/actions.go:50:	HTTPMethod        string            `json:"http_method"`
internal/engine/triggers/actions.go:51:	HTTPPostURL       string            `json:"http_post_url"` // pre-1.1 legacy: implies POST
internal/engine/triggers/actions.go:52:	HTTPStaticHeaders map[string]string `json:"http_static_headers"`
internal/engine/triggers/actions.go:53:	IgnoreSSLErrors   bool              `json:"ignore_ssl_errors"`
internal/engine/triggers/actions.go:54:	Template          string            `json:"template"`
internal/engine/triggers/actions.go:55:	TemplateType      string            `json:"template_type"`
internal/engine/triggers/actions.go:141:			AMQPExchange string `json:"amqp_exchange"`
internal/engine/triggers/events.go:47:		Timestamp   string `json:"timestamp"`
internal/engine/triggers/events.go:48:		DeviceID    string `json:"device_id"`
internal/engine/triggers/events.go:49:		Event       any    `json:"event"`
internal/engine/triggers/events.go:50:		TriggerName string `json:"trigger_name"`
internal/engine/triggers/events.go:65:	Type string `json:"type"`
internal/engine/triggers/events.go:67:	Interface string `json:"interface"`
internal/engine/triggers/events.go:69:	Path string `json:"path"`
internal/engine/triggers/events.go:71:	Value any `json:"value"`
internal/engine/triggers/events.go:85:	Type string `json:"type"`
internal/engine/triggers/events.go:87:	Interface string `json:"interface"`
internal/engine/triggers/events.go:89:	Path string `json:"path"`
internal/engine/triggers/events.go:91:	OldValue any `json:"old_value"`
internal/engine/triggers/events.go:93:	NewValue any `json:"new_value"`
internal/engine/triggers/events.go:110:	Type string `json:"type"`
internal/engine/triggers/events.go:112:	Interface string `json:"interface"`
internal/engine/triggers/events.go:114:	Path string `json:"path"`
internal/engine/triggers/events.go:116:	OldValue any `json:"old_value"`
internal/engine/triggers/events.go:118:	NewValue any `json:"new_value"`
internal/engine/triggers/events.go:133:	Type string `json:"type"`
internal/engine/triggers/events.go:135:	Interface string `json:"interface"`
internal/engine/triggers/events.go:137:	Path string `json:"path"`
internal/engine/triggers/events.go:139:	Value any `json:"value"`
internal/engine/triggers/events.go:151:	Type string `json:"type"`
internal/engine/triggers/events.go:153:	Interface string `json:"interface"`
internal/engine/triggers/events.go:155:	Path string `json:"path"`
internal/engine/triggers/events.go:167:	Type string `json:"type"`
internal/engine/triggers/events.go:169:	Interface string `json:"interface"`
internal/engine/triggers/events.go:171:	Path string `json:"path"`
internal/engine/triggers/events.go:173:	Value any `json:"value"`
internal/engine/triggers/events.go:184:	Type string `json:"type"`
internal/engine/triggers/events.go:195:	Type string `json:"type"`
internal/engine/triggers/events.go:197:	DeviceIPAddress string `json:"device_ip_address"`
internal/engine/triggers/events.go:208:	Type string `json:"type"`
internal/engine/triggers/events.go:219:	Type string `json:"type"`
internal/engine/triggers/events.go:230:	Type string `json:"type"`
internal/engine/triggers/events.go:241:	Type string `json:"type"`
internal/engine/triggers/events.go:244:	ErrorName string `json:"error_name"`
internal/engine/triggers/events.go:246:	Metadata map[string]string `json:"metadata"`
internal/engine/triggers/events.go:273:	Type string `json:"type"`
internal/engine/triggers/events.go:275:	Introspection string `json:"introspection"`
internal/engine/triggers/events.go:286:	Type string `json:"type"`
internal/engine/triggers/events.go:288:	Interface string `json:"interface"`
internal/engine/triggers/events.go:290:	MajorVersion int `json:"major_version"`
internal/engine/triggers/events.go:291:	MinorVersion int `json:"minor_version"`
internal/engine/triggers/events.go:302:	Type string `json:"type"`
internal/engine/triggers/events.go:304:	Interface string `json:"interface"`
internal/engine/triggers/events.go:306:	MajorVersion int `json:"major_version"`
internal/engine/triggers/policy.go:99:		Name          string `json:"name"`
internal/engine/triggers/policy.go:101:			On       json.RawMessage `json:"on"`
internal/engine/triggers/policy.go:102:			Strategy string          `json:"strategy"`
internal/engine/triggers/policy.go:103:		} `json:"error_handlers"`
internal/engine/triggers/policy.go:104:		MaximumCapacity int  `json:"maximum_capacity"`

→ Read internal/engine/triggers/actions.go [limit=160, offset=40]

$ rg -n 'ownership|aggregation|required|missing' pkg/interfaceschema/parse.go | head -40
94:	// enforced at runtime only for object aggregation (upstream 1.4
115:	Ownership     *Ownership     `json:"ownership,omitempty"`
116:	// Quality is the upstream legacy alias of ownership (issue #61): accepted
120:	// Aggregate is the upstream legacy alias of aggregation (issue #61).
122:	Aggregation *Aggregation  `json:"aggregation,omitempty"`
147:	Required                *bool                    `json:"required,omitempty"`
156:// aggregation and per-type field constraints) are rejected with an error
213:	// accepted only without ownership (both → violation on ownership);
214:	// aggregate only without aggregation (both → violation on aggregation).
218:			vc.add("ownership", "ownership and quality are mutually exclusive")
227:			vc.add("aggregation", "aggregation and aggregate are mutually exclusive")
240:		return nil, errors.New(`missing "interface_name"`)
242:		return nil, errors.New(`missing "version_major"`)
244:		return nil, errors.New(`missing "version_minor"`)
246:		return nil, errors.New(`missing "type"`)
248:		return nil, errors.New(`missing "ownership"`)
284:		return nil, errors.New("properties interfaces cannot use object aggregation")
354:		return m, nil, errors.New(`mapping is missing "endpoint"`)
357:		return m, nil, fmt.Errorf(`mapping %q is missing "type"`, *raw.Endpoint)
396:			return m, nil, fmt.Errorf(`mapping %q: "required" is not allowed on properties`, m.Endpoint)
595:// checkObjectAggregation enforces the object-aggregation shape: every

$ rg -n 'func \(s \*Service\) DeleteTrigger|func \(s \*Service\) DeletePolicy|func \(s \*Service\) CreateTrigger|func \(s \*Service\) CreatePolicy' internal/realm/service.go
419:func (s *Service) CreateTrigger(ctx context.Context, realm string, def []byte) (*store.Trigger, error) {
473:func (s *Service) DeleteTrigger(ctx context.Context, realm, name string) error {
515:func (s *Service) CreatePolicy(ctx context.Context, realm string, def []byte) (*store.TriggerPolicy, error) {
559:func (s *Service) DeletePolicy(ctx context.Context, realm, name string) error {

→ Read internal/realm/service.go [limit=80, offset=419]

$ ls docs/site/ && echo "---" && rg -ln 'realmmanagement|realm-management|Realm Management|1\.1\.0|Interface major' docs/site/ 2>/dev/null
DESIGN.md
ROADMAP.md
api
appengine-api.md
architecture.md
compatibility.md
configuration-reference.md
contributing.md
data-modeling.md
deployment.md
housekeeping-api.md
index.md
interface-schema.md
json-payload-profile.md
migration-from-astarte.md
mqtt-protocol.md
observability.md
operations.md
pairing-and-security.md
payload-formats.md
pokemon-agent.md
quickstart.md
realm-management-api.md
swagger-ui
swagger.md
triggers.md
troubleshooting.md
---
docs/site/pairing-and-security.md
docs/site/architecture.md
docs/site/deployment.md
docs/site/swagger.md
docs/site/index.md
docs/site/triggers.md
docs/site/operations.md
docs/site/interface-schema.md
docs/site/realm-management-api.md
docs/site/migration-from-astarte.md
docs/site/compatibility.md

$ rg -n '1\.1\.0|Interface major|major version|404|Not Found|not found' docs/site/realm-management-api.md | head -40
(no output)

→ Read docs/site/realm-management-api.md
→ Read internal/realm/service.go [limit=40, offset=620]

$ rg -n 'triggers/\{name\}|Update trigger|updateTrigger|PUT.*triggers' internal/ cmd/ docs/api/ | head
docs/api/astarte_realm_management_api.yaml:449:  /realmmanagement/v1/{realm}/triggers/{name}:
internal/realm/http.go:46:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers/{name}", h(a.getTrigger))
internal/realm/http.go:47:	mux.Handle("DELETE /realmmanagement/v1/{realm}/triggers/{name}", h(a.deleteTrigger))

→ Read internal/realm/service.go [limit=50, offset=253]

$ rg -n 'func.*DeleteInterface|ErrInterfaceMajorNotZero|ErrInterfaceInUse' internal/store/*.go
internal/store/interfaces_test.go:189:		if !errors.Is(err, ErrInterfaceMajorNotZero) {
internal/store/interfaces_test.go:190:			t.Fatalf("delete of major 1: got %v, want ErrInterfaceMajorNotZero", err)
internal/store/interfaces_test.go:210:		if !errors.Is(err, ErrInterfaceInUse) {
internal/store/interfaces_test.go:211:			t.Fatalf("delete of declared interface: got %v, want ErrInterfaceInUse", err)
internal/store/store.go:48:	// ErrInterfaceInUse reports an interface delete blocked because some
internal/store/store.go:50:	ErrInterfaceInUse = errors.New("store: interface is referenced by device introspection")
internal/store/store.go:51:	// ErrInterfaceMajorNotZero reports an interface delete blocked by the
internal/store/store.go:53:	ErrInterfaceMajorNotZero = errors.New("store: only major version 0 interfaces can be deleted")
internal/store/interfaces.go:190:func (s *Store) DeleteInterface(ctx context.Context, realmID int16, name string, major int) error {
internal/store/interfaces.go:192:		return fmt.Errorf("%w: interface %s v%d", ErrInterfaceMajorNotZero, name, major)
internal/store/interfaces.go:225:		return fmt.Errorf("%w: interface %s v%d", ErrInterfaceInUse, name, major)

→ Read internal/store/interfaces.go [limit=45, offset=190]

$ rg -n 'func Test.*(RealmManagement|RealmName|Interface404|NotFound|Conflict)' internal/swagger/swagger_test.go
131:func TestRealmManagement403(t *testing.T) {
249:func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
603:func TestRealmManagementErrorDetailExamples(t *testing.T) {
1088:func TestRealmManagement422ValidationDetails(t *testing.T) {
1285:func TestRealmManagementDeviceIDEncodingDocumented(t *testing.T) {
1364:func TestHousekeepingRealmNamePatternDocumented(t *testing.T) {
1470:func TestPairingRealmNamePatternDocumented(t *testing.T) {

→ Read internal/swagger/swagger_test.go [limit=55, offset=603]

$ tail -c 1200 .mule/todo.md; echo; echo "=== FORGIULIO TAIL ==="; tail -c 800 .mule/for-giulio.md
2-425) — its 422/409/404/default arms including that it maps the two interface sentinels to 422, not 403; (3) `renderViolationsBody` (http.go:443-533) byte-for-byte for a mixed input (interface-level field merged first-seen + the full-length index-aligned `mappings` array with `{}` for clean entries); (4) `validationDetail` — prefix stripped, non-prefix passthrough, exact-prefix NOT stripped. `*API{}` needs no Service (these methods never deref `a.svc`), and `interfaceschema.ViolationsError` is built from its exported `Violations`/`MappingCount` fields, so no store or Docker.
- [x] realm-detailed-properties-branch-test [auto]: extend `internal/realm/detailed_required_encrypted_test.go` to cover the properties branch of `writeDetailedMapping` (internal/realm/detailed.go:61-63) — currently only the datastream branch is pinned container-free (the properties path runs only in the integration-tagged `internal/realm/detailed_test.go`): assert a properties mapping renders exactly `{"endpoint":..,"type":..,"allow_unset":<bool>}` with no reliability/retention/expiry/explicit_timestamp/required/encrypted/database_retention_policy keys, for `allow_unset` both true and false. No Docker.

=== FORGIULIO TAIL ===
ng else stale: the two
`jwt_public_key*` keys the page quotes at lines 69-71 are real (`config.go:112-113`), and the
`master_key_file` row's `ASTRATE_MASTER_KEY`/`ASTRATE_MASTER_KEY_FILE` fallbacks are real too
(`internal/store/crypto.go:24-27`). And the one thing that looked like a live bug is not one:
`view()` returns only `JWTPublicKeysPEM[0]` (`internal/housekeeping/service.go:277-281`) while
the realm-management surface joins *all* keys (`internal/realm/service.go:627`) — but every
writer stores a single-element array (create at `internal/store/realms.go:54-61`, PATCH at
`internal/store/realms.go:165`, `putAuthConfig` at `internal/realm/service.go:636`), so a
two-key realm is unreachable and there is nothing to document. Worth noting because it reads
like a bug on first inspection.

→ Read internal/auth/middleware.go [limit=100, offset=48]

$ rg -n "pattern: '\^\[a-z\]\[a-z0-9\]\*\\\$'" docs/api/*.yaml
docs/api/astrate_native_api.yaml:420:            pattern: '^[a-z][a-z0-9]*$'
docs/api/astrate_native_api.yaml:467:        pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_housekeeping_api.yaml:246:        pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_housekeeping_api.yaml:280:            pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_housekeeping_api.yaml:289:          pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_housekeeping_api.yaml:311:          pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_pairing_api.yaml:384:        pattern: '^[a-z][a-z0-9]*$'

✗ Read .mule/todo.md failed [offset=322]
Error: Offset 322 is out of range for this file (304 lines)
→ Read .mule/todo.md [offset=302]
