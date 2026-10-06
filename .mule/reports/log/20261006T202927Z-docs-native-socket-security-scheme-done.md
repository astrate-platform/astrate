slug: docs-native-socket-security-scheme
verdict: done
at:  87c18cf
ran: 2026-10-06T20:29:27Z on DietPi in 299s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

$ sed -n '1,60p' docs/api/astrate_native_api.yaml; echo ---; rg -n "security|/socket" docs/api/astrate_native_api.yaml
openapi: "3.0.3"
info:
  title: Astrate Native API
  description: |
    Astrate-native operational and real-time API surface. Covers liveness,
    readiness, and Prometheus metrics (no auth), compat health and version
    endpoints for the Astarte Dashboard, and two WebSocket endpoints: the
    native live event stream and the Phoenix V2-compatible channels endpoint.
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
  # ── Health & Observability ──────────────────────────────────────────

  /astrate/v1/health:
    get:
      operationId: getHealth
      summary: Liveness probe
      description: |
        Returns `{"status": "ok"}` when the process is up and serving.
        Used by orchestrators (Docker, Kubernetes) as a liveness check.
      tags: [Observability]
      responses:
        "200":
          description: Service is alive.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/HealthStatus"
              example:
                status: ok

  /astrate/v1/readiness:
    get:
      operationId: getReadiness
      summary: Readiness probe
      description: |
        Runs every registered dependency check (database, broker). Returns 200
        when all checks pass, or 503 with per-check status when one or more
        dependencies are down.
      tags: [Observability]
      responses:
        "200":
          description: All dependencies are ready.
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/ReadinessStatus"
              example:
---
21:security: []
284:  /astrate/v1/{realm}/socket:
376:  /appengine/v1/socket/websocket:
397:        `/astrate/v1/{realm}/socket` above, whose example is snake_case

$ sed -n '1425,1460p' docs/api/astarte_appengine_api.yaml; echo ---; rg -n "securitySchemes|^security|a_aea" docs/api/astarte_appengine_api.yaml
$ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "404":
          $ref: "#/components/responses/NotFound"
        "405":
          $ref: "#/components/responses/MethodNotAllowed"
        "500":
          $ref: "#/components/responses/InternalServerError"

components:
  securitySchemes:
    a_aea:
      type: apiKey
      in: header
      name: Authorization
      description: |
        Realm-level JWT carrying the `a_aea` claim. Used to authenticate
        AppEngine operations (device data, groups).

  parameters:
    RealmName:
      name: realm
      in: path
      required: true
      description: The realm name.
      schema:
        type: string

    DeviceID:
      name: device
      in: path
      required: true
      description: The device hardware ID.
      schema:
        type: string
---
8:    All endpoints are guarded by a realm JWT carrying the `a_aea` claim.
21:security: []
41:        - a_aea: []
76:        - a_aea: []
148:        - a_aea: []
188:        - a_aea: []
216:        - a_aea: []
266:        - a_aea: []
299:        - a_aea: []
356:        - a_aea: []
391:        - a_aea: []
437:        - a_aea: []
488:        - a_aea: []
540:        - a_aea: []
593:        - a_aea: []
633:        - a_aea: []
665:        - a_aea: []
704:        - a_aea: []
750:        - a_aea: []
797:        - a_aea: []
845:        - a_aea: []
881:        - a_aea: []
910:        - a_aea: []
959:        - a_aea: []
991:        - a_aea: []
1038:        - a_aea: []
1087:        - a_aea: []
1116:        - a_aea: []
1166:        - a_aea: []
1192:        - a_aea: []
1223:        - a_aea: []
1263:        - a_aea: []
1310:        - a_aea: []
1358:        - a_aea: []
1407:        - a_aea: []
1436:  securitySchemes:
1437:    a_aea:
1442:        Realm-level JWT carrying the `a_aea` claim. Used to authenticate

