slug: docs-native-realm-name-pattern
verdict: done
at:  2ddbf7f
ran: 2026-10-05T18:14:28Z on DietPi in 174s


> build · big-pickle

I'll read the instruction files first.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read migrations/000002_metadata.up.sql
→ Read docs/api/astrate_native_api.yaml [offset=320, limit=80]
→ Read docs/api/astrate_native_api.yaml [limit=95, offset=397]

$ rg -n "a-z0-9" docs/ --glob '*.yaml' -B6 -A4 | head -80
docs/api/astarte_pairing_api.yaml-384-      required: [hw_id]
docs/api/astarte_pairing_api.yaml-385-      properties:
docs/api/astarte_pairing_api.yaml-386-        hw_id:
docs/api/astarte_pairing_api.yaml-387-          type: string
docs/api/astarte_pairing_api.yaml-388-          minLength: 22
docs/api/astarte_pairing_api.yaml-389-          maxLength: 22
docs/api/astarte_pairing_api.yaml:390:          pattern: '^[A-Za-z0-9_-]{22}$'
docs/api/astarte_pairing_api.yaml-391-          description: |
docs/api/astarte_pairing_api.yaml-392-            Device hardware ID: a 128-bit value in the 22-character unpadded
docs/api/astarte_pairing_api.yaml-393-            base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and
docs/api/astarte_pairing_api.yaml-394-            `_`, with no `=` padding. The standard-alphabet characters `+` and
--
docs/api/astarte_realm_management_api.yaml-939-        the same `404` and `{"errors":{"detail":"Device not found"}}` body as an
docs/api/astarte_realm_management_api.yaml-940-        unknown device.
docs/api/astarte_realm_management_api.yaml-941-      schema:
docs/api/astarte_realm_management_api.yaml-942-        type: string
docs/api/astarte_realm_management_api.yaml-943-        minLength: 22
docs/api/astarte_realm_management_api.yaml-944-        maxLength: 22
docs/api/astarte_realm_management_api.yaml:945:        pattern: '^[A-Za-z0-9_-]{22}$'
docs/api/astarte_realm_management_api.yaml-946-
docs/api/astarte_realm_management_api.yaml-947-    PolicyName:
docs/api/astarte_realm_management_api.yaml-948-      name: name
docs/api/astarte_realm_management_api.yaml-949-      in: path
--
docs/api/astarte_housekeeping_api.yaml-240-        as malformed: it is answered with the same `404` and
docs/api/astarte_housekeeping_api.yaml-241-        `{"errors":{"detail":"Not Found"}}` body as a realm that does not
docs/api/astarte_housekeeping_api.yaml-242-        exist, because the lookup matches the name verbatim with no format check
docs/api/astarte_housekeeping_api.yaml-243-        of its own.
docs/api/astarte_housekeeping_api.yaml-244-      schema:
docs/api/astarte_housekeeping_api.yaml-245-        type: string
docs/api/astarte_housekeeping_api.yaml:246:        pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_housekeeping_api.yaml-247-
docs/api/astarte_housekeeping_api.yaml-248-    AsyncOperation:
docs/api/astarte_housekeeping_api.yaml-249-      name: async_operation
docs/api/astarte_housekeeping_api.yaml-250-      in: query
--
docs/api/astarte_housekeeping_api.yaml-274-      required: [data]
docs/api/astarte_housekeeping_api.yaml-275-      properties:
docs/api/astarte_housekeeping_api.yaml-276-        data:
docs/api/astarte_housekeeping_api.yaml-277-          type: array
docs/api/astarte_housekeeping_api.yaml-278-          items:
docs/api/astarte_housekeeping_api.yaml-279-            type: string
docs/api/astarte_housekeeping_api.yaml:280:            pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_housekeeping_api.yaml-281-          description: List of realm names.
docs/api/astarte_housekeeping_api.yaml-282-
docs/api/astarte_housekeeping_api.yaml-283-    Realm:
docs/api/astarte_housekeeping_api.yaml-284-      type: object
docs/api/astarte_housekeeping_api.yaml-285-      required: [realm_name, jwt_public_key_pem]
docs/api/astarte_housekeeping_api.yaml-286-      properties:
docs/api/astarte_housekeeping_api.yaml-287-        realm_name:
docs/api/astarte_housekeeping_api.yaml-288-          type: string
docs/api/astarte_housekeeping_api.yaml:289:          pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_housekeeping_api.yaml-290-          description: Unique realm identifier.
docs/api/astarte_housekeeping_api.yaml-291-        jwt_public_key_pem:
docs/api/astarte_housekeeping_api.yaml-292-          type: string
docs/api/astarte_housekeeping_api.yaml-293-          description: PEM-encoded public key used to verify realm JWTs.
--
docs/api/astarte_housekeeping_api.yaml-305-    RealmCreate:
docs/api/astarte_housekeeping_api.yaml-306-      type: object
docs/api/astarte_housekeeping_api.yaml-307-      required: [realm_name, jwt_public_key_pem]
docs/api/astarte_housekeeping_api.yaml-308-      properties:
docs/api/astarte_housekeeping_api.yaml-309-        realm_name:
docs/api/astarte_housekeeping_api.yaml-310-          type: string
docs/api/astarte_housekeeping_api.yaml:311:          pattern: '^[a-z][a-z0-9]*$'
docs/api/astarte_housekeeping_api.yaml-312-          description: |
docs/api/astarte_housekeeping_api.yaml-313-            Unique realm identifier: a lowercase ASCII word starting with a
docs/api/astarte_housekeeping_api.yaml-314-            letter, so `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and
docs/api/astarte_housekeeping_api.yaml-315-            `test realm` are refused. The pattern is the realm row's own CHECK
--
docs/api/astarte_housekeeping_api.yaml-478-
docs/api/astarte_housekeeping_api.yaml-479-    ValidationError:
docs/api/astarte_housekeeping_api.yaml-480-      description: |
docs/api/astarte_housekeeping_api.yaml-481-        Unprocessable Entity — request body failed validation. The two
docs/api/astarte_housekeeping_api.yaml-482-        `realm_name` rejections partition the cases and neither overlaps: an
docs/api/astarte_housekeeping_api.yaml-483-        empty name is refused by the handler with `realm_name can't be blank`,
docs/api/astarte_housekeeping_api.yaml:484:        a non-empty name off the `^[a-z][a-z0-9]*$` pattern is refused by the
docs/api/astarte_housekeeping_api.yaml-485-        realm row's CHECK constraint with `realm_name is invalid`.
docs/api/astarte_housekeeping_api.yaml-486-      content:
docs/api/astarte_housekeeping_api.yaml-487-        application/json; charset=utf-8:
docs/api/astarte_housekeeping_api.yaml-488-          schema:

✗ Read docs/api/astrate_housekeeping_api.yaml failed [limit=30, offset=220]
Error: File not found: /root/astrate-mule/docs/api/astrate_housekeeping_api.yaml

$ rg -n "RealmName|/astrate/v1/\{realm\}/socket" docs/api/astrate_native_api.yaml
284:  /astrate/v1/{realm}/socket:
317:        - $ref: "#/components/parameters/RealmName"
370:        `/astrate/v1/{realm}/socket` above, whose example is snake_case
399:    RealmName:

$ ls docs/api/ && rg -n "RealmName:" -A22 docs/api/astrate_hk_api.yaml 2>/dev/null | head -40
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

$ rg -n "^    RealmName:" -A 20 docs/api/astrate_housekeeping_api.yaml
rg: docs/api/astrate_housekeeping_api.yaml: IO error for operation on docs/api/astrate_housekeeping_api.yaml: No such file or directory (os error 2)

$ ls -la docs/api/ ; rg --files docs/api/
total 180
drwxr-xr-x 2 root root  4096 Oct  4 19:17 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 61063 Oct  5 13:29 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 17437 Oct  4 19:36 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 21724 Sep 28 19:59 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 47623 Oct  2 22:09 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 16826 Oct  5 20:03 astrate_native_api.yaml
docs/api/astarte_pairing_api.yaml
docs/api/astarte_housekeeping_api.yaml
docs/api/astarte_appengine_api.yaml
docs/api/astrate_native_api.yaml
docs/api/astarte_realm_management_api.yaml

