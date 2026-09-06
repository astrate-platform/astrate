slug: recipe-astarte-upstream
verdict: proposed
at:  eaa2c37
ran: 2026-09-06T18:31:02Z on DietPi in 526s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

✱ Grep "APICompatVersion|target.*version|upstream.*version" in . · 11 matches

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
Sixth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.4

### Fixed

- [astarte_data_updater_plant] Make object values only validate type on selected interface
- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated

The latest stable release is **v1.3.3** (2026-08-07), ahead of `APICompatVersion = "1.2.2"`. Let me fetch the release notes.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name | startswith("v1.3")) | "\(.tag_name)\t\(.published_at)\n\(.body)\n---"'
v1.3.3	2026-08-07T14:33:58Z

---
v1.3.2	2026-07-14T10:07:26Z
Astarte v1.3.2 release.

## Changes since v1.3.1

## Fixed

- Avoid crash on invalid properties message

---
v1.3.1	2026-07-13T11:50:05Z
Astarte v1.3.1 release.

## Changes since v1.3.0

### Fixed

- Solve capabilities message stalling queues
- [astarte_housekeeping] Allow resuming partially completed realm deletions
- Use network topology as default realm replication
- Fix simple strategy realm creation without replication factor

### Changed

- Immediately delete devices in device deletion best case scenario

---
v1.3.0	2026-05-06T15:30:50Z
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


---
v1.3.0-rc.2	2026-04-08T16:34:59Z
Third Astarte v1.3.0 release candidate.

## Changes since v1.3.0-rc.1

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

---
v1.3.0-rc.1	2026-01-26T15:58:33Z
Second Astarte v1.3.0 release candidate.

# Changes since v1.3.0-rc.0

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

---
v1.3.0-rc.0	2025-11-21T14:58:46Z
First Astarte v1.3.0 release candidate.

## Changes since 1.2
### Added
- [astarte_data_updater_plant] Added separate AMQP triggers producer configuration (falls back to general AMQP producer settings if unset):
  - `DATA_UPDATER_PLANT_AMQP_TRIGGERS_PRODUCER_HOST` - Host for triggers producer connection (no default, falls back to producer host)
  - `DATA_UPDATER_PLANT_AMQP_TRIGGERS_PRODUCER_USERNAME` - Username for triggers producer (no default, falls back to producer username)
  - `DATA_UPDATER_PLANT_AMQP_TRIGGERS_PRODUCER_PASSWORD` - Password for triggers producer (no default, falls back to producer password)
  - `DATA_UPDATER_PLANT_AMQP_TRIGGERS_PRODUCER_VIRTUAL_HOST` - Virtual host for triggers producer (no default, falls back to producer virtual host)
  - `DATA_UPDATER_PLANT_AMQP_TRIGGERS_PRODUCER_PORT` - Port for triggers producer (no default, falls back to producer port)
  - `DATA_UPDATER_PLANT_AMQP_TRIGGERS_PRODUCER_SSL_ENABLED` - Enable SSL for triggers producer (no default, falls back to producer SSL setting)
  - `DATA_UPDATER_PLANT_AMQP_TRIGGERS_PRODUCER_SSL_CA_FILE` - CA certificate file for triggers producer SSL (falls back to producer CA file or bundled cURL certificates)
  - `DATA_UPDATER_PLANT_AMQP_TRIGGERS_PRODUCER_SSL_DISABLE_SNI` - Disable Server Name Indication for triggers producer (default: false)
  - `DATA_UPDATER_PLANT_AMQP_TRIGGERS_PRODUCER_SSL_CUSTOM_SNI` - Custom SNI hostname for triggers producer (falls back to producer host if unset)
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
- BREAKING: [astarte_pairing] AMQP Producer configuration is now mandatory using the `ASTARTE_EVENTS_PRODUCER_AMQP_*` environment variables
- BREAKING: [astarte_realm_management] AMQP Producer configuration is now mandatory using the `ASTARTE_EVENTS_PRODUCER_AMQP_*` environment variables
- More accurate health checks for astarte services

### Fixed
- [astarte_data_updater_plant] Do not crash when transient triggers are installed on devices with outdated introspection
- [astarte_data_updater_plant] Correctly encode values when sending properties to device on connection
- [astarte_realm_management] Allow to delete long-disconnected devices
---

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | "\(.tag_name)\t\(.published_at)\t\(.body | length)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z	256
v1.3.3	2026-08-07T14:33:58Z	0
v1.4.0-rc.4	2026-08-06T13:40:20Z	2054
v1.4.0-rc.3	2026-07-31T10:48:39Z	543
v1.4.0-rc.2	2026-07-14T10:16:22Z	39
v1.4.0-rc.1	2026-07-14T08:53:27Z	629
v1.3.2	2026-07-14T10:07:26Z	111
v1.3.1	2026-07-13T11:50:05Z	393
v1.3.0	2026-05-06T15:30:50Z	9079
v1.2.2	2026-04-27T16:39:28Z	395
v1.4.0-rc.0	2026-04-08T12:32:33Z	1523
v1.3.0-rc.2	2026-04-08T16:34:59Z	2282
v1.2.2-rc.0	2026-04-08T16:31:03Z	398
v1.2.1	2026-03-12T16:11:20Z	5043
v1.2.1-rc.1	2026-02-13T10:36:20Z	435
v1.3.0-rc.1	2026-01-26T15:58:33Z	1234
v1.3.0-rc.0	2025-11-21T14:58:46Z	4584
v1.2.1-rc.0	2025-08-27T07:03:28Z	3336
v1.2.1-alpha.0	2025-04-10T15:31:04Z	1243
v1.2.0	2024-07-02T14:09:07Z	2795
v1.2.0-rc.0	2024-06-11T13:02:33Z	2619
v1.0.6	2024-04-23T21:15:57Z	926
v1.1.1	2023-11-15T15:29:57Z	1387
v1.0.5	2023-09-26T10:27:57Z	1555
v1.1.0	2023-06-20T18:59:09Z	1459
v1.1.0-rc.0	2023-06-09T15:49:44Z	505
v1.1.0-alpha.0	2022-11-24T16:38:37Z	921
v1.0.4	2022-10-25T16:50:21Z	319
v1.0.3	2022-07-05T13:57:17Z	1622
v1.0.2	2022-04-01T17:18:15Z	967

