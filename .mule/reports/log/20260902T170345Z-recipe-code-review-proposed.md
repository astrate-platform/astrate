slug: recipe-code-review
verdict: proposed
at:  345e129
ran: 2026-09-02T17:03:45Z on DietPi in 453s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
✱ Glob ".mule/reviews/*.md" 0 matches
→ Read .mule/log.md
No previous reviews. Let me survey the packages and pick one.

$ ls /root/astrate-mule/internal/ /root/astrate-mule/pkg/ 2>/dev/null
/root/astrate-mule/internal/:
appengine
auth
broker
config
engine
flow
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil

/root/astrate-mule/pkg/:
astarteapi
deviceid
interfaceschema
payload

$ for dir in /root/astrate-mule/internal/*/ /root/astrate-mule/pkg/*/; do count=$(find "$dir" -name '*.go' -exec wc -l {} + 2>/dev/null | tail -1 | awk '{print $1}'); echo "$count $dir"; done | sort -rn
12413 /root/astrate-mule/internal/engine/
5348 /root/astrate-mule/internal/store/
4885 /root/astrate-mule/internal/appengine/
3441 /root/astrate-mule/internal/broker/
2777 /root/astrate-mule/internal/pairing/
2707 /root/astrate-mule/internal/flow/
2578 /root/astrate-mule/pkg/payload/
2345 /root/astrate-mule/pkg/interfaceschema/
1404 /root/astrate-mule/internal/auth/
1386 /root/astrate-mule/internal/realm/
697 /root/astrate-mule/internal/testutil/
647 /root/astrate-mule/internal/config/
500 /root/astrate-mule/internal/housekeeping/
441 /root/astrate-mule/pkg/astarteapi/
358 /root/astrate-mule/pkg/deviceid/
357 /root/astrate-mule/internal/httpx/
272 /root/astrate-mule/internal/observability/
43 /root/astrate-mule/internal/swagger/

I'll review `internal/flow` — 2707 lines, data flow layer, likely touching concurrency. Let me read all its Go files and tests.
✱ Glob "internal/flow/**/*.go" 11 matches

