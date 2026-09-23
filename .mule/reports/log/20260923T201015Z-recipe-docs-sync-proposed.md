slug: recipe-docs-sync
verdict: proposed
at:  55bba71
ran: 2026-09-23T20:10:15Z on DietPi in 779s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
→ Read .mule/todo.md
✱ Glob "docs/api/*.yaml" 5 matches
✱ Glob ".mule/**" 0 matches
→ Read .mule/todo.md [offset=174]
→ Read docs/api/astrate_native_api.yaml
Let me look at the docs-sync recipe surface. I'll pick the astrate-native surface. Let me examine the spec and how routes are registered.

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc|\.Handle\(|Methods' cmd/astrate/main.go internal/observability/ internal/swagger/ 2>/dev/null | head -80
cmd/astrate/main.go:433:	mux.Handle("GET /appengine/v1/{realm}/version",
cmd/astrate/main.go:435:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
internal/swagger/swagger.go:20:	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
internal/swagger/swagger.go:23:	mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(uiRoot))))
internal/swagger/swagger.go:27:	mux.Handle("GET /api/", http.StripPrefix("/api/", http.FileServer(http.FS(apiRoot))))
internal/observability/health.go:44:	mux.HandleFunc("GET /astrate/v1/health", h.handleHealth)
internal/observability/health.go:45:	mux.HandleFunc("GET /astrate/v1/readiness", h.handleReadiness)
internal/observability/health.go:46:	mux.Handle("GET /astrate/v1/metrics", h.metrics)
internal/swagger/swagger_test.go:35:		loc := resp.Header.Get("Location")
internal/observability/compat.go:20:	mux.HandleFunc("GET /"+service+"/health", func(w http.ResponseWriter, r *http.Request) {
internal/observability/compat.go:51:	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
internal/observability/compat_version_test.go:31:		if ct := resp.Header.Get("Content-Type"); ct == "" {
internal/observability/compat_version_test.go:39:	mux.HandleFunc("GET /v1/{realm}/version", VersionHandler("x"))

The native surface looks fully covered by recent tasks. Let me check the other surfaces for drift. Let me run the recipe's drift detection for pairing and realm-management.
$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/pairing/*.go
internal/pairing/http.go:81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/pairing/http.go:86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
internal/pairing/http.go:88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
internal/pairing/http.go:90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
internal/pairing/http.go:344:	header := r.Header.Get("Authorization")

$ rg -n '^\s{2}/' docs/api/astarte_pairing_api.yaml
27:  /pairing/v1/{realm}/health:
70:  /pairing/v1/{realm}/version:
97:  /pairing/v1/{realm}/agent/devices:
149:  /pairing/v1/{realm}/agent/devices/{deviceID}:
172:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
224:  /pairing/v1/{realm}/devices/{deviceID}:
258:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:

Pairing looks fully documented. Let me check realm-management routes vs spec, and housekeeping.
$ rg -n 'HandleFunc|Handle\(|\.(Get|Post|Put|Delete|Patch)\("' internal/realm/http.go | head -50
38:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces", h(a.listInterfaces))
39:	mux.Handle("POST /realmmanagement/v1/{realm}/interfaces", h(a.installInterface))
40:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}", h(a.listInterfaceMajors))
41:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.getInterface))
42:	mux.Handle("PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.updateInterface))
43:	mux.Handle("DELETE /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.deleteInterface))
44:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers", h(a.listTriggers))
45:	mux.Handle("POST /realmmanagement/v1/{realm}/triggers", h(a.createTrigger))
46:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers/{name}", h(a.getTrigger))
47:	mux.Handle("DELETE /realmmanagement/v1/{realm}/triggers/{name}", h(a.deleteTrigger))
48:	mux.Handle("GET /realmmanagement/v1/{realm}/config/auth", h(a.getAuth))
49:	mux.Handle("PUT /realmmanagement/v1/{realm}/config/auth", h(a.putAuth))
50:	mux.Handle("GET /realmmanagement/v1/{realm}/config/device_registration_limit", h(a.getRegistrationLimit))
51:	mux.Handle("GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention", h(a.getDatastreamMaximumStorageRetention))
52:	mux.Handle("GET /realmmanagement/v1/{realm}/version", h(a.getVersion))
53:	mux.Handle("DELETE /realmmanagement/v1/{realm}/devices/{device}", h(a.deleteDevice))
54:	mux.Handle("GET /realmmanagement/v1/{realm}/policies", h(a.listPolicies))
55:	mux.Handle("POST /realmmanagement/v1/{realm}/policies", h(a.createPolicy))
56:	mux.Handle("GET /realmmanagement/v1/{realm}/policies/{name}", h(a.getPolicy))
57:	mux.Handle("DELETE /realmmanagement/v1/{realm}/policies/{name}", h(a.deletePolicy))
150:	if r.URL.Query().Get("detailed") == "true" {

$ rg -n '^\s{2}/' docs/api/astarte_realm_management_api.yaml
27:  /realmmanagement/v1/{realm}/interfaces:
138:  /realmmanagement/v1/{realm}/interfaces/{name}:
167:  /realmmanagement/v1/{realm}/interfaces/{name}/{major}:
279:  /realmmanagement/v1/{realm}/triggers:
372:  /realmmanagement/v1/{realm}/triggers/{name}:
424:  /realmmanagement/v1/{realm}/config/auth:
493:  /realmmanagement/v1/{realm}/config/device_registration_limit:
524:  /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention:
559:  /realmmanagement/v1/{realm}/version:
592:  /realmmanagement/v1/{realm}/devices/{device}:
618:  /realmmanagement/v1/{realm}/policies:
688:  /realmmanagement/v1/{realm}/policies/{name}:

$ rg -n 'HandleFunc|Handle\(|\.(Get|Post|Put|Delete|Patch)\("' internal/housekeeping/*.go
internal/housekeeping/http.go:37:	mux.Handle("GET /housekeeping/v1/realms", h(a.listRealms))
internal/housekeeping/http.go:38:	mux.Handle("POST /housekeeping/v1/realms", h(a.createRealm))
internal/housekeeping/http.go:39:	mux.Handle("GET /housekeeping/v1/realms/{realm}", h(a.getRealm))
internal/housekeeping/http.go:40:	mux.Handle("PATCH /housekeeping/v1/realms/{realm}", h(a.patchRealm))
internal/housekeeping/http.go:41:	mux.Handle("DELETE /housekeeping/v1/realms/{realm}", h(a.deleteRealm))

$ rg -n '^\s{2}/' docs/api/astarte_housekeeping_api.yaml
26:  /housekeeping/v1/realms:
106:  /housekeeping/v1/realms/{realm}:

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

$ rg -n 'Handle\(' internal/appengine/*.go internal/appengine/*/*.go
internal/appengine/stream/ws.go:46:	mux.Handle("GET /astrate/v1/{realm}/socket", a.require(http.HandlerFunc(a.handle)))
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
internal/appengine/channels/ws.go:48:	mux.Handle("GET /appengine/v1/socket/websocket", http.HandlerFunc(a.handle))

All appengine routes are documented. Let me pick **realm-management** as the surface for this run and deep-check a few endpoints against the handlers.
→ Read internal/realm/http.go
Let me read the realm-management spec and check a few endpoints against the handlers.
→ Read docs/api/astarte_realm_management_api.yaml

$ rg -n 'APICompatVersion' internal/realm/service.go
580:// APICompatVersion is the upstream Realm Management API level Astrate
588:const APICompatVersion = "1.2.2"

The realm-management spec looks thorough. Let me check the auth middleware (to know which auth responses are possible) and the version example drift, plus a few handler specifics.
→ Read internal/auth/middleware.go
The getVersion example drift (1.1.0 vs 1.2.2) is already a blocked task in the queue (`docs-sync-rm-version-example`). Let me check the update path and how the appengine spec documents auth responses, to find fresh drift.

$ rg -n '403|Forbidden|404' docs/api/astarte_appengine_api.yaml | head -30
156:        "404":
201:        "404":
235:        "404":
285:        "404":
319:        "404":
362:        "404":
412:        "404":
460:        "404":
510:        "404":
548:        "404":
576:        "404":
616:        "404":
659:        "404":
702:        "404":
747:        "404":
780:        "404":
886:        "404":
932:        "404":
974:        "404":
1004:        "404":
1050:        "404":
1075:        "404":
1103:        "404":
1142:        "404":
1186:        "404":
1230:        "404":
1276:        "404":
1310:        "404":

$ rg -n 'UpdateInterface|func \(s \*Service\)' internal/realm/service.go | sed -n '1,60p'
100:func (s *Service) WithDisconnecter(d Disconnecter) *Service {
113:func (s *Service) DeleteDevice(ctx context.Context, realm, deviceID string) error {
139:func (s *Service) realmID(ctx context.Context, realm string) (int16, error) {
156:func (s *Service) InstallInterface(ctx context.Context, realm string, def []byte) (*store.StoredInterface, error) {
185:// UpdateInterface applies a minor upgrade, enforcing the additive-only
191:func (s *Service) UpdateInterface(ctx context.Context, realm, urlName string, urlMajor int, def []byte) (*store.StoredInterface, error) {
227:	si, err := s.st.UpdateInterface(ctx, r.ID, def)
238:func (s *Service) checkRetentionCeiling(r *store.Realm, iface *interfaceschema.Interface) error {
255:func (s *Service) DeleteInterface(ctx context.Context, realm, name string, major int) error {
274:func (s *Service) checkNameCollision(ctx context.Context, rid int16, name string) error {
309:func (s *Service) ListInterfaces(ctx context.Context, realm string) ([]string, error) {
335:func (s *Service) ListInterfacesDetailed(ctx context.Context, realm string) ([]json.RawMessage, error) {
364:func (s *Service) ListInterfaceMajors(ctx context.Context, realm, name string) ([]int, error) {
388:func (s *Service) GetInterface(ctx context.Context, realm, name string, major int) (json.RawMessage, error) {
419:func (s *Service) CreateTrigger(ctx context.Context, realm string, def []byte) (*store.Trigger, error) {
460:func (s *Service) GetTrigger(ctx context.Context, realm, name string) (json.RawMessage, error) {
473:func (s *Service) DeleteTrigger(ctx context.Context, realm, name string) error {
486:func (s *Service) ListTriggers(ctx context.Context, realm string) ([]string, error) {
515:func (s *Service) CreatePolicy(ctx context.Context, realm string, def []byte) (*store.TriggerPolicy, error) {
528:func (s *Service) GetPolicy(ctx context.Context, realm, name string) (json.RawMessage, error) {
541:func (s *Service) ListPolicies(ctx context.Context, realm string) ([]string, error) {
559:func (s *Service) DeletePolicy(ctx context.Context, realm, name string) error {
592:func (s *Service) GetDeviceRegistrationLimit(ctx context.Context, realm string) (*int32, error) {
605:func (s *Service) GetDatastreamMaximumStorageRetention(ctx context.Context, realm string) (int64, error) {
622:func (s *Service) GetAuthKey(ctx context.Context, realm string) (string, error) {
632:func (s *Service) SetAuthKey(ctx context.Context, realm, keyPEM string) error {
664:func (s *Service) interfacesChanged(ctx context.Context, rid int16, realm string) {
677:func (s *Service) triggersChanged(ctx context.Context, rid int16, realm string) {

→ Read internal/realm/service.go [offset=185, limit=90]
The RM surface has good coverage already. Let me check a few more details — the `InstallInterface` error path (violations envelope), and whether the appengine spec documents 403 (for comparison).

$ rg -n '"403"|Forbidden' docs/api/astarte_housekeeping_api.yaml docs/api/astarte_pairing_api.yaml docs/api/astarte_appengine_api.yaml
docs/api/astarte_housekeeping_api.yaml:45:        "403":
docs/api/astarte_housekeeping_api.yaml:46:          $ref: "#/components/responses/Forbidden"
docs/api/astarte_housekeeping_api.yaml:97:        "403":
docs/api/astarte_housekeeping_api.yaml:98:          $ref: "#/components/responses/Forbidden"
docs/api/astarte_housekeeping_api.yaml:133:        "403":
docs/api/astarte_housekeeping_api.yaml:134:          $ref: "#/components/responses/Forbidden"
docs/api/astarte_housekeeping_api.yaml:178:        "403":
docs/api/astarte_housekeeping_api.yaml:179:          $ref: "#/components/responses/Forbidden"
docs/api/astarte_housekeeping_api.yaml:201:        "403":
docs/api/astarte_housekeeping_api.yaml:202:          $ref: "#/components/responses/Forbidden"
docs/api/astarte_housekeeping_api.yaml:360:    Forbidden:
docs/api/astarte_housekeeping_api.yaml:361:      description: Forbidden — valid JWT but insufficient permissions.
docs/api/astarte_housekeeping_api.yaml:368:              detail: Forbidden
docs/api/astarte_pairing_api.yaml:138:        "403":
docs/api/astarte_pairing_api.yaml:139:          $ref: "#/components/responses/Forbidden"
docs/api/astarte_pairing_api.yaml:165:        "403":
docs/api/astarte_pairing_api.yaml:166:          $ref: "#/components/responses/Forbidden"
docs/api/astarte_pairing_api.yaml:215:        "403":
docs/api/astarte_pairing_api.yaml:216:          $ref: "#/components/responses/Forbidden"
docs/api/astarte_pairing_api.yaml:322:        "403":
docs/api/astarte_pairing_api.yaml:323:          $ref: "#/components/responses/Forbidden"
docs/api/astarte_pairing_api.yaml:568:    Forbidden:
docs/api/astarte_pairing_api.yaml:569:      description: Forbidden — valid credentials but insufficient permissions or device inhibited.
docs/api/astarte_pairing_api.yaml:576:              detail: Forbidden

→ Read internal/realm/service.go [limit=50, offset=156]

$ rg -n 'ViolationsError|violations' internal/realm/service.go; rg -n 'func ParseInterface|func ParseInterfaceCanonical|return.*ViolationsError|ViolationsError{' pkg/interfaceschema/parse.go | head
152:// a *interfaceschema.ViolationsError when the rejection has a probe-verified
rg: regex parse error:
    (?:func ParseInterface|func ParseInterfaceCanonical|return.*ViolationsError|ViolationsError{)
                                                                                                ^
error: repetition quantifier expects a valid decimal

Let me check the versions example (blocked task already exists for version example), and whether the parser emits `ViolationsError` (nested 422 envelope) on install.
$ rg -n 'createTrigger|trigger' internal/realm/http.go | head; rg -n 'WriteFieldErrors|WriteGateWay|error_name' docs/api/astarte_realm_management_api.yaml | head
11:	"github.com/astrate-platform/astrate/internal/engine/triggers"
18:// an interface definition (≤ 1024 mappings) or a trigger fits in 1 MiB.
44:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers", h(a.listTriggers))
45:	mux.Handle("POST /realmmanagement/v1/{realm}/triggers", h(a.createTrigger))
46:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers/{name}", h(a.getTrigger))
47:	mux.Handle("DELETE /realmmanagement/v1/{realm}/triggers/{name}", h(a.deleteTrigger))
60:// --- trigger delivery policies ------------------------------------------------
232:// --- triggers ---------------------------------------------------------------
243:func (a *API) createTrigger(w http.ResponseWriter, r *http.Request) {
251:		// Trigger validation errors are produced BY triggers and render as

$ rg -n 'ViolationsError' pkg/interfaceschema/parse.go | head
158:// probe-verified (issue #61) additionally carry a *ViolationsError in the
184:		var ve *ViolationsError
205:// probe-verified (issue #61) accumulate in a ViolationsError and abort the
206:// parse as a combined ErrInvalid + *ViolationsError value once every mapping
210:	vc := &ViolationsError{}
351:func buildMapping(raw *mappingJSON, ifaceType InterfaceType, index int, vc *ViolationsError) (Mapping, []endpointSegment, error) {

→ Read internal/realm/service.go [limit=30, offset=140]

$ rg -n 'errors.New|fmt.Errorf' pkg/interfaceschema/parse.go | head -5; rg -n 'ErrInvalid|ErrNameCollision|ErrNameMismatch|ErrMajorMismatch' pkg/interfaceschema/*.go | head; rg -n 'interface definition is invalid' -r . --glob '*.go'
38:var ErrInvalid = errors.New("invalid interface")
176:		return nil, nil, fmt.Errorf("%w: %w", ErrInvalid, err)
179:		return nil, nil, fmt.Errorf("%w: trailing data after interface document", ErrInvalid)
188:		return nil, nil, fmt.Errorf("%w: %w", ErrInvalid, err)
197:		return nil, nil, fmt.Errorf("%w: %w", ErrInvalid, err)
pkg/interfaceschema/violations.go:18:// (issue #61). ParseInterface returns it wrapped together with ErrInvalid,
pkg/interfaceschema/violations.go:19:// so errors.Is(err, ErrInvalid) and errors.As(err, &ve) both hit on the
pkg/interfaceschema/violations.go:28:// violationError couples the structured violations with the ErrInvalid
pkg/interfaceschema/violations.go:30:// own message (once), while Unwrap exposes both ErrInvalid and the
pkg/interfaceschema/violations.go:38:func (e *violationError) Unwrap() []error { return []error{ErrInvalid, e.ve} }
pkg/interfaceschema/parse.go:36:// ErrInvalid is wrapped by every ParseInterface failure, so callers can
pkg/interfaceschema/parse.go:38:var ErrInvalid = errors.New("invalid interface")
pkg/interfaceschema/parse.go:157:// wrapping ErrInvalid. Rejections for the rules whose upstream wire shape was
pkg/interfaceschema/parse.go:176:		return nil, nil, fmt.Errorf("%w: %w", ErrInvalid, err)
pkg/interfaceschema/parse.go:179:		return nil, nil, fmt.Errorf("%w: trailing data after interface document", ErrInvalid)

The RM spec is missing some responses. Let me verify the errors the interface handlers actually emit on the wire and confirm the version example status.
$ rg -n 'validationDetail|WriteError\(' internal/realm/*_test.go | head; rg -n 'ValidationError|errors.detail' docs/api/astarte_realm_management_api.yaml | head
134:          $ref: "#/components/responses/ValidationError"
232:          $ref: "#/components/responses/ValidationError"
489:          $ref: "#/components/responses/ValidationError"
684:          $ref: "#/components/responses/ValidationError"
736:          $ref: "#/components/responses/ValidationError"
1153:    ValidationError:

Let me verify the .mule directory contents (evidence/task files) and confirm the missing responses precisely.
$ ls -la .mule/
total 6728
drwxr-xr-x  7 root root    4096 Sep 23 21:52 .
drwxr-xr-x 17 root root    4096 Sep 23 19:04 ..
-rw-r--r--  1 root root      14 Sep 23 21:57 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 23 21:52 .heartbeat
-rw-r--r--  1 root root   20480 Sep 23 22:00 .last-output
-rw-r--r--  1 root root       0 Sep 23 21:57 .lock
-rw-r--r--  1 root root      10 Sep 23 21:57 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 6628054 Sep 23 22:00 cron.log
drwxr-xr-x  2 root root    4096 Sep 23 21:33 failed
-rw-r--r--  1 root root   56560 Sep 23 19:47 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    9017 Sep 23 21:52 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 23 20:39 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   79064 Sep 23 21:52 todo.md

$ ls .mule/tasks/ 2>/dev/null; ls .mule/reports/ 2>/dev/null; test -f .mule/for-giulio.md && echo "for-giulio exists"
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md
hygiene.md
log
for-giulio exists

$ rg -n 'ErrValidation =' internal/realm/*.go; rg -n 'validation failed' internal/realm/*.go | head
internal/realm/service.go:29:var ErrValidation = errors.New("realm: validation failed")
internal/realm/http.go:564:	const prefix = "realm: validation failed: "
internal/realm/service.go:29:var ErrValidation = errors.New("realm: validation failed")

Let me verify the exact wire messages for the interfaces 422 (nested envelope), the version example status, and the ValidationError example detail.
$ rg -n 'const validationPrefix|validationDetail' internal/realm/http.go internal/realm/*_test.go 2>/dev/null | head; rg -n 'WriteViolations|writeViolations' internal/realm/*_test.go 2>/dev/null | head
internal/realm/http.go:393:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
internal/realm/http.go:411:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
internal/realm/http.go:560:// validationDetail strips the ErrValidation sentinel prefix so the response
internal/realm/http.go:562:func validationDetail(err error) string {

$ rg -n 'interface definition is invalid|invalid interface' pkg/interfaceschema/parse.go docs/api/astarte_realm_management_api.yaml | head
pkg/interfaceschema/parse.go:38:var ErrInvalid = errors.New("invalid interface")
pkg/interfaceschema/parse.go:256:		return nil, fmt.Errorf("invalid interface name %q", name)
docs/api/astarte_realm_management_api.yaml:1161:              detail: "realm: validation failed: interface definition is invalid"

$ rg -n '403|Forbidden' docs/api/astarte_realm_management_api.yaml || echo "NO 403 IN RM SPEC"
257:        "403":
258:          description: Forbidden — major version not 0 or interface in use.

Let me confirm the RM spec has no 403 anywhere and verify the update path errors reach it, then check the appengine spec for comparison.
$ rg -n '403|Forbidden' docs/api/astarte_appengine_api.yaml || echo "NO 403 IN AE SPEC"
NO 403 IN AE SPEC

$ rg -n 'ErrNameMismatch|ErrMajorMismatch|ErrMinorNotIncreased|ErrDowngradeNotAllowed|ErrMissingEndpoints|ErrIncompatibleEndpointChange|ErrAlreadyExists' internal/realm/service.go
37:	// ErrNameMismatch is a PUT whose body interface_name disagrees with the
39:	ErrNameMismatch = errors.New("realm: interface name does not match")
40:	// ErrMajorMismatch is a PUT whose body version_major disagrees with the
42:	ErrMajorMismatch = errors.New("realm: interface major does not match")
151:// store.ErrAlreadyExists; a schema violation yields ErrValidation — carrying
170:	// install keeps its store.ErrAlreadyExists shape.
204:		return nil, fmt.Errorf("%w: body declares %q, URL says %q", ErrNameMismatch, next.Name, urlName)
207:		return nil, fmt.Errorf("%w: body declares %d, URL says %d", ErrMajorMismatch, next.Major, urlMajor)
273:// store.ErrAlreadyExists path.
416:// and installs it. A duplicate name yields store.ErrAlreadyExists.
514:// /policies). Duplicates yield store.ErrAlreadyExists.

RM spec looks mostly solid. Let me verify a few specific claims — the 409-on-update path and what a flat 422 install error actually displays.
$ rg -n 'validationDetail|interface definition|invalid interface|errors.detail|detail' internal/realm/http_test.go 2>/dev/null | head -20; ls internal/realm/*_test.go
164:		// answers with its measured 409 detail (#62).
165:		want := `{"errors":{"detail":"Interface update contains incompatible endpoint changes"}}`
186:		want := `{"errors":{"detail":"Interface can't be deleted"}}`
203:		wantInUse := `{"errors":{"detail":"Interface can't be deleted since it's currently used"}}`
332:// compares the exact {"errors":{"detail":...}} body, interleaved with success
338:		bodyDup           = `{"errors":{"detail":"Interface already exists"}}`
339:		bodyCollision     = `{"errors":{"detail":"Interface name collision detected. Make sure that the difference between two interface names is not limited to the casing or the presence of hyphens."}}`
340:		bodyNameMismatch  = `{"errors":{"detail":"Interface name doesn't match the one in the interface json"}}`
341:		bodyMajorMismatch = `{"errors":{"detail":"Interface major version doesn't match the one in the interface json"}}`
342:		bodyMinorSame     = `{"errors":{"detail":"Interface minor version was not increased"}}`
343:		bodyDowngrade     = `{"errors":{"detail":"Interface downgrade not allowed"}}`
344:		bodyMutated       = `{"errors":{"detail":"Interface update contains incompatible endpoint changes"}}`
345:		bodyDropped       = `{"errors":{"detail":"Interface update has missing endpoints"}}`
346:		bodyMajorNotFound = `{"errors":{"detail":"Interface major not found"}}`
347:		bodyNotFound      = `{"errors":{"detail":"Interface not found"}}`
348:		bodyCantDelete    = `{"errors":{"detail":"Interface can't be deleted"}}`
349:		bodyInUse         = `{"errors":{"detail":"Interface can't be deleted since it's currently used"}}`
350:		bodyTriggerDup    = `{"errors":{"detail":"Already exists"}}`
454:		{"duplicate_trigger_keeps_generic_detail", http.MethodPost, "/triggers", triggerJSON,
575:			t.Errorf("amqp trigger detail %q does not mention the rejection", body)
internal/realm/async_operation_test.go
internal/realm/dashboard_compat_test.go
internal/realm/detailed_required_encrypted_test.go
internal/realm/detailed_test.go
internal/realm/helpers_test.go
internal/realm/http_test.go
internal/realm/policy_test.go

→ Read internal/realm/http_test.go [limit=120, offset=330]

$ rg -n '422|errors|mappings|violation' internal/realm/http_test.go | sed -n '1,40p'
34:	ifaceV1  = `{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`
35:	ifaceV1b = `{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":1,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"},{"endpoint":"/count","type":"integer"}]}`
36:	ifaceV1x = `{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":2,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"integer"}]}`
37:	ifaceV2  = `{"interface_name":"com.ex.M7a.Sensors","version_major":2,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"string"}]}`
38:	draftV0  = `{"interface_name":"com.ex.M7a.Draft","version_major":0,"version_minor":1,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/x","type":"double"}]}`
165:		want := `{"errors":{"detail":"Interface update contains incompatible endpoint changes"}}`
186:		want := `{"errors":{"detail":"Interface can't be deleted"}}`
203:		wantInUse := `{"errors":{"detail":"Interface can't be deleted since it's currently used"}}`
225:			t.Errorf("invalid trigger: got %d, want 422", rec.Code)
270:// datastream_maximum_storage_retention answer upstream's named 422 envelope
279:			`"mappings":[{"endpoint":"/value","type":"double","database_retention_policy":"use_ttl","database_retention_ttl":%d}]}`,
282:	const overEnvelope = `{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`
301:	// Install above the ceiling → upstream's named 422, nothing installed.
304:		t.Errorf("install ttl 7200 under a 3600 ceiling: got %d %s, want 422 %s",
320:		t.Errorf("minor update to ttl 7200: got %d %s, want 422 %s",
332:// compares the exact {"errors":{"detail":...}} body, interleaved with success
338:		bodyDup           = `{"errors":{"detail":"Interface already exists"}}`
339:		bodyCollision     = `{"errors":{"detail":"Interface name collision detected. Make sure that the difference between two interface names is not limited to the casing or the presence of hyphens."}}`
340:		bodyNameMismatch  = `{"errors":{"detail":"Interface name doesn't match the one in the interface json"}}`
341:		bodyMajorMismatch = `{"errors":{"detail":"Interface major version doesn't match the one in the interface json"}}`
342:		bodyMinorSame     = `{"errors":{"detail":"Interface minor version was not increased"}}`
343:		bodyDowngrade     = `{"errors":{"detail":"Interface downgrade not allowed"}}`
344:		bodyMutated       = `{"errors":{"detail":"Interface update contains incompatible endpoint changes"}}`
345:		bodyDropped       = `{"errors":{"detail":"Interface update has missing endpoints"}}`
346:		bodyMajorNotFound = `{"errors":{"detail":"Interface major not found"}}`
347:		bodyNotFound      = `{"errors":{"detail":"Interface not found"}}`
348:		bodyCantDelete    = `{"errors":{"detail":"Interface can't be deleted"}}`
349:		bodyInUse         = `{"errors":{"detail":"Interface can't be deleted since it's currently used"}}`
350:		bodyTriggerDup    = `{"errors":{"detail":"Already exists"}}`
353:	sensors := func(minor int, mappings string) string {
355:			`"type":"datastream","ownership":"device","mappings":[%s]}`, minor, mappings)
359:			`"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`, name)
421:				`"ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`, rmIface),
441:				`"ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`, rmIface),
474:// TestRealmManagementViolationEnvelopes pins the probe-frozen 422 bodies of
476:// aligned mappings array.
483:			`"mappings":[{"endpoint":"/value","type":"double"}]}`, strings.Repeat("x", 1001))
484:		want := `{"errors":{"description":["should be at most 1000 character(s)"]}}`
487:			t.Errorf("got %d %s, want 422 %s", rec.Code, rec.Body, want)
493:			`"type":"datastream","ownership":"device","mappings":[`+

Good coverage already. Let me verify the violations-envelope on install and the actual flat-validation wire message.
$ rg -n 'Invalid|invalid|realm: validation' internal/realm/http_test.go | head
225:			t.Errorf("invalid trigger: got %d, want 422", rec.Code)
662:			wantBody:   `{"errors":{"action":{"http_method":["is invalid"]}}}`,

→ Read internal/realm/service.go [limit=75, offset=388]

$ rg -n 'interface definition is invalid|validation failed' docs/api/astarte_realm_management_api.yaml
1161:              detail: "realm: validation failed: interface definition is invalid"

$ rg -n 'ValidationError|error_name|maximum_database_retention|WriteFieldErrors' internal/realm/*_test.go | head
internal/realm/http_test.go:282:	const overEnvelope = `{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`

→ Read internal/realm/http_test.go [limit=80, offset=440]
The RM surface has good coverage. I've found a few real drifts. Let me verify the last pieces and check which branch/todo the append goes to.

$ git branch --show-current; git status --short | head
mule/queue

→ Read internal/realm/http_test.go [limit=45, offset=260]

$ rg -n 'a_rma|\*|Authorizes|func .*Authorizes' internal/auth/token.go | head -30
rg: internal/auth/token.go: IO error for operation on internal/auth/token.go: No such file or directory (os error 2)

Solid drift found. Let me verify one more candidate — whether a token-scoped 403 is genuinely reachable on RM routes — before finalizing proposals.
$ rg -n '"1\.1\.0"|1\.2\.2' docs/api/astarte_realm_management_api.yaml
584:                data: "1.1.0"

$ rg -l 'func.*Authorizes' internal/auth/; rg -n 'Authorizes|a_rma' internal/auth/*.go | head
internal/auth/claims_test.go
internal/auth/claims.go
internal/auth/cache_test.go:64:	if !tokA.Authorizes(ClaimPairing, "POST", "x") || tokA.Authorizes(ClaimAppEngine, "GET", "x") {
internal/auth/cache_test.go:67:	if !tokB.Authorizes(ClaimAppEngine, "GET", "x") || tokB.Authorizes(ClaimPairing, "POST", "x") {
internal/auth/cache_test.go:149:		if !tok.Authorizes(ClaimPairing, "POST", "x") {
internal/auth/jwt_test.go:88:			if !tok.Authorizes(ClaimPairing, "POST", "agent/devices") {
internal/auth/jwt_test.go:186:		"a_rma": []string{".*::.*"},
internal/auth/claims.go:32:	// ClaimRealmManagement authorizes the Realm Management API (a_rma).
internal/auth/claims.go:33:	ClaimRealmManagement Claim = "a_rma"
internal/auth/claims.go:75:	RealmManagement stringList `json:"a_rma,omitempty"`
internal/auth/claims.go:81:// AuthorizesChannel.
internal/auth/claims.go:159:// Authorizes reports whether the token grants `claim` for the given verb

→ Read internal/auth/claims.go [limit=70, offset=150]
The version example (1.1.0) is already a blocked task in the queue — skip it. Let me verify two remaining claims: the `ValidationError` example text, and the exact suite of 409 codes on PUT /interfaces.

$ rg -rn 'interface definition is invalid' --glob '*.go' --glob '*.yaml' . | head; echo ---; rg -n 'validation failed' internal/realm/service.go pkg/interfaceschema/parse.go | head
./docs/api/astarte_realm_management_api.yaml:              detail: "realm: validation failed: n"
---
internal/realm/service.go:29:var ErrValidation = errors.New("realm: validation failed")

$ rg -n 'StatusConflict|http.StatusConflict|writeInterfaceError' internal/realm/http_test.go | head; echo ---; rg -n 'discard detection|information' internal/realm/http_test.go | head -3
136:		if rec := r.req(t, http.MethodPost, "/interfaces", ifaceV1, r.rmaToken); rec.Code != http.StatusConflict {
167:		if rec.Code != http.StatusConflict || rec.Body.String() != want {
221:		if rec := r.req(t, http.MethodPost, "/triggers", triggerJSON, r.rmaToken); rec.Code != http.StatusConflict {
403:			nil, http.StatusConflict, bodyDup},
406:			nil, http.StatusConflict, bodyCollision},
413:			nil, http.StatusConflict, bodyCollision},
418:			named("com.ex.Mismatched"), nil, http.StatusConflict, bodyNameMismatch},
422:			nil, http.StatusConflict, bodyMajorMismatch},
425:			nil, http.StatusConflict, bodyMinorSame},
429:			nil, http.StatusConflict, bodyDowngrade},
---

Let me verify what an actual malformed-interface install emits on the wire (to ground the ValidationError example fix), and check the `422` mapping for generic invalid bodies.
$ rg -n 'StatusUnprocessableEntity|422' internal/realm/http_test.go | head; echo ---; rg -n 'invalid interface|does not parse|ErrInvalid' internal/realm/service.go
224:		if rec := r.req(t, http.MethodPost, "/triggers", `{"name":"bad","action":{},"simple_triggers":[]}`, r.rmaToken); rec.Code != http.StatusUnprocessableEntity {
225:			t.Errorf("invalid trigger: got %d, want 422", rec.Code)
270:// datastream_maximum_storage_retention answer upstream's named 422 envelope
301:	// Install above the ceiling → upstream's named 422, nothing installed.
303:	if rec.Code != http.StatusUnprocessableEntity || rec.Body.String() != overEnvelope {
304:		t.Errorf("install ttl 7200 under a 3600 ceiling: got %d %s, want 422 %s",
319:	if rec.Code != http.StatusUnprocessableEntity || rec.Body.String() != overEnvelope {
320:		t.Errorf("minor update to ttl 7200: got %d %s, want 422 %s",
474:// TestRealmManagementViolationEnvelopes pins the probe-frozen 422 bodies of
486:		if rec.Code != http.StatusUnprocessableEntity || rec.Body.String() != want {
---
221:		return nil, fmt.Errorf("realm: stored interface %s v%d does not parse: %w", next.Name, next.Major, err)
354:			return nil, fmt.Errorf("realm: stored interface %s v%d does not parse: %w", si.Name, si.Major, err)
426:		return nil, fmt.Errorf("%w: trigger does not parse: %v", ErrValidation, err)

$ rg -n 'ParseInterface|parse' internal/realm/service.go | head
161:	iface, canon, err := interfaceschema.ParseInterfaceCanonical(def)
189:// parsed body must agree with both (upstream 409s the mismatches, #62), and
196:	next, canon, err := interfaceschema.ParseInterfaceCanonical(def)
219:	prev, err := interfaceschema.ParseInterface(stored.Definition)
221:		return nil, fmt.Errorf("realm: stored interface %s v%d does not parse: %w", next.Name, next.Major, err)
354:			return nil, fmt.Errorf("realm: stored interface %s v%d does not parse: %w", si.Name, si.Major, err)
426:		return nil, fmt.Errorf("%w: trigger does not parse: %v", ErrValidation, err)
643:// parses back into individual keys.

The RM surface is mostly well-covered (many docs-sync tasks already landed). I've verified four genuine drifts. Let me record them in the todo. First, let me read the current tail of the file to append cleanly.
→ Read .mule/todo.md [offset=209, limit=6]
I've verified the realm-management surface thoroughly. The route list is fully covered (a lot of docs-sync work has already landed), but I found four genuine drifts. Let me append these as task lines.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -211,4 +211,9 @@
 - [x] flow-msg-json-integer-precision [auto]: `setDataFromWire`'s TypeInteger float64 branch does `int64(v)` on a value that plain `json.Unmarshal` (message.go:141-145) already rounded to float64, so integer wire values > 2^53 silently lose precision and wrap (`float64(MaxInt64)` → 2^63 → MinInt64 on decode), a fractional `3.7` silently truncates to 3, and `1e300` becomes garbage — the `json.Number` branch (message.go:268-273) that would reject these is unreachable in this path. Mirror the accepted `payload-longinteger-fraction-quantize` rule: decode with `json.Decoder`+`UseNumber` (parse via `json.Number.Int64`) or reject non-integral/out-of-int64 float64; add round-trip rows for `math.MaxInt64`, `1e300`, `3.7` to message_test.go. Container-free.
 - [x] flow-sort-bounded-buffer [auto]: `Sort` (internal/flow/blocks/sort.go:32-65) appends every message to an uncapped `buf` and only flushes while `len(buf) > 1 && buf[0].ts <= newest-windowUs` (sort.go:60) — a stream whose timestamps stay within `window_ms` of the newest (clustered bursts, or messages sharing one timestamp) never satisfies the condition, so the block emits nothing while the slice grows without bound, and the documented "newest is never emitted" mean the final message is held indefinitely; secondarily `window_ms * 1000` (sort.go:31) overflows for large configs and silently inverts the flush edge. Add a configurable max buffered count with a documented overflow policy (drop-oldest / error / force-flush) and a cap test feeding N same-timestamp messages; guard the multiply. Container-free.
 - [!] flow-randomsource-span-overflow [auto]: `randomSource.next()` computes `rand.Int64N(s.maxInt-s.minInt+1)` (internal/flow/blocks/randomsource.go:171) whose argument overflows int64 to a non-positive value when min/max span more than the int64 range — config like `min:-9e18, max:9e18` passes the single `min <= max` check (randomsource.go:106-108) but makes `rand.Int64N` panic (it panics on n <= 0) inside the source pump goroutine (flow.go:205-239, no recover), crashing the whole process. Reject a span wider than `math.MaxInt64` at construction alongside the existing check, and add a construct/Emit test with the wide span. Container-free. — BLOCKED: lint failed: internal/flow/blocks/randomsource.go:113:40: G115: integer overflow conversion int64 -> uint64 (gosec)
 - [x] flow-filter-key-contains-test [auto]: `Filter`'s `key_contains` rule (`strings.Contains`, internal/flow/blocks/transform.go:69) has zero test coverage — `transform_test.go` only exercises `key_prefix`, `type`, unknown-type and the at-least-one construct rule, so a regression silently dropping the Contains check would pass the whole suite. Add substring-match, substring-absent, and prefix+contains combined rows to the filter tests. Container-free.
+
+- [ ] docs-sync-rm-put-interface-409 [auto]: add the missing `409` Conflict response to `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` in docs/api/astarte_realm_management_api.yaml (yaml:222-234 lists only 204/400/401/404/422/500) — `updateInterface` answers 409 for every url/body disagreement and minor-upgrade incompatibility: `ErrNameMismatch` (service.go:204 → "Interface name doesn't match the one in the interface json", http.go:367-369), `ErrMajorMismatch` (service.go:207 → http.go:370-372), and the `CheckMinorUpgrade` sentinels `ErrMinorNotIncreased`/`ErrDowngradeNotAllowed`/`ErrMissingEndpoints`/`ErrIncompatibleEndpointChange` (service.go:223-226 → http.go:373-384) plus `store.ErrAlreadyExists` (http.go:385-386) — the exact set `TestRealmManagementErrorCodes` pins as StatusConflict (http_test.go:417-437). Reuse the `Conflict` response component (already used on POST, yaml:131-132) and add name-mismatch/major-mismatch examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` (yaml:231-232) in docs/api/astarte_realm_management_api.yaml as a oneOf — both handlers go through `writeInterfaceError`, which answers three distinct 422 bodies: the flat ErrorDetail (`ErrValidation` via `validationDetail`, http.go:392-393), the nested violations changeset envelope (`writeViolations`, http.go:358-359 + 434-533, e.g. `{"errors":{"description":["should be at most 1000 character(s)"]}}` and the aligned full-length `mappings` array, http_test.go:480-501), and the named FieldErrors envelope for `ErrMaximumDatabaseRetentionExceeded` (`{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, http.go:355-357, http_test.go:282-320); today only the flat ValidationError is referenced. Mirror the createTrigger oneOf pattern (yaml:344-368). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "realm: validation failed: interface definition is invalid"`, but `validationDetail` strips the `realm: validation failed: ` prefix on the wire (http.go:562-568) and the parser's real messages are `invalid interface: ...` (pkg/interfaceschema/parse.go:38, wrapped at service.go:163); the example should show the stripped, real message. Same class as the already-fixed docs-sync-hk-validation-example (housekeeping http.go:251-258). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

Done: ran the docs-sync recipe on the realm-management surface, found four concrete drifts in `docs/api/astarte_realm_management_api.yaml`, and appended them as `[auto]` task lines to `.mule/todo.md`.

The comparisons:
- Code routes (`internal/realm/http.go:38-57`) ↔ documented paths (yaml:27-739) match 1:1 — no missing/present-without-doc route drift.
- Deep-checked several endpoints against handlers (`updateInterface`, `installInterface`, `writeInterfaceError`, `writeError`, auth middleware); the version example (`1.1.0` vs `APICompatVersion "1.2.2"`) is already queued as the blocked `docs-sync-rm-version-example`, so I did not duplicate it.

Proposals added (4, under the max of 5):
1. `docs-sync-rm-put-interface-409` — PUT /interfaces/{name}/{major} missing 409 Conflict (name/major mismatch, minor-upgrade sentinels, already-exists; all pinned in http_test.go:417-437).
2. `docs-sync-rm-interface-422-shapes` — install/update 422 documented only as flat ValidationError, but handlers also emit the nested violations envelope and the `error_name` retention envelope (mirror the createTrigger oneOf).
3. `docs-sync-rm-validation-example-prefix` — ValidationError example shows the `realm: validation failed: ` prefix that `validationDetail` (http.go:562-568) strips on the wire.
4. `docs-sync-rm-auth-403` — a_rma ops document only 401, but the middleware emits 403 for verified-but-unauthorized claims; housekeeping/pairing specs already document it.

Each line names the file, the path, what the code actually does and where, and the `make -C docs build` verification.

```
Done: realm-management docs sync — 4 drift proposals appended to .mule/todo.md
Files: .mule/todo.md
Verified: rg/lit comparisons of internal/realm/http.go, service.go, http_test.go, auth/middleware.go, claims.go vs astarte_realm_management_api.yaml -> pass (no build run; proposals only)
Unsure: nothing
Follow-ups: docs-sync-rm-put-interface-409, docs-sync-rm-interface-422-shapes, docs-sync-rm-validation-example-prefix, docs-sync-rm-auth-403
```
