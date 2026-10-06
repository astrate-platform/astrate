slug: recipe-astarte-upstream
verdict: proposed
at:  59a31a2
ran: 2026-10-06T17:43:53Z on DietPi in 1064s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.3.5	2026-10-05T14:11:44Z
v1.4.0-rc.6	2026-09-30T07:28:27Z
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z

$ ls -la /root/astrate-mule/.mule/
total 13220
drwxr-xr-x  7 root root     4096 Oct  6 16:05 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root        0 Oct  6 12:59 .alarmed
-rw-r--r--  1 root root       13 Oct  6 19:26 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  5 20:14 .heartbeat
-rw-r--r--  1 root root        0 Oct  6 19:26 .last-output
-rw-r--r--  1 root root        0 Oct  6 19:26 .lock
-rw-r--r--  1 root root       17 Oct  6 19:26 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 13133888 Oct  6 19:26 cron.log
drwxr-xr-x  2 root root     4096 Oct  4 22:35 failed
-rw-r--r--  1 root root   145756 Oct  6 18:55 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    14600 Oct  5 20:14 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  4 21:08 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   149912 Oct  5 20:14 todo.md

$ grep -rn "v1\.\|1\.3\|1\.2\.2\|upstream" /root/astrate-mule/docs/UPSTREAM-EXPERIMENTAL.md 2>/dev/null | head -50; echo ---; ls /root/astrate-mule/docs/
3:Astarte's brain is the upstream Astarte team. When upstream ships a feature
5:behaviour is ready) but tracks it here until upstream promotes or deprecates
7:upstream version, check each row below against what that release did with the
16:upstream level. When that happens, every open row below tagged with that level
24:| Required + encrypted interface mapping fields | #67 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
25:| `async_operation=false` on mutating RM endpoints | #68 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
---
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
UPSTREAM-EXPERIMENTAL.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui

$ grep -n "1\.[0-9]\+\.[0-9]\+\|v1\.[0-9]" /root/astrate-mule/docs/COMPATIBILITY.md
52:- The Astarte Dashboard v1.2.2 runs unmodified against Astrate (compose
79:     Recorded against upstream v1.2.0 in
133:   v1.2.0 in `test/conformance/upstream/channels.json`, and two rows were
182:     upstream 1.2.0 on 2026-08-24 (`test/conformance/upstream/verify-versions.json`,
253:    reports `realm.APICompatVersion` (currently `1.2.2`), a compatibility
256:    (>= 1.1.1).
259:    1.2.2 serves time series newest-first and has no `sort` parameter; Astrate
270:13. **AMQP trigger actions are rejected at creation** — upstream 1.2.2 accepts
279:    1.2.0 (`test/conformance/upstream/verify-server-writes.json`, 2026-08-24,
304:    1.2.0, and Astrate's persist-then-publish behaves the same way.

$ grep -rn "astarte-upstream\|compat-note\|upstream movement" /root/astrate-mule/.mule/log.md /root/astrate-mule/.mule/todo.md | head -40
/root/astrate-mule/.mule/log.md:18:| 2026-09-05 | compat-note-v132 | done | 73s | 4bc3e1a |
/root/astrate-mule/.mule/log.md:23:| 2026-09-05 | compat-note-v133 | blocked | 44s | wrote nothing |
/root/astrate-mule/.mule/log.md:41:| 2026-09-08 | compat-note-v14-rc | blocked | 326s | wrote nothing |
/root/astrate-mule/.mule/log.md:101:| 2026-09-18 | compat-note-v134 | done | 154s | bcd4383 |
/root/astrate-mule/.mule/log.md:107:| 2026-09-20 | compat-note-custom-action-policy-boundary | done | 97s | 8cd17a8 |
/root/astrate-mule/.mule/log.md:152:| 2026-10-02 | compat-note-v14-rc6 | done | 398s | 0215f97 |
/root/astrate-mule/.mule/todo.md:3:- [x] compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).
/root/astrate-mule/.mule/todo.md:65:`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
/root/astrate-mule/.mule/todo.md:115:- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
/root/astrate-mule/.mule/todo.md:120:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
/root/astrate-mule/.mule/todo.md:138:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
/root/astrate-mule/.mule/todo.md:177:- [x] compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.
/root/astrate-mule/.mule/todo.md:248:- [x] compat-note-v1.4-rc.6 [auto]: record in `.mule/for-giulio.md` that v1.4.0-rc.6 is wire-inert for Astrate, as a delta on the open v1.3.2/v1.3.4 wording proposal rather than a re-derivation. Measured, the whole non-FDO delta is inert: DUP RPC-server reliability is OTP supervisor/GC tuning; the Dashboard's new `unknown_status` is computed client-side from `last_connection && !last_disconnection && !isConnected` (models/Device/index.ts:213-215) and both timestamps are already on Astrate's device shape (internal/appengine/service.go:118), so a device whose disconnection was never recorded reads as "Unknown status" rather than "Disconnected"; the `on_device_registered` dispatch moved from `Agent.register_device` to `Engine.register_device` (pairing/engine.ex:157-165, 186-193) so it also covers the FDO path, but the Agent path is byte-for-byte the same rule and Astrate already fires it on every successful registration (internal/pairing/service.go:221-223); the `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"); and `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma. So the proposal's "v1.4.0 is still rc.5-only" reference becomes rc.6, and its "FDO authentication (pairing, disabled by default)" line understates the 1.4 line, which ships a full ownership-voucher surface Astrate does not emulate (#78, parked). Do not edit docs/COMPATIBILITY.md and do not restate the standing proposal.

$ grep -n "v1\.3\.4\|v1\.3\.5\|rc\.6\|v1\.3\.2" /root/astrate-mule/.mule/for-giulio.md | tail -40
15:  **One detail to fold into that review rather than queue separately:** the rewritten comment at `internal/broker/aclhook.go` (the `control/keyAgreement` deny) says "…upstream's wire spec (082-key_agreement_protocol.md) only ships with v1.4.0-rc.5, v1.3.3 being the newest stable tag". The 082 citation, the "deliberately denied as Astrate's own choice" framing and the unchanged deny + `aclhook_test.go` case are exactly what #93 asked for, but the version reference has aged — upstream's newest stable is **v1.3.5** (2026-10-05) and newest overall `v1.4.0-rc.6` (2026-09-30). No separate task: one word in the same diff when someone finally reads it.
19:- **Milestone recipe run, 2026-10-06: v3.0 verified unchanged — still gated on upstream — and the master delta since the 2026-10-05 run is a closed set of 13 commits, every one wire-inert for Astrate. No issue filed, no `.mule/todo.md` line queued.** Release sweep is byte-identical to yesterday's (`gh api repos/astarte-platform/astarte/releases`): **no stable v1.4.0** — newest stable `v1.3.5` (2026-10-05), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition (your 2026-09-04 "wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) stays unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** — this milestone's whole ledger lives under `upstream-parity`, where the open set is unchanged: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`); **#47–#89** otherwise all closed. Non-alarm open issues total exactly four (plus the 25 `mule-alarm` #94–#118): the three above and **#1** (untouched per standing instruction). `APICompatVersion` still `"1.2.2"` (`internal/realm/service.go:588`, re-checked).
24:- **astarte-upstream recipe run, 2026-10-05 (evening): no upstream movement that touches anything Astrate implements, so nothing was queued — no `.mule/todo.md` line, no `gh issue create`, nothing commented, closed or edited on GitHub.** One thing *is* new since this recipe's last run (`2026-10-04T20:40:52+02:00` ran; the `2026-10-05T04:29:10+02:00` run timed out): stable **v1.3.5**, published `2026-10-05T14:11:44Z`. The milestone recipe already recorded it as inert today; I re-measured it rather than trusting that, because "a release note is a claim, not a fact" is this recipe's own rule. `gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5` is **4 commits, 35 files** — `8d47574f` + `4c4222cf` `chore: update horde`, `c4c63ebe` `chore(rm): bump xandra to 0.19.4`, `8f0f5ef1` `chore: prepare v1.3.5` — and **every app/lib file in it is a `mix.exs`/`mix.lock` version bump**. The rest is `CHANGELOG.md` (+6: one `### Fixed` entry, "[astarte_data_updater_plant] Increase RPC server reliability"), `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `README.md`, `astarte-dashboard/package.json`. Not one source file changed, so no route, MQTT topic, BSON encoding, control message or interface-schema field moved: Astrate implements no `data_updater_plant` and no Erlang RPC server. Wire-inert, same shape as v1.3.2/v1.3.3/v1.3.4. **Master since `v1.4.0-rc.6` is still the single commit `731bcddc`** ("chore: forward port release-1.4"; `gh api .../commits?since=2026-10-05T09:34:00Z` returns it and nothing newer), the migration-relocation commit the milestone run measured today as inert — so both wire-bearing post-rc.6 commits (`b6d46ad4`, `958bb9fe`) still have their owners in the queue and nothing new joins them.
25:  **The parked-decisions row for #92 is still unmet, re-checked against the tag list rather than the release list** (this host has no `.mule/waiting-on.md` — absent from `HEAD`, already reported; I read it out of `origin/main` read-only): `gh api repos/astarte-platform/astarte/tags?per_page=100` gives newest stable **v1.3.5** and newest overall **v1.4.0-rc.6**, with `v1.4.0-rc.0`…`rc.6` and no stable `v1.4.0`. The row asks for a **stable** tag, which an `-rc.N` does not satisfy, so the twenty-second run ends the same way as the twenty-one before it.
26:  **The one delta to your open COMPATIBILITY.md wording proposal (the v1.3.2 one, folding v1.3.3/v1.3.4, further down this file): its version-reference sentence should read "newest stable v1.3.5" where it currently says "v1.3.4". Nothing else in it changes** — v1.3.5 adds no capability, so the same decision and the same "not yet emulated" list stand. Not applied: `docs/COMPATIBILITY.md` is untouched and still targets **1.2.2**, `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). I deliberately did **not** queue a fresh `compat-note-v1.3.5` line: the proposal is already open here, and a second line would re-derive it for the next run to find.
36:- **Milestone recipe run, 2026-10-05: v3.0's gate is unchanged, but the release sweep finally has something new to swallow — upstream shipped stable v1.3.5 today, seven hours ago, and it is wire-inert for Astrate. No issue filed, no task line queued.**  **The new release, measured.** `gh api repos/astarte-platform/astarte/releases` now shows a stable **v1.3.5** at `2026-10-05T14:11:44Z` where every prior run of this recipe (2026-09-10 onward, 20+ entries in this file) saw `v1.3.4` (2026-09-18) as newest. `compare/v1.3.4...v1.3.5` is **4 commits, 35 files**: two `chore: update horde` (`8d47574f`, `4c4222cf`), `chore(rm): bump xandra to 0.19.4` (`c4c63ebe`), `chore: prepare v1.3.5` (`8f0f5ef1`). The file list is `mix.exs`/`mix.lock` version bumps across all apps and libs, `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `astarte-dashboard/package.json`, `doc/mix.exs`, `README.md` — plus a **single** CHANGELOG entry, `## [1.3.5] - 2026-10-05` / `### Fixed` / "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available". That is the whole delta and none of it is Astarte-facing: Astrate implements no data_updater_plant and no Erlang RPC server (the only `rpc` hit under `internal/` is a comment at `internal/pairing/http.go:300` about AMQP reply-code conflicts), so no HTTP route, MQTT topic, AMQP control message or schema field moved. Same shape as the v1.3.2/v1.3.3/v1.3.4 maintenance releases already recorded here, so it is **not a v3.0 gap** and there is nothing to file.
37:  **#92's parking condition is still unmet, now for the twenty-first run.** Your 2026-09-04 decision ("wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) — newest stable is `v1.3.5` (today), newest overall remains `v1.4.0-rc.6` (2026-09-30, prerelease). The post-rc.6 master delta has grown by exactly one commit: `731bcddc` "chore: forward port release-1.4" (2026-10-05T09:34Z), and it is the housekeeping/FDO migration relocation the 2026-10-03 run already measured inside rc.5→rc.6 (six `astarte_housekeeping/priv/migrations/{astarte,realm}/*.sql` files removed, table creation moved into migrations) plus a DUP-internal `apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0011_remove_replacement_data.ex` (+31) and `database.ex` (+2/-1). No Realm Management, AppEngine, Pairing or Trigger wire surface. **The milestone ledger is unchanged:** `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** (this milestone's whole ledger lives under `upstream-parity`, as the 2026-09-10 entry established), of #47–#89 **all remain closed**, and the authoritative open set is 29 issues — 25 alarms (#94–#118, one new today: **#118**, that recipe's business) and exactly four non-alarm: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** (untouched per standing instruction). `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`, re-checked), per **#90's** frozen decision.
38:  **The standing item is unchanged and I am not re-litigating it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries its two rows, both tagged `1.4 experimental` with trigger "upstream 1.4 final: promoted or removed" — the mis-tagging the 2026-10-04 entry measured out of upstream's CHANGELOG (row #68's `async_operation` feature is 1.0.2 and stable for four years; row #67's `required` half is a plain 1.4.0-rc.0 addition, only its `encrypted` half being genuinely experimental and already tracked by #92/#93). v1.3.5 touches neither feature, so it changes nothing about your four pending decisions (a)–(d) in that entry, and I have not edited the register.
39:  **Step 5 deliberately not taken.** I am *not* proposing "milestone v3.0 looks complete, verify and cut the tag", for the same reason as the twenty runs before this one: the section's declared scope is **1.2.2 → 1.3/1.4** and only the 1.3 half is delivered (now including v1.3.5, which is inert anyway). The 1.4 half is open — #92 unanswered, both experimental rows unreconciled, `APICompatVersion` not bumped — so the tag must not be cut. Your call when a stable v1.4.0 ships: answer #92, settle the register rows, run the final-phase bump, then cut the tag.
44:- **github-issues triage run, 2026-10-04: nothing proposable for the 22nd run, and the reason the alarm pile has never been actioned is finally located — `mule/queue` forked from `main` 27 minutes before your four commits landed there, so this Pi has been running for a month with a dead-man's switch, a dashboard and a parked-decision check that all point at things which do not exist here.** 28 open issues (authoritative `--limit 100` list: `1, 78, 92, 93, 94…117`; note the recipe's own `--limit 40` command printed 27 of these and silently dropped **#94**, which is open with `mule-alarm` — so treat "27 open" in the 2026-10-03 entry as 28). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`; `8c61268` re-verified still an ancestor of `origin/mule/queue` this run), **#92** keyAgreement (parked on a stable v1.4.0 — measured again today, still unmet: newest upstream stable `v1.3.4`, newest overall `v1.4.0-rc.6`, 2026-09-30), **#78** FDO (`milestone-4.0`, already escalated), **#1** untouched per standing instruction.
58:  **Gate re-verified, unchanged.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — nothing shipped upstream since the 2026-10-03 run, so #92's parking condition (your 2026-09-04 decision) is still unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** (this milestone's ledger lives under `upstream-parity`); of #47–#89 **all remain closed**; open non-alarm set unchanged — **#93** (`mule-review`, own recipe path), **#92** (parked), **#78** (`milestone-4.0`), **#1** (untouched). The post-rc.6 master delta is now a closed set of **8 commits** (`gh api .../commits -f since=2026-09-30T07:28:27Z`), of which exactly two carry wire surface and **both already have owners in the queue**: `b6d46ad4` (#2237 unexpected object keys) → done as `appengine-unexpected-object-key`, and `958bb9fe` (#2232 `device_empty_cache_received`) → queued as `device-empty-cache-received-trigger` (currently `[!]` blocked). Neither is from stable 1.3 or 1.4, so neither is a v3.0 gate.
79:  **Verified unchanged, no new gaps.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition, frozen by you on 2026-09-04 ("wait for a stable v1.4.0", an `-rc.N` does not satisfy it), is still unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Of #47–#89 **all are closed**; the open non-alarm set is unchanged — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** (untouched per standing instruction). rc.6 adds no v3.0 gap: it is already audited in the two entries below, and I re-measured its non-FDO half myself (`compare/v1.4.0-rc.5...v1.4.0-rc.6`, 27 commits) — outside FDO the only library changes are a `# coveralls-ignore-stop` line in `apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex:43` and `realms/queries.ex` *removing* the ownership_voucher / to2_sessions / unconfirmed_devices table creation (moved into migrations). No Realm Management, AppEngine or Pairing wire surface moved.
89:- **Delta on the two open COMPATIBILITY.md wording entries further down this file (the v1.3.4 one folding in v1.3.3, and the original v1.3.2 one): upstream `v1.4.0-rc.6` is wire-inert for Astrate, so nothing in that proposal is re-derived — only its version reference and one line of its "not yet emulated" list change.** Not applied: `docs/COMPATIBILITY.md` is untouched (still the 1.2.2 target, deviation 18's "which upstream 404s" at `docs/COMPATIBILITY.md:360`, `APICompatVersion = "1.2.2"` at `internal/realm/service.go:588`), and the FDO half of the rc.6 delta already has its own entry below rather than being folded in here.
93:  - **The two edits the proposal needs, both yours to approve:** (1) the 2026-09-18 entry's parenthetical "(v1.4.0 is still rc.5-only)" becomes "(v1.4.0 is still rc.6-only)" — re-verified this session that rc.6 is the newest 1.4 prerelease and there is still **no stable v1.4.0** (`gh api repos/astarte-platform/astarte/releases`); (2) its "**FDO authentication** (pairing, disabled by default)" item, and the same "experimental FDO pairing auth" phrase in the proposed doc sentence, understates the 1.4 line, which ships a full ownership-voucher surface (upload with mandatory `hw_id`, deletion, re-TO0, per-voucher `expiry`, a spec-mandated realm-less TO2 path) that Astrate does not emulate — #78, deliberately parked at `.mule/milestones.md:124`. "FDO ownership vouchers (1.4 line, not stable upstream)" says it without claiming the rc.5 shape.
94:  Unverified: the upstream side of the item list is the upstream-watch recipe's `v1.4.0-rc.5...v1.4.0-rc.6` comparison — the same measurement as the FDO entry below, taken there and not re-derived here (a raw fetch of `models/Device/index.ts` at the rc.6 tag did not resolve from this box). I re-verified only the Astrate-side claims above and the release list. No `.mule/todo.md` lines: the deliverable is the delta, and the proposal itself stays open.
98:- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published 2026-09-30T07:28Z, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and the one path #78's body states as fact no longer exists upstream. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix") deletes the `realm_name` path parameter from all seven TO2 `operation` specs in `fdo_onboarding_controller.ex`. **This breaking change is not in the rc.6 CHANGELOG entry** (which lists only the to0 route, voucher deletion, the `ASTARTE_BASE_URL_*` rename, device-deletion cascade, `device_id` on upload, immediate registration, the two x5chain/voucher-store fixes and the cross-realm delete fix) — so anyone tracking upstream by its changelog, including #78's body, misses it. The realm now comes from the voucher row: the global-keyspace query `fetch_device_id_and_ownership_voucher_and_realm/1` returns `{realm, device_id, voucher_data}` and the `:fdo` plug fetches the session by GUID alone (`plug/fdo_session.ex`: `Session.fetch(guid)`); `RvTO2Addr` no longer carries the realm in the hostname either ("the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure"). **Yours:** #78's acceptance is "a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end" — against rc.6 that SDK posts to `/fdo/101/msg/{60..70}` with no realm, so mirroring rc.5's shape would mean mirroring a path upstream has deleted. No action while parked; the body is just stale.
108:- **Not client-observable, so not part of the decision:** the keyspace move itself — vouchers and TO2 sessions now live in the global `astarte` keyspace (`Realm.astarte_keyspace_name()` throughout `fdo/queries.ex`) with a `realm` column on the voucher row and `ALLOW FILTERING` on the list query, so **realm isolation is now application-level rather than keyspace-level** (which is what makes `ensure_voucher_in_realm`-style cross-realm checks necessary, `b3bf731` "ensure users cannot delete ownership vouchers belonging to other realms"); `fdo_guid` added to the device row (migration `realm/0025`) so device removal can find the voucher; the `unconfirmed_devices` drop (`realm/0020`, schema file deleted, `remove_device_from_unconfirmed_devices!` removed); the dashboard voucher page (`astarte-dashboard/src/FdoVoucherPage.tsx`, `hooks/useFdo.ts`, a new cypress spec); and the `data_updater_plant` RPC/GC fixes plus the `v1.3.4` forward-ports.
109:- **One premise in #78's own comment is now false.** It says "upstream ships almost no public FDO documentation (code-level only); ours must exceed that." rc.6 adds two documented sections: `doc/pages/user/035-register_device.md` gained `## FIDO Device Onboard` (Owner Keys, Ownership Vouchers) and `## Credentials Secret Lifecycle`, and `doc/pages/architecture/050-pairing_mechanism.md` gained `## FIDO Device Onboarding` (TO0, TO2, Owner Keys management) — commits `162062e`, `d1c39fb`. "Docs must exceed upstream" is still the right requirement, but the gap is far smaller than #78 assumes, and those two pages are now the cheapest available scope input.
110:- **Where Astrate stands, for the size estimate:** it has **no FDO surface at all** — the only FDO reference in the tree is the comment on the realm health probe (`internal/pairing/http.go:79`, `internal/pairing/service.go:368`), and `.mule/milestones.md:124` still lists #78 as a deliberately parked item. So none of this is a delta to partially-built code; every item above is net-new scope for a "Size L" issue that is blocked on "who runs a rendezvous server?". The rc.6 delta does not change that call, but it adds four things any eventual implementation must decide: a spec-mandated device-side path with no realm (so any deviation is ours alone), a mandatory `hw_id` in the upload body (so the owner must supply the device id *before* the device has onboarded), an upload that creates a visible realm device as a side effect, and a per-voucher `expiry` the operator has to poll to know when the rendezvous stops serving.
111:  Unverified: read-only — every line above comes from the compare API and the rc.6-tagged sources, and the 409/500 mappings come from `FallbackController` source rather than a live call (no Docker on this box). Deliberately no `.mule/todo.md` lines: the deliverable is the scope update, and #78 stays parked until you answer the rendezvous question.
135:- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Status: not started", which is no longer true.** Release sweep (`gh api repos/astarte-platform/astarte/releases`, run this session) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.5` (2026-08-20, prerelease) — byte-identical to what the 2026-09-24 run saw, so nothing shipped upstream in the two days since. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Open non-alarm issues remain exactly **#93** (aclhook comment, `mule-review`, its own recipe path), **#92** (keyAgreement, parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** (FDO, `milestone-4.0`), **#1** (untouched per standing instruction). No new gaps, **no issues filed, no `.mule/todo.md` task lines added** — there is nothing machine-checkable left, and filing anything here would duplicate #92/#93/#78.
210:- **Milestone recipe run, 2026-09-24: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-23 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
218:- **Milestone recipe run, 2026-09-23: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-21 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
234:- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
242:- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
254:- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
258:- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
266:- **COMPATIBILITY.md wording update for upstream v1.3.4 (newest stable, 2026-09-18; v1.4.0 is still rc.5-only) — folds v1.3.3 and v1.3.4 into the open v1.3.2 proposal below; both are maintenance-only and wire-inert, so that proposal stands unchanged.** v1.3.3 (2026-08-07) was an empty-body patch and v1.3.4 (2026-09-18, `astarte_data_updater_plant` only — more frequent fullsweep GC on AMQPDataConsumer processes plus an RPC-availability fix; release body via `gh api repos/astarte-platform/astarte/releases` `v1.3.4`) introduce **no wire/API surface change**, so the v1.3.2 wording proposal stands complete: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"). The only delta a v1.3.4-aware doc carries is the version reference: the proposed §Infrastructure-differences sentence's "until the milestone that adopts v1.3.2 as the target" reads "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target". Raw: [v1.3.4](https://github.com/astarte-platform/astarte/releases/tag/v1.3.4).
274:- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) shows a new **stable v1.3.4** (2026-09-18) — maintenance only (data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change), so no new gap for Astrate — and still **no stable v1.4.0**: newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
419:- **COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable, 2026-07-14; v1.4.0 is still rc-only).** Astrate's doc and `APICompatVersion` still target upstream **1.2.2** (`internal/realm/service.go:588`); v1.3.0 (2026-05-06) introduced wire-surface changes Astrate does not yet emulate, so this is a decision — adopt v1.3.2 as the compatibility target (then update the doc + bump `APICompatVersion` together, per the bump rule) or keep 1.2.2 and add a "not yet emulated" note. Wire-relevant v1.3.0 deltas (release notes): **MQTT v1 capabilities** incl. `purge_properties_compression_format` (plaintext vs zlib purge — touches the `emptyCache`/`producer,properties` contract COMPATIBILITY.md deviation 1 documents); **empty introspection now allowed**; **device registration triggers** (pairing) and **device deletion started/completed triggers** (RM — the latter two already exist as Astrate deviation 9 emits both around the synchronous delete); **FDO authentication** (pairing, disabled by default); **realm-scoped health** — upstream v1.3 added `GET /pairing/v1/{realm}/health`, which Astrate already serves (`internal/pairing/http.go:78-81`, comment already says "upstream 1.3+"), so deviation 18's wording ("which upstream 404s") is now false against 1.3 and the note should be reworded either way. Proposed doc wording (for your approval, edit to taste): in §Infrastructure differences add a sentence — *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl.
421:plaintext `purge_properties_compression_format`, empty-introspection allowance, device registration/deletion triggers, experimental FDO pairing auth) are not yet emulated and are out of scope until the milestone that adopts v1.3.2 as the target."* — and reword deviation 18's realm-health note from "which upstream 404s" to "added by upstream v1.3 (Astrate serves it against a 1.2.2 target; kept, matching behavior)". Raw upstream changes: [v1.3.0](https://github.com/astarte-platform/astarte/releases/tag/v1.3.0), [v1.3.2](https://github.com/astarte-platform/astarte/releases/tag/v1.3.2).

