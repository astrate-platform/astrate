slug: docs-sync-ae-forbidden-403
verdict: done
at:  6a29ac7
ran: 2026-09-26T20:34:15Z on DietPi in 1037s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/swagger/swagger_test.go

$ ls /root/astrate-mule/docs/ /root/astrate-mule/docs/api/ /root/astrate-mule/internal/swagger/
/root/astrate-mule/docs/:
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

/root/astrate-mule/docs/api/:
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

/root/astrate-mule/internal/swagger/:
swagger.go
swagger_test.go

$ sed -n '1255,1300p' /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
responses:
    BadRequest:
      description: Bad Request — malformed or missing data envelope.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Bad Request

    Unauthorized:
      description: Unauthorized — missing or invalid JWT.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Unauthorized

    Forbidden:
      description: Forbidden — valid JWT but insufficient permissions.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Forbidden

    NotFound:
      description: Not Found — resource does not exist.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Not Found

    DeviceNotFound:
      description: Device not found — unknown device ID.
      content:
        application/json; charset=utf-8:

$ grep -n '"401"\|"403"\|Forbidden\|Unauthorized' /root/astrate-mule/docs/api/astarte_realm_management_api.yaml | head -80
88:        "401":
89:          $ref: "#/components/responses/Unauthorized"
90:        "403":
91:          $ref: "#/components/responses/Forbidden"
132:        "401":
133:          $ref: "#/components/responses/Unauthorized"
134:        "403":
135:          $ref: "#/components/responses/Forbidden"
194:        "401":
195:          $ref: "#/components/responses/Unauthorized"
196:        "403":
197:          $ref: "#/components/responses/Forbidden"
227:        "401":
228:          $ref: "#/components/responses/Unauthorized"
229:        "403":
230:          $ref: "#/components/responses/Forbidden"
266:        "401":
267:          $ref: "#/components/responses/Unauthorized"
268:        "403":
269:          $ref: "#/components/responses/Forbidden"
326:        "401":
327:          $ref: "#/components/responses/Unauthorized"
330:        "403":
331:          description: Forbidden — major version not 0 or interface in use.
372:        "401":
373:          $ref: "#/components/responses/Unauthorized"
374:        "403":
375:          $ref: "#/components/responses/Forbidden"
415:        "401":
416:          $ref: "#/components/responses/Unauthorized"
417:        "403":
418:          $ref: "#/components/responses/Forbidden"
472:        "401":
473:          $ref: "#/components/responses/Unauthorized"
474:        "403":
475:          $ref: "#/components/responses/Forbidden"
494:        "401":
495:          $ref: "#/components/responses/Unauthorized"
496:        "403":
497:          $ref: "#/components/responses/Forbidden"
532:        "401":
533:          $ref: "#/components/responses/Unauthorized"
534:        "403":
535:          $ref: "#/components/responses/Forbidden"
569:        "401":
570:          $ref: "#/components/responses/Unauthorized"
571:        "403":
572:          $ref: "#/components/responses/Forbidden"
604:        "401":
605:          $ref: "#/components/responses/Unauthorized"
606:        "403":
607:          $ref: "#/components/responses/Forbidden"
639:        "401":
640:          $ref: "#/components/responses/Unauthorized"
641:        "403":
642:          $ref: "#/components/responses/Forbidden"
674:        "401":
675:          $ref: "#/components/responses/Unauthorized"
676:        "403":
677:          $ref: "#/components/responses/Forbidden"
701:        "401":
702:          $ref: "#/components/responses/Unauthorized"
703:        "403":
704:          $ref: "#/components/responses/Forbidden"
732:        "401":
733:          $ref: "#/components/responses/Unauthorized"
734:        "403":
735:          $ref: "#/components/responses/Forbidden"
775:        "401":
776:          $ref: "#/components/responses/Unauthorized"
777:        "403":
778:          $ref: "#/components/responses/Forbidden"
809:        "401":
810:          $ref: "#/components/responses/Unauthorized"
811:        "403":
812:          $ref: "#/components/responses/Forbidden"
832:        "401":
833:          $ref: "#/components/responses/Unauthorized"
834:        "403":
835:          $ref: "#/components/responses/Forbidden"

