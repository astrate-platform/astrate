slug: observability-health-content-type-test
verdict: done
at:  f6c27f5
ran: 2026-09-16T20:28:28Z on DietPi in 250s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
✱ Glob "**/health.go" 1 match
✱ Glob "**/compat.go" 2 matches
✱ Glob "**/compat_version_test.go" 1 match
→ Read internal/observability/compat.go
→ Read internal/observability/health.go
→ Read internal/observability/compat_version_test.go
✱ Glob "internal/observability/*_test.go" 2 matches
→ Read internal/observability/observability_test.go
Task: add `Content-Type: application/json` assertions to the observability health tests so a regression dropping the header on the `/astrate/v1` surface (and per-service `/health` 200/503) can't pass. Let me check the repo history and todo for conventions around test-only tasks.
→ Read .mule/todo.md

$ git log --oneline -15 -- internal/observability
19cd4f1 mule: observability-dbpool-gauge-coverage [auto]: extend `TestMetricsExposesGauges` (internal/observability/observability_test.go:25-34) to assert all four `astrate_db_pool_*` gauges — it currently pins only `acquired_conns` and `max_conns`, so a regression deleting or mis-wiring `idle_conns` or `total_conns` (metrics.go:77-78) passes the suite; expect idle 5, total 7 from the supplied `DBPoolStats{2,5,7,10}`.
4ef2e23 mule: observability-readiness-wedge-test [auto]: make the readiness budget injectable and pinned — add a `timeout time.Duration` field on `Health` (internal/observability/health.go) defaulting to `readinessTimeout` (health.go:12, no test currently exercises the wedged-dependency path), and add a test in observability_test.go asserting a ctx-honoring wedged check makes `/astrate/v1/readiness` return 503 reporting it failing within a short injected budget, instead of hanging the endpoint.
b6e1ca4 compat: make the dashboard health indicators real, pin always-sync (#68, #91)
cda4010 feat(http): per-service version endpoints, upstream-parity (#77)
699cf11 Initial commit under updated license terms

Confirmed: test-only tasks are a normal pattern here (e.g. `pairing-bearer-secret-test`, `store-alias-values-taken-test`). Now I'll add `Content-Type: application/json` assertions to the tests that already exercise each path.

The four paths: `/astrate/v1/health` 200 (observability_test.go:52-55), `/astrate/v1/readiness` 503 (:57-62) and 200 (:108-112), per-service `/health` 200 and 503 (:115-153).
← Edit internal/observability/observability_test.go
Index: /root/astrate-mule/internal/observability/observability_test.go
===================================================================
--- /root/astrate-mule/internal/observability/observability_test.go
+++ /root/astrate-mule/internal/observability/observability_test.go
@@ -52,8 +52,11 @@
 mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/health", nil))
 if rec.Code != http.StatusOK {
 	t.Errorf("health = %d, want 200", rec.Code)
 }
+if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
+	t.Errorf("health Content-Type = %q, want application/json", ct)
+}
 
 // One failing dependency makes readiness 503.
 rec = httptest.NewRecorder()
 mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/readiness", nil))

← Edit internal/observability/observability_test.go
Index: /root/astrate-mule/internal/observability/observability_test.go
===================================================================
--- /root/astrate-mule/internal/observability/observability_test.go
+++ /root/astrate-mule/internal/observability/observability_test.go
@@ -65,8 +65,11 @@
 	}
 	if !strings.Contains(rec.Body.String(), "down") {
 		t.Errorf("readiness body should name the failing check: %s", rec.Body)
 	}
+	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
+		t.Errorf("readiness Content-Type = %q, want application/json", ct)
+	}
 }
 
 func TestReadinessWedgedCheckBoundedByBudget(t *testing.T) {
 	// A ctx-honoring wedged dependency blocks until the deadline; the

← Edit internal/observability/observability_test.go
Index: /root/astrate-mule/internal/observability/observability_test.go
===================================================================
--- /root/astrate-mule/internal/observability/observability_test.go
+++ /root/astrate-mule/internal/observability/observability_test.go
@@ -115,8 +115,11 @@
 	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/readiness", nil))
 	if rec.Code != http.StatusOK {
 		t.Errorf("readiness = %d, want 200", rec.Code)
 	}
+	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
+		t.Errorf("readiness Content-Type = %q, want application/json", ct)
+	}
 }
 
 func TestMountServiceCompat(t *testing.T) {
 	svcs := []string{"appengine", "realmmanagement", "pairing"}

← Edit internal/observability/observability_test.go
Index: /root/astrate-mule/internal/observability/observability_test.go
===================================================================
--- /root/astrate-mule/internal/observability/observability_test.go
+++ /root/astrate-mule/internal/observability/observability_test.go
@@ -137,8 +137,11 @@
 		}
 		if body := rec.Body.String(); body != `{"data":{"status":"ok"}}` {
 			t.Errorf("%s/health body = %s", svc, body)
 		}
+		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
+			t.Errorf("%s/health Content-Type = %q, want application/json", svc, ct)
+		}
 	}
 })
 
 // The Dashboard's per-service indicator is only useful if it can go red:

