slug: docs-sync-hk-patch-422-field-error-examples
verdict: blocked
at:  2c5d007
ran: 2026-10-04T17:17:09Z on DietPi in 1200s


> build · big-pickle

I'll start by reading the required context files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ ls /root/astrate-mule/docs/ && ls /root/astrate-mule/docs/api/ && cat /root/astrate-mule/docs/Makefile 2>/dev/null | head -40
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

→ Read internal/housekeeping/http.go [offset=150, limit=70]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=120, offset=150]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=80, offset=395]

$ rg -n "RealmPatch|^    [A-Za-z]+:|FieldErrorsDetail|ErrorDetail" /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml | head -60
13:    name: "Apache-2.0 with SECO Exclusion"
14:    url: https://www.apache.org/licenses/LICENSE-2.0
18:    variables:
27:    get:
50:    post:
109:    get:
142:    patch:
169:                  $ref: "#/components/schemas/RealmPatch"
194:    delete:
231:    RealmName:
248:    AsyncOperation:
265:    DataEnvelopeRealm:
272:    DataEnvelopeRealmNames:
282:    Realm:
303:    RealmCreate:
341:    RealmPatch:
368:    ErrorDetail:
382:    FieldErrorsDetail:
395:    BadRequest:
400:            $ref: "#/components/schemas/ErrorDetail"
405:    Unauthorized:
410:            $ref: "#/components/schemas/ErrorDetail"
415:    Forbidden:
420:            $ref: "#/components/schemas/ErrorDetail"
425:    NotFound:
430:            $ref: "#/components/schemas/ErrorDetail"
435:    MethodNotAllowed:
440:            $ref: "#/components/schemas/ErrorDetail"
445:    ValidationErrors:
450:            $ref: "#/components/schemas/FieldErrorsDetail"
456:    UnprocessableDeletion:
461:            $ref: "#/components/schemas/FieldErrorsDetail"
467:    Conflict:
472:            $ref: "#/components/schemas/ErrorDetail"
477:    ValidationError:
487:            $ref: "#/components/schemas/ErrorDetail"
500:    InternalServerError:
505:            $ref: "#/components/schemas/ErrorDetail"

→ Read internal/housekeeping/service.go [limit=100, offset=125]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=95, offset=300]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=60, offset=50]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=35, offset=474]

