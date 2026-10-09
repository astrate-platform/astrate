slug: recipe-docs-sync
verdict: proposed
at:  a086f26
ran: 2026-10-09T11:29:51Z on DietPi in 762s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls docs/api/ && echo "---SITE---" && ls docs/site/ 2>/dev/null | head -40 && echo "---MAKEFILE---" && ls docs/
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
---SITE---
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
---MAKEFILE---
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

$ ls -la .mule/ && echo "---TODO---" && cat .mule/todo.md && echo "---FORGIULIO---" && cat .mule/for-giulio.md 2>/dev/null | tail -50
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_1206203ad001Kcx5lad0oDAniT

- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` shows `detail: Bad Request` (yaml:1732) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1792), while the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58 via `WriteBadRequest`/`WriteInternalServerError`, golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json); every 400/500 on that surface goes through the same astarteapi writers. The other detail examples on the file are already correct — change only those two — and add the docs-half assertion `TestAppEngineErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go, which pins the pairing and realm-management halves. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDetail` schema-level example (yaml:361-362) and the `BadRequest` response (yaml:385) show `detail: Bad Request`, and `InternalServerError` (yaml:477) shows `detail: Internal Server Error`, but the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58, via WriteBadRequest/WriteInternalServerError, pinned by pkg/astarteapi/testdata/error_bad_request.json and error_internal.json). While you are in `ErrorDetail`: its `example` sits at the schema level and shows a bare `{detail: ...}` where the envelope the wire sends is `{"errors":{"detail":...}}` — move it under `errors` (or drop it, as the realm-management and pairing `ErrorDetail` schemas carry no example) so a generated client is not shown the wrong shape. Add `TestHousekeepingErrorDetailExamples` beside the pairing and realm-management halves in internal/swagger/swagger_test.go. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `"Internal server error"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implements `device_empty_cache_received` end to end — `SimpleEvent` field 23, an empty `DeviceEmptyCacheReceivedEvent` whose encoder emits exactly `{"type":"device_empty_cache_received"}` (astarte_core encoder.ex:54-62), fired by DUP at `control_handler.ex:161` *inside* the `with` that runs `send_control_consumer_properties` → `resend_all_properties` → `set_pending_empty_cache`, so it fires only when all three succeed and never on failure, and the commit deletes the `# TODO: implement on_empty_cache_received` from `astarte_events/lib/astarte_events/triggers/core.ex`. Astrate already accepts the trigger but never evaluates it: `OnDeviceEmptyCacheReceived` sits in `deviceOns` with value `false` (internal/engine/triggers/match.go:123, the `evaluated` flag read at match.go:406) and the comment's stated reason ("upstream defines no SimpleEvent variant for it") is now false. Add the event body constructor to internal/engine/triggers/events.go (the payload is `{"type":"device_empty_cache_received"}` and nothing else — no device id, no metadata), flip the `deviceOns` value to true, and fire it from `handleEmptyCache` (internal/engine/control.go:58) with the same idiom `handleLifecycle` uses for `device_connected` (`e.fireDevice` + `stream.Publish`, internal/engine/engine.go:346-353), at the point where `m.Ack()` runs — i.e. after `resendServerProperties` and `sendConsumerProperties` have both succeeded inside `retryStore`. Add a T1 test that a trigger on it fires exactly once per accepted emptyCache and does not fire when the resend fails. — BLOCKED: wrote nothing
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
---FORGIULIO---
**1. `docs/site/housekeeping-api.md:30` still talks about Cassandra.**

> Cassandra-specific fields (`replication_class`, `replication_factor`, etc.) from upstream are accepted but ignored.

Astrate stores PostgreSQL — `internal/store` on pgx, schema in `migrations/*.sql`; there is no
Cassandra anywhere in the tree. The wire shape is four fields (`internal/housekeeping/http.go:46-51`),
and the comment right above it says the opposite of "accepted": *"Astrate omits the
Cassandra-specific fields (replication factor/class) upstream carries"*
(`internal/housekeeping/http.go:44-45`). From a client's side the two readings are the same
(the request succeeds, the fields do nothing), so this is vocabulary rather than behaviour —
but it is the kind of sentence that sends an operator hunting for a Cassandra setting that
cannot exist.

**2. `PATCH /housekeeping/v1/realms/{realm}` is missing from the page entirely.**
`docs/site/housekeeping-api.md` goes Create → List → Get → Delete (lines 10-50), so the route
`docs-sync-hk-patch-endpoint` added and documented in the spec
(`docs/api/astarte_housekeeping_api.yaml:142-192`) has no prose at all. The same page also
predates `datastream_maximum_storage_retention`, which the spec has carried on both `Realm` and
`RealmCreate` since `docs-sync-hk-retention-field`, so its create example (lines 12-21) shows a
three-field body that is not the whole body.

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

→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ ls -la .mule/reports/ 2>/dev/null | tail -40; echo "---"; ls .mule/
total 64
drwxr-xr-x 3 root root  4096 Sep  6 19:40 .
drwxr-xr-x 7 root root  4096 Oct  9 12:56 ..
-rw-r--r-- 1 root root  2279 Sep  6 19:40 hygiene.md
drwxr-xr-x 2 root root 49152 Oct  9 12:58 log
---
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
recipes
reports
reviews
tasks
todo.md

