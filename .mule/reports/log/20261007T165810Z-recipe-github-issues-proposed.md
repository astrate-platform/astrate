slug: recipe-github-issues
verdict: proposed
at:  dde64e0
ran: 2026-10-07T16:58:10Z on DietPi in 601s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md

$ gh issue list --state open --limit 40 \
  --json number,title,labels,updatedAt,comments \
  --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
#120 mule: nothing has landed in 14h [mule-alarm ] []c
#119 mule: nothing has landed in 16h [mule-alarm ] []c
#118 mule: nothing has landed in 14h [mule-alarm ] []c
#117 mule: nothing has landed in 15h [mule-alarm ] []c
#116 mule: nothing has landed in 14h [mule-alarm ] []c
#115 mule: nothing has landed in 14h [mule-alarm ] []c
#114 mule: nothing has landed in 8h [mule-alarm ] []c
#113 mule: nothing has landed in 14h [mule-alarm ] []c
#112 mule: nothing has landed in 14h [mule-alarm ] []c
#111 mule: nothing has landed in 16h [mule-alarm ] []c
#110 mule: nothing has landed in 14h [mule-alarm ] []c
#109 mule: nothing has landed in 14h [mule-alarm ] []c
#108 mule: nothing has landed in 15h [mule-alarm ] []c
#107 mule: nothing has landed in 14h [mule-alarm ] []c
#106 mule: nothing has landed in 14h [mule-alarm ] []c
#105 mule: nothing has landed in 14h [mule-alarm ] []c
#104 mule: nothing has landed in 14h [mule-alarm ] []c
#103 mule: nothing has landed in 16h [mule-alarm ] []c
#102 mule: nothing has landed in 14h [mule-alarm ] []c
#101 mule: nothing has landed in 14h [mule-alarm ] []c
#100 mule: nothing has landed in 15h [mule-alarm ] []c
#99 mule: nothing has landed in 14h [mule-alarm ] []c
#98 mule: nothing has landed in 16h [mule-alarm ] []c
#97 mule: nothing has landed in 15h [mule-alarm ] []c
#96 mule: nothing has landed in 14h [mule-alarm ] []c
#95 mule: nothing has landed in 15h [mule-alarm ] []c
#94 mule: nothing has landed in 15h [mule-alarm ] []c
#93 In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged [mule-review upstream-parity ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:The mule pushed `8c61268` to `mule/queue` for this.

Unreviewed: the gates passed — it compiles, the tests pass, a new test was shown to fail without the change, and no frozen file was touched. That is not the same as the change being worth having. Nothing merges to `main` until someone reads it:

    bash tools/mule.sh review

Left open on purpose. Whether this actually resolves the issue is your call, not mine. createdAt:2026-09-04T19:54:58Z id:IC_kwDORmfsJs8AAAABSo2DbA includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/93#issuecomment-5545755500 viewerDidAuthor:true]]c
#92 keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds [enhancement upstream-parity upstream-experimental ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:Parked by Giulio on 2026-09-04: **wait for a stable \`v1.4.0\`.** The spec is complete but ships only in \`v1.4.0-rc.5\`, and building against an rc risks rebuilding it.

