slug: recipe-astarte-upstream
verdict: proposed
at:  28f3adf
ran: 2026-08-24T02:14:45Z on DietPi in 396s


> build · x-preview-f-free

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md
Recipe understood: check upstream Astarte's newest release tag against `docs/COMPATIBILITY.md`; if they match, report "no movement" and stop — the deliverable is either nothing or proposed task lines in `.mule/todo.md`.

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:20][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z
v1.4.0-rc.1	2026-07-14T08:53:27Z
v1.3.2	2026-07-14T10:07:26Z
v1.3.1	2026-07-13T11:50:05Z
v1.3.0	2026-05-06T15:30:50Z
v1.2.2	2026-04-27T16:39:28Z
v1.4.0-rc.0	2026-04-08T12:32:33Z
v1.3.0-rc.2	2026-04-08T16:34:59Z
v1.2.2-rc.0	2026-04-08T16:31:03Z
v1.2.1	2026-03-12T16:11:20Z
v1.2.1-rc.1	2026-02-13T10:36:20Z
v1.3.0-rc.1	2026-01-26T15:58:33Z
v1.3.0-rc.0	2025-11-21T14:58:46Z
v1.2.1-rc.0	2025-08-27T07:03:28Z
v1.2.1-alpha.0	2025-04-10T15:31:04Z
v1.2.0	2024-07-02T14:09:07Z

Upstream is ahead: `docs/COMPATIBILITY.md` targets upstream **1.2.x** (`APICompatVersion` = `1.2.2`, measurements recorded against `v1.2.0`, Dashboard `v1.2.2`), while upstream's newest tags are `v1.4.0-rc.5` (2026-08-20) and stable `v1.3.3` (2026-08-07). Per the recipe, I read the release notes before proposing anything.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
Sixth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.4

### Fixed

- [astarte_data_updater_plant] Make object values only validate type on selected interface
- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated

`v1.3.3`'s body is empty; `v1.4.0-rc.5` only lists two data-updater validation fixes since rc.4. Astrate targets `v1.2.2` (2026-04-27), so the real gap is the whole **v1.3 line** (v1.3.0, May 2026) plus the v1.4.0 RC series. Fetching those notes.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.0") | .body'
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

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.1") | .body'
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

The gap is exactly **v1.3.0's "Changes since 1.2.2"** plus the v1.4.0 RCs. Filtering strictly per the recipe (wire-visible / schema rules qualify; AMQP vhosts, clustering env vars, Cassandra replication, container merges, FDO plumbing are infra/Elixir-shaped and ignored), the qualifying items are: empty-introspection acceptance, the new MQTT v1 **capabilities** mechanism (`purge_properties_compression_format`), the v1.4.0-rc.5 data-validation fixes, and the v1.3.0 properties-on-connect encoding fix. Before writing task lines, three narrow greps so each line names a real file.
$ rg -n -i "purge|producer/properties" --type go -l | head -8
internal/engine/engine.go
internal/engine/control.go
internal/engine/data_test.go
test/security/security_test.go
internal/engine/serverdata.go
internal/engine/fakestore_test.go
internal/engine/control_test.go
internal/testutil/astartedevice.go

$ rg -n -i "object" internal/engine/*.go --type go | grep -i -E "aggregat|valid" | head -6; rg -n "binaryblob|BinaryBlob|binary_blob" --type go -l | head -5
internal/engine/serverdata_test.go:267:// TestServerObjectAggregate: object-aggregated server interfaces accept one
internal/engine/serverdata_test.go:269:func TestServerObjectAggregate(t *testing.T) {
internal/engine/serverdata_test.go:274:		"type": "datastream", "ownership": "server", "aggregation": "object",
internal/engine/engine_e2e_test.go:47:		"type": "datastream", "ownership": "device", "aggregation": "object",
internal/engine/engine_e2e_test.go:262:// BSON object-aggregated publishes all land in the right columns. QoS >= 1
internal/engine/engine_integration_test.go:264:// testObjectRows: object-aggregated publishes land as one jsonb document on
pkg/interfaceschema/compile_test.go
pkg/interfaceschema/parse_test.go
pkg/interfaceschema/types.go
pkg/payload/value.go
pkg/payload/json.go

$ rg -n "introspection" --type go -l internal/ pkg/ | head -8; rg -rn -i "empty introspection|len\(.*introspection.*\) == 0" --type go | head -4
internal/store/store.go
internal/engine/data.go
internal/engine/engine.go
internal/engine/control.go
internal/store/interfaces.go
internal/engine/data_test.go
internal/store/devices.go
internal/store/devices_test.go
internal/engine/introspection.go:// (docs/DESIGN.md §3.3). The empty payload is a valid n
internal/engine/introspection_test.go:			name: "empty payload is an n",

→ Read internal/engine/introspection.go
Astrate already accepts empty-introspection payloads by design (`internal/engine/introspection.go:96-103`, tested), so that v1.3.0 change needs no proposal — dropped. Checking `.mule/todo.md` before appending so slugs don't collide.
→ Read .mule/todo.md
Key finding: **a previous run of this recipe already queued most of this gap** — lines 93–100 of `.mule/todo.md` already cover the `purge_properties_compression_format` capability, per-realm retention (incl. the rc.1 instance default), the mapping-defaults fix, object value-type validation, the v1.3.3 emptyCache device-error fix, both deletion triggers, and both compat-note lines. Per the no-duplicates rule I will not re-propose any of those.

Remaining uncovered items from the notes: binaryblob validation (rc.5), property encoding on connect (v1.3.0 Fixed), realm-scoped pairing health checks (v1.3.0 Added), device **registration** triggers (deletion ones are queued, registration is not), and the capability-declaration wire syntax itself. One cheap grep so the trigger line names real symbols:

$ rg -n "On[A-Z][a-zA-Z]* =" internal/engine/triggers/match.go | head -20
27:	OnIncomingData = "incoming_data"
30:	OnValueChange = "value_change"
33:	OnValueChangeApplied = "value_change_applied"
36:	OnPathCreated = "path_created"
39:	OnPathRemoved = "path_removed"
42:	OnValueStored = "value_stored"
44:	OnDeviceRegistered = "device_registered"
46:	OnDeviceConnected = "device_connected"
48:	OnDeviceDisconnected = "device_disconnected"
51:	OnDeviceEmptyCacheReceived = "device_empty_cache_received"
53:	OnDeviceError = "device_error"
55:	OnIncomingIntrospection = "incoming_introspection"
58:	OnInterfaceAdded = "interface_added"
60:	OnInterfaceRemoved = "interface_removed"
63:	OnInterfaceMinorUpdated = "interface_minor_updated"

`OnDeviceRegistered` (`device_registered`) already exists at `internal/engine/triggers/match.go:44`, so upstream's new registration triggers are already matched by Astrate — dropped. Checking the remaining un-read RC notes (rc.0, rc.3, rc.4) so I don't miss anything before appending:
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.3" or .tag_name=="v1.4.0-rc.4") | .body'
Fifth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.3

### Added

- [config] Initialize `config` library. All external services declared using the `config` library helpers can be configured either using the `_URL` environment variable or by using single component variables `_SCHEME`, `_HOST`, `_PORT`, `_PATH`, `_QUERY`, `_FRAGMENT` (`_URL` has priority), together with normal `_SSL_*` variables.
- [secrets] Use `config` for vault configuration. Vault can now be configured using component variables instead of `ASTARTE_VAULT_URL`:
  - `ASTARTE_VAULT_SCHEME`: scheme part of the URL (binary, default: `"http"`, or `"https"` if `ASTARTE_VAULT_SSL_ENABLED` is `true`)
  - `ASTARTE_VAULT_HOST`: host part of the URL (binary, default: `"localhost"`)
  - `ASTARTE_VAULT_PORT`: port part of the url (integer, default `8200`)
  - `ASTARTE_VAULT_PATH`: path part of the url (string)
  - `ASTARTE_VAULT_QUERY`: query part of the url (string)
  - `ASTARTE_VAULT_FRAGMENT`: fragment part of the url (string)
- [fdo] Use `config` for rendezvous configuration. The rendezvous can now be configured using component variables instead of `PAIRING_FDO_RENDEZVOUS_URL`:
  - `PAIRING_FDO_RENDEZVOUS_SCHEME`: scheme part of the URL (binary, default: `"http"`, or `"https"` if `PAIRING_FDO_RENDEZVOUS_SSL_ENABLED` is `true`)
  - `PAIRING_FDO_RENDEZVOUS_HOST`: host part of the URL (binary, default: `"localhost"`)
  - `PAIRING_FDO_RENDEZVOUS_PORT`: port part of the url (integer, default `8041`)
  - `PAIRING_FDO_RENDEZVOUS_PATH`: path part of the url (string)
  - `PAIRING_FDO_RENDEZVOUS_QUERY`: query part of the url (string)
  - `PAIRING_FDO_RENDEZVOUS_FRAGMENT`: fragment part of the url (string)

### Fixed

- [astarte_data_updater_plant] Ensure RPC server is always available to clients. Resolved the issue where a temporary disconnection and reconnection of data_updater_plant to the cluster would make the RPC server inaccessible.

### Changed

- [astarte_data_updater_plant] Ensure memory is properly garbage collected

Fourth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.2

### Added

- [secrets] Configurable base namespace for Vault using `ASTARTE_VAULT_BASE_NAMESPACE`: all created namespaces will be children of the base namespace.

### Fixed

- [fdo] Support for Vault and Vault managed
- [fdo] Properly handle chunked service info messages
- [astarte_data_updater_plant] Avoid clogging AMQP queues when receiving invalid device messages
- [astarte_data_updater_plant] Restore correct type validation checks for incoming data

rc.3/rc.4 are all infra-shaped (config library, Vault, FDO, RPC/GC) — nothing new qualifies. Appending the four uncovered items to `.mule/todo.md`:
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -97,4 +97,8 @@
 - [ ] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto]
 - [ ] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto]
 - [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]
 - [ ] trigger-device-deletion-events: add `device_deletion_started` and `device_deletion_finished` trigger conditions to `internal/engine/triggers/match.go` (new constants + entries in `deviceOns`), emit them from the device-deletion path in `internal/engine/engine.go`, and add them to the `on` enum in `docs/api/astarte_realm_management_api.yaml` (upstream v1.3.0-rc.0) [auto]
+- [ ] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto]
+- [ ] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto]
+- [ ] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto]
+- [ ] probe-pairing-realm-health: record upstream v1.3.0's realm-scoped pairing health-check routes and response shapes (added for devices aware only of realm-scoped paths, e.g. FDO) and check whether internal/pairing serves any equivalent; report, do not patch [auto]

Upstream is ahead of the targeted 1.2.x (newest stable v1.3.3, newest tag v1.4.0-rc.5), but an earlier run already queued most qualifying changes (todo.md lines 93–100: purge compression capability + retention config, mapping-defaults and object value-type validation fixes, emptyCache device-error fix, deletion triggers, both compat-notes — so no compat-note line was re-added). Items dropped as already satisfied or non-qualifying: empty-introspection acceptance (Astrate already accepts it, `internal/engine/introspection.go:96`), registration triggers (`device_registered` exists at `internal/engine/triggers/match.go:44`), AMQP/clustering/Cassandra/FDO/Vault/config-library items (infra or Elixir-shaped).

```
Done: compared upstream releases to docs/COMPATIBILITY.md (targets 1.2.x; upstream now v1.3.3 stable + v1.4.0-rc.5) and proposed the 4 not-yet-queued compatibility items
Files: .mule/todo.md
Verified: gh api release list + release notes for v1.3.0/v1.3.3/v1.4.0-rc.* read -> n/a (proposal job); rg/read confirmed device_registered trigger and empty-introspection support already exist
Unsure: rc-series notes may compress older fixes oddly (rc.5 lists two validation fixes "since rc.4" that overlap queued probe-value-type-validation's rc.3 reference); capability-declaration wire syntax unverified by design — that is what the probe asks
Follow-ups: probe-mqtt-capabilities-declaration, probe-binaryblob-validation, probe-properties-on-connect-encoding, probe-pairing-realm-health
```
