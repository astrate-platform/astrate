slug: recipe-github-issues
verdict: proposed
at:  7bc8abd
ran: 2026-09-10T20:01:20Z on DietPi in 63s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md

$ cat /root/astrate-mule/.mule/todo.md 2>/dev/null; echo "---EXIT:$?"
# Mule queue

One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).
The mule runs the topmost open line. Reorder freely; the order is the priority.

If a line needs more detail than fits, write `.mule/tasks/<slug>.md` and keep the line short.

    tools/mule.sh add "<slug>: <outcome>"
    tools/mule.sh loop

A line tagged `[legion]` needs the Legion Go and is skipped automatically while it is asleep;
the queue moves on to the next runnable line rather than stalling behind it.

A line tagged `[readonly]` is a **standing check**, not a piece of work. It verifies something
and is expected to produce no diff at all — so an empty diff is success, not the `wrote
nothing` failure it would be for any other line. It never gets ticked off, because a gate is
never done. Its output lands in `.mule/reports/<slug>.md` with the sha it ran against, and it
is skipped while that sha is still HEAD: re-checking code that has not moved buys nothing and
spends a call on a free provider.

**`mule/queue` is authoritative for this file.** The mule ticks on the Pi and writes its
progress there, so editing the queue on `main` resurrects completed tasks the moment the two
are merged — which has happened once. Add and reorder tasks on `mule/queue`; let them reach
`main` only when that branch is merged.

A line marked `- [~]` is **parked**: real work, but not tick-sized. The mule only ever
picks up `- [ ]`. Benchmark runs live here because a single tier is 5-20 minutes of ingest
alone -- past the per-task budget -- and because they want someone watching. Run one by hand:

    tools/mule.sh legion bench-push
    ssh legion 'cd ~/astrate/bench && ./scripts/run-tier.sh small astrate -base-url ... -housekeeping-key ...'

## Where tasks come from

**This file is not the whole queue, and for real work it is not even the main part of it.**

The queue is: the standing lines below, plus **every open GitHub issue labelled `mule`**.
Issues are read live on each tick and are never copied into this file — a copy would be a
second place the same fact lives, on a branch the mule commits to and you edit on `main`,
and that produced three merge conflicts in one afternoon.

**To give the mule work, file an issue and label it `mule`.** From anywhere, by anyone,
including another model with repo access. No SSH, no editing this file:

    gh issue create --label mule --title "<slug>: <outcome>" --body "<the detail>"

Labels on the issue are the tags: `legion` and `readonly` mean what `[legion]` and
`[readonly]` mean here. State lives on the issue, as labels, because there is exactly one
copy of it there:

| label          | meaning                                                          |
|----------------|------------------------------------------------------------------|
| `mule`         | queued                                                            |
| `mule-review`  | the mule pushed something; **it is not merged and not reviewed**  |
| `mule-blocked` | it tried and could not; re-label `mule` to try again              |

The mule never closes an issue. Whether the work actually resolves it is a judgement about
intent, which is the reviewer's call.

When both sources are empty a tick runs a **proposal recipe** instead, rotating through
`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
stuck re-reviewing the same code. Lines it invents are tagged `[auto]`: nobody approved those.
A refill costs a tick from the daily budget and never runs what it just proposed — the lines
sit for one tick, which is your window to cut a bad one. `MULE_NO_REFILL=1` turns it off.

## Nothing merges on its own

Everything lands on `mule/queue`. The gates prove a change compiles, passes the tests, ships
a test that fails without it, and touches no frozen file — none of which means the change is
worth having. Before any of it reaches `main`:

    bash tools/mule.sh review

## Queue

- [x] flowapi-autorestart-terminal-failure: `restartWithBackoff` (internal/flowapi/service.go:585-621) retries forever on a permanent rebuild failure — a flow whose pipeline was deleted (service.go:360-363 wraps it as ErrValidation) or whose config/definition no longer resolves keeps marking itself failed and rebuilding every ≤30s with no exit short of flow deletion. Classify terminal vs transient (ErrValidation → permanent; store/manager errors → transient) via a pure helper and stop the loop on terminal failures; add a container-free unit test on the helper. [auto]
- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.go:560) with no stop signal, so a block-death racing process shutdown (cmd/astrate/main.go:482-488 drains the Manager but never cancels the goroutine) can rebuild the flow via `mgr.StartFlow` on a background pump (internal/flow/flow.go:173-175) after `Manager.Shutdown` and flip the durable status back to running. Thread a stop channel/context through the Service, checked in the loop's sleep, cancelled on shutdown; verification needs a `[legion]` integration or timing-based test. [auto]
- [x] flowapi-validationdetail-test: add a container-free unit test for `validationDetail` (internal/flowapi/http.go:297-304) — the 422-detail prefix-strip rule currently has zero coverage. Assert (1) `"flowapi: validation failed: bad block"` strips to `"bad block"`, (2) a non-prefix message passes through unchanged, (3) a message exactly equal to the prefix is NOT stripped (the guard is `len(msg) > len(prefix)`). Pure helper; no store. [auto]
- [x] flowapi-autorestart-default-test: pin the wire rule "POST /flows without `auto_restart` defaults to true" (internal/flowapi/http.go:137-144) — currently untested, so a regression to default-false would pass every existing suite. Extract the default into a pure helper (missing → true, present → honored) and assert all three cases container-free. [auto]
- [x] bench-tiers: create the tiered benchmark definitions per .mule/recipes/benchmarks.md (first run only — this task builds the harness, it does not run it)
- [~] bench-small-astrate: run `bench/scripts/run-tier.sh small astrate` against local Astrate, commit results (two runs minimum) [legion]
- [~] bench-medium-astrate: run `bench/scripts/run-tier.sh medium astrate` against local Astrate, commit results (two runs minimum) [legion]
- [~] bench-small-astarte: run `bench/scripts/run-tier.sh small astarte` against local Astarte, commit results (two runs minimum) [legion]
- [~] bench-medium-astarte: run `bench/scripts/run-tier.sh medium astarte` against local Astarte, commit results (two runs minimum) [legion]
- [~] bench-big-astrate [legion]: run `bench/scripts/run-tier.sh big astrate` against Legion Go Astrate, commit results (two runs minimum)
- [~] bench-giant-astrate [legion]: run `bench/scripts/run-tier.sh giant astrate` against Legion Go Astrate, commit results (two runs minimum)
- [~] bench-big-astarte [legion]: run `bench/scripts/run-tier.sh big astarte` against Legion Go Astarte, commit results (two runs minimum)
- [~] bench-giant-astarte [legion]: run `bench/scripts/run-tier.sh giant astarte` against Legion Go Astarte, commit results (two runs minimum)
- [ ] race-check-store: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/store/... ./internal/housekeeping/... ./migrations/...`. Report any failure to .mule/for-giulio.md with the full race report. Split out of the former single `race-check` line, which timed out running the whole tree at once. [legion] [readonly]
- [ ] race-check-engine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/engine/... ./internal/broker/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
- [ ] race-check-flow: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/flow/... ./internal/realm/... ./internal/pairing/... ./internal/auth/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
- [ ] race-check-appengine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/appengine/... ./internal/observability/... ./internal/httpx/... ./internal/config/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
- [ ] race-check-pkg: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./pkg/... ./cmd/... ./internal/testutil/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]

