slug: docs-sync-hk-async-operation-param
verdict: done
at:  ed1cd72
ran: 2026-09-25T19:35:04Z on DietPi in 248s


> build · big-pickle

I'll read the MULE.md and task.md files first.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_housekeeping_api.yaml

$ sed -n '330,350p' docs/COMPATIBILITY.md
`422` afterwards; Astrate's pairing door accepts re-registration of a
    not-yet-confirmed device (the minted secret is discarded) and refuses an
    already-confirmed one (`ErrAlreadyRegistered`, mapped by the block to a
    log-and-drop).

17. **Always synchronous where upstream 1.4 defaults to asynchronous** —
    upstream 1.4 runs realm create/delete (Housekeeping), interface
    install/update/delete (Realm Management) and trigger-delivery-policy
    delete in the background, answering immediately and letting the caller opt
    into synchronous execution with `?async_operation=false`. Astrate performs
    every one of these synchronously and answers only once the work is done —
    a strictly stronger guarantee, so a client written against upstream never
    observes a regression. The parameter is accepted and ignored on either
    value, on every one of those endpoints, so upstream clients that send it
    keep working unchanged; pinned by
    `TestRealmManagementAsyncOperationParam` and
    `TestHousekeepingAsyncOperationParam`. Tracked in
    `docs/UPSTREAM-EXPERIMENTAL.md` (#68) until upstream 1.4 is final.

18. **Per-service health endpoints answer 503 when a dependency is down** —
    upstream's unauthenticated `GET /{appengine,realmmanagement,pairing}/health`,

$ sed -n '50,100p' internal/housekeeping/http.go; echo ---; sed -n '1,60p' internal/housekeeping/async_operation_test.go
DatastreamMaximumStorageRetention *int64 `json:"datastream_maximum_storage_retention"`
}

func (a *API) listRealms(w http.ResponseWriter, r *http.Request) {
	names, err := a.svc.ListRealms(r.Context())
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, names)
}

func (a *API) createRealm(w http.ResponseWriter, r *http.Request) {
	var req realmBody
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &req); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	rv, err := a.svc.CreateRealm(r.Context(), req.RealmName, req.JWTPublicKeyPEM, req.DeviceRegistrationLimit, req.DatastreamMaximumStorageRetention)
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusCreated, viewBody(rv))
}

func (a *API) getRealm(w http.ResponseWriter, r *http.Request) {
	rv, err := a.svc.GetRealm(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, viewBody(rv))
}

