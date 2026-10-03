slug: docs-sync-hk-error-example-capitalisation
verdict: blocked
at:  4e37cc3
ran: 2026-10-03T11:30:21Z on DietPi in 1200s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

$ cat /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml | head -80
openapi: "3.0.3"
info:
  title: Astarte Housekeeping API
  description: |
    Instance-level administration API for Astrate. Manages realms — create, list,
    inspect, and delete them. All endpoints are guarded by an instance-level JWT
    carrying the `a_ha` claim.

    This API surface mirrors the upstream Astarte Housekeeping API paths and
    envelope conventions for wire compatibility.
  version: "0.0.1"
  license:
    name: "Apache-2.0 with SECO Exclusion"
    url: https://www.apache.org/licenses/LICENSE-2.0

servers:
  - url: "{base_url}"
    variables:
      base_url:
        default: http://localhost:8080
        description: Astrate base URL

security: []

paths:
  /housekeeping/v1/realms:
    get:
      operationId: listRealms
      summary: List all realms
      description: Returns the names of all configured realms, wrapped in a data envelope.
      tags: [Housekeeping]
      security:
        - a_ha: []
      responses:
        "200":
          description: A list of realm names.
          content:
            application/json; charset=utf-8:
              schema:
                $ref: "#/components/schemas/DataEnvelopeRealmNames"
              example:
                {"data":["test","production"]}
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "500":
          $ref: "#/components/responses/InternalServerError"

    post:
      operationId: createRealm
      summary: Create a new realm
      description: |
        Creates a new realm. The request body is wrapped in a data envelope and
        contains the realm name, the JWT public key PEM, and an optional device
        registration limit.
      tags: [Housekeeping]
      security:
        - a_ha: []
      parameters:
        - $ref: "#/components/parameters/AsyncOperation"
      requestBody:
        required: true
        content:
          application/json; charset=utf-8:
            schema:
              type: object
              required: [data]
              properties:
                data:
                  $ref: "#/components/schemas/RealmCreate"
            example:
              data:
                realm_name: test
                jwt_public_key_pem: |
                  -----BEGIN PUBLIC KEY-----
                  MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
                  -----END PUBLIC KEY-----
                device_registration_limit: 100
      responses:

$ sed -n '340,390p' /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
format: int64
          nullable: true
          description: |
            Maximum datastream storage retention in seconds. Null clears the
            retention. `0` clears it as well: 0 is folded to unset, so it cannot
            ask for a literal zero-second retention and the realm is left
            unlimited (upstream parity, measured on Astarte v1.2.0).

    # ── Error envelopes ──

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
              description: Human-readable error detail string.
          example:
            detail: Bad Request

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

$ sed -n '450,490p' /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
description: Conflict — resource already exists (e.g. realm name taken).
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Realm already exists

    ValidationError:
      description: Unprocessable Entity — request body failed validation.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: "realm_name can't be blank"

    InternalServerError:
      description: Internal Server Error — unexpected failure.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Internal Server Error

$ sed -n '1,80p' /root/astrate-mule/pkg/astarteapi/envelope.go
// Package astarteapi implements the JSON envelope conventions shared by every
// Astarte-compatible REST surface: success bodies are wrapped as
// {"data": ...}, error bodies as {"errors": {"detail": "..."}} (or the
// field-keyed changeset shape used by 422 validation failures), and request
// bodies arrive wrapped as {"data": ...}.
//
// These exact bytes are parsed by astartectl and the official device SDKs, so
// the golden fixtures under testdata/ are wire-frozen: changing any envelope
// produced here is a compatibility break, not a refactor.
package astarteapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ContentType is the Content-Type header value set on every envelope,
// matching what upstream Astarte's Phoenix endpoints emit.
const ContentType = "application/json; charset=utf-8"

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
	Data any `json:"data"`
}

// detailEnvelope is the error wrapper: {"errors": {"detail": "..."}}.
type detailEnvelope struct {
	Errors detailBody `json:"errors"`
}