- [x] broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [approved 2026-09-04]
- [!] broker-disconnect-device-zombie-session: in `internal/broker/broker.go` `DisconnectDevice` (lines 260-266), force-closing the live connection leaves the persisted bbolt session (subscriptions + offline queue, keyed by CN) in place — mochi derives `expire` from the Clean flag, not the reason code (server.go:485-491), so `dropClient` (sessionstore.go:500) is never reached and a later device with the same CN resurrects the deleted device's session. Wire the session store handle into `Broker` and call `dropClient` from `DisconnectDevice`, and add a T1 test in `broker_test.go`: connect clean=false + subscribe, call `DisconnectDevice`, reconnect — assert `session_present` is false and the stale subscription is gone. [approved 2026-09-04] — BLOCKED: wrote nothing
- [!] broker-offline-acl-tests: in `internal/broker/aclhook_test.go`, unit-test the offline-delivery ACL — `offlineACL.ownershipOf` cache hit within TTL, TTL expiry triggering a reload, and a load failure caching as empty and denying (aclhook.go:116-138) — plus a T1 `OnACLCheck` check that the offline branch (delivery to a session-present but disconnected device, aclhook.go:196-202) consults the store. The rule currently has zero direct tests; `TestCheckACLMatrix` exercises only the pure `checkACL`. [approved 2026-09-04] — BLOCKED: wrote nothing
- [!] broker-onconnect-doc-comment: in `internal/broker/authhook.go:314`, restore the missing first line of the `OnConnect` doc comment — it opens mid-sentence ("of the Astarte MQTT v1 protocol, and mochi publishes them without a publish-side ACL check…") and the Will-clearing security rationale (a retained LWT would escape the §3.2 matrix) is only half-documented. [approved 2026-09-04] — BLOCKED: tests failed: --- FAIL: TestMQTTSink_Retained (0.02s)

