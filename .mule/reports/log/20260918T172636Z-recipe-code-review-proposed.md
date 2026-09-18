slug: recipe-code-review
verdict: proposed
at:  216e5a3
ran: 2026-09-18T17:26:36Z on DietPi in 320s


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

$ ls .mule/reviews/ 2>/dev/null; echo "---LOG---"; tail -40 .mule/log.md 2>/dev/null
appengine-2026-09-07.md
astarteapi-2026-09-15.md
auth-2026-09-14.md
broker-2026-09-04.md
config-2026-09-13.md
engine-2026-09-06.md
flow-2026-09-05.md
flowapi-2026-09-10.md
housekeeping-2026-09-15.md
httpx-2026-09-17.md
interfaceschema-2026-09-11.md
observability-2026-09-16.md
pairing-2026-09-09.md
payload-2026-09-12.md
realm-2026-09-09.md
store-2026-09-08.md
---LOG---
| 2026-09-11 | docs-sync-rm-put-auth-422 | done | 145s | 8554720 |
| 2026-09-11 | docs-sync-rm-version-example | blocked | 97s | opencode exited 1 |
| 2026-09-11 | interfaceschema-compat-attrs-flip-test | done | 329s | 4fe1f4f |
| 2026-09-11 | interfaceschema-object-attrs-uniformity | done | 344s | c6c8362 |
| 2026-09-11 | interfaceschema-malformed-placeholder-fixtures | done | 130s | f6a44a2 |
| 2026-09-11 | interfaceschema-enum-roundtrip-coverage | done | 185s | f955958 |
| 2026-09-11 | interfaceschema-trie-interior-literal-param | done | 121s | 705795f |
| 2026-09-12 | docs-sync-ae-write-405 | done | 400s | 9ed08c3 |
| 2026-09-12 | docs-sync-ae-write-value-422 | done | 338s | 17f6f70 |
| 2026-09-12 | docs-sync-ae-post-groups-409 | done | 128s | d6865a1 |
| 2026-09-12 | payload-longinteger-fraction-quantize | done | 356s | d7babc6 |
| 2026-09-12 | payload-json-malformed-t-not-tolerated | done | 99s | 5b56557 |
| 2026-09-12 | docs-sync-ae-post-group-devices-409 | done | 179s | cbf81af |
| 2026-09-12 | docs-sync-ae-patch-merge-patch-media-type | done | 137s | 64361d9 |
| 2026-09-13 | store-todo-lttb-stale | done | 66s | 829f408 |
| 2026-09-13 | config-fail-loud-engine-shards | done | 150s | 83482e4 |
| 2026-09-13 | config-fail-loud-dev-mode | done | 227s | 27d8165 |
| 2026-09-13 | config-key-resolvers-test | done | 143s | 9e68a01 |
| 2026-09-13 | docs-sync-hk-delete-gating-responses | done | 331s | d07b6c7 |
| 2026-09-13 | docs-sync-hk-validation-example | done | 68s | d762f00 |
| 2026-09-14 | auth-iat-not-required-test | done | 257s | 6524c7f |
| 2026-09-14 | auth-empty-claims-403-test | done | 348s | b91247c |
| 2026-09-14 | auth-cache-default-size-test | done | 157s | 874797b |
| 2026-09-15 | hk-zero-retention-asymmetry | done | 904s | 015a706 |
| 2026-09-15 | hk-zero-reglimit-create-asymmetry | blocked | 37s | opencode exited 1 |
| 2026-09-15 | hk-zero-retention-asymmetry | done | 567s | 203af1c |
| 2026-09-15 | docs-sync-native-socket-query-token-auth | done | 98s | 8148aae |
| 2026-09-15 | astarteapi-metadata-envelope-test | done | 238s | 1365cb8 |
| 2026-09-16 | astarteapi-multikey-sorted-golden | done | 215s | bfffbed |
| 2026-09-16 | docs-sync-pairing-deviceendpoints-dead-404-403 | done | 379s | f0d6a4d |
| 2026-09-16 | docs-sync-pairing-version-404-unreachable | done | 200s | bc73f65 |
| 2026-09-16 | observability-readiness-wedge-test | done | 178s | 4ef2e23 |
| 2026-09-16 | observability-dbpool-gauge-coverage | done | 62s | 19cd4f1 |
| 2026-09-16 | observability-health-content-type-test | done | 250s | f6c27f5 |
| 2026-09-17 | docs-sync-pairing-version-value | done | 1059s | c69a164 |
| 2026-09-17 | docs-sync-native-version-value | done | 503s | 08048b3 |
| 2026-09-17 | httpx-cors-vary-origin-passthrough | done | 143s | 986c100 |
| 2026-09-17 | docs-sync-rm-policies-delete-422 | done | 117s | c0521bc |
| 2026-09-17 | docs-sync-rm-triggers-422-nested-envelope | done | 161s | 8bfedd6 |
| 2026-09-18 | compat-note-v134 | done | 154s | bcd4383 |