type detailBody struct {
	Detail string `json:"detail"`

$ cat /root/astrate-mule/pkg/astarteapi/testdata/error_bad_request.json
{"errors":{"detail":"Bad request"}}

$ cat /root/astrate-mule/pkg/astarteapi/testdata/error_internal.json
{"errors":{"detail":"Internal server error"}}

$ grep -rn "TestHousekeepingErrorDetailExamples\|TestPairingErrorDetailExamples\|TestRealmManagementErrorDetailExamples" /root/astrate-mule/internal/swagger
/root/astrate-mule/internal/swagger/swagger_test.go:460:// TestPairingErrorDetailExamples pins the error-detail examples in the pairing
/root/astrate-mule/internal/swagger/swagger_test.go:470:func TestPairingErrorDetailExamples(t *testing.T) {
/root/astrate-mule/internal/swagger/swagger_test.go:509:// TestRealmManagementErrorDetailExamples pins the error-detail examples in the
/root/astrate-mule/internal/swagger/swagger_test.go:512:// from TestPairingErrorDetailExamples. The BadRequest and InternalServerError
/root/astrate-mule/internal/swagger/swagger_test.go:521:func TestRealmManagementErrorDetailExamples(t *testing.T) {
/root/astrate-mule/internal/swagger/swagger_test.go:564:// management specs get from TestPairingErrorDetailExamples and
/root/astrate-mule/internal/swagger/swagger_test.go:565:// TestRealmManagementErrorDetailExamples. The BadRequest and InternalServerError

$ sed -n '440,550p' /root/astrate-mule/internal/swagger/swagger_test.go
for i, l := range block {
			if l == key {
				idx = i
				break
			}
		}
		if idx < 0 {
			t.Errorf("operation %s documents no 400 response", op)
			continue
		}
		if idx+1 >= len(block) || block[idx+1] != ref {
			t.Errorf("operation %s 400 is followed by %q, want the BadRequest $ref", op, block[idx+1])
		}
	}

	if !containsLine(lines, "    BadRequest:") {
		t.Error("components.responses defines no BadRequest response")
	}
}

// TestPairingErrorDetailExamples pins the error-detail examples in the pairing
// spec's components.responses to the frozen canonical strings the wire emits
// (astarteapi/envelope.go), so a generated client that copies an example does
// not learn a detail that never goes over the wire. Phoenix renders "Bad
// request" and "Internal server error" — the capital-R/capital-S spellings the
// spec used to carry were reconstructed, not measured on upstream 1.2.0, and a
// startectl/SDK error path matching on the canonical string would miss them.
// The 404 detail is deliberately absent here: pairing answers an unmatched
// path with DetailPageNotFound while DetailNotFound is the generic form, and
// the spec's own DeviceNotFound example is the one that carries a constant.
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

// TestRealmManagementErrorDetailExamples pins the error-detail examples in the
// realm management spec's components.responses to the frozen canonical strings
// the wire emits (astarteapi/envelope.go), the same guard the pairing spec gets
// from TestPairingErrorDetailExamples. The BadRequest and InternalServerError
// examples carried the capital-R/capital-S spellings, which upstream Phoenix
// never renders and no handler emits: every 400 on this surface goes through
// astarteapi.WriteBadRequest and every 500 through WriteInternalServerError
// (internal/realm/http.go), and the SDK/astartectl error paths match on the
// canonical strings. The Conflict example is an endpoint-specific detail with no
// constant behind it, so it stays free-form; the three 422s that used to share
// one ValidationError component are pinned by
// TestRealmManagement422ValidationDetails.
func TestRealmManagementErrorDetailExamples(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
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

I'll start by fixing the housekeeping API YAML error examples as specified in the task.
← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -358,9 +358,9 @@
         detail:
           type: string
           description: Human-readable error detail string.
       example:
-        detail: Bad Request
+        detail: Bad request
 
 FieldErrorsDetail:
   type: object
   required: [errors]

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -381,9 +381,9 @@
       schema:
         $ref: "#/components/schemas/ErrorDetail"
       example:
         errors:
-          detail: Bad Request
+          detail: Bad request
 
 Unauthorized:
   description: Unauthorized — missing or invalid JWT.
   content:

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -473,5 +473,5 @@
 schema:
   $ref: "#/components/schemas/ErrorDetail"
 example:
   errors:
-    detail: Internal Server Error
+    detail: Internal server error

$ sed -n '355,365p' /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
type: object
          required: [detail]
          properties:
            detail:
              type: string
              description: Human-readable error detail string.
          example:
            detail: Bad request

    FieldErrorsDetail:
      type: object

$ sed -n '380,390p' /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Bad request

    Unauthorized:
      description: Unauthorized — missing or invalid JWT.
      content:
        application/json; charset=utf-8:

$ sed -n '470,480p' /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
description: Internal Server Error — unexpected failure.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Internal server error

The ErrorDetail schema has a schema-level example with wrong shape. Let us remove it to match other specs.
← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -357,11 +357,10 @@
       properties:
         detail:
           type: string
           description: Human-readable error detail string.
-      example:
-        detail: Bad request
 
+
 FieldErrorsDetail:
   type: object
   required: [errors]
   properties:

$ sed -n '352,362p' /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
required: [errors]
      properties:
        errors:
          type: object
          required: [detail]
          properties:
            detail:
              type: string
              description: Human-readable error detail string.

$ grep -rn "TestHousekeepingErrorDetailExamples" /root/astrate-mule/internal/swagger
(no output)

Let us add the new test function.
$ tail -30 /root/astrate-mule/internal/swagger/swagger_test.go
}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".yaml") {
			names = append(names, strings.TrimPrefix(path, "api/"))
		}
		return nil
	})
	sort.Strings(names)
	return names, err
}

