slug: docs-sync-hk-patch-endpoint
verdict: done
at:  9fc9f33
ran: 2026-09-09T17:06:37Z on DietPi in 403s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/housekeeping/http.go
→ Read docs/api/astarte_housekeeping_api.yaml
→ Read docs/Makefile
Now let me check the docs build setup and the service's `UpdateRealm` to see which errors it can actually return.
✱ Grep "func.*UpdateRealm|ErrValidation|ErrDeletionDisabled|ErrConnectedDevicesPresent" in internal/housekeeping · 35 matches
→ Read internal/housekeeping/service.go [offset=190, limit=70]
→ Read docs/mkdocs.yml

$ ls /root/astrate-mule/docs/ && ls /root/astrate-mule/docs/api/
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

✗ Read docs/swagger.md failed
Error: File not found: /root/astrate-mule/docs/swagger.md
→ Read docs/embed.go
✱ Glob "docs/site/swagger.md" 1 match
→ Read docs/site/swagger.md
→ Read docs/requirements.txt

$ ls /root/astrate-mule/docs/swagger-ui/
index.html

→ Read docs/swagger-ui/index.html
✱ Grep "func WriteError|func WriteFieldErrors|func WriteNotFound|func WriteBadRequest|func WriteInternalServerError" in pkg/astarteapi · 5 matches
→ Read pkg/astarteapi/envelope.go [offset=140, limit=90]
✱ Grep "invalid_update_parameters|connected_devices_present|422|UnprocessableEntity|MethodNotAllowed|405" in docs/api · 24 matches
→ Read docs/api/astarte_pairing_api.yaml [limit=20, offset=110]
→ Read docs/api/astarte_appengine_api.yaml [limit=40, offset=1665]
✱ Grep "FieldErrorsDetail|FieldErrors\b" in docs/api/astarte_appengine_api.yaml · 2 matches
→ Read docs/api/astarte_appengine_api.yaml [limit=14, offset=360]
✱ Grep "ValidationErrors" in docs/api/astarte_appengine_api.yaml · 9 matches
→ Read docs/api/astarte_appengine_api.yaml [limit=20, offset=1655]
→ Read docs/api/astarte_appengine_api.yaml [limit=25, offset=1590]
Let me check what `datastream_maximum_storage_retention` means (units) so I don't invent wording.
✱ Grep "datastream_maximum_storage_retention|DatastreamMaximumStorageRetention|maximum_storage_retention" in . · 99 matches
I have everything I need. Now editing the YAML — first the PATCH operation between GET and DELETE.

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -136,8 +136,57 @@
       $ref: "#/components/responses/NotFound"
     "500":
       $ref: "#/components/responses/InternalServerError"
 
+patch:
+  operationId: patchRealm
+  summary: Update a realm
+  description: |
+    Updates one or more fields of an existing realm. Anything not named here
+    is rejected with `invalid_update_parameters`. A null
+    `device_registration_limit` or `datastream_maximum_storage_retention`
+    clears that field. Returns 200 with an empty body on success.
+  tags: [Housekeeping]
+  security:
+    - a_ha: []
+  parameters:
+    - $ref: "#/components/parameters/RealmName"
+  requestBody:
+    required: true
+    content:
+      application/json; charset=utf-8:
+        schema:
+          type: object
+          required: [data]
+          properties:
+            data:
+              $ref: "#/components/schemas/RealmPatch"
+        example:
+          data:
+            jwt_public_key_pem: |
+              -----BEGIN PUBLIC KEY-----
+              MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
+              -----END PUBLIC KEY-----
+            device_registration_limit: 100
+            datastream_maximum_storage_retention: 3600
+  responses:
+    "200":
+      description: Realm updated successfully. No body returned.
+    "400":
+      $ref: "#/components/responses/BadRequest"
+    "401":
+      $ref: "#/components/responses/Unauthorized"
+    "403":
+      $ref: "#/components/responses/Forbidden"
+    "404":
+      $ref: "#/components/responses/NotFound"
+    "405":
+      $ref: "#/components/responses/MethodNotAllowed"
+    "422":
+      $ref: "#/components/responses/ValidationErrors"
+    "500":
+      $ref: "#/components/responses/InternalServerError"
+
 delete:
   operationId: deleteRealm
   summary: Delete a realm
   description: Deletes a realm and all its data. Returns 204 on success.

Now add the `RealmPatch` schema after `RealmCreate`:
✗ Edit docs/api/astarte_housekeeping_api.yaml failed
Error: Could not find oldString in the file. It must match exactly, including whitespace, indentation, and line endings.
→ Read docs/api/astarte_housekeeping_api.yaml [limit=70, offset=228]

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -277,8 +277,26 @@
       format: int32
       nullable: true
       description: Maximum number of devices allowed in this realm. Null means unlimited.
 
+RealmPatch:
+  type: object
+  properties:
+    jwt_public_key_pem:
+      type: string
+      nullable: true
+      description: PEM-encoded public key used to verify realm JWTs. Omitted leaves the key unchanged.
+    device_registration_limit:
+      type: integer
+      format: int32
+      nullable: true
+      description: Maximum number of devices allowed in this realm. Null clears the limit (unlimited).
+    datastream_maximum_storage_retention:
+      type: integer
+      format: int64
+      nullable: true
+      description: Maximum datastream storage retention in seconds. Null clears the retention.
+
 # ── Error envelopes ──
 
 ErrorDetail:
   type: object

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -311,8 +311,20 @@
             description: Human-readable error detail string.
         example:
           detail: Bad Request
 
