slug: compat-note-custom-action-policy-boundary
verdict: done
at:  8cd17a8
ran: 2026-09-20T18:55:39Z on DietPi in 97s


> build · big-pickle

→ Read .mule/task.md
→ Read .mule/MULE.md
The task: propose (via `.mule/for-giulio.md`, without editing docs) a `docs/COMPATIBILITY.md` §8 sentence clarifying that trigger-policy retry/discard only apply to HTTP webhook actions, while custom forwarder actions are single-shot with only maximum_capacity/event_ttl.

Let me read the named files and verify the claim against the code.
→ Read docs/COMPATIBILITY.md
→ Read .mule/for-giulio.md
§8 is "Trigger delivery policies" (lines 212-242). Now let me verify the claim against the code — that retry/discard handlers apply only to HTTP webhook actions, while custom forwarder actions are single-shot receiving only `maximum_capacity`/`event_ttl`.

$ rg -rn "maximum_capacity|event_ttl|error_handlers|retry_times" --type go internal/ | rg -i "forward|webhook|http|policy|action" | head -50
internal/engine/triggers/actions_test.go:// admission/staleness knobs (n, n) reach custom
internal/engine/triggers/actions_test.go:		"n":[{"on":"any_error","strategy":"retry"}],
internal/engine/triggers/actions_test.go:		"n":3,
internal/engine/triggers/actions_test.go:		"n":1
internal/engine/triggers/actions_test.go:	// A Decide-loop forwarder would retry (any_error=retry, n=3):
internal/engine/triggers/actions_test.go:// 404 exactly n times then fails — the opposite of the default rule
internal/engine/triggers/actions_test.go:		"n":[{"on":"client_error","strategy":"retry"}],
internal/engine/triggers/actions_test.go:		"n":2,
internal/engine/triggers/actions_test.go:		"n":1
internal/engine/triggers/actions_test.go:		"n":[{"on":"server_error","strategy":"discard"}],
internal/engine/triggers/actions_test.go:		"n":1
internal/engine/triggers/actions_test.go:		"n":[{"on":[503],"strategy":"retry"}],
internal/engine/triggers/actions_test.go:		"n":1,
internal/engine/triggers/actions_test.go:		"n":1
internal/engine/triggers/actions_test.go:		"n":[{"on":[500],"strategy":"retry"}],
internal/engine/triggers/actions_test.go:		"n":2,
internal/engine/triggers/actions_test.go:		"n":1
internal/engine/triggers/actions_test.go:// TestPolicyTTExpiresDelivery: a policy with n=1s causes a queued
internal/engine/triggers/actions_test.go:		"n":[{"on":"any_error","strategy":"retry"}],
internal/engine/triggers/actions_test.go:		"n":1,
internal/engine/triggers/actions_test.go:		"n":2,
internal/engine/triggers/actions_test.go:		"n":1
internal/engine/triggers/actions_test.go:// TestPolicyCapacityDrops: with n=1, a second Enqueue for the
internal/engine/triggers/actions_test.go:		"n":[{"on":"any_error","strategy":"discard"}],
internal/engine/triggers/actions_test.go:		"n":1
internal/engine/triggers/actions_test.go:		"n":[{"on":"any_error","strategy":"discard"}],
internal/engine/triggers/actions_test.go:		"n":1
internal/engine/triggers/policy_test.go:			def:     `{"name":"test","n":[{"on":"any_error","strategy":"discard"}],"n":10}`,
internal/engine/triggers/policy_test.go:			def:     `{"name":"test","n":[{"on":"server_error","strategy":"retry"}],"n":10,"n":3}`,
internal/engine/triggers/policy_test.go:			name:    "valid retry with n",
internal/engine/triggers/policy_test.go:			def:     `{"name":"test","n":[{"on":"any_error","strategy":"retry"}],"n":5,"n":1,"n":60}`,
internal/engine/triggers/policy_test.go:			def:     `{"name":"k","n":[{"on":"any_error","strategy":"discard"}],"n":1}`,
internal/engine/triggers/policy_test.go:			def:     `{"name":"k","n":[{"on":"client_error","strategy":"discard"}],"n":1}`,
internal/engine/triggers/policy_test.go:			def:     `{"name":"k","n":[{"on":"server_error","strategy":"discard"}],"n":1}`,
internal/engine/triggers/policy_test.go:			def:     `{"name":"k","n":[{"on":[404,503],"strategy":"discard"}],"n":1}`,
internal/engine/triggers/policy_test.go:			def:     `{"name":"k","n":[{"on":{"custom_status_codes":[429]},"strategy":"retry"}],"n":1,"n":5}`,
internal/engine/triggers/policy_test.go:			name:    "n required when retry handler",
internal/engine/triggers/policy_test.go:			def:     `{"name":"k","n":[{"on":"any_error","strategy":"retry"}],"n":1}`,
internal/engine/triggers/policy_test.go:			wantErr: "n must be 1-100",
internal/engine/triggers/policy_test.go:			name:    "n forbidden when no retry handler",
internal/engine/triggers/policy_test.go:			def:     `{"name":"k","n":[{"on":"any_error","strategy":"discard"}],"n":1,"n":3}`,
internal/engine/triggers/policy_test.go:			wantErr: "n requires a retry handler",
internal/engine/triggers/policy_test.go:			name:    "n out of range low",
internal/engine/triggers/policy_test.go:			def:     `{"name":"k","n":[{"on":"any_error","strategy":"retry"}],"n":1,"n":0}`,
internal/engine/triggers/policy_test.go:			wantErr: "n must be 1-100",
internal/engine/triggers/policy_test.go:			name:    "n out of range high",
internal/engine/triggers/policy_test.go:			def:     `{"name":"k","n":[{"on":"any_error","strategy":"retry"}],"n":1,"n":101}`,
internal/engine/triggers/policy_test.go:			wantErr: "n must be 1-100",
internal/engine/triggers/policy_test.go:			name:    "n too low",
internal/engine/triggers/policy_test.go:			def:     `{"name":"k","n":[{"on":"any_error","strategy":"discard"}],"n":0}`,

