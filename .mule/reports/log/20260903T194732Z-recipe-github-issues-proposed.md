slug: recipe-github-issues
verdict: proposed
at:  477cc01
ran: 2026-09-03T19:47:32Z on DietPi in 138s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md
Running the issue survey now.

$ gh issue list --state open --limit 40 --json number,title,labels,updatedAt,comments --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
#91 Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health [mule-blocked upstream-parity ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:Scoped for the mule (2026-08-31).

**In scope:** serve the unauthenticated pairing health check at \`GET /pairing/health\` (no realm segment), matching the path shape measured upstream, and **keep** the existing \`GET /pairing/v1/{realm}/health\` route as-is. Same handler, same 200 payload; additive only, so nothing that already polls the v1 route breaks. Tests for both paths.

**Out of scope:** removing or moving the v1/{realm} route, any docs file, and re-running the original realmcfg-02 probe against upstream — that needs the Legion Go, which is off by default; it stays parked and is not a precondition for the additive route above. createdAt:2026-08-30T23:33:37Z id:IC_kwDORmfsJs8AAAABRidaKA includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/91#issuecomment-5471951400 viewerDidAuthor:true] map[author:map[login:atsetilam] authorAssociation:MEMBER body:The mule could not do this: **gates failed**

Taken out of its queue so it does not retry the same failure every half hour. Re-label `mule` to queue it again, ideally after making the task smaller or the requirement clearer. createdAt:2026-08-31T11:31:15Z id:IC_kwDORmfsJs8AAAABRn6yRg includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/91#issuecomment-5477675590 viewerDidAuthor:true]]c
#87 Flow block: lua_map — needs embedded Lua runtime (parked) [enhancement upstream-parity ] []c
#78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate) [enhancement milestone-4.0 upstream-parity ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:Reframed 2026-08-23 by Giulio's decision: no longer parked — **strategic feature** (zero-touch factory onboarding is a commercial-viability requirement for an IoT platform), milestone-4.0 candidate.

## Verified facts (2026-08-23, upstream master)

- Upstream integrates FDO **inside the Pairing service**: `fdo_onboarding_controller.ex`, session plugs (`setup_fdo`, `fdo_session`), CBOR codec + voucher queries in `astarte_data_access`, and two dedicated libs `astarte_fdo` / `astarte_fdo_core` (SECO Mind, 2025; bugfixes through 2026-05). Not an optional sidecar.
- Device side is the official SDK `astarte-device-fdo-rust` (crate on docs.rs, pushed 2026-08-21) — actively maintained.
- Upstream does NOT implement manufacturing / rendezvous services either; those come from the FIDO Alliance reference stack (`fdo-rs/fido-device-onboard-rs`, used by Red Hat/Fedora IoT).

## Scope decision

Astrate implements **only the last mile**, mirroring upstream's cut:
- owner-side protocol surface in our pairing service (TO1 redirect + TO2 onboarding, CBOR wire format, ownership-voucher storage, owner key management)
- reuse the existing open-source FDO ecosystem for manufacturing/rendezvous — never reimplement it.

**Acceptance:** a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end and lands as a provisioned realm device.

**Docs are a first-class deliverable:** upstream ships almost no public FDO documentation (code-level only); ours must exceed that — operator guide covering the full chain (manufacturing → rendezvous → owner) with Astrate-specific setup.

## Next step

Investigation phase before any implementation: read upstream's TO2 handling end-to-end (`libs/astarte_fdo*", pairing controllers/plugs, `fdo/queries`), inventory the exact endpoints, credential/voucher schema and key material required, and measure what our pairing service lacks. Output feeds the v4.0 scope decision in `.mule/milestones.md`. createdAt:2026-08-23T16:16:08Z id:IC_kwDORmfsJs8AAAABQRcugw includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/78#issuecomment-5387005571 viewerDidAuthor:true]]c
#68 Decide async_operation=false params vs documented always-sync deviation [enhancement mule-blocked upstream-parity upstream-experimental ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:Scoped for the mule (2026-08-31).

