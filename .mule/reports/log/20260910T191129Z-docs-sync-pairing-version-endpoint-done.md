slug: docs-sync-pairing-version-endpoint
verdict: done
at:  c867b2c
ran: 2026-09-10T19:11:29Z on DietPi in 95s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_pairing_api.yaml
→ Read internal/observability/compat.go
→ Read cmd/astrate/main.go [offset=385, limit=25]
Now let me check the realm management YAML to see how they documented the version endpoint:

$ rg -n 'version' /root/astrate-mule/docs/api/astarte_realm_management_api.yaml | head -30
10:  version: "0.0.1"
73:                        version_major: 1
74:                        version_minor: 0
139:      summary: List major versions of an interface
140:      description: Returns the installed major version numbers for the named interface.
149:          description: A list of major version numbers.
169:      description: Returns the full interface definition for a specific major version.
200:        Updates an existing interface definition (same major version). Returns
238:        Deletes an interface at a specific major version. Fails if the major
239:        version is not 0 or if the interface is still used by device introspection.
256:          description: Unprocessable Entity — major version not 0 or interface in use.
263:                  summary: Major version is not 0
266:                      detail: "Interface major version is not 0, can't be deleted"
532:  /realmmanagement/v1/{realm}/version:
535:      summary: Get emulated API version
537:        Returns the emulated upstream API compatibility version. The Astarte
546:          description: API compatibility version.
742:      description: The interface major version number (integer >= 0).
786:      required: [interface_name, version_major, version_minor, type, ownership, mappings]
793:        version_major:
796:        version_minor:

→ Read docs/api/astarte_realm_management_api.yaml [offset=532, limit=35]
Now I have all the context. Let me add the version endpoint to the pairing YAML.

← Edit docs/api/astarte_pairing_api.yaml
Index: /root/astrate-mule/docs/api/astarte_pairing_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_pairing_api.yaml
+++ /root/astrate-mule/docs/api/astarte_pairing_api.yaml
@@ -66,8 +66,44 @@
             example:
               errors:
                 detail: unhealthy
 
+/pairing/v1/{realm}/version:
+  get:
+    operationId: getVersion
+    summary: Get emulated API version
+    description: |
+      Returns the emulated upstream API compatibility version. Unauthenticated
+      — mirrors upstream's public pairing version endpoint.
+
+      The Astarte Dashboard gates feature UI on this value.
+    tags: [Pairing]
+    parameters:
+      - $ref: "#/components/parameters/RealmName"
+    responses:
+      "200":
+        description: API compatibility version.
+        content:
+          application/json; charset=utf-8:
+            schema:
+              type: object
+              required: [data]
+              properties:
+                data:
+                  type: string
+                  description: Version string.
+            example:
+              data: "1.1.0"
+      "404":
+        description: Unknown realm.
+        content:
+          application/json; charset=utf-8:
+            schema:
+              $ref: "#/components/schemas/ErrorDetail"
+            example:
+              errors:
+                detail: Not Found
+
 /pairing/v1/{realm}/agent/devices:
   post:
     operationId: registerDevice
     summary: Register a new device (flow A)

Now let me check the docs Makefile and run the build:
→ Read docs
→ Read docs/Makefile

$ which mkdocs 2>/dev/null || pip install -r /root/astrate-mule/docs/requirements.txt 2>&1 | tail -5
/bin/sh: 1: pip: not found

$ pip3 install -r /root/astrate-mule/docs/requirements.txt 2>&1 | tail -10
/bin/sh: 1: pip3: not found

$ python3 -c "import yaml; yaml.safe_load(open('/root/astrate-mule/docs/api/astarte_pairing_api.yaml')); print('YAML valid')" 2>&1
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; yaml.safe_load(open('/root/astrate-mule/docs/api/astarte_pairing_api.yaml')); print('YAML valid')
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ python3 -c "import json; print('python3 ok')" 2>&1
python3 ok

! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cat > /tmp/validate_yaml.go << 'EOF'
package main

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		os.Exit(1)
	}
	var out interface{}
	if err := yaml.Unmarshal(data, &out); err != nil {
		fmt.Fprintf(os.Stderr, "YAML parse error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("YAML valid")
}
EOF
cd /tmp && go run validate_yaml.go /root/astrate-mule/docs/api/astarte_pairing_api.yaml failed
Error: The user rejected permission to use this specific tool call.
