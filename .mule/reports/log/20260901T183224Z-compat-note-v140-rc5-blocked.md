slug: compat-note-v140-rc5
verdict: blocked
at:  54b406b
ran: 2026-09-01T18:32:24Z on DietPi in 282s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
I'll start by understanding the task: propose docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in `.mule/for-giulio.md` (without editing the COMPATIBILITY.md file). Let me explore the relevant files.
✱ Glob "docs/COMPATIBILITY.md" 1 match

$ cat .mule/for-giulio.md 2>/dev/null | head -50; echo "---TODO---"; cat .mule/todo.md 2>/dev/null | head -50
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **Milestones recipe has nothing to file: v2.0 is already DONE, and `mule/queue` is 120
  commits behind `origin/main`, so this branch's `.mule/milestones.md` is stale.** My branch
  still shows v2.0 "not started" (and v3.0 as CLEA TBD), so the recipe pointed me at v2.0 —
  but on `main` v2.0 is marked **DONE (2026-07-29)**: every `milestone-2.0` issue is CLOSED
  (#23–#27, #37, #39, #40, #41, #42, #43), and the factory/catalog/container-block/durable
  named flows/`/flow/v1` API all live on `main` (`git show origin/main:.mule/milestones.md`).
  `gh issue list --label milestone-2.0 --state open` is empty and there are no new v2.0 gaps
  to file, so per the recipe's step 5 I filed **nothing** rather than duplicate closed work.
  **Real current target is v3.0** ("upstream Astarte feature parity 1.3/1.4", scoped on main
  2026-08-22): the first recipe job there is to triage `upstream-parity` backlog **#47–#89**
  (source `.mule/research/upstream-parity-2026-08-22.md`) into a plan. But I should not start
  that triage from this stale branch. Action needed: **merge/rebuild `mule/queue` on top of
  `origin/main`**, then the next milestones run targets v3.0 with the correct research doc.
  (Milestones recipe run, 2026-09-01.)

---

- ~~`device_deletion_started`/`device_deletion_finished` trigger events are not emitted~~ —
  **decided 2026-07-27: emit both, back-to-back, around the synchronous delete.** Filed as
  issue #21 (`mule`). (Cross-project survey, 2026-07-27,
  `.mule/research/survey-2026-07-27.md` source 4.)
- ~~Mustache trigger-action templates are accepted but not rendered~~ — **decided
  2026-07-27: implement it.** Guiding principle clarified: Astarte compatibility means
  SDK/wire compatibility, not minimum dependency count — Astrate is allowed to be a
  compatible *superset*. Library picked: `github.com/cbroglie/mustache`. Filed as issue #22
  (`mule`). (Same survey, source 4.)
- **`value_change`/`value_change_applied`/`path_created`/`path_removed`/`value_stored` trigger
  types compile but never fire** (`internal/engine/triggers/match.go:30-42`). Decision
  deferred pending data: issue #20 (`mule`, `readonly`) asks Big Pickle to benchmark the
  ingest-path cost of a previous-value lookup on the Legion Go before this gets decided one
  way or the other. (Same survey, source 4.)
- **Group-scoped triggers (`group_name` on device/data triggers) compile but never match**
  (`internal/engine/triggers/match.go:11-12`). Decision deferred, tied to issue #17
  (group-WATCH-path reconciliation, trickle work, not mule): whatever group-membership
  mechanism comes out of that phase should also report the perf cost for this decision —
  noted in a comment on #17 so it isn't benchmarked twice. (Same survey, source 4.)

---
---TODO---
# Mule queue

One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).
The mule runs the topmost open line. Reorder freely; the order is the priority.

If a line needs more detail than fits, write `.mule/tasks/<slug>.md` and keep the line short.

    tools/mule.sh add "<slug>: <outcome>"
    tools/mule.sh loop

A line tagged `[legion]` needs the Legion Go and is skipped automatically while it is asleep;
the queue moves on to the next runnable line rather than stalling behind it.

A line tagged `[readonly]` is a **standing check**, not a piece of work. It verifies something
and is expected to produce no diff at all — so an empty diff is success, not the `wrote
nothing` failure it would be for any other line. It never gets ticked off, because a gate is
never done. Its output lands in `.mule/reports/<slug>.md` with the sha it ran against, and it
is skipped while that sha is still HEAD: re-checking code that has not moved buys nothing and
spends a call on a free provider.

**`mule/queue` is authoritative for this file.** The mule ticks on the Pi and writes its
progress there, so editing the queue on `main` resurrects completed tasks the moment the two
are merged — which has happened once. Add and reorder tasks on `mule/queue`; let them reach
`main` only when that branch is merged.

A line marked `- [~]` is **parked**: real work, but not tick-sized. The mule only ever
picks up `- [ ]`. Benchmark runs live here because a single tier is 5-20 minutes of ingest
alone -- past the per-task budget -- and because they want someone watching. Run one by hand:

    tools/mule.sh legion bench-push
    ssh legion 'cd ~/astrate/bench && ./scripts/run-tier.sh small astrate -base-url ... -housekeeping-key ...'

## Where tasks come from

**This file is not the whole queue, and for real work it is not even the main part of it.**

The queue is: the standing lines below, plus **every open GitHub issue labelled `mule`**.
Issues are read live on each tick and are never copied into this file — a copy would be a
second place the same fact lives, on a branch the mule commits to and you edit on `main`,
and that produced three merge conflicts in one afternoon.

**Open issues labelled `mule` are worked BEFORE the standing lines below** (changed
2026-08-31). An issue is work a human approved; a standing line is either a gate or something
the mule proposed to itself, and self-proposed work never runs ahead of approved work.

**To give the mule work, file an issue and label it `mule`.** From anywhere, by anyone,
including another model with repo access. No SSH, no editing this file:

    gh issue create --label mule --title "<slug>: <outcome>" --body "<the detail>"

Let me look at the current COMPATIBILITY.md to understand its structure, and check for any prior rc.5 or version-related context.
→ Read docs/COMPATIBILITY.md

$ git log --oneline -20 2>/dev/null; echo "---TAGS---"; git tag 2>/dev/null | tail -20; echo "---VERSION in for-giulio---"; rg -n "1\.4\.0|rc\.5|COMPATIBILITY" .mule/ docs/ --no-heading 2>/dev/null | head -40
54b406b mule: blocked probe-object-validation-selected-interface
4cd1e88 mule: refill the queue
709bc6a mule: recipe astarte-upstream ran (478s)
f4ad882 mule: blocked issue-68-async-operation-accepted
e0f8823 mule: blocked issue-91-pairing-health-serve-root
733f965 mule: refill the queue
b69e32e mule: recipe github-issues ran (167s)
b8424d8 mule: refill the queue
2777546 mule: recipe milestones ran (607s)
3a70ddb mule: refill the queue
22875f1 mule: recipe hygiene timed out
5006f6e mule: blocked docs-sync-pairing-status-enum
49fcd62 mule: refill the queue
45171b9 mule: recipe docs-sync ran (203s)
19ad0fe mule: blocked probe-properties-on-connect-encoding
8a983a3 mule: blocked probe-binaryblob-validation
25560b4 mule: blocked probe-mqtt-capabilities-declaration
21f74d8 mule: blocked compat-note-v133
2f072d3 mule: blocked probe-emptycache-resend-device-error
4c1228c mule: blocked compat-note-v140-rc3
---TAGS---
v0.1.0
---VERSION in for-giulio---
docs/api/astarte_realm_management_api.yaml:500:        deviation from upstream's async deletion (see COMPATIBILITY.md).
docs/ROADMAP.md:428:| 9.10 | `docs/COMPATIBILITY.md` | Deviation inventory (§3.6) + supported-SDK matrix with pinned versions tested |
docs/handoff/phase-8-handoff.md:39:   - Read DESIGN.md §6, OPERATIONS.md, COMPATIBILITY.md
docs/handoff/phase-8-handoff.md:43:   - Read DESIGN.md §0 (comparison), §1.1 (service mapping), COMPATIBILITY.md
docs/DESIGN.md:647:  COMPATIBILITY.md deviation 1.
.mule/milestones.md:10:`.mule/for-giulio.md`, but never edits it directly — same rule as `docs/COMPATIBILITY.md`.
docs/handoff/phase-3-memory.md:32:- `docs/COMPATIBILITY.md` — deviation inventory and tested client matrix used in compatibility page.
docs/handoff/phase-8-memory.md:27:- `docs/COMPATIBILITY.md` — deliberate deviations and supported clients
docs/OPERATIONS.md:98:`docs/COMPATIBILITY.md`). The compose `full` profile ships it pre-wired at
.mule/tasks/realm-config-datastream-retention.md:4:v1.2.2 → v1.4.0-rc.3 gap.
.mule/tasks/realm-config-datastream-retention.md:19:- **v1.4.0-rc.1** added `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`
.mule/tasks/realm-config-datastream-retention.md:38:- `router.ex@v1.3.0` and `@v1.4.0-rc.3` both list the route; `router.ex@v1.2.2` has no
.mule/tasks/realm-config-datastream-retention.md:40:- Release notes v1.3.0 and v1.4.0-rc.1 (both claims now diff-checked).
.mule/todo.md:99:- [!] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto] — BLOCKED: gates failed
.mule/todo.md:100:- [!] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
.mule/todo.md:102:- [!] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
.mule/todo.md:104:- [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
.mule/todo.md:109:- [!] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto] — BLOCKED: gates failed
.mule/todo.md:110:- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]
.mule/recipes/astarte-upstream.md:7:Compatibility state is recorded in `docs/COMPATIBILITY.md` — **read that first**, it tells
.mule/recipes/astarte-upstream.md:17:Compare the newest tag to the version named in `docs/COMPATIBILITY.md`. If they match, say
.mule/recipes/astarte-upstream.md:52:- [ ] compat-note-<tag>: propose the docs/COMPATIBILITY.md wording for <tag> in .mule/for-giulio.md (do not edit the file)
.mule/cron.log:1252:Size S decision.\n\n---\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).
.mule/cron.log:1259:{"comments":[{"id":"IC_kwDORmfsJs8AAAABRidW1w","author":{"login":"atsetilam"},"authorAssociation":"MEMBER","body":"Scoped for the mule (2026-08-31).\n\nDesign is already frozen in \\`.trickle/plans/MASTER-HANDOFF.md\\` fase 4c and is not up for relitigation: Astrate stays always-sync, and \\`async_operation=false\\` is accepted and ignored on the mutating endpoints upstream exposes it on — housekeeping realm create/delete; realm-management interface install/update/delete and trigger/policy delete (policy routes are mounted in \\`internal/realm/http.go\\`). An unparseable or \\`true\\` value must not change behaviour either.\n\n**In scope for the mule:** the code and its tests only.\n**Out of scope:** \\`docs/COMPATIBILITY.md\\` (deviation 17), \\`docs/UPSTREAM-EXPERIMENTAL.md\\`, and closing this issue — the docs closeout stays with the architect and lands once, together with #67.","createdAt":"2026-08-30T23:33:27Z","includesCreatedEdit":false,"isMinimized":false,"minimizedReason":"","reactionGroups":[],"url":"https://github.com/astrate-platform/astrate/issues/68#issuecomment-5471950551","viewerDidAuthor":true}]}
.mule/cron.log:1511:Actually — I bet the real situation is different. Let me reconsider: maybe there IS a real behavior in these handlers that should be gated on async. Upstream's `async_operation=false` is the *opt-in to synchronous*. If Astrate ignores it and is always sync, a client requesting sync gets sync — good. But what about a client that does NOT pass the param (relying on async default)? Upstream would return 202-ish async; Astrate returns the sync result immediately. That's a response-shape difference (upstream async returns 200 with a "status_url"/202?). But design says Astrate stays always-sync everywhere — this is already the accepted deviation (documented as COMPATIBILITY deviation 17, out of scope for mule).
.mule/cron.log:1515:Let me look at the actual git state — is there maybe a `COMPATIBILITY.md` deviation 17 that describes exactly the code contract? Let me check the current repo for it to ground the exact intended behavior.
.mule/cron.log:1517:$ grep -rn "async\|always-sync\|always sync\|deviation 17\|#17" /root/astrate-mule/docs/COMPATIBILITY.md 2>/dev/null | head -30
.mule/cron.log:1520:→ Read docs/COMPATIBILITY.md [offset=195, limit=45]
.mule/cron.log:1521:The COMPATIBILITY.md only has 10 numbered deviations (out of scope for me to edit anyway — but the comment said deviation 17 in a future/other doc; that's out of my scope).
.mule/cron.log:1547:// Astrate is always-sync (docs/COMPATIBILITY.md, async deviation), so callers
.mule/cron.log:1621:// Astrate is always-synchronous (design; see COMPATIBILITY async deviation),
.mule/cron.log:1694:+// design's async deviation, docs handled in COMPATIBILITY), so handlers must
.mule/cron.log:2194:+// design's async deviation, docs handled in COMPATIBILITY), so handlers must
.mule/cron.log:2224:Unsure: The endpoint-level acceptance tests assert behavior that Go already provided implicitly (unknown query params are ignored), so only the astarteapi unit test fails if the helper is removed. No endpoint behaviour actually changes by design — always-sync is preserved. Docs/COMPATIBILITY, docs/UPSTREAM-EXPERIMENTAL, and closing #68 are out of scope (architect/landing with #67).
.mule/cron.log:2862:>> task: probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto]
.mule/cron.log:2869:This is a `[readonly]` probe: compare Astrate's per-mapping value-type validation against upstream Astarte's restored `validate_value_type` from v1.4.0-rc.3. I must report, not patch.
.mule/cron.log:2881:>> task: compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto]
.mule/cron.log:2887:Task: propose wording for `docs/COMPATIBILITY.md` for v1.4.0-rc.3 in `.mule/for-giulio.md` without editing the compat file itself.
.mule/cron.log:2890:→ Read docs/COMPATIBILITY.md
.mule/cron.log:2892:The task: propose the `docs/COMPATIBILITY.md` wording for v1.4.0-rc.3 in `.mule/for-giulio.md` — without editing the compat file itself. So I need to research what changed in v1.4.0-rc.3 and draft the compat doc wording, writing it as a for-giulio entry.

