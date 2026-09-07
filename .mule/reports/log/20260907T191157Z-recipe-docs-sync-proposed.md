slug: recipe-docs-sync
verdict: proposed
at:  a42f2cf
ran: 2026-09-07T19:11:57Z on DietPi in 341s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

# Todos
[•] Pick a surface and find drift between code and docs/api/*.yaml
[ ] Compare documented paths vs actual handlers for 3-4 endpoints
[ ] Propose corrections as task lines in .mule/todo.md

Let me see what API surfaces exist and pick one.
✱ Glob "docs/api/*.yaml" 5 matches

$ ls /root/astrate-mule/internal/ | head -30
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

I'll pick **appengine**. Let me find the routes in both the code and the YAML spec.
$ rg -n '^\s{2}/' docs/api/astarte_appengine_api.yaml
26:  /appengine/v1/{realm}/devices:
93:  /appengine/v1/{realm}/stats/devices:
131:  /appengine/v1/{realm}/devices/{device}:
205:  /appengine/v1/{realm}/devices-by-alias/{alias}:
287:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces:
319:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
361:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
534:  /appengine/v1/{realm}/devices/{device}/interfaces:
562:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
599:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
752:  /appengine/v1/{realm}/groups:
824:  /appengine/v1/{realm}/groups/{group}:
852:  /appengine/v1/{realm}/groups/{group}/devices:
940:  /appengine/v1/{realm}/groups/{group}/devices/{device}:
1039:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces:
1067:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}:
1105:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}:

No match for the glob pattern. Let me search for the routing registration in the appengine package differently.
$ rg -n '\.(Get|Post|Put|Delete|Patch)\(' internal/appengine/ | head -40
internal/appengine/stream/ws_test.go:81:		if ct := rec.Header().Get("Content-Type"); ct != "text/event-stream" {
internal/appengine/stream/ws.go:72:		DeviceID:  r.URL.Query().Get("device_id"),
internal/appengine/stream/ws.go:73:		Interface: r.URL.Query().Get("interface"),
internal/appengine/stream/ws.go:152:	return r.Header.Get("Accept") == "text/event-stream" || r.URL.Query().Get("transport") == "sse"
internal/appengine/groups_parity_test.go:240:			if got := u.Query().Get("from_token"); got == "" {
internal/appengine/groups_parity_test.go:245:			if got := u.Query().Get("limit"); got != "2" {
internal/appengine/http.go:87:	limit, _ := strconv.Atoi(q.Get("limit"))
internal/appengine/http.go:88:	details, _ := strconv.ParseBool(q.Get("details"))
internal/appengine/http.go:89:	page, err := a.svc.ListDevices(r.Context(), r.PathValue("realm"), q.Get("from_token"), limit, details)
internal/appengine/http.go:116:		details := q.Get("details")
internal/appengine/http.go:120:		limit := q.Get("limit")
internal/appengine/http.go:209:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:230:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:438:	details, _ := strconv.ParseBool(q.Get("details"))
internal/appengine/http.go:439:	limit, _ := strconv.Atoi(q.Get("limit"))
internal/appengine/http.go:441:		details, q.Get("from_token"), limit)
internal/appengine/http.go:469:	limit := q.Get("limit")
internal/appengine/http.go:474:	if d := q.Get("details"); d != "" {
internal/appengine/http.go:521:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:564:		if v := q.Get(p.name); v != "" {
internal/appengine/http.go:573:	if q.Get("since") != "" && q.Get("since_after") != "" {
internal/appengine/http.go:576:	if v := q.Get("limit"); v != "" {
internal/appengine/http.go:587:	if v := q.Get("downsample_to"); v != "" {
internal/appengine/http.go:601:	switch v := q.Get("format"); v {
internal/appengine/http.go:615:		switch v := q.Get(p.name); v {
internal/appengine/http.go:624:	if v := q.Get("retrieve_metadata"); v != "" {
internal/appengine/http.go:633:	opts.DownsampleKey = q.Get("downsample_key")
internal/appengine/http.go:635:	opts.Descending = q.Get("sort") != "ascending"
internal/appengine/channels/ws.go:53:	realm := r.URL.Query().Get("realm")
internal/appengine/channels/ws.go:54:	token := r.URL.Query().Get("token")

✱ Grep "HandleFunc|\.Get\(|\.Post\(|\.Put\(|\.Delete\(|\.Patch\(" in internal/appengine · 30 matches
Let me find where the router registers the paths in appengine.
✱ Glob "internal/appengine/**/*.go" 23 matches
✱ Grep ""/|Handle\(|HandleFunc|registerRoutes|mux\." in internal/appengine · 100 matches
Now I have the full route list from `internal/appengine/http.go`. Let me compare it against the documented paths. 

From the code routes (all under base `/appengine/v1/{realm}`):
- GET /devices
- GET /stats/devices
- GET /devices/{device}
- PATCH /devices/{device}
- GET /devices-by-alias/{alias}
- GET /devices/{device}/interfaces
- GET /devices/{device}/interfaces/{interface}
- GET /devices/{device}/interfaces/{interface}/{path...}
- PUT /devices/{device}/interfaces/{interface}/{path...}
- POST /devices/{device}/interfaces/{interface}/{path...}
- DELETE /devices/{device}/interfaces/{interface}/{path...}
- PATCH /devices-by-alias/{alias}
- GET /devices-by-alias/{alias}/interfaces
- GET /devices-by-alias/{alias}/interfaces/{interface}
- GET /devices-by-alias/{alias}/interfaces/{interface}/{path...}
- PUT/POST/DELETE /devices-by-alias/.../path...
- GET /groups
- POST /groups
- GET /groups/{group}
- GET /groups/{group}/devices
- POST /groups/{group}/devices
- GET /groups/{group}/devices/{device}
- PATCH /groups/{group}/devices/{device}
- DELETE /groups/{group}/devices/{device}
- GET/PUT/POST/DELETE /groups/{group}/devices/{device}/interfaces...

Now let me look at the YAML documented paths. I need to see which methods are documented for each path. Let me look at the yaml more carefully.

$ rg -n '^\s{2}/|^    (get|post|put|delete|patch):|^\s{4}(get|post|put|delete|patch):' docs/api/astarte_appengine_api.yaml
26:  /appengine/v1/{realm}/devices:
27:    get:
93:  /appengine/v1/{realm}/stats/devices:
94:    get:
131:  /appengine/v1/{realm}/devices/{device}:
132:    get:
156:    patch:
205:  /appengine/v1/{realm}/devices-by-alias/{alias}:
206:    get:
235:    patch:
287:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces:
288:    get:
319:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
320:    get:
361:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
362:    get:
408:    put:
454:    post:
500:    delete:
534:  /appengine/v1/{realm}/devices/{device}/interfaces:
535:    get:
562:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
563:    get:
599:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
600:    get:
639:    put:
680:    post:
721:    delete:
752:  /appengine/v1/{realm}/groups:
753:    get:
778:    post:
824:  /appengine/v1/{realm}/groups/{group}:
825:    get:
852:  /appengine/v1/{realm}/groups/{group}/devices:
853:    get:
900:    post:
940:  /appengine/v1/{realm}/groups/{group}/devices/{device}:
941:    get:
968:    patch:
1018:    delete:
1039:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces:
1040:    get:
1067:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}:
1068:    get:
1105:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}:
1106:    get:
1146:    put:
1188:    post:
1230:    delete:

The documented paths and methods match the code routes exactly — all the paths line up with the routes in `http.go`, including the PATCH on devices-by-alias, POST on data paths, etc. Both surfaces are in sync in terms of path/method coverage.

Now let me check the details of a few endpoints against the handlers as the recipe instructs. Let me pick a few endpoints and verify their documented status codes and response fields. Let me first look at the actual handler responses in the code. Let me focus on a few representative endpoints.

Let me check the documented status codes and compare with what the handlers actually return. Let me look at some handlers in http.go.
→ Read internal/appengine/http.go [offset=80, limit=160]
Let me compare the documented `devices/{device}` patch status codes and the `devices` get responses in the YAML. Let me look at the YAML for a couple of endpoints to verify status codes.
→ Read docs/api/astarte_appengine_api.yaml [offset=131, limit=75]
The PATCH `/devices/{device}` handler maps the wrong Content-Type to an unmapped 500 (per the comment), but the YAML documents 400 for it. The handler explicitly returns `WriteInternalServerError` for a wrong content type. This is a discrepancy — the code comment even notes that upstream surfaces an unmapped 500 on the device path, but not on the alias path (which returns 400). The YAML documents 400 for `/devices/{device}` patch but the code returns 500 for that case.

However, this matches upstream behavior which the code deliberately replicates. This is worth checking carefully. Let me note this is a potential drift: the YAML documents 400 for patchDevice, but the code returns 500 for the content-type mismatch. Let me verify the YAML documents this correctly. Actually the code comment says the 500 is "replicated as-is" from upstream — so the YAML may be wrong (documenting 400 where code returns 500). This is a candidate fix task.

