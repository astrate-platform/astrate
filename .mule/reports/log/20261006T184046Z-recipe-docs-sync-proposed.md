slug: recipe-docs-sync
verdict: proposed
at:  8b1b7f2
ran: 2026-10-06T18:40:46Z on DietPi in 816s


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

$ cat .mule/todo.md 2>/dev/null; echo "---FORGIULIO---"; cat .mule/for-giulio.md 2>/dev/null | tail -40
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_11278bf62001GwQICfESDZliUG

- [x] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` path parameter (docs/api/astarte_pairing_api.yaml:360) and `RegisterRequest.hw_id` (yaml:373) — but `deviceid.Parse` accepts only the 22-character **unpadded base64url** form, in strict mode, and rejects the standard alphabet's `+` and `/` and any `=` padding with an error wrapping `ErrInvalid` (pkg/deviceid/deviceid.go:26, 37, 39-60; upstream parity cited at 36). Measured: with the same `encoding` the parser uses, `dT6hS2W9TT6LEnP25ks+lg` and `.../lg` fail at byte 19 and the padded 23-character form fails at byte 21, so a client generated from the current spec can send a well-formed standard-base64 ID and get 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}` (internal/pairing/http.go:307-309). State the exact encoding and the 22-character length in both descriptions and add a `pattern`/`minLength`/`maxLength` to `hw_id`. The existing `hw_id: dT6hS2W9TT6LEnP25ks_lg` example (yaml:123) is valid and stays. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in docs/api/astarte_pairing_api.yaml (yaml:374-376), an unconstrained string described only as "Optional Astrate extension specifying the device's initial payload format" — so a generated client can send any string, while `Service.Register` accepts only `"bson"` or `"json"` and rejects everything else with 422 `{"errors":{"initial_payload_format":["is invalid"]}}` (internal/pairing/service.go:179-181 -> internal/pairing/http.go:310-312; the allowed set is stated at service.go:173). Add `enum: [bson, json]` plus a note that absent/empty means unset, and add the three missing `ValidationError` examples while you are in the component: `initial_payload_format: ["is invalid"]` (http.go:310-312), `csr: ["can't be blank"]` (http.go:173-177) and `client_crt: ["can't be blank"]` (http.go:253-257) — the component currently shows only the four hw_id/error_name rows (yaml:594-618) and the `csr`/`client_crt` blank-field cases are what a device actually hits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` description in docs/api/astarte_pairing_api.yaml (yaml:153), "Removes a device from the realm. Returns 204 on success.", which reads as a device-and-data deletion. The row and every byte of its stored data survive: `store.UnregisterDevice` only clears the credential trail (`credentials_secret_hash = ''`, `first_credentials_request = NULL`, `cert_serial = NULL`, `cert_aki = NULL`) and flips `status` back to `'registered'` (internal/store/devices.go:93-110), and the service docstring says so outright ("its data is kept", internal/pairing/service.go:227-229). State the real semantics: the device becomes registrable again, its interfaces, datastream data and group memberships are retained, and a second DELETE is a 204 no-op (the row is still there, so RowsAffected is 1, devices.go:106-108) — not the documented 404. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `base64.RawURLEncoding.Strict()` (pkg/deviceid/deviceid.go:37), so of the 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with `ErrInvalid` while `..._lA` is accepted, and the pairing spec's `hw_id` pattern `^[A-Za-z0-9_-]{22}$` (added by docs-sync-pairing-deviceid-base64url) deliberately does not encode the rule because a pattern is a client-side contract and upstream's leniency is unmeasured. The docstring justifies strictness as byte-identical round-tripping, which is a reason, not a parity claim, and docs/COMPATIBILITY.md records no entry for it. Measure upstream (Elixir `Base.url_decode64!(padding: false)`, which does not check those bits) on the Legion Go, then either record the divergence in docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.
- [x] fdo-rc6-scope-delta-for-giulio [auto]: upstream v1.4.0-rc.6 (2026-09-30, 27 commits, 154 files, almost all FDO) grows the onboarding surface issue #78 is scoped against, and none of it is in #78's body (written against upstream 1.3 / rc.5): new `POST /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}/to0` re-runs TO0 and returns the refreshed rendezvous expiry (`{"data":{"expiry":…}}`, ownership_voucher.ex:63-73; 409 `device_already_onboarded` via fallback_controller.ex:120-125); voucher deletion, which revokes the rendezvous registration first and *keeps* the voucher if that fails (`:rendezvous_revocation_failed` → 500) and refuses another realm's voucher through `ensure_voucher_in_realm` (ownership_voucher.ex:40-60 — now expressible because the voucher row carries its realm); `device_id` accepted on voucher upload; the device registered immediately on upload rather than mid-TO2 (`add_unconfirmed_credentials`, device.ex:63-85); single-entry x5chain certs accepted; invalid vouchers no longer stored; vouchers and to2-sessions moved from the per-realm keyspace to the global `astarte` one with `fdo_guid` added to the device row (housekeeping/realms/queries.ex:784, 1129-1133); and RM device deletion now cascading to the voucher (device_removal/core.ex). Write the delta into `.mule/for-giulio.md` as a scope update for #78 and say which of these a client could actually observe — do not edit or comment on the issue (GitHub state is read-only for the mule) and do not implement anything.
- [x] compat-note-v1.4-rc.6 [auto]: record in `.mule/for-giulio.md` that v1.4.0-rc.6 is wire-inert for Astrate, as a delta on the open v1.3.2/v1.3.4 wording proposal rather than a re-derivation. Measured, the whole non-FDO delta is inert: DUP RPC-server reliability is OTP supervisor/GC tuning; the Dashboard's new `unknown_status` is computed client-side from `last_connection && !last_disconnection && !isConnected` (models/Device/index.ts:213-215) and both timestamps are already on Astrate's device shape (internal/appengine/service.go:118), so a device whose disconnection was never recorded reads as "Unknown status" rather than "Disconnected"; the `on_device_registered` dispatch moved from `Agent.register_device` to `Engine.register_device` (pairing/engine.ex:157-165, 186-193) so it also covers the FDO path, but the Agent path is byte-for-byte the same rule and Astrate already fires it on every successful registration (internal/pairing/service.go:221-223); the `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"); and `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma. So the proposal's "v1.4.0 is still rc.5-only" reference becomes rc.6, and its "FDO authentication (pairing, disabled by default)" line understates the 1.4 line, which ships a full ownership-voucher surface Astrate does not emulate (#78, parked). Do not edit docs/COMPATIBILITY.md and do not restate the standing proposal.

