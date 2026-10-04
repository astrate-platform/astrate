slug: recipe-astarte-upstream
verdict: proposed
at:  2f44477
ran: 2026-10-04T18:40:52Z on DietPi in 103s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"' 2>&1
v1.4.0-rc.6	2026-09-30T07:28:27Z
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z

$ grep -rn "compat-note-\|1\.4\.0-rc\|v1\.4\.0\|v1\.3\.4\|UPSTREAM-EXPERIMENTAL" .mule/ docs/ --include=*.md -l 2>&1 | head -20
.mule/log.md
.mule/for-giulio.md
.mule/milestones.md
.mule/recipes/milestones.md
.mule/recipes/astarte-upstream.md
.mule/reports/log/20260918T113930Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260923T181743Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20261002T181712Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260924T201012Z-recipe-code-review-timeout.md
.mule/reports/log/20260928T170059Z-docs-sync-pairing-deviceid-base64url-done.md
.mule/reports/log/20260913T171327Z-recipe-astarte-upstream-noop.md
.mule/reports/log/20260917T194954Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260920T185539Z-compat-note-custom-action-policy-boundary-done.md
.mule/reports/log/20261003T180922Z-recipe-github-issues-proposed.md
.mule/reports/log/20260908T164748Z-compat-note-v14-rc-blocked.md
.mule/reports/log/20260915T170704Z-hk-zero-retention-asymmetry-done.md
.mule/reports/log/20260925T201910Z-docs-sync-rm-delete-device-async-operation-param-done.md
.mule/reports/log/20260921T105709Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md
.mule/reports/log/20260915T111354Z-recipe-astarte-upstream-proposed.md

$ ls -la .mule/ && wc -l .mule/todo.md .mule/for-giulio.md 2>&1
total 12068
drwxr-xr-x  7 root root     4096 Oct  4 20:11 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       14 Oct  4 20:39 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  4 19:38 .heartbeat
-rw-r--r--  1 root root        0 Oct  4 20:39 .last-output
-rw-r--r--  1 root root        0 Oct  4 20:39 .lock
-rw-r--r--  1 root root       17 Oct  4 20:39 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 11977352 Oct  4 20:39 cron.log
drwxr-xr-x  2 root root     4096 Oct  4 19:22 failed
-rw-r--r--  1 root root   127766 Oct  4 20:31 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    13845 Oct  4 19:38 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  2 13:43 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   141312 Oct  4 19:38 todo.md
   272 .mule/todo.md
   604 .mule/for-giulio.md
   876 total

$ grep -n "astarte-upstream" .mule/log.md | tail -20
(no output)

$ ls -1 .mule/reports/log/ | grep astarte-upstream | sort | tail -12
20260927T203020Z-recipe-astarte-upstream-timeout.md
20260928T020420Z-recipe-astarte-upstream-proposed.md
20260928T203150Z-recipe-astarte-upstream-proposed.md
20260929T165352Z-recipe-astarte-upstream-proposed.md
20260929T180551Z-recipe-astarte-upstream-proposed.md
20260929T194353Z-recipe-astarte-upstream-proposed.md
20260930T171256Z-recipe-astarte-upstream-proposed.md
20260930T191053Z-recipe-astarte-upstream-proposed.md
20261001T111852Z-recipe-astarte-upstream-proposed.md
20261001T184353Z-recipe-astarte-upstream-proposed.md
20261001T203619Z-recipe-astarte-upstream-proposed.md
20261003T183445Z-recipe-astarte-upstream-proposed.md

→ Read .mule/reports/log/20261003T183445Z-recipe-astarte-upstream-proposed.md
Previous run (2026-10-03) already handled rc.6. Let me verify the state of the queue before concluding.

