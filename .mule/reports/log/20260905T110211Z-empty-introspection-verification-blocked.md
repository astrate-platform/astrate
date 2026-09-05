slug: empty-introspection-verification
verdict: blocked
at:  8f6dee4
ran: 2026-09-05T11:02:11Z on DietPi in 106s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
✱ Grep "introspection" 100 matches
Task: verify whether Astrate rejects empty introspection strings where upstream v1.3.0 now accepts them, and propose a fix if so.

Let me find the introspection handling code.
✱ Grep "empty.*introspection|introspection.*empty" 11 matches
→ Read internal/engine/introspection.go
→ Read internal/engine/introspection_test.go
→ Read internal/engine/data.go
→ Read internal/engine/mqttv1_parity_test.go
Let me look at the MQTT parity test for empty introspection and also check what upstream does.
→ Read internal/broker/broker_test.go [offset=430, limit=50]
✱ Grep "empty introspection|empty_introspection|1\.3\.0|allow devices with empty" 20 matches
Astrate already handles empty introspection correctly. Let me verify what upstream v1.3.0 actually changed. The task says "upstream v1.3.0 changed 'allow devices with empty introspection'". Let me check whether there's context about this in the repo, and understand what upstream's change entails to verify Astrate's current handling matches.

Let me look at the conformance recordchannels to understand what upstream says, and check if there's an issue or milestone context about v1.3.0.

$ cat /root/astrate-mule/.mule/todo.md 2>/dev/null | tail -30
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

- [x] broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [approved 2026-09-04]
- [!] broker-disconnect-device-zombie-session: in `internal/broker/broker.go` `DisconnectDevice` (lines 260-266), force-closing the live connection leaves the persisted bbolt session (subscriptions + offline queue, keyed by CN) in place — mochi derives `expire` from the Clean flag, not the reason code (server.go:485-491), so `dropClient` (sessionstore.go:500) is never reached and a later device with the same CN resurrects the deleted device's session. Wire the session store handle into `Broker` and call `dropClient` from `DisconnectDevice`, and add a T1 test in `broker_test.go`: connect clean=false + subscribe, call `DisconnectDevice`, reconnect — assert `session_present` is false and the stale subscription is gone. [approved 2026-09-04] — BLOCKED: wrote nothing
- [!] broker-offline-acl-tests: in `internal/broker/aclhook_test.go`, unit-test the offline-delivery ACL — `offlineACL.ownershipOf` cache hit within TTL, TTL expiry triggering a reload, and a load failure caching as empty and denying (aclhook.go:116-138) — plus a T1 `OnACLCheck` check that the offline branch (delivery to a session-present but disconnected device, aclhook.go:196-202) consults the store. The rule currently has zero direct tests; `TestCheckACLMatrix` exercises only the pure `checkACL`. [approved 2026-09-04] — BLOCKED: wrote nothing
- [!] broker-onconnect-doc-comment: in `internal/broker/authhook.go:314`, restore the missing first line of the `OnConnect` doc comment — it opens mid-sentence ("of the Astarte MQTT v1 protocol, and mochi publishes them without a publish-side ACL check…") and the Will-clearing security rationale (a retained LWT would escape the §3.2 matrix) is only half-documented. [approved 2026-09-04] — BLOCKED: tests failed: --- FAIL: TestMQTTSink_Retained (0.02s)

- [x] purge-properties-compression-capability [auto]: upstream v1.3.0 adds a `purge_properties_compression_format` device capability (`zlib`|`plaintext`, default `zlib`) — a wire-visible capability value. Check whether Astrate's capabilities handling (internal/broker, the `<realm>/<device_id>/capabilities` topic, issue #16) needs to recognize/honour it, or whether zlib-only is already the deliberate default; propose the change or note why not needed.
- [ ] empty-introspection-verification [auto]: upstream v1.3.0 changed "allow devices with empty introspection" — verify whether Astrate's device connection/introspection handling currently rejects an empty introspection string where upstream now accepts it, and propose a fix if so.
- [ ] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages" — investigate only: does Astrate have an analogous delay between trigger install and first delivery? Report, do not patch.
- [ ] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.

The task is essentially a verification task. Astrate already handles empty introspection: `parseIntrospection` at introspection.go:99-103 returns `map[string]store.InterfaceVersion{}` for an empty string, and there's already a passing e2e test `TestEmptyIntrospectionAccepted` in mqttv1_parity_test.go.

Let me verify the upstream v1.3.0 change to be thorough. The task asks me to verify whether Astrate rejects empty introspection "where upstream now accepts it". Astrate clearly accepts it (both the parse function and the packed e2e test prove it).

Let me quickly confirm the connection/introspection handling on the device connect side isn't rejecting empty introspection somewhere else. The task specifically mentions "device connection/introspection handling". Let me check if there's any other place that might reject empty introspection, like during the pairing/connection flow.

Let me check the recordchannels conformance to confirm upstream's exact behavior, and check whether there's a connection-time introspection validation.
✱ Grep "Connect|connect|Register|register" in internal/pairing · 100 matches
This is the connection (pairing) side, not the introspection parse side. The introspection arrives on the MQTT introspection topic, and the engine's `handleIntrospection` handles it — that's where empty introspection is handled. There's no connection-time introspection validation that would reject empty introspection.

