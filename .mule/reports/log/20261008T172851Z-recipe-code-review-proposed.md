slug: recipe-code-review
verdict: proposed
at:  6df81bc
ran: 2026-10-08T17:28:51Z on DietPi in 822s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md

$ ls .mule/reviews/ 2>/dev/null; echo "---LOG tail---"; tail -60 .mule/log.md 2>/dev/null; echo "---TODO---"; cat .mule/todo.md 2>/dev/null
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_11c8384a4001cmAfCscrtErTSD

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

$ ls internal/ pkg/ 2>/dev/null; echo "=== sizes ==="; rg -c '' --glob 'internal/**/*.go' | sort -t: -k2 -rn | head -40
internal/:
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

pkg/:
astarteapi
deviceid
interfaceschema
payload
=== sizes ===
internal/swagger/swagger_test.go:1499
internal/engine/triggers/actions_test.go:1199
internal/appengine/channels/ws_test.go:1005
internal/flowapi/service.go:992
internal/appengine/service.go:846
internal/store/datastreams_test.go:823
internal/engine/fakestore_test.go:818
internal/realm/http_test.go:782
internal/engine/triggers/match.go:769
internal/flow/blocks/container/block_test.go:766
internal/appengine/http.go:748
internal/realm/service.go:686
internal/pairing/service_test.go:668
internal/pairing/http_test.go:649
internal/broker/broker_test.go:640
internal/engine/triggers/actions.go:629
internal/broker/broker_e2e_test.go:618
internal/engine/triggers/match_test.go:594
internal/config/config_test.go:584
internal/flow/blocks/httpblocks_test.go:577
internal/realm/http.go:569
internal/flow/router_test.go:565
internal/engine/data.go:552
internal/broker/sessionstore.go:533
internal/appengine/data.go:514
internal/appengine/http_test.go:512
internal/housekeeping/http_test.go:511
internal/engine/cache.go:502
internal/broker/authhook.go:496
internal/flow/blocks/virtualdevicepool/virtualdevicepool_test.go:493
internal/engine/router.go:490
internal/store/devices.go:485
internal/appengine/formats_parity_test.go:467
internal/engine/data_test.go:461
internal/engine/batch.go:459
internal/engine/engine.go:456
internal/store/devices_test.go:446
internal/engine/batch_test.go:446
internal/appengine/channels/ws.go:442
internal/pairing/service.go:440

$ ls -la .mule/reviews/ 2>/dev/null; echo "=== log tail ==="; tail -40 .mule/log.md
total 224
drwxr-xr-x 2 root root  4096 Oct  7 19:43 .
drwxr-xr-x 7 root root  4096 Oct  8 13:03 ..
-rw-r--r-- 1 root root  4600 Sep  7 20:15 appengine-2026-09-07.md
-rw-r--r-- 1 root root  3443 Sep 15 22:09 astarteapi-2026-09-15.md
-rw-r--r-- 1 root root  5192 Sep 14 13:34 auth-2026-09-14.md
-rw-r--r-- 1 root root  8188 Sep  4 20:43 broker-2026-09-04.md
-rw-r--r-- 1 root root  6580 Sep 21 21:26 broker-2026-09-21.md
-rw-r--r-- 1 root root  6967 Sep 18 22:00 channels-2026-09-18.md
-rw-r--r-- 1 root root 10261 Sep 25 19:46 cmd-astrate-2026-09-25.md
-rw-r--r-- 1 root root  4340 Sep 13 19:26 config-2026-09-13.md
-rw-r--r-- 1 root root 10392 Oct  2 13:43 container-2026-10-02.md
-rw-r--r-- 1 root root  5264 Sep 18 19:25 deviceid-2026-09-18.md
-rw-r--r-- 1 root root  3141 Sep  6 20:57 engine-2026-09-06.md
-rw-r--r-- 1 root root  4692 Sep 20 19:30 engine-triggers-2026-09-20.md
-rw-r--r-- 1 root root  5516 Sep  5 21:19 flow-2026-09-05.md
-rw-r--r-- 1 root root  6931 Sep 23 20:39 flow-2026-09-23.md
-rw-r--r-- 1 root root  4837 Sep 10 19:00 flowapi-2026-09-10.md
-rw-r--r-- 1 root root  8949 Sep 26 19:42 forward-2026-09-26.md
-rw-r--r-- 1 root root  2595 Sep 15 13:29 housekeeping-2026-09-15.md
-rw-r--r-- 1 root root  4914 Sep 17 21:11 httpx-2026-09-17.md
-rw-r--r-- 1 root root  4712 Sep 11 21:40 interfaceschema-2026-09-11.md
-rw-r--r-- 1 root root  4594 Sep 16 20:48 observability-2026-09-16.md
-rw-r--r-- 1 root root  3506 Sep  9 13:07 pairing-2026-09-09.md
-rw-r--r-- 1 root root  4915 Oct  7 19:43 pairing-2026-10-07.md
-rw-r--r-- 1 root root  6002 Sep 12 20:24 payload-2026-09-12.md
-rw-r--r-- 1 root root  4259 Sep  9 21:13 realm-2026-09-09.md
-rw-r--r-- 1 root root  4605 Sep  8 18:59 store-2026-09-08.md
-rw-r--r-- 1 root root  5504 Sep 21 13:32 store-2026-09-21.md
-rw-r--r-- 1 root root  4511 Oct  4 21:08 store-2026-10-04.md
-rw-r--r-- 1 root root  2644 Sep 24 20:37 swagger-2026-09-24.md
=== log tail ===
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

$ ls -la .mule/
total 14936
drwxr-xr-x  7 root root     4096 Oct  8 13:03 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root        0 Oct  8 13:03 .alarmed
-rw-r--r--  1 root root       13 Oct  8 19:15 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  7 20:47 .heartbeat
-rw-r--r--  1 root root    49152 Oct  8 19:15 .last-output
-rw-r--r--  1 root root        0 Oct  8 19:15 .lock
-rw-r--r--  1 root root       12 Oct  8 19:15 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 14816771 Oct  8 19:15 cron.log
drwxr-xr-x  2 root root     4096 Oct  7 20:02 failed
-rw-r--r--  1 root root   158413 Oct  8 13:10 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    15361 Oct  7 20:47 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  7 19:43 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   160095 Oct  7 20:47 todo.md

$ grep -c '^- \[ \]' .mule/todo.md; echo "--- open slugs ---"; grep -o '^- \[ \] [a-z0-9-]*' .mule/todo.md | sed 's/- \[ \] //' | sort
14
--- open slugs ---
broker-external-bus-intake
container-orphan-boot-cleanup
deviceid-strict-trailing-bits
deviceid-trailing-bits-upstream-probe
flow-boot-resumes-stopped-durable-flows
flowapi-autorestart-shutdown-cancel
hygiene-govulncheck
race-check-appengine
race-check-engine
race-check-flow
race-check-pkg
race-check-store
store-devices-alias-lowest-id-test
store-devices-inhibit-re-register

$ ls internal/appengine/ internal/engine/ internal/realm/ | cat; echo "=== engine sizes ==="; rg -c '' --glob 'internal/engine/**/*.go' | sort -t: -k2 -rn
internal/appengine/:
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

internal/engine/:
batch.go
batch_test.go
bus_events_test.go
cache.go
cache_test.go
capabilities.go
capabilities_test.go
control.go
control_test.go
data.go
data_test.go
devicedata_test.go
engine.go
engine_e2e_test.go
engine_integration_test.go
fakestore_test.go
forward
introspection.go
introspection_test.go
mqttv1_parity_test.go
router.go
router_test.go
serverdata.go
serverdata_test.go
stream
testdata
topics.go
topics_test.go
triggers
triggers_change_integration_test.go
triggers_change_test.go

internal/realm/:
async_operation_test.go
dashboard_compat_test.go
detailed.go
detailed_required_encrypted_test.go
detailed_test.go
helpers_test.go
http.go
http_test.go
policy.go
policy_test.go
service.go
=== engine sizes ===
internal/engine/triggers/actions_test.go:1199
internal/engine/fakestore_test.go:818
internal/engine/triggers/match.go:769
internal/engine/triggers/actions.go:629
internal/engine/triggers/match_test.go:594
internal/engine/data.go:552
internal/engine/cache.go:502
internal/engine/router.go:490
internal/engine/data_test.go:461
internal/engine/batch.go:459
internal/engine/engine.go:456
internal/engine/batch_test.go:446
internal/engine/engine_e2e_test.go:440
internal/engine/forward/http_test.go:438
internal/engine/engine_integration_test.go:435
internal/engine/serverdata_test.go:434
internal/engine/triggers/policy_test.go:410
internal/engine/cache_test.go:363
internal/engine/router_test.go:362
internal/engine/serverdata.go:349
internal/engine/control_test.go:330
internal/engine/control.go:316
internal/engine/triggers/events.go:312
internal/engine/triggers/policy.go:297
internal/engine/triggers_change_test.go:293
internal/engine/capabilities_test.go:255
internal/engine/stream/bus_test.go:234
internal/engine/devicedata_test.go:223
internal/engine/stream/bus.go:209
internal/engine/forward/nats_test.go:196
internal/engine/forward/http.go:194
internal/engine/introspection_test.go:189
internal/engine/introspection.go:168
internal/engine/triggers/validation.go:163
internal/engine/triggers/change_test.go:161
internal/engine/triggers/events_test.go:155
internal/engine/mqttv1_parity_test.go:151
internal/engine/triggers/transient_test.go:142
internal/engine/triggers_change_integration_test.go:141
internal/engine/bus_events_test.go:140
internal/engine/triggers/errorname_test.go:112
internal/engine/topics_test.go:109
internal/engine/capabilities.go:102
internal/engine/topics.go:97
internal/engine/triggers/errorname.go:79
internal/engine/forward/envelope_test.go:74
internal/engine/forward/nats.go:53
internal/engine/forward/envelope.go:29
internal/engine/triggers/transient.go:22

→ Read .mule/reviews/engine-2026-09-06.md

$ rg -c '' --glob 'internal/engine/*.go' | sort -t: -k2 -rn
internal/engine/fakestore_test.go:818
internal/engine/data.go:552
internal/engine/cache.go:502
internal/engine/router.go:490
internal/engine/data_test.go:461
internal/engine/batch.go:459
internal/engine/engine.go:456
internal/engine/batch_test.go:446
internal/engine/engine_e2e_test.go:440
internal/engine/engine_integration_test.go:435
internal/engine/serverdata_test.go:434
internal/engine/cache_test.go:363
internal/engine/router_test.go:362
internal/engine/serverdata.go:349
internal/engine/control_test.go:330
internal/engine/control.go:316
internal/engine/triggers_change_test.go:293
internal/engine/capabilities_test.go:255
internal/engine/devicedata_test.go:223
internal/engine/introspection_test.go:189
internal/engine/introspection.go:168
internal/engine/mqttv1_parity_test.go:151
internal/engine/triggers_change_integration_test.go:141
internal/engine/bus_events_test.go:140
internal/engine/topics_test.go:109
internal/engine/capabilities.go:102
internal/engine/topics.go:97

→ Read .mule/reviews/appengine-2026-09-07.md
→ Read .mule/reviews/realm-2026-09-09.md

$ head -20 .mule/reviews/flow-2026-09-23.md; echo "======"; head -20 .mule/reviews/channels-2026-09-18.md; echo "======"; head -15 .mule/reviews/store-2026-10-04.md; echo "======"; head -12 .mule/reviews/pairing-2026-10-07.md
# Review — internal/flow, 2026-09-23

Last reviewed: 2026-09-05 (`flow`). That pass covered the core package;
`internal/flow/blocks/` (and its `astartesource/`, `container/`,
`virtualdevicepool/` subpackages) has never been reviewed, so this pass
focuses there plus the message wire-format decode path that feeds them.

## What I read

Core (`internal/flow/`):
- `pipeline.go` — Parse/Validate; the 09-05 dead source/sink branch was
  confirmed fixed (the duplicate hasSource/hasSink recompute is now gone,
  matching todo line 122).
- `flow.go`, `router.go`, `message.go`, `factory.go`, `composite.go`,
  `substitute.go`, `block.go`, `graph.go` — lifecycle, lanes, message
  wire format (Marshal/UnmarshalJSON + setDataFromWire), registry,
  composite expansion, block graph. Tests read: `message_test.go`
  (roundtrip + rejection), `factory_test.go` (instantiate linear /
  unknown type / construct-error / stop-on-failure / parse def /
  astarte source → null sink lifecycle).