$ ls .mule/reviews/ | sed 's/-[0-9].*//' | sort
appengine
astarteapi
auth
broker
config
engine
flow
flowapi
housekeeping
httpx
interfaceschema
observability
pairing
payload
realm
store

$ for pkg in internal/*/; do echo "$(find "$pkg" -name '*.go' | wc -l) $pkg"; done | sort -rn
49 internal/flow/
48 internal/engine/
32 internal/store/
23 internal/appengine/
13 internal/broker/
11 internal/realm/
10 internal/pairing/
8 internal/flowapi/
8 internal/auth/
6 internal/testutil/
5 internal/observability/
5 internal/housekeeping/
4 internal/httpx/
2 internal/swagger/
2 internal/config/

$ find pkg/deviceid -name '*.go' | head; echo "---"; rg -c '' pkg/deviceid/*.go 2>/dev/null; echo "---testutil---"; find internal/testutil -name '*.go'; echo "---swagger---"; find internal/swagger -name '*.go'
pkg/deviceid/deviceid.go
pkg/deviceid/deviceid_test.go
---
pkg/deviceid/deviceid_test.go:226
pkg/deviceid/deviceid.go:132
---testutil---
internal/testutil/pg_test.go
internal/testutil/golden.go
internal/testutil/golden_test.go
internal/testutil/pg.go
internal/testutil/mqttclient.go
internal/testutil/astartedevice.go
---swagger---
internal/swagger/swagger_test.go
internal/swagger/swagger.go

**deviceid** (pkg/deviceid) has not been reviewed yet. Reading it now.
→ Read pkg/deviceid/deviceid.go
→ Read pkg/deviceid/deviceid_test.go
Now let me check the callers to understand how the package is used and whether any rule lacks coverage.

$ rg -l 'deviceid\.' --glob '*.go' -g '!pkg/deviceid/*' | sort
cmd/astrate/main.go
cmd/astrate/main_test.go
internal/appengine/dashboard_compat_test.go
internal/appengine/data.go
internal/appengine/e2e_test.go
internal/appengine/formats_parity_test.go
internal/appengine/groups_parity_test.go
internal/appengine/http_test.go
internal/appengine/mirror_test.go
internal/appengine/patch_test.go
internal/appengine/service.go
internal/broker/authhook.go
internal/broker/broker.go
internal/broker/broker_e2e_test.go
internal/broker/broker_test.go
internal/broker/identity.go
internal/broker/identity_test.go
internal/broker/intake.go
internal/engine/cache.go
internal/engine/cache_test.go
internal/engine/control.go
internal/engine/data.go
internal/engine/data_test.go
internal/engine/devicedata_test.go
internal/engine/engine.go
internal/engine/engine_e2e_test.go
internal/engine/engine_integration_test.go
internal/engine/fakestore_test.go
internal/engine/router_test.go
internal/engine/serverdata.go
internal/engine/serverdata_test.go
internal/engine/topics.go
internal/engine/topics_test.go
internal/engine/triggers/match.go
internal/engine/triggers_change_integration_test.go
internal/housekeeping/http_test.go
internal/pairing/http_test.go
internal/pairing/service.go
internal/pairing/service_test.go
internal/realm/dashboard_compat_test.go
internal/realm/http_test.go
internal/realm/service.go
internal/store/datastreams.go
internal/store/datastreams_test.go
internal/store/devices.go
internal/store/devices_test.go
internal/store/groups.go
internal/store/groups_test.go
internal/store/properties.go
internal/store/store_test.go
internal/testutil/astartedevice.go
test/conformance/cpa/cpa_test.go
test/conformance/cpb/cpb_test.go
test/conformance/cpc/cpc_test.go
test/conformance/gosdk/gosdk_test.go
test/conformance/instance/instance.go
test/conformance/load/load_test.go
test/conformance/pysdk/pysdk_test.go
test/security/security_test.go

$ rg -n 'deviceid\.(Parse|FromUUID|Random|FromNamespace|EncodedLen|ErrInvalid)' --glob '*.go' -g '!pkg/deviceid/*' | head -40
cmd/astrate/main.go:162:			id, err := deviceid.Parse(deviceID)
cmd/astrate/main_test.go:159:		id, _ := deviceid.Random()
test/conformance/instance/instance.go:223:	id, err := deviceid.Random()
test/conformance/cpb/cpb_test.go:350:	id, err := deviceid.Random()
test/conformance/cpc/cpc_test.go:294:	id, err := deviceid.Random()
test/conformance/pysdk/pysdk_test.go:60:	id, err := deviceid.Random()
test/security/security_test.go:183:	id, _ := deviceid.Random()
test/security/security_test.go:228:		id, _ := deviceid.Random()
test/security/security_test.go:282:	id, _ := deviceid.Random()
test/conformance/cpa/cpa_test.go:224:	id, err := deviceid.Random()
test/conformance/cpa/cpa_test.go:362:	id, err := deviceid.Random()
test/conformance/gosdk/gosdk_test.go:119:	id, err := deviceid.Random()
internal/store/devices_test.go:152:		ghost, err := deviceid.Random()
internal/store/devices_test.go:216:		ghost, _ := deviceid.Random()
internal/engine/engine.go:382:	id, err := deviceid.Parse(hwID)
internal/engine/engine.go:418:	id, err := deviceid.Parse(hwID)
internal/store/store_test.go:246:	id, err := deviceid.Random()
internal/store/store_test.go:248:		t.Fatalf("deviceid.Random: %v", err)
internal/pairing/service_test.go:175:	id, err := deviceid.Random()
internal/pairing/service_test.go:219:	id, _ := deviceid.Parse(hwID)
internal/pairing/service_test.go:248:	id, _ := deviceid.Parse(hwID)
internal/pairing/service_test.go:303:	id, _ := deviceid.Parse(hwID)
internal/pairing/service_test.go:313:	id2, _ := deviceid.Parse(hwID2)
internal/pairing/service_test.go:353:	id, _ := deviceid.Parse(hwID)
internal/pairing/service_test.go:376:	id, _ := deviceid.Parse(hwID)
internal/pairing/service_test.go:464:	id, _ := deviceid.Parse(hwID)
internal/pairing/service_test.go:491:	id, _ := deviceid.Parse(hwID)
internal/pairing/service_test.go:532:	id, _ := deviceid.Parse(hwID)
internal/store/groups_test.go:59:	ghost, _ := deviceid.Random()
internal/pairing/http_test.go:361:		id, err := deviceid.Parse(hwID)
internal/pairing/http_test.go:616:	id, err := deviceid.Random()
internal/engine/triggers/match.go:493:		if _, err := deviceid.Parse(deviceID); err != nil {
internal/engine/triggers/match.go:507:	if len(s) <= deviceid.EncodedLen {
internal/pairing/service.go:175:	id, err := deviceid.Parse(hwID)
internal/pairing/service.go:231:	id, err := deviceid.Parse(deviceIDStr)
internal/pairing/service.go:381:	id, err := deviceid.Parse(deviceIDStr)
internal/engine/engine_integration_test.go:137:	rig.dev, err = deviceid.Random()
internal/engine/engine_integration_test.go:139:		t.Fatalf("deviceid.Random: %v", err)
internal/engine/triggers_change_integration_test.go:26:	dev, err = deviceid.Random()
internal/broker/broker_test.go:212:	devID, err := deviceid.Random()

→ Read internal/engine/engine.go [offset=370, limit=60]
→ Read internal/engine/triggers/match.go [offset=470, limit=60]

$ rg -n 'isExtendedDeviceID|ExtendedDevice|too long, device id must be 128' internal/engine/triggers/*_test.go | head -20
internal/engine/triggers/match_test.go:386:			want: "device_id=is too long, device id must be 128 bits"},

$ go version; ls pkg/deviceid/
go version go1.26.7 linux/arm64
deviceid.go
deviceid_test.go

$ which elixir erl iex docker 2>&1; echo "---"; git log --oneline -3 -- pkg/deviceid/ 2>/dev/null
---
699cf11 Initial commit under updated license terms

Let me check the current todo.md style and for-giulio format before writing.
$ rg -n '\[legion\]|^\- \[ \]' .mule/todo.md | head -20; echo "---tail---"; tail -12 .mule/todo.md
11:A line tagged `[legion]` needs the Legion Go and is skipped automatically while it is asleep;
47:Labels on the issue are the tags: `legion` and `readonly` mean what `[legion]` and
82:- [ ] broker-external-bus-intake [legion]: implement the second Intake implementation named in the TODO at internal/broker/intake.go:56-62 — an external-bus-backed (e.g. NATS JetStream) Intake for multi-instance deployment / restart survival, reproducing the embedded broker's per-device ordering, deferred-ack backpressure and QoS 0 drop semantics; an implementation of this interface needs containerised integration tests, so it runs on the Legion Go.
84:- [ ] hygiene-govulncheck [legion]: on the Legion Go, run govulncheck ./... on a fresh ~/astrate clone and file a task line for any REACHABLE advisory it reports (name the CVE and the call path, per the hygiene recipe); this box cannot build it — go install golang.org/x/vuln/cmd/govulncheck@latest is OOM-killed here (3.7GB, 4 cores), so the recipe's highest-priority check never runs unless done there.
85:- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.go:560) with no stop signal, so a block-death racing process shutdown (cmd/astrate/main.go:482-488 drains the Manager but never cancels the goroutine) can rebuild the flow via `mgr.StartFlow` on a background pump (internal/flow/flow.go:173-175) after `Manager.Shutdown` and flip the durable status back to running. Thread a stop channel/context through the Service, checked in the loop's sleep, cancelled on shutdown; verification needs a `[legion]` integration or timing-based test. [auto]
89:- [~] bench-small-astrate: run `bench/scripts/run-tier.sh small astrate` against local Astrate, commit results (two runs minimum) [legion]
90:- [~] bench-medium-astrate: run `bench/scripts/run-tier.sh medium astrate` against local Astrate, commit results (two runs minimum) [legion]
91:- [~] bench-small-astarte: run `bench/scripts/run-tier.sh small astarte` against local Astarte, commit results (two runs minimum) [legion]
92:- [~] bench-medium-astarte: run `bench/scripts/run-tier.sh medium astarte` against local Astarte, commit results (two runs minimum) [legion]
93:- [~] bench-big-astrate [legion]: run `bench/scripts/run-tier.sh big astrate` against Legion Go Astrate, commit results (two runs minimum)
94:- [~] bench-giant-astrate [legion]: run `bench/scripts/run-tier.sh giant astrate` against Legion Go Astrate, commit results (two runs minimum)
95:- [~] bench-big-astarte [legion]: run `bench/scripts/run-tier.sh big astarte` against Legion Go Astarte, commit results (two runs minimum)
96:- [~] bench-giant-astarte [legion]: run `bench/scripts/run-tier.sh giant astarte` against Legion Go Astarte, commit results (two runs minimum)
97:- [ ] race-check-store: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/store/... ./internal/housekeeping/... ./migrations/...`. Report any failure to .mule/for-giulio.md with the full race report. Split out of the former single `race-check` line, which timed out running the whole tree at once. [legion] [readonly]
98:- [ ] race-check-engine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/engine/... ./internal/broker/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
99:- [ ] race-check-flow: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/flow/... ./internal/realm/... ./internal/pairing/... ./internal/auth/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
100:- [ ] race-check-appengine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/appengine/... ./internal/observability/... ./internal/httpx/... ./internal/config/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
101:- [ ] race-check-pkg: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./pkg/... ./cmd/... ./internal/testutil/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
---tail---
- [x] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication is "via Authorization header or `?token=` query parameter" (yaml:293-294), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] astarteapi-metadata-envelope-test [auto]: add a container-free golden test for `WriteDataWithMetadata` (pkg/astarteapi/envelope.go:154-158) pinning the two wire branches — nil metadata omits the `metadata` key (`{"data": v}` only, via the `omitempty` at envelope.go:151), populated metadata emits `{"data": v, "metadata": {...}}` — with 200 + the package `Content-Type`, matching the TestGoldenEnvelopes style. Today the non-nil branch is only covered by the integration-tagged formats_parity_test.go (runs on the Legion Go, not the Pi gate) and the nil-omission branch is covered nowhere (internal/appengine/data.go:262-268 always sets metadata), so a regression dropping the omitempty passes every suite on this box.
- [x] astarteapi-multikey-sorted-golden [auto]: extend TestGoldenEnvelopes (pkg/astarteapi/envelope_test.go:101-121) with multi-key `WriteFieldErrors` and `WriteRawErrors` maps asserting the sorted key order on the wire — the docs claim "keys are emitted in sorted order" (envelope.go:170, 184) but every fixture uses a single key, so a change to the marshalling path that stopped sorting (e.g. switching encoders) passes the suite.
- [x] docs-sync-pairing-deviceendpoints-dead-404-403 [auto]: remove the unreachable `404` DeviceNotFound response from the three device-facing ops in docs/api/astarte_pairing_api.yaml — `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials` (yaml:226-227), `GET /pairing/v1/{realm}/devices/{deviceID}` (yaml:268-269), `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify` (yaml:339-340) — and the `403` Forbidden from the info GET (yaml:266-267). `Service.authenticateDevice` folds every failure (malformed/unknown deviceID, unknown realm, unregistered device, wrong secret) into uniform `ErrUnauthorized` (internal/pairing/service.go:380-411), so `store.ErrNotFound` never reaches `writeServiceError` (http.go:326-327) from those handlers; and `Service.Info` never returns `ErrInhibited` (inhibited devices still read info with 200 + status "inhibited", service.go:286-307, comment at 282-285). Reachable on these three are only 200/201/400/401/422/429/500 plus 403 on credentials and verify (both genuinely raise `ErrInhibited`, service.go:252-254 and 329-331). Keep 404 on the agent endpoints (register yaml:149-150, unregister yaml:176-177) where `GetRealmByName`/`deviceid.Parse` really do return `store.ErrNotFound`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-version-404-unreachable [auto]: remove the `404` "Unknown realm" response from `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml (yaml:96-104) — the route is served by `observability.VersionHandler` (cmd/astrate/main.go:398, internal/observability/compat.go:38-47), which answers static `200 {"data":"<version>"}` with no realm lookup, so a 404 can never be emitted there. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] observability-readiness-wedge-test [auto]: make the readiness budget injectable and pinned — add a `timeout time.Duration` field on `Health` (internal/observability/health.go) defaulting to `readinessTimeout` (health.go:12, no test currently exercises the wedged-dependency path), and add a test in observability_test.go asserting a ctx-honoring wedged check makes `/astrate/v1/readiness` return 503 reporting it failing within a short injected budget, instead of hanging the endpoint.
- [x] observability-dbpool-gauge-coverage [auto]: extend `TestMetricsExposesGauges` (internal/observability/observability_test.go:25-34) to assert all four `astrate_db_pool_*` gauges — it currently pins only `acquired_conns` and `max_conns`, so a regression deleting or mis-wiring `idle_conns` or `total_conns` (metrics.go:77-78) passes the suite; expect idle 5, total 7 from the supplied `DBPoolStats{2,5,7,10}`.
- [x] observability-health-content-type-test [auto]: assert `Content-Type: application/json` on `/astrate/v1/health` and `/astrate/v1/readiness` (set by `writeJSON`, health.go:81) and on both the 200 and 503 paths of the per-service `/health` (set by hand in compat.go:26,31) — only the version endpoint pins the header today (compat_version_test.go:31-33), so a regression dropping the header on the /astrate/v1 surface passes the suite.

- [x] docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] httpx-cors-vary-origin-passthrough [auto]: in internal/httpx/cors.go non-wildcard mode, the disallowed-origin pass-through (cors.go:26-29) returns without `Vary: Origin`, so responses whose CORS outcome depends on the Origin header reach shared caches without the RFC 9110 §12.5.5 variation marker (a cache could serve the CORS-approving variant to a disallowed origin, or the headerless variant to an allowed one); the stamp lives only in the allowed branch (cors.go:35-36) and only the preflight pins it (cors_test.go:42). Add `h.Add("Vary", "Origin")` on the pass-through branch too, and assert `Vary: Origin` in TestCORSActualRequestStamped (actual non-preflight, pinned only for Allow-Origin today) and TestCORSDisallowedOriginPassesThroughUnstamped. Additive header change; container-free.
- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ tail -25 .mule/for-giulio.md 2>/dev/null; echo "---count---"; wc -l .mule/for-giulio.md .mule/todo.md 2>/dev/null
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
- **The mule has been idle 17h.** Filed by the dead-man's switch; see journalctl on the Pi.
---count---
   304 .mule/for-giulio.md
   192 .mule/todo.md
   496 total

→ Read .mule/for-giulio.md [limit=30]
Analysis complete. This is a small, unusually well-tested package. My one substantial finding is a suspected wire-parity asymmetry in the strictness of the base64 decoding. Writing the review file, then appending the proposal.

← Write .mule/reviews/deviceid-2026-09-18.md
Wrote file successfully.

$ tail -2 .mule/todo.md
- [x] httpx-cors-vary-origin-passthrough [auto]: in internal/httpx/cors.go non-wildcard mode, the disallowed-origin pass-through (cors.go:26-29) returns without `Vary: Origin`, so responses whose CORS outcome depends on the Origin header reach shared caches without the RFC 9110 §12.5.5 variation marker (a cache could serve the CORS-approving variant to a disallowed origin, or the headerless variant to an allowed one); the stamp lives only in the allowed branch (cors.go:35-36) and only the preflight pins it (cors_test.go:42). Add `h.Add("Vary", "Origin")` on the pass-through branch too, and assert `Vary: Origin` in TestCORSActualRequestStamped (actual non-preflight, pinned only for Allow-Origin today) and TestCORSDisallowedOriginPassesThroughUnstamped. Additive header change; container-free.
- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

→ Read .mule/todo.md [offset=191]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -189,4 +189,5 @@
 
 - [x] docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] httpx-cors-vary-origin-passthrough [auto]: in internal/httpx/cors.go non-wildcard mode, the disallowed-origin pass-through (cors.go:26-29) returns without `Vary: Origin`, so responses whose CORS outcome depends on the Origin header reach shared caches without the RFC 9110 §12.5.5 variation marker (a cache could serve the CORS-approving variant to a disallowed origin, or the headerless variant to an allowed one); the stamp lives only in the allowed branch (cors.go:35-36) and only the preflight pins it (cors_test.go:42). Add `h.Add("Vary", "Origin")` on the pass-through branch too, and assert `Vary: Origin` in TestCORSActualRequestStamped (actual non-preflight, pinned only for Allow-Origin today) and TestCORSDisallowedOriginPassesThroughUnstamped. Additive header change; container-free.
 - [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a locally-run `elixir -e` with `Base.url_decode64!`) with the device ID `"AAAAAAAAAAAAAAAAAAAAAB"` — nonzero unused bits in the 22nd base64url char — and report whether upstream accepts it: `pkg/deviceid/deviceid.go:34-37` uses `base64.RawURLEncoding.Strict()` and deviceid_test.go:176 asserts rejection while the comment claims parity with `Elixir Base.url_decode64!(padding: false)`, which decodes and discards those bits rather than erroring. Astrate may therefore be stricter than upstream on the same wire form (a device registered with such an id on upstream would be rejected here). If upstream rejects it identically, the strict encoding is verified parity and the task is done; if upstream accepts it, escalate the relax-vs-strict call to `.mule/for-giulio.md` with the measurement. Probe first, no code change either way.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **deviceid wire-parity concern, from the 2026-09-18 code review: Astrate is likely stricter than upstream on the wire form's unused bits.** `pkg/deviceid/deviceid.go:37` uses `base64.RawURLEncoding.Strict()`, so `Parse` rejects `"AAAAAAAAAAAAAAAAAAAAAB"` (unused low 4 bits of the 22nd char set; rejection pinned at deviceid_test.go:176), while the comment at deviceid.go:34-36 claims parity with `Elixir Base.url_decode64!(padding: false)` — which decodes the 128 payload bits and discards the unused ones rather than erroring. If so, a device registered with a non-canonical-trailing-bits id on upstream Astarte is accepted there but rejected by Astrate's `Parse` — a wire-compat break in a package whose purpose is the wire form. Not verified on the Pi (no Docker/Elixir); a `[legion]` probe (`deviceid-trailing-bits-upstream-probe`) is queued to measure it. Once confirmed, the call is yours: relax `Parse` to upstream's decode-and-ignore, or keep strict canonicality as the package's contract.
+
+---
+
 - **COMPATIBILITY.md wording update for upstream v1.3.4 (newest stable, 2026-09-18; v1.4.0 is still rc.5-only) — folds v1.3.3 and v1.3.4 into the open v1.3.2 proposal below; both are maintenance-only and wire-inert, so that proposal stands unchanged.** v1.3.3 (2026-08-07) was an empty-body patch and v1.3.4 (2026-09-18, `astarte_data_updater_plant` only — more frequent fullsweep GC on AMQPDataConsumer processes plus an RPC-availability fix; release body via `gh api repos/astarte-platform/astarte/releases` `v1.3.4`) introduce **no wire/API surface change**, so the v1.3.2 wording proposal stands complete: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"). The only delta a v1.3.4-aware doc carries is the version reference: the proposed §Infrastructure-differences sentence's "until the milestone that adopts v1.3.2 as the target" reads "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target". Raw: [v1.3.4](https://github.com/astarte-platform/astarte/releases/tag/v1.3.4).
 
 ---
 

$ git diff --stat; echo "---"; git status --short
.mule/for-giulio.md | 4 ++++
 .mule/todo.md       | 1 +
 2 files changed, 5 insertions(+)
---
 M .mule/for-giulio.md
 M .mule/todo.md
?? .mule/reviews/deviceid-2026-09-18.md

Reviewed `pkg/deviceid` (the one never-reviewed package). Small, well-pinned package — the one real finding is a likely strictness asymmetry vs. upstream Elixir decode, queued as a `[legion]` probe, with the decision escaped to Giulio.

```
Done: reviewed pkg/deviceid; proposed one [legion] probe for a suspected base64 strictness asymmetry vs upstream
Files: .mule/reviews/deviceid-2026-09-18.md, .mule/todo.md, .mule/for-giulio.md
Verified: read deviceid.go + deviceid_test.go + Parse call sites (match.go:488-519, engine.go:377-426) -> no gate needed (no code written)
Unsure: Elixir Base.url_decode64! behaviour not run (no Docker/Elixir on Pi) — hedged in review and for-giulio, probe queued
Follow-ups: deviceid-trailing-bits-upstream-probe [legion]; (plus for-giulio decision note on relax-vs-strict)
```

No code written, no git operations.
