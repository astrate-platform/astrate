slug: compat-note-v133
verdict: blocked
at:  2f072d3
ran: 2026-08-31T19:05:59Z on DietPi in 199s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/COMPATIBILITY.md
→ Read .mule/for-giulio.md
The task is: propose the `docs/COMPATIBILITY.md` wording for v1.3.3 in `.mule/for-giulio.md` (do not edit the file itself).

Let me gather the necessary context.

$ cd /root/astrate-mule && git log --oneline -30
2f072d3 mule: blocked probe-emptycache-resend-device-error
4c1228c mule: blocked compat-note-v140-rc3
acfcb26 mule: blocked probe-value-type-validation
db8048e mule: blocked probe-interface-default-values
eeb69d4 mule: blocked control-producer-properties-compression
d044364 mule: blocked issue-68
8e976a6 mule: blocked issue-91
7b03846 mule: recipe astarte-upstream ran (224s)
6e30a1c queue: drop three lines already done or duplicated, approve the plaintext-purge line, take the issues-first ordering from main
9233c4c mule: pin the model (from main)
6d1055b mule: take the cron survival fix from main (PATH + cwd)
071d0a2 mule: recipe astarte-upstream ran (396s)
28f3adf mule: fix ticks dying outside the repo (cron cwd), recover unsaved queue state
83568dd mule: recipe astarte-upstream ran (333s)
3a4b693 mule: recipe astarte-upstream ran (431s)
3a746d1 mule: recipe astarte-upstream ran (485s)
4f6d95c mule: split race-check into 5 per-subsystem lines
4beddf0 mule: issue-13 passed on beecab8
beecab8 mule: issue-20 passed on 33d5be6
33d5be6 mule: log issue-12
89b1004 mule: issue-12: device-registered-trigger: emit DeviceRegisteredEvent from internal/pairing/service.go on registration
647b803 mule: issue-13 passed on 4af2073
4af2073 mule: blocked issue-14
7d87440 mule: issue-20 passed on c59a7eb
c59a7eb mule: log issue-23
a9481d0 mule: issue-23: flow-models: FlowMessage wire format and Block/Pipeline/Flow data models
04718fd mule: blocked issue-24
6b2bd3a mule: log issue-25
ee59732 mule: issue-25: flow-lifecycle: Flow start/stop lifecycle and runtime management
6ab3307 mule: log issue-26

$ cd /root/astrate-mule && git tag --list 'v1.3*' --sort=-version:refname
(no output)

Now let me understand what v1.3.3 entails — I need to find the relevant tag/commits and what changed.
$ cd /root/astrate-mule && git log --oneline v1.3.2..HEAD 2>/dev/null | head -60
(no output)