- [x] purge-properties-compression-capability [auto]: upstream v1.3.0 adds a `purge_properties_compression_format` device capability (`zlib`|`plaintext`, default `zlib`) — a wire-visible capability value. Check whether Astrate's capabilities handling (internal/broker, the `<realm>/<device_id>/capabilities` topic, issue #16) needs to recognize/honour it, or whether zlib-only is already the deliberate default; propose the change or note why not needed.
- [!] empty-introspection-verification [auto]: upstream v1.3.0 changed "allow devices with empty introspection" — verify whether Astrate's device connection/introspection handling currently rejects an empty introspection string where upstream now accepts it, and propose a fix if so. — BLOCKED: wrote nothing
- [!] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages" — investigate only: does Astrate have an analogous delay between trigger install and first delivery? Report, do not patch. — BLOCKED: wrote nothing
- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
- [x] docs-sync-pairing-health-path [auto]: add the undocumented `GET /pairing/v1/{realm}/health` route to docs/api/astarte_pairing_api.yaml — it exists in code (internal/pairing/http.go:81) and returns 200 `{"data":{"status":"ok"}}` for an existing realm, 404 for an unknown realm, 503 when unhealthy, unauthenticated and unrate-limited. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-register-404 [auto]: document the missing `404` (DeviceNotFound) on `POST /pairing/v1/{realm}/agent/devices` in docs/api/astarte_pairing_api.yaml — the handler returns it when `GetRealmByName` fails for an unknown realm (internal/pairing/service.go:183-186 → http.go:326-327), but the spec lists only 201/400/401/403/422/429/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] swagger-httptest-coverage [auto]: add a container-free httptest suite for internal/swagger, which currently has no test file — cover two rules: (1) `Mount` wires `GET /swagger` to a 302→/swagger/index.html, serves the embedded UI under `/swagger/`, and serves the OpenAPI YAML specs under `/api/`; (2) `Specs` returns exactly the `.yaml` filenames from `docs.APIYAML` with no path prefix or dirs. httptest only, no Docker.
- [!] probe-property-resend-encoding [auto]: upstream v1.3.0 fixed outbound server-property values sent to a device on connect/emptyCache (commit 522ccf4f — a raw stored binaryblob is now wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding, so it ships as a BSON binary, not a string). Investigate only: does Astrate's resendServerProperties/rehydrate path (internal/engine/control.go:144, pkg/payload) emit binaryblob → BSON subtype-0 binary and datetime → UTC datetime exactly like upstream 1.3? Report, do not patch. — BLOCKED: wrote nothing
- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
- [!] flow-validate-source-sink: `Pipeline.Validate`'s no-source/no-sink error branches (internal/flow/pipeline.go:148-153) are unreachable — the preceding Kahn cycle check (pipeline.go:100-115) guarantees any graph that gets past it is acyclic and therefore always has a source and a sink. The two test cases "no source (all blocks have incoming edges)" and "no sink (all blocks have outgoing edges)" (pipeline_test.go:134-163) both use a 2-cycle, so they pass via the cycle error, never the branch they name. Make the rule real (a cycle-independent formulation, e.g. reject a graph where a node has no path to a sink) or delete the vacuous check, and fix the two test cases to assert the actual error message so a regression to wrong-path coverage fails. — BLOCKED: tests failed: --- FAIL: TestMQTTSink_Retained (0.02s)
- [x] flow-validate-dead-source-sink-recompute: in `Pipeline.Validate`, the first `hasSource`/`hasSink` computation (pipeline.go:117-126) is overwritten by the second (pipeline.go:128-147) and its result discarded — dead code (the comment at pipeline.go:127 documents it). Delete the first block; behaviour is unchanged. Tests: existing `TestPipeline_Validate` passes.
- [x] docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` route to docs/api/astarte_realm_management_api.yaml — it exists in code (internal/realm/http.go:51, handler getDatastreamMaximumStorageRetention at http.go:135, returns 200 with a data-envelope retention value via Service.GetDatastreamMaximumStorageRetention) but the spec jumps straight from config/device_registration_limit (yaml:428) to /version. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
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
- [x] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
- [x] detailed-listing-required-encrypted-flags: emit the `required`/`encrypted` mapping flags in `writeDetailedMapping` (internal/realm/detailed.go) so `GET /realmmanagement/v1/{realm}/interfaces?detailed=true` round-trips them — the detailed listing currently materializes a fixed field set that drops both flags even when a datastream mapping declares one, while the plain GET preserves them via the stored definition. Datastream branch only (properties reject both at parse); add a T1 assertion in internal/realm/detailed_test.go.

- [x] store-alias-values-taken-test: add a store test for `AliasValuesTaken` (internal/store/devices.go:343, the `alias_already_in_use` ownership check behind PATCH /devices) pinning the three branches that currently have no direct coverage — no other device in the realm carries any of the values → false; another device carries one → true; the same device (excluded via `d.id <> $2`) carries the value → false (must not count itself). Add to internal/store/devices_test.go. [auto]
- [x] store-latest-individual-test: add a store test for `LatestIndividual` (internal/store/datastreams.go:247, the previous-value lookup feeding value_change/path_created triggers) — only exercised through the engine's fakeStore today; assert (a) of two samples on one series, `ORDER BY ts DESC LIMIT 1` returns the newest with the right value column, (b) an empty series yields ErrNotFound. Add to internal/store/datastreams_test.go. [auto]
- [x] docs-sync-appengine-device-status-schema [auto]: fix the `DeviceStatus` schema in docs/api/astrate_appengine_api.yaml — the wire emits `id` (internal/appengine/service.go:106 `json:"id"`, decoded as `ds.ID` in mirror_test.go:156 and dashboard_compat_test.go:100/149) but the spec names the field `device_id` (yaml:1409) in both the schema and the list/detail examples (yaml:67, 75); `introspection` is documented as `map<string,integer>` (yaml:1438-1441) but each value is a `{major, minor}` object (introspectionEntry, service.go:123-127, introspectionView service.go:833-836); and the schema omits the six fields `groups`, `total_received_msgs`, `total_received_bytes`, `first_credentials_request`, `last_seen_ip` (omitempty), `previous_interfaces` (omitempty) that deviceStatusView emits (service.go:254-279). Rename `device_id`→`id`, correct `introspection` to `map<string,{major:int,minor:int}>`, add the missing fields, and fix the examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-appengine-data-set-422 [auto]: add the missing `422` response to the six PUT/POST data-set operations in docs/api/astrate_appengine_api.yaml — `putData`, `putDataByAlias`, `putDataInGroup` (internal/appengine/http.go:292-309, 321-338, 340-357) return 422 `{"errors":{"detail":"Value size exceeds size limits"}}` when the body value exceeds 65536 bytes (maxValueBytes, http.go:24), both for PUT and its POST alias (routes registered at http.go:52-53, 63-64, 78-79), but the spec lists only 200/400/401/404/500 for all six operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-appengine-downsample-min [auto]: fix the `downsample_to` constraint in docs/api/astrate_appengine_api.yaml — the DataDownsample description says "(must be >= 2)" (yaml:1357) but parseQueryOpts rejects `downsample_to <= 2` with "must be greater than 2" (internal/appengine/http.go:592-593), so the documented minimum is off by one: readings of `2` (allowed per the doc) actually 422. Change the description to "(must be > 2)" (the 422 error shape is already documented). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

- [x] pairing-bearer-secret-test [auto]: add a container-free unit test for `bearerSecret` (internal/pairing/http.go:343) locking the upstream-verbatim permissive match `~r/bearer\:?\s+(.*)$/i` — the colon form `Bearer: <secret>`, the case variants (`bearer `, `BEARER `), and the reject paths (no whitespace before value, non-bearer scheme, empty value). Only the plain `Bearer <secret>` form is exercised today (integration http_test.go); a regression to case-sensitive or drop-colon matching would pass every existing test. Pure helper in a non-integration file; no Docker.
- [x] pairing-remoteip-fallback-test [auto]: add a container-free unit test for `remoteIP` (internal/pairing/http.go:359) — a valid addr:port yields `ap.Addr().Unmap()`, a bare address yields the mapped address, and an unparseable `RemoteAddr` (unix-socket/proxy case) falls back to `netip.IPv4Unspecified()`. That fallback branch is currently uncovered (only reachable through the integration HTTP stack). Pure helper in a non-integration file; no Docker.
- [x] docs-sync-hk-patch-endpoint [auto]: add the undocumented `PATCH /housekeeping/v1/realms/{realm}` route to docs/api/astarte_housekeeping_api.yaml — it exists in code (internal/housekeeping/http.go:40, handler patchRealm at http.go:163-229, returns 200 on success) but the spec jumps straight from GET to DELETE. Document the request body (`RealmPatch` — same three optional fields as `RealmCreate`: `jwt_public_key_pem`, `device_registration_limit`, `datastream_maximum_storage_retention`, all nullable), the success response, and the error responses: 400 (bad envelope), 401, 403, 404, 405 (deletion disabled), 422 (invalid update parameters for unknown fields, blank jwt_public_key_pem, negative limit/retention, connected devices present), 500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-retention-field [auto]: add the `datastream_maximum_storage_retention` field to the `Realm` and `RealmCreate` schemas in docs/api/astarte_housekeeping_api.yaml — the field is on the wire (realmBody struct, internal/housekeeping/http.go:50, returned by viewBody at http.go:99, accepted in create at http.go:68 and patch at http.go:216) but absent from both schemas (yaml:200-231). Type: integer, format int64, nullable, description matching the device_registration_limit style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

- [x] realm-pure-helper-tests [auto]: add container-free unit tests for the four untested pure helpers in internal/realm. Create `internal/realm/helpers_test.go` (no build tag) covering: (1) `normaliseIfaceName` — lowercase, hyphen strip, mixed case+hyphens, empty string, already-normalized passthrough; (2) `joinPEM` — empty slice, single key, two keys separated by `\n`; (3) `triggerErrorBody` — empty action+simple_triggers, action-only, simple_triggers-only with index-aligned `{}` substitution, both present; (4) `majorParam` — valid positive, zero, negative (→404), non-numeric (→404). All tests use httptest.ResponseRecorder for majorParam; the rest are direct function calls. No Docker.

- [x] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmmanagement/health`, `GET /pairing/health` — in docs/api/astrate_native_api.yaml, and widen the shared `DataEnvelopeHealthStatus.status` enum (yaml:301-303) from `[ok]` to `[ok, unhealthy]` — the handler answers 503 `{"data":{"status":"unhealthy"}}` when the readiness probe fails (internal/observability/compat.go:19-34, deviation 18), but the spec documents only 200 with `status: ok` for all three paths. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-status-enum [auto]: fix the `PairingInfo.status` enum in docs/api/astarte_pairing_api.yaml from `[confirmed, pending, denied, expired]` (yaml:412) to `[confirmed, pending, inhibited]` — the handler never emits `denied` or `expired`: `service.Info` maps only `DeviceStatusInhibited`→"inhibited", `DeviceStatusConfirmed`→"confirmed", and everything else→"pending" (internal/pairing/service.go:296-299, statuses defined in store/devices.go:19-27; upstream-parity per the comment at service.go:282-285). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-version-endpoint [auto]: add the undocumented `GET /pairing/v1/{realm}/version` route to docs/api/astarte_pairing_api.yaml — it exists in code (cmd/astrate/main.go:398, unauthenticated, upstream-parity per main.go:387-391) and answers 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38); mirror the RM spec's getVersion op (astarte_realm_management_api.yaml:532) without its a_rma security. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
---EXIT:0