$ wc -l /root/astrate-mule/internal/flow/*.go | sort -rn
2707 total
  402 /root/astrate-mule/internal/flow/router_test.go
  362 /root/astrate-mule/internal/flow/message.go
  348 /root/astrate-mule/internal/flow/message_test.go
  313 /root/astrate-mule/internal/flow/manager_test.go
  280 /root/astrate-mule/internal/flow/pipeline_test.go
  257 /root/astrate-mule/internal/flow/router.go
  245 /root/astrate-mule/internal/flow/flow.go
  196 /root/astrate-mule/internal/flow/flow_test.go
  170 /root/astrate-mule/internal/flow/pipeline.go
   91 /root/astrate-mule/internal/flow/block.go
   43 /root/astrate-mule/internal/flow/graph.go

→ Read internal/flow/flow.go
→ Read internal/flow/message.go
→ Read internal/flow/router.go
→ Read internal/flow/pipeline.go
→ Read internal/flow/block.go
→ Read internal/flow/graph.go
→ Read internal/flow/router_test.go
→ Read internal/flow/message_test.go
→ Read internal/flow/manager_test.go
Now let me read the test files.
→ Read internal/flow/pipeline_test.go
→ Read internal/flow/flow_test.go
✱ Grep ""internal/flow"" in . · 0 matches
Now let me check how the flow package is used from the rest of the codebase.

$ cd /root/astrate-mule && go test ./internal/flow/... 2>&1 | tail -5
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)

$ ls /root/astrate-mule/.mule/reviews/ 2>/dev/null || echo "(directory does not exist)"
realm-2026-09-01.md
store-2026-07-27.md

→ Read .mule/todo.md
Now I have a complete picture. Let me write the review and task lines.

← Write .mule/reviews/flow-2026-09-02.md
Wrote file successfully.

Now append task lines to `.mule/todo.md`:
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -109,4 +109,8 @@
 - [!] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto] — BLOCKED: gates failed
 - [!] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
 - [!] realm-policy-list-sorted: in `internal/realm/service.go`, sort the names returned by `ListPolicies` the way `ListTriggers`/`ListInterfaces`/`ListInterfaceMajors` already do (it is the only list method without `sort.Strings`), and extend `TestDashboardCompat.Policies` (or a new subtest in `http_test.go`) to install several policies in non-alphabetical order and assert `GET /policies` comes back sorted. [auto] — BLOCKED: gates failed
 - [!] realm-interface-lookup-404: in `internal/realm/http_test.go` `TestRealmManagement`, assert the 404 paths that nothing currently covers — `GET /interfaces/<nonexistent-name>` and `GET /interfaces/<name>/<no-such-major>` (service.go `ListInterfaceMajors`:211-213 and `GetInterface`, mapped to 404 via writeError). A rule with no test currently. [auto] — BLOCKED: gates failed
+- [ ] flow-setstatus-race: in `internal/flow/flow.go` `Manager.StartFlow`, wrap the `f.setStatus(FlowStatusFailed)` call at line 157 in `f.mu.Lock()`/`f.mu.Unlock()` to match the contract documented on `setStatus`. Add a test that starts a flow, confirms it is in the manager's map, then concurrently calls `f.Status()` while `StartFlow` sets the failed status — if the race detector is unavailable, at minimum confirm the lock is acquired. [auto]
+- [ ] flow-validate-deadcode: in `internal/flow/pipeline.go` `Validate`, remove the dead loop at lines 117-123 (the one whose body is a comment) and the redundant source/sink loop at lines 124-131, keeping only the recomputed inDeg2/outDeg2 check at lines 133-158. Add a test for a pipeline with no source and no sink (already covered) and one for a pipeline where a cycle also lacks sources (to confirm the error message is correct). [auto]
+- [ ] flow-unmarshal-error-tests: in `internal/flow/message_test.go`, add table-driven tests for `UnmarshalJSON` error paths: missing key, unknown type string, map type field with non-string value, and map data with a field absent from FieldTypes. Each should assert the expected error substring. [auto]
+- [ ] flow-datawirescalar-fallthrough: in `internal/flow/message.go` `dataWireScalar`, the default case (line 208) returns `m.Data` raw, which is correct after UnmarshalJSON but undocumented. Add a comment documenting the invariant that `Data` matches `Type` after construction, and add a test in `message_test.go` that constructs a FlowMessage with mismatched Type/Data and verifies the wire output is the raw Go value (proving the invariant, not fixing it). [auto]

Now add the flagged design question to `for-giulio.md`:
→ Read .mule/for-giulio.md

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,18 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **`Router.Submit` TOCTOU on `closed` flag** (`internal/flow/router.go:113-120`): Submit
+  reads `r.closed` under the mutex, drops the lock, then sends on the channel. A concurrent
+  `Drain` could close the channel between the unlock and the send, causing a send-on-closed
+  channel panic in the caller (not recovered by `processOne`'s defer). Fix options: (a) hold
+  the lock during channel send (hot-path perf hit), (b) use a select with `r.quit` instead of
+  sending on a closed channel, (c) accept the risk and document that Submit must not be called
+  after Drain. Needs a design call — not proposing a task. (Code review 2026-09-02.)
+
+---
+
 - **#87 `lua_map` — needs embedded Lua runtime, parked.** Design/implementation decision: embedding a Lua VM in Astrate is not machine-checkable by the mule. Consider closing if Lua flow support is not on any active roadmap.
 - **#78 FDO device onboarding — milestone-4.0, investigation phase.** Too large for a single mule task; the investigation work (reading upstream's TO2 handling, inventorying endpoints/schema/keys) is a multi-session project. Parking for now until the v3.0 queue clears and this becomes the next milestone target.
 - **#1 stale — "Provide an Open Source IoT Platform unironically" (wontfix).** Has been open since the repo's founding with no activity. Consider closing.
 

Let me verify everything looks right:
→ Read .mule/todo.md [offset=110]
All done. Let me verify I didn't touch git or any frozen paths:

$ cd /root/astrate-mule && git status --short
M .mule/cron.log
 M .mule/for-giulio.md
 M .mule/todo.md
?? .mule/reviews/flow-2026-09-02.md

Only the expected files changed: my review evidence file, todo.md (4 new task lines), and for-giulio.md (1 design flag). No git operations, no frozen files touched.

**Summary of the `internal/flow` review:**

| # | Finding | Category | Task line |
|---|---------|----------|-----------|
| 1 | `Flow.setStatus` called without `f.mu` in `StartFlow` failure path (`flow.go:157`) | unguarded race | `flow-setstatus-race` |
| 2 | `Pipeline.Validate` has dead loop + redundant recomputation (`pipeline.go:117-152`) | clarity/dead code | `flow-validate-deadcode` |
| 3 | No tests for `UnmarshalJSON` error paths | missing test | `flow-unmarshal-error-tests` |
| 4 | `dataWireScalar` default case undocumented invariant | clarity | `flow-datawirescalar-fallthrough` |
| 5 | `Router.Submit` TOCTOU on `closed` flag — send-on-closed-channel panic | design decision | flagged to Giulio, not queued |
