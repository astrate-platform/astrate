slug: docs-sync-pairing-health-path
verdict: done
at:  d73e225
ran: 2026-09-05T17:17:55Z on DietPi in 186s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_pairing_api.yaml
I'll read the existing YAML, the code reference, and check for a Makefile.
→ Read internal/pairing/http.go [limit=40, offset=70]
→ Read docs/Makefile
Let me read the health handler to confirm response shapes.
✱ Grep "handleHealth" in internal/pairing · 3 matches
→ Read internal/pairing/http.go [offset=278, limit=50]
Now I have everything. Let me add the health endpoint and then verify the YAML loads.

← Edit docs/api/astarte_pairing_api.yaml
Index: /root/astrate-mule/docs/api/astarte_pairing_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_pairing_api.yaml
+++ /root/astrate-mule/docs/api/astarte_pairing_api.yaml
@@ -23,8 +23,51 @@
 
 security: []
 
 paths:
+  /pairing/v1/{realm}/health:
+    get:
+      operationId: getHealth
+      summary: Realm health check
+      description: |
+        Returns the health status of a realm's backing store. Unauthenticated
+        and unrate-limited — FDO-flow devices probe it before they hold
+        credentials.
+
+        Returns 200 when the realm exists and the database answers, 404 for an
+        unknown realm, and 503 when the database is unhealthy.
+      tags: [Pairing]
+      parameters:
+        - $ref: "#/components/parameters/RealmName"
+      responses:
+        "200":
+          description: Realm is healthy.
+          content:
+            application/json; charset=utf-8:
+              schema:
+                $ref: "#/components/schemas/DataEnvelopeHealthStatus"
+              example:
+                data:
+                  status: ok
+        "404":
+          description: Unknown realm.
+          content:
+            application/json; charset=utf-8:
+              schema:
+                $ref: "#/components/schemas/ErrorDetail"
+              example:
+                errors:
+                  detail: Not Found
+        "503":
+          description: Realm exists but the backing store is unhealthy.
+          content:
+            application/json; charset=utf-8:
+              schema:
+                $ref: "#/components/schemas/ErrorDetail"
+              example:
+                errors:
+                  detail: unhealthy
+
   /pairing/v1/{realm}/agent/devices:
     post:
       operationId: registerDevice
       summary: Register a new device (flow A)

← Edit docs/api/astarte_pairing_api.yaml
Index: /root/astrate-mule/docs/api/astarte_pairing_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_pairing_api.yaml
+++ /root/astrate-mule/docs/api/astarte_pairing_api.yaml
@@ -468,8 +468,24 @@
   properties:
     data:
       $ref: "#/components/schemas/PairingInfo"
 
+HealthStatus:
+  type: object
+  required: [status]
+  properties:
+    status:
+      type: string
+      enum: [ok]
+      description: Realm health status.
+
+DataEnvelopeHealthStatus:
+  type: object
+  required: [data]
+  properties:
+    data:
+      $ref: "#/components/schemas/HealthStatus"
+
 DataEnvelopeVerifyValid:
   type: object
   required: [data]
   properties:

$ python3 -c "import yaml; yaml.safe_load(open('docs/api/astarte_pairing_api.yaml')); print('YAML OK')" 2>&1 || pip3 list 2>/dev/null | grep -i yaml
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; yaml.safe_load(open('docs/api/astarte_pairing_api.yaml')); print('YAML OK')
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

Now let me verify the YAML loads and run the docs build.
→ Read docs/mkdocs.yml
→ Read docs/requirements.txt
No pyyaml installed. Let me check requirements and try installing the doc deps to run the build.

$ python3 -m pip --version 2>&1; python3 -c "import mkdocs" 2>&1; which mkdocs 2>&1
/usr/bin/python3: No module named pip
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import mkdocs
ModuleNotFoundError: No module named 'mkdocs'

$ grep -i yaml go.mod go.sum 2>/dev/null | head; which go ruby node 2>&1; ls ~/.local/go/bin 2>/dev/null
go.mod:	go.yaml.in/yaml/v2 v2.4.2 // indirect
go.mod:	gopkg.in/yaml.v3 v3.0.1 // indirect
go.sum:go.yaml.in/yaml/v2 v2.4.2 h1:DzmwEr2rDGHl7lsFgAHxmNz/1NlQ7xLIrlN2h5d1eGI=
go.sum:go.yaml.in/yaml/v2 v2.4.2/go.mod h1:081UH+NErpNdqlCXm3TtEran0rJZGxAYx9hb/ELlsPU=
go.sum:gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
go.sum:gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
go.sum:gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
/usr/bin/go
/usr/local/bin/node

Go is available. Let me check how the swagger page references the YAML, then write a probe.
$ grep -rn "astarte_pairing_api" docs/site/swagger.md docs/swagger.md 2>/dev/null | head; ls docs/site 2>/dev/null | head
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

