slug: recipe-astarte-upstream
verdict: proposed
at:  59940df
ran: 2026-09-12T18:07:35Z on DietPi in 375s


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

→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
Sixth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.4

### Fixed

- [astarte_data_updater_plant] Make object values only validate type on selected interface
- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated

The latest release is v1.4.0-rc.5, and Astrate currently targets 1.2.2. Upstream is ahead. Let me check the release notes for the relevant recent releases.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.4") | .body'
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

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.3") | .body'
Fourth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.2

### Added

- [secrets] Configurable base namespace for Vault using `ASTARTE_VAULT_BASE_NAMESPACE`: all created namespaces will be children of the base namespace.

### Fixed

- [fdo] Support for Vault and Vault managed
- [fdo] Properly handle chunked service info messages
- [astarte_data_updater_plant] Avoid clogging AMQP queues when receiving invalid device messages
- [astarte_data_updater_plant] Restore correct type validation checks for incoming data

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.2") | .body'
Third Astarte v1.4.0 release candidate.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | "\(.tag_name) \(.body | length) chars"' 2>&1
v1.3.3 0 chars

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | keys' 2>&1
["assets","assets_url","author","body","created_at","draft","html_url","id","immutable","name","node_id","prerelease","published_at","tag_name","tarball_url","target_commitish","updated_at","upload_url","url","zipball_url"]

$ gh api repos/astarte-platform/astarte/releases --jq '.[].tag_name' 2>&1
v1.4.0-rc.5
v1.3.3
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
v1.3.2
v1.3.1
v1.3.0
v1.2.2
v1.4.0-rc.0
v1.3.0-rc.2
v1.2.2-rc.0
v1.2.1
v1.2.1-rc.1
v1.3.0-rc.1
v1.3.0-rc.0
v1.2.1-rc.0
v1.2.1-alpha.0
v1.2.0
v1.2.0-rc.0
v1.0.6
v1.1.1
v1.0.5
v1.1.0
v1.1.0-rc.0
v1.1.0-alpha.0
v1.0.4
v1.0.3
v1.0.2

The latest tag is v1.4.0-rc.5; Astrate targets 1.2.2. Upstream is ahead. Let me look at the wire-relevant changes more closely.
$ gh search code --owner astarte-platform 'object values only validate type on selected interface' --json path,repository,url 2>&1 | head -30
[{"path":"CHANGELOG.md","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/CHANGELOG.md"}]

$ gh search code --owner astarte-platform 'binaryblob data is correctly validated' --json path,repository,url 2>&1 | head -30
[{"path":"CHANGELOG.md","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/CHANGELOG.md"}]

$ gh search code --owner astarte-platform 'Avoid clogging AMQP queues when receiving invalid device messages' --json path,repository,url 2>&1 | head -30
[{"path":"CHANGELOG.md","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/CHANGELOG.md"}]

$ gh api 'repos/astarte-platform/astarte/contents/CHANGELOG.md?ref=v1.4.0-rc.5' --jq '.content' 2>&1 | base64 -d 2>/dev/null | head -80
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

## [1.4.0-rc.5] - 2026-08-20

### Fixed

- [astarte_data_updater_plant] Make object values only validate type on selected interface
- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated

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

$ gh api 'repos/astarte-platform/astarte/contents/CHANGELOG.md?ref=v1.4.0-rc.5' --jq '.content' 2>&1 | base64 -d 2>/dev/null | sed -n '80,220p'
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

## [1.3.3] - 2026-08-07

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

