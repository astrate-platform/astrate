slug: recipe-docs-sync
verdict: proposed
at:  4091bbf
ran: 2026-09-16T11:55:27Z on DietPi in 1088s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls docs/api/
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

$ tail -20 .mule/todo.md
- [x] docs-sync-ae-write-405 [auto]: add the missing `405` response to the nine data-write ops in docs/api/astarte_appengine_api.yaml — PUT/POST/DELETE `/devices/{device}/interfaces/{interface}/{path}`, the three `/devices-by-alias/{alias}` twins, and the three `/groups/{group}/devices/{device}` twins. `writeError` emits `405 Method Not Allowed` with "Cannot write to device owned resource" (`ErrNotServerOwned`) and "Cannot write to read-only resource" (`ErrNotAProperty`) at internal/appengine/http.go:673-676, reachable from putData/putDataByAlias/putDataInGroup (http.go:292-357) and deleteData/deleteDataByAlias/deleteDataInGroup (http.go:311-377); `"405"` currently appears nowhere in the spec (those ops document only 200/204, 400, 401, 404, 500). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-write-value-422 [auto]: add the missing `422` response to the six PUT/POST data-write ops in docs/api/astarte_appengine_api.yaml (device-scoped, by-alias, and in-group at `/interfaces/{interface}/{path}`). A value larger than `maxValueBytes` (64 KiB, internal/appengine/http.go:24) is answered `422` "Value size exceeds size limits" by putData/putDataByAlias/putDataInGroup (http.go:298-301, 327-329, 346-348), but the specs for PUT/POST (yaml:686-695, 726-736, 453-463, 499-509, 1176-1186, 1218-1228) list only 200/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-post-groups-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups` in docs/api/astarte_appengine_api.yaml — `Service.CreateGroup` returns `ErrGroupAlreadyExists` (internal/appengine/service.go:495, mapped at http.go:660-661) which the handler answers as `409` "Group already exists"; the spec (yaml:821-839) documents only 201/400/401/422/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] payload-longinteger-fraction-quantize [auto]: `jsonInt64` (pkg/payload/json.go:255-271) accepts fractional longinteger tokens that the ParseFloat fallback silently rounds onto an integer inside its own documented-safety window |x| < 2^53 (json.go:250-254 claims values there never round, "silently corrupting the value") — `{"v":4503599627370496.5}` decodes to 4503599627370496, `{"v":6000000000000000.9}` to 6000000000000001, and `{"v":5e-400}` underflows to 0, all without error, while the contract says fraction forms are "accepted only when integral". Verify the token is an exact decimal integer (`math/big.Rat` SetString + IsInt) on the float fallback — keeps `5.0`, `1e3`, `0.5e1` accepted, rejects the rounding cases as ReasonTypeMismatch; plain-integer ParseInt fast path and the > 2^53 window check unchanged (BSON unaffected: longinteger accepts int64/int32 only). Add failing rows to TestJSONProfile.
- [x] payload-json-malformed-t-not-tolerated [auto]: pin in pkg/payload/payload_test.go that a malformed `t` is rejected with ReasonBadTimestamp even on a non-explicit mapping — `decodeJSONEnvelope` (json.go:50-66) parses `t` before the facade learns the explicit-timestamp policy (payload.go:52-54), so only a *valid* `t` is "tolerated and ignored" while `{"v":1,"t":"yesterday"}` with explicit=false is still a rejection; only the valid-t-ignored row and the BSON "t wrong type" twin are asserted today, so a regression implementing the literal leniency reading (ignore a bad t on non-explicit mappings) passes every suite. Add one table row to TestJSONProfile.
- [x] docs-sync-ae-post-group-devices-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml — re-adding a member the group already holds answers `409` "Device already in group" (`Service.AddGroupDevice` returns `ErrDeviceAlreadyInGroup`, internal/appengine/service.go:636-638, mapped at http.go:662-663), but the spec (yaml:967-977) lists only 201/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-patch-merge-patch-media-type [auto]: fix the three PATCH request bodies in docs/api/astarte_appengine_api.yaml — PATCH /devices/{device} (yaml:176), PATCH /devices-by-alias/{alias} (yaml:260), PATCH /groups/{group}/devices/{device} (yaml:1023) declare `application/json; charset=utf-8`, but all three handlers enforce `Content-Type: application/merge-patch+json` exactly and reject every other value — 500 on the device and group paths (internal/appengine/http.go:209-211, 521-523), 400 on the alias path (http.go:230-232). Change the requestBody media type to `application/merge-patch+json` on all three. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
- [x] config-fail-loud-dev-mode [auto]: `applyEnv` silently discards a malformed `ASTRATE_MQTT_INSECURE_DEV_MODE` value (internal/config/config.go:255-257) — `strconv.ParseBool` error is discarded, the field becomes false, and TLS requirements are silently re-enabled (or a misleading "tls_cert_file required" error surfaces via validate). Return the parse error on bad input; add a failing test that asserts `ASTRATE_MQTT_INSECURE_DEV_MODE=yes` causes Load to error.
- [x] config-key-resolvers-test [auto]: add container-free unit tests for `HousekeepingKeys()` and `RealmJWTPublicKey()` (internal/config/config.go:383-411) — both file-reading helpers have zero test coverage. Cover: happy path (inline PEM + file PEM → concatenated/trimmed), file-not-found error (assert error message includes the path), inline-over-file priority (`RealmJWTPublicKey` returns the inline block when both are set), and TrimSpace behaviour (leading/trailing newlines stripped). `t.TempDir()` + `os.WriteFile` for fixtures; no Docker.
- [x] docs-sync-hk-delete-gating-responses [auto]: add the missing `405` and `422` responses to `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — `writeError` answers 405 "Realm deletion disabled" for `ErrDeletionDisabled` and 422 `{"errors":{"error_name":["connected_devices_present"]}}` for `ErrConnectedDevicesPresent` (internal/housekeeping/http.go:233-249), both reachable from `Service.DeleteRealm` (service.go:248-266), but the spec lists only 204/401/403/404/500 (yaml:198-208). Its sole `405` reference sits on PATCH (yaml:182-183), which can never return it (`UpdateRealm`, service.go:196-216, has no deletion-gating path) — re-point 405 to DELETE (reusing the MethodNotAllowed component) and add 422 with the FieldErrors shape and a connected_devices_present example. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_housekeeping_api.yaml (yaml:417) — it shows `detail: "housekeeping: validation failed: realm name must not be empty"`, but `validationDetail` strips the `housekeeping: validation failed: ` prefix on the wire (internal/housekeeping/http.go:251-258) and the real blank-name message is `"realm_name can't be blank"` (service.go:138). Change the example to `detail: "realm_name can't be blank"`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

