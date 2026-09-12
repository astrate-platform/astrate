slug: recipe-docs-sync
verdict: proposed
at:  90c485c
ran: 2026-09-12T11:04:49Z on DietPi in 390s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls docs/api/ && ls docs/site/ 2>/dev/null | head && ls internal/
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
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

$ git status --short | head -20; git log --oneline -3
90c485c mule: log interfaceschema-trie-interior-literal-param
705795f mule: interfaceschema-trie-interior-literal-param [auto]: pin the rule that an interior-only literal trie child does not shadow its parametric sibling at the matching depth (pkg/interfaceschema/trie.go:124-151) — purpose-build a trie with "/a/b/c" and "/a/%{p}", assert Match("/a/b") returns the /a/%{p} mapping and Match("/a/b/c") the /a/b/c one; statement coverage cannot see this branch (matchNode reports 100%), so a "fix" that misses on shallow paths through an interior literal node sails through the suite.
009a6a4 mule: log interfaceschema-enum-roundtrip-coverage

Working tree clean. I'll run the docs-sync recipe on one surface — `appengine`.
$ rg -n '^\s{2}/' docs/api/astarte_appengine_api.yaml
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

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*appengine*'
(no output)

The first grep found nothing — handlers likely live elsewhere. Let me locate them.
$ ls internal/appengine/ 2>/dev/null; ls internal/flowapi 2>/dev/null; rg -ln 'HandleFunc|r\.(Get|Post|Put|Delete|Patch)\(' internal/ | head -30
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
composite_wire_test.go
http.go
http_test.go
service.go
service_integration_test.go
service_test.go
userblocks.go
userblocks_test.go
internal/flow/blocks/httpblocks_test.go
internal/flow/blocks/httpblocks.go
internal/engine/forward/http_test.go
internal/httpx/notfound_test.go
internal/httpx/cors.go
internal/swagger/swagger.go
internal/swagger/swagger_test.go
internal/appengine/stream/ws.go
internal/appengine/http.go
internal/engine/triggers/actions_test.go
internal/auth/middleware.go
internal/pairing/http.go
internal/observability/compat.go
internal/observability/health.go
internal/observability/compat_version_test.go

✱ Grep "HandleFunc|\.(Get|Post|Put|Delete|Patch)\(" in internal/appengine · 22 matches

$ rg -n 'RegisterHandler|\.Handle|AddRoute|route|mux|http\.Method' internal/appengine/http.go | head -40
30:	require func(http.Handler) http.Handler
38:// Mount registers the routes on mux (paths wire-frozen to upstream
40:func (a *API) Mount(mux *http.ServeMux) {
41:	h := func(f http.HandlerFunc) http.Handler { return a.require(f) }
43:	mux.Handle("GET "+base+"/devices", h(a.listDevices))
44:	mux.Handle("GET "+base+"/stats/devices", h(a.devicesStats))
45:	mux.Handle("GET "+base+"/devices/{device}", h(a.getDevice))
46:	mux.Handle("PATCH "+base+"/devices/{device}", h(a.patchDevice))
47:	mux.Handle("GET "+base+"/devices-by-alias/{alias}", h(a.getDeviceByAlias))
48:	mux.Handle("GET "+base+"/devices/{device}/interfaces", h(a.listDeviceInterfaces))
50:	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}", h(a.getData))
51:	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.getData))
52:	mux.Handle("PUT "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
53:	mux.Handle("POST "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
54:	mux.Handle("DELETE "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteData))
57:	// ones; Go 1.22 precedence resolves the overlap with the group routes
59:	mux.Handle("PATCH "+base+"/devices-by-alias/{alias}", h(a.patchDeviceByAlias))
60:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces", h(a.listInterfacesByAlias))
61:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces/{interface}", h(a.getDataByAlias))
62:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.getDataByAlias))
63:	mux.Handle("PUT "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.putDataByAlias))
64:	mux.Handle("POST "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.putDataByAlias))
65:	mux.Handle("DELETE "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.deleteDataByAlias))
67:	mux.Handle("GET "+base+"/groups", h(a.listGroups))
68:	mux.Handle("POST "+base+"/groups", h(a.createGroup))
69:	mux.Handle("GET "+base+"/groups/{group}", h(a.getGroup))
70:	mux.Handle("GET "+base+"/groups/{group}/devices", h(a.listGroupDevices))
71:	mux.Handle("POST "+base+"/groups/{group}/devices", h(a.addGroupDevice))
72:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}", h(a.getGroupDevice))
73:	mux.Handle("PATCH "+base+"/groups/{group}/devices/{device}", h(a.patchGroupDevice))
74:	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}", h(a.removeGroupDevice))
75:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces", h(a.listInterfacesInGroup))
76:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}", h(a.getDataInGroup))
77:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.getDataInGroup))
78:	mux.Handle("PUT "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.putDataInGroup))
79:	mux.Handle("POST "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.putDataInGroup))
80:	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteDataInGroup))