func get(t *testing.T, url string) string {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s: status = %d, want %d", url, resp.StatusCode, http.StatusOK)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading %s: %v", url, err)
	}
	return string(b)
}

$ sed -n '560,570p' /root/astrate-mule/internal/swagger/swagger_test.go
// TestAppEngineErrorDetailExamples pins the error-detail examples in the app
// engine spec's components.responses to the frozen canonical strings the wire
// emits (astarteapi/envelope.go), the same guard the pairing and realm
// management specs get from TestPairingErrorDetailExamples and
// TestRealmManagementErrorDetailExamples. The BadRequest and InternalServerError
// examples carried the capital-R/capital-S spellings, which upstream Phoenix
// never renders and no handler emits: every 400 on this surface goes through
// astarteapi.WriteBadRequest and every 500 through WriteInternalServerError
// (internal/appengine/http.go), and the SDK/astartectl error paths match on the
// canonical strings. The responses whose example is an endpoint-specific detail

$ sed -n '600,650p' /root/astrate-mule/internal/swagger/swagger_test.go
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

// TestRealmManagement422ValidationDetails pins the 422 examples the three
// non-interface realm-management operations carry to the details the wire
// actually emits, so a generated client that copies one does not learn a
// message no handler can produce.
//
// The three operations used to share components.responses.ValidationError, whose
// single example read `invalid interface: missing "interface_name"`. That string
// is wrong for all three: it names an interface, and none of them installs or
// updates one. Each now carries its own inline 422, the pattern
// installInterface and updateInterface already use for their own details.
//
// The createPolicy examples are derived from the compiler here instead of being
// written out, so a change to the message fails here instead of drifting the
// spec — the coupling TestPairingDeviceIDEncodingDocumented gets from
// deviceid. The putAuthConfig and deletePolicy details have no exported
// constant behind them (Service.SetAuthKey and Service.DeletePolicy,
// internal/realm/service.go), so they are spelled out and cited here.
func TestRealmManagement422ValidationDetails(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	if strings.Contains(string(b), "#/components/responses/ValidationError") {
		t.Error("spec still $refs a shared ValidationError response; one example cannot be right for every surface that refs it")
	}
	if containsLine(lines, "    ValidationError:") {
		t.Error("spec still defines components.responses.ValidationError")
	}

	_, noHandlers := triggers.CompilePolicy([]byte(`{"name":"p","error_handlers":[]}`))
	if noHandlers == nil {
		t.Fatal("triggers.CompilePolicy accepted a policy with no error handler; this test's premise no longer holds")
	}
	_, badName := triggers.CompilePolicy(
		[]byte(`{"name":"","error_handlers":[{"on":"any_error","strategy":"discard"}]}`))
	if badName == nil {

$ sed -n '650,680p' /root/astrate-mule/internal/swagger/swagger_test.go
if badName == nil {
		t.Fatal("triggers.CompilePolicy accepted an empty policy name; this test's premise no longer holds")
	}

	for _, tc := range []struct {
		op   string
		want []string
	}{
		{"putAuthConfig", []string{`jwt_public_key_pem can't be blank`}},
		{"createPolicy", []string{noHandlers.Error(), badName.Error()}},
		{"deletePolicy", []string{`policy "audit" is still used by trigger "on_audit_failure"`}},
	} {
		block := strings.Join(responseBlock(t, operationBlock(t, lines, tc.op), "422"), "\n")
		for _, want := range tc.want {
			if !strings.Contains(block, want) {
				t.Errorf("%s 422 carries no example with detail %q", tc.op, want)
			}
		}
		if strings.Contains(block, "#/components/responses/") {
			t.Errorf("%s 422 still $refs a response component instead of carrying its own example", tc.op)
		}
	}
}

// responseBlock returns the lines of the operation's response with the given
// status code, up to the next status code or the end of the operation block.
func responseBlock(t *testing.T, block []string, status string) []string {
	t.Helper()
	key := `        "` + status + `":`
	start := -1
	for i, l := range block {

$ tail -15 /root/astrate-mule/internal/swagger/swagger_test.go
resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s: status = %d, want %d", url, resp.StatusCode, http.StatusOK)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading %s: %v", url, err)
	}
	return string(b)
}