$ wc -l /root/astrate-mule/docs/api/astarte_appengine_api.yaml; grep -n '^paths:\|^components:\|^    responses:\|^    schemas:\|^    parameters:\|^  [a-z]' /root/astrate-mule/docs/api/astarte_appengine_api.yaml
1766 /root/astrate-mule/docs/api/astarte_appengine_api.yaml
3:  title: Astarte AppEngine API
4:  description: |
9:  version: "0.0.1"
10:  license:
23:paths:
1317:components:
1318:  securitySchemes:
1327:  parameters:
1441:  schemas:
1646:  responses:

$ grep -n '"401"' /root/astrate-mule/docs/api/astarte_appengine_api.yaml | wc -l; grep -c '^    \(get\|post\|put\|delete\|patch\):$' /root/astrate-mule/docs/api/astarte_appengine_api.yaml
32
32

$ sed -n '23,140p' /root/astrate-mule/docs/api/astarte_appengine_api.yaml
paths:
  # ── Devices ─────────────────────────────────────────────────────────

  /appengine/v1/{realm}/devices:
    get:
      operationId: listDevices
      summary: List devices in a realm
      description: |
        Returns a paginated list of device IDs (or device status objects when
        `details=true`). Supports cursor-based pagination via `links.next`.
      tags: [AppEngine - Devices]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - name: details
          in: query
          schema:
            type: boolean
            default: false
          description: When true, return full device status objects instead of just IDs.
        - name: from_token
          in: query
          schema:
            type: string
          description: Pagination cursor — start after this device token.
        - name: limit
          in: query
          schema:
            type: integer
          description: Maximum number of devices to return.
      responses:
        "200":
          description: Paginated list of devices.
          content:
            application/json; charset=utf-8:
              schema:
                oneOf:
                  - $ref: "#/components/schemas/DataEnvelopeDeviceIDList"
                  - $ref: "#/components/schemas/DataEnvelopeDeviceStatusList"
              examples:
                ids:
                  summary: Simple device ID list
                  value:
                    data:
                      - dT6hS2W9TT6LEnP25ks_lg
                    links:
                      self: "/v1/test/devices?details=false&limit=1"
                details:
                  summary: Device status list with pagination
                  value:
                    data:
                      - id: dT6hS2W9TT6LEnP25ks_lg
                        aliases: {}
                        attributes: {}
                        connected: true
                        last_connection: "2024-05-30 13:49:57.045Z"
                        first_registration: "2024-01-15 10:00:00.000Z"
                        last_disconnection: null
                        credentials_inhibited: false
                        groups: []
                        total_received_msgs: 12
                        total_received_bytes: 4096
                        introspection:
                          org.example.Sensors:
                            major: 0
                            minor: 1
                    links:
                      self: "/v1/test/devices?details=true&limit=1"
                      next: "/v1/test/devices?details=true&from_token=dT6hS2W9TT6LEnP25ks_lg&limit=1"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "500":
          $ref: "#/components/responses/InternalServerError"

  /appengine/v1/{realm}/stats/devices:
    get:
      operationId: getDeviceStats
      summary: Get device statistics
      description: Returns total and connected device counts for the realm.
      tags: [AppEngine - Devices]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
      responses:
        "200":
          description: Device statistics.
          content:
            application/json; charset=utf-8:
              schema:
                type: object
                required: [data]
                properties:
                  data:
                    type: object
                    required: [connected_devices, total_devices]
                    properties:
                      connected_devices:
                        type: integer
                        format: int64
                      total_devices:
                        type: integer
                        format: int64
              example:
                data:
                  connected_devices: 42
                  total_devices: 150
        "401":
          $ref: "#/components/responses/Unauthorized"
        "500":
          $ref: "#/components/responses/InternalServerError"

  /appengine/v1/{realm}/devices/{device}:
    get:
      operationId: getDevice
      summary: Get device details
      description: Returns the full device status object.