======
# Code review: internal/appengine/channels (+ internal/appengine/stream) — 2026-09-18

Area chosen: the `channels/` subpackage (Phoenix V2 wire protocol at
`/appengine/v1/socket/websocket`) had no review file of its own — the
2026-09-07 appengine review read only `service.go`, `http.go`, `data.go`,
`downsample.go`, so this websocket/room concurrency code was never reviewed.
It is the kind of place (goroutines, mutexes, per-connection loops) a review
earns its keep, and it is container-free, so everything here is T1-checkable.

## Files read

- `internal/appengine/channels/ws.go` (442) — session loop, join/leave/watch/
  unwatch handlers, `watchAuthPath`, `pumpEvents`, `writeFrame`
- `internal/appengine/channels/room.go` (287) — Registry/room lifecycle,
  dispatch fan-out, Member mailbox/drop, retire
- `internal/appengine/channels/protocol.go` (128) — Frame codec, Reply/OK/Err
- `internal/appengine/stream/ws.go` (162) — Astrate-native socket +
  SSE fallback (the documented deviation from the Phoenix protocol)
- `internal/engine/stream/bus.go` (209) — subscription cancel/close semantics
- `internal/engine/triggers/match.go` (330-504) — CompileCondition /
======
# Code review: internal/store — 2026-10-04

Area chosen: internal/store. Previous reviews were 2026-09-08 and 2026-09-21.
The package is large (30 source files) and central; I read non-test source files
in descending size order, with particular attention to devices, datastreams,
interfaces, groups/realms/store plumbing, and their corresponding tests.

## What I read

- devices.go (485), devices_test.go (427) — register/unregister/get, inhibit,
  credentials stamp, connect/disconnect/stats, alias values taken, patch merge,
  payload hint, delete
- datastreams.go (434), datastreams_test.go (823) — COPY append, Series,
  LatestIndividual, IndividualSnapshot, ObjectSeries, Downsample, DownsampleLTTB,
  SeriesSpan
======
# Codebase review — internal/pairing (2026-10-07)

## What I read

- `internal/pairing/service.go` (full, 440 lines) — flows A/B/C, `authenticateDevice`,
  `loadCA`, `burnBcrypt`, `generateSecret`.
- `internal/pairing/http.go` (full, 367 lines) — routes, handlers, `writeServiceError`,
  `bearerSecret`, `remoteIP`, `formatDateTime`.
- `internal/pairing/ca/ca.go` (full, 287 lines) — `Generate`, `Load`, `SignCSR`, `Verify`,
  `parseCSR`, `randomSerial`.
- `internal/pairing/ratelimit.go` (full, 79 lines) — keyed token bucket over an LRU.
- Tests: `service_test.go` (all flows), `http_test.go` (golden flow tests, rate-limit,

$ rg -c '' --glob 'internal/flow/blocks/**/*.go' | sort -t: -k2 -rn | head -20; echo "=== stream ==="; rg -c '' --glob 'internal/appengine/stream/*.go' --glob 'internal/broker/*.go' | sort -t: -k2 -rn
internal/flow/blocks/container/block_test.go:766
internal/flow/blocks/httpblocks_test.go:577
internal/flow/blocks/virtualdevicepool/virtualdevicepool_test.go:493
internal/flow/blocks/mqtt_test.go:401
internal/flow/blocks/mqtt.go:388
internal/flow/blocks/httpblocks.go:344
internal/flow/blocks/container/block.go:320
internal/flow/blocks/jsonpathmap_test.go:284
internal/flow/blocks/transform.go:255
internal/flow/blocks/jsonpathmap.go:238
internal/flow/blocks/transform_test.go:236
internal/flow/blocks/randomsource.go:226
internal/flow/blocks/container/docker.go:223
internal/flow/blocks/virtualdevicepool/virtualdevicepool.go:214
internal/flow/blocks/splitmap_test.go:192
internal/flow/blocks/sort_test.go:190
internal/flow/blocks/container/httpbridge.go:178
internal/flow/blocks/randomsource_test.go:173
internal/flow/blocks/sort.go:171
internal/flow/blocks/astartesource/source.go:171
=== stream ===
internal/broker/broker_test.go:640
internal/broker/broker_e2e_test.go:618
internal/broker/sessionstore.go:533
internal/broker/authhook.go:496
internal/broker/broker.go:343
internal/broker/aclhook.go:246
internal/broker/sessionstore_test.go:227
internal/broker/aclhook_test.go:195
internal/appengine/stream/ws.go:162
internal/appengine/stream/ws_test.go:110
internal/broker/identity.go:104
internal/broker/identity_test.go:103
internal/broker/intake.go:100
internal/broker/lifecycle.go:90
internal/broker/publisher.go:62

→ Read .mule/reviews/flowapi-2026-09-10.md

$ rg -c '' --glob 'internal/flowapi/*.go' --glob 'internal/housekeeping/*.go' --glob 'internal/broker/*.go' | sort -t: -k2 -rn; echo "=== git log since 10-01 ==="; git log --oneline -12
internal/flowapi/service.go:992
internal/broker/broker_test.go:640
internal/broker/broker_e2e_test.go:618
internal/broker/sessionstore.go:533
internal/housekeeping/http_test.go:511
internal/broker/authhook.go:496
internal/broker/broker.go:343
internal/flowapi/service_integration_test.go:334
internal/flowapi/http.go:310
internal/housekeeping/service.go:288
internal/housekeeping/service_test.go:277
internal/housekeeping/http.go:259
internal/broker/aclhook.go:246
internal/broker/sessionstore_test.go:227
internal/flowapi/userblocks.go:209
internal/broker/aclhook_test.go:195
internal/flowapi/service_test.go:167
internal/flowapi/composite_wire_test.go:133
internal/broker/identity.go:104
internal/broker/identity_test.go:103
internal/broker/intake.go:100
internal/flowapi/userblocks_test.go:97
internal/flowapi/http_test.go:93
internal/broker/lifecycle.go:90
internal/broker/publisher.go:62
internal/housekeeping/async_operation_test.go:37
=== git log since 10-01 ===
6df81bc mule: recipe astarte-upstream ran (156s)
d0cab66 mule: refill the queue
bae4595 mule: recipe github-issues ran (463s)
0d2e197 mule: recipe milestones ran (824s)
0f9c65c mule: recipe hygiene timed out
d8124d6 mule: refill the queue
10efb94 mule: recipe docs-sync ran (518s)
ebd5fb0 mule: log pairing-unregister-bad-id-test
322242a mule: pairing-unregister-bad-id-test: cover the malformed-device-IDs branch of Unregister in internal/pairing/service_test.go — Unregister wraps deviceid.Parse failure as store.ErrNotFound (internal/pairing/service.go:231-234) but TestUnregister (service_test.go:340) only exercises a well-formed unknown ID; assert both bad-ID shapes and the well-formed unknown-ID 404 in one table.
403d725 mule: log pairing-onregistered-clock
4e71533 mule: pairing-onregistered-clock: route the registration event timestamp through the injectable clock — Register calls s.OnRegistered(realmName, hwID, time.Now()) with literal wall-clock time (internal/pairing/service.go:222) while every other timestamp in the package funnels through the s.now seam (VerifyCredentials, service.go:333), so a time-travel test cannot pin the DeviceRegisteredEvent timestamp fired via e.HandleDeviceRegistered (cmd/astrate/main.go:154). Change to s.now(), keep TestRegisterEmitsEvent.
b7c8fe8 mule: blocked pairing-burn-bcrypt-cost

$ for d in internal/engine internal/appengine internal/broker internal/housekeeping internal/observability internal/httpx internal/auth internal/config; do n=$(git log --oneline --since=2026-09-06 -- $d | wc -l); echo "$d: $n commits since 09-06"; done
internal/engine: 9 commits since 09-06
internal/appengine: 4 commits since 09-06
internal/broker: 2 commits since 09-06
internal/housekeeping: 2 commits since 09-06
internal/observability: 3 commits since 09-06
internal/httpx: 1 commits since 09-06
internal/auth: 3 commits since 09-06
internal/config: 3 commits since 09-06

$ git log --oneline --since=2026-09-06 -- internal/engine | head -30; echo "=== appengine ==="; git log --oneline --since=2026-09-07 -- internal/appengine | head -30
f7d60af mule: appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
f9db0a6 mule: webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.
ed869b6 mule: forward-envelope-bytes-test [auto]: `marshalEnvelope` (internal/engine/forward/envelope.go:19-29) has no direct test — `rg marshalEnvelope` finds exactly two callers (http.go:99, nats.go:45) and no test file — while envelope.go:6-8 makes the strong claim that "every Forwarder implementation in this package must produce byte-identical envelopes for the same inputs". Both suites check the shape only by unmarshalling into `bodyShape` (http_test.go:24-29), which cannot see key order or the exact bytes and would pass unchanged if a field were renamed on both sides. Add internal/engine/forward/envelope_test.go with a table pinning the exact output for (nil, nil) → `{"realm":"r","trigger":"t","action":null,"event":null}`, for the empty-non-nil pair the rule at 16-18 is about, and for a valued pair; container-free, no HTTP server needed.
915c260 mule: forward-status-error-body [auto]: a non-2xx forward returns `forward: status %d` and nothing else (internal/engine/forward/http.go:121-123) after copying the **entire** response body into `io.Discard` with no limit (http.go:117-120) — the one thing an operator needs to debug a 500 from the bus is destroyed, and the sibling webhook request in the same codebase does bound its drain with `io.LimitReader(resp.Body, 1<<20)` (internal/engine/triggers/actions.go:626), so the two near-identical request paths have drifted on exactly that line. Read a bounded prefix (e.g. 512 bytes) of the body into the returned error and bound the discard the way actions.go:626 does. Add a case beside `TestStatusTable` (internal/engine/forward/http_test.go:164) where the handler answers 500 with body `{"error":"boom"}` and assert the error mentions both `500` and `boom`; it fails today because the error carries only the status code.
3a0f153 mule: forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
1cbebd5 mule: forward-static-header-validation [auto]: `New` (internal/engine/forward/http.go:38-71) validates the URL and the method at boot precisely because of the rule its own comment states at 42-45 — "an unusable endpoint must fail here rather than surface once per delivery" — but `cfg.StaticHeaders` is stored unvalidated (http.go:68) and applied per delivery with `req.Header.Set` (http.go:110-112), and net/http checks header names and values only at write time. Measured on the real net/http path with a throwaway program: a name of `X Bad Name` or `X:Foo` fails per delivery with `net/http: invalid header field name`, and a value containing a newline fails with `net/http: invalid header field value for "X-Foo"` (a plain space in a value is legal). So one typo in `triggers.forward.static_headers` makes the process boot clean, accept triggers, and then count `astrate_engine_trigger_deliveries_total{outcome="failed"}` and log "custom trigger action failed" (internal/engine/triggers/actions.go:499) on every custom action forever. Validate in `New`: the RFC 7230 field-name grammar is exactly the token grammar `validMethod` (http.go:77-95) already implements, so reuse that predicate per name (rename it if you like) and add a CR/LF rejection per value. Add a table case to `TestNewRejectsBadConfig` (internal/engine/forward/http_test.go:271) asserting `New` returns an error for a bad name and a bad value and that no request is attempted; this test fails today because `New` returns nil error for both.
3c8b89a mule: triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
739df67 mule: errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
e46ce11 mule: server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.
=== appengine ===
f84bbb5 mule: appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
f7d60af mule: appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
d48e17f mule: channels-rejoin-joinref-tagging-test: pin in internal/appengine/channels/ws_test.go that a second `phx_join` on the same topic updates the session's `joinRef` (ws.go:284-291) and subsequent `new_event` pushes carry the fresh ref (read per event at ws.go:403) — no test drives a rejoin today (all `TestWireSession_*` join once), so the comment's "server pushes must be tagged with the current one" rots unasserted. Wire: join with ref "1", rejoin with "2", push an event, assert element 0 is `"2"`.
aeaeb47 mule: channels-group-watch-membership-scope: group watches authorize against `WATCH::groups/<name>/…` (ws.go:178-205, conformance channels.json) but `Room.dispatch` ever only filters by `watchEntry.deviceID` (empty for a group watch — room.go:36-40, 196-200) and the trigger compiles with no group/device scope (triggers/match.go:374-396, because upstream carries group_name at the payload top level), so a group-scoped token receives every realm event matching interface/path/on, members or not. Plumb a group-membership resolver into internal/appengine/channels (it knows only `Bus` today), resolve membership at Watch time and store it in `watchEntry`, filter deliveries by it. Add a container-free test with a fake resolver: watch group `probe`, publish an event from a member and a non-member, assert only the member's arrives. Verify the member-filter rule against upstream on the Legion if disputed; the boundary leak stands regardless.

$ git log --oneline --since=2026-09-09 -- internal/realm internal/housekeeping internal/observability internal/httpx internal/auth | cat; echo "=== realm/housekeeping/obs/httpx/auth todo slugs ==="; grep -o '^- \[[x!]\] [a-z0-9-]*' .mule/todo.md | grep -E 'realm-|hk-|housekeep|observ|metrics|httpx|auth-|jwt|token' | sort | head -40
f965f39 mule: docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
986c100 mule: httpx-cors-vary-origin-passthrough [auto]: in internal/httpx/cors.go non-wildcard mode, the disallowed-origin pass-through (cors.go:26-29) returns without `Vary: Origin`, so responses whose CORS outcome depends on the Origin header reach shared caches without the RFC 9110 §12.5.5 variation marker (a cache could serve the CORS-approving variant to a disallowed origin, or the headerless variant to an allowed one); the stamp lives only in the allowed branch (cors.go:35-36) and only the preflight pins it (cors_test.go:42). Add `h.Add("Vary", "Origin")` on the pass-through branch too, and assert `Vary: Origin` in TestCORSActualRequestStamped (actual non-preflight, pinned only for Allow-Origin today) and TestCORSDisallowedOriginPassesThroughUnstamped. Additive header change; container-free.
f6c27f5 mule: observability-health-content-type-test [auto]: assert `Content-Type: application/json` on `/astrate/v1/health` and `/astrate/v1/readiness` (set by `writeJSON`, health.go:81) and on both the 200 and 503 paths of the per-service `/health` (set by hand in compat.go:26,31) — only the version endpoint pins the header today (compat_version_test.go:31-33), so a regression dropping the header on the /astrate/v1 surface passes the suite.
19cd4f1 mule: observability-dbpool-gauge-coverage [auto]: extend `TestMetricsExposesGauges` (internal/observability/observability_test.go:25-34) to assert all four `astrate_db_pool_*` gauges — it currently pins only `acquired_conns` and `max_conns`, so a regression deleting or mis-wiring `idle_conns` or `total_conns` (metrics.go:77-78) passes the suite; expect idle 5, total 7 from the supplied `DBPoolStats{2,5,7,10}`.
4ef2e23 mule: observability-readiness-wedge-test [auto]: make the readiness budget injectable and pinned — add a `timeout time.Duration` field on `Health` (internal/observability/health.go) defaulting to `readinessTimeout` (health.go:12, no test currently exercises the wedged-dependency path), and add a test in observability_test.go asserting a ctx-honoring wedged check makes `/astrate/v1/readiness` return 503 reporting it failing within a short injected budget, instead of hanging the endpoint.
015a706 mule: hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto]
874797b mule: auth-cache-default-size-test [auto]: pin in internal/auth/cache_test.go that `NewCache(0)` and `NewCache(-1)` fall back to `DefaultCacheSize` (cache.go:29-31) — the `size < 1` branch exists purely so `lru.New` cannot panic and has zero coverage; a regression dropping the fallback (panic on `lru.New` with a non-tick size, or a silently tiny cache) would pass every test. Construct with both sizes, assert no panic and a working `Verify` with `lru.Len() == 0`.
b91247c mule: auth-empty-claims-403-test [auto]: pin in internal/auth/middleware_test.go that a validly signed token carrying NO Astarte claims is `403` on a guarded route, not `401` — the authenticated-but-grant-less distinction (middleware.go:149-151, `Authorizes` on a nil grants map → false → `WriteForbidden`) is currently untested: every 403 fixture has a claim on the wrong surface, so a regression mapping grant-less tokens to 401 (or skipping the forbidden path) passes the suite. Sign RS256 with empty `jwt.MapClaims`, hit the pairing route, expect 403 + golden `envelope_403.json`.
6524c7f mule: auth-iat-not-required-test [auto]: pin in internal/auth/jwt_test.go that a token with a *future* `iat` still verifies — the "`iat` is not required — upstream parity" rule (jwt.go:104-105, no `jwt.WithIssuedAt` in the parser, jwt.go:115-118) has zero coverage, so a regression that starts validating `iat` or `WithIssuedAt` would silently reject tokens from a skewed-clock issuer and pass every existing suite. Sign RS256 with `iat` = now+1h and no `exp`/`nbf`, assert `Verify` succeeds.
39eb70a mule: realm-pure-helper-tests [auto]: add container-free unit tests for the four untested pure helpers in internal/realm. Create `internal/realm/helpers_test.go` (no build tag) covering: (1) `normaliseIfaceName` — lowercase, hyphen strip, mixed case+hyphens, empty string, already-normalized passthrough; (2) `joinPEM` — empty slice, single key, two keys separated by `\n`; (3) `triggerErrorBody` — empty action+simple_triggers, action-only, simple_triggers-only with index-aligned `{}` substitution, both present; (4) `majorParam` — valid positive, zero, negative (→404), non-numeric (→404). All tests use httptest.ResponseRecorder for majorParam; the rest are direct function calls. No Docker.
=== realm/housekeeping/obs/httpx/auth todo slugs ===
- [!] appengine-group-token-roundtrip-test
- [!] docs-sync-hk-error-example-capitalisation
- [!] docs-sync-hk-patch-422-field-error-examples
- [!] docs-sync-hk-wrong-type-field-400
- [!] hk-zero-reglimit-create-asymmetry
- [x] auth-cache-default-size-test
- [x] auth-empty-claims-403-test
- [x] auth-iat-not-required-test
- [x] docs-native-ae-realm-version-endpoint
- [x] docs-native-hk-version-description
- [x] docs-native-metrics-content-negotiation
- [x] docs-native-realm-name-pattern
- [x] docs-sync-hk-async-operation-param
- [x] docs-sync-hk-delete-gating-responses
- [x] docs-sync-hk-error-detail-examples-split
- [x] docs-sync-hk-errordetail-schema-example
- [x] docs-sync-hk-patch-endpoint
- [x] docs-sync-hk-realm-name-pattern
- [x] docs-sync-hk-realm-name-response-schemas
- [x] docs-sync-hk-retention-field
- [x] docs-sync-hk-retention-zero-is-unset
- [x] docs-sync-hk-validation-example
- [x] docs-sync-native-metrics-example-fake-series
- [x] docs-sync-native-socket-query-token-auth
- [x] docs-sync-rm-auth-403
- [x] docs-sync-rm-put-auth-422
- [x] hk-zero-retention-asymmetry
- [x] housekeeping-tests
- [x] httpx-cors-vary-origin-passthrough
- [x] observability-dbpool-gauge-coverage
- [x] observability-health-content-type-test
- [x] observability-readiness-wedge-test
- [x] realm-pure-helper-tests

