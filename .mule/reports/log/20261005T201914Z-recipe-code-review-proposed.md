slug: recipe-code-review
verdict: proposed
at:  1e8d87f
ran: 2026-10-05T20:19:14Z on DietPi in 125s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md

$ ls -la /root/astrate-mule/.mule/ && ls -la /root/astrate-mule/.mule/reviews/ 2>/dev/null
total 12948
drwxr-xr-x  7 root root     4096 Oct  5 21:23 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       14 Oct  5 22:17 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  5 20:14 .heartbeat
-rw-r--r--  1 root root        0 Oct  5 22:17 .last-output
-rw-r--r--  1 root root        0 Oct  5 22:17 .lock
-rw-r--r--  1 root root       12 Oct  5 22:17 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 12860933 Oct  5 22:17 cron.log
drwxr-xr-x  2 root root     4096 Oct  4 22:35 failed
-rw-r--r--  1 root root   139887 Oct  5 22:16 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    14600 Oct  5 20:14 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  4 21:08 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   149912 Oct  5 20:14 todo.md
total 216
drwxr-xr-x 2 root root  4096 Oct  4 21:08 .
drwxr-xr-x 7 root root  4096 Oct  5 21:23 ..
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
-rw-r--r-- 1 root root  6002 Sep 12 20:24 payload-2026-09-12.md
-rw-r--r-- 1 root root  4259 Sep  9 21:13 realm-2026-09-09.md
-rw-r--r-- 1 root root  4605 Sep  8 18:59 store-2026-09-08.md
-rw-r--r-- 1 root root  5504 Sep 21 13:32 store-2026-09-21.md
-rw-r--r-- 1 root root  4511 Oct  4 21:08 store-2026-10-04.md
-rw-r--r-- 1 root root  2644 Sep 24 20:37 swagger-2026-09-24.md

$ ls /root/astrate-mule/ && echo "---INTERNAL---" && ls /root/astrate-mule/internal/ 2>/dev/null && echo "---PKG---" && ls /root/astrate-mule/pkg/ 2>/dev/null
AGENTS.md
Dockerfile
LICENSE
Makefile
README.md
astrate
bench
cmd
deploy
docker-compose.yml
docs
examples
flow-container-echo
go.mod
go.sum
internal
migrations
opencode.json
pkg
site-dist
test
tools
---INTERNAL---
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
---PKG---
astarteapi
deviceid
interfaceschema
payload

$ cd /root/astrate-mule && tail -40 .mule/log.md
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

$ cd /root/astrate-mule && for d in internal/*/ pkg/*/ cmd/ bench/ tools/ test/; do n=$(rg -c '' --glob "$d**/*.go" 2>/dev/null | awk -F: '{s+=$2} END {print s+0}'); f=$(rg --files --glob "$d**/*.go" 2>/dev/null | wc -l); echo "$n lines / $f files  $d"; done | sort -rn
15552 lines / 49 files  internal/engine/
11751 lines / 49 files  internal/flow/
7535 lines / 23 files  internal/appengine/
6874 lines / 32 files  internal/store/
4147 lines / 13 files  test/
3757 lines / 13 files  internal/broker/
3407 lines / 12 files  pkg/interfaceschema/
2942 lines / 10 files  internal/pairing/
2798 lines / 10 files  pkg/payload/
2730 lines / 11 files  internal/realm/
2335 lines / 8 files  internal/flowapi/
1658 lines / 10 files  bench/
1588 lines / 8 files  internal/auth/
1500 lines / 2 files  internal/swagger/
1394 lines / 10 files  cmd/
1372 lines / 5 files  internal/housekeeping/
1004 lines / 2 files  internal/config/
697 lines / 6 files  internal/testutil/
525 lines / 2 files  pkg/astarteapi/
451 lines / 5 files  internal/observability/
369 lines / 4 files  internal/httpx/
358 lines / 2 files  pkg/deviceid/
114 lines / 1 files  tools/

$ cd /root/astrate-mule && echo "=== cmd ===" && ls -R cmd/ && echo "=== bench ===" && ls -la bench/ && echo "=== tools ===" && ls -la tools/ && echo "=== test ===" && ls -la test/
=== cmd ===
cmd/:
astrate

