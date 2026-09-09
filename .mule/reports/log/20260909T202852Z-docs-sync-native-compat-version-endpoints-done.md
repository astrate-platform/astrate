slug: docs-sync-native-compat-version-endpoints
verdict: done
at:  d148f16
ran: 2026-09-09T20:28:52Z on DietPi in 480s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

# Todos
[•] Read docs/api/astrate_native_api.yaml to understand existing health endpoint pattern
[ ] Read observability code to confirm response shape
[ ] Add four /{service}/version endpoints to the YAML
[ ] Run make -C docs build to verify YAML loads
[ ] Report

→ Read docs/api/astrate_native_api.yaml
Now let me verify the code implementation of these endpoints.
✱ Glob "**/observability/compat.go" 1 match
→ Read cmd/astrate/main.go [offset=380, limit=30]
→ Read internal/observability/compat.go
Let me check how the docs build validates the YAML and how APICompatVersion is defined.
✱ Grep "APICompatVersion" 6 matches
→ Read docs/Makefile
→ Read docs/mkdocs.yml

$ ls /root/astrate-mule/docs/api /root/astrate-mule/docs/swagger-ui
/root/astrate-mule/docs/api:
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

/root/astrate-mule/docs/swagger-ui:
index.html

$ python3 -c "import yaml; print('pyyaml ok')" 2>&1; which mkdocs 2>&1; pip list 2>/dev/null | grep -i mkdocs
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; print('pyyaml ok')
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

→ Read docs/requirements.txt
→ Read docs/swagger-ui/index.html

$ make -C docs build 2>&1 | head -30
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

Let me try the docs build and check available tooling.
✱ Grep "yaml" 4 matches
mkdocs isn't installed. Let me try installing it, and meanwhile validate the YAML parses with the Go yaml library.

$ pip install -q "mkdocs>=1.6,<2" 2>&1 | tail -5; echo "---"; python3 -c "import mkdocs; print(mkdocs.__version__)" 2>&1
/bin/sh: 1: pip: not found
---
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import mkdocs; print(mkdocs.__version__)
    ^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'mkdocs'

$ which python3 pip3 python uv 2>&1; python3 --version 2>&1; ls /usr/lib/python3*/dist-packages 2>/dev/null | head; ls /root/.venv 2>/dev/null; ls /opt 2>/dev/null
/usr/bin/python3
Python 3.13.5
Cryptodome
packaging
packaging-25.0.dist-info
pycryptodomex-3.20.0.egg-info
dietpi-dashboard
pihole
rustdesk
scale-analytics

$ python3 -m pip --version 2>&1; python3 -m venv --help >/dev/null 2>&1 && echo "venv ok" || echo "no venv"
/usr/bin/python3: No module named pip
venv ok

$ timeout 15 getent hosts pypi.org >/dev/null 2>&1 && echo "pypi resolves" || echo "pypi no resolve"; timeout 15 curl -sI https://pypi.org/simple/ -o /dev/null -w "%{http_code}\n" 2>&1 | tail -1
pypi resolves
200

Network is available. Setting up a venv to run the docs build properly.
$ python3 -m venv /tmp/mule-docs-venv && /tmp/mule-docs-venv/bin/pip install -q -r requirements.txt 2>&1 | tail -3; /tmp/mule-docs-venv/bin/mkdocs --version
mkdocs, version 1.6.1 from /tmp/mule-docs-venv/lib/python3.13/site-packages/mkdocs (Python 3.13)

mkdocs ready. Let me check the build works at baseline before my change.
$ make -C docs build 2>&1 | tail -15
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

The Makefile needs mkdocs in PATH. Running with the venv prepended.
$ PATH="/tmp/mule-docs-venv/bin:$PATH" make -C docs build 2>&1 | tail -15
 │  × Closed contribution model – community members can't report bugs
 │  × Currently unlicensed – unsuitable for production use
 │
 │  Our full analysis:
 │
 │  https://squidfunk.github.io/mkdocs-material/blog/2026/02/18/mkdocs-2.0/

INFO    -  Cleaning site directory
INFO    -  Building documentation to directory: /root/astrate-mule/site-dist
INFO    -  The following pages exist in the docs directory, but are not included in the "nav" configuration:
  - DESIGN.md
  - ROADMAP.md
INFO    -  Doc file 'index.md' contains an unrecognized relative link 'api/', it was left as is.
INFO    -  Documentation built in 5.29 seconds
make: Leaving directory '/root/astrate-mule/docs'