Now I understand the structure of COMPATIBILITY.md. The task is to propose wording for a v1.4.0-rc.5 update in `.mule/for-giulio.md`, without editing the file itself.

Let me understand what version this repo is at, and whether there's an rc.5 that corresponds to something. Let me look at what "v1.4.0-rc.5" refers to — likely the Astrate release version. Let me search for version references and recent context.
$ rg -rn "1\.4\.0|rc\.5" . --no-heading -g '!*.sum' -g '!todo.md' 2>/dev/null | head -40
./go.mod:	github.com/rs/xid vn // indirect
./test/conformance/go.mod:	github.com/rs/xid vn // indirect

Let me look at the recipe `.mule/recipes/astarte-upstream.md` which contains the compat-note workflow, and check the previous compat-note entries and research to understand what rc.5 changed.
→ Read .mule/recipes/astarte-upstream.md
This is a compat-note task. I need to brainstorm the docs/COMPATIBILITY.md wording for v1.4.0-rc.5. Let me look at how previous compat-note tasks were handled and check the research on upstream-parity for v1.4.0-rc.5 changes. Let me look at prior reports and the milestones context.

$ ls .mule/reports/ 2>/dev/null; echo "---tasks---"; ls .mule/tasks/ 2>/dev/null; echo "---research---"; ls .mule/research/ 2>/dev/null
issue-13.md
issue-20.md
log
race-check.md
---tasks---
issue-16.md
realm-config-datastream-retention.md
---research---

Let me look at the upstream-parity research doc referenced in the cron.log and milestones, which covers v1.4.0-rc.5. Let me find it.
$ ls -la .mule/; echo "---"; find . -iname "*upstream*" 2>/dev/null; echo "---git show any research on main---"; git show origin/main:.mule/research/upstream-parity-2026-08-22.md 2>/dev/null | head -80
total 632
drwxr-xr-x  8 root root   4096 Sep  1 20:27 .
drwxr-xr-x 14 root root   4096 Aug 31 21:44 ..
-rw-r--r--  1 root root      0 Jul 28 22:45 .alarmed
-rw-r--r--  1 root root     13 Sep  1 20:27 .budget
-rw-r--r--  1 root root    364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root     11 Jul 28 07:36 .heartbeat
-rw-r--r--  1 root root  16384 Sep  1 20:28 .last-output
-rw-r--r--  1 root root      0 Sep  1 20:27 .lock
-rw-r--r--  1 root root     17 Sep  1 19:49 .rotation
-rw-r--r--  1 root root   6429 Jul 27 22:18 MULE.md
-rw-r--r--  1 root root   2762 Aug 31 01:45 config
-rw-r--r--  1 root root 518054 Sep  1 20:28 cron.log
drwxr-xr-x  2 root root   4096 Sep  1 20:13 failed
-rw-r--r--  1 root root   4901 Sep  1 13:43 for-giulio.md
-rw-r--r--  1 root root   1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root   1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root   4296 Sep  1 20:13 log.md
-rw-r--r--  1 root root   2643 Jul 27 22:18 milestones.md
drwxr-xr-x  2 root root   4096 Jul 27 21:20 plans
drwxr-xr-x  2 root root   4096 Jul 27 22:54 recipes
drwxr-xr-x  3 root root   4096 Jul 28 05:50 reports
drwxr-xr-x  2 root root   4096 Jul 27 14:57 reviews
-rw-r--r--  1 root root    145 Sep  1 20:27 task.md
drwxr-xr-x  2 root root   4096 Aug  3 04:11 tasks
-rw-r--r--  1 root root  12611 Sep  1 20:13 todo.md
---
./test/conformance/upstream
./.mule/recipes/astarte-upstream.md
./.mule/reports/log/20260824T021445Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260803T021148Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260901T175716Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260817T020856Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260810T021151Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md
./.git/refs/remotes/origin/worktree-m12-05-match-upstream
./.git/logs/refs/remotes/origin/worktree-m12-05-match-upstream
---git show any research on main---

Let me look at the most recent upstream recipe report and the milestones/for-giulio to understand what v1.4.0-rc.5 involves. Let me check the cron.log for the previous compat-note attempt for rc.5 and the upstream-parity research.
$ cat .mule/reports/log/20260901T175716Z-recipe-astarte-upstream-proposed.md 2>/dev/null
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_05e3ab9ca001QMNywrSEJjFXS1

081abe4 milestones: redefine v3.0 as upstream Astarte feature parity
83568dd mule: recipe astarte-upstream ran (333s)
3a4b693 mule: recipe astarte-upstream ran (431s)
3a746d1 mule: recipe astarte-upstream ran (485s)
d145c7d mule: propose upstream v1.3.0 wire-visible changes + compat-note-v1.3.2

$ git show 733f965 -- .mule/todo.md | head -40
commit 733f9651fd4a8e4f1e877c7a53e3df81914c7e53
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Tue Sep 1 19:02:01 2026 +0200

    mule: refill the queue

diff --git a/.mule/todo.md b/.mule/todo.md
index bc44d10..f189f2f 100644
--- a/.mule/todo.md
+++ b/.mule/todo.md
@@ -104,3 +104,5 @@ worth having. Before any of it reaches `main`:
 - [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
 - [!] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto] — BLOCKED: gates failed
 - [!] docs-sync-pairing-status-enum: in docs/api/astarte_pairing_api.yaml, fix the PairingInfo.status enum (line 367): it lists `confirmed, pending, denied, expired` but the handler only ever emits `confirmed`, `pending`, or `inhibited` (internal/pairing/service.go:289-297, via internal/pairing/http.go:211) — drop the dead `denied`/`expired` and add the undocumented `inhibited`. [auto] — BLOCKED: gates failed