cmd/astrate:
devcert_test.go
forward_nats_test.go
forward_test.go
main.go
main_test.go
newnats_default.go
newnats_nats.go
sealer_test.go
shutdown_test.go
version_test.go
=== bench ===
total 108
drwxr-xr-x  5 root root  4096 Jul 27 05:10 .
drwxr-xr-x 17 root root  4096 Sep 25 19:41 ..
-rw-r--r--  1 root root  5901 Jul 27 05:10 README.md
drwxr-xr-x  2 root root  4096 Jul 27 05:10 astrate
-rw-r--r--  1 root root 11444 Jul 27 05:10 client.go
-rw-r--r--  1 root root  4481 Jul 27 05:10 client_test.go
-rw-r--r--  1 root root  4567 Jul 27 05:10 connstorm.go
-rw-r--r--  1 root root  5068 Jul 27 05:10 device.go
-rw-r--r--  1 root root   656 Jul 27 05:10 go.mod
-rw-r--r--  1 root root  1364 Jul 27 05:10 go.sum
-rw-r--r--  1 root root 11158 Jul 27 05:10 ingest.go
drwxr-xr-x  2 root root  4096 Jul 27 05:10 interfaces
-rw-r--r--  1 root root  1973 Jul 27 05:10 main.go
-rw-r--r--  1 root root  5086 Jul 27 05:10 provision.go
-rw-r--r--  1 root root  3510 Jul 27 05:10 query.go
drwxr-xr-x  3 root root  4096 Jul 27 05:32 scripts
-rw-r--r--  1 root root  2406 Jul 27 05:10 state.go
-rw-r--r--  1 root root  2038 Jul 27 05:10 stats.go
=== tools ===
total 92
drwxr-xr-x  3 root root  4096 Sep  4 21:26 .
drwxr-xr-x 17 root root  4096 Sep 25 19:41 ..
drwxr-xr-x  2 root root  4096 Jul 27 05:10 bsoncapture
-rwxr-xr-x  1 root root  1961 Aug 31 01:29 mule-plan-day.sh
-rwxr-xr-x  1 root root 12902 Sep  4 20:46 mule-setup-pi.sh
-rwxr-xr-x  1 root root  6188 Sep  4 20:46 mule-survey.sh
-rwxr-xr-x  1 root root 48895 Sep  4 21:26 mule.sh
-rwxr-xr-x  1 root root  3216 Jul 27 22:18 reconcile.sh
=== test ===
total 16
drwxr-xr-x  4 root root 4096 Jul 27 05:10 .
drwxr-xr-x 17 root root 4096 Sep 25 19:41 ..
drwxr-xr-x 11 root root 4096 Sep  4 20:46 conformance
drwxr-xr-x  2 root root 4096 Jul 27 05:10 security

$ cd /root/astrate-mule && rg --files --glob 'bench/**/*.go' | xargs wc -l | sort -rn
1658 total
  362 bench/client.go
  349 bench/ingest.go
  163 bench/device.go
  159 bench/provision.go
  145 bench/connstorm.go
  134 bench/client_test.go
  125 bench/query.go
   82 bench/stats.go
   71 bench/main.go
   68 bench/state.go

$ cd /root/astrate-mule && echo "=== bench/astrate ===" && ls -la bench/astrate/ && echo "=== bench/interfaces ===" && ls -la bench/interfaces/ && echo "=== bench/scripts ===" && ls -laR bench/scripts/ && echo "=== bench/go.mod ===" && cat bench/go.mod
=== bench/astrate ===
total 12
drwxr-xr-x 2 root root 4096 Jul 27 05:10 .
drwxr-xr-x 5 root root 4096 Jul 27 05:10 ..
-rw-r--r-- 1 root root 2053 Jul 27 05:10 compose.yml
=== bench/interfaces ===
total 16
drwxr-xr-x 2 root root 4096 Jul 27 05:10 .
drwxr-xr-x 5 root root 4096 Jul 27 05:10 ..
-rw-r--r-- 1 root root  515 Jul 27 05:10 org.astrate.bench.Individual.json
-rw-r--r-- 1 root root  843 Jul 27 05:10 org.astrate.bench.Object.json
=== bench/scripts ===
bench/scripts/:
total 40
drwxr-xr-x 3 root root 4096 Jul 27 05:32 .
drwxr-xr-x 5 root root 4096 Jul 27 05:10 ..
-rw-r--r-- 1 root root  971 Jul 27 05:10 astarte-compose-override.yml
-rwxr-xr-x 1 root root 1284 Jul 27 05:10 disk-usage.sh
-rwxr-xr-x 1 root root 6944 Jul 27 05:45 run-tier.sh
-rwxr-xr-x 1 root root 1606 Jul 27 05:10 sample-stats.sh
drwxr-xr-x 2 root root 4096 Jul 27 05:32 tiers
-rwxr-xr-x 1 root root 3827 Jul 27 05:10 up-astarte.sh
-rwxr-xr-x 1 root root 2091 Jul 27 05:10 up-astrate.sh

