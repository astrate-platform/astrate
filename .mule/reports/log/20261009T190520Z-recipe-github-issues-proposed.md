slug: recipe-github-issues
verdict: proposed
at:  9660f70
ran: 2026-10-09T19:05:20Z on DietPi in 131s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md

$ gh issue list --state open --limit 40 --json number,title,labels,updatedAt,comments --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
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

→ Read .mule/todo.md

$ gh issue view 93 --json number,title,body,state,comments --template '{{.title}}{{"\n"}}{{.body}}' 2>&1 | head -80
In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
Comment-only change. The deny itself stays; only its stated justification is stale.

`internal/broker/aclhook.go` currently reads:

```go
// control/keyAgreement (080-mqtt-v1-protocol.md:172-176) is deliberately
// unsupported: upstream itself documents the handshake as "reserved and
// routed correctly, but not yet implemented", so falling through to the
// deny path below matches upstream's own status rather than being a gap.
```

That reasoning no longer holds on either half. Verified against upstream `master`:

- The quoted sentence, and the companion "accepted and acknowledged by Astarte, but no response is sent until the feature is fully implemented", are **gone from both** `080-mqtt-v1-protocol.md` and the new `082-key_agreement_protocol.md` (grep count 0 in each).
- `082-key_agreement_protocol.md` (13,832 bytes, added by upstream `d084308` on 2026-08-31) states the handshake "is implemented on the Astarte side (`astarte_data_updater_plant`)".
- The 080 topic table row also changed shape: a single `control/keyAgreement` at QoS 1 became `control/keyAgreement/<message id>` at QoS 2.

So the comment now tells a reader something false about upstream, and a future reader could reasonably "fix the gap" on the strength of it.

**Do:** rewrite the comment so it stands on Astrate's own decision — the topic is denied because full key agreement is out of scope for the current 1.2.2 target and the spec ships only in `v1.4.0-rc.5` (`v1.3.3` is the newest stable tag) — cite `082-key_agreement_protocol.md` as the current upstream reference rather than the deleted 080 sentence, and point at issue #92 where the implement-or-keep-parking decision lives.

**Do not:** change the ACL behaviour, and do not touch the expectation at `internal/broker/aclhook_test.go:53` (`pub control keyAgreement` still denies).

**Acceptance:** `go build ./...` and `go test ./internal/broker/...` pass unchanged; the comment no longer contains the phrase "not yet implemented"; `aclhook_test.go:53` is untouched.

Survey source: `survey-2026-08-31.md` on `mule/research`.

→ Read internal/broker/aclhook.go [offset=55, limit=30]

$ git log --oneline --all | grep -i 8c61268; echo "---branches---"; git branch -a 2>/dev/null; echo "---current---"; git log --oneline -5
8c61268 mule: log issue-93
---branches---
  main
* mule/queue
  remotes/origin/HEAD -> origin/main
  remotes/origin/feat/pokemon-agent
  remotes/origin/fix/blocked-mule-issues
  remotes/origin/m12-06b-watch-path-onto-main
  remotes/origin/main
  remotes/origin/mule/queue
  remotes/origin/mule/research
  remotes/origin/wip/DietPi-20260727T171543Z
  remotes/origin/worktree-m12-05-match-upstream
  remotes/origin/worktree-m12-06-ach-join
  remotes/origin/worktree-m12-06b-watch-path
---current---
9660f70 mule: recipe milestones ran (181s)
cf5d7b4 mule: recipe hygiene ran (811s)
67809c0 mule: blocked docs-sync-ae-notfound-component
d68d111 mule: log docs-sync-ae-group-create-devices-required
9a2cf24 mule: docs-sync-ae-group-create-devices-required [auto]: in docs/api/astarte_appengine_api.yaml, `GroupCreate` (yaml:1652-1663) marks only `group_name` required, but `createGroup` treats `devices` as required — `groupBody.Devices *[]string` (internal/appengine/http.go:391-397) makes a MISSING key a 422 `{"errors":{"devices":["can't be blank"]}}` and an EMPTY array a 422 `{"errors":{"devices":["should have at least 1 item(s)"]}}` (http.go:405-418), and only the MISSING-key case is unreachable from the schema as written. The shared `ValidationErrors` example (yaml:1841-1844) that POST /groups refs (yaml:948) shows `group_name: can't be blank` alone. Add `devices` to `GroupCreate`'s `required` list (the 201 echo always carries it, http.go:424) and add the two device messages to the POST 422 story. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ ls -la /root/astrate-mule/.mule/; echo "---"; ls /root/astrate-mule/.mule/recipes/
total 16420
drwxr-xr-x  7 root root     4096 Oct  9 20:07 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       13 Oct  9 21:03 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  9 19:51 .heartbeat
-rw-r--r--  1 root root     8192 Oct  9 21:03 .last-output
-rw-r--r--  1 root root        0 Oct  9 21:03 .lock
-rw-r--r--  1 root root       14 Oct  9 21:03 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 16361626 Oct  9 21:03 cron.log
drwxr-xr-x  2 root root     4096 Oct  9 20:07 failed
-rw-r--r--  1 root root   161185 Oct  9 13:28 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    16266 Oct  9 20:07 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  8 19:28 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   175392 Oct  9 20:07 todo.md
---
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

