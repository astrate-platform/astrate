slug: recipe-docs-sync
verdict: proposed
at:  f6dc460
ran: 2026-09-18T18:01:19Z on DietPi in 421s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls docs/api/ 2>/dev/null; echo ---; ls docs/site/ 2>/dev/null | head; echo ---; ls internal/
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

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*appengine*'
(no output)

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

$ ls internal/appengine/; echo ---; rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/appengine/
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
---
internal/appengine/http.go:209:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:230:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:521:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/stream/ws.go:152:	return r.Header.Get("Accept") == "text/event-stream" || r.URL.Query().Get("transport") == "sse"

$ rg -n 'Handle|\.(Get|Post|Put|Delete|Patch)\(' internal/appengine/http.go | head -60; echo ---; rg -n 'func.*handle|endpoint|registerRoutes|newRouter' internal/appengine/http.go
30:	require func(http.Handler) http.Handler
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
87:	limit, _ := strconv.Atoi(q.Get("limit"))
88:	details, _ := strconv.ParseBool(q.Get("details"))
89:	page, err := a.svc.ListDevices(r.Context(), r.PathValue("realm"), q.Get("from_token"), limit, details)
116:		details := q.Get("details")
120:		limit := q.Get("limit")
209:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
230:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
438:	details, _ := strconv.ParseBool(q.Get("details"))
439:	limit, _ := strconv.Atoi(q.Get("limit"))
441:		details, q.Get("from_token"), limit)
469:	limit := q.Get("limit")
474:	if d := q.Get("details"); d != "" {
521:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
564:		if v := q.Get(p.name); v != "" {
573:	if q.Get("since") != "" && q.Get("since_after") != "" {
576:	if v := q.Get("limit"); v != "" {
587:	if v := q.Get("downsample_to"); v != "" {
601:	switch v := q.Get("format"); v {
615:		switch v := q.Get(p.name); v {
624:	if v := q.Get("retrieve_metadata"); v != "" {
633:	opts.DownsampleKey = q.Get("downsample_key")
635:	opts.Descending = q.Get("sort") != "ascending"
---

$ rg -n '^    (get|post|put|delete|patch):' docs/api/astarte_appengine_api.yaml
27:    get:
99:    get:
137:    get:
161:    patch:
211:    get:
240:    patch:
293:    get:
325:    get:
370:    get:
419:    put:
469:    post:
519:    delete:
556:    get:
584:    get:
624:    get:
666:    put:
711:    post:
756:    delete:
790:    get:
815:    post:
864:    get:
892:    get:
939:    post:
982:    get:
1009:    patch:
1059:    delete:
1081:    get:
1109:    get:
1150:    get:
1193:    put:
1239:    post:
1285:    delete:

→ Read internal/appengine/http.go [offset=30, limit=55]
→ Read docs/api/astarte_appengine_api.yaml [offset=136, limit=75]
→ Read internal/appengine/http.go [offset=85, limit=130]
→ Read internal/appengine/http.go [limit=120, offset=214]
→ Read docs/api/astarte_appengine_api.yaml [limit=130, offset=660]
→ Read docs/api/astarte_appengine_api.yaml [limit=110, offset=26]

$ rg -n 'func \(a \*API\) (createGroup|addGroupDevice|removeGroupDevice|getGroup|listGroupDevices)' internal/appengine/http.go
398:func (a *API) createGroup(w http.ResponseWriter, r *http.Request) {
428:func (a *API) getGroup(w http.ResponseWriter, r *http.Request) {
436:func (a *API) listGroupDevices(w http.ResponseWriter, r *http.Request) {
486:func (a *API) addGroupDevice(w http.ResponseWriter, r *http.Request) {
499:func (a *API) removeGroupDevice(w http.ResponseWriter, r *http.Request) {
509:func (a *API) getGroupDevice(w http.ResponseWriter, r *http.Request) {

→ Read docs/api/astarte_appengine_api.yaml [limit=175, offset=863]
→ Read internal/appengine/http.go [limit=130, offset=398]
→ Read docs/api/astarte_appengine_api.yaml [limit=68, offset=789]
→ Read docs/api/astarte_appengine_api.yaml [limit=25, offset=1056]

$ rg -n 'DataEnvelopeDeviceStatus|DeviceStatus:' docs/api/astarte_appengine_api.yaml | head; echo ---; rg -n 'type DeviceStatus struct' internal/appengine/ -A 40
62:                  - $ref: "#/components/schemas/DataEnvelopeDeviceStatusList"
153:                $ref: "#/components/schemas/DataEnvelopeDeviceStatus"
196:                $ref: "#/components/schemas/DataEnvelopeDeviceStatus"
232:                $ref: "#/components/schemas/DataEnvelopeDeviceStatus"
280:                $ref: "#/components/schemas/DataEnvelopeDeviceStatus"
929:                  - $ref: "#/components/schemas/DataEnvelopeDeviceStatusList"
1001:                $ref: "#/components/schemas/DataEnvelopeDeviceStatus"
1045:                $ref: "#/components/schemas/DataEnvelopeDeviceStatus"
1444:    DeviceStatus:
1564:    DataEnvelopeDeviceStatus:
---
internal/appengine/service.go:105:type DeviceStatus struct {
internal/appengine/service.go-106-	ID                   string                        `json:"id"`
internal/appengine/service.go-107-	Connected            bool                          `json:"connected"`
internal/appengine/service.go-108-	Introspection        map[string]introspectionEntry `json:"introspection"`
internal/appengine/service.go-109-	Aliases              map[string]string             `json:"aliases"`
internal/appengine/service.go-110-	Attributes           map[string]string             `json:"attributes"`
internal/appengine/service.go-111-	Groups               []string                      `json:"groups"`
internal/appengine/service.go-112-	CredentialsInhibited bool                          `json:"credentials_inhibited"`
internal/appengine/service.go-113-	TotalReceivedMsgs    int64                         `json:"total_received_msgs"`
internal/appengine/service.go-114-	TotalReceivedBytes   int64                         `json:"total_received_bytes"`
internal/appengine/service.go-115-	FirstRegistration    *time.Time                    `json:"first_registration"`
internal/appengine/service.go-116-	FirstCredentialsReq  *time.Time                    `json:"first_credentials_request"`
internal/appengine/service.go-117-	LastConnection       *time.Time                    `json:"last_connection"`
internal/appengine/service.go-118-	LastDisconnection    *time.Time                    `json:"last_disconnection"`
internal/appengine/service.go-119-	LastSeenIP           string                        `json:"last_seen_ip,omitempty"`
internal/appengine/service.go-120-	PreviousInterfaces   map[string]introspectionEntry `json:"previous_interfaces,omitempty"`
internal/appengine/service.go-121-}
internal/appengine/service.go-122-
internal/appengine/service.go-123-// introspectionEntry renders one introspection pair (upstream shape).
internal/appengine/service.go-124-type introspectionEntry struct {
internal/appengine/service.go-125-	Major int `json:"major"`
internal/appengine/service.go-126-	Minor int `json:"minor"`
internal/appengine/service.go-127-}
internal/appengine/service.go-128-
internal/appengine/service.go-129-// DevicePage is one device-list page: IDs, or full statuses when details was
internal/appengine/service.go-130-// requested. Next is the next-page cursor ("" at the end).
internal/appengine/service.go-131-type DevicePage struct {
internal/appengine/service.go-132-	IDs      []string
internal/appengine/service.go-133-	Statuses []*DeviceStatus
internal/appengine/service.go-134-	Next     string
internal/appengine/service.go-135-}
internal/appengine/service.go-136-
internal/appengine/service.go-137-// ListDevices returns one page of devices (upstream GET /devices). after is
internal/appengine/service.go-138-// the cursor (exclusive); limit <= 0 selects DefaultDeviceLimit; details
internal/appengine/service.go-139-// switches the page from bare IDs to full status objects (the dashboard's
internal/appengine/service.go-140-// device list).
internal/appengine/service.go-141-func (s *Service) ListDevices(ctx context.Context, realm string, after string, limit int, details bool) (*DevicePage, error) {
internal/appengine/service.go-142-	rid, err := s.realmID(ctx, realm)
internal/appengine/service.go-143-	if err != nil {
internal/appengine/service.go-144-		return nil, err
internal/appengine/service.go-145-	}

→ Read docs/api/astarte_appengine_api.yaml [limit=125, offset=1444]
→ Read docs/api/astarte_appengine_api.yaml [limit=90, offset=584]
→ Read internal/appengine/http.go [limit=95, offset=545]

$ rg -n 'retrieve_metadata|downsample_key|allow_biginteger|allow_safe_biginteger|since_after' docs/api/astarte_appengine_api.yaml
1369:      name: since_after
1421:      name: allow_bigintegers
1431:      name: allow_safe_bigintegers

$ rg -rn 'retrieve_metadata|downsample_key' internal/ --glob '*.go' | head -30
internal/appengine/http.go:	if v := q.Get("n"); v != "" {
internal/appengine/http.go:			fe.addf("n", "is invalid")
internal/appengine/http.go:	opts.DownsampleKey = q.Get("n")
internal/appengine/query_opts_test.go:		{"non bool n", "n=maybe", FieldErrors{"n": {"is invalid"}}},
internal/appengine/query_opts_test.go:	r := httptest.NewRequest(http.MethodGet, "/x?format=table&n=true"+
internal/appengine/query_opts_test.go:		"&allow_bigintegers=false&allow_safe_bigintegers=true&n=abc&limit=5", nil)

That grep used `-r` as replace by mistake. Re-checking properly:
$ rg -n 'retrieve_metadata|downsample_key' internal/appengine/http.go internal/appengine/query_opts_test.go
internal/appengine/http.go:624:	if v := q.Get("retrieve_metadata"); v != "" {
internal/appengine/http.go:630:			fe.addf("retrieve_metadata", "is invalid")
internal/appengine/http.go:633:	opts.DownsampleKey = q.Get("downsample_key")
internal/appengine/query_opts_test.go:93:		{"non bool retrieve_metadata", "retrieve_metadata=maybe", FieldErrors{"retrieve_metadata": {"is invalid"}}},
internal/appengine/query_opts_test.go:124:	r := httptest.NewRequest(http.MethodGet, "/x?format=table&retrieve_metadata=true"+
internal/appengine/query_opts_test.go:125:		"&allow_bigintegers=false&allow_safe_bigintegers=true&downsample_key=abc&limit=5", nil)

→ Read internal/appengine/http.go [limit=55, offset=385]

$ ls .mule/ 2>/dev/null; echo ---; ls .mule/tasks/ 2>/dev/null | head; echo ---; ls .mule/evidence/ 2>/dev/null | head
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
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md
---

→ Read .mule/todo.md
→ Read .mule/tasks/docs-sync-appengine-data-output-params.md

$ rg -n 'DataFormat|DataAllowBigIntegers|DataAllowSafeBigIntegers|DataDownsample|DataSort|DataLimit|DataSince' docs/api/astarte_appengine_api.yaml
344:        - $ref: "#/components/parameters/DataSince"
345:        - $ref: "#/components/parameters/DataSinceAfter"
347:        - $ref: "#/components/parameters/DataLimit"
348:        - $ref: "#/components/parameters/DataDownsample"
349:        - $ref: "#/components/parameters/DataSort"
350:        - $ref: "#/components/parameters/DataFormat"
351:        - $ref: "#/components/parameters/DataAllowBigIntegers"
352:        - $ref: "#/components/parameters/DataAllowSafeBigIntegers"
394:        - $ref: "#/components/parameters/DataSince"
395:        - $ref: "#/components/parameters/DataSinceAfter"
397:        - $ref: "#/components/parameters/DataLimit"
398:        - $ref: "#/components/parameters/DataDownsample"
399:        - $ref: "#/components/parameters/DataSort"
400:        - $ref: "#/components/parameters/DataFormat"
401:        - $ref: "#/components/parameters/DataAllowBigIntegers"
402:        - $ref: "#/components/parameters/DataAllowSafeBigIntegers"
598:        - $ref: "#/components/parameters/DataSince"
599:        - $ref: "#/components/parameters/DataSinceAfter"
601:        - $ref: "#/components/parameters/DataLimit"
602:        - $ref: "#/components/parameters/DataDownsample"
603:        - $ref: "#/components/parameters/DataSort"
604:        - $ref: "#/components/parameters/DataFormat"
605:        - $ref: "#/components/parameters/DataAllowBigIntegers"
606:        - $ref: "#/components/parameters/DataAllowSafeBigIntegers"
641:        - $ref: "#/components/parameters/DataSince"
642:        - $ref: "#/components/parameters/DataSinceAfter"
644:        - $ref: "#/components/parameters/DataLimit"
645:        - $ref: "#/components/parameters/DataDownsample"
646:        - $ref: "#/components/parameters/DataSort"
647:        - $ref: "#/components/parameters/DataFormat"
648:        - $ref: "#/components/parameters/DataAllowBigIntegers"
649:        - $ref: "#/components/parameters/DataAllowSafeBigIntegers"
1124:        - $ref: "#/components/parameters/DataSince"
1125:        - $ref: "#/components/parameters/DataSinceAfter"
1127:        - $ref: "#/components/parameters/DataLimit"
1128:        - $ref: "#/components/parameters/DataDownsample"
1129:        - $ref: "#/components/parameters/DataSort"
1130:        - $ref: "#/components/parameters/DataFormat"
1131:        - $ref: "#/components/parameters/DataAllowBigIntegers"
1132:        - $ref: "#/components/parameters/DataAllowSafeBigIntegers"
1168:        - $ref: "#/components/parameters/DataSince"
1169:        - $ref: "#/components/parameters/DataSinceAfter"
1171:        - $ref: "#/components/parameters/DataLimit"
1172:        - $ref: "#/components/parameters/DataDownsample"
1173:        - $ref: "#/components/parameters/DataSort"
1174:        - $ref: "#/components/parameters/DataFormat"
1175:        - $ref: "#/components/parameters/DataAllowBigIntegers"
1176:        - $ref: "#/components/parameters/DataAllowSafeBigIntegers"
1360:    DataSince:
1368:    DataSinceAfter:
1384:    DataLimit:
1391:    DataDownsample:
1398:    DataSort:
1407:    DataFormat:
1420:    DataAllowBigIntegers:
1430:    DataAllowSafeBigIntegers:

Appengine surface checks out in sync (all prior docs-sync-appengine tasks landed and the two dead params are deliberately undocumented). Picking the **housekeeping** surface for this run as a fresh one.
$ rg -n 'Handle\(|mux\.' internal/housekeeping/http.go; echo ---; rg -n '^\s{2}/' docs/api/astarte_housekeeping_api.yaml
37:	mux.Handle("GET /housekeeping/v1/realms", h(a.listRealms))
38:	mux.Handle("POST /housekeeping/v1/realms", h(a.createRealm))
39:	mux.Handle("GET /housekeeping/v1/realms/{realm}", h(a.getRealm))
40:	mux.Handle("PATCH /housekeeping/v1/realms/{realm}", h(a.patchRealm))
41:	mux.Handle("DELETE /housekeeping/v1/realms/{realm}", h(a.deleteRealm))
---
26:  /housekeeping/v1/realms:
106:  /housekeeping/v1/realms/{realm}:

→ Read docs/api/astarte_housekeeping_api.yaml
→ Read internal/housekeeping/http.go
→ Read internal/housekeeping/service.go

$ rg -o '\bASTRATE_[A-Z_]+\b' internal/ -N | sort -u
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

Housekeeping is fully in sync (paths, status codes, and all three schemas match the handlers). Let me run the recipe's config-key escalation check before concluding.
$ rg -o '\bASTRATE_[A-Z_]+\b' docs/site/configuration-reference.md 2>/dev/null | sort -u
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

$ rg -n 'ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION|ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED|HOUSEKEEPING_DEFAULT|HOUSEKEEPING_REALM' internal/config/config.go
279:		"ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
280:		"HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
297:	if v, ok := os.LookupEnv("ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"); ok {

Two code config keys are missing from the doc. Verifying they are real env-var reads:
→ Read internal/config/config.go [limit=45, offset=265]

$ rg -ni 'housekeeping|realm_deletion|default_datastream' docs/site/configuration-reference.md | head -30; echo ---; rg -ni 'HOUSEKEEPING' docs/site/ | head
36:| `addr` | string | `":8080"` | `ASTRATE_HTTP_ADDR` | No | Bind address for the single REST listener (pairing, realm management, housekeeping, appengine, native endpoints). |
78:## `[housekeeping]` — Instance-admin keys
139:[housekeeping]
140:jwt_public_key_files = ["/etc/astrate/housekeeping.pub"]
---
docs/site/pairing-and-security.md:17:- Claims: `a_aea` (AppEngine), `a_rma` (Realm Management), `a_pa` (Pairing), `a_ha` (Housekeeping), `a_ch` (Channels/live stream).
docs/site/deployment.md:69:| `:8080` | HTTP/HTTPS | REST API: `/pairing/v1`, `/realmmanagement/v1`, `/housekeeping/v1`, `/appengine/v1`, `/astrate/v1/{health,readiness,metrics}` |
docs/site/troubleshooting.md:16:3. If you suspect CA compromise, re-key the realm CA via the Housekeeping API (delete and recreate the realm).
docs/site/troubleshooting.md:79:3. To recover: generate a new master key, delete and recreate affected realms via the Housekeeping API (this mints fresh CAs sealed under the new key).
docs/site/swagger.md:12:| Housekeeping | Realm lifecycle and cross-realm ops |
docs/site/housekeeping-api.md:1:# Housekeeping API
docs/site/housekeeping-api.md:3:The Housekeeping API manages realm lifecycle at the instance level. It is the top-level administrative surface.
docs/site/housekeeping-api.md:5:**Base path:** `/housekeeping/v1/`
docs/site/housekeeping-api.md:13:POST /housekeeping/v1/realms
docs/site/housekeeping-api.md:35:GET /housekeeping/v1/realms

→ Read docs/site/configuration-reference.md [limit=85, offset=70]
The `[housekeeping]` section documents only the two JWT keys — missing the two config keys the code actually reads. This is a docs/site (Giulio's-voice) contradiction, so it escalates to `.mule/for-giulio.md`. Checking the file's current shape:
→ Read .mule/for-giulio.md
Both surfaces I checked are in sync, and the one piece of drift is already escalated. No files to write, no task lines to append. Here is the run:

**Surface picked: housekeeping** (appengine was the 2026-09-07 run's surface and already has 15+ landed docs-sync tasks; verified it again anyway — still in sync).

**Housekeeping findings:**
- Paths: all 5 registered routes (internal/housekeeping/http.go:37-41) match the spec's 2 paths with all methods — no documented-but-absent, no present-but-undocumented.
- Status codes (verified handler vs spec): createRealm 201/400/409/422; getRealm 404; patchRealm 200/400/404/422 (FieldErrors `invalid_update_parameters`, `can't be blank`, `is invalid`); deleteRealm 405 (`ErrDeletionDisabled`), 422 (`connected_devices_present`), 404 — all match `writeError` (http.go:233-249).
- Response/request schemas: `Realm`/`RealmCreate`/`RealmPatch` exactly match `viewBody` (http.go:94-101) and the optional `device_registration_limit` / `datastream_maximum_storage_retention` pointer fields.
- Config-key check: the only hit — `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` and `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` read at internal/config/config.go:277-306 and missing from configuration-reference.md — is **already escalated** to `.mule/for-giulio.md` (line 47, 2026-09-12; also line 129). Not re-escalated per the standing instruction against re-deriving the same proposal.

An accurate spec is the good result here. The recipe's cap is five proposals; I have zero, so nothing is appended to `.mule/todo.md` and no issue is filed (this recipe never asks for `gh issue create`, so none is run).

```
Done: docs-sync housekeeping surface verified fully in sync — no drift, no proposals; the one config-key gap (two ASTRATE_HOUSEKEEPING_* env keys) was already escalated in for-giulio.md
Files: none (nothing to write; appengine surface re-verified in sync too)
Verified: route/method inventory (rg vs mux.Handle), status-code + schema read-throughs in internal/housekeeping/http.go:37-249 & service.go vs docs/api/astarte_housekeeping_api.yaml; config-key diff (rg 'ASTRATE_[A-Z_]+' internal/config vs configuration-reference.md) -> pass/no-new-drift
Unsure: nothing (405/401/403 come from the a_ha middleware, matched to the spec's documented list, not probe-verified)
Follow-ups: none
```
