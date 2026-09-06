slug: recipe-docs-sync
verdict: proposed
at:  296c092
ran: 2026-09-06T20:04:31Z on DietPi in 262s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls docs/api/ 2>/dev/null
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*appengine*' 2>/dev/null | head -80
(no output)

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(' internal/ --glob '*.go' -l 2>/dev/null | head -20
internal/flow/blocks/httpblocks.go
internal/flow/blocks/httpblocks_test.go
internal/observability/compat_version_test.go
internal/pairing/http.go
internal/engine/triggers/actions_test.go
internal/engine/forward/http_test.go
internal/httpx/cors.go
internal/swagger/swagger_test.go
internal/appengine/stream/ws.go
internal/appengine/http.go
internal/auth/middleware.go

$ rg -n '^\s{2}/' docs/api/astarte_appengine_api.yaml 2>/dev/null | head -40
26:  /appengine/v1/{realm}/devices:
93:  /appengine/v1/{realm}/stats/devices:
131:  /appengine/v1/{realm}/devices/{device}:
201:  /appengine/v1/{realm}/devices-by-alias/{alias}:
233:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
268:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
419:  /appengine/v1/{realm}/groups:
489:  /appengine/v1/{realm}/groups/{group}/devices:
564:  /appengine/v1/{realm}/groups/{group}/devices/{device}:

→ Read internal/appengine/http.go
→ Read docs/api/astarte_appengine_api.yaml
Now let me check a few specific status codes more carefully.
→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ rg -n 'listDeviceInterfaces' internal/appengine/ --include='*.go' | head -5
rg: unrecognized flag --include

similar flags that are available: --include-zero

$ rg -n 'listGroupDevices' internal/appengine/ --include='*.go' | head -5
rg: unrecognized flag --include

similar flags that are available: --include-zero

$ rg -n 'StatusConflict' internal/appengine/http.go
651:		_ = astarteapi.WriteError(w, http.StatusConflict, "Alias already in use")
661:		_ = astarteapi.WriteError(w, http.StatusConflict, "Group already exists")
663:		_ = astarteapi.WriteError(w, http.StatusConflict, "Device already in group")
665:		_ = astarteapi.WriteError(w, http.StatusConflict, "Already exists")

Now I have the full picture. Let me compare code vs docs for appengine and prepare the proposals.
$ rg -n 'from_token\|"limit"' /root/astrate-mule/internal/appengine/http.go | head -10
(no output)

$ rg -n 'from_token\|"limit"' /root/astrate-mule/internal/appengine/http.go
(no output)

$ rg -n 'from_token|limit' /root/astrate-mule/internal/appengine/http.go
21:// maxValueBytes is upstream's per-value limit (~64 KiB): a larger REST
22:// write is answered 422 "Value size exceeds size limits" (measured
87:	limit, _ := strconv.Atoi(q.Get("limit"))
89:	page, err := a.svc.ListDevices(r.Context(), r.PathValue("realm"), q.Get("from_token"), limit, details)
103:// always carries details, from_token, and limit: url.Values.Encode sorts
104:// keys, so "details" always precedes "from_token" — the dashboard parses
106:// corrupts the FIRST key/value pair, so from_token must never be first.
120:		limit := q.Get("limit")
121:		if limit == "" {
122:			limit = strconv.Itoa(DefaultDeviceLimit)
124:		nq := url.Values{"details": {details}, "from_token": {next}, "limit": {limit}}
299:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
328:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
347:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
439:	limit, _ := strconv.Atoi(q.Get("limit"))
441:		details, q.Get("from_token"), limit)
456:// from_token plus the effective limit (and details when the request had it),
458:// details precedes from_token precedes limit.
469:	limit := q.Get("limit")
470:	if limit == "" {
471:		limit = strconv.Itoa(DefaultDeviceLimit)
473:	nq := url.Values{"from_token": {next}, "limit": {limit}}
576:	if v := q.Get("limit"); v != "" {
580:			fe.addf("limit", "is invalid")
582:			fe.addf("limit", "must be greater than or equal to 0")

# Todos
[✓] Compare appengine code routes vs documented YAML paths
[✓] Check for undocumented query parameters and status code mismatches
[•] Append proposed task lines to .mule/todo.md

→ Read .mule/todo.md [offset=105]
Now I have the full comparison. Here are the findings:

**Surface: appengine.** 15 routes in code vs 9 documented paths. The missing ones are all the by-alias mirror endpoints, the group data endpoints, `GET /devices/{device}/interfaces`, and missing query params + status codes on existing paths.

Let me append the proposals to todo.md:

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -107,5 +107,8 @@
 - [x] docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` route to docs/api/astarte_realm_management_api.yaml — it exists in code (internal/realm/http.go:51, handler getDatastreamMaximumStorageRetention at http.go:135, returns 200 with a data-envelope retention value via Service.GetDatastreamMaximumStorageRetention) but the spec jumps straight from config/device_registration_limit (yaml:428) to /version. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] housekeeping-tests [auto]: add a container-free unit test file for `internal/housekeeping` (service.go, 262 lines, zero test coverage). Cover two rules: (1) `CreateRealm` validation — blank name, blank JWT key, negative regLimit, negative retention, and default retention injection when `WithDefaultDatastreamMaximumStorageRetention` is set and the caller passes nil; (2) `DeleteRealm` gating — `ErrDeletionDisabled` when `WithRealmDeletionDisabled(true)`, `ErrConnectedDevicesPresent` when `DeviceStats` reports connected > 0. Use httptest-style fakes for the store and sealer (no Docker).
 - [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
+- [ ] docs-sync-appengine-by-alias-endpoints [auto]: add the 8 undocumented by-alias mirror routes to docs/api/astarte_appengine_api.yaml — PATCH /devices-by-alias/{alias}, GET /devices-by-alias/{alias}/interfaces, GET/PUT/POST/DELETE /devices-by-alias/{alias}/interfaces/{interface}/{path} (plus the interface-scoped GET) — all registered in internal/appengine/http.go:59-65 and backed by handlers patchDeviceByAlias, listInterfacesByAlias, getDataByAlias, putDataByAlias, deleteDataByAlias (http.go:226-367). Match the existing spec style (same schemas, same response codes as the device-scoped counterparts). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-appengine-group-endpoints [auto]: add the 9 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (returns {data: {group_name: "<name>"}}, http.go:428-434), PATCH /groups/{group}/devices/{device} (DevicePatch body → updated status, http.go:518-537), and the 7 group data CRUD endpoints: GET/PUT/POST/DELETE /groups/{group}/devices/{device}/interfaces/{interface}/{path} plus the interface-level and path-level GETs (http.go:75-80, handlers listInterfacesInGroup, getDataInGroup, putDataInGroup, deleteDataInGroup at http.go:187-377). Also add GET /devices/{device}/interfaces (listDeviceInterfaces, http.go:48) which is missing from the spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups/{group}/devices is missing `from_token` and `limit` pagination params that code reads at http.go:439-441 (and returns via links.next), (2) POST /groups is missing 422 validation-error response (FieldErrors for blank group_name/empty devices, http.go:404-416), (3) PATCH /devices/{device} and PATCH /devices-by-alias/{alias} are missing 422 (parseQueryOpts FieldErrors), (4) GET /devices/{device}/interfaces/{interface}/{path} is missing 422 (invalid since/since_after/to/limit/downsample_to/format params, http.go:554-639). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.

Verifying the appended lines:
→ Read .mule/todo.md [offset=111, limit=10]
The count wording in line 112 is inconsistent. Let me fix it:

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -108,7 +108,7 @@
 - [x] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] housekeeping-tests [auto]: add a container-free unit test file for `internal/housekeeping` (service.go, 262 lines, zero test coverage). Cover two rules: (1) `CreateRealm` validation — blank name, blank JWT key, negative regLimit, negative retention, and default retention injection when `WithDefaultDatastreamMaximumStorageRetention` is set and the caller passes nil; (2) `DeleteRealm` gating — `ErrDeletionDisabled` when `WithRealmDeletionDisabled(true)`, `ErrConnectedDevicesPresent` when `DeviceStats` reports connected > 0. Use httptest-style fakes for the store and sealer (no Docker).
 - [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
 - [ ] docs-sync-appengine-by-alias-endpoints [auto]: add the 8 undocumented by-alias mirror routes to docs/api/astarte_appengine_api.yaml — PATCH /devices-by-alias/{alias}, GET /devices-by-alias/{alias}/interfaces, GET/PUT/POST/DELETE /devices-by-alias/{alias}/interfaces/{interface}/{path} (plus the interface-scoped GET) — all registered in internal/appengine/http.go:59-65 and backed by handlers patchDeviceByAlias, listInterfacesByAlias, getDataByAlias, putDataByAlias, deleteDataByAlias (http.go:226-367). Match the existing spec style (same schemas, same response codes as the device-scoped counterparts). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
-- [ ] docs-sync-appengine-group-endpoints [auto]: add the 9 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (returns {data: {group_name: "<name>"}}, http.go:428-434), PATCH /groups/{group}/devices/{device} (DevicePatch body → updated status, http.go:518-537), and the 7 group data CRUD endpoints: GET/PUT/POST/DELETE /groups/{group}/devices/{device}/interfaces/{interface}/{path} plus the interface-level and path-level GETs (http.go:75-80, handlers listInterfacesInGroup, getDataInGroup, putDataInGroup, deleteDataInGroup at http.go:187-377). Also add GET /devices/{device}/interfaces (listDeviceInterfaces, http.go:48) which is missing from the spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-appengine-group-endpoints [auto]: add the 8 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (returns {data: {group_name: "<name>"}}, http.go:428-434), PATCH /groups/{group}/devices/{device} (DevicePatch body → updated status, http.go:518-537), and the 6 group interface-data routes: the interface-level GET and the path-level GET/PUT/POST/DELETE on /groups/{group}/devices/{device}/interfaces (http.go:75-80, handlers listInterfacesInGroup/listInterfacesInGroup, getDataInGroup, putDataInGroup, deleteDataInGroup at http.go:187-377). Also add GET /devices/{device}/interfaces (listDeviceInterfaces, http.go:48) which is missing from the spec — that is the 9th route. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups/{group}/devices is missing `from_token` and `limit` pagination params that code reads at http.go:439-441 (and returns via links.next), (2) POST /groups is missing 422 validation-error response (FieldErrors for blank group_name/empty devices, http.go:404-416), (3) PATCH /devices/{device} and PATCH /devices-by-alias/{alias} are missing 422 (parseQueryOpts FieldErrors), (4) GET /devices/{device}/interfaces/{interface}/{path} is missing 422 (invalid since/since_after/to/limit/downsample_to/format params, http.go:554-639). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -108,7 +108,7 @@
 - [x] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] housekeeping-tests [auto]: add a container-free unit test file for `internal/housekeeping` (service.go, 262 lines, zero test coverage). Cover two rules: (1) `CreateRealm` validation — blank name, blank JWT key, negative regLimit, negative retention, and default retention injection when `WithDefaultDatastreamMaximumStorageRetention` is set and the caller passes nil; (2) `DeleteRealm` gating — `ErrDeletionDisabled` when `WithRealmDeletionDisabled(true)`, `ErrConnectedDevicesPresent` when `DeviceStats` reports connected > 0. Use httptest-style fakes for the store and sealer (no Docker).
 - [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
 - [ ] docs-sync-appengine-by-alias-endpoints [auto]: add the 8 undocumented by-alias mirror routes to docs/api/astarte_appengine_api.yaml — PATCH /devices-by-alias/{alias}, GET /devices-by-alias/{alias}/interfaces, GET/PUT/POST/DELETE /devices-by-alias/{alias}/interfaces/{interface}/{path} (plus the interface-scoped GET) — all registered in internal/appengine/http.go:59-65 and backed by handlers patchDeviceByAlias, listInterfacesByAlias, getDataByAlias, putDataByAlias, deleteDataByAlias (http.go:226-367). Match the existing spec style (same schemas, same response codes as the device-scoped counterparts). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
-- [ ] docs-sync-appengine-group-endpoints [auto]: add the 8 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (returns {data: {group_name: "<name>"}}, http.go:428-434), PATCH /groups/{group}/devices/{device} (DevicePatch body → updated status, http.go:518-537), and the 6 group interface-data routes: the interface-level GET and the path-level GET/PUT/POST/DELETE on /groups/{group}/devices/{device}/interfaces (http.go:75-80, handlers listInterfacesInGroup/listInterfacesInGroup, getDataInGroup, putDataInGroup, deleteDataInGroup at http.go:187-377). Also add GET /devices/{device}/interfaces (listDeviceInterfaces, http.go:48) which is missing from the spec — that is the 9th route. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-appengine-group-endpoints [auto]: add the 8 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (returns {data: {group_name: "<name>"}}, http.go:428-434), PATCH /groups/{group}/devices/{device} (DevicePatch body → updated status, http.go:518-537), and the 6 group interface-data routes: the interface-level GET and the path-level GET/PUT/POST/DELETE on /groups/{group}/devices/{device}/interfaces (http.go:75-80, handlers listInterfacesInGroup, getDataInGroup, putDataInGroup, deleteDataInGroup at http.go:187-377). Also add GET /devices/{device}/interfaces (listDeviceInterfaces, http.go:48) which is missing from the spec — that is the 9th route. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups/{group}/devices is missing `from_token` and `limit` pagination params that code reads at http.go:439-441 (and returns via links.next), (2) POST /groups is missing 422 validation-error response (FieldErrors for blank group_name/empty devices, http.go:404-416), (3) PATCH /devices/{device} and PATCH /devices-by-alias/{alias} are missing 422 (parseQueryOpts FieldErrors), (4) GET /devices/{device}/interfaces/{interface}/{path} is missing 422 (invalid since/since_after/to/limit/downsample_to/format params, http.go:554-639). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -107,8 +107,8 @@
 - [x] docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` route to docs/api/astarte_realm_management_api.yaml — it exists in code (internal/realm/http.go:51, handler getDatastreamMaximumStorageRetention at http.go:135, returns 200 with a data-envelope retention value via Service.GetDatastreamMaximumStorageRetention) but the spec jumps straight from config/device_registration_limit (yaml:428) to /version. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] housekeeping-tests [auto]: add a container-free unit test file for `internal/housekeeping` (service.go, 262 lines, zero test coverage). Cover two rules: (1) `CreateRealm` validation — blank name, blank JWT key, negative regLimit, negative retention, and default retention injection when `WithDefaultDatastreamMaximumStorageRetention` is set and the caller passes nil; (2) `DeleteRealm` gating — `ErrDeletionDisabled` when `WithRealmDeletionDisabled(true)`, `ErrConnectedDevicesPresent` when `DeviceStats` reports connected > 0. Use httptest-style fakes for the store and sealer (no Docker).
 - [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
-- [ ] docs-sync-appengine-by-alias-endpoints [auto]: add the 8 undocumented by-alias mirror routes to docs/api/astarte_appengine_api.yaml — PATCH /devices-by-alias/{alias}, GET /devices-by-alias/{alias}/interfaces, GET/PUT/POST/DELETE /devices-by-alias/{alias}/interfaces/{interface}/{path} (plus the interface-scoped GET) — all registered in internal/appengine/http.go:59-65 and backed by handlers patchDeviceByAlias, listInterfacesByAlias, getDataByAlias, putDataByAlias, deleteDataByAlias (http.go:226-367). Match the existing spec style (same schemas, same response codes as the device-scoped counterparts). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-appengine-by-alias-endpoints [auto]: add the 7 undocumented by-alias mirror routes to docs/api/astarte_appengine_api.yaml — PATCH /devices-by-alias/{alias}, GET /devices-by-alias/{alias}/interfaces, GET/PUT/POST/DELETE /devices-by-alias/{alias}/interfaces/{interface}/{path} (plus the interface-scoped GET) — all registered in internal/appengine/http.go:59-65 and backed by handlers patchDeviceByAlias, listInterfacesByAlias, getDataByAlias, putDataByAlias, deleteDataByAlias (http.go:226-367). Match the existing spec style (same schemas, same response codes as the device-scoped counterparts). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] docs-sync-appengine-group-endpoints [auto]: add the 8 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (returns {data: {group_name: "<name>"}}, http.go:428-434), PATCH /groups/{group}/devices/{device} (DevicePatch body → updated status, http.go:518-537), and the 6 group interface-data routes: the interface-level GET and the path-level GET/PUT/POST/DELETE on /groups/{group}/devices/{device}/interfaces (http.go:75-80, handlers listInterfacesInGroup, getDataInGroup, putDataInGroup, deleteDataInGroup at http.go:187-377). Also add GET /devices/{device}/interfaces (listDeviceInterfaces, http.go:48) which is missing from the spec — that is the 9th route. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups/{group}/devices is missing `from_token` and `limit` pagination params that code reads at http.go:439-441 (and returns via links.next), (2) POST /groups is missing 422 validation-error response (FieldErrors for blank group_name/empty devices, http.go:404-416), (3) PATCH /devices/{device} and PATCH /devices-by-alias/{alias} are missing 422 (parseQueryOpts FieldErrors), (4) GET /devices/{device}/interfaces/{interface}/{path} is missing 422 (invalid since/since_after/to/limit/downsample_to/format params, http.go:554-639). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.

→ Read .mule/todo.md [limit=4, offset=111]
Let me verify the 422 claim on the PATCH endpoints before finalizing:

$ rg -n 'PatchDevice|FieldErrors' /root/astrate-mule/internal/appengine/service.go | head -20
53:// FieldErrors carries per-field validation failures rendered as the
55:type FieldErrors map[string][]string
57:func (fe FieldErrors) Error() string { return "appengine: invalid request payload" }
60:func (fe FieldErrors) addf(field, format string, args ...any) {
290:// PatchDevice applies a device patch (upstream PATCH /devices/{id}).
291:func (s *Service) PatchDevice(ctx context.Context, realm, deviceID string, p DevicePatch) (*DeviceStatus, error) {
307:// PatchDeviceByAlias applies a device patch addressed by alias (upstream
309:func (s *Service) PatchDeviceByAlias(ctx context.Context, realm, alias string, p DevicePatch) (*DeviceStatus, error) {
360:		if err := s.st.PatchDeviceAliases(ctx, rid, id, p.Aliases); err != nil {
365:		if err := s.st.PatchDeviceAttributes(ctx, rid, id, p.Attributes); err != nil {
467:	fe := FieldErrors{}
550:			return nil, FieldErrors{"from_token": {"is invalid"}}
554:		return nil, FieldErrors{"limit": {"must be greater than or equal to 0"}}
633:		return FieldErrors{"device_id": {"is not a valid device id"}}

→ Read internal/appengine/service.go [offset=290, limit=60]
→ Read internal/appengine/service.go [limit=30, offset=350]
Confirmed the PATCH 422 comes from `ErrAttributeKeyNotFound`, not `parseQueryOpts`. Fixing that attribution in line 113 and folding in the also-missing 409:

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -109,6 +109,6 @@
 - [x] housekeeping-tests [auto]: add a container-free unit test file for `internal/housekeeping` (service.go, 262 lines, zero test coverage). Cover two rules: (1) `CreateRealm` validation — blank name, blank JWT key, negative regLimit, negative retention, and default retention injection when `WithDefaultDatastreamMaximumStorageRetention` is set and the caller passes nil; (2) `DeleteRealm` gating — `ErrDeletionDisabled` when `WithRealmDeletionDisabled(true)`, `ErrConnectedDevicesPresent` when `DeviceStats` reports connected > 0. Use httptest-style fakes for the store and sealer (no Docker).
 - [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
 - [ ] docs-sync-appengine-by-alias-endpoints [auto]: add the 7 undocumented by-alias mirror routes to docs/api/astarte_appengine_api.yaml — PATCH /devices-by-alias/{alias}, GET /devices-by-alias/{alias}/interfaces, GET/PUT/POST/DELETE /devices-by-alias/{alias}/interfaces/{interface}/{path} (plus the interface-scoped GET) — all registered in internal/appengine/http.go:59-65 and backed by handlers patchDeviceByAlias, listInterfacesByAlias, getDataByAlias, putDataByAlias, deleteDataByAlias (http.go:226-367). Match the existing spec style (same schemas, same response codes as the device-scoped counterparts). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] docs-sync-appengine-group-endpoints [auto]: add the 8 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (returns {data: {group_name: "<name>"}}, http.go:428-434), PATCH /groups/{group}/devices/{device} (DevicePatch body → updated status, http.go:518-537), and the 6 group interface-data routes: the interface-level GET and the path-level GET/PUT/POST/DELETE on /groups/{group}/devices/{device}/interfaces (http.go:75-80, handlers listInterfacesInGroup, getDataInGroup, putDataInGroup, deleteDataInGroup at http.go:187-377). Also add GET /devices/{device}/interfaces (listDeviceInterfaces, http.go:48) which is missing from the spec — that is the 9th route. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
-- [ ] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups/{group}/devices is missing `from_token` and `limit` pagination params that code reads at http.go:439-441 (and returns via links.next), (2) POST /groups is missing 422 validation-error response (FieldErrors for blank group_name/empty devices, http.go:404-416), (3) PATCH /devices/{device} and PATCH /devices-by-alias/{alias} are missing 422 (parseQueryOpts FieldErrors), (4) GET /devices/{device}/interfaces/{interface}/{path} is missing 422 (invalid since/since_after/to/limit/downsample_to/format params, http.go:554-639). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups/{group}/devices is missing `from_token` and `limit` pagination params that code reads at http.go:439-441 (and returns via links.next), (2) POST /groups is missing 422 validation-error response (FieldErrors for blank group_name/empty devices, http.go:404-416), (3) PATCH /devices/{device} and PATCH /devices-by-alias/{alias} are missing 422 ("Attribute key not found" via ErrAttributeKeyNotFound, service.go:352-357 → http.go:656-657), and the device-scoped PATCH is also missing 409 ("Alias already in use", http.go:651) and its documented 500 covers the merge-patch content-type mismatch which code answers with an unmapped 500 too, (4) GET /devices/{device}/interfaces/{interface}/{path} is missing 422 (invalid since/since_after/to/limit/downsample_to/format params, http.go:554-639). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.

→ Read .mule/todo.md [limit=3, offset=111]
Report per the mule convention:

Done: docs-sync proposal run on the appengine surface — compared all handler routes vs the spec, appended 3 proposals for the drift found
Files: .mule/todo.md
Verified: rg over internal/appengine/http.go vs docs/api/astarte_appengine_api.yaml; handler/status-code claims cross-checked against http.go:59-80, 428-537 and service.go:290-379 -> pass
Unsure: nothing
Follow-ups: 3 lines appended to .mule/todo.md — docs-sync-appengine-by-alias-endpoints (7 missing by-alias routes, http.go:59-65), docs-sync-appengine-group-endpoints (8 group routes + GET /devices/{device}/interfaces, http.go:75-80), docs-sync-appengine-query-params-status (missing `from_token`/`limit`, missing 422/409 on POST /groups and PATCH devices, missing 422 on GET data path)

Findings: the spec (9 paths) documents less than half the code's 16 paths. All drift is in `docs/api/astarte_appengine_api.yaml`, a generated artefact — no docs/site prose contradiction found, so no `for-giulio.md` entry. No `gh issue create` required (recipe deliverable is task lines). Git untouched.