→ Read .mule/for-giulio.md

$ cd /root/astrate-mule; echo "--- divergence ---"; git rev-list --left-right --count origin/main...HEAD; echo "--- 8c61268 ancestor of origin/main? ---"; git merge-base --is-ancestor 8c61268 origin/main && echo ANCESTOR || echo NOT-ANCESTOR; echo "--- branches containing 8c61268 ---"; git branch -a --contains 8c61268 2>/dev/null; echo "--- newest log rows ---"; tail -c 1200 /root/astrate-mule/.mule/log.md
--- divergence ---
4	747
--- 8c61268 ancestor of origin/main? ---
NOT-ANCESTOR
--- branches containing 8c61268 ---
* mule/queue
  remotes/origin/mule/queue
--- newest log rows ---

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

This is a recurring triage run. Let me gather the live facts.
$ date -u; echo "--- heartbeat ---"; cat /root/astrate-mule/.mule/.heartbeat; echo; echo "--- alarmed? ---"; ls -la /root/astrate-mule/.mule/.alarmed 2>&1; echo "--- open issues count ---"; gh issue list --state open --limit 100 --json number | wc -l; echo "--- all open numbers ---"; gh issue list --state open --limit 100 --json number,labels --template '{{range .}}#{{.number}} [{{range .labels}}{{.name}} {{end}}]{{"\n"}}{{end}}'
Fri Oct  9 19:04:13 UTC 2026
--- heartbeat ---
1791568277

--- alarmed? ---
ls: cannot access '/root/astrate-mule/.mule/.alarmed': No such file or directory
--- open issues count ---
1
--- all open numbers ---
#122 [mule-alarm ]
#121 [mule-alarm ]
#120 [mule-alarm ]
#119 [mule-alarm ]
#118 [mule-alarm ]
#117 [mule-alarm ]
#116 [mule-alarm ]
#115 [mule-alarm ]
#114 [mule-alarm ]
#113 [mule-alarm ]
#112 [mule-alarm ]
#111 [mule-alarm ]
#110 [mule-alarm ]
#109 [mule-alarm ]
#108 [mule-alarm ]
#107 [mule-alarm ]
#106 [mule-alarm ]
#105 [mule-alarm ]
#104 [mule-alarm ]
#103 [mule-alarm ]
#102 [mule-alarm ]
#101 [mule-alarm ]
#100 [mule-alarm ]
#99 [mule-alarm ]
#98 [mule-alarm ]
#97 [mule-alarm ]
#96 [mule-alarm ]
#95 [mule-alarm ]
#94 [mule-alarm ]
#93 [mule-review upstream-parity ]
#92 [enhancement upstream-parity upstream-experimental ]
#78 [enhancement milestone-4.0 upstream-parity ]
#1 [wontfix ]

$ gh api 'repos/astarte-platform/astarte/releases?per_page=8' --template '{{range .}}{{.tag_name}} prerelease={{.prerelease}} {{.published_at}}{{"\n"}}{{end}}' 2>&1; echo "--- tags ---"; gh api 'repos/astarte-platform/astarte/tags?per_page=8' --template '{{range .}}{{.name}}{{"\n"}}{{end}}' 2>&1
v1.3.5 prerelease=false 2026-10-05T14:11:44Z
v1.4.0-rc.6 prerelease=true 2026-09-30T07:28:27Z
v1.3.4 prerelease=false 2026-09-18T09:42:57Z
v1.4.0-rc.5 prerelease=true 2026-08-20T12:28:04Z
v1.3.3 prerelease=false 2026-08-07T14:33:58Z
v1.4.0-rc.4 prerelease=true 2026-08-06T13:40:20Z
v1.4.0-rc.3 prerelease=true 2026-07-31T10:48:39Z
v1.4.0-rc.2 prerelease=true 2026-07-14T10:16:22Z
--- tags ---
v1.4.0-rc.6
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
v1.4.0-rc.0
v1.3.5

