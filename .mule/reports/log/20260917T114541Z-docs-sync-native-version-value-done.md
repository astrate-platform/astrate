slug: docs-sync-native-version-value
verdict: done
at:  08048b3
ran: 2026-09-17T11:45:41Z on DietPi in 503s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astrate_native_api.yaml [offset=175, limit=95]
→ Read cmd/astrate/main.go [offset=40, limit=30]
→ Read cmd/astrate/main.go [limit=40, offset=375]
→ Read docs/api/astrate_native_api.yaml [limit=35, offset=405]
→ Read internal/realm/service.go [limit=30, offset=575]
✱ Grep "func MountVersionCompat" in internal · 1 match
✱ Grep "mountRealmVersion|APICompatVersion" in cmd · 7 matches
✗ Read docs/compat.md failed [limit=30, offset=60]
Error: File not found: /root/astrate-mule/docs/compat.md
→ Read test/conformance/upstream/verify-versions.json
✱ Glob "**/compat*.md" 1 match
→ Read internal/observability/compat.go
→ Read docs/site/compatibility.md [limit=30, offset=55]
✱ Grep "appengine/version|housekeeping/version|/pairing/version|MountVersionCompat|VersionHandler" in . · 12 matches
✱ Grep "version" in docs · 100 matches
✱ Grep "DataEnvelopeVersion" in docs/api/astrate_native_api.yaml · 5 matches
→ Read docs/Makefile
✱ Grep "astrate_native_api|embed|site/api" in internal/swagger · 9 matches
→ Read internal/swagger/swagger.go
→ Read internal/swagger/swagger_test.go
✱ Grep "go:embed" in docs · 17 matches
✱ Grep "astrate_native_api|native_api" in . · 0 matches
This is a doc-sync task: the four service-root version endpoints in the spec claim `data: "1.2.2"` but the wire serves the build `version` var. Let me confirm the two candidate facts before deciding.

Key evidence gathered:
- `main.go:393-395` serves the build `version` (default `"0.1.0-dev"`) on the four service-root endpoints.
- `mountRealmVersion` (main.go:432-435) serves `realm.APICompatVersion` (1.2.2) only on the realm-scoped endpoints.
- `docs/site/compatibility.md:75` (deviation 10) is scoped to `GET /v1/<realm>/version`.
- Upstream 1.2.0 serves the install version (1.2.0) on both classes; the astrate analogue of "install version" is the build version.

Decision: **document the build version** (doc-only, spec matches the wire). The service-root endpoints are what the dashboard displays to identify the running software, so an honest build version there is correct; the compat level already lives on the realm-scoped endpoints where feature-gating reads it. Now editing the YAML.

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -186,21 +186,22 @@
   get:
     operationId: getAppEngineVersion
     summary: AppEngine API version
     description: |
-      Upstream-parity unauthenticated version endpoint. Returns the API
-      compatibility version of the AppEngine service; the Astarte Dashboard
-      reads one per service at startup.
+      Upstream-parity unauthenticated version endpoint. Returns the Astrate
+      build version serving the AppEngine service; the Astarte Dashboard
+      reads one per service at startup. The emulated upstream API level is
+      served by the realm-scoped /v1/{realm}/version endpoints.
     tags: [Compat Health]
     responses:
       "200":
-        description: API compatibility version.
+        description: Astrate build version.
         content:
           application/json:
             schema:
               $ref: "#/components/schemas/DataEnvelopeVersion"
             example:
-              data: "1.2.2"
+              data: "0.1.0-dev"
 
 /realmmanagement/version:
   get:
     operationId: getRealmManagementVersion

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -206,21 +206,22 @@
   get:
     operationId: getRealmManagementVersion
     summary: Realm Management API version
     description: |
-      Upstream-parity unauthenticated version endpoint. Returns the API
-      compatibility version of the Realm Management service; the Astarte
-      Dashboard reads one per service at startup.
+      Upstream-parity unauthenticated version endpoint. Returns the Astrate
+      build version serving the Realm Management service; the Astarte
+      Dashboard reads one per service at startup. The emulated upstream API
+      level is served by the realm-scoped /v1/{realm}/version endpoints.
     tags: [Compat Health]
     responses:
       "200":