Let me verify the upstream v1.3.0 change to confirm the exact requirement. Let me search for the upstream change note.
◈ Exa Web Search "Astarte v1.3.0 changelog "empty introspection" allow devices with empty introspection"
The upstream v1.3.0 change "Allow devices with empty introspection" is confirmed. Astrate already accepts an empty introspection payload: `parseIntrospection` (introspection.go:99-103) returns an empty map for `""`, and `TestEmptyIntrospectionAccepted` in mqttv1_parity_test.go proves it end-to-end (acked, persisted, no rejects).

Let me also confirm the data-path gate doesn't misinterpret an empty introspection. When introspection is empty, `dev.declares` will be empty, so any data publish is rejected with `reasonInterfaceNotDeclared` — which is correct upstream behavior (a device with no introspected interface can't publish data).

Let me double check the deviceState `declares` field to make sure an empty introspection doesn't cause a panic or rejection at connect/state-load time.
✱ Grep "declares|func.*declares" in internal/engine · 26 matches
The connection/data-path handling of an empty introspection is correct: `declares` returns not-found for everything, so no data publish is wrongly accepted, and `matchInterface` handles the empty `declares` fine (returns not ok → rejects with `reasonInterfaceNotDeclared`).

Astrate fully supports empty introspection. This is a verification task, and there's nothing to fix. Based on MULE.md rules, a verification with no finding means there's no change to make. Let me run the gate to confirm the existing tests pass and report.

Let me run the relevant engine tests to confirm they pass.

$ go test ./internal/engine/... -run 'Introspection|EmptyIntrospection' -v 2>&1 | tail -40
--- PASS: TestParseIntrospection/non-numeric_major (0.00s)
    --- PASS: TestParseIntrospection/negative_major (0.00s)
    --- PASS: TestParseIntrospection/non-numeric_minor (0.00s)
    --- PASS: TestParseIntrospection/trailing_semicolon (0.00s)
    --- PASS: TestParseIntrospection/wildcard_metacharacter_in_name (0.00s)
    --- PASS: TestParseIntrospection/slash_in_name (0.00s)
=== RUN   TestIntrospectionHandler
--- PASS: TestIntrospectionHandler (0.00s)
=== RUN   TestIntrospectionRejects
--- PASS: TestIntrospectionRejects (0.00s)
=== RUN   TestIntrospectionParking
--- PASS: TestIntrospectionParking (0.10s)
=== RUN   TestEmptyIntrospectionAccepted
--- PASS: TestEmptyIntrospectionAccepted (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/engine	0.148s
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/engine/forward	0.013s [no tests to run]
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/engine/stream	0.018s [no tests to run]
=== RUN   TestCompileDeviceIntrospectionMatrix
=== RUN   TestCompileDeviceIntrospectionMatrix/name-less_added_matches_any_interface_of_the_major
=== RUN   TestCompileDeviceIntrospectionMatrix/name-less_added_rejects_other_majors
=== RUN   TestCompileDeviceIntrospectionMatrix/name-less_removed_matches_any_interface_of_the_major
=== RUN   TestCompileDeviceIntrospectionMatrix/name-less_removed_rejects_other_majors
=== RUN   TestCompileDeviceIntrospectionMatrix/wildcard_name_matches_any_interface
=== RUN   TestCompileDeviceIntrospectionMatrix/wildcard_name_ignores_the_major
=== RUN   TestCompileDeviceIntrospectionMatrix/connect_with_major-only_stays_unfiltered
--- PASS: TestCompileDeviceIntrospectionMatrix (0.00s)
    --- PASS: TestCompileDeviceIntrospectionMatrix/name-less_added_matches_any_interface_of_the_major (0.00s)
    --- PASS: TestCompileDeviceIntrospectionMatrix/name-less_added_rejects_other_majors (0.00s)
    --- PASS: TestCompileDeviceIntrospectionMatrix/name-less_removed_matches_any_interface_of_the_major (0.00s)
    --- PASS: TestCompileDeviceIntrospectionMatrix/name-less_removed_rejects_other_majors (0.00s)
    --- PASS: TestCompileDeviceIntrospectionMatrix/wildcard_name_matches_any_interface (0.00s)
    --- PASS: TestCompileDeviceIntrospectionMatrix/wildcard_name_ignores_the_major (0.00s)
    --- PASS: TestCompileDeviceIntrospectionMatrix/connect_with_major-only_stays_unfiltered (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	0.068s

All relevant tests pass. This is a verification task with **no finding** — Astrate already accepts empty introspection, matching upstream v1.3.0's "Allow devices with empty introspection". There is no fix to propose and no code change to make.

Let me run the full gate to be complete, since the task made no code change but I should confirm the verified state.
$ go test ./internal/engine/ ./internal/broker/ ./internal/store/ 2>&1 | tail -20 && gofmt -l internal/engine/introspection.go internal/engine/data.go internal/engine/introspection_test.go internal/engine/mqttv1_parity_test.go && echo "gofmt clean"
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	2.268s
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
gofmt clean

The verification is complete. Astrate already accepts empty introspection, matching upstream v1.3.0. No fix is needed. I made no code changes.

The task did not name a specific file to change, and there's nothing to fix — this is a `[readonly]`-style verification with a clean (good) outcome. No follow-up tasks to add.

```
Done: Empty introspection is already accepted by Astrate (parseIntrospection returns an empty map for "" and TestEmptyIntrospectionAccepted covers it end-to-end), so no fix is needed — upstream v1.3.0's "allow devices with empty introspection" is already satisfied.
Files: none (verification only)
Verified: go test ./internal/engine/ ./internal/broker/ ./internal/store/ -> pass; gofmt -l (engine introspection/data files) -> clean
Unsure: nothing
Follow-ups: none
```
