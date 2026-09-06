slug: docs-sync-rm-datastream-retention-endpoint
verdict: done
at:  0574fd9
ran: 2026-09-06T11:31:43Z on DietPi in 232s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
Task: Add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` endpoint to the OpenAPI spec. Let me read the relevant code and the YAML file.
→ Read internal/realm/http.go [offset=40, limit=110]
→ Read docs/api/astarte_realm_management_api.yaml [offset=415, limit=80]
Now let me check the return type of `GetDatastreamMaximumStorageRetention` to get the schema right.
✱ Grep "GetDatastreamMaximumStorageRetention" in . · 4 matches
→ Read internal/realm/service.go [limit=20, offset=595]
✱ Grep "datastream_maximum_storage_retention|int64|seconds" in docs/api · 5 matches
The return type is `int64` (seconds), with an unset ceiling rendering as `0`. Let me check if there's an upstream spec reference for this endpoint to match wire format.
✱ Grep "datastream_maximum_storage_retention|retention" in docs · 50 matches
docs/site/api has a matching file. Let me check if it's a symlink or a separate copy.

$ ls -la docs/site/api/astarte_realm_management_api.yaml docs/api/astarte_realm_management_api.yaml; ls -la docs/site/api/ | head -20
-rw-r--r-- 1 root root 33035 Sep  4 18:55 docs/api/astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 33035 Sep  5 19:21 docs/site/api/astarte_realm_management_api.yaml
total 116
drwxr-xr-x 2 root root  4096 Sep  5 19:21 .
drwxr-xr-x 4 root root  4096 Sep  5 19:21 ..
-rw-r--r-- 1 root root 26566 Sep  5 19:21 astarte_appengine_api.yaml
-rw-r--r-- 1 root root  9444 Sep  5 19:21 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 18942 Sep  5 19:21 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 33035 Sep  5 19:21 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 10378 Sep  5 19:21 astrate_native_api.yaml

→ Read docs/Makefile
`make -C docs build` syncs `docs/api` → `docs/site/api`, so I only edit `docs/api`. Now adding the endpoint between `device_registration_limit` and `version`.

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -455,8 +455,41 @@
         $ref: "#/components/responses/Unauthorized"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
+/realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention:
+  get:
+    operationId: getDatastreamMaximumStorageRetention
+    summary: Get datastream maximum storage retention
+    description: |
+      Returns the realm's datastream maximum storage retention in seconds
+      (upstream serves this since 1.2.0). An unset ceiling renders as 0 —
+      a recorded deviation from upstream, which answers null here (issue #60).
+    tags: [Realm Management - Config]
+    security:
+      - a_rma: []
+    parameters:
+      - $ref: "#/components/parameters/RealmName"
+    responses:
+      "200":
+        description: Datastream maximum storage retention in seconds.
+        content:
+          application/json; charset=utf-8:
+            schema:
+              type: object
+              required: [data]
+              properties:
+                data:
+                  type: integer
+                  format: int64
+                  description: Maximum datastream storage retention in seconds. 0 means unset.
+            example:
+              data: 31536000
+      "401":
+        $ref: "#/components/responses/Unauthorized"
+      "500":
+        $ref: "#/components/responses/InternalServerError"
+
 # ── Version ─────────────────────────────────────────────────────────
 
 /realmmanagement/v1/{realm}/version:
   get:

Now run the docs build gate and validate the YAML loads.
$ make -C docs build 2>&1 | tail -30
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

mkdocs is missing. Let me check requirements and try installing it, and validate the YAML independently.
$ cat docs/requirements.txt; python3 -c "import yaml,sys; yaml.safe_load(open('docs/api/astarte_realm_management_api.yaml')); print('YAML OK')"
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml,sys; yaml.safe_load(open('docs/api/astarte_realm_management_api.yaml')); print('YAML OK')
    ^^^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ pip install -q -r docs/requirements.txt 2>&1 | tail -5; which mkdocs
/bin/sh: 1: pip: not found

$ python3 -m pip --version 2>&1; python3 -m venv --help >/dev/null 2>&1 && echo venv-ok; uv --version 2>/dev/null || true
/usr/bin/python3: No module named pip
venv-ok

$ python3 -m ensurepip --user 2>&1 | tail -3
module, so use it on your own risk, or make sure to only use it in virtual
environments.

$ python3 -m ensurepip --user >/dev/null 2>&1; python3 -m pip --user --version 2>&1 | head -2
/usr/bin/python3: No module named pip

$ python3 -m ensurepip --user 2>&1 | head -20
ensurepip is disabled in Debian/Ubuntu for the system python.

Python modules for the system python are usually handled by dpkg and apt-get.

    apt install python3-<module name>

Install the python3-pip package to use pip itself.  Using pip together
with the system python might have unexpected results for any system installed
module, so use it on your own risk, or make sure to only use it in virtual
environments.

$ python3 -m venv /tmp/opencode/mkdocs-venv && /tmp/opencode/mkdocs-venv/bin/pip --version
pip 25.1.1 from /tmp/opencode/mkdocs-venv/lib/python3.13/site-packages/pip (python 3.13)

$ /tmp/opencode/mkdocs-venv/bin/pip install -q -r docs/requirements.txt 2>&1 | tail -5; /tmp/opencode/mkdocs-venv/bin/mkdocs --version
mkdocs, version 1.6.1 from /tmp/opencode/mkdocs-venv/lib/python3.13/site-packages/mkdocs (Python 3.13)

mkdocs installed in the venv. Now run the build with the venv on PATH.
$ PATH=/tmp/opencode/mkdocs-venv/bin:$PATH make -C docs build 2>&1 | tail -25
mkdocs build -f mkdocs.yml

 │  ⚠  Warning from the Material for MkDocs team
 │
 │  MkDocs 2.0, the underlying framework of Material for MkDocs,
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
INFO    -  Documentation built in 3.93 seconds
make: Leaving directory '/root/astrate-mule/docs'

$ grep -n "swagger\|api/\|yaml" docs/mkdocs.yml | head -20
63:  - API Explorer: swagger.md

Build succeeded. Let me confirm how the YAML is Swagger-referenced and verify my endpoint parses with a proper YAML/OpenAPI load.
$ /tmp/opencode/mkdocs-venv/bin/python -c "
import yaml
d = yaml.safe_load(open('docs/api/astarte_realm_management_api.yaml'))
op = d['paths']['/realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention']['get']
print('operationId:', op['operationId'])
print('responses:', sorted(op['responses']))
print('data type:', op['responses']['200']['content']['application/json; charset=utf-8']['schema']['properties']['data'])
print('total paths:', len(d['paths']))
"
operationId: getDatastreamMaximumStorageRetention
responses: ['200', '401', '500']
data type: {'type': 'integer', 'format': 'int64', 'description': 'Maximum datastream storage retention in seconds. 0 means unset.'}
total paths: 12

$ cat docs/swagger.md 2>/dev/null; ls docs/swagger-ui 2>/dev/null | head
index.html

$ grep -rn "swagger\|realm_management_api" docs/swagger.md docs/*.md 2>/dev/null | grep -i realm | head; ls docs/*.md
docs/COMPATIBILITY.md
docs/DESIGN.md
docs/JSON-PAYLOAD-PROFILE.md
docs/OPERATIONS.md
docs/ROADMAP.md
docs/UPSTREAM-EXPERIMENTAL.md

$ grep -n "astarte_realm\|urls\|yaml" docs/swagger-ui/index.html | head
44:      <option value="../api/astarte_housekeeping_api.yaml">Housekeeping API — Realm management, a_ha auth</option>
45:      <option value="../api/astarte_pairing_api.yaml">Pairing API — Device pairing flows A/B/C, a_pa auth</option>
46:      <option value="../api/astarte_realm_management_api.yaml">Realm Management API — Interfaces, triggers, policies, a_rma auth</option>
47:      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
48:      <option value="../api/astrate_native_api.yaml">Astrate Native API — Health, readiness, metrics, WebSocket</option>

$ grep -c "datastream_maximum_storage_retention" site-dist/api/astarte_realm_management_api.yaml docs/site/api/astarte_realm_management_api.yaml docs/api/astarte_realm_management_api.yaml; ls site-dist/swagger-ui/index.html
site-dist/api/astarte_realm_management_api.yaml:1
docs/site/api/astarte_realm_management_api.yaml:1
docs/api/astarte_realm_management_api.yaml:1
site-dist/swagger-ui/index.html

All three copies contain the new path and the Swagger UI reference resolves. Let me check working-tree state for the report.
$ git status --short 2>/dev/null; echo ---; git check-ignore site-dist docs/site 2>/dev/null || true
M docs/api/astarte_realm_management_api.yaml
---
site-dist

$ git ls-files docs/site | head -3; echo "---"; git ls-files docs/site | wc -l; echo "---"; git status --short docs/api docs/site | head
docs/site/appengine-api.md
docs/site/architecture.md
docs/site/compatibility.md
---
23
---
 M docs/api/astarte_realm_management_api.yaml

$ git check-ignore -v docs/site/api/astarte_realm_management_api.yaml docs/site/swagger-ui/index.html 2>/dev/null; git ls-files docs/site | grep -E "api|swagger" | head
.gitignore:34:/docs/site/api/	docs/site/api/astarte_realm_management_api.yaml
.gitignore:35:/docs/site/swagger-ui/	docs/site/swagger-ui/index.html
docs/site/appengine-api.md
docs/site/housekeeping-api.md
docs/site/realm-management-api.md
docs/site/swagger.md

Task complete. The only tracked change is `docs/api/astarte_realm_management_api.yaml`; the build regenerated only gitignored artifacts (`docs/site/api/`, `docs/site/swagger-ui/`, `site-dist/`).

```
Done: Documented GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention (getDatastreamMaximumStorageRetention, int64 data-envelope in seconds, 0-when-unset deviation) in docs/api/astarte_realm_management_api.yaml as a new Realm Management - Config path between device_registration_limit and version.
Files: docs/api/astarte_realm_management_api.yaml
Verified: make -C docs build -> pass (mkdocs 1.6.1 via throwaway venv at /tmp/opencode/mkdocs-venv since system python has no pip/mkdocs; full site build in 3.93s, only pre-existing warnings); yaml.safe_load of docs/api/astarte_realm_management_api.yaml -> pass, new path present with operationId getDatastreamMaximumStorageRetention, 200/401/500 responses, data type integer/int64; swagger-ui index.html reference ../api/astarte_realm_management_api.yaml resolves to the synced copy in site-dist which contains the new path. No Go code touched, so go vet/test not applicable.
Unsure: nothing (response schema modeled on the existing device_registration_limit block; unset-ceiling-as-0 described as documented deviation per service.go:603/#60, not guessed)
Follow-ups: none
```
