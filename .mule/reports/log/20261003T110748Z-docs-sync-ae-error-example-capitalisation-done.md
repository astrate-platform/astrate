slug: docs-sync-ae-error-example-capitalisation
verdict: done
at:  f3bcc6d
ran: 2026-10-03T11:07:48Z on DietPi in 298s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ sed -n 1700,1810p docs/api/astarte_appengine_api.yaml
ErrorDetail:
      type: object
      required: [errors]
      properties:
        errors:
          type: object
          required: [detail]
          properties:
            detail:
              type: string

    FieldErrorsDetail:
      type: object
      required: [errors]
      properties:
        errors:
          type: object
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
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Internal Server Error

    ValidationErrors:
      description: Validation Error — one or more fields failed validation.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/FieldErrorsDetail"
          example:
            errors:
              group_name:
                - can't be blank

    ValueTooLarge:
      description: Unprocessable Entity — the value is larger than the per-value size limit (64 KiB).
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"

$ sed -n 25,70p pkg/astarteapi/envelope.go; echo ---; cat pkg/astarteapi/testdata/error_bad_request.json pkg/astarteapi/testdata/error_internal.json
// DefaultMaxBodyBytes is the request body size cap callers are expected to
// pass to DecodeData unless an endpoint has a documented reason to differ
// (interface uploads, for example, may need more than a pairing request).
const DefaultMaxBodyBytes int64 = 1 << 20 // 1 MiB

// Canonical upstream error detail strings. These are frozen: SDK and
// astartectl error paths match on them.
const (
	// DetailBadRequest is the canonical 400 detail. Measured upstream 1.2.0
	// (2026-08-24, verify batch #57): Phoenix renders "Bad request" — the
	// capital-R form this constant carried before was reconstructed, not
	// observed, and every probed 400 row answers with the lowercase r.
	DetailBadRequest = "Bad request"
	// DetailUnauthorized is the canonical 401 detail.
	DetailUnauthorized = "Unauthorized"
	// DetailForbidden is the canonical 403 detail.
	DetailForbidden = "Forbidden"
	// DetailNotFound is the canonical generic 404 detail.
	DetailNotFound = "Not Found"
	// DetailDeviceNotFound is the canonical 404 detail for unknown devices
	// (upstream AppEngine/Pairing shape).
	DetailDeviceNotFound = "Device not found"
	// DetailRouteNotFound is the 404 detail upstream AppEngine and
	// RealmManagement emit for a path matching no route at all. Distinct from
	// DetailNotFound, which a handler emits for a route that matched and found
	// no resource; the capitalisation difference is upstream's, not a typo.
	DetailRouteNotFound = "Not found"
	// DetailPageNotFound is the 404 detail upstream Pairing emits for a path
	// matching no route at all.
	DetailPageNotFound = "Page not found"
	// DetailInternalServerError is the canonical 500 detail. Measured
	// upstream 1.2.0 (2026-08-24, verify batch #57): Phoenix renders
	// "Internal server error" — the capital-S form was reconstructed.
	DetailInternalServerError = "Internal server error"
)

// ErrMissingData is wrapped by DecodeData when the request body has no
// "data" key (or it is JSON null) — upstream rejects such bodies uniformly.
var ErrMissingData = errors.New(`missing "data" key in request body`)

// ErrBodyTooLarge is wrapped by DecodeData when the request body exceeds the
// caller-supplied size cap.
var ErrBodyTooLarge = errors.New("request body too large")

