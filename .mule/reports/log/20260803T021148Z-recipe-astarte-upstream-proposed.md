slug: recipe-astarte-upstream
verdict: proposed
at:  4beddf0
ran: 2026-08-03T02:11:48Z on DietPi in 485s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z
v1.4.0-rc.1	2026-07-14T08:53:27Z
v1.3.2	2026-07-14T10:07:26Z
v1.3.1	2026-07-13T11:50:05Z

→ Read .mule/todo.md
✱ Grep "v1\.\d+\.\d+" in .mule · 45 matches

$ ls .mule/reports/ 2>/dev/null; echo ---; ls .mule/reports/log/ 2>/dev/null | head -50
issue-13.md
issue-20.md
log
race-check.md
---
20260727T224629Z-race-check-blocked.md
20260727T233600Z-issue-27-blocked.md
20260728T001830Z-issue-26-done.md
20260728T005328Z-issue-25-done.md
20260728T013112Z-issue-24-blocked.md
20260728T021130Z-issue-23-done.md
20260728T024428Z-issue-20-checked.md
20260728T031805Z-issue-14-blocked.md
20260728T035046Z-issue-13-checked.md
20260728T043125Z-issue-12-done.md
20260728T050332Z-issue-20-checked.md
20260728T053605Z-issue-13-checked.md

✱ Grep "no upstream movement|astarte-upstream|upstream since" in . · 7 matches