+- [ ] issue-91-pairing-health-serve-root: in internal/pairing/http.go, add the unauthenticated `GET /pairing/health` route (no realm segment) sharing the handler/payload of the existing `GET /pairing/v1/{realm}/health`, keeping the v1 route as-is; tests for both paths (closes #91)
+- [ ] issue-68-async-operation-accepted: on the mutating endpoints upstream surfaces `async_operation` on (housekeeping realm create/delete; realm-management interface install/update/delete, trigger/policy delete in internal/realm/http.go), accept and ignore `?async_operation=false`, with unparseable/`true` values also not changing behaviour (Astrate stays always-sync); tests cover the flag parsing (closes #68)

$ git show 7b03846 --stat | head -20
commit 7b0384603d4f9515a3e89a457002a4aef3761c8b
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 31 04:12:36 2026 +0200

    mule: recipe astarte-upstream ran (224s)

 ...831T021236Z-recipe-astarte-upstream-proposed.md | 280 +++++++++++++++++++++
 1 file changed, 280 insertions(+)

$ git show 7b03846 -- .mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md | head -120
commit 7b0384603d4f9515a3e89a457002a4aef3761c8b
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 31 04:12:36 2026 +0200

    mule: recipe astarte-upstream ran (224s)

diff --git a/.mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md b/.mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md
new file mode 100644
index 0000000..b3a6d53
--- /dev/null
+++ b/.mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md
@@ -0,0 +1,280 @@
+slug: recipe-astarte-upstream
+verdict: proposed
+at:  6e30a1c
+ran: 2026-08-31T02:12:36Z on DietPi in 224s
+
+
+> build · big-pickle
+
+I'll start by reading the two files specified.
+→ Read .mule/MULE.md
+→ Read .mule/recipes/astarte-upstream.md
+Now let me read the compatibility doc and check upstream releases in parallel.
+→ Read docs/COMPATIBILITY.md
+
+$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
+v1.4.0-rc.5	2026-08-20T12:28:04Z
+v1.3.3	2026-08-07T14:33:58Z
+v1.4.0-rc.4	2026-08-06T13:40:20Z
+v1.4.0-rc.3	2026-07-31T10:48:39Z
+v1.4.0-rc.2	2026-07-14T10:16:22Z
+
+→ Read .mule/todo.md
+Upstream is ahead: newest tag is `v1.4.0-rc.5` (stable also has `v1.3.3`), while COMPATIBILITY.md records upstream v1.2.0/v1.2.2. Let me check the todo queue for prior runs and pull the release notes.
+
+$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
+Sixth Astarte v1.4.0 release candidate.
+
+## Changes since v1.4.0-rc.4
+
+### Fixed
+
+- [astarte_data_updater_plant] Make object values only validate type on selected interface
+- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated
+
+$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'
+The queue already covers v1.3.0/v1.3.3/v1.4.0-rc.3 and rc.5's binaryblob fix, but not rc.5's other change: "Make object values only validate type on selected interface." Let me look at that upstream change before proposing anything.
+
+$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | {tag_name, body, published_at}'
+{"body":"","published_at":"2026-08-07T14:33:58Z","tag_name":"v1.3.3"}
+
+$ gh search commits --owner astarte-platform 'Make object values only validate type on selected interface' --json sha,repository,commit --jq '.[] | "\(.repository.full_name)\t\(.sha[0:9])\t\(.commit.message | split("\n")[0])"'
+(no output)
+
+$ gh search commits --owner astarte-platform 'object values only validate type on selected interface' --json sha,repository,commit 2>/dev/null; echo "---"; gh search commits --owner astarte-platform 'object values' --json sha,repository,commit.url,commit.author --jq '.[] | "\(.repository.full_name)\t\(.sha[0:9])\t\(.commit.message | split("\n")[0])"' 2>&1 | head -20
+[]
+---
+Unknown JSON field: "commit.url"
+Available fields:
+  author
+  commit
+  committer
+  id
+  parents
+  repository
+  sha
+  url
+
+$ gh api repos/astarte-platform/astarte_data_updater_plant/commits --jq '.[0:15][] | "\(.sha[0:9])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"'
+8739d915f	2019-12-20T15:06:00Z	This repository has been archived.
+69612e95f	2019-12-18T16:49:09Z	Forward port changes from v0.10
+28ae84960	2019-12-09T11:32:32Z	Merge branch 'prepare-0.10.2' into 'release-0.10'
+c0087ddce	2019-12-09T11:21:25Z	Prepare 0.10.2 release
+64da6f368	2019-11-28T14:49:05Z	Merge pull request #17 from rbino/scylla-beta-support
+a40acb406	2019-11-28T14:18:26Z	.travis.yml: remove Scylla CI from allowed failures
+9a3198a26	2019-11-28T14:12:40Z	queries: don't use shorthand syntax to update maps
+bc11f0383	2019-11-27T15:16:08Z	Merge branch 'ci-upgrade-services' into 'master'
+59bd736d7	2019-11-27T10:49:58Z	.gitlab-ci.yml: upgrade Cassandra and RabbitMQ services
+a4156ce2b	2019-11-26T14:20:37Z	Forward port bugfixes from v0.10
+140f80393	2019-11-26T11:34:43Z	Merge branch 'fix-13' into 'release-0.10'
+745d9b792	2019-11-26T11:22:01Z	impl: load device volatile triggers when they're installed
+afe9d5b4f	2019-11-26T11:21:20Z	Merge branch 'fix-15' into 'release-0.10'
+466ca7c70	2019-11-26T10:49:12Z	impl: fix volatile triggers deletion bug
+6fcd3f5f0	2019-11-19T16:24:14Z	Merge pull request #12 from rbino/ci-wait-for-cassandra
+
+The data updater code now lives in the astarte monorepo. Let me find the actual commit for the rc.5 object-validation change.
+$ gh api 'repos/astarte-platform/astarte/commits?since=2026-08-15T00:00:00Z&until=2026-08-21T00:00:00Z' --jq '.[] | "\(.sha[0:9])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"'
+739724991	2026-08-19T13:06:29Z	chore: Update exandra
+643d06107	2026-08-19T12:36:15Z	ci: add release-please for automated versioning and changelog
+6e65c8aca	2026-08-17T12:25:10Z	refactor(dup): Move encryption logic into library (#2123)
+88004a2cf	2026-08-14T09:21:28Z	chore(core): use the local astarte_core across the monorepo
+
+$ gh search commits --owner astarte-platform --repository astarte 'binaryblob data is correctly validated' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo "==="; gh search commits --owner astarte-platform --repository astarte 'object values only validate type' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'
+unknown flag: --repository
+
+Usage:  gh search commits [<query>] [flags]
+
+Flags:
+      --author string            Filter by author
+      --author-date date         Filter based on authored date
+      --author-email string      Filter on author email
+      --author-name string       Filter on author name
+      --committer string         Filter by committer
+      --committer-date date      Filter based on committed date
+      --committer-email string   Filter on committer email
+      --committer-name string    Filter on committer name
+      --hash string              Filter by commit hash
+  -q, --jq expression            Filter JSON output using a jq expression
+      --json fields              Output JSON with the specified fields
+  -L, --limit int                Maximum number of commits to fetch (default 30)
+      --merge                    Filter on merge commits
+      --order string             Order of commits returned, ignored unless '--sort' flag is specified: {asc|desc} (default "desc")
+      --owner strings            Filter on repository owner
+      --parent string            Filter by parent hash
+  -R, --repo OWNER/REPO          Filter on repository, in OWNER/REPO format
+      --sort string              Sort fetched commits: {author-date|committer-date} (default "best-match")
+  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
+      --tree string              Filter by tree hash
+      --visibility strings       Filter based on repository visibility: {public|private|internal}

$ git show 7b03846 -- .mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md | sed -n '120,280p'
+      --visibility strings       Filter based on repository visibility: {public|private|internal}
+  -w, --web                      Open the search query in the web browser
+  
+===
+unknown flag: --repository
+
+Usage:  gh search commits [<query>] [flags]
+
+Flags:
+      --author string            Filter by author
+      --author-date date         Filter based on authored date
+      --author-email string      Filter on author email
+      --author-name string       Filter on author name
+      --committer string         Filter by committer
+      --committer-date date      Filter based on committed date
+      --committer-email string   Filter on committer email
+      --committer-name string    Filter on committer name
+      --hash string              Filter by commit hash
+  -q, --jq expression            Filter JSON output using a jq expression
+      --json fields              Output JSON with the specified fields
+  -L, --limit int                Maximum number of commits to fetch (default 30)
+      --merge                    Filter on merge commits
+      --order string             Order of commits returned, ignored unless '--sort' flag is specified: {asc|desc} (default "desc")
+      --owner strings            Filter on repository owner
+      --parent string            Filter by parent hash
+  -R, --repo OWNER/REPO          Filter on repository, in OWNER/REPO format
+      --sort string              Sort fetched commits: {author-date|committer-date} (default "best-match")
+  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
+      --tree string              Filter by tree hash
+      --visibility strings       Filter based on repository visibility: {public|private|internal}
+  -w, --web                      Open the search query in the web browser
+
+$ gh search commits -R astarte-platform/astarte 'binaryblob data is correctly validated' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo "==="; gh search commits -R astarte-platform/astarte 'object values only validate type on selected interface' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'
+aedc41483	fix(dup): ensure binaryblob data is correctly validated
+===
+
+$ gh search commits -R astarte-platform/astarte 'object values only validate type' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo "===loose==="; gh search commits -R astarte-platform/astarte 'object values' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"' | head -10
+===loose===
+cc89858b3	DUP: do not crash on binaryblob insertion
+219b9c79d	Merge branch 'object-aggregation' into 'master'
+084db3757	Handle aggregate interface in live events
+adf94edf8	Fix data parsing on datastream object interfaces w/ parametric endpoints
+
+$ gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.4...v1.4.0-rc.5 --jq '.commits[] | "\(.sha[0:9])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"'
+d8592091e	2026-08-07T06:55:11Z	fix(dup): stop remapping interface loading failed error (#2119)
+a816eaf0e	2026-08-07T06:57:11Z	chore: forward port release-1.3
+aedc41483	2026-08-07T07:21:33Z	fix(dup): ensure binaryblob data is correctly validated
+3e70d63f8	2026-08-07T14:28:36Z	chore: prepare v1.3.3 (#2128)
+fdc6d29b8	2026-08-19T08:14:18Z	fix(dup): Make object values only validate type on selected interface (#2141)
+4141a2ae4	2026-08-19T13:08:32Z	chore: forward-port release-1.3 into release-1.4
+3484096d2	2026-08-20T12:10:25Z	Prepare release v1.4.0-rc.5 (#2147)
+
+The new item is commit `fdc6d29b8` (#2141). Let me read the diff.
+$ gh api repos/astarte-platform/astarte/commits/fdc6d29b8 --jq '{sha, message: .commit.message, files: [.files[] | {filename, status: .status, additions, deletions}]}'
+{"files":[{"additions":1,"deletions":1,"filename":"apps/astarte_appengine_api/mix.lock","status":"modified"},{"additions":2,"deletions":1,"filename":"apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex","status":"modified"},{"additions":40,"deletions":0,"filename":"apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_test.exs","status":"modified"},{"additions":23,"deletions":0,"filename":"apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/interface_test.exs","status":"modified"},{"additions":58,"deletions":0,"filename":"apps/astarte_data_updater_plant/test/support/cases/device.ex","status":"modified"},{"additions":1,"deletions":1,"filename":"apps/astarte_pairing/mix.lock","status":"modified"}],"message":"fix(dup): Make object values only validate type on selected interface (#2141)\n\nSigned-off-by: Francesco Noacco \u003cfrancesco.noacco@secomind.com\u003e","sha":"fdc6d29b84ec06a12316340b60f1b43491ca3d5c"}
+
+$ gh api repos/astarte-platform/astarte/commits/fdc6d29b8 --jq '.files[] | select(.filename=="apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex") | .patch'
+@@ -316,11 +316,12 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.Interface do
+   end
+ 
+   def extract_mappings(
+-        %InterfaceDescriptor{aggregation: :object},
++        %InterfaceDescriptor{aggregation: :object, interface_id: interface_id},
+         _mapping,
+         mappings
+       ) do
+     mappings
++    |> Enum.filter(fn {_id, m} -> m.interface_id == interface_id end)
+     |> Map.new(fn {_id, m} ->
+       key = m.endpoint |> String.split("/") |> List.last()
+       {key, m}
+
+$ gh api repos/astarte-platform/astarte/commits/fdc6d29b8 --jq '.files[] | select(.filename=="apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/interface_test.exs") | .patch'
+@@ -209,6 +209,29 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.InterfaceTest do
+       end
+     end
+ 
++    @tag regression: true
++    test "extract_mappings/3 only considers mappings for the current object interface", context do
++      %{
++        fixed_object_datastream_1: fixed_object_datastream_1,
++        fixed_object_datastream_2: fixed_object_datastream_2,
++        state: state,
++        realm_name: realm_name
++      } = context
++
++      # ignored for objects
++      mapping = nil
++
++      descriptor = state.interfaces[fixed_object_datastream_1.name]
++
++      assert %{"first" => %{value_type: :string}, "second" => %{value_type: :integer}} =
++               Core.Interface.extract_mappings(descriptor, mapping, state.mappings)
++
++      descriptor = state.interfaces[fixed_object_datastream_2.name]
++
++      assert %{"first" => %{value_type: :boolean}, "second" => %{value_type: :binaryblob}} =
++               Core.Interface.extract_mappings(descriptor, mapping, state.mappings)
++    end
++
+     property "forget_interfaces/2 removes interfaces from state cache", context do
+       %{
+         interfaces: interfaces,
+
+✱ Grep "validateValueType|ValidateValueType|func.*[Vv]alidate.*[Vv]alue" in internal/engine · 1 match
+This is a distinct wire-visible fix: object datastream values were validated against whatever interface had the same endpoint names, since `extract_mappings` for an `:object` aggregation pulled mappings across all interfaces. Let me see where Astrate validates object values so the probe line targets the right files.
+✱ Grep "extract_mappings|extractMappings" in internal · 0 matches
+✱ Grep "[Vv]alue[Tt]ype" in . · 89 matches
+✱ Grep "Lookup|Trie|resolveMapping|Mapping.*by.*path|matchPath" in internal/engine · 13 matches
+Let me find how Astrate resolves the mapping set when decoding an object value — that determines whether the upstream bug has an Astrate mirror.
+→ Read pkg/payload/payload.go
+→ Read internal/engine/data.go [offset=230, limit=120]
+Astrate decodes object values against `ci.ObjectLeaves` — the topic-named interface's own leaves — so the rc.5 #2141 bug may or may not have a mirror; that's exactly what a probe verifies. Let me check rc.4's notes too, since it sits between what prior runs covered (rc.3, rc.5-partial), and confirm nothing else is pending.
+
+$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.4") | .body'
+Fifth Astarte v1.4.0 release candidate.
+
+## Changes since v1.4.0-rc.3
+
+### Added
+
+- [config] Initialize `config` library. All external services declared using the `config` library helpers can be configured either using the `_URL` environment variable or by using single component variables `_SCHEME`, `_HOST`, `_PORT`, `_PATH`, `_QUERY`, `_FRAGMENT` (`_URL` has priority), together with normal `_SSL_*` variables.
+- [secrets] Use `config` for vault configuration. Vault can now be configured using component variables instead of `ASTARTE_VAULT_URL`:
+  - `ASTARTE_VAULT_SCHEME`: scheme part of the URL (binary, default: `"http"`, or `"https"` if `ASTARTE_VAULT_SSL_ENABLED` is `true`)
+  - `ASTARTE_VAULT_HOST`: host part of the URL (binary, default: `"localhost"`)
+  - `ASTARTE_VAULT_PORT`: port part of the url (integer, default `8200`)
+  - `ASTARTE_VAULT_PATH`: path part of the url (string)
+  - `ASTARTE_VAULT_QUERY`: query part of the url (string)
+  - `ASTARTE_VAULT_FRAGMENT`: fragment part of the url (string)
+- [fdo] Use `config` for rendezvous configuration. The rendezvous can now be configured using component variables instead of `PAIRING_FDO_RENDEZVOUS_URL`:
+  - `PAIRING_FDO_RENDEZVOUS_SCHEME`: scheme part of the URL (binary, default: `"http"`, or `"https"` if `PAIRING_FDO_RENDEZVOUS_SSL_ENABLED` is `true`)
+  - `PAIRING_FDO_RENDEZVOUS_HOST`: host part of the URL (binary, default: `"localhost"`)
+  - `PAIRING_FDO_RENDEZVOUS_PORT`: port part of the url (integer, default `8041`)
+  - `PAIRING_FDO_RENDEZVOUS_PATH`: path part of the url (string)
+  - `PAIRING_FDO_RENDEZVOUS_QUERY`: query part of the url (string)
+  - `PAIRING_FDO_RENDEZVOUS_FRAGMENT`: fragment part of the url (string)
+
+### Fixed
+
+- [astarte_data_updater_plant] Ensure RPC server is always available to clients. Resolved the issue where a temporary disconnection and reconnection of data_updater_plant to the cluster would make the RPC server inaccessible.
+
+### Changed
+
+- [astarte_data_updater_plant] Ensure memory is properly garbage collected
+
+← Edit .mule/todo.md
+Index: /root/astrate-mule/.mule/todo.md
+===================================================================
+--- /root/astrate-mule/.mule/todo.md
++++ /root/astrate-mule/.mule/todo.md
+@@ -102,4 +102,6 @@
+ - [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]
+ - [ ] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto]
+ - [ ] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto]
+ - [ ] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto]
++- [ ] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto]
++- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]
+
+Nothing new in rc.4 (config-library and Elixir-infra changes — ignore). The only uncovered item since the last run is rc.5's #2141. Appending the two task lines to the queue.
+$ tail -n 5 .mule/todo.md
+- [ ] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto]

$ git branch --show-current; git log --all --oneline -1 -- .mule/todo.md; git log --oneline -3 -- .mule/todo.md
mule/queue
f4ad882 mule: blocked issue-68-async-operation-accepted
f4ad882 mule: blocked issue-68-async-operation-accepted
e0f8823 mule: blocked issue-91-pairing-health-serve-root
733f965 mule: refill the queue

The last upstream run (2026-08-31) proposed `probe-object-validation-selected-interface` and `compat-note-v1.4.0-rc.5`, but those aren't in the current todo.md on disk. Let me check how the queue state evolved — the file may have been overwritten by refills.
$ git log --all --oneline -S "probe-object-validation-selected-interface" -- .mule/todo.md
(no output)

$ git log --all --oneline --reverse -- .mule/todo.md
956d2a7 Add solo-mule mode: an unsupervised opencode worker driven by a task queue
0921b60 mule: make the Legion Go the concurrency gate
2e885ab mule: note mule/queue as authoritative for todo.md, make race-check recurring
8b333f2 mule: a standing check that finds nothing is passing, not failing
0b5b23f mule: refill its own queue, with GitHub issues as the front door
6f2a33e mule: refill the queue
8d98ab5 mule: issues are the queue, and a review step before main
95fdef5 mule: merge main
71877c3 mule: drop the copied issue line — issues are read live now
8606336 mule: refill the queue
16cc008 mule: log store-realm-cascade-policies
344d013 mule: log store-alias-lowest-id
ad335ec mule: log store-delete-device-objects
8a18476 mule: blocked race-check
4f6d95c mule: split race-check into 5 per-subsystem lines
d145c7d mule: propose upstream v1.3.0 wire-visible changes + compat-note-v1.3.2
28f3adf mule: fix ticks dying outside the repo (cron cwd), recover unsaved queue state
6e30a1c queue: drop three lines already done or duplicated, approve the plaintext-purge line, take the issues-first ordering from main
eeb69d4 mule: blocked control-producer-properties-compression
db8048e mule: blocked probe-interface-default-values
acfcb26 mule: blocked probe-value-type-validation
4c1228c mule: blocked compat-note-v140-rc3
2f072d3 mule: blocked probe-emptycache-resend-device-error
21f74d8 mule: blocked compat-note-v133
25560b4 mule: blocked probe-mqtt-capabilities-declaration
8a983a3 mule: blocked probe-binaryblob-validation
19ad0fe mule: blocked probe-properties-on-connect-encoding
49fcd62 mule: refill the queue
5006f6e mule: blocked docs-sync-pairing-status-enum
733f965 mule: refill the queue
e0f8823 mule: blocked issue-91-pairing-health-serve-root
f4ad882 mule: blocked issue-68-async-operation-accepted

$ git log --all --oneline -S "compat-note-v1.4.0-rc.5"
7b03846 mule: recipe astarte-upstream ran (224s)

$ git show d145c7d --stat | head -20; echo ===; git show d145c7d -- .mule/todo.md | head -40
commit d145c7d98cf19891485a0913928ff02e1db3585a
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Wed Jul 29 00:45:07 2026 +0200

    mule: propose upstream v1.3.0 wire-visible changes + compat-note-v1.3.2

 .mule/todo.md | 5 +++++
 1 file changed, 5 insertions(+)
===
commit d145c7d98cf19891485a0913928ff02e1db3585a
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Wed Jul 29 00:45:07 2026 +0200

    mule: propose upstream v1.3.0 wire-visible changes + compat-note-v1.3.2

diff --git a/.mule/todo.md b/.mule/todo.md
index af5ed5f..32b8fd6 100644
--- a/.mule/todo.md
+++ b/.mule/todo.md
@@ -83,3 +83,8 @@ worth having. Before any of it reaches `main`:
 - [~] bench-big-astarte [legion]: run `bench/scripts/run-tier.sh big astarte` against Legion Go Astarte, commit results (two runs minimum)
 - [~] bench-giant-astarte [legion]: run `bench/scripts/run-tier.sh giant astarte` against Legion Go Astarte, commit results (two runs minimum)
 - [ ] race-check: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...` (~40s). Report any failure to .mule/for-giulio.md with the full race report. This is the only race coverage that exists — the Pi cannot run -race. [legion] [readonly]
+
+- [ ] purge-properties-compression-capability [auto]: upstream v1.3.0 adds a `purge_properties_compression_format` device capability (`zlib`|`plaintext`, default `zlib`) — a wire-visible capability value. Check whether Astrate's capabilities handling (internal/broker, the `<realm>/<device_id>/capabilities` topic, issue #16) needs to recognize/honour it, or whether zlib-only is already the deliberate default; propose the change or note why not needed.
+- [ ] empty-introspection-verification [auto]: upstream v1.3.0 changed "allow devices with empty introspection" — verify whether Astrate's device connection/introspection handling currently rejects an empty introspection string where upstream now accepts it, and propose a fix if so.
+- [ ] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages" — investigate only: does Astrate have an analogous delay between trigger install and first delivery? Report, do not patch.
+- [ ] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.

$ git log --all --format='%h %ad %s' --date=iso -- .mule/todo.md | head -40
f4ad882 2026-09-01 19:46:55 +0200 mule: blocked issue-68-async-operation-accepted
e0f8823 2026-09-01 19:32:34 +0200 mule: blocked issue-91-pairing-health-serve-root
733f965 2026-09-01 19:02:01 +0200 mule: refill the queue
5006f6e 2026-08-31 22:32:39 +0200 mule: blocked docs-sync-pairing-status-enum
49fcd62 2026-08-31 22:10:36 +0200 mule: refill the queue
19ad0fe 2026-08-31 21:45:56 +0200 mule: blocked probe-properties-on-connect-encoding
8a983a3 2026-08-31 21:35:46 +0200 mule: blocked probe-binaryblob-validation
25560b4 2026-08-31 21:20:19 +0200 mule: blocked probe-mqtt-capabilities-declaration
21f74d8 2026-08-31 21:05:59 +0200 mule: blocked compat-note-v133
2f072d3 2026-08-31 20:47:54 +0200 mule: blocked probe-emptycache-resend-device-error
4c1228c 2026-08-31 20:42:21 +0200 mule: blocked compat-note-v140-rc3
acfcb26 2026-08-31 20:26:28 +0200 mule: blocked probe-value-type-validation
db8048e 2026-08-31 19:21:40 +0200 mule: blocked probe-interface-default-values
eeb69d4 2026-08-31 18:52:25 +0200 mule: blocked control-producer-properties-compression
6e30a1c 2026-08-31 02:10:57 +0200 queue: drop three lines already done or duplicated, approve the plaintext-purge line, take the issues-first ordering from main
28f3adf 2026-08-22 02:26:26 +0200 mule: fix ticks dying outside the repo (cron cwd), recover unsaved queue state
d145c7d 2026-07-29 00:45:07 +0200 mule: propose upstream v1.3.0 wire-visible changes + compat-note-v1.3.2
4f6d95c 2026-07-28 23:25:22 +0200 mule: split race-check into 5 per-subsystem lines
8a18476 2026-07-28 00:46:30 +0200 mule: blocked race-check
ad335ec 2026-07-27 17:44:03 +0200 mule: log store-delete-device-objects
344d013 2026-07-27 16:36:55 +0200 mule: log store-alias-lowest-id
16cc008 2026-07-27 15:31:24 +0200 mule: log store-realm-cascade-policies
8606336 2026-07-27 14:58:49 +0200 mule: refill the queue
71877c3 2026-07-27 11:44:22 +0200 mule: drop the copied issue line — issues are read live now
95fdef5 2026-07-27 11:44:22 +0200 mule: merge main
8d98ab5 2026-07-27 11:43:49 +0200 mule: issues are the queue, and a review step before main
6f2a33e 2026-07-27 11:09:20 +0200 mule: refill the queue
0b5b23f 2026-07-27 11:07:06 +0200 mule: refill its own queue, with GitHub issues as the front door
8b333f2 2026-07-27 10:49:27 +0200 mule: a standing check that finds nothing is passing, not failing
2e885ab 2026-07-27 06:12:54 +0200 mule: note mule/queue as authoritative for todo.md, make race-check recurring
0921b60 2026-07-27 06:07:17 +0200 mule: make the Legion Go the concurrency gate
956d2a7 2026-07-27 05:26:54 +0200 Add solo-mule mode: an unsupervised opencode worker driven by a task queue

Let me check what's in `.mule/for-giulio.md` (where compat-note proposals land) and research the v1.4.0-rc.0 "required flag for object-aggregated mappings" change, which is the one wire-visible schema item not yet covered.
→ Read .mule/for-giulio.md

$ gh search commits -R astarte-platform/astarte 'required flag for mappings of object' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo ===; gh search commits -R astarte-platform/astarte 'required flag' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.author.date)\t\(.commit.message|split("\n")[0])"' | head -10
===
82018e1c9	2026-05-20T10:33:00+02:00	feat: add required flag support to interface mappings (#524)
b9a4418d1	2026-05-20T10:33:00+02:00	feat: add required flag support to interface mappings (#524)

Let me examine the upstream `required` flag change (#524) to understand its wire effect before proposing.
$ gh api repos/astarte-platform/astarte/commits/82018e1c9 --jq '{message: .commit.message, files: [.files[] | {filename, status, additions, deletions}]}'
{"files":[{"additions":20,"deletions":2,"filename":"astarte-dashboard/cypress/e2e/interface_builder.cy.js","status":"modified"},{"additions":2,"deletions":1,"filename":"astarte-dashboard/cypress/fixtures/test.astarte.AggregatedObjectInterface.json","status":"modified"},{"additions":6,"deletions":0,"filename":"astarte-dashboard/src/astarte-client/models/Mapping/index.ts","status":"modified"},{"additions":2,"deletions":0,"filename":"astarte-dashboard/src/astarte-client/transforms/mapping.ts","status":"modified"},{"additions":1,"deletions":0,"filename":"astarte-dashboard/src/astarte-client/types/dto/mapping.d.ts","status":"modified"},{"additions":6,"deletions":0,"filename":"astarte-dashboard/src/components/InterfaceEditor.tsx","status":"modified"},{"additions":21,"deletions":0,"filename":"astarte-dashboard/src/components/MappingEditor.tsx","status":"modified"}],"message":"feat: add required flag support to interface mappings (#524)\n\nSigned-off-by: Riccardo Nalgi \u003criccardo.nalgi@secomind.com\u003e"}

$ gh api repos/astarte-platform/astarte/commits/b9a4418d1 --jq '{message: .commit.message, files: [.files[] | {filename, status, additions, deletions}]}'
{"files":[{"additions":20,"deletions":2,"filename":"astarte-dashboard/cypress/e2e/interface_builder.cy.js","status":"modified"},{"additions":2,"deletions":1,"filename":"astarte-dashboard/cypress/fixtures/test.astarte.AggregatedObjectInterface.json","status":"modified"},{"additions":6,"deletions":0,"filename":"astarte-dashboard/src/astarte-client/models/Mapping/index.ts","status":"modified"},{"additions":2,"deletions":0,"filename":"astarte-dashboard/src/astarte-client/transforms/mapping.ts","status":"modified"},{"additions":1,"deletions":0,"filename":"astarte-dashboard/src/astarte-client/types/dto/mapping.d.ts","status":"modified"},{"additions":6,"deletions":0,"filename":"astarte-dashboard/src/components/InterfaceEditor.tsx","status":"modified"},{"additions":21,"deletions":0,"filename":"astarte-dashboard/src/components/MappingEditor.tsx","status":"modified"}],"message":"feat: add required flag support to interface mappings (#524)\n\nSigned-off-by: Riccardo Nalgi \u003criccardo.nalgi@secomind.com\u003e"}

The `required` flag commit is actually a dashboard-UI change (#524), not a backend schema change. Let me check the older recipe-run logs to confirm what v1.4.0-rc.0 items were already considered and rejected, so I don't re-propose dropped items.
$ git show 83568dd --stat | head; echo ===; git show 3a4b693 --stat | head -5; echo ===; git show 3a746d1 --stat | head -5
commit 83568ddfa0a4ebdac26755099fbf9e38ddb0197e
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 17 04:08:56 2026 +0200

    mule: recipe astarte-upstream ran (333s)

 ...817T020856Z-recipe-astarte-upstream-proposed.md | 405 +++++++++++++++++++++
 1 file changed, 405 insertions(+)
===
commit 3a4b69370d6c405411fa7307c7cac3b0b7565812
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 10 04:11:51 2026 +0200

    mule: recipe astarte-upstream ran (431s)
===
commit 3a746d1c64d902f284321c9ca480f02c4ccfe4cb
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 3 04:11:49 2026 +0200

    mule: recipe astarte-upstream ran (485s)

$ git show 83568dd:.mule/reports/log/20260817T020856Z-recipe-astarte-upstream-proposed.md | rg -n "propose|todo|detailed|required|rc\.|v1\.4|Append|Edit" | head -40
2:verdict: proposed
15:v1.4.0-rc.4	2026-08-06T13:40:20Z
16:v1.4.0-rc.3	2026-07-31T10:48:39Z
17:v1.4.0-rc.2	2026-07-14T10:16:22Z
18:v1.4.0-rc.1	2026-07-14T08:53:27Z
58:v1.3.0-rc.2	v1.3.0-rc.2	Third Astarte v1.3.0 release candidate.
60:## Changes since v1.3.0-rc.1
71:v1.3.0-rc.1	v1.3.0-rc.1	Second Astarte v1.3.0 release candidate.
73:# Changes since v1.3.0-rc.0
80:v1.3.0-rc.0	v1.3.0-rc.0	First Astarte v1.3.0 release candidate.
102:## [1.4.0-rc.4] - 2026-08-06
130:## [1.4.0-rc.3] - 2026-07-31
143:## [1.4.0-rc.2] - 2026-07-14
145:## [1.4.0-rc.1] - 2026-07-13
163:## [1.4.0-rc.0] - 2026-04-08
167:- Add required flag for mappings of object aggregated interfaces
176:- [astarte_realm_management] Allow listing interfaces with mappings using the `detailed=true` parameter
221:## [1.3.0-rc.2] - 2026-04-08
232:  - `ASTARTE_BASE_URL_DOMAIN` - domain part of the base URL of astarte, used by devices to connect in TO2 phase (required if FDO enabled)
233:  - `ASTARTE_BASE_URL_PORT` - port of the base URL of astarte (required if FDO enabled)
234:  - `ASTARTE_BASE_URL_PROTOCOL` - protocol of the base URL of astarte (required if FDO enabled)
250:## [1.3.0-rc.1] - 2026-01-26
269:## [1.3.0-rc.0] - 2025-11-21
276:  - `HOUSEKEEPING_ASTARTE_KEYSPACE_NETWORK_REPLICATION_MAP` - Datacenter replication map when using NetworkTopologyStrategy (no default, required when using network strategy)
291:- BREAKING: [astarte_housekeeping] Housekeeping now creates an AMQP vhost for each created realm. Some required configuration was introduced:
292:  - `HOUSEKEEPING_AMQP_HOST` - AMQP host for housekeeping operations (required in production, defaults to "localhost" in dev/test)
297:## [1.3.0-rc.0] - 2025-11-21
304:  - `HOUSEKEEPING_ASTARTE_KEYSPACE_NETWORK_REPLICATION_MAP` - Datacenter replication map when using NetworkTopologyStrategy (no default, required when using network strategy)
319:- BREAKING: [astarte_housekeeping] Housekeeping now creates an AMQP vhost for each created realm. Some required configuration was introduced:
320:  - `HOUSEKEEPING_AMQP_HOST` - AMQP host for housekeeping operations (required in production, defaults to "localhost" in dev/test)
350:- [astarte_realm_management] Allow to delete long-disconnected devices. For this to work, Realm Management needs `CLUSTERING_STRATEGY`, `CLUSTERING_KUBERNETES_NAMESPACE` and `DATA_UPDATER_PLANT_CLUSTERING_KUBERNETES_SELECTOR` to be set, just like AppEngine and DUP. Refer to 1.2.1-rc.0 for additional information on the variables.
355:→ Read .mule/todo.md
360:$ grep -n 'device.deletion.*trigger\|deletion.*trigger\|probe-emptycache\|compat-note-v1.3.3\|compat-note-v1.3' /root/astrate-mule/.mule/todo.md
362:95:- [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]
365:→ Read .mule/todo.md [offset=90]
367:← Edit .mule/todo.md
368:Index: /root/astrate-mule/.mule/todo.md
370:--- /root/astrate-mule/.mule/todo.md
371:+++ /root/astrate-mule/.mule/todo.md
373: - [ ] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto]

$ git show 83568dd:.mule/reports/log/20260817T020856Z-recipe-astarte-upstream-proposed.md | sed -n '150,215p'
### Fixed

- [fdo] Do not crash when using rotated keys as owner keys
- [fdo] Delete devices from unfinished sessions
- [fdo] Delete old sessions from device
- [fdo] Properly respond to service info messages after IsDone=True
- [astarte_realm_management] Properly show default values for all mapping parameters

### Changed

- Enable FDO by default

## [1.4.0-rc.0] - 2026-04-08

### Added

- Add required flag for mappings of object aggregated interfaces
- [astarte_pairing] Support for FIDO Device Onboard v1.1 for device registration. The following environment variables have been added for the Vault connectivity:
  - `ASTARTE_VAULT_URL` - The URL to access Vault.
  - `ASTARTE_VAULT_AUTHENTICATION_MECHANISM` - The mechanism to use for authenticating with Vault (default: `""`, allowed: \["token"])
  - `ASTARTE_VAULT_TOKEN` - Token to authenticate with Vault, used if `ASTARTE_VAULT_AUTHENTICATION_MECHANISM` is `"token"` (default: `""`)
  - `ASTARTE_VAULT_SSL_CA_FILE` - path to the CA certificate file for the vault TLS connection; when not specified, the bundled cURL certificate bundle will be used
  - `ASTARTE_VAULT_SSL_CUSTOM_SNI` - custom SNI hostname for the vault TLS connection; when not specified, the vault hostname will be used
  - `ASTARTE_VAULT_SSL_DISABLE_SNI` - Disable Server Name Indication for Vault. (default: `false`)
  - `ASTARTE_VAULT_SSL_ENABLED` - Whether to enable SSL for the Vault connection. (default `false`)
- [astarte_realm_management] Allow listing interfaces with mappings using the `detailed=true` parameter
- [astarte_import] Added support for data types: `doublearray`, `integerarray`,
  `booleanarray`, `longintegerarray`, `stringarray`, `datetimearray`, `binaryblobarray`.

### Changed

- [astarte_data_updater_plant] Use mississippi consumer for data updater processes

## [1.3.3] - Unreleased

### Fixed

- [astarte_data_updater_plant] Resending properties to a device now reports distinct device error triggers depending on the failure: `interface_loading_failed` when an interface fails to load and `resend_interface_properties_failed` when sending the properties to the device itself fails. Previously, both cases were incorrectly reported as `resend_interface_properties_failed`.

## [1.3.2] - 2026-07-14

## Fixed

- Avoid crash on invalid properties message

## [1.3.1] - 2026-07-13

### Fixed

- Solve capabilities message stalling queues
- [astarte_housekeeping] Allow resuming partially completed realm deletions
- Use network topology as default realm replication
- Fix simple strategy realm creation without replication factor

### Changed

- Immediately delete devices in device deletion best case scenario

## [1.3.0] - 2026-05-06

### Changed

- [astarte_housekeeping] When `HOUSEKEEPING_ASTARTE_KEYSPACE_REPLICATION_STRATEGY` is not set, the `astarte` keyspace is now created using `NetworkTopologyStrategy` with a replication map derived from the current ScyllaDB network topology (one replica per node in each datacenter), instead of falling back to `SimpleStrategy` with replication factor 1.
- [astarte_housekeeping] When creating a realm without specifying its replication strategy, the strategy used in the `astarte` keyspace is applied as default.

$ git show 83568dd:.mule/reports/log/20260817T020856Z-recipe-astarte-upstream-proposed.md | sed -n '355,405p'
→ Read .mule/todo.md
• Check device triggers Astrate Explore Agent
✓ Check device triggers Astrate Explore Agent
Now I have the full picture. Let me check what's already tracked vs. what's new.

$ grep -n 'device.deletion.*trigger\|deletion.*trigger\|probe-emptycache\|compat-note-v1.3.3\|compat-note-v1.3' /root/astrate-mule/.mule/todo.md
94:- [ ] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto]
95:- [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]

Upstream is at v1.3.3, Astrate targets v1.2.x. Most v1.3.x changes are already tracked. The one wire-visible gap not yet in the queue is device deletion trigger events. Let me append the new task lines.
→ Read .mule/todo.md [offset=90]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -92,4 +92,5 @@
 - [ ] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto]
 - [ ] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto]
 - [ ] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto]
 - [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]
+- [ ] trigger-device-deletion-events: add `device_deletion_started` and `device_deletion_finished` trigger conditions to `internal/engine/triggers/match.go` (new constants + entries in `deviceOns`), emit them from the device-deletion path in `internal/engine/engine.go`, and add them to the `on` enum in `docs/api/astarte_realm_management_api.yaml` (upstream v1.3.0-rc.0) [auto]

Here's the summary:

**Upstream: v1.3.3 (released 2026-08-07) vs. Astrate target: v1.2.x**

Most v1.3.x changes are already tracked in the queue:
- `control-producer-properties-compression` — plaintext producer/properties (v1.3.0 capability)
- `realm-config-datastream-retention` — per-realm data retention (v1.3.0)
- `probe-emptycache-resend-device-error` — v1.3.3 error differentiation
- `compat-note-v1.3.3` / `compat-note-v1.4.0-rc.3` — compat notes

**New task added:**
- `trigger-device-deletion-events` — upstream v1.3.0-rc.0 added `device_deletion_started` and `device_deletion_finished` trigger conditions. Astrate has no equivalent. This is wire-visible (Realm Management trigger API `on` enum, device lifecycle events).

**Already implemented (no task needed):**
- Empty introspection accepted — already handled in `introspection.go:95-102`
- Capabilities + `purge_properties_compression_format` — fully implemented in `capabilities.go`
- Device registration triggers — already wired in `engine.go:283-301`
- Transient trigger crash fix — N/A (Astrate's transient triggers are matcher-only, never enter pipeline)

**Infrastructure-only / Elixir-only (ignored per recipe rules):**
- AMQP vhost per realm, merged API containers, env var changes, ScyllaDB replication strategy, capabilities queue stalling fix, volatile trigger encoding fix — none of these are wire-visible or applicable to Astrate's architecture.

Done: appended 1 task line (trigger-device-deletion-events) to .mule/todo.md
Files: .mule/todo.md
Verified: grep confirms no duplicate
Unsure: nothing
Follow-ups: trigger-device-deletion-events

Confirmed the picture now. Summary of findings:

- **Upstream movement:** newest tag is still `v1.4.0-rc.5` (2026-08-20); stable is `v1.3.3` (2026-08-07). No new releases since the last upstream run (2026-08-31).
- **The last run's queue append was lost:** commit `7b03846` logged its proposal (`probe-object-validation-selected-interface` for rc.5 #2141, plus `compat-note-v1.4.0-rc.5`) but only the log file landed — the todo.md edit never made it into the branch. Neither line exists in the current queue; the only rc.5 item tracked is the binaryblob fix.
- **Everything else in the v1.4.0 RCs** (config-library, vault, FDO, required-flag dashboard UI #524, housekeeping retention) was reviewed by prior runs as dashboard/infra-only and correctly not proposed.

So I'll restore the two lost rc.5 lines (not duplicates — they're absent from the queue) and nothing else.
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -105,4 +105,6 @@
 - [!] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto] — BLOCKED: gates failed
 - [!] docs-sync-pairing-status-enum: in docs/api/astarte_pairing_api.yaml, fix the PairingInfo.status enum (line 367): it lists `confirmed, pending, denied, expired` but the handler only ever emits `confirmed`, `pending`, or `inhibited` (internal/pairing/service.go:289-297, via internal/pairing/http.go:211) — drop the dead `denied`/`expired` and add the undocumented `inhibited`. [auto] — BLOCKED: gates failed
 - [!] issue-91-pairing-health-serve-root: in internal/pairing/http.go, add the unauthenticated `GET /pairing/health` route (no realm segment) sharing the handler/payload of the existing `GET /pairing/v1/{realm}/health`, keeping the v1 route as-is; tests for both paths (closes #91) — BLOCKED: gates failed
 - [!] issue-68-async-operation-accepted: on the mutating endpoints upstream surfaces `async_operation` on (housekeeping realm create/delete; realm-management interface install/update/delete, trigger/policy delete in internal/realm/http.go), accept and ignore `?async_operation=false`, with unparseable/`true` values also not changing behaviour (Astrate stays always-sync); tests cover the flag parsing (closes #68) — BLOCKED: gates failed
+- [ ] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto]
+- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]

$ grep -n "probe-object-validation-selected-interface\|compat-note-v1.4.0-rc.5" .mule/todo.md
109:- [ ] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto]
110:- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]

Done: upstream is at v1.4.0-rc.5 (stable v1.3.3) — no new releases since the last run; restored the two rc.5 proposal lines the previous run appended but lost from the queue (probe-object-validation-selected-interface, compat-note-v1.4.0-rc.5). Nothing new proposed; remaining v1.4.0-rc.x items are config/vault/FDO infra or the dashboard-frontend `required`-flag UI (#524), all correctly skipped.

Files: .mule/todo.md
Verified: `gh api .../releases` + `git log -S` on todo.md (confirm lost append) -> pass; grep confirms no duplicate of the two lines -> pass
Unsure: whether restoring lines whose prior append never committed counts as duplication — no live duplicates exist, so I restored rather than leave the gap
Follow-ups: probe-object-validation-selected-interface, compat-note-v1.4.0-rc.5 (both `[auto]`)

Now I have the full context. Let me understand the two rc.5 changes:

1. **`#2141` "Make object values only validate type on selected interface"** — wire-visible: object datastream values were validated against whatever interface had the same endpoint names; now only the selected (topic-named) interface's mappings are used.
2. **`aedc41483` "Ensure binaryblob data is correctly validated"** — already covered by `probe-binaryblob-validation`.

The COMPATIBILITY.md currently records upstream `v1.2.0`/`v1.2.2` (via Channels measurement references) and the emulated version `1.2.2`. Let me check the current header and how existing compat-note proposals looked, so I match the style. Let me look at the for-giulio.md header/format more and how past compat-note proposals were written (e.g., the rc.3 one that was blocked).

Let me look at how the for-giulio entries are formatted and check the current COMPATIBILITY.md's version references. The task is to *propose* the wording, writing it into `.mule/for-giulio.md`, without editing the COMPATIBILITY.md file itself.

Let me check the current for-giulio.md top entries to match style, and confirm where to add (newest at top).
$ rg -n "version|v1\.2|v1\.4|upstream" docs/COMPATIBILITY.md | head; echo "---for-giulio head (first 12 lines)---"; sed -n '1,12p' .mule/for-giulio.md
6:versions the conformance suite (`test/conformance`) pins and exercises.
14:| Client | Pinned version | Checkpoint(s) |
47:- The Astarte Dashboard v1.2.2 runs unmodified against Astrate (compose
55:1. **Astarte Channels: two sockets, one bus** — the upstream Phoenix socket is
74:     Recorded against upstream v1.2.0 in
75:     `test/conformance/upstream/channels.json`. Upstream partitions the `a_ch`
83:     which stays a verb-regex match because that is what upstream's REST plug
88:     (`test/conformance/upstream/channels.json`). A data trigger is authorized
93:     it refused claims upstream accepts and accepted claims upstream refuses.
95:     `groups/<name>`) come from the same upstream function but are **not
---for-giulio head (first 12 lines)---
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

Now I understand the structure. The COMPATIBILITY.md records upstream against v1.2.0/v1.2.2 (the channels and emulated API version references). The task asks me to propose wording for v1.4.0-rc.5 in `.mule/for-giulio.md`.

The two rc.5 changes are both DUP (data updater plant) validation fixes. The compat-note is about what, if anything, the COMPATIBILITY.md needs to say about v1.4.0-rc.5. Let me consider the substance:

The rc.5 changes are:
1. **#2141** — object values only validate type on selected interface. This is already covered by `probe-object-validation-selected-interface` (a `[readonly]` probe), which is the investigation line saying Astrate decodes against `ci.ObjectLeaves` so it may not have the bug.
2. **binaryblob validation** — already covered by `probe-binaryblob-validation`.

The compat-note's job is to propose the COMPATIBILITY.md wording. But here's the key insight: these rc.5 changes are **bugfixes to validation behaviour** whose Astrate equivalence is *unverified* — exactly the kind of thing the probes cover. The COMPATIBILITY.md records *deliberate deviations* and *supported tested clients*. An rc.5 candidate's bugfix changes wouldn't normally update the deviations inventory unless there's a new deviation.

However, the task is explicit: "propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md". The wording to propose should reflect the state after the probes run. Since both rc.5 changes are validation-behaviour fixes whose Astrate mirrors are still under probe, the honest proposal is: the COMPATIBILITY.md wording for rc.5 is contingent on the probes' outcome — it cannot be written as fact until those complete.

Let me check whether COMPATIBILITY.md already discusses value-type validation anywhere, so my proposed wording slots in correctly.

Let me look at how a blocked prior compat-note attempt (rc.3 or v1.3.3) in the cron.log was worded, to match the deliverable format. Let me check the cron.log for those attempts.
$ rg -n "compat-note-v1.4.0-rc.3|compat-note-v1.3.3|COMPATIBILITY.md wording|propose the docs/COMPATIBILITY" .mule/cron.log | head -20
2881:>> task: compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto]
3015:100:- [ ] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto]
3016:102:- [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]
3020:The task is: "compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto]"
3163:- [ ] compat-note-<tag>: propose the docs/COMPATIBILITY.md wording for <tag> in .mule/for-giulio.md (do not edit the file)
3196:+- [ ] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
3205:100:- [ ] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto]
3207:102:- [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]
3212:Now I have full context. My task `compat-note-v1.4.0-rc.3` is to propose the `docs/COMPATIBILITY.md` wording for upstream Astarte v1.4.0-rc.3 in `.mule/for-giulio.md`.
3382:Let me check the for-giulio as it stands and whether prior compat-note tasks completed by writing proposed wording there. Since the queue still has `compat-note-v1.3.3` and `compat-note-v1.4.0-rc.3` open, and no v1.3.x wording is in for-giulio, it seems these compat-note tasks haven't been resolved yet. I'm the one doing the v1.4.0-rc.3 one.
3384:My deliverable: **propose the docs/COMPATIBILITY.md wording** for v1.4.0-rc.3 in `.mule/for-giulio.md`. I should write a for-giulio entry containing the actual proposed prose for the COMPATIBILITY doc (a docs page that needs Giulio's voice — this is exactly the "docs page that needs your voice" escalation category).
3415:Wait — that's a *mule change proposal*, not what my task asks. My task is narrowly: **propose the COMPATIBILITY.md wording**, not fix code. But the wording needs to be grounded. Let me be disciplined about scope.
3417:The deliverable is a for-giulio entry proposing the COMPATIBILITY.md wording for v1.4.0-rc.3. The COMPATIBILITY.md doc needs to reflect, for v1.4.0-rc.3:
3430:Now let me write the for-giulio entry. I'll propose the actual COMPATIBILITY.md wording. Since COMPATIBILITY.md doesn't have a single "targets upstream version" line beyond the emulated-level note and measured-surfaces references, my proposed wording should focus on what a v1.4.0-rc.3 note would add.
3440:I'll insert my proposal at the top, right after line 11's `---` separator. The entry proposes the COMPATIBILITY.md wording for v1.4.0-rc.3.
3450:The proposed COMPATIBILITY.md wording: a short new "Target" / version-reference note plus a marked gap. Since the doc's structure is deviation-inventory + client matrix, the natural insertion is to (a) refresh the emulated-version and any "measured against" references where upstream moved, and (b) add a deviation/note entry for the `required`-flag gap and the validation-restore.
3514:Done: Proposed the docs/COMPATIBILITY.md wording for upstream v1.4.0-rc.3 in .mule/for-giulio.md (did not touch the doc itself), grounded in the v1.4.0-rc.0/rc.3 release notes and a code check of the `required`-mapping gap.
3604:>> task: compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]
3757:The task is `compat-note-v1.3.3`: propose the `docs/COMPATIBILITY.md` wording for v1.3.3 in `.mule/for-giulio.md`. The version v1.3.3 refers to an *upstream Astarte* release (the probes mention "upstream v1.3.3 (#2119)"). This is about documenting what changed in the compatibility document for the v1.3.3 conformance target.
3759:Let me understand what "v1.3.3" means here. Looking at the conformance pinned versions, and the precedent task `compat-note-v1.4.0-rc.3` which was blocked. Let me understand what these compat-note tasks actually do — I need to find the upstream v1.3.3 release changes relevant to compatibility.

Let me look at the rc.3 compat-note for-giulio entry in the cron.log to match the exact intended deliverable style, and see whether it still exists in for-giulio.md.
$ sed -n '3440,3514p' .mule/cron.log
I'll insert my proposal at the top, right after line 11's `---` separator. The entry proposes the COMPATIBILITY.md wording for v1.4.0-rc.3.

Let me draft the wording carefully, grounded in what I verified:
- COMPATIBILITY.md currently emulates `realm.APICompatVersion` = `1.2.2` and documents deviations measured against upstream v1.2.x.
- v1.4.0-rc.3 upstream added (cumulative since v1.3.x, the wire-visible ones):
  - `required` flag on object-aggregated interface mappings (rc.0). Astrate's `ParseInterface` uses `DisallowUnknownFields` (`pkg/interfaceschema/parse.go:121`) and has no `required` field (`pkg/interfaceschema/types.go` mappingJSON, `parse.go:100-112`), so such an interface is currently **rejected** — a v1.4.0 wire gap.
  - `astarte_data_updater_plant` "Restore correct type validation checks for incoming data" (rc.3) — the `probe-value-type-validation` item (queue line 99, blocked).
  - Realm-management `detailed=true` interface-list listing (rc.0) — a REST query-param addition, additive.
  - The rest (FDO/Vault, mississippi consumer, housekeeping default-retention env, astarte_import array types) are Elixir/infra/config — not wire-visible to unmodified SDKs.

The proposed COMPATIBILITY.md wording: a short new "Target" / version-reference note plus a marked gap. Since the doc's structure is deviation-inventory + client matrix, the natural insertion is to (a) refresh the emulated-version and any "measured against" references where upstream moved, and (b) add a deviation/note entry for the `required`-flag gap and the validation-restore.

I'll propose it as a concise block. Let me write the edit.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,38 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Proposed `docs/COMPATIBILITY.md` wording for upstream v1.4.0-rc.3** (keep/adapt/approve —
+  `docs/COMPATIBILITY.md` is on the never-touch list, so I propose, not edit). Upstream v1.4.0
+  series (rc.0–rc.3, released 2026-04-08..07-31) vs the v1.2.x the doc currently targets; the
+  only wire-visible deltas that touch what Astrate implements are below. Suggested insertion
+  under `Wire-identical surfaces` / `Deliberate deviations`:
+
+  > **Targets upstream v1.4.0-rc.3** (ahead of the v1.3.3 stable the previous note pinned).
+  > Two v1.4.0 items are wire-visible to Astrate's surfaces:
+  >
+  > - **`required` on object-aggregated interface mappings** (added v1.4.0-rc.0). Astrate's
+  >   `ParseInterface` decodes with `json.DisallowUnknownFields` and has no `required` mapping
+  >   field (`pkg/interfaceschema/parse.go:100-112,121`), so an interface whose object mapping
+  >   carries `"required": true` is rejected, where upstream v1.4.0 accepts it. Not yet covered
+  >   — a deviation, listed here as a gap until `probe-value-type-validation`/the derived task
+  >   lands.
+  > - **Ingest value-type validation restored** (rc.3, `astarte_data_updater_plant`). Matches
+  >   Astrate's per-mapping validation; parity claim stays as-is until the open
+  >   `probe-value-type-validation` records it.
+  >
+  > Everything else in the v1.4.0 notes (FDO/Vault provisioning, mississippi AMQP consumer,
+  > `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`, `astarte_import` array types,
+  > Realm-Management `detailed=true` interface listing) is infra/config or additive-REST, with
+  > no SDK impact; Astrate neither needs to nor intends to emulate them.
+
+  (Evidence: `gh api repos/astarte-platform/astarte/releases/tags/v1.4.0-rc.{0,3}` bodies;
+  verified the `required` field is absent and unknown fields are rejected in
+  `pkg/interfaceschema/parse.go`.)
+
+---
+
 - ~~`device_deletion_started`/`device_deletion_finished` trigger events are not emitted~~ —
   **decided 2026-07-27: emit both, back-to-back, around the synchronous delete.** Filed as
   issue #21 (`mule`). (Cross-project survey, 2026-07-27,
   `.mule/research/survey-2026-07-27.md` source 4.)

