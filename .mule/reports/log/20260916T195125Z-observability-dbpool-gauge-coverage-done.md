slug: observability-dbpool-gauge-coverage
verdict: done
at:  19cd4f1
ran: 2026-09-16T19:51:25Z on DietPi in 62s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/observability/observability_test.go
→ Read internal/observability/metrics.go
Task: add `idle_conns 5` and `total_conns 7` assertions. Straightforward.

← Edit internal/observability/observability_test.go
Index: /root/astrate-mule/internal/observability/observability_test.go
===================================================================
--- /root/astrate-mule/internal/observability/observability_test.go
+++ /root/astrate-mule/internal/observability/observability_test.go
@@ -25,8 +25,10 @@
 body := rec.Body.String()
 for _, want := range []string{
 	"astrate_broker_sessions 3",
 	"astrate_db_pool_acquired_conns 2",
+	"astrate_db_pool_idle_conns 5",
+	"astrate_db_pool_total_conns 7",
 	"astrate_db_pool_max_conns 10",
 	"go_goroutines", // runtime collector present
 } {
 	if !strings.Contains(body, want) {

$ go test ./internal/observability/ && gofmt -l internal/observability/
ok  	github.com/astrate-platform/astrate/internal/observability	0.076s

Done: pinned `idle_conns` (5) and `total_conns` (7) assertions alongside the existing acquired/max pins in `TestMetricsExposesGauges`.
Files: internal/observability/observability_test.go
Verified: go test ./internal/observability/ -> pass; gofmt -l internal/observability/ -> clean
Unsure: nothing
Follow-ups: none
