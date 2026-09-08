slug: recipe-astarte-upstream
verdict: proposed
at:  1e56fda
ran: 2026-09-08T11:34:22Z on DietPi in 305s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md
Now let me read the compatibility doc and check upstream releases.

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

Upstream is ahead: latest stable is **v1.3.3** and latest overall is **v1.4.0-rc.5**; Astrate targets 1.2.2. Let me get the release notes for the relevant versions.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
Sixth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.4

### Fixed

- [astarte_data_updater_plant] Make object values only validate type on selected interface
- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.2") | .body' 2>/dev/null | head -200
Astarte v1.3.2 release.

## Changes since v1.3.1

## Fixed

- Avoid crash on invalid properties message

v1.3.3 has empty body. Let me look at more releases to understand what changed since 1.2.2.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | "\(.tag_name)\t\(.published_at)"'
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
v1.2.0-rc.0	2024-06-11T13:02:33Z
v1.0.6	2024-04-23T21:15:57Z
v1.1.1	2023-11-15T15:29:57Z
v1.0.5	2023-09-26T10:27:57Z
v1.1.0	2023-06-20T18:59:09Z
v1.1.0-rc.0	2023-06-09T15:49:44Z
v1.1.0-alpha.0	2022-11-24T16:38:37Z
v1.0.4	2022-10-25T16:50:21Z
v1.0.3	2022-07-05T13:57:17Z
v1.0.2	2022-04-01T17:18:15Z

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

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.1") | .body'
Astarte v1.3.1 release.

## Changes since v1.3.0

### Fixed

- Solve capabilities message stalling queues
- [astarte_housekeeping] Allow resuming partially completed realm deletions
- Use network topology as default realm replication
- Fix simple strategy realm creation without replication factor

### Changed

- Immediately delete devices in device deletion best case scenario

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

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.0") | .body'
First Astarte v1.4.0 release candidate.

## Changes since v1.3

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

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.2") | .body'
Third Astarte v1.4.0 release candidate.

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

Now let me read the current todo.md to know what's already queued.
→ Read .mule/todo.md
Let me investigate a few more details before writing the proposals.

$ gh search code --owner astarte-platform "device_registration" --limit 5 --json repository,path 2>/dev/null | head -40
[{"path":"libs/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0005_add_device_registration_limit.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"}},{"path":"lib/astarte_rpc/protocol/proto/realm_management/get_device_registration_limit_reply.pb.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMTk3MTE4MzA=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte_rpc","url":"https://github.com/astarte-platform/astarte_rpc"}},{"path":"lib/astarte_rpc/protocol/proto/realm_management/get_device_registration_limit_reply.proto","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMTk3MTE4MzA=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte_rpc","url":"https://github.com/astarte-platform/astarte_rpc"}},{"path":"doc/pages/administrator/070-manage_realms.md","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"}},{"path":"doc/pages/user/020-accessing_and_exploring_a_realm.md","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"}}]

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .target_commitish'
master

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .html_url'
https://github.com/astarte-platform/astarte/releases/tag/v1.3.3

$ gh api repos/astarte-platform/astarte/compare/v1.3.2...v1.3.3 --jq '.commits[] | "\(.sha[0:8]) \(.commit.message | split("\n")[0])"'
d8592091 fix(dup): stop remapping interface loading failed error (#2119)
3e70d63f chore: prepare v1.3.3 (#2128)

$ gh search code --owner astarte-platform "device_registration_trigger" --limit 10 --json repository,path 2>/dev/null | head -60
[{"path":"apps/astarte_pairing/lib/astarte_pairing/agent/agent.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"}},{"path":"apps/astarte_pairing/test/support/helpers/triggers.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"}},{"path":"apps/astarte_pairing/test/astarte_pairing/agent/agent_test.exs","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"}}]

$ gh api repos/astarte-platform/astarte/contents/apps/astarte_pairing/lib/astarte_pairing/agent/agent.ex --jq '.content' | base64 -d 2>/dev/null | head -100
#
# This file is part of Astarte.
#
# Copyright 2017 - 2025 SECO Mind Srl
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

defmodule Astarte.Pairing.Agent do
  @moduledoc """
  The Agent context.
  """

  alias Astarte.Core.Device
  alias Astarte.Pairing.Agent.DeviceRegistrationRequest
  alias Astarte.Pairing.Agent.DeviceRegistrationResponse
  alias Astarte.Pairing.Engine

  alias Astarte.Core.Triggers.SimpleEvents.DeviceRegisteredEvent
  alias Astarte.Events.Triggers
  alias Astarte.Events.TriggersHandler

  def register_device(realm, attrs \\ %{}) do
    changeset =
      %DeviceRegistrationRequest{}
      |> DeviceRegistrationRequest.changeset(attrs)

    with {:ok,
          %DeviceRegistrationRequest{hw_id: hw_id, initial_introspection: initial_introspection}} <-
           Ecto.Changeset.apply_action(changeset, :insert),
         initial_introspection =
           Enum.map(initial_introspection, fn {interface_name,
                                               %{"major" => major, "minor" => minor}} ->
             %{
               interface_name: interface_name,
               major_version: major,
               minor_version: minor
             }
           end),
         {:ok, credentials_secret} <-
           Engine.register_device(realm, hw_id, initial_introspection: initial_introspection) do
      dispatch_device_registration_trigger(realm, hw_id)

      {:ok, %DeviceRegistrationResponse{credentials_secret: credentials_secret}}
    end
  end

  def unregister_device(realm, device_id) do
    with {:ok, _} <- Device.decode_device_id(device_id) do
      Engine.unregister_device(realm, device_id)
    end
  end

  defp dispatch_device_registration_trigger(realm_name, hw_id) do
    timestamp = DateTime.utc_now() |> DateTime.to_unix(:millisecond)
    event_key = :on_device_registered
    event_type = :device_registered_event
    {:ok, device_id} = Device.decode_device_id(hw_id, allow_extended_id: true)

    Triggers.find_device_trigger_targets(realm_name, device_id, event_key)
    |> dispatch_all(realm_name, hw_id, timestamp, event_type, %DeviceRegisteredEvent{})
  end

  defp dispatch_all(targets, realm_name, device_id, timestamp, event_type, event) do
    targets
    |> Enum.map(fn {target, policy} ->
      TriggersHandler.dispatch_event(
        event,
        event_type,
        target,
        realm_name,
        device_id,
        timestamp,
        policy
      )
    end)
  end
end

