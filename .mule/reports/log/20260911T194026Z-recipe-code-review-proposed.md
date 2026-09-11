slug: recipe-code-review
verdict: proposed
at:  ee76ba1
ran: 2026-09-11T19:40:26Z on DietPi in 430s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md

$ ls internal/ pkg/ 2>/dev/null
internal/:
appengine
auth
broker
config
engine
flow
flowapi
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil

pkg/:
astarteapi
deviceid
interfaceschema
payload

Reviewed both files. Now let me survey the codebase and check which areas were reviewed recently.
$ ls -la .mule/reviews/ 2>/dev/null; echo "---TODO---"; cat .mule/todo.md 2>/dev/null
total 64
drwxr-xr-x 2 root root 4096 Sep 10 19:00 .
drwxr-xr-x 7 root root 4096 Sep 11 20:29 ..
-rw-r--r-- 1 root root 4600 Sep  7 20:15 appengine-2026-09-07.md
-rw-r--r-- 1 root root 8188 Sep  4 20:43 broker-2026-09-04.md
-rw-r--r-- 1 root root 3141 Sep  6 20:57 engine-2026-09-06.md
-rw-r--r-- 1 root root 5516 Sep  5 21:19 flow-2026-09-05.md
-rw-r--r-- 1 root root 4837 Sep 10 19:00 flowapi-2026-09-10.md
-rw-r--r-- 1 root root 3506 Sep  9 13:07 pairing-2026-09-09.md
-rw-r--r-- 1 root root 4259 Sep  9 21:13 realm-2026-09-09.md
-rw-r--r-- 1 root root 4605 Sep  8 18:59 store-2026-09-08.md
---TODO---
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
- [x] docs-sync-rm-delete-interface-status [auto]: fix the `deleteInterface` response in docs/api/astarte_realm_management_api.yaml — the spec documents `422` with examples "Interface major version is not 0, can't be deleted" / "Cannot delete an interface that is used by a device introspection" (yaml:255-272), but the handler returns `403` with "Interface can't be deleted" / "Interface can't be deleted since it's currently used" (internal/realm/http.go:387-391, `writeInterfaceError`; both are 403 `StatusForbidden`). The spec's detail strings were lifted from the `writeError` branch (http.go:414-419) which DELETE /interfaces does not use. Fix status to `403` and update the two examples to match the actual wire messages. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-mapping-required-encrypted [auto]: add `required` (boolean, datastream-only) and `encrypted` (boolean, datastream-only) to the `InterfaceMapping` schema in docs/api/astarte_realm_management_api.yaml — the parser accepts and persists both on datastream mappings (internal/interfaceschema/parse.go:410-414), rejects them on properties mappings (parse.go:395-398), the detailed listing emits them (`detailed.go:69`), and the compatibility check compares them (interfaceschema/compat.go:121-123); the schema (yaml:820-862) currently has `additionalProperties: false` and neither field, so the OpenAPI would mark a valid datastream interface install body as invalid. Also add both fields to the detailed-listing example (yaml:70-85). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-put-auth-422 [auto]: add the missing `422` (ValidationError) response to `PUT /realmmanagement/v1/{realm}/config/auth` in docs/api/astarte_realm_management_api.yaml — `SetAuthKey` returns `ErrValidation` for a blank `jwt_public_key_pem` (internal/realm/service.go:632-640), which `writeError` maps to 422; the spec (yaml:431-464) documents only 204/400/401/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-rm-version-example [auto]: fix the `getVersion` example in docs/api/astarte_realm_management_api.yaml from `data: "1.1.0"` (yaml:557) to `data: "1.2.2"` — `APICompatVersion` is `"1.2.2"` (internal/realm/service.go:588); the native spec's four compat-version examples already use "1.2.2". Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1