Baseline build passes. Making the YAML changes now.
← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -179,8 +179,82 @@
             example:
               data:
                 status: unhealthy
 
+# ── Compat version (upstream parity) ────────────────────────────────
+
+/appengine/version:
+  get:
+    operationId: getAppEngineVersion
+    summary: AppEngine compat version
+    description: |
+      Upstream-parity unauthenticated version endpoint. Answers the
+      AppEngine service's API compatibility version.
+    tags: [Compat Health]
+    responses:
+      "200":
+        description: Service's API compatibility version.
+        content:
+          application/json:
+            schema:
+              $ref: "#/components/schemas/DataEnvelopeVersion"
+            example:
+              data: "1.0.0"
+
+/realmmanagement/version:
+  get:
+    operationId: getRealmManagementVersion
+    summary: Realm Management compat version
+    description: |
+      Upstream-parity unauthenticated version endpoint. Answers the
+      Realm Management service's API compatibility version.
+    tags: [Compat Health]
+    responses:
+      "200":
+        description: Service's API compatibility version.
+        content:
+          application/json:
+            schema:
+              $ref: "#/components/schemas/DataEnvelopeVersion"
+            example:
+              data: "1.2.2"
+
+/pairing/version:
+  get:
+    operationId: getPairingVersion
+    summary: Pairing compat version
+    description: |
+      Upstream-parity unauthenticated version endpoint. Answers the
+      Pairing service's API compatibility version.
+    tags: [Compat Health]
+    responses:
+      "200":
+        description: Service's API compatibility version.
+        content:
+          application/json:
+            schema:
+              $ref: "#/components/schemas/DataEnvelopeVersion"
+            example:
+              data: "1.0.0"
+
+/housekeeping/version:
+  get:
+    operationId: getHousekeepingVersion
+    summary: Housekeeping compat version
+    description: |
+      Upstream-parity unauthenticated version endpoint. Answers the
+      Housekeeping service's API compatibility version.
+    tags: [Compat Health]
+    responses:
+      "200":
+        description: Service's API compatibility version.
+        content:
+          application/json:
+            schema:
+              $ref: "#/components/schemas/DataEnvelopeVersion"
+            example:
+              data: "1.0.0"
+
 # ── Native WebSocket (live stream) ──────────────────────────────────
 
 /astrate/v1/{realm}/socket:
   get:

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -403,8 +403,16 @@
         status:
           type: string
           enum: [ok, unhealthy]
 
+DataEnvelopeVersion:
+  type: object
+  required: [data]
+  properties:
+    data:
+      type: string
+      description: The service's API compatibility version.
+
 ErrorDetail:
   type: object
   required: [errors]
   properties:

Now adding the `DataEnvelopeVersion` schema.
Let me verify the actual version string passed to `MountVersionCompat`.
✱ Grep "version\s*=|version :=|VersionCompat" in cmd/astrate · 3 matches
→ Read cmd/astrate/main.go [limit=10, offset=48]