→ Read internal/appengine/http.go [limit=85]

$ rg -n '^\s{2}/|^\s{6}(get|post|put|patch|delete):' docs/api/astarte_appengine_api.yaml
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

$ awk '/^  \//{path=$0} /^\s{4,6}(get|post|put|patch|delete):/{gsub(/[: ]/,""); print path, $0}' docs/api/astarte_appengine_api.yaml
(no output)

→ Read docs/api/astarte_appengine_api.yaml [limit=100]

$ rg -n '^  /appengine|^\s{4}(get|post|put|patch|delete):' docs/api/astarte_appengine_api.yaml
26:  /appengine/v1/{realm}/devices:
27:    get:
98:  /appengine/v1/{realm}/stats/devices:
99:    get:
136:  /appengine/v1/{realm}/devices/{device}:
137:    get:
161:    patch:
210:  /appengine/v1/{realm}/devices-by-alias/{alias}:
211:    get:
240:    patch:
292:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces:
293:    get:
324:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
325:    get:
369:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
370:    get:
419:    put:
465:    post:
511:    delete:
545:  /appengine/v1/{realm}/devices/{device}/interfaces:
546:    get:
573:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
574:    get:
613:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
614:    get:
656:    put:
697:    post:
738:    delete:
769:  /appengine/v1/{realm}/groups:
770:    get:
795:    post:
841:  /appengine/v1/{realm}/groups/{group}:
842:    get:
869:  /appengine/v1/{realm}/groups/{group}/devices:
870:    get:
917:    post:
957:  /appengine/v1/{realm}/groups/{group}/devices/{device}:
958:    get:
985:    patch:
1035:    delete:
1056:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces:
1057:    get:
1084:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}:
1085:    get:
1125:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}:
1126:    get:
1169:    put:
1211:    post:
1253:    delete:

→ Read docs/api/astarte_appengine_api.yaml [offset=136, limit=74]
All 17 documented paths and all methods match the code routes. Now checking status codes/response fields for a few endpoints.
→ Read docs/api/astarte_appengine_api.yaml [offset=613, limit=156]
Now let me check the handlers' actual status codes.