$ tail -50 .mule/log.md 2>/dev/null
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
| 2026-09-08 | compat-note-v14-rc | blocked | 326s | wrote nothing |
| 2026-09-08 | detailed-listing-required-encrypted-flags | done | 159s | 63b629e |
| 2026-09-08 | store-alias-values-taken-test | done | 85s | a07bd6a |
| 2026-09-08 | store-latest-individual-test | done | 134s | 741794c |
| 2026-09-08 | docs-sync-appengine-device-status-schema | done | 255s | 00da133 |
| 2026-09-08 | docs-sync-appengine-data-set-422 | blocked | 251s | wrote nothing |
| 2026-09-08 | docs-sync-appengine-downsample-min | done | 86s | db506f9 |
| 2026-09-09 | pairing-bearer-secret-test | done | 90s | 32544c0 |
| 2026-09-09 | pairing-remoteip-fallback-test | done | 469s | 3d6e790 |
| 2026-09-09 | docs-sync-hk-patch-endpoint | done | 403s | 9fc9f33 |
| 2026-09-09 | docs-sync-hk-retention-field | done | 142s | d6d7c7f |
| 2026-09-09 | realm-pure-helper-tests | done | 403s | 39eb70a |
| 2026-09-09 | docs-sync-native-compat-health-503 | done | 202s | 656e21e |
| 2026-09-09 | docs-sync-native-compat-version-endpoints | done | 480s | d148f16 |
| 2026-09-10 | flowapi-autorestart-terminal-failure | done | 251s | 8ce04d2 |
| 2026-09-10 | flowapi-validationdetail-test | done | 121s | 832d8fb |
| 2026-09-10 | flowapi-autorestart-default-test | done | 148s | e272cfa |
| 2026-09-10 | docs-sync-pairing-status-enum | done | 89s | 2f6176c |
| 2026-09-10 | docs-sync-pairing-version-endpoint | done | 95s | c867b2c |
| 2026-09-11 | docs-sync-rm-delete-interface-status | done | 123s | b37d11f |
| 2026-09-11 | docs-sync-rm-mapping-required-encrypted | done | 137s | 5f8cb25 |
| 2026-09-11 | docs-sync-rm-put-auth-422 | done | 145s | 8554720 |
| 2026-09-11 | docs-sync-rm-version-example | blocked | 97s | opencode exited 1 |

$ for d in internal/auth internal/config internal/housekeeping internal/httpx internal/observability internal/swagger pkg/astarteapi pkg/deviceid pkg/interfaceschema pkg/payload internal/testutil; do printf "%-30s " "$d"; rg -c '' --glob '*.go' "$d" 2>/dev/null | wc -l; done
internal/auth                  8
internal/config                2
internal/housekeeping          5
internal/httpx                 4
internal/observability         5
internal/swagger               2
pkg/astarteapi                 2
pkg/deviceid                   2
pkg/interfaceschema            12
pkg/payload                    10
internal/testutil              6

$ for d in pkg/interfaceschema pkg/payload internal/auth internal/httpx internal/observability internal/config internal/swagger; do printf "%-24s " "$d"; rg -c '' --glob '*.go' "$d" | awk -F: '{s+=$2} END{print s" lines"}'; done
pkg/interfaceschema      3176 lines
pkg/payload              2673 lines
internal/auth            1532 lines
internal/httpx           357 lines
internal/observability   398 lines
internal/config          818 lines
internal/swagger         171 lines

Reviewed-packages set: appengine, broker, engine, flow, flowapi, pairing, realm, store. Unreviewed candidates include `pkg/interfaceschema`, `pkg/payload`, `internal/auth`, `internal/httpx`, `internal/observability`. Let me size them by lines.
$ ls -la pkg/interfaceschema pkg/payload 2>/dev/null
pkg/interfaceschema:
total 136
drwxr-xr-x 3 root root  4096 Sep  4 20:46 .
drwxr-xr-x 6 root root  4096 Jul 27 05:10 ..
-rw-r--r-- 1 root root  4880 Sep  4 20:46 compat.go
-rw-r--r-- 1 root root  9730 Sep  4 20:46 compat_test.go
-rw-r--r-- 1 root root  4175 Sep  4 20:46 compile.go
-rw-r--r-- 1 root root  6752 Jul 27 05:10 compile_test.go
-rw-r--r-- 1 root root 24107 Sep  4 20:46 parse.go
-rw-r--r-- 1 root root 10954 Jul 27 05:10 parse_test.go
-rw-r--r-- 1 root root  6693 Sep  4 20:46 required_encrypted_test.go
drwxr-xr-x 4 root root  4096 Jul 27 05:10 testdata
-rw-r--r-- 1 root root  5283 Sep  4 20:46 trie.go
-rw-r--r-- 1 root root  6060 Jul 27 05:10 trie_test.go
-rw-r--r-- 1 root root 12186 Jul 27 05:10 types.go
-rw-r--r-- 1 root root  2710 Sep  4 20:46 violations.go
-rw-r--r-- 1 root root 11632 Sep  4 20:46 violations_test.go

