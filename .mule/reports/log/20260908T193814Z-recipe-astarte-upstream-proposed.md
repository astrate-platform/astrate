slug: recipe-astarte-upstream
verdict: proposed
at:  eca0e08
ran: 2026-09-08T19:38:14Z on DietPi in 237s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md
→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

→ Read .mule/todo.md
Upstream latest is v1.4.0-rc.5 (2026-08-20), with v1.3.3 (2026-08-07) as the newest stable. Let me check what the previous run recorded.

$ rg -n "1\.3\.3|1\.4\.0|astarte-upstream|compat-note|upstream movement" /root/astrate-mule/.mule/ /root/astrate-mule/docs/ 2>/dev/null | head -50
/root/astrate-mule/.mule/todo.md:61:`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
/root/astrate-mule/.mule/todo.md:99:- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
/root/astrate-mule/.mule/todo.md:104:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
/root/astrate-mule/.mule/todo.md:110:- [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
/root/astrate-mule/.mule/todo.md:121:- [x] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
/root/astrate-mule/.mule/todo.md:122:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
/root/astrate-mule/.mule/for-giulio.md:16:  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
/root/astrate-mule/.mule/for-giulio.md:31:  experimental decision, escalated 2026-09-05, gated on a stable upstream v1.4.0), **#78 FDO**
/root/astrate-mule/.mule/for-giulio.md:37:  call when v1.4.0 goes stable: answer #92 (implement against the rc, re-park, or narrow #93
/root/astrate-mule/.mule/for-giulio.md:57:  awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the waiting-on
/root/astrate-mule/.mule/for-giulio.md:65:  stable v1.4.0?** (issue #92, `upstream-parity`/`upstream-experimental`). #51 closed
/root/astrate-mule/.mule/for-giulio.md:73:  fixed the stale ACL comment that cited them). But the spec ships only in `v1.4.0-rc.5` —
/root/astrate-mule/.mule/for-giulio.md:74:  `v1.3.3` is still the newest stable tag and Astrate targets 1.2.2 — so the document has
/root/astrate-mule/.mule/for-giulio.md:78:  re-park until v1.4.0 is a stable tag, or take only the narrow #93 fix. (Escalated again
/root/astrate-mule/.mule/for-giulio.md:88:- **COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable, 2026-07-14; v1.4.0 is still rc-only).** Astrate's doc and `APICompatVersion` still target upstream **1.2.2** (`internal/realm/service.go:588`); v1.3.0 (2026-05-06) introduced wire-surface changes Astrate does not yet emulate, so this is a decision — adopt v1.3.2 as the compatibility target (then update the doc + bump `APICompatVersion` together, per the bump rule) or keep 1.2.2 and add a "not yet emulated" note. Wire-relevant v1.3.0 deltas (release notes): **MQTT v1 capabilities** incl. `purge_properties_compression_format` (plaintext vs zlib purge — touches the `emptyCache`/`producer,properties` contract COMPATIBILITY.md deviation 1 documents); **empty introspection now allowed**; **device registration triggers** (pairing) and **device deletion started/completed triggers** (RM — the latter two already exist as Astrate deviation 9 emits both around the synchronous delete); **FDO authentication** (pairing, disabled by default); **realm-scoped health** — upstream v1.3 added `GET /pairing/v1/{realm}/health`, which Astrate already serves (`internal/pairing/http.go:78-81`, comment already says "upstream 1.3+"), so deviation 18's wording ("which upstream 404s") is now false against 1.3 and the note should be reworded either way. Proposed doc wording (for your approval, edit to taste): in §Infrastructure differences add a sentence — *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl.
/root/astrate-mule/.mule/log.md:18:| 2026-09-05 | compat-note-v132 | done | 73s | 4bc3e1a |
/root/astrate-mule/.mule/log.md:23:| 2026-09-05 | compat-note-v133 | blocked | 44s | wrote nothing |
/root/astrate-mule/.mule/log.md:41:| 2026-09-08 | compat-note-v14-rc | blocked | 326s | wrote nothing |
/root/astrate-mule/.mule/reports/log/20260907T112447Z-docs-sync-appengine-query-params-status-done.md:615:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
/root/astrate-mule/.mule/reports/log/20260907T112447Z-docs-sync-appengine-query-params-status-done.md:621:- [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
/root/astrate-mule/.mule/reports/log/20260905T165512Z-recipe-docs-sync-proposed.md:371: - [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:1:slug: recipe-astarte-upstream
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:10:→ Read .mule/recipes/astarte-upstream.md
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:14:v1.4.0-rc.5	2026-08-20T12:28:04Z
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:15:v1.3.3	2026-08-07T14:33:58Z
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:16:v1.4.0-rc.4	2026-08-06T13:40:20Z
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:17:v1.4.0-rc.3	2026-07-31T10:48:39Z
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:18:v1.4.0-rc.2	2026-07-14T10:16:22Z
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:22:$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:23:Sixth Astarte v1.4.0 release candidate.
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:25:## Changes since v1.4.0-rc.4
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:32:$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:33:v1.3.3 has an empty body. The most recent tag is `v1.4.0-rc.5` (2026-08-20). Let me understand the current state.
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:35:First, what version does Astrate target? From `COMPATIBILITY.md`, deviations reference "upstream 1.2.2" and "upstream 1.4", and `UPSTREAM-EXPERIMENTAL.md` tracks 1.4 experimental features. The latest stable tag is `v1.3.3`; the latest release overall is `v1.4.0-rc.5`.
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:50:c3cad68 mule: recipe astarte-upstream ran (85s)
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:68:- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:73:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:79:- [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:88:The previous run (`c3cad68 mule: recipe astarte-upstream ran (85s)`) already produced several `[auto]` task lines from the upstream recipe covering v1.3.0, v1.3.2, v1.3.3. Let me understand the current state fully.
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:90:The latest tag is now `v1.4.0-rc.5` (2026-08-20), newer than the previous stable `v1.3.3` (2026-08-07). But wait — Astrate already targets 1.4 experimental features (`UPSTREAM-EXPERIMENTAL.md` tracks #67 and #68 at 1.4 experimental level, and deviation 17 mentions "upstream 1.4 defaults to asynchronous"). And from the todo, the previous run already proposed `compat-note-v1.3.2` and `compat-note-v1.3.3`.
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:92:Let me check what the previous astarte-upstream run already covered. The todo lines already present from that recipe:
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:96:- `compat-note-v1.3.2 [auto]` (v1.3.2)
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:98:- `compat-note-v1.3.3 [auto]` (v1.3.3)
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:99:- `probe-props-resend-error-triggers [auto]` (v1.3.3)
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:101:So the previous run already covered v1.3.0, v1.3.2, v1.3.3. Now the newest tag is `v1.4.0-rc.5`, released 2026-08-20, which is newer than v1.3.3. 
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:107:Let me think about what the previous runs did. The previous run compared against v1.3.3 and generated the compat-note. Now the newest is v1.4.0-rc.5. The v1.4.0-rc.5 changes are minimal (two fixes in astarte_data_updater_plant about object value type validation and binaryblob validation). These are Elixir data_updater_plant internals.
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:109:Let me check the release notes history — what did earlier runs already cover for v1.4.0-rc.x? The `UPSTREAM-EXPERIMENTAL.md` register tracks #67 and #68. The previous recipe runs mentioned v1.4.0 as "still rc-only". 
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:111:Since v1.4.0-rc.5 is newer than what was previously tracked (v1.3.3), let me look at the release notes diffs for the rc releases. The previous run apparently only looked up to v1.3.3 (the newest stable at the time). Actually wait — the previous run happened and proposed compat-note-v1.3.3 saying "newest stable, 2026-08-07... v1.4.0 is still rc.5-only". So at that point v1.4.0-rc.5 already existed! The previous run already knew about rc.5.
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:113:So the question: did the previous run already examine v1.4.0-rc.5? It said "v1.4.0 is still rc.5-only" meaning it knew rc.5 was the latest rc. But did it propose anything for the rc.5 changes?
/root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md:115:Let me look at the release notes for the rc versions more carefully to understand what's covered. The rc.5 body only has two fixes. Let me check what the previous run proposed regarding v1.4.0-rc.x changes. The todo doesn't show any v1.4.0-rc proposals other than the experimental register stuff already handled.

The previous run (2026-09-07) already ran this recipe. Let me read its full report to avoid duplicating its proposals.
→ Read .mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md
Let me read the conclusion of the previous run's report to see exactly what it proposed.
→ Read .mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md [offset=490]
The previous run examined everything through v1.4.0-rc.5 and v1.3.3. Let me verify there's been no new upstream movement since, and check the report-file conventions.

$ ls /root/astrate-mule/.mule/reports/log/ | sort | head -40
20260904T190911Z-purge-properties-compression-capabilityauto-done.md
20260904T195018Z-broker-acl-coldstart-introspection-miss-done.md
20260904T195529Z-issue-93-done.md
20260904T195814Z-broker-disconnect-device-zombie-session-blocked.md
20260904T203003Z-broker-offline-acl-tests-blocked.md
20260904T210532Z-broker-onconnect-doc-comment-blocked.md
20260905T110211Z-empty-introspection-verification-blocked.md
20260905T111428Z-probe-trigger-install-notification-delay-blocked.md
20260905T112539Z-compat-note-v132-done.md
20260905T165512Z-recipe-docs-sync-proposed.md
20260905T171755Z-docs-sync-pairing-health-path-done.md
20260905T172234Z-docs-sync-pairing-register-404-done.md
20260905T173657Z-recipe-hygiene-proposed.md
20260905T174528Z-swagger-httptest-coverage-done.md
20260905T181053Z-recipe-milestones-proposed.md
20260905T181219Z-recipe-github-issues-proposed.md
20260905T184625Z-recipe-astarte-upstream-proposed.md
20260905T184903Z-probe-property-resend-encoding-blocked.md
20260905T190855Z-compat-note-v133-blocked.md
20260905T192010Z-recipe-code-review-proposed.md
20260905T192911Z-flow-validate-source-sink-blocked.md
20260905T201529Z-flow-validate-dead-source-sink-recompute-done.md
20260906T112006Z-recipe-docs-sync-proposed.md
20260906T113143Z-docs-sync-rm-datastream-retention-endpoint-done.md
20260906T171045Z-docs-sync-rm-interfaces-detailed-param-done.md
20260906T174031Z-recipe-hygiene-proposed.md
20260906T175748Z-housekeeping-tests-done.md
20260906T180442Z-recipe-milestones-proposed.md
20260906T181352Z-recipe-github-issues-proposed.md
20260906T183102Z-recipe-astarte-upstream-proposed.md
20260906T183952Z-probe-props-resend-error-triggers-blocked.md
20260906T185735Z-recipe-code-review-proposed.md
20260906T195827Z-server-data-trigger-bus-done.md
20260906T200431Z-recipe-docs-sync-proposed.md
20260906T201002Z-docs-sync-appengine-by-alias-endpoints-done.md
20260906T202744Z-docs-sync-appengine-group-endpoints-done.md
20260906T204333Z-docs-sync-appengine-get-group-device-done.md
20260907T020359Z-recipe-astarte-upstream-proposed.md
20260907T112447Z-docs-sync-appengine-query-params-status-done.md
20260907T112654Z-docs-sync-appengine-group-patch-status-done.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'; echo "==="; gh api repos/astarte-platform/astarte/releases?per_page=100 --jq '[.[].tag_name] | .[0:10]'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z
===
["v1.4.0-rc.5","v1.3.3","v1.4.0-rc.4","v1.4.0-rc.3","v1.4.0-rc.2","v1.4.0-rc.1","v1.3.2","v1.3.1","v1.3.0","v1.2.2"]

$ git rev-parse --short HEAD; echo "==="; git status --short; echo "==="; date -u +%Y-%m-%dT%H:%M:%SZ
eca0e08
===
===
2026-09-08T19:36:14Z

No new upstream tag: newest is still v1.4.0-rc.5 (2026-08-20), same as the last run. Everything from v1.3.0 through v1.4.0-rc.5 is already covered by prior `[auto]` task lines and compat-notes in the queue. Let me confirm the current HEAD for the evidence file and that git is clean.
→ Read .mule/reports/log/20260907T020359Z-recipe-astarte-upstream-proposed.md
No upstream movement. Let me check the header convention used by a prior complete-run report to write the evidence file consistently.

$ head -8 /root/astrate-mule/.mule/reports/log/20260906T183102Z-recipe-astarte-upstream-proposed.md; echo "===TAIL==="; tail -5 /root/astrate-mule/.mule/reports/log/20260906T183102Z-recipe-astarte-upstream-proposed.md
slug: recipe-astarte-upstream
verdict: proposed
at:  eaa2c37
ran: 2026-09-06T18:31:02Z on DietPi in 526s

> build · big-pickle

===TAIL===
**Verified:** `git status --short` → only `.mule/todo.md` modified.

**Unsure:** whether `sendConsumerProperties` (the purge-list path) is considered a "resend" in the same semantic category as upstream's per-interface value-resend; the upstream fix (device.ex `resend_all_interface_properties`) operates on per-interface property values, while Astrate's counterpart `resendServerProperties` (control.go:144) does the same per-interface value work but `sendConsumerProperties` (control.go:191) emits the purge-list which has different failure semantics (silently drops unknown interfaces at line 199 rather than failing per-interface). The probe task's wording covers both.

**Follow-ups:** one new task line (`probe-props-resend-error-triggers` at todo.md:110).

$ ls /root/astrate-mule/.mule/reports/log/ | grep astarte-upstream; echo "==="; for f in /root/astrate-mule/.mule/reports/log/*recipe-astarte-upstream*.md; do echo "--- $f"; sed -n '1,4p' "$f"; done
20260905T184625Z-recipe-astarte-upstream-proposed.md
20260906T183102Z-recipe-astarte-upstream-proposed.md
20260907T020359Z-recipe-astarte-upstream-proposed.md
20260907T175524Z-recipe-astarte-upstream-proposed.md
20260908T113422Z-recipe-astarte-upstream-proposed.md
===
--- /root/astrate-mule/.mule/reports/log/20260905T184625Z-recipe-astarte-upstream-proposed.md
slug: recipe-astarte-upstream
verdict: proposed
at:  6217b3c
ran: 2026-09-05T18:46:25Z on DietPi in 730s
--- /root/astrate-mule/.mule/reports/log/20260906T183102Z-recipe-astarte-upstream-proposed.md
slug: recipe-astarte-upstream
verdict: proposed
at:  eaa2c37
ran: 2026-09-06T18:31:02Z on DietPi in 526s
--- /root/astrate-mule/.mule/reports/log/20260907T020359Z-recipe-astarte-upstream-proposed.md
slug: recipe-astarte-upstream
verdict: proposed
at:  818a2ef
ran: 2026-09-07T02:03:59Z on DietPi in 85s
--- /root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md
slug: recipe-astarte-upstream
verdict: proposed
at:  08dee34
ran: 2026-09-07T17:55:24Z on DietPi in 127s
--- /root/astrate-mule/.mule/reports/log/20260908T113422Z-recipe-astarte-upstream-proposed.md
slug: recipe-astarte-upstream
verdict: proposed
at:  1e56fda
ran: 2026-09-08T11:34:22Z on DietPi in 305s

There's already a run of this recipe today (2026-09-08T11:34). Let me read it to see what it found, so I don't duplicate its proposals.
→ Read .mule/reports/log/20260908T113422Z-recipe-astarte-upstream-proposed.md

$ git log --oneline -12
eca0e08 mule: refill the queue
426d906 mule: recipe github-issues ran (64s)
54e453d mule: recipe milestones ran (244s)
4e42aff mule: recipe hygiene timed out
2dd6eea mule: log docs-sync-appengine-downsample-min
db506f9 mule: docs-sync-appengine-downsample-min [auto]: fix the `downsample_to` constraint in docs/api/astrate_appengine_api.yaml — the DataDownsample description says "(must be >= 2)" (yaml:1357) but parseQueryOpts rejects `downsample_to <= 2` with "must be greater than 2" (internal/appengine/http.go:592-593), so the documented minimum is off by one: readings of `2` (allowed per the doc) actually 422. Change the description to "(must be > 2)" (the 422 error shape is already documented). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
ce9d2ee mule: blocked docs-sync-appengine-data-set-422
74987a3 mule: log docs-sync-appengine-device-status-schema
00da133 mule: docs-sync-appengine-device-status-schema [auto]: fix the `DeviceStatus` schema in docs/api/astrate_appengine_api.yaml — the wire emits `id` (internal/appengine/service.go:106 `json:"id"`, decoded as `ds.ID` in mirror_test.go:156 and dashboard_compat_test.go:100/149) but the spec names the field `device_id` (yaml:1409) in both the schema and the list/detail examples (yaml:67, 75); `introspection` is documented as `map<string,integer>` (yaml:1438-1441) but each value is a `{major, minor}` object (introspectionEntry, service.go:123-127, introspectionView service.go:833-836); and the schema omits the six fields `groups`, `total_received_msgs`, `total_received_bytes`, `first_credentials_request`, `last_seen_ip` (omitempty), `previous_interfaces` (omitempty) that deviceStatusView emits (service.go:254-279). Rename `device_id`→`id`, correct `introspection` to `map<string,{major:int,minor:int}>`, add the missing fields, and fix the examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
56b56b1 mule: refill the queue
643c61b mule: recipe docs-sync ran (711s)
95dc630 mule: log store-latest-individual-test

The recipe already ran today (2026-09-08T11:34, at sha `1e56fda`) and appended `probe-required-mapping-flag` and `compat-note-v1.4-rc`; both already sit in the queue. Let me confirm the git state since then and that nothing upstream has moved.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | "\(.tag_name)\t\(.published_at)"' | head -12
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z
v1.4.0-rc.1	2026-07-14T08:53:27Z
v1.3.2	2026-07-14T10:07:26Z
v1.3.1	2026-07-13T11:50:05Z
v1.3.0	2026-05-06T15:30:50Z
v1.2.2	2026-04-27T16:39:28Z
v1.4.0-rc.0	2026-04-08T12:32:33Z
v1.3.0-rc.2	2026-04-08T16:34:59Z

This recipe already ran earlier today (2026-09-08T11:34, sha `1e56fda`) and appended both proposals (`probe-required-mapping-flag`, `compat-note-v1.4-rc`); upstream is unchanged — newest tag still v1.4.0-rc.5 (2026-08-20). Let me check how no-movement outcomes are recorded, then write the evidence file.
$ ls /root/astrate-mule/.mule/reports/log/ | grep -vE "astarte-upstream|docs-sync|hygiene|milestones|github-issues|code-review" | head; echo "===verdicts==="; grep -rh "^verdict:" /root/astrate-mule/.mule/reports/log/*.md | sort | uniq -c
20260904T190911Z-purge-properties-compression-capabilityauto-done.md
20260904T195018Z-broker-acl-coldstart-introspection-miss-done.md
20260904T195529Z-issue-93-done.md
20260904T195814Z-broker-disconnect-device-zombie-session-blocked.md
20260904T203003Z-broker-offline-acl-tests-blocked.md
20260904T210532Z-broker-onconnect-doc-comment-blocked.md
20260905T110211Z-empty-introspection-verification-blocked.md
20260905T111428Z-probe-trigger-install-notification-delay-blocked.md
20260905T112539Z-compat-note-v132-done.md
20260905T174528Z-swagger-httptest-coverage-done.md
===verdicts===
     13 verdict: blocked
     25 verdict: done
     30 verdict: proposed
      1 verdict: timeout

✱ Grep "verdict" in tools · 11 matches
The recipe already ran today and proposed; this run found no upstream movement since. Let me check how the runner parses verdicts before writing the evidence file.
→ Read tools/mule.sh [offset=260, limit=80]

← Write .mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md
Wrote file successfully.

Nothing upstream has moved: the tag list is byte-identical to the run earlier today (`1e56fda`), which already proposed the only new wire-visible item (`probe-required-mapping-flag`, `compat-note-v1.4-rc` — both now in the queue as `- [x]`/`- [!]`). Re-proposing either would duplicate. Per the recipe, an empty result is the correct result here.

Writing the evidence file.
No upstream movement: the newest tag is still v1.4.0-rc.5 (2026-08-20), unchanged since the recipe's earlier run today (`1e56fda`), which already proposed the only qualified items (`probe-required-mapping-flag`, `compat-note-v1.4-rc`). Both are already in the queue, so there was nothing new to propose and I appended nothing to `.mule/todo.md` — the empty result the recipe treats as the good one.

Done: upstream watch — no upstream movement since v1.4.0-rc.5
Files: .mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md
Verified: `gh api repositories releases` -> newest v1.4.0-rc.5 (unchanged since 2026-08-20); git status clean
Unsure: nothing
Follow-ups: none
