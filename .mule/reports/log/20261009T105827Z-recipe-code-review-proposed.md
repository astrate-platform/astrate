slug: recipe-code-review
verdict: proposed
at:  705c467
ran: 2026-10-09T10:58:27Z on DietPi in 136s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md

$ ls .mule/reviews/ 2>/dev/null; echo "---"; ls internal/ pkg/ 2>/dev/null
appengine-2026-09-07.md
astarteapi-2026-09-15.md
auth-2026-09-14.md
broker-2026-09-04.md
broker-2026-09-21.md
channels-2026-09-18.md
cmd-astrate-2026-09-25.md
config-2026-09-13.md
container-2026-10-02.md
deviceid-2026-09-18.md
engine-2026-09-06.md
engine-triggers-2026-09-20.md
flow-2026-09-05.md
flow-2026-09-23.md
flowapi-2026-09-10.md
forward-2026-09-26.md
housekeeping-2026-09-15.md
httpx-2026-09-17.md
interfaceschema-2026-09-11.md
observability-2026-09-16.md
pairing-2026-09-09.md
pairing-2026-10-07.md
payload-2026-09-12.md
realm-2026-09-09.md
store-2026-09-08.md
store-2026-09-21.md
store-2026-10-04.md
swagger-2026-09-24.md
testutil-2026-10-08.md
---
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