$ sed -n '1317,1340p;1640,1700p' /root/astrate-mule/docs/api/astarte_appengine_api.yaml
components:
  securitySchemes:
    a_aea:
      type: apiKey
      in: header
      name: Authorization
      description: |
        Realm-level JWT carrying the `a_aea` claim. Used to authenticate
        AppEngine operations (device data, groups).

  parameters:
    RealmName:
      name: realm
      in: path
      required: true
      description: The realm name.
      schema:
        type: string

    DeviceID:
      name: device
      in: path
      required: true
      description: The device hardware ID.
          additionalProperties:
            type: array
            items:
              type: string
          description: Per-field error messages (changeset shape).

  responses:
    BadRequest:
      description: Bad Request — malformed or missing data envelope.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Bad Request

    Unauthorized:
      description: Unauthorized — missing or invalid JWT.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Unauthorized

    NotFound:
      description: Not Found — resource does not exist.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Not Found

    DeviceNotFound:
      description: Device not found — unknown device ID.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Device not found

    MethodNotAllowed:
      description: Method Not Allowed — cannot write to a device-owned or read-only resource.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Cannot write to read-only resource

    InternalServerError:
      description: Internal Server Error — unexpected failure.
      content:
        application/json; charset=utf-8:

$ grep -n '^        "[0-9][0-9][0-9]":' /root/astrate-mule/docs/api/astarte_appengine_api.yaml | sed 's/.*\(  *"[0-9]*"\).*/\1/' | sort | uniq -c
26  "200"
      2  "201"
      4  "204"
     11  "400"
     32  "401"
     28  "404"
      9  "405"
      4  "409"
     17  "422"
     32  "500"

$ grep -n -A2 '^        "401":' /root/astrate-mule/docs/api/astarte_appengine_api.yaml | grep -v '^--$' | grep -v '401' | sed 's/^[0-9]*-//' | sed 's/^ *//' | sort | uniq -c
28 "404":
      1 "409":
      3 "500":
     32 $ref: "#/components/responses/Unauthorized"

$ grep -n '"401":' /root/astrate-mule/docs/api/astarte_appengine_api.yaml | head -3; python3 -c "
import io
p='docs/api/astarte_appengine_api.yaml'
s=open(p,encoding='utf-8').read()
needle='        \"401\":\n          \$ref: \"#/components/responses/Unauthorized\"\n'
print('occurrences:', s.count(needle))
print('403 present:', '\"403\"' in s)
"
93:        "401":
131:        "401":
154:        "401":
occurrences: 32
403 present: False

$ python3 - <<'EOF'
p='docs/api/astarte_appengine_api.yaml'
s=open(p,encoding='utf-8').read()

comp_anchor = '''    Unauthorized:
      description: Unauthorized \u2014 missing or invalid JWT.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Unauthorized
'''
assert s.count(comp_anchor) == 1, s.count(comp_anchor)
forbidden = '''
    Forbidden:
      description: Forbidden \u2014 valid JWT but insufficient permissions.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Forbidden
'''
s = s.replace(comp_anchor, comp_anchor + forbidden, 1)

needle = '        "401":\n          $ref: "#/components/responses/Unauthorized"\n'
add = '        "403":\n          $ref: "#/components/responses/Forbidden"\n'
assert s.count(needle) == 32
s = s.replace(needle, needle + add)
open(p,'w',encoding='utf-8').write(s)
print('ok')
EOF
ok