Now let me investigate the wire-visible changes. The key candidate from v1.3.0 is the new **capabilities** feature. Let me check what Astrate already handles.
✱ Grep "capabilities|purge_properties_compression" in . · 54 matches

$ gh search code --owner astarte-platform 'capabilities' --match path --filename '*.ex' -L 10 2>&1 | head -30
(no output)

The capabilities feature is already implemented. The leading candidates from v1.3.0 are the new **device registration** and **device deletion** trigger types. Let me check what trigger event types Astrate supports, and what actually changed in v1.3.3 (empty release body).
✱ Grep "device_connected|device_disconnected|device_registration|device_deleted|device_created" in . · 89 matches

$ gh api 'repos/astarte-platform/astarte/compare/v1.3.2...v1.3.3' --jq '.commits[] | "\(.sha[0:8]) \(.commit.message | split("\n")[0])"'
d8592091 fix(dup): stop remapping interface loading failed error (#2119)
3e70d63f chore: prepare v1.3.3 (#2128)

v1.3.3 contains one fix: "stop remapping interface loading failed error" — directly relevant to Astrate's `UpstreamErrorName` mapping (COMPATIBILITY.md deviation 1). Let me look at the actual diff, and check the new trigger types upstream added.
$ gh api 'repos/astarte-platform/astarte/commits/d8592091' --jq '.files[] | "FILE: \(.filename) (+\(.additions) -\(.deletions))\n\(.patch // "no patch")\n"' 2>&1 | head -120
FILE: CHANGELOG.md (+6 -0)
@@ -5,6 +5,12 @@ All notable changes to this project will be documented in this file.
 The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
 and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).
 
+## Unreleased
+
+### Fixed
+
+- [astarte_data_updater_plant] Resending properties to a device now reports distinct device error triggers depending on the failure: `interface_loading_failed` when an interface fails to load and `resend_interface_properties_failed` when sending the properties to the device itself fails. Previously, both cases were incorrectly reported as `resend_interface_properties_failed`.
+
 ## [1.3.2] - 2026-07-14
 
 ## Fixed

FILE: apps/astarte_appengine_api/mix.lock (+1 -1)
@@ -4,7 +4,7 @@
   "astarte_core": {:git, "https://github.com/astarte-platform/astarte_core.git", "e6c245d2642cf499ba29a78d08e5d3d6de47badc", [tag: "v1.3.2"]},
   "astarte_rpc": {:hex, :astarte_rpc, "1.2.0", "dcef7434bf3f19ff30ff0bc245ef2d3b0f9abeb369405713cfd28916a5635926", [:mix], [{:amqp, "~> 3.3", [hex: :amqp, repo: "hexpm", optional: false]}, {:castore, "~> 1.0.0", [hex: :castore, repo: "hexpm", optional: false]}, {:protobuf, "~> 0.12", [hex: :protobuf, repo: "hexpm", optional: false]}, {:skogsra, "~> 2.2", [hex: :skogsra, repo: "hexpm", optional: false]}], "hexpm", "8470ed2f116fa8c9d70845f295f3738529aff123d9d3f98cbfa37107314763cb"},
   "bunt": {:hex, :bunt, "1.0.0", "081c2c665f086849e6d57900292b3a161727ab40431219529f13c4ddcf3e7a44", [:mix], [], "hexpm", "dc5f86aa08a5f6fa6b8096f0735c4e76d54ae5c9fa2c143e5a1fc7c1cd9bb6b5"},
-  "castore": {:hex, :castore, "1.0.19", "6903cabdfd9d1af46454126e7c8385186659dd33ecfb74a885cae52221ad6109", [:mix], [], "hexpm", "3669e6cab13f54c2df26b3e6833745d647f35b6e30d8ddd5975df0d5c842ca98"},
+  "castore": {:hex, :castore, "1.0.20", "455e48f7115eca98c9f2b0e7a152b5a2e8f2a8a4f964c96e95bd31645ee5fa59", [:mix], [], "hexpm", "940eafbfd8b14bee649f083bc11b3b54ec555b54c3e4ea8213351ff6fee39c10"},
   "certifi": {:hex, :certifi, "2.15.0", "0e6e882fcdaaa0a5a9f2b3db55b1394dba07e8d6d9bcad08318fb604c6839712", [:rebar3], [], "hexpm", "b147ed22ce71d72eafdad94f055165c1c182f61a2ff49df28bcc71d1d5b94a60"},
   "con_cache": {:hex, :con_cache, "1.1.1", "9f47a68dfef5ac3bbff8ce2c499869dbc5ba889dadde6ac4aff8eb78ddaf6d82", [:mix], [{:telemetry, "~> 1.0", [hex: :telemetry, repo: "hexpm", optional: false]}], "hexpm", "1def4d1bec296564c75b5bbc60a19f2b5649d81bfa345a2febcc6ae380e8ae15"},
   "cors_plug": {:hex, :cors_plug, "2.0.3", "316f806d10316e6d10f09473f19052d20ba0a0ce2a1d910ddf57d663dac402ae", [:mix], [{:plug, "~> 1.8", [hex: :plug, repo: "hexpm", optional: false]}], "hexpm", "ee4ae1418e6ce117fc42c2ba3e6cbdca4e95ecd2fe59a05ec6884ca16d469aea"},

FILE: apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex (+35 -2)
@@ -343,12 +343,22 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.ControlHandler do
 
         {:ok, new_state}
 
+      {:error, :interface_loading_failed} ->
+        # Track interface loading failure
+        :telemetry.execute(
+          [:astarte, :data_updater_plant, :control_handler, :properties_resend],
+          %{duration: System.monotonic_time() - resend_start},
+          %{realm: state.realm, result: "interface_loading_failed"}
+        )
+
+        interface_loading_error(state, message_id, timestamp)
+
       {:error, :sending_properties_to_interface_failed} ->
         # Track interface send failure
         :telemetry.execute(
           [:astarte, :data_updater_plant, :control_handler, :properties_resend],
           %{duration: System.monotonic_time() - resend_start},
-          %{realm: state.realm, result: "interface_failed"}
+          %{realm: state.realm, result: "sending_properties_to_interface_failed"}
         )
 
         sending_properties_error(state, message_id, timestamp)
@@ -386,8 +396,31 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.ControlHandler do
     new_state
   end
 
+  defp interface_loading_error(state, message_id, timestamp) do
+    Logger.warning("Cannot load interface while resending properties",
+      tag: "interface_loading_failed"
+    )
+
+    {:ok, new_state} = Core.Device.ask_clean_session(state, timestamp)
+    MessageTracker.discard(new_state.message_tracker, message_id)
+
+    :telemetry.execute(
+      [:astarte, :data_updater_plant, :data_updater, :discarded_message],
+      %{},
+      %{realm: new_state.realm}
+    )
+
+    Core.Trigger.execute_device_error_triggers(
+      new_state,
+      "interface_loading_failed",
+      timestamp
+    )
+
+    new_state
+  end
+
   defp sending_properties_error(state, message_id, timestamp) do
-    Logger.warning("Cannot resend properties to interface",
+    Logger.warning("Cannot send properties to device.",
       tag: "resend_interface_properties_failed"
     )
 

FILE: apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/device.ex (+28 -11)
@@ -424,9 +424,12 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.Device do
            :ok <- resend_all_interface_properties(new_state, interface_descriptor) do
         {:cont, {:ok, new_state}}
       else
-        {:error, :interface_loading_failed} ->
-          Logger.warning("Failed #{interface} interface loading.")
-          {:halt, {:error, :sending_properties_to_interface_failed}}
+        {:error, reason} = error ->
+          Logger.warning(
+            "Error while sending properties to device: #{inspect(reason)} on interface #{interface}"
+          )
+
+          {:halt, error}
       end
     end)
   end
@@ -437,14 +440,28 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.Device do
        ) do
     encoded_device_id = Device.encode_device_id(device_id)
 
