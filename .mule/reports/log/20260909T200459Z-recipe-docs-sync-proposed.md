slug: recipe-docs-sync
verdict: proposed
at:  304c126
ran: 2026-09-09T20:04:59Z on DietPi in 282s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ ls /root/astrate-mule/docs/api/
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' /root/astrate-mule/internal/ --glob '*native*'
(no output)

I'll pick the `astrate-native` surface — I don't see any prior docs-sync tasks for it. Let me find the routes and documented paths.
$ rg -n '^\s{2}/' /root/astrate-mule/docs/api/astrate_native_api.yaml
26:  /astrate/v1/health:
44:  /astrate/v1/readiness:
77:  /astrate/v1/metrics:
99:  /appengine/health:
118:  /realmmanagement/health:
137:  /pairing/health:
158:  /astrate/v1/{realm}/socket:
220:  /appengine/v1/socket/websocket:

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' /root/astrate-mule/internal/ --glob '*native*' --glob '*astrate_native*'
(no output)

$ ls /root/astrate-mule/internal/ | grep -i nat
(no output)

✱ Glob "internal/**/native*" 0 matches

$ rg -rn 'astrate/v1|/astrate' /root/astrate-mule/internal/ --glob '*.go' | head -30
/root/astrate-mule/internal/flow/factory_test.go:	"github.comn-platformn/internal/engine/stream"
/root/astrate-mule/internal/flow/factory_test.go:	"github.comn-platformn/internal/flow"
/root/astrate-mule/internal/flow/factory_test.go:	"github.comn-platformn/internal/flow/blocks"
/root/astrate-mule/internal/flow/factory_test.go:	"github.comn-platformn/internal/flow/blocks/astartesource"
/root/astrate-mule/internal/config/config.go:// cmdn/newnats_*.go).
/root/astrate-mule/internal/flow/blocks/mqtt.go:	"github.comn-platformn/internal/flow"
/root/astrate-mule/internal/flow/blocks/virtualdevicepool/virtualdevicepool_test.go:	"github.comn-platformn/internal/flow"
/root/astrate-mule/internal/flow/blocks/virtualdevicepool/virtualdevicepool_test.go:	"github.comn-platformn/internal/flow/blocks/virtualdevicepool"
/root/astrate-mule/internal/flow/blocks/virtualdevicepool/virtualdevicepool.go:	"github.comn-platformn/internal/flow"
/root/astrate-mule/internal/flow/blocks/transform_test.go:	"github.comn-platformn/internal/engine/stream"
/root/astrate-mule/internal/flow/blocks/transform_test.go:	"github.comn-platformn/internal/flow"
/root/astrate-mule/internal/flow/blocks/transform_test.go:	"github.comn-platformn/internal/flow/blocks"
/root/astrate-mule/internal/engine/triggers_change_test.go:	"github.comn-platformn/pkg/payload"
/root/astrate-mule/internal/engine/router_test.go:	"github.comn-platformn/internal/broker"
/root/astrate-mule/internal/engine/router_test.go:	"github.comn-platformn/internal/store"
/root/astrate-mule/internal/engine/router_test.go:	"github.comn-platformn/pkg/deviceid"
/root/astrate-mule/internal/engine/router_test.go:	"github.comn-platformn/pkg/payload"
/root/astrate-mule/internal/flow/blocks/httpblocks.go:	"github.comn-platformn/internal/flow"
/root/astrate-mule/internal/flow/blocks/jsonpathmap_test.go:	"github.comn-platformn/internal/flow"
/root/astrate-mule/internal/flow/blocks/jsonpathmap_test.go:	"github.comn-platformn/internal/flow/blocks"
/root/astrate-mule/internal/engine/data.go:	"github.comn-platformn/internal/broker"
/root/astrate-mule/internal/engine/data.go:	"github.comn-platformn/internal/engine/triggers"
/root/astrate-mule/internal/engine/data.go:	"github.comn-platformn/internal/store"
/root/astrate-mule/internal/engine/data.go:	"github.comn-platformn/pkg/deviceid"
/root/astrate-mule/internal/engine/data.go:	"github.comn-platformn/pkg/interfaceschema"
/root/astrate-mule/internal/engine/data.go:	"github.comn-platformn/pkg/payload"
/root/astrate-mule/internal/flow/blocks/sort.go:	"github.comn-platformn/internal/flow"
/root/astrate-mule/internal/flow/blocks/updatemetadata.go:	"github.comn-platformn/internal/flow"
/root/astrate-mule/internal/flow/blocks/mqtt_test.go:	"github.comn-platformn/internal/flow"
/root/astrate-mule/internal/flow/blocks/mqtt_test.go:	"github.comn-platformn/internal/flow/blocks"

