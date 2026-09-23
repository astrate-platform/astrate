slug: docs-sync-rm-put-interface-409
verdict: done
at:  7256af9
ran: 2026-09-23T20:17:46Z on DietPi in 356s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_realm_management_api.yaml [offset=100, limit=160]
✱ Grep "Conflict|BadRequest|NotFound" in docs/api/astarte_realm_management_api.yaml · 91 matches
→ Read docs/api/astarte_realm_management_api.yaml [limit=75, offset=1098]
✗ Read internal/realmmanagement/http.go failed [offset=360, limit=35]
Error: File not found: /root/astrate-mule/internal/realmmanagement/http.go
✗ Read internal/realmmanagement/service.go failed [offset=195, limit=40]
Error: File not found: /root/astrate-mule/internal/realmmanagement/service.go

$ rg -l "ErrNameMismatch|checkMinorUpgrade" --glob '*.go' | head -30
internal/realm/http.go
internal/realm/service.go

✱ Grep "ErrNameMismatch|ErrMajorMismatch|CheckMinorUpgrade|ErrAlreadyExists|StatusConflict" in internal/realm/http.go · 40 matches
✱ Grep "ErrNameMismatch|ErrMajorMismatch|CheckMinorUpgrade" in internal/realm/service.go · 10 matches
→ Read internal/realm/http.go [limit=40, offset=355]
→ Read internal/realm/http_test.go [limit=55, offset=388]
✱ Grep "bodyNameMismatch|bodyMajorMismatch" in internal/realm/http_test.go · 4 matches
→ Read docs/api/astarte_realm_management_api.yaml [limit=30]
→ Read docs/Makefile

$ ls docs/; echo ---; cat docs/requirements.txt; echo ---; rg -l "astarte_realm_management_api|realm_management_api" docs/ bench/ tools/ 2>/dev/null | head -20
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
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10
---
docs/swagger-ui/index.html
docs/handoff/phase-1-memory.md
docs/site/swagger-ui/index.html

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -227,8 +227,10 @@
 "401":
   $ref: "#/components/responses/Unauthorized"
 "404":
   $ref: "#/components/responses/NotFound"
+"409":
+  $ref: "#/components/responses/Conflict"
 "422":
   $ref: "#/components/responses/ValidationError"
 "500":
   $ref: "#/components/responses/InternalServerError"

Now let me make the edits. First, add the 409 response to the PUT operation:
← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -1142,16 +1142,29 @@
         errors:
           detail: Device not found
 
 Conflict:
-  description: Conflict — resource already exists.
+  description: Conflict — resource already exists or the submitted data disagrees with the resource.
   content:
     application/json; charset=utf-8:
       schema:
         $ref: "#/components/schemas/ErrorDetail"
-      example:
-        errors:
-          detail: Already exists
+      examples:
+        alreadyExists:
+          summary: The resource already exists.
+          value:
+            errors:
+              detail: Already exists
+        nameMismatch:
+          summary: The interface name in the body does not match the one in the URL.
+          value:
+            errors:
+              detail: Interface name doesn't match the one in the interface json
+        majorMismatch:
+          summary: The interface major version in the body does not match the one in the URL.
+          value:
+            errors:
+              detail: Interface major version doesn't match the one in the interface json
 
 ValidationError:
   description: Unprocessable Entity — request body failed validation.
   content:

$ which python3; python3 -c "import yaml; print('yaml ok')" 2>&1; python3 -m mkdocs --version 2>&1 | head -2
/usr/bin/python3
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; print('yaml ok')
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'
/usr/bin/python3: No module named mkdocs

$ python3 -m pip --version 2>&1 | head -1; python3 -m venv --help >/dev/null 2>&1 && echo "venv ok" || echo "no venv"
/usr/bin/python3: No module named pip
venv ok

$ go version 2>&1; rg -n "yaml" go.mod; rg -n "gopkg.in/yaml|sigs.k8s.io/yaml|url:" docs/swagger-ui/index.html | head
go version go1.26.7 linux/arm64
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect
65:        url: specUrl,

→ Read docs/swagger-ui/index.html [limit=20, offset=55]

$ ls "$(go env GOMODCACHE)/gopkg.in/yaml.v3@"* 2>/dev/null; ls "$(go env GOMODCACHE)/cache/download/gopkg.in/yaml.v3/@v/" 2>/dev/null | head
LICENSE
NOTICE
README.md
apic.go
decode.go
decode_test.go
emitterc.go
encode.go
encode_test.go
example_embedded_test.go
go.mod
limit_test.go
node_test.go
parserc.go
readerc.go
resolve.go
scannerc.go
sorter.go
suite_test.go
writerc.go
yaml.go
yamlh.go
yamlprivateh.go
list
v3.0.0-20200313102051-9f266ea9e77c.mod
v3.0.1.info
v3.0.1.lock
v3.0.1.mod
v3.0.1.zip
v3.0.1.ziphash

← Write /tmp/opencode/validate/main.go
Wrote file successfully.

$ cd /tmp/opencode/validate && go mod init validate >/dev/null 2>&1; go mod edit -require=gopkg.in/yaml.v3@v3.0.1; GOPROXY=off GOFLAGS=-mod=mod go run .
FAIL: root key: key "root" not found
exit status 1

