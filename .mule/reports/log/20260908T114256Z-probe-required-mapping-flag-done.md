slug: probe-required-mapping-flag
verdict: done
at:  db9d06e
ran: 2026-09-08T11:42:56Z on DietPi in 426s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
✱ Grep "required" in internal/realm · 0 matches
✱ Grep "Mapping|mapping" in internal/realm · 60 matches
This is a `[readonly]` probe — investigate whether Astrate handles the `required` mapping field, then report the gap without patching.

Let me look at the relevant code paths.
✱ Grep "struct" in internal · 100 matches
✱ Grep "type Mapping struct" 1 match
✱ Grep "type Interface struct" in internal · 0 matches
→ Read pkg/interfaceschema/parse.go [offset=60, limit=60]
✱ Grep "Required|required|Encrypted|encrypted" in pkg/interfaceschema · 14 matches
→ Read pkg/interfaceschema/parse.go [limit=120, offset=119]
→ Read pkg/interfaceschema/parse.go [limit=90, offset=340]
✱ Grep "Required" in . · 61 matches
✱ Grep "required" in . · 1 match
✱ Grep "compiled|Compile|trie" in internal · 100 matches
Already finding that `required` seems largely implemented. Let me verify the RM install/persist path and the enforcement path carefully.
✱ Grep "ParseInterface|Install|install" in internal/realm · 61 matches
→ Read internal/engine/data.go [limit=120, offset=100]
→ Read pkg/interfaceschema/compile.go [limit=120, offset=30]
→ Read pkg/payload/payload.go [limit=120, offset=90]
✱ Grep "DecodeObject|ObjectLeaves|Loader|NewDecoder" in internal/engine · 5 matches
→ Read internal/store/interfaces.go
→ Read internal/realm/detailed.go

$ git log --oneline -20
6c8dbf1 mule: refill the queue
f117601 mule: recipe astarte-upstream ran (305s)
1e56fda mule: refill the queue
c5f3ec2 mule: recipe github-issues ran (74s)
01cf228 mule: recipe milestones ran (215s)
ac78a09 mule: recipe hygiene ran (110s)
4dee74d mule: log docs-sync-appengine-data-output-params
262fdd5 mule: docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
e16cf72 mule: refill the queue
8b575c3 mule: recipe docs-sync ran (341s)
a42f2cf mule: blocked appengine-group-token-roundtrip-test
1f8c775 mule: blocked appengine-snapshot-ignores-query-params
7cbb4b1 mule: refill the queue
99d1a8c mule: recipe code-review ran (153s)
ab49dcc mule: recipe astarte-upstream ran (127s)
08dee34 mule: refill the queue
faf977e mule: recipe github-issues ran (55s)
17b595c mule: recipe milestones ran (312s)
9d822c3 mule: recipe hygiene ran (156s)
17c6cff mule: log docs-sync-appengine-data-422-interface-level

→ Read internal/realm/http.go [limit=80, offset=160]
✱ Grep "renderDetailed|detailed" in internal/realm · 15 matches
✱ Grep "payload.Object|DecodeObject|\.Object\(|Required" in internal/appengine · 0 matches
✱ Grep "Object|Decode|CompiledInterface|leaves|trie" in internal/appengine · 54 matches
→ Read internal/engine/serverdata.go [limit=60, offset=100]

$ ls .mule/recipes/ 2>/dev/null; ls /tmp | rg -i 'astarte|upstream' 2>/dev/null
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

