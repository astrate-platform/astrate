slug: probe-props-resend-error-triggers
verdict: blocked
at:  3e1b801
ran: 2026-09-06T18:39:52Z on DietPi in 331s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/engine/control.go [offset=130, limit=100]
This is a `[readonly]` investigation — read the resend-failure paths in Astrate, compare with upstream v1.3.3's two distinct device_error triggers, and report. No code changes.

Let me read the relevant source files.
→ Read internal/engine/engine.go [offset=345, limit=40]
✱ Grep "resendServerProperties" in internal/engine · 3 matches
✱ Grep "device_error" in internal · 29 matches
Let me also look at how `resendServerProperties` is called on the connect path, and check for any existing device_error trigger firing.
✱ Grep "DeviceEvent|fireDevice" in internal/engine · 46 matches
→ Read internal/engine/engine.go [limit=40, offset=295]
→ Read internal/engine/triggers/events.go [limit=40, offset=235]
→ Read internal/engine/control.go [limit=120, offset=50]
✱ Grep "UpstreamErrorName" in internal/engine/triggers · 18 matches
Now let me look at the upstream error name mapping and the reject reasons used by Astrate.
✱ Grep "reasonControl|reason[A-Z]" in internal/engine · 70 matches
→ Read internal/engine/triggers/errorname.go
→ Read internal/engine/router.go [limit=40, offset=110]
Now let me see the `reject` function and how it calls `onDeviceError`.
→ Read internal/engine/data.go [limit=60, offset=490]
Let me see the full connect path that calls `resendServerProperties`.
→ Read internal/engine/data.go [limit=80, offset=350]
✱ Grep "resendServerProperties|resend.*propert" in internal/engine · 5 matches
→ Read internal/engine/data.go [limit=50, offset=1]
→ Read internal/engine/cache.go [limit=80, offset=50]
✱ Grep "resend_interface_properties_failed|interface_loading_failed" in . · 36 matches
Let me check how upstream Astarte handles these error names — first look for vendored/parity references locally, then check the actual upstream commit.
✱ Grep "1\.3\.3|2119" in . · 66 matches
The prior upstream-watch report already analyzed this fix. Let me read the relevant section to verify against the actual upstream commit.
→ Read .mule/reports/log/20260906T183102Z-recipe-astarte-upstream-proposed.md [limit=120, offset=320]
→ Read .mule/reports/log/20260906T183102Z-recipe-astarte-upstream-proposed.md [limit=80, offset=440]
The prior report captured the upstream diff well. Let me verify against the actual upstream commit to be sure.

$ gh api repos/astarte-platform/astarte/commits/d8592091e1 --jq '.commit.message, .files[] | .filename' 2>&1 | head -30
expected an object but got: string ("fix(dup): stop remapping ...")