Now let me verify the wire-relevant changes by looking at the diffs. The key candidates: the 1.3.3 device-error distinction, 1.3.0 empty-introspection + capabilities + `purge_properties_compression_format`, and 1.4.0-rc.5 validation fixes.
$ gh api 'repos/astarte-platform/astarte/compare/v1.2.2...v1.3.3' --jq '.commits[] | "\(.sha[0:8]) \(.commit.message | split("\n")[0])"' 2>&1 | head -60
06030b5e chore(pairing): fix typespec for proveovhdr (#1691)
fdc51071 refactor(pairing): have service info start use its module for decode
28632228 fix(pairing): use correct fallback controller and corid (#1693)
1d74b75c fix(pairing): service info start message
e54d2348 feat(fdo): log internal server errors
7f7ba6a1 fix(pairing): service info messages
850e739d fix(pairing): fix done messages
2ca8213f fix(pairing): return 101 for invalid message
3392ee28 build(docker): rewrite all paths to pairing
af2926d1 feat(pairing): add realm scoped health check
63770e24 chore: forward port release-1.3
db412687 chore(pairing): separate realm and fdo scopes
043b15a9 fix(pairing): don't crash on msg Done error
3cea9668 test(pairing): FDO controller tests
04d885e8 ci: add fdo end to end tests
38c1ce90 chore(pairing): Rendezvous info struct
7e47b60a fix(pairing): only accept prove device messages with valid nonces  (#1704)
d6a1f5e5 style: format files
0dbd8a74 chore: remove leftover api services mentions
5d12121e chore(pairing): added correct error logs for ov_next_entry
193736e8 feat: release-1.3 support for astarte_import and astarte_export (#1596)
35b377e0 ci: optimize build caching with matrix strategy (#1715)
057bac4d chore(pairing): added correct error logs for service_info (#1714)
f2b9ef5f chore(pairing): hello_device returns correct errors to the device (#1706)
66a99695 chore(pairing): added correct error logs for prove_device (#1707)
0c651cf9 ci: fix e2e workflow to use Astarte images built in CI (#1718)
f755b582 build(docker): fix libs location
68e354f4 build(docker): optimize dockerignore
546eef43 chore: update mix.lock
20607bd5 build(docker): explicitly set hostnames
3fb57cfc build(docker): optimize Dockerfiles
232a4ca3 chore(fdo): ownership voucher validation (#1722)
ba1ddd7a fix(fdo): clean up old sessions (#1720)
5fd5637c ci(e2e): fix trigger pollution (#1719)
61ea879c fix(appengine): properly encode volatlie trigger requests (#1721)
df7650be feat(pairing): add support for SECP384R1 (ES384) device attestation
f4944514 fix: crashes on new scylladb versions (#1711)
226b4ad4 ci: avoid conflicting caches for libs/ (#1724)
d07bdc0f refactor(events): fix typespecs
23c9fdfb refactor(events): move more logic into core
fdde7aa7 refactor(events): align `fetch_device_event` and `fetch_data_event` impl
3d08747e feat: trigger installation notification (#1603)
cc124f6b chore: add typos spell checker (#1723)
3629b40b feat: trigger deletion notification (#1730)
71f65394 refactor: use phoenix pubsubs for trigger notifications (#1729)
b4b403c4 feat(pairing): handle chunked device service info (#1710)
ca551c83 feat(rm): trigger installation and deletion notification (#1731)
c63367a8 chore(rpc): do not start astarte_events
b97fa816 refactor: move clustering logic to astarte_rpc (#1732)
d7e18632 feat(pairing): handle trigger notifications (#1734)
315b1837 chore(rm): cluster realm management replicas
e3e8190e refactor: make realm management cluster with pairing (#1735)
34da830f perf: optimize trigger notifications
2f1acb90 chore(fdo): replacing stubs with real calls on controller_test
dbdd97f4 chore(fdo): add dedicated session module for tests
a6374bfb chore(fdo): entry_private_key directly returns the decoded owner public key
733873de build(docker): fix dev container build errors
642ec1f5 feat(pairing): implement OwnerServiceInfo chunking  (#1740)
735a48ce ci: use staple-actions (#1743)
508d26e2 fix(rpc): don't crash on rm+pairing clustering (#1746)

Let me find the 1.3.3 device-error fix commit specifically.
$ gh api 'repos/astarte-platform/astarte/commits?path=apps/astarte_data_updater_plant&since=2026-07-01&until=2026-08-21' --jq '.[] | "\(.sha[0:8]) \(.commit.message | split("\n")[0])"' 2>&1 | head -40
3484096d Prepare release v1.4.0-rc.5 (#2147)
73972499 chore: Update exandra
643d0610 ci: add release-please for automated versioning and changelog
fdc6d29b fix(dup): Make object values only validate type on selected interface (#2141)
6e65c8ac refactor(dup): Move encryption logic into library (#2123)
88004a2c chore(core): use the local astarte_core across the monorepo
0809af19 build: use elixir 1.20 and OTP 28 (#2071)
e26dfd03 chore: forward port release-1.4
aedc4148 fix(dup): ensure binaryblob data is correctly validated
a816eaf0 chore: forward port release-1.3
d8592091 fix(dup): stop remapping interface loading failed error (#2119)
171b83c4 chore: prepare v1.4.0-rc.4 (#2121)
208bfdb3 chore(dup): decrypt payloads incoming on encrypted endpoints (#2103)
85c05c67 refactor(generators)!: nuke to_changes + use adapters
b61a3ed3 chore(DUP): Refactor encrypted endpoints (#2112)
5a8e3591 fix: ensure RPC server reliability (#2115)
40e09e1f fix(dup): bound AMQP consumers memory growth (#2114)
bf9437ca chore: prepare v1.4.0-rc.3 (#2113)
4789d047 fix(dup): restore correct type validation checks for incoming data
5047681b feat(dup): Encrypt device shared secret (#2104)
1657e508 fix: properly cleanup rpc supervisor (#2109)
a4a841ed fix(dup): allow rpc server horizontal scalability (#2108)
6bc5e435 chore(dup): align mississippi to stable version (#2106)
646202eb fix(dup): stop re-queueing rejected messages (#2105)
0da8336f chore(dup): store shared secret during key agreement (#2101)
1c297aa3 refactor: purge exrabbitpool (#2050)
9424da20 chore(dup): set prod default log level to info (#2099)
49df96dc chore(dup): add rpc server startup log
ba2cd054 chore(dup): configure log level through env var at runtime
c8bbee87 chore(dup): add rpc server startup log (#2096)

$ gh api 'repos/astarte-platform/astarte/commits/d8592091' --jq '.files[] | "\(.filename) +\(.additions) -\(.deletions)"' 2>&1 | head -20
CHANGELOG.md +6 -0
apps/astarte_appengine_api/mix.lock +1 -1
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex +35 -2
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/device.ex +28 -11
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/control_handler_test.exs +29 -0
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/device_test.exs +14 -1
apps/astarte_pairing/mix.lock +1 -1

$ gh api 'repos/astarte-platform/astarte/commits/d8592091' --jq '.commit.message' 2>&1
fix(dup): stop remapping interface loading failed error (#2119)

Signed-off-by: Riccardo Nalgi <riccardo.nalgi@secomind.com>

$ gh api 'repos/astarte-platform/astarte/commits/d8592091' --jq '.files[] | select(.filename=="apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex") | .patch' 2>&1 | head -80
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

The 1.3.3 change distinguishes two device_error names (`interface_loading_failed` vs `resend_interface_properties_failed`) on the properties-resend path. Now let me check the 1.3.0 changes: empty introspection, capabilities, and `purge_properties_compression_format`.
$ gh api 'repos/astarte-platform/astarte/commits?path=apps/astarte_data_updater_plant&since=2026-01-01&until=2026-05-06&per_page=100' --jq '.[] | "\(.sha[0:8]) \(.commit.message | split("\n")[0])"' 2>&1 | head -60
26c9a5b9 chore: prepare v1.3.0
ff816f3f feat(docker): add health status (#1929)
5bf6243b chore: prepare v1.3.0-rc.2
c89cda80 chore: prepare v1.4.0-rc.0 (#1912)
958aa2be chore: forward-port release-1.3
d1c16ce7 chore(deps): update (#1861)
9320e9a5 refactor(fdo): add dedicated libraries (#1819)
d83bd921 chore: forward-port release-1.3
90a921e2 chore: forward-port release-1.2
cb66c7e3 chore: prepare v1.4.0-dev release cycle (#1853)
44dd26ea fix(dup): increase device process resiliency (#1852)
038ff72a feat(dup): Validate required interface mapping (#1849)
43b8d3b7 chore: prepare v1.2.1 (#1842)
e1b066ce chore: Add Required field to Endpoints table (#1846)
c2b53533 build(credo): update to 1.17 (#1841)
18000f06 fix(dup): ensure device deletion rpc always checks updated status (#1838)
8f1241fa chore: forward-port release-1.3 to master
6678c0e3 chore: merge branch 'feat/fido-device-onboard' into release-1.3
d50f1e38 build: update skogsra and lockfiles (#1832)
9f655414 chore: forward-port release-1.3 into master
71c35c48 chore: forward-port from 1.2 (solve issue #1072)
f06395d8 refactor(housekeeping)!: Consolidate AMQP management configuration (#1804)
07eecd09 fix(rm): detect conflicting options in intf aggregate mappings (#1605)
f4b99a23 chore(dup): Restore test functionality
f46529fa test(dup): Fix flakiness (#1789)
3fa2baeb refactor(dup): Fix credo complains (#1796)
5f397591 chore: prepare v1.2.1-rc.1
4347814c test(dup): update triggers in test
17465b3e refactor(events): fix credo complains (#1787)
dee633eb refactor(appengine): improve code quality and solve credo complains (#1770)
667da5c7 chore: Prepare v1.3.0-rc.1 (#1774)
18895720 test(events): Add tests to the Astarte Events
a2e0405c refactor: use correct typing for amqp static headers (#1751)
80d84394 ci: add credo (#1755)
735a48ce ci: use staple-actions (#1743)
733873de build(docker): fix dev container build errors

$ gh api 'repos/astarte-platform/astarte/commits?since=2025-10-01&until=2026-05-06&per_page=100' --jq '.[] | "\(.sha[0:8]) \(.commit.message | split("\n")[0])"' 2>&1 | rg -i 'capabilit|introspection|purge|empty'
(no output)

$ gh search code --owner astarte-platform 'purge_properties_compression_format' --json path,repository,url 2>&1 | head -40
[{"path":"CHANGELOG.md","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/CHANGELOG.md"},{"path":"libs/astarte_core/lib/astarte_core/device/capabilities.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/libs/astarte_core/lib/astarte_core/device/capabilities.ex"},{"path":"tools/astarte_import/lib/astarte/import.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/tools/astarte_import/lib/astarte/import.ex"},{"path":"doc/pages/architecture/090-database.md","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/doc/pages/architecture/090-database.md"},{"path":"doc/pages/architecture/080-mqtt-v1-protocol.md","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/doc/pages/architecture/080-mqtt-v1-protocol.md"},{"path":"lib/astarte_core/device/capabilities.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMTk3MDI0NTM=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte_core","url":"https://github.com/astarte-platform/astarte_core"},"url":"https://github.com/astarte-platform/astarte_core/blob/ceae2056ee37ed0b9413f74470f5e416db4f0cdf/lib/astarte_core/device/capabilities.ex"},{"path":"libs/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0008_create_capabilities_type.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/libs/astarte_data_access/lib/astarte_data_access/database/migrations/realm/0008_create_capabilities_type.ex"},{"path":"apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex"},{"path":"astarte/1.3/090-database.html","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMDU2MzM5NTA=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/docs","url":"https://github.com/astarte-platform/docs"},"url":"https://github.com/astarte-platform/docs/blob/4a34fe129ee3b7e1301bac0e6aa5109282eba77f/astarte/1.3/090-database.html"},{"path":"astarte/1.4/090-database.html","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMDU2MzM5NTA=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/docs","url":"https://github.com/astarte-platform/docs"},"url":"https://github.com/astarte-platform/docs/blob/4a34fe129ee3b7e1301bac0e6aa5109282eba77f/astarte/1.4/090-database.html"},{"path":"astarte/snapshot/090-database.html","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMDU2MzM5NTA=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/docs","url":"https://github.com/astarte-platform/docs"},"url":"https://github.com/astarte-platform/docs/blob/4a34fe129ee3b7e1301bac0e6aa5109282eba77f/astarte/snapshot/090-database.html"},{"path":"lib/astarte_device_sdk/device/session_manager.c","repository":{"id":"MDEwOlJlcG9zaXRvcnk3NDA1MjkzMTQ=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte-device-sdk-zephyr","url":"https://github.com/astarte-platform/astarte-device-sdk-zephyr"},"url":"https://github.com/astarte-platform/astarte-device-sdk-zephyr/blob/c2cc981bd79e3234336b7f6e573018fb64a60ceb/lib/astarte_device_sdk/device/session_manager.c"},{"path":"libs/astarte_core/test/astarte_core/device/capabilities_test.exs","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/libs/astarte_core/test/astarte_core/device/capabilities_test.exs"},{"path":"tools/astarte_export/test/astarte_export_test.exs","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/tools/astarte_export/test/astarte_export_test.exs"},{"path":"tools/astarte_export/test/database_test.exs","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/tools/astarte_export/test/database_test.exs"},{"path":"tools/astarte_import/test/astarte/import_test.exs","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/tools/astarte_import/test/astarte/import_test.exs"},{"path":"tools/astarte_import/test/astarte/populatedb_test.exs","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/tools/astarte_import/test/astarte/populatedb_test.exs"},{"path":"test/astarte_core/device/capabilities_test.exs","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMTk3MDI0NTM=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte_core","url":"https://github.com/astarte-platform/astarte_core"},"url":"https://github.com/astarte-platform/astarte_core/blob/ceae2056ee37ed0b9413f74470f5e416db4f0cdf/test/astarte_core/device/capabilities_test.exs"},{"path":"libs/astarte_rpc/test/support/helpers/database.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/libs/astarte_rpc/test/support/helpers/database.ex"},{"path":"libs/astarte_fdo/test/support/helpers/database.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/libs/astarte_fdo/test/support/helpers/database.ex"},{"path":"apps/astarte_trigger_engine/test/support/helpers/database.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/apps/astarte_trigger_engine/test/support/helpers/database.ex"},{"path":"tools/astarte_import/test/astarte/database_test.exs","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/tools/astarte_import/test/astarte/database_test.exs"},{"path":"apps/astarte_pairing/test/support/helpers/database.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/apps/astarte_pairing/test/support/helpers/database.ex"},{"path":"apps/astarte_appengine_api/test/support/helpers/database_v2.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/apps/astarte_appengine_api/test/support/helpers/database_v2.ex"},{"path":"apps/astarte_realm_management/test/support/helpers/database.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/apps/astarte_realm_management/test/support/helpers/database.ex"},{"path":"libs/astarte_data_access/test/support/helpers/database.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/libs/astarte_data_access/test/support/helpers/database.ex"},{"path":"apps/astarte_appengine_api/test/support/helpers/database.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/apps/astarte_appengine_api/test/support/helpers/database.ex"},{"path":"libs/astarte_events/test/support/helpers/database_test_helper.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/libs/astarte_events/test/support/helpers/database_test_helper.ex"},{"path":"apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/capabilities_handler_test.exs","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/capabilities_handler_test.exs"},{"path":"apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater_test.exs","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"},"url":"https://github.com/astarte-platform/astarte/blob/ad67a092755dfb7de4ec0be1c1e87d35f7e86993/apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater_test.exs"}]

→ Read .mule/todo.md
→ Read .mule/for-giulio.md
Let me verify there's no newer release/tag than what I found, to confirm whether upstream has moved at all.

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:8][] | "\(.tag_name)\t\(.published_at)\tprerelease=\(.prerelease)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z	prerelease=true
v1.3.3	2026-08-07T14:33:58Z	prerelease=false
v1.4.0-rc.4	2026-08-06T13:40:20Z	prerelease=true
v1.4.0-rc.3	2026-07-31T10:48:39Z	prerelease=true
v1.4.0-rc.2	2026-07-14T10:16:22Z	prerelease=true
v1.4.0-rc.1	2026-07-14T08:53:27Z	prerelease=true
v1.3.2	2026-07-14T10:07:26Z	prerelease=false
v1.3.1	2026-07-13T11:50:05Z	prerelease=false

$ gh api 'repos/astarte-platform/astarte/tags?per_page=10' --jq '.[] | "\(.name)\t\(.commit.sha[0:8])"'
v1.4.0-rc.5	3484096d
v1.4.0-rc.4	171b83c4
v1.4.0-rc.3	bf9437ca
v1.4.0-rc.2	08d7d5e6
v1.4.0-rc.1	c5fd781a
v1.4.0-rc.0	c89cda80
v1.3.3	3e70d63f
v1.3.2	41b4f4e8
v1.3.1	17b329a9
v1.3.0	26c9a5b9

No upstream movement since the last run of this recipe. The release sweep shows the same state the previous run recorded: newest release is still `v1.4.0-rc.5` (2026-08-20, prerelease), newest stable `v1.3.3` (2026-08-07) — nothing released since, and today's milestone run (2026-09-12) confirmed the same. Any wire-relevant proposal those releases would generate is already sitting in the queue (`probe-props-resend-error-triggers`, `probe-property-resend-encoding`, `empty-introspection-verification`, `probe-trigger-install-notification-delay`, `purge-properties-compression-capability`, `compat-note-v1.3.3`/`v1.4-rc`), and `for-giulio.md` already carries the compat-note wording proposal. Nothing new to propose.

```
Done: upstream unchanged — no releases since v1.4.0-rc.5 (2026-08-20) / v1.3.3 (2026-08-07); all wire-relevant items already proposed in the queue
Files: none
Verified: gh api repos/astarte-platform/astarte/releases -> all ≤ 2026-08-20; tags list -> same; diff v1.2.2...v1.3.3 -> nothing new
Unsure: nothing
Follow-ups: none
```
