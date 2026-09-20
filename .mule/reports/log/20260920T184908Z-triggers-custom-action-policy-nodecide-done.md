slug: triggers-custom-action-policy-nodecide
verdict: done
at:  3c8b89a
ran: 2026-09-20T18:49:08Z on DietPi in 459s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/engine/triggers/actions.go
✱ Grep "func \(.*\) Decide|StrategyRetry|StrategyDiscard|StatusTransport" in internal/engine/triggers · 44 matches
→ Read internal/engine/triggers/actions_test.go
→ Read internal/engine/triggers/policy.go
→ Read internal/engine/forward
✱ Grep "Forwarder|forward|Forward" in internal/config · 51 matches
→ Read internal/engine/forward/http.go
→ Read internal/engine/forward/nats.go
✱ Grep "forward|Forward|policy|Policy|retry|Retry|trigger|Trigger" in docs/DESIGN.md · 100 matches
→ Read docs/COMPATIBILITY.md [offset=208, limit=45]
✱ Grep "1\.1|forward|Forward" in docs/DESIGN.md · 68 matches
✱ Grep "func \(.*\) Policy\(\)|func \(.*\) AttachPolicy|policy" in internal/engine/triggers · 98 matches
→ Read .mule/todo.md
Now I understand the task. Let me look at the `Trigger.Policy()` helper and `match.go` context to confirm the policy accessor, then decide the intent.

The decision: the `Forwarder` seam returns only `error` (no HTTP status), while `Policy.Decide` is status-shaped (handlers match `client_error`/`server_error`/`any_error`/explicit codes). Routing `forward` through the Decide loop would force every failure into `StatusTransport` — matching only keyword handlers, silently leaving explicit status-list policies even more partially applied, and labelling every audit reason "transport failure". The task's own phrasing ("pin single-shot **as designed**") plus MULE's "never invent a behaviour" settle it: **pin single-shot as designed**, keep capacity/TTL applying (admission/staleness are action-kind-agnostic), and remove the *silence* — a failure on a trigger whose policy contains a retry rule must say so in the log.
→ Read internal/engine/triggers/match.go [offset=150, limit=25]
Now the changes. First, `forward()` in actions.go — declare the single-shot design and make the boundary non-silent:

← Edit internal/engine/triggers/actions.go
Index: /root/astrate-mule/internal/engine/triggers/actions.go
===================================================================
--- /root/astrate-mule/internal/engine/triggers/actions.go
+++ /root/astrate-mule/internal/engine/triggers/actions.go
@@ -464,9 +464,17 @@
 	}
 	return []byte(rendered), nil
 }
 