+  FieldErrorsDetail:
+    type: object
+    required: [errors]
+    properties:
+      errors:
+        type: object
+        additionalProperties:
+          type: array
+          items:
+            type: string
+        description: Per-field error messages (changeset shape).
+
 responses:
   BadRequest:
     description: Bad Request — malformed or missing data envelope.
     content:

Now add the `FieldErrorsDetail` schema and the two new response components:
← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -364,8 +364,29 @@
       example:
         errors:
           detail: Not Found
 
+MethodNotAllowed:
+  description: Method Not Allowed — realm deletion is disabled for this deployment.
+  content:
+    application/json; charset=utf-8:
+      schema:
+        $ref: "#/components/schemas/ErrorDetail"
+      example:
+        errors:
+          detail: Realm deletion disabled
+
+ValidationErrors:
+  description: Unprocessable Entity — update parameters failed validation.
+  content:
+    application/json; charset=utf-8:
+      schema:
+        $ref: "#/components/schemas/FieldErrorsDetail"
+      example:
+        errors:
+          error_name:
+            - invalid_update_parameters
+
 Conflict:
   description: Conflict — resource already exists (e.g. realm name taken).
   content:
     application/json; charset=utf-8:

Now verify. First check the build tooling is available, then run the doc build.
$ python3 -c "import yaml; print('pyyaml ok')" 2>&1; python3 -m mkdocs --version 2>&1 | head -1
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; print('pyyaml ok')
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'
/usr/bin/python3: No module named mkdocs

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

mkdocs isn't installed. Let me try installing the requirements so the gate can run.
$ pip3 install -q -r requirements.txt 2>&1 | tail -5; python3 -m mkdocs --version 2>&1 | head -1
/bin/sh: 1: pip3: not found
/usr/bin/python3: No module named mkdocs

$ python3 -m pip --version 2>&1; which python3 pip pip3 2>&1; ls /usr/lib/python3/dist-packages 2>/dev/null | head; apt list --installed 2>/dev/null | grep -i mkdocs
/usr/bin/python3: No module named pip
/usr/bin/python3
Cryptodome
packaging
packaging-25.0.dist-info
pycryptodomex-3.20.0.egg-info

$ python3 -m ensurepip --user 2>&1 | tail -3
module, so use it on your own risk, or make sure to only use it in virtual
environments.

$ ~/.local/bin/python3 -m pip --version 2>&1; ~/.local/bin/python3 -m pip install -q --user -r /root/astrate-mule/docs/requirements.txt 2>&1 | tail -5
/bin/sh: 1: /root/.local/bin/python3: not found
/bin/sh: 1: /root/.local/bin/python3: not found

$ python3 -m pip --version 2>&1; python3 -m pip install -q --user -r /root/astrate-mule/docs/requirements.txt 2>&1 | tail -6
/usr/bin/python3: No module named pip
/usr/bin/python3: No module named pip

$ ls ~/.local/bin 2>&1 | head -20; ls /usr/lib/python3/dist-packages 2>/dev/null | grep -i pip; /usr/local/bin/python3 -m pip --version 2>&1 | head -1
claude
/bin/sh: 1: /usr/local/bin/python3: not found

$ go version; ls $(go env GOMODCACHE)/gopkg.in 2>/dev/null | grep -i yaml
go version go1.26.7 linux/arm64
yaml.v2@v2.2.8
yaml.v3@v3.0.1

No pip/mkdocs and no PyPI access here, so I'll validate the YAML loads with Go (yaml.v3) + OpenAPI-ref resolution instead.
← Write /tmp/opencode/yamlcheck/main.go
Wrote file successfully.

← Write /tmp/opencode/yamlcheck/go.mod
Wrote file successfully.

$ GOFLAGS=-mod=mod go run . /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
YAML parses OK
42 external-component refs, all match a defined component key
realms/{realm} operations: [get patch delete]
RealmPatch schema present

YAML parses, all 42 `$ref`s resolve, and `{realm}` now has get/patch/delete. Let me check the working-tree state (read-only) and confirm the mkdocs gate note.
$ git status --short 2>&1 | head -30
M docs/api/astarte_housekeeping_api.yaml

Only the intended YAML is dirty. Let me do a final review of the patch operation block in context.
→ Read docs/api/astarte_housekeeping_api.yaml [limit=75, offset=136]

$ git diff --stat 2>&1
docs/api/astarte_housekeeping_api.yaml | 100 +++++++++++++++++++++++++++++++++
 1 file changed, 100 insertions(+)

Done: added the undocumented `PATCH /housekeeping/v1/realms/{realm}` route (openapi `patch` op, `RealmPatch` schema, `MethodNotAllowed` + `ValidationErrors` components) to docs/api/astarte_housekeeping_api.yaml.
Files: docs/api/astarte_housekeeping_api.yaml (+100 lines)
Verified: Go yaml.v3 load + all 42 `$ref`s resolve + path ops {get,patch,delete} -> pass; `make -C docs build` -> the sync step ran, `mkdocs build` failed (mkdocs/pip/network absent on this machine), so Swagger-load was checked structurally instead (swagger-ui/index.html:44 references exactly this YAML)
Unsure: 405 (deletion disabled) and the connected_devices_present 422 come only from the shared `writeError` (internal/housekeeping/http.go:233), not the current PATCH path — documented them because the task listed them; also noted `jwt_public_key_pem: null` is actually rejected as blank (http.go:187), despite being marked nullable per the task
Follow-ups: none
