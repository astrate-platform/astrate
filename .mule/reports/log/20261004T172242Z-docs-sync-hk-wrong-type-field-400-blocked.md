slug: docs-sync-hk-wrong-type-field-400
verdict: blocked
at:  d954574
ran: 2026-10-04T17:22:42Z on DietPi in 152s


> build · big-pickle

I'll read the two files first.
→ Read .mule/MULE.md
→ Read .mule/task.md
Task in one line: widen the `BadRequest` description in the housekeeping API spec for the two realm operations to cover wrong-typed `data` fields, without touching the frozen example.
→ Read docs/api/astarte_housekeeping_api.yaml [offset=300, limit=140]

$ ls /root/astrate-mule/docs; echo ---; ls /root/astrate-mule/docs/Makefile 2>/dev/null || ls /root/astrate-mule/Makefile
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
/root/astrate-mule/docs/Makefile

$ cd /root/astrate-mule && rg -n "BadRequest|responses/|^  /" docs/api/astarte_housekeeping_api.yaml
26:  /housekeeping/v1/realms:
44:          $ref: "#/components/responses/Unauthorized"
46:          $ref: "#/components/responses/Forbidden"
48:          $ref: "#/components/responses/InternalServerError"
96:          $ref: "#/components/responses/BadRequest"
98:          $ref: "#/components/responses/Unauthorized"
100:          $ref: "#/components/responses/Forbidden"
102:          $ref: "#/components/responses/Conflict"
104:          $ref: "#/components/responses/ValidationError"
106:          $ref: "#/components/responses/InternalServerError"
108:  /housekeeping/v1/realms/{realm}:
134:          $ref: "#/components/responses/Unauthorized"
136:          $ref: "#/components/responses/Forbidden"
138:          $ref: "#/components/responses/NotFound"
140:          $ref: "#/components/responses/InternalServerError"
182:          $ref: "#/components/responses/BadRequest"
184:          $ref: "#/components/responses/Unauthorized"
186:          $ref: "#/components/responses/Forbidden"
188:          $ref: "#/components/responses/NotFound"
190:          $ref: "#/components/responses/ValidationErrors"
192:          $ref: "#/components/responses/InternalServerError"
208:          $ref: "#/components/responses/Unauthorized"
210:          $ref: "#/components/responses/Forbidden"
212:          $ref: "#/components/responses/NotFound"
214:          $ref: "#/components/responses/MethodNotAllowed"
216:          $ref: "#/components/responses/UnprocessableDeletion"
218:          $ref: "#/components/responses/InternalServerError"
395:    BadRequest:

→ Read internal/housekeeping/http.go [offset=40, limit=160]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=130, offset=100]

$ cd /root/astrate-mule && cat docs/Makefile
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