Let me verify a few more endpoints. Let me check the group device patch and the stats devices. Let me look at the stats devices handler which returns a map, and check the YAML.

Let me look at more of the http.go, particularly around group operations, and check the group device patch status code.
→ Read internal/appengine/http.go [limit=130, offset=420]
So for the three PATCH endpoints (`patchDevice`, `patchGroupDevice`), wrong Content-Type maps to an unmapped 500 (handlers return `WriteInternalServerError` on line 210 and 522), while `patchDeviceByAlias` maps it to 400 (line 231).

Let me check what the YAML documents for each of these PATCH endpoints' 400 response. Let me look at the YAML for group device patch and the alias patch.
→ Read docs/api/astarte_appengine_api.yaml [limit=30, offset=968]
→ Read docs/api/astarte_appengine_api.yaml [limit=22, offset=998]
The YAML documents 400 for `patchGroupDevice` (line 1005), but the code handler returns **500** (`WriteInternalServerError`) for the content-type mismatch (line 522). Similarly for `patchDevice`, the YAML documents 400 (line 192) but the code returns 500 (line 210).

This is a documented drift. However, this is subtle — these are edge cases (wrong content type). The YAML doesn't distinguish "400 from content-type mismatch". The 400 in the YAML likely refers to bad request body (malformed JSON), which does return 400 (`WriteBadRequest`). So the 400 is still valid for the bad-body case. The content-type mismatch is a 500 edge case that's not explicitly documented. This is a minor nuance and arguably the docs aren't strictly wrong — they don't claim what happens on a wrong content type.

Let me not over-reach on that. Let me look at the `removeGroupDevice` delete handler which returns 204 NoContent, and check the YAML. Let me also verify the `addGroupDevice` returns 201.

Let me check the YAML for the group device delete and the add group device.

Actually, let me focus my checks. The recipe says check 3-4 endpoints' documented status codes and response fields. Let me verify a few concrete ones against handlers:

1. `GET /groups/{group}` returns `{"group_name": ...}` (http.go:433) — let me check YAML documents this.
2. `DELETE /groups/{group}/devices/{device}` returns 204 (http.go:504) — check YAML.
3. `POST /groups/{group}/devices` returns 201 (http.go:496) — check YAML.

Let me look at the YAML for the GET group path and the responses.
→ Read docs/api/astarte_appengine_api.yaml [limit=30, offset=824]
GET group matches (`group_name`). Let me check the delete group device and add group device responses in YAML.
→ Read docs/api/astarte_appengine_api.yaml [limit=45, offset=900]
POST addGroupDevice returns 201, matching the handler (http.go:496). Good. Let me check the DELETE group device response code in the YAML.
→ Read docs/api/astarte_appengine_api.yaml [limit=22, offset=1018]
DELETE group device returns 204, matching handler (http.go:504). All documented and code agree on the endpoints I've checked.

Now let me verify a couple of the data endpoints' response fields. Let me look at the getData handler and its documented response. Let me look at the interface value GET response schema in the YAML. Let me look at the data get endpoint.

Let me check the data endpoints 599 onwards where the YAML documents the data operations. Let me look at the getData handler to see the response format.
→ Read internal/appengine/http.go [limit=75, offset=550]
Let me look at the YAML data GET endpoint to check documented query params vs code. Let me look at the YAML for the data get endpoint responses and parameters, and check against the code's query params: since, since_after, to, limit, downsample_to, format, allow_bigintegers, allow_safe_bigintegers, retrieve_metadata, downsample_key, sort.
→ Read docs/api/astarte_appengine_api.yaml [limit=150, offset=599]
Interesting. The `getDataAtPath` GET documents only these query params: `DataSince`, `DataSinceAfter`, `DataTo`, `DataLimit`, `DataDownsample`, `DataSort`. 

