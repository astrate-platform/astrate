slug: docs-sync-pairing-unregister-description
verdict: done
at:  8ffb39e
ran: 2026-09-28T18:01:20Z on DietPi in 261s


> build · big-pickle

I'll read the required files first.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ sed -n 130,175p docs/api/astarte_pairing_api.yaml
$ref: "#/components/schemas/DataEnvelopeCredentialsSecret"
              example:
                data:
                  credentials_secret: TTkd0vDHgAs6FSyhmnnpfMZlDOdUWLNxNHbLeuTpegc=
        "400":
          $ref: "#/components/responses/BadRequest"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "404":
          $ref: "#/components/responses/DeviceNotFound"
        "422":
          $ref: "#/components/responses/ValidationError"
        "429":
          $ref: "#/components/responses/TooManyRequests"
        "500":
          $ref: "#/components/responses/InternalServerError"

  /pairing/v1/{realm}/agent/devices/{deviceID}:
    delete:
      operationId: unregisterDevice
      summary: Unregister a device
      description: Removes a device from the realm. Returns 204 on success.
      tags: [Pairing]
      security:
        - a_pa: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - $ref: "#/components/parameters/DeviceID"
      responses:
        "204":
          description: Device unregistered successfully. No body returned.
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "404":
          $ref: "#/components/responses/DeviceNotFound"
        "500":
          $ref: "#/components/responses/InternalServerError"

  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
    post:
      operationId: requestCredentials
      summary: Request device credentials (flow B)

$ sed -n 85,115p internal/store/devices.go; echo ---; sed -n 218,240p internal/pairing/service.go
return fmt.Errorf("store: registering device %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: device %s", ErrDeviceAlreadyConfirmed, id)
	}
	return nil
}

// UnregisterDevice makes a device registrable again without losing its data
// (DELETE /agent/devices parity, docs/DESIGN.md §4.4): the credentials
// secret and certificate trail are cleared, the row and all stored data stay.
func (s *Store) UnregisterDevice(ctx context.Context, realmID int16, id deviceid.ID) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE devices
		SET credentials_secret_hash = '', status = 'registered',
		    first_credentials_request = NULL, cert_serial = NULL, cert_aki = NULL
		WHERE realm_id = $1 AND id = $2`,
		realmID, uuidParam(id))
	if err != nil {
		return fmt.Errorf("store: unregistering device %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: device %s", ErrNotFound, id)
	}
	return nil
}

// GetDevice fetches one device.
func (s *Store) GetDevice(ctx context.Context, realmID int16, id deviceid.ID) (*Device, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT `+deviceColumns+` FROM devices WHERE realm_id = $1 AND id = $2`,
---
			return "", err
		}
	}
	if s.OnRegistered != nil {
		s.OnRegistered(realmName, hwID, time.Now())
	}
	return secret, nil
}

// Unregister implements the flow A DELETE: the device becomes registrable
// again, its data is kept (store.UnregisterDevice clears only the
// credential trail). store.ErrNotFound is returned for unknown devices.
func (s *Service) Unregister(ctx context.Context, realmName, deviceIDStr string) error {
	id, err := deviceid.Parse(deviceIDStr)
	if err != nil {
		return fmt.Errorf("%w: device %q", store.ErrNotFound, deviceIDStr)
	}
	realm, err := s.st.GetRealmByName(ctx, realmName)
	if err != nil {
		return err
	}
	return s.st.UnregisterDevice(ctx, realm.ID, id)
}

