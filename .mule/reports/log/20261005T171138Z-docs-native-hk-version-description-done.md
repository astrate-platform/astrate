slug: docs-native-hk-version-description
verdict: done
at:  cce1376
ran: 2026-10-05T17:11:38Z on DietPi in 360s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ rg -n "emulated upstream API level is served by the realm-scoped" -C 12 docs/api/astrate_native_api.yaml
(no output)

→ Read docs/api/astrate_native_api.yaml [offset=200, limit=120]

$ rg -n "emulated upstream API" -C 3 docs/api/astrate_native_api.yaml
189-      description: |
190-        Upstream-parity unauthenticated version endpoint. Returns the Astrate
191-        build version serving the AppEngine service; the Astarte Dashboard
192:        reads one per service at startup. The emulated upstream API level is
193-        served by the realm-scoped /v1/{realm}/version endpoints.
194-      tags: [Compat Health]
195-      responses:
--
209-      description: |
210-        Upstream-parity unauthenticated version endpoint. Returns the Astrate
211-        build version serving the Realm Management service; the Astarte
212:        Dashboard reads one per service at startup. The emulated upstream API
213-        level is served by the realm-scoped /v1/{realm}/version endpoints.
214-      tags: [Compat Health]
215-      responses:
--
229-      description: |
230-        Upstream-parity unauthenticated version endpoint. Returns the Astrate
231-        build version serving the Pairing service; the Astarte Dashboard reads
232:        one per service at startup. The emulated upstream API level is served
233-        by the realm-scoped /v1/{realm}/version endpoints.
234-      tags: [Compat Health]
235-      responses:
--
249-      description: |
250-        Upstream-parity unauthenticated version endpoint. Returns the Astrate
251-        build version serving the Housekeeping service; the Astarte Dashboard
252:        reads one per service at startup. The emulated upstream API level is
253-        served by the realm-scoped /v1/{realm}/version endpoints.
254-      tags: [Compat Health]
255-      responses:

$ rg -n "mountRealmVersion" -C 8 cmd/astrate/main.go
391-	// service answers an unauthenticated GET /{service}/version with
392-	// {"data": version}; AppEngine's realm-scoped variant requires auth,
393-	// Pairing's is public — both measured on upstream 1.2.0
394-	// (test/conformance/upstream/verify-versions.json). Realm Management's
395-	// realm-scoped route is served by its own API.
396-	for _, svc := range []string{"appengine", "realmmanagement", "pairing", "housekeeping"} {
397-		observability.MountVersionCompat(mux, svc, version)
398-	}
399:	mountRealmVersion(mux, mw)
400-
401-	metrics.RegisterBrokerSessions(func() float64 { return float64(b.SessionCount()) })
402-	metrics.RegisterDBPool(func() observability.DBPoolStats {
403-		s := st.Stat()
404-		return observability.DBPoolStats{
405-			AcquiredConns: s.AcquiredConns(),
406-			IdleConns:     s.IdleConns(),
407-			TotalConns:    s.TotalConns(),
--
421-	// client parsing the envelope would choke on.
422-	handler := httpx.NotFound(mux)
423-	if len(cfg.HTTP.CORSAllowedOrigins) > 0 {
424-		handler = httpx.CORS(cfg.HTTP.CORSAllowedOrigins)(handler)
425-	}
426-	return handler, hkSvc, nil
427-}
428-
429:// mountRealmVersion registers the realm-scoped version endpoints observed on
430-// upstream 1.2.0 (test/conformance/upstream/verify-versions.json): AppEngine's
431-// requires a realm token, pairing's is public. Both report the emulated
432-// upstream API level, the value the Astarte Dashboard feature-gates its UI on
433-// (COMPATIBILITY.md deviation 10) — not Astrate's build version. RM's
434-// realm-scoped route is served by its own API.
435:func mountRealmVersion(mux *http.ServeMux, mw *auth.Middleware) {
436-	mux.Handle("GET /appengine/v1/{realm}/version",
437-		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(realm.APICompatVersion)))
438-	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
439-}
440-
441-// brokerReadiness reports the broker listener as ready when a TCP connection to
442-// it is accepted (docs/DESIGN.md §5.2 readiness broker check).
443-func brokerReadiness(b *broker.Broker) observability.Check {

$ rg -n "version" internal/realm/http.go internal/housekeeping/http.go
internal/realm/http.go:52:	mux.Handle("GET /realmmanagement/v1/{realm}/version", h(a.getVersion))
internal/realm/http.go:372:			"Interface major version doesn't match the one in the interface json")
internal/realm/http.go:375:			"Interface minor version was not increased")
internal/realm/http.go:416:			"Interface major version is not 0, can't be deleted")

