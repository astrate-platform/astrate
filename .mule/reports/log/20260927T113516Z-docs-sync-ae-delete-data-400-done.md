slug: docs-sync-ae-delete-data-400
verdict: done
at:  66bc366
ran: 2026-09-27T11:35:16Z on DietPi in 330s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_appengine_api.yaml [offset=500, limit=70]

$ ls /root/astrate-mule && ls /root/astrate-mule/docs && ls /root/astrate-mule/docs/api 2>/dev/null
AGENTS.md
Dockerfile
LICENSE
Makefile
README.md
astrate
bench
cmd
deploy
docker-compose.yml
docs
examples
flow-container-echo
go.mod
go.sum
internal
migrations
opencode.json
pkg
site-dist
test
tools
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

$ rg -n "deleteDataByAlias|deleteData:|deleteDataInGroup" /root/astrate-mule/docs/api/astarte_appengine_api.yaml
544:      operationId: deleteDataByAlias
1352:      operationId: deleteDataInGroup

→ Read docs/api/astarte_appengine_api.yaml [limit=40, offset=565]

$ rg -n "^    delete:|^      operationId:" /root/astrate-mule/docs/api/astarte_appengine_api.yaml
28:      operationId: listDevices
102:      operationId: getDeviceStats
142:      operationId: getDevice
168:      operationId: patchDevice
220:      operationId: getDeviceByAlias
251:      operationId: patchDeviceByAlias
308:      operationId: listInterfacesByAlias
342:      operationId: getDataByAlias
389:      operationId: getDataAtPathByAlias
440:      operationId: putDataByAlias
492:      operationId: publishDataByAlias
543:    delete:
544:      operationId: deleteDataByAlias
583:      operationId: listDeviceInterfaces
613:      operationId: getData
655:      operationId: getDataAtPath
699:      operationId: putData
746:      operationId: publishData
792:    delete:
793:      operationId: deleteData
829:      operationId: listGroups
856:      operationId: createGroup
907:      operationId: getGroup
937:      operationId: listGroupDevices
986:      operationId: addGroupDevice
1033:      operationId: getGroupDevice
1062:      operationId: patchGroupDevice
1113:    delete:
1114:      operationId: removeGroupDevice
1138:      operationId: listInterfacesInGroup
1168:      operationId: getDataInGroup
1211:      operationId: getDataAtPathInGroup
1256:      operationId: putDataInGroup
1304:      operationId: publishDataInGroup
1351:    delete:
1352:      operationId: deleteDataInGroup

→ Read docs/api/astarte_appengine_api.yaml [limit=45, offset=786]
→ Read docs/api/astarte_appengine_api.yaml [limit=40, offset=1345]

$ sed -n '225,260p' /root/astrate-mule/internal/engine/serverdata.go; echo ---; sed -n '655,700p' /root/astrate-mule/internal/appengine/http.go
return fmt.Errorf("%w: %s", ErrRealmUnknown, realm)
	}
	view, err := e.deviceView(ctx, rs, id)
	if err != nil {
		return err
	}
	ci := resolveServerInterface(rs, view, ifaceName)
	if ci == nil {
		return fmt.Errorf("%w: %s", ErrInterfaceNotFound, ifaceName)
	}
	if ci.Ownership != interfaceschema.OwnershipServer {
		return fmt.Errorf("%w: %s", ErrNotServerOwned, ifaceName)
	}
	if ci.Type != interfaceschema.Properties {
		return fmt.Errorf("%w: %s", ErrNotAProperty, ifaceName)
	}
	mapping, ok := ci.Trie.Match(path)
	if !ok {
		return fmt.Errorf("%w: %q matches no endpoint of %s", ErrPathNotFound, path, ifaceName)
	}
	if !mapping.AllowUnset {
		return fmt.Errorf("%w: %s%s", ErrUnsetNotAllowed, ifaceName, path)
	}

	// Previous-value snapshot taken before the row is gone, so the
	// change-derived conditions (path_removed, value_change) evaluate against
	// the pre-unset value, then the op fires the same observers as any other
	// committed unset.
	op := PersistOp{
		Realm:     realm,
		RealmID:   rs.id,
		DeviceID:  id,
		Interface: ci,
		Mapping:   mapping,
		Path:      path,
		Kind:      OpPropertyUnset,
---
		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Invalid attributes")
	case errors.Is(err, ErrAttributeKeyNotFound):
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Attribute key not found")
	case errors.Is(err, ErrValidation):
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
	case errors.Is(err, ErrGroupAlreadyExists):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Group already exists")
	case errors.Is(err, ErrDeviceAlreadyInGroup):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Device already in group")
	case errors.Is(err, store.ErrAlreadyExists):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Already exists")
	case errors.Is(err, ErrGroupNotFound):
		_ = astarteapi.WriteError(w, http.StatusNotFound, "Group not found")
	case errors.Is(err, ErrPathNotFound):
		_ = astarteapi.WriteError(w, http.StatusNotFound, "Path not found")
	// Engine write sentinels must precede store.ErrNotFound: wrapped
	// ErrInterfaceNotFound errors are distinct from store's, and Go's switch
	// takes the first match (issue #57, measured taxonomy).
	case errors.Is(err, engine.ErrNotServerOwned):
		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Cannot write to device owned resource")
	case errors.Is(err, engine.ErrNotAProperty):
		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Cannot write to read-only resource")
	case errors.Is(err, engine.ErrInterfaceNotFound):
		_ = astarteapi.WriteError(w, http.StatusNotFound, "Interface not found in device introspection")
	case errors.Is(err, engine.ErrPathNotFound):
		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Endpoint not found")
	case errors.Is(err, store.ErrNotFound):
		_ = astarteapi.WriteDeviceNotFound(w)
	default:
		_ = astarteapi.WriteInternalServerError(w)
	}
}

// validationDetail strips the ErrValidation prefix for the response detail.
func validationDetail(err error) string {
	const prefix = "appengine: validation failed: "
	msg := err.Error()
	if strings.HasPrefix(msg, prefix) {
		return strings.TrimPrefix(msg, prefix)
	}
	return msg
}

→ Read docs/api/astarte_appengine_api.yaml [limit=50, offset=440]

$ rg -n "func \(s \*Server\) (deleteData|deleteDataByAlias|deleteDataInGroup)" /root/astrate-mule/internal/appengine/*.go; echo ---; rg -n "engine.ErrPathNotFound" /root/astrate-mule/internal/appengine/*.go | head -30; echo --- ; rg -n "Endpoint not found" /root/astrate-mule/internal/appengine/writeerror_test.go
---
/root/astrate-mule/internal/appengine/http.go:679:	case errors.Is(err, engine.ErrPathNotFound):
/root/astrate-mule/internal/appengine/writeerror_test.go:28:		{"unknown endpoint", fmt.Errorf("%w: context", engine.ErrPathNotFound),
---
29:			http.StatusBadRequest, "Endpoint not found"},

$ rg -n "deleteData" /root/astrate-mule/internal/appengine/*.go | head -30; echo ---; sed -n '1,60p' /root/astrate-mule/internal/appengine/writeerror_test.go
/root/astrate-mule/internal/appengine/http.go:54:	mux.Handle("DELETE "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteData))
/root/astrate-mule/internal/appengine/http.go:65:	mux.Handle("DELETE "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.deleteDataByAlias))
/root/astrate-mule/internal/appengine/http.go:80:	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteDataInGroup))
/root/astrate-mule/internal/appengine/http.go:311:func (a *API) deleteData(w http.ResponseWriter, r *http.Request) {
/root/astrate-mule/internal/appengine/http.go:359:func (a *API) deleteDataByAlias(w http.ResponseWriter, r *http.Request) {
/root/astrate-mule/internal/appengine/http.go:369:func (a *API) deleteDataInGroup(w http.ResponseWriter, r *http.Request) {
---
package appengine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astrate-platform/astrate/internal/engine"
)

func TestWriteErrorTaxonomy(t *testing.T) {
	a := &API{}
	cases := []struct {
		name   string
		err    error
		status int
		detail string
	}{
		{"device-owned write", fmt.Errorf("%w: context", engine.ErrNotServerOwned),
			http.StatusMethodNotAllowed, "Cannot write to device owned resource"},
		{"unset on server-owned datastream", fmt.Errorf("%w: context", engine.ErrNotAProperty),
			http.StatusMethodNotAllowed, "Cannot write to read-only resource"},
		{"unknown interface", fmt.Errorf("%w: context", engine.ErrInterfaceNotFound),
			http.StatusNotFound, "Interface not found in device introspection"},
		{"unknown endpoint", fmt.Errorf("%w: context", engine.ErrPathNotFound),
			http.StatusBadRequest, "Endpoint not found"},
		{"read-path path-not-found collision guard", ErrPathNotFound,
			http.StatusNotFound, "Path not found"},
		{"unknown cause stays 500", errors.New("boom"),
			http.StatusInternalServerError, "Internal server error"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			a.writeError(rec, tc.err)
			if rec.Code != tc.status {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tc.status, rec.Body.String())
			}
			var body struct {
				Errors struct {
					Detail string `json:"detail"`
				} `json:"errors"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("body %q: %v", rec.Body.String(), err)
			}
			if body.Errors.Detail != tc.detail {
				t.Errorf("detail = %q, want %q", body.Errors.Detail, tc.detail)
			}
		})
	}
}

$ sed -n '300,390p' /root/astrate-mule/internal/appengine/http.go
return
	}
	err := a.svc.PublishData(r.Context(), r.PathValue("realm"), r.PathValue("device"),
		r.PathValue("interface"), pathParam(r), value, nil)
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) deleteData(w http.ResponseWriter, r *http.Request) {
	err := a.svc.UnsetProperty(r.Context(), r.PathValue("realm"), r.PathValue("device"),
		r.PathValue("interface"), pathParam(r))
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) putDataByAlias(w http.ResponseWriter, r *http.Request) {
	var value json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &value); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if len(value) > maxValueBytes {
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
		return
	}
	err := a.svc.PublishDataByAlias(r.Context(), r.PathValue("realm"), r.PathValue("alias"),
		r.PathValue("interface"), pathParam(r), value, nil)
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) putDataInGroup(w http.ResponseWriter, r *http.Request) {
	var value json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &value); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if len(value) > maxValueBytes {
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
		return
	}
	err := a.svc.PublishDataInGroup(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device"),
		r.PathValue("interface"), pathParam(r), value, nil)
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) deleteDataByAlias(w http.ResponseWriter, r *http.Request) {
	err := a.svc.UnsetPropertyByAlias(r.Context(), r.PathValue("realm"), r.PathValue("alias"),
		r.PathValue("interface"), pathParam(r))
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteDataInGroup(w http.ResponseWriter, r *http.Request) {
	err := a.svc.UnsetPropertyInGroup(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device"),
		r.PathValue("interface"), pathParam(r))
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- groups -----------------------------------------------------------------

func (a *API) listGroups(w http.ResponseWriter, r *http.Request) {
	names, err := a.svc.ListGroups(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, names)
}

// groupBody is the POST /groups wire shape. Devices is a pointer so a MISSING

$ rg -rn "astarte_appengine_api.yaml" /root/astrate-mule --glob '!site-dist/**' -l | head -20; echo ---; cat /root/astrate-mule/docs/Makefile
/root/astrate-mule/docs/swagger-ui/index.html
/root/astrate-mule/docs/handoff/phase-1-memory.md
/root/astrate-mule/internal/swagger/swagger_test.go
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

→ Read internal/swagger/swagger_test.go

$ rg -n "    BadRequest:" -A 14 /root/astrate-mule/docs/api/astarte_appengine_api.yaml
1715:    BadRequest:
1716-      description: Bad Request — malformed or missing data envelope.
1717-      content:
1718-        application/json; charset=utf-8:
1719-          schema:
1720-            $ref: "#/components/schemas/ErrorDetail"
1721-          example:
1722-            errors:
1723-              detail: Bad Request
1724-
1725-    Unauthorized:
1726-      description: Unauthorized — missing or invalid JWT.
1727-      content:
1728-        application/json; charset=utf-8:
1729-          schema:

Now the three YAML edits.
✗ Edit docs/api/astarte_appengine_api.yaml failed
Error: Found multiple matches for oldString. Provide more surrounding context to make the match unique.

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -544,9 +544,10 @@
 operationId: deleteDataByAlias
 summary: Delete data at path by alias
 description: |
   Unsets the property at a specific path (properties interface only),
-  looked up by alias. Returns 204 on success.
+  looked up by alias. Returns 204 on success. A path matching no
+  endpoint mapping is a 400, not a 404.
 tags: [AppEngine - Data]
 security:
   - a_aea: []
 parameters:

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -567,8 +567,10 @@
       type: string
 responses:
   "204":
     description: Data deleted. No body returned.
+  "400":
+    $ref: "#/components/responses/BadRequest"
   "401":
     $ref: "#/components/responses/Unauthorized"
   "403":
     $ref: "#/components/responses/Forbidden"

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -796,9 +796,10 @@
 operationId: deleteData
 summary: Delete data at path
 description: |
   Unsets the property at a specific path (properties interface only).
-  Returns 204 on success.
+  Returns 204 on success. A path matching no endpoint mapping is a 400,
+  not a 404.
 tags: [AppEngine - Data]
 security:
   - a_aea: []
 parameters:
@@ -813,8 +814,10 @@
       type: string
 responses:
   "204":
     description: Data deleted. No body returned.
+  "400":
+    $ref: "#/components/responses/BadRequest"
   "401":
     $ref: "#/components/responses/Unauthorized"
   "403":
     $ref: "#/components/responses/Forbidden"

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -1358,9 +1358,10 @@
 operationId: deleteDataInGroup
 summary: Delete data at path in group
 description: |
   Unsets the property at a specific path (properties interface only),
-  within a group. Returns 204 on success.
+  within a group. Returns 204 on success. A path matching no endpoint
+  mapping is a 400, not a 404.
 tags: [AppEngine - Groups]
 security:
   - a_aea: []
 parameters:
@@ -1376,8 +1377,10 @@
       type: string
 responses:
   "204":
     description: Data deleted. No body returned.
+  "400":
+    $ref: "#/components/responses/BadRequest"
   "401":
     $ref: "#/components/responses/Unauthorized"
   "403":
     $ref: "#/components/responses/Forbidden"

Now the doc-contract test, following the existing pattern in `internal/swagger/swagger_test.go`.
← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -405,8 +405,54 @@
 		t.Error("components.responses.ValidationErrors does not carry the FieldErrorsDetail body")
 	}
 }
 
+// TestAppEngineDataDelete400Documented guards that the three data DELETE
+// operations document their 400: unsetting a path that matches no endpoint
+// mapping makes the trie lookup fail with engine.ErrPathNotFound
+// (UnsetServerProperty, internal/engine/serverdata.go), which writeError
+// answers as 400 "Endpoint not found" (internal/appengine/http.go) — the same
+// 400 the PUT/POST twins already document, because it is the same sentinel on
+// the same write path. Documenting only 404 would teach a generated client
+// that an unmatched path is a missing resource; the read path does use 404 for
+// that (appengine.ErrPathNotFound), and the two are deliberately distinct. This
+// is the documentation half of the behaviour pinned by
+// TestWriteErrorTaxonomy in internal/appengine.
+func TestAppEngineDataDelete400Documented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_appengine_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_appengine_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	const (
+		key = `        "400":`
+		ref = `          $ref: "#/components/responses/BadRequest"`
+	)
+
+	for _, op := range []string{"deleteDataByAlias", "deleteData", "deleteDataInGroup"} {
+		block := operationBlock(t, lines, op)
+		idx := -1
+		for i, l := range block {
+			if l == key {
+				idx = i
+				break
+			}
+		}
+		if idx < 0 {
+			t.Errorf("operation %s documents no 400 response", op)
+			continue
+		}
+		if idx+1 >= len(block) || block[idx+1] != ref {
+			t.Errorf("operation %s 400 is followed by %q, want the BadRequest $ref", op, block[idx+1])
+		}
+	}
+
+	if !containsLine(lines, "    BadRequest:") {
+		t.Error("components.responses defines no BadRequest response")
+	}
+}
+
 // operationBlock returns the lines of the operation with the given operationId,
 // up to the next operation, path, or the components section.
 func operationBlock(t *testing.T, lines []string, operationID string) []string {
 	t.Helper()

$ go test ./internal/swagger/ -run 'TestAppEngine' -v 2>&1 | tail -30
=== RUN   TestAppEngine403
--- PASS: TestAppEngine403 (0.00s)
=== RUN   TestAppEngineAddGroupDevice422Documented
--- PASS: TestAppEngineAddGroupDevice422Documented (0.00s)
=== RUN   TestAppEngineDataDelete400Documented
--- PASS: TestAppEngineDataDelete400Documented (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.041s

$ cp docs/api/astarte_appengine_api.yaml /tmp/opencode/api.yaml.bak && python3 - <<'EOF'
p='docs/api/astarte_appengine_api.yaml'
s=open(p).read()
old='''        "204":
          description: Data deleted. No body returned.
        "400":
          $ref: "#/components/responses/BadRequest"
        "401":'''
new='''        "204":
          description: Data deleted. No body returned.
        "401":'''
assert s.count(old)==3, s.count(old)
open(p,'w').write(s.replace(old,new))
EOF
go test ./internal/swagger/ -run TestAppEngineDataDelete400Documented 2>&1 | tail -12; cp /tmp/opencode/api.yaml.bak docs/api/astarte_appengine_api.yaml; echo "--- restored ---"; go test ./internal/swagger/ 2>&1 | tail -3
--- FAIL: TestAppEngineDataDelete400Documented (0.00s)
    swagger_test.go:442: operation deleteDataByAlias documents no 400 response
    swagger_test.go:442: operation deleteData documents no 400 response
    swagger_test.go:442: operation deleteDataInGroup documents no 400 response
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.025s
FAIL
--- restored ---
ok  	github.com/astrate-platform/astrate/internal/swagger	0.079s

$ cat .gitignore 2>/dev/null; echo "--- which mkdocs:"; which mkdocs || echo "no mkdocs"; python3 -c "import mkdocs; print('mkdocs module ok')" 2>&1 | tail -1
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
--- which mkdocs:
no mkdocs
ModuleNotFoundError: No module named 'mkdocs'

$ cat docs/requirements.txt; echo "--- try install:"; timeout 90 pip install -q -r docs/requirements.txt 2>&1 | tail -5; echo "exit=$?"; which mkdocs || echo "still no mkdocs"
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10
--- try install:
timeout: failed to run command ‘pip’: No such file or directory
exit=0
still no mkdocs

$ timeout 120 python3 -m pip install -q -r docs/requirements.txt 2>&1 | tail -5; which mkdocs || echo "still no mkdocs"
/usr/bin/python3: No module named pip
still no mkdocs

$ make -C docs sync 2>&1 | tail -5; echo "sync exit=$?"; ls docs/site/api/ && python3 -c "import yaml; print('pyyaml ok')" 2>&1 | tail -1
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
make: Leaving directory '/root/astrate-mule/docs'
sync exit=0
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
ModuleNotFoundError: No module named 'yaml'

$ rg -n "yaml" go.mod
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

$ cat > /tmp/opencode/yamlcheck.go <<'EOF'
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	fail := 0
	for _, path := range os.Args[1:] {
		b, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("%s: read: %v\n", path, err)
			fail++
			continue
		}
		var doc any
		if err := yaml.Unmarshal(b, &doc); err != nil {
			fmt.Printf("%s: PARSE FAIL: %v\n", path, err)
			fail++
			continue
		}
		var refs []string
		var walk func(any)
		walk = func(n any) {
			switch v := n.(type) {
			case map[string]any:
				for k, val := range v {
					if k == "$ref" {
						if s, ok := val.(string); ok {
							refs = append(refs, s)
						}
					}
					walk(val)
				}
			case []any:
				for _, e := range v {
					walk(e)
				}
			}
		}
		walk(doc)
		sort.Strings(refs)
		missing := 0
		for _, r := range refs {
			if !strings.HasPrefix(r, "#/") {
				fmt.Printf("%s: non-local $ref %q\n", path, r)
				missing++
				continue
			}
			cur := doc
			for _, seg := range strings.Split(strings.TrimPrefix(r, "#/"), "/") {
				m, ok := cur.(map[string]any)
				if !ok {
					missing++
					break
				}
				cur, ok = m[seg]
				if !ok {
					missing++
					break
				}
			}
			if missing > 0 {
				fmt.Printf("%s: UNRESOLVED $ref %q\n", path, r)
			}
		}
		if missing == 0 {
			fmt.Printf("%s: OK, parsed, %d $refs all resolve\n", path, len(refs))
		} else {
			fail++
		}
	}
	if fail > 0 {
		os.Exit(1)
	}
}
EOF
go run /tmp/opencode/yamlcheck.go docs/api/astrate_appengine_api.yaml docs/site/api/astarte_appengine_api.yaml; echo "exit=$?"
docs/api/astrate_appengine_api.yaml: read: open docs/api/astrate_appengine_api.yaml: no such file or directory
docs/site/api/astarte_appengine_api.yaml: OK, parsed, 346 $refs all resolve
exit=1
exit status 1

$ pwd; ls -la docs/api/astrate_appengine_api.yaml; go run /tmp/opencode/yamlcheck.go ./docs/api/astrate_appengine_api.yaml ./docs/site/api/astrate_appengine_api.yaml; echo "exit=$?"
/root/astrate-mule
ls: cannot access 'docs/api/astrate_appengine_api.yaml': No such file or directory
./docs/api/astrate_appengine_api.yaml: read: open ./docs/api/astrate_appengine_api.yaml: no such file or directory
./docs/site/api/astrate_appengine_api.yaml: read: open ./docs/site/api/astrate_appengine_api.yaml: no such file or directory
exit=1
exit status 1

$ ls -la docs/; echo ---; ls -la docs/api 2>&1; echo ---; ls -la docs/site 2>&1 | head
total 184
drwxr-xr-x  6 root root  4096 Sep  4 20:46 .
drwxr-xr-x 17 root root  4096 Sep 25 19:41 ..
-rw-r--r--  1 root root 23369 Sep 25 22:05 COMPATIBILITY.md
-rw-r--r--  1 root root 50906 Sep  4 20:46 DESIGN.md
-rw-r--r--  1 root root  6699 Jul 27 05:10 JSON-PAYLOAD-PROFILE.md
-rw-r--r--  1 root root   426 Sep  4 20:46 Makefile
-rw-r--r--  1 root root  6196 Jul 27 05:40 OPERATIONS.md
-rw-r--r--  1 root root 45928 Sep  4 20:46 ROADMAP.md
-rw-r--r--  1 root root  1336 Sep  4 20:46 UPSTREAM-EXPERIMENTAL.md
drwxr-xr-x  2 root root  4096 Sep  4 18:55 api
-rw-r--r--  1 root root   435 Jul 27 22:18 embed.go
drwxr-xr-x  2 root root  4096 Sep  4 20:46 handoff
-rw-r--r--  1 root root  1763 Sep  4 20:46 mkdocs.yml
-rw-r--r--  1 root root    40 Jul 27 05:10 requirements.txt
drwxr-xr-x  4 root root  4096 Sep 27 13:32 site
drwxr-xr-x  2 root root  4096 Jul 27 05:40 swagger-ui
---
total 168
drwxr-xr-x 2 root root  4096 Sep  4 18:55 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 59526 Sep 27 13:31 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 15547 Sep 25 21:48 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 19506 Sep 17 13:03 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 45387 Sep 25 22:07 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 14977 Sep 20 21:25 astrate_native_api.yaml
---
total 268
drwxr-xr-x 4 root root  4096 Sep 27 13:32 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 50906 Sep 27 13:32 DESIGN.md
-rw-r--r-- 1 root root 45928 Sep 27 13:32 ROADMAP.md
drwxr-xr-x 2 root root  4096 Sep 27 13:32 api
-rw-r--r-- 1 root root  3006 Sep  4 21:23 appengine-api.md
-rw-r--r-- 1 root root  6748 Jul 27 05:10 architecture.md
-rw-r--r-- 1 root root  6849 Sep  4 20:46 compatibility.md
-rw-r--r-- 1 root root  8501 Jul 27 05:40 configuration-reference.md

$ ls -b docs/api/; echo "--- retry:"; test -f docs/api/astarte_appengine_api.yaml && echo EXISTS || echo MISSING; stat -c '%n %s' docs/api/astarte_appengine_api.yaml
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
--- retry:
EXISTS
docs/api/astarte_appengine_api.yaml 59526

$ go run /tmp/opencode/yamlcheck.go docs/api/astarte_appengine_api.yaml docs/site/api/astarte_appengine_api.yaml; echo "exit=$?"; cmp docs/api/astarte_appengine_api.yaml docs/site/api/astarte_appengine_api.yaml && echo "site copy identical to source"
docs/api/astarte_appengine_api.yaml: OK, parsed, 346 $refs all resolve
docs/site/api/astarte_appengine_api.yaml: OK, parsed, 346 $refs all resolve
exit=0
site copy identical to source

$ make -C docs build 2>&1 | tail -3; echo "=== ensurepip attempt:"; timeout 60 python3 -m ensurepip --version 2>&1 | tail -2
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'
=== ensurepip attempt:
pip 25.1.1

$ timeout 150 python3 -m ensurepip >/dev/null 2>&1 && timeout 150 python3 -m pip install -q -r docs/requirements.txt 2>&1 | tail -3; which mkdocs || echo "no mkdocs"
no mkdocs

$ timeout 90 python3 -m pip install -r docs/requirements.txt 2>&1 | tail -6
/usr/bin/python3: No module named pip

$ git status --porcelain
M docs/api/astarte_appengine_api.yaml
 M internal/swagger/swagger_test.go

$ gofmt -l . 2>/dev/null; echo "--- gofmt done"; go vet ./... 2>&1 | tail -5; echo "--- vet exit=$?"
--- gofmt done
--- vet exit=0

$ go test ./... 2>&1 | rg -v "^ok|no test files" | head -20; echo "--- done"
--- done

$ timeout 900 go test -race ./internal/swagger/ ./internal/appengine/ 2>&1 | tail -5; echo "exit=$?"
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.055s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine	0.048s
FAIL
exit=0

$ timeout 300 go test -race ./internal/swagger/ 2>&1 | head -8
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.055s
FAIL

$ git diff --stat; echo ===; git diff docs/api/astarte_appengine_api.yaml
docs/api/astarte_appengine_api.yaml | 15 +++++++++---
 internal/swagger/swagger_test.go    | 46 +++++++++++++++++++++++++++++++++++++
 2 files changed, 58 insertions(+), 3 deletions(-)
===
diff --git a/docs/api/astarte_appengine_api.yaml b/docs/api/astarte_appengine_api.yaml
index 595da2d..d2ea100 100644
--- a/docs/api/astarte_appengine_api.yaml
+++ b/docs/api/astarte_appengine_api.yaml
@@ -545,7 +545,8 @@ paths:
       summary: Delete data at path by alias
       description: |
         Unsets the property at a specific path (properties interface only),
-        looked up by alias. Returns 204 on success.
+        looked up by alias. Returns 204 on success. A path matching no
+        endpoint mapping is a 400, not a 404.
       tags: [AppEngine - Data]
       security:
         - a_aea: []
@@ -567,6 +568,8 @@ paths:
       responses:
         "204":
           description: Data deleted. No body returned.
+        "400":
+          $ref: "#/components/responses/BadRequest"
         "401":
           $ref: "#/components/responses/Unauthorized"
         "403":
@@ -794,7 +797,8 @@ paths:
       summary: Delete data at path
       description: |
         Unsets the property at a specific path (properties interface only).
-        Returns 204 on success.
+        Returns 204 on success. A path matching no endpoint mapping is a 400,
+        not a 404.
       tags: [AppEngine - Data]
       security:
         - a_aea: []
@@ -811,6 +815,8 @@ paths:
       responses:
         "204":
           description: Data deleted. No body returned.
+        "400":
+          $ref: "#/components/responses/BadRequest"
         "401":
           $ref: "#/components/responses/Unauthorized"
         "403":
@@ -1353,7 +1359,8 @@ paths:
       summary: Delete data at path in group
       description: |
         Unsets the property at a specific path (properties interface only),
-        within a group. Returns 204 on success.
+        within a group. Returns 204 on success. A path matching no endpoint
+        mapping is a 400, not a 404.
       tags: [AppEngine - Groups]
       security:
         - a_aea: []
@@ -1371,6 +1378,8 @@ paths:
       responses:
         "204":
           description: Data deleted. No body returned.
+        "400":
+          $ref: "#/components/responses/BadRequest"
         "401":
           $ref: "#/components/responses/Unauthorized"
         "403":

$ tail -12 .mule/todo.md
- [x] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion Go — the device-delete leg added by docs-sync-rm-delete-device-async-operation-param (register a device, `DELETE /devices/{id}?async_operation=…` answers 204 and the row is gone on both values) compiles clean under `go vet -tags integration ./internal/realm/` but has never been executed, because the Pi has no database; until it runs, TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go is the only assertion of that leg that has actually been verified.
- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] forward-static-header-validation [auto]: `New` (internal/engine/forward/http.go:38-71) validates the URL and the method at boot precisely because of the rule its own comment states at 42-45 — "an unusable endpoint must fail here rather than surface once per delivery" — but `cfg.StaticHeaders` is stored unvalidated (http.go:68) and applied per delivery with `req.Header.Set` (http.go:110-112), and net/http checks header names and values only at write time. Measured on the real net/http path with a throwaway program: a name of `X Bad Name` or `X:Foo` fails per delivery with `net/http: invalid header field name`, and a value containing a newline fails with `net/http: invalid header field value for "X-Foo"` (a plain space in a value is legal). So one typo in `triggers.forward.static_headers` makes the process boot clean, accept triggers, and then count `astrate_engine_trigger_deliveries_total{outcome="failed"}` and log "custom trigger action failed" (internal/engine/triggers/actions.go:499) on every custom action forever. Validate in `New`: the RFC 7230 field-name grammar is exactly the token grammar `validMethod` (http.go:77-95) already implements, so reuse that predicate per name (rename it if you like) and add a CR/LF rejection per value. Add a table case to `TestNewRejectsBadConfig` (internal/engine/forward/http_test.go:271) asserting `New` returns an error for a bad name and a bad value and that no request is attempted; this test fails today because `New` returns nil error for both.
- [x] forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
- [x] forward-status-error-body [auto]: a non-2xx forward returns `forward: status %d` and nothing else (internal/engine/forward/http.go:121-123) after copying the **entire** response body into `io.Discard` with no limit (http.go:117-120) — the one thing an operator needs to debug a 500 from the bus is destroyed, and the sibling webhook request in the same codebase does bound its drain with `io.LimitReader(resp.Body, 1<<20)` (internal/engine/triggers/actions.go:626), so the two near-identical request paths have drifted on exactly that line. Read a bounded prefix (e.g. 512 bytes) of the body into the returned error and bound the discard the way actions.go:626 does. Add a case beside `TestStatusTable` (internal/engine/forward/http_test.go:164) where the handler answers 500 with body `{"error":"boom"}` and assert the error mentions both `500` and `boom`; it fails today because the error carries only the status code.
- [x] forward-envelope-bytes-test [auto]: `marshalEnvelope` (internal/engine/forward/envelope.go:19-29) has no direct test — `rg marshalEnvelope` finds exactly two callers (http.go:99, nats.go:45) and no test file — while envelope.go:6-8 makes the strong claim that "every Forwarder implementation in this package must produce byte-identical envelopes for the same inputs". Both suites check the shape only by unmarshalling into `bodyShape` (http_test.go:24-29), which cannot see key order or the exact bytes and would pass unchanged if a field were renamed on both sides. Add internal/engine/forward/envelope_test.go with a table pinning the exact output for (nil, nil) → `{"realm":"r","trigger":"t","action":null,"event":null}`, for the empty-non-nil pair the rule at 16-18 is about, and for a valued pair; container-free, no HTTP server needed.
- [x] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.
- [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml — the string `"403"` appears nowhere in the file and `components` defines no `Forbidden` response, yet every route is wrapped by `a.require = mw.RequireRealm(auth.ClaimAppEngine)` (internal/appengine/http.go:35), which answers `403 Forbidden` for a verified a_aea JWT whose grants do not authorize the method+path (`astarteapi.WriteForbidden`, internal/auth/middleware.go:90; `DetailForbidden`, pkg/astarteapi/envelope.go:41; pinned today by internal/appengine/http_test.go:184-185). Add a `Forbidden` response component copied from the realm-management one (docs/api/astarte_realm_management_api.yaml:1277-1285 — `ErrorDetail` schema, example `errors.detail: Forbidden`) and `$ref` it from all 32 operations, which is exactly the parity docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go:123-170) pins; add the same docs-half test for the AppEngine spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-patch-by-alias-409 [auto]: add the missing `409` response to `patchDeviceByAlias` in docs/api/astarte_appengine_api.yaml (yaml:241-291, currently 200/400/401/404/422/500) — the handler resolves the alias then calls the same `applyPatch` as the device- and group-scoped PATCH, and a rename to an alias already in use returns `ErrAliasAlreadyInUse` (internal/appengine/service.go:337) which `writeError` maps to `409 ConflictAliasInUse` (internal/appengine/http.go:650-651); the other two PATCH operations already document that response thanks to docs-sync-appengine-query-params-status and docs-sync-appengine-group-patch-status, so this is the last PATCH of the three with a hole. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml (yaml:939-980, currently 201/400/401/404/409/500) — a well-formed body whose `device_id` does not parse is rejected as a field error (`FieldErrors{"device_id": {"is not a valid device id"}}`, internal/appengine/service.go:631-634) and answered 422 with the `FieldErrorsDetail` body (`astarteapi.WriteFieldErrors`, internal/appengine/http.go:647), so `$ref` the existing `ValidationErrors` response (yaml:1707-1716) rather than `BadRequest` — the documented 400 only covers a body `DecodeData` cannot parse (internal/appengine/http.go:488-491). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data DELETE operations in docs/api/astarte_appengine_api.yaml (`deleteDataByAlias` yaml:519-547, `deleteData` yaml:756-784, `deleteDataInGroup` yaml:1285-1313 — each currently 204/401/404/405/500) — unsetting a path that matches no endpoint mapping returns `engine.ErrPathNotFound` (`UnsetServerProperty`, internal/engine/serverdata.go:241-244) which `writeError` answers as `400 "Endpoint not found"` (internal/appengine/http.go:679-680), the same 400 the PUT/POST twins already document via `BadRequest`; internal/appengine/writeerror_test.go pins the mapping. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 200/401/500) — a `from_token` that does not parse as a device id is refused with `ErrValidation: invalid cursor` (internal/appengine/service.go:150-155) which `writeError` maps to 422 with the flat `{"errors":{"detail":"invalid cursor"}}` body, i.e. the `ErrorDetail` shape and not `FieldErrorsDetail`, so `$ref` the same response the group twin uses for the identical parameter (`ValidationErrors`, yaml:934-935) only after fixing its example, or add a dedicated one. Extend the `from_token` parameter description (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ ls .mule/; echo ---; rg -n "todo.md|task.md" .mule/*.sh 2>/dev/null | head -20
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
recipes
reports
reviews
task.md
tasks
todo.md
---

$ ls .mule/hosts; rg -n "todo.md" -l . --glob '!.git/**' 2>/dev/null | head
.mule/hosts
./tools/mule.sh
./tools/mule-setup-pi.sh

$ rg -n "todo.md|\[x\]|\[!\]|task.md" tools/mule.sh | head -30
13:#   mule.sh next             run the first open task in .mule/todo.md
31:TODO="$MULE/todo.md"
261:# makes sed read the *script* as a filename, so the completed task was never marked [x] and
372:# A task is one `- [ ] slug: title` line in .mule/todo.md. If .mule/tasks/<slug>.md exists
455:  local taskfile="$MULE/task.md"
461:  local prompt='Read .mule/MULE.md, then read .mule/task.md, then do exactly that one task.
504:    # and marking it `- [!]` would retire a perfectly good task forever on a network blip.
557:    [ "$lineno" != 0 ] && sed_i "${lineno}s/^- \[ \]/- [x]/" "$TODO"
590:    sed_i "${lineno}s/^- \[ \]/- [!]/" "$TODO"
661:# The queue is the set of open issues labelled `mule`, plus the standing lines in todo.md.
662:# Issue-derived tasks deliberately do NOT get copied into todo.md: that file lives in git on
892:   tasks — it appends them to .mule/todo.md for you to approve
934:That recipe is a *proposal* job: its output is new task lines appended to .mule/todo.md,
940:and does not need the todo.md approval step. But then that is the *whole* deliverable for
965:  note "edit .mule/todo.md to approve/cut, then: tools/mule.sh loop"

→ Read internal/swagger/swagger_test.go [offset=405, limit=55]
Done.

```
Done: added the missing `400` -> BadRequest response to the three AppEngine data DELETE operations (`deleteDataByAlias`, `deleteData`, `deleteDataInGroup`) in docs/api/astarte_appengine_api.yaml, with a docs-half contract test `TestAppEngineDataDelete400Documented` (internal/swagger/swagger_test.go) that fails without the YAML change.
Files: docs/api/astarte_appengine_api.yaml, internal/swagger/swagger_test.go
Verified: `go vet ./...` -> pass; `go test ./...` -> pass (full suite); `gofmt -l .` -> clean; test proven to fail without the change (reverted the 3 blocks: 3 errors, then restored -> pass); YAML validity via throwaway Go program: both docs/api/ and docs/site/api/ copies parse and all 346 $refs resolve -> pass. `make -C docs build` could NOT be run: mkdocs is not installed on this host and there is no pip/network to install it (`mkdocs: No such file or directory`, `No module named pip`); `make -C docs sync` did run, and the staged site copy is byte-identical to the source. `go test -race` unavailable here — ThreadSanitizer "unsupported VMA range, Found 39 - Supported 48" (the known 39-bit-VMA host limitation); no concurrency was touched.
Unsure: the task line said the three ops currently list 204/401/404/405/500, but each also documented 403 (added earlier by docs-sync-ae-forbidden-403), so I inserted 400 in the 200/400/401/... position used by the PUT/POST twins. I also added one prose sentence per description ("A path matching no endpoint mapping is a 400, not a 404") because the read path genuinely answers 404 for the same situation (appengine.ErrPathNotFound) and the two are distinct sentinels; that wording is my judgement, not something the task required.
Follow-ups: none
```
