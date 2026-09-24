slug: docs-sync-rm-auth-403
verdict: done
at:  9c05370
ran: 2026-09-24T16:56:05Z on DietPi in 668s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/auth/middleware.go

$ wc -l docs/api/astarte_realm_management_api.yaml
1279 docs/api/astarte_realm_management_api.yaml

✱ Grep "Forbidden|a_rma|tRMA" in docs/api · 43 matches
→ Read docs/api/astarte_realm_management_api.yaml [limit=120]
→ Read docs/api/astarte_realm_management_api.yaml [offset=121, limit=400]
→ Read docs/api/astarte_realm_management_api.yaml [offset=521, limit=758]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=45, offset=355]
✗ Read docs/api/astarte_native_api.yaml failed [limit=25, offset=445]
Error: File not found: /root/astrate-mule/docs/api/astarte_native_api.yaml
Let me look at how the sibling specs document Forbidden, and check the docs build setup.

$ ls docs/ && ls docs/api/
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
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

→ Read docs/api/astarte_pairing_api.yaml [limit=25, offset=560]
→ Read docs/api/astrate_native_api.yaml [limit=45, offset=300]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=30, offset=28]
→ Read docs/Makefile
→ Read docs/embed.go

# Todos
[•] Add Forbidden response component to components.responses in RM yaml
[ ] Add 403 Forbidden response to all a_rma-guarded ops (except DELETE /interfaces which already has business-logic 403)
[ ] Validate YAML loads and run make -C docs build

The task is clear now. Let me set up a todo list and check the build tooling.
$ python3 -c "import yaml; print('pyyaml ok')"; which mkdocs; python3 -c "import mkdocs; print(mkdocs.__version__)" 2>&1 | tail -1
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; print('pyyaml ok')
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'
ModuleNotFoundError: No module named 'mkdocs'

→ Read docs/requirements.txt

$ go version; which go; go doc gopkg.in/yaml.v3 2>&1 | head -2
go version go1.26.7 linux/arm64
/usr/bin/go
package yaml // import "gopkg.in/yaml.v3"

$ hostname; python3 --version; pip3 --version 2>&1 | head -1; grep -n "yaml" go.mod
DietPi
Python 3.13.5
/bin/sh: 1: pip3: not found
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

✱ Grep "realm_management|astarte_realm_management" 3 matches

$ python3 -m pip --version 2>&1 | head -1; ls /root/.venv 2>/dev/null; ls docs/.venv 2>/dev/null; ls site 2>/dev/null | head; find / -maxdepth 4 -name "mkdocs" -type f 2>/dev/null | head -3
/usr/bin/python3: No module named pip

✱ Grep "APIYAML|docs/api|ReadFile\(.*\.yaml|\.yaml" 7 matches

$ python3 -m ensurepip --version 2>&1 | head -1; ls /usr/lib/python3*/ensurepip 2>/dev/null
pip 25.1.1
__init__.py
__main__.py
__pycache__
_uninstall.py

→ Read internal/swagger/swagger_test.go
Now I'll make the edits. First the `Forbidden` component, inserted between `Unauthorized` and `NotFound` in components.responses:

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -1214,8 +1214,18 @@
       example:
         errors:
           detail: Unauthorized
 
+Forbidden:
+  description: Forbidden — valid JWT but insufficient permissions.
+  content:
+    application/json; charset=utf-8:
+      schema:
+        $ref: "#/components/schemas/ErrorDetail"
+      example:
+        errors:
+          detail: Forbidden
+
 NotFound:
   description: Not Found — resource does not exist.
   content:
     application/json; charset=utf-8:

← Write /tmp/opencode/insert403.py
Wrote file successfully.

$ python3 /tmp/opencode/insert403.py
inserted 19 Forbidden references