$ cd /root/astrate-mule; echo "heartbeat epoch 1791568277 ->"; date -u -d @1791568277; echo "--- newest land commit time ---"; git log -1 --format='%h %ci %s' 9a2cf24; echo "--- HEAD time ---"; git log -1 --format='%h %ci %s' HEAD; echo "--- today's log rows ---"; grep '| 2026-10-09 ' /root/astrate-mule/.mule/log.md
heartbeat epoch 1791568277 ->
Fri Oct  9 17:51:17 UTC 2026
--- newest land commit time ---
9a2cf24 2026-10-09 19:51:14 +0200 mule: docs-sync-ae-group-create-devices-required [auto]: in docs/api/astarte_appengine_api.yaml, `GroupCreate` (yaml:1652-1663) marks only `group_name` required, but `createGroup` treats `devices` as required — `groupBody.Devices *[]string` (internal/appengine/http.go:391-397) makes a MISSING key a 422 `{"errors":{"devices":["can't be blank"]}}` and an EMPTY array a 422 `{"errors":{"devices":["should have at least 1 item(s)"]}}` (http.go:405-418), and only the MISSING-key case is unreachable from the schema as written. The shared `ValidationErrors` example (yaml:1841-1844) that POST /groups refs (yaml:948) shows `group_name: can't be blank` alone. Add `devices` to `GroupCreate`'s `required` list (the 201 echo always carries it, http.go:424) and add the two device messages to the POST 422 story. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
--- HEAD time ---
9660f70 2026-10-09 20:38:18 +0200 mule: recipe milestones ran (181s)
--- today's log rows ---
| 2026-10-09 | docs-sync-ae-data-envelope-metadata | done | 245s | 2ca634d |
| 2026-10-09 | docs-sync-ae-data-write-object-errors | done | 587s | e974103 |
| 2026-10-09 | docs-sync-ae-group-create-devices-required | done | 343s | 9a2cf24 |
| 2026-10-09 | docs-sync-ae-notfound-component | blocked | 220s | wrote nothing |

$ gh issue view 122 --json number,title,createdAt,author,body --template '{{.number}} {{.title}}{{"\n"}}created={{.createdAt}} author={{.author.login}}{{"\n"}}---body---{{"\n"}}{{.body}}' 2>&1
122 mule: nothing has landed in 15h
created=2026-10-09T10:55:45Z author=atsetilam
---body---
The mule has not landed any work in 15 hours (threshold 8h).

It is probably still ticking and exiting 0 — that is how this failure looks. Worth checking,
on the Pi:

    systemctl list-timers mule.timer
    journalctl -u mule.service --since '-1 day' | tail -50
    cd /root/astrate-mule && bash tools/mule.sh status

Common causes, all of which have happened: the provider is refusing every run, the queue has
nothing runnable in it, the working tree is dirty so every tick aborts, or push has been
failing and the work is only on the Pi.

Filed automatically. Close it — the mule reopens a new one if the silence continues.