$ grep -n "compat-note-v134\|compat-note-v14-rc6\|astarte-upstream" /root/astrate-mule/.mule/log.md
101:| 2026-09-18 | compat-note-v134 | done | 154s | bcd4383 |
152:| 2026-10-02 | compat-note-v14-rc6 | done | 398s | 0215f97 |

$ gh api repos/astarte-platform/astarte/commits --jq '.[0:8][] | "\(.sha[0:8])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"'
1da04832	2026-10-05T14:57:19Z	chore: forward port release-1.4
a8a26d0f	2026-10-05T14:40:35Z	chore: forward port release-1.3
d85e0ed0	2026-10-06T10:02:32Z	docs: sync FDO workflow to the current behavior (#2241)
a6bd2c21	2026-10-06T09:59:16Z	fix(fdo): reject owner key names not accepted by OpenBao (#2252)
a51ab2e1	2026-10-06T07:21:18Z	refactor: Update Umbrella checks (#2244)
7628d2ed	2026-10-05T15:49:18Z	test(generators): validate interfaces with Core changesets
ed1e99f4	2026-10-05T15:49:04Z	test(generators): validate mappings with Core changesets
bb17455a	2026-10-05T15:48:50Z	test(generators): validate policy handlers with Core changesets

$ gh api 'repos/astarte-platform/astarte/compare/v1.3.5...master' --jq '"\(.status) ahead_by=\(.ahead_by) commits=\(.commits|length) files=\(.files|length)"'
ahead ahead_by=1812 commits=250 files=300

$ gh api 'repos/astarte-platform/astarte/commits?since=2026-10-05T14:57:19Z&per_page=100' --jq 'length as $n | "count=\($n)", (.[] | "\(.sha[0:8])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])")'
count=9
1da04832	2026-10-05T14:57:19Z	chore: forward port release-1.4
a8a26d0f	2026-10-05T14:40:35Z	chore: forward port release-1.3
d85e0ed0	2026-10-06T10:02:32Z	docs: sync FDO workflow to the current behavior (#2241)
a6bd2c21	2026-10-06T09:59:16Z	fix(fdo): reject owner key names not accepted by OpenBao (#2252)
a51ab2e1	2026-10-06T07:21:18Z	refactor: Update Umbrella checks (#2244)
7628d2ed	2026-10-05T15:49:18Z	test(generators): validate interfaces with Core changesets
ed1e99f4	2026-10-05T15:49:04Z	test(generators): validate mappings with Core changesets
bb17455a	2026-10-05T15:48:50Z	test(generators): validate policy handlers with Core changesets
ec020bf8	2026-10-05T15:48:31Z	fix(generators): generate valid trigger policy names

→ Read .mule/todo.md [limit=70]

$ grep -n "^- \[ \] compat-note\|^- \[!\] compat-note\|^- \[ \] probe-" /root/astrate-mule/.mule/todo.md
120:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
138:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing

$ for s in 1da04832 a8a26d0f d85e0ed0 a6bd2c21 a51ab2e1 ec020bf8; do echo "== $s"; gh api repos/astarte-platform/astarte/commits/$s --jq '.files[] | "\(.filename) +\(.additions)/-\(.deletions)"'; done
== 1da04832
.credo.exs +2/-1
.dockerignore +8/-1
.formatter.exs +5/-0
.github/codecov.yml +20/-11
.github/labeler.yml +0/-17
.github/workflows/astarte-appengine-api-workflow.yaml +0/-39
.github/workflows/astarte-apps-build-workflow.yaml +1/-56
.github/workflows/astarte-apps-workflow.yaml +191/-0
.github/workflows/astarte-build-workflow.yaml +7/-7
.github/workflows/astarte-dashboard-compatibility.yml +1/-1
.github/workflows/astarte-dashboard-docker-build.yml +1/-1
.github/workflows/astarte-dashboard-tests.yml +1/-1
.github/workflows/astarte-data-updater-plant-workflow.yaml +0/-39
.github/workflows/astarte-e2e-lint-workflow.yaml +46/-0
.github/workflows/astarte-e2e-next-test-build-workflow.yaml +30/-0
.github/workflows/astarte-end-to-end-test-workflow.yaml +82/-2
.github/workflows/astarte-housekeeping-workflow.yaml +0/-39
.github/workflows/astarte-libs-build-workflow.yaml +0/-144
.github/workflows/astarte-libs-workflow.yml +0/-79
.github/workflows/astarte-pairing-workflow.yaml +0/-39
.github/workflows/astarte-realm-management-workflow.yaml +0/-39
.github/workflows/astarte-trigger-engine-workflow.yaml +0/-39
.github/workflows/docs-check.yml +4/-1
.github/workflows/docs-workflow.yaml +3/-1
.github/workflows/pr-container-build-workflow.yaml +1/-0
.github/workflows/publish-dashboard-snapshot-to-dockerhub.yaml +1/-1
.github/workflows/publish-snapshot-to-dockerhub-workflow.yaml +1/-1
.github/workflows/publish-tool-snapshot-to-dockerhub-workflow.yaml +1/-1
.github/workflows/publish-vmq-plugin-release-to-dockerhub-workflow.yaml +107/-0
.github/workflows/publish-vmq-plugin-snapshot-to-dockerhub-workflow.yaml +140/-0
.github/workflows/release-please.yml +27/-0
.gitignore +5/-0
.release-please-manifest.json +3/-0
.reuse/dep5 +61/-0
.tool-versions +2/-2
.typos.toml +1/-1
CHANGELOG.md +4/-0
Dockerfile +30/-30
LICENSES/Apache-2.0.txt +73/-0
LICENSES/CC-BY-4.0.txt +156/-0
LICENSES/CC0-1.0.txt +121/-0
apps/astarte_adapters/.formatter.exs +0/-0
apps/astarte_adapters/README.md +0/-0
apps/astarte_adapters/lib/astarte/adapters.ex +31/-8
apps/astarte_adapters/lib/astarte/adapters/engine.ex +0/-0
apps/astarte_adapters/lib/astarte/adapters/missing_field_error.ex +0/-0
apps/astarte_adapters/lib/astarte/core/interface.ex +65/-0
apps/astarte_adapters/lib/astarte/core/triggers/policy.ex +77/-0
apps/astarte_adapters/lib/astarte/data_access/interface.ex +72/-0
apps/astarte_adapters/mix.exs +31/-12
apps/astarte_adapters/test/astarte/adapters/engine_test.exs +0/-0
apps/astarte_adapters/test/astarte/adapters_test.exs +15/-3
apps/astarte_adapters/test/astarte/core/interface_test.exs +67/-0
apps/astarte_adapters/test/astarte/core/triggers/policy_test.exs +44/-0
apps/astarte_adapters/test/astarte/data_access/interface_test.exs +115/-0
apps/astarte_adapters/test/support/complex_struct.ex +0/-0
apps/astarte_adapters/test/support/generators.ex +1/-1
apps/astarte_adapters/test/support/mappings.ex +13/-0
apps/astarte_adapters/test/support/simple_struct.ex +0/-0
apps/astarte_adapters/test/test_helper.exs +0/-0
apps/astarte_appengine_api/config/ci.exs +0/-4
apps/astarte_appengine_api/config/config.exs +0/-66
apps/astarte_appengine_api/config/dev.exs +0/-37
apps/astarte_appengine_api/config/prod.exs +0/-45
apps/astarte_appengine_api/config/runtime.exs +0/-30
apps/astarte_appengine_api/config/test.exs +0/-70
apps/astarte_appengine_api/lib/astarte_appengine_api/application.ex +4/-0
apps/astarte_appengine_api/lib/astarte_appengine_api/auth/queries.ex +0/-2
apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex +297/-61
apps/astarte_appengine_api/lib/astarte_appengine_api/device/interface_value.ex +14/-1
apps/astarte_appengine_api/lib/astarte_appengine_api/device/queries.ex +185/-258
apps/astarte_appengine_api/lib/astarte_appengine_api/rooms/events_dispatcher.ex +3/-3
apps/astarte_appengine_api/lib/astarte_appengine_api/rooms/room.ex +19/-173
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant.ex +0/-69
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/behaviour.ex +0/-32
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/delete_volatile_trigger/request_data.ex +0/-36
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/install_volatile_trigger/request_data.ex +0/-46
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/vmq_plugin.ex +0/-1
apps/astarte_appengine_api/lib/astarte_appengine_api/stats/queries.ex +0/-2
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/device_status_by_alias_controller.ex +2/-2
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/device_status_by_group_controller.ex +2/-2
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/device_status_controller.ex +2/-2
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/fallback_controller.ex +2/-2
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/interface_values_by_group_controller.ex +12/-37
apps/astarte_appengine_api/lib/astarte_appengine_api_web/telemetry.ex +1/-1
apps/astarte_appengine_api/lib/astarte_appengine_api_web/telemetry/api_usage.ex +2/-2
apps/astarte_appengine_api/lib/astarte_appengine_api_web/telemetry/database_events.ex +0/-1
apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex +2/-2
apps/astarte_appengine_api/mix.exs +37/-40
apps/astarte_appengine_api/mix.lock +0/-97
apps/astarte_appengine_api/rel/config.exs +0/-49
apps/astarte_appengine_api/test/astarte_appengine_api/auth/auth_test.exs +1/-1
apps/astarte_appengine_api/test/astarte_appengine_api/data_transmitter_test.exs +1/-1
apps/astarte_appengine_api/test/astarte_appengine_api/device/device_reading_v2_test.exs +183/-23
apps/astarte_appengine_api/test/astarte_appengine_api/device/device_test.exs +2/-2
apps/astarte_appengine_api/test/astarte_appengine_api/device/device_v2_reading_test.exs +7/-11
apps/astarte_appengine_api/test/astarte_appengine_api/device/device_v2_test.exs +3/-3
apps/astarte_appengine_api/test/astarte_appengine_api/device/interface_value_test.exs +2/-1
apps/astarte_appengine_api/test/astarte_appengine_api/stats/stats_test.exs +2/-1
apps/astarte_appengine_api/test/astarte_appengine_api_web/auth/auth_test.exs +3/-3
apps/astarte_appengine_api/test/astarte_appengine_api_web/channels/rooms_channel_test.exs +105/-144
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/device_status_by_alias_controller_test.exs +3/-3
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/device_status_by_group_controller_test.exs +3/-3
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/device_status_controller_test.exs +3/-3
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/groups_controller_test.exs +8/-8
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_controller_test.exs +15/-17
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_values_by_device_alias_controller_test.exs +7/-4
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_values_by_group_controller_test.exs +7/-4
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_values_controller_test.exs +7/-4
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/stats_controller_test.exs +3/-3
apps/astarte_appengine_api/test/astarte_appengine_api_web/plug/group_name_decoder_test.exs +3/-3
apps/astarte_appengine_api/test/astarte_appengine_api_web/views/error_view_test.exs +1/-1
apps/astarte_appengine_api/test/support/cases/channel.ex +1/-1
apps/astarte_appengine_api/test/support/cases/conn.ex +1/-1
apps/astarte_appengine_api/test/support/cases/data.ex +4/-4
apps/astarte_appengine_api/test/support/cases/device.ex +182/-28
apps/astarte_appengine_api/test/support/generators/interface_update.ex +5/-5
apps/astarte_appengine_api/test/support/generators/interface_values_retrieveal.ex +1/-1
apps/astarte_appengine_api/test/support/helpers/database.ex +17/-4
apps/astarte_appengine_api/test/support/helpers/database_v2.ex +16/-2
apps/astarte_appengine_api/test/support/helpers/device.ex +11/-9
apps/astarte_appengine_api/test/support/helpers/interface.ex +1/-1
apps/astarte_appengine_api/test/support/helpers/jwt.ex +1/-1
apps/astarte_appengine_api/test/support/mocks.ex +0/-4
apps/astarte_appengine_api/test/test_helper.exs +3/-0
apps/astarte_config/.formatter.exs +0/-0
apps/astarte_config/README.md +0/-0
apps/astarte_config/lib/astarte_config/astarte_config.ex +0/-0
apps/astarte_config/lib/astarte_config/binding/binding.ex +0/-0
apps/astarte_config/lib/astarte_config/binding/request_opts.ex +0/-0
apps/astarte_config/lib/astarte_config/binding/scheme.ex +0/-0
apps/astarte_config/lib/astarte_config/binding/url.ex +0/-0
apps/astarte_config/lib/astarte_config/http_client.ex +3/-0
apps/astarte_config/mix.exs +10/-5
apps/astarte_config/test/astarte_config/astarte_config_test.exs +21/-3
apps/astarte_config/test/astarte_config/http_client_test.exs +6/-0
apps/astarte_config/test/support/sample_client.ex +0/-0
apps/astarte_config/test/support/sample_config.ex +0/-0
apps/astarte_config/test/test_helper.exs +0/-0
apps/astarte_core/.formatter.exs +2/-2
apps/astarte_core/.github/dco.yml +2/-0
apps/astarte_core/.github/workflows/build-workflow.yaml +112/-0
apps/astarte_core/.gitignore +7/-17
apps/astarte_core/AUTHORS.md +4/-0
apps/astarte_core/CHANGELOG.md +267/-0
apps/astarte_core/LICENSE +0/-0
apps/astarte_core/README.md +3/-0
apps/astarte_core/SECURITY.md +23/-0
apps/astarte_core/lib/astarte_core/cql_utils.ex +182/-0
apps/astarte_core/lib/astarte_core/device.ex +89/-0
apps/astarte_core/lib/astarte_core/device/capabilities.ex +22/-13
apps/astarte_core/lib/astarte_core/group.ex +19/-10
apps/astarte_core/lib/astarte_core/interface.ex +325/-0
apps/astarte_core/lib/astarte_core/interface/aggregation.ex +106/-0
apps/astarte_core/lib/astarte_core/interface/ownership.ex +117/-0
apps/astarte_core/lib/astarte_core/interface/type.ex +106/-0
apps/astarte_core/lib/astarte_core/interface_descriptor.ex +134/-0
apps/astarte_core/lib/astarte_core/mapping.ex +305/-0
apps/astarte_core/lib/astarte_core/mapping/database_retention_policy.ex +110/-0
apps/astarte_core/lib/astarte_core/mapping/endpoints_automaton.ex +186/-0
apps/astarte_core/lib/astarte_core/mapping/reliability.ex +111/-0
apps/astarte_core/lib/astarte_core/mapping/retention.ex +110/-0
apps/astarte_core/lib/astarte_core/mapping/value_type.ex +216/-0
apps/astarte_core/lib/astarte_core/proto/astarte_reference.pb.ex +8/-0
apps/astarte_core/lib/astarte_core/proto/astarte_reference.proto +24/-0
apps/astarte_core/lib/astarte_core/realm.ex +47/-0
apps/astarte_core/lib/astarte_core/storage_type.ex +140/-0
apps/astarte_core/lib/astarte_core/triggers/data_trigger.ex +56/-0
apps/astarte_core/lib/astarte_core/triggers/policy/error_type.ex +95/-0
apps/astarte_core/lib/astarte_core/triggers/policy/handler.ex +157/-0
apps/astarte_core/lib/astarte_core/triggers/policy/keyword_error.ex +42/-0
apps/astarte_core/lib/astarte_core/triggers/policy/policy.ex +167/-0
apps/astarte_core/lib/astarte_core/triggers/policy/range_error.ex +19/-12
apps/astarte_core/lib/astarte_core/triggers/policy_protobuf.ex +8/-7
apps/astarte_core/lib/astarte_core/triggers/policy_protobuf/error_keyword.pb.ex +20/-0
apps/astarte_core/lib/astarte_core/triggers/policy_protobuf/error_keyword.proto +31/-0
apps/astarte_core/lib/astarte_core/triggers/policy_protobuf/error_range.pb.ex +7/-0
apps/astarte_core/lib/astarte_core/triggers/policy_protobuf/error_range.proto +23/-0
apps/astarte_core/lib/astarte_core/triggers/policy_protobuf/handler.pb.ex +29/-0
apps/astarte_core/lib/astarte_core/triggers/policy_protobuf/handler.proto +36/-0
apps/astarte_core/lib/astarte_core/triggers/policy_protobuf/policy.pb.ex +17/-0
apps/astarte_core/lib/astarte_core/triggers/policy_protobuf/policy.proto +30/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_connected_event.pb.ex +7/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_connected_event.proto +23/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_deletion_finished_event.pb.ex +5/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_deletion_finished_event.proto +22/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_deletion_started_event.pb.ex +5/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_deletion_started_event.proto +22/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_disconnected_event.pb.ex +5/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_disconnected_event.proto +22/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_empty_cache_received_event.pb.ex +5/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_empty_cache_received_event.proto +4/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_error_event.pb.ex +21/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_error_event.proto +24/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_registered_event.pb.ex +5/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_registered_event.proto +22/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/encoder.ex +385/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/incoming_data_event.pb.ex +9/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/incoming_data_event.proto +25/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/incoming_introspection_event.pb.ex +31/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/incoming_introspection_event.proto +29/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/interface_added_event.pb.ex +9/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/interface_added_event.proto +25/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/interface_minor_updated_event.pb.ex +10/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/interface_minor_updated_event.proto +26/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/interface_removed_event.pb.ex +8/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/interface_removed_event.proto +24/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/path_created_event.pb.ex +9/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/path_created_event.proto +25/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/path_removed_event.pb.ex +8/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/path_removed_event.proto +24/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/simple_event.pb.ex +99/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/simple_event.proto +68/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/value_change_applied_event.pb.ex +10/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/value_change_applied_event.proto +26/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/value_change_event.pb.ex +10/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/value_change_event.proto +26/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/value_stored_event.pb.ex +9/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/value_stored_event.proto +25/-0
apps/astarte_core/lib/astarte_core/triggers/simple_trigger_config.ex +689/-0
apps/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/amqp_trigger_target.pb.ex +30/-0
apps/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/amqp_trigger_target.proto +35/-0
apps/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/data_trigger.pb.ex +56/-0
apps/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/data_trigger.proto +55/-0
apps/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/device_trigger.pb.ex +36/-0
apps/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/device_trigger.proto +44/-0
apps/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/simple_trigger_container.pb.ex +19/-0
apps/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/simple_trigger_container.proto +31/-0
apps/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/tagged_simple_trigger.pb.ex +13/-0
apps/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/tagged_simple_trigger.proto +29/-0
apps/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/trigger_target_container.pb.ex +14/-0
apps/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/trigger_target_container.proto +29/-0
apps/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/utils.ex +157/-0
apps/astarte_core/lib/astarte_core/triggers/trigger.pb.ex +12/-0
apps/astarte_core/lib/astarte_core/triggers/trigger.proto +31/-0
apps/astarte_core/mix.exs +96/-0
apps/astarte_core/specs/interface.json +56/-0
apps/astarte_core/specs/mapping.json +89/-0
apps/astarte_core/specs/policy.json +87/-0
apps/astarte_core/test/astarte_core/cql_utils_test.exs +177/-0
apps/astarte_core/test/astarte_core/device/capabilities_test.exs +73/-0
apps/astarte_core/test/astarte_core/device_test.exs +68/-0
apps/astarte_core/test/astarte_core/group_test.exs +26/-0
apps/astarte_core/test/astarte_core/interface/enums_test.exs +32/-0
apps/astarte_core/test/astarte_core/interface_descriptor_test.exs +103/-0
apps/astarte_core/test/astarte_core/interface_test.exs +539/-0
apps/astarte_core/test/astarte_core/mapping/endpoints_automaton_test.exs +383/-0
apps/astarte_core/test/astarte_core/mapping/enums_test.exs +74/-0
apps/astarte_core/test/astarte_core/mapping/value_type_test.exs +118/-0
apps/astarte_core/test/astarte_core/mapping_test.exs +377/-0
apps/astarte_core/test/astarte_core/proto_astarte_reference_test.exs +21/-0
apps/astarte_core/test/astarte_core/realm_test.exs +35/-0
apps/astarte_core/test/astarte_core/storage_type_test.exs +79/-0
apps/astarte_core/test/astarte_core/triggers/data_trigger_test.exs +82/-0
apps/astarte_core/test/astarte_core/triggers/policy/error_type_test.exs +34/-0
apps/astarte_core/test/astarte_core/triggers/policy/handler_test.exs +150/-0
apps/astarte_core/test/astarte_core/triggers/policy/policy_protobuf_test.exs +66/-0
apps/astarte_core/test/astarte_core/triggers/policy/policy_test.exs +343/-0
apps/astarte_core/test/astarte_core/triggers/simple_events_encoder_test.exs +574/-0
apps/astarte_core/test/astarte_core/triggers/simple_events_test.exs +482/-0
apps/astarte_core/test/astarte_core/triggers/simple_trigger_config_test.exs +882/-0
apps/astarte_core/test/astarte_core/triggers/simple_triggers_protobuf/utils_test.exs +75/-0
apps/astarte_core/test/astarte_core/triggers/simple_triggers_protobuf_test.exs +166/-0
apps/astarte_core/test/test_helper.exs +1/-0
apps/astarte_data_access/.formatter.exs +9/-0
apps/astarte_data_access/.gitignore +0/-0
apps/astarte_data_access/CHANGELOG.md +0/-0
apps/astarte_data_access/README.md +0/-0
apps/astarte_data_access/coveralls.json +0/-0
apps/astarte_data_access/lib/astarte_data_access.ex +0/-0
apps/astarte_data_access/lib/astarte_data_access/bigint.ex +0/-0
apps/astarte_data_access/lib/astarte_data_access/c_system.ex +0/-0
apps/astarte_data_access/lib/astarte_data_access/config.ex +0/-0
apps/astarte_data_access/lib/astarte_data_access/config/xandra_nodes.ex +0/-0
apps/astarte_data_access/lib/astarte_data_access/consistency.ex +0/-0
apps/astarte_data_access/lib/astarte_data_access/data.ex +3/-4
apps/astarte_data_access/lib/astarte_data_access/database/database.ex +189/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0001_migrate_schema_version.ex +57/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0002_init_database.ex +35/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0003_drop_astarte_schema.ex +8/-4
apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0004_remove_replication_factor.ex +29/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0005_add_device_registration_limit.ex +29/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0006_create_kv_store.ex +10/-11
apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0007_create_session_key_type.ex +8/-6
apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0008_create_device_sessions.ex +53/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0009_create_ownership_vouchers.ex +40/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0010_add_expiry_to_ownership_vouchers.ex +29/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0011_remove_replacement_data.ex +31/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0001_migrate_schema_version.ex +57/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0002_init_database.ex +127/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0003_create_grouped_devices.ex +31/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0004_add_groups_and_interface_stats.ex +8/-13
apps/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0005_add_database_retention_ttl_and_policy.ex +30/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0006_add_attributes.ex +29/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0007_create_deletion_in_progress.ex +32/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0008_create_capabilities_type.ex +31/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0009_add_capabilities.ex +29/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0010_add_groups_to_deletion_in_progress.ex +29/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0011_create_ownership_voucher.ex +31/-0
apps/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0012_create_device_sessions.ex +52/-0
== a8a26d0f
CHANGELOG.md +6/-0
README.md +2/-2
libs/astarte_adapters/mix.exs +1/-1
libs/astarte_test_suite/mix.exs +1/-1
== d85e0ed0
doc/pages/architecture/050-pairing_mechanism.md +4/-5
doc/pages/user/015-astarte_dashboard.md +2/-1
doc/pages/user/035-register_device.md +11/-0
== a6bd2c21
apps/astarte_pairing/test/astarte_pairing_web/controllers/owner_key_controller_test.exs +18/-0
libs/astarte_secrets/lib/astarte_secrets/owner_key_initialization_options.ex +8/-0
libs/astarte_secrets/test/astarte_secrets/core_test.exs +13/-0
libs/astarte_secrets/test/astarte_secrets/owner_key_initialization_test.exs +35/-0
== a51ab2e1
.github/workflows/astarte-apps-build-workflow.yaml +0/-73
.github/workflows/astarte-apps-workflow.yaml +72/-0
apps/astarte_adapters/mix.exs +1/-1
apps/astarte_generators/mix.exs +1/-1
apps/astarte_test_suite/mix.exs +1/-1
mix.exs +2/-1
== ec020bf8
apps/astarte_generators/lib/astarte/core/generators/triggers/policy.ex +2/-1
apps/astarte_generators/test/astarte/core/generators/triggers/policy/policy_test.exs +5/-2

$ gh api 'repos/astarte-platform/astarte/commits?since=2026-10-05T09:34:00Z&per_page=100' --jq 'length as $n | "count=\($n)", (.[] | "\(.sha[0:8])\t\(.commit.committer.date)\t\(.commit.message | split("\n")[0])")'; echo "--- ancestry of 731bcddc:"; gh api repos/astarte-platform/astarte/compare/731bcddc...master --jq '.status, .ahead_by'
count=12
1da04832	2026-10-06T11:01:17Z	chore: forward port release-1.4
a8a26d0f	2026-10-06T10:59:25Z	chore: forward port release-1.3
d85e0ed0	2026-10-06T10:02:32Z	docs: sync FDO workflow to the current behavior (#2241)
a6bd2c21	2026-10-06T09:59:16Z	fix(fdo): reject owner key names not accepted by OpenBao (#2252)
a51ab2e1	2026-10-06T07:21:18Z	refactor: Update Umbrella checks (#2244)
7628d2ed	2026-10-05T15:54:31Z	test(generators): validate interfaces with Core changesets
ed1e99f4	2026-10-05T15:54:31Z	test(generators): validate mappings with Core changesets
bb17455a	2026-10-05T15:54:31Z	test(generators): validate policy handlers with Core changesets
ec020bf8	2026-10-05T15:54:31Z	fix(generators): generate valid trigger policy names
8f0f5ef1	2026-10-05T14:09:21Z	chore: prepare v1.3.5 (#2245)
c4c63ebe	2026-10-05T09:55:53Z	chore(rm): bump xandra to 0.19.4 (#2246)
731bcddc	2026-10-05T09:34:08Z	chore: forward port release-1.4
--- ancestry of 731bcddc:
ahead
12

$ gh api repos/astarte-platform/astarte/contents/apps?ref=v1.4.0-rc.6 --jq '.[].name' ; echo "== rc.6 tree apps above"; gh api repos/astarte-platform/astarte/contents/apps?ref=v1.3.5 --jq '.[].name'
astarte_appengine_api
astarte_data_updater_plant
astarte_housekeeping
astarte_pairing
astarte_realm_management
astarte_trigger_engine
== rc.6 tree apps above
astarte_appengine_api
astarte_data_updater_plant
astarte_housekeeping
astarte_pairing
astarte_realm_management
astarte_trigger_engine

$ mkdir -p /tmp/opencode && gh api -H "Accept: application/vnd.github.v3.diff" repos/astarte-platform/astarte/commits/1da04832 > /tmp/opencode/1da04832.diff 2>/tmp/opencode/err.txt; echo "exit=$?"; wc -c /tmp/opencode/1da04832.diff; head -3 /tmp/opencode/err.txt
exit=0
3264364 /tmp/opencode/1da04832.diff

$ grep -c '^diff --git' /tmp/opencode/1da04832.diff; grep '^diff --git' /tmp/opencode/1da04832.diff | sed 's|diff --git a/||; s| b/.*||' > /tmp/opencode/files.txt; grep -E '^apps/astarte_(appengine_api|realm_management|pairing|trigger_engine|data_updater_plant|housekeeping)/lib' /tmp/opencode/files.txt
1232
apps/astarte_appengine_api/lib/astarte_appengine_api/application.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/auth/queries.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/device/interface_value.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/device/queries.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/rooms/events_dispatcher.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/rooms/room.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/behaviour.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/delete_volatile_trigger/request_data.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/install_volatile_trigger/request_data.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/vmq_plugin.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/stats/queries.ex
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/device_status_by_alias_controller.ex
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/device_status_by_group_controller.ex
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/device_status_controller.ex
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/fallback_controller.ex
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/interface_values_by_group_controller.ex
apps/astarte_appengine_api/lib/astarte_appengine_api_web/telemetry.ex
apps/astarte_appengine_api/lib/astarte_appengine_api_web/telemetry/api_usage.ex
apps/astarte_appengine_api/lib/astarte_appengine_api_web/telemetry/database_events.ex
apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/client.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/insert_context.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/application.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/config.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_pipeline_supervisor.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/capabilities_handler.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/data_handler.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/device.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/introspection_handler.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/key_agreement/exchange_failed.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/key_agreement/exchange_resp.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/key_agreement/handshake_state.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/key_agreement/hash_ok.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/key_agreement/init_exchange.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/key_agreement/secret_hash.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/key_agreement/shared_secret.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/trigger.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/impl.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/payloads_decoder.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/queries.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/state.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/time_based_actions.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/health.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/rpc/server.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/rpc/server/core.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/trigger_policy/queries.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/triggers_handler.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant_web/telemetry.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant_web/telemetry/database_events.ex
apps/astarte_housekeeping/lib/astarte_housekeeping/application.ex
apps/astarte_housekeeping/lib/astarte_housekeeping/migrator.ex
apps/astarte_housekeeping/lib/astarte_housekeeping/realms/core.ex
apps/astarte_housekeeping/lib/astarte_housekeeping/realms/queries.ex
apps/astarte_housekeeping/lib/astarte_housekeeping/release_task.ex
apps/astarte_housekeeping/lib/astarte_housekeeping_web/telemetry.ex
apps/astarte_housekeeping/lib/astarte_housekeeping_web/telemetry/api_usage.ex
apps/astarte_housekeeping/lib/astarte_housekeeping_web/telemetry/database_events.ex
apps/astarte_pairing/lib/astarte_pairing/engine.ex
apps/astarte_pairing/lib/astarte_pairing/info/info.ex
apps/astarte_pairing/lib/astarte_pairing_web/controllers/fallback_controller.ex
apps/astarte_pairing/lib/astarte_pairing_web/telemetry.ex
apps/astarte_pairing/lib/astarte_pairing_web/telemetry/api_usage.ex
apps/astarte_pairing/lib/astarte_pairing_web/telemetry/database_events.ex
apps/astarte_realm_management/lib/astarte_realm_management/auth/auth.ex
apps/astarte_realm_management/lib/astarte_realm_management/interfaces/create_datastream_individual_multi_interface.ex
apps/astarte_realm_management/lib/astarte_realm_management/interfaces/queries.ex
apps/astarte_realm_management/lib/astarte_realm_management/realm_config/realm_config.ex
apps/astarte_realm_management/lib/astarte_realm_management/triggers/core.ex
apps/astarte_realm_management/lib/astarte_realm_management/triggers/policies/policies.ex
apps/astarte_realm_management/lib/astarte_realm_management_web/api_spec/schemas/mapping.ex
apps/astarte_realm_management/lib/astarte_realm_management_web/plug/verify_header.ex
apps/astarte_realm_management/lib/astarte_realm_management_web/telemetry.ex
apps/astarte_realm_management/lib/astarte_realm_management_web/telemetry/api_usage.ex
apps/astarte_realm_management/lib/astarte_realm_management_web/telemetry/database_events.ex
apps/astarte_trigger_engine/lib/astarte_trigger_engine/amqp_consumer/amqp_consumer_tracker.ex
apps/astarte_trigger_engine/lib/astarte_trigger_engine/amqp_consumer/amqp_message_consumer/amqp_message_consumer.ex
apps/astarte_trigger_engine/lib/astarte_trigger_engine/amqp_consumer/amqp_message_consumer/impl.ex
apps/astarte_trigger_engine/lib/astarte_trigger_engine/config.ex
apps/astarte_trigger_engine/lib/astarte_trigger_engine/consumer_supervisor.ex
apps/astarte_trigger_engine/lib/astarte_trigger_engine/events_consumer.ex
apps/astarte_trigger_engine/lib/astarte_trigger_engine/health.ex
apps/astarte_trigger_engine/lib/astarte_trigger_engine/policy/impl.ex
apps/astarte_trigger_engine/lib/astarte_trigger_engine_web/telemetry.ex
apps/astarte_trigger_engine/lib/astarte_trigger_engine_web/telemetry/database_events.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/volatile_trigger.ex

$ gh api repos/astarte-platform/astarte/contents/apps?ref=master --jq '.[].name'; echo "== libs@master"; gh api repos/astarte-platform/astarte/contents/libs?ref=master --jq '.[].name'
astarte_adapters
astarte_appengine_api
astarte_config
astarte_core
astarte_data_access
astarte_data_updater_plant
astarte_events
astarte_fdo
astarte_fdo_core
astarte_generators
astarte_housekeeping
astarte_pairing
astarte_realm_management
astarte_rpc
astarte_secrets
astarte_test_suite
astarte_trigger_engine
astarte_vmq_plugin
== libs@master
{"message":"Not Found","documentation_url":"https://docs.github.com/rest/repos/contents#get-repository-content","status":"404"}gh: Not Found (HTTP 404)

$ cd /tmp/opencode && cat > extract.sh <<'EOF'
#!/bin/sh
# print the diff hunk for one file path from the big diff
awk -v target="$1" '
  /^diff --git / { p = ($0 == "diff --git a/" target " b/" target) }
  p { print }
' 1da04832.diff
EOF
chmod +x extract.sh
grep -E 'router' files.txt
echo "=== error_view ==="; ./extract.sh apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex | head -60
=== error_view ===
diff --git a/apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex b/apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex
index 8e6a692414..8316b0da7c 100644
--- a/apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex
+++ b/apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex
@@ -98,8 +98,8 @@ defmodule Astarte.AppEngine.APIWeb.ErrorView do
     %{errors: %{detail: "Invalid attributes"}}
   end
 
-  def render("422_unexpected_object_key.json", _assigns) do
-    %{errors: %{detail: "Unexpected object key"}}
+  def render("422_unexpected_object_key.json", %{keys: keys} = _assigns) do
+    %{errors: %{detail: "Unexpected object key", unexpected_keys: keys}}
   end
 
   def render("422_missing_required_mapping.json", _assigns) do

$ cd /tmp/opencode && grep -E '^apps/astarte_core/(specs|lib/astarte_core/interface|lib/astarte_core/mapping|lib/astarte_core/realm)' files.txt; echo "== count lib files under core:"; grep -c '^apps/astarte_core/' files.txt; echo "== fdo/rpc/events present at rc.6 libs?"; gh api repos/astarte-platform/astarte/contents/libs?ref=v1.4.0-rc.6 --jq '.[].name'
apps/astarte_core/lib/astarte_core/interface.ex
apps/astarte_core/lib/astarte_core/interface/aggregation.ex
apps/astarte_core/lib/astarte_core/interface/ownership.ex
apps/astarte_core/lib/astarte_core/interface/type.ex
apps/astarte_core/lib/astarte_core/interface_descriptor.ex
apps/astarte_core/lib/astarte_core/mapping.ex
apps/astarte_core/lib/astarte_core/mapping/database_retention_policy.ex
apps/astarte_core/lib/astarte_core/mapping/endpoints_automaton.ex
apps/astarte_core/lib/astarte_core/mapping/reliability.ex
apps/astarte_core/lib/astarte_core/mapping/retention.ex
apps/astarte_core/lib/astarte_core/mapping/value_type.ex
apps/astarte_core/lib/astarte_core/realm.ex
apps/astarte_core/specs/interface.json
apps/astarte_core/specs/mapping.json
apps/astarte_core/specs/policy.json
== count lib files under core:
118
== fdo/rpc/events present at rc.6 libs?
astarte_adapters
astarte_config
astarte_data_access
astarte_events
astarte_fdo
astarte_fdo_core
astarte_generators
astarte_rpc
astarte_secrets
astarte_test_suite

$ cd /tmp/opencode && for f in apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex apps/astarte_appengine_api/lib/astarte_appengine_api/rooms/room.ex apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/fallback_controller.ex apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/interface_values_by_group_controller.ex apps/astarte_appengine_api/lib/astarte_appengine_api/device/interface_value.ex; do n=$(echo $f | tr '/' '_'); ./extract.sh $f > "p_$n.diff"; wc -l "p_$n.diff"; done; ls -la p_*.diff
491 p_apps_astarte_appengine_api_lib_astarte_appengine_api_device_device.ex.diff
251 p_apps_astarte_appengine_api_lib_astarte_appengine_api_rooms_room.ex.diff
18 p_apps_astarte_appengine_api_lib_astarte_appengine_api_web_controllers_fallback_controller.ex.diff
116 p_apps_astarte_appengine_api_lib_astarte_appengine_api_web_controllers_interface_values_by_group_controller.ex.diff
33 p_apps_astarte_appengine_api_lib_astarte_appengine_api_device_interface_value.ex.diff
-rw-r--r-- 1 root root 15964 Oct  6 19:33 p_apps_astarte_appengine_api_lib_astarte_appengine_api_device_device.ex.diff
-rw-r--r-- 1 root root  1314 Oct  6 19:33 p_apps_astarte_appengine_api_lib_astarte_appengine_api_device_interface_value.ex.diff
-rw-r--r-- 1 root root  9268 Oct  6 19:33 p_apps_astarte_appengine_api_lib_astarte_appengine_api_rooms_room.ex.diff
-rw-r--r-- 1 root root   941 Oct  6 19:33 p_apps_astarte_appengine_api_lib_astarte_appengine_api_web_controllers_fallback_controller.ex.diff
-rw-r--r-- 1 root root  4575 Oct  6 19:33 p_apps_astarte_appengine_api_lib_astarte_appengine_api_web_controllers_interface_values_by_group_controller.ex.diff

→ Read /tmp/opencode/p_apps_astarte_appengine_api_lib_astarte_appengine_api_device_device.ex.diff
→ Read /tmp/opencode/p_apps_astarte_appengine_api_lib_astarte_appengine_api_web_controllers_fallback_controller.ex.diff

$ cd /tmp/opencode && ./extract.sh apps/astarte_appengine_api/lib/astarte_appengine_api/device/queries.ex > p_queries.diff; wc -l p_queries.diff; grep -n "encrypt\|Encrypt" p_queries.diff | head -30; echo "=== device.ex encrypt ctx in parent? ==="; grep -c "encrypt" p_queries.diff
557 p_queries.diff
27:+        select: [:value_type, :endpoint_id, :encrypted]
36:+        select: [:value_type, :endpoint, :encrypted]
46:+          :encrypted
319:+      if endpoint_rows |> Enum.any?(& &1.encrypted) do
320:+        [:encrypted_dek | columns]
343:+      if Map.get(endpoint_row, :encrypted) do
344:+        {:encryptedblob_value, [:encrypted_dek]}
374:+    query = from Endpoint, select: [:value_type, :encrypted]
385:+      if Map.get(endpoint_row, :encrypted) do
386:+        {:encryptedblob_value, [:encrypted_dek]}
400:+    if Enum.any?(rows, & &1.encrypted) do
412:+      if row.encrypted do
428:+        key = Map.fetch!(keys, value.encrypted_dek)
437:+  defp maybe_decrypt_values(values, %{encrypted: true}, realm_name, value_type) do
445:+  defp maybe_decrypt_result({:ok, row}, %{encrypted: true}, realm_name) do
446:+    with {:ok, dek} <- fetch_dek(row.encrypted_dek, realm_name),
447:+         {:ok, bin} <- Secrets.decrypt_with_dek(row.encrypted_value, dek) do
478:+    |> Enum.map(& &1.encrypted_dek)
480:+    |> Enum.reduce_while({:ok, %{}}, fn encrypted_dek, {:ok, cache} ->
481:+      case fetch_dek(encrypted_dek, realm_name) do
482:+        {:ok, dek} -> {:cont, {:ok, Map.put(cache, encrypted_dek, dek)}}
488:+  defp fetch_dek(encrypted_dek, realm_name) do
491:+        Secrets.unwrap_dek(key.name, encrypted_dek, key.namespace)
499:+    key = Map.fetch!(keys, row.encrypted_dek)
500:+    {:ok, bin} = Secrets.decrypt_with_dek(row.encryptedblob_value, key)
526:+        {:error, :device_not_ready_for_encryption}
=== device.ex encrypt ctx in parent? ===
26

$ cd /tmp/opencode && gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex?ref=a8a26d0f --jq '.content' | base64 -d > parent_device.ex; grep -n "encrypt\|Encrypt\|unexpected_keys\|wrap_to_bson_struct" parent_device.ex | head -20; echo "== parent queries:"; gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api/device/queries.ex?ref=a8a26d0f --jq '.content' | base64 -d > parent_queries.ex; grep -n "encrypt\|Encrypt" parent_queries.ex | head -20; echo "(no output above = none)"
205:         wrapped_value = wrap_to_bson_struct(mapping.value_type, value),
414:         wrapped_value = wrap_to_bson_struct(expected_types, value),
682:  defp wrap_to_bson_struct(:binaryblob, value) do
687:  defp wrap_to_bson_struct(:binaryblobarray, values) do
688:    Enum.map(values, &wrap_to_bson_struct(:binaryblob, &1))
691:  defp wrap_to_bson_struct(expected_types, values)
696:      {key, wrap_to_bson_struct(type, value)}
701:  defp wrap_to_bson_struct(_anytype, value) do
== parent queries:
(no output above = none)

$ cd /tmp/opencode && for ref in v1.4.0-rc.6 v1.4.0-rc.5; do echo "== $ref device.ex"; gh api "repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex?ref=$ref" --jq '.content' 2>/dev/null | base64 -d > d_$ref.ex; grep -c "maybe_apply_transport_encryption\|unexpected_keys" d_$ref.ex; echo "== $ref error_view"; gh api "repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex?ref=$ref" --jq '.content' | base64 -d | grep -n "unexpected_object_key" ; done
== v1.4.0-rc.6 device.ex
0
== v1.4.0-rc.6 error_view
101:  def render("422_unexpected_object_key.json", _assigns) do
== v1.4.0-rc.5 device.ex
0
== v1.4.0-rc.5 error_view
101:  def render("422_unexpected_object_key.json", _assigns) do

$ gh api 'repos/astarte-platform/astarte/compare/v1.4.0-rc.6...1da04832' --jq '"status=\(.status) ahead=\(.ahead_by) commits=\(.commits|length) files=\(.files|length)"' 2>&1 | head -5; echo "== was astarte_core present in rc.6 anywhere?"; gh api 'search/code?q=repo:astarte-platform/astarte+filename:interface.json+path:specs' --jq '.items[].path' 2>&1 | head
status=ahead ahead=1510 commits=250 files=300
== was astarte_core present in rc.6 anywhere?

$ gh api 'repos/astarte-platform/astarte/compare/v1.4.0-rc.6...master' --jq '{status, ahead_by, behind_by, merge_base: .merge_base_commit.sha[0:8], merge_base_date: .merge_base_commit.commit.committer.date}'; echo "== branches"; gh api repos/astarte-platform/astarte/branches --jq '.[].name' | head -20; echo "== tags list top"; gh api repos/astarte-platform/astarte/tags?per_page=15 --jq '.[].name'
{"ahead_by":1510,"behind_by":0,"merge_base":"20df0d32","merge_base_date":"2026-09-29T10:56:18Z","status":"ahead"}
== branches
fdo-pin
feature/erlang-rpc
feature/merge-backend-and-api-apps
graphQL
master
release-0.10
release-0.11
release-1.0
release-1.1
release-1.2
release-1.3
release-1.4
release-please--branches--master
renovate/elixir-1.x
renovate/elixir-dependencies-(non-major)
renovate/erlang-26.x
renovate/erlang-28.x
renovate/github-actions-(non-major)
renovate/lock-file-maintenance
renovate/major-elixir-dependencies
== tags list top
v1.4.0-rc.6
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
v1.4.0-rc.0
v1.3.5
v1.3.4
v1.3.3
v1.3.2
v1.3.1
v1.3.0
v1.3.0-rc.2
v1.3.0-rc.1

$ grep -n "encrypted\|unexpected.object.key\|unexpected_keys" /root/astrate-mule/.mule/todo.md | head -20; echo "=== for-giulio mentions:"; grep -n "unexpected_keys\|encrypted" /root/astrate-mule/.mule/for-giulio.md | head -20
1:- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
138:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
139:- [x] detailed-listing-required-encrypted-flags: emit the `required`/`encrypted` mapping flags in `writeDetailedMapping` (internal/realm/detailed.go) so `GET /realmmanagement/v1/{realm}/interfaces?detailed=true` round-trips them — the detailed listing currently materializes a fixed field set that drops both flags even when a datastream mapping declares one, while the plain GET preserves them via the stored definition. Datastream branch only (properties reject both at parse); add a T1 assertion in internal/realm/detailed_test.go.
159:- [x] docs-sync-rm-mapping-required-encrypted [auto]: add `required` (boolean, datastream-only) and `encrypted` (boolean, datastream-only) to the `InterfaceMapping` schema in docs/api/astarte_realm_management_api.yaml — the parser accepts and persists both on datastream mappings (internal/interfaceschema/parse.go:410-414), rejects them on properties mappings (parse.go:395-398), the detailed listing emits them (`detailed.go:69`), and the compatibility check compares them (interfaceschema/compat.go:121-123); the schema (yaml:820-862) currently has `additionalProperties: false` and neither field, so the OpenAPI would mark a valid datastream interface install body as invalid. Also add both fields to the detailed-listing example (yaml:70-85). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
162:- [x] interfaceschema-compat-attrs-flip-test [auto]: extend TestCheckMinorUpgradeTable (pkg/interfaceschema/compat_test.go) to assert the four currently-unpinned immutable attributes — retention, expiry, database_retention_policy/ttl (datastream old/next pair) and allow_unset (properties old/next pair) — each flipped on upgrade and rejected as ErrIncompatibleEndpointChange with an at-immutability acceptance twin; only type/reliability/explicit_timestamp/required/encrypted are tested today (sameMappingAttributes, compat.go:102, sits at 66.7% stmt coverage), so a regression relaxing expiry or DB-TTL immutability passes the suite.
265:- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
266:- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress. — BLOCKED: TIMEOUT after 1200s — task too big, split it
267:- [x] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
=== for-giulio mentions:
21:  **Standing item unchanged.** `docs/UPSTREAM-EXPERIMENTAL.md`'s two rows (#67 `required`/`encrypted` mapping fields, #68 `async_operation=false`) still wait on **upstream 1.4 final** for their promoted-vs-deprecated call; nothing in today's commit set promotes or removes either and none of it is released, so the register is untouched.
38:  **The standing item is unchanged and I am not re-litigating it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries its two rows, both tagged `1.4 experimental` with trigger "upstream 1.4 final: promoted or removed" — the mis-tagging the 2026-10-04 entry measured out of upstream's CHANGELOG (row #68's `async_operation` feature is 1.0.2 and stable for four years; row #67's `required` half is a plain 1.4.0-rc.0 addition, only its `encrypted` half being genuinely experimental and already tracked by #92/#93). v1.3.5 touches neither feature, so it changes nothing about your four pending decisions (a)–(d) in that entry, and I have not edited the register.
61:  - **Row #67 is half right.** Its `required` half is also a plain addition, not experimental: `CHANGELOG.md:113`, under `## [1.4.0-rc.0] - 2026-04-08` (header at :109), "Add required flag for mappings of object aggregated interfaces". Its `encrypted` half is the keyAgreement coupling, which *is* genuinely experimental and is already tracked where it belongs — **#51/#92** (parked) and **#93** — so the experimental half of row #67 is real, just filed in the wrong place.
64:  **Yours, four small decisions, none of them a code gap:** (a) row #68 — retag it `1.0.2, stable` with reconcile trigger "none, keep as a deviation reference", or **drop the row** and let deviation 17 be the only record, which is the cleaner reading of a register that exists to track *unstable upstream features*; (b) row #67 — split it, leaving `required` untagged (or dropped) and letting the `encrypted` half be tracked by #92/#93 alone, or keep one row that names only the keyAgreement-coupled part; (c) whether deviation 17's opening sentence should say 1.0.2 rather than 1.4, and its device-delete clause should say "extension" rather than "accepted and ignored"; (d) whether the register should gain a line saying what "reconcile" means when a feature turns out never to have been experimental, so the next run is not left guessing as this one was. I did not pick any of them and I have changed nothing: no `gh issue create`, nothing commented, closed or edited on GitHub, and **no `.mule/todo.md` lines** — each of these edits a file that carries your decisions, and queueing them would put your choice in front of the queue as if it were settled.
80:  **The recipe's standing item is satisfied and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental*, reconcile trigger "upstream 1.4 final". 1.4 is not final, so there is nothing to promote and nothing to deprecate, and the register is accurate as written. `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked this run) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision ("bumped only in the same change that completes that surface, after reconciling every row tagged with the level"). The milestone-bump issue the recipe asks me to file **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": v3.0's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is delivered, the 1.4 half is open (#92, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut.
137:  **The recipe's standing item is satisfied, and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still has exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental* with reconcile trigger "upstream 1.4 final". v1.4 is not final, so both rows stay as written: nothing to promote, nothing to deprecate, and the register is accurate. `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision; the milestone-bump issue the recipe asks for **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": the milestone's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is done, and the 1.4 half is unresolved (#92 open, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut. Also still true and already escalated on 2026-09-06, so not repeated: the section's source document `.mule/research/upstream-parity-2026-08-22.md` does not exist in the repo.
210:- **Milestone recipe run, 2026-09-24: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-23 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
218:- **Milestone recipe run, 2026-09-23: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-21 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
234:- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
242:- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
254:- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
258:- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
274:- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) shows a new **stable v1.3.4** (2026-09-18) — maintenance only (data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change), so no new gap for Astrate — and still **no stable v1.4.0**: newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
286:- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-13 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
290:- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-12 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
294:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
302:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
304:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
338:  UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait

$ cd /root/astrate-mule && grep -rn "encrypted" --include=*.go internal/ | grep -iv "encrypted_dek\|//" | head -15; echo "=== UPSTREAM-EXPERIMENTAL rows:"; sed -n '18,32p' docs/UPSTREAM-EXPERIMENTAL.md
internal/realm/detailed.go:69:	fmt.Fprintf(b, `,"expiry":%d,"explicit_timestamp":%t,"required":%t,"encrypted":%t,"database_retention_policy":`,
internal/realm/detailed_required_encrypted_test.go:15:		`"mappings":[{"endpoint":"/value","type":"double","required":true,"encrypted":true}]}`)
internal/realm/detailed_required_encrypted_test.go:24:		`"required":true,"encrypted":true,` +
=== UPSTREAM-EXPERIMENTAL rows:
remembered.

## Rows

| Feature | Issue | Upstream level | Adopted | Reconcile when |
|---|---|---|---|---|
| Required + encrypted interface mapping fields | #67 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
| `async_operation=false` on mutating RM endpoints | #68 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |

$ cd /root/astrate-mule && grep -rn "Unexpected object key\|unexpected_keys" --include=*.go internal/ pkg/ | head -10; echo "=== encrypt applied?"; grep -rln "SharedSecret\|EncryptValue\|encrypted_dek" --include=*.go internal/ pkg/ | head
internal/appengine/writeerror_test.go:38:			http.StatusBadRequest, "Unexpected object key", []string{"ghost"}},
internal/appengine/writeerror_test.go:43:			http.StatusBadRequest, "Unexpected object key", []string{"a", "b"}},
internal/appengine/writeerror_test.go:64:					Keys   []string `json:"unexpected_keys"`
internal/appengine/writeerror_test.go:74:				t.Fatalf("unexpected_keys = %v, want %v", body.Errors.Keys, tc.keys)
internal/appengine/writeerror_test.go:78:					t.Errorf("unexpected_keys[%d] = %q, want %q", i, body.Errors.Keys[i], k)
internal/appengine/writeerror_test.go:87:// {"errors":{"detail":"Unexpected object key","unexpected_keys":["ghost"]}}.
internal/appengine/writeerror_test.go:96:	want := `{"errors":{"detail":"Unexpected object key","unexpected_keys":["ghost"]}}`
internal/appengine/http.go:711:const detailUnexpectedObjectKey = "Unexpected object key"
internal/appengine/http.go:723:// upstream envelope: {"errors": {"detail": "Unexpected object key",
internal/appengine/http.go:724:// "unexpected_keys": [...]}}. The other object-shape failures — not a
=== encrypt applied?

$ cd /tmp/opencode && gh api repos/astarte-platform/astarte/contents/apps/astarte_core/specs/mapping.json?ref=master --jq '.content' | base64 -d > mapping.json; gh api repos/astarte-platform/astarte/contents/apps/astarte_core/specs/interface.json?ref=master --jq '.content' | base64 -d > interface.json; python3 -c "
import json
for f in ('mapping.json','interface.json'):
    d=json.load(open(f))
    print(f, '->', list(d.get('properties', d).keys()) if isinstance(d,dict) else type(d))
" 2>&1 | head; echo "== raw keys:"; python3 -c "
import json
for f in ('mapping.json','interface.json'):
    d=json.load(open(f))
    print(f, json.dumps(d, indent=0)[:600])
"
mapping.json -> ['endpoint', 'type', 'reliability', 'retention', 'expiry', 'database_retention_policy', 'database_retention_ttl', 'allow_unset', 'explicit_timestamp', 'description', 'doc', 'encrypted']
interface.json -> ['interface_name', 'version_major', 'version_minor', 'type', 'ownership', 'aggregation', 'description', 'doc', 'mappings']
== raw keys:
mapping.json {
"$schema": "http://json-schema.org/draft-04/schema",
"title": "Astarte Mapping Schema",
"type": "object",
"description": "Identifies a mapping for an interface. A mapping must consist at least of an endpoint and a type.",
"properties": {
"endpoint": {
"type": "string",
"pattern": "^(/(%{([a-zA-Z_][a-zA-Z0-9_]*)}|[a-zA-Z_][a-zA-Z0-9_]*)){1,64}$",
"minLength": 2,
"maxLength": 256,
"description": "The template of the path. This is a UNIX-like path (e.g. /my/path) and can be parametrized. Parameters are in the %{name} form, and can be used to create interfaces which represent dictionaries of map
interface.json {
"$schema": "http://json-schema.org/draft-04/schema",
"title": "Astarte Interface Schema",
"description": "This schema describes how an Astarte interface should be declared",
"type": "object",
"properties": {
"interface_name": {
"type": "string",
"pattern": "^([a-zA-Z][a-zA-Z0-9]*.([a-zA-Z0-9][a-zA-Z0-9-]*.)*)?[a-zA-Z][a-zA-Z0-9]*$",
"minLength": 1,
"maxLength": 128,
"description": "The name of the interface. This has to be an unique, alphanumeric reverse internet domain name, shorther than 128 characters."
},
"version_major": {
"type": "integer",
"description": "A Major version qualifier for

$ cd /tmp/opencode && awk '/^diff --git /{f=$3; sub("^a/","",f)} /^\+/ && tolower($0) ~ /encrypt/ {print f}' 1da04832.diff | sort | uniq -c | sort -rn | head -25
45 apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_test.exs
     41 apps/astarte_secrets/lib/astarte_secrets/astarte_secrets.ex
     31 apps/astarte_data_updater_plant/test/support/cases/device.ex
     31 apps/astarte_appengine_api/test/support/cases/device.ex
     28 apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex
     26 libs/astarte_secrets/test/astarte_secrets/astarte_secrets_test.exs
     26 apps/astarte_secrets/test/astarte_secrets/encrypted_messages_test.exs
     26 apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/control_handler_test.exs
     26 apps/astarte_appengine_api/lib/astarte_appengine_api/device/queries.ex
     23 apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_error_test.exs
     19 apps/astarte_appengine_api/test/astarte_appengine_api/device/device_reading_v2_test.exs
     18 apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex
     16 apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/data_handler.ex
     16 apps/astarte_core/test/astarte_core/mapping_test.exs
     15 doc/pages/architecture/082-key_agreement_protocol.md
     13 libs/astarte_data_access/lib/astarte_data_access/device.ex
     12 apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex
     10 libs/astarte_secrets/lib/astarte_secrets/core.ex
     10 apps/astarte_vmq_plugin/rel/overlays/bin/start_vernemq
      8 libs/astarte_data_access/test/mappings_test.exs
      8 apps/astarte_secrets/test/astarte_secrets/data_encryption_key_cache_test.exs
      8 apps/astarte_realm_management/lib/astarte_realm_management/interfaces/queries.ex
      7 libs/astarte_secrets/test/astarte_secrets/core_test.exs
      7 astarte-dashboard/src/components/MappingEditor.tsx
      6 apps/astarte_secrets/lib/astarte_secrets/data_encryption_key_cache.ex

$ cd /tmp/opencode && awk '/^diff --git /{f=$3; sub("^a/","",f)} /^-/ && tolower($0) ~ /encrypt/ {print f}' 1da04832.diff | sort | uniq -c | sort -rn | head -25; echo "=== any pre-existing encrypt in parent DUP:"; gh api repos/astarte-platform/astarte/contents/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/impl.ex?ref=a8a26d0f --jq '.content' | base64 -d | grep -ci encrypt
1 libs/astarte_secrets/lib/astarte_secrets/astarte_secrets.ex
=== any pre-existing encrypt in parent DUP:
0

$ gh api 'repos/astarte-platform/astarte/contents/libs/astarte_secrets/lib/astarte_secrets?ref=v1.4.0-rc.6' --jq '.[].name' 2>&1; echo "== rc.6 DUP data_handler encrypt:"; gh api 'repos/astarte-platform/astarte/contents/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/data_handler.ex?ref=v1.4.0-rc.6' --jq '.content' 2>/dev/null | base64 -d | grep -ci encrypt; echo "== rc.6 mapping struct encrypted field:"; gh api 'search/code?q=repo:astarte-platform/astarte+encrypted+in:file+filename:mapping.ex' --jq '.items[].path' 2>&1 | head -5
astarte_secrets.ex
client.ex
config
core.ex
key
owner_key_initialization.ex
owner_key_initialization_options.ex
== rc.6 DUP data_handler encrypt:
0
== rc.6 mapping struct encrypted field:
apps/astarte_core/lib/astarte_core/mapping.ex
apps/astarte_realm_management/lib/astarte_realm_management_web/api_spec/schemas/mapping.ex

$ grep -n "probe-" /root/astrate-mule/.mule/todo.md | head -20; echo "=== log.md probe entries:"; grep -n "probe-" /root/astrate-mule/.mule/log.md | head -10
114:- [!] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages" — investigate only: does Astrate have an analogous delay between trigger install and first delivery? Report, do not patch. — BLOCKED: wrote nothing
119:- [!] probe-property-resend-encoding [auto]: upstream v1.3.0 fixed outbound server-property values sent to a device on connect/emptyCache (commit 522ccf4f — a raw stored binaryblob is now wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding, so it ships as a BSON binary, not a string). Investigate only: does Astrate's resendServerProperties/rehydrate path (internal/engine/control.go:144, pkg/payload) emit binaryblob → BSON subtype-0 binary and datetime → UTC datetime exactly like upstream 1.3? Report, do not patch. — BLOCKED: wrote nothing
126:- [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
137:- [x] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
=== log.md probe entries:
17:| 2026-09-05 | probe-trigger-install-notification-delay | blocked | 308s | wrote nothing |
22:| 2026-09-05 | probe-property-resend-encoding | blocked | 100s | wrote nothing |
29:| 2026-09-06 | probe-props-resend-error-triggers | blocked | 331s | wrote nothing |
40:| 2026-09-08 | probe-required-mapping-flag | done | 426s | db9d06e |

$ grep -n "^## \|^# " /root/astrate-mule/.mule/todo.md; echo "=== tail:"; tail -25 /root/astrate-mule/.mule/todo.md | cut -c1-160
5:# Mule queue
37:## Where tasks come from
70:## Nothing merges on its own
78:## Queue
=== tail:
- [x] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstrained
- [!] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_realm_management_api.yaml — `data` is de
- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CL
- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` shows `d
- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDetail
- [x] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError`
- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec st
- [!] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implem
- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undecl
- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data wr
- [x] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mappin
- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z
- [x] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362)
- [!] docs-sync-hk-patch-422-field-error-examples [auto]: add the three missing PATCH 422 bodies to the `ValidationErrors` response component in docs/api/astart
- [!] docs-sync-hk-wrong-type-field-400 [auto]: widen the `BadRequest` component description in docs/api/astarte_housekeeping_api.yaml (yaml:378, "malformed or 
- [x] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `compon
- [!] store-register-inhibit-preserve: RegisterDevice must not overwrite status to 'registered' when a device was inhibited before first credentials request (pr
- [x] store-aliasvalues-self-exclusion-test: Add test for AliasValuesTaken to assert self-exclusion (same device's own alias values do not count as taken). File
- [!] store-latestindividual-empty-errnotfound: Add store-level test asserting LatestIndividual returns ErrNotFound for empty series and returns newest on multi
- [x] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.
- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is th
- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises 
- [x] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response co
- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, nul
- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName`

→ Read .mule/todo.md [offset=268]

$ grep -c "Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run .make -C docs build. and confirm the Swagger-referenced YAML still loads." /root/astrate-mule/.mule/todo.md
1

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -278,4 +278,6 @@
 - [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub.
+- [ ] compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (newest stable, 2026-10-05, maintenance-only: mix.exs/mix.lock version bumps plus a single CHANGELOG "astarte_data_updater_plant RPC server reliability" entry, no wire/API surface change; v1.4.0 is still rc.6-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2/v1.3.4 wording proposal there already stands and its version-reference delta was recorded on 2026-10-05 as reading "newest stable v1.3.5", so fold this in as that one-word delta rather than re-deriving the proposal.

$ tail -3 /root/astrate-mule/.mule/todo.md | cut -c1-130; echo "== total lines:"; wc -l /root/astrate-mule/.mule/todo.md
- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `com
- [ ] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 fil
- [ ] compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (newest stable, 2026-10-05, maintenance-only
== total lines:
283 /root/astrate-mule/.mule/todo.md

Done: watched upstream Astarte — no new release since v1.3.5 (measured wire-inert, as this recipe assessed it on 2026-10-05); the one wire-bearing master delta is today's unreleased forward-port (1da04832) that turns the `encrypted` mapping flag into enforcement, tracked by parked row #67/#92 — queued a probe line plus the recipe's mandated compat-note line.

Files: .mule/todo.md (two lines appended)

Verified: `gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] ...'` (newest stable v1.3.5 2026-10-05, newest overall rc.6) -> pass; `gh api repos/astarte-platform/astarte/commits?since=2026-10-05T09:34:00Z` = 12 commits, nothing for Astrate except 1da04832 -> pass; greps on rc.6/v1.3.5 tags for "encrypt"/`unexpected_keys` (both absent from rc.6, absent from master-before) -> pass; `internal/appengine/{http.go,writeerror_test.go}` already implement `unexpected_keys`; `docs/COMPATIBILITY.md` unchanged (target 1.2.2, never-touch).

Unsure: whether the compat-note-v1.3.5 line is redundant — the 2026-10-05 run deliberately skipped it because the v1.3.2/v1.3.4 wording proposal in for-giulio.md already records the v1.3.5 version-reference delta; I queued it anyway per the recipe's unconditional "always propose, as the last line", worded to fold rather than re-derive. Also unverifiable on this Mac (no Docker): live wire behaviour of the encrypted-mapping path — the probe line covers that.

Follow-ups: `.mule/todo.md` + `probe-encrypted-mapping-enforcement [auto]`, `compat-note-v1.3.5 [auto]` (no git touched, no GitHub state changed — the recipe asks for neither).