bench/scripts/tiers:
total 28
drwxr-xr-x 2 root root 4096 Jul 27 05:32 .
drwxr-xr-x 3 root root 4096 Jul 27 05:32 ..
-rw-r--r-- 1 root root 1551 Jul 27 05:32 README.md
-rw-r--r-- 1 root root  154 Jul 27 05:31 big.env
-rw-r--r-- 1 root root  159 Jul 27 05:31 giant.env
-rw-r--r-- 1 root root  154 Jul 27 05:31 medium.env
-rw-r--r-- 1 root root  147 Jul 27 05:31 small.env
=== bench/go.mod ===
// Standalone benchmark harness (see README.md). Deliberately a separate
// module that imports NOTHING from astrate: it speaks only the wire protocols
// (Astarte MQTT v1 + the REST APIs), so the same binary drives an upstream
// Astarte deployment and an Astrate deployment interchangeably.
module github.com/astrate-platform/astrate/bench

go 1.26.1

require (
	github.com/eclipse/paho.mqtt.golang v1.5.1 // same pin as the main module
	github.com/golang-jwt/jwt/v5 v5.3.1
	go.mongodb.org/mongo-driver/v2 v2.6.0
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
)

$ cd /root/astrate-mule && echo "=== Makefile ===" && cat Makefile && echo "=== gate refs in mule.sh ===" && rg -n 'bench' tools/mule.sh | head -40
=== Makefile ===
# Astrate build & verification entry points (docs/ROADMAP.md §1.2, verification tiers §0.2).
#
#   make build               static binary / compile check (CGO disabled)
#   make lint                golangci-lint (config: .golangci.yml)
#   make test                T1 unit tests (no Docker, no network)
#   make test-integration    T2 tests (testcontainers -> timescale/timescaledb:latest-pg16)
#   make test-e2e            T3 component/E2E tests (full wired binary; populated from M5/M6)
#   make test-conformance    T4 official-SDK conformance harness (lands in M9 under test/)
#   make up / make down      local TimescaleDB via docker compose

GO                     ?= go
GOLANGCI_LINT_VERSION  ?= v2.12.2
GOFLAGS                := -trimpath
LDFLAGS                := -s -w
DIST                   := dist

.DEFAULT_GOAL := build

.PHONY: build lint test test-integration test-e2e test-conformance tools up down clean

## build: compile every package statically; emits $(DIST)/astrate once cmd/astrate lands (M8).
build:
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' ./...
	@if [ -d cmd/astrate ]; then \
		mkdir -p $(DIST); \
		CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(DIST)/astrate ./cmd/astrate; \
		echo "built $(DIST)/astrate"; \
	fi

## lint: run the pinned golangci-lint (install with `make tools`).
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not found; run 'make tools' (installs $(GOLANGCI_LINT_VERSION))"; exit 1; }
	golangci-lint run ./...

## test: T1 — pure unit tests, race detector on.
test:
	$(GO) test -race ./...

## test-integration: T2 — requires a Docker daemon (or ASTRATE_TEST_DSN to reuse a database).
test-integration:
	$(GO) test -race -count=1 -tags integration ./...

## test-e2e: T3 — component tests behind the e2e build tag (suites land from M5 onward).
test-e2e:
	$(GO) test -race -count=1 -tags "integration e2e" ./...

## test-conformance: T4 — official-SDK harness; lives in test/ with its own go.mod (M9).
test-conformance:
	@if [ -f test/conformance/go.mod ]; then \
		cd test/conformance && $(GO) test -count=1 ./...; \
	else \
		echo "conformance harness lands in M9 (docs/ROADMAP.md §10); nothing to run yet"; \
	fi

## tools: install the pinned developer toolchain.
tools:
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

## up/down: local TimescaleDB for development and DSN-reuse test runs.
up:
	docker compose up -d

down:
	docker compose down

clean:
	rm -rf $(DIST)