$ cat /root/astrate-mule/.mule/for-giulio.md 2>/dev/null; echo "---EXIT:$?"
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

- **Extension point: external-bus Intake (NATS JetStream), 2026-09-10.** `internal/broker/intake.go:56-62` carries a `TODO(extension point)` for a second `Intake` implementation backed by a durable external bus (NATS JetStream) for multi-instance deployment or restart survival. The frozen design decision is in `docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4`; the consumer must reproduce per-device ordering, deferred-ack backpressure, and QoS 0 drop semantics. Not tracked by any open issue. Your call: file a milestone issue or park until the roadmap calls for it.

- **Extension point: timescaledb_toolkit lttb downsampling, 2026-09-10.** `internal/store/store.go:139-143` carries a `TODO(extension point)` for switching the `Downsample` method from `time_bucket+avg` to the toolkit's `lttb()` when `timescaledb_toolkit` is present — the probe at `store.go:144-150` already records availability in `s.hasToolkit`, but `datastreams.go` always uses the default path. Design reference: `docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5`. Not tracked by any open issue. Your call: file a milestone issue or park.

---

- **docs-sync pairing, 2026-09-10: docs/site names the pre-credentials wire status `registered`; the API returns `pending`.** `docs/site/pairing-and-security.md:56` ("flips status `registered -> confirmed`") and `:91` ("**registered** -- device registered, awaiting first credentials request"), plus `docs/site/data-modeling.md:94`, all give `registered` as a device status value. But `service.Info` emits `pending` for any device that is neither confirmed nor inhibited (internal/pairing/service.go:296-299 — upstream-parity per the comment at service.go:282-285; the DB value is `registered`, store/devices.go:19-27). Wire values are `pending`/`confirmed`/`inhibited`. Site prose is yours — reword to `pending`, or confirm the site intentionally describes the DB value.

---