$ ls .mule/recipes/; echo "=== rotation mention ==="; grep -rn "rotate\|rotation\|not been reviewed\|last reviewed" .mule/reviews/*.md | head -20
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md
=== rotation mention ===
.mule/reviews/auth-2026-09-14.md:16:key-rotation loop, the immutability that makes cache sharing safe, the
.mule/reviews/broker-2026-09-04.md:16:This is the first review of `broker`; the rotation had covered store, realm, flow, engine, pairing.
.mule/reviews/engine-triggers-2026-09-20.md:4:paths; this run targets the files the rotation had not read since: the
.mule/reviews/flow-2026-09-05.md:16:This is the first review of `flow`; the rotation checked broker last (2026-09-04).
.mule/reviews/interfaceschema-2026-09-11.md:6:no review file. The last review ran 09-10, so this is the next rotation slot.
.mule/reviews/pairing-2026-10-07.md:47:  `/credentials` + register". The hw_id comes from the request body and rotates freely, so a
.mule/reviews/pairing-2026-10-07.md:63:- **`Register` counts every device row against the registration limit, rotations included**
.mule/reviews/store-2026-09-21.md:9:- `devices.go` (485) — register/rotate, unregister, introspection diff + locking,
.mule/reviews/store-2026-09-21.md:56:only registered → rotated and confirmed → 422; the inhibit round-trip
.mule/reviews/store-2026-10-04.md:30:Untested: devices_test.Lifecycle covers re-register before credentials (rotates secret) but not the inhibit→re-register case. Uninhibit already restores 'registered' if never confirmed (SetDeviceInhibited logic), but the re-registration path overwrites status to 'registered' regardless of prior inhibited state — that’s the bug.

$ rg -c '' --glob 'pkg/**/*.go' --glob 'internal/config/*.go' | sort -t: -k2 -rn; echo "=== payload/intf/config review dates ==="; head -8 .mule/reviews/payload-2026-09-12.md; echo ---; head -8 .mule/reviews/interfaceschema-2026-09-11.md; echo ---; head -8 .mule/reviews/config-2026-09-13.md
pkg/payload/payload_test.go:856
pkg/interfaceschema/parse.go:664
internal/config/config_test.go:584
pkg/interfaceschema/compat_test.go:470
internal/config/config.go:420
pkg/payload/json.go:413
pkg/interfaceschema/types.go:406
pkg/interfaceschema/parse_test.go:388
pkg/payload/value.go:383
pkg/payload/bson.go:352
pkg/interfaceschema/violations_test.go:309
pkg/astarteapi/envelope_test.go:276
pkg/astarteapi/envelope.go:249
pkg/interfaceschema/trie_test.go:245
pkg/interfaceschema/compile_test.go:235
pkg/deviceid/deviceid_test.go:226
pkg/payload/value_test.go:223
pkg/payload/payload.go:217
pkg/interfaceschema/required_encrypted_test.go:194
pkg/interfaceschema/trie.go:161
pkg/deviceid/deviceid.go:132
pkg/interfaceschema/compat.go:131
pkg/payload/fuzz_test.go:122
pkg/interfaceschema/compile.go:122
pkg/payload/sniff_test.go:120
pkg/interfaceschema/violations.go:82
pkg/payload/required_test.go:71
pkg/payload/sniff.go:41
=== payload/intf/config review dates ===
# Code review — pkg/payload (2026-09-12)

Area chosen because recent reviews covered `flowapi` (09-10), `interfaceschema`
(09-11), `appengine` (09-07) and the store/engine/broker/flow/pairing/realm line —
but never `pkg/payload`, the wire-format codec, which sits at the top of the
upstream-compat risk list (BSON `{v,t}` parity with the official SDKs, the strict
JSON profile for constrained clients). I read every non-test source file plus all
six test files, largest first.
---
# Code review: pkg/interfaceschema — 2026-09-11

Area chosen because `appengine` (09-07), `store` (09-08), `pairing`/`realm`
(09-09), `flowapi` (09-10) have review files, `broker`/`flow`/`engine` were
09-04..09-06, and `interfaceschema` is the largest package (~3200 lines) with
no review file. The last review ran 09-10, so this is the next rotation slot.
I read every non-test file largest first — `parse.go` (664), `types.go` (406),
`compat.go` (131), `compile.go` (122), `trie.go` (161), `violations.go` (82) —
---
# Code review: internal/config — 2026-09-13

## What I read

`internal/config/config.go` (411 lines) in full, `internal/config/config_test.go` (407 lines) in full.
Traced callers of every unvalidated field to confirm downstream guard status via targeted grep.

## Findings

$ grep -o '^- \[[x!]\] \(payload\|interfaceschema\|config\|swagger\|cmd-\|astarteapi\|deviceid\)[a-z0-9-]*' .mule/todo.md | sort
- [!] cmd-healthcheck-contract-test
- [x] astarteapi-metadata-envelope-test
- [x] astarteapi-multikey-sorted-golden
- [x] cmd-devcert-fields-test
- [x] cmd-loadsealer-masterkeyfile-test
- [x] config-fail-loud-dev-mode
- [x] config-fail-loud-engine-shards
- [x] config-key-resolvers-test
- [x] interfaceschema-compat-attrs-flip-test
- [x] interfaceschema-enum-roundtrip-coverage
- [x] interfaceschema-malformed-placeholder-fixtures
- [x] interfaceschema-object-attrs-uniformity
- [x] interfaceschema-trie-interior-literal-param
- [x] payload-json-malformed-t-not-tolerated
- [x] payload-longinteger-fraction-quantize
- [x] swagger-httptest-coverage
- [x] swagger-sub-failfast

$ sed -n 1,60p .mule/reviews/flow-2026-09-23.md
# Review — internal/flow, 2026-09-23

Last reviewed: 2026-09-05 (`flow`). That pass covered the core package;
`internal/flow/blocks/` (and its `astartesource/`, `container/`,
`virtualdevicepool/` subpackages) has never been reviewed, so this pass
focuses there plus the message wire-format decode path that feeds them.

## What I read

Core (`internal/flow/`):
- `pipeline.go` — Parse/Validate; the 09-05 dead source/sink branch was
  confirmed fixed (the duplicate hasSource/hasSink recompute is now gone,
  matching todo line 122).
- `flow.go`, `router.go`, `message.go`, `factory.go`, `composite.go`,
  `substitute.go`, `block.go`, `graph.go` — lifecycle, lanes, message
  wire format (Marshal/UnmarshalJSON + setDataFromWire), registry,
  composite expansion, block graph. Tests read: `message_test.go`
  (roundtrip + rejection), `factory_test.go` (instantiate linear /
  unknown type / construct-error / stop-on-failure / parse def /
  astarte source → null sink lifecycle).

Blocks (`internal/flow/blocks/`):
- `mqtt.go` — MQTT source/sink, shared `newMQTTClient` (auto-reconnect),
  connection-lost latch.
- `httpblocks.go` — HTTP source/sink (opt, deadline derivation, JSON body
  rules).
- `transform.go` — Filter (key_prefix/key_contains/type/metadata) + Map.
- `jsonpathmap.go` — JSONPath w/ value extraction + arrays.
- `randomsource.go` — random integer/real/boolean source (min/max/interval).
- `sort.go` — windowed ascending-timestamp sort + dedup.
- `tojson.go`, `splitmap.go`, `updatemetadata.go`, `catalog.go`, `info.go`,
  `schema.go` — catalog wire, field-splitting, metadata updates, registry
  JSON metadata.
- `astartesource/source.go`, `container/block.go` +
  `container/httpbridge.go`, `virtualdevicepool/virtualdevicepool.go` —
  bus-backed source, container block + bridge, device pool.
- Tests sampled: `httpblocks_test.go`, `transform_test.go`,
  `sort_test.go`, `randomsource_test.go`, `mqtt_test.go`, `jsonpathmap_test.go`,
  `astartesource/source_test.go`, `virtualdevicepool_test.go`.

## What I found

Four behavior issues (unguarded/wrong) and one test gap, in value order.

### 1. `mqtt_source` connection-lost latch is never cleared — a transient broker blip kills the source forever

`newMQTTClient` sets `SetAutoReconnect(true)` (mqtt.go:54) and the docstring
promises background auto-reconnect (mqtt.go:35-36). But the source latches
`lost = true` in the ConnectionLost handler (mqtt.go:185-187) and nothing
ever clears it: there is no `SetOnConnectHandler`, and Emit/Process gate on
the sticky flag (mqtt.go:234-236, 252-254). A single transient disconnect
therefore makes `mqtt_source` answer "connection lost" for the rest of the
flow's life, even after paho finishes reconnecting and re-subscribing — the
flow stays `running` while going silently dark. `TestMQTTSource_ConnectionLost`
(mqtt_test.go:190) only asserts the error surfaces; it never proves recovery.

Fix: clear `lost` on a successful reconnect (e.g. `SetOnConnectHandler`, or
check `client.IsConnected()` in Emit/Process) and extend the embedded-mochi
test to stop+restart the broker and assert Emit resumes. Container-free.

$ grep -o '^- \[[x!]\] flow[a-z0-9-]*' .mule/todo.md | sort; echo "=== block slugs ==="; grep -o '^- \[[x!]\] [a-z0-9-]*block[a-z0-9-]*' .mule/todo.md | sort
- [!] flow-mqtt-source-reconnect-recovery
- [!] flow-randomsource-span-overflow
- [!] flow-validate-source-sink
- [x] flow-filter-key-contains-test
- [x] flow-msg-json-integer-precision
- [x] flow-sort-bounded-buffer
- [x] flow-validate-dead-source-sink-recompute
- [x] flowapi-autorestart-default-test
- [x] flowapi-autorestart-terminal-failure
- [x] flowapi-validationdetail-test
=== block slugs ===
- [x] store-pipelines-empty-name-zero-blocks-test

$ rg -c '' --glob 'internal/testutil/**/*.go' | sort -t: -k2 -rn; echo; ls internal/testutil/; echo "=== testutil referenced by ==="; rg -l 'testutil\.' --glob '*.go' | wc -l
internal/testutil/astartedevice.go:216
internal/testutil/mqttclient.go:163
internal/testutil/golden_test.go:117
internal/testutil/pg.go:106
internal/testutil/golden.go:59
internal/testutil/pg_test.go:36

