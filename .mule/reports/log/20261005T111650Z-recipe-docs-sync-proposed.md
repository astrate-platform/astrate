slug: recipe-docs-sync
verdict: proposed
at:  54672da
ran: 2026-10-05T11:16:50Z on DietPi in 761s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls docs/ docs/api/ 2>/dev/null; echo "---"; ls internal/ 2>/dev/null
docs/:
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

docs/api/:
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

$ ls -la .mule/ && ls .mule/recipes/ && ls .mule/tasks/ 2>/dev/null
total 12208
drwxr-xr-x  7 root root     4096 Oct  5 12:59 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       13 Oct  5 13:04 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  5 12:59 .heartbeat
-rw-r--r--  1 root root        0 Oct  5 13:04 .last-output
-rw-r--r--  1 root root        0 Oct  5 13:04 .lock
-rw-r--r--  1 root root       10 Oct  5 13:04 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 12121043 Oct  5 13:04 cron.log
drwxr-xr-x  2 root root     4096 Oct  4 22:35 failed
-rw-r--r--  1 root root   127766 Oct  4 20:31 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    14214 Oct  5 12:59 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  4 21:08 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   142187 Oct  5 12:59 todo.md
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md

$ cat .mule/.rotation; echo "=== rotation file above ==="; cat .mule/tasks/docs-sync-appengine-data-output-params.md
docs-sync
=== rotation file above ===
# docs-sync-appengine-data-output-params — document the output/format query params on the AppEngine data GETs

Proposed by the docs-sync recipe run on 2026-09-07 (surface: appengine).

## Drift

All six data GET operations in `docs/api/astarte_appengine_api.yaml` document only
`since`, `since_after`, `to`, `limit`, `downsample_to`, `sort`
(parameter refs `DataSince`/`DataSinceAfter`/`DataTo`/`DataLimit`/`DataDownsample`/`DataSort`).
The code accepts five more query parameters via `parseQueryOpts` (`internal/appengine/http.go:554-641`),
none of which the spec mentions. The six operations and their handlers:

- `GET /appengine/v1/{realm}/devices/{device}/interfaces/{interface}` — `getData` (http.go:50)
- `GET /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}` — `getData` (http.go:51)
- `GET /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}` — `getDataByAlias` (http.go:61)
- `GET /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}` — `getDataByAlias` (http.go:62)
- `GET /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}` — `getDataInGroup` (http.go:76)
- `GET /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}` — `getDataInGroup` (http.go:77)

All three handlers go through `parseQueryOpts` (http.go:254).

## What the code actually does (verified 2026-09-07)

Three of the five missing params genuinely change the response and must be documented:

| param | parsed | consumed |
|-------|--------|----------|
| `format` (`structured`/`table`/`disjoint_tables`, default `structured`) | http.go:601-607 | data.go:237, 244 — `renderObject`/`renderIndividual` |
| `allow_bigintegers` (`true`/`false`) | http.go:612 | data.go:230, 403 — biginteger handling |
| `allow_safe_bigintegers` (`true`/`false`) | http.go:613 | data.go:230, 405 — biginteger handling |

Two are parsed but have **no consumer** anywhere in `internal/` (grep for `RetrieveMetadata`/
`DownsampleKey` finds only the struct fields in data.go:53-54 and the parse sites):

- `retrieve_metadata` (http.go:624-627 → `opts.RetrieveMetadata`) — dead, silently ignored.
- `downsample_key` (http.go:633 → `opts.DownsampleKey`) — dead, silently ignored.

Do not document these two as if they had an effect. Options for the executor, in order of
likelihood: (a) document only the three live params, (b) wire the two dead ones up if there is
an issue that wants them, (c) confirm upstream parity on the Legion Go before deciding. At
minimum the spec must stop claiming the GET is described by its current six-param-only list.

## Scope note

This does **not** overlap the blocked `appengine-snapshot-ignores-query-params` line, which is
about the interface-root snapshot branch ignoring the *window* params (`since`/`since_after`/`to`/
`limit`/`sort`). This task is about *output* params never documented at all. `format=table` also
changes the response `data` shape (a table envelope rather than an object), so consider
documenting the response-shape consequence, not just the parameter list.

## Verification

