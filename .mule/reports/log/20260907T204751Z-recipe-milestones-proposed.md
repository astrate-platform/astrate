slug: recipe-milestones
verdict: proposed
at:  ac78a09
ran: 2026-09-07T20:47:51Z on DietPi in 215s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
✱ Glob ".mule/**/*.md" 0 matches
✱ Glob "**/*parity*.md" 1 match
✱ Grep "upstream-parity|2026-08-22" 100 matches

$ ls -la .mule/ && echo "---research---" && ls -la .mule/research/ 2>/dev/null; echo "---tasks---"; ls -la .mule/tasks/ 2>/dev/null
total 1496
drwxr-xr-x  7 root root    4096 Sep  7 22:12 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      14 Sep  7 22:44 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep  7 22:12 .heartbeat
-rw-r--r--  1 root root       0 Sep  7 22:44 .last-output
-rw-r--r--  1 root root       0 Sep  7 22:44 .lock
-rw-r--r--  1 root root      11 Sep  7 22:44 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 1400670 Sep  7 22:44 cron.log
drwxr-xr-x  2 root root    4096 Sep  7 21:03 failed
-rw-r--r--  1 root root   20080 Sep  7 19:48 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    2868 Sep  7 22:12 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep  7 20:15 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   23678 Sep  7 22:12 todo.md
---research---
---tasks---
total 60
drwxr-xr-x 2 root root 4096 Sep  7 21:11 .
drwxr-xr-x 7 root root 4096 Sep  7 22:12 ..
-rw-r--r-- 1 root root 3293 Sep  7 21:11 docs-sync-appengine-data-output-params.md
-rw-r--r-- 1 root root 5539 Jul 27 22:18 issue-16.md
-rw-r--r-- 1 root root 7510 Sep  4 20:46 issue-23.md
-rw-r--r-- 1 root root 6434 Sep  4 20:46 issue-24.md
-rw-r--r-- 1 root root 5665 Sep  4 20:46 issue-25.md
-rw-r--r-- 1 root root 6009 Sep  4 20:46 issue-26.md
-rw-r--r-- 1 root root 5683 Sep  4 20:46 issue-27.md

→ Read .mule/reports/log/20260907T171928Z-recipe-milestones-proposed.md
→ Read .mule/for-giulio.md

$ date -u; echo "---"; git log --oneline -5; echo "---"; git branch --show-current
Mon Sep  7 20:46:27 UTC 2026
---
ac78a09 mule: recipe hygiene ran (110s)
4dee74d mule: log docs-sync-appengine-data-output-params
262fdd5 mule: docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
e16cf72 mule: refill the queue
8b575c3 mule: recipe docs-sync ran (341s)
---
mule/queue

The v3.0 milestone is the current target. There have been several prior runs of this exact recipe (2026-09-05, 09-06, and earlier today 09-07T17:19Z), all concluding the milestone is gated on externals. Let me verify nothing has moved since this morning's run before concluding.
→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ gh issue list --label milestone-3.0 --state all --limit 50 2>&1; echo "=== upstream-parity open ==="; gh issue list --label upstream-parity --state open --limit 50 2>&1; echo "=== upstream-parity counts ==="; gh issue list --label upstream-parity --state all --limit 60 2>&1 | wc -l; gh issue list --label upstream-parity --state closed --limit 60 2>&1 | wc -l
=== upstream-parity open ===
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
=== upstream-parity counts ===
46
43

State is unchanged from this morning's milestone run. Let me re-verify the two gates that actually move the milestone: upstream v1.4.0 final status, and the missing source doc.
$ ls -la .mule/research/ 2>&1 | head -3; echo "---"; git ls-tree -r origin/mule/research --name-only 2>&1 | rg -i 'upstream-parity|parity' ; echo "exit=$?"
ls: cannot access '.mule/research/': No such file or directory
---
exit=1

