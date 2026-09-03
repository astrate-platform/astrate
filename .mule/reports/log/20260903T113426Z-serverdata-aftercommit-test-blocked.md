slug: serverdata-aftercommit-test
verdict: blocked
at:  bde026d
ran: 2026-09-03T11:34:26Z on DietPi in 114s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/engine/serverdata_test.go
→ Read internal/engine/serverdata.go
I'll start by reading the relevant source files to understand the code patterns.
→ Read internal/engine/batch.go [offset=155, limit=25]
Now let me understand the `afterCommit` observer and `newWiredRig` patterns.
✱ Grep "afterCommit" in internal/engine · 15 matches
✱ Grep "func newWiredRig" in internal/engine · 1 match
→ Read internal/engine/batch_test.go [offset=400, limit=50]
→ Read internal/engine/fakestore_test.go [offset=739, limit=60]
→ Read internal/engine/router.go [offset=100, limit=20]
Now I have enough context. `PublishServerValue` (serverdata.go:99-162) persists and publishes directly without going through the batch, so `afterCommit` should NOT fire. Let me write the test.
→ Read internal/engine/serverdata_test.go [limit=10, offset=303]

← Edit internal/engine/serverdata_test.go
Index: /root/astrate-mule/internal/engine/serverdata_test.go
===================================================================
--- /root/astrate-mule/internal/engine/serverdata_test.go
+++ /root/astrate-mule/internal/engine/serverdata_test.go
@@ -309,4 +309,25 @@
 	if !errors.Is(err, ErrPathNotFound) {
 		t.Errorf("bad prefix err = %v, want ErrPathNotFound", err)
 	}
 }
+
+// TestPublishServerValueNoAfterCommit: PublishServerValue persists and
+// publishes directly, bypassing the batch. The afterCommit observer must
+// not fire — only the batch path (batch.go finalize) feeds it.
+func TestPublishServerValueNoAfterCommit(t *testing.T) {
+	rig, _, _ := newWiredRig(t, Config{})
+	ctx := context.Background()
+
+	var observed []PersistOp
+	rig.e.afterCommit = func(ops []PersistOp) { observed = append(observed, ops...) }
+
+	if err := rig.e.PublishServerValue(ctx, realmAlpha, devAlpha,
+		"com.astrate.test.ServerProperties", "/limits/maxConnections",
+		json.RawMessage("42"), nil); err != nil {
+		t.Fatalf("PublishServerValue: %v", err)
+	}
+
+	if len(observed) != 0 {
+		t.Errorf("afterCommit fired %d ops from PublishServerValue, want 0 (bypasses batch)", len(observed))
+	}
+}

Now run the gate.
$ go test -run TestPublishServerValueNoAfterCommit ./internal/engine/ -v 2>&1
=== RUN   TestPublishServerValueNoAfterCommit
--- PASS: TestPublishServerValueNoAfterCommit (0.01s)
PASS
ok  	github.com/astrate-platform/astrate/internal/engine	0.044s

$ gofmt -l . 2>&1
(no output)

$ go vet ./... 2>&1
(no output)

$ go test ./... 2>&1
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	1.485s
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
?   	github.com/astrate-platform/astrate/internal/housekeeping	[no test files]
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
?   	github.com/astrate-platform/astrate/internal/swagger	[no test files]
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

All gates pass. Here's the report:

```
Done: added TestPublishServerValueNoAfterCommit confirming PublishServerValue bypasses the batch afterCommit hook
Files: internal/engine/serverdata_test.go
Verified: go vet ./... -> pass, go test ./... -> pass, gofmt -l . -> clean
Unsure: nothing
Follow-ups: none
```
