slug: flow-validate-deadcode
verdict: blocked
at:  4c636be
ran: 2026-09-02T17:42:30Z on DietPi in 223s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/flow/pipeline.go [offset=100, limit=70]
→ Read internal/flow/pipeline.go [offset=70, limit=35]
Now I'll remove the dead loop (117-123) and the redundant source/sink loop (124-131). Let me look at the exact text to remove.

← Edit internal/flow/pipeline.go
Index: /root/astrate-mule/internal/flow/pipeline.go
===================================================================
--- /root/astrate-mule/internal/flow/pipeline.go
+++ /root/astrate-mule/internal/flow/pipeline.go
@@ -113,23 +113,8 @@
 if visited != len(p.Blocks) {
 	return fmt.Errorf("flow: pipeline %q contains a cycle", p.ID)
 }
 
-hasSource := false
-hasSink := false
-for _, b := range p.Blocks {
-	if inDeg[b.Name] == 0 || outDeg[b.Name] == 0 {
-		// Re-check using original degrees (topo sort modified inDeg).
-	}
-}
-for _, b := range p.Blocks {
-	if outDeg[b.Name] == 0 {
-		hasSink = true
-	}
-	if inDeg[b.Name] == 0 {
-		hasSource = true
-	}
-}
 // Re-compute source/sink from original edges since topo sort zeroed inDeg.
 inDeg2 := make(map[string]int, len(p.Blocks))
 outDeg2 := make(map[string]int, len(p.Blocks))
 for _, b := range p.Blocks {

$ ls internal/flow/
block.go
flow.go
flow_test.go
graph.go
manager_test.go
message.go
message_test.go
pipeline.go
pipeline_test.go
router.go
router_test.go

Now let me look at the existing tests to see coverage and style.
✱ Grep "no source|no sink|no source block|no sink block|cycle" in internal/flow · 12 matches
→ Read internal/flow/pipeline_test.go [offset=85, limit=90]
The task asks to add a test for a pipeline with no source and no sink (already covered) and one where a cycle also lacks sources, to confirm the error message is correct.

Let me check the test structure more — the current table only checks `wantErr` boolean. The task wants to confirm the *error message* is correct. So I need tests that assert the specific error message.

For "no source and no sink (already covered)" — the "no source" test at 134 uses blocks a,b with connections a→b, b→a. That's actually a cycle too. When a cycle also lacks sources, which error comes first? Looking at Validate: the cycle check (line 113-114) runs *before* source/sink check. So for a cycle that lacks sources, the error would be "contains a cycle", not "no source block".

So the task wants to confirm/nail down the correct error message. Let me write tests that assert exact error strings. Let me check if there's an existing error-message assertion pattern in the test file.
→ Read internal/flow/pipeline_test.go [offset=1, limit=85]
Let me understand the semantics. Looking at Validate:

- Cycle check comes first (line 113-114): returns "contains a cycle" if `visited != len(blocks)`.
- Then source/sink check.

The existing "no source" and "no sink" tests both use a 2-cycle (a→b, b→a), which is itself a cycle. So those would hit the "contains a cycle" error, not the source/sink error. In a pure 2-cycle with all blocks in cycle, `visited` stays 0 (no in-degree 0 initially), so it errors with "contains a cycle".

The task asks for tests to *confirm the error message is correct*. The key scenario: a pipeline where a cycle also lacks sources. Because cycle check runs first, the error should be "contains a cycle".

Let me write a dedicated test asserting exact error messages. The current table tests use `wantErr bool`. I'll add error-message assertions. Let me add separate test cases with a `wantErrStr` field, or add a new dedicated test function. Given table style, I'll add a `wantSubstring` field to assert message.

Actually, the cleanest approach: add exact-message assertion tests. Let me add a new test function `TestPipeline_ValidateMessages` that verifies specific error strings for these scenarios. But better to keep it consistent. Let me modify the table to add a `wantErrMsg` field and assert it.

Let me make the change to add message assertions. The task says "Add a test for a pipeline with no source and no sink (already covered) and one for a pipeline where a cycle also lacks sources (to confirm the error message is correct)."

So the tests to confirm:
1. Pipeline with no source and no sink → confirm the message. Since blocks without connections would be both source and sink... wait "no source and no sink" means every block has incoming AND every block has outgoing. That means a graph where all blocks form cycles. That's a cycle → "contains a cycle" error (cycle check first).

2. Cycle that also lacks sources → "contains a cycle" error.

Both scenarios give "contains a cycle" since cycle check runs first. Let me write tests asserting the exact message.

Let me add a `wantErrMsg` field to the table to keep style consistent.The existing tests only check `wantErr bool`. I'll add an exact-message assertion to confirm which error is produced for these scenarios. Let me update the table struct and add relevant cases.

← Edit internal/flow/pipeline_test.go
Index: /root/astrate-mule/internal/flow/pipeline_test.go
===================================================================
--- /root/astrate-mule/internal/flow/pipeline_test.go
+++ /root/astrate-mule/internal/flow/pipeline_test.go
@@ -6,11 +6,12 @@
 )
 
 func TestPipeline_Validate(t *testing.T) {
 	tests := []struct {
-		name    string
-		p       Pipeline
-		wantErr bool
+		name       string
+		p          Pipeline
+		wantErr    bool
+		wantErrMsg string
 	}{
 		{
 			name: "valid linear pipeline",
 			p: Pipeline{

Now update the "no source", "no sink" cases and add a cycle-lacks-sources case. Since these existing cases use a 2-cycle, they actually produce "contains a cycle". Let me update them to assert the correct message, and add the new cycle-lacks-sources case.
← Edit internal/flow/pipeline_test.go
Index: /root/astrate-mule/internal/flow/pipeline_test.go
===================================================================
--- /root/astrate-mule/internal/flow/pipeline_test.go
+++ /root/astrate-mule/internal/flow/pipeline_test.go
@@ -144,9 +144,10 @@
 				{From: "a", To: "b"},
 				{From: "b", To: "a"},
 			},
 		},
-		wantErr: true,
+		wantErr:    true,
+		wantErrMsg: `flow: pipeline "no-source" contains a cycle`,
 	},
 	{
 		name: "no sink (all blocks have outgoing edges)",
 		p: Pipeline{
@@ -159,10 +160,28 @@
 				{From: "a", To: "b"},
 				{From: "b", To: "a"},
 			},
 		},
-		wantErr: true,
+		wantErr:    true,
+		wantErrMsg: `flow: pipeline "no-sink" contains a cycle`,
 	},