$ cd /root/astrate-mule && rg -n "TestHousekeepingErrorDetailExamples" -A 30 --glob '*_test.go' | head -50; echo ===; rg -n "func WriteBadRequest" -A 8 internal/astarteapi/
internal/swagger/swagger_test.go:666:// TestHousekeepingErrorDetailExamples pins the error-detail examples in the
internal/swagger/swagger_test.go-667-// housekeeping spec's components.responses to the frozen canonical strings the
internal/swagger/swagger_test.go-668-// wire emits (astarteapi/envelope.go), the same guard the pairing, realm
internal/swagger/swagger_test.go-669-// management, app engine and native specs get from
internal/swagger/swagger_test.go-670-// TestPairingErrorDetailExamples, TestRealmManagementErrorDetailExamples,
internal/swagger/swagger_test.go-671-// TestAppEngineErrorDetailExamples and TestNativeErrorDetailExamples. This is
internal/swagger/swagger_test.go-672-// the last spec still carrying a reconstructed capitalisation: its BadRequest
internal/swagger/swagger_test.go-673-// and InternalServerError examples read "Bad Request" and "Internal Server
internal/swagger/swagger_test.go-674-// Error", spellings no upstream Phoenix render answers with and no handler
internal/swagger/swagger_test.go-675-// emits — every 400 on this surface goes through astarteapi.WriteBadRequest
internal/swagger/swagger_test.go-676-// and the 500 through WriteInternalServerError (internal/housekeeping/http.go:
internal/swagger/swagger_test.go-677-// the writeError switch), and the SDK/astartectl error paths match on the
internal/swagger/swagger_test.go-678-// canonical strings. The 401 and 403 were already canonical. The responses
internal/swagger/swagger_test.go-679-// whose example is an endpoint-specific detail with no constant behind it
internal/swagger/swagger_test.go-680-// (MethodNotAllowed, Conflict, ValidationError) stay free-form and are not
internal/swagger/swagger_test.go-681-// pinned.
internal/swagger/swagger_test.go-682-//
internal/swagger/swagger_test.go-683-// There is no DeviceNotFound pin here, unlike the app engine and realm
internal/swagger/swagger_test.go-684-// management specs: this surface never answers a device 404, so
internal/swagger/swagger_test.go-685-// components.responses declares no such component to pin.
internal/swagger/swagger_test.go-686-//
internal/swagger/swagger_test.go-687-// The schema-level ErrorDetail example is pinned too, to the same canonical
internal/swagger/swagger_test.go-688-// string as BadRequest. It sits on the generic schema all eight error
internal/swagger/swagger_test.go-689-// responses $ref, so it is the example a generator falls back to for any
internal/swagger/swagger_test.go-690-// response whose own it cannot resolve — leaving it on the reconstructed
internal/swagger/swagger_test.go-691-// "Bad Request" taught a 404 and a 500 a detail no handler emits. Its shape is
internal/swagger/swagger_test.go-692-// already the wire's: the example belongs to the `errors` sub-schema, so it is
internal/swagger/swagger_test.go-693-// the value of the `errors` property and not a bare envelope, and it must not
internal/swagger/swagger_test.go-694-// be re-wrapped in another `errors:`.
internal/swagger/swagger_test.go:695:func TestHousekeepingErrorDetailExamples(t *testing.T) {
internal/swagger/swagger_test.go-696-	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
internal/swagger/swagger_test.go-697-	if err != nil {
internal/swagger/swagger_test.go-698-		t.Fatalf("reading astrate_housekeeping_api.yaml: %v", err)
internal/swagger/swagger_test.go-699-	}
internal/swagger/swagger_test.go-700-	lines := strings.Split(string(b), "\n")
internal/swagger/swagger_test.go-701-
internal/swagger/swagger_test.go-702-	const detailPrefix = "              detail: "
internal/swagger/swagger_test.go-703-
internal/swagger/swagger_test.go-704-	for _, tc := range []struct {
internal/swagger/swagger_test.go-705-		response string
internal/swagger/swagger_test.go-706-		want     string
internal/swagger/swagger_test.go-707-	}{
internal/swagger/swagger_test.go-708-		{"    BadRequest:", astarteapi.DetailBadRequest},
internal/swagger/swagger_test.go-709-		{"    Unauthorized:", astarteapi.DetailUnauthorized},
internal/swagger/swagger_test.go-710-		{"    Forbidden:", astarteapi.DetailForbidden},
internal/swagger/swagger_test.go-711-		{"    NotFound:", astarteapi.DetailNotFound},
internal/swagger/swagger_test.go-712-		{"    InternalServerError:", astarteapi.DetailInternalServerError},
internal/swagger/swagger_test.go-713-	} {
internal/swagger/swagger_test.go-714-		block := componentBlock(t, lines, tc.response)
internal/swagger/swagger_test.go-715-
===
rg: internal/astarteapi/: IO error for operation on internal/astarteapi/: No such file or directory (os error 2)

$ cd /root/astrate-mule && ls docs/swagger-ui && rg -n "yaml|json" docs/swagger-ui/* | head -20; echo ===; rg -rn "astarte_housekeeping_api" --glob '!docs/site/**' -l
index.html
44:      <option value="../api/astarte_housekeeping_api.yaml">Housekeeping API — Realm management, a_ha auth</option>
45:      <option value="../api/astarte_pairing_api.yaml">Pairing API — Device pairing flows A/B/C, a_pa auth</option>
46:      <option value="../api/astarte_realm_management_api.yaml">Realm Management API — Interfaces, triggers, policies, a_rma auth</option>
47:      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
48:      <option value="../api/astrate_native_api.yaml">Astrate Native API — Health, readiness, metrics, WebSocket</option>
===
docs/swagger-ui/index.html
docs/handoff/phase-1-handoff.md
docs/handoff/phase-1-memory.md
internal/swagger/swagger_test.go

$ cd /root/astrate-mule && rg -n "description" internal/swagger/swagger_test.go | head -30; echo === ; rg -ln "DecodeData" internal/ pkg/ 2>/dev/null | head; echo ===; rg -n "package astarteapi" -l -g '*.go' | head
337:			t.Errorf("patchRealm description does not say %q", want)
841:// and the description wording are derived from deviceid.EncodedLen rather than
856:			t.Errorf("DeviceID parameter description does not say %q", want)
867:			t.Errorf("RegisterRequest.hw_id description does not say %q", want)
919:// malformed — the description has to name that consequence, and the pattern has
920:// to reject the spellings the parser rejects. The bounds and the description
1106:			t.Errorf("unregisterDevice description does not say %q", want)
1110:		t.Error("unregisterDevice description still claims the device is removed from the realm")
1114:// operationDescription returns the text of the operation block's description,
1121:		if l == "      description: |" {
1133:		if v, ok := strings.CutPrefix(l, "      description: "); ok {
1137:	t.Fatal("operation declares no description")
===
internal/realm/http.go
pkg/astarteapi/envelope.go
pkg/astarteapi/envelope_test.go
internal/swagger/swagger_test.go
internal/housekeeping/http.go
internal/appengine/http.go
internal/flowapi/http.go
internal/pairing/http.go
===
pkg/astarteapi/envelope.go
pkg/astarteapi/envelope_test.go

$ cd /root/astrate-mule && ls internal/housekeeping/ && rg -n "httptest.NewServer|NewRouter|func Test" internal/housekeeping/*_test.go | head -40
async_operation_test.go
http.go
http_test.go
service.go
service_test.go
internal/housekeeping/async_operation_test.go:16:func TestHousekeepingAsyncOperationParam(t *testing.T) {
internal/housekeeping/http_test.go:116:func TestHousekeeping(t *testing.T) {
internal/housekeeping/http_test.go:197:func TestHousekeepingPatchRealm(t *testing.T) {
internal/housekeeping/http_test.go:309:func TestHousekeepingRetentionDefault(t *testing.T) {
internal/housekeeping/http_test.go:360:func TestHousekeepingDeleteGating(t *testing.T) {
internal/housekeeping/service_test.go:94:func TestCreateRealmValidation(t *testing.T) {
internal/housekeeping/service_test.go:130:func TestCreateRealmDefaultRetention(t *testing.T) {
internal/housekeeping/service_test.go:186:func TestCreateRealmZeroRetentionUnsets(t *testing.T) {
internal/housekeeping/service_test.go:237:func TestDeleteRealmGating(t *testing.T) {

→ Read internal/swagger/swagger_test.go [limit=70, offset=290]
→ Read internal/housekeeping/http_test.go [limit=120]

$ cd /root/astrate-mule && head -20 go.mod && echo === && rg -n "^import|^\t\"|^\t[a-z]* \"" pkg/astarteapi/envelope.go | head -20; echo === ; rg -n "DetailBadRequest|func WriteBadRequest|func DecodeData" -A 12 pkg/astarteapi/envelope.go | head -60
module github.com/astrate-platform/astrate

go 1.26.1

// Pinned dependency decisions — docs/ROADMAP.md §1.1 (frozen in docs/DESIGN.md).
// HTTP routing is stdlib net/http (Go >= 1.22 ServeMux patterns): deliberately no framework.
require (
	github.com/BurntSushi/toml v1.6.0 // TOML config (§5.1)
	github.com/coder/websocket v1.8.14 // live stream socket (M7)
	github.com/eclipse/paho.mqtt.golang v1.5.1 // test-only MQTT client (same family as official Go SDK)
	github.com/golang-jwt/jwt/v5 v5.3.1 // JWT verification, alg allowlist (§4.2)
	github.com/golang-migrate/migrate/v4 v4.19.1 // migrations via source/iofs + go:embed (§1.3)
	github.com/hashicorp/golang-lru/v2 v2.0.7 // verified-token LRU cache (§4.2)
	github.com/jackc/pgx/v5 v5.10.0 // Postgres/TimescaleDB driver + pgxpool (§1.3)
	github.com/mochi-mqtt/server/v2 v2.7.9 // embedded MQTT broker (§1.1)
	github.com/prometheus/client_golang v1.23.2 // metrics (§5.2)
	github.com/testcontainers/testcontainers-go v0.43.0 // T2 integration harness (§5.4)
	github.com/testcontainers/testcontainers-go/modules/postgres v0.42.0 // Timescale container module
	go.etcd.io/bbolt v1.4.3 // broker session store (§3.1; pinned over pebble)
	go.mongodb.org/mongo-driver/v2 v2.6.0 // BSON raw-document API only (§3.5.5)
===
12:import (
13:	"bytes"
14:	"encoding/json"
15:	"errors"
16:	"fmt"
17:	"io"
18:	"net/http"
===
33:	// DetailBadRequest is the canonical 400 detail. Measured upstream 1.2.0
34-	// (2026-08-24, verify batch #57): Phoenix renders "Bad request" — the
35-	// capital-R form this constant carried before was reconstructed, not
36-	// observed, and every probed 400 row answers with the lowercase r.
37:	DetailBadRequest = "Bad request"
38-	// DetailUnauthorized is the canonical 401 detail.
39-	DetailUnauthorized = "Unauthorized"
40-	// DetailForbidden is the canonical 403 detail.
41-	DetailForbidden = "Forbidden"
42-	// DetailNotFound is the canonical generic 404 detail.
43-	DetailNotFound = "Not Found"
44-	// DetailDeviceNotFound is the canonical 404 detail for unknown devices
45-	// (upstream AppEngine/Pairing shape).
46-	DetailDeviceNotFound = "Device not found"
47-	// DetailRouteNotFound is the 404 detail upstream AppEngine and
48-	// RealmManagement emit for a path matching no route at all. Distinct from
49-	// DetailNotFound, which a handler emits for a route that matched and found
--
191:func WriteBadRequest(w http.ResponseWriter) error {
192:	return WriteError(w, http.StatusBadRequest, DetailBadRequest)
193-}
194-
195-// WriteUnauthorized writes the canonical 401 envelope.
196-func WriteUnauthorized(w http.ResponseWriter) error {
197-	return WriteError(w, http.StatusUnauthorized, DetailUnauthorized)
198-}
199-
200-// WriteForbidden writes the canonical 403 envelope.
201-func WriteForbidden(w http.ResponseWriter) error {
202-	return WriteError(w, http.StatusForbidden, DetailForbidden)
203-}
204-
--
228:func DecodeData(r io.Reader, maxBytes int64, dst any) error {
229-	body, err := io.ReadAll(io.LimitReader(r, maxBytes+1))
230-	if err != nil {
231-		return fmt.Errorf("astarteapi: reading request body: %w", err)
232-	}
233-	if int64(len(body)) > maxBytes {
234-		return fmt.Errorf("astarteapi: %w: exceeds %d bytes", ErrBodyTooLarge, maxBytes)
235-	}
236-	var env struct {
237-		Data json.RawMessage `json:"data"`
238-	}
239-	if err := json.Unmarshal(body, &env); err != nil {
240-		return fmt.Errorf("astarteapi: decoding request envelope: %w", err)

→ Read pkg/astarteapi/envelope.go [offset=228, limit=40]

$ cd /root/astrate-mule && rg -n "maxBodyBytes" internal/housekeeping/*.go | head -5
internal/housekeeping/http.go:15:// maxBodyBytes caps Housekeeping request bodies (a realm-create body — name,
internal/housekeeping/http.go:17:const maxBodyBytes int64 = 64 << 10
internal/housekeeping/http.go:64:	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &req); err != nil {
internal/housekeeping/http.go:164:	raw, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes+1))
internal/housekeeping/http.go:165:	if err != nil || int64(len(raw)) > maxBodyBytes {

Now the probe, replicating both handlers' decode sequences verbatim:
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ mkdir -p /tmp/hkprobe && cd /tmp/hkprobe && printf 'module hkprobe\n\ngo 1.26.1\n\nrequire github.com/astrate-platform/astrate v0.0.0\n\nreplace github.com/astrate-platform/astrate => /root/astrate-mule\n' > go.mod && ls failed
Error: The user rejected permission to use this specific tool call.