→ Read docs/api/astarte_housekeeping_api.yaml [limit=30, offset=222]

$ sed -n 35,60p internal/appengine/stream/ws.go; echo ---; sed -n 45,75p internal/appengine/channels/ws.go; echo ---; sed -n 55,80p internal/auth/middleware.go
bus     Bus
	require func(http.Handler) http.Handler
}

// NewAPI wires the socket to the engine bus and the a_ch middleware.
func NewAPI(bus Bus, mw *auth.Middleware) *API {
	return &API{bus: bus, require: mw.RequireRealm(auth.ClaimChannels)}
}

// Mount registers the socket route on mux.
func (a *API) Mount(mux *http.ServeMux) {
	mux.Handle("GET /astrate/v1/{realm}/socket", a.require(http.HandlerFunc(a.handle)))
}

// wireEvent is the JSON shape pushed to subscribers.
type wireEvent struct {
	Event     string    `json:"event"`
	Realm     string    `json:"realm"`
	DeviceID  string    `json:"device_id"`
	Interface string    `json:"interface,omitempty"`
	Path      string    `json:"path,omitempty"`
	Value     any       `json:"value,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func toWire(ev *engstream.Event) wireEvent {
---

// Mount registers the WebSocket handler on the given mux.
func (a *API) Mount(mux *http.ServeMux) {
	mux.Handle("GET /appengine/v1/socket/websocket", http.HandlerFunc(a.handle))
}

// handle upgrades the HTTP connection and runs the session loop.
func (a *API) handle(w http.ResponseWriter, r *http.Request) {
	realm := r.URL.Query().Get("realm")
	token := r.URL.Query().Get("token")

	if realm == "" || token == "" {
		_ = astarteapi.WriteUnauthorized(w)
		return
	}

	row, err := a.keys.GetRealmByName(r.Context(), realm)
	switch {
	case errors.Is(err, store.ErrNotFound):
		_ = astarteapi.WriteUnauthorized(w)
		return
	case err != nil:
		_ = astarteapi.WriteInternalServerError(w)
		return
	}

	tok, err := a.cache.Verify(token, row.JWTPublicKeysPEM)
	if err != nil {
		_ = astarteapi.WriteUnauthorized(w)
		return
	}
---
// claim and an Astrate compatibility claim — e.g. Flow accepts a_f (upstream)
// and a_rma (Astrate's original operator claim).
func (m *Middleware) RequireRealmAny(claims ...Claim) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			realm := r.PathValue("realm")
			if realm == "" || len(claims) == 0 {
				_ = astarteapi.WriteUnauthorized(w)
				return
			}

			row, err := m.keys.GetRealmByName(r.Context(), realm)
			switch {
			case errors.Is(err, store.ErrNotFound):
				// An unknown realm has no keys: unauthenticated, not 404
				// (no existence oracle on auth failures).
				_ = astarteapi.WriteUnauthorized(w)
				return
			case err != nil:
				_ = astarteapi.WriteInternalServerError(w)
				return
			}

			for _, claim := range claims {
				switch m.tryAuthorize(w, r, next, claim, row.JWTPublicKeysPEM, realm) {
				case authServed:

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -399,11 +399,20 @@
   RealmName:
     name: realm
     in: path
     required: true
-    description: The realm name.
+    description: |
+      The realm name: a lowercase ASCII word starting with a letter, so
+      `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm`
+      are not realm names — no realm can be created under one. This socket is a
+      read path though, and it authenticates before it looks anything up, so an
+      off-pattern name is not reported as malformed: it is answered `401`
+      `Unauthorized`, the same body as an unknown realm or a bad token, because
+      the realm lookup (`GetRealmByName` returning `store.ErrNotFound`) is
+      turned into a 401 with no existence oracle.
     schema:
       type: string
+      pattern: '^[a-z][a-z0-9]*$'
 
 schemas:
   HealthStatus:
     type: object

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -376,11 +376,20 @@
 parameters:
   - name: realm
     in: query
     required: true
-    description: The realm name.
+    description: |
+      The realm name: a lowercase ASCII word starting with a letter, so
+      `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm`
+      are not realm names — no realm can be created under one. This socket is
+      a read path though, and it authenticates before it looks anything up,
+      so an off-pattern name is not reported as malformed: it is answered
+      `401` `Unauthorized`, the same body as an unknown realm or a bad
+      token, because the realm lookup (`GetRealmByName` returning
+      `store.ErrNotFound`) is turned into a 401 with no existence oracle.
     schema:
       type: string
+      pattern: '^[a-z][a-z0-9]*$'
   - name: token
     in: query
     required: true
     description: Realm JWT for authentication.

$ cat docs/Makefile 2>/dev/null | head -40; ls docs/
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

$ rg -n "site" .gitignore; echo ---; ls docs/swagger-ui; echo ---; rg -n '\$ref: *["'"'"']?[^#]' docs/api/*.yaml | grep -v '#/' | head
29:/site-dist/
31:# Docs site copies (refreshed by docs/make sync)
32:/docs/site/DESIGN.md
33:/docs/site/ROADMAP.md
34:/docs/site/api/
35:/docs/site/swagger-ui/
---
index.html
---

$ make -C docs build 2>&1 | tail -30
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

$ python3 -c "import mkdocs" 2>&1; python3 -c "
import yaml,sys
for f in ['docs/api/astrate_native_api.yaml','docs/site/api/astrate_native_api.yaml']:
    d=yaml.safe_load(open(f))
    print(f,'ok', d['info']['title'], d['info']['version'])
d=yaml.safe_load(open('docs/api/astrate_native_api.yaml'))
import re
txt=open('docs/site/swagger-ui/index.html').read()
for m in set(re.findall(r'[\w./-]+\.ya?ml', txt)): print('ref:',m)
"
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import mkdocs
ModuleNotFoundError: No module named 'mkdocs'
Traceback (most recent call last):
  File "<string>", line 2, in <module>
    import yaml,sys
ModuleNotFoundError: No module named 'yaml'

$ go list -m all 2>/dev/null | grep -i yaml; ls $(go env GOMODCACHE)/gopkg.in 2>/dev/null
go.yaml.in/yaml/v2 v2.4.2
gopkg.in/yaml.v2 v2.4.0
gopkg.in/yaml.v3 v3.0.1
check.v1@v0.0.0-20161208181325-20d25e280405
check.v1@v1.0.0-20201130134442-10cb98267c6c
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1

$ mkdir -p /tmp/opencode/yamlprobe && cd /tmp/opencode/yamlprobe && cat > go.mod <<'EOF'
module yamlprobe

go 1.24

require gopkg.in/yaml.v3 v3.0.1
EOF
cat > main.go <<'EOF'
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
			fmt.Println(f, "READ-ERR", err)
			continue
		}
		var v map[string]any
		if err := yaml.Unmarshal(b, &v); err != nil {
			fmt.Println(f, "PARSE-ERR", err)
			continue
		}
		info := v["info"].(map[string]any)
		fmt.Printf("%s ok title=%q version=%v paths=%d\n", f, info["title"], info["version"], len(v["paths"].(map[string]any)))
	}
}
EOF
GOFLAGS=-mod=mod go mod tidy >/dev/null 2>&1; go run . /root/astrate-mule/docs/api/astrate_native_api.yaml /root/astrate-mule/docs/site/api/astrate_native_api.yaml /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
/root/astrate-mule/docs/api/astrate_native_api.yaml ok title="Astrate Native API" version=0.0.1 paths=12
/root/astrate-mule/docs/site/api/astrate_native_api.yaml ok title="Astrate Native API" version=0.0.1 paths=12
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml ok title="Astarte Housekeeping API" version=0.0.1 paths=2

$ sed -n 374,400p /root/astrate-mule/docs/api/astrate_native_api.yaml; echo '=== phoenix ==='; sed -n 386,405p /root/astrate-mule/docs/api/astrate_native_api.yaml
rather than failing, so every field arrives empty with no error.
      tags: [WebSocket]
      parameters:
        - name: realm
          in: query
          required: true
          description: |
            The realm name: a lowercase ASCII word starting with a letter, so
            `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm`
            are not realm names — no realm can be created under one. This socket is
            a read path though, and it authenticates before it looks anything up,
            so an off-pattern name is not reported as malformed: it is answered
            `401` `Unauthorized`, the same body as an unknown realm or a bad
            token, because the realm lookup (`GetRealmByName` returning
            `store.ErrNotFound`) is turned into a 401 with no existence oracle.
          schema:
            type: string
            pattern: '^[a-z][a-z0-9]*$'
        - name: token
          in: query
          required: true
          description: Realm JWT for authentication.
          schema:
            type: string
      responses:
        "101":
          description: WebSocket upgrade successful.
=== phoenix ===
            `401` `Unauthorized`, the same body as an unknown realm or a bad
            token, because the realm lookup (`GetRealmByName` returning
            `store.ErrNotFound`) is turned into a 401 with no existence oracle.
          schema:
            type: string
            pattern: '^[a-z][a-z0-9]*$'
        - name: token
          in: query
          required: true
          description: Realm JWT for authentication.
          schema:
            type: string
      responses:
        "101":
          description: WebSocket upgrade successful.
        "401":
          $ref: "#/components/responses/Unauthorized"
        "500":
          $ref: "#/components/responses/InternalServerError"

$ gofmt -l . ; go vet ./... 2>&1 | tail -5; go test ./... 2>&1 | rg -v "^ok|no test files" | tail -20
(no output)

$ go vet ./... >/dev/null 2>&1; echo "vet=$?"; go test ./... 2>&1 | rg -c "^ok"; go test ./... 2>&1 | rg "FAIL" ; echo "fail_grep=$?"
vet=0
31
fail_grep=1

$ sed -n 397,415p docs/api/astrate_native_api.yaml
type: string
      responses:
        "101":
          description: WebSocket upgrade successful.
        "401":
          $ref: "#/components/responses/Unauthorized"
        "500":
          $ref: "#/components/responses/InternalServerError"

components:
  parameters:
    RealmName:
      name: realm
      in: path
      required: true
      description: |
        The realm name: a lowercase ASCII word starting with a letter, so
        `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm`
        are not realm names — no realm can be created under one. This socket is a

Both edits are in and the docs YAML still parses; `mkdocs` itself isn't installed here so `make -C docs build` stops at the last step (the `sync` copy succeeded).

```
Done: Added `pattern: '^[a-z][a-z0-9]*$'` plus a 401-not-404 note to both socket realm parameters in docs/api/astrate_native_api.yaml (components.parameters.RealmName at yaml:407-420 and the Phoenix `realm` query parameter at yaml:376-400).
Files: docs/api/astrate_native_api.yaml
Verified: `make -C docs build` -> sync succeeded, `mkdocs build` unavailable (mkdocs not installed); YAML load check via throwaway Go probe (gopkg.in/yaml.v3) on docs/api/astrate_native_api.yaml and its docs/site/api copy -> both parse, 12 paths; `go vet ./...` -> pass; `go test ./...` -> 31 ok, 0 FAIL; `gofmt -l .` -> empty
Unsure: Whether "docs build" running to completion is expected on this machine — I substituted a YAML parse check of the exact file docs/site/swagger-ui serves. `docs/site/api/` is gitignored so the sync copy left no tracked diff. Wording of the 401 explanation is mine; it names `GetRealmByName`/`store.ErrNotFound` in both places, matching internal/auth/middleware.go:66-72 and internal/appengine/channels/ws.go:56-69 as the task stated.
Follow-ups: none
```
