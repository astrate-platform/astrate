slug: docs-sync-hk-delete-gating-responses
verdict: done
at:  d07b6c7
ran: 2026-09-13T18:50:42Z on DietPi in 331s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

# Todos
[•] Read the OpenAPI spec for housekeeping DELETE /realms/{realm} and PATCH
[ ] Read internal/housekeeping/http.go to confirm writeError paths for 405/422
[ ] Read service.go to confirm DeleteRealm returns ErrDeletionDisabled and ErrConnectedDevicesPresent
[ ] Fix PATCH: remove incorrect 405 response
[ ] Add 405 response to DELETE using MethodNotAllowed component
[ ] Add 422 response to DELETE with FieldErrors shape + connected_devices_present example
[ ] Run make -C docs build and verify

→ Read docs/api/astarte_housekeeping_api.yaml
→ Read internal/housekeeping/http.go [offset=220, limit=50]
→ Read internal/housekeeping/service.go [offset=240, limit=40]
Now let me confirm `UpdateRealm` in service.go doesn't have a deletion-gating path:
✱ Grep "func \(s \*Service\) UpdateRealm" in internal/housekeeping · 1 match
→ Read internal/housekeeping/service.go [offset=196, limit=30]
Confirmed: `UpdateRealm` (service.go:196-216) returns only `ErrValidation` or store errors — no deletion-gating path. Now let me check the docs build setup:
→ Read docs/Makefile

$ ls /root/astrate-mule/docs
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

→ Read docs/embed.go

$ which mkdocs swagger-cli swagger-ui 2>/dev/null; ls /root/astrate-mule/docs/site | head; python3 -c "import mkdocs" 2>&1
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
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import mkdocs
ModuleNotFoundError: No module named 'mkdocs'

→ Read docs/requirements.txt

$ ls /root/astrate-mule/docs/site/api 2>/dev/null; ls /root/astrate-mule/venv 2>/dev/null; python3 -m pip --version 2>&1 | head -1
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
/usr/bin/python3: No module named pip

$ python3 -c "import yaml; print(yaml.__version__)" 2>&1; ls /etc/os-release 2>/dev/null && . /etc/os-release && echo $ID $VERSION_ID
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; print(yaml.__version__)
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'
/etc/os-release
debian 13

✱ Grep "APIYAML|astarte_housekeeping_api|swagger" in . · 30 matches

$ which ruby node go; go version 2>/dev/null; grep -c "yaml.v3" /root/astrate-mule/go.mod 2>/dev/null
/usr/local/bin/node
/usr/bin/go
go version go1.26.7 linux/arm64
1

$ grep -n "yaml" /root/astrate-mule/go.mod
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

$ which uv pip3 python3.12 python3.13 2>/dev/null; ls /root/.uv 2>/dev/null; go env GOMODCACHE
/usr/bin/python3.13
/root/go/pkg/mod

$ ls /root/go/pkg/mod/github.com/getkin 2>/dev/null; ls /root/go/pkg/mod/gopkg.in 2>/dev/null; ls /root/go/pkg/mod | head -30
check.v1@v0.0.0-20161208181325-20d25e280405
check.v1@v1.0.0-20201130134442-10cb98267c6c
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1
4d63.com
cache
charm.land
codeberg.org
dario.cat
dev.gaijin.team
github.com
gitlab.com
go-simpler.org
go.augendre.info
go.etcd.io
go.mongodb.org
go.opentelemetry.io
go.uber.org
go.yaml.in
golang.org
google.golang.org
gopkg.in
gotest.tools
honnef.co
mvdan.cc
pgregory.net

No pip/uv/mkdocs available. I'll make the YAML edits first, then validate the spec by loading it with a throwaway Go probe (yaml.v3 is in the module cache).

Now the edits. First, remove the wrong 405 from PATCH:
← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -178,10 +178,8 @@
 "403":
   $ref: "#/components/responses/Forbidden"
 "404":
   $ref: "#/components/responses/NotFound"
-"405":
-  $ref: "#/components/responses/MethodNotAllowed"
 "422":
   $ref: "#/components/responses/ValidationErrors"
 "500":
   $ref: "#/components/responses/InternalServerError"

Now add 405 and 422 to DELETE:
← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -201,8 +201,12 @@
         "403":
           $ref: "#/components/responses/Forbidden"
         "404":
           $ref: "#/components/responses/NotFound"
+        "405":
+          $ref: "#/components/responses/MethodNotAllowed"
+        "422":
+          $ref: "#/components/responses/UnprocessableDeletion"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
 components:

Now add the `UnprocessableDeletion` response component next to `ValidationErrors`:
← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -397,8 +397,19 @@
         errors:
           error_name:
             - invalid_update_parameters
 
+UnprocessableDeletion:
+  description: Unprocessable Entity — realm has devices still connected and cannot be deleted.
+  content:
+    application/json; charset=utf-8:
+      schema:
+        $ref: "#/components/schemas/FieldErrorsDetail"
+      example:
+        errors:
+          error_name:
+            - connected_devices_present
+
 Conflict:
   description: Conflict — resource already exists (e.g. realm name taken).
   content:
     application/json; charset=utf-8:

→ Read docs/api/astarte_housekeeping_api.yaml [offset=138, limit=85]
Now let me run the build as the task requires:

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

! permission requested: external_directory (/tmp/*, /tmp/specprobe/*); auto-rejecting
✗ rm -rf /tmp/specprobe && mkdir -p /tmp/specprobe && printf 'module specprobe\ngo 1.26\n\nrequire gopkg.in/yaml.v3 v3.0.1\n' > /tmp/specprobe/go.mod && cp /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml /tmp/specprobe/spec.yaml failed
Error: The user rejected permission to use this specific tool call.
`mkdocs` is not installed on this box, and there's no pip/uv to install it. The `sync` half of the build ran and pushed the updated spec into `site/api`. I'll confirm the spec still loads with a throwaway Go probe (yaml.v3 is in the module cache), which resolves the refs Swagger would follow:
