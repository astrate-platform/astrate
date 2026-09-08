slug: recipe-docs-sync
verdict: proposed
at:  95dc630
ran: 2026-09-08T18:21:08Z on DietPi in 711s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
I'll carry out the docs-sync recipe. It's a proposal job — I'll pick one surface, find drift, and propose fixes as new task lines in todo.md (or for-giulio.md for prose). Let me start by picking a surface and finding the endpoints.

Let me first look at the directory structure to understand the surfaces.

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

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*appengine*' 2>/dev/null | head -50
(no output)

I'll pick one surface per run. Let me start with appengine since it's likely the most substantive. Let me find the endpoints and documented paths.
→ Read docs/api/astarte_appengine_api.yaml
My `rg` for appengine routes returned nothing. Let me find how the routes are actually registered in the appengine package.

$ rg -n 'Get\(|Post\(|Delete\(|HandleFunc|Handle\(|/devices|/groups' internal/appengine/ 2>/dev/null | head -40; echo "==="; ls internal/appengine/
internal/appengine/mirror_test.go:25:	if rec := r.req(t, http.MethodPost, "/groups", body, r.token); rec.Code != http.StatusCreated {
internal/appengine/mirror_test.go:41:	decodeData(t, r.req(t, http.MethodGet, "/devices-by-alias/sensor-1/interfaces", "", r.token), &viaAlias)
internal/appengine/mirror_test.go:49:		"/groups/mg/devices/"+r.dev.String()+"/interfaces", "", r.token), &viaGroup)
internal/appengine/mirror_test.go:55:		{"unknown alias", "/devices-by-alias/nope/interfaces"},
internal/appengine/mirror_test.go:56:		{"unknown device", "/devices/" + unknownID(t) + "/interfaces"},
internal/appengine/mirror_test.go:68:	decodeData(t, r.reqCT(t, http.MethodPatch, "/devices-by-alias/sensor-1",
internal/appengine/mirror_test.go:74:	rec := r.reqCT(t, http.MethodPatch, "/devices-by-alias/sensor-1",
internal/appengine/mirror_test.go:83:	rec = r.reqCT(t, http.MethodPatch, "/devices-by-alias/nope",
internal/appengine/mirror_test.go:97:		"/devices-by-alias/sensor-1/interfaces/"+aeSensors+"/value", "", r.token), &viaAlias)
internal/appengine/mirror_test.go:108:		"/devices-by-alias/sensor-1/interfaces/"+aeSensors+"/value", "7.5", r.token); rec.Code != http.StatusOK {
internal/appengine/mirror_test.go:120:		"/devices-by-alias/sensor-1/interfaces/"+aeSensors+"/value", "", r.token); rec.Code != http.StatusNoContent {
internal/appengine/mirror_test.go:130:		"/devices-by-alias/nope/interfaces/"+aeSensors+"/value", "", r.token); rec.Code != http.StatusNotFound {
internal/appengine/mirror_test.go:155:		"/groups/"+gm+"/devices/"+r.dev.String(), "", r.token), &ds)
internal/appengine/mirror_test.go:160:	decodeData(t, r.reqCT(t, http.MethodPatch, "/groups/"+gm+"/devices/"+r.dev.String(),
internal/appengine/mirror_test.go:166:	rec := r.reqCT(t, http.MethodPatch, "/groups/"+gm+"/devices/"+r.dev.String(),
internal/appengine/mirror_test.go:177:		{"status", http.MethodGet, "/groups/" + gm + "/devices/" + stranger.String()},
internal/appengine/mirror_test.go:178:		{"data", http.MethodGet, "/groups/" + gm + "/devices/" + stranger.String() + "/interfaces/" + aeSensors + "/value"},
internal/appengine/mirror_test.go:188:	rec = r.reqCT(t, http.MethodPatch, "/groups/"+gm+"/devices/"+stranger.String(),
internal/appengine/mirror_test.go:197:		"/groups/"+gm+"/devices/"+r.dev.String()+"/interfaces/"+aeSensors+"/value", "", r.token), &inGroup)
internal/appengine/mirror_test.go:202:		"/groups/"+gm+"/devices/"+r.dev.String()+"/interfaces/"+aeSensors+"/value", "9.5", r.token); rec.Code != http.StatusOK {
internal/appengine/mirror_test.go:211:		"/groups/"+gm+"/devices/"+r.dev.String()+"/interfaces/"+aeSensors+"/value", "", r.token); rec.Code != http.StatusNoContent {
internal/appengine/mirror_test.go:226:		{"status", http.MethodGet, "/groups/nope/devices/" + r.dev.String()},
internal/appengine/mirror_test.go:227:		{"interfaces", http.MethodGet, "/groups/nope/devices/" + r.dev.String() + "/interfaces"},
internal/appengine/mirror_test.go:228:		{"data", http.MethodGet, "/groups/nope/devices/" + r.dev.String() + "/interfaces/" + aeSensors + "/value"},
internal/appengine/mirror_test.go:238:	rec = r.reqCT(t, http.MethodPatch, "/groups/nope/devices/"+r.dev.String(),
internal/appengine/stream/ws_test.go:81:		if ct := rec.Header().Get("Content-Type"); ct != "text/event-stream" {
internal/appengine/links_test.go:18:	if links.Self != "/v1/test/devices?details=true&limit=10" {
internal/appengine/links_test.go:50:	if links.Self != "/v1/test/devices" {
internal/appengine/links_test.go:58:	if links.Next != "/v1/test/devices?details=false&from_token=N&limit=5" {
internal/appengine/http_test.go:174:	return "/devices/" + r.dev.String() + suffix
internal/appengine/http_test.go:181:		if rec := r.req(t, http.MethodGet, "/devices", "", ""); rec.Code != http.StatusUnauthorized {
internal/appengine/http_test.go:184:		if rec := r.req(t, http.MethodGet, "/devices", "", r.wrong); rec.Code != http.StatusForbidden {
internal/appengine/http_test.go:198:		if rec := r.req(t, http.MethodGet, "/devices/"+unknownID(t), "", r.token); rec.Code != http.StatusNotFound {
internal/appengine/http_test.go:258:		path := "/devices/" + dev.String() + "/interfaces/" + aeSensors + "/value"
internal/appengine/http_test.go:316:			"/devices/"+dev.String()+"/interfaces/"+aeObj+"/12?downsample_to=5", "", r.token); rec.Code != http.StatusUnprocessableEntity {
internal/appengine/http_test.go:324:			"/devices/"+dev.String()+"/interfaces/"+aeSensors+"?downsample_to=5", "", r.token); rec.Code != http.StatusUnprocessableEntity {
internal/appengine/http_test.go:340:			"/devices/"+dev.String()+"/interfaces/"+aeSensors+"/big?downsample_to=5", "", r.token), &empty)
internal/appengine/http_test.go:424:		if rec := r.req(t, http.MethodPost, "/groups", body, r.token); rec.Code != http.StatusCreated {
internal/appengine/http_test.go:428:		decodeData(t, r.req(t, http.MethodGet, "/groups", "", r.token), &names)
internal/appengine/http_test.go:433:		decodeData(t, r.req(t, http.MethodGet, "/groups/g1/devices", "", r.token), &devs)
===
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

Let me find where the appengine routes are actually registered (the handler registration), not the test files.
$ rg -n 'appengine|/appengine' internal/ --glob '*.go' -l | grep -v _test; echo "==="; rg -n 'HandleFunc|\.Get\(|\.Post\(|\.Delete\(|\.Patch\(|\.Put\(|Handle\(' internal/appengine/*.go | grep -v _test
internal/observability/compat.go
internal/httpx/notfound.go
internal/appengine/service.go
internal/appengine/downsample.go
internal/appengine/channels/protocol.go
internal/appengine/channels/ws.go
internal/appengine/http.go
internal/appengine/data.go
===
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

Now I have the full route inventory for appengine. Let me read http.go:43-80 fully and also check the devices-by-alias PATCH route registration, since the doc shows a PATCH on `/devices-by-alias/{alias}` but I need to verify all registered routes vs. documented ones.

Let me list the full set of registered routes:

From http.go:
- GET /devices
- GET /stats/devices
- GET /devices/{device}
- PATCH /devices/{device}
- GET /devices-by-alias/{alias}
- GET /devices/{device}/interfaces
- (line 49 missing from output - let me check)
- GET /devices/{device}/interfaces/{interface}
- GET/PUT/POST/DELETE /devices/{device}/interfaces/{interface}/{path...}
- PATCH /devices-by-alias/{alias}
- GET /devices-by-alias/{alias}/interfaces
- GET /devices-by-alias/{alias}/interfaces/{interface}
- GET/PUT/POST/DELETE /devices-by-alias/{alias}/interfaces/{interface}/{path...}
- GET/POST /groups
- GET /groups/{group}
- GET/POST /groups/{group}/devices
- GET/PATCH/DELETE /groups/{group}/devices/{device}
- GET /groups/{group}/devices/{device}/interfaces
- GET /groups/{group}/devices/{device}/interfaces/{interface}
- GET/PUT/POST/DELETE /groups/{group}/devices/{device}/interfaces/{interface}/{path...}

Wait, the doc lists a PATCH on `/devices-by-alias/{alias}` (line 235) and `/devices/{device}` (line 156). Let me check line 49 and confirm what was between 48 and 50.

Let me check http.go lines 43-82 carefully.
→ Read internal/appengine/http.go [offset=40, limit=45]
✱ Grep "retrieve_metadata|downsample_key" in docs/api/astarte_appengine_api.yaml · 0 matches
Route sets match 1:1 between code and spec — no missing or extra routes. Now let me spot-check a few endpoints for status codes and response fields.

I noticed `retrieve_metadata` and `downsample_key` query params in the code (line 624, 633). Let me check if those are documented.
→ Read internal/appengine/http.go [offset=555, limit=90]
Found two undocumented query params: `retrieve_metadata` and `downsample_key`. Now let me spot-check a few endpoints for status code correctness.
→ Read internal/appengine/http.go [limit=130, offset=84]
→ Read internal/appengine/http.go [limit=100, offset=214]
→ Read internal/appengine/http.go [limit=100, offset=414]
Now let me check the remaining endpoint handlers to compare status codes. Let me also check the `patchDevice` handler for 409 handling, and the `deleteData` for status code.

Looking at what I've already found:
- `patchDeviceByAlias` returns 400 for wrong Content-Type (line 231), but the doc says 400 too ✓
- `patchDevice` returns 500 for wrong Content-Type (line 210), but the doc says 400 ✗

Wait, let me re-read the doc for `patchDevice`:
Line 191: `"400": $ref: "#/components/responses/BadRequest"` — but the code returns 500 (unmapped) when Content-Type is wrong.

Actually, looking more carefully at the code comment (line 206-208): "Upstream honors PATCH only with exactly this content type... Anything else fails its merge with :patch_mimetype_not_supported — which no FallbackController clause maps, so Phoenix surfaces an unmapped 500. Replicated as-is."

So the code intentionally returns 500 for wrong Content-Type on `patchDevice` (matching upstream behavior). But the doc only lists 400, 401, 404, 409, 422, 500 as possible responses. The 500 is listed, so technically it's covered. But the doc doesn't distinguish that the 500 can occur from a Content-Type mismatch — it's just listed as "InternalServerError". Actually, the doc does list 500. And there's no 400 returned for Content-Type on the `/devices/{device}` PATCH path — only for decode failures. Let me check the doc again...

Line 191 says:
```
"400": $ref: "#/components/responses/BadRequest"
```

In the code, `patchDevice` returns `WriteBadRequest(w)` only when `DecodeData` fails (line 215), not for Content-Type. For Content-Type mismatch, it returns 500 (line 210). So 400 is not incorrect — it IS returned for a malformed body. And 500 is also listed. The doc doesn't claim 500 comes from Content-Type specifically, just that 500 is a possible response. This seems fine.

Now the important findings:

1. **`retrieve_metadata` query param is undocumented** — exists in code at http.go:624, not in the OpenAPI spec. This is used on all data GET endpoints.

2. **`downsample_key` query param is undocumented** — exists in code at http.go:633, not in the OpenAPI spec.

Let me check if `retrieve_metadata` is a real query param that affects behavior...

Looking at http.go:624-631:
```go
if v := q.Get("retrieve_metadata"); v != "" {
    switch v {
    case "true":
        opts.RetrieveMetadata = true
    case "false":
    default:
        fe.addf("retrieve_metadata", "is invalid")
    }
}
```

And line 633:
```go
opts.DownsampleKey = q.Get("downsample_key")
```

Both are parsed by `parseQueryOpts` which is used in `serveData`, which serves `getData`, `getDataByAlias`, and `getDataInGroup`. So these query params are on all the interface data GET endpoints. They're not in the doc.

Let me also check the QueryOpts struct to understand what these do.</think>
✱ Grep "type QueryOpts|RetrieveMetadata|DownsampleKey" in internal/appengine · 11 matches
→ Read internal/appengine/http.go [limit=100, offset=314]
→ Read internal/appengine/http.go [limit=55, offset=642]
→ Read internal/appengine/http.go [limit=42, offset=513]
→ Read internal/appengine/data.go [limit=60, offset=40]
Let me check what these two params do, to write accurate proposals.
✱ Grep "DownsampleKey|RetrieveMetadata" in internal/appengine · 10 matches
✱ Grep "\.DownsampleKey|\.RetrieveMetadata" in internal · 6 matches

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

Let me check the config-key sweep and the docs/site configuration reference to see if there's config drift.
$ rg -o '\bASTRATE_[A-Z_]+' -N docs/site/ 2>/dev/null | sort -u
docs/site/configuration-reference.md:ASTRATE_DATABASE_DSN
docs/site/configuration-reference.md:ASTRATE_ENGINE_SHARDS
docs/site/configuration-reference.md:ASTRATE_HTTP_ADDR
docs/site/configuration-reference.md:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
docs/site/configuration-reference.md:ASTRATE_HTTP_TLS_CERT_FILE
docs/site/configuration-reference.md:ASTRATE_HTTP_TLS_KEY_FILE
docs/site/configuration-reference.md:ASTRATE_LOG_FORMAT
docs/site/configuration-reference.md:ASTRATE_LOG_LEVEL
docs/site/configuration-reference.md:ASTRATE_MASTER_KEY
docs/site/configuration-reference.md:ASTRATE_MASTER_KEY_FILE
docs/site/configuration-reference.md:ASTRATE_MQTT_ADDR
docs/site/configuration-reference.md:ASTRATE_MQTT_ADVERTISED_URL
docs/site/configuration-reference.md:ASTRATE_MQTT_INSECURE_DEV_MODE
docs/site/configuration-reference.md:ASTRATE_MQTT_SESSION_STORE_PATH
docs/site/configuration-reference.md:ASTRATE_MQTT_TLS_CERT_FILE
docs/site/configuration-reference.md:ASTRATE_MQTT_TLS_KEY_FILE
docs/site/configuration-reference.md:ASTRATE_REALM_JWT_PUBLIC_KEY
docs/site/configuration-reference.md:ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
docs/site/configuration-reference.md:ASTRATE_REALM_NAME
docs/site/configuration-reference.md:ASTRATE_SECURITY_MASTER_KEY_FILE
docs/site/contributing.md:ASTRATE_TEST_DSN
docs/site/deployment.md:ASTRATE_DATABASE_DSN
docs/site/deployment.md:ASTRATE_ENGINE_SHARDS
docs/site/deployment.md:ASTRATE_HTTP_ADDR
docs/site/deployment.md:ASTRATE_LOG_LEVEL
docs/site/deployment.md:ASTRATE_MASTER_KEY
docs/site/deployment.md:ASTRATE_MASTER_KEY_FILE
docs/site/deployment.md:ASTRATE_MQTT_ADDR
docs/site/observability.md:ASTRATE_LOG_LEVEL
docs/site/operations.md:ASTRATE_DATABASE_DSN
docs/site/operations.md:ASTRATE_ENGINE_SHARDS
docs/site/operations.md:ASTRATE_HTTP_ADDR
docs/site/operations.md:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
docs/site/operations.md:ASTRATE_LOG_LEVEL
docs/site/operations.md:ASTRATE_MASTER_KEY
docs/site/operations.md:ASTRATE_MASTER_KEY_FILE
docs/site/operations.md:ASTRATE_MQTT_ADDR
docs/site/operations.md:ASTRATE_MQTT_INSECURE_DEV_MODE
docs/site/operations.md:ASTRATE_REALM_NAME
docs/site/troubleshooting.md:ASTRATE_DATABASE_DSN
docs/site/troubleshooting.md:ASTRATE_LOG_LEVEL

$ rg -n 'HOUSEKEEPING_(DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION|REALM_DELETION_DISABLED)|ASTRATE_HOUSEKEEPING' internal/config/config.go
270:		"ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
271:		"HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
288:	if v, ok := os.LookupEnv("ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"); ok {

Two housekeeping config keys are in code but absent from `configuration-reference.md`. Let me verify they're real env vars config.go reads.
$ rg -n 'maxValueBytes|ASTRATE_FLOW_CONFIG' internal/appengine/*.go | grep -v _test
internal/appengine/http.go:21:// maxValueBytes is upstream's per-value limit (~64 KiB): a larger REST
internal/appengine/http.go:24:const maxValueBytes = 65536
internal/appengine/http.go:298:	if len(value) > maxValueBytes {
internal/appengine/http.go:327:	if len(value) > maxValueBytes {
internal/appengine/http.go:346:	if len(value) > maxValueBytes {

$ rg -n 'ASTRATE_FLOW_CONFIG' internal/flow/ | grep -v _test; echo "==="; ls .mule/
internal/flow/blocks/info.go:116:		Config:  "image (required), config (object→ASTRATE_FLOW_CONFIG), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",
internal/flow/blocks/schema.go:19:	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as ASTRATE_FLOW_CONFIG"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
internal/flow/blocks/container/block.go:45:	Nested        map[string]any // opaque JSON for ASTRATE_FLOW_CONFIG
internal/flow/blocks/container/docker.go:19:	// FlowConfigJSON is passed as ASTRATE_FLOW_CONFIG.
internal/flow/blocks/container/docker.go:99:		args = append(args, "-e", "ASTRATE_FLOW_CONFIG="+spec.FlowConfigJSON)
internal/flow/blocks/container/docker.go:201:// encodeFlowConfigJSON marshals nested config for ASTRATE_FLOW_CONFIG.
===
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

Let me verify the device status response fields against the code, and check the device stats keys.
✱ Grep "type DeviceStatus|device_id|first_registration|credentials_inhibited" in internal/appengine/service.go · 67 matches
→ Read internal/appengine/service.go [limit=30, offset=105]
The spec documents `DeviceStatus.device_id` but the struct serializes `json:"id"`. Let me verify how this renders on the wire — there may be a custom marshaler.

$ rg -n 'MarshalJSON|device_id|"id"' internal/appengine/service.go | head -30
106:	ID                   string                        `json:"id"`
633:		return FieldErrors{"device_id": {"is not a valid device id"}}

Let me verify the actual wire shape from tests — does the JSON use `id` or `device_id`?
✱ Grep ""id"|`json:"(id|device_id)"`" in internal/appengine · 2 matches
→ Read internal/appengine/e2e_test.go [limit=30, offset=300]
Let me check how device status is actually serialized to the wire — look at the mirror test that decodes a device status.
→ Read internal/appengine/mirror_test.go [limit=45, offset=140]
✱ Grep "\.ID\b|ID:\s|device_id"|"device_id"" in internal/appengine · 100 matches
✱ Grep "deviceStatusView|func.*DeviceStatus.*MarshalJSON|mirrorStatusView" in internal/appengine · 4 matches
→ Read internal/appengine/service.go [limit=50, offset=250]
✱ Grep "`json:"id"`|json:"id"|\.ID\b.*\bid\b" in internal/appengine · 2 matches
The `DeviceStatus` struct serializes the device ID as `"id"` (service.go:106), while the spec documents `device_id`. Let me verify what the wire emits and what the dashboard tests decode.
→ Read internal/appengine/dashboard_compat_test.go [limit=75, offset=85]
The `DeviceStatus` struct confirms a schema drift. Let me verify a couple more specifics: the introspection wire shape, and whether `listGroupDevices` can actually return 422.
✱ Grep "introspectionView" in internal/appengine · 3 matches
→ Read internal/appengine/service.go [limit=45, offset=555]
→ Read internal/appengine/service.go [limit=15, offset=833]

$ ls docs/; echo "---"; cat docs/Makefile 2>/dev/null | head -30
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
---
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

→ Read .mule/todo.md
Let me check the docs build so I can reference it correctly, and look at the Swagger UI/current todo.md and for-giulio.md.
→ Read .mule/for-giulio.md
The route sets match 1:1. Now let me check the previous docs-sync task detail file so I don't re-propose what was already decided about `retrieve_metadata`/`downsample_key`.

$ ls .mule/tasks/ | grep -i appengine; echo "---"
docs-sync-appengine-data-output-params.md
---

→ Read .mule/tasks/docs-sync-appengine-data-output-params.md
The previous run already decided `retrieve_metadata`/`downsample_key` are dead params (do not document). Let me verify the spec's `downsample_to` minimum and the 422 responses on the PUT/POST data endpoints against the code.

$ rg -n 'downsample_to|must be >= 2|downsam' docs/api/astarte_appengine_api.yaml | head; echo "==="; rg -n '"422"|UnprocessableEntity|reserved|Value size|Response' docs/api/astarte_appengine_api.yaml | rg -i '422|value' | head
325:        optional query parameters for time-range filtering and downsampling on
574:        parameters for time-range filtering and downsampling on datastream
1085:        optional query parameters for time-range filtering and downsampling on
1353:      name: downsample_to
1357:      description: Downsample to this many data points (must be >= 2).
===
200:        "422":
282:        "422":
359:        "422":
409:        "422":
603:        "422":
646:        "422":
831:        "422":
907:        "422":
1025:        "422":
1115:        "422":

The previous run (2026-09-07) already worked `retrieve_metadata`/`downsample_key` (dead params — decided not to document them), and the config-key drift is already escalated in for-giulio.md. Let me confirm `WriteData` doesn't transform the device status, and pin down the exact spec line numbers for my findings.
$ rg -n 'func WriteData\b|func WriteDataWithLinks' internal/httpx/*.go $(go list -f '{{.Dir}}' github.com/astarte-platform/astarte-go 2>/dev/null) 2>/dev/null | head; echo "==="; rg -rn 'func WriteData\b|func WriteDataWithLinks' internal/ --glob '*.go' | head
===

✱ Grep "func WriteData" in . · 7 matches
→ Read pkg/astarteapi/envelope.go [limit=40, offset=120]
Let me check whether `docs/site/` has an appengine page that could contradict the code (prose escalation vs docs/api fix), since the DeviceStatus schema is the key finding.

$ ls docs/site/; echo "==="; rg -ln 'device_id|DeviceStatus|introspection' docs/site/ 2>/dev/null
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
===
docs/site/pairing-and-security.md
docs/site/data-modeling.md
docs/site/troubleshooting.md
docs/site/appengine-api.md
docs/site/triggers.md
docs/site/architecture.md
docs/site/realm-management-api.md
docs/site/mqtt-protocol.md
docs/site/migration-from-astarte.md
docs/site/compatibility.md
docs/site/pokemon-agent.md
docs/site/observability.md
docs/site/json-payload-profile.md

$ rg -n 'device_id|introspection|data:' docs/site/appengine-api.md | head -40
21:GET /appengine/v1/<realm>/devices/<device_id>
24:Returns device status: introspection, connected flag, stats (total received msgs/bytes), timestamps, aliases, attributes.
29:PATCH /appengine/v1/<realm>/devices/<device_id>
47:GET /appengine/v1/<realm>/devices/<device_id>/interfaces/<interface>[/<path>]
55:GET /appengine/v1/<realm>/devices/<device_id>/interfaces/<interface>/<path>
70:PUT /appengine/v1/<realm>/devices/<device_id>/interfaces/<interface>/<path>
79:DELETE /appengine/v1/<realm>/devices/<device_id>/interfaces/<interface>/<path>
91:DELETE /appengine/v1/<realm>/groups/<name>/devices/<device_id>

→ Read docs/site/appengine-api.md [limit=45, offset=15]
Now I have a clear picture. Three genuine `docs/api` drift findings for the appengine surface; the config-key and dead-param items are already handled. Let me verify the spec's introspection example and the `DeviceStatus` field name claims once more against the spec text I read, then append the proposals.

The findings:

1. **`DeviceStatus` schema wrong**: spec documents `device_id` (yaml:1409, examples yaml:67/75/83-84) but the wire emits `id` (service.go:106, `json:"id"`; mirror/dashboard tests decode from it); `introspection` documented as `map<string,integer>` (yaml:1438-1441) but emits `{major,minor}` objects (service.go:123-127); six fields the code emits (`groups`, `total_received_msgs`, `total_received_bytes`, `first_credentials_request`, `last_seen_ip`, `previous_interfaces`; service.go:105-121) are missing from the schema.
2. **422 missing on PUT/POST data-set ops**: `putData`/`putDataByAlias`/`putDataInGroup` return 422 "Value size exceeds size limits" for bodies > 64 KiB (http.go:298-300, 327-329, 346-348), but the spec lists only 200/400/401/404/500 for all six PUT/POST ops.
3. **`downsample_to` minimum wrong**: spec says "(must be >= 2)" (yaml:1357) but code rejects `<= 2` with "must be greater than 2" (http.go:592-593).

Let me append the proposals to `.mule/todo.md`.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -123,4 +123,7 @@
 - [x] detailed-listing-required-encrypted-flags: emit the `required`/`encrypted` mapping flags in `writeDetailedMapping` (internal/realm/detailed.go) so `GET /realmmanagement/v1/{realm}/interfaces?detailed=true` round-trips them — the detailed listing currently materializes a fixed field set that drops both flags even when a datastream mapping declares one, while the plain GET preserves them via the stored definition. Datastream branch only (properties reject both at parse); add a T1 assertion in internal/realm/detailed_test.go.
 
 - [x] store-alias-values-taken-test: add a store test for `AliasValuesTaken` (internal/store/devices.go:343, the `alias_already_in_use` ownership check behind PATCH /devices) pinning the three branches that currently have no direct coverage — no other device in the realm carries any of the values → false; another device carries one → true; the same device (excluded via `d.id <> $2`) carries the value → false (must not count itself). Add to internal/store/devices_test.go. [auto]
 - [x] store-latest-individual-test: add a store test for `LatestIndividual` (internal/store/datastreams.go:247, the previous-value lookup feeding value_change/path_created triggers) — only exercised through the engine's fakeStore today; assert (a) of two samples on one series, `ORDER BY ts DESC LIMIT 1` returns the newest with the right value column, (b) an empty series yields ErrNotFound. Add to internal/store/datastreams_test.go. [auto]
+- [ ] docs-sync-appengine-device-status-schema [auto]: fix the `DeviceStatus` schema in docs/api/astrate_appengine_api.yaml — the wire emits `id` (internal/appengine/service.go:106 `json:"id"`, decoded as `ds.ID` in mirror_test.go:156 and dashboard_compat_test.go:100/149) but the spec names the field `device_id` (yaml:1409) in both the schema and the list/detail examples (yaml:67, 75); `introspection` is documented as `map<string,integer>` (yaml:1438-1441) but each value is a `{major, minor}` object (introspectionEntry, service.go:123-127, introspectionView service.go:833-836); and the schema omits the six fields `groups`, `total_received_msgs`, `total_received_bytes`, `first_credentials_request`, `last_seen_ip` (omitempty), `previous_interfaces` (omitempty) that deviceStatusView emits (service.go:254-279). Rename `device_id`→`id`, correct `introspection` to `map<string,{major:int,minor:int}>`, add the missing fields, and fix the examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-appengine-data-set-422 [auto]: add the missing `422` response to the six PUT/POST data-set operations in docs/api/astrate_appengine_api.yaml — `putData`, `putDataByAlias`, `putDataInGroup` (internal/appengine/http.go:292-309, 321-338, 340-357) return 422 `{"errors":{"detail":"Value size exceeds size limits"}}` when the body value exceeds 65536 bytes (maxValueBytes, http.go:24), both for PUT and its POST alias (routes registered at http.go:52-53, 63-64, 78-79), but the spec lists only 200/400/401/404/500 for all six operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-appengine-downsample-min [auto]: fix the `downsample_to` constraint in docs/api/astrate_appengine_api.yaml — the DataDownsample description says "(must be >= 2)" (yaml:1357) but parseQueryOpts rejects `downsample_to <= 2` with "must be greater than 2" (internal/appengine/http.go:592-593), so the documented minimum is off by one: readings of `2` (allowed per the doc) actually 422. Change the description to "(must be > 2)" (the 422 error shape is already documented). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

Now let me verify the todo append and confirm I haven't duplicated any existing entries.
→ Read .mule/todo.md [offset=124]
The three proposals are appended. Nothing was implemented and git was untouched. Surface checked: appengine.

**Findings**
- Routes match 1:1 (internal/appengine/http.go:43-81 vs all 32 documented paths) — no missing or extra routes.
- `DeviceStatus` schema drift: spec documents `device_id`, wire emits `id`; `introspection` documented as `map<string,integer>` but each value is a `{major,minor}` object; six emitted fields undocumented.
- PUT/POST data-set ops all return 422 for bodies > 64 KiB, spec lists only 200/400/401/404/500.
- `downsample_to` minimum off by one: doc says `>= 2`, code rejects `<= 2`.
- `retrieve_metadata`/`downsample_key`: parsed but dead — the 2026-09-07 run already decided not to document them (`.mule/tasks/docs-sync-appengine-data-output-params.md`), so not re-proposed.
- Config-key sweep: the two `ASTRATE_HOUSEKEEPING_*` gaps are already escalated in `.mule/for-giulio.md:47-52` — not duplicated.

Done: docs-sync — checked the appengine surface; appended 3 docs/api drift proposals to .mule/todo.md
Files: .mule/todo.md (3 added lines, none touched elsewhere)
Verified: rg route inventory + yaml path read -> match; handler reads (service.go:105-127,254-279, http.go:24,292-357,592-593) -> drift confirmed
Unsure: nothing
Follow-ups: docs-sync-appengine-device-status-schema, docs-sync-appengine-data-set-422, docs-sync-appengine-downsample-min (all [auto])