// dataEnvelope is the success wrapper: {"data": ...}.
type dataEnvelope struct {
---
{"errors":{"detail":"Bad request"}}{"errors":{"detail":"Internal server error"}}

$ rg -n "ErrorDetailExamples" -A 60 internal/swagger/swagger_test.go | head -160
460:// TestPairingErrorDetailExamples pins the error-detail examples in the pairing
461-// spec's components.responses to the frozen canonical strings the wire emits
462-// (astarteapi/envelope.go), so a generated client that copies an example does
463-// not learn a detail that never goes over the wire. Phoenix renders "Bad
464-// request" and "Internal server error" — the capital-R/capital-S spellings the
465-// spec used to carry were reconstructed, not measured on upstream 1.2.0, and a
466-// startectl/SDK error path matching on the canonical string would miss them.
467-// The 404 detail is deliberately absent here: pairing answers an unmatched
468-// path with DetailPageNotFound while DetailNotFound is the generic form, and
469-// the spec's own DeviceNotFound example is the one that carries a constant.
470:func TestPairingErrorDetailExamples(t *testing.T) {
471-	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
472-	if err != nil {
473-		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
474-	}
475-	lines := strings.Split(string(b), "\n")
476-
477-	const detailPrefix = "              detail: "
478-
479-	for _, tc := range []struct {
480-		response string
481-		want     string
482-	}{
483-		{"    BadRequest:", astarteapi.DetailBadRequest},
484-		{"    Unauthorized:", astarteapi.DetailUnauthorized},
485-		{"    Forbidden:", astarteapi.DetailForbidden},
486-		{"    DeviceNotFound:", astarteapi.DetailDeviceNotFound},
487-		{"    InternalServerError:", astarteapi.DetailInternalServerError},
488-	} {
489-		block := componentBlock(t, lines, tc.response)
490-
491-		got := ""
492-		for _, l := range block {
493-			if strings.HasPrefix(l, detailPrefix) {
494-				got = strings.TrimPrefix(l, detailPrefix)
495-				break
496-			}
497-		}
498-		if got == "" {
499-			t.Errorf("response %q carries no %q example", tc.response, "detail")
500-			continue
501-		}
502-		if got != tc.want {
503-			t.Errorf("response %q example detail = %q, want the canonical %q",
504-				strings.TrimSpace(tc.response), got, tc.want)
505-		}
506-	}
507-}
508-
509:// TestRealmManagementErrorDetailExamples pins the error-detail examples in the
510-// realm management spec's components.responses to the frozen canonical strings
511-// the wire emits (astarteapi/envelope.go), the same guard the pairing spec gets
512:// from TestPairingErrorDetailExamples. The BadRequest and InternalServerError
513-// examples carried the capital-R/capital-S spellings, which upstream Phoenix
514-// never renders and no handler emits: every 400 on this surface goes through
515-// astarteapi.WriteBadRequest and every 500 through WriteInternalServerError
516-// (internal/realm/http.go), and the SDK/astartectl error paths match on the
517-// canonical strings. The Conflict example is an endpoint-specific detail with no
518-// constant behind it, so it stays free-form; the three 422s that used to share
519-// one ValidationError component are pinned by
520-// TestRealmManagement422ValidationDetails.
521:func TestRealmManagementErrorDetailExamples(t *testing.T) {
522-	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
523-	if err != nil {
524-		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
525-	}
526-	lines := strings.Split(string(b), "\n")
527-
528-	const detailPrefix = "              detail: "
529-
530-	for _, tc := range []struct {
531-		response string
532-		want     string
533-	}{
534-		{"    BadRequest:", astarteapi.DetailBadRequest},
535-		{"    Unauthorized:", astarteapi.DetailUnauthorized},
536-		{"    Forbidden:", astarteapi.DetailForbidden},
537-		{"    NotFound:", astarteapi.DetailNotFound},
538-		{"    DeviceNotFound:", astarteapi.DetailDeviceNotFound},
539-		{"    InternalServerError:", astarteapi.DetailInternalServerError},
540-	} {
541-		block := componentBlock(t, lines, tc.response)
542-
543-		got := ""
544-		for _, l := range block {
545-			if strings.HasPrefix(l, detailPrefix) {
546-				got = strings.TrimPrefix(l, detailPrefix)
547-				break
548-			}
549-		}
550-		if got == "" {
551-			t.Errorf("response %q carries no %q example", tc.response, "detail")
552-			continue
553-		}
554-		if got != tc.want {
555-			t.Errorf("response %q example detail = %q, want the canonical %q",
556-				strings.TrimSpace(tc.response), got, tc.want)
557-		}
558-	}
559-}
560-
561-// TestRealmManagement422ValidationDetails pins the 422 examples the three
562-// non-interface realm-management operations carry to the details the wire
563-// actually emits, so a generated client that copies one does not learn a
564-// message no handler can produce.
565-//
566-// The three operations used to share components.responses.ValidationError, whose
567-// single example read `invalid interface: missing "interface_name"`. That string
568-// is wrong for all three: it names an interface, and none of them installs or
569-// updates one. Each now carries its own inline 422, the pattern
570-// installInterface and updateInterface already use for their own details.
571-//
572-// The createPolicy examples are derived from the compiler here instead of being
573-// written out, so a change to the message fails here instead of drifting the
574-// spec — the coupling TestPairingDeviceIDEncodingDocumented gets from
575-// deviceid. The putAuthConfig and deletePolicy details have no exported
576-// constant behind them (Service.SetAuthKey and Service.DeletePolicy,
577-// internal/realm/service.go), so they are spelled out and cited here.
578-func TestRealmManagement422ValidationDetails(t *testing.T) {
579-	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
580-	if err != nil {
581-		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)

$ wc -l internal/swagger/swagger_test.go; sed -n 1,60p internal/swagger/swagger_test.go
947 internal/swagger/swagger_test.go
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

$ rg -n "func componentBlock" -A 25 internal/swagger/swagger_test.go
881:func componentBlock(t *testing.T, lines []string, key string) []string {
882-	t.Helper()
883-	start := -1
884-	for i, l := range lines {
885-		if l == key {
886-			start = i
887-			break
888-		}
889-	}
890-	if start < 0 {
891-		t.Fatalf("spec declares no components entry %q", key)
892-	}
893-	for i := start + 1; i < len(lines); i++ {
894-		if strings.TrimSpace(lines[i]) == "" {
895-			continue
896-		}
897-		if !strings.HasPrefix(lines[i], "      ") {
898-			return lines[start+1 : i]
899-		}
900-	}
901-	return lines[start+1:]
902-}
903-
904-func containsLine(lines []string, want string) bool {
905-	for _, l := range lines {
906-		if l == want {

$ rg -n "Bad Request|Internal Server Error|Bad request|Internal server error" docs/api/astarte_appengine_api.yaml
1725:      description: Bad Request — malformed or missing data envelope.
1732:              detail: Bad Request
1785:      description: Internal Server Error — unexpected failure.
1792:              detail: Internal Server Error

$ rg -rn "WriteBadRequest|WriteInternalServerError" internal/appengine/ | head -20
internal/appengine/http.go:		_ = astarteapi.n(w)
internal/appengine/http.go:		_ = astarteapi.n(w)
internal/appengine/http.go:		_ = astarteapi.n(w)
internal/appengine/http.go:		_ = astarteapi.n(w)
internal/appengine/http.go:		_ = astarteapi.n(w)
internal/appengine/http.go:		_ = astarteapi.n(w)
internal/appengine/http.go:		_ = astarteapi.n(w)
internal/appengine/http.go:		_ = astarteapi.n(w)
internal/appengine/http.go:		_ = astarteapi.n(w)
internal/appengine/http.go:		_ = astarteapi.n(w)
internal/appengine/http.go:		_ = astarteapi.n(w)
internal/appengine/channels/ws.go:		_ = astarteapi.n(w)

$ sed -n 1720,1740p docs/api/astarte_appengine_api.yaml; echo ===; sed -n 1780,1800p docs/api/astarte_appengine_api.yaml
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
===
          example:
            errors:
              detail: Cannot write to read-only resource

    InternalServerError:
      description: Internal Server Error — unexpected failure.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Internal Server Error

    ValidationErrors:
      description: Validation Error — one or more fields failed validation.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/FieldErrorsDetail"
          example:

$ rg -c "astarteapi\.WriteBadRequest" internal/appengine/http.go; rg -c "astarteapi\.WriteInternalServerError" internal/appengine/http.go
8
3

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -1728,9 +1728,9 @@
       schema:
         $ref: "#/components/schemas/ErrorDetail"
       example:
         errors:
-          detail: Bad Request
+          detail: Bad request
 
 Unauthorized:
   description: Unauthorized — missing or invalid JWT.
   content:

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -1788,9 +1788,9 @@
       schema:
         $ref: "#/components/schemas/ErrorDetail"
       example:
         errors:
-          detail: Internal Server Error
+          detail: Internal server error
 
 ValidationErrors:
   description: Validation Error — one or more fields failed validation.
   content:

$ rg -n "detail: (Bad|Internal)" docs/api/astarte_appengine_api.yaml
1732:              detail: Bad request
1792:              detail: Internal server error

Now the test. Let me check how the appengine spec's responses section is laid out (other detail examples) and place the new test beside the other two.
→ Read internal/swagger/swagger_test.go [offset=440, limit=70]

$ awk 'NR>1700' docs/api/astarte_appengine_api.yaml | rg -n "^    [A-Za-z]+:|detail:" | head -60
8:            detail:
11:    FieldErrorsDetail:
24:    BadRequest:
32:              detail: Bad request
34:    Unauthorized:
42:              detail: Unauthorized
44:    Forbidden:
52:              detail: Forbidden
54:    NotFound:
62:              detail: Not Found
64:    DeviceNotFound:
72:              detail: Device not found
74:    MethodNotAllowed:
82:              detail: Cannot write to read-only resource
84:    InternalServerError:
92:              detail: Internal server error
94:    ValidationErrors:
105:    ValueTooLarge:
113:              detail: Value size exceeds size limits
115:    UnprocessableEntity:
123:              detail: Attribute key not found
125:    ConflictAliasInUse:
133:              detail: Alias already in use
135:    ConflictGroupAlreadyExists:
143:              detail: Group already exists
145:    ConflictDeviceAlreadyInGroup:
153:              detail: Device already in group

→ Read internal/swagger/swagger_test.go [limit=10, offset=555]

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -557,8 +557,60 @@
 		}
 	}
 }
 
+// TestAppEngineErrorDetailExamples pins the error-detail examples in the app
+// engine spec's components.responses to the frozen canonical strings the wire
+// emits (astarteapi/envelope.go), the same guard the pairing and realm
+// management specs get from TestPairingErrorDetailExamples and
+// TestRealmManagementErrorDetailExamples. The BadRequest and InternalServerError
+// examples carried the capital-R/capital-S spellings, which upstream Phoenix
+// never renders and no handler emits: every 400 on this surface goes through
+// astarteapi.WriteBadRequest and every 500 through WriteInternalServerError
+// (internal/appengine/http.go), and the SDK/astartectl error paths match on the
+// canonical strings. The responses whose example is an endpoint-specific detail
+// with no constant behind it (MethodNotAllowed, ValueTooLarge,
+// UnprocessableEntity, the three Conflicts) stay free-form and are not pinned.
+func TestAppEngineErrorDetailExamples(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_appengine_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_appengine_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	const detailPrefix = "              detail: "
+
+	for _, tc := range []struct {
+		response string
+		want     string
+	}{
+		{"    BadRequest:", astarteapi.DetailBadRequest},
+		{"    Unauthorized:", astarteapi.DetailUnauthorized},
+		{"    Forbidden:", astarteapi.DetailForbidden},
+		{"    NotFound:", astarteapi.DetailNotFound},
+		{"    DeviceNotFound:", astarteapi.DetailDeviceNotFound},
+		{"    InternalServerError:", astarteapi.DetailInternalServerError},
+	} {
+		block := componentBlock(t, lines, tc.response)
+
+		got := ""
+		for _, l := range block {
+			if strings.HasPrefix(l, detailPrefix) {
+				got = strings.TrimPrefix(l, detailPrefix)
+				break
+			}
+		}
+		if got == "" {
+			t.Errorf("response %q carries no %q example", tc.response, "detail")
+			continue
+		}
+		if got != tc.want {
+			t.Errorf("response %q example detail = %q, want the canonical %q",
+				strings.TrimSpace(tc.response), got, tc.want)
+		}
+	}
+}
+
 // TestRealmManagement422ValidationDetails pins the 422 examples the three
 // non-interface realm-management operations carry to the details the wire
 // actually emits, so a generated client that copies one does not learn a
 // message no handler can produce.

$ gofmt -l . ; go test ./internal/swagger/ -run 'ErrorDetailExamples' -v 2>&1 | tail -20
=== RUN   TestPairingErrorDetailExamples
--- PASS: TestPairingErrorDetailExamples (0.00s)
=== RUN   TestRealmManagementErrorDetailExamples
--- PASS: TestRealmManagementErrorDetailExamples (0.00s)
=== RUN   TestAppEngineErrorDetailExamples
--- PASS: TestAppEngineErrorDetailExamples (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.014s

$ cp docs/api/astarte_appengine_api.yaml /tmp/opencode/ae.yaml.bak && sed -i '1732s/Bad request/Bad Request/; 1792s/Internal server error/Internal Server Error/' docs/api/astarte_appengine_api.yaml && go test ./internal/swagger/ -run 'TestAppEngineErrorDetailExamples' 2>&1 | tail -8; cp /tmp/opencode/ae.yaml.bak docs/api/astarte_appengine_api.yaml && rg -n "detail: (Bad|Internal)" docs/api/astarte_appengine_api.yaml
--- FAIL: TestAppEngineErrorDetailExamples (0.00s)
    swagger_test.go:607: response "BadRequest:" example detail = "Bad Request", want the canonical "Bad request"
    swagger_test.go:607: response "InternalServerError:" example detail = "Internal Server Error", want the canonical "Internal server error"
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.032s
FAIL
1732:              detail: Bad request
1792:              detail: Internal server error

Test is confirmed to fail without the change. Now the docs build and the full gate.
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

$ rg -n "site" .gitignore docs/.gitignore 2>/dev/null; ls docs/site | head; which mkdocs python3 pip3 2>/dev/null
.gitignore:29:/site-dist/
.gitignore:31:# Docs site copies (refreshed by docs/make sync)
.gitignore:32:/docs/site/DESIGN.md
.gitignore:33:/docs/site/ROADMAP.md
.gitignore:34:/docs/site/api/
.gitignore:35:/docs/site/swagger-ui/
DESIGN.md
ROADMAP.md
api
appengine-api.md
architecture.md
compatibility.md
configuration-reference.md
contributing.md
data-modeling.md
deployment.md
/usr/bin/python3

$ rg -rn "astarte_appengine_api" --glob '!docs/site/**' -l .; echo ===; rg -n "yaml" go.mod | head; echo ===; rg -n "func Specs" -A 20 internal/swagger/swagger.go
./docs/swagger-ui/index.html
./docs/handoff/phase-1-memory.md
./internal/appengine/http.go
./internal/appengine/service.go
./internal/swagger/swagger_test.go
===
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect
===
57:func Specs() []string {
58-	var names []string
59-	_ = fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
60-		if err != nil || d.IsDir() {
61-			return nil
62-		}
63-		if strings.HasSuffix(path, ".yaml") {
64-			names = append(names, strings.TrimPrefix(path, "api/"))
65-		}
66-		return nil
67-	})
68-	return names
69-}

$ rg -n "astarte_appengine_api" docs/swagger-ui/index.html internal/swagger/swagger_test.go | head
docs/swagger-ui/index.html:47:      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
internal/swagger/swagger_test.go:140:	specDocuments403(t, "astarte_appengine_api.yaml")
internal/swagger/swagger_test.go:382:	b, err := docs.APIYAML.ReadFile("api/astarte_appengine_api.yaml")
internal/swagger/swagger_test.go:384:		t.Fatalf("reading astarte_appengine_api.yaml: %v", err)
internal/swagger/swagger_test.go:426:	b, err := docs.APIYAML.ReadFile("api/astarte_appengine_api.yaml")
internal/swagger/swagger_test.go:428:		t.Fatalf("reading astarte_appengine_api.yaml: %v", err)
internal/swagger/swagger_test.go:574:	b, err := docs.APIYAML.ReadFile("api/astarte_appengine_api.yaml")
internal/swagger/swagger_test.go:576:		t.Fatalf("reading astarte_appengine_api.yaml: %v", err)

$ python3 -c "import yaml; print(yaml.__version__)" 2>&1
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; print(yaml.__version__)
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ ls "$(go env GOMODCACHE)/gopkg.in" 2>/dev/null; ls "$(go env GOMODCACHE)/go.yaml.in" 2>/dev/null
check.v1@v0.0.0-20161208181325-20d25e280405
check.v1@v1.0.0-20201130134442-10cb98267c6c
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1
yaml

← Write /tmp/opencode/yamlprobe/main.go
Wrote file successfully.

$ cat > /tmp/opencode/yamlprobe/go.mod <<'EOF'
module yamlprobe

go 1.24

require gopkg.in/yaml.v3 v3.0.1
EOF
GOFLAGS=-mod=mod GOPROXY=off go mod tidy 2>&1 | tail -3; GOPROXY=off go run . /root/astrate-mule/docs/api/astarte_appengine_api.yaml
parsed ok
346 local refs, 0 dangling
BadRequest                     Bad request
ConflictAliasInUse             Alias already in use
ConflictDeviceAlreadyInGroup   Device already in group
ConflictGroupAlreadyExists     Group already exists
DeviceNotFound                 Device not found
Forbidden                      Forbidden
InternalServerError            Internal server error
MethodNotAllowed               Cannot write to read-only resource
NotFound                       Not Found
Unauthorized                   Unauthorized
UnprocessableEntity            Attribute key not found
ValidationErrors               <nil>
ValueTooLarge                  Value size exceeds size limits

$ rg -n "detail: (Bad|Internal)" docs/site/api/astarte_appengine_api.yaml; gofmt -l .; go vet ./... 2>&1 | tail -5
1732:              detail: Bad request
1792:              detail: Internal server error

$ go test ./... 2>&1 | rg -v "^ok|no test files" | tail -20; echo "EXIT_SUMMARY_DONE"
EXIT_SUMMARY_DONE

$ go test ./internal/swagger/ -count=1 2>&1 | tail -3; rm -rf /tmp/opencode/yamlprobe /tmp/opencode/ae.yaml.bak
ok  	github.com/astrate-platform/astrate/internal/swagger	0.050s

$ git status --short
M docs/api/astarte_appengine_api.yaml
 M internal/swagger/swagger_test.go

```
Done: fixed the two wrong app-engine spec error examples (`Bad request`, `Internal server error`) and added TestAppEngineErrorDetailExamples pinning the six constant-backed responses to astarteapi's frozen strings.
Files: docs/api/astarte_appengine_api.yaml (lines 1732, 1792), internal/swagger/swagger_test.go (new TestAppEngineErrorDetailExamples beside the pairing/realm-management halves)
Verified: `go test ./internal/swagger/ -run TestAppEngineErrorDetailExamples` -> pass; same test with the two YAML lines reverted to the old spellings -> FAIL on both (proves it bites); `make -C docs build` -> the cp/sync steps ran but `mkdocs build` failed with "mkdocs: No such file or directory" (not installed here), so I substituted a throwaway YAML probe (gopkg.in/yaml.v3 from the module cache, /tmp, now deleted) on docs/api/astarte_appengine_api.yaml: parses, 346 local $refs all resolve, and both details read the canonical constants; `go vet ./...` -> pass; `gofmt -l .` -> clean; `go test ./...` (no -race, Pi gate) -> pass
Unsure: whether `make -C docs build` was ever expected to succeed in this environment — mkdocs and pip3 are both absent, so I could not run the real docs build; I verified spec validity by parsing the same file the Swagger UI references (docs/swagger-ui/index.html:47 -> ../api/astarte_appengine_api.yaml) instead. I also left the MethodNotAllowed/ValueTooLarge/UnprocessableEntity/Conflict examples free-form since they have no constant behind them; that is my reading of "change only those two", not an instruction from the task.
Follow-ups: none
```
