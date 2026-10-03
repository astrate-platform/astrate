slug: appengine-missing-required-422
verdict: done
at:  f84bbb5
ran: 2026-10-03T19:50:23Z on DietPi in 584s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ rg -n "ReasonMissingRequired|ReasonBadObject|missing_required" --glob '*.go' | head -40
pkg/payload/required_test.go:13:// declares required is rejected with ReasonMissingRequired in both wire
pkg/payload/required_test.go:38:	if _, err := DecodeObject([]byte(`{"v":{"lon":9.0}}`), leaves); ReasonOf(err) != ReasonMissingRequired {
pkg/payload/required_test.go:39:		t.Fatalf("JSON rejection reason = %v (err %v); want %v", ReasonOf(err), err, ReasonMissingRequired)
pkg/payload/required_test.go:47:	if _, err := DecodeObject(bsonBad, leaves); ReasonOf(err) != ReasonMissingRequired {
pkg/payload/required_test.go:48:		t.Fatalf("BSON rejection reason = %v (err %v); want %v", ReasonOf(err), err, ReasonMissingRequired)
pkg/payload/value.go:145:	// ReasonBadObject: object-aggregation shape violation (not a
pkg/payload/value.go:147:	ReasonBadObject
pkg/payload/value.go:151:	// ReasonMissingRequired: an object-aggregation document omits a key
pkg/payload/value.go:153:	ReasonMissingRequired
pkg/payload/value.go:160:	"missing_required",
pkg/payload/value.go:165:	if r >= ReasonTooLarge && r <= ReasonMissingRequired {
pkg/payload/value.go:188:	// resolve to no declared leaf (ReasonBadObject only; nil otherwise).
pkg/payload/value.go:213:		Reason:         ReasonBadObject,
pkg/payload/payload.go:148:			return DecodedPayload{}, rejectf(ReasonMissingRequired,
pkg/payload/json.go:203:		return nil, rejectf(ReasonBadObject, "JSON value %s where an object-aggregation document was expected", clip(raw))
pkg/payload/json.go:210:		return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
pkg/payload/json.go:370:			return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
pkg/payload/json.go:375:				return nil, rejectf(ReasonBadObject, "key %q: nested objects are not encodable", k)
pkg/payload/payload_test.go:414:		`{"v":{"lat":1.0,"nope":2.0}}`: ReasonBadObject,    // undeclared key
pkg/payload/payload_test.go:415:		`{"v":{}}`:                     ReasonBadObject,    // empty document
pkg/payload/payload_test.go:416:		`{"v":5}`:                      ReasonBadObject,    // not a document
pkg/payload/payload_test.go:417:		`{"v":null}`:                   ReasonBadObject,    // null
pkg/payload/payload_test.go:418:		`{"v":[1]}`:                    ReasonBadObject,    // array
pkg/payload/payload_test.go:554:	if _, err := DecodeObject(dup, objectLeaves(false)); ReasonOf(err) != ReasonBadObject {
pkg/payload/payload_test.go:597:			if ReasonOf(err) != ReasonBadObject {
pkg/payload/payload_test.go:849:		if _, err := Encode(map[string]Value{}, nil, f); ReasonOf(err) != ReasonBadObject {
pkg/payload/payload_test.go:852:		if _, err := Encode(map[string]Value{"a": map[string]Value{"b": 1.0}}, nil, f); ReasonOf(err) != ReasonBadObject {
internal/appengine/http.go:684:	case payload.ReasonOf(err) == payload.ReasonBadObject:
pkg/payload/bson.go:140:		return nil, rejectf(ReasonBadObject, "BSON %s where an object-aggregation document was expected", rv.Type)
pkg/payload/bson.go:147:		return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
pkg/payload/bson.go:169:			return nil, rejectf(ReasonBadObject, "duplicate object key %q", key)
pkg/payload/bson.go:328:		return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
pkg/payload/bson.go:340:			return nil, rejectf(ReasonBadObject, "invalid object key %q", k)
pkg/payload/bson.go:343:			return nil, rejectf(ReasonBadObject, "key %q: nested objects are not encodable", k)
pkg/payload/value_test.go:196:		ReasonBadObject:       "bad_object",
pkg/payload/value_test.go:199:		ReasonMissingRequired: "missing_required",
internal/appengine/writeerror_test.go:35:			&payload.RejectError{Reason: payload.ReasonBadObject,
internal/appengine/writeerror_test.go:40:			fmt.Errorf("engine: %w", &payload.RejectError{Reason: payload.ReasonBadObject,
internal/appengine/writeerror_test.go:45:			&payload.RejectError{Reason: payload.ReasonBadObject, Detail: "object-aggregation document is empty"},
internal/appengine/writeerror_test.go:88:		Reason:         payload.ReasonBadObject,

$ ssh -o ConnectTimeout=5 legion 'echo ok' 2>&1 | tail -3
ok

The Legion Go answers. Reading the recipe and named files.
→ Read .mule/recipes/legion-go.md
→ Read internal/appengine/writeerror_test.go
→ Read pkg/payload/required_test.go
→ Read pkg/payload/payload.go [offset=120, limit=50]
→ Read internal/appengine/http.go [offset=630, limit=110]

$ ls /root/astrate-mule/.mule /root/astrate-mule/.mule/recipes /root/astrate-mule/.mule/tasks 2>&1
/root/astrate-mule/.mule:
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
recipes
reports
reviews
task.md
tasks
todo.md

/root/astrate-mule/.mule/recipes:
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

/root/astrate-mule/.mule/tasks:
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md

$ ssh -o ConnectTimeout=8 legion 'ls ~; echo ---; docker ps --format "{{.Names}}" 2>&1' 2>&1 | tail -30
drwxr-xr-x    - atsetilam 18 mar 01:00  Pubblici
drwxr-xr-x    - atsetilam 27 set 17:31  Scaricati
drwxr-xr-x    - atsetilam 20 lug 18:10  steamvr
drwxr-xr-x    - atsetilam  6 ago 02:42  test
drwxr-xr-x    - atsetilam  3 ott 16:16  toshy-src
drwxr-xr-x    - atsetilam  5 ago 17:13  unsloth-studio-test
drwxr-xr-x    - atsetilam  5 ago 16:06  unsloth_compiled_cache
drwxr-xr-x    - atsetilam 18 mar 01:00  Video
drwxr-xr-x    - atsetilam 30 mar 20:27  whisper.cpp
.rw-r--r--   21 atsetilam 18 mar 00:48 󱆃 .bash_logout
.rw-r--r--   57 atsetilam 18 mar 00:48 󱆃 .bash_profile
.rw-r--r--  172 atsetilam 18 mar 00:48 󱆃 .bashrc
.rw-------  74k atsetilam  3 ott 21:12  .claude.json
.rw-r--r--  133 atsetilam 21 lug 02:35 󰊢 .gitconfig
.rw-r--r--  348 atsetilam  3 ott 18:46  .gtkrc-2.0
.rw-------   30 atsetilam  6 ago 19:19  .npmrc
.rw-------  256 atsetilam 18 mar 01:00  .pulse-cookie
lrwxrwxrwx    - atsetilam  3 ott 16:45  .steampath -> /home/atsetilam/.steam/sdk32/steam
lrwxrwxrwx    - atsetilam  3 ott 16:45  .steampid -> /home/atsetilam/.steam/steam.pid
.rw-------  885 atsetilam 19 mar 20:11  .viminfo
.rw-r--r--  180 atsetilam 19 mar 20:20  .wget-hsts
.rw-r--r--  57k atsetilam  3 ott 16:21  .zcompdump-atsetilam-legion-5.9.2
.r--r--r-- 122k atsetilam  3 ott 16:21  .zcompdump-atsetilam-legion-5.9.2.zwc
.rw-r--r--   56 atsetilam 18 mar 00:48 󱆃 .zshrc
.rw-r--r--   64 atsetilam 20 lug 18:12  crash_log.txt
.rw-r--r--  22k atsetilam 30 mar 17:46  install-docker.sh
.rw-r--r--  19k atsetilam 21 lug 01:49  scale-gfx1103-llamacpp-report.md
.rw-r--r-- 1,4G atsetilam 18 lug 10:43  scale-latest-amd64.tar.xz
---
failed to connect to the docker API at unix:///home/atsetilam/.docker/desktop/docker.sock; check if the path is correct and if the daemon is running: dial unix /home/atsetilam/.docker/desktop/docker.sock: connect: no such file or directory

→ Read .mule/recipes/astarte-upstream.md

$ ssh -o ConnectTimeout=8 legion 'ls -d ~/*/ 2>/dev/null; echo "--- systemd docker:"; systemctl is-active docker 2>&1; echo "--- podman:"; command -v podman' 2>&1 | tail -25
drwxr-xr-x - atsetilam 26 set 00:29  /home/atsetilam/Games
drwxr-xr-x - atsetilam 27 lug 06:05  /home/atsetilam/go
drwxr-xr-x - atsetilam 21 lug 11:24  /home/atsetilam/Immagini
drwxr-xr-x - atsetilam 31 ago 00:44  /home/atsetilam/legion-gmind-backup-20260831-004404
drwxr-xr-x - atsetilam  4 ago 03:59  /home/atsetilam/listello
drwxr-xr-x - atsetilam  3 ago 18:09  /home/atsetilam/listello-wip
drwxr-xr-x - atsetilam  6 ago 02:57  /home/atsetilam/llama.cpp
drwxr-xr-x - atsetilam 18 mar 01:00  /home/atsetilam/Modelli
drwxr-xr-x - atsetilam 12 lug 07:16  /home/atsetilam/models
drwxr-xr-x - atsetilam 18 mar 01:00  /home/atsetilam/Musica
drwxr-xr-x - atsetilam 27 ago 03:38  /home/atsetilam/orange-cat
drwxr-xr-x - atsetilam 11 lug 18:30  /home/atsetilam/Progetti
drwxr-xr-x - atsetilam 21 lug 11:19  /home/atsetilam/prova-sito-phi
drwxr-xr-x - atsetilam 18 mar 01:00  /home/atsetilam/Pubblici
drwxr-xr-x - atsetilam 27 set 17:31  /home/atsetilam/Scaricati
drwxr-xr-x - atsetilam 20 lug 18:10  /home/atsetilam/steamvr
drwxr-xr-x - atsetilam  6 ago 02:42  /home/atsetilam/test
drwxr-xr-x - atsetilam  3 ott 16:16  /home/atsetilam/toshy-src
drwxr-xr-x - atsetilam  5 ago 17:13  /home/atsetilam/unsloth-studio-test
drwxr-xr-x - atsetilam  5 ago 16:06  /home/atsetilam/unsloth_compiled_cache
drwxr-xr-x - atsetilam 18 mar 01:00  /home/atsetilam/Video
drwxr-xr-x - atsetilam 30 mar 20:27  /home/atsetilam/whisper.cpp
--- systemd docker:
inactive
--- podman:

$ rg -n "b6d46ad4|missing_required_mapping|2237" -g '!*.sum' . 2>/dev/null | head -30
./pkg/payload/value.go:190:	// detail (master b6d46ad4, #2237), so the HTTP surfaces name the keys
./pkg/payload/value.go:208:// found (master b6d46ad4, #2237).
./pkg/payload/payload_test.go:561:// `unexpected_keys`, master b6d46ad4, #2237). JSON map iteration is random, so
./internal/engine/triggers/errorname.go:62:	// Upstream 1.4 emits missing_required_mapping, outside dashboard
./internal/engine/serverdata_test.go:318:// (master b6d46ad4, #2237). Nothing is persisted or published.
./internal/appengine/http.go:686:		// keys (master b6d46ad4, #2237); before that commit the shape was
./internal/appengine/http.go:698:// write carrying undeclared keys (master b6d46ad4, #2237). Frozen string:
./internal/appengine/writeerror_test.go:82:// an object write carrying an undeclared key (master b6d46ad4, #2237):

$ ssh -o ConnectTimeout=8 legion -T 'ls -d ~/*/ 2>/dev/null | head -20' 2>&1 | sed 's/\x1b\[[0-9;]*m//g'
drwxr-xr-x - atsetilam  3 ott 16:57  /home/atsetilam/agent-skills
drwxr-xr-x - atsetilam 25 lug 17:01  /home/atsetilam/Applications
drwxr-xr-x - atsetilam 22 ago 03:35  /home/atsetilam/astrate
drwxr-xr-x - atsetilam 11 lug 22:50 󱧼 /home/atsetilam/build
drwxr-xr-x - atsetilam 27 set 16:36  /home/atsetilam/Desktop
drwxr-xr-x - atsetilam  3 ott 17:50  /home/atsetilam/dev
drwxr-xr-x - atsetilam 18 mar 01:00  /home/atsetilam/Documenti
drwxr-xr-x - atsetilam  3 ott 16:44  /home/atsetilam/dot-steam.bak.1791038671
drwxr-xr-x - atsetilam 22 ago 03:56  /home/atsetilam/edgehog-test
drwxr-xr-x - atsetilam 19 mar 19:46  /home/atsetilam/emudeck
drwxr-xr-x - atsetilam 19 mar 20:26  /home/atsetilam/Emulation
drwxr-xr-x - atsetilam 19 mar 20:33  /home/atsetilam/ES-DE
drwxr-xr-x - atsetilam  3 ott 16:04  /home/atsetilam/g-mind
drwxr-xr-x - atsetilam 26 set 00:29  /home/atsetilam/Games
drwxr-xr-x - atsetilam 27 lug 06:05  /home/atsetilam/go
drwxr-xr-x - atsetilam 21 lug 11:24  /home/atsetilam/Immagini
drwxr-xr-x - atsetilam 31 ago 00:44  /home/atsetilam/legion-gmind-backup-20260831-004404
drwxr-xr-x - atsetilam  4 ago 03:59  /home/atsetilam/listello
drwxr-xr-x - atsetilam  3 ago 18:09  /home/atsetilam/listello-wip
drwxr-xr-x - atsetilam  6 ago 02:57  /home/atsetilam/llama.cpp

→ Read internal/engine/triggers/errorname.go [offset=40, limit=40]

$ gh search code --owner astarte-platform 'missing_required_mapping' --limit 20 2>&1 | head -30; echo "=== EXIT $?"
astarte-platform/astarte:apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex: def render("422_missing_required_mapping.json", _assigns) do
astarte-platform/astarte:apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex: tag: "missing_required_mapping"
astarte-platform/astarte:apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex: {:error, :missing_required_mapping}
astarte-platform/astarte:apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/fallback_controller.ex: def call(conn, {:error, :missing_required_mapping}) do
astarte-platform/astarte:apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/fallback_controller.ex: |> render(:"422_missing_required_mapping")
astarte-platform/astarte:apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/data_handler.ex: logger_metadata: [tag: "missing_required_mapping"],
astarte-platform/astarte:apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/data_handler.ex: error_name: "missing_required_mapping",
astarte-platform/astarte:apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/data_handler.ex: error: :missing_required_mapping
astarte-platform/astarte:apps/astarte_appengine_api/test/astarte_appengine_api/device/device_test.exs: ) == {:error, :missing_required_mapping}
astarte-platform/astarte:apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_error_test.exs: assert {:discard, :missing_required_mapping, ^state_with_required_mapping} =
=== EXIT 0

$ ssh -o ConnectTimeout=8 legion -T 'sudo -n true 2>&1 && echo sudo-ok; echo "--- astrate dir:"; ls ~/astrate | head -30; echo "--- compose:"; ls ~/astrate/*.y*ml 2>/dev/null' 2>&1 | sed 's/\x1b\[[0-9;]*m//g' | head -50
sudo: è necessaria una password
--- astrate dir:
drwxr-xr-x    - atsetilam 22 ago 03:35  .agents
drwxr-xr-x    - atsetilam 22 ago 03:39  .git
drwxr-xr-x    - atsetilam 26 lug 23:38  .github
drwxr-xr-x    - atsetilam 22 ago 03:35  .mule
drwxr-xr-x    - atsetilam 24 ago 20:25  bench
drwxr-xr-x    - atsetilam 26 lug 23:38  cmd
drwxr-xr-x    - atsetilam 26 lug 23:38  deploy
drwxr-xr-x    - atsetilam 22 ago 03:35  docs
drwxr-xr-x    - atsetilam 22 ago 03:35  examples
drwxr-xr-x    - atsetilam 22 ago 03:35  internal
drwxr-xr-x    - atsetilam 22 ago 03:35  migrations
drwxr-xr-x    - atsetilam 26 lug 23:38  pkg
drwxr-xr-x    - atsetilam 26 lug 23:38  test
drwxr-xr-x    - atsetilam 22 ago 03:35  tools
.rw-r--r--  207 atsetilam 22 ago 03:35  .dockerignore
.rw-r--r--  979 atsetilam 27 lug 06:05 󰊢 .gitignore
.rw-r--r-- 1,1k atsetilam 26 lug 23:38  .golangci.yml
.rw-r--r-- 4,3k atsetilam 27 lug 16:02  docker-compose.yml
.rw-r--r-- 1,3k atsetilam 26 lug 23:38  Dockerfile
.rw-r--r-- 5,0k atsetilam 22 ago 03:35  go.mod
.rw-r--r--  20k atsetilam 22 ago 03:35  go.sum
.rw-r--r--  14k atsetilam 26 lug 23:38  LICENSE
.rw-r--r-- 2,5k atsetilam 26 lug 23:38  Makefile
.rw-r--r--  102 atsetilam 22 ago 03:35  opencode.json
.rw-r--r-- 3,5k atsetilam 22 ago 03:35 󰂺 README.md
--- compose:
.rw-r--r-- 4,3k atsetilam 27 lug 16:02  /home/atsetilam/astrate/docker-compose.yml

$ gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex?ref=b6d46ad4 --jq '.content' 2>&1 | base64 -d | sed -n '1,200p'
#
# This file is part of Astarte.
#
# Copyright 2017-2023 SECO Mind Srl
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

defmodule Astarte.AppEngine.APIWeb.ErrorView do
  use Astarte.AppEngine.APIWeb, :view

  def render("400.json", _assigns) do
    %{errors: %{detail: "Bad request"}}
  end

  def render("422_unexpected_value_type.json", %{expected: expected} = _assigns) do
    %{errors: %{detail: "Unexpected value type", expected_type: expected}}
  end

  def render("422_value_size_exceeded.json", _assigns) do
    %{errors: %{detail: "Value size exceeds size limits"}}
  end

  def render("405_cannot_write_to_device_owned.json", _assigns) do
    %{errors: %{detail: "Cannot write to device owned resource"}}
  end

  def render("405_read_only_resource.json", _assigns) do
    %{errors: %{detail: "Cannot write to read-only resource"}}
  end

  def render("404.json", _assigns) do
    %{errors: %{detail: "Not found"}}
  end

  def render("404_device.json", _assigns) do
    %{errors: %{detail: "Device not found"}}
  end

  def render("404_endpoint_not_found.json", _assigns) do
    %{errors: %{detail: "Endpoint not found"}}
  end

  def render("404_interface_not_found.json", _assigns) do
    %{errors: %{detail: "Interface not found"}}
  end

  def render("404_interface_not_in_introspection.json", _assigns) do
    %{errors: %{detail: "Interface not found in device introspection"}}
  end

  def render("404_path.json", _assigns) do
    %{errors: %{detail: "Path not found"}}
  end

  def render("404_group.json", _assigns) do
    %{errors: %{detail: "Group not found"}}
  end

  def render("422_attribute_key_not_found.json", _assigns) do
    %{errors: %{detail: "Attribute key not found"}}
  end

  def render("409_group_already_exists.json", _assigns) do
    %{errors: %{detail: "Group already exists"}}
  end

  def render("409_device_already_in_group.json", _assigns) do
    %{errors: %{detail: "Device already in group"}}
  end

  def render("409_alias_already_in_use.json", _assigns) do
    %{errors: %{detail: "Alias already in use"}}
  end

  def render("422_unset_not_allowed.json", _assigns) do
    %{errors: %{detail: "Unset not allowed"}}
  end

  def render("422_alias_tag_not_found.json", _assigns) do
    %{errors: %{detail: "Alias tag not found"}}
  end

  def render("422_invalid_alias.json", _assigns) do
    %{errors: %{detail: "Invalid alias"}}
  end

  def render("422_invalid_attributes.json", _assigns) do
    %{errors: %{detail: "Invalid attributes"}}
  end

  def render("422_unexpected_object_key.json", %{keys: keys} = _assigns) do
    %{errors: %{detail: "Unexpected object key", unexpected_keys: keys}}
  end

  def render("422_missing_required_mapping.json", _assigns) do
    %{errors: %{detail: "Missing required mapping key"}}
  end

  def render("500.json", %{conn: %{assigns: %{reason: %Xandra.ConnectionError{}}}}) do
    %{errors: %{detail: "Database connection error"}}
  end

  def render("500.json", %{conn: %{assigns: %{reason: %Xandra.Error{message: message}}}}) do
    %{errors: %{detail: "Database error: #{message}"}}
  end

  def render("500.json", _assigns) do
    %{errors: %{detail: "Internal server error"}}
  end

  def render("401.json", _assigns) do
    %{errors: %{detail: "Unauthorized"}}
  end

  def render("403.json", _assigns) do
    %{errors: %{detail: "Forbidden"}}
  end

  def render("503_cannot_push_to_device.json", _assigns) do
    %{errors: %{detail: "Cannot push to device"}}
  end

  def render("503_service_unavailable.json", _assigns) do
    %{errors: %{detail: "Service unavailable"}}
  end

  def render("missing_token.json", _assigns) do
    %{errors: %{detail: "Missing authorization token"}}
  end

  def render("invalid_token.json", _assigns) do
    %{errors: %{detail: "Invalid JWT token"}}
  end

  def render("invalid_auth_path.json", _assigns) do
    %{
      errors: %{
        detail: "Authorization failed due to an invalid path"
      }
    }
  end

  def render("authorization_path_not_matched.json", %{method: method, path: path}) do
    %{
      errors: %{
        detail: "Unauthorized access to #{method} #{path}. Please verify your permissions"
      }
    }
  end

  # In case no render clause matches or no
  # template is found, let's render it as 500
  def template_not_found(_template, assigns) do
    render("500.json", assigns)
  end
end

$ ssh -o ConnectTimeout=8 legion -T 'ls ~/.local/bin 2>/dev/null | head; command -v dockerd dockerd-rootless.sh elixir mix 2>&1' 2>&1 | sed 's/\x1b\[[0-9;]*m//g' | head -20
lrwxrwxrwx    - atsetilam  3 ott 15:41 󰡯 claude -> /home/atsetilam/.local/share/claude/versions/2.1.288
lrwxrwxrwx    - atsetilam  6 ago 19:19 󰡯 opencode-telegram -> ../lib/node_modules/@grinev/opencode-telegram-bot/dist/cli.js
lrwxrwxrwx    - atsetilam 27 ago 02:17 󰡯 orange-cat -> /home/atsetilam/orange-cat/orange-cat
.rwxr-xr-x 3,6M atsetilam 27 ago 03:39 󰡯 orange-cat-tui
lrwxrwxrwx    - atsetilam  3 ott 16:16 󰡯 toshy-config-restart -> /home/atsetilam/.config/toshy/scripts/bin/toshy-config-restart.sh
lrwxrwxrwx    - atsetilam  3 ott 16:16 󰡯 toshy-config-start -> /home/atsetilam/.config/toshy/scripts/bin/toshy-config-start.sh
lrwxrwxrwx    - atsetilam  3 ott 16:16 󰡯 toshy-config-stop -> /home/atsetilam/.config/toshy/scripts/bin/toshy-config-stop.sh
lrwxrwxrwx    - atsetilam  3 ott 16:16 󰡯 toshy-cosmic-dbus-service -> /home/atsetilam/.config/toshy/scripts/bin/toshy-cosmic-dbus-service.sh
lrwxrwxrwx    - atsetilam  3 ott 16:16 󰡯 toshy-debug -> /home/atsetilam/.config/toshy/scripts/bin/toshy-config-start-verbose.sh
lrwxrwxrwx    - atsetilam  3 ott 16:16 󰡯 toshy-detector-check -> /home/atsetilam/.config/toshy/scripts/bin/toshy-detector-check.sh
/usr/bin/dockerd

$ gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/fallback_controller.ex?ref=b6d46ad4 --jq '.content' | base64 -d | sed -n '17,200p'
defmodule Astarte.AppEngine.APIWeb.FallbackController do
  require Logger

  @moduledoc """
  Translates controller action results into valid `Plug.Conn` responses.

  See `Phoenix.Controller.action_fallback/1` for more details.
  """
  use Astarte.AppEngine.APIWeb, :controller

  def call(conn, {:error, %Ecto.Changeset{} = changeset}) do
    conn
    |> put_status(:unprocessable_entity)
    |> put_view(Astarte.AppEngine.APIWeb.ChangesetView)
    |> render("error.json", changeset: changeset)
  end

  def call(conn, {:error, :cannot_write_to_device_owned}) do
    conn
    |> put_status(:method_not_allowed)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"405_cannot_write_to_device_owned")
  end

  def call(conn, {:error, :device_not_found}) do
    conn
    |> put_status(:not_found)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"404_device")
  end

  def call(conn, {:error, :cannot_push_to_device}) do
    conn
    |> put_status(:service_unavailable)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"503_cannot_push_to_device")
  end

  def call(conn, {:error, :endpoint_not_found}) do
    conn
    |> put_status(:bad_request)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"404_endpoint_not_found")
  end

  def call(conn, {:error, :extended_id_not_allowed}) do
    conn
    |> put_status(:bad_request)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"400")
  end

  def call(conn, {:error, :interface_not_found}) do
    conn
    |> put_status(:not_found)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"404_interface_not_found")
  end

  def call(conn, {:error, :interface_not_in_introspection}) do
    conn
    |> put_status(:not_found)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"404_interface_not_in_introspection")
  end

  def call(conn, {:error, :invalid_device_id}) do
    conn
    |> put_status(:bad_request)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"400")
  end

  def call(conn, {:error, :not_found}) do
    conn
    |> put_status(:not_found)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"404")
  end

  def call(conn, {:error, :path_not_found}) do
    conn
    |> put_status(:not_found)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"404_path")
  end

  def call(conn, {:error, :group_not_found}) do
    conn
    |> put_status(:not_found)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"404_group")
  end

  def call(conn, {:error, :group_already_exists}) do
    conn
    |> put_status(:conflict)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"409_group_already_exists")
  end

  def call(conn, {:error, :device_already_in_group}) do
    conn
    |> put_status(:conflict)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"409_device_already_in_group")
  end

  def call(conn, {:error, :read_only_resource}) do
    conn
    |> put_status(:method_not_allowed)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"405_read_only_resource")
  end

  def call(conn, {:error, :unauthorized}) do
    _ = Logger.info("Refusing unauthorized request.", tag: "unauthorized")

    conn
    |> put_status(:unauthorized)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"401")
  end

  def call(conn, {:error, :unexpected_value_type, expected: expected}) do
    conn
    |> put_status(:unprocessable_entity)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(
      :"422_unexpected_value_type",
      expected: expected
    )
  end

  def call(conn, {:error, :value_size_exceeded}) do
    conn
    |> put_status(:unprocessable_entity)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"422_value_size_exceeded")
  end

  def call(conn, {:error, :alias_already_in_use}) do
    conn
    |> put_status(:conflict)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"409_alias_already_in_use")
  end

  def call(conn, {:error, :attribute_key_not_found}) do
    conn
    |> put_status(:unprocessable_entity)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"422_attribute_key_not_found")
  end

  def call(conn, {:error, :mapping_not_found}) do
    conn
    |> put_status(:not_found)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"404_mapping_not_found")
  end

  def call(conn, {:error, :invalid_alias}) do
    conn
    |> put_status(:bad_request)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"422_invalid_alias")
  end

  def call(conn, {:error, :alias_tag_not_found}) do
    conn
    |> put_status(:bad_request)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"422_alias_tag_not_found")
  end

  def call(conn, {:error, :invalid_attributes}) do
    conn
    |> put_status(:bad_request)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"422_invalid_attributes")
  end

$ gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/fallback_controller.ex?ref=b6d46ad4 --jq '.content' | base64 -d | sed -n '200,300p'
def call(conn, {:error, :unexpected_object_key, keys: keys}) do
    conn
    |> put_status(:bad_request)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"422_unexpected_object_key", keys: keys)
  end

  def call(conn, {:error, :unset_not_allowed}) do
    conn
    |> put_status(:unprocessable_entity)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"422_unset_not_allowed")
  end

  def call(conn, {:error, :missing_required_mapping}) do
    conn
    |> put_status(:unprocessable_entity)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"422_missing_required_mapping")
  end

  def call(conn, {:error, :realm_not_found}) do
    conn
    |> put_status(:unauthorized)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:invalid_token)
  end

  # Invalid authorized path
  def call(conn, {:error, :invalid_auth_path}) do
    conn
    |> put_status(:unauthorized)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:invalid_auth_path)
  end

  # This is called when no JWT token is present
  def auth_error(conn, {:unauthenticated, :unauthenticated}, _opts) do
    conn
    |> put_status(:unauthorized)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:missing_token)
  end

  # Invalid JWT token
  def auth_error(conn, {:invalid_token, :invalid_token}, _opts) do
    conn
    |> put_status(:unauthorized)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:invalid_token)
  end

  # Path not authorized
  def auth_error(conn, {:unauthorized, :authorization_path_not_matched}, _opts) do
    conn
    |> put_status(:forbidden)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:authorization_path_not_matched, %{method: conn.method, path: conn.request_path})
  end

  def auth_error(conn, {:unauthenticated, _reason}, _opts) do
    conn
    |> put_status(:unauthorized)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"401")
  end

  # In all other cases, we reply with 403
  def auth_error(conn, _reason, _opts) do
    conn
    |> put_status(:forbidden)
    |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
    |> render(:"403")
  end