$ grep -rn "api/" docs/site/swagger.md | head -20
(no output)

$ ls docs/site/*.yaml 2>/dev/null; rg -n "astarte_pairing_api" docs/site -l 2>/dev/null | head; rg -n "pairing" docs/site/swagger.md 2>/dev/null | head
(no output)

$ ls docs/api/; rg -ln "astarte_pairing_api" docs/ 2>/dev/null | grep -v site/
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
docs/swagger-ui/index.html
docs/handoff/phase-1-memory.md

$ rg -n "astarte_pairing_api" docs/swagger-ui/index.html | head
45:      <option value="../api/astarte_pairing_api.yaml">Pairing API — Device pairing flows A/B/C, a_pa auth</option>

The build needs mkdocs/pip which aren't installed; I'll replicate the `sync` step (copies `api` → `site/api`, the path swagger-ui actually loads) and validate with a Go probe using `gopkg.in/yaml.v3` from the repo's module cache.
$ go env GOMODCACHE GOFLAGS GOPROXY 2>&1; ls $(go env GOMODCACHE)/gopkg.in/ 2>/dev/null | head
/root/go/pkg/mod

https://proxy.golang.org,direct
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1

$ rm -rf docs/site/api && cp -r docs/api docs/site/api && echo "sync OK"
sync OK

← Write /tmp/opencode/yamlcheck/main.go
Wrote file successfully.

$ cd /tmp/opencode/yamlcheck && printf 'module yamlcheck\n\ngo 1.24\n\nrequire gopkg.in/yaml.v3 v3.0.1\n' > go.mod && GOPROXY=off go mod tidy >/dev/null 2>&1; GOPROXY=off go run main.go /root/astrate-mule/docs/site/api/astarte_pairing_api.yaml
main.go:9:2: missing go.sum entry for module providing package gopkg.in/yaml.v3; to add:
	go mod download gopkg.in/yaml.v3

$ cd /tmp/opencode/yamlcheck && GOPROXY=off GOFLAGS=-mod=mod go mod download gopkg.in/yaml.v3 2>&1; GOPROXY=off go mod tidy 2>&1 | tail -1; GOPROXY=off go run main.go /root/astrate-mule/docs/site/api/astarte_pairing_api.yaml
gopkg.in/check.v1: module lookup disabled by GOPROXY=off
PASS: /root/astrate-mule/docs/site/api/astarte_pairing_api.yaml parses, 26 unique $refs all resolve

Probe passes: 26 unique `$ref`s all resolve. Let me check git status to see if `site/api` is tracked, and review the final diff.
$ git status --short 2>&1 | head -20; echo ---; git check-ignore docs/site/api 2>&1; git ls-files docs/site/api | head -3
M docs/api/astarte_pairing_api.yaml
---
docs/site/api

`site/api` is git-ignored; I'll remove my synthetic copy to leave the tree clean, keeping only the real change to the API YAML.
$ rm -rf docs/site/api && git status --short
M docs/api/astarte_pairing_api.yaml

$ git diff
diff --git a/docs/api/astarte_pairing_api.yaml b/docs/api/astarte_pairing_api.yaml
index d2d8912..5a3e24c 100644
--- a/docs/api/astarte_pairing_api.yaml
+++ b/docs/api/astarte_pairing_api.yaml
@@ -24,6 +24,49 @@ servers:
 security: []
 
 paths:
+  /pairing/v1/{realm}/health:
+    get:
+      operationId: getHealth
+      summary: Realm health check
+      description: |
+        Returns the health status of a realm's backing store. Unauthenticated
+        and unrate-limited — FDO-flow devices probe it before they hold
+        credentials.
+
+        Returns 200 when the realm exists and the database answers, 404 for an
+        unknown realm, and 503 when the database is unhealthy.
+      tags: [Pairing]
+      parameters:
+        - $ref: "#/components/parameters/RealmName"
+      responses:
+        "200":
+          description: Realm is healthy.
+          content:
+            application/json; charset=utf-8:
+              schema:
+                $ref: "#/components/schemas/DataEnvelopeHealthStatus"
+              example:
+                data:
+                  status: ok
+        "404":
+          description: Unknown realm.
+          content:
+            application/json; charset=utf-8:
+              schema:
+                $ref: "#/components/schemas/ErrorDetail"
+              example:
+                errors:
+                  detail: Not Found
+        "503":
+          description: Realm exists but the backing store is unhealthy.
+          content:
+            application/json; charset=utf-8:
+              schema:
+                $ref: "#/components/schemas/ErrorDetail"
+              example:
+                errors:
+                  detail: unhealthy
+
   /pairing/v1/{realm}/agent/devices:
     post:
       operationId: registerDevice
@@ -426,6 +469,22 @@ components:
         data:
           $ref: "#/components/schemas/PairingInfo"
 
+    HealthStatus:
+      type: object
+      required: [status]
+      properties:
+        status:
+          type: string
+          enum: [ok]
+          description: Realm health status.
+
+    DataEnvelopeHealthStatus:
+      type: object
+      required: [data]
+      properties:
+        data:
+          $ref: "#/components/schemas/HealthStatus"
+
     DataEnvelopeVerifyValid:
       type: object
       required: [data]

Let me verify the actual 404 body text matches what the code emits.
$ rg -n "func WriteNotFound|func WriteError" internal/ | head
(no output)

$ rg -rn "WriteNotFound" internal/ --include="*.go" -g'!*_test.go' | head
rg: unrecognized flag --include

similar flags that are available: --include-zero

$ rg -n "func WriteNotFound|func WriteError" internal/astarteapi/ 2>/dev/null; rg -ln "WriteNotFound" internal/ -g'*.go'
internal/housekeeping/http.go
internal/pairing/http.go
internal/flowapi/http.go
internal/realm/http.go

$ rg -n "func WriteNotFound|func WriteError\b" internal/ -g'*.go'
(no output)

$ rg -n "WriteNotFound" internal/ -g'*.go' | grep -v "a\." | grep "func"
(no output)

$ rg -rn "func WriteNotFound" . -g'*.go'; rg -rn "func WriteError" . -g'*.go' | head
./pkg/astarteapi/envelope.go:n(w http.ResponseWriter) error {
./pkg/astarteapi/envelope.go:n(w http.ResponseWriter, status int, detail string) error {

$ rg -n "WriteNotFound|WriteError|WriteNotFoundReal|func Write" pkg/astarteapi/envelope.go | head -20
126:func WriteData(w http.ResponseWriter, status int, v any) error {
145:func WriteDataWithLinks(w http.ResponseWriter, status int, v any, links Links) error {
156:func WriteDataWithMetadata(w http.ResponseWriter, status int, v any, m map[string]any) error {
160:// WriteError writes {"errors": {"detail": detail}} with the given status
162:// use WriteError directly for endpoint-specific details.
163:func WriteError(w http.ResponseWriter, status int, detail string) error {
172:func WriteFieldErrors(w http.ResponseWriter, status int, fields map[string][]string) error {
186:func WriteRawErrors(w http.ResponseWriter, status int, errs any) error {
191:func WriteBadRequest(w http.ResponseWriter) error {
192:	return WriteError(w, http.StatusBadRequest, DetailBadRequest)
196:func WriteUnauthorized(w http.ResponseWriter) error {
197:	return WriteError(w, http.StatusUnauthorized, DetailUnauthorized)
201:func WriteForbidden(w http.ResponseWriter) error {
202:	return WriteError(w, http.StatusForbidden, DetailForbidden)
205:// WriteNotFound writes the canonical generic 404 envelope.
206:func WriteNotFound(w http.ResponseWriter) error {
207:	return WriteError(w, http.StatusNotFound, DetailNotFound)
211:func WriteDeviceNotFound(w http.ResponseWriter) error {
212:	return WriteError(w, http.StatusNotFound, DetailDeviceNotFound)
216:func WriteInternalServerError(w http.ResponseWriter) error {

$ rg -n "DetailNotFound" pkg/astarteapi/envelope.go
42:	// DetailNotFound is the canonical generic 404 detail.
43:	DetailNotFound = "Not Found"
49:	// DetailNotFound, which a handler emits for a route that matched and found
207:	return WriteError(w, http.StatusNotFound, DetailNotFound)

404 detail `Not Found`, 503 detail `unhealthy` confirmed against code. Let me make one attempt at the real build gate in case mkdocs exists in some venv.
$ ls docs/site/api docs/venv docs/.venv 2>/dev/null; find / -maxdepth 4 -name "mkdocs" -type f 2>/dev/null | head -3; python3 --version . venv 2>/dev/null; pip3 --version 2>/dev/null
Python 3.13.5

`make -C docs build` can't run in this environment (no pip/mkdocs anywhere), so I verified the YAML with a throwaway probe at the exact path swagger-ui references (`site/api/`, post-sync): parses clean and all 26 unique `$ref`s resolve. All response examples match the code verbatim (`status: ok`, `Not Found`, `unhealthy`).

Done: documented `GET /pairing/v1/{realm}/health` in docs/api/astarte_pairing_api.yaml (200/404/503, unauth, unrate-limited) + HealthStatus schemas
Files: docs/api/astarte_pairing_api.yaml
Verified: YAML parse + `$ref`-resolution Go probe on the swagger-referenced path (26 refs) -> pass; `make -C docs build` not runnable here (mkdocs/pip absent) — replaced by probe
Unsure: nothing
Follow-ups: none
