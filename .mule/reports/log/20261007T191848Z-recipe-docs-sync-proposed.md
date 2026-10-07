slug: recipe-docs-sync
verdict: proposed
at:  ebd5fb0
ran: 2026-10-07T19:18:48Z on DietPi in 518s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls docs/api/ && ls .mule/ && ls .mule/tasks 2>/dev/null | head -30
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
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
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md

$ wc -l .mule/todo.md .mule/for-giulio.md 2>/dev/null; ls docs/site/ | head -50
290 .mule/todo.md
   649 .mule/for-giulio.md
   939 total
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

$ grep -n -i 'docs-sync\|docs/api\|astarte_.*_api.yaml\|astrate_native_api' .mule/todo.md
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_117c66c7d001Au5PFHUus3x2G8

217:- [x] docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` (yaml:231-232) in docs/api/astarte_realm_management_api.yaml as a oneOf — both handlers go through `writeInterfaceError`, which answers three distinct 422 bodies: the flat ErrorDetail (`ErrValidation` via `validationDetail`, http.go:392-393), the nested violations changeset envelope (`writeViolations`, http.go:358-359 + 434-533, e.g. `{"errors":{"description":["should be at most 1000 character(s)"]}}` and the aligned full-length `mappings` array, http_test.go:480-501), and the named FieldErrors envelope for `ErrMaximumDatabaseRetentionExceeded` (`{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, http.go:355-357, http_test.go:282-320); today only the flat ValidationError is referenced. Mirror the createTrigger oneOf pattern (yaml:344-368). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
218:- [x] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "realm: validation failed: interface definition is invalid"`, but `validationDetail` strips the `realm: validation failed: ` prefix on the wire (http.go:562-568) and the parser's real messages are `invalid interface: ...` (pkg/interfaceschema/parse.go:38, wrapped at service.go:163); the example should show the stripped, real message. Same class as the already-fixed docs-sync-hk-validation-example (housekeeping http.go:251-258). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
219:- [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
226:- [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
227:- [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
228:- [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
229:- [x] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion Go — the device-delete leg added by docs-sync-rm-delete-device-async-operation-param (register a device, `DELETE /devices/{id}?async_operation=…` answers 204 and the row is gone on both values) compiles clean under `go vet -tags integration ./internal/realm/` but has never been executed, because the Pi has no database; until it runs, TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go is the only assertion of that leg that has actually been verified.
230:- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
236:- [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml — the string `"403"` appears nowhere in the file and `components` defines no `Forbidden` response, yet every route is wrapped by `a.require = mw.RequireRealm(auth.ClaimAppEngine)` (internal/appengine/http.go:35), which answers `403 Forbidden` for a verified a_aea JWT whose grants do not authorize the method+path (`astarteapi.WriteForbidden`, internal/auth/middleware.go:90; `DetailForbidden`, pkg/astarteapi/envelope.go:41; pinned today by internal/appengine/http_test.go:184-185). Add a `Forbidden` response component copied from the realm-management one (docs/api/astarte_realm_management_api.yaml:1277-1285 — `ErrorDetail` schema, example `errors.detail: Forbidden`) and `$ref` it from all 32 operations, which is exactly the parity docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go:123-170) pins; add the same docs-half test for the AppEngine spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
237:- [x] docs-sync-ae-patch-by-alias-409 [auto]: add the missing `409` response to `patchDeviceByAlias` in docs/api/astarte_appengine_api.yaml (yaml:241-291, currently 200/400/401/404/422/500) — the handler resolves the alias then calls the same `applyPatch` as the device- and group-scoped PATCH, and a rename to an alias already in use returns `ErrAliasAlreadyInUse` (internal/appengine/service.go:337) which `writeError` maps to `409 ConflictAliasInUse` (internal/appengine/http.go:650-651); the other two PATCH operations already document that response thanks to docs-sync-appengine-query-params-status and docs-sync-appengine-group-patch-status, so this is the last PATCH of the three with a hole. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
238:- [x] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml (yaml:939-980, currently 201/400/401/404/409/500) — a well-formed body whose `device_id` does not parse is rejected as a field error (`FieldErrors{"device_id": {"is not a valid device id"}}`, internal/appengine/service.go:631-634) and answered 422 with the `FieldErrorsDetail` body (`astarteapi.WriteFieldErrors`, internal/appengine/http.go:647), so `$ref` the existing `ValidationErrors` response (yaml:1707-1716) rather than `BadRequest` — the documented 400 only covers a body `DecodeData` cannot parse (internal/appengine/http.go:488-491). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
239:- [x] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data DELETE operations in docs/api/astarte_appengine_api.yaml (`deleteDataByAlias` yaml:519-547, `deleteData` yaml:756-784, `deleteDataInGroup` yaml:1285-1313 — each currently 204/401/404/405/500) — unsetting a path that matches no endpoint mapping returns `engine.ErrPathNotFound` (`UnsetServerProperty`, internal/engine/serverdata.go:241-244) which `writeError` answers as `400 "Endpoint not found"` (internal/appengine/http.go:679-680), the same 400 the PUT/POST twins already document via `BadRequest`; internal/appengine/writeerror_test.go pins the mapping. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
240:- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 200/401/500) — a `from_token` that does not parse as a device id is refused with `ErrValidation: invalid cursor` (internal/appengine/service.go:150-155) which `writeError` maps to 422 with the flat `{"errors":{"detail":"invalid cursor"}}` body, i.e. the `ErrorDetail` shape and not `FieldErrorsDetail`, so `$ref` the same response the group twin uses for the identical parameter (`ValidationErrors`, yaml:934-935) only after fixing its example, or add a dedicated one. Extend the `from_token` parameter description (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
241:- [x] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:556) and `InternalServerError` shows `detail: Internal Server Error` (yaml:638), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37, golden fixture pkg/astarteapi/testdata/error_bad_request.json) and `"Internal server error"` (envelope.go:58, testdata/error_internal.json), both measured on upstream 1.2.0 (envelope.go:33-36, 55-57). The other five detail examples on this file (`Not Found`, `Device not found`, `Unauthorized`, `Forbidden`, `Too Many Requests`) are already correct — change only those two. The identical wrong pair is in the other four specs (appengine 1732/1792, housekeeping 362/385/477, realm_management 1265/1348, native 466); leave those for a run of their surface. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
242:- [x] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pairing_api.yaml, which shows `version: "0.1.0-test"` (yaml:252) — a value no build can produce. The field is filled from `pairing.Config.Version` (internal/pairing/service.go:302), wired from `main.version` (cmd/astrate/main.go:151), whose default is `"0.1.0-dev"` (main.go:53) and which is otherwise set only via `-ldflags "-X main.version=..."`. Use `"0.1.0-dev"`, matching the four compat-version examples docs-sync-native-version-value already corrected to that (docs/api/astrate_native_api.yaml:203, 223, 243, 263). Do NOT document the service's own fallback `DefaultVersion = "0.1.0-astrate"` (service.go:86-88): it is unreachable from `main`, which always passes a non-empty value. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
243:- [x] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` path parameter (docs/api/astarte_pairing_api.yaml:360) and `RegisterRequest.hw_id` (yaml:373) — but `deviceid.Parse` accepts only the 22-character **unpadded base64url** form, in strict mode, and rejects the standard alphabet's `+` and `/` and any `=` padding with an error wrapping `ErrInvalid` (pkg/deviceid/deviceid.go:26, 37, 39-60; upstream parity cited at 36). Measured: with the same `encoding` the parser uses, `dT6hS2W9TT6LEnP25ks+lg` and `.../lg` fail at byte 19 and the padded 23-character form fails at byte 21, so a client generated from the current spec can send a well-formed standard-base64 ID and get 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}` (internal/pairing/http.go:307-309). State the exact encoding and the 22-character length in both descriptions and add a `pattern`/`minLength`/`maxLength` to `hw_id`. The existing `hw_id: dT6hS2W9TT6LEnP25ks_lg` example (yaml:123) is valid and stays. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
244:- [x] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in docs/api/astarte_pairing_api.yaml (yaml:374-376), an unconstrained string described only as "Optional Astrate extension specifying the device's initial payload format" — so a generated client can send any string, while `Service.Register` accepts only `"bson"` or `"json"` and rejects everything else with 422 `{"errors":{"initial_payload_format":["is invalid"]}}` (internal/pairing/service.go:179-181 -> internal/pairing/http.go:310-312; the allowed set is stated at service.go:173). Add `enum: [bson, json]` plus a note that absent/empty means unset, and add the three missing `ValidationError` examples while you are in the component: `initial_payload_format: ["is invalid"]` (http.go:310-312), `csr: ["can't be blank"]` (http.go:173-177) and `client_crt: ["can't be blank"]` (http.go:253-257) — the component currently shows only the four hw_id/error_name rows (yaml:594-618) and the `csr`/`client_crt` blank-field cases are what a device actually hits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
245:- [x] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` description in docs/api/astarte_pairing_api.yaml (yaml:153), "Removes a device from the realm. Returns 204 on success.", which reads as a device-and-data deletion. The row and every byte of its stored data survive: `store.UnregisterDevice` only clears the credential trail (`credentials_secret_hash = ''`, `first_credentials_request = NULL`, `cert_serial = NULL`, `cert_aki = NULL`) and flips `status` back to `'registered'` (internal/store/devices.go:93-110), and the service docstring says so outright ("its data is kept", internal/pairing/service.go:227-229). State the real semantics: the device becomes registrable again, its interfaces, datastream data and group memberships are retained, and a second DELETE is a 204 no-op (the row is still there, so RowsAffected is 1, devices.go:106-108) — not the documented 404. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
246:- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `base64.RawURLEncoding.Strict()` (pkg/deviceid/deviceid.go:37), so of the 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with `ErrInvalid` while `..._lA` is accepted, and the pairing spec's `hw_id` pattern `^[A-Za-z0-9_-]{22}$` (added by docs-sync-pairing-deviceid-base64url) deliberately does not encode the rule because a pattern is a client-side contract and upstream's leniency is unmeasured. The docstring justifies strictness as byte-identical round-tripping, which is a reason, not a parity claim, and docs/COMPATIBILITY.md records no entry for it. Measure upstream (Elixir `Base.url_decode64!(padding: false)`, which does not check those bits) on the Legion Go, then either record the divergence in docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.
254:- [x] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, housekeeping and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
255:- [!] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy alias fields to `InterfaceDefinition` and `InterfaceMapping` in docs/api/astarte_realm_management_api.yaml — both declare `additionalProperties: false` (yaml:935, 972) yet omit interface-level `quality` (alias of `ownership`) and `aggregate` (alias of `aggregation`) and mapping-level `path` (alias of `endpoint`), all three of which the parser accepts and canonicalises: decoded into dedicated fields (pkg/interfaceschema/parse.go:116-122, 136-138), normalised before any rule runs (parse.go:212-236), and re-encoded canonically by `ParseInterfaceCanonical` so the store never sees them (parse.go:165-170, 307-311). A spec-conformant client — or any generated client with this schema baked in — therefore rejects a valid upstream-shaped install/update with a spurious 422 while the service would have accepted and stored it as `ownership`/`aggregation`/`endpoint`. Add the three as documented deprecated aliases (enums mirroring ownership/aggregation, `path` a plain string) or relax `additionalProperties`, and record the three rules the code enforces: `ownership` + `quality` is a violation (parse.go:218) and `aggregation` + `aggregate` is a violation (parse.go:227), while `endpoint` + `path` silently keeps `endpoint` (parse.go:307-311). None of the three may appear in a `required` list. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
256:- [x] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
257:- [x] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstrained string described only as "The device hardware ID." — `Service.DeleteDevice` runs the path value through `deviceid.Parse` (internal/realm/service.go:118-121), which accepts exactly 22 characters of unpadded base64url in strict mode and rejects padding, the standard alphabet's `+` and `/`, and non-canonical trailing bits (pkg/deviceid/deviceid.go:24-26, 33-37, 42-60); every rejection is folded into `store.ErrNotFound` (service.go:119-121) and answered `404 {"errors":{"detail":"Device not found"}}` (internal/realm/http.go:107-109), so a client that sends a well-formed standard-base64 or 23-character padded id is told the device does not exist rather than that its id is malformed. State the encoding and the 22-character length and add `pattern`/`minLength`/`maxLength`, exactly as docs-sync-pairing-deviceid-base64url did for the pairing twin. Do NOT encode the strict trailing-bits rule in the pattern: whether it is a deviation worth recording is still open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
258:- [!] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_realm_management_api.yaml — `data` is declared as a bare `type: object` with no properties (yaml:250-261) although the handler parses it with the same `interfaceschema.ParseInterfaceCanonical` as `installInterface` (internal/realm/service.go:196 vs 161) and then checks the body's `interface_name`/`version_major` against the URL's `{name}`/`{major}` (service.go:203-208), so the update body — the one carrying every interface rule — is the only interface body in the spec with nothing to validate against. `$ref` `#/components/schemas/InterfaceDefinition` as `installInterface` does (yaml:115-117), and say in the operation description that the body must name the same interface and major as the URL (disagreement is 409, already documented) and that `version_minor` must increase (409). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1201s — task too big, split it
260:- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` shows `detail: Bad Request` (yaml:1732) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1792), while the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58 via `WriteBadRequest`/`WriteInternalServerError`, golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json); every 400/500 on that surface goes through the same astarteapi writers. The other detail examples on the file are already correct — change only those two — and add the docs-half assertion `TestAppEngineErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go, which pins the pairing and realm-management halves. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
261:- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDetail` schema-level example (yaml:361-362) and the `BadRequest` response (yaml:385) show `detail: Bad Request`, and `InternalServerError` (yaml:477) shows `detail: Internal Server Error`, but the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58, via WriteBadRequest/WriteInternalServerError, pinned by pkg/astarteapi/testdata/error_bad_request.json and error_internal.json). While you are in `ErrorDetail`: its `example` sits at the schema level and shows a bare `{detail: ...}` where the envelope the wire sends is `{"errors":{"detail":...}}` — move it under `errors` (or drop it, as the realm-management and pairing `ErrorDetail` schemas carry no example) so a generated client is not shown the wrong shape. Add `TestHousekeepingErrorDetailExamples` beside the pairing and realm-management halves in internal/swagger/swagger_test.go. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
262:- [x] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `"Internal server error"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
263:- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
268:- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
269:- [x] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362) — the second half docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail: Bad Request`, the reconstructed title-case form that `TestHousekeepingErrorDetailExamples` (internal/swagger/swagger_test.go:687) now pins against everywhere else on this spec, while the wire emits the frozen `astarteapi.DetailBadRequest` = `"Bad request"` (pkg/astarteapi/envelope.go:37, via `WriteBadRequest` on every 400 — internal/housekeeping/http.go:65, 166, 171, 183). Because the example sits on the generic `ErrorDetail` schema, which all eight error responses `$ref` (yaml:377-477), it is the example a generator falls back to for any response whose own example it cannot resolve — so it teaches `"Bad Request"` for a 404 and a 500 too. Change it to `Bad request`, or drop it as the realm-management and pairing `ErrorDetail` schemas do (neither carries one). Then extend `TestHousekeepingErrorDetailExamples`, whose comment currently says the schema-level example is "deliberately not pinned" (swagger_test.go:684-686), to pin it like the responses. One thing not to "fix": the example is nested one level in — it is the `errors` sub-schema's example, not a bare envelope — so do not re-wrap it in another `errors:`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
270:- [!] docs-sync-hk-patch-422-field-error-examples [auto]: add the three missing PATCH 422 bodies to the `ValidationErrors` response component in docs/api/astarte_housekeeping_api.yaml (yaml:427-436, referenced once at yaml:190) — it carries a single example, `error_name: [invalid_update_parameters]` (internal/housekeeping/http.go:176-177), but `patchRealm` emits three more field-error envelopes: `jwt_public_key_pem: ["can't be blank"]` (http.go:188-189), `device_registration_limit: ["is invalid"]` (http.go:193-194) and `datastream_maximum_storage_retention: ["is invalid"]` (http.go:198-199). Add them as further `examples` on the media type (OpenAPI 3.0.3, yaml:1) or split the component per case. Also record that the *messages* differ from POST for the same conditions: create rejects a negative limit/retention with the flat `ErrValidation` detail `device_registration_limit must be non-negative` / `datastream_maximum_storage_retention must be non-negative` (internal/housekeeping/service.go:143-148, 422 `ErrorDetail` shape), while PATCH answers `is invalid` in the FieldErrors shape — and that the service's own three `ErrValidation` branches (service.go:205-213) are unreachable from REST, because `patchRealm` pre-checks the identical conditions and answers first. Say in the `RealmPatch` field descriptions that a negative value is a 422 rather than silently ignored. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
271:- [!] docs-sync-hk-wrong-type-field-400 [auto]: widen the `BadRequest` component description in docs/api/astarte_housekeeping_api.yaml (yaml:378, "malformed or missing data envelope") — measured today, a perfectly well-formed `data` envelope carrying a wrong-typed field answers **400**, not the documented 422. `astarteapi.DecodeData` into `map[string]json.RawMessage` (internal/housekeeping/http.go:170) cannot fail on a type mismatch, so the failure surfaces on the second decode into the typed body (http.go:182-184) and is answered `WriteBadRequest`. Probe output: `{"data":{"device_registration_limit":"100"}}` → `json: cannot unmarshal string into Go struct field patchBody.device_registration_limit of type int32`; `{"data":{"jwt_public_key_pem":123}}` → the string twin; `{"data":{"device_registration_limit":1.5}}` → `cannot unmarshal number 1.5`; `{"data":[]}` → `cannot unmarshal array into Go value of type patchBody`. Same on POST, which decodes straight into `realmBody` (http.go:63-67). This is unreachable from a spec-conformant client — `RealmCreate`/`RealmPatch` already type both fields — so document it as a one-line extension of the 400 description on both operations rather than inventing a 422, and do NOT re-shape the `BadRequest` example, whose `detail: Bad request` is the frozen constant `TestHousekeepingErrorDetailExamples` pins. One distinction worth stating in the same description: a wrong-typed field is *not* the `invalid_update_parameters` case — unknown keys are caught earlier by the `patchAllowedFields` sweep (http.go:174-180, 422 FieldErrors) on PATCH only, and silently dropped on POST. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
272:- [x] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `components.parameters.RealmName` and `RealmCreate.realm_name` both carry `pattern: '^[a-z][a-z0-9]*$'` and `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go) holds both against the realm name CHECK extracted from migrations/000002_metadata.up.sql — but the read side is still bare `type: string`: `Realm.realm_name` (components.schemas.Realm) and `DataEnvelopeRealmNames.items`. Decide whether a response field earns the same pattern (a generated client learns nothing new from it, but a reader diffing create against get sees one constrained and one not) and, if it does, add it to both and extend that test to hold the response-side pattern against the same CHECK. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
280:- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
281:- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
282:- [x] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
283:- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
284:- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
287:- [!] docs-native-socket-handshake-errors [auto]: add the two WebSocket-handshake failures to both socket operations in docs/api/astrate_native_api.yaml — `/astrate/v1/{realm}/socket` responses (yaml:335-345) and `/appengine/v1/socket/websocket` (yaml:398-404) document only 101/401/403/500 and 101/401/500 — but measured today with the repo's own dependency (probe in /tmp, coder/websocket accept.go `verifyClientRequest`): a plain GET with no `Connection: Upgrade` answers **426 Upgrade Required**, `Content-Type: text/plain; charset=utf-8`, body `WebSocket protocol violation: Connection header "" does not contain Upgrade\n`, and `Sec-WebSocket-Version: 14` answers **400** `unsupported WebSocket protocol version (only 13 is supported): "14"\n`. Add both to each operation and say two things explicitly: the bodies are plain text from `http.Error`, NOT the `{"errors":{"detail":...}}` envelope every other response on this spec uses (a client that JSON-parses them breaks), and they are only reached after authentication — native: the route's guard `mw.RequireRealm(auth.ClaimChannels)` wraps `handle` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112), Phoenix: the realm/token checks at internal/appengine/channels/ws.go:51-68 run before `websocket.Accept` at :71 — and only on the WebSocket branch (`wantsSSE` is tested first at ws.go:78), so 401/403 and the 200 SSE path win. Do NOT add 405: ServeMux's method-mismatch applies equally to every GET route in every spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
288:- [x] docs-native-socket-sse-exact-accept [auto]: fix the SSE negotiation prose on `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml — the description (yaml:289-291), the `transport` parameter (yaml:328-334) and the `200` response (yaml:338-339) read as if any client asking for `text/event-stream` gets SSE, but the predicate `wantsSSE` is exact string equality (internal/appengine/stream/ws.go:151-153): measured, `Accept: text/event-stream, text/plain` and `Accept: text/event-stream;q=1.0` both return false and fall through to `websocket.Accept`, where a non-upgrade client gets the 426 above; `transport` selects SSE only for the literal `sse`, and `transport=websocket` / `transport=` (empty) both mean WebSocket. Say "exactly" for the header, say any other `Accept` falls through to the upgrade attempt, and do NOT claim the server rejects out-of-enum `transport` values — the enum is a client-side contract the handler never enforces. While there, give the `200` its real media type: the handler sets `Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive` and writes one `data: {...}` frame per event followed by a blank line (ws.go:123-126, 156-162). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
289:- [x] docs-native-socket-security-scheme [auto]: `/astrate/v1/{realm}/socket` is documented as if public in docs/api/astrate_native_api.yaml — the spec has no `components.securitySchemes` at all and root `security: []` (yaml:21), and the operation declares no `security`, so a client generated from this spec sends no `Authorization` header and is answered 401 `{"errors":{"detail":"Unauthorized"}}` by `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112, bearer token read only from the Authorization header at middleware.go:131, 178-196). Add a scheme mirroring `a_aea` (docs/api/astarte_appengine_api.yaml:1436-1443: apiKey in header `Authorization`, described as a realm JWT carrying the `a_ch` claim) and set `security` on this operation only — root `security: []` stays so health/readiness/metrics/version remain documented as unauthenticated. If you model the Phoenix twin too, its credential is `?token=` in the query (internal/appengine/channels/ws.go:53-54), i.e. apiKey `in: query`, not header. Description notes worth carrying: the scheme is case-insensitive `Bearer` with an optional colon, and an unknown realm is 401 not 404 (middleware.go:68-72). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
290:- [!] docs-native-socket-event-schema [auto]: document the native socket's event payload in docs/api/astrate_native_api.yaml — currently only a fenced `incoming_data` example inside the operation description (yaml:300-310), with no `components.schemas` entry anywhere (yaml:425-478), and that example is not the general shape: `wireEvent` marks `interface`, `path` and `value` `omitempty` (internal/appengine/stream/ws.go:51-58), so measured on a verbatim copy a lifecycle event marshals to `{"event":"device_connected","realm":"test","device_id":"dev1","timestamp":"2024-05-30T13:49:57.045Z"}` — no interface/path/value at all (and `value: false`/`value: 0` ARE emitted, because omitempty on an `any` drops only nil). Add a `NativeEvent` schema (event, realm, device_id, optional interface/path/value, timestamp) referenced from the operation, state that `event` is one of the internal/engine/stream/bus.go Kind* constants (bus.go:20-26), and record that this socket drops `InterfaceMajor`, `IP`, `ErrorName` and `ErrorMetadata`, which `stream.Event` carries (bus.go:39-65) — so a `device_error` arrives here with no reason, unlike the Phoenix twin, which marshals the whole struct (internal/appengine/channels/ws.go:406-416). That last sentence documents a code gap, not a spec one: do not add the fields; if they are wanted on the wire, that is a separate code task. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1

