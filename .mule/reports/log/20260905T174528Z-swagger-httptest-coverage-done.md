slug: swagger-httptest-coverage
verdict: done
at:  ba79b34
ran: 2026-09-05T17:45:28Z on DietPi in 267s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
• Explore internal/swagger package Explore Agent
✓ Explore internal/swagger package Explore Agent
→ Read docs/embed.go
✱ Glob "docs/swagger-ui/*" 1 match
✱ Glob "docs/api/*.yaml" 5 matches
→ Read docs/swagger-ui/index.html
Now I understand the package fully. Let me write the test suite.

Task meaning: Add an httptest-based test file (`swagger_test.go`) for `internal/swagger`, covering Mount (302 redirect, /swagger/ UI, /api/ YAML specs) and Specs (exact yaml filenames).

Let me check the module name to write correct imports.

$ head -1 /root/astrate-mule/go.mod
module github.com/astrate-platform/astrate

Now I'll write the test file. I need tests that fail without the implementation change — but this is adding a test suite for existing behavior, so the tests should validate existing (already-correct) behavior. Let me write table-driven httptest tests.
← Write internal/swagger/swagger_test.go
Wrote file successfully.

I referenced `WalkYAML` which doesn't exist. Let me fix the helper to use `fs.WalkDir` directly to keep implementation-independent enumeration.
← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -1,8 +1,9 @@
 package swagger
 
 import (
 	"io"
+	"io/fs"
 	"net/http"
 	"net/http/httptest"
 	"reflect"
 	"sort"

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -89,10 +89,19 @@
 }
 
 func embeddedYAMLFilenames() ([]string, error) {
 	var names []string
-	err := WalkYAML(docs.APIYAML, func(path string) {
-		names = append(names, strings.TrimPrefix(path, "api/"))
+	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
+		if err != nil {
+			return err
+		}
+		if d.IsDir() {
+			return nil
+		}
+		if strings.HasSuffix(path, ".yaml") {
+			names = append(names, strings.TrimPrefix(path, "api/"))
+		}
+		return nil
 	})
 	sort.Strings(names)
 	return names, nil
 }

Now let me run the tests. Since I'm on a machine without the race detector issue being a blocker for this container-free test, `go test` alone is the container-free gate. Let me verify the task said "go test -race ./..." per step 5 but AGENTS.md says race detector doesn't work on the Pi. The task is container-free httptest — `go test -race` may work here depending on kernel. Let me run both but tolerate race failure.
$ gofmt -l internal/swagger/
(no output)

$ go test ./internal/swagger/ 2>&1
# github.com/astrate-platform/astrate/internal/swagger [github.com/astrate-platform/astrate/internal/swagger.test]
internal/swagger/swagger_test.go:93:2: declared and not used: err
FAIL	github.com/astrate-platform/astrate/internal/swagger [build failed]
FAIL

The `embeddedYAMLFilenames` returns error but I use `err := fs.WalkDir` handled inside. Line 93 — let me check. The function declares `err := fs.WalkDir(...)` and returns it. Wait, `err` is used in callback and return. Let me look.

Actually line 93 is in the callback `return err` — no. The declared-and-not-used is the outer `err := fs.WalkDir`. Let me re-read. The callback uses its own `err` param shadowing. The outer `err := fs.WalkDir(...)` is captured but... it IS returned: `return names, err`. Hmm. Let me check what's on line 93.

Actually wait — inside the callback I did `return err` but the callback signature has `err fs.DirEntry ... error`. Let me re-read my code. The callback: `func(path string, d fs.DirEntry, err error) error` then `if err != nil { return err }`. That's fine. 

Let me just look at the file around line 93.
→ Read internal/swagger/swagger_test.go [offset=88, limit=15]
The callback shadows `err` but doesn't use it (returns err but the shadowed `err` param — that IS the shadow). The outer `err` at line 93 is used in `return names, err`. Why "declared and not used"? The callback shadows `err` so within the callback the outer `err` is hidden, but the outer `err` is still assigned. Hmm, actually the issue: `err := fs.WalkDir(...)` — the `err` variable is used in `return names, err`. That should count as used.

Wait, unless I never return it. Let me read the return.
→ Read internal/swagger/swagger_test.go [offset=102, limit=25]
Line 106: `return names, nil` — the outer `err` is never used in the return, hence "declared and not used". Fix to return err.

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -102,9 +102,9 @@
 		}
 		return nil
 	})
 	sort.Strings(names)