$ gh api repos/astarte-platform/astarte/releases --paginate -q '.[] | select(.prerelease==false) | .tag_name' 2>&1 | head -5; echo "=== prerelease tail ==="; gh api repos/astarte-platform/astarte/releases --paginate -q '.[] | select(.prerelease==true) | .tag_name' 2>&1 | head -8
v1.3.3
v1.3.2
v1.3.1
v1.3.0
v1.2.2
=== prerelease tail ===
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
v1.4.0-rc.0
v1.3.0-rc.2
v1.2.2-rc.0

Everything is unchanged since this morning's run. Let me confirm the report/logging conventions and that the escalated state in `for-giulio.md` is current.
$ ls -la .mule/reports/ .mule/reports/log/ 2>/dev/null | head -40; echo "=== reviews ==="; ls -la .mule/reviews/ 2>/dev/null
.mule/reports/:
total 16
drwxr-xr-x 3 root root 4096 Sep  6 19:40 .
drwxr-xr-x 7 root root 4096 Sep  7 22:12 ..
-rw-r--r-- 1 root root 2279 Sep  6 19:40 hygiene.md
drwxr-xr-x 2 root root 4096 Sep  7 22:18 log

.mule/reports/log/:
total 1424
drwxr-xr-x 2 root root   4096 Sep  7 22:18 .
drwxr-xr-x 3 root root   4096 Sep  6 19:40 ..
-rw-r--r-- 1 root root  19332 Sep  4 21:09 20260904T190911Z-purge-properties-compression-capabilityauto-done.md
-rw-r--r-- 1 root root  25676 Sep  4 21:50 20260904T195018Z-broker-acl-coldstart-introspection-miss-done.md
-rw-r--r-- 1 root root   5576 Sep  4 21:55 20260904T195529Z-issue-93-done.md
-rw-r--r-- 1 root root    970 Sep  4 21:58 20260904T195814Z-broker-disconnect-device-zombie-session-blocked.md
-rw-r--r-- 1 root root   7195 Sep  4 22:30 20260904T203003Z-broker-offline-acl-tests-blocked.md
-rw-r--r-- 1 root root   3691 Sep  4 23:05 20260904T210532Z-broker-onconnect-doc-comment-blocked.md
-rw-r--r-- 1 root root  15724 Sep  5 13:02 20260905T110211Z-empty-introspection-verification-blocked.md
-rw-r--r-- 1 root root   4497 Sep  5 13:14 20260905T111428Z-probe-trigger-install-notification-delay-blocked.md
-rw-r--r-- 1 root root  12798 Sep  5 13:25 20260905T112539Z-compat-note-v132-done.md
-rw-r--r-- 1 root root  31514 Sep  5 18:55 20260905T165512Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root  13915 Sep  5 19:17 20260905T171755Z-docs-sync-pairing-health-path-done.md
-rw-r--r-- 1 root root   5035 Sep  5 19:22 20260905T172234Z-docs-sync-pairing-register-404-done.md
-rw-r--r-- 1 root root  32904 Sep  5 19:36 20260905T173657Z-recipe-hygiene-proposed.md
-rw-r--r-- 1 root root  11141 Sep  5 19:45 20260905T174528Z-swagger-httptest-coverage-done.md
-rw-r--r-- 1 root root 120837 Sep  5 20:10 20260905T181053Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  24739 Sep  5 20:12 20260905T181219Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root 152510 Sep  5 20:46 20260905T184625Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   9342 Sep  5 20:49 20260905T184903Z-probe-property-resend-encoding-blocked.md
-rw-r--r-- 1 root root   2052 Sep  5 21:08 20260905T190855Z-compat-note-v133-blocked.md
-rw-r--r-- 1 root root  38218 Sep  5 21:20 20260905T192010Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root  42265 Sep  5 21:29 20260905T192911Z-flow-validate-source-sink-blocked.md
-rw-r--r-- 1 root root   6007 Sep  5 22:15 20260905T201529Z-flow-validate-dead-source-sink-recompute-done.md
-rw-r--r-- 1 root root  41275 Sep  6 13:20 20260906T112006Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root  12532 Sep  6 13:31 20260906T113143Z-docs-sync-rm-datastream-retention-endpoint-done.md
-rw-r--r-- 1 root root  19509 Sep  6 19:10 20260906T171045Z-docs-sync-rm-interfaces-detailed-param-done.md
-rw-r--r-- 1 root root  12896 Sep  6 19:40 20260906T174031Z-recipe-hygiene-proposed.md
-rw-r--r-- 1 root root  17313 Sep  6 19:57 20260906T175748Z-housekeeping-tests-done.md
-rw-r--r-- 1 root root  63922 Sep  6 20:04 20260906T180442Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root   8576 Sep  6 20:13 20260906T181352Z-recipe-github-issues-proposed.md
=== reviews ===
total 36
drwxr-xr-x 2 root root 4096 Sep  7 20:15 .
drwxr-xr-x 7 root root 4096 Sep  7 22:12 ..
-rw-r--r-- 1 root root 4600 Sep  7 20:15 appengine-2026-09-07.md
-rw-r--r-- 1 root root 8188 Sep  4 20:43 broker-2026-09-04.md
-rw-r--r-- 1 root root 3141 Sep  6 20:57 engine-2026-09-06.md
-rw-r--r-- 1 root root 5516 Sep  5 21:19 flow-2026-09-05.md

