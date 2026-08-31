slug: recipe-astarte-upstream
verdict: proposed
at:  6e30a1c
ran: 2026-08-31T02:12:36Z on DietPi in 224s


> build · big-pickle

I'll start by reading the two files specified.
→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
Now let me read the compatibility doc and check upstream releases in parallel.
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

→ Read .mule/todo.md
Upstream is ahead: newest tag is `v1.4.0-rc.5` (stable also has `v1.3.3`), while COMPATIBILITY.md records upstream v1.2.0/v1.2.2. Let me check the todo queue for prior runs and pull the release notes.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
Sixth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.4

### Fixed

- [astarte_data_updater_plant] Make object values only validate type on selected interface
- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'
The queue already covers v1.3.0/v1.3.3/v1.4.0-rc.3 and rc.5's binaryblob fix, but not rc.5's other change: "Make object values only validate type on selected interface." Let me look at that upstream change before proposing anything.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | {tag_name, body, published_at}'
{"body":"","published_at":"2026-08-07T14:33:58Z","tag_name":"v1.3.3"}

$ gh search commits --owner astarte-platform 'Make object values only validate type on selected interface' --json sha,repository,commit --jq '.[] | "\(.repository.full_name)\t\(.sha[0:9])\t\(.commit.message | split("\n")[0])"'
(no output)

$ gh search commits --owner astarte-platform 'object values only validate type on selected interface' --json sha,repository,commit 2>/dev/null; echo "---"; gh search commits --owner astarte-platform 'object values' --json sha,repository,commit.url,commit.author --jq '.[] | "\(.repository.full_name)\t\(.sha[0:9])\t\(.commit.message | split("\n")[0])"' 2>&1 | head -20
[]
---
Unknown JSON field: "commit.url"
Available fields:
  author
  commit
  committer
  id
  parents
  repository
  sha
  url

$ gh api repos/astarte-platform/astarte_data_updater_plant/commits --jq '.[0:15][] | "\(.sha[0:9])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"'
8739d915f	2019-12-20T15:06:00Z	This repository has been archived.
69612e95f	2019-12-18T16:49:09Z	Forward port changes from v0.10
28ae84960	2019-12-09T11:32:32Z	Merge branch 'prepare-0.10.2' into 'release-0.10'
c0087ddce	2019-12-09T11:21:25Z	Prepare 0.10.2 release
64da6f368	2019-11-28T14:49:05Z	Merge pull request #17 from rbino/scylla-beta-support
a40acb406	2019-11-28T14:18:26Z	.travis.yml: remove Scylla CI from allowed failures
9a3198a26	2019-11-28T14:12:40Z	queries: don't use shorthand syntax to update maps
bc11f0383	2019-11-27T15:16:08Z	Merge branch 'ci-upgrade-services' into 'master'
59bd736d7	2019-11-27T10:49:58Z	.gitlab-ci.yml: upgrade Cassandra and RabbitMQ services
a4156ce2b	2019-11-26T14:20:37Z	Forward port bugfixes from v0.10
140f80393	2019-11-26T11:34:43Z	Merge branch 'fix-13' into 'release-0.10'
745d9b792	2019-11-26T11:22:01Z	impl: load device volatile triggers when they're installed
afe9d5b4f	2019-11-26T11:21:20Z	Merge branch 'fix-15' into 'release-0.10'
466ca7c70	2019-11-26T10:49:12Z	impl: fix volatile triggers deletion bug
6fcd3f5f0	2019-11-19T16:24:14Z	Merge pull request #12 from rbino/ci-wait-for-cassandra