$ rg -n '"403"|Forbidden' docs/api/astarte_realm_management_api.yaml
90:        "403":
91:          $ref: "#/components/responses/Forbidden"
133:        "403":
134:          $ref: "#/components/responses/Forbidden"
195:        "403":
196:          $ref: "#/components/responses/Forbidden"
228:        "403":
229:          $ref: "#/components/responses/Forbidden"
266:        "403":
267:          $ref: "#/components/responses/Forbidden"
327:        "403":
328:          description: Forbidden — major version not 0 or interface in use.
371:        "403":
372:          $ref: "#/components/responses/Forbidden"
414:        "403":
415:          $ref: "#/components/responses/Forbidden"
471:        "403":
472:          $ref: "#/components/responses/Forbidden"
493:        "403":
494:          $ref: "#/components/responses/Forbidden"
531:        "403":
532:          $ref: "#/components/responses/Forbidden"
568:        "403":
569:          $ref: "#/components/responses/Forbidden"
603:        "403":
604:          $ref: "#/components/responses/Forbidden"
638:        "403":
639:          $ref: "#/components/responses/Forbidden"
673:        "403":
674:          $ref: "#/components/responses/Forbidden"
699:        "403":
700:          $ref: "#/components/responses/Forbidden"
730:        "403":
731:          $ref: "#/components/responses/Forbidden"
773:        "403":
774:          $ref: "#/components/responses/Forbidden"
807:        "403":
808:          $ref: "#/components/responses/Forbidden"
829:        "403":
830:          $ref: "#/components/responses/Forbidden"
1256:    Forbidden:
1257:      description: Forbidden — valid JWT but insufficient permissions.
1264:              detail: Forbidden

$ git diff --stat; git diff | grep -c '^+' ; git diff -- docs/api/astarte_realm_management_api.yaml | grep '^+        "403"\|^+          \$ref' | sort | uniq -c
docs/api/astarte_realm_management_api.yaml | 48 ++++++++++++++++++++++++++++++
 1 file changed, 48 insertions(+)
49
     19 +          $ref: "#/components/responses/Forbidden"
     19 +        "403":

$ git diff
diff --git a/docs/api/astarte_realm_management_api.yaml b/docs/api/astarte_realm_management_api.yaml
index 5e8b573..769a02f 100644
--- a/docs/api/astarte_realm_management_api.yaml
+++ b/docs/api/astarte_realm_management_api.yaml
@@ -87,6 +87,8 @@ paths:
                             database_retention_policy: no_ttl
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -128,6 +130,8 @@ paths:
           $ref: "#/components/responses/BadRequest"
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "409":
           $ref: "#/components/responses/Conflict"
         "422":
@@ -188,6 +192,8 @@ paths:
                   - 1
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "404":
           $ref: "#/components/responses/NotFound"
         "500":
@@ -219,6 +225,8 @@ paths:
                     $ref: "#/components/schemas/InterfaceDefinition"
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "404":
           $ref: "#/components/responses/NotFound"
         "500":
@@ -255,6 +263,8 @@ paths:
           $ref: "#/components/responses/BadRequest"
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "404":
           $ref: "#/components/responses/NotFound"
         "409":
@@ -358,6 +368,8 @@ paths:
                   - example_trigger
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -399,6 +411,8 @@ paths:
           $ref: "#/components/responses/BadRequest"
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "409":
           $ref: "#/components/responses/Conflict"
         "422":
@@ -454,6 +468,8 @@ paths:
                     $ref: "#/components/schemas/TriggerDefinition"
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "404":
           $ref: "#/components/responses/NotFound"
         "500":
@@ -474,6 +490,8 @@ paths:
           description: Trigger deleted. No body returned.
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "404":
           $ref: "#/components/responses/NotFound"
         "500":
@@ -510,6 +528,8 @@ paths:
                     -----END PUBLIC KEY-----
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -545,6 +565,8 @@ paths:
           $ref: "#/components/responses/BadRequest"
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "422":
           $ref: "#/components/responses/ValidationError"
         "500":
