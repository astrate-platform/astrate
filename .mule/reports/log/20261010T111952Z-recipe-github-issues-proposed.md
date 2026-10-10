slug: recipe-github-issues
verdict: proposed
at:  4d78c9b
ran: 2026-10-10T11:19:52Z on DietPi in 161s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md
→ Read .mule/todo.md

$ gh issue list --state open --limit 40 \
  --json number,title,labels,updatedAt,comments \
  --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
#123 mule: nothing has landed in 14h [mule-alarm ] []c
#122 mule: nothing has landed in 15h [mule-alarm ] []c
#121 mule: nothing has landed in 16h [mule-alarm ] []c
#120 mule: nothing has landed in 14h [mule-alarm ] []c
#119 mule: nothing has landed in 16h [mule-alarm ] []c
#118 mule: nothing has landed in 14h [mule-alarm ] []c
#117 mule: nothing has landed in 15h [mule-alarm ] []c
#116 mule: nothing has landed in 14h [mule-alarm ] []c
#115 mule: nothing has landed in 14h [mule-alarm ] []c
#114 mule: nothing has landed in 8h [mule-alarm ] []c
#113 mule: nothing has landed in 14h [mule-alarm ] []c
#112 mule: nothing has landed in 14h [mule-alarm ] []c
#111 mule: nothing has landed in 16h [mule-alarm ] []c
#110 mule: nothing has landed in 14h [mule-alarm ] []c
#109 mule: nothing has landed in 14h [mule-alarm ] []c
#108 mule: nothing has landed in 15h [mule-alarm ] []c
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

→ Read .mule/for-giulio.md

$ gh issue list --state open --limit 60 --json number,title,labels,updatedAt --template '{{range .}}#{{.number}} [{{range .labels}}{{.name}},{{end}}] {{.title}}{{"\n"}}{{end}}'
#123 [mule-alarm,] mule: nothing has landed in 14h
#122 [mule-alarm,] mule: nothing has landed in 15h
#121 [mule-alarm,] mule: nothing has landed in 16h
#120 [mule-alarm,] mule: nothing has landed in 14h
#119 [mule-alarm,] mule: nothing has landed in 16h
#118 [mule-alarm,] mule: nothing has landed in 14h
#117 [mule-alarm,] mule: nothing has landed in 15h
#116 [mule-alarm,] mule: nothing has landed in 14h
#115 [mule-alarm,] mule: nothing has landed in 14h
#114 [mule-alarm,] mule: nothing has landed in 8h
#113 [mule-alarm,] mule: nothing has landed in 14h
#112 [mule-alarm,] mule: nothing has landed in 14h
#111 [mule-alarm,] mule: nothing has landed in 16h
#110 [mule-alarm,] mule: nothing has landed in 14h
#109 [mule-alarm,] mule: nothing has landed in 14h
#108 [mule-alarm,] mule: nothing has landed in 15h
#107 [mule-alarm,] mule: nothing has landed in 14h
#106 [mule-alarm,] mule: nothing has landed in 14h
#105 [mule-alarm,] mule: nothing has landed in 14h
#104 [mule-alarm,] mule: nothing has landed in 14h
#103 [mule-alarm,] mule: nothing has landed in 16h
#102 [mule-alarm,] mule: nothing has landed in 14h
#101 [mule-alarm,] mule: nothing has landed in 14h
#100 [mule-alarm,] mule: nothing has landed in 15h
#99 [mule-alarm,] mule: nothing has landed in 14h
#98 [mule-alarm,] mule: nothing has landed in 16h
#97 [mule-alarm,] mule: nothing has landed in 15h
#96 [mule-alarm,] mule: nothing has landed in 14h
#95 [mule-alarm,] mule: nothing has landed in 15h
#94 [mule-alarm,] mule: nothing has landed in 15h
#93 [mule-review,upstream-parity,] In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
#92 [enhancement,upstream-parity,upstream-experimental,] keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds
#78 [enhancement,milestone-4.0,upstream-parity,] FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
#1 [wontfix,] Provide an Open Source IoT Platform unironically, unironically