- **Flow code review, 2026-09-10: `DeletePipeline` semantics need your call.**
  `DeletePipeline` (internal/flowapi/service.go:199-206) removes a pipeline
  silently even when durable flows still reference it; running instances keep
  the old graph, but their next reload fails ("pipeline not found") and a
  block death after that starts the auto-restart loop retrying forever against
  the missing pipeline. `UpdatePipeline` already reports `referencing_flows`
  (service.go:162-182, issue #44) — should DELETE also refuse (409) while
  flows reference it, or return the referencing names and allow? Design choice,
  not a mule task.

---

- **github-issues triage run, 2026-09-10: nothing proposable, nothing stale.** Nine open
  issues: #99–#94 are mule-alarms (idle queue), not code issues. #93 has a pushed commit
  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated. #1
  untouched per standing instruction. No task lines proposed, no for-giulio close suggestions.

- **Milestone recipe run, 2026-09-09: v3.0 unchanged since 2026-09-06 — still waiting on
  upstream gates.** `milestone-3.0` label remains empty. The three open `upstream-parity`
  issues are the same: #92 keyAgreement (upstream decision, gated on stable v1.4.0), #93
  aclhook comment (pushed commit `8c61268`, mule-review), #78 FDO (milestone-4.0). Both
  UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait
  on upstream v1.4 final (currently rc.5). `APICompatVersion` stays 1.2.2 per #90's frozen
  decision. No new gaps found, no issues filed, no task lines proposed.

---

- **github-issues triage run, 2026-09-08: nothing proposable, nothing stale.** All six open
  issues were already covered by another mechanism. #97/#96/#95/#94 are mule-alarms (idle
  queue), not code issues — expected during low-activity windows. #93 has a pushed commit
  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated below.
  #1 untouched per standing instruction. No task lines proposed, no for-giulio close
  suggestions. (Identical to the 2026-09-07 run — issue list unchanged except the new
  alarm #97.)

---

- **Milestone recipe run, 2026-09-06: v3.0 (upstream 1.2.2 → 1.3/1.4) has no open
  machine-checkable gaps left to file; the milestone now waits on two upstream gates, and the
  `milestone-3.0` label is empty (this milestone's work lives under `upstream-parity`).** All
  the 1.3-surface gaps from the #47–#89 backlog are delivered/closed — retention ceiling #72,
  alias/group #52–#54 #59, query formats #55 #56, wire capabilities #47–#49, validation /
  error-code fidelity #57 #61 #62 #79, per-service version #77, housekeeping #73–#76 — and the
  remaining open items are already escalated, not fileable gaps: **#92 keyAgreement** (1.4
  experimental decision, escalated 2026-09-05, gated on a stable upstream v1.4.0), **#78 FDO**
  (milestone-4.0, escalated), **#93** aclhook comment rewrite (mule-review, pushed commit
  `8c61268`). `docs/UPSTREAM-EXPERIMENTAL.md` rows #67/#68 are "1.4 experimental" and reconcile
  only when upstream ships **v1.4 final** — not yet (still rc). Per **#90's frozen decision**
  (2026-08-23), `APICompatVersion` stays 1.2.2 until the FULL 1.3 surface is complete *and*
  every UPSTREAM-EXPERIMENTAL row at that level is reconciled, so the bump is not yet due. Your
  call when v1.4.0 goes stable: answer #92 (implement against the rc, re-park, or narrow #93
  only) — until then v3.0 has no open implementable work. This run filed **no issues** and
  proposes **no** "complete, cut the tag" line: the open decision set is real and is yours.

- **`milestones.md` v3.0 names a source doc that does not exist: `.mule/research/
  upstream-parity-2026-08-22.md` is absent from `main` and from `origin/mule/research`.** The
  investigation's *output* (issues #47–#89) is intact and reachable, so milestone work is
  unaffected, but if that doc is your canonical 1.3/1.4 diff it needs restoring (or the
  milestone paragraph re-pointed) before v4.0 drafting leans on it.

