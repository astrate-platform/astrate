slug: docs-sync-hk-validation-example
verdict: done
at:  d762f00
ran: 2026-09-13T19:19:10Z on DietPi in 68s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_housekeeping_api.yaml [offset=390, limit=60]
→ Read internal/housekeeping/http.go [offset=245, limit=20]
Task: fix the 422 ValidationError example in `docs/api/astarte_housekeeping_api.yaml` to match the real on-wire message. Let me read the named files first.
→ Read internal/housekeeping/service.go [offset=130, limit=15]
Confirmed: line 430 has the wrong detail. Fixing it.

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -426,9 +426,9 @@
       schema:
         $ref: "#/components/schemas/ErrorDetail"
       example:
         errors:
-          detail: "housekeeping: validation failed: realm name must not be empty"
+          detail: "realm_name can't be blank"
 
 InternalServerError:
   description: Internal Server Error — unexpected failure.
   content:

$ make -C docs build
make: Entering directory '/root/astrate-mule/docs'
cp DESIGN.md site/DESIGN.md
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: Leaving directory '/root/astrate-mule/docs'

The copy of the Swagger-referenced API YAML succeeded; `mkdocs` itself isn't installed on this machine. Let me check the Makefile for the right target and validate the YAML loads.
→ Read docs/Makefile

$ make -C docs sync 2>&1; python3 -c "import yaml; yaml.safe_load(open('docs/api/astarte_housekeeping_api.yaml')); print('YAML loads OK')"
make: Entering directory '/root/astrate-mule/docs'
cp DESIGN.md site/DESIGN.md
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
make: Leaving directory '/root/astrate-mule/docs'
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; yaml.safe_load(open('docs/api/astarte_housekeeping_api.yaml')); print('YAML loads OK')
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ grep -n "sigs.k8s.io/yaml\|ghodss\|gopkg.in/yaml" go.mod
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

No pyyaml on the Pi. Let me validate the YAML with a small Go probe using the repo's own yaml dependency.
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ mkdir -p /tmp/yamlprobe && cat > /tmp/yamlprobe/main.go <<'EOF'
package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	for _, p := range os.Args[1:] {
		b, err := os.ReadFile(p)
		if err != nil {
			fmt.Println(p, "ERROR", err)
			os.Exit(1)
		}
		var v any
		if err := yaml.Unmarshal(b, &v); err != nil {
			fmt.Println(p, "ERROR", err)
			os.Exit(1)
		}
		fmt.Println(p, "OK")
	}
}
EOF
go run /tmp/yamlprobe/main.go docs/api/astarte_housekeeping_api.yaml docs/site/api/astarte_housekeeping_api.yaml failed
Error: The user rejected permission to use this specific tool call.
