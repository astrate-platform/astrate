slug: recipe-docs-sync
verdict: proposed
at:  9625e13
ran: 2026-09-21T20:34:50Z on DietPi in 393s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls docs/api/ 2>/dev/null; echo "---"; ls internal/ 2>/dev/null
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

$ ls docs/site/ 2>/dev/null; echo "---"; ls docs/ 2>/dev/null
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
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
UPSTREAM-EXPERIMENTAL.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui

$ ls internal/appengine internal/housekeeping internal/pairing internal/realm 2>/dev/null
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

$ rg -n '^\s{2}/' docs/api/astrate_native_api.yaml | head -60
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

$ rg -n '^\s{2}/' docs/api/astarte_appengine_api.yaml | head -80
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

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/appengine/ internal/flowapi/ | head -80
internal/appengine/stream/ws.go:152:	return r.Header.Get("Accept") == "text/event-stream" || r.URL.Query().Get("transport") == "sse"
internal/appengine/http.go:209:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:230:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:521:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {

$ rg -n '\.(GET|POST|PUT|DELETE|PATCH)\("|Handle|Path\(|Methods\(|route|Route' internal/appengine/http.go | head -80
30:	require func(http.Handler) http.Handler
38:// Mount registers the routes on mux (paths wire-frozen to upstream
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

$ rg -n 'HandleFunc|\.Get\(|\.Post\(|\.Put\(|\.Delete\(|\.Patch\(' internal/appengine/http.go | head -80
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

→ Read internal/appengine/http.go [limit=90]

$ rg -n -A 12 '^\s{2}/appengine' docs/api/astarte_appengine_api.yaml | rg '^\s*\d+[-:]\s*(  /| {2,}(get|post|put|patch|delete|description))' | rg '^\s*\d+[-:]\s{2}/|^\s*\d+[-:]\s+(get|post|put|patch|delete):'
26:  /appengine/v1/{realm}/devices:
27-    get:
98:  /appengine/v1/{realm}/stats/devices:
99-    get:
136:  /appengine/v1/{realm}/devices/{device}:
137-    get:
210:  /appengine/v1/{realm}/devices-by-alias/{alias}:
211-    get:
292:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces:
293-    get:
324:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
325-    get:
369:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
370-    get:
555:  /appengine/v1/{realm}/devices/{device}/interfaces:
556-    get:
583:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
584-    get:
623:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
624-    get:
789:  /appengine/v1/{realm}/groups:
790-    get:
863:  /appengine/v1/{realm}/groups/{group}:
864-    get:
891:  /appengine/v1/{realm}/groups/{group}/devices:
892-    get:
981:  /appengine/v1/{realm}/groups/{group}/devices/{device}:
982-    get:
1080:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces:
1081-    get:
1108:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}:
1109-    get:
1149:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}:
1150-    get:

$ python3 - <<'EOF'
import re
lines = open('docs/api/astarte_appengine_api.yaml').read().splitlines()
path = None
for i, ln in enumerate(lines):
    m = re.match(r'^  (/.*):\s*$', ln)
    if m and not ln.startswith('      /'):
        path = m.group(1)
        methods = []
    elif re.match(r'^    (get|post|put|patch|delete):', ln):
        methods.append(re.match(r'^    (get|post|put|patch|delete):', ln).group(1))
    elif re.match(r'^  /\w', ln) and not re.match(r'^  /', ln):
        pass
    if i+1 < len(lines):
        nxt = lines[i+1]
        if path and re.match(r'^  /.*:$', nxt) and path != re.match(r'^  (/.*):$', nxt).group(1):
            print(path, '->', ','.join(methods) or 'NONE')
            path = None
    if i == len(lines)-1 and path:
        print(path, '->', ','.join(methods) or 'NONE')
EOF
/appengine/v1/{realm}/devices -> get
/appengine/v1/{realm}/stats/devices -> get
/appengine/v1/{realm}/devices/{device} -> get,patch
/appengine/v1/{realm}/devices-by-alias/{alias} -> get,patch
/appengine/v1/{realm}/devices-by-alias/{alias}/interfaces -> get
/appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface} -> get
/appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path} -> get,put,post,delete
/appengine/v1/{realm}/devices/{device}/interfaces -> get
/appengine/v1/{realm}/devices/{device}/interfaces/{interface} -> get
/appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path} -> get,put,post,delete
/appengine/v1/{realm}/groups -> get,post
/appengine/v1/{realm}/groups/{group} -> get
/appengine/v1/{realm}/groups/{group}/devices -> get,post
/appengine/v1/{realm}/groups/{group}/devices/{device} -> get,patch,delete
/appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces -> get
/appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface} -> get
/appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path} -> get,put,post,delete

The appengine path/method surface is fully in sync (all 32 code routes match the yaml). Now let me check documented status codes/responses against handlers for a few endpoints.
→ Read internal/appengine/http.go [offset=91, limit=320]
→ Read internal/appengine/http.go [offset=411]
→ Read docs/api/astarte_appengine_api.yaml [limit=180]
→ Read docs/api/astarte_appengine_api.yaml [offset=555, limit=310]
✱ Grep "func (DecodeData|WriteData|WriteDataWithLinks|WriteDataWithMetadata)" in pkg · 4 matches
→ Read docs/api/astarte_appengine_api.yaml [limit=45, offset=173]
→ Read pkg/astarteapi/envelope.go [offset=120, limit=140]

$ rg -n 'retrieve_metadata|downsample_key|DataRetrieveMetadata|DataDownsampleKey' docs/api/astarte_appengine_api.yaml; echo "=== rc:"; echo $?
=== rc:
0

`DecodeData` unwraps a `{"data": ...}` request envelope, so the yaml's `data`-wrapped request bodies are consistent. Let me check query-param coverage and the group endpoints.
$ rg -n 'DataSince|DataLimit|DataDownsample|DataSort|DataFormat|DataAllow|DataRetrieve|DataDownsampleKey|DataFromToken' docs/api/astarte_appengine_api.yaml | head
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

→ Read docs/api/astarte_appengine_api.yaml [limit=250, offset=981]

$ rg -n 'DataEnvelopeDeviceStatus:|DeviceStatus:|DataEnvelopeObject:|DataEnvelopeStringArray:|GroupCreate:' docs/api/astarte_appengine_api.yaml
1444:    DeviceStatus:
1534:    GroupCreate:
1564:    DataEnvelopeDeviceStatus:
1593:    DataEnvelopeObject:
1601:    DataEnvelopeStringArray:

→ Read docs/api/astarte_appengine_api.yaml [limit=200, offset=1444]

$ rg -n 'struct.*Device|type Device|"first_credentials_request"|"last_seen_ip"|"previous_interfaces"|"total_received_msgs"' internal/appengine/*.go | head -40
internal/appengine/service.go:105:type DeviceStatus struct {
internal/appengine/service.go:113:	TotalReceivedMsgs    int64                         `json:"total_received_msgs"`
internal/appengine/service.go:116:	FirstCredentialsReq  *time.Time                    `json:"first_credentials_request"`
internal/appengine/service.go:131:type DevicePage struct {
internal/appengine/service.go:284:type DevicePatch struct {

→ Read internal/appengine/service.go [offset=100, limit=45]
→ Read docs/api/astarte_appengine_api.yaml [limit=92, offset=210]

$ rg -n 'ErrInvalidAlias|ErrAliasTagNotFound|ErrAliasAlreadyInUse' internal/appengine/service.go | head
34:	ErrInvalidAlias         = errors.New("Invalid alias")           //nolint:staticcheck // ST1005: detail strings are upstream astarte_appengine_api 1.2.2 wire text
35:	ErrAliasAlreadyInUse    = errors.New("Alias already in use")    //nolint:staticcheck // ST1005: upstream wire text
36:	ErrAliasTagNotFound     = errors.New("Alias tag not found")     //nolint:staticcheck // ST1005: upstream wire text
328:			return nil, ErrInvalidAlias
337:			return nil, ErrAliasAlreadyInUse
348:				return nil, ErrAliasTagNotFound

$ rg -n -A 14 '"400"|responses:' docs/api/astarte_appengine_api.yaml | rg 'operationId: (getDataByAlias|getDataAtPathByAlias|listInterfacesByAlias)' ; rg -n -B 3 -A 20 'devices-by-alias/\{alias\}/interfaces/\{interface\}/\{path\}:|operationId: getDataAtPathByAlias|operationId: putDataByAlias' docs/api/astarte_appengine_api.yaml | head -70
294-      operationId: listInterfacesByAlias
366-        "500":
367-          $ref: "#/components/responses/InternalServerError"
368-
369:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
370-    get:
371:      operationId: getDataAtPathByAlias
372-      summary: Get data at path by alias
373-      description: |
374-        Returns data for a specific path within a device interface, looked up
375-        by alias.
376-      tags: [AppEngine - Data]
377-      security:
378-        - a_aea: []
379-      parameters:
380-        - $ref: "#/components/parameters/RealmName"
381-        - name: alias
382-          in: path
383-          required: true
384-          description: The alias value to look up.
385-          schema:
386-            type: string
387-        - $ref: "#/components/parameters/InterfaceName"
388-        - name: path
389-          in: path
390-          required: true
391-          description: The data path (e.g. `/sensors/temperature/value`).
--
417-          $ref: "#/components/responses/InternalServerError"
418-
419-    put:
420:      operationId: putDataByAlias
421-      summary: Set data at path by alias
422-      description: |
423-        Sets the value at a specific path in a device interface (server-owned
424-        data), looked up by alias. Returns 200 on success.
425-      tags: [AppEngine - Data]
426-      security:
427-        - a_aea: []
428-      parameters:
429-        - $ref: "#/components/parameters/RealmName"
430-        - name: alias
431-          in: path
432-          required: true
433-          description: The alias value to look up.
434-          schema:
435-            type: string
436-        - $ref: "#/components/parameters/InterfaceName"
437-        - name: path
438-          in: path
439-          required: true
440-          description: The data path.

→ Read internal/appengine/service.go [limit=55, offset=315]
→ Read docs/api/astarte_appengine_api.yaml [limit=29, offset=863]

$ rg -n -A 20 'DataFormat:|DataDownsample:|DataSort:|DataAllowBigIntegers:|DataSinceAfter:|DataSince:' docs/api/astarte_appengine_api.yaml | head -120
1360:    DataSince:
1361-      name: since
1362-      in: query
1363-      schema:
1364-        type: string
1365-        format: date-time
1366-      description: Return data points since this timestamp (RFC 3339).
1367-
1368:    DataSinceAfter:
1369-      name: since_after
1370-      in: query
1371-      schema:
1372-        type: string
1373-        format: date-time
1374-      description: Return data points strictly after this timestamp (RFC 3339).
1375-
1376-    DataTo:
1377-      name: to
1378-      in: query
1379-      schema:
1380-        type: string
1381-        format: date-time
1382-      description: Return data points up to this timestamp (RFC 3339).
1383-
1384-    DataLimit:
1385-      name: limit
1386-      in: query
1387-      schema:
1388-        type: integer
--
1391:    DataDownsample:
1392-      name: downsample_to
1393-      in: query
1394-      schema:
1395-        type: integer
1396-      description: Downsample to this many data points (must be > 2).
1397-
1398:    DataSort:
1399-      name: sort
1400-      in: query
1401-      schema:
1402-        type: string
1403-        enum: [ascending, descending]
1404-        default: descending
1405-      description: Sort order for datastream results.
1406-
1407:    DataFormat:
1408-      name: format
1409-      in: query
1410-      schema:
1411-        type: string
1412-        enum: [structured, table, disjoint_tables]
1413-        default: structured
1414-      description: |
1415-        Output format for the returned data. `structured` (the default) returns
1416-        flat objects; `table` returns a `{data, metadata}` envelope with rows
1417-        aligned to a sorted `table_header`; `disjoint_tables` pivots to a map of
1418-        column name to `[value, timestamp]` pairs.
1419-
1420:    DataAllowBigIntegers:
1421-      name: allow_bigintegers
1422-      in: query
1423-      schema:
1424-        type: boolean
1425-        default: false
1426-      description: |
1427-        When false, big-integer (longinteger) values on object interfaces are
1428-        serialized as decimal strings so no JavaScript caller loses precision.
1429-
1430-    DataAllowSafeBigIntegers:
1431-      name: allow_safe_bigintegers
1432-      in: query
1433-      schema:
1434-        type: boolean
1435-        default: false
1436-      description: |
1437-        When true, big-integer (longinteger) values on object interfaces stay
1438-        numbers when they fit JavaScript's safe integer range (|v| ≤ 9007199254740991)
1439-        and fall back to decimal strings otherwise.
1440-

The appengine path/method surface is fully in sync; I found one solid param-coverage gap. Let me also run the config-key check the recipe allows.
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

$ rg -o '\bASTRATE_[A-Z_]+' -N internal/ | sed 's/.*://' | sort -u
ASTRATE_DATABASE_DSN
ASTRATE_ENGINE_SHARDS
ASTRATE_FLOW_CONFIG
ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
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
ASTRATE_TEST_DSN

$ rg -n 'ASTRATE_FLOW_CONFIG|ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION|ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED|ASTRATE_TEST_DSN' internal/ | head -20
internal/flow/blocks/info.go:116:		Config:  "image (required), config (object→ASTRATE_FLOW_CONFIG), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",
internal/flow/blocks/schema.go:19:	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as ASTRATE_FLOW_CONFIG"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
internal/store/store_test.go:26:// ASTRATE_TEST_DSN database), one migrated Store, and one sub-suite per
internal/store/store_test.go:216:// suites stay rerunnable against an ASTRATE_TEST_DSN database.
internal/flow/blocks/container/block.go:45:	Nested        map[string]any // opaque JSON for ASTRATE_FLOW_CONFIG
internal/flow/blocks/container/docker.go:19:	// FlowConfigJSON is passed as ASTRATE_FLOW_CONFIG.
internal/flow/blocks/container/docker.go:99:		args = append(args, "-e", "ASTRATE_FLOW_CONFIG="+spec.FlowConfigJSON)
internal/flow/blocks/container/docker.go:201:// encodeFlowConfigJSON marshals nested config for ASTRATE_FLOW_CONFIG.
internal/config/config.go:279:		"ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
internal/config/config.go:297:	if v, ok := os.LookupEnv("ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"); ok {
internal/config/config_test.go:392:		prefixed = "ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION"
internal/config/config_test.go:533:	const env = "ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"
internal/engine/engine_integration_test.go:85:// TimescaleDB container (or ASTRATE_TEST_DSN), one migrated store, one
internal/testutil/pg.go:31:	// `make up` + ASTRATE_TEST_DSN=...). This is the fast local iteration
internal/testutil/pg.go:33:	EnvTestDSN = "ASTRATE_TEST_DSN"

$ rg -l 'ASTRATE_FLOW_CONFIG|ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION|ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED' docs/ || echo "absent from docs/"
docs/handoff/flow-lifecycle-followups-design-2026-08-22.md
docs/handoff/flow-design-b-container-block-2026-07-29.md

→ Read internal/config/config.go [limit=40, offset=270]

$ rg -n 'ErrPathNotFound|ErrNotFound|ErrInterfaceNotFound' internal/appengine/service.go internal/engine/*.go | head -30
internal/engine/serverdata.go:22:	// ErrInterfaceNotFound: the realm has no installed interface by that name.
internal/engine/serverdata.go:23:	ErrInterfaceNotFound = errors.New("engine: interface not installed")
internal/engine/serverdata.go:30:	// ErrPathNotFound: the path resolves no endpoint mapping.
internal/engine/serverdata.go:31:	ErrPathNotFound = errors.New("engine: path matches no endpoint")
internal/engine/serverdata.go:115:		return nil, fmt.Errorf("%w: %s", ErrInterfaceNotFound, ifaceName)
internal/engine/serverdata.go:137:			return nil, fmt.Errorf("%w: %q is not an aggregation prefix of %s", ErrPathNotFound, path, ifaceName)
internal/engine/serverdata.go:145:			return nil, fmt.Errorf("%w: %q matches no endpoint of %s", ErrPathNotFound, path, ifaceName)
internal/engine/serverdata.go:233:		return fmt.Errorf("%w: %s", ErrInterfaceNotFound, ifaceName)
internal/engine/serverdata.go:243:		return fmt.Errorf("%w: %q matches no endpoint of %s", ErrPathNotFound, path, ifaceName)
internal/engine/serverdata.go:290:// deviceView resolves the view; store.ErrNotFound propagates (M7 maps it to
internal/engine/fakestore_test.go:276:	return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, id)
internal/engine/fakestore_test.go:289:	return fmt.Errorf("%w: device %s", store.ErrNotFound, id)
internal/engine/fakestore_test.go:373:		return nil, fmt.Errorf("%w: property %s", store.ErrNotFound, path)
internal/engine/fakestore_test.go:386:		return nil, fmt.Errorf("%w: series %s", store.ErrNotFound, path)
internal/engine/fakestore_test.go:412:	return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, id)
internal/engine/engine_e2e_test.go:317:		return errors.Is(err, store.ErrNotFound)
internal/engine/engine_e2e_test.go:413:		return errors.Is(errB, store.ErrNotFound) && errors.Is(errG, store.ErrNotFound)
internal/engine/devicedata_test.go:141:			iface: stream, path: "/value", value: `"x"`, wantErr: store.ErrNotFound},
internal/engine/devicedata_test.go:143:			iface: "com.example.Nope", path: "/value", value: `"x"`, wantErr: ErrInterfaceNotFound},
internal/engine/devicedata_test.go:147:			iface: stream, path: "/nope", value: `"x"`, wantErr: ErrPathNotFound},
internal/engine/data.go:225:		if errors.Is(err, store.ErrNotFound) {
internal/engine/data.go:426:		case errors.Is(err, store.ErrNotFound):
internal/engine/data.go:436:		case errors.Is(err, store.ErrNotFound):
internal/appengine/service.go:65:// ErrGroupNotFound (→ 404 "Group not found") and store.ErrNotFound.
internal/appengine/service.go:67:	return fmt.Errorf("%w: %w: group %q", ErrGroupNotFound, store.ErrNotFound, name)
internal/appengine/service.go:217:		return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
internal/appengine/service.go:298:		return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
internal/appengine/service.go:404:// registered-but-not-member (or unparseable/unknown) device is store.ErrNotFound,
internal/appengine/service.go:416:		return 0, deviceid.ID{}, nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
internal/appengine/service.go:423:		return 0, deviceid.ID{}, nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)

$ ls -la .mule/todo.md .mule/for-giulio.md 2>/dev/null; echo ---; tail -25 .mule/todo.md 2>/dev/null; echo ---; tail -25 .mule/for-giulio.md 2>/dev/null
-rw-r--r-- 1 root root 52831 Sep 21 20:53 .mule/for-giulio.md
-rw-r--r-- 1 root root 73661 Sep 21 22:10 .mule/todo.md
---
- [x] auth-cache-default-size-test [auto]: pin in internal/auth/cache_test.go that `NewCache(0)` and `NewCache(-1)` fall back to `DefaultCacheSize` (cache.go:29-31) — the `size < 1` branch exists purely so `lru.New` cannot panic and has zero coverage; a regression dropping the fallback (panic on `lru.New` with a non-tick size, or a silently tiny cache) would pass every test. Construct with both sizes, assert no panic and a working `Verify` with `lru.Len() == 0`.
- [x] hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto] — already resolved and re-queued by the refill: commit 015a706 chose the upstream-parity behaviour (0 ≡ unset, not invalid, not a literal 0) and pinned it — create folds 0 → nil with the container-free TestCreateRealmZeroRetentionUnsets (internal/housekeeping/service_test.go:186) and PATCH 0 → ClearRetention in the integration suite (internal/housekeeping/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).
- [x] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication is "via Authorization header or `?token=` query parameter" (yaml:293-294), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
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
- [x] broker-acl-coldstart-fallback-flood [auto]: `syncOwnershipOf` (internal/broker/authhook.go:212-237, reached from the ACL miss path aclhook.go:196-206) does a full synchronous `GetDevice`+`GetInterface` store read for every distinct interface name a device publishes to, with no rate limit — `introspectionReloadDebounce` (authhook.go:42-45, "so an adversarial topic flood cannot hammer the database") throttles only `refreshIfStale`'s full reload, and the fallback's per-name cache write means N distinct bogus names in one window cause N synchronous full-device reads per second, an open anti-flood hole the pre-fix denied-miss path did not have. Gate the fallback to one sync read per session per debounce window (or fold it into a full reload that stamps `lastIntroLoad`), keep `TestBrokerACLColdStartIntrospectionMiss` (one name) green, and add a T1 test in broker_test.go with a GetDevice-counting fake: connect cold, publish to K>=2 distinct unknown-interface names, assert no more than one fallback read was made.
- [x] broker-offlineacl-entry-eviction [auto]: `offlineACL.entries` (internal/broker/aclhook.go:100-162) is append-only — `ownershipOf` upserts one entry per CN ever ACL-checked offline and nothing ever deletes, so a long-lived broker grows the map unboundedly across device churn even after the 10s `offlineACLCacheTTL` passes (the entry is kept and re-stamped, not reaped). Evict stale entries lazily on access (drop an entry whose `loadedAt` predates some multiple of the TTL before inserting another, or cap+LRU), add an injectable `now func() time.Time` in the same style as `lifecycleHook.now`, and cover in a container-free unit test: seed N entries, advance the clock, access one, assert the map stays bounded.
---
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
- **The mule has been idle 17h.** Filed by the dead-man's switch; see journalctl on the Pi.

$ tail -40 .mule/todo.md
- [x] docs-sync-ae-write-value-422 [auto]: add the missing `422` response to the six PUT/POST data-write ops in docs/api/astarte_appengine_api.yaml (device-scoped, by-alias, and in-group at `/interfaces/{interface}/{path}`). A value larger than `maxValueBytes` (64 KiB, internal/appengine/http.go:24) is answered `422` "Value size exceeds size limits" by putData/putDataByAlias/putDataInGroup (http.go:298-301, 327-329, 346-348), but the specs for PUT/POST (yaml:686-695, 726-736, 453-463, 499-509, 1176-1186, 1218-1228) list only 200/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-post-groups-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups` in docs/api/astarte_appengine_api.yaml — `Service.CreateGroup` returns `ErrGroupAlreadyExists` (internal/appengine/service.go:495, mapped at http.go:660-661) which the handler answers as `409` "Group already exists"; the spec (yaml:821-839) documents only 201/400/401/422/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] payload-longinteger-fraction-quantize [auto]: `jsonInt64` (pkg/payload/json.go:255-271) accepts fractional longinteger tokens that the ParseFloat fallback silently rounds onto an integer inside its own documented-safety window |x| < 2^53 (json.go:250-254 claims values there never round, "silently corrupting the value") — `{"v":4503599627370496.5}` decodes to 4503599627370496, `{"v":6000000000000000.9}` to 6000000000000001, and `{"v":5e-400}` underflows to 0, all without error, while the contract says fraction forms are "accepted only when integral". Verify the token is an exact decimal integer (`math/big.Rat` SetString + IsInt) on the float fallback — keeps `5.0`, `1e3`, `0.5e1` accepted, rejects the rounding cases as ReasonTypeMismatch; plain-integer ParseInt fast path and the > 2^53 window check unchanged (BSON unaffected: longinteger accepts int64/int32 only). Add failing rows to TestJSONProfile.
- [x] payload-json-malformed-t-not-tolerated [auto]: pin in pkg/payload/payload_test.go that a malformed `t` is rejected with ReasonBadTimestamp even on a non-explicit mapping — `decodeJSONEnvelope` (json.go:50-66) parses `t` before the facade learns the explicit-timestamp policy (payload.go:52-54), so only a *valid* `t` is "tolerated and ignored" while `{"v":1,"t":"yesterday"}` with explicit=false is still a rejection; only the valid-t-ignored row and the BSON "t wrong type" twin are asserted today, so a regression implementing the literal leniency reading (ignore a bad t on non-explicit mappings) passes every suite. Add one table row to TestJSONProfile.
- [x] docs-sync-ae-post-group-devices-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml — re-adding a member the group already holds answers `409` "Device already in group" (`Service.AddGroupDevice` returns `ErrDeviceAlreadyInGroup`, internal/appengine/service.go:636-638, mapped at http.go:662-663), but the spec (yaml:967-977) lists only 201/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-patch-merge-patch-media-type [auto]: fix the three PATCH request bodies in docs/api/astarte_appengine_api.yaml — PATCH /devices/{device} (yaml:176), PATCH /devices-by-alias/{alias} (yaml:260), PATCH /groups/{group}/devices/{device} (yaml:1023) declare `application/json; charset=utf-8`, but all three handlers enforce `Content-Type: application/merge-patch+json` exactly and reject every other value — 500 on the device and group paths (internal/appengine/http.go:209-211, 521-523), 400 on the alias path (http.go:230-232). Change the requestBody media type to `application/merge-patch+json` on all three. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
- [x] config-fail-loud-dev-mode [auto]: `applyEnv` silently discards a malformed `ASTRATE_MQTT_INSECURE_DEV_MODE` value (internal/config/config.go:255-257) — `strconv.ParseBool` error is discarded, the field becomes false, and TLS requirements are silently re-enabled (or a misleading "tls_cert_file required" error surfaces via validate). Return the parse error on bad input; add a failing test that asserts `ASTRATE_MQTT_INSECURE_DEV_MODE=yes` causes Load to error.
- [x] config-key-resolvers-test [auto]: add container-free unit tests for `HousekeepingKeys()` and `RealmJWTPublicKey()` (internal/config/config.go:383-411) — both file-reading helpers have zero test coverage. Cover: happy path (inline PEM + file PEM → concatenated/trimmed), file-not-found error (assert error message includes the path), inline-over-file priority (`RealmJWTPublicKey` returns the inline block when both are set), and TrimSpace behaviour (leading/trailing newlines stripped). `t.TempDir()` + `os.WriteFile` for fixtures; no Docker.
- [x] compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.
- [x] docs-sync-hk-delete-gating-responses [auto]: add the missing `405` and `422` responses to `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — `writeError` answers 405 "Realm deletion disabled" for `ErrDeletionDisabled` and 422 `{"errors":{"error_name":["connected_devices_present"]}}` for `ErrConnectedDevicesPresent` (internal/housekeeping/http.go:233-249), both reachable from `Service.DeleteRealm` (service.go:248-266), but the spec lists only 204/401/403/404/500 (yaml:198-208). Its sole `405` reference sits on PATCH (yaml:182-183), which can never return it (`UpdateRealm`, service.go:196-216, has no deletion-gating path) — re-point 405 to DELETE (reusing the MethodNotAllowed component) and add 422 with the FieldErrors shape and a connected_devices_present example. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_housekeeping_api.yaml (yaml:417) — it shows `detail: "housekeeping: validation failed: realm name must not be empty"`, but `validationDetail` strips the `housekeeping: validation failed: ` prefix on the wire (internal/housekeeping/http.go:251-258) and the real blank-name message is `"realm_name can't be blank"` (service.go:138). Change the example to `detail: "realm_name can't be blank"`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

- [x] auth-iat-not-required-test [auto]: pin in internal/auth/jwt_test.go that a token with a *future* `iat` still verifies — the "`iat` is not required — upstream parity" rule (jwt.go:104-105, no `jwt.WithIssuedAt` in the parser, jwt.go:115-118) has zero coverage, so a regression that starts validating `iat` or `WithIssuedAt` would silently reject tokens from a skewed-clock issuer and pass every existing suite. Sign RS256 with `iat` = now+1h and no `exp`/`nbf`, assert `Verify` succeeds.
- [x] auth-empty-claims-403-test [auto]: pin in internal/auth/middleware_test.go that a validly signed token carrying NO Astarte claims is `403` on a guarded route, not `401` — the authenticated-but-grant-less distinction (middleware.go:149-151, `Authorizes` on a nil grants map → false → `WriteForbidden`) is currently untested: every 403 fixture has a claim on the wrong surface, so a regression mapping grant-less tokens to 401 (or skipping the forbidden path) passes the suite. Sign RS256 with empty `jwt.MapClaims`, hit the pairing route, expect 403 + golden `envelope_403.json`.
- [x] auth-cache-default-size-test [auto]: pin in internal/auth/cache_test.go that `NewCache(0)` and `NewCache(-1)` fall back to `DefaultCacheSize` (cache.go:29-31) — the `size < 1` branch exists purely so `lru.New` cannot panic and has zero coverage; a regression dropping the fallback (panic on `lru.New` with a non-tick size, or a silently tiny cache) would pass every test. Construct with both sizes, assert no panic and a working `Verify` with `lru.Len() == 0`.
- [x] hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto] — already resolved and re-queued by the refill: commit 015a706 chose the upstream-parity behaviour (0 ≡ unset, not invalid, not a literal 0) and pinned it — create folds 0 → nil with the container-free TestCreateRealmZeroRetentionUnsets (internal/housekeeping/service_test.go:186) and PATCH 0 → ClearRetention in the integration suite (internal/housekeeping/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).
- [x] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication is "via Authorization header or `?token=` query parameter" (yaml:293-294), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
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
- [x] broker-acl-coldstart-fallback-flood [auto]: `syncOwnershipOf` (internal/broker/authhook.go:212-237, reached from the ACL miss path aclhook.go:196-206) does a full synchronous `GetDevice`+`GetInterface` store read for every distinct interface name a device publishes to, with no rate limit — `introspectionReloadDebounce` (authhook.go:42-45, "so an adversarial topic flood cannot hammer the database") throttles only `refreshIfStale`'s full reload, and the fallback's per-name cache write means N distinct bogus names in one window cause N synchronous full-device reads per second, an open anti-flood hole the pre-fix denied-miss path did not have. Gate the fallback to one sync read per session per debounce window (or fold it into a full reload that stamps `lastIntroLoad`), keep `TestBrokerACLColdStartIntrospectionMiss` (one name) green, and add a T1 test in broker_test.go with a GetDevice-counting fake: connect cold, publish to K>=2 distinct unknown-interface names, assert no more than one fallback read was made.
- [x] broker-offlineacl-entry-eviction [auto]: `offlineACL.entries` (internal/broker/aclhook.go:100-162) is append-only — `ownershipOf` upserts one entry per CN ever ACL-checked offline and nothing ever deletes, so a long-lived broker grows the map unboundedly across device churn even after the 10s `offlineACLCacheTTL` passes (the entry is kept and re-stamped, not reaped). Evict stale entries lazily on access (drop an entry whose `loadedAt` predates some multiple of the TTL before inserting another, or cap+LRU), add an injectable `now func() time.Time` in the same style as `lifecycleHook.now`, and cover in a container-free unit test: seed N entries, advance the clock, access one, assert the map stays bounded.

$ head -60 .mule/for-giulio.md
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **github-issues triage run, 2026-09-21: still nothing proposable — the mule-alarm pile is now 16 days (#94–#108).** Eighteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-20 evening run: **#108**, today's alarm (created 09-21 10:54Z, "nothing has landed in 15h"); **#107** has expired into the pile (superseded by #108). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-20 one: close #94–#107** — each was a one-day event, superseded by next day's alarm, never actionable; leave #108 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.

---

- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

---

- **github-issues triage re-run, 2026-09-20 (evening): nothing proposable, no change since this morning's 13:10 run.** Re-surveyed with the recipe command: the issue set is identical to the 2026-09-20 morning triage — **#94–#107** are still the daily `mule-alarm` noise (14 issues, zero comments, telemetry not code), **#93** aclhook comment rewrite is still mule-review with `8c61268` already pushed (its own recipe path), **#92** keyAgreement is still parked on a stable upstream v1.4.0 with its waiting-on.md row and the mule-upstream-watch escalation wired (an -rc does not satisfy it), **#78** FDO is still the milestone-4.0 design/investigation already escalated, and **#1** is untouched per standing instruction. No new issue has been filed since the morning run. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines proposed**. This morning's proposal stands unchanged: **close #94–#106** — each was a one-day event, superseded by the next day's alarm, never actionable — and leave **#107** (still the newest, no #108 yet) to self-expire. The queue idle streak since ~2026-09-04/05 persists; if that is not intentional the wedge is on the Pi (queue review = the fix, per the 13:10 line).

---

- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

---

- **COMPATIBILITY.md §8 wording proposal: trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (Forwarder) actions are single-shot by design, receiving only `maximum_capacity`/`event_ttl`.** Propose adding to the "Trigger delivery policies" deviation: *"Policy `error_handlers` (retry/discard) and the `retry_times` attempt cap govern HTTP webhook actions only: a custom action delivered through the Forwarder seam is single-shot by design, receiving from the policy solely its `maximum_capacity` in-flight bound and its `event_ttl` staleness drop, with a forward failure counted once and never retried."* Verified in `internal/engine/triggers/actions.go`: `Action.Custom` routes to `forward` (actions.go:435-439), which never consults `policy.Decide` — the Forwarder seam returns only an error, no HTTP status (actions.go:161-167, doc comment at actions.go:471-479) — while `Enqueue` bounds `maximum_capacity` on every delivery (actions.go:338-349) and `deliver` drops past `event_ttl` before the custom/webhook branch (actions.go:400-416). The existing warn-on-single-shot line already tells the operator when a policy's `retry_times` cannot apply to the forward path (actions.go:491-497).

---

- **github-issues triage run, 2026-09-20: still nothing proposable — the mule-alarm pile is now 15 days (#94–#107).** Seventeen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-18 run: **#107**, yesterday's alarm (created 09-19 10:49Z, "nothing has landed in 14h"); **#106** has expired into the pile (superseded by #107). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-18 one: close #94–#106** — each was a one-day event, superseded by next day's alarm, never actionable; leave #107 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.

---

- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

---

- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

---

- **deviceid wire-parity concern, from the 2026-09-18 code review: Astrate is likely stricter than upstream on the wire form's unused bits.** `pkg/deviceid/deviceid.go:37` uses `base64.RawURLEncoding.Strict()`, so `Parse` rejects `"AAAAAAAAAAAAAAAAAAAAAB"` (unused low 4 bits of the 22nd char set; rejection pinned at deviceid_test.go:176), while the comment at deviceid.go:34-36 claims parity with `Elixir Base.url_decode64!(padding: false)` — which decodes the 128 payload bits and discards the unused ones rather than erroring. If so, a device registered with a non-canonical-trailing-bits id on upstream Astarte is accepted there but rejected by Astrate's `Parse` — a wire-compat break in a package whose purpose is the wire form. Not verified on the Pi (no Docker/Elixir); a `[legion]` probe (`deviceid-trailing-bits-upstream-probe`) is queued to measure it. Once confirmed, the call is yours: relax `Parse` to upstream's decode-and-ignore, or keep strict canonicality as the package's contract.

---

- **COMPATIBILITY.md wording update for upstream v1.3.4 (newest stable, 2026-09-18; v1.4.0 is still rc.5-only) — folds v1.3.3 and v1.3.4 into the open v1.3.2 proposal below; both are maintenance-only and wire-inert, so that proposal stands unchanged.** v1.3.3 (2026-08-07) was an empty-body patch and v1.3.4 (2026-09-18, `astarte_data_updater_plant` only — more frequent fullsweep GC on AMQPDataConsumer processes plus an RPC-availability fix; release body via `gh api repos/astarte-platform/astarte/releases` `v1.3.4`) introduce **no wire/API surface change**, so the v1.3.2 wording proposal stands complete: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"). The only delta a v1.3.4-aware doc carries is the version reference: the proposed §Infrastructure-differences sentence's "until the milestone that adopts v1.3.2 as the target" reads "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target". Raw: [v1.3.4](https://github.com/astarte-platform/astarte/releases/tag/v1.3.4).

---

- **github-issues triage run, 2026-09-18: still nothing proposable — the mule-alarm pile is now 13 straight days (#94–#106).** Sixteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-17 run: **#106**, today's alarm (created 09-18 10:54Z, "nothing has landed in 14h"); **#105** has expired into the pile (superseded by #106). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-17 one: close #94–#105** — each was a one-day event, superseded by next day's alarm, never actionable; leave #106 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.

---

- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) shows a new **stable v1.3.4** (2026-09-18) — maintenance only (data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change), so no new gap for Astrate — and still **no stable v1.4.0**: newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

---

$ rg -n 'realmID|ErrRealmNotFound|GetRealmByName|missingRealm|ErrNotFound' internal/appengine/service.go | head -20
65:// ErrGroupNotFound (→ 404 "Group not found") and store.ErrNotFound.
67:	return fmt.Errorf("%w: %w: group %q", ErrGroupNotFound, store.ErrNotFound, name)
93:func (s *Service) realmID(ctx context.Context, realm string) (int16, error) {
94:	r, err := s.st.GetRealmByName(ctx, realm)
142:	rid, err := s.realmID(ctx, realm)
183:	rid, err := s.realmID(ctx, realm)
211:	rid, err := s.realmID(ctx, realm)
217:		return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
229:	rid, err := s.realmID(ctx, realm)
292:	rid, err := s.realmID(ctx, realm)
298:		return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
310:	rid, err := s.realmID(ctx, realm)
404:// registered-but-not-member (or unparseable/unknown) device is store.ErrNotFound,
407:	rid, err := s.realmID(ctx, realm)
416:		return 0, deviceid.ID{}, nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
423:		return 0, deviceid.ID{}, nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
463:	rid, err := s.realmID(ctx, realm)
477:			if errors.Is(err, store.ErrNotFound) {
509:	rid, err := s.realmID(ctx, realm)
521:	rid, err := s.realmID(ctx, realm)

→ Read internal/appengine/service.go [limit=18, offset=88]
→ Read internal/auth/middleware.go [limit=115, offset=20]

$ rg -n 'func \(s \*Service\) GetData\b|func \(s \*Service\) GetDataByAlias|func \(s \*Service\) GetDataInGroup|func \(s \*Service\) readData' internal/appengine/service.go
726:func (s *Service) GetDataByAlias(ctx context.Context, realm, alias, ifaceName, path string, opts QueryOpts) (any, error) {
736:func (s *Service) GetDataInGroup(ctx context.Context, realm, groupName, deviceID, ifaceName, path string, opts QueryOpts) (any, error) {

$ rg -n 'GetData\b' internal/appengine/service.go internal/appengine/data.go | head; echo ---; rg -n 'GetData|read' internal/appengine/data.go | head
internal/appengine/service.go:712:// the shared tail of GetData and its mirror variants.
internal/appengine/service.go:715:		// See GetData: reducing a snapshot is refused, not silently ignored.
internal/appengine/data.go:106:// GetData reads an interface endpoint (upstream GET .../interfaces/{iface}
internal/appengine/data.go:109:func (s *Service) GetData(ctx context.Context, realm, deviceID, ifaceName, path string, opts QueryOpts) (any, error) {
---
35:// to/limit/downsample_to plus the rendering options). The zero value reads
106:// GetData reads an interface endpoint (upstream GET .../interfaces/{iface}
109:func (s *Service) GetData(ctx context.Context, realm, deviceID, ifaceName, path string, opts QueryOpts) (any, error) {
126:// datastreamData reads a datastream endpoint. For an object-aggregated
414:// propertiesData reads a properties interface. With a concrete path it returns
417:// which astarte-go flattens via parsePropertiesMap). The stored jsonb already

→ Read internal/appengine/data.go [limit=40, offset=106]
I have a clear picture of the appengine surface. The path/method coverage is fully in sync; the real drift: two query params the code accepts are undocumented, plus three `ASTRATE_` config keys absent from the reference. Let me record the deliverables.

$ tail -3 .mule/todo.md | cat -A | tail -3
- [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (internal/store/devices.go:130-131, "if several devices share an alias the lowest device ID wins", `ORDER BY id LIMIT 1`) M-bM-^@M-^T no test drives it today; add a case in internal/store/devices_test.go that gives two devices the same alias and asserts the lookup resolves the lower ID. Needs the integration DB.$
- [x] broker-acl-coldstart-fallback-flood [auto]: `syncOwnershipOf` (internal/broker/authhook.go:212-237, reached from the ACL miss path aclhook.go:196-206) does a full synchronous `GetDevice`+`GetInterface` store read for every distinct interface name a device publishes to, with no rate limit M-bM-^@M-^T `introspectionReloadDebounce` (authhook.go:42-45, "so an adversarial topic flood cannot hammer the database") throttles only `refreshIfStale`'s full reload, and the fallback's per-name cache write means N distinct bogus names in one window cause N synchronous full-device reads per second, an open anti-flood hole the pre-fix denied-miss path did not have. Gate the fallback to one sync read per session per debounce window (or fold it into a full reload that stamps `lastIntroLoad`), keep `TestBrokerACLColdStartIntrospectionMiss` (one name) green, and add a T1 test in broker_test.go with a GetDevice-counting fake: connect cold, publish to K>=2 distinct unknown-interface names, assert no more than one fallback read was made.$
- [x] broker-offlineacl-entry-eviction [auto]: `offlineACL.entries` (internal/broker/aclhook.go:100-162) is append-only M-bM-^@M-^T `ownershipOf` upserts one entry per CN ever ACL-checked offline and nothing ever deletes, so a long-lived broker grows the map unboundedly across device churn even after the 10s `offlineACLCacheTTL` passes (the entry is kept and re-stamped, not reaped). Evict stale entries lazily on access (drop an entry whose `loadedAt` predates some multiple of the TTL before inserting another, or cap+LRU), add an injectable `now func() time.Time` in the same style as `lifecycleHook.now`, and cover in a container-free unit test: seed N entries, advance the clock, access one, assert the map stays bounded.$

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -204,4 +204,5 @@
 - [x] store-pipelines-empty-name-zero-blocks-test [auto]: extend internal/store/pipelines_validate_test.go with the two untested error branches of `validatePipelineGraph` (internal/store/pipelines.go:51-64) — a definition with zero blocks and one whose block has an empty name must both be rejected. Pure helper, no DB; the existing suite covers only acyclic/cyclic/unknown-ref/duplicate-name.
 - [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (internal/store/devices.go:130-131, "if several devices share an alias the lowest device ID wins", `ORDER BY id LIMIT 1`) — no test drives it today; add a case in internal/store/devices_test.go that gives two devices the same alias and asserts the lookup resolves the lower ID. Needs the integration DB.
 - [x] broker-acl-coldstart-fallback-flood [auto]: `syncOwnershipOf` (internal/broker/authhook.go:212-237, reached from the ACL miss path aclhook.go:196-206) does a full synchronous `GetDevice`+`GetInterface` store read for every distinct interface name a device publishes to, with no rate limit — `introspectionReloadDebounce` (authhook.go:42-45, "so an adversarial topic flood cannot hammer the database") throttles only `refreshIfStale`'s full reload, and the fallback's per-name cache write means N distinct bogus names in one window cause N synchronous full-device reads per second, an open anti-flood hole the pre-fix denied-miss path did not have. Gate the fallback to one sync read per session per debounce window (or fold it into a full reload that stamps `lastIntroLoad`), keep `TestBrokerACLColdStartIntrospectionMiss` (one name) green, and add a T1 test in broker_test.go with a GetDevice-counting fake: connect cold, publish to K>=2 distinct unknown-interface names, assert no more than one fallback read was made.
 - [x] broker-offlineacl-entry-eviction [auto]: `offlineACL.entries` (internal/broker/aclhook.go:100-162) is append-only — `ownershipOf` upserts one entry per CN ever ACL-checked offline and nothing ever deletes, so a long-lived broker grows the map unboundedly across device churn even after the 10s `offlineACLCacheTTL` passes (the entry is kept and re-stamped, not reaped). Evict stale entries lazily on access (drop an entry whose `loadedAt` predates some multiple of the TTL before inserting another, or cap+LRU), add an injectable `now func() time.Time` in the same style as `lifecycleHook.now`, and cover in a container-free unit test: seed N entries, advance the clock, access one, assert the map stays bounded.
+- [ ] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine_api.yaml — `retrieve_metadata` and `downsample_key` are parsed for every data read by `parseQueryOpts` (internal/appengine/http.go:624 `retrieve_metadata`, 633 `downsample_key`) yet appear on none of the read endpoints: device-scoped GET `{interface}` (yaml:583-606) and `{interface}/{path}` (yaml:623-649), by-alias (yaml:324-368, 369-417), and in-group (yaml:1108-1132, 1149-1176) — the parameter block on each lists only since/since_after/to/limit/downsample_to/sort/format/allow_bigintegers/allow_safe_bigintegers (components DataSince..DataAllowSafeBigIntegers, yaml:1360-1439). Add the two `$ref`s (new DataRetrieveMetadata and DataDownsampleKey components, or inline — mirror the existing style) to all six, and update the DataSort-independent `sort` default note only if verified. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`ASTRATE_FLOW_CONFIG`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
+
+---
+
 - **github-issues triage run, 2026-09-21: still nothing proposable — the mule-alarm pile is now 16 days (#94–#108).** Eighteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-20 evening run: **#108**, today's alarm (created 09-21 10:54Z, "nothing has landed in 15h"); **#107** has expired into the pile (superseded by #108). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-20 one: close #94–#107** — each was a one-day event, superseded by next day's alarm, never actionable; leave #108 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
 
 ---
 

$ tail -2 .mule/todo.md
- [x] broker-offlineacl-entry-eviction [auto]: `offlineACL.entries` (internal/broker/aclhook.go:100-162) is append-only — `ownershipOf` upserts one entry per CN ever ACL-checked offline and nothing ever deletes, so a long-lived broker grows the map unboundedly across device churn even after the 10s `offlineACLCacheTTL` passes (the entry is kept and re-stamped, not reaped). Evict stale entries lazily on access (drop an entry whose `loadedAt` predates some multiple of the TTL before inserting another, or cap+LRU), add an injectable `now func() time.Time` in the same style as `lifecycleHook.now`, and cover in a container-free unit test: seed N entries, advance the clock, access one, assert the map stays bounded.
- [ ] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine_api.yaml — `retrieve_metadata` and `downsample_key` are parsed for every data read by `parseQueryOpts` (internal/appengine/http.go:624 `retrieve_metadata`, 633 `downsample_key`) yet appear on none of the read endpoints: device-scoped GET `{interface}` (yaml:583-606) and `{interface}/{path}` (yaml:623-649), by-alias (yaml:324-368, 369-417), and in-group (yaml:1108-1132, 1149-1176) — the parameter block on each lists only since/since_after/to/limit/downsample_to/sort/format/allow_bigintegers/allow_safe_bigintegers (components DataSince..DataAllowSafeBigIntegers, yaml:1360-1439). Add the two `$ref`s (new DataRetrieveMetadata and DataDownsampleKey components, or inline — mirror the existing style) to all six, and update the DataSort-independent `sort` default note only if verified. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ head -22 .mule/for-giulio.md
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`ASTRATE_FLOW_CONFIG`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.

---

- **github-issues triage run, 2026-09-21: still nothing proposable — the mule-alarm pile is now 16 days (#94–#108).** Eighteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-20 evening run: **#108**, today's alarm (created 09-21 10:54Z, "nothing has landed in 15h"); **#107** has expired into the pile (superseded by #108). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-20 one: close #94–#107** — each was a one-day event, superseded by next day's alarm, never actionable; leave #108 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.

---

- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

Done. Proposal job carried out for the **appengine** surface.

Surface comparison (`internal/appengine/http.go:40-81` vs `docs/api/astarte_appengine_api.yaml`): all 32 routed methods fall inside the 18 documented path blocks — no path or method drift. Checked status codes/responses on 4 endpoints (GET /devices, PATCH /devices/{device}, POST /groups, GET/PUT/DELETE at interface path): all consistent, including `data`-wrapped request bodies (matches `DecodeData`, `pkg/astarteapi/envelope.go:228`), 409s, 422 ValueTooLarge, 405s, and request-media-type. One real gap found and one escalation proposed; spec is otherwise accurate.

Outcomes:
- `docs/api` fix task (1 line) appended to `.mule/todo.md`: `retrieve_metadata` and `downsample_key` are parsed for every data read (`parseQueryOpts`, internal/appengine/http.go:624, 633) but appear on none of the six data-read GET endpoints. A docs/api fix is a normal task per the recipe.
- `.mule/for-giulio.md` (newest-first): three `ASTRATE_` keys exist in code but are absent from `configuration-reference.md` (`ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`, `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` — real startup keys; `ASTRATE_FLOW_CONFIG` — per-theme flow-block env contract, your call where it belongs). Reverse direction clean.
- No code changed, no git touched, no `gh issue create` (recipe doesn't request one).

```
Done: appengine docs-sync proposal — routes fully in sync; 1 docs/api task (retrieve_metadata/downsample_key) to todo.md, 3 undocumented ASTRATE_ keys escalated to for-giulio.md
Files: .mule/todo.md, .mule/for-giulio.md
Verified: path/method diff (code vs yaml), 4 endpoint status-code/field checks, ASTRATE_ key rg sweep -> all pass, no drift beyond the reported gaps
Unsure: whether ASTRATE_FLOW_CONFIG belongs in configuration-reference or a flow page (flagged in the escalation)
Follow-ups: docs-sync-ae-read-query-params (appended); for-giulio.md entry added (no task line — docs/site is never-edit)
```
