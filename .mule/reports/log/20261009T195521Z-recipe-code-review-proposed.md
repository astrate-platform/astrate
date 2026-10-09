slug: recipe-code-review
verdict: proposed
at:  7642cbc
ran: 2026-10-09T19:55:21Z on DietPi in 310s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md

$ ls /root/astrate-mule/internal/ /root/astrate-mule/pkg/ 2>/dev/null
/root/astrate-mule/internal/:
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

/root/astrate-mule/pkg/:
astarteapi
deviceid
interfaceschema
payload

$ ls -la /root/astrate-mule/.mule/ && echo "---REVIEWS---" && ls -la /root/astrate-mule/.mule/reviews/ 2>/dev/null && echo "---LOG TAIL---" && tail -60 /root/astrate-mule/.mule/log.md 2>/dev/null
total 16452
drwxr-xr-x  7 root root     4096 Oct  9 20:07 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       14 Oct  9 21:50 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  9 19:51 .heartbeat
-rw-r--r--  1 root root        0 Oct  9 21:50 .last-output
-rw-r--r--  1 root root        0 Oct  9 21:50 .lock
-rw-r--r--  1 root root       12 Oct  9 21:50 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 16403404 Oct  9 21:50 cron.log
drwxr-xr-x  2 root root     4096 Oct  9 20:07 failed
-rw-r--r--  1 root root   163554 Oct  9 21:05 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    16266 Oct  9 20:07 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  8 19:28 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   175392 Oct  9 20:07 todo.md
---REVIEWS---
total 236
drwxr-xr-x 2 root root  4096 Oct  8 19:28 .
drwxr-xr-x 7 root root  4096 Oct  9 20:07 ..
-rw-r--r-- 1 root root  4600 Sep  7 20:15 appengine-2026-09-07.md
-rw-r--r-- 1 root root  3443 Sep 15 22:09 astarteapi-2026-09-15.md
-rw-r--r-- 1 root root  5192 Sep 14 13:34 auth-2026-09-14.md
-rw-r--r-- 1 root root  8188 Sep  4 20:43 broker-2026-09-04.md
-rw-r--r-- 1 root root  6580 Sep 21 21:26 broker-2026-09-21.md
-rw-r--r-- 1 root root  6967 Sep 18 22:00 channels-2026-09-18.md
-rw-r--r-- 1 root root 10261 Sep 25 19:46 cmd-astrate-2026-09-25.md
-rw-r--r-- 1 root root  4340 Sep 13 19:26 config-2026-09-13.md
-rw-r--r-- 1 root root 10392 Oct  2 13:43 container-2026-10-02.md
-rw-r--r-- 1 root root  5264 Sep 18 19:25 deviceid-2026-09-18.md
-rw-r--r-- 1 root root  3141 Sep  6 20:57 engine-2026-09-06.md
-rw-r--r-- 1 root root  4692 Sep 20 19:30 engine-triggers-2026-09-20.md
-rw-r--r-- 1 root root  5516 Sep  5 21:19 flow-2026-09-05.md
-rw-r--r-- 1 root root  6931 Sep 23 20:39 flow-2026-09-23.md
-rw-r--r-- 1 root root  4837 Sep 10 19:00 flowapi-2026-09-10.md
-rw-r--r-- 1 root root  8949 Sep 26 19:42 forward-2026-09-26.md
-rw-r--r-- 1 root root  2595 Sep 15 13:29 housekeeping-2026-09-15.md
-rw-r--r-- 1 root root  4914 Sep 17 21:11 httpx-2026-09-17.md
-rw-r--r-- 1 root root  4712 Sep 11 21:40 interfaceschema-2026-09-11.md
-rw-r--r-- 1 root root  4594 Sep 16 20:48 observability-2026-09-16.md
-rw-r--r-- 1 root root  3506 Sep  9 13:07 pairing-2026-09-09.md
-rw-r--r-- 1 root root  4915 Oct  7 19:43 pairing-2026-10-07.md
-rw-r--r-- 1 root root  6002 Sep 12 20:24 payload-2026-09-12.md
-rw-r--r-- 1 root root  4259 Sep  9 21:13 realm-2026-09-09.md
-rw-r--r-- 1 root root  4605 Sep  8 18:59 store-2026-09-08.md
-rw-r--r-- 1 root root  5504 Sep 21 13:32 store-2026-09-21.md
-rw-r--r-- 1 root root  4511 Oct  4 21:08 store-2026-10-04.md
-rw-r--r-- 1 root root  2644 Sep 24 20:37 swagger-2026-09-24.md
-rw-r--r-- 1 root root  8271 Oct  8 19:28 testutil-2026-10-08.md
---LOG TAIL---
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
| 2026-10-09 | docs-sync-ae-data-envelope-metadata | done | 245s | 2ca634d |
| 2026-10-09 | docs-sync-ae-data-write-object-errors | done | 587s | e974103 |
| 2026-10-09 | docs-sync-ae-group-create-devices-required | done | 343s | 9a2cf24 |
| 2026-10-09 | docs-sync-ae-notfound-component | blocked | 220s | wrote nothing |

→ Read .mule/todo.md

