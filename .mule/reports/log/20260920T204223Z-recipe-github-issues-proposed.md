slug: recipe-github-issues
verdict: proposed
at:  2e34867
ran: 2026-09-20T20:42:23Z on DietPi in 186s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md

$ wc -l .mule/todo.md; ls .mule/for-giulio.md 2>/dev/null
202 .mule/todo.md
.mule/for-giulio.md

$ gh issue list --state open --limit 40 \
  --json number,title,labels,updatedAt,comments \
  --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
#107 mule: nothing has landed in 14h [mule-alarm ] []c
#106 mule: nothing has landed in 14h [mule-alarm ] []c
#105 mule: nothing has landed in 14h [mule-alarm ] []c
#104 mule: nothing has landed in 14h [mule-alarm ] []c
#103 mule: nothing has landed in 16h [mule-alarm ] []c
#102 mule: nothing has landed in 14h [mule-alarm ] []c
#101 mule: nothing has landed in 14h [mule-alarm ] []c
#100 mule: nothing has landed in 15h [mule-alarm ] []c
#99 mule: nothing has landed in 14h [mule-alarm ] []c
#98 mule: nothing has landed in 16h [mule-alarm ] []c
#97 mule: nothing has landed in 15h [mule-alarm ] []c
#96 mule: nothing has landed in 14h [mule-alarm ] []c
#95 mule: nothing has landed in 15h [mule-alarm ] []c
#94 mule: nothing has landed in 15h [mule-alarm ] []c
#93 In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged [mule-review upstream-parity ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:The mule pushed `8c61268` to `mule/queue` for this.

Unreviewed: the gates passed — it compiles, the tests pass, a new test was shown to fail without the change, and no frozen file was touched. That is not the same as the change being worth having. Nothing merges to `main` until someone reads it:

    bash tools/mule.sh review

Left open on purpose. Whether this actually resolves the issue is your call, not mine. createdAt:2026-09-04T19:54:58Z id:IC_kwDORmfsJs8AAAABSo2DbA includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/93#issuecomment-5545755500 viewerDidAuthor:true]]c
#92 keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds [enhancement upstream-parity upstream-experimental ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:Parked by Giulio on 2026-09-04: **wait for a stable \`v1.4.0\`.** The spec is complete but ships only in \`v1.4.0-rc.5\`, and building against an rc risks rebuilding it.

So this issue is not waiting on anyone to remember it. It has a row in [`.mule/waiting-on.md`](https://github.com/astrate-platform/astrate/blob/main/.mule/waiting-on.md), and the weekly `mule-upstream-watch` job now reads that file as its first step every run: when a stable `v1.4.0` tag exists upstream, it escalates to `.mule/for-giulio.md` naming this issue. An `-rc.N` does not satisfy the row.

That mechanism was added because it was missing here: #51 was parked on the same condition, the condition fired on 2026-08-31, and nothing connected the new upstream tags to the parked issue — it took a hand triage four days later to catch it. createdAt:2026-09-04T19:52:21Z id:IC_kwDORmfsJs8AAAABSo0f8A includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/92#issuecomment-5545730032 viewerDidAuthor:true]]c
#78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate) [enhancement milestone-4.0 upstream-parity ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:Reframed 2026-08-23 by Giulio's decision: no longer parked — **strategic feature** (zero-touch factory onboarding is a commercial-viability requirement for an IoT platform), milestone-4.0 candidate.

## Verified facts (2026-08-23, upstream master)

- Upstream integrates FDO **inside the Pairing service**: `fdo_onboarding_controller.ex`, session plugs (`setup_fdo`, `fdo_session`), CBOR codec + voucher queries in `astarte_data_access`, and two dedicated libs `astarte_fdo` / `astarte_fdo_core` (SECO Mind, 2025; bugfixes through 2026-05). Not an optional sidecar.
- Device side is the official SDK `astarte-device-fdo-rust` (crate on docs.rs, pushed 2026-08-21) — actively maintained.
- Upstream does NOT implement manufacturing / rendezvous services either; those come from the FIDO Alliance reference stack (`fdo-rs/fido-device-onboard-rs`, used by Red Hat/Fedora IoT).

## Scope decision

Astrate implements **only the last mile**, mirroring upstream's cut:
- owner-side protocol surface in our pairing service (TO1 redirect + TO2 onboarding, CBOR wire format, ownership-voucher storage, owner key management)
- reuse the existing open-source FDO ecosystem for manufacturing/rendezvous — never reimplement it.

**Acceptance:** a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end and lands as a provisioned realm device.

**Docs are a first-class deliverable:** upstream ships almost no public FDO documentation (code-level only); ours must exceed that — operator guide covering the full chain (manufacturing → rendezvous → owner) with Astrate-specific setup.

## Next step

Investigation phase before any implementation: read upstream's TO2 handling end-to-end (`libs/astarte_fdo*", pairing controllers/plugs, `fdo/queries`), inventory the exact endpoints, credential/voucher schema and key material required, and measure what our pairing service lacks. Output feeds the v4.0 scope decision in `.mule/milestones.md`. createdAt:2026-08-23T16:16:08Z id:IC_kwDORmfsJs8AAAABQRcugw includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/78#issuecomment-5387005571 viewerDidAuthor:true]]c
#1 Provide an Open Source IoT Platform unironically, unironically [wontfix ] []c

→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ git log --oneline -5 -- .mule/for-giulio.md .mule/todo.md 2>/dev/null; echo "---"; git log --oneline -3 2>/dev/null
2e34867 mule: refill the queue
395d8bb mule: log docs-sync-native-socket-missing-403-500
ba8afaf mule: log docs-sync-native-metrics-example-fake-series
d2a8c51 mule: refill the queue
ea4f1b9 mule: log compat-note-custom-action-policy-boundary
---
2e34867 mule: refill the queue
fc5f492 mule: recipe milestones ran (193s)
9195e7f mule: recipe hygiene timed out

→ Read .mule/todo.md [offset=174]

$ git log --oneline -10 -- .mule/todo.md
395d8bb mule: log docs-sync-native-socket-missing-403-500
ba8afaf mule: log docs-sync-native-metrics-example-fake-series
d2a8c51 mule: refill the queue
ea4f1b9 mule: log compat-note-custom-action-policy-boundary
d68bc52 mule: log triggers-custom-action-policy-nodecide
3c8b89a mule: triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
7e2ceeb mule: log errorname-missing-required-test
e645b75 mule: refill the queue
054d305 mule: blocked channels-rejoin-authz-mismatch
ebc3854 mule: log channels-rejoin-joinref-tagging-test

$ git log --oneline -20 -- .mule/for-giulio.md
2e34867 mule: refill the queue
8cd17a8 mule: compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).
1dba3d5 mule: refill the queue
74021c8 mule: refill the queue
79439e6 mule: refill the queue
f6dc460 mule: refill the queue
bcd4383 mule: compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.
ff7905b mule: refill the queue
dbbf2e5 mule: refill the queue
fa326e1 mule: refill the queue
5a314c7 mule: refill the queue
ed715ec mule: refill the queue
a7fb528 mule: refill the queue
b36f50a mule: refill the queue
c8218e5 mule: refill the queue
bb70e4c mule: refill the queue
59940df mule: refill the queue
1c8d6f4 mule: refill the queue
a3b5716 mule: refill the queue
514a685 mule: refill the queue