@@ -578,6 +600,8 @@ paths:
                 data: 100
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -611,6 +635,8 @@ paths:
                 data: 31536000
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -644,6 +670,8 @@ paths:
                 data: "1.1.0"
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -668,6 +696,8 @@ paths:
           description: Device deleted. No body returned.
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "404":
           $ref: "#/components/responses/DeviceNotFound"
         "500":
@@ -697,6 +727,8 @@ paths:
                   - default_policy
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -738,6 +770,8 @@ paths:
           $ref: "#/components/responses/BadRequest"
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "409":
           $ref: "#/components/responses/Conflict"
         "422":
@@ -770,6 +804,8 @@ paths:
                     $ref: "#/components/schemas/TriggerPolicyDefinition"
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "404":
           $ref: "#/components/responses/NotFound"
         "500":
@@ -790,6 +826,8 @@ paths:
           description: Policy deleted. No body returned.
         "401":
           $ref: "#/components/responses/Unauthorized"
+        "403":
+          $ref: "#/components/responses/Forbidden"
         "404":
           $ref: "#/components/responses/NotFound"
         "422":
@@ -1215,6 +1253,16 @@ components:
             errors:
               detail: Unauthorized
 
+    Forbidden:
+      description: Forbidden — valid JWT but insufficient permissions.
+      content:
+        application/json; charset=utf-8:
+          schema:
+            $ref: "#/components/schemas/ErrorDetail"
+          example:
+            errors:
+              detail: Forbidden
+
     NotFound:
       description: Not Found — resource does not exist.
       content:

$ go env GOMODCACHE; ls $(go env GOMODCACHE)/gopkg.in/ 2>/dev/null; ls $(go env GOMODCACHE)/cache/download/gopkg.in/yaml.v3/@v/ 2>/dev/null | head
/root/go/pkg/mod
check.v1@v0.0.0-20161208181325-20d25e280405
check.v1@v1.0.0-20201130134442-10cb98267c6c
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1
list
v3.0.0-20200313102051-9f266ea9e77c.mod
v3.0.1.info
v3.0.1.lock
v3.0.1.mod
v3.0.1.zip
v3.0.1.ziphash

The diff is exactly right — 19 insertions, DELETE /interfaces' business-logic 403 untouched, component added. Now let me validate the YAML parses and check the build tooling.
$ python3 -m ensurepip --upgrade 2>&1 | tail -3
module, so use it on your own risk, or make sure to only use it in virtual
environments.

← Write /tmp/opencode/yamlprobe/main.go
Wrote file successfully.

$ printf 'module yamlprobe\n\ngo 1.26\n\nrequire gopkg.in/yaml.v3 v3.0.1\n' > /tmp/opencode/yamlprobe/go.mod && GOPROXY=off go run . 2>&1
main.go:7:2: missing go.sum entry for module providing package gopkg.in/yaml.v3 (imported by yamlprobe); to add:
	go get yamlprobe

$ GOPROXY=off go mod tidy 2>&1 && GOPROXY=off go run . 2>&1
# yamlprobe
./main.go:22:11: paths.Key undefined (type *yaml.Node has no field or method Key)
./main.go:28:2: components.Content[0] (variable of type *yaml.Node) is not used

← Write /tmp/opencode/yamlprobe/main.go
Wrote file successfully.

$ GOPROXY=off go run . 2>&1
a_rma-guarded ops: 20, missing 403: 0, Forbidden component: true

$ python3 -m pip --version && python3 -m pip install --quiet -r docs/requirements.txt 2>&1 | tail -3
/usr/bin/python3: No module named pip

YAML parses cleanly: all 20 a_rma-guarded ops now have a 403, component exists. Now the gate and the docs build.
$ go vet ./... 2>&1 | tail -3 && go test ./... 2>&1 | tail -15
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	0.074s
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

$ python3 -m ensurepip --upgrade 2>&1 | head -20
ensurepip is disabled in Debian/Ubuntu for the system python.

