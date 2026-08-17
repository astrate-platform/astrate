slug: recipe-astarte-upstream
verdict: proposed
at:  6e82c66
ran: 2026-08-17T02:08:56Z on DietPi in 333s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z
v1.4.0-rc.1	2026-07-14T08:53:27Z

Upstream is at v1.3.3, Astrate targets v1.2.x. Let me read the release notes.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name | startswith("v1.3")) | "\(.tag_name)\t\(.name)\t\(.body[0:500])"'
v1.3.3		
v1.3.2	v1.3.2	Astarte v1.3.2 release.

## Changes since v1.3.1

## Fixed

- Avoid crash on invalid properties message

v1.3.1	v1.3.1	Astarte v1.3.1 release.

## Changes since v1.3.0

### Fixed

- Solve capabilities message stalling queues
- [astarte_housekeeping] Allow resuming partially completed realm deletions
- Use network topology as default realm replication
- Fix simple strategy realm creation without replication factor

### Changed

- Immediately delete devices in device deletion best case scenario

v1.3.0	v1.3.0	Astarte v1.3.0 release.

## Changes since 1.2.2

### Added

- [astarte_realm_management] realm migrations now create vhost for amqp triggers
- New environment variables for trigger notifications between realm management replicas and realm management -> pairing. These variables are currently being used only by realm management
  - `REALM_MANAGEMENT_CLUSTERING_KUBERNETES_SELECTOR`. The Endpoint label to query to get realm management instances. Defaults to `app=astarte-realm-management`.
 
v1.3.0-rc.2	v1.3.0-rc.2	Third Astarte v1.3.0 release candidate.

## Changes since v1.3.0-rc.1

### Added

- [astarte_realm_management] realm migrations now create vhost for amqp triggers

### Changed