-// forward hands a non-HTTP action to the Forwarder extension point.
+// forward hands a non-HTTP action to the Forwarder extension point. It is
+// single-shot by design: the Forwarder seam returns only an error, never an
+// HTTP status, so policy.Decide's retry/discard handlers (which decide on a
+// status code, and are what webhook deliveries retry on) cannot apply to
+// custom actions. The admission/staleness policy knobs do still apply to the
+// seam: Enqueue bounds maximum_capacity in-flight and deliver drops past
+// event_ttl, exactly as for webhooks. A forwarder failure counts as one
+// failed attempt, and when the trigger's policy carries a retry rule the log
+// says the retry was not applied instead of abandoning it silently.
 func (x *Executor) forward(d Delivery, body []byte) {
 	if x.cfg.Forwarder == nil {
 		x.outcomes.WithLabelValues(outcomeSkipped).Inc()
 		x.log.Info("no forwarder configured; custom trigger action skipped",
@@ -476,8 +484,16 @@
 	ctx, cancel := context.WithTimeout(context.Background(), x.cfg.RequestTimeout)
 	defer cancel()
 	if err := x.cfg.Forwarder.Forward(ctx, d.Realm, d.Trigger.Name, d.Trigger.Action.Custom, body); err != nil {
 		x.outcomes.WithLabelValues(outcomeFailed).Inc()
+		if pol := d.Trigger.Policy(); pol != nil && pol.RetryTimes > 0 {
+			// RetryTimes is nonzero iff some handler retries (CompilePolicy
+			// rejects retry_times without a retry handler): the operator
+			// expects a retry this path cannot give.
+			x.log.Warn("custom action delivery failed; forward path is single-shot and does not apply policy retry",
+				"realm", d.Realm, "trigger", d.Trigger.Name, "policy", pol.Name, "err", err)
+			return
+		}
 		x.log.Warn("custom trigger action failed",
 			"realm", d.Realm, "trigger", d.Trigger.Name, "err", err)
 		return
 	}

← Edit internal/engine/triggers/actions.go
Index: /root/astrate-mule/internal/engine/triggers/actions.go
===================================================================
--- /root/astrate-mule/internal/engine/triggers/actions.go
+++ /root/astrate-mule/internal/engine/triggers/actions.go
@@ -158,9 +158,12 @@
 // internal/config and cmd/astrate. The default (nil) logs the event and
 // counts it as skipped, which is the designed behaviour, not a gap in this
 // code path.
 type Forwarder interface {
-	// Forward delivers one rendered event for a custom action.
+	// Forward delivers one rendered event for a custom action. It is a
+	// single-shot hand-off: Astrate never retries a returned error, and a
+	// trigger policy's retry/discard handlers do not apply to this seam
+	// (only maximum_capacity and event_ttl do — see forward).
 	Forward(ctx context.Context, realm, trigger string, action json.RawMessage, event []byte) error
 }
 
 // Delivery is one matched (trigger, event) pair queued for execution.

Now the test. Add the `bytes` import and the discriminating test:
← Edit internal/engine/triggers/actions_test.go
Index: /root/astrate-mule/internal/engine/triggers/actions_test.go
===================================================================
--- /root/astrate-mule/internal/engine/triggers/actions_test.go
+++ /root/astrate-mule/internal/engine/triggers/actions_test.go
@@ -1,7 +1,8 @@
 package triggers
 
 import (
+	"bytes"
 	"context"
 	"encoding/json"
 	"errors"
 	"io"

← Edit internal/engine/triggers/actions_test.go
Index: /root/astrate-mule/internal/engine/triggers/actions_test.go
===================================================================
--- /root/astrate-mule/internal/engine/triggers/actions_test.go
+++ /root/astrate-mule/internal/engine/triggers/actions_test.go
@@ -612,8 +612,72 @@
 	x.Enqueue(testDelivery(&Action{Custom: json.RawMessage(`{"amqp_exchange":"e"}`)}))
 	eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
 }
 
+// TestForwarderPolicySingleShot pins the designed boundary between delivery
+// policies and the custom-action (forward) path: the Forwarder seam returns
+// only an error, never an HTTP status, so policy.Decide's retry/discard
+// handlers cannot apply to it the way they do to webhook deliveries. Only the
+// admission/staleness knobs (maximum_capacity, event_ttl) reach custom
+// actions. A retry policy therefore still drives exactly one forward attempt,
+// counted as failed; routing forward through the Decide loop would break this
+// test by retrying. The failure is not silent: when the policy carries a
+// retry rule, the log states that the forward path is single-shot.
+func TestForwarderPolicySingleShot(t *testing.T) {
+	policy := mustCompilePolicy(t, `{
+		"name":"fwd-retry",
+		"error_handlers":[{"on":"any_error","strategy":"retry"}],
+		"retry_times":3,
+		"maximum_capacity":1
+	}`)
+
+	fwd := &recordingForwarder{err: io.ErrUnexpectedEOF}
+	var logBuf bytes.Buffer
+	logger := slog.New(slog.NewTextHandler(&logBuf, nil))
+	x := NewExecutor(ExecutorConfig{
+		Workers:      1,
+		BackoffStart: time.Millisecond,
+		BackoffCap:   5 * time.Millisecond,
+		Registerer:   prometheus.NewRegistry(),
+		Forwarder:    fwd,
+		Logger:       logger,
+	})
+	t.Cleanup(func() {
+		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+		defer cancel()
+		if err := x.Close(ctx); err != nil {
+			t.Errorf("executor close: %v", err)
+		}
+	})
+
+	d := testDelivery(&Action{Custom: json.RawMessage(`{"nats_subject":"events"}`)})
+	d.Trigger.AttachPolicy(policy)
+	x.Enqueue(d)
+
+	eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
+	if got := outcome(x, outcomeForwarded); got != 0 {
+		t.Errorf("forwarded = %v, want 0", got)
+	}
+	// A Decide-loop forwarder would retry (any_error=retry, retry_times=3):
+	// watch the attempt count across a window so any erroneous retry is seen.
+	deadline := time.Now().Add(100 * time.Millisecond)
+	for time.Now().Before(deadline) {
+		fwd.mu.Lock()
+		called := fwd.called
+		fwd.mu.Unlock()
+		if called != 1 {
+			t.Fatalf("forwarder called %d times, want 1 (forward path is single-shot)", called)
+		}
+		time.Sleep(time.Millisecond)
+	}
+	if got := promtest.ToFloat64(x.retries); got != 0 {
+		t.Errorf("retries = %v, want 0", got)
+	}
+	if !strings.Contains(logBuf.String(), "does not apply policy retry") {
+		t.Errorf("log %q missing the single-shot retry notice", logBuf.String())
+	}
+}
+
 // TestEnqueueAfterCloseDrops: enqueueing after Close drops with a metric and
 // never panics on the closed channel.
 func TestEnqueueAfterCloseDrops(t *testing.T) {
 	x := newTestExecutor(t, ExecutorConfig{Workers: 1})

Now append the follow-up task line to `.mule/todo.md` — the operator-facing policy boundary deserves a COMPATIBILITY §8 sentence, and that file is never-touch, so it goes through `.mule/for-giulio.md`:
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -1,6 +1,7 @@
 - [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
 - [ ] triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
+- [ ] compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).
 
 # Mule queue
 
 One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).

Now verify. Run the package tests, vet, and gofmt:
$ gofmt -l .
(no output)

$ go test ./internal/engine/triggers/... 2>&1
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	2.030s

The gate says my test must fail without my change. Let me prove it — temporarily revert just the log line, run only the new test, then restore:
← Edit internal/engine/triggers/actions.go
Index: /root/astrate-mule/internal/engine/triggers/actions.go
===================================================================
--- /root/astrate-mule/internal/engine/triggers/actions.go
+++ /root/astrate-mule/internal/engine/triggers/actions.go
@@ -491,9 +491,9 @@
 if pol := d.Trigger.Policy(); pol != nil && pol.RetryTimes > 0 {
 	// RetryTimes is nonzero iff some handler retries (CompilePolicy
 	// rejects retry_times without a retry handler): the operator
 	// expects a retry this path cannot give.
-	x.log.Warn("custom action delivery failed; forward path is single-shot and does not apply policy retry",
+	x.log.Warn("custom action delivery failed (no retry)",
 		"realm", d.Realm, "trigger", d.Trigger.Name, "policy", pol.Name, "err", err)
 	return
 }
 x.log.Warn("custom trigger action failed",

$ go test ./internal/engine/triggers/ -run 'TestForwarderPolicySingleShot' -count=1 2>&1; echo "exit=$?"
--- FAIL: TestForwarderPolicySingleShot (0.10s)
    actions_test.go:676: log "time=2026-09-20T20:44:39.423+02:00 level=WARN msg=\"custom action delivery failed (no retry)\" component=triggers realm=testrealm trigger=hook policy=fwd-retry err=\"unexpected EOF\"\n" missing the single-shot retry notice
FAIL
FAIL	github.com/astrate-platform/astrate/internal/engine/triggers	0.250s
FAIL
exit=1

Test fails without my change, as required. Restore the log:
← Edit internal/engine/triggers/actions.go
Index: /root/astrate-mule/internal/engine/triggers/actions.go
===================================================================
--- /root/astrate-mule/internal/engine/triggers/actions.go
+++ /root/astrate-mule/internal/engine/triggers/actions.go
@@ -491,9 +491,9 @@
 if pol := d.Trigger.Policy(); pol != nil && pol.RetryTimes > 0 {
 	// RetryTimes is nonzero iff some handler retries (CompilePolicy
 	// rejects retry_times without a retry handler): the operator
 	// expects a retry this path cannot give.
-	x.log.Warn("custom action delivery failed (no retry)",
+	x.log.Warn("custom action delivery failed; forward path is single-shot and does not apply policy retry",
 		"realm", d.Realm, "trigger", d.Trigger.Name, "policy", pol.Name, "err", err)
 	return
 }
 x.log.Warn("custom trigger action failed",

Now run the full gates:
$ go test ./internal/engine/triggers/... -count=1 2>&1
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	2.836s

$ go vet ./... 2>&1; echo "vet_exit=$?"
vet_exit=0

$ go test ./... 2>&1 | tail -40; echo "test_exit=${PIPESTATUS:-$?}"
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	2.194s
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
test_exit=0

All green. Let me try the race build on the triggers package (and note which I ran):
$ go test -race ./internal/engine/triggers/ -count=1 2>&1; echo "race_exit=$?"
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/engine/triggers	0.024s
FAIL
race_exit=1

Race detector is unsupported on this box (39-bit VMA). Final self-review of the working tree (read-only):
$ git status --short; echo "---"; git --no-pager diff --stat
M .mule/todo.md
 M internal/engine/triggers/actions.go
 M internal/engine/triggers/actions_test.go
---
 .mule/todo.md                            |  1 +
 internal/engine/triggers/actions.go      | 23 ++++++++++-
 internal/engine/triggers/actions_test.go | 65 ++++++++++++++++++++++++++++++++
 3 files changed, 87 insertions(+), 2 deletions(-)

$ git --no-pager diff internal/engine/triggers/actions.go
diff --git a/internal/engine/triggers/actions.go b/internal/engine/triggers/actions.go
index 9c6bb88..536e4e1 100644
--- a/internal/engine/triggers/actions.go
+++ b/internal/engine/triggers/actions.go
@@ -159,7 +159,10 @@ func parseAction(raw json.RawMessage) (*Action, []string, error) {
 // counts it as skipped, which is the designed behaviour, not a gap in this
 // code path.
 type Forwarder interface {
-	// Forward delivers one rendered event for a custom action.
+	// Forward delivers one rendered event for a custom action. It is a
+	// single-shot hand-off: Astrate never retries a returned error, and a
+	// trigger policy's retry/discard handlers do not apply to this seam
+	// (only maximum_capacity and event_ttl do — see forward).
 	Forward(ctx context.Context, realm, trigger string, action json.RawMessage, event []byte) error
 }
 
@@ -465,7 +468,15 @@ func renderMustache(template, realm string, event SimpleEvent) ([]byte, error) {
 	return []byte(rendered), nil
 }
 
-// forward hands a non-HTTP action to the Forwarder extension point.
+// forward hands a non-HTTP action to the Forwarder extension point. It is
+// single-shot by design: the Forwarder seam returns only an error, never an
+// HTTP status, so policy.Decide's retry/discard handlers (which decide on a
+// status code, and are what webhook deliveries retry on) cannot apply to
+// custom actions. The admission/staleness policy knobs do still apply to the
+// seam: Enqueue bounds maximum_capacity in-flight and deliver drops past
+// event_ttl, exactly as for webhooks. A forwarder failure counts as one
+// failed attempt, and when the trigger's policy carries a retry rule the log
+// says the retry was not applied instead of abandoning it silently.
 func (x *Executor) forward(d Delivery, body []byte) {
 	if x.cfg.Forwarder == nil {
 		x.outcomes.WithLabelValues(outcomeSkipped).Inc()
@@ -477,6 +488,14 @@ func (x *Executor) forward(d Delivery, body []byte) {
 	defer cancel()
 	if err := x.cfg.Forwarder.Forward(ctx, d.Realm, d.Trigger.Name, d.Trigger.Action.Custom, body); err != nil {
 		x.outcomes.WithLabelValues(outcomeFailed).Inc()
+		if pol := d.Trigger.Policy(); pol != nil && pol.RetryTimes > 0 {
+			// RetryTimes is nonzero iff some handler retries (CompilePolicy
+			// rejects retry_times without a retry handler): the operator
+			// expects a retry this path cannot give.
+			x.log.Warn("custom action delivery failed; forward path is single-shot and does not apply policy retry",
+				"realm", d.Realm, "trigger", d.Trigger.Name, "policy", pol.Name, "err", err)
+			return
+		}
 		x.log.Warn("custom trigger action failed",
 			"realm", d.Realm, "trigger", d.Trigger.Name, "err", err)
 		return

Done. I pinned custom-action forwarding as **single-shot by design** and removed the silence.

**Intent reasoning:** `Forwarder.Forward` returns only `error` (no HTTP status, actions.go:161-164), while `Policy.Decide` is status-shaped (actions.go:505-579). Routing `forward` through the Decide loop would force every failure into `StatusTransport`, which only matches `any_error`/`server_error` keyword handlers and never explicit status lists — leaving status-list policies *more* silently partial, and labeling every audit reason "transport failure" even for a bus 422/NATS marshall error. Since the task's phrasing offered "single-shot as designed", and enforcing two of four knobs on the seam is coherent (capacity/TTL are admission/staleness, action-kind-agnostic; retry/discard are webhook-status-shaped), I kept single-shot and made the boundary explicit and non-silent.

```
Done: Pinned custom-action forwards as single-shot by design; failed retry-policy forwards now log that the policy retry rule is not applied, and a test forces the choice
Files: internal/engine/triggers/actions.go, internal/engine/triggers/actions_test.go, .mule/todo.md
Verified: go test ./internal/engine/triggers/... -> pass; go test ./... -> pass; go vet ./... -> pass; gofmt -l . -> clean; new test temporarily reverted against -> fails (gate proof). go test -race not run: ThreadSanitizer unsupported on this box (39-bit VMA, "FATAL: ThreadSanitizer: unsupported VMA range"), so the test's concurrency (polling the forwarder call counter across a window) is unverified for races — reads are mutex-guarded via the existing recordingForwarder lock, and no new shared state was added
Unsure: The choice itself — Option A (route forward through the Decide loop) is the alternative; I rejected it because the seam returns no status code, so Decide would see only StatusTransport and silence/worse-label the status-case. Confirmed a retry policy still gets capacity/TTL. docs/COMPATIBILITY.md §8 still reads as if retry/discard apply to all actions; that file is never-touch, so the follow-up below proposes the wording instead
Follow-ups: appended "compat-note-custom-action-policy-boundary" to .mule/todo.md (propose a §8 sentence for the boundary via .mule/for-giulio.md)
```
