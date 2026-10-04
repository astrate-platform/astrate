slug: recipe-docs-sync
verdict: proposed
at:  87c3010
ran: 2026-10-04T11:25:05Z on DietPi in 474s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ tail -30 /root/astrate-mule/.mule/todo.md
- [x] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml (yaml:939-980, currently 201/400/401/404/409/500) — a well-formed body whose `device_id` does not parse is rejected as a field error (`FieldErrors{"device_id": {"is not a valid device id"}}`, internal/appengine/service.go:631-634) and answered 422 with the `FieldErrorsDetail` body (`astarteapi.WriteFieldErrors`, internal/appengine/http.go:647), so `$ref` the existing `ValidationErrors` response (yaml:1707-1716) rather than `BadRequest` — the documented 400 only covers a body `DecodeData` cannot parse (internal/appengine/http.go:488-491). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data DELETE operations in docs/api/astarte_appengine_api.yaml (`deleteDataByAlias` yaml:519-547, `deleteData` yaml:756-784, `deleteDataInGroup` yaml:1285-1313 — each currently 204/401/404/405/500) — unsetting a path that matches no endpoint mapping returns `engine.ErrPathNotFound` (`UnsetServerProperty`, internal/engine/serverdata.go:241-244) which `writeError` answers as `400 "Endpoint not found"` (internal/appengine/http.go:679-680), the same 400 the PUT/POST twins already document via `BadRequest`; internal/appengine/writeerror_test.go pins the mapping. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 200/401/500) — a `from_token` that does not parse as a device id is refused with `ErrValidation: invalid cursor` (internal/appengine/service.go:150-155) which `writeError` maps to 422 with the flat `{"errors":{"detail":"invalid cursor"}}` body, i.e. the `ErrorDetail` shape and not `FieldErrorsDetail`, so `$ref` the same response the group twin uses for the identical parameter (`ValidationErrors`, yaml:934-935) only after fixing its example, or add a dedicated one. Extend the `from_token` parameter description (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:556) and `InternalServerError` shows `detail: Internal Server Error` (yaml:638), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37, golden fixture pkg/astarteapi/testdata/error_bad_request.json) and `"Internal server error"` (envelope.go:58, testdata/error_internal.json), both measured on upstream 1.2.0 (envelope.go:33-36, 55-57). The other five detail examples on this file (`Not Found`, `Device not found`, `Unauthorized`, `Forbidden`, `Too Many Requests`) are already correct — change only those two. The identical wrong pair is in the other four specs (appengine 1732/1792, housekeeping 362/385/477, realm_management 1265/1348, native 466); leave those for a run of their surface. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
- [x] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pairing_api.yaml, which shows `version: "0.1.0-test"` (yaml:252) — a value no build can produce. The field is filled from `pairing.Config.Version` (internal/pairing/service.go:302), wired from `main.version` (cmd/astrate/main.go:151), whose default is `"0.1.0-dev"` (main.go:53) and which is otherwise set only via `-ldflags "-X main.version=..."`. Use `"0.1.0-dev"`, matching the four compat-version examples docs-sync-native-version-value already corrected to that (docs/api/astrate_native_api.yaml:203, 223, 243, 263). Do NOT document the service's own fallback `DefaultVersion = "0.1.0-astrate"` (service.go:86-88): it is unreachable from `main`, which always passes a non-empty value. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
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

$ ls /root/astrate-mule/docs/api/ && echo "---" && ls /root/astrate-mule/.mule/ && echo "---" && ls /root/astrate-mule/internal/ | head -50
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
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

# Todos
[•] Pick one surface and diff code routes vs docs/api paths
[ ] Spot-check 3-4 endpoints' status codes / response fields against handlers
[ ] Write up to 5 task lines to .mule/todo.md / .mule/for-giulio.md

$ cd /root/astrate-mule && rg -n '^\s{2}/' docs/api/astarte_housekeeping_api.yaml
26:  /housekeeping/v1/realms:
108:  /housekeeping/v1/realms/{realm}:

$ cd /root/astrate-mule && rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/housekeeping/ 2>/dev/null; echo "=== files ==="; ls internal/housekeeping/
=== files ===
async_operation_test.go
http.go
http_test.go
service.go
service_test.go

$ cd /root/astrate-mule && rg -n 'HandleFunc|Methods\(|Path\(' internal/housekeeping/http.go | head -40
(no output)

$ cd /root/astrate-mule && ls internal/httpx/ && rg -rn 'housekeeping' --glob '*.go' internal/ cmd/ -l | head -20
cors.go
cors_test.go
notfound.go
notfound_test.go
internal/config/config.go
internal/httpx/notfound_test.go
internal/httpx/notfound.go
internal/auth/middleware_test.go
internal/swagger/swagger_test.go
internal/broker/broker.go
internal/housekeeping/http.go
internal/broker/identity.go
internal/realm/http_test.go
internal/realm/service.go
internal/observability/compat_version_test.go
internal/housekeeping/service.go
internal/housekeeping/service_test.go
internal/housekeeping/async_operation_test.go
internal/housekeeping/http_test.go
internal/pairing/service.go
cmd/astrate/main.go