-        description: API compatibility version.
+        description: Astrate build version.
         content:
           application/json:
             schema:
               $ref: "#/components/schemas/DataEnvelopeVersion"
             example:
-              data: "1.2.2"
+              data: "0.1.0-dev"
 
 /pairing/version:
   get:
     operationId: getPairingVersion

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -226,21 +226,22 @@
   get:
     operationId: getPairingVersion
     summary: Pairing API version
     description: |
-      Upstream-parity unauthenticated version endpoint. Returns the API
-      compatibility version of the Pairing service; the Astarte Dashboard
-      reads one per service at startup.
+      Upstream-parity unauthenticated version endpoint. Returns the Astrate
+      build version serving the Pairing service; the Astarte Dashboard reads
+      one per service at startup. The emulated upstream API level is served
+      by the realm-scoped /v1/{realm}/version endpoints.
     tags: [Compat Health]
     responses:
       "200":
-        description: API compatibility version.
+        description: Astrate build version.
         content:
           application/json:
             schema:
               $ref: "#/components/schemas/DataEnvelopeVersion"
             example:
-              data: "1.2.2"
+              data: "0.1.0-dev"
 
 /housekeeping/version:
   get:
     operationId: getHousekeepingVersion

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -246,21 +246,22 @@
   get:
     operationId: getHousekeepingVersion
     summary: Housekeeping API version
     description: |
-      Upstream-parity unauthenticated version endpoint. Returns the API
-      compatibility version of the Housekeeping service; the Astarte
-      Dashboard reads one per service at startup.
+      Upstream-parity unauthenticated version endpoint. Returns the Astrate
+      build version serving the Housekeeping service; the Astarte Dashboard
+      reads one per service at startup. The emulated upstream API level is
+      served by the realm-scoped /v1/{realm}/version endpoints.
     tags: [Compat Health]
     responses:
       "200":
-        description: API compatibility version.
+        description: Astrate build version.
         content:
           application/json:
             schema:
               $ref: "#/components/schemas/DataEnvelopeVersion"
             example:
-              data: "1.2.2"
+              data: "0.1.0-dev"
 
 # ── Native WebSocket (live stream) ──────────────────────────────────
 
 /astrate/v1/{realm}/socket:

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -416,9 +416,9 @@
   required: [data]
   properties:
     data:
       type: string
-      description: The service's API compatibility version.
+      description: The Astrate build version.
 
 ErrorDetail:
   type: object
   required: [errors]

→ Read docs/api/astrate_native_api.yaml [limit=80, offset=183]
Now verifying. Let me re-read the changed sections and check tooling for validation.

$ python3 -c "import yaml; d=yaml.safe_load(open('/root/astrate-mule/docs/api/astrate_native_api.yaml')); print('yaml ok, paths:', list(d['paths'].keys()))" 2>&1; which mkdocs; python3 -c "import mkdocs; print(mkdocs.__version__)" 2>&1
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; d=yaml.safe_load(open('/root/astrate-mule/docs/api/astrate_native_api.yaml')); print('yaml ok, paths:', list(d['paths'].keys()))
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import mkdocs; print(mkdocs.__version__)
    ^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'mkdocs'

PyYAML/mkdocs not installed. Let me validate the YAML with Go (the repo likely vendors a yaml parser) and try installing mkdocs for the site build.
✱ Grep "yaml" in go.mod · 100 matches
→ Read docs/requirements.txt
→ Read .mule/recipes/docs-sync.md
This is a docs-sync task. Let me read the docs-sync recipe to follow its exact gate, and confirm the realm-scoped wiring.
→ Read cmd/astrate/main.go [limit=14, offset=426]
Code comment at main.go:426-431 confirms the intended dichotomy: realm-scoped = emulated API level, service-root = build version. Decision is well-supported. Now run the docs build gate — mkdocs isn't installed, so let me install it.

$ python3 --version; python3 -m pip --version 2>&1; pip3 --version 2>&1
Python 3.13.5
/usr/bin/python3: No module named pip
/bin/sh: 1: pip3: not found

