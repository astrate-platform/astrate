slug: recipe-docs-sync
verdict: proposed
at:  b200577
ran: 2026-09-21T17:46:46Z on DietPi in 328s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls docs/api/ 2>/dev/null; echo ---; ls internal/ 2>/dev/null
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

$ ls docs/site/ 2>/dev/null; echo ---; ls .mule/
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

$ rg -n '^\s{2}/' docs/api/astrate_native_api.yaml | head -80
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

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*flow*' --glob '*broker*' --glob '*engine*' | head -60
(no output)

$ rg -rn 'astrate/v1|appengine/health|realmmanagement|pairing/health|/socket' internal/ | head -40
internal/appengine/stream/ws_test.go:		req := httptest.NewRequest(http.MethodGet, "/n/testrealmn", nil)
internal/appengine/stream/ws_test.go:		req := httptest.NewRequest(http.MethodGet, "/n/testrealmn?device_id=dev1&interface=com.ex.S", nil)
internal/appengine/stream/ws.go:// at /n/{realm}n fed by the engine's in-process fan-out bus.
internal/appengine/stream/ws.go:	mux.Handle("GET /n/{realm}n", a.require(http.HandlerFunc(a.handle)))
internal/appengine/channels/ws_test.go:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1n/websocket?realm=testrealm", nil)
internal/appengine/channels/ws_test.go:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1n/websocket?token="+token, nil)
internal/appengine/channels/ws_test.go:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1n/websocket?realm=nosuchrealm&token="+token, nil)
internal/appengine/channels/ws_test.go:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1n/websocket?realm=testrealm&token=not.a.jwt", nil)
internal/appengine/channels/ws_test.go:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1n/websocket?realm=testrealm&token="+token, nil)
internal/appengine/channels/ws_test.go:		"/appengine/v1n/websocket?vsn=2.0.0&realm=" + testRealmName + "&token=" + token
internal/appengine/channels/protocol.go:// served at /appengine/v1n/websocket. It is the compatibility counterpart
internal/appengine/channels/ws.go:	mux.Handle("GET /appengine/v1n/websocket", http.HandlerFunc(a.handle))
internal/config/config.example.toml:# and the Astrate-native /n health/metricsn endpoints.
internal/realm/http.go:// API is the /n/v1 HTTP surface (docs/ROADMAP.md §8.1 file
internal/realm/http.go:	mux.Handle("GET /n/v1/{realm}/interfaces", h(a.listInterfaces))
internal/realm/http.go:	mux.Handle("POST /n/v1/{realm}/interfaces", h(a.installInterface))
internal/realm/http.go:	mux.Handle("GET /n/v1/{realm}/interfaces/{name}", h(a.listInterfaceMajors))
internal/realm/http.go:	mux.Handle("GET /n/v1/{realm}/interfaces/{name}/{major}", h(a.getInterface))
internal/realm/http.go:	mux.Handle("PUT /n/v1/{realm}/interfaces/{name}/{major}", h(a.updateInterface))
internal/realm/http.go:	mux.Handle("DELETE /n/v1/{realm}/interfaces/{name}/{major}", h(a.deleteInterface))
internal/realm/http.go:	mux.Handle("GET /n/v1/{realm}/triggers", h(a.listTriggers))
internal/realm/http.go:	mux.Handle("POST /n/v1/{realm}/triggers", h(a.createTrigger))
internal/realm/http.go:	mux.Handle("GET /n/v1/{realm}/triggers/{name}", h(a.getTrigger))
internal/realm/http.go:	mux.Handle("DELETE /n/v1/{realm}/triggers/{name}", h(a.deleteTrigger))
internal/realm/http.go:	mux.Handle("GET /n/v1/{realm}/config/auth", h(a.getAuth))
internal/realm/http.go:	mux.Handle("PUT /n/v1/{realm}/config/auth", h(a.putAuth))
internal/realm/http.go:	mux.Handle("GET /n/v1/{realm}/config/device_registration_limit", h(a.getRegistrationLimit))
internal/realm/http.go:	mux.Handle("GET /n/v1/{realm}/config/datastream_maximum_storage_retention", h(a.getDatastreamMaximumStorageRetention))
internal/realm/http.go:	mux.Handle("GET /n/v1/{realm}/version", h(a.getVersion))
internal/realm/http.go:	mux.Handle("DELETE /n/v1/{realm}/devices/{device}", h(a.deleteDevice))
internal/realm/http.go:	mux.Handle("GET /n/v1/{realm}/policies", h(a.listPolicies))
internal/realm/http.go:	mux.Handle("POST /n/v1/{realm}/policies", h(a.createPolicy))
internal/realm/http.go:	mux.Handle("GET /n/v1/{realm}/policies/{name}", h(a.getPolicy))
internal/realm/http.go:	mux.Handle("DELETE /n/v1/{realm}/policies/{name}", h(a.deletePolicy))
internal/realm/http_test.go:	httpReq := httptest.NewRequest(method, "/n/v1/"+r.realm+path, body)
internal/observability/compat.go:// (appengine, n, pairing) for its API status indicators; only
internal/observability/compat.go:// same dependency probe that backs /n/readiness — and answers 503
internal/observability/observability_test.go:	m.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/n/metrics", nil))
internal/observability/observability_test.go:	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/n/health", nil))
internal/observability/observability_test.go:	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/n/readiness", nil))