- [x] container-stop-deadline [auto]: make `cliInstance.Stop` (internal/flow/blocks/container/docker.go:167-175) derive its 15s bound whenever the passed context carries **no deadline**, not only when it is nil — both cleanup callers pass `context.Background()` (docker.go:122, the "Best-effort cleanup so we do not leave orphans on mapping failure" path, and block.go:121, the not-ready cleanup in `New`), so a wedged docker daemon leaves an unbounded `docker rm -f` inside `exec.CommandContext` holding flow instantiation, while `Block.Stop` (block.go:310-311) does the same job correctly with an explicit 15s timeout. Measured with the package's own `Run` hook: the `rm` invocation arrives with `ctx.Deadline()` unset on both paths. Add a case asserting `inst.Stop(context.Background())` reaches the injected `Run` with a deadline set (it does not today) and that the hostPort-failure path still issues `rm -f <id>` (docker.go:120-124, never executed by any test today). Container-free.
- [!] container-timeout-bounds [auto]: bound `timeout_ms` and `ready_timeout_ms` in `parseConfig` (internal/flow/blocks/container/block.go:205-219) — both are only checked `> 0` and then multiplied unchecked (`time.Duration(n) * time.Millisecond`), which wraps: measured `ready_timeout_ms: 18446744073709` gives `-551.616µs`, so `New` fails with "container: not ready after wait" without ever issuing a probe (the wrapped negative deadline fails the loop guard at httpbridge.go:64), while `ready_timeout_ms: 9223372036854` gives 2562047h47m16s and `New` never returns (still polling `/healthz` after 3s in the probe). The two consumers of a non-positive duration also disagree — `Bridge.timeout()` (httpbridge.go:39-44) treats it as "use the 5s default" while net/http treats a non-positive `Client.Timeout` as no deadline at all (measured: a client with `Timeout: -551µs` returned err=nil after a 3s handler). Reject a value that would overflow or exceed a documented maximum, naming the key and the value in the error, and add table cases asserting `New` refuses both (it accepts them today). Nothing bounds these on the wire either: internal/flow/blocks/schema.go:19 advertises both as bare `"type":"integer"`. Container-free. — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [x] container-parseconfig-rules-test [auto]: pin the four `parseConfig` rules that no test touches (internal/flow/blocks/container/block.go:188-219 — all `count=0` in `go test -coverprofile`, which is 71.2% for the package): `config` must be a JSON object (192-193), `port` an integer in 1–65535 including the 0 / 65536 / non-integer rejections (197-202), and `timeout_ms` and `ready_timeout_ms` positive (205-210, 215-217) — beside the existing `TestConstructor_RequiresImage` (block_test.go:99), which is the only config case the suite has. The port rule is the one with teeth: it is what keeps `127.0.0.1::99999` out of the `docker run` argv (docker.go:93). Container-free, pure.
- [!] container-response-cap-test [auto]: pin `Bridge.RoundTrip`'s 1 MiB response cap (internal/flow/blocks/container/httpbridge.go:133-140 — `io.LimitReader(resp.Body, maxBody+1)` plus the explicit rejection; lines 136, 138.16 and 138.40 are all uncovered), the only bound between a tenant-supplied container and process memory on the message path: assert a handler answering exactly `maxBody+1` bytes is rejected with the "exceeds" error and that exactly `maxBody` bytes still round-trip, driving a `Bridge{MaxBodyBytes: N}` so the fixture stays small. Also worth one line in the same change: the cap is absent from the operator-facing catalog entry (internal/flow/blocks/info.go:116) and from the schema, so an author whose container emits 2 MiB learns about it only from a runtime error. Container-free. — BLOCKED: lint failed: internal/flow/blocks/container/block_test.go:500:2: redefines-builtin-id: redefinition of the built-in function cap (rev
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
- [x] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.go
- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
- [x] compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (newest stable, 2026-10-05, maintenance-only: mix.exs/mix.lock version bumps plus a single CHANGELOG "astarte_data_updater_plant RPC server reliability" entry, no wire/API surface change; v1.4.0 is still rc.6-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2/v1.3.4 wording proposal there already stands and its version-reference delta was recorded on 2026-10-05 as reading "newest stable v1.3.5", so fold this in as that one-word delta rather than re-deriving the proposal.
---FORGIULIO---
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

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*native*' ; echo "=== SPEC PATHS ==="; rg -n '^\s{2}/' docs/api/astrate_native_api.yaml
=== SPEC PATHS ===
26:  /astrate/v1/health:
44:  /astrate/v1/readiness:
77:  /astrate/v1/metrics:
114:  /appengine/health:
142:  /realmmanagement/health:
170:  /pairing/health:
200:  /appengine/version:
220:  /realmmanagement/version:
240:  /pairing/version:
260:  /housekeeping/version:
284:  /astrate/v1/{realm}/socket:
349:  /appengine/v1/socket/websocket:

$ wc -l docs/api/*.yaml; echo ---; ls internal/swagger/
1894 docs/api/astarte_appengine_api.yaml
   510 docs/api/astarte_housekeeping_api.yaml
   686 docs/api/astarte_pairing_api.yaml
  1390 docs/api/astarte_realm_management_api.yaml
   509 docs/api/astrate_native_api.yaml
  4989 total
---
swagger.go
swagger_test.go

$ rg -n '"/astrate/v1|/appengine/health|/housekeeping/version|/appengine/v1/socket' cmd/ internal/ | head -50
cmd/astrate/main.go:551:	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+addr+"/astrate/v1/readiness", nil) //nolint:gosec // G704: self-probe of the local readiness endpoint
cmd/astrate/main_test.go:121:		if code, _ := get(t, base+"/astrate/v1/health"); code != http.StatusOK {
cmd/astrate/main_test.go:130:		code, body := get(t, base+"/astrate/v1/metrics")
cmd/astrate/main_test.go:204:		if code, _ := get(t, base+"/astrate/v1/readiness"); code == http.StatusOK {
internal/appengine/stream/ws_test.go:66:		req := httptest.NewRequest(http.MethodGet, "/astrate/v1/testrealm/socket", nil)
internal/appengine/stream/ws_test.go:76:		req := httptest.NewRequest(http.MethodGet, "/astrate/v1/testrealm/socket?device_id=dev1&interface=com.ex.S", nil)
internal/appengine/channels/ws_test.go:63:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?realm=testrealm", nil)
internal/appengine/channels/ws_test.go:76:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?token="+token, nil)
internal/appengine/channels/ws_test.go:89:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?realm=nosuchrealm&token="+token, nil)
internal/appengine/channels/ws_test.go:101:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?realm=testrealm&token=not.a.jwt", nil)
internal/appengine/channels/ws_test.go:114:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?realm=testrealm&token="+token, nil)
internal/appengine/channels/ws_test.go:128:		"/appengine/v1/socket/websocket?vsn=2.0.0&realm=" + testRealmName + "&token=" + token
internal/appengine/channels/protocol.go:2:// served at /appengine/v1/socket/websocket. It is the compatibility counterpart
internal/appengine/channels/ws.go:48:	mux.Handle("GET /appengine/v1/socket/websocket", http.HandlerFunc(a.handle))
internal/appengine/e2e_test.go:249:		"/appengine/v1/socket/websocket?vsn=2.0.0&realm=" + env.realm.Name + "&token=" + token
internal/swagger/swagger_test.go:800:	req := httptest.NewRequest(http.MethodGet, "/astrate/v1/metrics", nil)
internal/observability/observability_test.go:21:	m.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/metrics", nil))
internal/observability/observability_test.go:52:	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/health", nil))
internal/observability/observability_test.go:62:	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/readiness", nil))
internal/observability/observability_test.go:93:	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/readiness", nil))
internal/observability/observability_test.go:115:	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/readiness", nil))

$ rg -n 'HandleFunc|mux\.Handle|r\.Get|r\.Post' cmd/astrate/main.go | head -80
436:	mux.Handle("GET /appengine/v1/{realm}/version",
438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))

$ rg -n 'mux\.Handle|HandleFunc' cmd/astrate/main.go internal/observability/*.go internal/httpx/*.go
internal/httpx/notfound_test.go:13:	mux.HandleFunc("GET /appengine/v1/{realm}/devices", func(w http.ResponseWriter, _ *http.Request) {
internal/httpx/notfound_test.go:19:	mux.HandleFunc("GET /housekeeping/v1/realms", func(w http.ResponseWriter, _ *http.Request) {
internal/httpx/notfound_test.go:119:	mux.HandleFunc("GET /appengine/v1/{path...}", func(w http.ResponseWriter, _ *http.Request) {
internal/observability/compat_version_test.go:39:	mux.HandleFunc("GET /v1/{realm}/version", VersionHandler("x"))
internal/observability/compat.go:20:	mux.HandleFunc("GET /"+service+"/health", func(w http.ResponseWriter, r *http.Request) {
internal/observability/compat.go:51:	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
internal/observability/health.go:44:	mux.HandleFunc("GET /astrate/v1/health", h.handleHealth)
internal/observability/health.go:45:	mux.HandleFunc("GET /astrate/v1/readiness", h.handleReadiness)
internal/observability/health.go:46:	mux.Handle("GET /astrate/v1/metrics", h.metrics)
internal/httpx/notfound.go:16:		_, pattern := mux.Handler(r)
cmd/astrate/main.go:436:	mux.Handle("GET /appengine/v1/{realm}/version",
cmd/astrate/main.go:438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))

$ sed -n '1,120p' docs/api/astrate_native_api.yaml
openapi: "3.0.3"
info:
  title: Astrate Native API
  description: |
    Astrate-native operational and real-time API surface. Covers liveness,
    readiness, and Prometheus metrics (no auth), compat health and version
    endpoints for the Astarte Dashboard, and two WebSocket endpoints: the
    native live event stream and the Phoenix V2-compatible channels endpoint.
  version: "0.0.1"
  license:
    name: "Apache-2.0 with SECO Exclusion"
    url: https://www.apache.org/licenses/LICENSE-2.0

servers:
  - url: "{base_url}"
    variables:
      base_url:
        default: http://localhost:8080
        description: Astrate base URL

security: []

paths:
  # ── Health & Observability ──────────────────────────────────────────

  /astrate/v1/health:
    get:
      operationId: getHealth
      summary: Liveness probe
      description: |
        Returns `{"status": "ok"}` when the process is up and serving.
        Used by orchestrators (Docker, Kubernetes) as a liveness check.
      tags: [Observability]
      responses:
        "200":
          description: Service is alive.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/HealthStatus"
              example:
                status: ok

  /astrate/v1/readiness:
    get:
      operationId: getReadiness
      summary: Readiness probe
      description: |
        Runs every registered dependency check (database, broker). Returns 200
        when all checks pass, or 503 with per-check status when one or more
        dependencies are down.
      tags: [Observability]
      responses:
        "200":
          description: All dependencies are ready.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/ReadinessStatus"
              example:
                status: ok
                checks:
                  database: ok
                  broker: ok
        "503":
          description: One or more dependencies are unavailable.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/ReadinessStatus"
              example:
                status: unavailable
                checks:
                  database: ok
                  broker: "error: connection refused"

  /astrate/v1/metrics:
    get:
      operationId: getMetrics
      summary: Prometheus metrics endpoint
      description: |
        Returns metrics in Prometheus exposition format. The response content
        type is
        `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` —
        that is what a Prometheus scraper gets, since it asks for text (or
        sends no `Accept` at all).

        The handler negotiates `Accept`. A scraper that explicitly asks for the
        protobuf format gets the delimited binary `io.prometheus.client.MetricFamily`
        body instead, documented as a second media type below; every other
        `Accept` — including OpenMetrics and `application/zstd` — falls back to
        the text format. `Accept-Encoding: gzip` is honoured and answered with
        `Content-Encoding: gzip`, the only encoding offered: `zstd` and
        `deflate` are ignored and the body comes back uncompressed.
      tags: [Observability]
      responses:
        "200":
          description: Prometheus metrics.
          content:
            text/plain; version=0.0.4; charset=utf-8; escaping=underscores:
              schema:
                type: string
              example: |
                # HELP astrate_broker_sessions Live authenticated MQTT device sessions.
                # TYPE astrate_broker_sessions gauge
                astrate_broker_sessions 0
            application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores:
              schema:
                type: string
                format: binary

  # ── Compat health (Astarte Dashboard) ───────────────────────────────

  /appengine/health:
    get:
      operationId: getAppEngineHealth
      summary: AppEngine compat health
      description: |
        Upstream-parity unauthenticated health endpoint. The Astarte Dashboard
        polls this for its API status indicator.

$ rg -n --no-heading -o '(mux|sm|m)\.Handle(Func)?\("([A-Z]+ )?/[^"]+"' -r '$0' internal/ cmd/ | sed 's/.*HandleFunc\?("//' | sort -u
GET /appengine/v1/{path...}"
GET /appengine/v1/{realm}/devices"
GET /astrate/v1/health"
GET /astrate/v1/readiness"
GET /housekeeping/v1/realms"
GET /pairing/v1/{realm}/devices/{deviceID}"
GET /pairing/v1/{realm}/health"
GET /swagger"
GET /v1/{realm}/version"
POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials"
POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify"
cmd/astrate/main.go:436:mux.Handle("GET /appengine/v1/{realm}/version"
cmd/astrate/main.go:438:mux.Handle("GET /pairing/v1/{realm}/version"
internal/appengine/channels/ws.go:48:mux.Handle("GET /appengine/v1/socket/websocket"
internal/appengine/stream/ws.go:46:mux.Handle("GET /astrate/v1/{realm}/socket"
internal/auth/middleware_test.go:238:mux.Handle("GET /flow/v1/{realm}/pipelines"
internal/auth/middleware_test.go:41:mux.Handle("POST /pairing/v1/{realm}/agent/devices"
internal/auth/middleware_test.go:43:mux.Handle("GET /appengine/v1/{realm}/devices/{deviceID}"
internal/auth/middleware_test.go:45:mux.Handle("GET /housekeeping/v1/realms"
internal/flowapi/http.go:33:mux.Handle("GET /flow/v1/{realm}/pipelines"
internal/flowapi/http.go:34:mux.Handle("POST /flow/v1/{realm}/pipelines"
internal/flowapi/http.go:35:mux.Handle("GET /flow/v1/{realm}/pipelines/{name}"
internal/flowapi/http.go:36:mux.Handle("PUT /flow/v1/{realm}/pipelines/{name}"
internal/flowapi/http.go:37:mux.Handle("DELETE /flow/v1/{realm}/pipelines/{name}"
internal/flowapi/http.go:39:mux.Handle("GET /flow/v1/{realm}/flows"
internal/flowapi/http.go:40:mux.Handle("POST /flow/v1/{realm}/flows"
internal/flowapi/http.go:41:mux.Handle("GET /flow/v1/{realm}/flows/{name}"
internal/flowapi/http.go:42:mux.Handle("DELETE /flow/v1/{realm}/flows/{name}"
internal/flowapi/http.go:43:mux.Handle("POST /flow/v1/{realm}/flows/{name}/reload"
internal/flowapi/http.go:44:mux.Handle("PUT /flow/v1/{realm}/flows/{name}/config"
internal/flowapi/http.go:46:mux.Handle("GET /flow/v1/{realm}/blocks"
internal/flowapi/http.go:47:mux.Handle("POST /flow/v1/{realm}/blocks"
internal/flowapi/http.go:48:mux.Handle("GET /flow/v1/{realm}/blocks/{type}"
internal/flowapi/http.go:49:mux.Handle("PUT /flow/v1/{realm}/blocks/{name}"
internal/flowapi/http.go:50:mux.Handle("DELETE /flow/v1/{realm}/blocks/{name}"
internal/housekeeping/http.go:37:mux.Handle("GET /housekeeping/v1/realms"
internal/housekeeping/http.go:38:mux.Handle("POST /housekeeping/v1/realms"
internal/housekeeping/http.go:39:mux.Handle("GET /housekeeping/v1/realms/{realm}"
internal/housekeeping/http.go:40:mux.Handle("PATCH /housekeeping/v1/realms/{realm}"
internal/housekeeping/http.go:41:mux.Handle("DELETE /housekeeping/v1/realms/{realm}"
internal/observability/health.go:46:mux.Handle("GET /astrate/v1/metrics"
internal/pairing/http.go:82:mux.Handle("POST /pairing/v1/{realm}/agent/devices"
internal/pairing/http.go:84:mux.Handle("DELETE /pairing/v1/{realm}/agent/devices/{deviceID}"
internal/realm/http.go:38:mux.Handle("GET /realmmanagement/v1/{realm}/interfaces"
internal/realm/http.go:39:mux.Handle("POST /realmmanagement/v1/{realm}/interfaces"
internal/realm/http.go:40:mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}"
internal/realm/http.go:41:mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}/{major}"
internal/realm/http.go:42:mux.Handle("PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}"
internal/realm/http.go:43:mux.Handle("DELETE /realmmanagement/v1/{realm}/interfaces/{name}/{major}"
internal/realm/http.go:44:mux.Handle("GET /realmmanagement/v1/{realm}/triggers"
internal/realm/http.go:45:mux.Handle("POST /realmmanagement/v1/{realm}/triggers"
internal/realm/http.go:46:mux.Handle("GET /realmmanagement/v1/{realm}/triggers/{name}"
internal/realm/http.go:47:mux.Handle("DELETE /realmmanagement/v1/{realm}/triggers/{name}"
internal/realm/http.go:48:mux.Handle("GET /realmmanagement/v1/{realm}/config/auth"
internal/realm/http.go:49:mux.Handle("PUT /realmmanagement/v1/{realm}/config/auth"
internal/realm/http.go:50:mux.Handle("GET /realmmanagement/v1/{realm}/config/device_registration_limit"
internal/realm/http.go:51:mux.Handle("GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention"
internal/realm/http.go:52:mux.Handle("GET /realmmanagement/v1/{realm}/version"
internal/realm/http.go:53:mux.Handle("DELETE /realmmanagement/v1/{realm}/devices/{device}"
internal/realm/http.go:54:mux.Handle("GET /realmmanagement/v1/{realm}/policies"
internal/realm/http.go:55:mux.Handle("POST /realmmanagement/v1/{realm}/policies"
internal/realm/http.go:56:mux.Handle("GET /realmmanagement/v1/{realm}/policies/{name}"
internal/realm/http.go:57:mux.Handle("DELETE /realmmanagement/v1/{realm}/policies/{name}"
internal/swagger/swagger.go:33:mux.Handle("GET /swagger/"
internal/swagger/swagger.go:37:mux.Handle("GET /api/"

$ rg -n 'Handle\("|HandleFunc\("' internal/ cmd/ | rg -v '_test.go'
cmd/astrate/main.go:436:	mux.Handle("GET /appengine/v1/{realm}/version",
cmd/astrate/main.go:438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
internal/housekeeping/http.go:37:	mux.Handle("GET /housekeeping/v1/realms", h(a.listRealms))
internal/housekeeping/http.go:38:	mux.Handle("POST /housekeeping/v1/realms", h(a.createRealm))
internal/housekeeping/http.go:39:	mux.Handle("GET /housekeeping/v1/realms/{realm}", h(a.getRealm))
internal/housekeeping/http.go:40:	mux.Handle("PATCH /housekeeping/v1/realms/{realm}", h(a.patchRealm))
internal/housekeeping/http.go:41:	mux.Handle("DELETE /housekeeping/v1/realms/{realm}", h(a.deleteRealm))
internal/flowapi/http.go:33:	mux.Handle("GET /flow/v1/{realm}/pipelines", h(a.listPipelines))
internal/flowapi/http.go:34:	mux.Handle("POST /flow/v1/{realm}/pipelines", h(a.createPipeline))
internal/flowapi/http.go:35:	mux.Handle("GET /flow/v1/{realm}/pipelines/{name}", h(a.getPipeline))
internal/flowapi/http.go:36:	mux.Handle("PUT /flow/v1/{realm}/pipelines/{name}", h(a.updatePipeline))
internal/flowapi/http.go:37:	mux.Handle("DELETE /flow/v1/{realm}/pipelines/{name}", h(a.deletePipeline))
internal/flowapi/http.go:39:	mux.Handle("GET /flow/v1/{realm}/flows", h(a.listFlows))
internal/flowapi/http.go:40:	mux.Handle("POST /flow/v1/{realm}/flows", h(a.startFlow))
internal/flowapi/http.go:41:	mux.Handle("GET /flow/v1/{realm}/flows/{name}", h(a.getFlow))
internal/flowapi/http.go:42:	mux.Handle("DELETE /flow/v1/{realm}/flows/{name}", h(a.stopFlow))
internal/flowapi/http.go:43:	mux.Handle("POST /flow/v1/{realm}/flows/{name}/reload", h(a.reloadFlow))
internal/flowapi/http.go:44:	mux.Handle("PUT /flow/v1/{realm}/flows/{name}/config", h(a.updateFlowConfig))
internal/flowapi/http.go:46:	mux.Handle("GET /flow/v1/{realm}/blocks", h(a.listBlocks))
internal/flowapi/http.go:47:	mux.Handle("POST /flow/v1/{realm}/blocks", h(a.createUserBlock))
internal/flowapi/http.go:48:	mux.Handle("GET /flow/v1/{realm}/blocks/{type}", h(a.getBlock))
internal/flowapi/http.go:49:	mux.Handle("PUT /flow/v1/{realm}/blocks/{name}", h(a.updateUserBlock))
internal/flowapi/http.go:50:	mux.Handle("DELETE /flow/v1/{realm}/blocks/{name}", h(a.deleteUserBlock))
internal/realm/http.go:38:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces", h(a.listInterfaces))
internal/realm/http.go:39:	mux.Handle("POST /realmmanagement/v1/{realm}/interfaces", h(a.installInterface))
internal/realm/http.go:40:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}", h(a.listInterfaceMajors))
internal/realm/http.go:41:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.getInterface))
internal/realm/http.go:42:	mux.Handle("PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.updateInterface))
internal/realm/http.go:43:	mux.Handle("DELETE /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.deleteInterface))
internal/realm/http.go:44:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers", h(a.listTriggers))
internal/realm/http.go:45:	mux.Handle("POST /realmmanagement/v1/{realm}/triggers", h(a.createTrigger))
internal/realm/http.go:46:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers/{name}", h(a.getTrigger))
internal/realm/http.go:47:	mux.Handle("DELETE /realmmanagement/v1/{realm}/triggers/{name}", h(a.deleteTrigger))
internal/realm/http.go:48:	mux.Handle("GET /realmmanagement/v1/{realm}/config/auth", h(a.getAuth))
internal/realm/http.go:49:	mux.Handle("PUT /realmmanagement/v1/{realm}/config/auth", h(a.putAuth))
internal/realm/http.go:50:	mux.Handle("GET /realmmanagement/v1/{realm}/config/device_registration_limit", h(a.getRegistrationLimit))
internal/realm/http.go:51:	mux.Handle("GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention", h(a.getDatastreamMaximumStorageRetention))
internal/realm/http.go:52:	mux.Handle("GET /realmmanagement/v1/{realm}/version", h(a.getVersion))
internal/realm/http.go:53:	mux.Handle("DELETE /realmmanagement/v1/{realm}/devices/{device}", h(a.deleteDevice))
internal/realm/http.go:54:	mux.Handle("GET /realmmanagement/v1/{realm}/policies", h(a.listPolicies))
internal/realm/http.go:55:	mux.Handle("POST /realmmanagement/v1/{realm}/policies", h(a.createPolicy))
internal/realm/http.go:56:	mux.Handle("GET /realmmanagement/v1/{realm}/policies/{name}", h(a.getPolicy))
internal/realm/http.go:57:	mux.Handle("DELETE /realmmanagement/v1/{realm}/policies/{name}", h(a.deletePolicy))
internal/pairing/http.go:81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/pairing/http.go:82:	mux.Handle("POST /pairing/v1/{realm}/agent/devices",
internal/pairing/http.go:84:	mux.Handle("DELETE /pairing/v1/{realm}/agent/devices/{deviceID}",
internal/pairing/http.go:86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
internal/pairing/http.go:88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
internal/pairing/http.go:90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
internal/observability/compat.go:20:	mux.HandleFunc("GET /"+service+"/health", func(w http.ResponseWriter, r *http.Request) {
internal/observability/compat.go:51:	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
internal/observability/health.go:44:	mux.HandleFunc("GET /astrate/v1/health", h.handleHealth)
internal/observability/health.go:45:	mux.HandleFunc("GET /astrate/v1/readiness", h.handleReadiness)
internal/observability/health.go:46:	mux.Handle("GET /astrate/v1/metrics", h.metrics)
internal/appengine/stream/ws.go:46:	mux.Handle("GET /astrate/v1/{realm}/socket", a.require(http.HandlerFunc(a.handle)))
internal/appengine/http.go:44:	mux.Handle("GET "+base+"/devices", h(a.listDevices))
internal/appengine/http.go:45:	mux.Handle("GET "+base+"/stats/devices", h(a.devicesStats))
internal/appengine/http.go:46:	mux.Handle("GET "+base+"/devices/{device}", h(a.getDevice))
internal/appengine/http.go:47:	mux.Handle("PATCH "+base+"/devices/{device}", h(a.patchDevice))
internal/appengine/http.go:48:	mux.Handle("GET "+base+"/devices-by-alias/{alias}", h(a.getDeviceByAlias))
internal/appengine/http.go:49:	mux.Handle("GET "+base+"/devices/{device}/interfaces", h(a.listDeviceInterfaces))
internal/appengine/http.go:51:	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}", h(a.getData))
internal/appengine/http.go:52:	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.getData))
internal/appengine/http.go:53:	mux.Handle("PUT "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
internal/appengine/http.go:54:	mux.Handle("POST "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
internal/appengine/http.go:55:	mux.Handle("DELETE "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteData))
internal/appengine/http.go:60:	mux.Handle("PATCH "+base+"/devices-by-alias/{alias}", h(a.patchDeviceByAlias))
internal/appengine/http.go:61:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces", h(a.listInterfacesByAlias))
internal/appengine/http.go:62:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces/{interface}", h(a.getDataByAlias))
internal/appengine/http.go:63:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.getDataByAlias))
internal/appengine/http.go:64:	mux.Handle("PUT "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.putDataByAlias))
internal/appengine/http.go:65:	mux.Handle("POST "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.putDataByAlias))
internal/appengine/http.go:66:	mux.Handle("DELETE "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.deleteDataByAlias))
internal/appengine/http.go:68:	mux.Handle("GET "+base+"/groups", h(a.listGroups))
internal/appengine/http.go:69:	mux.Handle("POST "+base+"/groups", h(a.createGroup))
internal/appengine/http.go:70:	mux.Handle("GET "+base+"/groups/{group}", h(a.getGroup))
internal/appengine/http.go:71:	mux.Handle("GET "+base+"/groups/{group}/devices", h(a.listGroupDevices))
internal/appengine/http.go:72:	mux.Handle("POST "+base+"/groups/{group}/devices", h(a.addGroupDevice))
internal/appengine/http.go:73:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}", h(a.getGroupDevice))
internal/appengine/http.go:74:	mux.Handle("PATCH "+base+"/groups/{group}/devices/{device}", h(a.patchGroupDevice))
internal/appengine/http.go:75:	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}", h(a.removeGroupDevice))
internal/appengine/http.go:76:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces", h(a.listInterfacesInGroup))
internal/appengine/http.go:77:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}", h(a.getDataInGroup))
internal/appengine/http.go:78:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.getDataInGroup))
internal/appengine/http.go:79:	mux.Handle("PUT "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.putDataInGroup))
internal/appengine/http.go:80:	mux.Handle("POST "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.putDataInGroup))
internal/appengine/http.go:81:	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteDataInGroup))
internal/appengine/channels/ws.go:48:	mux.Handle("GET /appengine/v1/socket/websocket", http.HandlerFunc(a.handle))
internal/swagger/swagger.go:30:	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
internal/swagger/swagger.go:33:	mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(uiRoot))))
internal/swagger/swagger.go:37:	mux.Handle("GET /api/", http.StripPrefix("/api/", http.FileServer(http.FS(apiRoot))))

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

$ rg -n '^\s{2}/' docs/api/astarte_realm_management_api.yaml; echo "=== PAIRING ==="; rg -n '^\s{2}/' docs/api/astarte_pairing_api.yaml; echo "=== HK ==="; rg -n '^\s{2}/' docs/api/astarte_housekeeping_api.yaml
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
=== PAIRING ===
27:  /pairing/v1/{realm}/health:
70:  /pairing/v1/{realm}/version:
97:  /pairing/v1/{realm}/agent/devices:
149:  /pairing/v1/{realm}/agent/devices/{deviceID}:
183:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
235:  /pairing/v1/{realm}/devices/{deviceID}:
269:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:
=== HK ===
26:  /housekeeping/v1/realms:
108:  /housekeeping/v1/realms/{realm}:

$ rg -n 'func \(a \*API\)|Write[A-Z]|Status[A-Z]|http.Status' internal/appengine/http.go | sed -n '1,200p'
41:func (a *API) Mount(mux *http.ServeMux) {
86:func (a *API) listDevices(w http.ResponseWriter, r *http.Request) {
99:	_ = astarteapi.WriteDataWithLinks(w, http.StatusOK, body,
131:func (a *API) devicesStats(w http.ResponseWriter, r *http.Request) {
137:	_ = astarteapi.WriteData(w, http.StatusOK, map[string]int64{
143:func (a *API) getDevice(w http.ResponseWriter, r *http.Request) {
149:	_ = astarteapi.WriteData(w, http.StatusOK, st)
152:func (a *API) getDeviceByAlias(w http.ResponseWriter, r *http.Request) {
158:	_ = astarteapi.WriteData(w, http.StatusOK, st)
163:func (a *API) writeInterfaces(w http.ResponseWriter, names []string) {
167:	_ = astarteapi.WriteData(w, http.StatusOK, names)
170:func (a *API) listDeviceInterfaces(w http.ResponseWriter, r *http.Request) {
179:func (a *API) listInterfacesByAlias(w http.ResponseWriter, r *http.Request) {
188:func (a *API) listInterfacesInGroup(w http.ResponseWriter, r *http.Request) {
205:func (a *API) patchDevice(w http.ResponseWriter, r *http.Request) {
211:		_ = astarteapi.WriteInternalServerError(w)
216:		_ = astarteapi.WriteBadRequest(w)
224:	_ = astarteapi.WriteData(w, http.StatusOK, st)
227:func (a *API) patchDeviceByAlias(w http.ResponseWriter, r *http.Request) {
232:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Bad request")
237:		_ = astarteapi.WriteBadRequest(w)
245:	_ = astarteapi.WriteData(w, http.StatusOK, st)
254:func (a *API) serveData(w http.ResponseWriter, r *http.Request, read func(QueryOpts) (any, error)) {
266:		_ = astarteapi.WriteDataWithMetadata(w, http.StatusOK, t.Data, t.Metadata)
269:	_ = astarteapi.WriteData(w, http.StatusOK, data)
272:func (a *API) getData(w http.ResponseWriter, r *http.Request) {
279:func (a *API) getDataByAlias(w http.ResponseWriter, r *http.Request) {
286:func (a *API) getDataInGroup(w http.ResponseWriter, r *http.Request) {
293:func (a *API) putData(w http.ResponseWriter, r *http.Request) {
296:		_ = astarteapi.WriteBadRequest(w)
300:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
309:	w.WriteHeader(http.StatusOK)
312:func (a *API) deleteData(w http.ResponseWriter, r *http.Request) {
319:	w.WriteHeader(http.StatusNoContent)
322:func (a *API) putDataByAlias(w http.ResponseWriter, r *http.Request) {
325:		_ = astarteapi.WriteBadRequest(w)
329:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
338:	w.WriteHeader(http.StatusOK)
341:func (a *API) putDataInGroup(w http.ResponseWriter, r *http.Request) {
344:		_ = astarteapi.WriteBadRequest(w)
348:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
357:	w.WriteHeader(http.StatusOK)
360:func (a *API) deleteDataByAlias(w http.ResponseWriter, r *http.Request) {
367:	w.WriteHeader(http.StatusNoContent)
370:func (a *API) deleteDataInGroup(w http.ResponseWriter, r *http.Request) {
377:	w.WriteHeader(http.StatusNoContent)
382:func (a *API) listGroups(w http.ResponseWriter, r *http.Request) {
388:	_ = astarteapi.WriteData(w, http.StatusOK, names)
399:func (a *API) createGroup(w http.ResponseWriter, r *http.Request) {
402:		_ = astarteapi.WriteBadRequest(w)
407:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
411:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
415:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
424:	_ = astarteapi.WriteData(w, http.StatusCreated, body)
429:func (a *API) getGroup(w http.ResponseWriter, r *http.Request) {
434:	_ = astarteapi.WriteData(w, http.StatusOK, map[string]string{"group_name": r.PathValue("group")})
437:func (a *API) listGroupDevices(w http.ResponseWriter, r *http.Request) {
451:	_ = astarteapi.WriteDataWithLinks(w, http.StatusOK, body,
487:func (a *API) addGroupDevice(w http.ResponseWriter, r *http.Request) {
490:		_ = astarteapi.WriteBadRequest(w)
497:	w.WriteHeader(http.StatusCreated)
500:func (a *API) removeGroupDevice(w http.ResponseWriter, r *http.Request) {
505:	w.WriteHeader(http.StatusNoContent)
510:func (a *API) getGroupDevice(w http.ResponseWriter, r *http.Request) {
516:	_ = astarteapi.WriteData(w, http.StatusOK, st)
519:func (a *API) patchGroupDevice(w http.ResponseWriter, r *http.Request) {
523:		_ = astarteapi.WriteInternalServerError(w)
528:		_ = astarteapi.WriteBadRequest(w)
537:	_ = astarteapi.WriteData(w, http.StatusOK, st)
644:func (a *API) writeError(w http.ResponseWriter, err error) {
648:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity, fe)
650:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Invalid alias")
652:		_ = astarteapi.WriteError(w, http.StatusConflict, "Alias already in use")
654:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Alias tag not found")
656:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Invalid attributes")
658:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Attribute key not found")
660:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
662:		_ = astarteapi.WriteError(w, http.StatusConflict, "Group already exists")
664:		_ = astarteapi.WriteError(w, http.StatusConflict, "Device already in group")
666:		_ = astarteapi.WriteError(w, http.StatusConflict, "Already exists")
668:		_ = astarteapi.WriteError(w, http.StatusNotFound, "Group not found")
670:		_ = astarteapi.WriteError(w, http.StatusNotFound, "Path not found")
675:		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Cannot write to device owned resource")
677:		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Cannot write to read-only resource")
679:		_ = astarteapi.WriteError(w, http.StatusNotFound, "Interface not found in device introspection")
681:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Endpoint not found")
683:		_ = astarteapi.WriteDeviceNotFound(w)
702:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, detailMissingRequiredMapping)
704:		_ = astarteapi.WriteInternalServerError(w)
728:func (a *API) writeBadObjectError(w http.ResponseWriter, err error) {
731:		_ = astarteapi.WriteBadRequest(w)
734:	_ = astarteapi.WriteRawErrors(w, http.StatusBadRequest, struct {

$ sed -n '1,60p' internal/appengine/http.go
package appengine

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/engine"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// maxBodyBytes caps AppEngine request bodies.
const maxBodyBytes int64 = 1 << 20

// maxValueBytes is upstream's per-value limit (~64 KiB): a larger REST
// write is answered 422 "Value size exceeds size limits" (measured
// 2026-08-24 against 1.2.0, verify-server-writes.json).
const maxValueBytes = 65536

// API is the /appengine/v1 HTTP surface (docs/ROADMAP.md §8.2 file 7.8),
// guarded by a realm JWT carrying a_aea.
type API struct {
	svc     *Service
	require func(http.Handler) http.Handler
}

// NewAPI wires the AppEngine service to its HTTP surface.
func NewAPI(svc *Service, mw *auth.Middleware) *API {
	return &API{svc: svc, require: mw.RequireRealm(auth.ClaimAppEngine)}
}

// Mount registers the routes on mux (paths wire-frozen to upstream
// astarte_appengine_api).
func (a *API) Mount(mux *http.ServeMux) {
	h := func(f http.HandlerFunc) http.Handler { return a.require(f) }
	const base = "/appengine/v1/{realm}"
	mux.Handle("GET "+base+"/devices", h(a.listDevices))
	mux.Handle("GET "+base+"/stats/devices", h(a.devicesStats))
	mux.Handle("GET "+base+"/devices/{device}", h(a.getDevice))
	mux.Handle("PATCH "+base+"/devices/{device}", h(a.patchDevice))
	mux.Handle("GET "+base+"/devices-by-alias/{alias}", h(a.getDeviceByAlias))
	mux.Handle("GET "+base+"/devices/{device}/interfaces", h(a.listDeviceInterfaces))

	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}", h(a.getData))
	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.getData))
	mux.Handle("PUT "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
	mux.Handle("POST "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
	mux.Handle("DELETE "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteData))

	// Mirror surfaces (upstream registers these alongside the device-scoped
	// ones; Go 1.22 precedence resolves the overlap with the group routes
	// above/below without re-registering them).
	mux.Handle("PATCH "+base+"/devices-by-alias/{alias}", h(a.patchDeviceByAlias))

$ sed -n '697,880p' docs/api/astarte_appengine_api.yaml
/appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
    get:
      operationId: getDataAtPath
      summary: Get data at path
      description: Returns data for a specific path within a device interface.
      tags: [AppEngine - Data]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - $ref: "#/components/parameters/DeviceID"
        - $ref: "#/components/parameters/InterfaceName"
        - name: path
          in: path
          required: true
          description: The data path (e.g. `/sensors/temperature/value`).
          schema:
            type: string
        - $ref: "#/components/parameters/DataSince"
        - $ref: "#/components/parameters/DataSinceAfter"
        - $ref: "#/components/parameters/DataTo"
        - $ref: "#/components/parameters/DataLimit"
        - $ref: "#/components/parameters/DataDownsample"
        - $ref: "#/components/parameters/DataSort"
        - $ref: "#/components/parameters/DataFormat"
        - $ref: "#/components/parameters/DataAllowBigIntegers"
        - $ref: "#/components/parameters/DataAllowSafeBigIntegers"
      responses:
        "200":
          description: Path data.
          content:
            application/json; charset=utf-8:
              schema:
                $ref: "#/components/schemas/DataEnvelopeObject"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "404":
          $ref: "#/components/responses/DeviceNotFound"
        "422":
          $ref: "#/components/responses/ValidationErrors"
        "500":
          $ref: "#/components/responses/InternalServerError"

    put:
      operationId: putData
      summary: Set data at path
      description: |
        Sets the value at a specific path in a device interface (server-owned
        data). Returns 200 on success.
      tags: [AppEngine - Data]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - $ref: "#/components/parameters/DeviceID"
        - $ref: "#/components/parameters/InterfaceName"
        - name: path
          in: path
          required: true
          description: The data path.
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json; charset=utf-8:
            schema:
              type: object
              required: [data]
              properties:
                data:
                  description: The value to set.
      responses:
        "200":
          description: Data set successfully.
        "400":
          $ref: "#/components/responses/BadRequest"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "404":
          $ref: "#/components/responses/DeviceNotFound"
        "405":
          $ref: "#/components/responses/MethodNotAllowed"
        "422":
          $ref: "#/components/responses/ValueTooLarge"
        "500":
          $ref: "#/components/responses/InternalServerError"

    post:
      operationId: publishData
      summary: Publish data at path
      description: |
        Alias for PUT — publishes a value to a device interface path.
        Returns 200 on success.
      tags: [AppEngine - Data]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - $ref: "#/components/parameters/DeviceID"
        - $ref: "#/components/parameters/InterfaceName"
        - name: path
          in: path
          required: true
          description: The data path.
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json; charset=utf-8:
            schema:
              type: object
              required: [data]
              properties:
                data:
                  description: The value to publish.
      responses:
        "200":
          description: Data published successfully.
        "400":
          $ref: "#/components/responses/BadRequest"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "404":
          $ref: "#/components/responses/DeviceNotFound"
        "405":
          $ref: "#/components/responses/MethodNotAllowed"
        "422":
          $ref: "#/components/responses/ValueTooLarge"
        "500":
          $ref: "#/components/responses/InternalServerError"

    delete:
      operationId: deleteData
      summary: Delete data at path
      description: |
        Unsets the property at a specific path (properties interface only).
        Returns 204 on success. A path matching no endpoint mapping is a 400,
        not a 404.
      tags: [AppEngine - Data]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - $ref: "#/components/parameters/DeviceID"
        - $ref: "#/components/parameters/InterfaceName"
        - name: path
          in: path
          required: true
          description: The data path.
          schema:
            type: string
      responses:
        "204":
          description: Data deleted. No body returned.
        "400":
          $ref: "#/components/responses/BadRequest"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "404":
          $ref: "#/components/responses/DeviceNotFound"
        "405":
          $ref: "#/components/responses/MethodNotAllowed"
        "500":
          $ref: "#/components/responses/InternalServerError"

  # ── Groups ──────────────────────────────────────────────────────────

  /appengine/v1/{realm}/groups:
    get:
      operationId: listGroups
      summary: List groups
      description: Returns the names of all groups in the realm.
      tags: [AppEngine - Groups]
      security:

$ sed -n '254,380p' internal/appengine/http.go
func (a *API) serveData(w http.ResponseWriter, r *http.Request, read func(QueryOpts) (any, error)) {
	opts, err := parseQueryOpts(r)
	if err != nil {
		a.writeError(w, err)
		return
	}
	data, err := read(opts)
	if err != nil {
		a.writeError(w, err)
		return
	}
	if t, ok := data.(*Tabular); ok {
		_ = astarteapi.WriteDataWithMetadata(w, http.StatusOK, t.Data, t.Metadata)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, data)
}

func (a *API) getData(w http.ResponseWriter, r *http.Request) {
	a.serveData(w, r, func(opts QueryOpts) (any, error) {
		return a.svc.GetData(r.Context(), r.PathValue("realm"), r.PathValue("device"),
			r.PathValue("interface"), pathParam(r), opts)
	})
}

func (a *API) getDataByAlias(w http.ResponseWriter, r *http.Request) {
	a.serveData(w, r, func(opts QueryOpts) (any, error) {
		return a.svc.GetDataByAlias(r.Context(), r.PathValue("realm"), r.PathValue("alias"),
			r.PathValue("interface"), pathParam(r), opts)
	})
}

func (a *API) getDataInGroup(w http.ResponseWriter, r *http.Request) {
	a.serveData(w, r, func(opts QueryOpts) (any, error) {
		return a.svc.GetDataInGroup(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device"),
			r.PathValue("interface"), pathParam(r), opts)
	})
}

func (a *API) putData(w http.ResponseWriter, r *http.Request) {
	var value json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &value); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if len(value) > maxValueBytes {
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
		return
	}
	err := a.svc.PublishData(r.Context(), r.PathValue("realm"), r.PathValue("device"),
		r.PathValue("interface"), pathParam(r), value, nil)
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) deleteData(w http.ResponseWriter, r *http.Request) {
	err := a.svc.UnsetProperty(r.Context(), r.PathValue("realm"), r.PathValue("device"),
		r.PathValue("interface"), pathParam(r))
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) putDataByAlias(w http.ResponseWriter, r *http.Request) {
	var value json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &value); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if len(value) > maxValueBytes {
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
		return
	}
	err := a.svc.PublishDataByAlias(r.Context(), r.PathValue("realm"), r.PathValue("alias"),
		r.PathValue("interface"), pathParam(r), value, nil)
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) putDataInGroup(w http.ResponseWriter, r *http.Request) {
	var value json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &value); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if len(value) > maxValueBytes {
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
		return
	}
	err := a.svc.PublishDataInGroup(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device"),
		r.PathValue("interface"), pathParam(r), value, nil)
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) deleteDataByAlias(w http.ResponseWriter, r *http.Request) {
	err := a.svc.UnsetPropertyByAlias(r.Context(), r.PathValue("realm"), r.PathValue("alias"),
		r.PathValue("interface"), pathParam(r))
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteDataInGroup(w http.ResponseWriter, r *http.Request) {
	err := a.svc.UnsetPropertyInGroup(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device"),
		r.PathValue("interface"), pathParam(r))
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- groups -----------------------------------------------------------------

$ sed -n '141,181p' docs/api/astarte_appengine_api.yaml; echo "=== code ==="; sed -n '131,160p' internal/appengine/http.go
/appengine/v1/{realm}/stats/devices:
    get:
      operationId: getDeviceStats
      summary: Get device statistics
      description: Returns total and connected device counts for the realm.
      tags: [AppEngine - Devices]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
      responses:
        "200":
          description: Device statistics.
          content:
            application/json; charset=utf-8:
              schema:
                type: object
                required: [data]
                properties:
                  data:
                    type: object
                    required: [connected_devices, total_devices]
                    properties:
                      connected_devices:
                        type: integer
                        format: int64
                      total_devices:
                        type: integer
                        format: int64
              example:
                data:
                  connected_devices: 42
                  total_devices: 150
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "500":
          $ref: "#/components/responses/InternalServerError"

  /appengine/v1/{realm}/devices/{device}:
=== code ===
func (a *API) devicesStats(w http.ResponseWriter, r *http.Request) {
	total, connected, err := a.svc.DevicesStats(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, map[string]int64{
		"connected_devices": connected,
		"total_devices":     total,
	})
}

func (a *API) getDevice(w http.ResponseWriter, r *http.Request) {
	st, err := a.svc.GetDevice(r.Context(), r.PathValue("realm"), r.PathValue("device"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, st)
}

func (a *API) getDeviceByAlias(w http.ResponseWriter, r *http.Request) {
	st, err := a.svc.GetDeviceByAlias(r.Context(), r.PathValue("realm"), r.PathValue("alias"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, st)
}

$ sed -n '67,141p' docs/api/astarte_appengine_api.yaml; echo "=== code ==="; sed -n '86,131p' internal/appengine/http.go
/appengine/v1/{realm}/devices:
    get:
      operationId: listDevices
      summary: List devices in a realm
      description: |
        Returns a paginated list of device IDs (or device status objects when
        `details=true`). Supports cursor-based pagination via `links.next`.
      tags: [AppEngine - Devices]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - name: details
          in: query
          schema:
            type: boolean
            default: false
          description: When true, return full device status objects instead of just IDs.
        - name: from_token
          in: query
          schema:
            type: string
          description: Pagination cursor — start after this device token.
        - name: limit
          in: query
          schema:
            type: integer
          description: Maximum number of devices to return.
      responses:
        "200":
          description: Paginated list of devices.
          content:
            application/json; charset=utf-8:
              schema:
                oneOf:
                  - $ref: "#/components/schemas/DataEnvelopeDeviceIDList"
                  - $ref: "#/components/schemas/DataEnvelopeDeviceStatusList"
              examples:
                ids:
                  summary: Simple device ID list
                  value:
                    data:
                      - dT6hS2W9TT6LEnP25ks_lg
                    links:
                      self: "/v1/test/devices?details=false&limit=1"
                details:
                  summary: Device status list with pagination
                  value:
                    data:
                      - id: dT6hS2W9TT6LEnP25ks_lg
                        aliases: {}
                        attributes: {}
                        connected: true
                        last_connection: "2024-05-30 13:49:57.045Z"
                        first_registration: "2024-01-15 10:00:00.000Z"
                        last_disconnection: null
                        credentials_inhibited: false
                        groups: []
                        total_received_msgs: 12
                        total_received_bytes: 4096
                        introspection:
                          org.example.Sensors:
                            major: 0
                            minor: 1
                    links:
                      self: "/v1/test/devices?details=true&limit=1"
                      next: "/v1/test/devices?details=true&from_token=dT6hS2W9TT6LEnP25ks_lg&limit=1"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "500":
          $ref: "#/components/responses/InternalServerError"

  /appengine/v1/{realm}/stats/devices:
=== code ===
func (a *API) listDevices(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	details, _ := strconv.ParseBool(q.Get("details"))
	page, err := a.svc.ListDevices(r.Context(), r.PathValue("realm"), q.Get("from_token"), limit, details)
	if err != nil {
		a.writeError(w, err)
		return
	}
	var body any = page.IDs
	if details {
		body = page.Statuses
	}
	_ = astarteapi.WriteDataWithLinks(w, http.StatusOK, body,
		deviceListLinks(r.PathValue("realm"), q, page.Next))
}

// deviceListLinks builds the upstream pagination links object. links.next
// always carries details, from_token, and limit: url.Values.Encode sorts
// keys, so "details" always precedes "from_token" — the dashboard parses
// links.next with URLSearchParams over the raw path+query string, which
// corrupts the FIRST key/value pair, so from_token must never be first.
// Absent request values are filled with their effective defaults.
func deviceListLinks(realm string, q url.Values, next string) astarteapi.Links {
	base := "/v1/" + realm + "/devices"
	self := base
	if enc := q.Encode(); enc != "" {
		self += "?" + enc
	}
	links := astarteapi.Links{Self: self}
	if next != "" {
		details := q.Get("details")
		if details == "" {
			details = "false"
		}
		limit := q.Get("limit")
		if limit == "" {
			limit = strconv.Itoa(DefaultDeviceLimit)
		}
		nq := url.Values{"details": {details}, "from_token": {next}, "limit": {limit}}
		links.Next = base + "?" + nq.Encode()
	}
	return links
}

func (a *API) devicesStats(w http.ResponseWriter, r *http.Request) {

→ Read internal/observability/health.go

$ sed -n '120,290p' docs/api/astrate_native_api.yaml
polls this for its API status indicator.
      tags: [Compat Health]
      responses:
        "200":
          description: Service is healthy.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DataEnvelopeHealthStatus"
              example:
                data:
                  status: ok
        "503":
          description: Service is unhealthy — readiness probe failed.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DataEnvelopeHealthStatus"
              example:
                data:
                  status: unhealthy

  /realmmanagement/health:
    get:
      operationId: getRealmManagementHealth
      summary: Realm Management compat health
      description: |
        Upstream-parity unauthenticated health endpoint. The Astarte Dashboard
        polls this for its API status indicator.
      tags: [Compat Health]
      responses:
        "200":
          description: Service is healthy.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DataEnvelopeHealthStatus"
              example:
                data:
                  status: ok
        "503":
          description: Service is unhealthy — readiness probe failed.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DataEnvelopeHealthStatus"
              example:
                data:
                  status: unhealthy

  /pairing/health:
    get:
      operationId: getPairingHealth
      summary: Pairing compat health
      description: |
        Upstream-parity unauthenticated health endpoint. The Astarte Dashboard
        polls this for its API status indicator.
      tags: [Compat Health]
      responses:
        "200":
          description: Service is healthy.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DataEnvelopeHealthStatus"
              example:
                data:
                  status: ok
        "503":
          description: Service is unhealthy — readiness probe failed.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DataEnvelopeHealthStatus"
              example:
                data:
                  status: unhealthy

  # ── Compat version (upstream parity) ────────────────────────────────

  /appengine/version:
    get:
      operationId: getAppEngineVersion
      summary: AppEngine API version
      description: |
        Upstream-parity unauthenticated version endpoint. Returns the Astrate
        build version serving the AppEngine service; the Astarte Dashboard
        reads one per service at startup. The emulated upstream API level is
        served by the realm-scoped /v1/{realm}/version endpoints.
      tags: [Compat Health]
      responses:
        "200":
          description: Astrate build version.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DataEnvelopeVersion"
              example:
                data: "0.1.0-dev"

  /realmmanagement/version:
    get:
      operationId: getRealmManagementVersion
      summary: Realm Management API version
      description: |
        Upstream-parity unauthenticated version endpoint. Returns the Astrate
        build version serving the Realm Management service; the Astarte
        Dashboard reads one per service at startup. The emulated upstream API
        level is served by the realm-scoped /v1/{realm}/version endpoints.
      tags: [Compat Health]
      responses:
        "200":
          description: Astrate build version.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DataEnvelopeVersion"
              example:
                data: "0.1.0-dev"

  /pairing/version:
    get:
      operationId: getPairingVersion
      summary: Pairing API version
      description: |
        Upstream-parity unauthenticated version endpoint. Returns the Astrate
        build version serving the Pairing service; the Astarte Dashboard reads
        one per service at startup. The emulated upstream API level is served
        by the realm-scoped /v1/{realm}/version endpoints.
      tags: [Compat Health]
      responses:
        "200":
          description: Astrate build version.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DataEnvelopeVersion"
              example:
                data: "0.1.0-dev"

  /housekeeping/version:
    get:
      operationId: getHousekeepingVersion
      summary: Housekeeping API version
      description: |
        Upstream-parity unauthenticated version endpoint. Returns the Astrate
        build version serving the Housekeeping service; the Astarte Dashboard
        reads one per service at startup. Housekeeping is the one service with
        no realm-scoped twin: `/housekeeping/v1/{realm}/version` is not a route
        and answers 404, so unlike AppEngine, Realm Management and Pairing no
        emulated upstream API level is reported here.
      tags: [Compat Health]
      responses:
        "200":
          description: Astrate build version.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DataEnvelopeVersion"
              example:
                data: "0.1.0-dev"

  # ── Native WebSocket (live stream) ──────────────────────────────────

  /astrate/v1/{realm}/socket:
    get:
      operationId: nativeWebSocket
      summary: Native live event stream (WebSocket / SSE)
      description: |
        Astrate-native real-time event stream. Upgrades to WebSocket by
        default, or streams Server-Sent Events when the client sends

→ Read internal/observability/compat.go

$ sed -n '284,509p' docs/api/astrate_native_api.yaml
/astrate/v1/{realm}/socket:
    get:
      operationId: nativeWebSocket
      summary: Native live event stream (WebSocket / SSE)
      description: |
        Astrate-native real-time event stream. Upgrades to WebSocket by
        default, or streams Server-Sent Events when the client sends
        `Accept: text/event-stream` or `?transport=sse`.

        Subscribes to the realm's committed data and device lifecycle events.
        Optional `device_id` and `interface` query parameters narrow the
        event filter.

        **WebSocket protocol:** server-push only — client frames are drained
        and ignored. Events are JSON objects:

        ```json
        {
          "event": "incoming_data",
          "realm": "test",
          "device_id": "f0VMRgIBAQAAAAAAAAAAAA",
          "interface": "org.example.Sensors",
          "path": "/sensors/temperature/value",
          "value": 23.5,
          "timestamp": "2024-05-30T13:49:57.045Z"
        }
        ```

        **SSE protocol:** each event is a `data:` frame followed by a blank line.

        Authentication: `a_ch` realm JWT via Authorization header.
      tags: [WebSocket]
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - name: device_id
          in: query
          schema:
            type: string
          description: Filter events to this device ID.
        - name: interface
          in: query
          schema:
            type: string
          description: Filter events to this interface.
        - name: transport
          in: query
          schema:
            type: string
            enum: [websocket, sse]
            default: websocket
          description: Force SSE transport instead of WebSocket.
      responses:
        "101":
          description: WebSocket upgrade successful.
        "200":
          description: "SSE stream started (when `Accept: text/event-stream`)."
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "500":
          $ref: "#/components/responses/InternalServerError"

  # ── Phoenix V2 compat WebSocket ─────────────────────────────────────

  /appengine/v1/socket/websocket:
    get:
      operationId: phoenixWebSocket
      summary: Phoenix V2-compatible WebSocket endpoint
      description: |
        Upstream-compatible Phoenix Channels WebSocket. Authenticates via
        `?realm=` and `?token=` query parameters (realm JWT).

        Uses the Phoenix V2 frame protocol:
        - Client sends `phx_join` to join a room (`rooms:{realm}:{name}`)
        - Client sends `heartbeat` for keepalive
        - Client sends `watch` / `unwatch` to manage triggers
        - Server pushes `new_event` frames for matching events

        Frame format:
        ```json
        [null, null, "rooms:test:all", "new_event", {"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}]
        ```

        The payload is the bus event marshalled straight into the frame, so
        its keys are the capitalised Go field names — unlike
        `/astrate/v1/{realm}/socket` above, whose example is snake_case
        because that socket serialises through its own `wireEvent` struct.
        A client generated from any other shape (snake_case, `event_data`)
        reads nothing at all here: Go's unmarshal leaves unknown fields zero
        rather than failing, so every field arrives empty with no error.
      tags: [WebSocket]
      parameters:
        - name: realm
          in: query
          required: true
          description: |
            The realm name: a lowercase ASCII word starting with a letter, so
            `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm`
            are not realm names — no realm can be created under one. This socket is
            a read path though, and it authenticates before it looks anything up,
            so an off-pattern name is not reported as malformed: it is answered
            `401` `Unauthorized`, the same body as an unknown realm or a bad
            token, because the realm lookup (`GetRealmByName` returning
            `store.ErrNotFound`) is turned into a 401 with no existence oracle.
          schema:
            type: string
            pattern: '^[a-z][a-z0-9]*$'
        - name: token
          in: query
          required: true
          description: Realm JWT for authentication.
          schema:
            type: string
      responses:
        "101":
          description: WebSocket upgrade successful.
        "401":
          $ref: "#/components/responses/Unauthorized"
        "500":
          $ref: "#/components/responses/InternalServerError"

components:
  parameters:
    RealmName:
      name: realm
      in: path
      required: true
      description: |
        The realm name: a lowercase ASCII word starting with a letter, so
        `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm`
        are not realm names — no realm can be created under one. This socket is a
        read path though, and it authenticates before it looks anything up, so an
        off-pattern name is not reported as malformed: it is answered `401`
        `Unauthorized`, the same body as an unknown realm or a bad token, because
        the realm lookup (`GetRealmByName` returning `store.ErrNotFound`) is
        turned into a 401 with no existence oracle.
      schema:
        type: string
        pattern: '^[a-z][a-z0-9]*$'

  schemas:
    HealthStatus:
      type: object
      required: [status]
      properties:
        status:
          type: string
          enum: [ok]
          description: Liveness status.

    ReadinessStatus:
      type: object
      required: [status, checks]
      properties:
        status:
          type: string
          enum: [ok, unavailable]
          description: Overall readiness status.
        checks:
          type: object
          additionalProperties:
            type: string
          description: 'Per-dependency check results ("ok" or "error: ...").'

    DataEnvelopeHealthStatus:
      type: object
      required: [data]
      properties:
        data:
          type: object
          required: [status]
          properties:
            status:
              type: string
              enum: [ok, unhealthy]

    DataEnvelopeVersion:
      type: object
      required: [data]
      properties:
        data:
          type: string
          description: The Astrate build version.

    ErrorDetail:
      type: object
      required: [errors]
      properties:
        errors:
          type: object
          required: [detail]
          properties:
            detail:
              type: string

  responses:
    Unauthorized:
      description: Unauthorized — missing or invalid JWT/token.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Unauthorized

    Forbidden:
      description: Forbidden — the JWT is valid but does not authorize this path.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Forbidden

    InternalServerError:
      description: Internal Server Error — unexpected failure.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Internal server error

$ rg -n 'func |Query\(\)|Header\.Get|Accept|transport|device_id|interface|WriteHeader|Status' internal/appengine/stream/ws.go | sed -n '1,160p'
5:// as JSON; an a_ch JWT guards the endpoint, and device_id/interface query
29:type Bus interface {
40:func NewAPI(bus Bus, mw *auth.Middleware) *API {
45:func (a *API) Mount(mux *http.ServeMux) {
53:	DeviceID  string    `json:"device_id"`
54:	Interface string    `json:"interface,omitempty"`
60:func toWire(ev *engstream.Event) wireEvent {
69:func (a *API) handle(w http.ResponseWriter, r *http.Request) {
72:		DeviceID:  r.URL.Query().Get("device_id"),
73:		Interface: r.URL.Query().Get("interface"),
86:func (a *API) serveWebSocket(w http.ResponseWriter, r *http.Request, events <-chan engstream.Event) {
87:	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
89:		return // Accept already wrote the failure
102:				_ = conn.Close(websocket.StatusGoingAway, "bus closed")
117:func (a *API) serveSSE(w http.ResponseWriter, r *http.Request, events <-chan engstream.Event) {
120:		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
126:	w.WriteHeader(http.StatusOK)
151:func wantsSSE(r *http.Request) bool {
152:	return r.Header.Get("Accept") == "text/event-stream" || r.URL.Query().Get("transport") == "sse"
156:func sseFrame(payload []byte) []byte {

$ wc -l internal/appengine/stream/ws.go internal/appengine/channels/ws.go
162 internal/appengine/stream/ws.go
  442 internal/appengine/channels/ws.go
  604 total

→ Read internal/appengine/stream/ws.go

$ go env GOMODCACHE; rg -n 'StatusText|WriteHeader|badRequest|verifyClientRequest' "$(go env GOMODCACHE)/github.com/coder/websocket@"*/accept.go 2>/dev/null | head -40
/root/go/pkg/mod
109:	errCode, err := verifyClientRequest(w, r)
121:				err = errors.New(http.StatusText(http.StatusForbidden))
131:		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
151:	w.WriteHeader(http.StatusSwitchingProtocols)
154:		WriteHeaderNow()
156:		ginWriter.WriteHeaderNow()
162:		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
184:func verifyClientRequest(w http.ResponseWriter, r *http.Request) (errCode int, _ error) {

$ rg -n 'func \(m \*Middleware\) RequireRealm|Authorization|Bearer|token|ClaimChannels' internal/auth/middleware.go | sed -n '1,60p'
19:// tokenContextKey carries the verified *Token in the request context.
20:type tokenContextKey struct{}
22:// TokenFromContext returns the verified token stored by the middleware, if
25:	tok, ok := ctx.Value(tokenContextKey{}).(*Token)
31:// unverifiable token → 401, verified token whose claims do not authorize the
39:// DefaultCacheSize token cache.
46:// bearer token, and matches the claim's authorization strings against the
48:func (m *Middleware) RequireRealm(claim Claim) func(http.Handler) http.Handler {
57:func (m *Middleware) RequireRealmAny(claims ...Claim) func(http.Handler) http.Handler {
131:	tokenString, ok := bearerToken(r)
137:	tok, err := m.cache.Verify(tokenString, keysPEM)
153:	ctx := context.WithValue(r.Context(), tokenContextKey{}, tok)
175:// bearerToken extracts the credential from the Authorization header.
177:// "Bearer", optionally followed by a colon, then whitespace and the token.
179:	header := r.Header.Get("Authorization")
191:	token := strings.TrimSpace(rest)
192:	if token == "" {
195:	return token, true

$ sed -n '184,240p' /root/go/pkg/mod/github.com/coder/websocket@*/accept.go; echo "=== Accept fn ==="; sed -n '100,135p' /root/go/pkg/mod/github.com/coder/websocket@*/accept.go
func verifyClientRequest(w http.ResponseWriter, r *http.Request) (errCode int, _ error) {
	if !r.ProtoAtLeast(1, 1) {
		return http.StatusUpgradeRequired, fmt.Errorf("WebSocket protocol violation: handshake request must be at least HTTP/1.1: %q", r.Proto)
	}

	if !headerContainsTokenIgnoreCase(r.Header, "Connection", "Upgrade") {
		w.Header().Set("Connection", "Upgrade")
		w.Header().Set("Upgrade", "websocket")
		return http.StatusUpgradeRequired, fmt.Errorf("WebSocket protocol violation: Connection header %q does not contain Upgrade", r.Header.Get("Connection"))
	}

	if !headerContainsTokenIgnoreCase(r.Header, "Upgrade", "websocket") {
		w.Header().Set("Connection", "Upgrade")
		w.Header().Set("Upgrade", "websocket")
		return http.StatusUpgradeRequired, fmt.Errorf("WebSocket protocol violation: Upgrade header %q does not contain websocket", r.Header.Get("Upgrade"))
	}

	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, fmt.Errorf("WebSocket protocol violation: handshake request method is not GET but %q", r.Method)
	}

	if r.Header.Get("Sec-WebSocket-Version") != "13" {
		w.Header().Set("Sec-WebSocket-Version", "13")
		return http.StatusBadRequest, fmt.Errorf("unsupported WebSocket protocol version (only 13 is supported): %q", r.Header.Get("Sec-WebSocket-Version"))
	}

	websocketSecKeys := r.Header.Values("Sec-WebSocket-Key")
	if len(websocketSecKeys) == 0 {
		return http.StatusBadRequest, errors.New("WebSocket protocol violation: missing Sec-WebSocket-Key")
	}

	if len(websocketSecKeys) > 1 {
		return http.StatusBadRequest, errors.New("WebSocket protocol violation: multiple Sec-WebSocket-Key headers")
	}

	// The RFC states to remove any leading or trailing whitespace.
	websocketSecKey := strings.TrimSpace(websocketSecKeys[0])
	if v, err := base64.StdEncoding.DecodeString(websocketSecKey); err != nil || len(v) != 16 {
		return http.StatusBadRequest, fmt.Errorf("WebSocket protocol violation: invalid Sec-WebSocket-Key %q, must be a 16 byte base64 encoded string", websocketSecKey)
	}

	return 0, nil
}

func authenticateOrigin(r *http.Request, originHosts []string) error {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return nil
	}

	u, err := url.Parse(origin)
	if err != nil {
		return fmt.Errorf("failed to parse Origin header %q: %w", origin, err)
	}

	if strings.EqualFold(r.Host, u.Host) {
		return nil
=== Accept fn ===
// Note that using the http.Request Context after Accept returns may lead to
// unexpected behavior (see http.Hijacker).
func Accept(w http.ResponseWriter, r *http.Request, opts *AcceptOptions) (*Conn, error) {
	return accept(w, r, opts)
}

func accept(w http.ResponseWriter, r *http.Request, opts *AcceptOptions) (_ *Conn, err error) {
	defer errd.Wrap(&err, "failed to accept WebSocket connection")

	errCode, err := verifyClientRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), errCode)
		return nil, err
	}

	opts = opts.cloneWithDefaults()
	if !opts.InsecureSkipVerify {
		err = authenticateOrigin(r, opts.OriginPatterns)
		if err != nil {
			if errors.Is(err, path.ErrBadPattern) {
				log.Printf("websocket: %v", err)
				err = errors.New(http.StatusText(http.StatusForbidden))
			}
			http.Error(w, err.Error(), http.StatusForbidden)
			return nil, err
		}
	}

	hj, ok := hijacker(w)
	if !ok {
		err = errors.New("http.ResponseWriter does not implement http.Hijacker")
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return nil, err
	}

	w.Header().Set("Upgrade", "websocket")

$ sed -n '125,200p' internal/auth/middleware.go
authForbidden                   // claim did not grant the request; try next
)

// tryAuthorize is authorize without the terminal 403: it classifies the
// attempt. Callers write their own forbidden response on authForbidden.
func (m *Middleware) tryAuthorize(w http.ResponseWriter, r *http.Request, next http.Handler, claim Claim, keysPEM []string, base string) authResult {
	tokenString, ok := bearerToken(r)
	if !ok {
		_ = astarteapi.WriteUnauthorized(w)
		return authFailed
	}

	tok, err := m.cache.Verify(tokenString, keysPEM)
	if err != nil {
		_ = astarteapi.WriteUnauthorized(w)
		return authFailed
	}

	authPath, ok := RelativePath(r.URL.Path, base)
	if !ok {
		// Upstream parity: a path the authorizer cannot anchor is an
		// authorization failure (403), not an authentication one.
		return authForbidden
	}
	if !tok.Authorizes(claim, r.Method, authPath) {
		return authForbidden
	}

	ctx := context.WithValue(r.Context(), tokenContextKey{}, tok)
	next.ServeHTTP(w, r.WithContext(ctx))
	return authServed
}

// RelativePath computes the authorization path with upstream parity (Astarte's
// GuardianAuthorizePath plug): split the URL path into segments, drop
// everything up to and including the first segment equal to base, and join
// the rest with "/". ok is false when base does not appear in the path.
//
// Example: RelativePath("/pairing/v1/test/agent/devices", "test") returns
// ("agent/devices", true).
func RelativePath(urlPath, base string) (string, bool) {
	segments := strings.Split(strings.Trim(urlPath, "/"), "/")
	for i, s := range segments {
		if s == base {
			return strings.Join(segments[i+1:], "/"), true
		}
	}
	return "", false
}

// bearerToken extracts the credential from the Authorization header.
// Upstream parity (~r/bearer\:?\s+(.*)$/i): the scheme is case-insensitive
// "Bearer", optionally followed by a colon, then whitespace and the token.
func bearerToken(r *http.Request) (string, bool) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", false
	}
	scheme, rest, found := strings.Cut(header, " ")
	if !found {
		return "", false
	}
	scheme = strings.TrimSuffix(scheme, ":")
	if !strings.EqualFold(scheme, "bearer") {
		return "", false
	}
	token := strings.TrimSpace(rest)
	if token == "" {
		return "", false
	}
	return token, true
}