- [astarte_pairing] FDO authentication (EXPERIMENTAL feature, disabled by default). New environment variables are needed in order to use FDO:
  - `PAIRING_ENABLE_FDO` - whether the FDO feature is enabled or not (default: false)
  - `PAIRING_FDO_RENDEZVOUS_URL` - URL of the rendezvous server (default: "http://r
v1.3.0-rc.1	v1.3.0-rc.1	Second Astarte v1.3.0 release candidate.

# Changes since v1.3.0-rc.0

### Added

- New environment variables for trigger notifications between realm management replicas and realm management -> pairing. These variables are currently being used only by realm management
  - `REALM_MANAGEMENT_CLUSTERING_KUBERNETES_SELECTOR`. The Endpoint label to query to get realm management instances. Defaults to `app=astarte-realm-management`.
  - `PAIRING_CLUSTERING_KUBERNETES_SELECTOR`. The Endpoint la
v1.3.0-rc.0	v1.3.0-rc.0	First Astarte v1.3.0 release candidate.

## Changes since 1.2
### Added
- [astarte_data_updater_plant] Added separate AMQP triggers producer configuration (falls back to general AMQP producer settings if unset):
  - `DATA_UPDATER_PLANT_AMQP_TRIGGERS_PRODUCER_HOST` - Host for triggers producer connection (no default, falls back to producer host)
  - `DATA_UPDATER_PLANT_AMQP_TRIGGERS_PRODUCER_USERNAME` - Username for triggers producer (no default, falls back to producer username)
  - `DATA_

$ gh search code --owner astarte-platform "capabilities" --language Elixir -L 10 --json path,repository,textMatches
[]

$ gh api repos/astarte-platform/astarte/contents/CHANGELOG.md --jq '.content' | base64 -d | head -200
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

## [1.4.0-rc.4] - 2026-08-06

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

## [1.4.0-rc.3] - 2026-07-31

### Added

- [secrets] Configurable base namespace for Vault using `ASTARTE_VAULT_BASE_NAMESPACE`: all created namespaces will be children of the base namespace.

### Fixed

- [fdo] Support for Vault and Vault managed
- [fdo] Properly handle chunked service info messages
- [astarte_data_updater_plant] Avoid clogging AMQP queues when receiving invalid device messages
- [astarte_data_updater_plant] Restore correct type validation checks for incoming data

## [1.4.0-rc.2] - 2026-07-14

## [1.4.0-rc.1] - 2026-07-13

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

## [1.4.0-rc.0] - 2026-04-08

### Added

- Add required flag for mappings of object aggregated interfaces
- [astarte_pairing] Support for FIDO Device Onboard v1.1 for device registration. The following environment variables have been added for the Vault connectivity:
  - `ASTARTE_VAULT_URL` - The URL to access Vault.
  - `ASTARTE_VAULT_AUTHENTICATION_MECHANISM` - The mechanism to use for authenticating with Vault (default: `""`, allowed: \["token"])
  - `ASTARTE_VAULT_TOKEN` - Token to authenticate with Vault, used if `ASTARTE_VAULT_AUTHENTICATION_MECHANISM` is `"token"` (default: `""`)
  - `ASTARTE_VAULT_SSL_CA_FILE` - path to the CA certificate file for the vault TLS connection; when not specified, the bundled cURL certificate bundle will be used
  - `ASTARTE_VAULT_SSL_CUSTOM_SNI` - custom SNI hostname for the vault TLS connection; when not specified, the vault hostname will be used
  - `ASTARTE_VAULT_SSL_DISABLE_SNI` - Disable Server Name Indication for Vault. (default: `false`)
  - `ASTARTE_VAULT_SSL_ENABLED` - Whether to enable SSL for the Vault connection. (default `false`)
- [astarte_realm_management] Allow listing interfaces with mappings using the `detailed=true` parameter
- [astarte_import] Added support for data types: `doublearray`, `integerarray`,
  `booleanarray`, `longintegerarray`, `stringarray`, `datetimearray`, `binaryblobarray`.

### Changed

- [astarte_data_updater_plant] Use mississippi consumer for data updater processes

## [1.3.3] - Unreleased

### Fixed

- [astarte_data_updater_plant] Resending properties to a device now reports distinct device error triggers depending on the failure: `interface_loading_failed` when an interface fails to load and `resend_interface_properties_failed` when sending the properties to the device itself fails. Previously, both cases were incorrectly reported as `resend_interface_properties_failed`.

## [1.3.2] - 2026-07-14

## Fixed

- Avoid crash on invalid properties message

## [1.3.1] - 2026-07-13

### Fixed

- Solve capabilities message stalling queues
- [astarte_housekeeping] Allow resuming partially completed realm deletions
- Use network topology as default realm replication
- Fix simple strategy realm creation without replication factor

### Changed

- Immediately delete devices in device deletion best case scenario

## [1.3.0] - 2026-05-06

### Changed

- [astarte_housekeeping] When `HOUSEKEEPING_ASTARTE_KEYSPACE_REPLICATION_STRATEGY` is not set, the `astarte` keyspace is now created using `NetworkTopologyStrategy` with a replication map derived from the current ScyllaDB network topology (one replica per node in each datacenter), instead of falling back to `SimpleStrategy` with replication factor 1.
- [astarte_housekeeping] When creating a realm without specifying its replication strategy, the strategy used in the `astarte` keyspace is applied as default.

### Fixed

- [astarte_housekeeping] Changed logger config to make the metadata fields appear in logfmt log output; in particular the msg field
- [astarte_pairing] Changed logger config to make the metadata fields appear in logfmt log output; in particular the msg field

## [1.3.0-rc.2] - 2026-04-08

### Added

- [astarte_realm_management] realm migrations now create vhost for amqp triggers

### Changed

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

## [1.3.0-rc.1] - 2026-01-26

### Added

- New environment variables for trigger notifications between realm management replicas and realm management -> pairing. These variables are currently being used only by realm management
  - `REALM_MANAGEMENT_CLUSTERING_KUBERNETES_SELECTOR`. The Endpoint label to query to get realm management instances. Defaults to `app=astarte-realm-management`.
  - `PAIRING_CLUSTERING_KUBERNETES_SELECTOR`. The Endpoint label to query to get pairing instances. Defaults to `app=astarte-pairing`.
- [astarte_pairing] Cluster with realm management using `CLUSTERING_STRATEGY` and `CLUSTERING_KUBERNETES_NAMESPACE`
- [astarte_pairing] Add realm-scoped health checks, which can be used to test astarte health or connectivity by devices who are only aware of realm-scoped paths, such as during FDO

### Fixed

- Fix crashes in new ScyllaDB versions
- [astarte_appengine_api] Volatile trigger requests to Data Updater Plant now properly encode the volatile trigger type

### Changed

- Services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages

## [1.3.0-rc.0] - 2025-11-21

### Added

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

- BREAKING: Merged API services into main services, eliminating separate containers:
  - `astarte-housekeeping-api` merged into `astarte-housekeeping`
  - `astarte-pairing-api` merged into `astarte-pairing`
  - `astarte-realm-management-api` merged into `astarte-realm-management`
- BREAKING: [astarte_housekeeping] Housekeeping now creates an AMQP vhost for each created realm. Some required configuration was introduced:
  - `HOUSEKEEPING_AMQP_HOST` - AMQP host for housekeeping operations (required in production, defaults to "localhost" in dev/test)

• Research Astrate capabilities Explore Agent

$ gh api repos/astarte-platform/astarte/contents/CHANGELOG.md --jq '.content' | base64 -d | sed -n '/## \[1\.3\.0-rc\.0\]/,/## \[1\.2/p' | head -150
## [1.3.0-rc.0] - 2025-11-21

### Added

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

- [astarte_data_updater_plant] Do not crash when transient triggers are installed on devices with outdated introspection
- [astarte_data_updater_plant] Correctly encode values when sending properties to device on connection
- [astarte_realm_management] Allow to delete long-disconnected devices. For this to work, Realm Management needs `CLUSTERING_STRATEGY`, `CLUSTERING_KUBERNETES_NAMESPACE` and `DATA_UPDATER_PLANT_CLUSTERING_KUBERNETES_SELECTOR` to be set, just like AppEngine and DUP. Refer to 1.2.1-rc.0 for additional information on the variables.

## [1.2.2] - 2026-04-27

✓ Research Astrate capabilities Explore Agent
→ Read .mule/todo.md
• Check device triggers Astrate Explore Agent
✓ Check device triggers Astrate Explore Agent
Now I have the full picture. Let me check what's already tracked vs. what's new.

$ grep -n 'device.deletion.*trigger\|deletion.*trigger\|probe-emptycache\|compat-note-v1.3.3\|compat-note-v1.3' /root/astrate-mule/.mule/todo.md
94:- [ ] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto]
95:- [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]

Upstream is at v1.3.3, Astrate targets v1.2.x. Most v1.3.x changes are already tracked. The one wire-visible gap not yet in the queue is device deletion trigger events. Let me append the new task lines.
→ Read .mule/todo.md [offset=90]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -92,4 +92,5 @@
 - [ ] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto]
 - [ ] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto]
 - [ ] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto]
 - [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]
+- [ ] trigger-device-deletion-events: add `device_deletion_started` and `device_deletion_finished` trigger conditions to `internal/engine/triggers/match.go` (new constants + entries in `deviceOns`), emit them from the device-deletion path in `internal/engine/engine.go`, and add them to the `on` enum in `docs/api/astarte_realm_management_api.yaml` (upstream v1.3.0-rc.0) [auto]

Here's the summary:

**Upstream: v1.3.3 (released 2026-08-07) vs. Astrate target: v1.2.x**

Most v1.3.x changes are already tracked in the queue:
- `control-producer-properties-compression` — plaintext producer/properties (v1.3.0 capability)
- `realm-config-datastream-retention` — per-realm data retention (v1.3.0)
- `probe-emptycache-resend-device-error` — v1.3.3 error differentiation
- `compat-note-v1.3.3` / `compat-note-v1.4.0-rc.3` — compat notes

**New task added:**
- `trigger-device-deletion-events` — upstream v1.3.0-rc.0 added `device_deletion_started` and `device_deletion_finished` trigger conditions. Astrate has no equivalent. This is wire-visible (Realm Management trigger API `on` enum, device lifecycle events).

**Already implemented (no task needed):**
- Empty introspection accepted — already handled in `introspection.go:95-102`
- Capabilities + `purge_properties_compression_format` — fully implemented in `capabilities.go`
- Device registration triggers — already wired in `engine.go:283-301`
- Transient trigger crash fix — N/A (Astrate's transient triggers are matcher-only, never enter pipeline)

**Infrastructure-only / Elixir-only (ignored per recipe rules):**
- AMQP vhost per realm, merged API containers, env var changes, ScyllaDB replication strategy, capabilities queue stalling fix, volatile trigger encoding fix — none of these are wire-visible or applicable to Astrate's architecture.

Done: appended 1 task line (trigger-device-deletion-events) to .mule/todo.md
Files: .mule/todo.md
Verified: grep confirms no duplicate
Unsure: nothing
Follow-ups: trigger-device-deletion-events
