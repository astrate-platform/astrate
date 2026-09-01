slug: probe-object-validation-selected-interface
verdict: blocked
at:  4cd1e88
ran: 2026-09-01T18:13:03Z on DietPi in 201s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
Let me understand this task. It's a `[auto]` probe task - read-only, report findings. I need to check whether Astrate decodes and validates an object datastream against the topic-named interface's own mappings (specifically `ci.ObjectLeaves` in internal/engine/data.go), analogous to the upstream v1.4.0-rc.5 fix #2141.

Let me read the relevant code.
→ Read internal/engine/data.go
The key line is at data.go:269:
```go
dp, err = e.dec.Object(m.Payload, ci.ObjectLeaves)
```

Here `ci` is the compiled interface from the topic that was matched (via `matchInterface` and `realm.iface(name, ver.Major)`). So the object datastream is decoded against `ci.ObjectLeaves`, where `ci` is the **topic-named interface's own** compiled interface. This is exactly what the task asks — the decode uses the selected interface's own mappings.

Let me look at the upstream fix #2141 context. Let me understand what `e.dec.Object` does and how `ObjectLeaves` is used. Let me look at the payload decoder and the CompiledInterface structure.

Let me examine the payload `Object` function and `ObjectLeaves`.

$ rg -n "func.*Object|ObjectLeaves|func.*dec.Object|DecodedPayload" /root/astrate-mule/pkg/payload/ /root/astrate-mule/pkg/interfaceschema/ | head -50
/root/astrate-mule/pkg/interfaceschema/compile.go:38:	// ObjectLeaves maps, for object aggregation, each last-level name to its
/root/astrate-mule/pkg/interfaceschema/compile.go:40:	ObjectLeaves map[string]*CompiledMapping
/root/astrate-mule/pkg/interfaceschema/compile.go:68:		ci.ObjectLeaves = make(map[string]*CompiledMapping, len(iface.Mappings))
/root/astrate-mule/pkg/interfaceschema/compile.go:92:		if ci.ObjectLeaves != nil {
/root/astrate-mule/pkg/interfaceschema/compile.go:97:			if _, dup := ci.ObjectLeaves[leaf]; dup {
/root/astrate-mule/pkg/interfaceschema/compile.go:100:			ci.ObjectLeaves[leaf] = cm
/root/astrate-mule/pkg/interfaceschema/compile_test.go:72:	if ci.ObjectLeaves != nil {
/root/astrate-mule/pkg/interfaceschema/compile_test.go:73:		t.Errorf("ObjectLeaves = %v, want nil for individual aggregation", ci.ObjectLeaves)
/root/astrate-mule/pkg/interfaceschema/compile_test.go:100:func TestCompileObjectAggregation(t *testing.T) {
/root/astrate-mule/pkg/interfaceschema/compile_test.go:110:	if len(ci.ObjectLeaves) != len(wantLeaves) {
/root/astrate-mule/pkg/interfaceschema/compile_test.go:111:		t.Fatalf("len(ObjectLeaves) = %d, want %d", len(ci.ObjectLeaves), len(wantLeaves))
/root/astrate-mule/pkg/interfaceschema/compile_test.go:114:		m, ok := ci.ObjectLeaves[leaf]
/root/astrate-mule/pkg/interfaceschema/compile_test.go:116:			t.Errorf("ObjectLeaves missing %q", leaf)
/root/astrate-mule/pkg/interfaceschema/compile_test.go:128:	if m != ci.ObjectLeaves["latitude"] {
/root/astrate-mule/pkg/interfaceschema/compile_test.go:129:		t.Error("trie leaf and ObjectLeaves entry are not the same mapping")
/root/astrate-mule/pkg/interfaceschema/parse_test.go:234:func TestParseRejectsNonObject(t *testing.T) {
/root/astrate-mule/pkg/payload/value.go:73:// DecodedPayload is the result of decoding one data payload.
/root/astrate-mule/pkg/payload/value.go:74:type DecodedPayload struct {
/root/astrate-mule/pkg/payload/value.go:88:func (d DecodedPayload) IsUnset() bool { return d.Value == nil }
/root/astrate-mule/pkg/interfaceschema/parse.go:436:func checkObjectAggregation(iface *Interface, endpoints [][]endpointSegment) error {
/root/astrate-mule/pkg/interfaceschema/parse.go:463:func sameObjectAttributes(ref, m Mapping) error {
/root/astrate-mule/pkg/payload/payload.go:53:func (d Decoder) Individual(p []byte, m *interfaceschema.CompiledMapping) (DecodedPayload, error) {
/root/astrate-mule/pkg/payload/payload.go:55:		return DecodedPayload{}, errors.New("payload: nil mapping")
/root/astrate-mule/pkg/payload/payload.go:59:		return DecodedPayload{}, err
/root/astrate-mule/pkg/payload/payload.go:69:			return DecodedPayload{}, rejectf(ReasonUnsetNotAllowed, "empty payload on a mapping without allow_unset")
/root/astrate-mule/pkg/payload/payload.go:71:		return DecodedPayload{Format: FormatEmpty}, nil
/root/astrate-mule/pkg/payload/payload.go:75:			return DecodedPayload{}, err
/root/astrate-mule/pkg/payload/payload.go:78:			return DecodedPayload{}, err
/root/astrate-mule/pkg/payload/payload.go:84:			return DecodedPayload{}, err
/root/astrate-mule/pkg/payload/payload.go:87:			return DecodedPayload{}, err
/root/astrate-mule/pkg/payload/payload.go:96:// (CompiledInterface.ObjectLeaves). Object aggregation exists only on
/root/astrate-mule/pkg/payload/payload.go:100:func (d Decoder) Object(p []byte, leaves map[string]*interfaceschema.CompiledMapping) (DecodedPayload, error) {
/root/astrate-mule/pkg/payload/payload.go:102:		return DecodedPayload{}, errors.New("payload: no object leaves")
/root/astrate-mule/pkg/payload/payload.go:106:		return DecodedPayload{}, err
/root/astrate-mule/pkg/payload/payload.go:115:		return DecodedPayload{}, rejectf(ReasonUnsetNotAllowed, "empty payload on an object-aggregated datastream")
/root/astrate-mule/pkg/payload/payload.go:119:			return DecodedPayload{}, err
/root/astrate-mule/pkg/payload/payload.go:122:			return DecodedPayload{}, err
/root/astrate-mule/pkg/payload/payload.go:128:			return DecodedPayload{}, err
/root/astrate-mule/pkg/payload/payload.go:131:			return DecodedPayload{}, err
/root/astrate-mule/pkg/payload/payload.go:140:func finishDecode(val Value, ts *time.Time, f Format, explicit bool) (DecodedPayload, error) {
/root/astrate-mule/pkg/payload/payload.go:142:		return DecodedPayload{}, rejectf(ReasonBadTimestamp, "mapping declares explicit_timestamp but payload carries no t")
/root/astrate-mule/pkg/payload/payload.go:147:	return DecodedPayload{Value: val, Timestamp: ts, Format: f}, nil
/root/astrate-mule/pkg/payload/payload.go:163:func Decode(p []byte, m *interfaceschema.CompiledMapping) (DecodedPayload, error) {
/root/astrate-mule/pkg/payload/payload.go:169:func DecodeObject(p []byte, leaves map[string]*interfaceschema.CompiledMapping) (DecodedPayload, error) {
/root/astrate-mule/pkg/payload/json.go:201:func decodeJSONObject(raw json.RawMessage, leaves map[string]*interfaceschema.CompiledMapping) (Value, error) {
/root/astrate-mule/pkg/payload/payload_test.go:207:func reencodeAndCompare(t *testing.T, golden []byte, dec DecodedPayload, m *interfaceschema.CompiledMapping, object bool) {
/root/astrate-mule/pkg/payload/payload_test.go:251:	var redec DecodedPayload
/root/astrate-mule/pkg/payload/payload_test.go:389:func TestJSONObjectAggregation(t *testing.T) {
/root/astrate-mule/pkg/payload/bson.go:137:func decodeBSONObject(rv bson.RawValue, leaves map[string]*interfaceschema.CompiledMapping) (Value, error) {
/root/astrate-mule/pkg/payload/bson.go:317:func appendBSONObject(dst []byte, key string, obj map[string]Value) ([]byte, error) {

Let me look at the key coupling: how `ci` (the compiled interface used to decode) is derived, and understand the upstream fix #2141's concern. The concern in upstream was that the datastream could be decoded against the *wrong* interface's mappings when there were multiple object interfaces with the same last-level endpoint names.

In Astrate's code at data.go:264-278, `ci := realm.iface(name, ver.Major)` where `name`/`ver` come from `matchInterface(rest, dev.declares)` — the interface named in the **topic**. And `ci.ObjectLeaves` is that same `ci`'s leaves. So the object datastream is decoded against the topic-named interface's own mappings.

Wait — let me double check that `realm.iface` returns a per-realm interface, and that `ObjectLeaves` is per-interface. Also I need to verify there isn't some other code path where an object value is decoded against leaves from a different interface. Let me look at matchInterface and realm.iface.

Let me also check the upstream #2141 to understand the bug precisely.
$ rg -n "func matchInterface|func (e \*Engine) dec|func (r \*realmSchema) iface|ObjectLeaves" /root/astrate-mule/internal/engine/ | head -30
/root/astrate-mule/internal/engine/data.go:269:		dp, err = e.dec.Object(m.Payload, ci.ObjectLeaves)
/root/astrate-mule/internal/engine/data.go:331:	for leaf := range ci.ObjectLeaves {
/root/astrate-mule/internal/engine/topics.go:85:func matchInterface(rest string, declares func(string) (store.InterfaceVersion, bool)) (name, path string, ver store.InterfaceVersion, ok bool) {
/root/astrate-mule/internal/engine/cache_test.go:49:	if obj := alpha.iface("com.astrate.test.ObjectFlat", 1); obj == nil || obj.ObjectLeaves["latitude"] == nil {
/root/astrate-mule/internal/engine/serverdata.go:82:		dp, err = dec.Object(envelope, ci.ObjectLeaves)
/root/astrate-mule/internal/engine/serverdata.go:261:	for _, m := range ci.ObjectLeaves {

$ rg -n "func.*iface\b|ObjectLeaves" /root/astrate-mule/internal/engine/*.go
/root/astrate-mule/internal/engine/engine_e2e_test.go:215:func (env *e2eEnv) decodeIndividual(t *testing.T, iface string, major int, path string, body []byte) payload.DecodedPayload {
/root/astrate-mule/internal/engine/engine_e2e_test.go:236:func (env *e2eEnv) series(t *testing.T, dev *testutil.AstarteDevice, iface, path string) []store.IndividualRow {
/root/astrate-mule/internal/engine/engine_integration_test.go:68:func (r *engineRig) msg(iface, path string, qos byte, body []byte, ack *ackCounter) broker.InboundMessage {
/root/astrate-mule/internal/engine/cache_test.go:49:	if obj := alpha.iface("com.astrate.test.ObjectFlat", 1); obj == nil || obj.ObjectLeaves["latitude"] == nil {
/root/astrate-mule/internal/engine/data.go:269:		dp, err = e.dec.Object(m.Payload, ci.ObjectLeaves)
/root/astrate-mule/internal/engine/data.go:331:	for leaf := range ci.ObjectLeaves {
/root/astrate-mule/internal/engine/fakestore_test.go:590:func deviceMsg(iface, path string, qos byte, body []byte, ack *ackCounter) broker.InboundMessage {
/root/astrate-mule/internal/engine/fakestore_test.go:595:func deviceMsgFor(dev deviceid.ID, iface, path string, qos byte, body []byte, ack *ackCounter) broker.InboundMessage {
/root/astrate-mule/internal/engine/cache.go:93:func (r *realmSchema) iface(name string, major int) *interfaceschema.CompiledInterface {
/root/astrate-mule/internal/engine/serverdata.go:82:		dp, err = dec.Object(envelope, ci.ObjectLeaves)
/root/astrate-mule/internal/engine/serverdata.go:261:	for _, m := range ci.ObjectLeaves {

The key question is whether `ci` is the topic-named interface. Let me look at `realm.iface` and `matchInterface` to confirm, and also check `serverdata.go` since that's a separate object decode path (server-owned interfaces, though those wouldn't matter for type mismatches between two device object interfaces).

Let me look at cache.go:iface and topics.go:matchInterface.
→ Read internal/engine/cache.go [offset=80, limit=25]
→ Read internal/engine/topics.go [offset=85, limit=60]
→ Read pkg/payload/payload.go
`matchInterface` resolves the topic-named interface, `realm.iface(name, major)` fetches that interface's compiled form, and `e.dec.Object(m.Payload, ci.ObjectLeaves)` uses that interface's own leaves. Let me confirm the decoder validates each key's type against its own leaf mapping, and check the upstream issue to verify the exact concern.
→ Read pkg/interfaceschema/compile.go [limit=75]
Confirmed: each interface compiles its own `ObjectLeaves` map (compile.go:68,100), and the decoder type-checks each key against that map. The engine passes the topic-named interface's own leaves. Let me confirm the upstream issue context to characterize how this differs.

$ gh issue view 2141 2>&1 | head -80
GraphQL: Could not resolve to an issue or pull request with the number of 2141. (repository.issue)

% WebFetch https://github.com/astarte-platform/astarte/issues/2141
The PR was merged into `release-1.4` as commit `fdc6d29`. Let me look at the actual diff to characterize exactly what was fixed upstream.

$ curl -sL https://github.com/astarte-platform/astarte/commit/fdc6d29b84ec06a12316340b60f1b43491ca3d5c.diff | head -200
diff --git a/apps/astarte_appengine_api/mix.lock b/apps/astarte_appengine_api/mix.lock
index 0af1a900f2..5df8cbf7c6 100644
--- a/apps/astarte_appengine_api/mix.lock
+++ b/apps/astarte_appengine_api/mix.lock
@@ -5,7 +5,7 @@
   "astarte_rpc": {:hex, :astarte_rpc, "1.2.0", "dcef7434bf3f19ff30ff0bc245ef2d3b0f9abeb369405713cfd28916a5635926", [:mix], [{:amqp, "~> 3.3", [hex: :amqp, repo: "hexpm", optional: false]}, {:castore, "~> 1.0.0", [hex: :castore, repo: "hexpm", optional: false]}, {:protobuf, "~> 0.12", [hex: :protobuf, repo: "hexpm", optional: false]}, {:skogsra, "~> 2.2", [hex: :skogsra, repo: "hexpm", optional: false]}], "hexpm", "8470ed2f116fa8c9d70845f295f3738529aff123d9d3f98cbfa37107314763cb"},
   "b58": {:hex, :b58, "1.0.3", "d300d6ae5a3de956a54b9e8220e924e4fee1a349de983df2340fe61e0e464202", [:mix], [], "hexpm", "af62a98a8661fd89978cf3a3a4b5b2ebe82209de6ac6164f0b112e36af72fc59"},
   "bunt": {:hex, :bunt, "1.0.0", "081c2c665f086849e6d57900292b3a161727ab40431219529f13c4ddcf3e7a44", [:mix], [], "hexpm", "dc5f86aa08a5f6fa6b8096f0735c4e76d54ae5c9fa2c143e5a1fc7c1cd9bb6b5"},
-  "castore": {:hex, :castore, "1.0.20", "455e48f7115eca98c9f2b0e7a152b5a2e8f2a8a4f964c96e95bd31645ee5fa59", [:mix], [], "hexpm", "940eafbfd8b14bee649f083bc11b3b54ec555b54c3e4ea8213351ff6fee39c10"},
+  "castore": {:hex, :castore, "1.0.21", "0a0e8330dc267a40a3b7ad86d39302764bb71758172904e6a59d5ad6443ce307", [:mix], [], "hexpm", "e42e22723e25dbd46876d056a03f685513d6e98f6b5e555dc551321decd76c5c"},
   "cbor": {:hex, :cbor, "1.0.2", "9b0af85af291a556e10a0ffd48ba9a21a75e711828fafd3af193d56d95f0907f", [:mix], [], "hexpm", "edbc9b4a16eb93a582437b9b249c340a75af03958e338fb43d8c1be9fc65b864"},
   "certifi": {:hex, :certifi, "2.15.0", "0e6e882fcdaaa0a5a9f2b3db55b1394dba07e8d6d9bcad08318fb604c6839712", [:rebar3], [], "hexpm", "b147ed22ce71d72eafdad94f055165c1c182f61a2ff49df28bcc71d1d5b94a60"},
   "combine": {:hex, :combine, "0.10.0", "eff8224eeb56498a2af13011d142c5e7997a80c8f5b97c499f84c841032e429f", [:mix], [], "hexpm", "1b1dbc1790073076580d0d1d64e42eae2366583e7aecd455d1215b0d16f2451b"},
diff --git a/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex b/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex
index 0b3a902846..1674a2d430 100644
--- a/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex
+++ b/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex
@@ -316,11 +316,12 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.Interface do
   end
 
   def extract_mappings(
-        %InterfaceDescriptor{aggregation: :object},
+        %InterfaceDescriptor{aggregation: :object, interface_id: interface_id},
         _mapping,
         mappings
       ) do
     mappings
+    |> Enum.filter(fn {_id, m} -> m.interface_id == interface_id end)
     |> Map.new(fn {_id, m} ->
       key = m.endpoint |> String.split("/") |> List.last()
       {key, m}
diff --git a/apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_test.exs b/apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_test.exs
index e6b179133a..a37fbb757b 100644
--- a/apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_test.exs
+++ b/apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_test.exs
@@ -26,6 +26,7 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.DataHandlerTest do
 
   alias Astarte.Core.Mapping
   alias Astarte.DataUpdaterPlant.DataUpdater.Core.DataHandler
+  alias Astarte.DataUpdaterPlant.DataUpdater.Core.Interface
 
   import Astarte.InterfaceUpdateGenerators
 
@@ -91,6 +92,45 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.DataHandlerTest do
                  )
       end
     end
+
+    @tag regression: true
+    test "correctly validates object aggregate values for current interface", context do
+      %{
+        fixed_object_datastream_1: fixed_object_datastream_1,
+        fixed_object_datastream_2: fixed_object_datastream_2
+      } = context
+
+      timestamp = System.system_time(:microsecond) * 10
+      start = System.monotonic_time()
+
+      {:ok, _, state} =
+        Interface.maybe_handle_cache_miss(nil, fixed_object_datastream_1.name, context.state)
+
+      {:ok, _, state} =
+        Interface.maybe_handle_cache_miss(nil, fixed_object_datastream_2.name, state)
+
+      ctx = gen_context(state, fixed_object_datastream_1) |> Enum.at(0)
+
+      %{
+        interface: interface_name,
+        path: path,
+        payload: payload
+      } = ctx
+
+      assert {:ack, :ok, _, _} =
+               DataHandler.handle_data(state, interface_name, path, payload, timestamp, start)
+
+      ctx = gen_context(state, fixed_object_datastream_2) |> Enum.at(0)
+
+      %{
+        interface: interface_name,
+        path: path,
+        payload: payload
+      } = ctx
+
+      assert {:ack, :ok, _, _} =
+               DataHandler.handle_data(state, interface_name, path, payload, timestamp, start)
+    end
   end
 
   describe "validate_value_type/2" do
diff --git a/apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/interface_test.exs b/apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/interface_test.exs
index fef135e950..c0ae763ccb 100644
--- a/apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/interface_test.exs
+++ b/apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/interface_test.exs
@@ -209,6 +209,29 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.InterfaceTest do
       end
     end
 
+    @tag regression: true
+    test "extract_mappings/3 only considers mappings for the current object interface", context do
+      %{
+        fixed_object_datastream_1: fixed_object_datastream_1,
+        fixed_object_datastream_2: fixed_object_datastream_2,
+        state: state,
+        realm_name: realm_name
+      } = context
+
+      # ignored for objects
+      mapping = nil
+
+      descriptor = state.interfaces[fixed_object_datastream_1.name]
+
+      assert %{"first" => %{value_type: :string}, "second" => %{value_type: :integer}} =
+               Core.Interface.extract_mappings(descriptor, mapping, state.mappings)
+
+      descriptor = state.interfaces[fixed_object_datastream_2.name]
+
+      assert %{"first" => %{value_type: :boolean}, "second" => %{value_type: :binaryblob}} =
+               Core.Interface.extract_mappings(descriptor, mapping, state.mappings)
+    end
+
     property "forget_interfaces/2 removes interfaces from state cache", context do
       %{
         interfaces: interfaces,
diff --git a/apps/astarte_data_updater_plant/test/support/cases/device.ex b/apps/astarte_data_updater_plant/test/support/cases/device.ex
index 986093c7b6..e0d81b78ac 100644
--- a/apps/astarte_data_updater_plant/test/support/cases/device.ex
+++ b/apps/astarte_data_updater_plant/test/support/cases/device.ex
@@ -238,6 +238,10 @@ defmodule Astarte.Cases.Device do
        end},
       {:fixed_endpoint_interface,
        fn acc -> [new_interfaces(fixed_endpoint_interface(), acc, :single)] end},
+      {:fixed_object_datastream_1,
+       fn acc -> [new_interfaces(fixed_object_datastream_1(), acc, :single)] end},
+      {:fixed_object_datastream_2,
+       fn acc -> [new_interfaces(fixed_object_datastream_2(), acc, :single)] end},
       {:other_interfaces,
        fn acc -> new_interfaces(InterfaceGenerator.interface(), acc, :list) end}
     ]
@@ -252,6 +256,8 @@ defmodule Astarte.Cases.Device do
 
     %{
       interfaces: all_interfaces,
+      fixed_object_datastream_1: named_interfaces.fixed_object_datastream_1,
+      fixed_object_datastream_2: named_interfaces.fixed_object_datastream_2,
       fixed_endpoint_interface: named_interfaces.fixed_endpoint_interface,
       server_property_with_all_endpoint_types:
         named_interfaces.server_property_with_all_endpoint_types,
@@ -288,6 +294,58 @@ defmodule Astarte.Cases.Device do
     end)
   end
 
+  defp fixed_object_datastream_1 do
+    InterfaceGenerator.interface(
+      name: "test.FixedObjectDatastream1",
+      ownership: :device,
+      type: :datastream,
+      aggregation: :object
+    )
+    |> map(fn interface ->
+      mapping = Enum.at(interface.mappings, 0)
+
+      mapping_first = %{
+        mapping
+        | endpoint: "/data/first",
+          value_type: :string
+      }
+
+      mapping_second = %{
+        mapping
+        | endpoint: "/data/second",
+          value_type: :integer
+      }
+
+      %{interface | mappings: [mapping_first, mapping_second]}
+    end)
+  end
+
+  defp fixed_object_datastream_2 do
+    InterfaceGenerator.interface(
+      name: "test.FixedObjectDatastream2",
+      ownership: :device,
+      type: :datastream,
+      aggregation: :object
+    )
+    |> map(fn interface ->
+      mapping = Enum.at(interface.mappings, 0)
+
+      mapping_first = %{
+        mapping
+        | endpoint: "/data/first",
+          value_type: :boolean
+      }
+
+      mapping_second = %{
+        mapping
+        | endpoint: "/data/second",
+          value_type: :binaryblob
+      }
+
+      %{interface | mappings: [mapping_first, mapping_second]}

Now I have the upstream bug precisely. The fix is at interface.ex:316-325: upstream's `extract_mappings` built the object leaf table from a *flat* `state.mappings` map keyed by last path segment (across all interfaces), so a second object interface with the same last-level endpoint name shadowed the first's type. The fix filters to `m.interface_id == interface_id`.

Let me confirm Astrate has no such merged/flat table: check the rest of compile.go and grep for any other ObjectLeaves construction.
→ Read pkg/interfaceschema/compile.go [offset=76]

$ rg -n "ObjectLeaves\[|ObjectLeaves\s*=|ObjectLeaves=" /root/astrate-mule --glob '!**/*_test.go'
/root/astrate-mule/pkg/interfaceschema/compile.go:68:		ci.ObjectLeaves = make(map[string]*CompiledMapping, len(iface.Mappings))
/root/astrate-mule/pkg/interfaceschema/compile.go:97:			if _, dup := ci.ObjectLeaves[leaf]; dup {
/root/astrate-mule/pkg/interfaceschema/compile.go:100:			ci.ObjectLeaves[leaf] = cm

The only `ObjectLeaves` population is per-interface in compile.go:68-101, from that interface's own `iface.Mappings`. Let me verify the JSON/BSON object decoders actually type-check each key against its own leaf, and check whether any existing test covers two colliding object interfaces.
→ Read pkg/payload/json.go [offset=195, limit=45]

$ rg -ln "AggregationObject|aggregation: *object|object aggregate" /root/astrate-mule/internal/engine/*_test.go /root/astrate-mule/pkg/payload/*_test.go /root/astrate-mule/pkg/interfaceschema/*_test.go 2>/dev/null
/root/astrate-mule/pkg/interfaceschema/parse_test.go

Confirmed: decodeJSONObject type-checks each key against `leaves[key].ValueType`. The engine passes the topic-named interface's own leaves. Let me check the todo.md format to see if a follow-up is warranted — no existing engine test exercises two object interfaces with colliding leaf names (only an interfaceschema parse test exists).
$ tail -15 /root/astrate-mule/.mule/todo.md
- [x] store-delete-device-objects: in `internal/store/devices_test.go` `StatsAndDelete`, insert object datastream rows for the device before deleting it, and assert they are gone after the delete. `DeleteDevice` explicitly sweeps `object_datastreams` but the test only checks individual rows. [auto]
- [!] control-producer-properties-compression: in `internal/engine/control.go`, accept a plaintext (and the 4-zero-byte empty-frame) device→server `producer/properties` purge list when the device declared `purge_properties_compression_format: plaintext` — `inflateProperties` currently rejects every non-zlib payload while upstream's `control_handler.ex` decodes per-capability **Approved by Giulio 2026-08-31**: yes, devices that cannot compress must be able to talk to us, following upstream. — BLOCKED: gates failed
- [!] probe-interface-default-values: does `GET /realmmanagement/v1/<realm>/interfaces/<name>/<major>` return the same defaulted mapping-parameter values upstream reports after 3f0b864 ("Properly show default values for all mapping parameters")? report, do not patch [auto] — BLOCKED: gates failed
- [!] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto] — BLOCKED: gates failed
- [!] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
- [!] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto] — BLOCKED: gates failed
- [!] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
- [!] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto] — BLOCKED: gates failed
- [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
- [!] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto] — BLOCKED: gates failed
- [!] docs-sync-pairing-status-enum: in docs/api/astarte_pairing_api.yaml, fix the PairingInfo.status enum (line 367): it lists `confirmed, pending, denied, expired` but the handler only ever emits `confirmed`, `pending`, or `inhibited` (internal/pairing/service.go:289-297, via internal/pairing/http.go:211) — drop the dead `denied`/`expired` and add the undocumented `inhibited`. [auto] — BLOCKED: gates failed
- [!] issue-91-pairing-health-serve-root: in internal/pairing/http.go, add the unauthenticated `GET /pairing/health` route (no realm segment) sharing the handler/payload of the existing `GET /pairing/v1/{realm}/health`, keeping the v1 route as-is; tests for both paths (closes #91) — BLOCKED: gates failed
- [!] issue-68-async-operation-accepted: on the mutating endpoints upstream surfaces `async_operation` on (housekeeping realm create/delete; realm-management interface install/update/delete, trigger/policy delete in internal/realm/http.go), accept and ignore `?async_operation=false`, with unparseable/`true` values also not changing behaviour (Astrate stays always-sync); tests cover the flag parsing (closes #68) — BLOCKED: gates failed
- [ ] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto]
- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]

The probe is complete and the code answers the question. Let me append one follow-up since no engine-level regression test covers the exact #2141 scenario here.
$ cat >> /root/astrate-mule/.mule/todo.md <<'EOF'
- [ ] object-selected-interface-regression-test: add an internal/engine test with two object-aggregated interfaces sharing identical last-level endpoint names but different value types, asserting each publish is type-checked against the topic-named interface's own CompiledInterface.ObjectLeaves (the upstream #2141 scenario, currently uncovered at engine level) [auto]
EOF
(no output)

$ tail -2 /root/astrate-mule/.mule/todo.md
- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]
- [ ] object-selected-interface-regression-test: add an internal/engine test with two object-aggregated interfaces sharing identical last-level endpoint names but different value types, asserting each publish is type-checked against the topic-named interface's own CompiledInterface.ObjectLeaves (the upstream #2141 scenario, currently uncovered at engine level) [auto]

Task understood: probe whether Astrate's object-datastream decode uses the topic-named interface's own mappings, mirroring upstream #2141's fix. Findings:

**Upstream #2141 bug:** `extract_mappings/3` built an object interface's leaf table from the *flat* `state.mappings` cache keyed by last path segment across all interfaces. Two object interfaces with identical leaf names but different types collided, so a publish could be type-checked against the wrong interface's mapping; the fix filters `m.interface_id == interface_id`.

**Astrate — not affected, by construction:**
- `CompiledInterface.ObjectLeaves` is built per-interface in `Compile`, iterating only that interface's own `iface.Mappings` (compile.go:67-101); `rg` confirms it is constructed nowhere else. No flat/merged-by-leaf-key table exists anywhere.
- At decode time data.go:269 passes `ci.ObjectLeaves` where `ci = realm.iface(name, ver.Major)` and `name` is the interface prefix `matchInterface` (topics.go:85) resolves from the **topic** — the selected interface's own leaves.
- `decodeJSONObject`/`decodeBSONObject` type-check each document key against its leaf's `ValueType` (json.go:213-223), so wrong-type values on the selected interface are rejected.

Identical last-level names in a second object interface live in a separate map and cannot shadow the first's type-checking. No code change made (task says report, do not patch). No engine-level test currently covers the two-colliding-object-interfaces scenario.

```
Done: Astrate decodes and type-checks object datastreams against the topic-named interface's own CompiledInterface.ObjectLeaves, so the upstream #2141 wrong-interface type-checking bug cannot occur — no patch needed
Files: (report only, nothing touched except .mule/todo.md)
Verified: read internal/engine/data.go:238-278, topics.go:85-97, cache.go:93-98, pkg/interfaceschema/compile.go:67-101, pkg/payload/payload.go:100-136, json.go:201-225; upstream commit fdc6d29 diff -> pass (no identical bug found)
Unsure: nothing
Follow-ups: object-selected-interface-regression-test: add an internal/engine test with two object-aggregated interfaces sharing identical last-level endpoint names but different value types, asserting each publish is type-checked against the topic-named interface's own CompiledInterface.ObjectLeaves (the upstream #2141 scenario, currently uncovered at engine level) [auto]
```