$ cat /root/astrate-mule/test/conformance/upstream/verify-versions.json 2>/dev/null; echo "---"; rg -l "appengine/version|verify-versions" /root/astrate-mule/test/conformance 2>/dev/null
{
  "recorded_at": "2026-08-24T14:30:00Z",
  "astarte_version": "1.2.0",
  "realm": "zzprobe08",
  "recorder_command": "scratch probes (fase 1 verify batch, issues #69/#77); see verify-versions.transcript.txt",
  "observations": [
    {
      "name": "service-root version, housekeeping", "why": "#77: the dashboard polls one version endpoint per service at startup",
      "method": "GET", "path": "/housekeeping/version", "auth": "none",
      "status": 200, "body": "{\"data\":\"1.2.0\"}"
    },
    {
      "name": "service-root version, appengine", "why": "#77",
      "method": "GET", "path": "/appengine/version", "auth": "none",
      "status": 200, "body": "{\"data\":\"1.2.0\"}"
    },
    {
      "name": "service-root version, realmmanagement", "why": "#77",
      "method": "GET", "path": "/realmmanagement/version", "auth": "none",
      "status": 200, "body": "{\"data\":\"1.2.0\"}"
    },
    {
      "name": "service-root version, pairing", "why": "#77",
      "method": "GET", "path": "/pairing/version", "auth": "none",
      "status": 200, "body": "{\"data\":\"1.2.0\"}"
    },
    {
      "name": "global version route absent", "why": "#77: bounds the surface — there is no instance-wide /version",
      "method": "GET", "path": "/version", "auth": "none",
      "status": 404, "body": "404 page not found"
    },
    {
      "name": "flow service absent from compose", "why": "#77: the dashboard's fourth target is not routed on upstream 1.2 self-hosted",
      "method": "GET", "path": "/flow/version", "auth": "none",
      "status": 404, "body": "{\"errors\":{\"detail\":\"Page not found\"}}"
    },
    {
      "name": "appengine realm-scoped version, unauthenticated", "why": "#77: realm-scoped AE/RM versions REQUIRE auth (401), unlike pairing's",
      "method": "GET", "path": "/appengine/v1/zzprobe08/version", "auth": "none",
      "status": 401, "body": "{\"errors\":{\"detail\":\"Unauthorized\"}}"
    },
    {
      "name": "appengine realm-scoped version, valid token", "why": "#77 acceptance row: the paired acceptance proving 401 is auth, not absence",
      "method": "GET", "path": "/appengine/v1/zzprobe08/version", "auth": "Bearer <valid a_aea .*::.* token>",
      "status": 200, "body": "{\"data\":\"1.2.0\"}"
    },
    {
      "name": "pairing realm-scoped version, unauthenticated", "why": "#77: upstream serves pairing's realm version PUBLIC — the odd one out of the three services",
      "method": "GET", "path": "/pairing/v1/zzprobe08/version", "auth": "none",
      "status": 200, "body": "{\"data\":\"1.2.0\"}"
    },
    {
      "name": "unknown realm on RM, no token", "why": "#69 baseline row",
      "method": "GET", "path": "/realmmanagement/v1/nosuchrealm/interfaces", "auth": "none",
      "status": 401, "body": "{\"errors\":{\"detail\":\"Unauthorized\"}}"
    },
    {
      "name": "unknown realm on RM, bearer token present", "why": "#69: ANY non-absent token → 403 Forbidden (unknown realm has no key; every verification failure maps to 403). Astrate answers uniform 401 (deliberate deviation-5 extension)",
      "method": "GET", "path": "/realmmanagement/v1/nosuchrealm/interfaces", "auth": "Bearer <token signed by any key>",
      "status": 403, "body": "{\"errors\":{\"detail\":\"Forbidden\"}}"
    },
    {
      "name": "unknown realm on AppEngine", "why": "#69: same class as RM, confirming non-uniformity is per-service wording only",
      "method": "GET", "path": "/appengine/v1/nosuchrealm/devices", "auth": "Bearer <valid bench-realm a_aea token>",
      "status": 403, "body": "{\"errors\":{\"detail\":\"Forbidden\"}}"
    }
  ]
}
---
/root/astrate-mule/test/conformance/upstream/verify-versions.transcript.txt
/root/astrate-mule/test/conformance/upstream/verify-versions.json
/root/astrate-mule/test/conformance/upstream/README.md

Let me read issue #77 for the intended version semantics.
$ gh issue view 77 2>&1 | head -60
title:	Verify per-service version endpoints (GET /version, GET /v1/{realm}/version) served everywhere
state:	CLOSED
author:	atsetilam (Giulio Malitesta)
labels:	enhancement, upstream-parity
comments:	1
assignees:	
projects:	
milestone:	
issue-type:	
parent:	
sub-issues:	
sub-issues-completed:	
blocked-by:	
blocking:	
number:	77
--
Each upstream service exposes unauthenticated `GET /version` and realm-scoped `GET /v1/{realm}/version` returning `{"data": "<version>"}`. The dashboard hits these at login/startup on all four services (RM, AppEngine, pairing, flow). Only RM's realm-scoped endpoint is known-good in Astrate. Audit and fill gaps; keep the APICompatVersion convention consistent.
Size S verification.\n\n---\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).