pkg/payload:
total 124
drwxr-xr-x 3 root root  4096 Sep  4 20:46 .
drwxr-xr-x 6 root root  4096 Jul 27 05:10 ..
-rw-r--r-- 1 root root 12064 Jul 27 05:10 bson.go
-rw-r--r-- 1 root root  3804 Jul 27 05:10 fuzz_test.go
-rw-r--r-- 1 root root 13391 Jul 27 05:10 json.go
-rw-r--r-- 1 root root  6932 Sep  4 20:46 payload.go
-rw-r--r-- 1 root root 34751 Jul 27 05:10 payload_test.go
-rw-r--r-- 1 root root  2590 Sep  4 20:46 required_test.go
-rw-r--r-- 1 root root  1220 Jul 27 05:10 sniff.go
-rw-r--r-- 1 root root  4091 Jul 27 05:10 sniff_test.go
drwxr-xr-x 3 root root  4096 Jul 27 05:10 testdata
-rw-r--r-- 1 root root 13274 Sep  4 20:46 value.go
-rw-r--r-- 1 root root  7975 Sep  4 20:46 value_test.go

I'll review `pkg/interfaceschema` (3176 lines, unreviewed, wire-critical for interface installs). Reading largest files first.
→ Read pkg/interfaceschema/parse.go
→ Read pkg/interfaceschema/compat.go
→ Read pkg/interfaceschema/compile.go
→ Read pkg/interfaceschema/trie.go
→ Read pkg/interfaceschema/violations.go
→ Read pkg/interfaceschema/types.go
→ Read pkg/interfaceschema/trie_test.go
Now the tests, to see what's already covered.
→ Read pkg/interfaceschema/compat_test.go
→ Read pkg/interfaceschema/required_encrypted_test.go
→ Read pkg/interfaceschema/compile_test.go
Now the remaining tests to map coverage.
→ Read pkg/interfaceschema/parse_test.go
→ Read pkg/interfaceschema/violations_test.go

$ cat pkg/interfaceschema/testdata/invalid/manifest.json
{
  "01-bad-name.json": "invalid interface name",
  "02-name-too-long.json": "exceeds 128 characters",
  "03-version-zero.json": "cannot both be 0",
  "04-negative-major.json": "version_major must be >= 0",
  "05-no-mappings.json": "at least one mapping",
  "06-missing-type.json": "missing \"type\"",
  "07-bad-type.json": "unknown interface type \"stream\"",
  "08-bad-ownership.json": "unknown ownership \"cloud\"",
  "09-bad-aggregation.json": "unknown aggregation \"objects\"",
  "10-bad-value-type.json": "unknown value type \"float\"",
  "11-endpoint-not-rooted.json": "must start with '/'",
  "12-endpoint-partial-placeholder.json": "placeholders must span a whole segment",
  "13-endpoint-duplicate.json": "duplicate endpoint",
  "14-endpoint-conflict.json": "conflicting endpoints",
  "15-object-on-properties.json": "properties interfaces cannot use object aggregation",
  "16-object-mixed-depth.json": "must all have the same depth",
  "17-object-mixed-attrs.json": "share the same reliability",
  "18-object-divergent-prefix.json": "must share the same prefix",
  "19-object-param-leaf.json": "must end in a literal segment",
  "20-properties-explicit-timestamp.json": "\"explicit_timestamp\" is not allowed on properties",
  "21-allow-unset-datastream.json": "\"allow_unset\" is only allowed on properties",
  "22-ttl-without-policy.json": "database_retention_ttl requires database_retention_policy \"use_ttl\"",
  "23-use-ttl-missing-ttl.json": "requires database_retention_ttl >= 1",
  "24-unknown-field.json": "unknown field",
  "25-retention-on-properties.json": "\"retention\" is not allowed on properties",
  "26-negative-expiry.json": "expiry must be >= 0",
  "27-trailing-garbage.json": "trailing data"
}