$ wc -l docs/api/astrate_native_api.yaml; rg -n "^components:|^  [a-zA-Z]+:" docs/api/astrate_native_api.yaml | tail -30
536 docs/api/astrate_native_api.yaml
3:  title: Astrate Native API
4:  description: |
9:  version: "0.0.1"
10:  license:
433:components:
434:  parameters:
452:  schemas:
507:  responses:

→ Read docs/api/astrate_native_api.yaml [offset=270, limit=150]
→ Read docs/api/astrate_native_api.yaml [offset=419]

$ sed -n '40,80p;120,200p' internal/auth/middleware.go; echo === WS ===; sed -n '35,60p' internal/appengine/stream/ws.go; echo === CH ===; sed -n '45,60p' internal/appengine/channels/ws.go
func NewMiddleware(keys KeySource) *Middleware {
	return &Middleware{keys: keys, cache: NewCache(DefaultCacheSize)}
}

// RequireRealm guards a realm-scoped route (path pattern must carry a
// {realm} segment): it resolves the realm's JWT public keys, verifies the
// bearer token, and matches the claim's authorization strings against the
// method and the path relative to the realm base.
func (m *Middleware) RequireRealm(claim Claim) func(http.Handler) http.Handler {
	return m.RequireRealmAny(claim)
}

// RequireRealmAny guards a realm-scoped route accepting ANY of the given
// claims (OR-ed): the first claim whose authorization strings match the
// method/path grants access. Used by surfaces that honour both their upstream
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
type authResult int

const (
	authServed    authResult = iota // request authorized and next served
	authFailed                      // authentication failed (401 written); stop
	authForbidden                   // claim did not grant the request; try next
)

// tryAuthorize is authorize without the terminal 403: it classifies the
// attempt. Callers write their own forbidden response on authForbidden.
func (m *Middleware) tryAuthorize(w http.ResponseWriter, r *http.Request, next http.Handler, claim Claim, keysPEM []string, base string) authResult {
	tokenString, ok := bearerToken(r)
	if !ok {
		_ = astarteapi.WriteUnauthorized(w)
		return authFailed
	}

	tok, err := m.cache.Verify(tokenString, keysPEM)
	if err != nil {
		_ = astarteapi.WriteUnauthorized(w)
		return authFailed
	}

	authPath, ok := RelativePath(r.URL.Path, base)
	if !ok {
		// Upstream parity: a path the authorizer cannot anchor is an
		// authorization failure (403), not an authentication one.
		return authForbidden
	}
	if !tok.Authorizes(claim, r.Method, authPath) {
		return authForbidden
	}

	ctx := context.WithValue(r.Context(), tokenContextKey{}, tok)
	next.ServeHTTP(w, r.WithContext(ctx))
	return authServed
}

// RelativePath computes the authorization path with upstream parity (Astarte's
// GuardianAuthorizePath plug): split the URL path into segments, drop
// everything up to and including the first segment equal to base, and join
// the rest with "/". ok is false when base does not appear in the path.
//
// Example: RelativePath("/pairing/v1/test/agent/devices", "test") returns
// ("agent/devices", true).
func RelativePath(urlPath, base string) (string, bool) {
	segments := strings.Split(strings.Trim(urlPath, "/"), "/")
	for i, s := range segments {
		if s == base {
			return strings.Join(segments[i+1:], "/"), true
		}
	}
	return "", false
}

// bearerToken extracts the credential from the Authorization header.
// Upstream parity (~r/bearer\:?\s+(.*)$/i): the scheme is case-insensitive
// "Bearer", optionally followed by a colon, then whitespace and the token.
func bearerToken(r *http.Request) (string, bool) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", false
	}
	scheme, rest, found := strings.Cut(header, " ")
	if !found {
		return "", false
	}
	scheme = strings.TrimSuffix(scheme, ":")
	if !strings.EqualFold(scheme, "bearer") {
		return "", false
	}
	token := strings.TrimSpace(rest)
	if token == "" {
		return "", false
	}
	return token, true
}
=== WS ===
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
=== CH ===

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

$ cat docs/Makefile 2>/dev/null || ls docs; echo ===; ls docs
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
===
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

