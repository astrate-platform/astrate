slug: recipe-docs-sync
verdict: proposed
at:  c6ab452
ran: 2026-09-17T19:49:54Z on DietPi in 393s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
→ Read .mule/todo.md

$ ls internal/ && echo "---" && ls docs/api/ 2>/dev/null && echo "---" && ls docs/site/ 2>/dev/null
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
---
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
---
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

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*appengine*'
(no output)

$ rg -n '^\s{2}/' docs/api/astrate_appengine_api.yaml
rg: docs/api/astrate_appengine_api.yaml: IO error for operation on docs/api/astrate_appengine_api.yaml: No such file or directory (os error 2)

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc|\.Handle\(' internal/appengine/
internal/appengine/http.go:43:	mux.Handle("GET "+base+"/devices", h(a.listDevices))
internal/appengine/http.go:44:	mux.Handle("GET "+base+"/stats/devices", h(a.devicesStats))
internal/appengine/http.go:45:	mux.Handle("GET "+base+"/devices/{device}", h(a.getDevice))
internal/appengine/http.go:46:	mux.Handle("PATCH "+base+"/devices/{device}", h(a.patchDevice))
internal/appengine/http.go:47:	mux.Handle("GET "+base+"/devices-by-alias/{alias}", h(a.getDeviceByAlias))
internal/appengine/http.go:48:	mux.Handle("GET "+base+"/devices/{device}/interfaces", h(a.listDeviceInterfaces))
internal/appengine/http.go:50:	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}", h(a.getData))
internal/appengine/http.go:51:	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.getData))
internal/appengine/http.go:52:	mux.Handle("PUT "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
internal/appengine/http.go:53:	mux.Handle("POST "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
internal/appengine/http.go:54:	mux.Handle("DELETE "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteData))
internal/appengine/http.go:59:	mux.Handle("PATCH "+base+"/devices-by-alias/{alias}", h(a.patchDeviceByAlias))
internal/appengine/http.go:60:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces", h(a.listInterfacesByAlias))
internal/appengine/http.go:61:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces/{interface}", h(a.getDataByAlias))
internal/appengine/http.go:62:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.getDataByAlias))
internal/appengine/http.go:63:	mux.Handle("PUT "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.putDataByAlias))
internal/appengine/http.go:64:	mux.Handle("POST "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.putDataByAlias))
internal/appengine/http.go:65:	mux.Handle("DELETE "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.deleteDataByAlias))
internal/appengine/http.go:67:	mux.Handle("GET "+base+"/groups", h(a.listGroups))
internal/appengine/http.go:68:	mux.Handle("POST "+base+"/groups", h(a.createGroup))
internal/appengine/http.go:69:	mux.Handle("GET "+base+"/groups/{group}", h(a.getGroup))
internal/appengine/http.go:70:	mux.Handle("GET "+base+"/groups/{group}/devices", h(a.listGroupDevices))
internal/appengine/http.go:71:	mux.Handle("POST "+base+"/groups/{group}/devices", h(a.addGroupDevice))
internal/appengine/http.go:72:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}", h(a.getGroupDevice))
internal/appengine/http.go:73:	mux.Handle("PATCH "+base+"/groups/{group}/devices/{device}", h(a.patchGroupDevice))
internal/appengine/http.go:74:	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}", h(a.removeGroupDevice))
internal/appengine/http.go:75:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces", h(a.listInterfacesInGroup))
internal/appengine/http.go:76:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}", h(a.getDataInGroup))
internal/appengine/http.go:77:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.getDataInGroup))
internal/appengine/http.go:78:	mux.Handle("PUT "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.putDataInGroup))
internal/appengine/http.go:79:	mux.Handle("POST "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.putDataInGroup))
internal/appengine/http.go:80:	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteDataInGroup))
internal/appengine/http.go:209:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:230:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:521:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/stream/ws.go:46:	mux.Handle("GET /astrate/v1/{realm}/socket", a.require(http.HandlerFunc(a.handle)))
internal/appengine/stream/ws.go:152:	return r.Header.Get("Accept") == "text/event-stream" || r.URL.Query().Get("transport") == "sse"
internal/appengine/channels/ws.go:48:	mux.Handle("GET /appengine/v1/socket/websocket", http.HandlerFunc(a.handle))