Coverage looks thorough. Let me verify my one candidate gap (interior-literal + param backtracking in the trie) against real coverage before proposing anything.
$ go test ./pkg/interfaceschema/ -coverprofile=/tmp/iface_cover.out 2>&1 | tail -5 && go tool cover -func=/tmp/iface_cover.out
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	0.153s	coverage: 87.3% of statements
github.com/astrate-platform/astrate/pkg/interfaceschema/compat.go:32:		CheckMinorUpgrade		93.9%
github.com/astrate-platform/astrate/pkg/interfaceschema/compat.go:92:		normalizedKey			75.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/compat.go:102:		sameMappingAttributes		66.7%
github.com/astrate-platform/astrate/pkg/interfaceschema/compile.go:47:		Compile				85.7%
github.com/astrate-platform/astrate/pkg/interfaceschema/compile.go:111:		lastLiteralSegment		71.4%
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:160:		ParseInterface			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:171:		ParseInterfaceCanonical		95.5%
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:209:		buildInterface			92.8%
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:342:		tooLongMessage			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:351:		buildMapping			89.1%
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:473:		splitEndpoint			90.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:505:		parseSegment			78.6%
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:530:		validIdentifier			77.8%
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:551:		endsWithSlash			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:560:		checkEndpointUniqueness		100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:599:		checkObjectAggregation		100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:626:		sameObjectAttributes		54.5%
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:653:		normalizeEndpoint		100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/trie.go:61:		NewEndpointTrie			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/trie.go:71:		Add				100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/trie.go:113:		Match				100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/trie.go:124:		matchNode			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/trie.go:156:		placeholderValueOK		100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:31:		ParseInterfaceType		100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:43:		String				75.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:55:		UnmarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:58:		MarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:73:		ParseOwnership			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:85:		String				75.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:97:		UnmarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:100:		MarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:118:		ParseAggregation		100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:130:		String				75.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:142:		UnmarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:145:		MarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:164:		ParseReliability		80.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:178:		String				80.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:193:		QoS				100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:196:		UnmarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:199:		MarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:218:		ParseRetention			60.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:232:		String				0.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:246:		UnmarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:249:		MarshalJSON			0.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:267:		ParseDatabaseRetentionPolicy	75.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:279:		String				0.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:291:		UnmarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:296:		MarshalJSON			0.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:344:		ParseValueType			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:354:		String				66.7%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:362:		IsArray				100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:366:		Elem				100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:374:		UnmarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:377:		MarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:387:		unmarshalEnum			87.5%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:401:		marshalEnum			66.7%
github.com/astrate-platform/astrate/pkg/interfaceschema/violations.go:36:	Error				0.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/violations.go:38:	Unwrap				100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/violations.go:41:	add				100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/violations.go:46:	addMapping			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/violations.go:51:	empty				100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/violations.go:55:	errOrNil			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/violations.go:66:	Error				91.7%
total:										(statements)			87.3%

$ cat pkg/interfaceschema/testdata/invalid/12-endpoint-partial-placeholder.json; echo; ls pkg/interfaceschema/testdata/invalid/ | head -40
{
  "interface_name": "com.astrate.test.PartialPlaceholder",
  "version_major": 1,
  "version_minor": 0,
  "type": "datastream",
  "ownership": "device",
  "mappings": [{"endpoint": "/foo%{bar}/value", "type": "double"}]
}

01-bad-name.json
02-name-too-long.json
03-version-zero.json
04-negative-major.json
05-no-mappings.json
06-missing-type.json
07-bad-type.json
08-bad-ownership.json
09-bad-aggregation.json
10-bad-value-type.json
11-endpoint-not-rooted.json
12-endpoint-partial-placeholder.json
13-endpoint-duplicate.json
14-endpoint-conflict.json
15-object-on-properties.json
16-object-mixed-depth.json
17-object-mixed-attrs.json
18-object-divergent-prefix.json
19-object-param-leaf.json
20-properties-explicit-timestamp.json
21-allow-unset-datastream.json
22-ttl-without-policy.json
23-use-ttl-missing-ttl.json
24-unknown-field.json
25-retention-on-properties.json
26-negative-expiry.json
27-trailing-garbage.json
manifest.json

