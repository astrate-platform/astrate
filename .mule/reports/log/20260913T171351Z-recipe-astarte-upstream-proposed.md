slug: recipe-astarte-upstream
verdict: proposed
at:  bb70e4c
ran: 2026-09-13T17:13:51Z on DietPi in 334s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md
I'll read the compatibility doc first, then check upstream releases.

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

Newest upstream tag is `v1.4.0-rc.5`. Let me check the current state in the todo and what version Astrate targets.
→ Read .mule/todo.md
✱ Grep "1\.4|targets|currently|v1\.2|1\.2\.2|APICompatVersion" in docs · 73 matches
Newest is `v1.4.0-rc.5`. Let me check for anything newer (stable v1.4.0?), the full tag list, and the current for-giulio proposal state.
✱ Grep "v1\.4|v1\.3|COMPATIBILITY|compat-note|compat" in .mule · 26 matches

$ gh api 'repos/astarte-platform/astarte/releases?per_page=20' --jq '.[] | "\(.tag_name)\t\(.published_at)\tprerelease=\(.prerelease)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z	prerelease=true
v1.3.3	2026-08-07T14:33:58Z	prerelease=false
v1.4.0-rc.4	2026-08-06T13:40:20Z	prerelease=true
v1.4.0-rc.3	2026-07-31T10:48:39Z	prerelease=true
v1.4.0-rc.2	2026-07-14T10:16:22Z	prerelease=true
v1.4.0-rc.1	2026-07-14T08:53:27Z	prerelease=true
v1.3.2	2026-07-14T10:07:26Z	prerelease=false
v1.3.1	2026-07-13T11:50:05Z	prerelease=false
v1.3.0	2026-05-06T15:30:50Z	prerelease=false
v1.2.2	2026-04-27T16:39:28Z	prerelease=false
v1.4.0-rc.0	2026-04-08T12:32:33Z	prerelease=true
v1.3.0-rc.2	2026-04-08T16:34:59Z	prerelease=true
v1.2.2-rc.0	2026-04-08T16:31:03Z	prerelease=true
v1.2.1	2026-03-12T16:11:20Z	prerelease=false
v1.2.1-rc.1	2026-02-13T10:36:20Z	prerelease=true
v1.3.0-rc.1	2026-01-26T15:58:33Z	prerelease=true
v1.3.0-rc.0	2025-11-21T14:58:46Z	prerelease=true
v1.2.1-rc.0	2025-08-27T07:03:28Z	prerelease=true
v1.2.1-alpha.0	2025-04-10T15:31:04Z	prerelease=true
v1.2.0	2024-07-02T14:09:07Z	prerelease=false

Nothing newer than rc.5 since the milestone runs. Let me read the v1.4.0-rc.5 release notes, and check the rc.5↔rc.4 delta since previous proposal runs.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
Sixth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.4

### Fixed

- [astarte_data_updater_plant] Make object values only validate type on selected interface
- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated

Two bugs fixed upstream in rc.5. Let me check whether Astrate has analogous issues.
✱ Grep "binaryblob.*valid|valid.*binaryblob|object.*type.*valid|validateObject|objectMapValid" in internal · 0 matches

$ gh search code --owner astarte-platform "object values only validate type on selected interface" --limit 5 2>/dev/null || gh api 'repos/astarte-platform/astarte-data-updater-plant/commits?per_page=10' --jq '.[] | select(.commit.message | test("object|binaryblob|validate"; "i")) | "\(.sha[0:7]) \(.commit.message | split("\n")[0])"'
astarte-platform/astarte:CHANGELOG.md: - [astarte_data_updater_plant] Make object values only validate type on selected interface

✱ Grep "binaryblob|objectAggregate|validateType|matchValue|valueCheck|validateValue" in internal/engine · 10 matches