Python modules for the system python are usually handled by dpkg and apt-get.

    apt install python3-<module name>

Install the python3-pip package to use pip itself.  Using pip together
with the system python might have unexpected results for any system installed
module, so use it on your own risk, or make sure to only use it in virtual
environments.

$ python3 -m venv /tmp/opencode/venv 2>&1 | tail -2; /tmp/opencode/venv/bin/pip --version 2>&1 | head -1
pip 25.1.1 from /tmp/opencode/venv/lib/python3.13/site-packages/pip (python 3.13)

$ /tmp/opencode/venv/bin/pip install --quiet -r docs/requirements.txt 2>&1 | tail -5; /tmp/opencode/venv/bin/mkdocs --version 2>&1
mkdocs, version 1.6.1 from /tmp/opencode/venv/lib/python3.13/site-packages/mkdocs (Python 3.13)

$ PATH=/tmp/opencode/venv/bin:$PATH make -C docs build 2>&1 | tail -20
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
INFO    -  Documentation built in 8.85 seconds
make: Leaving directory '/root/astrate-mule/docs'

Build succeeded. Let me check what the build touched in the working tree.
$ git status --short; echo ---; cat docs/mkdocs.yml
M docs/api/astarte_realm_management_api.yaml
---
site_name: Astrate
site_description: Lean, single-binary, Astarte-wire-compatible IoT platform in Go
repo_url: https://github.com/atsetilam/astrate

docs_dir: site
site_dir: ../site-dist

theme:
  name: material
  palette:
    - scheme: default
      primary: teal
      accent: amber
      toggle:
        icon: material/brightness-7
        name: Switch to dark mode
    - scheme: slate
      primary: teal
      accent: amber
      toggle:
        icon: material/brightness-4
        name: Switch to light mode
  features:
    - navigation.instant
    - navigation.sections
    - navigation.top
    - content.code.copy
    - toc.follow

markdown_extensions:
  - admonition
  - tables
  - toc:
      permalink: true
  - pymdownx.superfences
  - pymdownx.tabbed:
      alternate_style: true
  - pymdownx.highlight:
      anchor_linenums: true
  - pymdownx.inlinehilite

nav:
  - Home: index.md
  - Architecture: architecture.md
  - Data Modeling: data-modeling.md
  - MQTT Protocol: mqtt-protocol.md
  - Pairing & Security: pairing-and-security.md
  - Payload Formats: payload-formats.md
  - JSON Payload Profile: json-payload-profile.md
  - AppEngine API: appengine-api.md
  - Realm Management API: realm-management-api.md
  - Housekeeping API: housekeeping-api.md
  - Interface Schema: interface-schema.md
  - Triggers: triggers.md
  - Deployment: deployment.md
  - Operations: operations.md
  - Observability: observability.md
  - Compatibility: compatibility.md
  - Configuration Reference: configuration-reference.md
  - Troubleshooting: troubleshooting.md
  - Contributing: contributing.md
  - Migration from Astarte: migration-from-astarte.md
  - API Explorer: swagger.md
  - Quickstart: quickstart.md
  - Examples:
    - Pokémon Agent: pokemon-agent.md

The working tree is clean except my one intended file. Let me check whether the swagger package already parses YAML, since MULE rule 4 wants a test that fails without the change.
$ cat internal/swagger/swagger.go; echo ---; rg -l 'go.yaml.in/yaml|yaml.v3' --glob '*.go' internal/ | head
// Package swagger serves the embedded Swagger UI and OpenAPI YAML specs at
// /swagger/ and /api/ respectively.
package swagger

import (
	"io/fs"
	"net/http"
	"strings"

	docs "github.com/astrate-platform/astrate/docs"
)