$ python3 -m ensurepip --version 2>&1; python3 -m venv /tmp/mkvenv 2>&1 && /tmp/mkvenv/bin/pip --version 2>&1
pip 25.1.1
pip 25.1.1 from /tmp/mkvenv/lib/python3.13/site-packages/pip (python 3.13)

$ git status --short; echo "---"; git status | head -20
M docs/api/astrate_native_api.yaml
---
On branch mule/queue
Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
	modified:   docs/api/astrate_native_api.yaml

no changes added to commit (use "git add" and/or "git commit -a")

Tree was clean before my edit. Now installing mkdocs for the build gate.
$ /tmp/mkvenv/bin/pip install -q "mkdocs>=1.6,<2" "mkdocs-material>=9.5,<10" 2>&1 | tail -5; /tmp/mkvenv/bin/mkdocs --version 2>&1
mkdocs, version 1.6.1 from /tmp/mkvenv/lib/python3.13/site-packages/mkdocs (Python 3.13)

$ PATH=/tmp/mkvenv/bin:$PATH make -C docs build 2>&1 | tail -20; echo "EXIT=$?"
 │  will introduce backward-incompatible changes, including:
 │
 │  × All plugins will stop working – the plugin system has been removed
 │  × All theme overrides will break – the theming system has been rewritten
 │  × No migration path exists – existing projects cannot be upgraded
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
INFO    -  Documentation built in 9.05 seconds
make: Leaving directory '/root/astrate-mule/docs'
EXIT=0

$ git status --short
M docs/api/astrate_native_api.yaml

$ git diff --stat; echo "---"; git status --short --ignored | head; echo "---"; rg -n 'data: "0.1.0-dev"|data: "1.2.2"|Astrate build version|API compatibility version' docs/site/api/astrate_native_api.yaml | head
docs/api/astrate_native_api.yaml | 46 ++++++++++++++++++++++------------------
 1 file changed, 25 insertions(+), 21 deletions(-)
---
 M docs/api/astrate_native_api.yaml
!! .mule/.budget
!! .mule/.heartbeat
!! .mule/.last-output
!! .mule/.lock
!! .mule/.rotation
!! .mule/cron.log
!! .mule/failed/
!! .mule/hosts
!! .mule/task.md
---
197:          description: Astrate build version.
203:                data: "0.1.0-dev"
217:          description: Astrate build version.
223:                data: "0.1.0-dev"
237:          description: Astrate build version.
243:                data: "0.1.0-dev"
257:          description: Astrate build version.
263:                data: "0.1.0-dev"
420:          description: The Astrate build version.

$ git ls-files docs/site/api | head; echo "---"; git check-ignore -v docs/site/api/astrate_native_api.yaml; echo "---"; rg -n 'site/api|site-dist' .gitignore 2>/dev/null; git show HEAD:docs/api/astrate_native_api.yaml >/dev/null 2>&1 && echo "api tracked"
---
.gitignore:34:/docs/site/api/	docs/site/api/astrate_native_api.yaml
---
29:/site-dist/
34:/docs/site/api/
api tracked

Build passed and the tree is clean (site-dist and site/api are git-ignored build artefacts). Now validating the YAML parses and running the Go gate.
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ mkdir -p /tmp/yamlprobe && printf 'module yamlprobe\n\ngo 1.22\n\nrequire gopkg.in/yaml.v3 v3.0.1\n' > /tmp/yamlprobe/go.mod && printf 'package main\n\nimport (\n\t"fmt"\n\t"os"\n\n\t"gopkg.in/yaml.v3"\n)\n\nfunc main() {\n\tf, err := os.ReadFile(os.Args[1])\n\tif err != nil { panic(err) }\n\tvar v any\n\tif err := yaml.Unmarshal(f, &v); err != nil {\n\t\tfmt.Println("PARSE ERROR:", err)\n\t\tos.Exit(1)\n\t}\n\tfmt.Println("yaml parse OK")\n}\n' > /tmp/yamlprobe/main.go && cd /tmp/yamlprobe && GOPROXY=off GOFLAGS=-mod=mod go run . /root/astrate-mule/docs/api/astrate_native_api.yaml 2>&1 failed
Error: The user rejected permission to use this specific tool call.