$ rg -n "Unregister" --glob '!docs/**' -g '!*.md' | head -40
internal/flow/flow.go:185:// UnregisterFlow removes an instance key from the manager map. Call after
internal/flow/flow.go:187:func (m *Manager) UnregisterFlow(instanceID string) {
internal/store/devices.go:93:// UnregisterDevice makes a device registrable again without losing its data
internal/store/devices.go:96:func (s *Store) UnregisterDevice(ctx context.Context, realmID int16, id deviceid.ID) error {
internal/store/devices_test.go:134:		// Unregister keeps the row but resets the credential trail.
internal/store/devices_test.go:135:		if err := s.UnregisterDevice(ctx, realm.ID, id); err != nil {
internal/store/devices_test.go:162:		if err := s.UnregisterDevice(ctx, realm.ID, ghost); !errors.Is(err, ErrNotFound) {
internal/store/devices_test.go:163:			t.Errorf("UnregisterDevice: %v", err)
internal/pairing/http.go:85:		a.requireAgent(http.HandlerFunc(a.handleUnregister)))
internal/pairing/http.go:134:func (a *API) handleUnregister(w http.ResponseWriter, r *http.Request) {
internal/pairing/http.go:135:	err := a.svc.Unregister(r.Context(), r.PathValue("realm"), r.PathValue("deviceID"))
internal/pairing/service_test.go:75:func (f *fakeStore) UnregisterDevice(_ context.Context, realmID int16, id deviceid.ID) error {
internal/pairing/service_test.go:340:func TestUnregister(t *testing.T) {
internal/pairing/service_test.go:345:	if err := svc.Unregister(ctx, "test", hwID); !errors.Is(err, store.ErrNotFound) {
internal/pairing/service_test.go:358:	if err := svc.Unregister(ctx, "test", hwID); err != nil {
internal/pairing/service_test.go:359:		t.Fatalf("Unregister: %v", err)
internal/pairing/http_test.go:502:	t.Run("FlowA_Unregister", func(t *testing.T) {
internal/pairing/service.go:31:	UnregisterDevice(ctx context.Context, realmID int16, id deviceid.ID) error
internal/pairing/service.go:227:// Unregister implements the flow A DELETE: the device becomes registrable
internal/pairing/service.go:228:// again, its data is kept (store.UnregisterDevice clears only the
internal/pairing/service.go:230:func (s *Service) Unregister(ctx context.Context, realmName, deviceIDStr string) error {
internal/pairing/service.go:239:	return s.st.UnregisterDevice(ctx, realm.ID, id)
internal/pairing/service.go:403:		// Unregistered (or never-registered) device: no valid hash exists.
internal/flowapi/service_integration_test.go:151:	svc.Manager().UnregisterFlow(flow.InstanceID(realm, "f1"))
internal/flowapi/service.go:459:	s.mgr.UnregisterFlow(instanceID)
internal/flowapi/service.go:481:	s.mgr.UnregisterFlow(instanceID)
internal/flowapi/service.go:545:		s.mgr.UnregisterFlow(instanceID)

$ sed -n 128,150p internal/pairing/http.go
a.writeServiceError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusCreated, registerResponse{CredentialsSecret: secret})
}

func (a *API) handleUnregister(w http.ResponseWriter, r *http.Request) {
	err := a.svc.Unregister(r.Context(), r.PathValue("realm"), r.PathValue("deviceID"))
	if err != nil {
		a.writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- flow B/C: device endpoints ---------------------------------------------

// credentialsRequest is the flow B body.
type credentialsRequest struct {
	CSR string `json:"csr"`
}

// credentialsResponse carries the issued client certificate.

$ rg -n "description: (\||>|\")" docs/api/astarte_pairing_api.yaml | head -20
4:  description: |
31:      description: |
74:      description: |
101:      description: |
176:      description: |
228:      description: |
262:      description: |
335:      description: |
342:      description: |
360:      description: |
380:          description: |
391:          description: |

→ Read docs/api/astarte_pairing_api.yaml [offset=24, limit=60]
Now the edit, matching the block-scalar style used elsewhere in the file:

← Edit docs/api/astarte_pairing_api.yaml
Index: /root/astrate-mule/docs/api/astarte_pairing_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_pairing_api.yaml
+++ /root/astrate-mule/docs/api/astarte_pairing_api.yaml
@@ -149,9 +149,20 @@
 /pairing/v1/{realm}/agent/devices/{deviceID}:
   delete:
     operationId: unregisterDevice
     summary: Unregister a device
-    description: Removes a device from the realm. Returns 204 on success.
+    description: |
+      Makes a device registrable again. This is not a deletion: the device
+      row survives together with everything stored for it — its interfaces,
+      its datastream data and its group memberships. Only the credential
+      trail is cleared (credentials secret, first-credentials-request
+      timestamp and client certificate) and `status` goes back to
+      `registered`, so the device has to be paired again to obtain a
+      credentials secret or a certificate.
+
+      Returns 204 on success. Because the row is still there, unregistering
+      an already-unregistered device is also a 204 no-op; 404 is returned
+      only for a device that does not exist in the realm.
     tags: [Pairing]
     security:
       - a_pa: []
     parameters:

$ python3 -c "
import yaml,sys
d=yaml.safe_load(open('docs/api/astarte_pairing_api.yaml'))
p=d['paths']['/pairing/v1/{realm}/agent/devices/{deviceID}']['delete']
print(p['operationId']); print(p['description'])
"
Traceback (most recent call last):
  File "<string>", line 2, in <module>
    import yaml,sys
ModuleNotFoundError: No module named 'yaml'

$ cat docs/Makefile
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

$ cat .gitignore 2>/dev/null; echo "--- site tracked?"; git ls-files docs/site | head -5; echo "--- mkdocs?"; command -v mkdocs; pip --version 2>/dev/null
# Build artifacts
dist/
/astrate
*.exe
*.test
*.out

# Coverage
coverage.*
*.cover

# bbolt session stores (broker persistence, §3.1)
*.db

# Local env / editor
.env
.envrc
.idea/
.vscode/
*.swp
.DS_Store

# Local AI documentation workflow. The handoff files under docs/handoff/ ARE tracked:
# they are the project's phase history, and the agents that need it most are the cheap ones
# that cannot re-derive it.
/docs/AI-DOCUMENTATION-WORKFLOW.md

# MkDocs build output
/site-dist/

# Docs site copies (refreshed by docs/make sync)
/docs/site/DESIGN.md
/docs/site/ROADMAP.md
/docs/site/api/
/docs/site/swagger-ui/

# Local-only compose overrides
docker-compose.override.yml
__pycache__/

# bench harness artifacts (bench/README.md). Note bench/results/ is deliberately NOT here:
# a benchmark number without its host and its history is worthless, so results are evidence
# and evidence is versioned.
/bench/.astarte/
/bench/keys/
/bench/astrate/conf/
/bench/*.json
/bench/bench

# The mule host's tick transcript: local to whichever machine ran it, grows without
# bound, and got committed by accident on 2026-09-04.
.mule/cron.log
--- site tracked?
docs/site/appengine-api.md
docs/site/architecture.md
docs/site/compatibility.md
docs/site/configuration-reference.md
docs/site/contributing.md
--- mkdocs?

$ make -C docs build 2>&1 | tail -20
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

$ rg -n "yaml" go.mod; ls /root/go/pkg/mod/gopkg.in 2>/dev/null
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect
check.v1@v0.0.0-20161208181325-20d25e280405
check.v1@v1.0.0-20201130134442-10cb98267c6c
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1

$ timeout 40 pip install -q -r docs/requirements.txt 2>&1 | tail -5; echo "exit=$?"
timeout: failed to run command ‘pip’: No such file or directory
exit=0

$ ls docs/swagger-ui/; rg -n "yaml" docs/swagger-ui/* docs/mkdocs.yml 2>/dev/null | head
index.html
docs/swagger-ui/index.html:44:      <option value="../api/astarte_housekeeping_api.yaml">Housekeeping API — Realm management, a_ha auth</option>
docs/swagger-ui/index.html:45:      <option value="../api/astarte_pairing_api.yaml">Pairing API — Device pairing flows A/B/C, a_pa auth</option>
docs/swagger-ui/index.html:46:      <option value="../api/astarte_realm_management_api.yaml">Realm Management API — Interfaces, triggers, policies, a_rma auth</option>
docs/swagger-ui/index.html:47:      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
docs/swagger-ui/index.html:48:      <option value="../api/astrate_native_api.yaml">Astrate Native API — Health, readiness, metrics, WebSocket</option>

$ cat > /tmp/opencode/yamlprobe.go <<'EOF'
package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	for _, f := range os.Args[1:] {
		b, err := os.ReadFile(f)
		if err != nil {
			fmt.Println("READ FAIL", f, err)
			os.Exit(1)
		}
		var doc struct {
			Paths map[string]map[string]struct {
				OperationID string `yaml:"operationId"`
				Description string `yaml:"description"`
			} `yaml:"paths"`
		}
		if err := yaml.Unmarshal(b, &doc); err != nil {
			fmt.Println("PARSE FAIL", f, err)
			os.Exit(1)
		}
		for p, ops := range doc.Paths {
			for m, op := range ops {
				if op.OperationID == "unregisterDevice" {
					fmt.Printf("OK %s\n%s %s\n--- description ---\n%s\n", f, m, p, op.Description)
				}
			}
		}
	}
}
EOF
go run /tmp/opencode/yamlprobe.go docs/api/astarte_pairing_api.yaml
OK docs/api/astarte_pairing_api.yaml
delete /pairing/v1/{realm}/agent/devices/{deviceID}
--- description ---
Makes a device registrable again. This is not a deletion: the device
row survives together with everything stored for it — its interfaces,
its datastream data and its group memberships. Only the credential
trail is cleared (credentials secret, first-credentials-request
timestamp and client certificate) and `status` goes back to
`registered`, so the device has to be paired again to obtain a
credentials secret or a certificate.

Returns 204 on success. Because the row is still there, unregistering
an already-unregistered device is also a 204 no-op; 404 is returned
only for a device that does not exist in the realm.

$ rg -ln "api/astrate|astarte_pairing_api|docs/api" --glob '*.go' | head
internal/swagger/swagger_test.go

$ sed -n 1,80p internal/swagger/swagger_test.go
package swagger

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	docs "github.com/astrate-platform/astrate/docs"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

func TestMount(t *testing.T) {
	mux := http.NewServeMux()
	Mount(mux)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	t.Run("GET /swagger redirects to /swagger/index.html", func(t *testing.T) {
		client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}}
		resp, err := client.Get(srv.URL + "/swagger")
		if err != nil {
			t.Fatalf("GET /swagger: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusFound {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusFound)
		}
		loc := resp.Header.Get("Location")
		if loc != "/swagger/index.html" {
			t.Errorf("Location = %q, want %q", loc, "/swagger/index.html")
		}
	})

	t.Run("GET /swagger/ serves the embedded UI", func(t *testing.T) {
		body := get(t, srv.URL+"/swagger/index.html")
		want, err := docs.SwaggerUI.ReadFile("swagger-ui/index.html")
		if err != nil {
			t.Fatalf("reading embedded index.html: %v", err)
		}
		if body != string(want) {
			t.Errorf("served index.html does not match embedded copy")
		}
	})

	t.Run("GET /api/ serves every OpenAPI YAML spec", func(t *testing.T) {
		for _, name := range Specs() {
			body := get(t, srv.URL+"/api/"+name)
			want, err := docs.APIYAML.ReadFile("api/" + name)
			if err != nil {
				t.Fatalf("reading embedded %s: %v", name, err)
			}
			if body != string(want) {
				t.Errorf("served /api/%s does not match embedded copy", name)
			}
		}
	})
}

// emptyEmbed simulates the day the docs embed layout drops the swagger-ui/ or
// api/ tree: an embed.FS declared without a //go:embed directive is empty, and
// rising fs.Sub over it silently yields an empty tree, exactly what Mount must
// fail fast on instead of serving 404s at /swagger/ and /api/.
var emptyEmbed embed.FS

// TestMountSubPanicsOnBrokenFS guards the fail-fast contract: Mount must panic
// when a docs embed sub-tree is absent rather than silently serve an empty
// tree that 404s at /swagger/ and /api/.

$ rg -n "unregister|Unregister|func Test" internal/swagger/swagger_test.go | head -40
22:func TestMount(t *testing.T) {
81:func TestMountSubPanicsOnBrokenFS(t *testing.T) {
94:func TestSpecs(t *testing.T) {
127:func TestRealmManagement403(t *testing.T) {
138:func TestAppEngine403(t *testing.T) {
202:func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
245:func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
320:func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {
380:func TestAppEngineAddGroupDevice422Documented(t *testing.T) {
424:func TestAppEngineDataDelete400Documented(t *testing.T) {
469:func TestPairingErrorDetailExamples(t *testing.T) {
518:func TestPairingDeviceIDEncodingDocumented(t *testing.T) {

$ sed -n 469,560p internal/swagger/swagger_test.go
func TestPairingErrorDetailExamples(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	const detailPrefix = "              detail: "

	for _, tc := range []struct {
		response string
		want     string
	}{
		{"    BadRequest:", astarteapi.DetailBadRequest},
		{"    Unauthorized:", astarteapi.DetailUnauthorized},
		{"    Forbidden:", astarteapi.DetailForbidden},
		{"    DeviceNotFound:", astarteapi.DetailDeviceNotFound},
		{"    InternalServerError:", astarteapi.DetailInternalServerError},
	} {
		block := componentBlock(t, lines, tc.response)

		got := ""
		for _, l := range block {
			if strings.HasPrefix(l, detailPrefix) {
				got = strings.TrimPrefix(l, detailPrefix)
				break
			}
		}
		if got == "" {
			t.Errorf("response %q carries no %q example", tc.response, "detail")
			continue
		}
		if got != tc.want {
			t.Errorf("response %q example detail = %q, want the canonical %q",
				strings.TrimSpace(tc.response), got, tc.want)
		}
	}
}

// TestPairingDeviceIDEncodingDocumented guards that the pairing spec states the
// exact wire form of a device identifier — the 22-character unpadded base64url
// string deviceid.Parse accepts — instead of the "base64-encoded 128-bit"
// wording, which read as standard base64 and let a generated client send a
// well-formed `+`/`/` identifier or a padded 23-character spelling and collect
// a 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}`
// (internal/pairing/http.go). The shipped pattern is compiled here and held to
// the same spellings the parser is, so the two cannot drift apart: the bounds
// and the description wording are derived from deviceid.EncodedLen rather than
// written out, so a change to the accepted length fails this test too.
func TestPairingDeviceIDEncodingDocumented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	wantLen := strconv.Itoa(deviceid.EncodedLen)
	wantDesc := fmt.Sprintf("%s-character unpadded", wantLen)

	param := strings.Join(componentBlock(t, lines, "    DeviceID:"), "\n")
	for _, want := range []string{wantDesc, "base64url"} {
		if !strings.Contains(param, want) {
			t.Errorf("DeviceID parameter description does not say %q", want)
		}
	}
	if strings.Contains(param, "(base64-encoded 128-bit)") {
		t.Error("DeviceID parameter still describes the ID as plain base64")
	}

	block := propertyBlock(t, lines, "hw_id")
	hwID := strings.Join(block, "\n")
	for _, want := range []string{wantDesc, "base64url"} {
		if !strings.Contains(hwID, want) {
			t.Errorf("RegisterRequest.hw_id description does not say %q", want)
		}
	}
	for _, want := range []string{
		"          minLength: " + wantLen,
		"          maxLength: " + wantLen,
	} {
		if !containsLine(block, want) {
			t.Errorf("RegisterRequest.hw_id is missing line %q", want)
		}
	}

	pattern := ""
	for _, l := range block {
		if v, ok := strings.CutPrefix(l, "          pattern: "); ok {
			pattern = strings.Trim(v, `'"`)
			break
		}
	}

$ rg -n "^func (componentBlock|propertyBlock|containsLine|operationBlock)" -A 25 internal/swagger/swagger_test.go | head -80; wc -l internal/swagger/swagger_test.go
284:func propertyBlock(t *testing.T, lines []string, name string) []string {
285-	t.Helper()
286-	key := "        " + name + ":"
287-	start := -1
288-	for i, l := range lines {
289-		if l == key {
290-			start = i
291-			break
292-		}
293-	}
294-	if start < 0 {
295-		t.Fatalf("schema declares no property %q", name)
296-	}
297-	for i := start + 1; i < len(lines); i++ {
298-		if strings.TrimSpace(lines[i]) == "" {
299-			continue
300-		}
301-		if !strings.HasPrefix(lines[i], "          ") {
302-			return lines[start+1 : i]
303-		}
304-	}
305-	return lines[start+1:]
306-}
307-
308-// TestHousekeepingRetentionZeroFoldDocumented guards that the housekeeping spec
309-// tells clients that `datastream_maximum_storage_retention: 0` means unset, not
--
588:func operationBlock(t *testing.T, lines []string, operationID string) []string {
589-	t.Helper()
590-	marker := "      operationId: " + operationID
591-	start := -1
592-	for i, l := range lines {
593-		if l == marker {
594-			start = i
595-			break
596-		}
597-	}
598-	if start < 0 {
599-		t.Fatalf("spec declares no operation %q", operationID)
600-	}
601-	endRe := regexp.MustCompile(`^      operationId: |^  \S|^components:`)
602-	for i := start + 1; i < len(lines); i++ {
603-		if endRe.MatchString(lines[i]) {
604-			return lines[start+1 : i]
605-		}
606-	}
607-	return lines[start+1:]
608-}
609-
610-// componentBlock returns the lines of the components entry starting with the
611-// given key line, up to the next entry at the same or shallower indentation.
612:func componentBlock(t *testing.T, lines []string, key string) []string {
613-	t.Helper()
614-	start := -1
615-	for i, l := range lines {
616-		if l == key {
617-			start = i
618-			break
619-		}
620-	}
621-	if start < 0 {
622-		t.Fatalf("spec declares no components entry %q", key)
623-	}
624-	for i := start + 1; i < len(lines); i++ {
625-		if strings.TrimSpace(lines[i]) == "" {
626-			continue
627-		}
628-		if !strings.HasPrefix(lines[i], "      ") {
629-			return lines[start+1 : i]
630-		}
631-	}
632-	return lines[start+1:]
633-}
634-
635:func containsLine(lines []string, want string) bool {
636-	for _, l := range lines {
637-		if l == want {
638-			return true
639-		}
640-	}
678 internal/swagger/swagger_test.go

$ rg -n "operationBlock" internal/swagger/swagger_test.go; sed -n 200,284p internal/swagger/swagger_test.go
211:		block := operationBlock(t, lines, op)
257:		block := operationBlock(t, lines, op)
327:	patch := strings.Join(operationBlock(t, lines, "patchRealm"), "\n")
392:	block := operationBlock(t, lines, "addGroupDevice")
437:		block := operationBlock(t, lines, op)
586:// operationBlock returns the lines of the operation with the given operationId,
588:func operationBlock(t *testing.T, lines []string, operationID string) []string {
// the documentation half of the behaviour pinned by
// TestHousekeepingAsyncOperationParam in internal/housekeeping.
func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	const ref = `        - $ref: "#/components/parameters/AsyncOperation"`
	for _, op := range []string{"createRealm", "deleteRealm"} {
		block := operationBlock(t, lines, op)
		if !containsLine(block, ref) {
			t.Errorf("operation %s does not $ref the AsyncOperation parameter", op)
		}
	}

	comp := componentBlock(t, lines, "    AsyncOperation:")
	for _, want := range []string{
		"      name: async_operation",
		"      in: query",
		"      required: false",
		"      schema:",
		"        type: boolean",
		"        default: false",
	} {
		if !containsLine(comp, want) {
			t.Errorf("components parameter AsyncOperation is missing line %q", want)
		}
	}
	if !strings.Contains(strings.Join(comp, "\n"), "Accepted and ignored") {
		t.Error("components parameter AsyncOperation does not say the value is accepted and ignored")
	}
}

// TestRealmManagementAsyncOperationParamDocumented guards that the realm
// management spec tells clients the `?async_operation` parameter exists on the
// five operations upstream 1.4 runs in the background: interface
// install/update/delete, device deletion and trigger-delivery-policy delete.
// Astrate accepts and ignores it there (deviation 17, docs/COMPATIBILITY.md) —
// none of the five handlers reads the query string, which is the whole of the
// acceptance — so an upstream client keeps working, but a client generated from
// the spec has to learn the parameter from the spec. This is the documentation
// half of the behaviour pinned by TestRealmManagementAsyncOperationParam in
// internal/realm.
func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	const ref = `        - $ref: "#/components/parameters/AsyncOperation"`
	for _, op := range []string{
		"installInterface", "updateInterface", "deleteInterface",
		"deleteDevice", "deletePolicy",
	} {
		block := operationBlock(t, lines, op)
		if !containsLine(block, ref) {
			t.Errorf("operation %s does not $ref the AsyncOperation parameter", op)
		}
	}

	comp := componentBlock(t, lines, "    AsyncOperation:")
	for _, want := range []string{
		"      name: async_operation",
		"      in: query",
		"      required: false",
		"      schema:",
		"        type: boolean",
		"        default: false",
	} {
		if !containsLine(comp, want) {
			t.Errorf("components parameter AsyncOperation is missing line %q", want)
		}
	}
	if !strings.Contains(strings.Join(comp, "\n"), "Accepted and ignored") {
		t.Error("components parameter AsyncOperation does not say the value is accepted and ignored")
	}
}

// propertyBlock returns the lines of the schema property with the given name
// inside a components.schemas entry, up to the next property or the end of the
// entry.
func propertyBlock(t *testing.T, lines []string, name string) []string {

→ Read internal/swagger/swagger_test.go [offset=558, limit=32]

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -582,8 +582,67 @@
 		}
 	}
 }
 
+// TestPairingUnregisterDeviceSemanticsDocumented guards that the pairing spec
+// describes `unregisterDevice` for what it does. "Removes a device from the
+// realm" reads as a device-and-data deletion; store.UnregisterDevice clears the
+// credential trail and flips status back to 'registered' and touches nothing
+// else (internal/store/devices.go), so the row, its interfaces, its datastream
+// data and its group memberships all survive, and a second DELETE still matches
+// the row — a 204 no-op, not the 404 a deletion would have produced.
+func TestPairingUnregisterDeviceSemanticsDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	desc := operationDescription(t, operationBlock(t, lines, "unregisterDevice"))
+	for _, want := range []string{
+		"not a deletion",
+		"registrable again",
+		"interfaces",
+		"datastream data",
+		"group memberships",
+		"credential",
+		"204 no-op",
+		"404 is returned only for a device that does not exist",
+	} {
+		if !strings.Contains(desc, want) {
+			t.Errorf("unregisterDevice description does not say %q", want)
+		}
+	}
+	if strings.Contains(desc, "Removes a device from the realm") {
+		t.Error("unregisterDevice description still claims the device is removed from the realm")
+	}
+}
+
+// operationDescription returns the text of the operation block's block-scalar
+// `description`, with the wrapping flattened so assertions do not depend on it.
+func operationDescription(t *testing.T, block []string) string {
+	t.Helper()
+	var text []string
+	inBlock := false
+	for _, l := range block {
+		if !inBlock {
+			inBlock = l == "      description: |"
+			continue
+		}
+		if strings.TrimSpace(l) == "" {
+			continue
+		}
+		if !strings.HasPrefix(l, "        ") {
+			break
+		}
+		text = append(text, strings.TrimSpace(l))
+	}
+	if !inBlock {
+		t.Fatal("operation declares no block-scalar description")
+	}
+	return strings.Join(text, " ")
+}
+
 // operationBlock returns the lines of the operation with the given operationId,
 // up to the next operation, path, or the components section.
 func operationBlock(t *testing.T, lines []string, operationID string) []string {
 	t.Helper()

$ gofmt -l . 2>&1 | grep -v '^site-dist' ; go vet ./... && go test ./internal/swagger/ 2>&1 | tail -5
ok  	github.com/astrate-platform/astrate/internal/swagger	0.105s

Now confirm the test actually fails against the old description:
$ cp docs/api/astarte_pairing_api.yaml /tmp/opencode/pairing.new.yaml && python3 - <<'EOF'
import re
p='docs/api/astarte_pairing_api.yaml'
s=open(p).read()
i=s.index('      operationId: unregisterDevice')
j=s.index('      tags: [Pairing]', i)
old=s[i:j]
new=old.replace(old[old.index('      description: |'):], '      description: Removes a device from the realm. Returns 204 on success.\n')
open(p,'w').write(s[:i]+new+s[j:])
EOF
go test ./internal/swagger/ -run UnregisterDeviceSemantics 2>&1 | tail -14; cp /tmp/opencode/pairing.new.yaml docs/api/astarte_pairing_api.yaml
--- FAIL: TestPairingUnregisterDeviceSemanticsDocumented (0.01s)
    swagger_test.go:600: operation declares no block-scalar description
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.030s
FAIL

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -616,31 +616,33 @@
 		t.Error("unregisterDevice description still claims the device is removed from the realm")
 	}
 }
 
-// operationDescription returns the text of the operation block's block-scalar
-// `description`, with the wrapping flattened so assertions do not depend on it.
+// operationDescription returns the text of the operation block's description,
+// block-scalar or single-line, with the wrapping flattened so assertions do not
+// depend on it.
 func operationDescription(t *testing.T, block []string) string {
 	t.Helper()
 	var text []string
-	inBlock := false
-	for _, l := range block {
-		if !inBlock {
-			inBlock = l == "      description: |"
-			continue
+	for i, l := range block {
+		if l == "      description: |" {
+			for _, b := range block[i+1:] {
+				if strings.TrimSpace(b) == "" {
+					continue
+				}
+				if !strings.HasPrefix(b, "        ") {
+					break
+				}
+				text = append(text, strings.TrimSpace(b))
+			}
+			return strings.Join(text, " ")
 		}
-		if strings.TrimSpace(l) == "" {
-			continue
+		if v, ok := strings.CutPrefix(l, "      description: "); ok {
+			return strings.Trim(strings.TrimSpace(v), `"'`)
 		}
-		if !strings.HasPrefix(l, "        ") {
-			break
-		}
-		text = append(text, strings.TrimSpace(l))
 	}
-	if !inBlock {
-		t.Fatal("operation declares no block-scalar description")
-	}
-	return strings.Join(text, " ")
+	t.Fatal("operation declares no description")
+	return ""
 }
 
 // operationBlock returns the lines of the operation with the given operationId,
 // up to the next operation, path, or the components section.

$ go test ./internal/swagger/ -run UnregisterDeviceSemantics 2>&1 | tail -3 && python3 - <<'EOF'
p='docs/api/astarte_pairing_api.yaml'
s=open(p).read()
i=s.index('      operationId: unregisterDevice')
j=s.index('      tags: [Pairing]', i)
old=s[i:j]
new=old.replace(old[old.index('      description: |'):], '      description: Removes a device from the realm. Returns 204 on success.\n')
open(p,'w').write(s[:i]+new+s[j:])
EOF
go test ./internal/swagger/ -run UnregisterDeviceSemantics 2>&1 | tail -14; cp /tmp/opencode/pairing.new.yaml docs/api/astarte_pairing_api.yaml
ok  	github.com/astrate-platform/astrate/internal/swagger	0.034s
--- FAIL: TestPairingUnregisterDeviceSemanticsDocumented (0.00s)
    swagger_test.go:612: unregisterDevice description does not say "not a deletion"
    swagger_test.go:612: unregisterDevice description does not say "registrable again"
    swagger_test.go:612: unregisterDevice description does not say "interfaces"
    swagger_test.go:612: unregisterDevice description does not say "datastream data"
    swagger_test.go:612: unregisterDevice description does not say "group memberships"
    swagger_test.go:612: unregisterDevice description does not say "credential"
    swagger_test.go:612: unregisterDevice description does not say "204 no-op"
    swagger_test.go:612: unregisterDevice description does not say "404 is returned only for a device that does not exist"
    swagger_test.go:616: unregisterDevice description still claims the device is removed from the realm
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.022s
FAIL

$ gofmt -l . | grep -v '^site-dist/' ; go vet ./... && go test ./... 2>&1 | grep -v "^ok\|no test files" | head -20; echo "GATE_DONE"
GATE_DONE

$ git status --short; echo ---; git diff --stat
M docs/api/astarte_pairing_api.yaml
 M internal/swagger/swagger_test.go
---
 docs/api/astarte_pairing_api.yaml | 13 ++++++++-
 internal/swagger/swagger_test.go  | 61 +++++++++++++++++++++++++++++++++++++++
 2 files changed, 73 insertions(+), 1 deletion(-)

Done: rewrote the `unregisterDevice` description in `docs/api/astarte_pairing_api.yaml` to state the real semantics (device becomes registrable again; row, interfaces, datastream data and group memberships retained; only the credential trail cleared and `status` back to `registered`; repeat DELETE is a 204 no-op, 404 only for an unknown device), and added `TestPairingUnregisterDeviceSemanticsDocumented` in `internal/swagger/swagger_test.go` to lock the wording to `store.UnregisterDevice`.
Files: docs/api/astarte_pairing_api.yaml, internal/swagger/swagger_test.go
Verified: `gofmt -l .` -> clean; `go vet ./... && go test ./...` -> pass (no `-race`, docs/test-only change); new test fails on the old wording and passes on the new one; YAML re-parsed with a throwaway `gopkg.in/yaml.v3` probe (the file `docs/swagger-ui/index.html:45` points at) -> loads; `make -C docs build` -> **fail**: `mkdocs: No such file or directory` (exit 127) after its `sync` step — this host has neither `mkdocs` nor `pip`, so the docs site cannot be rendered here, and the YAML-load check was done with the Go probe instead
Unsure: whether `make -C docs build` passing matters to the gate; also whether the wording should name `credentials_secret_hash`/`cert_serial` explicitly rather than in prose
Follow-ups: none