$ cd /root/astrate-mule && rg -n 'GET|POST|DELETE|http\.' internal/housekeeping/http.go | head -60
23:	require func(http.Handler) http.Handler
35:func (a *API) Mount(mux *http.ServeMux) {
36:	h := func(f http.HandlerFunc) http.Handler { return a.require(f) }
37:	mux.Handle("GET /housekeeping/v1/realms", h(a.listRealms))
38:	mux.Handle("POST /housekeeping/v1/realms", h(a.createRealm))
39:	mux.Handle("GET /housekeeping/v1/realms/{realm}", h(a.getRealm))
41:	mux.Handle("DELETE /housekeeping/v1/realms/{realm}", h(a.deleteRealm))
53:func (a *API) listRealms(w http.ResponseWriter, r *http.Request) {
59:	_ = astarteapi.WriteData(w, http.StatusOK, names)
62:func (a *API) createRealm(w http.ResponseWriter, r *http.Request) {
73:	_ = astarteapi.WriteData(w, http.StatusCreated, viewBody(rv))
76:func (a *API) getRealm(w http.ResponseWriter, r *http.Request) {
82:	_ = astarteapi.WriteData(w, http.StatusOK, viewBody(rv))
85:func (a *API) deleteRealm(w http.ResponseWriter, r *http.Request) {
90:	w.WriteHeader(http.StatusNoContent)
163:func (a *API) patchRealm(w http.ResponseWriter, r *http.Request) {
176:			_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
188:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
193:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
198:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
229:	w.WriteHeader(http.StatusOK)
233:func (a *API) writeError(w http.ResponseWriter, err error) {
236:		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Realm deletion disabled")
238:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
241:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
243:		_ = astarteapi.WriteError(w, http.StatusConflict, "Realm already exists")

→ Read internal/housekeeping/http.go [limit=100]
→ Read docs/api/astarte_housekeeping_api.yaml

$ cd /root/astrate-mule && rg -n 'hk-|ErrorDetail|housekeeping' .mule/todo.md | head -20
35:    ssh legion 'cd ~/astrate/bench && ./scripts/run-tier.sh small astrate -base-url ... -housekeeping-key ...'
81:- [x] docs-sync-rm-triggers-422-nested-envelope [auto]: document the nested changeset envelope on `POST /realmmanagement/v1/{realm}/triggers` in docs/api/astarte_realm_management_api.yaml — when trigger compile fails, `createTrigger` answers 422 with the upstream nested body `{"errors":{"action":{...},"simple_triggers":[...]}}` (internal/realm/http.go:252-258, triggerErrorBody at http.go:268-285, issues #63/#70), but the spec's createTrigger 422 (yaml:344-345) refs only the flat `ValidationError` component (`{"errors":{"detail":...}}`) — which the handler also emits for plain ErrValidation (http.go:259-260, e.g. "trigger requires a name", service.go:429). Model the response as oneOf across the flat ErrorDetail and the nested envelope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
83:- [!] hk-zero-reglimit-create-asymmetry: POST /realms persists device_registration_limit 0 as a literal 0 (service.go:143 only rejects < 0) while upstream Astarte v1.2.0 folds 0 to nil at create (`insert_realm` does `if device_limit == 0, do: nil`, so GET returns null); PATCH 0 already matches upstream (sets 0, http.go:208-214). Fold 0 → nil in CreateRealm for parity, same pattern as the retention fix, and pin with a container-free test. — BLOCKED: opencode exited 1
101:- [ ] race-check-store: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/store/... ./internal/housekeeping/... ./migrations/...`. Report any failure to .mule/for-giulio.md with the full race report. Split out of the former single `race-check` line, which timed out running the whole tree at once. [legion] [readonly]
125:- [x] housekeeping-tests [auto]: add a container-free unit test file for `internal/housekeeping` (service.go, 262 lines, zero test coverage). Cover two rules: (1) `CreateRealm` validation — blank name, blank JWT key, negative regLimit, negative retention, and default retention injection when `WithDefaultDatastreamMaximumStorageRetention` is set and the caller passes nil; (2) `DeleteRealm` gating — `ErrDeletionDisabled` when `WithRealmDeletionDisabled(true)`, `ErrConnectedDevicesPresent` when `DeviceStats` reports connected > 0. Use httptest-style fakes for the store and sealer (no Docker).
149:- [x] docs-sync-hk-patch-endpoint [auto]: add the undocumented `PATCH /housekeeping/v1/realms/{realm}` route to docs/api/astarte_housekeeping_api.yaml — it exists in code (internal/housekeeping/http.go:40, handler patchRealm at http.go:163-229, returns 200 on success) but the spec jumps straight from GET to DELETE. Document the request body (`RealmPatch` — same three optional fields as `RealmCreate`: `jwt_public_key_pem`, `device_registration_limit`, `datastream_maximum_storage_retention`, all nullable), the success response, and the error responses: 400 (bad envelope), 401, 403, 404, 405 (deletion disabled), 422 (invalid update parameters for unknown fields, blank jwt_public_key_pem, negative limit/retention, connected devices present), 500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
150:- [x] docs-sync-hk-retention-field [auto]: add the `datastream_maximum_storage_retention` field to the `Realm` and `RealmCreate` schemas in docs/api/astarte_housekeeping_api.yaml — the field is on the wire (realmBody struct, internal/housekeeping/http.go:50, returned by viewBody at http.go:99, accepted in create at http.go:68 and patch at http.go:216) but absent from both schemas (yaml:200-231). Type: integer, format int64, nullable, description matching the device_registration_limit style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
157:- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
178:- [x] docs-sync-hk-delete-gating-responses [auto]: add the missing `405` and `422` responses to `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — `writeError` answers 405 "Realm deletion disabled" for `ErrDeletionDisabled` and 422 `{"errors":{"error_name":["connected_devices_present"]}}` for `ErrConnectedDevicesPresent` (internal/housekeeping/http.go:233-249), both reachable from `Service.DeleteRealm` (service.go:248-266), but the spec lists only 204/401/403/404/500 (yaml:198-208). Its sole `405` reference sits on PATCH (yaml:182-183), which can never return it (`UpdateRealm`, service.go:196-216, has no deletion-gating path) — re-point 405 to DELETE (reusing the MethodNotAllowed component) and add 422 with the FieldErrors shape and a connected_devices_present example. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
179:- [x] docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_housekeeping_api.yaml (yaml:417) — it shows `detail: "housekeeping: validation failed: realm name must not be empty"`, but `validationDetail` strips the `housekeeping: validation failed: ` prefix on the wire (internal/housekeeping/http.go:251-258) and the real blank-name message is `"realm_name can't be blank"` (service.go:138). Change the example to `detail: "realm_name can't be blank"`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
184:- [x] hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto] — already resolved and re-queued by the refill: commit 015a706 chose the upstream-parity behaviour (0 ≡ unset, not invalid, not a literal 0) and pinned it — create folds 0 → nil with the container-free TestCreateRealmZeroRetentionUnsets (internal/housekeeping/service_test.go:186) and PATCH 0 → ClearRetention in the integration suite (internal/housekeeping/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).
196:- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
217:- [x] docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` (yaml:231-232) in docs/api/astarte_realm_management_api.yaml as a oneOf — both handlers go through `writeInterfaceError`, which answers three distinct 422 bodies: the flat ErrorDetail (`ErrValidation` via `validationDetail`, http.go:392-393), the nested violations changeset envelope (`writeViolations`, http.go:358-359 + 434-533, e.g. `{"errors":{"description":["should be at most 1000 character(s)"]}}` and the aligned full-length `mappings` array, http_test.go:480-501), and the named FieldErrors envelope for `ErrMaximumDatabaseRetentionExceeded` (`{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, http.go:355-357, http_test.go:282-320); today only the flat ValidationError is referenced. Mirror the createTrigger oneOf pattern (yaml:344-368). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
218:- [x] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "realm: validation failed: interface definition is invalid"`, but `validationDetail` strips the `realm: validation failed: ` prefix on the wire (http.go:562-568) and the parser's real messages are `invalid interface: ...` (pkg/interfaceschema/parse.go:38, wrapped at service.go:163); the example should show the stripped, real message. Same class as the already-fixed docs-sync-hk-validation-example (housekeeping http.go:251-258). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
219:- [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
226:- [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
227:- [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
228:- [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
236:- [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml — the string `"403"` appears nowhere in the file and `components` defines no `Forbidden` response, yet every route is wrapped by `a.require = mw.RequireRealm(auth.ClaimAppEngine)` (internal/appengine/http.go:35), which answers `403 Forbidden` for a verified a_aea JWT whose grants do not authorize the method+path (`astarteapi.WriteForbidden`, internal/auth/middleware.go:90; `DetailForbidden`, pkg/astarteapi/envelope.go:41; pinned today by internal/appengine/http_test.go:184-185). Add a `Forbidden` response component copied from the realm-management one (docs/api/astarte_realm_management_api.yaml:1277-1285 — `ErrorDetail` schema, example `errors.detail: Forbidden`) and `$ref` it from all 32 operations, which is exactly the parity docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go:123-170) pins; add the same docs-half test for the AppEngine spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
240:- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 200/401/500) — a `from_token` that does not parse as a device id is refused with `ErrValidation: invalid cursor` (internal/appengine/service.go:150-155) which `writeError` maps to 422 with the flat `{"errors":{"detail":"invalid cursor"}}` body, i.e. the `ErrorDetail` shape and not `FieldErrorsDetail`, so `$ref` the same response the group twin uses for the identical parameter (`ValidationErrors`, yaml:934-935) only after fixing its example, or add a dedicated one. Extend the `from_token` parameter description (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it

→ Read internal/housekeeping/http.go [offset=100]

$ cd /root/astrate-mule && rg -n '^- \[ \]' .mule/todo.md | rg -i 'hk-|housekeeping|errordetail|docs-sync'
101:- [ ] race-check-store: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/store/... ./internal/housekeeping/... ./migrations/...`. Report any failure to .mule/for-giulio.md with the full race report. Split out of the former single `race-check` line, which timed out running the whole tree at once. [legion] [readonly]
246:- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `base64.RawURLEncoding.Strict()` (pkg/deviceid/deviceid.go:37), so of the 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with `ErrInvalid` while `..._lA` is accepted, and the pairing spec's `hw_id` pattern `^[A-Za-z0-9_-]{22}$` (added by docs-sync-pairing-deviceid-base64url) deliberately does not encode the rule because a pattern is a client-side contract and upstream's leniency is unmeasured. The docstring justifies strictness as byte-identical round-tripping, which is a reason, not a parity claim, and docs/COMPATIBILITY.md records no entry for it. Measure upstream (Elixir `Base.url_decode64!(padding: false)`, which does not check those bits) on the Legion Go, then either record the divergence in docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.

$ cd /root/astrate-mule && rg -n 'ErrorDetail' internal/swagger/swagger_test.go | head -30; echo "=== test names ==="; rg -n '^func Test' internal/swagger/swagger_test.go
460:// TestPairingErrorDetailExamples pins the error-detail examples in the pairing
470:func TestPairingErrorDetailExamples(t *testing.T) {
509:// TestRealmManagementErrorDetailExamples pins the error-detail examples in the
512:// from TestPairingErrorDetailExamples. The BadRequest and InternalServerError
521:func TestRealmManagementErrorDetailExamples(t *testing.T) {
561:// TestAppEngineErrorDetailExamples pins the error-detail examples in the app
564:// management specs get from TestPairingErrorDetailExamples and
565:// TestRealmManagementErrorDetailExamples. The BadRequest and InternalServerError
573:func TestAppEngineErrorDetailExamples(t *testing.T) {
613:// TestNativeErrorDetailExamples pins the error-detail examples in the native
616:// app engine specs get from TestPairingErrorDetailExamples,
617:// TestRealmManagementErrorDetailExamples and TestAppEngineErrorDetailExamples.
628:func TestNativeErrorDetailExamples(t *testing.T) {
665:// TestHousekeepingErrorDetailExamples pins the error-detail examples in the
669:// TestPairingErrorDetailExamples, TestRealmManagementErrorDetailExamples,
670:// TestAppEngineErrorDetailExamples and TestNativeErrorDetailExamples. This is
685:// ErrorDetail example is deliberately not pinned either — it is not a
687:func TestHousekeepingErrorDetailExamples(t *testing.T) {
=== test names ===
23:func TestMount(t *testing.T) {
82:func TestMountSubPanicsOnBrokenFS(t *testing.T) {
95:func TestSpecs(t *testing.T) {
128:func TestRealmManagement403(t *testing.T) {
139:func TestAppEngine403(t *testing.T) {
203:func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
246:func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
321:func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {
381:func TestAppEngineAddGroupDevice422Documented(t *testing.T) {
425:func TestAppEngineDataDelete400Documented(t *testing.T) {
470:func TestPairingErrorDetailExamples(t *testing.T) {
521:func TestRealmManagementErrorDetailExamples(t *testing.T) {
573:func TestAppEngineErrorDetailExamples(t *testing.T) {
628:func TestNativeErrorDetailExamples(t *testing.T) {
687:func TestHousekeepingErrorDetailExamples(t *testing.T) {
743:func TestRealmManagement422ValidationDetails(t *testing.T) {
820:func TestPairingDeviceIDEncodingDocumented(t *testing.T) {
902:func TestRealmManagementDeviceIDEncodingDocumented(t *testing.T) {
966:func TestPairingUnregisterDeviceSemanticsDocumented(t *testing.T) {

$ cd /root/astrate-mule && rg -n 'func |ErrValidation|ErrDeletionDisabled|ErrConnectedDevices' internal/housekeeping/service.go
19:// ErrValidation wraps a well-formed request that violates a realm-creation
21:var ErrValidation = errors.New("housekeeping: validation failed")
27:	// ErrDeletionDisabled reports DELETE against a deployment whose
29:	ErrDeletionDisabled = errors.New("housekeeping: realm deletion disabled")
30:	// ErrConnectedDevicesPresent reports DELETE against a realm that has
32:	ErrConnectedDevicesPresent = errors.New("housekeeping: connected devices present")
83:	// (#75): when set, DeleteRealm answers ErrDeletionDisabled without
91:func NewService(st realmStore, sealer sealer, reloader Reloader, log *slog.Logger) *Service {
103:func (s *Service) WithDefaultDatastreamMaximumStorageRetention(defaultRetention *int64) *Service {
111:func (s *Service) WithRealmDeletionDisabled(disabled bool) *Service {
119:func (s *Service) notifyBrokerReload(ctx context.Context, realm string) {
131:// A blank/invalid name or missing JWT key yields ErrValidation; a duplicate
136:func (s *Service) CreateRealm(ctx context.Context, name, jwtPublicKeyPEM string, regLimit *int32, retention *int64) (*RealmView, error) {
138:		return nil, fmt.Errorf("%w: realm_name can't be blank", ErrValidation)
141:		return nil, fmt.Errorf("%w: jwt_public_key_pem can't be blank", ErrValidation)
144:		return nil, fmt.Errorf("%w: device_registration_limit must be non-negative", ErrValidation)
147:		return nil, fmt.Errorf("%w: datastream_maximum_storage_retention must be non-negative", ErrValidation)
180:			return nil, fmt.Errorf("%w: realm_name is invalid", ErrValidation)
204:func (s *Service) UpdateRealm(ctx context.Context, name string, p RealmUpdate) error {
206:		return fmt.Errorf("%w: jwt_public_key_pem can't be blank", ErrValidation)
209:		return fmt.Errorf("%w: device_registration_limit must be non-negative", ErrValidation)
212:		return fmt.Errorf("%w: datastream_maximum_storage_retention must be non-negative", ErrValidation)
228:func (s *Service) GetRealm(ctx context.Context, name string) (*RealmView, error) {
238:func (s *Service) ListRealms(ctx context.Context) ([]string, error) {
253:// Gated upstream-style (#75): ErrDeletionDisabled when the deployment flag is
254:// set (405), ErrConnectedDevicesPresent when the realm has devices still
256:func (s *Service) DeleteRealm(ctx context.Context, name string) error {
258:		return ErrDeletionDisabled
267:		return ErrConnectedDevicesPresent
277:func view(r *store.Realm) *RealmView {

→ Read internal/auth/middleware.go
→ Read internal/housekeeping/service.go [offset=126, limit=150]

$ cd /root/astrate-mule && rg -n 'ErrInvalidRealmName|validRealmName|realmName' internal/store/realms.go | head -20
69:		return nil, fmt.Errorf("%w: %q", ErrInvalidRealmName, nr.Name)

→ Read internal/swagger/swagger_test.go [offset=660, limit=90]

$ cd /root/astrate-mule && rg -n 'ErrorDetail|schema-level|yaml:361|:361' .mule/todo.md .mule/for-giulio.md | rg -i 'hk|housekeep|shape' | head; echo "=== realms.go ==="; sed -n '50,85p' internal/store/realms.go
.mule/todo.md:217:- [x] docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` (yaml:231-232) in docs/api/astarte_realm_management_api.yaml as a oneOf — both handlers go through `writeInterfaceError`, which answers three distinct 422 bodies: the flat ErrorDetail (`ErrValidation` via `validationDetail`, http.go:392-393), the nested violations changeset envelope (`writeViolations`, http.go:358-359 + 434-533, e.g. `{"errors":{"description":["should be at most 1000 character(s)"]}}` and the aligned full-length `mappings` array, http_test.go:480-501), and the named FieldErrors envelope for `ErrMaximumDatabaseRetentionExceeded` (`{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, http.go:355-357, http_test.go:282-320); today only the flat ValidationError is referenced. Mirror the createTrigger oneOf pattern (yaml:344-368). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:219:- [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:240:- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 200/401/500) — a `from_token` that does not parse as a device id is refused with `ErrValidation: invalid cursor` (internal/appengine/service.go:150-155) which `writeError` maps to 422 with the flat `{"errors":{"detail":"invalid cursor"}}` body, i.e. the `ErrorDetail` shape and not `FieldErrorsDetail`, so `$ref` the same response the group twin uses for the identical parameter (`ValidationErrors`, yaml:934-935) only after fixing its example, or add a dedicated one. Extend the `from_token` parameter description (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
.mule/todo.md:261:- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDetail` schema-level example (yaml:361-362) and the `BadRequest` response (yaml:385) show `detail: Bad Request`, and `InternalServerError` (yaml:477) shows `detail: Internal Server Error`, but the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58, via WriteBadRequest/WriteInternalServerError, pinned by pkg/astarteapi/testdata/error_bad_request.json and error_internal.json). While you are in `ErrorDetail`: its `example` sits at the schema level and shows a bare `{detail: ...}` where the envelope the wire sends is `{"errors":{"detail":...}}` — move it under `errors` (or drop it, as the realm-management and pairing `ErrorDetail` schemas carry no example) so a generated client is not shown the wrong shape. Add `TestHousekeepingErrorDetailExamples` beside the pairing and realm-management halves in internal/swagger/swagger_test.go. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
.mule/todo.md:263:- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
=== realms.go ===

// CreateRealm inserts the realm row together with its CA material in one
// transaction (docs/ROADMAP.md §3.1 file 2.8) and returns the stored realm.
func (s *Store) CreateRealm(ctx context.Context, nr NewRealm) (*Realm, error) {
	keys := nr.JWTPublicKeysPEM
	if keys == nil {
		keys = []string{}
	}
	row := s.pool.QueryRow(ctx, `
		INSERT INTO realms (name, jwt_public_keys, ca_certificate, ca_private_key, device_registration_limit, datastream_maximum_storage_retention)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING `+realmColumns,
		nr.Name, keys, nr.CACertificatePEM, nr.CAPrivateKeySealed, nr.DeviceRegistrationLimit, nr.DatastreamMaximumStorageRetention)

	realm, err := scanRealm(row)
	switch pgErrCode(err) {
	case pgCodeUniqueViolation:
		return nil, fmt.Errorf("%w: realm %q", ErrAlreadyExists, nr.Name)
	case pgCodeCheckViolation:
		return nil, fmt.Errorf("%w: %q", ErrInvalidRealmName, nr.Name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: creating realm %q: %w", nr.Name, err)
	}
	return realm, nil
}

// GetRealm fetches a realm by ID.
func (s *Store) GetRealm(ctx context.Context, id int16) (*Realm, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+realmColumns+` FROM realms WHERE id = $1`, id)
	realm, err := scanRealm(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: realm id %d", ErrNotFound, id)
	}
	if err != nil {
		return nil, fmt.Errorf("store: getting realm id %d: %w", id, err)

$ cd /root/astrate-mule && rg -rn 'realms' migrations/ -l | head && echo "=== check constraints ===" && rg -n 'CHECK|check' migrations/*realm* 2>/dev/null | head -20
migrations/000012_user_blocks.up.sql
migrations/000009_flows.up.sql
migrations/000006_policies.up.sql
migrations/000011_realm_retention.down.sql
migrations/000011_realm_retention.up.sql
migrations/000008_pipelines.up.sql
migrations/000002_metadata.down.sql
migrations/000002_metadata.up.sql
=== check constraints ===

$ cd /root/astrate-mule && rg -n 'CREATE TABLE realms' -A 20 migrations/000002_metadata.up.sql
5:CREATE TABLE realms (
6-    id               smallint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
7-    name             text NOT NULL UNIQUE CHECK (name ~ '^[a-z][a-z0-9]*$'),
8-    jwt_public_keys  jsonb NOT NULL DEFAULT '[]',   -- array of PEM strings (RSA/EC)
9-    ca_certificate   text NOT NULL,                  -- realm CA cert, PEM
10-    ca_private_key   bytea NOT NULL,                 -- encrypted at rest (AES-256-GCM,
11-                                                     -- key from config/KMS env var)
12-    device_registration_limit integer,
13-    created_at       timestamptz NOT NULL DEFAULT now()
14-);
15-
16--- Interfaces (Realm Management domain). The raw JSON is the source of truth;
17--- generated columns lift the routing-critical fields out for indexing.
18-CREATE TABLE interfaces (
19-    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
20-    realm_id      smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
21-    definition    jsonb NOT NULL,
22-    name          text     GENERATED ALWAYS AS (definition->>'interface_name') STORED,
23-    major_version integer  GENERATED ALWAYS AS ((definition->>'version_major')::int) STORED,
24-    minor_version integer  GENERATED ALWAYS AS ((definition->>'version_minor')::int) STORED,
25-    type          text     GENERATED ALWAYS AS (definition->>'type') STORED,          -- datastream|properties

$ cd /root/astrate-mule && rg -o '\bASTRATE_[A-Z_]+' -N internal/ | sort -u -t: -k2 | head -60
internal/config/config.go:ASTRATE_DATABASE_DSN
internal/config/config.go:ASTRATE_ENGINE_SHARDS
internal/flow/blocks/info.go:ASTRATE_FLOW_CONFIG
internal/config/config.go:ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
internal/config/config.go:ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
internal/config/config.go:ASTRATE_HTTP_ADDR
internal/config/config.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
internal/config/config.go:ASTRATE_HTTP_TLS_CERT_FILE
internal/config/config.go:ASTRATE_HTTP_TLS_KEY_FILE
internal/config/config.go:ASTRATE_LOG_FORMAT
internal/config/config.go:ASTRATE_LOG_LEVEL
internal/store/crypto.go:ASTRATE_MASTER_KEY
internal/store/crypto.go:ASTRATE_MASTER_KEY_FILE
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
internal/store/store_test.go:ASTRATE_TEST_DSN

$ mkdir -p /tmp/opencode/probe && cd /tmp/opencode/probe && cat > go.mod <<'EOF'
module probe

go 1.24

require github.com/astrate-platform/astrate v0.0.0
replace github.com/astrate-platform/astrate => /root/astrate-mule
EOF
cat > main.go <<'EOF'
package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

func main() {
	for _, body := range []string{
		`{"data":{"device_registration_limit":"100"}}`,
		`{"data":{"jwt_public_key_pem":123}}`,
		`{"data":{"unknown_field":1}}`,
		`{"data":{}}`,
		`{}`,
		`{"data":null}`,
		`not json`,
	} {
		var fields map[string]json.RawMessage
		err := astarteapi.DecodeData(bytes.NewReader([]byte(body)), 64<<10, &fields)
		fmt.Printf("map  %-45s -> err=%v keys=%v\n", body, err, keysOf(fields))
	}
}

func keysOf(m map[string]json.RawMessage) []string {
	out := []string{}
	for k := range m {
		out = append(out, k)
	}
	return out
}
EOF
GOFLAGS=-mod=mod go run . 2>&1 | tail -20
map  {"data":{"device_registration_limit":"100"}}  -> err=<nil> keys=[device_registration_limit]
map  {"data":{"jwt_public_key_pem":123}}           -> err=<nil> keys=[jwt_public_key_pem]
map  {"data":{"unknown_field":1}}                  -> err=<nil> keys=[unknown_field]
map  {"data":{}}                                   -> err=<nil> keys=[]
map  {}                                            -> err=astarteapi: missing "data" key in request body keys=[]
map  {"data":null}                                 -> err=astarteapi: missing "data" key in request body keys=[]
map  not json                                      -> err=astarteapi: decoding request envelope: invalid character 'o' in literal null (expecting 'u') keys=[]

$ cd /root/astrate-mule && rg -n 'ASTRATE_SECURITY_MASTER_KEY_FILE|ASTRATE_MASTER_KEY_FILE|ASTRATE_MASTER_KEY\b' docs/site/*.md internal/ -r '' | head -20
docs/site/operations.md:23:band, never in the config body: set `` (64 hex chars or
docs/site/operations.md:24:base64 of 32 bytes), ``, or `security.master_key_file`.
docs/site/operations.md:62:- A strong, secret `[_FILE]`.
docs/site/deployment.md:55:- `` (64 hex chars or base64 of 32 bytes)
docs/site/deployment.md:56:- ``
docs/site/configuration-reference.md:95:| `master_key_file` | string | `""` | `` | No | File holding the AES-256 master key that seals realm CA private keys. When empty, falls back to `` (64 hex chars) or ``. |
docs/site/configuration-reference.md:155:This sets `ASTRATE_MQTT_INSECURE_DEV_MODE=true`, `` (throwaway), and auto-provisions a `test` realm. See [Deployment](deployment.md) for details.
internal/config/config.go:135://  / .
internal/config/config.go:248:	str("", &cfg.Security.MasterKeyFile)
internal/config/config.example.toml:87:# empty, the store reads  or .
internal/store/crypto.go:24:	EnvMasterKey = ""
internal/store/crypto.go:27:	EnvMasterKeyFile = ""

$ cd /tmp/opencode/probe && cat > main.go <<'EOF'
package main

import (
	"bytes"
	"fmt"

	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

type optString struct {
	present bool
	val     string
}

func (o *optString) UnmarshalJSON(b []byte) error {
	o.present = true
	return jsonUnmarshal(b, &o.val)
}

type optInt32 struct {
	present bool
	null    bool
	val     int32
}

func (o *optInt32) UnmarshalJSON(b []byte) error {
	o.present = true
	if bytes.Equal(b, []byte("null")) {
		o.null = true
		return nil
	}
	return jsonUnmarshal(b, &o.val)
}

type patchBody struct {
	JWTPublicKeyPEM     optString `json:"jwt_public_key_pem"`
	DeviceRegistrationL optInt32  `json:"device_registration_limit"`
}

func main() {
	for _, body := range []string{
		`{"data":{"device_registration_limit":"100"}}`,
		`{"data":{"jwt_public_key_pem":123}}`,
		`{"data":{"jwt_public_key_pem":""}}`,
		`{"data":{"device_registration_limit":1.5}}`,
		`{"data":[]}`,
	} {
		var pb patchBody
		err := astarteapi.DecodeData(bytes.NewReader([]byte(body)), 64<<10, &pb)
		fmt.Printf("%-48s -> err=%v pem.present=%v pem=%q limit=%d\n", body, err, pb.JWTPublicKeyPEM.present, pb.JWTPublicKeyPEM.val, pb.DeviceRegistrationL.val)
	}
}
EOF
cat > helper.go <<'EOF'
package main

import "encoding/json"

func jsonUnmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }
EOF
GOFLAGS=-mod=mod go run . 2>&1 | tail -10
{"data":{"device_registration_limit":"100"}}     -> err=astarteapi: decoding request data: json: cannot unmarshal string into Go struct field patchBody.device_registration_limit of type int32 pem.present=false pem="" limit=0
{"data":{"jwt_public_key_pem":123}}              -> err=astarteapi: decoding request data: json: cannot unmarshal number into Go struct field patchBody.jwt_public_key_pem of type string pem.present=true pem="" limit=0
{"data":{"jwt_public_key_pem":""}}               -> err=<nil> pem.present=true pem="" limit=0
{"data":{"device_registration_limit":1.5}}       -> err=astarteapi: decoding request data: json: cannot unmarshal number 1.5 into Go struct field patchBody.device_registration_limit of type int32 pem.present=false pem="" limit=0
{"data":[]}                                      -> err=astarteapi: decoding request data: json: cannot unmarshal array into Go value of type main.patchBody pem.present=false pem="" limit=0

$ cd /root/astrate-mule && rg -n 'JWTPublicKeysPEM' internal/store/realms.go internal/realm/service.go | head -20; echo "=== GetRealmByName ==="; rg -n 'func \(s \*Store\) GetRealmByName' -A 12 internal/store/realms.go
internal/store/realms.go:19:	JWTPublicKeysPEM                  []string
internal/store/realms.go:31:	JWTPublicKeysPEM                  []string
internal/store/realms.go:54:	keys := nr.JWTPublicKeysPEM
internal/store/realms.go:244:	if err := row.Scan(&r.ID, &r.Name, &r.JWTPublicKeysPEM, &r.CACertificatePEM,
internal/realm/service.go:627:	return joinPEM(r.JWTPublicKeysPEM), nil
=== GetRealmByName ===
91:func (s *Store) GetRealmByName(ctx context.Context, name string) (*Realm, error) {
92-	row := s.pool.QueryRow(ctx, `SELECT `+realmColumns+` FROM realms WHERE name = $1`, name)
93-	realm, err := scanRealm(row)
94-	if errors.Is(err, pgx.ErrNoRows) {
95-		return nil, fmt.Errorf("%w: realm %q", ErrNotFound, name)
96-	}
97-	if err != nil {
98-		return nil, fmt.Errorf("store: getting realm %q: %w", name, err)
99-	}
100-	return realm, nil
101-}
102-
103-// ListRealms returns every realm ordered by name.

→ Read internal/housekeeping/service.go [offset=276]

$ cd /root/astrate-mule && rg -n 'jwt_public_keys' internal/store/*.go internal/realm/service.go migrations/*.sql | head -20
migrations/000002_metadata.up.sql:8:    jwt_public_keys  jsonb NOT NULL DEFAULT '[]',   -- array of PEM strings (RSA/EC)
internal/store/realms.go:38:const realmColumns = `id, name, jwt_public_keys, ca_certificate, ca_private_key, device_registration_limit, datastream_maximum_storage_retention, created_at`
internal/store/realms.go:59:		INSERT INTO realms (name, jwt_public_keys, ca_certificate, ca_private_key, device_registration_limit, datastream_maximum_storage_retention)
internal/store/realms.go:131:		`UPDATE realms SET jwt_public_keys = $2 WHERE name = $1`, name, keysPEM)
internal/store/realms.go:165:		add("jwt_public_keys", []string{p.SetJWTPublicKeyPEM})

$ cd /root/astrate-mule && sed -n '120,180p' internal/store/realms.go
return nil, fmt.Errorf("store: listing realms: %w", err)
	}
	return realms, nil
}

// SetRealmJWTPublicKeys replaces the realm's JWT public key set (PEM strings).
func (s *Store) SetRealmJWTPublicKeys(ctx context.Context, name string, keysPEM []string) error {
	if keysPEM == nil {
		keysPEM = []string{}
	}
	tag, err := s.pool.Exec(ctx,
		`UPDATE realms SET jwt_public_keys = $2 WHERE name = $1`, name, keysPEM)
	if err != nil {
		return fmt.Errorf("store: updating realm %q JWT keys: %w", name, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: realm %q", ErrNotFound, name)
	}
	return nil
}

// RealmPatch carries optional realm updates; each Patch* flag gates its
// value field, and Clear* flags write NULL (Clear beats Set).
type RealmPatch struct {
	PatchJWTPublicKeyPEM   bool
	SetJWTPublicKeyPEM     string
	PatchRegistrationLimit bool
	SetRegistrationLimit   int32
	ClearRegistrationLimit bool
	PatchRetention         bool
	SetRetention           int64 // seconds; SetRetention=0 with PatchRetention means clear too
	ClearRetention         bool
}

// UpdateRealm applies the patched fields to one realm in a single UPDATE built
// from only the touched columns. An unknown realm yields ErrNotFound wrapped
// like the other realm methods.
func (s *Store) UpdateRealm(ctx context.Context, name string, p RealmPatch) error {
	sets := []string{}
	args := []any{name} // $1 is the WHERE name
	add := func(col string, v any) {
		sets = append(sets, fmt.Sprintf("%s = $%d", col, len(args)+1))
		args = append(args, v)
	}
	if p.PatchJWTPublicKeyPEM {
		add("jwt_public_keys", []string{p.SetJWTPublicKeyPEM})
	}
	if p.ClearRegistrationLimit {
		add("device_registration_limit", nil)
	} else if p.PatchRegistrationLimit {
		add("device_registration_limit", p.SetRegistrationLimit)
	}
	if p.ClearRetention || (p.PatchRetention && p.SetRetention <= 0) {
		add("datastream_maximum_storage_retention", nil)
	} else if p.PatchRetention {
		add("datastream_maximum_storage_retention", p.SetRetention)
	}
	if len(sets) == 0 {
		return nil
	}

$ cd /root/astrate-mule && rg -n 'SetRealmJWTPublicKeys' --glob '*.go' . | grep -v _test | head; echo "=== SetAuthKey ==="; rg -n 'func \(s \*Service\) SetAuthKey' -A 30 internal/realm/service.go
./internal/store/realms.go:125:// SetRealmJWTPublicKeys replaces the realm's JWT public key set (PEM strings).
./internal/store/realms.go:126:func (s *Store) SetRealmJWTPublicKeys(ctx context.Context, name string, keysPEM []string) error {
./internal/realm/service.go:636:	if err := s.st.SetRealmJWTPublicKeys(ctx, realm, []string{keyPEM}); err != nil {
=== SetAuthKey ===
632:func (s *Service) SetAuthKey(ctx context.Context, realm, keyPEM string) error {
633-	if keyPEM == "" {
634-		return fmt.Errorf("%w: jwt_public_key_pem can't be blank", ErrValidation)
635-	}
636-	if err := s.st.SetRealmJWTPublicKeys(ctx, realm, []string{keyPEM}); err != nil {
637-		return err
638-	}
639-	return nil
640-}
641-
642-// joinPEM concatenates PEM blocks with a blank line, the form the verifier
643-// parses back into individual keys.
644-func joinPEM(keys []string) string {
645-	switch len(keys) {
646-	case 0:
647-		return ""
648-	case 1:
649-		return keys[0]
650-	default:
651-		out := keys[0]
652-		for _, k := range keys[1:] {
653-			out += "\n" + k
654-		}
655-		return out
656-	}
657-}
658-
659-// --- invalidation -----------------------------------------------------------
660-
661-// interfacesChanged emits the store NOTIFY and refreshes the engine snapshot
662-// after an interface mutation. Failures are logged, never fatal: the change

$ cd /root/astrate-mule && rg -n 'Detail[A-Za-z]* *=' pkg/astarteapi/envelope.go
37:	DetailBadRequest = "Bad request"
39:	DetailUnauthorized = "Unauthorized"
41:	DetailForbidden = "Forbidden"
43:	DetailNotFound = "Not Found"
46:	DetailDeviceNotFound = "Device not found"
51:	DetailRouteNotFound = "Not found"
54:	DetailPageNotFound = "Page not found"
58:	DetailInternalServerError = "Internal server error"

$ cd /root/astrate-mule && rg -n 'func Generate' -A 30 pkg/ca/*.go | head -50
rg: pkg/ca/*.go: IO error for operation on pkg/ca/*.go: No such file or directory (os error 2)

$ cd /root/astrate-mule && rg -n '"ca"|/ca"' internal/housekeeping/service.go && rg -rn 'func Generate' --glob '*.go' . | grep -v _test
15:	"github.com/astrate-platform/astrate/internal/pairing/ca"
./internal/pairing/ca/ca.go:n(realm string, lifetime time.Duration) (*CA, error) {

$ cd /root/astrate-mule && rg -n 'func Generate' -A 32 internal/pairing/ca/ca.go
71:func Generate(realm string, lifetime time.Duration) (*CA, error) {
72-	if lifetime == 0 {
73-		lifetime = DefaultCALifetime
74-	}
75-	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
76-	if err != nil {
77-		return nil, fmt.Errorf("ca: generating CA key: %w", err)
78-	}
79-	serial, err := randomSerial()
80-	if err != nil {
81-		return nil, err
82-	}
83-
84-	now := time.Now()
85-	template := &x509.Certificate{
86-		SerialNumber: serial,
87-		Subject: pkix.Name{
88-			CommonName:   "Astrate Realm " + realm + " CA",
89-			Organization: []string{"Astrate"},
90-		},
91-		NotBefore:             now.Add(-clockSkewBackdate),
92-		NotAfter:              now.Add(lifetime),
93-		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
94-		BasicConstraintsValid: true,
95-		IsCA:                  true,
96-		MaxPathLenZero:        true, // issues leaves only
97-	}
98-	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
99-	if err != nil {
100-		return nil, fmt.Errorf("ca: self-signing CA certificate: %w", err)
101-	}
102-	cert, err := x509.ParseCertificate(der)
103-	if err != nil {

$ cd /tmp/opencode/probe && cat > main.go <<'EOF'
package main

import (
	"fmt"

	"github.com/astrate-platform/astrate/internal/pairing/ca"
)

func main() {
	for _, n := range []string{"test", "TEST", "test_realm", "test-realm", "1test", "", "test realm", "tëst"} {
		c, err := ca.Generate(n, 0)
		if err != nil {
			fmt.Printf("%-12q -> ca.Generate err=%v\n", n, err)
			continue
		}
		fmt.Printf("%-12q -> ok cn=%q\n", n, c.Certificate().Subject.CommonName)
	}
}
EOF
GOFLAGS=-mod=mod go run . 2>&1 | tail -12
package probe
	main.go:6:2: use of internal package github.com/astrate-platform/astrate/internal/pairing/ca not allowed

$ cd /tmp/opencode/probe && cat > main.go <<'EOF'
package main

import (
	"fmt"
	"regexp"
)

// The CHECK constraint as written in migrations/000002_metadata.up.sql:7.
var re = regexp.MustCompile(`^[a-z][a-z0-9]*$`)

func main() {
	for _, n := range []string{"test", "TEST", "test_realm", "test-realm", "1test", "", "t", "test1", "tëst", "test realm"} {
		fmt.Printf("%-12q -> accepted=%v\n", n, re.MatchString(n))
	}
}
EOF
GOFLAGS=-mod=mod go run . 2>&1 | tail -12
"test"       -> accepted=true
"TEST"       -> accepted=false
"test_realm" -> accepted=false
"test-realm" -> accepted=false
"1test"      -> accepted=false
""           -> accepted=false
"t"          -> accepted=true
"test1"      -> accepted=true
"tëst"       -> accepted=false
"test realm" -> accepted=false

$ cd /root/astrate-mule && rg -n -i 'realm name|realm_name|lowercase' docs/site/*.md | head -15; echo "=== todo lines ==="; wc -l .mule/todo.md
docs/site/migration-from-astarte.md:115:3. **Deploy Astrate** with the same realm names.
docs/site/housekeeping-api.md:16:    "realm_name": "<name>",
docs/site/housekeeping-api.md:41:GET /housekeeping/v1/realms/<realm_name>
docs/site/housekeeping-api.md:47:DELETE /housekeeping/v1/realms/<realm_name>
docs/site/configuration-reference.md:104:| `name` | string | `""` | `ASTRATE_REALM_NAME` | No | Realm name. Empty disables auto-provisioning. |
docs/site/architecture.md:89:Realms become cheap: a `realms` row + per-realm CA + per-realm JWT public keys. A single-realm install is just a realm named e.g. `home` created at first boot.
docs/site/compatibility.md:98:Upstream's Phoenix renders `Bad request` and `Internal server error`; Astrate uses the measured lowercase forms rather than the reconstructed title-case ones.
docs/site/operations.md:31:`ASTRATE_REALM_NAME`.
docs/site/DESIGN.md:163:keys. A single-realm install is just a realm named e.g. `home` created at first boot via
=== todo lines ===
267 .mule/todo.md

→ Read docs/site/housekeeping-api.md

$ cd /root/astrate-mule && rg -n 'housekeeping|Housekeeping' internal/config/config.go | head -30
49:	Housekeeping HousekeepingConfig `toml:"housekeeping"`
109:// HousekeepingConfig carries the instance-admin JWT public keys (a_ha): either
111:type HousekeepingConfig struct {
119:	// RealmDeletionDisabled gates DELETE /housekeeping/v1/realms/{realm}
289:			return fmt.Errorf("config: housekeeping.default_datastream_maximum_storage_retention %q must be a non-negative integer", retention)
291:		cfg.Housekeeping.DefaultDatastreamMaximumStorageRetention = &n
300:			cfg.Housekeeping.RealmDeletionDisabled = true
302:			cfg.Housekeeping.RealmDeletionDisabled = false
304:			return fmt.Errorf("config: housekeeping.realm_deletion_disabled %q must be a boolean (1|true|TRUE|True|0|false)", v)
392:// HousekeepingKeys resolves the instance-admin JWT public keys, reading any
394:func (c *Config) HousekeepingKeys() ([]string, error) {
395:	keys := append([]string(nil), c.Housekeeping.JWTPublicKeys...)
396:	for _, f := range c.Housekeeping.JWTPublicKeyFiles {
399:			return nil, fmt.Errorf("config: reading housekeeping key %s: %w", f, err)

$ cd /root/astrate-mule && rg -n 'realm_deletion_disabled|default_datastream_maximum_storage_retention|jwt_public_key_files' docs/site/configuration-reference.md docs/site/*.md | head -10; echo "=== reference housekeeping block ==="; rg -n '^\| `housekeeping' -A 2 docs/site/configuration-reference.md | head -20
docs/site/operations.md:70:  `jwt_public_key_files`) — the instance-admin JWT (`a_ha`) public keys.
docs/site/configuration-reference.md:83:| `jwt_public_key_files` | list[string] | `[]` | — | No | File paths to PEM public keys. Both inline and file references are concatenated. |
docs/site/configuration-reference.md:140:jwt_public_key_files = ["/etc/astrate/housekeeping.pub"]
docs/site/housekeeping-api.md:71:housekeeping.jwt_public_key_files = ["/etc/astrate/admin_public.pem"]
docs/site/configuration-reference.md:83:| `jwt_public_key_files` | list[string] | `[]` | — | No | File paths to PEM public keys. Both inline and file references are concatenated. |
docs/site/configuration-reference.md:140:jwt_public_key_files = ["/etc/astrate/housekeeping.pub"]
=== reference housekeeping block ===

$ cd /root/astrate-mule && sed -n '74,100p' docs/site/configuration-reference.md
| `credentials_rate` | float | `5.0` | — | No | Token-bucket rate limit for credential requests (per IP, requests/sec). |
| `credentials_burst` | int | `10` | — | No | Token-bucket burst for credential requests. |
| `bcrypt_cost` | int | `10` | — | No | bcrypt cost for hashing credentials secrets. |

## `[housekeeping]` — Instance-admin keys

| Key | Type | Default | Env override | Required | Description |
|---|---|---|---|---|---|
| `jwt_public_keys` | list[string] | `[]` | — | No | PEM blocks inline for instance-admin JWT public keys (claim `a_ha`). |
| `jwt_public_key_files` | list[string] | `[]` | — | No | File paths to PEM public keys. Both inline and file references are concatenated. |

## `[storage]` — Retention

| Key | Type | Default | Env override | Required | Description |
|---|---|---|---|---|---|
| `retention` | duration | `""` (disabled) | — | No | Global drop-chunks retention. Empty/0 disables it; per-endpoint TTL still applies. Example: `"8760h"` (one year). |

## `[security]` — Master key

| Key | Type | Default | Env override | Required | Description |
|---|---|---|---|---|---|
| `master_key_file` | string | `""` | `ASTRATE_SECURITY_MASTER_KEY_FILE` | No | File holding the AES-256 master key that seals realm CA private keys. When empty, falls back to `ASTRATE_MASTER_KEY` (64 hex chars) or `ASTRATE_MASTER_KEY_FILE`. |

!!! warning "Master key"
    Losing the master key means re-issuing realm CAs. Devices re-pair automatically at their next credential rotation since their credentials secret still works. Keep a separate backup of the master key.

## `[realm]` — Auto-provision realm

$ cd /root/astrate-mule && sed -n '105,126p' internal/config/config.go; echo "=== env binding ==="; sed -n '278,308p' internal/config/config.go
CredentialsBurst  int      `toml:"credentials_burst"`
	BcryptCost        int      `toml:"bcrypt_cost"`
}

// HousekeepingConfig carries the instance-admin JWT public keys (a_ha): either
// inline PEM blocks or file references; both are concatenated.
type HousekeepingConfig struct {
	JWTPublicKeys     []string `toml:"jwt_public_keys"`
	JWTPublicKeyFiles []string `toml:"jwt_public_key_files"`
	// DefaultDatastreamMaximumStorageRetention is the realm-level datastream
	// storage ceiling (seconds) injected at realm creation when the caller
	// omits the field (#73). When unset (nil) no default is injected and
	// existing deployments behave exactly as before.
	DefaultDatastreamMaximumStorageRetention *int64 `toml:"default_datastream_maximum_storage_retention"`
	// RealmDeletionDisabled gates DELETE /housekeeping/v1/realms/{realm}
	// (#75): when true, deletion answers upstream's 405 "Realm deletion
	// disabled". Opt-in — unset (false) keeps the historical
	// always-delete behavior.
	RealmDeletionDisabled bool `toml:"realm_deletion_disabled"`
}

// StorageConfig holds runtime storage policy. Retention applies a global
=== env binding ===
	for _, env := range []string{
		"ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
		"HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
	} {
		if v, ok := os.LookupEnv(env); ok && v != "" {
			retention = v
		}
	}
	if retention != "" {
		n, err := strconv.ParseInt(retention, 10, 64)
		if err != nil || n < 0 {
			return fmt.Errorf("config: housekeeping.default_datastream_maximum_storage_retention %q must be a non-negative integer", retention)
		}
		cfg.Housekeeping.DefaultDatastreamMaximumStorageRetention = &n
	}

	// The realm-deletion gate (#75) is a fail-loud boolean: only the exact
	// spellings below are accepted, anything else refuses to load rather
	// than silently meaning "off". Absent (or empty) leaves the default.
	if v, ok := os.LookupEnv("ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"); ok {
		switch v {
		case "1", "true", "TRUE", "True":
			cfg.Housekeeping.RealmDeletionDisabled = true
		case "0", "false", "":
			cfg.Housekeeping.RealmDeletionDisabled = false
		default:
			return fmt.Errorf("config: housekeeping.realm_deletion_disabled %q must be a boolean (1|true|TRUE|True|0|false)", v)
		}
	}
	return nil
}

$ cd /root/astrate-mule && tail -40 .mule/for-giulio.md
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

$ cd /root/astrate-mule && tail -c 300 .mule/todo.md | cat -A | tail -3 | cut -c1-120
ody, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to t

$ cd /root/astrate-mule && tail -c 1 .mule/todo.md | xxd | tail -1; tail -c 1 .mule/for-giulio.md | xxd | tail -1
/bin/sh: 1: xxd: not found
/bin/sh: 1: xxd: not found

$ cd /root/astrate-mule && for f in .mule/todo.md .mule/for-giulio.md; do printf '%s: last byte = ' "$f"; tail -c 1 "$f" | od -c | head -1; done
.mule/todo.md: last byte = 0000000  \n
.mule/for-giulio.md: last byte = 0000000  \n

Both files end cleanly. Appending the four spec-fix proposals to `.mule/todo.md`:
$ cd /root/astrate-mule && cat >> .mule/todo.md <<'EOF'
- [ ] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362) — the second half docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail: Bad Request`, the reconstructed title-case form that `TestHousekeepingErrorDetailExamples` (internal/swagger/swagger_test.go:687) now pins against everywhere else on this spec, while the wire emits the frozen `astarteapi.DetailBadRequest` = `"Bad request"` (pkg/astarteapi/envelope.go:37, via `WriteBadRequest` on every 400 — internal/housekeeping/http.go:65, 166, 171, 183). Because the example sits on the generic `ErrorDetail` schema, which all eight error responses `$ref` (yaml:377-477), it is the example a generator falls back to for any response whose own example it cannot resolve — so it teaches `"Bad Request"` for a 404 and a 500 too. Change it to `Bad request`, or drop it as the realm-management and pairing `ErrorDetail` schemas do (neither carries one). Then extend `TestHousekeepingErrorDetailExamples`, whose comment currently says the schema-level example is "deliberately not pinned" (swagger_test.go:684-686), to pin it like the responses. One thing not to "fix": the example is nested one level in — it is the `errors` sub-schema's example, not a bare envelope — so do not re-wrap it in another `errors:`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-hk-patch-422-field-error-examples [auto]: add the three missing PATCH 422 bodies to the `ValidationErrors` response component in docs/api/astarte_housekeeping_api.yaml (yaml:427-436, referenced once at yaml:190) — it carries a single example, `error_name: [invalid_update_parameters]` (internal/housekeeping/http.go:176-177), but `patchRealm` emits three more field-error envelopes: `jwt_public_key_pem: ["can't be blank"]` (http.go:188-189), `device_registration_limit: ["is invalid"]` (http.go:193-194) and `datastream_maximum_storage_retention: ["is invalid"]` (http.go:198-199). Add them as further `examples` on the media type (OpenAPI 3.0.3, yaml:1) or split the component per case. Also record that the *messages* differ from POST for the same conditions: create rejects a negative limit/retention with the flat `ErrValidation` detail `device_registration_limit must be non-negative` / `datastream_maximum_storage_retention must be non-negative` (internal/housekeeping/service.go:143-148, 422 `ErrorDetail` shape), while PATCH answers `is invalid` in the FieldErrors shape — and that the service's own three `ErrValidation` branches (service.go:205-213) are unreachable from REST, because `patchRealm` pre-checks the identical conditions and answers first. Say in the `RealmPatch` field descriptions that a negative value is a 422 rather than silently ignored. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-hk-wrong-type-field-400 [auto]: widen the `BadRequest` component description in docs/api/astarte_housekeeping_api.yaml (yaml:378, "malformed or missing data envelope") — measured today, a perfectly well-formed `data` envelope carrying a wrong-typed field answers **400**, not the documented 422. `astarteapi.DecodeData` into `map[string]json.RawMessage` (internal/housekeeping/http.go:170) cannot fail on a type mismatch, so the failure surfaces on the second decode into the typed body (http.go:182-184) and is answered `WriteBadRequest`. Probe output: `{"data":{"device_registration_limit":"100"}}` → `json: cannot unmarshal string into Go struct field patchBody.device_registration_limit of type int32`; `{"data":{"jwt_public_key_pem":123}}` → the string twin; `{"data":{"device_registration_limit":1.5}}` → `cannot unmarshal number 1.5`; `{"data":[]}` → `cannot unmarshal array into Go value of type patchBody`. Same on POST, which decodes straight into `realmBody` (http.go:63-67). This is unreachable from a spec-conformant client — `RealmCreate`/`RealmPatch` already type both fields — so document it as a one-line extension of the 400 description on both operations rather than inventing a 422, and do NOT re-shape the `BadRequest` example, whose `detail: Bad request` is the frozen constant `TestHousekeepingErrorDetailExamples` pins. One distinction worth stating in the same description: a wrong-typed field is *not* the `invalid_update_parameters` case — unknown keys are caught earlier by the `patchAllowedFields` sweep (http.go:174-180, 422 FieldErrors) on PATCH only, and silently dropped on POST. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
EOF
wc -l .mule/todo.md
271 .mule/todo.md

Now the prose/config escalations, which need Giulio (prose is never-touch):
$ cd /root/astrate-mule && cat >> .mule/for-giulio.md <<'EOF'

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
EOF
tail -3 .mule/for-giulio.md
`internal/store/realms.go:165`, `putAuthConfig` at `internal/realm/service.go:636`), so a
two-key realm is unreachable and there is nothing to document. Worth noting because it reads
like a bug on first inspection.

$ cd /root/astrate-mule && rm -rf /tmp/opencode/probe && git status --short && echo "--- new task lines ---" && rg -c '^- \[ \] docs-sync-hk-(realm-name-pattern|errordetail-schema-example|patch-422-field-error-examples|wrong-type-field-400)' .mule/todo.md
M .mule/for-giulio.md
 M .mule/todo.md
--- new task lines ---
4

# Todos
[✓] Pick one surface and diff code routes vs docs/api paths
[✓] Spot-check 3-4 endpoints' status codes / response fields against handlers
[✓] Write up to 5 task lines to .mule/todo.md / .mule/for-giulio.md

Ran the docs-sync recipe on the **housekeeping** surface (one per run). Route sets match (5 routes, 2 paths — `docs-sync-hk-patch-endpoint` had already closed the only gap), so the drift was in status codes, examples and constraints.

**`.mule/todo.md` — 4 new task lines** (docs/api fixes, none executed):
- `docs-sync-hk-realm-name-pattern` — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and a bad name is a 422 `realm_name is invalid` out of the DB (store/realms.go:67-70 → service.go:178-181); both the `RealmName` param and `RealmCreate.realm_name` are bare strings with no `pattern`, and on GET/PATCH/DELETE the same name is a plain 404. Regex probe measured.
- `docs-sync-hk-errordetail-schema-example` — the half `docs-sync-hk-error-detail-examples-split` explicitly deferred: the schema-level example still reads `Bad Request` on the generic `ErrorDetail` all 8 responses `$ref`.
- `docs-sync-hk-patch-422-field-error-examples` — 3 of the 4 PATCH 422 bodies are undocumented, and the messages differ from POST for the same condition.
- `docs-sync-hk-wrong-type-field-400` — measured: a wrong-typed field in a valid envelope is a **400**, not the documented 422 (probe output included).

**`.mule/for-giulio.md`** — 3 prose/config escalations (`housekeeping-api.md:30` still says "Cassandra"; the page has no PATCH section; `configuration-reference.md` omits two real `[housekeeping]` keys, one of which gates the 405).

Recipe asked for no `gh issue create`, so none was run; no git touched (only the two `.mule` files are modified). Two things I checked and cleared rather than proposed: the config-key sweep, and `view()` returning only `JWTPublicKeysPEM[0]` — that reads like a bug but no writer can ever store two keys.