$ cd /root/astrate-mule && for f in .mule/reviews/*.md; do echo "== $(basename $f)"; head -3 "$f" | tail -2; done 2>/dev/null
== appengine-2026-09-07.md

Area chosen because recent reviews were `broker` (2026-09-04), `flow`
== astarteapi-2026-09-15.md

Date: 2026-09-15
== auth-2026-09-14.md

## What I read
== broker-2026-09-04.md

## Files read
== broker-2026-09-21.md

Last reviewed: 2026-09-04. Store was reviewed 09-08 and again today, so this
== channels-2026-09-18.md

Area chosen: the `channels/` subpackage (Phoenix V2 wire protocol at
== cmd-astrate-2026-09-25.md

Area: the composition root. Every `internal/*` and `pkg/*` package had been reviewed within
== config-2026-09-13.md

## What I read
== container-2026-10-02.md

Rotation check first: every `internal/` and `pkg/` package has a review file except
== deviceid-2026-09-18.md

Package never reviewed before (checked `.mule/reviews/`; `internal/deviceid`
== engine-2026-09-06.md

Area chosen because `broker` (2026-09-04) and `flow` (2026-09-05) were recent reviews and `engine` is the largest un-reviewed package (48 files).
== engine-triggers-2026-09-20.md

Second pass over `internal/engine`. The 09-06 review covered the core data
== flow-2026-09-05.md

## Files read
== flow-2026-09-23.md

Last reviewed: 2026-09-05 (`flow`). That pass covered the core package;
== flowapi-2026-09-10.md

## What I read
== forward-2026-09-26.md

First review of this area, so no prior file to extend. Chosen because it is the only
== housekeeping-2026-09-15.md

Date: 2026-09-15
== httpx-2026-09-17.md

Date: 2026-09-17
== interfaceschema-2026-09-11.md

Area chosen because `appengine` (09-07), `store` (09-08), `pairing`/`realm`
== observability-2026-09-16.md

Date: 2026-09-16
== pairing-2026-09-09.md

## What I read
== pairing-2026-10-07.md

## What I read
== payload-2026-09-12.md

Area chosen because recent reviews covered `flowapi` (09-10), `interfaceschema`
== realm-2026-09-09.md

## What I read
== store-2026-09-08.md

Area chosen because recent reviews were `broker` (09-04), `flow` (09-05),
== store-2026-09-21.md

Last reviewed: 2026-09-08 (`store`). Rotated back in ahead of engine/flow/appengine.
== store-2026-10-04.md

Area chosen: internal/store. Previous reviews were 2026-09-08 and 2026-09-21.
== swagger-2026-09-24.md

First review of this area, so no prior file to extend.
== testutil-2026-10-08.md

Rotation check: `container-2026-10-02` recorded that every `internal/` and `pkg/`

$ cd /root/astrate-mule && go list ./... 2>/dev/null
github.com/astrate-platform/astrate/cmd/astrate
github.com/astrate-platform/astrate/docs
github.com/astrate-platform/astrate/examples/flow-container-echo
github.com/astrate-platform/astrate/internal/appengine
github.com/astrate-platform/astrate/internal/appengine/channels
github.com/astrate-platform/astrate/internal/appengine/stream
github.com/astrate-platform/astrate/internal/auth
github.com/astrate-platform/astrate/internal/broker
github.com/astrate-platform/astrate/internal/config
github.com/astrate-platform/astrate/internal/engine
github.com/astrate-platform/astrate/internal/engine/forward
github.com/astrate-platform/astrate/internal/engine/stream
github.com/astrate-platform/astrate/internal/engine/triggers
github.com/astrate-platform/astrate/internal/flow
github.com/astrate-platform/astrate/internal/flow/blocks
github.com/astrate-platform/astrate/internal/flow/blocks/astartesource
github.com/astrate-platform/astrate/internal/flow/blocks/container
github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool
github.com/astrate-platform/astrate/internal/flowapi
github.com/astrate-platform/astrate/internal/housekeeping
github.com/astrate-platform/astrate/internal/httpx
github.com/astrate-platform/astrate/internal/observability
github.com/astrate-platform/astrate/internal/pairing
github.com/astrate-platform/astrate/internal/pairing/ca
github.com/astrate-platform/astrate/internal/realm
github.com/astrate-platform/astrate/internal/store
github.com/astrate-platform/astrate/internal/swagger
github.com/astrate-platform/astrate/internal/testutil
github.com/astrate-platform/astrate/migrations
github.com/astrate-platform/astrate/pkg/astarteapi
github.com/astrate-platform/astrate/pkg/deviceid
github.com/astrate-platform/astrate/pkg/interfaceschema
github.com/astrate-platform/astrate/pkg/payload

$ cd /root/astrate-mule && for p in internal/appengine/channels internal/appengine/stream internal/engine/stream internal/engine/forward internal/flow/blocks internal/flow/blocks/astartesource internal/flow/blocks/container internal/flow/blocks/virtualdevicepool internal/pairing/ca; do echo "== $p"; find "$p" -name '*.go' -not -name '*_test.go' | xargs wc -l 2>/dev/null | tail -1; done
== internal/appengine/channels
  912 total
== internal/appengine/stream
162 internal/appengine/stream/ws.go
== internal/engine/stream
209 internal/engine/stream/bus.go
== internal/engine/forward
 276 total
== internal/flow/blocks
  3281 total
== internal/flow/blocks/astartesource
171 internal/flow/blocks/astartesource/source.go
== internal/flow/blocks/container
  721 total
== internal/flow/blocks/virtualdevicepool
214 internal/flow/blocks/virtualdevicepool/virtualdevicepool.go
== internal/pairing/ca
287 internal/pairing/ca/ca.go

$ cd /root/astrate-mule && ls -la internal/pairing/ca internal/flow/blocks/virtualdevicepool internal/flow/blocks/astartesource internal/engine/stream internal/appengine/stream 2>/dev/null
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

internal/pairing/ca:
total 32
drwxr-xr-x 2 root root  4096 Jul 27 05:10 .
drwxr-xr-x 4 root root  4096 Oct  7 20:08 ..
-rw-r--r-- 1 root root 10148 Jul 27 05:10 ca.go
-rw-r--r-- 1 root root  9831 Jul 27 05:10 ca_test.go

$ cd /root/astrate-mule && rg -l 'pairing/ca|internal/pairing/ca|virtualdevicepool|astartesource|appengine/stream|engine/stream' .mule/reviews/ 2>/dev/null; echo "---"; rg -n 'pairing/ca|/ca\b' .mule/reviews/pairing-2026-10-07.md .mule/reviews/pairing-2026-09-09.md 2>/dev/null | head
.mule/reviews/pairing-2026-10-07.md
.mule/reviews/pairing-2026-09-09.md
.mule/reviews/flow-2026-09-23.md
.mule/reviews/engine-2026-09-06.md
.mule/reviews/channels-2026-09-18.md
.mule/reviews/auth-2026-09-14.md
.mule/reviews/cmd-astrate-2026-09-25.md
---
.mule/reviews/pairing-2026-10-07.md:9:- `internal/pairing/ca/ca.go` (full, 287 lines) — `Generate`, `Load`, `SignCSR`, `Verify`,
.mule/reviews/pairing-2026-09-09.md:13:- `internal/pairing/ca/ca.go` + `ca_test.go` — embedded per-realm CA: issue,

→ Read internal/pairing/ca/ca.go
→ Read .mule/reviews/pairing-2026-10-07.md
→ Read .mule/reviews/pairing-2026-09-09.md

$ cd /root/astrate-mule && echo "=== appengine/stream mentions ==="; rg -l 'appengine/stream|internal/appengine/stream' .mule/reviews/ 2>/dev/null; echo "=== engine/stream mentions ==="; rg -l 'engine/stream' .mule/reviews/ 2>/dev/null; echo "=== virtualdevicepool mentions ==="; rg -l 'virtualdevicepool' .mule/reviews/ 2>/dev/null; echo "=== astartesource mentions ==="; rg -l 'astartesource' .mule/reviews/ 2>/dev/null; echo "=== flow/blocks core ==="; rg -l 'flow/blocks' .mule/reviews/ 2>/dev/null
=== appengine/stream mentions ===
.mule/reviews/cmd-astrate-2026-09-25.md
.mule/reviews/auth-2026-09-14.md
.mule/reviews/channels-2026-09-18.md
=== engine/stream mentions ===
.mule/reviews/engine-2026-09-06.md
.mule/reviews/channels-2026-09-18.md
.mule/reviews/cmd-astrate-2026-09-25.md
.mule/reviews/auth-2026-09-14.md
=== virtualdevicepool mentions ===
.mule/reviews/flow-2026-09-23.md
=== astartesource mentions ===
.mule/reviews/flow-2026-09-23.md
=== flow/blocks core ===
.mule/reviews/flow-2026-09-23.md
.mule/reviews/container-2026-10-02.md

→ Read .mule/reviews/flow-2026-09-23.md
→ Read .mule/reviews/channels-2026-09-18.md

$ cd /root/astrate-mule && for p in internal/realm internal/flowapi internal/config internal/auth internal/housekeeping internal/observability internal/httpx internal/interfaceschema pkg/payload internal/deviceid internal/appengine internal/engine; do n=$(find "$p" -name '*.go' -not -name '*_test.go' | xargs wc -l 2>/dev/null | tail -1 | awk '{print $1}'); files=$(find "$p" -name '*.go' -not -name '*_test.go' | wc -l); echo "$p: $n lines, $files files"; done
internal/realm: 1348 lines, 4 files
internal/flowapi: 1511 lines, 3 files
internal/config: 420 lines, 1 files
internal/auth: 620 lines, 4 files
internal/housekeeping: 547 lines, 2 files
internal/observability: 216 lines, 3 files
internal/httpx: 108 lines, 2 files
internal/interfaceschema: 0 lines, 0 files
pkg/payload: 1406 lines, 5 files
internal/deviceid: 0 lines, 0 files
internal/appengine: 3211 lines, 8 files
find: ‘internal/interfaceschema’: No such file or directory
find: ‘internal/interfaceschema’: No such file or directory
find: ‘internal/deviceid’: No such file or directory
find: ‘internal/deviceid’: No such file or directory
internal/engine: 6247 lines, 21 files

$ cd /root/astrate-mule && ls -la internal/realm/ && wc -l internal/realm/*.go
total 132
drwxr-xr-x  2 root root  4096 Sep  9 21:46 .
drwxr-xr-x 17 root root  4096 Sep  4 20:46 ..
-rw-r--r--  1 root root  4344 Sep 25 22:06 async_operation_test.go
-rw-r--r--  1 root root  8885 Sep  4 20:46 dashboard_compat_test.go
-rw-r--r--  1 root root  2674 Sep  8 18:56 detailed.go
-rw-r--r--  1 root root  1315 Sep  8 18:54 detailed_required_encrypted_test.go
-rw-r--r--  1 root root  2666 Sep  4 20:46 detailed_test.go
-rw-r--r--  1 root root  3765 Sep  9 21:46 helpers_test.go
-rw-r--r--  1 root root 19956 Sep  9 21:49 http.go
-rw-r--r--  1 root root 35006 Sep  4 20:46 http_test.go
-rw-r--r--  1 root root   427 Jul 27 05:10 policy.go
-rw-r--r--  1 root root  3593 Jul 27 05:10 policy_test.go
-rw-r--r--  1 root root 24167 Sep  9 21:49 service.go
    89 internal/realm/async_operation_test.go
   213 internal/realm/dashboard_compat_test.go
    76 internal/realm/detailed.go
    29 internal/realm/detailed_required_encrypted_test.go
    63 internal/realm/detailed_test.go
   136 internal/realm/helpers_test.go
   569 internal/realm/http.go
   782 internal/realm/http_test.go
    17 internal/realm/policy.go
    70 internal/realm/policy_test.go
   686 internal/realm/service.go
  2730 total

→ Read .mule/reviews/realm-2026-09-09.md
I'll review `internal/realm` — last reviewed 2026-09-09, and its follow-ups are already done. Let me read it.
→ Read internal/realm/service.go
→ Read internal/realm/http.go
→ Read internal/realm/helpers_test.go
→ Read internal/realm/policy.go
→ Read internal/realm/policy_test.go

$ cd /root/astrate-mule && rg -n 'renderViolationsBody|validationDetail|writeInterfaceError|writeViolations|quoteJSON|marshalStrings|OnDeletionStart|OnDeletionFinish|ListPolicies|GetDatastreamMaximumStorageRetention|SetAuthKey' internal/realm/*_test.go internal/realm/*.go | grep -v 'func ' | head -50
internal/realm/dashboard_compat_test.go:92:		r.svc.OnDeletionStart = func(realmName, deviceID string, _ time.Time) {
internal/realm/dashboard_compat_test.go:95:		r.svc.OnDeletionFinish = func(realmName, deviceID string, _ time.Time) {
internal/realm/service.go:83:	// OnDeletionStart / OnDeletionFinish bookend DeleteDevice (nil-safe).
internal/realm/service.go:84:	OnDeletionStart  OnDeletionFunc
internal/realm/service.go:85:	OnDeletionFinish OnDeletionFunc
internal/realm/service.go:126:	if s.OnDeletionStart != nil {
internal/realm/service.go:127:		s.OnDeletionStart(realm, deviceID, at)
internal/realm/service.go:130:	if s.OnDeletionFinish != nil {
internal/realm/service.go:131:		s.OnDeletionFinish(realm, deviceID, time.Now().UTC())
internal/realm/service.go:540:// ListPolicies returns the realm's policy names.
internal/realm/service.go:600:// GetDatastreamMaximumStorageRetention returns the realm's datastream maximum
internal/realm/service.go:630:// SetAuthKey rotates the realm's JWT public key (upstream PUT /config/auth).
internal/realm/dashboard_compat_test.go:92:		r.svc.OnDeletionStart = func(realmName, deviceID string, _ time.Time) {
internal/realm/dashboard_compat_test.go:95:		r.svc.OnDeletionFinish = func(realmName, deviceID string, _ time.Time) {
internal/realm/detailed.go:26:	b.Write(quoteJSON(iface.Name))
internal/realm/detailed.go:29:	b.Write(quoteJSON(iface.Type.String()))
internal/realm/detailed.go:31:	b.Write(quoteJSON(iface.Ownership.String()))
internal/realm/detailed.go:33:	b.Write(quoteJSON(iface.Aggregation.String()))
internal/realm/detailed.go:36:		b.Write(quoteJSON(iface.Description))
internal/realm/detailed.go:40:		b.Write(quoteJSON(iface.Doc))
internal/realm/detailed.go:58:	b.Write(quoteJSON(m.Endpoint))
internal/realm/detailed.go:60:	b.Write(quoteJSON(m.Type.String()))
internal/realm/detailed.go:66:	b.Write(quoteJSON(m.Reliability.String()))
internal/realm/detailed.go:68:	b.Write(quoteJSON(m.Retention.String()))
internal/realm/detailed.go:71:	b.Write(quoteJSON(m.DatabaseRetentionPolicy.String()))
internal/realm/http.go:63:	names, err := a.svc.ListPolicies(r.Context(), r.PathValue("realm"))
internal/realm/http.go:134:// seconds (upstream 1.2.0+; see Service.GetDatastreamMaximumStorageRetention).
internal/realm/http.go:136:	retention, err := a.svc.GetDatastreamMaximumStorageRetention(r.Context(), r.PathValue("realm"))
internal/realm/http.go:153:			a.writeInterfaceError(w, err)
internal/realm/http.go:161:		a.writeInterfaceError(w, err)
internal/realm/http.go:175:		a.writeInterfaceError(w, err)
internal/realm/http.go:184:		a.writeInterfaceError(w, err)
internal/realm/http.go:197:		a.writeInterfaceError(w, err)
internal/realm/http.go:214:		a.writeInterfaceError(w, err)
internal/realm/http.go:226:		a.writeInterfaceError(w, err)
internal/realm/http.go:327:	if err := a.svc.SetAuthKey(r.Context(), r.PathValue("realm"), cfg.JWTPublicKeyPEM); err != nil {
internal/realm/http.go:347:// writeInterfaceError maps interface-service/store errors onto the upstream
internal/realm/http.go:359:		writeViolations(w, ve)
internal/realm/http.go:393:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
internal/realm/http.go:409:		writeViolations(w, ve)
internal/realm/http.go:411:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
internal/realm/http.go:427:// writeViolations renders a *ViolationsError as upstream's Phoenix-changeset
internal/realm/http.go:435:	body := renderViolationsBody(ve)
internal/realm/http.go:441:// renderViolationsBody encodes the {"errors": {...}} object. Interface-level
internal/realm/http.go:452:		buf.Write(quoteJSON(key))
internal/realm/http.go:477:		piece(f.field, marshalStrings(f.msgs))
internal/realm/http.go:520:					arr.Write(quoteJSON(f))
internal/realm/http.go:522:					arr.Write(marshalStrings(e.msgs[f]))
internal/realm/http.go:535:// quoteJSON encodes s as a JSON string without HTML escaping, matching the
internal/realm/http.go:546:// marshalStrings renders msgs as a JSON array of strings.

$ cd /root/astrate-mule && rg -n 'go:build|package realm' internal/realm/*_test.go | head -40
internal/realm/helpers_test.go:1:package realm
internal/realm/async_operation_test.go:1://go:build integration
internal/realm/async_operation_test.go:3:package realm
internal/realm/dashboard_compat_test.go:1://go:build integration
internal/realm/dashboard_compat_test.go:3:package realm
internal/realm/http_test.go:1://go:build integration
internal/realm/http_test.go:3:package realm
internal/realm/detailed_required_encrypted_test.go:1://go:build !integration
internal/realm/detailed_required_encrypted_test.go:3:package realm
internal/realm/policy_test.go:1:package realm
internal/realm/detailed_test.go:1://go:build integration
internal/realm/detailed_test.go:3:package realm

$ cd /root/astrate-mule && rg -n 'func .*ListTriggerPolicies|func .*ListTriggers|func .*LoadRealmInterfaces' internal/store/*.go
internal/store/triggers.go:67:func (s *Store) ListTriggers(ctx context.Context, realmID int16) ([]Trigger, error) {
internal/store/interfaces.go:266:func (s *Store) LoadRealmInterfaces(ctx context.Context, realmID int16) ([]*StoredInterface, error) {
internal/store/policies.go:67:func (s *Store) ListTriggerPolicies(ctx context.Context, realmID int16) ([]TriggerPolicy, error) {

$ cd /root/astrate-mule && rg -n 'ListTriggerPolicies' -A 18 internal/store/*.go | head -60
internal/store/policies_test.go:56:	list, err := s.ListTriggerPolicies(ctx, realm.ID)
internal/store/policies_test.go-57-	if err != nil {
internal/store/policies_test.go:58:		t.Fatalf("ListTriggerPolicies: %v", err)
internal/store/policies_test.go-59-	}
internal/store/policies_test.go-60-	if len(list) != 2 || list[0].Name != "discard" || list[1].Name != "retry5xx" {
internal/store/policies_test.go-61-		t.Errorf("list = %+v, want [discard retry5xx]", list)
internal/store/policies_test.go-62-	}
internal/store/policies_test.go-63-
internal/store/policies_test.go-64-	if err := s.DeleteTriggerPolicy(ctx, realm.ID, "retry5xx"); err != nil {
internal/store/policies_test.go-65-		t.Fatalf("DeleteTriggerPolicy: %v", err)
internal/store/policies_test.go-66-	}
internal/store/policies_test.go-67-	if err := s.DeleteTriggerPolicy(ctx, realm.ID, "retry5xx"); !errors.Is(err, ErrNotFound) {
internal/store/policies_test.go-68-		t.Errorf("double delete = %v, want ErrNotFound", err)
internal/store/policies_test.go-69-	}
internal/store/policies_test.go-70-}
--
internal/store/policies.go:66:// ListTriggerPolicies returns every policy of a realm ordered by name.
internal/store/policies.go:67:func (s *Store) ListTriggerPolicies(ctx context.Context, realmID int16) ([]TriggerPolicy, error) {
internal/store/policies.go-68-	rows, err := s.pool.Query(ctx,
internal/store/policies.go-69-		`SELECT id, realm_id, name, definition FROM trigger_policies WHERE realm_id = $1 ORDER BY name`,
internal/store/policies.go-70-		realmID)
internal/store/policies.go-71-	if err != nil {
internal/store/policies.go-72-		return nil, fmt.Errorf("store: listing policies: %w", err)
internal/store/policies.go-73-	}
internal/store/policies.go-74-	defer rows.Close()
internal/store/policies.go-75-
internal/store/policies.go-76-	var out []TriggerPolicy
internal/store/policies.go-77-	for rows.Next() {
internal/store/policies.go-78-		var p TriggerPolicy
internal/store/policies.go-79-		if err := rows.Scan(&p.ID, &p.RealmID, &p.Name, &p.Definition); err != nil {
internal/store/policies.go-80-			return nil, fmt.Errorf("store: scanning policy: %w", err)
internal/store/policies.go-81-		}
internal/store/policies.go-82-		out = append(out, p)
internal/store/policies.go-83-	}
internal/store/policies.go-84-	if err := rows.Err(); err != nil {
internal/store/policies.go-85-		return nil, fmt.Errorf("store: listing policies: %w", err)

$ cd /root/astrate-mule && rg -n 'type ViolationsError|type Violation struct|MappingIndex|MappingCount|Field ' pkg/interfaceschema/*.go | head -30
pkg/interfaceschema/violations.go:8:// Violation is one upstream-shaped field rejection. MappingIndex < 0 marks
pkg/interfaceschema/violations.go:9:// an interface-level field; >= 0 marks a field of mappings[MappingIndex].
pkg/interfaceschema/violations.go:10:type Violation struct {
pkg/interfaceschema/violations.go:11:	Field        string
pkg/interfaceschema/violations.go:13:	MappingIndex int
pkg/interfaceschema/violations.go:21:type ViolationsError struct {
pkg/interfaceschema/violations.go:23:	// MappingCount is the number of declared mappings, needed to render
pkg/interfaceschema/violations.go:25:	MappingCount int
pkg/interfaceschema/violations.go:42:	e.Violations = append(e.Violations, Violation{Field: field, Messages: msgs, MappingIndex: -1})
pkg/interfaceschema/violations.go:47:	e.Violations = append(e.Violations, Violation{Field: field, Messages: msgs, MappingIndex: index})
pkg/interfaceschema/violations.go:74:		if v.MappingIndex >= 0 {
pkg/interfaceschema/violations.go:75:			fmt.Fprintf(&b, "mappings[%d].", v.MappingIndex)
pkg/interfaceschema/parse.go:295:	vc.MappingCount = len(raw.Mappings)
pkg/interfaceschema/violations_test.go:44:		if v.Field == field && v.MappingIndex == index {
pkg/interfaceschema/violations_test.go:221:		if ve.MappingCount != 2 {
pkg/interfaceschema/violations_test.go:222:			t.Errorf("MappingCount = %d, want 2 (declared mappings)", ve.MappingCount)

$ cd /root/astrate-mule && rg -n 'ViolationsError|Violation{' pkg/interfaceschema/violations.go 2>/dev/null; ls pkg/interfaceschema/
compat.go
compat_test.go
compile.go
compile_test.go
parse.go
parse_test.go
required_encrypted_test.go
testdata
trie.go
trie_test.go
types.go
violations.go
violations_test.go

$ cd /root/astrate-mule && git log --oneline --since=2026-09-09 -- internal/realm/ 2>/dev/null | head -40
f965f39 mule: docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cd /root/astrate-mule && rg -n 'realm-|realm:' .mule/todo.md | head -40
152:- [x] realm-pure-helper-tests [auto]: add container-free unit tests for the four untested pure helpers in internal/realm. Create `internal/realm/helpers_test.go` (no build tag) covering: (1) `normaliseIfaceName` — lowercase, hyphen strip, mixed case+hyphens, empty string, already-normalized passthrough; (2) `joinPEM` — empty slice, single key, two keys separated by `\n`; (3) `triggerErrorBody` — empty action+simple_triggers, action-only, simple_triggers-only with index-aligned `{}` substitution, both present; (4) `majorParam` — valid positive, zero, negative (→404), non-numeric (→404). All tests use httptest.ResponseRecorder for majorParam; the rest are direct function calls. No Docker.
157:- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
194:- [x] docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
196:- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
218:- [x] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "realm: validation failed: interface definition is invalid"`, but `validationDetail` strips the `realm: validation failed: ` prefix on the wire (http.go:562-568) and the parser's real messages are `invalid interface: ...` (pkg/interfaceschema/parse.go:38, wrapped at service.go:163); the example should show the stripped, real message. Same class as the already-fixed docs-sync-hk-validation-example (housekeeping http.go:251-258). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
228:- [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
232:- [x] forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
236:- [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml — the string `"403"` appears nowhere in the file and `components` defines no `Forbidden` response, yet every route is wrapped by `a.require = mw.RequireRealm(auth.ClaimAppEngine)` (internal/appengine/http.go:35), which answers `403 Forbidden` for a verified a_aea JWT whose grants do not authorize the method+path (`astarteapi.WriteForbidden`, internal/auth/middleware.go:90; `DetailForbidden`, pkg/astarteapi/envelope.go:41; pinned today by internal/appengine/http_test.go:184-185). Add a `Forbidden` response component copied from the realm-management one (docs/api/astarte_realm_management_api.yaml:1277-1285 — `ErrorDetail` schema, example `errors.detail: Forbidden`) and `$ref` it from all 32 operations, which is exactly the parity docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go:123-170) pins; add the same docs-half test for the AppEngine spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
254:- [x] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, housekeeping and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
256:- [x] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
260:- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` shows `detail: Bad Request` (yaml:1732) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1792), while the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58 via `WriteBadRequest`/`WriteInternalServerError`, golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json); every 400/500 on that surface goes through the same astarteapi writers. The other detail examples on the file are already correct — change only those two — and add the docs-half assertion `TestAppEngineErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go, which pins the pairing and realm-management halves. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
261:- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDetail` schema-level example (yaml:361-362) and the `BadRequest` response (yaml:385) show `detail: Bad Request`, and `InternalServerError` (yaml:477) shows `detail: Internal Server Error`, but the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58, via WriteBadRequest/WriteInternalServerError, pinned by pkg/astarteapi/testdata/error_bad_request.json and error_internal.json). While you are in `ErrorDetail`: its `example` sits at the schema level and shows a bare `{detail: ...}` where the envelope the wire sends is `{"errors":{"detail":...}}` — move it under `errors` (or drop it, as the realm-management and pairing `ErrorDetail` schemas carry no example) so a generated client is not shown the wrong shape. Add `TestHousekeepingErrorDetailExamples` beside the pairing and realm-management halves in internal/swagger/swagger_test.go. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
268:- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
269:- [x] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362) — the second half docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail: Bad Request`, the reconstructed title-case form that `TestHousekeepingErrorDetailExamples` (internal/swagger/swagger_test.go:687) now pins against everywhere else on this spec, while the wire emits the frozen `astarteapi.DetailBadRequest` = `"Bad request"` (pkg/astarteapi/envelope.go:37, via `WriteBadRequest` on every 400 — internal/housekeeping/http.go:65, 166, 171, 183). Because the example sits on the generic `ErrorDetail` schema, which all eight error responses `$ref` (yaml:377-477), it is the example a generator falls back to for any response whose own example it cannot resolve — so it teaches `"Bad Request"` for a 404 and a 500 too. Change it to `Bad request`, or drop it as the realm-management and pairing `ErrorDetail` schemas do (neither carries one). Then extend `TestHousekeepingErrorDetailExamples`, whose comment currently says the schema-level example is "deliberately not pinned" (swagger_test.go:684-686), to pin it like the responses. One thing not to "fix": the example is nested one level in — it is the `errors` sub-schema's example, not a bare envelope — so do not re-wrap it in another `errors:`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
272:- [x] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `components.parameters.RealmName` and `RealmCreate.realm_name` both carry `pattern: '^[a-z][a-z0-9]*$'` and `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go) holds both against the realm name CHECK extracted from migrations/000002_metadata.up.sql — but the read side is still bare `type: string`: `Realm.realm_name` (components.schemas.Realm) and `DataEnvelopeRealmNames.items`. Decide whether a response field earns the same pattern (a generated client learns nothing new from it, but a reader diffing create against get sees one constrained and one not) and, if it does, add it to both and extend that test to hold the response-side pattern against the same CHECK. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
280:- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
281:- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
284:- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
297:- [x] docs-sync-pairing-realm-name-pattern [auto]: constrain `components.parameters.RealmName` in docs/api/astarte_pairing_api.yaml (yaml:359-365, a bare `type: string` described only as "The realm name.") with `pattern: '^[a-z][a-z0-9]*$'` — the realm name is CHECK-constrained to exactly that (migrations/000002_metadata.up.sql:7, read-only, never edit) and the two surfaces whose tasks already ran carry it on their twins (docs/api/astarte_housekeeping_api.yaml:246, docs/api/astrate_native_api.yaml:420 and 467). Because this one parameter is shared by all seven operations, the description must record that an off-pattern name gets three different answers, not one: the agent and device routes answer **401** (the realm is resolved by the middleware, which treats a miss as unauthenticated rather than 404 — internal/auth/middleware.go:66-72); `getHealth` answers **404** (service.Health returns store.ErrNotFound → WriteNotFound, internal/pairing/service.go:369-371, http.go:289-290); and `getVersion` answers **200** `{"data":"1.2.2"}` for any realm string whatsoever, because `GET /pairing/v1/{realm}/version` is registered with no middleware and no realm lookup (cmd/astrate/main.go:438 → observability.VersionHandler, internal/observability/compat.go:38-47) — say that in the version operation so nobody reads the shared parameter as a guarantee the handler enforces. Mirror `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go:1247) with a pairing twin only if it stays a few lines; the YAML edit is the deliverable. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
298:- [x] docs-sync-pairing-deviceid-path-param [auto]: give `components.parameters.DeviceID` the constraints `RegisterRequest.hw_id` already carries in docs/api/astarte_pairing_api.yaml — the parameter (yaml:367-377) is prose-only, while `hw_id` (yaml:386-398) got `pattern: '^[A-Za-z0-9_-]{22}$'` plus `minLength`/`maxLength` from docs-sync-pairing-deviceid-base64url, and `TestPairingDeviceIDEncodingDocumented` (internal/swagger/swagger_test.go:1086) pins only the two descriptions and hw_id's bounds, never the parameter's schema; the realm-management twin was constrained later by docs-sync-rm-deviceid-param and its test (swagger_test.go:1168) does hold the parameter's pattern against `deviceid.Parse`. Add the same three schema lines to the parameter, and record in its description the asymmetry the prose omits: a malformed `deviceID` is **401** on all three device-authenticated operations — `deviceid.Parse` failing inside `authenticateDevice` burns bcrypt and returns `ErrUnauthorized` (internal/pairing/service.go:381-385) → `WriteUnauthorized` (http.go:322-323), exactly as an unknown device or a wrong secret is (service.go:395-397), so the client is told its credentials are bad rather than its ID malformed — but **404** on the agent `unregisterDevice` (service.go:231-234 → http.go:326-327). Extend TestPairingDeviceIDEncodingDocumented to compile the new pattern and hold it against `deviceid.Parse` the way the realm-management test does. Do not encode the strict trailing-bits rule in the pattern; that stays open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).

→ Read internal/realm/detailed.go
→ Read internal/realm/detailed_required_encrypted_test.go

$ cd /root/astrate-mule && rg -n 'internal/realm|realm\.' .mule/for-giulio.md | tail -20
36:- **COMPATIBILITY.md wording proposal for upstream v1.3.5 (newest stable, 2026-10-05; v1.4.0 is still rc.6-only) — a one-word delta on the open v1.3.2/v1.3.4 proposal below, nothing re-derived.** v1.3.5 is maintenance-only: `gh api repos/astarte-platform/astarte/releases` shows it published `2026-10-05T14:11:44Z`, not a prerelease, with a body of a single `### Fixed` line — "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available" — and `compare/v1.3.4...v1.3.5` is 4 commits / 35 files whose only non-`mix.exs`/`mix.lock` files are `.github/workflows/astarte-build-workflow.yaml`, `CHANGELOG.md`, `README.md`, `astarte-dashboard/package.json` and `docker-compose.yml`. Not one source file changed, so no route, MQTT topic, control message or interface-schema field moved (Astrate implements no `data_updater_plant` and no Erlang RPC server). **The only change a v1.3.5-aware doc carries: the v1.3.4 entry's §Infrastructure-differences sentence — "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target" — reads "…(**newest stable v1.3.5**) as the target".** Everything else stands unchanged: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"); that entry's "v1.4.0 is still rc.5-only" header is already superseded to rc.6 by the open rc.6 delta entry below, not here. This is the same one-word delta recorded on 2026-10-05 in this file ("newest stable v1.3.5"), folded in as asked rather than re-derived. Applied to nothing: `docs/COMPATIBILITY.md` untouched and still 1.2.2-targeted (it contains no `1.3` reference at all), `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). Raw: [v1.3.5](https://github.com/astarte-platform/astarte/releases/tag/v1.3.5).
50:- **Milestone recipe run, 2026-10-06: v3.0 verified unchanged — still gated on upstream — and the master delta since the 2026-10-05 run is a closed set of 13 commits, every one wire-inert for Astrate. No issue filed, no `.mule/todo.md` line queued.** Release sweep is byte-identical to yesterday's (`gh api repos/astarte-platform/astarte/releases`): **no stable v1.4.0** — newest stable `v1.3.5` (2026-10-05), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition (your 2026-09-04 "wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) stays unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** — this milestone's whole ledger lives under `upstream-parity`, where the open set is unchanged: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`); **#47–#89** otherwise all closed. Non-alarm open issues total exactly four (plus the 25 `mule-alarm` #94–#118): the three above and **#1** (untouched per standing instruction). `APICompatVersion` still `"1.2.2"` (`internal/realm/service.go:588`, re-checked).
57:  **The one delta to your open COMPATIBILITY.md wording proposal (the v1.3.2 one, folding v1.3.3/v1.3.4, further down this file): its version-reference sentence should read "newest stable v1.3.5" where it currently says "v1.3.4". Nothing else in it changes** — v1.3.5 adds no capability, so the same decision and the same "not yet emulated" list stand. Not applied: `docs/COMPATIBILITY.md` is untouched and still targets **1.2.2**, `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). I deliberately did **not** queue a fresh `compat-note-v1.3.5` line: the proposal is already open here, and a second line would re-derive it for the next run to find.
68:  **#92's parking condition is still unmet, now for the twenty-first run.** Your 2026-09-04 decision ("wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) — newest stable is `v1.3.5` (today), newest overall remains `v1.4.0-rc.6` (2026-09-30, prerelease). The post-rc.6 master delta has grown by exactly one commit: `731bcddc` "chore: forward port release-1.4" (2026-10-05T09:34Z), and it is the housekeeping/FDO migration relocation the 2026-10-03 run already measured inside rc.5→rc.6 (six `astarte_housekeeping/priv/migrations/{astarte,realm}/*.sql` files removed, table creation moved into migrations) plus a DUP-internal `apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0011_remove_replacement_data.ex` (+31) and `database.ex` (+2/-1). No Realm Management, AppEngine, Pairing or Trigger wire surface. **The milestone ledger is unchanged:** `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** (this milestone's whole ledger lives under `upstream-parity`, as the 2026-09-10 entry established), of #47–#89 **all remain closed**, and the authoritative open set is 29 issues — 25 alarms (#94–#118, one new today: **#118**, that recipe's business) and exactly four non-alarm: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** (untouched per standing instruction). `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`, re-checked), per **#90's** frozen decision.
93:  **Why this is worth your attention rather than a note.** The register is wired into the bump rule by its own header: "It is bumped only in the same change that completes the full surface of that upstream level. When that happens, **every open row below tagged with that level must be reconciled first** — the register is the checklist" (`docs/UPSTREAM-EXPERIMENTAL.md:16-18`), and `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`). So these two mis-tagged rows are standing between v3.0's final phase and the bump, and their only listed action — "upstream 1.4 final: promoted or removed" — is a condition upstream can never satisfy in the literal sense, because there is nothing experimental to promote or remove. It *is* dischargeable at 1.4 final under a generous reading ("present and unmarked in the final CHANGELOG" = promoted, so per the header "keep and drop the row"), but that reading is mine, not the file's, and the register does not say it. Fifteen prior runs of this recipe checked *whether 1.4 is final* and *whether the feature still exists*; none checked *whether upstream calls it experimental*, which is why this sat. I have not edited `docs/UPSTREAM-EXPERIMENTAL.md`, `docs/COMPATIBILITY.md` or anything else.
111:  **The recipe's standing item is satisfied and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental*, reconcile trigger "upstream 1.4 final". 1.4 is not final, so there is nothing to promote and nothing to deprecate, and the register is accurate as written. `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked this run) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision ("bumped only in the same change that completes that surface, after reconciling every row tagged with the level"). The milestone-bump issue the recipe asks me to file **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": v3.0's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is delivered, the 1.4 half is open (#92, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut.
120:- **Delta on the two open COMPATIBILITY.md wording entries further down this file (the v1.3.4 one folding in v1.3.3, and the original v1.3.2 one): upstream `v1.4.0-rc.6` is wire-inert for Astrate, so nothing in that proposal is re-derived — only its version reference and one line of its "not yet emulated" list change.** Not applied: `docs/COMPATIBILITY.md` is untouched (still the 1.2.2 target, deviation 18's "which upstream 404s" at `docs/COMPATIBILITY.md:360`, `APICompatVersion = "1.2.2"` at `internal/realm/service.go:588`), and the FDO half of the rc.6 delta already has its own entry below rather than being folded in here.
168:  **The recipe's standing item is satisfied, and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still has exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental* with reconcile trigger "upstream 1.4 final". v1.4 is not final, so both rows stay as written: nothing to promote, nothing to deprecate, and the register is accurate. `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision; the milestone-bump issue the recipe asks for **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": the milestone's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is done, and the 1.4 half is unresolved (#92 open, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut. Also still true and already escalated on 2026-09-06, so not repeated: the section's source document `.mule/research/upstream-parity-2026-08-22.md` does not exist in the repo.
309:- **docs-sync realm-management, 2026-09-17: `docs/site/realm-management-api.md` documents an "Update trigger" endpoint Astrate does not serve.** Lines 81-86 show `PUT /realmmanagement/v1/<realm>/triggers/<name>` with a `{ "data": <updated trigger JSON> }` body, but no such route is registered — internal/realm/http.go:44-47 wires only GET/POST `/triggers` and GET/DELETE `/triggers/{name}` (an update op is deliberately absent; not even a stub or TODO exists). Site prose is yours — drop the "Update trigger" section or confirm the route is planned.
450:- **COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable, 2026-07-14; v1.4.0 is still rc-only).** Astrate's doc and `APICompatVersion` still target upstream **1.2.2** (`internal/realm/service.go:588`); v1.3.0 (2026-05-06) introduced wire-surface changes Astrate does not yet emulate, so this is a decision — adopt v1.3.2 as the compatibility target (then update the doc + bump `APICompatVersion` together, per the bump rule) or keep 1.2.2 and add a "not yet emulated" note. Wire-relevant v1.3.0 deltas (release notes): **MQTT v1 capabilities** incl. `purge_properties_compression_format` (plaintext vs zlib purge — touches the `emptyCache`/`producer,properties` contract COMPATIBILITY.md deviation 1 documents); **empty introspection now allowed**; **device registration triggers** (pairing) and **device deletion started/completed triggers** (RM — the latter two already exist as Astrate deviation 9 emits both around the synchronous delete); **FDO authentication** (pairing, disabled by default); **realm-scoped health** — upstream v1.3 added `GET /pairing/v1/{realm}/health`, which Astrate already serves (`internal/pairing/http.go:78-81`, comment already says "upstream 1.3+"), so deviation 18's wording ("which upstream 404s") is now false against 1.3 and the note should be reworded either way. Proposed doc wording (for your approval, edit to taste): in §Infrastructure differences add a sentence — *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl.
662:the realm-management surface joins *all* keys (`internal/realm/service.go:627`) — but every
664:`internal/store/realms.go:165`, `putAuthConfig` at `internal/realm/service.go:636`), so a

$ cd /root/astrate-mule && go test ./internal/realm/ 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)

→ Read .mule/for-giulio.md [offset=640, limit=50]
→ Read pkg/interfaceschema/violations.go

$ cd /root/astrate-mule && rg -n 'writeInterfaceError|realm-error|realm-detailed|renderViolationsBody|allow_unset' .mule/todo.md .mule/tasks/ 2>/dev/null | head
.mule/todo.md:158:- [x] docs-sync-rm-delete-interface-status [auto]: fix the `deleteInterface` response in docs/api/astarte_realm_management_api.yaml — the spec documents `422` with examples "Interface major version is not 0, can't be deleted" / "Cannot delete an interface that is used by a device introspection" (yaml:255-272), but the handler returns `403` with "Interface can't be deleted" / "Interface can't be deleted since it's currently used" (internal/realm/http.go:387-391, `writeInterfaceError`; both are 403 `StatusForbidden`). The spec's detail strings were lifted from the `writeError` branch (http.go:414-419) which DELETE /interfaces does not use. Fix status to `403` and update the two examples to match the actual wire messages. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/todo.md:162:- [x] interfaceschema-compat-attrs-flip-test [auto]: extend TestCheckMinorUpgradeTable (pkg/interfaceschema/compat_test.go) to assert the four currently-unpinned immutable attributes — retention, expiry, database_retention_policy/ttl (datastream old/next pair) and allow_unset (properties old/next pair) — each flipped on upgrade and rejected as ErrIncompatibleEndpointChange with an at-immutability acceptance twin; only type/reliability/explicit_timestamp/required/encrypted are tested today (sameMappingAttributes, compat.go:102, sits at 66.7% stmt coverage), so a regression relaxing expiry or DB-TTL immutability passes the suite.
.mule/todo.md:217:- [x] docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` (yaml:231-232) in docs/api/astarte_realm_management_api.yaml as a oneOf — both handlers go through `writeInterfaceError`, which answers three distinct 422 bodies: the flat ErrorDetail (`ErrValidation` via `validationDetail`, http.go:392-393), the nested violations changeset envelope (`writeViolations`, http.go:358-359 + 434-533, e.g. `{"errors":{"description":["should be at most 1000 character(s)"]}}` and the aligned full-length `mappings` array, http_test.go:480-501), and the named FieldErrors envelope for `ErrMaximumDatabaseRetentionExceeded` (`{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, http.go:355-357, http_test.go:282-320); today only the flat ValidationError is referenced. Mirror the createTrigger oneOf pattern (yaml:344-368). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cd /root/astrate-mule && wc -l .mule/todo.md && tail -25 .mule/todo.md
302 .mule/todo.md
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
- [x] docs-sync-ae-data-envelope-metadata [auto]: add the `metadata` property to `DataEnvelopeObject` in docs/api/astarte_appengine_api.yaml (yaml:1711-1717) — the schema requires only `data` and types it `object`, but every `format=table` data response carries a sibling `metadata` object. Code: `astarteapi.WriteDataWithMetadata` renders `{"data": v, "metadata": m}` through `metadataEnvelope` (`Metadata map[string]any json:"metadata,omitempty"`, pkg/astarteapi/envelope.go:149-158), the `Tabular` carrier sets it (internal/appengine/data.go:59-64), and both renderers populate it — `renderIndividual` (data.go:256-268, `columns` keyed on the path's last segment + `table_header`) and `renderObject` (data.go:293-321, `columns` + sorted `table_header`). Add `metadata` as an optional object property (`columns: map<string,int>`, `table_header: array of string`) described as present only when `format=table`, on all six data GET 200s that `$ref` this schema (yaml:416, 468, 685, 730, 1244, 1290). While there, correct `data`'s type/description: it is not always an object — `format=structured` individual is `[{value,timestamp}]` (data.go:275-281) and `format=disjoint_tables` is an object keyed by path (data.go:269-274, 322-331). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-data-write-object-errors [auto]: document the two object-aggregated write rejections on the six data PUT/POST ops (putData/publishData, by-alias, in-group) in docs/api/astarte_appengine_api.yaml — each refs `BadRequest` for 400 (yaml:518, 570, 775, 822, 1336, 1384) and `ValueTooLarge` for 422 (yaml:528, 580, 785, 832, 1346, 1394), but three bodies are missing. (a) 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` from `writeBadObjectError` (internal/appengine/http.go:722-748, const `detailUnexpectedObjectKey` :711), reached by the `payload.ReasonOf(err) == payload.ReasonBadObject` case (http.go:684-691) — so the shared `ErrorDetail` schema (yaml:1741-1750) needs an optional `unexpected_keys` array, or a dedicated response. (b) 422 `{"errors":{"detail":"Missing required mapping key"}}` (const `detailMissingRequiredMapping`, http.go:713-720) from the `ReasonMissingRequired` case (http.go:692-702) — add an example beside the value-size one already on those ops. Both were implemented by appengine-unexpected-object-key/appengine-missing-required-422 (upstream master b6d46ad4, #2237) with no spec update. Do NOT invent other `payload.ReasonOf` statuses — that family is the separate blocked line appengine-payload-reason-status-map. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-group-create-devices-required [auto]: in docs/api/astarte_appengine_api.yaml, `GroupCreate` (yaml:1652-1663) marks only `group_name` required, but `createGroup` treats `devices` as required — `groupBody.Devices *[]string` (internal/appengine/http.go:391-397) makes a MISSING key a 422 `{"errors":{"devices":["can't be blank"]}}` and an EMPTY array a 422 `{"errors":{"devices":["should have at least 1 item(s)"]}}` (http.go:405-418), and only the MISSING-key case is unreachable from the schema as written. The shared `ValidationErrors` example (yaml:1841-1844) that POST /groups refs (yaml:948) shows `group_name: can't be blank` alone. Add `devices` to `GroupCreate`'s `required` list (the 201 echo always carries it, http.go:424) and add the two device messages to the POST 422 story. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-ae-notfound-component [auto]: the `NotFound` response component in docs/api/astarte_appengine_api.yaml (yaml:1795-1803, example `errors.detail: Not Found`) is `$ref`d by all twelve group-route 404s (yaml:978, 1026, 1070, 1104, 1152, 1179, 1209, 1250, 1296, 1342, 1390, 1429), yet `astarteapi.WriteNotFound`/`DetailNotFound` ("Not Found", pkg/astarteapi/envelope.go:42-43, 205-207) is emitted nowhere in internal/appengine (rg finds only `WriteDeviceNotFound`, http.go:683). A missing group answers 404 `{"errors":{"detail":"Group not found"}}` (`ErrGroupNotFound`, internal/appengine/service.go:41-42, 508-516; writeError http.go:667-668) and a missing member device answers `"Device not found"` (writeError store.ErrNotFound → WriteDeviceNotFound, http.go:682-683). Repoint/rename the component so the group routes show `Group not found` (add a `GroupNotFound` component, keep `DeviceNotFound` on the routes that can 404 for a missing member) and state which routes can emit which. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing

$ cd /root/astrate-mule && git log --pretty='%ad %h %s' --date=short -- internal/realm/*.go | head -20
2026-09-25 f965f39 mule: docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
2026-09-09 39eb70a mule: realm-pure-helper-tests [auto]: add container-free unit tests for the four untested pure helpers in internal/realm. Create `internal/realm/helpers_test.go` (no build tag) covering: (1) `normaliseIfaceName` — lowercase, hyphen strip, mixed case+hyphens, empty string, already-normalized passthrough; (2) `joinPEM` — empty slice, single key, two keys separated by `\n`; (3) `triggerErrorBody` — empty action+simple_triggers, action-only, simple_triggers-only with index-aligned `{}` substitution, both present; (4) `majorParam` — valid positive, zero, negative (→404), non-numeric (→404). All tests use httptest.ResponseRecorder for majorParam; the rest are direct function calls. No Docker.
2026-09-08 63b629e mule: detailed-listing-required-encrypted-flags: emit the `required`/`encrypted` mapping flags in `writeDetailedMapping` (internal/realm/detailed.go) so `GET /realmmanagement/v1/{realm}/interfaces?detailed=true` round-trips them — the detailed listing currently materializes a fixed field set that drops both flags even when a datastream mapping declares one, while the plain GET preserves them via the stored definition. Datastream branch only (properties reject both at parse); add a T1 assertion in internal/realm/detailed_test.go.
2026-09-04 b6e1ca4 compat: make the dashboard health indicators real, pin always-sync (#68, #91)
2026-08-23 dbd0267 docs: upstream-experimental register + APICompatVersion bump rule (#67, #68, #90)
2026-08-23 ed0342d feat(triggers): upstream wildcard semantics + introspection interface rules (#70)
2026-08-22 edf4e1b feat(triggers): upstream-measured HTTP action validation limits (#63)
2026-08-22 029e0b6 feat(triggers): pairwise-disjoint policy handlers, prefetch_count bounds, AMQP action rejection (issues #65, #64)
2026-08-22 0b5a35c feat(realm,interfaceschema): upstream-measured install/update/delete error taxonomy (issue #62)
2026-08-22 c073835 feat(realm,interfaceschema): validation matrix gaps + detailed interface listing (issues #61, #66)
2026-08-22 834d536 feat(realm,store): datastream_maximum_storage_retention ceiling — install-time reject + retention sweep (issue #72)
2026-08-22 5d74ad8 feat(realm): GET config/datastream_maximum_storage_retention (#60)
2026-07-29 aba2d18 fix: land blocked mule issues #14 #21 #22 #24 #27
2026-07-26 699cf11 Initial commit under updated license terms

$ cd /root/astrate-mule && tail -c 400 .mule/todo.md | od -c | tail -5
0000520   S   w   a   g   g   e   r   -   r   e   f   e   r   e   n   c
0000540   e   d       Y   A   M   L       s   t   i   l   l       l   o
0000560   a   d   s   .     342 200 224       B   L   O   C   K   E   D
0000600   :       w   r   o   t   e       n   o   t   h   i   n   g  \n
0000620

Now I'll write the review file and append the task lines.
← Write .mule/reviews/realm-2026-10-09.md
Wrote file successfully.

→ Read .mule/todo.md [offset=300]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -299,4 +299,6 @@
 - [x] docs-sync-ae-data-envelope-metadata [auto]: add the `metadata` property to `DataEnvelopeObject` in docs/api/astarte_appengine_api.yaml (yaml:1711-1717) — the schema requires only `data` and types it `object`, but every `format=table` data response carries a sibling `metadata` object. Code: `astarteapi.WriteDataWithMetadata` renders `{"data": v, "metadata": m}` through `metadataEnvelope` (`Metadata map[string]any json:"metadata,omitempty"`, pkg/astarteapi/envelope.go:149-158), the `Tabular` carrier sets it (internal/appengine/data.go:59-64), and both renderers populate it — `renderIndividual` (data.go:256-268, `columns` keyed on the path's last segment + `table_header`) and `renderObject` (data.go:293-321, `columns` + sorted `table_header`). Add `metadata` as an optional object property (`columns: map<string,int>`, `table_header: array of string`) described as present only when `format=table`, on all six data GET 200s that `$ref` this schema (yaml:416, 468, 685, 730, 1244, 1290). While there, correct `data`'s type/description: it is not always an object — `format=structured` individual is `[{value,timestamp}]` (data.go:275-281) and `format=disjoint_tables` is an object keyed by path (data.go:269-274, 322-331). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-ae-data-write-object-errors [auto]: document the two object-aggregated write rejections on the six data PUT/POST ops (putData/publishData, by-alias, in-group) in docs/api/astarte_appengine_api.yaml — each refs `BadRequest` for 400 (yaml:518, 570, 775, 822, 1336, 1384) and `ValueTooLarge` for 422 (yaml:528, 580, 785, 832, 1346, 1394), but three bodies are missing. (a) 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` from `writeBadObjectError` (internal/appengine/http.go:722-748, const `detailUnexpectedObjectKey` :711), reached by the `payload.ReasonOf(err) == payload.ReasonBadObject` case (http.go:684-691) — so the shared `ErrorDetail` schema (yaml:1741-1750) needs an optional `unexpected_keys` array, or a dedicated response. (b) 422 `{"errors":{"detail":"Missing required mapping key"}}` (const `detailMissingRequiredMapping`, http.go:713-720) from the `ReasonMissingRequired` case (http.go:692-702) — add an example beside the value-size one already on those ops. Both were implemented by appengine-unexpected-object-key/appengine-missing-required-422 (upstream master b6d46ad4, #2237) with no spec update. Do NOT invent other `payload.ReasonOf` statuses — that family is the separate blocked line appengine-payload-reason-status-map. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-ae-group-create-devices-required [auto]: in docs/api/astarte_appengine_api.yaml, `GroupCreate` (yaml:1652-1663) marks only `group_name` required, but `createGroup` treats `devices` as required — `groupBody.Devices *[]string` (internal/appengine/http.go:391-397) makes a MISSING key a 422 `{"errors":{"devices":["can't be blank"]}}` and an EMPTY array a 422 `{"errors":{"devices":["should have at least 1 item(s)"]}}` (http.go:405-418), and only the MISSING-key case is unreachable from the schema as written. The shared `ValidationErrors` example (yaml:1841-1844) that POST /groups refs (yaml:948) shows `group_name: can't be blank` alone. Add `devices` to `GroupCreate`'s `required` list (the 201 echo always carries it, http.go:424) and add the two device messages to the POST 422 story. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [!] docs-sync-ae-notfound-component [auto]: the `NotFound` response component in docs/api/astarte_appengine_api.yaml (yaml:1795-1803, example `errors.detail: Not Found`) is `$ref`d by all twelve group-route 404s (yaml:978, 1026, 1070, 1104, 1152, 1179, 1209, 1250, 1296, 1342, 1390, 1429), yet `astarteapi.WriteNotFound`/`DetailNotFound` ("Not Found", pkg/astarteapi/envelope.go:42-43, 205-207) is emitted nowhere in internal/appengine (rg finds only `WriteDeviceNotFound`, http.go:683). A missing group answers 404 `{"errors":{"detail":"Group not found"}}` (`ErrGroupNotFound`, internal/appengine/service.go:41-42, 508-516; writeError http.go:667-668) and a missing member device answers `"Device not found"` (writeError store.ErrNotFound → WriteDeviceNotFound, http.go:682-683). Repoint/rename the component so the group routes show `Group not found` (add a `GroupNotFound` component, keep `DeviceNotFound` on the routes that can 404 for a missing member) and state which routes can emit which. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
+- [ ] realm-error-mapping-test [auto]: add a container-free test file `internal/realm/errors_test.go` pinning the RM error-mapping layer, which is asserted today only by the integration-tagged `internal/realm/http_test.go` (TimescaleDB) so the Pi cannot catch a regression: (1) `writeInterfaceError` (http.go:352-399) — `ErrMaximumDatabaseRetentionExceeded`→422 `{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, a `*interfaceschema.ViolationsError`→422 nested envelope, `errMajorNotFound`→404 "Interface major not found", `ErrNameCollision`/`ErrNameMismatch`/`ErrMajorMismatch`/`interfaceschema.ErrMinorNotIncreased`/`ErrDowngradeNotAllowed`/`ErrMissingEndpoints`/`ErrIncompatibleEndpointChange`/`store.ErrAlreadyExists`→409, `store.ErrInterfaceMajorNotZero`/`store.ErrInterfaceInUse`→403, plain `ErrValidation`→422 with the `realm: validation failed: ` prefix stripped, `store.ErrNotFound`→404 "Interface not found", default→500; (2) `writeError` (http.go:402-425) — its 422/409/404/default arms including that it maps the two interface sentinels to 422, not 403; (3) `renderViolationsBody` (http.go:443-533) byte-for-byte for a mixed input (interface-level field merged first-seen + the full-length index-aligned `mappings` array with `{}` for clean entries); (4) `validationDetail` — prefix stripped, non-prefix passthrough, exact-prefix NOT stripped. `*API{}` needs no Service (these methods never deref `a.svc`), and `interfaceschema.ViolationsError` is built from its exported `Violations`/`MappingCount` fields, so no store or Docker.
+- [ ] realm-detailed-properties-branch-test [auto]: extend `internal/realm/detailed_required_encrypted_test.go` to cover the properties branch of `writeDetailedMapping` (internal/realm/detailed.go:61-63) — currently only the datastream branch is pinned container-free (the properties path runs only in the integration-tagged `internal/realm/detailed_test.go`): assert a properties mapping renders exactly `{"endpoint":..,"type":..,"allow_unset":<bool>}` with no reliability/retention/expiry/explicit_timestamp/required/encrypted/database_retention_policy keys, for `allow_unset` both true and false. No Docker.

$ cd /root/astrate-mule && tail -2 .mule/todo.md && echo "---" && ls -la .mule/reviews/realm-2026-10-09.md
- [ ] realm-error-mapping-test [auto]: add a container-free test file `internal/realm/errors_test.go` pinning the RM error-mapping layer, which is asserted today only by the integration-tagged `internal/realm/http_test.go` (TimescaleDB) so the Pi cannot catch a regression: (1) `writeInterfaceError` (http.go:352-399) — `ErrMaximumDatabaseRetentionExceeded`→422 `{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, a `*interfaceschema.ViolationsError`→422 nested envelope, `errMajorNotFound`→404 "Interface major not found", `ErrNameCollision`/`ErrNameMismatch`/`ErrMajorMismatch`/`interfaceschema.ErrMinorNotIncreased`/`ErrDowngradeNotAllowed`/`ErrMissingEndpoints`/`ErrIncompatibleEndpointChange`/`store.ErrAlreadyExists`→409, `store.ErrInterfaceMajorNotZero`/`store.ErrInterfaceInUse`→403, plain `ErrValidation`→422 with the `realm: validation failed: ` prefix stripped, `store.ErrNotFound`→404 "Interface not found", default→500; (2) `writeError` (http.go:402-425) — its 422/409/404/default arms including that it maps the two interface sentinels to 422, not 403; (3) `renderViolationsBody` (http.go:443-533) byte-for-byte for a mixed input (interface-level field merged first-seen + the full-length index-aligned `mappings` array with `{}` for clean entries); (4) `validationDetail` — prefix stripped, non-prefix passthrough, exact-prefix NOT stripped. `*API{}` needs no Service (these methods never deref `a.svc`), and `interfaceschema.ViolationsError` is built from its exported `Violations`/`MappingCount` fields, so no store or Docker.
- [ ] realm-detailed-properties-branch-test [auto]: extend `internal/realm/detailed_required_encrypted_test.go` to cover the properties branch of `writeDetailedMapping` (internal/realm/detailed.go:61-63) — currently only the datastream branch is pinned container-free (the properties path runs only in the integration-tagged `internal/realm/detailed_test.go`): assert a properties mapping renders exactly `{"endpoint":..,"type":..,"allow_unset":<bool>}` with no reliability/retention/expiry/explicit_timestamp/required/encrypted/database_retention_policy keys, for `allow_unset` both true and false. No Docker.
---
-rw-r--r-- 1 root root 6329 Oct  9 21:54 .mule/reviews/realm-2026-10-09.md

Done: reviewed `internal/realm` (last done 2026-09-09, code unchanged since), wrote the review file and appended 2 `[auto]` task lines.
Files: .mule/reviews/realm-2026-10-09.md, .mule/todo.md
Verified: `go test ./internal/realm/` (non-integration subset) -> pass; confirmed via rg that http_test.go/dashboard_compat_test.go are the only tests touching the error mappers and both are `//go:build integration`.
Unsure: no category-1 bug found; the two proposals are category-2 test gaps. `writeError`'s two interface-sentinel branches look dead but I left them (recorded in the review).
Follow-ups: `realm-error-mapping-test`, `realm-detailed-properties-branch-test` (appended to .mule/todo.md). No code changed, no git touched.
