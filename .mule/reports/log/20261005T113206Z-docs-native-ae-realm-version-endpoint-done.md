slug: docs-native-ae-realm-version-endpoint
verdict: done
at:  7331e6c
ran: 2026-10-05T11:32:06Z on DietPi in 344s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

$ rg -n "RealmName|Unauthorized|Forbidden|InternalServerError" docs/api/astrate_native_api.yaml | head -40
300:        - $ref: "#/components/parameters/RealmName"
324:          $ref: "#/components/responses/Unauthorized"
326:          $ref: "#/components/responses/Forbidden"
328:          $ref: "#/components/responses/InternalServerError"
368:          $ref: "#/components/responses/Unauthorized"
370:          $ref: "#/components/responses/InternalServerError"
374:    RealmName:
438:    Unauthorized:
439:      description: Unauthorized — missing or invalid JWT/token.
446:              detail: Unauthorized
448:    Forbidden:
449:      description: Forbidden — the JWT is valid but does not authorize this path.
456:              detail: Forbidden
458:    InternalServerError:

→ Read docs/api/astarte_pairing_api.yaml [offset=60, limit=45]

$ rg -n "^  /|^paths:|^  [a-z_]+:" docs/api/astrate_native_api.yaml
3:  title: Astrate Native API
4:  description: |
9:  version: "0.0.1"
10:  license:
23:paths:
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
373:  parameters:
382:  schemas:
437:  responses:

→ Read docs/api/astrate_native_api.yaml [offset=180, limit=180]

$ rg -n "^  /|^paths:|^  [a-z_]+:" docs/api/astarte_appengine_api.yaml
3:  title: Astarte AppEngine API
4:  description: |
9:  version: "0.0.1"
10:  license:
23:paths:
26:  /appengine/v1/{realm}/devices:
100:  /appengine/v1/{realm}/stats/devices:
140:  /appengine/v1/{realm}/devices/{device}:
218:  /appengine/v1/{realm}/devices-by-alias/{alias}:
306:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces:
340:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
387:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
584:  /appengine/v1/{realm}/devices/{device}/interfaces:
614:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
656:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
833:  /appengine/v1/{realm}/groups:
911:  /appengine/v1/{realm}/groups/{group}:
941:  /appengine/v1/{realm}/groups/{group}/devices:
1037:  /appengine/v1/{realm}/groups/{group}/devices/{device}:
1142:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces:
1172:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}:
1215:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}:
1404:  parameters:
1518:  schemas:
1723:  responses:

→ Read docs/api/astarte_realm_management_api.yaml [limit=50, offset=650]
→ Read docs/api/astarte_appengine_api.yaml [limit=80, offset=1330]
→ Read docs/api/astarte_appengine_api.yaml [limit=120, offset=1404]

$ rg -n "tags: \[" docs/api/astarte_appengine_api.yaml | sort -u -t: -k3 | head -20; echo ---; rg -n "mountRealmVersion|RequireRealmAny|ClaimAppEngine" cmd/astrate/main.go internal/auth/middleware.go | head; echo ---; rg -n "version" internal/observability/compat.go
313:      tags: [AppEngine - Data]
33:      tags: [AppEngine - Devices]
838:      tags: [AppEngine - Groups]
---
cmd/astrate/main.go:399:	mountRealmVersion(mux, mw)
cmd/astrate/main.go:429:// mountRealmVersion registers the realm-scoped version endpoints observed on
cmd/astrate/main.go:435:func mountRealmVersion(mux *http.ServeMux, mw *auth.Middleware) {
cmd/astrate/main.go:437:		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(realm.APICompatVersion)))
internal/auth/middleware.go:49:	return m.RequireRealmAny(claim)
internal/auth/middleware.go:52:// RequireRealmAny guards a realm-scoped route accepting ANY of the given
internal/auth/middleware.go:57:func (m *Middleware) RequireRealmAny(claims ...Claim) func(http.Handler) http.Handler {
---
36:// VersionHandler answers GET .../version with the upstream envelope
37:// {"data":"<version>"} (measured upstream 1.2.0, verify-versions.json).
38:func VersionHandler(version string) http.HandlerFunc {
39:	body, err := json.Marshal(map[string]string{"data": version})
49:// MountVersionCompat registers GET /{service}/version for one service.
50:func MountVersionCompat(mux *http.ServeMux, service, version string) {
51:	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))

→ Read internal/observability/compat.go
→ Read internal/auth/middleware.go [limit=100, offset=52]
→ Read cmd/astrate/main.go [limit=25, offset=425]

$ rg -n "APICompatVersion" --glob '!*_test.go' | head -20; echo ---; rg -n "^var version|version = \"|version =" cmd/astrate/main.go | head
cmd/astrate/main.go:437:		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(realm.APICompatVersion)))
cmd/astrate/main.go:438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
docs/COMPATIBILITY.md:253:    reports `realm.APICompatVersion` (currently `1.2.2`), a compatibility
docs/UPSTREAM-EXPERIMENTAL.md:13:`APICompatVersion` (internal/realm/service.go) must never claim a level higher
internal/realm/http.go:119:// feature UI on it — see APICompatVersion).
internal/realm/http.go:121:	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
internal/realm/service.go:580:// APICompatVersion is the upstream Realm Management API level Astrate
internal/realm/service.go:588:const APICompatVersion = "1.2.2"
---
53:var version = "0.1.0-dev"

