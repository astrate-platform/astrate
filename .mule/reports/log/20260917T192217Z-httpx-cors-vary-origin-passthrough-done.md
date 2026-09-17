slug: httpx-cors-vary-origin-passthrough
verdict: done
at:  986c100
ran: 2026-09-17T19:22:17Z on DietPi in 143s


> build · big-pickle

$ ls -la .mule/ 2>/dev/null || find . -name "MULE.md" -not -path "./.git/*" 2>/dev/null
total 4904
drwxr-xr-x  7 root root    4096 Sep 17 21:18 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      13 Sep 17 21:18 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 17 13:45 .heartbeat
-rw-r--r--  1 root root       0 Sep 17 21:18 .last-output
-rw-r--r--  1 root root       0 Sep 17 21:18 .lock
-rw-r--r--  1 root root      12 Sep 17 21:08 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 4816728 Sep 17 21:18 cron.log
drwxr-xr-x  2 root root    4096 Sep 15 19:13 failed
-rw-r--r--  1 root root   37387 Sep 17 21:03 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    7213 Sep 17 13:45 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 17 21:11 reviews
-rw-r--r--  1 root root     806 Sep 17 21:18 task.md
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   61943 Sep 17 21:12 todo.md

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/httpx/cors.go
→ Read internal/httpx/cors_test.go
Reading the task: add `Vary: Origin` to the disallowed-origin pass-through branch in non-wildcard mode, and pin it in tests. The branch covers both no-Origin and disallowed-origin requests; the task's stated harm includes "the headerless variant to an allowed one", so the whole pass-through branch should carry the variation marker.

← Edit internal/httpx/cors.go
Index: /root/astrate-mule/internal/httpx/cors.go
===================================================================
--- /root/astrate-mule/internal/httpx/cors.go
+++ /root/astrate-mule/internal/httpx/cors.go
@@ -23,8 +23,11 @@
 return func(next http.Handler) http.Handler {
 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 		origin := r.Header.Get("Origin")
 		if origin == "" || (!wildcard && !allowed[origin]) {
+			if !wildcard {
+				w.Header().Add("Vary", "Origin")
+			}
 			next.ServeHTTP(w, r)
 			return
 		}
 

← Edit internal/httpx/cors_test.go
Index: /root/astrate-mule/internal/httpx/cors_test.go
===================================================================
--- /root/astrate-mule/internal/httpx/cors_test.go
+++ /root/astrate-mule/internal/httpx/cors_test.go
@@ -56,8 +56,11 @@
 	}
 	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:4040" {
 		t.Errorf("Allow-Origin = %q", got)
 	}
+	if got := rec.Header().Get("Vary"); got != "Origin" {
+		t.Errorf("Vary = %q, want Origin", got)
+	}
 }
 
 func TestCORSDisallowedOriginPassesThroughUnstamped(t *testing.T) {
 	h := corsHandler([]string{"http://localhost:4040"})
@@ -68,8 +71,11 @@
 	}
 	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
 		t.Errorf("Allow-Origin stamped for disallowed origin: %q", got)
 	}
+	if got := rec.Header().Get("Vary"); got != "Origin" {
+		t.Errorf("Vary = %q, want Origin on the disallowed-origin pass-through", got)
+	}
 	// A preflight from a disallowed origin reaches the mux (which 404s/405s it);
 	// no CORS approval is expressed either way.
 	rec = doReq(t, h, http.MethodOptions, "http://evil.example", "GET")
 	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {

The pass-through branch also handles no-Origin requests (the "headerless variant" the task's harm statement names), so I'll pin that too.
← Edit internal/httpx/cors_test.go
Index: /root/astrate-mule/internal/httpx/cors_test.go
===================================================================
--- /root/astrate-mule/internal/httpx/cors_test.go
+++ /root/astrate-mule/internal/httpx/cors_test.go
@@ -105,8 +105,11 @@
 	}
 	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
 		t.Errorf("Allow-Origin stamped without Origin header: %q", got)
 	}
+	if got := rec.Header().Get("Vary"); got != "Origin" {
+		t.Errorf("Vary = %q, want Origin on the headerless pass-through", got)
+	}
 }
 
 // OPTIONS without Access-Control-Request-Method is not a preflight and must
 // reach the mux.

$ go test ./internal/httpx/ && gofmt -l . && go vet ./internal/httpx/
ok  	github.com/astrate-platform/astrate/internal/httpx	0.022s

$ go test ./... 2>&1 | grep -Ev "^ok|no test files" || true
(no output)

Done: Added `Vary: Origin` to the CORS pass-through branch in non-wildcard mode, so unstamped responses (disallowed-origin and headerless) carry the RFC 9110 §12.5.5 variation marker; pinned in tests.
Files: internal/httpx/cors.go, internal/httpx/cors_test.go
Verified: `go test ./...`, `go vet ./internal/httpx/`, `gofmt -l .` -> pass
Unsure: none (variation applies to the whole pass-through branch, including no-Origin, matching the "headerless variant" harm the task names)
Follow-ups: none