Run `make -C docs build` and confirm every YAML the Swagger UI references still loads
(the docs-sync recipe's build gate). No code changes; `docs/api/*.yaml` only.

$ tail -60 .mule/todo.md
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_10bbcc739001jdoEUu5At4nAyS

- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] forward-static-header-validation [auto]: `New` (internal/engine/forward/http.go:38-71) validates the URL and the method at boot precisely because of the rule its own comment states at 42-45 — "an unusable endpoint must fail here rather than surface once per delivery" — but `cfg.StaticHeaders` is stored unvalidated (http.go:68) and applied per delivery with `req.Header.Set` (http.go:110-112), and net/http checks header names and values only at write time. Measured on the real net/http path with a throwaway program: a name of `X Bad Name` or `X:Foo` fails per delivery with `net/http: invalid header field name`, and a value containing a newline fails with `net/http: invalid header field value for "X-Foo"` (a plain space in a value is legal). So one typo in `triggers.forward.static_headers` makes the process boot clean, accept triggers, and then count `astrate_engine_trigger_deliveries_total{outcome="failed"}` and log "custom trigger action failed" (internal/engine/triggers/actions.go:499) on every custom action forever. Validate in `New`: the RFC 7230 field-name grammar is exactly the token grammar `validMethod` (http.go:77-95) already implements, so reuse that predicate per name (rename it if you like) and add a CR/LF rejection per value. Add a table case to `TestNewRejectsBadConfig` (internal/engine/forward/http_test.go:271) asserting `New` returns an error for a bad name and a bad value and that no request is attempted; this test fails today because `New` returns nil error for both.
- [x] forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
- [x] forward-status-error-body [auto]: a non-2xx forward returns `forward: status %d` and nothing else (internal/engine/forward/http.go:121-123) after copying the **entire** response body into `io.Discard` with no limit (http.go:117-120) — the one thing an operator needs to debug a 500 from the bus is destroyed, and the sibling webhook request in the same codebase does bound its drain with `io.LimitReader(resp.Body, 1<<20)` (internal/engine/triggers/actions.go:626), so the two near-identical request paths have drifted on exactly that line. Read a bounded prefix (e.g. 512 bytes) of the body into the returned error and bound the discard the way actions.go:626 does. Add a case beside `TestStatusTable` (internal/engine/forward/http_test.go:164) where the handler answers 500 with body `{"error":"boom"}` and assert the error mentions both `500` and `boom`; it fails today because the error carries only the status code.
- [x] forward-envelope-bytes-test [auto]: `marshalEnvelope` (internal/engine/forward/envelope.go:19-29) has no direct test — `rg marshalEnvelope` finds exactly two callers (http.go:99, nats.go:45) and no test file — while envelope.go:6-8 makes the strong claim that "every Forwarder implementation in this package must produce byte-identical envelopes for the same inputs". Both suites check the shape only by unmarshalling into `bodyShape` (http_test.go:24-29), which cannot see key order or the exact bytes and would pass unchanged if a field were renamed on both sides. Add internal/engine/forward/envelope_test.go with a table pinning the exact output for (nil, nil) → `{"realm":"r","trigger":"t","action":null,"event":null}`, for the empty-non-nil pair the rule at 16-18 is about, and for a valued pair; container-free, no HTTP server needed.
- [x] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.
- [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml — the string `"403"` appears nowhere in the file and `components` defines no `Forbidden` response, yet every route is wrapped by `a.require = mw.RequireRealm(auth.ClaimAppEngine)` (internal/appengine/http.go:35), which answers `403 Forbidden` for a verified a_aea JWT whose grants do not authorize the method+path (`astarteapi.WriteForbidden`, internal/auth/middleware.go:90; `DetailForbidden`, pkg/astarteapi/envelope.go:41; pinned today by internal/appengine/http_test.go:184-185). Add a `Forbidden` response component copied from the realm-management one (docs/api/astarte_realm_management_api.yaml:1277-1285 — `ErrorDetail` schema, example `errors.detail: Forbidden`) and `$ref` it from all 32 operations, which is exactly the parity docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go:123-170) pins; add the same docs-half test for the AppEngine spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-patch-by-alias-409 [auto]: add the missing `409` response to `patchDeviceByAlias` in docs/api/astarte_appengine_api.yaml (yaml:241-291, currently 200/400/401/404/422/500) — the handler resolves the alias then calls the same `applyPatch` as the device- and group-scoped PATCH, and a rename to an alias already in use returns `ErrAliasAlreadyInUse` (internal/appengine/service.go:337) which `writeError` maps to `409 ConflictAliasInUse` (internal/appengine/http.go:650-651); the other two PATCH operations already document that response thanks to docs-sync-appengine-query-params-status and docs-sync-appengine-group-patch-status, so this is the last PATCH of the three with a hole. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
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
- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362) — the second half docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail: Bad Request`, the reconstructed title-case form that `TestHousekeepingErrorDetailExamples` (internal/swagger/swagger_test.go:687) now pins against everywhere else on this spec, while the wire emits the frozen `astarteapi.DetailBadRequest` = `"Bad request"` (pkg/astarteapi/envelope.go:37, via `WriteBadRequest` on every 400 — internal/housekeeping/http.go:65, 166, 171, 183). Because the example sits on the generic `ErrorDetail` schema, which all eight error responses `$ref` (yaml:377-477), it is the example a generator falls back to for any response whose own example it cannot resolve — so it teaches `"Bad Request"` for a 404 and a 500 too. Change it to `Bad request`, or drop it as the realm-management and pairing `ErrorDetail` schemas do (neither carries one). Then extend `TestHousekeepingErrorDetailExamples`, whose comment currently says the schema-level example is "deliberately not pinned" (swagger_test.go:684-686), to pin it like the responses. One thing not to "fix": the example is nested one level in — it is the `errors` sub-schema's example, not a bare envelope — so do not re-wrap it in another `errors:`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-patch-422-field-error-examples [auto]: add the three missing PATCH 422 bodies to the `ValidationErrors` response component in docs/api/astarte_housekeeping_api.yaml (yaml:427-436, referenced once at yaml:190) — it carries a single example, `error_name: [invalid_update_parameters]` (internal/housekeeping/http.go:176-177), but `patchRealm` emits three more field-error envelopes: `jwt_public_key_pem: ["can't be blank"]` (http.go:188-189), `device_registration_limit: ["is invalid"]` (http.go:193-194) and `datastream_maximum_storage_retention: ["is invalid"]` (http.go:198-199). Add them as further `examples` on the media type (OpenAPI 3.0.3, yaml:1) or split the component per case. Also record that the *messages* differ from POST for the same conditions: create rejects a negative limit/retention with the flat `ErrValidation` detail `device_registration_limit must be non-negative` / `datastream_maximum_storage_retention must be non-negative` (internal/housekeeping/service.go:143-148, 422 `ErrorDetail` shape), while PATCH answers `is invalid` in the FieldErrors shape — and that the service's own three `ErrValidation` branches (service.go:205-213) are unreachable from REST, because `patchRealm` pre-checks the identical conditions and answers first. Say in the `RealmPatch` field descriptions that a negative value is a 422 rather than silently ignored. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [!] docs-sync-hk-wrong-type-field-400 [auto]: widen the `BadRequest` component description in docs/api/astarte_housekeeping_api.yaml (yaml:378, "malformed or missing data envelope") — measured today, a perfectly well-formed `data` envelope carrying a wrong-typed field answers **400**, not the documented 422. `astarteapi.DecodeData` into `map[string]json.RawMessage` (internal/housekeeping/http.go:170) cannot fail on a type mismatch, so the failure surfaces on the second decode into the typed body (http.go:182-184) and is answered `WriteBadRequest`. Probe output: `{"data":{"device_registration_limit":"100"}}` → `json: cannot unmarshal string into Go struct field patchBody.device_registration_limit of type int32`; `{"data":{"jwt_public_key_pem":123}}` → the string twin; `{"data":{"device_registration_limit":1.5}}` → `cannot unmarshal number 1.5`; `{"data":[]}` → `cannot unmarshal array into Go value of type patchBody`. Same on POST, which decodes straight into `realmBody` (http.go:63-67). This is unreachable from a spec-conformant client — `RealmCreate`/`RealmPatch` already type both fields — so document it as a one-line extension of the 400 description on both operations rather than inventing a 422, and do NOT re-shape the `BadRequest` example, whose `detail: Bad request` is the frozen constant `TestHousekeepingErrorDetailExamples` pins. One distinction worth stating in the same description: a wrong-typed field is *not* the `invalid_update_parameters` case — unknown keys are caught earlier by the `patchAllowedFields` sweep (http.go:174-180, 422 FieldErrors) on PATCH only, and silently dropped on POST. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `components.parameters.RealmName` and `RealmCreate.realm_name` both carry `pattern: '^[a-z][a-z0-9]*$'` and `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go) holds both against the realm name CHECK extracted from migrations/000002_metadata.up.sql — but the read side is still bare `type: string`: `Realm.realm_name` (components.schemas.Realm) and `DataEnvelopeRealmNames.items`. Decide whether a response field earns the same pattern (a generated client learns nothing new from it, but a reader diffing create against get sees one constrained and one not) and, if it does, add it to both and extend that test to hold the response-side pattern against the same CHECK. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] store-register-inhibit-preserve: RegisterDevice must not overwrite status to 'registered' when a device was inhibited before first credentials request (preserve 'inhibited' state). File: internal/store/devices.go — BLOCKED: gates failed
- [x] store-aliasvalues-self-exclusion-test: Add test for AliasValuesTaken to assert self-exclusion (same device's own alias values do not count as taken). File: internal/store/devices_test.go
- [!] store-latestindividual-empty-errnotfound: Add store-level test asserting LatestIndividual returns ErrNotFound for empty series and returns newest on multiple samples. File: internal/store/datastreams_test.go — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [x] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.go

$ cat .mule/config; echo "=== log tail ==="; tail -40 .mule/log.md
# mule config for this repo. Sourced as bash; env vars override it.
#
# This file is read on the Mac AND on the Pi, which do not have the same toolchain, so it
# detects rather than hardcodes. Measured facts behind the detection:
#
#   Mac (arm64 darwin) : go test -race ./...  green, golangci-lint pinned in ~/go/bin
#   Pi  (arm64 linux)  : go test ./...        green in ~3m over 20 packages
#                        go test -race        IMPOSSIBLE — see below
#
# **The race detector cannot run on the Pi.** ThreadSanitizer needs a 48-bit VMA and this
# kernel is built with 39 ("FATAL: ThreadSanitizer: unsupported VMA range / Found 39 -
# Supported 48"). That is a kernel config, not a flake, and no amount of retrying changes it.
# So the unattended gate is weaker than the Mac's by construction: a data race passes on the
# Pi. Race coverage lives on the Legion Go instead (Go 1.26.5 at ~/.local/go, full suite
# green in ~40s on 16 cores) via the standing [legion] race-check task — which is therefore
# not optional housekeeping but the concurrency gate itself.

_gcl="$(command -v golangci-lint || echo "$HOME/go/bin/golangci-lint")"

if go test -race -run XXX_nonexistent ./pkg/payload >/dev/null 2>&1; then
  MULE_TEST_CMD="go test -race ./..."
else
  # No race detector here. `go vet` is not a replacement, but it is free and it catches a
  # different class of real defect, so it is not nothing.
  MULE_TEST_CMD="go vet ./... && go test ./..."
fi

if [ -x "$_gcl" ]; then
  MULE_LINT_CMD="$_gcl run ./..."
  MULE_FIX_CMD="gofmt -w . && $_gcl run --fix ./... >/dev/null 2>&1 || true"
else
  MULE_LINT_CMD=""            # absent on the Pi; gofmt still runs
  MULE_FIX_CMD="gofmt -w ."
fi

# Paths the mule must never modify: frozen spec, schema decisions, Giulio's voice, build pins.
MULE_NEVER="docs/DESIGN.md docs/ROADMAP.md docs/site/* migrations/* .github/* deploy/* go.mod go.sum Dockerfile docker-compose.yml .trickle/*"

# The mule never commits to main.
MULE_BRANCH="mule/queue"

# Per-task budget in seconds. A task needing more than this was specced too big. Generous
# here because the Pi is slow: the gate alone is 3 minutes of it.
MULE_TIMEOUT="1200"

# How many tasks the timer may run per day. The provider is free; do not hammer it. At one
# task per 30-minute tick the ceiling would be 48, so this is a deliberate ~3x reduction.
MULE_DAILY_MAX="16"

MULE_AGENT="build"
# Pin the model explicitly. Leaving this empty lets opencode resolve the agent's default,
# and since 2026-08-27 that resolution answers "UnknownError: Unexpected server error" every
# single time — it killed four nights of the daily survey. `--model opencode/big-pickle`
# succeeds on the same box, same agent, same second.
MULE_MODEL="opencode/big-pickle"
=== log tail ===
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

$ rg -c 'docs-sync' .mule/todo.md; echo "--- native ---"; rg -n 'docs-sync-native' .mule/todo.md | wc -l; echo "--- hk ---"; rg -n 'docs-sync-hk' .mule/todo.md | wc -l
76
--- native ---
8
--- hk ---
14

# Todos
[•] Pick one surface and list handler routes vs documented paths (native)
[ ] Verify status codes / response fields for 3-4 endpoints against handlers
[ ] Check ASTRATE_ config keys vs docs/site/configuration-reference.md
[ ] Append up to 5 task lines to .mule/todo.md (+ for-giulio.md if prose contradicts code)

$ rg -n '^\s{2}/' docs/api/astrate_native_api.yaml
26:  /astrate/v1/health:
44:  /astrate/v1/readiness:
77:  /astrate/v1/metrics:
99:  /appengine/health:
127:  /realmmanagement/health:
155:  /pairing/health:
185:  /appengine/version:
205:  /realmmanagement/version:
225:  /pairing/version:
245:  /housekeeping/version:
267:  /astrate/v1/{realm}/socket:
332:  /appengine/v1/socket/websocket:

$ ls internal/flowapi/ && echo "=== routes ===" && rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/flowapi/
composite_wire_test.go
http.go
http_test.go
service.go
service_integration_test.go
service_test.go
userblocks.go
userblocks_test.go
=== routes ===

$ rg -n 'Handle|Get\(|Post\(|mux\.|Route' internal/flowapi/http.go | head -60
17:// API is the /flow/v1 HTTP surface. Routes accept a realm JWT carrying the
22:	require func(http.Handler) http.Handler
30:// Mount registers the routes on mux.
32:	h := func(f http.HandlerFunc) http.Handler { return a.require(f) }
33:	mux.Handle("GET /flow/v1/{realm}/pipelines", h(a.listPipelines))
34:	mux.Handle("POST /flow/v1/{realm}/pipelines", h(a.createPipeline))
35:	mux.Handle("GET /flow/v1/{realm}/pipelines/{name}", h(a.getPipeline))
36:	mux.Handle("PUT /flow/v1/{realm}/pipelines/{name}", h(a.updatePipeline))
37:	mux.Handle("DELETE /flow/v1/{realm}/pipelines/{name}", h(a.deletePipeline))
39:	mux.Handle("GET /flow/v1/{realm}/flows", h(a.listFlows))
40:	mux.Handle("POST /flow/v1/{realm}/flows", h(a.startFlow))
41:	mux.Handle("GET /flow/v1/{realm}/flows/{name}", h(a.getFlow))
42:	mux.Handle("DELETE /flow/v1/{realm}/flows/{name}", h(a.stopFlow))
43:	mux.Handle("POST /flow/v1/{realm}/flows/{name}/reload", h(a.reloadFlow))
44:	mux.Handle("PUT /flow/v1/{realm}/flows/{name}/config", h(a.updateFlowConfig))
46:	mux.Handle("GET /flow/v1/{realm}/blocks", h(a.listBlocks))
47:	mux.Handle("POST /flow/v1/{realm}/blocks", h(a.createUserBlock))
48:	mux.Handle("GET /flow/v1/{realm}/blocks/{type}", h(a.getBlock))
49:	mux.Handle("PUT /flow/v1/{realm}/blocks/{name}", h(a.updateUserBlock))
50:	mux.Handle("DELETE /flow/v1/{realm}/blocks/{name}", h(a.deleteUserBlock))

$ rg -rn 'astrate/v1|appengine/v1/socket' cmd/ internal/ --glob '!*_test.go' | head -40
cmd/astrate/main.go:	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+addr+"/n/readiness", nil) //nolint:gosec // G704: self-probe of the local readiness endpoint
internal/appengine/stream/ws.go:// at /n/{realm}/socket fed by the engine's in-process fan-out bus.
internal/appengine/stream/ws.go:	mux.Handle("GET /n/{realm}/socket", a.require(http.HandlerFunc(a.handle)))
internal/appengine/channels/protocol.go:// served at /n/websocket. It is the compatibility counterpart
internal/appengine/channels/ws.go:	mux.Handle("GET /n/websocket", http.HandlerFunc(a.handle))
internal/config/config.example.toml:# and the Astrate-native /n health/metrics/socket endpoints.
internal/observability/compat.go:// same dependency probe that backs /n/readiness — and answers 503
internal/observability/metrics.go:// health/readiness/metrics HTTP surface under /n (docs/DESIGN.md
internal/observability/health.go:// /n so they never collide with the upstream API namespace.
internal/observability/health.go:	mux.HandleFunc("GET /n/health", h.handleHealth)
internal/observability/health.go:	mux.HandleFunc("GET /n/readiness", h.handleReadiness)
internal/observability/health.go:	mux.Handle("GET /n/metrics", h.metrics)