+	{
+		name: "cycle that also lacks sources reports the cycle error",
+		p: Pipeline{
+			ID: "cycle-no-source",
+			Blocks: []PipelineNode{
+				{Name: "a", BlockType: "x"},
+				{Name: "b", BlockType: "x"},
+				{Name: "c", BlockType: "x"},
+			},
+			Connections: []Connection{
+				{From: "b", To: "c"},
+				{From: "c", To: "a"},
+			},
+		},
+		wantErr:    true,
+		wantErrMsg: `flow: pipeline "cycle-no-source" contains a cycle`,
+	},
 }
 
 for _, tt := range tests {
 	t.Run(tt.name, func(t *testing.T) {

Wait, let me reconsider the "cycle that also lacks sources" case. Blocks a, b, c with edges b→c, c→a. That's not a cycle (a has no incoming... wait, a has incoming from c). Let me trace: edges b→c, c→a. inDeg: a=1 (from c), b=0, c=1 (from b). So b has in-degree 0. Topo sort: start b, then c, then a. visited=3 = len. No cycle. It has a source (b). This isn't a cycle nor lacks a source.

Let me fix this case to actually be a cycle. A cycle that lacks sources means all blocks are in a cycle (no in-degree 0). That's e.g. a→b, b→c, c→a (the existing "cycle" test at line 117). But that doesn't "lack sources" in a meaningful additional way — it IS the cycle case.

Hmm, the task says "a pipeline where a cycle also lacks sources". Let me think about what makes this distinct. Consider a graph that has BOTH a valid acyclic source path AND a cycle somewhere else. E.g., a source feeding a cycle. Then there IS a source, but there's a cycle → cycle error. That's not "lacks sources."

Actually "lack sources" in the sense that every block has incoming edges. Consider: a→b, b→c, c→a, plus... to also lack a source while being a cycle, it's exactly the full-cycle case. But that's the same as the existing "cycle" test.

Wait — let me reconsider. Maybe the point of the task: the existing "no source" test uses a 2-cycle (a→b, b→a), where the actual reported error is "contains a cycle" not "no source block". The task wants to confirm the message. And "a pipeline where a cycle also lacks sources" — to confirm the cycle error fires before/over the source error.

Let me reconsider the cycle-no-source case to make it genuinely a cycle. A cycle where every block also has no source = full cycle. So use a→b→c→a. But that's identical to the existing "cycle" test. The distinction the task wants is: the case that VERIFIES cycle error wins over source error.

Actually re-reading: "one for a pipeline where a cycle also lacks sources (to confirm the error message is correct)". I think the intent is a graph that is a cycle AND has no source/sink, confirming the cycle message. But since cycle-check runs before source/sink regardless of whether it lacks sources, any cycle confirms it.

Let me make the cycle-no-source case an actual cycle (a→b, b→c, c→a) but ALSO verify the message. However, the existing "cycle" test case at line 117 doesn't assert a message. Let me check it — it's separate. Let me make my new case a genuine cycle and assert its message.

Actually, let me reconsider: to make it also "lack sources", I want every block with incoming edges (a full cycle achieves that). Let me use the full cycle a→b→c→a for the new case. But that duplicates the "cycle" test structurally. That's fine — the point is the message assertion.

Let me fix my case to be a real cycle.
← Edit internal/flow/pipeline_test.go
Index: /root/astrate-mule/internal/flow/pipeline_test.go
===================================================================
--- /root/astrate-mule/internal/flow/pipeline_test.go
+++ /root/astrate-mule/internal/flow/pipeline_test.go
@@ -173,8 +173,9 @@
 		{Name: "b", BlockType: "x"},
 		{Name: "c", BlockType: "x"},
 	},
 	Connections: []Connection{
+		{From: "a", To: "b"},
 		{From: "b", To: "c"},
 		{From: "c", To: "a"},
 	},
 },

Now update the test loop to assert the message when `wantErrMsg` is set.
← Edit internal/flow/pipeline_test.go
Index: /root/astrate-mule/internal/flow/pipeline_test.go
===================================================================
--- /root/astrate-mule/internal/flow/pipeline_test.go
+++ /root/astrate-mule/internal/flow/pipeline_test.go
@@ -189,8 +189,13 @@
 			err := tt.p.Validate()
 			if (err != nil) != tt.wantErr {
 				t.Errorf("Validate() err = %v, wantErr %v", err, tt.wantErr)
 			}
+			if tt.wantErrMsg != "" {
+				if err == nil || err.Error() != tt.wantErrMsg {
+					t.Errorf("Validate() err = %v, want %q", err, tt.wantErrMsg)
+				}
+			}
 		})
 	}
 }
 