$ tail -60 .mule/log.md 2>/dev/null
| 2026-09-27 | docs-sync-ae-patch-by-alias-409 | done | 119s | d1af059 |
| 2026-09-27 | docs-sync-ae-add-group-device-422 | done | 274s | dff3be8 |
| 2026-09-27 | docs-sync-ae-delete-data-400 | done | 330s | 66bc366 |
| 2026-09-27 | docs-sync-ae-list-devices-422 | transient | 344s |  > build · big-pickle  Error: Error from provider (Console): Rate limit exceeded. Please try again later.  |
| 2026-09-27 | docs-sync-ae-list-devices-422 | transient | 248s |  > build · big-pickle  Error: Error from provider (Console): Rate limit exceeded. Please try again later.  |
| 2026-09-27 | docs-sync-ae-list-devices-422 | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-09-28 | docs-sync-pairing-error-example-capitalisation | done | 166s | 7ed78cb |
| 2026-09-28 | docs-sync-pairing-info-version-example | done | 333s | d9ec98f |
| 2026-09-28 | docs-sync-pairing-deviceid-base64url | done | 480s | 32693cf |
| 2026-09-28 | docs-sync-pairing-initial-payload-format-enum | done | 95s | 1f6a6c2 |
| 2026-09-28 | docs-sync-pairing-unregister-description | done | 261s | 8ffb39e |
| 2026-10-01 | fdo-rc6-scope-delta-for-giulio | done | 479s | 41af3c4 |
| 2026-10-02 | compat-note-v14-rc6 | done | 398s | 0215f97 |
| 2026-10-02 | container-stop-deadline | done | 410s | f9af73b |
| 2026-10-02 | container-timeout-bounds | blocked | 1201s | TIMEOUT after 1201s — task too big, split it |
| 2026-10-02 | container-parseconfig-rules-test | done | 274s | 32c85af |
| 2026-10-02 | container-response-cap-test | blocked | 271s | lint failed: internal/flow/blocks/container/block_test.go:500:2: redefines-builtin-id: redefinition of the built-in function cap (rev |
| 2026-10-02 | docs-sync-rm-error-example-capitalisation | done | 751s | 1d1e6f4 |
| 2026-10-02 | docs-sync-rm-legacy-alias-fields | blocked | 329s | wrote nothing |
| 2026-10-02 | docs-sync-rm-validationerror-example | done | 732s | 8a09a1a |
| 2026-10-02 | docs-sync-rm-deviceid-param | done | 303s | 2d1a1c1 |
| 2026-10-02 | docs-sync-rm-update-interface-body | blocked | 1201s | TIMEOUT after 1201s — task too big, split it |
| 2026-10-03 | docs-sync-ae-error-example-capitalisation | done | 298s | f3bcc6d |
| 2026-10-03 | docs-sync-hk-error-example-capitalisation | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-10-03 | docs-sync-native-error-example-capitalisation | done | 284s | 7cc237f |
| 2026-10-03 | docs-sync-hk-error-detail-examples-split | done | 114s | 23006d4 |
| 2026-10-03 | device-empty-cache-received-trigger | blocked | 202s | wrote nothing |
| 2026-10-03 | appengine-unexpected-object-key | done | 956s | f7d60af |
| 2026-10-03 | appengine-payload-reason-status-map | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-10-03 | appengine-missing-required-422 | done | 584s | f84bbb5 |
| 2026-10-04 | docs-sync-hk-realm-name-pattern | done | 525s | e728217 |
| 2026-10-04 | docs-sync-hk-errordetail-schema-example | done | 184s | 5044143 |
| 2026-10-04 | docs-sync-hk-patch-422-field-error-examples | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-10-04 | docs-sync-hk-wrong-type-field-400 | blocked | 152s | wrote nothing |
| 2026-10-04 | docs-sync-hk-realm-name-response-schemas | done | 178s | def0af6 |
| 2026-10-04 | store-register-inhibit-preserve | blocked | 245s | gates failed |
| 2026-10-04 | store-aliasvalues-self-exclusion-test | done | 717s | a7267ec |
| 2026-10-04 | store-latestindividual-empty-errnotfound | blocked | 1201s | TIMEOUT after 1201s — task too big, split it |
| 2026-10-05 | store-validatepipelinegraph-error-branches | done | 316s | 73e7387 |
| 2026-10-05 | docs-native-ae-realm-version-endpoint | done | 344s | 7331e6c |
| 2026-10-05 | docs-native-hk-version-description | done | 360s | cce1376 |
| 2026-10-05 | docs-native-metrics-content-negotiation | done | 346s | 3962921 |
| 2026-10-05 | docs-native-phoenix-newevent-payload | done | 267s | a16b5fe |
| 2026-10-05 | docs-native-realm-name-pattern | done | 174s | 2ddbf7f |
| 2026-10-06 | probe-encrypted-mapping-enforcement | blocked | 21s | wrote nothing |
| 2026-10-06 | compat-note-v135 | done | 190s | 3a8ea6d |
| 2026-10-06 | docs-native-socket-handshake-errors | blocked | 91s | wrote nothing |
| 2026-10-06 | docs-native-socket-sse-exact-accept | done | 410s | 1322214 |
| 2026-10-06 | docs-native-socket-security-scheme | done | 299s | 87c18cf |
| 2026-10-06 | docs-native-socket-event-schema | blocked | 662s | opencode exited 1 |
| 2026-10-07 | pairing-burn-bcrypt-cost | blocked | 359s | lint failed: internal/pairing/service.go:95:7: G101: Potential hardcoded credentials (gosec) |
| 2026-10-07 | pairing-onregistered-clock | done | 143s | 4e71533 |
| 2026-10-07 | pairing-unregister-bad-id-test | done | 225s | 322242a |
| 2026-10-08 | testutil-wait-since-cursor | done | 319s | 3654fda |
| 2026-10-08 | testutil-control-frame-test | blocked | 102s | wrote nothing |
| 2026-10-08 | engine-introspection-producer-roundtrip | done | 234s | d070345 |
| 2026-10-08 | docs-sync-pairing-register-404-unreachable | done | 148s | cfa6e76 |
| 2026-10-08 | docs-sync-pairing-body-cap-400 | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-10-08 | docs-sync-pairing-realm-name-pattern | done | 410s | ed301fb |
| 2026-10-08 | docs-sync-pairing-deviceid-path-param | done | 395s | 1cda0fb |

$ rg -c '' --glob 'internal/**/*.go' | awk -F: '{n=$1; sub(/\/[^\/]*$/,"",n); c[n]+=$2} END {for (k in c) print c[k], k}' | sort -rn | head -40
9124 internal/engine
6874 internal/store
5044 internal/engine/triggers
4730 internal/appengine
4657 internal/flow/blocks
4576 internal/flow
3757 internal/broker
2730 internal/realm
2533 internal/appengine/channels
2377 internal/pairing
2335 internal/flowapi
1648 internal/swagger
1588 internal/auth
1487 internal/flow/blocks/container
1372 internal/housekeeping
1004 internal/config
984 internal/engine/forward
777 internal/testutil
707 internal/flow/blocks/virtualdevicepool
596 internal/pairing/ca
451 internal/observability
443 internal/engine/stream
369 internal/httpx
324 internal/flow/blocks/astartesource
272 internal/appengine/stream

$ rg -c '' --glob 'pkg/**/*.go' | awk -F: '{n=$1; sub(/\/[^\/]*$/,"",n); c[n]+=$2} END {for (k in c) print c[k], k}' | sort -rn
3407 pkg/interfaceschema
2798 pkg/payload
525 pkg/astarteapi
358 pkg/deviceid

$ grep -l -i 'engine/stream\|appengine/stream\|virtualdevicepool\|astartesource\|pairing/ca\|blocks/' .mule/reviews/*.md
.mule/reviews/auth-2026-09-14.md
.mule/reviews/channels-2026-09-18.md
.mule/reviews/cmd-astrate-2026-09-25.md
.mule/reviews/container-2026-10-02.md
.mule/reviews/engine-2026-09-06.md
.mule/reviews/flow-2026-09-23.md
.mule/reviews/pairing-2026-09-09.md
.mule/reviews/pairing-2026-10-07.md

$ head -30 .mule/reviews/testutil-2026-10-08.md; echo ====; head -25 .mule/reviews/store-2026-10-04.md
# Code review — internal/testutil — 2026-10-08

Rotation check: `container-2026-10-02` recorded that every `internal/` and `pkg/`
package has a review file *except* `internal/testutil`, which it dismissed in
passing as "test helpers, no logic worth proposing on" without reading it. That
dismissal is why this pass exists: the package has no review of its own, it is
imported by 23 files across broker/engine/appengine/store tests plus the
conformance harness, and it encodes three wire rules (control-frame framing,
introspection rendering, golden-file byte comparison) that the default gate
either does or does not assert. It is also small enough to read whole, which is
what I did — 100% of the package, no sampling.

## What I read

- `internal/testutil/astartedevice.go` (all 216) — `AstarteDevice` capture and
  waits (75-175), `Introspection` (93-105), `DeflateControlList` /
  `InflateControlList` (181-216).
- `internal/testutil/mqttclient.go` (all 163) — `ServerTLSCert`, `DeviceCSR`,
  `DeviceTLSConfig`, `MQTTTryConnect`/`MQTTConnect`, `WaitToken`.
- `internal/testutil/golden.go` (all 59) + `golden_test.go` (all 117) — the
  `-update` fixture helper and its three tests.
- `internal/testutil/pg.go` (all 106) + `pg_test.go` (all 36, `integration`) —
  Timescale testcontainer.
- Callers, for reach: `rg` over every `testutil.` use in `internal/`,
  `cmd/`, `test/` (23 files), plus `internal/engine/control.go:291-315`,
  `internal/engine/introspection.go:99-127`,
  `internal/engine/introspection_test.go:17-71`,
  `internal/engine/engine_e2e_test.go:349-380, 423-430`.

## What I found (worth proposing)
====
# Code review: internal/store — 2026-10-04

Area chosen: internal/store. Previous reviews were 2026-09-08 and 2026-09-21.
The package is large (30 source files) and central; I read non-test source files
in descending size order, with particular attention to devices, datastreams,
interfaces, groups/realms/store plumbing, and their corresponding tests.

## What I read

- devices.go (485), devices_test.go (427) — register/unregister/get, inhibit,
  credentials stamp, connect/disconnect/stats, alias values taken, patch merge,
  payload hint, delete
- datastreams.go (434), datastreams_test.go (823) — COPY append, Series,
  LatestIndividual, IndividualSnapshot, ObjectSeries, Downsample, DownsampleLTTB,
  SeriesSpan
- interfaces.go (398), interfaces_test.go (249) — install/update/delete/get/load,
  StoredInterface resolver
- groups.go (265), groups_test.go; realms.go (250), realms_test.go; store.go (204),
  store_test.go; flows.go (207), flows_test.go; pipelines.go (193),
  pipelines_validate_test.go; properties.go (183), properties_test.go;
  crypto.go (136), crypto_test.go; userblocks.go (134), userblocks_test.go;
  notify.go (109), notify_test.go; triggers.go (88); policies.go (88);
  retention.go (54), retention_test.go

## What I found (worth proposing)

$ for f in .mule/reviews/*.md; do echo "== $f"; grep -o '`\(internal\|pkg\)/[a-zA-Z0-9_/]*\.go`' $f | sort -u | tr '\n' ' '; echo; done 2>/dev/null | head -120
== .mule/reviews/appengine-2026-09-07.md
`internal/appengine/data.go` `internal/appengine/downsample.go` `internal/appengine/http.go` `internal/appengine/service.go` 
== .mule/reviews/astarteapi-2026-09-15.md

== .mule/reviews/auth-2026-09-14.md
`internal/auth/claims.go` 
== .mule/reviews/broker-2026-09-04.md
`internal/broker/aclhook.go` `internal/broker/authhook.go` `internal/broker/broker.go` `internal/broker/identity.go` `internal/broker/intake.go` `internal/broker/lifecycle.go` `internal/broker/publisher.go` `internal/broker/sessionstore.go` 
== .mule/reviews/broker-2026-09-21.md

== .mule/reviews/channels-2026-09-18.md
`internal/appengine/channels/protocol.go` `internal/appengine/channels/room.go` `internal/appengine/channels/ws.go` `internal/appengine/stream/ws.go` `internal/engine/stream/bus.go` `internal/engine/triggers/match.go` 
== .mule/reviews/cmd-astrate-2026-09-25.md
`internal/engine/forward/nats.go` 
== .mule/reviews/config-2026-09-13.md
`internal/config/config.go` `internal/config/config_test.go` 
== .mule/reviews/container-2026-10-02.md
`internal/flow/blocks/container/block.go` `internal/flow/blocks/container/block_test.go` `internal/flow/blocks/container/docker.go` `internal/flow/blocks/container/httpbridge.go` 
== .mule/reviews/deviceid-2026-09-18.md
`pkg/deviceid/deviceid.go` `pkg/deviceid/deviceid_test.go` 
== .mule/reviews/engine-2026-09-06.md
`internal/engine/batch.go` `internal/engine/cache.go` `internal/engine/control.go` `internal/engine/data.go` `internal/engine/engine.go` `internal/engine/introspection.go` `internal/engine/router.go` `internal/engine/stream/bus.go` 
== .mule/reviews/engine-triggers-2026-09-20.md
`internal/config/config.go` `internal/engine/triggers/actions.go` `internal/engine/triggers/errorname.go` `internal/engine/triggers/match.go` `internal/engine/triggers/transient.go` `internal/engine/triggers/validation.go` 
== .mule/reviews/flow-2026-09-05.md
`internal/flow/block.go` `internal/flow/composite.go` `internal/flow/factory.go` `internal/flow/flow.go` `internal/flow/graph.go` `internal/flow/message.go` `internal/flow/pipeline.go` `internal/flow/router.go` `internal/flow/substitute.go` 
== .mule/reviews/flow-2026-09-23.md

== .mule/reviews/flowapi-2026-09-10.md
`internal/flow/flow.go` `internal/flowapi/http.go` `internal/flowapi/service.go` `internal/flowapi/userblocks.go` 
== .mule/reviews/forward-2026-09-26.md
`internal/engine/forward/envelope.go` `internal/engine/forward/http.go` `internal/engine/forward/http_test.go` `internal/engine/forward/nats.go` 
== .mule/reviews/housekeeping-2026-09-15.md

== .mule/reviews/httpx-2026-09-17.md

== .mule/reviews/interfaceschema-2026-09-11.md

== .mule/reviews/observability-2026-09-16.md

== .mule/reviews/pairing-2026-09-09.md
`internal/pairing/ca/ca.go` `internal/pairing/http.go` `internal/pairing/http_test.go` `internal/pairing/ratelimit.go` `internal/pairing/service.go` `internal/pairing/service_test.go` 
== .mule/reviews/pairing-2026-10-07.md
`internal/config/config.go` `internal/pairing/ca/ca.go` `internal/pairing/http.go` `internal/pairing/ratelimit.go` `internal/pairing/service.go` 
== .mule/reviews/payload-2026-09-12.md
`pkg/interfaceschema/types.go` `pkg/payload/bson.go` `pkg/payload/json.go` `pkg/payload/payload.go` `pkg/payload/sniff.go` `pkg/payload/value.go` 
== .mule/reviews/realm-2026-09-09.md
`internal/realm/async_operation_test.go` `internal/realm/dashboard_compat_test.go` `internal/realm/detailed.go` `internal/realm/detailed_required_encrypted_test.go` `internal/realm/detailed_test.go` `internal/realm/http.go` `internal/realm/http_test.go` `internal/realm/policy.go` `internal/realm/policy_test.go` `internal/realm/service.go` 
== .mule/reviews/store-2026-09-08.md
`internal/store/datastreams.go` `internal/store/devices.go` `internal/store/interfaces.go` `internal/store/realms.go` `internal/store/store.go` 
== .mule/reviews/store-2026-09-21.md

== .mule/reviews/store-2026-10-04.md

== .mule/reviews/swagger-2026-09-24.md
`internal/swagger/swagger.go` `internal/swagger/swagger_test.go` 
== .mule/reviews/testutil-2026-10-08.md
`internal/testutil/astartedevice.go` `internal/testutil/golden.go` `internal/testutil/mqttclient.go` `internal/testutil/pg.go`

$ ls -la internal/flow/blocks/virtualdevicepool internal/flow/blocks/astartesource internal/engine/stream internal/appengine/stream; echo ---; rg -c '' internal/flow/blocks/virtualdevicepool/*.go internal/flow/blocks/astartesource/*.go internal/engine/stream/*.go internal/appengine/stream/*.go
internal/appengine/stream:
total 20
drwxr-xr-x 2 root root 4096 Jul 27 05:10 .
drwxr-xr-x 4 root root 4096 Oct  3 21:50 ..
-rw-r--r-- 1 root root 4914 Jul 27 05:10 ws.go
-rw-r--r-- 1 root root 3545 Jul 27 05:10 ws_test.go

internal/engine/stream:
total 24
drwxr-xr-x 2 root root 4096 Sep  4 20:46 .
drwxr-xr-x 6 root root 4096 Sep  6 21:58 ..
-rw-r--r-- 1 root root 6103 Sep  4 20:46 bus.go
-rw-r--r-- 1 root root 6210 Jul 27 05:10 bus_test.go

internal/flow/blocks/astartesource:
total 24
drwxr-xr-x 2 root root 4096 Sep  4 20:46 .
drwxr-xr-x 5 root root 4096 Oct  2 19:50 ..
-rw-r--r-- 1 root root 5075 Sep  4 20:46 source.go
-rw-r--r-- 1 root root 4405 Sep  4 20:46 source_test.go

internal/flow/blocks/virtualdevicepool:
total 32
drwxr-xr-x 2 root root  4096 Sep  4 20:46 .
drwxr-xr-x 5 root root  4096 Oct  2 19:50 ..
-rw-r--r-- 1 root root  7369 Sep  4 20:46 virtualdevicepool.go
-rw-r--r-- 1 root root 14944 Sep  4 20:46 virtualdevicepool_test.go
---
internal/engine/stream/bus.go:209
internal/flow/blocks/virtualdevicepool/virtualdevicepool.go:214
internal/engine/stream/bus_test.go:234
internal/flow/blocks/virtualdevicepool/virtualdevicepool_test.go:493
internal/appengine/stream/ws.go:162
internal/flow/blocks/astartesource/source.go:171
internal/flow/blocks/astartesource/source_test.go:153
internal/appengine/stream/ws_test.go:110

$ grep -n 'virtualdevicepool\|astartesource' .mule/todo.md; echo ---; tail -40 .mule/todo.md
---
- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.
- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` shows `detail: Bad Request` (yaml:1732) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1792), while the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58 via `WriteBadRequest`/`WriteInternalServerError`, golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json); every 400/500 on that surface goes through the same astarteapi writers. The other detail examples on the file are already correct — change only those two — and add the docs-half assertion `TestAppEngineErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go, which pins the pairing and realm-management halves. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDetail` schema-level example (yaml:361-362) and the `BadRequest` response (yaml:385) show `detail: Bad Request`, and `InternalServerError` (yaml:477) shows `detail: Internal Server Error`, but the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58, via WriteBadRequest/WriteInternalServerError, pinned by pkg/astarteapi/testdata/error_bad_request.json and error_internal.json). While you are in `ErrorDetail`: its `example` sits at the schema level and shows a bare `{detail: ...}` where the envelope the wire sends is `{"errors":{"detail":...}}` — move it under `errors` (or drop it, as the realm-management and pairing `ErrorDetail` schemas carry no example) so a generated client is not shown the wrong shape. Add `TestHousekeepingErrorDetailExamples` beside the pairing and realm-management halves in internal/swagger/swagger_test.go. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `"Internal server error"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implements `device_empty_cache_received` end to end — `SimpleEvent` field 23, an empty `DeviceEmptyCacheReceivedEvent` whose encoder emits exactly `{"type":"device_empty_cache_received"}` (astarte_core encoder.ex:54-62), fired by DUP at `control_handler.ex:161` *inside* the `with` that runs `send_control_consumer_properties` → `resend_all_properties` → `set_pending_empty_cache`, so it fires only when all three succeed and never on failure, and the commit deletes the `# TODO: implement on_empty_cache_received` from `astarte_events/lib/astarte_events/triggers/core.ex`. Astrate already accepts the trigger but never evaluates it: `OnDeviceEmptyCacheReceived` sits in `deviceOns` with value `false` (internal/engine/triggers/match.go:123, the `evaluated` flag read at match.go:406) and the comment's stated reason ("upstream defines no SimpleEvent variant for it") is now false. Add the event body constructor to internal/engine/triggers/events.go (the payload is `{"type":"device_empty_cache_received"}` and nothing else — no device id, no metadata), flip the `deviceOns` value to true, and fire it from `handleEmptyCache` (internal/engine/control.go:58) with the same idiom `handleLifecycle` uses for `device_connected` (`e.fireDevice` + `stream.Publish`, internal/engine/engine.go:346-353), at the point where `m.Ack()` runs — i.e. after `resendServerProperties` and `sendConsumerProperties` have both succeeded inside `retryStore`. Add a T1 test that a trigger on it fires exactly once per accepted emptyCache and does not fire when the resend fails. — BLOCKED: wrote nothing
- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362) — the second half docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail: Bad Request`, the reconstructed title-case form that `TestHousekeepingErrorDetailExamples` (internal/swagger/swagger_test.go:687) now pins against everywhere else on this spec, while the wire emits the frozen `astarteapi.DetailBadRequest` = `"Bad request"` (pkg/astarteapi/envelope.go:37, via `WriteBadRequest` on every 400 — internal/housekeeping/http.go:65, 166, 171, 183). Because the example sits on the generic `ErrorDetail` schema, which all eight error responses `$ref` (yaml:377-477), it is the example a generator falls back to for any response whose own example it cannot resolve — so it teaches `"Bad Request"` for a 404 and a 500 too. Change it to `Bad request`, or drop it as the realm-management and pairing `ErrorDetail` schemas do (neither carries one). Then extend `TestHousekeepingErrorDetailExamples`, whose comment currently says the schema-level example is "deliberately not pinned" (swagger_test.go:684-686), to pin it like the responses. One thing not to "fix": the example is nested one level in — it is the `errors` sub-schema's example, not a bare envelope — so do not re-wrap it in another `errors:`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-patch-422-field-error-examples [auto]: add the three missing PATCH 422 bodies to the `ValidationErrors` response component in docs/api/astarte_housekeeping_api.yaml (yaml:427-436, referenced once at yaml:190) — it carries a single example, `error_name: [invalid_update_parameters]` (internal/housekeeping/http.go:176-177), but `patchRealm` emits three more field-error envelopes: `jwt_public_key_pem: ["can't be blank"]` (http.go:188-189), `device_registration_limit: ["is invalid"]` (http.go:193-194) and `datastream_maximum_storage_retention: ["is invalid"]` (http.go:198-199). Add them as further `examples` on the media type (OpenAPI 3.0.3, yaml:1) or split the component per case. Also record that the *messages* differ from POST for the same conditions: create rejects a negative limit/retention with the flat `ErrValidation` detail `device_registration_limit must be non-negative` / `datastream_maximum_storage_retention must be non-negative` (internal/housekeeping/service.go:143-148, 422 `ErrorDetail` shape), while PATCH answers `is invalid` in the FieldErrors shape — and that the service's own three `ErrValidation` branches (service.go:205-213) are unreachable from REST, because `patchRealm` pre-checks the identical conditions and answers first. Say in the `RealmPatch` field descriptions that a negative value is a 422 rather than silently ignored. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [!] docs-sync-hk-wrong-type-field-400 [auto]: widen the `BadRequest` component description in docs/api/astarte_housekeeping_api.yaml (yaml:378, "malformed or missing data envelope") — measured today, a perfectly well-formed `data` envelope carrying a wrong-typed field answers **400**, not the documented 422. `astarteapi.DecodeData` into `map[string]json.RawMessage` (internal/housekeeping/http.go:170) cannot fail on a type mismatch, so the failure surfaces on the second decode into the typed body (http.go:182-184) and is answered `WriteBadRequest`. Probe output: `{"data":{"device_registration_limit":"100"}}` → `json: cannot unmarshal string into Go struct field patchBody.device_registration_limit of type int32`; `{"data":{"jwt_public_key_pem":123}}` → the string twin; `{"data":{"device_registration_limit":1.5}}` → `cannot unmarshal number 1.5`; `{"data":[]}` → `cannot unmarshal array into Go value of type patchBody`. Same on POST, which decodes straight into `realmBody` (http.go:63-67). This is unreachable from a spec-conformant client — `RealmCreate`/`RealmPatch` already type both fields — so document it as a one-line extension of the 400 description on both operations rather than inventing a 422, and do NOT re-shape the `BadRequest` example, whose `detail: Bad request` is the frozen constant `TestHousekeepingErrorDetailExamples` pins. One distinction worth stating in the same description: a wrong-typed field is *not* the `invalid_update_parameters` case — unknown keys are caught earlier by the `patchAllowedFields` sweep (http.go:174-180, 422 FieldErrors) on PATCH only, and silently dropped on POST. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `components.parameters.RealmName` and `RealmCreate.realm_name` both carry `pattern: '^[a-z][a-z0-9]*$'` and `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go) holds both against the realm name CHECK extracted from migrations/000002_metadata.up.sql — but the read side is still bare `type: string`: `Realm.realm_name` (components.schemas.Realm) and `DataEnvelopeRealmNames.items`. Decide whether a response field earns the same pattern (a generated client learns nothing new from it, but a reader diffing create against get sees one constrained and one not) and, if it does, add it to both and extend that test to hold the response-side pattern against the same CHECK. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] store-register-inhibit-preserve: RegisterDevice must not overwrite status to 'registered' when a device was inhibited before first credentials request (preserve 'inhibited' state). File: internal/store/devices.go — BLOCKED: gates failed
- [x] store-aliasvalues-self-exclusion-test: Add test for AliasValuesTaken to assert self-exclusion (same device's own alias values do not count as taken). File: internal/store/devices_test.go
- [!] store-latestindividual-empty-errnotfound: Add store-level test asserting LatestIndividual returns ErrNotFound for empty series and returns newest on multiple samples. File: internal/store/datastreams_test.go — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [!] pairing-burn-bcrypt-cost [auto]: make the dummy-secret timing burn use the configured BcryptCost — burnBcrypt compares against a package-level hash generated once at bcrypt.DefaultCost (internal/pairing/service.go:94-100, 428-430) while real auth hashes and compares at cfg.BcryptCost (service.go:204, 407), and bcrypt_cost is a documented, plumbed lever (internal/config/config.go:106, cmd/astrate/main.go:152, docs/site/configuration-reference.md:76, default 10), so with bcrypt_cost < 10 the unknown-device/empty-hash/malformed-ID failure paths each burn cost-10 work against cost-4 real compares, making failure responses ~an order of magnitude slower and a latency oracle for device existence — the inverse of the DESIGN §4.4 uniform-auth-time property. Keep one dummy hash per configured cost (generate lazily per Service at s.cfg.BcryptCost) and burn against it; add a deterministic test asserting bcrypt.Cost of the burn hash equals cfg.BcryptCost (e.g. a Service built with BcryptCost:4 must burn at cost 4), which fails today. — BLOCKED: lint failed: internal/pairing/service.go:95:7: G101: Potential hardcoded credentials (gosec)
- [x] pairing-onregistered-clock: route the registration event timestamp through the injectable clock — Register calls s.OnRegistered(realmName, hwID, time.Now()) with literal wall-clock time (internal/pairing/service.go:222) while every other timestamp in the package funnels through the s.now seam (VerifyCredentials, service.go:333), so a time-travel test cannot pin the DeviceRegisteredEvent timestamp fired via e.HandleDeviceRegistered (cmd/astrate/main.go:154). Change to s.now(), keep TestRegisterEmitsEvent.
- [x] pairing-unregister-bad-id-test: cover the malformed-device-IDs branch of Unregister in internal/pairing/service_test.go — Unregister wraps deviceid.Parse failure as store.ErrNotFound (internal/pairing/service.go:231-234) but TestUnregister (service_test.go:340) only exercises a well-formed unknown ID; assert both bad-ID shapes and the well-formed unknown-ID 404 in one table.
- [x] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.go
- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
- [x] compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (newest stable, 2026-10-05, maintenance-only: mix.exs/mix.lock version bumps plus a single CHANGELOG "astarte_data_updater_plant RPC server reliability" entry, no wire/API surface change; v1.4.0 is still rc.6-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2/v1.3.4 wording proposal there already stands and its version-reference delta was recorded on 2026-10-05 as reading "newest stable v1.3.5", so fold this in as that one-word delta rather than re-deriving the proposal.
- [!] docs-native-socket-handshake-errors [auto]: add the two WebSocket-handshake failures to both socket operations in docs/api/astrate_native_api.yaml — `/astrate/v1/{realm}/socket` responses (yaml:335-345) and `/appengine/v1/socket/websocket` (yaml:398-404) document only 101/401/403/500 and 101/401/500 — but measured today with the repo's own dependency (probe in /tmp, coder/websocket accept.go `verifyClientRequest`): a plain GET with no `Connection: Upgrade` answers **426 Upgrade Required**, `Content-Type: text/plain; charset=utf-8`, body `WebSocket protocol violation: Connection header "" does not contain Upgrade\n`, and `Sec-WebSocket-Version: 14` answers **400** `unsupported WebSocket protocol version (only 13 is supported): "14"\n`. Add both to each operation and say two things explicitly: the bodies are plain text from `http.Error`, NOT the `{"errors":{"detail":...}}` envelope every other response on this spec uses (a client that JSON-parses them breaks), and they are only reached after authentication — native: the route's guard `mw.RequireRealm(auth.ClaimChannels)` wraps `handle` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112), Phoenix: the realm/token checks at internal/appengine/channels/ws.go:51-68 run before `websocket.Accept` at :71 — and only on the WebSocket branch (`wantsSSE` is tested first at ws.go:78), so 401/403 and the 200 SSE path win. Do NOT add 405: ServeMux's method-mismatch applies equally to every GET route in every spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-native-socket-sse-exact-accept [auto]: fix the SSE negotiation prose on `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml — the description (yaml:289-291), the `transport` parameter (yaml:328-334) and the `200` response (yaml:338-339) read as if any client asking for `text/event-stream` gets SSE, but the predicate `wantsSSE` is exact string equality (internal/appengine/stream/ws.go:151-153): measured, `Accept: text/event-stream, text/plain` and `Accept: text/event-stream;q=1.0` both return false and fall through to `websocket.Accept`, where a non-upgrade client gets the 426 above; `transport` selects SSE only for the literal `sse`, and `transport=websocket` / `transport=` (empty) both mean WebSocket. Say "exactly" for the header, say any other `Accept` falls through to the upgrade attempt, and do NOT claim the server rejects out-of-enum `transport` values — the enum is a client-side contract the handler never enforces. While there, give the `200` its real media type: the handler sets `Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive` and writes one `data: {...}` frame per event followed by a blank line (ws.go:123-126, 156-162). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-socket-security-scheme [auto]: `/astrate/v1/{realm}/socket` is documented as if public in docs/api/astrate_native_api.yaml — the spec has no `components.securitySchemes` at all and root `security: []` (yaml:21), and the operation declares no `security`, so a client generated from this spec sends no `Authorization` header and is answered 401 `{"errors":{"detail":"Unauthorized"}}` by `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112, bearer token read only from the Authorization header at middleware.go:131, 178-196). Add a scheme mirroring `a_aea` (docs/api/astarte_appengine_api.yaml:1436-1443: apiKey in header `Authorization`, described as a realm JWT carrying the `a_ch` claim) and set `security` on this operation only — root `security: []` stays so health/readiness/metrics/version remain documented as unauthenticated. If you model the Phoenix twin too, its credential is `?token=` in the query (internal/appengine/channels/ws.go:53-54), i.e. apiKey `in: query`, not header. Description notes worth carrying: the scheme is case-insensitive `Bearer` with an optional colon, and an unknown realm is 401 not 404 (middleware.go:68-72). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-native-socket-event-schema [auto]: document the native socket's event payload in docs/api/astrate_native_api.yaml — currently only a fenced `incoming_data` example inside the operation description (yaml:300-310), with no `components.schemas` entry anywhere (yaml:425-478), and that example is not the general shape: `wireEvent` marks `interface`, `path` and `value` `omitempty` (internal/appengine/stream/ws.go:51-58), so measured on a verbatim copy a lifecycle event marshals to `{"event":"device_connected","realm":"test","device_id":"dev1","timestamp":"2024-05-30T13:49:57.045Z"}` — no interface/path/value at all (and `value: false`/`value: 0` ARE emitted, because omitempty on an `any` drops only nil). Add a `NativeEvent` schema (event, realm, device_id, optional interface/path/value, timestamp) referenced from the operation, state that `event` is one of the internal/engine/stream/bus.go Kind* constants (bus.go:20-26), and record that this socket drops `InterfaceMajor`, `IP`, `ErrorName` and `ErrorMetadata`, which `stream.Event` carries (bus.go:39-65) — so a `device_error` arrives here with no reason, unlike the Phoenix twin, which marshals the whole struct (internal/appengine/channels/ws.go:406-416). That last sentence documents a code gap, not a spec one: do not add the fields; if they are wanted on the wire, that is a separate code task. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1
- [x] testutil-wait-since-cursor [auto]: give AstarteDevice's message waits a cursor so a second wait for the same topic cannot re-match an older message — `WaitForMessage`/`WaitForTopic` (internal/testutil/astartedevice.go:149-170) scan the capture buffer from index 0 on every poll, so a repeated wait returns the first match again, and in-tree that makes testE2EEmptyCache's resync wait vacuous (see the companion line). Move the paho capture handler out of the `ConnectAstarteDevice` closure (astartedevice.go:61-67) into an injectable method, add `Mark() int` and a `From(start int, …)` wait variant, and add a container-free internal/testutil/astartedevice_test.go that drives two captures on one topic and asserts a wait started after the second returns the second message — it cannot even compile today, which is the failure proof. Keep every new read of `received` under the existing `d.mu`, and say in the report that the change is unverified for races (the Pi gate runs no `-race`). The consumer of this API is engine-e2e-emptycache-resync-assert.
- [ ] engine-e2e-emptycache-resync-assert [legion] [auto]: make the second wait in `testE2EEmptyCache` (internal/engine/engine_e2e_test.go:367) wait for the resync re-publish instead of re-matching the retained set captured at :362 — mark the cursor before `dev.EmptyCache(t)` (:364) and take only messages after the mark, so the "resent property value = int32(42)" assertion at :368-369 fails when the resend never happens; today it cannot fail, because `WaitForTopic` returns the pre-EmptyCache message, which carries the same value 42. Needs `-tags e2e` plus Docker, so this half is Legion-only — the container-free helper and its test land in testutil-wait-since-cursor.
- [!] testutil-control-frame-test [auto]: pin the producer/properties control frame that `SendProducerProperties` puts on the wire (internal/testutil/astartedevice.go:181-216) — `DeflateControlList`/`InflateControlList` have zero tests and every caller of the helpers sits behind an `e2e`/`integration`/conformance build tag, so nothing on the default gate asserts the 4-byte big-endian uncompressed-length prefix plus zlib framing, while the parser side is pinned in internal/engine/control_test.go:89-130 against a CPython golden frame. Add a container-free internal/testutil/astartedevice_test.go with exact-bytes, round-trip, empty-list→nil, sub-4-byte header and declared≠inflated cases; bound the inflate with `io.LimitReader(zr, int64(declared)+1)` the way internal/engine/control.go:301-315 does (astartedevice.go:205 reads with unbounded `io.ReadAll` today, and a lying header from the code under test makes the harness allocate it before failing) and assert the truncated byte count in the mismatch message so the bound itself is pinned. One extra line in internal/engine/control_test.go feeding `testutil.DeflateControlList` into `inflateProperties` is the cross-check that closes the loop. — BLOCKED: wrote nothing
- [x] engine-introspection-producer-roundtrip [auto]: pin the producer half of the docs/DESIGN.md §3.3 introspection rule on the default gate — `testutil.Introspection` (internal/testutil/astartedevice.go:93-105) is called from six places, all behind `e2e`, `integration && e2e` or test/conformance build tags, while the consumer `parseIntrospection` (internal/engine/introspection.go:99-127) is tested only against hand-written literals in internal/engine/introspection_test.go:17-71, so nothing container-free checks that the string the harness renders is the one the engine accepts or that the documented "deterministic (sorted) ordering" holds. Add a case to internal/engine/introspection_test.go that builds the string from an unsorted multi-entry map via `testutil.Introspection`, asserts `parseIntrospection` returns the same map, and pins the exact output (sorted, `name:major:minor;…`, no trailing `;`, empty map → empty string).
- [x] docs-sync-pairing-register-404-unreachable [auto]: drop the `404` `DeviceNotFound` response from `registerDevice` in docs/api/astarte_pairing_api.yaml (yaml:140-141) — it is documented but absent: the route is wrapped by `requireAgent` (internal/pairing/http.go:82-83), i.e. `mw.RequireRealm(auth.ClaimPairing)`, which resolves the realm *before* the handler and answers `401 {"errors":{"detail":"Unauthorized"}}` for an unknown one (internal/auth/middleware.go:66-72, "no existence oracle on auth failures"), so `Service.Register`'s own `GetRealmByName` (internal/pairing/service.go:183-186) can return `store.ErrNotFound` only if the realm is deleted between the middleware's lookup and the handler's — a race, not a contract. Nothing else in Register returns `store.ErrNotFound`: `RegisterDevice`'s `ErrDeviceAlreadyConfirmed` becomes 422 (service.go:209-213) and the CountDevices/SetPayloadFormatHint failures are DB errors → 500. No test covers it (internal/pairing/http_test.go asserts 404 only for the agent DELETE at :527 and for health at :595-598), and the `DeviceNotFound` component's own description — "unknown device ID" — does not match what that path would be even if it fired. The premise came from docs-sync-pairing-register-404, which read service.go without the middleware in front of it. Keep the `404` on `unregisterDevice` (reachable: unknown device, and a malformed `deviceID` via service.go:231-234) and on `getHealth`. State the reachability argument in the report; if you find a reachable path instead, keep the response and fix its description rather than deleting it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
- [!] docs-sync-pairing-body-cap-400 [auto]: document the 64 KiB request-body cap on the three body-taking pairing ops — `registerDevice`, `requestCredentials` and `verifyCredentials` all decode through `astarteapi.DecodeData(r.Body, maxBodyBytes, …)` with `maxBodyBytes = 64 << 10` (internal/pairing/http.go:17, 116, 169, 249), and any body over 65536 bytes fails with an error wrapping `ErrBodyTooLarge` (pkg/astarteapi/envelope.go:233-235) that each handler turns into a plain 400 `{"errors":{"detail":"Bad request"}}` (http.go:117-118, 170-171, 250-251). The shared `BadRequest` component (yaml:578-586) says only "malformed or missing data envelope" and no request-body schema carries a size hint, so a client is told nothing about the cap. Widen that one shared description to name the 64 KiB cap and the failure shape (a single edit covers all three operations) and say explicitly that an over-size body is 400, not 413 — do not invent a 413. One clause on `csr` (yaml:414) and `client_crt` (yaml:422) too: the cap is on the whole envelope, not on the PEM. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML). — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-pairing-realm-name-pattern [auto]: constrain `components.parameters.RealmName` in docs/api/astarte_pairing_api.yaml (yaml:359-365, a bare `type: string` described only as "The realm name.") with `pattern: '^[a-z][a-z0-9]*$'` — the realm name is CHECK-constrained to exactly that (migrations/000002_metadata.up.sql:7, read-only, never edit) and the two surfaces whose tasks already ran carry it on their twins (docs/api/astarte_housekeeping_api.yaml:246, docs/api/astrate_native_api.yaml:420 and 467). Because this one parameter is shared by all seven operations, the description must record that an off-pattern name gets three different answers, not one: the agent and device routes answer **401** (the realm is resolved by the middleware, which treats a miss as unauthenticated rather than 404 — internal/auth/middleware.go:66-72); `getHealth` answers **404** (service.Health returns store.ErrNotFound → WriteNotFound, internal/pairing/service.go:369-371, http.go:289-290); and `getVersion` answers **200** `{"data":"1.2.2"}` for any realm string whatsoever, because `GET /pairing/v1/{realm}/version` is registered with no middleware and no realm lookup (cmd/astrate/main.go:438 → observability.VersionHandler, internal/observability/compat.go:38-47) — say that in the version operation so nobody reads the shared parameter as a guarantee the handler enforces. Mirror `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go:1247) with a pairing twin only if it stays a few lines; the YAML edit is the deliverable. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
- [x] docs-sync-pairing-deviceid-path-param [auto]: give `components.parameters.DeviceID` the constraints `RegisterRequest.hw_id` already carries in docs/api/astarte_pairing_api.yaml — the parameter (yaml:367-377) is prose-only, while `hw_id` (yaml:386-398) got `pattern: '^[A-Za-z0-9_-]{22}$'` plus `minLength`/`maxLength` from docs-sync-pairing-deviceid-base64url, and `TestPairingDeviceIDEncodingDocumented` (internal/swagger/swagger_test.go:1086) pins only the two descriptions and hw_id's bounds, never the parameter's schema; the realm-management twin was constrained later by docs-sync-rm-deviceid-param and its test (swagger_test.go:1168) does hold the parameter's pattern against `deviceid.Parse`. Add the same three schema lines to the parameter, and record in its description the asymmetry the prose omits: a malformed `deviceID` is **401** on all three device-authenticated operations — `deviceid.Parse` failing inside `authenticateDevice` burns bcrypt and returns `ErrUnauthorized` (internal/pairing/service.go:381-385) → `WriteUnauthorized` (http.go:322-323), exactly as an unknown device or a wrong secret is (service.go:395-397), so the client is told its credentials are bad rather than its ID malformed — but **404** on the agent `unregisterDevice` (service.go:231-234 → http.go:326-327). Extend TestPairingDeviceIDEncodingDocumented to compile the new pattern and hold it against `deviceid.Parse` the way the realm-management test does. Do not encode the strict trailing-bits rule in the pattern; that stays open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).

→ Read internal/flow/blocks/virtualdevicepool/virtualdevicepool.go

$ rg -n 'IngestFunc|RegisterFunc|ErrVirtualDeviceRegistered|func NewSinkBlock|type Message|type Deps' internal/flow/*.go
internal/flow/block.go:106:func NewSinkBlock(name string, fn SinkFunc) Block {
internal/flow/message.go:81:type Message struct {
internal/flow/factory.go:20:type Deps struct {
internal/flow/factory.go:38:	Ingest IngestFunc
internal/flow/factory.go:43:	Register RegisterFunc
internal/flow/factory.go:46:// ErrVirtualDeviceRegistered reports that a first-seen virtual device could
internal/flow/factory.go:50:var ErrVirtualDeviceRegistered = errors.New("flow: virtual device already registered")
internal/flow/factory.go:52:// RegisterFunc registers a first-seen virtual device through the pairing
internal/flow/factory.go:53:// door (#84 auto_register). Implementations return ErrVirtualDeviceRegistered
internal/flow/factory.go:55:type RegisterFunc func(ctx context.Context, realm, deviceID string) error
internal/flow/factory.go:57:// IngestFunc lands data as if a registered device had produced it: full
internal/flow/factory.go:59:type IngestFunc func(ctx context.Context, realm, deviceID, ifaceName, path string, payload json.RawMessage, ts *time.Time) error

$ rg -n 'virtualdevicepool|virtual_device_pool' --glob '*.go' -l; echo ---; rg -n 'virtual_device_pool' .mule/todo.md .mule/for-giulio.md 2>/dev/null | head
internal/flow/blocks/virtualdevicepool/virtualdevicepool_test.go
internal/flow/blocks/virtualdevicepool/virtualdevicepool.go
internal/flow/blocks/catalog.go
internal/flowapi/service_test.go
internal/flowapi/service.go
internal/flow/factory.go
internal/flow/composite_test.go
---
.mule/for-giulio.md:485:  json_path_map, pure-transform set, virtual_device_pool, container block MVP, flow API,

! permission requested: external_directory (/root/arate-mule/internal/flow/*); auto-rejecting
✗ Read /root/arate-mule/internal/flow/message.go failed
Error: The user rejected permission to use this specific tool call.