Now let me verify the gate. This is a docs-only proposal in `.mule/for-giulio.md` — no code changed. Per MULE.md, for a non-code task like this, the gate is essentially that nothing broke. Let me run `gofmt -l .` and `go vet ./... && go test ./...` to satisfy the gate (they should pass since I only touched a markdown file).

Given my context budget, let me run gofmt -l quickly (fast) and go build-vet on the touched package. Actually MULE.md says run the gate. Since I edited only a markdown file, the Go gate is unaffected. Let me run gofmt -l . and a quick vet/test to confirm.
$ cd /root/astrate-mule && gofmt -l . 2>/dev/null; echo "gofmt-done"; go build ./... 2>&1 | tail -5; echo "build-done"
gofmt-done
build-done

gofmt reports nothing (clean), and `go build ./...` succeeds. I edited only a markdown file, so no code is affected. The gate passes.

Let me verify my for-giulio edit reads correctly.
→ Read .mule/for-giulio.md [offset=10, limit=30]
The entry is in place at the top of the queue.

## Report

Done: Proposed the docs/COMPATIBILITY.md wording for upstream v1.4.0-rc.3 in .mule/for-giulio.md (did not touch the doc itself), grounded in the v1.4.0-rc.0/rc.3 release notes and a code check of the `required`-mapping gap.

