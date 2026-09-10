slug: recipe-astarte-upstream
verdict: proposed
at:  eb44d02
ran: 2026-09-10T20:29:49Z on DietPi in 152s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

$ ls -la /root/astrate-mule/.mule/
total 2584
drwxr-xr-x  7 root root    4096 Sep 10 21:11 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      14 Sep 10 22:27 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 10 21:11 .heartbeat
-rw-r--r--  1 root root       0 Sep 10 22:27 .last-output
-rw-r--r--  1 root root       0 Sep 10 22:27 .lock
-rw-r--r--  1 root root      17 Sep 10 22:27 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 2484781 Sep 10 22:27 cron.log
drwxr-xr-x  2 root root    4096 Sep  8 20:45 failed
-rw-r--r--  1 root root   25105 Sep 10 21:52 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    4339 Sep 10 21:11 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 10 19:00 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   36222 Sep 10 21:11 todo.md

→ Read .mule/todo.md
→ Read .mule/for-giulio.md
→ Read .mule/log.md
Now let me check what the last `astarte-upstream` proposal run concluded, to avoid duplicating.

$ rg -n "astarte-upstream|astarte upstream|RETAIL|no upstream movement|proposal" /root/astrate-mule/.mule/cron.log | tail -80
24428:/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:109:Let me check the release notes history — what did earlier runs already cover for v1.4.0-rc.x? The `UPSTREAM-EXPERIMENTAL.md` register tracks #67 and #68. The previous recipe runs mentioned v1.4.0 as "still rc-only". 
24429:/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:111:Since v1.4.0-rc.5 is newer than what was previously tracked (v1.3.3), let me look at the release notes diffs for the rc releases. The previous run apparently only looked up to v1.3.3 (the newest stable at the time). Actually wait — the previous run happened and proposed compat-note-v1.3.3 saying "newest stable, 2026-08-07... v1.4.0 is still rc.5-only". So at that point v1.4.0-rc.5 already existed! The previous run already knew about rc.5.
24430:/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:113:So the question: did the previous run already examine v1.4.0-rc.5? It said "v1.4.0 is still rc.5-only" meaning it knew rc.5 was the latest rc. But did it propose anything for the rc.5 changes?
24431:/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:115:Let me look at the release notes for the rc versions more carefully to understand what's covered. The rc.5 body only has two fixes. Let me check what the previous run proposed regarding v1.4.0-rc.x changes. The todo doesn't show any v1.4.0-rc proposals other than the experimental register stuff already handled.
24433:The previous run (2026-09-07) already ran this recipe. Let me read its full report to avoid duplicating its proposals.
24434:→ Read .mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md
24436:→ Read .mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md [offset=490]
24456:20260905T184625Z-recipe-astarte-upstream-proposed.md
24469:20260906T183102Z-recipe-astarte-upstream-proposed.md
24477:20260907T020359Z-recipe-astarte-upstream-proposed.md
24497:→ Read .mule/reports/log/20260907T020359Z-recipe-astarte-upstream-proposed.md
24500:$ head -8 /root/astrate-mule/.mule/reports/log/20260906T183102Z-recipe-astarte-upstream-proposed.md; echo "===TAIL==="; tail -5 /root/astrate-mule/.mule/reports/log/20260906T183102Z-recipe-astarte-upstream-proposed.md
24501:slug: recipe-astarte-upstream
24515:$ ls /root/astrate-mule/.mule/reports/log/ | grep astarte-upstream; echo "==="; for f in /root/astrate-mule/.mule/reports/log/*recipe-astarte-upstream*.md; do echo "--- $f"; sed -n '1,4p' "$f"; done
24516:20260905T184625Z-recipe-astarte-upstream-proposed.md
24517:20260906T183102Z-recipe-astarte-upstream-proposed.md
24518:20260907T020359Z-recipe-astarte-upstream-proposed.md
24519:20260907T175524Z-recipe-astarte-upstream-proposed.md
24520:20260908T113422Z-recipe-astarte-upstream-proposed.md
24522:--- /root/astrate-mule/.mule/reports/log/20260905T184625Z-recipe-astarte-upstream-proposed.md
24523:slug: recipe-astarte-upstream
24527:--- /root/astrate-mule/.mule/reports/log/20260906T183102Z-recipe-astarte-upstream-proposed.md
24528:slug: recipe-astarte-upstream
24532:--- /root/astrate-mule/.mule/reports/log/20260907T020359Z-recipe-astarte-upstream-proposed.md
24533:slug: recipe-astarte-upstream
24537:--- /root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md
24538:slug: recipe-astarte-upstream
24542:--- /root/astrate-mule/.mule/reports/log/20260908T113422Z-recipe-astarte-upstream-proposed.md
24543:slug: recipe-astarte-upstream
24548:There's already a run of this recipe today (2026-09-08T11:34). Let me read it to see what it found, so I don't duplicate its proposals.
24549:→ Read .mule/reports/log/20260908T113422Z-recipe-astarte-upstream-proposed.md
24580:This recipe already ran earlier today (2026-09-08T11:34, sha `1e56fda`) and appended both proposals (`probe-required-mapping-flag`, `compat-note-v1.4-rc`); upstream is unchanged — newest tag still v1.4.0-rc.5 (2026-08-20). Let me check how no-movement outcomes are recorded, then write the evidence file.
24581:$ ls /root/astrate-mule/.mule/reports/log/ | grep -vE "astarte-upstream|docs-sync|hygiene|milestones|github-issues|code-review" | head; echo "===verdicts==="; grep -rh "^verdict:" /root/astrate-mule/.mule/reports/log/*.md | sort | uniq -c
24599:The recipe already ran today and proposed; this run found no upstream movement since. Let me check how the runner parses verdicts before writing the evidence file.
24602:← Write .mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md
24610:Done: upstream watch — no upstream movement since v1.4.0-rc.5
24611:Files: .mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md
24751:1. **The rate limiter double-Allow in `handleCredentials`** (http.go:158). Order matters: `!a.credLimiter.Allow("ip|"+ip.String()) || !a.credLimiter.Allow("dev|"+realm+"/"+deviceID)`. If IP bucket over limit → short circuit, dev not checked, rejected. If IP allows but dev over → rejected but IP token consumed. This is a minor token-waste issue. Not a strong proposal.
24816:Let me now consolidate my findings. This package is very thoroughly tested (unit tests cover service, ca, ratelimit; integration covers HTTP). Let me identify genuine, defensible proposals.
24828:Let me think about **missing tests** — highest-value proposals. What rules exist but aren't asserted in the container-free unit tests?
24845:Let me now identify genuine, defensible proposals. This package is among the most thoroughly tested I've seen. Let me settle on the strongest findings.
24857:Let me front-load the genuinely most valuable proposal. The recipe's #1 value is "a behaviour that is wrong or unguarded." I haven't found a clear wrong behaviour. Let me look harder at the rate limiter issue and the `Register`/`SetPayloadFormatHint` non-atomicity.
24873:But let me identify at least the defensible minor proposals that a careful reviewer would raise, in descending order:
24885:I think my two missing-test proposals are the right items, plus I should record in the review file what I decided *not* to propose (the cosmetic/nit items and the nil-CertSerial non-issue).
24903:These are both fully machine-checkable, container-free, and currently untested rules. Both are good "missing test" proposals (recipe value tier #2).
24979:When both sources are empty a tick runs a **proposal recipe** instead, rotating through
24980:`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
25023:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
25041:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
25050:The todo queue is long with existing blocked tasks. My code-review proposals should be appended. I have two solid missing-test proposals. Let me pick the strongest 2 and keep them tight.
25052:Let me now write the review file. First, let me double check `remoteIP` behavior once more for the review notes, and confirm my proposal about `bearerSecret` - both pure helpers in a non-integration file.
25068:Now let me verify my claims about `bearerSecret` once more against upstream — I want to make sure the proposal task line describes correct expected behavior. The upstream Astarte pairing regex is `~r/bearer\:?\s+(.*)$/i`. Let me confirm this is what the code comment claims (it does, at http.go:342). Since the code comment explicitly documents this regex, and the code implements it, my proposal to add tests locking that behavior is sound.
25109:The package was well-tested (service/CA/limiter unit coverage is deep). The only real gaps are two pure helpers in the non-integration `http.go` — the permissive `bearerSecret` bearer match and the unparseable-address fallback in `remoteIP` — each currently untested and each asserted only by an untestable-on-this-Mac integration suite. I recorded both as the sole proposals, and logged in the review file what I looked at and deliberately did not propose (the limiter token asymmetry, the nil-CertSerial edge, the hint non-atomicity) so the next run doesn't repeat the exercise.
26122:>> queue is empty — running the 'astarte-upstream' recipe to propose work
26127:→ Read .mule/recipes/astarte-upstream.md
27607:I've compared the astrate-native surface (the recipe's fifth option — no prior docs-sync run covers it). Code routes match the spec's documented paths one-for-one, so the drift is in response accuracy, not missing/phantom paths. Let me note the findings and append the proposals.
27625:Findings → 2 task lines appended (tagged `[auto]`, as todo.md requires for proposal output):
28805:a701b05 mule: recipe astarte-upstream ran (239s)
29051:.mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md
29101:20260908T193614Z-recipe-astarte-upstream-noop.md
29102:20260908T193814Z-recipe-astarte-upstream-proposed.md
29112:20260909T184315Z-recipe-astarte-upstream-proposed.md
29171:- **Gap buckets:** no machine-checkable gaps → nothing to file (the recipe cap "file at most five" is irrelevant — zero candidates). No new design decisions → nothing to escalate. The "complete, cut the tag" proposal (`for-giulio.md` step 5) is **not** warranted: the milestone's open decision set (#92 keyAgreement gated on stable v1.4.0) is real and is Giulio's call.
29523:>> queue is empty — running the 'astarte-upstream' recipe to propose work
29529:→ Read .mule/recipes/astarte-upstream.md
29864:There's a prior report of this recipe run. Let me check it to avoid duplicating proposals.
29865:→ Read .mule/reports/log/20260909T184315Z-recipe-astarte-upstream-proposed.md
29893:→ Read .mule/reports/log/20260909T184315Z-recipe-astarte-upstream-proposed.md [offset=1, limit=4]
29902:← Write .mule/reports/log/20260910T112515Z-recipe-astarte-upstream-proposed.md
29909:Done: no upstream movement since v1.4.0-rc.5; analyzed and nothing proposeable; appended no task lines
29910:Files: .mule/reports/log/20260910T112515Z-recipe-astarte-upstream-proposed.md
30032:Let me check what the `flow` package's tests cover to avoid duplicating proposals.
30468:Done: code review of internal/flowapi — 4 proposals: auto-restart retries forever on permanent failure, restart goroutine unreachable on shutdown, validationDetail untested, auto_restart default untested
30470:Verified: none run — proposal run, no code changed
31180:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
31941:When both sources are empty a tick runs a **proposal recipe** instead, rotating through
31942:`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
31989:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
32007:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
32360:>> queue is empty — running the 'astarte-upstream' recipe to propose work

Let me read the most recent `astarte-upstream` proposal run report (2026-09-10) to confirm the state and avoid duplication.
→ Read .mule/reports/log/20260910T112515Z-recipe-astarte-upstream-proposed.md
Last run (2026-09-10, sha 9a28cd1) already determined no upstream movement since v1.4.0-rc.5 and proposed nothing. Let me confirm the current state hasn't changed since then.

$ git rev-parse HEAD && echo "---deps-untouched check---" && gh api repos/astarte-platform/astarte/releases --jq '.[0:1][] | "\(.tag_name)\t\(.published_at)"'
eb44d02b00c7cd8fcc3db6d938748d5437085304
---deps-untouched check---
v1.4.0-rc.5	2026-08-20T12:28:04Z

→ Read .mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md
Newest tag is unchanged: still `v1.4.0-rc.5` (2026-08-20). The last recipe run (2026-09-10, sha 9a28cd1) already analyzed rc.4→rc.5 and every prior gap; the codebase has since moved to eb44d02 but upstream has not. Let me check the convention for a no-movement evidence file.

$ ls /root/astrate-mule/.mule/reports/log/ | tail -20
20260909T191420Z-recipe-code-review-proposed.md
20260909T195057Z-realm-pure-helper-tests-done.md
20260909T200459Z-recipe-docs-sync-proposed.md
20260909T201715Z-docs-sync-native-compat-health-503-done.md
20260909T202852Z-docs-sync-native-compat-version-endpoints-done.md
20260909T230006Z-recipe-hygiene-timeout.md
20260909T230411Z-recipe-milestones-proposed.md
20260910T110956Z-recipe-github-issues-proposed.md
20260910T112515Z-recipe-astarte-upstream-proposed.md
20260910T112545Z-recipe-astarte-upstream-proposed.md
20260910T170149Z-recipe-code-review-proposed.md
20260910T171629Z-flowapi-autorestart-terminal-failure-done.md
20260910T172606Z-flowapi-validationdetail-test-done.md
20260910T182250Z-flowapi-autorestart-default-test-done.md
20260910T184706Z-recipe-docs-sync-proposed.md
20260910T190522Z-docs-sync-pairing-status-enum-done.md
20260910T191129Z-docs-sync-pairing-version-endpoint-done.md
20260910T192802Z-recipe-hygiene-proposed.md
20260910T195304Z-recipe-milestones-proposed.md
20260910T200120Z-recipe-github-issues-proposed.md

Upstream has not moved since the last check (still `v1.4.0-rc.5`, 2026-08-20), and the prior run already analyzed everything and appended nothing. Per the recipe, an empty result is the correct result: append nothing, write the evidence file.
$ date -u +'%Y%m%dT%H%M%SZ'
20260910T202936Z

← Write .mule/reports/log/20260910T202936Z-recipe-astarte-upstream-noop.md
Wrote file successfully.

Done: upstream watch — no upstream movement since v1.4.0-rc.5; appended nothing to `.mule/todo.md`

Files: `.mule/reports/log/20260910T202936Z-recipe-astarte-upstream-noop.md` (evidence; todo.md untouched)
Verified: `gh api repos/astarte-platform/astarte/releases` -> newest remains v1.4.0-rc.5 (2026-08-20), unchanged since the last run; nothing new to analyze or propose
Unsure: nothing
Follow-ups: none