$ sed -n '1,120p' internal/appengine/channels/ws.go
package channels

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

// writeTimeout bounds a single frame write, so one stalled viewer cannot hold
// the session's write mutex indefinitely.
const writeTimeout = 5 * time.Second

// Event-name constants for the watch/unwatch/new_event cycle.
const (
	EventWatch    = "watch"
	EventUnwatch  = "unwatch"
	EventNewEvent = "new_event"
)

// API serves the Phoenix V2 WebSocket endpoint.
type API struct {
	reg   *Registry
	keys  auth.KeySource
	cache *auth.Cache
}

// NewAPI creates an API backed by the given bus, key source and group resolver.
func NewAPI(bus Bus, keys auth.KeySource, groups GroupMembers) *API {
	return &API{
		reg:   NewRegistry(bus, groups),
		keys:  keys,
		cache: auth.NewCache(auth.DefaultCacheSize),
	}
}

// Mount registers the WebSocket handler on the given mux.
func (a *API) Mount(mux *http.ServeMux) {
	mux.Handle("GET /appengine/v1/socket/websocket", http.HandlerFunc(a.handle))
}

// handle upgrades the HTTP connection and runs the session loop.
func (a *API) handle(w http.ResponseWriter, r *http.Request) {
	realm := r.URL.Query().Get("realm")
	token := r.URL.Query().Get("token")

	if realm == "" || token == "" {
		_ = astarteapi.WriteUnauthorized(w)
		return
	}

	row, err := a.keys.GetRealmByName(r.Context(), realm)
	switch {
	case errors.Is(err, store.ErrNotFound):
		_ = astarteapi.WriteUnauthorized(w)
		return
	case err != nil:
		_ = astarteapi.WriteInternalServerError(w)
		return
	}

	tok, err := a.cache.Verify(token, row.JWTPublicKeysPEM)
	if err != nil {
		_ = astarteapi.WriteUnauthorized(w)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return
	}

	s := &session{
		conn:  conn,
		ctx:   r.Context(),
		realm: realm,
		tok:   tok,
		reg:   a.reg,
		rooms: make(map[string]*joined),
	}
	s.loop(r.Context())
}