-    Core.Interface.each_interface_mapping(mappings, interface_descriptor, fn mapping ->
-      resend_interface_mapping_properties(
-        realm,
-        device_id,
-        encoded_device_id,
-        interface_descriptor,
-        mapping
-      )
+    mappings
+    |> Enum.filter(fn {_endpoint_id, mapping} ->
+      mapping.interface_id == interface_descriptor.interface_id
+    end)
+    |> Enum.reduce_while(:ok, fn {_endpoint_id, mapping}, :ok ->
+      case resend_interface_mapping_properties(

$ gh search code --owner astarte-platform 'device_deletion' --filename '*.ex' -L 20 2>&1 | head -40
(no output)

✱ Grep "UpstreamErrorName|interface_loading_failed|resend_interface_properties_failed" in . · 31 matches
The v1.3.3 fix is about distinguishing `interface_loading_failed` vs `resend_interface_properties_failed` during properties resend. Let me check Astrate's `UpstreamErrorName` mapping and where properties resend happens.

$ gh search code --owner astarte-platform 'device_created' --filename '*.ex' -L 20 2>&1 | head -40
(no output)

→ Read internal/engine/triggers/errorname.go

$ gh search code --owner astarte-platform 'device_registration' --filename '*.ex' -L 20 2>&1 | head -40
(no output)

→ Read internal/engine/cache.go [offset=340, limit=80]
Astrate already carries both names and routes interface problems to `interface_loading_failed` — already aligned with 1.3.3's semantics. Let me check Astrate's properties-resend code to see what it emits, and search upstream for the new trigger types properly.

$ gh search code --owner astarte-platform 'device_deletion_started' -L 20 2>&1 | head -40; echo "===REG==="; gh search code --owner astarte-platform 'device_registered' -L 20 2>&1 | head -20
astarte-platform/astarte:apps/astarte_realm_management/lib/astarte_realm_management/triggers_handler/triggers_handler.ex: def device_deletion_started(realm_name, device_id) do
astarte-platform/astarte:apps/astarte_realm_management/lib/astarte_realm_management/triggers_handler/triggers_handler.ex: event_key = :on_device_deletion_started
astarte-platform/astarte:libs/astarte_generators/lib/astarte/core/generators/triggers/simple_events/device_deletion_started_event.ex: @spec device_deletion_started_event() :: StreamData.t(DeviceDeletionStartedEvent.t())
astarte-platform/astarte:libs/astarte_generators/lib/astarte/core/generators/triggers/simple_events/device_deletion_started_event.ex: def device_deletion_started_event, do: constant(%DeviceDeletionStartedEvent{})
astarte-platform/astarte_generators:lib/astarte/core/generators/triggers/simple_events/device_deletion_started_event.ex: @spec device_deletion_started_event() :: StreamData.t(DeviceDeletionStartedEvent.t())
astarte-platform/astarte_generators:lib/astarte/core/generators/triggers/simple_events/device_deletion_started_event.ex: def device_deletion_started_event, do: constant(%DeviceDeletionStartedEvent{})
astarte-platform/astarte:tools/astarte-e2e-next/astarte-e2e/examples/regen.rs: ("device_deletion_started_event", "device_deletion_started"),
astarte-platform/astarte:doc/pages/architecture/060-triggers.md: - `device_deletion_started`: triggered when the deletion of a device is started.
astarte-platform/astarte:doc/pages/architecture/060-triggers.md: "type": "device_deletion_started"
astarte-platform/astarte:libs/astarte_core/lib/astarte_core/triggers/simple_events/simple_event.proto: import "lib/astarte_core/triggers/simple_events/device_deletion_started_event.proto";
astarte-platform/astarte:libs/astarte_core/lib/astarte_core/triggers/simple_events/simple_event.proto: DeviceDeletionStartedEvent device_deletion_started_event = 21;
astarte-platform/astarte:astarte-dashboard/src/astarte-client/types/events/AstarteDeviceDeletionStartedEvent.ts: type: 'device_deletion_started';
astarte-platform/astarte:astarte-dashboard/src/astarte-client/types/events/AstarteDeviceDeletionStartedEvent.ts: type: yup.string().oneOf(['device_deletion_started']).required(),
astarte-platform/astarte:apps/astarte_realm_management/lib/astarte_realm_management/devices/devices.ex: TriggersHandler.device_deletion_started(realm_name, device_id)
astarte-platform/astarte:libs/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/device_trigger.proto: DEVICE_DELETION_STARTED = 10;
astarte-platform/astarte:libs/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/device_trigger.pb.ex: field :DEVICE_DELETION_STARTED, 10
astarte-platform/astarte:apps/astarte_realm_management/lib/astarte_realm_management/application.ex: trigger_types = [:DEVICE_DELETION_STARTED, :DEVICE_DELETION_FINISHED]
astarte-platform/astarte:libs/astarte_events/lib/astarte_events/triggers/core.ex: | :on_device_deletion_started
astarte-platform/astarte:libs/astarte_events/lib/astarte_events/triggers/core.ex: DEVICE_DELETION_STARTED: :on_device_deletion_started,
astarte-platform/astarte:libs/astarte_core/lib/astarte_core/triggers/simple_events/simple_event.pb.ex: field :device_deletion_started_event, 21,
astarte-platform/astarte:tools/astarte-e2e-next/astarte-e2e/src/simple_events.rs: #[serde(rename = "device_deletion_started")]
astarte-platform/astarte:libs/astarte_generators/lib/astarte/core/generators/triggers/simple_events/simple_event.ex: {:device_deletion_started_event, device_deletion_started_event()},
astarte-platform/astarte:astarte-dashboard/src/astarte-client/models/Trigger/index.ts: | 'device_deletion_started'
astarte-platform/astarte:astarte-dashboard/src/astarte-client/models/Trigger/index.ts: 'device_deletion_started',
astarte-platform/astarte:astarte-dashboard/cypress/e2e/trigger_builder.cy.js: device_deletion_started: 'Device Deletion Started',
astarte-platform/astarte:libs/astarte_core/lib/astarte_core/triggers/simple_trigger_config.ex: "device_deletion_started" => :DEVICE_DELETION_STARTED,
astarte-platform/astarte_core:lib/astarte_core/triggers/simple_events/simple_event.proto: import "lib/astarte_core/triggers/simple_events/device_deletion_started_event.proto";
astarte-platform/astarte_core:lib/astarte_core/triggers/simple_events/simple_event.proto: DeviceDeletionStartedEvent device_deletion_started_event = 21;
astarte-platform/astarte-dashboard:src/astarte-client/types/events/AstarteDeviceDeletionStartedEvent.ts: type: 'device_deletion_started';
astarte-platform/astarte-dashboard:src/astarte-client/types/events/AstarteDeviceDeletionStartedEvent.ts: type: yup.string().oneOf(['device_deletion_started']).required(),
===REG===
astarte-platform/astarte:libs/astarte_generators/lib/astarte/core/generators/triggers/simple_events/device_registered_event.ex: @spec device_registered_event() :: StreamData.t(DeviceRegisteredEvent.t())
astarte-platform/astarte:libs/astarte_generators/lib/astarte/core/generators/triggers/simple_events/device_registered_event.ex: def device_registered_event, do: constant(%DeviceRegisteredEvent{})
astarte-platform/astarte_generators:lib/astarte/core/generators/triggers/simple_events/device_registered_event.ex: @spec device_registered_event() :: StreamData.t(DeviceRegisteredEvent.t())
astarte-platform/astarte_generators:lib/astarte/core/generators/triggers/simple_events/device_registered_event.ex: def device_registered_event, do: constant(%DeviceRegisteredEvent{})
astarte-platform/astarte:apps/astarte_pairing/lib/astarte_pairing/agent/agent.ex: event_key = :on_device_registered
astarte-platform/astarte:apps/astarte_pairing/lib/astarte_pairing/agent/agent.ex: event_type = :device_registered_event
astarte-platform/astarte:tools/astarte-e2e-next/astarte-e2e/examples/regen.rs: ("device_registered_event", "device_registered"),
astarte-platform/astarte:apps/astarte_pairing/lib/astarte_pairing/application.ex: {Astarte.RPC.Triggers.Client, types: [:DEVICE_REGISTERED]},
astarte-platform/astarte:doc/pages/architecture/060-triggers.md: - `device_registered`: triggered when a device is registered to Astarte.
astarte-platform/astarte:doc/pages/architecture/060-triggers.md: "type": "device_registered"
astarte-platform/astarte:libs/astarte_core/lib/astarte_core/triggers/simple_events/simple_event.proto: import "lib/astarte_core/triggers/simple_events/device_registered_event.proto";
astarte-platform/astarte:libs/astarte_core/lib/astarte_core/triggers/simple_events/simple_event.proto: DeviceRegisteredEvent device_registered_event = 20;
astarte-platform/astarte:astarte-dashboard/src/astarte-client/types/events/AstarteDeviceRegistrationEvent.ts: type: 'device_registered';
astarte-platform/astarte:astarte-dashboard/src/astarte-client/types/events/AstarteDeviceRegistrationEvent.ts: type: yup.string().oneOf(['device_registered']).required(),
astarte-platform/astarte:libs/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/device_trigger.proto: DEVICE_REGISTERED = 9;
astarte-platform/astarte:libs/astarte_core/lib/astarte_core/triggers/simple_triggers_protobuf/device_trigger.pb.ex: field :DEVICE_REGISTERED, 9
astarte-platform/astarte:libs/astarte_events/lib/astarte_events/triggers/core.ex: | :on_device_registered
astarte-platform/astarte:libs/astarte_events/lib/astarte_events/triggers/core.ex: DEVICE_REGISTERED: :on_device_registered,
astarte-platform/astarte:libs/astarte_core/lib/astarte_core/triggers/simple_events/simple_event.pb.ex: field :device_registered_event, 20,
astarte-platform/astarte:tools/astarte-e2e-next/astarte-e2e/src/simple_events.rs: #[serde(rename = "device_registered")]

Upstream v1.3.0 added three new trigger `on` values: `device_registered`, `device_deletion_started`, `device_deletion_finished`. Astrate only knows `device_connected`/`device_disconnected`. Let me confirm how Astrate handles an unknown `on` value.
→ Read internal/engine/triggers/match.go

$ rg -n 'OnDeviceRegistered|OnDeviceDeletionStarted|OnDeviceDeletionFinished|device_registered|device_deletion_started|device_deletion_finished' --glob '*.go' -l
internal/engine/stream/bus.go
internal/engine/engine.go
internal/engine/triggers/match_test.go
internal/engine/triggers/match.go
internal/engine/triggers/events_test.go
internal/engine/triggers/events.go
internal/realm/service.go

The trigger *constants* already exist. Let me check whether the events actually get fired — search for where these are wired.
$ rg -n 'fireDevice|DeviceEvent\{' internal/engine --glob '*.go' | head -40
internal/engine/engine.go:75:	e.onDeviceError = e.fireDeviceError
internal/engine/engine.go:284:// fireDevice evaluates device triggers for one device-scoped event.
internal/engine/engine.go:285:func (e *Engine) fireDevice(rs *realmSchema, id deviceid.ID, at time.Time, match triggers.DeviceEvent, body any) {
internal/engine/engine.go:302:// fireDeviceError is the onDeviceError observer: every pipeline rejection
internal/engine/engine.go:305:func (e *Engine) fireDeviceError(m broker.InboundMessage, reason, detail string) {
internal/engine/engine.go:322:	e.fireDevice(rs, m.DeviceID, m.ReceivedAt,
internal/engine/engine.go:323:		triggers.DeviceEvent{DeviceID: m.DeviceID.String(), On: triggers.OnDeviceError},
internal/engine/engine.go:347:		e.fireDevice(rs, ev.DeviceID, ev.At,
internal/engine/engine.go:348:			triggers.DeviceEvent{DeviceID: deviceID, On: triggers.OnDeviceConnected},
internal/engine/engine.go:365:		e.fireDevice(rs, ev.DeviceID, ev.At,
internal/engine/engine.go:366:			triggers.DeviceEvent{DeviceID: deviceID, On: triggers.OnDeviceDisconnected},
internal/engine/engine.go:386:	e.fireDevice(rs, id, at,
internal/engine/engine.go:387:		triggers.DeviceEvent{DeviceID: hwID, On: triggers.OnDeviceRegistered},
internal/engine/engine.go:422:	e.fireDevice(rs, id, at, triggers.DeviceEvent{DeviceID: hwID, On: on}, body)
internal/engine/triggers/match_test.go:244:			ev: DeviceEvent{DeviceID: "x", On: OnDeviceConnected}, want: true},
internal/engine/triggers/match_test.go:246:			ev: DeviceEvent{DeviceID: "x", On: OnDeviceDisconnected}, want: false},
internal/engine/triggers/match_test.go:251:			ev: DeviceEvent{DeviceID: "x", On: OnDeviceDeletionStarted}, want: true},
internal/engine/triggers/match_test.go:256:			ev: DeviceEvent{DeviceID: "x", On: OnDeviceDeletionFinished}, want: true},
internal/engine/triggers/match_test.go:261:			ev: DeviceEvent{DeviceID: "x", On: OnDeviceDeletionFinished}, want: false},
internal/engine/triggers/match_test.go:263:			ev: DeviceEvent{DeviceID: "f0VMRgIBAQAAAAAAAAAAAA", On: OnDeviceDisconnected}, want: true},
internal/engine/triggers/match_test.go:265:			ev: DeviceEvent{DeviceID: "AAAAAAAAAAAAAAAAAAAAAQ", On: OnDeviceDisconnected}, want: false},
internal/engine/triggers/match_test.go:267:			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.Sensors", Major: 2}, want: true},
internal/engine/triggers/match_test.go:269:			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.Sensors", Major: 1}, want: false},
internal/engine/triggers/match_test.go:271:			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.Other", Major: 2}, want: false},
internal/engine/triggers/match_test.go:273:			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceRemoved, Interface: "com.example.Sensors", Major: 2}, want: false},
internal/engine/triggers/match_test.go:492:			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.Other", Major: 2}, want: true},
internal/engine/triggers/match_test.go:494:			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.Other", Major: 1}, want: false},
internal/engine/triggers/match_test.go:496:			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceRemoved, Interface: "org.a.B", Major: 3}, want: true},
internal/engine/triggers/match_test.go:498:			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceRemoved, Interface: "org.a.B", Major: 4}, want: false},
internal/engine/triggers/match_test.go:500:			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.X", Major: 9}, want: true},
internal/engine/triggers/match_test.go:502:			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.X", Major: 123}, want: true},
internal/engine/triggers/match_test.go:504:			ev: DeviceEvent{DeviceID: "x", On: OnDeviceConnected}, want: true},
internal/engine/triggers/transient_test.go:18:	if !tg.MatchesDevice(DeviceEvent{DeviceID: "abc123", On: OnDeviceConnected}) {
internal/engine/bus_events_test.go:96:		rig.e.fireDeviceError(broker.InboundMessage{
internal/engine/bus_events_test.go:127:		rig.e.fireDeviceError(broker.InboundMessage{
internal/engine/introspection.go:76:	e.fireDevice(realm, m.DeviceID, at,
internal/engine/introspection.go:77:		triggers.DeviceEvent{DeviceID: m.DeviceID.String(), On: triggers.OnIncomingIntrospection},
internal/engine/introspection.go:81:			e.fireDevice(realm, m.DeviceID, at,
internal/engine/introspection.go:82:				triggers.DeviceEvent{DeviceID: m.DeviceID.String(), On: triggers.OnInterfaceAdded, Interface: name, Major: v.Major},
internal/engine/introspection.go:87:		e.fireDevice(realm, m.DeviceID, at,

→ Read internal/engine/engine.go [offset=370, limit=70]
→ Read internal/engine/introspection.go
Astrate has already implemented all three v1.3.0 trigger types. Let me check the remaining candidates: empty-introspection allowance (v1.3.0) and the properties-resend error naming (v1.3.3).

$ rg -n 'consumer/properties|producer/properties|resend' internal/engine --glob '*.go' | head -30
internal/engine/mqttv1_parity_test.go:101:		t.Fatalf("resend publications: %d, want 1", len(pubs))
internal/engine/mqttv1_parity_test.go:105:		t.Errorf("resend qos = %d, want 2", pub.qos)
internal/engine/mqttv1_parity_test.go:108:		t.Error("resend published with retain set")
internal/engine/mqttv1_parity_test.go:149:		t.Fatal(`no binary element with key "v" found in the raw resend payload`)
internal/engine/cache.go:62:	// PurgeDeviceOwnedExcept implements the `producer/properties` resync
internal/engine/cache.go:66:	// (emptyCache resend + `consumer/properties` payload, docs/DESIGN.md §3.4).
internal/engine/cache.go:365:	// consumer/properties payloads ("zlib" or "plaintext"); set by the
internal/engine/cache.go:404:// setPurgeCompression stores the device's preferred consumer/properties
internal/engine/cache.go:413:// purgeCompressionFormat returns the device's preferred consumer/properties
internal/engine/engine_e2e_test.go:351:// consumer/properties purge list — both decoded and asserted (docs/DESIGN.md
internal/engine/engine_e2e_test.go:372:	// ...and publishes the consumer/properties purge list naming it.
internal/engine/engine_e2e_test.go:374:	cp := dev.WaitForMessage(t, 5*time.Second, "consumer/properties listing "+want, func(m testutil.ServerMessage) bool {
internal/engine/engine_e2e_test.go:375:		if m.Topic != dev.Base()+"/control/consumer/properties" {
internal/engine/engine_e2e_test.go:386:		t.Errorf("consumer/properties entries = %v, want [%s]", got, want)
internal/engine/engine_e2e_test.go:390:// testE2EProducerProperties: control/producer/properties is the device's
internal/engine/router.go:96:	// bg tracks asynchronous control sends (consumer/properties after a
internal/engine/capabilities.go:27:// server sends consumer/properties payloads compressed (zlib, the default) or
internal/engine/capabilities.go:89:// purgeCompressionFor returns the consumer/properties compression format for
internal/engine/serverdata.go:208:// the `consumer/properties` purge message follows so offline-window state
internal/engine/bus_events_test.go:27:	rig.e.broker = &fakePort{} // handleLifecycle's connect case publishes consumer/properties
internal/engine/topics_test.go:55:		{"control/producer/properties", kindControl, "producer/properties"},
internal/engine/topics.go:61:// "producer/properties"); a bare "control" yields an empty subpath, which
internal/engine/capabilities_test.go:173:	// emptyCache triggers consumer/properties send.
internal/engine/capabilities_test.go:181:	purges := port.publishedTo(base + "/control/consumer/properties")
internal/engine/capabilities_test.go:183:		t.Fatalf("consumer/properties messages: %d, want 1", len(purges))
internal/engine/capabilities_test.go:216:	purges := port.publishedTo(base + "/control/consumer/properties")
internal/engine/capabilities_test.go:218:		t.Fatalf("consumer/properties messages: %d, want 1", len(purges))
internal/engine/serverdata_test.go:211:// consumer/properties purge follow-up (docs/DESIGN.md §3.4).
internal/engine/serverdata_test.go:245:	purges := port.publishedTo(base + "/control/consumer/properties")
internal/engine/control_test.go:20:// goldenFrame loads the frozen producer/properties blob: framed with an

Empty introspection is already accepted (`introspection.go:101`). Let me check Astrate's properties-resend-on-connect path (the v1.3.3 fix's subject).
→ Read internal/engine/serverdata.go [limit=70, offset=180]
Let me look at the v1.3.2 crash-fix diff and Astrate's properties-resend on connect (the lifecycle path).

$ rg -n 'consumer/properties|PurgeDeviceOwnedExcept|handleLifecycle|OnConnect' internal/engine/lifecycle.go internal/engine/cache.go 2>/dev/null | head -30; echo =====; rg -n 'consumer/properties|Purging|purge' internal/engine/lifecycle.go 2>/dev/null | head
internal/engine/cache.go:62:	// PurgeDeviceOwnedExcept implements the `producer/properties` resync
internal/engine/cache.go:64:	PurgeDeviceOwnedExcept(ctx context.Context, realmID int16, deviceID deviceid.ID, keep []store.PropertyRef) (int64, error)
internal/engine/cache.go:66:	// (emptyCache resend + `consumer/properties` payload, docs/DESIGN.md §3.4).
internal/engine/cache.go:365:	// consumer/properties payloads ("zlib" or "plaintext"); set by the
internal/engine/cache.go:404:// setPurgeCompression stores the device's preferred consumer/properties
internal/engine/cache.go:413:// purgeCompressionFormat returns the device's preferred consumer/properties
=====

$ gh api 'repos/astarte-platform/astarte/compare/v1.3.1...v1.3.2' --jq '.commits[] | "\(.sha[0:8]) \(.commit.message | split("\n")[0])"'
0808aa58 fix: Avoid crash on invalid properties message
41b4f4e8 chore: prepare v1.3.2

✗ Read internal/engine/lifecycle.go failed
Error: File not found: /root/astrate-mule/internal/engine/lifecycle.go

$ gh api 'repos/astarte-platform/astarte/commits/0808aa58' --jq '.files[] | "FILE: \(.filename)\n\(.patch // "no patch")"' 2>&1 | head -80
FILE: apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex
@@ -69,41 +69,40 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.ControlHandler do
 
     decompression_start = System.monotonic_time()
 
-    case decode_payload(state, payload) do
-      {:ok, decoded_payload} ->
-        # Track successful decompression
-        :telemetry.execute(
-          [:astarte, :data_updater_plant, :control_handler, :payload_decompression],
-          %{
-            duration: System.monotonic_time() - decompression_start,
-            compressed_size: byte_size(payload),
-            uncompressed_size: byte_size(decoded_payload)
-          },
-          %{realm: new_state.realm, result: "success"}
-        )
-
-        :ok = Core.Device.prune_device_properties(new_state, decoded_payload, timestamp_ms)
-        MessageTracker.ack_delivery(new_state.message_tracker, message_id)
+    with {:ok, decoded_payload} <- decode_payload(state, payload),
+         :ok <- Core.Device.prune_device_properties(new_state, decoded_payload, timestamp_ms) do
+      # Track successful decompression
+      :telemetry.execute(
+        [:astarte, :data_updater_plant, :control_handler, :payload_decompression],
+        %{
+          duration: System.monotonic_time() - decompression_start,
+          compressed_size: byte_size(payload),
+          uncompressed_size: byte_size(decoded_payload)
+        },
+        %{realm: new_state.realm, result: "success"}
+      )
 
-        # Track properties prune with payload
-        :telemetry.execute(
-          [:astarte, :data_updater_plant, :control_handler, :properties_prune],
-          %{
-            duration: System.monotonic_time() - start_time,
-            payload_size: byte_size(payload)
-          },
-          %{realm: new_state.realm, prune_type: "with_payload"}
-        )
+      MessageTracker.ack_delivery(new_state.message_tracker, message_id)
 
+      # Track properties prune with payload
+      :telemetry.execute(
+        [:astarte, :data_updater_plant, :control_handler, :properties_prune],
         %{
-          new_state
-          | total_received_msgs: new_state.total_received_msgs + 1,
-            total_received_bytes:
-              new_state.total_received_bytes + byte_size(payload) +
-                byte_size("/producer/properties")
-        }
-
-      :error ->
+          duration: System.monotonic_time() - start_time,
+          payload_size: byte_size(payload)
+        },
+        %{realm: new_state.realm, prune_type: "with_payload"}
+      )
+
+      %{
+        new_state
+        | total_received_msgs: new_state.total_received_msgs + 1,
+          total_received_bytes:
+            new_state.total_received_bytes + byte_size(payload) +
+              byte_size("/producer/properties")
+      }
+    else
+      _ ->
         # Track failed decompression
         :telemetry.execute(
           [:astarte, :data_updater_plant, :control_handler, :payload_decompression],
FILE: apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/device.ex
@@ -402,15 +402,15 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.Device do
   end
 
   def prune_device_properties(state, decoded_payload, timestamp) do
-    {:ok, paths_set} =

✱ Grep "func.*handleConnect|func.*(e \*Engine) .*[Cc]onnect|consumer/properties" in internal/engine · 32 matches
→ Read internal/engine/engine.go [limit=50, offset=325]
→ Read internal/engine/control.go [limit=55, offset=183]
I have the full picture now. Let me check the current `.mule/todo.md` and `.mule/for-giulio.md` to see what's already queued and how prior runs framed this.
→ Read .mule/todo.md
→ Read .mule/for-giulio.md
The for-giulio.md already carries the v1.3.2/v1.3.3 compat-note wording (lines 77–79) and the 2026-09-06 milestone escalation covers the 1.3-surface status, and `compat-note-v1.3.3` is already queued (blocked, todo.md:104). Let me verify one last candidate — Astrate's handling of an invalid `producer/properties` purge payload (v1.3.2's crash fix) — before writing the proposal.
→ Read internal/engine/control.go [limit=120, offset=1]
Astrate's `producer/properties` handler already rejects malformed payloads gracefully (`control.go:86`) — the v1.3.2 crash fix has no Astrate counterpart. Let me read `resendServerProperties` to pin down the exact resend-failure path for the probe line.
→ Read internal/engine/control.go [limit=45, offset=136]
All candidates verified. Here's the picture:

- **v1.3.0 wire items already implemented**: capabilities (incl. `purge_properties_compression_format`), empty-introspection allowance, the three new trigger types (`device_registered`, `device_deletion_started/finished`), realm-scoped pairing health.
- **v1.3.2** crash fix: no Astrate counterpart (producer/properties rejects malformed payloads).
- **v1.3.3** (empty release body, but one substantive fix #2119): upstream now fires **distinct device_error triggers** on the properties-resend path — `interface_loading_failed` vs `resend_interface_properties_failed`. Astrate's `resendServerProperties` (`control.go:144`) and `sendConsumerProperties` (`control.go:191`) skip bad rows or log a Warn (`engine.go:359-362`) but **never fire a device_error trigger**.
- **compat-note-v1.3.3** (todo.md:104) and its for-giulio wording (for-giulio.md:77-79) already exist — re-proposing would duplicate.

Let me check git state is clean and append the one new task line.

$ git status --short
(no output)

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -106,4 +106,5 @@
 - [x] flow-validate-dead-source-sink-recompute: in `Pipeline.Validate`, the first `hasSource`/`hasSink` computation (pipeline.go:117-126) is overwritten by the second (pipeline.go:128-147) and its result discarded — dead code (the comment at pipeline.go:127 documents it). Delete the first block; behaviour is unchanged. Tests: existing `TestPipeline_Validate` passes.
 - [x] docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` route to docs/api/astarte_realm_management_api.yaml — it exists in code (internal/realm/http.go:51, handler getDatastreamMaximumStorageRetention at http.go:135, returns 200 with a data-envelope retention value via Service.GetDatastreamMaximumStorageRetention) but the spec jumps straight from config/device_registration_limit (yaml:428) to /version. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] housekeeping-tests [auto]: add a container-free unit test file for `internal/housekeeping` (service.go, 262 lines, zero test coverage). Cover two rules: (1) `CreateRealm` validation — blank name, blank JWT key, negative regLimit, negative retention, and default retention injection when `WithDefaultDatastreamMaximumStorageRetention` is set and the caller passes nil; (2) `DeleteRealm` gating — `ErrDeletionDisabled` when `WithRealmDeletionDisabled(true)`, `ErrConnectedDevicesPresent` when `DeviceStats` reports connected > 0. Use httptest-style fakes for the store and sealer (no Docker).
+- [ ] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch.

→ Read .mule/todo.md [offset=108, limit=10]
Task line appended cleanly. Working tree still shows no untracked files beyond the expected `.mule/todo.md` modification. Let me confirm.

$ git status --short
M .mule/todo.md

---

**Done:** Upstream watch run (v1.3.3 vs `APICompatVersion = "1.2.2"`): investigated the three v1.3.3/1.3.x wire-visible candidates; one new probe task appended to `.mule/todo.md`; compat-note line not re-proposed (already queued at line 104, wording already in for-giulio.md:77-79 and superseded by the 2026-09-06 milestone escalation).

**Files:** `.mule/todo.md:110`

**Verified:** `git status --short` → only `.mule/todo.md` modified.

**Unsure:** whether `sendConsumerProperties` (the purge-list path) is considered a "resend" in the same semantic category as upstream's per-interface value-resend; the upstream fix (device.ex `resend_all_interface_properties`) operates on per-interface property values, while Astrate's counterpart `resendServerProperties` (control.go:144) does the same per-interface value work but `sendConsumerProperties` (control.go:191) emits the purge-list which has different failure semantics (silently drops unknown interfaces at line 199 rather than failing per-interface). The probe task's wording covers both.

**Follow-ups:** one new task line (`probe-props-resend-error-triggers` at todo.md:110).
