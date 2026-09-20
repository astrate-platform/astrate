slug: recipe-astarte-upstream
verdict: proposed
at:  1dba3d5
ran: 2026-09-20T11:36:11Z on DietPi in 115s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z

→ Read .mule/todo.md
→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.4") | .body'
### Fixed

- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
  uncollected binaries over time by forcing more frequent full sweep garbage collections
  on them.
- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.

$ gh api repos/astarte-platform/astarte/commits?per_page=20 --jq '.[] | "\(.sha[0:10])\t\(.commit.committer.date)\t\(.commit.message | split("\n")[0])"' | rg -i '1\.3|changelog|release|fix|feat'
257086cb2e	2026-09-18T08:43:47Z	feat(generators): add trigger generator
cd51834ee3	2026-09-18T08:43:47Z	feat(generators): add AMQP trigger target generator
2ef8178815	2026-09-18T08:43:19Z	feat(generators): add simple trigger config generator
461c8a64e2	2026-09-18T08:43:19Z	feat(generators): add device trigger config generator
4f99a33f07	2026-09-18T08:42:35Z	feat(generators): add data trigger config generator
d2caff5c99	2026-09-18T08:42:13Z	feat(generators): add simple trigger scope generator
9372eb0e75	2026-09-17T19:45:02Z	feat(generators): add UUID generator
425dfef0cc	2026-09-17T19:45:02Z	fix(generators): support production environment

$ gh api repos/astarte-platform/astarte/compare/v1.3.3...v1.3.4 --jq '.commits[] | "\(.sha[0:10])\t\(.commit.message | split("\n")[0])"'
7790c741ae	fix(dup): prevent AMQP consumer process to accumulate uncollected bin… (#2194)
37d8fd2ba2	chore!(dup): remove fullsweep_after configuration (#2198)
16befa35f3	fix(dup): ensure RPC server is always available to clients (#2203)
76b21535d9	chore: prepare v1.3.4 (#2208)

$ ls -la .mule/ 2>/dev/null && echo "---reports---" && ls .mule/reports/ 2>/dev/null
total 5476
drwxr-xr-x  7 root root    4096 Sep 19 19:08 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root       0 Sep 19 12:49 .alarmed
-rw-r--r--  1 root root      13 Sep 20 13:34 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 18 22:43 .heartbeat
-rw-r--r--  1 root root       0 Sep 20 13:34 .last-output
-rw-r--r--  1 root root       0 Sep 20 13:34 .lock
-rw-r--r--  1 root root      17 Sep 20 13:34 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 5392234 Sep 20 13:34 cron.log
drwxr-xr-x  2 root root    4096 Sep 19 12:52 failed
-rw-r--r--  1 root root   46506 Sep 20 13:10 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    7742 Sep 19 12:52 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 18 22:00 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   66763 Sep 19 12:52 todo.md
---reports---
hygiene.md
log

→ Read .mule/for-giulio.md

$ rg -n "compat-note|astarte-upstream|v1\.3\.4|v1\.3\.3" .mule/todo.md
61:`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
111:- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
116:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
122:- [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
134:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
173:- [x] compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.

Newest upstream tag is **v1.3.4** (2026-09-18) — same as the last handled release. The `compat-note-v1.3.4` line is already done (`todo.md:173`) and its for-giulio wording proposal is in place (`for-giulio.md:29`); v1.4.0 is still rc-only (rc.5, 2026-08-20), already covered by the pending v1.4-rc note. No release since last time, nothing changed to investigate, so per the recipe I append nothing to the queue.

```
Done: upstream Astarte watch — newest tag v1.3.4 matches the last-handled version; no movement to digest. No task lines appended.
Files: none (nothing changed)
Verified: gh api repos/astarte-platform/astarte/releases (top 5) -> v1.3.4 still newest; compare v1.3.3...v1.3.4 -> 3 commits, all data_updater_plant BEAM/Elixir internals (fullsweep GC, RPC availability), no wire surface; pass
Unsure: nothing
Follow-ups: none — no upstream movement since v1.3.4
```