$ git log --oneline -15; echo "---"; git log -1 --format='%ci' HEAD
2e34867 mule: refill the queue
fc5f492 mule: recipe milestones ran (193s)
9195e7f mule: recipe hygiene timed out
395d8bb mule: log docs-sync-native-socket-missing-403-500
f9d8b27 mule: docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
ba8afaf mule: log docs-sync-native-metrics-example-fake-series
43c4a8b mule: docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate_devices_total{realm="test"} 42`, but no such series is ever registered — the registry carries `astrate_broker_sessions` (internal/observability/metrics.go:54-62), `astrate_db_pool_{acquired,idle,total,max}_conns` (metrics.go:65-79), plus the engine/flow/trigger families (`astrate_engine_*` engine/router.go:411-444, `astrate_engine_stream_dropped_total` engine/stream/bus.go:119, `astrate_engine_trigger_*` engine/triggers/actions.go:296-300, `astrate_flow_router_*` flow/router.go:260-276) and the standard `go_*`/`process_*` runtime collectors — so the example advertises a scrape result the wire can never return. Swap it for a real series (e.g. `astrate_broker_sessions 0`), keeping the HELP/TYPE comment style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
d2a8c51 mule: refill the queue
14f73da mule: recipe docs-sync ran (218s)
ea4f1b9 mule: log compat-note-custom-action-policy-boundary
8cd17a8 mule: compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).
d68bc52 mule: log triggers-custom-action-policy-nodecide
3c8b89a mule: triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
7e2ceeb mule: log errorname-missing-required-test
739df67 mule: errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
---
2026-09-20 22:30:31 +0200

$ git status --short; echo "---untracked---"; git diff --stat
---untracked---

$ git show --stat HEAD
commit 2e348674b3cfbd77ea1fe3ba42ee8b1e4ddd619b
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Sun Sep 20 22:30:31 2026 +0200

    mule: refill the queue

 .mule/for-giulio.md | 4 ++++
 1 file changed, 4 insertions(+)

$ git show HEAD
commit 2e348674b3cfbd77ea1fe3ba42ee8b1e4ddd619b
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Sun Sep 20 22:30:31 2026 +0200

    mule: refill the queue

diff --git a/.mule/for-giulio.md b/.mule/for-giulio.md
index c83684e..fb54a47 100644
--- a/.mule/for-giulio.md
+++ b/.mule/for-giulio.md
@@ -10,6 +10,10 @@ line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+
+---
+
 - **COMPATIBILITY.md §8 wording proposal: trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (Forwarder) actions are single-shot by design, receiving only `maximum_capacity`/`event_ttl`.** Propose adding to the "Trigger delivery policies" deviation: *"Policy `error_handlers` (retry/discard) and the `retry_times` attempt cap govern HTTP webhook actions only: a custom action delivered through the Forwarder seam is single-shot by design, receiving from the policy solely its `maximum_capacity` in-flight bound and its `event_ttl` staleness drop, with a forward failure counted once and never retried."* Verified in `internal/engine/triggers/actions.go`: `Action.Custom` routes to `forward` (actions.go:435-439), which never consults `policy.Decide` — the Forwarder seam returns only an error, no HTTP status (actions.go:161-167, doc comment at actions.go:471-479) — while `Enqueue` bounds `maximum_capacity` on every delivery (actions.go:338-349) and `deliver` drops past `event_ttl` before the custom/webhook branch (actions.go:400-416). The existing warn-on-single-shot line already tells the operator when a policy's `retry_times` cannot apply to the forward path (actions.go:491-497).
 
 ---

$ git show fc5f492 --stat; echo "==="; git show 9195e7f --stat
commit fc5f4927a07332af5c6123702d65ca91d25a0228
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Sun Sep 20 22:30:30 2026 +0200

    mule: recipe milestones ran (193s)

 .../20260920T203030Z-recipe-milestones-proposed.md | 1196 ++++++++++++++++++++
 1 file changed, 1196 insertions(+)
===
commit 9195e7faca366be3797a471fe648a782132d658a
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Sun Sep 20 21:54:49 2026 +0200

    mule: recipe hygiene timed out

 .../log/20260920T195443Z-recipe-hygiene-timeout.md | 157 +++++++++++++++++++++
 1 file changed, 157 insertions(+)

$ git log --oneline -S "github-issues" --all -- .mule/for-giulio.md | head -10
1dba3d5 mule: refill the queue
ff7905b mule: refill the queue
5a314c7 mule: refill the queue
ed715ec mule: refill the queue
b36f50a mule: refill the queue
bb70e4c mule: refill the queue
59940df mule: refill the queue
514a685 mule: refill the queue
9a28cd1 mule: refill the queue
08dee34 mule: refill the queue

$ git log --oneline -S "github-issues triage run, 2026-09-20" -- .mule/for-giulio.md | head -5
1dba3d5 mule: refill the queue

$ git show 1dba3d5 --format='%ci %s' --stat | head -20
2026-09-20 13:10:57 +0200 mule: refill the queue

 .mule/for-giulio.md | 4 ++++
 1 file changed, 4 insertions(+)

Nothing new in the working tree and the issue set is unchanged. Last github-issues run was this morning (13:10, commit `1dba3d5`): same 17/18 open issues, same triage result. Let me do the current survey's analysis and check the for-giulio top entry.

Survey result (18 open): #94–#107 are the daily `mule-alarm` pile; #93 aclhook comment (mule-review, `8c61268` pushed); #92 keyAgreement (parked on stable v1.4.0, waiting-on.md mechanism); #78 FDO (milestone-4.0, already escalated); #1 untouched per standing instruction. No machine-checkable code-fix candidates → no todo.md task lines.

I'll append one for-giulio line reflecting this run, keeping the alarm-closure proposal from this morning's run without duplicating it.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **github-issues triage re-run, 2026-09-20 (evening): nothing proposable, no change since this morning's 13:10 run.** Re-surveyed with the recipe command: the issue set is identical to the 2026-09-20 morning triage — **#94–#107** are still the daily `mule-alarm` noise (14 issues, zero comments, telemetry not code), **#93** aclhook comment rewrite is still mule-review with `8c61268` already pushed (its own recipe path), **#92** keyAgreement is still parked on a stable upstream v1.4.0 with its waiting-on.md row and the mule-upstream-watch escalation wired (an -rc does not satisfy it), **#78** FDO is still the milestone-4.0 design/investigation already escalated, and **#1** is untouched per standing instruction. No new issue has been filed since the morning run. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines proposed**. This morning's proposal stands unchanged: **close #94–#106** — each was a one-day event, superseded by the next day's alarm, never actionable — and leave **#107** (still the newest, no #108 yet) to self-expire. The queue idle streak since ~2026-09-04/05 persists; if that is not intentional the wedge is on the Pi (queue review = the fix, per the 13:10 line).
+
+---
+
 - **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
 
 ---
 

→ Read .mule/for-giulio.md [limit=40]
Recipe complete. Survey ran, nothing new.

Done: github-issues triage re-run (evening 2026-09-20) — nothing proposable, issue set unchanged since this morning's 13:10 run; appended one escalation line, no task lines.
Files: .mule/for-giulio.md
Verified: gh issue list --state open --limit 40 -> pass (18 open: #94–#107 alarms, #93 mule-review, #92 parked, #78 FDO, #1 untouched)
Unsure: nothing
Follow-ups: none — no machine-checkable candidates; close-#94–#106 proposal already filed this morning and re-affirmed. No git touched.