- **docs-sync realm-management: two ASTRATE_ config keys exist in code but are missing from
  `docs/site/configuration-reference.md`** (site prose is yours, never edited by the mule).
  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` and
  `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` are consumed in `internal/config/config.go`
  (lines 270-271 and 288); neither appears in `docs/site/configuration-reference.md`.
  Add them or confirm they are meant to be internal-only.

- **github-issues triage run, 2026-09-06: nothing proposable, nothing stale.** All six open
  issues were already covered by another mechanism. #95 and #94 are mule-alarms (idle queue),
  not code issues — expected during low-activity windows. #93 has a pushed commit (`8c61268`)
  awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the waiting-on
  row. #78 is the milestone-4.0 FDO design/investigation, already escalated below. #1
  untouched per standing instruction. No task lines proposed, no for-giulio close suggestions.
  (Identical to 2026-09-05 run — issue list unchanged except new alarm #95.)

---

- **keyAgreement: the parking condition from #51 has fired — implement now, or wait for a
  stable v1.4.0?** (issue #92, `upstream-parity`/`upstream-experimental`). #51 closed
  2026-08-22 with "parked until the upstream 1.4 experimental spec stabilizes — reopen or
  file fresh when it does". The document side has now stabilized: upstream `d084308`
  (2026-08-31) published `082-key_agreement_protocol.md`, a full 267-line wire protocol
  (topics `control/keyAgreement/0..4` — InitExchange/ExchangeResp/SecretHash/HashOk/
  ExchangeFailed at QoS 2, CBOR bodies with CDDL, `alg` 0 ECDH_P256-HKDF_SHA256-AES_256_GCM,
  CBOR-wrapped COSE_Key + 32-byte HkdfSalt, session-scoped keys, enumerated `ExchangeFailed`
  codes), and deleted the "not yet implemented" sentences the old parking quoted (#93 already
  fixed the stale ACL comment that cited them). But the spec ships only in `v1.4.0-rc.5` —
  `v1.3.3` is still the newest stable tag and Astrate targets 1.2.2 — so the document has
  stabilized and the release has not. Implementing it is the largest surface in the parity
  backlog (CBOR codec, X25519/P-256, HKDF, AES-256-GCM, a 5-state handshake machine,
  shared-secret persistence, five new error names). Your call: build against the rc now,
  re-park until v1.4.0 is a stable tag, or take only the narrow #93 fix. (Escalated again
  2026-09-05 — a prior escalation from the 2026-09-04 milestone run was lost in the queue
  rebuild.)

---

- **govulncheck GO-2026-5970: reachable DoS in golang.org/x/text (infinite loop on invalid input, fixed in v0.39.0, available v0.41.0).** Astrate pins `x/text` indirect at v0.38.0 (go.mod:97) and pgx pulls it into production: `internal/store/notify.go:59` `store.Listen` → `pgx.ConnectConfig` → `unicode/norm.*`. This is the only govulncheck symbol finding that is not test-harness-only: GO-2026-6355/6354 (x/crypto/ssh deadlocked-channel DoS) and GO-2026-6253 (moby/go-archive tar path traversal) are reachable only through testcontainers in `internal/testutil/pg.go`, i.e. never in the deployed binary. `x/text` keeps API compatibility minor-to-minor and the modules Astrate exercises (`unicode/norm` via pgx, `text/language` via jsonschema) are unchanged, so this is a fix Astrate actually needs — the hygiene recipe's highest-priority category. Not a mule task (go.mod never-touch): your decision to bump ≥v0.39.0 now or fold into the next milestone-boundary sweep. Raw: https://pkg.go.dev/vuln/GO-2026-5970. (The 2026-09-04 dep sweep did not list x/text.)

---

- **COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable, 2026-07-14; v1.4.0 is still rc-only).** Astrate's doc and `APICompatVersion` still target upstream **1.2.2** (`internal/realm/service.go:588`); v1.3.0 (2026-05-06) introduced wire-surface changes Astrate does not yet emulate, so this is a decision — adopt v1.3.2 as the compatibility target (then update the doc + bump `APICompatVersion` together, per the bump rule) or keep 1.2.2 and add a "not yet emulated" note. Wire-relevant v1.3.0 deltas (release notes): **MQTT v1 capabilities** incl. `purge_properties_compression_format` (plaintext vs zlib purge — touches the `emptyCache`/`producer,properties` contract COMPATIBILITY.md deviation 1 documents); **empty introspection now allowed**; **device registration triggers** (pairing) and **device deletion started/completed triggers** (RM — the latter two already exist as Astrate deviation 9 emits both around the synchronous delete); **FDO authentication** (pairing, disabled by default); **realm-scoped health** — upstream v1.3 added `GET /pairing/v1/{realm}/health`, which Astrate already serves (`internal/pairing/http.go:78-81`, comment already says "upstream 1.3+"), so deviation 18's wording ("which upstream 404s") is now false against 1.3 and the note should be reworded either way. Proposed doc wording (for your approval, edit to taste): in §Infrastructure differences add a sentence — *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl.

plaintext `purge_properties_compression_format`, empty-introspection allowance, device registration/deletion triggers, experimental FDO pairing auth) are not yet emulated and are out of scope until the milestone that adopts v1.3.2 as the target."* — and reword deviation 18's realm-health note from "which upstream 404s" to "added by upstream v1.3 (Astrate serves it against a 1.2.2 target; kept, matching behavior)". Raw upstream changes: [v1.3.0](https://github.com/astarte-platform/astarte/releases/tag/v1.3.0), [v1.3.2](https://github.com/astarte-platform/astarte/releases/tag/v1.3.2).

---
- **The mule's 42 straight failures were ten lint errors in its own base, not a stale branch.**
  `mule/queue` stopped taking `main` on 2026-07-27 and drifted 120 commits behind, which is
  real and is why it is being rebuilt — but it is not what blocked the work. The lint gate runs
  `golangci-lint run ./...` over the whole repo, and the branch's own `internal/flow/*` code
  carried 10 findings (1 goimports, 1 gosec, 8 revive). So every task failed the gate no matter
  how good the change was, including the four tasks queued to fix those very findings: each
  fixed one and died on the other nine. Base tests and `go vet` pass on that checkout, and
  `golangci-lint` and `govulncheck` have been installed on the Pi the whole time — the two
  things that looked broken were not. `main` is lint-clean, so rebuilding from it clears the
  deadlock. `tools/mule.sh preflight` checks the lint baseline and would have said so on day
  one; nobody ran it between 2026-08-31 and 2026-09-04. (Diagnosed 2026-09-04.)

---

- **Dependency sweep corrected: direct (pinned) deps DO have newer versions** — the 2026-09-02 note said the `go list -m -u` sweep showed "only version-skew on transitive deps", but that run hit the recipe's `head -20` cutoff (all cloud/azure/transitive) and never reached the directly-required modules. Full sweep, 2026-09-04. None of these is a fix this repo *needs*, so no bump is proposed — recorded for the decision. Per module (current → available; breaking change; repo use):
  - `github.com/coder/websocket` v1.8.14 → v1.8.15 — no breaking (patch); used in `internal/appengine/stream/ws.go`, `channels/ws.go`; worth it only for the "transmit in single frame when compression enabled" fix + read-path alloc reduction.
  - `go.etcd.io/bbolt` v1.4.3 → v1.5.0 — bbolt's semver promises no API change between patch/minor, so additive-only; used in `internal/broker/sessionstore.go`; v1.5 adds a data-file size limit and panic-recovery hardening, nothing Astrate needs.
  - `go.mongodb.org/mongo-driver/v2` v2.6.0 → v2.8.2 — the 2.8.0 breaking changes are confined to Queryable Encryption string-query options (`options.Text()`→`String()`); Astrate uses only the raw BSON API (`pkg/payload/bson.go`, `internal/engine/capabilities.go`, `bench/`) and is unaffected.
  - `github.com/nats-io/nats.go` v1.52.0 → v1.53.1 — no breaking; the headline fixes (JetStream `resetOrderedConsumer` race, KV dot-rejection) are paths Astrate does not use — `internal/engine/forward/nats.go` is core NATS publish only.
  - `github.com/prometheus/client_golang` v1.23.2 → v1.24.1 — requires Go ≥1.25 (fine, repo is 1.26.1); the breaking `LabelNames`/remote-api renames don't touch repo usage (`prometheus`/`collectors`/`promhttp` in `internal/observability/metrics.go`, flow/engine metrics); would buy `Gather()` panic-recovery and opt-in `CoalesceGather` scrape-pile-up protection.
  - `github.com/testcontainers/testcontainers-go` v0.43.0 → v0.44.0 (modules/postgres v0.42.0, modules/nats v0.43.0) — breaking in `wait.ForSQL` (callback now takes `network.Port`) and `ImageProvider` (new `PullImageWithPlatform`); Astrate's `internal/testutil/pg.go` looks unaffected but it is test-only anyway.
  - `golang.org/x/crypto` v0.53.0 → v0.56.0 — x/crypto keeps API compatibility; used only for bcrypt in `internal/auth`.
  Note (corrected 2026-09-04): `govulncheck` and `golangci-lint` **are** installed on the Pi (`/root/go/bin`, since 2026-07-28 and 2026-09-01) and `.mule/config` finds them there, so both checks were available; the sweep that produced this list simply ran without invoking govulncheck.

  **Decided 2026-09-04: no bumps.** None of the seven fixes anything this repo has, and each
  one costs a full test run to land. The standing rule instead: re-run this sweep at every
  milestone boundary (the point where `APICompatVersion` or a milestone tag moves), and bump
  only what carries a fix Astrate actually needs. `go.mod` stays on the never-touch list.

---

- **milestone 2.0 looks complete, verify and cut the tag** — all 11 `milestone-2.0` issues
  CLOSED (#23–#27, #37, #39–#43), no open issues, no new gaps after re-checking upstream
  astarte_flow block catalog against `internal/flow/` + git log (MQTT/HTTP source/sink,
  json_path_map, pure-transform set, virtual_device_pool, container block MVP, flow API,
  durable named flows all landed).   (Milestones recipe run, 2026-09-03.)

  **Decided 2026-09-04: the tag is `v0.2.0`, cut on the same day.** The `v2.0`/`v3.0` names
  are milestone names, not release versions — the project is still pre-1.0 and the version
  number keeps its own line rather than jumping two majors to match a milestone label.

---
- ~~`docs/site/appengine-api.md` documents `GET` and `DELETE` on `/appengine/v1/<realm>/groups/<name>`~~
  — **resolved 2026-09-04, and the original note was half wrong.** `GET /groups/{group}` does
  exist (`internal/appengine/http.go:69`); only `DELETE` was absent, and it is absent from
  upstream's own spec too (`docs/api/astarte_appengine_api.yaml` has no `/groups/{group}`
  path at all). So the docs line was spurious, not the code. Line dropped.
---
- ~~`Router.Submit` TOCTOU on the `closed` flag~~ — **fixed 2026-09-04**, and none of the three
  options in the original note was taken: (b) was wrong (a `select` does not save a send on a
  closed channel — it panics either way). The fix is that `Drain` no longer closes the lane
  channels at all; it closes `quit` under a write lock while `Submit` holds the read lock
  across its send, so a lane can never be retired underneath an in-flight sender, and the lanes
  exit on `quit` after draining what is buffered. `TestRouter_SubmitParkedWhenDrainRuns`
  reproduces the old panic deterministically and passes on the fix.
---
- ~~**#87 `lua_map` — needs embedded Lua runtime, parked.**~~ — **closed 2026-09-04.**
  Embedding a Lua VM is not machine-checkable by the mule and Lua flow support is on no active
  roadmap. Reopen in a minute if that changes.
- **#78 FDO device onboarding — milestone-4.0, investigation phase.** Too large for a single
  mule task; the investigation work (reading upstream's TO2 handling, inventorying endpoints,
  schema and keys) is a multi-session project. Parked until the v3.0 queue clears and this
  becomes the next milestone target.
- **#1 is never to be raised again.** Giulio's standing instruction, 2026-09-04: it stays open
  permanently and is not a candidate for closing, triage, or a for-giulio entry. Do not
  propose it again.
---

- ~~**Flow v2.0: named multi-instance flows + pipeline config?**~~ — **decided
  2026-07-29: (b) named multi-instance + config.** Design then implement (#40).
  Doc: `docs/handoff/flow-v2-decisions-2026-07-29.md`.
- ~~**Flow v2.0: durable `flows` table vs in-memory only?**~~ — **decided
  2026-07-29: (b) durable records; `auto_restart` default true, optional never;
  process rehydrates on boot; fail loudly on bad pipeline/start.** Filed #41;
  edge-case follow-ups #42. Same design package as #40.
- ~~**Flow v2.0: containers in scope?**~~ — **decided 2026-07-29: yes, phased
  PoC → MVP (#43).** Native Lua/MQTT blocks are *not* a v2.0 gate (containers cover
  custom logic). Doc: `docs/handoff/flow-v2-decisions-2026-07-29.md`.
- ~~`device_deletion_started`/`device_deletion_finished` trigger events are not emitted~~ —
  **decided 2026-07-27: emit both, back-to-back, around the synchronous delete.** Filed as
  issue #21 (`mule`). (Cross-project survey, 2026-07-27,
  `.mule/research/survey-2026-07-27.md` source 4.)
- ~~Mustache trigger-action templates are accepted but not rendered~~ — **decided
  2026-07-27: implement it.** Guiding principle clarified: Astarte compatibility means
  SDK/wire compatibility, not minimum dependency count — Astrate is allowed to be a
  compatible *superset*. Library picked: `github.com/cbroglie/mustache`. Filed as issue #22
  (`mule`). (Same survey, source 4.)
- ~~`value_change`/`value_change_applied`/`path_created`/`path_removed`/`value_stored` trigger
  types compile but never fire~~ — **resolved: implemented in commit 6bd14a7 (2026-08-22)**,
  and the cost question answered by #20's bench on the Legion Go (2026-08-24 closeout):
  ~0.13 ms keyed read per message, −0.7% throughput even paid once per message over REST.
  Performance is not the constraint for these types (nor, a fortiori, for the group-scoped
  line below). (Same survey, source 4.)
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
- ~~golangci-lint is not installed on the Pi~~ — **resolved 2026-07-28**: installed v2.12.2
  via `go install` (the prebuilt-binary installer's published sha256 for linux/arm64 did not
  verify, reproducibly, so built from source instead). `mule.service` now sets
  `MULE_LINT_CMD=golangci-lint run ./...` and has `/root/go/bin` on `PATH`; the lint gate
  runs starting with the next tick.
- ~~`/root/astrate` on the Pi has uncommitted work~~ — **resolved 2026-07-27** with the new
  `tools/reconcile.sh`: rescued onto `origin/wip/DietPi-20260727T171543Z` (pushed, not
  reviewed — read the diff before merging anything from it) and `/root/astrate` is now clean
  on `main`. Also had to set `commit.gpgsign false` locally on that clone first — it had
  signing on with no working gpg-agent for a non-interactive session, same class of problem
  `/root/astrate-mule` already solved the same way (your call, confirmed 2026-07-27: align it
  with the existing unattended-Pi-clone convention rather than fix gpg-agent or bypass signing
  silently inside the script).

## 2026-08-23 — FDO promoted to milestone-4.0 candidate (Giulio's decision, recorded)

#78 is no longer parked: zero-touch onboarding is strategic for commercial
viability. Scope frozen on the issue: owner-side TO1/TO2 in our Pairing
service only (last mile, like upstream), reuse fdo-rs for
manufacturing/rendezvous, acceptance = official `astarte-device-fdo-rust`
SDK completes onboarding against Astrate, docs as a first-class deliverable.
When v3.0 is marked DONE, the v4.0 section of `.mule/milestones.md` should be
drafted with this investigation as its first item (issue #78 has the full
verified context).
---EXIT:0

$ gh issue list --state open --limit 40 \
  --json number,title,labels,updatedAt,comments \
  --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
#99 mule: nothing has landed in 14h [mule-alarm ] []c
#98 mule: nothing has landed in 16h [mule-alarm ] []c
#97 mule: nothing has landed in 15h [mule-alarm ] []c
#96 mule: nothing has landed in 14h [mule-alarm ] []c
#95 mule: nothing has landed in 15h [mule-alarm ] []c
#94 mule: nothing has landed in 15h [mule-alarm ] []c
#93 In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged [mule-review upstream-parity ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:The mule pushed `8c61268` to `mule/queue` for this.

Unreviewed: the gates passed — it compiles, the tests pass, a new test was shown to fail without the change, and no frozen file was touched. That is not the same as the change being worth having. Nothing merges to `main` until someone reads it:

    bash tools/mule.sh review

Left open on purpose. Whether this actually resolves the issue is your call, not mine. createdAt:2026-09-04T19:54:58Z id:IC_kwDORmfsJs8AAAABSo2DbA includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/93#issuecomment-5545755500 viewerDidAuthor:true]]c
#92 keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds [enhancement upstream-parity upstream-experimental ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:Parked by Giulio on 2026-09-04: **wait for a stable \`v1.4.0\`.** The spec is complete but ships only in \`v1.4.0-rc.5\`, and building against an rc risks rebuilding it.

So this issue is not waiting on anyone to remember it. It has a row in [`.mule/waiting-on.md`](https://github.com/astrate-platform/astrate/blob/main/.mule/waiting-on.md), and the weekly `mule-upstream-watch` job now reads that file as its first step every run: when a stable `v1.4.0` tag exists upstream, it escalates to `.mule/for-giulio.md` naming this issue. An `-rc.N` does not satisfy the row.

That mechanism was added because it was missing here: #51 was parked on the same condition, the condition fired on 2026-08-31, and nothing connected the new upstream tags to the parked issue — it took a hand triage four days later to catch it. createdAt:2026-09-04T19:52:21Z id:IC_kwDORmfsJs8AAAABSo0f8A includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/92#issuecomment-5545730032 viewerDidAuthor:true]]c
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
#1 Provide an Open Source IoT Platform unironically, unironically [wontfix ] []c

Now I need to read the existing for-giulio.md entry for this exact recipe run, since I can see the top entry is already from today's github-issues triage. Let me check whether it already covers this run.
→ Read .mule/for-giulio.md [limit=20]
The github-issues triage for today (2026-09-10) is already recorded in `.mule/for-giulio.md` and the result is unchanged — all 9 open issues are already covered by other mechanisms. No task lines to propose, no new for-giulio entry needed.

Done: github-issues triage run — nothing proposable, nothing stale
Files: none modified
Verified: `gh issue list` -> pass (9 issues returned)
Unsure: nothing
Follow-ups: none