But the code (`parseQueryOpts` in http.go) also supports: `format` (structured/table/disjoint_tables), `allow_bigintegers`, `allow_safe_bigintegers`, `retrieve_metadata`, `downsample_key`, and implicitly `downsample_to`. Wait, `DataDownsample` is the downsample_to param. Let me check the definitions of these parameter refs.

Let me check the component parameter definitions for what query param names they map to.

$ grep -n 'DataSince:\|DataSinceAfter:\|DataTo:\|DataLimit:\|DataDownsample:\|DataSort:' docs/api/astarte_appengine_api.yaml
1303:    DataSince:
1311:    DataSinceAfter:
1319:    DataTo:
1327:    DataLimit:
1334:    DataDownsample:
1341:    DataSort:

→ Read docs/api/astarte_appengine_api.yaml [limit=75, offset=1303]
The GET data path endpoint in the YAML documents only 6 query params (`since`, `since_after`, `to`, `limit`, `downsample_to`, `sort`), but the code (`parseQueryOpts` in http.go:554-641) supports several more: `format`, `allow_bigintegers`, `allow_safe_bigintegers`, `retrieve_metadata`, `downsample_key`. 

This is a real drift: **present but undocumented** query parameters. The code supports these additional query params on the data GET endpoints (`getDataAtPath`, `getDataByAlias`, `getDataInGroup`, all of which use `parseQueryOpts`). They are not documented in the YAML.

This is a legitimate `docs/api/` fix task per the recipe ("present but undocumented").

Let me verify the getData handler uses parseQueryOpts and confirm which endpoints these params apply to.