$ rg -l "error_handlers|retry_times|maximum_capacity|event_ttl" --type go internal/ pkg/ 2>/dev/null
internal/engine/triggers/actions.go
internal/engine/triggers/policy.go
internal/realm/dashboard_compat_test.go
internal/realm/async_operation_test.go
internal/realm/http_test.go
internal/realm/policy_test.go
internal/store/policies_test.go
internal/engine/cache_test.go
internal/engine/triggers/policy_test.go
internal/engine/triggers/actions_test.go

→ Read internal/engine/triggers/actions.go
→ Read internal/engine/triggers/policy.go
The claim checks out against the code:

- `Action.Custom` routes to `forward()` (actions.go:435-439); only `webhook()` consults `policy.Decide` (actions.go:524)
- `forward()` never applies retry/discard — the Forwarder seam returns only an error (actions.go:161-167, 471-479)
- `maximum_capacity` is enforced at `Enqueue` on every delivery (actions.go:338-349) and `event_ttl` at dequeue before the custom/webhook branch (actions.go:400-416)

Now the deliverable: a for-giulio entry (top of file, newest first) proposing the §8 sentence, without touching `docs/COMPATIBILITY.md`.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **COMPATIBILITY.md §8 wording proposal: trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (Forwarder) actions are single-shot by design, receiving only `maximum_capacity`/`event_ttl`.** Propose adding to the "Trigger delivery policies" deviation: *"Policy `error_handlers` (retry/discard) and the `retry_times` attempt cap govern HTTP webhook actions only: a custom action delivered through the Forwarder seam is single-shot by design, receiving from the policy solely its `maximum_capacity` in-flight bound and its `event_ttl` staleness drop, with a forward failure counted once and never retried."* Verified in `internal/engine/triggers/actions.go`: `Action.Custom` routes to `forward` (actions.go:435-439), which never consults `policy.Decide` — the Forwarder seam returns only an error, no HTTP status (actions.go:161-167, doc comment at actions.go:471-479) — while `Enqueue` bounds `maximum_capacity` on every delivery (actions.go:338-349) and `deliver` drops past `event_ttl` before the custom/webhook branch (actions.go:400-416). The existing warn-on-single-shot line already tells the operator when a policy's `retry_times` cannot apply to the forward path (actions.go:491-497).
+
+---
+
 - **github-issues triage run, 2026-09-20: still nothing proposable — the mule-alarm pile is now 15 days (#94–#107).** Seventeen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-18 run: **#107**, yesterday's alarm (created 09-19 10:49Z, "nothing has landed in 14h"); **#106** has expired into the pile (superseded by #107). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-18 one: close #94–#106** — each was a one-day event, superseded by next day's alarm, never actionable; leave #107 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
 
 ---
 

→ Read .mule/for-giulio.md [limit=30]
Done: proposed the §8 boundary sentence in a new top-of-queue for-giulio entry, backed by the code, without touching `docs/COMPATIBILITY.md` or git.

```
Done: Proposed docs/COMPATIBILITY.md §8 sentence (trigger-policy retry/discard govern HTTP webhooks only; custom Forwarder actions single-shot with only maximum_capacity/event_ttl) via a new .mule/for-giulio.md entry
Files: .mule/for-giulio.md
Verified: claim checked against internal/engine/triggers/actions.go (deliver branches custom->forward, forward never calls policy.Decide, Enqueue/deliver apply maximum_capacity/event_ttl on the shared path) -> pass; no Go/docs edits, so no gate run
Unsure: nothing
Follow-ups: none
```