← Edit internal/observability/observability_test.go
Index: /root/astrate-mule/internal/observability/observability_test.go
===================================================================
--- /root/astrate-mule/internal/observability/observability_test.go
+++ /root/astrate-mule/internal/observability/observability_test.go
@@ -160,8 +160,11 @@
 		}
 		if body := rec.Body.String(); body != `{"data":{"status":"unhealthy"}}` {
 			t.Errorf("%s/health body = %s", svc, body)
 		}
+		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
+			t.Errorf("%s/health Content-Type = %q, want application/json", svc, ct)
+		}
 	}
 })
 
 // A nil check keeps upstream's static 200 for callers that have no probe.

Now run the gate. First targeted checks, then the full suite.
$ gofmt -l . && go vet ./internal/observability/ && go test ./internal/observability/
ok  	github.com/astrate-platform/astrate/internal/observability	0.115s

Package clean. Let me check how the runner verifies a test-only change (no behavior diff here), so I know this task shape passes the gate.
$ git diff --stat
internal/observability/observability_test.go | 15 +++++++++++++++
 1 file changed, 15 insertions(+)

✱ Grep "pace|implementation|take their|break|gate" in tools · 56 matches
→ Read tools/mule.sh [offset=160, limit=55]
Proof gate only fires when a diff has both impl and test Go changes; a test-only commit is explicitly "not this gate's business" (mule.sh:169-170). Now the full gate.

$ gofmt -l . && go vet ./... && go test ./... 2>&1 | tail -20
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
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

$ nproc; uname -r; grep -m1 MemTotal /proc/meminfo
4
6.18.39+rpt-rpi-v8
MemTotal:        3886904 kB

On the Pi (4 cores, 3.7GB); per AGENTS.md/MULE.md the gate here is `go vet ./... && go test ./...` — no `-race`. Final diff review:
$ git diff
diff --git a/internal/observability/observability_test.go b/internal/observability/observability_test.go
index c52d567..ed2be3e 100644
--- a/internal/observability/observability_test.go
+++ b/internal/observability/observability_test.go
@@ -53,6 +53,9 @@ func TestReadiness(t *testing.T) {
 	if rec.Code != http.StatusOK {
 		t.Errorf("health = %d, want 200", rec.Code)
 	}
+	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
+		t.Errorf("health Content-Type = %q, want application/json", ct)
+	}
 
 	// One failing dependency makes readiness 503.
 	rec = httptest.NewRecorder()
@@ -63,6 +66,9 @@ func TestReadiness(t *testing.T) {
 	if !strings.Contains(rec.Body.String(), "down") {
 		t.Errorf("readiness body should name the failing check: %s", rec.Body)
 	}
+	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
+		t.Errorf("readiness Content-Type = %q, want application/json", ct)
+	}
 }
 
 func TestReadinessWedgedCheckBoundedByBudget(t *testing.T) {
@@ -110,6 +116,9 @@ func TestReadinessAllOK(t *testing.T) {
 	if rec.Code != http.StatusOK {
 		t.Errorf("readiness = %d, want 200", rec.Code)
 	}
+	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
+		t.Errorf("readiness Content-Type = %q, want application/json", ct)
+	}
 }
 
 func TestMountServiceCompat(t *testing.T) {
@@ -129,6 +138,9 @@ func TestMountServiceCompat(t *testing.T) {
 			if body := rec.Body.String(); body != `{"data":{"status":"ok"}}` {
 				t.Errorf("%s/health body = %s", svc, body)
 			}
+			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
+				t.Errorf("%s/health Content-Type = %q, want application/json", svc, ct)
+			}
 		}
 	})
 
@@ -149,6 +161,9 @@ func TestMountServiceCompat(t *testing.T) {
 			if body := rec.Body.String(); body != `{"data":{"status":"unhealthy"}}` {
 				t.Errorf("%s/health body = %s", svc, body)
 			}
+			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
+				t.Errorf("%s/health Content-Type = %q, want application/json", svc, ct)
+			}
 		}
 	})

Done. Added `Content-Type: application/json` assertions on both `/astrate/v1/health` (200) and `/astrate/v1/readiness` (200 + 503) paths — which flow through `writeJSON` (health.go:81) — plus the 200 and 503 paths of the per-service `/health` (headers set by hand in compat.go:26,31). All assertions match the exact `application/json` value, not just non-empty, matching the task wording.

```
Done: pinned Content-Type: application/json on /astrate/v1/health, /astrate/v1/readiness (200+503) and per-service /health (200+503) via exact-value assertions in the existing suite
Files: internal/observability/observability_test.go
Verified: gofmt -l . -> pass (clean); go vet ./... -> pass; go test ./... -> pass (all packages ok, incl. internal/observability). No -race: this box is the Pi (4 cores, 3.7GB, kernel 6.18.39+rpt-rpi-v8), where the race detector cannot run and the standing gate is vet+test; change touches no goroutines/locks.
Unsure: nothing
Follow-ups: none
```
