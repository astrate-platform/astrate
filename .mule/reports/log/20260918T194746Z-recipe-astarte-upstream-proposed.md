slug: recipe-astarte-upstream
verdict: proposed
at:  1c7d617
ran: 2026-09-18T19:47:46Z on DietPi in 209s


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

The newest upstream release is **v1.3.4** (2026-09-18, today). Let me check what version Astrate currently considers its target and how past runs handled this.
→ Read .mule/todo.md
→ Read docs/UPSTREAM-EXPERIMENTAL.md
Upstream is ahead — newest is **v1.3.4**, while the last tracked stable was v1.3.3. Let me read the release notes and check for-giulio.md for existing wording proposals.
→ Read .mule/for-giulio.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.4") | .body'
### Fixed

- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
  uncollected binaries over time by forcing more frequent full sweep garbage collections
  on them.
- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.

$ rg -n "compat-note|v1\.3\.4|v1\.3\.3" .mule/todo.md
111:- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
116:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
122:- [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
134:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
173:- [x] compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.

The milestone-vs-upstream picture is already recorded, but let me verify what's in the queue and what today's runs did, so I don't duplicate.
$ git log --oneline -20 -- .mule/for-giulio.md .mule/todo.md
79439e6 mule: refill the queue
f6dc460 mule: refill the queue
216e5a3 mule: log compat-note-v134
bcd4383 mule: compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.
194ac6a mule: refill the queue
ff7905b mule: refill the queue
dbbf2e5 mule: refill the queue
bb8bd4e mule: log docs-sync-rm-triggers-422-nested-envelope
4807952 mule: log docs-sync-rm-policies-delete-422
fa326e1 mule: refill the queue
c6ab452 mule: log httpx-cors-vary-origin-passthrough
0af7cfe mule: refill the queue
5a314c7 mule: refill the queue
e33ab46 mule: log docs-sync-native-version-value
30665a5 mule: log docs-sync-pairing-version-value
4c12a38 mule: refill the queue
3f687a0 mule: log observability-health-content-type-test
2de9e46 mule: log observability-dbpool-gauge-coverage
25a3fde mule: log observability-readiness-wedge-test
4534f94 mule: refill the queue

`compat-note-v1.3.4` already exists (todo.md:173, done). Let me check whether today's run already executed this recipe, and confirm the branch state.
$ git branch --show-current; git log -1 --format='%h %ci %s'
mule/queue
1c7d617 2026-09-18 21:42:16 +0200 mule: recipe github-issues ran (119s)

$ git log --format='%h %ci %s' -15
1c7d617 2026-09-18 21:42:16 +0200 mule: recipe github-issues ran (119s)
79439e6 2026-09-18 21:36:42 +0200 mule: refill the queue
b3c4bad 2026-09-18 21:36:42 +0200 mule: recipe milestones ran (176s)
eb95cdd 2026-09-18 20:41:19 +0200 mule: recipe hygiene timed out
4c4cf5b 2026-09-18 20:01:19 +0200 mule: recipe docs-sync ran (421s)
f6dc460 2026-09-18 19:26:36 +0200 mule: refill the queue
42acc40 2026-09-18 19:26:36 +0200 mule: recipe code-review ran (320s)
216e5a3 2026-09-18 18:53:51 +0200 mule: log compat-note-v134
bcd4383 2026-09-18 18:53:51 +0200 mule: compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.
194ac6a 2026-09-18 13:41:01 +0200 mule: refill the queue
f6d958b 2026-09-18 13:41:01 +0200 mule: recipe astarte-upstream ran (222s)
ff7905b 2026-09-18 13:34:42 +0200 mule: refill the queue
f3dd191 2026-09-18 13:34:42 +0200 mule: recipe github-issues ran (206s)
dbbf2e5 2026-09-18 13:00:39 +0200 mule: refill the queue
8190c4e 2026-09-18 13:00:39 +0200 mule: recipe milestones ran (379s)

`compat-note-v1.3.4` already exists (todo.md:173, marked done) and the wording is already in for-giulio.md. Let me confirm the recipe's history — the astarte-upstream recipe already ran today at 13:41 (commit `f6d958b`), proposed it, and the mule executed it at 18:53. Let me verify what the two v1.3.4 changes actually touch upstream to be strict about wire-visible qualification.
$ gh api repos/astarte-platform/astarte/releases/tags/v1.3.4 --jq '.body'; echo '---TAGS---'; gh api 'repos/astarte-platform/astarte/compare/v1.3.3...v1.3.4' --jq '.commits[] | "\(.sha[0:9]) \(.commit.message | split("\n")[0])"'
### Fixed

- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
  uncollected binaries over time by forcing more frequent full sweep garbage collections
  on them.
- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.

---TAGS---
7790c741a fix(dup): prevent AMQP consumer process to accumulate uncollected bin… (#2194)
37d8fd2ba chore!(dup): remove fullsweep_after configuration (#2198)
16befa35f fix(dup): ensure RPC server is always available to clients (#2203)
76b21535d chore: prepare v1.3.4 (#2208)

The recipe already ran today. Let me verify what the 13:41 run proposed, to avoid duplicating.
$ git show 194ac6a --stat | head -20; echo '==='; git show f6dc460 --stat | head -20; echo '=== diff of refill commits for todo.md new lines ==='; git diff 8190c4e 194ac6a -- .mule/todo.md | rg '^\+' | rg -v '^\+\+\+'
commit 194ac6a24428abbe7757c649144b215790d906e6
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Fri Sep 18 13:41:01 2026 +0200

    mule: refill the queue

 .mule/todo.md | 1 +
 1 file changed, 1 insertion(+)
===
commit f6dc46064f6886b38623b805196bebcf9bab37d7
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Fri Sep 18 19:26:36 2026 +0200

    mule: refill the queue

 .mule/for-giulio.md                  |   4 ++
 .mule/reviews/deviceid-2026-09-18.md | 100 +++++++++++++++++++++++++++++++++++
 .mule/todo.md                        |   1 +
 3 files changed, 105 insertions(+)
=== diff of refill commits for todo.md new lines ===
+- [ ] compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.

The recipe already ran today and v1.3.4 is fully handled. Verifying no new movement between the morning run and now — the newest release is still the same v1.3.4 (2026-09-18 09:42Z, before today's 13:41 recipe run).

Both v1.3.4 changes are DUP-internal (`fix(dup): prevent AMQP consumer process to accumulate uncollected binaries` #2194, `chore!(dup): remove fullsweep_after configuration` #2198, `fix(dup): ensure RPC server is always available to clients` #2203) — Erlang process/GC/RPC internals on the RabbitMQ consumer path, which Astrate does not have (in-process sharded pipeline instead). Per the recipe's rules these are Elixir-shaped and ignored. The mandatory `compat-note-v1.3.4` line is already in `.mule/todo.md:173` and marked done; its wording proposal is already in `.mule/for-giulio.md:21`. Nothing new qualifies, so nothing is appended.

```
Done: Upstream newest stable is still v1.3.4 (2026-09-18), already handled by today's 13:41 astarte-upstream run — both 1.3.4 changes are DUP Elixir-internal (AMQP consumer GC, RPC availability), no wire-visible surface, and compat-note-v1.3.4 is already proposed and delivered; no upstream movement since v1.3.4.
Files: none (read-only; nothing appended)
Verified: gh api release sweep -> v1.3.4 latest stable / rc.5; gh compare v1.3.3...v1.3.4 -> 4 commits, all DUP-internal -> pass (no qualifyable items)
Unsure: nothing
Follow-ups: none
```
