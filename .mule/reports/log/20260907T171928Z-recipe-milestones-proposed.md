slug: recipe-milestones
verdict: proposed
at:  9d822c3
ran: 2026-09-07T17:19:28Z on DietPi in 312s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
Now I need to read `.mule/milestones.md` and `.mule/todo.md` (to see current state), plus check existing issues for the milestone label.
→ Read .mule/todo.md
→ Read .mule/for-giulio.md
The first non-DONE milestone is **v3.0**. The for-giulio.md entry from 2026-09-06 says there are no open machine-checkable gaps left to file. Let me verify by checking the current issue state.

$ cd /root/astrate-mule && gh issue list --label milestone-3.0 --state all --limit 50 2>&1
(no output)

$ cd /root/astrate-mule && gh issue list --label upstream-parity --state open --limit 50 2>&1
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z

→ Read docs/UPSTREAM-EXPERIMENTAL.md
v3.0 is the first non-DONE milestone. The `milestone-3.0` label is empty (that run's work lives under `upstream-parity`). Per the 2026-09-06 entry, the remaining open items (#92, #78, #93) are already escalated rather than fileable gaps. Let me verify nothing has moved since and check the standing UPSTREAM-EXPERIMENTAL item.

$ cd /root/astrate-mule && gh issue list --label upstream-parity --state all --limit 60 2>&1 | wc -l && gh issue list --label upstream-parity --state closed --limit 60 2>&1 | wc -l
46
43

State matches the 2026-09-06 run exactly: 43 of 46 `upstream-parity` issues closed, 3 open (#93 mule-review, #92 and #78 escalated), `milestone-3.0` label empty. UPSTREAM-EXPERIMENTAL rows #67/#68 still "1.4 experimental", gated on upstream v1.4.0 final. Let me check whether upstream v1.4.0 has gone stable since the last run — that is the gate the whole milestone waits on.
◈ Exa Web Search "astarte-platform astarte v1.4.0 release stable tag 2026"
v1.4.0 is still rc-only (rc.0 published 2026-04-08; latest tags still pre-release). Let me confirm the newest non-prerelease tag directly.

$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/releases --paginate -q '.[] | select(.prerelease==false) | .tag_name' 2>&1 | head -5
v1.3.3
v1.3.2
v1.3.1
v1.3.0
v1.2.2

v1.4.0 is still rc-only (newest stable remains v1.3.3). Let me confirm the overall open-issue picture and that nothing milestone-related is lurking elsewhere.
$ cd /root/astrate-mule && gh issue list --label mule --state open --limit 50 2>&1
(no output)

$ cd /root/astrate-mule && gh issue list --label milestone-4.0 --state all --limit 20 2>&1; echo "---"; gh issue list --label upstream-experimental --state all --limit 20 2>&1
90	CLOSED	milestone-4.0: decide on bumping APICompatVersion from "1.2.2" to "1.3.0-rc.0" (or making it configurable)	enhancement, milestone-4.0	2026-08-23T21:06:34Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
28	CLOSED	milestone-4.0: investigate Edgehog compatibility with Astrate (readonly)	mule, readonly, milestone-4.0	2026-08-23T15:48:58Z
---
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
68	CLOSED	Decide async_operation=false params vs documented always-sync deviation	enhancement, mule-blocked, upstream-parity, upstream-experimental	2026-09-04T18:34:35Z
67	CLOSED	Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4)	enhancement, upstream-parity, upstream-experimental	2026-09-04T18:34:19Z

The picture matches the 2026-09-06 run exactly: no `mule` issues open, `milestone-3.0` label empty, all 43 `upstream-parity` backlog issues closed, the 3 remaining open items (#92, #93, #78) already escalated or in review, and UPSTREAM-EXPERIMENTAL rows #67/#68 still gated on a v1.4.0 final that hasn't shipped (newest stable remains v1.3.3). Let me verify the research source-doc situation the milestone references before concluding.
$ cd /root/astrate-mule && ls -la .mule/research/ 2>&1; echo "---"; git branch -a 2>&1 | head -20; echo "---"; git log --oneline -3 2>&1
ls: cannot access '.mule/research/': No such file or directory
---
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
---
9d822c3 mule: recipe hygiene ran (156s)
17c6cff mule: log docs-sync-appengine-data-422-interface-level
f53d1ee mule: docs-sync-appengine-data-422-interface-level [auto]: the interface-level data GETs — GET /devices/{device}/interfaces/{interface}, GET /devices-by-alias/{alias}/interfaces/{interface}, GET /groups/{group}/devices/{device}/interfaces/{interface} — share parseQueryOpts (internal/appengine/http.go:554-639) and answer 422 FieldErrors for invalid since/since_after/to/limit/downsample_to/format params, but the spec documents only the path-level GET with 422 — add 422 to the interface-level GETs.

$ cd /root/astrate-mule && git ls-tree -r origin/mule/research --name-only 2>&1 | rg 'upstream-parity' ; echo "exit=$?"
exit=1

All checks converge: nothing has moved since the 2026-09-06 milestone run, and this run has nothing to propose or file.

**Investigation:**
- First non-DONE milestone: **v3.0** (upstream 1.2.2 → 1.3/1.4). Scope is decided (not TBD), so no candidate-list for-giulio entry needed.
- `gh issue list --label milestone-3.0 --state all`: **empty** — this milestone's work lives under `upstream-parity`.
- `upstream-parity`: 43/46 closed (all of #47–#89). The 3 open are all already escalated or in review, none a fileable gap:
  - **#92** keyAgreement — escalated, gated on a stable upstream v1.4.0 (still rc, newest stable = v1.3.3).
  - **#93** aclhook comment rewrite — `mule-review`, commit `8c61268` pushed, awaiting review.
  - **#78** FDO — milestone-4.0, escalated.
- **Standing item:** `docs/UPSTREAM-EXPERIMENTAL.md` rows #67/#68 are "1.4 experimental" → reconcile only when upstream 1.4 final ships; the v1.4.0 release check confirms it hasn't. Per #90's frozen decision, the `APICompatVersion` bump is not due.
- Source doc `.mule/research/upstream-parity-2026-08-22.md` is still absent (already escalated 2026-09-06; investigation output #47–#89 intact, so work unaffected).
- No open `mule` issues.

**Actions taken:** none. Filed no issues (no machine-checkable gaps), appended no `.mule/todo.md` lines (nothing proposable), wrote no new for-giulio entry (the 2026-09-06 escalation already states this exact state; the "verify and cut the tag" line is deliberately not proposed because the open decision set — #92, v1.4.0 gates — is real and is Giulio's). No evidence file required by the recipe.

The milestone remains genuinely gated on upstream v1.4.0 final + Giulio's decisions, unchanged from yesterday.
