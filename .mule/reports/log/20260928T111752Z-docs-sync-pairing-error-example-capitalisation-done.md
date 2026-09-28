slug: docs-sync-pairing-error-example-capitalisation
verdict: done
at:  7ed78cb
ran: 2026-09-28T11:17:52Z on DietPi in 166s


> build · big-pickle

I'll start by reading the two files.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read pkg/astarteapi/envelope.go [offset=25, limit=45]

$ ls docs/ && ls docs/api/ 2>/dev/null | head -20 && which mkdocs; python3 -c "import yaml; print('pyyaml ok')" 2>&1 | tail -1
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
ModuleNotFoundError: No module named 'yaml'

→ Read docs/api/astarte_pairing_api.yaml [offset=540, limit=30]
→ Read docs/api/astarte_pairing_api.yaml [offset=620, limit=30]

$ cat pkg/astarteapi/testdata/error_bad_request.json pkg/astarteapi/testdata/error_internal.json 2>&1; echo "---"; rg -n "astarte_pairing_api|docs/api" --glob '!docs/api/*' -l | head -20
{"errors":{"detail":"Bad request"}}{"errors":{"detail":"Internal server error"}}---
docs/swagger-ui/index.html
docs/handoff/phase-1-handoff.md
docs/handoff/phase-5-handoff.md
docs/handoff/phase-7-memory.md
docs/handoff/phase-2-handoff.md
docs/handoff/phase-1-memory.md
docs/site/quickstart.md
docs/handoff/phase-4-handoff.md
docs/handoff/phase-3-memory.md
docs/handoff/phase-2-memory.md
docs/handoff/phase-6-handoff.md
docs/handoff/phase-3-handoff.md
docs/handoff/phase-4-memory.md

$ cat docs/Makefile; echo "--- mkdocs.yml refs ---"; rg -n "api/astarte" docs/mkdocs.yml | head -20
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
--- mkdocs.yml refs ---

$ cat docs/embed.go; echo "--- who embeds docs ---"; rg -n "docs\"|docs/api|astarte_pairing" --glob '*.go' | head -20; echo "--- mkdocs.yml head ---"; head -40 docs/mkdocs.yml
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
--- who embeds docs ---
pkg/interfaceschema/compat_test.go:128:					 "description": "after", "doc": "new docs"},
internal/flow/blocks/jsonpathmap_test.go:279:		t.Fatal("LookupInfo missing json_path_map docs")
internal/swagger/swagger.go:11:	docs "github.com/astrate-platform/astrate/docs"
internal/swagger/swagger_test.go:15:	docs "github.com/astrate-platform/astrate/docs"
--- mkdocs.yml head ---
site_name: Astrate
site_description: Lean, single-binary, Astarte-wire-compatible IoT platform in Go
repo_url: https://github.com/atsetilam/astrate

docs_dir: site
site_dir: ../site-dist

theme:
  name: material
  palette:
    - scheme: default
      primary: teal
      accent: amber
      toggle:
        icon: material/brightness-7
        name: Switch to dark mode
    - scheme: slate
      primary: teal
      accent: amber
      toggle:
        icon: material/brightness-4
        name: Switch to light mode
  features:
    - navigation.instant
    - navigation.sections
    - navigation.top
    - content.code.copy
    - toc.follow

markdown_extensions:
  - admonition
  - tables
  - toc:
      permalink: true
  - pymdownx.superfences
  - pymdownx.tabbed:
      alternate_style: true
  - pymdownx.highlight:
      anchor_linenums: true
  - pymdownx.inlinehilite

$ rg -n "api" docs/mkdocs.yml | head -20; echo "--- requirements ---"; cat docs/requirements.txt
50:  - AppEngine API: appengine-api.md
51:  - Realm Management API: realm-management-api.md
52:  - Housekeeping API: housekeeping-api.md
--- requirements ---
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10

→ Read internal/swagger/swagger_test.go
→ Read docs/api/astarte_pairing_api.yaml [offset=566, limit=56]

