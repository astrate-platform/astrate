slug: flow-validate-dead-source-sink-recompute
verdict: done
at:  7445176
ran: 2026-09-05T20:15:29Z on DietPi in 200s


> build · big-pickle

I'll read the two files to understand the task.
→ Read .mule/MULE.md
→ Read .mule/task.md
✗ Read internal/pipeline/pipeline.go failed [offset=90, limit=80]
Error: File not found: /root/astrate-mule/internal/pipeline/pipeline.go
✱ Glob "**/pipeline.go" 1 match
→ Read internal/flow/pipeline.go [limit=80, offset=90]
I'll delete the first dead block (lines 117-127) as the task specifies.

← Edit internal/flow/pipeline.go
Index: /root/astrate-mule/internal/flow/pipeline.go
===================================================================
--- /root/astrate-mule/internal/flow/pipeline.go
+++ /root/astrate-mule/internal/flow/pipeline.go
@@ -113,19 +113,10 @@
 if visited != len(p.Blocks) {
 	return fmt.Errorf("flow: pipeline %q contains a cycle", p.ID)
 }
 
-hasSource := false
-hasSink := false
-for _, b := range p.Blocks {
-	if outDeg[b.Name] == 0 {
-		hasSink = true
-	}
-	if inDeg[b.Name] == 0 {
-		hasSource = true
-	}
-}
-// Re-compute source/sink from original edges since topo sort zeroed inDeg.
+// Source/sink are computed from the original edges below since the
+// topological sort above zeroed inDeg.
 inDeg2 := make(map[string]int, len(p.Blocks))
 outDeg2 := make(map[string]int, len(p.Blocks))
 for _, b := range p.Blocks {
 	inDeg2[b.Name] = 0

Let me keep the original comment form for a minimal diff.
← Edit internal/flow/pipeline.go
Index: /root/astrate-mule/internal/flow/pipeline.go
===================================================================
--- /root/astrate-mule/internal/flow/pipeline.go
+++ /root/astrate-mule/internal/flow/pipeline.go
@@ -113,10 +113,9 @@
 if visited != len(p.Blocks) {
 	return fmt.Errorf("flow: pipeline %q contains a cycle", p.ID)
 }
 
-// Source/sink are computed from the original edges below since the
-// topological sort above zeroed inDeg.
+// Re-compute source/sink from original edges since topo sort zeroed inDeg.
 inDeg2 := make(map[string]int, len(p.Blocks))
 outDeg2 := make(map[string]int, len(p.Blocks))
 for _, b := range p.Blocks {
 	inDeg2[b.Name] = 0

Now verify the result and run the gates.
✱ Grep "TestPipeline_Validate" in internal/flow · 1 match
→ Read internal/flow/pipeline.go [limit=40, offset=110]
The second block's `hasSource = false` now refers to undeclared vars — that must become a declaration. Fixing:

← Edit internal/flow/pipeline.go
Index: /root/astrate-mule/internal/flow/pipeline.go
===================================================================
--- /root/astrate-mule/internal/flow/pipeline.go
+++ /root/astrate-mule/internal/flow/pipeline.go
@@ -124,10 +124,10 @@
 for _, c := range p.Connections {
 	outDeg2[c.From]++
 	inDeg2[c.To]++
 }
-hasSource = false
-hasSink = false
+hasSource := false
+hasSink := false
 for _, b := range p.Blocks {
 	if inDeg2[b.Name] == 0 {
 		hasSource = true
 	}

→ Read internal/flow/pipeline.go [limit=25, offset=116]
Now run the gates.

$ go vet ./... && gofmt -l . && go test ./internal/flow/ -count=1
ok  	github.com/astrate-platform/astrate/internal/flow	0.499s

$ go test ./... -count=1
ok  	github.com/astrate-platform/astrate/cmd/astrate	0.546s
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
ok  	github.com/astrate-platform/astrate/internal/appengine	0.109s
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	41.874s
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	2.502s
ok  	github.com/astrate-platform/astrate/internal/auth	4.609s
ok  	github.com/astrate-platform/astrate/internal/broker	2.651s
ok  	github.com/astrate-platform/astrate/internal/config	0.071s
ok  	github.com/astrate-platform/astrate/internal/engine	2.075s
ok  	github.com/astrate-platform/astrate/internal/engine/forward	0.210s
ok  	github.com/astrate-platform/astrate/internal/engine/stream	0.088s
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	2.644s
ok  	github.com/astrate-platform/astrate/internal/flow	0.829s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	1.476s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	0.115s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.508s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	0.084s
ok  	github.com/astrate-platform/astrate/internal/flowapi	0.183s
?   	github.com/astrate-platform/astrate/internal/housekeeping	[no test files]
ok  	github.com/astrate-platform/astrate/internal/httpx	0.055s
ok  	github.com/astrate-platform/astrate/internal/observability	0.113s
ok  	github.com/astrate-platform/astrate/internal/pairing	16.288s
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	27.724s
ok  	github.com/astrate-platform/astrate/internal/realm	0.046s
ok  	github.com/astrate-platform/astrate/internal/store	0.028s
ok  	github.com/astrate-platform/astrate/internal/swagger	0.036s
ok  	github.com/astrate-platform/astrate/internal/testutil	0.203s
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	0.206s
ok  	github.com/astrate-platform/astrate/pkg/deviceid	0.033s
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	0.167s
ok  	github.com/astrate-platform/astrate/pkg/payload	0.200s

All green. Task complete.

```
Done: Deleted the dead first hasSource/hasSink computation in Pipeline.Validate; the recompute block now declares and fills them (was the only truly-relevant issue to fix, so the delete is behaviour-identical).
Files: internal/flow/pipeline.go
Verified: go vet ./...; gofmt -l .; go test ./... -count=1 -> pass (all three)
Unsure: nothing (the second block's assignment needed to become a declaration once the first block was gone; confirming it was the only residual reference — it was)
Follow-ups: none
```