=== gate refs in mule.sh ===
128:    # nobody made. bench/ imports nothing from Astrate and needs no cgo, so cross-compiling
130:    bench-push)
132:      local out="$MULE/.bench-linux-amd64"
133:      note "cross-compiling bench for linux/amd64"
134:      ( cd "$REPO/bench" && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$out" . ) \
136:      legion_sh "mkdir -p '$MULE_LEGION_REPO/bench'" || die "mkdir failed on the Legion Go"
138:      scp $MULE_SSH_OPTS "$out" "$MULE_LEGION_SSH:$MULE_LEGION_REPO/bench/bench" \
140:      legion_sh "chmod +x '$MULE_LEGION_REPO/bench/bench'"
141:      ok "bench binary is on the Legion Go at $MULE_LEGION_REPO/bench/bench"
143:    *) die "legion: check | sh <cmd> | bench-push";;
884:   4  bench      run the tiered benchmarks (small/medium/big/giant)
925:    4|bench)    name=benchmarks;;

→ Read bench/client.go
→ Read bench/ingest.go
→ Read bench/stats.go
→ Read bench/device.go
→ Read bench/provision.go
→ Read bench/connstorm.go
→ Read bench/state.go
→ Read bench/query.go
→ Read bench/main.go
→ Read bench/client_test.go

