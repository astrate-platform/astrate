slug: recipe-github-issues
verdict: proposed
at:  c7deb30
ran: 2026-09-03T10:55:57Z on DietPi in 103s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md

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

Let me check what's already queued and read the bodies of the most actionable issues.
→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ gh issue view 91 --json body,comments --jq '.body, "\n---LAST COMMENT---\n", (.comments[-1].body // "none")'
Found during the 2026-08-24 verify batch (#57/#69/#77/#79 probes): upstream Astarte 1.2.0 answers

- `GET /pairing/health` → 200
- `GET /pairing/v1/{realm}/health` → 404 Page not found (bench and scratch realms both)

Astrate implements the unauthenticated health route at `GET /pairing/v1/{realm}/health` (landed in realmcfg-02, issue #71, commit 5111e3c) — the path shape does not exist upstream. The dashboard polls per-service health at the service root, so the current route may be invisible to it.

Decide: move/duplicate the handler at `GET /pairing/health` (keeping or dropping the v1/{realm} variant), and re-check whatever probe led realmcfg-02 to the v1 shape — that session's stack had a port override in play, so the original measurement may have hit a different router.

---LAST COMMENT---

The mule could not do this: **gates failed**

Taken out of its queue so it does not retry the same failure every half hour. Re-label `mule` to queue it again, ideally after making the task smaller or the requirement clearer.

$ gh issue view 68 --json body,comments --jq '.body, "\n---LAST COMMENT---\n", (.comments[-1].body // "none")'
Upstream defaults to async execution with `?async_operation=false` opt-in for: realm create/delete (housekeeping), interface install/update/delete (RM), policy delete (master). Astrate is always-sync everywhere — a superset clients tolerate, but a behavioral difference worth either matching or recording as deviation #11.
Size S decision.\n\n---\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).

---LAST COMMENT---

The mule could not do this: **gates failed**

Taken out of its queue so it does not retry the same failure every half hour. Re-label `mule` to queue it again, ideally after making the task smaller or the requirement clearer.

$ gh issue list --state open --limit 40 --label mule --json number,title,labels,comments --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
(no output)

All open issues are already triaged:

| Issue | Status |
|-------|--------|
| #91 pairing-health | Already in `todo.md` (blocked) + `for-giulio.md` |
| #68 async_operation | Already in `todo.md` (blocked) |
| #87 lua_map | Already in `for-giulio.md` (parked, design decision) |
| #78 FDO onboarding | Already in `for-giulio.md` (milestone-4.0, too large) |
| #1 stale wontfix | Already in `for-giulio.md` (suggested close) |

No issues are labelled `mule` currently. Nothing new to propose — all actionable issues are already queued (blocked) or escalated. No lines appended to `todo.md` or `for-giulio.md`.
