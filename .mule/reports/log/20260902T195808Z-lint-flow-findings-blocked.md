slug: lint-flow-findings
verdict: blocked
at:  154ee25
ran: 2026-09-02T19:58:08Z on DietPi in 260s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/flow/router.go [offset=170, limit=60]
Let me read the files I need to change.
→ Read internal/flow/message.go [offset=15, limit=15]
→ Read internal/flow/flow_test.go [offset=110, limit=20]
Four changes. Let me make them all.

← Edit internal/flow/router.go
Index: /root/astrate-mule/internal/flow/router.go
===================================================================
--- /root/astrate-mule/internal/flow/router.go
+++ /root/astrate-mule/internal/flow/router.go
@@ -211,9 +211,9 @@
 		prime64  = 1099511628211
 	)
 	h := uint64(offset64)
 	for _, b := range key {
-		h ^= uint64(b)
+		h ^= uint64(b) // #nosec G115 -- rune range fits uint64
 		h *= prime64
 	}
 	return int(h % uint64(n)) // #nosec G115 -- value already reduced mod n
 }

→ Read internal/flow/router.go [offset=213, limit=8]

← Edit internal/flow/flow_test.go
Index: /root/astrate-mule/internal/flow/flow_test.go
===================================================================
--- /root/astrate-mule/internal/flow/flow_test.go
+++ /root/astrate-mule/internal/flow/flow_test.go
@@ -114,11 +114,8 @@
 	}
 	if f.StartedAt().IsZero() {
 		t.Error("StartedAt is zero")
 	}
-	if f.StoppedAt().IsZero() {
-		// Not stopped yet; StoppedAt should be zero.
-	}
 }
 
 func TestFlow_StatusAfterStop(t *testing.T) {
 	mgr := NewManager()

← Edit internal/flow/router.go
Index: /root/astrate-mule/internal/flow/router.go
===================================================================
--- /root/astrate-mule/internal/flow/router.go
+++ /root/astrate-mule/internal/flow/router.go
@@ -176,9 +176,9 @@
 }
 
 // runLane drains one lane's channel through the block graph, recovering
 // panics per-message so one block's bug cannot crash the router.
-func (r *Router) runLane(ctx context.Context, l *lane) {
+func (r *Router) runLane(_ context.Context, l *lane) {
 	defer r.laneWG.Done()
 	for fm := range l.ch {
 		r.processOne(fm)
 	}

Good, that's correct. Now the remaining three changes:
← Edit internal/flow/message.go
Index: /root/astrate-mule/internal/flow/message.go
===================================================================
--- /root/astrate-mule/internal/flow/message.go
+++ /root/astrate-mule/internal/flow/message.go
@@ -17,8 +17,9 @@
 // DataType enumerates the wire-level value types Astarte Flow supports.
 type DataType uint8
 
 const (
+	// TypeInteger indicates the message carries an integer payload.
 	TypeInteger DataType = iota
 	TypeReal
 	TypeBoolean
 	TypeDatetime

Now run the gates.
$ cd /root/astrate-mule && gofmt -l . 2>&1
(no output)

$ cd /root/astrate-mule && go vet ./... 2>&1 && go test ./... 2>&1
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
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
ok  	github.com/astrate-platform/astrate/internal/flow	0.250s
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

All green. Report:

```
Done: cleaned up four revive/gosec lint findings in internal/flow
Files: internal/flow/router.go, internal/flow/flow_test.go, internal/flow/message.go
Verified: go vet ./... && go test ./... -> pass; gofmt -l . -> none
Unsure: nothing
Follow-ups: none
```