$ cd /root/astrate-mule && cat bench/scripts/tiers/*.env && echo "=== run-tier.sh ===" && cat bench/scripts/run-tier.sh
# Tier: big — Legion Go
# 10000 devices × 0.1 msg/s = 1000 aggregate msg/s
TIER_NAME=big
DEVICES=10000
RATE=0.1
INGEST_DURATION=15m
STORM_DEVICES=5000
# Tier: giant — Legion Go
# 50000 devices × 0.1 msg/s = 5000 aggregate msg/s
TIER_NAME=giant
DEVICES=50000
RATE=0.1
INGEST_DURATION=20m
STORM_DEVICES=20000
# Tier: medium — laptop
# 1000 devices × 0.2 msg/s = 200 aggregate msg/s
TIER_NAME=medium
DEVICES=1000
RATE=0.2
INGEST_DURATION=10m
STORM_DEVICES=1000
# Tier: small — laptop
# 100 devices × 0.2 msg/s = 20 aggregate msg/s
TIER_NAME=small
DEVICES=100
RATE=0.2
INGEST_DURATION=5m
STORM_DEVICES=100
=== run-tier.sh ===
#!/usr/bin/env bash
# Runs a named tier against a deployment and writes all results to a single
# timestamped directory.  Sourced tier env vars drive the device counts and
# rates; the base URL and housekeeping key are passed on the command line.
#
#   bench/scripts/run-tier.sh <tier> <target> -base-url <url> -housekeeping-key <file>
#
# <tier>   small | medium | big | giant  (must have a matching tiers/<tier>.env)
# <target> label for the system under test, e.g. "astrate" or "astarte"
#
# Resumable: if the state file for this tier already has enough devices,
# provision is skipped rather than re-registering tens of thousands of them.
# Never overwrites an existing results directory — the benchmark result is
# evidence, and a duplicate means something is wrong.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BENCH_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# ── parse args ──────────────────────────────────────────────────────────────
TIER="${1:?usage: run-tier.sh <tier> <target> -base-url <url> -housekeeping-key <file>}"
TARGET="${2:?usage: run-tier.sh <tier> <target> -base-url <url> -housekeeping-key <file>}"
shift 2

BASE_URL=""
HK_KEY=""
while [[ $# -gt 0 ]]; do
    case "$1" in
        -base-url)     BASE_URL="$2"; shift 2 ;;
        -housekeeping-key) HK_KEY="$2"; shift 2 ;;
        *) echo "run-tier.sh: unknown arg $1" >&2; exit 2 ;;
    esac
done
if [[ -z "$BASE_URL" || -z "$HK_KEY" ]]; then
    echo "run-tier.sh: -base-url and -housekeeping-key are required" >&2
    exit 2
fi

TIER_ENV="$SCRIPT_DIR/tiers/${TIER}.env"
if [[ ! -f "$TIER_ENV" ]]; then
    echo "run-tier.sh: tier file not found: $TIER_ENV" >&2
    exit 1
fi
# shellcheck source=tiers/small.env
source "$TIER_ENV"

echo "tier=$TIER target=$TARGET devices=$DEVICES rate=$RATE duration=$INGEST_DURATION"

# ── build bench once ────────────────────────────────────────────────────────
# Use a prebuilt binary when there is one: the Legion Go, which runs the heavy
# tiers, has no Go toolchain -- the binary is cross-compiled elsewhere and
# copied in by `mule.sh legion bench-push`. Only build when we actually can.
BENCH_BIN="$BENCH_DIR/bench"
if [[ -x "$BENCH_BIN" ]]; then
    echo "using prebuilt bench binary: $BENCH_BIN"
elif command -v go >/dev/null; then
    echo "building bench..."
    go build -o "$BENCH_BIN" "$BENCH_DIR"
    trap 'rm -f "$BENCH_BIN"' EXIT
else
    echo "run-tier.sh: no bench binary and no Go toolchain here." >&2
    echo "  Cross-compile and copy one in:  tools/mule.sh legion bench-push" >&2
    exit 1
fi

# ── results directory ───────────────────────────────────────────────────────
TS="$(date -u +%Y%m%d-%H%M%S)"
RESULTS_DIR="$BENCH_DIR/results/${TIER}-${TARGET}-${TS}"
if [[ -d "$RESULTS_DIR" ]]; then
    echo "run-tier.sh: results directory already exists — refusing to overwrite:" >&2
    echo "  $RESULTS_DIR" >&2
    exit 1
fi
mkdir -p "$RESULTS_DIR"

# ── host info ───────────────────────────────────────────────────────────────
{
    echo "=== host ==="
    uname -a
    echo ""
    echo "=== cpu ==="
    nproc
    echo ""
    echo "=== memory ==="
    free -h 2>/dev/null || vm_stat 2>/dev/null || echo "(unknown)"
} > "$RESULTS_DIR/host.txt"

cp "$TIER_ENV" "$RESULTS_DIR/"

# ── background stats sampler ────────────────────────────────────────────────
STATS_CSV="$RESULTS_DIR/stats.csv"
if [[ -x "$SCRIPT_DIR/sample-stats.sh" ]]; then
    "$SCRIPT_DIR/sample-stats.sh" "$STATS_CSV" 5 &
    STATS_PID=$!
    trap 'kill "$STATS_PID" 2>/dev/null; rm -f "$BENCH_BIN"' EXIT
else
    echo "note: sample-stats.sh not found or not executable — skipping resource sampling"
fi

# ── provision ───────────────────────────────────────────────────────────────
# Skip if a state file already exists with enough devices (resumable).
STATE="$RESULTS_DIR/state.json"
if [[ -f "$STATE" ]]; then
    existing=$(python3 -c "import json; print(len(json.load(open('$STATE'))['devices']))" 2>/dev/null || echo 0)
    if (( existing >= DEVICES )); then
        echo "provision: state file has $existing devices (need $DEVICES) — reusing"
    else
        echo "provision: state file has $existing devices (need $DEVICES) — re-provisioning"
        "$BENCH_BIN" provision \
            -base-url "$BASE_URL" \
            -housekeeping-key "$HK_KEY" \
            -devices "$DEVICES" \
            -state "$STATE"
    fi
else
    "$BENCH_BIN" provision \
        -base-url "$BASE_URL" \
        -housekeeping-key "$HK_KEY" \
        -devices "$DEVICES" \
        -state "$STATE"
fi

# ── ingest ──────────────────────────────────────────────────────────────────
"$BENCH_BIN" ingest \
    -state "$STATE" \
    -devices "$DEVICES" \
    -rate "$RATE" \
    -duration "$INGEST_DURATION" \
    -out "$RESULTS_DIR"

# ── connstorm ───────────────────────────────────────────────────────────────
"$BENCH_BIN" connstorm \
    -state "$STATE" \
    -devices "$STORM_DEVICES" \
    -out "$RESULTS_DIR"

# ── query ───────────────────────────────────────────────────────────────────
"$BENCH_BIN" query \
    -state "$STATE" \
    -devices "$DEVICES" \
    -out "$RESULTS_DIR"

# ── done ────────────────────────────────────────────────────────────────────
if [[ -n "${STATS_PID:-}" ]]; then
    kill "$STATS_PID" 2>/dev/null || true
    wait "$STATS_PID" 2>/dev/null || true
fi

echo ""
echo "done — results in $RESULTS_DIR"
echo "  host:        $RESULTS_DIR/host.txt"
echo "  tier env:    $RESULTS_DIR/${TIER}.env"
echo "  stats:       $STATS_CSV"
echo "  provision:   $RESULTS_DIR/provision-*.json"
echo "  ingest:      $RESULTS_DIR/ingest-*.json"
echo "  connstorm:   $RESULTS_DIR/connstorm-*.json"
echo "  query:       $RESULTS_DIR/query-*.json"

$ cd /root/astrate-mule && go version && echo "--- root go list (does it descend into bench?) ---" && go list ./... 2>&1 | head -30 && echo "--- bench vet ---" && (cd bench && go vet ./... 2>&1 | head -20)
go version go1.26.7 linux/arm64
--- root go list (does it descend into bench?) ---
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
--- bench vet ---

$ cd /root/astrate-mule && rg -n 'go test|go vet|gofmt|golangci|make test|make lint' tools/mule.sh | head -30
195:  ( cd "$REPO" && eval "${MULE_TEST_CMD:-go test ./...}" ) >/dev/null 2>&1 && still_passes=1
221:# errors, test failures and golangci-lint findings all carry a file:line prefix.

$ cd /root/astrate-mule && rg -n -i 'loss|percentile|p99|nearest|quantil' bench/README.md | head -40
11:bench ingest      N devices × R msg/s: e2e latency, PUBACK latency, loss
60:**Loss** counts rows through the AppEngine API against messages sent. On
61:Astarte, rows lag PUBACKs by design — if loss is nonzero, first increase
70:`-devices`/`-rate` and find the knee where e2e p99 degrades or `behind_ticks`
88:  ingest p99 < 250 ms — DESIGN §5.4) make useful chart annotations.
107:Divide by the row count (`loss_check_found` totals, or query the stores).
120:- The loss check pages with `since_after` above ~9000 rows per device. Row

$ cd /root/astrate-mule && sed -n '50,130p' bench/README.md
platforms. It is quantized by `-probe-poll` (default 25 ms) — report it with
that caveat, or lower the interval and note the extra read load you injected.

**PUBACK latency is NOT comparable across platforms — never headline it.**
Astrate withholds PUBACK until the row is committed (ack-after-commit);
Astarte acks at VerneMQ and persists asynchronously (RabbitMQ → data updater
plant → Cassandra). The same number means "durably stored" on one platform
and "accepted by the broker" on the other. It is reported because it *is*
each platform's device-visible backpressure behavior.

**Loss** counts rows through the AppEngine API against messages sent. On
Astarte, rows lag PUBACKs by design — if loss is nonzero, first increase
`-grace` before concluding anything.

**Closed loop per device:** each simulated device keeps one message in
flight (publish, wait PUBACK, next tick). If acks are slower than the tick,
the device falls behind — `behind_ticks` > 0 means the platform, not the
generator, set the achieved rate. That is a result, not an error.

**Throughput ceiling:** there is no auto-ramp; run `ingest` at increasing
`-devices`/`-rate` and find the knee where e2e p99 degrades or `behind_ticks`
explodes. Change one variable per run.

Other traps the harness cannot absorb for you:

- **Durability knobs decide the winner silently.** PG `synchronous_commit`
  vs Cassandra commitlog sync — run defaults, but document both alongside
  results.
- **Warm-up:** Cassandra behaves differently once compaction starts. Run
  ≥ 15 min steady state; discard the first minutes; report the steady tail.
- **Host sizing is the headline experiment.** Run the matrix on a
  comfortable host (8–16 GB) for a fair performance comparison, *and* on the
  1–2 GB VPS target where "Astarte does not fit" is itself the result.
- **Generator placement:** run `bench` on a separate machine (or CPU-pin it)
  so it does not steal cycles from the system under test. For a remote
  generator against Astarte, map `*.astarte.localhost` to the Docker host in
  `/etc/hosts` and pass `-broker-url` if needed.
- **Reference lines:** Astrate's own budgets (RSS ≤ 150 MB app, ≤ 768 MB PG,
  ingest p99 < 250 ms — DESIGN §5.4) make useful chart annotations.

## Storage efficiency (bytes per datapoint)

After a large ingest, force both platforms to their at-rest representation,
then measure:

```sh
# Astrate: compress the chunks, then measure
docker exec -it astrate-bench-timescaledb-1 psql -U astrate -c \
  "SELECT compress_chunk(c, true) FROM show_chunks('individual_datastreams') c;"
./scripts/disk-usage.sh astrate

# Astarte: flush + major compaction, then measure
docker exec -it <cassandra-container> nodetool flush
docker exec -it <cassandra-container> nodetool compact
./scripts/disk-usage.sh astarte
```

Divide by the row count (`loss_check_found` totals, or query the stores).
Measure a *delta* between two ingests of known size if the baseline schema
overhead should be excluded.

## Mechanics worth knowing

- Both dev stacks use self-signed broker certs; devices default to skipping
  TLS verification (`-insecure-tls`, recorded in the state file). mTLS client
  auth is always on — that part is the real protocol on both platforms.
- Device keys are ECDSA P-256 by default so thousands generate quickly;
  `-rsa` switches to RSA-2048 if an Astarte deployment rejects EC CSRs.
- Provision retries HTTP 429 with backoff (Astrate rate-limits pairing);
  registration parallelism is `-concurrency`.
- The loss check pages with `since_after` above ~9000 rows per device. Row
  timestamps are ms precision on the wire — above ~100 msg/s/device,
  same-millisecond collisions can undercount the paged variant; keep
  per-device rates below that or extend the run instead.
- Interfaces publish with `reliability: guaranteed` (QoS 1) and explicit
  timestamps, so both platforms index identical time series.
- The state file contains the realm private key and device secrets:
  benchmark material, not secrets — but don't commit it (gitignored).

$ cd /root/astrate-mule && sed -n '185,215p' tools/mule.sh && echo "=== MULE_TEST_CMD default / hosts ===" && rg -n 'MULE_TEST_CMD' tools/mule.sh .mule/hosts .mule/config 2>/dev/null | head
note "gate: do the new tests actually fail without the change?"
  if ! git -C "$REPO" apply -R "$patch" 2>/dev/null; then
    # Could not isolate the implementation — say so rather than pretend the gate ran.
    note "  could not un-apply the implementation; skipping this gate"
    rm -f "$patch"; return 0
  fi

  local still_passes=0
  # A compile failure here is a pass, not an error: the test referencing something that no
  # longer exists is exactly the dependency we are trying to demonstrate.
  ( cd "$REPO" && eval "${MULE_TEST_CMD:-go test ./...}" ) >/dev/null 2>&1 && still_passes=1

  git -C "$REPO" apply "$patch" 2>/dev/null || {
    bad "could not restore the implementation after the proof gate — leaving the tree alone"
    rm -f "$patch"; return 1
  }
  rm -f "$patch"

  if [ "$still_passes" = 1 ]; then
    bad "the new tests pass with the implementation removed — they do not prove the change"
    return 1
  fi
  ok "the tests fail without the change, as they should"
  return 0
}

# A failing gate used to discard its own output, so every failure looked alike: the log row
# said "gates failed" and nothing said why. Forty-two consecutive tasks were blocked that way
# over five weeks by ten pre-existing lint errors in the checkout's own base — the gate lints
# the whole repo, so no task could ever pass, and the tasks queued to fix those very errors
# failed on the other nine. The output now survives in .mule/failed/gate.log and its first
=== MULE_TEST_CMD default / hosts ===
tools/mule.sh:195:  ( cd "$REPO" && eval "${MULE_TEST_CMD:-go test ./...}" ) >/dev/null 2>&1 && still_passes=1
tools/mule.sh:238:  if [ -n "${MULE_TEST_CMD:-}" ]; then
tools/mule.sh:239:    note "gate: $MULE_TEST_CMD"
tools/mule.sh:240:    ( cd "$REPO" && eval "$MULE_TEST_CMD" ) >>"$GATE_LOG" 2>&1 \
tools/mule.sh:300:  if [ -n "${MULE_TEST_CMD:-}" ]; then
tools/mule.sh:301:    if ( cd "$REPO" && eval "$MULE_TEST_CMD" ) >/dev/null 2>&1; then ok "test baseline is green"
tools/mule.sh:303:  else bad "MULE_TEST_CMD is empty — the mule would have no gate at all"; fail=1; fi
.mule/config:21:  MULE_TEST_CMD="go test -race ./..."
.mule/config:25:  MULE_TEST_CMD="go vet ./... && go test ./..."

! permission requested: external_directory (/tmp/*); auto-rejecting
✗ mkdir -p /tmp/benchprobe && cat > /tmp/benchprobe/main.go <<'EOF'
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"time"
)

type Device struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}
type Endpoints struct {
	Housekeeping string `json:"housekeeping"`
	RealmMgmt    string `json:"realm_management"`
	Pairing      string `json:"pairing"`
	AppEngine    string `json:"appengine"`
}
type State struct {
	BaseURL       string    `json:"base_url"`
	Endpoints     Endpoints `json:"endpoints"`
	Realm         string    `json:"realm"`
	RealmKeyPEM   string    `json:"realm_key_pem"`
	BrokerURL     string    `json:"broker_url"`
	TLSSkipVerify bool      `json:"tls_skip_verify"`
	Probes        int       `json:"probes"`
	Devices       []Device  `json:"devices"`
}

// verbatim copy of bench/state.go workload()
func (s *State) workload(n int) (probes, workers []Device, err error) {
	if s.Probes+n > len(s.Devices) {
		return nil, nil, fmt.Errorf("state has %d devices (%d probes + %d workload available), need %d workload — re-provision with more",
			len(s.Devices), s.Probes, len(s.Devices)-s.Probes, n)
	}
	return s.Devices[:s.Probes], s.Devices[s.Probes : s.Probes+n], nil
}

// verbatim of the interval expression + tick decision from bench/ingest.go:145-166
func tickDecision(rate float64) (interval time.Duration, behind bool, note string) {
	interval = time.Duration(float64(time.Second) / rate)
	next := time.Now()
	now := time.Now()
	wait := next.Add(interval).Sub(now)
	switch {
	case wait > 0:
		note = "sleep until tick"
	case -wait > interval:
		note = "behind-tick: reset next=now (SPIN)"
	default:
		note = "publish immediately"
	}
	return interval, -wait > interval, note
}

func main() {
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	st := &State{
		BaseURL: "http://127.0.0.1:8080", Realm: "bench",
		RealmKeyPEM: string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)})),
		BrokerURL:  "mqtts://127.0.0.1:8883", TLSSkipVerify: true, Probes: 3,
		Devices: []Device{{ID: "aaaaaaaaaaaaaaaaaaaaaa"}, {ID: "bbbbbbbbbbbbbbbbbbbbbb"}, {ID: "cccccccccccccccccccccc"},
			{ID: "dddddddddddddddddddddd"}, {ID: "eeeeeeeeeeeeeeeeeeeeee"}},
	}
	raw, _ := json.MarshalIndent(st, "", "  ")
	if err := os.WriteFile("/tmp/benchprobe/state.json", raw, 0o600); err != nil {
		panic(err)
	}
	fmt.Println("wrote /tmp/benchprobe/state.json")

	fmt.Println("\n--- interval arithmetic (bench/ingest.go:145) ---")
	for _, r := range []float64{0.1, 0.2, 1.0, 0, -1} {
		iv, behind, note := tickDecision(r)
		fmt.Printf("rate=%-5g interval=%-16v behind=%-5v %s\n", r, iv, behind, note)
	}

	fmt.Println("\n--- workload(-1) (bench/state.go:62) ---")
	func() {
		defer func() { if p := recover(); p != nil { fmt.Printf("PANIC: %v\n", p) } }()
		_, w, err := st.workload(-1)
		fmt.Printf("no panic: workers=%v err=%v\n", w, err)
	}()

	fmt.Println("\n--- workload with Probes > len(Devices) (state file with probes=9, devices=5) ---")
	bad := *st
	bad.Probes = 9
	func() {
		defer func() { if p := recover(); p != nil { fmt.Printf("PANIC: %v\n", p) } }()
		_, _, err := bad.workload(0)
		fmt.Printf("no panic: err=%v\n", err)
	}()

	fmt.Println("\n--- rand.Intn(len(workers)) with -devices 0 (bench/query.go:103) ---")
	_, workers, err := st.workload(0)
	fmt.Printf("workload(0): len(workers)=%d err=%v\n", len(workers), err)
	func() {
		defer func() { if p := recover(); p != nil { fmt.Printf("PANIC: %v\n", p) } }()
		_ = workers[0]
	}()
}
EOF
cd /tmp/benchprobe && go mod init probe >/dev/null 2>&1; go run . 2>&1 failed
Error: The user rejected permission to use this specific tool call.

$ cd /root/astrate-mule && cat .mule/.gitignore && echo "=== root gitignore ===" && cat .gitignore 2>/dev/null | head -20 && echo "=== results dir? ===" && ls bench/results 2>&1 | head -5
# The mule's own working files. None of these belong in git, and — more sharply — any of
# them that is not listed here makes the tree dirty, which makes `tick` refuse to run and so
# turns the timer into a permanent silent no-op.
failed/
task.md
.timeout
.lock
.budget
.bench-linux-amd64
.last-output
.rotation
.heartbeat
.alarmed
.impl.patch
hosts
research/
=== root gitignore ===
# Build artifacts
dist/
/astrate
*.exe
*.test
*.out

# Coverage
coverage.*
*.cover

# bbolt session stores (broker persistence, §3.1)
*.db

# Local env / editor
.env
.envrc
.idea/
.vscode/
*.swp
=== results dir? ===
ls: cannot access 'bench/results': No such file or directory