Nothing machine-checkable among the 33 open issues; the bulk are the recurring one-day `mule-alarm` pile. I'll record this run's survey and escalate the standing closure proposal in `for-giulio.md` (no `gh` state changes, no todo lines).
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,11 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **github-issues triage run, 2026-10-09: nothing proposable for the 28th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub.** The recipe's `--limit 40` command prints all **33** open issues: **29** `mule-alarm` **#94–#122** plus the same four non-alarm, none proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** untouched per standing instruction. New since the 2026-10-08 run: **#122**, created `2026-10-09T10:55:45Z`, "nothing has landed in 15h" — a genuine overnight window, pinned not inferred: it measured back to `.mule/.heartbeat` `1791488167` = `2026-10-08T19:36:07Z` (the last `done` land, the value the 10-08 evening run recorded), and `10:55:45Z − 19:36:07Z` = **15h19m**, so `age=$(( (now-last)/3600 ))` prints 15 — the title verbatim. The latch has since cleared: `.mule/.alarmed` is absent and heartbeat is `1791568277` = `2026-10-09T17:51:17Z`, because work landed today — `.mule/log.md` carries four `2026-10-09` rows, three `done` (`2ca634d`, `e974103`, `9a2cf24`) and one blocked. **Proposal, extending the 2026-10-08 evening one: close #94–#121, all twenty-eight** — each is a one-day alarm superseded by the next and nothing reads issue state; leave **#122** (today's, real window, already ended) to self-expire.
+  **#93 is now 35 days unreviewed and its in-tree comment's version reference has aged another tag.** The aclhook rewrite (`8c61268`, "mule: log issue-93") is still not an ancestor of `origin/main` (`git merge-base --is-ancestor 8c61268 origin/main` → false; `git branch -a --contains 8c61268` lists only `mule/queue` / `origin/mule/queue`), divergence `git rev-list --left-right --count origin/main...HEAD` = **4 / 747**. The rewritten comment at the `control/keyAgreement` deny still reads "v1.4.0-rc.5 … v1.3.3 being the newest stable tag", but upstream's newest overall is `v1.4.0-rc.6` and newest stable is `v1.3.5` (`gh api .../releases`, re-checked today) — one word in the same diff when someone finally runs `bash tools/mule.sh review`, not a task line. **#92's parking condition re-verified unmet** by the same release sweep: still no stable `v1.4.0`, so an `-rc.N` does not satisfy your 2026-09-04 decision. No evidence file to write; the recipe names none.
+
 - **docs-sync run, 2026-10-09, surface: appengine — `docs/site/appengine-api.md` mislabels the `downsample_to` datastream query parameter as a duration.** The datastream-query snippet writes it as `&downsample_to=<bucket_duration>` (line 60) and the note below says "`downsample_to` maps onto Timescale `time_bucket()`" (line 65), but the code reads it as a point COUNT, not a time span: `parseQueryOpts` parses an integer and requires it `> 2`, storing it in `opts.DownsamplePoints` (internal/appengine/http.go:588-598), and the bucket interval is computed *from* that count downstream (`bucketFor(last.Sub(first), opts.DownsamplePoints)` then `s.Downsample(ctx, q, bucket)`, internal/appengine/data.go:187-196) — which is what the OpenAPI already says ("Downsample to this many data points (must be > 2)", docs/api/astarte_appengine_api.yaml:1509-1514). A reader following the page sends a duration like `1h` and gets a 422 (`downsample_to is invalid`, http.go:591-592). Proposed wording (your voice): `&downsample_to=<max_points>` and "`downsample_to` is the target number of points (must be greater than 2); the server chooses the Timescale `time_bucket()` interval." Page untouched.
 
 - **github-issues triage run, 2026-10-08 (evening, 20:37Z): nothing proposable for the 27th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub.** The recipe's `--limit 40` command prints all **32** open issues: still **28** `mule-alarm` **#94–#121** (no #122 — nothing new filed since this morning's run) plus the same four non-alarm, none proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), one excluded by standing instruction. **The alarm latch is working again and #121's window has ended**: `.mule/.alarmed` is gone, `.mule/.heartbeat` = `1791488167` = **2026-10-08T19:36:07Z**, and `.mule/log.md` now carries four `2026-10-08` `done` rows the morning run could not see (`d070345`, `cfa6e76`, `ed301fb`, `1cda0fb`, plus one 1200s timeout) — the very land that clears the latch (mule.sh:578→737), so #121 needs no action either. **Proposal unchanged from earlier today: close #94–#120, all twenty-seven** — each is a one-day alarm superseded by the next and nothing reads issue state; leave #121 to self-expire. Re-verified today, not inherited: **#92's parking condition still unmet** (`gh api .../tags` → newest `v1.4.0-rc.6`, `/releases` → newest stable `v1.3.5`; an `-rc.N` does not satisfy "wait for a stable v1.4.0"), and **#93's rewrite still entirely off main** — `git branch -a --contains 24ad5b8` lists only `mule/queue` / `origin/mule/queue`, `origin/main...HEAD` = 4 / **732** (4 / 711 this morning).
 

→ Read .mule/for-giulio.md [limit=18]
Surveyed all 33 open issues via the recipe's `gh issue list`: 29 are the one-day `mule-alarm` pile (#94–#122), plus #93 (`mule-review`), #92 (parked), #78 (`milestone-4.0`), and #1 (never touched). Zero machine-checkable candidates, so no `todo.md` lines and no `gh issue create` — this is a triage-only run.

Done: triaged 2026-10-09; nothing proposable, escalated the standing alarm-close proposal in `for-giulio.md`.
Files: `.mule/for-giulio.md` (new top entry only). `gh issue create` deliberately not run (recipe doesn't ask for it); no todo.md lines (none warranted).
Verified: `gh issue list` 33 open; `gh api .../releases` → newest stable v1.3.5, newest overall v1.4.0-rc.6; `git merge-base --is-ancestor 8c61268 origin/main` → false; divergence 4/747; `.mule/.alarmed` absent, heartbeat 2026-10-09T17:51:17Z.
Unsure: whether the alarm-#94–#122 pile should keep getting a fresh closure proposal each run rather than resting on the existing one — I extended it, matching prior practice.
Follow-ups: none added to `.mule/todo.md`.

No `gh issue edit/close/comment` and no git mutation was performed.
