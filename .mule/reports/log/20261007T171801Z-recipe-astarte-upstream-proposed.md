slug: recipe-astarte-upstream
verdict: proposed
at:  98a5a94
ran: 2026-10-07T17:18:01Z on DietPi in 172s


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

$ grep -rn "v1\.3\|v1\.2\|1\.4\.0\|upstream version\|targets" docs/COMPATIBILITY.md | head -50; echo ---; ls .mule; echo ---; ls .mule/recipes
52:- The Astarte Dashboard v1.2.2 runs unmodified against Astrate (compose
79:     Recorded against upstream v1.2.0 in
133:   v1.2.0 in `test/conformance/upstream/channels.json`, and two rows were
---
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
recipes
reports
reviews
tasks
todo.md
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

$ grep -rn "astarte-upstream\|upstream" .mule/*.md | head -60
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_1175d718d001HBk4TtINNwee2z

.mule/for-giulio.md:71:- **Milestone recipe run, 2026-10-04: the gate is unchanged as always, but the recipe's standing item finally produced a real finding — both rows of `docs/UPSTREAM-EXPERIMENTAL.md` are tagged "1.4 experimental" and neither of the two features upstream calls experimental. It is a wording/policy call on two files that are yours, so this run filed no issue and queued no task.**
.mule/for-giulio.md:72:  **Gate re-verified, unchanged.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — nothing shipped upstream since the 2026-10-03 run, so #92's parking condition (your 2026-09-04 decision) is still unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** (this milestone's ledger lives under `upstream-parity`); of #47–#89 **all remain closed**; open non-alarm set unchanged — **#93** (`mule-review`, own recipe path), **#92** (parked), **#78** (`milestone-4.0`), **#1** (untouched). The post-rc.6 master delta is now a closed set of **8 commits** (`gh api .../commits -f since=2026-09-30T07:28:27Z`), of which exactly two carry wire surface and **both already have owners in the queue**: `b6d46ad4` (#2237 unexpected object keys) → done as `appengine-unexpected-object-key`, and `958bb9fe` (#2232 `device_empty_cache_received`) → queued as `device-empty-cache-received-trigger` (currently `[!]` blocked). Neither is from stable 1.3 or 1.4, so neither is a v3.0 gate.
.mule/for-giulio.md:73:  **The finding, and it is a measurement, not an inference: upstream's CHANGELOG marks exactly one feature EXPERIMENTAL, and it is neither of ours.** Read this run from `gh api repos/astarte-platform/astarte/contents/CHANGELOG.md` (master, 1052 lines). `rg -i "experimental|unstable"` over that file has **one** hit — `CHANGELOG.md:184`, "[astarte_pairing] FDO authentication (**EXPERIMENTAL feature**, disabled by default)", under `## [1.4.0-rc.1]` — and the string `experimental` matches only `CHANGELOG.md` and `doc/pages/user/062-using_trigger_delivery_policies.md` across the entire upstream repo. Nothing in upstream's *code* annotates either feature either. Our two rows, verbatim (`docs/UPSTREAM-EXPERIMENTAL.md:24-25`), both read `| 1.4 experimental |` in the *Upstream level* column and `upstream 1.4 final: promoted or removed` in the *Reconcile when* column. Neither survives contact with the provenance:
.mule/for-giulio.md:76:  **Why this is worth your attention rather than a note.** The register is wired into the bump rule by its own header: "It is bumped only in the same change that completes the full surface of that upstream level. When that happens, **every open row below tagged with that level must be reconciled first** — the register is the checklist" (`docs/UPSTREAM-EXPERIMENTAL.md:16-18`), and `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`). So these two mis-tagged rows are standing between v3.0's final phase and the bump, and their only listed action — "upstream 1.4 final: promoted or removed" — is a condition upstream can never satisfy in the literal sense, because there is nothing experimental to promote or remove. It *is* dischargeable at 1.4 final under a generous reading ("present and unmarked in the final CHANGELOG" = promoted, so per the header "keep and drop the row"), but that reading is mine, not the file's, and the register does not say it. Fifteen prior runs of this recipe checked *whether 1.4 is final* and *whether the feature still exists*; none checked *whether upstream calls it experimental*, which is why this sat. I have not edited `docs/UPSTREAM-EXPERIMENTAL.md`, `docs/COMPATIBILITY.md` or anything else.
.mule/for-giulio.md:77:  **The knock-on: `docs/COMPATIBILITY.md` deviation 17 repeats the same error, and its wording about one endpoint is now measurably wrong.** It opens "**Always synchronous where upstream 1.4 defaults to asynchronous** — upstream 1.4 runs realm create/delete … and lets the caller opt into synchronous execution with `?async_operation=false`" (docs/COMPATIBILITY.md:335-347). Per the CHANGELOG the async default and the opt-out both date to 1.0.2, not 1.4 — the behaviour is not a moving target, which is the opposite of how that sentence reads. Separately, the sentence lists **device deletion** among the endpoints where the caller can opt into synchronous execution, and upstream has **no such opt-out there**: in `apps/astarte_realm_management/lib/astarte_realm_management_web/controllers/device_controller.ex`, the `operation :delete` doc comment says "Device deletion happens asynchronously, and receiving a 204 response" (:31-38) and the handler is bare — `with :ok <- Devices.delete_device(realm_name, device_id)` (:52-53), with `async_operation` never read. The policy twin does honour it, read as a bare string compare with no parsing (`trigger_policy_controller.ex:155-167`: `if Map.get(params, "async_operation") == "false"`), so upstream accepts only the literal `"false"` and treats absent/`0`/`no` as async. Astrate accepting-and-ignoring the parameter on device delete is therefore an **Astrate extension, not parity** — harmless and stronger, but the record calls it parity. For whoever picks up the queued `docs-sync-rm-delete-device-async-operation-param` line: documenting that parameter on the device-delete operation is still right for an Astrate client, it just has to be recorded as an extension, because upstream's spec has no such parameter on that operation.
.mule/for-giulio.md:78:  **Yours, four small decisions, none of them a code gap:** (a) row #68 — retag it `1.0.2, stable` with reconcile trigger "none, keep as a deviation reference", or **drop the row** and let deviation 17 be the only record, which is the cleaner reading of a register that exists to track *unstable upstream features*; (b) row #67 — split it, leaving `required` untagged (or dropped) and letting the `encrypted` half be tracked by #92/#93 alone, or keep one row that names only the keyAgreement-coupled part; (c) whether deviation 17's opening sentence should say 1.0.2 rather than 1.4, and its device-delete clause should say "extension" rather than "accepted and ignored"; (d) whether the register should gain a line saying what "reconcile" means when a feature turns out never to have been experimental, so the next run is not left guessing as this one was. I did not pick any of them and I have changed nothing: no `gh issue create`, nothing commented, closed or edited on GitHub, and **no `.mule/todo.md` lines** — each of these edits a file that carries your decisions, and queueing them would put your choice in front of the queue as if it were settled.
.mule/for-giulio.md:82:- **The Legion Go answers ssh but its Docker is not running, so every "probe live upstream" line is blocked on you, not on the hardware (measured 2026-10-03).** `ssh -o ConnectTimeout=5 legion 'echo ok'` returns `ok`, but `docker ps` fails with `dial unix /home/atsetilam/.docker/desktop/docker.sock: no such file or directory`, `systemctl is-active docker` answers `inactive`, and `sudo -n true` answers `è necessaria una password` — so the mule cannot start it. It is Docker Desktop's socket, not a rootful `dockerd` the mule could launch as a user, and no rootless extras are installed. Consequence for the queue: `appengine-payload-reason-status-map` (the remaining eight `payload.Reason` values that still answer 500) cannot be measured the way its line demands. Today's `appengine-missing-required-422` was therefore mapped from the frozen upstream source at `b6d46ad4` (`Device.validate_required_mappings/2` → `fallback_controller.ex:215-219` → `ErrorView "422_missing_required_mapping.json"`), and the commit says so plainly. Nothing about that sibling line needs rewriting — only Docker up on the Legion Go.
.mule/for-giulio.md:93:  **Verified unchanged, no new gaps.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition, frozen by you on 2026-09-04 ("wait for a stable v1.4.0", an `-rc.N` does not satisfy it), is still unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Of #47–#89 **all are closed**; the open non-alarm set is unchanged — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** (untouched per standing instruction). rc.6 adds no v3.0 gap: it is already audited in the two entries below, and I re-measured its non-FDO half myself (`compare/v1.4.0-rc.5...v1.4.0-rc.6`, 27 commits) — outside FDO the only library changes are a `# coveralls-ignore-stop` line in `apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex:43` and `realms/queries.ex` *removing* the ownership_voucher / to2_sessions / unconfirmed_devices table creation (moved into migrations). No Realm Management, AppEngine or Pairing wire surface moved.
.mule/for-giulio.md:94:  **The recipe's standing item is satisfied and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental*, reconcile trigger "upstream 1.4 final". 1.4 is not final, so there is nothing to promote and nothing to deprecate, and the register is accurate as written. `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked this run) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision ("bumped only in the same change that completes that surface, after reconciling every row tagged with the level"). The milestone-bump issue the recipe asks me to file **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": v3.0's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is delivered, the 1.4 half is open (#92, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut.
.mule/for-giulio.md:95:  **New, and the reason this run is not a copy of the twenty before it: `.mule/research/` is listed in `.mule/.gitignore`, so the source document this milestone is built on — `.mule/research/upstream-parity-2026-08-22.md`, named at `.mule/milestones.md:114` — cannot exist in any clone, and the Pi's copy is gone.** Evidence: the last line of `.mule/.gitignore` is `research/`, a bare directory pattern that matches `.mule/research/` at any depth; the directory does not exist here; `git ls-files .mule` has no `research` entry; and `grep -rl 'upstream-parity-2026-08-22' --include='*.md'` over the tree matches **exactly one file, `.mule/milestones.md:114` itself** — there is no tracked equivalent (the only parity document under `docs/handoff/` is the v2.0 `flow-parity-audit-2026-07-29.md`). This is the sharper form of the 2026-09-06 note ("the section's source document does not exist in the repo"), and it is also why that note never got resolved: dropping the file back in place as-is re-dirties the tree on the next tick, because `.mule/.gitignore`'s own header states that any of its unlisted siblings "makes the tree dirty, which makes `tick` refuse to run and so turns the timer into a permanent silent no-op". **Yours, three ways:** (a) commit it once behind a negation (`!.mule/research/upstream-parity-2026-08-22.md`), so the milestone keeps a written reference; (b) move the source document to a tracked path (`docs/handoff/upstream-parity-2026-08-22.md`) and fix the `:114` pointer to match; (c) declare it lost and let the #47–#89 issue bodies plus the per-release CHANGELOG diffs stand as the reference. I did not pick one — every option edits a file that carries your decisions, and (a) versus (b) differ on a question only you can answer: is that document mutable working state or a frozen record?
.mule/for-giulio.md:96:  **Read together with the 2026-09-26 `waiting-on.md` entry, not as a repeat of it: the same rule explains why that file still does not exist either.** `.mule/waiting-on.md` is *not* in `.mule/.gitignore`, so under the header's rule it cannot be created here without wedging `tick`, and it is not tracked (`git ls-files .mule` has no `waiting-on`), so it is not here at all — and nothing in the repo references it (`rg waiting-on` over `.mule/` and `tools/` matches only prose in this file and copies under `.mule/reports/log/`). So #92's escalation is currently parked on a condition the repository cannot express. That entry's one-liner — `systemctl list-timers 'mule-upstream*' --all` — is still what settles whether the mechanism lives in a Pi-side unit outside this repo; **I could not run it from this box** (no systemd here), so the 2026-09-26 finding stands unchanged rather than refuted, and this is only the repo-side half of the answer.
.mule/for-giulio.md:103:- **Delta on the two open COMPATIBILITY.md wording entries further down this file (the v1.3.4 one folding in v1.3.3, and the original v1.3.2 one): upstream `v1.4.0-rc.6` is wire-inert for Astrate, so nothing in that proposal is re-derived — only its version reference and one line of its "not yet emulated" list change.** Not applied: `docs/COMPATIBILITY.md` is untouched (still the 1.2.2 target, deviation 18's "which upstream 404s" at `docs/COMPATIBILITY.md:360`, `APICompatVersion = "1.2.2"` at `internal/realm/service.go:588`), and the FDO half of the rc.6 delta already has its own entry below rather than being folded in here.
.mule/for-giulio.md:105:  - **The Dashboard's new `unknown_status` is computed client-side, not a new field** — upstream `models/Device/index.ts:213-215` derives it as `last_connection && !last_disconnection && !isConnected`, and both timestamps are already on Astrate's device shape (`internal/appengine/service.go:117-118`), so nothing is added to the `/appengine` device. The one consequence is a display change: a device whose disconnection was never recorded reads "Unknown status" where it used to read "Disconnected". Astrate ships no dashboard, so this is not a compat item — it only matters if #78's eventual work mirrors that dashboard.
.mule/for-giulio.md:107:  - **The two edits the proposal needs, both yours to approve:** (1) the 2026-09-18 entry's parenthetical "(v1.4.0 is still rc.5-only)" becomes "(v1.4.0 is still rc.6-only)" — re-verified this session that rc.6 is the newest 1.4 prerelease and there is still **no stable v1.4.0** (`gh api repos/astarte-platform/astarte/releases`); (2) its "**FDO authentication** (pairing, disabled by default)" item, and the same "experimental FDO pairing auth" phrase in the proposed doc sentence, understates the 1.4 line, which ships a full ownership-voucher surface (upload with mandatory `hw_id`, deletion, re-TO0, per-voucher `expiry`, a spec-mandated realm-less TO2 path) that Astrate does not emulate — #78, deliberately parked at `.mule/milestones.md:124`. "FDO ownership vouchers (1.4 line, not stable upstream)" says it without claiming the rc.5 shape.
.mule/for-giulio.md:108:  Unverified: the upstream side of the item list is the upstream-watch recipe's `v1.4.0-rc.5...v1.4.0-rc.6` comparison — the same measurement as the FDO entry below, taken there and not re-derived here (a raw fetch of `models/Device/index.ts` at the rc.6 tag did not resolve from this box). I re-verified only the Astrate-side claims above and the release list. No `.mule/todo.md` lines: the deliverable is the delta, and the proposal itself stays open.
.mule/for-giulio.md:112:- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published 2026-09-30T07:28Z, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and the one path #78's body states as fact no longer exists upstream. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix") deletes the `realm_name` path parameter from all seven TO2 `operation` specs in `fdo_onboarding_controller.ex`. **This breaking change is not in the rc.6 CHANGELOG entry** (which lists only the to0 route, voucher deletion, the `ASTARTE_BASE_URL_*` rename, device-deletion cascade, `device_id` on upload, immediate registration, the two x5chain/voucher-store fixes and the cross-realm delete fix) — so anyone tracking upstream by its changelog, including #78's body, misses it. The realm now comes from the voucher row: the global-keyspace query `fetch_device_id_and_ownership_voucher_and_realm/1` returns `{realm, device_id, voucher_data}` and the `:fdo` plug fetches the session by GUID alone (`plug/fdo_session.ex`: `Session.fetch(guid)`); `RvTO2Addr` no longer carries the realm in the hostname either ("the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure"). **Yours:** #78's acceptance is "a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end" — against rc.6 that SDK posts to `/fdo/101/msg/{60..70}` with no realm, so mirroring rc.5's shape would mean mirroring a path upstream has deleted. No action while parked; the body is just stale.
.mule/for-giulio.md:121:- **Two behaviours that look like one feature but are not**, and #78 should not treat them as one: owner-side voucher deletion revokes the rendezvous registration **first** and keeps the voucher if that fails, whereas RM device removal deletes the voucher row **without** revoking — `RealmManagement.DeviceRemoval.Core.delete_ownership_voucher!/2` logs `"Deleting ownership voucher without revoking its rendezvous registration"` (tag `fdo_voucher_deleted_without_revocation`) and its own doc says "N.B.: only the database row is removed." So removing a device can leave a live rendezvous registration pointing at a voucher Astrate/upstream no longer holds. Not a blocker, but it is the kind of asymmetry an operator guide must state.
.mule/for-giulio.md:123:- **One premise in #78's own comment is now false.** It says "upstream ships almost no public FDO documentation (code-level only); ours must exceed that." rc.6 adds two documented sections: `doc/pages/user/035-register_device.md` gained `## FIDO Device Onboard` (Owner Keys, Ownership Vouchers) and `## Credentials Secret Lifecycle`, and `doc/pages/architecture/050-pairing_mechanism.md` gained `## FIDO Device Onboarding` (TO0, TO2, Owner Keys management) — commits `162062e`, `d1c39fb`. "Docs must exceed upstream" is still the right requirement, but the gap is far smaller than #78 assumes, and those two pages are now the cheapest available scope input.
.mule/for-giulio.md:135:  **`triggers.forward.static_headers` gets no validation at all, while trigger *action* headers get a blocklist and a size cap.** `config.validate` checks kind, url, method and subject (internal/config/config.go:339-375) and nothing else; the action path runs `validateStaticHeaders` (internal/engine/triggers/validation.go:140-155) against upstream's hop-by-hop/sensitive list (validation.go:21-34) plus an 8 KiB total. Measured that the practical harm is smaller than it looks — net/http ignores a static `Host` and derives `Content-Length` from the body, so a blocked name is silently dropped rather than corrupting the request. **Yours:** should bus-forwarder headers carry the same policy as action headers, and if so is a boot-time rejection of a currently-"working" config acceptable?
.mule/for-giulio.md:145:  **Correction on a fact 19 entries have repeated as settled: there is no `.mule/waiting-on.md`.** #92's own comment (2026-09-04) states the issue "has a row in `.mule/waiting-on.md`" and that the weekly `mule-upstream-watch` job reads that file as its first step. The file does not exist here, and `rg 'waiting-on|upstream-watch'` over the whole tree — including gitignored `.mule/` and `tools/` — matches only prose in this file and copies of it under `.mule/reports/log/`. So the escalation that was supposed to catch a stable upstream `v1.4.0` is not in the repo: either it lives in a Pi-side unit outside it, or it was never wired. One command settles it: `systemctl list-timers 'mule-upstream*' --all; grep -rl 'upstream-watch\|waiting-on' /etc/systemd/system /root 2>/dev/null`. If it was never wired, #92 is parked on a condition nothing is watching, and the 2026-09-04 claim that "the mechanism was added because it was missing here" is false. The gate is being caught by hand in the meantime — today's milestone run re-derived it from `gh api .../releases`.
.mule/for-giulio.md:149:- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Status: not started", which is no longer true.** Release sweep (`gh api repos/astarte-platform/astarte/releases`, run this session) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.5` (2026-08-20, prerelease) — byte-identical to what the 2026-09-24 run saw, so nothing shipped upstream in the two days since. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Open non-alarm issues remain exactly **#93** (aclhook comment, `mule-review`, its own recipe path), **#92** (keyAgreement, parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** (FDO, `milestone-4.0`), **#1** (untouched per standing instruction). No new gaps, **no issues filed, no `.mule/todo.md` task lines added** — there is nothing machine-checkable left, and filing anything here would duplicate #92/#93/#78.
.mule/for-giulio.md:150:  **New, and the reason this run is not a copy of the last thirteen: `.mule/milestones.md:127-128` still reads "Status: not started. First recipe job: triage #47–#89 into an ordered plan (which are audits vs features vs decisions), file sub-issues where work splits, escalate the 'decide' set in one batch."** All of that happened: of #47–#89 every issue is closed except the deliberately-parked #78. The section describes as future work a triage that is finished, and it is the exact text every future run of this recipe reads to pick its target — so the recipe has been re-deriving "the backlog is done, the gate is upstream" from scratch on a timer for three weeks. Suggest rewriting those two lines to state the real position: *1.3 surface delivered (#47–#89 closed); milestone blocked on upstream v1.4.0 going stable; final phase = reconcile the two `docs/UPSTREAM-EXPERIMENTAL.md` rows, answer #92, bump `APICompatVersion`.* That also makes the blocking condition greppable instead of inferred.
.mule/for-giulio.md:151:  **The recipe's standing item is satisfied, and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still has exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental* with reconcile trigger "upstream 1.4 final". v1.4 is not final, so both rows stay as written: nothing to promote, nothing to deprecate, and the register is accurate. `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision; the milestone-bump issue the recipe asks for **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": the milestone's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is done, and the 1.4 half is unresolved (#92 open, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut. Also still true and already escalated on 2026-09-06, so not repeated: the section's source document `.mule/research/upstream-parity-2026-08-22.md` does not exist in the repo.
.mule/for-giulio.md:174:  also accepts the bare upstream name `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`
.mule/for-giulio.md:220:- **github-issues triage run, 2026-09-25: still nothing proposable — the mule-alarm pile is 18 issues (#94–#111), and two corrections to what earlier runs inferred about it.** Twenty-two open issues, same set as every run since 2026-09-05: **#111–#94** are the daily `mule-alarm` noise (zero comments, telemetry not code), **#93** aclhook comment (mule-review, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable upstream v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-24 run: **#111** (today, 09-25 11:01:07Z, "nothing has landed in 16h"). Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines added**. **Proposal, extending the 2026-09-24 one: close #94–#110**; leave #111 (live today) to self-expire. Two corrections, both from reading `tools/mule.sh` rather than the issue titles: **(1) the "missed days" in earlier entries (09-15, 09-20, 09-23) are evidence of activity, not its absence** — I previously called them timer gaps; the alarm is gated on a `.alarmed` sentinel that `beat` deletes on every land (`mule.sh:737,748-749`), so a day with no new issue is a day something landed, not a day the timer slept. **(2) the alarm body's closing line is wrong**: "Close it — the mule reopens a new one if the silence continues" (`mule.sh:764`) implies issue state drives the alarm; it does not. Nothing reads or reconciles the 18 open issues, so the pile can only grow and closing the old ones is safe but inert — and `cmd_refill` commits `mule: refill the queue` (mule.sh:719-720) **without calling `beat`**, so a doc-only refill never counts as landing work. The titles' ages also imply a beat roughly every day around 19:00–21:00Z (#111 filed 09-25 11:01Z at 16h → last beat ≈ 09-24 19:01Z; #110 filed 09-24 10:54Z at 15h → ≈ 09-23 19:54Z) while the ~11:00Z daily check has never once found under the 8h threshold — **what calls `beat` at that time is unverified**, and one command on the Pi settles it (`grep -E 'landed|checked:' /root/astrate-mule/.mule/log | tail -20`). The queue has still landed nothing since ~2026-09-04/05; if that idle streak is not intentional, the wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:224:- **Milestone recipe run, 2026-09-24: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-23 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:228:- **github-issues triage run, 2026-09-24: still nothing proposable — the mule-alarm pile is now 17 days (#94–#110).** Twenty-one open issues, same set: **#93** aclhook comment (mule-review, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable upstream v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-20 runs: the alarms **#108** (09-21), **#109** (09-22) and **#110** (today, 09-24 10:54Z, the live one) — a four-day gap on 09-23 not withstanding. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines added**. **Proposal, extending the 2026-09-20 one: close #94–#109** — each was a one-day event, superseded by the next day's alarm, never actionable; leave #110 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:232:- **Milestone recipe run, 2026-09-23: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-21 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:240:- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`ASTRATE_FLOW_CONFIG`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
.mule/for-giulio.md:248:- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:252:- **github-issues triage re-run, 2026-09-20 (evening): nothing proposable, no change since this morning's 13:10 run.** Re-surveyed with the recipe command: the issue set is identical to the 2026-09-20 morning triage — **#94–#107** are still the daily `mule-alarm` noise (14 issues, zero comments, telemetry not code), **#93** aclhook comment rewrite is still mule-review with `8c61268` already pushed (its own recipe path), **#92** keyAgreement is still parked on a stable upstream v1.4.0 with its waiting-on.md row and the mule-upstream-watch escalation wired (an -rc does not satisfy it), **#78** FDO is still the milestone-4.0 design/investigation already escalated, and **#1** is untouched per standing instruction. No new issue has been filed since the morning run. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines proposed**. This morning's proposal stands unchanged: **close #94–#106** — each was a one-day event, superseded by the next day's alarm, never actionable — and leave **#107** (still the newest, no #108 yet) to self-expire. The queue idle streak since ~2026-09-04/05 persists; if that is not intentional the wedge is on the Pi (queue review = the fix, per the 13:10 line).
.mule/for-giulio.md:256:- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:268:- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:272:- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:276:- **deviceid wire-parity concern, from the 2026-09-18 code review: Astrate is likely stricter than upstream on the wire form's unused bits.** `pkg/deviceid/deviceid.go:37` uses `base64.RawURLEncoding.Strict()`, so `Parse` rejects `"AAAAAAAAAAAAAAAAAAAAAB"` (unused low 4 bits of the 22nd char set; rejection pinned at deviceid_test.go:176), while the comment at deviceid.go:34-36 claims parity with `Elixir Base.url_decode64!(padding: false)` — which decodes the 128 payload bits and discards the unused ones rather than erroring. If so, a device registered with a non-canonical-trailing-bits id on upstream Astarte is accepted there but rejected by Astrate's `Parse` — a wire-compat break in a package whose purpose is the wire form. Not verified on the Pi (no Docker/Elixir); a `[legion]` probe (`deviceid-trailing-bits-upstream-probe`) is queued to measure it. Once confirmed, the call is yours: relax `Parse` to upstream's decode-and-ignore, or keep strict canonicality as the package's contract.
.mule/for-giulio.md:280:- **COMPATIBILITY.md wording update for upstream v1.3.4 (newest stable, 2026-09-18; v1.4.0 is still rc.5-only) — folds v1.3.3 and v1.3.4 into the open v1.3.2 proposal below; both are maintenance-only and wire-inert, so that proposal stands unchanged.** v1.3.3 (2026-08-07) was an empty-body patch and v1.3.4 (2026-09-18, `astarte_data_updater_plant` only — more frequent fullsweep GC on AMQPDataConsumer processes plus an RPC-availability fix; release body via `gh api repos/astarte-platform/astarte/releases` `v1.3.4`) introduce **no wire/API surface change**, so the v1.3.2 wording proposal stands complete: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"). The only delta a v1.3.4-aware doc carries is the version reference: the proposed §Infrastructure-differences sentence's "until the milestone that adopts v1.3.2 as the target" reads "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target". Raw: [v1.3.4](https://github.com/astarte-platform/astarte/releases/tag/v1.3.4).
.mule/for-giulio.md:288:- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) shows a new **stable v1.3.4** (2026-09-18) — maintenance only (data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change), so no new gap for Astrate — and still **no stable v1.4.0**: newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:300:- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-13 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:304:- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-12 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:308:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:310:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
.mule/for-giulio.md:316:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:318:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

$ grep -rn "1\.3\.4\|1\.3\.5\|v1\.3" .mule/todo.md | head -40; echo "=== tail of todo ==="; tail -40 .mule/todo.md
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_1175d9baa001Q4oClPoksfiLsX

282:- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
283:- [x] compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (newest stable, 2026-10-05, maintenance-only: mix.exs/mix.lock version bumps plus a single CHANGELOG "astarte_data_updater_plant RPC server reliability" entry, no wire/API surface change; v1.4.0 is still rc.6-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2/v1.3.4 wording proposal there already stands and its version-reference delta was recorded on 2026-10-05 as reading "newest stable v1.3.5", so fold this in as that one-word delta rather than re-deriving the proposal.
=== tail of todo ===
- [x] compat-note-v1.4-rc.6 [auto]: record in `.mule/for-giulio.md` that v1.4.0-rc.6 is wire-inert for Astrate, as a delta on the open v1.3.2/v1.3.4 wording proposal rather than a re-derivation. Measured, the whole non-FDO delta is inert: DUP RPC-server reliability is OTP supervisor/GC tuning; the Dashboard's new `unknown_status` is computed client-side from `last_connection && !last_disconnection && !isConnected` (models/Device/index.ts:213-215) and both timestamps are already on Astrate's device shape (internal/appengine/service.go:118), so a device whose disconnection was never recorded reads as "Unknown status" rather than "Disconnected"; the `on_device_registered` dispatch moved from `Agent.register_device` to `Engine.register_device` (pairing/engine.ex:157-165, 186-193) so it also covers the FDO path, but the Agent path is byte-for-byte the same rule and Astrate already fires it on every successful registration (internal/pairing/service.go:221-223); the `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"); and `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma. So the proposal's "v1.4.0 is still rc.5-only" reference becomes rc.6, and its "FDO authentication (pairing, disabled by default)" line understates the 1.4 line, which ships a full ownership-voucher surface Astrate does not emulate (#78, parked). Do not edit docs/COMPATIBILITY.md and do not restate the standing proposal.

- [x] container-stop-deadline [auto]: make `cliInstance.Stop` (internal/flow/blocks/container/docker.go:167-175) derive its 15s bound whenever the passed context carries **no deadline**, not only when it is nil — both cleanup callers pass `context.Background()` (docker.go:122, the "Best-effort cleanup so we do not leave orphans on mapping failure" path, and block.go:121, the not-ready cleanup in `New`), so a wedged docker daemon leaves an unbounded `docker rm -f` inside `exec.CommandContext` holding flow instantiation, while `Block.Stop` (block.go:310-311) does the same job correctly with an explicit 15s timeout. Measured with the package's own `Run` hook: the `rm` invocation arrives with `ctx.Deadline()` unset on both paths. Add a case asserting `inst.Stop(context.Background())` reaches the injected `Run` with a deadline set (it does not today) and that the hostPort-failure path still issues `rm -f <id>` (docker.go:120-124, never executed by any test today). Container-free.
- [!] container-timeout-bounds [auto]: bound `timeout_ms` and `ready_timeout_ms` in `parseConfig` (internal/flow/blocks/container/block.go:205-219) — both are only checked `> 0` and then multiplied unchecked (`time.Duration(n) * time.Millisecond`), which wraps: measured `ready_timeout_ms: 18446744073709` gives `-551.616µs`, so `New` fails with "container: not ready after wait" without ever issuing a probe (the wrapped negative deadline fails the loop guard at httpbridge.go:64), while `ready_timeout_ms: 9223372036854` gives 2562047h47m16s and `New` never returns (still polling `/healthz` after 3s in the probe). The two consumers of a non-positive duration also disagree — `Bridge.timeout()` (httpbridge.go:39-44) treats it as "use the 5s default" while net/http treats a non-positive `Client.Timeout` as no deadline at all (measured: a client with `Timeout: -551µs` returned err=nil after a 3s handler). Reject a value that would overflow or exceed a documented maximum, naming the key and the value in the error, and add table cases asserting `New` refuses both (it accepts them today). Nothing bounds these on the wire either: internal/flow/blocks/schema.go:19 advertises both as bare `"type":"integer"`. Container-free. — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [x] container-parseconfig-rules-test [auto]: pin the four `parseConfig` rules that no test touches (internal/flow/blocks/container/block.go:188-219 — all `count=0` in `go test -coverprofile`, which is 71.2% for the package): `config` must be a JSON object (192-193), `port` an integer in 1–65535 including the 0 / 65536 / non-integer rejections (197-202), and `timeout_ms` and `ready_timeout_ms` positive (205-210, 215-217) — beside the existing `TestConstructor_RequiresImage` (block_test.go:99), which is the only config case the suite has. The port rule is the one with teeth: it is what keeps `127.0.0.1::99999` out of the `docker run` argv (docker.go:93). Container-free, pure.
- [!] container-response-cap-test [auto]: pin `Bridge.RoundTrip`'s 1 MiB response cap (internal/flow/blocks/container/httpbridge.go:133-140 — `io.LimitReader(resp.Body, maxBody+1)` plus the explicit rejection; lines 136, 138.16 and 138.40 are all uncovered), the only bound between a tenant-supplied container and process memory on the message path: assert a handler answering exactly `maxBody+1` bytes is rejected with the "exceeds" error and that exactly `maxBody` bytes still round-trip, driving a `Bridge{MaxBodyBytes: N}` so the fixture stays small. Also worth one line in the same change: the cap is absent from the operator-facing catalog entry (internal/flow/blocks/info.go:116) and from the schema, so an author whose container emits 2 MiB learns about it only from a runtime error. Container-free. — BLOCKED: lint failed: internal/flow/blocks/container/block_test.go:500:2: redefines-builtin-id: redefinition of the built-in function cap (rev
- [x] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, housekeeping and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
- [!] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy alias fields to `InterfaceDefinition` and `InterfaceMapping` in docs/api/astarte_realm_management_api.yaml — both declare `additionalProperties: false` (yaml:935, 972) yet omit interface-level `quality` (alias of `ownership`) and `aggregate` (alias of `aggregation`) and mapping-level `path` (alias of `endpoint`), all three of which the parser accepts and canonicalises: decoded into dedicated fields (pkg/interfaceschema/parse.go:116-122, 136-138), normalised before any rule runs (parse.go:212-236), and re-encoded canonically by `ParseInterfaceCanonical` so the store never sees them (parse.go:165-170, 307-311). A spec-conformant client — or any generated client with this schema baked in — therefore rejects a valid upstream-shaped install/update with a spurious 422 while the service would have accepted and stored it as `ownership`/`aggregation`/`endpoint`. Add the three as documented deprecated aliases (enums mirroring ownership/aggregation, `path` a plain string) or relax `additionalProperties`, and record the three rules the code enforces: `ownership` + `quality` is a violation (parse.go:218) and `aggregation` + `aggregate` is a violation (parse.go:227), while `endpoint` + `path` silently keeps `endpoint` (parse.go:307-311). None of the three may appear in a `required` list. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstrained string described only as "The device hardware ID." — `Service.DeleteDevice` runs the path value through `deviceid.Parse` (internal/realm/service.go:118-121), which accepts exactly 22 characters of unpadded base64url in strict mode and rejects padding, the standard alphabet's `+` and `/`, and non-canonical trailing bits (pkg/deviceid/deviceid.go:24-26, 33-37, 42-60); every rejection is folded into `store.ErrNotFound` (service.go:119-121) and answered `404 {"errors":{"detail":"Device not found"}}` (internal/realm/http.go:107-109), so a client that sends a well-formed standard-base64 or 23-character padded id is told the device does not exist rather than that its id is malformed. State the encoding and the 22-character length and add `pattern`/`minLength`/`maxLength`, exactly as docs-sync-pairing-deviceid-base64url did for the pairing twin. Do NOT encode the strict trailing-bits rule in the pattern: whether it is a deviation worth recording is still open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_realm_management_api.yaml — `data` is declared as a bare `type: object` with no properties (yaml:250-261) although the handler parses it with the same `interfaceschema.ParseInterfaceCanonical` as `installInterface` (internal/realm/service.go:196 vs 161) and then checks the body's `interface_name`/`version_major` against the URL's `{name}`/`{major}` (service.go:203-208), so the update body — the one carrying every interface rule — is the only interface body in the spec with nothing to validate against. `$ref` `#/components/schemas/InterfaceDefinition` as `installInterface` does (yaml:115-117), and say in the operation description that the body must name the same interface and major as the URL (disagreement is 409, already documented) and that `version_minor` must increase (409). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.
- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` shows `detail: Bad Request` (yaml:1732) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1792), while the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58 via `WriteBadRequest`/`WriteInternalServerError`, golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json); every 400/500 on that surface goes through the same astarteapi writers. The other detail examples on the file are already correct — change only those two — and add the docs-half assertion `TestAppEngineErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go, which pins the pairing and realm-management halves. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDetail` schema-level example (yaml:361-362) and the `BadRequest` response (yaml:385) show `detail: Bad Request`, and `InternalServerError` (yaml:477) shows `detail: Internal Server Error`, but the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58, via WriteBadRequest/WriteInternalServerError, pinned by pkg/astarteapi/testdata/error_bad_request.json and error_internal.json). While you are in `ErrorDetail`: its `example` sits at the schema level and shows a bare `{detail: ...}` where the envelope the wire sends is `{"errors":{"detail":...}}` — move it under `errors` (or drop it, as the realm-management and pairing `ErrorDetail` schemas carry no example) so a generated client is not shown the wrong shape. Add `TestHousekeepingErrorDetailExamples` beside the pairing and realm-management halves in internal/swagger/swagger_test.go. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `"Internal server error"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implements `device_empty_cache_received` end to end — `SimpleEvent` field 23, an empty `DeviceEmptyCacheReceivedEvent` whose encoder emits exactly `{"type":"device_empty_cache_received"}` (astarte_core encoder.ex:54-62), fired by DUP at `control_handler.ex:161` *inside* the `with` that runs `send_control_consumer_properties` → `resend_all_properties` → `set_pending_empty_cache`, so it fires only when all three succeed and never on failure, and the commit deletes the `# TODO: implement on_empty_cache_received` from `astarte_events/lib/astarte_events/triggers/core.ex`. Astrate already accepts the trigger but never evaluates it: `OnDeviceEmptyCacheReceived` sits in `deviceOns` with value `false` (internal/engine/triggers/match.go:123, the `evaluated` flag read at match.go:406) and the comment's stated reason ("upstream defines no SimpleEvent variant for it") is now false. Add the event body constructor to internal/engine/triggers/events.go (the payload is `{"type":"device_empty_cache_received"}` and nothing else — no device id, no metadata), flip the `deviceOns` value to true, and fire it from `handleEmptyCache` (internal/engine/control.go:58) with the same idiom `handleLifecycle` uses for `device_connected` (`e.fireDevice` + `stream.Publish`, internal/engine/engine.go:346-353), at the point where `m.Ack()` runs — i.e. after `resendServerProperties` and `sendConsumerProperties` have both succeeded inside `retryStore`. Add a T1 test that a trigger on it fires exactly once per accepted emptyCache and does not fire when the resend fails. — BLOCKED: wrote nothing
- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362) — the second half docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail: Bad Request`, the reconstructed title-case form that `TestHousekeepingErrorDetailExamples` (internal/swagger/swagger_test.go:687) now pins against everywhere else on this spec, while the wire emits the frozen `astarteapi.DetailBadRequest` = `"Bad request"` (pkg/astarteapi/envelope.go:37, via `WriteBadRequest` on every 400 — internal/housekeeping/http.go:65, 166, 171, 183). Because the example sits on the generic `ErrorDetail` schema, which all eight error responses `$ref` (yaml:377-477), it is the example a generator falls back to for any response whose own example it cannot resolve — so it teaches `"Bad Request"` for a 404 and a 500 too. Change it to `Bad request`, or drop it as the realm-management and pairing `ErrorDetail` schemas do (neither carries one). Then extend `TestHousekeepingErrorDetailExamples`, whose comment currently says the schema-level example is "deliberately not pinned" (swagger_test.go:684-686), to pin it like the responses. One thing not to "fix": the example is nested one level in — it is the `errors` sub-schema's example, not a bare envelope — so do not re-wrap it in another `errors:`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-patch-422-field-error-examples [auto]: add the three missing PATCH 422 bodies to the `ValidationErrors` response component in docs/api/astarte_housekeeping_api.yaml (yaml:427-436, referenced once at yaml:190) — it carries a single example, `error_name: [invalid_update_parameters]` (internal/housekeeping/http.go:176-177), but `patchRealm` emits three more field-error envelopes: `jwt_public_key_pem: ["can't be blank"]` (http.go:188-189), `device_registration_limit: ["is invalid"]` (http.go:193-194) and `datastream_maximum_storage_retention: ["is invalid"]` (http.go:198-199). Add them as further `examples` on the media type (OpenAPI 3.0.3, yaml:1) or split the component per case. Also record that the *messages* differ from POST for the same conditions: create rejects a negative limit/retention with the flat `ErrValidation` detail `device_registration_limit must be non-negative` / `datastream_maximum_storage_retention must be non-negative` (internal/housekeeping/service.go:143-148, 422 `ErrorDetail` shape), while PATCH answers `is invalid` in the FieldErrors shape — and that the service's own three `ErrValidation` branches (service.go:205-213) are unreachable from REST, because `patchRealm` pre-checks the identical conditions and answers first. Say in the `RealmPatch` field descriptions that a negative value is a 422 rather than silently ignored. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [!] docs-sync-hk-wrong-type-field-400 [auto]: widen the `BadRequest` component description in docs/api/astarte_housekeeping_api.yaml (yaml:378, "malformed or missing data envelope") — measured today, a perfectly well-formed `data` envelope carrying a wrong-typed field answers **400**, not the documented 422. `astarteapi.DecodeData` into `map[string]json.RawMessage` (internal/housekeeping/http.go:170) cannot fail on a type mismatch, so the failure surfaces on the second decode into the typed body (http.go:182-184) and is answered `WriteBadRequest`. Probe output: `{"data":{"device_registration_limit":"100"}}` → `json: cannot unmarshal string into Go struct field patchBody.device_registration_limit of type int32`; `{"data":{"jwt_public_key_pem":123}}` → the string twin; `{"data":{"device_registration_limit":1.5}}` → `cannot unmarshal number 1.5`; `{"data":[]}` → `cannot unmarshal array into Go value of type patchBody`. Same on POST, which decodes straight into `realmBody` (http.go:63-67). This is unreachable from a spec-conformant client — `RealmCreate`/`RealmPatch` already type both fields — so document it as a one-line extension of the 400 description on both operations rather than inventing a 422, and do NOT re-shape the `BadRequest` example, whose `detail: Bad request` is the frozen constant `TestHousekeepingErrorDetailExamples` pins. One distinction worth stating in the same description: a wrong-typed field is *not* the `invalid_update_parameters` case — unknown keys are caught earlier by the `patchAllowedFields` sweep (http.go:174-180, 422 FieldErrors) on PATCH only, and silently dropped on POST. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `components.parameters.RealmName` and `RealmCreate.realm_name` both carry `pattern: '^[a-z][a-z0-9]*$'` and `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go) holds both against the realm name CHECK extracted from migrations/000002_metadata.up.sql — but the read side is still bare `type: string`: `Realm.realm_name` (components.schemas.Realm) and `DataEnvelopeRealmNames.items`. Decide whether a response field earns the same pattern (a generated client learns nothing new from it, but a reader diffing create against get sees one constrained and one not) and, if it does, add it to both and extend that test to hold the response-side pattern against the same CHECK. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] store-register-inhibit-preserve: RegisterDevice must not overwrite status to 'registered' when a device was inhibited before first credentials request (preserve 'inhibited' state). File: internal/store/devices.go — BLOCKED: gates failed
- [x] store-aliasvalues-self-exclusion-test: Add test for AliasValuesTaken to assert self-exclusion (same device's own alias values do not count as taken). File: internal/store/devices_test.go
- [!] store-latestindividual-empty-errnotfound: Add store-level test asserting LatestIndividual returns ErrNotFound for empty series and returns newest on multiple samples. File: internal/store/datastreams_test.go — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [x] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.go
- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
- [x] compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (newest stable, 2026-10-05, maintenance-only: mix.exs/mix.lock version bumps plus a single CHANGELOG "astarte_data_updater_plant RPC server reliability" entry, no wire/API surface change; v1.4.0 is still rc.6-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2/v1.3.4 wording proposal there already stands and its version-reference delta was recorded on 2026-10-05 as reading "newest stable v1.3.5", so fold this in as that one-word delta rather than re-deriving the proposal.
- [!] docs-native-socket-handshake-errors [auto]: add the two WebSocket-handshake failures to both socket operations in docs/api/astrate_native_api.yaml — `/astrate/v1/{realm}/socket` responses (yaml:335-345) and `/appengine/v1/socket/websocket` (yaml:398-404) document only 101/401/403/500 and 101/401/500 — but measured today with the repo's own dependency (probe in /tmp, coder/websocket accept.go `verifyClientRequest`): a plain GET with no `Connection: Upgrade` answers **426 Upgrade Required**, `Content-Type: text/plain; charset=utf-8`, body `WebSocket protocol violation: Connection header "" does not contain Upgrade\n`, and `Sec-WebSocket-Version: 14` answers **400** `unsupported WebSocket protocol version (only 13 is supported): "14"\n`. Add both to each operation and say two things explicitly: the bodies are plain text from `http.Error`, NOT the `{"errors":{"detail":...}}` envelope every other response on this spec uses (a client that JSON-parses them breaks), and they are only reached after authentication — native: the route's guard `mw.RequireRealm(auth.ClaimChannels)` wraps `handle` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112), Phoenix: the realm/token checks at internal/appengine/channels/ws.go:51-68 run before `websocket.Accept` at :71 — and only on the WebSocket branch (`wantsSSE` is tested first at ws.go:78), so 401/403 and the 200 SSE path win. Do NOT add 405: ServeMux's method-mismatch applies equally to every GET route in every spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-native-socket-sse-exact-accept [auto]: fix the SSE negotiation prose on `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml — the description (yaml:289-291), the `transport` parameter (yaml:328-334) and the `200` response (yaml:338-339) read as if any client asking for `text/event-stream` gets SSE, but the predicate `wantsSSE` is exact string equality (internal/appengine/stream/ws.go:151-153): measured, `Accept: text/event-stream, text/plain` and `Accept: text/event-stream;q=1.0` both return false and fall through to `websocket.Accept`, where a non-upgrade client gets the 426 above; `transport` selects SSE only for the literal `sse`, and `transport=websocket` / `transport=` (empty) both mean WebSocket. Say "exactly" for the header, say any other `Accept` falls through to the upgrade attempt, and do NOT claim the server rejects out-of-enum `transport` values — the enum is a client-side contract the handler never enforces. While there, give the `200` its real media type: the handler sets `Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive` and writes one `data: {...}` frame per event followed by a blank line (ws.go:123-126, 156-162). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-socket-security-scheme [auto]: `/astrate/v1/{realm}/socket` is documented as if public in docs/api/astrate_native_api.yaml — the spec has no `components.securitySchemes` at all and root `security: []` (yaml:21), and the operation declares no `security`, so a client generated from this spec sends no `Authorization` header and is answered 401 `{"errors":{"detail":"Unauthorized"}}` by `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112, bearer token read only from the Authorization header at middleware.go:131, 178-196). Add a scheme mirroring `a_aea` (docs/api/astarte_appengine_api.yaml:1436-1443: apiKey in header `Authorization`, described as a realm JWT carrying the `a_ch` claim) and set `security` on this operation only — root `security: []` stays so health/readiness/metrics/version remain documented as unauthenticated. If you model the Phoenix twin too, its credential is `?token=` in the query (internal/appengine/channels/ws.go:53-54), i.e. apiKey `in: query`, not header. Description notes worth carrying: the scheme is case-insensitive `Bearer` with an optional colon, and an unknown realm is 401 not 404 (middleware.go:68-72). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-native-socket-event-schema [auto]: document the native socket's event payload in docs/api/astrate_native_api.yaml — currently only a fenced `incoming_data` example inside the operation description (yaml:300-310), with no `components.schemas` entry anywhere (yaml:425-478), and that example is not the general shape: `wireEvent` marks `interface`, `path` and `value` `omitempty` (internal/appengine/stream/ws.go:51-58), so measured on a verbatim copy a lifecycle event marshals to `{"event":"device_connected","realm":"test","device_id":"dev1","timestamp":"2024-05-30T13:49:57.045Z"}` — no interface/path/value at all (and `value: false`/`value: 0` ARE emitted, because omitempty on an `any` drops only nil). Add a `NativeEvent` schema (event, realm, device_id, optional interface/path/value, timestamp) referenced from the operation, state that `event` is one of the internal/engine/stream/bus.go Kind* constants (bus.go:20-26), and record that this socket drops `InterfaceMajor`, `IP`, `ErrorName` and `ErrorMetadata`, which `stream.Event` carries (bus.go:39-65) — so a `device_error` arrives here with no reason, unlike the Phoenix twin, which marshals the whole struct (internal/appengine/channels/ws.go:406-416). That last sentence documents a code gap, not a spec one: do not add the fields; if they are wanted on the wire, that is a separate code task. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.5") | .body'
Astarte v1.3.5 release.

## Changes since v1.3.4

### Fixed

- [astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available

$ head -70 .mule/for-giulio.md
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **github-issues triage run, 2026-10-07: nothing proposable for the 25th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub.** The recipe's `--limit 40` command prints all **31** open issues (authoritative `--limit 100` agrees): **27** `mule-alarm` **#94–#120** plus the same four non-alarm, none of them proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (`milestone-4.0`, parked until the v3.0 queue clears), and one the standing instruction excludes from triage entirely. New since the 2026-10-06 run: **#120**, created 2026-10-07T11:03:45Z, "nothing has landed in 14h" — arithmetic pinned, not inferred: `.mule/.heartbeat` holds `1791318570` = **2026-10-06T20:29:30Z**, written 3s after the `docs-native-socket-security-scheme` done land (`87c18cf`, 20:29:27Z), and 11:03:45Z − 20:29:30Z = **14h34m**, so `age=$(( (now-last)/3600 ))` prints 14 — the title verbatim; `.mule/.alarmed` carries mtime 11:04Z, i.e. set by that filing. Genuine window, the same 14–16h overnight band #117–#119 measured. **Proposal, extending the 2026-10-06 one: close #94–#119, all twenty-six** — each is a one-day alarm superseded by the next and nothing reads issue state; leave **#120** (today's, real window) to self-expire.
  **New this run, and it is the 2026-10-01 pattern recurring: the switch is latched off through a silence that is real.** At 2026-10-07T16:55Z `.mule/.heartbeat` still read 2026-10-06T20:29:30Z — **20h**, with zero `2026-10-07` rows in `.mule/log.md` — but `.mule/.alarmed` (set by #120) short-circuits `check_pulse` at `tools/mule.sh:748`, so nothing can file again until a `done` land (mule.sh:578) or a `checked` run (mule.sh:540) calls `beat`, which is the only thing that clears the latch (mule.sh:737). Recipe ticks are committing today (`becf3cd` 11:25:35Z, `dde64e0` 11:38:44Z) but a recipe verdict is neither of those two paths, so they neither beat nor un-latch — measured today, not assumed. The `check_pulse`-before-`beat` ordering fix (move the call at `tools/mule.sh:633` to after the tick's `beat`) remains escalated above, still yours, still unactioned; #120 is the same arithmetic as #119, not fresh evidence for it either way.
  **#93 at 33 days.** `git rev-list --left-right --count origin/main...HEAD` = **4 / 697** today (4 / 679 on 10-06), and `git merge-base --is-ancestor 24ad5b8 origin/main` still fails — the rewrite is entirely on `mule/queue` / `origin/mule/queue`. The version-drift detail recorded on 2026-10-06 (the comment says `v1.4.0-rc.5` / "v1.3.3 being the newest stable tag"; upstream is `v1.4.0-rc.6` / `v1.3.5`, re-checked via `gh api repos/astarte-platform/astarte/tags` today) folds into that same review rather than becoming a line of its own.

---

- **COMPATIBILITY.md wording proposal for upstream v1.3.5 (newest stable, 2026-10-05; v1.4.0 is still rc.6-only) — a one-word delta on the open v1.3.2/v1.3.4 proposal below, nothing re-derived.** v1.3.5 is maintenance-only: `gh api repos/astarte-platform/astarte/releases` shows it published `2026-10-05T14:11:44Z`, not a prerelease, with a body of a single `### Fixed` line — "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available" — and `compare/v1.3.4...v1.3.5` is 4 commits / 35 files whose only non-`mix.exs`/`mix.lock` files are `.github/workflows/astarte-build-workflow.yaml`, `CHANGELOG.md`, `README.md`, `astarte-dashboard/package.json` and `docker-compose.yml`. Not one source file changed, so no route, MQTT topic, control message or interface-schema field moved (Astrate implements no `data_updater_plant` and no Erlang RPC server). **The only change a v1.3.5-aware doc carries: the v1.3.4 entry's §Infrastructure-differences sentence — "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target" — reads "…(**newest stable v1.3.5**) as the target".** Everything else stands unchanged: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"); that entry's "v1.4.0 is still rc.5-only" header is already superseded to rc.6 by the open rc.6 delta entry below, not here. This is the same one-word delta recorded on 2026-10-05 in this file ("newest stable v1.3.5"), folded in as asked rather than re-derived. Applied to nothing: `docs/COMPATIBILITY.md` untouched and still 1.2.2-targeted (it contains no `1.3` reference at all), `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). Raw: [v1.3.5](https://github.com/astarte-platform/astarte/releases/tag/v1.3.5).

---

- **`docs/site/appengine-api.md:114` — "All responses use the Astarte envelope format" is contradicted by the endpoints the page itself lists above it.** The "Response envelope" section (lines 112-126) opens with that sentence and then the `{ "data": ... }` / `{ "errors": { "detail": ... } }` shapes, but the live-stream endpoints documented at lines 94-110 never return one: the native socket pushes bare event JSON as WebSocket text frames or `data: {json}` SSE frames (internal/appengine/stream/ws.go:105-110, 156-162), `/astrate/v1/health` answers `{"status":"ok"}` (internal/observability/health.go:51), `/astrate/v1/readiness` answers `{"status":...,"checks":{...}}` (health.go:77), `/astrate/v1/metrics` answers Prometheus text under `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (metrics.go:42), and — measured today with the repo's own `coder/websocket` (probe in /tmp) — a WebSocket handshake failure on either socket is plain text under `text/plain; charset=utf-8` (`WebSocket protocol violation: Connection header "" does not contain Upgrade`), not JSON. REST appengine paths do use the envelope, so the sentence is right for the endpoints above "Device data" and wrong unqualified. Proposed wording (your voice): scope it — "All REST responses use the Astarte envelope format" — with one line noting the socket and observability surfaces answer un-enveloped. Page untouched.

- **`docs/site/appengine-api.md:102` — "Honours `a_ch` claims as room filters" imports upstream's room semantics onto a socket that has no rooms.** The native socket subscribes to the whole realm bus, narrowed only by the `device_id`/`interface` query parameters (internal/appengine/stream/ws.go:69-75); there is no room concept and no per-room filtering. What a_ch actually does here is gate the whole route: `mw.RequireRealm(auth.ClaimChannels)` requires the token to authorize `GET socket` via the REST verb-regex rule (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112 -> internal/auth/claims.go:164-171). The upstream-measured JOIN/WATCH partition rule DESIGN.md §4.2 describes (`AuthorizesChannel`, claims.go:197-204) is used only by the Phoenix socket's per-room join/watch checks (internal/appengine/channels/ws.go:278, 354). Consequence worth stating: a token whose a_ch list is a blanket `".*::.*"` — which upstream's Channels rule treats as authorizing nothing — authorizes the native socket, while the Phoenix socket, where room filters actually live, still reads JOIN/WATCH literally. Proposed wording (your voice): "guarded by the `a_ch` claim (GET on `socket`); `device_id`/`interface` narrow the stream". Page untouched.

- **github-issues triage run, 2026-10-06: nothing proposable for the 24th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented or closed on GitHub.** The recipe's `--limit 40` command prints all **30** open issues: **26** `mule-alarm` **#94–#119** plus the same four non-alarm, none of them proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (`milestone-4.0`, already escalated), and one the standing instruction excludes from triage entirely. New since the 2026-10-05 run: **#119**, created 2026-10-06T10:58:45Z, "nothing has landed in 15h" — a real window, not an artefact: the land before it was `230dd11` at 2026-10-05T21:04:09+02:00 (19:04Z), 15.9h earlier, and it was ended by today's first land `cbdd3e8` at 16:05:39+02:00 (14:05Z), three hours after the alarm. So it belongs to the genuine 14–16h band the 2026-10-05 entry measured, i.e. overnight silence between evening and afternoon batches. **Proposal, extending the 2026-10-05 one: close #94–#118, all twenty-five** — each is a one-day alarm superseded by the next, and nothing reads issue state; leave **#119** (today's, real window, already ended) to self-expire. The `check_pulse`-before-`beat` ordering fix (move the call at `tools/mule.sh:633` to after the tick's `beat`) is still escalated above, still yours, and still unactioned — today's alarm is not evidence for it either way.
  **The review queue moved 9 commits in 24h and main's four did not move again: `git rev-list --left-right --count origin/main...HEAD` is 4 / 679 today against 4 / 670 yesterday.** The four are main-side commits this branch lacks (`f1d0069`, `ca47b35`, `99743ca`, `2a033d4` — all ancestors of `origin/main`), and **#93's work is still entirely off main**: `8c61268` ("mule: log issue-93") and `24ad5b8` (the aclhook comment rewrite, whose commit message is the issue title verbatim) both satisfy `git merge-base --is-ancestor <sha> origin/main` → false, and `git branch -a --contains 24ad5b8` lists only `mule/queue` / `origin/mule/queue`. That is **32 days** of `mule-review` with `bash tools/mule.sh review` unread. I have not touched git and will not.
  **One detail to fold into that review rather than queue separately:** the rewritten comment at `internal/broker/aclhook.go` (the `control/keyAgreement` deny) says "…upstream's wire spec (082-key_agreement_protocol.md) only ships with v1.4.0-rc.5, v1.3.3 being the newest stable tag". The 082 citation, the "deliberately denied as Astrate's own choice" framing and the unchanged deny + `aclhook_test.go` case are exactly what #93 asked for, but the version reference has aged — upstream's newest stable is **v1.3.5** (2026-10-05) and newest overall `v1.4.0-rc.6` (2026-09-30). No separate task: one word in the same diff when someone finally reads it.

---

- **Milestone recipe run, 2026-10-06: v3.0 verified unchanged — still gated on upstream — and the master delta since the 2026-10-05 run is a closed set of 13 commits, every one wire-inert for Astrate. No issue filed, no `.mule/todo.md` line queued.** Release sweep is byte-identical to yesterday's (`gh api repos/astarte-platform/astarte/releases`): **no stable v1.4.0** — newest stable `v1.3.5` (2026-10-05), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition (your 2026-09-04 "wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) stays unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** — this milestone's whole ledger lives under `upstream-parity`, where the open set is unchanged: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`); **#47–#89** otherwise all closed. Non-alarm open issues total exactly four (plus the 25 `mule-alarm` #94–#118): the three above and **#1** (untouched per standing instruction). `APICompatVersion` still `"1.2.2"` (`internal/realm/service.go:588`, re-checked).
  **The new master commits — measured, and none is a gap.** `gh api .../commits?since=2026-10-05T14:20:00Z` returns 9 today (13 with the four landing between yesterday's run and that mark): `1da04832` + `a8a26d0f` "chore: forward port release-1.4 / release-1.3" are release-branch merges — CI/workflow consolidation, CHANGELOG, `mix.exs` bumps, and the appengine RPC refactor that *deletes* the `astarte_data_updater_plant` RPC client (Astrate implements no `data_updater_plant`); a router-level sweep of the 300-file forward-port finds **no new route or channel** — the only `_web/(router)` hit is `rooms_channel_test.exs`. `d85e0ed0` (docs) + `a6bd2c21` `fix(fdo): reject owner key names not accepted by OpenBao` are both FDO, i.e. **#78**'s parked body. `a51ab2e1` is an umbrella-CI refactor; `7628d2ed`/`ed1e99f4`/`bb17455a`/`ec020bf8` are `test(generators)`/`fix(generators)` Core-changeset validation — upstream test tooling, not Astrate surface. No HTTP route, MQTT topic, AMQP control message or interface-schema field moved.
  **Standing item unchanged.** `docs/UPSTREAM-EXPERIMENTAL.md`'s two rows (#67 `required`/`encrypted` mapping fields, #68 `async_operation=false`) still wait on **upstream 1.4 final** for their promoted-vs-deprecated call; nothing in today's commit set promotes or removes either and none of it is released, so the register is untouched.
  **Step 5 deliberately not taken**, same reason as the 21 runs before: the section's declared scope is **1.2.2 → 1.3/1.4** and only the 1.3 half is delivered. The 1.4 half stays open — #92 unanswered, both register rows unreconciled, `APICompatVersion` unbumped — so the tag must not be cut. Your call when a stable v1.4.0 ships: answer #92, reconcile the register, run the final-phase bump, then cut the tag.

- **astarte-upstream recipe run, 2026-10-05 (evening): no upstream movement that touches anything Astrate implements, so nothing was queued — no `.mule/todo.md` line, no `gh issue create`, nothing commented, closed or edited on GitHub.** One thing *is* new since this recipe's last run (`2026-10-04T20:40:52+02:00` ran; the `2026-10-05T04:29:10+02:00` run timed out): stable **v1.3.5**, published `2026-10-05T14:11:44Z`. The milestone recipe already recorded it as inert today; I re-measured it rather than trusting that, because "a release note is a claim, not a fact" is this recipe's own rule. `gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5` is **4 commits, 35 files** — `8d47574f` + `4c4222cf` `chore: update horde`, `c4c63ebe` `chore(rm): bump xandra to 0.19.4`, `8f0f5ef1` `chore: prepare v1.3.5` — and **every app/lib file in it is a `mix.exs`/`mix.lock` version bump**. The rest is `CHANGELOG.md` (+6: one `### Fixed` entry, "[astarte_data_updater_plant] Increase RPC server reliability"), `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `README.md`, `astarte-dashboard/package.json`. Not one source file changed, so no route, MQTT topic, BSON encoding, control message or interface-schema field moved: Astrate implements no `data_updater_plant` and no Erlang RPC server. Wire-inert, same shape as v1.3.2/v1.3.3/v1.3.4. **Master since `v1.4.0-rc.6` is still the single commit `731bcddc`** ("chore: forward port release-1.4"; `gh api .../commits?since=2026-10-05T09:34:00Z` returns it and nothing newer), the migration-relocation commit the milestone run measured today as inert — so both wire-bearing post-rc.6 commits (`b6d46ad4`, `958bb9fe`) still have their owners in the queue and nothing new joins them.
  **The parked-decisions row for #92 is still unmet, re-checked against the tag list rather than the release list** (this host has no `.mule/waiting-on.md` — absent from `HEAD`, already reported; I read it out of `origin/main` read-only): `gh api repos/astarte-platform/astarte/tags?per_page=100` gives newest stable **v1.3.5** and newest overall **v1.4.0-rc.6**, with `v1.4.0-rc.0`…`rc.6` and no stable `v1.4.0`. The row asks for a **stable** tag, which an `-rc.N` does not satisfy, so the twenty-second run ends the same way as the twenty-one before it.
  **The one delta to your open COMPATIBILITY.md wording proposal (the v1.3.2 one, folding v1.3.3/v1.3.4, further down this file): its version-reference sentence should read "newest stable v1.3.5" where it currently says "v1.3.4". Nothing else in it changes** — v1.3.5 adds no capability, so the same decision and the same "not yet emulated" list stand. Not applied: `docs/COMPATIBILITY.md` is untouched and still targets **1.2.2**, `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). I deliberately did **not** queue a fresh `compat-note-v1.3.5` line: the proposal is already open here, and a second line would re-derive it for the next run to find.

- **github-issues triage run, 2026-10-05: nothing proposable for the 23rd run, and the alarm pile's arithmetic is now fully pinned — 27 gaps over 8h since 2026-08-20 and **not one of them between 8h and 13h**, so the threshold cannot be a number. Two things are new: today's #118 was silenced by a land 5 minutes later, and #93's commit has now been unreviewed for 31 days.** 29 open issues (authoritative `--limit 100`: 25 alarms **#94–#118**, one new today, plus exactly four non-alarm — **#93** aclhook comment `mule-review`, **#92** keyAgreement parked, **#78** FDO `milestone-4.0`, **#1** untouched per standing instruction). Note the recipe's own `--limit 40` command printed **all 29** today, so the "silently dropped #94" parenthetical in the 2026-10-04 entry did not reproduce — treat the limit as adequate at this pile size. Zero machine-checkable candidates, so **no `.mule/todo.md` lines added, no `gh issue create`, and nothing commented, closed or edited on GitHub**.
  **My recount reproduces yesterday's numbers rather than replacing them, which is the useful part: the distribution is bimodal with an empty middle.** Over all **139** `mule: log` lands since 2026-08-20 (138 gaps): median **0.8h**, and the counts `> 8h` = `> 12h` = `> 13h` = **27** — so no gap anywhere in the record falls in (8h, 13h). Each of the 25 open alarms therefore measured a real 14.1–16.2h window (the only outlier is **#114** at 8.6h, filed 2026-09-27T20:09Z off the 11:35Z land that same day, which is the bimodality showing its other edge: a 20:00Z check before the evening batch starts). The 2026-10-04 entry's conclusion holds and hardens — a threshold in [8h, 13h) fires on every healthy day by construction; the real stall in the record (2026-09-28 → 10-01) is **74.7h**. **The proposed fix is unchanged and remains yours: move the `check_pulse` call from `tools/mule.sh:633` to after the tick's `beat`, so the switch cannot file against the work it is about to watch land.** Today's instance is the cleanest yet: **#118 filed 10:54:06Z** measuring a genuine 14.8h back to the 2026-10-04T20:06Z land, and the next land, `store-validatepipelinegraph-error-branches`, committed **10:59:49Z — 5 minutes later** (#117 was 23). Six lands today (10:59Z → 19:34Z), five `done` in `.mule/log.md`.
  **The 2026-10-01 entry's "latched off for three days" diagnosis is now superseded, and `.alarmed` explains how: `beat()` at `tools/mule.sh:737` does `rm -f "$MULE/.alarmed"` on every land, so one land per day is enough to re-arm the alarm every single morning.** `.mule/.alarmed` does not exist now, the switch re-armed on **2026-10-02**, and it has filed four more (#115–#118) since. The queue has not been idle for a month — my `git log --since=2026-08-20` count is 139 lands. Do not act on the September idle-streak framing any more; it described four real days and then the branch simply started moving again.
  **Proposal, extending the 2026-09-25, 09-26, 10-01, 10-03 and 10-04 ones: close #94–#117, all twenty-four.** Each is a one-day alarm superseded by the next, and as established on 2026-09-25 nothing reads issue state, so closing stays inert — but 24 open "the mule is dead" issues whose bodies point at a timer that does not exist is the noise that buried a real three-day stall in September. Leave **#118** (today's) open: like #117 it is the one whose window was unambiguously real.
  **New, and it is the same root cause as yesterday's merge finding: the review queue is not draining, so #93 has been waiting 31 days.** `8c61268` ("mule: log issue-93", aclhook comment rewrite) landed **2026-09-04T21:55:29+02:00**, is still **not** an ancestor of `origin/main`, and `git branch -a --contains 8c61268` lists only `mule/queue` / `origin/mule/queue` — a month of `mule-review` with no reviewer, while `bash tools/mule.sh review` sits unread. Divergence is now **4 / 670** (`git rev-list --left-right --count origin/main...HEAD`), the same 4 unpushed-here commits as yesterday (`f1d0069`, `ca47b35`, `99743ca`, `2a033d4`, all 2026-09-04 21:40–21:52) with `git merge-base --is-ancestor f1d0069 HEAD` → false. Yesterday measured this at 4 / 643: **this branch moved 27 commits in 24h and main's four did not move at all**, i.e. the review step has now missed a full day of landings on top of the month it has already missed. That single merge is what retires the pile's misleading bodies *and* unblocks #93; I have not touched git and will not.

---

- **Milestone recipe run, 2026-10-05: v3.0's gate is unchanged, but the release sweep finally has something new to swallow — upstream shipped stable v1.3.5 today, seven hours ago, and it is wire-inert for Astrate. No issue filed, no task line queued.**  **The new release, measured.** `gh api repos/astarte-platform/astarte/releases` now shows a stable **v1.3.5** at `2026-10-05T14:11:44Z` where every prior run of this recipe (2026-09-10 onward, 20+ entries in this file) saw `v1.3.4` (2026-09-18) as newest. `compare/v1.3.4...v1.3.5` is **4 commits, 35 files**: two `chore: update horde` (`8d47574f`, `4c4222cf`), `chore(rm): bump xandra to 0.19.4` (`c4c63ebe`), `chore: prepare v1.3.5` (`8f0f5ef1`). The file list is `mix.exs`/`mix.lock` version bumps across all apps and libs, `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `astarte-dashboard/package.json`, `doc/mix.exs`, `README.md` — plus a **single** CHANGELOG entry, `## [1.3.5] - 2026-10-05` / `### Fixed` / "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available". That is the whole delta and none of it is Astarte-facing: Astrate implements no data_updater_plant and no Erlang RPC server (the only `rpc` hit under `internal/` is a comment at `internal/pairing/http.go:300` about AMQP reply-code conflicts), so no HTTP route, MQTT topic, AMQP control message or schema field moved. Same shape as the v1.3.2/v1.3.3/v1.3.4 maintenance releases already recorded here, so it is **not a v3.0 gap** and there is nothing to file.
  **#92's parking condition is still unmet, now for the twenty-first run.** Your 2026-09-04 decision ("wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) — newest stable is `v1.3.5` (today), newest overall remains `v1.4.0-rc.6` (2026-09-30, prerelease). The post-rc.6 master delta has grown by exactly one commit: `731bcddc` "chore: forward port release-1.4" (2026-10-05T09:34Z), and it is the housekeeping/FDO migration relocation the 2026-10-03 run already measured inside rc.5→rc.6 (six `astarte_housekeeping/priv/migrations/{astarte,realm}/*.sql` files removed, table creation moved into migrations) plus a DUP-internal `apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0011_remove_replacement_data.ex` (+31) and `database.ex` (+2/-1). No Realm Management, AppEngine, Pairing or Trigger wire surface. **The milestone ledger is unchanged:** `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** (this milestone's whole ledger lives under `upstream-parity`, as the 2026-09-10 entry established), of #47–#89 **all remain closed**, and the authoritative open set is 29 issues — 25 alarms (#94–#118, one new today: **#118**, that recipe's business) and exactly four non-alarm: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** (untouched per standing instruction). `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`, re-checked), per **#90's** frozen decision.
  **The standing item is unchanged and I am not re-litigating it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries its two rows, both tagged `1.4 experimental` with trigger "upstream 1.4 final: promoted or removed" — the mis-tagging the 2026-10-04 entry measured out of upstream's CHANGELOG (row #68's `async_operation` feature is 1.0.2 and stable for four years; row #67's `required` half is a plain 1.4.0-rc.0 addition, only its `encrypted` half being genuinely experimental and already tracked by #92/#93). v1.3.5 touches neither feature, so it changes nothing about your four pending decisions (a)–(d) in that entry, and I have not edited the register.
  **Step 5 deliberately not taken.** I am *not* proposing "milestone v3.0 looks complete, verify and cut the tag", for the same reason as the twenty runs before this one: the section's declared scope is **1.2.2 → 1.3/1.4** and only the 1.3 half is delivered (now including v1.3.5, which is inert anyway). The 1.4 half is open — #92 unanswered, both experimental rows unreconciled, `APICompatVersion` not bumped — so the tag must not be cut. Your call when a stable v1.4.0 ships: answer #92, settle the register rows, run the final-phase bump, then cut the tag.
  No `gh issue create`, nothing commented, closed or edited on GitHub, and **no `.mule/todo.md` lines**: there is nothing machine-checkable left in the 1.3 line, every open item already has an owner, and the one new fact is a wire-inert maintenance release.

---

- **github-issues triage run, 2026-10-04: nothing proposable for the 22nd run, and the reason the alarm pile has never been actioned is finally located — `mule/queue` forked from `main` 27 minutes before your four commits landed there, so this Pi has been running for a month with a dead-man's switch, a dashboard and a parked-decision check that all point at things which do not exist here.** 28 open issues (authoritative `--limit 100` list: `1, 78, 92, 93, 94…117`; note the recipe's own `--limit 40` command printed 27 of these and silently dropped **#94**, which is open with `mule-alarm` — so treat "27 open" in the 2026-10-03 entry as 28). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`; `8c61268` re-verified still an ancestor of `origin/mule/queue` this run), **#92** keyAgreement (parked on a stable v1.4.0 — measured again today, still unmet: newest upstream stable `v1.3.4`, newest overall `v1.4.0-rc.6`, 2026-09-30), **#78** FDO (`milestone-4.0`, already escalated), **#1** untouched per standing instruction.
  **The branch divergence, measured.** `git rev-list --left-right --count origin/main...HEAD` → **4 643**: `main` is 4 commits ahead, this branch is 643 ahead, merge-base `215d409` at **2026-09-04T21:25:31+02:00**. Your four commits are all dated 2026-09-04 21:40–21:52 — `f1d0069`, `ca47b35`, `99743ca`, `2a033d4`, the last titled "parking a decision on an upstream release now has something that notices". They landed on `main` **after** this branch forked and were never reconciled. `git merge-base --is-ancestor 2a033d4 HEAD` → false. So for 30 days this working tree has been missing:
  - **The fixed alarm body.** `tools/mule.sh:755-756` on this host still tells whoever opens an idle alarm to run `systemctl list-timers mule.timer` and `journalctl -u mule.service --since '-1 day'`. Measured here: `systemctl list-timers --all` returns exactly three mule units — `mule-survey.timer`, `mule-upstream-watch.timer`, `mule-planner.timer`. **`mule.timer` does not exist**, and `mule.service` is not how ticks run. Ticks are the one-shot cron entries inside the `BEGIN mule-daily-schedule` block (`crontab -l`; today `31 13`, `36 13`, `17 13` CEST). Your `f1d0069` message says the reason in one line — "its absence reads as a dead mule to anyone following the map, which is exactly what happened" — and the fix has been sitting on `main` since the night it was written. **All 24 open alarm bodies point a reader at two commands that cannot work on this machine.** That is the mechanical reason the pile has sat for a month: the alarm told the only human who would read it to go and look at a timer that is not there, and he found nothing, and it filed again the next morning.
  - **The dashboard fix, same commit.** `.agents/skills/astrate-dashboard/SKILL.md` on this host still says `mule.timer` (30 min) at **:35, :126**, still runs `systemctl list-timers mule.timer …` at **:66, :69**, and still offers `systemctl stop|start mule.timer` at **:111** as the manual recovery for *"ferma/riattiva il mulo"*. On this Pi that command is a no-op. The map Giulio follows when the mule looks stuck points at a timer that does not exist.
  - **`.mule/waiting-on.md` and the recipe step that reads it** (`f1d0069`→`2a033d4`; `origin/main:.mule/recipes/astarte-upstream.md:17` is the line "*First, before anything else, check the parked decisions*"). The file is in `origin/main` and **absent from `HEAD`**; `rg waiting-on tools/mule.sh .mule/recipes/ .agents/` returns nothing on this host. So #92's own comment — "the weekly `mule-upstream-watch` job now reads that file as its first step every run" — is **not true of this host**. `mule-upstream-watch.timer` does exist: LAST **Mon 2026-09-28 04:02:35 CEST**, NEXT **2026-10-05 04:08:35 CEST**. It has therefore already run once without the parked-decisions step, and it runs again tomorrow morning without it. #92's re-parking has survived a month on hand-checks by the milestone and triage runs, not on the mechanism the issue was told about.
  - `tools/mule-survey.sh`'s push-with-retry (the fix for reports stranding unpushed, 2026-08-31→09-04) and `.mule/recipes/hygiene.md`'s dep-sweep honesty fix are likewise missing here.
  **This is not the mule's to fix and I have not touched it.** `cmd_review` states "The mule never merges" (`tools/mule.sh:794`); the `git checkout main && git merge --no-ff` at **:839** is text *printed for you* to run, not something the script executes. Reconciling `main` into `mule/queue` is a human merge, and both sides have moved (643 vs 4), so it is yours: `bash tools/mule.sh review` and merge at your convenience.
  **On the alarms themselves: they are not "false", and I should correct that framing.** The 2026-10-03 entry called #115/#116 "provably false" on the grounds that work had landed and been pushed on all three days. The gaps were real; what was false is the alarm's *conclusion*, not its arithmetic. Measured from the commit dates of every `[auto]` land since 2026-09-05 (115 lands, 114 gaps): **median gap 0.8h, and exactly 24 gaps exceed the 8h threshold — one per open alarm issue, #94–#117, a 1:1 correspondence.** The distribution is bimodal with nothing in the middle: the **smallest** gap above 8h is **14.6h** and the largest sub-8h gap is far below it. The mechanism is now exact rather than inferred: `check_pulse` is called at **`tools/mule.sh:633`, before the tick attempts any task**, and today's first cron entry is `17 13 4 10` = **13:17 CEST = 11:17Z** — the minute #117 was filed (11:17:06Z). **The switch fires on the very tick that is about to do the work which would have silenced it**, and the `beat` at the end of that tick (mule.sh:737) then deletes `.alarmed`. #117 is the clean case: it fired at 11:17:06Z measuring a genuine 15h window back to `f84bbb5` (2026-10-03T19:50:23Z), and the next land, `e728217`, committed at **11:40:11Z — 23 minutes later**. The four-day stall the 2026-10-01 entry diagnosed is over: five task attempts today, three `done` (`e728217`, `5044143`, `def0af6`), two `blocked`.
  **This should let you retire the threshold question rather than answer it.** The 2026-09-26 entry put two options to you (raise `MULE_IDLE_ALARM_HOURS` to ~24, or pulse `.mule/log.md` for a whole day with no `done` row). The gap histogram now says a threshold **cannot** work: healthy days reach **23.1h** (09-24→09-25) and **20.5h** (09-17→09-18) with no stall, while the three real stalls were **44.9h, 49.3h and 74.7h**. The only separating cut is ~36h, which would tolerate three days of silence before speaking — so the option worth taking is the second one, or the cheaper structural fix: **move the `check_pulse` call from mule.sh:633 to the end of the tick**, after the task has landed, so a healthy first-tick-of-the-day can never file against itself. That is a one-line change plus a shell test, and it is the only one of the three that does not trade a false positive for a false negative. It stays yours because it changes the mule's own safety behaviour, and because it lives in the same branch that needs the merge above — say the word and it becomes a queue line.
  **Proposal, extending the 2026-09-25, 09-26, 10-01 and 10-03 ones: close #94–#116, all twenty-three.** Each was a one-day alarm superseded by the next, and as established on 2026-09-25 nothing reads issue state, so closing stays inert — but 24 open "the mule is dead" issues is the noise that buried a real three-day stall in September, and the body of each one currently sends you to a timer that does not exist. Leave **#117** (today's) open: it is the one alarm whose window was unambiguously real, so it is the one worth reading before closing the other twenty-three.

---

$ grep -n "astarte-upstream\|upstream" .mule/log.md | tail -20; echo "=== log head ==="; head -40 .mule/log.md
=== log head ===
# Mule log

One row per task attempt, written by `tools/mule.sh`. This is the record of what the cheap
layer can actually do — read it before deciding whether a kind of task is worth delegating.

`secs` is the honest signal: a task that used most of its 900s budget was too big.

| date | task | outcome | secs | note |
| --- | --- | --- | --- | --- |
| 2026-09-04 | purge-properties-compression-capabilityauto | done | 86s | 6f7a3d6 |
| 2026-09-04 | broker-acl-coldstart-introspection-miss | done | 629s | d201db4 |
| 2026-09-04 | issue-93 | done | 108s | 24ad5b8 |
| 2026-09-04 | broker-disconnect-device-zombie-session | blocked | 117s | wrote nothing |
| 2026-09-04 | broker-offline-acl-tests | blocked | 208s | wrote nothing |
| 2026-09-04 | broker-onconnect-doc-comment | blocked | 438s | tests failed: --- FAIL: TestMQTTSink_Retained (0.02s) |
| 2026-09-05 | empty-introspection-verification | blocked | 106s | wrote nothing |
| 2026-09-05 | probe-trigger-install-notification-delay | blocked | 308s | wrote nothing |
| 2026-09-05 | compat-note-v132 | done | 73s | 4bc3e1a |
| 2026-09-05 | docs-sync-pairing-health-path | done | 186s | d73e225 |
| 2026-09-05 | docs-sync-pairing-register-404 | done | 106s | 44a7cae |
| 2026-09-05 | swagger-httptest-coverage | done | 267s | ba79b34 |
| 2026-09-05 | probe-property-resend-encoding | blocked | 100s | wrote nothing |
| 2026-09-05 | compat-note-v133 | blocked | 44s | wrote nothing |
| 2026-09-05 | flow-validate-source-sink | blocked | 360s | tests failed: --- FAIL: TestMQTTSink_Retained (0.02s) |
| 2026-09-05 | flow-validate-dead-source-sink-recompute | done | 200s | 7445176 |
| 2026-09-06 | docs-sync-rm-datastream-retention-endpoint | done | 232s | 0574fd9 |
| 2026-09-06 | docs-sync-rm-interfaces-detailed-param | done | 352s | 4094145 |
| 2026-09-06 | housekeeping-tests | done | 455s | 7583a7d |
| 2026-09-06 | probe-props-resend-error-triggers | blocked | 331s | wrote nothing |
| 2026-09-06 | server-data-trigger-bus | done | 827s | e46ce11 |
| 2026-09-06 | docs-sync-appengine-by-alias-endpoints | done | 144s | 35359df |
| 2026-09-06 | docs-sync-appengine-group-endpoints | done | 365s | 6ba74b8 |
| 2026-09-06 | docs-sync-appengine-get-group-device | done | 114s | 1aa98d6 |
| 2026-09-07 | docs-sync-appengine-query-params-status | done | 547s | bff7e1b |
| 2026-09-07 | docs-sync-appengine-group-patch-status | done | 75s | 81ced8f |
| 2026-09-07 | docs-sync-appengine-data-422-interface-level | done | 153s | f53d1ee |
| 2026-09-07 | appengine-snapshot-ignores-query-params | blocked | 327s | wrote nothing |
| 2026-09-07 | appengine-group-token-roundtrip-test | blocked | 306s | lint failed: internal/appengine/groups_token_test.go:14:5: redefines-builtin-id: redefinition of the built-in function max (revive) |
| 2026-09-07 | docs-sync-appengine-data-output-params | done | 215s | 262fdd5 |
| 2026-09-08 | probe-required-mapping-flag | done | 426s | db9d06e |

$ grep -n "compat-note\|probe-encrypted" .mule/todo.md | head; echo "=== check for-giulio for 10-07 upstream run ==="; grep -n "astarte-upstream recipe run, 2026-10" .mule/for-giulio.md | head
3:- [x] compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).
115:- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
120:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
138:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
177:- [x] compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.
248:- [x] compat-note-v1.4-rc.6 [auto]: record in `.mule/for-giulio.md` that v1.4.0-rc.6 is wire-inert for Astrate, as a delta on the open v1.3.2/v1.3.4 wording proposal rather than a re-derivation. Measured, the whole non-FDO delta is inert: DUP RPC-server reliability is OTP supervisor/GC tuning; the Dashboard's new `unknown_status` is computed client-side from `last_connection && !last_disconnection && !isConnected` (models/Device/index.ts:213-215) and both timestamps are already on Astrate's device shape (internal/appengine/service.go:118), so a device whose disconnection was never recorded reads as "Unknown status" rather than "Disconnected"; the `on_device_registered` dispatch moved from `Agent.register_device` to `Engine.register_device` (pairing/engine.ex:157-165, 186-193) so it also covers the FDO path, but the Agent path is byte-for-byte the same rule and Astrate already fires it on every successful registration (internal/pairing/service.go:221-223); the `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"); and `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma. So the proposal's "v1.4.0 is still rc.5-only" reference becomes rc.6, and its "FDO authentication (pairing, disabled by default)" line understates the 1.4 line, which ships a full ownership-voucher surface Astrate does not emulate (#78, parked). Do not edit docs/COMPATIBILITY.md and do not restate the standing proposal.
282:- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
283:- [x] compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (newest stable, 2026-10-05, maintenance-only: mix.exs/mix.lock version bumps plus a single CHANGELOG "astarte_data_updater_plant RPC server reliability" entry, no wire/API surface change; v1.4.0 is still rc.6-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2/v1.3.4 wording proposal there already stands and its version-reference delta was recorded on 2026-10-05 as reading "newest stable v1.3.5", so fold this in as that one-word delta rather than re-deriving the proposal.
=== check for-giulio for 10-07 upstream run ===
38:- **astarte-upstream recipe run, 2026-10-05 (evening): no upstream movement that touches anything Astrate implements, so nothing was queued — no `.mule/todo.md` line, no `gh issue create`, nothing commented, closed or edited on GitHub.** One thing *is* new since this recipe's last run (`2026-10-04T20:40:52+02:00` ran; the `2026-10-05T04:29:10+02:00` run timed out): stable **v1.3.5**, published `2026-10-05T14:11:44Z`. The milestone recipe already recorded it as inert today; I re-measured it rather than trusting that, because "a release note is a claim, not a fact" is this recipe's own rule. `gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5` is **4 commits, 35 files** — `8d47574f` + `4c4222cf` `chore: update horde`, `c4c63ebe` `chore(rm): bump xandra to 0.19.4`, `8f0f5ef1` `chore: prepare v1.3.5` — and **every app/lib file in it is a `mix.exs`/`mix.lock` version bump**. The rest is `CHANGELOG.md` (+6: one `### Fixed` entry, "[astarte_data_updater_plant] Increase RPC server reliability"), `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `README.md`, `astarte-dashboard/package.json`. Not one source file changed, so no route, MQTT topic, BSON encoding, control message or interface-schema field moved: Astrate implements no `data_updater_plant` and no Erlang RPC server. Wire-inert, same shape as v1.3.2/v1.3.3/v1.3.4. **Master since `v1.4.0-rc.6` is still the single commit `731bcddc`** ("chore: forward port release-1.4"; `gh api .../commits?since=2026-10-05T09:34:00Z` returns it and nothing newer), the migration-relocation commit the milestone run measured today as inert — so both wire-bearing post-rc.6 commits (`b6d46ad4`, `958bb9fe`) still have their owners in the queue and nothing new joins them.

$ gh api 'repos/astarte-platform/astarte/commits?since=2026-10-06T10:00:00Z&per_page=30' --jq '.[] | "\(.sha[0:8])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"'
171c1bd7	2026-10-06T15:57:45Z	chore: forward-port astarte_vmq_plugin from release-1.4
1da04832	2026-10-05T14:57:19Z	chore: forward port release-1.4
a8a26d0f	2026-10-05T14:40:35Z	chore: forward port release-1.3
d85e0ed0	2026-10-06T10:02:32Z	docs: sync FDO workflow to the current behavior (#2241)

$ gh api repos/astarte-platform/astarte/commits/171c1bd7 --jq '.files[] | "\(.status)\t+\(.additions)/-\(.deletions)\t\(.filename)"' | head -40; echo "=== msg ==="; gh api repos/astarte-platform/astarte/commits/171c1bd7 --jq '.commit.message'
modified	+85/-3	CHANGELOG.md
removed	+0/-255	apps/astarte_vmq_plugin/CHANGELOG.md
modified	+3/-1	apps/astarte_vmq_plugin/Dockerfile
modified	+2/-1	apps/astarte_vmq_plugin/lib/astarte_vmq_plugin/application.ex
modified	+11/-0	apps/astarte_vmq_plugin/lib/astarte_vmq_plugin/config.ex
added	+37/-0	apps/astarte_vmq_plugin/lib/astarte_vmq_plugin/health_http/handler.ex
added	+110/-0	apps/astarte_vmq_plugin/lib/astarte_vmq_plugin/health_http/route_injection_service.ex
modified	+5/-1	apps/astarte_vmq_plugin/mix.exs
modified	+4/-0	apps/astarte_vmq_plugin/priv/astarte_vmq_plugin.schema
modified	+3/-0	apps/astarte_vmq_plugin/rel/overlays/bin/start_vernemq
added	+59/-0	apps/astarte_vmq_plugin/test/astarte_vmq_plugin_health_http_handler_test.exs
added	+86/-0	apps/astarte_vmq_plugin/test/astarte_vmq_plugin_route_injection_service_test.exs
modified	+1/-0	apps/astarte_vmq_plugin/test/test_helper.exs
modified	+3/-1	mix.lock
=== msg ===
chore: forward-port astarte_vmq_plugin from release-1.4

Signed-off-by: frossq <francesco.asquini@secomind.com>

$ gh api 'repos/astarte-platform/astarte/compare/v1.4.0-rc.6...master' --jq '{ahead: .ahead_by, commits: (.commits|length), files: (.files|length)}'
{"ahead":1511,"commits":250,"files":300}

$ gh api repos/astarte-platform/astarte/commits/171c1bd7 --jq '.files[] | select(.filename|test("application.ex|config.ex|schema|start_vernemq|CHANGELOG.md$")) | "### \(.filename)\n\(.patch)"'
### CHANGELOG.md
@@ -10,6 +10,11 @@ and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.
 ### Added
 
 - [astarte_data_updater_plant] Generate `device_empty_cache_received` trigger events when a device sends `/emptyCache`
+- [astarte_vmq_plugin] Expose a new /healtz endpoint to check whether the AMQP connection toward RabbitMQ instance is ok and all the expected Mississippi EventsProducer workers are running.
+
+### Changed
+
+- [astarte_vmq_plugin] Pin vernemq to 2.2.0
 
 ## [1.4.0-rc.6] - 2026-09-28
 
@@ -30,6 +35,8 @@ Includes changes from v1.3.4
 - [fdo] Allow specifying device_id on ownership voucher upload
 - [fdo] The device is now immediately registered on ownership voucher upload
 - [astarte_data_updater_plant] Improve RPC server reliability
+- [astarte_vmq_plugin] Shard AMQP producers
+- [astarte_vmq_plugin] Increase RPC server reliability
 
 ### Fixed
 
@@ -68,10 +75,12 @@ Includes changes from v1.3.4
 ### Fixed
 
 - [astarte_data_updater_plant] Ensure RPC server is always available to clients. Resolved the issue where a temporary disconnection and reconnection of data_updater_plant to the cluster would make the RPC server inaccessible.
+- [astarte_vmq_plugin] Ensure RPC server is always available to clients. Resolved the issue where a temporary disconnection and reconnection of vmq_plugin to the cluster would make the RPC server inaccessible.
 
 ### Changed
 
 - [astarte_data_updater_plant] Ensure memory is properly garbage collected
+- [astarte_vmq_plugin] Pin vernemq to 2.1.2
 
 ## [1.4.0-rc.3] - 2026-07-31
 
@@ -126,12 +135,15 @@ Includes changes from v1.3.4
 ### Changed
 
 - [astarte_data_updater_plant] Use mississippi consumer for data updater processes
+- [astarte_vmq_plugin] Update VerneMQ to 2.0.1
+- [astarte_vmq_plugin] Use mississippi as AMQP publisher
 
 ## [1.3.5] - 2026-10-05
 
 ### Fixed
 
 - [astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available
+- [astarte_vmq_plugin] Increase RPC server reliability. A new corner case has been fixed which would've resulted in the rpc server not being available
 
 ## [1.3.4] - 2026-09-17
 
@@ -141,6 +153,7 @@ Includes changes from v1.3.4
   uncollected binaries over time by forcing more frequent full sweep garbage collections
   on them.
 - [astarte_data_updater_plant] Ensure the RPC server is always available to clients.
+- [astarte_vmq_plugin] Ensure the RPC server is always available to clients. Previously, a temporary disconnection and reconnection of VerneMQ to the cluster would make the RPC server inaccessible.
 
 ## [1.3.3] - 2026-08-07
 
@@ -150,7 +163,7 @@ Includes changes from v1.3.4
 
 ## [1.3.2] - 2026-07-14
 
-## Fixed
+### Fixed
 
 - Avoid crash on invalid properties message
 
@@ -242,6 +255,7 @@ Includes changes from v1.3.4
 - Allow devices with empty introspection
 - Devices can now declare support for optional Astarte MQTT v1 features to Astarte via capabilities
 - Support for `purge_properties_compression_format` capability. possible values are `zlib` (default) and `plaintext`
+- [astarte_vmq_plugin] Add the option to enable keepalive for scylladb connections, using the environment variable `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__ENABLE_KEEPALIVE`. Defaults to `true`
 
 ### Changed
 
@@ -291,6 +305,7 @@ Includes changes from v1.3.4
 - [astarte_data_updater_plant] Increase device process resiliency: avoid restarting the whole supervision tree when one device/amqp connection crashes
 - [astarte_housekeeping] Fix crashes in migrator
 - [astarte_realm_management] Fix corner case during for the installation of interfaces without data retention ttl
+- [astarte_vmq_plugin] Avoid race conditions preventing correct processing of device deletion requests
 
 ## [1.2.1] - 2026-03-12
 
@@ -302,6 +317,10 @@ Includes changes from v1.3.4
 
 ## [1.2.1-rc.1] - 2026-02-13
 
+### Added
+
+- [astarte_vmq_plugin] Add the option to enable keepalive for scylladb connections, using the environment variable `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__ENABLE_KEEPALIVE`. Defaults to `true`
+
 ### Fixed
 
 - [astarte_realm_management] Bug where devices got stuck in the "in deletion" status: [#1493](https://github.com/astarte-platform/astarte/issues/1493).
@@ -337,6 +356,7 @@ Includes changes from v1.3.4
 - Changed the database driver from CQEx (unmantained) to (E)xandra
 - [astarte_trigger_engine] avoid exposing **unknown_fields** in mustache templates
 - [astarte_trigger_engine] properly handle incoming introspection events
+- [astarte_vmq_plugin] RPC now uses erlang clustering instead of `astarte_rpc`
 
 ### Fixed
 
@@ -358,12 +378,16 @@ Includes changes from v1.3.4
 
 ## [1.2.1-alpha.0] - 2025-04-10
 
+### Added
+
+- [astarte_vmq_plugin] Allow to set the Erlang cookie via the `RELEASE_COOKIE` env var. Default to `vmq` for backwards compatibility.
+
 ### Changed
 
 - Update the docker-compose configuration to allow both physical and virtual devices
   to connect to Astarte, provided that the devices and the host are on the same LAN.
 
-## Fixed
+### Fixed
 
 - [astarte_appengine_api] Correctly handle Cassandra `varchar`s.
 - [astarte_data_updater_plant] Correctly handle Cassandra `varchar`s.
@@ -388,7 +412,7 @@ Includes changes from v1.3.4
 - Forward port changes from release-1.1 (connection failure when delivering
   triggers is handled as an error).
 
-## [1.2.0-rc.0] 11-06-2024
+## [1.2.0-rc.0] - 2024-06-11
 
 ### Added
 
@@ -420,6 +444,31 @@ Includes changes from v1.3.4
 - [astarte_realm_management_api] Allow to read realm's maximum datastream
   storage retention period with the `/config/datastream_maximum_storage_retention`
   endpoint.
+- [astarte_vmq_plugin] The plugin now accesses the Astarte database. The following
+  env variables have been added:
+  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__NODES`
+    (defaults to `localhost:9042`)
+  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__USERNAME`
+    (defaults to `cassandra`)
+  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__PASSWORD`
+    (defaults to `cassandra`)
+  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__POOL_SIZE`
+    (defaults to 10)
+  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__SSL_ENABLED`
+    (defaults to `false`)
+  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__SSL_DISABLE_SNI`
+    (defaults to `true`)
+  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__SSL_CUSTOM_SNI`
+  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__SSL_CA_FILE`
+- [astarte_vmq_plugin] Added support for device deletion. During deletion, a device is
+  disconnected and not allowed to reconnect until deletion ends.
+  Inflight messages are discarded. After deletion, a device must be
+  registered again in order to connect to Astarte.
+- [astarte_vmq_plugin] Added support for multiple Astarte instances sharing the same database,
+  the following env variable has been added:
+  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__ASTARTE_INSTANCE_ID`
+    (defaults to ``)
+- [astarte_vmq_plugin] Added support for `capabilities` message topic at `/<realm name>/<device name>/capabilities`
 
 ### Changed
 
@@ -438,6 +487,7 @@ Includes changes from v1.3.4
 - BREAKING: [astarte_realm_management] do not allow installation of interfaces
   where database_retention_ttl exceeds the realm's maximum datastream storage
   retention period, if set.
+- [astarte_vmq_plugin] Update VerneMQ to master (1cc57fa) to support OTP 26.
 
 ## [1.1.2] - Unreleased
 
@@ -501,6 +551,7 @@ Includes changes from v1.3.4
 - Update Elixir to 1.14.5 and Erlang/OTP to 25.3.2.
 - [astarte_data_updater_plant] Use the `internal` event type for Astarte
   internal messages. (e.g. device heartbeat).
+- [astarte_vmq_plugin] Use the `internal` event type for device heartbeat.
 
 ### Fixed
 
@@ -526,6 +577,7 @@ Includes changes from v1.3.4
 
 - [astarte_appengine_api] Return empty data instead of error when querying `properties` interfaces
   which are not fully populated. Fix [531](astarte-platform#531).
+- [astarte_vmq_plugin] Correctly serialize disconnection/reconnection events if VerneMQ hooks are called in the wrong order. Fix https://github.com/astarte-platform/astarte/issues/668.
 
 ## [1.0.6] - 2024-04-23
 
@@ -580,6 +632,15 @@ Includes changes from v1.3.4
 - [astarte_data_updater_plant] Check for device existence before installation or deletion
   of volatile triggers.
 
+### Fixed
+
+- [astarte_vmq_plugin] Do not let VerneMQ container start unless the CA cert is retrieved from CFSSL.
+- [astarte_vmq_plugin] Prevent the connection from timing out when the client takes more than 5 seconds to perform the SSL handshake
+
+### Security
+
+- [astarte_vmq_plugin] Rebuild official docker image (updates OTP to 23.3.4.17), in order to fix CVE-2022-37026.
+
 ## [1.0.3] - 2022-07-04
 
 ### Fixed
@@ -656,6 +717,7 @@ Includes changes from v1.3.4
 - [astarte_data_updater_plant] Don't crash when receiving `binaryblobarray` and `datetimearray`
   values.
 - Update Cyanide BSON library, in order to fix crash when handling ill-formed BSON arrays.
+- [astarte_vmq_plugin] Do not override VerneMQ config `max_message_rate` value.
 
 ## [1.0.0] - 2021-06-30
 
@@ -667,6 +729,7 @@ Includes changes from v1.3.4
 
 - Document future removal of Astarte Operator's support for Cassandra.
 - Log application version when starting.
+- [astarte_vmq_plugin] Log plugin version when the application is starting.
 
 ### Fixed
 
@@ -719,6 +782,7 @@ Includes changes from v1.3.4
 - Rename device `metadata` to `attributes`. _This requires a manual intervention on the database_,
   see the [Schema Changes](https://docs.astarte-platform.org/1.0/090-database.html#schema-changes)
   documentation for additional information.
+- [astarte_vmq_plugin] Do not authorize non-devices blindly in `auth_on_publish` and `auth_on_subscribe`.
 
 ## [1.0.0-beta.1] - 2021-02-16
 
@@ -754,6 +818,7 @@ Includes changes from v1.3.4
   bumping an interface minor.
 - Remove postgresql dependency in `docker-compose`, make CFSSL stateless.
 - Update Operator's documentation for install/upgrade/uninstall procedures.
+- [astarte_vmq_plugin] Default data_queue_count to 128.
 
 ## [1.0.0-alpha.1] - 2020-06-19
 
@@ -791,6 +856,10 @@ Includes changes from v1.3.4
   pass the `device_id` or `group_name` key inside the `simple_trigger`.
 - [data_updater_plant] Add support for device-specific and group-specific triggers.
 - Add support for device error triggers.
+- [astarte_vmq_plugin] Send a periodic heartbeat for every connected device.
+- [astarte_vmq_plugin] Support SSL for RabbitMQ connections.
+- [astarte_vmq_plugin] Reply with local and remote matches when a publish is requested.
+- [astarte_vmq_plugin] Allow configuring `max_offline_messages` and `persistent_client_expiration` with Docker env variables
 
 ### Removed
 
@@ -820,9 +889,14 @@ Includes changes from v1.3.4
 - [realm_management] Do not allow `/*` as match path when using `value_change` and
   `value_change_applied`. (workaround to https://github.com/astarte-platform/astarte/issues/513).
 - [trigger_engine] Update certifi to 2.5.3 (includes 2020-11-13 mkcert.org full CA bundle).
+- [astarte_vmq_plugin] Fix a bug where the plugin would remain unfunctional after suddenly disconnecting from RabbitMQ.
 
 ## [0.11.3] - 2020-09-24
 
+### Fixed
+
+- [astarte_vmq_plugin] Fix bug that prevented property unset
+
 ## [0.11.2] - 2020-08-14
 
 ### Added
@@ -844,6 +918,7 @@ Includes changes from v1.3.4
   this must be equal to the total number of queues in the Astarte instance.
 - [trigger_engine] Add `TRIGGER_ENGINE_AMQP_PREFETCH_COUNT` environment variable to set the
   prefetech count of AMQPEventsConsumer, avoiding excessive memory usage.
+- [astarte_vmq_plugin] Enhance docker build process
 
 ### Fixed
 
@@ -944,6 +1019,7 @@ Includes changes from v1.3.4
 - [housekeeping] Add database retention ttl and policy related columns (schema has been changed).
 - Allow specifying initial introspection when registering a device.
 - [realm_management] Trigger validation, checks that the interface is existing and performs validation on object aggregation triggers.
+- [astarte_vmq_plugin] Add support to multiple queues with consistent hashing
 
 ### Changed
 
@@ -1006,6 +1082,12 @@ Includes changes from v1.3.4
 
 ## [0.10.0] - 2019-04-16
 
+## [0.10.0-rc.1] - 2019-04-10
+
+### Fixed
+
+- [astarte_vmq_plugin] Re-enable SSL listener, which broke Docker Compose.
+
 ## [0.10.0-rc.0] - 2019-04-03
 
 ### Added
### apps/astarte_vmq_plugin/CHANGELOG.md
@@ -1,255 +0,0 @@
-# Changelog
-
-All notable changes to this project will be documented in this file.
-
-The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
-and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).
-
-## Unreleased
-
-## [1.4.0-rc.4] - 2026-08-06
-
-### Fixed
-
-- Ensure RPC server is always available to clients. Resolved the issue where a temporary disconnection and reconnection of vmq_plugin to the cluster would make the RPC server inaccessible.
-
-### Changed
-
-- Pin vernemq to 2.1.2
-
-## [1.4.0-rc.3] - 2026-07-31
-
-## [1.4.0-rc.2] - 2026-07-14
-
-## [1.4.0-rc.1] - 2026-07-13
-
-## [1.4.0-rc.0] - 2026-04-07
-
-### Changed
-
-- Update VerneMQ to 2.0.1
-- Use mississippi as AMQP publisher
-
-## [1.3.2] - 2026-07-14
-
-## [1.3.1] - 2026-07-10
-
-## [1.3.0] - 2026-05-05
-
-## [1.3.0-rc.2] - 2026-01-23
-
-## [1.3.0-rc.1] - 2026-01-23
-
-## [1.3.0-rc.0] - 2025-11-21
-
-### Added
-
-- Add the option to enable keepalive for scylladb connections, using the environment variable
-  `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__ENABLE_KEEPALIVE`. Defaults to `true`
-- Devices can now declare support for optional Astarte MQTT v1 features to Astarte via capabilities
-
-## [1.2.2] - 2026-04-27
-
-## [1.2.2-rc.0] 2026-04-08
-
-### Fixed
-
-- Avoid race conditions preventing correct processing of device deletion requests
-
-## [1.2.1] 2026-03-06
-
-## [1.2.1-rc.1] 2026-02-13
-
-### Added
-
-- Add the option to enable keepalive for scylladb connections, using the environment variable
-  `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__ENABLE_KEEPALIVE`. Defaults to `true`
-
-## [1.2.1-rc.0] 2025-08-22
-
-### Changed
-
-- RPC now uses erlang clustering instead of `astarte_rpc`
-
-## [1.2.1-alpha.0] - 2025-04-10
-
-### Added
-
-- Allow to set the Erlang cookie via the `RELEASE_COOKIE`
-  env var. Default to `vmq` for backwards compatibility.
-
-## [1.2.0] - 2024-07-01
-
-## [1.2.0-rc.0] - 2024-05-29
-
-### Added
-
-- The plugin now accesses the Astarte database. The following
-  env variables have been added:
-  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__NODES`
-    (defaults to `localhost:9042`)
-  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__USERNAME`
-    (defaults to `cassandra`)
-  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__PASSWORD`
-    (defaults to `cassandra`)
-  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__POOL_SIZE`
-    (defaults to 10)
-  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__SSL_ENABLED`
-    (defaults to `false`)
-  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__SSL_DISABLE_SNI`
-    (defaults to `true`)
-  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__SSL_CUSTOM_SNI`
-  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__CASSANDRA__SSL_CA_FILE`
-- Added support for device deletion. During deletion, a device is
-  disconnected and not allowed to reconnect until deletion ends.
-  Inflight messages are discarded. After deletion, a device must be
-  registered again in order to connect to Astarte.
-- Added support for multiple Astarte instances sharing the same database,
-  the following env variable has been added:
-  - `DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__ASTARTE_INSTANCE_ID`
-    (defaults to ``)
-- Added support for `capabilities` message topic at `/<realm name>/<device name>/capabilities`
-
-### Changed
-
-- Update Elixir to 1.15.7.
-- Update Erlang/OTP to 26.1.
-- Update VerneMQ to master (1cc57fa) to support OTP 26.
-
-## [1.1.1] - 2023-10-03
-
-## [1.1.0] - 2023-06-20
-
-## [1.1.0-rc.0] - 2023-06-09
-
-### Changed
-
-- Use the `internal` event type for device heartbeat.
-- Update Elixir to 1.14.5 and Erlang/OTP to 25.3.2.
-
-## [1.1.0-alpha.0] - 2022-11-24
-
-### Fixed
-
-- Correctly serialize disconnection/reconnection events if VerneMQ hooks are called in
-  the wrong order. Fix https://github.com/astarte-platform/astarte/issues/668.
-
-## [1.0.6] - 2024-04-18
-
-## [1.0.5] - 2023-09-25
-
-## [1.0.4] - 2022-09-26
-
-### Fixed
-
-- Do not let VerneMQ container start unless the CA cert is retrieved from CFSSL.
-- Prevent the connection from timing out when the client takes more than 5 seconds to perform the
-  SSL handshake
-
-### Security
-
-- Rebuild official docker image (updates OTP to 23.3.4.17), in order to fix CVE-2022-37026.
-
-## [1.0.3] - 2022-04-07
-
-## [1.0.2] - 2022-03-30
-
-## [1.0.1] - 2021-12-16
-
-### Fixed
-
-- Do not override VerneMQ config `max_message_rate` value.
-
-## [1.0.0] - 2021-06-30
-
-### Changed
-
-- Log plugin version when the application is starting.
-
-## [1.0.0-rc.0] - 2021-05-05
-
-## [1.0.0-beta.2] - 2021-03-24
-
-### Changed
-
-- Update Elixir to 1.11.4 and Erlang/OTP to 23.2
-- Do not authorize non-devices blindly in `auth_on_publish` and `auth_on_subscribe`.
-
-## [1.0.0-beta.1] - 2021-02-16
-
-### Changed
-
-- Default data_queue_count to 128.
-
-## [1.0.0-alpha.1] - 2020-06-19
-
-### Added
-
-- Send a periodic heartbeat for every connected device.
-- Support SSL for RabbitMQ connections.
-- Default max certificate chain length to 10.
-- Reply with local and remote matches when a publish is requested.
-- Allow configuring `max_offline_messages` and `persistent_client_expiration` with Docker env
-  variables
-
-## [0.11.4] - 2021-01-26
-
-### Fixed
-
-- Fix a bug where the plugin would remain unfunctional after suddenly disconnecting from RabbitMQ.
-
-## [0.11.3] - 2020-09-24
-
-### Fixed
-
-- Fix bug that prevented property unset
-
-## [0.11.2] - 2020-08-14
-
-### Added
-
-- Update Elixir to 1.8.2
-
-## [0.11.1] - 2020-05-18
-
-### Added
-
-- Enhance docker build process
-
-## [0.11.0] - 2020-04-13
-
-## [0.11.0-rc.1] - 2020-03-26
-
-## [0.11.0-rc.0] - 2020-02-26
-
-## [0.11.0-beta.2] - 2020-01-24
-
-## [0.11.0-beta.1] - 2019-12-26
-
-### Added
-
-- Add support to multiple queues with consistent hashing
-
-## [0.10.2] - 2019-12-09
-
-## [0.10.1] - 2019-10-02
-
-## [0.10.0] - 2019-04-16
-
-## [0.10.0-rc.1] - 2019-04-10
-
-### Fixed
-
-- Re-enable SSL listener, which broke Docker Compose.
-
-## [0.10.0-rc.0] - 2019-04-03
-
-## [0.10.0-beta.3] - 2018-12-19
-
-## [0.10.0-beta.2] - 2018-10-19
-
-## [0.10.0-beta.1] - 2018-08-27
-
-### Added
-
-- First Astarte release.
### apps/astarte_vmq_plugin/lib/astarte_vmq_plugin/application.ex
@@ -42,7 +42,8 @@ defmodule Astarte.VMQ.Plugin.Application do
       {Astarte.VMQ.Plugin.Publisher, [Config.registry_mfa()]},
       {Horde.Registry, [name: Registry.VMQPluginRPC, keys: :unique, members: :auto]},
       Astarte.VMQ.Plugin.RPC.Supervisor,
-      {Xandra.Cluster, Config.xandra_options!()}
+      {Xandra.Cluster, Config.xandra_options!()},
+      {Astarte.VMQ.Plugin.HealthHttp.RouteInjectionService, []}
     ]
 
     # See https://hexdocs.pm/elixir/Supervisor.html
### apps/astarte_vmq_plugin/lib/astarte_vmq_plugin/config.ex
@@ -84,6 +84,17 @@ defmodule Astarte.VMQ.Plugin.Config do
       end
 
     Application.put_env(:astarte_rpc, :amqp_connection, astarte_rpc_amqp_opts)
+
+    # Check if we want to inject custom HTTP endpoints in the Verne cowboy listener
+    inject_custom_routes =
+      Application.get_env(:astarte_vmq_plugin, :inject_custom_routes, "false")
+      |> to_string()
+      |> case do
+        "true" -> true
+        _ -> false
+      end
+
+    Application.put_env(:astarte_vmq_plugin, :inject_custom_routes, inject_custom_routes)
   end
 
   defp init_ssl_options(amqp_options) do
### apps/astarte_vmq_plugin/priv/astarte_vmq_plugin.schema
@@ -170,3 +170,7 @@
     {datatype, string}
   ]}.
   
+{mapping, "astarte_vmq_plugin.inject_custom_routes", "astarte_vmq_plugin.inject_custom_routes", [
+  {default, true},
+  {datatype, {enum, [true, false]}}
+]}.
### apps/astarte_vmq_plugin/rel/overlays/bin/start_vernemq
@@ -1,5 +1,8 @@
 #!/usr/bin/env bash
 
+# allow injection of custom HTTP routes when running inside the VerneMQ release container
+export DOCKER_VERNEMQ_ASTARTE_VMQ_PLUGIN__INJECT_CUSTOM_ROUTES="true"
+
 NET_INTERFACE=$(route | grep '^default' | grep -o '[^ ]*$')
 NET_INTERFACE=${DOCKER_NET_INTERFACE:-${NET_INTERFACE}}
 IP_ADDRESS=$(ip -4 addr show ${NET_INTERFACE} | grep -oE '[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}' | sed -e "s/^[[:space:]]*//" | head -n 1)

Read the rules, the recipe, `docs/COMPATIBILITY.md`, and the prior runs. Result:

**Releases** (`gh api repos/astarte-platform/astarte/releases`): newest stable `v1.3.5` (2026-10-05), newest overall `v1.4.0-rc.6` (2026-09-30). Both already fully handled — `v1.3.5` by the 2026-10-05 run of this recipe (measured wire-inert: 4 commits, 35 files, all `mix.exs`/`mix.lock` + one DUP RPC CHANGELOG line) and its `compat-note-v1.3.5` line is already in `.mule/todo.md` (`[x]`) with the wording delta already recorded in `.mule/for-giulio.md`. Re-queueing it would duplicate.

**Master since the last checks**: `gh api .../commits?since=2026-10-06T10:00:00Z` returns four, three already measured by earlier runs (`1da04832` → `probe-encrypted-mapping-enforcement`, `a8a26d0f`, `d85e0ed0`). One is new: `171c1bd7` "chore: forward-port astarte_vmq_plugin from release-1.4" — inspected the file list and patch: entirely inside `apps/astarte_vmq_plugin`, adding a `/healtz` health endpoint injected into VerneMQ's cowboy listener (`inject_custom_routes`), pinning VerneMQ 2.2.0, and backfilling that app's old CHANGELOG into the root one. Astrate runs an embedded mochi-mqtt broker with no VerneMQ and no AMQP/RabbitMQ, so no MQTT topic, payload encoding, pairing flow, API shape, trigger payload or interface-schema field moved. Not a problem Astrate has, so no `probe-` line either.

Nothing qualifies under the recipe's proposal rules, so I appended nothing.

```
Done: no upstream movement since v1.3.5 — nothing appended to .mule/todo.md
Files: none
Verified: gh api repos/astarte-platform/astarte/releases (v1.3.5 / v1.4.0-rc.6, both already handled) ; gh api .../commits?since=2026-10-06T10:00:00Z (4, 3 known + 171c1bd7) ; gh api .../commits/171c1bd7 file+patch inspection (vmq_plugin-only) -> pass
Unsure: nothing
Follow-ups: none
```