$ ls docs/; echo ---; grep -rn "v1\.2\|v1\.3\|v1\.4\|targets\|targeting" docs/ 2>/dev/null | head -40
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui
---
docs/DESIGN.md:164:**Decision: shared tables + `realm_id`.** Astrate targets small installs (1–5 realms typical);
docs/DESIGN.md:523:Astarte mandates BSON (`{v, t}` documents). On ultra-constrained targets — in particular
docs/JSON-PAYLOAD-PROFILE.md:5:targets on ESP32/RP2040-class MCUs) who want to be first-class Astarte devices
docs/site/operations.md:106:For your own deployment: run `astarte/astarte-dashboard:v1.2.2` with a
docs/site/data-modeling.md:7:**Decision: shared tables + `realm_id`.** Row-level tenancy with composite keys is simpler and faster at the 1-5 realm scale Astrate targets. Realm deletion is a transactional cascade.
docs/site/compatibility.md:14:| Astarte Dashboard v1.2.2 | Runs unmodified |
docs/site/compatibility.md:27:- Astarte Dashboard v1.2.2 Device Live Events (since M11 Channels socket)
docs/site/json-payload-profile.md:5:targets on ESP32/RP2040-class MCUs) who want to be first-class Astarte devices
docs/COMPATIBILITY.md:47:- The Astarte Dashboard v1.2.2 runs unmodified against Astrate (compose
docs/COMPATIBILITY.md:74:     Recorded against upstream v1.2.0 in
docs/COMPATIBILITY.md:120:   v1.2.0 in `test/conformance/upstream/channels.json`, and two rows were
docs/OPERATIONS.md:107:For your own deployment: run `astarte/astarte-dashboard:v1.2.2` with a
docs/handoff/phase-8-memory.md:29:- `Makefile` — build/test targets
docs/handoff/phase-5-memory.md:17:- `docs/make sync` and `mkdocs build` steps match existing Makefile targets.
docs/handoff/phase-6-memory.md:13:- `docs/Makefile` — sync and build targets work
docs/handoff/phase-4-memory.md:15:- `docs/Makefile` — targets: install, sync, serve, build, clean

$ ls .mule/; echo ---; ls .mule/recipes/ 2>/dev/null
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
plans
recipes
reports
reviews
tasks
todo.md
---
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

$ tail -40 .mule/log.md 2>/dev/null; echo ===; tail -60 .mule/for-giulio.md 2>/dev/null
| 2026-07-27 | race-check | blocked | 60s | wrote nothing |
| 2026-07-27 | race-check | checked | 54s | 8cc751c |
| 2026-07-27 | race-check | checked | 51s | 886ff24 |
| 2026-07-27 | race-check | checked | 42s | 57338db |
| 2026-07-27 | race-check | checked | 56s | ccb4450 |
| 2026-07-27 | race-check | checked | 74s | a0ae9f4 |
| 2026-07-27 | race-check | checked | 82s | 7099f80 |
| 2026-07-27 | issue-6 | transient | 19s |  > build · big-pickle  Error: No provider available  |
| 2026-07-27 | issue-6 | transient | 20s |  > build · big-pickle  Error: No provider available  |
| 2026-07-27 | race-check | transient | 21s |  > build · big-pickle  Error: No provider available  |
| 2026-07-27 | race-check | transient | 22s |  > build · big-pickle  Error: No provider available  |
| 2026-07-27 | race-check | checked | 59s | 13cd7f1 |
| 2026-07-27 | race-check | checked | 46s | 452b303 |
| 2026-07-27 | race-check | checked | 65s | 5ac1fdb |
| 2026-07-27 | issue-6 | done | 401s | 997cdef |
| 2026-07-27 | race-check | checked | 88s | 70bc4be |
| 2026-07-27 | store-realm-cascade-policies | done | 122s | 0480501 |
| 2026-07-27 | race-check | checked | 108s | 16cc008 |
| 2026-07-27 | store-alias-lowest-id | done | 194s | 4af622f |
| 2026-07-27 | race-check | checked | 48s | 344d013 |
| 2026-07-27 | store-delete-device-objects | done | 321s | deb01ac |
| 2026-07-27 | race-check | checked | 71s | ad335ec |
| 2026-07-27 | issue-16 | done | 1092s | d670932 |
| 2026-07-27 | race-check | checked | 72s | cc13242 |
| 2026-07-27 | issue-20 | checked | 211s | 446b806 |
| 2026-07-27 | issue-22 | blocked | 340s | touched a never-touch path |
| 2026-07-27 | issue-21 | done | 297s | 9d85191 |
| 2026-07-27 | issue-15 | done | 209s | 5ea4c1b |
| 2026-07-28 | race-check | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-07-28 | issue-27 | blocked | 1157s | TIMEOUT after 1157s — task too big, split it |
| 2026-07-28 | issue-26 | done | 731s | e250381 |
| 2026-07-28 | issue-25 | done | 261s | ee59732 |
| 2026-07-28 | issue-24 | blocked | 433s | touched a never-touch path |
| 2026-07-28 | issue-23 | done | 567s | a9481d0 |
| 2026-07-28 | issue-20 | checked | 152s | c59a7eb |
| 2026-07-28 | issue-14 | blocked | 156s | gates failed |
| 2026-07-28 | issue-13 | checked | 136s | 4af2073 |
| 2026-07-28 | issue-12 | done | 607s | 89b1004 |
| 2026-07-28 | issue-20 | checked | 103s | 33d5be6 |
| 2026-07-28 | issue-13 | checked | 128s | beecab8 |
===
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

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
- **golangci-lint is not installed on the Pi**, so the mule's second gate is silently absent
  there — `gofmt` still runs, the linter does not. `go install
  github.com/golangci/golangci-lint/cmd/golangci-lint@<the pinned version>` on the Pi would
  close it; I did not pick a version for you, since the pin is a decision.
- ~~`/root/astrate` on the Pi has uncommitted work~~ — **resolved 2026-07-27** with the new
  `tools/reconcile.sh`: rescued onto `origin/wip/DietPi-20260727T171543Z` (pushed, not
  reviewed — read the diff before merging anything from it) and `/root/astrate` is now clean
  on `main`. Also had to set `commit.gpgsign false` locally on that clone first — it had
  signing on with no working gpg-agent for a non-interactive session, same class of problem
  `/root/astrate-mule` already solved the same way (your call, confirmed 2026-07-27: align it
  with the existing unattended-Pi-clone convention rather than fix gpg-agent or bypass signing
  silently inside the script).
- **The mule has been idle 15h.** Filed by the dead-man's switch; see journalctl on the Pi.

No prior run of this recipe is recorded, so this is a fresh comparison. The doc targets effectively v1.2.x (measured against v1.2.0, emulated API + Dashboard at v1.2.2); newest tag is now v1.4.0-rc.3, so upstream is ahead. Let me read the release notes for the gap.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.3" or .tag_name=="v1.4.0-rc.2" or .tag_name=="v1.4.0-rc.1" or .tag_name=="v1.3.2" or .tag_name=="v1.3.1" or .tag_name=="v1.3.0") | "=== \(.tag_name) (\(.published_at)) ===\n\(.body)\n"'
=== v1.4.0-rc.3 (2026-07-31T10:48:39Z) ===
Fourth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.2

### Added

- [secrets] Configurable base namespace for Vault using `ASTARTE_VAULT_BASE_NAMESPACE`: all created namespaces will be children of the base namespace.

### Fixed

- [fdo] Support for Vault and Vault managed
- [fdo] Properly handle chunked service info messages
- [astarte_data_updater_plant] Avoid clogging AMQP queues when receiving invalid device messages
- [astarte_data_updater_plant] Restore correct type validation checks for incoming data

=== v1.4.0-rc.2 (2026-07-14T10:16:22Z) ===
Third Astarte v1.4.0 release candidate.

=== v1.4.0-rc.1 (2026-07-14T08:53:27Z) ===
Second Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.0

### Added

- [astarte_housekeeping] Added HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION to set default, instance wise, realm default realm data retention, expressed in seconds

### Fixed

- [fdo] Do not crash when using rotated keys as owner keys
- [fdo] Delete devices from unfinished sessions
- [fdo] Delete old sessions from device
- [fdo] Properly respond to service info messages after IsDone=True
- [astarte_realm_management] Properly show default values for all mapping parameters

### Changed

- Enable FDO by default

=== v1.3.2 (2026-07-14T10:07:26Z) ===
Astarte v1.3.2 release.

## Changes since v1.3.1

## Fixed

- Avoid crash on invalid properties message

=== v1.3.1 (2026-07-13T11:50:05Z) ===
Astarte v1.3.1 release.

## Changes since v1.3.0

### Fixed

- Solve capabilities message stalling queues
- [astarte_housekeeping] Allow resuming partially completed realm deletions
- Use network topology as default realm replication
- Fix simple strategy realm creation without replication factor

### Changed

- Immediately delete devices in device deletion best case scenario

=== v1.3.0 (2026-05-06T15:30:50Z) ===
Astarte v1.3.0 release.

## Changes since 1.2.2

### Added

- [astarte_realm_management] realm migrations now create vhost for amqp triggers
- New environment variables for trigger notifications between realm management replicas and realm management -> pairing. These variables are currently being used only by realm management
  - `REALM_MANAGEMENT_CLUSTERING_KUBERNETES_SELECTOR`. The Endpoint label to query to get realm management instances. Defaults to `app=astarte-realm-management`.
  - `PAIRING_CLUSTERING_KUBERNETES_SELECTOR`. The Endpoint label to query to get pairing instances. Defaults to `app=astarte-pairing`.
- [astarte_pairing] Cluster with realm management using `CLUSTERING_STRATEGY` and `CLUSTERING_KUBERNETES_NAMESPACE`
- [astarte_pairing] Add realm-scoped health checks, which can be used to test astarte health or connectivity by devices who are only aware of realm-scoped paths, such as during FDO
- [astarte_housekeeping] support network topology replication strategy for the `astarte` keyspace, with the following env vars:
  - `HOUSEKEEPING_ASTARTE_KEYSPACE_REPLICATION_STRATEGY` - Replication strategy for the `astarte` keyspace: "SimpleStrategy" or "NetworkTopologyStrategy" (default: "SimpleStrategy")
  - `HOUSEKEEPING_ASTARTE_KEYSPACE_REPLICATION_FACTOR` - Replication factor when using SimpleStrategy (default: 1)
  - `HOUSEKEEPING_ASTARTE_KEYSPACE_NETWORK_REPLICATION_MAP` - Datacenter replication map when using NetworkTopologyStrategy (no default, required when using network strategy)
- Added database events handling configuration across all services:
  - `DATABASE_EVENTS_HANDLING_METHOD` - Controls how database events are handled: "expose" (via telemetry) or "log" (to logs) (default: "expose")
- [astarte_pairing] Added device registration triggers
- [astarte_realm_management] Added device deletion started and device deletion completed triggers
- Allow devices with empty introspection
- Devices can now declare support for optional Astarte MQTT v1 features to Astarte via capabilities
- Support for `purge_properties_compression_format` capability. possible values are `zlib` (default) and `plaintext`


### Changed

- [astarte_housekeeping] When `HOUSEKEEPING_ASTARTE_KEYSPACE_REPLICATION_STRATEGY` is not set, the `astarte` keyspace is now created using `NetworkTopologyStrategy` with a replication map derived from the current ScyllaDB network topology (one replica per node in each datacenter), instead of falling back to `SimpleStrategy` with replication factor 1.
- [astarte_housekeeping] When creating a realm without specifying its replication strategy, the strategy used in the `astarte` keyspace is applied as default.
- [astarte_pairing] FDO authentication (EXPERIMENTAL feature, disabled by default). New environment variables are needed in order to use FDO:
  - `PAIRING_ENABLE_FDO` - whether the FDO feature is enabled or not (default: false)
  - `PAIRING_FDO_RENDEZVOUS_URL` - URL of the rendezvous server (default: "http://rendezvous:8041")
  - `ASTARTE_BASE_URL_DOMAIN` - domain part of the base URL of astarte, used by devices to connect in TO2 phase (required if FDO enabled)
  - `ASTARTE_BASE_URL_PORT` - port of the base URL of astarte (required if FDO enabled)
  - `ASTARTE_BASE_URL_PROTOCOL` - protocol of the base URL of astarte (required if FDO enabled)
  - `PAIRING_FDO_RENDEZVOUS_SSL_ENABLED` - whether SSL is enabled for the connection to the rendezvous server (default: false)
  - `PAIRING_FDO_RENDEZVOUS_SSL_CA_FILE` - path to the CA certificate file for the rendezvous server TLS connection; when not specified, the bundled cURL certificate bundle will be used
  - `PAIRING_FDO_RENDEZVOUS_SSL_DISABLE_SNI` - disable Server Name Indication for the rendezvous server TLS connection (default: false)
  - `PAIRING_FDO_RENDEZVOUS_SSL_CUSTOM_SNI` - custom SNI hostname for the rendezvous server TLS connection; when not specified, the Rendezvous hostname will be used
- [astarte_housekeeping] AMQP management configuration moved to `astarte_events` library. Environment variables changed:
  - `HOUSEKEEPING_AMQP_SSL_ENABLED` is now `ASTARTE_EVENTS_AMQP_MANAGEMENT_SSL_ENABLED`
  - `HOUSEKEEPING_AMQP_SSL_CA_FILE` is now `ASTARTE_EVENTS_AMQP_MANAGEMENT_SSL_CA_FILE`
  - `HOUSEKEEPING_AMQP_SSL_DISABLE_SNI` is now `ASTARTE_EVENTS_AMQP_MANAGEMENT_SSL_DISABLE_SNI`
  - `HOUSEKEEPING_AMQP_SSL_CUSTOM_SNI` is now `ASTARTE_EVENTS_AMQP_MANAGEMENT_SSL_CUSTOM_SNI`
  - `HOUSEKEEPING_AMQP_HOST` is now `ASTARTE_EVENTS_AMQP_MANAGEMENT_HOST`
  - `HOUSEKEEPING_AMQP_MANAGEMENT_PORT` is now `ASTARTE_EVENTS_AMQP_MANAGEMENT_PORT`
  - `HOUSEKEEPING_AMQP_USERNAME` is now `ASTARTE_EVENTS_AMQP_MANAGEMENT_USERNAME`
  - `HOUSEKEEPING_AMQP_PASSWORD` is now `ASTARTE_EVENTS_AMQP_MANAGEMENT_PASSWORD`
 - Services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages
 - BREAKING: Merged API services into main services, eliminating separate containers:
  - `astarte-housekeeping-api` merged into `astarte-housekeeping`
  - `astarte-pairing-api` merged into `astarte-pairing`
  - `astarte-realm-management-api` merged into `astarte-realm-management`
- BREAKING: [astarte_housekeeping] Housekeeping now creates an AMQP vhost for each created realm. Some required configuration was introduced:
  - `HOUSEKEEPING_AMQP_HOST` - AMQP host for housekeeping operations (required in production, defaults to "localhost" in dev/test)
  - `HOUSEKEEPING_AMQP_SSL_ENABLED` - Enable SSL for AMQP connections (default: false)
  - `HOUSEKEEPING_AMQP_SSL_CA_FILE` - CA certificate file for AMQP SSL connections (default: bundled cURL certificate bundle)
  - `HOUSEKEEPING_AMQP_SSL_DISABLE_SNI` - Disable Server Name Indication (default: false)
  - `HOUSEKEEPING_AMQP_SSL_CUSTOM_SNI` - Custom SNI hostname (defaults to AMQP host if unset)
  - `HOUSEKEEPING_AMQP_USERNAME` - AMQP username (default: guest)
  - `HOUSEKEEPING_AMQP_PASSWORD` - AMQP password (default: guest)
  - `HOUSEKEEPING_AMQP_MANAGEMENT_PORT` - AMQP management API port (default: 15672)
- BREAKING: [astarte_data_updater_plant] Added a separate the AMQP producer configuration, which is mandatory:
  - `ASTARTE_EVENTS_PRODUCER_AMQP_HOST` - Host for producer connection (default: "localhost")
  - `ASTARTE_EVENTS_PRODUCER_AMQP_USERNAME` - Username for producer (default: "guest")
  - `ASTARTE_EVENTS_PRODUCER_AMQP_PASSWORD` - Password for producer (default: "guest")
  - `ASTARTE_EVENTS_PRODUCER_AMQP_VIRTUAL_HOST` - Virtual host for internal events (default: "/")
  - `ASTARTE_EVENTS_PRODUCER_AMQP_PORT` - Port for producer (default: 5672)
  - `ASTARTE_EVENTS_PRODUCER_AMQP_SSL_ENABLED` - Enable SSL for producer (default: false)
  - `ASTARTE_EVENTS_PRODUCER_AMQP_SSL_CA_FILE` - CA certificate file for producer SSL (default: bundled cURL certificates)
  - `ASTARTE_EVENTS_PRODUCER_AMQP_SSL_DISABLE_SNI` - Disable Server Name Indication for producer (default: false)
  - `ASTARTE_EVENTS_PRODUCER_AMQP_SSL_CUSTOM_SNI` - Custom SNI hostname for producer (falls back to the value of ASTARTE_EVENTS_PRODUCER_AMQP_HOST)
  - `ASTARTE_EVENTS_PRODUCER_AMQP_CONNECTION_NUMBER` - The number of open connections to RabbitMQ (default: 10)
  - `ASTARTE_EVENTS_PRODUCER_AMQP_DATA_QUEUE_TOTAL_COUNT` - The total number of data queues in the astarte cluster (default: 128)
  - `ASTARTE_EVENTS_PRODUCER_AMQP_EVENTS_EXCHANGE_NAME` - The exchange for internal events (default: "astarte_events")
- BREAKING: AMQP triggers are now always sent to an ad-hoc vhost for each realm. the vhost is `[astarte_instance_id]_[realm_name]`, which under normal circumstances is just `_[realm_name]`
- BREAKING: [astarte_pairing] AMQP Producer configuration is now mandatory using the `ASTARTE_EVENTS_PRODUCER_AMQP_*` environment variables
- BREAKING: [astarte_realm_management] AMQP Producer configuration is now mandatory using the `ASTARTE_EVENTS_PRODUCER_AMQP_*` environment variables
- More accurate health checks for astarte services


### Fixed

- [astarte_housekeeping] Changed logger config to make the metadata fields appear in logfmt log output; in particular the msg field
- [astarte_pairing] Changed logger config to make the metadata fields appear in logfmt log output; in particular the msg field
- Fix crashes in new ScyllaDB versions
- [astarte_appengine_api] Volatile trigger requests to Data Updater Plant now properly encode the volatile trigger type
- [astarte_data_updater_plant] Do not crash when transient triggers are installed on devices with outdated introspection
- [astarte_data_updater_plant] Correctly encode values when sending properties to device on connection
- [astarte_realm_management] Allow to delete long-disconnected devices. For this to work, Realm Management needs `CLUSTERING_STRATEGY`, `CLUSTERING_KUBERNETES_NAMESPACE` and `DATA_UPDATER_PLANT_CLUSTERING_KUBERNETES_SELECTOR` to be set, just like AppEngine and DUP. Refer to 1.2.1-rc.0 for additional information on the variables.

• Find Astrate pairing/trigger/properties code Explore Agent
No prior run is recorded, and the gap spans v1.2.2 → v1.4.0-rc.3. Strong wire-visible candidates from the notes: empty introspection, MQTT v1 capabilities (`purge_properties_compression_format: plaintext`), new trigger event types, and realm data retention. The recipe requires checking the diff for anything that matters — let me look at both upstream changes and Astrate's current surfaces.

$ gh api "repos/astarte-platform/astarte/commits?path=apps/astarte_pairing&since=2026-04-24T00:00:00Z&per_page=50" --jq '.[] | "\(.sha[0:7])\t\(.commit.message | split("\n")[0])"' 2>/dev/null | head -40
5047681	feat(dup): Encrypt device shared secret (#2104)
6e58d66	test(pairing): fix flakiness (#2100)
1c297aa	refactor: purge exrabbitpool (#2050)
6d0c0d6	chore(deps): bump httpoison and hackney
4424b08	chore(dup): update cose elixir dep
36a70c2	chore: forward port release-1.4
897a383	feat: Store shared secret as session_key
08d7d5e	chore: prepare v1.4.0-rc.2
c5fd781	chore: prepare v1.4.0-rc.1
9593d55	fix(fdo): Ensure data from previous session is deleted
9a4b375	chore: forward-port release-1.3
17b329a	chore: prepare v1.3.1
20a6ba9	chore(dup): persist shared secret
d0fa44e	feat(appengine): Decrypt stored values (#2029)
f347184	test(pairing): renew expired self-signed certificate used in credentials tests
d8ab739	test(pairing): renew expired self-signed certificate used in credentials tests
5938a8a	test(pairing): renew expired self-signed certificate used in credentials tests
95d22ed	chore: forward port release-1.4
815add7	refactor: use ad-hoc table for unconfirmed devices
43de307	fix(fdo): delete devices from previous sessions (#2028)
cebf13f	chore(dup): add cache for storage of data encryption keys (DEKs)
a540858	chore(dup): add InitExchange logic for Astarte-Device handshake
11a8de1	refactor(data_access)!: obliterate `astarte_data_access_generator` (#2003)
8f87e80	chore(deps): update protobuf to 0.17.0 (#2016) (#2019)
036172b	fix(housekeeping): Fix realm creation (#2018)
e7401ec	chore(deps): update protobuf to 0.17.0 (#2016)
980e487	chore: forward-port release-1.3
30bfbf7	fix: Remove vulnerable cowlib dependency (#2006)
9357d52	ci: Alert for security vulnerabilities in app/lib dependencies (#1978)
ab334e4	chore(data_access): remove warning
8313325	chore: forward port release-1.4
3e8c8f9	test(pairing): Fix failing tests (#1997)
6cc81bc	chore: add data_access support for encrypted endpoints (#1993)
963829c	style: fix formatting for mix compliance
9475a72	build(docker): properly include local dependencies
c9ca09f	chore: forward port release-1.4
5f11907	chore: forward port release-1.3
be65367	docs(appengine): Fix swagger endpoint
08963e5	chore: prepare 1.5 dev cycle
528489d	chore: forward-port release-1.3

✓ Find Astrate pairing/trigger/properties code Explore Agent
Astrate already implements some of this (empty introspection accepted, `purge_properties_compression_format` capability outbound, `device_registered` trigger, global retention config). Let me verify the upstream diffs for the items that matter, and check issue #21's status.

$ gh issue view 21 --json title,state,labels 2>/dev/null; echo ---; gh issue view 21 --json body --jq '.body' 2>/dev/null | head -20
{"labels":[{"id":"LA_kwDORmfsJs8AAAACte8Saw","name":"mule-review","description":"Mule pushed a change for this; waiting for a human or a strong model to review it","color":"0e8a16"}],"state":"CLOSED","title":"device-deletion-events: emit device_deletion_started + device_deletion_finished around synchronous delete"}
---
Astrate performs synchronous device deletion (`COMPATIBILITY.md` deviation 9), while upstream
fires `device_deletion_started` then `device_deletion_finished` around an async deletion
(`060-triggers.md:109-110`). Currently Astrate emits neither event, so a trigger config
imported from a real Astarte setup silently never fires on device deletion.

**Decision (Giulio, 2026-07-27):** emit both events back-to-back around the synchronous
delete, even though there's no real gap between them here — the goal is wire/behavioral
compatibility for anyone plugging in an existing trigger config, not literal fidelity to
upstream's async timing.

**Work needed**: find the device-deletion path (`internal/realm/service.go` /
`internal/realm/http.go` / `internal/store/devices.go` handle `DeleteDevice`-shaped calls),
emit `device_deletion_started` immediately before the delete and `device_deletion_finished`
immediately after (same event/payload shape as upstream's SimpleEvent, matching how
`device_connected`/`device_disconnected` are already emitted in the trigger engine). Add a
test asserting both fire in order around a delete.

Found by a cross-project survey (`.mule/research/survey-2026-07-27.md`, source 4).

$ gh api "repos/astarte-platform/astarte/commits?since=2026-04-24T00:00:00Z&per_page=100" --jq '.[] | "\(.sha[0:7])\t\(.commit.message | split("\n")[0])"' 2>/dev/null | grep -iE "introspection|capabil|empty|registration|retention|purge_properties|compression" 
2a9f009	fix: Solve capabilities message stalling queues

$ git branch --show-current; git log --oneline -3 2>/dev/null; echo ---; git status -s | head
mule/queue
4beddf0 mule: issue-13 passed on beecab8
beecab8 mule: issue-20 passed on 33d5be6
33d5be6 mule: log issue-12
---
 M .mule/for-giulio.md
?? .mule/cron.log
?? tools/mule-plan-day.sh

$ rg -l "device_deletion_started|device_deletion_finished|DeviceDeletion" --glob '!test/**' --glob '!docs/**' --glob '!.mule/**' . 2>/dev/null | head
./internal/realm/dashboard_compat_test.go

→ Read internal/engine/capabilities.go
On `mule/queue`; issue #21 (deletion events) is closed with `mule-review` — already in flight, so not for me to re-propose. Let me verify the upstream diffs for the remaining candidates and read Astrate's capabilities handling.

$ gh search prs --repo astarte-platform/astarte "empty introspection" --limit 10 --json number,title,state,mergedAt --jq '.[] | "\(.number)\t\(.title)\t\(.state)"' 2>/dev/null; echo ===; gh search prs --repo astarte-platform/astarte "capabilities" --limit 15 --json number,title,state,mergedAt --jq '.[] | "\(.number)\t\(.title)\t\(.state)"' 2>/dev/null
===

→ Read internal/engine/control.go [offset=1, limit=120]
Astrate already implements the capabilities mechanism, empty-introspection acceptance, and `device_registered`. Let me pin down the two open candidates (upstream registration-trigger event name, realm data-retention surface) and check Astrate's control-payload directions.

$ gh api "repos/astarte-platform/astarte/commits?since=2026-04-20T00:00:00Z&per_page=200&path=apps/astarte_trigger_engine" --jq '.[] | "\(.sha[0:7])\t\(.commit.message | split("\n")[0])"' 2>/dev/null | head -30
1c297aa	refactor: purge exrabbitpool (#2050)
6d0c0d6	chore(deps): bump httpoison and hackney
4424b08	chore(dup): update cose elixir dep
36a70c2	chore: forward port release-1.4
897a383	feat: Store shared secret as session_key
08d7d5e	chore: prepare v1.4.0-rc.2
c5fd781	chore: prepare v1.4.0-rc.1
9a4b375	chore: forward-port release-1.3
17b329a	chore: prepare v1.3.1
20a6ba9	chore(dup): persist shared secret
d0fa44e	feat(appengine): Decrypt stored values (#2029)
a540858	chore(dup): add InitExchange logic for Astarte-Device handshake
05115d5	chore: forward port release-1.4
11a8de1	refactor(data_access)!: obliterate `astarte_data_access_generator` (#2003)
8f87e80	chore(deps): update protobuf to 0.17.0 (#2016) (#2019)
e7401ec	chore(deps): update protobuf to 0.17.0 (#2016)
30bfbf7	fix: Remove vulnerable cowlib dependency (#2006)
9357d52	ci: Alert for security vulnerabilities in app/lib dependencies (#1978)
6cc81bc	chore: add data_access support for encrypted endpoints (#1993)
9475a72	build(docker): properly include local dependencies
08963e5	chore: prepare 1.5 dev cycle
528489d	chore: forward-port release-1.3
26c9a5b	chore: prepare v1.3.0

$ gh api "repos/astarte-platform/astarte/commits?since=2026-04-24T00:00:00Z&until=2026-05-20T00:00:00Z&per_page=100" --jq '.[] | "\(.sha[0:7])\t\(.commit.message | split("\n")[0])"' 2>/dev/null | grep -viE "chore|deps|test:|ci|merge|forward|refactor|build|revert|docs" | head -60
bcfbd83	feat(adapters): introduce `transform` macro
a7718be	feat(housekeeping): default realm keyspace replication uses data from db
6523b2a	feat(housekeeping): default astarte keyspace replication to NetworkTopologyStrategy

$ gh api "repos/astarte-platform/astarte/git/trees/v1.3.0?recursive=1" --jq '.tree[].path' 2>/dev/null | grep -iE "trigger|pairing" | head -40
.github/workflows/astarte-pairing-workflow.yaml
.github/workflows/astarte-trigger-engine-workflow.yaml
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/delete_volatile_trigger
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/delete_volatile_trigger/request_data.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/install_volatile_trigger
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/install_volatile_trigger/request_data.ex
apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/volatile_trigger.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/data_trigger.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/trigger.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/trigger_policy
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/trigger_policy/queries.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/triggers_handler.ex
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_trigger_test.exs
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/trigger_handler_test.exs
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/trigger_test.exs
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/triggers_handler_test.exs
apps/astarte_data_updater_plant/test/support/cases/trigger.ex
apps/astarte_pairing
apps/astarte_pairing/.formatter.exs
apps/astarte_pairing/Dockerfile
apps/astarte_pairing/README.md
apps/astarte_pairing/config
apps/astarte_pairing/config/ci.exs
apps/astarte_pairing/config/config.exs
apps/astarte_pairing/config/dev.exs
apps/astarte_pairing/config/prod.exs
apps/astarte_pairing/config/releases.exs
apps/astarte_pairing/config/test.exs
apps/astarte_pairing/lib
apps/astarte_pairing/lib/astarte_pairing.ex
apps/astarte_pairing/lib/astarte_pairing
apps/astarte_pairing/lib/astarte_pairing/agent
apps/astarte_pairing/lib/astarte_pairing/agent/agent.ex
apps/astarte_pairing/lib/astarte_pairing/agent/device_registration_request.ex
apps/astarte_pairing/lib/astarte_pairing/agent/device_registration_response.ex
apps/astarte_pairing/lib/astarte_pairing/agent/realm.ex
apps/astarte_pairing/lib/astarte_pairing/application.ex
apps/astarte_pairing/lib/astarte_pairing/auth
apps/astarte_pairing/lib/astarte_pairing/auth/auth.ex
apps/astarte_pairing/lib/astarte_pairing/auth/user.ex

$ gh search code --repo astarte-platform/astarte "device_registered" --json path,textMatches --jq '.[] | "\(.path)\t\(.textMatches[0].fragment | gsub("\\s+";" ") | .[0:160])"' 2>/dev/null | head -30
libs/astarte_generators/lib/astarte/core/generators/triggers/simple_events/device_registered_event.ex	 use ExUnitProperties alias Astarte.Core.Triggers.SimpleEvents.DeviceRegisteredEvent @spec device_registered_event() :: StreamData.t(DeviceRegisteredEvent.t()) 
apps/astarte_pairing/lib/astarte_pairing/agent/agent.ex	 defp dispatch_device_registration_trigger(realm_name, hw_id) do timestamp = DateTime.utc_now() |> DateTime.to_unix(:millisecond) event_key = :on_device_registe
apps/astarte_pairing/lib/astarte_pairing/application.ex	 children = [ Astarte.PairingWeb.Telemetry, {Astarte.Pairing.CredentialsSecret.Cache, []}, {Astarte.RPC.Triggers.Client, types: [:DEVICE_REGISTERED]}, Astarte.P
doc/pages/architecture/060-triggers.md	- `device_error`: triggered when data from a device causes an error. - `device_registered`: triggered when a device is registered to Astarte. - `device_deletion
astarte-dashboard/src/astarte-client/types/dto/trigger.d.ts	 on: | 'device_disconnected' | 'device_connected' | 'device_registered' | 'device_deletion_finished' | 'device_deletion_started' | 'device_error'
astarte-dashboard/src/astarte-client/types/events/AstarteDeviceRegistrationEvent.ts	 event: { type: 'device_registered'; };
libs/astarte_events/lib/astarte_events/triggers/core.ex	 | {:on_interface_minor_updated, interface()} | :on_device_registered | :on_device_deletion_started
astarte-dashboard/src/astarte-client/models/Trigger/index.ts	 | 'device_connected' | 'device_registered' | 'device_deletion_finished'
astarte-dashboard/cypress/e2e/trigger_builder.cy.js	const triggerConditionToLabel = { device_connected: 'Device Connected', device_disconnected: 'Device Disconnected', device_registered: 'Device Registered', devi
libs/astarte_generators/lib/astarte/core/generators/triggers/simple_events/simple_event.ex	 DeviceDeletionStartedEventGenerator.device_deletion_started_event()}, {:device_deletion_finished_event, DeviceDeletionFinishedEventGenerator.device_deletion_fi
astarte-dashboard/src/components/TriggerEditor/SimpleTriggerForm.tsx	 device_connected: 'Device Connected', device_registered: 'Device Registered', device_deletion_finished: 'Device Deletion Finished',
libs/astarte_generators/test/astarte/core/generators/triggers/simple_events/device_registered_event_test.exs	 @moduletag :simple_event @moduletag :device_registered_event 
apps/astarte_pairing/test/support/helpers/triggers.ex	 def register_device_registration_trigger(realm_name, conditions) do register_device_trigger( realm_name, "device_registered", :device_registered_event, conditi
libs/astarte_events/test/astarte_events/triggers/core_test.exs	 :on_interface_minor_updated end test "maps DEVICE_REGISTERED to :on_device_registered" do assert Core.pretty_device_event_type(:DEVICE_REGISTERED) == :on_devic

$ gh api "repos/astarte-platform/astarte/contents/apps/astarte_pairing/lib/astarte_pairing/agent/agent.ex?ref=v1.2.2" --jq '.content' 2>/dev/null | base64 -d 2>/dev/null | grep -n "device_registration_trigger\|register_trigger\|on_device_registered" | head
(no output)

$ gh search code --repo astarte-platform/astarte "dispatch_device_registration_trigger" --json path --jq '.[] | .path' 2>/dev/null | head; echo === v1.2.2 tree check ===; gh api "repos/astarte-platform/astarte/git/trees/v1.2.2?recursive=1" --jq '.tree[].path' 2>/dev/null | grep -c "astarte_pairing/lib"
apps/astarte_pairing/lib/astarte_pairing/agent/agent.ex
=== v1.2.2 tree check ===
22

`device_registered` trigger was indeed new in v1.3.0 and Astrate already emits it (issue #21 already covers the deletion pair). Now let me verify the remaining candidates against upstream code.
$ gh search code --repo astarte-platform/astarte "purge_properties_compression_format" --json path,textMatches --jq '.[] | "\(.path)\t\(.textMatches[0].fragment | gsub("\\s+";" ") | .[0:140])"' 2>/dev/null | head -20
CHANGELOG.md	- [astarte_realm_management] Added device deletion started and device deletion completed triggers - Allow devices with empty introspection -
tools/astarte_import/lib/astarte/import.ex	 end defp xml_event({:startElement, _uri, _l_name, {_prefix, ~c"capabilities"}, attr}, _loc, state) do with {:ok, purge_properties_compressi
doc/pages/architecture/090-database.md	 ```sql CREATE TYPE <realm name>.capabilities ( purge_properties_compression_format int ); ``` 
doc/pages/architecture/080-mqtt-v1-protocol.md	 | Name | Values | Purpose | | ------------------------------------- | ----------------------- | -------------------------------------------
libs/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0008_create_capabilities_type.ex	 use Ecto.Migration def up do execute "CREATE TYPE #{prefix()}.capabilities (purge_properties_compression_format int)" end def down do
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex	 defp decode_payload(%State{capabilities: capabilities} = _state, payload) do case Map.get(capabilities, :purge_properties_compression_forma
tools/astarte_export/test/astarte_export_test.exs	 <protocol revision="0" pending_empty_cache="false"/> <registration credentials_secret="$2b$12$bKly9EEKmxfVyDeXjXu1vOebWgr34C8r4IHd9Cd.34Ozm
tools/astarte_export/test/database_test.exs	 CREATE TYPE test.capabilities ( purge_properties_compression_format int );
tools/astarte_import/test/astarte/import_test.exs	 last_credentials_request_ip="198.51.100.1" /> <capabilities purge_properties_compression_format="0" /> <stats total_received_msgs="64" tota
tools/astarte_import/test/astarte/populatedb_test.exs	 cert_aki="a8eaf08a797f0b10bb9e7b5dca027ec2571c5ea6" first_credentials_request="2019-05-30T13:49:57.355Z" last_credentials_request_ip="198.5
libs/astarte_rpc/test/support/helpers/database.ex	 @create_capabilities_type """ CREATE TYPE :keyspace.capabilities ( purge_properties_compression_format int ); """ 
libs/astarte_fdo/test/support/helpers/database.ex	 @create_capabilities_type """ CREATE TYPE :keyspace.capabilities ( purge_properties_compression_format int ); """ 
apps/astarte_trigger_engine/test/support/helpers/database.ex	 @create_capabilities_type """ CREATE TYPE :keyspace.capabilities ( purge_properties_compression_format int ); """ 
tools/astarte_import/test/astarte/database_test.exs	 @create_capabilities_type """ CREATE TYPE test.capabilities ( purge_properties_compression_format int ); """ 
apps/astarte_pairing/test/support/helpers/database.ex	 @create_capabilities_type """ CREATE TYPE :keyspace.capabilities ( purge_properties_compression_format int ); """ 
apps/astarte_appengine_api/test/support/helpers/database_v2.ex	 @create_capabilities_type """ CREATE TYPE :keyspace.capabilities ( purge_properties_compression_format int ); """ 
apps/astarte_realm_management/test/support/helpers/database.ex	 @create_capabilities_type """ CREATE TYPE IF NOT EXISTS :keyspace.capabilities ( purge_properties_compression_format int ); """ 
libs/astarte_data_access/test/support/helpers/database.ex	 @create_capabilities_type """ CREATE TYPE :keyspace.capabilities ( purge_properties_compression_format int ); """ 
apps/astarte_appengine_api/test/support/helpers/database.ex	 @create_capabilities_type """ CREATE TYPE #{Realm.keyspace_name(@test_realm)}.capabilities ( purge_properties_compression_format int ); """
libs/astarte_events/test/support/helpers/database_test_helper.ex	 @create_capabilities_type """ CREATE TYPE :keyspace.capabilities ( purge_properties_compression_format int ); """

$ gh search code --repo astarte-platform/astarte "DATASTREAM_MAXIMUM_STORAGE_RETENTION" --json path,textMatches --jq '.[] | "\(.path)\t\(.textMatches[0].fragment | gsub("\\s+";" ") | .[0:140])"' 2>/dev/null | head -20
CHANGELOG.md	 - [astarte_housekeeping] Added HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION to set default, instance wise, realm default realm
apps/astarte_housekeeping/lib/astarte_housekeeping/config.ex	 Can be overridden at any time during realm creation or update [Check docs for more info](https://docs.astarte-platform.org/astarte/latest/0
apps/astarte_housekeeping/lib/astarte_housekeeping/realms/core.ex	 Queries.set_device_registration_limit(realm.realm_name, limit) end defp update_realm_field(realm, :datastream_maximum_storage_retention, :u
apps/astarte_housekeeping/lib/astarte_housekeeping/realms/realm.ex	 :device_registration_limit, :datastream_maximum_storage_retention | @required_create_fields
apps/astarte_housekeeping/lib/astarte_housekeeping/realms/realms.ex	 defp maybe_inject_default_retention(attrs) do case Config.default_datastream_maximum_storage_retention() do {:ok, nil} ->
apps/astarte_realm_management/lib/astarte_realm_management_web/router.ex	 get "/config/datastream_maximum_storage_retention", RealmConfigController, :show_datastream_maximum_storage_retention resources "/triggers"
apps/astarte_housekeeping/lib/astarte_housekeeping_web/views/error_view.ex	 %{errors: %{detail: "Realm already exists"}} end def render("set_datastream_maximum_storage_retention_fail.json", _assigns) do %{errors: %{
apps/astarte_housekeeping/lib/astarte_housekeeping_web/views/realm_view.ex	 replication_class: "SimpleStrategy", replication_factor: realm.replication_factor, device_registration_limit: realm.device_registration_lim
doc/pages/architecture/090-database.md	 ```sql INSERT INTO <realm name>.kv_store (group, key, value) VALUES ('realm_config', 'datastream_maximum_storage_retention', intAsBlob(<max
apps/astarte_housekeeping/lib/astarte_housekeeping_web/api_spec/schemas/realm.ex	 }, datastream_maximum_storage_retention: %Schema{ type: :integer,
apps/astarte_realm_management/lib/astarte_realm_management/realm_config/queries.ex	 Returns either `{:ok, limit}` or `{:error, reason}`. The limit is a strictly positive integer (if set), 0 if unset. """ @spec get_datastrea
apps/astarte_housekeeping/lib/astarte_housekeeping/realms/queries.ex	 {:ok, device_registration_limit} <- fetch_device_registration_limit(realm_name) do max_retention = get_datastream_maximum_storage_retention
apps/astarte_realm_management/lib/astarte_realm_management/realm_config/realm_config.ex	 end def get_datastream_maximum_storage_retention(realm_name) do Queries.get_datastream_maximum_storage_retention(realm_name) end def update
apps/astarte_realm_management/lib/astarte_realm_management_web/views/realm_config_view.ex	 } end def render("show.json", %{datastream_maximum_storage_retention: ttl}) do %{ data: ttl }
apps/astarte_housekeeping/lib/astarte_housekeeping_web/controllers/fallback_controller.ex	 def call(conn, {:error, :delete_datastream_maximum_storage_retention_fail}) do conn
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/time_based_actions.ex	 |> reload_device_deletion_status_on_expiry(timestamp) |> reload_datastream_maximum_storage_retention_on_expiry(timestamp) end
apps/astarte_realm_management/lib/astarte_realm_management_web/controllers/realm_config_controller.ex	 operation :show_datastream_maximum_storage_retention, summary: "Get datastream maximum storage retention",
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/impl.ex	 # TODO this could be a bang! {:ok, ttl} = Queries.get_datastream_maximum_storage_retention(state.realm) 
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/state.ex	 field :last_seen_message, non_neg_integer(), default: 0 field :last_device_triggers_refresh, non_neg_integer(), default: 0 field :last_grou
apps/astarte_appengine_api/lib/astarte_appengine_api/device/queries.ex	 def fetch_datastream_maximum_storage_retention(realm_name) do keyspace = Realm.keyspace_name(realm_name)

Realm-scoped datastream retention (`datastream_maximum_storage_retention`) is real and new — a Realm Management `/config/datastream_maximum_storage_retention` endpoint plus housekeeping field. Astrate only has global retention. Let me confirm when it appeared, its verbs, and whether the purge-properties capability affects the device→server direction.
$ gh api "repos/astarte-platform/astarte/contents/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex?ref=v1.4.0-rc.3" --jq '.content' 2>/dev/null | base64 -d 2>/dev/null | grep -n "decode_payload\|purge\|compression\|capabilit" | head -20
68:    decompression_start = System.monotonic_time()
70:    with {:ok, decoded_payload} <- decode_payload(state, payload),
72:      # Track successful decompression
74:        [:astarte, :data_updater_plant, :control_handler, :payload_decompression],
76:          duration: System.monotonic_time() - decompression_start,
104:        # Track failed decompression
106:          [:astarte, :data_updater_plant, :control_handler, :payload_decompression],
108:            duration: System.monotonic_time() - decompression_start,
122:          message: "Invalid purge_properties payload",
123:          logger_metadata: [tag: "purge_properties_error"],
124:          error_name: "purge_properties_error",
125:          error: :purge_properties_error
172:  defp decode_payload(%State{capabilities: capabilities} = _state, payload) do
173:    case Map.get(capabilities, :purge_properties_compression_format) do
214:    compression_format = state.capabilities.purge_properties_compression_format
221:           compression_format
229:  defp send_consumer_properties_payload(realm, device_id, abs_paths_list, compression_format) do
235:      case compression_format do

$ gh api "repos/astarte-platform/astarte/contents/apps/astarte_realm_management/lib/astarte_realm_management_web/router.ex?ref=v1.4.0-rc.3" --jq '.content' 2>/dev/null | base64 -d 2>/dev/null | grep -n "datastream\|config" | head; echo === v1.3.0 router ===; gh api "repos/astarte-platform/astarte/contents/apps/astarte_realm_management/lib/astarte_realm_management_web/router.ex?ref=v1.3.0" --jq '.content' 2>/dev/null | base64 -d 2>/dev/null | grep -n "datastream" | head
39:    get "/config/auth", RealmConfigController, :show_auth
40:    put "/config/auth", RealmConfigController, :update_auth
42:    get "/config/device_registration_limit",
46:    get "/config/datastream_maximum_storage_retention",
48:        :show_datastream_maximum_storage_retention
=== v1.3.0 router ===
46:    get "/config/datastream_maximum_storage_retention",
48:        :show_datastream_maximum_storage_retention

$ rg -n "config/|datastream_maximum|device_registration_limit|config/auth" internal/realm internal/realmmanagement internal/ 2>/dev/null | head -20
internal/realm/http.go:45:	mux.Handle("GET /realmmanagement/v1/{realm}/config/auth", h(a.getAuth))
internal/realm/http.go:46:	mux.Handle("PUT /realmmanagement/v1/{realm}/config/auth", h(a.putAuth))
internal/realm/http.go:47:	mux.Handle("GET /realmmanagement/v1/{realm}/config/device_registration_limit", h(a.getRegistrationLimit))
internal/realm/http.go:246:// --- config/auth ------------------------------------------------------------
internal/realm/http.go:248:// authConfig is the GET/PUT /config/auth body shape (upstream
internal/realm/dashboard_compat_test.go:26:// Dashboard v1.2.2 requires: version, device_registration_limit, delivery
internal/realm/dashboard_compat_test.go:41:		rec := r.req(t, http.MethodGet, "/config/device_registration_limit", "", r.rmaToken)
internal/realm/http_test.go:239:		if rec := r.req(t, http.MethodPut, "/config/auth", `{"jwt_public_key_pem":`+jsonStr(rotated)+`}`, r.rmaToken); rec.Code != http.StatusNoContent {
internal/realm/http_test.go:240:			t.Fatalf("put config/auth: got %d, want 204 (%s)", rec.Code, rec.Body)
internal/realm/http_test.go:243:		decodeData(t, r.req(t, http.MethodGet, "/config/auth", "", r.rmaToken), &cfg)
internal/realm/http_test.go:245:			t.Errorf("config/auth key not rotated")
internal/realm/service.go:412:// (nil = unlimited), upstream GET /config/device_registration_limit.
internal/realm/service.go:421:// --- config/auth ------------------------------------------------------------
internal/realm/service.go:424:// /config/auth → {"jwt_public_key_pem": "..."}). Astrate stores a list for
internal/realm/service.go:435:// SetAuthKey rotates the realm's JWT public key (upstream PUT /config/auth).
internal/config/config.go:136:	DeviceRegistrationLimit *int32 `toml:"device_registration_limit"`
internal/housekeeping/http.go:45:	DeviceRegistrationLimit *int32 `json:"device_registration_limit"`
internal/config/config.example.toml:94:# device_registration_limit = 0
internal/housekeeping/service.go:81:		return nil, fmt.Errorf("%w: device_registration_limit must be non-negative", ErrValidation)
internal/store/devices.go:186:// device_registration_limit in the pairing layer).

The retention endpoint predates v1.4.0 (present at v1.3.0) — the v1.4.0 addition is the instance-wide default. The capability clearly affects both directions: upstream `decode_payload` (device→server) honors the compression format. Let me confirm the Astrate realm-config surface and read upstream's control handler.
$ gh api "repos/astarte-platform/astarte/contents/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex?ref=v1.4.0-rc.3" --jq '.content' 2>/dev/null | base64 -d 2>/dev/null | sed -n '120,260p'
error = %{
          message: "Invalid purge_properties payload",
          logger_metadata: [tag: "purge_properties_error"],
          error_name: "purge_properties_error",
          error: :purge_properties_error
        }

        opts = [execute_error_triggers: false, update_stats: false]

        Core.Error.handle_error(context, error, opts)
    end
  end

  def handle_control(state, "/emptyCache", _payload, timestamp) do
    state = TimeBasedActions.execute_time_based_actions(state, timestamp)

    with :ok <- send_control_consumer_properties(state, timestamp),
         {:ok, state} <- resend_all_properties(state, timestamp),
         :ok <- set_pending_empty_cache(state, timestamp) do
      :telemetry.execute(
        [:astarte, :data_updater_plant, :data_updater, :processed_empty_cache],
        %{},
        %{realm: state.realm}
      )

      {:ack, :ok, state}
    end
  end

  def handle_control(state, path, payload, timestamp) do
    # Track unexpected control messages
    :telemetry.execute(
      [:astarte, :data_updater_plant, :control_handler, :unexpected_control],
      %{payload_size: byte_size(payload)},
      %{realm: state.realm, control_path: path}
    )

    context = %{state: state, payload: payload, path: path, timestamp: timestamp}

    error = %{
      message:
        "Unexpected control on #{path}, base64-encoded payload: #{inspect(Base.encode64(payload))}",
      logger_metadata: [tag: "unexpected_control_message"],
      error_name: "unexpected_control_message",
      error: :unexpected_control_message,
      telemetry_name: [:astarte, :data_updater_plant, :data_updater, :discarded_control_message]
    }

    Core.Error.handle_error(context, error)
  end

  defp decode_payload(%State{capabilities: capabilities} = _state, payload) do
    case Map.get(capabilities, :purge_properties_compression_format) do
      :zlib ->
        <<_size_header::size(32), zlib_payload::binary>> = payload
        PayloadsDecoder.safe_inflate(zlib_payload)

      :plaintext ->
        {:ok, payload}
    end
  end

  defp send_control_consumer_properties(state, timestamp) do
    Logger.debug("Device introspection: #{inspect(state.introspection)}.")

    abs_paths_list =
      Enum.flat_map(state.introspection, fn {interface, _} ->
        descriptor = Map.get(state.interfaces, interface)

        case Core.Interface.maybe_handle_cache_miss(descriptor, interface, state) do
          {:ok, interface_descriptor, new_state} ->
            # Track successful interface loading
            :telemetry.execute(
              [:astarte, :data_updater_plant, :control_handler, :interface_loading],
              %{},
              %{realm: state.realm, interface: interface, result: "success"}
            )

            Core.Interface.gather_interface_property_paths(new_state.realm, interface_descriptor)

          {:error, :interface_loading_failed} ->
            # Track failed interface loading
            :telemetry.execute(
              [:astarte, :data_updater_plant, :control_handler, :interface_loading],
              %{},
              %{realm: state.realm, interface: interface, result: "failed"}
            )

            Logger.warning("Failed #{interface} interface loading.")
            []
        end
      end)

    compression_format = state.capabilities.purge_properties_compression_format

    # TODO: use the returned byte count in stats
    case send_consumer_properties_payload(
           state.realm,
           state.device_id,
           abs_paths_list,
           compression_format
         ) do
      {:ok, _bytes} -> :ok
      {:error, :session_not_found} -> session_not_found_error(state, timestamp)
      {:error, reason} -> generic_error(state, timestamp, reason)
    end
  end

  defp send_consumer_properties_payload(realm, device_id, abs_paths_list, compression_format) do
    topic = "#{realm}/#{Device.encode_device_id(device_id)}/control/consumer/properties"

    uncompressed_payload = Enum.join(abs_paths_list, ";")

    payload =
      case compression_format do
        :zlib ->
          payload_size = byte_size(uncompressed_payload)
          compressed_payload = :zlib.compress(uncompressed_payload)
          <<payload_size::unsigned-big-integer-size(32), compressed_payload::binary>>

        :plaintext ->
          uncompressed_payload
      end

    publish_start = System.monotonic_time()

    case VMQPlugin.publish(topic, payload, 2) do
      {:ok, %{local_matches: local, remote_matches: remote}} when local + remote == 1 ->
        # Track successful publish
        :telemetry.execute(
          [:astarte, :data_updater_plant, :control_handler, :vmq_publish],
          %{
            duration: System.monotonic_time() - publish_start,
            payload_size: byte_size(payload)
          },
          %{realm: realm, result: "success", matches: "single"}
        )

        {:ok, byte_size(topic) + byte_size(payload)}

Confirmed: upstream decodes the device→server purge-properties payload per the capability (plaintext allowed). Let me verify the v1.2.2 baseline for the retention endpoint, the producer/properties topic registration, and the "show default values" fix.
$ gh api "repos/astarte-platform/astarte/contents/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex?ref=v1.4.0-rc.3" --jq '.content' 2>/dev/null | base64 -d 2>/dev/null | sed -n '1,60p'
#
# This file is part of Astarte.
#
# Copyright 2025 SECO Mind Srl
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# SPDX-License-Identifier: Apache-2.0
#

defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.ControlHandler do
  @moduledoc """
  This module is responsible for handling the control messages.
  """
  alias Astarte.DataUpdaterPlant.DataUpdater.Core

  alias Astarte.Core.Device
  alias Astarte.DataUpdaterPlant.DataUpdater.Core
  alias Astarte.DataUpdaterPlant.DataUpdater.PayloadsDecoder
  alias Astarte.DataUpdaterPlant.DataUpdater.Queries
  alias Astarte.DataUpdaterPlant.DataUpdater.State
  alias Astarte.DataUpdaterPlant.RPC.VMQPlugin
  alias Astarte.DataUpdaterPlant.TimeBasedActions

  require Logger

  def handle_control(%State{discard_messages: true} = state, _, _, _) do
    {:discard, :discard_messages, state}
  end

  def handle_control(state, "/producer/properties", <<0, 0, 0, 0>>, timestamp) do
    new_state = TimeBasedActions.execute_time_based_actions(state, timestamp)

    timestamp_ms = div(timestamp, 10_000)

    :ok = Core.Device.prune_device_properties(new_state, "", timestamp_ms)

    final_state = %{
      new_state
      | total_received_msgs: new_state.total_received_msgs + 1,
        total_received_bytes:
          new_state.total_received_bytes + byte_size(<<0, 0, 0, 0>>) +
            byte_size("/producer/properties")
    }

    {:ack, :ok, final_state}
  end

  def handle_control(state, "/producer/properties", payload, timestamp) do
    start_time = System.monotonic_time()

$ gh api "repos/astarte-platform/astarte/contents/apps/astarte_realm_management/lib/astarte_realm_management_web/router.ex?ref=v1.2.2" --jq '.content' 2>/dev/null | base64 -d 2>/dev/null | grep -n "config" | head; echo === rc3 router config block ===; gh api "repos/astarte-platform/astarte/contents/apps/astarte_realm_management/lib/astarte_realm_management_web/router.ex?ref=v1.4.0-rc.3" --jq '.content' 2>/dev/null | base64 -d 2>/dev/null | sed -n '35,60p'
=== rc3 router config block ===
    resources "/interfaces", InterfaceController, only: [:index, :create]
    get "/interfaces/:interface_name/:major_version", InterfaceController, :show
    put "/interfaces/:interface_name/:major_version", InterfaceController, :update
    delete "/interfaces/:interface_name/:major_version", InterfaceController, :delete
    get "/config/auth", RealmConfigController, :show_auth
    put "/config/auth", RealmConfigController, :update_auth

    get "/config/device_registration_limit",
        RealmConfigController,
        :show_device_registration_limit

    get "/config/datastream_maximum_storage_retention",
        RealmConfigController,
        :show_datastream_maximum_storage_retention

    resources "/triggers", TriggerController,
      except: [:new, :edit, :update],
      param: "trigger_name"

    resources "/policies", TriggerPolicyController,
      except: [:new, :edit, :update],
      param: "policy_name"

    delete "/devices/:device_id", DeviceController, :delete
  end

v1.2.2 grep returned nothing — need to confirm whether the `/config/*` block even existed there, and find the "show default values" change.
$ gh api "repos/astarte-platform/astarte/contents/apps/astarte_realm_management/lib/astarte_realm_management_web/router.ex?ref=v1.2.2" --jq '.content' 2>/dev/null | base64 -d 2>/dev/null | grep -nE "config|datastream|device_registration|auth" | head; echo "=== exit: $? ==="
=== exit: 0 ===

$ gh search commits --repo astarte-platform/astarte "default values" --json sha,commit --jq '.[] | "\(.sha[0:7])\t\(.commit.message | split("\n")[0])"' 2>/dev/null | head -5; echo ===; gh api "repos/astarte-platform/astarte/commits?since=2026-07-01T00:00:00Z&until=2026-07-20T00:00:00Z&path=apps/astarte_realm_management&per_page=50" --jq '.[] | "\(.sha[0:7])\t\(.commit.message | split("\n")[0])"' 2>/dev/null | head -20
e98ca77	fix: Allow default values for simple strategy
3f0b864	fix(rm): Properly show default values for all mapping parameters (#2061)
4c1884b	fix: Use network topology as default realm replication
55fc84c	fix(interfaces): resolve explicit_timestamp injection and modal form submission-m
c6a4c38	feat(housekeeping): add default datastream maximum storage retention for realms
===
4424b08	chore(dup): update cose elixir dep
36a70c2	chore: forward port release-1.4
897a383	feat: Store shared secret as session_key
08d7d5e	chore: prepare v1.4.0-rc.2
c5fd781	chore: prepare v1.4.0-rc.1
9a4b375	chore: forward-port release-1.3
17b329a	chore: prepare v1.3.1
20a6ba9	chore(dup): persist shared secret
c6cf7fc	fix: ensure devices are eventually deleted (#2064)
d0fa44e	feat(appengine): Decrypt stored values (#2029)
0ef4c5c	chore: forward port release-1.4
3f0b864	fix(rm): Properly show default values for all mapping parameters (#2061)

$ gh api "repos/astarte-platform/astarte/commits/3f0b864" --jq '.files[] | "\(.filename)\t+\(.additions)/-\(.deletions)"' 2>/dev/null | head; echo === v1.2.2 file check ===; gh api "repos/astarte-platform/astarte/contents/apps/astarte_realm_management/lib/astarte_realm_management_web/router.ex?ref=v1.2.2" --jq '.name, .size' 2>/dev/null
apps/astarte_realm_management/lib/astarte_realm_management/interfaces/queries.ex	+3/-10
=== v1.2.2 file check ===
router.ex
885

$ gh api "repos/astarte-platform/astarte/commits/3f0b864" --jq '.files[] | .patch' 2>/dev/null | head -60
@@ -500,7 +500,7 @@ defmodule Astarte.RealmManagement.Interfaces.Queries do
 
       with {:ok, endpoints} <-
              Repo.fetch_all(endpoints_query, prefix: keyspace, consistency: consistency) do
-        mappings = Enum.map(endpoints, &to_mapping(&1, interface))
+        mappings = Enum.map(endpoints, &to_mapping(&1))
 
         interface =
           interface
@@ -520,15 +520,8 @@ defmodule Astarte.RealmManagement.Interfaces.Queries do
     end
   end
 
-  defp to_mapping(endpoint, interface) do
-    %Mapping{}
-    |> Mapping.changeset(Map.from_struct(endpoint),
-      interface_name: interface.name,
-      interface_id: interface.interface_id,
-      interface_major: interface.major_version,
-      interface_type: interface.type
-    )
-    |> Ecto.Changeset.apply_changes()
+  defp to_mapping(endpoint) do
+    Mapping.from_db_result!(endpoint)
     |> Map.from_struct()
     |> Map.put(:type, endpoint.value_type)
   end

The v1.2.2 router had no `/config/*` at all — those routes are all v1.3.0+. Let me pin down the remaining two claims: the empty-introspection change and the rc.3 type-validation restore.
$ gh search commits --repo astarte-platform/astarte "empty introspection" --json sha,commit --jq '.[] | "\(.sha[0:7])\t\(.commit.message | split("\n")[0])"' 2>/dev/null | head -3; echo ===; gh api "repos/astarte-platform/astarte/commits?since=2026-07-20T00:00:00Z&per_page=60" --jq '.[] | "\(.sha[0:7])\t\(.commit.message | split("\n")[0])"' 2>/dev/null | grep -iE "type valid|invalid|data" | head
9e1a4ed	Merge branch 'handle-empty-introspection' into 'master'
7c39e77	Merge branch 'do-not-crash-with-empty-introspection' into 'master'
3a838d3	data_updater_test: add empty introspection test
===

$ gh api "repos/astarte-platform/astarte/commits?since=2026-07-15T00:00:00Z&per_page=100&path=apps/astarte_data_updater_plant" --jq '.[] | "\(.sha[0:7])\t\(.commit.message | split("\n")[0])"' 2>/dev/null | head -30
5047681	feat(dup): Encrypt device shared secret (#2104)
0da8336	chore(dup): store shared secret during key agreement (#2101)
1c297aa	refactor: purge exrabbitpool (#2050)
c8bbee8	chore(dup): add rpc server startup log (#2096)
ed7cd4f	chore(dup): configure log level through env var at runtime (#2094)
6d0c0d6	chore(deps): bump httpoison and hackney
4424b08	chore(dup): update cose elixir dep
cf3cab4	chore(dup): add send_exchange_failed in missing failure branches
36a70c2	chore: forward port release-1.4

$ gh api "repos/astarte-platform/astarte/commits/9e1a4ed" --jq '.files[] | "\(.filename)\t+\(.additions)/-\(.deletions)"' 2>/dev/null | head; echo ===; gh api "repos/astarte-platform/astarte/commits/9e1a4ed" --jq '.commit.message' 2>/dev/null | head -5
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/impl.ex	+26/-2
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater_test.exs	+58/-5
apps/astarte_data_updater_plant/test/support/database_test_helper.exs	+36/-3
===
Merge branch 'handle-empty-introspection' into 'master'

Handle empty introspection

See merge request Astarte-NG/astarte_data_updater_plant!20

$ gh api "repos/astarte-platform/astarte/commits?since=2026-07-14T10:00:00Z&per_page=150" --jq '.[] | "\(.sha[0:7])\t\(.commit.message | split("\n")[0])"' 2>/dev/null | grep -viE "chore|deps|test:|ci|forward|merge|build|docs" | head -30
5047681	feat(dup): Encrypt device shared secret (#2104)
6e58d66	test(pairing): fix flakiness (#2100)
1c297aa	refactor: purge exrabbitpool (#2050)
fe5ad3c	fix(fdo): make use of vault_base_namespace in find_key
e3d47b7	feat(secrets): add configurable base namespace through `ASTARTE_VAULT_BASE_NAMESPACE` environment variable (default: `""`)
897a383	feat: Store shared secret as session_key

$ gh api "repos/astarte-platform/astarte/compare/v1.4.0-rc.2...v1.4.0-rc.3" --jq '.files[] | select(.filename | contains("data_updater")) | "\(.filename)\t+\(.additions)/-\(.deletions)"' 2>/dev/null | head -30
apps/astarte_data_updater_plant/config/config.exs	+2/-0
apps/astarte_data_updater_plant/config/prod.exs	+2/-5
apps/astarte_data_updater_plant/config/runtime.exs	+24/-0
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/data_handler.ex	+16/-17
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/rpc/server.ex	+1/-3
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/rpc/supervisor.ex	+63/-16
apps/astarte_data_updater_plant/mix.exs	+4/-4
apps/astarte_data_updater_plant/mix.lock	+2/-2
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_test.exs	+113/-0
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/rpc/server_test.exs	+6/-0
apps/astarte_data_updater_plant/test/support/cases/device.ex	+28/-5

$ gh api "repos/astarte-platform/astarte/compare/v1.4.0-rc.2...v1.4.0-rc.3" --jq '.files[] | select(.filename | endswith("data_handler.ex")) | .patch' 2>/dev/null
@@ -500,12 +500,14 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.DataHandler do
   end
 
   # TODO: We need tests for this function
-  def validate_value_type(expected_type, %DateTime{} = value) do
+  @spec validate_value_type(Mapping.t() | %{String.t() => Mapping.t()}, any()) ::
+          :ok | {:error, atom()}
+  def validate_value_type(%Mapping{value_type: expected_type}, %DateTime{} = value) do
     ValueType.validate_value(expected_type, value)
   end
 
   # From Cyanide 2.0, binaries are decoded as %Cyanide.Binary{}
-  def validate_value_type(expected_type, %Cyanide.Binary{} = value) do
+  def validate_value_type(%Mapping{value_type: expected_type}, %Cyanide.Binary{} = value) do
     %Cyanide.Binary{subtype: _subtype, data: bin} = value
     validate_value_type(expected_type, bin)
   end
@@ -515,10 +517,21 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.DataHandler do
     {:error, :unexpected_value_type}
   end
 
+  # TODO: we should test for this kind of unexpected messages
+  # We expected an individual value, but we received an aggregated
+  def validate_value_type(%Mapping{} = _expected_types, %{} = _object) do
+    {:error, :unexpected_value_type}
+  end
+
   def validate_value_type(%Mapping{value_type: expected_type}, value) do
-    validate_value_type(expected_type, value)
+    if value != nil do
+      ValueType.validate_value(expected_type, value)
+    else
+      :ok
+    end
   end
 
+  # object interface: validate all the types of the contained mappings
   def validate_value_type(%{} = mappings_by_key, %{} = object) do
     Enum.reduce_while(object, :ok, fn {key, value}, _acc ->
       with {:ok, %Mapping{value_type: expected_type}} <- Map.fetch(mappings_by_key, key),
@@ -538,26 +551,12 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.DataHandler do
     end)
   end
 
-  # TODO: we should test for this kind of unexpected messages
-  # We expected an individual value, but we received an aggregated
-  def validate_value_type(_expected_types, %{} = _object) do
-    {:error, :unexpected_value_type}
-  end
-
   # TODO: we should test for this kind of unexpected messages
   # We expected an aggregated, but we received an individual
   def validate_value_type(%{} = _expected_types, _object) do
     {:error, :unexpected_value_type}
   end
 
-  def validate_value_type(expected_type, value) do
-    if value != nil do
-      ValueType.validate_value(expected_type, value)
-    else
-      :ok
-    end
-  end
-
   def update_stats(state, interface, major, path, payload) do
     exchanged_bytes = byte_size(payload) + byte_size(interface) + byte_size(path)

$ ls .mule/tasks/ 2>/dev/null | head; echo "---"; ls .mule/plans/ 2>/dev/null | head
issue-16.md
---
2026-07-27-device-deletion-events.md

I now have all the evidence I need. Let me verify the `.mule/tasks/` convention and Astrate's housekeeping body handling before writing the queue lines.
→ Read internal/housekeeping/http.go [limit=70, offset=30]

← Write .mule/tasks/realm-config-datastream-retention.md
Wrote file successfully.

All evidence gathered. Now appending the proposal lines to `.mule/todo.md` and a detail file for the one complex task.
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -85,4 +85,9 @@
 - [!] race-check: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...` (~40s). Report any failure to .mule/for-giulio.md with the full race report. This is the only race coverage that exists — the Pi cannot run -race. [legion] [readonly] — BLOCKED: TIMEOUT after 1200s — task too big, split it
 - [x] store-realm-cascade-policies: in `internal/store/realms_test.go` `CascadeDelete`, add `trigger_policies` to the post-delete verification loop (query `SELECT count(*) FROM trigger_policies WHERE realm_id = $1`). The migration 000006 has ON DELETE CASCADE but nothing asserts on it. [auto]
 - [x] store-alias-lowest-id: in `internal/store/devices_test.go`, add a subtest that registers two devices in the same realm, sets the same alias tag on both, and asserts `GetDeviceByAlias` returns the one with the lower UUID. The SQL uses `ORDER BY id LIMIT 1` but no test proves it. [auto]
 - [x] store-delete-device-objects: in `internal/store/devices_test.go` `StatsAndDelete`, insert object datastream rows for the device before deleting it, and assert they are gone after the delete. `DeleteDevice` explicitly sweeps `object_datastreams` but the test only checks individual rows. [auto]
+- [ ] control-producer-properties-compression: in `internal/engine/control.go`, accept a plaintext (and the 4-zero-byte empty-frame) device→server `producer/properties` purge list when the device declared `purge_properties_compression_format: plaintext` — `inflateProperties` currently rejects every non-zlib payload while upstream's `control_handler.ex` decodes per-capability [auto]
+- [ ] realm-config-datastream-retention: add `GET /realmmanagement/v1/<realm>/config/datastream_maximum_storage_retention` to `internal/realm/http.go` and carry `datastream_maximum_storage_retention` through `internal/housekeeping/http.go` + `internal/store/realms.go` (upstream v1.3.0 per-realm retention, v1.4.0-rc.1 adds an instance-wide default); detail in `.mule/tasks/realm-config-datastream-retention.md` [auto]
+- [ ] probe-interface-default-values: does `GET /realmmanagement/v1/<realm>/interfaces/<name>/<major>` return the same defaulted mapping-parameter values upstream reports after 3f0b864 ("Properly show default values for all mapping parameters")? report, do not patch [auto]
+- [ ] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto]
+- [ ] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto]

→ Read .mule/todo.md [offset=86]
Done: proposed 5 new queue lines from the v1.2.x → v1.4.0-rc.3 upstream gap (4 tasks/probes + compat-note line); no code changed, no git touched.

**Findings** — upstream is ahead (`v1.4.0-rc.3`, 2026-07-31). After diff-checking the release notes for the items that matter:

Proposed (new wire-visible gaps):
1. `control-producer-properties-compression` — upstream's `control_handler.ex` decodes the device→server purge list per the `purge_properties_compression_format` capability (plaintext + all-zero empty frame); Astrate's `inflateProperties` (`control.go:290`) rejects every non-zlib payload.
2. `realm-config-datastream-retention` — per-realm datastream retention: RM `GET /config/datastream_maximum_storage_retention` (no v1.2.2 baseline; v1.3.0+) + housekeeping body field; Astrate has only global retention. Detail file written.
3. `probe-interface-default-values` — upstream fix `3f0b864` ("show default values for all mapping parameters") in the interface GET response; probe only.
4. `probe-value-type-validation` — upstream rc.3 restored `validate_value_type` (aggregated-on-individual now a handled error); probe only.
5. `compat-note-v1.4.0-rc.3` (per recipe, always last).

Checked and **not** proposed (already covered): empty introspection (Astrate accepts), `device_registered` trigger (Astrate emits it), deletion-trigger pair (issue #21 closed/mule-review — already in flight, so not re-proposed), capabilities topic + outbound plaintext consumer/properties (Astrate implements). Dropped as non-applicable: FDO, AMQP vhosts/queues, replication strategies.

Files: `.mule/todo.md`, `.mule/tasks/realm-config-datastream-retention.md`
Verified: `gh api .../releases`, `gh api` tree/contents/commits/compare fetches, `rg` over the repo -> pass
Unsure: whether real SDKs ever send plaintext producer lists (upstream accepts it regardless, so it's still an obligation); the retention task may exceed one tick — suggested a split in the detail file; compat-note tag is a pre-release, wording is a human call.
Follow-ups: the 5 `- [ ]` lines appended to `.mule/todo.md` (lines 89-93)