I have the full context. The rc.3 compat-note's for-giulio entry didn't survive (that task was eventually blocked/reverted), so mine is the entry that lands. Let me verify the key claim I'll cite — that Astrate decodes object datastreams against the selected interface's own leaves — so the proposed wording is grounded.
$ rg -n "ObjectLeaves|objectLeaves" internal/engine/data.go | head; echo "---"; rg -rn "ObjectLeaves" internal/ pkg/ --no-heading | head -10; echo "---payload map object---"; rg -n "func.*[Mm]ap|Object" pkg/payload/payload.go | head
269:		dp, err = e.dec.Object(m.Payload, ci.ObjectLeaves)
331:	for leaf := range ci.ObjectLeaves {
---
internal/engine/data.go:		dp, err = e.dec.Object(m.Payload, ci.n)
internal/engine/data.go:	for leaf := range ci.n {
pkg/interfaceschema/compile.go:	// n maps, for object aggregation, each last-level name to its
pkg/interfaceschema/compile.go:	n map[string]*CompiledMapping
pkg/interfaceschema/compile.go:		ci.n = make(map[string]*CompiledMapping, len(iface.Mappings))
pkg/interfaceschema/compile.go:		if ci.n != nil {
pkg/interfaceschema/compile.go:			if _, dup := ci.n[leaf]; dup {
pkg/interfaceschema/compile.go:			ci.n[leaf] = cm
pkg/interfaceschema/compile_test.go:	if ci.n != nil {
pkg/interfaceschema/compile_test.go:		t.Errorf("n = %v, want nil for individual aggregation", ci.n)
---payload map object---
13:// through Decoder.Individual / Decoder.Object (or the package-level
53:func (d Decoder) Individual(p []byte, m *interfaceschema.CompiledMapping) (DecodedPayload, error) {
94:// Object decodes p against an object-aggregated interface: `v` must be a
96:// (CompiledInterface.ObjectLeaves). Object aggregation exists only on
100:func (d Decoder) Object(p []byte, leaves map[string]*interfaceschema.CompiledMapping) (DecodedPayload, error) {
121:		if val, err = decodeBSONObject(rv, leaves); err != nil {
130:		if val, err = decodeJSONObject(raw, leaves); err != nil {
153:func objectExplicitTimestamp(leaves map[string]*interfaceschema.CompiledMapping) bool {
163:func Decode(p []byte, m *interfaceschema.CompiledMapping) (DecodedPayload, error) {
167:// DecodeObject decodes an object-aggregation payload with the default

I now have what I need. My deliverable: write the proposed `docs/COMPATIBILITY.md` wording for upstream v1.4.0-rc.5 into `.mule/for-giulio.md` (the doc itself stays untouched). The rc.5 wire-visible deltas are both DUP validation fixes that are already under probe (`probe-object-validation-selected-interface`, `probe-binaryblob-validation`), plus the forward-ported #2119 covered by `probe-emptycache-resend-device-error`. Let me write the proposal at the top of the for-giulio queue.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,50 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Proposed `docs/COMPATIBILITY.md` wording for upstream v1.4.0-rc.5** (keep/adapt/approve —
+  `docs/COMPATIBILITY.md` is on the never-touch list, so I propose, not edit). Upstream is at
+  v1.4.0-rc.5 (2026-08-20; stable v1.3.3, 2026-08-07), ahead of the v1.2.x/1.2.2 the doc
+  currently targets (§10 emulated API version; §1 measurements). The only wire-visible deltas
+  the rc series adds over v1.3.x that touch Astrate's surfaces are validation-behaviour fixes,
+  both still under probe, so the wording can't be written as fact until they report. Suggested
+  insertion under `Wire-identical surfaces` / `Deliberate deviations`:
+
+  > **Targets upstream v1.4.0-rc.5** (ahead of the v1.3.3 stable). Three wire-visible items,
+  > all returned-to-correctness fixes to upstream's own validation whose parity is being
+  > verified rather than claimed:
+  >
+  > - **Object values validate type on the selected interface only** (rc.5, #2141) — upstream
+  >   used to type-check an object datastream against any interface sharing endpoint names;
+  >   `extract_mappings` now keeps only the topic-named interface's mappings. Astrate already
+  >   decodes against the topic-named interface's own leaves
+  >   (`internal/engine/data.go:269`, `ci.ObjectLeaves`), so the cross-interface bug likely has
+  >   no mirror — but parity is recorded by `probe-object-validation-selected-interface`
+  >   (todo line 109), not asserted here.
+  > - **binaryblob values correctly validated** (rc.5, `aedc41483`) — upstream's validator now
+  >   brackets binaryblob acceptance where previously it did not. Astrate's boundaries are
+  >   measured by `probe-binaryblob-validation` (todo line 104, against
+  >   `pkg/payload/value.go` and `internal/engine/serverdata.go`); §1's
+  >   "value_size_exceeded unverified" note is superseded by whatever that records.
+  > - **No remap of `interface_loading_failed` on property resend** (rc.5, forward-ported
+  >   v1.3.3 #2119) — already tracked by `probe-emptycache-resend-device-error` (todo line 94)
+  >   against the §1 measured `device_error.error_name` mapping.
+  >
+  > Everything else in the v1.4.0 notes (rc.4 `config` library and vault/rendezvous env
+  > variables, RPC-server availability, GC work; rc.0-rc.3 FDO/Vault, mississippi AMQP
+  > consumer, housekeeping retention, `astarte_import` array types, Realm-Management
+  > `detailed=true` interface listing, dashboard `required`-flag UI #524) is infra/config,
+  > dashboard-frontend, or additive-REST with no unmodified-SDK impact; Astrate neither needs
+  > to nor intends to emulate them.
+
+  (Evidence: `gh api repos/astarte-platform/astarte/releases/tags/v1.4.0-rc.5` body + compare
+  `v1.4.0-rc.4...v1.4.0-rc.5`: `fdc6d29b8` #2141, `aedc41483` binaryblob, `d8592091e` #2119
+  forward-port; Astrate's selected-interface decode verified at `internal/engine/data.go:269`
+  via `ci.ObjectLeaves`.)
+
+---
+
 - **Milestones recipe has nothing to file: v2.0 is already DONE, and `mule/queue` is 120
   commits behind `origin/main`, so this branch's `.mule/milestones.md` is stale.** My branch
   still shows v2.0 "not started" (and v3.0 as CLEA TBD), so the recipe pointed me at v2.0 —
   but on `main` v2.0 is marked **DONE (2026-07-29)**: every `milestone-2.0` issue is CLOSED

Now let me run the gate since I edited a file in the repo. Per MULE.md step 5, I run `gofmt -l .` and `go test ./...` (markdown-only change, so Go is unaffected, but the gate is the habit).
$ gofmt -l . 2>/dev/null; echo "gofmt-exit:$?"
gofmt-exit:0

Markdown-only change in `.mule/`; `gofmt` is clean. Go code untouched, but let me sanity-check the build.
$ go build ./... 2>&1 | tail -3; echo "build-exit:${PIPESTATUS[0]}"
/bin/sh: 1: Bad substitution

$ go build ./... 2>&1 | tail -3; echo "build-done"
build-done

Let me verify the entry is correctly placed at the top of the queue.
→ Read .mule/for-giulio.md [limit=45]
The entry landed. I cited todo.md line numbers (109, 104, 94) — the queue drifts, so let me verify they're still accurate before leaving them in.

$ rg -n "probe-object-validation-selected-interface|probe-binaryblob-validation|probe-emptycache-resend-device-error|compat-note-v1.4.0-rc.5" .mule/todo.md
101:- [!] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto] — BLOCKED: gates failed
104:- [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
109:- [!] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto] — BLOCKED: gates failed
110:- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]

Line numbers drift in the queue; I cited stale ones. Let me make the references line-number-free so they stay accurate.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -26,10 +26,10 @@
 >   used to type-check an object datastream against any interface sharing endpoint names;
 >   `extract_mappings` now keeps only the topic-named interface's mappings. Astrate already
 >   decodes against the topic-named interface's own leaves
 >   (`internal/engine/data.go:269`, `ci.ObjectLeaves`), so the cross-interface bug likely has
->   no mirror — but parity is recorded by `probe-object-validation-selected-interface`
->   (todo line 109), not asserted here.
+>   no mirror — but parity is recorded by `probe-object-validation-selected-interface`, not
+>   asserted here.
 > - **binaryblob values correctly validated** (rc.5, `aedc41483`) — upstream's validator now
 >   brackets binaryblob acceptance where previously it did not. Astrate's boundaries are
 >   measured by `probe-binaryblob-validation` (todo line 104, against
 >   `pkg/payload/value.go` and `internal/engine/serverdata.go`); §1's

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -30,14 +30,14 @@
 >   no mirror — but parity is recorded by `probe-object-validation-selected-interface`, not
 >   asserted here.
 > - **binaryblob values correctly validated** (rc.5, `aedc41483`) — upstream's validator now
 >   brackets binaryblob acceptance where previously it did not. Astrate's boundaries are
->   measured by `probe-binaryblob-validation` (todo line 104, against
+>   measured by `probe-binaryblob-validation` (against
 >   `pkg/payload/value.go` and `internal/engine/serverdata.go`); §1's
 >   "value_size_exceeded unverified" note is superseded by whatever that records.
 > - **No remap of `interface_loading_failed` on property resend** (rc.5, forward-ported
->   v1.3.3 #2119) — already tracked by `probe-emptycache-resend-device-error` (todo line 94)
->   against the §1 measured `device_error.error_name` mapping.
+>   v1.3.3 #2119) — already tracked by `probe-emptycache-resend-device-error` against the §1
+>   measured `device_error.error_name` mapping.
 >
 > Everything else in the v1.4.0 notes (rc.4 `config` library and vault/rendezvous env
 > variables, RPC-server availability, GC work; rc.0-rc.3 FDO/Vault, mississippi AMQP
 > consumer, housekeeping retention, `astarte_import` array types, Realm-Management

Now the entry's stable. Let me also check whether this compat-note should be ticked off in the queue — the task line was `- [ ]`, but the mule script handles committing/ticking. Per MULE.md I leave the working tree; the script commits if gates pass. I should not tick anything myself.

Final verification of the entry and a quick gate check.
$ gofmt -l . 2>/dev/null; echo "gofmt-clean"; go vet ./... 2>&1 | tail -3; echo "vet-done"
gofmt-clean
vet-done

→ Read .mule/for-giulio.md [limit=46]
The proposal is in place at the top of `.mule/for-giulio.md`, `docs/COMPATIBILITY.md` untouched. I did not touch git.

```
Done: Proposed docs/COMPATIBILITY.md wording for upstream v1.4.0-rc.5 in .mule/for-giulio.md — rc.5's only wire-visible deltas are three validation fixes already under probe (object-on-selected-interface #2141, binaryblob aedc41483, #2119 forward-port); wording is deferred to the probes, everything else in the rc's is infra/config and skipped
Files: .mule/for-giulio.md
Verified: gofmt -l . -> pass (nothing); go vet ./... -> pass; gopark go build ./... -> pass (md-only change, no tests run)
Unsure: whether to assert rc.5 parity as fact — chose not to, since the two validation changes are still queued probes; the entry flags this for Giulio to approve the deferred wording
Follow-ups: none
```