- [x] auth-iat-not-required-test [auto]: pin in internal/auth/jwt_test.go that a token with a *future* `iat` still verifies — the "`iat` is not required — upstream parity" rule (jwt.go:104-105, no `jwt.WithIssuedAt` in the parser, jwt.go:115-118) has zero coverage, so a regression that starts validating `iat` or `WithIssuedAt` would silently reject tokens from a skewed-clock issuer and pass every existing suite. Sign RS256 with `iat` = now+1h and no `exp`/`nbf`, assert `Verify` succeeds.
- [x] auth-empty-claims-403-test [auto]: pin in internal/auth/middleware_test.go that a validly signed token carrying NO Astarte claims is `403` on a guarded route, not `401` — the authenticated-but-grant-less distinction (middleware.go:149-151, `Authorizes` on a nil grants map → false → `WriteForbidden`) is currently untested: every 403 fixture has a claim on the wrong surface, so a regression mapping grant-less tokens to 401 (or skipping the forbidden path) passes the suite. Sign RS256 with empty `jwt.MapClaims`, hit the pairing route, expect 403 + golden `envelope_403.json`.
- [x] auth-cache-default-size-test [auto]: pin in internal/auth/cache_test.go that `NewCache(0)` and `NewCache(-1)` fall back to `DefaultCacheSize` (cache.go:29-31) — the `size < 1` branch exists purely so `lru.New` cannot panic and has zero coverage; a regression dropping the fallback (panic on `lru.New` with a non-tick size, or a silently tiny cache) would pass every test. Construct with both sizes, assert no panic and a working `Verify` with `lru.Len() == 0`.
- [x] hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto] — already resolved and re-queued by the refill: commit 015a706 chose the upstream-parity behaviour (0 ≡ unset, not invalid, not a literal 0) and pinned it — create folds 0 → nil with the container-free TestCreateRealmZeroRetentionUnsets (internal/housekeeping/service_test.go:186) and PATCH 0 → ClearRetention in the integration suite (internal/housekeeping/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).
- [x] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication is "via Authorization header or `?token=` query parameter" (yaml:293-294), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] astarteapi-metadata-envelope-test [auto]: add a container-free golden test for `WriteDataWithMetadata` (pkg/astarteapi/envelope.go:154-158) pinning the two wire branches — nil metadata omits the `metadata` key (`{"data": v}` only, via the `omitempty` at envelope.go:151), populated metadata emits `{"data": v, "metadata": {...}}` — with 200 + the package `Content-Type`, matching the TestGoldenEnvelopes style. Today the non-nil branch is only covered by the integration-tagged formats_parity_test.go (runs on the Legion Go, not the Pi gate) and the nil-omission branch is covered nowhere (internal/appengine/data.go:262-268 always sets metadata), so a regression dropping the omitempty passes every suite on this box.
- [x] astarteapi-multikey-sorted-golden [auto]: extend TestGoldenEnvelopes (pkg/astarteapi/envelope_test.go:101-121) with multi-key `WriteFieldErrors` and `WriteRawErrors` maps asserting the sorted key order on the wire — the docs claim "keys are emitted in sorted order" (envelope.go:170, 184) but every fixture uses a single key, so a change to the marshalling path that stopped sorting (e.g. switching encoders) passes the suite.

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/pairing/ docs/api/astarte_pairing_api.yaml 2>/dev/null | head -60
internal/pairing/http.go:81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/pairing/http.go:86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
internal/pairing/http.go:88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
internal/pairing/http.go:90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
internal/pairing/http.go:344:	header := r.Header.Get("Authorization")

This is the docs-sync recipe. I'll pick one surface per run — appengine and housekeeping were just covered in the queue, so I'll check the untouched ones. Let me survey the surfaces.
$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/realmmanagement/ internal/realm_management/ docs/api/astarte_realm_management_api.yaml 2>/dev/null | head -80
(no output)

$ ls internal/ | rg -i 'realm|management|native|housekeep|appengine|pairing'
appengine
housekeeping
pairing
realm

