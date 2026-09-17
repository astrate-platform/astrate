slug: recipe-astarte-upstream
verdict: proposed
at:  5a314c7
ran: 2026-09-17T19:07:03Z on DietPi in 106s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

$ git log --oneline -10
5a314c7 mule: refill the queue
b3e8361 mule: recipe github-issues ran (126s)
9d75bee mule: recipe milestones timed out
d9a534b mule: recipe hygiene timed out
e33ab46 mule: log docs-sync-native-version-value
08048b3 mule: docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
30665a5 mule: log docs-sync-pairing-version-value
c69a164 mule: docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
4c12a38 mule: refill the queue
ee22e79 mule: recipe docs-sync ran (588s)

→ Read .mule/todo.md
→ Read .mule/todo.md [offset=173]
→ Read .mule/for-giulio.md
→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ git log --all --oneline --grep='upstream' | head -30
c69a164 mule: docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
f27dda1 mule: recipe astarte-upstream ran (121s)
9471dab mule: recipe astarte-upstream ran (370s)
203af1c mule: hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto]
015a706 mule: hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto]
f0863f2 mule: recipe astarte-upstream ran (97s)
6524c7f mule: auth-iat-not-required-test [auto]: pin in internal/auth/jwt_test.go that a token with a *future* `iat` still verifies — the "`iat` is not required — upstream parity" rule (jwt.go:104-105, no `jwt.WithIssuedAt` in the parser, jwt.go:115-118) has zero coverage, so a regression that starts validating `iat` or `WithIssuedAt` would silently reject tokens from a skewed-clock issuer and pass every existing suite. Sign RS256 with `iat` = now+1h and no `exp`/`nbf`, assert `Verify` succeeds.
6bca5cc mule: recipe astarte-upstream ran (100s)
d29e382 mule: recipe astarte-upstream ran (119s)
a609274 mule: recipe astarte-upstream ran (334s)
3cbd3d6 mule: recipe astarte-upstream ran (375s)
ee76ba1 mule: recipe astarte-upstream ran (204s)
321fc79 mule: recipe astarte-upstream ran (152s)
c867b2c mule: docs-sync-pairing-version-endpoint [auto]: add the undocumented `GET /pairing/v1/{realm}/version` route to docs/api/astarte_pairing_api.yaml — it exists in code (cmd/astrate/main.go:398, unauthenticated, upstream-parity per main.go:387-391) and answers 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38); mirror the RM spec's getVersion op (astarte_realm_management_api.yaml:532) without its a_rma security. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
2f6176c mule: docs-sync-pairing-status-enum [auto]: fix the `PairingInfo.status` enum in docs/api/astarte_pairing_api.yaml from `[confirmed, pending, denied, expired]` (yaml:412) to `[confirmed, pending, inhibited]` — the handler never emits `denied` or `expired`: `service.Info` maps only `DeviceStatusInhibited`→"inhibited", `DeviceStatusConfirmed`→"confirmed", and everything else→"pending" (internal/pairing/service.go:296-299, statuses defined in store/devices.go:19-27; upstream-parity per the comment at service.go:282-285). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
f463039 mule: recipe astarte-upstream ran (268s)
a701b05 mule: recipe astarte-upstream ran (239s)
32544c0 mule: pairing-bearer-secret-test [auto]: add a container-free unit test for `bearerSecret` (internal/pairing/http.go:343) locking the upstream-verbatim permissive match `~r/bearer\:?\s+(.*)$/i` — the colon form `Bearer: <secret>`, the case variants (`bearer `, `BEARER `), and the reject paths (no whitespace before value, non-bearer scheme, empty value). Only the plain `Bearer <secret>` form is exercised today (integration http_test.go); a regression to case-sensitive or drop-colon matching would pass every existing test. Pure helper in a non-integration file; no Docker.
95ef44e mule: recipe astarte-upstream ran (237s)
db9d06e mule: probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
f117601 mule: recipe astarte-upstream ran (305s)
ab49dcc mule: recipe astarte-upstream ran (127s)
c3cad68 mule: recipe astarte-upstream ran (85s)
4637135 mule: recipe astarte-upstream ran (526s)
4071671 mule: recipe astarte-upstream ran (730s)
4bc3e1a mule: compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
24ad5b8 mule: issue-93: In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
2a033d4 mule: parking a decision on an upstream release now has something that notices
99743ca mule: the keyAgreement parking condition has fired — record the scope decision
f1d0069 mule: the survey stops stranding reports, and the map stops lying about timers