$ rg -n "^ +detail: " docs/api/astarte_pairing_api.yaml | cat -A | sed 's/\$$//' | head
59:                  detail: Not Found
68:                  detail: unhealthy
556:              detail: Bad Request
566:              detail: Unauthorized
576:              detail: Forbidden
586:              detail: Device not found
628:              detail: Too Many Requests
638:              detail: Internal Server Error

← Edit docs/api/astarte_pairing_api.yaml
Index: /root/astrate-mule/docs/api/astarte_pairing_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_pairing_api.yaml
+++ /root/astrate-mule/docs/api/astarte_pairing_api.yaml
@@ -552,9 +552,9 @@
       schema:
         $ref: "#/components/schemas/ErrorDetail"
       example:
         errors:
-          detail: Bad Request
+          detail: Bad request
 
 Unauthorized:
   description: Unauthorized — missing or invalid JWT/Bearer token.
   content:

← Edit docs/api/astarte_pairing_api.yaml
Index: /root/astrate-mule/docs/api/astarte_pairing_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_pairing_api.yaml
+++ /root/astrate-mule/docs/api/astarte_pairing_api.yaml
@@ -634,5 +634,5 @@
 schema:
   $ref: "#/components/schemas/ErrorDetail"
 example:
   errors:
-    detail: Internal Server Error
+    detail: Internal server error

Now the test that pins these to the frozen constants:
← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -451,8 +451,57 @@
 		t.Error("components.responses defines no BadRequest response")
 	}
 }
 
+// TestPairingErrorDetailExamples pins the error-detail examples in the pairing
+// spec's components.responses to the frozen canonical strings the wire emits
+// (astarteapi/envelope.go), so a generated client that copies an example does
+// not learn a detail that never goes over the wire. Phoenix renders "Bad
+// request" and "Internal server error" — the capital-R/capital-S spellings the
+// spec used to carry were reconstructed, not measured on upstream 1.2.0, and a
+//startectl/SDK error path matching on the canonical string would miss them.
+// The 404 detail is deliberately absent here: pairing answers an unmatched
+// path with DetailPageNotFound while DetailNotFound is the generic form, and
+// the spec's own DeviceNotFound example is the one that carries a constant.
+func TestPairingErrorDetailExamples(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
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
 // operationBlock returns the lines of the operation with the given operationId,
 // up to the next operation, path, or the components section.
 func operationBlock(t *testing.T, lines []string, operationID string) []string {
 	t.Helper()

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -11,8 +11,9 @@
 	"sort"
 	"strings"
 	"testing"
 
+	"github.com/astrate-platform/astrate/pkg/astarteapi"
 	docs "github.com/astrate-platform/astrate/docs"
 )
 
 func TestMount(t *testing.T) {

$ gofmt -l . ; echo "gofmt done"; go vet ./internal/swagger/ && go test ./internal/swagger/ -run TestPairingErrorDetailExamples -v 2>&1 | tail -15
internal/swagger/swagger_test.go
gofmt done
=== RUN   TestPairingErrorDetailExamples
--- PASS: TestPairingErrorDetailExamples (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.023s

$ gofmt -w internal/swagger/swagger_test.go && gofmt -l . ; echo "gofmt clean"; sed -n '1,20p' internal/swagger/swagger_test.go
gofmt clean
package swagger

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	docs "github.com/astrate-platform/astrate/docs"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

func TestMount(t *testing.T) {
	mux := http.NewServeMux()

Now verify the test actually fails without the fix, and run the full gate:
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cp docs/api/astarte_pairing_api.yaml /tmp/fixed.yaml && sed -i '556s/Bad request/Bad Request/; 638s/Internal server error/Internal Server Error/' docs/api/astarte_pairing_api.yaml && go test ./internal/swagger/ -run TestPairingErrorDetailExamples 2>&1 | tail -6; cp /tmp/fixed.yaml docs/api/astarte_pairing_api.yaml; echo "--- restored ---"; sed -n '556p;638p' docs/api/astarte_pairing_api.yaml failed
Error: The user rejected permission to use this specific tool call.