$ echo "=== heartbeat/alarmed ==="; ls -la .mule/.heartbeat .mule/.alarmed 2>&1; echo "--- heartbeat value:"; cat .mule/.heartbeat 2>&1; echo; echo "--- date now:"; date -u +%Y-%m-%dT%H:%M:%SZ; echo "=== log.md tail ==="; tail -n 20 .mule/log.md 2>&1; echo "=== git divergence ==="; git rev-list --left-right --count origin/main...HEAD 2>&1; echo "=== HEAD ==="; git log --oneline -3 2>&1
=== heartbeat/alarmed ===
-rw-r--r-- 1 root root  0 Oct 10 12:50 .mule/.alarmed
-rw-r--r-- 1 root root 11 Oct  9 22:26 .mule/.heartbeat
--- heartbeat value:
1791577584

--- date now:
2026-10-10T11:18:35Z
=== log.md tail ===
| 2026-10-06 | docs-native-socket-handshake-errors | blocked | 91s | wrote nothing |
| 2026-10-06 | docs-native-socket-sse-exact-accept | done | 410s | 1322214 |
| 2026-10-06 | docs-native-socket-security-scheme | done | 299s | 87c18cf |
| 2026-10-06 | docs-native-socket-event-schema | blocked | 662s | opencode exited 1 |
| 2026-10-07 | pairing-burn-bcrypt-cost | blocked | 359s | lint failed: internal/pairing/service.go:95:7: G101: Potential hardcoded credentials (gosec) |
| 2026-10-07 | pairing-onregistered-clock | done | 143s | 4e71533 |
| 2026-10-07 | pairing-unregister-bad-id-test | done | 225s | 322242a |
| 2026-10-08 | testutil-wait-since-cursor | done | 319s | 3654fda |
| 2026-10-08 | testutil-control-frame-test | blocked | 102s | wrote nothing |
| 2026-10-08 | engine-introspection-producer-roundtrip | done | 234s | d070345 |
| 2026-10-08 | docs-sync-pairing-register-404-unreachable | done | 148s | cfa6e76 |
| 2026-10-08 | docs-sync-pairing-body-cap-400 | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-10-08 | docs-sync-pairing-realm-name-pattern | done | 410s | ed301fb |
| 2026-10-08 | docs-sync-pairing-deviceid-path-param | done | 395s | 1cda0fb |
| 2026-10-09 | docs-sync-ae-data-envelope-metadata | done | 245s | 2ca634d |
| 2026-10-09 | docs-sync-ae-data-write-object-errors | done | 587s | e974103 |
| 2026-10-09 | docs-sync-ae-group-create-devices-required | done | 343s | 9a2cf24 |
| 2026-10-09 | docs-sync-ae-notfound-component | blocked | 220s | wrote nothing |
| 2026-10-09 | realm-error-mapping-test | done | 175s | aef0914 |
| 2026-10-09 | realm-detailed-properties-branch-test | done | 88s | 76f80e9 |
=== git divergence ===
4	758
=== HEAD ===
4d78c9b mule: recipe milestones ran (403s)
69f3cff mule: recipe docs-sync timed out
1114f61 mule: log realm-detailed-properties-branch-test

$ echo "=== upstream releases ==="; gh api repos/astarte-platform/astarte/releases --jq '.[0:4][] | "\(.tag_name) prerelease=\(.prerelease) \(.published_at)"' 2>&1; echo "=== upstream tags ==="; gh api 'repos/astarte-platform/astarte/tags?per_page=100' --jq '.[0:6][].name' 2>&1
=== upstream releases ===
v1.3.5 prerelease=false 2026-10-05T14:11:44Z
v1.4.0-rc.6 prerelease=true 2026-09-30T07:28:27Z
v1.3.4 prerelease=false 2026-09-18T09:42:57Z
v1.4.0-rc.5 prerelease=true 2026-08-20T12:28:04Z
=== upstream tags ===
v1.4.0-rc.6
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1

$ echo "=== issue counts ==="; gh issue list --state open --limit 100 --json number --jq 'length' 2>&1; echo "=== label mule-alarm count ==="; gh issue list --state open --label mule-alarm --limit 100 --json number --jq 'length' 2>&1; echo "=== newest alarm ==="; gh issue list --state open --label mule-alarm --limit 1 --json number,title,createdAt --jq '.[]' 2>&1
=== issue counts ===
34
=== label mule-alarm count ===
30
=== newest alarm ===
{"createdAt":"2026-10-10T10:49:50Z","number":123,"title":"mule: nothing has landed in 14h"}

$ echo "=== #93 ancestry ==="; git merge-base --is-ancestor 8c61268 origin/main && echo "8c61268 IS ancestor of origin/main" || echo "8c61268 NOT ancestor of origin/main"; git branch -a --contains 8c61268 2>&1 | head; echo "=== search aclhook rewrite commits ==="; git log --oneline --all --grep 'issue-93' 2>&1 | head
=== #93 ancestry ===
8c61268 NOT ancestor of origin/main
* mule/queue
  remotes/origin/mule/queue
