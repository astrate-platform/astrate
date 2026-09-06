slug: docs-sync-appengine-by-alias-endpoints
verdict: done
at:  35359df
ran: 2026-09-06T20:10:02Z on DietPi in 144s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/appengine/http.go [offset=55, limit=30]
→ Read internal/appengine/http.go [offset=226, limit=150]
I'll read the relevant source files to understand the routes and existing spec style.

$ wc -l /root/astrate-mule/docs/api/astarte_appengine_api.yaml
872 /root/astrate-mule/docs/api/astarte_appengine_api.yaml

→ Read docs/api/astarte_appengine_api.yaml
Now let me read the YAML spec to understand the existing style for device-scoped routes, and also find the `listInterfacesByAlias` handler.
✱ Grep "func \(a \*API\) listInterfaces" 3 matches
→ Read internal/appengine/http.go [limit=20, offset=170]
Let me check `listInterfacesByAlias` and also the `listInterfaces` device-scoped handler to understand the response format, and also the `patchDevice` handler for the patch endpoint.
✱ Grep "func \(a \*API\) listInterfaces\b" 0 matches
→ Read internal/appengine/http.go [limit=20, offset=158]
✱ Grep "func \(a \*API\) writeInterfaces" 0 matches
✱ Grep "listInterfaces|interfaces" 3 matches
Now I have everything I need. Let me check if there's a device-scoped `listInterfaces` endpoint in the YAML (to mirror the `listInterfacesByAlias` route), and also check for the `patchDevice` route registration in the Go code.
✱ Grep "patchDevice" in internal/appengine · 6 matches

# Todos
[•] Add PATCH to existing /devices-by-alias/{alias} path
[ ] Add /devices-by-alias/{alias}/interfaces path (listInterfacesByAlias)
[ ] Add /devices-by-alias/{alias}/interfaces/{interface} path (getDataByAlias interface-scoped)
[ ] Add /devices-by-alias/{alias}/interfaces/{interface}/{path} path (GET/PUT/POST/DELETE)
[ ] Run make -C docs build to validate
[ ] Report

I now have the full picture. The YAML has:
- `GET /devices-by-alias/{alias}` (line 201) but no PATCH
- No `/devices-by-alias/{alias}/interfaces/...` paths at all

I need to add:
1. `PATCH` on existing `/devices-by-alias/{alias}`
2. New path: `GET /devices-by-alias/{alias}/interfaces`
3. New path: `GET/PUT/POST/DELETE /devices-by-alias/{alias}/interfaces/{interface}/{path}`
4. New path: `GET /devices-by-alias/{alias}/interfaces/{interface}`

Let me make the edits.
Now I'll add all the routes. First, add PATCH to the existing `/devices-by-alias/{alias}` path, then add the three new paths before the groups section.
← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -227,8 +227,301 @@
         $ref: "#/components/responses/DeviceNotFound"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