$ git log --oneline -30 2>/dev/null; echo "---BRANCH---"; git branch --show-current
a086f26 mule: recipe code-review ran (136s)
705c467 mule: recipe astarte-upstream ran (108s)
1ca603c mule: refill the queue
26c12e4 mule: recipe github-issues ran (155s)
271423d mule: recipe milestones ran (574s)
c1b43cd mule: recipe hygiene timed out
e30b5fd mule: log docs-sync-pairing-deviceid-path-param
1cda0fb mule: docs-sync-pairing-deviceid-path-param [auto]: give `components.parameters.DeviceID` the constraints `RegisterRequest.hw_id` already carries in docs/api/astarte_pairing_api.yaml — the parameter (yaml:367-377) is prose-only, while `hw_id` (yaml:386-398) got `pattern: '^[A-Za-z0-9_-]{22}$'` plus `minLength`/`maxLength` from docs-sync-pairing-deviceid-base64url, and `TestPairingDeviceIDEncodingDocumented` (internal/swagger/swagger_test.go:1086) pins only the two descriptions and hw_id's bounds, never the parameter's schema; the realm-management twin was constrained later by docs-sync-rm-deviceid-param and its test (swagger_test.go:1168) does hold the parameter's pattern against `deviceid.Parse`. Add the same three schema lines to the parameter, and record in its description the asymmetry the prose omits: a malformed `deviceID` is **401** on all three device-authenticated operations — `deviceid.Parse` failing inside `authenticateDevice` burns bcrypt and returns `ErrUnauthorized` (internal/pairing/service.go:381-385) → `WriteUnauthorized` (http.go:322-323), exactly as an unknown device or a wrong secret is (service.go:395-397), so the client is told its credentials are bad rather than its ID malformed — but **404** on the agent `unregisterDevice` (service.go:231-234 → http.go:326-327). Extend TestPairingDeviceIDEncodingDocumented to compile the new pattern and hold it against `deviceid.Parse` the way the realm-management test does. Do not encode the strict trailing-bits rule in the pattern; that stays open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
555b6e8 mule: log docs-sync-pairing-realm-name-pattern
ed301fb mule: docs-sync-pairing-realm-name-pattern [auto]: constrain `components.parameters.RealmName` in docs/api/astarte_pairing_api.yaml (yaml:359-365, a bare `type: string` described only as "The realm name.") with `pattern: '^[a-z][a-z0-9]*$'` — the realm name is CHECK-constrained to exactly that (migrations/000002_metadata.up.sql:7, read-only, never edit) and the two surfaces whose tasks already ran carry it on their twins (docs/api/astarte_housekeeping_api.yaml:246, docs/api/astrate_native_api.yaml:420 and 467). Because this one parameter is shared by all seven operations, the description must record that an off-pattern name gets three different answers, not one: the agent and device routes answer **401** (the realm is resolved by the middleware, which treats a miss as unauthenticated rather than 404 — internal/auth/middleware.go:66-72); `getHealth` answers **404** (service.Health returns store.ErrNotFound → WriteNotFound, internal/pairing/service.go:369-371, http.go:289-290); and `getVersion` answers **200** `{"data":"1.2.2"}` for any realm string whatsoever, because `GET /pairing/v1/{realm}/version` is registered with no middleware and no realm lookup (cmd/astrate/main.go:438 → observability.VersionHandler, internal/observability/compat.go:38-47) — say that in the version operation so nobody reads the shared parameter as a guarantee the handler enforces. Mirror `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go:1247) with a pairing twin only if it stays a few lines; the YAML edit is the deliverable. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
bf5aadb mule: blocked docs-sync-pairing-body-cap-400
d7e14cf mule: log docs-sync-pairing-register-404-unreachable
cfa6e76 mule: docs-sync-pairing-register-404-unreachable [auto]: drop the `404` `DeviceNotFound` response from `registerDevice` in docs/api/astarte_pairing_api.yaml (yaml:140-141) — it is documented but absent: the route is wrapped by `requireAgent` (internal/pairing/http.go:82-83), i.e. `mw.RequireRealm(auth.ClaimPairing)`, which resolves the realm *before* the handler and answers `401 {"errors":{"detail":"Unauthorized"}}` for an unknown one (internal/auth/middleware.go:66-72, "no existence oracle on auth failures"), so `Service.Register`'s own `GetRealmByName` (internal/pairing/service.go:183-186) can return `store.ErrNotFound` only if the realm is deleted between the middleware's lookup and the handler's — a race, not a contract. Nothing else in Register returns `store.ErrNotFound`: `RegisterDevice`'s `ErrDeviceAlreadyConfirmed` becomes 422 (service.go:209-213) and the CountDevices/SetPayloadFormatHint failures are DB errors → 500. No test covers it (internal/pairing/http_test.go asserts 404 only for the agent DELETE at :527 and for health at :595-598), and the `DeviceNotFound` component's own description — "unknown device ID" — does not match what that path would be even if it fired. The premise came from docs-sync-pairing-register-404, which read service.go without the middleware in front of it. Keep the `404` on `unregisterDevice` (reachable: unknown device, and a malformed `deviceID` via service.go:231-234) and on `getHealth`. State the reachability argument in the report; if you find a reachable path instead, keep the response and fix its description rather than deleting it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
e06b1d8 mule: refill the queue
1753d91 mule: recipe docs-sync ran (693s)
59c62e8 mule: log engine-introspection-producer-roundtrip
d070345 mule: engine-introspection-producer-roundtrip [auto]: pin the producer half of the docs/DESIGN.md §3.3 introspection rule on the default gate — `testutil.Introspection` (internal/testutil/astartedevice.go:93-105) is called from six places, all behind `e2e`, `integration && e2e` or test/conformance build tags, while the consumer `parseIntrospection` (internal/engine/introspection.go:99-127) is tested only against hand-written literals in internal/engine/introspection_test.go:17-71, so nothing container-free checks that the string the harness renders is the one the engine accepts or that the documented "deterministic (sorted) ordering" holds. Add a case to internal/engine/introspection_test.go that builds the string from an unsorted multi-entry map via `testutil.Introspection`, asserts `parseIntrospection` returns the same map, and pins the exact output (sorted, `name:major:minor;…`, no trailing `;`, empty map → empty string).
4693262 mule: blocked testutil-control-frame-test
80ae7bc mule: log testutil-wait-since-cursor
3654fda mule: testutil-wait-since-cursor [auto]: give AstarteDevice's message waits a cursor so a second wait for the same topic cannot re-match an older message — `WaitForMessage`/`WaitForTopic` (internal/testutil/astartedevice.go:149-170) scan the capture buffer from index 0 on every poll, so a repeated wait returns the first match again, and in-tree that makes testE2EEmptyCache's resync wait vacuous (see the companion line). Move the paho capture handler out of the `ConnectAstarteDevice` closure (astartedevice.go:61-67) into an injectable method, add `Mark() int` and a `From(start int, …)` wait variant, and add a container-free internal/testutil/astartedevice_test.go that drives two captures on one topic and asserts a wait started after the second returns the second message — it cannot even compile today, which is the failure proof. Keep every new read of `received` under the existing `d.mu`, and say in the report that the change is unverified for races (the Pi gate runs no `-race`). The consumer of this API is engine-e2e-emptycache-resync-assert.
fcdeb2e mule: refill the queue
1fe59e3 mule: recipe code-review ran (822s)
6df81bc mule: recipe astarte-upstream ran (156s)
d0cab66 mule: refill the queue
bae4595 mule: recipe github-issues ran (463s)
0d2e197 mule: recipe milestones ran (824s)
0f9c65c mule: recipe hygiene timed out
d8124d6 mule: refill the queue
10efb94 mule: recipe docs-sync ran (518s)
ebd5fb0 mule: log pairing-unregister-bad-id-test
---BRANCH---
mule/queue

$ ls -lt .mule/reports/log/ | head -30
total 14828
-rw-r--r-- 1 root root  67075 Oct  9 12:58 20261009T105827Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root  89855 Oct  8 22:42 20261008T204257Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  61863 Oct  8 22:38 20261008T203843Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root 186642 Oct  8 22:21 20261008T202151Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  99670 Oct  8 22:06 20261008T200623Z-recipe-hygiene-timeout.md
-rw-r--r-- 1 root root  27083 Oct  8 21:36 20261008T193604Z-docs-sync-pairing-deviceid-path-param-done.md
-rw-r--r-- 1 root root  23561 Oct  8 21:28 20261008T192820Z-docs-sync-pairing-realm-name-pattern-done.md
-rw-r--r-- 1 root root  14077 Oct  8 21:13 20261008T191310Z-docs-sync-pairing-body-cap-400-blocked.md
-rw-r--r-- 1 root root   3621 Oct  8 20:31 20261008T183153Z-docs-sync-pairing-register-404-unreachable-done.md
-rw-r--r-- 1 root root 134121 Oct  8 20:26 20261008T182643Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root   6525 Oct  8 19:58 20261008T175838Z-engine-introspection-producer-roundtrip-done.md
-rw-r--r-- 1 root root   3972 Oct  8 19:47 20261008T174752Z-testutil-control-frame-test-blocked.md
-rw-r--r-- 1 root root  12117 Oct  8 19:37 20261008T173721Z-testutil-wait-since-cursor-done.md
-rw-r--r-- 1 root root 159555 Oct  8 19:28 20261008T172851Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root  55153 Oct  8 13:39 20261008T113945Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root 151100 Oct  8 13:10 20261008T111054Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  42189 Oct  8 00:44 20261007T224444Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root   2634 Oct  8 00:30 20261007T223040Z-recipe-hygiene-timeout.md
-rw-r--r-- 1 root root 123125 Oct  7 21:18 20261007T191848Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root   6767 Oct  7 20:47 20261007T184717Z-pairing-unregister-bad-id-test-done.md
-rw-r--r-- 1 root root  11653 Oct  7 20:08 20261007T180823Z-pairing-onregistered-clock-done.md
-rw-r--r-- 1 root root  10631 Oct  7 20:02 20261007T180202Z-pairing-burn-bcrypt-cost-blocked.md
-rw-r--r-- 1 root root  66419 Oct  7 19:44 20261007T174416Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root 172374 Oct  7 19:18 20261007T171801Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root 156608 Oct  7 18:58 20261007T165810Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  75616 Oct  7 13:38 20261007T113844Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root   2700 Oct  7 13:25 20261007T112533Z-recipe-hygiene-timeout.md
-rw-r--r-- 1 root root  42462 Oct  6 22:49 20261006T204912Z-docs-native-socket-event-schema-blocked.md
-rw-r--r-- 1 root root  41769 Oct  6 22:29 20261006T202927Z-docs-native-socket-security-scheme-done.md

$ grep -l "docs-sync" .mule/reports/log/* 2>/dev/null | head; echo "---"; ls .mule/reports/log/ | grep -i docs
.mule/reports/log/20260905T165512Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260905T171755Z-docs-sync-pairing-health-path-done.md
.mule/reports/log/20260905T172234Z-docs-sync-pairing-register-404-done.md
.mule/reports/log/20260905T173657Z-recipe-hygiene-proposed.md
.mule/reports/log/20260905T181053Z-recipe-milestones-proposed.md
.mule/reports/log/20260905T184625Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260905T192010Z-recipe-code-review-proposed.md
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260906T113143Z-docs-sync-rm-datastream-retention-endpoint-done.md
.mule/reports/log/20260906T171045Z-docs-sync-rm-interfaces-detailed-param-done.md
---
20260905T165512Z-recipe-docs-sync-proposed.md
20260905T171755Z-docs-sync-pairing-health-path-done.md
20260905T172234Z-docs-sync-pairing-register-404-done.md
20260906T112006Z-recipe-docs-sync-proposed.md
20260906T113143Z-docs-sync-rm-datastream-retention-endpoint-done.md
20260906T171045Z-docs-sync-rm-interfaces-detailed-param-done.md
20260906T200431Z-recipe-docs-sync-proposed.md
20260906T201002Z-docs-sync-appengine-by-alias-endpoints-done.md
20260906T202744Z-docs-sync-appengine-group-endpoints-done.md
20260906T204333Z-docs-sync-appengine-get-group-device-done.md
20260907T112447Z-docs-sync-appengine-query-params-status-done.md
20260907T112654Z-docs-sync-appengine-group-patch-status-done.md
20260907T170024Z-docs-sync-appengine-data-422-interface-level-done.md
20260907T191157Z-recipe-docs-sync-proposed.md
20260907T201232Z-docs-sync-appengine-data-output-params-done.md
20260908T182108Z-recipe-docs-sync-proposed.md
20260908T183808Z-docs-sync-appengine-device-status-schema-done.md
20260908T184535Z-docs-sync-appengine-data-set-422-blocked.md
20260908T185820Z-docs-sync-appengine-downsample-min-done.md
20260909T165735Z-recipe-docs-sync-proposed.md
20260909T170637Z-docs-sync-hk-patch-endpoint-done.md
20260909T171614Z-docs-sync-hk-retention-field-done.md
20260909T200459Z-recipe-docs-sync-proposed.md
20260909T201715Z-docs-sync-native-compat-health-503-done.md
20260909T202852Z-docs-sync-native-compat-version-endpoints-done.md
20260910T184706Z-recipe-docs-sync-proposed.md
20260910T190522Z-docs-sync-pairing-status-enum-done.md
20260910T191129Z-docs-sync-pairing-version-endpoint-done.md
20260911T110803Z-recipe-docs-sync-proposed.md
20260911T112356Z-docs-sync-rm-delete-interface-status-done.md
20260911T113809Z-docs-sync-rm-mapping-required-encrypted-done.md
20260911T170716Z-docs-sync-rm-put-auth-422-done.md
20260911T182955Z-docs-sync-rm-version-example-blocked.md
20260912T110449Z-recipe-docs-sync-proposed.md
20260912T112234Z-docs-sync-ae-write-405-done.md
20260912T113032Z-docs-sync-ae-write-value-422-done.md
20260912T164501Z-docs-sync-ae-post-groups-409-done.md
20260912T192246Z-recipe-docs-sync-proposed.md
20260912T192341Z-recipe-docs-sync-proposed.md
20260912T203456Z-docs-sync-ae-post-group-devices-409-done.md
20260912T204618Z-docs-sync-ae-patch-merge-patch-media-type-done.md
20260913T182321Z-recipe-docs-sync-proposed.md
20260913T185042Z-docs-sync-hk-delete-gating-responses-done.md
20260913T191910Z-docs-sync-hk-validation-example-done.md
20260914T174510Z-recipe-docs-sync-proposed.md
20260915T180720Z-recipe-docs-sync-proposed.md
20260915T181748Z-docs-sync-native-socket-query-token-auth-done.md
20260916T115527Z-recipe-docs-sync-proposed.md
20260916T165431Z-docs-sync-pairing-deviceendpoints-dead-404-403-done.md
20260916T171527Z-docs-sync-pairing-version-404-unreachable-done.md
20260916T204810Z-recipe-docs-sync-proposed.md
20260917T111017Z-docs-sync-pairing-version-value-done.md
20260917T114541Z-docs-sync-native-version-value-done.md
20260917T194954Z-recipe-docs-sync-proposed.md
20260917T200755Z-docs-sync-rm-policies-delete-422-done.md
20260917T202137Z-docs-sync-rm-triggers-422-nested-envelope-done.md
20260918T180119Z-recipe-docs-sync-proposed.md
20260919T113513Z-recipe-docs-sync-proposed.md
20260920T185956Z-recipe-docs-sync-proposed.md
20260920T191642Z-docs-sync-native-metrics-example-fake-series-done.md
20260920T192702Z-docs-sync-native-socket-missing-403-500-done.md
20260921T174646Z-recipe-docs-sync-proposed.md
20260921T203450Z-recipe-docs-sync-proposed.md
20260922T105805Z-docs-sync-ae-read-query-params-blocked.md
20260922T173055Z-recipe-docs-sync-proposed.md
20260922T200147Z-recipe-docs-sync-proposed.md
20260923T112816Z-recipe-docs-sync-proposed.md
20260923T201015Z-recipe-docs-sync-proposed.md
20260923T201746Z-docs-sync-rm-put-interface-409-done.md
20260924T110303Z-docs-sync-rm-interface-422-shapes-done.md
20260924T111049Z-docs-sync-rm-validation-example-prefix-done.md
20260924T165605Z-docs-sync-rm-auth-403-done.md
20260924T185047Z-recipe-docs-sync-proposed.md
20260924T202134Z-recipe-docs-sync-proposed.md
20260925T192553Z-recipe-docs-sync-proposed.md
20260925T193504Z-docs-sync-hk-async-operation-param-done.md
20260925T195121Z-docs-sync-hk-retention-zero-is-unset-done.md
20260925T200001Z-docs-sync-rm-async-operation-param-done.md
20260925T201910Z-docs-sync-rm-delete-device-async-operation-param-done.md
20260925T204521Z-docs-sync-rm-delete-device-async-operation-param-blocked.md
20260926T194649Z-recipe-docs-sync-proposed.md
20260926T203415Z-docs-sync-ae-forbidden-403-done.md
20260927T111051Z-docs-sync-ae-patch-by-alias-409-done.md
20260927T111624Z-docs-sync-ae-add-group-device-422-done.md
20260927T113516Z-docs-sync-ae-delete-data-400-done.md
20260927T164857Z-docs-sync-ae-list-devices-422-transient.md
20260927T171520Z-docs-sync-ae-list-devices-422-transient.md
20260927T182324Z-docs-sync-ae-list-devices-422-blocked.md
20260928T105915Z-recipe-docs-sync-proposed.md
20260928T111752Z-docs-sync-pairing-error-example-capitalisation-done.md
20260928T113726Z-docs-sync-pairing-info-version-example-done.md
20260928T170059Z-docs-sync-pairing-deviceid-base64url-done.md
20260928T172627Z-docs-sync-pairing-initial-payload-format-enum-done.md
20260928T180120Z-docs-sync-pairing-unregister-description-done.md
20260928T204651Z-recipe-docs-sync-proposed.md
20260929T170954Z-recipe-docs-sync-proposed.md
20260929T183752Z-recipe-docs-sync-proposed.md
20260930T113453Z-recipe-docs-sync-proposed.md
20260930T175957Z-recipe-docs-sync-proposed.md
20260930T192852Z-recipe-docs-sync-proposed.md
20261001T164758Z-recipe-docs-sync-proposed.md
20261001T185956Z-recipe-docs-sync-proposed.md
20261002T181712Z-recipe-docs-sync-proposed.md
20261002T183434Z-docs-sync-rm-error-example-capitalisation-done.md
20261002T192357Z-docs-sync-rm-legacy-alias-fields-blocked.md
20261002T195225Z-docs-sync-rm-validationerror-example-done.md
20261002T201307Z-docs-sync-rm-deviceid-param-done.md
20261002T203814Z-docs-sync-rm-update-interface-body-blocked.md
20261003T110748Z-docs-sync-ae-error-example-capitalisation-done.md
20261003T113021Z-docs-sync-hk-error-example-capitalisation-blocked.md
20261003T165209Z-docs-sync-native-error-example-capitalisation-done.md
20261003T170519Z-docs-sync-hk-error-detail-examples-split-done.md
20261004T112505Z-recipe-docs-sync-proposed.md
20261004T114011Z-docs-sync-hk-realm-name-pattern-done.md
20261004T165229Z-docs-sync-hk-errordetail-schema-example-done.md
20261004T171709Z-docs-sync-hk-patch-422-field-error-examples-blocked.md
20261004T172242Z-docs-sync-hk-wrong-type-field-400-blocked.md
20261004T173828Z-docs-sync-hk-realm-name-response-schemas-done.md
20261005T111650Z-recipe-docs-sync-proposed.md
20261005T113206Z-docs-native-ae-realm-version-endpoint-done.md
20261005T171138Z-docs-native-hk-version-description-done.md
20261005T175822Z-docs-native-metrics-content-negotiation-done.md
20261005T180602Z-docs-native-phoenix-newevent-payload-done.md
20261005T181428Z-docs-native-realm-name-pattern-done.md
20261005T210409Z-recipe-docs-sync-timeout.md
20261006T184046Z-recipe-docs-sync-proposed.md
20261006T185740Z-docs-native-socket-handshake-errors-blocked.md
20261006T201112Z-docs-native-socket-sse-exact-accept-done.md
20261006T202927Z-docs-native-socket-security-scheme-done.md
20261006T204912Z-docs-native-socket-event-schema-blocked.md
20261007T191848Z-recipe-docs-sync-proposed.md
20261008T182643Z-recipe-docs-sync-proposed.md
20261008T183153Z-docs-sync-pairing-register-404-unreachable-done.md
20261008T191310Z-docs-sync-pairing-body-cap-400-blocked.md
20261008T192820Z-docs-sync-pairing-realm-name-pattern-done.md
20261008T193604Z-docs-sync-pairing-deviceid-path-param-done.md

$ head -60 .mule/reports/log/20261008T182643Z-recipe-docs-sync-proposed.md
slug: recipe-docs-sync
verdict: proposed
at:  59c62e8
ran: 2026-10-08T18:26:43Z on DietPi in 693s

> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls docs/api/ && ls internal/
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
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

$ cat .mule/todo.md; echo ---; cat .mule/for-giulio.md 2>/dev/null | tail -40
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_11cba6257001zqXZ3qwLn16ems

- [x] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, housekeeping and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
- [!] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy alias fields to `InterfaceDefinition` and `InterfaceMapping` in docs/api/astarte_realm_management_api.yaml — both declare `additionalProperties: false` (yaml:935, 972) yet omit interface-level `quality` (alias of `ownership`) and `aggregate` (alias of `aggregation`) and mapping-level `path` (alias of `endpoint`), all three of which the parser accepts and canonicalises: decoded into dedicated fields (pkg/interfaceschema/parse.go:116-122, 136-138), normalised before any rule runs (parse.go:212-236), and re-encoded canonically by `ParseInterfaceCanonical` so the store never sees them (parse.go:165-170, 307-311). A spec-conformant client — or any generated client with this schema baked in — therefore rejects a valid upstream-shaped install/update with a spurious 422 while the service would have accepted and stored it as `ownership`/`aggregation`/`endpoint`. Add the three as documented deprecated aliases (enums mirroring ownership/aggregation, `path` a plain string) or relax `additionalProperties`, and record the three rules the code enforces: `ownership` + `quality` is a violation (parse.go:218) and `aggregation` + `aggregate` is a violation (parse.go:227), while `endpoint` + `path` silently keeps `endpoint` (parse.go:307-311). None of the three may appear in a `required` list. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstrained string described only as "The device hardware ID." — `Service.DeleteDevice` runs the path value through `deviceid.Parse` (internal/realm/service.go:118-121), which accepts exactly 22 characters of unpadded base64url in strict mode and rejects padding, the standard alphabet's `+` and `/`, and non-canonical trailing bits (pkg/deviceid/deviceid.go:24-26, 33-37, 42-60); every rejection is folded into `store.ErrNotFound` (service.go:119-121) and answered `404 {"errors":{"detail":"Device not found"}}` (internal/realm/http.go:107-109), so a client that sends a well-formed standard-base64 or 23-character padded id is told the device does not exist rather than that its id is malformed. State the encoding and the 22-character length and add `pattern`/`minLength`/`maxLength`, exactly as docs-sync-pairing-deviceid-base64url did for the pairing twin. Do NOT encode the strict trailing-bits rule in the pattern: whether it is a deviation worth recording is still open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_realm_management_api.yaml — `data` is declared as a bare `type: object` with no properties (yaml:250-261) although the handler parses it with the same `interfaceschema.ParseInterfaceCanonical` as `installInterface` (internal/realm/service.go:196 vs 161) and then checks the body's `interface_name`/`version_major` against the URL's `{name}`/`{major}` (service.go:203-208), so the update body — the one carrying every interface rule — is the only interface body in the spec with nothing to validate against. `$ref` `#/components/schemas/InterfaceDefinition` as `installInterface` does (yaml:115-117), and say in the operation description that the body must name the same interface and major as the URL (disagreement is 409, already documented) and that `version_minor` must increase (409). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.
- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` shows `detail: Bad Request` (yaml:1732) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1792), while the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58 via `WriteBadRequest`/`WriteInternalServerError`, golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json); every 400/500 on that surface goes through the same astarteapi writers. The other detail examples on the file are already correct — change only those two — and add the docs-half assertion `TestAppEngineErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go, which pins the pairing and realm-management halves. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDetail` schema-level example (yaml:361-362) and the `BadRequest` response (yaml:385) show `detail: Bad Request`, and `InternalServerError` (yaml:477) shows `detail: Internal Server Error`, but the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58, via WriteBadRequest/WriteInternalServerError, pinned by pkg/astarteapi/testdata/error_bad_request.json and error_internal.json). While you are in `ErrorDetail`: its `example` sits at the schema level and shows a bare `{detail: ...}` where the envelope the wire sends is `{"errors":{"detail":...}}` — move it under `errors` (or drop it, as the realm-management and pairing `ErrorDetail` schemas carry no example) so a generated client is not shown the wrong shape. Add `TestHousekeepingErrorDetailExamples` beside the pairing and realm-management halves in internal/swagger/swagger_test.go. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `"Internal server error"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implements `device_empty_cache_received` end to end — `SimpleEvent` field 23, an empty `DeviceEmptyCacheReceivedEvent` whose encoder emits exactly `{"type":"device_empty_cache_received"}` (astarte_core encoder.ex:54-62), fired by DUP at `control_handler.ex:161` *inside* the `with` that runs `send_control_consumer_properties` → `resend_all_properties` → `set_pending_empty_cache`, so it fires only when all three succeed and never on failure, and the commit deletes the `# TODO: implement on_empty_cache_received` from `astarte_events/lib/astarte_events/triggers/core.ex`. Astrate already accepts the trigger but never evaluates it: `OnDeviceEmptyCacheReceived` sits in `deviceOns` with value `false` (internal/engine/triggers/match.go:123, the `evaluated` flag read at match.go:406) and the comment's stated reason ("upstream defines no SimpleEvent variant for it") is now false. Add the event body constructor to internal/engine/triggers/events.go (the payload is `{"type":"device_empty_cache_received"}` and nothing else — no device id, no metadata), flip the `deviceOns` value to true, and fire it from `handleEmptyCache` (internal/engine/control.go:58) with the same idiom `handleLifecycle` uses for `device_connected` (`e.fireDevice` + `stream.Publish`, internal/engine/engine.go:346-353), at the point where `m.Ack()` runs — i.e. after `resendServerProperties` and `sendConsumerProperties` have both succeeded inside `retryStore`. Add a T1 test that a trigger on it fires exactly once per accepted emptyCache and does not fire when the resend fails. — BLOCKED: wrote nothing
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

$ tail -80 .mule/reports/log/20261008T182643Z-recipe-docs-sync-proposed.md
221:- [ ] flow-boot-resumes-stopped-durable-flows [legion] [auto]: RehydrateAutoRestart (cmd/astrate/main.go:216) restarts every flow with auto_restart=true whatever its last status — `ListAutoRestartFlows` filters on `f.auto_restart = true` alone (internal/store/flows.go:125) and nothing ever clears the flag, which is written once at create (flowapi/http.go:141-146 -> service.go:318) and never on stop — so a durable flow stopped through the API is running again after the next boot, while the shutdown mark that says otherwise (MarkRunningFlowsStopped, main.go:497 -> internal/flowapi/service.go:885-901, which writes status="stopped") is never read by that query. The two halves disagree: service.go:325 documents "every durable flow", the shutdown call assumes "only what was running". Say which is the intent, then make them agree — either add `AND f.status = 'running'` to flows.go:125 (making the shutdown mark load-bearing) or delete the mark as dead code — and add a case in internal/store/flows_test.go pinning the chosen rule (a stopped auto_restart row must, or must not, come back). Needs the DB.
246:- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `base64.RawURLEncoding.Strict()` (pkg/deviceid/deviceid.go:37), so of the 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with `ErrInvalid` while `..._lA` is accepted, and the pairing spec's `hw_id` pattern `^[A-Za-z0-9_-]{22}$` (added by docs-sync-pairing-deviceid-base64url) deliberately does not encode the rule because a pattern is a client-side contract and upstream's leniency is unmeasured. The docstring justifies strictness as byte-identical round-tripping, which is a reason, not a parity claim, and docs/COMPATIBILITY.md records no entry for it. Measure upstream (Elixir `Base.url_decode64!(padding: false)`, which does not check those bits) on the Legion Go, then either record the divergence in docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.
259:- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.
292:- [ ] engine-e2e-emptycache-resync-assert [legion] [auto]: make the second wait in `testE2EEmptyCache` (internal/engine/engine_e2e_test.go:367) wait for the resync re-publish instead of re-matching the retained set captured at :362 — mark the cursor before `dev.EmptyCache(t)` (:364) and take only messages after the mark, so the "resent property value = int32(42)" assertion at :368-369 fails when the resend never happens; today it cannot fail, because `WaitForTopic` returns the pre-EmptyCache message, which carries the same value 42. Needs `-tags e2e` plus Docker, so this half is Legion-only — the container-free helper and its test land in testutil-wait-since-cursor.