$ rg -n 'HandleFunc|\.Get\(|\.Post\(' /root/astrate-mule/internal/httpx/ /root/astrate-mule/internal/observability/ 2>/dev/null | head -20
/root/astrate-mule/internal/observability/compat.go:20:	mux.HandleFunc("GET /"+service+"/health", func(w http.ResponseWriter, r *http.Request) {
/root/astrate-mule/internal/observability/compat.go:51:	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
/root/astrate-mule/internal/httpx/cors_test.go:44:		if got := rec.Header().Get(header); got != want {
/root/astrate-mule/internal/httpx/cors_test.go:57:	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:4040" {
/root/astrate-mule/internal/httpx/cors_test.go:69:	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
/root/astrate-mule/internal/httpx/cors_test.go:75:	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
/root/astrate-mule/internal/httpx/cors_test.go:86:	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
/root/astrate-mule/internal/httpx/cors_test.go:89:	if got := rec.Header().Get("Vary"); got != "" {
/root/astrate-mule/internal/httpx/cors_test.go:100:	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
/root/astrate-mule/internal/httpx/notfound_test.go:13:	mux.HandleFunc("GET /appengine/v1/{realm}/devices", func(w http.ResponseWriter, _ *http.Request) {
/root/astrate-mule/internal/httpx/notfound_test.go:19:	mux.HandleFunc("GET /housekeeping/v1/realms", func(w http.ResponseWriter, _ *http.Request) {
/root/astrate-mule/internal/httpx/notfound_test.go:106:				ct := rr.Header().Get("Content-Type")
/root/astrate-mule/internal/httpx/notfound_test.go:119:	mux.HandleFunc("GET /appengine/v1/{path...}", func(w http.ResponseWriter, _ *http.Request) {
/root/astrate-mule/internal/httpx/notfound_test.go:136:	if rr.Header().Get("X-Custom") != "yes" {
/root/astrate-mule/internal/observability/health.go:43:	mux.HandleFunc("GET /astrate/v1/health", h.handleHealth)
/root/astrate-mule/internal/observability/health.go:44:	mux.HandleFunc("GET /astrate/v1/readiness", h.handleReadiness)
/root/astrate-mule/internal/httpx/cors.go:25:			origin := r.Header.Get("Origin")
/root/astrate-mule/internal/httpx/cors.go:39:			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
/root/astrate-mule/internal/observability/compat_version_test.go:19:		resp, err := http.Get(srv.URL + "/" + svc + "/version")
/root/astrate-mule/internal/observability/compat_version_test.go:31:		if ct := resp.Header.Get("Content-Type"); ct == "" {

$ rg -rn 'astrate/v1|health|readiness|metrics|/socket|/websocket' /root/astrate-mule/internal/observability/ --glob '*.go' | head -30
/root/astrate-mule/internal/observability/compat.go:// GET /{service}/n endpoint. The Astarte Dashboard polls one per service
/root/astrate-mule/internal/observability/compat.go:// same dependency probe that backs /n/n — and answers 503
/root/astrate-mule/internal/observability/compat.go:	mux.HandleFunc("GET /"+service+"/n", func(w http.ResponseWriter, r *http.Request) {
/root/astrate-mule/internal/observability/compat.go:			ctx, cancel := context.WithTimeout(r.Context(), nTimeout)
/root/astrate-mule/internal/observability/compat.go:				_, _ = w.Write([]byte(`{"data":{"status":"unny"}}`))
/root/astrate-mule/internal/observability/health.go:// nTimeout bounds the whole n probe so a wedged dependency
/root/astrate-mule/internal/observability/health.go:const nTimeout = 3 * time.Second
/root/astrate-mule/internal/observability/health.go:// Check reports a dependency's n; a nil error means ready. The store
/root/astrate-mule/internal/observability/health.go:// DB ping and a broker-listener check are the two §5.2 n backends.
/root/astrate-mule/internal/observability/health.go:// liveness, n, and the Prometheus scrape endpoint, all under
/root/astrate-mule/internal/observability/health.go:// /n so they never collide with the upstream API namespace.
/root/astrate-mule/internal/observability/health.go:	n http.Handler
/root/astrate-mule/internal/observability/health.go:// NewHealth builds the surface around an existing n handler.
/root/astrate-mule/internal/observability/health.go:func NewHealth(n http.Handler) *Health {
/root/astrate-mule/internal/observability/health.go:	return &Health{n: n}
/root/astrate-mule/internal/observability/health.go:// AddReadiness registers a named dependency probed by /n.
/root/astrate-mule/internal/observability/health.go:	mux.HandleFunc("GET /n/n", h.handleHealth)
/root/astrate-mule/internal/observability/health.go:	mux.HandleFunc("GET /n/n", h.handleReadiness)
/root/astrate-mule/internal/observability/health.go:	mux.Handle("GET /n/n", h.n)
/root/astrate-mule/internal/observability/health.go:	ctx, cancel := context.WithTimeout(r.Context(), nTimeout)
/root/astrate-mule/internal/observability/metrics.go:// n/n/n HTTP surface under /n (docs/DESIGN.md
/root/astrate-mule/internal/observability/metrics.go:// nNamespace prefixes every Astrate-defined series (astrate_*).
/root/astrate-mule/internal/observability/metrics.go:const nNamespace = "astrate"
/root/astrate-mule/internal/observability/metrics.go:		Namespace: nNamespace,
/root/astrate-mule/internal/observability/metrics.go:			Namespace: nNamespace,
/root/astrate-mule/internal/observability/observability_test.go:	m.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/n/n", nil))
/root/astrate-mule/internal/observability/observability_test.go:	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/n/n", nil))
/root/astrate-mule/internal/observability/observability_test.go:		t.Errorf("n = %d, want 200", rec.Code)
/root/astrate-mule/internal/observability/observability_test.go:	// One failing dependency makes n 503.
/root/astrate-mule/internal/observability/observability_test.go:	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/n/n", nil))

$ rg -rn 'HandleFunc|\.Get\(|\.Post\(' /root/astrate-mule/internal/observability/ --glob '*.go' | grep -v _test
/root/astrate-mule/internal/observability/compat.go:	mux.n("GET /"+service+"/health", func(w http.ResponseWriter, r *http.Request) {
/root/astrate-mule/internal/observability/compat.go:	mux.n("GET /"+service+"/version", VersionHandler(version))
/root/astrate-mule/internal/observability/health.go:	mux.n("GET /astrate/v1/health", h.handleHealth)
/root/astrate-mule/internal/observability/health.go:	mux.n("GET /astrate/v1/readiness", h.handleReadiness)

→ Read internal/observability/health.go
→ Read internal/observability/metrics.go
→ Read internal/observability/compat.go
Now let me check the websocket routes and the compat.go routes.

$ rg -rn 'HandleFunc|\.Get\(|\.Post\(' /root/astrate-mule/internal/appengine/ --glob '*.go' | grep -v _test | grep -i 'socket\|websocket'
(no output)

$ rg -rn 'socket|websocket' /root/astrate-mule/internal/appengine/ --glob '*.go' | grep -v _test | head -20
/root/astrate-mule/internal/appengine/stream/ws.go:// Package stream is the AppEngine live event n (docs/DESIGN.md §3.7,
/root/astrate-mule/internal/appengine/stream/ws.go:// at /astrate/v1/{realm}/n fed by the engine's in-process fan-out bus.
/root/astrate-mule/internal/appengine/stream/ws.go:	"github.com/coder/n"
/root/astrate-mule/internal/appengine/stream/ws.go:// subscribeBuffer is the per-connection bus buffer; a slow n drops
/root/astrate-mule/internal/appengine/stream/ws.go:// Bus is the subset of the engine fan-out bus the n consumes.
/root/astrate-mule/internal/appengine/stream/ws.go:// API is the live-n HTTP surface, guarded by an a_ch realm JWT.
/root/astrate-mule/internal/appengine/stream/ws.go:// NewAPI wires the n to the engine bus and the a_ch middleware.
/root/astrate-mule/internal/appengine/stream/ws.go:// Mount registers the n route on mux.
/root/astrate-mule/internal/appengine/stream/ws.go:	mux.Handle("GET /astrate/v1/{realm}/n", a.require(http.HandlerFunc(a.handle)))
/root/astrate-mule/internal/appengine/stream/ws.go:	conn, err := n.Accept(w, r, &n.AcceptOptions{InsecureSkipVerify: true})
/root/astrate-mule/internal/appengine/stream/ws.go:	// the client disconnects (the n is server-push only).
/root/astrate-mule/internal/appengine/stream/ws.go:				_ = conn.Close(n.StatusGoingAway, "bus closed")
/root/astrate-mule/internal/appengine/stream/ws.go:			if err := conn.Write(ctx, n.MessageText, payload); err != nil {
/root/astrate-mule/internal/appengine/service.go:// server-owned publishing, groups, and the live event n. It is wire-shaped
/root/astrate-mule/internal/appengine/channels/ws.go:	"github.com/coder/n"
/root/astrate-mule/internal/appengine/channels/ws.go:	mux.Handle("GET /appengine/v1/n/n", http.HandlerFunc(a.handle))
/root/astrate-mule/internal/appengine/channels/ws.go:	conn, err := n.Accept(w, r, &n.AcceptOptions{InsecureSkipVerify: true})
/root/astrate-mule/internal/appengine/channels/ws.go:	conn    *n.Conn
/root/astrate-mule/internal/appengine/channels/ws.go:		if typ != n.MessageText {
/root/astrate-mule/internal/appengine/channels/ws.go:		// A rejoin on the same n carries a fresh join_ref; server pushes

→ Read docs/api/astrate_native_api.yaml
→ Read internal/appengine/stream/ws.go
Let me check the actual code routes against the spec. I see the spec documents `/astrate/v1/{realm}/socket` and `/appengine/v1/socket/websocket`, but the code registers different paths.
→ Read internal/appengine/channels/ws.go

$ rg -n 'MountVersionCompat|MountServiceCompat|Mount' /root/astrate-mule/cmd/ --glob '*.go' | head -30
/root/astrate-mule/cmd/astrate/main.go:359:	}).Mount(mux)
/root/astrate-mule/cmd/astrate/main.go:368:	housekeeping.NewAPI(hkSvc, mw, hkKeys).Mount(mux)
/root/astrate-mule/cmd/astrate/main.go:372:	realm.NewAPI(realmSvc, mw).Mount(mux)
/root/astrate-mule/cmd/astrate/main.go:373:	appengine.NewAPI(appengine.NewService(st, e, log), mw).Mount(mux)
/root/astrate-mule/cmd/astrate/main.go:374:	apstream.NewAPI(e.Bus(), mw).Mount(mux)
/root/astrate-mule/cmd/astrate/main.go:377:	channels.NewAPI(e.Bus(), st).Mount(mux)
/root/astrate-mule/cmd/astrate/main.go:379:	flowapi.NewAPI(flowSvc, mw).Mount(mux)
/root/astrate-mule/cmd/astrate/main.go:384:		observability.MountServiceCompat(mux, svc, st.Health)
/root/astrate-mule/cmd/astrate/main.go:394:		observability.MountVersionCompat(mux, svc, version)
/root/astrate-mule/cmd/astrate/main.go:413:	health.Mount(mux)
/root/astrate-mule/cmd/astrate/main.go:416:	swagger.Mount(mux)

$ rg -n 'version' /root/astrate-mule/docs/api/astrate_native_api.yaml
9:  version: "0.0.1"
83:        type is `text/plain; version=0.0.4`.
89:            text/plain; version=0.0.4:

→ Read cmd/astrate/main.go [offset=378, limit=40]

$ rg -n 'version' /root/astrate-mule/docs/api/astrate_housekeeping_api.yaml | head -5
rg: /root/astrate-mule/docs/api/astrate_housekeeping_api.yaml: IO error for operation on /root/astrate-mule/docs/api/astrate_housekeeping_api.yaml: No such file or directory (os error 2)

$ for f in /root/astrate-mule/docs/api/*.yaml; do echo "=== $(basename $f) ==="; rg -n 'version' "$f" | head -5; done
=== astarte_appengine_api.yaml ===
9:  version: "0.0.1"
188:                  firmware_version: "2.1.0"
272:                  firmware_version: "2.1.0"
1013:                  firmware_version: "2.1.0"
1467:          description: Map of interface name to its major and minor version.
=== astarte_housekeeping_api.yaml ===
11:  version: "0.0.1"
=== astarte_pairing_api.yaml ===
12:  version: "0.0.1"
205:        CA certificate), and version. The device authenticates with its
227:                  version: "0.1.0-test"
395:      required: [protocols, status, version]
414:        version:
=== astarte_realm_management_api.yaml ===
10:  version: "0.0.1"
73:                        version_major: 1
74:                        version_minor: 0
139:      summary: List major versions of an interface
140:      description: Returns the installed major version numbers for the named interface.
=== astrate_native_api.yaml ===
9:  version: "0.0.1"
83:        type is `text/plain; version=0.0.4`.
89:            text/plain; version=0.0.4:

$ for f in /root/astrate-mule/docs/api/*.yaml; do rg -n '/version' "$f" | head -3 && echo "--- $(basename $f) ---"; done
--- astarte_appengine_api.yaml ---
--- astarte_housekeeping_api.yaml ---
--- astarte_pairing_api.yaml ---
532:  /realmmanagement/v1/{realm}/version:
--- astarte_realm_management_api.yaml ---
--- astrate_native_api.yaml ---

$ rg -n 'version' /root/astrate-mule/docs/api/astarte_pairing_api.yaml
12:  version: "0.0.1"
205:        CA certificate), and version. The device authenticates with its
227:                  version: "0.1.0-test"
395:      required: [protocols, status, version]
414:        version:
416:          description: Astrate version string.

→ Read docs/api/astarte_pairing_api.yaml [offset=525, limit=20]

$ rg -n 'realmmanagement.*version|pairing.*version|/version' /root/astrate-mule/docs/api/astarte_pairing_api.yaml
(no output)

$ for f in /root/astrate-mule/docs/api/*.yaml; do echo "=== $(basename $f) ==="; rg -n '^\s{2}/' "$f"; done
=== astarte_appengine_api.yaml ===
26:  /appengine/v1/{realm}/devices:
98:  /appengine/v1/{realm}/stats/devices:
136:  /appengine/v1/{realm}/devices/{device}:
210:  /appengine/v1/{realm}/devices-by-alias/{alias}:
292:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces:
324:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
369:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
545:  /appengine/v1/{realm}/devices/{device}/interfaces:
573:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
613:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
769:  /appengine/v1/{realm}/groups:
841:  /appengine/v1/{realm}/groups/{group}:
869:  /appengine/v1/{realm}/groups/{group}/devices:
957:  /appengine/v1/{realm}/groups/{group}/devices/{device}:
1056:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces:
1084:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}:
1125:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}:
=== astarte_housekeeping_api.yaml ===
26:  /housekeeping/v1/realms:
106:  /housekeeping/v1/realms/{realm}:
=== astarte_pairing_api.yaml ===
27:  /pairing/v1/{realm}/health:
70:  /pairing/v1/{realm}/agent/devices:
122:  /pairing/v1/{realm}/agent/devices/{deviceID}:
145:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
199:  /pairing/v1/{realm}/devices/{deviceID}:
237:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:
=== astarte_realm_management_api.yaml ===
27:  /realmmanagement/v1/{realm}/interfaces:
136:  /realmmanagement/v1/{realm}/interfaces/{name}:
165:  /realmmanagement/v1/{realm}/interfaces/{name}/{major}:
277:  /realmmanagement/v1/{realm}/triggers:
347:  /realmmanagement/v1/{realm}/triggers/{name}:
399:  /realmmanagement/v1/{realm}/config/auth:
466:  /realmmanagement/v1/{realm}/config/device_registration_limit:
497:  /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention:
532:  /realmmanagement/v1/{realm}/version:
565:  /realmmanagement/v1/{realm}/devices/{device}:
591:  /realmmanagement/v1/{realm}/policies:
661:  /realmmanagement/v1/{realm}/policies/{name}:
=== astrate_native_api.yaml ===
26:  /astrate/v1/health:
44:  /astrate/v1/readiness:
77:  /astrate/v1/metrics:
99:  /appengine/health:
118:  /realmmanagement/health:
137:  /pairing/health:
158:  /astrate/v1/{realm}/socket:
220:  /appengine/v1/socket/websocket:

$ rg -n 'version' /root/astrate-mule/internal/realm/http.go /root/astrate-mule/internal/pairing/http.go /root/astrate-mule/internal/appengine/http.go | grep -i -v firmware
/root/astrate-mule/internal/realm/http.go:52:	mux.Handle("GET /realmmanagement/v1/{realm}/version", h(a.getVersion))
/root/astrate-mule/internal/realm/http.go:372:			"Interface major version doesn't match the one in the interface json")
/root/astrate-mule/internal/realm/http.go:375:			"Interface minor version was not increased")
/root/astrate-mule/internal/realm/http.go:416:			"Interface major version is not 0, can't be deleted")
/root/astrate-mule/internal/pairing/http.go:198:	Version string `json:"version"`

$ sed -n '520,565p' /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
type: integer
                    format: int64
                    description: Maximum datastream storage retention in seconds. 0 means unset.
              example:
                data: 31536000
        "401":
          $ref: "#/components/responses/Unauthorized"
        "500":
          $ref: "#/components/responses/InternalServerError"

  # ── Version ─────────────────────────────────────────────────────────

  /realmmanagement/v1/{realm}/version:
    get:
      operationId: getVersion
      summary: Get emulated API version
      description: |
        Returns the emulated upstream API compatibility version. The Astarte
        Dashboard gates feature UI on this value.
      tags: [Realm Management - Misc]
      security:
        - a_rma: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
      responses:
        "200":
          description: API compatibility version.
          content:
            application/json; charset=utf-8:
              schema:
                type: object
                required: [data]
                properties:
                  data:
                    type: string
                    description: Version string.
              example:
                data: "1.1.0"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "500":
          $ref: "#/components/responses/InternalServerError"

  # ── Device lifecycle ────────────────────────────────────────────────

  /realmmanagement/v1/{realm}/devices/{device}:

$ rg -n 'func.*RequireRealm|WriteUnauthorized|StatusUnauthorized' /root/astrate-mule/internal/auth/ --glob '*.go' | head
/root/astrate-mule/internal/auth/middleware.go:48:func (m *Middleware) RequireRealm(claim Claim) func(http.Handler) http.Handler {
/root/astrate-mule/internal/auth/middleware.go:57:func (m *Middleware) RequireRealmAny(claims ...Claim) func(http.Handler) http.Handler {
/root/astrate-mule/internal/auth/middleware.go:62:				_ = astarteapi.WriteUnauthorized(w)
/root/astrate-mule/internal/auth/middleware.go:71:				_ = astarteapi.WriteUnauthorized(w)
/root/astrate-mule/internal/auth/middleware.go:133:		_ = astarteapi.WriteUnauthorized(w)
/root/astrate-mule/internal/auth/middleware.go:139:		_ = astarteapi.WriteUnauthorized(w)
/root/astrate-mule/internal/auth/middleware_test.go:85:		if w.Code != http.StatusUnauthorized {
/root/astrate-mule/internal/auth/middleware_test.go:93:		if w.Code != http.StatusUnauthorized {
/root/astrate-mule/internal/auth/middleware_test.go:101:		if w.Code != http.StatusUnauthorized {
/root/astrate-mule/internal/auth/middleware_test.go:155:		if w.Code != http.StatusUnauthorized {

Now let me verify the auth middleware behavior for the socket endpoint and check what the compat version endpoints return exactly.
$ rg -rn 'MountVersionCompat|VersionHandler' /root/astrate-mule/internal/observability/ /root/astrate-mule/cmd/ --glob '*.go' | grep -v _test
/root/astrate-mule/internal/observability/compat.go:// n answers GET .../version with the upstream envelope
/root/astrate-mule/internal/observability/compat.go:func n(version string) http.HandlerFunc {
/root/astrate-mule/internal/observability/compat.go:// n registers GET /{service}/version for one service.
/root/astrate-mule/internal/observability/compat.go:func n(mux *http.ServeMux, service, version string) {
/root/astrate-mule/internal/observability/compat.go:	mux.HandleFunc("GET /"+service+"/version", n(version))
/root/astrate-mule/cmd/astrate/main.go:		observability.n(mux, svc, version)
/root/astrate-mule/cmd/astrate/main.go:		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.n(version)))
/root/astrate-mule/cmd/astrate/main.go:	mux.Handle("GET /pairing/v1/{realm}/version", observability.n(version))

$ rg -o '\bASTRATE_[A-Z_]+' -N /root/astrate-mule/internal/ /root/astrate-mule/cmd/ --glob '*.go' | sort -u
/root/astrate-mule/cmd/astrate/main.go:ASTRATE_HTTP_ADDR
/root/astrate-mule/cmd/astrate/main.go:ASTRATE_MASTER_KEY
/root/astrate-mule/cmd/astrate/main.go:ASTRATE_MASTER_KEY_FILE
/root/astrate-mule/cmd/astrate/main_test.go:ASTRATE_MASTER_KEY
/root/astrate-mule/cmd/astrate/main_test.go:ASTRATE_TEST_DSN
/root/astrate-mule/internal/config/config.go:ASTRATE_DATABASE_DSN
/root/astrate-mule/internal/config/config.go:ASTRATE_ENGINE_SHARDS
/root/astrate-mule/internal/config/config.go:ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
/root/astrate-mule/internal/config/config.go:ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
/root/astrate-mule/internal/config/config.go:ASTRATE_HTTP_ADDR
/root/astrate-mule/internal/config/config.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
/root/astrate-mule/internal/config/config.go:ASTRATE_HTTP_TLS_CERT_FILE
/root/astrate-mule/internal/config/config.go:ASTRATE_HTTP_TLS_KEY_FILE
/root/astrate-mule/internal/config/config.go:ASTRATE_LOG_FORMAT
/root/astrate-mule/internal/config/config.go:ASTRATE_LOG_LEVEL
/root/astrate-mule/internal/config/config.go:ASTRATE_MASTER_KEY
/root/astrate-mule/internal/config/config.go:ASTRATE_MASTER_KEY_FILE
/root/astrate-mule/internal/config/config.go:ASTRATE_MQTT_ADDR
/root/astrate-mule/internal/config/config.go:ASTRATE_MQTT_ADVERTISED_URL
/root/astrate-mule/internal/config/config.go:ASTRATE_MQTT_INSECURE_DEV_MODE
/root/astrate-mule/internal/config/config.go:ASTRATE_MQTT_SESSION_STORE_PATH
/root/astrate-mule/internal/config/config.go:ASTRATE_MQTT_TLS_CERT_FILE
/root/astrate-mule/internal/config/config.go:ASTRATE_MQTT_TLS_KEY_FILE
/root/astrate-mule/internal/config/config.go:ASTRATE_REALM_JWT_PUBLIC_KEY
/root/astrate-mule/internal/config/config.go:ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
/root/astrate-mule/internal/config/config.go:ASTRATE_REALM_NAME
/root/astrate-mule/internal/config/config.go:ASTRATE_SECURITY_MASTER_KEY_FILE
/root/astrate-mule/internal/config/config_test.go:ASTRATE_DATABASE_DSN
/root/astrate-mule/internal/config/config_test.go:ASTRATE_ENGINE_SHARDS
/root/astrate-mule/internal/config/config_test.go:ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
/root/astrate-mule/internal/config/config_test.go:ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
/root/astrate-mule/internal/config/config_test.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
/root/astrate-mule/internal/engine/engine_integration_test.go:ASTRATE_TEST_DSN
/root/astrate-mule/internal/flow/blocks/container/block.go:ASTRATE_FLOW_CONFIG
/root/astrate-mule/internal/flow/blocks/container/docker.go:ASTRATE_FLOW_CONFIG
/root/astrate-mule/internal/flow/blocks/info.go:ASTRATE_FLOW_CONFIG
/root/astrate-mule/internal/flow/blocks/schema.go:ASTRATE_FLOW_CONFIG
/root/astrate-mule/internal/store/crypto.go:ASTRATE_MASTER_KEY
/root/astrate-mule/internal/store/crypto.go:ASTRATE_MASTER_KEY_FILE
/root/astrate-mule/internal/store/store_test.go:ASTRATE_TEST_DSN
/root/astrate-mule/internal/testutil/pg.go:ASTRATE_TEST_DSN

Now let me run the config-key sweep the recipe also calls for.
$ rg -o '\bASTRATE_[A-Z_]+' -N /root/astrate-mule/docs/site/configuration-reference.md | sort -u
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

$ rg -n 'ASTRATE_FLOW_CONFIG|FLOW_CONFIG' /root/astrate-mule/internal/flow/blocks/container/block.go | head
45:	Nested        map[string]any // opaque JSON for ASTRATE_FLOW_CONFIG

$ rg -n 'FLOW_CONFIG' /root/astrate-mule/internal/flow/blocks/container/docker.go /root/astrate-mule/internal/flow/blocks/schema.go /root/astrate-mule/internal/flow/blocks/info.go
/root/astrate-mule/internal/flow/blocks/container/docker.go:19:	// FlowConfigJSON is passed as ASTRATE_FLOW_CONFIG.
/root/astrate-mule/internal/flow/blocks/container/docker.go:99:		args = append(args, "-e", "ASTRATE_FLOW_CONFIG="+spec.FlowConfigJSON)
/root/astrate-mule/internal/flow/blocks/container/docker.go:201:// encodeFlowConfigJSON marshals nested config for ASTRATE_FLOW_CONFIG.
/root/astrate-mule/internal/flow/blocks/schema.go:19:	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as ASTRATE_FLOW_CONFIG"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
/root/astrate-mule/internal/flow/blocks/info.go:116:		Config:  "image (required), config (object→ASTRATE_FLOW_CONFIG), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",

$ cat /root/astrate-mule/docs/Makefile
.PHONY: serve build clean install sync

install:
	pip install -r requirements.txt

sync:
	cp DESIGN.md site/DESIGN.md
	cp ROADMAP.md site/ROADMAP.md
	cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
	rm -rf site/api && cp -r api site/api
	rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui

serve: sync
	mkdocs serve -f mkdocs.yml

build: sync
	mkdocs build -f mkdocs.yml

clean:
	rm -rf ../site-dist

Let me verify the compat version response shape and check the docs Makefile for the build step.
→ Read internal/observability/compat.go [offset=36]
I've compared the astrate-native surface (the recipe's fifth option — no prior docs-sync run covers it). Code routes match the spec's documented paths one-for-one, so the drift is in response accuracy, not missing/phantom paths. Let me note the findings and append the proposals.
→ Read .mule/todo.md [offset=130]

$ cat >> /root/astrate-mule/.mule/todo.md << 'EOF'

- [ ] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmmanagement/health`, `GET /pairing/health` — in docs/api/astrate_native_api.yaml, and widen the shared `DataEnvelopeHealthStatus.status` enum (yaml:301-303) from `[ok]` to `[ok, unhealthy]` — the handler answers 503 `{"data":{"status":"unhealthy"}}` when the readiness probe fails (internal/observability/compat.go:19-34, deviation 18), but the spec documents only 200 with `status: ok` for all three paths. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
EOF
(no output)

$ tail -4 /root/astrate-mule/.mule/todo.md
- [x] realm-pure-helper-tests [auto]: add container-free unit tests for the four untested pure helpers in internal/realm. Create `internal/realm/helpers_test.go` (no build tag) covering: (1) `normaliseIfaceName` — lowercase, hyphen strip, mixed case+hyphens, empty string, already-normalized passthrough; (2) `joinPEM` — empty slice, single key, two keys separated by `\n`; (3) `triggerErrorBody` — empty action+simple_triggers, action-only, simple_triggers-only with index-aligned `{}` substitution, both present; (4) `majorParam` — valid positive, zero, negative (→404), non-numeric (→404). All tests use httptest.ResponseRecorder for majorParam; the rest are direct function calls. No Docker.

- [ ] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmmanagement/health`, `GET /pairing/health` — in docs/api/astrate_native_api.yaml, and widen the shared `DataEnvelopeHealthStatus.status` enum (yaml:301-303) from `[ok]` to `[ok, unhealthy]` — the handler answers 503 `{"data":{"status":"unhealthy"}}` when the readiness probe fails (internal/observability/compat.go:19-34, deviation 18), but the spec documents only 200 with `status: ok` for all three paths. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

Ran the docs-sync recipe on the **astrate-native** surface (the only one of the five never previously covered). Route comparison: all 9 documented paths match code one-for-one (`astrate/v1` health/readiness/metrics, the three `/{service}/health` compat endpoints, both WebSockets) — no phantom or absent routes. Drift was in response accuracy, not presence.

Findings → 2 task lines appended (tagged `[auto]`, as todo.md requires for proposal output):
- Compat-health endpoints: code answers `503 {"data":{"status":"unhealthy"}}` on probe failure (compat.go:19-34), spec documents only 200 and pins the shared enum to `[ok]`.
- Four unauthenticated `/{service}/version` compat endpoints registered for all four services (main.go:393-395) are documented nowhere.

Noted but not proposed (different surfaces / already escalated): the `ASTRATE_HOUSEKEEPING_*` missing config keys are already in for-giulio.md; `ASTRATE_FLOW_CONFIG` is an env var Astrate *sets for spawned containers*, not a config key; realm-scoped appengine/pairing version twins belong to their own surfaces. No evidence file was required, no `gh issue create` asked for, so deliverable is just the task lines.

Done: docs-sync for astrate-native surface — two drift-fix tasks proposed
Files: .mule/todo.md
Verified: rg route/path comparison (internal/* vs docs/api/astrate_native_api.yaml) -> all paths match, drift confined to response accuracy
Unsure: nothing
Follow-ups: docs-sync-native-compat-health-503, docs-sync-native-compat-version-endpoints