$ rg -rn "docs-sync-native-version-value" .mule/ | head
.mule/todo.md:- [x] n [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:- [x] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pairing_api.yaml, which shows `version: "0.1.0-test"` (yaml:252) — a value no build can produce. The field is filled from `pairing.Config.Version` (internal/pairing/service.go:302), wired from `main.version` (cmd/astrate/main.go:151), whose default is `"0.1.0-dev"` (main.go:53) and which is otherwise set only via `-ldflags "-X main.version=..."`. Use `"0.1.0-dev"`, matching the four compat-version examples n already corrected to that (docs/api/astrate_native_api.yaml:203, 223, 243, 263). Do NOT document the service's own fallback `DefaultVersion = "0.1.0-astrate"` (service.go:86-88): it is unreachable from `main`, which always passes a non-empty value. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:- [ ] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value n settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/log.md:| 2026-09-17 | n | done | 503s | 08048b3 |
.mule/reports/log/20260917T194954Z-recipe-docs-sync-proposed.md:e33ab46 mule: log n
.mule/reports/log/20260917T194954Z-recipe-docs-sync-proposed.md:08048b3 mule: n [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/reports/log/20260925T201910Z-docs-sync-rm-delete-device-async-operation-param-done.md:- [x] n [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/reports/log/20260928T170059Z-docs-sync-pairing-deviceid-base64url-done.md:- [x] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pairing_api.yaml, which shows `version: "0.1.0-test"` (yaml:252) — a value no build can produce. The field is filled from `pairing.Config.Version` (internal/pairing/service.go:302), wired from `main.version` (cmd/astrate/main.go:151), whose default is `"0.1.0-dev"` (main.go:53) and which is otherwise set only via `-ldflags "-X main.version=..."`. Use `"0.1.0-dev"`, matching the four compat-version examples n already corrected to that (docs/api/astrate_native_api.yaml:203, 223, 243, 263). Do NOT document the service's own fallback `DefaultVersion = "0.1.0-astrate"` (service.go:86-88): it is unreachable from `main`, which always passes a non-empty value. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/reports/log/20260928T105915Z-recipe-docs-sync-proposed.md:196:- [x] n [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/re...
.mule/reports/log/20260928T105915Z-recipe-docs-sync-proposed.md:- [ ] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pairing_api.yaml, which shows `version: "0.1.0-test"` (yaml:252) — a value no build can produce. The field is filled from `pairing.Config.Version` (internal/pairing/service.go:302), wired from `main.version` (cmd/astrate/main.go:151), whose default is `"0.1.0-dev"` (main.go:53) and which is otherwise set only via `-ldflags "-X main.version=..."`. Use `"0.1.0-dev"`, matching the four compat-version examples n already corrected to that (docs/api/astrate_native_api.yaml:203, 223, 243, 263). Do NOT document the service's own fallback `DefaultVersion = "0.1.0-astrate"` (service.go:86-88): it is unreachable from `main`, which always passes a non-empty value. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

→ Read internal/realm/service.go [limit=30, offset=575]
→ Read cmd/astrate/main.go [limit=30, offset=385]
→ Read docs/api/astarte_appengine_api.yaml [limit=80, offset=1723]
→ Read docs/api/astarte_appengine_api.yaml [limit=40]

$ rg -n "version" docs/api/astrate_pairing_api.yaml | head -20
rg: docs/api/astrate_pairing_api.yaml: IO error for operation on docs/api/astrate_pairing_api.yaml: No such file or directory (os error 2)

$ rg -rn "realm/version|mountRealmVersion" --glob '*_test.go' | head -20
cmd/astrate/version_test.go:	n(mux, auth.NewMiddleware(nil))
cmd/astrate/version_test.go:	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/pairing/v1/somen", nil))
internal/observability/compat_version_test.go:	resp, err := http.Get(srv.URL + "/v1/somen")
internal/observability/compat_version_test.go:		t.Fatalf("/v1/somen: %v", err)
internal/observability/compat_version_test.go:		t.Errorf("/v1/somen = %d, want 200", resp.StatusCode)
internal/observability/compat_version_test.go:		t.Errorf("/v1/somen body = %s", body)

$ pwd; ls docs/api/ | head -20
/root/astrate-mule
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

→ Read cmd/astrate/version_test.go
→ Read docs/api/astarte_pairing_api.yaml [offset=1, limit=30]

$ rg -n "tags:|^- name:" docs/api/astarte_appengine_api.yaml | rg -v "tags: \[AppEngine" | head; echo ---; ls docs; cat docs/Makefile 2>/dev/null | head -30
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

$ rg -n "astarte_appengine_api|swagger" docs/swagger-ui/*.html docs/swagger-ui/* 2>/dev/null | head -20
docs/swagger-ui/index.html:7:  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5.32.11/favicon-32x32.png">
docs/swagger-ui/index.html:8:  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui.css">
docs/swagger-ui/index.html:32:    #swagger-ui { max-width: 100%; }
docs/swagger-ui/index.html:47:      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
docs/swagger-ui/index.html:53:  <div id="swagger-ui"></div>
docs/swagger-ui/index.html:55:  <script src="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui-bundle.js"></script>
docs/swagger-ui/index.html:56:  <script src="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui-standalone-preset.js"></script>
docs/swagger-ui/index.html:60:    const uiContainer = document.getElementById('swagger-ui');
docs/swagger-ui/index.html:66:        dom_id: '#swagger-ui',
docs/swagger-ui/index.html:7:  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5.32.11/favicon-32x32.png">
docs/swagger-ui/index.html:8:  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui.css">
docs/swagger-ui/index.html:32:    #swagger-ui { max-width: 100%; }
docs/swagger-ui/index.html:47:      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
docs/swagger-ui/index.html:53:  <div id="swagger-ui"></div>
docs/swagger-ui/index.html:55:  <script src="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui-bundle.js"></script>
docs/swagger-ui/index.html:56:  <script src="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui-standalone-preset.js"></script>
docs/swagger-ui/index.html:60:    const uiContainer = document.getElementById('swagger-ui');
docs/swagger-ui/index.html:66:        dom_id: '#swagger-ui',

$ rg -n "tags:" docs/api/astarte_pairing_api.yaml docs/api/astarte_realm_management_api.yaml | sed 's/.*tags: //' | sort -u
[Pairing]
[Realm Management - Config]
[Realm Management - Devices]
[Realm Management - Interfaces]
[Realm Management - Misc]
[Realm Management - Policies]
[Realm Management - Triggers]

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -20,8 +20,49 @@
 
 security: []
 
 paths:
+  # ── Version ─────────────────────────────────────────────────────────
+
+  /appengine/v1/{realm}/version:
+    get:
+      operationId: getVersion
+      summary: Get emulated API version
+      description: |
+        Returns the emulated upstream API compatibility level. The Astarte
+        Dashboard gates feature UI on this value.
+
+        The authorization path is exactly `version` — the realm segment is
+        stripped before the JWT grant is matched — so the token needs a grant
+        for GET on `version`. An unknown realm is answered 401, not 404: a
+        realm with no keys is indistinguishable from an unauthenticated
+        caller.
+      tags: [AppEngine - Misc]
+      security:
+        - a_aea: []
+      parameters:
+        - $ref: "#/components/parameters/RealmName"
+      responses:
+        "200":
+          description: API compatibility version.
+          content:
+            application/json; charset=utf-8:
+              schema:
+                type: object
+                required: [data]
+                properties:
+                  data:
+                    type: string
+                    description: Version string.
+              example:
+                data: "1.2.2"
+        "401":
+          $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
+        "500":
+          $ref: "#/components/responses/InternalServerError"
+
   # ── Devices ─────────────────────────────────────────────────────────
 
   /appengine/v1/{realm}/devices:
     get:

$ python3 -c "
import yaml,sys
for f in ['docs/api/astarte_appengine_api.yaml']:
    d=yaml.safe_load(open(f))
    p=d['paths']['/appengine/v1/{realm}/version']['get']
    print(p['operationId'], p['tags'], p['security'], [r for r in p['responses']])
    comps=set(d['components']['parameters'])|set(d['components']['responses'])
    import json
    s=json.dumps(d)
    import re
    refs=set(re.findall(r'\"#/components/(\w+)/([\w-]+)\"', s))
    print('missing refs:', [r for r in refs if r[1] not in comps.get(r[0],[])])
"
Traceback (most recent call last):
  File "<string>", line 2, in <module>
    import yaml,sys
ModuleNotFoundError: No module named 'yaml'

$ which yq swagger-cli redocly npx 2>/dev/null; pip list 2>/dev/null | head -5; python3 -c "import json;print('ok')"
/usr/local/bin/npx
ok

$ npx --yes @redocly/cli@latest lint /root/astrate-mule/docs/api/astarte_appengine_api.yaml 2>&1 | tail -30
|   ^
105 |     ids:
106 |       summary: Simple device ID list

referenced from ../../root/astrate-mule/docs/api/astarte_appengine_api.yaml:100:15 at #/paths/~1appengine~1v1~1{realm}~1devices/get/responses/200/content/application~1json; charset=utf-8 

Warning was generated by the no-invalid-media-type-examples rule.

Reference: https://redocly.com/docs/cli/rules/oas/no-invalid-media-type-examples

[8] ../../root/astrate-mule/docs/api/astarte_appengine_api.yaml:1018:17 at #/paths/~1appengine~1v1~1{realm}~1groups~1{group}~1devices/get/responses/200/content/application~1json; charset=utf-8/schema/oneOf

Schemas in `oneOf` must be mutually exclusive. Found overlapping schemas: `#/components/schemas/DataEnvelopeDeviceIDList` and `#/components/schemas/DataEnvelopeDeviceStatusList`. Both schemas define `links` without constraints that exclude each other. Add a discriminator, or constrain the shared properties to different values.

1016 | application/json; charset=utf-8:
1017 |   schema:
1018 |     oneOf:
     |     ^^^^^
1019 |       - $ref: "#/components/schemas/DataEnvelopeDeviceIDList"
1020 |       - $ref: "#/components/schemas/DataEnvelopeDeviceStatusList"

Warning was generated by the no-illogical-composition-keywords rule.

../../root/astrate-mule/docs/api/astarte_appengine_api.yaml: validated in 891ms

Woohoo! Your API description is valid. 🎉
You have 8 warnings.

$ npx --yes @redocly/cli@latest lint /root/astrate-mule/docs/api/astarte_appengine_api.yaml 2>&1 | rg "^\[[0-9]+\]"
[1] ../../root/astrate-mule/docs/api/astarte_appengine_api.yaml:101:17 at #/paths/~1appengine~1v1~1{realm}~1devices/get/responses/200/content/application~1json; charset=utf-8/schema/oneOf
[2] ../../root/astrate-mule/docs/api/astarte_appengine_api.yaml:1614:11 at #/components/schemas/DeviceStatus/properties/previous_interfaces/nullable
[3] ../../root/astrate-mule/docs/api/astarte_appengine_api.yaml:1615:11 at #/components/schemas/DeviceStatus/properties/previous_interfaces/description
[4] ../../root/astrate-mule/docs/api/astarte_appengine_api.yaml:1618:11 at #/components/schemas/DeviceStatus/properties/introspection/nullable
[5] ../../root/astrate-mule/docs/api/astarte_appengine_api.yaml:1619:11 at #/components/schemas/DeviceStatus/properties/introspection/description
[6] ../../root/astrate-mule/docs/api/astarte_appengine_api.yaml:101:17 at #/paths/~1appengine~1v1~1{realm}~1devices/get/responses/200/content/application~1json; charset=utf-8/schema
[7] ../../root/astrate-mule/docs/api/astarte_appengine_api.yaml:101:17 at #/paths/~1appengine~1v1~1{realm}~1devices/get/responses/200/content/application~1json; charset=utf-8/schema
[8] ../../root/astrate-mule/docs/api/astarte_appengine_api.yaml:1018:17 at #/paths/~1appengine~1v1~1{realm}~1groups~1{group}~1devices/get/responses/200/content/application~1json; charset=utf-8/schema/oneOf

$ make -C docs build 2>&1 | tail -20
make: Entering directory '/root/astrate-mule/docs'
cp DESIGN.md site/DESIGN.md
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'

$ git status --porcelain; python3 -c "
import json,sys
print(open('/root/astrate-mule/docs/site/api/astarte_appengine_api.yaml').read().count('/appengine/v1/{realm}/version'))
"
M docs/api/astarte_appengine_api.yaml
1

$ python3 -m pip --version 2>&1 | head -2
/usr/bin/python3: No module named pip

→ Read .mule/todo.md [offset=1, limit=15]

$ rg -n "docs-native-ae-realm-version-endpoint" .mule/todo.md
277:- [ ] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
278:- [ ] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -273,9 +273,9 @@
 - [!] store-register-inhibit-preserve: RegisterDevice must not overwrite status to 'registered' when a device was inhibited before first credentials request (preserve 'inhibited' state). File: internal/store/devices.go — BLOCKED: gates failed
 - [x] store-aliasvalues-self-exclusion-test: Add test for AliasValuesTaken to assert self-exclusion (same device's own alias values do not count as taken). File: internal/store/devices_test.go
 - [!] store-latestindividual-empty-errnotfound: Add store-level test asserting LatestIndividual returns ErrNotFound for empty series and returns newest on multiple samples. File: internal/store/datastreams_test.go — BLOCKED: TIMEOUT after 1201s — task too big, split it
 - [x] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.go
-- [ ] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ go vet ./... 2>&1 | tail -5; echo "vet=$?"; go test ./... 2>&1 | rg -v "^ok|no test files" | tail -15; gofmt -l . | tail -5
vet=0

→ Read docs/api/astarte_appengine_api.yaml [limit=48, offset=22]
Done: documented `GET /appengine/v1/{realm}/version` in `docs/api/astarte_appengine_api.yaml:26-63` (op `getVersion`, `security: [a_aea]`, `$ref` RealmName + Unauthorized/Forbidden/InternalServerError, inline `{data: string}`, no 404, 401-for-unknown-realm and authorization-path-`version` notes), and marked the todo line done.

Two deviations from the task line, both deliberate:
- The task said the served value is the build `version` var (`0.1.0-dev`). It is not: `cmd/astrate/main.go:437` passes `realm.APICompatVersion` (`internal/realm/service.go:588` = `"1.2.2"`), and `cmd/astrate/version_test.go:25` pins exactly that for the realm-scoped routes. The build version only feeds the unauthenticated per-service `/appengine/version` via `MountVersionCompat` (main.go:397). So the example is `data: "1.2.2"`, matching the pairing twin.
- The appengine spec had no `Misc` tag (only Devices/Data/Groups), so I introduced `AppEngine - Misc` mirroring RM's `Realm Management - Misc`. The specs carry no top-level `tags:` block, so nothing else needed.

Files: `docs/api/astarte_appengine_api.yaml`, `.mule/todo.md`
Verified: `go vet ./...` -> pass; `go test ./...` -> pass (no race build attempted on this Pi-per-rules, no code changed anyway); `gofmt -l .` -> clean; `npx @redocly/cli lint docs/api/astarte_appengine_api.yaml` -> valid, 8 warnings all pre-existing and none in the new operation; `make -C docs build` -> **failed at the mkdocs step**, `mkdocs: No such file or directory` and no pip in this Python, so the Swagger page itself was not rendered here. The `sync` half of the target did run, so `docs/site/api/astarte_appengine_api.yaml` carries the new path, and redocly parsing that same file is the closest available proof it still loads.
Unsure: the `AppEngine - Misc` tag name is my choice — Giulio may prefer folding the op into `AppEngine - Devices`. Also the docs-site render is unverified on this machine for the mkdocs reason above.
Follow-ups: none