$ git log --format='%h %ad %s' -1 f27dda1
f27dda1 Wed Sep 16 20:42:18 2026 +0200 mule: recipe astarte-upstream ran (121s)

$ git show f27dda1 --stat
commit f27dda1adfdb7aa95a6d1d6975961dbc49c92a5f
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Wed Sep 16 20:42:18 2026 +0200

    mule: recipe astarte-upstream ran (121s)

 ...0260916T184159Z-recipe-astarte-upstream-noop.md |  45 ++++++++
 ...916T184218Z-recipe-astarte-upstream-proposed.md | 126 +++++++++++++++++++++
 2 files changed, 171 insertions(+)

$ git show f27dda1
commit f27dda1adfdb7aa95a6d1d6975961dbc49c92a5f
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Wed Sep 16 20:42:18 2026 +0200

    mule: recipe astarte-upstream ran (121s)

diff --git a/.mule/reports/log/20260916T184159Z-recipe-astarte-upstream-noop.md b/.mule/reports/log/20260916T184159Z-recipe-astarte-upstream-noop.md
new file mode 100644
index 0000000..ec6e9d3
--- /dev/null
+++ b/.mule/reports/log/20260916T184159Z-recipe-astarte-upstream-noop.md
@@ -0,0 +1,45 @@
+slug: recipe-astarte-upstream
+verdict: noop
+at:  ed715ec
+ran: 2026-09-16T18:41:59Z on linux x86_64
+
+
+> build · big-pickle
+
+→ Read .mule/MULE.md
+→ Read .mule/recipes/astarte-upstream.md
+→ Read docs/COMPATIBILITY.md
+→ Read .mule/todo.md
+
+$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
+v1.4.0-rc.5	2026-08-20T12:28:04Z
+v1.3.3	2026-08-07T14:33:58Z
+v1.4.0-rc.4	2026-08-06T13:40:20Z
+v1.4.0-rc.3	2026-07-31T10:48:39Z
+v1.4.0-rc.2	2026-07-14T10:16:22Z
+
+The release list is unchanged since the last astarte-upstream recipe run (2026-09-08, report
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md): no new stable, no new
+rc, still `v1.4.0-rc.5` (2026-08-20) as the newest tag, `v1.3.3` (2026-08-07) as the newest
+stable. No movement on any repo in this recipe's scope since then.
+
+Astrate targets upstream 1.2.2 (docs/COMPATIBILITY.md, `APICompatVersion` 1.2.2, #90 frozen
+decision); everything upstream has shipped since that target is already accounted for in the
+queue from prior runs:
+
+- v1.3.x: purge-properties-compression-capability [x], empty-introspection-verification [!],
+  probe-trigger-install-notification-delay [!], probe-property-resend-encoding [!],
+  probe-props-resend-error-triggers [!], compat-note-v1.3.2 [x], compat-note-v1.3.3 [!]
+- v1.4.0-rc line: probe-required-mapping-flag [x], compat-note-v1.4-rc [!]
+
+No new proposal is qualified: nothing to add, and re-proposing the compat-note or a probe
+line already in the queue would duplicate. Recipe rule honoured: "If nothing changed since
+last time, append nothing to the queue and say 'no upstream movement since <tag>'."
+
+No upstream movement since v1.4.0-rc.5.
+
+Done: upstream watch — no upstream movement since v1.4.0-rc.5; appended nothing to .mule/todo.md
+Files: .mule/reports/log/20260916T184159Z-recipe-astarte-upstream-noop.md (evidence; todo.md untouched)
+Verified: gh api repos/astarte-platform/astarte/releases -> v1.4.0-rc.5 newest, unchanged since 2026-08-20
+Unsure: nothing
+Follow-ups: none
\ No newline at end of file
diff --git a/.mule/reports/log/20260916T184218Z-recipe-astarte-upstream-proposed.md b/.mule/reports/log/20260916T184218Z-recipe-astarte-upstream-proposed.md
new file mode 100644
index 0000000..6a08f85
--- /dev/null
+++ b/.mule/reports/log/20260916T184218Z-recipe-astarte-upstream-proposed.md
@@ -0,0 +1,126 @@
+slug: recipe-astarte-upstream
+verdict: proposed
+at:  ed715ec
+ran: 2026-09-16T18:42:18Z on DietPi in 121s
+
+
+> build · big-pickle
+
+→ Read .mule/MULE.md
+→ Read .mule/recipes/astarte-upstream.md
+→ Read docs/COMPATIBILITY.md
+Recipe read. First: read `docs/COMPATIBILITY.md` and check the latest upstream releases in parallel.
+
+$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
+v1.4.0-rc.5	2026-08-20T12:28:04Z
+v1.3.3	2026-08-07T14:33:58Z
+v1.4.0-rc.4	2026-08-06T13:40:20Z
+v1.4.0-rc.3	2026-07-31T10:48:39Z
+v1.4.0-rc.2	2026-07-14T10:16:22Z
+
+→ Read .mule/todo.md
+Now let me check the current todo and prior survey context to see what version is being tracked.
+
+$ rg -n "1\.4|1\.3|target|COMPATIBILITY" .mule/for-giulio.md .mule/milestones.md docs/UPSTREAM-EXPERIMENTAL.md 2>/dev/null | head -40
+docs/UPSTREAM-EXPERIMENTAL.md:6:it. **Review trigger: every time a future parity milestone targets a newer
+docs/UPSTREAM-EXPERIMENTAL.md:24:| Required + encrypted interface mapping fields | #67 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
+docs/UPSTREAM-EXPERIMENTAL.md:25:| `async_operation=false` on mutating RM endpoints | #68 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
+.mule/for-giulio.md:13:- **github-issues triage run, 2026-09-16: still nothing proposable — the mule-alarm pile is now 11 straight days (#94–#104).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), plus the daily alarm. New since the 2026-09-15 run: **#104**, today's alarm (created 09-16 11:12Z, "nothing has landed in 14h"); **#103** has expired into the pile (superseded by #104). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-15 one: close #94–#103** — each was a one-day event, superseded by next day's alarm, never actionable; leave #104 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
+.mule/for-giulio.md:15:- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-13 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+.mule/for-giulio.md:17:- **github-issues triage run, 2026-09-15: still nothing proposable — the mule-alarm pile is now 10 straight days (#94–#103).** Fourteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-13 run: **#103**, the alarm created 2026-09-14 11:27Z ("nothing has landed in 16h"); **#102** has expired into the pile (superseded by #103). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-13 one: close #94–#102** — each was a one-day event, superseded by next day's alarm, never actionable; leave #103 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
+.mule/for-giulio.md:19:- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-12 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+.mule/for-giulio.md:21:- **github-issues triage run, 2026-09-13: still nothing proposable — the mule-alarm pile is now 9 straight days (#94–#102).** Thirteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-12 run: **#102**, today's alarm (created 11:09Z, "nothing has landed in 14h"); **#101** has expired into the pile (superseded by #102). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-12 one: close #94–#101** — each was a one-day event, superseded by next day's alarm, never actionable; leave #102 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
+.mule/for-giulio.md:23:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+.mule/for-giulio.md:27:- **github-issues triage run, 2026-09-12: still nothing proposable — the mule-alarm pile is now 8 straight days (#94–#101).** Twelve open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-11 run: **#101**, today's alarm (created 10:58Z, "nothing has landed in 14h"); **#100** has expired into the pile (superseded by #101). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-11 one: close #94–#100** — each was a one-day event, superseded by next day's alarm, never actionable; leave #101 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
+.mule/for-giulio.md:29:- **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pile is now 7 straight days and worth a look.** Twelve open issues. Still no machine-checkable fix candidates, so no task lines: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-10 run: **#100**, today's alarm (created 10:55Z, "nothing has landed in 15h"). #99–#94 are the same `mule-alarm` (one per day 2026-09-05→09-10, ~11:00Z each, zero comments) — telemetry, not code issues, so never proposable. But seven consecutive alarms is past the "low-activity window" wording earlier runs used: the mule has landed no commit since ~2026-09-04/05, and `main`'s queue copy shows the top line as `- [ ]` while the landed work below is all `[!] BLOCKED` (`wrote nothing` / `tests failed`) — though `mule/queue` is the authoritative copy. **Proposal: close #94–#99** — each was a one-day event, superseded by next day's alarm, never actionable; leave #100 (live today) to self-expire. If the idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
+.mule/for-giulio.md:31:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+.mule/for-giulio.md:33:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+.mule/for-giulio.md:35:- **Extension point: external-bus Intake (NATS JetStream), 2026-09-10.** `internal/broker/intake.go:56-62` carries a `TODO(extension point)` for a second `Intake` implementation backed by a durable external bus (NATS JetStream) for multi-instance deployment or restart survival. The frozen design decision is in `docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4`; the consumer must reproduce per-device ordering, deferred-ack backpressure, and QoS 0 drop semantics. Not tracked by any open issue. Your call: file a milestone issue or park until the roadmap calls for it.
+.mule/for-giulio.md:59:  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
+.mule/for-giulio.md:65:  issues are the same: #92 keyAgreement (upstream decision, gated on stable v1.4.0), #93
+.mule/for-giulio.md:68:  on upstream v1.4 final (currently rc.5). `APICompatVersion` stays 1.2.2 per #90's frozen
+.mule/for-giulio.md:76:  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
+.mule/for-giulio.md:84:- **Milestone recipe run, 2026-09-06: v3.0 (upstream 1.2.2 → 1.3/1.4) has no open
+.mule/for-giulio.md:87:  the 1.3-surface gaps from the #47–#89 backlog are delivered/closed — retention ceiling #72,
+.mule/for-giulio.md:90:  remaining open items are already escalated, not fileable gaps: **#92 keyAgreement** (1.4
+.mule/for-giulio.md:91:  experimental decision, escalated 2026-09-05, gated on a stable upstream v1.4.0), **#78 FDO**
+.mule/for-giulio.md:93:  `8c61268`). `docs/UPSTREAM-EXPERIMENTAL.md` rows #67/#68 are "1.4 experimental" and reconcile
+.mule/for-giulio.md:94:  only when upstream ships **v1.4 final** — not yet (still rc). Per **#90's frozen decision**
+.mule/for-giulio.md:95:  (2026-08-23), `APICompatVersion` stays 1.2.2 until the FULL 1.3 surface is complete *and*
+.mule/for-giulio.md:97:  call when v1.4.0 goes stable: answer #92 (implement against the rc, re-park, or narrow #93
+.mule/for-giulio.md:104:  unaffected, but if that doc is your canonical 1.3/1.4 diff it needs restoring (or the
+.mule/for-giulio.md:117:  awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the waiting-on
+.mule/for-giulio.md:125:  stable v1.4.0?** (issue #92, `upstream-parity`/`upstream-experimental`). #51 closed
+.mule/for-giulio.md:126:  2026-08-22 with "parked until the upstream 1.4 experimental spec stabilizes — reopen or
+.mule/for-giulio.md:133:  fixed the stale ACL comment that cited them). But the spec ships only in `v1.4.0-rc.5` —
+.mule/for-giulio.md:134:  `v1.3.3` is still the newest stable tag and Astrate targets 1.2.2 — so the document has
+.mule/for-giulio.md:138:  re-park until v1.4.0 is a stable tag, or take only the narrow #93 fix. (Escalated again
+.mule/for-giulio.md:148:- **COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable, 2026-07-14; v1.4.0 is still rc-only).** Astrate's doc and `APICompatVersion` still target upstream **1.2.2** (`internal/realm/service.go:588`); v1.3.0 (2026-05-06) introduced wire-surface changes Astrate does not yet emulate, so this is a decision — adopt v1.3.2 as the compatibility target (then update the doc + bump `APICompatVersion` together, per the bump rule) or keep 1.2.2 and add a "not yet emulated" note. Wire-relevant v1.3.0 deltas (release notes): **MQTT v1 capabilities** incl. `purge_properties_compression_format` (plaintext vs zlib purge — touches the `emptyCache`/`producer,properties` contract COMPATIBILITY.md deviation 1 documents); **empty introspection now allowed**; **device registration triggers** (pairing) and **device deletion started/completed triggers** (RM — the latter two already exist as Astrate deviation 9 emits both around the synchronous delete); **FDO authentication** (pairing, disabled by default); **realm-scoped health** — upstream v1.3 added `GET /pairing/v1/{realm}/health`, which Astrate already serves (`internal/pairing/http.go:78-81`, comment already says "upstream 1.3+"), so deviation 18's wording ("which upstream 404s") is now false against 1.3 and the note should be reworded either way. Proposed doc wording (for your approval, edit to taste): in §Infrastructure differences add a sentence — *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl.
+.mule/for-giulio.md:150:plaintext `purge_properties_compression_format`, empty-introspection allowance, device registration/deletion triggers, experimental FDO pairing auth) are not yet emulated and are out of scope until the milestone that adopts v1.3.2 as the target."* — and reword deviation 18's realm-health note from "which upstream 404s" to "added by upstream v1.3 (Astrate serves it against a 1.2.2 target; kept, matching behavior)". Raw upstream changes: [v1.3.0](https://github.com/astarte-platform/astarte/releases/tag/v1.3.0), [v1.3.2](https://github.com/astarte-platform/astarte/releases/tag/v1.3.2).
+.mule/for-giulio.md:170:  - `go.etcd.io/bbolt` v1.4.3 → v1.5.0 — bbolt's semver promises no API change between patch/minor, so additive-only; used in `internal/broker/sessionstore.go`; v1.5 adds a data-file size limit and panic-recovery hardening, nothing Astrate needs.
+.mule/for-giulio.md:216:  becomes the next milestone target.
+.mule/milestones.md:5:target** — the first section below not marked `DONE` — and works toward it: investigating,
+.mule/milestones.md:10:`.mule/for-giulio.md`, but never edits it directly — same rule as `docs/COMPATIBILITY.md`.
+.mule/milestones.md:107:## v3.0 — upstream Astarte feature parity (1.2.2 → 1.3/1.4)
+
+Let me check the history of prior astarte-upstream runs to know what was last proposed and when.
+$ rg -n "astarte-upstream|upstream movement|rc\.5|v1\.4\.0-rc|newest tag|newest stable" .mule/ | head -40
+.mule/todo.md:61:`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
+.mule/todo.md:113:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
+.mule/todo.md:130:- [x] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
+.mule/todo.md:131:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
+.mule/recipes/astarte-upstream.md:17:Compare the newest tag to the version named in `docs/COMPATIBILITY.md`. If they match, say
+.mule/recipes/cross-project-survey.md:28:   is a good result; say so and stop, exactly as `.mule/recipes/astarte-upstream.md` already
+.mule/recipes/cross-project-survey.md:33:   upstream clone) to know what to focus on, the same way `.mule/recipes/astarte-upstream.md`
+.mule/recipes/milestones.md:21:     milestone actually requires, same care as `.mule/recipes/astarte-upstream.md` — a
+.mule/milestones.md:28:`.mule/recipes/astarte-upstream.md`'s rule: port the idea, restated in Go, never the code.
+.mule/milestones.md:115:  v1.4.0-rc.5 diffed service-by-service against our emulated 1.2.2 level).
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:1:slug: recipe-astarte-upstream
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:10:→ Read .mule/recipes/astarte-upstream.md
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:16:v1.4.0-rc.5	2026-08-20T12:28:04Z
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:18:v1.4.0-rc.4	2026-08-06T13:40:20Z
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:19:v1.4.0-rc.3	2026-07-31T10:48:39Z
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:20:v1.4.0-rc.2	2026-07-14T10:16:22Z
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:22:The astarte-upstream recipe already ran earlier today (2026-09-08T11:34:22Z, at 1e56fda,
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:23:report .mule/reports/log/20260908T113422Z-recipe-astarte-upstream-proposed.md). The release
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:24:list is unchanged since the 2026-08-20 v1.4.0-rc.5 tag: no new stable, no new rc, no movement
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:33:- v1.4.0-rc line: probe-required-mapping-flag [x], compat-note-v1.4-rc [!]
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:37:last time, append nothing to the queue and say 'no upstream movement since <tag>'."
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:39:No upstream movement since v1.4.0-rc.5.
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:41:Done: upstream watch — no upstream movement since v1.4.0-rc.5; appended nothing to .mule/todo.md
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:42:Files: .mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md (evidence; todo.md untouched)
+.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md:43:Verified: gh api repos/astarte-platform/astarte/releases -> v1.4.0-rc.5 newest, unchanged since 2026-08-20; git status clean
+.mule/for-giulio.md:15:- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-13 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+.mule/for-giulio.md:19:- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-12 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+.mule/for-giulio.md:23:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+.mule/for-giulio.md:31:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+.mule/for-giulio.md:33:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+.mule/for-giulio.md:68:  on upstream v1.4 final (currently rc.5). `APICompatVersion` stays 1.2.2 per #90's frozen
+.mule/for-giulio.md:133:  fixed the stale ACL comment that cited them). But the spec ships only in `v1.4.0-rc.5` —
+.mule/for-giulio.md:134:  `v1.3.3` is still the newest stable tag and Astrate targets 1.2.2 — so the document has
+.mule/reports/log/20260915T170704Z-hk-zero-retention-asymmetry-done.md:46:astarte-upstream.md
+.mule/reports/log/20260915T170704Z-hk-zero-retention-asymmetry-done.md:133:→ Read .mule/recipes/astarte-upstream.md
+.mule/reports/log/20260915T170704Z-hk-zero-retention-asymmetry-done.md:2674:"name": "v1.4.0-rc.5"
+.mule/reports/log/20260915T170704Z-hk-zero-retention-asymmetry-done.md:2675:"name": "v1.4.0-rc.4"
+.mule/reports/log/20260915T170704Z-hk-zero-retention-asymmetry-done.md:2676:"name": "v1.4.0-rc.3"
+.mule/reports/log/20260915T170704Z-hk-zero-retention-asymmetry-done.md:2677:"name": "v1.4.0-rc.2"
+.mule/reports/log/20260915T170704Z-hk-zero-retention-asymmetry-done.md:2678:"name": "v1.4.0-rc.1"
+
+→ Read .mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md
+The prior astarte-upstream run (2026-09-08) already found **no upstream movement since v1.4.0-rc.5**, and today's release list is identical — newest is still `v1.4.0-rc.5` (2026-08-20), no new stable or rc. All post-1.2.2 items are already in the queue (`compat-note-v1.3.3`, `compat-note-v1.4-rc`, probes). Let me grab the current timestamp and sha to write the evidence report like prior runs.
+
+$ date -u +%Y%m%dT%H%M%SZ; git rev-parse --short HEAD 2>/dev/null
+20260916T184159Z
+ed715ec
+
+← Write .mule/reports/log/20260916T184159Z-recipe-astarte-upstream-noop.md
+Wrote file successfully.
+
+Done: upstream watch — no upstream movement since v1.4.0-rc.5; appended nothing to .mule/todo.md
+
+The release list is unchanged since the prior run (2026-09-08): newest tag `v1.4.0-rc.5` (2026-08-20), newest stable `v1.3.3`. Upstream is ahead of Astrate's 1.2.2 target, but every post-1.2.2 item is already proposed in the queue from earlier runs (`compat-note-v1.3.3`, `compat-note-v1.4-rc`, and the associated probes) — re-proposing would duplicate. Per the recipe's own rule, an empty result is the good result.
+
+Files: `.mule/reports/log/20260916T184159Z-recipe-astarte-upstream-noop.md` (evidence; `todo.md` untouched)
+Verified: `gh api repos/astarte-platform/astarte/releases` -> `v1.4.0-rc.5` newest, unchanged since 2026-08-20 -> pass
+Unsure: nothing
+Follow-ups: none

$ date -u +%Y%m%dT%H%M%SZ; git rev-parse --short HEAD 2>/dev/null; ls .mule/reports/log/ | tail -5
20260917T190650Z
5a314c7
20260917T111017Z-docs-sync-pairing-version-value-done.md
20260917T114541Z-docs-sync-native-version-value-done.md
20260917T182237Z-recipe-hygiene-timeout.md
20260917T184319Z-recipe-milestones-timeout.md
20260917T190323Z-recipe-github-issues-proposed.md

Upstream is unchanged since the last run: still `v1.4.0-rc.5` newest (2026-08-20), `v1.3.3` newest stable. Recipe rule says append nothing. Writing the evidence file.
← Write .mule/reports/log/20260917T190650Z-recipe-astarte-upstream-noop.md
Wrote file successfully.

No upstream movement since last run.

Done: upstream watch — no upstream movement since v1.4.0-rc.5 (2026-08-20); newest stable still v1.3.3; appended nothing to .mule/todo.md
Files: .mule/reports/log/20260917T190650Z-recipe-astarte-upstream-noop.md
Verified: `gh api repos/astarte-platform/astarte/releases` -> v1.4.0-rc.5 newest, unchanged since 2026-08-20 -> pass
Unsure: nothing
Follow-ups: none