Design is already frozen in \`.trickle/plans/MASTER-HANDOFF.md\` fase 4c and is not up for relitigation: Astrate stays always-sync, and \`async_operation=false\` is accepted and ignored on the mutating endpoints upstream exposes it on — housekeeping realm create/delete; realm-management interface install/update/delete and trigger/policy delete (policy routes are mounted in \`internal/realm/http.go\`). An unparseable or \`true\` value must not change behaviour either.

**In scope for the mule:** the code and its tests only.
**Out of scope:** \`docs/COMPATIBILITY.md\` (deviation 17), \`docs/UPSTREAM-EXPERIMENTAL.md\`, and closing this issue — the docs closeout stays with the architect and lands once, together with #67. createdAt:2026-08-30T23:33:27Z id:IC_kwDORmfsJs8AAAABRidW1w includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/68#issuecomment-5471950551 viewerDidAuthor:true] map[author:map[login:atsetilam] authorAssociation:MEMBER body:The mule could not do this: **gates failed**

Taken out of its queue so it does not retry the same failure every half hour. Re-label `mule` to queue it again, ideally after making the task smaller or the requirement clearer. createdAt:2026-08-31T11:47:36Z id:IC_kwDORmfsJs8AAAABRoFeQg includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/68#issuecomment-5477850690 viewerDidAuthor:true]]c
#67 Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4) [enhancement upstream-parity upstream-experimental ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:Not delegated to the mule (assessed 2026-08-31).

Both code phases have landed (4a schema layer \`9c4d411\`, 4b runtime enforcement of \`required\` \`7a244b2\`). What is left is the docs closeout — the \`docs/COMPATIBILITY.md\` row explaining why astrate’s \`missing_required\` reason is translated to \`unexpected_object_key\`, and the \`docs/UPSTREAM-EXPERIMENTAL.md\` Adopted dates — which is prose the mule does not write, and it lands together with #68's closeout in fase 4c. Keeping this issue out of the mule queue on purpose.

\`encrypted\` runtime behaviour remains tied to keyAgreement (#51) and is out of this issue’s scope. createdAt:2026-08-30T23:33:48Z id:IC_kwDORmfsJs8AAAABRiddQA includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/67#issuecomment-5471952192 viewerDidAuthor:true]]c
#1 Provide an Open Source IoT Platform unironically, unironically [wontfix ] []c

Let me check the existing todo.md and for-giulio.md, and look at the two mule-blocked issues more closely.
→ Read .mule/todo.md

$ gh issue view 91 --json number,title,body,labels,state --jq '{number,title,state,labels:[.labels[].name],body_preview: .body[:800]}'
{"body_preview":"Found during the 2026-08-24 verify batch (#57/#69/#77/#79 probes): upstream Astarte 1.2.0 answers\n\n- `GET /pairing/health` → 200\n- `GET /pairing/v1/{realm}/health` → 404 Page not found (bench and scratch realms both)\n\nAstrate implements the unauthenticated health route at `GET /pairing/v1/{realm}/health` (landed in realmcfg-02, issue #71, commit 5111e3c) — the path shape does not exist upstream. The dashboard polls per-service health at the service root, so the current route may be invisible to it.\n\nDecide: move/duplicate the handler at `GET /pairing/health` (keeping or dropping the v1/{realm} variant), and re-check whatever probe led realmcfg-02 to the v1 shape — that session's stack had a port override in play, so the original measurement may have hit a different router.","labels":["mule-blocked","upstream-parity"],"number":91,"state":"OPEN","title":"Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health"}

$ gh issue view 68 --json number,title,body,labels,state --jq '{number,title,state,labels:[.labels[].name],body_preview: .body[:800]}'
{"body_preview":"Upstream defaults to async execution with `?async_operation=false` opt-in for: realm create/delete (housekeeping), interface install/update/delete (RM), policy delete (master). Astrate is always-sync everywhere — a superset clients tolerate, but a behavioral difference worth either matching or recording as deviation #11.\nSize S decision.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":["enhancement","mule-blocked","upstream-parity","upstream-experimental"],"number":68,"state":"OPEN","title":"Decide async_operation=false params vs documented always-sync deviation"}

$ ls -la /root/astrate-mule/.mule/for-giulio.md 2>/dev/null; echo "---"; cat /root/astrate-mule/.mule/for-giulio.md 2>/dev/null || echo "(does not exist)"
-rw-r--r-- 1 root root 8721 Sep  3 21:28 /root/astrate-mule/.mule/for-giulio.md
---
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **milestone 2.0 looks complete, verify and cut the tag** — all 11 `milestone-2.0` issues
  CLOSED (#23–#27, #37, #39–#43), no open issues, no new gaps after re-checking upstream
  astarte_flow block catalog against `internal/flow/` + git log (MQTT/HTTP source/sink,
  json_path_map, pure-transform set, virtual_device_pool, container block MVP, flow API,
  durable named flows all landed). `mule/queue` branch is stale — see entry below.
  (Milestones recipe run, 2026-09-03.)

---

- **v3.0 decide #67** (`upstream-parity`): `required` and `encrypted` mapping fields exist on upstream master but not at 1.2.2/1.3.0 (`required` = missing key at runtime → `missing_required_mapping` on object-aggregated mappings; `encrypted` ties into keyAgreement #51). Decide now: reject (strict 1.2.2), ignore-and-accept, or implement `required` early. Issue is OPEN, no decision on record, not yet escalated here.

---

- **v3.0 decide #91** (`upstream-parity`, `mule-blocked`): upstream serves the unauthenticated health check at `GET /pairing/health` and returns 404 for `GET /pairing/v1/{realm}/health`, but Astrate only implemented the `v1/{realm}` shape (issue #71, commit 5111e3c) — the dashboard polls per-service health at the service root, so the current route may be invisible to it. Decide: move/duplicate the handler at `GET /pairing/health` (keep or drop the v1/{realm} variant), and re-check the original realmcfg-02 measurement (a port override may have skewed it).

---

- **`govulncheck` is not installed on the Pi**, so the hygiene recipe's highest-priority check (reachable dependency vulnerabilities) could not run here. `go install golang.org/x/vuln/cmd/govulncheck@latest` would close it — I did not pick a version because the pin is a decision. The `go list -m -u` sweep shows only version-skew on transitive deps (cloud/azure SDKs, cel/expr, clickhouse); none corresponds to an identified fix this repo needs, so I proposed no bumps. (Hygiene recipe run, 2026-09-02.)

---

- **`docs/site/appengine-api.md:87-88` documents `GET` and `DELETE /appengine/v1/<realm>/groups/<name>`**, but no such routes exist in `internal/appengine/http.go` (only `/groups/{group}/devices` and the nested device ops — there is no `GET` or `DELETE` on `/groups/{group}`). The endpoint the page describes is absent from the code. Decide: drop the two lines, mark them not-yet-implemented, or have the code grow them. (Docs-sync recipe run, 2026-09-02.)

---

- **`Router.Submit` TOCTOU on `closed` flag** (`internal/flow/router.go:113-120`): Submit
  reads `r.closed` under the mutex, drops the lock, then sends on the channel. A concurrent
  `Drain` could close the channel between the unlock and the send, causing a send-on-closed
  channel panic in the caller (not recovered by `processOne`'s defer). Fix options: (a) hold
  the lock during channel send (hot-path perf hit), (b) use a select with `r.quit` instead of
  sending on a closed channel, (c) accept the risk and document that Submit must not be called
  after Drain. Needs a design call — not proposing a task. (Code review 2026-09-02.)

---

- **#87 `lua_map` — needs embedded Lua runtime, parked.** Design/implementation decision: embedding a Lua VM in Astrate is not machine-checkable by the mule. Consider closing if Lua flow support is not on any active roadmap.
- **#78 FDO device onboarding — milestone-4.0, investigation phase.** Too large for a single mule task; the investigation work (reading upstream's TO2 handling, inventorying endpoints/schema/keys) is a multi-session project. Parking for now until the v3.0 queue clears and this becomes the next milestone target.
- **#1 stale — "Provide an Open Source IoT Platform unironically" (wontfix).** Has been open since the repo's founding with no activity. Consider closing.

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

- ~~The Pi cannot run the race detector~~ — **resolved 2026-07-27** by installing Go 1.26.5
  as a userland toolchain on the Legion Go (`~/.local/go`, no root, `rm -rf` to undo). The
  Pi still cannot run `-race` (39-bit VMA kernel vs the 48 ThreadSanitizer needs), so its
  gate remains `go vet ./... && go test ./...` — but race coverage now exists on the Legion
  Go, where the full suite runs clean in ~40s on 16 cores. The standing `race-check` task is
  the concurrency gate. Concurrency work is queueable again, provided the race-check runs
  after it.
- **golangci-lint is not installed on the Pi**, so the mule's second gate is silently absent
  there — `gofmt` still runs, the linter does not. `go install
  github.com/golangci/golangci-lint/cmd/golangci-lint@<the pinned version>` on the Pi would
  close it; I did not pick a version for you, since the pin is a decision.
- ~~`/root/astrate` on the Pi has uncommitted work~~ — **resolved 2026-07-27** with the new
  `tools/reconcile.sh`: rescued onto `origin/wip/DietPi-20260727T171543Z` (pushed, not
  reviewed — read the diff before merging anything from it) and `/root/astrate` is now clean
  on `main`. Also had to set `commit.gpgsign false` locally on that clone first — it had
  signing on with no working gpg-agent for a non-interactive session, same class of problem
  `/root/astrate-mule` already solved the same way (your call, confirmed 2026-07-27: align it
  with the existing unattended-Pi-clone convention rather than fix gpg-agent or bypass signing
  silently inside the script).
- **The mule has been idle 15h.** Filed by the dead-man's switch; see journalctl on the Pi.

Now let me verify which issues are already in todo.md or for-giulio.md to avoid duplicates.

The survey shows 5 open issues:

| # | Title | Labels | In todo.md | In for-giulio.md |
|---|-------|--------|-----------|-----------------|
| 91 | Pairing health endpoint divergence | mule-blocked | `issue-91-pairing-health-serve-root` (blocked) | Yes |
| 87 | Flow block: lua_map (parked) | enhancement | No | Yes |
| 78 | FDO device onboarding | milestone-4.0 | No | Yes |
| 68 | async_operation=false decision | mule-blocked | `issue-68-async-operation-accepted` (blocked) | **No** |
| 67 | Interfaces required/encrypted | upstream-parity | No | Yes |
| 1 | wontfix | wontfix | No | Yes |

**Analysis per triage rules:**

- **#91** — machine-checkable, scoped for the mule, but already queued (blocked: gates failed). Cannot re-propose.
- **#68** — machine-checkable, scoped for the mule, but already queued (blocked: gates failed). Cannot re-propose. Missing from for-giulio.md though.
- **#87** — needs embedded Lua runtime; not machine-checkable by the mule. Already in for-giulio.md.
- **#78** — strategic milestone-4.0 feature, investigation phase. Already in for-giulio.md.
- **#67** — "not delegated to the mule" per its own comment. Already in for-giulio.md.
- **#1** — wontfix, no activity. Already in for-giulio.md.

No new task lines to append to todo.md — every candidate is either already queued or not machine-checkable.

One gap: #68 is not in for-giulio.md yet. I'll add it.
→ Read .mule/for-giulio.md [limit=25]

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -16,8 +16,10 @@
   json_path_map, pure-transform set, virtual_device_pool, container block MVP, flow API,
   durable named flows all landed). `mule/queue` branch is stale — see entry below.
   (Milestones recipe run, 2026-09-03.)
 
+- **v3.0 decide #68** (`upstream-parity`, `mule-blocked`): `async_operation=false` acceptance is already decided (Astrate stays always-sync; accept and ignore on mutating endpoints: housekeeping realm create/delete, realm-management interface install/update/delete, trigger/policy delete in `internal/realm/http.go`). The mule's code+test task (`issue-68-async-operation-accepted`) failed gates and is parked. Needs re-queue or gate fix.
+
 ---
 
 - **v3.0 decide #67** (`upstream-parity`): `required` and `encrypted` mapping fields exist on upstream master but not at 1.2.2/1.3.0 (`required` = missing key at runtime → `missing_required_mapping` on object-aggregated mappings; `encrypted` ties into keyAgreement #51). Decide now: reject (strict 1.2.2), ignore-and-accept, or implement `required` early. Issue is OPEN, no decision on record, not yet escalated here.
 

→ Read .mule/for-giulio.md [offset=20, limit=15]
Done: GitHub issues triage — no new task lines (all open issues already queued, blocked, or not machine-checkable)
Files: .mule/for-giulio.md (added #68 escalation)
Verified: 5 open issues surveyed; #91/#68 already in todo.md blocked; #87/#78/#67/#1 already in for-giulio.md; added missing #68
Unsure: nothing
Follow-ups: none