→ Read internal/swagger/swagger_test.go [limit=90, offset=1086]

$ cat docs/Makefile; echo ===; ls docs/
.PHONY: serve build clean install sync

install:
	pip install -r requirements.txt

sync:
	cp DESIGN.md site/DESIGN.md
	cp ROADMAP.md site/ROADMAP.md
	cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
	rm -rf site/api && cp -r api site/api
	rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui

serve: sync
	mkdocs serve -f mkdocs.yml

build: sync
	mkdocs build -f mkdocs.yml

clean:
	rm -rf ../site-dist
===
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

$ tail -c 600 .mule/todo.md | cat -A | tail -5 | cut -c1-200; echo; wc -l .mule/todo.md; tail -1 .mule/todo.md | cut -c1-120
ection` (internal/engine/introspection.go:99-127) is tested only against hand-written literals in internal/engine/introspection_test.go:17-71, so nothing container-free checks that the string the harn

294 .mule/todo.md
- [x] engine-introspection-producer-roundtrip [auto]: pin the producer half of the docs/DESIGN.md §3.3 introspection ru

$ tail -c 300 .mule/todo.md
olds. Add a case to internal/engine/introspection_test.go that builds the string from an unsorted multi-entry map via `testutil.Introspection`, asserts `parseIntrospection` returns the same map, and pins the exact output (sorted, `name:major:minor;…`, no trailing `;`, empty map → empty string).

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -291,4 +291,8 @@
 - [x] testutil-wait-since-cursor [auto]: give AstarteDevice's message waits a cursor so a second wait for the same topic cannot re-match an older message — `WaitForMessage`/`WaitForTopic` (internal/testutil/astartedevice.go:149-170) scan the capture buffer from index 0 on every poll, so a repeated wait returns the first match again, and in-tree that makes testE2EEmptyCache's resync wait vacuous (see the companion line). Move the paho capture handler out of the `ConnectAstarteDevice` closure (astartedevice.go:61-67) into an injectable method, add `Mark() int` and a `From(start int, …)` wait variant, and add a container-free internal/testutil/astartedevice_test.go that drives two captures on one topic and asserts a wait started after the second returns the second message — it cannot even compile today, which is the failure proof. Keep every new read of `received` under the existing `d.mu`, and say in the report that the change is unverified for races (the Pi gate runs no `-race`). The consumer of this API is engine-e2e-emptycache-resync-assert.
 - [ ] engine-e2e-emptycache-resync-assert [legion] [auto]: make the second wait in `testE2EEmptyCache` (internal/engine/engine_e2e_test.go:367) wait for the resync re-publish instead of re-matching the retained set captured at :362 — mark the cursor before `dev.EmptyCache(t)` (:364) and take only messages after the mark, so the "resent property value = int32(42)" assertion at :368-369 fails when the resend never happens; today it cannot fail, because `WaitForTopic` returns the pre-EmptyCache message, which carries the same value 42. Needs `-tags e2e` plus Docker, so this half is Legion-only — the container-free helper and its test land in testutil-wait-since-cursor.
 - [!] testutil-control-frame-test [auto]: pin the producer/properties control frame that `SendProducerProperties` puts on the wire (internal/testutil/astartedevice.go:181-216) — `DeflateControlList`/`InflateControlList` have zero tests and every caller of the helpers sits behind an `e2e`/`integration`/conformance build tag, so nothing on the default gate asserts the 4-byte big-endian uncompressed-length prefix plus zlib framing, while the parser side is pinned in internal/engine/control_test.go:89-130 against a CPython golden frame. Add a container-free internal/testutil/astartedevice_test.go with exact-bytes, round-trip, empty-list→nil, sub-4-byte header and declared≠inflated cases; bound the inflate with `io.LimitReader(zr, int64(declared)+1)` the way internal/engine/control.go:301-315 does (astartedevice.go:205 reads with unbounded `io.ReadAll` today, and a lying header from the code under test makes the harness allocate it before failing) and assert the truncated byte count in the mismatch message so the bound itself is pinned. One extra line in internal/engine/control_test.go feeding `testutil.DeflateControlList` into `inflateProperties` is the cross-check that closes the loop. — BLOCKED: wrote nothing
 - [x] engine-introspection-producer-roundtrip [auto]: pin the producer half of the docs/DESIGN.md §3.3 introspection rule on the default gate — `testutil.Introspection` (internal/testutil/astartedevice.go:93-105) is called from six places, all behind `e2e`, `integration && e2e` or test/conformance build tags, while the consumer `parseIntrospection` (internal/engine/introspection.go:99-127) is tested only against hand-written literals in internal/engine/introspection_test.go:17-71, so nothing container-free checks that the string the harness renders is the one the engine accepts or that the documented "deterministic (sorted) ordering" holds. Add a case to internal/engine/introspection_test.go that builds the string from an unsorted multi-entry map via `testutil.Introspection`, asserts `parseIntrospection` returns the same map, and pins the exact output (sorted, `name:major:minor;…`, no trailing `;`, empty map → empty string).
+- [ ] docs-sync-pairing-register-404-unreachable [auto]: drop the `404` `DeviceNotFound` response from `registerDevice` in docs/api/astarte_pairing_api.yaml (yaml:140-141) — it is documented but absent: the route is wrapped by `requireAgent` (internal/pairing/http.go:82-83), i.e. `mw.RequireRealm(auth.ClaimPairing)`, which resolves the realm *before* the handler and answers `401 {"errors":{"detail":"Unauthorized"}}` for an unknown one (internal/auth/middleware.go:66-72, "no existence oracle on auth failures"), so `Service.Register`'s own `GetRealmByName` (internal/pairing/service.go:183-186) can return `store.ErrNotFound` only if the realm is deleted between the middleware's lookup and the handler's — a race, not a contract. Nothing else in Register returns `store.ErrNotFound`: `RegisterDevice`'s `ErrDeviceAlreadyConfirmed` becomes 422 (service.go:209-213) and the CountDevices/SetPayloadFormatHint failures are DB errors → 500. No test covers it (internal/pairing/http_test.go asserts 404 only for the agent DELETE at :527 and for health at :595-598), and the `DeviceNotFound` component's own description — "unknown device ID" — does not match what that path would be even if it fired. The premise came from docs-sync-pairing-register-404, which read service.go without the middleware in front of it. Keep the `404` on `unregisterDevice` (reachable: unknown device, and a malformed `deviceID` via service.go:231-234) and on `getHealth`. State the reachability argument in the report; if you find a reachable path instead, keep the response and fix its description rather than deleting it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
+- [ ] docs-sync-pairing-body-cap-400 [auto]: document the 64 KiB request-body cap on the three body-taking pairing ops — `registerDevice`, `requestCredentials` and `verifyCredentials` all decode through `astarteapi.DecodeData(r.Body, maxBodyBytes, …)` with `maxBodyBytes = 64 << 10` (internal/pairing/http.go:17, 116, 169, 249), and any body over 65536 bytes fails with an error wrapping `ErrBodyTooLarge` (pkg/astarteapi/envelope.go:233-235) that each handler turns into a plain 400 `{"errors":{"detail":"Bad request"}}` (http.go:117-118, 170-171, 250-251). The shared `BadRequest` component (yaml:578-586) says only "malformed or missing data envelope" and no request-body schema carries a size hint, so a client is told nothing about the cap. Widen that one shared description to name the 64 KiB cap and the failure shape (a single edit covers all three operations) and say explicitly that an over-size body is 400, not 413 — do not invent a 413. One clause on `csr` (yaml:414) and `client_crt` (yaml:422) too: the cap is on the whole envelope, not on the PEM. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
+- [ ] docs-sync-pairing-realm-name-pattern [auto]: constrain `components.parameters.RealmName` in docs/api/astarte_pairing_api.yaml (yaml:359-365, a bare `type: string` described only as "The realm name.") with `pattern: '^[a-z][a-z0-9]*$'` — the realm name is CHECK-constrained to exactly that (migrations/000002_metadata.up.sql:7, read-only, never edit) and the two surfaces whose tasks already ran carry it on their twins (docs/api/astarte_housekeeping_api.yaml:246, docs/api/astrate_native_api.yaml:420 and 467). Because this one parameter is shared by all seven operations, the description must record that an off-pattern name gets three different answers, not one: the agent and device routes answer **401** (the realm is resolved by the middleware, which treats a miss as unauthenticated rather than 404 — internal/auth/middleware.go:66-72); `getHealth` answers **404** (service.Health returns store.ErrNotFound → WriteNotFound, internal/pairing/service.go:369-371, http.go:289-290); and `getVersion` answers **200** `{"data":"1.2.2"}` for any realm string whatsoever, because `GET /pairing/v1/{realm}/version` is registered with no middleware and no realm lookup (cmd/astrate/main.go:438 → observability.VersionHandler, internal/observability/compat.go:38-47) — say that in the version operation so nobody reads the shared parameter as a guarantee the handler enforces. Mirror `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go:1247) with a pairing twin only if it stays a few lines; the YAML edit is the deliverable. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
+- [ ] docs-sync-pairing-deviceid-path-param [auto]: give `components.parameters.DeviceID` the constraints `RegisterRequest.hw_id` already carries in docs/api/astarte_pairing_api.yaml — the parameter (yaml:367-377) is prose-only, while `hw_id` (yaml:386-398) got `pattern: '^[A-Za-z0-9_-]{22}$'` plus `minLength`/`maxLength` from docs-sync-pairing-deviceid-base64url, and `TestPairingDeviceIDEncodingDocumented` (internal/swagger/swagger_test.go:1086) pins only the two descriptions and hw_id's bounds, never the parameter's schema; the realm-management twin was constrained later by docs-sync-rm-deviceid-param and its test (swagger_test.go:1168) does hold the parameter's pattern against `deviceid.Parse`. Add the same three schema lines to the parameter, and record in its description the asymmetry the prose omits: a malformed `deviceID` is **401** on all three device-authenticated operations — `deviceid.Parse` failing inside `authenticateDevice` burns bcrypt and returns `ErrUnauthorized` (internal/pairing/service.go:381-385) → `WriteUnauthorized` (http.go:322-323), exactly as an unknown device or a wrong secret is (service.go:395-397), so the client is told its credentials are bad rather than its ID malformed — but **404** on the agent `unregisterDevice` (service.go:231-234 → http.go:326-327). Extend TestPairingDeviceIDEncodingDocumented to compile the new pattern and hold it against `deviceid.Parse` the way the realm-management test does. Do not encode the strict trailing-bits rule in the pattern; that stays open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).

$ wc -l .mule/todo.md; grep -c '^- \[ \] docs-sync-pairing-\(register-404-unreachable\|body-cap-400\|realm-name-pattern\|deviceid-path-param\)' .mule/todo.md; awk 'NR>294 && !/^- \[ \]/ {print NR": NOT-A-LINE: "substr($0,1,80)}' .mule/todo.md; git status --porcelain
298 .mule/todo.md
4
 M .mule/todo.md

Read `.mule/MULE.md` and `.mule/recipes/docs-sync.md`, then ran the recipe for the **pairing** surface (the one least recently swept — last run 2026-09-28): compared the seven routes in `internal/pairing/http.go:81-91` against `docs/api/astarte_pairing_api.yaml`, and checked health, register, unregister, info and verify status/field mappings plus `docs/site/pairing-and-security.md`. No `gh issue create` in this recipe; no `docs/site` contradiction found, so nothing appended to `for-giulio.md` (the config-key sweep is already recorded there from an earlier run).

Done: appended 4 proposal task lines to .mule/todo.md (pairing surface drift; nothing implemented, no git touched)
Files: .mule/todo.md
Verified: `wc -l` 294→298, 4 new `- [ ]` lines, each single-line, no duplicate slugs, no other file modified (`git status --porcelain` = only `.mule/todo.md`)
Unsure: whether `registerDevice`'s documented `404` is truly unreachable — argued from read code only (middleware resolves the realm before the handler); the task tells the executor to prove it or keep and reword instead
Follow-ups: docs-sync-pairing-register-404-unreachable, docs-sync-pairing-body-cap-400, docs-sync-pairing-realm-name-pattern, docs-sync-pairing-deviceid-path-param