-	return names, nil
+	return names, err
 }
 
 func get(t *testing.T, url string) string {
 	t.Helper()

$ go test ./internal/swagger/ 2>&1 && echo "---VET---" && go vet ./internal/swagger/ 2>&1 && echo "---GOFMT---" && gofmt -l internal/swagger/
--- FAIL: TestMount (0.03s)
    --- FAIL: TestMount/GET_/swagger_redirects_to_/swagger/index.html (0.02s)
        swagger_test.go:30: status = 200, want 302
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.035s
FAIL

The redirect returns 200 instead of 302. This is likely because `http.Get` follows redirects — the request to `/swagger` gets redirected to `/swagger/index.html` and `http.Get` follows it, resulting in a 200. I need to disable redirect following to test the 302. Let me use a `CheckRedirect` that returns `http.ErrUseLastResponse`.
← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -19,9 +19,12 @@
 srv := httptest.NewServer(mux)
 defer srv.Close()
 
 t.Run("GET /swagger redirects to /swagger/index.html", func(t *testing.T) {
-	resp, err := http.Get(srv.URL + "/swagger")
+	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
+		return http.ErrUseLastResponse
+	}}
+	resp, err := client.Get(srv.URL + "/swagger")
 	if err != nil {
 		t.Fatalf("GET /swagger: %v", err)
 	}
 	defer resp.Body.Close()

$ go test ./internal/swagger/ -v 2>&1 && echo "---VET---" && go vet ./internal/swagger/ 2>&1 && echo "---GOFMT---" && gofmt -l ./internal/swagger/
=== RUN   TestMount
=== RUN   TestMount/GET_/swagger_redirects_to_/swagger/index.html
=== RUN   TestMount/GET_/swagger/_serves_the_embedded_UI
=== RUN   TestMount/GET_/api/_serves_every_OpenAPI_YAML_spec
--- PASS: TestMount (0.06s)
    --- PASS: TestMount/GET_/swagger_redirects_to_/swagger/index.html (0.01s)
    --- PASS: TestMount/GET_/swagger/_serves_the_embedded_UI (0.03s)
    --- PASS: TestMount/GET_/api/_serves_every_OpenAPI_YAML_spec (0.01s)
=== RUN   TestSpecs
--- PASS: TestSpecs (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.077s
---VET---
---GOFMT---

All tests pass. Let me run the broader package test and race to complete the gate. Per AGENTS.md, `-race` may not work on this kernel — the task is container-free so `go test` is fine. Let me try `go test -race` but tolerate a race-build failure since the Pi can't run it.
$ go test ./internal/... 2>&1 | tail -30
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
?   	github.com/astrate-platform/astrate/internal/housekeeping	[no test files]
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	0.079s
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)

All `internal/...` tests pass. Now let me verify the test suite would actually catch a broken implementation (the runner validates this). Let me confirm via the runner's mental model — my `TestMount` asserts exact redirect path, UI byte-match on embedded copy, and each spec byte-match; `TestSpecs` asserts no-path-prefix and exact enumeration. These genuinely break if the implementation regresses.

Let me try the `-race` gate since this might be a machine where it works, and report honestly if not.
$ go test -race ./internal/swagger/ 2>&1; echo "EXIT: $?"
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.048s
FAIL
EXIT: 1

Race detector fails to build on this kernel (39-bit VMA) as AGENTS.md documents; the suite is httptest-only and single-goroutine, so `go test` is the honest gate here. Let me do a final review of the test file.
→ Read internal/swagger/swagger_test.go
Done: added a container-free httptest suite for internal/swagger covering Mount (302→/swagger/index.html, embedded UI under /swagger/, YAML specs under /api/) and Specs (exact .yaml filenames, no path prefix or dirs)
Files: internal/swagger/swagger_test.go
Verified: go test ./internal/swagger/ -> pass; go vet ./internal/swagger/ -> pass; gofmt -l ./internal/swagger/ -> clean; go test ./internal/... -> pass; go test -race ./internal/swagger/ -> fail (39-bit VMA kernel, race detector cannot launch — per AGENTS.md, said loudly)
Unsure: nothing
Follow-ups: none