$ gh api repos/astarte-platform/astarte/commits/d8592091e1 --jq '.commit.message' 2>&1; echo "---FILES---"; gh api repos/astarte-platform/astarte/commits/d8592091e1 --jq '.files[] | .filename' 2>&1
fix(dup): stop remapping interface loading failed error (#2119)

Signed-off-by: Riccardo Nalgi <riccardo.nalgi@secomind.com>
---FILES---
CHANGELOG.md
apps/astarte_appengine_api/mix.lock
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/device.ex
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/control_handler_test.exs
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/device_test.exs
apps/astarte_pairing/mix.lock

$ gh api repos/astarte-platform/astarte/commits/d8592091e1 --jq '.files[] | select(.filename | contains("device.ex")) | .patch' 2>&1
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

Let me see the full upstream `resend_interface_properties` and `resend_interface_mapping_properties` functions in v1.3.3 to identify where each error originates.
$ gh api "repos/astarte-platform/astarte/contents/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/device.ex?ref=v1.3.3" --jq '.content' 2>/dev/null | tr -d '\n' | base64 -d 2>/dev/null | rg -n "resend_interface_properties|resend_all_interface_properties|resend_interface_mapping_properties|interface_loading_failed"
424:           :ok <- resend_all_interface_properties(new_state, interface_descriptor) do
437:  defp resend_all_interface_properties(
448:      case resend_interface_mapping_properties(
468:  defp resend_all_interface_properties(_state, %InterfaceDescriptor{} = _descriptor) do
472:  defp resend_interface_mapping_properties(

$ gh api "repos/astarte-platform/astarte/contents/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/device.ex?ref=v1.3.3" --jq '.content' 2>/dev/null | tr -d '\n' | base64 -d 2>/dev/null | sed -n '400,540p'
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

        {:ok, byte_size(topic) + byte_size(bson_value)}

      {:ok, %{local_matches: local, remote_matches: remote}} when local + remote == 0 ->
        {:error, :session_not_found}

      {:error, reason} ->
        {:error, reason}
    end
  end

Now let me confirm where `:interface_loading_failed` originates — `Core.Interface.maybe_handle_cache_miss`.
$ gh api "repos/astarte-platform/astarte/contents/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex?ref=v1.3.3" --jq '.content' 2>/dev/null | tr -d '\n' | base64 -d 2>/dev/null | rg -n -A15 "maybe_handle_cache_miss|interface_loading_failed"
43:  def maybe_handle_cache_miss(nil, interface_name, state) do
44-    with {:ok, major_version} <-
45-           Device.interface_version(state.realm, state.device_id, interface_name),
46-         {:ok, interface_row} <-
47-           Interface.retrieve_interface_row(state.realm, interface_name, major_version),
48-         %InterfaceDescriptor{interface_id: interface_id} = interface_descriptor <-
49-           InterfaceDescriptor.from_db_result!(interface_row),
50-         {:ok, mappings} <-
51-           Mappings.fetch_interface_mappings_map(state.realm, interface_id) do
52-      new_interfaces_by_expiry =
53-        state.interfaces_by_expiry ++
54-          [{state.last_seen_message + @interface_lifespan_decimicroseconds, interface_name}]
55-
56-      new_state = %State{
57-        state
58-        | interfaces: Map.put(state.interfaces, interface_name, interface_descriptor),
--
74:        {:error, :interface_loading_failed}
75-
76-      {:error, :device_not_found} ->
77:        {:error, :interface_loading_failed}
78-
79-      {:error, :database_error} ->
80:        {:error, :interface_loading_failed}
81-
82-      {:error, :interface_not_found} ->
83:        {:error, :interface_loading_failed}
84-
85-      other ->
86:        Logger.warning("maybe_handle_cache_miss failed: #{inspect(other)}")
87:        {:error, :interface_loading_failed}
88-    end
89-  end
90-
91:  def maybe_handle_cache_miss(interface_descriptor, _interface_name, state) do
92-    {:ok, interface_descriptor, state}
93-  end
94-
95-  def prune_interface(state, interface, all_paths_set, timestamp) do
96-    with {:ok, interface_descriptor, new_state} <-
97:           maybe_handle_cache_miss(
98-             Map.get(state.interfaces, interface),
99-             interface,
100-             state
101-           ) do
102-      cond do
103-        interface_descriptor.type != :properties ->
104-          # TODO: nobody uses new_state
105-          {:ok, new_state}
106-
107-        interface_descriptor.ownership != :device ->
108-          Logger.warning("Tried to prune server owned interface: #{interface}.")
109-          {:error, :maybe_outdated_introspection}
110-
111-        true ->
112-          do_prune(new_state, interface_descriptor, all_paths_set, timestamp)

Let me check the full v1.3.3 control_handler emptyCache path to see the error handling context.
$ gh api "repos/astarte-platform/astarte/contents/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex?ref=v1.3.3" --jq '.content' 2>/dev/null | tr -d '\n' | base64 -d 2>/dev/null | sed -n '300,430p'
[:astarte, :data_updater_plant, :control_handler, :vmq_publish],
          %{
            duration: System.monotonic_time() - publish_start,
            payload_size: byte_size(payload)
          },
          %{realm: realm, result: "no_matches", matches: "none"}
        )

        {:error, :session_not_found}

      {:error, reason} ->
        # Track publish errors
        :telemetry.execute(
          [:astarte, :data_updater_plant, :control_handler, :vmq_publish],
          %{
            duration: System.monotonic_time() - publish_start,
            payload_size: byte_size(payload)
          },
          %{realm: realm, result: "error", matches: "unknown"}
        )

        {:error, reason}
    end
  end

  defp set_pending_empty_cache(state, message_id, timestamp) do
    case Queries.set_pending_empty_cache(state.realm, state.device_id, false) do
      :ok -> :ok
      {:error, reason} -> generic_error(state, message_id, timestamp, reason)
    end
  end

  defp resend_all_properties(state, message_id, timestamp) do
    resend_start = System.monotonic_time()

    case Core.Device.resend_all_properties(state) do
      {:ok, new_state} ->
        # Track successful properties resend
        :telemetry.execute(
          [:astarte, :data_updater_plant, :control_handler, :properties_resend],
          %{duration: System.monotonic_time() - resend_start},
          %{realm: state.realm, result: "success"}
        )

        {:ok, new_state}

      {:error, :interface_loading_failed} ->
        # Track interface loading failure
        :telemetry.execute(
          [:astarte, :data_updater_plant, :control_handler, :properties_resend],
          %{duration: System.monotonic_time() - resend_start},
          %{realm: state.realm, result: "interface_loading_failed"}
        )

        interface_loading_error(state, message_id, timestamp)

      {:error, :sending_properties_to_interface_failed} ->
        # Track interface send failure
        :telemetry.execute(
          [:astarte, :data_updater_plant, :control_handler, :properties_resend],
          %{duration: System.monotonic_time() - resend_start},
          %{realm: state.realm, result: "sending_properties_to_interface_failed"}
        )

        sending_properties_error(state, message_id, timestamp)

      {:error, reason} ->
        # Track other resend failures
        :telemetry.execute(
          [:astarte, :data_updater_plant, :control_handler, :properties_resend],
          %{duration: System.monotonic_time() - resend_start},
          %{realm: state.realm, result: "error"}
        )

        generic_error(state, message_id, timestamp, reason)
    end
  end

  defp session_not_found_error(state, message_id, timestamp) do
    Logger.warning("Cannot push data to device.", tag: "device_session_not_found")

    {:ok, new_state} = Core.Device.ask_clean_session(state, timestamp)
    MessageTracker.discard(new_state.message_tracker, message_id)

    :telemetry.execute(
      [:astarte, :data_updater_plant, :data_updater, :discarded_message],
      %{},
      %{realm: new_state.realm}
    )

    Core.Trigger.execute_device_error_triggers(
      new_state,
      "device_session_not_found",
      timestamp
    )

    new_state
  end

  defp interface_loading_error(state, message_id, timestamp) do
    Logger.warning("Cannot load interface while resending properties",
      tag: "interface_loading_failed"
    )

    {:ok, new_state} = Core.Device.ask_clean_session(state, timestamp)
    MessageTracker.discard(new_state.message_tracker, message_id)

    :telemetry.execute(
      [:astarte, :data_updater_plant, :data_updater, :discarded_message],
      %{},
      %{realm: new_state.realm}
    )

    Core.Trigger.execute_device_error_triggers(
      new_state,
      "interface_loading_failed",
      timestamp
    )

    new_state
  end

  defp sending_properties_error(state, message_id, timestamp) do
    Logger.warning("Cannot send properties to device.",
      tag: "resend_interface_properties_failed"
    )

    {:ok, new_state} = Core.Device.ask_clean_session(state, timestamp)
    MessageTracker.discard(new_state.message_tracker, message_id)

    :telemetry.execute(

Now let me check Astrate's `retryStore` to fully describe what happens on store/publish errors in the emptyCache path.
✱ Grep "func \(e \*Engine\) retryStore|consumerPropertiesSendTimeout" in internal/engine · 4 matches
→ Read internal/engine/introspection.go [limit=30, offset=130]
→ Read internal/engine/engine.go [limit=15, offset=45]
Let me check Astrate's broker `Publish` semantics — whether a publish to a topic with no live subscriber returns an error or is silently dropped.
✱ Grep "func.*Publish|Publish\(" in internal/broker · 23 matches
→ Read internal/broker/publisher.go [limit=30, offset=30]
→ Read internal/broker/publisher.go [limit=10, offset=58]
Let me check what upstream `generic_error` fires, to complete the DB-error case mapping.

$ gh api "repos/astarte-platform/astarte/contents/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex?ref=v1.3.3" --jq '.content' 2>/dev/null | tr -d '\n' | base64 -d 2>/dev/null | rg -n -A20 "defp generic_error"
445:  defp generic_error(state, message_id, timestamp, reason) do
446-    Logger.warning("Unhandled error during emptyCache: #{inspect(reason)}",
447-      tag: "empty_cache_error"
448-    )
449-
450-    {:ok, new_state} = Core.Device.ask_clean_session(state, timestamp)
451-    MessageTracker.discard(new_state.message_tracker, message_id)
452-
453-    :telemetry.execute(
454-      [:astarte, :data_updater_plant, :data_updater, :discarded_message],
455-      %{},
456-      %{realm: new_state.realm}
457-    )
458-
459-    error_metadata = %{"reason" => inspect(reason)}
460-
461-    Core.Trigger.execute_device_error_triggers(
462-      new_state,
463-      "empty_cache_error",
464-      error_metadata,
465-      timestamp

✱ Grep "sendConsumerProperties" in internal/engine · 5 matches
→ Read internal/engine/serverdata.go [limit=30, offset=235]
Let me verify the no-subscriber publish behavior against Astrate's broker tests.
→ Read internal/broker/broker_test.go [limit=30, offset=305]
Investigation complete. Here is the report.

## Findings

**Upstream v1.3.3 fix (#2119, `d8592091`).** Verified against the commit and the v1.3.3 tree. In upstream, property resend happens in a single batch (`resend_all_properties`): for each introspected interface it (a) *loads* the compiled interface descriptor (`Interface.maybe_handle_cache_miss` — any failure there is folded to `:interface_loading_failed`) and then (b) *publishes* each stored server-owned property (`send_value` — VMQ publish error or `local + remote == 0`, i.e. **no reachable device session**, folds to `:sending_properties_to_interface_failed`). Before 1.3.3 both were remapped to `:sending_properties_to_interface_failed`, so both fired `resend_interface_properties_failed`. The fix propagates the original error and the control handler fires `interface_loading_failed` for (a), `resend_interface_properties_failed` for (b), and generic `empty_cache_error` for anything else — each once per batch, after discarding the control message.

**Mapping of Astrate's `resendServerProperties` modes (control.go:144):**
- `ci == nil` (uninstalled interface, control.go:152), path not in trie (control.go:157) → closest to **`interface_loading_failed`** (upstream has no analog for the trie miss: it iterates mappings, so a stale property row is never queried and silently skipped — no trigger).
- rehydrate failure (control.go:163), encode failure (control.go:169) → a value that cannot be sent → **`resend_interface_properties_failed`**.
- `Publish` error (control.go:176) → genuine send-to-device failure → **`resend_interface_properties_failed`**.
- store read error (control.go:145) → propagated into `retryStore` (park/retry); upstream's analog (DB error) fires **`empty_cache_error`**.

**`sendConsumerProperties` (control.go:191)** failures (store read, deflate build, `Publish`) map to `resend_interface_properties_failed`: on emptyCache they retry-store; on connect they log Warn (engine.go:359-362); on server-owned unset (serverdata.go:250) they propagate to the API caller.

**Should Astrate fire a device_error?** Three structural facts matter:
1. Both names are already in Astrate's closed upstream set (errorname.go:11,20) and `interface_loading_failed` is reachable — but `resend_interface_properties_failed` is emitted **nowhere** today.
2. The device_error seam (`fireDeviceError`, engine.go:305) only fires from `reject()`, i.e. for inbound messages. Only the emptyCache resend has an inbound message; the connect and unset purge paths don't, so firing there needs a new message-less seam.
3. The case 1.3.3 is most eager to surface — a device with no reachable session — is **undetectable** via Astrate's `Publisher.Publish` (publisher.go:37 wraps mochi `InjectPacket`, which succeeds when zero subscribers match; QoS 2 queues or drops silently). So `Publish` errors would only ever mean broker-internal failures, a far narrower and less useful event than upstream's name implies. Wire-faithful `resend_interface_properties_failed` requires an explicit session-liveness check Astrate does not have.

Recommendation: if matching 1.3.3 wire behaviour for Channels/Dashboard consumers is desired, the only faithful point is a batch-level decision on the emptyCache resend (interface-resolution skips → `interface_loading_failed`; rehydrate/encode/publish failures → `resend_interface_properties_failed`). But Astrate would need session-liveness detection first, or `resend_interface_properties_failed` would fire only for rare broker plumbing errors. Today the store/publish errors are already loud via retryStore park/retry + attacked logs — not silent.

---

Done: readonly probe — mapped Astrate's resend-failure modes to upstream v1.3.3's split device_error names and assessed trigger worthiness; no patch
Files: internal/engine/control.go, internal/engine/engine.go, internal/engine/triggers/errorname.go (read-only; nothing changed)
Verified: upstream commit d8592091 + v1.3.3 tree fetched via gh api -> pass
Unsure: no claim on whether Giulio wants 1.3.3 wire parity for device_error on resend — that is his call
Follow-ups: none