end

$ gh search code --owner astarte-platform 'Missing required mapping key' --limit 20 2>&1 | head -20
astarte-platform/astarte:apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex: %{errors: %{detail: "Missing required mapping key"}}
astarte-platform/astarte:apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex: "Missing required mapping key #{inspect(key)} in object interface value.",
astarte-platform/astarte:apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/data_handler.ex: message: "Missing required mapping key in object sent to #{interface}#{path}.",

$ rg -n "writeError\(" internal/appengine/*.go | head -20; echo "=== engine reason flow:"; rg -n "ReasonOf|RejectError" internal/engine/*.go internal/appengine/*.go | head -20
internal/appengine/writeerror_test.go:53:			a.writeError(rec, tc.err)
internal/appengine/writeerror_test.go:87:	a.writeError(rec, &payload.RejectError{
internal/appengine/http.go:92:		a.writeError(w, err)
internal/appengine/http.go:134:		a.writeError(w, err)
internal/appengine/http.go:146:		a.writeError(w, err)
internal/appengine/http.go:155:		a.writeError(w, err)
internal/appengine/http.go:173:		a.writeError(w, err)
internal/appengine/http.go:182:		a.writeError(w, err)
internal/appengine/http.go:192:		a.writeError(w, err)
internal/appengine/http.go:221:		a.writeError(w, err)
internal/appengine/http.go:242:		a.writeError(w, err)
internal/appengine/http.go:257:		a.writeError(w, err)
internal/appengine/http.go:262:		a.writeError(w, err)
internal/appengine/http.go:306:		a.writeError(w, err)
internal/appengine/http.go:316:		a.writeError(w, err)
internal/appengine/http.go:335:		a.writeError(w, err)
internal/appengine/http.go:354:		a.writeError(w, err)
internal/appengine/http.go:364:		a.writeError(w, err)
internal/appengine/http.go:374:		a.writeError(w, err)
internal/appengine/http.go:385:		a.writeError(w, err)
=== engine reason flow:
internal/engine/serverdata.go:18:// statuses; payload validation failures surface as *payload.RejectError.
internal/engine/data.go:289:		reason := payload.ReasonOf(err)
internal/appengine/writeerror_test.go:35:			&payload.RejectError{Reason: payload.ReasonBadObject,
internal/appengine/writeerror_test.go:40:			fmt.Errorf("engine: %w", &payload.RejectError{Reason: payload.ReasonBadObject,
internal/appengine/writeerror_test.go:45:			&payload.RejectError{Reason: payload.ReasonBadObject, Detail: "object-aggregation document is empty"},
internal/appengine/writeerror_test.go:87:	a.writeError(rec, &payload.RejectError{
internal/engine/serverdata_test.go:196:		if payload.ReasonOf(err) == payload.ReasonNone {
internal/engine/serverdata_test.go:338:	if payload.ReasonOf(err) != payload.ReasonBadObject {
internal/engine/serverdata_test.go:339:		t.Fatalf("reason = %v (err %v); want bad_object", payload.ReasonOf(err), err)
internal/engine/serverdata_test.go:341:	var re *payload.RejectError
internal/engine/serverdata_test.go:343:		t.Fatalf("err = %v; want a *payload.RejectError", err)
internal/appengine/http.go:684:	case payload.ReasonOf(err) == payload.ReasonBadObject:
internal/appengine/http.go:709:	var re *payload.RejectError
internal/engine/devicedata_test.go:164:		if got := payload.ReasonOf(err); got != payload.ReasonTypeMismatch {

$ gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex?ref=b6d46ad4 --jq '.content' | base64 -d | grep -n -B18 'missing_required_mapping' | head -80
482-      validate_required_mappings(mappings_by_key, value)
483-    end
484-  end
485-
486-  defp validate_value(_mappings_by_key, _value) do
487-    {:error, :unexpected_value_type, expected: :object}
488-  end
489-
490-  defp validate_required_mappings(mappings_by_key, value) when is_map(value) do
491-    case Enum.find(mappings_by_key, fn {key, mapping} ->
492-           mapping.required and not Map.has_key?(value, key)
493-         end) do
494-      nil ->
495-        :ok
496-
497-      {key, _mapping} ->
498-        Logger.warning(
499-          "Missing required mapping key #{inspect(key)} in object interface value.",
500:          tag: "missing_required_mapping"
501-        )
502-
503:        {:error, :missing_required_mapping}
--
567-
568-      {:error, :unexpected_object_key, keys: keys} ->
569-        Logger.warning("Unexpected object keys #{inspect(keys)} in object interface update.",
570-          tag: "unexpected_object_key"
571-        )
572-
573-        {:error, :unexpected_object_key, keys: keys}
574-
575-      {:error, :invalid_object_aggregation_path} ->
576-        Logger.warning("Error while trying to publish on path for object aggregated interface.",
577-          tag: "invalid_object_aggregation_path"
578-        )
579-
580-        {:error, :invalid_object_aggregation_path}
581-
582-      {:error, :mapping_not_found} ->
583-        {:error, :mapping_not_found}
584-
585:      {:error, :missing_required_mapping} ->
586-        Logger.warning("Missing required mapping in object interface update.",
587:          tag: "missing_required_mapping"
588-        )
589-
590:        {:error, :missing_required_mapping}

$ rg -n "docker|compose" tools/mule.sh | head -20; echo "=== legion helpers:"; rg -n "^\s*(legion_|[a-z-]+\(\))" tools/mule.sh | head -40
121:      if legion_sh 'docker info >/dev/null 2>&1'; then
122:        ok "docker is up ($(legion_sh 'docker ps -q | wc -l' 2>/dev/null | tr -d ' ') containers)"
123:      else bad "docker is not up — heavy tasks cannot run"; return 1; fi
=== legion helpers:
34:note() { printf '\033[36m>>\033[0m %s\n' "$*" >&2; }
35:bad()  { printf '\033[31m!!\033[0m %s\n' "$*" >&2; }
36:ok()   { printf '\033[32mok\033[0m %s\n' "$*" >&2; }
37:die()  { bad "$*"; exit 1; }
100:legion_sh() {
110:legion_up() {
112:  legion_sh "${MULE_LEGION_PROBE:-true}" >/dev/null 2>&1
118:      legion_sh 'echo up' >/dev/null 2>&1 \
124:      legion_up && ok "[legion] tasks are runnable" || { bad "probe failed"; return 1; }
131:      legion_up || die "Legion Go is not available"
136:      legion_sh "mkdir -p '$MULE_LEGION_REPO/bench'" || die "mkdir failed on the Legion Go"
140:      legion_sh "chmod +x '$MULE_LEGION_REPO/bench/bench'"
231:gates() {
737:beat() { date +%s > "$MULE/.heartbeat"; rm -f "$MULE/.alarmed"; }

$ rg -n "missing_required|unexpected-object-key|Missing required mapping" .mule/ docs/ 2>/dev/null | head -20
.mule/todo.md:1:- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
.mule/todo.md:265:- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
.mule/todo.md:266:- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress. — BLOCKED: TIMEOUT after 1200s — task too big, split it
.mule/todo.md:267:- [ ] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
.mule/reviews/engine-triggers-2026-09-20.md:28:### 1. `missing_required` is the only `astrateToUpstream` key with no test row
.mule/reviews/engine-triggers-2026-09-20.md:31:`missing_required` → `unexpected_object_key`. `TestUpstreamErrorNameMapping`
.mule/log.md:167:| 2026-10-03 | appengine-unexpected-object-key | done | 956s | f7d60af |
.mule/reports/log/20260926T113138Z-recipe-milestones-proposed.md:265:Both code phases have landed (4a schema layer \`9c4d411\`, 4b runtime enforcement of \`required\` \`7a244b2\`). What is left is the docs closeout — the \`docs/COMPATIBILITY.md\` row explaining why astrate’s \`missing_required\` reason is translated to \`unexpected_object_key\`, and the \`docs/UPSTREAM-EXPERIMENTAL.md\` Adopted dates — which is prose the mule does not write, and it lands together with #68's closeout in fase 4c. Keeping this issue out of the mule queue on purpose.
.mule/reports/log/20260926T113138Z-recipe-milestones-proposed.md:274:Decided and shipped on 2026-08-23, before this issue was revisited: `required` and `encrypted` are parsed, stored, carried through interface compatibility checks, and `required` is enforced at runtime — an object-aggregated document missing a required key is rejected (`missing_required`, translated to the closest 1.2.2 error_name for the dashboard's closed set). `encrypted` is parsed and stored only, pending the keyAgreement work.
.mule/reports/log/20260928T020420Z-recipe-astarte-upstream-proposed.md:38:- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
.mule/reports/log/20260924T201012Z-recipe-code-review-timeout.md:264:739df67 mule: errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
.mule/reports/log/20260924T201012Z-recipe-code-review-timeout.md:266:fd562f7 payload: pre-wire missing_required reject reason label fixture and error-name translation (fase 4b prep)
.mule/reports/log/20260926T164829Z-recipe-github-issues-proposed.md:135:- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
.mule/reports/log/20261001T203619Z-recipe-astarte-upstream-proposed.md:94:1:- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
.mule/reports/log/20260920T184908Z-triggers-custom-action-policy-nodecide-done.md:194: - [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
.mule/reports/log/20261003T165209Z-docs-sync-native-error-example-capitalisation-done.md:224:- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
.mule/reports/log/20261003T190400Z-appengine-unexpected-object-key-done.md:1:slug: appengine-unexpected-object-key
.mule/reports/log/20261003T190400Z-appengine-unexpected-object-key-done.md:1024:zz_probe_test.go:18: missing_required   -> 500 {"errors":{"detail":"Internal server error"}}
.mule/reports/log/20261003T190400Z-appengine-unexpected-object-key-done.md:1056:- [ ] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
.mule/reports/log/20261003T190400Z-appengine-unexpected-object-key-done.md:1064:- [ ] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress.

$ rg -n "legion|docker|probe" .mule/reports/log/20261003T190400Z-appengine-unexpected-object-key-done.md | head -40
1001:$ cat > internal/appengine/zz_probe_test.go <<'EOF'
1018:		(&API{}).writeError(rec, &payload.RejectError{Reason: r, Detail: "probe"})
1023:go test ./internal/appengine/ -run TestZZProbe -count=1 -v 2>&1 | grep -E "zz_probe|PASS|FAIL" | head -20; rm internal/appengine/zz_probe_test.go
1024:zz_probe_test.go:18: missing_required   -> 500 {"errors":{"detail":"Internal server error"}}
1025:    zz_probe_test.go:18: too_large          -> 500 {"errors":{"detail":"Internal server error"}}
1026:    zz_probe_test.go:18: malformed          -> 500 {"errors":{"detail":"Internal server error"}}
1027:    zz_probe_test.go:18: no_value           -> 500 {"errors":{"detail":"Internal server error"}}
1028:    zz_probe_test.go:18: type_mismatch      -> 500 {"errors":{"detail":"Internal server error"}}
1029:    zz_probe_test.go:18: value_too_large    -> 500 {"errors":{"detail":"Internal server error"}}
1030:    zz_probe_test.go:18: unset_not_allowed  -> 500 {"errors":{"detail":"Internal server error"}}
1031:    zz_probe_test.go:18: bad_timestamp      -> 500 {"errors":{"detail":"Internal server error"}}
1032:    zz_probe_test.go:18: unknown_format     -> 500 {"errors":{"detail":"Internal server error"}}
1050:- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.
1064:- [ ] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress.
1065:- [ ] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
1270:Verified: `go vet ./...` -> pass (exit 0); `go test ./...` -> pass, 0 FAIL, 0 failures anywhere (run **without** `-race`: this host has the same 39-bit VMA kernel as the Pi, `go test -race` dies with `ThreadSanitizer: unsupported VMA range`); `gofmt -l .` -> clean; `go vet -tags integration ./internal/appengine/ ./internal/engine/` -> clean (builds, not run); mutation check — disabling only the new `writeError` case makes `TestWriteErrorTaxonomy` and `TestWriteErrorBadObjectEnvelope` fail with `status = 500, want 400`, confirming the old behaviour and that the tests are not vacuous; the probe measured the report claim below.
1272:Unsure: two things I did not measure on upstream (no container-free way to reach it, and the Legion Go is off per AGENTS.md, so I did not probe it):
1277:Report half (measured, not changed): the hole is not specific to `bad_object`. A throwaway probe over `writeError` returned **500 "Internal server error"** for every one of the nine remaining reasons — `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp`, `unknown_format`. So `ReasonMissingRequired` has the same hole, and the upstream 422 `missing_required_mapping` status/detail must be probed on the Legion Go before mapping — left as a follow-up line, not guessed.
1279:Follow-ups: appended to `.mule/todo.md` — `appengine-payload-reason-status-map` (map the remaining nine reasons, each probed on the Legion Go; notes that `too_large` may be unreachable from REST because `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner`) and `appengine-missing-required-422` (probe the status/detail, then map `ReasonMissingRequired`).

$ gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/device_controller_test.exs?ref=b6d46ad4 --jq '.content' 2>&1 | base64 -d | grep -n -B12 "missing_required" | head -60
base64: invalid input

$ gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/fallback_controller.ex?ref=b6d46ad4 --jq '.content' | base64 -d | grep -n "def call(conn"
28:  def call(conn, {:error, %Ecto.Changeset{} = changeset}) do
35:  def call(conn, {:error, :cannot_write_to_device_owned}) do
42:  def call(conn, {:error, :device_not_found}) do
49:  def call(conn, {:error, :cannot_push_to_device}) do
56:  def call(conn, {:error, :endpoint_not_found}) do
63:  def call(conn, {:error, :extended_id_not_allowed}) do
70:  def call(conn, {:error, :interface_not_found}) do
77:  def call(conn, {:error, :interface_not_in_introspection}) do
84:  def call(conn, {:error, :invalid_device_id}) do
91:  def call(conn, {:error, :not_found}) do
98:  def call(conn, {:error, :path_not_found}) do
105:  def call(conn, {:error, :group_not_found}) do
112:  def call(conn, {:error, :group_already_exists}) do
119:  def call(conn, {:error, :device_already_in_group}) do
126:  def call(conn, {:error, :read_only_resource}) do
133:  def call(conn, {:error, :unauthorized}) do
142:  def call(conn, {:error, :unexpected_value_type, expected: expected}) do
152:  def call(conn, {:error, :value_size_exceeded}) do
159:  def call(conn, {:error, :alias_already_in_use}) do
166:  def call(conn, {:error, :attribute_key_not_found}) do
173:  def call(conn, {:error, :mapping_not_found}) do
180:  def call(conn, {:error, :invalid_alias}) do
187:  def call(conn, {:error, :alias_tag_not_found}) do
194:  def call(conn, {:error, :invalid_attributes}) do
201:  def call(conn, {:error, :unexpected_object_key, keys: keys}) do
208:  def call(conn, {:error, :unset_not_allowed}) do
215:  def call(conn, {:error, :missing_required_mapping}) do
222:  def call(conn, {:error, :realm_not_found}) do
230:  def call(conn, {:error, :invalid_auth_path}) do

$ gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers?ref=b6d46ad4 --jq '.[].name' 2>&1 | head -20
device_status_by_alias_controller_test.exs
device_status_by_group_controller_test.exs
device_status_controller_test.exs
groups_controller_test.exs
interface_controller_test.exs
interface_values_by_device_alias_controller_test.exs
interface_values_by_group_controller_test.exs
interface_values_controller_test.exs
stats_controller_test.exs

$ gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex?ref=b6d46ad4 --jq '.content' | base64 -d | grep -n "validate_value\|def send\|action_fallback" | head -20; echo "=== controller:"; gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/device_controller.ex?ref=b6d46ad4 --jq '.content' | base64 -d | grep -n "action_fallback\|def send\|Device.send\|put_status" | head -20
209:         :ok <- validate_value_type(mapping.value_type, value),
480:  defp validate_value(mappings_by_key, value) when is_map(value) do
481:    with :ok <- validate_value_type(mappings_by_key, value) do
486:  defp validate_value(_mappings_by_key, _value) do
528:         :ok <- validate_value(mappings_by_key, value),
871:  defp validate_value_type(mappings_by_key, object)
875:           :ok <- validate_value_type(expected_type, value) do
891:  defp validate_value_type(value_type, value) do
892:    case ValueType.validate_value(value_type, value) do
916:      # We can be sure this exists since we validated it in validate_value_type
=== controller:
gh: Not Found (HTTP 404)
base64: invalid input

$ gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_values_controller_test.exs?ref=b6d46ad4 --jq '.content' | base64 -d > /tmp/opencode/ivc_test.exs; grep -n -B20 "missing_required" /tmp/opencode/ivc_test.exs | head -70
(no output)

$ grep -rn "Missing required mapping\|missing_required" /tmp/opencode/ivc_test.exs | head; echo "=== file size:"; wc -l /tmp/opencode/ivc_test.exs; echo "=== action_fallback in interface_values_controller:"; gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/interface_values_controller.ex?ref=b6d46ad4 --jq '.content' | base64 -d | grep -n "action_fallback\|def send\|with\b" | head -20
=== file size:
185 /tmp/opencode/ivc_test.exs
=== action_fallback in interface_values_controller:
7:# you may not use this file except in compliance with the License.
27:  action_fallback Astarte.AppEngine.APIWeb.FallbackController
167:      unprocessable_entity: "Tried to unset a property with `allow_unset` false"
196:    with {:ok, interfaces} <- Device.list_interfaces(realm_name, device_id) do
218:    with {:ok, %InterfaceValues{} = interface_values} <-
229:    with {:ok, %InterfaceValues{} = interface_values} <-
245:    with {:ok, %InterfaceValues{} = interface_values} <-
264:    with :ok <- Device.delete_interface_values(realm_name, device_id, interface, path) do

$ gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/interface_values_controller.ex?ref=b6d46ad4 --jq '.content' | base64 -d | sed -n '140,185p'
]

  operation :update,
    summary: "Update and push a value on a path",
    description:
      "Update and push a property value to the device on a certain endpoint path. interface should be an individual server owned property interface. It mustn't be used to stream data to a device or to update single properties that are members of an object aggregated interface.",
    operation_id: "updatePathValue",
    parameters: [
      realm_name: [
        in: :path,
        description: "Name of the realm which the device belongs to.",
        required: true,
        type: :string
      ],
      device_id: [in: :path, description: "Device ID", required: true, type: :string],
      interface: [in: :path, description: "Interface name", required: true, type: :string],
      path: [in: :path, description: "Endpoint Path", required: true, type: :string]
    ],
    responses: [
      ok: "Success",
      bad_request: "Bad request",
      unauthorized: %Reference{"$ref": "#/components/responses/Unauthorized"},
      forbidden: %Reference{"$ref": "#/components/responses/AuthorizationPathNotMatched"},
      not_found:
        {"Endpoint not found or interface not found in introspection or device not found.",
         "application/json", Errors.NotFoundError},
      method_not_allowed: "Invalid object",
      unprocessable_entity: "Tried to unset a property with `allow_unset` false"
    ]

  operation :delete,
    summary: "Delete path and push an unset value message",
    description:
      "Unset a value on a certain path, path is also deleted. Endpoint must support unset.",
    operation_id: "deletePathValue",
    parameters: [
      realm_name: [
        in: :path,
        description: "Name of the realm which the device belongs to.",
        required: true,
        type: :string
      ],
      device_id: [in: :path, description: "Device ID", required: true, type: :string],
      interface: [in: :path, description: "Interface name", required: true, type: :string],
      path: [in: :path, description: "Endpoint Path", required: true, type: :string]
    ],

$ gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex?ref=b6d46ad4 --jq '.content' | base64 -d | sed -n '515,535p'
now = DateTime.utc_now()

    with {:ok, mappings} <-
           Mappings.fetch_interface_mappings(
             realm_name,
             interface_descriptor.interface_id
           ),
         {:ok, {endpoint_id, mapping, current_interface_mappings}} <-
           resolve_object_aggregation_path(path, interface_descriptor, mappings),
         mappings_by_key = extract_mappings(mappings),
         expected_types =
           Map.new(mappings_by_key, fn {k, %Mapping{value_type: t}} -> {k, t} end),
         {:ok, value} <- InterfaceValue.cast_value(expected_types, raw_value),
         :ok <- validate_value(mappings_by_key, value),
         reliability = extract_aggregate_reliability(mappings),
         interface_type = interface_descriptor.type,
         publish_opts = build_publish_opts(interface_type, reliability),
         interface_name = interface_descriptor.name,
         publish_value =
           maybe_apply_transport_encryption(
             realm_name,

$ rg -n "func ReasonOf" -A18 pkg/payload/value.go
229:func ReasonOf(err error) RejectReason {
230-	var re *RejectError
231-	if errors.As(err, &re) {
232-		return re.Reason
233-	}
234-	return ReasonNone
235-}
236-
237-// ---------------------------------------------------------------------------
238-// Coercion table (docs/DESIGN.md §2.6 step 5). Each helper returns a
239-// *RejectError on failure so decoders can pass errors through unchanged.
240-// ---------------------------------------------------------------------------
241-
242-// maxInt64Float is 2^63 as a float64 (exactly representable); float values
243-// in [-2^63, 2^63) convert to int64 without implementation-defined results.
244-const maxInt64Float = 9223372036854775808.0
245-
246-// checkDouble rejects NaN and ±Inf, which neither Astarte's BSON profile
247-// nor JSON can represent.

$ rg -n "func WriteError|func WriteBadRequest|func WriteRawErrors|func WriteInternalServerError" -A12 pkg/astarteapi/*.go | head -60
pkg/astarteapi/envelope.go:163:func WriteError(w http.ResponseWriter, status int, detail string) error {
pkg/astarteapi/envelope.go-164-	return write(w, status, detailEnvelope{Errors: detailBody{Detail: detail}})
pkg/astarteapi/envelope.go-165-}
pkg/astarteapi/envelope.go-166-
pkg/astarteapi/envelope.go-167-// WriteFieldErrors writes the Phoenix-changeset-shaped error envelope
pkg/astarteapi/envelope.go-168-// {"errors": {"<field>": ["<message>", ...]}} used by upstream 422 validation
pkg/astarteapi/envelope.go-169-// failures (for example {"errors": {"hw_id": ["is invalid"]}}). Keys are
pkg/astarteapi/envelope.go-170-// emitted in sorted order (Go map marshalling), which matches the
pkg/astarteapi/envelope.go-171-// deterministic bodies upstream produces for single-field failures.
pkg/astarteapi/envelope.go-172-func WriteFieldErrors(w http.ResponseWriter, status int, fields map[string][]string) error {
pkg/astarteapi/envelope.go-173-	return write(w, status, fieldsEnvelope{Errors: fields})
pkg/astarteapi/envelope.go-174-}
pkg/astarteapi/envelope.go-175-
--
pkg/astarteapi/envelope.go:186:func WriteRawErrors(w http.ResponseWriter, status int, errs any) error {
pkg/astarteapi/envelope.go-187-	return write(w, status, rawEnvelope{Errors: errs})
pkg/astarteapi/envelope.go-188-}
pkg/astarteapi/envelope.go-189-
pkg/astarteapi/envelope.go-190-// WriteBadRequest writes the canonical 400 envelope.
pkg/astarteapi/envelope.go:191:func WriteBadRequest(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go-192-	return WriteError(w, http.StatusBadRequest, DetailBadRequest)
pkg/astarteapi/envelope.go-193-}
pkg/astarteapi/envelope.go-194-
pkg/astarteapi/envelope.go-195-// WriteUnauthorized writes the canonical 401 envelope.
pkg/astarteapi/envelope.go-196-func WriteUnauthorized(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go-197-	return WriteError(w, http.StatusUnauthorized, DetailUnauthorized)
pkg/astarteapi/envelope.go-198-}
pkg/astarteapi/envelope.go-199-
pkg/astarteapi/envelope.go-200-// WriteForbidden writes the canonical 403 envelope.
pkg/astarteapi/envelope.go-201-func WriteForbidden(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go-202-	return WriteError(w, http.StatusForbidden, DetailForbidden)
pkg/astarteapi/envelope.go-203-}
--
pkg/astarteapi/envelope.go:216:func WriteInternalServerError(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go-217-	return WriteError(w, http.StatusInternalServerError, DetailInternalServerError)
pkg/astarteapi/envelope.go-218-}
pkg/astarteapi/envelope.go-219-
pkg/astarteapi/envelope.go-220-// DecodeData reads at most maxBytes bytes from r, unwraps the mandatory
pkg/astarteapi/envelope.go-221-// {"data": ...} request envelope, and unmarshals the "data" value into dst.
pkg/astarteapi/envelope.go-222-//
pkg/astarteapi/envelope.go-223-// It fails with an error wrapping ErrBodyTooLarge when the body exceeds
pkg/astarteapi/envelope.go-224-// maxBytes, with one wrapping ErrMissingData when the "data" key is absent or
pkg/astarteapi/envelope.go-225-// null, and with the underlying JSON error for malformed bodies (including
pkg/astarteapi/envelope.go-226-// trailing garbage after the top-level value). Sibling keys next to "data"
pkg/astarteapi/envelope.go-227-// are ignored, matching upstream parameter handling.
pkg/astarteapi/envelope.go-228-func DecodeData(r io.Reader, maxBytes int64, dst any) error {

Upstream chain confirmed. Now the change.
← Edit internal/appengine/http.go
Index: /root/astrate-mule/internal/appengine/http.go
===================================================================
--- /root/astrate-mule/internal/appengine/http.go
+++ /root/astrate-mule/internal/appengine/http.go
@@ -688,8 +688,19 @@
 		// after the sentinels because a payload rejection carries none of
 		// them, and last because every other reason stays unmapped (report
 		// only — see .mule/todo.md).
 		a.writeBadObjectError(w, err)
+	case payload.ReasonOf(err) == payload.ReasonMissingRequired:
+		// A distinct upstream failure from the object shape above, and the
+		// only other mapped payload reason: 422, because the document parsed
+		// and every key is declared — one of them the mapping declares
+		// required and the document omits (master b6d46ad4,
+		// Device.validate_required_mappings/2 answers
+		// :missing_required_mapping, the AppEngine fallback controller puts
+		// :unprocessable_entity for it, ErrorView renders the detail alone).
+		// Placed after bad_object because that one names the offending keys
+		// and answers 400; both precede default: so neither falls to the 500.
+		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, detailMissingRequiredMapping)
 	default:
 		_ = astarteapi.WriteInternalServerError(w)
 	}
 }

← Edit internal/appengine/http.go
Index: /root/astrate-mule/internal/appengine/http.go
===================================================================
--- /root/astrate-mule/internal/appengine/http.go
+++ /root/astrate-mule/internal/appengine/http.go
@@ -709,8 +709,17 @@
 // write carrying undeclared keys (master b6d46ad4, #2237). Frozen string:
 // astartectl and the SDKs match on it.
 const detailUnexpectedObjectKey = "Unexpected object key"
 
+// detailMissingRequiredMapping is upstream's 422 detail for an
+// object-aggregated write whose document omits a key the mapping declares
+// required (master b6d46ad4, ErrorView "422_missing_required_mapping.json"
+// renders `%{errors: %{detail: "Missing required mapping key"}}`). Frozen
+// string; unlike the undeclared-key case the body carries no key list —
+// upstream logs the offending key (device.ex:498-502) and renders the detail
+// alone.
+const detailMissingRequiredMapping = "Missing required mapping key"
+
 // writeBadObjectError renders the object-shape rejection in the measured
 // upstream envelope: {"errors": {"detail": "Unexpected object key",
 // "unexpected_keys": [...]}}. The other object-shape failures — not a
 // document, an empty document, a duplicate key — carry no key list and answer

$ rg -n "fmt.Errorf" internal/engine/serverdata.go | head -10
57:		return fmt.Errorf("engine: encoding server value: %w", err)
61:		return fmt.Errorf("engine: publishing %s: %w", topic, err)
107:		return nil, fmt.Errorf("%w: %s", ErrRealmUnknown, realm)
115:		return nil, fmt.Errorf("%w: %s", ErrInterfaceNotFound, ifaceName)
118:		return nil, fmt.Errorf("%w: %s", notOwned, ifaceName)
137:			return nil, fmt.Errorf("%w: %q is not an aggregation prefix of %s", ErrPathNotFound, path, ifaceName)
145:			return nil, fmt.Errorf("%w: %q matches no endpoint of %s", ErrPathNotFound, path, ifaceName)
153:		return nil, fmt.Errorf("%w: null value (use the property DELETE path for unset)", ErrUnsetNotAllowed)
225:		return fmt.Errorf("%w: %s", ErrRealmUnknown, realm)
233:		return fmt.Errorf("%w: %s", ErrInterfaceNotFound, ifaceName)

→ Read internal/engine/serverdata.go [limit=30, offset=125]

← Edit internal/appengine/writeerror_test.go
Index: /root/astrate-mule/internal/appengine/writeerror_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/writeerror_test.go
+++ /root/astrate-mule/internal/appengine/writeerror_test.go
@@ -43,8 +43,12 @@
 		http.StatusBadRequest, "Unexpected object key", []string{"a", "b"}},
 	{"object write with no key list stays the canonical 400",
 		&payload.RejectError{Reason: payload.ReasonBadObject, Detail: "object-aggregation document is empty"},
 		http.StatusBadRequest, "Bad request", nil},
+	{"object write omitting a required key answers 422",
+		&payload.RejectError{Reason: payload.ReasonMissingRequired,
+			Detail: "object-aggregation document is missing required key(s) lat"},
+		http.StatusUnprocessableEntity, "Missing required mapping key", nil},
 	{"unknown cause stays 500", errors.New("boom"),
 		http.StatusInternalServerError, "Internal server error", nil},
 }
 for _, tc := range cases {

$ gofmt -l . && go vet ./internal/appengine/ && go test ./internal/appengine/ -run 'TestWriteError' -count=1 2>&1 | tail -5
ok  	github.com/astrate-platform/astrate/internal/appengine	0.019s

← Edit internal/appengine/http.go
Index: /root/astrate-mule/internal/appengine/http.go
===================================================================
--- /root/astrate-mule/internal/appengine/http.go
+++ /root/astrate-mule/internal/appengine/http.go
@@ -688,9 +688,9 @@
 	// after the sentinels because a payload rejection carries none of
 	// them, and last because every other reason stays unmapped (report
 	// only — see .mule/todo.md).
 	a.writeBadObjectError(w, err)
-case payload.ReasonOf(err) == payload.ReasonMissingRequired:
+case payload.ReasonOf(err) == payload.ReasonMissingRequired && false:
 	// A distinct upstream failure from the object shape above, and the
 	// only other mapped payload reason: 422, because the document parsed
 	// and every key is declared — one of them the mapping declares
 	// required and the document omits (master b6d46ad4,

$ go test ./internal/appengine/ -run 'TestWriteError' -count=1 2>&1 | head -12
--- FAIL: TestWriteErrorTaxonomy (0.00s)
    --- FAIL: TestWriteErrorTaxonomy/object_write_omitting_a_required_key_answers_422 (0.00s)
        writeerror_test.go:59: status = 500, want 422 (body: {"errors":{"detail":"Internal server error"}})
FAIL
FAIL	github.com/astrate-platform/astrate/internal/appengine	0.020s
FAIL

← Edit internal/appengine/http.go
Index: /root/astrate-mule/internal/appengine/http.go
===================================================================
--- /root/astrate-mule/internal/appengine/http.go
+++ /root/astrate-mule/internal/appengine/http.go
@@ -688,9 +688,9 @@
 	// after the sentinels because a payload rejection carries none of
 	// them, and last because every other reason stays unmapped (report
 	// only — see .mule/todo.md).
 	a.writeBadObjectError(w, err)
-case payload.ReasonOf(err) == payload.ReasonMissingRequired && false:
+case payload.ReasonOf(err) == payload.ReasonMissingRequired:
 	// A distinct upstream failure from the object shape above, and the
 	// only other mapped payload reason: 422, because the document parsed
 	// and every key is declared — one of them the mapping declares
 	// required and the document omits (master b6d46ad4,

$ gofmt -l . ; go vet ./... && go test ./... 2>&1 | grep -v "^ok\|no test files" | head -10; echo "EXIT=$?"
EXIT=0

$ uname -m; go test -race ./internal/appengine/ -run TestWriteError -count=1 2>&1 | tail -3
aarch64
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine	0.024s
FAIL

$ tail -30 .mule/for-giulio.md
#78 is no longer parked: zero-touch onboarding is strategic for commercial
viability. Scope frozen on the issue: owner-side TO1/TO2 in our Pairing
service only (last mile, like upstream), reuse fdo-rs for
manufacturing/rendezvous, acceptance = official `astarte-device-fdo-rust`
SDK completes onboarding against Astrate, docs as a first-class deliverable.
When v3.0 is marked DONE, the v4.0 section of `.mule/milestones.md` should be
drafted with this investigation as its first item (issue #78 has the full
verified context).
- **The mule has been idle 17h.** Filed by the dead-man's switch; see journalctl on the Pi.

## 2026-09-25 — the `-tags nats` half of the tree is compiled by no gate

`cmd/astrate/newnats_nats.go`, `internal/engine/forward/nats.go` and
`internal/engine/forward/nats_test.go` are all behind `//go:build nats`, and nothing
builds or tests that tag: not `make build`/`make test`/`make test-integration`/
`make test-e2e` (Makefile:23, 38, 42, 46), not the mule gate, not the lint target.
So the whole NATS trigger-forwarding path is invisible to CI — while
`config.validate` accepts `triggers.forward.kind = "nats"` (internal/config/config.go:339-343)
and the untagged binary then fails at boot with the build-tag message
(cmd/astrate/newnats_default.go). I measured it today: `go build -tags nats ./cmd/...`
and `go vet -tags nats ./cmd/...` both exit 0, so there is nothing broken — just nothing
keeping it unbroken. Worth a `make build-nats` (and the NATS test tier, which needs a
container, so it would be a `[legion]` line) in the gate; not a queue line because the
only test I can imagine for it shells out to `go build`, and the mule's own check strips
the implementation and the guard together, so it would never fail.
Not queued for the same reason: `IdleTimeout` is missing on the HTTP server
(cmd/astrate/main.go:221) and `-healthcheck` ignores the config file
(cmd/astrate/main.go:531) — both real, both yours to call, neither has a test worth
writing. See `.mule/reviews/cmd-astrate-2026-09-25.md`.
- **The mule has been idle 16h.** Filed by the dead-man's switch; see journalctl on the Pi.

$ grep -c '^- \[ \]' .mule/todo.md; echo "=== last unchecked:"; grep -n '^- \[ \]' .mule/todo.md | tail -3 | cut -c1-160; echo "=== legion-worded precedent:"; grep -n 'Giulio\|legion' .mule/todo.md | grep '^- \[ \]' | cut -c1-180 | tail -5
15
=== last unchecked:
246:- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording �
259:- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) —
267:- [ ] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a ma
=== legion-worded precedent:

$ head -20 .mule/for-giulio.md; echo "=== date headings order:"; grep -n '^## ' .mule/for-giulio.md | head -5
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **github-issues triage run, 2026-10-03: nothing proposable for the 21st run — and the alarm pile finally produced a falsifiable demonstration of what the 2026-09-26 entry diagnosed, on both sides of two fresh alarms.** 27 open issues, the same set plus **#115** (10-02 11:22:50Z) and **#116** (today, 10-03 11:02:14Z). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`; `8c61268` re-verified still reachable from `origin/mule/queue` this run, so it is pushed and waiting only on your read), **#92** keyAgreement (parked on a stable v1.4.0 per your 2026-09-04 decision — today's milestone run re-verified there is still no stable v1.4.0, so the condition is unmet and the `waiting-on.md` gap is unchanged; see the entry below), **#78** FDO (`milestone-4.0`, already escalated), **#1** untouched per standing instruction.
  **Both new alarms are provably false, and this time the falsifier is in the log rather than inferred.** #115 says "nothing has landed in 14h", so it measured a gap ending ~10-01T21Z — the tail of 10-01's one `done` land (`41af3c4`, 479s, log.md:151). #116 says the same 14h, measuring back to ~10-02T21Z — the tail of 10-02's **six** `done` lands (`0215f97`, `f9af73b`, `32c85af`, `1d1e6f4`, `8a09a1a`, `2d1a1c1`, log.md:152-161, four `blocked` alongside). Every one of the ten shas from 10-01 to 10-03 is reachable from `origin/mule/queue` (`git branch -r --contains`), so they landed *and* pushed — the precise failure the alarm exists to catch, satisfied on all three days.
  **The 2026-10-01 "latched off, and the silence is real" reading is now superseded: the latch cleared itself the moment work landed, and the false alarm resumed on schedule.** That entry found `.alarmed` present with no land since 09-28T18:01Z, and that was true — the queue had genuinely stalled, on twelve `[legion]` lines with the Legion Go off. The 10-01 land beat, `beat` deleted `.alarmed` (mule.sh:737), `check_pulse` stopped short-circuiting (mule.sh:748), and the switch filed again the next morning. `.alarmed` is absent right now and the heartbeat reads 2026-10-03T17:05:22Z. **So the latch was a symptom of the stall, not a fault in the switch** — I withdraw the implicit "the switch cannot report a real outage" worry from that entry: nothing in `mule.sh` is broken, the 8h default (mule.sh:740) is simply wrong for a cadence that lands in the evening and polls in the morning, where 14h > 8h every day by construction.
  I did not re-argue the threshold: the 2026-09-26 entry already put both options to you (raise `MULE_IDLE_ALARM_HOURS` to ~24, or pulse `.mule/log.md` for a whole day with no `done` row instead of the heartbeat), and that question is still yours and still unanswered. This entry adds two measured data points to it and nothing else.
  **Proposal, extending the 2026-10-01 one: close #94–#116, all twenty-three.** Each was a one-day false alarm superseded by the next, and as established 2026-09-25 nothing reads issue state, so closing stays inert — but 23 open "the mule is dead" issues is precisely what buried a *real* three-day stall, and it will bury the next one. Leave #116 (live today) to self-expire.

---

=== date headings order:
85:## 2026-09-25 — docs-sync run, surface: **housekeeping** (the least-covered of the five)
490:## 2026-08-23 — FDO promoted to milestone-4.0 candidate (Giulio's decision, recorded)
502:## 2026-09-25 — the `-tags nats` half of the tree is compiled by no gate

$ sed -n '231,260p' tools/mule.sh
gates() {
  local rc=0
  mkdir -p "$MULE/failed"
  GATE_LOG="$MULE/failed/gate.log"
  : > "$GATE_LOG"
  GATE_REASON=""
  if [ -n "${MULE_FIX_CMD:-}" ]; then ( cd "$REPO" && eval "$MULE_FIX_CMD" ) >/dev/null 2>&1; fi
  if [ -n "${MULE_TEST_CMD:-}" ]; then
    note "gate: $MULE_TEST_CMD"
    ( cd "$REPO" && eval "$MULE_TEST_CMD" ) >>"$GATE_LOG" 2>&1 \
      || { gate_reason "tests failed"; bad "$GATE_REASON"; rc=1; }
  fi
  if [ "$rc" = 0 ] && [ -n "${MULE_LINT_CMD:-}" ]; then
    note "gate: lint"
    ( cd "$REPO" && eval "$MULE_LINT_CMD" ) >>"$GATE_LOG" 2>&1 \
      || { gate_reason "lint failed"; bad "$GATE_REASON"; rc=1; }
  fi
  if [ "$rc" = 0 ] && [ -z "${MULE_NO_PROOF_GATE:-}" ]; then
    proof_gate || rc=1
    # proof_gate marks untracked files intent-to-add so diff can see them, which leaves index
    # entries behind on every exit path. Nothing is deliberately staged during a task — the
    # runner stages after the gates — so clearing the index here is safe, and not clearing it
    # leaves a phantom deletion that the revert path would act on.
    git -C "$REPO" reset -q >/dev/null 2>&1 || true
  fi
  return $rc
}

# GNU sed takes -i with no argument; BSD/macOS sed requires an empty suffix after it. This
# script runs on both. Getting it wrong on Linux is silent and expensive: `sed -i '' 32s/...`

$ grep -n "MULE_TEST_CMD\|MULE_LINT_CMD" tools/mule.sh | head -8
195:  ( cd "$REPO" && eval "${MULE_TEST_CMD:-go test ./...}" ) >/dev/null 2>&1 && still_passes=1
238:  if [ -n "${MULE_TEST_CMD:-}" ]; then
239:    note "gate: $MULE_TEST_CMD"
240:    ( cd "$REPO" && eval "$MULE_TEST_CMD" ) >>"$GATE_LOG" 2>&1 \
243:  if [ "$rc" = 0 ] && [ -n "${MULE_LINT_CMD:-}" ]; then
245:    ( cd "$REPO" && eval "$MULE_LINT_CMD" ) >>"$GATE_LOG" 2>&1 \
300:  if [ -n "${MULE_TEST_CMD:-}" ]; then
301:    if ( cd "$REPO" && eval "$MULE_TEST_CMD" ) >/dev/null 2>&1; then ok "test baseline is green"

$ grep -rn "MULE_LINT_CMD\|MULE_TEST_CMD" .mule/config tools/*.sh 2>/dev/null | grep -v "^tools/mule.sh" | head -8
.mule/config:21:  MULE_TEST_CMD="go test -race ./..."
.mule/config:25:  MULE_TEST_CMD="go vet ./... && go test ./..."
.mule/config:29:  MULE_LINT_CMD="$_gcl run ./..."
.mule/config:32:  MULE_LINT_CMD=""            # absent on the Pi; gofmt still runs

→ Read .mule/config

$ command -v golangci-lint || ls ~/go/bin/golangci-lint 2>/dev/null || echo "absent"
/root/go/bin/golangci-lint

$ /root/go/bin/golangci-lint run ./internal/appengine/ 2>&1 | tail -10; echo "EXIT=$?"
0 issues.
EXIT=0

→ Read .mule/for-giulio.md [limit=4, offset=12]

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,10 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **The Legion Go answers ssh but its Docker is not running, so every "probe live upstream" line is blocked on you, not on the hardware (measured 2026-10-03).** `ssh -o ConnectTimeout=5 legion 'echo ok'` returns `ok`, but `docker ps` fails with `dial unix /home/atsetilam/.docker/desktop/docker.sock: no such file or directory`, `systemctl is-active docker` answers `inactive`, and `sudo -n true` answers `è necessaria una password` — so the mule cannot start it. It is Docker Desktop's socket, not a rootful `dockerd` the mule could launch as a user, and no rootless extras are installed. Consequence for the queue: `appengine-payload-reason-status-map` (the remaining eight `payload.Reason` values that still answer 500) cannot be measured the way its line demands. Today's `appengine-missing-required-422` was therefore mapped from the frozen upstream source at `b6d46ad4` (`Device.validate_required_mappings/2` → `fallback_controller.ex:215-219` → `ErrorView "422_missing_required_mapping.json"`), and the commit says so plainly. Nothing about that sibling line needs rewriting — only Docker up on the Legion Go.
+
 - **github-issues triage run, 2026-10-03: nothing proposable for the 21st run — and the alarm pile finally produced a falsifiable demonstration of what the 2026-09-26 entry diagnosed, on both sides of two fresh alarms.** 27 open issues, the same set plus **#115** (10-02 11:22:50Z) and **#116** (today, 10-03 11:02:14Z). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`; `8c61268` re-verified still reachable from `origin/mule/queue` this run, so it is pushed and waiting only on your read), **#92** keyAgreement (parked on a stable v1.4.0 per your 2026-09-04 decision — today's milestone run re-verified there is still no stable v1.4.0, so the condition is unmet and the `waiting-on.md` gap is unchanged; see the entry below), **#78** FDO (`milestone-4.0`, already escalated), **#1** untouched per standing instruction.
   **Both new alarms are provably false, and this time the falsifier is in the log rather than inferred.** #115 says "nothing has landed in 14h", so it measured a gap ending ~10-01T21Z — the tail of 10-01's one `done` land (`41af3c4`, 479s, log.md:151). #116 says the same 14h, measuring back to ~10-02T21Z — the tail of 10-02's **six** `done` lands (`0215f97`, `f9af73b`, `32c85af`, `1d1e6f4`, `8a09a1a`, `2d1a1c1`, log.md:152-161, four `blocked` alongside). Every one of the ten shas from 10-01 to 10-03 is reachable from `origin/mule/queue` (`git branch -r --contains`), so they landed *and* pushed — the precise failure the alarm exists to catch, satisfied on all three days.
   **The 2026-10-01 "latched off, and the silence is real" reading is now superseded: the latch cleared itself the moment work landed, and the false alarm resumed on schedule.** That entry found `.alarmed` present with no land since 09-28T18:01Z, and that was true — the queue had genuinely stalled, on twelve `[legion]` lines with the Legion Go off. The 10-01 land beat, `beat` deleted `.alarmed` (mule.sh:737), `check_pulse` stopped short-circuiting (mule.sh:748), and the switch filed again the next morning. `.alarmed` is absent right now and the heartbeat reads 2026-10-03T17:05:22Z. **So the latch was a symptom of the stall, not a fault in the switch** — I withdraw the implicit "the switch cannot report a real outage" worry from that entry: nothing in `mule.sh` is broken, the 8h default (mule.sh:740) is simply wrong for a cadence that lands in the evening and polls in the morning, where 14h > 8h every day by construction.
   I did not re-argue the threshold: the 2026-09-26 entry already put both options to you (raise `MULE_IDLE_ALARM_HOURS` to ~24, or pulse `.mule/log.md` for a whole day with no `done` row instead of the heartbeat), and that question is still yours and still unanswered. This entry adds two measured data points to it and nothing else.

$ git -C /root/astrate-mule diff -- internal/appengine/ | head -80
diff --git a/internal/appengine/http.go b/internal/appengine/http.go
index a39f559..532032a 100644
--- a/internal/appengine/http.go
+++ b/internal/appengine/http.go
@@ -689,6 +689,17 @@ func (a *API) writeError(w http.ResponseWriter, err error) {
 		// them, and last because every other reason stays unmapped (report
 		// only — see .mule/todo.md).
 		a.writeBadObjectError(w, err)
+	case payload.ReasonOf(err) == payload.ReasonMissingRequired:
+		// A distinct upstream failure from the object shape above, and the
+		// only other mapped payload reason: 422, because the document parsed
+		// and every key is declared — one of them the mapping declares
+		// required and the document omits (master b6d46ad4,
+		// Device.validate_required_mappings/2 answers
+		// :missing_required_mapping, the AppEngine fallback controller puts
+		// :unprocessable_entity for it, ErrorView renders the detail alone).
+		// Placed after bad_object because that one names the offending keys
+		// and answers 400; both precede default: so neither falls to the 500.
+		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, detailMissingRequiredMapping)
 	default:
 		_ = astarteapi.WriteInternalServerError(w)
 	}
@@ -699,6 +710,15 @@ func (a *API) writeError(w http.ResponseWriter, err error) {
 // astartectl and the SDKs match on it.
 const detailUnexpectedObjectKey = "Unexpected object key"
 
+// detailMissingRequiredMapping is upstream's 422 detail for an
+// object-aggregated write whose document omits a key the mapping declares
+// required (master b6d46ad4, ErrorView "422_missing_required_mapping.json"
+// renders `%{errors: %{detail: "Missing required mapping key"}}`). Frozen
+// string; unlike the undeclared-key case the body carries no key list —
+// upstream logs the offending key (device.ex:498-502) and renders the detail
+// alone.
+const detailMissingRequiredMapping = "Missing required mapping key"
+
 // writeBadObjectError renders the object-shape rejection in the measured
 // upstream envelope: {"errors": {"detail": "Unexpected object key",
 // "unexpected_keys": [...]}}. The other object-shape failures — not a
diff --git a/internal/appengine/writeerror_test.go b/internal/appengine/writeerror_test.go
index de8f821..bc3540f 100644
--- a/internal/appengine/writeerror_test.go
+++ b/internal/appengine/writeerror_test.go
@@ -44,6 +44,10 @@ func TestWriteErrorTaxonomy(t *testing.T) {
 		{"object write with no key list stays the canonical 400",
 			&payload.RejectError{Reason: payload.ReasonBadObject, Detail: "object-aggregation document is empty"},
 			http.StatusBadRequest, "Bad request", nil},
+		{"object write omitting a required key answers 422",
+			&payload.RejectError{Reason: payload.ReasonMissingRequired,
+				Detail: "object-aggregation document is missing required key(s) lat"},
+			http.StatusUnprocessableEntity, "Missing required mapping key", nil},
 		{"unknown cause stays 500", errors.New("boom"),
 			http.StatusInternalServerError, "Internal server error", nil},
 	}

$ git -C /root/astrate-mule status --porcelain; echo "=== gate:"; go vet ./... && go test ./... 2>&1 | grep -Ev "no test files" | tail -25
M .mule/for-giulio.md
 M internal/appengine/http.go
 M internal/appengine/writeerror_test.go
=== gate:
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

Done: mapped `payload.ReasonMissingRequired` to upstream's 422 `{"errors":{"detail":"Missing required mapping key"}}` in `writeError`, beside the `bad_object` case, with a table row.

Files: internal/appengine/http.go, internal/appengine/writeerror_test.go, .mule/for-giulio.md

Verified: `go vet ./... && go test ./...` -> pass (0 failures, run **without** `-race`: this host is aarch64 with the same 39-bit VMA as the Pi, `go test -race` dies `FATAL: Found 39 - Supported 48`; the change adds no concurrency); `gofmt -l .` -> clean; `golangci-lint run ./internal/appengine/` -> 0 issues; mutation check — disabling only the new case makes the new row fail with `status = 500, want 422`, matching the runner's proof gate. Upstream evidence, pinned at master b6d46ad4: `device.ex:490-503` returns `{:error, :missing_required_mapping}` for an object value missing a `required` key (the key goes to `Logger.warning`, not the body), `fallback_controller.ex:215-219` is the only clause matching it and puts `:unprocessable_entity`, `error_view.ex` `"422_missing_required_mapping.json"` renders `%{errors: %{detail: "Missing required mapping key"}}`; `interface_values_controller.ex:27` routes it there via `action_fallback`.

Unsure: **the live probe did not happen** — the Legion Go answers ssh but its Docker is down (`docker ps` fails on the missing `~/.docker/desktop/docker.sock`, `systemctl is-active docker` = inactive, `sudo -n true` = password required), so no upstream request was made and the status/detail come from reading the pinned commit, not from a measured response. Every intermediate is a literal frozen string and the one judgement call is that Phoenix/Jason renders the single-key map as `{"errors":{"detail":"Missing required mapping key"}}` with no extra field, which the source fixes unambiguously.

Follow-ups: none (nothing new to queue — `appengine-payload-reason-status-map` already covers the other eight reasons and needs no rewording; the Docker-down blocker is written up in `.mule/for-giulio.md` for you, since only you can start it).