$ for f in $(ls -t .mule/reports/log/*recipe-docs-sync-proposed.md | head -12); do echo "=== $f"; grep -io 'surface[^.]*' "$f" | head -3; done
=== .mule/reports/log/20261008T182643Z-recipe-docs-sync-proposed.md
surface"
surfaces
surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service
=== .mule/reports/log/20261007T191848Z-recipe-docs-sync-proposed.md
surface
surface"
surfaces
=== .mule/reports/log/20261006T184046Z-recipe-docs-sync-proposed.md
surface issue #78 is scoped against, and none of it is in #78's body (written against upstream 1
surface Astrate does not emulate (#78, parked)
surface"
=== .mule/reports/log/20261005T111650Z-recipe-docs-sync-proposed.md
surface: appengine)
surface once per delivery" — but `cfg
surface
=== .mule/reports/log/20261004T112505Z-recipe-docs-sync-proposed.md
surface
surface issue #78 is scoped against, and none of it is in #78's body (written against upstream 1
surface Astrate does not emulate (#78, parked)
=== .mule/reports/log/20261002T181712Z-recipe-docs-sync-proposed.md
surfaces, never recovery
surface once per delivery" — but `cfg
surface
=== .mule/reports/log/20261001T185956Z-recipe-docs-sync-proposed.md
=== .mule/reports/log/20261001T164758Z-recipe-docs-sync-proposed.md
=== .mule/reports/log/20260930T192852Z-recipe-docs-sync-proposed.md
=== .mule/reports/log/20260930T175957Z-recipe-docs-sync-proposed.md
=== .mule/reports/log/20260930T113453Z-recipe-docs-sync-proposed.md
=== .mule/reports/log/20260929T183752Z-recipe-docs-sync-proposed.md

$ for f in $(ls -t .mule/reports/log/*recipe-docs-sync-proposed.md | head -10); do echo "=== $f"; grep -ioE '(surface: [a-z-]+|surface [a-z-]+ surface|least recently swept[^.]*|last run [0-9-]+)' "$f" | head -4; done
=== .mule/reports/log/20261008T182643Z-recipe-docs-sync-proposed.md
least recently swept — last run 2026-09-28): compared the seven routes in `internal/pairing/http
=== .mule/reports/log/20261007T191848Z-recipe-docs-sync-proposed.md
surface: pairing
=== .mule/reports/log/20261006T184046Z-recipe-docs-sync-proposed.md
=== .mule/reports/log/20261005T111650Z-recipe-docs-sync-proposed.md
surface: appengine
=== .mule/reports/log/20261004T112505Z-recipe-docs-sync-proposed.md
=== .mule/reports/log/20261002T181712Z-recipe-docs-sync-proposed.md
=== .mule/reports/log/20261001T185956Z-recipe-docs-sync-proposed.md
=== .mule/reports/log/20261001T164758Z-recipe-docs-sync-proposed.md
=== .mule/reports/log/20260930T192852Z-recipe-docs-sync-proposed.md
=== .mule/reports/log/20260930T175957Z-recipe-docs-sync-proposed.md

$ for f in 20261007T191848Z 20261006T184046Z 20261005T111650Z 20261004T112505Z 20261003T165209Z 20261002T181712Z; do echo "=== $f"; tail -12 .mule/reports/log/${f}-recipe-docs-sync-proposed.md 2>/dev/null | grep -iE 'surface|Follow-ups|Done:'; done
=== 20261007T191848Z
Done: docs-sync pairing run — spec accurate (no task lines); escalated the site rate-limiting contradiction to for-giulio.md
Follow-ups: none
=== 20261006T184046Z
Done: docs-sync (astrate-native surface) — proposed 4 docs/api fix tasks + 2 for-giulio.md escalations, nothing implemented
Follow-ups: (todo.md) docs-native-socket-handshake-errors, docs-native-socket-sse-exact-accept, docs-native-socket-security-scheme, docs-native-socket-event-schema; (for-giulio.md) appengine-api.md envelope-format and a_ch-room-filters contradictions
=== 20261005T111650Z
Done: docs-sync recipe run on surface **astrate-native** — five `docs/api/` fix tasks appended to `.mule/todo.md` (no code touched, no git, no `gh`; the recipe's deliverable here is task lines, not an issue).
Follow-ups: five lines appended — `docs-native-ae-realm-version-endpoint`, `docs-native-hk-version-description`, `docs-native-metrics-content-negotiation`, `docs-native-phoenix-newevent-payload`, `docs-native-realm-name-pattern`. The gate (`make -C docs build`) was **not** run: nothing under `docs/` changed in this run.
=== 20261004T112505Z
Ran the docs-sync recipe on the **housekeeping** surface (one per run). Route sets match (5 routes, 2 paths — `docs-sync-hk-patch-endpoint` had already closed the only gap), so the drift was in status codes, examples and constraints.
=== 20261003T165209Z
=== 20261002T181712Z
Done: docs-sync recipe on the realm-management surface — compared all 20 mounted routes against the 12 documented paths (every route documented, no dead paths), then hand-checked status codes and field schemas against the handlers; five `docs/api/` corrections appended to `.mule/todo.md`, nothing implemented, no GitHub or git action taken.
Follow-ups: five lines appended — `docs-sync-rm-error-example-capitalisation`, `docs-sync-rm-legacy-alias-fields`, `docs-sync-rm-validationerror-example`, `docs-sync-rm-deviceid-param`, `docs-sync-rm-update-interface-body`

# Todos
[•] Pick surface (appengine) and compare routes vs documented paths
[ ] Check 3-4 endpoints' status codes and response fields vs handlers
[ ] Check config keys / docs-site contradictions
[ ] Append up to 5 task lines to .mule/todo.md
[ ] Write report to .mule/reports if recipe requires; final report

→ Read internal/appengine/http.go [limit=260]

$ rg -n '^\s{2}/' docs/api/astarte_appengine_api.yaml
26:  /appengine/v1/{realm}/version:
67:  /appengine/v1/{realm}/devices:
141:  /appengine/v1/{realm}/stats/devices:
181:  /appengine/v1/{realm}/devices/{device}:
259:  /appengine/v1/{realm}/devices-by-alias/{alias}:
347:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces:
381:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
428:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
625:  /appengine/v1/{realm}/devices/{device}/interfaces:
655:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
697:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
874:  /appengine/v1/{realm}/groups:
952:  /appengine/v1/{realm}/groups/{group}:
982:  /appengine/v1/{realm}/groups/{group}/devices:
1078:  /appengine/v1/{realm}/groups/{group}/devices/{device}:
1183:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces:
1213:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}:
1256:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}:

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*appengine*'
(no output)

→ Read internal/appengine/http.go [offset=260, limit=300]
→ Read internal/appengine/http.go [offset=560, limit=188]

$ rg -n 'operationId:|^\s{4}(get|post|put|delete|patch):|^\s{6}\x27?[0-9]{3}\x27?:' docs/api/astarte_appengine_api.yaml | head -200
27:    get:
28:      operationId: getVersion
68:    get:
69:      operationId: listDevices
142:    get:
143:      operationId: getDeviceStats
182:    get:
183:      operationId: getDevice
208:    patch:
209:      operationId: patchDevice
260:    get:
261:      operationId: getDeviceByAlias
291:    patch:
292:      operationId: patchDeviceByAlias
348:    get:
349:      operationId: listInterfacesByAlias
382:    get:
383:      operationId: getDataByAlias
429:    get:
430:      operationId: getDataAtPathByAlias
480:    put:
481:      operationId: putDataByAlias
532:    post:
533:      operationId: publishDataByAlias
584:    delete:
585:      operationId: deleteDataByAlias
626:    get:
627:      operationId: listDeviceInterfaces
656:    get:
657:      operationId: getData
698:    get:
699:      operationId: getDataAtPath
742:    put:
743:      operationId: putData
789:    post:
790:      operationId: publishData
836:    delete:
837:      operationId: deleteData
875:    get:
876:      operationId: listGroups
902:    post:
903:      operationId: createGroup
953:    get:
954:      operationId: getGroup
983:    get:
984:      operationId: listGroupDevices
1032:    post:
1033:      operationId: addGroupDevice
1079:    get:
1080:      operationId: getGroupDevice
1108:    patch:
1109:      operationId: patchGroupDevice
1160:    delete:
1161:      operationId: removeGroupDevice
1184:    get:
1185:      operationId: listInterfacesInGroup
1214:    get:
1215:      operationId: getDataInGroup
1257:    get:
1258:      operationId: getDataAtPathInGroup
1302:    put:
1303:      operationId: putDataInGroup
1350:    post:
1351:      operationId: publishDataInGroup
1398:    delete:
1399:      operationId: deleteDataInGroup

→ Read docs/api/astarte_appengine_api.yaml [offset=62, limit=200]
→ Read docs/api/astarte_appengine_api.yaml [offset=870, limit=340]
→ Read docs/api/astarte_appengine_api.yaml [offset=625, limit=245]
→ Read docs/api/astarte_appengine_api.yaml [offset=1400, limit=494]
→ Read docs/api/astarte_appengine_api.yaml [offset=259, limit=370]
→ Read docs/api/astarte_appengine_api.yaml [offset=1210, limit=195]

$ rg -n 'ErrGroupNotFound|ErrGroupAlreadyExists|ErrDeviceAlreadyInGroup|ErrValidation|ErrInvalidAlias|ErrAliasTagNotFound|ErrInvalidAttributes|ErrAttributeKeyNotFound|return.*ErrNotFound|func \(s \*Service\)' internal/appengine/service.go | head -120
27:// ErrValidation wraps a well-formed request that violates a rule (maps to 422).
28:var ErrValidation = errors.New("appengine: validation failed")
34:	ErrInvalidAlias         = errors.New("Invalid alias")           //nolint:staticcheck // ST1005: detail strings are upstream astarte_appengine_api 1.2.2 wire text
36:	ErrAliasTagNotFound     = errors.New("Alias tag not found")     //nolint:staticcheck // ST1005: upstream wire text
37:	ErrInvalidAttributes    = errors.New("Invalid attributes")      //nolint:staticcheck // ST1005: upstream wire text
38:	ErrAttributeKeyNotFound = errors.New("Attribute key not found") //nolint:staticcheck // ST1005: upstream wire text
41:// ErrGroupNotFound marks a missing group (maps to 404 "Group not found").
42:var ErrGroupNotFound = errors.New("Group not found") //nolint:staticcheck // ST1005: upstream wire text
44:// ErrGroupAlreadyExists marks a duplicate group name (maps to 409 "Group
47:var ErrGroupAlreadyExists = fmt.Errorf("Group already exists: %w", store.ErrAlreadyExists) //nolint:staticcheck // ST1005: upstream wire text
49:// ErrDeviceAlreadyInGroup marks re-adding an existing member (maps to 409
51:var ErrDeviceAlreadyInGroup = fmt.Errorf("Device already in group: %w", store.ErrAlreadyExists) //nolint:staticcheck // ST1005: upstream wire text
65:// ErrGroupNotFound (→ 404 "Group not found") and store.ErrNotFound.
67:	return fmt.Errorf("%w: %w: group %q", ErrGroupNotFound, store.ErrNotFound, name)
93:func (s *Service) realmID(ctx context.Context, realm string) (int16, error) {
141:func (s *Service) ListDevices(ctx context.Context, realm string, after string, limit int, details bool) (*DevicePage, error) {
153:			return nil, fmt.Errorf("%w: invalid cursor", ErrValidation)
182:func (s *Service) DevicesStats(ctx context.Context, realm string) (total, connected int64, err error) {
193:func (s *Service) deviceStatusBatch(ctx context.Context, rid int16, devs []store.Device) ([]*DeviceStatus, error) {
210:func (s *Service) GetDevice(ctx context.Context, realm, deviceID string) (*DeviceStatus, error) {
217:		return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
228:func (s *Service) GetDeviceByAlias(ctx context.Context, realm, alias string) (*DeviceStatus, error) {
242:func (s *Service) deviceStatus(ctx context.Context, rid int16, d *store.Device) (*DeviceStatus, error) {
291:func (s *Service) PatchDevice(ctx context.Context, realm, deviceID string, p DevicePatch) (*DeviceStatus, error) {
298:		return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
309:func (s *Service) PatchDeviceByAlias(ctx context.Context, realm, alias string, p DevicePatch) (*DeviceStatus, error) {
325:func (s *Service) applyPatch(ctx context.Context, rid int16, id deviceid.ID, d *store.Device, p DevicePatch) (*DeviceStatus, error) {
328:			return nil, ErrInvalidAlias
342:			return nil, ErrInvalidAttributes
348:				return nil, ErrAliasTagNotFound
355:				return nil, ErrAttributeKeyNotFound
382:// PATCH /groups/{group}/devices/{device}). Unknown group → ErrGroupNotFound;
384:func (s *Service) PatchGroupDevice(ctx context.Context, realm, groupName, deviceID string, p DevicePatch) (*DeviceStatus, error) {
394:func (s *Service) GetGroupDevice(ctx context.Context, realm, groupName, deviceID string) (*DeviceStatus, error) {
403:// ordering: unknown group first (ErrGroupNotFound), then membership — a
406:func (s *Service) groupMember(ctx context.Context, realm, groupName, deviceID string) (int16, deviceid.ID, *store.Device, error) {
416:		return 0, deviceid.ID{}, nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
423:		return 0, deviceid.ID{}, nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
436:func (s *Service) aliasValuesTaken(ctx context.Context, rid int16, id deviceid.ID,
462:func (s *Service) CreateGroup(ctx context.Context, realm, name string, devices []string) error {
495:			return fmt.Errorf("%w: %w", ErrGroupAlreadyExists, err)
508:func (s *Service) GetGroup(ctx context.Context, realm, name string) error {
520:func (s *Service) ListGroups(ctx context.Context, realm string) ([]string, error) {
541:func (s *Service) ListGroupDevices(ctx context.Context, realm, name string, details bool, fromToken string, limit int) (*DevicePage, error) {
626:func (s *Service) AddGroupDevice(ctx context.Context, realm, name, deviceID string) error {
637:			return fmt.Errorf("%w: %w", ErrDeviceAlreadyInGroup, err)
645:func (s *Service) RemoveGroupDevice(ctx context.Context, realm, name, deviceID string) error {
652:		return fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
658:func (s *Service) group(ctx context.Context, realm, name string) (int16, *store.Group, error) {
674:func (s *Service) resolveOnDevice(ctx context.Context, rid int16, d *store.Device, ifaceName string) (*resolved, error) {
677:		return nil, fmt.Errorf("%w: interface %s not in device introspection", store.ErrNotFound, ifaceName)
687:func (s *Service) resolveByAlias(ctx context.Context, realm, alias, ifaceName string) (*resolved, error) {
703:func (s *Service) resolveInGroup(ctx context.Context, realm, groupName, deviceID, ifaceName string) (*resolved, error) {
713:func (s *Service) readResolved(ctx context.Context, r *resolved, path string, opts QueryOpts) (any, error) {
717:			return nil, fmt.Errorf("%w: downsample_to is not supported on properties interfaces", ErrValidation)
726:func (s *Service) GetDataByAlias(ctx context.Context, realm, alias, ifaceName, path string, opts QueryOpts) (any, error) {
736:func (s *Service) GetDataInGroup(ctx context.Context, realm, groupName, deviceID, ifaceName, path string, opts QueryOpts) (any, error) {
746:func (s *Service) PublishDataByAlias(ctx context.Context, realm, alias, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
755:func (s *Service) PublishDataInGroup(ctx context.Context, realm, groupName, deviceID, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
764:func (s *Service) UnsetPropertyByAlias(ctx context.Context, realm, alias, ifaceName, path string) error {
773:func (s *Service) UnsetPropertyInGroup(ctx context.Context, realm, groupName, deviceID, ifaceName, path string) error {
793:func (s *Service) ListInterfaces(ctx context.Context, realm, deviceID string) ([]string, error) {
800:		return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
810:func (s *Service) ListInterfacesByAlias(ctx context.Context, realm, alias string) ([]string, error) {
823:func (s *Service) ListInterfacesInGroup(ctx context.Context, realm, groupName, deviceID string) ([]string, error) {

→ Read internal/appengine/service.go [offset=141, limit=60]
→ Read internal/appengine/service.go [offset=541, limit=135]

$ grep -n 'docs-sync-ae-list-devices-422' .mule/todo.md; echo "---report---"; cat .mule/reports/log/20260927T182324Z-docs-sync-ae-list-devices-422-blocked.md 2>/dev/null | head -60
240:- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 200/401/500) — a `from_token` that does not parse as a device id is refused with `ErrValidation: invalid cursor` (internal/appengine/service.go:150-155) which `writeError` maps to 422 with the flat `{"errors":{"detail":"invalid cursor"}}` body, i.e. the `ErrorDetail` shape and not `FieldErrorsDetail`, so `$ref` the same response the group twin uses for the identical parameter (`ValidationErrors`, yaml:934-935) only after fixing its example, or add a dedicated one. Extend the `from_token` parameter description (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
---report---
slug: docs-sync-ae-list-devices-422
verdict: blocked
at:  4055896
ran: 2026-09-27T18:23:24Z on DietPi in 1200s

> build · big-pickle

→ Read internal/appengine/service.go [offset=674, limit=173]

$ rg -n 'engine\.(ErrNotServerOwned|ErrNotAProperty|ErrInterfaceNotFound|ErrPathNotFound)|func \(s \*Service\) (PublishData|UnsetProperty|GetData)\b|func \(e \*Engine\)' internal/appengine/service.go internal/engine/*.go | head -60
internal/engine/router.go:198:func (e *Engine) start(ctx context.Context) error {
internal/engine/router.go:218:func (e *Engine) Submit(m broker.InboundMessage) {
internal/engine/router.go:259:func (e *Engine) OnLifecycleEvent(ev broker.LifecycleEvent) {
internal/engine/router.go:271:func (e *Engine) RefreshInterfaces(ctx context.Context, realmID int16) error {
internal/engine/router.go:277:func (e *Engine) runInvalidation(ctx context.Context, ch <-chan store.Notification) {
internal/engine/router.go:297:func (e *Engine) runShard(ctx context.Context, sh *shard) {
internal/engine/router.go:306:func (e *Engine) shardPass(ctx context.Context, sh *shard) (restart bool) {
internal/engine/router.go:347:func (e *Engine) flushShard(ctx context.Context, sh *shard) {
internal/engine/router.go:361:func (e *Engine) drain(ctx context.Context) error {
internal/engine/control.go:40:func (e *Engine) handleControl(ctx context.Context, m broker.InboundMessage, realm *realmSchema, subpath string) {
internal/engine/control.go:58:func (e *Engine) handleEmptyCache(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
internal/engine/control.go:83:func (e *Engine) handleProducerProperties(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
internal/engine/control.go:113:func (e *Engine) resolvePropertyRefs(realm *realmSchema, dev *deviceState, entries []string) []store.PropertyRef {
internal/engine/control.go:144:func (e *Engine) resendServerProperties(ctx context.Context, realm *realmSchema, id deviceid.ID, format payload.Format) error {
internal/engine/control.go:191:func (e *Engine) sendConsumerProperties(ctx context.Context, realm *realmSchema, id deviceid.ID) error {
internal/engine/control.go:230:func (e *Engine) rehydrateValue(stored []byte, m *interfaceschema.CompiledMapping) (payload.Value, error) {
internal/engine/introspection.go:35:func (e *Engine) handleIntrospection(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
internal/engine/introspection.go:135:func (e *Engine) retryStore(ctx context.Context, m broker.InboundMessage, what string, fn func() error) bool {
internal/engine/engine.go:82:func (e *Engine) AttachBroker(bp BrokerPort) {
internal/engine/engine.go:89:func (e *Engine) Start(ctx context.Context) error {
internal/engine/engine.go:101:func (e *Engine) Drain(ctx context.Context) error {
internal/engine/engine.go:126:func (e *Engine) Bus() *stream.Bus {
internal/engine/engine.go:133:func (e *Engine) RefreshTriggers(ctx context.Context, realmID int16) error {
internal/engine/engine.go:140:func (e *Engine) fireCommitted(ops []PersistOp) {
internal/engine/engine.go:153:func (e *Engine) fireData(rs *realmSchema, op *PersistOp) {
internal/engine/engine.go:227:func (e *Engine) fireDataChanges(rs *realmSchema, op *PersistOp, value any) {
internal/engine/engine.go:285:func (e *Engine) fireDevice(rs *realmSchema, id deviceid.ID, at time.Time, match triggers.DeviceEvent, body any) {
internal/engine/engine.go:305:func (e *Engine) fireDeviceError(m broker.InboundMessage, reason, detail string) {
internal/engine/engine.go:334:func (e *Engine) handleLifecycle(ev broker.LifecycleEvent) {
internal/engine/engine.go:377:func (e *Engine) HandleDeviceRegistered(realmName string, hwID string, at time.Time) {
internal/engine/engine.go:398:func (e *Engine) HandleDeviceDeletionStarted(realmName string, hwID string, at time.Time) {
internal/engine/engine.go:407:func (e *Engine) HandleDeviceDeletionFinished(realmName string, hwID string, at time.Time) {
internal/engine/engine.go:413:func (e *Engine) handleDeviceDeletion(realmName, hwID string, at time.Time, on, kind string, body any) {
internal/engine/engine.go:444:func (e *Engine) eventValue(v any) any {
internal/engine/capabilities.go:29:func (e *Engine) handleCapabilities(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
internal/engine/data.go:153:func (e *Engine) handle(ctx context.Context, sh *shard, m broker.InboundMessage) {
internal/engine/data.go:218:func (e *Engine) deviceState(ctx context.Context, m broker.InboundMessage, realm *realmSchema) (*deviceState, bool) {
internal/engine/data.go:245:func (e *Engine) processData(ctx context.Context, sh *shard, m broker.InboundMessage, realm *realmSchema, dev *deviceState, rest string) {
internal/engine/data.go:361:func (e *Engine) trackPrevious(ctx context.Context, sh *shard, rs *realmSchema, op *PersistOp) {
internal/engine/data.go:419:func (e *Engine) lookupPrev(ctx context.Context, op *PersistOp) (prev []byte, found bool, err error) {
internal/engine/data.go:454:func (e *Engine) capturePrevious(ctx context.Context, rs *realmSchema, op *PersistOp) {
internal/engine/data.go:485:func (e *Engine) lookupFailed(op *PersistOp, err error) {
internal/engine/data.go:509:func (e *Engine) noteFormat(ctx context.Context, id deviceid.ID, realm *realmSchema, dev *deviceState, f payload.Format) {
internal/engine/data.go:543:func (e *Engine) reject(m broker.InboundMessage, reason, detail string) {
internal/engine/serverdata.go:50:func (e *Engine) PublishServerValue(ctx context.Context, realm string, id deviceid.ID, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
internal/engine/serverdata.go:72:func (e *Engine) PublishDeviceValue(ctx context.Context, realm string, id deviceid.ID, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
internal/engine/serverdata.go:98:func (e *Engine) publishAsOwner(ctx context.Context, realm string, id deviceid.ID,
internal/engine/serverdata.go:222:func (e *Engine) UnsetServerProperty(ctx context.Context, realm string, id deviceid.ID, ifaceName, path string) error {
internal/engine/serverdata.go:292:func (e *Engine) deviceView(ctx context.Context, rs *realmSchema, id deviceid.ID) (deviceView, error) {

$ rg -n 'ErrNotServerOwned|ErrNotAProperty|ErrInterfaceNotFound|ErrPathNotFound' internal/engine/ | head -40
internal/engine/serverdata_test.go:174:			iface: "com.example.Nope", path: "/x", value: "42", wantErr: ErrInterfaceNotFound},
internal/engine/serverdata_test.go:176:			iface: "com.astrate.test.Minimal", path: "/value", value: "1.5", wantErr: ErrNotServerOwned},
internal/engine/serverdata_test.go:178:			iface: propsIface, path: "/limits/nope", value: "42", wantErr: ErrPathNotFound},
internal/engine/serverdata_test.go:263:	if err := rig.e.UnsetServerProperty(ctx, realmAlpha, devAlpha, "com.astrate.test.ServerStream", "/cmd"); !errors.Is(err, ErrNotAProperty) {
internal/engine/serverdata_test.go:264:		t.Errorf("unset on datastream: %v, want ErrNotAProperty", err)
internal/engine/serverdata_test.go:310:	if !errors.Is(err, ErrPathNotFound) {
internal/engine/serverdata_test.go:311:		t.Errorf("bad prefix err = %v, want ErrPathNotFound", err)
internal/engine/serverdata.go:22:	// ErrInterfaceNotFound: the realm has no installed interface by that name.
internal/engine/serverdata.go:23:	ErrInterfaceNotFound = errors.New("engine: interface not installed")
internal/engine/serverdata.go:24:	// ErrNotServerOwned: the interface is device-owned (docs/DESIGN.md §2.6
internal/engine/serverdata.go:26:	ErrNotServerOwned = errors.New("engine: interface is not server-owned")
internal/engine/serverdata.go:30:	// ErrPathNotFound: the path resolves no endpoint mapping.
internal/engine/serverdata.go:31:	ErrPathNotFound = errors.New("engine: path matches no endpoint")
internal/engine/serverdata.go:32:	// ErrNotAProperty: a property operation on a datastream interface.
internal/engine/serverdata.go:33:	ErrNotAProperty = errors.New("engine: interface is not a properties interface")
internal/engine/serverdata.go:103:		notOwned = ErrNotServerOwned
internal/engine/serverdata.go:115:		return nil, fmt.Errorf("%w: %s", ErrInterfaceNotFound, ifaceName)
internal/engine/serverdata.go:137:			return nil, fmt.Errorf("%w: %q is not an aggregation prefix of %s", ErrPathNotFound, path, ifaceName)
internal/engine/serverdata.go:145:			return nil, fmt.Errorf("%w: %q matches no endpoint of %s", ErrPathNotFound, path, ifaceName)
internal/engine/serverdata.go:233:		return fmt.Errorf("%w: %s", ErrInterfaceNotFound, ifaceName)
internal/engine/serverdata.go:236:		return fmt.Errorf("%w: %s", ErrNotServerOwned, ifaceName)
internal/engine/serverdata.go:239:		return fmt.Errorf("%w: %s", ErrNotAProperty, ifaceName)
internal/engine/serverdata.go:243:		return fmt.Errorf("%w: %q matches no endpoint of %s", ErrPathNotFound, path, ifaceName)
internal/engine/devicedata_test.go:143:			iface: "com.example.Nope", path: "/value", value: `"x"`, wantErr: ErrInterfaceNotFound},
internal/engine/devicedata_test.go:147:			iface: stream, path: "/nope", value: `"x"`, wantErr: ErrPathNotFound},

$ rg -n 'func \(s \*Service\) (PublishData|UnsetProperty|GetData)\b' internal/appengine/; echo "---"; rg -n 'PublishServerValue|PublishDeviceValue|UnsetServerProperty' internal/appengine/
internal/appengine/data.go:109:func (s *Service) GetData(ctx context.Context, realm, deviceID, ifaceName, path string, opts QueryOpts) (any, error) {
internal/appengine/data.go:466:func (s *Service) PublishData(ctx context.Context, realm, deviceID, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
internal/appengine/data.go:479:func (s *Service) UnsetProperty(ctx context.Context, realm, deviceID, ifaceName, path string) error {
---
internal/appengine/data.go:474:	return s.sd.PublishServerValue(ctx, realm, id, ifaceName, path, value, ts)
internal/appengine/data.go:487:	return s.sd.UnsetServerProperty(ctx, realm, id, ifaceName, path)
internal/appengine/http_test.go:46:func (f *fakeServerData) PublishServerValue(_ context.Context, realm string, id deviceid.ID, iface, path string, value json.RawMessage, _ *time.Time) error {
internal/appengine/http_test.go:53:func (f *fakeServerData) UnsetServerProperty(_ context.Context, realm string, id deviceid.ID, iface, path string) error {
internal/appengine/service.go:73:	PublishServerValue(ctx context.Context, realm string, id deviceid.ID, iface, path string, value json.RawMessage, ts *time.Time) error
internal/appengine/service.go:74:	UnsetServerProperty(ctx context.Context, realm string, id deviceid.ID, iface, path string) error

→ Read internal/appengine/data.go [offset=100, limit=90]
→ Read internal/appengine/data.go [offset=460, limit=30]

$ rg -n 'func WriteDataWithMetadata|func WriteData\b|func WriteDataWithLinks|func WriteRawErrors' pkg/astarteapi/*.go
pkg/astarteapi/envelope.go:126:func WriteData(w http.ResponseWriter, status int, v any) error {
pkg/astarteapi/envelope.go:145:func WriteDataWithLinks(w http.ResponseWriter, status int, v any, links Links) error {
pkg/astarteapi/envelope.go:156:func WriteDataWithMetadata(w http.ResponseWriter, status int, v any, m map[string]any) error {
pkg/astarteapi/envelope.go:186:func WriteRawErrors(w http.ResponseWriter, status int, errs any) error {

$ rg -n 'metadata' docs/api/astarte_appengine_api.yaml
1534:        flat objects; `table` returns a `{data, metadata}` envelope with rows

$ rg -n 'Tabular|Metadata' internal/appengine/data.go | head -40
54:	RetrieveMetadata     bool
59:// Tabular carries a format=table payload plus its metadata object; serveData
61:type Tabular struct {
63:	Metadata map[string]any `json:"metadata"`
262:		return &Tabular{
264:			Metadata: map[string]any{
318:		return &Tabular{
320:			Metadata: map[string]any{"columns": columns, "table_header": header},

→ Read internal/appengine/service.go [offset=100, limit=40]
→ Read internal/appengine/service.go [offset=242, limit=50]
→ Read docs/site/appengine-api.md
→ Read pkg/astarteapi/envelope.go [offset=120, limit=80]
→ Read internal/appengine/data.go [offset=190, limit=150]
→ Read docs/api/astarte_appengine_api.yaml [offset=1, limit=62]

$ grep -niE 'metadata|unexpected_keys|GroupCreate|bucket_duration|required mapping|Missing required' .mule/todo.md | head -40
136:- [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
186:- [x] astarteapi-metadata-envelope-test [auto]: add a container-free golden test for `WriteDataWithMetadata` (pkg/astarteapi/envelope.go:154-158) pinning the two wire branches — nil metadata omits the `metadata` key (`{"data": v}` only, via the `omitempty` at envelope.go:151), populated metadata emits `{"data": v, "metadata": {...}}` — with 200 + the package `Content-Type`, matching the TestGoldenEnvelopes style. Today the non-nil branch is only covered by the integration-tagged formats_parity_test.go (runs on the Legion Go, not the Pi gate) and the nil-omission branch is covered nowhere (internal/appengine/data.go:262-268 always sets metadata), so a regression dropping the omitempty passes every suite on this box.
208:- [!] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine_api.yaml — `retrieve_metadata` and `downsample_key` are parsed for every data read by `parseQueryOpts` (internal/appengine/http.go:624 `retrieve_metadata`, 633 `downsample_key`) yet appear on none of the read endpoints: device-scoped GET `{interface}` (yaml:583-606) and `{interface}/{path}` (yaml:623-649), by-alias (yaml:324-368, 369-417), and in-group (yaml:1108-1132, 1149-1176) — the parameter block on each lists only since/since_after/to/limit/downsample_to/sort/format/allow_bigintegers/allow_safe_bigintegers (components DataSince..DataAllowSafeBigIntegers, yaml:1360-1439). Add the two `$ref`s (new DataRetrieveMetadata and DataDownsampleKey components, or inline — mirror the existing style) to all six, and update the DataSort-independent `sort` default note only if verified. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
209:- [x] examples-echo-container-contract-test [auto]: add a container-free httptest suite for examples/flow-container-echo — the only no-test package in the repo that contains logic (`handleMessage`, main.go:49-75; `sanitizeLog`, main.go:40-47) — pinning (1) POST `/v1/message` echoes a valid JSON-object body verbatim with 200, answers 400 on invalid JSON and 204 on an empty body (the 1 MiB read cap), and (2) metadata `echo_drop` == "1" forces the 204 filter/drop path; the example is the canonical reference container contract (README.md, docs/handoff/flow-design-b-container-block-2026-07-29.md), so a regression silently misleads every container author who copies it.
256:- [x] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
264:- [!] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implements `device_empty_cache_received` end to end — `SimpleEvent` field 23, an empty `DeviceEmptyCacheReceivedEvent` whose encoder emits exactly `{"type":"device_empty_cache_received"}` (astarte_core encoder.ex:54-62), fired by DUP at `control_handler.ex:161` *inside* the `with` that runs `send_control_consumer_properties` → `resend_all_properties` → `set_pending_empty_cache`, so it fires only when all three succeed and never on failure, and the commit deletes the `# TODO: implement on_empty_cache_received` from `astarte_events/lib/astarte_events/triggers/core.ex`. Astrate already accepts the trigger but never evaluates it: `OnDeviceEmptyCacheReceived` sits in `deviceOns` with value `false` (internal/engine/triggers/match.go:123, the `evaluated` flag read at match.go:406) and the comment's stated reason ("upstream defines no SimpleEvent variant for it") is now false. Add the event body constructor to internal/engine/triggers/events.go (the payload is `{"type":"device_empty_cache_received"}` and nothing else — no device id, no metadata), flip the `deviceOns` value to true, and fire it from `handleEmptyCache` (internal/engine/control.go:58) with the same idiom `handleLifecycle` uses for `device_connected` (`e.fireDevice` + `stream.Publish`, internal/engine/engine.go:346-353), at the point where `m.Ack()` runs — i.e. after `resendServerProperties` and `sendConsumerProperties` have both succeeded inside `retryStore`. Add a T1 test that a trigger on it fires exactly once per accepted emptyCache and does not fire when the resend fails. — BLOCKED: wrote nothing
265:- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
268:- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
272:- [x] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `components.parameters.RealmName` and `RealmCreate.realm_name` both carry `pattern: '^[a-z][a-z0-9]*$'` and `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go) holds both against the realm name CHECK extracted from migrations/000002_metadata.up.sql — but the read side is still bare `type: string`: `Realm.realm_name` (components.schemas.Realm) and `DataEnvelopeRealmNames.items`. Decide whether a response field earns the same pattern (a generated client learns nothing new from it, but a reader diffing create against get sees one constrained and one not) and, if it does, add it to both and extend that test to hold the response-side pattern against the same CHECK. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
283:- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
284:- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
285:- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
290:- [!] docs-native-socket-event-schema [auto]: document the native socket's event payload in docs/api/astrate_native_api.yaml — currently only a fenced `incoming_data` example inside the operation description (yaml:300-310), with no `components.schemas` entry anywhere (yaml:425-478), and that example is not the general shape: `wireEvent` marks `interface`, `path` and `value` `omitempty` (internal/appengine/stream/ws.go:51-58), so measured on a verbatim copy a lifecycle event marshals to `{"event":"device_connected","realm":"test","device_id":"dev1","timestamp":"2024-05-30T13:49:57.045Z"}` — no interface/path/value at all (and `value: false`/`value: 0` ARE emitted, because omitempty on an `any` drops only nil). Add a `NativeEvent` schema (event, realm, device_id, optional interface/path/value, timestamp) referenced from the operation, state that `event` is one of the internal/engine/stream/bus.go Kind* constants (bus.go:20-26), and record that this socket drops `InterfaceMajor`, `IP`, `ErrorName` and `ErrorMetadata`, which `stream.Event` carries (bus.go:39-65) — so a `device_error` arrives here with no reason, unlike the Phoenix twin, which marshals the whole struct (internal/appengine/channels/ws.go:406-416). That last sentence documents a code gap, not a spec one: do not add the fields; if they are wanted on the wire, that is a separate code task. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1
297:- [x] docs-sync-pairing-realm-name-pattern [auto]: constrain `components.parameters.RealmName` in docs/api/astarte_pairing_api.yaml (yaml:359-365, a bare `type: string` described only as "The realm name.") with `pattern: '^[a-z][a-z0-9]*$'` — the realm name is CHECK-constrained to exactly that (migrations/000002_metadata.up.sql:7, read-only, never edit) and the two surfaces whose tasks already ran carry it on their twins (docs/api/astarte_housekeeping_api.yaml:246, docs/api/astrate_native_api.yaml:420 and 467). Because this one parameter is shared by all seven operations, the description must record that an off-pattern name gets three different answers, not one: the agent and device routes answer **401** (the realm is resolved by the middleware, which treats a miss as unauthenticated rather than 404 — internal/auth/middleware.go:66-72); `getHealth` answers **404** (service.Health returns store.ErrNotFound → WriteNotFound, internal/pairing/service.go:369-371, http.go:289-290); and `getVersion` answers **200** `{"data":"1.2.2"}` for any realm string whatsoever, because `GET /pairing/v1/{realm}/version` is registered with no middleware and no realm lookup (cmd/astrate/main.go:438 → observability.VersionHandler, internal/observability/compat.go:38-47) — say that in the version operation so nobody reads the shared parameter as a guarantee the handler enforces. Mirror `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go:1247) with a pairing twin only if it stays a few lines; the YAML edit is the deliverable. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).

$ grep -niE 'DataEnvelope|unexpected|group_name.*devices|downsample_to.*bucket|metadata' .mule/todo.md | head; echo "---forgiulio---"; grep -niE 'bucket_duration|downsample_to' .mule/for-giulio.md | head
1:- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
128:- [x] docs-sync-appengine-group-endpoints [auto]: add the 8 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (returns {data: {group_name: "<name>"}}, http.go:428-434), PATCH /groups/{group}/devices/{device} (DevicePatch body → updated status, http.go:518-537), and the 6 group interface-data routes: the interface-level GET and the path-level GET/PUT/POST/DELETE on /groups/{group}/devices/{device}/interfaces (http.go:75-80, handlers listInterfacesInGroup, getDataInGroup, putDataInGroup, deleteDataInGroup at http.go:187-377). Also add GET /devices/{device}/interfaces (listDeviceInterfaces, http.go:48) which is missing from the spec — that is the 9th route. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
130:- [x] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups/{group}/devices is missing `from_token` and `limit` pagination params that code reads at http.go:439-441 (and returns via links.next), (2) POST /groups is missing 422 validation-error response (FieldErrors for blank group_name/empty devices, http.go:404-416), (3) PATCH /devices/{device} and PATCH /devices-by-alias/{alias} are missing 422 ("Attribute key not found" via ErrAttributeKeyNotFound, service.go:352-357 → http.go:656-657), and the device-scoped PATCH is also missing 409 ("Alias already in use", http.go:651) and its documented 500 covers the merge-patch content-type mismatch which code answers with an unmapped 500 too, (4) GET /devices/{device}/interfaces/{interface}/{path} is missing 422 (invalid since/since_after/to/limit/downsample_to/format params, http.go:554-639). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
136:- [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
154:- [x] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmmanagement/health`, `GET /pairing/health` — in docs/api/astrate_native_api.yaml, and widen the shared `DataEnvelopeHealthStatus.status` enum (yaml:301-303) from `[ok]` to `[ok, unhealthy]` — the handler answers 503 `{"data":{"status":"unhealthy"}}` when the readiness probe fails (internal/observability/compat.go:19-34, deviation 18), but the spec documents only 200 with `status: ok` for all three paths. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
186:- [x] astarteapi-metadata-envelope-test [auto]: add a container-free golden test for `WriteDataWithMetadata` (pkg/astarteapi/envelope.go:154-158) pinning the two wire branches — nil metadata omits the `metadata` key (`{"data": v}` only, via the `omitempty` at envelope.go:151), populated metadata emits `{"data": v, "metadata": {...}}` — with 200 + the package `Content-Type`, matching the TestGoldenEnvelopes style. Today the non-nil branch is only covered by the integration-tagged formats_parity_test.go (runs on the Legion Go, not the Pi gate) and the nil-omission branch is covered nowhere (internal/appengine/data.go:262-268 always sets metadata), so a regression dropping the omitempty passes every suite on this box.
196:- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
208:- [!] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine_api.yaml — `retrieve_metadata` and `downsample_key` are parsed for every data read by `parseQueryOpts` (internal/appengine/http.go:624 `retrieve_metadata`, 633 `downsample_key`) yet appear on none of the read endpoints: device-scoped GET `{interface}` (yaml:583-606) and `{interface}/{path}` (yaml:623-649), by-alias (yaml:324-368, 369-417), and in-group (yaml:1108-1132, 1149-1176) — the parameter block on each lists only since/since_after/to/limit/downsample_to/sort/format/allow_bigintegers/allow_safe_bigintegers (components DataSince..DataAllowSafeBigIntegers, yaml:1360-1439). Add the two `$ref`s (new DataRetrieveMetadata and DataDownsampleKey components, or inline — mirror the existing style) to all six, and update the DataSort-independent `sort` default note only if verified. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
209:- [x] examples-echo-container-contract-test [auto]: add a container-free httptest suite for examples/flow-container-echo — the only no-test package in the repo that contains logic (`handleMessage`, main.go:49-75; `sanitizeLog`, main.go:40-47) — pinning (1) POST `/v1/message` echoes a valid JSON-object body verbatim with 200, answers 400 on invalid JSON and 204 on an empty body (the 1 MiB read cap), and (2) metadata `echo_drop` == "1" forces the 204 filter/drop path; the example is the canonical reference container contract (README.md, docs/handoff/flow-design-b-container-block-2026-07-29.md), so a regression silently misleads every container author who copies it.
264:- [!] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implements `device_empty_cache_received` end to end — `SimpleEvent` field 23, an empty `DeviceEmptyCacheReceivedEvent` whose encoder emits exactly `{"type":"device_empty_cache_received"}` (astarte_core encoder.ex:54-62), fired by DUP at `control_handler.ex:161` *inside* the `with` that runs `send_control_consumer_properties` → `resend_all_properties` → `set_pending_empty_cache`, so it fires only when all three succeed and never on failure, and the commit deletes the `# TODO: implement on_empty_cache_received` from `astarte_events/lib/astarte_events/triggers/core.ex`. Astrate already accepts the trigger but never evaluates it: `OnDeviceEmptyCacheReceived` sits in `deviceOns` with value `false` (internal/engine/triggers/match.go:123, the `evaluated` flag read at match.go:406) and the comment's stated reason ("upstream defines no SimpleEvent variant for it") is now false. Add the event body constructor to internal/engine/triggers/events.go (the payload is `{"type":"device_empty_cache_received"}` and nothing else — no device id, no metadata), flip the `deviceOns` value to true, and fire it from `handleEmptyCache` (internal/engine/control.go:58) with the same idiom `handleLifecycle` uses for `device_connected` (`e.fireDevice` + `stream.Publish`, internal/engine/engine.go:346-353), at the point where `m.Ack()` runs — i.e. after `resendServerProperties` and `sendConsumerProperties` have both succeeded inside `retryStore`. Add a T1 test that a trigger on it fires exactly once per accepted emptyCache and does not fire when the resend fails. — BLOCKED: wrote nothing
---forgiulio---

$ grep -nE '^- \[.\] docs-sync-(appengine|ae)-' .mule/todo.md | sed 's/\(.\{160\}\).*/\1/'
127:- [x] docs-sync-appengine-by-alias-endpoints [auto]: add the 7 undocumented by-alias mirror routes to docs/api/astarte_appengine_api.yaml — PATCH /devices-b
128:- [x] docs-sync-appengine-group-endpoints [auto]: add the 8 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (ret
129:- [x] docs-sync-appengine-get-group-device [auto]: add the undocumented `GET /groups/{group}/devices/{device}` route to docs/api/astarte_appengine_api.yaml 
130:- [x] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups
132:- [x] docs-sync-appengine-group-patch-status [auto]: PATCH /groups/{group}/devices/{device} (internal/appengine/http.go:518-537, patchGroupDevice) shares ap
133:- [x] docs-sync-appengine-data-422-interface-level [auto]: the interface-level data GETs — GET /devices/{device}/interfaces/{interface}, GET /devices-by-ali
136:- [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accept
143:- [x] docs-sync-appengine-device-status-schema [auto]: fix the `DeviceStatus` schema in docs/api/astrate_appengine_api.yaml — the wire emits `id` (internal/
144:- [!] docs-sync-appengine-data-set-422 [auto]: add the missing `422` response to the six PUT/POST data-set operations in docs/api/astrate_appengine_api.yaml
145:- [x] docs-sync-appengine-downsample-min [auto]: fix the `downsample_to` constraint in docs/api/astrate_appengine_api.yaml — the DataDownsample description 
167:- [x] docs-sync-ae-write-405 [auto]: add the missing `405` response to the nine data-write ops in docs/api/astarte_appengine_api.yaml — PUT/POST/DELETE `/de
168:- [x] docs-sync-ae-write-value-422 [auto]: add the missing `422` response to the six PUT/POST data-write ops in docs/api/astarte_appengine_api.yaml (device-
169:- [x] docs-sync-ae-post-groups-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups` in docs/api/astarte_appengine_api.yaml — `S
172:- [x] docs-sync-ae-post-group-devices-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_
173:- [x] docs-sync-ae-patch-merge-patch-media-type [auto]: fix the three PATCH request bodies in docs/api/astarte_appengine_api.yaml — PATCH /devices/{device} 
208:- [!] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine
236:- [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml —
237:- [x] docs-sync-ae-patch-by-alias-409 [auto]: add the missing `409` response to `patchDeviceByAlias` in docs/api/astarte_appengine_api.yaml (yaml:241-291, c
238:- [x] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_ap
239:- [x] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data DELETE operations in docs/api/astarte_appengine_api.yaml (`delet
240:- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (ya
260:- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` shows

## Objective
- Carry out `.mule/recipes/docs-sync.md` as a *proposal* job: find docs/code drift on one API surface, append ≤5 new task lines to `.mule/todo.md`, and escalate `docs/site/` contradictions to `.mule/for-giulio.md`. Do not implement proposed tasks; do not touch git.

## Important Details
- Output of this recipe = task lines in `.mule/todo.md` (+ any evidence file). No code changes.
- Hard constraints from user: no git (no commit/branch/checkout/stash); never `gh issue edit/close/comment`; only run `gh issue create` if a recipe explicitly names it as the deliverable (docs-sync does NOT — deliverable is task lines, not an issue). Never both issue + task line for same work.
- MULE.md rules: everything in English; ~30k-token useful context; read only named files; use `rg`; stop when you have a working answer; if bigger than it looked, do the first useful part and split the rest into `.mule/todo.md`.
- docs-sync recipe specifics: prose in `docs/site/` is on the never-touch list (Giulio's voice) — propose only; `docs/api/*.yaml` are generated artefacts and ARE fixable. Fixing a wrong path/status/field there is a normal task.
- Recipe loop: pick ONE surface (appengine, housekeeping, pairing, realm-management, astrate-native), compare handler routes vs documented paths, then hand-check 3-4 endpoints' status codes/response fields. Max 5 proposals. `docs/site/` contradictions and stale config keys (`rg -o '\bASTRATE_[A-Z_]+' -N internal/ | sort -u`) go to `.mule/for-giulio.md` (never edit the page).
- Site-build gate: only if a task touches `docs/api/` or `docs/mkdocs.yml`, run `make -C docs build` (mkdocs not installed on the Pi — may need `pip install -r docs/requirements.txt`, or at minimum confirm YAML still parses).
- Rotation evidence gathered: rm run 2026-10-02, housekeeping 2026-10-04, astrate-native 2026-10-05 and 2026-10-06, pairing 2026-10-07 (spec accurate, nothing proposed) and 2026-10-08. → **appengine is the least-recently-swept surface; chosen for this run.**

## Work State
### Completed
- Read `.mule/MULE.md`, `.mule/recipes/docs-sync.md`, `.mule/todo.md`, `.mule/for-giulio.md`.
- Determined rotation; selected **appengine** surface.
- Route-set comparison: `internal/appengine/http.go:41` `Mount` registers 18 paths; `docs/api/astarte_appengine_api.yaml` documents paths at lines 26, 67, 141, 181, 259, 347, 381, 428, 625, 655, 697, 874, 952, 982, 1078, 1183, 1213, 1256. **No route gap** (all code routes documented, no dead paths). `/version` documented at yaml:26.
- Read handlers (`http.go` 41-560+, `data.go`, `service.go`) and spec sections for devices/stats/groups/data.

### Active
- Hand-checking status codes/fields for appengine endpoints. Candidate drifts found so far (need triage to ≤5 tasks):
  1. **ALREADY QUEUED — do NOT re-propose:** `docs-sync-ae-list-devices-422` sits at `.mule/todo.md:240` as `- [!] ... BLOCKED: TIMEOUT`. It covers `GET /devices` missing its 422 for an invalid `from_token` (`internal/appengine/service.go:141-155`, `ErrValidation: invalid cursor`).
  2. `DataEnvelopeObject` (appengine spec) is missing the `metadata` property. Code returns `{"data":..., "metadata":...}` for `format=table` via `astarteapi.WriteDataWithMetadata` (`pkg/astarteapi/envelope.go:156`) / `Tabular` struct (`internal/appengine/data.go:59-63`). Only referenced in the `DataFormat` description at `docs/api/astarte_appengine_api.yaml:1534`. Affects the six data-GET 200 responses.
  3. 422 on `putData`/`publishData` (+by-alias, +in-group) is documented only as `ValueTooLarge`, but after `appengine-missing-required-422` merged a required-mapping miss now returns 422 `ErrorDetail` "Missing required mapping key". Also the 400 can carry an `unexpected_keys` array (written via `WriteRawErrors`, `pkg/astarteapi/envelope.go:186`), which the `ErrorDetail`/`BadRequest` example omits.
  4. `GroupCreate` schema marks only `group_name` required, but the code requires `devices` (422 on missing → `{"devices":["can't be blank"]}`, on empty → `["should have at least 1 item(s)"]` — see `internal/appengine/service.go` CreateGroup path). Wrong required-ness.
  5. Group routes' 404 use generic `NotFound` (example `detail: Not Found`), but code returns `"Group not found"` (`ErrGroupNotFound`, `internal/appengine/service.go:42`; `missingGroup`, `service.go:65-67`).
  6. `docs/site/appengine-api.md:60` `&downsample_to=<bucket_duration>` and `:65` "maps onto Timescale `time_bucket()`" contradict code: `downsample_to` is a point count (`internal/appengine/http.go:588-598`, must be > 2; spec `DataDownsample` "Downsample to this many data points"). → for-giulio.md candidate.
- Do NOT duplicate already-escalated for-giulio items: appengine-api.md envelope-format contradiction and `a_ch` room-filters contradiction (from the 2026-10-06 astrate-native run).
- Key error taxonomy read: `internal/appengine/service.go:34-67` (`ErrInvalidAlias`, `ErrAliasTagNotFound`, `ErrInvalidAttributes`, `ErrAttributeKeyNotFound`, `ErrGroupNotFound`, `ErrGroupAlreadyExists`, `ErrDeviceAlreadyInGroup`, `ErrValidation`, `missingGroup`); `DeviceStatus` at `service.go:105-121`, `deviceStatusView` at `service.go:254-279`.

### Blocked
- (none yet) — nothing failing; only open question is which ≤5 candidates to promote to task lines.

## Next Move
1. Finish triaging the six candidates: finalize the ≤5 task lines (exclude candidate 1 — already queued), draft each as `docs-sync-ae-<slug> [auto]: <outcome>` naming the yaml file+line and the code `file:line` that proves it.
2. Append approved task lines to `.mule/todo.md`; append the `downsample_to` site contradiction to `.mule/for-giulio.md` (newest at top, evidence inline); write the run report to `.mule/reports/log/<UTC>-recipe-docs-sync-proposed.md`; do not touch git and do not run `make -C docs build` (no `docs/` files change in a proposal run).

## Relevant Files
- `.mule/recipes/docs-sync.md`: the recipe being executed (surface rotation, rules, gate).
- `.mule/MULE.md`: standing safety envelope + Giulio's standing instructions (issue #1 never raised, etc.).
- `.mule/todo.md`: deliverable target; `- [!] docs-sync-ae-list-devices-422` already at line 240 (BLOCKED) — do not re-propose.
- `.mule/for-giulio.md`: escalation channel (newest-at-top, one line each, file:line evidence inline).
- `.mule/reports/log/`: where the run's `*-recipe-docs-sync-proposed.md` report goes; recent rotation evidence lives here.
- `docs/api/astarte_appengine_api.yaml`: the spec under review (paths, `DataEnvelopeObject`, `GroupCreate`, `ValueTooLarge`, `NotFound` components).
- `internal/appengine/http.go`: `Mount` (line 41) route set; query-param validation incl. `downsample_to` (588-598); `writeError`.
- `internal/appengine/service.go`: error sentinels (34-67), `ListDevices` (141-179), `ListGroupDevices` (541-597), `DeviceStatus`/`deviceStatusView` (105-121, 254-279).
- `internal/appengine/data.go`: `GetData` (109), `Tabular` (59-63), `PublishData` (466), `UnsetProperty` (479).
- `pkg/astarteapi/envelope.go`: `WriteData` (126), `WriteDataWithLinks` (145), `WriteDataWithMetadata` (156), `WriteError` (163), `WriteRawErrors` (186).
- `docs/site/appengine-api.md`: prose page; `downsample_to` contradiction at lines 60/65 (escalate, never edit).
- `docs/Makefile` (invoked as `make -C docs build`): site build gate, only for tasks touching `docs/api/` or `docs/mkdocs.yml`.
Continue if you have next steps, or stop and ask for clarification if you are unsure how to proceed.
→ Read .mule/todo.md [offset=127, limit=48]
→ Read .mule/todo.md [offset=236, limit=30]
→ Read internal/appengine/http.go [offset=640, limit=90]

$ cd /root/astrate-mule && grep -nE 'DataEnvelopeObject|BadRequest|ValueTooLarge|ValidationErrors|NotFound:|GroupNotFound|DeviceNotFound' docs/api/astarte_appengine_api.yaml | head -80
204:          $ref: "#/components/responses/DeviceNotFound"
245:          $ref: "#/components/responses/BadRequest"
251:          $ref: "#/components/responses/DeviceNotFound"
287:          $ref: "#/components/responses/DeviceNotFound"
333:          $ref: "#/components/responses/BadRequest"
339:          $ref: "#/components/responses/DeviceNotFound"
377:          $ref: "#/components/responses/DeviceNotFound"
416:                $ref: "#/components/schemas/DataEnvelopeObject"
422:          $ref: "#/components/responses/DeviceNotFound"
424:          $ref: "#/components/responses/ValidationErrors"
468:                $ref: "#/components/schemas/DataEnvelopeObject"
474:          $ref: "#/components/responses/DeviceNotFound"
476:          $ref: "#/components/responses/ValidationErrors"
518:          $ref: "#/components/responses/BadRequest"
524:          $ref: "#/components/responses/DeviceNotFound"
528:          $ref: "#/components/responses/ValueTooLarge"
570:          $ref: "#/components/responses/BadRequest"
576:          $ref: "#/components/responses/DeviceNotFound"
580:          $ref: "#/components/responses/ValueTooLarge"
613:          $ref: "#/components/responses/BadRequest"
619:          $ref: "#/components/responses/DeviceNotFound"
649:          $ref: "#/components/responses/DeviceNotFound"
685:                $ref: "#/components/schemas/DataEnvelopeObject"
691:          $ref: "#/components/responses/DeviceNotFound"
693:          $ref: "#/components/responses/ValidationErrors"
730:                $ref: "#/components/schemas/DataEnvelopeObject"
736:          $ref: "#/components/responses/DeviceNotFound"
738:          $ref: "#/components/responses/ValidationErrors"
775:          $ref: "#/components/responses/BadRequest"
781:          $ref: "#/components/responses/DeviceNotFound"
785:          $ref: "#/components/responses/ValueTooLarge"
822:          $ref: "#/components/responses/BadRequest"
828:          $ref: "#/components/responses/DeviceNotFound"
832:          $ref: "#/components/responses/ValueTooLarge"
860:          $ref: "#/components/responses/BadRequest"
866:          $ref: "#/components/responses/DeviceNotFound"
940:          $ref: "#/components/responses/BadRequest"
948:          $ref: "#/components/responses/ValidationErrors"
1028:          $ref: "#/components/responses/ValidationErrors"
1064:          $ref: "#/components/responses/BadRequest"
1074:          $ref: "#/components/responses/ValidationErrors"
1146:          $ref: "#/components/responses/BadRequest"
1244:                $ref: "#/components/schemas/DataEnvelopeObject"
1252:          $ref: "#/components/responses/ValidationErrors"
1290:                $ref: "#/components/schemas/DataEnvelopeObject"
1298:          $ref: "#/components/responses/ValidationErrors"
1336:          $ref: "#/components/responses/BadRequest"
1346:          $ref: "#/components/responses/ValueTooLarge"
1384:          $ref: "#/components/responses/BadRequest"
1394:          $ref: "#/components/responses/ValueTooLarge"
1423:          $ref: "#/components/responses/BadRequest"
1711:    DataEnvelopeObject:
1765:    BadRequest:
1795:    NotFound:
1805:    DeviceNotFound:
1835:    ValidationErrors:
1846:    ValueTooLarge:

→ Read docs/api/astarte_appengine_api.yaml [offset=1690, limit=204]
→ Read docs/api/astarte_appengine_api.yaml [offset=440, limit=150]
→ Read internal/appengine/data.go [offset=50, limit=20]
→ Read pkg/astarteapi/envelope.go [offset=140, limit=60]

$ cd /root/astrate-mule && grep -nE 'group_name|devices|required:|can.t be blank|at least 1 item' docs/api/astarte_appengine_api.yaml | sed -n '1,60p'
51:                required: [data]
67:  /appengine/v1/{realm}/devices:
70:      summary: List devices in a realm
94:          description: Maximum number of devices to return.
97:          description: Paginated list of devices.
111:                      self: "/v1/test/devices?details=false&limit=1"
132:                      self: "/v1/test/devices?details=true&limit=1"
133:                      next: "/v1/test/devices?details=true&from_token=dT6hS2W9TT6LEnP25ks_lg&limit=1"
141:  /appengine/v1/{realm}/stats/devices:
158:                required: [data]
162:                    required: [connected_devices, total_devices]
164:                      connected_devices:
167:                      total_devices:
172:                  connected_devices: 42
173:                  total_devices: 150
181:  /appengine/v1/{realm}/devices/{device}:
221:        required: true
226:              required: [data]
259:  /appengine/v1/{realm}/devices-by-alias/{alias}:
271:          required: true
304:          required: true
309:        required: true
314:              required: [data]
347:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces:
361:          required: true
381:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
396:          required: true
428:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
442:          required: true
449:          required: true
493:          required: true
500:          required: true
505:        required: true
510:              required: [data]
545:          required: true
552:          required: true
557:        required: true
562:              required: [data]
598:          required: true
605:          required: true
625:  /appengine/v1/{realm}/devices/{device}/interfaces:
655:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
697:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
711:          required: true
757:          required: true
762:        required: true
767:              required: [data]
804:          required: true
809:        required: true
814:              required: [data]
852:          required: true
906:        Creates a new device group with an optional initial list of devices.
914:        required: true
919:              required: [data]
925:                group_name: sensors
926:                devices:
935:                required: [data]
972:                  group_name: sensors
982:  /appengine/v1/{realm}/groups/{group}/devices:
985:      summary: List devices in a group

$ cd /root/astrate-mule && grep -nE 'func .*CreateGroup|func .*GetGroup\b|func .*ListGroups|ErrGroupNotFound|missingGroup|"devices"|should have at least|can.t be blank' internal/appengine/service.go
41:// ErrGroupNotFound marks a missing group (maps to 404 "Group not found").
42:var ErrGroupNotFound = errors.New("Group not found") //nolint:staticcheck // ST1005: upstream wire text
64:// missingGroup wraps a failed group resolution so it satisfies BOTH
65:// ErrGroupNotFound (→ 404 "Group not found") and store.ErrNotFound.
66:func missingGroup(name string) error {
67:	return fmt.Errorf("%w: %w: group %q", ErrGroupNotFound, store.ErrNotFound, name)
382:// PATCH /groups/{group}/devices/{device}). Unknown group → ErrGroupNotFound;
403:// ordering: unknown group first (ErrGroupNotFound), then membership — a
412:		return 0, deviceid.ID{}, nil, missingGroup(groupName)
462:func (s *Service) CreateGroup(ctx context.Context, realm, name string, devices []string) error {
473:			fe.addf("devices", "must exist (%s not found)", raw)
478:				fe.addf("devices", "must exist (%s not found)", raw)
508:func (s *Service) GetGroup(ctx context.Context, realm, name string) error {
514:		return missingGroup(name)
520:func (s *Service) ListGroups(ctx context.Context, realm string) ([]string, error) {
561:		return nil, missingGroup(name)

→ Read internal/appengine/service.go [offset=462, limit=60]
→ Read docs/api/astarte_appengine_api.yaml [offset=900, limit=150]

$ cd /root/astrate-mule && grep -nE 'GroupCreate:|GroupName:|responses/NotFound|\$ref: "#/components/responses/NotFound"' docs/api/astarte_appengine_api.yaml; echo ---; grep -n 'GroupCreate' docs/api/astarte_appengine_api.yaml
978:          $ref: "#/components/responses/NotFound"
1026:          $ref: "#/components/responses/NotFound"
1070:          $ref: "#/components/responses/NotFound"
1104:          $ref: "#/components/responses/NotFound"
1152:          $ref: "#/components/responses/NotFound"
1179:          $ref: "#/components/responses/NotFound"
1209:          $ref: "#/components/responses/NotFound"
1250:          $ref: "#/components/responses/NotFound"
1296:          $ref: "#/components/responses/NotFound"
1342:          $ref: "#/components/responses/NotFound"
1390:          $ref: "#/components/responses/NotFound"
1429:          $ref: "#/components/responses/NotFound"
1470:    GroupName:
1652:    GroupCreate:
---
922:                  $ref: "#/components/schemas/GroupCreate"
938:                    $ref: "#/components/schemas/GroupCreate"
1652:    GroupCreate:

→ Read internal/appengine/http.go [offset=388, limit=40]
→ Read docs/api/astarte_appengine_api.yaml [offset=1652, limit=40]

$ cd /root/astrate-mule && rg -n 'WriteNotFound|DetailNotFound|WriteDeviceNotFound' internal/appengine pkg/astarteapi/envelope.go
pkg/astarteapi/envelope.go:42:	// DetailNotFound is the canonical generic 404 detail.
pkg/astarteapi/envelope.go:43:	DetailNotFound = "Not Found"
pkg/astarteapi/envelope.go:49:	// DetailNotFound, which a handler emits for a route that matched and found
pkg/astarteapi/envelope.go:205:// WriteNotFound writes the canonical generic 404 envelope.
pkg/astarteapi/envelope.go:206:func WriteNotFound(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:207:	return WriteError(w, http.StatusNotFound, DetailNotFound)
pkg/astarteapi/envelope.go:210:// WriteDeviceNotFound writes the canonical 404 "Device not found" envelope.
pkg/astarteapi/envelope.go:211:func WriteDeviceNotFound(w http.ResponseWriter) error {
internal/appengine/http.go:683:		_ = astarteapi.WriteDeviceNotFound(w)

→ Read internal/appengine/http.go [offset=505, limit=40]
→ Read docs/site/appengine-api.md [offset=40, limit=40]
→ Read internal/appengine/http.go [offset=580, limit=45]
→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
→ Read .mule/for-giulio.md [limit=50]
→ Read .mule/todo.md [offset=266, limit=40]

$ cd /root/astrate-mule && grep -nE 'DataDownsample:|Downsample to|bucketFor|func .*Downsample' docs/api/astarte_appengine_api.yaml internal/appengine/data.go internal/store/*.go pkg/*/*.go 2>/dev/null | head -20
docs/api/astarte_appengine_api.yaml:1509:    DataDownsample:
docs/api/astarte_appengine_api.yaml:1514:      description: Downsample to this many data points (must be > 2).
internal/appengine/data.go:187:		bucket := bucketFor(last.Sub(first), opts.DownsamplePoints)
internal/store/datastreams.go:326:func (s *Store) Downsample(ctx context.Context, q SeriesQuery, bucket time.Duration) ([]DownsamplePoint, error) {
internal/store/datastreams.go:375:func (s *Store) DownsampleLTTB(ctx context.Context, q SeriesQuery, points int) ([]DownsamplePoint, error) {

I have enough to finalize. Appending the four appengine drift proposals to `.mule/todo.md` and the site contradiction to `.mule/for-giulio.md`.
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -295,4 +295,8 @@
 - [x] docs-sync-pairing-register-404-unreachable [auto]: drop the `404` `DeviceNotFound` response from `registerDevice` in docs/api/astarte_pairing_api.yaml (yaml:140-141) — it is documented but absent: the route is wrapped by `requireAgent` (internal/pairing/http.go:82-83), i.e. `mw.RequireRealm(auth.ClaimPairing)`, which resolves the realm *before* the handler and answers `401 {"errors":{"detail":"Unauthorized"}}` for an unknown one (internal/auth/middleware.go:66-72, "no existence oracle on auth failures"), so `Service.Register`'s own `GetRealmByName` (internal/pairing/service.go:183-186) can return `store.ErrNotFound` only if the realm is deleted between the middleware's lookup and the handler's — a race, not a contract. Nothing else in Register returns `store.ErrNotFound`: `RegisterDevice`'s `ErrDeviceAlreadyConfirmed` becomes 422 (service.go:209-213) and the CountDevices/SetPayloadFormatHint failures are DB errors → 500. No test covers it (internal/pairing/http_test.go asserts 404 only for the agent DELETE at :527 and for health at :595-598), and the `DeviceNotFound` component's own description — "unknown device ID" — does not match what that path would be even if it fired. The premise came from docs-sync-pairing-register-404, which read service.go without the middleware in front of it. Keep the `404` on `unregisterDevice` (reachable: unknown device, and a malformed `deviceID` via service.go:231-234) and on `getHealth`. State the reachability argument in the report; if you find a reachable path instead, keep the response and fix its description rather than deleting it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
 - [!] docs-sync-pairing-body-cap-400 [auto]: document the 64 KiB request-body cap on the three body-taking pairing ops — `registerDevice`, `requestCredentials` and `verifyCredentials` all decode through `astarteapi.DecodeData(r.Body, maxBodyBytes, …)` with `maxBodyBytes = 64 << 10` (internal/pairing/http.go:17, 116, 169, 249), and any body over 65536 bytes fails with an error wrapping `ErrBodyTooLarge` (pkg/astarteapi/envelope.go:233-235) that each handler turns into a plain 400 `{"errors":{"detail":"Bad request"}}` (http.go:117-118, 170-171, 250-251). The shared `BadRequest` component (yaml:578-586) says only "malformed or missing data envelope" and no request-body schema carries a size hint, so a client is told nothing about the cap. Widen that one shared description to name the 64 KiB cap and the failure shape (a single edit covers all three operations) and say explicitly that an over-size body is 400, not 413 — do not invent a 413. One clause on `csr` (yaml:414) and `client_crt` (yaml:422) too: the cap is on the whole envelope, not on the PEM. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML). — BLOCKED: TIMEOUT after 1200s — task too big, split it
 - [x] docs-sync-pairing-realm-name-pattern [auto]: constrain `components.parameters.RealmName` in docs/api/astarte_pairing_api.yaml (yaml:359-365, a bare `type: string` described only as "The realm name.") with `pattern: '^[a-z][a-z0-9]*$'` — the realm name is CHECK-constrained to exactly that (migrations/000002_metadata.up.sql:7, read-only, never edit) and the two surfaces whose tasks already ran carry it on their twins (docs/api/astarte_housekeeping_api.yaml:246, docs/api/astrate_native_api.yaml:420 and 467). Because this one parameter is shared by all seven operations, the description must record that an off-pattern name gets three different answers, not one: the agent and device routes answer **401** (the realm is resolved by the middleware, which treats a miss as unauthenticated rather than 404 — internal/auth/middleware.go:66-72); `getHealth` answers **404** (service.Health returns store.ErrNotFound → WriteNotFound, internal/pairing/service.go:369-371, http.go:289-290); and `getVersion` answers **200** `{"data":"1.2.2"}` for any realm string whatsoever, because `GET /pairing/v1/{realm}/version` is registered with no middleware and no realm lookup (cmd/astrate/main.go:438 → observability.VersionHandler, internal/observability/compat.go:38-47) — say that in the version operation so nobody reads the shared parameter as a guarantee the handler enforces. Mirror `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go:1247) with a pairing twin only if it stays a few lines; the YAML edit is the deliverable. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
 - [x] docs-sync-pairing-deviceid-path-param [auto]: give `components.parameters.DeviceID` the constraints `RegisterRequest.hw_id` already carries in docs/api/astarte_pairing_api.yaml — the parameter (yaml:367-377) is prose-only, while `hw_id` (yaml:386-398) got `pattern: '^[A-Za-z0-9_-]{22}$'` plus `minLength`/`maxLength` from docs-sync-pairing-deviceid-base64url, and `TestPairingDeviceIDEncodingDocumented` (internal/swagger/swagger_test.go:1086) pins only the two descriptions and hw_id's bounds, never the parameter's schema; the realm-management twin was constrained later by docs-sync-rm-deviceid-param and its test (swagger_test.go:1168) does hold the parameter's pattern against `deviceid.Parse`. Add the same three schema lines to the parameter, and record in its description the asymmetry the prose omits: a malformed `deviceID` is **401** on all three device-authenticated operations — `deviceid.Parse` failing inside `authenticateDevice` burns bcrypt and returns `ErrUnauthorized` (internal/pairing/service.go:381-385) → `WriteUnauthorized` (http.go:322-323), exactly as an unknown device or a wrong secret is (service.go:395-397), so the client is told its credentials are bad rather than its ID malformed — but **404** on the agent `unregisterDevice` (service.go:231-234 → http.go:326-327). Extend TestPairingDeviceIDEncodingDocumented to compile the new pattern and hold it against `deviceid.Parse` the way the realm-management test does. Do not encode the strict trailing-bits rule in the pattern; that stays open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
+- [ ] docs-sync-ae-data-envelope-metadata [auto]: add the `metadata` property to `DataEnvelopeObject` in docs/api/astarte_appengine_api.yaml (yaml:1711-1717) — the schema requires only `data` and types it `object`, but every `format=table` data response carries a sibling `metadata` object. Code: `astarteapi.WriteDataWithMetadata` renders `{"data": v, "metadata": m}` through `metadataEnvelope` (`Metadata map[string]any json:"metadata,omitempty"`, pkg/astarteapi/envelope.go:149-158), the `Tabular` carrier sets it (internal/appengine/data.go:59-64), and both renderers populate it — `renderIndividual` (data.go:256-268, `columns` keyed on the path's last segment + `table_header`) and `renderObject` (data.go:293-321, `columns` + sorted `table_header`). Add `metadata` as an optional object property (`columns: map<string,int>`, `table_header: array of string`) described as present only when `format=table`, on all six data GET 200s that `$ref` this schema (yaml:416, 468, 685, 730, 1244, 1290). While there, correct `data`'s type/description: it is not always an object — `format=structured` individual is `[{value,timestamp}]` (data.go:275-281) and `format=disjoint_tables` is an object keyed by path (data.go:269-274, 322-331). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-ae-data-write-object-errors [auto]: document the two object-aggregated write rejections on the six data PUT/POST ops (putData/publishData, by-alias, in-group) in docs/api/astarte_appengine_api.yaml — each refs `BadRequest` for 400 (yaml:518, 570, 775, 822, 1336, 1384) and `ValueTooLarge` for 422 (yaml:528, 580, 785, 832, 1346, 1394), but three bodies are missing. (a) 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` from `writeBadObjectError` (internal/appengine/http.go:722-748, const `detailUnexpectedObjectKey` :711), reached by the `payload.ReasonOf(err) == payload.ReasonBadObject` case (http.go:684-691) — so the shared `ErrorDetail` schema (yaml:1741-1750) needs an optional `unexpected_keys` array, or a dedicated response. (b) 422 `{"errors":{"detail":"Missing required mapping key"}}` (const `detailMissingRequiredMapping`, http.go:713-720) from the `ReasonMissingRequired` case (http.go:692-702) — add an example beside the value-size one already on those ops. Both were implemented by appengine-unexpected-object-key/appengine-missing-required-422 (upstream master b6d46ad4, #2237) with no spec update. Do NOT invent other `payload.ReasonOf` statuses — that family is the separate blocked line appengine-payload-reason-status-map. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-ae-group-create-devices-required [auto]: in docs/api/astarte_appengine_api.yaml, `GroupCreate` (yaml:1652-1663) marks only `group_name` required, but `createGroup` treats `devices` as required — `groupBody.Devices *[]string` (internal/appengine/http.go:391-397) makes a MISSING key a 422 `{"errors":{"devices":["can't be blank"]}}` and an EMPTY array a 422 `{"errors":{"devices":["should have at least 1 item(s)"]}}` (http.go:405-418), and only the MISSING-key case is unreachable from the schema as written. The shared `ValidationErrors` example (yaml:1841-1844) that POST /groups refs (yaml:948) shows `group_name: can't be blank` alone. Add `devices` to `GroupCreate`'s `required` list (the 201 echo always carries it, http.go:424) and add the two device messages to the POST 422 story. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-ae-notfound-component [auto]: the `NotFound` response component in docs/api/astarte_appengine_api.yaml (yaml:1795-1803, example `errors.detail: Not Found`) is `$ref`d by all twelve group-route 404s (yaml:978, 1026, 1070, 1104, 1152, 1179, 1209, 1250, 1296, 1342, 1390, 1429), yet `astarteapi.WriteNotFound`/`DetailNotFound` ("Not Found", pkg/astarteapi/envelope.go:42-43, 205-207) is emitted nowhere in internal/appengine (rg finds only `WriteDeviceNotFound`, http.go:683). A missing group answers 404 `{"errors":{"detail":"Group not found"}}` (`ErrGroupNotFound`, internal/appengine/service.go:41-42, 508-516; writeError http.go:667-668) and a missing member device answers `"Device not found"` (writeError store.ErrNotFound → WriteDeviceNotFound, http.go:682-683). Repoint/rename the component so the group routes show `Group not found` (add a `GroupNotFound` component, keep `DeviceNotFound` on the routes that can 404 for a missing member) and state which routes can emit which. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,10 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **docs-sync run, 2026-10-09, surface: appengine — `docs/site/appengine-api.md` mislabels the `downsample_to` datastream query parameter as a duration.** The datastream-query snippet writes it as `&downsample_to=<bucket_duration>` (line 60) and the note below says "`downsample_to` maps onto Timescale `time_bucket()`" (line 65), but the code reads it as a point COUNT, not a time span: `parseQueryOpts` parses an integer and requires it `> 2`, storing it in `opts.DownsamplePoints` (internal/appengine/http.go:588-598), and the bucket interval is computed *from* that count downstream (`bucketFor(last.Sub(first), opts.DownsamplePoints)` then `s.Downsample(ctx, q, bucket)`, internal/appengine/data.go:187-196) — which is what the OpenAPI already says ("Downsample to this many data points (must be > 2)", docs/api/astarte_appengine_api.yaml:1509-1514). A reader following the page sends a duration like `1h` and gets a 422 (`downsample_to is invalid`, http.go:591-592). Proposed wording (your voice): `&downsample_to=<max_points>` and "`downsample_to` is the target number of points (must be greater than 2); the server chooses the Timescale `time_bucket()` interval." Page untouched.
+
 - **github-issues triage run, 2026-10-08 (evening, 20:37Z): nothing proposable for the 27th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub.** The recipe's `--limit 40` command prints all **32** open issues: still **28** `mule-alarm` **#94–#121** (no #122 — nothing new filed since this morning's run) plus the same four non-alarm, none proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), one excluded by standing instruction. **The alarm latch is working again and #121's window has ended**: `.mule/.alarmed` is gone, `.mule/.heartbeat` = `1791488167` = **2026-10-08T19:36:07Z**, and `.mule/log.md` now carries four `2026-10-08` `done` rows the morning run could not see (`d070345`, `cfa6e76`, `ed301fb`, `1cda0fb`, plus one 1200s timeout) — the very land that clears the latch (mule.sh:578→737), so #121 needs no action either. **Proposal unchanged from earlier today: close #94–#120, all twenty-seven** — each is a one-day alarm superseded by the next and nothing reads issue state; leave #121 to self-expire. Re-verified today, not inherited: **#92's parking condition still unmet** (`gh api .../tags` → newest `v1.4.0-rc.6`, `/releases` → newest stable `v1.3.5`; an `-rc.N` does not satisfy "wait for a stable v1.4.0"), and **#93's rewrite still entirely off main** — `git branch -a --contains 24ad5b8` lists only `mule/queue` / `origin/mule/queue`, `origin/main...HEAD` = 4 / **732** (4 / 711 this morning).
 
 ---
 

$ cd /root/astrate-mule && ls .mule/reports/log/ | tail -15; echo "---date---"; date -u +%Y-%m-%dT%H%M%SZ
20261008T113945Z-recipe-astarte-upstream-proposed.md
20261008T172851Z-recipe-code-review-proposed.md
20261008T173721Z-testutil-wait-since-cursor-done.md
20261008T174752Z-testutil-control-frame-test-blocked.md
20261008T175838Z-engine-introspection-producer-roundtrip-done.md
20261008T182643Z-recipe-docs-sync-proposed.md
20261008T183153Z-docs-sync-pairing-register-404-unreachable-done.md
20261008T191310Z-docs-sync-pairing-body-cap-400-blocked.md
20261008T192820Z-docs-sync-pairing-realm-name-pattern-done.md
20261008T193604Z-docs-sync-pairing-deviceid-path-param-done.md
20261008T200623Z-recipe-hygiene-timeout.md
20261008T202151Z-recipe-milestones-proposed.md
20261008T203843Z-recipe-github-issues-proposed.md
20261008T204257Z-recipe-astarte-upstream-proposed.md
20261009T105827Z-recipe-code-review-proposed.md
---date---
2026-10-09T112852Z

→ Read .mule/reports/log/20261008T182643Z-recipe-docs-sync-proposed.md

$ cd /root/astrate-mule && git rev-parse --short HEAD && hostname && date -u +%Y-%m-%dT%H:%M:%SZ
a086f26
DietPi
2026-10-09T11:29:11Z

← Write .mule/reports/log/20261009T112852Z-recipe-docs-sync-proposed.md
Wrote file successfully.

$ cd /root/astrate-mule && grep -nE 'docs-sync-ae-(data-envelope-metadata|data-write-object-errors|group-create-devices-required|notfound-component)' .mule/todo.md; echo "--- for-giulio ---"; sed -n '11,15p' .mule/for-giulio.md; echo "--- git status ---"; git status --short
299:- [ ] docs-sync-ae-data-envelope-metadata [auto]: add the `metadata` property to `DataEnvelopeObject` in docs/api/astarte_appengine_api.yaml (yaml:1711-1717) — the schema requires only `data` and types it `object`, but every `format=table` data response carries a sibling `metadata` object. Code: `astarteapi.WriteDataWithMetadata` renders `{"data": v, "metadata": m}` through `metadataEnvelope` (`Metadata map[string]any json:"metadata,omitempty"`, pkg/astarteapi/envelope.go:149-158), the `Tabular` carrier sets it (internal/appengine/data.go:59-64), and both renderers populate it — `renderIndividual` (data.go:256-268, `columns` keyed on the path's last segment + `table_header`) and `renderObject` (data.go:293-321, `columns` + sorted `table_header`). Add `metadata` as an optional object property (`columns: map<string,int>`, `table_header: array of string`) described as present only when `format=table`, on all six data GET 200s that `$ref` this schema (yaml:416, 468, 685, 730, 1244, 1290). While there, correct `data`'s type/description: it is not always an object — `format=structured` individual is `[{value,timestamp}]` (data.go:275-281) and `format=disjoint_tables` is an object keyed by path (data.go:269-274, 322-331). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
300:- [ ] docs-sync-ae-data-write-object-errors [auto]: document the two object-aggregated write rejections on the six data PUT/POST ops (putData/publishData, by-alias, in-group) in docs/api/astarte_appengine_api.yaml — each refs `BadRequest` for 400 (yaml:518, 570, 775, 822, 1336, 1384) and `ValueTooLarge` for 422 (yaml:528, 580, 785, 832, 1346, 1394), but three bodies are missing. (a) 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` from `writeBadObjectError` (internal/appengine/http.go:722-748, const `detailUnexpectedObjectKey` :711), reached by the `payload.ReasonOf(err) == payload.ReasonBadObject` case (http.go:684-691) — so the shared `ErrorDetail` schema (yaml:1741-1750) needs an optional `unexpected_keys` array, or a dedicated response. (b) 422 `{"errors":{"detail":"Missing required mapping key"}}` (const `detailMissingRequiredMapping`, http.go:713-720) from the `ReasonMissingRequired` case (http.go:692-702) — add an example beside the value-size one already on those ops. Both were implemented by appengine-unexpected-object-key/appengine-missing-required-422 (upstream master b6d46ad4, #2237) with no spec update. Do NOT invent other `payload.ReasonOf` statuses — that family is the separate blocked line appengine-payload-reason-status-map. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
301:- [ ] docs-sync-ae-group-create-devices-required [auto]: in docs/api/astarte_appengine_api.yaml, `GroupCreate` (yaml:1652-1663) marks only `group_name` required, but `createGroup` treats `devices` as required — `groupBody.Devices *[]string` (internal/appengine/http.go:391-397) makes a MISSING key a 422 `{"errors":{"devices":["can't be blank"]}}` and an EMPTY array a 422 `{"errors":{"devices":["should have at least 1 item(s)"]}}` (http.go:405-418), and only the MISSING-key case is unreachable from the schema as written. The shared `ValidationErrors` example (yaml:1841-1844) that POST /groups refs (yaml:948) shows `group_name: can't be blank` alone. Add `devices` to `GroupCreate`'s `required` list (the 201 echo always carries it, http.go:424) and add the two device messages to the POST 422 story. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
302:- [ ] docs-sync-ae-notfound-component [auto]: the `NotFound` response component in docs/api/astarte_appengine_api.yaml (yaml:1795-1803, example `errors.detail: Not Found`) is `$ref`d by all twelve group-route 404s (yaml:978, 1026, 1070, 1104, 1152, 1179, 1209, 1250, 1296, 1342, 1390, 1429), yet `astarteapi.WriteNotFound`/`DetailNotFound` ("Not Found", pkg/astarteapi/envelope.go:42-43, 205-207) is emitted nowhere in internal/appengine (rg finds only `WriteDeviceNotFound`, http.go:683). A missing group answers 404 `{"errors":{"detail":"Group not found"}}` (`ErrGroupNotFound`, internal/appengine/service.go:41-42, 508-516; writeError http.go:667-668) and a missing member device answers `"Device not found"` (writeError store.ErrNotFound → WriteDeviceNotFound, http.go:682-683). Repoint/rename the component so the group routes show `Group not found` (add a `GroupNotFound` component, keep `DeviceNotFound` on the routes that can 404 for a missing member) and state which routes can emit which. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
--- for-giulio ---
---

- **docs-sync run, 2026-10-09, surface: appengine — `docs/site/appengine-api.md` mislabels the `downsample_to` datastream query parameter as a duration.** The datastream-query snippet writes it as `&downsample_to=<bucket_duration>` (line 60) and the note below says "`downsample_to` maps onto Timescale `time_bucket()`" (line 65), but the code reads it as a point COUNT, not a time span: `parseQueryOpts` parses an integer and requires it `> 2`, storing it in `opts.DownsamplePoints` (internal/appengine/http.go:588-598), and the bucket interval is computed *from* that count downstream (`bucketFor(last.Sub(first), opts.DownsamplePoints)` then `s.Downsample(ctx, q, bucket)`, internal/appengine/data.go:187-196) — which is what the OpenAPI already says ("Downsample to this many data points (must be > 2)", docs/api/astarte_appengine_api.yaml:1509-1514). A reader following the page sends a duration like `1h` and gets a 422 (`downsample_to is invalid`, http.go:591-592). Proposed wording (your voice): `&downsample_to=<max_points>` and "`downsample_to` is the target number of points (must be greater than 2); the server chooses the Timescale `time_bucket()` interval." Page untouched.

- **github-issues triage run, 2026-10-08 (evening, 20:37Z): nothing proposable for the 27th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub.** The recipe's `--limit 40` command prints all **32** open issues: still **28** `mule-alarm` **#94–#121** (no #122 — nothing new filed since this morning's run) plus the same four non-alarm, none proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), one excluded by standing instruction. **The alarm latch is working again and #121's window has ended**: `.mule/.alarmed` is gone, `.mule/.heartbeat` = `1791488167` = **2026-10-08T19:36:07Z**, and `.mule/log.md` now carries four `2026-10-08` `done` rows the morning run could not see (`d070345`, `cfa6e76`, `ed301fb`, `1cda0fb`, plus one 1200s timeout) — the very land that clears the latch (mule.sh:578→737), so #121 needs no action either. **Proposal unchanged from earlier today: close #94–#120, all twenty-seven** — each is a one-day alarm superseded by the next and nothing reads issue state; leave #121 to self-expire. Re-verified today, not inherited: **#92's parking condition still unmet** (`gh api .../tags` → newest `v1.4.0-rc.6`, `/releases` → newest stable `v1.3.5`; an `-rc.N` does not satisfy "wait for a stable v1.4.0"), and **#93's rewrite still entirely off main** — `git branch -a --contains 24ad5b8` lists only `mule/queue` / `origin/mule/queue`, `origin/main...HEAD` = 4 / **732** (4 / 711 this morning).
--- git status ---
 M .mule/for-giulio.md
 M .mule/todo.md
?? .mule/reports/log/20261009T112852Z-recipe-docs-sync-proposed.md

Done: ran the docs-sync recipe on the **appengine** surface as a proposal job; appended 4 task lines to `.mule/todo.md`, escalated 1 `docs/site` contradiction to `.mule/for-giulio.md`, wrote the run report.

Files:
- `.mule/todo.md:299-302` (4 new `docs/api/` proposals)
- `.mule/for-giulio.md:13` (newest-at-top escalation for `docs/site/appengine-api.md:60,65`)
- `.mule/reports/log/20261009T112852Z-recipe-docs-sync-proposed.md` (report)

Proposals:
1. `docs-sync-ae-data-envelope-metadata` — `DataEnvelopeObject` omits the `metadata` object that every `format=table` data GET returns.
2. `docs-sync-ae-data-write-object-errors` — 400 `unexpected_keys` and 422 `Missing required mapping key` bodies absent from the six data-write ops.
3. `docs-sync-ae-group-create-devices-required` — `GroupCreate.devices` is de facto required but not `required`; 422 device messages undocumented.
4. `docs-sync-ae-notfound-component` — the `NotFound` component ("Not Found") is unreachable; group routes emit "Group not found".

Route-set comparison found no path gaps; existing blocked lines (`docs-sync-ae-list-devices-422`, `docs-sync-appengine-data-set-422`, `docs-sync-ae-read-query-params`) were not re-proposed.

Verified: `grep`/`read` cross-checks against `internal/appengine/{http,data,service}.go` and `pkg/astarteapi/envelope.go`; `git status --short` shows only the three intended files. Not run: `make -C docs build` (no `docs/` file changed in a proposal run).

Unsure: which group routes can 404 for a missing member — proposal 4 leaves that to the executor; the key claim (the `NotFound` body is emitted nowhere in appengine) is pinned by `rg`.

Follow-ups: none beyond the 4 lines above. git untouched (no commit/branch/checkout/stash), no GitHub mutation.