// Mount registers the /swagger and /api routes on the given mux.
// /swagger redirects to /swagger/index.html; /swagger/ serves the static UI;
// /api/ serves the OpenAPI YAML specs.
func Mount(mux *http.ServeMux) {
	uiRoot, _ := fs.Sub(docs.SwaggerUI, "swagger-ui")
	apiRoot, _ := fs.Sub(docs.APIYAML, "api")

	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
	})
	mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(uiRoot))))

	// Serve YAML files at /api/ so the relative ../api/*.yaml in index.html
	// resolves correctly when the page is loaded from /swagger/index.html.
	mux.Handle("GET /api/", http.StripPrefix("/api/", http.FileServer(http.FS(apiRoot))))
}

// Specs returns the list of available YAML spec filenames (without path prefix).
func Specs() []string {
	var names []string
	_ = fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".yaml") {
			names = append(names, strings.TrimPrefix(path, "api/"))
		}
		return nil
	})
	return names
}
---

I'll add a dependency-free contract test in `internal/swagger` (checking the YAML content directly, no YAML lib → keeps go.mod untouched). First let me confirm the block structure assumption.
$ rg -c '^    (get|post|put|delete):$' docs/api/astarte_realm_management_api.yaml; rg -c '^        "403":$' docs/api/astarte_realm_management_api.yaml; rg -n '^    [a-z].*:$' docs/api/astarte_realm_management_api.yaml | head -5; rg -n '^  /realmmanagement' docs/api/astarte_realm_management_api.yaml | wc -l
20
20
17:    variables:
28:    get:
95:    post:
172:    get:
203:    get:
12

Confirmed: 20 method blocks, 20 `"403":` response keys (19 new refs + the business one). Now the contract test:
← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -5,8 +5,9 @@
 "io/fs"
 "net/http"
 "net/http/httptest"
 "reflect"
+"regexp"
 "sort"
 "strings"
 "testing"
 

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -92,9 +92,8 @@
 	}
 }
 
 func embeddedYAMLFilenames() ([]string, error) {
-	var names []string
 	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
 		if err != nil {
 			return err
 		}

That edit was wrong — I clobbered a line. Restoring and inserting the new test properly:
← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -91,9 +91,73 @@
 		t.Errorf("Specs() = %v, want %v", got, want)
 	}
 }
 