func (a *API) deleteRealm(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.DeleteRealm(r.Context(), r.PathValue("realm")); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// viewBody renders a RealmView as its wire shape.
func viewBody(rv *RealmView) realmBody {
	return realmBody{
		RealmName:                         rv.Name,
		JWTPublicKeyPEM:                   rv.JWTPublicKeyPEM,
		DeviceRegistrationLimit:           rv.DeviceRegistrationLimit,
		DatastreamMaximumStorageRetention: rv.DatastreamMaximumStorageRetention,
	}
---
//go:build integration

package housekeeping

import (
	"net/http"
	"testing"
)

// TestHousekeepingAsyncOperationParam pins deviation 17
// (docs/COMPATIBILITY.md): upstream 1.4 runs realm create and delete
// asynchronously unless the caller opts out with `?async_operation=false`.
// Astrate is always synchronous, so the parameter is accepted and ignored on
// both values — a client that sends it must never see a 4xx for it, and the
// response must be identical to the parameterless call.
func TestHousekeepingAsyncOperationParam(t *testing.T) {
	r := newHKRig(t)

	for _, q := range []string{"?async_operation=false", "?async_operation=true"} {
		realmName := "as" + randSuffix(t)
		createBody := `{"realm_name":` + jsonStr(realmName) + `,"jwt_public_key_pem":` + jsonStr(r.realmKey) + `}`

		if rec := r.req(t, http.MethodPost, "/realms"+q, createBody, r.haToken); rec.Code != http.StatusCreated {
			t.Fatalf("create realm %s: got %d, want 201 (%s)", q, rec.Code, rec.Body)
		}
		// Synchronous means done on return: the realm reads back immediately.
		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusOK {
			t.Errorf("get realm after create %s: got %d, want 200", q, rec.Code)
		}
		if rec := r.req(t, http.MethodDelete, "/realms/"+realmName+q, "", r.haToken); rec.Code != http.StatusNoContent {
			t.Fatalf("delete realm %s: got %d, want 204 (%s)", q, rec.Code, rec.Body)
		}
		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusNotFound {
			t.Errorf("get realm after delete %s: got %d, want 404", q, rec.Code)
		}
	}
}

$ rg -n "async_operation|AsyncOperation" --glob '!*.md' -g '!internal/housekeeping/*' . | head -40
./internal/realm/async_operation_test.go:11:// TestRealmManagementAsyncOperationParam pins deviation 17
./internal/realm/async_operation_test.go:14:// `?async_operation=false`. Astrate is always synchronous, so the parameter is
./internal/realm/async_operation_test.go:17:func TestRealmManagementAsyncOperationParam(t *testing.T) {
./internal/realm/async_operation_test.go:20:	for i, q := range []string{"?async_operation=false", "?async_operation=true"} {

$ ls docs; echo ---; cat docs/Makefile 2>/dev/null | head -40; echo ---; rg -rn "openapi|swagger" --glob '*.go' -l internal cmd 2>/dev/null | head; echo ---; rg -n "docs/api" -g '!*.md' . | head -20
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
---
cmd/astrate/main.go
internal/swagger/swagger.go
internal/swagger/swagger_test.go
---

$ rg -n "Accepted and ignored|accepted and ignored|deviation" docs/api/*.yaml | head -20; echo ---; rg -n "Async|async" docs/api/*.yaml | head -20
docs/api/astarte_realm_management_api.yaml:615:        a recorded deviation from upstream, which answers null here (issue #60).
docs/api/astarte_realm_management_api.yaml:686:        deviation from upstream's async deletion (see COMPATIBILITY.md).
---
docs/api/astarte_realm_management_api.yaml:686:        deviation from upstream's async deletion (see COMPATIBILITY.md).

$ sed -n '600,700p' docs/api/astarte_realm_management_api.yaml; echo === ; cat internal/swagger/swagger.go; echo ===; cat docs/embed.go
data: 100
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "500":
          $ref: "#/components/responses/InternalServerError"

  /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention:
    get:
      operationId: getDatastreamMaximumStorageRetention
      summary: Get datastream maximum storage retention
      description: |
        Returns the realm's datastream maximum storage retention in seconds
        (upstream serves this since 1.2.0). An unset ceiling renders as 0 —
        a recorded deviation from upstream, which answers null here (issue #60).
      tags: [Realm Management - Config]
      security:
        - a_rma: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
      responses:
        "200":
          description: Datastream maximum storage retention in seconds.
          content:
            application/json; charset=utf-8:
              schema:
                type: object
                required: [data]
                properties:
                  data:
                    type: integer
                    format: int64
                    description: Maximum datastream storage retention in seconds. 0 means unset.
              example:
                data: 31536000
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "500":
          $ref: "#/components/responses/InternalServerError"

  # ── Version ─────────────────────────────────────────────────────────

  /realmmanagement/v1/{realm}/version:
    get:
      operationId: getVersion
      summary: Get emulated API version
      description: |
        Returns the emulated upstream API compatibility version. The Astarte
        Dashboard gates feature UI on this value.
      tags: [Realm Management - Misc]
      security:
        - a_rma: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
      responses:
        "200":
          description: API compatibility version.
          content:
            application/json; charset=utf-8:
              schema:
                type: object
                required: [data]
                properties:
                  data:
                    type: string
                    description: Version string.
              example:
                data: "1.1.0"
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
      responses:
        "204":
          description: Device deleted. No body returned.
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
===
// Package swagger serves the embedded Swagger UI and OpenAPI YAML specs at
// /swagger/ and /api/ respectively.
package swagger

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	docs "github.com/astrate-platform/astrate/docs"
)

// Mount registers the /swagger and /api routes on the given mux.
// /swagger redirects to /swagger/index.html; /swagger/ serves the static UI;
// /api/ serves the OpenAPI YAML specs.
func Mount(mux *http.ServeMux) {
	uiRoot := mustSub(docs.SwaggerUI, "swagger-ui")
	apiRoot := mustSub(docs.APIYAML, "api")

	MountWithFS(mux, uiRoot, apiRoot)
}

// MountWithFS registers the /swagger and /api routes on the given mux from
// already-subtreed file systems. Mount computes the sub-trees from the
// embedded docs and delegates here; injecting fs.FS values makes the served
// roots testable and lets a broken embed fail fast instead of silently
// serving an empty tree.
func MountWithFS(mux *http.ServeMux, uiRoot, apiRoot fs.FS) {
	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
	})
	mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(uiRoot))))

	// Serve YAML files at /api/ so the relative ../api/*.yaml in index.html
	// resolves correctly when the page is loaded from /swagger/index.html.
	mux.Handle("GET /api/", http.StripPrefix("/api/", http.FileServer(http.FS(apiRoot))))
}

// mustSub returns fs.Sub(fsys, name), panicking if the sub-tree is absent so
// a broken docs embed fails fast instead of silently serving an empty tree.
// fs.Sub itself only errors on an invalid path (embed.FS has no Sub method),
// so the sub-root is also opened to catch an embed whose layout dropped the
// tree entirely.
func mustSub(fsys fs.FS, name string) fs.FS {
	sub, err := fs.Sub(fsys, name)
	if err != nil {
		panic(fmt.Sprintf("swagger: fs.Sub(%q): %v", name, err))
	}
	if _, err := sub.Open("."); err != nil {
		panic(fmt.Sprintf("swagger: %q missing from embedded fs: %v", name, err))
	}
	return sub
}

// Specs returns the list of available YAML spec filenames (without path prefix).
func Specs() []string {
	var names []string
	_ = fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".yaml") {
			names = append(names, strings.TrimPrefix(path, "api/"))
		}
		return nil
	})
	return names
}
===
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

$ rg -ln "astarte_housekeeping_api" . | head; echo ===; ls docs/api
./docs/swagger-ui/index.html
./docs/handoff/phase-1-handoff.md
./docs/handoff/phase-1-memory.md
===
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

$ rg -n "yaml|api/|Specs\(\)" internal/swagger/swagger_test.go | head -30; echo ===; which mkdocs python3; python3 -c "import mkdocs, yaml; print('ok')" 2>&1 | tail -2
54:	t.Run("GET /api/ serves every OpenAPI YAML spec", func(t *testing.T) {
55:		for _, name := range Specs() {
56:			body := get(t, srv.URL+"/api/"+name)
57:			want, err := docs.APIYAML.ReadFile("api/" + name)
62:				t.Errorf("served /api/%s does not match embedded copy", name)
69:// api/ tree: an embed.FS declared without a //go:embed directive is empty, and
71:// fail fast on instead of serving 404s at /swagger/ and /api/.
76:// tree that 404s at /swagger/ and /api/.
91:	got := Specs()
94:		t.Fatal("Specs() returned no filenames")
98:		if strings.HasPrefix(name, "api/") || strings.Contains(name, "/") {
99:			t.Errorf("Specs() entry %q should have no path prefix or dirs", name)
101:		if !strings.HasSuffix(name, ".yaml") {
102:			t.Errorf("Specs() entry %q does not end in .yaml", name)
104:		if _, err := docs.APIYAML.ReadFile("api/" + name); err != nil {
105:			t.Errorf("Specs() entry %q cannot be read from docs.APIYAML: %v", name, err)
114:		t.Errorf("Specs() = %v, want %v", got, want)
124:	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
126:		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
190:		if strings.HasSuffix(path, ".yaml") {
191:			names = append(names, strings.TrimPrefix(path, "api/"))
===
/usr/bin/python3
    ^^^^^^^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'mkdocs'

→ Read internal/swagger/swagger_test.go [offset=115, limit=60]

$ rg -ln "docs.APIYAML|astarte_housekeeping_api.yaml" --glob '*_test.go' .; echo ===; rg -n "COMPATIBILITY" --glob '*_test.go' -l . | head
./internal/swagger/swagger_test.go
===
./cmd/astrate/version_test.go
./internal/realm/async_operation_test.go
./internal/observability/observability_test.go
./internal/housekeeping/async_operation_test.go
./internal/appengine/e2e_test.go

$ rg -n "yaml" go.mod go.sum | head; echo ===; rg -n "yaml" --glob '*.go' -l . | head
go.mod:93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
go.mod:99:	gopkg.in/yaml.v3 v3.0.1 // indirect
go.sum:202:go.yaml.in/yaml/v2 v2.4.2 h1:DzmwEr2rDGHl7lsFgAHxmNz/1NlQ7xLIrlN2h5d1eGI=
go.sum:203:go.yaml.in/yaml/v2 v2.4.2/go.mod h1:081UH+NErpNdqlCXm3TtEran0rJZGxAYx9hb/ELlsPU=
go.sum:225:gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
go.sum:226:gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
go.sum:227:gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
===
./docs/embed.go
./internal/swagger/swagger.go
./internal/swagger/swagger_test.go

$ python3 -c "import yaml; print(yaml.__version__)" 2>&1 | tail -1; echo ===; rg -n "yaml" docs/swagger-ui/index.html | head; echo ===; cat docs/requirements.txt
ModuleNotFoundError: No module named 'yaml'
===
44:      <option value="../api/astarte_housekeeping_api.yaml">Housekeeping API — Realm management, a_ha auth</option>
45:      <option value="../api/astarte_pairing_api.yaml">Pairing API — Device pairing flows A/B/C, a_pa auth</option>
46:      <option value="../api/astarte_realm_management_api.yaml">Realm Management API — Interfaces, triggers, policies, a_rma auth</option>
47:      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
48:      <option value="../api/astrate_native_api.yaml">Astrate Native API — Health, readiness, metrics, WebSocket</option>
===
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10

$ command -v pip3 pip; timeout 60 pip3 install -q -r docs/requirements.txt 2>&1 | tail -5; echo "exit=$?"
timeout: failed to run command ‘pip3’: No such file or directory
exit=0

$ ls "$(go env GOMODCACHE)/gopkg.in" 2>/dev/null | head; echo ---; ls "$(go env GOMODCACHE)/go.yaml.in/yaml" 2>/dev/null | head; echo ---; go env GOFLAGS GOPROXY
check.v1@v0.0.0-20161208181325-20d25e280405
check.v1@v1.0.0-20201130134442-10cb98267c6c
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1
---
v2@v2.4.2
v3@v3.0.4
---

https://proxy.golang.org,direct

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -56,8 +56,10 @@
   registration limit.
 tags: [Housekeeping]
 security:
   - a_ha: []
+parameters:
+  - $ref: "#/components/parameters/AsyncOperation"
 requestBody:
   required: true
   content:
     application/json; charset=utf-8:

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -194,8 +194,9 @@
 security:
   - a_ha: []
 parameters:
   - $ref: "#/components/parameters/RealmName"
+  - $ref: "#/components/parameters/AsyncOperation"
 responses:
   "204":
     description: Realm deleted successfully. No body returned.
   "401":

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -230,8 +230,22 @@
     description: The realm name.
     schema:
       type: string
 
+  AsyncOperation:
+    name: async_operation
+    in: query
+    required: false
+    description: |
+      Accepted and ignored on either value. Upstream 1.4 runs this operation in
+      the background and lets the caller opt into synchronous execution with
+      `async_operation=false`; Astrate always works synchronously and answers
+      only once the work is done, so the value changes nothing and the
+      response is identical either way (deviation 17 in COMPATIBILITY.md).
+    schema:
+      type: boolean
+      default: false
+
 schemas:
   # ── Envelope variants ──
 
   DataEnvelopeRealm:

→ Read internal/swagger/swagger_test.go [limit=55]
→ Read internal/swagger/swagger_test.go [offset=170]

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -168,8 +168,97 @@
 		t.Error("components.responses defines no Forbidden response")
 	}
 }
 
+// TestHousekeepingAsyncOperationParamDocumented guards that the housekeeping
+// spec tells clients the `?async_operation` parameter exists. Astrate accepts
+// and ignores it on realm create and delete (deviation 17,
+// docs/COMPATIBILITY.md), so an upstream client keeps working — but a client
+// generated from the spec has to learn the parameter from the spec. This is
+// the documentation half of the behaviour pinned by
+// TestHousekeepingAsyncOperationParam in internal/housekeeping.
+func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	const ref = `        - $ref: "#/components/parameters/AsyncOperation"`
+	for _, op := range []string{"createRealm", "deleteRealm"} {
+		block := operationBlock(t, lines, op)
+		if !containsLine(block, ref) {
+			t.Errorf("operation %s does not $ref the AsyncOperation parameter", op)
+		}
+	}
+
+	comp := componentBlock(t, lines, "    AsyncOperation:")
+	for _, want := range []string{
+		"      name: async_operation",
+		"      in: query",
+		"      required: false",
+		"      schema:",
+		"        type: boolean",
+		"        default: false",
+	} {
+		if !containsLine(comp, want) {
+			t.Errorf("components parameter AsyncOperation is missing line %q", want)
+		}
+	}
+	if !strings.Contains(strings.Join(comp, "\n"), "Accepted and ignored") {
+		t.Error("components parameter AsyncOperation does not say the value is accepted and ignored")
+	}
+}
+
+// operationBlock returns the lines of the operation with the given operationId,
+// up to the next operation, path, or the components section.
+func operationBlock(t *testing.T, lines []string, operationID string) []string {
+	t.Helper()
+	marker := "      operationId: " + operationID
+	start := -1
+	for i, l := range lines {
+		if l == marker {
+			start = i
+			break
+		}
+	}
+	if start < 0 {
+		t.Fatalf("spec declares no operation %q", operationID)
+	}
+	endRe := regexp.MustCompile(`^      operationId: |^  \S|^components:`)
+	for i := start + 1; i < len(lines); i++ {
+		if endRe.MatchString(lines[i]) {
+			return lines[start+1 : i]
+		}
+	}
+	return lines[start+1:]
+}
+
+// componentBlock returns the lines of the components entry starting with the
+// given key line, up to the next entry at the same or shallower indentation.
+func componentBlock(t *testing.T, lines []string, key string) []string {
+	t.Helper()
+	start := -1
+	for i, l := range lines {
+		if l == key {
+			start = i
+			break
+		}
+	}
+	if start < 0 {
+		t.Fatalf("spec declares no components entry %q", key)
+	}
+	for i := start + 1; i < len(lines); i++ {
+		if strings.TrimSpace(lines[i]) == "" {
+			continue
+		}
+		if !strings.HasPrefix(lines[i], "      ") {
+			return lines[start+1 : i]
+		}
+	}
+	return lines[start+1:]
+}
+
 func containsLine(lines []string, want string) bool {
 	for _, l := range lines {
 		if l == want {
 			return true

$ gofmt -l . 2>&1 | grep -v '^$' ; echo "gofmt done"; go vet ./internal/swagger/ && go test ./internal/swagger/ 2>&1 | tail -20
gofmt done
ok  	github.com/astrate-platform/astrate/internal/swagger	0.063s

Now verifying the test actually fails without the spec change (reverting the YAML temporarily):
$ cp docs/api/astarte_housekeeping_api.yaml /tmp/opencode/hk.new.yaml && git show HEAD:docs/api/astarte_housekeeping_api.yaml > docs/api/astarte_housekeeping_api.yaml && go test ./internal/swagger/ -run TestHousekeepingAsyncOperationParamDocumented 2>&1 | tail -20; cp /tmp/opencode/hk.new.yaml docs/api/astarte_housekeeping_api.yaml; git diff --stat
--- FAIL: TestHousekeepingAsyncOperationParamDocumented (0.00s)
    swagger_test.go:190: operation createRealm does not $ref the AsyncOperation parameter
    swagger_test.go:190: operation deleteRealm does not $ref the AsyncOperation parameter
    swagger_test.go:194: spec declares no components entry "    AsyncOperation:"
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.036s
FAIL
 docs/api/astarte_housekeeping_api.yaml | 17 +++++++
 internal/swagger/swagger_test.go       | 89 ++++++++++++++++++++++++++++++++++
 2 files changed, 106 insertions(+)

$ mkdir -p /tmp/opencode/yamlprobe && cat > /tmp/opencode/yamlprobe/go.mod <<'EOF'
module yamlprobe

go 1.24

require gopkg.in/yaml.v3 v3.0.1
EOF
cat > /tmp/opencode/yamlprobe/main.go <<'EOF'
package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	bad := 0
	for _, f := range os.Args[1:] {
		b, err := os.ReadFile(f)
		if err != nil {
			fmt.Printf("READ FAIL %s: %v\n", f, err)
			bad++
			continue
		}
		var root map[string]any
		if err := yaml.Unmarshal(b, &root); err != nil {
			fmt.Printf("PARSE FAIL %s: %v\n", f, err)
			bad++
			continue
		}
		var walk func(prefix string, v any) int
		walk = func(prefix string, v any) int {
			n := 0
			switch t := v.(type) {
			case map[string]any:
				for k, val := range t {
					if s, ok := val.(string); ok && strings.HasPrefix(s, "#/") {
						cur := root
						for _, seg := range strings.Split(strings.TrimPrefix(s, "#/"), "/") {
							nx, ok := cur[seg]
							if !ok {
								fmt.Printf("DANGLING REF %s -> %s\n", f, s)
								n++
								bad++
								break
							}
							cur, _ = nx.(map[string]any)
						}
						continue
					}
					n += walk(prefix+"/"+k, val)
				}
			case []any:
				for i, val := range t {
					n += walk(fmt.Sprintf("%s/%d", prefix, i), val)
				}
			}
			return n
		}
		refs := walk("", root)
		paths, _ := root["paths"].(map[string]any)
		fmt.Printf("OK %-50s paths=%d refs=%d\n", f, len(paths), refs)
	}
	if bad > 0 {
		os.Exit(1)
	}
}
EOF
cd /tmp/opencode/yamlprobe && GOFLAGS=-mod=mod GOPROXY=off go mod tidy >/dev/null 2>&1; GOPROXY=off go run . /root/astrate-mule/docs/api/*.yaml
OK /root/astrate-mule/docs/api/astarte_appengine_api.yaml paths=17 refs=0
OK /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml paths=2 refs=0
OK /root/astrate-mule/docs/api/astarte_pairing_api.yaml paths=7 refs=0
OK /root/astrate-mule/docs/api/astarte_realm_management_api.yaml paths=12 refs=0
OK /root/astrate-mule/docs/api/astrate_native_api.yaml paths=12 refs=0

$ cd /tmp/opencode/yamlprobe && python3 - <<'EOF'
import re
s = open('main.go').read()
s = s.replace("""							cur, _ = nx.(map[string]any)
						}
						continue""", """							cur, _ = nx.(map[string]any)
						}
						n++
						continue""")
open('main.go','w').write(s)
EOF
sed 's|AsyncOperation|AsyncOperationTypo|' /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml > /tmp/opencode/hk_broken.yaml
GOPROXY=off go run . /root/astrate-mule/docs/api/*.yaml /tmp/opencode/hk_broken.yaml; echo "exit=$?"
OK /root/astrate-mule/docs/api/astarte_appengine_api.yaml paths=17 refs=308
OK /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml paths=2 refs=46
OK /root/astrate-mule/docs/api/astarte_pairing_api.yaml paths=7 refs=59
OK /root/astrate-mule/docs/api/astarte_realm_management_api.yaml paths=12 refs=150
OK /root/astrate-mule/docs/api/astrate_native_api.yaml paths=12 refs=22
OK /tmp/opencode/hk_broken.yaml                       paths=2 refs=46
exit=0

$ cd /tmp/opencode/yamlprobe && sed 's|parameters/AsyncOperation|parameters/AsyncOperationTypo|' /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml > /tmp/opencode/hk_broken.yaml && GOPROXY=off go run . /tmp/opencode/hk_broken.yaml; echo "exit=$?"; rm -f /tmp/opencode/hk_broken.yaml
DANGLING REF /tmp/opencode/hk_broken.yaml -> #/components/parameters/AsyncOperationTypo
DANGLING REF /tmp/opencode/hk_broken.yaml -> #/components/parameters/AsyncOperationTypo
OK /tmp/opencode/hk_broken.yaml                       paths=2 refs=48
exit=1
exit status 1

$ go vet ./... && go test ./... 2>&1 | grep -v "^ok\|no test files" | head -20; echo "gate exit=$?"
gate exit=0

$ make -C docs build 2>&1 | tail -5; echo "exit=${PIPESTATUS[0]}"
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'
/bin/sh: 1: Bad substitution

$ git status --porcelain | head -20; echo ===; git check-ignore -v docs/site/DESIGN.md docs/site/api/astarte_housekeeping_api.yaml 2>&1 | head
M docs/api/astarte_housekeeping_api.yaml
 M internal/swagger/swagger_test.go
===
.gitignore:32:/docs/site/DESIGN.md	docs/site/DESIGN.md
.gitignore:34:/docs/site/api/	docs/site/api/astarte_housekeeping_api.yaml

$ git diff; echo ===; tail -8 .mule/todo.md; echo ===; wc -l .mule/todo.md
diff --git a/docs/api/astarte_housekeeping_api.yaml b/docs/api/astarte_housekeeping_api.yaml
index f7b4f6d..bea5e74 100644
--- a/docs/api/astarte_housekeeping_api.yaml
+++ b/docs/api/astarte_housekeeping_api.yaml
@@ -57,6 +57,8 @@ paths:
       tags: [Housekeeping]
       security:
         - a_ha: []
+      parameters:
+        - $ref: "#/components/parameters/AsyncOperation"
       requestBody:
         required: true
         content:
@@ -193,6 +195,7 @@ paths:
         - a_ha: []
       parameters:
         - $ref: "#/components/parameters/RealmName"
+        - $ref: "#/components/parameters/AsyncOperation"
       responses:
         "204":
           description: Realm deleted successfully. No body returned.
@@ -228,6 +231,20 @@ components:
       schema:
         type: string
 
+    AsyncOperation:
+      name: async_operation
+      in: query
+      required: false
+      description: |
+        Accepted and ignored on either value. Upstream 1.4 runs this operation in
+        the background and lets the caller opt into synchronous execution with
+        `async_operation=false`; Astrate always works synchronously and answers
+        only once the work is done, so the value changes nothing and the
+        response is identical either way (deviation 17 in COMPATIBILITY.md).
+      schema:
+        type: boolean
+        default: false
+
   schemas:
     # ── Envelope variants ──
 
diff --git a/internal/swagger/swagger_test.go b/internal/swagger/swagger_test.go
index 626ca9c..c3926df 100644
--- a/internal/swagger/swagger_test.go
+++ b/internal/swagger/swagger_test.go
@@ -169,6 +169,95 @@ func TestRealmManagement403(t *testing.T) {
 	}
 }
 
+// TestHousekeepingAsyncOperationParamDocumented guards that the housekeeping
+// spec tells clients the `?async_operation` parameter exists. Astrate accepts
+// and ignores it on realm create and delete (deviation 17,
+// docs/COMPATIBILITY.md), so an upstream client keeps working — but a client
+// generated from the spec has to learn the parameter from the spec. This is
+// the documentation half of the behaviour pinned by
+// TestHousekeepingAsyncOperationParam in internal/housekeeping.
+func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	const ref = `        - $ref: "#/components/parameters/AsyncOperation"`
+	for _, op := range []string{"createRealm", "deleteRealm"} {
+		block := operationBlock(t, lines, op)
+		if !containsLine(block, ref) {
+			t.Errorf("operation %s does not $ref the AsyncOperation parameter", op)
+		}
+	}
+
+	comp := componentBlock(t, lines, "    AsyncOperation:")
+	for _, want := range []string{
+		"      name: async_operation",
+		"      in: query",
+		"      required: false",
+		"      schema:",
+		"        type: boolean",
+		"        default: false",
+	} {
+		if !containsLine(comp, want) {
+			t.Errorf("components parameter AsyncOperation is missing line %q", want)
+		}
+	}
+	if !strings.Contains(strings.Join(comp, "\n"), "Accepted and ignored") {
+		t.Error("components parameter AsyncOperation does not say the value is accepted and ignored")
+	}
+}
+
+// operationBlock returns the lines of the operation with the given operationId,
+// up to the next operation, path, or the components section.
+func operationBlock(t *testing.T, lines []string, operationID string) []string {
+	t.Helper()
+	marker := "      operationId: " + operationID
+	start := -1
+	for i, l := range lines {
+		if l == marker {
+			start = i
+			break
+		}
+	}
+	if start < 0 {
+		t.Fatalf("spec declares no operation %q", operationID)
+	}
+	endRe := regexp.MustCompile(`^      operationId: |^  \S|^components:`)
+	for i := start + 1; i < len(lines); i++ {
+		if endRe.MatchString(lines[i]) {
+			return lines[start+1 : i]
+		}
+	}
+	return lines[start+1:]
+}
+
+// componentBlock returns the lines of the components entry starting with the
+// given key line, up to the next entry at the same or shallower indentation.
+func componentBlock(t *testing.T, lines []string, key string) []string {
+	t.Helper()
+	start := -1
+	for i, l := range lines {
+		if l == key {
+			start = i
+			break
+		}
+	}
+	if start < 0 {
+		t.Fatalf("spec declares no components entry %q", key)
+	}
+	for i := start + 1; i < len(lines); i++ {
+		if strings.TrimSpace(lines[i]) == "" {
+			continue
+		}
+		if !strings.HasPrefix(lines[i], "      ") {
+			return lines[start+1 : i]
+		}
+	}
+	return lines[start+1:]
+}
+
 func containsLine(lines []string, want string) bool {
 	for _, l := range lines {
 		if l == want {
===
- [x] swagger-sub-failfast [auto]: make `Mount` fail fast (panic) when `fs.Sub` over the `docs.SwaggerUI`/`docs.APIYAML` embeds errors — `uiRoot, _ := fs.Sub(...)` and `apiRoot, _ := fs.Sub(...)` (internal/swagger/swagger.go:17-18) silently serve an empty tree → 404s at /swagger/ and /api/ with no log the day the docs embed layout changes; add an fs-injecting internal variant (`MountWithFS(mux, uiRoot, apiRoot fs.FS)` called by `Mount` after the subs) so a broken-FS unit test in swagger_test.go can prove the fail-fast, keeping the public `Mount(mux)` signature and its single main.go:414 caller untouched. Container-free.
- [ ] flow-boot-resumes-stopped-durable-flows [legion] [auto]: RehydrateAutoRestart (cmd/astrate/main.go:216) restarts every flow with auto_restart=true whatever its last status — `ListAutoRestartFlows` filters on `f.auto_restart = true` alone (internal/store/flows.go:125) and nothing ever clears the flag, which is written once at create (flowapi/http.go:141-146 -> service.go:318) and never on stop — so a durable flow stopped through the API is running again after the next boot, while the shutdown mark that says otherwise (MarkRunningFlowsStopped, main.go:497 -> internal/flowapi/service.go:885-901, which writes status="stopped") is never read by that query. The two halves disagree: service.go:325 documents "every durable flow", the shutdown call assumes "only what was running". Say which is the intent, then make them agree — either add `AND f.status = 'running'` to flows.go:125 (making the shutdown mark load-bearing) or delete the mark as dead code — and add a case in internal/store/flows_test.go pinning the chosen rule (a stopped auto_restart row must, or must not, come back). Needs the DB.
- [x] drain-per-stage-budget [auto]: `shutdown` spends a single 30s sctx (cmd/astrate/main.go:479) across three sequential stages — srv.Shutdown (484), b.Close (489), flowSvc.Manager().Shutdown (494) and MarkRunningFlowsStopped (497) — while `drainEngine` alone opens a second one (504), so the constant's own comment "bounds the whole graceful drain" (main.go:55-56) is wrong in both directions, and one slow in-flight HTTP request (the server sets no per-request deadline, 221) burns the whole first budget and leaves the flow stage on an expired context: StopFlow then returns at `f.router.Drain(ctx)` before step 3 releases block resources (internal/flow/flow.go:265-273) and MarkRunningFlowsStopped skips every flow (flowapi/service.go:895-897), both surfacing as one log.Warn. Give each stage its own budget through an extracted helper (shutdownTimeout has to become an injectable var for the test) and pin it container-free: a second stage still receives a live context after a first stage exhausts its own.
- [x] cmd-devcert-fields-test [auto]: add a container-free test for `selfSignedDevCert` (cmd/astrate/main.go:320-345), the throwaway mTLS identity `insecure_dev_mode` boots on, which no test touches — pin NotBefore = now-1h (331; a plain now would reject a device whose clock is behind), NotAfter = now+365d (332), DNSNames localhost and IPAddresses 127.0.0.1/::1 (336-337, dropping which breaks hostname verification for a device dialling mqtts://127.0.0.1:8883), ExtKeyUsage serverAuth (334), and that the leaf is self-signed and parses (x509.ParseCertificate on Certificate[0]).
- [!] cmd-healthcheck-contract-test [auto]: add a container-free test for `runHealthcheck` (cmd/astrate/main.go:531-558) — the `HEALTHCHECK CMD ["/astrate", "-healthcheck"]` contract (Dockerfile:35-36) has no coverage. Assert it probes /astrate/v1/readiness and not /health (543), returns 0 only on 200 so a 503 readiness exits 1 (554), and maps a bare ":PORT" in ASTRATE_HTTP_ADDR to 127.0.0.1 (536-538). httptest.NewServer plus t.Setenv("ASTRATE_HTTP_ADDR", srv.Listener.Addr().String()). — BLOCKED: tests failed: --- FAIL: TestRunHealthcheckProbesReadiness (0.01s)
- [x] cmd-loadsealer-masterkeyfile-test [auto]: add a container-free test for `loadSealer` (cmd/astrate/main.go:515-526) — the `security.master_key_file` to env hand-off (`os.Setenv(store.EnvMasterKeyFile, ...)`, 517) is the only wiring that makes a configured key file reach the store loader, and dropping or misspelling that branch leaves the file silently ignored and boot dying with the generic "master key" error (523) even though the config named a valid file. Assert a 32-byte key file under t.TempDir() loads, the inline ASTRATE_MASTER_KEY path still wins, and a missing key returns the error naming all three env vars. t.Setenv both ASTRATE_MASTER_KEY and ASTRATE_MASTER_KEY_FILE first so the os.Setenv is restored at cleanup.
- [ ] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
===
227 .mule/todo.md

$ cat >> .mule/todo.md <<'EOF'
- [ ] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
EOF
tail -1 .mule/todo.md | head -c 120; echo; git status --porcelain
- [ ] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four rea
 M .mule/todo.md
 M docs/api/astarte_housekeeping_api.yaml
 M internal/swagger/swagger_test.go

$ go test -race ./internal/swagger/ 2>&1 | tail -3; gofmt -l . | head
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.047s
FAIL

```
Done: documented the accepted-and-ignored `?async_operation` query parameter as an `AsyncOperation` component parameter in the housekeeping OpenAPI spec, `$ref`'d from `createRealm` (post) and `deleteRealm` (delete), with a description naming deviation 17, plus a container-free spec test that fails without the change.
Files: docs/api/astarte_housekeeping_api.yaml, internal/swagger/swagger_test.go, .mule/todo.md (follow-up line)
Verified: `go vet ./... && go test ./...` -> pass; `gofmt -l .` -> clean; `go test ./internal/swagger/ -run TestHousekeepingAsyncOperationParamDocumented` against the pre-change yaml (git show HEAD:...) -> FAIL as required, restored after; `/tmp` probe with yaml.v3 parsing all five docs/api/*.yaml -> all parse, all 46 refs in the housekeeping spec resolve (negative control with a typo'd $ref correctly reported DANGLING); `go test -race ./internal/swagger/` -> FATAL "Found 39 - Supported 48" (race detector unusable on this 39-bit VMA machine, so the gate is the non-race run above); `make -C docs build` -> could not run: `mkdocs: No such file or directory` and there is no pip3 on this box to install docs/requirements.txt, so the `sync` half did run (docs/site/ is gitignored, no tracked files touched)
Unsure: the mkdocs build itself is unverified — my substitute was "YAML parses, every $ref resolves, and internal/swagger serves the embedded copy byte-for-byte at /api/ (TestMount)", not an actual mkdocs render; the rm spec's four async-capable operations were left alone as instructed, so the RM gap is only recorded as a follow-up line
Follow-ups: - [ ] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
```