$ rg -rln 'astrate/v1' . --glob '!docs/**' --glob '!.mule/**' --glob '!go.sum' | head
./README.md:curl localhost:8080/ln/readiness        # {"status":"ok",...}
./README.md:curl localhost:8080/ln/metrics          # Prometheus metrics
./cmd/astrate/main.go:	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+addr+"/ln/readiness", nil) //nolint:gosec // G704: self-probe of the local readiness endpoint
./cmd/astrate/main_test.go:		if code, _ := get(t, base+"/ln/health"); code != http.StatusOK {
./cmd/astrate/main_test.go:		code, body := get(t, base+"/ln/metrics")
./cmd/astrate/main_test.go:		if code, _ := get(t, base+"/ln/readiness"); code == http.StatusOK {
./bench/scripts/up-astrate.sh:    if curl -fsS http://127.0.0.1:8080/ln/readiness >/dev/null 2>&1; then
./examples/pokemon-agent/docs/TESTING.md:# readiness: curl -s localhost:8080/ln/readiness
./examples/pokemon-agent/docs/TESTING.md:# GET /ln/$REALM/socket?device_id=$DEVICE_ID&interface=org.pokemon.emulator.GameState
./examples/pokemon-agent/docs/TESTING.md:Astrate’s live stream (`GET /ln/{realm}/socket`) authorizes with

→ Read internal/realm/http.go [limit=60]
→ Read docs/api/astrate_native_api.yaml
→ Read internal/appengine/stream/ws.go
→ Read internal/appengine/channels/ws.go
→ Read internal/observability

$ rg -rn 'health|readiness|metrics|version' internal/observability/compat.go | head -40
// GET /{service}/n endpoint. The Astarte Dashboard polls one per service
// same dependency probe that backs /astrate/v1/n — and answers 503
	mux.HandleFunc("GET /"+service+"/n", func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), nTimeout)
				_, _ = w.Write([]byte(`{"data":{"status":"unny"}}`))
// VersionHandler answers GET .../n with the upstream envelope
// {"data":"<n>"} (measured upstream 1.2.0, verify-ns.json).
func VersionHandler(n string) http.HandlerFunc {
	body, err := json.Marshal(map[string]string{"data": n})
// MountVersionCompat registers GET /{service}/n for one service.
func MountVersionCompat(mux *http.ServeMux, service, n string) {
	mux.HandleFunc("GET /"+service+"/n", VersionHandler(n))

→ Read internal/observability/compat.go
→ Read internal/observability/health.go
→ Read internal/observability/metrics.go

$ rg -rn 'WriteUnauthorized|WriteForbidden|WriteInternalServerError' pkg/astarteapi/ internal/auth/ --glob '*.go' | head -30
pkg/astarteapi/envelope.go:// n writes the canonical 401 envelope.
pkg/astarteapi/envelope.go:func n(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:// n writes the canonical 403 envelope.
pkg/astarteapi/envelope.go:func n(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:// n writes the canonical 500 envelope.
pkg/astarteapi/envelope.go:func n(w http.ResponseWriter) error {
internal/auth/middleware.go:				_ = astarteapi.n(w)
internal/auth/middleware.go:				_ = astarteapi.n(w)
internal/auth/middleware.go:				_ = astarteapi.n(w)
internal/auth/middleware.go:			_ = astarteapi.n(w)
internal/auth/middleware.go:		_ = astarteapi.n(w)
internal/auth/middleware.go:		_ = astarteapi.n(w)
internal/auth/middleware.go:		_ = astarteapi.n(w)
pkg/astarteapi/envelope_test.go:			write: astarteapi.n,
pkg/astarteapi/envelope_test.go:			write: astarteapi.n,
pkg/astarteapi/envelope_test.go:			write: astarteapi.n,

$ rg -rn 'Kind|incoming_data|device_connected|device_disconnected|intf' internal/engine/stream/*.go 2>/dev/null | head -30
internal/engine/stream/bus.go:	// nIncomingData is a committed device data operation (set or unset).
internal/engine/stream/bus.go:	nIncomingData = "n"
internal/engine/stream/bus.go:	// nDeviceRegistered is a device registration via the Pairing API.
internal/engine/stream/bus.go:	nDeviceRegistered = "device_registered"
internal/engine/stream/bus.go:	// nDeviceConnected is a device connection.
internal/engine/stream/bus.go:	nDeviceConnected = "n"
internal/engine/stream/bus.go:	// nDeviceDisconnected is a device disconnection.
internal/engine/stream/bus.go:	nDeviceDisconnected = "n"
internal/engine/stream/bus.go:	// nDeviceDeletionStarted is emitted immediately before a synchronous
internal/engine/stream/bus.go:	nDeviceDeletionStarted = "device_deletion_started"
internal/engine/stream/bus.go:	// nDeviceDeletionFinished is emitted immediately after a synchronous
internal/engine/stream/bus.go:	nDeviceDeletionFinished = "device_deletion_finished"
internal/engine/stream/bus.go:	// nDeviceError is a rejected device message (docs/DESIGN.md §2.6
internal/engine/stream/bus.go:	nDeviceError = "device_error"
internal/engine/stream/bus.go:	// n discriminates the event (n* constants).
internal/engine/stream/bus.go:	n string
internal/engine/stream/bus.go:	// IP is the peer address of a n event, empty otherwise.
internal/engine/stream/bus_test.go:// dataEvent builds an n event for one realm/device/interface.
internal/engine/stream/bus_test.go:		n: nIncomingData, Realm: realm, DeviceID: device,
internal/engine/stream/bus_test.go:		b.Publish(Event{n: nDeviceConnected, Realm: "r", DeviceID: "d"})
internal/engine/stream/bus_test.go:		b.Publish(Event{n: nDeviceConnected, Realm: "r", DeviceID: "d"})
internal/engine/stream/bus_test.go:		if ev := recv(t, ch); ev.n != nDeviceConnected {
internal/engine/stream/bus_test.go:			t.Errorf("got kind %q", ev.n)

→ Read internal/engine/stream/bus.go [limit=70]
→ Read pkg/astarteapi/envelope.go

$ rg -rn 'MountServiceCompat|MountVersionCompat|Health|NewHealth|Mount\(|\.Mount' cmd/astrate/main.go | head -30
os.Exit(runncheck())
	})n(mux)
	housekeeping.NewAPI(hkSvc, mw, hkKeys)n(mux)
	realm.NewAPI(realmSvc, mw)n(mux)
	appengine.NewAPI(appengine.NewService(st, e, log), mw)n(mux)
	apstream.NewAPI(e.Bus(), mw)n(mux)
	channels.NewAPI(e.Bus(), st, st)n(mux)
	flowapi.NewAPI(flowSvc, mw)n(mux)
		observabilitynServiceCompat(mux, svc, st.n)
		observabilitynVersionCompat(mux, svc, version)
	health := observability.n(metrics.Handler())
	health.AddReadiness("database", st.n)
	healthn(mux)
	swaggern(mux)
// runncheck probes the local readiness endpoint and returns a process
func runncheck() int {

→ Read internal/auth/middleware.go
The bash tool is heavily redacting output; the Read/Grep tools render content faithfully. Let me use those instead.
✱ Grep "MountServiceCompat|MountVersionCompat|NewHealth|\.Mount\(|AddReadiness" in cmd/astrate · 14 matches
→ Read cmd/astrate/main.go [offset=350, limit=70]
✱ Grep "ClaimChannels|ClaimRealmManagement|a_ch" in internal/auth · 25 matches
✱ Grep "ASTRATE_" in docs/site/configuration-reference.md · 41 matches
✱ Grep ""ASTRATE_[A-Z_]+"" in internal · 30 matches
✱ Grep "housekeeping|HOUSEKEEPING" in docs/site/configuration-reference.md · 22 matches
→ Read internal/config/config.go [offset=258, limit=50]
→ Read docs/site/configuration-reference.md [limit=40, offset=76]
✱ Grep "DefaultDatastreamMaximumStorageRetention|RealmDeletionDisabled|housekeeping" in internal/config · 23 matches
→ Read .mule/for-giulio.md
Let me check the current state of for-giulio.md and todo.md before writing anything.

$ tail -20 .mule/todo.md
- [x] astarteapi-metadata-envelope-test [auto]: add a container-free golden test for `WriteDataWithMetadata` (pkg/astarteapi/envelope.go:154-158) pinning the two wire branches — nil metadata omits the `metadata` key (`{"data": v}` only, via the `omitempty` at envelope.go:151), populated metadata emits `{"data": v, "metadata": {...}}` — with 200 + the package `Content-Type`, matching the TestGoldenEnvelopes style. Today the non-nil branch is only covered by the integration-tagged formats_parity_test.go (runs on the Legion Go, not the Pi gate) and the nil-omission branch is covered nowhere (internal/appengine/data.go:262-268 always sets metadata), so a regression dropping the omitempty passes every suite on this box.
- [x] astarteapi-multikey-sorted-golden [auto]: extend TestGoldenEnvelopes (pkg/astarteapi/envelope_test.go:101-121) with multi-key `WriteFieldErrors` and `WriteRawErrors` maps asserting the sorted key order on the wire — the docs claim "keys are emitted in sorted order" (envelope.go:170, 184) but every fixture uses a single key, so a change to the marshalling path that stopped sorting (e.g. switching encoders) passes the suite.
- [x] docs-sync-pairing-deviceendpoints-dead-404-403 [auto]: remove the unreachable `404` DeviceNotFound response from the three device-facing ops in docs/api/astarte_pairing_api.yaml — `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials` (yaml:226-227), `GET /pairing/v1/{realm}/devices/{deviceID}` (yaml:268-269), `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify` (yaml:339-340) — and the `403` Forbidden from the info GET (yaml:266-267). `Service.authenticateDevice` folds every failure (malformed/unknown deviceID, unknown realm, unregistered device, wrong secret) into uniform `ErrUnauthorized` (internal/pairing/service.go:380-411), so `store.ErrNotFound` never reaches `writeServiceError` (http.go:326-327) from those handlers; and `Service.Info` never returns `ErrInhibited` (inhibited devices still read info with 200 + status "inhibited", service.go:286-307, comment at 282-285). Reachable on these three are only 200/201/400/401/422/429/500 plus 403 on credentials and verify (both genuinely raise `ErrInhibited`, service.go:252-254 and 329-331). Keep 404 on the agent endpoints (register yaml:149-150, unregister yaml:176-177) where `GetRealmByName`/`deviceid.Parse` really do return `store.ErrNotFound`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-version-404-unreachable [auto]: remove the `404` "Unknown realm" response from `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml (yaml:96-104) — the route is served by `observability.VersionHandler` (cmd/astrate/main.go:398, internal/observability/compat.go:38-47), which answers static `200 {"data":"<version>"}` with no realm lookup, so a 404 can never be emitted there. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] observability-readiness-wedge-test [auto]: make the readiness budget injectable and pinned — add a `timeout time.Duration` field on `Health` (internal/observability/health.go) defaulting to `readinessTimeout` (health.go:12, no test currently exercises the wedged-dependency path), and add a test in observability_test.go asserting a ctx-honoring wedged check makes `/astrate/v1/readiness` return 503 reporting it failing within a short injected budget, instead of hanging the endpoint.
- [x] observability-dbpool-gauge-coverage [auto]: extend `TestMetricsExposesGauges` (internal/observability/observability_test.go:25-34) to assert all four `astrate_db_pool_*` gauges — it currently pins only `acquired_conns` and `max_conns`, so a regression deleting or mis-wiring `idle_conns` or `total_conns` (metrics.go:77-78) passes the suite; expect idle 5, total 7 from the supplied `DBPoolStats{2,5,7,10}`.
- [x] observability-health-content-type-test [auto]: assert `Content-Type: application/json` on `/astrate/v1/health` and `/astrate/v1/readiness` (set by `writeJSON`, health.go:81) and on both the 200 and 503 paths of the per-service `/health` (set by hand in compat.go:26,31) — only the version endpoint pins the header today (compat_version_test.go:31-33), so a regression dropping the header on the /astrate/v1 surface passes the suite.

- [x] docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] httpx-cors-vary-origin-passthrough [auto]: in internal/httpx/cors.go non-wildcard mode, the disallowed-origin pass-through (cors.go:26-29) returns without `Vary: Origin`, so responses whose CORS outcome depends on the Origin header reach shared caches without the RFC 9110 §12.5.5 variation marker (a cache could serve the CORS-approving variant to a disallowed origin, or the headerless variant to an allowed one); the stamp lives only in the allowed branch (cors.go:35-36) and only the preflight pins it (cors_test.go:42). Add `h.Add("Vary", "Origin")` on the pass-through branch too, and assert `Vary: Origin` in TestCORSActualRequestStamped (actual non-preflight, pinned only for Allow-Origin today) and TestCORSDisallowedOriginPassesThroughUnstamped. Additive header change; container-free.
- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] channels-group-watch-membership-scope: group watches authorize against `WATCH::groups/<name>/…` (ws.go:178-205, conformance channels.json) but `Room.dispatch` ever only filters by `watchEntry.deviceID` (empty for a group watch — room.go:36-40, 196-200) and the trigger compiles with no group/device scope (triggers/match.go:374-396, because upstream carries group_name at the payload top level), so a group-scoped token receives every realm event matching interface/path/on, members or not. Plumb a group-membership resolver into internal/appengine/channels (it knows only `Bus` today), resolve membership at Watch time and store it in `watchEntry`, filter deliveries by it. Add a container-free test with a fake resolver: watch group `probe`, publish an event from a member and a non-member, assert only the member's arrives. Verify the member-filter rule against upstream on the Legion if disputed; the boundary leak stands regardless.
- [x] channels-rejoin-joinref-tagging-test: pin in internal/appengine/channels/ws_test.go that a second `phx_join` on the same topic updates the session's `joinRef` (ws.go:284-291) and subsequent `new_event` pushes carry the fresh ref (read per event at ws.go:403) — no test drives a rejoin today (all `TestWireSession_*` join once), so the comment's "server pushes must be tagged with the current one" rots unasserted. Wire: join with ref "1", rejoin with "2", push an event, assert element 0 is `"2"`.
- [!] channels-rejoin-authz-mismatch: `handleJoin`'s comment says upstream authorizes every join including a rejoin (ws.go:276-277), but the rejoin branch (ws.go:284-291) replies OK without calling `s.tok.AuthorizesChannel(auth.VerbJoin, name)` — behaviourally identical today (a token is immutable for the session) but the comment lies about what the code does. Either move the authz check to cover the rejoin branch too, or fix the comment; mechanical, no behaviour change, choose one. — BLOCKED: wrote nothing
- [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a locally-run `elixir -e` with `Base.url_decode64!`) with the device ID `"AAAAAAAAAAAAAAAAAAAAAB"` — nonzero unused bits in the 22nd base64url char — and report whether upstream accepts it: `pkg/deviceid/deviceid.go:34-37` uses `base64.RawURLEncoding.Strict()` and deviceid_test.go:176 asserts rejection while the comment claims parity with `Elixir Base.url_decode64!(padding: false)`, which decodes and discards those bits rather than erroring. Astrate may therefore be stricter than upstream on the same wire form (a device registered with such an id on upstream would be rejected here). If upstream rejects it identically, the strict encoding is verified parity and the task is done; if upstream accepts it, escalate the relax-vs-strict call to `.mule/for-giulio.md` with the measurement. Probe first, no code change either way.
- [x] docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate_devices_total{realm="test"} 42`, but no such series is ever registered — the registry carries `astrate_broker_sessions` (internal/observability/metrics.go:54-62), `astrate_db_pool_{acquired,idle,total,max}_conns` (metrics.go:65-79), plus the engine/flow/trigger families (`astrate_engine_*` engine/router.go:411-444, `astrate_engine_stream_dropped_total` engine/stream/bus.go:119, `astrate_engine_trigger_*` engine/triggers/actions.go:296-300, `astrate_flow_router_*` flow/router.go:260-276) and the standard `go_*`/`process_*` runtime collectors — so the example advertises a scrape result the wire can never return. Swap it for a real series (e.g. `astrate_broker_sessions 0`), keeping the HELP/TYPE comment style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] store-devices-inhibit-re-register [legion] [auto]: `RegisterDevice` (internal/store/devices.go:75-91) silently clears the inhibit flag — its `ON CONFLICT ... DO UPDATE SET status = 'registered'` fires for any device with `first_credentials_request IS NULL`, including one the admin inhibited via `SetDeviceInhibited` (devices.go:251-268), which sets `status='inhibited'` on unconfirmed devices too; §5.3 says an inhibited device blocks new credentials and connections, so the re-registration re-opens it. Preserve `'inhibited'` in the SET (e.g. `status = CASE WHEN devices.status = 'inhibited' THEN 'inhibited' ELSE 'registered' END`) and add a Lifecycle case in internal/store/devices_test.go: inhibit an unconfirmed device, re-register, assert status stays inhibited and the secret still rotates. Verify the assert against upstream's register-not-touching-inhibit on the Legion while the integration suite runs.
- [x] store-pipelines-empty-name-zero-blocks-test [auto]: extend internal/store/pipelines_validate_test.go with the two untested error branches of `validatePipelineGraph` (internal/store/pipelines.go:51-64) — a definition with zero blocks and one whose block has an empty name must both be rejected. Pure helper, no DB; the existing suite covers only acyclic/cyclic/unknown-ref/duplicate-name.
- [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (internal/store/devices.go:130-131, "if several devices share an alias the lowest device ID wins", `ORDER BY id LIMIT 1`) — no test drives it today; add a case in internal/store/devices_test.go that gives two devices the same alias and asserts the lookup resolves the lower ID. Needs the integration DB.

Good context. Previous runs already synced this exact surface (the `[x]` lines referencing astrate_native_api.yaml). Let me check for any *present-but-undocumented* routes and verify the previously-proposed fixes are truly matched by the code.
✱ Grep "astrate/v1" in . · 21 matches
All `/astrate/v1/*` endpoints in the code are covered by the yaml; the previously-filed native-sync tasks are all landed and matched. The `ASTRATE_HOUSEKEEPING_*` config-drift I found was already escalated in `.mule/for-giulio.md` (line 71, 2026-09-12, and again line 153), so re-reporting it would duplicate the queue.

Done: Surface astrate-native (docs/api/astrate_native_api.yaml) is an accurate spec — no drift found, no proposals.

- Documented paths vs code: all 12 match — `/astrate/v1/{health,readiness,metrics}` (internal/observability/health.go:44-46), compat health×3 and version×4 (cmd/astrate/main.go:383-395), `/astrate/v1/{realm}/socket` (internal/appengine/stream/ws.go:46), `/appengine/v1/socket/websocket` (internal/appengine/channels/ws.go:48). No present-but-undocumented routes on this surface.
- Checks (past the required 3-4): readiness 200/503 + `{status,checks}` (health.go:56-78); native socket 101/200/401/403/500 vs guards (middleware.go:62-90); phoenix socket 101/401/500 (ws.go:56-75); compat health 200/503 `{data:{status:ok|unhealthy}}` (compat.go:20-33); metrics example `astrate_broker_sessions 0` is a real registered series (metrics.go:54-62); event wire shape `{event,realm,device_id,interface,path,value,timestamp}` matches wireEvent (ws.go:50-58).
- The three prior native-sync fixes are all landed and verified against the current yaml (build-version values, real metric series, socket 403/500).
- No docs build required: nothing touches `docs/api/` or `docs/mkdocs.yml`.

Files: none changed
Verified: no commands needed (static comparison via grep/read); no `docs/` touched so the build rule does not apply
Unsure: nothing
Follow-ups: none

Done: docs-sync astrate-native found no drift — spec is accurate, previous fixes verified landed, nothing to propose or escalate.