$ rg -n '^\s{2}/' docs/api/astarte_appengine_api.yaml
26:  /appengine/v1/{realm}/devices:
98:  /appengine/v1/{realm}/stats/devices:
136:  /appengine/v1/{realm}/devices/{device}:
210:  /appengine/v1/{realm}/devices-by-alias/{alias}:
292:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces:
324:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
369:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
555:  /appengine/v1/{realm}/devices/{device}/interfaces:
583:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
623:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
789:  /appengine/v1/{realm}/groups:
863:  /appengine/v1/{realm}/groups/{group}:
891:  /appengine/v1/{realm}/groups/{group}/devices:
981:  /appengine/v1/{realm}/groups/{group}/devices/{device}:
1080:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces:
1108:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}:
1149:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}:

$ ls .mule/reports/ 2>/dev/null | head -50; echo "---"; git log --oneline -15 -- .mule/ 2>/dev/null
hygiene.md
log
---
c6ab452 mule: log httpx-cors-vary-origin-passthrough
0af7cfe mule: refill the queue
1514f97 mule: recipe code-review ran (246s)
66f4932 mule: recipe astarte-upstream ran (106s)
5a314c7 mule: refill the queue
b3e8361 mule: recipe github-issues ran (126s)
9d75bee mule: recipe milestones timed out
d9a534b mule: recipe hygiene timed out
e33ab46 mule: log docs-sync-native-version-value
30665a5 mule: log docs-sync-pairing-version-value
4c12a38 mule: refill the queue
ee22e79 mule: recipe docs-sync ran (588s)
3f687a0 mule: log observability-health-content-type-test
2de9e46 mule: log observability-dbpool-gauge-coverage
25a3fde mule: log observability-readiness-wedge-test

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc|\.Handle\(' internal/pairing/ internal/housekeeping/ internal/realm/ cmd/ | head -60
internal/pairing/http.go:81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/pairing/http.go:82:	mux.Handle("POST /pairing/v1/{realm}/agent/devices",
internal/pairing/http.go:84:	mux.Handle("DELETE /pairing/v1/{realm}/agent/devices/{deviceID}",
internal/pairing/http.go:86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
internal/pairing/http.go:88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
internal/pairing/http.go:90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
internal/pairing/http.go:344:	header := r.Header.Get("Authorization")
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
cmd/astrate/main.go:433:	mux.Handle("GET /appengine/v1/{realm}/version",
cmd/astrate/main.go:435:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
cmd/astrate/main_test.go:153:		if got := resp.Header.Get("Access-Control-Allow-Origin"); got != "http://dashboard.example:4040" {
internal/housekeeping/http.go:37:	mux.Handle("GET /housekeeping/v1/realms", h(a.listRealms))
internal/housekeeping/http.go:38:	mux.Handle("POST /housekeeping/v1/realms", h(a.createRealm))
internal/housekeeping/http.go:39:	mux.Handle("GET /housekeeping/v1/realms/{realm}", h(a.getRealm))
internal/housekeeping/http.go:40:	mux.Handle("PATCH /housekeeping/v1/realms/{realm}", h(a.patchRealm))
internal/housekeeping/http.go:41:	mux.Handle("DELETE /housekeeping/v1/realms/{realm}", h(a.deleteRealm))

$ git log --oneline -8 -- docs/api/ 2>/dev/null; echo "---"; git log -1 --format='%ci' 2>/dev/null
08048b3 mule: docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
c69a164 mule: docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
bc73f65 mule: docs-sync-pairing-version-404-unreachable [auto]: remove the `404` "Unknown realm" response from `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml (yaml:96-104) — the route is served by `observability.VersionHandler` (cmd/astrate/main.go:398, internal/observability/compat.go:38-47), which answers static `200 {"data":"<version>"}` with no realm lookup, so a 404 can never be emitted there. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
f0d6a4d mule: docs-sync-pairing-deviceendpoints-dead-404-403 [auto]: remove the unreachable `404` DeviceNotFound response from the three device-facing ops in docs/api/astarte_pairing_api.yaml — `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials` (yaml:226-227), `GET /pairing/v1/{realm}/devices/{deviceID}` (yaml:268-269), `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify` (yaml:339-340) — and the `403` Forbidden from the info GET (yaml:266-267). `Service.authenticateDevice` folds every failure (malformed/unknown deviceID, unknown realm, unregistered device, wrong secret) into uniform `ErrUnauthorized` (internal/pairing/service.go:380-411), so `store.ErrNotFound` never reaches `writeServiceError` (http.go:326-327) from those handlers; and `Service.Info` never returns `ErrInhibited` (inhibited devices still read info with 200 + status "inhibited", service.go:286-307, comment at 282-285). Reachable on these three are only 200/201/400/401/422/429/500 plus 403 on credentials and verify (both genuinely raise `ErrInhibited`, service.go:252-254 and 329-331). Keep 404 on the agent endpoints (register yaml:149-150, unregister yaml:176-177) where `GetRealmByName`/`deviceid.Parse` really do return `store.ErrNotFound`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
8148aae mule: docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication is "via Authorization header or `?token=` query parameter" (yaml:293-294), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
d762f00 mule: docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_housekeeping_api.yaml (yaml:417) — it shows `detail: "housekeeping: validation failed: realm name must not be empty"`, but `validationDetail` strips the `housekeeping: validation failed: ` prefix on the wire (internal/housekeeping/http.go:251-258) and the real blank-name message is `"realm_name can't be blank"` (service.go:138). Change the example to `detail: "realm_name can't be blank"`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
d07b6c7 mule: docs-sync-hk-delete-gating-responses [auto]: add the missing `405` and `422` responses to `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — `writeError` answers 405 "Realm deletion disabled" for `ErrDeletionDisabled` and 422 `{"errors":{"error_name":["connected_devices_present"]}}` for `ErrConnectedDevicesPresent` (internal/housekeeping/http.go:233-249), both reachable from `Service.DeleteRealm` (service.go:248-266), but the spec lists only 204/401/403/404/500 (yaml:198-208). Its sole `405` reference sits on PATCH (yaml:182-183), which can never return it (`UpdateRealm`, service.go:196-216, has no deletion-gating path) — re-point 405 to DELETE (reusing the MethodNotAllowed component) and add 422 with the FieldErrors shape and a connected_devices_present example. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
64361d9 mule: docs-sync-ae-patch-merge-patch-media-type [auto]: fix the three PATCH request bodies in docs/api/astarte_appengine_api.yaml — PATCH /devices/{device} (yaml:176), PATCH /devices-by-alias/{alias} (yaml:260), PATCH /groups/{group}/devices/{device} (yaml:1023) declare `application/json; charset=utf-8`, but all three handlers enforce `Content-Type: application/merge-patch+json` exactly and reject every other value — 500 on the device and group paths (internal/appengine/http.go:209-211, 521-523), 400 on the alias path (http.go:230-232). Change the requestBody media type to `application/merge-patch+json` on all three. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
---
2026-09-17 21:22:17 +0200

$ rg -n '^\s{2}/' docs/api/astrate_native_api.yaml docs/api/astarte_pairing_api.yaml docs/api/astarte_realm_management_api.yaml docs/api/astarte_housekeeping_api.yaml
docs/api/astrate_native_api.yaml:26:  /astrate/v1/health:
docs/api/astrate_native_api.yaml:44:  /astrate/v1/readiness:
docs/api/astrate_native_api.yaml:77:  /astrate/v1/metrics:
docs/api/astrate_native_api.yaml:99:  /appengine/health:
docs/api/astrate_native_api.yaml:127:  /realmmanagement/health:
docs/api/astrate_native_api.yaml:155:  /pairing/health:
docs/api/astrate_native_api.yaml:185:  /appengine/version:
docs/api/astrate_native_api.yaml:205:  /realmmanagement/version:
docs/api/astrate_native_api.yaml:225:  /pairing/version:
docs/api/astrate_native_api.yaml:245:  /housekeeping/version:
docs/api/astrate_native_api.yaml:267:  /astrate/v1/{realm}/socket:
docs/api/astrate_native_api.yaml:328:  /appengine/v1/socket/websocket:
docs/api/astarte_pairing_api.yaml:27:  /pairing/v1/{realm}/health:
docs/api/astarte_pairing_api.yaml:70:  /pairing/v1/{realm}/version:
docs/api/astarte_pairing_api.yaml:97:  /pairing/v1/{realm}/agent/devices:
docs/api/astarte_pairing_api.yaml:149:  /pairing/v1/{realm}/agent/devices/{deviceID}:
docs/api/astarte_pairing_api.yaml:172:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
docs/api/astarte_pairing_api.yaml:224:  /pairing/v1/{realm}/devices/{deviceID}:
docs/api/astarte_pairing_api.yaml:258:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:
docs/api/astarte_housekeeping_api.yaml:26:  /housekeeping/v1/realms:
docs/api/astarte_housekeeping_api.yaml:106:  /housekeeping/v1/realms/{realm}:
docs/api/astarte_realm_management_api.yaml:27:  /realmmanagement/v1/{realm}/interfaces:
docs/api/astarte_realm_management_api.yaml:138:  /realmmanagement/v1/{realm}/interfaces/{name}:
docs/api/astarte_realm_management_api.yaml:167:  /realmmanagement/v1/{realm}/interfaces/{name}/{major}:
docs/api/astarte_realm_management_api.yaml:279:  /realmmanagement/v1/{realm}/triggers:
docs/api/astarte_realm_management_api.yaml:349:  /realmmanagement/v1/{realm}/triggers/{name}:
docs/api/astarte_realm_management_api.yaml:401:  /realmmanagement/v1/{realm}/config/auth:
docs/api/astarte_realm_management_api.yaml:470:  /realmmanagement/v1/{realm}/config/device_registration_limit:
docs/api/astarte_realm_management_api.yaml:501:  /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention:
docs/api/astarte_realm_management_api.yaml:536:  /realmmanagement/v1/{realm}/version:
docs/api/astarte_realm_management_api.yaml:569:  /realmmanagement/v1/{realm}/devices/{device}:
docs/api/astarte_realm_management_api.yaml:595:  /realmmanagement/v1/{realm}/policies:
docs/api/astarte_realm_management_api.yaml:665:  /realmmanagement/v1/{realm}/policies/{name}:

All routes are documented in each spec — the drift is now in the details. I picked **realm-management**: its `policies` and `triggers` ops were never touched by earlier docs-sync runs. Let me check those against the handlers.
→ Read docs/api/astarte_realm_management_api.yaml [offset=569, limit=130]
→ Read internal/realm/http.go
→ Read docs/api/astarte_realm_management_api.yaml [limit=75, offset=695]
→ Read docs/api/astarte_realm_management_api.yaml [limit=125, offset=279]
✱ Grep "func \(s \*Service\) (CreatePolicy|GetPolicy|DeletePolicy|ListPolicies|CreateTrigger|GetTrigger|DeleteTrigger|ListTriggers)" in internal/realm · 8 matches
→ Read internal/realm/service.go [limit=160, offset=419]
✱ Grep "func validatePolicy" in internal/realm · 1 match
Now verify the deletePolicy 422 path — the spec documents only 204/401/404/500 for it.
✱ Grep "func \(s \*Service\) realmID" in internal/realm · 1 match
→ Read internal/realm/policy.go
✱ Grep "func Write(Data|Error|NotFound|BadRequest|InternalServerError|FieldErrors)" in pkg/astarteapi · 8 matches
Both `policies` and `triggers` handlers route through the generic `writeError`. Let me check `WriteError`/`WriteNotFound`/`WriteData` shapes and the spec's response components.
→ Read docs/api/astarte_realm_management_api.yaml [limit=90, offset=771]
→ Read docs/api/astarte_realm_management_api.yaml [limit=261, offset=861]
→ Read internal/realm/service.go [limit=30, offset=139]
✱ Grep "func WriteError|func WriteNotFound|func WriteInternalServerError|WriteRawErrors" in pkg/astarteapi/envelope.go · 7 matches
→ Read docs/api/astarte_realm_management_api.yaml [limit=112, offset=27]
✱ Grep "func .*GetRealmByName" in internal/store · 1 match
The trigger/policy surface of realm-management has never been covered by earlier docs-sync runs. I've verified the handlers and spec. Let me run the recipe's cheap config-key drift check too, plus verify the docs build is intact.

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

$ rg -n 'HOUSEKEEPING' docs/site/configuration-reference.md
(no output)

Two code-side config keys are absent from the docs. Let me verify they're genuine, user-facing keys (not test-only) and check how the doc renames them.
✱ Grep "ASTRATE_HOUSEKEEPING|HousekeepingDefault|RealmDeletionDisabled" in internal/config/config.go · 13 matches

$ rg -n 'realm_deletion|datastream_maximum_storage_retention|housekeeping' docs/site/configuration-reference.md
36:| `addr` | string | `":8080"` | `ASTRATE_HTTP_ADDR` | No | Bind address for the single REST listener (pairing, realm management, housekeeping, appengine, native endpoints). |
78:## `[housekeeping]` — Instance-admin keys
139:[housekeeping]
140:jwt_public_key_files = ["/etc/astrate/housekeeping.pub"]

→ Read docs/site/configuration-reference.md [limit=75, offset=78]
Both config keys exist in `internal/config/config.go` (real fields, not test-only): `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (config.go:297-302) and `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (config.go:279). Neither appears anywhere in `configuration-reference.md`, whose `[housekeeping]` section documents only the two JWT key fields. Let me confirm the TOML field names and defaults to make a precise for-giulio entry.
→ Read internal/config/config.go [limit=40, offset=110]
→ Read internal/config/config.go [limit=40, offset=270]
Now let me pin down whether "unknown realm → 404" is genuinely reachable on these handlers, or the middleware blocks that path first.
✱ Grep "func .*RequireRealm|ClaimRealmManagement|GetRealmByName|ParseRealm|realm" in internal/auth/middleware.go · 55 matches
→ Read internal/auth/middleware.go [limit=50, offset=48]
Let me confirm `ErrValidation`'s string so the proposed detail matches what the wire would show.
✱ Grep "ErrValidation = |ErrValidation" in internal/realm · 21 matches
Let me check whether `CompilePolicy` can produce the nested `TriggerErrors` envelope (which would affect the createPolicy 422 shape), and read for-giulio.md.
✱ Grep "func CompilePolicy" in internal/engine/triggers · 1 match
→ Read .mule/for-giulio.md
Two solid findings so far. Let me check the docs/site RM page for any policy/trigger prose contradiction.

$ rg -n 'polic|trigger' docs/site/realm-management-api.md
3:The Realm Management API manages interface (schema) and trigger definitions per realm. It is the operator's interface to the data model.
60:### Install trigger
63:POST /realmmanagement/v1/<realm>/triggers
64:{ "data": <trigger JSON definition> }
67:Validates the trigger definition against the Astarte trigger schema (simple triggers + action).
69:### List triggers
72:GET /realmmanagement/v1/<realm>/triggers
75:### Get trigger
78:GET /realmmanagement/v1/<realm>/triggers/<name>
81:### Update trigger
84:PUT /realmmanagement/v1/<realm>/triggers/<name>
85:{ "data": <updated trigger JSON> }
88:### Delete trigger
91:DELETE /realmmanagement/v1/<realm>/triggers/<name>
107:Every mutation (interface install/update/delete, trigger CRUD) issues a Postgres `NOTIFY astrate_interfaces`. The engine listens and rebuilds the affected realm's compiled cache. In a single process this is a belt-and-suspenders measure; it exists so an optional hot-standby instance stays coherent.

The site doc says "Update trigger" with `PUT /triggers/{name}` — let me verify the code really has no such route.
$ rg -n 'UpdateTrigger|PUT .*triggers|/triggers/\{name\}' internal/ --glob '!**/*_test.go'
internal/realm/http.go:46:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers/{name}", h(a.getTrigger))
internal/realm/http.go:47:	mux.Handle("DELETE /realmmanagement/v1/{realm}/triggers/{name}", h(a.deleteTrigger))

The site page documents an "Update trigger" PUT route. Confirmed absent from the code. Let me verify the triggered `TriggerErrors` path produces a nested envelope on a real malformed trigger, to state the 422 task line accurately.
$ rg -n 'TriggerErrors' internal/engine/triggers/*_test.go | head; echo ---; rg -n 'func Compile\(' internal/engine/triggers/
internal/engine/triggers/match_test.go:449:	var te *TriggerErrors
internal/engine/triggers/match_test.go:451:		t.Fatalf("err %v is not *TriggerErrors", err)
---
internal/engine/triggers/match.go:259:func Compile(name string, def []byte) (*Trigger, error) {

→ Read internal/engine/triggers/match.go [limit=45, offset=259]
Let me verify an "Update trigger" PUT is genuinely absent (it would be the site-vs-code contradiction), and confirm the docs build target exists so the task lines' instruction is accurate.

$ rg -rn 'UpdateTrigger|"PUT .*triggers"' internal/ --glob '!**/*_test.go'; echo "exit: $?"; echo ---; rg -n 'build:' docs/Makefile
exit: 1
---
16:build: sync

→ Read docs/site/realm-management-api.md [limit=40, offset=60]
I have my findings. Let me verify the docs api 422 shapes and confirm the trigger 422 nested envelope is reachable (both shapes), which I've already confirmed via http.go:252-260. Now let me write the deliverables: two todo.md task lines and one for-giulio.md entry.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -72,8 +72,11 @@
     bash tools/mule.sh review
 
 ## Queue
 
+- [ ] docs-sync-rm-policies-delete-422 [auto]: add the missing `422` response to `DELETE /realmmanagement/v1/{realm}/policies/{name}` in docs/api/astarte_realm_management_api.yaml — `Service.DeletePolicy` refuses to delete a policy still referenced by a trigger (`"policy %q is still used by trigger %q"`, internal/realm/service.go:573-574, wrapped as ErrValidation) and `writeError` answers 422 for it (internal/realm/http.go:410-411, detail stripped by validationDetail), but the spec lists only 204/401/404/500 for the delete op (yaml:695-713). Add the 422 ValidationError response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-rm-triggers-422-nested-envelope [auto]: document the nested changeset envelope on `POST /realmmanagement/v1/{realm}/triggers` in docs/api/astarte_realm_management_api.yaml — when trigger compile fails, `createTrigger` answers 422 with the upstream nested body `{"errors":{"action":{...},"simple_triggers":[...]}}` (internal/realm/http.go:252-258, triggerErrorBody at http.go:268-285, issues #63/#70), but the spec's createTrigger 422 (yaml:344-345) refs only the flat `ValidationError` component (`{"errors":{"detail":...}}`) — which the handler also emits for plain ErrValidation (http.go:259-260, e.g. "trigger requires a name", service.go:429). Model the response as oneOf across the flat ErrorDetail and the nested envelope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+
 - [!] hk-zero-reglimit-create-asymmetry: POST /realms persists device_registration_limit 0 as a literal 0 (service.go:143 only rejects < 0) while upstream Astarte v1.2.0 folds 0 to nil at create (`insert_realm` does `if device_limit == 0, do: nil`, so GET returns null); PATCH 0 already matches upstream (sets 0, http.go:208-214). Fold 0 → nil in CreateRealm for parity, same pattern as the retention fix, and pin with a container-free test. — BLOCKED: opencode exited 1
 
 - [x] flowapi-autorestart-terminal-failure: `restartWithBackoff` (internal/flowapi/service.go:585-621) retries forever on a permanent rebuild failure — a flow whose pipeline was deleted (service.go:360-363 wraps it as ErrValidation) or whose config/definition no longer resolves keeps marking itself failed and rebuilding every ≤30s with no exit short of flow deletion. Classify terminal vs transient (ErrValidation → permanent; store/manager errors → transient) via a pure helper and stop the loop on terminal failures; add a container-free unit test on the helper. [auto]
 - [ ] broker-external-bus-intake [legion]: implement the second Intake implementation named in the TODO at internal/broker/intake.go:56-62 — an external-bus-backed (e.g. NATS JetStream) Intake for multi-instance deployment / restart survival, reproducing the embedded broker's per-device ordering, deferred-ack backpressure and QoS 0 drop semantics; an implementation of this interface needs containerised integration tests, so it runs on the Legion Go.

Now the for-giulio.md entry — the docs/site page contradicts the code (documents a PUT /triggers route + "Update trigger" that doesn't exist).
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **docs-sync realm-management, 2026-09-17: `docs/site/realm-management-api.md` documents an "Update trigger" endpoint Astrate does not serve.** Lines 81-86 show `PUT /realmmanagement/v1/<realm>/triggers/<name>` with a `{ "data": <updated trigger JSON> }` body, but no such route is registered — internal/realm/http.go:44-47 wires only GET/POST `/triggers` and GET/DELETE `/triggers/{name}` (an update op is deliberately absent; not even a stub or TODO exists). Site prose is yours — drop the "Update trigger" section or confirm the route is planned.
+
+---
+
 - **github-issues triage run, 2026-09-17: still nothing proposable — the mule-alarm pile is now 12 straight days (#94–#105).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-16 run: **#105**, today's alarm (created 09-17 10:51Z, "nothing has landed in 14h"); **#104** has expired into the pile (superseded by #105). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-16 one: close #94–#104** — each was a one-day event, superseded by next day's alarm, never actionable; leave #105 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
 
 - **github-issues triage run, 2026-09-16: still nothing proposable — the mule-alarm pile is now 11 straight days (#94–#104).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), plus the daily alarm. New since the 2026-09-15 run: **#104**, today's alarm (created 09-16 11:12Z, "nothing has landed in 14h"); **#103** has expired into the pile (superseded by #104). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-15 one: close #94–#103** — each was a one-day event, superseded by next day's alarm, never actionable; leave #104 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
 

Let me confirm the docs build baseline so the proposed task lines reference a working build.
→ Read docs/Makefile
This is a proposal job with no docs/api edits, but the tasks I propose must name a working build gate. Let me verify the docs build baseline is green.

$ make -C docs build 2>&1 | tail -15
make: Entering directory '/root/astrate-mule/docs'
cp DESIGN.md site/DESIGN.md
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'

The docs build tool (`mkdocs`) isn't installed on this box, so I can't execute the build gate myself — I'll note that in the report; the fixing task lines already carry the `make -C docs build` instruction for wherever they run. Let me confirm my edits landed cleanly.
$ sed -n '74,80p' .mule/todo.md; echo ---; sed -n '9,15p' .mule/for-giulio.md
## Queue

- [ ] docs-sync-rm-policies-delete-422 [auto]: add the missing `422` response to `DELETE /realmmanagement/v1/{realm}/policies/{name}` in docs/api/astarte_realm_management_api.yaml — `Service.DeletePolicy` refuses to delete a policy still referenced by a trigger (`"policy %q is still used by trigger %q"`, internal/realm/service.go:573-574, wrapped as ErrValidation) and `writeError` answers 422 for it (internal/realm/http.go:410-411, detail stripped by validationDetail), but the spec lists only 204/401/404/500 for the delete op (yaml:695-713). Add the 422 ValidationError response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-rm-triggers-422-nested-envelope [auto]: document the nested changeset envelope on `POST /realmmanagement/v1/{realm}/triggers` in docs/api/astarte_realm_management_api.yaml — when trigger compile fails, `createTrigger` answers 422 with the upstream nested body `{"errors":{"action":{...},"simple_triggers":[...]}}` (internal/realm/http.go:252-258, triggerErrorBody at http.go:268-285, issues #63/#70), but the spec's createTrigger 422 (yaml:344-345) refs only the flat `ValidationError` component (`{"errors":{"detail":...}}`) — which the handler also emits for plain ErrValidation (http.go:259-260, e.g. "trigger requires a name", service.go:429). Model the response as oneOf across the flat ErrorDetail and the nested envelope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

- [!] hk-zero-reglimit-create-asymmetry: POST /realms persists device_registration_limit 0 as a literal 0 (service.go:143 only rejects < 0) while upstream Astarte v1.2.0 folds 0 to nil at create (`insert_realm` does `if device_limit == 0, do: nil`, so GET returns null); PATCH 0 already matches upstream (sets 0, http.go:208-214). Fold 0 → nil in CreateRealm for parity, same pattern as the retention fix, and pin with a container-free test. — BLOCKED: opencode exited 1

---
line once you have dealt with it — this file is a queue, not a log.

---

- **docs-sync realm-management, 2026-09-17: `docs/site/realm-management-api.md` documents an "Update trigger" endpoint Astrate does not serve.** Lines 81-86 show `PUT /realmmanagement/v1/<realm>/triggers/<name>` with a `{ "data": <updated trigger JSON> }` body, but no such route is registered — internal/realm/http.go:44-47 wires only GET/POST `/triggers` and GET/DELETE `/triggers/{name}` (an update op is deliberately absent; not even a stub or TODO exists). Site prose is yours — drop the "Update trigger" section or confirm the route is planned.

---

Proposal run complete. I checked the realm-management surface (triggers + policies ops, never covered by earlier docs-sync runs), verified every route in the spec exists in code and vice versa, then checked the trigger/policy handlers' status codes and response shapes against the spec.

**Findings, all on `docs/api` + `docs/site`:**

1. `DELETE /policies/{name}` — spec lists only 204/401/404/500, but `DeletePolicy` answers 422 when a trigger still references the policy (service.go:573-574 → http.go:410-411). **Missing 422** → task line.
2. `POST /triggers` 422 — handler emits the upstream nested envelope `{"errors":{"action":...,"simple_triggers":[...]}}` for compile failures (http.go:252-258), but the spec's 422 refs only the flat `ValidationError` component. **Response shape under-documented** → task line.
3. `docs/site/realm-management-api.md:81-86` documents `PUT /triggers/{name}` "Update trigger" — **no such route exists** in code (http.go:44-47). Site prose → escalated to for-giulio.md.

The `ASTRATE_HOUSEKEEPING_*` config-key gap the recipe's sweep flags is already escalated twice in for-giulio.md (2026-09-12 and earlier), so I did not duplicate it.

```
Done: Realm-management docs-sync: policies delete missing 422, triggers 422 nested envelope, site-update-trigger contradiction
Files: .mule/todo.md (2 task lines), .mule/for-giulio.md (1 entry)
Verified: rg route/status-code/shape cross-checks against internal/realm/http.go + service.go + docs/api yaml -> pass; make -C docs build -> fail (mkdocs not installed on this box, so could not run the Swagger-YAML-load gate here; task lines carry that instruction)
Unsure: nothing — every claimed status code and response body is read directly from the handler/service code
Follow-ups: 2 lines appended to .mule/todo.md: docs-sync-rm-policies-delete-422, docs-sync-rm-triggers-422-nested-envelope; 1 entry in .mule/for-giulio.md
```