$ rg -n 'mux\.(Handle|HandleFunc)' internal/ cmd/ --glob '!*_test.go'
cmd/astrate/main.go:436:	mux.Handle("GET /appengine/v1/{realm}/version",
cmd/astrate/main.go:438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
internal/httpx/notfound.go:16:		_, pattern := mux.Handler(r)
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
internal/observability/compat.go:20:	mux.HandleFunc("GET /"+service+"/health", func(w http.ResponseWriter, r *http.Request) {
internal/observability/compat.go:51:	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
internal/pairing/http.go:81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/pairing/http.go:82:	mux.Handle("POST /pairing/v1/{realm}/agent/devices",
internal/pairing/http.go:84:	mux.Handle("DELETE /pairing/v1/{realm}/agent/devices/{deviceID}",
internal/pairing/http.go:86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
internal/pairing/http.go:88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
internal/pairing/http.go:90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
internal/observability/health.go:44:	mux.HandleFunc("GET /astrate/v1/health", h.handleHealth)
internal/observability/health.go:45:	mux.HandleFunc("GET /astrate/v1/readiness", h.handleReadiness)
internal/observability/health.go:46:	mux.Handle("GET /astrate/v1/metrics", h.metrics)
internal/housekeeping/http.go:37:	mux.Handle("GET /housekeeping/v1/realms", h(a.listRealms))
internal/housekeeping/http.go:38:	mux.Handle("POST /housekeeping/v1/realms", h(a.createRealm))
internal/housekeeping/http.go:39:	mux.Handle("GET /housekeeping/v1/realms/{realm}", h(a.getRealm))
internal/housekeeping/http.go:40:	mux.Handle("PATCH /housekeeping/v1/realms/{realm}", h(a.patchRealm))
internal/housekeeping/http.go:41:	mux.Handle("DELETE /housekeeping/v1/realms/{realm}", h(a.deleteRealm))

$ rg -n 'flow' docs/api/astrate_native_api.yaml | head -20; echo "=== paths in spec ==="; rg -n 'operationId:' docs/api/astrate_native_api.yaml
=== paths in spec ===
28:      operationId: getHealth
46:      operationId: getReadiness
79:      operationId: getMetrics
101:      operationId: getAppEngineHealth
129:      operationId: getRealmManagementHealth
157:      operationId: getPairingHealth
187:      operationId: getAppEngineVersion
207:      operationId: getRealmManagementVersion
227:      operationId: getPairingVersion
247:      operationId: getHousekeepingVersion
269:      operationId: nativeWebSocket
334:      operationId: phoenixWebSocket

$ wc -l docs/api/astrate_native_api.yaml; echo ---; sed -n '1,120p' docs/api/astrate_native_api.yaml
466 docs/api/astrate_native_api.yaml
---
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
        type is `text/plain; version=0.0.4`.
      tags: [Observability]
      responses:
        "200":
          description: Prometheus metrics.
          content:
            text/plain; version=0.0.4:
              schema:
                type: string
              example: |
                # HELP astrate_broker_sessions Live authenticated MQTT device sessions.
                # TYPE astrate_broker_sessions gauge
                astrate_broker_sessions 0

  # ── Compat health (Astarte Dashboard) ───────────────────────────────

  /appengine/health:
    get:
      operationId: getAppEngineHealth
      summary: AppEngine compat health
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

$ sed -n '420,450p' cmd/astrate/main.go; echo "=== compat.go ==="; cat internal/observability/compat.go
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
=== compat.go ===
package observability

import (
	"context"
	"encoding/json"
	"net/http"
)

// MountServiceCompat registers the upstream-parity unauthenticated
// GET /{service}/health endpoint. The Astarte Dashboard polls one per service
// (appengine, realmmanagement, pairing) for its API status indicators; only
// the status code matters to it, but the body mirrors upstream's envelope.
//
// Deliberate deviation 18 (docs/COMPATIBILITY.md): upstream answers a static 200
// here, so the indicator it feeds cannot ever go red. Astrate runs check — the
// same dependency probe that backs /astrate/v1/readiness — and answers 503
// when it fails, which is strictly safer for whoever is reading the Dashboard.
// A nil check keeps upstream's static 200.
func MountServiceCompat(mux *http.ServeMux, service string, check Check) {
	mux.HandleFunc("GET /"+service+"/health", func(w http.ResponseWriter, r *http.Request) {
		if check != nil {
			ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
			defer cancel()
			if err := check(ctx); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte(`{"data":{"status":"unhealthy"}}`))
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"status":"ok"}}`))
	})
}

// VersionHandler answers GET .../version with the upstream envelope
// {"data":"<version>"} (measured upstream 1.2.0, verify-versions.json).
func VersionHandler(version string) http.HandlerFunc {
	body, err := json.Marshal(map[string]string{"data": version})
	if err != nil {
		body = []byte(`{"data":""}`)
	}
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
}

// MountVersionCompat registers GET /{service}/version for one service.
func MountVersionCompat(mux *http.ServeMux, service, version string) {
	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
}

→ Read docs/api/astrate_native_api.yaml [offset=120, limit=240]

$ rg -n 'version' docs/api/astrate_appengine_api.yaml docs/api/astrate_pairing_api.yaml docs/api/astrate_realm_management_api.yaml docs/api/astrate_housekeeping_api.yaml | head -30
rg: docs/api/astrate_appengine_api.yaml: No such file or directory (os error 2)
rg: docs/api/astrate_pairing_api.yaml: No such file or directory (os error 2)
rg: docs/api/astrate_realm_management_api.yaml: No such file or directory (os error 2)
rg: docs/api/astrate_housekeeping_api.yaml: No such file or directory (os error 2)

→ Read docs/api/astrate_native_api.yaml [offset=359]

$ rg -n 'version' docs/api/astrate_appengine_api.yaml docs/api/astrate_pairing_api.yaml docs/api/astrate_realm_management_api.yaml docs/api/astrate_housekeeping_api.yaml
rg: docs/api/astrate_appengine_api.yaml: No such file or directory (os error 2)
rg: docs/api/astrate_pairing_api.yaml: No such file or directory (os error 2)
rg: docs/api/astrate_realm_management_api.yaml: No such file or directory (os error 2)
rg: docs/api/astrate_housekeeping_api.yaml: No such file or directory (os error 2)

$ rg -n 'MountServiceCompat|MountVersionCompat|mountRealmVersion' cmd/ internal/ --glob '!*_test.go'
cmd/astrate/main.go:387:		observability.MountServiceCompat(mux, svc, st.Health)
cmd/astrate/main.go:397:		observability.MountVersionCompat(mux, svc, version)
cmd/astrate/main.go:399:	mountRealmVersion(mux, mw)
cmd/astrate/main.go:429:// mountRealmVersion registers the realm-scoped version endpoints observed on
cmd/astrate/main.go:435:func mountRealmVersion(mux *http.ServeMux, mw *auth.Middleware) {
internal/observability/compat.go:9:// MountServiceCompat registers the upstream-parity unauthenticated
internal/observability/compat.go:19:func MountServiceCompat(mux *http.ServeMux, service string, check Check) {
internal/observability/compat.go:49:// MountVersionCompat registers GET /{service}/version for one service.
internal/observability/compat.go:50:func MountVersionCompat(mux *http.ServeMux, service, version string) {

$ ls -la docs/api/
total 176
drwxr-xr-x 2 root root  4096 Oct  4 19:17 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 59526 Oct  3 13:05 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 17437 Oct  4 19:36 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 21724 Sep 28 19:59 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 47623 Oct  2 22:09 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 14977 Oct  3 18:49 astrate_native_api.yaml

✱ Grep "/version" in docs/api · 10 matches
→ Read cmd/astrate/main.go [offset=370, limit=45]
→ Read internal/observability/health.go
→ Read internal/appengine/stream/ws.go

$ rg -n 'func Test.*Native|astrate_native' internal/swagger/swagger_test.go
629:func TestNativeErrorDetailExamples(t *testing.T) {
630:	b, err := docs.APIYAML.ReadFile("api/astrate_native_api.yaml")
632:		t.Fatalf("reading astrate_native_api.yaml: %v", err)

→ Read docs/api/astarte_pairing_api.yaml [limit=50, offset=68]

$ rg -n 'APICompatVersion' internal/realm/*.go | head
internal/realm/dashboard_compat_test.go:35:		if v != APICompatVersion {
internal/realm/dashboard_compat_test.go:36:			t.Errorf("version = %q, want %q", v, APICompatVersion)
internal/realm/http.go:119:// feature UI on it — see APICompatVersion).
internal/realm/http.go:121:	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
internal/realm/service.go:580:// APICompatVersion is the upstream Realm Management API level Astrate
internal/realm/service.go:588:const APICompatVersion = "1.2.2"

→ Read docs/api/astarte_realm_management_api.yaml [limit=40, offset=657]

$ ls docs/site/; echo "=== flow mentions in site ==="; rg -ln 'flow/v1' docs/ --glob '!docs/api/*'
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
=== flow mentions in site ===
docs/handoff/session-2026-07-29-handoff.md
docs/DESIGN.md
docs/handoff/session-2026-07-29-flow-factory-handoff.md
docs/handoff/flow-parity-audit-2026-07-29.md
docs/handoff/flow-lifecycle-followups-design-2026-08-22.md
docs/handoff/README.md
docs/handoff/flow-design-a-named-durable-flows-2026-07-29.md
docs/handoff/flow-design-b-container-block-2026-07-29.md
docs/site/DESIGN.md

$ rg -o '\bASTRATE_[A-Z_]+' -N internal/ cmd/ 2>/dev/null | sed 's/.*://' | sort -u > /tmp/opencode/code_keys.txt; wc -l /tmp/opencode/code_keys.txt; cat /tmp/opencode/code_keys.txt
24 /tmp/opencode/code_keys.txt
ASTRATE_DATABASE_DSN
ASTRATE_ENGINE_SHARDS
ASTRATE_FLOW_CONFIG
ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
ASTRATE_HTTP_ADDR
ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
ASTRATE_HTTP_TLS_CERT_FILE
ASTRATE_HTTP_TLS_KEY_FILE
ASTRATE_LOG_FORMAT
ASTRATE_LOG_LEVEL
ASTRATE_MASTER_KEY
ASTRATE_MASTER_KEY_FILE
ASTRATE_MQTT_ADDR
ASTRATE_MQTT_ADVERTISED_URL
ASTRATE_MQTT_INSECURE_DEV_MODE
ASTRATE_MQTT_SESSION_STORE_PATH
ASTRATE_MQTT_TLS_CERT_FILE
ASTRATE_MQTT_TLS_KEY_FILE
ASTRATE_REALM_JWT_PUBLIC_KEY
ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
ASTRATE_REALM_NAME
ASTRATE_SECURITY_MASTER_KEY_FILE
ASTRATE_TEST_DSN

$ rg -n 'ASTRATE_FLOW_CONFIG|ASTRATE_TEST_DSN|ASTRATE_ENGINE_SHARDS|ASTRATE_SECURITY_MASTER_KEY_FILE' docs/ | head -20
docs/OPERATIONS.md:30:`ASTRATE_MQTT_INSECURE_DEV_MODE`, `ASTRATE_ENGINE_SHARDS`, `ASTRATE_LOG_LEVEL`,
docs/ROADMAP.md:112:| 0.5 | `internal/testutil/pg.go` | `StartTimescale(t) *pgxpool.Pool`: container (or `ASTRATE_TEST_DSN` reuse for speed), waits ready, returns pool + teardown | 90 |
docs/site/configuration-reference.md:60:| `shards` | int | `16` | `ASTRATE_ENGINE_SHARDS` | No | Number of ordered ingestion shards. Each shard processes messages strictly in order. |
docs/site/configuration-reference.md:95:| `master_key_file` | string | `""` | `ASTRATE_SECURITY_MASTER_KEY_FILE` | No | File holding the AES-256 master key that seals realm CA private keys. When empty, falls back to `ASTRATE_MASTER_KEY` (64 hex chars) or `ASTRATE_MASTER_KEY_FILE`. |
docs/site/deployment.md:63:Named `ASTRATE_<SECTION>_<FIELD>` -- e.g. `ASTRATE_HTTP_ADDR`, `ASTRATE_MQTT_ADDR`, `ASTRATE_ENGINE_SHARDS`, `ASTRATE_LOG_LEVEL`.
docs/site/operations.md:30:`ASTRATE_MQTT_INSECURE_DEV_MODE`, `ASTRATE_ENGINE_SHARDS`, `ASTRATE_LOG_LEVEL`,
docs/site/ROADMAP.md:112:| 0.5 | `internal/testutil/pg.go` | `StartTimescale(t) *pgxpool.Pool`: container (or `ASTRATE_TEST_DSN` reuse for speed), waits ready, returns pool + teardown | 90 |
docs/site/contributing.md:76:ASTRATE_TEST_DSN="postgres://user:pass@host:5432/dbname?sslmode=disable" make test-integration
docs/handoff/flow-design-b-container-block-2026-07-29.md:132:  - `ASTRATE_FLOW_CONFIG=<json>` of nested `config`
docs/handoff/flow-lifecycle-followups-design-2026-08-22.md:126:  `ASTRATE_FLOW_CONFIG` at container start (`docker.go:90`) — unchangeable without a

$ ls docs/site/api/; echo "=== documented keys ==="; rg -o '\bASTRATE_[A-Z_]+' -N docs/site/configuration-reference.md | sort -u
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
=== documented keys ===
ASTRATE_DATABASE_DSN
ASTRATE_ENGINE_SHARDS
ASTRATE_HTTP_ADDR
ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
ASTRATE_HTTP_TLS_CERT_FILE
ASTRATE_HTTP_TLS_KEY_FILE
ASTRATE_LOG_FORMAT
ASTRATE_LOG_LEVEL
ASTRATE_MASTER_KEY
ASTRATE_MASTER_KEY_FILE
ASTRATE_MQTT_ADDR
ASTRATE_MQTT_ADVERTISED_URL
ASTRATE_MQTT_INSECURE_DEV_MODE
ASTRATE_MQTT_SESSION_STORE_PATH
ASTRATE_MQTT_TLS_CERT_FILE
ASTRATE_MQTT_TLS_KEY_FILE
ASTRATE_REALM_JWT_PUBLIC_KEY
ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
ASTRATE_REALM_NAME
ASTRATE_SECURITY_MASTER_KEY_FILE

$ rg -n 'ASTRATE_HOUSEKEEPING|ASTRATE_FLOW_CONFIG' internal/ cmd/
internal/config/config.go:279:		"ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
internal/config/config.go:297:	if v, ok := os.LookupEnv("ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"); ok {
internal/config/config_test.go:392:		prefixed = "ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION"
internal/config/config_test.go:533:	const env = "ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"
internal/flow/blocks/info.go:116:		Config:  "image (required), config (object→ASTRATE_FLOW_CONFIG), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",
internal/flow/blocks/schema.go:19:	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as ASTRATE_FLOW_CONFIG"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
internal/flow/blocks/container/block.go:45:	Nested        map[string]any // opaque JSON for ASTRATE_FLOW_CONFIG
internal/flow/blocks/container/block_test.go:125:	// silently dropped from ASTRATE_FLOW_CONFIG, so it must be refused.
internal/flow/blocks/container/docker.go:19:	// FlowConfigJSON is passed as ASTRATE_FLOW_CONFIG.
internal/flow/blocks/container/docker.go:105:		args = append(args, "-e", "ASTRATE_FLOW_CONFIG="+spec.FlowConfigJSON)
internal/flow/blocks/container/docker.go:213:// encodeFlowConfigJSON marshals nested config for ASTRATE_FLOW_CONFIG.

✱ Grep "housekeeping|realm_deletion_disabled|datastream_maximum_storage_retention|default_datastream" in docs/site/configuration-reference.md · 22 matches
→ Read docs/site/configuration-reference.md [limit=30, offset=74]
→ Read internal/config/config.go [limit=40, offset=265]

$ rg -n 'realm_deletion_disabled|default_datastream|Retain|DeletionDisabled' internal/config/config.example.toml
(no output)

$ rg -n -A14 'Housekeeping struct' internal/config/config.go
(no output)

$ rg -n -B2 -A18 'Housekeeping ' internal/config/config.go | head -40
47-	Engine       EngineConfig       `toml:"engine"`
48-	Pairing      PairingConfig      `toml:"pairing"`
49:	Housekeeping HousekeepingConfig `toml:"housekeeping"`
50-	Storage      StorageConfig      `toml:"storage"`
51-	Security     SecurityConfig     `toml:"security"`
52-	Realm        RealmConfig        `toml:"realm"`
53-	Log          LogConfig          `toml:"log"`
54-	Triggers     TriggersConfig     `toml:"triggers"`
55-}
56-
57-// HTTPConfig is the single REST listener (§3.7). TLS is optional: leave the
58-// files empty to serve plaintext behind a TLS-terminating reverse proxy.
59-// CORSAllowedOrigins enables cross-origin browser clients (the Astarte
60-// Dashboard SPA): each entry is "*" or an absolute origin such as
61-// "http://localhost:4040"; empty (the default) disables CORS entirely.
62-type HTTPConfig struct {
63-	Addr               string   `toml:"addr"`
64-	TLSCertFile        string   `toml:"tls_cert_file"`
65-	TLSKeyFile         string   `toml:"tls_key_file"`
66-	CORSAllowedOrigins []string `toml:"cors_allowed_origins"`
67-}

$ rg -n -A16 'type HousekeepingConfig struct' internal/config/config.go
111:type HousekeepingConfig struct {
112-	JWTPublicKeys     []string `toml:"jwt_public_keys"`
113-	JWTPublicKeyFiles []string `toml:"jwt_public_key_files"`
114-	// DefaultDatastreamMaximumStorageRetention is the realm-level datastream
115-	// storage ceiling (seconds) injected at realm creation when the caller
116-	// omits the field (#73). When unset (nil) no default is injected and
117-	// existing deployments behave exactly as before.
118-	DefaultDatastreamMaximumStorageRetention *int64 `toml:"default_datastream_maximum_storage_retention"`
119-	// RealmDeletionDisabled gates DELETE /housekeeping/v1/realms/{realm}
120-	// (#75): when true, deletion answers upstream's 405 "Realm deletion
121-	// disabled". Opt-in — unset (false) keeps the historical
122-	// always-delete behavior.
123-	RealmDeletionDisabled bool `toml:"realm_deletion_disabled"`
124-}
125-
126-// StorageConfig holds runtime storage policy. Retention applies a global
127-// TimescaleDB drop-chunks policy; zero (the default) disables it, leaving

$ rg -n 'appengine/v1/\{realm\}/version|/version' .mule/todo.md | head -20; echo "=== hk version in for-giulio ==="; rg -n 'housekeeping/version|/housekeeping/health|realm_deletion_disabled|default_datastream_maximum_storage_retention|HOUSEKEEPING_DEFAULT' .mule/for-giulio.md | head -20
123:- [x] docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` route to docs/api/astarte_realm_management_api.yaml — it exists in code (internal/realm/http.go:51, handler getDatastreamMaximumStorageRetention at http.go:135, returns 200 with a data-envelope retention value via Service.GetDatastreamMaximumStorageRetention) but the spec jumps straight from config/device_registration_limit (yaml:428) to /version. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
156:- [x] docs-sync-pairing-version-endpoint [auto]: add the undocumented `GET /pairing/v1/{realm}/version` route to docs/api/astarte_pairing_api.yaml — it exists in code (cmd/astrate/main.go:398, unauthenticated, upstream-parity per main.go:387-391) and answers 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38); mirror the RM spec's getVersion op (astarte_realm_management_api.yaml:532) without its a_rma security. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
157:- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
189:- [x] docs-sync-pairing-version-404-unreachable [auto]: remove the `404` "Unknown realm" response from `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml (yaml:96-104) — the route is served by `observability.VersionHandler` (cmd/astrate/main.go:398, internal/observability/compat.go:38-47), which answers static `200 {"data":"<version>"}` with no realm lookup, so a 404 can never be emitted there. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
194:- [x] docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
196:- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
=== hk version in for-giulio ===
126:  `default_datastream_maximum_storage_retention` (`*int64`, #73) and `realm_deletion_disabled`
128:  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (config.go:277-292, which
129:  also accepts the bare upstream name `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`
133:  (`realm_deletion_disabled` makes `DELETE /housekeeping/v1/realms/{realm}` answer 405) and the
195:- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`ASTRATE_FLOW_CONFIG`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
265:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
349:  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` and
581:- `default_datastream_maximum_storage_retention` (`*int64`, env
582:  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`, validated as a
586:- `realm_deletion_disabled` (bool, env `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`, a

$ rg -n 'appengine/v1/\{realm\}/version|appengine.*version' .mule/for-giulio.md | head -20
(no output)

$ rg -n 'appengine.*realm-scoped|realm-scoped.*version|getAppEngineVersion|ClaimAppEngine' .mule/todo.md | head -20
157:- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
196:- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
236:- [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml — the string `"403"` appears nowhere in the file and `components` defines no `Forbidden` response, yet every route is wrapped by `a.require = mw.RequireRealm(auth.ClaimAppEngine)` (internal/appengine/http.go:35), which answers `403 Forbidden` for a verified a_aea JWT whose grants do not authorize the method+path (`astarteapi.WriteForbidden`, internal/auth/middleware.go:90; `DetailForbidden`, pkg/astarteapi/envelope.go:41; pinned today by internal/appengine/http_test.go:184-185). Add a `Forbidden` response component copied from the realm-management one (docs/api/astarte_realm_management_api.yaml:1277-1285 — `ErrorDetail` schema, example `errors.detail: Forbidden`) and `$ref` it from all 32 operations, which is exactly the parity docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go:123-170) pins; add the same docs-half test for the AppEngine spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ rg -n -B4 -A30 'func (m \*Middleware) RequireRealmAny' internal/auth/middleware.go
(no output)

$ rg -n 'ClaimChannels|ClaimAppEngine|ClaimRealmManagement|ClaimPairingAdmin|ClaimHousekeeping' internal/auth/*.go | grep -v _test | head -20
internal/auth/claims.go:20:	// ClaimAppEngine authorizes the AppEngine API (a_aea).
internal/auth/claims.go:21:	ClaimAppEngine Claim = "a_aea"
internal/auth/claims.go:22:	// ClaimChannels authorizes Astarte Channels; Astrate honours it on the
internal/auth/claims.go:24:	ClaimChannels Claim = "a_ch"
internal/auth/claims.go:28:	// ClaimHousekeeping authorizes the Housekeeping API (a_ha).
internal/auth/claims.go:29:	ClaimHousekeeping Claim = "a_ha"
internal/auth/claims.go:32:	// ClaimRealmManagement authorizes the Realm Management API (a_rma).
internal/auth/claims.go:33:	ClaimRealmManagement Claim = "a_rma"
internal/auth/claims.go:134:		ClaimAppEngine:       c.AppEngine,
internal/auth/claims.go:135:		ClaimChannels:        c.Channels,
internal/auth/claims.go:137:		ClaimHousekeeping:    c.Housekeeping,
internal/auth/claims.go:139:		ClaimRealmManagement: c.RealmManagement,
internal/auth/claims.go:198:	for _, g := range t.grants[ClaimChannels] {

✱ Grep "func .*RequireRealm" in internal/auth · 3 matches
→ Read internal/auth/middleware.go [limit=70, offset=30]

$ rg -n -B6 -A30 'func .*tryAuthorize' internal/auth/middleware.go
124-	authFailed                      // authentication failed (401 written); stop
125-	authForbidden                   // claim did not grant the request; try next
126-)
127-
128-// tryAuthorize is authorize without the terminal 403: it classifies the
129-// attempt. Callers write their own forbidden response on authForbidden.
130:func (m *Middleware) tryAuthorize(w http.ResponseWriter, r *http.Request, next http.Handler, claim Claim, keysPEM []string, base string) authResult {
131-	tokenString, ok := bearerToken(r)
132-	if !ok {
133-		_ = astarteapi.WriteUnauthorized(w)
134-		return authFailed
135-	}
136-
137-	tok, err := m.cache.Verify(tokenString, keysPEM)
138-	if err != nil {
139-		_ = astarteapi.WriteUnauthorized(w)
140-		return authFailed
141-	}
142-
143-	authPath, ok := RelativePath(r.URL.Path, base)
144-	if !ok {
145-		// Upstream parity: a path the authorizer cannot anchor is an
146-		// authorization failure (403), not an authentication one.
147-		return authForbidden
148-	}
149-	if !tok.Authorizes(claim, r.Method, authPath) {
150-		return authForbidden
151-	}
152-
153-	ctx := context.WithValue(r.Context(), tokenContextKey{}, tok)
154-	next.ServeHTTP(w, r.WithContext(ctx))
155-	return authServed
156-}
157-
158-// RelativePath computes the authorization path with upstream parity (Astarte's
159-// GuardianAuthorizePath plug): split the URL path into segments, drop
160-// everything up to and including the first segment equal to base, and join

→ Read internal/auth/middleware.go [limit=60, offset=136]

$ rg -n -B4 -A25 'func (h \*Health) metrics|metrics http.Handler|func.*Handler\(\)' internal/observability/metrics.go | head -60
37-// bus collectors attach here).
38-func (m *Metrics) Registerer() prometheus.Registerer { return m.reg }
39-
40-// Handler serves the registry in the Prometheus text exposition format.
41:func (m *Metrics) Handler() http.Handler {
42-	return promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})
43-}
44-
45-// DBPoolStats is the subset of pgxpool.Stat exposed as gauges; the caller maps
46-// store.Stat() onto it (keeping pgx out of this package).
47-type DBPoolStats struct {
48-	AcquiredConns int32
49-	IdleConns     int32
50-	TotalConns    int32
51-	MaxConns      int32
52-}
53-
54-// RegisterBrokerSessions publishes astrate_broker_sessions, read on scrape from
55-// fn (broker.SessionCount).
56-func (m *Metrics) RegisterBrokerSessions(fn func() float64) {
57-	m.reg.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
58-		Namespace: metricsNamespace,
59-		Subsystem: "broker",
60-		Name:      "sessions",
61-		Help:      "Live authenticated MQTT device sessions.",
62-	}, fn))
63-}
64-
65-// RegisterDBPool publishes the astrate_db_pool_* gauges, read on scrape from fn
66-// (a snapshot of store.Stat()).

$ rg -n 'rooms:|new_event|event_data|phx_join|heartbeat|"watch"|"unwatch"' internal/appengine/channels/*.go | grep -v _test | head -40
internal/appengine/channels/ws.go:23:// Event-name constants for the watch/unwatch/new_event cycle.
internal/appengine/channels/ws.go:25:	EventWatch    = "watch"
internal/appengine/channels/ws.go:26:	EventUnwatch  = "unwatch"
internal/appengine/channels/ws.go:27:	EventNewEvent = "new_event"
internal/appengine/channels/ws.go:88:		rooms: make(map[string]*joined),
internal/appengine/channels/ws.go:266:	prefix := "rooms:" + s.realm + ":"
internal/appengine/channels/ws.go:385:// pumpEvents ranges over member events and pushes new_event frames to the client.
internal/appengine/channels/protocol.go:13:	EventPhxJoin   = "phx_join"
internal/appengine/channels/protocol.go:18:	EventHeartbeat = "heartbeat"
internal/appengine/channels/protocol.go:20:	// TopicHeartbeat is the reserved topic for heartbeat messages.
internal/appengine/channels/room.go:66:	return &Registry{bus: bus, groups: groups, rooms: make(map[string]*Room)}

$ rg -n 'client_golang' go.mod
16:	github.com/prometheus/client_golang v1.23.2 // metrics (§5.2)

$ rg -rn 'openmetrics|OpenMetrics' --glob '!go.sum' . 2>/dev/null | head; rg -n 'openmetrics|OpenMetrics' internal/ docs/ cmd/ | head
(no output)

$ ls $(go env GOMODCACHE)/github.com/prometheus/ 2>/dev/null; echo "---"; rg -n 'EnableOpenMetrics' $(go env GOMODCACHE)/github.com/prometheus/client_golang@v1.23.2/prometheus/promhttp/http.go 2>/dev/null | head
client_golang@v1.12.1
client_golang@v1.23.2
client_model@v0.2.0
client_model@v0.6.2
common@v0.32.1
common@v0.66.1
procfs@v0.16.1
procfs@v0.7.3
---
194:		if opts.EnableOpenMetrics {
218:		if opts.EnableOpenMetricsTextCreatedSamples {
422:	EnableOpenMetrics bool
423:	// EnableOpenMetricsTextCreatedSamples specifies if this handler should add, extra, synthetic
437:	EnableOpenMetricsTextCreatedSamples bool

$ sed -n '185,215p' $(go env GOMODCACHE)/github.com/prometheus/client_golang@v1.23.2/prometheus/promhttp/http.go
return
				}
			case HTTPErrorOnError:
				httpError(rsp, err)
				return
			}
		}

		var contentType expfmt.Format
		if opts.EnableOpenMetrics {
			contentType = expfmt.NegotiateIncludingOpenMetrics(req.Header)
		} else {
			contentType = expfmt.Negotiate(req.Header)
		}
		rsp.Header().Set(contentTypeHeader, string(contentType))

		w, encodingHeader, closeWriter, err := negotiateEncodingWriter(req, rsp, compressions)
		if err != nil {
			if opts.ErrorLog != nil {
				opts.ErrorLog.Println("error getting writer", err)
			}
			w = io.Writer(rsp)
			encodingHeader = string(Identity)
		}

		defer closeWriter()

		// Set Content-Encoding only when data is compressed
		if encodingHeader != string(Identity) {
			rsp.Header().Set(contentEncodingHeader, encodingHeader)
		}

$ rg -n 'compressions' $(go env GOMODCACHE)/github.com/prometheus/client_golang@v1.23.2/prometheus/promhttp/http.go | head
145:	var compressions []string
152:			compressions = append(compressions, string(comp))
201:		w, encodingHeader, closeWriter, err := negotiateEncodingWriter(req, rsp, compressions)
390:	// OfferedCompressions is a set of encodings (compressions) handler will
463:// compressions. It returns a writer implementing the compression and an the
465:func negotiateEncodingWriter(r *http.Request, rw io.Writer, compressions []string) (_ io.Writer, encodingHeaderValue string, closeWriter func(), _ error) {
466:	if len(compressions) == 0 {
471:	selected := httputil.NegotiateContentEncoding(r, compressions)

$ rg -n -A20 '^func Negotiate' $(go env GOMODCACHE)/github.com/prometheus/common@v0.66.1/expfmt/encode.go; echo "=== FmtText ==="; rg -n 'FmtText|TypeTextPlain' $(go env GOMODCACHE)/github.com/prometheus/common@v0.66.1/expfmt/expfmt.go
63:func Negotiate(h http.Header) Format {
64-	escapingScheme := Format(fmt.Sprintf("; escaping=%s", Format(model.NameEscapingScheme.String())))
65-	for _, ac := range goautoneg.ParseAccept(h.Get(hdrAccept)) {
66-		if escapeParam := ac.Params[model.EscapingKey]; escapeParam != "" {
67-			switch Format(escapeParam) {
68-			case model.AllowUTF8, model.EscapeUnderscores, model.EscapeDots, model.EscapeValues:
69-				escapingScheme = Format("; escaping=" + escapeParam)
70-			default:
71-				// If the escaping parameter is unknown, ignore it.
72-			}
73-		}
74-		ver := ac.Params["version"]
75-		if ac.Type+"/"+ac.SubType == ProtoType && ac.Params["proto"] == ProtoProtocol {
76-			switch ac.Params["encoding"] {
77-			case "delimited":
78-				return FmtProtoDelim + escapingScheme
79-			case "text":
80-				return FmtProtoText + escapingScheme
81-			case "compact-text":
82-				return FmtProtoCompact + escapingScheme
83-			}
--
96:func NegotiateIncludingOpenMetrics(h http.Header) Format {
97-	escapingScheme := Format(fmt.Sprintf("; escaping=%s", Format(model.NameEscapingScheme.String())))
98-	for _, ac := range goautoneg.ParseAccept(h.Get(hdrAccept)) {
99-		if escapeParam := ac.Params[model.EscapingKey]; escapeParam != "" {
100-			switch Format(escapeParam) {
101-			case model.AllowUTF8, model.EscapeUnderscores, model.EscapeDots, model.EscapeValues:
102-				escapingScheme = Format("; escaping=" + escapeParam)
103-			default:
104-				// If the escaping parameter is unknown, ignore it.
105-			}
106-		}
107-		ver := ac.Params["version"]
108-		if ac.Type+"/"+ac.SubType == ProtoType && ac.Params["proto"] == ProtoProtocol {
109-			switch ac.Params["encoding"] {
110-			case "delimited":
111-				return FmtProtoDelim + escapingScheme
112-			case "text":
113-				return FmtProtoText + escapingScheme
114-			case "compact-text":
115-				return FmtProtoCompact + escapingScheme
116-			}
=== FmtText ===
50:	// Deprecated: Use expfmt.NewFormat(expfmt.TypeTextPlain) instead.
51:	FmtText Format = `text/plain; version=` + TextVersion + `; charset=utf-8`
82:	TypeTextPlain
98:	case TypeTextPlain:
99:		return FmtText
180:			return TypeTextPlain
183:			return TypeTextPlain

$ sed -n '63,95p' $(go env GOMODCACHE)/github.com/prometheus/common@v0.66.1/expfmt/encode.go; echo "=== compressions block ==="; sed -n '140,160p' $(go env GOMODCACHE)/github.com/prometheus/client_golang@v1.23.2/prometheus/promhttp/http.go
func Negotiate(h http.Header) Format {
	escapingScheme := Format(fmt.Sprintf("; escaping=%s", Format(model.NameEscapingScheme.String())))
	for _, ac := range goautoneg.ParseAccept(h.Get(hdrAccept)) {
		if escapeParam := ac.Params[model.EscapingKey]; escapeParam != "" {
			switch Format(escapeParam) {
			case model.AllowUTF8, model.EscapeUnderscores, model.EscapeDots, model.EscapeValues:
				escapingScheme = Format("; escaping=" + escapeParam)
			default:
				// If the escaping parameter is unknown, ignore it.
			}
		}
		ver := ac.Params["version"]
		if ac.Type+"/"+ac.SubType == ProtoType && ac.Params["proto"] == ProtoProtocol {
			switch ac.Params["encoding"] {
			case "delimited":
				return FmtProtoDelim + escapingScheme
			case "text":
				return FmtProtoText + escapingScheme
			case "compact-text":
				return FmtProtoCompact + escapingScheme
			}
		}
		if ac.Type == "text" && ac.SubType == "plain" && (ver == TextVersion || ver == "") {
			return FmtText + escapingScheme
		}
	}
	return FmtText + escapingScheme
}

// NegotiateIncludingOpenMetrics works like Negotiate but includes
// FmtOpenMetrics as an option for the result. Note that this function is
// temporary and will disappear once FmtOpenMetrics is fully supported and as
// such may be negotiated by the normal Negotiate function.
=== compressions block ===
			}
		}
	}

	// Select compression formats to offer based on default or user choice.
	var compressions []string
	if !opts.DisableCompression {
		offers := defaultCompressionFormats()
		if len(opts.OfferedCompressions) > 0 {
			offers = opts.OfferedCompressions
		}
		for _, comp := range offers {
			compressions = append(compressions, string(comp))
		}
	}

	h := http.HandlerFunc(func(rsp http.ResponseWriter, req *http.Request) {
		if !opts.ProcessStartTime.IsZero() {
			rsp.Header().Set(processStartTimeHeader, strconv.FormatInt(opts.ProcessStartTime.Unix(), 10))
		}
		if inFlightSem != nil {

$ rg -n -B3 -A12 'NameEscapingScheme' $(go env GOMODCACHE)/github.com/prometheus/common@v0.66.1/model/metric.go | head -40; echo "=== defaultCompressionFormats ==="; rg -n -A12 'func defaultCompressionFormats' $(go env GOMODCACHE)/github.com/prometheus/client_golang@v1.23.2/prometheus/promhttp/http.go
50-	// with this warning and learning for the future.
51-	NameValidationScheme = UTF8Validation
52-
53:	// NameEscapingScheme defines the default way that names will be escaped when
54-	// presented to systems that do not support UTF-8 names. If the Content-Type
55-	// "escaping" term is specified, that will override this value.
56:	// NameEscapingScheme should not be set to the NoEscaping value. That string
57-	// is used in content negotiation to indicate that a system supports UTF-8 and
58-	// has that feature enabled.
59:	NameEscapingScheme = UnderscoreEscaping
60-)
61-
62-// ValidationScheme is a Go enum for determining how metric and label names will
63-// be validated by this library.
64-type ValidationScheme int
65-
66-const (
67-	// UnsetValidation represents an undefined ValidationScheme.
68-	// Should not be used in practice.
69-	UnsetValidation ValidationScheme = iota
70-
71-	// LegacyValidation is a setting that requires that all metric and label names
--
235-	// metric and label names that do not conform to the legacy character
236-	// requirements should be escaped when being scraped by a legacy prometheus
237-	// system. If a system does not explicitly pass an escaping parameter in the
238:	// Accept header, the default NameEscapingScheme will be used.
239-	EscapingKey = "escaping"
240-
241-	// Possible values for Escaping Key.
242-	AllowUTF8         = "allow-utf-8" // No escaping required.
243-	EscapeUnderscores = "underscores"
244-	EscapeDots        = "dots"
245-	EscapeValues      = "values"
246-)
247-
248-// MetricNameRE is a regular expression matching valid metric
249-// names. Note that the IsValidMetricName function performs the same
250-// check but faster than a match with this regular expression.
=== defaultCompressionFormats ===
68:func defaultCompressionFormats() []Compression {
69-	if internal.NewZstdWriter != nil {
70-		return []Compression{Identity, Gzip, Zstd}
71-	} else {
72-		return []Compression{Identity, Gzip}
73-	}
74-}
75-
76-var gzipPool = sync.Pool{
77-	New: func() interface{} {
78-		return gzip.NewWriter(nil)
79-	},
80-}

$ mkdir -p /tmp/opencode/mprobe && cp /root/astrate-mule/go.mod /root/astrate-mule/go.sum /tmp/opencode/mprobe/ && sed -i '1s|.*|module mprobe|' /tmp/opencode/mprobe/go.mod && head -3 /tmp/opencode/mprobe/go.mod
module mprobe

go 1.26.1

← Write /tmp/opencode/mprobe/main.go
Wrote file successfully.

$ GOFLAGS=-mod=mod GOPROXY=off GOSUMDB=off go run . 2>&1 | tail -20
no-accept    status=200 content-type="text/plain; version=0.0.4; charset=utf-8; escaping=underscores" content-encoding="" bodylen=135
prom-default status=200 content-type="text/plain; version=0.0.4; charset=utf-8; escaping=underscores" content-encoding="" bodylen=135
protobuf     status=200 content-type="application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores" content-encoding="" bodylen=83
gzip         status=200 content-type="text/plain; version=0.0.4; charset=utf-8; escaping=underscores" content-encoding="gzip" bodylen=109
zstd         status=200 content-type="text/plain; version=0.0.4; charset=utf-8; escaping=underscores" content-encoding="" bodylen=135

→ Read internal/appengine/channels/protocol.go

$ sed -n '380,430p' internal/appengine/channels/ws.go
j.room.Unwatch(payload.Name)
	rep, _ := OK(f, nil)
	s.writeFrame(rep)
}

// pumpEvents ranges over member events and pushes new_event frames to the client.
func (s *session) pumpEvents(topic string) {
	s.roomsMu.Lock()
	j, ok := s.rooms[topic]
	if !ok {
		s.roomsMu.Unlock()
		return
	}
	events := j.member.Events()
	s.roomsMu.Unlock()

	for ev := range events {
		s.roomsMu.Lock()
		j, ok := s.rooms[topic]
		if !ok {
			s.roomsMu.Unlock()
			return
		}
		joinRef := j.joinRef
		s.roomsMu.Unlock()

		b, err := json.Marshal(ev)
		if err != nil {
			continue
		}
		s.writeFrame(Frame{
			JoinRef: joinRef,
			Ref:     nil,
			Topic:   topic,
			Event:   EventNewEvent,
			Payload: b,
		})
	}
}

func (s *session) sendErr(f Frame, reason string) {
	rep, _ := Err(f, reason)
	s.writeFrame(rep)
}

// writeFrame serialises and sends a frame under mu so that two goroutines
// never corrupt the WebSocket stream.
//
// The write is bounded by writeTimeout rather than running on a background
// context: a viewer that has stopped reading must not wedge this mutex, which
// the event pump also takes (docs/DESIGN.md §1.4 — a slow viewer drops frames,

$ rg -n -A30 'func (s \*Server) handle|func.*handle\(w http' internal/appengine/channels/ws.go | head -70
52:func (a *API) handle(w http.ResponseWriter, r *http.Request) {
53-	realm := r.URL.Query().Get("realm")
54-	token := r.URL.Query().Get("token")
55-
56-	if realm == "" || token == "" {
57-		_ = astarteapi.WriteUnauthorized(w)
58-		return
59-	}
60-
61-	row, err := a.keys.GetRealmByName(r.Context(), realm)
62-	switch {
63-	case errors.Is(err, store.ErrNotFound):
64-		_ = astarteapi.WriteUnauthorized(w)
65-		return
66-	case err != nil:
67-		_ = astarteapi.WriteInternalServerError(w)
68-		return
69-	}
70-
71-	tok, err := a.cache.Verify(token, row.JWTPublicKeysPEM)
72-	if err != nil {
73-		_ = astarteapi.WriteUnauthorized(w)
74-		return
75-	}
76-
77-	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
78-	if err != nil {
79-		return
80-	}
81-
82-	s := &session{

$ rg -n -B6 -A25 'type Event struct' internal/engine/stream/*.go
internal/engine/stream/bus.go-34-	// "failures are never silent"). Unlike the trigger path, it is published
internal/engine/stream/bus.go-35-	// whether or not the realm declares any trigger.
internal/engine/stream/bus.go-36-	KindDeviceError = "device_error"
internal/engine/stream/bus.go-37-)
internal/engine/stream/bus.go-38-
internal/engine/stream/bus.go-39-// Event is one live event.
internal/engine/stream/bus.go:40:type Event struct {
internal/engine/stream/bus.go-41-	// Kind discriminates the event (Kind* constants).
internal/engine/stream/bus.go-42-	Kind string
internal/engine/stream/bus.go-43-	// Realm is the tenant.
internal/engine/stream/bus.go-44-	Realm string
internal/engine/stream/bus.go-45-	// DeviceID is the encoded device ID.
internal/engine/stream/bus.go-46-	DeviceID string
internal/engine/stream/bus.go-47-	// Interface and Path locate data events; empty for lifecycle events.
internal/engine/stream/bus.go-48-	Interface string
internal/engine/stream/bus.go-49-	Path      string
internal/engine/stream/bus.go-50-	// InterfaceMajor is the data event's interface major version. Consumers
internal/engine/stream/bus.go-51-	// that rebuild upstream trigger-event bodies need it; Filter ignores it.
internal/engine/stream/bus.go-52-	InterfaceMajor int
internal/engine/stream/bus.go-53-	// IP is the peer address of a device_connected event, empty otherwise.
internal/engine/stream/bus.go-54-	IP string
internal/engine/stream/bus.go-55-	// ErrorName and ErrorMetadata describe a device_error event: the §2.6
internal/engine/stream/bus.go-56-	// reject-reason label and its free-form diagnostic strings.
internal/engine/stream/bus.go-57-	ErrorName     string
internal/engine/stream/bus.go-58-	ErrorMetadata map[string]string
internal/engine/stream/bus.go-59-	// Value is the JSON-friendly rendering of a data event's value (nil for
internal/engine/stream/bus.go-60-	// property unset and lifecycle events).
internal/engine/stream/bus.go-61-	Value any
internal/engine/stream/bus.go-62-	// Timestamp is the event instant (the effective sample timestamp for
internal/engine/stream/bus.go-63-	// data events).
internal/engine/stream/bus.go-64-	Timestamp time.Time
internal/engine/stream/bus.go-65-}

$ rg -n 'Kind|DeviceID|event_data|InterfaceMajor' internal/appengine/channels/*_test.go | head -20
internal/appengine/channels/room_test.go:39:	validDeviceIDA = "f0VMRgIBAQAAAAAAAAAAAA"
internal/appengine/channels/room_test.go:40:	validDeviceIDB = "f0VMRgIBAQAAAAAAAAAAAQ"
internal/appengine/channels/room_test.go:50:		DeviceID:      validDeviceIDA,
internal/appengine/channels/room_test.go:51:		SimpleTrigger: json.RawMessage(`{"type":"device_trigger","on":"device_connected","device_id":"` + validDeviceIDA + `"}`),
internal/appengine/channels/room_test.go:55:		DeviceID:      validDeviceIDA,
internal/appengine/channels/room_test.go:56:		SimpleTrigger: json.RawMessage(`{"type":"device_trigger","on":"device_disconnected","device_id":"` + validDeviceIDA + `"}`),
internal/appengine/channels/room_test.go:60:		DeviceID:      validDeviceIDA,
internal/appengine/channels/room_test.go:61:		SimpleTrigger: json.RawMessage(`{"type":"device_trigger","on":"device_error","device_id":"` + validDeviceIDA + `"}`),
internal/appengine/channels/room_test.go:65:		DeviceID:      validDeviceIDA,
internal/appengine/channels/room_test.go:84:	bus.ch <- stream.Event{Kind: stream.KindDeviceConnected, DeviceID: validDeviceIDA, IP: "1.2.3.4", Timestamp: now}
internal/appengine/channels/room_test.go:85:	bus.ch <- stream.Event{Kind: stream.KindDeviceDisconnected, DeviceID: validDeviceIDA, Timestamp: now}
internal/appengine/channels/room_test.go:86:	bus.ch <- stream.Event{Kind: stream.KindDeviceError, DeviceID: validDeviceIDA, ErrorName: "unknown_format", ErrorMetadata: map[string]string{"detail": "slow"}, Timestamp: now}
internal/appengine/channels/room_test.go:87:	bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: validDeviceIDA, Interface: "org.example.V1", Path: "/sensor/temp", Value: 42.0, InterfaceMajor: 1, Timestamp: now}
internal/appengine/channels/room_test.go:150:func TestDeviceIDFilterIsAdditive(t *testing.T) {
internal/appengine/channels/room_test.go:157:		DeviceID:      validDeviceIDA,
internal/appengine/channels/room_test.go:168:	bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: validDeviceIDB, Interface: "org.example.V1", Path: "/sensor/temp", Value: 42.0, InterfaceMajor: 1, Timestamp: now}
internal/appengine/channels/room_test.go:180:		"probe": {validDeviceIDA},
internal/appengine/channels/room_test.go:207:	bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: validDeviceIDB, Interface: "org.example.V1", Path: "/sensor/temp", Value: 1.0, InterfaceMajor: 1, Timestamp: now}
internal/appengine/channels/room_test.go:214:	bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: validDeviceIDA, Interface: "org.example.V1", Path: "/sensor/temp", Value: 2.0, InterfaceMajor: 1, Timestamp: now}
internal/appengine/channels/room_test.go:217:		if ev.DeviceID != validDeviceIDA {

← Write /tmp/opencode/mprobe/main.go
Wrote file successfully.

$ GOFLAGS=-mod=mod GOPROXY=off GOSUMDB=off go run . 2>&1 | tail
new_event payload: {"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}
spec payload unmarshals into event_data string? true
real payload unmarshals into event_data string? true

$ rg -n 'version=0.0.4|text/plain|metrics|Content-Type|readiness|health' docs/site/observability.md | head -40
3:Astrate exposes Prometheus metrics, structured JSON logging, and health/readiness probes.
9:| `GET /astrate/v1/health` | Liveness probe. Returns 200 if the process is running. |
10:| `GET /astrate/v1/readiness` | Readiness probe. Pings the database and checks the broker listener. Returns 503 with per-check status if any dependency is down. Use for load-balancer and orchestrator probes. |
11:| `GET /astrate/v1/metrics` | Prometheus scrape endpoint. |
15:## Prometheus metrics
17:Exposed at `/astrate/v1/metrics`. All counters and histograms use the `astrate_` namespace.
86:- [Deployment](deployment.md) — health probe usage in Docker Compose and orchestrators
88:- [Troubleshooting](troubleshooting.md) — debugging with health/readiness/metrics endpoints

$ rg -n 'astrate/v1|appengine/v1/\{realm\}/version|/version' docs/site/observability.md docs/site/api.md 2>/dev/null | head -30; echo "=== files with compat health ==="; rg -ln 'appengine/health' docs/site/
docs/site/observability.md:9:| `GET /astrate/v1/health` | Liveness probe. Returns 200 if the process is running. |
docs/site/observability.md:10:| `GET /astrate/v1/readiness` | Readiness probe. Pings the database and checks the broker listener. Returns 503 with per-check status if any dependency is down. Use for load-balancer and orchestrator probes. |
docs/site/observability.md:11:| `GET /astrate/v1/metrics` | Prometheus scrape endpoint. |
docs/site/observability.md:17:Exposed at `/astrate/v1/metrics`. All counters and histograms use the `astrate_` namespace.
=== files with compat health ===

$ rg -n 'docs-sync-native' .mule/todo.md | sed 's/\[auto\].*//'
154:- [x] docs-sync-native-compat-health-503 
157:- [x] docs-sync-native-compat-version-endpoints 
185:- [x] docs-sync-native-socket-query-token-auth 
196:- [x] docs-sync-native-version-value 
201:- [x] docs-sync-native-metrics-example-fake-series 
202:- [x] docs-sync-native-socket-missing-403-500 
242:- [x] docs-sync-pairing-info-version-example 
262:- [x] docs-sync-native-error-example-capitalisation

$ rg -in 'event_data|new_event|text/plain|escaping=underscores|realm-name-pattern|realmname' .mule/todo.md | head
198:- [x] channels-rejoin-joinref-tagging-test: pin in internal/appengine/channels/ws_test.go that a second `phx_join` on the same topic updates the session's `joinRef` (ws.go:284-291) and subsequent `new_event` pushes carry the fresh ref (read per event at ws.go:403) — no test drives a rejoin today (all `TestWireSession_*` join once), so the comment's "server pushes must be tagged with the current one" rots unasserted. Wire: join with ref "1", rejoin with "2", push an event, assert element 0 is `"2"`.
226:- [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
268:- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
272:- [x] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `components.parameters.RealmName` and `RealmCreate.realm_name` both carry `pattern: '^[a-z][a-z0-9]*$'` and `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go) holds both against the realm name CHECK extracted from migrations/000002_metadata.up.sql — but the read side is still bare `type: string`: `Realm.realm_name` (components.schemas.Realm) and `DataEnvelopeRealmNames.items`. Decide whether a response field earns the same pattern (a generated client learns nothing new from it, but a reader diffing create against get sees one constrained and one not) and, if it does, add it to both and extend that test to hold the response-side pattern against the same CHECK. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ sed -n '154p;185p;201p;202p' .mule/todo.md | cut -c1-1400
- [x] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmmanagement/health`, `GET /pairing/health` — in docs/api/astrate_native_api.yaml, and widen the shared `DataEnvelopeHealthStatus.status` enum (yaml:301-303) from `[ok]` to `[ok, unhealthy]` — the handler answers 503 `{"data":{"status":"unhealthy"}}` when the readiness probe fails (internal/observability/compat.go:19-34, deviation 18), but the spec documents only 200 with `status: ok` for all three paths. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication is "via Authorization header or `?token=` query parameter" (yaml:293-294), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate_devices_total{realm="test"} 42`, but no such series is ever registered — the registry carries `astrate_broker_sessions` (internal/observability/metrics.go:54-62), `astrate_db_pool_{acquired,idle,total,max}_conns` (metrics.go:65-79), plus the engine/flow/trigger families (`astrate_engine_*` engine/router.go:411-444, `astrate_engine_stream_dropped_total` engine/stream/bus.go:119, `astrate_engine_trigger_*` engine/triggers/actions.go:296-300, `astrate_flow_router_*` flow/router.go:260-276) and the standard `go_*`/`process_*` runtime collectors — so the example advertises a scrape result the wire can never return. Swap it for a real series (e.g. `astrate_broker_sessions 0`), keeping the HELP/TYPE comment style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ rg -n 'securitySchemes|^  [A-Za-z]+:$|^    [A-Za-z]+:$' docs/api/astarte_appengine_api.yaml | sed -n '1,60p'
10:  license:
16:    variables:
27:    get:
101:    get:
141:    get:
167:    patch:
219:    get:
250:    patch:
307:    get:
341:    get:
388:    get:
439:    put:
491:    post:
543:    delete:
585:    get:
615:    get:
657:    get:
701:    put:
748:    post:
795:    delete:
834:    get:
861:    post:
912:    get:
942:    get:
991:    post:
1038:    get:
1067:    patch:
1119:    delete:
1143:    get:
1173:    get:
1216:    get:
1261:    put:
1309:    post:
1357:    delete:
1395:  securitySchemes:
1404:  parameters:
1405:    RealmName:
1413:    DeviceID:
1421:    InterfaceName:
1429:    GroupName:
1437:    DataSince:
1445:    DataSinceAfter:
1453:    DataTo:
1461:    DataLimit:
1468:    DataDownsample:
1475:    DataSort:
1484:    DataFormat:
1497:    DataAllowBigIntegers:
1507:    DataAllowSafeBigIntegers:
1518:  schemas:
1521:    DeviceStatus:
1580:    Introspection:
1591:    DevicePatch:
1611:    GroupCreate:
1624:    GroupInfo:
1634:    DataEnvelopeGroupInfo:
1641:    DataEnvelopeDeviceStatus:
1648:    DataEnvelopeDeviceIDList:
1659:    DataEnvelopeDeviceStatusList:
1670:    DataEnvelopeObject:

→ Read docs/api/astarte_appengine_api.yaml [limit=28, offset=1393]

$ rg -n -A3 '^  responses:' docs/api/astarte_appengine_api.yaml; rg -n '^    [A-Z][A-Za-z]+:$' docs/api/astarte_appengine_api.yaml | sed -n '/1[5-9][0-9][0-9]:/,$p' | head -30
1723:  responses:
1724-    BadRequest:
1725-      description: Bad Request — malformed or missing data envelope.
1726-      content:
1507:    DataAllowSafeBigIntegers:
1521:    DeviceStatus:
1580:    Introspection:
1591:    DevicePatch:
1611:    GroupCreate:
1624:    GroupInfo:
1634:    DataEnvelopeGroupInfo:
1641:    DataEnvelopeDeviceStatus:
1648:    DataEnvelopeDeviceIDList:
1659:    DataEnvelopeDeviceStatusList:
1670:    DataEnvelopeObject:
1678:    DataEnvelopeStringArray:
1687:    Links:
1700:    ErrorDetail:
1711:    FieldErrorsDetail:
1724:    BadRequest:
1734:    Unauthorized:
1744:    Forbidden:
1754:    NotFound:
1764:    DeviceNotFound:
1774:    MethodNotAllowed:
1784:    InternalServerError:
1794:    ValidationErrors:
1805:    ValueTooLarge:
1815:    UnprocessableEntity:
1825:    ConflictAliasInUse:
1835:    ConflictGroupAlreadyExists:
1845:    ConflictDeviceAlreadyInGroup:

$ rg -n 'json.Marshal\(ev\)|WriteFrame|writeFrame\(Frame\{|Payload: b,' internal/appengine/channels/ws.go; echo "=== phoenix auth lines ==="; rg -n 'realm == "" \|\| token == ""|GetRealmByName|cache.Verify' internal/appengine/channels/ws.go; echo "=== migrations check ==="; sed -n '1,12p' migrations/000002_metadata.up.sql
406:		b, err := json.Marshal(ev)
410:		s.writeFrame(Frame{
415:			Payload: b,
=== phoenix auth lines ===
56:	if realm == "" || token == "" {
61:	row, err := a.keys.GetRealmByName(r.Context(), realm)
71:	tok, err := a.cache.Verify(token, row.JWTPublicKeysPEM)
=== migrations check ===
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

$ rg -n 'version = "0.1.0-dev"|handler := httpx.NotFound|mountRealmVersion\(mux|MountVersionCompat\(mux, svc, version\)' cmd/astrate/main.go
53:var version = "0.1.0-dev"
397:		observability.MountVersionCompat(mux, svc, version)
399:	mountRealmVersion(mux, mw)
422:	handler := httpx.NotFound(mux)
435:func mountRealmVersion(mux *http.ServeMux, mw *auth.Middleware) {

Line numbers verified. Appending the five proposals.
$ tail -c 400 .mule/todo.md; echo; echo "=== bytes ==="; wc -c .mule/todo.md
rnotfound: Add store-level test asserting LatestIndividual returns ErrNotFound for empty series and returns newest on multiple samples. File: internal/store/datastreams_test.go — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [x] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.go

=== bytes ===
142187 .mule/todo.md

✗ Read .mule/todo.md failed [offset=296]
Error: Offset 296 is out of range for this file (276 lines)
→ Read .mule/todo.md [offset=273]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -273,4 +273,9 @@
 - [!] store-register-inhibit-preserve: RegisterDevice must not overwrite status to 'registered' when a device was inhibited before first credentials request (preserve 'inhibited' state). File: internal/store/devices.go — BLOCKED: gates failed
 - [x] store-aliasvalues-self-exclusion-test: Add test for AliasValuesTaken to assert self-exclusion (same device's own alias values do not count as taken). File: internal/store/devices_test.go
 - [!] store-latestindividual-empty-errnotfound: Add store-level test asserting LatestIndividual returns ErrNotFound for empty series and returns newest on multiple samples. File: internal/store/datastreams_test.go — BLOCKED: TIMEOUT after 1201s — task too big, split it
 - [x] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.go
+- [ ] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ rg -c '^- \[ \] docs-native' .mule/todo.md; echo "--- new lines only, unescaped view ---"; rg -n '^- \[ \] docs-native-(ae-realm|hk-version|metrics-content|phoenix-newevent|realm-name)' .mule/todo.md | cut -c1-120
5
--- new lines only, unescaped view ---
277:- [ ] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs
278:- [ ] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_nat
279:- [ ] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native
280:- [ ] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.
281:- [ ] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api

$ git status --porcelain; echo "--- diff stat ---"; git diff --stat
M .mule/todo.md
--- diff stat ---
 .mule/todo.md | 5 +++++
 1 file changed, 5 insertions(+)

# Todos
[✓] Pick one surface and list handler routes vs documented paths (native)
[✓] Verify status codes / response fields for 3-4 endpoints against handlers
[✓] Check ASTRATE_ config keys vs docs/site/configuration-reference.md
[✓] Append up to 5 task lines to .mule/todo.md (+ for-giulio.md if prose contradicts code)

Done: docs-sync recipe run on surface **astrate-native** — five `docs/api/` fix tasks appended to `.mule/todo.md` (no code touched, no git, no `gh`; the recipe's deliverable here is task lines, not an issue).

Files: `.mule/todo.md` (+5 lines, 277-281)

Verified: `rg 'mux\.(Handle|HandleFunc)' internal/ cmd/` vs `rg '^\s{2}/' docs/api/astrate_native_api.yaml`; read every handler behind the 12 documented native paths; two throwaway probes in `/tmp/opencode/mprobe` against the repo's own go.mod/go.sum (`promhttp.HandlerFor(reg, HandlerOpts{})` content negotiation → `text/plain; version=0.0.4; charset=utf-8; escaping=underscores`, protobuf + gzip paths real, zstd absent; verbatim copy of `stream.Event` marshalled → capitalised Go field names, no `event_data`). Every `file:line` cited in the task lines re-checked against the tree.

Unsure: two judgement calls I left to the executor rather than deciding — (1) whether the undocumented `/appengine/v1/{realm}/version` belongs in `astarte_appengine_api.yaml` (I said so) or the native spec; (2) whether to document the protobuf/gzip metrics paths or just correct the content-type string. Also: the config-key sweep found **no new** drift — the two `ASTRATE_HOUSEKEEPING_*` keys missing from `configuration-reference.md` are already escalated three times in `for-giulio.md` (lines 126, 195, 265), so I did not re-report them and wrote nothing to `for-giulio.md`.

Follow-ups: five lines appended — `docs-native-ae-realm-version-endpoint`, `docs-native-hk-version-description`, `docs-native-metrics-content-negotiation`, `docs-native-phoenix-newevent-payload`, `docs-native-realm-name-pattern`. The gate (`make -C docs build`) was **not** run: nothing under `docs/` changed in this run.