So this issue is not waiting on anyone to remember it. It has a row in [`.mule/waiting-on.md`](https://github.com/astrate-platform/astrate/blob/main/.mule/waiting-on.md), and the weekly `mule-upstream-watch` job now reads that file as its first step every run: when a stable `v1.4.0` tag exists upstream, it escalates to `.mule/for-giulio.md` naming this issue. An `-rc.N` does not satisfy the row.

That mechanism was added because it was missing here: #51 was parked on the same condition, the condition fired on 2026-08-31, and nothing connected the new upstream tags to the parked issue — it took a hand triage four days later to catch it. createdAt:2026-09-04T19:52:21Z id:IC_kwDORmfsJs8AAAABSo0f8A includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/92#issuecomment-5545730032 viewerDidAuthor:true]]c
#78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate) [enhancement milestone-4.0 upstream-parity ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:Reframed 2026-08-23 by Giulio's decision: no longer parked — **strategic feature** (zero-touch factory onboarding is a commercial-viability requirement for an IoT platform), milestone-4.0 candidate.

## Verified facts (2026-08-23, upstream master)

- Upstream integrates FDO **inside the Pairing service**: `fdo_onboarding_controller.ex`, session plugs (`setup_fdo`, `fdo_session`), CBOR codec + voucher queries in `astarte_data_access`, and two dedicated libs `astarte_fdo` / `astarte_fdo_core` (SECO Mind, 2025; bugfixes through 2026-05). Not an optional sidecar.
- Device side is the official SDK `astarte-device-fdo-rust` (crate on docs.rs, pushed 2026-08-21) — actively maintained.
- Upstream does NOT implement manufacturing / rendezvous services either; those come from the FIDO Alliance reference stack (`fdo-rs/fido-device-onboard-rs`, used by Red Hat/Fedora IoT).

## Scope decision

Astrate implements **only the last mile**, mirroring upstream's cut:
- owner-side protocol surface in our pairing service (TO1 redirect + TO2 onboarding, CBOR wire format, ownership-voucher storage, owner key management)
- reuse the existing open-source FDO ecosystem for manufacturing/rendezvous — never reimplement it.

**Acceptance:** a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end and lands as a provisioned realm device.

**Docs are a first-class deliverable:** upstream ships almost no public FDO documentation (code-level only); ours must exceed that — operator guide covering the full chain (manufacturing → rendezvous → owner) with Astrate-specific setup.

## Next step

Investigation phase before any implementation: read upstream's TO2 handling end-to-end (`libs/astarte_fdo*", pairing controllers/plugs, `fdo/queries`), inventory the exact endpoints, credential/voucher schema and key material required, and measure what our pairing service lacks. Output feeds the v4.0 scope decision in `.mule/milestones.md`. createdAt:2026-08-23T16:16:08Z id:IC_kwDORmfsJs8AAAABQRcugw includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/78#issuecomment-5387005571 viewerDidAuthor:true]]c
#1 Provide an Open Source IoT Platform unironically, unironically [wontfix ] []c

$ grep -n "issue-" /root/astrate-mule/.mule/todo.md; echo "---TAIL---"; tail -30 /root/astrate-mule/.mule/todo.md; echo "---LINES---"; wc -l /root/astrate-mule/.mule/todo.md
---TAIL---
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
---LINES---
287 /root/astrate-mule/.mule/todo.md

$ ls -la /root/astrate-mule/.mule/; echo "=== for-giulio ==="; cat /root/astrate-mule/.mule/for-giulio.md 2>/dev/null || echo "(absent)"
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_117451544001s31J1LUvbu68ol

---

- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

---

- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

---

- **deviceid wire-parity concern, from the 2026-09-18 code review: Astrate is likely stricter than upstream on the wire form's unused bits.** `pkg/deviceid/deviceid.go:37` uses `base64.RawURLEncoding.Strict()`, so `Parse` rejects `"AAAAAAAAAAAAAAAAAAAAAB"` (unused low 4 bits of the 22nd char set; rejection pinned at deviceid_test.go:176), while the comment at deviceid.go:34-36 claims parity with `Elixir Base.url_decode64!(padding: false)` — which decodes the 128 payload bits and discards the unused ones rather than erroring. If so, a device registered with a non-canonical-trailing-bits id on upstream Astarte is accepted there but rejected by Astrate's `Parse` — a wire-compat break in a package whose purpose is the wire form. Not verified on the Pi (no Docker/Elixir); a `[legion]` probe (`deviceid-trailing-bits-upstream-probe`) is queued to measure it. Once confirmed, the call is yours: relax `Parse` to upstream's decode-and-ignore, or keep strict canonicality as the package's contract.

---

- **COMPATIBILITY.md wording update for upstream v1.3.4 (newest stable, 2026-09-18; v1.4.0 is still rc.5-only) — folds v1.3.3 and v1.3.4 into the open v1.3.2 proposal below; both are maintenance-only and wire-inert, so that proposal stands unchanged.** v1.3.3 (2026-08-07) was an empty-body patch and v1.3.4 (2026-09-18, `astarte_data_updater_plant` only — more frequent fullsweep GC on AMQPDataConsumer processes plus an RPC-availability fix; release body via `gh api repos/astarte-platform/astarte/releases` `v1.3.4`) introduce **no wire/API surface change**, so the v1.3.2 wording proposal stands complete: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"). The only delta a v1.3.4-aware doc carries is the version reference: the proposed §Infrastructure-differences sentence's "until the milestone that adopts v1.3.2 as the target" reads "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target". Raw: [v1.3.4](https://github.com/astarte-platform/astarte/releases/tag/v1.3.4).

---

- **github-issues triage run, 2026-09-18: still nothing proposable — the mule-alarm pile is now 13 straight days (#94–#106).** Sixteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-17 run: **#106**, today's alarm (created 09-18 10:54Z, "nothing has landed in 14h"); **#105** has expired into the pile (superseded by #106). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-17 one: close #94–#105** — each was a one-day event, superseded by next day's alarm, never actionable; leave #106 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.

---

- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) shows a new **stable v1.3.4** (2026-09-18) — maintenance only (data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change), so no new gap for Astrate — and still **no stable v1.4.0**: newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

---

- **docs-sync realm-management, 2026-09-17: `docs/site/realm-management-api.md` documents an "Update trigger" endpoint Astrate does not serve.** Lines 81-86 show `PUT /realmmanagement/v1/<realm>/triggers/<name>` with a `{ "data": <updated trigger JSON> }` body, but no such route is registered — internal/realm/http.go:44-47 wires only GET/POST `/triggers` and GET/DELETE `/triggers/{name}` (an update op is deliberately absent; not even a stub or TODO exists). Site prose is yours — drop the "Update trigger" section or confirm the route is planned.

---

- **github-issues triage run, 2026-09-17: still nothing proposable — the mule-alarm pile is now 12 straight days (#94–#105).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-16 run: **#105**, today's alarm (created 09-17 10:51Z, "nothing has landed in 14h"); **#104** has expired into the pile (superseded by #105). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-16 one: close #94–#104** — each was a one-day event, superseded by next day's alarm, never actionable; leave #105 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.

- **github-issues triage run, 2026-09-16: still nothing proposable — the mule-alarm pile is now 11 straight days (#94–#104).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), plus the daily alarm. New since the 2026-09-15 run: **#104**, today's alarm (created 09-16 11:12Z, "nothing has landed in 14h"); **#103** has expired into the pile (superseded by #104). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-15 one: close #94–#103** — each was a one-day event, superseded by next day's alarm, never actionable; leave #104 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.

- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-13 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

- **github-issues triage run, 2026-09-15: still nothing proposable — the mule-alarm pile is now 10 straight days (#94–#103).** Fourteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-13 run: **#103**, the alarm created 2026-09-14 11:27Z ("nothing has landed in 16h"); **#102** has expired into the pile (superseded by #103). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-13 one: close #94–#102** — each was a one-day event, superseded by next day's alarm, never actionable; leave #103 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.

- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-12 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

- **github-issues triage run, 2026-09-13: still nothing proposable — the mule-alarm pile is now 9 straight days (#94–#102).** Thirteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-12 run: **#102**, today's alarm (created 11:09Z, "nothing has landed in 14h"); **#101** has expired into the pile (superseded by #102). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-12 one: close #94–#101** — each was a one-day event, superseded by next day's alarm, never actionable; leave #102 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.

- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.

- **github-issues triage run, 2026-09-12: still nothing proposable — the mule-alarm pile is now 8 straight days (#94–#101).** Twelve open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-11 run: **#101**, today's alarm (created 10:58Z, "nothing has landed in 14h"); **#100** has expired into the pile (superseded by #101). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-11 one: close #94–#100** — each was a one-day event, superseded by next day's alarm, never actionable; leave #101 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.

- **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pile is now 7 straight days and worth a look.** Twelve open issues. Still no machine-checkable fix candidates, so no task lines: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-10 run: **#100**, today's alarm (created 10:55Z, "nothing has landed in 15h"). #99–#94 are the same `mule-alarm` (one per day 2026-09-05→09-10, ~11:00Z each, zero comments) — telemetry, not code issues, so never proposable. But seven consecutive alarms is past the "low-activity window" wording earlier runs used: the mule has landed no commit since ~2026-09-04/05, and `main`'s queue copy shows the top line as `- [ ]` while the landed work below is all `[!] BLOCKED` (`wrote nothing` / `tests failed`) — though `mule/queue` is the authoritative copy. **Proposal: close #94–#99** — each was a one-day event, superseded by next day's alarm, never actionable; leave #100 (live today) to self-expire. If the idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.

- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

- **Extension point: external-bus Intake (NATS JetStream), 2026-09-10.** `internal/broker/intake.go:56-62` carries a `TODO(extension point)` for a second `Intake` implementation backed by a durable external bus (NATS JetStream) for multi-instance deployment or restart survival. The frozen design decision is in `docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4`; the consumer must reproduce per-device ordering, deferred-ack backpressure, and QoS 0 drop semantics. Not tracked by any open issue. Your call: file a milestone issue or park until the roadmap calls for it.

- **Extension point: timescaledb_toolkit lttb downsampling, 2026-09-10.** `internal/store/store.go:139-143` carries a `TODO(extension point)` for switching the `Downsample` method from `time_bucket+avg` to the toolkit's `lttb()` when `timescaledb_toolkit` is present — the probe at `store.go:144-150` already records availability in `s.hasToolkit`, but `datastreams.go` always uses the default path. Design reference: `docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5`. Not tracked by any open issue. Your call: file a milestone issue or park.

---

- **docs-sync pairing, 2026-09-10: docs/site names the pre-credentials wire status `registered`; the API returns `pending`.** `docs/site/pairing-and-security.md:56` ("flips status `registered -> confirmed`") and `:91` ("**registered** -- device registered, awaiting first credentials request"), plus `docs/site/data-modeling.md:94`, all give `registered` as a device status value. But `service.Info` emits `pending` for any device that is neither confirmed nor inhibited (internal/pairing/service.go:296-299 — upstream-parity per the comment at service.go:282-285; the DB value is `registered`, store/devices.go:19-27). Wire values are `pending`/`confirmed`/`inhibited`. Site prose is yours — reword to `pending`, or confirm the site intentionally describes the DB value.

---

- **Flow code review, 2026-09-10: `DeletePipeline` semantics need your call.**
  `DeletePipeline` (internal/flowapi/service.go:199-206) removes a pipeline
  silently even when durable flows still reference it; running instances keep
  the old graph, but their next reload fails ("pipeline not found") and a
  block death after that starts the auto-restart loop retrying forever against
  the missing pipeline. `UpdatePipeline` already reports `referencing_flows`
  (service.go:162-182, issue #44) — should DELETE also refuse (409) while
  flows reference it, or return the referencing names and allow? Design choice,
  not a mule task.

---

- **github-issues triage run, 2026-09-10: nothing proposable, nothing stale.** Nine open
  issues: #99–#94 are mule-alarms (idle queue), not code issues. #93 has a pushed commit
  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated. #1
  untouched per standing instruction. No task lines proposed, no for-giulio close suggestions.

- **Milestone recipe run, 2026-09-09: v3.0 unchanged since 2026-09-06 — still waiting on
  upstream gates.** `milestone-3.0` label remains empty. The three open `upstream-parity`
  issues are the same: #92 keyAgreement (upstream decision, gated on stable v1.4.0), #93
  aclhook comment (pushed commit `8c61268`, mule-review), #78 FDO (milestone-4.0). Both
  UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait
  on upstream v1.4 final (currently rc.5). `APICompatVersion` stays 1.2.2 per #90's frozen
  decision. No new gaps found, no issues filed, no task lines proposed.

---

- **github-issues triage run, 2026-09-08: nothing proposable, nothing stale.** All six open
  issues were already covered by another mechanism. #97/#96/#95/#94 are mule-alarms (idle
  queue), not code issues — expected during low-activity windows. #93 has a pushed commit
  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated below.
  #1 untouched per standing instruction. No task lines proposed, no for-giulio close
  suggestions. (Identical to the 2026-09-07 run — issue list unchanged except the new
  alarm #97.)

---

- **Milestone recipe run, 2026-09-06: v3.0 (upstream 1.2.2 → 1.3/1.4) has no open
  machine-checkable gaps left to file; the milestone now waits on two upstream gates, and the
  `milestone-3.0` label is empty (this milestone's work lives under `upstream-parity`).** All
  the 1.3-surface gaps from the #47–#89 backlog are delivered/closed — retention ceiling #72,
  alias/group #52–#54 #59, query formats #55 #56, wire capabilities #47–#49, validation /
  error-code fidelity #57 #61 #62 #79, per-service version #77, housekeeping #73–#76 — and the
  remaining open items are already escalated, not fileable gaps: **#92 keyAgreement** (1.4
  experimental decision, escalated 2026-09-05, gated on a stable upstream v1.4.0), **#78 FDO**
  (milestone-4.0, escalated), **#93** aclhook comment rewrite (mule-review, pushed commit
  `8c61268`). `docs/UPSTREAM-EXPERIMENTAL.md` rows #67/#68 are "1.4 experimental" and reconcile
  only when upstream ships **v1.4 final** — not yet (still rc). Per **#90's frozen decision**
  (2026-08-23), `APICompatVersion` stays 1.2.2 until the FULL 1.3 surface is complete *and*
  every UPSTREAM-EXPERIMENTAL row at that level is reconciled, so the bump is not yet due. Your
  call when v1.4.0 goes stable: answer #92 (implement against the rc, re-park, or narrow #93
  only) — until then v3.0 has no open implementable work. This run filed **no issues** and
  proposes **no** "complete, cut the tag" line: the open decision set is real and is yours.

- **`milestones.md` v3.0 names a source doc that does not exist: `.mule/research/
  upstream-parity-2026-08-22.md` is absent from `main` and from `origin/mule/research`.** The
  investigation's *output* (issues #47–#89) is intact and reachable, so milestone work is
  unaffected, but if that doc is your canonical 1.3/1.4 diff it needs restoring (or the
  milestone paragraph re-pointed) before v4.0 drafting leans on it.

- **docs-sync realm-management: two ASTRATE_ config keys exist in code but are missing from
  `docs/site/configuration-reference.md`** (site prose is yours, never edited by the mule).
  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` and
  `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` are consumed in `internal/config/config.go`
  (lines 270-271 and 288); neither appears in `docs/site/configuration-reference.md`.
  Add them or confirm they are meant to be internal-only.

- **github-issues triage run, 2026-09-06: nothing proposable, nothing stale.** All six open
  issues were already covered by another mechanism. #95 and #94 are mule-alarms (idle queue),
  not code issues — expected during low-activity windows. #93 has a pushed commit (`8c61268`)
  awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the waiting-on
  row. #78 is the milestone-4.0 FDO design/investigation, already escalated below. #1
  untouched per standing instruction. No task lines proposed, no for-giulio close suggestions.
  (Identical to 2026-09-05 run — issue list unchanged except new alarm #95.)

---

- **keyAgreement: the parking condition from #51 has fired — implement now, or wait for a
  stable v1.4.0?** (issue #92, `upstream-parity`/`upstream-experimental`). #51 closed
  2026-08-22 with "parked until the upstream 1.4 experimental spec stabilizes — reopen or
  file fresh when it does". The document side has now stabilized: upstream `d084308`
  (2026-08-31) published `082-key_agreement_protocol.md`, a full 267-line wire protocol
  (topics `control/keyAgreement/0..4` — InitExchange/ExchangeResp/SecretHash/HashOk/
  ExchangeFailed at QoS 2, CBOR bodies with CDDL, `alg` 0 ECDH_P256-HKDF_SHA256-AES_256_GCM,
  CBOR-wrapped COSE_Key + 32-byte HkdfSalt, session-scoped keys, enumerated `ExchangeFailed`
  codes), and deleted the "not yet implemented" sentences the old parking quoted (#93 already
  fixed the stale ACL comment that cited them). But the spec ships only in `v1.4.0-rc.5` —
  `v1.3.3` is still the newest stable tag and Astrate targets 1.2.2 — so the document has
  stabilized and the release has not. Implementing it is the largest surface in the parity
  backlog (CBOR codec, X25519/P-256, HKDF, AES-256-GCM, a 5-state handshake machine,
  shared-secret persistence, five new error names). Your call: build against the rc now,
  re-park until v1.4.0 is a stable tag, or take only the narrow #93 fix. (Escalated again
  2026-09-05 — a prior escalation from the 2026-09-04 milestone run was lost in the queue
  rebuild.)

---

- **govulncheck GO-2026-5970: reachable DoS in golang.org/x/text (infinite loop on invalid input, fixed in v0.39.0, available v0.41.0).** Astrate pins `x/text` indirect at v0.38.0 (go.mod:97) and pgx pulls it into production: `internal/store/notify.go:59` `store.Listen` → `pgx.ConnectConfig` → `unicode/norm.*`. This is the only govulncheck symbol finding that is not test-harness-only: GO-2026-6355/6354 (x/crypto/ssh deadlocked-channel DoS) and GO-2026-6253 (moby/go-archive tar path traversal) are reachable only through testcontainers in `internal/testutil/pg.go`, i.e. never in the deployed binary. `x/text` keeps API compatibility minor-to-minor and the modules Astrate exercises (`unicode/norm` via pgx, `text/language` via jsonschema) are unchanged, so this is a fix Astrate actually needs — the hygiene recipe's highest-priority category. Not a mule task (go.mod never-touch): your decision to bump ≥v0.39.0 now or fold into the next milestone-boundary sweep. Raw: https://pkg.go.dev/vuln/GO-2026-5970. (The 2026-09-04 dep sweep did not list x/text.)

---

- **COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable, 2026-07-14; v1.4.0 is still rc-only).** Astrate's doc and `APICompatVersion` still target upstream **1.2.2** (`internal/realm/service.go:588`); v1.3.0 (2026-05-06) introduced wire-surface changes Astrate does not yet emulate, so this is a decision — adopt v1.3.2 as the compatibility target (then update the doc + bump `APICompatVersion` together, per the bump rule) or keep 1.2.2 and add a "not yet emulated" note. Wire-relevant v1.3.0 deltas (release notes): **MQTT v1 capabilities** incl. `purge_properties_compression_format` (plaintext vs zlib purge — touches the `emptyCache`/`producer,properties` contract COMPATIBILITY.md deviation 1 documents); **empty introspection now allowed**; **device registration triggers** (pairing) and **device deletion started/completed triggers** (RM — the latter two already exist as Astrate deviation 9 emits both around the synchronous delete); **FDO authentication** (pairing, disabled by default); **realm-scoped health** — upstream v1.3 added `GET /pairing/v1/{realm}/health`, which Astrate already serves (`internal/pairing/http.go:78-81`, comment already says "upstream 1.3+"), so deviation 18's wording ("which upstream 404s") is now false against 1.3 and the note should be reworded either way. Proposed doc wording (for your approval, edit to taste): in §Infrastructure differences add a sentence — *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl.

plaintext `purge_properties_compression_format`, empty-introspection allowance, device registration/deletion triggers, experimental FDO pairing auth) are not yet emulated and are out of scope until the milestone that adopts v1.3.2 as the target."* — and reword deviation 18's realm-health note from "which upstream 404s" to "added by upstream v1.3 (Astrate serves it against a 1.2.2 target; kept, matching behavior)". Raw upstream changes: [v1.3.0](https://github.com/astarte-platform/astarte/releases/tag/v1.3.0), [v1.3.2](https://github.com/astarte-platform/astarte/releases/tag/v1.3.2).

---
- **The mule's 42 straight failures were ten lint errors in its own base, not a stale branch.**
  `mule/queue` stopped taking `main` on 2026-07-27 and drifted 120 commits behind, which is
  real and is why it is being rebuilt — but it is not what blocked the work. The lint gate runs
  `golangci-lint run ./...` over the whole repo, and the branch's own `internal/flow/*` code
  carried 10 findings (1 goimports, 1 gosec, 8 revive). So every task failed the gate no matter
  how good the change was, including the four tasks queued to fix those very findings: each
  fixed one and died on the other nine. Base tests and `go vet` pass on that checkout, and
  `golangci-lint` and `govulncheck` have been installed on the Pi the whole time — the two
  things that looked broken were not. `main` is lint-clean, so rebuilding from it clears the
  deadlock. `tools/mule.sh preflight` checks the lint baseline and would have said so on day
  one; nobody ran it between 2026-08-31 and 2026-09-04. (Diagnosed 2026-09-04.)

---

- **Dependency sweep corrected: direct (pinned) deps DO have newer versions** — the 2026-09-02 note said the `go list -m -u` sweep showed "only version-skew on transitive deps", but that run hit the recipe's `head -20` cutoff (all cloud/azure/transitive) and never reached the directly-required modules. Full sweep, 2026-09-04. None of these is a fix this repo *needs*, so no bump is proposed — recorded for the decision. Per module (current → available; breaking change; repo use):
  - `github.com/coder/websocket` v1.8.14 → v1.8.15 — no breaking (patch); used in `internal/appengine/stream/ws.go`, `channels/ws.go`; worth it only for the "transmit in single frame when compression enabled" fix + read-path alloc reduction.
  - `go.etcd.io/bbolt` v1.4.3 → v1.5.0 — bbolt's semver promises no API change between patch/minor, so additive-only; used in `internal/broker/sessionstore.go`; v1.5 adds a data-file size limit and panic-recovery hardening, nothing Astrate needs.
  - `go.mongodb.org/mongo-driver/v2` v2.6.0 → v2.8.2 — the 2.8.0 breaking changes are confined to Queryable Encryption string-query options (`options.Text()`→`String()`); Astrate uses only the raw BSON API (`pkg/payload/bson.go`, `internal/engine/capabilities.go`, `bench/`) and is unaffected.
  - `github.com/nats-io/nats.go` v1.52.0 → v1.53.1 — no breaking; the headline fixes (JetStream `resetOrderedConsumer` race, KV dot-rejection) are paths Astrate does not use — `internal/engine/forward/nats.go` is core NATS publish only.
  - `github.com/prometheus/client_golang` v1.23.2 → v1.24.1 — requires Go ≥1.25 (fine, repo is 1.26.1); the breaking `LabelNames`/remote-api renames don't touch repo usage (`prometheus`/`collectors`/`promhttp` in `internal/observability/metrics.go`, flow/engine metrics); would buy `Gather()` panic-recovery and opt-in `CoalesceGather` scrape-pile-up protection.
  - `github.com/testcontainers/testcontainers-go` v0.43.0 → v0.44.0 (modules/postgres v0.42.0, modules/nats v0.43.0) — breaking in `wait.ForSQL` (callback now takes `network.Port`) and `ImageProvider` (new `PullImageWithPlatform`); Astrate's `internal/testutil/pg.go` looks unaffected but it is test-only anyway.
  - `golang.org/x/crypto` v0.53.0 → v0.56.0 — x/crypto keeps API compatibility; used only for bcrypt in `internal/auth`.
  Note (corrected 2026-09-04): `govulncheck` and `golangci-lint` **are** installed on the Pi (`/root/go/bin`, since 2026-07-28 and 2026-09-01) and `.mule/config` finds them there, so both checks were available; the sweep that produced this list simply ran without invoking govulncheck.

  **Decided 2026-09-04: no bumps.** None of the seven fixes anything this repo has, and each
  one costs a full test run to land. The standing rule instead: re-run this sweep at every
  milestone boundary (the point where `APICompatVersion` or a milestone tag moves), and bump
  only what carries a fix Astrate actually needs. `go.mod` stays on the never-touch list.

---

- **milestone 2.0 looks complete, verify and cut the tag** — all 11 `milestone-2.0` issues
  CLOSED (#23–#27, #37, #39–#43), no open issues, no new gaps after re-checking upstream
  astarte_flow block catalog against `internal/flow/` + git log (MQTT/HTTP source/sink,
  json_path_map, pure-transform set, virtual_device_pool, container block MVP, flow API,
  durable named flows all landed).   (Milestones recipe run, 2026-09-03.)

  **Decided 2026-09-04: the tag is `v0.2.0`, cut on the same day.** The `v2.0`/`v3.0` names
  are milestone names, not release versions — the project is still pre-1.0 and the version
  number keeps its own line rather than jumping two majors to match a milestone label.

---
- ~~`docs/site/appengine-api.md` documents `GET` and `DELETE` on `/appengine/v1/<realm>/groups/<name>`~~
  — **resolved 2026-09-04, and the original note was half wrong.** `GET /groups/{group}` does
  exist (`internal/appengine/http.go:69`); only `DELETE` was absent, and it is absent from
  upstream's own spec too (`docs/api/astarte_appengine_api.yaml` has no `/groups/{group}`
  path at all). So the docs line was spurious, not the code. Line dropped.
---
- ~~`Router.Submit` TOCTOU on the `closed` flag~~ — **fixed 2026-09-04**, and none of the three
  options in the original note was taken: (b) was wrong (a `select` does not save a send on a
  closed channel — it panics either way). The fix is that `Drain` no longer closes the lane
  channels at all; it closes `quit` under a write lock while `Submit` holds the read lock
  across its send, so a lane can never be retired underneath an in-flight sender, and the lanes
  exit on `quit` after draining what is buffered. `TestRouter_SubmitParkedWhenDrainRuns`
  reproduces the old panic deterministically and passes on the fix.
---
- ~~**#87 `lua_map` — needs embedded Lua runtime, parked.**~~ — **closed 2026-09-04.**
  Embedding a Lua VM is not machine-checkable by the mule and Lua flow support is on no active
  roadmap. Reopen in a minute if that changes.
- **#78 FDO device onboarding — milestone-4.0, investigation phase.** Too large for a single
  mule task; the investigation work (reading upstream's TO2 handling, inventorying endpoints,
  schema and keys) is a multi-session project. Parked until the v3.0 queue clears and this
  becomes the next milestone target.
- **#1 is never to be raised again.** Giulio's standing instruction, 2026-09-04: it stays open
  permanently and is not a candidate for closing, triage, or a for-giulio entry. Do not
  propose it again.
---

- ~~**Flow v2.0: named multi-instance flows + pipeline config?**~~ — **decided
  2026-07-29: (b) named multi-instance + config.** Design then implement (#40).
  Doc: `docs/handoff/flow-v2-decisions-2026-07-29.md`.
- ~~**Flow v2.0: durable `flows` table vs in-memory only?**~~ — **decided
  2026-07-29: (b) durable records; `auto_restart` default true, optional never;
  process rehydrates on boot; fail loudly on bad pipeline/start.** Filed #41;
  edge-case follow-ups #42. Same design package as #40.
- ~~**Flow v2.0: containers in scope?**~~ — **decided 2026-07-29: yes, phased
  PoC → MVP (#43).** Native Lua/MQTT blocks are *not* a v2.0 gate (containers cover
  custom logic). Doc: `docs/handoff/flow-v2-decisions-2026-07-29.md`.
- ~~`device_deletion_started`/`device_deletion_finished` trigger events are not emitted~~ —
  **decided 2026-07-27: emit both, back-to-back, around the synchronous delete.** Filed as
  issue #21 (`mule`). (Cross-project survey, 2026-07-27,
  `.mule/research/survey-2026-07-27.md` source 4.)
- ~~Mustache trigger-action templates are accepted but not rendered~~ — **decided
  2026-07-27: implement it.** Guiding principle clarified: Astarte compatibility means
  SDK/wire compatibility, not minimum dependency count — Astrate is allowed to be a
  compatible *superset*. Library picked: `github.com/cbroglie/mustache`. Filed as issue #22
  (`mule`). (Same survey, source 4.)
- ~~`value_change`/`value_change_applied`/`path_created`/`path_removed`/`value_stored` trigger
  types compile but never fire~~ — **resolved: implemented in commit 6bd14a7 (2026-08-22)**,
  and the cost question answered by #20's bench on the Legion Go (2026-08-24 closeout):
  ~0.13 ms keyed read per message, −0.7% throughput even paid once per message over REST.
  Performance is not the constraint for these types (nor, a fortiori, for the group-scoped
  line below). (Same survey, source 4.)
- **Group-scoped triggers (`group_name` on device/data triggers) compile but never match**
  (`internal/engine/triggers/match.go:11-12`). Decision deferred, tied to issue #17
  (group-WATCH-path reconciliation, trickle work, not mule): whatever group-membership
  mechanism comes out of that phase should also report the perf cost for this decision —
  noted in a comment on #17 so it isn't benchmarked twice. (Same survey, source 4.)

---

- ~~The Pi cannot run the race detector~~ — **resolved 2026-07-27** by installing Go 1.26.5
  as a userland toolchain on the Legion Go (`~/.local/go`, no root, `rm -rf` to undo). The
  Pi still cannot run `-race` (39-bit VMA kernel vs the 48 ThreadSanitizer needs), so its
  gate remains `go vet ./... && go test ./...` — but race coverage now exists on the Legion
  Go, where the full suite runs clean in ~40s on 16 cores. The standing `race-check` task is
  the concurrency gate. Concurrency work is queueable again, provided the race-check runs
  after it.
- ~~golangci-lint is not installed on the Pi~~ — **resolved 2026-07-28**: installed v2.12.2
  via `go install` (the prebuilt-binary installer's published sha256 for linux/arm64 did not
  verify, reproducibly, so built from source instead). `mule.service` now sets
  `MULE_LINT_CMD=golangci-lint run ./...` and has `/root/go/bin` on `PATH`; the lint gate
  runs starting with the next tick.
- ~~`/root/astrate` on the Pi has uncommitted work~~ — **resolved 2026-07-27** with the new
  `tools/reconcile.sh`: rescued onto `origin/wip/DietPi-20260727T171543Z` (pushed, not
  reviewed — read the diff before merging anything from it) and `/root/astrate` is now clean
  on `main`. Also had to set `commit.gpgsign false` locally on that clone first — it had
  signing on with no working gpg-agent for a non-interactive session, same class of problem
  `/root/astrate-mule` already solved the same way (your call, confirmed 2026-07-27: align it
  with the existing unattended-Pi-clone convention rather than fix gpg-agent or bypass signing
  silently inside the script).

## 2026-08-23 — FDO promoted to milestone-4.0 candidate (Giulio's decision, recorded)

#78 is no longer parked: zero-touch onboarding is strategic for commercial
viability. Scope frozen on the issue: owner-side TO1/TO2 in our Pairing
service only (last mile, like upstream), reuse fdo-rs for
manufacturing/rendezvous, acceptance = official `astarte-device-fdo-rust`
SDK completes onboarding against Astrate, docs as a first-class deliverable.
When v3.0 is marked DONE, the v4.0 section of `.mule/milestones.md` should be
drafted with this investigation as its first item (issue #78 has the full
verified context).
- **The mule has been idle 17h.** Filed by the dead-man's switch; see journalctl on the Pi.

## 2026-09-25 — the `-tags nats` half of the tree is compiled by no gate

`cmd/astrate/newnats_nats.go`, `internal/engine/forward/nats.go` and
`internal/engine/forward/nats_test.go` are all behind `//go:build nats`, and nothing
builds or tests that tag: not `make build`/`make test`/`make test-integration`/
`make test-e2e` (Makefile:23, 38, 42, 46), not the mule gate, not the lint target.
So the whole NATS trigger-forwarding path is invisible to CI — while
`config.validate` accepts `triggers.forward.kind = "nats"` (internal/config/config.go:339-343)
and the untagged binary then fails at boot with the build-tag message
(cmd/astrate/newnats_default.go). I measured it today: `go build -tags nats ./cmd/...`
and `go vet -tags nats ./cmd/...` both exit 0, so there is nothing broken — just nothing
keeping it unbroken. Worth a `make build-nats` (and the NATS test tier, which needs a
container, so it would be a `[legion]` line) in the gate; not a queue line because the
only test I can imagine for it shells out to `go build`, and the mule's own check strips
the implementation and the guard together, so it would never fail.
Not queued for the same reason: `IdleTimeout` is missing on the HTTP server
(cmd/astrate/main.go:221) and `-healthcheck` ignores the config file
(cmd/astrate/main.go:531) — both real, both yours to call, neither has a test worth
writing. See `.mule/reviews/cmd-astrate-2026-09-25.md`.
- **The mule has been idle 16h.** Filed by the dead-man's switch; see journalctl on the Pi.

## 2026-10-04 — three housekeeping drifts in `docs/site/` the recipe is not allowed to touch

Ran the docs-sync recipe over the housekeeping surface (code routes vs
`docs/api/astarte_housekeeping_api.yaml`). The spec half is queued as four lines in
`.mule/todo.md`; the prose half needs you, because `docs/site/` is on the never-touch list.

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

→ Read .mule/todo.md [limit=40]

$ git log --oneline -8 main 2>/dev/null; echo "=== 8c61268 ==="; git log --oneline -1 8c61268 2>/dev/null || echo "not found"; git branch -a --contains 8c61268 2>/dev/null
ffa5ce9 Track the docs phase history and benchmark results
156f28d Serve Swagger UI and the OpenAPI specs from the astrate binary
283f966 docs: add operations, troubleshooting, configuration and migration pages
3758110 mule: ignore tick's own working files, and reach the Legion Go by ssh alias
956d2a7 Add solo-mule mode: an unsupervised opencode worker driven by a task queue
3ec2753 Select the trigger forwarder from [triggers.forward] config
9c6b22b Add an HTTP bus Forwarder that posts a custom trigger action verbatim
cc34e6a docs: add MkDocs site, Swagger UI, GitHub Pages workflow
=== 8c61268 ===
8c61268 mule: log issue-93
* mule/queue
  remotes/origin/mule/queue

$ gh issue list --state open --limit 40 --json number,title,createdAt --template '{{range .}}#{{.number}} {{.createdAt}} {{.title}}{{"\n"}}{{end}}' | head -35
#120 2026-10-07T11:03:45Z mule: nothing has landed in 14h
#119 2026-10-06T10:58:45Z mule: nothing has landed in 16h
#118 2026-10-05T10:54:06Z mule: nothing has landed in 14h
#117 2026-10-04T11:17:06Z mule: nothing has landed in 15h
#116 2026-10-03T11:02:14Z mule: nothing has landed in 14h
#115 2026-10-02T11:22:50Z mule: nothing has landed in 14h
#114 2026-09-27T20:09:31Z mule: nothing has landed in 8h
#113 2026-09-27T11:07:48Z mule: nothing has landed in 14h
#112 2026-09-26T11:04:46Z mule: nothing has landed in 14h
#111 2026-09-25T11:01:07Z mule: nothing has landed in 16h
#110 2026-09-24T10:54:08Z mule: nothing has landed in 14h
#109 2026-09-22T10:54:47Z mule: nothing has landed in 14h
#108 2026-09-21T10:54:46Z mule: nothing has landed in 15h
#107 2026-09-19T10:49:12Z mule: nothing has landed in 14h
#106 2026-09-18T10:54:09Z mule: nothing has landed in 14h
#105 2026-09-17T10:51:07Z mule: nothing has landed in 14h
#104 2026-09-16T11:12:14Z mule: nothing has landed in 14h
#103 2026-09-14T11:27:08Z mule: nothing has landed in 16h
#102 2026-09-13T11:09:08Z mule: nothing has landed in 14h
#101 2026-09-12T10:58:09Z mule: nothing has landed in 14h
#100 2026-09-11T10:55:09Z mule: nothing has landed in 15h
#99 2026-09-10T11:08:07Z mule: nothing has landed in 14h
#98 2026-09-09T11:05:08Z mule: nothing has landed in 16h
#97 2026-09-08T11:19:07Z mule: nothing has landed in 15h
#96 2026-09-07T11:15:07Z mule: nothing has landed in 14h
#95 2026-09-06T11:17:07Z mule: nothing has landed in 15h
#94 2026-09-05T11:00:06Z mule: nothing has landed in 15h
#93 2026-09-04T19:44:12Z In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
#92 2026-09-04T19:43:54Z keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds
#78 2026-08-22T02:24:44Z FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
#1 2026-03-13T21:54:29Z Provide an Open Source IoT Platform unironically, unironically

$ git status -sb | head -3; echo "=== main dates ==="; git log -6 --format='%h %ad %s' --date=short main; echo "=== mule/queue dates ==="; git log -6 --format='%h %ad %s' --date=short origin/mule/queue; echo "=== is 8c61268 in main? ==="; git merge-base --is-ancestor 8c61268 main && echo YES || echo NO
## mule/queue
=== main dates ===
ffa5ce9 2026-07-27 Track the docs phase history and benchmark results
156f28d 2026-07-27 Serve Swagger UI and the OpenAPI specs from the astrate binary
283f966 2026-07-27 docs: add operations, troubleshooting, configuration and migration pages
3758110 2026-07-27 mule: ignore tick's own working files, and reach the Legion Go by ssh alias
956d2a7 2026-07-27 Add solo-mule mode: an unsupervised opencode worker driven by a task queue
3ec2753 2026-07-27 Select the trigger forwarder from [triggers.forward] config
=== mule/queue dates ===
dde64e0 2026-10-07 mule: recipe milestones ran (628s)
becf3cd 2026-10-07 mule: recipe hygiene timed out
7a24500 2026-10-06 mule: blocked docs-native-socket-event-schema
bc3dc9a 2026-10-06 mule: log docs-native-socket-security-scheme
87c18cf 2026-10-06 mule: docs-native-socket-security-scheme [auto]: `/astrate/v1/{realm}/socket` is documented as if public in docs/api/astrate_native_api.yaml — the spec has no `components.securitySchemes` at all and root `security: []` (yaml:21), and the operation declares no `security`, so a client generated from this spec sends no `Authorization` header and is answered 401 `{"errors":{"detail":"Unauthorized"}}` by `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112, bearer token read only from the Authorization header at middleware.go:131, 178-196). Add a scheme mirroring `a_aea` (docs/api/astarte_appengine_api.yaml:1436-1443: apiKey in header `Authorization`, described as a realm JWT carrying the `a_ch` claim) and set `security` on this operation only — root `security: []` stays so health/readiness/metrics/version remain documented as unauthenticated. If you model the Phoenix twin too, its credential is `?token=` in the query (internal/appengine/channels/ws.go:53-54), i.e. apiKey `in: query`, not header. Description notes worth carrying: the scheme is case-insensitive `Bearer` with an optional colon, and an unknown realm is 401 not 404 (middleware.go:68-72). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
3676ad9 2026-10-06 mule: log docs-native-socket-sse-exact-accept
=== is 8c61268 in main? ===
NO

$ rg -l "nothing has landed" --hidden -g '!.git' . 2>/dev/null
./tools/mule.sh
./.mule/for-giulio.md
./.mule/reports/log/20261003T180922Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260920T185539Z-compat-note-custom-action-policy-boundary-done.md
./.mule/reports/log/20260917T194954Z-recipe-docs-sync-proposed.md
./.mule/reports/log/20260923T181743Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20261004T183240Z-recipe-github-issues-proposed.md
./.mule/reports/log/20261005T193452Z-recipe-milestones-proposed.md
./.mule/reports/log/20260916T175801Z-recipe-milestones-proposed.md
./.mule/reports/log/20260920T111057Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260911T191541Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260923T174341Z-recipe-milestones-proposed.md
./.mule/reports/log/20260912T171256Z-recipe-milestones-proposed.md
./.mule/reports/log/20260913T170550Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260926T173418Z-recipe-astarte-upstream-timeout.md
./.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md
./.mule/reports/log/20260908T112032Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260918T114100Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260926T164829Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260908T193127Z-recipe-milestones-proposed.md
./.mule/reports/log/20260915T193626Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260913T203209Z-recipe-milestones-proposed.md
./.mule/reports/log/20260926T174259Z-recipe-code-review-proposed.md
./.mule/reports/log/20260917T190323Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260912T172901Z-recipe-github-issues-proposed.md
./.mule/reports/log/20261001T204604Z-fdo-rc6-scope-delta-for-giulio-done.md
./.mule/reports/log/20260917T190703Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260925T200001Z-docs-sync-rm-async-operation-param-done.md
./.mule/reports/log/20260923T174805Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260924T193516Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260918T165351Z-compat-note-v134-done.md
./.mule/reports/log/20260925T192553Z-recipe-docs-sync-proposed.md
./.mule/reports/log/20261001T201010Z-recipe-github-issues-proposed.md
./.mule/reports/log/20261003T175916Z-recipe-milestones-proposed.md
./.mule/reports/log/20260917T184319Z-recipe-milestones-timeout.md
./.mule/reports/log/20261004T181934Z-recipe-milestones-proposed.md
./.mule/reports/log/20260915T185804Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260905T181219Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260918T194216Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260918T113442Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260913T203800Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md
./.mule/reports/log/20260908T193321Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260913T203830Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260925T165408Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260920T204223Z-recipe-github-issues-proposed.md
./.mule/reports/log/20261005T195440Z-recipe-github-issues-proposed.md
./.mule/reports/log/20261002T113048Z-compat-note-v14-rc6-done.md
./.mule/reports/log/20260915T184406Z-recipe-milestones-proposed.md
./.mule/reports/log/20260907T174812Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260910T110956Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260925T115624Z-recipe-milestones-timeout.md
./.mule/reports/log/20260924T192153Z-recipe-milestones-proposed.md
./.mule/reports/log/20260911T190417Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260916T184218Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260916T183922Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260915T105740Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260908T193026Z-recipe-milestones-proposed.md
./.mule/reports/log/20260906T181352Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260925T170557Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260910T200120Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260918T110039Z-recipe-milestones-proposed.md
./.mule/reports/log/20260924T182956Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260921T185351Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md
./.mule/reports/log/20261003T195023Z-appengine-missing-required-422-done.md
./.mule/reports/log/20261006T175547Z-compat-note-v135-done.md
./.mule/reports/log/20261006T165612Z-recipe-github-issues-proposed.md
./.mule/reports/log/20261006T184046Z-recipe-docs-sync-proposed.md

$ rg -l "dead-man|deadman|dead_man" --hidden -g '!.git' . 2>/dev/null
./.agents/skills/astrate-dashboard/SKILL.md
./.mule/reports/log/20260926T113138Z-recipe-milestones-proposed.md
./.mule/reports/log/20261003T180922Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260928T105915Z-recipe-docs-sync-proposed.md
./.mule/reports/log/20261005T193452Z-recipe-milestones-proposed.md
./.mule/reports/log/20260918T200041Z-recipe-code-review-proposed.md
./.mule/reports/log/20260909T230411Z-recipe-milestones-proposed.md
./.mule/reports/log/20260911T185303Z-recipe-milestones-proposed.md
./.mule/for-giulio.md
./.mule/reports/log/20261005T201617Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260926T164829Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260920T185956Z-recipe-docs-sync-proposed.md
./tools/mule.sh
./.mule/reports/log/20260924T175556Z-recipe-milestones-proposed.md
./.mule/reports/log/20261006T174353Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260926T194649Z-recipe-docs-sync-proposed.md
./.mule/reports/log/20261004T183240Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260918T172636Z-recipe-code-review-proposed.md
./.mule/reports/log/20260916T115527Z-recipe-docs-sync-proposed.md
./.mule/reports/log/20260906T180442Z-recipe-milestones-proposed.md
./.mule/reports/log/20261006T141129Z-recipe-milestones-proposed.md
./.mule/reports/log/20260912T171256Z-recipe-milestones-proposed.md
./.mule/reports/log/20261001T204604Z-fdo-rc6-scope-delta-for-giulio-done.md
./.mule/reports/log/20260923T174805Z-recipe-github-issues-proposed.md
./.mule/reports/log/20261001T201010Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260924T182956Z-recipe-github-issues-proposed.md
./.mule/reports/log/20261007T113844Z-recipe-milestones-proposed.md
./.mule/reports/log/20261002T114427Z-recipe-code-review-proposed.md
./.mule/reports/log/20260917T184319Z-recipe-milestones-timeout.md
./.mule/reports/log/20261004T181934Z-recipe-milestones-proposed.md
./.mule/reports/log/20260915T185804Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260916T183922Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260925T112307Z-recipe-hygiene-timeout.md
./.mule/reports/log/20260924T174100Z-recipe-hygiene-proposed.md
./.mule/reports/log/20261003T183445Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20261005T195440Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260910T195304Z-recipe-milestones-proposed.md
./.mule/reports/log/20260925T192553Z-recipe-docs-sync-proposed.md
./.mule/reports/log/20261002T113048Z-compat-note-v14-rc6-done.md
./.mule/reports/log/20260923T112257Z-recipe-code-review-proposed.md
./.mule/reports/log/20261004T112505Z-recipe-docs-sync-proposed.md
./.mule/reports/log/20261006T165612Z-recipe-github-issues-proposed.md
./.mule/reports/log/20261003T175916Z-recipe-milestones-proposed.md
./.mule/reports/log/20260925T174727Z-recipe-code-review-proposed.md
./.mule/reports/log/20261003T195023Z-appengine-missing-required-422-done.md
./.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md

→ Read .mule/reports/log/20261006T165612Z-recipe-github-issues-proposed.md

$ grep -n "triage run, 2026-10" .mule/for-giulio.md; echo "=== all 2026-10 dated bullet lines ==="; grep -n "^- \*\*.*2026-10" .mule/for-giulio.md | head -20; echo "=== wc ==="; wc -l .mule/for-giulio.md
21:- **github-issues triage run, 2026-10-06: nothing proposable for the 24th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented or closed on GitHub.** The recipe's `--limit 40` command prints all **30** open issues: **26** `mule-alarm` **#94–#119** plus the same four non-alarm, none of them proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (`milestone-4.0`, already escalated), and one the standing instruction excludes from triage entirely. New since the 2026-10-05 run: **#119**, created 2026-10-06T10:58:45Z, "nothing has landed in 15h" — a real window, not an artefact: the land before it was `230dd11` at 2026-10-05T21:04:09+02:00 (19:04Z), 15.9h earlier, and it was ended by today's first land `cbdd3e8` at 16:05:39+02:00 (14:05Z), three hours after the alarm. So it belongs to the genuine 14–16h band the 2026-10-05 entry measured, i.e. overnight silence between evening and afternoon batches. **Proposal, extending the 2026-10-05 one: close #94–#118, all twenty-five** — each is a one-day alarm superseded by the next, and nothing reads issue state; leave **#119** (today's, real window, already ended) to self-expire. The `check_pulse`-before-`beat` ordering fix (move the call at `tools/mule.sh:633` to after the tick's `beat`) is still escalated above, still yours, and still unactioned — today's alarm is not evidence for it either way.
36:- **github-issues triage run, 2026-10-05: nothing proposable for the 23rd run, and the alarm pile's arithmetic is now fully pinned — 27 gaps over 8h since 2026-08-20 and **not one of them between 8h and 13h**, so the threshold cannot be a number. Two things are new: today's #118 was silenced by a land 5 minutes later, and #93's commit has now been unreviewed for 31 days.** 29 open issues (authoritative `--limit 100`: 25 alarms **#94–#118**, one new today, plus exactly four non-alarm — **#93** aclhook comment `mule-review`, **#92** keyAgreement parked, **#78** FDO `milestone-4.0`, **#1** untouched per standing instruction). Note the recipe's own `--limit 40` command printed **all 29** today, so the "silently dropped #94" parenthetical in the 2026-10-04 entry did not reproduce — treat the limit as adequate at this pile size. Zero machine-checkable candidates, so **no `.mule/todo.md` lines added, no `gh issue create`, and nothing commented, closed or edited on GitHub**.
52:- **github-issues triage run, 2026-10-04: nothing proposable for the 22nd run, and the reason the alarm pile has never been actioned is finally located — `mule/queue` forked from `main` 27 minutes before your four commits landed there, so this Pi has been running for a month with a dead-man's switch, a dashboard and a parked-decision check that all point at things which do not exist here.** 28 open issues (authoritative `--limit 100` list: `1, 78, 92, 93, 94…117`; note the recipe's own `--limit 40` command printed 27 of these and silently dropped **#94**, which is open with `mule-alarm` — so treat "27 open" in the 2026-10-03 entry as 28). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`; `8c61268` re-verified still an ancestor of `origin/mule/queue` this run), **#92** keyAgreement (parked on a stable v1.4.0 — measured again today, still unmet: newest upstream stable `v1.3.4`, newest overall `v1.4.0-rc.6`, 2026-09-30), **#78** FDO (`milestone-4.0`, already escalated), **#1** untouched per standing instruction.
78:- **github-issues triage run, 2026-10-03: nothing proposable for the 21st run — and the alarm pile finally produced a falsifiable demonstration of what the 2026-09-26 entry diagnosed, on both sides of two fresh alarms.** 27 open issues, the same set plus **#115** (10-02 11:22:50Z) and **#116** (today, 10-03 11:02:14Z). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`; `8c61268` re-verified still reachable from `origin/mule/queue` this run, so it is pushed and waiting only on your read), **#92** keyAgreement (parked on a stable v1.4.0 per your 2026-09-04 decision — today's milestone run re-verified there is still no stable v1.4.0, so the condition is unmet and the `waiting-on.md` gap is unchanged; see the entry below), **#78** FDO (`milestone-4.0`, already escalated), **#1** untouched per standing instruction.
121:- **github-issues triage run, 2026-10-01: nothing proposable, and the 2026-09-26 reading of the alarm pile is now superseded — this silence is real, and the dead-man's switch is latched off so it cannot report it.** 25 open issues, the same set plus **#113** (09-27 11:07Z) and **#114** (09-27 20:09Z, "nothing has landed in 8h"). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction.
=== all 2026-10 dated bullet lines ===
13:- **COMPATIBILITY.md wording proposal for upstream v1.3.5 (newest stable, 2026-10-05; v1.4.0 is still rc.6-only) — a one-word delta on the open v1.3.2/v1.3.4 proposal below, nothing re-derived.** v1.3.5 is maintenance-only: `gh api repos/astarte-platform/astarte/releases` shows it published `2026-10-05T14:11:44Z`, not a prerelease, with a body of a single `### Fixed` line — "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available" — and `compare/v1.3.4...v1.3.5` is 4 commits / 35 files whose only non-`mix.exs`/`mix.lock` files are `.github/workflows/astarte-build-workflow.yaml`, `CHANGELOG.md`, `README.md`, `astarte-dashboard/package.json` and `docker-compose.yml`. Not one source file changed, so no route, MQTT topic, control message or interface-schema field moved (Astrate implements no `data_updater_plant` and no Erlang RPC server). **The only change a v1.3.5-aware doc carries: the v1.3.4 entry's §Infrastructure-differences sentence — "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target" — reads "…(**newest stable v1.3.5**) as the target".** Everything else stands unchanged: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"); that entry's "v1.4.0 is still rc.5-only" header is already superseded to rc.6 by the open rc.6 delta entry below, not here. This is the same one-word delta recorded on 2026-10-05 in this file ("newest stable v1.3.5"), folded in as asked rather than re-derived. Applied to nothing: `docs/COMPATIBILITY.md` untouched and still 1.2.2-targeted (it contains no `1.3` reference at all), `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). Raw: [v1.3.5](https://github.com/astarte-platform/astarte/releases/tag/v1.3.5).
21:- **github-issues triage run, 2026-10-06: nothing proposable for the 24th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented or closed on GitHub.** The recipe's `--limit 40` command prints all **30** open issues: **26** `mule-alarm` **#94–#119** plus the same four non-alarm, none of them proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (`milestone-4.0`, already escalated), and one the standing instruction excludes from triage entirely. New since the 2026-10-05 run: **#119**, created 2026-10-06T10:58:45Z, "nothing has landed in 15h" — a real window, not an artefact: the land before it was `230dd11` at 2026-10-05T21:04:09+02:00 (19:04Z), 15.9h earlier, and it was ended by today's first land `cbdd3e8` at 16:05:39+02:00 (14:05Z), three hours after the alarm. So it belongs to the genuine 14–16h band the 2026-10-05 entry measured, i.e. overnight silence between evening and afternoon batches. **Proposal, extending the 2026-10-05 one: close #94–#118, all twenty-five** — each is a one-day alarm superseded by the next, and nothing reads issue state; leave **#119** (today's, real window, already ended) to self-expire. The `check_pulse`-before-`beat` ordering fix (move the call at `tools/mule.sh:633` to after the tick's `beat`) is still escalated above, still yours, and still unactioned — today's alarm is not evidence for it either way.
27:- **Milestone recipe run, 2026-10-06: v3.0 verified unchanged — still gated on upstream — and the master delta since the 2026-10-05 run is a closed set of 13 commits, every one wire-inert for Astrate. No issue filed, no `.mule/todo.md` line queued.** Release sweep is byte-identical to yesterday's (`gh api repos/astarte-platform/astarte/releases`): **no stable v1.4.0** — newest stable `v1.3.5` (2026-10-05), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition (your 2026-09-04 "wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) stays unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** — this milestone's whole ledger lives under `upstream-parity`, where the open set is unchanged: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`); **#47–#89** otherwise all closed. Non-alarm open issues total exactly four (plus the 25 `mule-alarm` #94–#118): the three above and **#1** (untouched per standing instruction). `APICompatVersion` still `"1.2.2"` (`internal/realm/service.go:588`, re-checked).
32:- **astarte-upstream recipe run, 2026-10-05 (evening): no upstream movement that touches anything Astrate implements, so nothing was queued — no `.mule/todo.md` line, no `gh issue create`, nothing commented, closed or edited on GitHub.** One thing *is* new since this recipe's last run (`2026-10-04T20:40:52+02:00` ran; the `2026-10-05T04:29:10+02:00` run timed out): stable **v1.3.5**, published `2026-10-05T14:11:44Z`. The milestone recipe already recorded it as inert today; I re-measured it rather than trusting that, because "a release note is a claim, not a fact" is this recipe's own rule. `gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5` is **4 commits, 35 files** — `8d47574f` + `4c4222cf` `chore: update horde`, `c4c63ebe` `chore(rm): bump xandra to 0.19.4`, `8f0f5ef1` `chore: prepare v1.3.5` — and **every app/lib file in it is a `mix.exs`/`mix.lock` version bump**. The rest is `CHANGELOG.md` (+6: one `### Fixed` entry, "[astarte_data_updater_plant] Increase RPC server reliability"), `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `README.md`, `astarte-dashboard/package.json`. Not one source file changed, so no route, MQTT topic, BSON encoding, control message or interface-schema field moved: Astrate implements no `data_updater_plant` and no Erlang RPC server. Wire-inert, same shape as v1.3.2/v1.3.3/v1.3.4. **Master since `v1.4.0-rc.6` is still the single commit `731bcddc`** ("chore: forward port release-1.4"; `gh api .../commits?since=2026-10-05T09:34:00Z` returns it and nothing newer), the migration-relocation commit the milestone run measured today as inert — so both wire-bearing post-rc.6 commits (`b6d46ad4`, `958bb9fe`) still have their owners in the queue and nothing new joins them.
36:- **github-issues triage run, 2026-10-05: nothing proposable for the 23rd run, and the alarm pile's arithmetic is now fully pinned — 27 gaps over 8h since 2026-08-20 and **not one of them between 8h and 13h**, so the threshold cannot be a number. Two things are new: today's #118 was silenced by a land 5 minutes later, and #93's commit has now been unreviewed for 31 days.** 29 open issues (authoritative `--limit 100`: 25 alarms **#94–#118**, one new today, plus exactly four non-alarm — **#93** aclhook comment `mule-review`, **#92** keyAgreement parked, **#78** FDO `milestone-4.0`, **#1** untouched per standing instruction). Note the recipe's own `--limit 40` command printed **all 29** today, so the "silently dropped #94" parenthetical in the 2026-10-04 entry did not reproduce — treat the limit as adequate at this pile size. Zero machine-checkable candidates, so **no `.mule/todo.md` lines added, no `gh issue create`, and nothing commented, closed or edited on GitHub**.
44:- **Milestone recipe run, 2026-10-05: v3.0's gate is unchanged, but the release sweep finally has something new to swallow — upstream shipped stable v1.3.5 today, seven hours ago, and it is wire-inert for Astrate. No issue filed, no task line queued.**  **The new release, measured.** `gh api repos/astarte-platform/astarte/releases` now shows a stable **v1.3.5** at `2026-10-05T14:11:44Z` where every prior run of this recipe (2026-09-10 onward, 20+ entries in this file) saw `v1.3.4` (2026-09-18) as newest. `compare/v1.3.4...v1.3.5` is **4 commits, 35 files**: two `chore: update horde` (`8d47574f`, `4c4222cf`), `chore(rm): bump xandra to 0.19.4` (`c4c63ebe`), `chore: prepare v1.3.5` (`8f0f5ef1`). The file list is `mix.exs`/`mix.lock` version bumps across all apps and libs, `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `astarte-dashboard/package.json`, `doc/mix.exs`, `README.md` — plus a **single** CHANGELOG entry, `## [1.3.5] - 2026-10-05` / `### Fixed` / "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available". That is the whole delta and none of it is Astarte-facing: Astrate implements no data_updater_plant and no Erlang RPC server (the only `rpc` hit under `internal/` is a comment at `internal/pairing/http.go:300` about AMQP reply-code conflicts), so no HTTP route, MQTT topic, AMQP control message or schema field moved. Same shape as the v1.3.2/v1.3.3/v1.3.4 maintenance releases already recorded here, so it is **not a v3.0 gap** and there is nothing to file.
52:- **github-issues triage run, 2026-10-04: nothing proposable for the 22nd run, and the reason the alarm pile has never been actioned is finally located — `mule/queue` forked from `main` 27 minutes before your four commits landed there, so this Pi has been running for a month with a dead-man's switch, a dashboard and a parked-decision check that all point at things which do not exist here.** 28 open issues (authoritative `--limit 100` list: `1, 78, 92, 93, 94…117`; note the recipe's own `--limit 40` command printed 27 of these and silently dropped **#94**, which is open with `mule-alarm` — so treat "27 open" in the 2026-10-03 entry as 28). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`; `8c61268` re-verified still an ancestor of `origin/mule/queue` this run), **#92** keyAgreement (parked on a stable v1.4.0 — measured again today, still unmet: newest upstream stable `v1.3.4`, newest overall `v1.4.0-rc.6`, 2026-09-30), **#78** FDO (`milestone-4.0`, already escalated), **#1** untouched per standing instruction.
65:- **Milestone recipe run, 2026-10-04: the gate is unchanged as always, but the recipe's standing item finally produced a real finding — both rows of `docs/UPSTREAM-EXPERIMENTAL.md` are tagged "1.4 experimental" and neither of the two features upstream calls experimental. It is a wording/policy call on two files that are yours, so this run filed no issue and queued no task.**
76:- **The Legion Go answers ssh but its Docker is not running, so every "probe live upstream" line is blocked on you, not on the hardware (measured 2026-10-03).** `ssh -o ConnectTimeout=5 legion 'echo ok'` returns `ok`, but `docker ps` fails with `dial unix /home/atsetilam/.docker/desktop/docker.sock: no such file or directory`, `systemctl is-active docker` answers `inactive`, and `sudo -n true` answers `è necessaria una password` — so the mule cannot start it. It is Docker Desktop's socket, not a rootful `dockerd` the mule could launch as a user, and no rootless extras are installed. Consequence for the queue: `appengine-payload-reason-status-map` (the remaining eight `payload.Reason` values that still answer 500) cannot be measured the way its line demands. Today's `appengine-missing-required-422` was therefore mapped from the frozen upstream source at `b6d46ad4` (`Device.validate_required_mappings/2` → `fallback_controller.ex:215-219` → `ErrorView "422_missing_required_mapping.json"`), and the commit says so plainly. Nothing about that sibling line needs rewriting — only Docker up on the Legion Go.
78:- **github-issues triage run, 2026-10-03: nothing proposable for the 21st run — and the alarm pile finally produced a falsifiable demonstration of what the 2026-09-26 entry diagnosed, on both sides of two fresh alarms.** 27 open issues, the same set plus **#115** (10-02 11:22:50Z) and **#116** (today, 10-03 11:02:14Z). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`; `8c61268` re-verified still reachable from `origin/mule/queue` this run, so it is pushed and waiting only on your read), **#92** keyAgreement (parked on a stable v1.4.0 per your 2026-09-04 decision — today's milestone run re-verified there is still no stable v1.4.0, so the condition is unmet and the `waiting-on.md` gap is unchanged; see the entry below), **#78** FDO (`milestone-4.0`, already escalated), **#1** untouched per standing instruction.
86:- **Milestone recipe run, 2026-10-03: v3.0's gate is unchanged, and this run's one new finding is not a gap in Astrate's code — it is that this milestone's declared source document can never be committed, so every run re-derives the milestone from the issue bodies.**
121:- **github-issues triage run, 2026-10-01: nothing proposable, and the 2026-09-26 reading of the alarm pile is now superseded — this silence is real, and the dead-man's switch is latched off so it cannot report it.** 25 open issues, the same set plus **#113** (09-27 11:07Z) and **#114** (09-27 20:09Z, "nothing has landed in 8h"). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction.
=== wc ===
643 .mule/for-giulio.md

$ grep -n "MULE_ALARM\|alarm" tools/mule.sh | head -30
55:# Deliberately NOT the queue label: an alarm about the mule being idle must never become a
57:MULE_ALARM_LABEL="${MULE_ALARM_LABEL:-mule-alarm}"
737:beat() { date +%s > "$MULE/.heartbeat"; rm -f "$MULE/.alarmed"; }
747:  # One alarm per silence, not one per tick — the alarm is only useful if it is rare.
748:  [ -f "$MULE/.alarmed" ] && return 0
749:  : > "$MULE/.alarmed"
750:  bad "no work has landed in ${age}h — raising the alarm"
767:      --label "$MULE_ALARM_LABEL" >/dev/null 2>&1 \
768:      && note "filed an idle alarm on GitHub" \

$ git log --format='%h %aI %s' --since='2026-10-05' origin/mule/queue | head -30
dde64e0 2026-10-07T13:38:44+02:00 mule: recipe milestones ran (628s)
becf3cd 2026-10-07T13:25:35+02:00 mule: recipe hygiene timed out
7a24500 2026-10-06T22:49:12+02:00 mule: blocked docs-native-socket-event-schema
bc3dc9a 2026-10-06T22:29:27+02:00 mule: log docs-native-socket-security-scheme
87c18cf 2026-10-06T22:29:27+02:00 mule: docs-native-socket-security-scheme [auto]: `/astrate/v1/{realm}/socket` is documented as if public in docs/api/astrate_native_api.yaml — the spec has no `components.securitySchemes` at all and root `security: []` (yaml:21), and the operation declares no `security`, so a client generated from this spec sends no `Authorization` header and is answered 401 `{"errors":{"detail":"Unauthorized"}}` by `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112, bearer token read only from the Authorization header at middleware.go:131, 178-196). Add a scheme mirroring `a_aea` (docs/api/astarte_appengine_api.yaml:1436-1443: apiKey in header `Authorization`, described as a realm JWT carrying the `a_ch` claim) and set `security` on this operation only — root `security: []` stays so health/readiness/metrics/version remain documented as unauthenticated. If you model the Phoenix twin too, its credential is `?token=` in the query (internal/appengine/channels/ws.go:53-54), i.e. apiKey `in: query`, not header. Description notes worth carrying: the scheme is case-insensitive `Bearer` with an optional colon, and an unknown realm is 401 not 404 (middleware.go:68-72). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
3676ad9 2026-10-06T22:11:12+02:00 mule: log docs-native-socket-sse-exact-accept
1322214 2026-10-06T22:11:12+02:00 mule: docs-native-socket-sse-exact-accept [auto]: fix the SSE negotiation prose on `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml — the description (yaml:289-291), the `transport` parameter (yaml:328-334) and the `200` response (yaml:338-339) read as if any client asking for `text/event-stream` gets SSE, but the predicate `wantsSSE` is exact string equality (internal/appengine/stream/ws.go:151-153): measured, `Accept: text/event-stream, text/plain` and `Accept: text/event-stream;q=1.0` both return false and fall through to `websocket.Accept`, where a non-upgrade client gets the 426 above; `transport` selects SSE only for the literal `sse`, and `transport=websocket` / `transport=` (empty) both mean WebSocket. Say "exactly" for the header, say any other `Accept` falls through to the upgrade attempt, and do NOT claim the server rejects out-of-enum `transport` values — the enum is a client-side contract the handler never enforces. While there, give the `200` its real media type: the handler sets `Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive` and writes one `data: {...}` frame per event followed by a blank line (ws.go:123-126, 156-162). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
fc702e3 2026-10-06T20:57:40+02:00 mule: blocked docs-native-socket-handshake-errors
e3f5780 2026-10-06T20:40:46+02:00 mule: refill the queue
3fe292e 2026-10-06T20:40:46+02:00 mule: recipe docs-sync ran (816s)
8b1b7f2 2026-10-06T20:25:10+02:00 mule: recipe code-review timed out
d567596 2026-10-06T19:55:47+02:00 mule: log compat-note-v135
3a8ea6d 2026-10-06T19:55:47+02:00 mule: compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (newest stable, 2026-10-05, maintenance-only: mix.exs/mix.lock version bumps plus a single CHANGELOG "astarte_data_updater_plant RPC server reliability" entry, no wire/API surface change; v1.4.0 is still rc.6-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2/v1.3.4 wording proposal there already stands and its version-reference delta was recorded on 2026-10-05 as reading "newest stable v1.3.5", so fold this in as that one-word delta rather than re-deriving the proposal.
bec0dfb 2026-10-06T19:49:31+02:00 mule: blocked probe-encrypted-mapping-enforcement
ec1b11e 2026-10-06T19:43:53+02:00 mule: refill the queue
c7db5c0 2026-10-06T19:43:53+02:00 mule: recipe astarte-upstream ran (1064s)
59a31a2 2026-10-06T18:56:12+02:00 mule: refill the queue
9945b5c 2026-10-06T18:56:12+02:00 mule: recipe github-issues ran (543s)
4a91f16 2026-10-06T16:11:29+02:00 mule: refill the queue
71806b5 2026-10-06T16:11:29+02:00 mule: recipe milestones ran (327s)
cbdd3e8 2026-10-06T16:05:39+02:00 mule: recipe hygiene timed out
230dd11 2026-10-05T23:04:09+02:00 mule: recipe docs-sync timed out
02003c1 2026-10-05T22:19:14+02:00 mule: recipe code-review ran (125s)
1e8d87f 2026-10-05T22:16:18+02:00 mule: refill the queue
28ffb8d 2026-10-05T22:16:18+02:00 mule: recipe astarte-upstream ran (187s)
732cc7c 2026-10-05T21:54:40+02:00 mule: refill the queue
5c05c31 2026-10-05T21:54:40+02:00 mule: recipe github-issues ran (325s)
5fe2775 2026-10-05T21:34:52+02:00 mule: refill the queue
6ec6c12 2026-10-05T21:34:52+02:00 mule: recipe milestones ran (210s)
c03cfe4 2026-10-05T21:23:13+02:00 mule: recipe hygiene timed out

$ grep -n "check_pulse\|pulse\|dead-man\|alarm" .mule/todo.md | head; echo "=== for-giulio check_pulse entries ==="; grep -n "check_pulse" .mule/for-giulio.md | head
=== for-giulio check_pulse entries ===
21:- **github-issues triage run, 2026-10-06: nothing proposable for the 24th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented or closed on GitHub.** The recipe's `--limit 40` command prints all **30** open issues: **26** `mule-alarm` **#94–#119** plus the same four non-alarm, none of them proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (`milestone-4.0`, already escalated), and one the standing instruction excludes from triage entirely. New since the 2026-10-05 run: **#119**, created 2026-10-06T10:58:45Z, "nothing has landed in 15h" — a real window, not an artefact: the land before it was `230dd11` at 2026-10-05T21:04:09+02:00 (19:04Z), 15.9h earlier, and it was ended by today's first land `cbdd3e8` at 16:05:39+02:00 (14:05Z), three hours after the alarm. So it belongs to the genuine 14–16h band the 2026-10-05 entry measured, i.e. overnight silence between evening and afternoon batches. **Proposal, extending the 2026-10-05 one: close #94–#118, all twenty-five** — each is a one-day alarm superseded by the next, and nothing reads issue state; leave **#119** (today's, real window, already ended) to self-expire. The `check_pulse`-before-`beat` ordering fix (move the call at `tools/mule.sh:633` to after the tick's `beat`) is still escalated above, still yours, and still unactioned — today's alarm is not evidence for it either way.
37:  **My recount reproduces yesterday's numbers rather than replacing them, which is the useful part: the distribution is bimodal with an empty middle.** Over all **139** `mule: log` lands since 2026-08-20 (138 gaps): median **0.8h**, and the counts `> 8h` = `> 12h` = `> 13h` = **27** — so no gap anywhere in the record falls in (8h, 13h). Each of the 25 open alarms therefore measured a real 14.1–16.2h window (the only outlier is **#114** at 8.6h, filed 2026-09-27T20:09Z off the 11:35Z land that same day, which is the bimodality showing its other edge: a 20:00Z check before the evening batch starts). The 2026-10-04 entry's conclusion holds and hardens — a threshold in [8h, 13h) fires on every healthy day by construction; the real stall in the record (2026-09-28 → 10-01) is **74.7h**. **The proposed fix is unchanged and remains yours: move the `check_pulse` call from `tools/mule.sh:633` to after the tick's `beat`, so the switch cannot file against the work it is about to watch land.** Today's instance is the cleanest yet: **#118 filed 10:54:06Z** measuring a genuine 14.8h back to the 2026-10-04T20:06Z land, and the next land, `store-validatepipelinegraph-error-branches`, committed **10:59:49Z — 5 minutes later** (#117 was 23). Six lands today (10:59Z → 19:34Z), five `done` in `.mule/log.md`.
59:  **On the alarms themselves: they are not "false", and I should correct that framing.** The 2026-10-03 entry called #115/#116 "provably false" on the grounds that work had landed and been pushed on all three days. The gaps were real; what was false is the alarm's *conclusion*, not its arithmetic. Measured from the commit dates of every `[auto]` land since 2026-09-05 (115 lands, 114 gaps): **median gap 0.8h, and exactly 24 gaps exceed the 8h threshold — one per open alarm issue, #94–#117, a 1:1 correspondence.** The distribution is bimodal with nothing in the middle: the **smallest** gap above 8h is **14.6h** and the largest sub-8h gap is far below it. The mechanism is now exact rather than inferred: `check_pulse` is called at **`tools/mule.sh:633`, before the tick attempts any task**, and today's first cron entry is `17 13 4 10` = **13:17 CEST = 11:17Z** — the minute #117 was filed (11:17:06Z). **The switch fires on the very tick that is about to do the work which would have silenced it**, and the `beat` at the end of that tick (mule.sh:737) then deletes `.alarmed`. #117 is the clean case: it fired at 11:17:06Z measuring a genuine 15h window back to `f84bbb5` (2026-10-03T19:50:23Z), and the next land, `e728217`, committed at **11:40:11Z — 23 minutes later**. The four-day stall the 2026-10-01 entry diagnosed is over: five task attempts today, three `done` (`e728217`, `5044143`, `def0af6`), two `blocked`.
60:  **This should let you retire the threshold question rather than answer it.** The 2026-09-26 entry put two options to you (raise `MULE_IDLE_ALARM_HOURS` to ~24, or pulse `.mule/log.md` for a whole day with no `done` row). The gap histogram now says a threshold **cannot** work: healthy days reach **23.1h** (09-24→09-25) and **20.5h** (09-17→09-18) with no stall, while the three real stalls were **44.9h, 49.3h and 74.7h**. The only separating cut is ~36h, which would tolerate three days of silence before speaking — so the option worth taking is the second one, or the cheaper structural fix: **move the `check_pulse` call from mule.sh:633 to the end of the tick**, after the task has landed, so a healthy first-tick-of-the-day can never file against itself. That is a one-line change plus a shell test, and it is the only one of the three that does not trade a false positive for a false negative. It stays yours because it changes the mule's own safety behaviour, and because it lives in the same branch that needs the merge above — say the word and it becomes a queue line.
80:  **The 2026-10-01 "latched off, and the silence is real" reading is now superseded: the latch cleared itself the moment work landed, and the false alarm resumed on schedule.** That entry found `.alarmed` present with no land since 09-28T18:01Z, and that was true — the queue had genuinely stalled, on twelve `[legion]` lines with the Legion Go off. The 10-01 land beat, `beat` deleted `.alarmed` (mule.sh:737), `check_pulse` stopped short-circuiting (mule.sh:748), and the switch filed again the next morning. `.alarmed` is absent right now and the heartbeat reads 2026-10-03T17:05:22Z. **So the latch was a symptom of the stall, not a fault in the switch** — I withdraw the implicit "the switch cannot report a real outage" worry from that entry: nothing in `mule.sh` is broken, the 8h default (mule.sh:740) is simply wrong for a cadence that lands in the evening and polls in the morning, where 14h > 8h every day by construction.
122:  **The alarm stopped firing because of a bug in the switch, not because the mule recovered.** Last heartbeat **2026-09-28T18:01Z** (`.mule/.heartbeat`), and the last row in `.mule/log.md` is 2026-09-28 (`docs-sync-pairing-unregister-description`, `8ffb39e`) — no task attempt in the ~74h since, while recipe ticks are still committing (today 19:46–21:41Z), so the timer is alive and only the task queue is silent. On **2026-09-29T10:49Z** the pulse did fire (16h) but took the no-`gh` branch: the last line of this file is the fallback text "The mule has been idle 16h" (`mule.sh:771-773`) and **no issue was filed** — #114 is the last alarm ever filed. That same call created `.mule/.alarmed`, which only `beat` clears (`mule.sh:737`), and `check_pulse` returns immediately when it exists (`mule.sh:748`), so the switch has been **latched off for three days** and will file nothing until some tick lands work. The longest silence in the log is the one the switch reports least. That also gives the 2026-09-25 "the body is wrong" note a consequence at last: "Close it — the mule reopens a new one if the silence continues" (`mule.sh:764`) is false by construction — one alarm per silence, never a second.
136:  **What the alarm really measures, and why it fires daily by construction.** `check_pulse` (mule.sh:739-774) compares `.heartbeat` against `MULE_IDLE_ALARM_HOURS` (default 8, mule.sh:740), and `beat` is called only on a `done` land or a `checked` run (mule.sh:578, 540) — a refill does not beat. Last beat: **2026-09-25T20:19:13Z**; the pulse that files the issue runs ~10:54–11:12Z. 20:19Z → 11:04Z is ~14.75h, comfortably over 8h, so the arithmetic is right and the conclusion is wrong: with lands bunched into an evening window, the gap to the next morning's tick is >8h *every day*. It is a mis-set threshold, not a dead mule. The two days in this window that landed nothing (09-19, 09-22) are the only ones where the alarm had anything real to say.

$ cat .mule/.heartbeat 2>/dev/null | xargs -I{} date -u -d @{} '+last beat: %Y-%m-%dT%H:%M:%SZ' 2>/dev/null; echo "--- raw: $(cat .mule/.heartbeat 2>/dev/null)"; ls -la .mule/.alarmed 2>/dev/null || echo "(no .alarmed)"; echo "=== log.md tail ==="; tail -5 .mule/log.md
last beat: 2026-10-06T20:29:30Z
--- raw: 1791318570
-rw-r--r-- 1 root root 0 Oct  7 13:04 .mule/.alarmed
=== log.md tail ===
| 2026-10-06 | compat-note-v135 | done | 190s | 3a8ea6d |
| 2026-10-06 | docs-native-socket-handshake-errors | blocked | 91s | wrote nothing |
| 2026-10-06 | docs-native-socket-sse-exact-accept | done | 410s | 1322214 |
| 2026-10-06 | docs-native-socket-security-scheme | done | 299s | 87c18cf |
| 2026-10-06 | docs-native-socket-event-schema | blocked | 662s | opencode exited 1 |

$ sed -n '570,582p;735,775p' tools/mule.sh
"The mule pushed \`$(git -C "$REPO" rev-parse --short HEAD)\` to \`$MULE_BRANCH\` for this.

Unreviewed: the gates passed — it compiles, the tests pass, a new test was shown to fail without the change, and no frozen file was touched. That is not the same as the change being worth having. Nothing merges to \`main\` until someone reads it:

    bash tools/mule.sh review

Left open on purpose. Whether this actually resolves the issue is your call, not mine.";;
    esac
    beat
    ok "landed: $text (${elapsed}s)"
    return 0
  fi

# So: any tick that lands work touches the heartbeat. A tick that finds it stale says so
# somewhere Giulio will actually see, exactly once per silence.
beat() { date +%s > "$MULE/.heartbeat"; rm -f "$MULE/.alarmed"; }

check_pulse() {
  local max="${MULE_IDLE_ALARM_HOURS:-8}" now last age
  [ "$max" = 0 ] && return 0
  now="$(date +%s)"
  if [ ! -f "$MULE/.heartbeat" ]; then beat; return 0; fi
  last="$(cat "$MULE/.heartbeat" 2>/dev/null)"; [ -n "$last" ] || { beat; return 0; }
  age=$(( (now - last) / 3600 ))
  [ "$age" -lt "$max" ] && return 0
  # One alarm per silence, not one per tick — the alarm is only useful if it is rare.
  [ -f "$MULE/.alarmed" ] && return 0
  : > "$MULE/.alarmed"
  bad "no work has landed in ${age}h — raising the alarm"
  local body="The mule has not landed any work in ${age} hours (threshold ${max}h).

It is probably still ticking and exiting 0 — that is how this failure looks. Worth checking,
on the Pi:

    systemctl list-timers mule.timer
    journalctl -u mule.service --since '-1 day' | tail -50
    cd /root/astrate-mule && bash tools/mule.sh status

Common causes, all of which have happened: the provider is refusing every run, the queue has
nothing runnable in it, the working tree is dirty so every tick aborts, or push has been
failing and the work is only on the Pi.

Filed automatically. Close it — the mule reopens a new one if the silence continues."
  if command -v gh >/dev/null && [ -n "${MULE_PUSH:-}" ]; then
    gh issue create --title "mule: nothing has landed in ${age}h" --body "$body" \
      --label "$MULE_ALARM_LABEL" >/dev/null 2>&1 \
      && note "filed an idle alarm on GitHub" \
      && return 0
  fi
  # No gh, or it failed: the escalation file is the fallback and it is committed anyway.
  printf -- '- **The mule has been idle %sh.** %s\n' "$age" \
    "Filed by the dead-man's switch; see journalctl on the Pi." >> "$MULE/for-giulio.md"
}

$ sed -n '144,152p' internal/broker/aclhook.go; echo "=== test row ==="; grep -n "keyAgreement" internal/broker/aclhook_test.go; echo "=== upstream tags ==="; gh api repos/astarte-platform/astarte/tags --jq '.[].name' | head -6
e = &offlineEntry{ownership: map[string]interfaceschema.Ownership{}, loadedAt: now}
	o.entries[cn] = e
	o.mu.Unlock()

	ownership := o.load(cn)

	o.mu.Lock()
	e.ownership = ownership
	o.mu.Unlock()
=== test row ===
56:		{"pub control keyAgreement", base + "/control/keyAgreement", true, false},
=== upstream tags ===
v1.4.0-rc.6
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1

$ gh issue list --state open --limit 100 --json number --template '{{range .}}{{.number}} {{end}}'; echo; echo "count: $(gh issue list --state open --limit 100 --json number --jq 'length')"
120 119 118 117 116 115 114 113 112 111 110 109 108 107 106 105 104 103 102 101 100 99 98 97 96 95 94 93 92 78 1 
count: 31

$ grep -n "keyAgreement\|082-key\|rc\.\|newest stable" internal/broker/aclhook.go
72:		// control/keyAgreement is deliberately denied: full key agreement is
74:		// spec (082-key_agreement_protocol.md) only ships with v1.4.0-rc.5,
75:		// v1.3.3 being the newest stable tag. It stays denied until the
181:	dev, err := o.st.GetDevice(ctx, rc.id, identity.DeviceID)
186:	return loadOwnership(ctx, o.st, rc.id, dev.Introspection, o.log)

→ Read .mule/for-giulio.md [limit=24]

$ date -u '+now: %Y-%m-%dT%H:%M:%SZ'; echo "heartbeat age: $(( ( $(date +%s) - $(cat .mule/.heartbeat) ) / 3600 ))h"; echo "=== log.md rows for 2026-10-07 ==="; grep -c "2026-10-07" .mule/log.md || echo 0; echo "=== divergence ==="; git rev-list --left-right --count origin/main...HEAD; echo "=== #93 ancestry ==="; git merge-base --is-ancestor 24ad5b8 origin/main && echo "24ad5b8 ON main" || echo "24ad5b8 NOT on main"
now: 2026-10-07T16:55:20Z
heartbeat age: 20h
=== log.md rows for 2026-10-07 ===
0
0
=== divergence ===
4	697
=== #93 ancestry ===
24ad5b8 NOT on main

$ grep -n "\bbeat\b" tools/mule.sh
322:  # cmd_refill once sliced beat() and check_pulse() out of the file and left their call sites
327:  for fn in beat check_pulse first_open code_id line_slug proof_gate gates issue_tasks \
540:            beat
578:    beat
737:beat() { date +%s > "$MULE/.heartbeat"; rm -f "$MULE/.alarmed"; }
743:  if [ ! -f "$MULE/.heartbeat" ]; then beat; return 0; fi
744:  last="$(cat "$MULE/.heartbeat" 2>/dev/null)"; [ -n "$last" ] || { beat; return 0; }

$ sed -n '532,545p' tools/mule.sh
echo; cat "$outlog"
        } > "$MULE/reports/$slug.md"
        log_row "$slug" "checked" "$elapsed" "$(git -C "$REPO" rev-parse --short HEAD)"
        save_report "$slug" "checked" "$elapsed" "$outlog"
        git -C "$REPO" add "$MULE/reports/$slug.md" "$MULE/reports/log" "$LOG" >/dev/null 2>&1
        git -C "$REPO" commit -q -m "mule: $slug passed on $(git -C "$REPO" rev-parse --short HEAD)"
        [ -n "${MULE_PUSH:-}" ] && { git -C "$REPO" push -q origin "$MULE_BRANCH" 2>/dev/null \
          && ok "pushed $MULE_BRANCH" || note "push failed — the work is safe locally"; }
            beat
    ok "checked: $slug clean (${elapsed}s) — stays queued for the next change"
        return 0
        ;;
    esac
    verdict=blocked; reason="wrote nothing"

$ grep -n "check_pulse$\|check_pulse " tools/mule.sh; echo "---"; sed -n '746p' tools/mule.sh; sed -n '748p' tools/mule.sh; sed -n '633p' tools/mule.sh
327:  for fn in beat check_pulse first_open code_id line_slug proof_gate gates issue_tasks \
633:  check_pulse
---
  [ "$age" -lt "$max" ] && return 0
  [ -f "$MULE/.alarmed" ] && return 0
  check_pulse

Zero proposable issues (25th consecutive run); writing the evidence entry to `.mule/for-giulio.md`.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,14 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **github-issues triage run, 2026-10-07: nothing proposable for the 25th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub.** The recipe's `--limit 40` command prints all **31** open issues (authoritative `--limit 100` agrees): **27** `mule-alarm` **#94–#120** plus the same four non-alarm, none of them proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (`milestone-4.0`, parked until the v3.0 queue clears), and one the standing instruction excludes from triage entirely. New since the 2026-10-06 run: **#120**, created 2026-10-07T11:03:45Z, "nothing has landed in 14h" — arithmetic pinned, not inferred: `.mule/.heartbeat` holds `1791318570` = **2026-10-06T20:29:30Z**, written 3s after the `docs-native-socket-security-scheme` done land (`87c18cf`, 20:29:27Z), and 11:03:45Z − 20:29:30Z = **14h34m**, so `age=$(( (now-last)/3600 ))` prints 14 — the title verbatim; `.mule/.alarmed` carries mtime 11:04Z, i.e. set by that filing. Genuine window, the same 14–16h overnight band #117–#119 measured. **Proposal, extending the 2026-10-06 one: close #94–#119, all twenty-six** — each is a one-day alarm superseded by the next and nothing reads issue state; leave **#120** (today's, real window) to self-expire.
+  **New this run, and it is the 2026-10-01 pattern recurring: the switch is latched off through a silence that is real.** At 2026-10-07T16:55Z `.mule/.heartbeat` still read 2026-10-06T20:29:30Z — **20h**, with zero `2026-10-07` rows in `.mule/log.md` — but `.mule/.alarmed` (set by #120) short-circuits `check_pulse` at `tools/mule.sh:748`, so nothing can file again until a `done` land (mule.sh:578) or a `checked` run (mule.sh:540) calls `beat`, which is the only thing that clears the latch (mule.sh:737). Recipe ticks are committing today (`becf3cd` 11:25:35Z, `dde64e0` 11:38:44Z) but a recipe verdict is neither of those two paths, so they neither beat nor un-latch — measured today, not assumed. The `check_pulse`-before-`beat` ordering fix (move the call at `tools/mule.sh:633` to after the tick's `beat`) remains escalated above, still yours, still unactioned; #120 is the same arithmetic as #119, not fresh evidence for it either way.
+  **#93 at 33 days.** `git rev-list --left-right --count origin/main...HEAD` = **4 / 697** today (4 / 679 on 10-06), and `git merge-base --is-ancestor 24ad5b8 origin/main` still fails — the rewrite is entirely on `mule/queue` / `origin/mule/queue`. The version-drift detail recorded on 2026-10-06 (the comment says `v1.4.0-rc.5` / "v1.3.3 being the newest stable tag"; upstream is `v1.4.0-rc.6` / `v1.3.5`, re-checked via `gh api repos/astarte-platform/astarte/tags` today) folds into that same review rather than becoming a line of its own.
+
+---
+
 - **COMPATIBILITY.md wording proposal for upstream v1.3.5 (newest stable, 2026-10-05; v1.4.0 is still rc.6-only) — a one-word delta on the open v1.3.2/v1.3.4 proposal below, nothing re-derived.** v1.3.5 is maintenance-only: `gh api repos/astarte-platform/astarte/releases` shows it published `2026-10-05T14:11:44Z`, not a prerelease, with a body of a single `### Fixed` line — "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available" — and `compare/v1.3.4...v1.3.5` is 4 commits / 35 files whose only non-`mix.exs`/`mix.lock` files are `.github/workflows/astarte-build-workflow.yaml`, `CHANGELOG.md`, `README.md`, `astarte-dashboard/package.json` and `docker-compose.yml`. Not one source file changed, so no route, MQTT topic, control message or interface-schema field moved (Astrate implements no `data_updater_plant` and no Erlang RPC server). **The only change a v1.3.5-aware doc carries: the v1.3.4 entry's §Infrastructure-differences sentence — "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target" — reads "…(**newest stable v1.3.5**) as the target".** Everything else stands unchanged: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"); that entry's "v1.4.0 is still rc.5-only" header is already superseded to rc.6 by the open rc.6 delta entry below, not here. This is the same one-word delta recorded on 2026-10-05 in this file ("newest stable v1.3.5"), folded in as asked rather than re-derived. Applied to nothing: `docs/COMPATIBILITY.md` untouched and still 1.2.2-targeted (it contains no `1.3` reference at all), `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). Raw: [v1.3.5](https://github.com/astarte-platform/astarte/releases/tag/v1.3.5).
 
 ---
 

$ sed -n '9,22p' .mule/for-giulio.md | cut -c1-140; echo "=== git status (read-only check) ==="; git status --porcelain | head; echo "=== todo.md issue- lines ==="; grep -c "issue-9[0-9]\|issue-1[01][0-9]\|issue-78" .mule/todo.md
line once you have dealt with it — this file is a queue, not a log.

---

- **github-issues triage run, 2026-10-07: nothing proposable for the 25th run — no `.mule/todo.md` lines, no `gh issue create`, nothing co
  **New this run, and it is the 2026-10-01 pattern recurring: the switch is latched off through a silence that is real.** At 2026-10-07T16:5
  **#93 at 33 days.** `git rev-list --left-right --count origin/main...HEAD` = **4 / 697** today (4 / 679 on 10-06), and `git merge-base --i

---

- **COMPATIBILITY.md wording proposal for upstream v1.3.5 (newest stable, 2026-10-05; v1.4.0 is still rc.6-only) — a one-word delta on the

---

=== git status (read-only check) ===
 M .mule/for-giulio.md
=== todo.md issue- lines ===
0

Recipe carried out. Zero issues qualified for a task line — 25th consecutive run.

```
Done: github-issues triage run 2026-10-07 — 31 open issues surveyed, none machine-checkable, so no todo.md lines; evidence appended to .mule/for-giulio.md
Files: .mule/for-giulio.md (one new top entry, 3 paragraphs)
Verified: gh issue list --limit 40 and --limit 100 -> both print 31 (27 mule-alarm #94–#120 + #93/#92/#78/#1) pass; .mule/.heartbeat=1791318570 = 2026-10-06T20:29:30Z vs #120 created 11:03:45Z = 14h34m, matches title pass; git rev-list --left-right --count origin/main...HEAD = 4/697, git merge-base --is-ancestor 24ad5b8 origin/main = false pass; grep issue- .mule/todo.md -> 0 lines (untouched) pass. No Go file changed, so the go test/gofmt gate does not apply and was not run.
Unsure: whether #120's window will still read as "real" once a done land clears the latch — I measured the current state (heartbeat 20h stale, .alarmed set at 11:04Z, zero log.md rows for 2026-10-07) but cannot see the next beat
Follow-ups: none
```

Dispositions, unchanged from the 2026-10-06 entry: **#93** implemented on `mule/queue`, 33 days awaiting `bash tools/mule.sh review`; **#92** parked on a stable v1.4.0; **#78** parked until the v3.0 queue clears; **#1** untouched per standing instruction. Close proposal extended to **#94–#119**; **#120** left to self-expire. No `gh issue create` — the recipe does not ask for it — and no GitHub state changed.