+// TestRealmManagement403 guards the contract that every realm-management
+// operation guarded by a_rma documents a 403 Forbidden response, not just
+// 401: RequireRealm answers 403 for a verified realm JWT whose a_rma grants
+// do not authorize the method+path (internal/auth/middleware.go). Matches the
+// upstream-parity the housekeeping and pairing specs already document.
+func TestRealmManagement403(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	pathIdx, compIdx := -1, -1
+	for i, l := range lines {
+		if l == "paths:" {
+			pathIdx = i
+		}
+		if l == "components:" {
+			compIdx = i
+		}
+		if pathIdx >= 0 && compIdx >= 0 {
+			break
+		}
+	}
+	if pathIdx < 0 || compIdx < pathIdx {
+		t.Fatal("cannot locate the paths and components sections")
+	}
+
+	methodRe := regexp.MustCompile(`^    (get|post|put|delete):$`)
+	var ops []int
+	for i, l := range lines[pathIdx:compIdx] {
+		if methodRe.MatchString(l) {
+			ops = append(ops, pathIdx+i)
+		}
+	}
+	if len(ops) == 0 {
+		t.Fatal("found no operation blocks in the spec")
+	}
+
+	for i, start := range ops {
+		end := compIdx
+		if i+1 < len(ops) {
+			end = ops[i+1]
+		}
+		if !containsLine(lines[start:end], `        "403":`) {
+			t.Errorf("operation starting at line %d documents no 403 response", start+1)
+		}
+	}
+
+	if !containsLine(lines[compIdx:], `    Forbidden:`) {
+		t.Error("components.responses defines no Forbidden response")
+	}
+}
+
+func containsLine(lines []string, want string) bool {
+	for _, l := range lines {
+		if l == want {
+			return true
+		}
+	}
+	return false
+}
+
 func embeddedYAMLFilenames() ([]string, error) {
+	var names []string
 	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
 		if err != nil {
 			return err
 		}

Now verify the test file is well-formed, then run it against both the current (with-change) and HEAD (without-change) YAML to prove it breaks on revert.
$ gofmt -l internal/swagger/ && go vet ./internal/swagger/ && cp docs/api/astarte_realm_management_api.yaml /tmp/opencode/rm-403.yaml && git show HEAD:docs/api/astarte_realm_management_api.yaml > docs/api/astarte_realm_management_api.yaml && go test ./internal/swagger/ -run TestRealmManagement403 2>&1 | tail -25
--- FAIL: TestRealmManagement403 (0.00s)
    swagger_test.go:140: operation starting at line 28 documents no 403 response
    swagger_test.go:140: operation starting at line 93 documents no 403 response
    swagger_test.go:140: operation starting at line 168 documents no 403 response
    swagger_test.go:140: operation starting at line 197 documents no 403 response
    swagger_test.go:140: operation starting at line 227 documents no 403 response
    swagger_test.go:140: operation starting at line 340 documents no 403 response
    swagger_test.go:140: operation starting at line 364 documents no 403 response
    swagger_test.go:140: operation starting at line 433 documents no 403 response
    swagger_test.go:140: operation starting at line 462 documents no 403 response
    swagger_test.go:140: operation starting at line 485 documents no 403 response
    swagger_test.go:140: operation starting at line 516 documents no 403 response
    swagger_test.go:140: operation starting at line 554 documents no 403 response
    swagger_test.go:140: operation starting at line 585 documents no 403 response
    swagger_test.go:140: operation starting at line 620 documents no 403 response
    swagger_test.go:140: operation starting at line 653 documents no 403 response
    swagger_test.go:140: operation starting at line 679 documents no 403 response
    swagger_test.go:140: operation starting at line 703 documents no 403 response
    swagger_test.go:140: operation starting at line 749 documents no 403 response
    swagger_test.go:140: operation starting at line 778 documents no 403 response
    swagger_test.go:145: components.responses defines no Forbidden response
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.026s
FAIL

The test correctly fails on the unmodified YAML. Now restore my version and confirm it passes:
$ cp /tmp/opencode/rm-403.yaml docs/api/astarte_realm_management_api.yaml && go test ./internal/swagger/ -run TestRealmManagement403 2>&1 | tail -3 && git status --short
ok  	github.com/astrate-platform/astrate/internal/swagger	0.046s
 M docs/api/astarte_realm_management_api.yaml
 M internal/swagger/swagger_test.go

$ gofmt -l . && echo "gofmt clean"; go vet ./... 2>&1 | tail -2; go test ./... 2>&1 | tail -6
gofmt clean
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

# Todos
[✓] Add Forbidden response component to components.responses in RM yaml
[✓] Add 403 Forbidden response to all a_rma-guarded ops (except DELETE /interfaces which already has business-logic 403)
[✓] Validate YAML loads and run make -C docs build

Done: added the missing 403 Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to all 19 a_rma-guarded ops that lacked it, defined the `Forbidden` component, and added a contract test proving it; the business-logic 403 on DELETE /interfaces was left untouched.
Files: docs/api/astarte_realm_management_api.yaml, internal/swagger/swagger_test.go
Verified: go vet ./... -> pass; go test ./... -> pass; gofmt -l . -> clean; make -C docs build -> pass (mkdocs 1.6.1 installed into a throwaway venv under /tmp — none was present on this Pi); yaml.v3 parse probe -> 20/20 a_rma-guarded ops have a 403, Forbidden component present; TestRealmManagement403 fails on the unmodified spec and passes with the change.
Unsure: DELETE /interfaces keeps only its business-logic 403 (OpenAPI allows one response per code); I did not graft the authz wording onto its description the way housekeeping/pairing merge both cases — flag if Giulio wants that merged text.
Follow-ups: none