→ Read internal/housekeeping/http.go [limit=60, offset=20]

$ rg -n "MountVersionCompat" -A 25 internal/observability/*.go | head -60
internal/observability/compat.go:49:// MountVersionCompat registers GET /{service}/version for one service.
internal/observability/compat.go:50:func MountVersionCompat(mux *http.ServeMux, service, version string) {
internal/observability/compat.go-51-	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
internal/observability/compat.go-52-}
--
internal/observability/compat_version_test.go:10:func TestMountVersionCompat(t *testing.T) {
internal/observability/compat_version_test.go-11-	mux := http.NewServeMux()
internal/observability/compat_version_test.go-12-	for _, svc := range []string{"housekeeping", "appengine", "realmmanagement", "pairing"} {
internal/observability/compat_version_test.go:13:		MountVersionCompat(mux, svc, "0.0.0-test")
internal/observability/compat_version_test.go-14-	}
internal/observability/compat_version_test.go-15-	srv := httptest.NewServer(mux)
internal/observability/compat_version_test.go-16-	defer srv.Close()
internal/observability/compat_version_test.go-17-
internal/observability/compat_version_test.go-18-	for _, svc := range []string{"housekeeping", "appengine", "realmmanagement", "pairing"} {
internal/observability/compat_version_test.go-19-		resp, err := http.Get(srv.URL + "/" + svc + "/version")
internal/observability/compat_version_test.go-20-		if err != nil {
internal/observability/compat_version_test.go-21-			t.Fatalf("%s/version: %v", svc, err)
internal/observability/compat_version_test.go-22-		}
internal/observability/compat_version_test.go-23-		body, _ := io.ReadAll(resp.Body)
internal/observability/compat_version_test.go-24-		resp.Body.Close()
internal/observability/compat_version_test.go-25-		if resp.StatusCode != http.StatusOK {
internal/observability/compat_version_test.go-26-			t.Errorf("%s/version = %d, want 200", svc, resp.StatusCode)
internal/observability/compat_version_test.go-27-		}
internal/observability/compat_version_test.go-28-		if string(body) != `{"data":"0.0.0-test"}` {
internal/observability/compat_version_test.go-29-			t.Errorf("%s/version body = %s", svc, body)
internal/observability/compat_version_test.go-30-		}
internal/observability/compat_version_test.go-31-		if ct := resp.Header.Get("Content-Type"); ct == "" {
internal/observability/compat_version_test.go-32-			t.Errorf("%s/version missing Content-Type", svc)
internal/observability/compat_version_test.go-33-		}
internal/observability/compat_version_test.go-34-	}
internal/observability/compat_version_test.go-35-}
internal/observability/compat_version_test.go-36-
internal/observability/compat_version_test.go-37-func TestVersionHandlerRealmScoped(t *testing.T) {
internal/observability/compat_version_test.go-38-	mux := http.NewServeMux()

$ rg -n "astrate_native_api" --glob '!docs/api/*' -l; echo ---; rg -n "docs" Makefile docs/Makefile 2>/dev/null | head -40
docs/swagger-ui/index.html
docs/handoff/phase-2-handoff.md
docs/handoff/phase-2-memory.md
docs/handoff/phase-1-memory.md
internal/swagger/swagger_test.go
---
docs/Makefile:9:	cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
docs/Makefile:14:	mkdocs serve -f mkdocs.yml
docs/Makefile:17:	mkdocs build -f mkdocs.yml
Makefile:1:# Astrate build & verification entry points (docs/ROADMAP.md §1.2, verification tiers §0.2).
Makefile:53:		echo "conformance harness lands in M9 (docs/ROADMAP.md §10); nothing to run yet"; \

→ Read docs/Makefile
→ Read internal/swagger/swagger_test.go

$ rg -n "getHousekeepingVersion|realm-scoped|api/astrate_native" internal/swagger/ docs/swagger-ui/index.html | head -30; echo ---; rg -n "yaml|url:" docs/swagger-ui/index.html | head -20
docs/swagger-ui/index.html:48:      <option value="../api/astrate_native_api.yaml">Astrate Native API — Health, readiness, metrics, WebSocket</option>
internal/swagger/swagger_test.go:630:	b, err := docs.APIYAML.ReadFile("api/astrate_native_api.yaml")
---
44:      <option value="../api/astarte_housekeeping_api.yaml">Housekeeping API — Realm management, a_ha auth</option>
45:      <option value="../api/astarte_pairing_api.yaml">Pairing API — Device pairing flows A/B/C, a_pa auth</option>
46:      <option value="../api/astarte_realm_management_api.yaml">Realm Management API — Interfaces, triggers, policies, a_rma auth</option>
47:      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
48:      <option value="../api/astrate_native_api.yaml">Astrate Native API — Health, readiness, metrics, WebSocket</option>
65:        url: specUrl,

$ which mkdocs; python3 -c "import mkdocs; print(mkdocs.__version__)" 2>&1 | tail -1; echo ---; cat .gitignore 2>/dev/null | head -30; echo ---; git -C /root/astrate-mule status --porcelain | head
ModuleNotFoundError: No module named 'mkdocs'
---
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

---

$ rg -n "site" .gitignore; echo "--- requirements:"; cat docs/requirements.txt 2>/dev/null; echo "--- go yaml deps:"; rg -n "yaml" go.mod
29:/site-dist/
31:# Docs site copies (refreshed by docs/make sync)
32:/docs/site/DESIGN.md
33:/docs/site/ROADMAP.md
34:/docs/site/api/
35:/docs/site/swagger-ui/
--- requirements:
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10
--- go yaml deps:
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

$ rg -n "yaml\.(Unmarshal|NewDecoder)" --glob '*.go' -l | head -20; echo ---; python3 -c "import yaml; print('pyyaml ok', yaml.__version__)"
---
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; print('pyyaml ok', yaml.__version__)
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ rg -n "deviation 10|Deviation 10|^### 10|^## 10" -A 20 docs/COMPATIBILITY.md | head -60
(no output)

$ cat test/conformance/upstream/verify-versions.json 2>/dev/null | head -60
{
  "recorded_at": "2026-08-24T14:30:00Z",
  "astarte_version": "1.2.0",
  "realm": "zzprobe08",
  "recorder_command": "scratch probes (fase 1 verify batch, issues #69/#77); see verify-versions.transcript.txt",
  "observations": [
    {
      "name": "service-root version, housekeeping", "why": "#77: the dashboard polls one version endpoint per service at startup",
      "method": "GET", "path": "/housekeeping/version", "auth": "none",
      "status": 200, "body": "{\"data\":\"1.2.0\"}"
    },
    {
      "name": "service-root version, appengine", "why": "#77",
      "method": "GET", "path": "/appengine/version", "auth": "none",
      "status": 200, "body": "{\"data\":\"1.2.0\"}"
    },
    {
      "name": "service-root version, realmmanagement", "why": "#77",
      "method": "GET", "path": "/realmmanagement/version", "auth": "none",
      "status": 200, "body": "{\"data\":\"1.2.0\"}"
    },
    {
      "name": "service-root version, pairing", "why": "#77",
      "method": "GET", "path": "/pairing/version", "auth": "none",
      "status": 200, "body": "{\"data\":\"1.2.0\"}"
    },
    {
      "name": "global version route absent", "why": "#77: bounds the surface — there is no instance-wide /version",
      "method": "GET", "path": "/version", "auth": "none",
      "status": 404, "body": "404 page not found"
    },
    {
      "name": "flow service absent from compose", "why": "#77: the dashboard's fourth target is not routed on upstream 1.2 self-hosted",
      "method": "GET", "path": "/flow/version", "auth": "none",
      "status": 404, "body": "{\"errors\":{\"detail\":\"Page not found\"}}"
    },
    {
      "name": "appengine realm-scoped version, unauthenticated", "why": "#77: realm-scoped AE/RM versions REQUIRE auth (401), unlike pairing's",
      "method": "GET", "path": "/appengine/v1/zzprobe08/version", "auth": "none",
      "status": 401, "body": "{\"errors\":{\"detail\":\"Unauthorized\"}}"
    },
    {
      "name": "appengine realm-scoped version, valid token", "why": "#77 acceptance row: the paired acceptance proving 401 is auth, not absence",
      "method": "GET", "path": "/appengine/v1/zzprobe08/version", "auth": "Bearer <valid a_aea .*::.* token>",
      "status": 200, "body": "{\"data\":\"1.2.0\"}"
    },
    {
      "name": "pairing realm-scoped version, unauthenticated", "why": "#77: upstream serves pairing's realm version PUBLIC — the odd one out of the three services",
      "method": "GET", "path": "/pairing/v1/zzprobe08/version", "auth": "none",
      "status": 200, "body": "{\"data\":\"1.2.0\"}"
    },
    {
      "name": "unknown realm on RM, no token", "why": "#69 baseline row",
      "method": "GET", "path": "/realmmanagement/v1/nosuchrealm/interfaces", "auth": "none",
      "status": 401, "body": "{\"errors\":{\"detail\":\"Unauthorized\"}}"
    },
    {
      "name": "unknown realm on RM, bearer token present", "why": "#69: ANY non-absent token → 403 Forbidden (unknown realm has no key; every verification failure maps to 403). Astrate answers uniform 401 (deliberate deviation-5 extension)",
      "method": "GET", "path": "/realmmanagement/v1/nosuchrealm/interfaces", "auth": "Bearer <token signed by any key>",
      "status": 403, "body": "{\"errors\":{\"detail\":\"Forbidden\"}}"

$ rg -n "housekeeping" -A 5 test/conformance/upstream/verify-versions.json | rg -n "realm-scoped|housekeeping/v1"
(no output)

$ rg -n "APICompatVersion|realm-scoped|no realm-scoped" docs/COMPATIBILITY.md | head -20
253:    reports `realm.APICompatVersion` (currently `1.2.2`), a compatibility
359:    additionally keeps the realm-scoped `GET /pairing/v1/{realm}/health` (#71),

→ Read docs/COMPATIBILITY.md [offset=240, limit=40]

$ rg -n "operationId: get.*Version" docs/api/astrate_native_api.yaml; echo ---; rg -rn "emulated upstream API level" --glob '*.go' .; echo ---; rg -n "func NotFound" -A 12 internal/httpx/*.go
187:      operationId: getAppEngineVersion
207:      operationId: getRealmManagementVersion
227:      operationId: getPairingVersion
247:      operationId: getHousekeepingVersion
---
./cmd/astrate/version_test.go:// endpoints serve the n the Dashboard feature-gates
./internal/realm/http.go:// getVersion reports the n (the dashboard gates
---
internal/httpx/notfound.go:14:func NotFound(mux *http.ServeMux) http.Handler {
internal/httpx/notfound.go-15-	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
internal/httpx/notfound.go-16-		_, pattern := mux.Handler(r)
internal/httpx/notfound.go-17-		if pattern != "" {
internal/httpx/notfound.go-18-			mux.ServeHTTP(w, r)
internal/httpx/notfound.go-19-			return
internal/httpx/notfound.go-20-		}
internal/httpx/notfound.go-21-
internal/httpx/notfound.go-22-		rec := httptest.NewRecorder()
internal/httpx/notfound.go-23-		mux.ServeHTTP(rec, r)
internal/httpx/notfound.go-24-
internal/httpx/notfound.go-25-		// Replay what the mux produced, byte for byte.
internal/httpx/notfound.go-26-		replay := func() {

$ rg -n "realm-scoped|APICompatVersion" -C 6 cmd/astrate/version_test.go | head -60
6-	"testing"
7-
8-	"github.com/astrate-platform/astrate/internal/auth"
9-	"github.com/astrate-platform/astrate/internal/realm"
10-)
11-
12:// TestRealmScopedVersionReportsAPICompatLevel proves the realm-scoped version
13-// endpoints serve the emulated upstream API level the Dashboard feature-gates
14-// its UI on (COMPATIBILITY.md deviation 10, issue #77), not Astrate's build
15-// version. The pairing route is public, so it is reachable without a token.
16-func TestRealmScopedVersionReportsAPICompatLevel(t *testing.T) {
17-	mux := http.NewServeMux()
18-	mountRealmVersion(mux, auth.NewMiddleware(nil))
19-
20-	rec := httptest.NewRecorder()
21-	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/pairing/v1/somerealm/version", nil))
22-	if rec.Code != http.StatusOK {
23-		t.Fatalf("pairing version status = %d, want 200", rec.Code)
24-	}
25:	want := `{"data":"` + realm.APICompatVersion + `"}`
26-	if got := rec.Body.String(); got != want {
27-		t.Errorf("pairing version body = %s, want %s", got, want)
28-	}
29-}

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -248,10 +248,12 @@
 summary: Housekeeping API version
 description: |
   Upstream-parity unauthenticated version endpoint. Returns the Astrate
   build version serving the Housekeeping service; the Astarte Dashboard
-  reads one per service at startup. The emulated upstream API level is
-  served by the realm-scoped /v1/{realm}/version endpoints.
+  reads one per service at startup. Housekeeping is the one service with
+  no realm-scoped twin: `/housekeeping/v1/{realm}/version` is not a route
+  and answers 404, so unlike AppEngine, Realm Management and Pairing no
+  emulated upstream API level is reported here.
 tags: [Compat Health]
 responses:
   "200":
     description: Astrate build version.

→ Read docs/api/astrate_native_api.yaml [limit=26, offset=243]
Now the guard test, following the file's existing docs-guard style:

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -662,8 +662,52 @@
 		}
 	}
 }
 
+// TestNativeVersionDescriptions guards the four compat-version operations of the
+// native spec, which used to carry one identical sentence: "The emulated
+// upstream API level is served by the realm-scoped /v1/{realm}/version
+// endpoints". That is true for three of them and false for the fourth.
+// mountRealmVersion registers only appengine and pairing (cmd/astrate/main.go),
+// realm management's realm-scoped route is served by its own API
+// (internal/realm/http.go), and housekeeping's Mount registers just the five
+// realm-lifecycle routes (internal/housekeeping/http.go), so
+// /housekeeping/v1/{realm}/version matches no pattern and falls through to
+// httpx.NotFound (cmd/astrate/main.go). A generated client that believed the
+// sentence would build a call that can only ever get a 404, and the difference
+// is invisible in the response schema — both are {"data": string}. So
+// housekeeping now says plainly that it is the one service with no realm-scoped
+// twin, and the other three keep the pointer.
+func TestNativeVersionDescriptions(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astrate_native_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astrate_native_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	const twin = "realm-scoped /v1/{realm}/version"
+
+	for _, op := range []string{"getAppEngineVersion", "getRealmManagementVersion", "getPairingVersion"} {
+		if desc := operationDescription(t, operationBlock(t, lines, op)); !strings.Contains(desc, twin) {
+			t.Errorf("%s description does not point at the %s endpoints", op, twin)
+		}
+	}
+
+	desc := operationDescription(t, operationBlock(t, lines, "getHousekeepingVersion"))
+	for _, want := range []string{
+		"no realm-scoped twin",
+		"/housekeeping/v1/{realm}/version",
+		"404",
+	} {
+		if !strings.Contains(desc, want) {
+			t.Errorf("getHousekeepingVersion description does not say %q", want)
+		}
+	}
+	if strings.Contains(desc, twin) {
+		t.Errorf("getHousekeepingVersion description still promises the %s endpoints, which housekeeping does not serve", twin)
+	}
+}
+
 // TestHousekeepingErrorDetailExamples pins the error-detail examples in the
 // housekeeping spec's components.responses to the frozen canonical strings the
 // wire emits (astarteapi/envelope.go), the same guard the pairing, realm
 // management, app engine and native specs get from

$ gofmt -l . ; go vet ./internal/swagger/ && go test ./internal/swagger/ -run 'TestNative' -v 2>&1 | tail -20
=== RUN   TestNativeErrorDetailExamples
--- PASS: TestNativeErrorDetailExamples (0.00s)
=== RUN   TestNativeVersionDescriptions
--- PASS: TestNativeVersionDescriptions (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.014s

$ cp docs/api/astrate_native_api.yaml /tmp/opencode/native_api.yaml.bak && ls -l /tmp/opencode/native_api.yaml.bak
-rw-r--r-- 1 root root 15125 Oct  5 19:08 /tmp/opencode/native_api.yaml.bak

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -248,12 +248,10 @@
 summary: Housekeeping API version
 description: |
   Upstream-parity unauthenticated version endpoint. Returns the Astrate
   build version serving the Housekeeping service; the Astarte Dashboard
-  reads one per service at startup. Housekeeping is the one service with
-  no realm-scoped twin: `/housekeeping/v1/{realm}/version` is not a route
-  and answers 404, so unlike AppEngine, Realm Management and Pairing no
-  emulated upstream API level is reported here.
+  reads one per service at startup. The emulated upstream API level is
+  served by the realm-scoped /v1/{realm}/version endpoints.
 tags: [Compat Health]
 responses:
   "200":
     description: Astrate build version.

$ go test ./internal/swagger/ -run 'TestNativeVersionDescriptions' 2>&1 | tail -12
--- FAIL: TestNativeVersionDescriptions (0.00s)
    swagger_test.go:702: getHousekeepingVersion description does not say "no realm-scoped twin"
    swagger_test.go:702: getHousekeepingVersion description does not say "/housekeeping/v1/{realm}/version"
    swagger_test.go:702: getHousekeepingVersion description does not say "404"
    swagger_test.go:706: getHousekeepingVersion description still promises the realm-scoped /v1/{realm}/version endpoints, which housekeeping does not serve
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.017s
FAIL

$ cp /tmp/opencode/native_api.yaml.bak docs/api/astrate_native_api.yaml && go test ./internal/swagger/ 2>&1 | tail -3
ok  	github.com/astrate-platform/astrate/internal/swagger	0.052s

Now the docs build and the YAML-load check:
$ make -C docs build 2>&1 | tail -15
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

$ timeout 20 pip install --quiet "mkdocs>=1.6,<2" "mkdocs-material>=9.5,<10" 2>&1 | tail -5; echo "exit=$?"
timeout: failed to run command ‘pip’: No such file or directory
exit=0

$ command -v pip3 python3 node npx ruby perl 2>&1; echo ---; timeout 25 pip3 install --quiet "mkdocs>=1.6,<2" "mkdocs-material>=9.5,<10" 2>&1 | tail -4
---
timeout: failed to run command ‘pip3’: No such file or directory

$ go env GOMODCACHE GOFLAGS GOPROXY; ls "$(go env GOMODCACHE)/gopkg.in" 2>&1 | head; ls -d "$(go env GOMODCACHE)"/gopkg.in/yaml.v3@* 2>&1 | head
/root/go/pkg/mod

https://proxy.golang.org,direct
check.v1@v0.0.0-20161208181325-20d25e280405
check.v1@v1.0.0-20201130134442-10cb98267c6c
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1
/root/go/pkg/mod/gopkg.in/yaml.v3@v3.0.1

← Write /tmp/opencode/yamlprobe/main.go
Wrote file successfully.

$ cd /tmp/opencode/yamlprobe && cat > go.mod <<'EOF'
module yamlprobe

go 1.24

require gopkg.in/yaml.v3 v3.0.1

replace gopkg.in/yaml.v3 => /root/go/pkg/mod/gopkg.in/yaml.v3@v3.0.1
EOF
GOFLAGS=-mod=mod GOPROXY=off go build -o probe . 2>&1 | tail -5 && echo BUILD_OK
BUILD_OK

$ /tmp/opencode/yamlprobe/probe /root/astrate-mule/docs/api/*.yaml; echo "probe exit=$?"
ok   /root/astrate-mule/docs/api/astarte_appengine_api.yaml: openapi=3.0.3 paths=18 components=4
ok   /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml: openapi=3.0.3 paths=2 components=4
ok   /root/astrate-mule/docs/api/astarte_pairing_api.yaml: openapi=3.0.3 paths=7 components=4
ok   /root/astrate-mule/docs/api/astarte_realm_management_api.yaml: openapi=3.0.3 paths=12 components=4
ok   /root/astrate-mule/docs/api/astrate_native_api.yaml: openapi=3.0.3 paths=12 components=3

--- getHousekeepingVersion description as loaded ---
Upstream-parity unauthenticated version endpoint. Returns the Astrate
build version serving the Housekeeping service; the Astarte Dashboard
reads one per service at startup. Housekeeping is the one service with
no realm-scoped twin: `/housekeeping/v1/{realm}/version` is not a route
and answers 404, so unlike AppEngine, Realm Management and Pairing no
emulated upstream API level is reported here.

---

probe exit=0

$ gofmt -l . ; echo "gofmt done"; go vet ./... 2>&1 | tail -5; echo "vet exit=$?"
gofmt done
vet exit=0

$ go test ./... 2>&1 | grep -v "^ok\|no test files" | head -20; echo "plain test exit=${PIPESTATUS[0]}"
/bin/sh: 1: Bad substitution

$ go test ./... > /tmp/opencode/gotest.log 2>&1; echo "exit=$?"; grep -v "^ok\|no test files" /tmp/opencode/gotest.log | head -20
exit=0

$ go test -race ./internal/swagger/ ./cmd/astrate/ > /tmp/opencode/gorace.log 2>&1; echo "exit=$?"; tail -5 /tmp/opencode/gorace.log
exit=1
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.025s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/cmd/astrate	0.022s
FAIL

$ git status --porcelain; echo "--- diffstat:"; git diff --stat
M docs/api/astrate_native_api.yaml
 M internal/swagger/swagger_test.go
--- diffstat:
 docs/api/astrate_native_api.yaml |  6 ++++--
 internal/swagger/swagger_test.go | 44 ++++++++++++++++++++++++++++++++++++++++
 2 files changed, 48 insertions(+), 2 deletions(-)

$ rg -n "docs-native-ae-realm-version-endpoint" .mule/todo.md | head -3
277:- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
278:- [ ] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ git diff -- docs/api/astrate_native_api.yaml
diff --git a/docs/api/astrate_native_api.yaml b/docs/api/astrate_native_api.yaml
index cbd02c8..4352661 100644
--- a/docs/api/astrate_native_api.yaml
+++ b/docs/api/astrate_native_api.yaml
@@ -249,8 +249,10 @@ paths:
       description: |
         Upstream-parity unauthenticated version endpoint. Returns the Astrate
         build version serving the Housekeeping service; the Astarte Dashboard
-        reads one per service at startup. The emulated upstream API level is
-        served by the realm-scoped /v1/{realm}/version endpoints.
+        reads one per service at startup. Housekeeping is the one service with
+        no realm-scoped twin: `/housekeeping/v1/{realm}/version` is not a route
+        and answers 404, so unlike AppEngine, Realm Management and Pairing no
+        emulated upstream API level is reported here.
       tags: [Compat Health]
       responses:
         "200":

```
Done: /housekeeping/version in docs/api/astrate_native_api.yaml no longer promises a realm-scoped /v1/{realm}/version twin that housekeeping does not serve; it now says plainly that housekeeping is the one service with none, and the three sibling compat-version descriptions are pinned by a new guard test.
Files: docs/api/astrate_native_api.yaml (yaml:249-255), internal/swagger/swagger_test.go (new TestNativeVersionDescriptions, before TestHousekeepingErrorDetailExamples)
Verified: claims checked first against cmd/astrate/main.go:396-439 (mountRealmVersion registers only appengine + pairing), internal/realm/http.go:52, internal/housekeeping/http.go:37-41 (five realm-lifecycle routes, no version); `go test ./internal/swagger/ -run TestNativeVersionDescriptions` -> FAIL on the old YAML (4 assertions), PASS restored -> the test fails without the change; `go vet ./...` -> pass; `go test ./...` (no -race: this box is 39-bit VMA, `go test -race` dies with ThreadSanitizer: unsupported VMA range) -> pass; `gofmt -l .` -> empty; `make -C docs build` -> `sync` ran, then FAILED at `mkdocs build` (Error 127, mkdocs not installed; no pip/pip3/python3/node on this machine to install it), so instead I parsed all five specs the Swagger UI references with a throwaway gopkg.in/yaml.v3 probe (/tmp/opencode/yamlprobe) -> all load as OpenAPI 3.0.3 and the new description round-trips verbatim
Unsure: whether the mkdocs half of `make -C docs build` would pass — unverifiable here. The sync step (the only part that touches the API YAML) copies it unmodified, and mkdocs does not parse it, so I expect no impact, but I did not run it.
Follow-ups: none
```
