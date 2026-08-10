slug: recipe-astarte-upstream
verdict: proposed
at:  7bef58c
ran: 2026-08-10T02:11:51Z on DietPi in 431s


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

$ rg -n -i "upstream|astarte v1|target|1\.2\.2|1\.3" /root/astrate-mule/.mule/for-giulio.md 2>/dev/null | head -50; echo "---"; ls /root/astrate-mule/.mule/
---
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

→ Read .mule/todo.md

$ rg -n -i "1\.4\.0|1\.3\.3|1\.2\.2|compat-note|COMPATIBILITY" /root/astrate-mule/.mule/for-giulio.md | head -40; echo "=== reports ==="; ls /root/astrate-mule/.mule/reports/ | head -40
18:  2026-07-27: implement it.** Guiding principle clarified: Astarte compatibility means
19:  SDK/wire compatibility, not minimum dependency count — Astrate is allowed to be a
=== reports ===
issue-13.md
issue-20.md
log
race-check.md

→ Read .mule/for-giulio.md
→ Read .mule/log.md
The queue shows the last upstream review was around v1.4.0-rc.3 (probes and compat-note for it are still open). Newer tags exist now: v1.3.3 and v1.4.0-rc.4. Let me look at both release notes.

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

v1.3.3 has no release notes. Let me see what changed in it, and the full diff between rc.3 (last reviewed) and rc.4.
$ gh api "repos/astarte-platform/astarte/compare/v1.4.0-rc.3...v1.4.0-rc.4" --jq '.commits[] | "\(.sha[0:8]) \(.commit.message | split("\n")[0])"'
40e09e1f fix(dup): bound AMQP consumers memory growth (#2114)
5a8e3591 fix: ensure RPC server reliability (#2115)
745f6884 chore: init astarte_config
10728896 feat(config): add url_env and http client
f1faceeb refactor(fdo): use astarte_config for rendezvous configuration
625e1f67 refactor(secrets): use astarte_config for vault configuration
171b83c4 chore: prepare v1.4.0-rc.4 (#2121)

$ gh api "repos/astarte-platform/astarte/compare/v1.3.2...v1.3.3" --jq '.commits[] | "\(.sha[0:8]) \(.commit.message | split("\n")[0])"' 2>&1 | head -30; echo "=== v1.3.3 vs previous 1.3 tag ==="; gh api repos/astarte-platform/astarte/tags --jq '.[] | .name' | head -20
d8592091 fix(dup): stop remapping interface loading failed error (#2119)
3e70d63f chore: prepare v1.3.3 (#2128)
=== v1.3.3 vs previous 1.3 tag ===
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
v1.4.0-rc.0
v1.3.3
v1.3.2
v1.3.1
v1.3.0
v1.3.0-rc.2
v1.3.0-rc.1
v1.3.0-rc.0
v1.2.2
v1.2.2-rc.0
v1.2.1
v1.2.1-rc.1
v1.2.1-rc.0
v1.2.1-alpha.0
v1.2.0
v1.2.0-rc.0

The v1.3.3 fix "stop remapping interface loading failed error" looks directly relevant — Astrate's COMPATIBILITY.md documents an `interface_loading_failed` fallback in the error-name mapping. Let me look at it.
$ gh api repos/astarte-platform/astarte/commits/d8592091 --jq '.commit.message'
fix(dup): stop remapping interface loading failed error (#2119)

Signed-off-by: Riccardo Nalgi <riccardo.nalgi@secomind.com>

$ gh api "repos/astarte-platform/astarte/commits/d8592091" --jq '.files[] | "\(.filename)\t+\(.additions)/-\(.deletions)"'
CHANGELOG.md	+6/-0
apps/astarte_appengine_api/mix.lock	+1/-1
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex	+35/-2
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/device.ex	+28/-11
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/control_handler_test.exs	+29/-0
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/device_test.exs	+14/-1
apps/astarte_pairing/mix.lock	+1/-1

$ gh api repos/astarte-platform/astarte/commits/d8592091 --jq '.files[] | select(.filename=="apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/device.ex") | .patch'
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
+             realm,
+             device_id,
+             encoded_device_id,
+             interface_descriptor,
+             mapping
+           ) do
+        :ok ->
+          {:cont, :ok}
+
+        {:error, reason} ->
+          Logger.warning(
+            "Failed sending properties for interface #{interface_descriptor.name}: #{inspect(reason)}."
+          )
+
+          {:halt, {:error, :sending_properties_to_interface_failed}}
+      end
     end)
   end

$ gh api repos/astarte-platform/astarte/commits/d8592091 --jq '.files[] | select(.filename=="apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex") | .patch'
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

This is directly relevant to Astrate: upstream now emits a `device_error` with `error_name: "interface_loading_failed"` for a case it used to remap to `resend_interface_properties_failed`. Let me check Astrate's mapping.
$ rg -n "interface_loading_failed|resend_interface_properties_failed|sending_properties_to_interface_failed|UpstreamErrorName" /root/astrate-mule/internal /root/astrate-mule/test 2>/dev/null | head -40
/root/astrate-mule/test/conformance/upstream/channels.json:180:      "error_name": "interface_loading_failed",
/root/astrate-mule/test/conformance/upstream/channels.json:193:        "[null,null,\"rooms:bench:probe\",\"new_event\",{\"device_id\":\"KqYJ16Cs7lyeJ3PPzg7xQg\",\"event\":{\"error_name\":\"interface_loading_failed\",\"metadata\":{\"base64_payload\":\"EAAAAAF2AAAAAAAAAPA/AA==\",\"interface\":\"\\\"org.astrate.nonexistent.Iface\\\"\",\"path\":\"\\\"/value\\\"\"},\"type\":\"device_error\"},\"timestamp\":\"2026-07-26T20:54:39.413Z\"}]",
/root/astrate-mule/test/conformance/upstream/channels.json:194:        "[null,null,\"rooms:bench:probe\",\"new_event\",{\"device_id\":\"KqYJ16Cs7lyeJ3PPzg7xQg\",\"event\":{\"error_name\":\"interface_loading_failed\",\"metadata\":{\"base64_payload\":\"EAAAAAF2AAAAAAAAAPA/AA==\",\"interface\":\"\\\"org.astrate.nonexistent.Iface\\\"\",\"path\":\"\\\"/value\\\"\"},\"type\":\"device_error\"},\"timestamp\":\"2026-07-26T20:55:07.678Z\"}]",
/root/astrate-mule/test/conformance/upstream/channels.json:195:        "[null,null,\"rooms:bench:probe\",\"new_event\",{\"device_id\":\"KqYJ16Cs7lyeJ3PPzg7xQg\",\"event\":{\"error_name\":\"interface_loading_failed\",\"metadata\":{\"base64_payload\":\"EAAAAAF2AAAAAAAAAPA/AA==\",\"interface\":\"\\\"org.astrate.nonexistent.Iface\\\"\",\"path\":\"\\\"/value\\\"\"},\"type\":\"device_error\"},\"timestamp\":\"2026-07-26T20:55:35.939Z\"}]"
/root/astrate-mule/internal/engine/triggers/errorname.go:11:	"interface_loading_failed",
/root/astrate-mule/internal/engine/triggers/errorname.go:20:	"resend_interface_properties_failed",
/root/astrate-mule/internal/engine/triggers/errorname.go:27:const upstreamFallback = "interface_loading_failed"
/root/astrate-mule/internal/engine/triggers/errorname.go:37:// UpstreamErrorNames returns a copy of the closed set, so a caller cannot
/root/astrate-mule/internal/engine/triggers/errorname.go:39:func UpstreamErrorNames() []string {
/root/astrate-mule/internal/engine/triggers/errorname.go:47:	"interface_not_in_introspection": "interface_loading_failed",
/root/astrate-mule/internal/engine/triggers/errorname.go:48:	"interface_not_installed":        "interface_loading_failed",
/root/astrate-mule/internal/engine/triggers/errorname.go:65:// UpstreamErrorName maps one of Astrate's reject-reason labels to the upstream
/root/astrate-mule/internal/engine/triggers/errorname.go:67:// returned unchanged. Any other input maps to interface_loading_failed.
/root/astrate-mule/internal/engine/triggers/errorname.go:68:func UpstreamErrorName(reason string) string {
/root/astrate-mule/test/conformance/upstream/README.md:90:| interface not in introspection | `interface_not_in_introspection` → `invalid_interface` | `interface_loading_failed` | **wrong** |
/root/astrate-mule/test/conformance/upstream/README.md:100:than from what upstream does. `interface_loading_failed` reads like an internal failure and was
/root/astrate-mule/test/conformance/upstream/README.md:217:  loads no interface at all, and *every* publish then comes back as `interface_loading_failed` —
/root/astrate-mule/internal/engine/triggers/errorname_test.go:7:func TestUpstreamErrorNameMapping(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/errorname_test.go:15:		{"interface_not_in_introspection", "interface_loading_failed"},
/root/astrate-mule/internal/engine/triggers/errorname_test.go:16:		{"interface_not_installed", "interface_loading_failed"},
/root/astrate-mule/internal/engine/triggers/errorname_test.go:34:			got := UpstreamErrorName(tt.input)
/root/astrate-mule/internal/engine/triggers/errorname_test.go:36:				t.Errorf("UpstreamErrorName(%q) = %q, want %q", tt.input, got, tt.expected)
/root/astrate-mule/internal/engine/triggers/errorname_test.go:42:func TestUpstreamErrorNamePassthrough(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/errorname_test.go:43:	for _, name := range UpstreamErrorNames() {
/root/astrate-mule/internal/engine/triggers/errorname_test.go:45:			got := UpstreamErrorName(name)
/root/astrate-mule/internal/engine/triggers/errorname_test.go:47:				t.Errorf("UpstreamErrorName(%q) = %q, want %q (passthrough)", name, got, name)
/root/astrate-mule/internal/engine/triggers/errorname_test.go:53:func TestUpstreamErrorNameClosedSetInvariant(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/errorname_test.go:54:	fallback := "interface_loading_failed"
/root/astrate-mule/internal/engine/triggers/errorname_test.go:57:	for _, n := range UpstreamErrorNames() {
/root/astrate-mule/internal/engine/triggers/errorname_test.go:62:			got := UpstreamErrorName(input)
/root/astrate-mule/internal/engine/triggers/errorname_test.go:64:				t.Errorf("UpstreamErrorName(%q) = %q is not in UpstreamErrorNames()", input, got)
/root/astrate-mule/internal/engine/triggers/errorname_test.go:67:				t.Errorf("UpstreamErrorName(%q) = %q, want %q", input, got, fallback)
/root/astrate-mule/internal/engine/triggers/testdata/device_error.json:1:{"timestamp":"2026-06-12T10:00:00.123Z","device_id":"f0VMRgIBAQAAAAAAAAAAAA","event":{"type":"device_error","error_name":"interface_loading_failed","metadata":{"astrate_reason":"interface_not_in_introspection","detail":"no introspected interface matches x"}},"trigger_name":"example_trigger"}
/root/astrate-mule/internal/engine/triggers/events.go:127:// error_name enum by UpstreamErrorName, because consumers validate error_name
/root/astrate-mule/internal/engine/triggers/events.go:134:	upstream := UpstreamErrorName(errorName)
/root/astrate-mule/internal/appengine/channels/room_test.go:113:				// validates against (triggers.UpstreamErrorName).

→ Read internal/engine/triggers/errorname.go
→ Read test/conformance/upstream/README.md [offset=70, limit=160]
The upstream v1.3.3 change is directly on a surface Astrate implements. Let me find Astrate's properties-resend path to see how it emits `device_error` there.

$ rg -n -i "resend" /root/astrate-mule/internal --include="*.go" | rg -v "_test.go" | head -30
rg: unrecognized flag --include

similar flags that are available: --include-zero

$ rg -n "NewDeviceErrorEvent|device_error" /root/astrate-mule/internal --include="*.go" | rg -v "_test.go" | head -30
rg: unrecognized flag --include

similar flags that are available: --include-zero

$ rg -n -i "resend" /root/astrate-mule/internal -g "*.go" | rg -v "_test.go" | head -30
/root/astrate-mule/internal/engine/control.go:66:		if err := e.resendServerProperties(ctx, realm, m.DeviceID, format); err != nil {
/root/astrate-mule/internal/engine/control.go:138:// resendServerProperties publishes every server-owned property of the
/root/astrate-mule/internal/engine/control.go:143:func (e *Engine) resendServerProperties(ctx context.Context, realm *realmSchema, id deviceid.ID, format payload.Format) error {
/root/astrate-mule/internal/engine/control.go:152:			e.log.Warn("server-owned property references an uninstalled interface; skipping resend",
/root/astrate-mule/internal/engine/control.go:158:			e.log.Warn("server-owned property path matches no endpoint; skipping resend",
/root/astrate-mule/internal/engine/control.go:164:			e.log.Error("stored property value does not rehydrate; skipping resend",
/root/astrate-mule/internal/engine/control.go:170:			e.log.Error("property value does not encode; skipping resend",
/root/astrate-mule/internal/engine/triggers/errorname.go:20:	"resend_interface_properties_failed",
/root/astrate-mule/internal/engine/cache.go:60:	// (emptyCache resend + `consumer/properties` payload, docs/DESIGN.md §3.4).
/root/astrate-mule/internal/broker/sessionstore.go:209:// OnQosPublish upserts an inflight (or offline-queued) message. Resends

$ rg -n "NewDeviceErrorEvent|device_error" /root/astrate-mule/internal -g "*.go" | rg -v "_test.go" | head -30
/root/astrate-mule/internal/engine/data.go:18:// metric labels and the M6b device_error trigger error names.
/root/astrate-mule/internal/engine/data.go:375:// device_error seam (M6b), and acknowledgment — devices must not stall on
/root/astrate-mule/internal/engine/stream/bus.go:30:	KindDeviceError = "device_error"
/root/astrate-mule/internal/engine/stream/bus.go:49:	// ErrorName and ErrorMetadata describe a device_error event: the §2.6
/root/astrate-mule/internal/engine/engine.go:212:// becomes a device_error trigger event (docs/DESIGN.md §2.6 "failures are
/root/astrate-mule/internal/engine/engine.go:216:	// declares no trigger at all — a Channels viewer watching device_error must
/root/astrate-mule/internal/engine/engine.go:233:		triggers.NewDeviceErrorEvent(reason, map[string]string{"detail": detail}))
/root/astrate-mule/internal/engine/triggers/errorname.go:26:// — NewDeviceErrorEvent keeps it under metadata["astrate_reason"].
/root/astrate-mule/internal/engine/router.go:115:	// onDeviceError observes every rejection (M6b wires device_error trigger
/root/astrate-mule/internal/engine/triggers/match.go:53:	OnDeviceError = "device_error"
/root/astrate-mule/internal/appengine/channels/room.go:224:				event = triggers.NewDeviceErrorEvent(ev.ErrorName, ev.ErrorMetadata)
/root/astrate-mule/internal/engine/triggers/events.go:114:// DeviceErrorEvent is the device_error event body.
/root/astrate-mule/internal/engine/triggers/events.go:116:	// Type is always "device_error".
/root/astrate-mule/internal/engine/triggers/events.go:125:// NewDeviceErrorEvent builds a device_error event body from one of Astrate's
/root/astrate-mule/internal/engine/triggers/events.go:133:func NewDeviceErrorEvent(errorName string, metadata map[string]string) DeviceErrorEvent {

→ Read internal/engine/control.go [limit=200]
Let me get more context on the upstream change — which functions the hunks belong to.

$ gh api repos/astarte-platform/astarte/contents/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/device.ex?ref=v1.3.3 --jq '.content' | base64 -d > /tmp/upstream_device.ex; rg -n "def resend|defp resend|def update_trigger|def handle_message|interface_loading_failed|sending_properties_to_interface_failed" /tmp/upstream_device.ex | head -30
416:  def resend_all_properties(state) do
437:  defp resend_all_interface_properties(
463:          {:halt, {:error, :sending_properties_to_interface_failed}}
468:  defp resend_all_interface_properties(_state, %InterfaceDescriptor{} = _descriptor) do
472:  defp resend_interface_mapping_properties(

$ sed -n '400,530p' /tmp/upstream_device.ex
{:error, reason}
    end
  end

  def prune_device_properties(state, decoded_payload, timestamp) do
    with {:ok, paths_set} <-
           PayloadsDecoder.parse_device_properties_payload(decoded_payload, state.introspection) do
      Enum.each(state.introspection, fn {interface, _} ->
        # TODO: check result here
        Core.Interface.prune_interface(state, interface, paths_set, timestamp)
      end)

      :ok
    end
  end

  def resend_all_properties(state) do
    Logger.debug("Device introspection: #{inspect(state.introspection)}")

    Enum.reduce_while(state.introspection, {:ok, state}, fn {interface, _}, {:ok, state_acc} ->
      maybe_descriptor = Map.get(state_acc.interfaces, interface)

      with {:ok, interface_descriptor, new_state} <-
             Core.Interface.maybe_handle_cache_miss(maybe_descriptor, interface, state_acc),
           :ok <- resend_all_interface_properties(new_state, interface_descriptor) do
        {:cont, {:ok, new_state}}
      else
        {:error, reason} = error ->
          Logger.warning(
            "Error while sending properties to device: #{inspect(reason)} on interface #{interface}"
          )

          {:halt, error}
      end
    end)
  end

  defp resend_all_interface_properties(
         %State{realm: realm, device_id: device_id, mappings: mappings} = _state,
         %InterfaceDescriptor{type: :properties, ownership: :server} = interface_descriptor
       ) do
    encoded_device_id = Device.encode_device_id(device_id)

    mappings
    |> Enum.filter(fn {_endpoint_id, mapping} ->
      mapping.interface_id == interface_descriptor.interface_id
    end)
    |> Enum.reduce_while(:ok, fn {_endpoint_id, mapping}, :ok ->
      case resend_interface_mapping_properties(
             realm,
             device_id,
             encoded_device_id,
             interface_descriptor,
             mapping
           ) do
        :ok ->
          {:cont, :ok}

        {:error, reason} ->
          Logger.warning(
            "Failed sending properties for interface #{interface_descriptor.name}: #{inspect(reason)}."
          )

          {:halt, {:error, :sending_properties_to_interface_failed}}
      end
    end)
  end

  defp resend_all_interface_properties(_state, %InterfaceDescriptor{} = _descriptor) do
    :ok
  end

  defp resend_interface_mapping_properties(
         realm,
         device_id,
         encoded_device_id,
         interface_descriptor,
         %Mapping{value_type: value_type} = mapping
       ) do
    column_name =
      CQLUtils.type_to_db_column_name(value_type) |> String.to_existing_atom()

    Queries.retrieve_property_values(realm, device_id, interface_descriptor, mapping)
    |> Enum.reduce_while(:ok, fn %{:path => path, ^column_name => value}, _acc ->
      handle_property_send(
        realm,
        encoded_device_id,
        interface_descriptor.name,
        path,
        value_type,
        value
      )
    end)
  end

  defp handle_property_send(realm, encoded_device_id, interface_name, path, value_type, value) do
    case send_value(
           realm,
           encoded_device_id,
           interface_name,
           path,
           value_type,
           value
         ) do
      {:ok, _bytes} ->
        # TODO: use the returned bytes count in stats
        {:cont, :ok}

      {:error, reason} ->
        {:halt, {:error, reason}}
    end
  end

  defp send_value(realm, device_id_string, interface_name, path, value_type, value) do
    topic = "#{realm}/#{device_id_string}/#{interface_name}#{path}"
    encapsulated_value = %{v: cast_bson_value(value_type, value)}

    bson_value = Cyanide.encode!(encapsulated_value)

    Logger.debug("Going to publish #{inspect(encapsulated_value)} on #{topic}.")

    case VMQPlugin.publish(topic, bson_value, 2) do
      {:ok, %{local_matches: local, remote_matches: remote}} when local + remote == 1 ->
        {:ok, byte_size(topic) + byte_size(bson_value)}

      {:ok, %{local_matches: local, remote_matches: remote}} when local + remote > 1 ->
        # This should not happen so we print a warning, but we consider it a successful publish
        Logger.warning(
          "Multiple match while publishing #{inspect(encapsulated_value)} on #{topic}.",
          tag: "publish_multiple_matches"
        )

The v1.3.3 fix is on the emptyCache property-resend path. Let me confirm what events upstream emits there now and what the old `sending_properties_error` did.
$ gh api repos/astarte-platform/astarte/contents/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex?ref=v1.3.3 --jq '.content' | base64 -d > /tmp/upstream_control_handler.ex; rg -n "interface_loading_error|sending_properties_error|execute_device_error_triggers|resend_interface_properties_failed|interface_loading_failed|resend_all_properties|resend_start" /tmp/upstream_control_handler.ex
136:         {:ok, state} <- resend_all_properties(state, message_id, timestamp),
179:    Core.Trigger.execute_device_error_triggers(
218:          {:error, :interface_loading_failed} ->
332:  defp resend_all_properties(state, message_id, timestamp) do
333:    resend_start = System.monotonic_time()
335:    case Core.Device.resend_all_properties(state) do
340:          %{duration: System.monotonic_time() - resend_start},
346:      {:error, :interface_loading_failed} ->
350:          %{duration: System.monotonic_time() - resend_start},
351:          %{realm: state.realm, result: "interface_loading_failed"}
354:        interface_loading_error(state, message_id, timestamp)
360:          %{duration: System.monotonic_time() - resend_start},
364:        sending_properties_error(state, message_id, timestamp)
370:          %{duration: System.monotonic_time() - resend_start},
390:    Core.Trigger.execute_device_error_triggers(
399:  defp interface_loading_error(state, message_id, timestamp) do
401:      tag: "interface_loading_failed"
413:    Core.Trigger.execute_device_error_triggers(
415:      "interface_loading_failed",
422:  defp sending_properties_error(state, message_id, timestamp) do
424:      tag: "resend_interface_properties_failed"
436:    Core.Trigger.execute_device_error_triggers(
438:      "resend_interface_properties_failed",
461:    Core.Trigger.execute_device_error_triggers(

$ sed -n '100,230p' /tmp/upstream_control_handler.ex
total_received_bytes:
            new_state.total_received_bytes + byte_size(payload) +
              byte_size("/producer/properties")
      }
    else
      _ ->
        # Track failed decompression
        :telemetry.execute(
          [:astarte, :data_updater_plant, :control_handler, :payload_decompression],
          %{
            duration: System.monotonic_time() - decompression_start,
            compressed_size: byte_size(payload),
            uncompressed_size: 0
          },
          %{realm: new_state.realm, result: "failed"}
        )

        Logger.warning("Invalid purge_properties payload", tag: "purge_properties_error")

        {:ok, new_state} = Core.Device.ask_clean_session(new_state, timestamp)
        MessageTracker.discard(new_state.message_tracker, message_id)

        :telemetry.execute(
          [:astarte, :data_updater_plant, :data_updater, :discarded_message],
          %{},
          %{realm: new_state.realm}
        )

        new_state
    end
  end

  def handle_control(state, "/emptyCache", _payload, message_id, timestamp) do
    state = TimeBasedActions.execute_time_based_actions(state, timestamp)

    with :ok <- send_control_consumer_properties(state, message_id, timestamp),
         {:ok, state} <- resend_all_properties(state, message_id, timestamp),
         :ok <- set_pending_empty_cache(state, message_id, timestamp) do
      MessageTracker.ack_delivery(state.message_tracker, message_id)

      :telemetry.execute(
        [:astarte, :data_updater_plant, :data_updater, :processed_empty_cache],
        %{},
        %{realm: state.realm}
      )

      state
    end
  end

  def handle_control(state, path, payload, message_id, timestamp) do
    # Track unexpected control messages
    :telemetry.execute(
      [:astarte, :data_updater_plant, :control_handler, :unexpected_control],
      %{payload_size: byte_size(payload)},
      %{realm: state.realm, control_path: path}
    )

    Logger.warning(
      "Unexpected control on #{path}, base64-encoded payload: #{inspect(Base.encode64(payload))}",
      tag: "unexpected_control_message"
    )

    {:ok, new_state} = Core.Device.ask_clean_session(state, timestamp)
    MessageTracker.discard(new_state.message_tracker, message_id)

    :telemetry.execute(
      [:astarte, :data_updater_plant, :data_updater, :discarded_control_message],
      %{},
      %{realm: new_state.realm}
    )

    base64_payload = Base.encode64(payload)

    error_metadata = %{
      "path" => inspect(path),
      "base64_payload" => base64_payload
    }

    Core.Trigger.execute_device_error_triggers(
      new_state,
      "unexpected_control_message",
      error_metadata,
      timestamp
    )

    Core.DataHandler.update_stats(new_state, "", nil, path, payload)
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

  defp send_control_consumer_properties(state, message_id, timestamp) do
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

Confirmed the full picture. Now I have everything I need to write the proposal.

**Findings:**
- Newest tag is **v1.3.3** (2026-08-07), newer than the last-reviewed v1.4.0-rc.3. v1.4.0-rc.4 (2026-08-06) is also newer.
- **v1.4.0-rc.4** is all Elixir-shaped infra (config/env-var library, vault, FDO rendezvous, AMQP consumer memory bound, DUP RPC reliability, GC) — nothing Astrate implements.
- **v1.3.3** contains exactly one real change, `fix(dup): stop remapping interface loading failed error` (#2119): on `emptyCache`, when the server-owned-property resend fails because an interface can't be loaded, upstream now emits a `device_error` trigger event with `error_name: "interface_loading_failed"` (previously remapped to `resend_interface_properties_failed`). It also makes the per-mapping resend halt-on-first-error instead of silently swallowing.
- Astrate's counterpart `resendServerProperties` (`internal/engine/control.go:143`) only logs-and-skips on unloadable interface/path/encode failures and emits **no** `device_error` event at all on this path — a wire-visible gap vs upstream's trigger payloads, but which Astrate failure maps to which upstream name is unmeasured (the recorded `channels.json` fixture never covered the resend path). That makes a probe the strict first step.

Appending the proposal lines to the queue:
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -90,4 +90,6 @@
 - [ ] realm-config-datastream-retention: add `GET /realmmanagement/v1/<realm>/config/datastream_maximum_storage_retention` to `internal/realm/http.go` and carry `datastream_maximum_storage_retention` through `internal/housekeeping/http.go` + `internal/store/realms.go` (upstream v1.3.0 per-realm retention, v1.4.0-rc.1 adds an instance-wide default); detail in `.mule/tasks/realm-config-datastream-retention.md` [auto]
 - [ ] probe-interface-default-values: does `GET /realmmanagement/v1/<realm>/interfaces/<name>/<major>` return the same defaulted mapping-parameter values upstream reports after 3f0b864 ("Properly show default values for all mapping parameters")? report, do not patch [auto]
 - [ ] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto]
 - [ ] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto]
+- [ ] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto]
+- [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]

→ Read .mule/todo.md [offset=88]
```
Done: reviewed upstream against v1.4.0-rc.3 (last checked) — v1.3.3 (2026-08-07) and v1.4.0-rc.4 (2026-08-06) are newer; proposed 2 queue lines
Files: .mule/todo.md
Verified: gh api releases/compare reads (no code gates run — no code changed)
Unsure: whether upstream's device_error-on-resend delivery is reliable on this path (unmeasured — the recorded fixture only covered publish-path provocations); which Astrate resend failure maps to which upstream error_name is a guess, which is why the probe precedes any patch
Follow-ups: probe-emptycache-resend-device-error, compat-note-v1.3.3 (both tagged [auto])
```

Details: v1.4.0-rc.4's seven commits are all Elixir-shaped infra (config/env-var library, vault, FDO rendezvous, AMQP consumer memory bound, DUP RPC reliability, GC) — nothing Astrate implements, nothing proposed. v1.3.3's only change, `fix(dup): stop remapping interface loading failed error` (#2119), is wire-visible: on `emptyCache` the server-owned-property resend now emits a `device_error` trigger event with `error_name: interface_loading_failed` where it previously remapped to `resend_interface_properties_failed`, and per-mapping failures now halt and emit `resend_interface_properties_failed` instead of being swallowed. Astrate's `resendServerProperties` (internal/engine/control.go:143) logs-and-skips on unloadable-interface/path/encode failures and never emits a `device_error` on this path.
