slug: docs-sync-rm-validation-example-prefix
verdict: done
at:  02e269b
ran: 2026-09-24T11:10:49Z on DietPi in 238s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_realm_management_api.yaml [offset=1130, limit=60]
✗ Read housekeeping/http.go failed [offset=240, limit=30]
Error: File not found: /root/astrate-mule/housekeeping/http.go
✱ Grep "realm: validation failed" in . · 18 matches
✱ Grep "validationDetail|stripValidation|realm: validation failed" in . · 15 matches
→ Read internal/realm/http.go [limit=20, offset=556]
Line numbers are stale (the yaml is now line 1269). Let me read the actual code and the example.
→ Read docs/api/astarte_realm_management_api.yaml [limit=40, offset=1240]
→ Read pkg/interfaceschema/parse.go [limit=40, offset=20]
→ Read internal/realm/service.go [limit=30, offset=150]
✱ Grep "ErrInvalid" in pkg/interfaceschema · 24 matches
→ Read internal/housekeeping/http.go [limit=20, offset=245]
✱ Grep "ErrInvalid, |%w" in pkg/interfaceschema/parse.go · 22 matches
✱ Grep "hk-validation-example" in .mule · 32 matches
→ Read .mule/reports/log/20260913T191910Z-docs-sync-hk-validation-example-done.md
→ Read pkg/interfaceschema/parse.go [limit=50, offset=160]
✱ Grep "interface definition" in pkg/interfaceschema · 1 match
"interface definition is invalid" doesn't exist in the parser — real messages are `invalid interface: ...`. Let me probe for a realistic real message.

← Write /tmp/opencode/dmprobe/main.go
Wrote file successfully.

$ go run /tmp/opencode/dmprobe/main.go
/tmp/opencode/dmprobe/main.go:6:2: no required module provides package github.com/astarte-platform/astarte-internal-dev/pkg/interfaceschema; to add it:
	go get github.com/astarte-platform/astarte-internal-dev/pkg/interfaceschema

$ head -1 go.mod
module github.com/astrate-platform/astrate

← Edit /tmp/opencode/dmprobe/main.go
Index: /tmp/opencode/dmprobe/main.go
===================================================================
--- /tmp/opencode/dmprobe/main.go
+++ /tmp/opencode/dmprobe/main.go
@@ -2,9 +2,9 @@
 
 import (
 	"fmt"
 
-	"github.com/astarte-platform/astarte-internal-dev/pkg/interfaceschema"
+	"github.com/astrate-platform/astrate/pkg/interfaceschema"
 )
 
 func main() {
 	docs := []string{

$ go run /tmp/opencode/dmprobe/main.go
"{}" => invalid interface: missing "interface_name"
"{\"interface_name\": \"com.example.Test\", \"major_version\": 0,\n   \"type\": \"properties\", \"ownership\": \"device\", \"mappings\": []}" => invalid interface: json: unknown field "major_version"
"not json" => invalid interface: invalid character 'o' in literal null (expecting 'u')
"{\"aggregation\": \"individual\"}" => invalid interface: missing "interface_name"

✱ Grep "interface_name|major_version|valid interface|InterfaceValid|\{"interface_name"" in pkg/interfaceschema/parse_test.go · 100 matches
✱ Grep "ValidationError|#\/components\/responses\/" in docs/api/astarte_realm_management_api.yaml · 100 matches

$ python3 -c "import yaml" 2>&1 || true; ls docs/Makefile && grep -n "" docs/Makefile
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml
ModuleNotFoundError: No module named 'yaml'
docs/Makefile
1:.PHONY: serve build clean install sync
2:
3:install:
4:	pip install -r requirements.txt
5:
6:sync:
7:	cp DESIGN.md site/DESIGN.md
8:	cp ROADMAP.md site/ROADMAP.md
9:	cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
10:	rm -rf site/api && cp -r api site/api
11:	rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
12:
13:serve: sync
14:	mkdocs serve -f mkdocs.yml
15:
16:build: sync
17:	mkdocs build -f mkdocs.yml
18:
19:clean:
20:	rm -rf ../site-dist

`make build` needs mkdocs (not installed — same as the prior task); the YAML copy step still runs. pyyaml is also absent, so I'll validate the YAML with a Go probe using the repo's `gopkg.in/yaml.v3`. Now the edit:
← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -1265,9 +1265,9 @@
       schema:
         $ref: "#/components/schemas/ErrorDetail"
       example:
         errors:
-          detail: "realm: validation failed: interface definition is invalid"
+          detail: "invalid interface: missing \"interface_name\""
 
 InternalServerError:
   description: Internal Server Error — unexpected failure.
   content:

$ make -C docs build 2>&1; echo "EXIT=$?"
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
EXIT=2

← Write /tmp/opencode/yamlprobe/main.go
Wrote file successfully.

mkdocs missing (expected, same as prior hk task) — the Swagger-referenced `site/api` copy succeeded. Now validating the YAML loads with a probe:
$ go run /tmp/opencode/yamlprobe/main.go docs/api/astarte_realm_management_api.yaml docs/site/api/astarte_realm_management_api.yaml
docs/api/astarte_realm_management_api.yaml OK
docs/site/api/astarte_realm_management_api.yaml OK

$ git diff --stat docs/api/astarte_realm_management_api.yaml; rg -n 'realm: validation failed' docs/api/astarte_realm_management_api.yaml || echo "no remaining unstripped prefix in rm yaml"
docs/api/astarte_realm_management_api.yaml | 2 +-
 1 file changed, 1 insertion(+), 1 deletion(-)
no remaining unstripped prefix in rm yaml

$ git status --short
M docs/api/astarte_realm_management_api.yaml

Done: ValidationError example in realm_management yaml now shows the real, prefix-stripped wire message (`invalid interface: missing "interface_name"`), measured from `ParseInterfaceCanonical` on a body missing `interface_name`.

```
Done: Fixed the realm-management ValidationError example to show the real stripped wire message "invalid interface: missing \"interface_name\"" instead of the fake prefixed one.
Files: docs/api/astarte_realm_management_api.yaml
Verified: make -C docs build -> api copy step ok, fails only at missing mkdocs (same as prior hk task); yaml.v3 probe loads docs/api/... and docs/site/api/... yaml -> pass
Unsure: chose "invalid interface: missing \"interface_name\"" (probe-measured) as the representative real message; any concrete parser message would do
Follow-ups: none
```