$ rg -rn "astrate_native_api|securitySchemes|openapi" --glob '!docs/api/*' -l . | head -20
./docs/swagger-ui/index.html
./docs/handoff/phase-2-handoff.md
./docs/handoff/phase-2-memory.md
./docs/handoff/phase-1-memory.md
./internal/swagger/swagger_test.go

→ Read internal/swagger/swagger_test.go

$ cat internal/swagger/*.go | head -60; echo ===; cat docs/swagger-ui/index.html | head -60
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
===
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Astrate API Documentation</title>
  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5.32.11/favicon-32x32.png">
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui.css">
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #fafafa; }
    #header {
      background: #1b1b1b; color: #fff; padding: 16px 24px;
      display: flex; align-items: center; gap: 24px; flex-wrap: wrap;
    }
    #header h1 { font-size: 18px; font-weight: 600; white-space: nowrap; }
    #header .subtitle { font-size: 13px; color: #a0a0a0; }
    #nav {
      background: #fff; border-bottom: 1px solid #e0e0e0; padding: 12px 24px;
      display: flex; align-items: center; gap: 12px; flex-wrap: wrap;
    }
    #nav label { font-size: 13px; font-weight: 600; color: #555; }
    #nav select {
      padding: 6px 12px; border: 1px solid #ccc; border-radius: 4px;
      font-size: 13px; background: #fff; min-width: 280px;
    }
    #nav .server-label { font-size: 13px; color: #555; margin-left: 16px; }
    #nav input[type="text"] {
      padding: 6px 12px; border: 1px solid #ccc; border-radius: 4px;
      font-size: 13px; width: 320px;
    }
    #swagger-ui { max-width: 100%; }
    .topbar { display: none; }
  </style>
</head>
<body>
  <div id="header">
    <h1>Astrate API</h1>
    <span class="subtitle">OpenAPI 3.0 Reference</span>
  </div>
  <div id="nav">
    <label for="spec-select">Specification:</label>
    <select id="spec-select">
      <option value="../api/astarte_housekeeping_api.yaml">Housekeeping API — Realm management, a_ha auth</option>
      <option value="../api/astarte_pairing_api.yaml">Pairing API — Device pairing flows A/B/C, a_pa auth</option>
      <option value="../api/astarte_realm_management_api.yaml">Realm Management API — Interfaces, triggers, policies, a_rma auth</option>
      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
      <option value="../api/astrate_native_api.yaml">Astrate Native API — Health, readiness, metrics, WebSocket</option>
    </select>
    <span class="server-label">Server:</span>
    <input type="text" id="server-url" value="http://localhost:8080" placeholder="e.g. http://localhost:8080">
  </div>
  <div id="swagger-ui"></div>

  <script src="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui-standalone-preset.js"></script>
  <script>
    const specSelect = document.getElementById('spec-select');
    const serverUrl = document.getElementById('server-url');
    const uiContainer = document.getElementById('swagger-ui');

$ rg -n "^func (operationBlock|operationDescription|componentBlock|containsLine)" -A 40 internal/swagger/swagger_test.go | sed -n '1,140p'
1315:func operationDescription(t *testing.T, block []string) string {
1316-	t.Helper()
1317-	var text []string
1318-	for i, l := range block {
1319-		if l == "      description: |" {
1320-			for _, b := range block[i+1:] {
1321-				if strings.TrimSpace(b) == "" {
1322-					continue
1323-				}
1324-				if !strings.HasPrefix(b, "        ") {
1325-					break
1326-				}
1327-				text = append(text, strings.TrimSpace(b))
1328-			}
1329-			return strings.Join(text, " ")
1330-		}
1331-		if v, ok := strings.CutPrefix(l, "      description: "); ok {
1332-			return strings.Trim(strings.TrimSpace(v), `"'`)
1333-		}
1334-	}
1335-	t.Fatal("operation declares no description")
1336-	return ""
1337-}
1338-
1339-// operationBlock returns the lines of the operation with the given operationId,
1340-// up to the next operation, path, or the components section.
1341:func operationBlock(t *testing.T, lines []string, operationID string) []string {
1342-	t.Helper()
1343-	marker := "      operationId: " + operationID
1344-	start := -1
1345-	for i, l := range lines {
1346-		if l == marker {
1347-			start = i
1348-			break
1349-		}
1350-	}
1351-	if start < 0 {
1352-		t.Fatalf("spec declares no operation %q", operationID)
1353-	}
1354-	endRe := regexp.MustCompile(`^      operationId: |^  \S|^components:`)
1355-	for i := start + 1; i < len(lines); i++ {
1356-		if endRe.MatchString(lines[i]) {
1357-			return lines[start+1 : i]
1358-		}
1359-	}
1360-	return lines[start+1:]
1361-}
1362-
1363-// componentBlock returns the lines of the components entry starting with the
1364-// given key line, up to the next entry at the same or shallower indentation.
1365:func componentBlock(t *testing.T, lines []string, key string) []string {
1366-	t.Helper()
1367-	start := -1
1368-	for i, l := range lines {
1369-		if l == key {
1370-			start = i
1371-			break
1372-		}
1373-	}
1374-	if start < 0 {
1375-		t.Fatalf("spec declares no components entry %q", key)
1376-	}
1377-	for i := start + 1; i < len(lines); i++ {
1378-		if strings.TrimSpace(lines[i]) == "" {
1379-			continue
1380-		}
1381-		if !strings.HasPrefix(lines[i], "      ") {
1382-			return lines[start+1 : i]
1383-		}
1384-	}
1385-	return lines[start+1:]
1386-}
1387-
1388:func containsLine(lines []string, want string) bool {
1389-	for _, l := range lines {
1390-		if l == want {
1391-			return true
1392-		}
1393-	}
1394-	return false
1395-}
1396-
1397-func embeddedYAMLFilenames() ([]string, error) {
1398-	var names []string
1399-	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
1400-		if err != nil {
1401-			return err
1402-		}
1403-		if d.IsDir() {
1404-			return nil
1405-		}
1406-		if strings.HasSuffix(path, ".yaml") {
1407-			names = append(names, strings.TrimPrefix(path, "api/"))
1408-		}
1409-		return nil
1410-	})
1411-	sort.Strings(names)
1412-	return names, err
1413-}
1414-
1415-func get(t *testing.T, url string) string {
1416-	t.Helper()
1417-	resp, err := http.Get(url)
1418-	if err != nil {
1419-		t.Fatalf("GET %s: %v", url, err)
1420-	}
1421-	defer resp.Body.Close()
1422-
1423-	if resp.StatusCode != http.StatusOK {
1424-		t.Fatalf("GET %s: status = %d, want %d", url, resp.StatusCode, http.StatusOK)
1425-	}
1426-	b, err := io.ReadAll(resp.Body)
1427-	if err != nil {
1428-		t.Fatalf("reading %s: %v", url, err)

$ rg -n "security" internal/swagger/swagger_test.go; echo === ; rg -rn "yaml" --glob '*.go' -l | head; echo ===; rg -n "gopkg.in/yaml|sigs.k8s.io/yaml" go.mod
===
docs/embed.go
internal/swagger/swagger.go
internal/swagger/swagger_test.go
===
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

$ cat docs/mkdocs.yml
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

nav:
  - Home: index.md
  - Architecture: architecture.md
  - Data Modeling: data-modeling.md
  - MQTT Protocol: mqtt-protocol.md
  - Pairing & Security: pairing-and-security.md
  - Payload Formats: payload-formats.md
  - JSON Payload Profile: json-payload-profile.md
  - AppEngine API: appengine-api.md
  - Realm Management API: realm-management-api.md
  - Housekeeping API: housekeeping-api.md
  - Interface Schema: interface-schema.md
  - Triggers: triggers.md
  - Deployment: deployment.md
  - Operations: operations.md
  - Observability: observability.md
  - Compatibility: compatibility.md
  - Configuration Reference: configuration-reference.md
  - Troubleshooting: troubleshooting.md
  - Contributing: contributing.md
  - Migration from Astarte: migration-from-astarte.md
  - API Explorer: swagger.md
  - Quickstart: quickstart.md
  - Examples:
    - Pokémon Agent: pokemon-agent.md

$ go version; which mkdocs python3; python3 -c "import yaml; print('pyyaml ok')" 2>&1 | tail -1; cat docs/requirements.txt; echo ===; git check-ignore -v docs/site site-dist 2>&1; git status --porcelain | head
go version go1.26.7 linux/arm64
/usr/bin/python3
ModuleNotFoundError: No module named 'yaml'
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10
===
.gitignore:29:/site-dist/	site-dist

$ git ls-files docs/site | head; echo "---count:"; git ls-files docs/site | wc -l; echo ===; git ls-files docs/api
docs/site/appengine-api.md
docs/site/architecture.md
docs/site/compatibility.md
docs/site/configuration-reference.md
docs/site/contributing.md
docs/site/data-modeling.md
docs/site/deployment.md
docs/site/housekeeping-api.md
docs/site/index.md
docs/site/interface-schema.md
---count:
23
===
docs/api/astarte_appengine_api.yaml
docs/api/astarte_housekeeping_api.yaml
docs/api/astarte_pairing_api.yaml
docs/api/astarte_realm_management_api.yaml
docs/api/astrate_native_api.yaml

$ python3 -m pip --version 2>&1 | tail -2; ls /usr/bin | rg -i "pip|mkdocs" ; echo ===; go env GOFLAGS GOPATH GOMODCACHE
/usr/bin/python3: No module named pip
pipewire
pipewire-aes67
pipewire-avb
===

/root/go
/root/go/pkg/mod

$ git ls-files docs/site | rg "api/|swagger"; echo "---"; git status --porcelain --ignored docs/site | head; echo ===; ls ~/.local/bin 2>/dev/null; ls /usr/local/bin 2>/dev/null | head -20; find / -name "mkdocs*" -maxdepth 6 -not -path "*/proc/*" 2>/dev/null | head
docs/site/swagger.md
---
!! docs/site/DESIGN.md
!! docs/site/ROADMAP.md
!! docs/site/api/
!! docs/site/swagger-ui/
===
claude
heif-convert
heif-dec
heif-enc
heif-info
heif-thumbnailer
hermes
net-watchdog.sh
node
npm
npx
opencode-telegram
pihole
/usr/share/fish/completions/mkdocs.fish
/root/astrate-mule/docs/mkdocs.yml
/root/astrate/docs/mkdocs.yml
/root/astrate-survey/docs/mkdocs.yml
/tmp/mkvenv/bin/mkdocs
/tmp/mkvenv/bin/mkdocs-get-deps
/tmp/mkvenv/lib/python3.13/site-packages/mkdocs_material-9.7.7.dist-info
/tmp/mkvenv/lib/python3.13/site-packages/mkdocs-1.6.1.dist-info
/tmp/mkvenv/lib/python3.13/site-packages/mkdocs
/tmp/mkvenv/lib/python3.13/site-packages/mkdocs_get_deps-0.2.2.dist-info

$ sed -n '30,50p' docs/api/astarte_appengine_api.yaml; echo "=== housekeeping op ==="; rg -n "^  /|^    (get|post|put|delete|patch):|^      security|^      tags|security:" docs/api/astarte_housekeeping_api.yaml | sed -n '1,25p'
description: |
        Returns the emulated upstream API compatibility level. The Astarte
        Dashboard gates feature UI on this value.

        The authorization path is exactly `version` — the realm segment is
        stripped before the JWT grant is matched — so the token needs a grant
        for GET on `version`. An unknown realm is answered 401, not 404: a
        realm with no keys is indistinguishable from an unauthenticated
        caller.
      tags: [AppEngine - Misc]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
      responses:
        "200":
          description: API compatibility version.
          content:
            application/json; charset=utf-8:
              schema:
                type: object
=== housekeeping op ===
23:security: []
26:  /housekeeping/v1/realms:
27:    get:
31:      tags: [Housekeeping]
32:      security:
50:    post:
57:      tags: [Housekeeping]
58:      security:
108:  /housekeeping/v1/realms/{realm}:
109:    get:
113:      tags: [Housekeeping]
114:      security:
142:    patch:
155:      tags: [Housekeeping]
156:      security:
194:    delete:
198:      tags: [Housekeeping]
199:      security:

$ rg -n "securitySchemes" -A 12 docs/api/astarte_housekeeping_api.yaml docs/api/astarte_pairing_api.yaml
docs/api/astarte_housekeeping_api.yaml:221:  securitySchemes:
docs/api/astarte_housekeeping_api.yaml-222-    a_ha:
docs/api/astarte_housekeeping_api.yaml-223-      type: apiKey
docs/api/astarte_housekeeping_api.yaml-224-      in: header
docs/api/astarte_housekeeping_api.yaml-225-      name: Authorization
docs/api/astarte_housekeeping_api.yaml-226-      description: |
docs/api/astarte_housekeeping_api.yaml-227-        Instance-level JWT carrying the `a_ha` claim. Used to authenticate
docs/api/astarte_housekeeping_api.yaml-228-        housekeeping operations (realm CRUD).
docs/api/astarte_housekeeping_api.yaml-229-
docs/api/astarte_housekeeping_api.yaml-230-  parameters:
docs/api/astarte_housekeeping_api.yaml-231-    RealmName:
docs/api/astarte_housekeeping_api.yaml-232-      name: realm
docs/api/astarte_housekeeping_api.yaml-233-      in: path
--
docs/api/astarte_pairing_api.yaml:341:  securitySchemes:
docs/api/astarte_pairing_api.yaml-342-    a_pa:
docs/api/astarte_pairing_api.yaml-343-      type: apiKey
docs/api/astarte_pairing_api.yaml-344-      in: header
docs/api/astarte_pairing_api.yaml-345-      name: Authorization
docs/api/astarte_pairing_api.yaml-346-      description: |
docs/api/astarte_pairing_api.yaml-347-        Realm-level JWT carrying the `a_pa` claim. Used to authenticate
docs/api/astarte_pairing_api.yaml-348-        pairing agent endpoints (register/unregister devices).
docs/api/astarte_pairing_api.yaml-349-
docs/api/astarte_pairing_api.yaml-350-    bearerSecret:
docs/api/astarte_pairing_api.yaml-351-      type: http
docs/api/astarte_pairing_api.yaml-352-      scheme: bearer
docs/api/astarte_pairing_api.yaml-353-      description: |

$ rg -n "operationId:" docs/api/astrate_native_api.yaml
28:      operationId: getHealth
46:      operationId: getReadiness
79:      operationId: getMetrics
116:      operationId: getAppEngineHealth
144:      operationId: getRealmManagementHealth
172:      operationId: getPairingHealth
202:      operationId: getAppEngineVersion
222:      operationId: getRealmManagementVersion
242:      operationId: getPairingVersion
262:      operationId: getHousekeepingVersion
286:      operationId: nativeWebSocket
378:      operationId: phoenixWebSocket

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -318,8 +318,10 @@
   **SSE protocol:** each event is a `data:` frame followed by a blank line.
 
   Authentication: `a_ch` realm JWT via Authorization header.
 tags: [WebSocket]
+security:
+  - a_ch: []
 parameters:
   - $ref: "#/components/parameters/RealmName"
   - name: device_id
     in: query

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -432,8 +432,23 @@
         "500":
           $ref: "#/components/responses/InternalServerError"
 
 components:
+  securitySchemes:
+    a_ch:
+      type: apiKey
+      in: header
+      name: Authorization
+      description: |
+        Realm-level JWT carrying the `a_ch` claim, sent as a bearer credential
+        in the Authorization header: the scheme is a case-insensitive `Bearer`
+        optionally followed by a colon, then whitespace and the token (upstream
+        parity, `~r/bearer\:?\s+(.*)$/i`). A missing, empty or differently
+        spelled header, a bad token and an unknown realm are all answered 401
+        `{"errors":{"detail":"Unauthorized"}}` — an unknown realm is 401 and
+        not 404, because a realm with no keys is indistinguishable from an
+        unauthenticated caller (no existence oracle).
+
   parameters:
     RealmName:
       name: realm
       in: path

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -439,15 +439,15 @@
     in: header
     name: Authorization
     description: |
       Realm-level JWT carrying the `a_ch` claim, sent as a bearer credential
-      in the Authorization header: the scheme is a case-insensitive `Bearer`
-      optionally followed by a colon, then whitespace and the token (upstream
-      parity, `~r/bearer\:?\s+(.*)$/i`). A missing, empty or differently
-      spelled header, a bad token and an unknown realm are all answered 401
+      in the Authorization header: a case-insensitive `Bearer` optionally
+      followed by a colon, then whitespace and the token (upstream parity,
+      `~r/bearer\:?\s+(.*)$/i`). A missing, empty or differently spelled
+      header, a bad token and an unknown realm are all answered 401
       `{"errors":{"detail":"Unauthorized"}}` — an unknown realm is 401 and
       not 404, because a realm with no keys is indistinguishable from an
-      unauthenticated caller (no existence oracle).
+      unauthenticated caller.
 
 parameters:
   RealmName:
     name: realm

$ rg -n "^security|^components:|^  securitySchemes|^    a_ch:|^      security:|^        - a_ch: \[\]|^      tags: \[WebSocket\]" docs/api/astrate_native_api.yaml
21:security: []
321:      tags: [WebSocket]
322:      security:
323:        - a_ch: []
404:      tags: [WebSocket]
435:components:
436:  securitySchemes:
437:    a_ch:

$ /tmp/mkvenv/bin/python3 - <<'EOF'
import yaml, sys
for f in ["docs/api/astrate_native_api.yaml"]:
    d = yaml.safe_load(open(f))
    print(f, "parsed ok")
    print("root security:", d.get("security"))
    print("schemes:", list(d["components"].get("securitySchemes", {})))
    sec = {k: v.get("get", {}).get("security") for k, v in d["paths"].items()}
    for k, v in sec.items():
        print(" ", k, "->", v)
    s = d["components"]["securitySchemes"]["a_ch"]
    print("a_ch:", s["type"], s["in"], s["name"])
EOF
docs/api/astrate_native_api.yaml parsed ok
root security: []
schemes: ['a_ch']
  /astrate/v1/health -> None
  /astrate/v1/readiness -> None
  /astrate/v1/metrics -> None
  /appengine/health -> None
  /realmmanagement/health -> None
  /pairing/health -> None
  /appengine/version -> None
  /realmmanagement/version -> None
  /pairing/version -> None
  /housekeeping/version -> None
  /astrate/v1/{realm}/socket -> [{'a_ch': []}]
  /appengine/v1/socket/websocket -> None
a_ch: apiKey header Authorization

→ Read internal/swagger/swagger_test.go [offset=700, limit=20]

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -708,8 +708,76 @@
 		t.Errorf("getHousekeepingVersion description still promises the %s endpoints, which housekeeping does not serve", twin)
 	}
 }
 
+// TestNativeSocketSecurityDocumented guards that the native spec declares the
+// authentication of /astrate/v1/{realm}/socket as an OpenAPI security
+// requirement instead of prose in the description alone. The socket is wrapped
+// by mw.RequireRealm(auth.ClaimChannels) (internal/appengine/stream/ws.go) and
+// bearerToken reads the credential only from the Authorization header
+// (internal/auth/middleware.go): a request with no usable header is answered
+// 401 {"errors":{"detail":"Unauthorized"}} before any handler runs, so a
+// client generated from a spec with no security requirement on the operation
+// never sends the header and can only ever collect that 401. The scheme
+// mirrors upstream's a_aea (astarte_appengine_api.yaml): apiKey in the
+// Authorization header.
+//
+// The other half is asserted too: the document-level `security: []` and the
+// absence of an operation-level `security` on every other operation are what
+// document health, readiness, metrics, the compat version endpoints and the
+// Phoenix socket as unauthenticated — the Phoenix twin needs no header scheme
+// because it already declares `?token=` as a required query parameter, and a
+// document-level requirement would teach a generated client to send a
+// credential those handlers never read.
+func TestNativeSocketSecurityDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astrate_native_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astrate_native_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	if !containsLine(lines, "security: []") {
+		t.Error("root `security: []` is gone; it is what documents the observability endpoints as unauthenticated")
+	}
+
+	scheme := componentBlock(t, lines, "    a_ch:")
+	for _, want := range []string{
+		"      type: apiKey",
+		"      in: header",
+		"      name: Authorization",
+	} {
+		if !containsLine(scheme, want) {
+			t.Errorf("securityScheme a_ch is missing line %q", want)
+		}
+	}
+	desc := strings.Join(strings.Fields(strings.Join(scheme, "\n")), " ")
+	for _, want := range []string{
+		"a_ch",
+		"case-insensitive `Bearer` optionally followed by a colon",
+		"an unknown realm is 401 and not 404",
+	} {
+		if !strings.Contains(desc, want) {
+			t.Errorf("securityScheme a_ch description does not say %q", want)
+		}
+	}
+
+	sock := operationBlock(t, lines, "nativeWebSocket")
+	if !containsLine(sock, "      security:") || !containsLine(sock, "        - a_ch: []") {
+		t.Errorf("nativeWebSocket declares no `security: - a_ch: []`; a generated client sends no Authorization header")
+	}
+
+	for _, op := range []string{
+		"getHealth", "getReadiness", "getMetrics",
+		"getAppEngineHealth", "getRealmManagementHealth", "getPairingHealth",
+		"getAppEngineVersion", "getRealmManagementVersion", "getPairingVersion",
+		"getHousekeepingVersion", "phoenixWebSocket",
+	} {
+		if containsLine(operationBlock(t, lines, op), "      security:") {
+			t.Errorf("%s declares an operation-level security requirement, but it is documented as unauthenticated", op)
+		}
+	}
+}
+
 // protobufAccept is the Accept header a scraper sends to ask for the
 // delimited binary MetricFamily format instead of the text exposition format.
 const protobufAccept = "application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited"
 

$ gofmt -l . 2>&1 | head; echo "gofmt-exit:$?"; go vet ./... 2>&1 | tail -5; echo "vet-exit:$?"
gofmt-exit:0
vet-exit:0

$ go test -race ./... 2>&1 | tail -40
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/housekeeping	0.058s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/httpx	0.077s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/observability	0.026s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/pairing	0.033s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/pairing/ca	0.042s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/realm	0.034s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/store	0.028s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.076s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/testutil	0.053s
?   	github.com/astrate-platform/astrate/migrations	[no test files]
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/astarteapi	0.066s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/deviceid	0.072s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/interfaceschema	0.041s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/payload	0.036s
FAIL

$ go test ./... 2>&1 | tail -50
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
ok  	github.com/astrate-platform/astrate/examples/flow-container-echo	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	0.075s
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cp docs/api/astrate_native_api.yaml /tmp/native.yaml.bak && python3 - <<'EOF'
import re
p="docs/api/astrate_native_api.yaml"
s=open(p).read()
s=s.replace("""      security:
        - a_ch: []
""","",1)
i=s.index("components:\n")
j=s.index("  parameters:",i)
s=s[:i]+"components:\n"+s[j:]
open(p,"w").write(s)
EOF
go test ./internal/swagger/ -run TestNativeSocketSecurityDocumented 2>&1 | tail -15; echo "=== restore ==="; cp /tmp/native.yaml.bak docs/api/astrate_native_api.yaml && go test ./internal/swagger/ -run TestNativeSocketSecurityDocumented 2>&1 | tail -5 failed
Error: The user rejected permission to use this specific tool call.