$ gh api 'repos/astarte-platform/astarte-data-updater-plant/commits?per_page=30' --jq '.[] | "\(.sha[0:12]) \(.commit.message | split("\n")[0])"'
{"message":"Not Found","documentation_url":"https://docs.github.com/rest/commits/commits#list-commits","status":"404"}gh: Not Found (HTTP 404)

$ gh api 'repos/astarte-platform/astarte/commits?per_page=15' --jq '.[] | "\(.sha[0:12]) \(.commit.message | split("\n")[0])"' 2>&1 | head -20
ad67a092755d refactor(vmq_plugin): use Umbrella config
6c2293138243 refactor(rm): use Umbrella config
da3c31816c59 refactor(dup): use Umbrella config
9ddf8a6d582f build(test_suite): remove obsolete dependency overrides (#2190)
ea1ae252b5ad refactor(pairing): use Umbrella config
0fe32a9cc9a7 refactor(housekeeping): use Umbrella configuration
ce0599c3f0b3 chore: forward-port release-1.4
682ab62dc55c build(e2e): fix docker build (#2185)
a1edc76de097 build: use proper logger level for umbrella configuration
2935c1e75046 refactor(te): use umbrella configuration
d98f18247793 build: enable dashboard typos (#2180)
2b1689aea7f0 chore(e2e): Init end to end next project (#2149)
1493be504cc3 chore: forward port release-1.4
2cb87896f7e6 refactor(appengine): use umbrella configuration
a8767f926c16 refactor: make existing apps/ build part of the umbrella (#2169)

$ gh api 'repos/astarte-platform/astarte/compare/v1.4.0-rc.4...v1.4.0-rc.5' --jq '.total_commits, (.commits[] | "\(.sha[0:12]) \(.commit.message | split("\n")[0])")' 2>&1 | head -40
7
d8592091e1a4 fix(dup): stop remapping interface loading failed error (#2119)
a816eaf0e6b0 chore: forward port release-1.3
aedc41483304 fix(dup): ensure binaryblob data is correctly validated
3e70d63f8935 chore: prepare v1.3.3 (#2128)
fdc6d29b84ec fix(dup): Make object values only validate type on selected interface (#2141)
4141a2ae456c chore: forward-port release-1.3 into release-1.4
3484096d2e70 Prepare release v1.4.0-rc.5 (#2147)

$ gh api 'repos/astarte-platform/astarte/commits/fdc6d29b84ec' --jq '.commit.message, "---FILES---", (.files[] | "\(.filename): \(.patch // "(no patch)")")' 2>&1 | head -120
fix(dup): Make object values only validate type on selected interface (#2141)

Signed-off-by: Francesco Noacco <francesco.noacco@secomind.com>
---FILES---
apps/astarte_appengine_api/mix.lock: @@ -5,7 +5,7 @@
   "astarte_rpc": {:hex, :astarte_rpc, "1.2.0", "dcef7434bf3f19ff30ff0bc245ef2d3b0f9abeb369405713cfd28916a5635926", [:mix], [{:amqp, "~> 3.3", [hex: :amqp, repo: "hexpm", optional: false]}, {:castore, "~> 1.0.0", [hex: :castore, repo: "hexpm", optional: false]}, {:protobuf, "~> 0.12", [hex: :protobuf, repo: "hexpm", optional: false]}, {:skogsra, "~> 2.2", [hex: :skogsra, repo: "hexpm", optional: false]}], "hexpm", "8470ed2f116fa8c9d70845f295f3738529aff123d9d3f98cbfa37107314763cb"},
   "b58": {:hex, :b58, "1.0.3", "d300d6ae5a3de956a54b9e8220e924e4fee1a349de983df2340fe61e0e464202", [:mix], [], "hexpm", "af62a98a8661fd89978cf3a3a4b5b2ebe82209de6ac6164f0b112e36af72fc59"},
   "bunt": {:hex, :bunt, "1.0.0", "081c2c665f086849e6d57900292b3a161727ab40431219529f13c4ddcf3e7a44", [:mix], [], "hexpm", "dc5f86aa08a5f6fa6b8096f0735c4e76d54ae5c9fa2c143e5a1fc7c1cd9bb6b5"},
-  "castore": {:hex, :castore, "1.0.20", "455e48f7115eca98c9f2b0e7a152b5a2e8f2a8a4f964c96e95bd31645ee5fa59", [:mix], [], "hexpm", "940eafbfd8b14bee649f083bc11b3b54ec555b54c3e4ea8213351ff6fee39c10"},
+  "castore": {:hex, :castore, "1.0.21", "0a0e8330dc267a40a3b7ad86d39302764bb71758172904e6a59d5ad6443ce307", [:mix], [], "hexpm", "e42e22723e25dbd46876d056a03f685513d6e98f6b5e555dc551321decd76c5c"},
   "cbor": {:hex, :cbor, "1.0.2", "9b0af85af291a556e10a0ffd48ba9a21a75e711828fafd3af193d56d95f0907f", [:mix], [], "hexpm", "edbc9b4a16eb93a582437b9b249c340a75af03958e338fb43d8c1be9fc65b864"},
   "certifi": {:hex, :certifi, "2.15.0", "0e6e882fcdaaa0a5a9f2b3db55b1394dba07e8d6d9bcad08318fb604c6839712", [:rebar3], [], "hexpm", "b147ed22ce71d72eafdad94f055165c1c182f61a2ff49df28bcc71d1d5b94a60"},
   "combine": {:hex, :combine, "0.10.0", "eff8224eeb56498a2af13011d142c5e7997a80c8f5b97c499f84c841032e429f", [:mix], [], "hexpm", "1b1dbc1790073076580d0d1d64e42eae2366583e7aecd455d1215b0d16f2451b"},
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex: @@ -316,11 +316,12 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.Interface do
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
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_test.exs: @@ -26,6 +26,7 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.DataHandlerTest do
 
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
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/interface_test.exs: @@ -209,6 +209,29 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.InterfaceTest do
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
apps/astarte_data_updater_plant/test/support/cases/device.ex: @@ -238,6 +238,10 @@ defmodule Astarte.Cases.Device do
        end},
       {:fixed_endpoint_interface,
        fn acc -> [new_interfaces(fixed_endpoint_interface(), acc, :single)] end},
+      {:fixed_object_datastream_1,
+       fn acc -> [new_interfaces(fixed_object_datastream_1(), acc, :single)] end},
+      {:fixed_object_datastream_2,
+       fn acc -> [new_interfaces(fixed_object_datastream_2(), acc, :single)] end},
       {:other_interfaces,

$ gh api 'repos/astarte-platform/astarte/commits/aedc41483304' --jq '.commit.message, "---FILES---", (.files[] | "\(.filename): \(.patch // "(no patch)")")' 2>&1 | head -150
fix(dup): ensure binaryblob data is correctly validated

Signed-off-by: Francesco Noacco <francesco.noacco@secomind.com>
---FILES---
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/data_handler.ex: @@ -509,7 +509,12 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.DataHandler do
   # From Cyanide 2.0, binaries are decoded as %Cyanide.Binary{}
   def validate_value_type(%Mapping{value_type: expected_type}, %Cyanide.Binary{} = value) do
     %Cyanide.Binary{subtype: _subtype, data: bin} = value
-    validate_value_type(expected_type, bin)
+    ValueType.validate_value(expected_type, bin)
+  end
+
+  # All binaries should be %Cyanide.Binary{}, if we're here something went wrong
+  def validate_value_type(%Mapping{value_type: :binaryblob}, binary) when is_binary(binary) do
+    {:error, :unexpected_value_type}
   end
 
   # Explicitly match on all other structs to avoid pattern matching them as maps below
@@ -534,8 +539,8 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.DataHandler do
   # object interface: validate all the types of the contained mappings
   def validate_value_type(%{} = mappings_by_key, %{} = object) do
     Enum.reduce_while(object, :ok, fn {key, value}, _acc ->
-      with {:ok, %Mapping{value_type: expected_type}} <- Map.fetch(mappings_by_key, key),
-           :ok <- ValueType.validate_value(expected_type, value) do
+      with {:ok, mapping} <- Map.fetch(mappings_by_key, key),
+           :ok <- validate_value_type(mapping, value) do
         {:cont, :ok}
       else
         {:error, reason} ->
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_test.exs: @@ -24,6 +24,7 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.DataHandlerTest do
   use Astarte.Cases.DataUpdater
   use ExUnitProperties
 
+  alias Astarte.Core.Mapping
   alias Astarte.DataUpdaterPlant.DataUpdater.Core.DataHandler
 
   import Astarte.InterfaceUpdateGenerators
@@ -92,8 +93,21 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.DataHandlerTest do
     end
   end
 
+  describe "validate_value_type/2" do
+    test "returns ok for valid binaryblob" do
+      binary = %Cyanide.Binary{subtype: :generic, data: <<1, 2, 3, 4>>}
+
+      assert :ok = DataHandler.validate_value_type(%Mapping{value_type: :binaryblob}, binary)
+    end
+
+    test "returns error for raw binaries in binaryblobs" do
+      assert {:error, :unexpected_value_type} =
+               DataHandler.validate_value_type(%Mapping{value_type: :binaryblob}, <<1, 2, 3, 4>>)
+    end
+  end
+
   defp gen_context(state, interface) do
-    gen all update <- valid_mapping_update_for(interface),
+    gen all update <- valid_complete_mapping_update_for(interface),
             timestamp <- repeatedly(fn -> DateTime.utc_now(:millisecond) end) do
       payload =
         %{
apps/astarte_data_updater_plant/test/support/interface_update_generators.ex: @@ -368,7 +368,9 @@ defmodule Astarte.InterfaceUpdateGenerators do
         |> map(&DateTime.from_unix!/1)
       ])
 
-    binaryblob = binary(max_length: 65_535, min_length: min_length)
+    binaryblob =
+      binary(max_length: 65_535, min_length: min_length)
+      |> map(&%Cyanide.Binary{subtype: :generic, data: &1})
 
     scalar_generators = %{
       double: float(),

Now let me check how Astrate validates object-aggregated values and binaryblobs.
✱ Grep "func.*[Vv]alidate|object|Object" in internal/engine · 100 matches
✱ Grep "func.*Object|ObjectLeaves|Leaves" in pkg/payload · 19 matches
→ Read pkg/payload/payload.go [offset=60, limit=140]
→ Read pkg/payload/bson.go [offset=130, limit=80]
→ Read pkg/payload/bson.go [limit=95, offset=36]
✱ Grep "ValueBinaryblob|Binaryblob|BinaryBlob|subtype" in pkg/payload · 15 matches
✱ Grep "astarte-upstream" in .mule · 100 matches
Let me check the last time this specific recipe ran and what it concluded.
✱ Glob ".mule/reports/*" 0 matches

$ ls -la /root/astrate-mule/.mule/reports/log/ | grep -i upstream
-rw-r--r-- 1 root root 152510 Sep  5 20:46 20260905T184625Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  57962 Sep  6 20:31 20260906T183102Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root    218 Sep  7 04:03 20260907T020359Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  82600 Sep  7 19:55 20260907T175524Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  85176 Sep  8 13:34 20260908T113422Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   2013 Sep  8 21:38 20260908T193614Z-recipe-astarte-upstream-noop.md
-rw-r--r-- 1 root root  28135 Sep  8 21:38 20260908T193814Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  27766 Sep  9 20:43 20260909T184315Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   1602 Sep 10 13:25 20260910T112515Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  17437 Sep 10 13:25 20260910T112545Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   2160 Sep 10 22:29 20260910T202936Z-recipe-astarte-upstream-noop.md
-rw-r--r-- 1 root root  17656 Sep 10 22:29 20260910T202949Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  19257 Sep 11 21:15 20260911T191541Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  45876 Sep 12 20:07 20260912T180735Z-recipe-astarte-upstream-proposed.md

$ tail -60 /root/astrate-mule/.mule/reports/log/20260910T202949Z-recipe-astarte-upstream-proposed.md
29902:← Write .mule/reports/log/20260910T112515Z-recipe-astarte-upstream-proposed.md
29909:Done: no upstream movement since v1.4.0-rc.5; analyzed and nothing proposeable; appended no task lines
29910:Files: .mule/reports/log/20260910T112515Z-recipe-astarte-upstream-proposed.md
30032:Let me check what the `flow` package's tests cover to avoid duplicating proposals.
30468:Done: code review of internal/flowapi — 4 proposals: auto-restart retries forever on permanent failure, restart goroutine unreachable on shutdown, validationDetail untested, auto_restart default untested
30470:Verified: none run — proposal run, no code changed
31180:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
31941:When both sources are empty a tick runs a **proposal recipe** instead, rotating through
31942:`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
31989:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
32007:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
32360:>> queue is empty — running the 'astarte-upstream' recipe to propose work

Let me read the most recent `astarte-upstream` proposal run report (2026-09-10) to confirm the state and avoid duplication.
→ Read .mule/reports/log/20260910T112515Z-recipe-astarte-upstream-proposed.md
Last run (2026-09-10, sha 9a28cd1) already determined no upstream movement since v1.4.0-rc.5 and proposed nothing. Let me confirm the current state hasn't changed since then.

$ git rev-parse HEAD && echo "---deps-untouched check---" && gh api repos/astarte-platform/astarte/releases --jq '.[0:1][] | "\(.tag_name)\t\(.published_at)"'
eb44d02b00c7cd8fcc3db6d938748d5437085304
---deps-untouched check---
v1.4.0-rc.5	2026-08-20T12:28:04Z

→ Read .mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md
Newest tag is unchanged: still `v1.4.0-rc.5` (2026-08-20). The last recipe run (2026-09-10, sha 9a28cd1) already analyzed rc.4→rc.5 and every prior gap; the codebase has since moved to eb44d02 but upstream has not. Let me check the convention for a no-movement evidence file.

$ ls /root/astrate-mule/.mule/reports/log/ | tail -20
20260909T191420Z-recipe-code-review-proposed.md
20260909T195057Z-realm-pure-helper-tests-done.md
20260909T200459Z-recipe-docs-sync-proposed.md
20260909T201715Z-docs-sync-native-compat-health-503-done.md
20260909T202852Z-docs-sync-native-compat-version-endpoints-done.md
20260909T230006Z-recipe-hygiene-timeout.md
20260909T230411Z-recipe-milestones-proposed.md
20260910T110956Z-recipe-github-issues-proposed.md
20260910T112515Z-recipe-astarte-upstream-proposed.md
20260910T112545Z-recipe-astarte-upstream-proposed.md
20260910T170149Z-recipe-code-review-proposed.md
20260910T171629Z-flowapi-autorestart-terminal-failure-done.md
20260910T172606Z-flowapi-validationdetail-test-done.md
20260910T182250Z-flowapi-autorestart-default-test-done.md
20260910T184706Z-recipe-docs-sync-proposed.md
20260910T190522Z-docs-sync-pairing-status-enum-done.md
20260910T191129Z-docs-sync-pairing-version-endpoint-done.md
20260910T192802Z-recipe-hygiene-proposed.md
20260910T195304Z-recipe-milestones-proposed.md
20260910T200120Z-recipe-github-issues-proposed.md

Upstream has not moved since the last check (still `v1.4.0-rc.5`, 2026-08-20), and the prior run already analyzed everything and appended nothing. Per the recipe, an empty result is the correct result: append nothing, write the evidence file.
$ date -u +'%Y%m%dT%H%M%SZ'
20260910T202936Z

← Write .mule/reports/log/20260910T202936Z-recipe-astarte-upstream-noop.md
Wrote file successfully.

Done: upstream watch — no upstream movement since v1.4.0-rc.5; appended nothing to `.mule/todo.md`

Files: `.mule/reports/log/20260910T202936Z-recipe-astarte-upstream-noop.md` (evidence; todo.md untouched)
Verified: `gh api repos/astarte-platform/astarte/releases` -> newest remains v1.4.0-rc.5 (2026-08-20), unchanged since the last run; nothing new to analyze or propose
Unsure: nothing
Follow-ups: none

$ for f in /root/astrate-mule/.mule/reports/log/20260911T191541Z-recipe-astarte-upstream-proposed.md /root/astrate-mule/.mule/reports/log/20260912T180735Z-recipe-astarte-upstream-proposed.md; do echo "===== $f ====="; tail -8 "$f"; done
===== /root/astrate-mule/.mule/reports/log/20260911T191541Z-recipe-astarte-upstream-proposed.md =====

```
Done: no upstream movement since v1.4.0-rc.5 (newest stable v1.3.3, unchanged since 2026-08-07); nothing appended to .mule/todo.md
Files: none
Verified: gh api repos/astarte-platform/astarte/releases; git ls-remote --tags -> pass (newest tag v1.4.0-rc.5)
Unsure: nothing
Follow-ups: none
```
===== /root/astrate-mule/.mule/reports/log/20260912T180735Z-recipe-astarte-upstream-proposed.md =====

```
Done: upstream unchanged — no releases since v1.4.0-rc.5 (2026-08-20) / v1.3.3 (2026-08-07); all wire-relevant items already proposed in the queue
Files: none
Verified: gh api repos/astarte-platform/astarte/releases -> all ≤ 2026-08-20; tags list -> same; diff v1.2.2...v1.3.3 -> nothing new
Unsure: nothing
Follow-ups: none
```

The 2026-09-10 run already found nothing to propose. Let me check the most recent runs (09-11, 09-12) for their verdicts.
Upstream unchanged since all prior runs. Let me generate the evidence file per the recipe's established noop convention.
$ date -u +'%Y%m%dT%H%M%SZ' && gh api 'repos/astarte-platform/astarte/tags?per_page=3' --jq '.[] | .name'
20260913T171327Z
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3

Upstream confirmed unchanged. Writing the noop evidence file.
← Write .mule/reports/log/20260913T171327Z-recipe-astarte-upstream-noop.md
Wrote file successfully.

Upstream unchanged: newest tag still `v1.4.0-rc.5` (2026-08-20), newest stable `v1.3.3`. No releases since the last proposal runs (09-11, 09-12), so I re-verified only the newest unreviewed delta (rc.4→rc.5, two DUP validation fixes) and found neither maps to Astrate: objects are validated against the current interface's leaves only (internal/engine/data.go:278, pkg/payload/payload.go:102), and binaryblob accepts only a BSON binary subtype, rejecting others (pkg/payload/bson.go:120-125). The v1.4-rc compat-note is already queued. Appended nothing to todo.md.

```
Done: upstream watch — no upstream movement since v1.4.0-rc.5; nothing appended to .mule/todo.md
Files: .mule/reports/log/20260913T171327Z-recipe-astarte-upstream-noop.md (evidence; todo.md untouched)
Verified: gh api repos/astarte-platform/astarte/releases + tags -> newest still v1.4.0-rc.5 (2026-08-20), v1.3.3 stable (2026-08-07); rc.4→rc.5 commits read (fdc6d29, aedc414) and both mapped against Astrate code paths -> pass/not applicable
Unsure: nothing
Follow-ups: none
```