$ rg -n "0\.1\.0-dev|\.version|/version|VersionHandler|version " /root/astrate-mule/internal/observability /root/astrate-mule/cmd/astrate 2>/dev/null | head -20
/root/astrate-mule/internal/observability/compat.go:36:// VersionHandler answers GET .../version with the upstream envelope
/root/astrate-mule/internal/observability/compat.go:38:func VersionHandler(version string) http.HandlerFunc {
/root/astrate-mule/internal/observability/compat.go:49:// MountVersionCompat registers GET /{service}/version for one service.
/root/astrate-mule/internal/observability/compat.go:50:func MountVersionCompat(mux *http.ServeMux, service, version string) {
/root/astrate-mule/internal/observability/compat.go:51:	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
/root/astrate-mule/cmd/astrate/main.go:51:// version is the reported build version (override with
/root/astrate-mule/cmd/astrate/main.go:52:// -ldflags "-X main.version=vX.Y.Z").
/root/astrate-mule/cmd/astrate/main.go:53:var version = "0.1.0-dev"
/root/astrate-mule/cmd/astrate/main.go:60:	showVersion := flag.Bool("version", false, "print the version and exit")
/root/astrate-mule/cmd/astrate/main.go:387:	// Upstream-parity per-service version endpoints (issue #77): every
/root/astrate-mule/cmd/astrate/main.go:388:	// service answers an unauthenticated GET /{service}/version with
/root/astrate-mule/cmd/astrate/main.go:396:	mux.Handle("GET /appengine/v1/{realm}/version",
/root/astrate-mule/cmd/astrate/main.go:397:		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(version)))
/root/astrate-mule/cmd/astrate/main.go:398:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(version))
/root/astrate-mule/internal/observability/compat_version_test.go:19:		resp, err := http.Get(srv.URL + "/" + svc + "/version")
/root/astrate-mule/internal/observability/compat_version_test.go:21:			t.Fatalf("%s/version: %v", svc, err)
/root/astrate-mule/internal/observability/compat_version_test.go:26:			t.Errorf("%s/version = %d, want 200", svc, resp.StatusCode)
/root/astrate-mule/internal/observability/compat_version_test.go:29:			t.Errorf("%s/version body = %s", svc, body)
/root/astrate-mule/internal/observability/compat_version_test.go:32:			t.Errorf("%s/version missing Content-Type", svc)
/root/astrate-mule/internal/observability/compat_version_test.go:37:func TestVersionHandlerRealmScoped(t *testing.T) {

Let me check how the RM spec documents its realm-scoped version endpoint, to match its convention.
$ rg -n -A 15 "realm.*/version" /root/astrate-mule/docs/api/astarte_realm_management_api.yaml | head -60
532:  /realmmanagement/v1/{realm}/version:
533-    get:
534-      operationId: getVersion
535-      summary: Get emulated API version
536-      description: |
537-        Returns the emulated upstream API compatibility version. The Astarte
538-        Dashboard gates feature UI on this value.
539-      tags: [Realm Management - Misc]
540-      security:
541-        - a_rma: []
542-      parameters:
543-        - $ref: "#/components/parameters/RealmName"
544-      responses:
545-        "200":
546-          description: API compatibility version.
547-          content:

→ Read docs/api/astarte_realm_management_api.yaml [limit=20, offset=545]
The RM spec's sibling endpoint documents it as "API compatibility version". The const `APICompatVersion = "1.2.2"` matches the task's literal `<APICompatVersion>`. Let me align all four examples to that real value and mirror the RM spec's wording.

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -184,34 +184,36 @@
 
 /appengine/version:
   get:
     operationId: getAppEngineVersion
-    summary: AppEngine compat version
+    summary: AppEngine API version
     description: |
-      Upstream-parity unauthenticated version endpoint. Answers the
-      AppEngine service's API compatibility version.
+      Upstream-parity unauthenticated version endpoint. Returns the API
+      compatibility version of the AppEngine service; the Astarte Dashboard
+      reads one per service at startup.
     tags: [Compat Health]
     responses:
       "200":
-        description: Service's API compatibility version.
+        description: API compatibility version.
         content:
           application/json:
             schema:
               $ref: "#/components/schemas/DataEnvelopeVersion"
             example:
-              data: "1.0.0"
+              data: "1.2.2"
 
 /realmmanagement/version:
   get:
     operationId: getRealmManagementVersion
-    summary: Realm Management compat version
+    summary: Realm Management API version
     description: |
-      Upstream-parity unauthenticated version endpoint. Answers the
-      Realm Management service's API compatibility version.
+      Upstream-parity unauthenticated version endpoint. Returns the API
+      compatibility version of the Realm Management service; the Astarte
+      Dashboard reads one per service at startup.
     tags: [Compat Health]
     responses:
       "200":
-        description: Service's API compatibility version.
+        description: API compatibility version.
         content:
           application/json:
             schema:
               $ref: "#/components/schemas/DataEnvelopeVersion"
@@ -220,40 +222,42 @@
 
 /pairing/version:
   get:
     operationId: getPairingVersion
-    summary: Pairing compat version
+    summary: Pairing API version
     description: |
-      Upstream-parity unauthenticated version endpoint. Answers the
-      Pairing service's API compatibility version.
+      Upstream-parity unauthenticated version endpoint. Returns the API
+      compatibility version of the Pairing service; the Astarte Dashboard
+      reads one per service at startup.
     tags: [Compat Health]
     responses:
       "200":
-        description: Service's API compatibility version.
+        description: API compatibility version.
         content:
           application/json:
             schema:
               $ref: "#/components/schemas/DataEnvelopeVersion"
             example:
-              data: "1.0.0"
+              data: "1.2.2"
 
 /housekeeping/version:
   get:
     operationId: getHousekeepingVersion
-    summary: Housekeeping compat version
+    summary: Housekeeping API version
     description: |
-      Upstream-parity unauthenticated version endpoint. Answers the
-      Housekeeping service's API compatibility version.
+      Upstream-parity unauthenticated version endpoint. Returns the API
+      compatibility version of the Housekeeping service; the Astarte
+      Dashboard reads one per service at startup.
     tags: [Compat Health]
     responses:
       "200":
-        description: Service's API compatibility version.
+        description: API compatibility version.
         content:
           application/json:
             schema:
               $ref: "#/components/schemas/DataEnvelopeVersion"
             example:
-              data: "1.0.0"
+              data: "1.2.2"
 
 # ── Native WebSocket (live stream) ──────────────────────────────────
 
 /astrate/v1/{realm}/socket:

✗ Edit docs/api/astrate_native_api.yaml failed
Error: Could not find oldString in the file. It must match exactly, including whitespace, indentation, and line endings.
→ Read docs/api/astrate_native_api.yaml [limit=120, offset=155]
The YAML looks right. Now let me verify the schema addition and run the build plus a YAML-parse check.
→ Read docs/api/astrate_native_api.yaml [limit=30, offset=400]
Schema is in place. Now the build and a YAML/OpenAPI parse validation.

$ PATH="/tmp/mule-docs-venv/bin:$PATH" make -C docs build 2>&1 | grep -Ev "^\[3[0-9]m|squidfunk|^\[0m|^ \[0m" | tail -8
INFO    -  Cleaning site directory
INFO    -  Building documentation to directory: /root/astrate-mule/site-dist
INFO    -  The following pages exist in the docs directory, but are not included in the "nav" configuration:
  - DESIGN.md
  - ROADMAP.md
INFO    -  Doc file 'index.md' contains an unrecognized relative link 'api/', it was left as is.
INFO    -  Documentation built in 7.09 seconds
make: Leaving directory '/root/astrate-mule/docs'

Build succeeded. Let me validate the YAML is structurally valid OpenAPI.
$ /tmp/mule-docs-venv/bin/pip install -q openapi-spec-validator 2>&1 | tail -2; /tmp/mule-docs-venv/bin/python3 -c "
import yaml, json
from openapi_spec_validator import validate
spec = yaml.safe_load(open('docs/api/astrate_native_api.yaml'))
validate(spec)
paths = spec['paths']
for p in ['/appengine/version','/realmmanagement/version','/pairing/version','/housekeeping/version']:
    assert p in paths, p
    op = paths[p]['get']
    assert '200' in op['responses'], p
    assert 'schemas' in spec['components']
for svc in ['appengine','realmmanagement','pairing','housekeeping']:
    assert 'DataEnvelopeVersion' in spec['components']['schemas']
    print('OK', '/%s/version' % svc, '-> 200', json.dumps(op['responses']['200']['content'].get('application/json', {}).get('example', {})))
print('VALID OPENAPI 3.0')
"
OK /appengine/version -> 200 {"data": "1.2.2"}
OK /realmmanagement/version -> 200 {"data": "1.2.2"}
OK /pairing/version -> 200 {"data": "1.2.2"}
OK /housekeeping/version -> 200 {"data": "1.2.2"}
VALID OPENAPI 3.0

Spec is valid. Now the Go gate (docs-only change, but MULE requires it) plus a check that no other doc inventory lists these endpoints.
$ rg -n "appengine/version|housekeeping/version|pairing/version|realmmanagement/version" docs --glob '!api/astrate_native_api.yaml' 2>/dev/null
docs/api/astrate_native_api.yaml:185:  /appengine/version:
docs/api/astrate_native_api.yaml:204:  /realmmanagement/version:
docs/api/astrate_native_api.yaml:223:  /pairing/version:
docs/api/astrate_native_api.yaml:242:  /housekeeping/version:
docs/site/api/astrate_native_api.yaml:185:  /appengine/version:
docs/site/api/astrate_native_api.yaml:204:  /realmmanagement/version:
docs/site/api/astrate_native_api.yaml:223:  /pairing/version:
docs/site/api/astrate_native_api.yaml:242:  /housekeeping/version:

$ go vet ./... 2>&1 | tail -5 && echo "VET OK"
VET OK

$ git status --short 2>&1 | head; echo "---"; git check-ignore docs/site 2>&1; echo "---"; git config --get core.excludesfile
M docs/api/astrate_native_api.yaml
---
---

Only the intended file is modified; `docs/site/` is a gitignored build artifact. Running the test gate.
$ go test ./... 2>&1 | grep -Ev "^(ok|---)" | head -20; echo "EXIT: ${PIPESTATUS[0]}"
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]
/bin/sh: 1: Bad substitution