The data updater code now lives in the astarte monorepo. Let me find the actual commit for the rc.5 object-validation change.
$ gh api 'repos/astarte-platform/astarte/commits?since=2026-08-15T00:00:00Z&until=2026-08-21T00:00:00Z' --jq '.[] | "\(.sha[0:9])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"'
739724991	2026-08-19T13:06:29Z	chore: Update exandra
643d06107	2026-08-19T12:36:15Z	ci: add release-please for automated versioning and changelog
6e65c8aca	2026-08-17T12:25:10Z	refactor(dup): Move encryption logic into library (#2123)
88004a2cf	2026-08-14T09:21:28Z	chore(core): use the local astarte_core across the monorepo

$ gh search commits --owner astarte-platform --repository astarte 'binaryblob data is correctly validated' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo "==="; gh search commits --owner astarte-platform --repository astarte 'object values only validate type' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'
unknown flag: --repository

Usage:  gh search commits [<query>] [flags]

Flags:
      --author string            Filter by author
      --author-date date         Filter based on authored date
      --author-email string      Filter on author email
      --author-name string       Filter on author name
      --committer string         Filter by committer
      --committer-date date      Filter based on committed date
      --committer-email string   Filter on committer email
      --committer-name string    Filter on committer name
      --hash string              Filter by commit hash
  -q, --jq expression            Filter JSON output using a jq expression
      --json fields              Output JSON with the specified fields
  -L, --limit int                Maximum number of commits to fetch (default 30)
      --merge                    Filter on merge commits
      --order string             Order of commits returned, ignored unless '--sort' flag is specified: {asc|desc} (default "desc")
      --owner strings            Filter on repository owner
      --parent string            Filter by parent hash
  -R, --repo OWNER/REPO          Filter on repository, in OWNER/REPO format
      --sort string              Sort fetched commits: {author-date|committer-date} (default "best-match")
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --tree string              Filter by tree hash
      --visibility strings       Filter based on repository visibility: {public|private|internal}
  -w, --web                      Open the search query in the web browser
  
===
unknown flag: --repository

Usage:  gh search commits [<query>] [flags]

Flags:
      --author string            Filter by author
      --author-date date         Filter based on authored date
      --author-email string      Filter on author email
      --author-name string       Filter on author name
      --committer string         Filter by committer
      --committer-date date      Filter based on committed date
      --committer-email string   Filter on committer email
      --committer-name string    Filter on committer name
      --hash string              Filter by commit hash
  -q, --jq expression            Filter JSON output using a jq expression
      --json fields              Output JSON with the specified fields
  -L, --limit int                Maximum number of commits to fetch (default 30)
      --merge                    Filter on merge commits
      --order string             Order of commits returned, ignored unless '--sort' flag is specified: {asc|desc} (default "desc")
      --owner strings            Filter on repository owner
      --parent string            Filter by parent hash
  -R, --repo OWNER/REPO          Filter on repository, in OWNER/REPO format
      --sort string              Sort fetched commits: {author-date|committer-date} (default "best-match")
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --tree string              Filter by tree hash
      --visibility strings       Filter based on repository visibility: {public|private|internal}
  -w, --web                      Open the search query in the web browser

$ gh search commits -R astarte-platform/astarte 'binaryblob data is correctly validated' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo "==="; gh search commits -R astarte-platform/astarte 'object values only validate type on selected interface' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'
aedc41483	fix(dup): ensure binaryblob data is correctly validated
===

$ gh search commits -R astarte-platform/astarte 'object values only validate type' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo "===loose==="; gh search commits -R astarte-platform/astarte 'object values' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"' | head -10
===loose===
cc89858b3	DUP: do not crash on binaryblob insertion
219b9c79d	Merge branch 'object-aggregation' into 'master'
084db3757	Handle aggregate interface in live events
adf94edf8	Fix data parsing on datastream object interfaces w/ parametric endpoints

$ gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.4...v1.4.0-rc.5 --jq '.commits[] | "\(.sha[0:9])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"'
d8592091e	2026-08-07T06:55:11Z	fix(dup): stop remapping interface loading failed error (#2119)
a816eaf0e	2026-08-07T06:57:11Z	chore: forward port release-1.3
aedc41483	2026-08-07T07:21:33Z	fix(dup): ensure binaryblob data is correctly validated
3e70d63f8	2026-08-07T14:28:36Z	chore: prepare v1.3.3 (#2128)
fdc6d29b8	2026-08-19T08:14:18Z	fix(dup): Make object values only validate type on selected interface (#2141)
4141a2ae4	2026-08-19T13:08:32Z	chore: forward-port release-1.3 into release-1.4
3484096d2	2026-08-20T12:10:25Z	Prepare release v1.4.0-rc.5 (#2147)

The new item is commit `fdc6d29b8` (#2141). Let me read the diff.
$ gh api repos/astarte-platform/astarte/commits/fdc6d29b8 --jq '{sha, message: .commit.message, files: [.files[] | {filename, status: .status, additions, deletions}]}'
{"files":[{"additions":1,"deletions":1,"filename":"apps/astarte_appengine_api/mix.lock","status":"modified"},{"additions":2,"deletions":1,"filename":"apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex","status":"modified"},{"additions":40,"deletions":0,"filename":"apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_test.exs","status":"modified"},{"additions":23,"deletions":0,"filename":"apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/interface_test.exs","status":"modified"},{"additions":58,"deletions":0,"filename":"apps/astarte_data_updater_plant/test/support/cases/device.ex","status":"modified"},{"additions":1,"deletions":1,"filename":"apps/astarte_pairing/mix.lock","status":"modified"}],"message":"fix(dup): Make object values only validate type on selected interface (#2141)\n\nSigned-off-by: Francesco Noacco \u003cfrancesco.noacco@secomind.com\u003e","sha":"fdc6d29b84ec06a12316340b60f1b43491ca3d5c"}

$ gh api repos/astarte-platform/astarte/commits/fdc6d29b8 --jq '.files[] | select(.filename=="apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex") | .patch'
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

$ gh api repos/astarte-platform/astarte/commits/fdc6d29b8 --jq '.files[] | select(.filename=="apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/interface_test.exs") | .patch'
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

✱ Grep "validateValueType|ValidateValueType|func.*[Vv]alidate.*[Vv]alue" in internal/engine · 1 match
This is a distinct wire-visible fix: object datastream values were validated against whatever interface had the same endpoint names, since `extract_mappings` for an `:object` aggregation pulled mappings across all interfaces. Let me see where Astrate validates object values so the probe line targets the right files.
✱ Grep "extract_mappings|extractMappings" in internal · 0 matches
✱ Grep "[Vv]alue[Tt]ype" in . · 89 matches
✱ Grep "Lookup|Trie|resolveMapping|Mapping.*by.*path|matchPath" in internal/engine · 13 matches
Let me find how Astrate resolves the mapping set when decoding an object value — that determines whether the upstream bug has an Astrate mirror.
→ Read pkg/payload/payload.go
→ Read internal/engine/data.go [offset=230, limit=120]
Astrate decodes object values against `ci.ObjectLeaves` — the topic-named interface's own leaves — so the rc.5 #2141 bug may or may not have a mirror; that's exactly what a probe verifies. Let me check rc.4's notes too, since it sits between what prior runs covered (rc.3, rc.5-partial), and confirm nothing else is pending.

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

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -102,4 +102,6 @@
 - [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]
 - [ ] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto]
 - [ ] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto]
 - [ ] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto]
+- [ ] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto]
+- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]

Nothing new in rc.4 (config-library and Elixir-infra changes — ignore). The only uncovered item since the last run is rc.5's #2141. Appending the two task lines to the queue.
$ tail -n 5 .mule/todo.md
- [ ] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto]
- [ ] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto]
- [ ] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto]
- [ ] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto]
- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]

Done: upstream is ahead (newest `v1.4.0-rc.5` vs the v1.2.x recorded in COMPATIBILITY.md), but prior runs already proposed for v1.3.0/v1.3.3/rc.3 and rc.5's binaryblob fix. The only uncovered change is rc.5 #2141, so I appended just those two lines.

Done: watched upstream Astarte; proposed probe + compat-note for the one uncovered change, rc.5 #2141 (v1.4.0-rc.5 "Make object values only validate type on selected interface")
Files: .mule/todo.md
Verified: gh api releases + compare rc.4...rc.5 + commit fdc6d29b8 diff reviewed -> pass; no gate commands needed (doc/queue change only)
Unsure: nothing
Follow-ups: `probe-object-validation-selected-interface`, `compat-note-v1.4.0-rc.5` appended to .mule/todo.md