$ rg -n 'func \(a \*API\) (putData|deleteData|getData|getDevice|createGroup|getGroupDevice)\(|WriteHeader\(|http\.Status' internal/appengine/http.go
98:	_ = astarteapi.WriteDataWithLinks(w, http.StatusOK, body,
136:	_ = astarteapi.WriteData(w, http.StatusOK, map[string]int64{
142:func (a *API) getDevice(w http.ResponseWriter, r *http.Request) {
148:	_ = astarteapi.WriteData(w, http.StatusOK, st)
157:	_ = astarteapi.WriteData(w, http.StatusOK, st)
166:	_ = astarteapi.WriteData(w, http.StatusOK, names)
223:	_ = astarteapi.WriteData(w, http.StatusOK, st)
231:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Bad request")
244:	_ = astarteapi.WriteData(w, http.StatusOK, st)
265:		_ = astarteapi.WriteDataWithMetadata(w, http.StatusOK, t.Data, t.Metadata)
268:	_ = astarteapi.WriteData(w, http.StatusOK, data)
271:func (a *API) getData(w http.ResponseWriter, r *http.Request) {
292:func (a *API) putData(w http.ResponseWriter, r *http.Request) {
299:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
308:	w.WriteHeader(http.StatusOK)
311:func (a *API) deleteData(w http.ResponseWriter, r *http.Request) {
318:	w.WriteHeader(http.StatusNoContent)
328:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
337:	w.WriteHeader(http.StatusOK)
347:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
356:	w.WriteHeader(http.StatusOK)
366:	w.WriteHeader(http.StatusNoContent)
376:	w.WriteHeader(http.StatusNoContent)
387:	_ = astarteapi.WriteData(w, http.StatusOK, names)
398:func (a *API) createGroup(w http.ResponseWriter, r *http.Request) {
406:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
410:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
414:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
423:	_ = astarteapi.WriteData(w, http.StatusCreated, body)
433:	_ = astarteapi.WriteData(w, http.StatusOK, map[string]string{"group_name": r.PathValue("group")})
450:	_ = astarteapi.WriteDataWithLinks(w, http.StatusOK, body,
496:	w.WriteHeader(http.StatusCreated)
504:	w.WriteHeader(http.StatusNoContent)
509:func (a *API) getGroupDevice(w http.ResponseWriter, r *http.Request) {
515:	_ = astarteapi.WriteData(w, http.StatusOK, st)
536:	_ = astarteapi.WriteData(w, http.StatusOK, st)
647:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity, fe)
649:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Invalid alias")
651:		_ = astarteapi.WriteError(w, http.StatusConflict, "Alias already in use")
653:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Alias tag not found")
655:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Invalid attributes")
657:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Attribute key not found")
659:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
661:		_ = astarteapi.WriteError(w, http.StatusConflict, "Group already exists")
663:		_ = astarteapi.WriteError(w, http.StatusConflict, "Device already in group")
665:		_ = astarteapi.WriteError(w, http.StatusConflict, "Already exists")
667:		_ = astarteapi.WriteError(w, http.StatusNotFound, "Group not found")
669:		_ = astarteapi.WriteError(w, http.StatusNotFound, "Path not found")
674:		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Cannot write to device owned resource")
676:		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Cannot write to read-only resource")
678:		_ = astarteapi.WriteError(w, http.StatusNotFound, "Interface not found in device introspection")
680:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Endpoint not found")

→ Read internal/appengine/http.go [limit=180, offset=290]
→ Read internal/appengine/http.go [limit=70, offset=630]

$ rg -n '^\s{8}422|^\s{8}400|^\s{8}404|^\s{8}401|^\s{8}500' docs/api/astarte_appengine_api.yaml | head -60
(no output)

$ rg -n '"4[0-9][0-9]"|"5[0-9][0-9]"|"2[0-9][0-9]"' docs/api/astarte_appengine_api.yaml | rg -v '^docs/api/.+:\s*[#-]' | awk -F'"' '{print $2}' | sort | uniq -c
26 200
      2 201
      4 204
     11 400
     32 401
     28 404
      2 409
     11 422
     32 500

$ rg -n 'WriteBadRequest|WriteDeviceNotFound|WriteUnauthorized' internal/httpx/ internal/astarteapi/ pkg/astarteapi/ 2>/dev/null | head; rg -n 'func WriteBadRequest|func WriteDeviceNotFound|func WriteUnauthorized|func WriteError|func WriteData\b' -g '*.go' | head
pkg/astarteapi/envelope.go:190:// WriteBadRequest writes the canonical 400 envelope.
pkg/astarteapi/envelope.go:191:func WriteBadRequest(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:195:// WriteUnauthorized writes the canonical 401 envelope.
pkg/astarteapi/envelope.go:196:func WriteUnauthorized(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:210:// WriteDeviceNotFound writes the canonical 404 "Device not found" envelope.
pkg/astarteapi/envelope.go:211:func WriteDeviceNotFound(w http.ResponseWriter) error {
pkg/astarteapi/envelope_test.go:72:			write: astarteapi.WriteBadRequest,
pkg/astarteapi/envelope_test.go:76:			write: astarteapi.WriteUnauthorized,
pkg/astarteapi/envelope_test.go:88:			write: astarteapi.WriteDeviceNotFound,
pkg/astarteapi/envelope.go:126:func WriteData(w http.ResponseWriter, status int, v any) error {
pkg/astarteapi/envelope.go:163:func WriteError(w http.ResponseWriter, status int, detail string) error {
pkg/astarteapi/envelope.go:191:func WriteBadRequest(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:196:func WriteUnauthorized(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:211:func WriteDeviceNotFound(w http.ResponseWriter) error {

→ Read docs/api/astarte_appengine_api.yaml [limit=100, offset=957]
405 is never documented in the appengine YAML, yet `writeError` emits it (http.go:673-676). Let me verify the mirror surfaces and the group-device endpoints before proposing.
→ Read docs/api/astarte_appengine_api.yaml [limit=177, offset=369]
→ Read docs/api/astarte_appengine_api.yaml [limit=73, offset=769]

$ rg -n 'UnprocessableEntity|ValidationErrors|MethodNotAllowed|"405"|Value size' docs/api/astarte_appengine_api.yaml
206:          $ref: "#/components/responses/UnprocessableEntity"
288:          $ref: "#/components/responses/UnprocessableEntity"
365:          $ref: "#/components/responses/ValidationErrors"
415:          $ref: "#/components/responses/ValidationErrors"
609:          $ref: "#/components/responses/ValidationErrors"
652:          $ref: "#/components/responses/ValidationErrors"
837:          $ref: "#/components/responses/ValidationErrors"
913:          $ref: "#/components/responses/ValidationErrors"
1031:          $ref: "#/components/responses/UnprocessableEntity"
1121:          $ref: "#/components/responses/ValidationErrors"
1165:          $ref: "#/components/responses/ValidationErrors"
1663:    ValidationErrors:
1674:    UnprocessableEntity:

$ rg -n 'func \(s \*Service\) CreateGroup|ErrGroupAlreadyExists' internal/appengine/service.go
44:// ErrGroupAlreadyExists marks a duplicate group name (maps to 409 "Group
47:var ErrGroupAlreadyExists = fmt.Errorf("Group already exists: %w", store.ErrAlreadyExists) //nolint:staticcheck // ST1005: upstream wire text
462:func (s *Service) CreateGroup(ctx context.Context, realm, name string, devices []string) error {
495:			return fmt.Errorf("%w: %w", ErrGroupAlreadyExists, err)

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
internal/engine/engine_integration_test.go:ASTRATE_TEST_DSN
internal/flow/blocks/container/block.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/container/docker.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/info.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/schema.go:ASTRATE_FLOW_CONFIG
internal/store/crypto.go:ASTRATE_MASTER_KEY
internal/store/crypto.go:ASTRATE_MASTER_KEY_FILE
internal/store/store_test.go:ASTRATE_TEST_DSN
internal/testutil/pg.go:ASTRATE_TEST_DSN

Three solid docs/api drift items confirmed. Let me also run the recipe's cheap config-key sweep before writing proposals.
$ rg -o '\bASTRATE_[A-Z_]+' docs/site/configuration-reference.md | sort -u
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

Config sweep: code reads two `ASTRATE_HOUSEKEEPING_*` keys absent from `configuration-reference.md`. Verifying they're live in config.go.
$ rg -n 'HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION|HOUSEKEEPING_REALM_DELETION_DISABLED' internal/config/config.go
270:		"ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
271:		"HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
288:	if v, ok := os.LookupEnv("ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"); ok {

$ ls .mule/; tail -30 .mule/for-giulio.md 2>/dev/null || echo "NO FILE"
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
  as a userland toolchain on the Legion Go (`~/.local/go`, no root, `rm -rf` to undo). The
  Pi still cannot run `-race` (39-bit VMA kernel vs the 48 ThreadSanitizer needs), so its
  gate remains `go vet ./... && go test ./...` — but race coverage now exists on the Legion
  Go, where the full suite runs clean in ~40s on 16 cores. The standing `race-check` task is
  the concurrency gate. Concurrency work is queueable again, provided the race-check runs
  after it.
- ~~golangci-lint is not installed on the Pi~~ — **resolved 2026-07-28**: installed v2.12.2
  via `go install` (the prebuilt-binary installer's published sha256 for linux/arm64 did not
  verify, reproducibly, so built from source instead). `mule.service` now sets
  `MULE_LINT_CMD=golangci-lint run ./...` and has `/root/go/bin` on `PATH`; the lint gate
  runs starting with the next tick.
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

$ tail -15 .mule/todo.md
- [x] realm-pure-helper-tests [auto]: add container-free unit tests for the four untested pure helpers in internal/realm. Create `internal/realm/helpers_test.go` (no build tag) covering: (1) `normaliseIfaceName` — lowercase, hyphen strip, mixed case+hyphens, empty string, already-normalized passthrough; (2) `joinPEM` — empty slice, single key, two keys separated by `\n`; (3) `triggerErrorBody` — empty action+simple_triggers, action-only, simple_triggers-only with index-aligned `{}` substitution, both present; (4) `majorParam` — valid positive, zero, negative (→404), non-numeric (→404). All tests use httptest.ResponseRecorder for majorParam; the rest are direct function calls. No Docker.

- [x] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmmanagement/health`, `GET /pairing/health` — in docs/api/astrate_native_api.yaml, and widen the shared `DataEnvelopeHealthStatus.status` enum (yaml:301-303) from `[ok]` to `[ok, unhealthy]` — the handler answers 503 `{"data":{"status":"unhealthy"}}` when the readiness probe fails (internal/observability/compat.go:19-34, deviation 18), but the spec documents only 200 with `status: ok` for all three paths. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-status-enum [auto]: fix the `PairingInfo.status` enum in docs/api/astarte_pairing_api.yaml from `[confirmed, pending, denied, expired]` (yaml:412) to `[confirmed, pending, inhibited]` — the handler never emits `denied` or `expired`: `service.Info` maps only `DeviceStatusInhibited`→"inhibited", `DeviceStatusConfirmed`→"confirmed", and everything else→"pending" (internal/pairing/service.go:296-299, statuses defined in store/devices.go:19-27; upstream-parity per the comment at service.go:282-285). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-version-endpoint [auto]: add the undocumented `GET /pairing/v1/{realm}/version` route to docs/api/astarte_pairing_api.yaml — it exists in code (cmd/astrate/main.go:398, unauthenticated, upstream-parity per main.go:387-391) and answers 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38); mirror the RM spec's getVersion op (astarte_realm_management_api.yaml:532) without its a_rma security. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-delete-interface-status [auto]: fix the `deleteInterface` response in docs/api/astarte_realm_management_api.yaml — the spec documents `422` with examples "Interface major version is not 0, can't be deleted" / "Cannot delete an interface that is used by a device introspection" (yaml:255-272), but the handler returns `403` with "Interface can't be deleted" / "Interface can't be deleted since it's currently used" (internal/realm/http.go:387-391, `writeInterfaceError`; both are 403 `StatusForbidden`). The spec's detail strings were lifted from the `writeError` branch (http.go:414-419) which DELETE /interfaces does not use. Fix status to `403` and update the two examples to match the actual wire messages. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-mapping-required-encrypted [auto]: add `required` (boolean, datastream-only) and `encrypted` (boolean, datastream-only) to the `InterfaceMapping` schema in docs/api/astarte_realm_management_api.yaml — the parser accepts and persists both on datastream mappings (internal/interfaceschema/parse.go:410-414), rejects them on properties mappings (parse.go:395-398), the detailed listing emits them (`detailed.go:69`), and the compatibility check compares them (interfaceschema/compat.go:121-123); the schema (yaml:820-862) currently has `additionalProperties: false` and neither field, so the OpenAPI would mark a valid datastream interface install body as invalid. Also add both fields to the detailed-listing example (yaml:70-85). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-put-auth-422 [auto]: add the missing `422` (ValidationError) response to `PUT /realmmanagement/v1/{realm}/config/auth` in docs/api/astarte_realm_management_api.yaml — `SetAuthKey` returns `ErrValidation` for a blank `jwt_public_key_pem` (internal/realm/service.go:632-640), which `writeError` maps to 422; the spec (yaml:431-464) documents only 204/400/401/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-rm-version-example [auto]: fix the `getVersion` example in docs/api/astarte_realm_management_api.yaml from `data: "1.1.0"` (yaml:557) to `data: "1.2.2"` — `APICompatVersion` is `"1.2.2"` (internal/realm/service.go:588); the native spec's four compat-version examples already use "1.2.2". Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1
- [x] interfaceschema-compat-attrs-flip-test [auto]: extend TestCheckMinorUpgradeTable (pkg/interfaceschema/compat_test.go) to assert the four currently-unpinned immutable attributes — retention, expiry, database_retention_policy/ttl (datastream old/next pair) and allow_unset (properties old/next pair) — each flipped on upgrade and rejected as ErrIncompatibleEndpointChange with an at-immutability acceptance twin; only type/reliability/explicit_timestamp/required/encrypted are tested today (sameMappingAttributes, compat.go:102, sits at 66.7% stmt coverage), so a regression relaxing expiry or DB-TTL immutability passes the suite.
- [x] interfaceschema-object-attrs-uniformity [auto]: pin the object-aggregation attribute-uniformity rule for the four unasserted attributes — retention, expiry, explicit_timestamp, database_retention_policy/ttl — via ParseInterface rejections ("must all share the same <attr>", pkg/interfaceschema/parse.go:626-647) each paired with an at-uniformity acceptance twin; only reliability is pinned today (testdata/invalid/17-object-mixed-attrs.json, sameObjectAttributes at 54.5% stmt coverage).
- [x] interfaceschema-malformed-placeholder-fixtures [auto]: add invalid fixtures + manifest entries for parseSegment's four untested param branches (pkg/interfaceschema/parse.go:509-517) — "/%{", "/%ax", "/%{x", "/%{1}" — each rejected with the malformed-placeholder / invalid-placeholder-name message; the only placeholder-syntax fixture today is the mid-segment 12-endpoint-partial-placeholder.json (ContainsAny branch), so a regression accepting a dangling "/%{" passes every test.
- [x] interfaceschema-enum-roundtrip-coverage [auto]: extend TestEnumRoundTrips (pkg/interfaceschema/parse_test.go:244) to cover Ownership, Aggregation, Retention, DatabaseRetentionPolicy — the four enums it omits (Retention.String/MarshalJSON and DatabaseRetentionPolicy.String at 0% stmt coverage in the package's own tests) — so a typo'd wire string like "retention":"descard" or "aggregation":"obilver" fails the suite instead of passing.
- [x] interfaceschema-trie-interior-literal-param [auto]: pin the rule that an interior-only literal trie child does not shadow its parametric sibling at the matching depth (pkg/interfaceschema/trie.go:124-151) — purpose-build a trie with "/a/b/c" and "/a/%{p}", assert Match("/a/b") returns the /a/%{p} mapping and Match("/a/b/c") the /a/b/c one; statement coverage cannot see this branch (matchNode reports 100%), so a "fix" that misses on shallow paths through an interior literal node sails through the suite.

Confirmed. Now writing the deliverable: 3 `docs/api` fix tasks in todo.md + one `for-giulio.md` escalation for the config-key drift. Let me check for-giulio.md's format first.
→ Read .mule/for-giulio.md [limit=40]
Let me verify what those two housekeeping keys do before writing the escalation.
→ Read internal/config/config.go [limit=30, offset=264]
Both keys confirmed live (#73, #75). Writing the deliverables now.
→ Read .mule/todo.md [offset=1, limit=8]
→ Read .mule/todo.md [offset=145]
Now appending the three proposal lines:

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -151,4 +151,7 @@
 - [x] interfaceschema-object-attrs-uniformity [auto]: pin the object-aggregation attribute-uniformity rule for the four unasserted attributes — retention, expiry, explicit_timestamp, database_retention_policy/ttl — via ParseInterface rejections ("must all share the same <attr>", pkg/interfaceschema/parse.go:626-647) each paired with an at-uniformity acceptance twin; only reliability is pinned today (testdata/invalid/17-object-mixed-attrs.json, sameObjectAttributes at 54.5% stmt coverage).
 - [x] interfaceschema-malformed-placeholder-fixtures [auto]: add invalid fixtures + manifest entries for parseSegment's four untested param branches (pkg/interfaceschema/parse.go:509-517) — "/%{", "/%ax", "/%{x", "/%{1}" — each rejected with the malformed-placeholder / invalid-placeholder-name message; the only placeholder-syntax fixture today is the mid-segment 12-endpoint-partial-placeholder.json (ContainsAny branch), so a regression accepting a dangling "/%{" passes every test.
 - [x] interfaceschema-enum-roundtrip-coverage [auto]: extend TestEnumRoundTrips (pkg/interfaceschema/parse_test.go:244) to cover Ownership, Aggregation, Retention, DatabaseRetentionPolicy — the four enums it omits (Retention.String/MarshalJSON and DatabaseRetentionPolicy.String at 0% stmt coverage in the package's own tests) — so a typo'd wire string like "retention":"descard" or "aggregation":"obilver" fails the suite instead of passing.
 - [x] interfaceschema-trie-interior-literal-param [auto]: pin the rule that an interior-only literal trie child does not shadow its parametric sibling at the matching depth (pkg/interfaceschema/trie.go:124-151) — purpose-build a trie with "/a/b/c" and "/a/%{p}", assert Match("/a/b") returns the /a/%{p} mapping and Match("/a/b/c") the /a/b/c one; statement coverage cannot see this branch (matchNode reports 100%), so a "fix" that misses on shallow paths through an interior literal node sails through the suite.
+- [ ] docs-sync-ae-write-405 [auto]: add the missing `405` response to the nine data-write ops in docs/api/astarte_appengine_api.yaml — PUT/POST/DELETE `/devices/{device}/interfaces/{interface}/{path}`, the three `/devices-by-alias/{alias}` twins, and the three `/groups/{group}/devices/{device}` twins. `writeError` emits `405 Method Not Allowed` with "Cannot write to device owned resource" (`ErrNotServerOwned`) and "Cannot write to read-only resource" (`ErrNotAProperty`) at internal/appengine/http.go:673-676, reachable from putData/putDataByAlias/putDataInGroup (http.go:292-357) and deleteData/deleteDataByAlias/deleteDataInGroup (http.go:311-377); `"405"` currently appears nowhere in the spec (those ops document only 200/204, 400, 401, 404, 500). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-ae-write-value-422 [auto]: add the missing `422` response to the six PUT/POST data-write ops in docs/api/astarte_appengine_api.yaml (device-scoped, by-alias, and in-group at `/interfaces/{interface}/{path}`). A value larger than `maxValueBytes` (64 KiB, internal/appengine/http.go:24) is answered `422` "Value size exceeds size limits" by putData/putDataByAlias/putDataInGroup (http.go:298-301, 327-329, 346-348), but the specs for PUT/POST (yaml:686-695, 726-736, 453-463, 499-509, 1176-1186, 1218-1228) list only 200/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-ae-post-groups-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups` in docs/api/astarte_appengine_api.yaml — `Service.CreateGroup` returns `ErrGroupAlreadyExists` (internal/appengine/service.go:495, mapped at http.go:660-661) which the handler answers as `409` "Group already exists"; the spec (yaml:821-839) documents only 201/400/401/422/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

Now the for-giulio.md escalation for the two undocumented config keys:
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,10 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
+
 - **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pile is now 7 straight days and worth a look.** Twelve open issues. Still no machine-checkable fix candidates, so no task lines: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-10 run: **#100**, today's alarm (created 10:55Z, "nothing has landed in 15h"). #99–#94 are the same `mule-alarm` (one per day 2026-09-05→09-10, ~11:00Z each, zero comments) — telemetry, not code issues, so never proposable. But seven consecutive alarms is past the "low-activity window" wording earlier runs used: the mule has landed no commit since ~2026-09-04/05, and `main`'s queue copy shows the top line as `- [ ]` while the landed work below is all `[!] BLOCKED` (`wrote nothing` / `tests failed`) — though `mule/queue` is the authoritative copy. **Proposal: close #94–#99** — each was a one-day event, superseded by next day's alarm, never actionable; leave #100 (live today) to self-expire. If the idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
 
 - **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
 

$ git status --short; git diff --stat
M .mule/for-giulio.md
 M .mule/todo.md
 .mule/for-giulio.md | 2 ++
 .mule/todo.md       | 3 +++
 2 files changed, 5 insertions(+)

Done: docs-sync recipe on the appengine surface + config-key sweep. Compared all 17 documented paths/methods against `internal/appengine/http.go:43-80` (full match), checked status codes on 4 endpoint groups, ran the `ASTRATE_*` key sweep. Found three `docs/api` glitches and one config-key omission.

Drift found:
- **405 not documented anywhere** — `writeError` emits 405 "Cannot write to device owned resource" / "Cannot write to read-only resource" (http.go:673-676) from all six data-write handlers (http.go:292-377).
- **422 omitted on PUT/POST data writes** — `>64KiB` value answers 422 "Value size exceeds size limits" (http.go:298-301 etc), specs show only 200/400/401/404/500.
- **409 omitted on POST /groups** — duplicate group answers 409 "Group already exists" (service.go:495 → http.go:660-661).
- **Two config keys read but absent from configuration-reference.md** (`ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`, `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`, config.go:268-293).

Clean surfaces noted: all path/method registrations match; getDevice, patchDevice, group-device status codes check out; GET data-query params' 422 matches.

Done: compared appengine spec vs handlers; proposed 3 docs/api fixes, escalated 1 config-key drift; no code changed
Files: .mule/todo.md, .mule/for-giulio.md
Verified: rg path/method/status-code cross-check + config key sweep (`internal/config` vs docs) -> found 4 items, no tests run (proposal-only)
Unsure: whether the two undocumented config keys are deliberately a subset (flagged in for-giulio.md); whether 409 can surface on non-group write ops (not claimed)
Follow-ups: docs-sync-ae-write-405, docs-sync-ae-write-value-422, docs-sync-ae-post-groups-409 appended to .mule/todo.md