$ gh api repos/astarte-platform/astarte/compare/v1.2.2...v1.3.0 --jq '[.commits[] | .commit.message | split("\n")[0]] | .[:30] | .[]'
fix(pairing) fix setup device payload (#1622)
chore(pairing): ov next entry encoding (#1625)
refactor(pairing): rework prove ovh hdr
chore(pairing): eb_sig_info for ecc (#1632)
chore: init eatoken module (#1633)
fix(pairing): correctly build rv_to2_addr entries (#1634)
fix(pairing): owner_sign message (#1635)
test: solve trigger cache flakiness
chore: forward-port release-1.3 into feat/fido-device-onboard
feat(pairing): implement TO2.ProveDevice (Type 64) decoding
fix(pairing): add appropriate variables for rv entry (#1641)
refactor(pairing): init fdo view (#1642)
chore: forward-port release-1.3
fix(rm): deletion for long disconnected devices
chore(pairing) handle messages 66 and 67 (#1628)
chore(pairing): device -> owner secure tunnel (#1643)
chore: owner -> device secure tunnel (#1644)
fix(pairing): cipher suite name parsing (#1645)
fix(pairing): correctly parse signature info (#1646)
docs(trigger): document missing device trigger events (#1653)
fix(pairing): correctly parse hello device binary fields (#1652)
refactor(pairing): add fallback controller for fdo (#1655)
fix(pairing): align query return format
fix(pairing): send all required headers during fdo
fix(pairing): fix matching owner public key
feat(pairing): TO2 Done/Done2 msgs
feat(pairing): save SetupDvNonce in TO2 session
fix(pairing): put euphnonce in unprotected headers
build(docker): add the FDO traefik path labels (#1666)
fix(pairing): signing of ProveOVHdr

✱ Grep "device_registered|device_deletion_started|device_deletion_completed" in internal · 20 matches

$ gh api repos/astarte-platform/astarte/commits/d8592091 --jq '.commit.message'
fix(dup): stop remapping interface loading failed error (#2119)

Signed-off-by: Riccardo Nalgi <riccardo.nalgi@secomind.com>

✱ Grep "device_deletion_finished|device_deletion_completed" in internal · 10 matches

$ gh api repos/astarte-platform/astarte/commits?path=apps/astarte_data_updater_plant&per_page=5 --jq '.[].commit.message' | head -20
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_080c9b72b001q93aWYE12DNzjr

vQ\nhWJYqgCTlXiMBpsLKkQnU1Tgp3nvAMC+MywAzuGThBZxOd/eFkZBsW8z9hp2FuhT\noaBYpmlsVK06BNIsOVa1gSe+luat784VGVRlaxH/YZbqVpMYyXQtCsFA+DjisnE/\nWnYt+qhJ2ptVseIgLJZJowd+PCH798IqBDDjIoQspMH6mruolS1U9ZVWHK8deuCz\ngaGmKxe4wbSXcpxPb3VioqKVYpTz9R5QQsTDoi9hQokZ2rT+ALqgsMk7shU8j5PX\nOeF7utM01q6LtgF4i9C+ZeDdzqQI2AkDkccXX+bz9z/vVShZSfDRs8MAU0DDX13G\nOzOLmvqsB+HrMh7UOELkdV8jGW49dH5VbZV+ziVp5232DJ2VSSVijWx6UgUJ959o\nvPdHU7idHscKkB90rMkg4zarJDppxYbm8OV7x9xM50d/l1H+yk9Y8w0Mhc+hCnkA\nmdCfuVMzWwaJJLQJRkrIAy+1gO5KrXnK0ss9okUVgxELcz++YBIKNMJpfFqnQnfb\n/4zfc9RnFeK6M4XzUXJI\n=jnEW\n-----END PGP SIGNATURE-----\n","payload":"tree 2fabf8ea5e2f445fec42653e0fe9f4cdcab493e2\nparent 41b4f4e8923535a1a575ffb6af42c605f2d1f4fb\nauthor Riccardo Nalgi <64850100+mizzet1@users.noreply.github.com> 1786085711 +0200\ncommitter GitHub <noreply@github.com> 1786085711 +0200\n\nfix(dup): stop remapping interface loading failed error (#2119)\n\nSigned-off-by: Riccardo Nalgi <riccardo.nalgi@secomind.com>","verified_at":"2026-08-07T06:55:15Z"}},"url":"https://api.github.com/repos/astarte-platform/astarte/commits/d8592091e1a48d5a50dee923e3a7d2cc487e55b2","html_url":"https://github.com/astarte-platform/astarte/commit/d8592091e1a48d5a50dee923e3a7d2cc487e55b2","comments_url":"https://api.github.com/repos/astarte-platform/astarte/commits/d8592091e1a48d5a50dee923e3a7d2cc487e55b2/comments","author":{"login":"mizzet1","id":64850100,"node_id":"MDQ6VXNlcjY0ODUwMTAw","avatar_url":"https://avatars.githubusercontent.com/u/64850100?v=4","gravatar_id":"","url":"https://api.github.com/users/mizzet1","html_url":"https://github.com/mizzet1","followers_url":"https://api.github.com/users/mizzet1/followers","following_url":"https://api.github.com/users/mizzet1/following{/other_user}","gists_url":"https://api.github.com/users/mizzet1/gists{/gist_id}","starred_url":"https://api.github.com/users/mizzet1/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/mizzet1/subscriptions","organizations_url":"https://api.github.com/users/mizzet1/orgs","repos_url":"https://api.github.com/users/mizzet1/repos","events_url":"https://api.github.com/users/mizzet1/events{/privacy}","received_events_url":"https://api.github.com/users/mizzet1/received_events","type":"User","user_view_type":"public","site_admin":false},"committer":{"login":"web-flow","id":19864447,"node_id":"MDQ6VXNlcjE5ODY0NDQ3","avatar_url":"https://avatars.githubusercontent.com/u/19864447?v=4","gravatar_id":"","url":"https://api.github.com/users/web-flow","html_url":"https://github.com/web-flow","followers_url":"https://api.github.com/users/web-flow/followers","following_url":"https://api.github.com/users/web-flow/following{/other_user}","gists_url":"https://api.github.com/users/web-flow/gists{/gist_id}","starred_url":"https://api.github.com/users/web-flow/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/web-flow/subscriptions","organizations_url":"https://api.github.com/users/web-flow/orgs","repos_url":"https://api.github.com/users/web-flow/repos","events_url":"https://api.github.com/users/web-flow/events{/privacy}","received_events_url":"https://api.github.com/users/web-flow/received_events","type":"User","user_view_type":"public","site_admin":false},"parents":[{"sha":"41b4f4e8923535a1a575ffb6af42c605f2d1f4fb","url":"https://api.github.com/repos/astarte-platform/astarte/commits/41b4f4e8923535a1a575ffb6af42c605f2d1f4fb","html_url":"https://github.com/astarte-platform/astarte/commit/41b4f4e8923535a1a575ffb6af42c605f2d1f4fb"}]},{"sha":"171b83c481413c74f693728cce6ee48e4c16aa6f","node_id":"C_kwDOB5ht_doAKDE3MWI4M2M0ODE0MTNjNzRmNjkzNzI4Y2NlNmVlNDhlNGMxNmFhNmY","commit":{"author":{"name":"Francesco Noacco","email":"francesco.noacco@secomind.com","date":"2026-08-06T13:38:40Z"},"committer":{"name":"GitHub","email":"noreply@github.com","date":"2026-08-06T13:38:40Z"},"message":"chore: prepare v1.4.0-rc.4 (#2121)\n\nSigned-off-by: Francesco Noacco <francesco.noacco@secomind.com>","tree":{"sha":"373e3fbb73d4c6f38dd83e41e0ee7d6f5dbd1420","url":"https://api.github.com/repos/astarte-platform/astarte/git/trees/373e3fbb73d4c6f38dd83e41e0ee7d6f5dbd1420"},"url":"https://api.github.com/repos/astarte-platform/astarte/git/commits/171b83c481413c74f693728cce6ee48e4c16aa6f","comment_count":0,"verification":{"verified":true,"reason":"valid","signature":"-----BEGIN PGP SIGNATURE-----\n\nwsFcBAABCAAQBQJqdI5gCRC1aQ7uu5UhlAAAMPoQACuEHWHom8RB/w5cwekWCxMU\nAmtxp9cEEhcrN8+ttTzL280SlhWzS1Ukp0ouqIkyZpH/J075cnyzp3q2M/pbW3cy\nrwghtVIaGAWIzWvgweve4TRwgOYshnMW4zMw6HF/Wb+lWWXCz5rjT5DUFrfvTgqy\n0jrD0NrRhIFbqAWLi9iw6e0ebrAImGK2iJnWr8QfdORJ2XH1iCUU9ykH0PemZX4Y\nY/NM2Co1ZO+TWvLjhark3Wbr9Hg5a+s1Qk+qwS9gMIeg1iHj9vsHd+SEDmNSAoeQ\nnEYuP9zng9SGhd9/B6h+hA4wnP6xhoZFUiQBRzmbRYE7rlEK/tx5JMji4A7TKOCL\njhPZ//m8eTCZPJ2uMz0srKVdD4jmgh2OQvLdDEz7farCF3XPClbZZX1iOKmIFKPs\noq5oLJND3Rr75DhCuB1TE8ar8K6idADnc+ow2INWQxVcOnoZGWQa+PxWiA0IGJcB\nb95K3w6ILyMXiA7FfSIzksBvPWk9UTW0O20040EPdwpEfBmnYyY9L23MpKid4BU/\n4owA+OnGMMcgZuNYQCB99SySqjcCrJrx6BHl40KIFwmUWvnmIH/+ygRmfMyI046o\nzR1FmfEFvNEQiOV/tybCU7Mjuh5ngoiwpueJ7FgpHQZhlfyJF5JsU9CfGeR262KQ\nt0oFZKNrKyBoMuUvmrT9\n=d0L1\n-----END PGP SIGNATURE-----\n","payload":"tree 373e3fbb73d4c6f38dd83e41e0ee7d6f5dbd1420\nparent 625e1f6780c982ca19318480d9af5b20eda34a15\nauthor Francesco Noacco <francesco.noacco@secomind.com> 1786023520 +0200\ncommitter GitHub <noreply@github.com> 1786023520 +0200\n\nchore: prepare v1.4.0-rc.4 (#2121)\n\nSigned-off-by: Francesco Noacco <francesco.noacco@secomind.com>","verified_at":"2026-08-06T13:38:44Z"}},"url":"https://api.github.com/repos/astarte-platform/astarte/commits/171b83c481413c74f693728cce6ee48e4c16aa6f","html_url":"https://github.com/astarte-platform/astarte/commit/171b83c481413c74f693728cce6ee48e4c16aa6f","comments_url":"https://api.github.com/repos/astarte-platform/astarte/commits/171b83c481413c74f693728cce6ee48e4c16aa6f/comments","author":{"login":"noaccOS","id":24324352,"node_id":"MDQ6VXNlcjI0MzI0MzUy","avatar_url":"https://avatars.githubusercontent.com/u/24324352?v=4","gravatar_id":"","url":"https://api.github.com/users/noaccOS","html_url":"https://github.com/noaccOS","followers_url":"https://api.github.com/users/noaccOS/followers","following_url":"https://api.github.com/users/noaccOS/following{/other_user}","gists_url":"https://api.github.com/users/noaccOS/gists{/gist_id}","starred_url":"https://api.github.com/users/noaccOS/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/noaccOS/subscriptions","organizations_url":"https://api.github.com/users/noaccOS/orgs","repos_url":"https://api.github.com/users/noaccOS/repos","events_url":"https://api.github.com/users/noaccOS/events{/privacy}","received_events_url":"https://api.github.com/users/noaccOS/received_events","type":"User","user_view_type":"public","site_admin":false},"committer":{"login":"web-flow","id":19864447,"node_id":"MDQ6VXNlcjE5ODY0NDQ3","avatar_url":"https://avatars.githubusercontent.com/u/19864447?v=4","gravatar_id":"","url":"https://api.github.com/users/web-flow","html_url":"https://github.com/web-flow","followers_url":"https://api.github.com/users/web-flow/followers","following_url":"https://api.github.com/users/web-flow/following{/other_user}","gists_url":"https://api.github.com/users/web-flow/gists{/gist_id}","starred_url":"https://api.github.com/users/web-flow/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/web-flow/subscriptions","organizations_url":"https://api.github.com/users/web-flow/orgs","repos_url":"https://api.github.com/users/web-flow/repos","events_url":"https://api.github.com/users/web-flow/events{/privacy}","received_events_url":"https://api.github.com/users/web-flow/received_events","type":"User","user_view_type":"public","site_admin":false},"parents":[{"sha":"625e1f6780c982ca19318480d9af5b20eda34a15","url":"https://api.github.com/repos/astarte-platform/astarte/commits/625e1f6780c982ca19318480d9af5b20eda34a15","html_url":"https://github.com/astarte-platform/astarte/commit/625e1f6780c982ca19318480d9af5b20eda34a15"}]},{"sha":"208bfdb360d92999711c79f807de3599cac4c15e","node_id":"C_kwDOB5ht_doAKDIwOGJmZGIzNjBkOTI5OTk3MTFjNzlmODA3ZGUzNTk5Y2FjNGMxNWU","commit":{"author":{"name":"Francesco Asq","email":"152388339+frossq@users.noreply.github.com","date":"2026-08-06T09:11:43Z"},"committer":{"name":"GitHub","email":"noreply@github.com","date":"2026-08-06T09:11:43Z"},"message":"chore(dup): decrypt payloads incoming on encrypted endpoints (#2103)\n\nDecrypt CBOR-encrypted payloads sent from devices on encrypted endpoints,\nusing the shared secret previously derived and found in DUP state.\n\nEnsure that the key status gets set to 'established' when the shared key\nis correctly retrieved from db.\n\nAdded tests for payloads decryption and related errors.\nTest setup now also imports shared secret into DUP state.\n\nSigned-off-by: frossq <francesco.asquini@secomind.com>","tree":{"sha":"574eabb1537301ab0f7f5109281aae617d7d7833","url":"https://api.github.com/repos/astarte-platform/astarte/git/trees/574eabb1537301ab0f7f5109281aae617d7d7833"},"url":"https://api.github.com/repos/astarte-platform/astarte/git/commits/208bfdb360d92999711c79f807de3599cac4c15e","comment_count":0,"verification":{"verified":true,"reason":"valid","signature":"-----BEGIN PGP SIGNATURE-----\n\nwsFcBAABCAAQBQJqdE/PCRC1aQ7uu5UhlAAAybwQAIIEtH9PCzoVuBsBFPtOrcsF\n/gb9XJKy3mizpnFfMwZCfPTiukwmrP7rxujAQsmFRhgJW2ecEeo5GqZFTy07rMXO\nKKE9sMLP5AauxM3Q9QKXKKrAEnXUVPmQLMqAUCNsvgRc+jC09P2wK4LpJgw5KqVw\npDLnn4u+he7oSvGm0g/+lJxFh+SsX7Lza/iiWqW7zj/A5yitQJ0Sk+Zc3vhp7QGw\n1RmqZIYucGIztU38c5nc4r6teYauR39RAh1pFLGjGzBHLG7pj5rdCbWd9WcGcmMw\n6uL35/92b9OhLrmKtNCbFYOIvDZ2IQRYFIEREImmx/e6CA6ODH7Dv3GYDv7ALgS2\naRy4sp86Sj9wwOdzo6jFLZAylnWF/MIzm/ikWu07OBX7PUNYoRz8rl+zU36odW10\nZpC04SHVEw6VBcTjGixK8zNRSOvbm2ox6dmp4fXNEhTNoJM+HhdNU24+Nds3Fojm\nARNI3EfZRpTnHt24FFwUr7OO5pjD08T3MNsjGFGbiNkObVWbIb1xBdB+H+EnwjEz\nEQhr8mCca+Fr7j31h9zB8bo2VkcRcDBxj1Rjo+RKvi4Za5f62VhUtOTVEnHiII5h\nMQ+k/4ChOIYuDgLV35bxNmYWJuJAYjCIlJbEeiPC3Ak01Jlg8hhEqQYjUi/yZ8nE\ntSURAXQyV3mY07LNffOm\n=D2pN\n-----END PGP SIGNATURE-----\n","payload":"tree 574eabb1537301ab0f7f5109281aae617d7d7833\nparent 85c05c67ccbf971c90c6c73a5f35facdb27cdf6f\nauthor Francesco Asq <152388339+frossq@users.noreply.github.com> 1786007503 +0200\ncommitter GitHub <noreply@github.com> 1786007503 +0200\n\nchore(dup): decrypt payloads incoming on encrypted endpoints (#2103)\n\nDecrypt CBOR-encrypted payloads sent from devices on encrypted endpoints,\nusing the shared secret previously derived and found in DUP state.\n\nEnsure that the key status gets set to 'established' when the shared key\nis correctly retrieved from db.\n\nAdded tests for payloads decryption and related errors.\nTest setup now also imports shared secret into DUP state.\n\nSigned-off-by: frossq <francesco.asquini@secomind.com>","verified_at":"2026-08-06T09:11:48Z"}},"url":"https://api.github.com/repos/astarte-platform/astarte/commits/208bfdb360d92999711c79f807de3599cac4c15e","html_url":"https://github.com/astarte-platform/astarte/commit/208bfdb360d92999711c79f807de3599cac4c15e","comments_url":"https://api.github.com/repos/astarte-platform/astarte/commits/208bfdb360d92999711c79f807de3599cac4c15e/comments","author":{"login":"frossq","id":152388339,"node_id":"U_kgDOCRVC8w","avatar_url":"https://avatars.githubusercontent.com/u/152388339?v=4","gravatar_id":"","url":"https://api.github.com/users/frossq","html_url":"https://github.com/frossq","followers_url":"https://api.github.com/users/frossq/followers","following_url":"https://api.github.com/users/frossq/following{/other_user}","gists_url":"https://api.github.com/users/frossq/gists{/gist_id}","starred_url":"https://api.github.com/users/frossq/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/frossq/subscriptions","organizations_url":"https://api.github.com/users/frossq/orgs","repos_url":"https://api.github.com/users/frossq/repos","events_url":"https://api.github.com/users/frossq/events{/privacy}","received_events_url":"https://api.github.com/users/frossq/received_events","type":"User","user_view_type":"public","site_admin":false},"committer":{"login":"web-flow","id":19864447,"node_id":"MDQ6VXNlcjE5ODY0NDQ3","avatar_url":"https://avatars.githubusercontent.com/u/19864447?v=4","gravatar_id":"","url":"https://api.github.com/users/web-flow","html_url":"https://github.com/web-flow","followers_url":"https://api.github.com/users/web-flow/followers","following_url":"https://api.github.com/users/web-flow/following{/other_user}","gists_url":"https://api.github.com/users/web-flow/gists{/gist_id}","starred_url":"https://api.github.com/users/web-flow/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/web-flow/subscriptions","organizations_url":"https://api.github.com/users/web-flow/orgs","repos_url":"https://api.github.com/users/web-flow/repos","events_url":"https://api.github.com/users/web-flow/events{/privacy}","received_events_url":"https://api.github.com/users/web-flow/received_events","type":"User","user_view_type":"public","site_admin":false},"parents":[{"sha":"85c05c67ccbf971c90c6c73a5f35facdb27cdf6f","url":"https://api.github.com/repos/astarte-platform/astarte/commits/85c05c67ccbf971c90c6c73a5f35facdb27cdf6f","html_url":"https://github.com/astarte-platform/astarte/commit/85c05c67ccbf971c90c6c73a5f35facdb27cdf6f"}]},{"sha":"85c05c67ccbf971c90c6c73a5f35facdb27cdf6f","node_id":"C_kwDOB5ht_doAKDg1YzA1YzY3Y2NiZjk3MWM5MGM2YzczYTVmMzVmYWNkYjI3Y2RmNmY","commit":{"author":{"name":"Gabriele Ghio","email":"gabriele.ghio@secomind.com","date":"2026-06-12T11:30:16Z"},"committer":{"name":"Gabriele Ghio","email":"gabriele.ghio@secomind.com","date":"2026-08-05T14:34:34Z"},"message":"refactor(generators)!: nuke to_changes + use adapters\n\nThis is a major refactoring of the `astarte_generators` project,\ndesigned to facilitate the use of `astarte_adapters` and to simplify\nthe use of generators\n\n- the names of the generator functions have been standardised for use with\n`import`, thereby avoiding conflicts\n- the `to_changes` functions have been removed\n- the libraries have been updated\n- the applications have been updated\n\nSigned-off-by: Gabriele Ghio <gabriele.ghio@secomind.com>","tree":{"sha":"e2406a7aa1d8a6755e4d34ed41727c106de17988","url":"https://api.github.com/repos/astarte-platform/astarte/git/trees/e2406a7aa1d8a6755e4d34ed41727c106de17988"},"url":"https://api.github.com/repos/astarte-platform/astarte/git/commits/85c05c67ccbf971c90c6c73a5f35facdb27cdf6f","comment_count":0,"verification":{"verified":true,"reason":"valid","signature":"-----BEGIN PGP SIGNATURE-----\n\niQGzBAABCAAdFiEEJrs2ID58lAdyGg5Ut8sDhP9VdbkFAmpzSfsACgkQt8sDhP9V\ndblCBAv/XAWDJ+74a7rGeUw4yLgxJPicjm/t60d43g4KLT3QLhSn3lwaVtI6NDSE\nWGQMa5lSEwjmRf1eAJjCNNiepHai76PTCrujS0BabsKWbepcjymTjPGkD/96y+W+\nPX6suHJKetJX+cXCUi+ipLqapLhVj2gWE/JHUf4F3SCCiXSMUo1rLnqoafSYNiW+\nSO8mDyzX84GoFkwcGNTGP7k5Y592v7PmbU3EjX/neNznmkgM3PqOmQqjQfPP3tcX\nLJC9HFohDGe2w5t6ZyHBn8NvQu7JQR40MytOITBfDkQLcu6RVyx/exyLRajtBMhT\n9pxFq5NikAMPhhnas4C9r+OjeovQ8T1Sj9S59Q6vpSTFbDcmcUknkK8dkOST51X7\nAyrmos5k7b98uPDaS2vjf+LqWvuizS53tcGo2z4HJv5ksg/+VNKwYcSWcSU9D2EV\nHH6lFAHSRarQVb3VMwt5fbHqC/FWwK7/2ShczKfOxAz/Gf2n7Y0k4zxLDMkxcqDa\nVLQjizBo\n=PBaj\n-----END PGP SIGNATURE-----","payload":"tree e2406a7aa1d8a6755e4d34ed41727c106de17988\nparent b61a3ed33dd2cf8eb7f55c55446e01d3d6e2a9e9\nauthor Gabriele Ghio <gabriele.ghio@secomind.com> 1781263816 +0200\ncommitter Gabriele Ghio <gabriele.ghio@secomind.com> 1785940474 +0200\n\nrefactor(generators)!: nuke to_changes + use adapters\n\nThis is a major refactoring of the `astarte_generators` project,\ndesigned to facilitate the use of `astarte_adapters` and to simplify\nthe use of generators\n\n- the names of the generator functions have been standardised for use with\n`import`, thereby avoiding conflicts\n- the `to_changes` functions have been removed\n- the libraries have been updated\n- the applications have been updated\n\nSigned-off-by: Gabriele Ghio <gabriele.ghio@secomind.com>\n","verified_at":"2026-08-05T14:34:45Z"}},"url":"https://api.github.com/repos/astarte-platform/astarte/commits/85c05c67ccbf971c90c6c73a5f35facdb27cdf6f","html_url":"https://github.com/astarte-platform/astarte/commit/85c05c67ccbf971c90c6c73a5f35facdb27cdf6f","comments_url":"https://api.github.com/repos/astarte-platform/astarte/commits/85c05c67ccbf971c90c6c73a5f35facdb27cdf6f/comments","author":{"login":"shinnokdisengir","id":11438857,"node_id":"MDQ6VXNlcjExNDM4ODU3","avatar_url":"https://avatars.githubusercontent.com/u/11438857?v=4","gravatar_id":"","url":"https://api.github.com/users/shinnokdisengir","html_url":"https://github.com/shinnokdisengir","followers_url":"https://api.github.com/users/shinnokdisengir/followers","following_url":"https://api.github.com/users/shinnokdisengir/following{/other_user}","gists_url":"https://api.github.com/users/shinnokdisengir/gists{/gist_id}","starred_url":"https://api.github.com/users/shinnokdisengir/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/shinnokdisengir/subscriptions","organizations_url":"https://api.github.com/users/shinnokdisengir/orgs","repos_url":"https://api.github.com/users/shinnokdisengir/repos","events_url":"https://api.github.com/users/shinnokdisengir/events{/privacy}","received_events_url":"https://api.github.com/users/shinnokdisengir/received_events","type":"User","user_view_type":"public","site_admin":false},"committer":{"login":"shinnokdisengir","id":11438857,"node_id":"MDQ6VXNlcjExNDM4ODU3","avatar_url":"https://avatars.githubusercontent.com/u/11438857?v=4","gravatar_id":"","url":"https://api.github.com/users/shinnokdisengir","html_url":"https://github.com/shinnokdisengir","followers_url":"https://api.github.com/users/shinnokdisengir/followers","following_url":"https://api.github.com/users/shinnokdisengir/following{/other_user}","gists_url":"https://api.github.com/users/shinnokdisengir/gists{/gist_id}","starred_url":"https://api.github.com/users/shinnokdisengir/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/shinnokdisengir/subscriptions","organizations_url":"https://api.github.com/users/shinnokdisengir/orgs","repos_url":"https://api.github.com/users/shinnokdisengir/repos","events_url":"https://api.github.com/users/shinnokdisengir/events{/privacy}","received_events_url":"https://api.github.com/users/shinnokdisengir/received_events","type":"User","user_view_type":"public","site_admin":false},"parents":[{"sha":"b61a3ed33dd2cf8eb7f55c55446e01d3d6e2a9e9","url":"https://api.github.com/repos/astarte-platform/astarte/commits/b61a3ed33dd2cf8eb7f55c55446e01d3d6e2a9e9","html_url":"https://github.com/astarte-platform/astarte/commit/b61a3ed33dd2cf8eb7f55c55446e01d3d6e2a9e9"}]},{"sha":"b61a3ed33dd2cf8eb7f55c55446e01d3d6e2a9e9","node_id":"C_kwDOB5ht_doAKGI2MWEzZWQzM2RkMmNmOGViN2Y1NWM1NTQ0NmUwMWQzZDZlMmE5ZTk","commit":{"author":{"name":"nedimtokic","email":"64848448+nedimtokic@users.noreply.github.com","date":"2026-08-05T13:20:08Z"},"committer":{"name":"GitHub","email":"noreply@github.com","date":"2026-08-05T13:20:08Z"},"message":"chore(DUP): Refactor encrypted endpoints (#2112)\n\nUpdate encrypted interfaces to support\nparameters\n\nSigned-off-by: nedimtokic <nedim.tokic@secomind.com>","tree":{"sha":"eaacce67f898a97650a613be6ea3b7bdcd5eb9f3","url":"https://api.github.com/repos/astarte-platform/astarte/git/trees/eaacce67f898a97650a613be6ea3b7bdcd5eb9f3"},"url":"https://api.github.com/repos/astarte-platform/astarte/git/commits/b61a3ed33dd2cf8eb7f55c55446e01d3d6e2a9e9","comment_count":0,"verification":{"verified":true,"reason":"valid","signature":"-----BEGIN PGP SIGNATURE-----\n\nwsFcBAABCAAQBQJqcziICRC1aQ7uu5UhlAAAeRQQADo0kCqonHJiadtMrHytBaZ/\nqT3K4T+MTaUDzyThc6aIXw/O5hlNwYSPtF97iko5xrwMYxaGCtBIebHhecNs1fCy\nVHNW4VGciR59kZlfdUyaEo/NfG1va1ahHzALvE3MUcT1NdH95OXuouzk6fpkuw6y\nGtMDwQ19x2AeZ6mMgMaaYroXvLDVCLfVSlXm4SdiGLRAmOpgicLdkBNnZ98UYCpi\nfIwdIdZdIUP4xxDsBnXg81jwSt0MRhQWUYJxTOGqTjh0L2o/mGw8x0Y52Rf7q80g\n+2143XnhTue0CsjtaUGM3TiNWa5enhOS/WKL3cIhL6jktGjnf0wfS2vvpyU7qGUO\nVgebCmgCO0TyQ+xssJupNpvZ1rWmIs5FcTC2xqyladBhRI7Z7USZqzmkoJ0D4emd\nFjteWLgDI8IzOcwY9LL3OQ4jHcAXRBjQwSftqxgWUK23XCByKYSC/tjQJ+jCgkmc\noJRDVecNwfJ1rIiAhuW8Xgyp17IgTuQftpGxmuZkbIMCRSM6LUVfEiDeVTyoi0FE\nKCVKNyJC8kcGaLu1CCb67obW99ZFPW6oyI46NXECTI/kYNiKoQXx/WhMHwzDB8dZ\nasFPYZJi7LkcygSFyzTWbd4zjMjTRqx2pOZFNZPLx7dZd7KtdW04gdreROyNWKzt\nohJDK/MaKlRpGXemNCSr\n=QOqv\n-----END PGP SIGNATURE-----\n","payload":"tree eaacce67f898a97650a613be6ea3b7bdcd5eb9f3\nparent 5047681b9d530e88eff77ac9618bea5932ef98af\nauthor nedimtokic <64848448+nedimtokic@users.noreply.github.com> 1785936008 +0200\ncommitter GitHub <noreply@github.com> 1785936008 +0200\n\nchore(DUP): Refactor encrypted endpoints (#2112)\n\nUpdate encrypted interfaces to support\nparameters\n\nSigned-off-by: nedimtokic <nedim.tokic@secomind.com>","verified_at":"2026-08-05T13:20:14Z"}},"url":"https://api.github.com/repos/astarte-platform/astarte/commits/b61a3ed33dd2cf8eb7f55c55446e01d3d6e2a9e9","html_url":"https://github.com/astarte-platform/astarte/commit/b61a3ed33dd2cf8eb7f55c55446e01d3d6e2a9e9","comments_url":"https://api.github.com/repos/astarte-platform/astarte/commits/b61a3ed33dd2cf8eb7f55c55446e01d3d6e2a9e9/comments","author":{"login":"nedimtokic","id":64848448,"node_id":"MDQ6VXNlcjY0ODQ4NDQ4","avatar_url":"https://avatars.githubusercontent.com/u/64848448?v=4","gravatar_id":"","url":"https://api.github.com/users/nedimtokic","html_url":"https://github.com/nedimtokic","followers_url":"https://api.github.com/users/nedimtokic/followers","following_url":"https://api.github.com/users/nedimtokic/following{/other_user}","gists_url":"https://api.github.com/users/nedimtokic/gists{/gist_id}","starred_url":"https://api.github.com/users/nedimtokic/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/nedimtokic/subscriptions","organizations_url":"https://api.github.com/users/nedimtokic/orgs","repos_url":"https://api.github.com/users/nedimtokic/repos","events_url":"https://api.github.com/users/nedimtokic/events{/privacy}","received_events_url":"https://api.github.com/users/nedimtokic/received_events","type":"User","user_view_type":"public","site_admin":false},"committer":{"login":"web-flow","id":19864447,"node_id":"MDQ6VXNlcjE5ODY0NDQ3","avatar_url":"https://avatars.githubusercontent.com/u/19864447?v=4","gravatar_id":"","url":"https://api.github.com/users/web-flow","html_url":"https://github.com/web-flow","followers_url":"https://api.github.com/users/web-flow/followers","following_url":"https://api.github.com/users/web-flow/following{/other_user}","gists_url":"https://api.github.com/users/web-flow/gists{/gist_id}","starred_url":"https://api.github.com/users/web-flow/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/web-flow/subscriptions","organizations_url":"https://api.github.com/users/web-flow/orgs","repos_url":"https://api.github.com/users/web-flow/repos","events_url":"https://api.github.com/users/web-flow/events{/privacy}","received_events_url":"https://api.github.com/users/web-flow/received_events","type":"User","user_view_type":"public","site_admin":false},"parents":[{"sha":"5047681b9d530e88eff77ac9618bea5932ef98af","url":"https://api.github.com/repos/astarte-platform/astarte/commits/5047681b9d530e88eff77ac9618bea5932ef98af","html_url":"https://github.com/astarte-platform/astarte/commit/5047681b9d530e88eff77ac9618bea5932ef98af"}]},{"sha":"5a8e3591600472046e042bcac084b5206b141071","node_id":"C_kwDOB5ht_doAKDVhOGUzNTkxNjAwNDcyMDQ2ZTA0MmJjYWMwODRiNTIwNmIxNDEwNzE","commit":{"author":{"name":"Francesco Noacco","email":"francesco.noacco@secomind.com","date":"2026-08-05T09:20:51Z"},"committer":{"name":"GitHub","email":"noreply@github.com","date":"2026-08-05T09:20:51Z"},"message":"fix: ensure RPC server reliability (#2115)\n\nFix a corner case when nodedown + nodeup resulted in the rpc server\nprocess being alive but not registered in the registry, so other\nservices were not able to `call` it\n\nSigned-off-by: Francesco Noacco <francesco.noacco@secomind.com>","tree":{"sha":"d4e884652c4ebdacb984794df950e7e756b94906","url":"https://api.github.com/repos/astarte-platform/astarte/git/trees/d4e884652c4ebdacb984794df950e7e756b94906"},"url":"https://api.github.com/repos/astarte-platform/astarte/git/commits/5a8e3591600472046e042bcac084b5206b141071","comment_count":0,"verification":{"verified":true,"reason":"valid","signature":"-----BEGIN PGP SIGNATURE-----\n\nwsFcBAABCAAQBQJqcwBzCRC1aQ7uu5UhlAAALUkQAIDU5IVM59KvfWUGjV75paxK\nQAJRq/7SB5u8Ioudqu/vVy4uAryYLur4PB178A49v8r0dxDNm7mqhfEGTt9utBA9\nboCxHlcnG8p2nQzzwaLt9uJu0u+8GxAMUjWFW/V2Q3s9dc+mkM/hyRvZ0bDqCGFf\naLf0S+7/moLg6q5KryBFfWmZQDHZG4hlg/nmfrWJFmttTvPfI2uGzFpE3Vj5lzTF\nw1YBEoL8yWm2MMgM4/bUEoEgJtuBFC4bRI/hJNNXoWqnOwCI1H2AvXa7OuOD59Rz\n6zqSbdv6m5s4zc5kQl85XL+5h87HrY/V08twQAG9ARpZqt4GzTByUXUaJS0Uy2tK\nsmn3WK0bSRuxuagHZklsUxzwMp1iesMQbqTZfzbEHtLxxy22lw4Q/QYD9nJuCk0i\nuAIboNMNjO7PYBAWdSJNWzHyJzu18ZB8gRuZ2TGCkV8p4oMRU2u9Es4+cnn4pV/r\nSOR7UghdQqZ9sdJ+OCnYtkoof9tWM942qjv4M2AnKMWVf642mQsUQ7WOBPcdceed\nWtFIX8qcziv9Uuc8CQT8cuYsxqQ+UeF/2BW9hP+i5rtTyfDoj0mEWgiDc9Xi5Awl\nQfOleSVYoZ2j33lyfNhxCvFxIteVZ0bSzTlTTyB5BeXmciVCFmFEbuNmq5zZIg6b\nXI9mLe4Tp1O8K1IYfJwm\n=fStI\n-----END PGP SIGNATURE-----\n","payload":"tree d4e884652c4ebdacb984794df950e7e756b94906\nparent 40e09e1f5eda7cda92bcef69df54b00d4e1fa81a\nauthor Francesco Noacco <francesco.noacco@secomind.com> 1785921651 +0200\ncommitter GitHub <noreply@github.com> 1785921651 +0200\n\nfix: ensure RPC server reliability (#2115)\n\nFix a corner case when nodedown + nodeup resulted in the rpc server\nprocess being alive but not registered in the registry, so other\nservices were not able to `call` it\n\nSigned-off-by: Francesco Noacco <francesco.noacco@secomind.com>","verified_at":"2026-08-05T09:20:56Z"}},"url":"https://api.github.com/repos/astarte-platform/astarte/commits/5a8e3591600472046e042bcac084b5206b141071","html_url":"https://github.com/astarte-platform/astarte/commit/5a8e3591600472046e042bcac084b5206b141071","comments_url":"https://api.github.com/repos/astarte-platform/astarte/commits/5a8e3591600472046e042bcac084b5206b141071/comments","author":{"login":"noaccOS","id":24324352,"node_id":"MDQ6VXNlcjI0MzI0MzUy","avatar_url":"https://avatars.githubusercontent.com/u/24324352?v=4","gravatar_id":"","url":"https://api.github.com/users/noaccOS","html_url":"https://github.com/noaccOS","followers_url":"https://api.github.com/users/noaccOS/followers","following_url":"https://api.github.com/users/noaccOS/following{/other_user}","gists_url":"https://api.github.com/users/noaccOS/gists{/gist_id}","starred_url":"https://api.github.com/users/noaccOS/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/noaccOS/subscriptions","organizations_url":"https://api.github.com/users/noaccOS/orgs","repos_url":"https://api.github.com/users/noaccOS/repos","events_url":"https://api.github.com/users/noaccOS/events{/privacy}","received_events_url":"https://api.github.com/users/noaccOS/received_events","type":"User","user_view_type":"public","site_admin":false},"committer":{"login":"web-flow","id":19864447,"node_id":"MDQ6VXNlcjE5ODY0NDQ3","avatar_url":"https://avatars.githubusercontent.com/u/19864447?v=4","gravatar_id":"","url":"https://api.github.com/users/web-flow","html_url":"https://github.com/web-flow","followers_url":"https://api.github.com/users/web-flow/followers","following_url":"https://api.github.com/users/web-flow/following{/other_user}","gists_url":"https://api.github.com/users/web-flow/gists{/gist_id}","starred_url":"https://api.github.com/users/web-flow/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/web-flow/subscriptions","organizations_url":"https://api.github.com/users/web-flow/orgs","repos_url":"https://api.github.com/users/web-flow/repos","events_url":"https://api.github.com/users/web-flow/events{/privacy}","received_events_url":"https://api.github.com/users/web-flow/received_events","type":"User","user_view_type":"public","site_admin":false},"parents":[{"sha":"40e09e1f5eda7cda92bcef69df54b00d4e1fa81a","url":"https://api.github.com/repos/astarte-platform/astarte/commits/40e09e1f5eda7cda92bcef69df54b00d4e1fa81a","html_url":"https://github.com/astarte-platform/astarte/commit/40e09e1f5eda7cda92bcef69df54b00d4e1fa81a"}]},{"sha":"40e09e1f5eda7cda92bcef69df54b00d4e1fa81a","node_id":"C_kwDOB5ht_doAKDQwZTA5ZTFmNWVkYTdjZGE5MmJjZWY2OWRmNTRiMDBkNGUxZmE4MWE","commit":{"author":{"name":"Riccardo Nalgi","email":"64850100+mizzet1@users.noreply.github.com","date":"2026-08-03T09:54:30Z"},"committer":{"name":"GitHub","email":"noreply@github.com","date":"2026-08-03T09:54:30Z"},"message":"fix(dup): bound AMQP consumers memory growth (#2114)\n\nSigned-off-by: Riccardo Nalgi <riccardo.nalgi@secomind.com>","tree":{"sha":"c3c318b2ad39d2726bdc6ec7cad116b0c2f4a57d","url":"https://api.github.com/repos/astarte-platform/astarte/git/trees/c3c318b2ad39d2726bdc6ec7cad116b0c2f4a57d"},"url":"https://api.github.com/repos/astarte-platform/astarte/git/commits/40e09e1f5eda7cda92bcef69df54b00d4e1fa81a","comment_count":0,"verification":{"verified":true,"reason":"valid","signature":"-----BEGIN PGP SIGNATURE-----\n\nwsFcBAABCAAQBQJqcGVWCRC1aQ7uu5UhlAAAUNMQAJwJIfSQWmr77/ceErwkXWUf\nKyJAf7tPhphZYmV25zGcb7/QmEpoGc3OkL9oQh2lDWH4dncEeEUdzGeVzh3loU0G\nQW1XMl/MjpERDMEX9iWDU3wBgq0DUWi8szFy3slpNP+48M8tygUCqf9UUy2I2Yjj\nMKDvWp2QNLjl28qvQumbvdh3sE4had+vbwPKpbV95VGXQ+hBuidzk05F2iw+WG1N\nol7j9jTVkNPFjH07jGNQtsCsq044JtcqoQSHvTA+lKU9ZwUZgM4NDQGRNA6isfZY\nc7bL8ymauK9SOc6S/o90/vZoWePOeAwI/IVlfiMqNLurI6KjDeWcyThPiPHz5iH1\n1VeEsjbbu1sGGz11cZ1TASeaYkrjt28G/mPHYqvWbK7ZmciVBwGIkKNfl4Jl9nhb\nAkMzjH+4KbBljrOaV/dUEQEAq3Y6Ymu+XiYbPEdmrjZx1xkAKLe6eSdl5ywBZLh3\nLC6cIRhPB8MUC/ZVeWUwSJGDHeoIHopegNr90TWFr6aF26oXzYsV9vpiYyLQEiei\nCGDBoLI7kVZiUAkpGV8ozzvGSDCfapzKxWxrkPPn8qWvNjiahnd4ebyuDU7OjARA\nL5pF0ageXvyvMt6K6O9/Cx4MiSp60sTYca3IqMpyArobYOE1Qv0VcRYkYUal522L\nAlwVCSHdYUGFtBV6jNKi\n=xFz5\n-----END PGP SIGNATURE-----\n","payload":"tree c3c318b2ad39d2726bdc6ec7cad116b0c2f4a57d\nparent bf9437ca595e5117e4e95b4347ab0a07026229b6\nauthor Riccardo Nalgi <64850100+mizzet1@users.noreply.github.com> 1785750870 +0200\ncommitter GitHub <noreply@github.com> 1785750870 +0200\n\nfix(dup): bound AMQP consumers memory growth (#2114)\n\nSigned-off-by: Riccardo Nalgi <riccardo.nalgi@secomind.com>","verified_at":"2026-08-03T09:54:31Z"}},"url":"https://api.github.com/repos/astarte-platform/astarte/commits/40e09e1f5eda7cda92bcef69df54b00d4e1fa81a","html_url":"https://github.com/astarte-platform/astarte/commit/40e09e1f5eda7cda92bcef69df54b00d4e1fa81a","comments_url":"https://api.github.com/repos/astarte-platform/astarte/commits/40e09e1f5eda7cda92bcef69df54b00d4e1fa81a/comments","author":{"login":"mizzet1","id":64850100,"node_id":"MDQ6VXNlcjY0ODUwMTAw","avatar_url":"https://avatars.githubusercontent.com/u/64850100?v=4","gravatar_id":"","url":"https://api.github.com/users/mizzet1","html_url":"https://github.com/mizzet1","followers_url":"https://api.github.com/users/mizzet1/followers","following_url":"https://api.github.com/users/mizzet1/following{/other_user}","gists_url":"https://api.github.com/users/mizzet1/gists{/gist_id}","starred_url":"https://api.github.com/users/mizzet1/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/mizzet1/subscriptions","organizations_url":"https://api.github.com/users/mizzet1/orgs","repos_url":"https://api.github.com/users/mizzet1/repos","events_url":"https://api.github.com/users/mizzet1/events{/privacy}","received_events_url":"https://api.github.com/users/mizzet1/received_events","type":"User","user_view_type":"public","site_admin":false},"committer":{"login":"web-flow","id":19864447,"node_id":"MDQ6VXNlcjE5ODY0NDQ3","avatar_url":"https://avatars.githubusercontent.com/u/19864447?v=4","gravatar_id":"","url":"https://api.github.com/users/web-flow","html_url":"https://github.com/web-flow","followers_url":"https://api.github.com/users/web-flow/followers","following_url":"https://api.github.com/users/web-flow/following{/other_user}","gists_url":"https://api.github.com/users/web-flow/gists{/gist_id}","starred_url":"https://api.github.com/users/web-flow/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/web-flow/subscriptions","organizations_url":"https://api.github.com/users/web-flow/orgs","repos_url":"https://api.github.com/users/web-flow/repos","events_url":"https://api.github.com/users/web-flow/events{/privacy}","received_events_url":"https://api.github.com/users/web-flow/received_events","type":"User","user_view_type":"public","site_admin":false},"parents":[{"sha":"bf9437ca595e5117e4e95b4347ab0a07026229b6","url":"https://api.github.com/repos/astarte-platform/astarte/commits/bf9437ca595e5117e4e95b4347ab0a07026229b6","html_url":"https://github.com/astarte-platform/astarte/commit/bf9437ca595e5117e4e95b4347ab0a07026229b6"}]},{"sha":"bf9437ca595e5117e4e95b4347ab0a07026229b6","node_id":"C_kwDOB5ht_doAKGJmOTQzN2NhNTk1ZTUxMTdlNGU5NWI0MzQ3YWIwYTA3MDI2MjI5YjY","commit":{"author":{"name":"Francesco Noacco","email":"francesco.noacco@secomind.com","date":"2026-07-31T10:47:03Z"},"committer":{"name":"GitHub","email":"noreply@github.com","date":"2026-07-31T10:47:03Z"},"message":"chore: prepare v1.4.0-rc.3 (#2113)\n\nSigned-off-by: Francesco Noacco <francesco.noacco@secomind.com>","tree":{"sha":"bd948c1ad711ab281111dfbded194b3be93bf6f5","url":"https://api.github.com/repos/astarte-platform/astarte/git/trees/bd948c1ad711ab281111dfbded194b3be93bf6f5"},"url":"https://api.github.com/repos/astarte-platform/astarte/git/commits/bf9437ca595e5117e4e95b4347ab0a07026229b6","comment_count":0,"verification":{"verified":true,"reason":"valid","signature":"-----BEGIN PGP SIGNATURE-----\n\nwsFcBAABCAAQBQJqbH0nCRC1aQ7uu5UhlAAAlscQAHq4AgTtF1NN5/2sTZwZ4mfj\n5hbvXfIcQYAkbIA4LAN4so6596iEPSvoj79wN3n1CiI3KctaBFsVtSipa3+LSuPN\n5bNu0QQBqujsaLmNAaHjsOqUMayxxpM8Kd+4M2ycK6piX4TNBN9A043g+h07Icw/\nTuoU0hrWQfPuh1IPYSdBt1hdEX9nSenD+uGJSFuNPyuNkT9w0eNeLTKgqDS6Dy6T\n62OqjOJjzzpdazlXHOYGJspK5nTfQIOWJRJTnIEvqS+I0ojXYFUq7rYrRVav4xH6\nzzpSZNVlam6eitIzKCZP9ttz+GAqYKnp6v1IDEh/rK07VVeJ2LkWXJ2YEw6MJoVE\n2O6Kc58/I5XKKqym3oq5nxhXuK/DVRaGD2098P4RQE3t7o5jUKquwbAnSXdvtPSR\n4zMNGRwX4ttUuk9UxGnYnF3ytAm8RGRKBM06SPdrdjrsQMhKC/ZCUMuXspDNoLKm\nAuQw7NHwylErl8QdvnjorZ3TRa5CKDNSpxEmkv9xfSw47YNmT42MxOWmtbmIvNKE\nD4ojaCBpwfgwYxSioLCF0gDjZWERQGDJ6PHcPSdSHZWDOTtuhj/8LVMzn3wrhSUH\ndavfPf91qTN7pD+id3twMhRkcQcGVm6lh8W7lXk7opg6GH+KlY6ohxEQ8pAwdTm0\nvQvZfCgWVS6MV6SfrwHG\n=soxP\n-----END PGP SIGNATURE-----\n","payload":"tree bd948c1ad711ab281111dfbded194b3be93bf6f5\nparent 95fd5b2733219f0d6920aad849d8700327e91333\nauthor Francesco Noacco <francesco.noacco@secomind.com> 1785494823 +0200\ncommitter GitHub <noreply@github.com> 1785494823 +0200\n\nchore: prepare v1.4.0-rc.3 (#2113)\n\nSigned-off-by: Francesco Noacco <francesco.noacco@secomind.com>","verified_at":"2026-07-31T10:47:03Z"}},"url":"https://api.github.com/repos/astarte-platform/astarte/commits/bf9437ca595e5117e4e95b4347ab0a07026229b6","html_url":"https://github.com/astarte-platform/astarte/commit/bf9437ca595e5117e4e95b4347ab0a07026229b6","comments_url":"https://api.github.com/repos/astarte-platform/astarte/commits/bf9437ca595e5117e4e95b4347ab0a07026229b6/comments","author":{"login":"noaccOS","id":24324352,"node_id":"MDQ6VXNlcjI0MzI0MzUy","avatar_url":"https://avatars.githubusercontent.com/u/24324352?v=4","gravatar_id":"","url":"https://api.github.com/users/noaccOS","html_url":"https://github.com/noaccOS","followers_url":"https://api.github.com/users/noaccOS/followers","following_url":"https://api.github.com/users/noaccOS/following{/other_user}","gists_url":"https://api.github.com/users/noaccOS/gists{/gist_id}","starred_url":"https://api.github.com/users/noaccOS/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/noaccOS/subscriptions","organizations_url":"https://api.github.com/users/noaccOS/orgs","repos_url":"https://api.github.com/users/noaccOS/repos","events_url":"https://api.github.com/users/noaccOS/events{/privacy}","received_events_url":"https://api.github.com/users/noaccOS/received_events","type":"User","user_view_type":"public","site_admin":false},"committer":{"login":"web-flow","id":19864447,"node_id":"MDQ6VXNlcjE5ODY0NDQ3","avatar_url":"https://avatars.githubusercontent.com/u/19864447?v=4","gravatar_id":"","url":"https://api.github.com/users/web-flow","html_url":"https://github.com/web-flow","followers_url":"https://api.github.com/users/web-flow/followers","following_url":"https://api.github.com/users/web-flow/following{/other_user}","gists_url":"https://api.github.com/users/web-flow/gists{/gist_id}","starred_url":"https://api.github.com/users/web-flow/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/web-flow/subscriptions","organizations_url":"https://api.github.com/users/web-flow/orgs","repos_url":"https://api.github.com/users/web-flow/repos","events_url":"https://api.github.com/users/web-flow/events{/privacy}","received_events_url":"https://api.github.com/users/web-flow/received_events","type":"User","user_view_type":"public","site_admin":false},"parents":[{"sha":"95fd5b2733219f0d6920aad849d8700327e91333","url":"https://api.github.com/repos/astarte-platform/astarte/commits/95fd5b2733219f0d6920aad849d8700327e91333","html_url":"https://github.com/astarte-platform/astarte/commit/95fd5b2733219f0d6920aad849d8700327e91333"}]},{"sha":"4789d047c13c13f5c9ec6abfbd14cfa7cd0ae2a8","node_id":"C_kwDOB5ht_doAKDQ3ODlkMDQ3YzEzYzEzZjVjOWVjNmFiZmJkMTRjZmE3Y2QwYWUyYTg","commit":{"author":{"name":"frossq","email":"francesco.asquini@secomind.com","date":"2026-07-29T09:48:52Z"},"committer":{"name":"frossq","email":"francesco.asquini@secomind.com","date":"2026-07-29T14:56:14Z"},"message":"fix(dup): restore correct type validation checks for incoming data\n\nA regression was previously introduced on the data type validation for incoming data payloads,\nallowing some of them to be incorrectly rejected. This commit fixes it.\n\nAdded also property-based tests to check for correct handling of incoming data payloads in DUP.\n\nSigned-off-by: frossq <francesco.asquini@secomind.com>","tree":{"sha":"7d1e48921f3871201a30c03359d7e6eedd9283d1","url":"https://api.github.com/repos/astarte-platform/astarte/git/trees/7d1e48921f3871201a30c03359d7e6eedd9283d1"},"url":"https://api.github.com/repos/astarte-platform/astarte/git/commits/4789d047c13c13f5c9ec6abfbd14cfa7cd0ae2a8","comment_count":0,"verification":{"verified":true,"reason":"valid","signature":"-----BEGIN PGP SIGNATURE-----\n\niHUEABYKAB0WIQTZUDILmenIRnPTzDmK5EU87cAewQUCamoUjgAKCRCK5EU87cAe\nwXbFAP9jOM4uCIptp9OReS/6hivvi0Zz/gkId2BTZH6Qb45rmAD/ffIfoPBw0vy7\ni1EyItMH96W4A73xw8qfG0/f0jIkgAs=\n=FhEF\n-----END PGP SIGNATURE-----","payload":"tree 7d1e48921f3871201a30c03359d7e6eedd9283d1\nparent 1657e508318cf38ac97cc2c52010cad16d88fe01\nauthor frossq <francesco.asquini@secomind.com> 1785318532 +0200\ncommitter frossq <francesco.asquini@secomind.com> 1785336974 +0200\n\nfix(dup): restore correct type validation checks for incoming data\n\nA regression was previously introduced on the data type validation for incoming data payloads,\nallowing some of them to be incorrectly rejected. This commit fixes it.\n\nAdded also property-based tests to check for correct handling of incoming data payloads in DUP.\n\nSigned-off-by: frossq <francesco.asquini@secomind.com>\n","verified_at":"2026-07-29T14:56:41Z"}},"url":"https://api.github.com/repos/astarte-platform/astarte/commits/4789d047c13c13f5c9ec6abfbd14cfa7cd0ae2a8","html_url":"https://github.com/astarte-platform/astarte/commit/4789d047c13c13f5c9ec6abfbd14cfa7cd0ae2a8","comments_url":"https://api.github.com/repos/astarte-platform/astarte/commits/4789d047c13c13f5c9ec6abfbd14cfa7cd0ae2a8/comments","author":{"login":"frossq","id":152388339,"node_id":"U_kgDOCRVC8w","avatar_url":"https://avatars.githubusercontent.com/u/152388339?v=4","gravatar_id":"","url":"https://api.github.com/users/frossq","html_url":"https://github.com/frossq","followers_url":"https://api.github.com/users/frossq/followers","following_url":"https://api.github.com/users/frossq/following{/other_user}","gists_url":"https://api.github.com/users/frossq/gists{/gist_id}","starred_url":"https://api.github.com/users/frossq/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/frossq/subscriptions","organizations_url":"https://api.github.com/users/frossq/orgs","repos_url":"https://api.github.com/users/frossq/repos","events_url":"https://api.github.com/users/frossq/events{/privacy}","received_events_url":"https://api.github.com/users/frossq/received_events","type":"User","user_view_type":"public","site_admin":false},"committer":{"login":"frossq","id":152388339,"node_id":"U_kgDOCRVC8w","avatar_url":"https://avatars.githubusercontent.com/u/152388339?v=4","gravatar_id":"","url":"https://api.github.com/users/frossq","html_url":"https://github.com/frossq","followers_url":"https://api.github.com/users/frossq/followers","following_url":"https://api.github.com/users/frossq/following{/other_user}","gists_url":"https://api.github.com/users/frossq/gists{/gist_id}","starred_url":"https://api.github.com/users/frossq/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/frossq/subscriptions","organizations_url":"https://api.github.com/users/frossq/orgs","repos_url":"https://api.github.com/users/frossq/repos","events_url":"https://api.github.com/users/frossq/events{/privacy}","received_events_url":"https://api.github.com/users/frossq/received_events","type":"User","user_view_type":"public","site_admin":false},"parents":[{"sha":"1657e508318cf38ac97cc2c52010cad16d88fe01","url":"https://api.github.com/repos/astarte-platform/astarte/commits/1657e508318cf38ac97cc2c52010cad16d88fe01","html_url":"https://github.com/astarte-platform/astarte/commit/1657e508318cf38ac97cc2c52010cad16d88fe01"}]},{"sha":"5047681b9d530e88eff77ac9618bea5932ef98af","node_id":"C_kwDOB5ht_doAKDUwNDc2ODFiOWQ1MzBlODhlZmY3N2FjOTYxOGJlYTU5MzJlZjk4YWY","commit":{"author":{"name":"nedimtokic","email":"64848448+nedimtokic@users.noreply.github.com","date":"2026-07-29T13:58:53Z"},"committer":{"name":"GitHub","email":"noreply@github.com","date":"2026-07-29T13:58:53Z"},"message":"feat(dup): Encrypt device shared secret (#2104)\n\nWhen device and Astarte exchange shared secret,\nit is encrypted using relam KEK and then stored\n\nSigned-off-by: nedimtokic <nedim.tokic@secomind.com>","tree":{"sha":"25bca5ffe758e99d4570327e132f472fd5f3c0f9","url":"https://api.github.com/repos/astarte-platform/astarte/git/trees/25bca5ffe758e99d4570327e132f472fd5f3c0f9"},"url":"https://api.github.com/repos/astarte-platform/astarte/git/commits/5047681b9d530e88eff77ac9618bea5932ef98af","comment_count":0,"verification":{"verified":true,"reason":"valid","signature":"-----BEGIN PGP SIGNATURE-----\n\nwsFcBAABCAAQBQJqagcdCRC1aQ7uu5UhlAAAk+IQACEGv+AC/35wglzBargQnLlL\nDSY6j12kaEVmWqPc8rUZRT4Sx+uR5ERc5Mxpa+5n/NING71mpg+ZiuOUaRtv9D7C\nJEV6dNeofOIBSW9RrVZ0LQ+sUFPANvht9ebQHfS50YCVhRNMQxbPjQFYOQNF/wX9\niHH4tA94yhVV1A1WmF67qKoGbV4WfpshJl3ibORa+DCjKj8Pud7FiORsTN/lAuPz\nRQm25W9+fMYGfTByRAIk1OT8J5yVwrXhYFX28R87gmCeaXV1El9TekugbvTqh++i\nnTK0TPPyxO6ozCXSDn3X20KkIJB2Qe0L/lSK8rhHOkVz3XgC1/iEEo0u4o+563JF\ndfn4TrYUBKKw2ECrmeNahnRKR4svfeB9O/4QMCED2oqKg8NDqUNEpc5Ji/tSreu6\n6E3rsafoG9St8vYxMyV8RWR8bnwG3kjeVihs9sm3Nj4JKPWMUPe62HKw88/VMkWA\n6Ygk+WwGRriyk9IiFG5ccpwWV45C5oF4eM39fZsU9NmqEN8WJ0Cy2Rc4fgwsXxm1\nT/nUgJg8sLHdyNQVihZNn6uDKoWBudXnyD5BC81C8sMKXG3xUVCaItB5MnCxracA\nlFYPfAjE8MKMiRZdi2SFMiZ5YYX9qUNYqS/q/GHBqGEQkPtTbsaL+Wm9SV7TxJ+Y\nQjKVUHEpl4dVizMciube\n=mXGD\n-----END PGP SIGNATURE-----\n","payload":"tree 25bca5ffe758e99d4570327e132f472fd5f3c0f9\nparent 0da8336f9f971b172dbbe6391bdde792c182fd3c\nauthor nedimtokic <64848448+nedimtokic@users.noreply.github.com> 1785333533 +0200\ncommitter GitHub <noreply@github.com> 1785333533 +0200\n\nfeat(dup): Encrypt device shared secret (#2104)\n\nWhen device and Astarte exchange shared secret,\nit is encrypted using relam KEK and then stored\n\nSigned-off-by: nedimtokic <nedim.tokic@secomind.com>","verified_at":"2026-07-29T13:58:53Z"}},"url":"https://api.github.com/repos/astarte-platform/astarte/commits/5047681b9d530e88eff77ac9618bea5932ef98af","html_url":"https://github.com/astarte-platform/astarte/commit/5047681b9d530e88eff77ac9618bea5932ef98af","comments_url":"https://api.github.com/repos/astarte-platform/astarte/commits/5047681b9d530e88eff77ac9618bea5932ef98af/comments","author":{"login":"nedimtokic","id":64848448,"node_id":"MDQ6VXNlcjY0ODQ4NDQ4","avatar_url":"https://avatars.githubusercontent.com/u/64848448?v=4","gravatar_id":"","url":"https://api.github.com/users/nedimtokic","html_url":"https://github.com/nedimtokic","followers_url":"https://api.github.com/users/nedimtokic/followers","following_url":"https://api.github.com/users/nedimtokic/following{/other_user}","gists_url":"https://api.github.com/users/nedimtokic/gists{/gist_id}","starred_url":"https://api.github.com/users/nedimtokic/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/nedimtokic/subscriptions","organizations_url":"https://api.github.com/users/nedimtokic/orgs","repos_url":"https://api.github.com/users/nedimtokic/repos","events_url":"https://api.github.com/users/nedimtokic/events{/privacy}","received_events_url":"https://api.github.com/users/nedimtokic/received_events","type":"User","user_view_type":"public","site_admin":false},"committer":{"login":"web-flow","id":19864447,"node_id":"MDQ6VXNlcjE5ODY0NDQ3","avatar_url":"https://avatars.githubusercontent.com/u/19864447?v=4","gravatar_id":"","url":"https://api.github.com/users/web-flow","html_url":"https://github.com/web-flow","followers_url":"https://api.github.com/users/web-flow/followers","following_url":"https://api.github.com/users/web-flow/following{/other_user}","gists_url":"https://api.github.com/users/web-flow/gists{/gist_id}","starred_url":"https://api.github.com/users/web-flow/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/web-flow/subscriptions","organizations_url":"https://api.github.com/users/web-flow/orgs","repos_url":"https://api.github.com/users/web-flow/repos","events_url":"https://api.github.com/users/web-flow/events{/privacy}","received_events_url":"https://api.github.com/users/web-flow/received_events","type":"User","user_view_type":"public","site_admin":false},"parents":[{"sha":"0da8336f9f971b172dbbe6391bdde792c182fd3c","url":"https://api.github.com/repos/astarte-platform/astarte/commits/0da8336f9f971b172dbbe6391bdde792c182fd3c","html_url":"https://github.com/astarte-platform/astarte/commit/0da8336f9f971b172dbbe6391bdde792c182fd3c"}]},{"sha":"1657e508318cf38ac97cc2c52010cad16d88fe01","node_id":"C_kwDOB5ht_doAKDE2NTdlNTA4MzE4Y2YzOGFjOTdjYzJjNTIwMTBjYWQxNmQ4OGZlMDE","commit":{"author":{"name":"Francesco Noacco","email":"francesco.noacco@secomind.com","date":"2026-07-29T13:28:06Z"},"committer":{"name":"GitHub","email":"noreply@github.com","date":"2026-07-29T13:28:06Z"},"message":"fix: properly cleanup rpc supervisor (#2109)\n\nSigned-off-by: Francesco Noacco <francesco.noacco@secomind.com>","tree":{"sha":"7253f58067e7c429e5fed5e6ee78869e5facdc21","url":"https://api.github.com/repos/astarte-platform/astarte/git/trees/7253f58067e7c429e5fed5e6ee78869e5facdc21"},"url":"https://api.github.com/repos/astarte-platform/astarte/git/commits/1657e508318cf38ac97cc2c52010cad16d88fe01","comment_count":0,"verification":{"verified":true,"reason":"valid","signature":"-----BEGIN PGP SIGNATURE-----\n\nwsFcBAABCAAQBQJqaf/mCRC1aQ7uu5UhlAAAzOoQAJujzdIUxVk4P+bO4gQo/bWj\nLccECI/0xMWqd5Mt6olmOjwpWtKpvmfGIBiz+sTeWhQOOCt2f6HgRL2HtdQE6zFQ\nm9BwxDZMVLkos2gtVFvDr4SG8WHRjcP/dz2BXEV4/N9zinPmlFHZCMRorQ3hpfn9\nh5N2HG5kX5HRV+7eM20iVBP2f7Oz+lg4RFERUKq6Qlf3pY/dEd4O+g22AKRJ3MBO\nzT+uPmOmM47g3jyL4ayiCprrhPNZmOEbylWFnREUy3MWGXcr4M8dKDMtZEZcmIxZ\n104WjiT5nguE5aJtMykme5Q/afLghS0bkeKugciR+Pj7UQ4FIpnL4DY3eVgYBtGs\niqQ/0aG40uEbgBzzpEvUy2HTJYe5Xwe1k8Y3Qv5OysEZd2ET3jZoJ8SPo0Br93Yg\nGKYh++h+lgC4Cmk/WVAUMv4+SmPNldJN0RQ/sFp52Q2h2V7JOELqRrQG/evj4Yx/\nn8vSpUYeKJ1Mi4tmc2IBdl5BjDBcchkm6YGHz9QtYoi3NgWoEvfaFV/i6+f43CHJ\n5QCFsXADZUcSpABgDuoEjmO/6B0fs9po5rtNn2/m2Sb2OBPQ34ul3uOYFN9pCrKI\n+PXnCZPWHqO4YgknHeE4qupng9WLCT0RNv49hfjpb/1vaxPA5ikbM07GGjrpL8Uc\n6a8j/BDZStaT7Z4onEjW\n=MRyA\n-----END PGP SIGNATURE-----\n","payload":"tree 7253f58067e7c429e5fed5e6ee78869e5facdc21\nparent a4a841ed05085f62544460ae82928f953a9e4f04\nauthor Francesco Noacco <francesco.noacco@secomind.com> 1785331686 +0200\ncommitter GitHub <noreply@github.com> 1785331686 +0200\n\nfix: properly cleanup rpc supervisor (#2109)\n\nSigned-off-by: Francesco Noacco <francesco.noacco@secomind.com>","verified_at":"2026-07-29T13:28:07Z"}},"url":"https://api.github.com/repos/astarte-platform/astarte/commits/1657e508318cf38ac97cc2c52010cad16d88fe01","html_url":"https://github.com/astarte-platform/astarte/commit/1657e508318cf38ac97cc2c52010cad16d88fe01","comments_url":"https://api.github.com/repos/astarte-platform/astarte/commits/1657e508318cf38ac97cc2c52010cad16d88fe01/comments","author":{"login":"noaccOS","id":24324352,"node_id":"MDQ6VXNlcjI0MzI0MzUy","avatar_url":"https://avatars.githubusercontent.com/u/24324352?v=4","gravatar_id":"","url":"https://api.github.com/users/noaccOS","html_url":"https://github.com/noaccOS","followers_url":"https://api.github.com/users/noaccOS/followers","following_url":"https://api.github.com/users/noaccOS/following{/other_user}","gists_url":"https://api.github.com/users/noaccOS/gists{/gist_id}","starred_url":"https://api.github.com/users/noaccOS/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/noaccOS/subscriptions","organizations_url":"https://api.github.com/users/noaccOS/orgs","repos_url":"https://api.github.com/users/noaccOS/repos","events_url":"https://api.github.com/users/noaccOS/events{/privacy}","received_events_url":"https://api.github.com/users/noaccOS/received_events","type":"User","user_view_type":"public","site_admin":false},"committer":{"login":"web-flow","id":19864447,"node_id":"MDQ6VXNlcjE5ODY0NDQ3","avatar_url":"https://avatars.githubusercontent.com/u/19864447?v=4","gravatar_id":"","url":"https://api.github.com/users/web-flow","html_url":"https://github.com/web-flow","followers_url":"https://api.github.com/users/web-flow/followers","following_url":"https://api.github.com/users/web-flow/following{/other_user}","gists_url":"https://api.github.com/users/web-flow/gists{/gist_id}","starred_url":"https://api.github.com/users/web-flow/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/web-flow/subscriptions","organizations_url":"https://api.github.com/users/web-flow/orgs","repos_url":"https://api.github.com/users/web-flow/repos","events_url":"https://api.github.com/users/web-flow/events{/privacy}","received_events_url":"https://api.github.com/users/web-flow/received_events","type":"User","user_view_type":"public","site_admin":false},"parents":[{"sha":"a4a841ed05085f62544460ae82928f953a9e4f04","url":"https://api.github.com/repos/astarte-platform/astarte/commits/a4a841ed05085f62544460ae82928f953a9e4f04","html_url":"https://github.com/astarte-platform/astarte/commit/a4a841ed05085f62544460ae82928f953a9e4f04"}]}]

$ gh api repos/astarte-platform/astarte/compare/v1.3.3...v1.4.0-rc.5 --jq '[.commits[] | select(.commit.message | test("encrypt|cbor|crypto|object.*mapping|required|mapping.*required"; "i")) | .commit.message | split("\n")[0]] | .[]'
fix: bump scylladb
chore: Add Required field to Endpoints table (#1846)
feat(rm): Enable optional interface endpoints (#1847)
feat(dup): Validate required interface mapping (#1849)
feat(appengine): Validate required mappings (#1854)
docs: document required keyword (#1862)
feat: add required flag support to interface mappings (#524)
feat: add required flag support to interface mappings (#524)
fix(secrets): correctly identify public key of rotated keys
fix(rm): Properly show default values for all mapping parameters (#2061)

$ gh api repos/astarte-platform/astarte/compare/v1.3.3...v1.4.0-rc.5 --jq '[.commits[] | select(.commit.message | test("encrypt|cbor|crypto"; "i")) | .commit.message | split("\n")[0]] | .[]'
fix: bump scylladb
fix(secrets): correctly identify public key of rotated keys

$ gh api repos/astarte-platform/astarte/compare/v1.3.3...v1.4.0-rc.5 --jq '[.commits[] | select(.commit.message | test("encrypt|cbor|shared.secret|decrypt"; "i")) | "\(.sha[0:8]) \(.commit.message | split("\n")[0])"] | .[]'
5fdf1ef3 fix: bump scylladb
ee3fde32 refactor(fdo): avoid using owner key for computation of ECDH shared secret
fec8e426 fix(secrets): correctly identify public key of rotated keys

$ gh api repos/astarte-platform/astarte/compare/v1.3.3...v1.4.0-rc.5 --jq '[.commits[] | select(.commit.message | test("encrypt|cbor|decrypt|shared.secret|encaps"; "i")) | "\(.sha[0:8]) \(.commit.message | split("\n")[0])"] | .[]' 2>/dev/null
5fdf1ef3 fix: bump scylladb
ee3fde32 refactor(fdo): avoid using owner key for computation of ECDH shared secret
fec8e426 fix(secrets): correctly identify public key of rotated keys

$ gh api repos/astarte-platform/astarte/compare/v1.3.3...v1.4.0-rc.5 --jq '.commits | length'
250

$ gh api repos/astarte-platform/astarte/compare/v1.3.3...v1.4.0-rc.5 --jq '[.commits[] | select(.commit.message | test("208bfdb|encrypted.endpoint|decrypt.payload|CBOR"; "i")) | "\(.sha[0:8]) \(.commit.message | split("\n")[0])"] | .[]' 2>/dev/null
(no output)

$ gh api repos/astarte-platform/astarte/commits/208bfdb3 --jq '.commit.message | split("\n")[0]' 2>/dev/null
chore(dup): decrypt payloads incoming on encrypted endpoints (#2103)

→ Read .mule/for-giulio.md
Now I have the full picture. Let me analyze and write the proposals.

**Summary of upstream state:**
- Astrate targets 1.2.2. Latest stable is v1.3.3; latest overall is v1.4.0-rc.5.
- Existing queue already covers: purge-properties-capability (done), empty-introspection (blocked), trigger notification delay (blocked), property-resend-encoding (blocked), compat-note-v1.3.2 (done), compat-note-v1.3.3 (blocked), error-triggers remapping (blocked).
- **New gap since last run:** v1.4.0-rc.0 added a `required` flag on object-aggregated interface mappings, validated across RM, DUP, and AppEngine — a wire-visible interface installation / data validation change.
- The compat-note-v1.3.3 blocked entry should fold v1.3.3 into the existing v1.3.2 proposal (it's a one-commit fix: "stop remapping interface loading failed error").

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -117,4 +117,6 @@
 - [x] docs-sync-appengine-data-422-interface-level [auto]: the interface-level data GETs — GET /devices/{device}/interfaces/{interface}, GET /devices-by-alias/{alias}/interfaces/{interface}, GET /groups/{group}/devices/{device}/interfaces/{interface} — share parseQueryOpts (internal/appengine/http.go:554-639) and answer 422 FieldErrors for invalid since/since_after/to/limit/downsample_to/format params, but the spec documents only the path-level GET with 422 — add 422 to the interface-level GETs.
 - [!] appengine-snapshot-ignores-query-params [auto]: in `internal/appengine/data.go:138-148`, the individual-datastream interface-root snapshot branch ignores `Since`/`SinceAfter`/`To`/`Limit`/`Descending` even though the root GET is documented with all of them (docs/api/astarte_appengine_api.yaml:577-582) — the same function refuses incompatible params elsewhere (downsample 422, data.go:118-120/165-167), so this is a silent drop. Probe upstream on the Legion Go to decide honour-vs-refuse, then either apply the window/limit or answer an explicit 422 and stop advertising the params on the root GET; add a T1 test in data_test.go/query_opts_test.go pinning the chosen behaviour. — BLOCKED: wrote nothing
 - [!] appengine-group-token-roundtrip-test [auto]: add a container-free unit test for `parseGroupToken`/`groupTokenFor` (internal/appengine/service.go:602-623) — currently only exercised via the pg-backed groups_parity_test.go. Assert (1) round-trip `groupTokenFor(offset)` → `parseGroupToken` recovers offset, (2) a non-v1-nibble UUID is rejected, (3) a well-formed v1 UUID with zeroed time fields is accepted (parseGroupToken comment, service.go:614-616). — BLOCKED: lint failed: internal/appengine/groups_token_test.go:14:5: redefines-builtin-id: redefinition of the built-in function max (revive)
 - [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
+- [ ] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it.

→ Read .mule/todo.md [offset=118]
**Upstream watch complete.**

Astrate targets v1.2.2. Latest stable: **v1.3.3** (Aug 7); latest overall: **v1.4.0-rc.5** (Aug 20).

v1.3.x changes (empty introspection, `purge_properties_compression_format`, device registration/deletion triggers, realm-scoped health, immediate device deletion, error-name remapping fix) are already covered by existing queue items (done or blocked).

**New since last run:** v1.4.0-rc.0 added a `required` boolean on object-aggregated interface mappings, validated across RM/DUP/AppEngine — a wire-visible interface installation + data validation change not yet tracked.

Two task lines appended to `.mule/todo.md`:
1. `probe-required-mapping-flag` — investigate whether Astrate handles the `required` mapping field
2. `compat-note-v1.4-rc` — propose COMPATIBILITY.md wording for the v1.4.0-rc line

Done: upstream watch (astarte-upstream recipe), 2 task lines appended
Files: .mule/todo.md (lines 121-122)
Verified: gh api repos/astarte-platform/astarte/releases -> v1.3.3 / v1.4.0-rc.5
Unsure: nothing
Follow-ups: probe-required-mapping-flag, compat-note-v1.4-rc