Now let me run the gate.
$ go vet ./... && go test ./... && gofmt -l .
# github.com/astrate-platform/astrate/internal/flow
# [github.com/astrate-platform/astrate/internal/flow]
vet: internal/flow/pipeline.go:128:2: undefined: hasSource

I left a reference to `hasSource` somewhere. Let me find it.
→ Read internal/flow/pipeline.go [offset=113, limit=35]
The reassignment `hasSource = false`/`hasSink = false` (lines 128-129) previously relied on the declarations in the removed loops. I need to change these to declarations.

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

$ go vet ./... && go test ./... && gofmt -l .
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
ok  	github.com/astrate-platform/astrate/internal/flow	0.253s
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

All pass, gofmt clean. Let me verify the tests actually exercise the new assertion meaningfully—confirm the new cycle-lacks-sources test was hit. Let me run the flow tests with verbose for the specific titles.
$ go test ./internal/flow/ -run TestPipeline_Validate -v 2>&1 | grep -E "no source|no sink|cycle"
=== RUN   TestPipeline_Validate/cycle
=== RUN   TestPipeline_Validate/cycle_that_also_lacks_sources_reports_the_cycle_error
    --- PASS: TestPipeline_Validate/cycle (0.00s)
    --- PASS: TestPipeline_Validate/cycle_that_also_lacks_sources_reports_the_cycle_error (0.00s)