$ grep -n "compat-note-\|fdo-rc6\|v1\.4\.0-rc" .mule/todo.md
3:- [x] compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).
115:- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
120:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
137:- [x] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
138:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
177:- [x] compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.
247:- [x] fdo-rc6-scope-delta-for-giulio [auto]: upstream v1.4.0-rc.6 (2026-09-30, 27 commits, 154 files, almost all FDO) grows the onboarding surface issue #78 is scoped against, and none of it is in #78's body (written against upstream 1.3 / rc.5): new `POST /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}/to0` re-runs TO0 and returns the refreshed rendezvous expiry (`{"data":{"expiry":…}}`, ownership_voucher.ex:63-73; 409 `device_already_onboarded` via fallback_controller.ex:120-125); voucher deletion, which revokes the rendezvous registration first and *keeps* the voucher if that fails (`:rendezvous_revocation_failed` → 500) and refuses another realm's voucher through `ensure_voucher_in_realm` (ownership_voucher.ex:40-60 — now expressible because the voucher row carries its realm); `device_id` accepted on voucher upload; the device registered immediately on upload rather than mid-TO2 (`add_unconfirmed_credentials`, device.ex:63-85); single-entry x5chain certs accepted; invalid vouchers no longer stored; vouchers and to2-sessions moved from the per-realm keyspace to the global `astarte` one with `fdo_guid` added to the device row (housekeeping/realms/queries.ex:784, 1129-1133); and RM device deletion now cascading to the voucher (device_removal/core.ex). Write the delta into `.mule/for-giulio.md` as a scope update for #78 and say which of these a client could actually observe — do not edit or comment on the issue (GitHub state is read-only for the mule) and do not implement anything.
248:- [x] compat-note-v1.4-rc.6 [auto]: record in `.mule/for-giulio.md` that v1.4.0-rc.6 is wire-inert for Astrate, as a delta on the open v1.3.2/v1.3.4 wording proposal rather than a re-derivation. Measured, the whole non-FDO delta is inert: DUP RPC-server reliability is OTP supervisor/GC tuning; the Dashboard's new `unknown_status` is computed client-side from `last_connection && !last_disconnection && !isConnected` (models/Device/index.ts:213-215) and both timestamps are already on Astrate's device shape (internal/appengine/service.go:118), so a device whose disconnection was never recorded reads as "Unknown status" rather than "Disconnected"; the `on_device_registered` dispatch moved from `Agent.register_device` to `Engine.register_device` (pairing/engine.ex:157-165, 186-193) so it also covers the FDO path, but the Agent path is byte-for-byte the same rule and Astrate already fires it on every successful registration (internal/pairing/service.go:221-223); the `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"); and `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma. So the proposal's "v1.4.0 is still rc.5-only" reference becomes rc.6, and its "FDO authentication (pairing, disabled by default)" line understates the 1.4 line, which ships a full ownership-voucher surface Astrate does not emulate (#78, parked). Do not edit docs/COMPATIBILITY.md and do not restate the standing proposal.
264:- [!] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implements `device_empty_cache_received` end to end — `SimpleEvent` field 23, an empty `DeviceEmptyCacheReceivedEvent` whose encoder emits exactly `{"type":"device_empty_cache_received"}` (astarte_core encoder.ex:54-62), fired by DUP at `control_handler.ex:161` *inside* the `with` that runs `send_control_consumer_properties` → `resend_all_properties` → `set_pending_empty_cache`, so it fires only when all three succeed and never on failure, and the commit deletes the `# TODO: implement on_empty_cache_received` from `astarte_events/lib/astarte_events/triggers/core.ex`. Astrate already accepts the trigger but never evaluates it: `OnDeviceEmptyCacheReceived` sits in `deviceOns` with value `false` (internal/engine/triggers/match.go:123, the `evaluated` flag read at match.go:406) and the comment's stated reason ("upstream defines no SimpleEvent variant for it") is now false. Add the event body constructor to internal/engine/triggers/events.go (the payload is `{"type":"device_empty_cache_received"}` and nothing else — no device id, no metadata), flip the `deviceOns` value to true, and fire it from `handleEmptyCache` (internal/engine/control.go:58) with the same idiom `handleLifecycle` uses for `device_connected` (`e.fireDevice` + `stream.Publish`, internal/engine/engine.go:346-353), at the point where `m.Ack()` runs — i.e. after `resendServerProperties` and `sendConsumerProperties` have both succeeded inside `retryStore`. Add a T1 test that a trigger on it fires exactly once per accepted emptyCache and does not fire when the resend fails. — BLOCKED: wrote nothing

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | "\(.tag_name)\t\(.prerelease)\t\(.published_at)"' 2>&1 | head -8
v1.4.0-rc.6	true	2026-09-30T07:28:27Z
v1.3.4	false	2026-09-18T09:42:57Z
v1.4.0-rc.5	true	2026-08-20T12:28:04Z
v1.3.3	false	2026-08-07T14:33:58Z
v1.4.0-rc.4	true	2026-08-06T13:40:20Z
v1.4.0-rc.3	true	2026-07-31T10:48:39Z
v1.4.0-rc.2	true	2026-07-14T10:16:22Z
v1.4.0-rc.1	true	2026-07-14T08:53:27Z

Release list is byte-identical to the 2026-10-03 run, and the rc.6 items are already delivered. One cheap check for untagged movement since then:
$ gh api "repos/astarte-platform/astarte/compare/v1.4.0-rc.6...main" --jq '"ahead_by=\(.ahead_by) behind_by=\(.behind_by) files=\(.files|length) commits=\(.commits|length)", (.commits[] | "\(.sha[0:8]) \(.commit.author.date) \(.commit.message | split("\n")[0])")' 2>&1 | head -40
{"message":"Not Found","documentation_url":"https://docs.github.com/rest/commits/commits#compare-two-commits","status":"404"}gh: Not Found (HTTP 404)