← Edit /tmp/opencode/validate/main.go
Index: /tmp/opencode/validate/main.go
===================================================================
--- /tmp/opencode/validate/main.go
+++ /tmp/opencode/validate/main.go
@@ -88,12 +88,8 @@
 doc := &yaml.Node{}
 if err := yaml.Unmarshal(data, doc); err != nil {
 	check(err, "yaml parse")
 }
-root, err := mapGet(doc, "root")
-if err != nil {
-	check(err, "root key")
-}
 openapi := doc.Content[0]
 if openapi.Kind != yaml.MappingNode {
 	fmt.Println("FAIL: top-level node is not a mapping; parse may have produced multiple docs")
 	os.Exit(1)

$ cd /tmp/opencode/validate && GOPROXY=off GOFLAGS=-mod=mod go run .
# validate
./main.go:99:20: undefined: root

← Edit /tmp/opencode/validate/main.go
Index: /tmp/opencode/validate/main.go
===================================================================
--- /tmp/opencode/validate/main.go
+++ /tmp/opencode/validate/main.go
@@ -95,9 +95,9 @@
 	os.Exit(1)
 }
 
 missing := []string{}
-walkRefs(openapi, root, &missing)
+walkRefs(openapi, openapi, &missing)
 if len(missing) > 0 {
 	fmt.Println("FAIL: unresolvable $refs:", missing)
 	os.Exit(1)
 }

$ cd /tmp/opencode/validate && GOPROXY=off GOFLAGS=-mod=mod go run .
OK: all 126 $refs resolve
OK: PUT 409 -> #/components/responses/Conflict
OK: Conflict component has nameMismatch and majorMismatch examples with expected detail
ALL CHECKS PASSED

The YAML validates. Now let me attempt the documented gate:
$ make -C docs build 2>&1 | head -20
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

$ git status --short; git check-ignore docs/site 2>/dev/null; echo "ignore-exit=$?"
M docs/api/astarte_realm_management_api.yaml
ignore-exit=1

$ python3 -m ensurepip --version 2>&1 | head -1
pip 25.1.1

$ rg -n "theme|material|highlight" docs/mkdocs.yml | head; wc -l docs/mkdocs.yml
8:theme:
9:  name: material
15:        icon: material/brightness-7
21:        icon: material/brightness-4
38:  - pymdownx.highlight:
66 docs/mkdocs.yml

$ python3 -m venv /tmp/opencode/mkdocs-venv && /tmp/opencode/mkdocs-venv/bin/pip install -q -r /root/astrate-mule/docs/requirements.txt 2>&1 | tail -5; echo "pip-exit=$?"
pip-exit=0

$ PATH="/tmp/opencode/mkdocs-venv/bin:$PATH" make -C docs build 2>&1 | tail -25
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
INFO    -  Documentation built in 8.67 seconds
make: Leaving directory '/root/astrate-mule/docs'

$ git status --short; git check-ignore site-dist >/dev/null 2>&1 && echo "site-dist ignored" || echo "site-dist NOT ignored"
M docs/api/astarte_realm_management_api.yaml
site-dist ignored

Gate passed. Let me review the final diff:
$ git diff
diff --git a/docs/api/astarte_realm_management_api.yaml b/docs/api/astarte_realm_management_api.yaml
index 57d9981..ba1a8f3 100644
--- a/docs/api/astarte_realm_management_api.yaml
+++ b/docs/api/astarte_realm_management_api.yaml
@@ -228,6 +228,8 @@ paths:
           $ref: "#/components/responses/Unauthorized"
         "404":
           $ref: "#/components/responses/NotFound"
+        "409":
+          $ref: "#/components/responses/Conflict"
         "422":
           $ref: "#/components/responses/ValidationError"
         "500":
@@ -1141,14 +1143,27 @@ components:
               detail: Device not found
 
     Conflict:
-      description: Conflict — resource already exists.
+      description: Conflict — resource already exists or the submitted data disagrees with the resource.
       content:
         application/json; charset=utf-8:
           schema:
             $ref: "#/components/schemas/ErrorDetail"
-          example:
-            errors:
-              detail: Already exists
+          examples:
+            alreadyExists:
+              summary: The resource already exists.
+              value:
+                errors:
+                  detail: Already exists
+            nameMismatch:
+              summary: The interface name in the body does not match the one in the URL.
+              value:
+                errors:
+                  detail: Interface name doesn't match the one in the interface json
+            majorMismatch:
+              summary: The interface major version in the body does not match the one in the URL.
+              value:
+                errors:
+                  detail: Interface major version doesn't match the one in the interface json
 
     ValidationError:
       description: Unprocessable Entity — request body failed validation.

Done: added 409 Conflict to `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` referencing the shared Conflict component, and added nameMismatch/majorMismatch examples (plus preserved alreadyExists) with the exact detail strings the code pins.
Files: docs/api/astarte_realm_management_api.yaml
Verified: Go yaml probe (parse + all 126 $refs resolve + PUT 409->Conflict + example details match http_test.go strings) -> pass; `make -C docs build` -> pass (both run; mkdocs built site in 8.67s, no pip/mkdocs preinstalled so installed requirements into a /tmp venv)
Unsure: nothing
Follow-ups: none
