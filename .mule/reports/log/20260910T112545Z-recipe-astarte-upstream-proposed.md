slug: recipe-astarte-upstream
verdict: proposed
at:  9a28cd1
ran: 2026-09-10T11:25:45Z on DietPi in 268s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md
→ Read .mule/todo.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

The newest tag is `v1.4.0-rc.5`. The COMPATIBILITY.md declares `APICompatVersion` `1.2.2`. Upstream is ahead. Let me read the release notes.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
Sixth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.4

### Fixed

- [astarte_data_updater_plant] Make object values only validate type on selected interface
- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'
Let me check what these fixes actually changed in the upstream code.

$ gh search code --owner astarte-platform "validate_type" --repo astarte-platform/astarte_data_updater_plant --limit 5
(no output)

$ gh search code --owner astarte-platform "binaryblob" --repo astarte-platform/astarte_data_updater_plant --limit 5
(no output)

$ gh api repos/astarte-platform/astarte_data_updater_plant/commits?path=lib/data_updater_plant/data_updater.ex --jq '.[0:3][] | "\(.sha[0:7]) \(.commit.message | split("\n")[0])"'
(no output)

$ gh api repos/astarte-platform/astarte_data_updater_plant/tags --jq '.[0:10][] | "\(.name)\t\(.commit.sha[0:7])"'
v0.10.2	28ae849
v0.10.1	de1fe5f
v0.10.0	d1e8099
v0.10.0-rc.0	1b18dee
v0.10.0-beta.3	ba581c5
v0.10.0-beta.2	3d12ba1
v0.10.0-beta.1	819f5e6

$ rg -c 'binaryblob\|validate.*object\|object.*validate' /root/astrate-mule/internal/engine/ /root/astrate-mule/pkg/payload/ 2>/dev/null | head -20
(no output)

$ gh api repos/astarte-platform/astarte_data_updater_plant/commits --jq '.[] | select(.commit.message | test("validate type|binaryblob"; "i")) | "\(.sha[0:7]) \(.commit.message | split("\n")[0])"'
(no output)

$ rg -l 'binaryblob\|BinaryBlob' /root/astrate-mule/internal/ /root/astrate-mule/pkg/ 2>/dev/null
(no output)

$ rg -l 'object.*aggregate\|ObjectAggregate\|object_aggregat' /root/astrate-mule/internal/ 2>/dev/null
(no output)

$ rg -l 'validateValue\|ValidateValue\|validate.*value\|value.*valid' /root/astrate-mule/internal/engine/ 2>/dev/null
(no output)

$ rg -rl 'bson\|BSON' /root/astrate-mule/internal/engine/ /root/astrate-mule/pkg/ 2>/dev/null | head -10
(no output)

$ gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.4...v1.4.0-rc.5 --jq '.commits[] | "\(.sha[0:7]) \(.commit.message | split("\n")[0])"'
d859209 fix(dup): stop remapping interface loading failed error (#2119)
a816eaf chore: forward port release-1.3
aedc414 fix(dup): ensure binaryblob data is correctly validated
3e70d63 chore: prepare v1.3.3 (#2128)
fdc6d29 fix(dup): Make object values only validate type on selected interface (#2141)
4141a2a chore: forward-port release-1.3 into release-1.4
3484096 Prepare release v1.4.0-rc.5 (#2147)

$ rg -l 'value.*valid\|valid.*value\|ValidateMapping\|validateMapping' /root/astrate-mule/internal/ /root/astrate-mule/pkg/ 2>/dev/null | head -10
(no output)

$ rg -l 'data_updater\|DataUpdater\|DUP\|dataUpdater' /root/astrate-mule/internal/ 2>/dev/null | head -10
(no output)

$ rg -l 'validateData\|ValidateData\|validatePayload\|ValidatePayload\|validateIncoming\|incoming.*data\|processData\|handleData' /root/astrate-mule/internal/ 2>/dev/null | head -10
(no output)

$ rg -l 'maxValueSize\|MaxPayloadSize\|maxPayload\|MaxPayloadBytes\|payloadSize\|value_size' /root/astrate-mule/internal/ 2>/dev/null | head -10
(no output)

$ gh api repos/astarte-platform/astarte_data_updater_plant/commits/aedc414 --jq '.files[] | "\(.filename) +\(.additions) -\(.deletions)"'
{"message":"No commit found for SHA: aedc414","documentation_url":"https://docs.github.com/rest/commits/commits#get-a-commit","status":"422"}gh: No commit found for SHA: aedc414 (HTTP 422)

$ gh api repos/astarte-platform/astarte_data_updater_plant/commits/fdc6d29 --jq '.files[] | "\(.filename) +\(.additions) -\(.deletions)"'
{"message":"No commit found for SHA: fdc6d29","documentation_url":"https://docs.github.com/rest/commits/commits#get-a-commit","status":"422"}gh: No commit found for SHA: fdc6d29 (HTTP 422)

$ gh api repos/astarte-platform/astarte/pulls/2141 --jq '.title, .body'
fix(dup): Make object values only validate type on selected interface
#### What this PR does / why we need it:

#### Additional documentation e.g. usage docs, diagrams, reviewer notes, etc.:

<!--
This section can be blank if this pull request does not require additional resources.
-->

---

<details>
<summary><i>Thanks for sending a pull request! If this is your first time, here are some tips for you:</i></summary>

1. You can take a look at our [developer guide] for an introduction on Astarte development!
2. Make sure to read [CONTRIBUTING.md] and [CODE_OF_CONDUCT.md]
3. If the PR is unfinished or you're actively working on it, mark it as draft

When fixing existing issues, use [github's syntax to link your pull request] to it

> `fixes #<issue number>`

We also have a syntax to signal dependencies to other open pull requests

> `depends on #<pr number>`
> `depends on https://github.com/...`

In case of stacked PRs, you may add the PR number in the last commit's title instead:

> ```mermaid
> gitGraph
>     commit id: "Current master"
>     branch feat1
>     checkout feat1
>     commit id: "feat: add something"
>     commit id: "feat: add something else (#100)"
>     branch feat2
>     checkout feat2
>     commit id: "refactor: do something"
>     commit id: "fix: solve issue"
>     commit id: "feat: add a feature (#101)"
>     branch feat3
>     checkout feat3
>     commit id: "feat: feat without pr number"
> ```

</details>

[github's syntax to link your pull request]: https://docs.github.com/en/issues/tracking-your-work-with-issues/using-issues/linking-a-pull-request-to-an-issue#about-linked-issues-and-pull-requests
[developer guide]: https://docs.astarte-platform.org/astarte/snapshot/001-dev_guide.html
[CONTRIBUTING.md]: ../CONTRIBUTING.md
[CODE_OF_CONDUCT.md]: ../CODE_OF_CONDUCT.md

$ gh api repos/astarte-platform/astarte/commits/fdc6d29 --jq '.commit.message'
fix(dup): Make object values only validate type on selected interface (#2141)

Signed-off-by: Francesco Noacco <francesco.noacco@secomind.com>

$ gh api repos/astarte-platform/astarte/pulls/2141/files --jq '.[] | "\(.filename): +\(.additions) -\(.deletions)"'
apps/astarte_appengine_api/mix.lock: +1 -1
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex: +2 -1
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_test.exs: +40 -0
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/interface_test.exs: +23 -0
apps/astarte_data_updater_plant/test/support/cases/device.ex: +58 -0
apps/astarte_pairing/mix.lock: +1 -1

$ gh api repos/astarte-platform/astarte/pulls/2141/files --jq '.[].patch' 2>/dev/null | head -100
@@ -5,7 +5,7 @@
   "astarte_rpc": {:hex, :astarte_rpc, "1.2.0", "dcef7434bf3f19ff30ff0bc245ef2d3b0f9abeb369405713cfd28916a5635926", [:mix], [{:amqp, "~> 3.3", [hex: :amqp, repo: "hexpm", optional: false]}, {:castore, "~> 1.0.0", [hex: :castore, repo: "hexpm", optional: false]}, {:protobuf, "~> 0.12", [hex: :protobuf, repo: "hexpm", optional: false]}, {:skogsra, "~> 2.2", [hex: :skogsra, repo: "hexpm", optional: false]}], "hexpm", "8470ed2f116fa8c9d70845f295f3738529aff123d9d3f98cbfa37107314763cb"},
   "b58": {:hex, :b58, "1.0.3", "d300d6ae5a3de956a54b9e8220e924e4fee1a349de983df2340fe61e0e464202", [:mix], [], "hexpm", "af62a98a8661fd89978cf3a3a4b5b2ebe82209de6ac6164f0b112e36af72fc59"},
   "bunt": {:hex, :bunt, "1.0.0", "081c2c665f086849e6d57900292b3a161727ab40431219529f13c4ddcf3e7a44", [:mix], [], "hexpm", "dc5f86aa08a5f6fa6b8096f0735c4e76d54ae5c9fa2c143e5a1fc7c1cd9bb6b5"},
-  "castore": {:hex, :castore, "1.0.20", "455e48f7115eca98c9f2b0e7a152b5a2e8f2a8a4f964c96e95bd31645ee5fa59", [:mix], [], "hexpm", "940eafbfd8b14bee649f083bc11b3b54ec555b54c3e4ea8213351ff6fee39c10"},
+  "castore": {:hex, :castore, "1.0.21", "0a0e8330dc267a40a3b7ad86d39302764bb71758172904e6a59d5ad6443ce307", [:mix], [], "hexpm", "e42e22723e25dbd46876d056a03f685513d6e98f6b5e555dc551321decd76c5c"},
   "cbor": {:hex, :cbor, "1.0.2", "9b0af85af291a556e10a0ffd48ba9a21a75e711828fafd3af193d56d95f0907f", [:mix], [], "hexpm", "edbc9b4a16eb93a582437b9b249c340a75af03958e338fb43d8c1be9fc65b864"},
   "certifi": {:hex, :certifi, "2.15.0", "0e6e882fcdaaa0a5a9f2b3db55b1394dba07e8d6d9bcad08318fb604c6839712", [:rebar3], [], "hexpm", "b147ed22ce71d72eafdad94f055165c1c182f61a2ff49df28bcc71d1d5b94a60"},
   "combine": {:hex, :combine, "0.10.0", "eff8224eeb56498a2af13011d142c5e7997a80c8f5b97c499f84c841032e429f", [:mix], [], "hexpm", "1b1dbc1790073076580d0d1d64e42eae2366583e7aecd455d1215b0d16f2451b"},
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

$ rg -n 'object.*aggregate\|ObjectAggregate\|object_aggregate\|aggregate.*object' /root/astrate-mule/internal/ /root/astrate-mule/pkg/ 2>/dev/null | head -20
(no output)

$ rg -rn 'mapping.*interface\|interface.*mapping\|extractMapping\|ExtractMapping' /root/astrate-mule/internal/engine/ 2>/dev/null | head -20
(no output)

$ rg -rn 'aggregation\|Aggregation' /root/astrate-mule/internal/engine/ 2>/dev/null | head -20
(no output)

$ rg -rn 'value.*type\|type.*check\|validateType\|ValidateType' /root/astrate-mule/internal/engine/ 2>/dev/null | head -20
(no output)

$ ls /root/astrate-mule/internal/engine/
batch.go
batch_test.go
bus_events_test.go
cache.go
cache_test.go
capabilities.go
capabilities_test.go
control.go
control_test.go
data.go
data_test.go
devicedata_test.go
engine.go
engine_e2e_test.go
engine_integration_test.go
fakestore_test.go
forward
introspection.go
introspection_test.go
mqttv1_parity_test.go
router.go
router_test.go
serverdata.go
serverdata_test.go
stream
testdata
topics.go
topics_test.go
triggers
triggers_change_integration_test.go
triggers_change_test.go

$ rg -n 'extractMappings\|ExtractMappings\|mappingFor\|MappingFor\|mappingsFor\|MappingsFor\|lookupMapping\|LookupMapping\|matchMapping\|MatchMapping\|filterMapping\|FilterMapping' /root/astrate-mule/internal/ /root/astrate-mule/pkg/ 2>/dev/null | head -20
(no output)

$ rg -n 'object.*aggregate\|ObjectAggregate\|objectAggregate\|aggregat' /root/astrate-mule/internal/ /root/astrate-mule/pkg/ 2>/dev/null | head -20
(no output)

$ rg -rn 'incoming_data\|IncomingData\|handleData\|HandleData\|processData\|ProcessData' /root/astrate-mule/internal/ 2>/dev/null | head -20
(no output)

→ Read internal/engine/data.go

$ ls /root/astrate-mule/pkg/interfaceschema/
compat.go
compat_test.go
compile.go
compile_test.go
parse.go
parse_test.go
required_encrypted_test.go
testdata
trie.go
trie_test.go
types.go
violations.go
violations_test.go

$ rg -n 'ObjectLeaves\|objectLeaves' /root/astrate-mule/pkg/interfaceschema/ /root/astrate-mule/internal/ 2>/dev/null | head -30
(no output)

$ rg -n 'ObjectLeaves\|Leaves' /root/astrate-mule/pkg/interfaceschema/*.go | head -30
(no output)

✱ Grep "ObjectLeaves" in . · 28 matches
There's a prior report of this recipe run. Let me check it to avoid duplicating proposals.
→ Read .mule/reports/log/20260909T184315Z-recipe-astarte-upstream-proposed.md

$ git -C /root/astrate-mule log --oneline -5 && git -C /root/astrate-mule status --short
9a28cd1 mule: refill the queue
1285956 mule: recipe github-issues ran (97s)
7977031 mule: recipe milestones ran (214s)
e8ebeb6 mule: recipe hygiene timed out
979872f mule: log docs-sync-native-compat-version-endpoints

The previous run on 2026-09-09 already analyzed rc.5 thoroughly and appended nothing. Let me verify the git state and whether upstream has moved since then.
$ gh api repos/astarte-platform/astarte/releases --jq '.[0:3][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z

Upstream has not moved since the last recipe run: the newest tag is still `v1.4.0-rc.5` (2026-08-20). Let me confirm there's no newer tag my 5-entry window might have missed.
$ gh api repos/astarte-platform/astarte/tags --jq '.[0:3][] | "\(.name)"'
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3

$ ls -t /root/astrate-mule/.mule/reports/log/ | head -5
20260910T110956Z-recipe-github-issues-proposed.md
20260909T230411Z-recipe-milestones-proposed.md
20260909T230006Z-recipe-hygiene-timeout.md
20260909T202852Z-docs-sync-native-compat-version-endpoints-done.md
20260909T201715Z-docs-sync-native-compat-health-503-done.md

→ Read .mule/reports/log/20260909T184315Z-recipe-astarte-upstream-proposed.md [offset=1, limit=4]
Nothing changed since the last run. Let me check the report-log convention before closing.
Nothing has moved since the last run. I'll write the report-log evidence file, as the recipe convention does, and append nothing to the queue.

$ git -C /root/astrate-mule rev-parse HEAD && date -u +%Y%m%dT%H%M%SZ
9a28cd1a9f9af2451c441c7a0af43ee037154b09
20260910T112515Z

Upstream has not moved since the last run. The prior report (2026-09-09) already analyzed v1.4.0-rc.5 in full. I'll write the report log and append nothing.
← Write .mule/reports/log/20260910T112515Z-recipe-astarte-upstream-proposed.md
Wrote file successfully.

The recipe is complete. Upstream has not moved since the last run — newest tag is still `v1.4.0-rc.5` (2026-08-20), which the 2026-09-09 run already analyzed in full and found no wire-visible gap for Astrate across any of the rc.5/rc.4/rc.3/V1.3.3 changes. The standing compat-note lines are already in the queue. Per the recipe, an empty result is the good result: appended nothing.

No upstream movement since v1.4.0-rc.5.

Done: no upstream movement since v1.4.0-rc.5; analyzed and nothing proposeable; appended no task lines
Files: .mule/reports/log/20260910T112515Z-recipe-astarte-upstream-proposed.md
Verified: gh releases/tags queries -> pass (newest still v1.4.0-rc.5)
Unsure: nothing
Follow-ups: none