`matchNode` is statement-covered; my trie corner is masked. Let me confirm which parseSegment param-error branches and enum round-trips are genuinely unpinned.
$ go tool cover -html=/tmp/iface_cover.out -o /tmp/iface_cover.html 2>/dev/null; rg -n "String|MarshalJSON" pkg/interfaceschema/types.go | head -30
42:// String returns the wire form.
43:func (t InterfaceType) String() string {
57:// MarshalJSON encodes the wire form.
58:func (t InterfaceType) MarshalJSON() ([]byte, error) { return marshalEnum(t, Datastream, Properties) }
84:// String returns the wire form.
85:func (o Ownership) String() string {
99:// MarshalJSON encodes the wire form.
100:func (o Ownership) MarshalJSON() ([]byte, error) {
129:// String returns the wire form.
130:func (a Aggregation) String() string {
144:// MarshalJSON encodes the wire form.
145:func (a Aggregation) MarshalJSON() ([]byte, error) {
177:// String returns the wire form.
178:func (r Reliability) String() string {
198:// MarshalJSON encodes the wire form.
199:func (r Reliability) MarshalJSON() ([]byte, error) {
231:// String returns the wire form.
232:func (r Retention) String() string {
248:// MarshalJSON encodes the wire form.
249:func (r Retention) MarshalJSON() ([]byte, error) {
278:// String returns the wire form.
279:func (p DatabaseRetentionPolicy) String() string {
295:// MarshalJSON encodes the wire form.
296:func (p DatabaseRetentionPolicy) MarshalJSON() ([]byte, error) { return marshalEnum(p, NoTTL, UseTTL) }
314:	// String is a UTF-8 string (≤ 64 KiB).
315:	String
328:	// StringArray is a homogeneous array of String.
329:	StringArray
353:// String returns the wire form.
354:func (v ValueType) String() string {

All four gaps confirmed. Let me check an existing review's style before writing mine.
→ Read .mule/reviews/store-2026-09-08.md
Style confirmed. Writing the review file, then appending task lines.

← Write .mule/reviews/interfaceschema-2026-09-11.md
Wrote file successfully.

Now appending the five task lines to the queue.
$ tail -3 .mule/todo.md
- [x] docs-sync-rm-mapping-required-encrypted [auto]: add `required` (boolean, datastream-only) and `encrypted` (boolean, datastream-only) to the `InterfaceMapping` schema in docs/api/astarte_realm_management_api.yaml — the parser accepts and persists both on datastream mappings (internal/interfaceschema/parse.go:410-414), rejects them on properties mappings (parse.go:395-398), the detailed listing emits them (`detailed.go:69`), and the compatibility check compares them (interfaceschema/compat.go:121-123); the schema (yaml:820-862) currently has `additionalProperties: false` and neither field, so the OpenAPI would mark a valid datastream interface install body as invalid. Also add both fields to the detailed-listing example (yaml:70-85). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-put-auth-422 [auto]: add the missing `422` (ValidationError) response to `PUT /realmmanagement/v1/{realm}/config/auth` in docs/api/astarte_realm_management_api.yaml — `SetAuthKey` returns `ErrValidation` for a blank `jwt_public_key_pem` (internal/realm/service.go:632-640), which `writeError` maps to 422; the spec (yaml:431-464) documents only 204/400/401/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-rm-version-example [auto]: fix the `getVersion` example in docs/api/astarte_realm_management_api.yaml from `data: "1.1.0"` (yaml:557) to `data: "1.2.2"` — `APICompatVersion` is `"1.2.2"` (internal/realm/service.go:588); the native spec's four compat-version examples already use "1.2.2". Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1

$ cat >> .mule/todo.md <<'EOF'
- [ ] interfaceschema-compat-attrs-flip-test [auto]: extend TestCheckMinorUpgradeTable (pkg/interfaceschema/compat_test.go) to assert the four currently-unpinned immutable attributes — retention, expiry, database_retention_policy/ttl (datastream old/next pair) and allow_unset (properties old/next pair) — each flipped on upgrade and rejected as ErrIncompatibleEndpointChange with an at-immutability acceptance twin; only type/reliability/explicit_timestamp/required/encrypted are tested today (sameMappingAttributes, compat.go:102, sits at 66.7% stmt coverage), so a regression relaxing expiry or DB-TTL immutability passes the suite.
- [ ] interfaceschema-object-attrs-uniformity [auto]: pin the object-aggregation attribute-uniformity rule for the four unasserted attributes — retention, expiry, explicit_timestamp, database_retention_policy/ttl — via ParseInterface rejections ("must all share the same <attr>", pkg/interfaceschema/parse.go:626-647) each paired with an at-uniformity acceptance twin; only reliability is pinned today (testdata/invalid/17-object-mixed-attrs.json, sameObjectAttributes at 54.5% stmt coverage).
- [ ] interfaceschema-malformed-placeholder-fixtures [auto]: add invalid fixtures + manifest entries for parseSegment's four untested param branches (pkg/interfaceschema/parse.go:509-517) — "/%{", "/%ax", "/%{x", "/%{1}" — each rejected with the malformed-placeholder / invalid-placeholder-name message; the only placeholder-syntax fixture today is the mid-segment 12-endpoint-partial-placeholder.json (ContainsAny branch), so a regression accepting a dangling "/%{" passes every test.
- [ ] interfaceschema-enum-roundtrip-coverage [auto]: extend TestEnumRoundTrips (pkg/interfaceschema/parse_test.go:244) to cover Ownership, Aggregation, Retention, DatabaseRetentionPolicy — the four enums it omits (Retention.String/MarshalJSON and DatabaseRetentionPolicy.String at 0% stmt coverage in the package's own tests) — so a typo'd wire string like "retention":"descard" or "aggregation":"obilver" fails the suite instead of passing.
- [ ] interfaceschema-trie-interior-literal-param [auto]: pin the rule that an interior-only literal trie child does not shadow its parametric sibling at the matching depth (pkg/interfaceschema/trie.go:124-151) — purpose-build a trie with "/a/b/c" and "/a/%{p}", assert Match("/a/b") returns the /a/%{p} mapping and Match("/a/b/c") the /a/b/c one; statement coverage cannot see this branch (matchNode reports 100%), so a "fix" that misses on shallow paths through an interior literal node sails through the suite.
EOF
echo appended
appended

$ git status --porcelain; echo "---"; tail -6 .mule/todo.md | cut -c1-80
M .mule/todo.md
?? .mule/reviews/interfaceschema-2026-09-11.md
---
- [!] docs-sync-rm-version-example [auto]: fix the `getVersion` example in docs/
- [ ] interfaceschema-compat-attrs-flip-test [auto]: extend TestCheckMinorUpgrad
- [ ] interfaceschema-object-attrs-uniformity [auto]: pin the object-aggregation
- [ ] interfaceschema-malformed-placeholder-fixtures [auto]: add invalid fixture
- [ ] interfaceschema-enum-roundtrip-coverage [auto]: extend TestEnumRoundTrips 
- [ ] interfaceschema-trie-interior-literal-param [auto]: pin the rule that an i

Done: reviewed pkg/interfaceschema (unreviewed rotation pick) and proposed five test-only tasks.
Files: .mule/reviews/interfaceschema-2026-09-11.md (new), .mule/todo.md (5 lines appended)
Verified: `go test ./pkg/interfaceschema/ -coverprofile` -> pass, 87.3% stmt; coverage output grounded every proposal. No git ops performed.
Unsure: nothing — I did not probe upstream, so the interfaceNameRe shape is noted in the review as deliberately un-proposed.
Follow-ups: 5 `[auto]` lines in .mule/todo.md — compat-attrs-flip-test, object-attrs-uniformity, malformed-placeholder-fixtures, enum-roundtrip-coverage, trie-interior-literal-param (details in the review file).

Findings were all missing-test gaps on existing rules (no behavioural defect found): sameMappingAttributes leaves 4 attrs unasserted, object-aggregation uniformity pins only reliability, parseSegment's param-error branches untested, enum round-trip omits 4 enums, and the trie's interior-literal-vs-param corner is invisible to statement coverage.
