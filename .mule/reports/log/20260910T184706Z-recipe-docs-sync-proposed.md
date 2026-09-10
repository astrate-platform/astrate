slug: recipe-docs-sync
verdict: proposed
at:  407f0c8
ran: 2026-09-10T18:47:06Z on DietPi in 289s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
✱ Glob "docs/api/*.yaml" in . · 5 matches

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' /root/astrate-mule/internal/ --glob '*pairing*' 2>/dev/null | head -60
(no output)

$ ls /root/astrate-mule/internal/
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

$ rg -ln 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' /root/astrate-mule/internal/ | head -40
/root/astrate-mule/internal/appengine/stream/ws.go
/root/astrate-mule/internal/appengine/http.go
/root/astrate-mule/internal/flow/blocks/httpblocks.go
/root/astrate-mule/internal/observability/compat.go
/root/astrate-mule/internal/observability/health.go
/root/astrate-mule/internal/observability/compat_version_test.go
/root/astrate-mule/internal/auth/middleware.go
/root/astrate-mule/internal/flow/blocks/httpblocks_test.go
/root/astrate-mule/internal/swagger/swagger.go
/root/astrate-mule/internal/swagger/swagger_test.go
/root/astrate-mule/internal/pairing/http.go
/root/astrate-mule/internal/engine/triggers/actions_test.go
/root/astrate-mule/internal/engine/forward/http_test.go
/root/astrate-mule/internal/httpx/notfound_test.go
/root/astrate-mule/internal/httpx/cors.go

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' /root/astrate-mule/internal/pairing/http.go
81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
344:	header := r.Header.Get("Authorization")

$ rg -n '^\s{2}/' /root/astrate-mule/docs/api/astarte_pairing_api.yaml
27:  /pairing/v1/{realm}/health:
70:  /pairing/v1/{realm}/agent/devices:
122:  /pairing/v1/{realm}/agent/devices/{deviceID}:
145:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
199:  /pairing/v1/{realm}/devices/{deviceID}:
237:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:

$ rg -n 'agent' /root/astrate-mule/internal/pairing/ | head -20
/root/astrate-mule/internal/pairing/http_test.go:112:// agentJWT signs a realm JWT with the given a_pa authorization strings (nil
/root/astrate-mule/internal/pairing/http_test.go:114:func (f *httpFixture) agentJWT(t *testing.T, claim auth.Claim, grants []string) string {
/root/astrate-mule/internal/pairing/http_test.go:205:	agentToken := f.agentJWT(t, auth.ClaimPairing, nil)
/root/astrate-mule/internal/pairing/http_test.go:217:		status, body := f.request(t, "POST", base+"/agent/devices", agentToken,
/root/astrate-mule/internal/pairing/http_test.go:230:		status, body := f.request(t, "POST", base+"/agent/devices", agentToken,
/root/astrate-mule/internal/pairing/http_test.go:248:			status, body := f.request(t, "POST", base+"/agent/devices", agentToken,
/root/astrate-mule/internal/pairing/http_test.go:256:		status, body := f.request(t, "POST", base+"/agent/devices", agentToken,
/root/astrate-mule/internal/pairing/http_test.go:265:		status, body := f.request(t, "POST", base+"/agent/devices", "",
/root/astrate-mule/internal/pairing/http_test.go:272:		wrongClaim := f.agentJWT(t, auth.ClaimAppEngine, nil)
/root/astrate-mule/internal/pairing/http_test.go:273:		status, body = f.request(t, "POST", base+"/agent/devices", wrongClaim,
/root/astrate-mule/internal/pairing/http_test.go:319:		status, body := f.request(t, "POST", base+"/agent/devices", agentToken,
/root/astrate-mule/internal/pairing/http_test.go:503:		status, body := f.request(t, "DELETE", base+"/agent/devices/"+hwID, agentToken, nil)
/root/astrate-mule/internal/pairing/http_test.go:520:		status, body = f.request(t, "POST", base+"/agent/devices", agentToken,
/root/astrate-mule/internal/pairing/http_test.go:527:		status, body = f.request(t, "DELETE", base+"/agent/devices/"+mustRandomDeviceID(t), agentToken, nil)
/root/astrate-mule/internal/pairing/http_test.go:546:		status, body := f.request(t, "POST", "/pairing/v1/"+limited.Name+"/agent/devices",
/root/astrate-mule/internal/pairing/http_test.go:547:			agentToken, map[string]string{"hw_id": mustRandomDeviceID(t)})
/root/astrate-mule/internal/pairing/http.go:82:	mux.Handle("POST /pairing/v1/{realm}/agent/devices",
/root/astrate-mule/internal/pairing/http.go:84:	mux.Handle("DELETE /pairing/v1/{realm}/agent/devices/{deviceID}",
/root/astrate-mule/internal/pairing/http.go:94:// --- flow A: agent endpoints ------------------------------------------------

→ Read internal/pairing/http.go [offset=70, limit=30]
→ Read docs/api/astarte_pairing_api.yaml
→ Read internal/pairing/http.go [offset=99, limit=270]

$ rg -n 'Cause|EXPIRED|REVOKED|INVALID|Status' /root/astrate-mule/internal/pairing/service.go | head -40
71:	// CauseExpired marks a certificate outside its validity window.
72:	CauseExpired = "EXPIRED"
73:	// CauseInvalid marks a certificate that does not parse or does not
75:	CauseInvalid = "INVALID"
76:	// CauseRevoked marks a certificate superseded by a newer issuance
78:	CauseRevoked = "REVOKED"
111:	// (docs/DESIGN.md §4.3): verify reports REVOKED for certificates whose
252:	if dev.Status == store.DeviceStatusInhibited {
276:	Status    string
282:// authenticated by credentials secret. Status strings are upstream parity:
293:	switch dev.Status {
294:	case store.DeviceStatusInhibited:
296:	case store.DeviceStatusConfirmed:
303:		Status:    status,
310:// Until carries the certificate expiry; otherwise Cause carries one of the
311:// Cause* constants. Timestamp is the verification instant.
316:	Cause     string
321:// validity window report EXPIRED; certificates that fail to parse or to
322:// chain to the realm CA report INVALID; valid-but-superseded certificates
323:// report REVOKED when EnforceLatestCert is on.
329:	if dev.Status == store.DeviceStatusInhibited {
342:		res.Cause = CauseExpired
345:		res.Cause = CauseInvalid
352:			res.Cause = CauseInvalid
356:			res.Cause = CauseRevoked

Now let me check the pairing service to verify the `cause` values and status behavior.
$ rg -n 'pairing|Pairing' /root/astrate-mule/internal/flowapi/ /root/astrate-mule/cmd/ 2>/dev/null | rg -i 'status|Version' | head
/root/astrate-mule/cmd/astrate/main.go:398:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(version))

→ Read internal/pairing/service.go [limit=50, offset=260]

$ rg -n 'confirmed|pending|denied|expired|inhibited' /root/astrate-mule/docs/site/ | head -30
/root/astrate-mule/docs/site/pairing-and-security.md:56:- First successful call stamps `first_credentials_request` and flips status `registered -> confirmed`.
/root/astrate-mule/docs/site/pairing-and-security.md:64:    "status": "confirmed",
/root/astrate-mule/docs/site/pairing-and-security.md:92:2. **confirmed** -- first credentials request succeeded.
/root/astrate-mule/docs/site/pairing-and-security.md:93:3. **inhibited** -- blocked from new credentials and connections (via `PATCH /appengine/v1/<realm>/devices/<id>` with `credentials_inhibited: true`).
/root/astrate-mule/docs/site/troubleshooting.md:20:**Symptom:** Device was connecting fine, now gets rejected. Logs show "certificate has expired".
/root/astrate-mule/docs/site/troubleshooting.md:31:**Symptom:** Device is rejected even though its cert is not expired. Logs show serial mismatch.
/root/astrate-mule/docs/site/troubleshooting.md:87:1. **401:** The JWT token is missing, expired, or signed with an unknown key. Generate a fresh token:
/root/astrate-mule/docs/site/troubleshooting.md:128:3. Verify the device's client cert hasn't expired (`pairing.cert_ttl`).
/root/astrate-mule/docs/site/observability.md:52:| `astrate_engine_trigger_deliveries_total` | counter | `outcome` | Trigger action deliveries by outcome: `delivered`, `failed`, `dropped`, `expired`, `forwarded`, `skipped`. |
/root/astrate-mule/docs/site/mqtt-protocol.md:28:Everything else is denied and logged. Server-side publishes use mochi's inline client and bypass ACLs.
/root/astrate-mule/docs/site/data-modeling.md:94:Device status values: `registered` (awaiting first credentials request), `confirmed` (has connected at least once), `inhibited` (blocked from new credentials and connections).
/root/astrate-mule/docs/site/appengine-api.md:32:Fields: `aliases`, `attributes`, `credentials_inhibited`.
/root/astrate-mule/docs/site/appengine-api.md:34:Setting `credentials_inhibited: true` blocks the device from obtaining new credentials and connecting to the broker.

Found drift in the status enum. Let me check whether `docs/site` prose repeats the wrong enum, and check the flow B rate-limit claim.
→ Read docs/api/astarte_pairing_api.yaml [offset=410, limit=10]
The status enum is genuinely wrong. Let me read the docs/site page to see if the prose contradicts the code too.
→ Read docs/site/pairing-and-security.md [limit=50, offset=50]

$ rg -n 'DeviceStatus' /root/astrate-mule/internal/store/*.go | rg -v '_test' | head -20
/root/astrate-mule/internal/store/devices.go:19:	// DeviceStatusRegistered marks a device that has a credentials secret
/root/astrate-mule/internal/store/devices.go:21:	DeviceStatusRegistered = "registered"
/root/astrate-mule/internal/store/devices.go:22:	// DeviceStatusConfirmed marks a device that has requested credentials
/root/astrate-mule/internal/store/devices.go:24:	DeviceStatusConfirmed = "confirmed"
/root/astrate-mule/internal/store/devices.go:25:	// DeviceStatusInhibited marks a device blocked from new credentials and
/root/astrate-mule/internal/store/devices.go:27:	DeviceStatusInhibited = "inhibited"

Found one more: `/pairing/v1/{realm}/version` is registered in main.go but absent from the yaml. Let me confirm it and check the observability handler.
$ rg -n 'pairing/v1/\{realm\}/version|/version' /root/astrate-mule/docs/api/ /root/astrate-mule/internal/observability/ /root/astrate-mule/cmd/astrate/main.go | head
/root/astrate-mule/cmd/astrate/main.go:388:	// service answers an unauthenticated GET /{service}/version with
/root/astrate-mule/cmd/astrate/main.go:396:	mux.Handle("GET /appengine/v1/{realm}/version",
/root/astrate-mule/cmd/astrate/main.go:398:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(version))
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml:532:  /realmmanagement/v1/{realm}/version:
/root/astrate-mule/internal/observability/compat.go:36:// VersionHandler answers GET .../version with the upstream envelope
/root/astrate-mule/internal/observability/compat.go:49:// MountVersionCompat registers GET /{service}/version for one service.
/root/astrate-mule/internal/observability/compat.go:51:	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
/root/astrate-mule/internal/observability/compat_version_test.go:19:		resp, err := http.Get(srv.URL + "/" + svc + "/version")
/root/astrate-mule/internal/observability/compat_version_test.go:21:			t.Fatalf("%s/version: %v", svc, err)
/root/astrate-mule/internal/observability/compat_version_test.go:26:			t.Errorf("%s/version = %d, want 200", svc, resp.StatusCode)

$ rg -n 'version' /root/astrate-mule/docs/site/pairing-and-security.md
(no output)

$ rg -n 'version' /root/astrate-mule/docs/api/astarte_appengine_api.yaml
9:  version: "0.0.1"
188:                  firmware_version: "2.1.0"
272:                  firmware_version: "2.1.0"
1013:                  firmware_version: "2.1.0"
1467:          description: Map of interface name to its major and minor version.

→ Read docs/api/astarte_realm_management_api.yaml [limit=25, offset=532]
→ Read internal/observability/compat.go [limit=40, offset=30]
→ Read cmd/astrate/main.go [limit=25, offset=384]
→ Read .mule/for-giulio.md
Now let me read the current state of the output files.

$ tail -30 /root/astrate-mule/.mule/todo.md
- [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
- [x] docs-sync-appengine-by-alias-endpoints [auto]: add the 7 undocumented by-alias mirror routes to docs/api/astarte_appengine_api.yaml — PATCH /devices-by-alias/{alias}, GET /devices-by-alias/{alias}/interfaces, GET/PUT/POST/DELETE /devices-by-alias/{alias}/interfaces/{interface}/{path} (plus the interface-scoped GET) — all registered in internal/appengine/http.go:59-65 and backed by handlers patchDeviceByAlias, listInterfacesByAlias, getDataByAlias, putDataByAlias, deleteDataByAlias (http.go:226-367). Match the existing spec style (same schemas, same response codes as the device-scoped counterparts). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-appengine-group-endpoints [auto]: add the 8 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (returns {data: {group_name: "<name>"}}, http.go:428-434), PATCH /groups/{group}/devices/{device} (DevicePatch body → updated status, http.go:518-537), and the 6 group interface-data routes: the interface-level GET and the path-level GET/PUT/POST/DELETE on /groups/{group}/devices/{device}/interfaces (http.go:75-80, handlers listInterfacesInGroup, getDataInGroup, putDataInGroup, deleteDataInGroup at http.go:187-377). Also add GET /devices/{device}/interfaces (listDeviceInterfaces, http.go:48) which is missing from the spec — that is the 9th route. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-appengine-get-group-device [auto]: add the undocumented `GET /groups/{group}/devices/{device}` route to docs/api/astarte_appengine_api.yaml — it exists in code (internal/appengine/http.go:72, handler getGroupDevice at http.go:509-516, returns the device status projection behind the membership gate) but the spec path only documents delete and patch. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups/{group}/devices is missing `from_token` and `limit` pagination params that code reads at http.go:439-441 (and returns via links.next), (2) POST /groups is missing 422 validation-error response (FieldErrors for blank group_name/empty devices, http.go:404-416), (3) PATCH /devices/{device} and PATCH /devices-by-alias/{alias} are missing 422 ("Attribute key not found" via ErrAttributeKeyNotFound, service.go:352-357 → http.go:656-657), and the device-scoped PATCH is also missing 409 ("Alias already in use", http.go:651) and its documented 500 covers the merge-patch content-type mismatch which code answers with an unmapped 500 too, (4) GET /devices/{device}/interfaces/{interface}/{path} is missing 422 (invalid since/since_after/to/limit/downsample_to/format params, http.go:554-639). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.
- [x] docs-sync-appengine-group-patch-status [auto]: PATCH /groups/{group}/devices/{device} (internal/appengine/http.go:518-537, patchGroupDevice) shares applyPatch with the device-scoped PATCH and also answers 422 ("Attribute key not found", ErrAttributeKeyNotFound, service.go:352-357) and 409 ("Alias already in use", http.go:651), but docs/api/astarte_appengine_api.yaml lists only 200/400/401/404/500 for it — add the two missing responses.
- [x] docs-sync-appengine-data-422-interface-level [auto]: the interface-level data GETs — GET /devices/{device}/interfaces/{interface}, GET /devices-by-alias/{alias}/interfaces/{interface}, GET /groups/{group}/devices/{device}/interfaces/{interface} — share parseQueryOpts (internal/appengine/http.go:554-639) and answer 422 FieldErrors for invalid since/since_after/to/limit/downsample_to/format params, but the spec documents only the path-level GET with 422 — add 422 to the interface-level GETs.
- [!] appengine-snapshot-ignores-query-params [auto]: in `internal/appengine/data.go:138-148`, the individual-datastream interface-root snapshot branch ignores `Since`/`SinceAfter`/`To`/`Limit`/`Descending` even though the root GET is documented with all of them (docs/api/astarte_appengine_api.yaml:577-582) — the same function refuses incompatible params elsewhere (downsample 422, data.go:118-120/165-167), so this is a silent drop. Probe upstream on the Legion Go to decide honour-vs-refuse, then either apply the window/limit or answer an explicit 422 and stop advertising the params on the root GET; add a T1 test in data_test.go/query_opts_test.go pinning the chosen behaviour. — BLOCKED: wrote nothing
- [!] appengine-group-token-roundtrip-test [auto]: add a container-free unit test for `parseGroupToken`/`groupTokenFor` (internal/appengine/service.go:602-623) — currently only exercised via the pg-backed groups_parity_test.go. Assert (1) round-trip `groupTokenFor(offset)` → `parseGroupToken` recovers offset, (2) a non-v1-nibble UUID is rejected, (3) a well-formed v1 UUID with zeroed time fields is accepted (parseGroupToken comment, service.go:614-616). — BLOCKED: lint failed: internal/appengine/groups_token_test.go:14:5: redefines-builtin-id: redefinition of the built-in function max (revive)
- [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
- [x] detailed-listing-required-encrypted-flags: emit the `required`/`encrypted` mapping flags in `writeDetailedMapping` (internal/realm/detailed.go) so `GET /realmmanagement/v1/{realm}/interfaces?detailed=true` round-trips them — the detailed listing currently materializes a fixed field set that drops both flags even when a datastream mapping declares one, while the plain GET preserves them via the stored definition. Datastream branch only (properties reject both at parse); add a T1 assertion in internal/realm/detailed_test.go.

- [x] store-alias-values-taken-test: add a store test for `AliasValuesTaken` (internal/store/devices.go:343, the `alias_already_in_use` ownership check behind PATCH /devices) pinning the three branches that currently have no direct coverage — no other device in the realm carries any of the values → false; another device carries one → true; the same device (excluded via `d.id <> $2`) carries the value → false (must not count itself). Add to internal/store/devices_test.go. [auto]
- [x] store-latest-individual-test: add a store test for `LatestIndividual` (internal/store/datastreams.go:247, the previous-value lookup feeding value_change/path_created triggers) — only exercised through the engine's fakeStore today; assert (a) of two samples on one series, `ORDER BY ts DESC LIMIT 1` returns the newest with the right value column, (b) an empty series yields ErrNotFound. Add to internal/store/datastreams_test.go. [auto]
- [x] docs-sync-appengine-device-status-schema [auto]: fix the `DeviceStatus` schema in docs/api/astrate_appengine_api.yaml — the wire emits `id` (internal/appengine/service.go:106 `json:"id"`, decoded as `ds.ID` in mirror_test.go:156 and dashboard_compat_test.go:100/149) but the spec names the field `device_id` (yaml:1409) in both the schema and the list/detail examples (yaml:67, 75); `introspection` is documented as `map<string,integer>` (yaml:1438-1441) but each value is a `{major, minor}` object (introspectionEntry, service.go:123-127, introspectionView service.go:833-836); and the schema omits the six fields `groups`, `total_received_msgs`, `total_received_bytes`, `first_credentials_request`, `last_seen_ip` (omitempty), `previous_interfaces` (omitempty) that deviceStatusView emits (service.go:254-279). Rename `device_id`→`id`, correct `introspection` to `map<string,{major:int,minor:int}>`, add the missing fields, and fix the examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-appengine-data-set-422 [auto]: add the missing `422` response to the six PUT/POST data-set operations in docs/api/astrate_appengine_api.yaml — `putData`, `putDataByAlias`, `putDataInGroup` (internal/appengine/http.go:292-309, 321-338, 340-357) return 422 `{"errors":{"detail":"Value size exceeds size limits"}}` when the body value exceeds 65536 bytes (maxValueBytes, http.go:24), both for PUT and its POST alias (routes registered at http.go:52-53, 63-64, 78-79), but the spec lists only 200/400/401/404/500 for all six operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-appengine-downsample-min [auto]: fix the `downsample_to` constraint in docs/api/astrate_appengine_api.yaml — the DataDownsample description says "(must be >= 2)" (yaml:1357) but parseQueryOpts rejects `downsample_to <= 2` with "must be greater than 2" (internal/appengine/http.go:592-593), so the documented minimum is off by one: readings of `2` (allowed per the doc) actually 422. Change the description to "(must be > 2)" (the 422 error shape is already documented). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

- [x] pairing-bearer-secret-test [auto]: add a container-free unit test for `bearerSecret` (internal/pairing/http.go:343) locking the upstream-verbatim permissive match `~r/bearer\:?\s+(.*)$/i` — the colon form `Bearer: <secret>`, the case variants (`bearer `, `BEARER `), and the reject paths (no whitespace before value, non-bearer scheme, empty value). Only the plain `Bearer <secret>` form is exercised today (integration http_test.go); a regression to case-sensitive or drop-colon matching would pass every existing test. Pure helper in a non-integration file; no Docker.
- [x] pairing-remoteip-fallback-test [auto]: add a container-free unit test for `remoteIP` (internal/pairing/http.go:359) — a valid addr:port yields `ap.Addr().Unmap()`, a bare address yields the mapped address, and an unparseable `RemoteAddr` (unix-socket/proxy case) falls back to `netip.IPv4Unspecified()`. That fallback branch is currently uncovered (only reachable through the integration HTTP stack). Pure helper in a non-integration file; no Docker.
- [x] docs-sync-hk-patch-endpoint [auto]: add the undocumented `PATCH /housekeeping/v1/realms/{realm}` route to docs/api/astarte_housekeeping_api.yaml — it exists in code (internal/housekeeping/http.go:40, handler patchRealm at http.go:163-229, returns 200 on success) but the spec jumps straight from GET to DELETE. Document the request body (`RealmPatch` — same three optional fields as `RealmCreate`: `jwt_public_key_pem`, `device_registration_limit`, `datastream_maximum_storage_retention`, all nullable), the success response, and the error responses: 400 (bad envelope), 401, 403, 404, 405 (deletion disabled), 422 (invalid update parameters for unknown fields, blank jwt_public_key_pem, negative limit/retention, connected devices present), 500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-retention-field [auto]: add the `datastream_maximum_storage_retention` field to the `Realm` and `RealmCreate` schemas in docs/api/astarte_housekeeping_api.yaml — the field is on the wire (realmBody struct, internal/housekeeping/http.go:50, returned by viewBody at http.go:99, accepted in create at http.go:68 and patch at http.go:216) but absent from both schemas (yaml:200-231). Type: integer, format int64, nullable, description matching the device_registration_limit style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

- [x] realm-pure-helper-tests [auto]: add container-free unit tests for the four untested pure helpers in internal/realm. Create `internal/realm/helpers_test.go` (no build tag) covering: (1) `normaliseIfaceName` — lowercase, hyphen strip, mixed case+hyphens, empty string, already-normalized passthrough; (2) `joinPEM` — empty slice, single key, two keys separated by `\n`; (3) `triggerErrorBody` — empty action+simple_triggers, action-only, simple_triggers-only with index-aligned `{}` substitution, both present; (4) `majorParam` — valid positive, zero, negative (→404), non-numeric (→404). All tests use httptest.ResponseRecorder for majorParam; the rest are direct function calls. No Docker.

- [x] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmmanagement/health`, `GET /pairing/health` — in docs/api/astrate_native_api.yaml, and widen the shared `DataEnvelopeHealthStatus.status` enum (yaml:301-303) from `[ok]` to `[ok, unhealthy]` — the handler answers 503 `{"data":{"status":"unhealthy"}}` when the readiness probe fails (internal/observability/compat.go:19-34, deviation 18), but the spec documents only 200 with `status: ok` for all three paths. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

→ Read .mule/todo.md [limit=30]

$ git -C /root/astrate-mule branch --show-current 2>/dev/null; git -C /root/astrate-mule status --short 2>/dev/null | head
mule/queue

$ rg -n 'pairing' /root/astrate-mule/.mule/todo.md
91:- [ ] race-check-flow: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/flow/... ./internal/realm/... ./internal/pairing/... ./internal/auth/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
104:- [x] docs-sync-pairing-health-path [auto]: add the undocumented `GET /pairing/v1/{realm}/health` route to docs/api/astarte_pairing_api.yaml — it exists in code (internal/pairing/http.go:81) and returns 200 `{"data":{"status":"ok"}}` for an existing realm, 404 for an unknown realm, 503 when unhealthy, unauthenticated and unrate-limited. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
105:- [x] docs-sync-pairing-register-404 [auto]: document the missing `404` (DeviceNotFound) on `POST /pairing/v1/{realm}/agent/devices` in docs/api/astarte_pairing_api.yaml — the handler returns it when `GetRealmByName` fails for an unknown realm (internal/pairing/service.go:183-186 → http.go:326-327), but the spec lists only 201/400/401/403/422/429/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
135:- [x] pairing-bearer-secret-test [auto]: add a container-free unit test for `bearerSecret` (internal/pairing/http.go:343) locking the upstream-verbatim permissive match `~r/bearer\:?\s+(.*)$/i` — the colon form `Bearer: <secret>`, the case variants (`bearer `, `BEARER `), and the reject paths (no whitespace before value, non-bearer scheme, empty value). Only the plain `Bearer <secret>` form is exercised today (integration http_test.go); a regression to case-sensitive or drop-colon matching would pass every existing test. Pure helper in a non-integration file; no Docker.
136:- [x] pairing-remoteip-fallback-test [auto]: add a container-free unit test for `remoteIP` (internal/pairing/http.go:359) — a valid addr:port yields `ap.Addr().Unmap()`, a bare address yields the mapped address, and an unparseable `RemoteAddr` (unix-socket/proxy case) falls back to `netip.IPv4Unspecified()`. That fallback branch is currently uncovered (only reachable through the integration HTTP stack). Pure helper in a non-integration file; no Docker.
142:- [x] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmmanagement/health`, `GET /pairing/health` — in docs/api/astrate_native_api.yaml, and widen the shared `DataEnvelopeHealthStatus.status` enum (yaml:301-303) from `[ok]` to `[ok, unhealthy]` — the handler answers 503 `{"data":{"status":"unhealthy"}}` when the readiness probe fails (internal/observability/compat.go:19-34, deviation 18), but the spec documents only 200 with `status: ok` for all three paths. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
143:- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

→ Read .mule/todo.md [offset=135]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -139,5 +139,7 @@
 
 - [x] realm-pure-helper-tests [auto]: add container-free unit tests for the four untested pure helpers in internal/realm. Create `internal/realm/helpers_test.go` (no build tag) covering: (1) `normaliseIfaceName` — lowercase, hyphen strip, mixed case+hyphens, empty string, already-normalized passthrough; (2) `joinPEM` — empty slice, single key, two keys separated by `\n`; (3) `triggerErrorBody` — empty action+simple_triggers, action-only, simple_triggers-only with index-aligned `{}` substitution, both present; (4) `majorParam` — valid positive, zero, negative (→404), non-numeric (→404). All tests use httptest.ResponseRecorder for majorParam; the rest are direct function calls. No Docker.
 
 - [x] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmmanagement/health`, `GET /pairing/health` — in docs/api/astrate_native_api.yaml, and widen the shared `DataEnvelopeHealthStatus.status` enum (yaml:301-303) from `[ok]` to `[ok, unhealthy]` — the handler answers 503 `{"data":{"status":"unhealthy"}}` when the readiness probe fails (internal/observability/compat.go:19-34, deviation 18), but the spec documents only 200 with `status: ok` for all three paths. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-pairing-status-enum [auto]: fix the `PairingInfo.status` enum in docs/api/astarte_pairing_api.yaml from `[confirmed, pending, denied, expired]` (yaml:412) to `[confirmed, pending, inhibited]` — the handler never emits `denied` or `expired`: `service.Info` maps only `DeviceStatusInhibited`→"inhibited", `DeviceStatusConfirmed`→"confirmed", and everything else→"pending" (internal/pairing/service.go:296-299, statuses defined in store/devices.go:19-27; upstream-parity per the comment at service.go:282-285). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-pairing-version-endpoint [auto]: add the undocumented `GET /pairing/v1/{realm}/version` route to docs/api/astarte_pairing_api.yaml — it exists in code (cmd/astrate/main.go:398, unauthenticated, upstream-parity per main.go:387-391) and answers 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38); mirror the RM spec's getVersion op (astarte_realm_management_api.yaml:532) without its a_rma security. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

Both findings confirmed. Now appending the two proposed tasks and the escalation.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **docs-sync pairing, 2026-09-10: docs/site names the pre-credentials wire status `registered`; the API returns `pending`.** `docs/site/pairing-and-security.md:56` ("flips status `registered -> confirmed`") and `:91` ("**registered** -- device registered, awaiting first credentials request"), plus `docs/site/data-modeling.md:94`, all give `registered` as a device status value. But `service.Info` emits `pending` for any device that is neither confirmed nor inhibited (internal/pairing/service.go:296-299 — upstream-parity per the comment at service.go:282-285; the DB value is `registered`, store/devices.go:19-27). Wire values are `pending`/`confirmed`/`inhibited`. Site prose is yours — reword to `pending`, or confirm the site intentionally describes the DB value.
+
+---
+
 - **Flow code review, 2026-09-10: `DeletePipeline` semantics need your call.**
   `DeletePipeline` (internal/flowapi/service.go:199-206) removes a pipeline
   silently even when durable flows still reference it; running instances keep
   the old graph, but their next reload fails ("pipeline not found") and a

$ tail -6 /root/astrate-mule/.mule/todo.md; echo '---'; sed -n '11,20p' /root/astrate-mule/.mule/for-giulio.md
- [x] realm-pure-helper-tests [auto]: add container-free unit tests for the four untested pure helpers in internal/realm. Create `internal/realm/helpers_test.go` (no build tag) covering: (1) `normaliseIfaceName` — lowercase, hyphen strip, mixed case+hyphens, empty string, already-normalized passthrough; (2) `joinPEM` — empty slice, single key, two keys separated by `\n`; (3) `triggerErrorBody` — empty action+simple_triggers, action-only, simple_triggers-only with index-aligned `{}` substitution, both present; (4) `majorParam` — valid positive, zero, negative (→404), non-numeric (→404). All tests use httptest.ResponseRecorder for majorParam; the rest are direct function calls. No Docker.

- [x] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmmanagement/health`, `GET /pairing/health` — in docs/api/astrate_native_api.yaml, and widen the shared `DataEnvelopeHealthStatus.status` enum (yaml:301-303) from `[ok]` to `[ok, unhealthy]` — the handler answers 503 `{"data":{"status":"unhealthy"}}` when the readiness probe fails (internal/observability/compat.go:19-34, deviation 18), but the spec documents only 200 with `status: ok` for all three paths. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-pairing-status-enum [auto]: fix the `PairingInfo.status` enum in docs/api/astarte_pairing_api.yaml from `[confirmed, pending, denied, expired]` (yaml:412) to `[confirmed, pending, inhibited]` — the handler never emits `denied` or `expired`: `service.Info` maps only `DeviceStatusInhibited`→"inhibited", `DeviceStatusConfirmed`→"confirmed", and everything else→"pending" (internal/pairing/service.go:296-299, statuses defined in store/devices.go:19-27; upstream-parity per the comment at service.go:282-285). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-pairing-version-endpoint [auto]: add the undocumented `GET /pairing/v1/{realm}/version` route to docs/api/astarte_pairing_api.yaml — it exists in code (cmd/astrate/main.go:398, unauthenticated, upstream-parity per main.go:387-391) and answers 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38); mirror the RM spec's getVersion op (astarte_realm_management_api.yaml:532) without its a_rma security. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
---
---

- **docs-sync pairing, 2026-09-10: docs/site names the pre-credentials wire status `registered`; the API returns `pending`.** `docs/site/pairing-and-security.md:56` ("flips status `registered -> confirmed`") and `:91` ("**registered** -- device registered, awaiting first credentials request"), plus `docs/site/data-modeling.md:94`, all give `registered` as a device status value. But `service.Info` emits `pending` for any device that is neither confirmed nor inhibited (internal/pairing/service.go:296-299 — upstream-parity per the comment at service.go:282-285; the DB value is `registered`, store/devices.go:19-27). Wire values are `pending`/`confirmed`/`inhibited`. Site prose is yours — reword to `pending`, or confirm the site intentionally describes the DB value.

---

- **Flow code review, 2026-09-10: `DeletePipeline` semantics need your call.**
  `DeletePipeline` (internal/flowapi/service.go:199-206) removes a pipeline
  silently even when durable flows still reference it; running instances keep
  the old graph, but their next reload fails ("pipeline not found") and a

Done: docs-sync recipe on the pairing surface — found one spec enum bug, one missing route, one prose contradiction.
Files: .mule/todo.md (2 new task lines), .mule/for-giulio.md (1 escalation)
Verified: cross-checked yaml paths/methods/status/fields vs internal/pairing/http.go:77-92, service.go:286-307, store/devices.go:19-27, compat.go:38, main.go:393-398 -> all other endpoints match
Unsure: nothing (DB value `registered` vs wire `pending` noted in the escalation for Giulio's call)
Follow-ups:
- docs-sync-pairing-status-enum: status enum `[confirmed, pending, denied, expired]` -> the wire only emits `pending`/`confirmed`/`inhibited` (service.go:296-299); denied/expired never produced, inhibited missing.
- docs-sync-pairing-version-endpoint: `GET /pairing/v1/{realm}/version` exists in code (main.go:398) but is absent from the yaml.
- for-giulio.md: docs/site names the pre-credentials wire status `registered`; API returns `pending`.