$ go test ./... > /tmp/gotest.log 2>&1; printf 'exit=%s\n' "$?"; grep -cE "^ok" /tmp/gotest.log; grep -E "^(FAIL|--- FAIL)" /tmp/gotest.log | head; tail -3 /tmp/gotest.log
exit=0
30
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

$ gofmt -l . 2>&1; echo "gofmt-exit=$?"
gofmt-exit=0

gofmt clean. Let me refresh the top-of-file description, which enumerates covered endpoints, so it reflects the new version endpoints too.
← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -2,11 +2,11 @@
 info:
   title: Astrate Native API
   description: |
     Astrate-native operational and real-time API surface. Covers liveness,
-    readiness, and Prometheus metrics (no auth), compat health endpoints for
-    the Astarte Dashboard, and two WebSocket endpoints: the native live
-    event stream and the Phoenix V2-compatible channels endpoint.
+    readiness, and Prometheus metrics (no auth), compat health and version
+    endpoints for the Astarte Dashboard, and two WebSocket endpoints: the
+    native live event stream and the Phoenix V2-compatible channels endpoint.
   version: "0.0.1"
   license:
     name: "Apache-2.0 with SECO Exclusion"
     url: https://www.apache.org/licenses/LICENSE-2.0

$ PATH="/tmp/mule-docs-venv/bin:$PATH" make -C docs build >/dev/null 2>&1 && echo "BUILD OK"; /tmp/mule-docs-venv/bin/python3 -c "
from openapi_spec_validator import validate
import yaml
validate(yaml.safe_load(open('docs/api/astrate_native_api.yaml')))
print('VALID OPENAPI 3.0')
"
BUILD OK
VALID OPENAPI 3.0