+  patch:
+    operationId: patchDeviceByAlias
+    summary: Update device attributes by alias
+    description: |
+      Partially updates a device's aliases, attributes, or credentials
+      inhibited status, looked up by alias. Returns the updated device status.
+    tags: [AppEngine - Devices]
+    security:
+      - a_aea: []
+    parameters:
+      - $ref: "#/components/parameters/RealmName"
+      - name: alias
+        in: path
+        required: true
+        description: The alias value to look up.
+        schema:
+          type: string
+    requestBody:
+      required: true
+      content:
+        application/json; charset=utf-8:
+          schema:
+            type: object
+            required: [data]
+            properties:
+              data:
+                $ref: "#/components/schemas/DevicePatch"
+          example:
+            data:
+              aliases:
+                serial_number: "ABC123"
+              attributes:
+                firmware_version: "2.1.0"
+              credentials_inhibited: false
+    responses:
+      "200":
+        description: Updated device status.
+        content:
+          application/json; charset=utf-8:
+            schema:
+              $ref: "#/components/schemas/DataEnvelopeDeviceStatus"
+      "400":
+        $ref: "#/components/responses/BadRequest"
+      "401":
+        $ref: "#/components/responses/Unauthorized"
+      "404":
+        $ref: "#/components/responses/DeviceNotFound"
+      "500":
+        $ref: "#/components/responses/InternalServerError"
+
+/appengine/v1/{realm}/devices-by-alias/{alias}/interfaces:
+  get:
+    operationId: listInterfacesByAlias
+    summary: List interfaces by alias
+    description: |
+      Returns the names of all interfaces currently exposed by a device
+      looked up by alias.
+    tags: [AppEngine - Data]
+    security:
+      - a_aea: []
+    parameters:
+      - $ref: "#/components/parameters/RealmName"
+      - name: alias
+        in: path
+        required: true
+        description: The alias value to look up.
+        schema:
+          type: string
+    responses:
+      "200":
+        description: List of interface names.
+        content:
+          application/json; charset=utf-8:
+            schema:
+              $ref: "#/components/schemas/DataEnvelopeStringArray"
+      "401":
+        $ref: "#/components/responses/Unauthorized"
+      "404":
+        $ref: "#/components/responses/DeviceNotFound"
+      "500":
+        $ref: "#/components/responses/InternalServerError"
+
+/appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
+  get:
+    operationId: getDataByAlias
+    summary: Get interface data by alias
+    description: |
+      Returns the data for a device interface, looked up by alias. Supports
+      optional query parameters for time-range filtering and downsampling on
+      datastream interfaces.
+    tags: [AppEngine - Data]
+    security:
+      - a_aea: []
+    parameters:
+      - $ref: "#/components/parameters/RealmName"
+      - name: alias
+        in: path
+        required: true
+        description: The alias value to look up.
+        schema:
+          type: string
+      - $ref: "#/components/parameters/InterfaceName"
+      - $ref: "#/components/parameters/DataSince"
+      - $ref: "#/components/parameters/DataSinceAfter"
+      - $ref: "#/components/parameters/DataTo"
+      - $ref: "#/components/parameters/DataLimit"
+      - $ref: "#/components/parameters/DataDownsample"
+      - $ref: "#/components/parameters/DataSort"
+    responses:
+      "200":
+        description: Interface data.
+        content:
+          application/json; charset=utf-8:
+            schema:
+              $ref: "#/components/schemas/DataEnvelopeObject"
+      "401":
+        $ref: "#/components/responses/Unauthorized"
+      "404":
+        $ref: "#/components/responses/DeviceNotFound"
+      "500":
+        $ref: "#/components/responses/InternalServerError"
+
+/appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
+  get:
+    operationId: getDataAtPathByAlias
+    summary: Get data at path by alias
+    description: |
+      Returns data for a specific path within a device interface, looked up
+      by alias.
+    tags: [AppEngine - Data]
+    security:
+      - a_aea: []
+    parameters:
+      - $ref: "#/components/parameters/RealmName"
+      - name: alias
+        in: path
+        required: true
+        description: The alias value to look up.
+        schema:
+          type: string
+      - $ref: "#/components/parameters/InterfaceName"
+      - name: path
+        in: path
+        required: true
+        description: The data path (e.g. `/sensors/temperature/value`).
+        schema:
+          type: string
+      - $ref: "#/components/parameters/DataSince"
+      - $ref: "#/components/parameters/DataSinceAfter"
+      - $ref: "#/components/parameters/DataTo"
+      - $ref: "#/components/parameters/DataLimit"
+      - $ref: "#/components/parameters/DataDownsample"
+      - $ref: "#/components/parameters/DataSort"
+    responses:
+      "200":
+        description: Path data.
+        content:
+          application/json; charset=utf-8:
+            schema:
+              $ref: "#/components/schemas/DataEnvelopeObject"
+      "401":
+        $ref: "#/components/responses/Unauthorized"
+      "404":
+        $ref: "#/components/responses/DeviceNotFound"
+      "500":
+        $ref: "#/components/responses/InternalServerError"
+
+  put:
+    operationId: putDataByAlias
+    summary: Set data at path by alias
+    description: |
+      Sets the value at a specific path in a device interface (server-owned
+      data), looked up by alias. Returns 200 on success.
+    tags: [AppEngine - Data]
+    security:
+      - a_aea: []
+    parameters:
+      - $ref: "#/components/parameters/RealmName"
+      - name: alias
+        in: path
+        required: true
+        description: The alias value to look up.
+        schema:
+          type: string
+      - $ref: "#/components/parameters/InterfaceName"
+      - name: path
+        in: path
+        required: true
+        description: The data path.
+        schema:
+          type: string
+    requestBody:
+      required: true
+      content:
+        application/json; charset=utf-8:
+          schema:
+            type: object
+            required: [data]
+            properties:
+              data:
+                description: The value to set.
+    responses:
+      "200":
+        description: Data set successfully.
+      "400":
+        $ref: "#/components/responses/BadRequest"
+      "401":
+        $ref: "#/components/responses/Unauthorized"
+      "404":
+        $ref: "#/components/responses/DeviceNotFound"
+      "500":
+        $ref: "#/components/responses/InternalServerError"
+
+  post:
+    operationId: publishDataByAlias
+    summary: Publish data at path by alias
+    description: |
+      Alias for PUT — publishes a value to a device interface path, looked
+      up by alias. Returns 200 on success.
+    tags: [AppEngine - Data]
+    security:
+      - a_aea: []
+    parameters:
+      - $ref: "#/components/parameters/RealmName"
+      - name: alias
+        in: path
+        required: true
+        description: The alias value to look up.
+        schema:
+          type: string
+      - $ref: "#/components/parameters/InterfaceName"
+      - name: path
+        in: path
+        required: true
+        description: The data path.
+        schema:
+          type: string
+    requestBody:
+      required: true
+      content:
+        application/json; charset=utf-8:
+          schema:
+            type: object
+            required: [data]
+            properties:
+              data:
+                description: The value to publish.
+    responses:
+      "200":
+        description: Data published successfully.
+      "400":
+        $ref: "#/components/responses/BadRequest"
+      "401":
+        $ref: "#/components/responses/Unauthorized"
+      "404":
+        $ref: "#/components/responses/DeviceNotFound"
+      "500":
+        $ref: "#/components/responses/InternalServerError"
+
+  delete:
+    operationId: deleteDataByAlias
+    summary: Delete data at path by alias
+    description: |
+      Unsets the property at a specific path (properties interface only),
+      looked up by alias. Returns 204 on success.
+    tags: [AppEngine - Data]
+    security:
+      - a_aea: []
+    parameters:
+      - $ref: "#/components/parameters/RealmName"
+      - name: alias
+        in: path
+        required: true
+        description: The alias value to look up.
+        schema:
+          type: string
+      - $ref: "#/components/parameters/InterfaceName"
+      - name: path
+        in: path
+        required: true
+        description: The data path.
+        schema:
+          type: string
+    responses:
+      "204":
+        description: Data deleted. No body returned.
+      "401":
+        $ref: "#/components/responses/Unauthorized"
+      "404":
+        $ref: "#/components/responses/DeviceNotFound"
+      "500":
+        $ref: "#/components/responses/InternalServerError"
+
 # ── Interface data ──────────────────────────────────────────────────
 
 /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
   get:

# Todos
[✓] Add PATCH to existing /devices-by-alias/{alias} path
[✓] Add /devices-by-alias/{alias}/interfaces path (listInterfacesByAlias)
[✓] Add /devices-by-alias/{alias}/interfaces/{interface} path (getDataByAlias interface-scoped)
[✓] Add /devices-by-alias/{alias}/interfaces/{interface}/{path} path (GET/PUT/POST/DELETE)
[•] Run make -C docs build to validate
[ ] Report

Now let me update the todos and run the docs build.
Now let me check how `make -C docs build` works and run it.
$ make -C docs build 2>&1 | head -50
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

→ Read docs/Makefile
`mkdocs` is not installed, so the full build can't run. But the task asks to "confirm the Swagger-referenced YAML still loads" — I should verify the YAML is valid. Let me check what tools are available for YAML validation, and check the Makefile/docs for how the spec is referenced.

$ which python3 yq python 2>/dev/null; python3 -c "import yaml; print('pyyaml ok')" 2>&1
/usr/bin/python3
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; print('pyyaml ok')
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

`mkdocs` isn't installed so the build itself can't fully run. But the task's key verification is that the YAML still loads validly. Let me install pyyaml or validate the YAML another way. Let me check if there's a Go-based YAML validator available, or install pyyaml.
$ pip install pyyaml 2>&1 | tail -3
/bin/sh: 1: pip: not found

$ python3 -m pip install pyyaml 2>&1 | tail -5
/usr/bin/python3: No module named pip

No pip available. Let me validate the YAML with an alternative approach — write a small Go program (mule rules recommend throwing a probe program in /tmp for verification).
$ cd /root/astrate-mule && go env GOPATH; ls $(go env GOPATH)/pkg/mod/gopkg.in 2>/dev/null | grep -i yaml
/root/go
yaml.v2@v2.2.8
yaml.v3@v3.0.1

Go's yaml.v3 is available in the module cache. Let me write a small probe program in /tmp to validate the YAML.
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ mkdir -p /tmp/yamlprobe && cat > /tmp/yamlprobe/main.go <<'EOF'
package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	data, err := os.ReadFile("/root/astrate-mule/docs/api/astarte_appengine_api.yaml")
	if err != nil {
		fmt.Println("read error:", err)
		os.Exit(1)
	}
	var doc any
	if err := yaml.Unmarshal(data, &doc); err != nil {
		fmt.Println("YAML PARSE ERROR:", err)
		os.Exit(1)
	}
	fmt.Println("YAML parses OK")
}
EOF
cat > /tmp/yamlprobe/go.mod <<'EOF'
module yamlprobe

go 1.21

require gopkg.in/yaml.v3 v3.0.1
EOF
cd /tmp/yamlprobe && GOFLAGS=-mod=mod go mod tidy 2>&1 | tail -3 && go run . 2>&1 failed
Error: The user rejected permission to use this specific tool call.
