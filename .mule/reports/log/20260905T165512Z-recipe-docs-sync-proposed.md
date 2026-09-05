slug: recipe-docs-sync
verdict: proposed
at:  9a0da50
ran: 2026-09-05T16:55:12Z on DietPi in 176s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
I'll pick one surface for this run. Let me start with **pairing** — it's likely smaller and faster. Let me grep the code and docs in parallel.

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*pairing*'
(no output)

$ rg -n '^\s{2}/' docs/api/astrate_pairing_api.yaml 2>/dev/null || echo "FILE NOT FOUND"
FILE NOT FOUND

$ ls docs/api/*.yaml 2>/dev/null; echo "---"; ls internal/ 2>/dev/null
docs/api/astarte_appengine_api.yaml
docs/api/astarte_housekeeping_api.yaml
docs/api/astarte_pairing_api.yaml
docs/api/astarte_realm_management_api.yaml
docs/api/astrate_native_api.yaml
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

No pairing surface found. Let me check what surfaces actually exist in the code and docs.
$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*.go' | head -80
internal/observability/compat.go:20:	mux.HandleFunc("GET /"+service+"/health", func(w http.ResponseWriter, r *http.Request) {
internal/observability/compat.go:51:	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
internal/observability/health.go:43:	mux.HandleFunc("GET /astrate/v1/health", h.handleHealth)
internal/observability/health.go:44:	mux.HandleFunc("GET /astrate/v1/readiness", h.handleReadiness)
internal/observability/compat_version_test.go:31:		if ct := resp.Header.Get("Content-Type"); ct == "" {
internal/observability/compat_version_test.go:39:	mux.HandleFunc("GET /v1/{realm}/version", VersionHandler("x"))
internal/flow/blocks/httpblocks.go:197:		Subtype:   mediaType(resp.Header.Get("Content-Type")),
internal/flow/blocks/httpblocks.go:330:	if req.Header.Get("Content-Type") == "" {
internal/flow/blocks/httpblocks_test.go:281:		contentType: r.Header.Get("Content-Type"),
internal/flow/blocks/httpblocks_test.go:438:	if v := got[0].header.Get("X-Test"); v != "yes" {
internal/swagger/swagger.go:20:	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
internal/engine/forward/http_test.go:42:		gotCT = r.Header.Get("Content-Type")
internal/engine/forward/http_test.go:43:		gotRealm = r.Header.Get("Astarte-Realm")
internal/engine/forward/http_test.go:44:		gotTrigger = r.Header.Get("Astrate-Trigger-Name")
internal/engine/forward/http_test.go:88:		if v := r.Header.Get("X-Custom"); v != "custom-val" {
internal/engine/forward/http_test.go:91:		if v := r.Header.Get("X-Other"); v != "other-val" {
internal/engine/forward/http_test.go:95:		if v := r.Header.Get("Content-Type"); v != "application/json" {
internal/engine/forward/http_test.go:98:		if v := r.Header.Get("Astarte-Realm"); v != "r" {
internal/engine/forward/http_test.go:101:		if v := r.Header.Get("Astrate-Trigger-Name"); v != "t" {
internal/engine/triggers/actions_test.go:363:		gotCT = r.Header.Get("Content-Type")
internal/engine/triggers/actions_test.go:364:		gotRealm = r.Header.Get("Astarte-Realm")
internal/engine/triggers/actions_test.go:365:		gotFoo = r.Header.Get("X-Foo")
internal/httpx/notfound_test.go:13:	mux.HandleFunc("GET /appengine/v1/{realm}/devices", func(w http.ResponseWriter, _ *http.Request) {
internal/httpx/notfound_test.go:19:	mux.HandleFunc("GET /housekeeping/v1/realms", func(w http.ResponseWriter, _ *http.Request) {
internal/httpx/notfound_test.go:119:	mux.HandleFunc("GET /appengine/v1/{path...}", func(w http.ResponseWriter, _ *http.Request) {
internal/httpx/cors.go:25:			origin := r.Header.Get("Origin")
internal/httpx/cors.go:39:			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
internal/appengine/stream/ws.go:152:	return r.Header.Get("Accept") == "text/event-stream" || r.URL.Query().Get("transport") == "sse"
internal/appengine/http.go:209:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:230:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:521:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/pairing/http.go:81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/pairing/http.go:86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
internal/pairing/http.go:88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
internal/pairing/http.go:90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
internal/pairing/http.go:344:	header := r.Header.Get("Authorization")
internal/auth/middleware.go:179:	header := r.Header.Get("Authorization")

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/pairing/http.go
81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
344:	header := r.Header.Get("Authorization")

Found pairing routes. Let me get the full list and the documented paths.
$ rg -n '^\s{2}/' docs/api/astarte_pairing_api.yaml
27:  /pairing/v1/{realm}/agent/devices:
77:  /pairing/v1/{realm}/agent/devices/{deviceID}:
100:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
154:  /pairing/v1/{realm}/devices/{deviceID}:
192:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:

$ rg -rn 'agent/devices' internal/ --include='*.go' | head -20
rg: unrecognized flag --include

similar flags that are available: --include-zero

Interesting — the docs have `/agent/devices` paths not in the code, and the code has a `/health` path not in the docs. Let me dig deeper.
→ Read docs/api/astarte_pairing_api.yaml
→ Read internal/pairing/http.go

$ rg -n 'writeServiceError|WriteError|WriteBadRequest|WriteNotFound|WriteForbidden|WriteUnauthorized|WriteInternalServerError|WriteFieldErrors|WriteData|StatusCreated|StatusNoContent|StatusOK|StatusForbidden|StatusNotFound|StatusUnauthorized|StatusServiceUnavailable|StatusTooManyRequests|StatusUnprocessableEntity' internal/pairing/http.go
111:		_ = astarteapi.WriteError(w, http.StatusTooManyRequests, detailTooManyRequests)
117:		_ = astarteapi.WriteBadRequest(w)
121:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
128:		a.writeServiceError(w, err)
131:	_ = astarteapi.WriteData(w, http.StatusCreated, registerResponse{CredentialsSecret: secret})
137:		a.writeServiceError(w, err)
140:	w.WriteHeader(http.StatusNoContent)
159:		_ = astarteapi.WriteError(w, http.StatusTooManyRequests, detailTooManyRequests)
165:		_ = astarteapi.WriteUnauthorized(w)
170:		_ = astarteapi.WriteBadRequest(w)
174:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
181:		a.writeServiceError(w, err)
184:	_ = astarteapi.WriteData(w, http.StatusCreated, credentialsResponse{ClientCrt: clientCrt})
204:		_ = astarteapi.WriteUnauthorized(w)
209:		a.writeServiceError(w, err)
217:	_ = astarteapi.WriteData(w, http.StatusOK, resp)
245:		_ = astarteapi.WriteUnauthorized(w)
250:		_ = astarteapi.WriteBadRequest(w)
254:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
261:		a.writeServiceError(w, err)
265:		_ = astarteapi.WriteData(w, http.StatusOK, verifyValidResponse{
272:	_ = astarteapi.WriteData(w, http.StatusOK, verifyInvalidResponse{
288:		_ = astarteapi.WriteData(w, http.StatusOK, map[string]string{"status": "ok"})
290:		_ = astarteapi.WriteNotFound(w)
292:		_ = astarteapi.WriteError(w, http.StatusServiceUnavailable, "unhealthy")
296:// writeServiceError maps service errors onto upstream statuses and bodies.
305:func (a *API) writeServiceError(w http.ResponseWriter, err error) {
308:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
311:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
314:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
317:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
320:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
323:		_ = astarteapi.WriteUnauthorized(w)
325:		_ = astarteapi.WriteForbidden(w)
329:		_ = astarteapi.WriteInternalServerError(w)

Good — the code actually registers 6 routes, the YAML has 5. Let me now spot-check status codes and fields for a few endpoints. Let me also check the verify endpoint's response and the register endpoint's status codes more closely.
$ rg -n 'WriteDeviceNotFound|WriteForbidden|WriteNotFound|WriteUnauthorized|WriteInternalServerError|WriteBadRequest|WriteFieldErrors|WriteError|WriteData|WriteTooManyRequests' internal/astarteapi/ 2>/dev/null || rg -rn 'func Write' internal/ --glob '*.go' | head -30
(no output)

$ rg -rn 'func Write|StatusNotFound|StatusForbidden|StatusUnauthorized|StatusInternalServerError|StatusBadRequest|StatusUnprocessableEntity|StatusTooManyRequests' internal/pkg/ 2>/dev/null; rg -rn 'func WriteDeviceNotFound|func WriteForbidden|func WriteNotFound|func WriteUnauthorized|func WriteInternalServerError|func WriteBadRequest' internal/ --glob '*.go' | head -20
(no output)

$ rg -rln 'func WriteData|func WriteDeviceNotFound|func WriteNotFound|func WriteForbidden' internal/ pkg/ --glob '*.go' 2>/dev/null
pkg/astarteapi/envelope.go:ln(w http.ResponseWriter, status int, v any) error {
pkg/astarteapi/envelope.go:lnWithLinks(w http.ResponseWriter, status int, v any, links Links) error {
pkg/astarteapi/envelope.go:lnWithMetadata(w http.ResponseWriter, status int, v any, m map[string]any) error {
pkg/astarteapi/envelope.go:ln(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:ln(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:ln(w http.ResponseWriter) error {

→ Read pkg/astarteapi/envelope.go
Now I have a comprehensive picture. Let me compare the documented spec against the actual code for the pairing surface.

**Code routes (internal/pairing/http.go:81-91):**
1. `GET /pairing/v1/{realm}/health` → handleHealth (200/404/503)
2. `POST /pairing/v1/{realm}/agent/devices` → handleRegister
3. `DELETE /pairing/v1/{realm}/agent/devices/{deviceID}` → handleUnregister
4. `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials` → handleCredentials
5. `GET /pairing/v1/{realm}/devices/{deviceID}` → handleInfo
6. `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify` → handleVerify

**Documented paths (yaml:27-263):** 5 paths, all matching except:
- **`/health` is not documented in the YAML at all** (present in code, absent in docs).

Let me now check the verify endpoint's field detail. The code's `verifyInvalidResponse` has `Cause` mapped to `cause` etc. The YAML documents `cause, details, timestamp, valid` — matches. Let me verify the register endpoint's status codes against the doc.

Looking at code's `writeServiceError` (http.go:305-330), the register endpoint can return:
- 422 (`ErrInvalidHWID`, `ErrInvalidPayloadFormat`, `ErrAlreadyRegistered`, `ErrRegistrationLimitReached`)
- 401 (`ErrUnauthorized`)
- 403 (`ErrInhibited`)
- 404 (`store.ErrNotFound`)
- 500

The YAML documents for register (lines 54-75): 201, 400, 401, 403, 422, 429, 500. All present. **Missing 404** — but the code doesn't produce a 404 for register either (writeServiceError's ErrNotFound would be... but for a new registration, the realm would need to exist). Actually, `store.ErrNotFound` from `svc.Register` — the realm not found case. That's a valid 404 the code provides but the docs don't list. That's a minor drift.

Now for the **verify endpoint** (lines 192-263): documented 200, 400, 401, 403, 404, 422, 500. The code produces the same.

Let me check the **agent endpoints' auth**: the code guards `/agent/devices` with `requireAgent` (a_pa claim JWT) — matches doc's `a_pa` security. Device endpoints use bearerSecret — matches doc's `bearerSecret`.

Now let me check the register endpoint's `initial_payload_format` — documented as "Optional Astrate extension". Code has it. ✓.

Let me check the health endpoint more — it's the clear gap. Let me also check the "Four proposals maximum" limit (recipe says "Five proposals maximum").

Let me now also check the `PairingInfo` schema in the YAML vs the code's `infoResponse`. Code: `protocols.astarte_mqtt_v1` with `broker_url` and `ca_crt`, plus `status`, `version`. YAML PairingInfo has same. ✓.

Let me now compile the drift findings. Let me check the status codes once more for the getDeviceInfo/flow C endpoint. Doc documents 200, 401, 403, 404, 500. Code: 200 (WriteData), 401 (Unauthorized), then writeServiceError → 401/403/404/500, plus 422? No — getDeviceInfo has no body decode, so no 400/422. Matches doc (no 400/422 documented). ✓.

Now, the health endpoint is undocumented. That's the main drift. Let me check whether there are documentation tasks and for-giulio escalations.

Let me also verify one endpoint's actual status codes against a doc claim — the register 404. Actually the doc doesn't list 404 for register, and code CAN emit 404 via `store.ErrNotFound` → `WriteDeviceNotFound`. Let me confirm realm-not-found maps to ErrNotFound in Register.

Actually, let me keep this focused (max 5 proposals). The clear, defensible findings for the pairing surface:

1. **`GET /pairing/v1/{realm}/health` is present in code but absent in the spec.** This is a docs/api fix task (add the path).

Let me look for other concrete drift. Let me check the "device_registration_limit_reached" and "already_registered" field-name in the doc — documented under ValidationError examples (lines 521-532). Code emits these same error_name fields. ✓

Let me verify the register 404 more concretely, and check whether the spec is missing any documented-but-absent items. All 5 documented paths exist in code. So no "documented but absent" drift — only "present but undocumented" for `/health`.

Let me check whether the config-key sweep needs to run for pairing. The recipe says config keys go to for-giulio, but that's a separate surface of the recipe. Since I picked pairing as the surface, I'll focus there but the conf key check is part of the same recipe. Let me keep focused on pairing per instructions to pick one surface.

Let me now double check the docs/site for any pairing contradiction. Actually the recipe says site prose → for-giulio. Let me check if there's a page that contradicts. But to stay within the surface and cheap, I'll focus on the api yaml which has clear drift.

Let me verify the `/health` path exists in the site docs or swagger references.

$ rg -rn 'pairing/v1.*health|/health' docs/site/ 2>/dev/null | head; echo "---mkdocs swagger refs---"; rg -rn 'pairing' docs/mkdocs.yml docs/*.yml 2>/dev/null | head
docs/site/troubleshooting.md:- **Health check:** `curl localhost:8080/astrate/v1n` — returns 200 if the process is alive.
docs/site/observability.md:| `GET /astrate/v1n` | Liveness probe. Returns 200 if the process is running. |
docs/site/quickstart.md:curl http://localhost:8080/astrate/v1n
docs/site/compatibility.md:Upstream's unauthenticated `GET /{appengine,realmmanagement,pairing}n`, which the Astarte Dashboard polls for its per-service status indicators, returns a static `200` -- the indicator cannot go red whatever state the instance is in. Astrate serves the same route and the same `200` envelope, but runs its database probe first and answers `503` when it fails. Astrate additionally keeps the realm-scoped `GET /n`, which upstream 404s: it resolves the realm too, and devices can probe it before they hold credentials.
docs/site/operations.md:`/astrate/v1n` is liveness; `/astrate/v1/readiness` pings the database
docs/site/deployment.md:| `:8080` | HTTP/HTTPS | REST API: `/n,readiness,metrics}` |
---mkdocs swagger refs---
docs/mkdocs.yml:  - Pairing & Security: n-and-security.md
docs/mkdocs.yml:  - Pairing & Security: n-and-security.md

Interesting! It seems `rg` output is replacing some text with "n" — the pattern matches may be mangled by the terminal (there appears to be some artifact). Let me look at the docs/site pairing page to check for contradictions, and check the swagger config.
$ rg -rln 'pairing' docs/site/ | head; echo "---"; ls docs/site/ | grep -i -E 'pair|sec'
docs/site/pairing-and-security.md:- Revocation: new credentials record new serial; broker rejects certs whose serial differs from latest (when `ln.enforce_latest_cert` is enabled).
docs/site/pairing-and-security.md:POST /ln/v1/<realm>/agent/devices
docs/site/pairing-and-security.md:POST /ln/v1/<realm>/devices/<device_id>/protocols/astarte_mqtt_v1/credentials
docs/site/pairing-and-security.md:GET /ln/v1/<realm>/devices/<device_id>
docs/site/pairing-and-security.md:Agent                Astrate(ln)          Device                 Astrate(broker)
docs/site/troubleshooting.md:1. The device trusts the per-realm CA returned by `GET /ln/v1/<realm>/devices/<id>` (the `ca_crt` field). Verify this CA matches what you expect.
docs/site/troubleshooting.md:**Cause:** Client certs expire after `ln.cert_ttl` (default 30 days).
docs/site/troubleshooting.md:1. The SDK should call `POST /ln/v1/<realm>/devices/<id>/protocols/astarte_mqtt_v1/credentials/verify` on boot. If it returns `"valid": false, "cause": "EXPIRED"`, the SDK re-CSRs automatically.
docs/site/troubleshooting.md:2. If the SDK doesn't re-CSR, check that `ln.cert_ttl` is not set too short for your device reboot frequency.
docs/site/troubleshooting.md:**Cause:** `ln.enforce_latest_cert = true` and the device presented an older certificate.
---
pairing-and-security.md

There's a display artifact — "pairing" is being rendered as "ln" and "/v1/" as "/ln/v1/". This is a terminal rendering issue in my view, but the actual file content is fine. Let me read the actual site page (it's on the never-touch list, so only read for contradictions to escalate).
→ Read docs/site/pairing-and-security.md
The site page documents Flow A/B/C but not the `/health` endpoint, and doesn't contradict code. No site contradiction to escalate for pairing.

Now let me verify the YAML cvf build references the pairing spec correctly, and check the docs build to confirm my proposal wouldn't break anything (per recipe, if I touch docs/api the executing task must run docs build). Let me look at the docs Makefile and swagger references.

$ rg -rn 'pairing' docs/ --include='*.yml' --include='*.yaml' --include='Makefile' 2>/dev/null | grep -iv 'site/' | head -20
(no output)

$ rg -rn 'astarte_pairing|swagger' docs/mkdocs.yml 2>/dev/null; echo "---makefile---"; cat docs/Makefile 2>/dev/null; echo "---mkdocs yaml refs---"; rg -n 'yaml|api' docs/mkdocs.yml | head
- API Explorer: n.md
---makefile---
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
---mkdocs yaml refs---
50:  - AppEngine API: appengine-api.md
51:  - Realm Management API: realm-management-api.md
52:  - Housekeeping API: housekeeping-api.md

The pairing API spec isn't referenced from mkdocs nav directly (only appengine/realmmanagement/housekeeping appear as API pages), but the `sync` copies the whole `docs/api/` to site/api. Let me check the API Explorer / swagger page that references the pairing yaml.

The recipe's "verify the site builds if you touch docs/api/..." applies to the executing task. Since the doc work here is a *proposal* (a task line for someone to execute later), the docs build check becomes a follow-up instruction for the executing task. I should include it in the proposed task description.

Let me check what swagger UI references to know if adding a path would break anything — but since this is a proposal, I just need to note that the executing task must run `make build`.

Let me now also double-check the config-key sweep direction, since the spec documents some pairing config. Actually, the recipe lists the config-key drift check as part of the same recipe but I picked "one surface per run". I selected pairing. I'll keep focus on the pairing API spec drift.

Let me finalize my drift findings for the pairing surface:

**Findings:**

1. **`GET /pairing/v1/{realm}/health` — present in code (http.go:81), absent in the spec** (`astarte_pairing_api.yaml`). This is a docs/api fix task. Clear, defensible.

2. **Register endpoint `POST /agent/devices` — code can return 404 (realm not found → `store.ErrNotFound` → `WriteDeviceNotFound` at http.go:326-327), but the YAML (lines 54-75) documents only 201/400/401/403/422/429/500.** Missing 404. Minor but real drift.

Let me confirm register can actually produce a 404. `a.svc.Register` — if the realm is unknown, `store.ErrNotFound`. This is plausible. But let me not over-assert. Let me check the service.
$ rg -rn 'func.*Register' internal/pairing/*.go | grep -v _test
internal/pairing/service.go:n(ctx context.Context, realmName, hwID, initialFormat string) (string, error) {
internal/pairing/http.go:n(w http.ResponseWriter, r *http.Request) {

The `rg` output is mangling text again (showing "n" instead of "Register(handle"). This is a tool output display quirk. Let me read the service Register directly.
✱ Grep "func .*Register\(" in internal/pairing · 3 matches
→ Read internal/pairing/service.go [offset=174, limit=60]
`Register` calls `s.st.GetRealmByName(ctx, realmName)` — if the realm doesn't exist, that returns `store.ErrNotFound`, and `writeServiceError` maps it to 404 `WriteDeviceNotFound`. So yes, register can return 404, which the YAML doesn't document. Same logic applies to credentials and verify endpoints (all call GetRealmByName) — but the YAML does document 404 for those. Only register's 404 is missing.

Now I have a solid, defensible set of drift findings for the **pairing** surface. Let me settle the proposals (max five; I have two solid API-spec ones):

1. **`/health` path absent from spec** — add `GET /pairing/v1/{realm}/health` to `astarte_pairing_api.yaml`. Present in code at http.go:81. Returns 200 `{"data":{"status":"ok"}}`, 404, 503. Unauthenticated (no security scheme).
2. **Register 404 undocumented** — add `404` (DeviceNotFound) to the register path's responses. Code path: service.go:183-186 → `store.ErrNotFound` → http.go:326-327 → `WriteDeviceNotFound`.

These two are genuine, mechanical drift. Both are `docs/api/` fixes (touchable, proposable as tasks). Both are "present but undocumented" / "response missing" — clean proposals.

The recipe says pick one surface, check three or four endpoints. I did. Five proposals max. I have two solid ones. I'll present those two.

Let me look at .mule/todo.md and .mule/for-giulio.md to append properly.
→ Read .mule/todo.md
I understand the format now. The recipe is a proposal job, and the deliverable is new task lines appended to `.mule/todo.md`.

Note: the todo.md at line 21-24 says `mule/queue` is authoritative and edits should go there. But my task instructions explicitly say to append new task lines to `.mule/todo.md`. Let me check where mule/queue is.

The instruction from the user says "its output is new task lines appended to .mule/todo.md". I'll follow the explicit instruction. However, I should check if `.mule/queue` exists as a separate file, since the todo.md itself says that branch is authoritative. But the user's instruction is clear: append to `.mule/todo.md`. Let me check whether there's a queue file to decide.

$ ls -la .mule/ 2>/dev/null; echo "---git branch---"; git branch --show-current 2>/dev/null; git status --short 2>/dev/null | head
total 240
drwxr-xr-x  7 root root   4096 Sep  5 13:25 .
drwxr-xr-x 16 root root   4096 Sep  4 21:26 ..
-rw-r--r--  1 root root     13 Sep  5 18:52 .budget
-rw-r--r--  1 root root    364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root     11 Sep  5 13:25 .heartbeat
-rw-r--r--  1 root root  20480 Sep  5 18:54 .last-output
-rw-r--r--  1 root root      0 Sep  5 18:52 .lock
-rw-r--r--  1 root root     10 Sep  5 18:52 .rotation
-rw-r--r--  1 root root   7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root   2762 Aug 31 01:45 config
-rw-r--r--  1 root root 108624 Sep  5 18:54 cron.log
drwxr-xr-x  2 root root   4096 Sep  5 13:14 failed
-rw-r--r--  1 root root  13842 Sep  5 13:24 for-giulio.md
-rw-r--r--  1 root root   1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root   1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root   1084 Sep  5 13:25 log.md
-rw-r--r--  1 root root   8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root   4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root   4096 Sep  4 21:09 reports
drwxr-xr-x  2 root root   4096 Sep  4 20:46 reviews
drwxr-xr-x  2 root root   4096 Sep  4 20:46 tasks
-rw-r--r--  1 root root  10864 Sep  5 13:25 todo.md
---git branch---
mule/queue

I'm on branch `mule/queue` already, so editing `.mule/todo.md` here is correct (this is the authoritative branch where the mule works). Good — that aligns with the todo.md note and the user's instruction.

Now let me append the two proposal task lines. Both are `[auto]` proposal tasks (nobody approved them). Let me create the proposal lines.

The two proposals:

1. **docs-sync pairing: add undocumented `/health` path** — `GET /pairing/v1/{realm}/health` is present in code (internal/pairing/http.go:81) but absent from docs/api/astarte_pairing_api.yaml. Add it: returns 200 `{"data":{"status":"ok"}}`, 404 (unknown realm), 503 (unhealthy), and is unauthenticated (no security scheme, no rate limit). Executing task must run `make -C docs build` to confirm the Swagger-referenced YAML still loads.

2. **docs-sync pairing: register endpoint missing 404** — `POST /pairing/v1/{realm}/agent/devices` can return 404 (DeviceNotFound) when `GetRealmByName` fails for an unknown realm (service.go:183-186 → http.go:326-327), but the YAML documents only 201/400/401/403/422/429/500.

Both are `docs/api/*.yaml` fixes (touchable, proposable as tasks, not prose).

Let me append them to the queue. Per the format, new proposal tasks get tagged `[auto]`. Let me append these two lines at the end of the Queue section.
→ Read .mule/todo.md [offset=96, limit=4]
I'll append my two proposal lines at the end.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -96,4 +96,6 @@
 - [x] purge-properties-compression-capability [auto]: upstream v1.3.0 adds a `purge_properties_compression_format` device capability (`zlib`|`plaintext`, default `zlib`) — a wire-visible capability value. Check whether Astrate's capabilities handling (internal/broker, the `<realm>/<device_id>/capabilities` topic, issue #16) needs to recognize/honour it, or whether zlib-only is already the deliberate default; propose the change or note why not needed.
 - [!] empty-introspection-verification [auto]: upstream v1.3.0 changed "allow devices with empty introspection" — verify whether Astrate's device connection/introspection handling currently rejects an empty introspection string where upstream now accepts it, and propose a fix if so. — BLOCKED: wrote nothing
 - [!] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages" — investigate only: does Astrate have an analogous delay between trigger install and first delivery? Report, do not patch. — BLOCKED: wrote nothing
 - [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
+- [ ] docs-sync-pairing-health-path [auto]: add the undocumented `GET /pairing/v1/{realm}/health` route to docs/api/astarte_pairing_api.yaml — it exists in code (internal/pairing/http.go:81) and returns 200 `{"data":{"status":"ok"}}` for an existing realm, 404 for an unknown realm, 503 when unhealthy, unauthenticated and unrate-limited. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-pairing-register-404 [auto]: document the missing `404` (DeviceNotFound) on `POST /pairing/v1/{realm}/agent/devices` in docs/api/astarte_pairing_api.yaml — the handler returns it when `GetRealmByName` fails for an unknown realm (internal/pairing/service.go:183-186 → http.go:326-327), but the spec lists only 201/400/401/403/422/429/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

Now let me verify the deliverable. The recipe also mentions escalating docs/site contradictions to for-giulio.md — I checked the pairing-and-security.md page and found no contradiction (it documents flows A/B/C consistently with code; it just doesn't mention health, which is a docs/api gap, not a site contradiction). No for-giulio escalation needed.

Let me do a final review of what I proposed and confirm nothing in my proposals touches the never-touch list improperly. `docs/api/*.yaml` is explicitly touchable per the recipe. `docs/site/` prose is untouched. Good.

I have not touched git, not implemented the tasks, and have not created any issue (the recipe does not ask for `gh issue create` — it asks for task lines, which I've done).

Let me verify the edit is clean.
→ Read .mule/todo.md [limit=4, offset=99]
The queue is updated. The docs-sync recipe is carried out: I picked the **pairing** surface, compared code routes against the spec, spot-checked endpoints/status codes/fields, and proposed two concrete `docs/api/` fixes. No docs/site contradictions found (the site page agrees with code on flows A/B/C), so no for-giulio escalation. I did not implement anything, did not touch git, and created no GitHub issues (the recipe doesn't ask for them).

```
Done: Compared pairing surface (code vs docs/api/astarte_pairing_api.yaml) and appended two docs-api gap proposals to .mule/todo.md
Files: .mule/todo.md
Verified: rg route extraction (internal/pairing/http.go:81-91), spec path list (astarte_pairing_api.yaml), status-code cross-check vs handler+envelope helpers -> pass
Unsure: nothing
Follow-ups: docs-sync-pairing-health-path, docs-sync-pairing-register-404 (both appended above)
```

Findings in brief:
- `GET /pairing/v1/{realm}/health` (http.go:81) is in code but absent from the spec — propose adding it.
- Register's `404` (unknown realm -> `store.ErrNotFound` -> WriteDeviceNotFound, service.go:183/326) is undocumented — propose adding it.

I deliberately did **not** flag the auth note line in the spec's `a_pa`/`bearerSecret` descriptions since they match the requireAgent/bearerSecret code. I also checked the site's payload examples against the spec (broker_url, ca_crt, status enum) — consistent, no contradiction to escalate.
