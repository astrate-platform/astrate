slug: recipe-docs-sync
verdict: proposed
at:  5ba4eb0
ran: 2026-09-23T11:28:16Z on DietPi in 238s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
This is the docs-sync recipe — a proposal job. Let me pick one surface and find the drift.

$ ls docs/api/*.yaml
docs/api/astarte_appengine_api.yaml
docs/api/astarte_housekeeping_api.yaml
docs/api/astarte_pairing_api.yaml
docs/api/astarte_realm_management_api.yaml
docs/api/astrate_native_api.yaml

$ ls docs/api/ && ls .mule/ && ls internal/
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

→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ ls internal/appengine internal/flowapi internal/housekeeping internal/pairing internal/realm 2>/dev/null
internal/appengine:
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

internal/flowapi:
composite_wire_test.go
http.go
http_test.go
service.go
service_integration_test.go
service_test.go
userblocks.go
userblocks_test.go

internal/housekeeping:
async_operation_test.go
http.go
http_test.go
service.go
service_test.go

internal/pairing:
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

internal/realm:
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

$ rg -n 'HandleFunc|\.Get\(|\.Post\(|\.Put\(|\.Delete\(|\.Patch\(|\.Mount\(' internal/flowapi/http.go internal/api* internal/observability/ 2>/dev/null
internal/observability/compat.go:20:	mux.HandleFunc("GET /"+service+"/health", func(w http.ResponseWriter, r *http.Request) {
internal/observability/compat.go:51:	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
internal/observability/observability_test.go:48:	h.Mount(mux)
internal/observability/observability_test.go:56:	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
internal/observability/observability_test.go:69:	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
internal/observability/observability_test.go:89:	h.Mount(mux)
internal/observability/observability_test.go:112:	h.Mount(mux)
internal/observability/observability_test.go:119:	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
internal/observability/observability_test.go:141:			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
internal/observability/observability_test.go:164:			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
internal/observability/health.go:44:	mux.HandleFunc("GET /astrate/v1/health", h.handleHealth)
internal/observability/health.go:45:	mux.HandleFunc("GET /astrate/v1/readiness", h.handleReadiness)
internal/observability/compat_version_test.go:19:		resp, err := http.Get(srv.URL + "/" + svc + "/version")
internal/observability/compat_version_test.go:31:		if ct := resp.Header.Get("Content-Type"); ct == "" {
internal/observability/compat_version_test.go:39:	mux.HandleFunc("GET /v1/{realm}/version", VersionHandler("x"))
internal/observability/compat_version_test.go:43:	resp, err := http.Get(srv.URL + "/v1/somerealm/version")

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

Now let me survey the surfaces. I'll start by examining the registry of routes in code vs the docs, per the recipe's cheap-drift commands.
$ wc -l docs/api/astrate_native_api.yaml docs/api/astarte_*.yaml
466 docs/api/astrate_native_api.yaml
  1766 docs/api/astarte_appengine_api.yaml
   440 docs/api/astarte_housekeeping_api.yaml
   638 docs/api/astarte_pairing_api.yaml
  1171 docs/api/astarte_realm_management_api.yaml
  4481 total

→ Read docs/api/astrate_native_api.yaml
→ Read internal/observability/compat.go
→ Read internal/observability/health.go

$ rg -n 'MountServiceCompat|MountVersionCompat|Mount\(' cmd/astrate/main.go
359:	}).Mount(mux)
368:	housekeeping.NewAPI(hkSvc, mw, hkKeys).Mount(mux)
372:	realm.NewAPI(realmSvc, mw).Mount(mux)
373:	appengine.NewAPI(appengine.NewService(st, e, log), mw).Mount(mux)
374:	apstream.NewAPI(e.Bus(), mw).Mount(mux)
377:	channels.NewAPI(e.Bus(), st, st).Mount(mux)
379:	flowapi.NewAPI(flowSvc, mw).Mount(mux)
384:		observability.MountServiceCompat(mux, svc, st.Health)
394:		observability.MountVersionCompat(mux, svc, version)
411:	health.Mount(mux)
414:	swagger.Mount(mux)

$ rg -n 'socket|websocket|Socket|channels|stream' cmd/astrate/main.go internal/appengine/http.go
cmd/astrate/main.go:30:	"github.com/astrate-platform/astrate/internal/appengine/channels"
cmd/astrate/main.go:31:	apstream "github.com/astrate-platform/astrate/internal/appengine/stream"
cmd/astrate/main.go:154:	// the same device events as the stream socket (v2.0 process wiring).
cmd/astrate/main.go:190:	// Realm datastream retention ceilings (#72): sweep hourly — and once
cmd/astrate/main.go:192:	// (upstream clamps at write time; Astrate stores no per-row TTL).
cmd/astrate/main.go:366:		WithDefaultDatastreamMaximumStorageRetention(cfg.Housekeeping.DefaultDatastreamMaximumStorageRetention).
cmd/astrate/main.go:374:	apstream.NewAPI(e.Bus(), mw).Mount(mux)
cmd/astrate/main.go:375:	// Phoenix Channels socket (phoenix.js V2), alongside the Astrate-native
cmd/astrate/main.go:376:	// socket above: two protocols, one bus.
cmd/astrate/main.go:377:	channels.NewAPI(e.Bus(), st, st).Mount(mux)
cmd/astrate/main.go:381:	// Upstream-parity per-service health endpoints (the dashboard's API
cmd/astrate/main.go:387:	// Upstream-parity per-service version endpoints (issue #77): every
cmd/astrate/main.go:390:	// Pairing's is public — both measured on upstream 1.2.0
cmd/astrate/main.go:391:	// (test/conformance/upstream/verify-versions.json). Realm Management's
cmd/astrate/main.go:427:// upstream 1.2.0 (test/conformance/upstream/verify-versions.json): AppEngine's
cmd/astrate/main.go:429:// upstream API level, the value the Astarte Dashboard feature-gates its UI on
internal/appengine/http.go:21:// maxValueBytes is upstream's per-value limit (~64 KiB): a larger REST
internal/appengine/http.go:38:// Mount registers the routes on mux (paths wire-frozen to upstream
internal/appengine/http.go:56:	// Mirror surfaces (upstream registers these alongside the device-scoped
internal/appengine/http.go:102:// deviceListLinks builds the upstream pagination links object. links.next
internal/appengine/http.go:205:	// Upstream honors PATCH only with exactly this content type (a header
internal/appengine/http.go:228:	// fallback maps the mismatch to 400 "Bad request" upstream (probed live),
internal/appengine/http.go:251:// Parse failures render as 422 changeset field errors, matching upstream
internal/appengine/http.go:391:// key is distinguishable from an empty array (upstream rejects both, with
internal/appengine/http.go:454:// groupListLinks builds the upstream pagination links for a group-device
internal/appengine/http.go:457:// matching the probed upstream bodies — url.Values.Encode sorts the keys, so
internal/appengine/http.go:519:	// Same header comparison as patchDevice; upstream's group-scoped fallback
internal/appengine/http.go:551:// parseQueryOpts reads the datastream query parameters. Malformed values come
internal/appengine/http.go:552:// back as FieldErrors so writeError renders upstream's changeset-shaped 422
internal/appengine/http.go:634:	// Upstream default ordering for datastreams is descending (newest first).
internal/appengine/http.go:642:// writeError maps service/store errors onto upstream-shaped responses.

→ Read cmd/astrate/main.go [offset=355, limit=70]

$ rg -n 'HandleFunc|mux\.|Mount|GET |socket' internal/appengine/stream/*.go internal/appengine/channels/*.go
internal/appengine/channels/room.go:312:// Member is a joined socket with its own mailbox.
internal/appengine/stream/ws_test.go:60:	NewAPI(bus, auth.NewMiddleware(fakeKeys{realm})).Mount(mux)
internal/appengine/stream/ws_test.go:66:		req := httptest.NewRequest(http.MethodGet, "/astrate/v1/testrealm/socket", nil)
internal/appengine/stream/ws_test.go:68:		mux.ServeHTTP(rec, req)
internal/appengine/stream/ws_test.go:76:		req := httptest.NewRequest(http.MethodGet, "/astrate/v1/testrealm/socket?device_id=dev1&interface=com.ex.S", nil)
internal/appengine/stream/ws_test.go:79:		mux.ServeHTTP(rec, req)
internal/appengine/channels/protocol.go:2:// served at /appengine/v1/socket/websocket. It is the compatibility counterpart
internal/appengine/channels/protocol.go:3:// to the Astrate-native socket in internal/appengine/stream.
internal/appengine/channels/ws.go:12:	"github.com/coder/websocket"
internal/appengine/channels/ws.go:46:// Mount registers the WebSocket handler on the given mux.
internal/appengine/channels/ws.go:47:func (a *API) Mount(mux *http.ServeMux) {
internal/appengine/channels/ws.go:48:	mux.Handle("GET /appengine/v1/socket/websocket", http.HandlerFunc(a.handle))
internal/appengine/channels/ws.go:77:	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
internal/appengine/channels/ws.go:95:	conn    *websocket.Conn
internal/appengine/channels/ws.go:216:		if typ != websocket.MessageText {
internal/appengine/channels/ws.go:285:		// A rejoin on the same socket carries a fresh join_ref; server pushes
internal/appengine/channels/ws.go:441:	_ = s.conn.Write(ctx, websocket.MessageText, b)
internal/appengine/channels/ws_test.go:17:	"github.com/coder/websocket"
internal/appengine/channels/ws_test.go:61:	api.Mount(mux)
internal/appengine/channels/ws_test.go:63:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?realm=testrealm", nil)
internal/appengine/channels/ws_test.go:64:	mux.ServeHTTP(rec, req)
internal/appengine/channels/ws_test.go:74:	api.Mount(mux)
internal/appengine/channels/ws_test.go:76:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?token="+token, nil)
internal/appengine/channels/ws_test.go:77:	mux.ServeHTTP(rec, req)
internal/appengine/channels/ws_test.go:87:	api.Mount(mux)
internal/appengine/channels/ws_test.go:89:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?realm=nosuchrealm&token="+token, nil)
internal/appengine/channels/ws_test.go:90:	mux.ServeHTTP(rec, req)
internal/appengine/channels/ws_test.go:99:	api.Mount(mux)
internal/appengine/channels/ws_test.go:101:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?realm=testrealm&token=not.a.jwt", nil)
internal/appengine/channels/ws_test.go:102:	mux.ServeHTTP(rec, req)
internal/appengine/channels/ws_test.go:112:	api.Mount(mux)
internal/appengine/channels/ws_test.go:114:	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?realm=testrealm&token="+token, nil)
internal/appengine/channels/ws_test.go:115:	mux.ServeHTTP(rec, req)
internal/appengine/channels/ws_test.go:121:func dialSession(t *testing.T, mux *http.ServeMux, token string) (*websocket.Conn, context.CancelFunc) {
internal/appengine/channels/ws_test.go:128:		"/appengine/v1/socket/websocket?vsn=2.0.0&realm=" + testRealmName + "&token=" + token
internal/appengine/channels/ws_test.go:129:	conn, _, err := websocket.Dial(ctx, url, nil)
internal/appengine/channels/ws_test.go:141:	api.Mount(mux)
internal/appengine/channels/ws_test.go:147:	if err := conn.Write(ctx, websocket.MessageText, []byte(`[null,"1","phoenix","heartbeat",{}]`)); err != nil {
internal/appengine/channels/ws_test.go:194:	api.Mount(mux)
internal/appengine/channels/ws_test.go:201:	if err := conn.Write(ctx, websocket.MessageText, []byte(`["1","1","`+topic+`","phx_join",{}]`)); err != nil {
internal/appengine/channels/ws_test.go:234:	api.Mount(mux)
internal/appengine/channels/ws_test.go:255:			if err := conn.Write(ctx, websocket.MessageText, []byte(msg)); err != nil {
internal/appengine/channels/ws_test.go:290:	api.Mount(mux)
internal/appengine/channels/ws_test.go:299:	if err := conn.Write(ctx, websocket.MessageText, []byte(`["1","1","`+topic+`","phx_join",{}]`)); err != nil {
internal/appengine/channels/ws_test.go:310:	if err := conn.Write(ctx, websocket.MessageText, []byte(`["2","2","`+topic+`","phx_leave",{}]`)); err != nil {
internal/appengine/channels/ws_test.go:339:	api.Mount(mux)
internal/appengine/channels/ws_test.go:346:	if err := conn.Write(ctx, websocket.MessageText, []byte(`["1","1","`+topic+`","phx_leave",{}]`)); err != nil {
internal/appengine/channels/ws_test.go:379:func joinTopic(t *testing.T, conn *websocket.Conn, topic, ref string) {
internal/appengine/channels/ws_test.go:382:	if err := conn.Write(context.Background(), websocket.MessageText, []byte(msg)); err != nil {
internal/appengine/channels/ws_test.go:387:func readFramePayloadStatus(t *testing.T, conn *websocket.Conn) string {
internal/appengine/channels/ws_test.go:410:	api.Mount(mux)
internal/appengine/channels/ws_test.go:424:	if err := conn.Write(ctx, websocket.MessageText, []byte(msg)); err != nil {
internal/appengine/channels/ws_test.go:444:	api.Mount(mux)
internal/appengine/channels/ws_test.go:501:			if err := conn.Write(ctx, websocket.MessageText, []byte(msg)); err != nil {
internal/appengine/channels/ws_test.go:530:	if err := conn.Write(ctx, websocket.MessageText, []byte(msg)); err != nil {
internal/appengine/channels/ws_test.go:546:	api.Mount(mux)
internal/appengine/channels/ws_test.go:564:	if err := conn.Write(ctx, websocket.MessageText, []byte(watchMsg)); err != nil {
internal/appengine/channels/ws_test.go:637:	api.Mount(mux)
internal/appengine/channels/ws_test.go:663:	if err := conn.Write(ctx, websocket.MessageText, []byte(watchMsg)); err != nil {
internal/appengine/channels/ws_test.go:710:	api.Mount(mux)
internal/appengine/channels/ws_test.go:718:	if err := conn.Write(ctx, websocket.MessageText, []byte(`["1","1","`+topic+`","phx_join",{}]`)); err != nil {
internal/appengine/channels/ws_test.go:727:	if err := conn.Write(ctx, websocket.MessageText, []byte(watch)); err != nil {
internal/appengine/channels/ws_test.go:734:	// Drain the socket so the pump is never the thing that blocks.
internal/appengine/channels/ws_test.go:761:	if err := conn.Write(ctx, websocket.MessageText, []byte(`["1","3","`+topic+`","phx_leave",{}]`)); err != nil {
internal/appengine/channels/ws_test.go:769:// The unit table in internal/auth proves the matcher; this proves the socket
internal/appengine/channels/ws_test.go:772:// The refusals are paired with acceptances on the same room, so a socket that
internal/appengine/channels/ws_test.go:802:			api.Mount(mux)
internal/appengine/channels/ws_test.go:809:			if err := conn.Write(ctx, websocket.MessageText,
internal/appengine/channels/ws_test.go:975:			api.Mount(mux)
internal/appengine/channels/ws_test.go:990:			if err := conn.Write(ctx, websocket.MessageText,
internal/appengine/stream/ws.go:1:// Package stream is the AppEngine live event socket (docs/DESIGN.md §3.7,
internal/appengine/stream/ws.go:3:// at /astrate/v1/{realm}/socket fed by the engine's in-process fan-out bus.
internal/appengine/stream/ws.go:18:	"github.com/coder/websocket"
internal/appengine/stream/ws.go:24:// subscribeBuffer is the per-connection bus buffer; a slow socket drops
internal/appengine/stream/ws.go:28:// Bus is the subset of the engine fan-out bus the socket consumes.
internal/appengine/stream/ws.go:33:// API is the live-socket HTTP surface, guarded by an a_ch realm JWT.
internal/appengine/stream/ws.go:39:// NewAPI wires the socket to the engine bus and the a_ch middleware.
internal/appengine/stream/ws.go:44:// Mount registers the socket route on mux.
internal/appengine/stream/ws.go:45:func (a *API) Mount(mux *http.ServeMux) {
internal/appengine/stream/ws.go:46:	mux.Handle("GET /astrate/v1/{realm}/socket", a.require(http.HandlerFunc(a.handle)))
internal/appengine/stream/ws.go:87:	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
internal/appengine/stream/ws.go:93:	// the client disconnects (the socket is server-push only).
internal/appengine/stream/ws.go:102:				_ = conn.Close(websocket.StatusGoingAway, "bus closed")
internal/appengine/stream/ws.go:109:			if err := conn.Write(ctx, websocket.MessageText, payload); err != nil {

→ Read internal/appengine/stream/ws.go
→ Read internal/observability/metrics.go
→ Read internal/appengine/channels/ws.go [limit=95]

$ rg -n 'func .*RequireRealm|WriteUnauthorized|StatusUnauthorized|StatusForbidden|StatusInternalServerError|StatusBadRequest' internal/auth/*.go | head -40
internal/auth/middleware.go:48:func (m *Middleware) RequireRealm(claim Claim) func(http.Handler) http.Handler {
internal/auth/middleware.go:57:func (m *Middleware) RequireRealmAny(claims ...Claim) func(http.Handler) http.Handler {
internal/auth/middleware.go:62:				_ = astarteapi.WriteUnauthorized(w)
internal/auth/middleware.go:71:				_ = astarteapi.WriteUnauthorized(w)
internal/auth/middleware.go:133:		_ = astarteapi.WriteUnauthorized(w)
internal/auth/middleware.go:139:		_ = astarteapi.WriteUnauthorized(w)
internal/auth/middleware_test.go:85:		if w.Code != http.StatusUnauthorized {
internal/auth/middleware_test.go:93:		if w.Code != http.StatusUnauthorized {
internal/auth/middleware_test.go:101:		if w.Code != http.StatusUnauthorized {
internal/auth/middleware_test.go:110:		if w.Code != http.StatusForbidden {
internal/auth/middleware_test.go:127:		if w.Code != http.StatusForbidden {
internal/auth/middleware_test.go:140:		if w.Code != http.StatusForbidden {
internal/auth/middleware_test.go:168:		if w.Code != http.StatusUnauthorized {
internal/auth/middleware_test.go:192:		if w.Code != http.StatusUnauthorized {
internal/auth/middleware_test.go:226:func TestRequireRealmAnyFlowClaims(t *testing.T) {
internal/auth/middleware_test.go:270:		if w := do(mux, "GET", "/flow/v1/test/pipelines", narrow); w.Code != http.StatusForbidden {
internal/auth/middleware_test.go:275:		if w := do(mux, "GET", "/flow/v1/test/pipelines", otherTok); w.Code != http.StatusForbidden {
internal/auth/middleware_test.go:280:		if w := do(mux, "GET", "/flow/v1/test/pipelines", "not.a.token"); w.Code != http.StatusUnauthorized {

→ Read internal/auth/middleware.go [limit=120, offset=40]

$ rg -n 'func Write' pkg/astarteapi/*.go
pkg/astarteapi/envelope.go:126:func WriteData(w http.ResponseWriter, status int, v any) error {
pkg/astarteapi/envelope.go:145:func WriteDataWithLinks(w http.ResponseWriter, status int, v any, links Links) error {
pkg/astarteapi/envelope.go:156:func WriteDataWithMetadata(w http.ResponseWriter, status int, v any, m map[string]any) error {
pkg/astarteapi/envelope.go:163:func WriteError(w http.ResponseWriter, status int, detail string) error {
pkg/astarteapi/envelope.go:172:func WriteFieldErrors(w http.ResponseWriter, status int, fields map[string][]string) error {
pkg/astarteapi/envelope.go:186:func WriteRawErrors(w http.ResponseWriter, status int, errs any) error {
pkg/astarteapi/envelope.go:191:func WriteBadRequest(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:196:func WriteUnauthorized(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:201:func WriteForbidden(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:206:func WriteNotFound(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:211:func WriteDeviceNotFound(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:216:func WriteInternalServerError(w http.ResponseWriter) error {

→ Read pkg/astarteapi/envelope.go [limit=40, offset=185]

$ rg -n 'DetailUnauthorized|DetailForbidden|DetailInternalServerError' pkg/astarteapi/*.go
pkg/astarteapi/envelope.go:38:	// DetailUnauthorized is the canonical 401 detail.
pkg/astarteapi/envelope.go:39:	DetailUnauthorized = "Unauthorized"
pkg/astarteapi/envelope.go:40:	// DetailForbidden is the canonical 403 detail.
pkg/astarteapi/envelope.go:41:	DetailForbidden = "Forbidden"
pkg/astarteapi/envelope.go:55:	// DetailInternalServerError is the canonical 500 detail. Measured
pkg/astarteapi/envelope.go:58:	DetailInternalServerError = "Internal server error"
pkg/astarteapi/envelope.go:110:		body, _ = marshal(detailEnvelope{Errors: detailBody{Detail: DetailInternalServerError}})
pkg/astarteapi/envelope.go:197:	return WriteError(w, http.StatusUnauthorized, DetailUnauthorized)
pkg/astarteapi/envelope.go:202:	return WriteError(w, http.StatusForbidden, DetailForbidden)
pkg/astarteapi/envelope.go:217:	return WriteError(w, http.StatusInternalServerError, DetailInternalServerError)

$ ls docs/site/ && rg -n -i 'flow' docs/site/mkdocs.yml 2>/dev/null; rg -n 'APIYAML|yaml' internal/swagger/*.go | head
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
internal/swagger/swagger.go:18:	apiRoot, _ := fs.Sub(docs.APIYAML, "api")
internal/swagger/swagger.go:25:	// Serve YAML files at /api/ so the relative ../api/*.yaml in index.html
internal/swagger/swagger.go:33:	_ = fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
internal/swagger/swagger.go:37:		if strings.HasSuffix(path, ".yaml") {
internal/swagger/swagger_test.go:55:			want, err := docs.APIYAML.ReadFile("api/" + name)
internal/swagger/swagger_test.go:77:		if !strings.HasSuffix(name, ".yaml") {
internal/swagger/swagger_test.go:78:			t.Errorf("Specs() entry %q does not end in .yaml", name)
internal/swagger/swagger_test.go:80:		if _, err := docs.APIYAML.ReadFile("api/" + name); err != nil {
internal/swagger/swagger_test.go:81:			t.Errorf("Specs() entry %q cannot be read from docs.APIYAML: %v", name, err)
internal/swagger/swagger_test.go:96:	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {

$ ls -la docs/site/api/ && rg -n -i 'flow api|/flow/v1|flow operator' docs/site/*.md docs/mkdocs.yml 2>/dev/null | head -20
total 156
drwxr-xr-x 2 root root  4096 Sep 20 21:25 .
drwxr-xr-x 4 root root  4096 Sep 20 21:25 ..
-rw-r--r-- 1 root root 56583 Sep 20 21:25 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 13640 Sep 20 21:25 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 19506 Sep 20 21:25 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 38625 Sep 20 21:25 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 14977 Sep 20 21:25 astrate_native_api.yaml
docs/site/DESIGN.md:64:| **Astarte Flow** | Dataflow processing framework | `internal/flow`, `internal/flowapi`, `internal/flow/blocks` | **v2.0 milestone.** Operator surface `/flow/v1/{realm}/…` (pipelines, flows, blocks). Built-in catalog + planned durable named flows and container blocks. Capability parity, not an Elixir port. |

→ Read pkg/astarteapi/envelope.go [limit=70, offset=60]

$ rg -n 'ContentType' pkg/astarteapi/envelope.go
21:// ContentType is the Content-Type header value set on every envelope,
23:const ContentType = "application/json; charset=utf-8"
112:		w.Header().Set("Content-Type", ContentType)
117:	w.Header().Set("Content-Type", ContentType)

! permission requested: external_directory (/tmp/*); auto-rejecting
Checking the metrics endpoint's exact Content-Type via a quick probe:
✗ Write /tmp/probe_metrics.go failed
Error: The user rejected permission to use this specific tool call.
