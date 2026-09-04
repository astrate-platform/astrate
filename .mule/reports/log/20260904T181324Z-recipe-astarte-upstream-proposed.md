slug: recipe-astarte-upstream
verdict: proposed
at:  c085169
ran: 2026-09-04T18:13:24Z on DietPi in 72s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
I'll read COMPATIBILITY.md first, then check upstream releases.
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

The recipe says to compare the newest tag to the version named in COMPATIBILITY.md. COMPATIBILITY.md doesn't name a specific Astarte version tag directly — it references upstream behavior "against upstream v1.2.0" and "v1.2.2" for the dashboard/API compat. Let me check the todo.md to see the last recorded upstream movement and what tag was last noted, plus check if there's a specific version reference.

Let me check the todo.md and any for-giulio notes to understand the current target version.
$ cat /root/astrate-mule/.mule/todo.md
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

**Open issues labelled `mule` are worked BEFORE the standing lines below** (changed
2026-08-31). An issue is work a human approved; a standing line is either a gate or something
the mule proposed to itself, and self-proposed work never runs ahead of approved work.

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
- [x] store-realm-cascade-policies: in `internal/store/realms_test.go` `CascadeDelete`, add `trigger_policies` to the post-delete verification loop (query `SELECT count(*) FROM trigger_policies WHERE realm_id = $1`). The migration 000006 has ON DELETE CASCADE but nothing asserts on it. [auto]
- [x] store-alias-lowest-id: in `internal/store/devices_test.go`, add a subtest that registers two devices in the same realm, sets the same alias tag on both, and asserts `GetDeviceByAlias` returns the one with the lower UUID. The SQL uses `ORDER BY id LIMIT 1` but no test proves it. [auto]
- [x] store-delete-device-objects: in `internal/store/devices_test.go` `StatsAndDelete`, insert object datastream rows for the device before deleting it, and assert they are gone after the delete. `DeleteDevice` explicitly sweeps `object_datastreams` but the test only checks individual rows. [auto]
- [!] control-producer-properties-compression: in `internal/engine/control.go`, accept a plaintext (and the 4-zero-byte empty-frame) device→server `producer/properties` purge list when the device declared `purge_properties_compression_format: plaintext` — `inflateProperties` currently rejects every non-zlib payload while upstream's `control_handler.ex` decodes per-capability **Approved by Giulio 2026-08-31**: yes, devices that cannot compress must be able to talk to us, following upstream. — BLOCKED: gates failed
- [!] probe-interface-default-values: does `GET /realmmanagement/v1/<realm>/interfaces/<name>/<major>` return the same defaulted mapping-parameter values upstream reports after 3f0b864 ("Properly show default values for all mapping parameters")? report, do not patch [auto] — BLOCKED: gates failed
- [!] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto] — BLOCKED: gates failed
- [!] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
- [!] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto] — BLOCKED: gates failed
- [!] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
- [!] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto] — BLOCKED: gates failed
- [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
- [!] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto] — BLOCKED: gates failed
- [!] docs-sync-pairing-status-enum: in docs/api/astarte_pairing_api.yaml, fix the PairingInfo.status enum (line 367): it lists `confirmed, pending, denied, expired` but the handler only ever emits `confirmed`, `pending`, or `inhibited` (internal/pairing/service.go:289-297, via internal/pairing/http.go:211) — drop the dead `denied`/`expired` and add the undocumented `inhibited`. [auto] — BLOCKED: gates failed
- [!] issue-91-pairing-health-serve-root: in internal/pairing/http.go, add the unauthenticated `GET /pairing/health` route (no realm segment) sharing the handler/payload of the existing `GET /pairing/v1/{realm}/health`, keeping the v1 route as-is; tests for both paths (closes #91) — BLOCKED: gates failed
- [!] issue-68-async-operation-accepted: on the mutating endpoints upstream surfaces `async_operation` on (housekeeping realm create/delete; realm-management interface install/update/delete, trigger/policy delete in internal/realm/http.go), accept and ignore `?async_operation=false`, with unparseable/`true` values also not changing behaviour (Astrate stays always-sync); tests cover the flag parsing (closes #68) — BLOCKED: gates failed
- [!] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto] — BLOCKED: gates failed
- [!] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
- [!] realm-policy-list-sorted: in `internal/realm/service.go`, sort the names returned by `ListPolicies` the way `ListTriggers`/`ListInterfaces`/`ListInterfaceMajors` already do (it is the only list method without `sort.Strings`), and extend `TestDashboardCompat.Policies` (or a new subtest in `http_test.go`) to install several policies in non-alphabetical order and assert `GET /policies` comes back sorted. [auto] — BLOCKED: gates failed
- [!] realm-interface-lookup-404: in `internal/realm/http_test.go` `TestRealmManagement`, assert the 404 paths that nothing currently covers — `GET /interfaces/<nonexistent-name>` and `GET /interfaces/<name>/<no-such-major>` (service.go `ListInterfaceMajors`:211-213 and `GetInterface`, mapped to 404 via writeError). A rule with no test currently. [auto] — BLOCKED: gates failed
- [!] flow-setstatus-race: in `internal/flow/flow.go` `Manager.StartFlow`, wrap the `f.setStatus(FlowStatusFailed)` call at line 157 in `f.mu.Lock()`/`f.mu.Unlock()` to match the contract documented on `setStatus`. Add a test that starts a flow, confirms it is in the manager's map, then concurrently calls `f.Status()` while `StartFlow` sets the failed status — if the race detector is unavailable, at minimum confirm the lock is acquired. [auto] — BLOCKED: gates failed
- [!] flow-validate-deadcode: in `internal/flow/pipeline.go` `Validate`, remove the dead loop at lines 117-123 (the one whose body is a comment) and the redundant source/sink loop at lines 124-131, keeping only the recomputed inDeg2/outDeg2 check at lines 133-158. Add a test for a pipeline with no source and no sink (already covered) and one for a pipeline where a cycle also lacks sources (to confirm the error message is correct). [auto] — BLOCKED: gates failed
- [!] flow-unmarshal-error-tests: in `internal/flow/message_test.go`, add table-driven tests for `UnmarshalJSON` error paths: missing key, unknown type string, map type field with non-string value, and map data with a field absent from FieldTypes. Each should assert the expected error substring. [auto] — BLOCKED: gates failed
- [!] flow-datawirescalar-fallthrough: in `internal/flow/message.go` `dataWireScalar`, the default case (line 208) returns `m.Data` raw, which is correct after UnmarshalJSON but undocumented. Add a comment documenting the invariant that `Data` matches `Type` after construction, and add a test in `message_test.go` that constructs a FlowMessage with mismatched Type/Data and verifies the wire output is the raw Go value (proving the invariant, not fixing it). [auto] — BLOCKED: gates failed
- [!] docs-sync-appengine-422: in docs/api/astarte_appengine_api.yaml, document the 422 Unprocessable Entity response on the AppEngine operations that reach it — POST /groups (blank group_name, internal/appengine/service.go:292), PATCH /devices/{device}, GET /devices/{device}/interfaces/{interface} (downsample_to on non-datastream, internal/appengine/data.go:90/126/134), GET /devices (invalid cursor, service.go:110) — returned by writeError for ErrValidation (internal/appengine/http.go:318-320); pairing/housekeeping/realm-management yamls already document 422, appengine is the only surface that does not. [auto] — BLOCKED: gates failed
- [!] docs-sync-appengine-socket: in docs/api/astarte_appengine_api.yaml, add the undocumented compat Channels endpoint GET /appengine/v1/socket/websocket (internal/appengine/channels/ws.go:48, Phoenix WebSocket V2 wire format for Dashboard live events) or explicitly note in the yaml why it is excluded, matching how upstream keeps Channels outside the appengine spec. [auto] — BLOCKED: gates failed
- [!] lint-engine-goimports: run goimports -w on `internal/engine/capabilities.go` — its import block is not goimports-formatted (the broker import is not grouped under the local local-prefix), which the golangci goimports formatter flags at every run. [auto] — BLOCKED: gates failed
- [!] lint-flow-findings: clean up the unexcluded revive/gosec findings in `internal/flow` — add `// #nosec G115` to the false-positive rune→uint64 conversion at router.go:215 (a rune can never overflow uint64, mirroring the existing nosec on line 218), delete the empty `if` block at flow_test.go:118, rename the unused `ctx` param of `runLane` to `_` (router.go:180), and add the missing doc comment for the `TypeInteger` constant (message.go:21). [auto] — BLOCKED: gates failed
- [!] lint-pairing-test-param: in `internal/pairing/service_test.go:192`, rename the unused `at` parameter of the `OnRegistered` func to `_` to silence revive's unused-parameter. [auto] — BLOCKED: gates failed
- [!] lint-flow-stutter: in `internal/flow`, rename exported `FlowStatus`/`FlowConfig`/`FlowMessage` to `Status`/`Config`/`Message` to satisfy revive's stutter rule; touches flow.go, message.go, block.go, graph.go, router.go and their tests — a public API rename, so no wire/JSON change (the Type names are unchanged), it is purely Go identifier renaming. [auto] — BLOCKED: gates failed
- [!] serverdata-aftercommit-test: in `internal/engine/serverdata_test.go`, add a test that wires an `afterCommit` observer on a `newWiredRig` engine, calls `PublishServerValue` for a server-owned property, and asserts whether the observer fires. The batch path fires `afterCommit` for every committed op (batch.go:168), but `PublishServerValue` bypasses the batch entirely (serverdata.go:99-162) — the current tests never check this. The test documents the intended behaviour (exclusion or parity with the data path) and fails if the contract changes. [auto] — BLOCKED: gates failed
- [!] consumer-properties-plaintext-fallback: in `internal/engine/control.go` `sendConsumerProperties` (line 204), when `devices.peek` returns nil (device not cached after eviction), the code defaults to zlib compression — but a device that declared `purge_properties_compression_format: plaintext` via the capabilities topic cannot decompress zlib. Fix by falling back to a store read of the device's capability state (or persisting the capability), and add a test in `internal/engine/capabilities_test.go` that evicts the cache entry for a plaintext device and asserts the purge payload is uncompressed. [auto] — BLOCKED: gates failed
- [!] object-path-ok-comment: in `internal/engine/data.go` `objectPathOK` (line 330), the loop over `ci.ObjectLeaves` returns after checking the first entry. Add a comment explaining why one leaf suffices (upstream validation guarantees uniform prefix depth), so future readers do not mistake the early return for a bug. [auto] — BLOCKED: gates failed
- [!] housekeeping-unit-tests: in `internal/housekeeping`, add a unit test file (no integration build tag) that exercises `CreateRealm` validation paths — blank name returns ErrValidation, blank JWT key returns ErrValidation, negative registration limit returns ErrValidation — using a mock store and sealer. The package currently has only integration tests gated behind `//go:build integration`, so these paths are never tested on the Pi. [auto] — BLOCKED: gates failed
- [!] store-downsample-lttb: in `internal/store`, implement the TODO at store.go:135 — when `s.hasToolkit` is true, switch `Downsample` from `time_bucket+avg` to toolkit `lttb()` downsampling. The probe already records availability; the time_bucket path in datastreams.go is the always-working default. Add a test that exercises both code paths (mocking `hasToolkit`). [auto] — BLOCKED: gates failed
- [!] pairing-info-verify-ratelimit: in `internal/pairing/http.go`, add per-IP (and per-device) token-bucket rate limits to `handleInfo` (line 197) and `handleVerify` (line 238), which currently run a full bcrypt `CompareHashAndPassword` (via `authenticateDevice`, service.go:394) on any bearer string with no throttle — the sibling `handleCredentials` that does identical bcrypt work is limited at burst-5, and DESIGN §4.5 says "pairing endpoints" should be rate-limited. Reuse the credentials limiter defaults (1 req/s, burst 5) and extend `TestPairingHTTPRateLimit` to assert a burst of fake-secret requests to `GET /devices/{id}` and `POST .../verify` returns 429. [auto] — BLOCKED: gates failed
- [!] docs-sync-realm-policies-422: in docs/api/astarte_realm_management_api.yaml, add the undocumented 422 Unprocessable Entity response to DELETE /policies/{name} — the code (internal/realm/service.go:398) returns ErrValidation when a policy is still referenced by a trigger, which maps to 422 via writeError; the YAML only lists 204, 401, 404, 500. [auto] — BLOCKED: gates failed

✱ Grep "upstream|astarte-platform|targets|version" in docs/COMPATIBILITY.md · 100 matches
The COMPATIBILITY.md doesn't name a specific upstream target tag in a simple seekable way — but from the todo.md history, the last tracked upstream versions were v1.4.0-rc.5, v1.4.0-rc.3, v1.3.3. The newest tag is now v1.4.0-rc.5 (2026-08-20). Let me check what's been recorded in for-giulio.md and check whether v1.4.0-rc.5 was already the last tracked one (there's a blocked `compat-note-v1.4.0-rc.5` line).

The newest tag is v1.4.0-rc.5. Looking at the todo.md, there's already a blocked line `compat-note-v1.4.0-rc.5` and several blocked probe lines referencing v1.4.0-rc.5. So the last upstream movement tracked was v1.4.0-rc.5. The upcoming changes batch (rc.4, rc.5) were already triaged.

Let me check the for-giulio.md to see if there were already notes, and check the release notes for the newest tag v1.4.0-rc.5 to see if there's new content beyond what was already proposed.

Let me look at the release body for v1.4.0-rc.5 and check for-giulio.md.

$ cat /root/astrate-mule/.mule/for-giulio.md 2>/dev/null | head -100
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **Dependency sweep corrected: direct (pinned) deps DO have newer versions** — the 2026-09-02 note said the `go list -m -u` sweep showed "only version-skew on transitive deps", but that run hit the recipe's `head -20` cutoff (all cloud/azure/transitive) and never reached the directly-required modules. Full sweep, 2026-09-04. None of these is a fix this repo *needs*, so no bump is proposed — recorded for the decision. Per module (current → available; breaking change; repo use):
  - `github.com/coder/websocket` v1.8.14 → v1.8.15 — no breaking (patch); used in `internal/appengine/stream/ws.go`, `channels/ws.go`; worth it only for the "transmit in single frame when compression enabled" fix + read-path alloc reduction.
  - `go.etcd.io/bbolt` v1.4.3 → v1.5.0 — bbolt's semver promises no API change between patch/minor, so additive-only; used in `internal/broker/sessionstore.go`; v1.5 adds a data-file size limit and panic-recovery hardening, nothing Astrate needs.
  - `go.mongodb.org/mongo-driver/v2` v2.6.0 → v2.8.2 — the 2.8.0 breaking changes are confined to Queryable Encryption string-query options (`options.Text()`→`String()`); Astrate uses only the raw BSON API (`pkg/payload/bson.go`, `internal/engine/capabilities.go`, `bench/`) and is unaffected.
  - `github.com/nats-io/nats.go` v1.52.0 → v1.53.1 — no breaking; the headline fixes (JetStream `resetOrderedConsumer` race, KV dot-rejection) are paths Astrate does not use — `internal/engine/forward/nats.go` is core NATS publish only.
  - `github.com/prometheus/client_golang` v1.23.2 → v1.24.1 — requires Go ≥1.25 (fine, repo is 1.26.1); the breaking `LabelNames`/remote-api renames don't touch repo usage (`prometheus`/`collectors`/`promhttp` in `internal/observability/metrics.go`, flow/engine metrics); would buy `Gather()` panic-recovery and opt-in `CoalesceGather` scrape-pile-up protection.
  - `github.com/testcontainers/testcontainers-go` v0.43.0 → v0.44.0 (modules/postgres v0.42.0, modules/nats v0.43.0) — breaking in `wait.ForSQL` (callback now takes `network.Port`) and `ImageProvider` (new `PullImageWithPlatform`); Astrate's `internal/testutil/pg.go` looks unaffected but it is test-only anyway.
  - `golang.org/x/crypto` v0.53.0 → v0.56.0 — x/crypto keeps API compatibility; used only for bcrypt in `internal/auth`.
  Note: `govulncheck` and `golangci-lint` are still not installed on the Pi (existing entries above), so the vulnerability and lint checks remain unavailable.

---

- **milestone 2.0 looks complete, verify and cut the tag** — all 11 `milestone-2.0` issues
  CLOSED (#23–#27, #37, #39–#43), no open issues, no new gaps after re-checking upstream
  astarte_flow block catalog against `internal/flow/` + git log (MQTT/HTTP source/sink,
  json_path_map, pure-transform set, virtual_device_pool, container block MVP, flow API,
  durable named flows all landed). `mule/queue` branch is stale — see entry below.
  (Milestones recipe run, 2026-09-03.)

- **v3.0 decide #68** (`upstream-parity`, `mule-blocked`): `async_operation=false` acceptance is already decided (Astrate stays always-sync; accept and ignore on mutating endpoints: housekeeping realm create/delete, realm-management interface install/update/delete, trigger/policy delete in `internal/realm/http.go`). The mule's code+test task (`issue-68-async-operation-accepted`) failed gates and is parked. Needs re-queue or gate fix.

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

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
Sixth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.4

### Fixed

- [astarte_data_updater_plant] Make object values only validate type on selected interface
- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated

Now I have a clear picture. The newest tag is **v1.4.0-rc.5**. Looking at the todo.md, there are already blocked lines for `probe-object-validation-selected-interface` and `probe-binaryblob-validation` which correspond to the two changes in this rc.5 release notes, plus `compat-note-v1.4.0-rc.5`.

So v1.4.0-rc.5 was already triaged in a previous run — the todo.md already has the corresponding probe lines and compat note. There are no *new* upstream changes since that triage (rc.5 is still the newest tag, and its two changes are already proposed).

Wait, but the recipe says "If nothing changed since last time, append nothing." The question is what was "last time." The newest tag v1.4.0-rc.5 published 2026-08-20 already has its items in the queue (blocked as gates-failed, but present). So there is no *new* upstream movement to triage — the rc.5 changes are already covered by existing queue lines.

Let me verify this reasoning: the existing queue lines are BLOCKED because "gates failed" on the mule/queue branch — these are proposal/known-probe tasks that failed gates, likely because of the branch state issues noted in for-giulio.md. Regardless, the *content* (the astarte-upstream proposals for rc.5) already exists in the queue. Adding them again would duplicate.

So the correct conclusion here: **no new upstream movement since v1.4.0-rc.5** — the rc.5 items are already in the queue. But wait — I should make sure the recipe actually wants me to append anything. The rc.5 release notes changes:
1. "Make object values only validate type on selected interface" → `probe-object-validation-selected-interface` (exists, blocked)
2. "Ensure binaryblob data is correctly validated" → `probe-binaryblob-validation` (exists, blocked)
3. compat note → `compat-note-v1.4.0-rc.5` (exists, blocked)

All three rc.5 items are already in the queue. So there's nothing new to append for upstream movement since rc.5 — the re-run of the proposal job has already covered rc.5.

Therefore the correct output per the recipe's "If nothing changed since last time, append nothing and say 'no upstream movement since <tag>'" — I should say no upstream movement since v1.4.0-rc.5 and append nothing, since the rc.5 proposals already exist in the queue.

Let me double check there's no tag newer than rc.5 I might have missed, and confirm the newest is indeed rc.5.

The first query showed the 5 newest: v1.4.0-rc.5, v1.3.3, v1.4.0-rc.4, v1.4.0-rc.3, v1.4.0-rc.2. So newest is v1.4.0-rc.5.

All the items for rc.5 already exist in the todo.md queue (probe-object-validation-selected-interface, probe-binaryblob-validation, compat-note-v1.4.0-rc.5). Appending them again would duplicate. So the correct action is to append nothing and report no upstream movement since v1.4.0-rc.5.

I'm done. No file writes needed, no git, no task lines appended.
$ gh api repos/astarte-platform/astarte/releases --jq '.[0:1][] | .tag_name'
v1.4.0-rc.5

Newest tag confirmed as **v1.4.0-rc.5** (2026-08-20). Its full set of upstream items is already in the queue:

- `probe-object-validation-selected-interface` ← "Make object values only validate type on selected interface"
- `probe-binaryblob-validation` ← "Ensure binaryblob data is correctly validated"
- `compat-note-v1.4.0-rc.5`

Appending these again would duplicate the previous run's proposals. So there is **no upstream movement since v1.4.0-rc.5**, and I append nothing per the recipe's "empty result is a good result" rule.

```
Done: checked upstream Astarte against newest tag v1.4.0-rc.5; its three items are already proposed in .mule/todo.md, so appended nothing
Files: none
Verified: gh api repos/astarte-platform/astarte/releases -> pass (newest tag v1.4.0-rc.5, published 2026-08-20)
Unsure: nothing
Follow-ups: none
```