$ git status --short
M docs/api/astrate_native_api.yaml

# Todos
[✓] Read docs/api/astrate_native_api.yaml to understand existing health endpoint pattern
[✓] Read observability code to confirm response shape
[✓] Add four /{service}/version endpoints to the YAML
[✓] Run make -C docs build to verify YAML loads
[✓] Report

Done: documented the four unauthenticated `GET /{service}/version` endpoints (`appengine`, `realmmanagement`, `pairing`, `housekeeping`) in astrate_native_api.yaml — 200 `{"data":"<version>"}` — with a new `DataEnvelopeVersion` schema; document-level `security: []` keeps them unauthenticated.
Files: docs/api/astrate_native_api.yaml (also touched only the build artifact docs/site/, which is gitignored)
Verified: `make -C docs build` (mkdocs seeded into /tmp venv) -> pass; openapi-spec-validator on the Swagger-referenced YAML -> valid; `go vet ./...` -> pass; `go test ./...` (no -race) -> pass, 30 pkgs; `gofmt -l .` -> clean
Unsure: task wording says the body is `<APICompatVersion>`; code actually serves the build `version` var (`0.1.0-dev` default, main.go:53). I documented the envelope per task and used the repo's real `APICompatVersion` "1.2.2" as the example value, mirroring the RM spec's convention. No behaviour changed, so no proving test applies (docs-only diff).
Follow-ups: none
