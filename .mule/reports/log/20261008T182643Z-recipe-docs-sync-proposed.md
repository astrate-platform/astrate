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
---
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

$ grep -n '^- \[.\] docs-sync' .mule/todo.md | sed 's/:.*\(docs-sync[^:]*\):.*/ \1/' | head -80; echo ===; wc -l .mule/todo.md
80 docs-sync-rm-policies-delete-422 [auto]
81 docs-sync-rm-triggers-422-nested-envelope [auto]
116 docs-sync-pairing-health-path [auto]
117 docs-sync-pairing-register-404 [auto]
123 docs-sync-rm-datastream-retention-endpoint [auto]
124 docs-sync-rm-interfaces-detailed-param [auto]
127 docs-sync-appengine-by-alias-endpoints [auto]
128 docs-sync-appengine-group-endpoints [auto]
129 docs-sync-appengine-get-group-device [auto]
130 docs-sync-appengine-query-params-status [auto]
132 docs-sync-appengine-group-patch-status [auto]
133 docs-sync-appengine-data-422-interface-level [auto]
136 docs-sync-appengine-data-output-params [auto]
143 docs-sync-appengine-device-status-schema [auto]
144 docs-sync-appengine-data-set-422 [auto]
145 docs-sync-appengine-downsample-min [auto]
149 docs-sync-hk-patch-endpoint [auto]
150 docs-sync-hk-retention-field [auto]
154 docs-sync-native-compat-health-503 [auto]
155 docs-sync-pairing-status-enum [auto]
156 docs-sync-pairing-version-endpoint [auto]
157 docs-sync-native-compat-version-endpoints [auto]
158 docs-sync-rm-delete-interface-status [auto]
159 docs-sync-rm-mapping-required-encrypted [auto]
160 docs-sync-rm-put-auth-422 [auto]
161 docs-sync-rm-version-example [auto]
167 docs-sync-ae-write-405 [auto]
168 docs-sync-ae-write-value-422 [auto]
169 docs-sync-ae-post-groups-409 [auto]
172 docs-sync-ae-post-group-devices-409 [auto]
173 docs-sync-ae-patch-merge-patch-media-type [auto]
178 docs-sync-hk-delete-gating-responses [auto]
179 docs-sync-hk-validation-example [auto]
185 docs-sync-native-socket-query-token-auth [auto]
188 docs-sync-pairing-deviceendpoints-dead-404-403 [auto]
189 docs-sync-pairing-version-404-unreachable [auto]
194 docs-sync-pairing-version-value [auto]
196 docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go
201 docs-sync-native-metrics-example-fake-series [auto]
202 docs-sync-native-socket-missing-403-500 [auto]
208 docs-sync-ae-read-query-params
216 docs-sync-rm-put-interface-409 [auto]
217 docs-sync-rm-interface-422-shapes [auto]
218 docs-sync-hk-validation-example (housekeeping http.go
219 docs-sync-rm-auth-403 [auto]
226 docs-sync-hk-async-operation-param [auto]
227 docs-sync-hk-retention-zero-is-unset [auto]
228 docs-sync-rm-async-operation-param [auto]
230 docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first
236 docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go
237 docs-sync-ae-patch-by-alias-409 [auto]
238 docs-sync-ae-add-group-device-422 [auto]
239 docs-sync-ae-delete-data-400 [auto]
240 docs-sync-ae-list-devices-422 [auto]
241 docs-sync-pairing-error-example-capitalisation [auto]
242 docs-sync-native-version-value already corrected to that (docs/api/astrate_native_api.yaml
243 docs-sync-pairing-deviceid-base64url [auto]
244 docs-sync-pairing-initial-payload-format-enum [auto]
245 docs-sync-pairing-unregister-description [auto]
254 docs-sync-rm-error-example-capitalisation [auto]
255 docs-sync-rm-legacy-alias-fields [auto]
256 docs-sync-rm-validationerror-example [auto]
257 docs-sync-pairing-deviceid-base64url did for the pairing twin. Do NOT encode the strict trailing-bits rule in the pattern
258 docs-sync-rm-update-interface-body [auto]
260 docs-sync-ae-error-example-capitalisation [auto]
261 docs-sync-hk-error-example-capitalisation [auto]
262 docs-sync-native-error-example-capitalisation [auto]
263 docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured
268 docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted
269 docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail
270 docs-sync-hk-patch-422-field-error-examples [auto]
271 docs-sync-hk-wrong-type-field-400 [auto]
272 docs-sync-hk-realm-name-response-schemas [auto]
===
294 .mule/todo.md

$ ls .mule/ .mule/recipes/ .mule/tasks/ 2>/dev/null
.mule/:
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

.mule/recipes/:
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

.mule/tasks/:
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md

$ ls -t .mule/reports | head -20; echo ===; ls -t .mule/failed | head
log
hygiene.md
===
gate.log
testutil-control-frame-test.diff
pairing-burn-bcrypt-cost.gate.log
pairing-burn-bcrypt-cost.diff
docs-native-socket-event-schema.diff
docs-native-socket-handshake-errors.diff
probe-encrypted-mapping-enforcement.diff
store-latestindividual-empty-errnotfound.diff
store-register-inhibit-preserve.gate.log
store-register-inhibit-preserve.diff

$ grep -n 'surface' .mule/log.md | tail -20; echo ===; tail -60 .mule/log.md
===
| 2026-09-26 | forward-status-error-body | done | 209s | 915c260 |
| 2026-09-26 | forward-envelope-bytes-test | done | 213s | ed869b6 |
| 2026-09-26 | webhook-static-headers-override-test | done | 178s | f9db0a6 |
| 2026-09-26 | docs-sync-ae-forbidden-403 | done | 1037s | 6a29ac7 |
| 2026-09-27 | docs-sync-ae-patch-by-alias-409 | done | 119s | d1af059 |
| 2026-09-27 | docs-sync-ae-add-group-device-422 | done | 274s | dff3be8 |
| 2026-09-27 | docs-sync-ae-delete-data-400 | done | 330s | 66bc366 |
| 2026-09-27 | docs-sync-ae-list-devices-422 | transient | 344s |  > build · big-pickle  Error: Error from provider (Console): Rate limit exceeded. Please try again later.  |
| 2026-09-27 | docs-sync-ae-list-devices-422 | transient | 248s |  > build · big-pickle  Error: Error from provider (Console): Rate limit exceeded. Please try again later.  |
| 2026-09-27 | docs-sync-ae-list-devices-422 | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-09-28 | docs-sync-pairing-error-example-capitalisation | done | 166s | 7ed78cb |
| 2026-09-28 | docs-sync-pairing-info-version-example | done | 333s | d9ec98f |
| 2026-09-28 | docs-sync-pairing-deviceid-base64url | done | 480s | 32693cf |
| 2026-09-28 | docs-sync-pairing-initial-payload-format-enum | done | 95s | 1f6a6c2 |
| 2026-09-28 | docs-sync-pairing-unregister-description | done | 261s | 8ffb39e |
| 2026-10-01 | fdo-rc6-scope-delta-for-giulio | done | 479s | 41af3c4 |
| 2026-10-02 | compat-note-v14-rc6 | done | 398s | 0215f97 |
| 2026-10-02 | container-stop-deadline | done | 410s | f9af73b |
| 2026-10-02 | container-timeout-bounds | blocked | 1201s | TIMEOUT after 1201s — task too big, split it |
| 2026-10-02 | container-parseconfig-rules-test | done | 274s | 32c85af |
| 2026-10-02 | container-response-cap-test | blocked | 271s | lint failed: internal/flow/blocks/container/block_test.go:500:2: redefines-builtin-id: redefinition of the built-in function cap (rev |
| 2026-10-02 | docs-sync-rm-error-example-capitalisation | done | 751s | 1d1e6f4 |
| 2026-10-02 | docs-sync-rm-legacy-alias-fields | blocked | 329s | wrote nothing |
| 2026-10-02 | docs-sync-rm-validationerror-example | done | 732s | 8a09a1a |
| 2026-10-02 | docs-sync-rm-deviceid-param | done | 303s | 2d1a1c1 |
| 2026-10-02 | docs-sync-rm-update-interface-body | blocked | 1201s | TIMEOUT after 1201s — task too big, split it |
| 2026-10-03 | docs-sync-ae-error-example-capitalisation | done | 298s | f3bcc6d |
| 2026-10-03 | docs-sync-hk-error-example-capitalisation | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-10-03 | docs-sync-native-error-example-capitalisation | done | 284s | 7cc237f |
| 2026-10-03 | docs-sync-hk-error-detail-examples-split | done | 114s | 23006d4 |
| 2026-10-03 | device-empty-cache-received-trigger | blocked | 202s | wrote nothing |
| 2026-10-03 | appengine-unexpected-object-key | done | 956s | f7d60af |
| 2026-10-03 | appengine-payload-reason-status-map | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-10-03 | appengine-missing-required-422 | done | 584s | f84bbb5 |
| 2026-10-04 | docs-sync-hk-realm-name-pattern | done | 525s | e728217 |
| 2026-10-04 | docs-sync-hk-errordetail-schema-example | done | 184s | 5044143 |
| 2026-10-04 | docs-sync-hk-patch-422-field-error-examples | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-10-04 | docs-sync-hk-wrong-type-field-400 | blocked | 152s | wrote nothing |
| 2026-10-04 | docs-sync-hk-realm-name-response-schemas | done | 178s | def0af6 |
| 2026-10-04 | store-register-inhibit-preserve | blocked | 245s | gates failed |
| 2026-10-04 | store-aliasvalues-self-exclusion-test | done | 717s | a7267ec |
| 2026-10-04 | store-latestindividual-empty-errnotfound | blocked | 1201s | TIMEOUT after 1201s — task too big, split it |
| 2026-10-05 | store-validatepipelinegraph-error-branches | done | 316s | 73e7387 |
| 2026-10-05 | docs-native-ae-realm-version-endpoint | done | 344s | 7331e6c |
| 2026-10-05 | docs-native-hk-version-description | done | 360s | cce1376 |
| 2026-10-05 | docs-native-metrics-content-negotiation | done | 346s | 3962921 |
| 2026-10-05 | docs-native-phoenix-newevent-payload | done | 267s | a16b5fe |
| 2026-10-05 | docs-native-realm-name-pattern | done | 174s | 2ddbf7f |
| 2026-10-06 | probe-encrypted-mapping-enforcement | blocked | 21s | wrote nothing |
| 2026-10-06 | compat-note-v135 | done | 190s | 3a8ea6d |
| 2026-10-06 | docs-native-socket-handshake-errors | blocked | 91s | wrote nothing |
| 2026-10-06 | docs-native-socket-sse-exact-accept | done | 410s | 1322214 |
| 2026-10-06 | docs-native-socket-security-scheme | done | 299s | 87c18cf |
| 2026-10-06 | docs-native-socket-event-schema | blocked | 662s | opencode exited 1 |
| 2026-10-07 | pairing-burn-bcrypt-cost | blocked | 359s | lint failed: internal/pairing/service.go:95:7: G101: Potential hardcoded credentials (gosec) |
| 2026-10-07 | pairing-onregistered-clock | done | 143s | 4e71533 |
| 2026-10-07 | pairing-unregister-bad-id-test | done | 225s | 322242a |
| 2026-10-08 | testutil-wait-since-cursor | done | 319s | 3654fda |
| 2026-10-08 | testutil-control-frame-test | blocked | 102s | wrote nothing |
| 2026-10-08 | engine-introspection-producer-roundtrip | done | 234s | d070345 |

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*pairing*'
(no output)

$ rg -n '^\s{2}/' docs/api/astarte_pairing_api.yaml
27:  /pairing/v1/{realm}/health:
70:  /pairing/v1/{realm}/version:
97:  /pairing/v1/{realm}/agent/devices:
149:  /pairing/v1/{realm}/agent/devices/{deviceID}:
183:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
235:  /pairing/v1/{realm}/devices/{deviceID}:
269:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:

$ ls internal/pairing/; echo ===; rg -n 'Handle|mux\.|Route|Get\(|Post\(|Delete\(' internal/pairing/*.go | head -60
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
===
internal/pairing/http.go:46:	requireAgent func(http.Handler) http.Handler
internal/pairing/http.go:74:// Mount registers the pairing routes on mux. Paths are wire-frozen
internal/pairing/http.go:81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/pairing/http.go:82:	mux.Handle("POST /pairing/v1/{realm}/agent/devices",
internal/pairing/http.go:83:		a.requireAgent(http.HandlerFunc(a.handleRegister)))
internal/pairing/http.go:84:	mux.Handle("DELETE /pairing/v1/{realm}/agent/devices/{deviceID}",
internal/pairing/http.go:85:		a.requireAgent(http.HandlerFunc(a.handleUnregister)))
internal/pairing/http.go:86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
internal/pairing/http.go:88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
internal/pairing/http.go:90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
internal/pairing/http.go:344:	header := r.Header.Get("Authorization")
internal/pairing/ratelimit.go:59:	b, ok := l.buckets.Get(key)

→ Read docs/api/astarte_pairing_api.yaml [limit=300]
→ Read internal/pairing/http.go
→ Read docs/api/astarte_pairing_api.yaml [offset=300]

$ rg -n 'func \(s \*Service\)|return .*Err|ErrNotFound|ErrUnauthorized|ErrInhibited|ErrAlreadyRegistered|ErrInvalidHWID|ErrInvalidCSR|ErrInvalidPayloadFormat|ErrRegistrationLimitReached' internal/pairing/service.go
44:	// ErrInvalidHWID reports a hw_id that is not a 22-character unpadded
46:	ErrInvalidHWID = errors.New("pairing: invalid hw_id")
47:	// ErrInvalidPayloadFormat reports an initial_payload_format outside
49:	ErrInvalidPayloadFormat = errors.New("pairing: invalid initial_payload_format")
50:	// ErrAlreadyRegistered reports re-registration of a device that has
53:	ErrAlreadyRegistered = errors.New("pairing: device already registered")
54:	// ErrRegistrationLimitReached reports the realm's
57:	ErrRegistrationLimitReached = errors.New("pairing: device registration limit reached")
58:	// ErrUnauthorized is the uniform device-authentication failure: unknown
61:	ErrUnauthorized = errors.New("pairing: unauthorized")
62:	// ErrInhibited reports a device blocked by credentials_inhibited (403).
63:	ErrInhibited = errors.New("pairing: credentials request inhibited")
64:	// ErrInvalidCSR reports an unusable certificate signing request (422).
65:	ErrInvalidCSR = errors.New("pairing: invalid CSR")
163:		return "", nil, fmt.Errorf("pairing: sealing CA key: %w", err)
172:// secret; afterwards it fails with ErrAlreadyRegistered. initialFormat is
174:func (s *Service) Register(ctx context.Context, realmName, hwID, initialFormat string) (string, error) {
177:		return "", fmt.Errorf("%w: %v", ErrInvalidHWID, err)
180:		return "", fmt.Errorf("%w: %q", ErrInvalidPayloadFormat, initialFormat)
196:			return "", ErrRegistrationLimitReached
206:		return "", fmt.Errorf("pairing: hashing credentials secret: %w", err)
211:			return "", fmt.Errorf("%w: %s", ErrAlreadyRegistered, hwID)
229:// credential trail). store.ErrNotFound is returned for unknown devices.
230:func (s *Service) Unregister(ctx context.Context, realmName, deviceIDStr string) error {
233:		return fmt.Errorf("%w: device %q", store.ErrNotFound, deviceIDStr)
243:// device by credentials secret (uniform ErrUnauthorized on any mismatch),
247:func (s *Service) Credentials(ctx context.Context, realmName, deviceIDStr, secret, csrPEM string, ip netip.Addr) (string, error) {
253:		return "", ErrInhibited
262:		if errors.Is(err, ca.ErrInvalidCSR) {
263:			return "", fmt.Errorf("%w: %v", ErrInvalidCSR, err)
286:func (s *Service) Info(ctx context.Context, realmName, deviceIDStr, secret string) (*Info, error) {
324:func (s *Service) VerifyCredentials(ctx context.Context, realmName, deviceIDStr, secret, clientCrtPEM string) (*VerifyResult, error) {
330:		return nil, ErrInhibited
369:func (s *Service) Health(ctx context.Context, realm string) error {
371:		return err // store.ErrNotFound → 404
378:// device, wrong secret) return the same ErrUnauthorized after a comparable
380:func (s *Service) authenticateDevice(ctx context.Context, realmName, deviceIDStr, secret string) (*store.Realm, *store.Device, error) {
384:		return nil, nil, ErrUnauthorized
387:	if errors.Is(err, store.ErrNotFound) {
389:		return nil, nil, ErrUnauthorized
395:	if errors.Is(err, store.ErrNotFound) {
397:		return nil, nil, ErrUnauthorized
405:		return nil, nil, ErrUnauthorized
408:		return nil, nil, ErrUnauthorized
414:func (s *Service) loadCA(realm *store.Realm) (*ca.CA, error) {
417:		return nil, fmt.Errorf("pairing: opening sealed CA key of realm %q: %w", realm.Name, err)
421:		return nil, fmt.Errorf("pairing: loading CA of realm %q: %w", realm.Name, err)
437:		return "", fmt.Errorf("pairing: gathering secret randomness: %w", err)

→ Read internal/pairing/service.go [offset=170, limit=200]

$ grep -n 'docs-sync-pairing-version-value\|docs-sync-pairing-version-404-unreachable\|docs-sync-pairing-health-path\|docs-sync-pairing-register-404\|pairing-realm-name\|docs-sync-pairing-status-enum\|docs-sync-pairing-info-version-example' .mule/todo.md
116:- [x] docs-sync-pairing-health-path [auto]: add the undocumented `GET /pairing/v1/{realm}/health` route to docs/api/astarte_pairing_api.yaml — it exists in code (internal/pairing/http.go:81) and returns 200 `{"data":{"status":"ok"}}` for an existing realm, 404 for an unknown realm, 503 when unhealthy, unauthenticated and unrate-limited. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
117:- [x] docs-sync-pairing-register-404 [auto]: document the missing `404` (DeviceNotFound) on `POST /pairing/v1/{realm}/agent/devices` in docs/api/astarte_pairing_api.yaml — the handler returns it when `GetRealmByName` fails for an unknown realm (internal/pairing/service.go:183-186 → http.go:326-327), but the spec lists only 201/400/401/403/422/429/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
155:- [x] docs-sync-pairing-status-enum [auto]: fix the `PairingInfo.status` enum in docs/api/astarte_pairing_api.yaml from `[confirmed, pending, denied, expired]` (yaml:412) to `[confirmed, pending, inhibited]` — the handler never emits `denied` or `expired`: `service.Info` maps only `DeviceStatusInhibited`→"inhibited", `DeviceStatusConfirmed`→"confirmed", and everything else→"pending" (internal/pairing/service.go:296-299, statuses defined in store/devices.go:19-27; upstream-parity per the comment at service.go:282-285). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
189:- [x] docs-sync-pairing-version-404-unreachable [auto]: remove the `404` "Unknown realm" response from `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml (yaml:96-104) — the route is served by `observability.VersionHandler` (cmd/astrate/main.go:398, internal/observability/compat.go:38-47), which answers static `200 {"data":"<version>"}` with no realm lookup, so a 404 can never be emitted there. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
194:- [x] docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
242:- [x] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pairing_api.yaml, which shows `version: "0.1.0-test"` (yaml:252) — a value no build can produce. The field is filled from `pairing.Config.Version` (internal/pairing/service.go:302), wired from `main.version` (cmd/astrate/main.go:151), whose default is `"0.1.0-dev"` (main.go:53) and which is otherwise set only via `-ldflags "-X main.version=..."`. Use `"0.1.0-dev"`, matching the four compat-version examples docs-sync-native-version-value already corrected to that (docs/api/astrate_native_api.yaml:203, 223, 243, 263). Do NOT document the service's own fallback `DefaultVersion = "0.1.0-astrate"` (service.go:86-88): it is unreachable from `main`, which always passes a non-empty value. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
280:- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ grep -rn 'pattern' docs/api/astarte_pairing_api.yaml docs/api/astarte_realm_management_api.yaml | head -20; echo ===; rg -n 'RealmName:' -A8 docs/api/astarte_realm_management_api.yaml | head -30
docs/api/astarte_pairing_api.yaml:390:          pattern: '^[A-Za-z0-9_-]{22}$'
docs/api/astarte_realm_management_api.yaml:945:        pattern: '^[A-Za-z0-9_-]{22}$'
docs/api/astarte_realm_management_api.yaml:992:          pattern: '^([a-zA-Z][a-zA-Z0-9]*\.([a-zA-Z0-9][a-zA-Z0-9-]*\.)*)?[a-zA-Z][a-zA-Z0-9]*$'
docs/api/astarte_realm_management_api.yaml:1158:          pattern: '^[a-zA-Z0-9_.~-]{1,128}$'
===
896:    RealmName:
897-      name: realm
898-      in: path
899-      required: true
900-      description: The realm name.
901-      schema:
902-        type: string
903-
904-    InterfaceName:

→ Read internal/auth/middleware.go

$ sed -n '243p;244p;245p;116,120p' .mule/todo.md; echo ===; sed -n '241p' .mule/todo.md
- [x] docs-sync-pairing-health-path [auto]: add the undocumented `GET /pairing/v1/{realm}/health` route to docs/api/astarte_pairing_api.yaml — it exists in code (internal/pairing/http.go:81) and returns 200 `{"data":{"status":"ok"}}` for an existing realm, 404 for an unknown realm, 503 when unhealthy, unauthenticated and unrate-limited. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-register-404 [auto]: document the missing `404` (DeviceNotFound) on `POST /pairing/v1/{realm}/agent/devices` in docs/api/astarte_pairing_api.yaml — the handler returns it when `GetRealmByName` fails for an unknown realm (internal/pairing/service.go:183-186 → http.go:326-327), but the spec lists only 201/400/401/403/422/429/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] swagger-httptest-coverage [auto]: add a container-free httptest suite for internal/swagger, which currently has no test file — cover two rules: (1) `Mount` wires `GET /swagger` to a 302→/swagger/index.html, serves the embedded UI under `/swagger/`, and serves the OpenAPI YAML specs under `/api/`; (2) `Specs` returns exactly the `.yaml` filenames from `docs.APIYAML` with no path prefix or dirs. httptest only, no Docker.
- [!] probe-property-resend-encoding [auto]: upstream v1.3.0 fixed outbound server-property values sent to a device on connect/emptyCache (commit 522ccf4f — a raw stored binaryblob is now wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding, so it ships as a BSON binary, not a string). Investigate only: does Astrate's resendServerProperties/rehydrate path (internal/engine/control.go:144, pkg/payload) emit binaryblob → BSON subtype-0 binary and datetime → UTC datetime exactly like upstream 1.3? Report, do not patch. — BLOCKED: wrote nothing
- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
- [x] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` path parameter (docs/api/astarte_pairing_api.yaml:360) and `RegisterRequest.hw_id` (yaml:373) — but `deviceid.Parse` accepts only the 22-character **unpadded base64url** form, in strict mode, and rejects the standard alphabet's `+` and `/` and any `=` padding with an error wrapping `ErrInvalid` (pkg/deviceid/deviceid.go:26, 37, 39-60; upstream parity cited at 36). Measured: with the same `encoding` the parser uses, `dT6hS2W9TT6LEnP25ks+lg` and `.../lg` fail at byte 19 and the padded 23-character form fails at byte 21, so a client generated from the current spec can send a well-formed standard-base64 ID and get 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}` (internal/pairing/http.go:307-309). State the exact encoding and the 22-character length in both descriptions and add a `pattern`/`minLength`/`maxLength` to `hw_id`. The existing `hw_id: dT6hS2W9TT6LEnP25ks_lg` example (yaml:123) is valid and stays. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in docs/api/astarte_pairing_api.yaml (yaml:374-376), an unconstrained string described only as "Optional Astrate extension specifying the device's initial payload format" — so a generated client can send any string, while `Service.Register` accepts only `"bson"` or `"json"` and rejects everything else with 422 `{"errors":{"initial_payload_format":["is invalid"]}}` (internal/pairing/service.go:179-181 -> internal/pairing/http.go:310-312; the allowed set is stated at service.go:173). Add `enum: [bson, json]` plus a note that absent/empty means unset, and add the three missing `ValidationError` examples while you are in the component: `initial_payload_format: ["is invalid"]` (http.go:310-312), `csr: ["can't be blank"]` (http.go:173-177) and `client_crt: ["can't be blank"]` (http.go:253-257) — the component currently shows only the four hw_id/error_name rows (yaml:594-618) and the `csr`/`client_crt` blank-field cases are what a device actually hits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` description in docs/api/astarte_pairing_api.yaml (yaml:153), "Removes a device from the realm. Returns 204 on success.", which reads as a device-and-data deletion. The row and every byte of its stored data survive: `store.UnregisterDevice` only clears the credential trail (`credentials_secret_hash = ''`, `first_credentials_request = NULL`, `cert_serial = NULL`, `cert_aki = NULL`) and flips `status` back to `'registered'` (internal/store/devices.go:93-110), and the service docstring says so outright ("its data is kept", internal/pairing/service.go:227-229). State the real semantics: the device becomes registrable again, its interfaces, datastream data and group memberships are retained, and a second DELETE is a 204 no-op (the row is still there, so RowsAffected is 1, devices.go:106-108) — not the documented 404. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
===
- [x] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:556) and `InternalServerError` shows `detail: Internal Server Error` (yaml:638), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37, golden fixture pkg/astarteapi/testdata/error_bad_request.json) and `"Internal server error"` (envelope.go:58, testdata/error_internal.json), both measured on upstream 1.2.0 (envelope.go:33-36, 55-57). The other five detail examples on this file (`Not Found`, `Device not found`, `Unauthorized`, `Forbidden`, `Too Many Requests`) are already correct — change only those two. The identical wrong pair is in the other four specs (appengine 1732/1792, housekeeping 362/385/477, realm_management 1265/1348, native 466); leave those for a run of their surface. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).

$ rg -n 'func DecodeData' -A 25 pkg/astarteapi/*.go
pkg/astarteapi/envelope.go:228:func DecodeData(r io.Reader, maxBytes int64, dst any) error {
pkg/astarteapi/envelope.go-229-	body, err := io.ReadAll(io.LimitReader(r, maxBytes+1))
pkg/astarteapi/envelope.go-230-	if err != nil {
pkg/astarteapi/envelope.go-231-		return fmt.Errorf("astarteapi: reading request body: %w", err)
pkg/astarteapi/envelope.go-232-	}
pkg/astarteapi/envelope.go-233-	if int64(len(body)) > maxBytes {
pkg/astarteapi/envelope.go-234-		return fmt.Errorf("astarteapi: %w: exceeds %d bytes", ErrBodyTooLarge, maxBytes)
pkg/astarteapi/envelope.go-235-	}
pkg/astarteapi/envelope.go-236-	var env struct {
pkg/astarteapi/envelope.go-237-		Data json.RawMessage `json:"data"`
pkg/astarteapi/envelope.go-238-	}
pkg/astarteapi/envelope.go-239-	if err := json.Unmarshal(body, &env); err != nil {
pkg/astarteapi/envelope.go-240-		return fmt.Errorf("astarteapi: decoding request envelope: %w", err)
pkg/astarteapi/envelope.go-241-	}
pkg/astarteapi/envelope.go-242-	if len(env.Data) == 0 || bytes.Equal(env.Data, []byte("null")) {
pkg/astarteapi/envelope.go-243-		return fmt.Errorf("astarteapi: %w", ErrMissingData)
pkg/astarteapi/envelope.go-244-	}
pkg/astarteapi/envelope.go-245-	if err := json.Unmarshal(env.Data, dst); err != nil {
pkg/astarteapi/envelope.go-246-		return fmt.Errorf("astarteapi: decoding request data: %w", err)
pkg/astarteapi/envelope.go-247-	}
pkg/astarteapi/envelope.go-248-	return nil
pkg/astarteapi/envelope.go-249-}

$ rg -n 'version|Version' cmd/astrate/main.go | sed -n '1,60p'
51:// version is the reported build version (override with
52:// -ldflags "-X main.version=vX.Y.Z").
53:var version = "0.1.0-dev"
63:	showVersion := flag.Bool("version", false, "print the version and exit")
67:	if *showVersion {
68:		fmt.Println("astrate", version)
151:		Version:           version,
227:		log.Info("astrate listening", "http", cfg.HTTP.Addr, "mqtt", b.TLSAddr(), "version", version)
390:	// Upstream-parity per-service version endpoints (issue #77): every
391:	// service answers an unauthenticated GET /{service}/version with
392:	// {"data": version}; AppEngine's realm-scoped variant requires auth,
394:	// (test/conformance/upstream/verify-versions.json). Realm Management's
397:		observability.MountVersionCompat(mux, svc, version)
399:	mountRealmVersion(mux, mw)
429:// mountRealmVersion registers the realm-scoped version endpoints observed on
430:// upstream 1.2.0 (test/conformance/upstream/verify-versions.json): AppEngine's
433:// (COMPATIBILITY.md deviation 10) — not Astrate's build version. RM's
435:func mountRealmVersion(mux *http.ServeMux, mw *auth.Middleware) {
436:	mux.Handle("GET /appengine/v1/{realm}/version",
437:		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(realm.APICompatVersion)))
438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))

$ rg -n 'func VersionHandler' -A 20 internal/observability/compat.go; echo ===; rg -n 'APICompatVersion' -r . --glob '*.go' | head
38:func VersionHandler(version string) http.HandlerFunc {
39-	body, err := json.Marshal(map[string]string{"data": version})
40-	if err != nil {
41-		body = []byte(`{"data":""}`)
42-	}
43-	return func(w http.ResponseWriter, _ *http.Request) {
44-		w.Header().Set("Content-Type", "application/json")
45-		_, _ = w.Write(body)
46-	}
47-}
48-
49-// MountVersionCompat registers GET /{service}/version for one service.
50-func MountVersionCompat(mux *http.ServeMux, service, version string) {
51-	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
52-}
===
cmd/astrate/main.go:437:		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(realm..)))
cmd/astrate/main.go:438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm..))
cmd/astrate/version_test.go:25:	want := `{"data":"` + realm.. + `"}`
internal/realm/http.go:119:// feature UI on it — see .).
internal/realm/http.go:121:	_ = astarteapi.WriteData(w, http.StatusOK, .)
internal/realm/dashboard_compat_test.go:35:		if v != . {
internal/realm/dashboard_compat_test.go:36:			t.Errorf("version = %q, want %q", v, .)
internal/realm/service.go:580:// . is the upstream Realm Management API level Astrate
internal/realm/service.go:588:const . = "1.2.2"

→ Read internal/pairing/service.go [offset=376, limit=64]

$ rg -n 'ErrNotFound|WriteDeviceNotFound|WriteNotFound|StatusNotFound' internal/pairing/*.go pkg/astarteapi/envelope.go | head -30
pkg/astarteapi/envelope.go:205:// WriteNotFound writes the canonical generic 404 envelope.
pkg/astarteapi/envelope.go:206:func WriteNotFound(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:207:	return WriteError(w, http.StatusNotFound, DetailNotFound)
pkg/astarteapi/envelope.go:210:// WriteDeviceNotFound writes the canonical 404 "Device not found" envelope.
pkg/astarteapi/envelope.go:211:func WriteDeviceNotFound(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:212:	return WriteError(w, http.StatusNotFound, DetailDeviceNotFound)
internal/pairing/http.go:289:	case errors.Is(err, store.ErrNotFound):
internal/pairing/http.go:290:		_ = astarteapi.WriteNotFound(w)
internal/pairing/http.go:326:	case errors.Is(err, store.ErrNotFound):
internal/pairing/http.go:327:		_ = astarteapi.WriteDeviceNotFound(w)
internal/pairing/service.go:229:// credential trail). store.ErrNotFound is returned for unknown devices.
internal/pairing/service.go:233:		return fmt.Errorf("%w: device %q", store.ErrNotFound, deviceIDStr)
internal/pairing/service.go:371:		return err // store.ErrNotFound → 404
internal/pairing/service.go:387:	if errors.Is(err, store.ErrNotFound) {
internal/pairing/service.go:395:	if errors.Is(err, store.ErrNotFound) {
internal/pairing/http_test.go:528:		if status != http.StatusNotFound {
internal/pairing/http_test.go:598:		if status != http.StatusNotFound {
internal/pairing/service_test.go:50:		return nil, fmt.Errorf("%w: realm %q", store.ErrNotFound, name)
internal/pairing/service_test.go:78:		return fmt.Errorf("%w: device %s", store.ErrNotFound, id)
internal/pairing/service_test.go:90:		return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, id)
internal/pairing/service_test.go:99:		return fmt.Errorf("%w: device %s", store.ErrNotFound, id)
internal/pairing/service_test.go:364:	// Unregister turns every lookup failure into store.ErrNotFound: both
internal/pairing/service_test.go:375:		if err := svc.Unregister(ctx, "test", tc.id); !errors.Is(err, store.ErrNotFound) {
internal/pairing/service_test.go:376:			t.Errorf("Unregister(%s = %q): got %v, want store.ErrNotFound", tc.name, tc.id, err)

$ rg -n 'StatusNotFound' -B 25 internal/pairing/http_test.go | sed -n '1,70p'
503-		status, body := f.request(t, "DELETE", base+"/agent/devices/"+hwID, agentToken, nil)
504-		if status != http.StatusNoContent {
505-			t.Fatalf("status: got %d, want 204 (body %s)", status, body)
506-		}
507-		if len(body) != 0 {
508-			t.Errorf("204 body must be empty, got %q", body)
509-		}
510-
511-		// Old secret is dead.
512-		status, _ = f.request(t, "POST",
513-			base+"/devices/"+hwID+"/protocols/astarte_mqtt_v1/credentials",
514-			secret, map[string]string{"csr": deviceCSR(t)})
515-		if status != http.StatusUnauthorized {
516-			t.Fatalf("old secret after unregister: got %d, want 401", status)
517-		}
518-
519-		// The device is registrable again (its row and data survived).
520-		status, body = f.request(t, "POST", base+"/agent/devices", agentToken,
521-			map[string]string{"hw_id": hwID})
522-		if status != http.StatusCreated {
523-			t.Fatalf("re-register after unregister: got %d (body %s)", status, body)
524-		}
525-
526-		// Unknown device → upstream 404 shape.
527-		status, body = f.request(t, "DELETE", base+"/agent/devices/"+mustRandomDeviceID(t), agentToken, nil)
528:		if status != http.StatusNotFound {
--
573-		t.Fatalf("status: got %d, want 429 (body %s)", status, body)
574-	}
575-	testutil.Golden(t, "http/envelope_429.json", body)
576-}
577-
578-// TestPairingRealmHealth covers the realm-scoped health probe (#71,
579-// upstream 1.3+): unauthenticated by design — none of these requests carries
580-// an Authorization header.
581-func TestPairingRealmHealth(t *testing.T) {
582-	f := newHTTPFixture(t, Config{}, APIConfig{})
583-	path := "/pairing/v1/" + f.realmName + "/health"
584-
585-	t.Run("KnownRealm200", func(t *testing.T) {
586-		status, got := f.request(t, "GET", path, "", nil)
587-		if status != http.StatusOK {
588-			t.Fatalf("status: got %d, want 200 (body %s)", status, got)
589-		}
590-		if want := `{"data":{"status":"ok"}}`; string(got) != want {
591-			t.Errorf("body: got %s, want %s", got, want)
592-		}
593-	})
594-
595-	t.Run("UnknownRealm404", func(t *testing.T) {
596-		status, _ := f.request(t, "GET",
597-			"/pairing/v1/nope"+strconv.FormatInt(time.Now().UnixNano(), 36)+"/health", "", nil)
598:		if status != http.StatusNotFound {

$ rg -n 'func newHTTPFixture' -A 60 internal/pairing/http_test.go
48:func newHTTPFixture(t *testing.T, cfg Config, apiCfg APIConfig) *httpFixture {
49-	t.Helper()
50-	ctx := context.Background()
51-
52-	pool := testutil.StartTimescale(t)
53-	st, err := store.New(ctx, pool.Config().ConnString())
54-	if err != nil {
55-		t.Fatalf("store.New: %v", err)
56-	}
57-	t.Cleanup(st.Close)
58-
59-	key := make([]byte, store.MasterKeySize)
60-	if _, err := rand.Read(key); err != nil {
61-		t.Fatal(err)
62-	}
63-	sealer, err := store.NewKeySealer(key)
64-	if err != nil {
65-		t.Fatal(err)
66-	}
67-
68-	jwtKey, err := rsa.GenerateKey(rand.Reader, 2048)
69-	if err != nil {
70-		t.Fatal(err)
71-	}
72-	pubDER, err := x509.MarshalPKIXPublicKey(&jwtKey.PublicKey)
73-	if err != nil {
74-		t.Fatal(err)
75-	}
76-	pubPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
77-
78-	realmName := "p" + strconv.FormatInt(time.Now().UnixNano(), 36)
79-	caPEM, sealedKey, err := ProvisionCA(realmName, 0, sealer)
80-	if err != nil {
81-		t.Fatalf("ProvisionCA: %v", err)
82-	}
83-	realm, err := st.CreateRealm(ctx, store.NewRealm{
84-		Name:               realmName,
85-		JWTPublicKeysPEM:   []string{pubPEM},
86-		CACertificatePEM:   caPEM,
87-		CAPrivateKeySealed: sealedKey,
88-	})
89-	if err != nil {
90-		t.Fatalf("CreateRealm: %v", err)
91-	}
92-
93-	if cfg.BrokerURL == "" {
94-		cfg.BrokerURL = "mqtts://broker.test.example:8883"
95-	}
96-	if cfg.Version == "" {
97-		cfg.Version = "0.1.0-test"
98-	}
99-	svc := New(st, sealer, cfg)
100-	api := NewAPI(svc, auth.NewMiddleware(st), apiCfg)
101-	mux := http.NewServeMux()
102-	api.Mount(mux)
103-	server := httptest.NewServer(mux)
104-	t.Cleanup(server.Close)
105-
106-	return &httpFixture{
107-		st: st, sealer: sealer, realm: realm, realmName: realmName,
108-		server: server, jwtKey: jwtKey, svc: svc,

$ rg -n 'Header\(\)|Content-Type|func Write' pkg/astarteapi/envelope.go | head -40
21:// ContentType is the Content-Type header value set on every envelope,
112:		w.Header().Set("Content-Type", ContentType)
117:	w.Header().Set("Content-Type", ContentType)
126:func WriteData(w http.ResponseWriter, status int, v any) error {
145:func WriteDataWithLinks(w http.ResponseWriter, status int, v any, links Links) error {
156:func WriteDataWithMetadata(w http.ResponseWriter, status int, v any, m map[string]any) error {
163:func WriteError(w http.ResponseWriter, status int, detail string) error {
172:func WriteFieldErrors(w http.ResponseWriter, status int, fields map[string][]string) error {
186:func WriteRawErrors(w http.ResponseWriter, status int, errs any) error {
191:func WriteBadRequest(w http.ResponseWriter) error {
196:func WriteUnauthorized(w http.ResponseWriter) error {
201:func WriteForbidden(w http.ResponseWriter) error {
206:func WriteNotFound(w http.ResponseWriter) error {
211:func WriteDeviceNotFound(w http.ResponseWriter) error {
216:func WriteInternalServerError(w http.ResponseWriter) error {

$ ls docs/site/; echo ===; rg -n 'pairing' docs/site/*.md -l
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
===
docs/site/troubleshooting.md
docs/site/quickstart.md
docs/site/operations.md
docs/site/json-payload-profile.md
docs/site/deployment.md
docs/site/compatibility.md
docs/site/pairing-and-security.md
docs/site/migration-from-astarte.md
docs/site/configuration-reference.md
docs/site/payload-formats.md
docs/site/pokemon-agent.md
docs/site/index.md
docs/site/contributing.md
docs/site/architecture.md
docs/site/DESIGN.md
docs/site/ROADMAP.md

$ wc -l docs/site/pairing-and-security.md; rg -n '40[134]|429|rate|inhibit|unregister|422|401|403' docs/site/pairing-and-security.md | head -40
99 docs/site/pairing-and-security.md
3:Astrate reproduces Astarte's three credential planes exactly, replacing the CFSSL sidecar with an embedded per-realm CA.
43:- After first credentials request, re-registration conflicts (422).
74:Agent                Astrate(pairing)          Device                 Astrate(broker)
93:3. **inhibited** -- blocked from new credentials and connections (via `PATCH /appengine/v1/<realm>/devices/<id>` with `credentials_inhibited: true`).

$ rg -n "a-z\]\[a-z0-9" docs/api/*.yaml; echo ===; sed -n '1,12p' migrations/000002_metadata.up.sql; echo ===; rg -n 'DetailNotFound|DetailDeviceNotFound' pkg/astarteapi/envelope.go
docs/api/astrate_native_api.yaml:420:            pattern: '^[a-z][a-z0-9]*$'
docs/api/astrate_native_api.yaml:467:        pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_housekeeping_api.yaml:246:        pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_housekeeping_api.yaml:280:            pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_housekeeping_api.yaml:289:          pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_housekeeping_api.yaml:311:          pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_housekeeping_api.yaml:484:        a non-empty name off the `^[a-z][a-z0-9]*$` pattern is refused by the
===
-- 000002: relational metadata schema, transcribed verbatim from
-- docs/DESIGN.md §2.2 (docs/ROADMAP.md §3.1 file 2.2).

-- Realms (Housekeeping domain)
CREATE TABLE realms (
    id               smallint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name             text NOT NULL UNIQUE CHECK (name ~ '^[a-z][a-z0-9]*$'),
    jwt_public_keys  jsonb NOT NULL DEFAULT '[]',   -- array of PEM strings (RSA/EC)
    ca_certificate   text NOT NULL,                  -- realm CA cert, PEM
    ca_private_key   bytea NOT NULL,                 -- encrypted at rest (AES-256-GCM,
                                                     -- key from config/KMS env var)
    device_registration_limit integer,
===
42:	// DetailNotFound is the canonical generic 404 detail.
43:	DetailNotFound = "Not Found"
44:	// DetailDeviceNotFound is the canonical 404 detail for unknown devices
46:	DetailDeviceNotFound = "Device not found"
49:	// DetailNotFound, which a handler emits for a route that matched and found
207:	return WriteError(w, http.StatusNotFound, DetailNotFound)
212:	return WriteError(w, http.StatusNotFound, DetailDeviceNotFound)

$ rg -n 'realm' docs/api/astarte_pairing_api.yaml | sed -n '1,5p'; rg -n 'pattern' docs/api/astrate_native_api.yaml | head; echo ===; rg -n 'ASTRATE_[A-Z_]+' -o internal/ | sort -u | head -50
9:    Agent endpoints (`/agent/`) require a realm JWT with the `a_pa` claim.
27:  /pairing/v1/{realm}/health:
32:        Returns the health status of a realm's backing store. Unauthenticated
36:        Returns 200 when the realm exists and the database answers, 404 for an
37:        unknown realm, and 503 when the database is unhealthy.
414:            so an off-pattern name is not reported as malformed: it is answered
420:            pattern: '^[a-z][a-z0-9]*$'
461:        off-pattern name is not reported as malformed: it is answered `401`
467:        pattern: '^[a-z][a-z0-9]*$'
===
internal/config/config.example.toml:10:ASTRATE_ENGINE_SHARDS
internal/config/config.example.toml:10:ASTRATE_LOG_LEVEL
internal/config/config.example.toml:10:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/config/config.example.toml:87:ASTRATE_MASTER_KEY
internal/config/config.example.toml:87:ASTRATE_MASTER_KEY_FILE
internal/config/config.example.toml:9:ASTRATE_DATABASE_DSN
internal/config/config.example.toml:9:ASTRATE_HTTP_ADDR
internal/config/config.go:135:ASTRATE_MASTER_KEY
internal/config/config.go:135:ASTRATE_MASTER_KEY_FILE
internal/config/config.go:230:ASTRATE_HTTP_ADDR
internal/config/config.go:231:ASTRATE_HTTP_TLS_CERT_FILE
internal/config/config.go:232:ASTRATE_HTTP_TLS_KEY_FILE
internal/config/config.go:233:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
internal/config/config.go:242:ASTRATE_MQTT_ADDR
internal/config/config.go:243:ASTRATE_MQTT_TLS_CERT_FILE
internal/config/config.go:244:ASTRATE_MQTT_TLS_KEY_FILE
internal/config/config.go:245:ASTRATE_MQTT_SESSION_STORE_PATH
internal/config/config.go:246:ASTRATE_MQTT_ADVERTISED_URL
internal/config/config.go:247:ASTRATE_DATABASE_DSN
internal/config/config.go:248:ASTRATE_SECURITY_MASTER_KEY_FILE
internal/config/config.go:249:ASTRATE_REALM_NAME
internal/config/config.go:250:ASTRATE_REALM_JWT_PUBLIC_KEY
internal/config/config.go:251:ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
internal/config/config.go:252:ASTRATE_LOG_LEVEL
internal/config/config.go:253:ASTRATE_LOG_FORMAT
internal/config/config.go:255:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/config/config.go:263:ASTRATE_ENGINE_SHARDS
internal/config/config.go:265:ASTRATE_ENGINE_SHARDS
internal/config/config.go:279:ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
internal/config/config.go:297:ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
internal/config/config.go:313:ASTRATE_DATABASE_DSN
internal/config/config_test.go:140:ASTRATE_DATABASE_DSN
internal/config/config_test.go:141:ASTRATE_ENGINE_SHARDS
internal/config/config_test.go:142:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
internal/config/config_test.go:392:ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
internal/config/config_test.go:459:ASTRATE_ENGINE_SHARDS
internal/config/config_test.go:490:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/config/config_test.go:533:ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
internal/engine/engine_integration_test.go:85:ASTRATE_TEST_DSN
internal/flow/blocks/container/block.go:45:ASTRATE_FLOW_CONFIG
internal/flow/blocks/container/block_test.go:125:ASTRATE_FLOW_CONFIG
internal/flow/blocks/container/docker.go:105:ASTRATE_FLOW_CONFIG
internal/flow/blocks/container/docker.go:19:ASTRATE_FLOW_CONFIG
internal/flow/blocks/container/docker.go:213:ASTRATE_FLOW_CONFIG
internal/flow/blocks/info.go:116:ASTRATE_FLOW_CONFIG
internal/flow/blocks/schema.go:19:ASTRATE_FLOW_CONFIG
internal/store/crypto.go:24:ASTRATE_MASTER_KEY
internal/store/crypto.go:27:ASTRATE_MASTER_KEY_FILE
internal/store/store_test.go:216:ASTRATE_TEST_DSN
internal/store/store_test.go:26:ASTRATE_TEST_DSN

$ rg -n '64|KiB|kilobyte|body|too large' docs/api/astarte_pairing_api.yaml | head -20; echo ===; rg -n 'maxBodyBytes|BodyTooLarge' -r '' internal/ pkg/ | head -20
103:        secret. The request body is wrapped in a data envelope.
173:          description: Device unregistered successfully. No body returned.
373:        base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and `_`,
393:            base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and
397:            `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}`
398:            body.
405:            `{"errors":{"initial_payload_format":["is invalid"]}}` body.
432:          description: One-time device credentials secret (base64-encoded).
648:                    - "is not a valid base64 encoded 128 bits id"
===
pkg/astarteapi/envelope.go:65:// Err is wrapped by DecodeData when the request body exceeds the
pkg/astarteapi/envelope.go:67:var Err = errors.New("request body too large")
pkg/astarteapi/envelope.go:223:// It fails with an error wrapping Err when the body exceeds
pkg/astarteapi/envelope.go:234:		return fmt.Errorf("astarteapi: %w: exceeds %d bytes", Err, maxBytes)
pkg/astarteapi/envelope_test.go:242:		{"oversized", `{"data": "` + strings.Repeat("a", 64) + `"}`, 32, astarteapi.Err},
pkg/astarteapi/envelope_test.go:243:		{"oversized by one", `{"data": 11}`, 11, astarteapi.Err},
internal/pairing/http.go:15://  caps pairing request bodies (docs/DESIGN.md §4.5 input
internal/pairing/http.go:17:const  int64 = 64 << 10
internal/pairing/http.go:116:	if err := astarteapi.DecodeData(r.Body, , &req); err != nil {
internal/pairing/http.go:169:	if err := astarteapi.DecodeData(r.Body, , &req); err != nil {
internal/pairing/http.go:249:	if err := astarteapi.DecodeData(r.Body, , &req); err != nil {
internal/flowapi/http.go:14://  caps Flow API request bodies (pipeline DAGs are small JSON).
internal/flowapi/http.go:15:const  int64 = 1 << 20
internal/flowapi/http.go:77:	if err := astarteapi.DecodeData(r.Body, , &body); err != nil {
internal/flowapi/http.go:102:	if err := astarteapi.DecodeData(r.Body, , &def); err != nil {
internal/flowapi/http.go:133:	if err := astarteapi.DecodeData(r.Body, , &body); err != nil {
internal/flowapi/http.go:187:	if err := astarteapi.DecodeData(r.Body, , &config); err != nil {
internal/flowapi/http.go:220:	if err := astarteapi.DecodeData(r.Body, , &body); err != nil {
internal/flowapi/http.go:234:	if err := astarteapi.DecodeData(r.Body, , &body); err != nil {
internal/realm/http.go:17://  caps Realm Management request bodies (docs/DESIGN.md §4.5):

$ rg -n '403|401|wrong secret|credentials secret' docs/site/*.md | head -30
docs/site/troubleshooting.md:78:2. **Losing it does not brick existing devices.** Devices keep their credentials secret and will re-pair automatically at their next credential rotation.
docs/site/troubleshooting.md:82:## Dashboard 401 / 403
docs/site/troubleshooting.md:84:**Symptom:** Astarte Dashboard shows "401 Unauthorized" or "403 Forbidden".
docs/site/troubleshooting.md:87:1. **401:** The JWT token is missing, expired, or signed with an unknown key. Generate a fresh token:
docs/site/troubleshooting.md:91:2. **403:** The token is valid but lacks the required `a_aea` claim for AppEngine access.
docs/site/operations.md:26:next credential rotation, since their credentials secret still works).
docs/site/operations.md:133:2. Devices keep their credentials secret, so they re-pair automatically: on the
docs/site/json-payload-profile.md:140:   credentials secret.
docs/site/compatibility.md:53:### 5. Uniform 401 vs 403
docs/site/compatibility.md:55:Astrate returns `401` for every failure to establish identity (wrong secret, unknown device, bad token) and `403` only for authorization refusal after identity is established. Stricter/safer -- eliminates enumeration oracles.
docs/site/compatibility.md:102:Upstream's dynamic pool registers each first-seen id through Pairing and then spawns a real MQTT device, keeping the credentials secret in a local store -- losing that store permanently bricks every id whose certificate was issued, with no recovery path. Astrate keeps the observable contract (key grammar, first-seen registration through the pairing door, rows queryable like any device-owned datastream) but lands values through the engine ingest path and keeps secrets server-side, so that failure mode does not exist.
docs/site/DESIGN.md:56:| **Pairing API** + CFSSL | Device registration, credentials secret issuance, X.509 CSR signing, broker info | `internal/pairing` | CA is embedded (`crypto/x509`); per-realm CA key in DB (encrypted) or on-disk PEM. Same REST surface (§4). |
docs/site/DESIGN.md:675:  their credentials secret still works).
docs/site/DESIGN.md:765:- Secrets handling: credentials secrets bcrypt-hashed; CA keys AES-GCM-encrypted; JWT public
docs/site/configuration-reference.md:76:| `bcrypt_cost` | int | `10` | — | No | bcrypt cost for hashing credentials secrets. |
docs/site/configuration-reference.md:98:    Losing the master key means re-issuing realm CAs. Devices re-pair automatically at their next credential rotation since their credentials secret still works. Keep a separate backup of the master key.
docs/site/ROADMAP.md:149:| 1.3 | `pkg/astarteapi/envelope.go` | `WriteData(w, status, v)` → `{"data": v}`; `WriteError(w, status, detail)` → `{"errors":{"detail":"..."}}`; canonical upstream error constructors (404 `"Device not found"`, 401, 403, 422 shapes); request body `{"data": ...}` unwrapper with size cap | 120 |
docs/site/ROADMAP.md:248:| 3.4 | `internal/auth/middleware.go` | `Require(api Claim, keys KeySource) func(http.Handler) http.Handler` — extracts realm from path, 401/403 via `astarteapi` envelopes; housekeeping variant (instance-level keys) | 110 |
docs/site/ROADMAP.md:256:- Middleware: 401 (no/bad token) vs 403 (valid token, claim mismatch) with golden envelopes.
docs/site/ROADMAP.md:276:- **T2 HTTP flow tests (golden bodies):** Flow A 201 + 44-char secret; second register pre-credentials returns a *different* secret; post-credentials register → 422 upstream-shaped; bad `hw_id` (21 chars, padded, non-url) → 422. Flow B: 201 `client_crt` that (a) parses, (b) has CN `<realm>/<device>`, (c) chains to the realm CA returned by Flow C — verified with `crypto/x509` *and* an `openssl verify` exec smoke; wrong secret → 401 with the same body/timing class as unknown device (uniform error assertion); rate limit → 429. Flow C: info golden body (§4.4 shape exactly); verify on fresh cert → `valid:true` + `until`; on expired (issue with 1 s TTL, sleep) → `EXPIRED`; on cert from a different CA → `INVALID`.
docs/site/ROADMAP.md:277:- Inhibited device: credentials → 403 parity.
docs/site/ROADMAP.md:379:- **T2 golden suites per service:** envelopes + status codes byte-compared (404 device, 409/422 interface conflicts, 401/403 split). Interfaces: minor-bump additive accepted, mapping-mutation rejected (CheckMinorUpgrade), major coexistence, delete-while-introspected rejected. Housekeeping: created realm immediately serves pairing (cross-domain test with M4). AppEngine: datastream query boundaries (`since` inclusive vs `since_after` exclusive), `limit`, descending default ordering parity, `downsample_to` bucket-count correctness; longinteger > 2^53 round-trips as string through publish→query; property tree shape golden.
docs/site/ROADMAP.md:427:| 9.8 | `test/security/*_test.go` | Zip-bomb on `producer/properties`; pairing brute-force rate-limit; authz sweep (every mux route 401s without token — route-table-driven); TLS config assertions (min version, no client-auth bypass on :8883); oversize topic/payload/introspection bounds (§4.5) |
docs/site/migration-from-astarte.md:108:3. A valid credentials secret (obtained at registration).
docs/site/migration-from-astarte.md:119:   - If devices store their credentials secret, they can re-CSR on next boot.
docs/site/migration-from-astarte.md:135:2. **Devices will re-pair** with Astarte on their next connection attempt (they lost their Astrate-issued certs, but their original credentials secret still works with Astarte if you re-imported device registrations).
docs/site/appengine-api.md:126:Status codes match upstream: 401 (no/bad token), 403 (claim mismatch), 404 (not found), 409/422 (conflict/validation).

$ rg -n 'MaxBytesReader' -r '' . --glob '*.go' | head; echo ===; rg -n '^func Test' internal/swagger/swagger_test.go
===
26:func TestMount(t *testing.T) {
85:func TestMountSubPanicsOnBrokenFS(t *testing.T) {
98:func TestSpecs(t *testing.T) {
131:func TestRealmManagement403(t *testing.T) {
142:func TestAppEngine403(t *testing.T) {
206:func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
249:func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
324:func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {
384:func TestAppEngineAddGroupDevice422Documented(t *testing.T) {
428:func TestAppEngineDataDelete400Documented(t *testing.T) {
473:func TestPairingErrorDetailExamples(t *testing.T) {
524:func TestRealmManagementErrorDetailExamples(t *testing.T) {
576:func TestAppEngineErrorDetailExamples(t *testing.T) {
631:func TestNativeErrorDetailExamples(t *testing.T) {
682:func TestNativeVersionDescriptions(t *testing.T) {
731:func TestNativeSocketSecurityDocumented(t *testing.T) {
802:func TestNativeMetricsContentNegotiationDocumented(t *testing.T) {
938:func TestHousekeepingErrorDetailExamples(t *testing.T) {
1009:func TestRealmManagement422ValidationDetails(t *testing.T) {
1086:func TestPairingDeviceIDEncodingDocumented(t *testing.T) {
1168:func TestRealmManagementDeviceIDEncodingDocumented(t *testing.T) {
1247:func TestHousekeepingRealmNamePatternDocumented(t *testing.T) {
1353:func TestPairingUnregisterDeviceSemanticsDocumented(t *testing.T) {

$ grep -n 'docs-sync-pairing-register-404-unreachable\|docs-sync-pairing-body-cap\|docs-sync-pairing-realm-name\|docs-sync-pairing-deviceid-path' .mule/todo.md; echo "none-found-ok"; grep -c '^- \[ \]' .mule/todo.md; grep -n '^- \[ \]' .mule/todo.md | head -20
none-found-ok
15
86:- [ ] broker-external-bus-intake [legion]: implement the second Intake implementation named in the TODO at internal/broker/intake.go:56-62 — an external-bus-backed (e.g. NATS JetStream) Intake for multi-instance deployment / restart survival, reproducing the embedded broker's per-device ordering, deferred-ack backpressure and QoS 0 drop semantics; an implementation of this interface needs containerised integration tests, so it runs on the Legion Go.
88:- [ ] hygiene-govulncheck [legion]: on the Legion Go, run govulncheck ./... on a fresh ~/astrate clone and file a task line for any REACHABLE advisory it reports (name the CVE and the call path, per the hygiene recipe); this box cannot build it — go install golang.org/x/vuln/cmd/govulncheck@latest is OOM-killed here (3.7GB, 4 cores), so the recipe's highest-priority check never runs unless done there.
89:- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.go:560) with no stop signal, so a block-death racing process shutdown (cmd/astrate/main.go:482-488 drains the Manager but never cancels the goroutine) can rebuild the flow via `mgr.StartFlow` on a background pump (internal/flow/flow.go:173-175) after `Manager.Shutdown` and flip the durable status back to running. Thread a stop channel/context through the Service, checked in the loop's sleep, cancelled on shutdown; verification needs a `[legion]` integration or timing-based test. [auto]
101:- [ ] race-check-store: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/store/... ./internal/housekeeping/... ./migrations/...`. Report any failure to .mule/for-giulio.md with the full race report. Split out of the former single `race-check` line, which timed out running the whole tree at once. [legion] [readonly]
102:- [ ] race-check-engine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/engine/... ./internal/broker/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
103:- [ ] race-check-flow: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/flow/... ./internal/realm/... ./internal/pairing/... ./internal/auth/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
104:- [ ] race-check-appengine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/appengine/... ./internal/observability/... ./internal/httpx/... ./internal/config/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
105:- [ ] race-check-pkg: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./pkg/... ./cmd/... ./internal/testutil/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
200:- [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a locally-run `elixir -e` with `Base.url_decode64!`) with the device ID `"AAAAAAAAAAAAAAAAAAAAAB"` — nonzero unused bits in the 22nd base64url char — and report whether upstream accepts it: `pkg/deviceid/deviceid.go:34-37` uses `base64.RawURLEncoding.Strict()` and deviceid_test.go:176 asserts rejection while the comment claims parity with `Elixir Base.url_decode64!(padding: false)`, which decodes and discards those bits rather than erroring. Astrate may therefore be stricter than upstream on the same wire form (a device registered with such an id on upstream would be rejected here). If upstream rejects it identically, the strict encoding is verified parity and the task is done; if upstream accepts it, escalate the relax-vs-strict call to `.mule/for-giulio.md` with the measurement. Probe first, no code change either way.
203:- [ ] store-devices-inhibit-re-register [legion] [auto]: `RegisterDevice` (internal/store/devices.go:75-91) silently clears the inhibit flag — its `ON CONFLICT ... DO UPDATE SET status = 'registered'` fires for any device with `first_credentials_request IS NULL`, including one the admin inhibited via `SetDeviceInhibited` (devices.go:251-268), which sets `status='inhibited'` on unconfirmed devices too; §5.3 says an inhibited device blocks new credentials and connections, so the re-registration re-opens it. Preserve `'inhibited'` in the SET (e.g. `status = CASE WHEN devices.status = 'inhibited' THEN 'inhibited' ELSE 'registered' END`) and add a Lifecycle case in internal/store/devices_test.go: inhibit an unconfirmed device, re-register, assert status stays inhibited and the secret still rotates. Verify the assert against upstream's register-not-touching-inhibit on the Legion while the integration suite runs.
205:- [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (internal/store/devices.go:130-131, "if several devices share an alias the lowest device ID wins", `ORDER BY id LIMIT 1`) — no test drives it today; add a case in internal/store/devices_test.go that gives two devices the same alias and asserts the lookup resolves the lower ID. Needs the integration DB.
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
