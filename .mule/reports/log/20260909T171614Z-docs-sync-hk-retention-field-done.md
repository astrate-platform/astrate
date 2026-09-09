slug: docs-sync-hk-retention-field
verdict: done
at:  d6d7c7f
ran: 2026-09-09T17:16:14Z on DietPi in 142s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

# Todos
[•] Read the YAML schema file to understand existing Realm/RealmCreate schemas
[ ] Read the Go code to understand the field on the wire
[ ] Add datastream_maximum_storage_retention to Realm and RealmCreate schemas
[ ] Run make -C docs build to verify
[ ] Report

→ Read docs/api/astarte_housekeeping_api.yaml [offset=180, limit=80]
→ Read internal/housekeeping/http.go [offset=40, limit=80]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=100, offset=259]

# Todos
[✓] Read the YAML schema file to understand existing Realm/RealmCreate schemas
[✓] Read the Go code to understand the field on the wire
[•] Add datastream_maximum_storage_retention to Realm and RealmCreate schemas
[ ] Run make -C docs build to verify
[ ] Report

Now I'll add the field to both schemas. The `RealmPatch` already has it, so I match its style.
← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -260,8 +260,13 @@
       type: integer
       format: int32
       nullable: true
       description: Maximum number of devices allowed in this realm. Null means unlimited.
+    datastream_maximum_storage_retention:
+      type: integer
+      format: int64
+      nullable: true
+      description: Maximum datastream storage retention in seconds. Null means unlimited.
 
 RealmCreate:
   type: object
   required: [realm_name, jwt_public_key_pem]

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -281,8 +281,13 @@
       type: integer
       format: int32
       nullable: true
       description: Maximum number of devices allowed in this realm. Null means unlimited.
+    datastream_maximum_storage_retention:
+      type: integer
+      format: int64
+      nullable: true
+      description: Maximum datastream storage retention in seconds. Null means unlimited.
 
 RealmPatch:
   type: object
   properties:

# Todos
[✓] Read the YAML schema file to understand existing Realm/RealmCreate schemas
[✓] Read the Go code to understand the field on the wire
[✓] Add datastream_maximum_storage_retention to Realm and RealmCreate schemas
[•] Run make -C docs build to verify
[ ] Report

→ Read docs/api/astarte_housekeeping_api.yaml [limit=50, offset=249]

$ make -C docs build 2>&1
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

$ python3 -c "import yaml; yaml.safe_load(open('docs/api/astarte_housekeeping_api.yaml')); print('YAML valid')" 2>&1
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; yaml.safe_load(open('docs/api/astarte_housekeeping_api.yaml')); print('YAML valid')
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ which swagger-cli redocly 2>&1; pip3 list 2>/dev/null | grep -i -E 'swagger|openapi|yaml' 2>&1
(no output)

$ pip3 install pyyaml -q 2>&1 && python3 -c "import yaml; data = yaml.safe_load(open('docs/api/astarte_housekeeping_api.yaml')); realm = data['components']['schemas']['Realm']; realm_create = data['components']['schemas']['RealmCreate']; print('Realm props:', list(realm['properties'].keys())); print('RealmCreate props:', list(realm_create['properties'].keys())); print('YAML valid, schemas load OK')"
/bin/sh: 1: pip3: not found

$ which go && go version
/usr/bin/go
go version go1.26.7 linux/arm64

$ cat > /tmp/yamlcheck.go << 'EOF'
package main

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil { fmt.Fprintf(os.Stderr, "read: %v\n", err); os.Exit(1) }
	var doc map[string]interface{}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		fmt.Fprintf(os.Stderr, "parse: %v\n", err); os.Exit(1)
	}
	schemas := doc["components"].(map[string]interface{})["schemas"].(map[string]interface{})
	for _, name := range []string{"Realm", "RealmCreate"} {
		s := schemas[name].(map[string]interface{})
		props := s["properties"].(map[string]interface{})
		fmt.Printf("%s properties: %v\n", name, keys(props))
		if _, ok := props["datastream_maximum_storage_retention"]; !ok {
			fmt.Fprintf(os.Stderr, "MISSING datastream_maximum_storage_retention in %s\n", name); os.Exit(1)
		}
	}
	fmt.Println("YAML parse OK, field present in both schemas")
}

func keys(m map[string]interface{}) []string {
	var ks []string
	for k := range m { ks = append(ks, k) }
	return ks
}
EOF
go run /tmp/yamlcheck.go docs/api/astarte_housekeeping_api.yaml 2>&1
Realm properties: [realm_name jwt_public_key_pem device_registration_limit datastream_maximum_storage_retention]
RealmCreate properties: [realm_name jwt_public_key_pem device_registration_limit datastream_maximum_storage_retention]
YAML parse OK, field present in both schemas

# Todos
[✓] Read the YAML schema file to understand existing Realm/RealmCreate schemas
[✓] Read the Go code to understand the field on the wire
[✓] Add datastream_maximum_storage_retention to Realm and RealmCreate schemas
[✓] Run make -C docs build to verify
[•] Report

Done: added `datastream_maximum_storage_retention` to `Realm` and `RealmCreate` schemas in the housekeeping API YAML
Files: docs/api/astarte_housekeeping_api.yaml
Verified: Go YAML parser validated both schemas parse cleanly, field present in `Realm` (line 264) and `RealmCreate` (line 285). `make -C docs build` failed (mkdocs not installed); full Swagger build needs the Pi's toolchain.
Unsure: nothing
Follow-ups: none