$ grep -c '"403":' docs/api/astarte_appengine_api.yaml; grep -n -A1 '"403":' docs/api/astarte_appengine_api.yaml | grep -c 'responses/Forbidden'; git diff --stat
32
32
 docs/api/astarte_appengine_api.yaml | 74 +++++++++++++++++++++++++++++++++++++
 1 file changed, 74 insertions(+)

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -120,11 +120,30 @@
 // 401: RequireRealm answers 403 for a verified realm JWT whose a_rma grants
 // do not authorize the method+path (internal/auth/middleware.go). Matches the
 // upstream-parity the housekeeping and pairing specs already document.
 func TestRealmManagement403(t *testing.T) {
-	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
+	specDocuments403(t, "astarte_realm_management_api.yaml")
+}
+
+// TestAppEngine403 is the AppEngine half of the same contract: every route is
+// wrapped by RequireRealm(auth.ClaimAppEngine) (internal/appengine/http.go),
+// so a verified a_aea JWT whose grants do not authorize the method+path gets
+// 403 Forbidden, not 401 (astarteapi.WriteForbidden,
+// internal/auth/middleware.go) — the behaviour pinned by
+// internal/appengine/http_test.go. Documenting only 401 would teach a
+// generated client that 403 cannot happen.
+func TestAppEngine403(t *testing.T) {
+	specDocuments403(t, "astarte_appengine_api.yaml")
+}
+
+// specDocuments403 checks that every operation in the given spec documents a
+// 403 response and that components.responses defines the Forbidden response
+// they reference.
+func specDocuments403(t *testing.T, filename string) {
+	t.Helper()
+	b, err := docs.APIYAML.ReadFile("api/" + filename)
 	if err != nil {
-		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
+		t.Fatalf("reading %s: %v", filename, err)
 	}
 	lines := strings.Split(string(b), "\n")
 
 	pathIdx, compIdx := -1, -1

→ Read internal/swagger/swagger_test.go [offset=116, limit=70]

$ grep -c '^    patch:$' docs/api/astarte_appengine_api.yaml; go test ./internal/swagger/ -run 'TestAppEngine403|TestRealmManagement403' -v 2>&1 | tail -20
3
=== RUN   TestRealmManagement403
--- PASS: TestRealmManagement403 (0.00s)
=== RUN   TestAppEngine403
--- PASS: TestAppEngine403 (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.029s

$ for f in docs/api/*.yaml; do echo "== $f"; awk 'NR>1 && /^paths:/{p=1} /^components:/{p=0} p && /^    [a-z]+:$/ {print $1}' "$f" | sort | uniq -c; done
== docs/api/astarte_appengine_api.yaml
      4 delete:
     17 get:
      3 patch:
      5 post:
      3 put:
== docs/api/astarte_housekeeping_api.yaml
      1 delete:
      2 get:
      1 patch:
      1 post:
== docs/api/astarte_pairing_api.yaml
      1 delete:
      3 get:
      3 post:
== docs/api/astarte_realm_management_api.yaml
      4 delete:
     11 get:
      3 post:
      2 put:
== docs/api/astrate_native_api.yaml
     12 get:

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -161,9 +161,9 @@
 if pathIdx < 0 || compIdx < pathIdx {
 	t.Fatal("cannot locate the paths and components sections")
 }
 
-methodRe := regexp.MustCompile(`^    (get|post|put|delete):$`)
+methodRe := regexp.MustCompile(`^    (get|post|put|delete|patch):$`)
 var ops []int
 for i, l := range lines[pathIdx:compIdx] {
 	if methodRe.MatchString(l) {
 		ops = append(ops, pathIdx+i)

! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cp docs/api/astarte_appengine_api.yaml /tmp/ae.yaml.bak && python3 - <<'EOF'
p='docs/api/astarte_appengine_api.yaml'
s=open(p,encoding='utf-8').read()
# drop the 403 ref of the first patch operation
i=s.index('    patch:')
j=s.index('        "403":\n          $ref: "#/components/responses/Forbidden"\n', i)
s=s[:j]+s[j+len('        "403":\n          $ref: "#/components/responses/Forbidden"\n'):]
open(p,'w',encoding='utf-8').write(s)
EOF
go test ./internal/swagger/ -run 'TestAppEngine403|TestRealmManagement403' 2>&1 | tail -8 failed
Error: The user rejected permission to use this specific tool call.