$ grep -n -i 'docs-sync\|docs/site\|mkdocs' .mule/for-giulio.md | tail -40
23:- **`docs/site/appengine-api.md:114` — "All responses use the Astarte envelope format" is contradicted by the endpoints the page itself lists above it.** The "Response envelope" section (lines 112-126) opens with that sentence and then the `{ "data": ... }` / `{ "errors": { "detail": ... } }` shapes, but the live-stream endpoints documented at lines 94-110 never return one: the native socket pushes bare event JSON as WebSocket text frames or `data: {json}` SSE frames (internal/appengine/stream/ws.go:105-110, 156-162), `/astrate/v1/health` answers `{"status":"ok"}` (internal/observability/health.go:51), `/astrate/v1/readiness` answers `{"status":...,"checks":{...}}` (health.go:77), `/astrate/v1/metrics` answers Prometheus text under `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (metrics.go:42), and — measured today with the repo's own `coder/websocket` (probe in /tmp) — a WebSocket handshake failure on either socket is plain text under `text/plain; charset=utf-8` (`WebSocket protocol violation: Connection header "" does not contain Upgrade`), not JSON. REST appengine paths do use the envelope, so the sentence is right for the endpoints above "Device data" and wrong unqualified. Proposed wording (your voice): scope it — "All REST responses use the Astarte envelope format" — with one line noting the socket and observability surfaces answer un-enveloped. Page untouched.
25:- **`docs/site/appengine-api.md:102` — "Honours `a_ch` claims as room filters" imports upstream's room semantics onto a socket that has no rooms.** The native socket subscribes to the whole realm bus, narrowed only by the `device_id`/`interface` query parameters (internal/appengine/stream/ws.go:69-75); there is no room concept and no per-room filtering. What a_ch actually does here is gate the whole route: `mw.RequireRealm(auth.ClaimChannels)` requires the token to authorize `GET socket` via the REST verb-regex rule (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112 -> internal/auth/claims.go:164-171). The upstream-measured JOIN/WATCH partition rule DESIGN.md §4.2 describes (`AuthorizesChannel`, claims.go:197-204) is used only by the Phoenix socket's per-room join/watch checks (internal/appengine/channels/ws.go:278, 354). Consequence worth stating: a token whose a_ch list is a blanket `".*::.*"` — which upstream's Channels rule treats as authorizing nothing — authorizes the native socket, while the Phoenix socket, where room filters actually live, still reads JOIN/WATCH literally. Proposed wording (your voice): "guarded by the `a_ch` claim (GET on `socket`); `device_id`/`interface` narrow the stream". Page untouched.
77:  **The knock-on: `docs/COMPATIBILITY.md` deviation 17 repeats the same error, and its wording about one endpoint is now measurably wrong.** It opens "**Always synchronous where upstream 1.4 defaults to asynchronous** — upstream 1.4 runs realm create/delete … and lets the caller opt into synchronous execution with `?async_operation=false`" (docs/COMPATIBILITY.md:335-347). Per the CHANGELOG the async default and the opt-out both date to 1.0.2, not 1.4 — the behaviour is not a moving target, which is the opposite of how that sentence reads. Separately, the sentence lists **device deletion** among the endpoints where the caller can opt into synchronous execution, and upstream has **no such opt-out there**: in `apps/astarte_realm_management/lib/astarte_realm_management_web/controllers/device_controller.ex`, the `operation :delete` doc comment says "Device deletion happens asynchronously, and receiving a 204 response" (:31-38) and the handler is bare — `with :ok <- Devices.delete_device(realm_name, device_id)` (:52-53), with `async_operation` never read. The policy twin does honour it, read as a bare string compare with no parsing (`trigger_policy_controller.ex:155-167`: `if Map.get(params, "async_operation") == "false"`), so upstream accepts only the literal `"false"` and treats absent/`0`/`no` as async. Astrate accepting-and-ignoring the parameter on device delete is therefore an **Astrate extension, not parity** — harmless and stronger, but the record calls it parity. For whoever picks up the queued `docs-sync-rm-delete-device-async-operation-param` line: documenting that parameter on the device-delete operation is still right for an Astrate client, it just has to be recorded as an extension, because upstream's spec has no such parameter on that operation.
128:  **The alarm stopped firing because of a bug in the switch, not because the mule recovered.** Last heartbeat **2026-09-28T18:01Z** (`.mule/.heartbeat`), and the last row in `.mule/log.md` is 2026-09-28 (`docs-sync-pairing-unregister-description`, `8ffb39e`) — no task attempt in the ~74h since, while recipe ticks are still committing (today 19:46–21:41Z), so the timer is alive and only the task queue is silent. On **2026-09-29T10:49Z** the pulse did fire (16h) but took the no-`gh` branch: the last line of this file is the fallback text "The mule has been idle 16h" (`mule.sh:771-773`) and **no issue was filed** — #114 is the last alarm ever filed. That same call created `.mule/.alarmed`, which only `beat` clears (`mule.sh:737`), and `check_pulse` returns immediately when it exists (`mule.sh:748`), so the switch has been **latched off for three days** and will file nothing until some tick lands work. The longest silence in the log is the one the switch reports least. That also gives the 2026-09-25 "the body is wrong" note a consequence at last: "Close it — the mule reopens a new one if the silence continues" (`mule.sh:764`) is false by construction — one alarm per silence, never a second.
156:## 2026-09-25 — docs-sync run, surface: **housekeeping** (the least-covered of the five)
159:`.mule/todo.md` (`docs-sync-hk-async-operation-param`, `docs-sync-hk-retention-zero-is-unset`);
167:- **`docs/site/configuration-reference.md` `[housekeeping]` lists 2 of the 4 keys it has, and
183:- **`docs/site/housekeeping-api.md` has no PATCH section.** The page enumerates Create / List /
214:(`yaml.Unmarshal`, 2 paths) — `mkdocs` is not installed on this box, so `make -C docs build`
216:earlier docs-sync line does.
240:- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`ASTRATE_FLOW_CONFIG`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
292:- **docs-sync realm-management, 2026-09-17: `docs/site/realm-management-api.md` documents an "Update trigger" endpoint Astrate does not serve.** Lines 81-86 show `PUT /realmmanagement/v1/<realm>/triggers/<name>` with a `{ "data": <updated trigger JSON> }` body, but no such route is registered — internal/realm/http.go:44-47 wires only GET/POST `/triggers` and GET/DELETE `/triggers/{name}` (an update op is deliberately absent; not even a stub or TODO exists). Site prose is yours — drop the "Update trigger" section or confirm the route is planned.
310:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
326:- **docs-sync pairing, 2026-09-10: docs/site names the pre-credentials wire status `registered`; the API returns `pending`.** `docs/site/pairing-and-security.md:56` ("flips status `registered -> confirmed`") and `:91` ("**registered** -- device registered, awaiting first credentials request"), plus `docs/site/data-modeling.md:94`, all give `registered` as a device status value. But `service.Info` emits `pending` for any device that is neither confirmed nor inhibited (internal/pairing/service.go:296-299 — upstream-parity per the comment at service.go:282-285; the DB value is `registered`, store/devices.go:19-27). Wire values are `pending`/`confirmed`/`inhibited`. Site prose is yours — reword to `pending`, or confirm the site intentionally describes the DB value.
392:- **docs-sync realm-management: two ASTRATE_ config keys exist in code but are missing from
393:  `docs/site/configuration-reference.md`** (site prose is yours, never edited by the mule).
396:  (lines 270-271 and 288); neither appears in `docs/site/configuration-reference.md`.
481:- ~~`docs/site/appengine-api.md` documents `GET` and `DELETE` on `/appengine/v1/<realm>/groups/<name>`~~
594:## 2026-10-04 — three housekeeping drifts in `docs/site/` the recipe is not allowed to touch
596:Ran the docs-sync recipe over the housekeeping surface (code routes vs
598:`.mule/todo.md`; the prose half needs you, because `docs/site/` is on the never-touch list.
600:**1. `docs/site/housekeeping-api.md:30` still talks about Cassandra.**
614:`docs/site/housekeeping-api.md` goes Create → List → Get → Delete (lines 10-50), so the route
615:`docs-sync-hk-patch-endpoint` added and documented in the spec
618:`RealmCreate` since `docs-sync-hk-retention-field`, so its create example (lines 12-21) shows a
622:The table at `docs/site/configuration-reference.md:79-84` has exactly two rows and the section

$ grep -n 'docs-sync\|surface' .mule/log.md | tail -40
99:| 2026-09-17 | docs-sync-rm-policies-delete-422 | done | 117s | c0521bc |
100:| 2026-09-17 | docs-sync-rm-triggers-422-nested-envelope | done | 161s | 8bfedd6 |
108:| 2026-09-20 | docs-sync-native-metrics-example-fake-series | done | 93s | 43c4a8b |
109:| 2026-09-20 | docs-sync-native-socket-missing-403-500 | done | 116s | f9d8b27 |
113:| 2026-09-22 | docs-sync-ae-read-query-params | blocked | 159s | wrote nothing |
120:| 2026-09-23 | docs-sync-rm-put-interface-409 | done | 356s | 7256af9 |
121:| 2026-09-24 | docs-sync-rm-interface-422-shapes | done | 490s | 93e784e |
122:| 2026-09-24 | docs-sync-rm-validation-example-prefix | done | 238s | 02e269b |
123:| 2026-09-24 | docs-sync-rm-auth-403 | done | 668s | 9c05370 |
129:| 2026-09-25 | docs-sync-hk-async-operation-param | done | 248s | ed1cd72 |
130:| 2026-09-25 | docs-sync-hk-retention-zero-is-unset | done | 325s | c2ad442 |
131:| 2026-09-25 | docs-sync-rm-async-operation-param | done | 425s | 3f0eca8 |
132:| 2026-09-25 | docs-sync-rm-delete-device-async-operation-param | done | 903s | f965f39 |
133:| 2026-09-25 | docs-sync-rm-delete-device-async-operation-param | blocked | 177s | wrote nothing |
139:| 2026-09-26 | docs-sync-ae-forbidden-403 | done | 1037s | 6a29ac7 |
140:| 2026-09-27 | docs-sync-ae-patch-by-alias-409 | done | 119s | d1af059 |
141:| 2026-09-27 | docs-sync-ae-add-group-device-422 | done | 274s | dff3be8 |
142:| 2026-09-27 | docs-sync-ae-delete-data-400 | done | 330s | 66bc366 |
143:| 2026-09-27 | docs-sync-ae-list-devices-422 | transient | 344s |  > build · big-pickle  Error: Error from provider (Console): Rate limit exceeded. Please try again later.  |
144:| 2026-09-27 | docs-sync-ae-list-devices-422 | transient | 248s |  > build · big-pickle  Error: Error from provider (Console): Rate limit exceeded. Please try again later.  |
145:| 2026-09-27 | docs-sync-ae-list-devices-422 | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
146:| 2026-09-28 | docs-sync-pairing-error-example-capitalisation | done | 166s | 7ed78cb |
147:| 2026-09-28 | docs-sync-pairing-info-version-example | done | 333s | d9ec98f |
148:| 2026-09-28 | docs-sync-pairing-deviceid-base64url | done | 480s | 32693cf |
149:| 2026-09-28 | docs-sync-pairing-initial-payload-format-enum | done | 95s | 1f6a6c2 |
150:| 2026-09-28 | docs-sync-pairing-unregister-description | done | 261s | 8ffb39e |
157:| 2026-10-02 | docs-sync-rm-error-example-capitalisation | done | 751s | 1d1e6f4 |
158:| 2026-10-02 | docs-sync-rm-legacy-alias-fields | blocked | 329s | wrote nothing |
159:| 2026-10-02 | docs-sync-rm-validationerror-example | done | 732s | 8a09a1a |
160:| 2026-10-02 | docs-sync-rm-deviceid-param | done | 303s | 2d1a1c1 |
161:| 2026-10-02 | docs-sync-rm-update-interface-body | blocked | 1201s | TIMEOUT after 1201s — task too big, split it |
162:| 2026-10-03 | docs-sync-ae-error-example-capitalisation | done | 298s | f3bcc6d |
163:| 2026-10-03 | docs-sync-hk-error-example-capitalisation | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
164:| 2026-10-03 | docs-sync-native-error-example-capitalisation | done | 284s | 7cc237f |
165:| 2026-10-03 | docs-sync-hk-error-detail-examples-split | done | 114s | 23006d4 |
170:| 2026-10-04 | docs-sync-hk-realm-name-pattern | done | 525s | e728217 |
171:| 2026-10-04 | docs-sync-hk-errordetail-schema-example | done | 184s | 5044143 |
172:| 2026-10-04 | docs-sync-hk-patch-422-field-error-examples | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
173:| 2026-10-04 | docs-sync-hk-wrong-type-field-400 | blocked | 152s | wrote nothing |
174:| 2026-10-04 | docs-sync-hk-realm-name-response-schemas | done | 178s | def0af6 |

$ grep -n '^## ' .mule/for-giulio.md | tail -25; echo ---; ls -la .mule/reports | tail -20
156:## 2026-09-25 — docs-sync run, surface: **housekeeping** (the least-covered of the five)
561:## 2026-08-23 — FDO promoted to milestone-4.0 candidate (Giulio's decision, recorded)
573:## 2026-09-25 — the `-tags nats` half of the tree is compiled by no gate
594:## 2026-10-04 — three housekeeping drifts in `docs/site/` the recipe is not allowed to touch
---
total 56
drwxr-xr-x 3 root root  4096 Sep  6 19:40 .
drwxr-xr-x 7 root root  4096 Oct  7 20:47 ..
-rw-r--r-- 1 root root  2279 Sep  6 19:40 hygiene.md
drwxr-xr-x 2 root root 40960 Oct  7 20:47 log

$ grep -n '^- \[.\] docs' .mule/todo.md | cut -c1-160
80:- [x] docs-sync-rm-policies-delete-422 [auto]: add the missing `422` response to `DELETE /realmmanagement/v1/{realm}/policies/{name}` in docs/api/astarte_rea
81:- [x] docs-sync-rm-triggers-422-nested-envelope [auto]: document the nested changeset envelope on `POST /realmmanagement/v1/{realm}/triggers` in docs/api/ast
116:- [x] docs-sync-pairing-health-path [auto]: add the undocumented `GET /pairing/v1/{realm}/health` route to docs/api/astarte_pairing_api.yaml — it exists i
117:- [x] docs-sync-pairing-register-404 [auto]: document the missing `404` (DeviceNotFound) on `POST /pairing/v1/{realm}/agent/devices` in docs/api/astarte_pai
123:- [x] docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` 
124:- [x] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/as
127:- [x] docs-sync-appengine-by-alias-endpoints [auto]: add the 7 undocumented by-alias mirror routes to docs/api/astarte_appengine_api.yaml — PATCH /devices
128:- [x] docs-sync-appengine-group-endpoints [auto]: add the 8 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (r
129:- [x] docs-sync-appengine-get-group-device [auto]: add the undocumented `GET /groups/{group}/devices/{device}` route to docs/api/astarte_appengine_api.yaml 
130:- [x] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /grou
132:- [x] docs-sync-appengine-group-patch-status [auto]: PATCH /groups/{group}/devices/{device} (internal/appengine/http.go:518-537, patchGroupDevice) shares ap
133:- [x] docs-sync-appengine-data-422-interface-level [auto]: the interface-level data GETs — GET /devices/{device}/interfaces/{interface}, GET /devices-by-a
136:- [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accept
143:- [x] docs-sync-appengine-device-status-schema [auto]: fix the `DeviceStatus` schema in docs/api/astrate_appengine_api.yaml — the wire emits `id` (interna
144:- [!] docs-sync-appengine-data-set-422 [auto]: add the missing `422` response to the six PUT/POST data-set operations in docs/api/astrate_appengine_api.yaml
145:- [x] docs-sync-appengine-downsample-min [auto]: fix the `downsample_to` constraint in docs/api/astrate_appengine_api.yaml — the DataDownsample descriptio
149:- [x] docs-sync-hk-patch-endpoint [auto]: add the undocumented `PATCH /housekeeping/v1/realms/{realm}` route to docs/api/astarte_housekeeping_api.yaml — i
150:- [x] docs-sync-hk-retention-field [auto]: add the `datastream_maximum_storage_retention` field to the `Realm` and `RealmCreate` schemas in docs/api/astarte
154:- [x] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /real
155:- [x] docs-sync-pairing-status-enum [auto]: fix the `PairingInfo.status` enum in docs/api/astarte_pairing_api.yaml from `[confirmed, pending, denied, expire
156:- [x] docs-sync-pairing-version-endpoint [auto]: add the undocumented `GET /pairing/v1/{realm}/version` route to docs/api/astarte_pairing_api.yaml — it ex
157:- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GE
158:- [x] docs-sync-rm-delete-interface-status [auto]: fix the `deleteInterface` response in docs/api/astarte_realm_management_api.yaml — the spec documents `
159:- [x] docs-sync-rm-mapping-required-encrypted [auto]: add `required` (boolean, datastream-only) and `encrypted` (boolean, datastream-only) to the `Interface
160:- [x] docs-sync-rm-put-auth-422 [auto]: add the missing `422` (ValidationError) response to `PUT /realmmanagement/v1/{realm}/config/auth` in docs/api/astart
161:- [!] docs-sync-rm-version-example [auto]: fix the `getVersion` example in docs/api/astarte_realm_management_api.yaml from `data: "1.1.0"` (yaml:557) to `da
167:- [x] docs-sync-ae-write-405 [auto]: add the missing `405` response to the nine data-write ops in docs/api/astarte_appengine_api.yaml — PUT/POST/DELETE `/
168:- [x] docs-sync-ae-write-value-422 [auto]: add the missing `422` response to the six PUT/POST data-write ops in docs/api/astarte_appengine_api.yaml (device-
169:- [x] docs-sync-ae-post-groups-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups` in docs/api/astarte_appengine_api.yaml — 
172:- [x] docs-sync-ae-post-group-devices-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_
173:- [x] docs-sync-ae-patch-merge-patch-media-type [auto]: fix the three PATCH request bodies in docs/api/astarte_appengine_api.yaml — PATCH /devices/{device
178:- [x] docs-sync-hk-delete-gating-responses [auto]: add the missing `405` and `422` responses to `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte
179:- [x] docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_housekeeping_api.yaml (yaml:417) — it shows `detai
185:- [x] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication 
188:- [x] docs-sync-pairing-deviceendpoints-dead-404-403 [auto]: remove the unreachable `404` DeviceNotFound response from the three device-facing ops in docs/a
189:- [x] docs-sync-pairing-version-404-unreachable [auto]: remove the `404` "Unknown realm" response from `GET /pairing/v1/{realm}/version` in docs/api/astarte
194:- [x] docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Retu
196:- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/
201:- [x] docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate
202:- [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yam
208:- [!] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine
216:- [x] docs-sync-rm-put-interface-409 [auto]: add the missing `409` Conflict response to `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` in docs/
217:- [x] docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1
218:- [x] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it
219:- [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `Req
226:- [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /hous
227:- [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions 
228:- [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.
230:- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no 
236:- [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml �
237:- [x] docs-sync-ae-patch-by-alias-409 [auto]: add the missing `409` response to `patchDeviceByAlias` in docs/api/astarte_appengine_api.yaml (yaml:241-291, c
238:- [x] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_ap
239:- [x] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data DELETE operations in docs/api/astarte_appengine_api.yaml (`delet
240:- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (ya
241:- [x] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadReques
242:- [x] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pairing_api.yaml, which shows `version: "0
243:- [x] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` pa
244:- [x] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in docs/api/astarte_pairing_api.yaml (yaml:37
245:- [x] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` description in docs/api/astarte_pairing_api.yaml (yaml:153), "Removes a de
254:- [x] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRe
255:- [!] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy alias fields to `InterfaceDefinition` and `InterfaceMapping` in docs/a
256:- [x] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yam
257:- [x] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstra
258:- [!] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_realm_management_api.yaml — `data` i
260:- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` show
261:- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDe
262:- [x] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerEr
263:- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spe
268:- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^
269:- [x] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-
270:- [!] docs-sync-hk-patch-422-field-error-examples [auto]: add the three missing PATCH 422 bodies to the `ValidationErrors` response component in docs/api/as
271:- [!] docs-sync-hk-wrong-type-field-400 [auto]: widen the `BadRequest` component description in docs/api/astarte_housekeeping_api.yaml (yaml:378, "malformed
272:- [x] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `co
280:- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it i
281:- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promi
282:- [x] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The respons
283:- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null,
284:- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmNa
287:- [!] docs-native-socket-handshake-errors [auto]: add the two WebSocket-handshake failures to both socket operations in docs/api/astrate_native_api.yaml —
288:- [x] docs-native-socket-sse-exact-accept [auto]: fix the SSE negotiation prose on `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml — the 
289:- [x] docs-native-socket-security-scheme [auto]: `/astrate/v1/{realm}/socket` is documented as if public in docs/api/astrate_native_api.yaml — the spec ha
290:- [!] docs-native-socket-event-schema [auto]: document the native socket's event payload in docs/api/astrate_native_api.yaml — currently only a fenced `in

$ grep -c '' .mule/todo.md; grep -n '^- \[ \]\|^- \[!\]' .mule/todo.md | wc -l; grep -n '^- \[ \]\|^- \[!\]' .mule/todo.md | cut -c1-120
290
51
83:- [!] hk-zero-reglimit-create-asymmetry: POST /realms persists device_registration_limit 0 as a literal 0 (service.go
86:- [ ] broker-external-bus-intake [legion]: implement the second Intake implementation named in the TODO at internal/b
88:- [ ] hygiene-govulncheck [legion]: on the Legion Go, run govulncheck ./... on a fresh ~/astrate clone and file a tas
89:- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.g
101:- [ ] race-check-store: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go 
102:- [ ] race-check-engine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go
103:- [ ] race-check-flow: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go t
104:- [ ] race-check-appengine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main &&
105:- [ ] race-check-pkg: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go te
108:- [!] broker-disconnect-device-zombie-session: in `internal/broker/broker.go` `DisconnectDevice` (lines 260-266), fo
109:- [!] broker-offline-acl-tests: in `internal/broker/aclhook_test.go`, unit-test the offline-delivery ACL — `offlin
110:- [!] broker-onconnect-doc-comment: in `internal/broker/authhook.go:314`, restore the missing first line of the `OnC
113:- [!] empty-introspection-verification [auto]: upstream v1.3.0 changed "allow devices with empty introspection" — 
114:- [!] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installati
119:- [!] probe-property-resend-encoding [auto]: upstream v1.3.0 fixed outbound server-property values sent to a device 
120:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, em
121:- [!] flow-validate-source-sink: `Pipeline.Validate`'s no-source/no-sink error branches (internal/flow/pipeline.go:1
126:- [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) p
134:- [!] appengine-snapshot-ignores-query-params [auto]: in `internal/appengine/data.go:138-148`, the individual-datast
135:- [!] appengine-group-token-roundtrip-test [auto]: add a container-free unit test for `parseGroupToken`/`groupTokenF
138:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoi
144:- [!] docs-sync-appengine-data-set-422 [auto]: add the missing `422` response to the six PUT/POST data-set operation
161:- [!] docs-sync-rm-version-example [auto]: fix the `getVersion` example in docs/api/astarte_realm_management_api.yam
199:- [!] channels-rejoin-authz-mismatch: `handleJoin`'s comment says upstream authorizes every join including a rejoin 
200:- [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a loc
203:- [ ] store-devices-inhibit-re-register [legion] [auto]: `RegisterDevice` (internal/store/devices.go:75-91) silently
205:- [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (intern
208:- [!] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-da
210:- [!] flow-mqtt-source-reconnect-recovery [auto]: `mqttSource` latches `lost` on paho's ConnectionLost handler (inte
213:- [!] flow-randomsource-span-overflow [auto]: `randomSource.next()` computes `rand.Int64N(s.maxInt-s.minInt+1)` (int
221:- [ ] flow-boot-resumes-stopped-durable-flows [legion] [auto]: RehydrateAutoRestart (cmd/astrate/main.go:216) restar
224:- [!] cmd-healthcheck-contract-test [auto]: add a container-free test for `runHealthcheck` (cmd/astrate/main.go:531-
230:- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_real
240:- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in
246:- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailin
251:- [!] container-timeout-bounds [auto]: bound `timeout_ms` and `ready_timeout_ms` in `parseConfig` (internal/flow/blo
253:- [!] container-response-cap-test [auto]: pin `Bridge.RoundTrip`'s 1 MiB response cap (internal/flow/blocks/containe
255:- [!] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy alias fields to `InterfaceDefi
258:- [!] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astart
259:- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-con
261:- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astart
264:- [!] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in
266:- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still 
270:- [!] docs-sync-hk-patch-422-field-error-examples [auto]: add the three missing PATCH 422 bodies to the `ValidationE
271:- [!] docs-sync-hk-wrong-type-field-400 [auto]: widen the `BadRequest` component description in docs/api/astarte_hou
273:- [!] store-register-inhibit-preserve: RegisterDevice must not overwrite status to 'registered' when a device was in
275:- [!] store-latestindividual-empty-errnotfound: Add store-level test asserting LatestIndividual returns ErrNotFound 
276:- [!] pairing-burn-bcrypt-cost [auto]: make the dummy-secret timing burn use the configured BcryptCost — burnBcryp
285:- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-1
287:- [!] docs-native-socket-handshake-errors [auto]: add the two WebSocket-handshake failures to both socket operations
290:- [!] docs-native-socket-event-schema [auto]: document the native socket's event payload in docs/api/astrate_native_

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

$ rg -n 'Handle|mux\.|Route|http\.(Get|Post)|\.Get\(|\.Post\(|\.Delete\(|\.Put\(' internal/pairing/http.go | head -60
46:	requireAgent func(http.Handler) http.Handler
74:// Mount registers the pairing routes on mux. Paths are wire-frozen
81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
82:	mux.Handle("POST /pairing/v1/{realm}/agent/devices",
83:		a.requireAgent(http.HandlerFunc(a.handleRegister)))
84:	mux.Handle("DELETE /pairing/v1/{realm}/agent/devices/{deviceID}",
85:		a.requireAgent(http.HandlerFunc(a.handleUnregister)))
86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
344:	header := r.Header.Get("Authorization")

$ rg -n 'pairing' cmd/astrate/main.go | head -40; echo ---; ls internal/pairing/
2:// it wires the store, ingestion engine, embedded MQTT broker, pairing, the M7
44:	"github.com/astrate-platform/astrate/internal/pairing"
147:	pairer := pairing.New(st, sealer, pairing.Config{
160:	// register first-seen ids through the pairing door; an id already taken
173:			if errors.Is(err, pairing.ErrAlreadyRegistered) || errors.Is(err, store.ErrDeviceAlreadyConfirmed) {
353:func mountAPIs(cfg config.Config, st *store.Store, e *engine.Engine, b *broker.Broker, sealer *store.KeySealer, metrics *observability.Metrics, flowSvc *flowapi.Service, pairer *pairing.Service, log *slog.Logger) (http.Handler, *housekeeping.Service, error) {
357:	pairing.NewAPI(pairer, mw, pairing.APIConfig{
386:	for _, svc := range []string{"appengine", "realmmanagement", "pairing"} {
396:	for _, svc := range []string{"appengine", "realmmanagement", "pairing", "housekeeping"} {
431:// requires a realm token, pairing's is public. Both report the emulated
438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
---
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

$ rg -n '^\s{2,4}(/|  [a-z]+:|operationId:|summary:)' docs/api/astarte_pairing_api.yaml | sed -n '1,400p'
14:    name: "Apache-2.0 with SECO Exclusion"
15:    url: https://www.apache.org/licenses/LICENSE-2.0
19:    variables:
27:  /pairing/v1/{realm}/health:
28:    get:
30:      summary: Realm health check
31:      description: |
38:      tags: [Pairing]
39:      parameters:
41:      responses:
70:  /pairing/v1/{realm}/version:
71:    get:
73:      summary: Get emulated API version
74:      description: |
79:      tags: [Pairing]
80:      parameters:
82:      responses:
97:  /pairing/v1/{realm}/agent/devices:
98:    post:
100:      summary: Register a new device (flow A)
101:      description: |
106:      tags: [Pairing]
107:      security:
109:      parameters:
124:      responses:
149:  /pairing/v1/{realm}/agent/devices/{deviceID}:
150:    delete:
152:      summary: Unregister a device
153:      description: |
165:      tags: [Pairing]
166:      security:
168:      parameters:
171:      responses:
183:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
184:    post:
186:      summary: Request device credentials (flow B)
187:      description: |
193:      tags: [Pairing]
194:      security:
196:      parameters:
212:      responses:
235:  /pairing/v1/{realm}/devices/{deviceID}:
236:    get:
238:      summary: Get device pairing info (flow C)
239:      description: |
243:      tags: [Pairing]
244:      security:
246:      parameters:
249:      responses:
269:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:
270:    post:
272:      summary: Verify device credentials
273:      description: |
277:      tags: [Pairing]
278:      security:
280:      parameters:
296:      responses:
343:      type: apiKey
344:      in: header
345:      name: Authorization
346:      description: |
351:      type: http
352:      scheme: bearer
353:      description: |
360:      name: realm
361:      in: path
362:      required: true
363:      description: The realm name.
364:      schema:
368:      name: deviceID
369:      in: path
370:      required: true
371:      description: |
376:      schema:
383:      type: object
384:      required: [hw_id]
385:      properties:
409:      type: object
410:      required: [csr]
411:      properties:
417:      type: object
418:      required: [client_crt]
419:      properties:
427:      type: object
428:      required: [credentials_secret]
429:      properties:
435:      type: object
436:      required: [client_crt]
437:      properties:
443:      type: object
444:      required: [protocols, status, version]
445:      properties:
468:      type: object
469:      required: [timestamp, until, valid]
470:      properties:
482:      type: object
483:      required: [cause, details, timestamp, valid]
484:      properties:
503:      type: object
504:      required: [data]
505:      properties:
510:      type: object
511:      required: [data]
512:      properties:
517:      type: object
518:      required: [data]
519:      properties:
524:      type: object
525:      required: [status]
526:      properties:
533:      type: object
534:      required: [data]
535:      properties:
540:      type: object
541:      required: [data]
542:      properties:
547:      type: object
548:      required: [data]
549:      properties:
556:      type: object
557:      required: [errors]
558:      properties:
567:      type: object
568:      required: [errors]
569:      properties:
579:      description: Bad Request — malformed or missing data envelope.
580:      content:
589:      description: Unauthorized — missing or invalid JWT/Bearer token.
590:      content:
599:      description: Forbidden — valid credentials but insufficient permissions or device inhibited.
600:      content:
609:      description: Device not found — unknown device ID.
610:      content:
619:      description: Unprocessable Entity — field-level validation errors.
620:      content:
669:      description: Too Many Requests — rate limit exceeded.
670:      content:
679:      description: Internal Server Error — unexpected failure.
680:      content:

$ rg -n '^\s{6}[0-9]{3}:$' docs/api/astarte_pairing_api.yaml
(no output)

$ rg -n 'Write[A-Z][A-Za-z]*\(|Status[A-Za-z]*|\.Status =|WriteHeader\(' internal/pairing/http.go
111:		_ = astarteapi.WriteError(w, http.StatusTooManyRequests, detailTooManyRequests)
117:		_ = astarteapi.WriteBadRequest(w)
121:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
131:	_ = astarteapi.WriteData(w, http.StatusCreated, registerResponse{CredentialsSecret: secret})
140:	w.WriteHeader(http.StatusNoContent)
159:		_ = astarteapi.WriteError(w, http.StatusTooManyRequests, detailTooManyRequests)
165:		_ = astarteapi.WriteUnauthorized(w)
170:		_ = astarteapi.WriteBadRequest(w)
174:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
184:	_ = astarteapi.WriteData(w, http.StatusCreated, credentialsResponse{ClientCrt: clientCrt})
197:	Status  string `json:"status"`
204:		_ = astarteapi.WriteUnauthorized(w)
215:	resp.Status = info.Status
217:	_ = astarteapi.WriteData(w, http.StatusOK, resp)
226:// CredentialsStatusView shapes: timestamps are rendered as Elixir
245:		_ = astarteapi.WriteUnauthorized(w)
250:		_ = astarteapi.WriteBadRequest(w)
254:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
265:		_ = astarteapi.WriteData(w, http.StatusOK, verifyValidResponse{
272:	_ = astarteapi.WriteData(w, http.StatusOK, verifyInvalidResponse{
288:		_ = astarteapi.WriteData(w, http.StatusOK, map[string]string{"status": "ok"})
290:		_ = astarteapi.WriteNotFound(w)
292:		_ = astarteapi.WriteError(w, http.StatusServiceUnavailable, "unhealthy")
308:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
311:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
314:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
317:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
320:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
323:		_ = astarteapi.WriteUnauthorized(w)
325:		_ = astarteapi.WriteForbidden(w)
327:		_ = astarteapi.WriteDeviceNotFound(w)
329:		_ = astarteapi.WriteInternalServerError(w)

$ rg -n '^\s*[0-9]{3}:$' docs/api/astarte_pairing_api.yaml
(no output)

→ Read internal/pairing/http.go

$ rg -n '^\s*(2|4|5)[0-9]{2}:$' docs/api/astarte_pairing_api.yaml
(no output)

→ Read docs/api/astarte_pairing_api.yaml [offset=27, limit=70]

$ rg -n '"[2-5][0-9]{2}":' docs/api/astarte_pairing_api.yaml
42:        "200":
51:        "404":
60:        "503":
83:        "200":
125:        "201":
134:        "400":
136:        "401":
138:        "403":
140:        "404":
142:        "422":
144:        "429":
146:        "500":
172:        "204":
174:        "401":
176:        "403":
178:        "404":
180:        "500":
213:        "201":
222:        "400":
224:        "401":
226:        "403":
228:        "422":
230:        "429":
232:        "500":
250:        "200":
264:        "401":
266:        "500":
297:        "200":
329:        "400":
331:        "401":
333:        "403":
335:        "422":
337:        "500":

→ Read docs/api/astarte_pairing_api.yaml [limit=130, offset=97]

$ rg -n 'func \(s \*Service\)' internal/pairing/service.go
174:func (s *Service) Register(ctx context.Context, realmName, hwID, initialFormat string) (string, error) {
230:func (s *Service) Unregister(ctx context.Context, realmName, deviceIDStr string) error {
247:func (s *Service) Credentials(ctx context.Context, realmName, deviceIDStr, secret, csrPEM string, ip netip.Addr) (string, error) {
286:func (s *Service) Info(ctx context.Context, realmName, deviceIDStr, secret string) (*Info, error) {
324:func (s *Service) VerifyCredentials(ctx context.Context, realmName, deviceIDStr, secret, clientCrtPEM string) (*VerifyResult, error) {
369:func (s *Service) Health(ctx context.Context, realm string) error {
380:func (s *Service) authenticateDevice(ctx context.Context, realmName, deviceIDStr, secret string) (*store.Realm, *store.Device, error) {
414:func (s *Service) loadCA(realm *store.Realm) (*ca.CA, error) {

$ rg -n 'ErrNotFound|ErrInhibited|ErrUnauthorized|ErrInvalidCSR|ErrAlreadyRegistered|ErrRegistrationLimit' internal/pairing/service.go
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
172:// secret; afterwards it fails with ErrAlreadyRegistered. initialFormat is
196:			return "", ErrRegistrationLimitReached
211:			return "", fmt.Errorf("%w: %s", ErrAlreadyRegistered, hwID)
229:// credential trail). store.ErrNotFound is returned for unknown devices.
233:		return fmt.Errorf("%w: device %q", store.ErrNotFound, deviceIDStr)
243:// device by credentials secret (uniform ErrUnauthorized on any mismatch),
253:		return "", ErrInhibited
262:		if errors.Is(err, ca.ErrInvalidCSR) {
263:			return "", fmt.Errorf("%w: %v", ErrInvalidCSR, err)
330:		return nil, ErrInhibited
371:		return err // store.ErrNotFound → 404
378:// device, wrong secret) return the same ErrUnauthorized after a comparable
384:		return nil, nil, ErrUnauthorized
387:	if errors.Is(err, store.ErrNotFound) {
389:		return nil, nil, ErrUnauthorized
395:	if errors.Is(err, store.ErrNotFound) {
397:		return nil, nil, ErrUnauthorized
405:		return nil, nil, ErrUnauthorized
408:		return nil, nil, ErrUnauthorized

→ Read internal/pairing/service.go [limit=245, offset=174]
→ Read docs/api/astarte_pairing_api.yaml [limit=115, offset=383]
→ Read docs/api/astarte_pairing_api.yaml [limit=140, offset=227]

$ rg -n 'APICompatVersion' internal/observability/compat.go internal/realm/*.go | head
internal/realm/dashboard_compat_test.go:35:		if v != APICompatVersion {
internal/realm/dashboard_compat_test.go:36:			t.Errorf("version = %q, want %q", v, APICompatVersion)
internal/realm/http.go:119:// feature UI on it — see APICompatVersion).
internal/realm/http.go:121:	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
internal/realm/service.go:580:// APICompatVersion is the upstream Realm Management API level Astrate
internal/realm/service.go:588:const APICompatVersion = "1.2.2"

→ Read docs/api/astarte_pairing_api.yaml [limit=120, offset=556]
→ Read docs/api/astarte_pairing_api.yaml [limit=58, offset=498]

$ rg -n -i 'unregister|un-register|removes a device|deletes a device|registration limit|device_registration' docs/site/pairing-and-security.md docs/site/architecture.md docs/site/data-modeling.md 2>/dev/null
docs/site/data-modeling.md:20:    device_registration_limit integer,

$ rg -n -i 'base64|hw_id|initial_payload|credentials_secret|42x|status' docs/site/pairing-and-security.md | head -40
38:{ "data": { "hw_id": "<22-char base64url device ID>" } }
39:→ 201 { "data": { "credentials_secret": "<44-char base64>" } }
44:- Optional extension: `"initial_payload_format": "json"` for JSON-profile devices.
50:Authorization: Bearer <credentials_secret>
56:- First successful call stamps `first_credentials_request` and flips status `registered -> confirmed`.
62:Authorization: Bearer <credentials_secret>
64:    "status": "confirmed",
77:  |<--credentials_secret--|   (secret delivered   |                        |
89:## Device status lifecycle

→ Read internal/auth/middleware.go [limit=120, offset=44]
→ Read docs/site/pairing-and-security.md

$ rg -n -i 'rate.?limit|per-device|per.IP' .mule/for-giulio.md
320:- **Extension point: external-bus Intake (NATS JetStream), 2026-09-10.** `internal/broker/intake.go:56-62` carries a `TODO(extension point)` for a second `Intake` implementation backed by a durable external bus (NATS JetStream) for multi-instance deployment or restart survival. The frozen design decision is in `docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4`; the consumer must reproduce per-device ordering, deferred-ack backpressure, and QoS 0 drop semantics. Not tracked by any open issue. Your call: file a milestone issue or park until the roadmap calls for it.

$ rg -n 'rotates the secret|re-registr|version|first_credentials_request' .mule/for-giulio.md | head
15:  **#93 at 33 days.** `git rev-list --left-right --count origin/main...HEAD` = **4 / 697** today (4 / 679 on 10-06), and `git merge-base --is-ancestor 24ad5b8 origin/main` still fails — the rewrite is entirely on `mule/queue` / `origin/mule/queue`. The version-drift detail recorded on 2026-10-06 (the comment says `v1.4.0-rc.5` / "v1.3.3 being the newest stable tag"; upstream is `v1.4.0-rc.6` / `v1.3.5`, re-checked via `gh api repos/astarte-platform/astarte/tags` today) folds into that same review rather than becoming a line of its own.
23:- **`docs/site/appengine-api.md:114` — "All responses use the Astarte envelope format" is contradicted by the endpoints the page itself lists above it.** The "Response envelope" section (lines 112-126) opens with that sentence and then the `{ "data": ... }` / `{ "errors": { "detail": ... } }` shapes, but the live-stream endpoints documented at lines 94-110 never return one: the native socket pushes bare event JSON as WebSocket text frames or `data: {json}` SSE frames (internal/appengine/stream/ws.go:105-110, 156-162), `/astrate/v1/health` answers `{"status":"ok"}` (internal/observability/health.go:51), `/astrate/v1/readiness` answers `{"status":...,"checks":{...}}` (health.go:77), `/astrate/v1/metrics` answers Prometheus text under `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (metrics.go:42), and — measured today with the repo's own `coder/websocket` (probe in /tmp) — a WebSocket handshake failure on either socket is plain text under `text/plain; charset=utf-8` (`WebSocket protocol violation: Connection header "" does not contain Upgrade`), not JSON. REST appengine paths do use the envelope, so the sentence is right for the endpoints above "Device data" and wrong unqualified. Proposed wording (your voice): scope it — "All REST responses use the Astarte envelope format" — with one line noting the socket and observability surfaces answer un-enveloped. Page untouched.
29:  **One detail to fold into that review rather than queue separately:** the rewritten comment at `internal/broker/aclhook.go` (the `control/keyAgreement` deny) says "…upstream's wire spec (082-key_agreement_protocol.md) only ships with v1.4.0-rc.5, v1.3.3 being the newest stable tag". The 082 citation, the "deliberately denied as Astrate's own choice" framing and the unchanged deny + `aclhook_test.go` case are exactly what #93 asked for, but the version reference has aged — upstream's newest stable is **v1.3.5** (2026-10-05) and newest overall `v1.4.0-rc.6` (2026-09-30). No separate task: one word in the same diff when someone finally reads it.
38:- **astarte-upstream recipe run, 2026-10-05 (evening): no upstream movement that touches anything Astrate implements, so nothing was queued — no `.mule/todo.md` line, no `gh issue create`, nothing commented, closed or edited on GitHub.** One thing *is* new since this recipe's last run (`2026-10-04T20:40:52+02:00` ran; the `2026-10-05T04:29:10+02:00` run timed out): stable **v1.3.5**, published `2026-10-05T14:11:44Z`. The milestone recipe already recorded it as inert today; I re-measured it rather than trusting that, because "a release note is a claim, not a fact" is this recipe's own rule. `gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5` is **4 commits, 35 files** — `8d47574f` + `4c4222cf` `chore: update horde`, `c4c63ebe` `chore(rm): bump xandra to 0.19.4`, `8f0f5ef1` `chore: prepare v1.3.5` — and **every app/lib file in it is a `mix.exs`/`mix.lock` version bump**. The rest is `CHANGELOG.md` (+6: one `### Fixed` entry, "[astarte_data_updater_plant] Increase RPC server reliability"), `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `README.md`, `astarte-dashboard/package.json`. Not one source file changed, so no route, MQTT topic, BSON encoding, control message or interface-schema field moved: Astrate implements no `data_updater_plant` and no Erlang RPC server. Wire-inert, same shape as v1.3.2/v1.3.3/v1.3.4. **Master since `v1.4.0-rc.6` is still the single commit `731bcddc`** ("chore: forward port release-1.4"; `gh api .../commits?since=2026-10-05T09:34:00Z` returns it and nothing newer), the migration-relocation commit the milestone run measured today as inert — so both wire-bearing post-rc.6 commits (`b6d46ad4`, `958bb9fe`) still have their owners in the queue and nothing new joins them.
40:  **The one delta to your open COMPATIBILITY.md wording proposal (the v1.3.2 one, folding v1.3.3/v1.3.4, further down this file): its version-reference sentence should read "newest stable v1.3.5" where it currently says "v1.3.4". Nothing else in it changes** — v1.3.5 adds no capability, so the same decision and the same "not yet emulated" list stand. Not applied: `docs/COMPATIBILITY.md` is untouched and still targets **1.2.2**, `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). I deliberately did **not** queue a fresh `compat-note-v1.3.5` line: the proposal is already open here, and a second line would re-derive it for the next run to find.
50:- **Milestone recipe run, 2026-10-05: v3.0's gate is unchanged, but the release sweep finally has something new to swallow — upstream shipped stable v1.3.5 today, seven hours ago, and it is wire-inert for Astrate. No issue filed, no task line queued.**  **The new release, measured.** `gh api repos/astarte-platform/astarte/releases` now shows a stable **v1.3.5** at `2026-10-05T14:11:44Z` where every prior run of this recipe (2026-09-10 onward, 20+ entries in this file) saw `v1.3.4` (2026-09-18) as newest. `compare/v1.3.4...v1.3.5` is **4 commits, 35 files**: two `chore: update horde` (`8d47574f`, `4c4222cf`), `chore(rm): bump xandra to 0.19.4` (`c4c63ebe`), `chore: prepare v1.3.5` (`8f0f5ef1`). The file list is `mix.exs`/`mix.lock` version bumps across all apps and libs, `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `astarte-dashboard/package.json`, `doc/mix.exs`, `README.md` — plus a **single** CHANGELOG entry, `## [1.3.5] - 2026-10-05` / `### Fixed` / "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available". That is the whole delta and none of it is Astarte-facing: Astrate implements no data_updater_plant and no Erlang RPC server (the only `rpc` hit under `internal/` is a comment at `internal/pairing/http.go:300` about AMQP reply-code conflicts), so no HTTP route, MQTT topic, AMQP control message or schema field moved. Same shape as the v1.3.2/v1.3.3/v1.3.4 maintenance releases already recorded here, so it is **not a v3.0 gap** and there is nothing to file.
94:  **The recipe's standing item is satisfied and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental*, reconcile trigger "upstream 1.4 final". 1.4 is not final, so there is nothing to promote and nothing to deprecate, and the register is accurate as written. `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked this run) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision ("bumped only in the same change that completes that surface, after reconciling every row tagged with the level"). The milestone-bump issue the recipe asks me to file **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": v3.0's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is delivered, the 1.4 half is open (#92, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut.
103:- **Delta on the two open COMPATIBILITY.md wording entries further down this file (the v1.3.4 one folding in v1.3.3, and the original v1.3.2 one): upstream `v1.4.0-rc.6` is wire-inert for Astrate, so nothing in that proposal is re-derived — only its version reference and one line of its "not yet emulated" list change.** Not applied: `docs/COMPATIBILITY.md` is untouched (still the 1.2.2 target, deviation 18's "which upstream 404s" at `docs/COMPATIBILITY.md:360`, `APICompatVersion = "1.2.2"` at `internal/realm/service.go:588`), and the FDO half of the rc.6 delta already has its own entry below rather than being folded in here.
116:  - **The voucher-upload body is not backwards compatible.** `hw_id` is now **required** and `initial_introspection` is accepted (`api_spec/schemas/ownership_voucher.ex`: `required: [:hw_id, :ownership_voucher, :key_name, :key_algorithm]`), and `LoadRequest.changeset/2` adds three rejectable cases that surface as `400`: `guid has already been claimed`, device `already exists`, and invalid `initial_introspection` (negative or malformed versions). `403` (Forbidden / AuthorizationPathNotMatched) is also now documented on `register`. **Correction to the CHANGELOG wording "Allow specifying device_id on ownership voucher upload": on the wire the client sends `hw_id`, and `device_id` is a `virtual: true` Ecto field derived from it** via `Device.decode_device_id(hw_id, allow_extended_id: true)` — so it is the 128-bit device id, not an arbitrary one, and the new "already exists" check is a realm-device collision check. Astrate must derive it the same way (`pkg/deviceid`) or it will disagree on the id space.
151:  **The recipe's standing item is satisfied, and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still has exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental* with reconcile trigger "upstream 1.4 final". v1.4 is not final, so both rows stay as written: nothing to promote, nothing to deprecate, and the register is accurate. `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision; the milestone-bump issue the recipe asks for **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": the milestone's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is done, and the 1.4 half is unresolved (#92 open, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut. Also still true and already escalated on 2026-09-06, so not repeated: the section's source document `.mule/research/upstream-parity-2026-08-22.md` does not exist in the repo.

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

$ rg -o '\bASTRATE_[A-Z_]+' -N docs/site/configuration-reference.md | sort -u
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

→ Read pkg/astarteapi/envelope.go
→ Read .mule/for-giulio.md [limit=40]
Now I've verified the pairing surface thoroughly. Let me record the one docs/site contradiction in `for-giulio.md`, since `docs/site/` is on the never-touch list.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,10 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **docs-sync run, 2026-10-07, surface: pairing — the `docs/api` spec is accurate (routes, status codes, request/response schemas and examples all match the handlers; no docs/api task to propose), but `docs/site/pairing-and-security.md:97` carries a rate-limiting claim the code contradicts: "Pairing endpoints: per-IP and per-device token buckets."** Actual per endpoint, read from the handlers — register is **per-IP only** (`a.regLimiter.Allow("ip|"+remoteIP(r))`, internal/pairing/http.go:110; one bucket, defaults 5.0/20 at http.go:23-24), credentials is the only one that is **per-IP AND per-device** (`credLimiter.Allow("ip|"+...) || credLimiter.Allow("dev|"+realm+"/"+deviceID)`, http.go:158), and **unregister, info, verify and health are not rate-limited at all** (handleUnregister http.go:134-141, handleInfo :201-218, handleVerify :242-277, handleHealth :284-294 — the health mount comment says so outright, "No rate limiter either; it is one cheap query", http.go:78-80). The API reference already documents the true per-op scope (registerDevice yaml:105 "Rate-limited per IP."; requestCredentials yaml:192 "Rate-limited per IP and per device."; the other four ops say nothing), so a client that reads the site for the surface's security posture is shown a limit that unregister/info/verify/health do not have. Proposed wording (your voice): perf-`endpoint`, e.g. "Register: per-IP token bucket. Credentials: per-IP and per-device token buckets. Unregister, info, verify, health: not rate-limited." Page untouched.
+
 - **github-issues triage run, 2026-10-07: nothing proposable for the 25th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub.** The recipe's `--limit 40` command prints all **31** open issues (authoritative `--limit 100` agrees): **27** `mule-alarm` **#94–#120** plus the same four non-alarm, none of them proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (`milestone-4.0`, parked until the v3.0 queue clears), and one the standing instruction excludes from triage entirely. New since the 2026-10-06 run: **#120**, created 2026-10-07T11:03:45Z, "nothing has landed in 14h" — arithmetic pinned, not inferred: `.mule/.heartbeat` holds `1791318570` = **2026-10-06T20:29:30Z**, written 3s after the `docs-native-socket-security-scheme` done land (`87c18cf`, 20:29:27Z), and 11:03:45Z − 20:29:30Z = **14h34m**, so `age=$(( (now-last)/3600 ))` prints 14 — the title verbatim; `.mule/.alarmed` carries mtime 11:04Z, i.e. set by that filing. Genuine window, the same 14–16h overnight band #117–#119 measured. **Proposal, extending the 2026-10-06 one: close #94–#119, all twenty-six** — each is a one-day alarm superseded by the next and nothing reads issue state; leave **#120** (today's, real window) to self-expire.
   **New this run, and it is the 2026-10-01 pattern recurring: the switch is latched off through a silence that is real.** At 2026-10-07T16:55Z `.mule/.heartbeat` still read 2026-10-06T20:29:30Z — **20h**, with zero `2026-10-07` rows in `.mule/log.md` — but `.mule/.alarmed` (set by #120) short-circuits `check_pulse` at `tools/mule.sh:748`, so nothing can file again until a `done` land (mule.sh:578) or a `checked` run (mule.sh:540) calls `beat`, which is the only thing that clears the latch (mule.sh:737). Recipe ticks are committing today (`becf3cd` 11:25:35Z, `dde64e0` 11:38:44Z) but a recipe verdict is neither of those two paths, so they neither beat nor un-latch — measured today, not assumed. The `check_pulse`-before-`beat` ordering fix (move the call at `tools/mule.sh:633` to after the tick's `beat`) remains escalated above, still yours, still unactioned; #120 is the same arithmetic as #119, not fresh evidence for it either way.
   **#93 at 33 days.** `git rev-list --left-right --count origin/main...HEAD` = **4 / 697** today (4 / 679 on 10-06), and `git merge-base --is-ancestor 24ad5b8 origin/main` still fails — the rewrite is entirely on `mule/queue` / `origin/mule/queue`. The version-drift detail recorded on 2026-10-06 (the comment says `v1.4.0-rc.5` / "v1.3.3 being the newest stable tag"; upstream is `v1.4.0-rc.6` / `v1.3.5`, re-checked via `gh api repos/astarte-platform/astarte/tags` today) folds into that same review rather than becoming a line of its own.
 

→ Read .mule/for-giulio.md [limit=6]
Run complete. Verified the **pairing** surface of docs-sync (routes vs handlers, then all seven endpoints' status codes and request/response schemas against the code).

**Result: the pairing OpenAPI spec is accurate — no `docs/api` task to propose, nothing appended to `todo.md`.** Route lists match 7/7; every documented status code maps to a real code path; schemas (`PairingInfo`, `VerifyValid/Invalid` shared shapes, `ErrorDetail`, field-error envelopes) and examples all match the wire (e.g. `verifyInvalidResponse.details` marshals `null`, `formatDateTime` space-separated timestamps, `APICompatVersion "1.2.2"`). So I stopped, per the recipe.

**One genuine `docs/site` contradiction found** (site prose is on the never-touch list, so escalated to `.mule/for-giulio.md`): `pairing-and-security.md:97` says "Pairing endpoints: per-IP and per-device token buckets" — but register is per-IP only, credentials is the only per-IP-and-per-device endpoint, and unregister/info/verify/health are not rate-limited at all.

No code changes, no git operations.

```
Done: docs-sync pairing run — spec accurate (no task lines); escalated the site rate-limiting contradiction to for-giulio.md
Files: .mule/for-giulio.md
Verified: route comparison (rg Mount handlers vs documented paths) -> pass; status/schema check of 7 endpoints vs internal/pairing/http.go + service.go -> pass
Unsure: nothing
Follow-ups: none
```