// session holds per-connection state for the WebSocket session loop.
type session struct {
	conn    *websocket.Conn
	ctx     context.Context
	realm   string
	tok     *auth.Token
	reg     *Registry
	mu      sync.Mutex
	roomsMu sync.Mutex
	rooms   map[string]*joined
}

// joined tracks a topic this session has joined.
type joined struct {
	room    *Room
	member  *Member
	joinRef *string
}

// triggerShape decodes the fields needed to derive the WATCH authorization path.
type triggerShape struct {
	Type          string `json:"type"`
	InterfaceName string `json:"interface_name"`
	MatchPath     string `json:"match_path"`
	GroupName     string `json:"group_name"`
	DeviceID      string `json:"device_id"`
}

$ rg -n 'WriteHeader|Status[A-Z]+|http\.Error' internal/appengine/channels/*.go
internal/appengine/channels/ws_test.go:65:	if rec.Code != http.StatusUnauthorized {
internal/appengine/channels/ws_test.go:78:	if rec.Code != http.StatusUnauthorized {
internal/appengine/channels/ws_test.go:91:	if rec.Code != http.StatusUnauthorized {
internal/appengine/channels/ws_test.go:103:	if rec.Code != http.StatusUnauthorized {
internal/appengine/channels/ws_test.go:116:	if rec.Code == http.StatusUnauthorized {

$ sed -n '28,125p' internal/auth/middleware.go
// Middleware authenticates and authorizes REST requests with realm JWTs
// (docs/DESIGN.md §4.2). Status mapping is upstream parity: missing or
// unverifiable token → 401, verified token whose claims do not authorize the
// request → 403, both with the canonical envelopes.
type Middleware struct {
	keys  KeySource
	cache *Cache
}

// NewMiddleware builds a Middleware over the given key source with a
// DefaultCacheSize token cache.
func NewMiddleware(keys KeySource) *Middleware {
	return &Middleware{keys: keys, cache: NewCache(DefaultCacheSize)}
}

// RequireRealm guards a realm-scoped route (path pattern must carry a
// {realm} segment): it resolves the realm's JWT public keys, verifies the
// bearer token, and matches the claim's authorization strings against the
// method and the path relative to the realm base.
func (m *Middleware) RequireRealm(claim Claim) func(http.Handler) http.Handler {
	return m.RequireRealmAny(claim)
}

// RequireRealmAny guards a realm-scoped route accepting ANY of the given
// claims (OR-ed): the first claim whose authorization strings match the
// method/path grants access. Used by surfaces that honour both their upstream
// claim and an Astrate compatibility claim — e.g. Flow accepts a_f (upstream)
// and a_rma (Astrate's original operator claim).
func (m *Middleware) RequireRealmAny(claims ...Claim) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			realm := r.PathValue("realm")
			if realm == "" || len(claims) == 0 {
				_ = astarteapi.WriteUnauthorized(w)
				return
			}

			row, err := m.keys.GetRealmByName(r.Context(), realm)
			switch {
			case errors.Is(err, store.ErrNotFound):
				// An unknown realm has no keys: unauthenticated, not 404
				// (no existence oracle on auth failures).
				_ = astarteapi.WriteUnauthorized(w)
				return
			case err != nil:
				_ = astarteapi.WriteInternalServerError(w)
				return
			}

			for _, claim := range claims {
				switch m.tryAuthorize(w, r, next, claim, row.JWTPublicKeysPEM, realm) {
				case authServed:
					return
				case authFailed:
					// Authentication itself failed (401 already written):
					// no other claim can succeed.
					return
				}
				// authForbidden: this claim did not grant the request; try
				// the next one.
			}
			_ = astarteapi.WriteForbidden(w)
		})
	}
}

// RequireStatic guards an instance-level route (Housekeeping) with a fixed
// key set instead of per-realm keys. The authorization path is the request
// path relative to the service base (the segment after "v1").
func (m *Middleware) RequireStatic(claim Claim, keysPEM []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.authorize(w, r, next, claim, keysPEM, "v1")
		})
	}
}

// authorize runs the shared bearer-extract → verify → claim-match pipeline,
// writing the error response itself. base is the path segment after which the
// authorization path starts (the realm name, or "v1" for instance-level
// routes).
func (m *Middleware) authorize(w http.ResponseWriter, r *http.Request, next http.Handler, claim Claim, keysPEM []string, base string) {
	switch m.tryAuthorize(w, r, next, claim, keysPEM, base) {
	case authServed, authFailed:
		// Served, or authentication failed with the 401 already written.
	default:
		_ = astarteapi.WriteForbidden(w)
	}
}

// authResult classifies one claim attempt by tryAuthorize.
type authResult int

const (
	authServed    authResult = iota // request authorized and next served
	authFailed                      // authentication failed (401 written); stop
	authForbidden                   // claim did not grant the request; try next

$ rg -n 'docs-sync-native|docs-native' .mule/todo.md | wc -l; rg -n '^\- \[.\] docs-(sync-)?native' .mule/todo.md
13
154:- [x] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmmanagement/health`, `GET /pairing/health` — in docs/api/astrate_native_api.yaml, and widen the shared `DataEnvelopeHealthStatus.status` enum (yaml:301-303) from `[ok]` to `[ok, unhealthy]` — the handler answers 503 `{"data":{"status":"unhealthy"}}` when the readiness probe fails (internal/observability/compat.go:19-34, deviation 18), but the spec documents only 200 with `status: ok` for all three paths. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
157:- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
185:- [x] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication is "via Authorization header or `?token=` query parameter" (yaml:293-294), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
196:- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
201:- [x] docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate_devices_total{realm="test"} 42`, but no such series is ever registered — the registry carries `astrate_broker_sessions` (internal/observability/metrics.go:54-62), `astrate_db_pool_{acquired,idle,total,max}_conns` (metrics.go:65-79), plus the engine/flow/trigger families (`astrate_engine_*` engine/router.go:411-444, `astrate_engine_stream_dropped_total` engine/stream/bus.go:119, `astrate_engine_trigger_*` engine/triggers/actions.go:296-300, `astrate_flow_router_*` flow/router.go:260-276) and the standard `go_*`/`process_*` runtime collectors — so the example advertises a scrape result the wire can never return. Swap it for a real series (e.g. `astrate_broker_sessions 0`), keeping the HELP/TYPE comment style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
202:- [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
262:- [x] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `"Internal server error"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
277:- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
278:- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
279:- [x] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
280:- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
281:- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ rg -n '426|Upgrade Required|Sec-WebSocket|Connection: Upgrade' .mule/todo.md .mule/for-giulio.md docs/api/ || echo NONE
NONE

$ rg -n 'transport|text/event-stream|SSE' docs/api/astrate_native_api.yaml
287:      summary: Native live event stream (WebSocket / SSE)
291:        `Accept: text/event-stream` or `?transport=sse`.
312:        **SSE protocol:** each event is a `data:` frame followed by a blank line.
328:        - name: transport
334:          description: Force SSE transport instead of WebSocket.
339:          description: "SSE stream started (when `Accept: text/event-stream`)."

$ mkdir -p /tmp/opencode/probe && cat > /tmp/opencode/probe/main.go <<'EOF'
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/coder/websocket"
)

func main() {
	r := httptest.NewRequest(http.MethodGet, "/astrate/v1/test/socket", nil)
	w := httptest.NewRecorder()
	_, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	res := w.Result()
	fmt.Printf("plain GET: status=%d contentType=%q body=%q err=%v\n", res.StatusCode, res.Header.Get("Content-Type"), w.Body.String(), err)

	r2 := httptest.NewRequest(http.MethodGet, "/astrate/v1/test/socket", nil)
	r2.Header.Set("Connection", "Upgrade")
	r2.Header.Set("Upgrade", "websocket")
	r2.Header.Set("Sec-WebSocket-Version", "14")
	r2.Header.Set("Sec-WebSocket-Key", "AAAAAAAAAAAAAAAAAAAAAA==")
	w2 := httptest.NewRecorder()
	_, err = websocket.Accept(w2, r2, &websocket.AcceptOptions{InsecureSkipVerify: true})
	res2 := w2.Result()
	fmt.Printf("bad version: status=%d body=%q err=%v\n", res2.StatusCode, w2.Body.String(), err)
}
EOF
go run /tmp/opencode/probe/main.go
plain GET: status=426 contentType="text/plain; charset=utf-8" body="WebSocket protocol violation: Connection header \"\" does not contain Upgrade\n" err=failed to accept WebSocket connection: WebSocket protocol violation: Connection header "" does not contain Upgrade
bad version: status=400 body="unsupported WebSocket protocol version (only 13 is supported): \"14\"\n" err=failed to accept WebSocket connection: unsupported WebSocket protocol version (only 13 is supported): "14"

$ cat > /tmp/opencode/probe/main2.go <<'EOF'
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

type wireEvent struct {
	Event     string    `json:"event"`
	Realm     string    `json:"realm"`
	DeviceID  string    `json:"device_id"`
	Interface string    `json:"interface,omitempty"`
	Path      string    `json:"path,omitempty"`
	Value     any       `json:"value,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func wantsSSE(r *http.Request) bool {
	return r.Header.Get("Accept") == "text/event-stream" || r.URL.Query().Get("transport") == "sse"
}

func main() {
	ts := time.Date(2024, 5, 30, 13, 49, 57, 45000000, time.UTC)
	for _, c := range []struct{ name, accept, q string }{
		{"exact", "text/event-stream", ""},
		{"list", "text/event-stream, text/plain", ""},
		{"qvalue", "text/event-stream;q=1.0", ""},
		{"none", "", ""},
		{"query-sse", "", "transport=sse"},
		{"query-other", "", "transport=websocket"},
		{"query-empty", "", "transport="},
	} {
		r := httptest.NewRequest(http.MethodGet, "/x?"+c.q, nil)
		if c.accept != "" {
			r.Header.Set("Accept", c.accept)
		}
		fmt.Printf("sse %-13s accept=%-32q -> %v\n", c.name, c.accept, wantsSSE(r))
	}

	b1, _ := json.Marshal(wireEvent{Event: "device_connected", Realm: "test", DeviceID: "dev1", Timestamp: ts})
	fmt.Println("lifecycle:", string(b1))
	b2, _ := json.Marshal(wireEvent{Event: "incoming_data", Realm: "test", DeviceID: "dev1", Interface: "org.ex.S", Path: "/p", Value: false, Timestamp: ts})
	fmt.Println("value false:", string(b2))
	b3, _ := json.Marshal(wireEvent{Event: "incoming_data", Realm: "test", DeviceID: "dev1", Interface: "org.ex.S", Path: "/p", Value: 0, Timestamp: ts})
	fmt.Println("value 0:", string(b3))
}
EOF
go run /tmp/opencode/probe/main2.go
sse exact         accept="text/event-stream"              -> true
sse list          accept="text/event-stream, text/plain"  -> false
sse qvalue        accept="text/event-stream;q=1.0"        -> false
sse none          accept=""                               -> false
sse query-sse     accept=""                               -> true
sse query-other   accept=""                               -> false
sse query-empty   accept=""                               -> false
lifecycle: {"event":"device_connected","realm":"test","device_id":"dev1","timestamp":"2024-05-30T13:49:57.045Z"}
value false: {"event":"incoming_data","realm":"test","device_id":"dev1","interface":"org.ex.S","path":"/p","value":false,"timestamp":"2024-05-30T13:49:57.045Z"}
value 0: {"event":"incoming_data","realm":"test","device_id":"dev1","interface":"org.ex.S","path":"/p","value":0,"timestamp":"2024-05-30T13:49:57.045Z"}

$ rg -n 'Kind:' internal/engine/ --glob '!*_test.go' | head -30; echo ---; rg -n 'Kind\s+string' internal/engine/stream/bus.go; sed -n '35,70p' internal/engine/stream/bus.go
internal/engine/engine.go:208:		Kind:           stream.KindIncomingData,
internal/engine/engine.go:310:		Kind:          stream.KindDeviceError,
internal/engine/engine.go:351:			Kind: stream.KindDeviceConnected, Realm: rs.name, DeviceID: deviceID, Timestamp: ev.At, IP: ip,
internal/engine/engine.go:369:			Kind: stream.KindDeviceDisconnected, Realm: rs.name, DeviceID: deviceID, Timestamp: ev.At,
internal/engine/engine.go:390:		Kind: stream.KindDeviceRegistered, Realm: rs.name, DeviceID: hwID, Timestamp: at,
internal/engine/engine.go:424:		Kind: kind, Realm: rs.name, DeviceID: hwID, Timestamp: at,
internal/engine/serverdata.go:260:		Kind:      OpPropertyUnset,
internal/engine/triggers/policy.go:195:			return compiledHandler{onKind: onKeyword, keyword: keyword}, nil
internal/engine/triggers/policy.go:219:	return compiledHandler{onKind: onExplicit, codes: codes}, nil
---
42:	Kind string
	// whether or not the realm declares any trigger.
	KindDeviceError = "device_error"
)

// Event is one live event.
type Event struct {
	// Kind discriminates the event (Kind* constants).
	Kind string
	// Realm is the tenant.
	Realm string
	// DeviceID is the encoded device ID.
	DeviceID string
	// Interface and Path locate data events; empty for lifecycle events.
	Interface string
	Path      string
	// InterfaceMajor is the data event's interface major version. Consumers
	// that rebuild upstream trigger-event bodies need it; Filter ignores it.
	InterfaceMajor int
	// IP is the peer address of a device_connected event, empty otherwise.
	IP string
	// ErrorName and ErrorMetadata describe a device_error event: the §2.6
	// reject-reason label and its free-form diagnostic strings.
	ErrorName     string
	ErrorMetadata map[string]string
	// Value is the JSON-friendly rendering of a data event's value (nil for
	// property unset and lifecycle events).
	Value any
	// Timestamp is the event instant (the effective sample timestamp for
	// data events).
	Timestamp time.Time
}

// Filter narrows a subscription; zero-value fields match everything.
type Filter struct {
	// DeviceID keeps only one device's events when set.
	DeviceID string

$ rg -n 'device_connected|device_disconnected|incoming_data|value_change' internal/engine/stream/*.go internal/engine/*.go | rg -v _test | head -30
internal/engine/cache.go:490:// store. Wired to device_disconnected lifecycle events.
internal/engine/stream/bus.go:20:	KindIncomingData = "incoming_data"
internal/engine/stream/bus.go:24:	KindDeviceConnected = "device_connected"
internal/engine/stream/bus.go:26:	KindDeviceDisconnected = "device_disconnected"
internal/engine/stream/bus.go:53:	// IP is the peer address of a device_connected event, empty otherwise.
internal/engine/engine.go:223://   - value_change fires when the value differs from the previous one;
internal/engine/engine.go:226://   - value_change_applied mirrors value_change (same comparison, same seam).
internal/engine/engine.go:328:// eviction): device_connected / device_disconnected triggers, live bus
internal/engine/engine.go:376:// event on the live bus — the same path device_connected already uses.
internal/engine/data.go:135:	// events but still fire incoming_data).
internal/engine/data.go:355:// for the change-derived conditions (value_change*, path_created/removed).
internal/engine/data.go:359:// Failures degrade to untracked — incoming_data still fires, change events
internal/engine/data.go:451:// "previous". Change-derived conditions (value_change*, path_created/removed)
internal/engine/serverdata.go:250:	// change-derived conditions (path_removed, value_change) evaluate against

$ rg -n 'MountServiceCompat|MountVersionCompat|Mount\(' cmd/astrate/main.go | sed -n '1,60p'
362:	}).Mount(mux)
371:	housekeeping.NewAPI(hkSvc, mw, hkKeys).Mount(mux)
375:	realm.NewAPI(realmSvc, mw).Mount(mux)
376:	appengine.NewAPI(appengine.NewService(st, e, log), mw).Mount(mux)
377:	apstream.NewAPI(e.Bus(), mw).Mount(mux)
380:	channels.NewAPI(e.Bus(), st, st).Mount(mux)
382:	flowapi.NewAPI(flowSvc, mw).Mount(mux)
387:		observability.MountServiceCompat(mux, svc, st.Health)
397:		observability.MountVersionCompat(mux, svc, version)
414:	health.Mount(mux)
417:	swagger.Mount(mux)

$ sed -n '380,450p' cmd/astrate/main.go
channels.NewAPI(e.Bus(), st, st).Mount(mux)
	// Flow operator API: pipelines CRUD + start/stop/status (v2.0).
	flowapi.NewAPI(flowSvc, mw).Mount(mux)

	// Upstream-parity per-service health endpoints (the dashboard's API
	// status indicators poll them).
	for _, svc := range []string{"appengine", "realmmanagement", "pairing"} {
		observability.MountServiceCompat(mux, svc, st.Health)
	}

	// Upstream-parity per-service version endpoints (issue #77): every
	// service answers an unauthenticated GET /{service}/version with
	// {"data": version}; AppEngine's realm-scoped variant requires auth,
	// Pairing's is public — both measured on upstream 1.2.0
	// (test/conformance/upstream/verify-versions.json). Realm Management's
	// realm-scoped route is served by its own API.
	for _, svc := range []string{"appengine", "realmmanagement", "pairing", "housekeeping"} {
		observability.MountVersionCompat(mux, svc, version)
	}
	mountRealmVersion(mux, mw)

	metrics.RegisterBrokerSessions(func() float64 { return float64(b.SessionCount()) })
	metrics.RegisterDBPool(func() observability.DBPoolStats {
		s := st.Stat()
		return observability.DBPoolStats{
			AcquiredConns: s.AcquiredConns(),
			IdleConns:     s.IdleConns(),
			TotalConns:    s.TotalConns(),
			MaxConns:      s.MaxConns(),
		}
	})
	health := observability.NewHealth(metrics.Handler())
	health.AddReadiness("database", st.Health)
	health.AddReadiness("broker", brokerReadiness(b))
	health.Mount(mux)

	// Swagger UI + OpenAPI YAML specs embedded in the binary.
	swagger.Mount(mux)

	// Unmatched routes under a service prefix answer that service's JSON error
	// envelope rather than Go's plain-text "404 page not found", which any
	// client parsing the envelope would choke on.
	handler := httpx.NotFound(mux)
	if len(cfg.HTTP.CORSAllowedOrigins) > 0 {
		handler = httpx.CORS(cfg.HTTP.CORSAllowedOrigins)(handler)
	}
	return handler, hkSvc, nil
}

// mountRealmVersion registers the realm-scoped version endpoints observed on
// upstream 1.2.0 (test/conformance/upstream/verify-versions.json): AppEngine's
// requires a realm token, pairing's is public. Both report the emulated
// upstream API level, the value the Astarte Dashboard feature-gates its UI on
// (COMPATIBILITY.md deviation 10) — not Astrate's build version. RM's
// realm-scoped route is served by its own API.
func mountRealmVersion(mux *http.ServeMux, mw *auth.Middleware) {
	mux.Handle("GET /appengine/v1/{realm}/version",
		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(realm.APICompatVersion)))
	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
}

// brokerReadiness reports the broker listener as ready when a TCP connection to
// it is accepted (docs/DESIGN.md §5.2 readiness broker check).
func brokerReadiness(b *broker.Broker) observability.Check {
	return func(ctx context.Context) error {
		addr := b.TLSAddr()
		if addr == "" {
			return errors.New("broker listener not bound")
		}
		var d net.Dialer
		conn, err := d.DialContext(ctx, "tcp", addr)

$ rg -n 'vsn' internal/appengine/channels/ internal/appengine/ --glob '!*_test.go' | head; echo "--- protocol"; sed -n '1,60p' internal/appengine/channels/protocol.go
--- protocol
// Package channels implements the upstream Phoenix Channels V2 wire protocol
// served at /appengine/v1/socket/websocket. It is the compatibility counterpart
// to the Astrate-native socket in internal/appengine/stream.
package channels

import (
	"encoding/json"
	"fmt"
)

// Event-name constants used by the Phoenix V2 protocol.
const (
	EventPhxJoin   = "phx_join"
	EventPhxReply  = "phx_reply"
	EventPhxLeave  = "phx_leave"
	EventPhxClose  = "phx_close"
	EventPhxError  = "phx_error"
	EventHeartbeat = "heartbeat"

	// TopicHeartbeat is the reserved topic for heartbeat messages.
	TopicHeartbeat = "phoenix"
)

// Frame is one Phoenix V2 message on the wire, encoded as a five-element JSON
// array: [join_ref, ref, topic, event, payload]. A nil JoinRef or Ref is
// marshalled as JSON null; the two are distinguished from the empty string.
type Frame struct {
	JoinRef *string
	Ref     *string
	Topic   string
	Event   string
	Payload json.RawMessage
}

// MarshalJSON renders the frame as a five-element JSON array. A nil JoinRef or
// Ref becomes JSON null. A nil or zero-length Payload becomes {}.
func (f Frame) MarshalJSON() ([]byte, error) {
	payload := json.RawMessage(`{}`)
	if len(f.Payload) > 0 {
		payload = f.Payload
	}
	return json.Marshal([5]interface{}{f.JoinRef, f.Ref, f.Topic, f.Event, payload})
}

// UnmarshalJSON parses a five-element JSON array back into the frame. It
// rejects anything that is not a JSON array of exactly five elements.
func (f *Frame) UnmarshalJSON(b []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("channels.Frame: invalid JSON: %w", err)
	}
	if len(raw) != 5 {
		return fmt.Errorf("channels.Frame: expected a 5-element JSON array, got %d elements", len(raw))
	}
	// Decode join_ref and ref as nullable strings.
	jr, err := unmarshalNullString(raw[0])
	if err != nil {
		return fmt.Errorf("channels.Frame: join_ref: %w", err)
	}
	f.JoinRef = jr

$ rg -n 'securitySchemes|security:' docs/api/astrate_native_api.yaml; echo "=== appengine schemes ==="; rg -n 'securitySchemes' -A 20 docs/api/astarte_appengine_api.yaml | sed -n '1,40p'
21:security: []
=== appengine schemes ===
1436:  securitySchemes:
1437-    a_aea:
1438-      type: apiKey
1439-      in: header
1440-      name: Authorization
1441-      description: |
1442-        Realm-level JWT carrying the `a_aea` claim. Used to authenticate
1443-        AppEngine operations (device data, groups).
1444-
1445-  parameters:
1446-    RealmName:
1447-      name: realm
1448-      in: path
1449-      required: true
1450-      description: The realm name.
1451-      schema:
1452-        type: string
1453-
1454-    DeviceID:
1455-      name: device
1456-      in: path

$ grep -n 'responses:\|        "101"\|        "200"\|        "401"\|        "403"\|        "500"\|Accept: text/event-stream\|transport\|components:\|  schemas:\|  responses:' docs/api/astrate_native_api.yaml | sed -n '1,80p'
34:      responses:
35:        "200":
53:      responses:
54:        "200":
96:      responses:
97:        "200":
122:      responses:
123:        "200":
150:      responses:
151:        "200":
178:      responses:
179:        "200":
210:      responses:
211:        "200":
230:      responses:
231:        "200":
250:      responses:
251:        "200":
272:      responses:
273:        "200":
291:        `Accept: text/event-stream` or `?transport=sse`.
328:        - name: transport
334:          description: Force SSE transport instead of WebSocket.
335:      responses:
336:        "101":
338:        "200":
339:          description: "SSE stream started (when `Accept: text/event-stream`)."
340:        "401":
342:        "403":
344:        "500":
398:      responses:
399:        "101":
401:        "401":
403:        "500":
406:components:
425:  schemas:
480:  responses:

$ sed -n '336,360p' docs/api/astrate_native_api.yaml; echo "=== 400-470 ==="; sed -n '400,440p' docs/api/astrate_native_api.yaml
"101":
          description: WebSocket upgrade successful.
        "200":
          description: "SSE stream started (when `Accept: text/event-stream`)."
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "500":
          $ref: "#/components/responses/InternalServerError"

  # ── Phoenix V2 compat WebSocket ─────────────────────────────────────

  /appengine/v1/socket/websocket:
    get:
      operationId: phoenixWebSocket
      summary: Phoenix V2-compatible WebSocket endpoint
      description: |
        Upstream-compatible Phoenix Channels WebSocket. Authenticates via
        `?realm=` and `?token=` query parameters (realm JWT).

        Uses the Phoenix V2 frame protocol:
        - Client sends `phx_join` to join a room (`rooms:{realm}:{name}`)
        - Client sends `heartbeat` for keepalive
        - Client sends `watch` / `unwatch` to manage triggers
=== 400-470 ===
          description: WebSocket upgrade successful.
        "401":
          $ref: "#/components/responses/Unauthorized"
        "500":
          $ref: "#/components/responses/InternalServerError"

components:
  parameters:
    RealmName:
      name: realm
      in: path
      required: true
      description: |
        The realm name: a lowercase ASCII word starting with a letter, so
        `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm`
        are not realm names — no realm can be created under one. This socket is a
        read path though, and it authenticates before it looks anything up, so an
        off-pattern name is not reported as malformed: it is answered `401`
        `Unauthorized`, the same body as an unknown realm or a bad token, because
        the realm lookup (`GetRealmByName` returning `store.ErrNotFound`) is
        turned into a 401 with no existence oracle.
      schema:
        type: string
        pattern: '^[a-z][a-z0-9]*$'

  schemas:
    HealthStatus:
      type: object
      required: [status]
      properties:
        status:
          type: string
          enum: [ok]
          description: Liveness status.

    ReadinessStatus:
      type: object
      required: [status, checks]
      properties:
        status:
          type: string

$ sed -n '1,40p' docs/Makefile
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

$ rg -n 'astrate_native|TestNative' internal/swagger/swagger_test.go | head -20
616:// TestNativeErrorDetailExamples pins the error-detail examples in the native
631:func TestNativeErrorDetailExamples(t *testing.T) {
632:	b, err := docs.APIYAML.ReadFile("api/astrate_native_api.yaml")
634:		t.Fatalf("reading astrate_native_api.yaml: %v", err)
668:// TestNativeVersionDescriptions guards the four compat-version operations of the
682:func TestNativeVersionDescriptions(t *testing.T) {
683:	b, err := docs.APIYAML.ReadFile("api/astrate_native_api.yaml")
685:		t.Fatalf("reading astrate_native_api.yaml: %v", err)
716:// TestNativeMetricsContentNegotiationDocumented guards that the native spec's
734:func TestNativeMetricsContentNegotiationDocumented(t *testing.T) {
735:	b, err := docs.APIYAML.ReadFile("api/astrate_native_api.yaml")
737:		t.Fatalf("reading astrate_native_api.yaml: %v", err)
846:// TestAppEngineErrorDetailExamples and TestNativeErrorDetailExamples. This is

$ rg -ln 'astrate/v1/.*/socket|socket/websocket|event-stream' docs/site/ | head; echo "---"; rg -n 'socket' docs/site/*.md | head -20
docs/site/appengine-api.md
docs/site/pokemon-agent.md
docs/site/compatibility.md
---
docs/site/operations.md:110:key. The Device Live Events card requires the Channels socket (planned M11);
docs/site/migration-from-astarte.md:33:| Channels (Phoenix socket) | WebSocket/SSE endpoint (additive) |
docs/site/migration-from-astarte.md:144:- **Dashboard:** The upstream Astarte Dashboard works against Astrate, but Device Live Events requires the Channels socket (see [Compatibility](compatibility.md) for details).
docs/site/compatibility.md:27:- Astarte Dashboard v1.2.2 Device Live Events (since M11 Channels socket)
docs/site/compatibility.md:33:### 1. Astarte Channels: two sockets, one bus
docs/site/compatibility.md:35:The upstream Phoenix socket is served at `/appengine/v1/socket/websocket` (phoenix.js V2 wire format) for Dashboard compatibility. Astrate keeps its own plain WebSocket/SSE endpoint at `/astrate/v1/<realm>/socket`.
docs/site/compatibility.md:51:Health, readiness, metrics, and the live-stream socket under `/astrate/v1/...` -- a namespace that cannot collide with upstream.
docs/site/ROADMAP.md:64:| **M7** | `internal/realm`, `internal/housekeeping`, `internal/appengine` (+ stream socket) | T2/T3 green; **CP-C** |
docs/site/ROADMAP.md:98:| WebSocket | `github.com/coder/websocket` | Minimal, context-native |
docs/site/ROADMAP.md:355:**M7a** (realm + housekeeping) and **M7b** (appengine + stream socket).
docs/site/ROADMAP.md:374:| 7.9 | `internal/appengine/stream/ws.go` | `/astrate/v1/{realm}/socket` WebSocket + SSE fallback fed by `engine/stream.Bus`; `a_ch` claims honoured as room filters (§1.1 deviation: not Phoenix-wire) | 180 |
docs/site/ROADMAP.md:380:- **T3 cross-domain:** PATCH `credentials_inhibited=true` ⇒ broker rejects next CONNECT and pairing rejects credentials. Server-owned PUT ⇒ connected test device receives it (correct QoS + format); offline device receives on reconnect (retention path). Stream socket: test device publishes ⇒ WebSocket subscriber receives the event JSON.
docs/site/DESIGN.md:35:criterion), and Astarte Channels' full Phoenix-socket protocol (we provide a simpler
docs/site/DESIGN.md:63:| **Astarte Channels** | Phoenix WebSocket rooms for live data | `internal/appengine/stream` | Simplified WebSocket + SSE endpoint fed by the engine's fan-out bus. Not wire-compatible with Phoenix sockets in v1 (documented deviation). |
docs/site/DESIGN.md:603:1. **Astarte Channels** Phoenix-socket protocol is replaced by plain WebSocket/SSE (different
docs/site/DESIGN.md:627:- Astrate-native additions under `/astrate/v1/...` (health, metrics, live stream socket) to
docs/site/DESIGN.md:652:  (Housekeeping), `a_ch` (Channels → honoured by Astrate's stream socket). Each is a list of
docs/site/DESIGN.md:658:- **`a_ch` is the exception and does not follow the rule above.** Upstream's Channels socket
docs/site/appengine-api.md:99:GET /astrate/v1/<realm>/socket
docs/site/appengine-api.md:107:GET /appengine/v1/socket/websocket?realm=<realm>&token=<jwt>

$ sed -n '85,140p' docs/site/appengine-api.md
GET    /appengine/v1/<realm>/groups
POST   /appengine/v1/<realm>/groups
GET    /appengine/v1/<realm>/groups/<name>

GET    /appengine/v1/<realm>/groups/<name>/devices
POST   /appengine/v1/<realm>/groups/<name>/devices
DELETE /appengine/v1/<realm>/groups/<name>/devices/<device_id>
```

## Live stream

### WebSocket/SSE endpoint

```
GET /astrate/v1/<realm>/socket
```

Feeds real-time device events to consumers. Honours `a_ch` claims as room filters.

### Astarte Channels (Dashboard compatibility)

```
GET /appengine/v1/socket/websocket?realm=<realm>&token=<jwt>
```

Phoenix WebSocket V2 wire format for upstream Astarte Dashboard compatibility. Answers `phx_join`, `watch`, `phx_leave`, and the `phoenix` heartbeat. Events pushed as `new_event`.

## Response envelope

All responses use the Astarte envelope format:

```json
{ "data": <value> }
```

Errors:

```json
{ "errors": { "detail": "<message>" } }
```

Status codes match upstream: 401 (no/bad token), 403 (claim mismatch), 404 (not found), 409/422 (conflict/validation).

$ sed -n '1,60p' docs/site/appengine-api.md; echo "=== DESIGN 645-670 ==="; sed -n '645,672p' docs/DESIGN.md
# AppEngine API

The AppEngine API provides REST endpoints for reading device data, publishing server-owned data, managing groups and aliases, and a live-stream WebSocket/SSE endpoint.

**Base path:** `/appengine/v1/<realm>/`
**Authentication:** JWT with `a_aea` claim.

## Device management

### List devices

```
GET /appengine/v1/<realm>/devices?details=true&limit=N&from_token=<cursor>
```

Returns a paginated list. Pagination uses `body.links.next` cursor (upstream parity).

### Get device

```
GET /appengine/v1/<realm>/devices/<device_id>
```

Returns device status: introspection, connected flag, stats (total received msgs/bytes), timestamps, aliases, attributes.

### Update device

```
PATCH /appengine/v1/<realm>/devices/<device_id>
```

Fields: `aliases`, `attributes`, `credentials_inhibited`.

Setting `credentials_inhibited: true` blocks the device from obtaining new credentials and connecting to the broker.

### Devices by alias

```
GET /appengine/v1/<realm>/devices-by-alias/<alias>
```

## Device data

### Properties

```
GET /appengine/v1/<realm>/devices/<device_id>/interfaces/<interface>[/<path>]
```

Returns the properties snapshot tree for a device's interfaces.

### Datastream queries

```
GET /appengine/v1/<realm>/devices/<device_id>/interfaces/<interface>/<path>
    ?since=<timestamp>
    &since_after=<timestamp>
    &to=<timestamp>
    &limit=N
    &downsample_to=<bucket_duration>
=== DESIGN 645-670 ===
   device-generated CSR; the only credential the broker accepts.

### 4.2 JWT validation & Astarte authorization claims (`internal/auth`)

Astarte's claim model is reproduced exactly so existing tokens/tooling (`astartectl`) work:

- Claims: `a_aea` (AppEngine), `a_rma` (Realm Management), `a_pa` (Pairing), `a_ha`
  (Housekeeping), `a_ch` (Channels → honoured by Astrate's stream socket). Each is a list of
  authorization strings `"<HTTP-verb-regex>::<path-regex>"` (e.g. `"^POST$::^devices/.*$"`,
  or the catch-all `".*::.*"`), matched against the method and the path *relative to the realm
  base* (e.g. `devices/h4-Dx_RYTU-RbpDOTabhRg/interfaces/...`).
- Regexes are compiled once per token (LRU cache keyed by token hash) and are implicitly
  anchored as upstream does; `exp` honoured if present; `iat` not required (parity).
- **`a_ch` is the exception and does not follow the rule above.** Upstream's Channels socket
  matches the verb field *literally* against `JOIN` or `WATCH` and discards every entry that
  is neither, so the catch-all `".*::.*"` authorizes nothing there. `Token.AuthorizesChannel`
  implements that rule; `Token.Authorizes` keeps the verb-regex behaviour for the REST
  surfaces. Measured, not inferred — see `test/conformance/upstream/channels.json` and
  COMPATIBILITY.md deviation 1.
- Multiple realm public keys allow zero-downtime key rotation (`PUT
  /realmmanagement/v1/<realm>/config/auth` parity endpoint).

### 4.3 Embedded per-realm CA (`internal/pairing/ca`)

Replaces CFSSL:

- On realm creation, generate an ECDSA P-256 CA key + self-signed CA cert (configurable
  lifetime, default 10 years), or import an operator-provided pair. Private key encrypted at

$ rg -n 'AuthorizesChannel|Authorizes\(' internal/ --glob '!*_test.go'
internal/auth/middleware.go:149:	if !tok.Authorizes(claim, r.Method, authPath) {
internal/auth/claims.go:81:// AuthorizesChannel.
internal/auth/claims.go:164:func (t *Token) Authorizes(claim Claim, verb, authPath string) bool {
internal/auth/claims.go:180:// AuthorizesChannel reports whether the a_ch claim grants verb (VerbJoin or
internal/auth/claims.go:197:func (t *Token) AuthorizesChannel(verb, authPath string) bool {
internal/appengine/channels/ws.go:278:	if !s.tok.AuthorizesChannel(auth.VerbJoin, name) {
internal/appengine/channels/ws.go:354:	if !s.tok.AuthorizesChannel(auth.VerbWatch, path) {

$ rg -n 'func \(t \*Token\) Authorizes' -A 30 internal/auth/*.go | sed -n '1,80p'
internal/auth/claims.go:164:func (t *Token) Authorizes(claim Claim, verb, authPath string) bool {
internal/auth/claims.go-165-	for _, g := range t.grants[claim] {
internal/auth/claims.go-166-		if g.verb.MatchString(verb) && g.path.MatchString(authPath) {
internal/auth/claims.go-167-			return true
internal/auth/claims.go-168-		}
internal/auth/claims.go-169-	}
internal/auth/claims.go-170-	return false
internal/auth/claims.go-171-}
internal/auth/claims.go-172-
internal/auth/claims.go-173-// Channels authorization verbs. Upstream recognises exactly these two and
internal/auth/claims.go-174-// silently discards an a_ch entry carrying anything else.
internal/auth/claims.go-175-const (
internal/auth/claims.go-176-	VerbJoin  = "JOIN"
internal/auth/claims.go-177-	VerbWatch = "WATCH"
internal/auth/claims.go-178-)
internal/auth/claims.go-179-
internal/auth/claims.go-180-// AuthorizesChannel reports whether the a_ch claim grants verb (VerbJoin or
internal/auth/claims.go-181-// VerbWatch) on authPath — the room name for a join, the trigger's target for
internal/auth/claims.go-182-// a watch.
internal/auth/claims.go-183-//
internal/auth/claims.go-184-// It deliberately does not reuse Authorizes, because upstream reads a_ch by a
internal/auth/claims.go-185-// different rule than the REST claims and the difference is observable. The
internal/auth/claims.go-186-// REST plug compiles the verb field into a regex and matches it against the
internal/auth/claims.go-187-// HTTP method, so "GET|POST" works and ".*::.*" grants everything. The
internal/auth/claims.go-188-// Channels socket instead *partitions* the a_ch list by an exact match on the
internal/auth/claims.go-189-// verb field, keeping only entries whose first field is literally "JOIN" or
internal/auth/claims.go-190-// "WATCH" and discarding the rest, then matches the path regex within the
internal/auth/claims.go-191-// chosen bucket.
internal/auth/claims.go-192-//
internal/auth/claims.go-193-// So a blanket ".*::.*" — which authorizes every REST surface — authorizes
internal/auth/claims.go-194-// nothing here: ".*" is not the string "JOIN". This is measured, not inferred;
--
internal/auth/claims.go:197:func (t *Token) AuthorizesChannel(verb, authPath string) bool {
internal/auth/claims.go-198-	for _, g := range t.grants[ClaimChannels] {
internal/auth/claims.go-199-		if g.verbLiteral == verb && g.path.MatchString(authPath) {
internal/auth/claims.go-200-			return true
internal/auth/claims.go-201-		}
internal/auth/claims.go-202-	}
internal/auth/claims.go-203-	return false
internal/auth/claims.go-204-}

$ grep -n 'Response envelope\|All responses use\|Status codes match\|Live stream\|Honours\|Phoenix WebSocket V2' docs/site/appengine-api.md
94:## Live stream
102:Feeds real-time device events to consumers. Honours `a_ch` claims as room filters.
110:Phoenix WebSocket V2 wire format for upstream Astarte Dashboard compatibility. Answers `phx_join`, `watch`, `phx_leave`, and the `phoenix` heartbeat. Events pushed as `new_event`.
112:## Response envelope
114:All responses use the Astarte envelope format:
126:Status codes match upstream: 401 (no/bad token), 403 (claim mismatch), 404 (not found), 409/422 (conflict/validation).

$ tail -c 400 .mule/for-giulio.md; echo; echo "=== head ==="; sed -n '1,15p' .mule/for-giulio.md; echo "=== lines ==="; wc -l .mule/for-giulio.md .mule/todo.md
the realm-management surface joins *all* keys (`internal/realm/service.go:627`) — but every
writer stores a single-element array (create at `internal/store/realms.go:54-61`, PATCH at
`internal/store/realms.go:165`, `putAuthConfig` at `internal/realm/service.go:636`), so a
two-key realm is unreachable and there is nothing to document. Worth noting because it reads
like a bug on first inspection.

=== head ===
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **COMPATIBILITY.md wording proposal for upstream v1.3.5 (newest stable, 2026-10-05; v1.4.0 is still rc.6-only) — a one-word delta on the open v1.3.2/v1.3.4 proposal below, nothing re-derived.** v1.3.5 is maintenance-only: `gh api repos/astarte-platform/astarte/releases` shows it published `2026-10-05T14:11:44Z`, not a prerelease, with a body of a single `### Fixed` line — "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available" — and `compare/v1.3.4...v1.3.5` is 4 commits / 35 files whose only non-`mix.exs`/`mix.lock` files are `.github/workflows/astarte-build-workflow.yaml`, `CHANGELOG.md`, `README.md`, `astarte-dashboard/package.json` and `docker-compose.yml`. Not one source file changed, so no route, MQTT topic, control message or interface-schema field moved (Astrate implements no `data_updater_plant` and no Erlang RPC server). **The only change a v1.3.5-aware doc carries: the v1.3.4 entry's §Infrastructure-differences sentence — "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target" — reads "…(**newest stable v1.3.5**) as the target".** Everything else stands unchanged: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"); that entry's "v1.4.0 is still rc.5-only" header is already superseded to rc.6 by the open rc.6 delta entry below, not here. This is the same one-word delta recorded on 2026-10-05 in this file ("newest stable v1.3.5"), folded in as asked rather than re-derived. Applied to nothing: `docs/COMPATIBILITY.md` untouched and still 1.2.2-targeted (it contains no `1.3` reference at all), `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). Raw: [v1.3.5](https://github.com/astarte-platform/astarte/releases/tag/v1.3.5).

---
=== lines ===
   639 .mule/for-giulio.md
   283 .mule/todo.md
   922 total

$ grep -n 'appengine-api.md\|Response envelope\|envelope format' .mule/for-giulio.md || echo "no forgiulio dup"; echo "==="; grep -n '426\|securitySchemes\|NativeEvent\|event-schema\|sse-exact\|handshake' .mule/todo.md || echo "no todo dup"
471:- ~~`docs/site/appengine-api.md` documents `GET` and `DELETE` on `/appengine/v1/<realm>/groups/<name>`~~
===
no todo dup

$ sed -n '284,346p' docs/api/astrate_native_api.yaml | cat -n | sed 's/^/   /' | awk '{printf "%d: %s\n", $1+283, substr($0, index($0,$2))}' | sed -n '1,70p'
284: /astrate/v1/{realm}/socket:
285: get:
286: operationId: nativeWebSocket
287: summary: Native live event stream (WebSocket / SSE)
288: description: |
289: Astrate-native real-time event stream. Upgrades to WebSocket by
290: default, or streams Server-Sent Events when the client sends
291: `Accept: text/event-stream` or `?transport=sse`.
292:         9	
293: Subscribes to the realm's committed data and device lifecycle events.
294: Optional `device_id` and `interface` query parameters narrow the
295: event filter.
296:        13	
297: **WebSocket protocol:** server-push only — client frames are drained
298: and ignored. Events are JSON objects:
299:        16	
300: ```json
301: {
302: "event": "incoming_data",
303: "realm": "test",
304: "device_id": "f0VMRgIBAQAAAAAAAAAAAA",
305: "interface": "org.example.Sensors",
306: "path": "/sensors/temperature/value",
307: "value": 23.5,
308: "timestamp": "2024-05-30T13:49:57.045Z"
309: }
310: ```
311:        28	
312: **SSE protocol:** each event is a `data:` frame followed by a blank line.
313:        30	
314: Authentication: `a_ch` realm JWT via Authorization header.
315: tags: [WebSocket]
316: parameters:
317: - $ref: "#/components/parameters/RealmName"
318: - name: device_id
319: in: query
320: schema:
321: type: string
322: description: Filter events to this device ID.
323: - name: interface
324: in: query
325: schema:
326: type: string
327: description: Filter events to this interface.
328: - name: transport
329: in: query
330: schema:
331: type: string
332: enum: [websocket, sse]
333: default: websocket
334: description: Force SSE transport instead of WebSocket.
335: responses:
336: "101":
337: description: WebSocket upgrade successful.
338: "200":
339: description: "SSE stream started (when `Accept: text/event-stream`)."
340: "401":
341: $ref: "#/components/responses/Unauthorized"
342: "403":
343: $ref: "#/components/responses/Forbidden"
344: "500":
345: $ref: "#/components/responses/InternalServerError"
346:        63

$ tail -c 300 .mule/todo.md | cat -A | tail -6
still rc.6-only) in .mule/for-giulio.md (do not edit the file) M-bM-^@M-^T the open v1.3.2/v1.3.4 wording proposal there already stands and its version-reference delta was recorded on 2026-10-05 as reading "newest stable v1.3.5", so fold this in as that one-word delta rather than re-deriving the proposal.$

$ tail -c 80 .mule/for-giulio.md | cat -A
hing to document. Worth noting because it reads$
like a bug on first inspection.$

$ sed -n '14,16p' .mule/for-giulio.md | cat -A
$
---$
$

$ sed -n '17p' .mule/for-giulio.md | cut -c1-60
- **github-issues triage run, 2026-10-06: nothing proposable

Now appending the deliverables — the two for-giulio escalations, then the four task lines.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -13,8 +13,12 @@
 - **COMPATIBILITY.md wording proposal for upstream v1.3.5 (newest stable, 2026-10-05; v1.4.0 is still rc.6-only) — a one-word delta on the open v1.3.2/v1.3.4 proposal below, nothing re-derived.** v1.3.5 is maintenance-only: `gh api repos/astarte-platform/astarte/releases` shows it published `2026-10-05T14:11:44Z`, not a prerelease, with a body of a single `### Fixed` line — "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available" — and `compare/v1.3.4...v1.3.5` is 4 commits / 35 files whose only non-`mix.exs`/`mix.lock` files are `.github/workflows/astarte-build-workflow.yaml`, `CHANGELOG.md`, `README.md`, `astarte-dashboard/package.json` and `docker-compose.yml`. Not one source file changed, so no route, MQTT topic, control message or interface-schema field moved (Astrate implements no `data_updater_plant` and no Erlang RPC server). **The only change a v1.3.5-aware doc carries: the v1.3.4 entry's §Infrastructure-differences sentence — "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target" — reads "…(**newest stable v1.3.5**) as the target".** Everything else stands unchanged: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"); that entry's "v1.4.0 is still rc.5-only" header is already superseded to rc.6 by the open rc.6 delta entry below, not here. This is the same one-word delta recorded on 2026-10-05 in this file ("newest stable v1.3.5"), folded in as asked rather than re-derived. Applied to nothing: `docs/COMPATIBILITY.md` untouched and still 1.2.2-targeted (it contains no `1.3` reference at all), `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). Raw: [v1.3.5](https://github.com/astarte-platform/astarte/releases/tag/v1.3.5).
 
 ---
 
+- **`docs/site/appengine-api.md:114` — "All responses use the Astarte envelope format" is contradicted by the endpoints the page itself lists above it.** The "Response envelope" section (lines 112-126) opens with that sentence and then the `{ "data": ... }` / `{ "errors": { "detail": ... } }` shapes, but the live-stream endpoints documented at lines 94-110 never return one: the native socket pushes bare event JSON as WebSocket text frames or `data: {json}` SSE frames (internal/appengine/stream/ws.go:105-110, 156-162), `/astrate/v1/health` answers `{"status":"ok"}` (internal/observability/health.go:51), `/astrate/v1/readiness` answers `{"status":...,"checks":{...}}` (health.go:77), `/astrate/v1/metrics` answers Prometheus text under `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (metrics.go:42), and — measured today with the repo's own `coder/websocket` (probe in /tmp) — a WebSocket handshake failure on either socket is plain text under `text/plain; charset=utf-8` (`WebSocket protocol violation: Connection header "" does not contain Upgrade`), not JSON. REST appengine paths do use the envelope, so the sentence is right for the endpoints above "Device data" and wrong unqualified. Proposed wording (your voice): scope it — "All REST responses use the Astarte envelope format" — with one line noting the socket and observability surfaces answer un-enveloped. Page untouched.
+
+- **`docs/site/appengine-api.md:102` — "Honours `a_ch` claims as room filters" imports upstream's room semantics onto a socket that has no rooms.** The native socket subscribes to the whole realm bus, narrowed only by the `device_id`/`interface` query parameters (internal/appengine/stream/ws.go:69-75); there is no room concept and no per-room filtering. What a_ch actually does here is gate the whole route: `mw.RequireRealm(auth.ClaimChannels)` requires the token to authorize `GET socket` via the REST verb-regex rule (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112 -> internal/auth/claims.go:164-171). The upstream-measured JOIN/WATCH partition rule DESIGN.md §4.2 describes (`AuthorizesChannel`, claims.go:197-204) is used only by the Phoenix socket's per-room join/watch checks (internal/appengine/channels/ws.go:278, 354). Consequence worth stating: a token whose a_ch list is a blanket `".*::.*"` — which upstream's Channels rule treats as authorizing nothing — authorizes the native socket, while the Phoenix socket, where room filters actually live, still reads JOIN/WATCH literally. Proposed wording (your voice): "guarded by the `a_ch` claim (GET on `socket`); `device_id`/`interface` narrow the stream". Page untouched.
+
 - **github-issues triage run, 2026-10-06: nothing proposable for the 24th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented or closed on GitHub.** The recipe's `--limit 40` command prints all **30** open issues: **26** `mule-alarm` **#94–#119** plus the same four non-alarm, none of them proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (`milestone-4.0`, already escalated), and one the standing instruction excludes from triage entirely. New since the 2026-10-05 run: **#119**, created 2026-10-06T10:58:45Z, "nothing has landed in 15h" — a real window, not an artefact: the land before it was `230dd11` at 2026-10-05T21:04:09+02:00 (19:04Z), 15.9h earlier, and it was ended by today's first land `cbdd3e8` at 16:05:39+02:00 (14:05Z), three hours after the alarm. So it belongs to the genuine 14–16h band the 2026-10-05 entry measured, i.e. overnight silence between evening and afternoon batches. **Proposal, extending the 2026-10-05 one: close #94–#118, all twenty-five** — each is a one-day alarm superseded by the next, and nothing reads issue state; leave **#119** (today's, real window, already ended) to self-expire. The `check_pulse`-before-`beat` ordering fix (move the call at `tools/mule.sh:633` to after the tick's `beat`) is still escalated above, still yours, and still unactioned — today's alarm is not evidence for it either way.
   **The review queue moved 9 commits in 24h and main's four did not move again: `git rev-list --left-right --count origin/main...HEAD` is 4 / 679 today against 4 / 670 yesterday.** The four are main-side commits this branch lacks (`f1d0069`, `ca47b35`, `99743ca`, `2a033d4` — all ancestors of `origin/main`), and **#93's work is still entirely off main**: `8c61268` ("mule: log issue-93") and `24ad5b8` (the aclhook comment rewrite, whose commit message is the issue title verbatim) both satisfy `git merge-base --is-ancestor <sha> origin/main` → false, and `git branch -a --contains 24ad5b8` lists only `mule/queue` / `origin/mule/queue`. That is **32 days** of `mule-review` with `bash tools/mule.sh review` unread. I have not touched git and will not.
   **One detail to fold into that review rather than queue separately:** the rewritten comment at `internal/broker/aclhook.go` (the `control/keyAgreement` deny) says "…upstream's wire spec (082-key_agreement_protocol.md) only ships with v1.4.0-rc.5, v1.3.3 being the newest stable tag". The 082 citation, the "deliberately denied as Astrate's own choice" framing and the unchanged deny + `aclhook_test.go` case are exactly what #93 asked for, but the version reference has aged — upstream's newest stable is **v1.3.5** (2026-10-05) and newest overall `v1.4.0-rc.6` (2026-09-30). No separate task: one word in the same diff when someone finally reads it.
 

✗ Edit .mule/todo.md failed
Error: Could not find oldString in the file. It must match exactly, including whitespace, indentation, and line endings.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -280,4 +280,8 @@
 - [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
 - [x] compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (newest stable, 2026-10-05, maintenance-only: mix.exs/mix.lock version bumps plus a single CHANGELOG "astarte_data_updater_plant RPC server reliability" entry, no wire/API surface change; v1.4.0 is still rc.6-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2/v1.3.4 wording proposal there already stands and its version-reference delta was recorded on 2026-10-05 as reading "newest stable v1.3.5", so fold this in as that one-word delta rather than re-deriving the proposal.
+- [ ] docs-native-socket-handshake-errors [auto]: add the two WebSocket-handshake failures to both socket operations in docs/api/astrate_native_api.yaml — `/astrate/v1/{realm}/socket` responses (yaml:335-345) and `/appengine/v1/socket/websocket` (yaml:398-404) document only 101/401/403/500 and 101/401/500 — but measured today with the repo's own dependency (probe in /tmp, coder/websocket accept.go `verifyClientRequest`): a plain GET with no `Connection: Upgrade` answers **426 Upgrade Required**, `Content-Type: text/plain; charset=utf-8`, body `WebSocket protocol violation: Connection header "" does not contain Upgrade\n`, and `Sec-WebSocket-Version: 14` answers **400** `unsupported WebSocket protocol version (only 13 is supported): "14"\n`. Add both to each operation and say two things explicitly: the bodies are plain text from `http.Error`, NOT the `{"errors":{"detail":...}}` envelope every other response on this spec uses (a client that JSON-parses them breaks), and they are only reached after authentication — native: the route's guard `mw.RequireRealm(auth.ClaimChannels)` wraps `handle` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112), Phoenix: the realm/token checks at internal/appengine/channels/ws.go:51-68 run before `websocket.Accept` at :71 — and only on the WebSocket branch (`wantsSSE` is tested first at ws.go:78), so 401/403 and the 200 SSE path win. Do NOT add 405: ServeMux's method-mismatch applies equally to every GET route in every spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-native-socket-sse-exact-accept [auto]: fix the SSE negotiation prose on `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml — the description (yaml:289-291), the `transport` parameter (yaml:328-334) and the `200` response (yaml:338-339) read as if any client asking for `text/event-stream` gets SSE, but the predicate `wantsSSE` is exact string equality (internal/appengine/stream/ws.go:151-153): measured, `Accept: text/event-stream, text/plain` and `Accept: text/event-stream;q=1.0` both return false and fall through to `websocket.Accept`, where a non-upgrade client gets the 426 above; `transport` selects SSE only for the literal `sse`, and `transport=websocket` / `transport=` (empty) both mean WebSocket. Say "exactly" for the header, say any other `Accept` falls through to the upgrade attempt, and do NOT claim the server rejects out-of-enum `transport` values — the enum is a client-side contract the handler never enforces. While there, give the `200` its real media type: the handler sets `Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive` and writes one `data: {...}` frame per event followed by a blank line (ws.go:123-126, 156-162). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-native-socket-security-scheme [auto]: `/astrate/v1/{realm}/socket` is documented as if public in docs/api/astrate_native_api.yaml — the spec has no `components.securitySchemes` at all and root `security: []` (yaml:21), and the operation declares no `security`, so a client generated from this spec sends no `Authorization` header and is answered 401 `{"errors":{"detail":"Unauthorized"}}` by `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112, bearer token read only from the Authorization header at middleware.go:131, 178-196). Add a scheme mirroring `a_aea` (docs/api/astarte_appengine_api.yaml:1436-1443: apiKey in header `Authorization`, described as a realm JWT carrying the `a_ch` claim) and set `security` on this operation only — root `security: []` stays so health/readiness/metrics/version remain documented as unauthenticated. If you model the Phoenix twin too, its credential is `?token=` in the query (internal/appengine/channels/ws.go:53-54), i.e. apiKey `in: query`, not header. Description notes worth carrying: the scheme is case-insensitive `Bearer` with an optional colon, and an unknown realm is 401 not 404 (middleware.go:68-72). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-native-socket-event-schema [auto]: document the native socket's event payload in docs/api/astrate_native_api.yaml — currently only a fenced `incoming_data` example inside the operation description (yaml:300-310), with no `components.schemas` entry anywhere (yaml:425-478), and that example is not the general shape: `wireEvent` marks `interface`, `path` and `value` `omitempty` (internal/appengine/stream/ws.go:51-58), so measured on a verbatim copy a lifecycle event marshals to `{"event":"device_connected","realm":"test","device_id":"dev1","timestamp":"2024-05-30T13:49:57.045Z"}` — no interface/path/value at all (and `value: false`/`value: 0` ARE emitted, because omitempty on an `any` drops only nil). Add a `NativeEvent` schema (event, realm, device_id, optional interface/path/value, timestamp) referenced from the operation, state that `event` is one of the internal/engine/stream/bus.go Kind* constants (bus.go:20-26), and record that this socket drops `InterfaceMajor`, `IP`, `ErrorName` and `ErrorMetadata`, which `stream.Event` carries (bus.go:39-65) — so a `device_error` arrives here with no reason, unlike the Phoenix twin, which marshals the whole struct (internal/appengine/channels/ws.go:406-416). That last sentence documents a code gap, not a spec one: do not add the fields; if they are wanted on the wire, that is a separate code task. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ grep -c 'docs-native-socket' .mule/todo.md; tail -n 8 .mule/todo.md | cut -c1-90
4
- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in do
- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in d
- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward
- [x] compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (new
- [ ] docs-native-socket-handshake-errors [auto]: add the two WebSocket-handshake failures
- [ ] docs-native-socket-sse-exact-accept [auto]: fix the SSE negotiation prose on `/astra
- [ ] docs-native-socket-security-scheme [auto]: `/astrate/v1/{realm}/socket` is documente
- [ ] docs-native-socket-event-schema [auto]: document the native socket's event payload i

$ sed -n '13,20p' .mule/for-giulio.md | cut -c1-80; grep -c 'appengine-api.md' .mule/for-giulio.md
- **COMPATIBILITY.md wording proposal for upstream v1.3.5 (newest stable, 2026-1

---

- **`docs/site/appengine-api.md:114` — "All responses use the Astarte envelope

- **`docs/site/appengine-api.md:102` — "Honours `a_ch` claims as room filters"

3

Recipe carried out on the **astrate-native** surface. No code changes, no git, nothing on GitHub.

**What I found (measured, probe in /tmp/opencode/probe):**
- Both sockets' specs omit the **426 / 400 WebSocket-handshake failures** (`verifyClientRequest`): a plain GET answers 426 `text/plain` "WebSocket protocol violation...", `Sec-WebSocket-Version: 14` answers 400. Bodies are plain text, not the `{"errors":...}` envelope every other response on the spec uses.
- SSE negotiation is **exact-string** (`wantsSSE`): `Accept: text/event-stream, text/plain` and `;q=1.0` fall through to the upgrade attempt; `transport` means SSE only for the literal `sse`.
- The native socket op has **no `security` requirement** and the spec has no `securitySchemes` — a generated client sends no `Authorization` header and gets 401.
- The event payload has **no schema**; lifecycle events omit `interface`/`path`/`value`, and `wireEvent` drops `InterfaceMajor`/`IP`/`ErrorName`/`ErrorMetadata` the Phoenix twin pushes.
- `docs/site/appengine-api.md` says "All responses use the Astarte envelope format" and "Honours a_ch claims as room filters" — both contradicted by code → escalated.

```
Done: docs-sync (astrate-native surface) — proposed 4 docs/api fix tasks + 2 for-giulio.md escalations, nothing implemented
Files: .mule/todo.md (4 new lines), .mule/for-giulio.md (2 new entries)
Verified: go run /tmp/opencode/probe/*
   -> pass (426/400 handshake, SSE exact-match table, wireEvent omitempty shapes measured)
Unsure: the a_ch "room filters" and dropped-event-fields readings are inference from code plus DESIGN.md; the CREATE-WS 403 path on the native socket is not reachable for real (origin check is disabled, InsecureSkipVerify: true) but 403 is still correctly documented from the auth middleware
Follow-ups: (todo.md) docs-native-socket-handshake-errors, docs-native-socket-sse-exact-accept, docs-native-socket-security-scheme, docs-native-socket-event-schema; (for-giulio.md) appengine-api.md envelope-format and a_ch-room-filters contradictions
```