=== search aclhook rewrite commits ===
8c61268 mule: log issue-93
24ad5b8 mule: issue-93: In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged

$ echo "=== heartbeat as ISO ==="; date -u -d @1791577584 +%Y-%m-%dT%H:%M:%SZ; echo "=== #123 window ==="; python3 -c "
from datetime import datetime,timezone
h=datetime.fromtimestamp(1791577584,timezone.utc)
c=datetime.fromisoformat('2026-10-10T10:49:50+00:00')
print('heartbeat',h.isoformat())
print('delta',c-h)
print('age_h',int((c-h).total_seconds()//3600))
"
=== heartbeat as ISO ===
2026-10-09T20:26:24Z
=== #123 window ===
heartbeat 2026-10-09T20:26:24+00:00
delta 14:23:26
age_h 14

→ Read internal/broker/aclhook.go [offset=55, limit=25]

$ gh issue view 93 --json number,createdAt,updatedAt,labels --jq '"- #\(.number) created \(.createdAt) updated \(.updatedAt) labels \(.labels|map(.name))"' 2>&1; gh issue view 92 --json number,createdAt,updatedAt --jq '"- #\(.number) created \(.createdAt) updated \(.updatedAt)"' 2>&1; echo "=== days for 93 ==="; python3 -c "
from datetime import datetime,timezone
c=datetime.fromisoformat('2026-09-04T19:54:58+00:00')
now=datetime.fromisoformat('2026-10-10T11:18:35+00:00')
print((now-c).days,'days')"
- #93 created 2026-09-04T19:44:12Z updated 2026-09-04T19:54:58Z labels ["mule-review","upstream-parity"]
- #92 created 2026-09-04T19:43:54Z updated 2026-09-04T19:52:21Z
=== days for 93 ===
35 days

No machine-checkable issue is proposable (the pile is all `mule-alarm` noise; #93/#92/#78 are review/parked/design). Following the recipe's escalation channel, I'll append the triage result to `.mule/for-giulio.md` — no `todo.md` lines, no GitHub mutations.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,11 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **github-issues triage run, 2026-10-10: nothing proposable for the 29th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub.** The recipe's `--limit 40` command prints all **34** open issues: **30** `mule-alarm` **#94–#123** plus the same four non-alarm, none proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** untouched per standing instruction. New since the 2026-10-09 run: **#123**, created `2026-10-10T10:49:50Z`, "nothing has landed in 14h" — window pinned, not inferred: `.mule/.heartbeat` holds `1791577584` = `2026-10-09T20:26:24Z` (the last `done` land, `realm-detailed-properties-branch-test`/`76f80e9`, per `.mule/log.md`), and `10:49:50Z − 20:26:24Z` = **14h23m**, so `age=$(( (now-last)/3600 ))` prints 14 — the title verbatim. The latch is set again: `.mule/.alarmed` mtime `10:50Z`, written by that filing, and heartbeat is unmoved, so the dead-man's switch is off through today's silence exactly as measured on 2026-10-07/08; today's ticks were recipe verdicts (`4d78c9b` milestones, `69f3cff` docs-sync timeout), which neither beat nor un-latch. **Proposal, extending the 2026-10-09 one: close #94–#122, all twenty-nine** — each is a one-day alarm superseded by the next and nothing reads issue state; leave **#123** (today's, real window) to self-expire.
+  **#93 is now 35 days unreviewed and its in-tree comment's version reference has aged another day.** The aclhook rewrite (`24ad5b8`, "mule: issue-93: …") and its log commit `8c61268` are still not ancestors of `origin/main` (`git merge-base --is-ancestor 8c61268 origin/main` → false; `git branch -a --contains 8c61268` lists only `mule/queue` / `origin/mule/queue`), divergence `git rev-list --left-right --count origin/main...HEAD` = **4 / 758** (4 / 747 on 2026-10-09). The rewritten comment at the `control/keyAgreement` deny (internal/broker/aclhook.go:72-76) still reads "v1.4.0-rc.5 … v1.3.3 being the newest stable tag", but upstream's newest stable is **v1.3.5** and newest overall **v1.4.0-rc.6** (`gh api …/releases`, re-checked today) — one word in the same diff when someone finally runs `bash tools/mule.sh review`, not a task line. **#92's parking condition re-verified unmet** by the same sweep: still no stable `v1.4.0`, so an `-rc.N` does not satisfy your 2026-09-04 decision. No evidence file to write; the recipe names none.
+
 - **github-issues triage run, 2026-10-09: nothing proposable for the 28th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub.** The recipe's `--limit 40` command prints all **33** open issues: **29** `mule-alarm` **#94–#122** plus the same four non-alarm, none proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** untouched per standing instruction. New since the 2026-10-08 run: **#122**, created `2026-10-09T10:55:45Z`, "nothing has landed in 15h" — a genuine overnight window, pinned not inferred: it measured back to `.mule/.heartbeat` `1791488167` = `2026-10-08T19:36:07Z` (the last `done` land, the value the 10-08 evening run recorded), and `10:55:45Z − 19:36:07Z` = **15h19m**, so `age=$(( (now-last)/3600 ))` prints 15 — the title verbatim. The latch has since cleared: `.mule/.alarmed` is absent and heartbeat is `1791568277` = `2026-10-09T17:51:17Z`, because work landed today — `.mule/log.md` carries four `2026-10-09` rows, three `done` (`2ca634d`, `e974103`, `9a2cf24`) and one blocked. **Proposal, extending the 2026-10-08 evening one: close #94–#121, all twenty-eight** — each is a one-day alarm superseded by the next and nothing reads issue state; leave **#122** (today's, real window, already ended) to self-expire.
   **#93 is now 35 days unreviewed and its in-tree comment's version reference has aged another tag.** The aclhook rewrite (`8c61268`, "mule: log issue-93") is still not an ancestor of `origin/main` (`git merge-base --is-ancestor 8c61268 origin/main` → false; `git branch -a --contains 8c61268` lists only `mule/queue` / `origin/mule/queue`), divergence `git rev-list --left-right --count origin/main...HEAD` = **4 / 747**. The rewritten comment at the `control/keyAgreement` deny still reads "v1.4.0-rc.5 … v1.3.3 being the newest stable tag", but upstream's newest overall is `v1.4.0-rc.6` and newest stable is `v1.3.5` (`gh api .../releases`, re-checked today) — one word in the same diff when someone finally runs `bash tools/mule.sh review`, not a task line. **#92's parking condition re-verified unmet** by the same release sweep: still no stable `v1.4.0`, so an `-rc.N` does not satisfy your 2026-09-04 decision. No evidence file to write; the recipe names none.
 
 - **docs-sync run, 2026-10-09, surface: appengine — `docs/site/appengine-api.md` mislabels the `downsample_to` datastream query parameter as a duration.** The datastream-query snippet writes it as `&downsample_to=<bucket_duration>` (line 60) and the note below says "`downsample_to` maps onto Timescale `time_bucket()`" (line 65), but the code reads it as a point COUNT, not a time span: `parseQueryOpts` parses an integer and requires it `> 2`, storing it in `opts.DownsamplePoints` (internal/appengine/http.go:588-598), and the bucket interval is computed *from* that count downstream (`bucketFor(last.Sub(first), opts.DownsamplePoints)` then `s.Downsample(ctx, q, bucket)`, internal/appengine/data.go:187-196) — which is what the OpenAPI already says ("Downsample to this many data points (must be > 2)", docs/api/astarte_appengine_api.yaml:1509-1514). A reader following the page sends a duration like `1h` and gets a 422 (`downsample_to is invalid`, http.go:591-592). Proposed wording (your voice): `&downsample_to=<max_points>` and "`downsample_to` is the target number of points (must be greater than 2); the server chooses the Timescale `time_bucket()` interval." Page untouched.

$ git status --porcelain 2>&1
M .mule/for-giulio.md

→ Read .mule/for-giulio.md [limit=18]
Done: github-issues triage run #29 — nothing machine-checkable to propose; result escalated to `.mule/for-giulio.md`.
Files: `.mule/for-giulio.md` (only file changed)
Verified: `gh issue list` (34 open = 30 `mule-alarm` #94–#123 + #93/#92/#78/#1, all review/parked/design/excluded) → no proposable issue; no `gh issue create`, no edits/comments/closes; `git status --porcelain` shows only `for-giulio.md`.
Unsure: nothing.
Follow-ups: none in `.mule/todo.md` (no candidate qualified); one for-giulio entry — 29th consecutive "nothing proposable", re-proposing close of #94–#122 and noting #93 still unreviewed, #92 still parked.