$ grep -n 'parseQueryOpts\|getData\|getDataByAlias\|getDataInGroup' internal/appengine/http.go | head -30
50:	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}", h(a.getData))
51:	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.getData))
61:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces/{interface}", h(a.getDataByAlias))
62:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.getDataByAlias))
76:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}", h(a.getDataInGroup))
77:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.getDataInGroup))
254:	opts, err := parseQueryOpts(r)
271:func (a *API) getData(w http.ResponseWriter, r *http.Request) {
278:func (a *API) getDataByAlias(w http.ResponseWriter, r *http.Request) {
285:func (a *API) getDataInGroup(w http.ResponseWriter, r *http.Request) {
551:// parseQueryOpts reads the datastream query parameters. Malformed values come
554:func parseQueryOpts(r *http.Request) (QueryOpts, error) {

Confirmed. `getData`, `getDataByAlias`, `getDataInGroup` all use `parseQueryOpts`, which parses `format`, `allow_bigintegers`, `allow_safe_bigintegers`, `retrieve_metadata`, `downsample_key` in addition to the six documented params.

Let me verify the three data GET specs (`devices-by-alias/interfaces/{interface}`, `groups/.../interfaces/{interface}`) also document only the six, to confirm the scope of the drift. Let me check those sections of the YAML.
→ Read docs/api/astarte_appengine_api.yaml [limit=45, offset=319]
The data GET endpoints consistently document only six query params while the code supports five more. Let me check the devices list GET response contract in the YAML to complete my endpoint checks.
→ Read docs/api/astarte_appengine_api.yaml [limit=67, offset=26]
GET /devices is well documented. Let me check the group devices GET spec for its query params vs the handler.
→ Read docs/api/astarte_appengine_api.yaml [limit=48, offset=852]
→ Read docs/api/astarte_appengine_api.yaml [limit=37, offset=562]
→ Read .mule/todo.md
Drift confirmed across all six data GET operations. Let me check the current state of `.mule/todo.md` and `.mule/for-giulio.md`.

$ ls -la /root/astrate-mule/.mule/ 2>/dev/null && echo "---" && test -f /root/astrate-mule/.mule/for-giulio.md && tail -5 /root/astrate-mule/.mule/for-giulio.md || echo "no for-giulio.md"
total 1460
drwxr-xr-x  7 root root    4096 Sep  7 21:03 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      14 Sep  7 21:06 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep  7 19:00 .heartbeat
-rw-r--r--  1 root root   16384 Sep  7 21:08 .last-output
-rw-r--r--  1 root root       0 Sep  7 21:06 .lock
-rw-r--r--  1 root root      10 Sep  7 21:06 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 1344485 Sep  7 21:08 cron.log
drwxr-xr-x  2 root root    4096 Sep  7 21:03 failed
-rw-r--r--  1 root root   20080 Sep  7 19:48 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    2788 Sep  7 21:03 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep  7 20:15 reviews
drwxr-xr-x  2 root root    4096 Sep  4 20:46 tasks
-rw-r--r--  1 root root   23075 Sep  7 21:03 todo.md
---
manufacturing/rendezvous, acceptance = official `astarte-device-fdo-rust`
SDK completes onboarding against Astrate, docs as a first-class deliverable.
When v3.0 is marked DONE, the v4.0 section of `.mule/milestones.md` should be
drafted with this investigation as its first item (issue #78 has the full
verified context).

$ rg -n 'Format|AllowBigIntegers|AllowSafeBigIntegers|RetrieveMetadata|DownsampleKey' internal/appengine/*.go | rg -v '_test' | head -40
internal/appengine/data.go:42:// Format is one of structured/table/disjoint_tables ("structured" when the
internal/appengine/data.go:43:// request omitted it). AllowBigIntegers/AllowSafeBigIntegers are nil when the
internal/appengine/data.go:52:	Format               string
internal/appengine/data.go:53:	DownsampleKey        string
internal/appengine/data.go:54:	RetrieveMetadata     bool
internal/appengine/data.go:55:	AllowBigIntegers     *bool
internal/appengine/data.go:56:	AllowSafeBigIntegers *bool
internal/appengine/data.go:136:	// via parseDatastreamMap, keyed on the "value" leaf field). Format is
internal/appengine/data.go:230:		if opts.AllowBigIntegers != nil || opts.AllowSafeBigIntegers != nil {
internal/appengine/data.go:237:		return renderObject(docs, stamps, opts.Format), nil
internal/appengine/data.go:244:	return renderIndividual(rows, path, opts.Format), nil
internal/appengine/data.go:403:			case opts.AllowBigIntegers != nil && !*opts.AllowBigIntegers:
internal/appengine/data.go:405:			case opts.AllowSafeBigIntegers != nil && *opts.AllowSafeBigIntegers:
internal/appengine/data.go:500:		return strconv.FormatInt(*r.ValueLonginteger, 10)
internal/appengine/data.go:508:		return r.ValueDatetime.UTC().Format(datetimeLayout)
internal/appengine/http.go:600:	opts.Format = "structured"
internal/appengine/http.go:604:		opts.Format = v
internal/appengine/http.go:612:		{"allow_bigintegers", &opts.AllowBigIntegers},
internal/appengine/http.go:613:		{"allow_safe_bigintegers", &opts.AllowSafeBigIntegers},
internal/appengine/http.go:627:			opts.RetrieveMetadata = true
internal/appengine/http.go:633:	opts.DownsampleKey = q.Get("downsample_key")

Looking at the queue history, prior `docs-sync-appengine-*` tasks already added missing routes, status codes, and pagination/422 params. None ever added the `format`/`allow_bigintegers`/`allow_safe_bigintegers`/`retrieve_metadata`/`downsample_key` params. Let me verify these params actually have effect (not parsed-and-dropped) before proposing.
$ rg -n 'RetrieveMetadata|DownsampleKey' internal/appengine/downsample.go internal/appengine/data.go internal/appengine/service.go | rg -v '_test' | head
internal/appengine/data.go:53:	DownsampleKey        string
internal/appengine/data.go:54:	RetrieveMetadata     bool

$ rg -n 'QueryOpts|opts\.' internal/appengine/service.go | head -30
713:func (s *Service) readResolved(ctx context.Context, r *resolved, path string, opts QueryOpts) (any, error) {
716:		if opts.DownsamplePoints > 0 {
726:func (s *Service) GetDataByAlias(ctx context.Context, realm, alias, ifaceName, path string, opts QueryOpts) (any, error) {
736:func (s *Service) GetDataInGroup(ctx context.Context, realm, groupName, deviceID, ifaceName, path string, opts QueryOpts) (any, error) {

$ rg -n 'RetrieveMetadata|DownsampleKey|downsample_key' internal/ --glob '!**/*_test.go' | head -20
internal/appengine/data.go:53:	DownsampleKey        string
internal/appengine/data.go:54:	RetrieveMetadata     bool
internal/appengine/http.go:627:			opts.RetrieveMetadata = true
internal/appengine/http.go:633:	opts.DownsampleKey = q.Get("downsample_key")

$ cd /root/astrate-mule && rg -o '\bASTRATE_[A-Z_]+' -N internal/ | sort -u
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

`retrieve_metadata` and `downsample_key` are parsed (http.go:627,633) but have no consumer — a separate nuance. Let me also run the recipe's config-key check cheaply to see if there's drift there before I finalize proposals.
$ cd /root/astrate-mule && rg -o '\bASTRATE_[A-Z_]+' docs/site/configuration-reference.md 2>/dev/null | sort -u || echo "no config reference page"
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

The config-key check found two `ASTRATE_` keys in code that are absent from `docs/site/configuration-reference.md`. Let me confirm they're real env keys in config.go.
$ cd /root/astrate-mule && rg -n 'HOUSEKEEPING' internal/config/config.go
270:		"ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
271:		"HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
288:	if v, ok := os.LookupEnv("ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"); ok {

→ Read internal/config/config.go [limit=45, offset=255]
Both keys are real, functional env vars. This is config-key drift → escalation to `.mule/for-giulio.md` per the recipe. Let me check the file's format and current header.
→ Read .mule/for-giulio.md [limit=60]
The config-key drift is **already escalated** in for-giulio.md (line 46-52) — no duplicate needed. My one genuinely new finding is the undocumented data-query params. Let me check the `.mule/tasks/` detail-file style before writing my proposal.

$ ls /root/astrate-mule/.mule/tasks/ 2>/dev/null
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md

→ Read .mule/tasks/issue-27.md
The prior config-key escalation already covers my config finding. Let me verify one last response code — the data DELETE handler's success code — before finalizing.

$ cd /root/astrate-mule && rg -n 'func \(a \*API\) (putData|deleteData|publishData)' internal/appengine/http.go
292:func (a *API) putData(w http.ResponseWriter, r *http.Request) {
311:func (a *API) deleteData(w http.ResponseWriter, r *http.Request) {
321:func (a *API) putDataByAlias(w http.ResponseWriter, r *http.Request) {
340:func (a *API) putDataInGroup(w http.ResponseWriter, r *http.Request) {
359:func (a *API) deleteDataByAlias(w http.ResponseWriter, r *http.Request) {
369:func (a *API) deleteDataInGroup(w http.ResponseWriter, r *http.Request) {

→ Read internal/appengine/http.go [limit=65, offset=292]

# Todos
[✓] Pick a surface and find drift between code and docs/api/*.yaml
[✓] Compare documented paths vs actual handlers for 3-4 endpoints
[•] Write .mule/tasks/docs-sync-appengine-data-output-params.md detail file
[ ] Append task line to .mule/todo.md

All checks confirm the remaining drift is exactly one thing. Let me write the proposal: a detail file and one todo line.
← Write .mule/tasks/docs-sync-appengine-data-output-params.md
Wrote file successfully.

← Edit .mule/tasks/docs-sync-appengine-data-output-params.md
Index: /root/astrate-mule/.mule/tasks/docs-sync-appengine-data-output-params.md
===================================================================
--- /root/astrate-mule/.mule/tasks/docs-sync-appengine-data-output-params.md
+++ /root/astrate-mule/.mule/tasks/docs-sync-appengine-data-output-params.md
@@ -20,9 +20,9 @@
 All three handlers go through `parseQueryOpts` (http.go:254).
 
 ## What the code actually does (verified 2026-09-07)
 
-Three of the five miss ing params genuinely change the response and must be documented:
+Three of the five missing params genuinely change the response and must be documented:
 
 | param | parsed | consumed |
 |-------|--------|----------|
 | `format` (`structured`/`table`/`disjoint_tables`, default `structured`) | http.go:601-607 | data.go:237, 244 — `renderObject`/`renderIndividual` |

Now append the task line to `.mule/todo.md`, matching the `[auto]` tag style of prior proposal-recipe entries.
→ Read .mule/todo.md [offset=116]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -116,4 +116,5 @@
 - [x] docs-sync-appengine-group-patch-status [auto]: PATCH /groups/{group}/devices/{device} (internal/appengine/http.go:518-537, patchGroupDevice) shares applyPatch with the device-scoped PATCH and also answers 422 ("Attribute key not found", ErrAttributeKeyNotFound, service.go:352-357) and 409 ("Alias already in use", http.go:651), but docs/api/astarte_appengine_api.yaml lists only 200/400/401/404/500 for it — add the two missing responses.
 - [x] docs-sync-appengine-data-422-interface-level [auto]: the interface-level data GETs — GET /devices/{device}/interfaces/{interface}, GET /devices-by-alias/{alias}/interfaces/{interface}, GET /groups/{group}/devices/{device}/interfaces/{interface} — share parseQueryOpts (internal/appengine/http.go:554-639) and answer 422 FieldErrors for invalid since/since_after/to/limit/downsample_to/format params, but the spec documents only the path-level GET with 422 — add 422 to the interface-level GETs.
 - [!] appengine-snapshot-ignores-query-params [auto]: in `internal/appengine/data.go:138-148`, the individual-datastream interface-root snapshot branch ignores `Since`/`SinceAfter`/`To`/`Limit`/`Descending` even though the root GET is documented with all of them (docs/api/astarte_appengine_api.yaml:577-582) — the same function refuses incompatible params elsewhere (downsample 422, data.go:118-120/165-167), so this is a silent drop. Probe upstream on the Legion Go to decide honour-vs-refuse, then either apply the window/limit or answer an explicit 422 and stop advertising the params on the root GET; add a T1 test in data_test.go/query_opts_test.go pinning the chosen behaviour. — BLOCKED: wrote nothing
 - [!] appengine-group-token-roundtrip-test [auto]: add a container-free unit test for `parseGroupToken`/`groupTokenFor` (internal/appengine/service.go:602-623) — currently only exercised via the pg-backed groups_parity_test.go. Assert (1) round-trip `groupTokenFor(offset)` → `parseGroupToken` recovers offset, (2) a non-v1-nibble UUID is rejected, (3) a well-formed v1 UUID with zeroed time fields is accepted (parseGroupToken comment, service.go:614-616). — BLOCKED: lint failed: internal/appengine/groups_token_test.go:14:5: redefines-builtin-id: redefinition of the built-in function max (revive)
+- [ ] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

# Todos
[✓] Pick a surface and find drift between code and docs/api/*.yaml
[✓] Compare documented paths vs actual handlers for 3-4 endpoints
[✓] Write .mule/tasks/docs-sync-appengine-data-output-params.md detail file
[✓] Append task line to .mule/todo.md

Done: docs-sync proposal for surface appengine — appended one task line to .mule/todo.md; all path/method coverage in docs/api/astarte_appengine_api.yaml matches the routes in internal/appengine/http.go.

Drift found: the six data GET operations document only `since/since_after/to/limit/downsample_to/sort`, but `parseQueryOpts` (http.go:554-641) also accepts `format`, `allow_bigintegers`, `allow_safe_bigintegers` (all three change the response — data.go:230-244,403-405) plus `retrieve_metadata`/`downsample_key`, which are parsed but have no consumer. The detail file records the verified facts and the "three live, two dead" nuance so the executor doesn't document dead params as functional. Checked status codes against handlers too (200/201/204 all agree); the config-key check surfaced two missing keys but that escalation already exists verbatim in `.mule/for-giulio.md` (line 46-52), so no duplicate.

This was a proposal run only — no implementation, no git operations.

Files: .mule/todo.md, .mule/tasks/docs-sync-appengine-data-output-params.md
Verified: path/method sets compared (code http.go:43-80 vs yaml paths) -> pass; status codes for GET /groups/{group}, addGroupDevice, removeGroupDevice, put/delete data -> all match; `rg -o 'ASTRATE_[A-Z_]+'` code-vs-config-reference -> two missing keys, already escalated
Unsure: nothing
Follow-ups: `- [ ] docs-sync-appengine-data-output-params [auto]: ...` appended to .mule/todo.md