$ gh api repos/astarte-platform/astarte --jq '.default_branch' 2>&1; gh api "repos/astarte-platform/astarte/compare/v1.4.0-rc.6...master" --jq '"ahead_by=\(.ahead_by) behind_by=\(.behind_by) files=\(.files|length) commits=\(.commits|length)", (.commits[] | "\(.sha[0:8]) \(.commit.author.date) \(.commit.message | split("\n")[0])")' 2>&1 | head -40
master
ahead_by=1495 behind_by=0 files=300 commits=250
29fccd66 2025-11-20T11:36:31Z fix: reject interfaces with conflicting options in aggregate mappings (#132)
f44f6aad 2025-11-17T15:59:00Z chore: correcly format CHANGELOG.md dates
bfa493fb 2025-11-17T15:59:29Z chore: prepare release 1.3.0-rc.0
4af59d6c 2025-11-21T14:37:04Z Merge pull request #120 from Annopaolo/prepare-1.3.0-rc.0
43e96a2b 2026-01-15T13:30:19Z fix!: trigger target decoding is now compliant to the typespec (#133)
ed0f1317 2026-01-23T16:12:52Z chore: prepare v1.3.0-rc.1 (#134)
cd6644c8 2026-01-23T10:37:37Z chore: prepare v1.3.0-rc.1
f89a00e3 2026-01-26T11:07:57Z Merge pull request #121 from noaccOS/push-myvsqqtvllkk
34d03020 2025-09-17T12:23:03Z feat: add keepalive option
995633d9 2026-02-13T09:51:17Z Merge pull request #122 from noaccOS/push-woonrwmpumlv
22f447ae 2026-02-13T08:36:25Z chore: prepare v1.2.1-rc.1
eba413b4 2026-02-13T10:01:18Z Merge pull request #123 from noaccOS/push-pnksonqtmmst
3208c1f3 2026-02-24T13:20:00Z docs: Include latest changes for 1.2 (#136)
73850794 2026-02-24T14:39:38Z chore: forward-port release-1.2
42b92a5c 2026-02-24T15:19:07Z chore: merge #135
01ec2281 2026-03-02T16:26:19Z chore: forward-port release-1.3
4b31ced0 2026-03-06T11:22:17Z chore: prepare v1.2.1 (#137)
1e896b2c 2026-03-06T12:50:54Z chore: prepare v1.2.1 (#138)
ce8b62a6 2026-03-06T13:46:20Z chore: prepare v1.2.1
86fcbaf6 2026-03-06T14:58:23Z Merge pull request #125 from noaccOS/push-ttwomwpynqvv
bc05aac2 2026-03-10T08:21:32Z chore: forward port release-1.2
e845830c 2026-03-10T14:44:07Z Merge pull request #139 from noaccOS/push-xkzvyyquuzvo
0b9b2184 2026-03-10T08:26:04Z chore: forward-port release-1.3
61155fd8 2026-03-10T14:47:21Z Merge pull request #140 from noaccOS/push-mmptuxorplps
40c707af 2026-03-10T16:13:45Z feat(mapping): Add required field (#141)
a6efa2e8 2026-03-11T16:31:25Z fix(mappings): allow required key for object aggregated interface (#142)
19bbc241 2026-03-13T15:26:57Z fix: avoid race condition with queue termination
df38d632 2026-03-16T08:53:21Z Merge pull request #126 from noaccOS/push-loosoqqotrrm
ecf84a8b 2026-03-16T09:17:28Z chore: forward port release-1.2
6ad73f2c 2026-03-17T08:48:54Z Merge pull request #127 from noaccOS/push-pzptxrtqmqun
61a4c708 2026-03-17T10:47:18Z chore: forward-port release-1.3
1367e187 2026-03-19T13:16:56Z test: increase coverage (#129)
2195e685 2026-04-07T08:28:47Z chore: prepare v1.4.0-rc.0 (#143)
273fecd3 2026-04-07T08:39:53Z chore: prepare v1.4.0-rc.0 (#130)
e2144499 2026-04-07T09:05:10Z Merge pull request #526 from davidebriani/chore/add-ai-skills
49aaf9c0 2026-04-08T14:24:42Z chore: prepare v1.2.2-rc.0
f9eaa3c4 2026-04-08T14:28:15Z chore: forward-port release-1.2
7e8db21f 2026-04-08T14:28:53Z chore: prepare v1.3.0-rc.2
78b6a9d3 2026-04-08T14:37:16Z chore: prepare v1.2.2-rc.0

$ gh api "repos/astarte-platform/astarte/commits?since=2026-09-30T07:28:27Z&per_page=20" --jq '.[] | "\(.sha[0:8]) \(.commit.author.date) \(.commit.message | split("\n")[0])"' 2>&1 | head -25
b6d46ad4 2026-10-02T07:49:35Z feat(appengine): include unexpected keys in object key errors (#2237)
5fc16d40 2026-09-30T07:34:42Z chore: forward port release-1.4
958bb9fe 2026-10-01T15:03:27Z feat: implement the device_empty_cache_received trigger (#2232)
b0d1ca71 2026-10-01T15:02:14Z build: allow credo and format from umbrella root (#2233)
90efd5ee 2026-10-01T15:01:39Z chore: remove unused labeler config (#2235)
acb5b379 2026-10-01T10:45:16Z fleet: update libs paths (#2231)
b3c13d98 2026-10-01T10:44:59Z ci: update paths in Release please (#2230)
5ae2b936 2026-09-30T07:54:22Z refactor: Update libs paths in tools (#2229)

$ grep -rn "2237\|unexpected key\|unexpected_object_key\|unexpected key" .mule/todo.md | cut -c1-400
1:- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToU
265:- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :u

Two untagged commits since rc.6. One is already queued (#2232 at todo.md:264); checking the other:
$ gh api repos/astarte-platform/astarte/commits/b6d46ad4 --jq '.commit.message, "---FILES---", (.files[] | "\(.filename) +\(.additions)/-\(.deletions)")' 2>&1 | head -40
feat(appengine): include unexpected keys in object key errors (#2237)

Signed-off-by: Riccardo Nalgi <riccardo.nalgi@secomind.com>
---FILES---
apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex +9/-1
apps/astarte_appengine_api/lib/astarte_appengine_api/device/interface_value.ex +14/-1
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/fallback_controller.ex +2/-2
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/interface_values_by_group_controller.ex +12/-37
apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex +2/-2
apps/astarte_appengine_api/test/astarte_appengine_api/device/device_test.exs +1/-1
apps/astarte_appengine_api/test/astarte_appengine_api/device/interface_value_test.exs +2/-1
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_values_by_device_alias_controller_test.exs +4/-1
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_values_by_group_controller_test.exs +4/-1
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_values_controller_test.exs +4/-1

$ git status --porcelain 2>&1 | head
(no output)

Nothing to propose — this is a noop run, which the recipe calls the expected outcome.

**What I checked**
- `docs/COMPATIBILITY.md` still targets `APICompatVersion = 1.2.2` (deviation 10); the 1.3.x wording proposal and the 1.4-rc delta live in `.mule/for-giulio.md`.
- Release list is **byte-identical** to the 2026-10-03 run: newest overall `v1.4.0-rc.6` (2026-09-30, prerelease), newest stable `v1.3.4` (2026-09-18). No `rc.7`, no stable v1.4.0.
- Both rc.6 items were already delivered by that run and are ticked: `todo.md:247` (`fdo-rc6-scope-delta-for-giulio`) and `todo.md:248` (`compat-note-v1.4-rc.6`, for-giulio.md entry written).
- Cheap untagged sweep (`commits?since=rc.6`): only two wire-relevant commits, **both already queued** — `958bb9fe` (#2232, `device_empty_cache_received`) at `todo.md:264`, and `b6d46ad4` (#2237, unexpected keys in object-key errors, 400 + `unexpected_keys`) at `todo.md:265`. Rest are forward-ports, CI and tooling paths.

**Deliberately did not append anything**
- No `compat-note-<tag>` line: re-proposing it would duplicate the already-ticked `compat-note-v1.4-rc.6`, whose wording proposal is already open in `for-giulio.md`.
- No `gh issue create`: this recipe proposes task lines, it does not file issues as a deliverable. No `gh issue edit/close/comment` run, no git touched, working tree clean.

```
Done: no upstream movement since v1.4.0-rc.6 — release list identical to the 2026-10-03 run, rc.6 items already delivered, two untagged commits already queued; appended nothing
Files: none (read-only run, `git status --porcelain` clean)
Verified: gh api repos/astarte-platform/astarte/releases (+prerelease flags), repos/astarte-platform/astarte/commits?since=2026-09-30T07:28:27Z, grep of .mule/todo.md:247-248/264-265 -> pass (no change warranted)
Unsure: nothing — the default branch `master` compare against rc.6 returns a divergent history (1495 commits, oldest first), so I did not use it; the untagged sweep via the commits API is the trustworthy read
Follow-ups: none
```