$ rg -n '^\s{2}/' docs/api/astarte_pairing_api.yaml
27:  /pairing/v1/{realm}/health:
70:  /pairing/v1/{realm}/version:
106:  /pairing/v1/{realm}/agent/devices:
158:  /pairing/v1/{realm}/agent/devices/{deviceID}:
181:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
235:  /pairing/v1/{realm}/devices/{deviceID}:
273:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:

$ rg -n 'HandleFunc|mux\.' internal/pairing/http.go
74:// Mount registers the pairing routes on mux. Paths are wire-frozen
81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
82:	mux.Handle("POST /pairing/v1/{realm}/agent/devices",
84:	mux.Handle("DELETE /pairing/v1/{realm}/agent/devices/{deviceID}",
86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",

$ rg -rn 'HandleFunc' internal/realm/ | head -40
(no output)

Pairing looks clean. Let me check the full route list in the pairing handler and the realm-management surface.
$ rg -n '^\s{2}/' docs/api/astarte_realm_management_api.yaml
27:  /realmmanagement/v1/{realm}/interfaces:
138:  /realmmanagement/v1/{realm}/interfaces/{name}:
167:  /realmmanagement/v1/{realm}/interfaces/{name}/{major}:
279:  /realmmanagement/v1/{realm}/triggers:
349:  /realmmanagement/v1/{realm}/triggers/{name}:
401:  /realmmanagement/v1/{realm}/config/auth:
470:  /realmmanagement/v1/{realm}/config/device_registration_limit:
501:  /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention:
536:  /realmmanagement/v1/{realm}/version:
569:  /realmmanagement/v1/{realm}/devices/{device}:
595:  /realmmanagement/v1/{realm}/policies:
665:  /realmmanagement/v1/{realm}/policies/{name}:

$ rg -n 'agent/devices' docs/api/astarte_pairing_api.yaml
106:  /pairing/v1/{realm}/agent/devices:
158:  /pairing/v1/{realm}/agent/devices/{deviceID}:

$ rg -rn 'HandleFunc|Handle\(' internal/realm/ | head -60
internal/realm/http.go:	mux.n"GET /realmmanagement/v1/{realm}/interfaces", h(a.listInterfaces))
internal/realm/http.go:	mux.n"POST /realmmanagement/v1/{realm}/interfaces", h(a.installInterface))
internal/realm/http.go:	mux.n"GET /realmmanagement/v1/{realm}/interfaces/{name}", h(a.listInterfaceMajors))
internal/realm/http.go:	mux.n"GET /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.getInterface))
internal/realm/http.go:	mux.n"PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.updateInterface))
internal/realm/http.go:	mux.n"DELETE /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.deleteInterface))
internal/realm/http.go:	mux.n"GET /realmmanagement/v1/{realm}/triggers", h(a.listTriggers))
internal/realm/http.go:	mux.n"POST /realmmanagement/v1/{realm}/triggers", h(a.createTrigger))
internal/realm/http.go:	mux.n"GET /realmmanagement/v1/{realm}/triggers/{name}", h(a.getTrigger))
internal/realm/http.go:	mux.n"DELETE /realmmanagement/v1/{realm}/triggers/{name}", h(a.deleteTrigger))
internal/realm/http.go:	mux.n"GET /realmmanagement/v1/{realm}/config/auth", h(a.getAuth))
internal/realm/http.go:	mux.n"PUT /realmmanagement/v1/{realm}/config/auth", h(a.putAuth))
internal/realm/http.go:	mux.n"GET /realmmanagement/v1/{realm}/config/device_registration_limit", h(a.getRegistrationLimit))
internal/realm/http.go:	mux.n"GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention", h(a.getDatastreamMaximumStorageRetention))
internal/realm/http.go:	mux.n"GET /realmmanagement/v1/{realm}/version", h(a.getVersion))
internal/realm/http.go:	mux.n"DELETE /realmmanagement/v1/{realm}/devices/{device}", h(a.deleteDevice))
internal/realm/http.go:	mux.n"GET /realmmanagement/v1/{realm}/policies", h(a.listPolicies))
internal/realm/http.go:	mux.n"POST /realmmanagement/v1/{realm}/policies", h(a.createPolicy))
internal/realm/http.go:	mux.n"GET /realmmanagement/v1/{realm}/policies/{name}", h(a.getPolicy))
internal/realm/http.go:	mux.n"DELETE /realmmanagement/v1/{realm}/policies/{name}", h(a.deletePolicy))

$ ls docs/api/ astrate-mule/ 2>/dev/null; echo ---; ls /root/astrate-mule/.mule/; echo ---; tail -15 /root/astrate-mule/.mule/current-task.md 2>/dev/null
docs/api/:
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

$ gh issue view 123 --repo astrate/astrate --json title,body,labels 2>&1 | head -20
GraphQL: Could not resolve to a Repository with the name 'astrate/astrate'. (repository)

Let me examine the pairing surface closely — the doc lists `/version` which I don't see in the code.
→ Read internal/pairing/http.go [offset=74, limit=40]
→ Read docs/api/astarte_pairing_api.yaml [offset=106, limit=180]

$ rg -n 'version' internal/pairing/http.go | head -20
198:	Version string `json:"version"`

$ rg -rn 'handleVersion\|GET.*version\|/version' internal/pairing/ | head -20
(no output)