$ cd /root/astrate-mule && git log --oneline --all | head -40
2f072d3 mule: blocked probe-emptycache-resend-device-error
4c1228c mule: blocked compat-note-v140-rc3
acfcb26 mule: blocked probe-value-type-validation
db8048e mule: blocked probe-interface-default-values
eeb69d4 mule: blocked control-producer-properties-compression
d044364 mule: blocked issue-68
8e976a6 mule: blocked issue-91
7b03846 mule: recipe astarte-upstream ran (224s)
6e30a1c queue: drop three lines already done or duplicated, approve the plaintext-purge line, take the issues-first ordering from main
60a8ec6 mule: approved issues outrank the mule's own proposals
ce759e1 survey: recover 16 daily reports the broken commit check left staged (2026-07-28..2026-08-26), and take the tooling fix from main
271c055 survey: actually commit what it wrote
9233c4c mule: pin the model (from main)
4e0ac73 mule: pin the model explicitly — the agent's default no longer resolves
6d1055b mule: take the cron survival fix from main (PATH + cwd)
306adc0 mule: survive cron — run from the repo, find opencode under $HOME
c522db1 docs(agents): Legion Go is off by default — probe once, then ask Giulio
72b025e chore(lint): clear golangci-lint debt to zero
7a244b2 feat(payload): reject object documents missing a required key (#67 fase 4b)
fd562f7 payload: pre-wire missing_required reject reason label fixture and error-name translation (fase 4b prep)
9c4d411 feat(interfaceschema): required+encrypted mapping fields parse and gate (#67 fase 4a)
84ae2b6 flow: dynamic virtual_device_pool — first-seen auto-registration (#84)
0376b83 docs: resolve dormant-trigger decision line — measured in #20, activated by 6bd14a7
1e4d4d5 docs: probe-craft lessons from the 2026-08-24 verify batch; clarify the Mac test rule
ab31319 docs(compat): measured write taxonomy + unknown-realm deviation; fix Internal server error casing to measured upstream form
ad4eb45 fix(appengine): server-write error taxonomy matches measured upstream (#57)
cda4010 feat(http): per-service version endpoints, upstream-parity (#77)
907b6e7 test(conformance): record fase-1 upstream verify fixtures (#69/#57); fix Bad request casing to measured upstream form
071d0a2 mule: recipe astarte-upstream ran (396s)
14708db feat(flowapi): composites accepted and expanded in the flow lifecycle (issue #85, phase 85d)
4b6dfc7 feat(flowapi): user-block CRUD endpoints with schema validation (issue #85, phase 85c)
4dec328 build: add santhosh-tekuri/jsonschema/v6 (issue #85, prep 85c)
9025404 feat(store): per-realm user_blocks CRUD (issue #85, phase 85b)
1cee640 mule: record FDO milestone-4.0 decision in for-giulio
8c8804d migrations: user_blocks table for composite blocks (issue #85, prep 85b)
78e293e feat(flow): ExpandComposites — inline expansion of user blocks with cycle guard (issue #85, phase 85a)
6d60193 feat(flowapi): plumb the engine ingest path into flow Deps (issue #84, phase 84c)
ad55cab feat(flow/blocks): virtual_device_pool block over the engine ingest path (issue #84, phase 84b)
970ce71 feat(flow/blocks): register virtual_device_pool type (stub) and Ingest dep (issue #84, prep 84b)
5e1aa1d feat(engine): PublishDeviceValue — device-owned ingest with no MQTT delivery (issue #84, phase 84a)

→ Read .mule/milestones.md

$ cd /root/astrate-mule && rg -n "1\.3\.3\|v1\.3\.3\|1.3.3" --type md -l
(no output)

$ cd /root/astrate-mule && rg -rn "1\.3\." --type md | head -40
(no output)

$ cd /root/astrate-mule && rg -n "compat" --type md -l | head -20
docs/ROADMAP.md
bench/README.md
README.md
pkg/interfaceschema/testdata/README.md
docs/handoff/phase-9-memory.md
docs/handoff/phase-10-memory.md
docs/handoff/phase-3-memory.md
docs/site/compatibility.md
docs/site/deployment.md
docs/handoff/phase-1-memory.md
docs/site/appengine-api.md
docs/COMPATIBILITY.md
docs/site/index.md
docs/site/operations.md
docs/site/realm-management-api.md
docs/site/migration-from-astarte.md
docs/DESIGN.md

→ Read .mule/todo.md

$ cd /root/astrate-mule && ls .mule/recipes/ && echo "---" && ls .mule/ && echo "---TASKS---" && ls .mule/tasks/ 2>/dev/null
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md
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
task.md
tasks
todo.md
---TASKS---
issue-16.md
realm-config-datastream-retention.md

The task is `compat-note-v1.3.3`: propose the `docs/COMPATIBILITY.md` wording for v1.3.3 in `.mule/for-giulio.md`. The version v1.3.3 refers to an *upstream Astarte* release (the probes mention "upstream v1.3.3 (#2119)"). This is about documenting what changed in the compatibility document for the v1.3.3 conformance target.

Let me understand what "v1.3.3" means here. Looking at the conformance pinned versions, and the precedent task `compat-note-v1.4.0-rc.3` which was blocked. Let me understand what these compat-note tasks actually do — I need to find the upstream v1.3.3 release changes relevant to compatibility.

Let me look at what changed. The v1.3.3 reference here is upstream Astarte. Let me check the upstream release and what's relevant. But wait — the task says "propose the docs/COMPATIBILITY.md wording for v1.3.3". 

Let me look at how previous compat-note tasks were handled, and what the docs-sync recipe does. Let me check the recipes.
→ Read .mule/recipes/astarte-upstream.md
→ Read .mule/recipes/docs-sync.md
Now I understand the context. The `compat-note-v1.3.3` task asks me to figure out what `docs/COMPATIBILITY.md` wording needs to change for the v1.3.3 conformance/upstream target, and propose it in `.mule/for-giulio.md`.

Let me look at what "v1.3.3" refers to within this repo. The COMPATIBILITY.md documents conformance checkpoints pinned against upstream. Let me check how the conformance suite references versions and what upstream v1.3.3 entailed.

Let me check the docs-sync recipe and how COMPATIBILITY.md mentions versions, and what the v1.3.3 upstream release changed.

$ cd /root/astrate-mule && rg -rn "v1\.3\.3\|1\.3\.3\|2119" docs/ .mule/ test/ 2>/dev/null | head -40
(no output)

v1.3.3 doesn't appear in the repo yet. Let me understand what the "compat-note" task should produce. The task (in the queue) says: "propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]".

The purpose: upstream Astarte v1.3.3 was released, and the COMPATIBILITY.md should be updated to reflect it. The conformance suite pins checkpoints against upstream versions. The "wording" task asks me to figure out what the COMPATIBILITY.md should say for v1.3.3.

Note the blocked sibling `compat-note-v1.4.0-rc.3` and the probes referencing upstream v1.3.3 (#2119) — probe-emptycache-resend-device-error mentions "upstream v1.3.3 (#2119)". So v1.3.3 is a real upstream Astarte release that the compatibility doc should mention.

Let me look at the upstream Astarte v1.3.3 release to understand what changed. First, let me check the upstream to understand what v1.3.3 is and whether it's relevant. Let me look at how the COMPATIBILITY.md and the conformance checkpoints reference specific upstream versions.

Actually, looking more carefully at COMPATIBILITY.md, the version references are the *client* versions (SDKs, astartectl), not the upstream *server* version. The upstream server version isn't pinned in the doc other than in prose (e.g., "astarte-dashboard v1.2.2", "upstream v1.2.0", "v1.4.0-rc.3"). 

Let me check what upstream v1.3.3 was and what it changed relevant to compatibility. Let me query the GitHub API.
$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/releases --jq '.[0:10][] | "\(.tag_name)\t\(.published_at)"' 2>&1 | head -20
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

Upstream v1.3.3 was released 2026-08-07. Let me get its release notes and the diff from v1.3.2, focusing on what Astrate implements (wire-visible compatibility obligations).
$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body' 2>&1

$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/compare/v1.3.2...v1.3.3 --jq '.commits[].commit.message' 2>&1 | head -80
fix(dup): stop remapping interface loading failed error (#2119)

Signed-off-by: Riccardo Nalgi <riccardo.nalgi@secomind.com>
chore: prepare v1.3.3 (#2128)

Signed-off-by: Eddy Babetto <eddy.babetto@secomind.com>

So upstream v1.3.3 contains essentially one change: `fix(dup): stop remapping interface loading failed error (#2119)`. This is directly the subject of the probe task `probe-emptycache-resend-device-error` which references "upstream v1.3.3 (#2119)".

Let me look at that change to understand its wire-visible behaviour. Let me get the details of PR #2119.
$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/commits?path= --jq '.[]' . 2>/dev/null | head -5; echo "---commits for v1.3.3 compare detail---"; gh api repos/astarte-platform/astarte/compare/v1.3.2...v1.3.3 --jq '.commits[] | {sha: .sha[0:8], msg: .commit.message}' 2>&1 | head -40
---commits for v1.3.3 compare detail---
{"msg":"fix(dup): stop remapping interface loading failed error (#2119)\n\nSigned-off-by: Riccardo Nalgi \u003criccardo.nalgi@secomind.com\u003e","sha":"d8592091"}
{"msg":"chore: prepare v1.3.3 (#2128)\n\nSigned-off-by: Eddy Babetto \u003ceddy.babetto@secomind.com\u003e","sha":"3e70d63f"}

$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/issues/2119 --jq '{title, body, state}' 2>&1 | head -80
{"body":"When resending properties to a device whose introspection includes an interface that isn't installed in Astarte, the device error trigger fired resend_interface_properties_failed instead of the actual cause, interface_loading_failed\r\n\r\nThis PR fixes this [issue](http://github.com/astarte-platform/astarte/issues/2118) by making Astarte correctly report interface_loading_failed.\r\n\r\n\r\n","state":"closed","title":"fix(dup): stop remapping interface loading failed error"}

Now I understand. Let me look at the actual upstream change in #2119 to understand the wire-visible behaviour precisely, then check what Astrate currently does (`internal/engine/control.go` `resendServerProperties`) — which is what the `probe-emptycache-resend-device-error` task covers.

Let me look at Astrate's current control.go `resendServerProperties` behaviour.
$ cd /root/astrate-mule && rg -n "resendServerProperties\|interface_loading_failed\|resend_interface_properties_failed\|emptyCache" internal/engine/control.go
(no output)

$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/pulls/2119/files --jq '.[] | {filename, patch}' 2>&1 | head -120
{"filename":"CHANGELOG.md","patch":"@@ -5,6 +5,12 @@ All notable changes to this project will be documented in this file.\n The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)\n and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).\n \n+## Unreleased\n+\n+### Fixed\n+\n+- [astarte_data_updater_plant] Resending properties to a device now reports distinct device error triggers depending on the failure: `interface_loading_failed` when an interface fails to load and `resend_interface_properties_failed` when sending the properties to the device itself fails. Previously, both cases were incorrectly reported as `resend_interface_properties_failed`.\n+\n ## [1.3.2] - 2026-07-14\n \n ## Fixed"}
{"filename":"apps/astarte_appengine_api/mix.lock","patch":"@@ -4,7 +4,7 @@\n   \"astarte_core\": {:git, \"https://github.com/astarte-platform/astarte_core.git\", \"e6c245d2642cf499ba29a78d08e5d3d6de47badc\", [tag: \"v1.3.2\"]},\n   \"astarte_rpc\": {:hex, :astarte_rpc, \"1.2.0\", \"dcef7434bf3f19ff30ff0bc245ef2d3b0f9abeb369405713cfd28916a5635926\", [:mix], [{:amqp, \"~\u003e 3.3\", [hex: :amqp, repo: \"hexpm\", optional: false]}, {:castore, \"~\u003e 1.0.0\", [hex: :castore, repo: \"hexpm\", optional: false]}, {:protobuf, \"~\u003e 0.12\", [hex: :protobuf, repo: \"hexpm\", optional: false]}, {:skogsra, \"~\u003e 2.2\", [hex: :skogsra, repo: \"hexpm\", optional: false]}], \"hexpm\", \"8470ed2f116fa8c9d70845f295f3738529aff123d9d3f98cbfa37107314763cb\"},\n   \"bunt\": {:hex, :bunt, \"1.0.0\", \"081c2c665f086849e6d57900292b3a161727ab40431219529f13c4ddcf3e7a44\", [:mix], [], \"hexpm\", \"dc5f86aa08a5f6fa6b8096f0735c4e76d54ae5c9fa2c143e5a1fc7c1cd9bb6b5\"},\n-  \"castore\": {:hex, :castore, \"1.0.19\", \"6903cabdfd9d1af46454126e7c8385186659dd33ecfb74a885cae52221ad6109\", [:mix], [], \"hexpm\", \"3669e6cab13f54c2df26b3e6833745d647f35b6e30d8ddd5975df0d5c842ca98\"},\n+  \"castore\": {:hex, :castore, \"1.0.20\", \"455e48f7115eca98c9f2b0e7a152b5a2e8f2a8a4f964c96e95bd31645ee5fa59\", [:mix], [], \"hexpm\", \"940eafbfd8b14bee649f083bc11b3b54ec555b54c3e4ea8213351ff6fee39c10\"},\n   \"certifi\": {:hex, :certifi, \"2.15.0\", \"0e6e882fcdaaa0a5a9f2b3db55b1394dba07e8d6d9bcad08318fb604c6839712\", [:rebar3], [], \"hexpm\", \"b147ed22ce71d72eafdad94f055165c1c182f61a2ff49df28bcc71d1d5b94a60\"},\n   \"con_cache\": {:hex, :con_cache, \"1.1.1\", \"9f47a68dfef5ac3bbff8ce2c499869dbc5ba889dadde6ac4aff8eb78ddaf6d82\", [:mix], [{:telemetry, \"~\u003e 1.0\", [hex: :telemetry, repo: \"hexpm\", optional: false]}], \"hexpm\", \"1def4d1bec296564c75b5bbc60a19f2b5649d81bfa345a2febcc6ae380e8ae15\"},\n   \"cors_plug\": {:hex, :cors_plug, \"2.0.3\", \"316f806d10316e6d10f09473f19052d20ba0a0ce2a1d910ddf57d663dac402ae\", [:mix], [{:plug, \"~\u003e 1.8\", [hex: :plug, repo: \"hexpm\", optional: false]}], \"hexpm\", \"ee4ae1418e6ce117fc42c2ba3e6cbdca4e95ecd2fe59a05ec6884ca16d469aea\"},"}
{"filename":"apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex","patch":"@@ -343,12 +343,22 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.ControlHandler do\n \n         {:ok, new_state}\n \n+      {:error, :interface_loading_failed} -\u003e\n+        # Track interface loading failure\n+        :telemetry.execute(\n+          [:astarte, :data_updater_plant, :control_handler, :properties_resend],\n+          %{duration: System.monotonic_time() - resend_start},\n+          %{realm: state.realm, result: \"interface_loading_failed\"}\n+        )\n+\n+        interface_loading_error(state, message_id, timestamp)\n+\n       {:error, :sending_properties_to_interface_failed} -\u003e\n         # Track interface send failure\n         :telemetry.execute(\n           [:astarte, :data_updater_plant, :control_handler, :properties_resend],\n           %{duration: System.monotonic_time() - resend_start},\n-          %{realm: state.realm, result: \"interface_failed\"}\n+          %{realm: state.realm, result: \"sending_properties_to_interface_failed\"}\n         )\n \n         sending_properties_error(state, message_id, timestamp)\n@@ -386,8 +396,31 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.ControlHandler do\n     new_state\n   end\n \n+  defp interface_loading_error(state, message_id, timestamp) do\n+    Logger.warning(\"Cannot load interface while resending properties\",\n+      tag: \"interface_loading_failed\"\n+    )\n+\n+    {:ok, new_state} = Core.Device.ask_clean_session(state, timestamp)\n+    MessageTracker.discard(new_state.message_tracker, message_id)\n+\n+    :telemetry.execute(\n+      [:astarte, :data_updater_plant, :data_updater, :discarded_message],\n+      %{},\n+      %{realm: new_state.realm}\n+    )\n+\n+    Core.Trigger.execute_device_error_triggers(\n+      new_state,\n+      \"interface_loading_failed\",\n+      timestamp\n+    )\n+\n+    new_state\n+  end\n+\n   defp sending_properties_error(state, message_id, timestamp) do\n-    Logger.warning(\"Cannot resend properties to interface\",\n+    Logger.warning(\"Cannot send properties to device.\",\n       tag: \"resend_interface_properties_failed\"\n     )\n "}
{"filename":"apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/device.ex","patch":"@@ -424,9 +424,12 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.Device do\n            :ok \u003c- resend_all_interface_properties(new_state, interface_descriptor) do\n         {:cont, {:ok, new_state}}\n       else\n-        {:error, :interface_loading_failed} -\u003e\n-          Logger.warning(\"Failed #{interface} interface loading.\")\n-          {:halt, {:error, :sending_properties_to_interface_failed}}\n+        {:error, reason} = error -\u003e\n+          Logger.warning(\n+            \"Error while sending properties to device: #{inspect(reason)} on interface #{interface}\"\n+          )\n+\n+          {:halt, error}\n       end\n     end)\n   end\n@@ -437,14 +440,28 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.Device do\n        ) do\n     encoded_device_id = Device.encode_device_id(device_id)\n \n-    Core.Interface.each_interface_mapping(mappings, interface_descriptor, fn mapping -\u003e\n-      resend_interface_mapping_properties(\n-        realm,\n-        device_id,\n-        encoded_device_id,\n-        interface_descriptor,\n-        mapping\n-      )\n+    mappings\n+    |\u003e Enum.filter(fn {_endpoint_id, mapping} -\u003e\n+      mapping.interface_id == interface_descriptor.interface_id\n+    end)\n+    |\u003e Enum.reduce_while(:ok, fn {_endpoint_id, mapping}, :ok -\u003e\n+      case resend_interface_mapping_properties(\n+             realm,\n+             device_id,\n+             encoded_device_id,\n+             interface_descriptor,\n+             mapping\n+           ) do\n+        :ok -\u003e\n+          {:cont, :ok}\n+\n+        {:error, reason} -\u003e\n+          Logger.warning(\n+            \"Failed sending properties for interface #{interface_descriptor.name}: #{inspect(reason)}.\"\n+          )\n+\n+          {:halt, {:error, :sending_properties_to_interface_failed}}\n+      end\n     end)\n   end\n "}
{"filename":"apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/control_handler_test.exs","patch":"@@ -122,6 +122,35 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.ControlHandlerTest do\n         {:ok, %{local_matches: 1, remote_matches: 0}}\n       end)\n \n+      expect(Core.Device, :resend_all_properties, fn _state -\u003e\n+        {:error, :interface_loading_failed}\n+      end)\n+\n+      expect(MessageTracker, :discard, fn ^message_tracker, ^message_id -\u003e :ok end)\n+      expect(Core.Device, :ask_clean_session, fn _state, _timestamp -\u003e {:ok, state} end)\n+\n+      expect(Core.Trigger, :execute_device_error_triggers, fn _state,\n+                                                              \"interface_loading_failed\",\n+                                                              _ts -\u003e\n+        :ok\n+      end)\n+\n+      ControlHandler.handle_control(\n+        state,\n+        \"/emptyCache\",\n+        \"\",\n+        message_id,\n+        0\n+      )\n+    end\n+\n+    test \"discards the message if sending properties to the device fails\", context do\n+      %{state: state, message_tracker: message_tracker, message_id: message_id} = context\n+\n+      Mox.expect(ClientMock, :publish, fn _data -\u003e\n+        {:ok, %{local_matches: 1, remote_matches: 0}}\n+      end)\n+\n       expect(Core.Device, :resend_all_properties, fn _state -\u003e\n         {:error, :sending_properties_to_interface_failed}\n       end)"}
{"filename":"apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/device_test.exs","patch":"@@ -232,7 +232,7 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.DeviceTest do\n         {:error, :interface_loading_failed}\n       end)\n \n-      assert {:error, :sending_properties_to_interface_failed} ==\n+      assert {:error, :interface_loading_failed} ==\n                Core.Device.resend_all_properties(state)\n     end\n \n@@ -268,6 +268,19 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.DeviceTest do\n \n       Core.Device.resend_all_properties(state)\n     end\n+\n+    test \"returns error when sending properties to the device fails\", context do\n+      %{state: state, server_property_with_all_endpoint_types: interface} = context\n+\n+      state = put_in(state.introspection, %{interface.name =\u003e interface.major_version})\n+\n+      Mimic.expect(VMQPlugin, :publish, 1, fn _topic, _bson, _qos -\u003e\n+        {:ok, %{local_matches: 0, remote_matches: 0}}\n+      end)\n+\n+      assert {:error, :sending_properties_to_interface_failed} ==\n+               Core.Device.resend_all_properties(state)\n+    end\n   end\n \n   describe \"set_device_disconnected/2\" do"}
{"filename":"apps/astarte_pairing/mix.lock","patch":"@@ -6,7 +6,7 @@\n   \"b58\": {:hex, :b58, \"1.0.3\", \"d300d6ae5a3de956a54b9e8220e924e4fee1a349de983df2340fe61e0e464202\", [:mix], [], \"hexpm\", \"af62a98a8661fd89978cf3a3a4b5b2ebe82209de6ac6164f0b112e36af72fc59\"},\n   \"bcrypt_elixir\": {:hex, :bcrypt_elixir, \"2.3.1\", \"5114d780459a04f2b4aeef52307de23de961b69e13a5cd98a911e39fda13f420\", [:make, :mix], [{:comeonin, \"~\u003e 5.3\", [hex: :comeonin, repo: \"hexpm\", optional: false]}, {:elixir_make, \"~\u003e 0.6\", [hex: :elixir_make, repo: \"hexpm\", optional: false]}], \"hexpm\", \"42182d5f46764def15bf9af83739e3bf4ad22661b1c34fc3e88558efced07279\"},\n   \"bunt\": {:hex, :bunt, \"1.0.0\", \"081c2c665f086849e6d57900292b3a161727ab40431219529f13c4ddcf3e7a44\", [:mix], [], \"hexpm\", \"dc5f86aa08a5f6fa6b8096f0735c4e76d54ae5c9fa2c143e5a1fc7c1cd9bb6b5\"},\n-  \"castore\": {:hex, :castore, \"1.0.19\", \"6903cabdfd9d1af46454126e7c8385186659dd33ecfb74a885cae52221ad6109\", [:mix], [], \"hexpm\", \"3669e6cab13f54c2df26b3e6833745d647f35b6e30d8ddd5975df0d5c842ca98\"},\n+  \"castore\": {:hex, :castore, \"1.0.20\", \"455e48f7115eca98c9f2b0e7a152b5a2e8f2a8a4f964c96e95bd31645ee5fa59\", [:mix], [], \"hexpm\", \"940eafbfd8b14bee649f083bc11b3b54ec555b54c3e4ea8213351ff6fee39c10\"},\n   \"cbor\": {:hex, :cbor, \"1.0.1\", \"39511158e8ea5a57c1fcb9639aaa7efde67129678fee49ebbda780f6f24959b0\", [:mix], [], \"hexpm\", \"5431acbe7a7908f17f6a9cd43311002836a34a8ab01876918d8cfb709cd8b6a2\"},\n   \"certifi\": {:hex, :certifi, \"2.15.0\", \"0e6e882fcdaaa0a5a9f2b3db55b1394dba07e8d6d9bcad08318fb604c6839712\", [:rebar3], [], \"hexpm\", \"b147ed22ce71d72eafdad94f055165c1c182f61a2ff49df28bcc71d1d5b94a60\"},\n   \"cfxxl\": {:git, \"https://github.com/ispirata/cfxxl.git\", \"98dc50b1cfe5a682b051b38b83cebc644bc08488\", []},"}

$ cd /root/astrate-mule && rg -n "resendServerProperties\|interface_loading_failed\|resend_interface_properties_failed\|emptyCache\|resendServer" internal/engine/*.go
(no output)

Now I have the upstream v1.3.3 change fully understood. The change (#2119) fixes the emptyCache server-property resend path in the data-updater-plant's `control_handler.ex` to emit a **distinct device-error trigger** depending on the failure:

- **`interface_loading_failed`** when a stored property's interface fails to load (this was previously remapped to `resend_interface_properties_failed`).
- **`resend_interface_properties_failed`** when sending the properties to the device itself fails.

This is exactly what the companion probe task `probe-emptycache-resend-device-error` asks about: "does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure?"

So the compat-note v1.3.3 task's deliverable is the wording for `docs/COMPATIBILITY.md` reflecting the v1.3.3 upstream change. Let me look at Astrate's `resendServerProperties` (in control.go) to understand the current behaviour and see whether the COMPATIBILITY.md doc has a section about device_error on the emptyCache resend.

Let me look at Astrate's control.go to understand where the relevant doc wording would go.
$ cd /root/astrate-mule && rg -n "device_error\|interface_loading_failed\|resend_interface_properties_failed\|UpstreamErrorName\|emptyCache\|server property\|server-property" docs/COMPATIBILITY.md
(no output)

$ cd /root/astrate-mule && rg -rn "resendServerProperties\|emptyCache" internal/
(no output)

There's no `resendServerProperties` in internal/engine/*.go per that grep. Let me search more broadly.
$ cd /root/astrate-mule && rg -rln "resend.*[Pp]roperties\|empty_cache\|EmptyCache" internal/ | head
(no output)

Interesting — the probe task references `internal/engine/control.go resendServerProperties`, but there's no such function in the current tree. Let me check what control.go actually contains and how the emptyCache control is handled currently.
$ cd /root/astrate-mule && ls internal/engine/ && echo "---" && rg -rn "emptyCache\|empty_cache\|producer/properties\|consumer/properties\|control" internal/engine/ -l
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
engine.go
engine_e2e_test.go
engine_integration_test.go
fakestore_test.go
forward
introspection.go
introspection_test.go
router.go
router_test.go
serverdata.go
serverdata_test.go
stream
testdata
topics.go
topics_test.go
triggers
---

✱ Grep "emptyCache|empty_cache|server property|serverProperty|resend" in internal/engine · 50 matches
Let me read `internal/engine/control.go` around `resendServerProperties` to see how failures are currently handled (whether device_error triggers are emitted at all).
→ Read internal/engine/control.go [offset=130, limit=60]
→ Read internal/engine/triggers/errorname.go
Now I understand the current Astrate state. In Astrate's `resendServerProperties`, failures are only logged and skipped (lines 151-173), and the `broker.Publish` failure propagates the error up (line 176). **No `device_error` trigger event is emitted** during the emptyCache resend — the failure cases are just `Warn`/`Error` logs with `continue`.

So the v1.3.3 upstream change (#2119) is a wire-visible behavioural requirement that Astrate currently does **not** implement: upstream now emits a distinct `device_error` trigger on the emptyCache property-resend path — `interface_loading_failed` when a stored property's interface is unloadable, and `resend_interface_properties_failed` when sending to the device fails. Astrate currently only logs and skips.

This is exactly what the companion probe `probe-emptycache-resend-device-error` covers (report, do not patch), and this compat-note task asks me to propose the `docs/COMPATIBILITY.md` wording for v1.3.3.

Let me check how COMPATIBILITY.md already handles the device_error topic (deviation 1, sub-point about `device_error.error_name` enum). The wording for v1.3.3 should note this new distinct-error behaviour. Let me check whether COMPATIBILITY.md currently mentions the emptyCache resend at all, and verify there's no reference.

The COMPATIBILITY.md's deviation 1 last sub-point discusses `device_error.error_name` enum (measured against upstream v1.2.0 in `channels.json`). The v1.3.3 change adds a distinct `interface_loading_failed` vs `resend_interface_properties_failed` on the emptyCache resend — a wire-visible trigger payload difference.

Now, what should the wording be? Since this is a *proposal* task (do not edit the file), the deliverable is a `.mule/for-giulio.md` entry. But the task is `[auto]` and produces no code change — wait, let me re-read.

The task: "compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]"

So the deliverable is a `.mule/for-giulio.md` entry proposing the wording. This is like the `[readonly]`/milestone case — a docs proposal, produces no code diff, and the for-giulio.md entry is the deliverable.

But wait — `.mule/for-giulio.md` is it a file I can edit? Let me check. MULE.md's never-touch list: `docs/DESIGN.md`, `docs/ROADMAP.md`, `migrations/`, `.github/`, `go.mod`, `go.sum`, `Dockerfile`, `docker-compose.yml`, `.mule/milestones.md`, `.trickle/`. `.mule/for-giulio.md` is NOT on the never-touch list — in fact the escalation channel explicitly says "The mule writes here". And the previous `compat-note-v1.4.0-rc.3` task was blocked, but the intent is clearly to append a proposal line to `.mule/for-giulio.md`.

Wait, let me re-check — is editing `.mule/for-giulio.md` the right deliverable? Let me re-read the recipe. The task explicitly says "propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file)". "do not edit the file" refers to `docs/COMPATIBILITY.md`. So yes, I write the proposal into `.mule/for-giulio.md`.

But wait — the compat-note task produces no code diff. Let me check: does the runner require a test failing gate for this? The MULE.md says for `[readonly]` and milestone-issue tasks, the deliverable is the report itself, "leaves the working tree clean". For this `[auto]` compat-note task, the deliverable is writing to `.mule/for-giulio.md`. That's a working-tree change but it's a docs proposal, like for-giulio.escalations.

Let me look at how the previous (now blocked) compat-note and for-giulio entries look, and whether there's precedent for what "COMPATIBILITY.md wording" means. Let me check the for-giulio.md format — one line each, newest at top. The entry should contain the evidence inline.

Let me confirm the current COMPATIBILITY.md's stated upstream target and whether v1.3.3 needs to be reflected there. The doc says deviation 1 mentions "measured, not reconstructed: recorded against upstream Astarte v1.2.0". The conformance matrix pins client SDK versions, not upstream server version. The upstream server version appears only in prose.

So the v1.3.3 compat-note wording should be a deviation/prose note about the emptyCache resend now distinguishing `interface_loading_failed` from `resend_interface_properties_failed` (v1.3.3 #2119) — and the current Astrate state (logs and skips, emits no device_error on this path).

Let me verify what the current COMPATIBILITY.md says about emptyCache resend's device_error — checking if there's any existing mention of empty_cache_error or the resend path in the doc.

There is none (grep above returned empty). Let me also double-check the probe task's framing to align my wording. Let me look at the `handleEmptyCache` function fully, and how `def`/device_error events are fired elsewhere (e.g., is there any device_error on the resend path).

I have enough context now. Let me also verify the CHANGELOG wording to base my proposal on upstream's own language:

Upstream CHANGELOG v1.3.3: "Resending properties to a device now reports distinct device error triggers depending on the failure: `interface_loading_failed` when an interface fails to load and `resend_interface_properties_failed` when sending the properties to the device itself fails. Previously, both cases were incorrectly reported as `resend_interface_properties_failed`."

Now let me write the proposal to `.mule/for-giulio.md`. The entry should propose the COMPATIBILITY.md wording for v1.3.3. Since the requirement is "propose the wording", I'll append a line (newest at top per the file convention) that names the wording I'd propose and notes Astrate's current non-conforming state (which is the subject of the separate probe task).

Let me look at the exact top of for-giulio.md again to match the insertion point (newest at top, right after the `---` separator at line 11).

The format: one line each, newest at the top, with evidence inline. Let me compose a concise proposal line. Given "propose the docs/COMPATIBILITY.md wording", I should give the actual wording text. This might warrant the verdict-style like other decided/undecided lines. Since it's a proposal (not decided), I'll present it as the proposed wording with evidence.

Let me write the proposal. I'll insert it right after the first `---` (line 11).

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Proposed COMPATIBILITY.md wording for upstream v1.3.3** — upstream v1.3.3 (#2119, released 2026-08-07) makes the emptyCache server-property resend report a **distinct device-error trigger per failure**: `interface_loading_failed` when a stored property's interface will not load, `resend_interface_properties_failed` when publishing to the device itself fails (previously both were remapped to `resend_interface_properties_failed`). Suggested deviation-1 wording: "the emptyCache property resend (`internal/engine/control.go` `resendServerProperties`) is a logged-and-skipped path, not a device-error trigger source: an unloadable interface, a path that matches no endpoint, a failed rehydrate and a failed encode are each logged and skipped rather than surfaced as `interface_loading_failed`, and a publish failure propagates to the retry loop instead of emitting `resend_interface_properties_failed` — upstream v1.3.3 (#2119) distinguishes these two device-error triggers, so this surface is not yet measured parity." (compat-note v1.3.3; Astrate's gap is the separate `probe-emptycache-resend-device-error` probe.)
+
+---
+
 - ~~`device_deletion_started`/`device_deletion_finished` trigger events are not emitted~~ —
   **decided 2026-07-27: emit both, back-to-back, around the synchronous delete.** Filed as
   issue #21 (`mule`). (Cross-project survey, 2026-07-27,
   `.mule/research/survey-2026-07-27.md` source 4.)

$ cd /root/astrate-mule && git status --porcelain
M .mule/for-giulio.md
?? .mule/cron.log
?? .mule/tasks/realm-config-datastream-retention.md

Now the proposal is added. This task produces no code change and no test is applicable (it's a docs-proposal deliverable; the doc file itself is never touched, and `.mule/for-giulio.md` is the escalation channel). I don't need to run the Go gates since I changed no Go code. Let me verify the edit and confirm the working tree state.
→ Read .mule/for-giulio.md [limit=16]
The proposal is in place. The only file I changed is `.mule/for-giulio.md` (the escalation channel I'm meant to write to). The other untracked files (`.mule/cron.log`, `.mule/tasks/realm-config-datastream-retention.md`) were pre-existing, not mine — I leave them alone.

Since I changed no Go code, the build gates don't apply. Let me confirm the report format. The task is done: proposed the wording in `.mule/for-giulio.md`, did not edit `docs/COMPATIBILITY.md`.

No test applies — this is a docs-proposal deliverable, analogous to a `milestone-...-issue` task where the deliverable is the proposed wording itself, and the runner would not expect a breaking test for it (there's no behaviour change to gate).

```
Done: Proposed docs/COMPATIBILITY.md wording for upstream v1.3.3 (#2119 — distinct interface_loading_failed vs resend_interface_properties_failed on the emptyCache resend) in .mule/for-giulio.md; did not edit the file.
Files: .mule/for-giulio.md
Verified: (no Go code changed, so no gate run; docs-only proposal, source = upstream v1.3.2..v1.3.3 compare + PR #2119)
Unsure: Astrate's exact gap phrasing is the probe-emptycache-resend-device-error task's to confirm; my COMPATIBILITY.md wording states the current logged-and-skipped behaviour as read from control.go.
Follow-ups: none
```