$ cat >> /root/astrate-mule/internal/swagger/swagger_test.go << 'EOF'

// TestHousekeepingErrorDetailExamples pins the error-detail examples in the
// housekeeping spec's components.responses to the frozen canonical strings the
// wire emits (astarteapi/envelope.go), so a generated client that copies an
// example does not learn a detail that never goes over the wire. Phoenix
// renders "Bad request" and "Internal server error" — the capital-R/capital-S
// spellings the spec used to carry were reconstructed, not measured on upstream
// 1.2.0, and every 400 on this surface goes through astarteapi.WriteBadRequest
// and every 500 through astarteapi.WriteInternalServerError (internal/housekeeping/http.go),
// and the SDK/astartectl error paths match on the canonical strings.
func TestHousekeepingErrorDetailExamples(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
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
EOF

(no output)

$ cd /root/astrate-mule && make -C docs build 2>&1 | tail -30
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

$ which mkdocs pip3 python3 2>&1 | head -5
/usr/bin/python3

$ python3 -m mkdocs --version 2>&1 | head -2
/usr/bin/python3: No module named mkdocs

$ pip3 show mkdocs 2>&1 | head -2; pip3 show mkdocs-material 2>&1 | head -2
/bin/sh: 1: pip3: not found
/bin/sh: 1: pip3: not found

$ ls /usr/local/bin/*mkdocs* /usr/bin/*mkdocs* 2>&1 | head -3
ls: cannot access '/usr/local/bin/*mkdocs*': No such file or directory
ls: cannot access '/usr/bin/*mkdocs*': No such file or directory

$ find /usr -name 'mkdocs' -not -path '*/site-packages/*' 2>&1 | head -5
(no output)

$ find /usr -path '*mkdocs*' 2>&1 | head -10
/usr/share/fish/completions/mkdocs.fish