The new tests pass with the implementation in place. Per MULE.md, the runner will remove my implementation and verify my message assertions break. Let me confirm that my new message assertions genuinely depend on the change being correct (they do — they assert exact strings produced by the kept recomputed check). But note: my dead-code removal doesn't change behaviour — the source/sink error messages come from the kept recomputed loop (lines 133-158), which I didn't change. The tests assert messages that were already produced before my change.

This is a concern per MULE.md rule 4: "a test that passes without your change is rejected." My change only removes dead/redundant loops — it changes NO behaviour. The error messages the tests assert were true before the change too. So these tests would pass without my change, which the runner would reject.

The task explicitly asks: "Add a test for a pipeline with no source and no sink (already covered) and one for a pipeline where a cycle also lacks sources (to confirm the error message is correct)." The task itself frames these as confirming already-correct behavior. This is a deadcode removal task — the "behaviour change" is the removal of dead loops.

Hmm. Let me reconsider whether any prior code path actually produced a different result. Before my change, the dead loop (117-123) did nothing. The redundant loop (124-131) used `inDeg`/`outDeg` which were zeroed by topo sort — so `hasSource`/`hasSink` from that loop: after topo sort, inDeg are all 0 (zeroed), so `hasSource` would be true, `hasSink` true if any outDeg==0 (outDeg not modified by topo). Then the recomputed loop overwrote `hasSource`/`hasSink`. So the intermediate was irrelevant.

But wait — before my change, in the cycle case, Validate returns at line 114 (cycle error) before reaching the source/sink recompute. So message is "contains a cycle". Same after. The "no source" test (2-cycle) → "contains a cycle". This was true before my change too.

So indeed these message tests pass without my change. The task is fundamentally a deadcode-removal task where no behaviour changes, so no test can fail-without/change. MULE.md says: "If you cannot write such a test, say so explicitly in the report."

I'll do exactly that — note in the report that the removal is behavior-preserving, so the message-assertion tests verify correctness but pass with or without the change (they can't be made to depend on the removal).

Let me double check gofmt is clean and the final state of the decoder. Also verify nothing else references removed vars.

The gate passed (go vet, go test, gofmt -l returned nothing). Work complete. Let me write the report.

```
Done: Removed the dead comment-only loop and redundant source/sink loop in internal/flow/pipeline.go Validate, keeping only the recomputed inDeg2/outDeg2 check; added exact error-message tests for no-source/no-sink and a cycle-that-also-lacks-sources.
Files: internal/flow/pipeline.go, internal/flow/pipeline_test.go
Verified: go vet ./... && go test ./... && gofmt -l . -> pass
Unsure: The removed loops were behaviorally dead (the recomputed inDeg2/outDeg2 loop already overwrote their results, and valid pipelines never reached them), so no test can fail without the removal — the message-assertion tests confirm correct output but pass with or without this change.
Follow-ups: none
```
