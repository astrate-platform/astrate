slug: docs-sync-rm-deviceid-param
verdict: done
at:  2d1a1c1
ran: 2026-10-02T20:13:07Z on DietPi in 303s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
I'll start by reading the required files.

$ sed -n 860,920p docs/api/astrate_realm_management_api.yaml
sed: can't read docs/api/astrate_realm_management_api.yaml: No such file or directory

$ rg -n "minLength|maxLength|pattern" docs/api/*.yaml | head -50
docs/api/astarte_pairing_api.yaml:388:          minLength: 22
docs/api/astarte_pairing_api.yaml:389:          maxLength: 22
docs/api/astarte_pairing_api.yaml:390:          pattern: '^[A-Za-z0-9_-]{22}$'
docs/api/astarte_realm_management_api.yaml:981:          maxLength: 128
docs/api/astarte_realm_management_api.yaml:982:          pattern: '^([a-zA-Z][a-zA-Z0-9]*\.([a-zA-Z0-9][a-zA-Z0-9-]*\.)*)?[a-zA-Z][a-zA-Z0-9]*$'
docs/api/astarte_realm_management_api.yaml:1146:          minLength: 1
docs/api/astarte_realm_management_api.yaml:1147:          maxLength: 128
docs/api/astarte_realm_management_api.yaml:1148:          pattern: '^[a-zA-Z0-9_.~-]{1,128}$'

$ sed -n 375,395p docs/api/astarte_pairing_api.yaml
a padded 23-character spelling of the same value, are not accepted.
      schema:
        type: string

  schemas:
    # ── Request types ──

    RegisterRequest:
      type: object
      required: [hw_id]
      properties:
        hw_id:
          type: string
          minLength: 22
          maxLength: 22
          pattern: '^[A-Za-z0-9_-]{22}$'
          description: |
            Device hardware ID: a 128-bit value in the 22-character unpadded
            base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and
            `_`, with no `=` padding. The standard-alphabet characters `+` and
            `/`, and a padded 23-character spelling of the same value, are

$ rg -n "hardware ID" docs/api/astarte_realm_management_api.yaml
933:      description: The device hardware ID.

$ sed -n 360,400p docs/api/astarte_pairing_api.yaml
name: realm
      in: path
      required: true
      description: The realm name.
      schema:
        type: string

    DeviceID:
      name: deviceID
      in: path
      required: true
      description: |
        The device hardware ID: a 128-bit value in the 22-character unpadded
        base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and `_`,
        with no `=` padding. The standard-alphabet characters `+` and `/`, and
        a padded 23-character spelling of the same value, are not accepted.
      schema:
        type: string

  schemas:
    # ── Request types ──

    RegisterRequest:
      type: object
      required: [hw_id]
      properties:
        hw_id:
          type: string
          minLength: 22
          maxLength: 22
          pattern: '^[A-Za-z0-9_-]{22}$'
          description: |
            Device hardware ID: a 128-bit value in the 22-character unpadded
            base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and
            `_`, with no `=` padding. The standard-alphabet characters `+` and
            `/`, and a padded 23-character spelling of the same value, are
            rejected with 422 and the
            `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}`
            body.
        initial_payload_format:
          type: string

$ sed -n 915,950p docs/api/astarte_realm_management_api.yaml
required: true
      description: The interface major version number (integer >= 0).
      schema:
        type: integer
        minimum: 0

    TriggerName:
      name: name
      in: path
      required: true
      description: The trigger name.
      schema:
        type: string

    DeviceID:
      name: device
      in: path
      required: true
      description: The device hardware ID.
      schema:
        type: string

    PolicyName:
      name: name
      in: path
      required: true
      description: The trigger delivery policy name.
      schema:
        type: string

    AsyncOperation:
      name: async_operation
      in: query
      required: false
      description: |
        Accepted and ignored on either value. Upstream 1.4 runs interface

$ sed -n 110,128p internal/realm/service.go; echo ---; sed -n 100,112p internal/realm/http.go
// #21: back-to-back around the sync path so imported trigger configs still
// fire). finished is emitted even when the store delete fails, so a started
// lifecycle always closes.
func (s *Service) DeleteDevice(ctx context.Context, realm, deviceID string) error {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return err
	}
	id, err := deviceid.Parse(deviceID)
	if err != nil {
		return fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
	}
	if s.disc != nil {
		s.disc.DisconnectDevice(realm, id)
	}
	at := time.Now().UTC()
	if s.OnDeletionStart != nil {
		s.OnDeletionStart(realm, deviceID, at)
	}
---
}

// deleteDevice synchronously removes a device and its data (the dashboard's
// Delete button; docs/COMPATIBILITY.md notes the deviation from upstream's
// async deletion).
func (a *API) deleteDevice(w http.ResponseWriter, r *http.Request) {
	err := a.svc.DeleteDevice(r.Context(), r.PathValue("realm"), r.PathValue("device"))
	if errors.Is(err, store.ErrNotFound) {
		_ = astarteapi.WriteDeviceNotFound(w)
		return
	}
	if err != nil {
		a.writeError(w, err)

→ Read pkg/deviceid/deviceid.go [limit=70]

$ rg -n "DeviceID" docs/api/astarte_realm_management_api.yaml | head -30
707:        - $ref: "#/components/parameters/DeviceID"
929:    DeviceID:

$ rg -n "Device not found" docs/api/astarte_realm_management_api.yaml | head -20; echo ---; git log --oneline -12
1340:      description: Device not found — unknown device ID.
1347:              detail: Device not found
---
d450919 mule: log docs-sync-rm-validationerror-example
8a09a1a mule: docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
5c30118 mule: blocked docs-sync-rm-legacy-alias-fields
a52889d mule: log docs-sync-rm-error-example-capitalisation
1d1e6f4 mule: docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, housekeeping and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
22cfaee mule: refill the queue
c145591 mule: recipe docs-sync ran (530s)
94b0d28 mule: blocked container-response-cap-test
e5753dd mule: log container-parseconfig-rules-test
32c85af mule: container-parseconfig-rules-test [auto]: pin the four `parseConfig` rules that no test touches (internal/flow/blocks/container/block.go:188-219 — all `count=0` in `go test -coverprofile`, which is 71.2% for the package): `config` must be a JSON object (192-193), `port` an integer in 1–65535 including the 0 / 65536 / non-integer rejections (197-202), and `timeout_ms` and `ready_timeout_ms` positive (205-210, 215-217) — beside the existing `TestConstructor_RequiresImage` (block_test.go:99), which is the only config case the suite has. The port rule is the one with teeth: it is what keeps `127.0.0.1::99999` out of the `docker run` argv (docker.go:93). Container-free, pure.
141e131 mule: blocked container-timeout-bounds
927aa1b mule: log container-stop-deadline

$ sed -n 685,712p docs/api/astarte_realm_management_api.yaml; echo '--- 1330-1355 ---'; sed -n 1328,1355p docs/api/astarte_realm_management_api.yaml
"401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "500":
          $ref: "#/components/responses/InternalServerError"

  # ── Device lifecycle ────────────────────────────────────────────────

  /realmmanagement/v1/{realm}/devices/{device}:
    delete:
      operationId: deleteDevice
      summary: Delete a device
      description: |
        Synchronously deletes a device and all its data. This is a deliberate
        deviation from upstream's async deletion (see COMPATIBILITY.md).
        Returns 204 on success.
      tags: [Realm Management - Devices]
      security:
        - a_rma: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - $ref: "#/components/parameters/DeviceID"
        - $ref: "#/components/parameters/AsyncOperation"
      responses:
        "204":
          description: Device deleted. No body returned.
        "401":
--- 1330-1355 ---

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

    Conflict:
      description: Conflict — resource already exists or the submitted data disagrees with the resource.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          examples:

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -929,11 +929,21 @@
 DeviceID:
   name: device
   in: path
   required: true
-  description: The device hardware ID.
+  description: |
+    The device hardware ID: a 128-bit value in the 22-character unpadded
+    base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and `_`,
+    with no `=` padding. The standard-alphabet characters `+` and `/`, and
+    a padded 23-character spelling of the same value, are not accepted —
+    and a malformed ID is not reported as malformed: it is answered with
+    the same `404` and `{"errors":{"detail":"Device not found"}}` body as an
+    unknown device.
   schema:
     type: string
+    minLength: 22
+    maxLength: 22
+    pattern: '^[A-Za-z0-9_-]{22}$'
 
 PolicyName:
   name: name
   in: path

$ cat docs/Makefile 2>/dev/null || ls docs
.PHONY: serve build clean install sync

install:
	pip install -r requirements.txt

sync:
	cp DESIGN.md site/DESIGN.md
	cp ROADMAP.md site/ROADMAP.md
	cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
	rm -rf site/api && cp -r api site/api
	rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui

serve: sync
	mkdocs serve -f mkdocs.yml

build: sync
	mkdocs build -f mkdocs.yml

clean:
	rm -rf ../site-dist

$ rg -ln "astarte_realm_management_api|docs/api" --glob '*.go' . | head
./internal/swagger/swagger_test.go

→ Read internal/swagger/swagger_test.go
Now the test, mirroring the pairing twin's guard:

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -719,8 +719,79 @@
 		}
 	}
 }
 
+// TestRealmManagementDeviceIDEncodingDocumented is the realm-management twin of
+// TestPairingDeviceIDEncodingDocumented: the `DeviceID` path parameter of
+// deleteDevice must state the same wire form, and constrain itself to it. The
+// parameter used to be a bare `type: string` described only as "The device
+// hardware ID", which read as any string and let a client send a well-formed
+// standard-base64 or padded 23-character ID. Service.DeleteDevice folds every
+// deviceid.Parse failure into store.ErrNotFound (internal/realm/service.go), so
+// such a client is answered 404 Device not found rather than told its ID is
+// malformed — the description has to name that consequence, and the pattern has
+// to reject the spellings the parser rejects. The bounds and the description
+// wording are derived from deviceid.EncodedLen, and the shipped pattern is
+// compiled and held against deviceid.Parse here, so the spec cannot drift from
+// the parser. The strict trailing-bits rule is deliberately not encoded in the
+// pattern: whether recording it is worth it is still open.
+func TestRealmManagementDeviceIDEncodingDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	wantLen := strconv.Itoa(deviceid.EncodedLen)
+	wantDesc := fmt.Sprintf("%s-character unpadded", wantLen)
+
+	param := componentBlock(t, lines, "    DeviceID:")
+	for _, want := range []string{wantDesc, "base64url", astarteapi.DetailDeviceNotFound} {
+		if !containsLine(param, want) && !strings.Contains(strings.Join(param, "\n"), want) {
+			t.Errorf("DeviceID parameter does not say %q", want)
+		}
+	}
+	for _, want := range []string{
+		"        minLength: " + wantLen,
+		"        maxLength: " + wantLen,
+	} {
+		if !containsLine(param, want) {
+			t.Errorf("DeviceID parameter is missing line %q", want)
+		}
+	}
+
+	pattern := ""
+	for _, l := range param {
+		if v, ok := strings.CutPrefix(l, "        pattern: "); ok {
+			pattern = strings.Trim(v, `'"`)
+			break
+		}
+	}
+	if pattern == "" {
+		t.Fatal("DeviceID parameter carries no pattern")
+	}
+	re, err := regexp.Compile(pattern)
+	if err != nil {
+		t.Fatalf("DeviceID pattern %q does not compile: %v", pattern, err)
+	}
+
+	if !re.MatchString("dT6hS2W9TT6LEnP25ks_lg") {
+		t.Errorf("DeviceID pattern %q rejects the canonical ID dT6hS2W9TT6LEnP25ks_lg", pattern)
+	}
+	for _, spelling := range []string{
+		"dT6hS2W9TT6LEnP25ks+lg",  // standard-alphabet '+'
+		"dT6hS2W9TT6LEnP25ks/lg",  // standard-alphabet '/'
+		"dT6hS2W9TT6LEnP25ks_lg=", // padded
+	} {
+		if _, err := deviceid.Parse(spelling); err == nil {
+			t.Errorf("deviceid.Parse accepts %q; this test's premise no longer holds", spelling)
+		}
+		if re.MatchString(spelling) {
+			t.Errorf("DeviceID pattern %q accepts %q, which deviceid.Parse rejects", pattern, spelling)
+		}
+	}
+}
+
 // TestPairingUnregisterDeviceSemanticsDocumented guards that the pairing spec
 // describes `unregisterDevice` for what it does. "Removes a device from the
 // realm" reads as a device-and-data deletion; store.UnregisterDevice clears the
 // credential trail and flips status back to 'registered' and touches nothing

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -745,9 +745,9 @@
 wantDesc := fmt.Sprintf("%s-character unpadded", wantLen)
 
 param := componentBlock(t, lines, "    DeviceID:")
 for _, want := range []string{wantDesc, "base64url", astarteapi.DetailDeviceNotFound} {
-	if !containsLine(param, want) && !strings.Contains(strings.Join(param, "\n"), want) {
+	if !strings.Contains(strings.Join(param, "\n"), want) {
 		t.Errorf("DeviceID parameter does not say %q", want)
 	}
 }
 for _, want := range []string{

$ gofmt -l . && go vet ./internal/swagger/ && go test ./internal/swagger/ -run 'DeviceID' -v 2>&1 | tail -20
=== RUN   TestPairingDeviceIDEncodingDocumented
--- PASS: TestPairingDeviceIDEncodingDocumented (0.00s)
=== RUN   TestRealmManagementDeviceIDEncodingDocumented
--- PASS: TestRealmManagementDeviceIDEncodingDocumented (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.036s

$ cp docs/api/astarte_realm_management_api.yaml /tmp/opencode/rm.yaml.new && python3 - <<'EOF'
import re
p='docs/api/astarte_realm_management_api.yaml'
s=open(p).read()
new=s.replace("""    DeviceID:
      name: device
      in: path
      required: true
      description: |
        The device hardware ID: a 128-bit value in the 22-character unpadded
        base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and `_`,
        with no `=` padding. The standard-alphabet characters `+` and `/`, and
        a padded 23-character spelling of the same value, are not accepted —
        and a malformed ID is not reported as malformed: it is answered with
        the same `404` and `{"errors":{"detail":"Device not found"}}` body as an
        unknown device.
      schema:
        type: string
        minLength: 22
        maxLength: 22
        pattern: '^[A-Za-z0-9_-]{22}$'
""","""    DeviceID:
      name: device
      in: path
      required: true
      description: The device hardware ID.
      schema:
        type: string
""")
assert new!=s
open(p,'w').write(new)
EOF
go test ./internal/swagger/ -run TestRealmManagementDeviceIDEncodingDocumented 2>&1 | tail -20; cp /tmp/opencode/rm.yaml.new docs/api/astarte_realm_management_api.yaml
--- FAIL: TestRealmManagementDeviceIDEncodingDocumented (0.00s)
    swagger_test.go:750: DeviceID parameter does not say "22-character unpadded"
    swagger_test.go:750: DeviceID parameter does not say "base64url"
    swagger_test.go:750: DeviceID parameter does not say "Device not found"
    swagger_test.go:758: DeviceID parameter is missing line "        minLength: 22"
    swagger_test.go:758: DeviceID parameter is missing line "        maxLength: 22"
    swagger_test.go:770: DeviceID parameter carries no pattern
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.031s
FAIL

$ go test ./internal/swagger/ 2>&1 | tail -3; which mkdocs || echo "mkdocs missing"
ok  	github.com/astrate-platform/astrate/internal/swagger	0.113s
mkdocs missing

$ python3 -c "import yaml,sys; d=yaml.safe_load(open('docs/api/astarte_realm_management_api.yaml')); import json; print(json.dumps(d['components']['parameters']['DeviceID'], indent=2))" 2>&1 | tail -20
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml,sys; d=yaml.safe_load(open('docs/api/astarte_realm_management_api.yaml')); import json; print(json.dumps(d['components']['parameters']['DeviceID'], indent=2))
    ^^^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ timeout 60 pip install --quiet pyyaml 2>&1 | tail -3; python3 -c "import yaml; print('ok')" 2>&1 | tail -1
timeout: failed to run command ‘pip’: No such file or directory
ModuleNotFoundError: No module named 'yaml'

$ timeout 90 python3 -m pip install --quiet pyyaml 2>&1 | tail -3; python3 -c "import yaml; print('yaml ok')" 2>&1 | tail -1
/usr/bin/python3: No module named pip
ModuleNotFoundError: No module named 'yaml'

$ rg -n "yaml" go.mod
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

$ mkdir -p /tmp/opencode/yamlprobe && cat > /tmp/opencode/yamlprobe/main_test.go <<'EOF'
package yamlprobe

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoads(t *testing.T) {
	b, err := os.ReadFile("/root/astrate-mule/docs/api/astarte_realm_management_api.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		OpenAPI string `yaml:"openapi"`
		Paths   map[string]map[string]struct {
			Parameters []struct {
				Ref string `yaml:"$ref"`
			} `yaml:"parameters"`
		} `yaml:"paths"`
		Components struct {
			Parameters map[string]struct {
				Name        string `yaml:"name"`
				In          string `yaml:"in"`
				Required    bool   `yaml:"required"`
				Description string `yaml:"description"`
				Schema      struct {
					Type       string `yaml:"type"`
					MinLength  int    `yaml:"minLength"`
					MaxLength  int    `yaml:"maxLength"`
					Pattern    string `yaml:"pattern"`
				} `yaml:"schema"`
			} `yaml:"parameters"`
		} `yaml:"components"`
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		t.Fatalf("spec does not parse: %v", err)
	}
	t.Logf("openapi=%s paths=%d", doc.OpenAPI, len(doc.Paths))
	d, ok := doc.Components.Parameters["DeviceID"]
	if !ok {
		t.Fatal("no DeviceID parameter")
	}
	out, _ := json.MarshalIndent(d, "", "  ")
	t.Logf("DeviceID: %s", out)
	re := regexp.MustCompile(d.Schema.Pattern)
	for _, s := range []string{"dT6hS2W9TT6LEnP25ks_lg", "dT6hS2W9TT6LEnP25ks+lg", "dT6hS2W9TT6LEnP25ks/lg", "dT6hS2W9TT6LEnP25ks_lg="} {
		t.Logf("pattern %s on %q -> %v", d.Schema.Pattern, s, re.MatchString(s))
	}
}
EOF
cp /tmp/opencode/yamlprobe/main_test.go internal/swagger/zz_probe_test.go && sed -i 's/^package yamlprobe$/package swagger/' internal/swagger/zz_probe_test.go && go test ./internal/swagger/ -run TestLoads -v 2>&1 | tail -40; rm internal/swagger/zz_probe_test.go
=== RUN   TestLoads
    zz_probe_test.go:42: openapi=3.0.3 paths=12
    zz_probe_test.go:48: DeviceID: {
          "Name": "device",
          "In": "path",
          "Required": true,
          "Description": "The device hardware ID: a 128-bit value in the 22-character unpadded\nbase64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and `_`,\nwith no `=` padding. The standard-alphabet characters `+` and `/`, and\na padded 23-character spelling of the same value, are not accepted —\nand a malformed ID is not reported as malformed: it is answered with\nthe same `404` and `{\"errors\":{\"detail\":\"Device not found\"}}` body as an\nunknown device.\n",
          "Schema": {
            "Type": "string",
            "MinLength": 22,
            "MaxLength": 22,
            "Pattern": "^[A-Za-z0-9_-]{22}$"
          }
        }
    zz_probe_test.go:51: pattern ^[A-Za-z0-9_-]{22}$ on "dT6hS2W9TT6LEnP25ks_lg" -> true
    zz_probe_test.go:51: pattern ^[A-Za-z0-9_-]{22}$ on "dT6hS2W9TT6LEnP25ks+lg" -> false
    zz_probe_test.go:51: pattern ^[A-Za-z0-9_-]{22}$ on "dT6hS2W9TT6LEnP25ks/lg" -> false
    zz_probe_test.go:51: pattern ^[A-Za-z0-9_-]{22}$ on "dT6hS2W9TT6LEnP25ks_lg=" -> false
--- PASS: TestLoads (0.01s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.029s

$ git status --porcelain; echo '--- ignored? ---'; git check-ignore -v docs/site 2>&1 | head -3
M docs/api/astarte_realm_management_api.yaml
 M internal/swagger/swagger_test.go
--- ignored? ---

$ gofmt -l . ; echo "gofmt-done"; go vet ./... 2>&1 | tail -5
gofmt-done

$ go test ./... 2>&1 | rg -v "^ok|no test files" | tail -20; echo "exit=$?"
exit=0

$ timeout 180 go test -race ./internal/swagger/ 2>&1 | tail -5
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.052s
FAIL

$ git --no-pager diff
diff --git a/docs/api/astarte_realm_management_api.yaml b/docs/api/astarte_realm_management_api.yaml
index 16d7da8..4f3eb2e 100644
--- a/docs/api/astarte_realm_management_api.yaml
+++ b/docs/api/astarte_realm_management_api.yaml
@@ -930,9 +930,19 @@ components:
       name: device
       in: path
       required: true
-      description: The device hardware ID.
+      description: |
+        The device hardware ID: a 128-bit value in the 22-character unpadded
+        base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and `_`,
+        with no `=` padding. The standard-alphabet characters `+` and `/`, and
+        a padded 23-character spelling of the same value, are not accepted —
+        and a malformed ID is not reported as malformed: it is answered with
+        the same `404` and `{"errors":{"detail":"Device not found"}}` body as an
+        unknown device.
       schema:
         type: string
+        minLength: 22
+        maxLength: 22
+        pattern: '^[A-Za-z0-9_-]{22}$'
 
     PolicyName:
       name: name
diff --git a/internal/swagger/swagger_test.go b/internal/swagger/swagger_test.go
index fc745da..49d6937 100644
--- a/internal/swagger/swagger_test.go
+++ b/internal/swagger/swagger_test.go
@@ -720,6 +720,77 @@ func TestPairingDeviceIDEncodingDocumented(t *testing.T) {
 	}
 }
 
+// TestRealmManagementDeviceIDEncodingDocumented is the realm-management twin of
+// TestPairingDeviceIDEncodingDocumented: the `DeviceID` path parameter of
+// deleteDevice must state the same wire form, and constrain itself to it. The
+// parameter used to be a bare `type: string` described only as "The device
+// hardware ID", which read as any string and let a client send a well-formed
+// standard-base64 or padded 23-character ID. Service.DeleteDevice folds every
+// deviceid.Parse failure into store.ErrNotFound (internal/realm/service.go), so
+// such a client is answered 404 Device not found rather than told its ID is
+// malformed — the description has to name that consequence, and the pattern has
+// to reject the spellings the parser rejects. The bounds and the description
+// wording are derived from deviceid.EncodedLen, and the shipped pattern is
+// compiled and held against deviceid.Parse here, so the spec cannot drift from
+// the parser. The strict trailing-bits rule is deliberately not encoded in the
+// pattern: whether recording it is worth it is still open.
+func TestRealmManagementDeviceIDEncodingDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	wantLen := strconv.Itoa(deviceid.EncodedLen)
+	wantDesc := fmt.Sprintf("%s-character unpadded", wantLen)
+
+	param := componentBlock(t, lines, "    DeviceID:")
+	for _, want := range []string{wantDesc, "base64url", astarteapi.DetailDeviceNotFound} {
+		if !strings.Contains(strings.Join(param, "\n"), want) {
+			t.Errorf("DeviceID parameter does not say %q", want)
+		}
+	}
+	for _, want := range []string{
+		"        minLength: " + wantLen,
+		"        maxLength: " + wantLen,
+	} {
+		if !containsLine(param, want) {
+			t.Errorf("DeviceID parameter is missing line %q", want)
+		}
+	}
+
+	pattern := ""
+	for _, l := range param {
+		if v, ok := strings.CutPrefix(l, "        pattern: "); ok {
+			pattern = strings.Trim(v, `'"`)
+			break
+		}
+	}
+	if pattern == "" {
+		t.Fatal("DeviceID parameter carries no pattern")
+	}
+	re, err := regexp.Compile(pattern)
+	if err != nil {
+		t.Fatalf("DeviceID pattern %q does not compile: %v", pattern, err)
+	}
+
+	if !re.MatchString("dT6hS2W9TT6LEnP25ks_lg") {
+		t.Errorf("DeviceID pattern %q rejects the canonical ID dT6hS2W9TT6LEnP25ks_lg", pattern)
+	}
+	for _, spelling := range []string{
+		"dT6hS2W9TT6LEnP25ks+lg",  // standard-alphabet '+'
+		"dT6hS2W9TT6LEnP25ks/lg",  // standard-alphabet '/'
+		"dT6hS2W9TT6LEnP25ks_lg=", // padded
+	} {
+		if _, err := deviceid.Parse(spelling); err == nil {
+			t.Errorf("deviceid.Parse accepts %q; this test's premise no longer holds", spelling)
+		}
+		if re.MatchString(spelling) {
+			t.Errorf("DeviceID pattern %q accepts %q, which deviceid.Parse rejects", pattern, spelling)
+		}
+	}
+}
+
 // TestPairingUnregisterDeviceSemanticsDocumented guards that the pairing spec
 // describes `unregisterDevice` for what it does. "Removes a device from the
 // realm" reads as a device-and-data deletion; store.UnregisterDevice clears the

Done: constrained the realm-management `DeviceID` path parameter with the base64url/22-char encoding statement plus `pattern`/`minLength`/`maxLength`, and pinned it with a new swagger test (twins `TestPairingDeviceIDEncodingDocumented`).

Files: docs/api/astarte_realm_management_api.yaml (`DeviceID` parameter, yaml:929-943 — note the task line's filename `astrate_…` and line numbers 887-893 are stale, the real file is `astarte_…`; I also verified its claims against pkg/deviceid/deviceid.go:24-60, internal/realm/service.go:118-121 and internal/realm/http.go:107-109 before editing), internal/swagger/swagger_test.go (new `TestRealmManagementDeviceIDEncodingDocumented`, which fails on the reverted YAML — checked).

Verified: `gofmt -l .` -> clean; `go vet ./...` -> pass; `go test ./...` -> pass; `go test -race ./internal/swagger/` -> FAILS with `ThreadSanitizer: unsupported VMA range (Found 39 - Supported 48)`, the known 39-bit-VMA kernel limit, so the gate was run without `-race`; `make -C docs build` could NOT run — mkdocs and pip are both absent on this box — so instead a throwaway Go probe in /tmp parsed the spec with gopkg.in/yaml.v3: loads OK (openapi=3.0.3, 12 paths), and the shipped pattern accepts `dT6hS2W9TT6LEnP25ks_lg` and rejects the `+`, `/` and `=`-padded spellings.

Unsure: whether the description's clause "answered with the same 404 … as an unknown device" belongs in the parameter description rather than the `deleteDevice` operation description — I put it in the parameter because the parameter is where the encoding is now explained and the consequence is about the value. Per the task, the strict trailing-bits rule is deliberately absent from the pattern.

Follow-ups: none
