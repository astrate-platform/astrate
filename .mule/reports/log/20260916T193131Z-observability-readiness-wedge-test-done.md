slug: observability-readiness-wedge-test
verdict: done
at:  4ef2e23
ran: 2026-09-16T19:31:31Z on DietPi in 178s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/observability/health.go
→ Read internal/observability/observability_test.go
✱ Grep "NewHealth|&Health\{|Health\{" in . · 14 matches
✱ Grep "readinessTimeout" in . · 8 matches

← Edit internal/observability/health.go
Index: /root/astrate-mule/internal/observability/health.go
===================================================================
--- /root/astrate-mule/internal/observability/health.go
+++ /root/astrate-mule/internal/observability/health.go
@@ -20,8 +20,9 @@
 // /astrate/v1 so they never collide with the upstream API namespace.
 type Health struct {
 	metrics http.Handler
 	ready   []namedCheck
+	timeout time.Duration
 }
 
 type namedCheck struct {
 	name  string

← Edit internal/observability/health.go
Index: /root/astrate-mule/internal/observability/health.go
===================================================================
--- /root/astrate-mule/internal/observability/health.go
+++ /root/astrate-mule/internal/observability/health.go
@@ -30,9 +30,9 @@
 }
 
 // NewHealth builds the surface around an existing metrics handler.
 func NewHealth(metrics http.Handler) *Health {
-	return &Health{metrics: metrics}
+	return &Health{metrics: metrics, timeout: readinessTimeout}
 }
 
 // AddReadiness registers a named dependency probed by /readiness.
 func (h *Health) AddReadiness(name string, check Check) {

The task: add an injectable `timeout` field on `Health` defaulting to `readinessTimeout`, and prove a ctx-honoring wedged check returns 503 within budget instead of hanging. `Health` is only built via `NewHealth`, so this is safe.
← Edit internal/observability/health.go
Index: /root/astrate-mule/internal/observability/health.go
===================================================================
--- /root/astrate-mule/internal/observability/health.go
+++ /root/astrate-mule/internal/observability/health.go
@@ -53,9 +53,9 @@
 
 // handleReadiness runs every dependency check; any failure yields 503 with the
 // per-check status so operators see which dependency is down.
 func (h *Health) handleReadiness(w http.ResponseWriter, r *http.Request) {
-	ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
+	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
 	defer cancel()
 
 	checks := make(map[string]string, len(h.ready))
 	ready := true

← Edit internal/observability/observability_test.go
Index: /root/astrate-mule/internal/observability/observability_test.go
===================================================================
--- /root/astrate-mule/internal/observability/observability_test.go
+++ /root/astrate-mule/internal/observability/observability_test.go
@@ -6,8 +6,9 @@
 	"net/http"
 	"net/http/httptest"
 	"strings"
 	"testing"
+	"time"
 )
 
 func TestMetricsExposesGauges(t *testing.T) {
 	m := NewMetrics()

✗ Edit internal/observability/observability_test.go failed
Error: Could not find oldString in the file. It must match exactly, including whitespace, indentation, and line endings.
Now the test — a ctx-honoring wedged check that only unblocks on ctx cancellation, with a short injected budget:

← Edit internal/observability/observability_test.go
Index: /root/astrate-mule/internal/observability/observability_test.go
===================================================================
--- /root/astrate-mule/internal/observability/observability_test.go
+++ /root/astrate-mule/internal/observability/observability_test.go
@@ -62,8 +62,42 @@
 		t.Errorf("readiness body should name the failing check: %s", rec.Body)
 	}
 }
 
+func TestReadinessWedgedCheckBoundedByBudget(t *testing.T) {
+	// A ctx-honoring wedged dependency blocks until the deadline; the
+	// readiness endpoint must return 503 within the injected budget instead
+	// of hanging for the full default readinessTimeout.
+	const budget = 50 * time.Millisecond
+
+	h := NewHealth(NewMetrics().Handler())
+	h.timeout = budget
+	// Context deadline exceeded is the report a wedged-but-cooperative
+	// dependency gives when its probe is cancelled.
+	h.AddReadiness("database", func(ctx context.Context) error {
+		<-ctx.Done()
+		return ctx.Err()
+	})
+	mux := http.NewServeMux()
+	h.Mount(mux)
+
+	start := time.Now()
+	rec := httptest.NewRecorder()
+	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/readiness", nil))
+	elapsed := time.Since(start)
+
+	if rec.Code != http.StatusServiceUnavailable {
+		t.Errorf("readiness = %d, want 503", rec.Code)
+	}
+	if !strings.Contains(rec.Body.String(), "database") ||
+		!strings.Contains(rec.Body.String(), "deadline exceeded") {
+		t.Errorf("readiness body should report the wedged check: %s", rec.Body)
+	}
+	if elapsed > budget*5 {
+		t.Errorf("readiness took %s, want it bounded by injected budget %s", elapsed, budget)
+	}
+}
+
 func TestReadinessAllOK(t *testing.T) {
 	h := NewHealth(NewMetrics().Handler())
 	h.AddReadiness("database", func(context.Context) error { return nil })
 	mux := http.NewServeMux()

I swapped the parameters. Let me redo it correctly:
$ gofmt -l . && go vet ./internal/observability/ && go test ./internal/observability/
ok  	github.com/astrate-platform/astrate/internal/observability	0.131s

Local package clean. Now the full gate:
$ go vet ./... && go test ./...
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
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
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

Gate is green. Done.

```
Done: Made the readiness budget injectable — added `timeout time.Duration` on `Health` defaulting to `readinessTimeout`, wired into `handleReadiness`, plus a test proving a ctx-honoring wedged check yields 503 within a short injected budget instead of hanging.
Files: internal/observability/health.go, internal/observability/observability_test.go
Verified: gofmt -l . -> clean; go vet ./... && go test ./... -> pass (go test, no -race — no container/race env needed; no goroutines touched)
Unsure: nothing
Follow-ups: none
```