→ Read .mule/log.md

$ rg -n 'milestone-3\.0|v3\.0|upstream-parity' .mule/todo.md 2>/dev/null | head -20; echo "exit=$?"
exit=0

Verified. Nothing has moved since this morning's milestone run (2026-09-07T17:19Z); every check re-confirms the same state:

**Current milestone: v3.0** (upstream 1.2.2 → 1.3/1.4). Scope is decided, not TBD, so no candidate-list entry is needed. Investigation per the recipe:

- `gh issue list --label milestone-3.0 --state all` → **empty** (this milestone's work lives under `upstream-parity`).
- `upstream-parity`: 43/46 closed, all of #47–#89. The 3 open are all already escalated, none a fileable gap:
  - **#92** keyAgreement — escalated, gated on a stable upstream v1.4.0 (newest stable is still v1.3.3; newest prerelease v1.4.0-rc.5).
  - **#93** aclhook comment rewrite — `mule-review`, commit `8c61268` awaiting review, its own recipe path.
  - **#78** FDO — milestone-4.0 candidate, escalated.
- **Standing UPSTREAM-EXPERIMENTAL item:** rows #67/#68 reconcile only when upstream 1.4 final ships — confirmed it hasn't (rc-only), and per #90's frozen decision the `APICompatVersion` bump is not due.
- Source doc `.mule/research/upstream-parity-2026-08-22.md` still absent from `main` and `origin/mule/research` — already escalated 2026-09-06; output #47–#89 intact.

**Proposed:** nothing. No machine-checkable gaps → filed no issues. Appended no `todo.md` lines (no implementable work is open on v3.0). Wrote no new for-giulio entry — the 2026-09-06 escalation already states this exact state, and the "complete, verify and cut the tag" line is deliberately **not** proposed because open items #92/#93 in the milestone's own backlog and the v1.4.0 gate are real decisions that are Giulio's. No git touched.

```
Done: v3.0 milestone recipe — state unchanged since 09-07T17:19Z; nothing to file, propose, or escalate (milestone gated on upstream v1.4.0 final + Giulio's #92 decision)
Files: none
Verified: gh issue list --label milestone-3.0|upstream-parity -> pass; gh api astarte releases stable->v1.3.3, prerelease->v1.4.0-rc.5; git ls-tree origin/mule/research parity doc absent -> pass
Unsure: nothing
Follow-ups: none
```