astartedevice.go
golden.go
golden_test.go
mqttclient.go
pg.go
pg_test.go
=== testutil referenced by ===
23

$ sed -n 1,80p .mule/reviews/broker-2026-09-21.md
# Code review: internal/broker — 2026-09-21

Last reviewed: 2026-09-04. Store was reviewed 09-08 and again today, so this
run returned to broker — the most stale package — to look specifically at what
changed since the last pass.

## What changed since 2026-09-04

Two commits landed in internal/broker:

- `d201db4` — the coldstart-introspection-miss fix: the ACL miss path now falls
  back to a synchronous per-interface store read (`sess.syncOwnershipOf`,
  authhook.go:212-237) when the reload debounce skips `refreshIfStale`.
- `24ad5b8` — the `control/keyAgreement` deny comment rewritten to cite
  upstream's 082 spec instead of the deleted 080 sentence. Comment only.

## What I read in detail

- `authhook.go:144-237` — the device-session cache, `refresh`/`refreshIfStale`
  debounce, and the new `syncOwnershipOf` fallback.
- `aclhook.go:95-162` — `offlineACL` and its `entries` map lifecycle.
- `aclhook.go:182-221` — the miss-path wiring that calls the new fallback.
- `broker_test.go:434-497` — the T1 test `TestBrokerACLColdStartIntrospectionMiss`
  that locks the new behaviour, and the fake store (GetDevice/GetInterface).
- Re-read the seams the 09-04 review already covered to confirm they are
  unchanged: `intakeHook` deferred-ack, `Publisher` expiry, `lifecycleHook`
  takeover guards, `sessionstore` Seq normalization. All still as reviewed;
  their queued tasks are still open/blocked.

## Findings (proposed)

### 1. The cold-start fallback defeats the reload debounce — one store read per attacker-chosen name, unbounded in a window