$ rg -n "astarte-upstream|upstream" .mule/recipes/*.md | head -30
.mule/recipes/milestones.md:20:   - Read the reference (upstream repo docs/README, or the linked doc site) for what the
.mule/recipes/milestones.md:21:     milestone actually requires, same care as `.mule/recipes/astarte-upstream.md` — a
.mule/recipes/milestones.md:28:     reconcile each row against the targeted upstream level (promoted → keep,
.mule/recipes/milestones.md:29:     deprecated → deprecate in sync); if the milestone completes a full upstream
.mule/recipes/milestones.md:77:- Never port upstream code or its language-specific structure — port the capability,
.mule/recipes/milestones.md:79:- If the reference named in the milestone section (upstream repos, docs) is unreachable,
.mule/recipes/astarte-upstream.md:1:# Recipe — watch upstream Astarte
.mule/recipes/astarte-upstream.md:8:you which upstream version Astrate currently targets. It is also on the never-touch list:
.mule/recipes/astarte-upstream.md:21:If upstream is ahead:
.mule/recipes/astarte-upstream.md:28:Astrate implements, look at the actual upstream change:
.mule/recipes/astarte-upstream.md:43:  structure, a limit that upstream learned it needed. Propose these as *investigation* tasks:
.mule/recipes/astarte-upstream.md:44:  `- [ ] probe-<slug>: does Astrate have the problem upstream's <change> fixes? report, do not patch`
.mule/recipes/astarte-upstream.md:59:- Never copy upstream code into this repo. Different licence, different language, different
.mule/recipes/astarte-upstream.md:61:- If nothing changed since last time, append nothing to the queue and say "no upstream
.mule/recipes/benchmarks.md:4:against **either** Astrate or upstream Astarte (`bench/README.md`, ~40 lines — read it).
.mule/recipes/benchmarks.md:61:- Comparing against upstream Astarte means the *same tier on the same host*, and
.mule/recipes/cross-project-survey.md:4:This is a bigger, slower job than a normal recipe: it clones two upstream repos and reads
.mule/recipes/cross-project-survey.md:20:   - `git -C .mule/research/upstream/astarte log -1` after a fresh shallow pull (or
.mule/recipes/cross-project-survey.md:22:     re-clone yet) — has upstream's default branch moved?
.mule/recipes/cross-project-survey.md:28:   is a good result; say so and stop, exactly as `.mule/recipes/astarte-upstream.md` already
.mule/recipes/cross-project-survey.md:32:   trigger docs changed upstream. Read the diff (`git log -p <old>..<new> -- doc/` in the
.mule/recipes/cross-project-survey.md:33:   upstream clone) to know what to focus on, the same way `.mule/recipes/astarte-upstream.md`
.mule/recipes/cross-project-survey.md:45:3. **Upstream `astarte-platform`** — shallow clone into `.mule/research/upstream/astarte`
.mule/recipes/cross-project-survey.md:48:5. **Upstream `atomvm/AtomVM`** — shallow clone into `.mule/research/upstream/AtomVM`, same
.mule/recipes/cross-project-survey.md:60:- **Never copy upstream code into this repo.** Read the rule or behaviour, propose that
.mule/recipes/cross-project-survey.md:64:  design-decision) and a check for the kind of trap the live-upstream-recording candidates
.mule/recipes/cross-project-survey.md:76:Before finishing, `rm -rf .mule/research/upstream` — the clones are disposable and the Pi's
.mule/recipes/legion-go.md:4:It already runs the full upstream Astarte stack in Docker — Scylla, VerneMQ, RabbitMQ and

$ ls .mule/research/upstream/ 2>/dev/null; ls .mule/research/ 2>/dev/null | head
(no output)

✱ Grep "issue-67|#67|required" in .mule · 100 matches
→ Read pkg/interfaceschema/required_encrypted_test.go

$ rg -n "required" docs/UPSTREAM-EXPERIMENTAL.md 2>/dev/null | head; rg -n "required" docs/COMPATIBILITY.md 2>/dev/null | head
(no output)

The implementation clearly landed as issue #67 (closed 2026-09-04). Let me verify two candidate gaps against upstream: install-time validation scope (object-aggregation only?) and the detailed RM listing re-render. First let me check the stored docs.
◈ Parallel Web Search "Astarte interface JSON mapping required flag object aggregation 1.4 validation docs"
→ Read internal/engine/data.go [limit=70, offset=255]
→ Read pkg/interfaceschema/compat.go [limit=30, offset=100]
→ Read internal/realm/service.go [limit=70, offset=149]
→ Read pkg/payload/required_test.go

$ tail -15 .mule/todo.md
- [x] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] housekeeping-tests [auto]: add a container-free unit test file for `internal/housekeeping` (service.go, 262 lines, zero test coverage). Cover two rules: (1) `CreateRealm` validation — blank name, blank JWT key, negative regLimit, negative retention, and default retention injection when `WithDefaultDatastreamMaximumStorageRetention` is set and the caller passes nil; (2) `DeleteRealm` gating — `ErrDeletionDisabled` when `WithRealmDeletionDisabled(true)`, `ErrConnectedDevicesPresent` when `DeviceStats` reports connected > 0. Use httptest-style fakes for the store and sealer (no Docker).
- [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
- [x] docs-sync-appengine-by-alias-endpoints [auto]: add the 7 undocumented by-alias mirror routes to docs/api/astarte_appengine_api.yaml — PATCH /devices-by-alias/{alias}, GET /devices-by-alias/{alias}/interfaces, GET/PUT/POST/DELETE /devices-by-alias/{alias}/interfaces/{interface}/{path} (plus the interface-scoped GET) — all registered in internal/appengine/http.go:59-65 and backed by handlers patchDeviceByAlias, listInterfacesByAlias, getDataByAlias, putDataByAlias, deleteDataByAlias (http.go:226-367). Match the existing spec style (same schemas, same response codes as the device-scoped counterparts). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-appengine-group-endpoints [auto]: add the 8 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (returns {data: {group_name: "<name>"}}, http.go:428-434), PATCH /groups/{group}/devices/{device} (DevicePatch body → updated status, http.go:518-537), and the 6 group interface-data routes: the interface-level GET and the path-level GET/PUT/POST/DELETE on /groups/{group}/devices/{device}/interfaces (http.go:75-80, handlers listInterfacesInGroup, getDataInGroup, putDataInGroup, deleteDataInGroup at http.go:187-377). Also add GET /devices/{device}/interfaces (listDeviceInterfaces, http.go:48) which is missing from the spec — that is the 9th route. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-appengine-get-group-device [auto]: add the undocumented `GET /groups/{group}/devices/{device}` route to docs/api/astarte_appengine_api.yaml — it exists in code (internal/appengine/http.go:72, handler getGroupDevice at http.go:509-516, returns the device status projection behind the membership gate) but the spec path only documents delete and patch. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups/{group}/devices is missing `from_token` and `limit` pagination params that code reads at http.go:439-441 (and returns via links.next), (2) POST /groups is missing 422 validation-error response (FieldErrors for blank group_name/empty devices, http.go:404-416), (3) PATCH /devices/{device} and PATCH /devices-by-alias/{alias} are missing 422 ("Attribute key not found" via ErrAttributeKeyNotFound, service.go:352-357 → http.go:656-657), and the device-scoped PATCH is also missing 409 ("Alias already in use", http.go:651) and its documented 500 covers the merge-patch content-type mismatch which code answers with an unmapped 500 too, (4) GET /devices/{device}/interfaces/{interface}/{path} is missing 422 (invalid since/since_after/to/limit/downsample_to/format params, http.go:554-639). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.
- [x] docs-sync-appengine-group-patch-status [auto]: PATCH /groups/{group}/devices/{device} (internal/appengine/http.go:518-537, patchGroupDevice) shares applyPatch with the device-scoped PATCH and also answers 422 ("Attribute key not found", ErrAttributeKeyNotFound, service.go:352-357) and 409 ("Alias already in use", http.go:651), but docs/api/astarte_appengine_api.yaml lists only 200/400/401/404/500 for it — add the two missing responses.
- [x] docs-sync-appengine-data-422-interface-level [auto]: the interface-level data GETs — GET /devices/{device}/interfaces/{interface}, GET /devices-by-alias/{alias}/interfaces/{interface}, GET /groups/{group}/devices/{device}/interfaces/{interface} — share parseQueryOpts (internal/appengine/http.go:554-639) and answer 422 FieldErrors for invalid since/since_after/to/limit/downsample_to/format params, but the spec documents only the path-level GET with 422 — add 422 to the interface-level GETs.
- [!] appengine-snapshot-ignores-query-params [auto]: in `internal/appengine/data.go:138-148`, the individual-datastream interface-root snapshot branch ignores `Since`/`SinceAfter`/`To`/`Limit`/`Descending` even though the root GET is documented with all of them (docs/api/astarte_appengine_api.yaml:577-582) — the same function refuses incompatible params elsewhere (downsample 422, data.go:118-120/165-167), so this is a silent drop. Probe upstream on the Legion Go to decide honour-vs-refuse, then either apply the window/limit or answer an explicit 422 and stop advertising the params on the root GET; add a T1 test in data_test.go/query_opts_test.go pinning the chosen behaviour. — BLOCKED: wrote nothing
- [!] appengine-group-token-roundtrip-test [auto]: add a container-free unit test for `parseGroupToken`/`groupTokenFor` (internal/appengine/service.go:602-623) — currently only exercised via the pg-backed groups_parity_test.go. Assert (1) round-trip `groupTokenFor(offset)` → `parseGroupToken` recovers offset, (2) a non-v1-nibble UUID is rejected, (3) a well-formed v1 UUID with zeroed time fields is accepted (parseGroupToken comment, service.go:614-616). — BLOCKED: lint failed: internal/appengine/groups_token_test.go:14:5: redefines-builtin-id: redefinition of the built-in function max (revive)
- [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
- [ ] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it.

The probe is complete: both behaviors are already implemented (issue #67, closed 2026-09-04). Let me append the one genuine residual gap as a follow-up task line.
$ tail -1 .mule/todo.md | cut -c1-40
- [ ] compat-note-v1.4-rc [auto]: propos

→ Read .mule/todo.md [offset=117, limit=10]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -119,4 +119,5 @@
 - [!] appengine-group-token-roundtrip-test [auto]: add a container-free unit test for `parseGroupToken`/`groupTokenFor` (internal/appengine/service.go:602-623) — currently only exercised via the pg-backed groups_parity_test.go. Assert (1) round-trip `groupTokenFor(offset)` → `parseGroupToken` recovers offset, (2) a non-v1-nibble UUID is rejected, (3) a well-formed v1 UUID with zeroed time fields is accepted (parseGroupToken comment, service.go:614-616). — BLOCKED: lint failed: internal/appengine/groups_token_test.go:14:5: redefines-builtin-id: redefinition of the built-in function max (revive)
 - [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
 - [ ] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it.
+- [ ] detailed-listing-required-encrypted-flags: emit the `required`/`encrypted` mapping flags in `writeDetailedMapping` (internal/realm/detailed.go) so `GET /realmmanagement/v1/{realm}/interfaces?detailed=true` round-trips them — the detailed listing currently materializes a fixed field set that drops both flags even when a datastream mapping declares one, while the plain GET preserves them via the stored definition. Datastream branch only (properties reject both at parse); add a T1 assertion in internal/realm/detailed_test.go.

Done: `required` mapping flag already fully implemented (issue #67, closed 2026-09-04) — RM accepts+persists it at install, and the engine data path rejects object writes omitting it; no gap in the two probed behaviors, one cosmetic gap found.
Files: .mule/todo.md (one follow-up line appended; no code touched)
Verified: read-only probe of install/parse/compile/persist/enforce paths (`rg`, file reads, git log to confirm #67 landed) -> pass; no commands that change state run
Unsure: (1) whether upstream rejects `required` on *individual*-aggregation datastreams at install — Astrate's parser accepts it on any datastream (only properties are rejected), runtime effect is nil since enforcement runs only via ObjectLeaves; not probe-verified against live upstream. (2) the `?detailed=true` listing drop is cosmetic (plain GET/persist are unaffected).
Follow-ups: `detailed-listing-required-encrypted-flags` appended to .mule/todo.md

Findings detail:
- RM install accept+persist: `required` parses on datastream mappings and is rejected on properties (pkg/interfaceschema/parse.go:395,410); install stores the definition verbatim (internal/store/interfaces.go:70-71) and plain GET returns it (internal/realm/http.go:178,200); minor-upgrade treats a `required` flip as incompatible (pkg/interfaceschema/compat.go:121). Pinned by required_encrypted_test.go.
- Write rejection: `payload.Object` rejects a document missing any Required leaf with ReasonMissingRequired in both JSON and BSON (pkg/payload/payload.go:137-152, payload/required_test.go); both ingest paths route through it — device publishes (internal/engine/data.go:278) and server-owned writes (internal/engine/serverdata.go:139).
- Residual gap: `writeDetailedMapping` (internal/realm/detailed.go:56-75) materializes a fixed field set for `GET /interfaces?detailed=true` and drops `required`/`encrypted` even when declared; plain GET still shows them. Flagged as the follow-up line.