The YAML has a `/pairing/v1/{realm}/version` path at line 70 that doesn't appear in the code's `Mount` function. Let me verify.
→ Read docs/api/astarte_pairing_api.yaml [offset=65, limit=50]
This run I'll pick the **pairing** surface (appengine and housekeeping were covered in recent runs). Found one candidate already: YAML documents `GET /pairing/v1/{realm}/version` but `Mount` has no such route. Let me check the handlers against the documented endpoints.
→ Read docs/api/astarte_pairing_api.yaml [limit=40, offset=27]
→ Read internal/pairing/http.go [limit=160, offset=109]

$ rg -n 'handleHealth' internal/pairing/http.go
81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
281:// handleHealth serves the realm-scoped health probe (#71): 200 when the
284:func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {

→ Read internal/pairing/http.go [offset=269, limit=99]
Now let me look at `writeServiceError` and the verify endpoint YAML, plus the health handler status codes.
→ Read docs/api/astarte_pairing_api.yaml [offset=286, limit=120]

$ rg -rn 'pairing/v1.*version|/version' internal/ cmd/ pkg/ --glob '*.go' | rg -i 'pair|version' | head -20
internal/realm/http.go:	mux.Handle("GET /realmmanagement/v1/{realm}n", h(a.getVersion))
internal/observability/compat.go:// VersionHandler answers GET ...n with the upstream envelope
internal/observability/compat.go:// MountVersionCompat registers GET /{service}n for one service.
internal/observability/compat.go:	mux.HandleFunc("GET /"+service+"n", VersionHandler(version))
internal/observability/compat_version_test.go:		resp, err := http.Get(srv.URL + "/" + svc + "n")
internal/observability/compat_version_test.go:			t.Fatalf("%sn: %v", svc, err)
internal/observability/compat_version_test.go:			t.Errorf("%sn = %d, want 200", svc, resp.StatusCode)
internal/observability/compat_version_test.go:			t.Errorf("%sn body = %s", svc, body)
internal/observability/compat_version_test.go:			t.Errorf("%sn missing Content-Type", svc)
internal/observability/compat_version_test.go:	mux.HandleFunc("GET /v1/{realm}n", VersionHandler("x"))
internal/observability/compat_version_test.go:	resp, err := http.Get(srv.URL + "/v1/somerealmn")
internal/observability/compat_version_test.go:		t.Fatalf("/v1/somerealmn: %v", err)
internal/observability/compat_version_test.go:		t.Errorf("/v1/somerealmn = %d, want 200", resp.StatusCode)
internal/observability/compat_version_test.go:		t.Errorf("/v1/somerealmn body = %s", body)

$ rg -n 'DataEnvelopeHealthStatus' docs/api/astarte_pairing_api.yaml
47:                $ref: "#/components/schemas/DataEnvelopeHealthStatus"
519:    DataEnvelopeHealthStatus:

The bash output is being mangled (globs in pattern). Let me read the compatibility version handler and the realm http.go directly.
→ Read internal/observability/compat.go
✱ Grep "MountVersionCompat|MountServiceCompat|VersionHandler" 19 matches
→ Read docs/api/astarte_pairing_api.yaml [limit=160, offset=406]
The `/version` route does exist — it's mounted in `cmd/astrate/main.go:398`, not in pairing's `Mount`. Let me check the schemas and the main.go wiring.
→ Read cmd/astrate/main.go [limit=40, offset=375]
→ Read pkg/astarteapi/envelope.go
Let me verify the error envelope shapes and check the health 503 body against what `WriteError` actually produces.
Let me read the rest of the pairing YAML (responses section) to compare exact shapes.
→ Read docs/api/astarte_pairing_api.yaml [offset=564]

$ rg -n 'func \(s \*service\) (Register|Unregister|Credentials|Info|VerifyCredentials|Health)' internal/pairing/*.go
(no output)

The pairing surface looks mostly consistent. Let me check the service layer for what errors `Unregister` and `Info` can emit, since DELETE unregister documents 204/401/403/404/500 only.
$ rg -n 'ErrInvalidHWID|ErrRegistrationLimitReached|ErrAlreadyRegistered|ErrInvalidPayloadFormat|ErrUnauthorized|ErrInhibited' internal/pairing/*.go
internal/pairing/http.go:307:	case errors.Is(err, ErrInvalidHWID):
internal/pairing/http.go:310:	case errors.Is(err, ErrInvalidPayloadFormat):
internal/pairing/http.go:313:	case errors.Is(err, ErrAlreadyRegistered):
internal/pairing/http.go:316:	case errors.Is(err, ErrRegistrationLimitReached):
internal/pairing/http.go:322:	case errors.Is(err, ErrUnauthorized):
internal/pairing/http.go:324:	case errors.Is(err, ErrInhibited):
internal/pairing/service.go:44:	// ErrInvalidHWID reports a hw_id that is not a 22-character unpadded
internal/pairing/service.go:46:	ErrInvalidHWID = errors.New("pairing: invalid hw_id")
internal/pairing/service.go:47:	// ErrInvalidPayloadFormat reports an initial_payload_format outside
internal/pairing/service.go:49:	ErrInvalidPayloadFormat = errors.New("pairing: invalid initial_payload_format")
internal/pairing/service.go:50:	// ErrAlreadyRegistered reports re-registration of a device that has
internal/pairing/service.go:53:	ErrAlreadyRegistered = errors.New("pairing: device already registered")
internal/pairing/service.go:54:	// ErrRegistrationLimitReached reports the realm's
internal/pairing/service.go:57:	ErrRegistrationLimitReached = errors.New("pairing: device registration limit reached")
internal/pairing/service.go:58:	// ErrUnauthorized is the uniform device-authentication failure: unknown
internal/pairing/service.go:61:	ErrUnauthorized = errors.New("pairing: unauthorized")
internal/pairing/service.go:62:	// ErrInhibited reports a device blocked by credentials_inhibited (403).
internal/pairing/service.go:63:	ErrInhibited = errors.New("pairing: credentials request inhibited")
internal/pairing/service.go:172:// secret; afterwards it fails with ErrAlreadyRegistered. initialFormat is
internal/pairing/service.go:177:		return "", fmt.Errorf("%w: %v", ErrInvalidHWID, err)
internal/pairing/service.go:180:		return "", fmt.Errorf("%w: %q", ErrInvalidPayloadFormat, initialFormat)
internal/pairing/service.go:196:			return "", ErrRegistrationLimitReached
internal/pairing/service.go:211:			return "", fmt.Errorf("%w: %s", ErrAlreadyRegistered, hwID)
internal/pairing/service.go:243:// device by credentials secret (uniform ErrUnauthorized on any mismatch),
internal/pairing/service.go:253:		return "", ErrInhibited
internal/pairing/service.go:330:		return nil, ErrInhibited
internal/pairing/service.go:378:// device, wrong secret) return the same ErrUnauthorized after a comparable
internal/pairing/service.go:384:		return nil, nil, ErrUnauthorized
internal/pairing/service.go:389:		return nil, nil, ErrUnauthorized
internal/pairing/service.go:397:		return nil, nil, ErrUnauthorized
internal/pairing/service.go:405:		return nil, nil, ErrUnauthorized
internal/pairing/service.go:408:		return nil, nil, ErrUnauthorized
internal/pairing/service_test.go:223:	if _, err := svc.Register(ctx, "test", hwID, ""); !errors.Is(err, ErrAlreadyRegistered) {
internal/pairing/service_test.go:224:		t.Errorf("post-credentials re-register: got %v, want ErrAlreadyRegistered", err)
internal/pairing/service_test.go:274:	if _, err := svc.Register(ctx, "test", hwID, ""); !errors.Is(err, ErrAlreadyRegistered) {
internal/pairing/service_test.go:275:		t.Errorf("post-credentials re-register: got %v, want ErrAlreadyRegistered", err)
internal/pairing/service_test.go:289:		if _, err := svc.Register(context.Background(), "test", hwID, ""); !errors.Is(err, ErrInvalidHWID) {
internal/pairing/service_test.go:290:			t.Errorf("Register(%q): got %v, want ErrInvalidHWID", hwID, err)
internal/pairing/service_test.go:318:	if _, err := svc.Register(ctx, "test", randomDeviceID(t), "cbor"); !errors.Is(err, ErrInvalidPayloadFormat) {
internal/pairing/service_test.go:319:		t.Errorf("invalid format: got %v, want ErrInvalidPayloadFormat", err)
internal/pairing/service_test.go:332:	if _, err := svc.Register(ctx, "test", randomDeviceID(t), ""); !errors.Is(err, ErrRegistrationLimitReached) {
internal/pairing/service_test.go:333:		t.Errorf("over-limit Register: got %v, want ErrRegistrationLimitReached", err)
internal/pairing/service_test.go:365:	if _, err := svc.Credentials(ctx, "test", hwID, secret, deviceCSR(t), testIP); !errors.Is(err, ErrUnauthorized) {
internal/pairing/service_test.go:366:		t.Errorf("old secret after unregister: got %v, want ErrUnauthorized", err)
internal/pairing/service_test.go:450:		if err := call(); !errors.Is(err, ErrUnauthorized) {
internal/pairing/service_test.go:451:			t.Errorf("%s: got %v, want ErrUnauthorized", name, err)
internal/pairing/service_test.go:467:	if _, err := svc.Credentials(ctx, "test", hwID, secret, deviceCSR(t), testIP); !errors.Is(err, ErrInhibited) {
internal/pairing/service_test.go:468:		t.Errorf("inhibited: got %v, want ErrInhibited", err)
internal/pairing/service_test.go:523:	if _, err := svc.Info(ctx, "test", hwID, "wrong"); !errors.Is(err, ErrUnauthorized) {
internal/pairing/service_test.go:524:		t.Errorf("wrong secret on info: got %v, want ErrUnauthorized", err)
internal/pairing/service_test.go:633:		if _, err := svc.VerifyCredentials(ctx, "test", hwID, secret, certPEM); !errors.Is(err, ErrInhibited) {
internal/pairing/service_test.go:634:			t.Errorf("inhibited verify: got %v, want ErrInhibited", err)

→ Read internal/pairing/service.go [limit=130, offset=240]
→ Read internal/pairing/service.go [limit=70, offset=370]

$ rg -n 'func \(a \*API\) requireAgent|requireAgent' internal/pairing/*.go
internal/pairing/http.go:46:	requireAgent func(http.Handler) http.Handler
internal/pairing/http.go:68:		requireAgent: mw.RequireRealm(auth.ClaimPairing),
internal/pairing/http.go:83:		a.requireAgent(http.HandlerFunc(a.handleRegister)))
internal/pairing/http.go:85:		a.requireAgent(http.HandlerFunc(a.handleUnregister)))

Let me check the `Unregister` service function to see what statuses DELETE can produce.
$ rg -n 'func \(s \*Service\) Unregister' internal/pairing/service.go
230:func (s *Service) Unregister(ctx context.Context, realmName, deviceIDStr string) error {

→ Read internal/pairing/service.go [limit=30, offset=225]

$ rg -n 'Cause(Expired|Invalid|Revoked)\s*=' internal/pairing/*.go
internal/pairing/service.go:72:	CauseExpired = "EXPIRED"
internal/pairing/service.go:75:	CauseInvalid = "INVALID"
internal/pairing/service.go:78:	CauseRevoked = "REVOKED"

Let me verify the cause constants and double-check the flow C endpoints' reachable status codes.
$ rg -n 'GetRealmByName|store.ErrNotFound' internal/pairing/service.go | head -20
29:	GetRealmByName(ctx context.Context, name string) (*store.Realm, error)
183:	realm, err := s.st.GetRealmByName(ctx, realmName)
229:// credential trail). store.ErrNotFound is returned for unknown devices.
233:		return fmt.Errorf("%w: device %q", store.ErrNotFound, deviceIDStr)
235:	realm, err := s.st.GetRealmByName(ctx, realmName)
370:	if _, err := s.st.GetRealmByName(ctx, realm); err != nil {
371:		return err // store.ErrNotFound → 404
386:	realm, err := s.st.GetRealmByName(ctx, realmName)
387:	if errors.Is(err, store.ErrNotFound) {
395:	if errors.Is(err, store.ErrNotFound) {

$ wc -l .mule/todo.md
179 .mule/todo.md

I have enough to assess the pairing surface. Let me confirm the exact reachable error paths once more, then check the end of todo.md for appending.

The key finding: `Service.Info` and `Service.Credentials`/`VerifyCredentials` only ever return `ErrUnauthorized` (via `authenticateDevice`, pairing/service.go:380-411 which folds parse/unknown/secret failures all into 401), `ErrInhibited` (403), or an opaque error (500) — never `store.ErrNotFound`. And the realm-scoped `/version` route is a static 200 handler (main.go:398). So several documented responses on flow C endpoints are unreachable.
→ Read .mule/todo.md [offset=170]
Now let me run the config-key sweep from the recipe and check the register service path for any `ErrInhibited`.

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
internal/flow/blocks/container/docker.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/info.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/schema.go:ASTRATE_FLOW_CONFIG
internal/store/crypto.go:ASTRATE_MASTER_KEY
internal/store/crypto.go:ASTRATE_MASTER_KEY_FILE
internal/store/store_test.go:ASTRATE_TEST_DSN
internal/testutil/pg.go:ASTRATE_TEST_DSN

→ Read internal/pairing/service.go [limit=60, offset=165]

$ rg -n 'HOUSEKEEPING' internal/config/config.go
279:		"ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
280:		"HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
297:	if v, ok := os.LookupEnv("ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"); ok {

$ rg -n 'HOUSEKEEPING' docs/site/configuration-reference.md; echo "exit=$?"
exit=1

$ tail -20 .mule/for-giulio.md 2>/dev/null; echo "=== lines:"; wc -l .mule/for-giulio.md
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
=== lines:
284 .mule/for-giulio.md

The config sweep shows two key candidates. Let me verify they're real env vars in `config.go` and check `for-giulio.md` and the docs page.
$ rg -n '240|279|280|297' internal/config/config.go | head; echo ---; rg -n 'storage|retention|deletion|HOUSEKEEP' docs/site/configuration-reference.md | head -20
---
85:## `[storage]` — Retention
89:| `retention` | duration | `""` (disabled) | — | No | Global drop-chunks retention. Empty/0 disables it; per-endpoint TTL still applies. Example: `"8760h"` (one year). |

Let me confirm the two API paths' documented 4xx responses and check the config reference file more broadly.
→ Read internal/config/config.go [limit=45, offset=270]

$ rg -n 'ASTRATE_|Env|environment' docs/site/configuration-reference.md | head -20
3:Astrate reads one TOML file (`astrate -config path.toml`) with `ASTRATE_*` environment overrides. Precedence: **built-in defaults < TOML file < environment variables**.
11:- `database.dsn` (or `ASTRATE_DATABASE_DSN`) — PostgreSQL/TimescaleDB connection string.
19:built-in defaults  <  TOML file  <  ASTRATE_* environment
22:Every value below is the built-in default unless noted. Environment overrides exist for operationally critical fields, named `ASTRATE_<SECTION>_<FIELD>` (e.g. `ASTRATE_DATABASE_DSN`, `ASTRATE_HTTP_ADDR`).
28:| Key | Type | Default | Env override | Required | Description |
30:| `dsn` | string | `""` | `ASTRATE_DATABASE_DSN` | **Yes** | PostgreSQL + TimescaleDB connection string. |
34:| Key | Type | Default | Env override | Required | Description |
36:| `addr` | string | `":8080"` | `ASTRATE_HTTP_ADDR` | No | Bind address for the single REST listener (pairing, realm management, housekeeping, appengine, native endpoints). |
37:| `tls_cert_file` | string | `""` | `ASTRATE_HTTP_TLS_CERT_FILE` | No | Optional in-binary TLS cert. Leave empty to terminate TLS at a reverse proxy. |
38:| `tls_key_file` | string | `""` | `ASTRATE_HTTP_TLS_KEY_FILE` | No | Optional in-binary TLS key. Must be set together with `tls_cert_file`. |
39:| `cors_allowed_origins` | list[string] | `[]` | `ASTRATE_HTTP_CORS_ALLOWED_ORIGINS` (comma-separated) | No | Browser origins allowed for CORS. `"*"` allows any; empty disables CORS. The Astarte Dashboard SPA needs its origin here (e.g. `"http://localhost:4040"`). |
43:| Key | Type | Default | Env override | Required | Description |
45:| `addr` | string | `":8883"` | `ASTRATE_MQTT_ADDR` | No | Bind address for the mTLS MQTT listener. |
46:| `tls_cert_file` | string | `""` | `ASTRATE_MQTT_TLS_CERT_FILE` | **Yes*** | Broker server cert. Required unless `insecure_dev_mode` is true. |
47:| `tls_key_file` | string | `""` | `ASTRATE_MQTT_TLS_KEY_FILE` | **Yes*** | Broker server key. Required unless `insecure_dev_mode` is true. |
48:| `insecure_dev_mode` | bool | `true` | `ASTRATE_MQTT_INSECURE_DEV_MODE` | No | Bind a plaintext listener on `dev_addr` that trusts the claimed client ID. **Never enable in production.** |
50:| `session_store_path` | string | `"sessions.db"` | `ASTRATE_MQTT_SESSION_STORE_PATH` | No | bbolt/pebble file for MQTT session persistence. On a read-only container FS, must point at the writable session volume. |
51:| `advertised_url` | string | `""` | `ASTRATE_MQTT_ADVERTISED_URL` | No | Broker URL handed to devices by the pairing info endpoint. Empty derives `"mqtts://<addr>"`; set when devices reach the broker by another host. |
58:| Key | Type | Default | Env override | Required | Description |
60:| `shards` | int | `16` | `ASTRATE_ENGINE_SHARDS` | No | Number of ordered ingestion shards. Each shard processes messages strictly in order. |

Let me check how the config reference doc presents env vars generally.
→ Read docs/site/configuration-reference.md [limit=60]
→ Read docs/site/configuration-reference.md [offset=61]

$ rg -rn 'HOUSEKEEPING|realm_deletion|default_datastream' docs/site/ internal/config/config.example.toml | head
(no output)

Now let me verify the two config keys against the full docs site and the example TOML, and re-confirm the flow C endpoints' reachable statuses.
→ Read internal/pairing/http.go [limit=40, offset=40]

$ rg -rn 'HOUSEKEEPING' docs/ | head; echo "exit=$?"
exit=0

The config sweep found two env vars in code that `docs/site/configuration-reference.md` doesn't document. Let me confirm they aren't referenced anywhere in `docs/`, and check the for-giulio.md format.
→ Read .mule/for-giulio.md [limit=25]
The config-reference key gap is already escalated (for-giulio 2026-09-12 line, both `ASTRATE_HOUSEKEEPING_*` keys). Let me check for prior pairing docs tasks to avoid duplication.

$ rg -n 'docs-sync.*pair|pairing.*yaml|pairing.*api' .mule/todo.md | head -20
109:- [x] docs-sync-pairing-health-path [auto]: add the undocumented `GET /pairing/v1/{realm}/health` route to docs/api/astarte_pairing_api.yaml — it exists in code (internal/pairing/http.go:81) and returns 200 `{"data":{"status":"ok"}}` for an existing realm, 404 for an unknown realm, 503 when unhealthy, unauthenticated and unrate-limited. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
110:- [x] docs-sync-pairing-register-404 [auto]: document the missing `404` (DeviceNotFound) on `POST /pairing/v1/{realm}/agent/devices` in docs/api/astarte_pairing_api.yaml — the handler returns it when `GetRealmByName` fails for an unknown realm (internal/pairing/service.go:183-186 → http.go:326-327), but the spec lists only 201/400/401/403/422/429/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
147:- [x] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmmanagement/health`, `GET /pairing/health` — in docs/api/astrate_native_api.yaml, and widen the shared `DataEnvelopeHealthStatus.status` enum (yaml:301-303) from `[ok]` to `[ok, unhealthy]` — the handler answers 503 `{"data":{"status":"unhealthy"}}` when the readiness probe fails (internal/observability/compat.go:19-34, deviation 18), but the spec documents only 200 with `status: ok` for all three paths. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
148:- [x] docs-sync-pairing-status-enum [auto]: fix the `PairingInfo.status` enum in docs/api/astarte_pairing_api.yaml from `[confirmed, pending, denied, expired]` (yaml:412) to `[confirmed, pending, inhibited]` — the handler never emits `denied` or `expired`: `service.Info` maps only `DeviceStatusInhibited`→"inhibited", `DeviceStatusConfirmed`→"confirmed", and everything else→"pending" (internal/pairing/service.go:296-299, statuses defined in store/devices.go:19-27; upstream-parity per the comment at service.go:282-285). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
149:- [x] docs-sync-pairing-version-endpoint [auto]: add the undocumented `GET /pairing/v1/{realm}/version` route to docs/api/astarte_pairing_api.yaml — it exists in code (cmd/astrate/main.go:398, unauthenticated, upstream-parity per main.go:387-391) and answers 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38); mirror the RM spec's getVersion op (astarte_realm_management_api.yaml:532) without its a_rma security. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
150:- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

No prior task covers these; the version endpoint was just added (with a claimed 404). The config-key gap is already escalated (for-giulio line 23). My pairing-surface findings are three unreachable documented responses. Let me verify my key claim once more — that `authenticateDevice` (the only error path for flow B/C endpoints) never surfaces `store.ErrNotFound` or `ErrInhibited` on the info route.
$ rg -n 'Inhibited|ErrNotFound|ErrUnauthorized' internal/pairing/service.go | sed -n '1,40p'
58:	// ErrUnauthorized is the uniform device-authentication failure: unknown
61:	ErrUnauthorized = errors.New("pairing: unauthorized")
62:	// ErrInhibited reports a device blocked by credentials_inhibited (403).
63:	ErrInhibited = errors.New("pairing: credentials request inhibited")
229:// credential trail). store.ErrNotFound is returned for unknown devices.
233:		return fmt.Errorf("%w: device %q", store.ErrNotFound, deviceIDStr)
243:// device by credentials secret (uniform ErrUnauthorized on any mismatch),
252:	if dev.Status == store.DeviceStatusInhibited {
253:		return "", ErrInhibited
294:	case store.DeviceStatusInhibited:
329:	if dev.Status == store.DeviceStatusInhibited {
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

Confirmed. The pairing surface has three over-documented, unreachable responses. Let me check where new proposal lines go in todo.md.
→ Read .mule/todo.md [limit=12]
→ Read .mule/todo.md [limit=30, offset=13]
→ Read .mule/todo.md [limit=85, offset=43]

$ rg -n 'build:' docs/Makefile | head -5
16:build: sync

Now let me verify the docs build target exists so task lines can reference it accurately.
→ Read .mule/todo.md [offset=177]
The pairing surface checks out. I found two real spec drift clusters (unreachable documented status codes). The config-key gap (two `ASTRATE_HOUSEKEEPING_*` keys) is already escalated in `.mule/for-giulio.md` line 23, so no new for-giulio entry is needed. Appending the proposals:

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -176,4 +176,6 @@
 - [x] hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto] — already resolved and re-queued by the refill: commit 015a706 chose the upstream-parity behaviour (0 ≡ unset, not invalid, not a literal 0) and pinned it — create folds 0 → nil with the container-free TestCreateRealmZeroRetentionUnsets (internal/housekeeping/service_test.go:186) and PATCH 0 → ClearRetention in the integration suite (internal/housekeeping/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).
 - [x] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication is "via Authorization header or `?token=` query parameter" (yaml:293-294), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] astarteapi-metadata-envelope-test [auto]: add a container-free golden test for `WriteDataWithMetadata` (pkg/astarteapi/envelope.go:154-158) pinning the two wire branches — nil metadata omits the `metadata` key (`{"data": v}` only, via the `omitempty` at envelope.go:151), populated metadata emits `{"data": v, "metadata": {...}}` — with 200 + the package `Content-Type`, matching the TestGoldenEnvelopes style. Today the non-nil branch is only covered by the integration-tagged formats_parity_test.go (runs on the Legion Go, not the Pi gate) and the nil-omission branch is covered nowhere (internal/appengine/data.go:262-268 always sets metadata), so a regression dropping the omitempty passes every suite on this box.
 - [x] astarteapi-multikey-sorted-golden [auto]: extend TestGoldenEnvelopes (pkg/astarteapi/envelope_test.go:101-121) with multi-key `WriteFieldErrors` and `WriteRawErrors` maps asserting the sorted key order on the wire — the docs claim "keys are emitted in sorted order" (envelope.go:170, 184) but every fixture uses a single key, so a change to the marshalling path that stopped sorting (e.g. switching encoders) passes the suite.
+- [ ] docs-sync-pairing-deviceendpoints-dead-404-403 [auto]: remove the unreachable `404` DeviceNotFound response from the three device-facing ops in docs/api/astarte_pairing_api.yaml — `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials` (yaml:226-227), `GET /pairing/v1/{realm}/devices/{deviceID}` (yaml:268-269), `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify` (yaml:339-340) — and the `403` Forbidden from the info GET (yaml:266-267). `Service.authenticateDevice` folds every failure (malformed/unknown deviceID, unknown realm, unregistered device, wrong secret) into uniform `ErrUnauthorized` (internal/pairing/service.go:380-411), so `store.ErrNotFound` never reaches `writeServiceError` (http.go:326-327) from those handlers; and `Service.Info` never returns `ErrInhibited` (inhibited devices still read info with 200 + status "inhibited", service.go:286-307, comment at 282-285). Reachable on these three are only 200/201/400/401/422/429/500 plus 403 on credentials and verify (both genuinely raise `ErrInhibited`, service.go:252-254 and 329-331). Keep 404 on the agent endpoints (register yaml:149-150, unregister yaml:176-177) where `GetRealmByName`/`deviceid.Parse` really do return `store.ErrNotFound`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-pairing-version-404-unreachable [auto]: remove the `404` "Unknown realm" response from `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml (yaml:96-104) — the route is served by `observability.VersionHandler` (cmd/astrate/main.go:398, internal/observability/compat.go:38-47), which answers static `200 {"data":"<version>"}` with no realm lookup, so a 404 can never be emitted there. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

Done: docs-sync pairing surface — two unreachable-documented-response tasks proposed, no for-giulio entry needed (config-key gap already escalated 2026-09-12).
Files: .mule/todo.md
Verified: rg route/path comparison (internal/pairing/http.go:77-92 vs docs/api/astarte_pairing_api.yaml) -> pass; handler status audit (http.go, service.go:174-411, main.go:393-398, compat.go:38-47) -> pass; config-key sweep (config.go vs configuration-reference.md) -> no new gaps
Unsure: nothing
Follow-ups: `docs-sync-pairing-deviceendpoints-dead-404-403`, `docs-sync-pairing-version-404-unreachable` appended to .mule/todo.md