`d201db4` gave the ACL miss path a synchronous fallback
(`syncOwnershipOf`, authhook.go:212-237) that, for **each distinct interface
name** the topic carries, does a full synchronous `GetDevice` + `GetInterface`
store read. The `introspectionReloadDebounce` (1s, authhook.go:42-45, "so an
adversarial topic flood cannot hammer the database") only throttles
`refreshIfStale`'s full-cache reload; the new fallback is on a rawer path:

- `aclhook.go:196-206`: miss → `refreshIfStale` returns false (debounce
  pending) → `syncOwnershipOf` runs **unconditionally**, once per name.
- `syncOwnershipOf` writes only that one name into the session map, so repeat
  publishes to the *same* name are cheap after the first hit — but an
  authenticated device publishing QoS0 to N **distinct** bogus/new interface
  names in one second causes N synchronous full-device store reads. The window
  reopens every second via the debounced full reload, so this is a standing
  per-session throughput of ~N reads/sec, attacker-controlled (QoS0 rate is
  unbounded apart from MaxPacketBytes).

Before the fix a denied miss cost nothing (the debounce-skipped reload re-checked
a cold map and denied); the fix traded a specific cold-start race for
re-opening the exact anti-flood property the debounce comment (authhook.go:
44-45) exists to guarantee. The fix's own comment ("the per-interface store
read happens once per introspection refresh cycle") is only true for names
already seen.

Proposal: throttle the synchronous fallback to the same budget as the debounce
— one sync store read per session per `introspectionReloadDebounce` window
(determine policy: fold `syncOwnershipOf` into a single full reload that also
stamps `lastIntroLoad`, or gate the fallback on its own per-window stamp).
Acceptable trade: a heavy device flooding distinct names within a window may
get the extra ones denied until the debounced reload (≤1s later; the reload
repopulates all real introspected names). The existing
`TestBrokerACLColdStartIntrospectionMiss` (one name) must stay green.

**Machine-checkable:** a T1 test in broker_test.go using the existing fake
store with a GetDevice-call counter: connect cold (stamp `lastIntroLoad` fresh
as the existing test does), publish to K ≥ 2 distinct **unknown** interface
names, and assert the fallback caused at most one more GetDevice (not K).

### 2. `offlineACL.entries` is append-only — unbounded growth across device churn

`offlineACL` (aclhook.go:100-162) keeps one `offlineEntry` per client CN ever
ACL-checked offline. `ownershipOf` upserts on every miss; nothing ever deletes
from `o.entries`, so over a long-lived broker the map grows one entry per
distinct CN seen, never reaped even after the 10s TTL expires (the old entry is
kept and its payload re-stamped). Bounded only by total devices ever registered
in the process lifetime.

→ Read internal/testutil/golden.go
→ Read internal/testutil/golden_test.go
Picking an area: `internal/testutil` is the only package with no review file, and it is small enough to read whole.
→ Read internal/testutil/pg.go
→ Read internal/testutil/astartedevice.go
→ Read internal/testutil/mqttclient.go

$ rg -n 'DeflateControlList|InflateControlList|Introspection\(' --glob '*.go' | grep -v 'internal/testutil/' | head -30; echo "=== testutil tests ==="; cat internal/testutil/pg_test.go
test/conformance/atomvm/atomvm_test.go:45:	dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{
test/conformance/cpc/cpc_test.go:318:	dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{valuesName: {1, 0}}))
test/conformance/load/load_test.go:51:		d.PublishIntrospection(t, testutil.Introspection(map[string][2]int{sensor: {1, 0}}))
internal/engine/data.go:173:			e.onIntrospection(ctx, m, realm)
test/security/security_test.go:295:	dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{"org.astrate.sec.Conf": {1, 0}}))
internal/engine/engine.go:26:	RefreshIntrospection(ctx context.Context, realm string, id deviceid.ID) error
internal/engine/engine.go:40:func (f brokerFacade) RefreshIntrospection(ctx context.Context, realm string, id deviceid.ID) error {
internal/engine/engine.go:41:	return f.b.RefreshIntrospection(ctx, realm, id)
internal/engine/engine_integration_test.go:152:	if _, err := raw.UpdateIntrospection(ctx, realm.ID, rig.dev, intro); err != nil {
internal/appengine/formats_parity_test.go:215:	if _, err := r.st.UpdateIntrospection(ctx, r.realmID, dev, map[string]store.InterfaceVersion{
internal/appengine/formats_parity_test.go:243:	if _, err := r.st.UpdateIntrospection(ctx, r.realmID, dev, map[string]store.InterfaceVersion{
internal/appengine/e2e_test.go:170:		dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{
internal/appengine/mirror_test.go:147:	if _, err := r.st.UpdateIntrospection(ctx, r.realmID, stranger, map[string]store.InterfaceVersion{
internal/appengine/http_test.go:107:	if _, err := st.UpdateIntrospection(ctx, realm.ID, dev, map[string]store.InterfaceVersion{
internal/appengine/http_test.go:242:		if _, err := r.st.UpdateIntrospection(ctx, r.realmID, dev, map[string]store.InterfaceVersion{
internal/engine/introspection.go:35:func (e *Engine) handleIntrospection(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
internal/engine/introspection.go:42:	intro, err := parseIntrospection(raw)
internal/engine/introspection.go:57:		removed, err = e.st.UpdateIntrospection(ctx, realm.id, m.DeviceID, intro)
internal/engine/introspection.go:62:	dev.setIntrospection(intro)
internal/engine/introspection.go:69:		if err := e.broker.RefreshIntrospection(ctx, m.Realm, m.DeviceID); err != nil {
internal/engine/introspection.go:99:func parseIntrospection(s string) (map[string]store.InterfaceVersion, error) {
internal/engine/introspection_test.go:17:func TestParseIntrospection(t *testing.T) {
internal/engine/introspection_test.go:59:			got, err := parseIntrospection(tc.in)
internal/engine/introspection_test.go:62:					t.Fatalf("parseIntrospection(%q) accepted, want error", tc.in)
internal/engine/introspection_test.go:67:				t.Fatalf("parseIntrospection(%q): %v", tc.in, err)
internal/engine/introspection_test.go:70:				t.Fatalf("parseIntrospection(%q) = %v, want %v", tc.in, got, tc.want)
internal/engine/triggers_change_integration_test.go:38:	_, err = r.raw.UpdateIntrospection(ctx, r.realmID, dev, intro)
internal/engine/fakestore_test.go:392:func (f *fakeStore) UpdateIntrospection(_ context.Context, realmID int16, id deviceid.ID, intro map[string]store.InterfaceVersion) (map[string]store.InterfaceVersion, error) {
internal/engine/fakestore_test.go:751:func (p *fakePort) RefreshIntrospection(_ context.Context, realm string, id deviceid.ID) error {
internal/engine/cache.go:61:	UpdateIntrospection(ctx context.Context, realmID int16, id deviceid.ID, intro map[string]store.InterfaceVersion) (map[string]store.InterfaceVersion, error)
=== testutil tests ===
//go:build integration

package testutil

import (
	"context"
	"testing"
	"time"
)

// TestStartTimescale is the M0 T2 gate (docs/ROADMAP.md §1.2 file 0.7): the
// production-parity container boots and the timescaledb extension is present.
func TestStartTimescale(t *testing.T) {
	pool := StartTimescale(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var version string
	err := pool.QueryRow(ctx,
		"SELECT extversion FROM pg_extension WHERE extname = 'timescaledb'",
	).Scan(&version)
	if err != nil {
		t.Fatalf("querying timescaledb extension: %v", err)
	}
	if version == "" {
		t.Fatal("timescaledb extension version is empty")
	}
	t.Logf("timescaledb extension version %s on %s", version, TimescaleImage)

	// The pool must be usable for ordinary work, not just the ready ping.
	var one int
	if err := pool.QueryRow(ctx, "SELECT 1").Scan(&one); err != nil || one != 1 {
		t.Fatalf("SELECT 1 = %d, err %v", one, err)
	}
}

$ rg -n 'zlib|Deflate|producer/properties|consumer/properties' --glob 'internal/engine/*.go' | head -30; echo "=== engine control inflate ==="; rg -n 'func.*[Ii]nflate|zlib.NewReader|LimitReader' --glob 'internal/engine/**/*.go' --glob 'internal/broker/*.go'
internal/engine/capabilities.go:19:	compressionZlib      = "zlib"
internal/engine/capabilities.go:27:// server sends consumer/properties payloads compressed (zlib, the default) or
internal/engine/capabilities.go:89:// purgeCompressionFor returns the consumer/properties compression format for
internal/engine/capabilities.go:90:// the device: "zlib" or "plaintext". The store fallback covers the case
internal/engine/data.go:44:	// reasonControlInvalid: a control payload with a bad zlib frame, a
internal/engine/topics_test.go:55:		{"control/producer/properties", kindControl, "producer/properties"},
internal/engine/topics.go:61:// "producer/properties"); a bare "control" yields an empty subpath, which
internal/engine/engine.go:329:// events, and — on connect — the asynchronous `consumer/properties` purge
internal/engine/engine.go:360:				e.log.Warn("consumer/properties send after connect failed",
internal/engine/control_test.go:5:	"compress/zlib"
internal/engine/control_test.go:20:// goldenFrame loads the frozen producer/properties blob: framed with an
internal/engine/control_test.go:21:// independent zlib implementation (CPython), exactly like an official SDK
internal/engine/control_test.go:22:// would produce it (docs/ROADMAP.md §7.3 "zlib golden producer/properties
internal/engine/control_test.go:92:		zw := zlib.NewWriter(&buf)
internal/engine/control_test.go:114:		{name: "not zlib", in: frame(8, []byte("garbage!"))},
internal/engine/control_test.go:132:// data topic (QoS 2, hint format) followed by the consumer/properties purge
internal/engine/control_test.go:182:	purges := port.publishedTo(base + "/control/consumer/properties")
internal/engine/control_test.go:184:		t.Fatalf("consumer/properties messages: %d, want 1", len(purges))
internal/engine/control_test.go:217:// emptyCache resends as JSON documents; the control frame stays zlib
internal/engine/control_test.go:250:	purges := port.publishedTo(base + "/control/consumer/properties")
internal/engine/control_test.go:252:		t.Fatalf("consumer/properties messages: %d, want 1", len(purges))
internal/engine/control_test.go:255:		t.Errorf("control frame must stay zlib for JSON devices: %v", err)
internal/engine/control_test.go:282:	rig.handle(deviceMsg("control", "/producer/properties", 2, frame, ack))
internal/engine/control_test.go:284:		t.Fatal("producer/properties not acknowledged")
internal/engine/control_test.go:314:	rig.handle(deviceMsg("control", "/producer/properties", 2, []byte{0, 0}, ack))
internal/engine/capabilities_test.go:30:			name:    "zlib compression capability",
internal/engine/capabilities_test.go:31:			payload: mustBSON(t, bson.M{"purge_properties_compression_format": "zlib"}),
internal/engine/capabilities_test.go:32:			want:    map[string]string{"purge_properties_compression_format": "zlib"},
internal/engine/capabilities_test.go:46:			payload: mustBSON(t, bson.M{"purge_properties_compression_format": "zlib", "numeric": 42}),
internal/engine/capabilities_test.go:47:			want:    map[string]string{"purge_properties_compression_format": "zlib"},
=== engine control inflate ===
internal/engine/control.go:291:func inflateProperties(p []byte) ([]string, error) {
internal/engine/control.go:299:	zr, err := zlib.NewReader(bytes.NewReader(p[controlFrameHeader:]))
internal/engine/control.go:305:	plain, err := io.ReadAll(io.LimitReader(zr, int64(declared)+1))
internal/engine/forward/http.go:173:		b, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
internal/engine/forward/http.go:184:	if _, err := io.Copy(io.Discard, io.LimitReader(resp.Body, maxDrain)); err != nil {
internal/engine/triggers/actions.go:626:	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
internal/engine/control_test.go:89:func TestInflatePropertiesRejects(t *testing.T) {
internal/engine/engine_e2e_test.go:427:	zr, err := zlib.NewReader(bytes.NewReader(frame[4:]))

$ rg -n 'InflateControlList|DeflateControlList|MQTTTryConnect|WaitForTopic|ServerTLSCert|DeviceCSR|DeviceTLSConfig|WaitToken|MQTTConnect' --glob '*.go' | grep -v '^internal/testutil/'
cmd/astrate/main.go:304:		bcfg.ServerTLSCert = cert
cmd/astrate/main.go:314:		bcfg.ServerTLSCert = cert
test/conformance/instance/instance.go:152:	serverCert, roots := testutil.ServerTLSCert(t)
test/conformance/instance/instance.go:154:		TLSAddr: "127.0.0.1:0", ServerTLSCert: serverCert,
test/conformance/instance/instance.go:210:	key, csr := testutil.DeviceCSR(t)
test/conformance/instance/instance.go:216:	return testutil.DeviceTLSConfig(t, crt, key, in.Roots)
test/conformance/cpb/cpb_test.go:186:	serverCert, _ := testutil.ServerTLSCert(t)
test/conformance/cpb/cpb_test.go:189:		ServerTLSCert:    serverCert,
test/conformance/atomvm/atomvm_test.go:95:		msg := dev.WaitForTopic(t, 15*time.Second, dev.Base()+"/"+ifServerData+"/value")
test/conformance/cpc/cpc_test.go:188:	serverCert, roots := testutil.ServerTLSCert(t)
test/conformance/cpc/cpc_test.go:191:		ServerTLSCert:    serverCert,
test/conformance/cpc/cpc_test.go:310:	key, csr := testutil.DeviceCSR(t)
test/conformance/cpc/cpc_test.go:315:	tlsCfg := testutil.DeviceTLSConfig(t, crt, key, f.roots)
test/security/security_test.go:114:	serverCert, roots := testutil.ServerTLSCert(t)
test/security/security_test.go:116:		TLSAddr: "127.0.0.1:0", ServerTLSCert: serverCert,
test/security/security_test.go:287:	key, csr := testutil.DeviceCSR(t)
test/security/security_test.go:292:	tlsCfg := testutil.DeviceTLSConfig(t, crt, key, env.roots)
test/security/security_test.go:313:	testutil.WaitToken(t, dev.Client.Publish(dev.Base()+"/control/producer/properties", 2, false, frame), 5*time.Second)
internal/engine/engine_e2e_test.go:140:	serverCert, roots := testutil.ServerTLSCert(t)
internal/engine/engine_e2e_test.go:143:		ServerTLSCert:    serverCert,
internal/engine/engine_e2e_test.go:192:	key, csrPEM := testutil.DeviceCSR(t)
internal/engine/engine_e2e_test.go:197:	tlsCfg := testutil.DeviceTLSConfig(t, certPEM, key, env.roots)
internal/engine/engine_e2e_test.go:333:	msg := dev.WaitForTopic(t, 5*time.Second, topic)
internal/engine/engine_e2e_test.go:362:	dev.WaitForTopic(t, 5*time.Second, propTopic) // initial retained set
internal/engine/engine_e2e_test.go:367:	msg := dev.WaitForTopic(t, 5*time.Second, propTopic)
internal/broker/broker.go:42:	// ServerTLSCert is the broker's server-side TLS identity, required for
internal/broker/broker.go:45:	ServerTLSCert tls.Certificate
internal/broker/broker.go:107:	if len(cfg.ServerTLSCert.Certificate) == 0 {
internal/broker/broker.go:108:		return nil, errors.New("broker: ServerTLSCert is required for the TLS listener")
internal/broker/broker.go:128:	b.pools = newRealmPools(st, cfg.ServerTLSCert, log)
internal/broker/broker_test.go:198:	cert, _ := testutil.ServerTLSCert(t)
internal/broker/broker_test.go:200:	if _, err := New(ctx, Config{SessionStorePath: "x", ServerTLSCert: cert}, st, nil, nil); err == nil {
internal/broker/broker_test.go:203:	if _, err := New(ctx, Config{ServerTLSCert: cert}, st, intake, nil); err == nil {
internal/broker/broker_test.go:207:		t.Error("New without ServerTLSCert: expected error")
internal/broker/broker_test.go:239:	_, csrPEM := testutil.DeviceCSR(t)
internal/broker/broker_test.go:254:	serverCert, roots := testutil.ServerTLSCert(t)
internal/broker/broker_test.go:258:		ServerTLSCert:    serverCert,
internal/broker/broker_test.go:270:	devKeyPriv, csrPEM := testutil.DeviceCSR(t)
internal/broker/broker_test.go:275:	tlsCfg := testutil.DeviceTLSConfig(t, certPEM, devKeyPriv, roots)
internal/broker/broker_test.go:278:	client, sessionPresent := testutil.MQTTConnect(t, url, identity.CN(), true, tlsCfg)
internal/broker/broker_test.go:299:		testutil.WaitToken(t, token, 5*time.Second)
internal/broker/broker_test.go:305:		testutil.WaitToken(t, token, 5*time.Second)
internal/broker/broker_test.go:318:		testutil.WaitToken(t, token, 5*time.Second)
internal/broker/broker_test.go:322:		testutil.WaitToken(t, pubToken, 5*time.Second)
internal/broker/broker_test.go:388:	serverCert, roots := testutil.ServerTLSCert(t)
internal/broker/broker_test.go:392:		ServerTLSCert:    serverCert,
internal/broker/broker_test.go:404:	devKeyPriv, csrPEM := testutil.DeviceCSR(t)
internal/broker/broker_test.go:409:	tlsCfg := testutil.DeviceTLSConfig(t, certPEM, devKeyPriv, roots)
internal/broker/broker_test.go:414:		client, _ := testutil.MQTTConnect(t, url, clientID, true, tlsCfg)
internal/broker/broker_test.go:422:		testutil.WaitToken(t, token, 5*time.Second)
internal/broker/broker_test.go:455:	serverCert, roots := testutil.ServerTLSCert(t)
internal/broker/broker_test.go:459:		ServerTLSCert:    serverCert,
internal/broker/broker_test.go:471:	devKeyPriv, csrPEM := testutil.DeviceCSR(t)
internal/broker/broker_test.go:476:	tlsCfg := testutil.DeviceTLSConfig(t, certPEM, devKeyPriv, roots)
internal/broker/broker_test.go:477:	client, _ := testutil.MQTTConnect(t, "ssl://"+b.TLSAddr(), identity.CN(), true, tlsCfg)
internal/broker/broker_test.go:501:	testutil.WaitToken(t, token, 5*time.Second)
internal/broker/broker_test.go:529:	serverCert, roots := testutil.ServerTLSCert(t)
internal/broker/broker_test.go:533:		ServerTLSCert:    serverCert,
internal/broker/broker_test.go:545:	devKeyPriv, csrPEM := testutil.DeviceCSR(t)
internal/broker/broker_test.go:550:	tlsCfg := testutil.DeviceTLSConfig(t, certPEM, devKeyPriv, roots)
internal/broker/broker_test.go:551:	client, _ := testutil.MQTTConnect(t, "ssl://"+b.TLSAddr(), identity.CN(), true, tlsCfg)
internal/broker/broker_test.go:574:		testutil.WaitToken(t, token, 5*time.Second)
internal/broker/broker_test.go:594:	serverCert, roots := testutil.ServerTLSCert(t)
internal/broker/broker_test.go:598:		ServerTLSCert:    serverCert,
internal/broker/broker_test.go:621:	devKeyPriv, csrPEM := testutil.DeviceCSR(t)
internal/broker/broker_test.go:626:	tlsCfg := testutil.DeviceTLSConfig(t, certPEM, devKeyPriv, roots)
internal/broker/broker_test.go:631:	if _, _, err := testutil.MQTTTryConnect(t, url, cn, true, tlsCfg); err == nil {
internal/broker/broker_test.go:637:	if _, _, err := testutil.MQTTTryConnect(t, url, cn, true, tlsCfg); err != nil {
internal/appengine/e2e_test.go:123:	serverCert, roots := testutil.ServerTLSCert(t)
internal/appengine/e2e_test.go:125:		TLSAddr: "127.0.0.1:0", ServerTLSCert: serverCert,
internal/appengine/e2e_test.go:162:		key, csr := testutil.DeviceCSR(t)
internal/appengine/e2e_test.go:167:		tlsCfg := testutil.DeviceTLSConfig(t, crt, key, roots)
internal/appengine/e2e_test.go:196:		dev.WaitForTopic(t, 5*time.Second, dev.Base()+"/"+cdServerData+"/value")
internal/appengine/e2e_test.go:213:		if _, _, err := testutil.MQTTTryConnect(t, env.sslURL, cn, true, dev.tlsCfg); err == nil {
internal/appengine/e2e_test.go:219:		_, csr := testutil.DeviceCSR(t)
internal/broker/broker_e2e_test.go:123:	serverCert, roots := testutil.ServerTLSCert(t)
internal/broker/broker_e2e_test.go:150:	key, csrPEM := testutil.DeviceCSR(t)
internal/broker/broker_e2e_test.go:164:		tlsCfg:   testutil.DeviceTLSConfig(t, certPEM, key, e.roots),
internal/broker/broker_e2e_test.go:200:		ServerTLSCert:     e.serverCert,
internal/broker/broker_e2e_test.go:240:	client, sessionPresent := testutil.MQTTConnect(t, sslURL(b), dev.identity.CN(), false, dev.tlsCfg)
internal/broker/broker_e2e_test.go:286:		key, csrPEM := testutil.DeviceCSR(t)
internal/broker/broker_e2e_test.go:291:		cfg := testutil.DeviceTLSConfig(t, certPEM, key, env.roots)
internal/broker/broker_e2e_test.go:292:		if _, _, err := testutil.MQTTTryConnect(t, url, dev.identity.CN(), true, cfg); err == nil {
internal/broker/broker_e2e_test.go:299:		key, csrPEM := testutil.DeviceCSR(t)
internal/broker/broker_e2e_test.go:304:		cfg := testutil.DeviceTLSConfig(t, certPEM, key, env.roots)
internal/broker/broker_e2e_test.go:306:		if _, _, err := testutil.MQTTTryConnect(t, url, cn, true, cfg); err == nil {
internal/broker/broker_e2e_test.go:319:		client, _, err := testutil.MQTTTryConnect(t, url, spoofCN, true, dev.tlsCfg)
internal/broker/broker_e2e_test.go:339:		if _, _, err := testutil.MQTTTryConnect(t, url, dev.identity.CN(), true, dev.tlsCfg); err == nil {
internal/broker/broker_e2e_test.go:348:		key2, csrPEM2 := testutil.DeviceCSR(t)
internal/broker/broker_e2e_test.go:354:		if _, _, err := testutil.MQTTTryConnect(t, url, dev.identity.CN(), true, dev.tlsCfg); err == nil {
internal/broker/broker_e2e_test.go:357:		freshCfg := testutil.DeviceTLSConfig(t, certPEM2, key2, env.roots)
internal/broker/broker_e2e_test.go:358:		if _, _, err := testutil.MQTTTryConnect(t, url, dev.identity.CN(), true, freshCfg); err != nil {
internal/broker/broker_e2e_test.go:364:		if _, _, err := testutil.MQTTTryConnect(t, sslURL(lenient), dev.identity.CN(), true, dev.tlsCfg); err != nil {
internal/broker/broker_e2e_test.go:373:		if _, _, err := testutil.MQTTTryConnect(t, "tcp://"+b.TLSAddr(), dev.identity.CN(), true, nil); err == nil {
internal/broker/broker_e2e_test.go:386:	if _, _, err := testutil.MQTTTryConnect(t, "tcp://"+b.DevAddr(), dev.identity.CN(), true, nil); err != nil {
internal/broker/broker_e2e_test.go:400:	client, _ := testutil.MQTTConnect(t, sslURL(b), dev.identity.CN(), true, dev.tlsCfg)
internal/broker/broker_e2e_test.go:412:	testutil.WaitToken(t, token, 5*time.Second)
internal/broker/broker_e2e_test.go:423:	client, sessionPresent := testutil.MQTTConnect(t, sslURL(b1), dev.identity.CN(), false, dev.tlsCfg)
internal/broker/broker_e2e_test.go:427:	testutil.WaitToken(t, client.Subscribe(subTopic, 1, nil), 5*time.Second)
internal/broker/broker_e2e_test.go:452:	_, sessionPresent = testutil.MQTTConnect(t, sslURL(b2), dev.identity.CN(), false, dev.tlsCfg, collector)
internal/broker/broker_e2e_test.go:486:	client, _ := testutil.MQTTConnect(t, sslURL(b), dev.identity.CN(), true, dev.tlsCfg, collector)
internal/broker/broker_e2e_test.go:492:	testutil.WaitToken(t, client.Publish(topic, 2, false, []byte("qos2-up")), 10*time.Second)
internal/broker/broker_e2e_test.go:501:	testutil.WaitToken(t, client.Subscribe(subTopic, 2, nil), 5*time.Second)
internal/broker/broker_e2e_test.go:528:	client, _ := testutil.MQTTConnect(t, sslURL(b1), dev.identity.CN(), false, dev.tlsCfg)
internal/broker/broker_e2e_test.go:529:	testutil.WaitToken(t, client.Subscribe(subTopic, 1, nil), 5*time.Second)
internal/broker/broker_e2e_test.go:555:	_, sessionPresent := testutil.MQTTConnect(t, sslURL(b2), dev.identity.CN(), false, dev.tlsCfg, collector)
internal/broker/broker_e2e_test.go:582:	c1, _ := testutil.MQTTConnect(t, url, dev.identity.CN(), false, dev.tlsCfg)
internal/broker/broker_e2e_test.go:588:	c2, _ := testutil.MQTTConnect(t, url, dev.identity.CN(), false, dev.tlsCfg)

$ sed -n 95,140p internal/engine/introspection.go; echo "=== engine_e2e uses testutil? ==="; rg -n 'testutil\.' internal/engine/engine_e2e_test.go | head
// parseIntrospection parses the `;`-separated `name:major:minor` triples
// (docs/DESIGN.md §3.3). The empty payload is a valid empty introspection
// (the device declares no interfaces); any malformed entry rejects the whole
// message — partial introspections would desynchronize the ACL state.
func parseIntrospection(s string) (map[string]store.InterfaceVersion, error) {
	intro := map[string]store.InterfaceVersion{}
	if s == "" {
		return intro, nil
	}
	for _, entry := range strings.Split(s, ";") {
		name, versions, ok := strings.Cut(entry, ":")
		if !ok || name == "" {
			return nil, fmt.Errorf("malformed introspection entry %q", entry)
		}
		if strings.ContainsAny(name, "/+#") {
			return nil, fmt.Errorf("interface name %q contains topic metacharacters", name)
		}
		majorStr, minorStr, ok := strings.Cut(versions, ":")
		if !ok {
			return nil, fmt.Errorf("malformed introspection entry %q", entry)
		}
		major, err := strconv.Atoi(majorStr)
		if err != nil || major < 0 {
			return nil, fmt.Errorf("invalid major version in introspection entry %q", entry)
		}
		minor, err := strconv.Atoi(minorStr)
		if err != nil || minor < 0 {
			return nil, fmt.Errorf("invalid minor version in introspection entry %q", entry)
		}
		intro[name] = store.InterfaceVersion{Major: major, Minor: minor}
	}
	return intro, nil
}

// retryStore runs fn with shard-parking retries for transient store
// failures, mirroring the data path's DB-outage behaviour (docs/DESIGN.md
// §5.3): while parked the shard does not consume its channel, so the §1.4
// backpressure builds. It reports whether fn eventually succeeded; on every
// false return the message has been fully handled (rejected and consumed, or
// abandoned unacknowledged at shutdown).
func (e *Engine) retryStore(ctx context.Context, m broker.InboundMessage, what string, fn func() error) bool {
	backoff := parkBackoffStart
	for {
		err := fn()
		if err == nil {
			return true
=== engine_e2e uses testutil? ===
85:	pool := testutil.StartTimescale(t)
140:	serverCert, roots := testutil.ServerTLSCert(t)
180:func (env *e2eEnv) connectDevice(t *testing.T) *testutil.AstarteDevice {
192:	key, csrPEM := testutil.DeviceCSR(t)
197:	tlsCfg := testutil.DeviceTLSConfig(t, certPEM, key, env.roots)
199:	dev := testutil.ConnectAstarteDevice(t, env.sslURL, env.realm.Name, id, tlsCfg, true)
202:	dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{
237:func (env *e2eEnv) series(t *testing.T, dev *testutil.AstarteDevice, iface, path string) []store.IndividualRow {
374:	cp := dev.WaitForMessage(t, 5*time.Second, "consumer/properties listing "+want, func(m testutil.ServerMessage) bool {

$ sed -n 350,395p internal/engine/engine_e2e_test.go; echo "=== e2e property unset usage ==="; rg -n 'WaitForMessage|WaitForTopic' --glob '*.go' | grep -v testutil/
// makes Astrate re-send that property on the data topic and publish the
// consumer/properties purge list — both decoded and asserted (docs/DESIGN.md
// §3.3–3.4).
func testE2EEmptyCache(t *testing.T, env *e2eEnv) {
	ctx := context.Background()
	dev := env.connectDevice(t)

	if err := env.e.PublishServerValue(ctx, env.realm.Name, dev.ID, e2eSrvCfg, "/brightness",
		json.RawMessage("42"), nil); err != nil {
		t.Fatalf("PublishServerValue: %v", err)
	}
	propTopic := dev.Base() + "/" + e2eSrvCfg + "/brightness"
	dev.WaitForTopic(t, 5*time.Second, propTopic) // initial retained set

	dev.EmptyCache(t)

	// The resync re-publishes the property value on its data topic...
	msg := dev.WaitForTopic(t, 5*time.Second, propTopic)
	if dp := env.decodeIndividual(t, e2eSrvCfg, 1, "/brightness", msg.Payload); dp.Value != int32(42) {
		t.Errorf("resent property value = %v (%T), want int32(42)", dp.Value, dp.Value)
	}

	// ...and publishes the consumer/properties purge list naming it.
	want := e2eSrvCfg + "/brightness"
	cp := dev.WaitForMessage(t, 5*time.Second, "consumer/properties listing "+want, func(m testutil.ServerMessage) bool {
		if m.Topic != dev.Base()+"/control/consumer/properties" {
			return false
		}
		for _, e := range controlEntries(m.Payload) {
			if e == want {
				return true
			}
		}
		return false
	})
	if got := controlEntries(cp.Payload); len(got) != 1 || got[0] != want {
		t.Errorf("consumer/properties entries = %v, want [%s]", got, want)
	}
}

// testE2EProducerProperties: control/producer/properties is the device's
// exhaustive list of held device-owned properties; Astrate purges every
// device-owned property not in it.
func testE2EProducerProperties(t *testing.T, env *e2eEnv) {
	ctx := context.Background()
	dev := env.connectDevice(t)
=== e2e property unset usage ===
test/conformance/atomvm/atomvm_test.go:95:		msg := dev.WaitForTopic(t, 15*time.Second, dev.Base()+"/"+ifServerData+"/value")
internal/appengine/e2e_test.go:196:		dev.WaitForTopic(t, 5*time.Second, dev.Base()+"/"+cdServerData+"/value")
internal/engine/engine_e2e_test.go:333:	msg := dev.WaitForTopic(t, 5*time.Second, topic)
internal/engine/engine_e2e_test.go:362:	dev.WaitForTopic(t, 5*time.Second, propTopic) // initial retained set
internal/engine/engine_e2e_test.go:367:	msg := dev.WaitForTopic(t, 5*time.Second, propTopic)
internal/engine/engine_e2e_test.go:374:	cp := dev.WaitForMessage(t, 5*time.Second, "consumer/properties listing "+want, func(m testutil.ServerMessage) bool {

$ rg -n 'func .*Disconnect|\.Disconnect\(\)' --glob '*.go' | grep -v 'internal/testutil/' | head; echo "=== flag update ==="; rg -n 'flag\.(Bool|String)\("update' --glob '*.go'
internal/store/devices.go:308:func (s *Store) SetDeviceDisconnected(ctx context.Context, realmID int16, id deviceid.ID, at time.Time) error {
internal/broker/broker.go:260:func (b *Broker) DisconnectDevice(realm string, id deviceid.ID) {
internal/broker/broker_test.go:173:func (f *fakeStore) SetDeviceDisconnected(_ context.Context, realmID int16, id deviceid.ID, _ time.Time) error {
internal/broker/sessionstore.go:154:func (s *sessionStore) OnDisconnect(cl *mqtt.Client, _ error, expire bool) {
internal/broker/lifecycle.go:66:func (h *lifecycleHook) OnDisconnect(cl *mqtt.Client, _ error, _ bool) {
internal/engine/triggers/events.go:212:func NewDeviceDisconnectedEvent() DeviceDisconnectedEvent {
internal/realm/dashboard_compat_test.go:21:func (f *fakeDisconnecter) DisconnectDevice(realm string, id deviceid.ID) {
internal/realm/service.go:100:func (s *Service) WithDisconnecter(d Disconnecter) *Service {
internal/appengine/e2e_test.go:211:		dev.Disconnect()
=== flag update ===
internal/testutil/golden.go:14:var updateGolden = flag.Bool("update", false, "rewrite golden files with the received content")

$ sed -n 1,30p internal/engine/engine_e2e_test.go; echo "=== tags on e2e files ==="; rg -l 'go:build' --glob 'internal/engine/*_test.go' | while read f; do echo "-- $f"; head -3 "$f"; done
//go:build e2e

package engine

import (
	"bytes"
	"compress/zlib"
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/netip"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/pairing"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// Interface names installed in the e2e realm.
=== tags on e2e files ===
-- internal/engine/engine_integration_test.go
//go:build integration

package engine
-- internal/engine/triggers_change_integration_test.go
//go:build integration

package engine
-- internal/engine/engine_e2e_test.go
//go:build e2e

package engine

$ grep -n 'WaitForMessage\|Messages()' -r --glob '*.go' | grep -v 'internal/testutil/'
grep: unrecognized option '--glob'
Usage: grep [OPTION]... PATTERNS [FILE]...
Try 'grep --help' for more information.

✱ Grep "WaitForMessage|\.Messages\(\)" 5 matches

$ for f in internal/appengine/e2e_test.go test/security/security_test.go test/conformance/cpc/cpc_test.go test/conformance/atomvm/atomvm_test.go test/conformance/load/load_test.go; do printf '%s: ' "$f"; head -1 "$f"; done; echo "=== does any non-tagged test call testutil.Introspection/Deflate? ==="; rg -n 'testutil\.Introspection|testutil\.DeflateControlList|testutil\.InflateControlList' --glob '*.go'
internal/appengine/e2e_test.go: //go:build integration && e2e
test/security/security_test.go: //go:build integration && e2e
test/conformance/cpc/cpc_test.go: // Package cpc is conformance checkpoint CP-C (docs/ROADMAP.md §0.3, §8.3): the
test/conformance/atomvm/atomvm_test.go: // Package atomvm is the CP-D AtomVM-profile conformance runner (docs/ROADMAP.md
test/conformance/load/load_test.go: // Package load is the T5 load/footprint smoke (docs/ROADMAP.md §10 file 9.7).
=== does any non-tagged test call testutil.Introspection/Deflate? ===
test/conformance/atomvm/atomvm_test.go:45:	dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{
test/conformance/cpc/cpc_test.go:318:	dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{valuesName: {1, 0}}))
test/conformance/load/load_test.go:51:		d.PublishIntrospection(t, testutil.Introspection(map[string][2]int{sensor: {1, 0}}))
test/security/security_test.go:295:	dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{"org.astrate.sec.Conf": {1, 0}}))
internal/appengine/e2e_test.go:170:		dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{
internal/engine/engine_e2e_test.go:202:	dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{

$ rg -n 'Introspection|Deflate|controlFrameHeader' internal/engine/introspection_test.go internal/engine/control_test.go | head -30; echo "=== control.go 285-315 ==="; sed -n 285,315p internal/engine/control.go
internal/engine/introspection_test.go:14:// TestParseIntrospection covers the docs/DESIGN.md §3.3 introspection
internal/engine/introspection_test.go:17:func TestParseIntrospection(t *testing.T) {
internal/engine/introspection_test.go:59:			got, err := parseIntrospection(tc.in)
internal/engine/introspection_test.go:62:					t.Fatalf("parseIntrospection(%q) accepted, want error", tc.in)
internal/engine/introspection_test.go:67:				t.Fatalf("parseIntrospection(%q): %v", tc.in, err)
internal/engine/introspection_test.go:70:				t.Fatalf("parseIntrospection(%q) = %v, want %v", tc.in, got, tc.want)
internal/engine/introspection_test.go:81:// TestIntrospectionHandler drives the wired handler end to end: persistence,
internal/engine/introspection_test.go:84:func TestIntrospectionHandler(t *testing.T) {
internal/engine/introspection_test.go:97:	persisted := len(dev.Introspection)
internal/engine/introspection_test.go:98:	_, hasMinimal := dev.Introspection["com.astrate.test.Minimal"]
internal/engine/introspection_test.go:131:// TestIntrospectionRejects: oversized and malformed payloads are rejected
internal/engine/introspection_test.go:132:// under reasonIntrospectionInvalid and consumed.
internal/engine/introspection_test.go:133:func TestIntrospectionRejects(t *testing.T) {
internal/engine/introspection_test.go:137:	rig.handle(deviceMsg("", "", 2, []byte(strings.Repeat("x", maxIntrospectionBytes+1)), ack))
internal/engine/introspection_test.go:148:	if got := promtest.ToFloat64(rig.e.met.rejects.WithLabelValues(reasonIntrospectionInvalid)); got != 2 {
internal/engine/introspection_test.go:149:		t.Errorf("rejects[%s] = %v, want 2", reasonIntrospectionInvalid, got)
internal/engine/introspection_test.go:155:	entries := len(dev.Introspection)
internal/engine/introspection_test.go:162:// TestIntrospectionParking: a transient UpdateIntrospection failure parks
internal/engine/introspection_test.go:164:func TestIntrospectionParking(t *testing.T) {
internal/engine/introspection_test.go:184:	entries := len(fs.devices[deviceKey{realm: realmAlpha, id: devAlpha}].Introspection)
=== control.go 285-315 ===
// inflateProperties parses a zlib control payload with the docs/DESIGN.md
// §4.5 bounds: the declared size is capped by maxControlInflated, and the
// stream may not inflate beyond what it declared (a lying header is a
// zip-bomb attempt, not a tolerable client quirk). It returns the entry
// list; an empty payload yields no entries (the SDK deliberately sends an
// empty list when the device holds no properties).
func inflateProperties(p []byte) ([]string, error) {
	if len(p) < controlFrameHeader {
		return nil, fmt.Errorf("control payload is %d bytes, below the %d byte size prefix", len(p), controlFrameHeader)
	}
	declared := binary.BigEndian.Uint32(p[:controlFrameHeader])
	if declared > maxControlInflated {
		return nil, fmt.Errorf("control payload declares %d inflated bytes, above the %d byte ceiling", declared, maxControlInflated)
	}
	zr, err := zlib.NewReader(bytes.NewReader(p[controlFrameHeader:]))
	if err != nil {
		return nil, fmt.Errorf("control payload is not a zlib stream: %w", err)
	}
	defer func() { _ = zr.Close() }()

	plain, err := io.ReadAll(io.LimitReader(zr, int64(declared)+1))
	if err != nil {
		return nil, fmt.Errorf("inflating control payload: %w", err)
	}
	if len(plain) > int(declared) {
		return nil, fmt.Errorf("control payload inflates beyond its declared %d bytes", declared)
	}
	if len(plain) == 0 {
		return nil, nil
	}
	return strings.Split(string(plain), ";"), nil

$ rg -n 'maxControlInflated\s*=' --glob '*.go'; echo "=== zlib.NewReader call sites ==="; rg -n 'zlib.NewReader' --glob '*.go'; echo "=== controlEntries ==="; rg -n 'func controlEntries|func .*[Ii]nflate' --glob '*.go'
internal/engine/control.go:30:const maxControlInflated = 1 << 20
=== zlib.NewReader call sites ===
internal/engine/control.go:299:	zr, err := zlib.NewReader(bytes.NewReader(p[controlFrameHeader:]))
internal/engine/engine_e2e_test.go:427:	zr, err := zlib.NewReader(bytes.NewReader(frame[4:]))
internal/testutil/astartedevice.go:200:	zr, err := zlib.NewReader(bytes.NewReader(frame[4:]))
=== controlEntries ===
internal/engine/control.go:291:func inflateProperties(p []byte) ([]string, error) {
internal/engine/control_test.go:89:func TestInflatePropertiesRejects(t *testing.T) {
internal/engine/engine_e2e_test.go:423:func controlEntries(frame []byte) []string {
internal/testutil/astartedevice.go:194:func InflateControlList(t testing.TB, frame []byte) []string {

$ rg -n 'go:build' test/ | head; echo "=== conformance tags ==="; ls test/
test/security/security_test.go:1://go:build integration && e2e
=== conformance tags ===
conformance
security

$ grep -n 'testutil\|introspection\|control-frame\|WaitFor' .mule/todo.md | head -20; echo "=== done tasks mentioning testutil in log ==="; grep -n 'testutil' .mule/log.md | head
105:- [ ] race-check-pkg: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./pkg/... ./cmd/... ./internal/testutil/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
107:- [x] broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [approved 2026-09-04]
113:- [!] empty-introspection-verification [auto]: upstream v1.3.0 changed "allow devices with empty introspection" — verify whether Astrate's device connection/introspection handling currently rejects an empty introspection string where upstream now accepts it, and propose a fix if so. — BLOCKED: wrote nothing
143:- [x] docs-sync-appengine-device-status-schema [auto]: fix the `DeviceStatus` schema in docs/api/astrate_appengine_api.yaml — the wire emits `id` (internal/appengine/service.go:106 `json:"id"`, decoded as `ds.ID` in mirror_test.go:156 and dashboard_compat_test.go:100/149) but the spec names the field `device_id` (yaml:1409) in both the schema and the list/detail examples (yaml:67, 75); `introspection` is documented as `map<string,integer>` (yaml:1438-1441) but each value is a `{major, minor}` object (introspectionEntry, service.go:123-127, introspectionView service.go:833-836); and the schema omits the six fields `groups`, `total_received_msgs`, `total_received_bytes`, `first_credentials_request`, `last_seen_ip` (omitempty), `previous_interfaces` (omitempty) that deviceStatusView emits (service.go:254-279). Rename `device_id`→`id`, correct `introspection` to `map<string,{major:int,minor:int}>`, add the missing fields, and fix the examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
158:- [x] docs-sync-rm-delete-interface-status [auto]: fix the `deleteInterface` response in docs/api/astarte_realm_management_api.yaml — the spec documents `422` with examples "Interface major version is not 0, can't be deleted" / "Cannot delete an interface that is used by a device introspection" (yaml:255-272), but the handler returns `403` with "Interface can't be deleted" / "Interface can't be deleted since it's currently used" (internal/realm/http.go:387-391, `writeInterfaceError`; both are 403 `StatusForbidden`). The spec's detail strings were lifted from the `writeError` branch (http.go:414-419) which DELETE /interfaces does not use. Fix status to `403` and update the two examples to match the actual wire messages. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
206:- [x] broker-acl-coldstart-fallback-flood [auto]: `syncOwnershipOf` (internal/broker/authhook.go:212-237, reached from the ACL miss path aclhook.go:196-206) does a full synchronous `GetDevice`+`GetInterface` store read for every distinct interface name a device publishes to, with no rate limit — `introspectionReloadDebounce` (authhook.go:42-45, "so an adversarial topic flood cannot hammer the database") throttles only `refreshIfStale`'s full reload, and the fallback's per-name cache write means N distinct bogus names in one window cause N synchronous full-device reads per second, an open anti-flood hole the pre-fix denied-miss path did not have. Gate the fallback to one sync read per session per debounce window (or fold it into a full reload that stamps `lastIntroLoad`), keep `TestBrokerACLColdStartIntrospectionMiss` (one name) green, and add a T1 test in broker_test.go with a GetDevice-counting fake: connect cold, publish to K>=2 distinct unknown-interface names, assert no more than one fallback read was made.
=== done tasks mentioning testutil in log ===

$ grep -rn 'testutil' .mule/reviews/*.md | head; echo "=== race-check-pkg line ==="; grep -n 'race-check-pkg' .mule/todo.md
.mule/reviews/container-2026-10-02.md:4:`internal/testutil` (test helpers, no logic worth proposing on). The 2026-09-23 flow
=== race-check-pkg line ===
105:- [ ] race-check-pkg: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./pkg/... ./cmd/... ./internal/testutil/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]

$ sed -n 1,20p .mule/reviews/container-2026-10-02.md; echo "=== grep all reviews for testutil context ==="; grep -rn -B2 -A3 'internal/testutil' .mule/reviews/*.md
# Review — internal/flow/blocks/container, 2026-10-02

Rotation check first: every `internal/` and `pkg/` package has a review file except
`internal/testutil` (test helpers, no logic worth proposing on). The 2026-09-23 flow
pass read `container/block.go` and `container/httpbridge.go` but **not `container/docker.go`**
(211 lines, the whole Docker-CLI surface: `CLIRunner`, `cliInstance`, `Spec`, `Waiter`,
`encodeFlowConfigJSON`) and drew no findings from this package at all. This pass is that
missing half.

## What I read

- `internal/flow/blocks/container/block.go` (all 320) — `Config`, `New` (Start, bridge
  wiring, `WaitReady`, death-watch goroutine), `watchExit`, `parseConfig`, `asInt`,
  `dockerName`, `Process`, `Stop`.
- `internal/flow/blocks/container/docker.go` (all 211) — `Spec`/`Instance`/`Runner`/
  `Waiter`, `CLIRunner.Start`/`run`/`hostPort`, `cliInstance.{BaseURL,ID,Stop,Wait}`,
  `encodeFlowConfigJSON`. **First review of this file.**
- `internal/flow/blocks/container/httpbridge.go` (all 178) — `Bridge`, `WaitReady`,
  `RoundTrip`, the body cap.
- `internal/flow/blocks/container/block_test.go` (all 424) — `fakeInstance`/`fakeRunner`
=== grep all reviews for testutil context ===
.mule/reviews/container-2026-10-02.md-2-
.mule/reviews/container-2026-10-02.md-3-Rotation check first: every `internal/` and `pkg/` package has a review file except
.mule/reviews/container-2026-10-02.md:4:`internal/testutil` (test helpers, no logic worth proposing on). The 2026-09-23 flow
.mule/reviews/container-2026-10-02.md-5-pass read `container/block.go` and `container/httpbridge.go` but **not `container/docker.go`**
.mule/reviews/container-2026-10-02.md-6-(211 lines, the whole Docker-CLI surface: `CLIRunner`, `cliInstance`, `Spec`, `Waiter`,
.mule/reviews/container-2026-10-02.md-7-`encodeFlowConfigJSON`) and drew no findings from this package at all. This pass is that

$ rg -n 'WaitForTopic|EmptyCache|initial retained|resent property' internal/engine/engine_e2e_test.go
258:	t.Run("EmptyCacheResync", func(t *testing.T) { testE2EEmptyCache(t, env) })
333:	msg := dev.WaitForTopic(t, 5*time.Second, topic)
349:// testE2EEmptyCache: after a server-owned property is set, control/emptyCache
353:func testE2EEmptyCache(t *testing.T, env *e2eEnv) {
362:	dev.WaitForTopic(t, 5*time.Second, propTopic) // initial retained set
364:	dev.EmptyCache(t)
367:	msg := dev.WaitForTopic(t, 5*time.Second, propTopic)
369:		t.Errorf("resent property value = %v (%T), want int32(42)", dp.Value, dp.Value)

$ rg -n 'func \(d \*AstarteDevice\)|func Introspection|func DeflateControlList|func InflateControlList|io.ReadAll' internal/testutil/astartedevice.go; echo "=== introspection_test head ==="; rg -n 'func Test' internal/engine/introspection_test.go
75:func (d *AstarteDevice) Base() string { return d.base }
80:func (d *AstarteDevice) DataTopic(iface, path string) string {
86:func (d *AstarteDevice) PublishIntrospection(t testing.TB, introspection string) {
93:func Introspection(entries map[string][2]int) string {
109:func (d *AstarteDevice) PublishValue(t testing.TB, iface, path string, v payload.Value, ts *time.Time, format payload.Format, qos byte) {
120:func (d *AstarteDevice) PublishRaw(t testing.TB, iface, path string, body []byte, qos byte) {
127:func (d *AstarteDevice) EmptyCache(t testing.TB) {
135:func (d *AstarteDevice) SendProducerProperties(t testing.TB, entries []string) {
141:func (d *AstarteDevice) Messages() []ServerMessage {
149:func (d *AstarteDevice) WaitForMessage(t testing.TB, timeout time.Duration, what string, pred func(ServerMessage) bool) ServerMessage {
165:func (d *AstarteDevice) WaitForTopic(t testing.TB, timeout time.Duration, topic string) ServerMessage {
173:func (d *AstarteDevice) Disconnect() {
181:func DeflateControlList(entries []string) []byte {
194:func InflateControlList(t testing.TB, frame []byte) []string {
205:	plain, err := io.ReadAll(zr)
=== introspection_test head ===
17:func TestParseIntrospection(t *testing.T) {
84:func TestIntrospectionHandler(t *testing.T) {
133:func TestIntrospectionRejects(t *testing.T) {
164:func TestIntrospectionParking(t *testing.T) {

Read complete — all six files of `internal/testutil` plus their callers. Writing the review:
← Write .mule/reviews/testutil-2026-10-08.md
Wrote file successfully.

$ tail -c 600 .mule/todo.md; echo; echo "=== count ==="; wc -l .mule/todo.md
ernal/engine/stream/bus.go Kind* constants (bus.go:20-26), and record that this socket drops `InterfaceMajor`, `IP`, `ErrorName` and `ErrorMetadata`, which `stream.Event` carries (bus.go:39-65) — so a `device_error` arrives here with no reason, unlike the Phoenix twin, which marshals the whole struct (internal/appengine/channels/ws.go:406-416). That last sentence documents a code gap, not a spec one: do not add the fields; if they are wanted on the wire, that is a separate code task. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1

=== count ===
290 .mule/todo.md

$ cat >> .mule/todo.md <<'EOF'
- [ ] testutil-wait-since-cursor [auto]: give AstarteDevice's message waits a cursor so a second wait for the same topic cannot re-match an older message — `WaitForMessage`/`WaitForTopic` (internal/testutil/astartedevice.go:149-170) scan the capture buffer from index 0 on every poll, so a repeated wait returns the first match again, and in-tree that makes testE2EEmptyCache's resync wait vacuous (see the companion line). Move the paho capture handler out of the `ConnectAstarteDevice` closure (astartedevice.go:61-67) into an injectable method, add `Mark() int` and a `From(start int, …)` wait variant, and add a container-free internal/testutil/astartedevice_test.go that drives two captures on one topic and asserts a wait started after the second returns the second message — it cannot even compile today, which is the failure proof. Keep every new read of `received` under the existing `d.mu`, and say in the report that the change is unverified for races (the Pi gate runs no `-race`). The consumer of this API is engine-e2e-emptycache-resync-assert.
- [ ] engine-e2e-emptycache-resync-assert [legion] [auto]: make the second wait in `testE2EEmptyCache` (internal/engine/engine_e2e_test.go:367) wait for the resync re-publish instead of re-matching the retained set captured at :362 — mark the cursor before `dev.EmptyCache(t)` (:364) and take only messages after the mark, so the "resent property value = int32(42)" assertion at :368-369 fails when the resend never happens; today it cannot fail, because `WaitForTopic` returns the pre-EmptyCache message, which carries the same value 42. Needs `-tags e2e` plus Docker, so this half is Legion-only — the container-free helper and its test land in testutil-wait-since-cursor.
- [ ] testutil-control-frame-test [auto]: pin the producer/properties control frame that `SendProducerProperties` puts on the wire (internal/testutil/astartedevice.go:181-216) — `DeflateControlList`/`InflateControlList` have zero tests and every caller of the helpers sits behind an `e2e`/`integration`/conformance build tag, so nothing on the default gate asserts the 4-byte big-endian uncompressed-length prefix plus zlib framing, while the parser side is pinned in internal/engine/control_test.go:89-130 against a CPython golden frame. Add a container-free internal/testutil/astartedevice_test.go with exact-bytes, round-trip, empty-list→nil, sub-4-byte header and declared≠inflated cases; bound the inflate with `io.LimitReader(zr, int64(declared)+1)` the way internal/engine/control.go:301-315 does (astartedevice.go:205 reads with unbounded `io.ReadAll` today, and a lying header from the code under test makes the harness allocate it before failing) and assert the truncated byte count in the mismatch message so the bound itself is pinned. One extra line in internal/engine/control_test.go feeding `testutil.DeflateControlList` into `inflateProperties` is the cross-check that closes the loop.
- [ ] engine-introspection-producer-roundtrip [auto]: pin the producer half of the docs/DESIGN.md §3.3 introspection rule on the default gate — `testutil.Introspection` (internal/testutil/astartedevice.go:93-105) is called from six places, all behind `e2e`, `integration && e2e` or test/conformance build tags, while the consumer `parseIntrospection` (internal/engine/introspection.go:99-127) is tested only against hand-written literals in internal/engine/introspection_test.go:17-71, so nothing container-free checks that the string the harness renders is the one the engine accepts or that the documented "deterministic (sorted) ordering" holds. Add a case to internal/engine/introspection_test.go that builds the string from an unsorted multi-entry map via `testutil.Introspection`, asserts `parseIntrospection` returns the same map, and pins the exact output (sorted, `name:major:minor;…`, no trailing `;`, empty map → empty string).
EOF
tail -5 .mule/todo.md | cut -c1-160; echo "=== open count now ==="; grep -c '^- \[ \]' .mule/todo.md
- [!] docs-native-socket-event-schema [auto]: document the native socket's event payload in docs/api/astrate_native_api.yaml — currently only a fenced `incomi
- [ ] testutil-wait-since-cursor [auto]: give AstarteDevice's message waits a cursor so a second wait for the same topic cannot re-match an older message — `W
- [ ] engine-e2e-emptycache-resync-assert [legion] [auto]: make the second wait in `testE2EEmptyCache` (internal/engine/engine_e2e_test.go:367) wait for the res
- [ ] testutil-control-frame-test [auto]: pin the producer/properties control frame that `SendProducerProperties` puts on the wire (internal/testutil/astartedev
- [ ] engine-introspection-producer-roundtrip [auto]: pin the producer half of the docs/DESIGN.md §3.3 introspection rule on the default gate — `testutil.Int
=== open count now ===
18

Done: code-review recipe run on `internal/testutil` (the only package with no review file), wrote the review and queued 4 proposals.
Files: .mule/reviews/testutil-2026-10-08.md, .mule/todo.md (+4 lines)
Verified: read 100% of internal/testutil (6 files) + all 23 callers; rg for build tags/callers -> no run of go test needed (proposal-only run, no code changed)
Unsure: whether `WaitForTopic`'s vacuity at engine_e2e_test.go:367 matters to you enough to spend a `[legion]` tick — the e2e suite is `-tags e2e`, so only the Legion can prove the fix
Follow-ups: testutil-wait-since-cursor, engine-e2e-emptycache-resync-assert [legion], testutil-control-frame-test, engine-introspection-producer-roundtrip