$ cd /root/astrate-mule && rg -n "func WriteFieldErrors|func WriteValidationError" -A 18 internal/astarteapi/*.go | head -50
rg: internal/astarteapi/*.go: IO error for operation on internal/astarteapi/*.go: No such file or directory (os error 2)

$ cd /root/astrate-mule && rg -ln "invalid_update_parameters" --glob '*_test.go' | head; echo ---; rg -n "jwt_public_key_pem.*can't be blank|device_registration_limit.*is invalid" --glob '*_test.go' | head -20
internal/housekeeping/http_test.go
---
internal/swagger/swagger_test.go:794:		{"putAuthConfig", []string{`jwt_public_key_pem can't be blank`}},
internal/housekeeping/http_test.go:288:			`{"errors":{"device_registration_limit":["is invalid"]}}`)
internal/housekeeping/http_test.go:290:			`{"errors":{"jwt_public_key_pem":["can't be blank"]}}`)

$ cd /root/astrate-mule && rg -n "func WriteFieldErrors" -A 12 --glob '*.go' | head -30
pkg/astarteapi/envelope.go:172:func WriteFieldErrors(w http.ResponseWriter, status int, fields map[string][]string) error {
pkg/astarteapi/envelope.go-173-	return write(w, status, fieldsEnvelope{Errors: fields})
pkg/astarteapi/envelope.go-174-}
pkg/astarteapi/envelope.go-175-
pkg/astarteapi/envelope.go-176-// rawEnvelope is the caller-shaped error wrapper: {"errors": <anything>}.
pkg/astarteapi/envelope.go-177-type rawEnvelope struct {
pkg/astarteapi/envelope.go-178-	Errors any `json:"errors"`
pkg/astarteapi/envelope.go-179-}
pkg/astarteapi/envelope.go-180-
pkg/astarteapi/envelope.go-181-// WriteRawErrors writes {"errors": <anything>} — caller-shaped error bodies
pkg/astarteapi/envelope.go-182-// for the nested changeset envelopes upstream uses inside sub-objects (for
pkg/astarteapi/envelope.go-183-// example a trigger's {"action": {"http_url": [...]}}). Map keys are emitted
pkg/astarteapi/envelope.go-184-// in sorted order, without HTML escaping, and without a trailing newline,

→ Read internal/housekeeping/http_test.go [limit=30, offset=270]

$ cd /root/astrate-mule && rg -n "type optString|func \(o \*optString\)" -A 22 internal/housekeeping/http.go | head -60
105:type optString struct {
106-	present bool
107-	val     string
108-}
109-
110:func (o *optString) UnmarshalJSON(b []byte) error {
111-	o.present = true
112-	return json.Unmarshal(b, &o.val)
113-}
114-
115-// optInt32 tracks presence and nullness of an optional integer field.
116-type optInt32 struct {
117-	present bool
118-	null    bool
119-	val     int32
120-}
121-
122-func (o *optInt32) UnmarshalJSON(b []byte) error {
123-	o.present = true
124-	if bytes.Equal(b, []byte("null")) {
125-		o.null = true
126-		return nil
127-	}
128-	return json.Unmarshal(b, &o.val)
129-}
130-
131-// optInt64 tracks presence and nullness of an optional integer field.
132-type optInt64 struct {

$ cd /root/astrate-mule && rg -n "func wantFieldErrors|func patch" -B4 -A18 internal/housekeeping/http_test.go | head -70
(no output)

$ cd /root/astrate-mule && ls internal/swagger/ && rg -n "housekeeping" internal/swagger/*.go | head -20; echo ---; cat docs/embed.go
swagger.go
swagger_test.go
internal/swagger/swagger_test.go:128:// upstream-parity the housekeeping and pairing specs already document.
internal/swagger/swagger_test.go:197:// TestHousekeepingAsyncOperationParamDocumented guards that the housekeeping
internal/swagger/swagger_test.go:203:// TestHousekeepingAsyncOperationParam in internal/housekeeping.
internal/swagger/swagger_test.go:205:	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
internal/swagger/swagger_test.go:207:		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
internal/swagger/swagger_test.go:310:// TestHousekeepingRetentionZeroFoldDocumented guards that the housekeeping spec
internal/swagger/swagger_test.go:313:// `null || val == 0` to ClearRetention (internal/housekeeping/http.go) and
internal/swagger/swagger_test.go:315:// (internal/housekeeping/service.go) — upstream parity measured on v1.2.0. A
internal/swagger/swagger_test.go:321:// in internal/housekeeping.
internal/swagger/swagger_test.go:323:	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
internal/swagger/swagger_test.go:325:		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
internal/swagger/swagger_test.go:667:// housekeeping spec's components.responses to the frozen canonical strings the
internal/swagger/swagger_test.go:676:// and the 500 through WriteInternalServerError (internal/housekeeping/http.go:
internal/swagger/swagger_test.go:696:	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
internal/swagger/swagger_test.go:698:		t.Fatalf("reading astrate_housekeeping_api.yaml: %v", err)
internal/swagger/swagger_test.go:982:// TestHousekeepingRealmNamePatternDocumented guards that the housekeeping spec
internal/swagger/swagger_test.go:989:// ErrValidation `realm_name is invalid` (internal/housekeeping/service.go) and
internal/swagger/swagger_test.go:991:// (internal/housekeeping/http.go). Both `RealmCreate.realm_name` and the
internal/swagger/swagger_test.go:1005:	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
internal/swagger/swagger_test.go:1007:		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
---
// Package docs embeds the Swagger UI static files and OpenAPI YAML specs so
// they are served directly from the astrate binary without external files.
package docs

import "embed"

// SwaggerUI holds the static Swagger UI files (index.html, CSS, JS refs).
//
//go:embed swagger-ui/*
var SwaggerUI embed.FS

// APIYAML holds the OpenAPI 3.0 YAML specifications for all five API surfaces.
//
//go:embed api/*.yaml
var APIYAML embed.FS

$ cd /root/astrate-mule && rg -n "wantFieldErrors" internal/housekeeping/*_test.go | head
internal/housekeeping/http_test.go:215:	wantFieldErrors := func(t *testing.T, rawBody, want string) {
internal/housekeeping/http_test.go:285:		wantFieldErrors(t, `{"datastream_maximum_storage_retention":-5}`,
internal/housekeeping/http_test.go:287:		wantFieldErrors(t, `{"device_registration_limit":-1}`,
internal/housekeeping/http_test.go:289:		wantFieldErrors(t, `{"jwt_public_key_pem":""}`,
internal/housekeeping/http_test.go:291:		wantFieldErrors(t, `{"replication_factor":3}`,

→ Read internal/swagger/swagger_test.go [limit=80, offset=305]
→ Read internal/swagger/swagger_test.go [limit=130, offset=175]

$ cd /root/astrate-mule && rg -n "^func (containsLine|componentBlock|operationBlock)" -A 30 internal/swagger/swagger_test.go | head -80
1143:func operationBlock(t *testing.T, lines []string, operationID string) []string {
1144-	t.Helper()
1145-	marker := "      operationId: " + operationID
1146-	start := -1
1147-	for i, l := range lines {
1148-		if l == marker {
1149-			start = i
1150-			break
1151-		}
1152-	}
1153-	if start < 0 {
1154-		t.Fatalf("spec declares no operation %q", operationID)
1155-	}
1156-	endRe := regexp.MustCompile(`^      operationId: |^  \S|^components:`)
1157-	for i := start + 1; i < len(lines); i++ {
1158-		if endRe.MatchString(lines[i]) {
1159-			return lines[start+1 : i]
1160-		}
1161-	}
1162-	return lines[start+1:]
1163-}
1164-
1165-// componentBlock returns the lines of the components entry starting with the
1166-// given key line, up to the next entry at the same or shallower indentation.
1167:func componentBlock(t *testing.T, lines []string, key string) []string {
1168-	t.Helper()
1169-	start := -1
1170-	for i, l := range lines {
1171-		if l == key {
1172-			start = i
1173-			break
1174-		}
1175-	}
1176-	if start < 0 {
1177-		t.Fatalf("spec declares no components entry %q", key)
1178-	}
1179-	for i := start + 1; i < len(lines); i++ {
1180-		if strings.TrimSpace(lines[i]) == "" {
1181-			continue
1182-		}
1183-		if !strings.HasPrefix(lines[i], "      ") {
1184-			return lines[start+1 : i]
1185-		}
1186-	}
1187-	return lines[start+1:]
1188-}
1189-
1190:func containsLine(lines []string, want string) bool {
1191-	for _, l := range lines {
1192-		if l == want {
1193-			return true
1194-		}
1195-	}
1196-	return false
1197-}
1198-
1199-func embeddedYAMLFilenames() ([]string, error) {
1200-	var names []string
1201-	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
1202-		if err != nil {
1203-			return err
1204-		}
1205-		if d.IsDir() {
1206-			return nil
1207-		}
1208-		if strings.HasSuffix(path, ".yaml") {
1209-			names = append(names, strings.TrimPrefix(path, "api/"))
1210-		}
1211-		return nil
1212-	})
1213-	sort.Strings(names)
1214-	return names, err
1215-}
1216-
1217-func get(t *testing.T, url string) string {
1218-	t.Helper()
1219-	resp, err := http.Get(url)
1220-	if err != nil {

→ Read internal/housekeeping/http_test.go [limit=70, offset=175]

← Write internal/housekeeping/zz_probe_test.go
Wrote file successfully.

$ cd /root/astrate-mule && go test ./internal/housekeeping/ -run 'TestZZProbe' -v 2>&1 | head -60
# github.com/astrate-platform/astrate/internal/housekeeping [github.com/astrate-platform/astrate/internal/housekeeping.test]
internal/housekeeping/zz_probe_test.go:9:7: undefined: newHKRig
internal/housekeeping/zz_probe_test.go:10:22: undefined: randSuffix
internal/housekeeping/zz_probe_test.go:11:76: undefined: jsonStr
internal/housekeeping/zz_probe_test.go:34:7: undefined: newHKRig
internal/housekeeping/zz_probe_test.go:36:26: undefined: randSuffix
internal/housekeeping/zz_probe_test.go:36:70: undefined: jsonStr
internal/housekeeping/zz_probe_test.go:37:26: undefined: randSuffix
internal/housekeeping/zz_probe_test.go:37:70: undefined: jsonStr
internal/housekeeping/zz_probe_test.go:38:26: undefined: randSuffix
internal/housekeeping/zz_probe_test.go:39:26: undefined: randSuffix
internal/housekeeping/zz_probe_test.go:39:26: too many errors
FAIL	github.com/astrate-platform/astrate/internal/housekeeping [build failed]
FAIL

$ cd /root/astrate-mule && head -3 internal/housekeeping/http_test.go
//go:build integration

package housekeeping

$ rm /root/astrate-mule/internal/housekeeping/zz_probe_test.go && mkdir -p /tmp/opencode/nullprobe && cat > /tmp/opencode/nullprobe/main.go <<'EOF'
package main

import (
	"encoding/json"
	"fmt"
)

type optString struct {
	present bool
	val     string
}

func (o *optString) UnmarshalJSON(b []byte) error {
	o.present = true
	return json.Unmarshal(b, &o.val)
}

type pb struct {
	PEM optString `json:"jwt_public_key_pem"`
}

func main() {
	for _, raw := range []string{`{"jwt_public_key_pem":null}`, `{"jwt_public_key_pem":""}`, `{}`} {
		var p pb
		err := json.Unmarshal([]byte(raw), &p)
		fmt.Printf("%-32s err=%v present=%v val=%q -> 422-blank=%v\n",
			raw, err, p.PEM.present, p.PEM.val, p.PEM.present && p.PEM.val == "")
	}
}
EOF
cd /tmp/opencode/nullprobe && go mod init nullprobe >/dev/null 2>&1; go run main.go
{"jwt_public_key_pem":null}      err=<nil> present=true val="" -> 422-blank=true
{"jwt_public_key_pem":""}        err=<nil> present=true val="" -> 422-blank=true
{}                               err=<nil> present=false val="" -> 422-blank=false

$ cd /root/astrate-mule && rg -n "UpdateRealm\(" --glob '*.go' | grep -v "_test.go" | head -20; echo ---; rg -n "ErrValidation" --glob '*.go' -l | grep -v _test | head
internal/store/realms.go:157:func (s *Store) UpdateRealm(ctx context.Context, name string, p RealmPatch) error {
internal/housekeeping/http.go:225:	if err := a.svc.UpdateRealm(r.Context(), r.PathValue("realm"), u); err != nil {
internal/housekeeping/service.go:57:	UpdateRealm(ctx context.Context, name string, p store.RealmPatch) error
internal/housekeeping/service.go:204:func (s *Service) UpdateRealm(ctx context.Context, name string, p RealmUpdate) error {
internal/housekeeping/service.go:214:	return s.st.UpdateRealm(ctx, name, store.RealmPatch{
---
internal/housekeeping/http.go
internal/housekeeping/service.go
internal/appengine/data.go
internal/appengine/http.go
internal/appengine/service.go
internal/flowapi/http.go
internal/flowapi/userblocks.go
internal/realm/http.go
internal/flowapi/service.go
internal/realm/service.go

$ cd /root/astrate-mule && rg -c "internal/" docs/api/*.yaml; echo ---; rg -n "internal/" docs/api/astarte_housekeeping_api.yaml | head
docs/api/astarte_realm_management_api.yaml:2
---

$ cd /root/astrate-mule && rg -n -B6 -A3 "internal/" docs/api/astarte_realm_management_api.yaml | head -40; echo === ; rg -n "COMPATIBILITY.md" docs/api/*.yaml | head
572-          $ref: "#/components/responses/Forbidden"
573-        "422":
574-          description: |
575-            Unprocessable Entity — the auth configuration failed validation.
576-            A body with an empty `jwt_public_key_pem` is answered 422, not 400:
577-            it parses, and SetAuthKey rejects a blank key (Service.SetAuthKey,
578:            internal/realm/service.go) as an ErrValidation detail.
579-          content:
580-            application/json; charset=utf-8:
581-              schema:
--
868-          $ref: "#/components/responses/NotFound"
869-        "422":
870-          description: |
871-            Unprocessable Entity — the policy is still referenced by a trigger.
872-            Deletion is refused while any trigger in the realm names this
873-            policy; the detail names both (Service.DeletePolicy,
874:            internal/realm/service.go).
875-          content:
876-            application/json; charset=utf-8:
877-              schema:
===
docs/api/astarte_housekeeping_api.yaml:257:        response is identical either way (deviation 17 in COMPATIBILITY.md).
docs/api/astarte_realm_management_api.yaml:700:        deviation from upstream's async deletion (see COMPATIBILITY.md).
docs/api/astarte_realm_management_api.yaml:966:        (deviation 17 in COMPATIBILITY.md).

$ cd /root/astrate-mule && rg -n "ValidationError|examples" internal/swagger/swagger_test.go | head -20; echo ===; rg -n "yaml" go.mod; echo ===; sed -n 1,60p internal/swagger/swagger_test.go
378:// $ref the ValidationErrors response — pointing it at BadRequest would teach a
391:		ref = `          $ref: "#/components/responses/ValidationErrors"`
406:		t.Errorf("addGroupDevice 422 is followed by %q, want the ValidationErrors $ref", block[idx+1])
409:	comp := strings.Join(componentBlock(t, lines, "    ValidationErrors:"), "\n")
411:		t.Error("components.responses.ValidationErrors does not carry the FieldErrorsDetail body")
461:// TestPairingErrorDetailExamples pins the error-detail examples in the pairing
510:// TestRealmManagementErrorDetailExamples pins the error-detail examples in the
514:// examples carried the capital-R/capital-S spellings, which upstream Phoenix
520:// one ValidationError component are pinned by
562:// TestAppEngineErrorDetailExamples pins the error-detail examples in the app
567:// examples carried the capital-R/capital-S spellings, which upstream Phoenix
614:// TestNativeErrorDetailExamples pins the error-detail examples in the native
666:// TestHousekeepingErrorDetailExamples pins the error-detail examples in the
673:// and InternalServerError examples read "Bad Request" and "Internal Server
680:// (MethodNotAllowed, Conflict, ValidationError) stay free-form and are not
749:// TestRealmManagement422ValidationDetails pins the 422 examples the three
754:// The three operations used to share components.responses.ValidationError, whose
760:// The createPolicy examples are derived from the compiler here instead of being
773:	if strings.Contains(string(b), "#/components/responses/ValidationError") {
774:		t.Error("spec still $refs a shared ValidationError response; one example cannot be right for every surface that refs it")
===
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect
===
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
	"github.com/astrate-platform/astrate/internal/engine/triggers"
	"github.com/astrate-platform/astrate/migrations"
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

$ cd /root/astrate-mule && sed -n 650,720p internal/swagger/swagger_test.go | head -90
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

// TestHousekeepingErrorDetailExamples pins the error-detail examples in the
// housekeeping spec's components.responses to the frozen canonical strings the
// wire emits (astarteapi/envelope.go), the same guard the pairing, realm
// management, app engine and native specs get from
// TestPairingErrorDetailExamples, TestRealmManagementErrorDetailExamples,
// TestAppEngineErrorDetailExamples and TestNativeErrorDetailExamples. This is
// the last spec still carrying a reconstructed capitalisation: its BadRequest
// and InternalServerError examples read "Bad Request" and "Internal Server
// Error", spellings no upstream Phoenix render answers with and no handler
// emits — every 400 on this surface goes through astarteapi.WriteBadRequest
// and the 500 through WriteInternalServerError (internal/housekeeping/http.go:
// the writeError switch), and the SDK/astartectl error paths match on the
// canonical strings. The 401 and 403 were already canonical. The responses
// whose example is an endpoint-specific detail with no constant behind it
// (MethodNotAllowed, Conflict, ValidationError) stay free-form and are not
// pinned.
//
// There is no DeviceNotFound pin here, unlike the app engine and realm
// management specs: this surface never answers a device 404, so
// components.responses declares no such component to pin.
//
// The schema-level ErrorDetail example is pinned too, to the same canonical
// string as BadRequest. It sits on the generic schema all eight error
// responses $ref, so it is the example a generator falls back to for any
// response whose own it cannot resolve — leaving it on the reconstructed
// "Bad Request" taught a 404 and a 500 a detail no handler emits. Its shape is
// already the wire's: the example belongs to the `errors` sub-schema, so it is
// the value of the `errors` property and not a bare envelope, and it must not
// be re-wrapped in another `errors:`.
func TestHousekeepingErrorDetailExamples(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
	if err != nil {
		t.Fatalf("reading astrate_housekeeping_api.yaml: %v", err)
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
		{"    NotFound:", astarteapi.DetailNotFound},
		{"    InternalServerError:", astarteapi.DetailInternalServerError},
	} {
		block := componentBlock(t, lines, tc.response)

		got := ""
		for _, l := range block {
			if strings.HasPrefix(l, detailPrefix) {
				got = strings.TrimPrefix(l, detailPrefix)
				break

$ cd /root/astrate-mule && python3 -c "
# Look up the constants
import sys
sys.path.insert(0, '.')
" 2>&1 | tail -2; echo "---"; grep -n "DetailBadRequest" pkg/astarteapi/* 2>&1 | head -5
---
pkg/astarteapi/envelope.go:33:	// DetailBadRequest is the canonical 400 detail. Measured upstream 1.2.0
pkg/astarteapi/envelope.go:37:	DetailBadRequest = "Bad request"
pkg/astarteapi/envelope.go:192:	return WriteError(w, http.StatusBadRequest, DetailBadRequest)
grep: pkg/astarteapi/testdata: Is a directory

$ cd /root/astrate-mule && python3 -c "
# find bad request capitalization
pass
" 2>&1 | tail -0; echo "---"
grep -rn '"Bad Request"' docs/api/astarte_housekeeping_api.yaml 2>&1 | head -5; echo "---"
grep -rn '"Internal Server Error"' docs/api/astarte_housekeeping_api.yaml 2>&1 | head -5
---
---

Let us look at the ValidationErrors response again and also look at ValidationError component's style for examples.
→ Read docs/api/astarte_housekeeping_api.yaml [limit=15, offset=445]

$ sed -n 440,460p /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml | head -15
$ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Realm deletion disabled

    ValidationErrors:
      description: Unprocessable Entity — update parameters failed validation.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/FieldErrorsDetail"
          example:
            errors:
              error_name:
                - invalid_update_parameters

I'll update the ValidationErrors response to include all four field error cases for PATCH.
← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -442,17 +442,42 @@
         errors:
           detail: Realm deletion disabled
 
 ValidationErrors:
-  description: Unprocessable Entity — update parameters failed validation.
+  description: |
+    Unprocessable Entity — update parameters failed validation. PATCH responses
+    use the FieldErrors shape with field-specific messages that differ from POST
+    create for the same numeric conditions (PATCH says `is invalid`, while POST
+    rejects a negative limit/retention with a flat `ErrorDetail` message).
   content:
     application/json; charset=utf-8:
       schema:
         $ref: "#/components/schemas/FieldErrorsDetail"
-      example:
-        errors:
-          error_name:
-            - invalid_update_parameters
+      examples:
+        unknown_field:
+          summary: Unknown field — rejected by PATCH allowed fields
+          value:
+            errors:
+              error_name:
+                - invalid_update_parameters
+        blank_jwt:
+          summary: Empty jwt_public_key_pem
+          value:
+            errors:
+              jwt_public_key_pem:
+                - "can't be blank"
+        invalid_limit:
+          summary: Negative device_registration_limit
+          value:
+            errors:
+              device_registration_limit:
+                - is invalid
+        invalid_retention:
+          summary: Negative datastream_maximum_storage_retention
+          value:
+            errors:
+              datastream_maximum_storage_retention:
+                - is invalid
 
 UnprocessableDeletion:
   description: Unprocessable Entity — realm has devices still connected and cannot be deleted.
   content:

$ python3 -c "import yaml; yaml.safe_load(open('/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml'))" 2>&1 | tail -30
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; yaml.safe_load(open('/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml'))
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ pip show pyyaml 2>&1 | tail -2; python3 -c "import sys; print(sys.path)" 2>&1 | tail -1
/bin/sh: 1: pip: not found
['', '/usr/lib/python313.zip', '/usr/lib/python3.13', '/usr/lib/python3.13/lib-dynload', '/usr/local/lib/python3.13/dist-packages', '/usr/lib/python3/dist-packages']

$ python3 -m yaml 2>&1 | tail -1; python3 -c "
# try to load with json-like? no. just check for syntax issues by parsing structure
pass
" 2>&1 | tail -0
/usr/bin/python3: No module named yaml

$ python3 -c "
# basic yaml validation: count colons
with open('/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml') as f:
    s = f.read()
# quick check
pass
" 2>&1 | tail -0; echo "---"
sed -n 445,485p /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml | head -40
---
    ValidationErrors:
      description: |
        Unprocessable Entity — update parameters failed validation. PATCH responses
        use the FieldErrors shape with field-specific messages that differ from POST
        create for the same numeric conditions (PATCH says `is invalid`, while POST
        rejects a negative limit/retention with a flat `ErrorDetail` message).
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/FieldErrorsDetail"
          examples:
            unknown_field:
              summary: Unknown field — rejected by PATCH allowed fields
              value:
                errors:
                  error_name:
                    - invalid_update_parameters
            blank_jwt:
              summary: Empty jwt_public_key_pem
              value:
                errors:
                  jwt_public_key_pem:
                    - "can't be blank"
            invalid_limit:
              summary: Negative device_registration_limit
              value:
                errors:
                  device_registration_limit:
                    - is invalid
            invalid_retention:
              summary: Negative datastream_maximum_storage_retention
              value:
                errors:
                  datastream_maximum_storage_retention:
                    - is invalid

    UnprocessableDeletion:
      description: Unprocessable Entity — realm has devices still connected and cannot be deleted.
      content:
        application/json; charset=utf-8:

$ cd /root/astrate-mule && python3 -c "import sys; sys.path.insert(0, ''); import yaml; print('ok')" 2>&1 | tail -5; python3 -c "import yaml; yaml.safe_load(open('docs/api/astarte_housekeeping_api.yaml'))" 2>&1 | tail -15
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import sys; sys.path.insert(0, ''); import yaml; print('ok')
                                        ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; yaml.safe_load(open('docs/api/astarte_housekeeping_api.yaml'))
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ find /usr/lib -name '*yaml*' -not -path '*site-packages*' 2>&1 | tail -10; find / -name 'PyYAML*' 2>&1 | tail -10
/usr/lib/aarch64-linux-gnu/libyaml-0.so.2
/usr/lib/aarch64-linux-gnu/libyaml-0.so.2.0.9
find: ‘/proc/2824790/task/2824790/net’: Invalid argument
find: ‘/proc/2824790/net’: Invalid argument

$ cd /root/astrate-mule && go run ./tools/... 2>&1 | tail -5
go: warning: "./tools/..." matched no packages
go: no packages loaded from ./tools/...

$ ls /root/astrate-mule/tools 2>&1 | tail -20; echo "---"; find /root/astrate-mule -name 'make*' -path '*/docs*' 2>&1 | tail -10
bsoncapture
mule-plan-day.sh
mule-setup-pi.sh
mule-survey.sh
mule.sh
reconcile.sh
---

$ cd /root/astrate-mule/docs && pip install -r requirements.txt 2>&1 | tail -30
/bin/sh: 1: pip: not found

$ python3 -m pip --version 2>&1 | tail -5
/usr/bin/python3: No module named pip

