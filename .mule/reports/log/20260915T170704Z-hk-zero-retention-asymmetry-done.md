slug: hk-zero-retention-asymmetry
verdict: done
at:  015a706
ran: 2026-09-15T17:07:04Z on DietPi in 904s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

$ ls .mule/
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

→ Read internal/housekeeping/http.go
→ Read internal/housekeeping/service.go
Let me look at the existing tests and check upstream parity reference material.

$ ls internal/housekeeping/
async_operation_test.go
http.go
http_test.go
service.go
service_test.go

✱ Grep "retention" in internal/housekeeping · 43 matches
→ Read internal/housekeeping/http_test.go
→ Read internal/housekeeping/service_test.go

$ ls .mule/recipes/ .mule/tasks/ 2>/dev/null; ls .trickle* 2>/dev/null; ls -d .trickle 2>/dev/null
.mule/recipes/:
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

.mule/tasks/:
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md

Now let me verify upstream parity. The http_test.go references `.trickle-phase.md` for upstream measurements, and there's a `.mule/recipes` directory. Let me search for upstream evidence about retention=0.
✱ Grep "retention" in .trickle · 100 matches

$ ls -a /root/astrate-mule | head -50
.
..
.agents
.dockerignore
.git
.github
.gitignore
.golangci.yml
.mule
AGENTS.md
Dockerfile
LICENSE
Makefile
README.md
bench
cmd
deploy
docker-compose.yml
docs
examples
go.mod
go.sum
internal
migrations
opencode.json
pkg
site-dist
test
tools

→ Read docs/api/astarte_housekeeping_api.yaml

$ ls docs/handoff/ 2>/dev/null; rg -l "trickle" docs/ 2>/dev/null
README.md
flow-design-a-named-durable-flows-2026-07-29.md
flow-design-b-container-block-2026-07-29.md
flow-lifecycle-followups-design-2026-08-22.md
flow-parity-audit-2026-07-29.md
flow-v2-decisions-2026-07-29.md
phase-1-handoff.md
phase-1-memory.md
phase-10-handoff.md
phase-10-memory.md
phase-11-handoff.md
phase-2-handoff.md
phase-2-memory.md
phase-3-handoff.md
phase-3-memory.md
phase-4-handoff.md
phase-4-memory.md
phase-5-handoff.md
phase-5-memory.md
phase-6-handoff.md
phase-6-memory.md
phase-7-memory.md
phase-8-handoff.md
phase-8-memory.md
phase-9-handoff.md
phase-9-memory.md
pokemon-agent-handoff.md
pokemon-agent-memory.md
session-2026-07-29-flow-factory-handoff.md
session-2026-07-29-handoff.md

→ Read .mule/recipes/astarte-upstream.md

$ rg -n -i "upstream|version" docs/COMPATIBILITY.md | head -30
6:versions the conformance suite (`test/conformance`) pins and exercises.
14:| Client | Pinned version | Checkpoint(s) |
48:  application/merge-patch+json` (anything else reproduces upstream's unmapped
60:1. **Astarte Channels: two sockets, one bus** — the upstream Phoenix socket is
79:     Recorded against upstream v1.2.0 in
80:     `test/conformance/upstream/channels.json`. Upstream partitions the `a_ch`
88:     which stays a verb-regex match because that is what upstream's REST plug
93:     (`test/conformance/upstream/channels.json`). A data trigger is authorized
98:     it refused claims upstream accepts and accepted claims upstream refuses.
104:     upstream's changeset before authorization is consulted — and that a group
109:     reconciliation, which means an upstream-shaped group watch silently
112:     equal the request's own.** Upstream refuses a watch whose `device_id` is
122:   One further point, where the constraint is upstream's and Astrate conforms to
127:   surface with their own stability requirements, are translated to the upstream
128:   enum by `triggers.UpstreamErrorName` when an event body is built. The
132:   is **measured**, not reconstructed: it was recorded against upstream Astarte
133:   v1.2.0 in `test/conformance/upstream/channels.json`, and two rows were
135:   `write_on_server_owned_interface` is exact — upstream emits precisely that
137:   `value_size_exceeded` is real in upstream's value validator (astarte_core
146:   where upstream says nothing at all — kept as a deliberate divergence, since
148:   name itself is upstream's own enum value.
153:   superset; see [JSON-PAYLOAD-PROFILE.md](JSON-PAYLOAD-PROFILE.md). Upstream SDK
161:   with the upstream API surface.
163:5. **Uniform `401` where upstream distinguishes with `403`** — Astrate answers
169:     device both return `401`; upstream returns `403` via its RPC "forbidden"
174:     `auth.Middleware`. Upstream answers `403` here too. Observed against
176:     kind, so it is a recorded upstream response rather than a reading of its
178:     Astarte version other than the one that instance was running.
182:     upstream 1.2.0 on 2026-08-24 (`test/conformance/upstream/verify-versions.json`,
183:     issue #69): upstream answers `403 Forbidden` for any non-absent

$ gh --version 2>/dev/null; curl -sI --max-time 8 https://raw.githubusercontent.com/astarte-platform/astarte_housekeeping/master/README.md 2>&1 | head -3
gh version 2.100.0 (2026-09-03)
https://github.com/cli/cli/releases/tag/v2.100.0
HTTP/2 200 
cache-control: max-age=300
content-security-policy: default-src 'none'; style-src 'unsafe-inline'; sandbox

$ curl -s --max-time 15 "https://api.github.com/repos/astarte-platform/astarte_housekeeping/git/trees/master?recursive=1" | rg -o '"path": "[^"]*"' | rg -i "realm" | head -40
"path": "priv/migrations/realm"
"path": "priv/migrations/realm/0001_create_grouped_devices_table.sql"
"path": "priv/migrations/realm/0002_add_groups_and_interface_stats_columns.sql"
"path": "priv/migrations/realm/0003_add_database_retention_ttl_and_policy.sql"

$ curl -s --max-time 15 "https://api.github.com/repos/astarte-platform/astarte_housekeeping/git/trees/master?recursive=1" | rg -o '"path": "[^"]*"' | rg -iv "priv/migrations|rel/" | head -60
"path": ".formatter.exs"
"path": ".gitignore"
"path": ".gitlab-ci.yml"
"path": ".travis.yml"
"path": "AUTHORS.md"
"path": "CHANGELOG.md"
"path": "Dockerfile"
"path": "LICENSE"
"path": "README.md"
"path": "config"
"path": "config/astarte_housekeeping.dev.conf"
"path": "config/astarte_housekeeping.prod.conf"
"path": "config/astarte_housekeeping.schema.exs"
"path": "config/ci.exs"
"path": "config/config.exs"
"path": "config/dev.exs"
"path": "config/prod.exs"
"path": "config/test.exs"
"path": "coveralls.json"
"path": "entrypoint.sh"
"path": "lib"
"path": "lib/astarte_housekeeping.ex"
"path": "lib/astarte_housekeeping"
"path": "lib/astarte_housekeeping/config.ex"
"path": "lib/astarte_housekeeping/engine.ex"
"path": "lib/astarte_housekeeping/migrator.ex"
"path": "lib/astarte_housekeeping/queries.ex"
"path": "lib/astarte_housekeeping/release_tasks.ex"
"path": "lib/astarte_housekeeping/rpc"
"path": "lib/astarte_housekeeping/rpc/handler.ex"
"path": "mix.exs"
"path": "mix.lock"
"path": "priv"
"path": "rel"
"path": "test"
"path": "test/astarte_housekeeping"
"path": "test/astarte_housekeeping/queries_test.exs"
"path": "test/astarte_housekeeping/rpc"
"path": "test/astarte_housekeeping/rpc/amqp_server_test.exs"
"path": "test/support"
"path": "test/support/database_test_helper.ex"
"path": "test/test_helper.exs"

$ curl -s --max-time 15 "https://api.github.com/repos/astarte-platform/astarte_housekeeping_api/git/trees/master?recursive=1" | rg -o '"path": "[^"]*"' | rg -i "controller|realm" | head -40
"path": "lib/astarte_housekeeping_api/realms"
"path": "lib/astarte_housekeeping_api/realms/realm.ex"
"path": "lib/astarte_housekeeping_api/realms/realms.ex"
"path": "lib/astarte_housekeeping_api_web/controllers"
"path": "lib/astarte_housekeeping_api_web/controllers/fallback_controller.ex"
"path": "lib/astarte_housekeeping_api_web/controllers/health_controller.ex"
"path": "lib/astarte_housekeeping_api_web/controllers/realm_controller.ex"
"path": "lib/astarte_housekeeping_api_web/views/realm_view.ex"
"path": "test/astarte_housekeeping_api/realms"
"path": "test/astarte_housekeeping_api/realms/realms_test.exs"
"path": "test/astarte_housekeeping_api_web/controllers"
"path": "test/astarte_housekeeping_api_web/controllers/realm_controller_test.exs"

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte_housekeeping_api/master/lib/astarte_housekeeping_api/realms/realm.ex"
#
# This file is part of Astarte.
#
# Copyright 2017-2018 Ispirata Srl
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

defmodule Astarte.Housekeeping.API.Realms.Realm do
  use Ecto.Schema
  import Ecto.Changeset

  @default_replication_factor 1
  @datacenter_name_regex ~r/^[a-z][a-zA-Z0-9_-]*$/

  @required_fields [:realm_name, :jwt_public_key_pem]
  @allowed_fields [
    :replication_factor,
    :replication_class,
    :datacenter_replication_factors
    | @required_fields
  ]

  @primary_key false
  @derive {Phoenix.Param, key: :realm_name}

  embedded_schema do
    field :realm_name
    field :jwt_public_key_pem
    field :replication_factor, :integer
    field :replication_class, :string, default: "SimpleStrategy"
    field :datacenter_replication_factors, {:map, :integer}
  end

  def changeset(realm, params \\ %{}) do
    realm
    |> cast(params, @allowed_fields)
    |> validate_required(@required_fields)
    |> validate_format(:realm_name, ~r/^[a-z][a-z0-9]*$/)
    |> validate_number(:replication_factor, greater_than: 0)
    |> validate_pem_public_key(:jwt_public_key_pem)
    |> validate_inclusion(:replication_class, ["SimpleStrategy", "NetworkTopologyStrategy"])
    |> validate_replication()
  end

  def error_changeset(realm, params \\ %{}) do
    changeset =
      realm
      |> cast(params, @required_fields)

    %{changeset | valid?: false}
  end

  defp datacenter_map_validator(field, datacenter_map) when map_size(datacenter_map) == 0 do
    [{field, "must not be empty"}]
  end

  defp datacenter_map_validator(field, datacenter_map) do
    Enum.reduce(datacenter_map, [], fn {datacenter_name, replication_factor}, errors_acc ->
      cond do
        not Regex.match?(@datacenter_name_regex, datacenter_name) ->
          [{field, "has invalid datacenter name: #{datacenter_name}"} | errors_acc]

        not is_number(replication_factor) or replication_factor <= 0 ->
          [{field, "has invalid replication factor: #{replication_factor}"} | errors_acc]

        true ->
          errors_acc
      end
    end)
  end

  defp validate_replication(changeset) do
    replication_class = get_field(changeset, :replication_class)

    cond do
      replication_class == "NetworkTopologyStrategy" ->
        changeset
        |> validate_required(:datacenter_replication_factors,
          message: "needs to be provided with NetworkTopologyStrategy"
        )
        |> validate_change(:datacenter_replication_factors, &datacenter_map_validator/2)

      # Here we're implicitly not in NetworkTopologyStrategy
      get_field(changeset, :datacenter_replication_factors) ->
        changeset
        |> add_error(
          :datacenter_replication_factors,
          "must be used with replication_class NetworkTopologyStrategy"
        )

      get_field(changeset, :replication_factor) == nil ->
        changeset
        |> put_change(:replication_factor, @default_replication_factor)

      true ->
        changeset
    end
  end

  defp validate_pem_public_key(%Ecto.Changeset{valid?: false} = changeset, _field), do: changeset

  defp validate_pem_public_key(changeset, field) do
    pem = get_field(changeset, field, "")

    try do
      case :public_key.pem_decode(pem) do
        [{:SubjectPublicKeyInfo, _, _}] ->
          changeset

        _ ->
          changeset
          |> add_error(field, "is not a valid PEM public key")
      end
    rescue
      _ ->
        changeset
        |> add_error(field, "is not a valid PEM public key")
    end
  end
end

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte_housekeeping_api/master/lib/astarte_housekeeping_api_web/controllers/realm_controller.ex"
#
# This file is part of Astarte.
#
# Copyright 2017 Ispirata Srl
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

defmodule Astarte.Housekeeping.APIWeb.RealmController do
  use Astarte.Housekeeping.APIWeb, :controller

  alias Astarte.Housekeeping.API.Realms
  alias Astarte.Housekeeping.API.Realms.Realm

  action_fallback Astarte.Housekeeping.APIWeb.FallbackController

  plug Astarte.Housekeeping.APIWeb.Plug.AuthorizePath

  def index(conn, _params) do
    realms = Realms.list_realms()
    render(conn, "index.json", realms: realms)
  end

  def create(conn, %{"data" => realm_params}) do
    with {:ok, %Realm{} = realm} <- Realms.create_realm(realm_params) do
      conn
      |> put_status(:created)
      |> put_resp_header("location", realm_path(conn, :show, realm))
      |> render("show.json", realm: realm)
    end
  end

  def show(conn, %{"id" => id}) do
    with {:ok, %Realm{} = realm} <- Realms.get_realm(id) do
      render(conn, "show.json", realm: realm)
    end
  end

  def update(conn, %{"id" => id, "data" => realm_params}) do
    with {:ok, %Realm{} = realm} <- Realms.get_realm(id),
         {:ok, %Realm{} = realm} <- Realms.update_realm(realm, realm_params) do
      render(conn, "show.json", realm: realm)
    end
  end

  def delete(conn, %{"id" => id}) do
    with {:ok, %Realm{} = realm} <- Realms.get_realm(id),
         {:ok, %Realm{}} <- Realms.delete_realm(realm) do
      send_resp(conn, :no_content, "")
    end
  end
end

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte_housekeeping_api/master/lib/astarte_housekeeping_api/realms/realms.ex"
#
# This file is part of Astarte.
#
# Copyright 2017-2018 Ispirata Srl
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

defmodule Astarte.Housekeeping.API.Realms do
  @moduledoc """
  The boundary for the Realms system.
  """

  alias Astarte.Housekeeping.API.Realms.Realm
  alias Astarte.Housekeeping.API.RPC.Housekeeping

  @doc """
  Returns the list of realms.

  ## Examples

      iex> list_realms()
      [%Realm{}, ...]

  """
  def list_realms do
    Housekeeping.list_realms()
  end

  @doc """
  Gets a single realm.

  ## Examples

      iex> get_realm!(123)
      %Realm{}

  """
  def get_realm(realm_name) do
    Housekeeping.get_realm(realm_name)
  end

  @doc """
  Creates a realm.

  ## Examples

      iex> create_realm(%{field: value})
      {:ok, %Realm{}}

      iex> create_realm(%{field: bad_value})
      {:error, ...}

  """
  def create_realm(attrs \\ %{}) do
    changeset =
      %Realm{}
      |> Realm.changeset(attrs)

    with {:ok, %Realm{} = realm} <- Ecto.Changeset.apply_action(changeset, :insert) do
      case Housekeeping.create_realm(realm) do
        :ok -> {:ok, realm}
        {:ok, :started} -> {:ok, realm}
        {:error, reason} -> {:error, reason}
      end
    end
  end

  @doc """
  Updates a realm.

  ## Examples

      iex> update_realm(realm, %{field: new_value})
      {:ok, %Realm{}}

      iex> update_realm(realm, %{field: bad_value})
      {:error, ...}

  """
  def update_realm(%Realm{} = _realm, _attrs) do
    raise "TODO"
  end

  @doc """
  Deletes a Realm.

  ## Examples

      iex> delete_realm(realm)
      {:ok, %Realm{}}

      iex> delete_realm(realm)
      {:error, ...}

  """
  def delete_realm(%Realm{} = _realm) do
    raise "TODO"
  end

  @doc """
  Returns a datastructure for tracking realm changes.

  ## Examples

      iex> change_realm(realm)
      %Todo{...}

  """
  def change_realm(%Realm{} = _realm) do
    raise "TODO"
  end
end

$ curl -s --max-time 15 "https://api.github.com/repos/astarte-platform/astarte_housekeeping_api/git/trees/master?recursive=1" | rg -o '"path": "[^"]*"' | rg -i "swagger|openapi|spec"
(no output)

$ curl -s --max-time 15 "https://api.github.com/repos/astarte-platform/astarte/git/trees/master?recursive=1" | rg -o '"path": "[^"]*"' | rg -i "realm_management.*(config|handler)|housekeeping.*realm" | head
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping/realms"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping/realms/core.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping/realms/queries.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping/realms/realm.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping/realms/realms.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web/api_spec/schemas/realm.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web/controllers/realm_controller.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web/views/realm_view.ex"
"path": "apps/astarte_housekeeping/test/astarte_housekeeping/realms"

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/master/apps/astarte_housekeeping/lib/astarte_housekeeping/realms/realm.ex"
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

defmodule Astarte.Housekeeping.Realms.Realm do
  @moduledoc false
  use TypedEctoSchema

  import Ecto.Changeset

  alias Astarte.Core.Realm, as: CoreRealm
  alias Astarte.Housekeeping.Realms.NonNegativeIntegerOrUnsetType

  require Logger

  @required_create_fields [:realm_name, :jwt_public_key_pem]
  @allowed_create_fields [
    :replication_factor,
    :replication_class,
    :datacenter_replication_factors,
    :device_registration_limit,
    :datastream_maximum_storage_retention
    | @required_create_fields
  ]

  @allowed_update_fields [
    :jwt_public_key_pem,
    :device_registration_limit,
    :datastream_maximum_storage_retention
  ]

  @primary_key false
  @derive {Phoenix.Param, key: :realm_name}

  typed_embedded_schema do
    field :realm_name, :string
    field :jwt_public_key_pem, :string
    field :replication_factor, :integer
    field :replication_class, :string
    field :datacenter_replication_factors, {:map, :integer}
    field :device_registration_limit, NonNegativeIntegerOrUnsetType
    field :datastream_maximum_storage_retention, NonNegativeIntegerOrUnsetType
  end

  def changeset(realm, params \\ %{}) do
    realm
    |> cast(params, @allowed_create_fields)
    |> validate_required(@required_create_fields)
    |> validate_change(:realm_name, &validate_realm_name/2)
    |> validate_number(:replication_factor, greater_than: 0)
    |> validate_change(:jwt_public_key_pem, &validate_pem_public_key/2)
    |> validate_inclusion(:replication_class, [
      "SimpleStrategy",
      "NetworkTopologyStrategy"
    ])
    |> validate_replication()
    |> maybe_put_default_replication_factor()
  end

  def update_changeset(realm, params \\ %{}) do
    realm
    |> cast(params, @allowed_update_fields)
    |> validate_change(:jwt_public_key_pem, &validate_pem_public_key/2)
  end

  def error_changeset(realm, params \\ %{}) do
    changeset = cast(realm, params, @required_create_fields)

    %{changeset | valid?: false}
  end

  defp datacenter_map_validator(field, datacenter_map) when map_size(datacenter_map) == 0 do
    [{field, "must not be empty"}]
  end

  defp datacenter_map_validator(field, datacenter_map) do
    Enum.reduce(datacenter_map, [], fn {_datacenter_name, replication_factor}, errors_acc ->
      if is_number(replication_factor) and replication_factor > 0 do
        errors_acc
      else
        [{field, "has invalid replication factor: #{replication_factor}"} | errors_acc]
      end
    end)
  end

  defp validate_replication(changeset) do
    replication_class = get_field(changeset, :replication_class)

    cond do
      replication_class == "NetworkTopologyStrategy" ->
        changeset
        |> validate_required(:datacenter_replication_factors,
          message: "needs to be provided with NetworkTopologyStrategy"
        )
        |> validate_change(:datacenter_replication_factors, &datacenter_map_validator/2)

      # Here we're implicitly not in NetworkTopologyStrategy
      get_field(changeset, :datacenter_replication_factors) ->
        add_error(
          changeset,
          :datacenter_replication_factors,
          "must be used with replication_class NetworkTopologyStrategy"
        )

      true ->
        changeset
    end
  end

  def apply_replication(%Astarte.Housekeeping.Realms.Realm{} = realm, %{
        strategy: :network_topology,
        dc_factors: dc_factors
      }) do
    dc_factors_map = Enum.into(Enum.to_list(dc_factors), %{})

    %{
      realm
      | replication_class: "NetworkTopologyStrategy",
        datacenter_replication_factors: dc_factors_map
    }
  end

  def apply_replication(%Astarte.Housekeeping.Realms.Realm{} = realm, %{
        strategy: :simple,
        factor: factor
      }) do
    %{realm | replication_class: "SimpleStrategy", replication_factor: factor}
  end

  defp maybe_put_default_replication_factor(changeset) do
    replication_class = get_field(changeset, :replication_class)

    if replication_class == "SimpleStrategy" do
      put_default_if_missing(changeset, :replication_factor, 1)
    else
      changeset
    end
  end

  def validate_realm_name(field, realm_name) do
    case CoreRealm.valid_name?(realm_name) do
      true ->
        []

      false ->
        Logger.warning("Invalid realm name.", tag: "invalid_realm_name", realm: realm_name)
        [{field, "has invalid format"}]
    end
  end

  defp validate_pem_public_key(field, pem) do
    case :public_key.pem_decode(pem) do
      [{:SubjectPublicKeyInfo, _, _}] ->
        []

      _ ->
        [{field, "is not a valid PEM public key"}]
    end
  rescue
    _ ->
      [{field, "is not a valid PEM public key"}]
  end

  defp put_default_if_missing(changeset, field, default) do
    if field_missing?(changeset, field) do
      put_change(changeset, field, default)
    else
      changeset
    end
  end
end

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/master/apps/astarte_housekeeping_web/api_spec/schemas/realm.ex"
404: Not Found

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/master/apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex"
#
# This file is part of Astarte.
#
# Copyright 2024 - 2025 SECO Mind srl
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

defmodule Astarte.Housekeeping.Realms.NonNegativeIntegerOrUnsetType do
  @moduledoc false
  use Ecto.Type

  @type t() :: non_neg_integer() | :unset

  def type, do: :any

  def cast(:unset), do: {:ok, :unset}

  def cast(n) when is_integer(n) do
    if n >= 0, do: {:ok, n}, else: :error
  end

  def cast(_), do: :error

  # load and dump are not meaningful for our use-case
  # coveralls-ignore-start
  def load(:unset), do: {:ok, :unset}
  def load(n) when is_integer(n), do: {:ok, n}
  def load(_), do: :error
  def dump(:unset), do: {:ok, :unset}
  def dump(n) when is_integer(n), do: {:ok, n}
  def dump(_), do: :error
end

$ curl -s --max-time 15 "https://api.github.com/repos/astarte-platform/astarte/git/trees/master?recursive=1" | rg -o '"path": "[^"]*"' | rg -i "housekeeping_web/api_spec|housekeeping.*schemas"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web/api_spec"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web/api_spec/api_spec.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web/api_spec/schemas"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web/api_spec/schemas/errors.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web/api_spec/schemas/realm.ex"
"path": "apps/astarte_housekeeping/test/astarte_housekeeping_web/api_spec"
"path": "apps/astarte_housekeeping/test/astarte_housekeeping_web/api_spec/api_spec_test.exs"

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/master/apps/astarte_housekeeping/lib/astarte_housekeeping_web/api_spec/schemas/realm.ex"
defmodule Astarte.HousekeepingWeb.ApiSpec.Schemas.Realm do
  @moduledoc false

  defmodule Realm do
    @moduledoc false

    require OpenApiSpex

    alias OpenApiSpex.Schema

    @public_key_example """
    -----BEGIN PUBLIC KEY-----
    MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsj7/Ci5Nx+ApLNW7+DyE
    eTzQ68KEJT/gPW73Kpa2uyvxDwY669z/rP4hMj16wv4Ku3bI6C1ZIqT5SVuF8pDo
    1Y1SF0GRIeslupm9KV1aFqIu1/srLz18LQHucQYUSa99PStFUJY2V83wneaeAArY
    4VKDuQYtRZOd2VeD5Cbn602ksLLWCQc9HfL3VUHXTw6DuthnMMJARcVem8RAMScm
    htGi6YRPFzvHtkb1WQCNGjw5gAmHX5/37ouwbBdnXOa9deiFv+1UIdcCVwMTyP/4
    f9jgaxW4oQV85enS/OJrrC9jU11agRc4bDv1h4s2t+ETWb4llTVk3HMIHbC3EvKJ
    VwIDAQAB
    -----END PUBLIC KEY-----
    """

    OpenApiSpex.schema(%{
      title: "Realm",
      type: :object,
      required: [:realm_name, :jwt_public_key_pem],
      properties: %{
        realm_name: %Schema{type: :string, example: "myrealm"},
        jwt_public_key_pem: %Schema{
          type: :string,
          example: @public_key_example,
          description: "PEM-encoded public key of the realm."
        },
        replication_class: %Schema{
          type: :string,
          example: "SimpleStrategy",
          description: "Replication Class of the keyspace that holds the realm's data."
        },
        replication_factor: %Schema{
          type: :integer,
          example: 2,
          description:
            "Replication factor of the keyspace that holds the realm's data (only if replication_class is \"SimpleStrategy\")."
        },
        datacenter_replication_factor: %Schema{
          type: :object,
          additionalProperties: %Schema{type: :integer},
          example: %{
            datacenter_1: 1,
            datacenter_2: 3
          },
          description:
            "Datacenter replication factor of the keyspace that holds the realm's data (only if replication_class is \"NetworkTopologyStrategy\")."
        },
        device_registration_limit: %Schema{
          type: :integer,
          minimum: 0,
          example: 100,
          description:
            "Optional upper bound to the number of devices that can be registered in the realm."
        },
        datastream_maximum_storage_retention: %Schema{
          type: :integer,
          minimum: 1,
          example: 100,
          description:
            "Optional upper bound to the retention period of all datastreams in the realm, in seconds."
        }
      },
      example: %{
        realm_name: "myrealm",
        jwt_public_key_pem: @public_key_example,
        replication_class: "SimpleStrategy",
        replication_factor: 2
      }
    })
  end

  defmodule RealmPatch do
    @moduledoc false

    require OpenApiSpex

    alias OpenApiSpex.Schema

    @public_key_example """
    -----BEGIN PUBLIC KEY-----
    MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsj7/Ci5Nx+ApLNW7+DyE
    eTzQ68KEJT/gPW73Kpa2uyvxDwY669z/rP4hMj16wv4Ku3bI6C1ZIqT5SVuF8pDo
    1Y1SF0GRIeslupm9KV1aFqIu1/srLz18LQHucQYUSa99PStFUJY2V83wneaeAArY
    4VKDuQYtRZOd2VeD5Cbn602ksLLWCQc9HfL3VUHXTw6DuthnMMJARcVem8RAMScm
    htGi6YRPFzvHtkb1WQCNGjw5gAmHX5/37ouwbBdnXOa9deiFv+1UIdcCVwMTyP/4
    f9jgaxW4oQV85enS/OJrrC9jU11agRc4bDv1h4s2t+ETWb4llTVk3HMIHbC3EvKJ
    VwIDAQAB
    -----END PUBLIC KEY-----
    """

    OpenApiSpex.schema(%{
      title: "RealmPatch",
      type: :object,
      properties: %{
        jwt_public_key_pem: %Schema{
          type: :string,
          example: @public_key_example,
          description: "PEM-encoded public key of the realm."
        },
        device_registration_limit: %Schema{
          type: :integer,
          minimum: 0,
          example: 100,
          description:
            "Optional upper bound to the number of devices that can be registered in the realm."
        },
        datastream_maximum_storage_retention: %Schema{
          type: :integer,
          minimum: 1,
          example: 100,
          description:
            "Optional upper bound to the retention period of all datastreams in the realm, in seconds."
        }
      },
      example: %{
        jwt_public_key_pem: @public_key_example
      }
    })
  end
end

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/master/apps/astarte_housekeeping/lib/astarte_housekeeping_web/controllers/realm_controller.ex"
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
#

defmodule Astarte.HousekeepingWeb.RealmController do
  use Astarte.HousekeepingWeb, :controller
  use OpenApiSpex.ControllerSpecs

  alias Astarte.Housekeeping.Config
  alias Astarte.Housekeeping.Realms
  alias Astarte.Housekeeping.Realms.Realm
  alias OpenApiSpex.{Reference, Schema}

  action_fallback Astarte.HousekeepingWeb.FallbackController

  tags(["realm"])
  security([%{"JWT" => []}])

  operation :index,
    summary: "Get all realms",
    description: "Returns a list of all existing realms.",
    operation_id: "getRealms",
    responses: [
      created:
        {"Success", "application/json",
         %Schema{
           type: :object,
           properties: %{
             data: %Schema{
               type: :array,
               items: %Schema{type: :string}
             }
           },
           example: %{
             data: ["arealm", "anotherrealm"]
           }
         }}
    ]

  operation :create,
    summary: "Create a realm",
    description: "Creates a new realm, based on the provided realm configuration.
        Realm creation will be executed asynchronously by default - it is not
        guaranteed that the requested realm will be available as soon as the
        API call returns, but it is guaranteed that it will be eventually created
        if no errors are returned and Astarte is operating normally.
        You can perform the call synchronously by setting the async_operation query
        param to false.",
    operation_id: "createRealm",
    parameters: [
      async_operation: [
        in: :query,
        description: "Whether the operation should be carried out asynchronously.",
        required: false,
        schema: %Schema{type: :boolean, default: true}
      ]
    ],
    request_body: %Reference{"$ref": "#/components/requestBodies/createRealmBody"},
    responses: [
      ok:
        {"Success", "application/json",
         %Schema{
           type: :object,
           properties: %{
             data: %Reference{"$ref": "#/components/schemas/Realm"}
           }
         }}
    ]

  operation :show,
    summary: "Get realm configuration",
    description: "Retrieves a realm's configuration.",
    operation_id: "getRealmConfiguration",
    parameters: [
      realm_name: [
        in: :path,
        description: "Realm name",
        required: true,
        schema: %Schema{type: :string}
      ]
    ],
    responses: [
      ok:
        {"Success", "application/json",
         %Schema{
           type: :object,
           properties: %{
             data: %Reference{"$ref": "#/components/schemas/Realm"}
           }
         }},
      unauthorized: %Reference{"$ref": "#/components/responses/Unauthorized"},
      forbidden: %Reference{"$ref": "#/components/responses/AuthorizationPathNotMatched"}
    ]

  operation :update,
    summary: "Update a realm",
    description: "Updates a realm's configuration.",
    operation_id: "updateRealm",
    parameters: [
      realm_name: [
        name: :realm_name,
        in: :path,
        description: "Realm name",
        required: true,
        schema: %Schema{type: :string}
      ]
    ],
    request_body: %Reference{"$ref": "#/components/requestBodies/updateRealmBody"},
    responses: [
      ok:
        {"Success", "application/json",
         %Schema{
           type: :object,
           properties: %{
             data: %Reference{"$ref": "#/components/schemas/Realm"}
           }
         }},
      bad_request: "Bad request",
      unauthorized: %Reference{"$ref": "#/components/responses/Unauthorized"},
      forbidden: %Reference{"$ref": "#/components/responses/AuthorizationPathNotMatched"},
      not_found:
        {"Realm not found", "application/json",
         %Reference{"$ref": "#/components/schemas/GenericError"}}
    ]

  operation :delete,
    summary: "Delete realm",
    description: "Deletes a realm from Astarte. This feature must be explicitly enabled
        in the cluster, if it's disabled a 405 status code will be returned.
        If there are connected devices present in the realm, a 422 status
        code will be returned. Realm deletion will be executed asynchronously
        by default - it is not guaranteed that the realm will be deleted as
        soon as the API call returns, but it is guaranteed that it will be
        eventually removed if no errors are returned and Astarte is
        operating normally. You can perform the call synchronously by setting
        the async_operation parameter to false.",
    operation_id: "deleteRealm",
    parameters: [
      realm_name: [
        name: :realm_name,
        in: :path,
        description: "Realm name",
        required: true,
        schema: %Schema{type: :string}
      ],
      async_operation: [
        in: :query,
        description: "Whether the operation should be carried out asynchronously.",
        required: false,
        schema: %Schema{type: :boolean, default: true}
      ]
    ],
    responses: [
      no_content: "Success",
      unauthorized: %Reference{"$ref": "#/components/responses/Unauthorized"},
      forbidden: %Reference{"$ref": "#/components/responses/AuthorizationPathNotMatched"},
      method_not_allowed: "Realm deletion disabled",
      unprocessable_entity: "Connected devices present"
    ]

  operation :get_default_replication,
    summary: "Get default Realm replication settings",
    description: "Retrieves the default DB replication settings applied to new Realms.",
    operation_id: "getDefaultReplication",
    responses: [
      ok:
        {"Success", "application/json",
         %Reference{"$ref": "#/components/schemas/DefaultReplication"}},
      unauthorized: %Reference{"$ref": "#/components/responses/Unauthorized"},
      forbidden: %Reference{"$ref": "#/components/responses/AuthorizationPathNotMatched"},
      service_unavailable: "Astarte is unable to retrieve the replication settings."
    ]

  def index(conn, _params) do
    with {:ok, realms} <- Realms.list_realms() do
      render(conn, "index.json", realms: realms)
    end
  end

  def create(conn, %{"data" => realm_params} = params) do
    async_operation =
      if Map.get(params, "async_operation") == "false" do
        false
      else
        true
      end

    with {:ok, %Realm{} = realm} <-
           Realms.create_realm(realm_params, async_operation: async_operation) do
      conn
      |> put_status(:created)
      |> put_resp_header("location", realm_path(conn, :show, realm))
      |> render("show.json", realm: realm)
    end
  end

  def show(conn, %{"realm_name" => realm_name}) do
    with {:ok, %Realm{} = realm} <- Realms.get_realm(realm_name) do
      render(conn, "show.json", realm: realm)
    end
  end

  def update(%Plug.Conn{method: "PATCH"} = conn, %{
        "realm_name" => realm_name,
        "data" => realm_params
      }) do
    update_params = normalize_update_attrs(realm_params)

    with {:ok, %Realm{} = updated_realm} <- Realms.update_realm(realm_name, update_params) do
      render(conn, "show.json", realm: updated_realm)
    end
  end

  def delete(conn, %{"realm_name" => realm_name} = params) do
    async_operation =
      if Map.get(params, "async_operation") == "false" do
        false
      else
        true
      end

    with :ok <- Realms.delete_realm(realm_name, async_operation: async_operation) do
      send_resp(conn, :no_content, "")
    end
  end

  def get_default_replication(conn, _params) do
    # TODO: this function needs to retrieve the default replication from the db
    # to change once the appropriate pr has been merged
    replication =
      case Config.astarte_keyspace_replication_strategy!() do
        :simple_strategy ->
          %{
            replication_class: :simple_strategy,
            replication_factor: Config.astarte_keyspace_replication_factor!()
          }

        :network_topology_strategy ->
          %{
            replication_class: :network_topology_strategy,
            datacenter_replication_factor: Config.astarte_keyspace_network_replication_map!()
          }
      end

    render(conn, "replication.json", replication: replication)
  end

  defp normalize_update_attrs(update_attrs) when is_map(update_attrs) do
    update_attrs
    |> Map.replace_lazy(:device_registration_limit, &normalize_integer_or_nil/1)
    |> Map.replace_lazy("device_registration_limit", &normalize_integer_or_nil/1)
    |> Map.replace_lazy(:datastream_maximum_storage_retention, &normalize_integer_or_nil/1)
    |> Map.replace_lazy("datastream_maximum_storage_retention", &normalize_integer_or_nil/1)
  end

  defp normalize_integer_or_nil(value) when is_nil(value), do: :unset
  defp normalize_integer_or_nil(value), do: value
end

Let me check upstream's router, view, and queries to confirm how 0 is handled end-to-end.
$ for f in "apps/astarte_housekeeping/lib/astarte_housekeeping_web/router.ex" "apps/astarte_housekeeping/lib/astarte_housekeeping_web/views/realm_view.ex" "apps/astarte_housekeeping/lib/astarte_housekeeping/realms/queries.ex"; do echo "===== $f ====="; curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/master/$f"; echo; done
===== apps/astarte_housekeeping/lib/astarte_housekeeping_web/router.ex =====
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
#

defmodule Astarte.HousekeepingWeb.Router do
  use Astarte.HousekeepingWeb, :router

  pipeline :api do
    plug :accepts, ["json"]
    plug Astarte.HousekeepingWeb.Plug.AuthorizePath
  end

  scope "/v1", Astarte.HousekeepingWeb do
    pipe_through :api

    get "/version", VersionController, :show
    get "/realm-defaults/replication", RealmController, :get_default_replication

    resources "/realms", RealmController,
      except: [:new, :edit, :update],
      param: "realm_name"

    patch "/realms/:realm_name", RealmController, :update
  end

  scope "/version", Astarte.HousekeepingWeb do
    get "/", VersionController, :show
  end
end

===== apps/astarte_housekeeping/lib/astarte_housekeeping_web/views/realm_view.ex =====
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

defmodule Astarte.HousekeepingWeb.RealmView do
  use Astarte.HousekeepingWeb, :view

  alias Astarte.Housekeeping.Realms.Realm
  alias Astarte.HousekeepingWeb.RealmView

  def render("index.json", %{realms: realms}) do
    %{data: render_many(realms, RealmView, "realm_name_only.json")}
  end

  def render("show.json", %{realm: realm}) do
    render_one(realm, RealmView, "realm.json")
  end

  def render("realm_name_only.json", %{realm: realm}) do
    realm.realm_name
  end

  def render("replication.json", %{replication: replication}) do
    %{data: replication}
  end

  def render("realm.json", %{realm: %Realm{replication_class: "SimpleStrategy"} = realm}) do
    %{
      data: %{
        realm_name: realm.realm_name,
        jwt_public_key_pem: realm.jwt_public_key_pem,
        replication_class: "SimpleStrategy",
        replication_factor: realm.replication_factor,
        device_registration_limit: realm.device_registration_limit,
        datastream_maximum_storage_retention: realm.datastream_maximum_storage_retention
      }
    }
  end

  def render("realm.json", %{realm: %Realm{replication_class: "NetworkTopologyStrategy"} = realm}) do
    %{
      data: %{
        realm_name: realm.realm_name,
        jwt_public_key_pem: realm.jwt_public_key_pem,
        replication_class: "NetworkTopologyStrategy",
        datacenter_replication_factors: realm.datacenter_replication_factors,
        device_registration_limit: realm.device_registration_limit,
        datastream_maximum_storage_retention: realm.datastream_maximum_storage_retention
      }
    }
  end
end

===== apps/astarte_housekeeping/lib/astarte_housekeeping/realms/queries.ex =====
#
# This file is part of Astarte.
#
# Copyright 2025 SECO Mind Srl
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

defmodule Astarte.Housekeeping.Realms.Queries do
  @moduledoc false
  import Ecto.Query

  alias Astarte.DataAccess.Consistency
  alias Astarte.DataAccess.CSystem
  alias Astarte.DataAccess.Database
  alias Astarte.DataAccess.Devices.Device
  alias Astarte.DataAccess.KvStore
  alias Astarte.DataAccess.Realms.Realm
  alias Astarte.DataAccess.Repo
  alias Astarte.Housekeeping.Config
  alias Astarte.Housekeeping.Realms.Realm, as: HKRealm

  require Logger

  @default_replication_factor {:simple_strategy, 1}

  def realm_existing?(realm_name) do
    keyspace_name = Realm.astarte_keyspace_name()

    query =
      from r in Realm,
        prefix: ^keyspace_name,
        where: r.realm_name == ^realm_name,
        select: count()

    consistency = Consistency.domain_model(:read)

    case Repo.safe_fetch_one(query, consistency: consistency) do
      {:ok, count} ->
        {:ok, count > 0}

      {:error, reason} ->
        Logger.warning("Cannot check if realm exists: #{inspect(reason)}.",
          tag: "realm_existing_error",
          realm: realm_name
        )

        {:error, reason}
    end
  end

  def fetch_keyspace_replication do
    astarte_keyspace = Realm.astarte_keyspace_name()
    consistency = Consistency.domain_model(:read)

    result =
      KvStore.fetch_value("astarte", "db_default_replication", :binary,
        error: :replication_not_found,
        consistency: consistency,
        prefix: astarte_keyspace
      )

    case result do
      {:ok, binary_replication} ->
        try do
          {:ok, :erlang.binary_to_term(binary_replication)}
        rescue
          _ ->
            Logger.error("Failed to deserialize replication data from KvStore",
              tag: "corrupted_replication_data"
            )

            {:error, :corrupted_replication_data}
        end

      {:error, reason} ->
        {:error, reason}
    end
  end

  def get_realm(realm_name) do
    keyspace_name = Realm.keyspace_name(realm_name)
    do_get_realm(realm_name, keyspace_name)
  rescue
    ArgumentError ->
      _ = Logger.warning("Invalid realm name.", tag: "invalid_realm_name", realm: realm_name)

      {:error, :realm_not_allowed}
  end

  defp do_get_realm(realm_name, keyspace_name) do
    with :ok <- verify_realm_exists(realm_name),
         {:ok, public_key} <- fetch_public_key(keyspace_name),
         {:ok, replication_map} <- fetch_realm_replication(keyspace_name),
         {:ok, device_registration_limit} <- fetch_device_registration_limit(realm_name) do
      max_retention = get_datastream_maximum_storage_retention(keyspace_name)

      {replication_class, replication_factor, datacenter_replication_factors} =
        replication_values(replication_map)

      realm = %HKRealm{
        realm_name: realm_name,
        jwt_public_key_pem: public_key,
        replication_class: replication_class,
        replication_factor: replication_factor,
        datacenter_replication_factors: datacenter_replication_factors,
        device_registration_limit: device_registration_limit,
        datastream_maximum_storage_retention: max_retention
      }

      {:ok, realm}
    end
  end

  def set_device_registration_limit(realm_name, device_registration_limit) do
    opts = [prefix: Realm.astarte_keyspace_name(), consistency: Consistency.domain_model(:write)]

    query =
      from Realm,
        where: [realm_name: ^realm_name]

    Repo.update_all(query, [set: [device_registration_limit: device_registration_limit]], opts)

    :ok
  end

  def delete_device_registration_limit(realm_name) do
    keyspace = Realm.astarte_keyspace_name()

    query =
      from Realm,
        where: [realm_name: ^realm_name]

    opts = [consistency: Consistency.domain_model(:write), prefix: keyspace]

    # nil inserts a tombstone
    # https://hexdocs.pm/xandra/Xandra.html#module-values
    Repo.update_all(query, [set: [device_registration_limit: nil]], opts)

    :ok
  end

  def set_datastream_maximum_storage_retention(realm_name, retention) do
    opts = [
      consistency: Consistency.domain_model(:write),
      prefix: Realm.keyspace_name(realm_name)
    ]

    kv_store = %{
      group: "realm_config",
      key: "datastream_maximum_storage_retention",
      value: retention,
      value_type: :integer
    }

    with {:error, xandra_error} <- KvStore.insert(kv_store, opts) do
      raise xandra_error
    end
  end

  def delete_datastream_maximum_storage_retention(realm_name) do
    opts = [
      consistency: Consistency.domain_model(:write),
      prefix: Realm.keyspace_name(realm_name)
    ]

    query =
      from KvStore,
        where: [group: "realm_config", key: "datastream_maximum_storage_retention"]

    Repo.delete_all(query, opts)

    :ok
  end

  def update_public_key(realm_name, jwt_public_key_pem) do
    opts = [
      consistency: Consistency.domain_model(:write),
      prefix: Realm.keyspace_name(realm_name)
    ]

    kv_store = %{
      group: "auth",
      key: "jwt_public_key_pem",
      value: jwt_public_key_pem,
      value_type: :string
    }

    with {:error, xandra_error} <- KvStore.insert(kv_store, opts) do
      raise xandra_error
    end
  end

  def create_realm(
        realm_name,
        public_key_pem,
        nil = _replication_factor,
        device_limit,
        max_retention,
        opts
      ) do
    create_realm(
      realm_name,
      public_key_pem,
      @default_replication_factor,
      device_limit,
      max_retention,
      opts
    )
  end

  def create_realm(realm_name, public_key_pem, replication, device_limit, max_retention, opts) do
    keyspace_name = Realm.keyspace_name(realm_name)

    with {:ok, replication_map_str} <- build_replication_map_str(replication) do
      execute_realm_creation(
        realm_name,
        keyspace_name,
        public_key_pem,
        replication_map_str,
        device_limit,
        max_retention,
        opts
      )
    end
  end

  defp execute_realm_creation(
         realm_name,
         keyspace_name,
         public_key_pem,
         replication_map_str,
         device_limit,
         max_retention,
         opts
       ) do
    if opts[:async] do
      {:ok, _pid} =
        Task.start(fn ->
          do_execute_realm_creation(
            realm_name,
            keyspace_name,
            public_key_pem,
            replication_map_str,
            device_limit,
            max_retention
          )
        end)

      {:ok, :started}
    else
      do_execute_realm_creation(
        realm_name,
        keyspace_name,
        public_key_pem,
        replication_map_str,
        device_limit,
        max_retention
      )
    end
  end

  defp do_execute_realm_creation(
         realm_name,
         keyspace_name,
         public_key_pem,
         replication_map_str,
         device_limit,
         max_retention
       ) do
    Repo.checkout(fn ->
      do_create_realm(
        realm_name,
        keyspace_name,
        public_key_pem,
        replication_map_str,
        device_limit,
        max_retention
      )
    end)
  end

  defp do_create_realm(
         realm_name,
         keyspace_name,
         public_key_pem,
         replication_map_str,
         device_limit,
         max_retention
       ) do
    with :ok <- create_realm_keyspace(keyspace_name, replication_map_str),
         :ok <- Database.migrate_realm(realm_name),
         :ok <- insert_realm_public_key(keyspace_name, public_key_pem),
         :ok <- insert_realm(realm_name, device_limit),
         :ok <- insert_datastream_max_retention(keyspace_name, max_retention) do
      :ok
    else
      {:error, %Xandra.Error{} = err} ->
        _ = Logger.warning("Database error: #{inspect(err)}.", tag: "database_error")
        {:error, :database_error}

      {:error, %Xandra.ConnectionError{} = err} ->
        _ =
          Logger.warning("Database connection error: #{inspect(err)}.",
            tag: "database_connection_error"
          )

        {:error, :database_connection_error}

      {:error, reason} ->
        _ =
          Logger.warning("Cannot create realm: #{inspect(reason)}.",
            tag: "realm_creation_failed",
            realm: realm_name
          )

        {:error, reason}
    end
  end

  defp check_replication(1) do
    :ok
  end

  # If replication factor is an integer, we're using SimpleStrategy
  # Check that the replication factor is <= the number of nodes in the same datacenter
  defp check_replication(replication_factor)
       when is_integer(replication_factor) and replication_factor > 1 do
    with {:ok, local_datacenter} <- get_local_datacenter() do
      check_replication_for_datacenter(local_datacenter, replication_factor, local: true)
    end
  end

  defp check_replication(datacenter_replication_factors)
       when is_map(datacenter_replication_factors) do
    with {:ok, local_datacenter} <- get_local_datacenter() do
      Enum.reduce_while(
        datacenter_replication_factors,
        :ok,
        &check_datacenter_replication(&1, &2, local_datacenter)
      )
    end
  end

  defp check_datacenter_replication({datacenter, replication_factor}, _acc, local_datacenter) do
    opts =
      if datacenter == local_datacenter do
        [local: true]
      else
        []
      end

    case check_replication_for_datacenter(datacenter, replication_factor, opts) do
      :ok -> {:cont, :ok}
      {:error, reason} -> {:halt, {:error, reason}}
    end
  end

  defp build_replication_map_str({:simple_strategy, replication_factor}) do
    with :ok <- check_replication(replication_factor) do
      replication_map_str =
        "{'class': 'SimpleStrategy', 'replication_factor': #{replication_factor}}"

      {:ok, replication_map_str}
    end
  end

  defp build_replication_map_str({:network_topology_strategy, datacenter_replication_factors}) do
    with :ok <- check_replication(datacenter_replication_factors) do
      datacenter_replications_str =
        Enum.map_join(datacenter_replication_factors, ",", fn {datacenter, replication_factor} ->
          "'#{datacenter}': #{replication_factor}"
        end)

      replication_map_str = "{'class': 'NetworkTopologyStrategy', #{datacenter_replications_str}}"

      {:ok, replication_map_str}
    end
  end

  defp build_replication_map_str(_invalid_replication) do
    {:error, :invalid_replication}
  end

  defp create_realm_keyspace(keyspace_name, replication_map_str) do
    base_query = """
    CREATE KEYSPACE #{keyspace_name}
    WITH replication = #{replication_map_str}
    AND durable_writes = true
    """

    tablets_query = base_query <> " AND tablets = { 'enabled': false }"

    case CSystem.execute_schema_change(tablets_query) do
      {:ok, %{rows: nil, num_rows: 1}} ->
        :ok

      {:ok, res} ->
        "Unexpected ok result from database while creating astarte keyspace: #{inspect(res)}"
        |> Logger.warning()

        {:error, :astarte_keyspace_creation_failed}

      {:error, %Xandra.Error{reason: :already_exists}} ->
        "Tried to create already existing realm"
        |> Logger.warning(realm: keyspace_name)

        {:error, :conflicting_realm_name}

      {:error, _reason} ->
        case CSystem.execute_schema_change(base_query) do
          {:ok, %{rows: nil, num_rows: 1}} ->
            :ok

          {:ok, res} ->
            "Unexpected ok result from database while creating astarte keyspace: #{inspect(res)}"
            |> Logger.warning()

            {:error, :astarte_keyspace_creation_failed}

          {:error, %Xandra.Error{reason: :already_exists}} ->
            "Tried to create already existing realm"
            |> Logger.warning(realm: keyspace_name)

            {:error, :conflicting_realm_name}

          {:error, reason} ->
            _ =
              Logger.warning("Cannot create keyspace: #{inspect(reason)}.",
                tag: "build_keyspace_error",
                realm: keyspace_name
              )

            {:error, reason}
        end
    end
  end

  # ScyllaDB considers TTL=0 as unset, see
  # https://opensource.docs.scylladb.com/stable/cql/time-to-live.html#notes
  defp insert_datastream_max_retention(_keyspace_name, 0) do
    :ok
  end

  # Apparently, before when the field was nil, it was encoded as zero (not optional on protobuff),
  # so treat it the same as zero
  defp insert_datastream_max_retention(_keyspace_name, nil) do
    :ok
  end

  defp insert_datastream_max_retention(keyspace_name, max_retention) do
    consistency = Consistency.domain_model(:write)

    opts = [
      consistency: consistency,
      prefix: keyspace_name
    ]

    kv_store_map = %{
      group: "realm_config",
      key: "datastream_maximum_storage_retention",
      value: max_retention,
      value_type: :integer
    }

    KvStore.insert(kv_store_map, opts)
  end

  defp insert_realm(realm_name, device_limit) do
    keyspace_name = Realm.astarte_keyspace_name()

    device_registration_limit = if device_limit == 0, do: nil, else: device_limit

    realm_attrs = %Realm{
      realm_name: realm_name,
      device_registration_limit: device_registration_limit
    }

    consistency = Consistency.domain_model(:write)

    opts = [
      consistency: consistency,
      prefix: keyspace_name
    ]

    case Repo.insert(realm_attrs, opts) do
      {:ok, _realm} -> :ok
      {:error, reason} -> {:error, reason}
    end
  end

  defp insert_realm_public_key(keyspace_name, public_key_pem) do
    consistency = Consistency.domain_model(:write)

    opts = [
      consistency: consistency,
      prefix: keyspace_name
    ]

    kv_store_map = %{
      group: "auth",
      key: "jwt_public_key_pem",
      value: public_key_pem,
      value_type: :string
    }

    KvStore.insert(kv_store_map, opts)
  end

  @doc false
  def get_local_datacenter do
    query =
      from sl in "system.local",
        select: sl.data_center

    opts = [consistency: Consistency.domain_model(:read)]

    case Repo.safe_fetch_one(query, opts) do
      {:ok, datacenter} ->
        {:ok, datacenter}

      {:error, :not_found} ->
        Logger.error(
          "Empty dataset while getting local datacenter, something is really wrong.",
          tag: "get_local_datacenter_error"
        )

        {:error, :local_datacenter_not_found}

      {:error, error} ->
        Logger.error("Error while getting local datacenter: #{inspect(error)}.",
          tag: "get_local_datacenter_error"
        )

        {:error, error}
    end
  end

  defp check_replication_for_datacenter(datacenter, replication_factor, opts) do
    query =
      from sp in "system.peers",
        hints: ["ALLOW FILTERING"],
        where: sp.data_center == ^datacenter,
        select: count()

    case Repo.safe_fetch_one(query, opts) do
      {:ok, dc_node_count} ->
        # If we're querying the datacenter of the local node, add 1 (itself) to the count
        actual_node_count =
          if opts[:local] do
            dc_node_count + 1
          else
            dc_node_count
          end

        if replication_factor <= actual_node_count do
          :ok
        else
          Logger.warning(
            "Trying to set replication_factor #{replication_factor} " <>
              "in datacenter #{datacenter} that has #{actual_node_count} nodes.",
            tag: "invalid_replication_factor",
            datacenter: datacenter,
            replication_factor: replication_factor
          )

          error_message =
            "replication_factor #{replication_factor} is >= #{actual_node_count} nodes " <>
              "in datacenter #{datacenter}"

          {:error, {:invalid_replication, error_message}}
        end

      {:error, :not_found} ->
        Logger.warning("Cannot retrieve node count for datacenter #{datacenter}.",
          tag: "datacenter_not_found",
          datacenter: datacenter
        )

        {:error, :datacenter_not_found}

      {:error, reason} ->
        Logger.warning(
          "Database error while checking replication for datacenter #{datacenter}: #{inspect(reason)}.",
          tag: "database_error",
          datacenter: datacenter
        )

        {:error, reason}
    end
  end

  defp replication_values(
         %{"class" => "org.apache.cassandra.locator.NetworkTopologyStrategy"} = replication_map
       ) do
    datacenter_map =
      replication_map
      |> Map.delete("class")
      |> Map.new(fn {datacenter, replication_factor_string} ->
        {replication_factor, ""} = Integer.parse(replication_factor_string)
        {datacenter, replication_factor}
      end)

    {"NetworkTopologyStrategy", nil, datacenter_map}
  end

  defp replication_values(
         %{
           "class" => "org.apache.cassandra.locator.SimpleStrategy",
           "replication_factor" => replication_factor_string
         } = _replication_map
       ) do
    {replication_factor, ""} = Integer.parse(replication_factor_string)

    {"SimpleStrategy", replication_factor, nil}
  end

  defp verify_realm_exists(realm_name) do
    case realm_existing?(realm_name) do
      {:ok, true} -> :ok
      {:ok, false} -> {:error, :realm_not_found}
      {:error, reason} -> {:error, reason}
    end
  end

  defp fetch_public_key(keyspace) do
    consistency = Consistency.domain_model(:read)

    KvStore.fetch_value("auth", "jwt_public_key_pem", :string,
      error: :public_key_not_found,
      consistency: consistency,
      prefix: keyspace
    )
  end

  defp fetch_realm_replication(keyspace) do
    opts = [consistency: Consistency.domain_model(:read), error: :realm_replication_not_found]

    query =
      from k in "keyspaces",
        prefix: "system_schema",
        select: k.replication,
        limit: 1

    Repo.fetch_by(query, %{keyspace_name: keyspace}, opts)
  end

  defp fetch_device_registration_limit(realm_name) do
    astarte_keyspace = Realm.astarte_keyspace_name()

    opts = [
      consistency: Consistency.domain_model(:read),
      error: :realm_device_registration_limit_not_found
    ]

    query =
      from r in Realm,
        prefix: ^astarte_keyspace,
        select: r.device_registration_limit

    Repo.fetch(query, realm_name, opts)
  end

  defp get_datastream_maximum_storage_retention(keyspace) do
    opts = [consistency: Consistency.domain_model(:read), prefix: keyspace]

    case KvStore.fetch_value(
           "realm_config",
           "datastream_maximum_storage_retention",
           :integer,
           opts
         ) do
      {:ok, value} -> value
      {:error, :not_found} -> nil
    end
  end

  def list_realms do
    keyspace_name = Realm.astarte_keyspace_name()

    query =
      from r in Realm,
        prefix: ^keyspace_name,
        select: r.realm_name

    consistency = Consistency.domain_model(:read)

    case Repo.fetch_all(query, consistency: consistency) do
      {:ok, realm_names} ->
        realms =
          Enum.map(realm_names, fn realm_name ->
            %HKRealm{realm_name: realm_name}
          end)

        {:ok, realms}

      {:error, reason} ->
        Logger.warning("Failed to list realms: #{inspect(reason)}.",
          tag: "list_realms_error"
        )

        {:error, reason}
    end
  end

  def delete_realm(realm_name, opts \\ []) do
    if Config.enable_realm_deletion!() do
      Logger.info("Deleting realm", tag: "delete_realm", realm: realm_name)
      do_delete_realm_with_options(realm_name, opts)
    else
      Logger.info("HOUSEKEEPING_ENABLE_REALM_DELETION is disabled, realm will not be deleted.",
        tag: "realm_deletion_disabled",
        realm: realm_name
      )

      {:error, :realm_deletion_disabled}
    end
  end

  defp do_delete_realm_with_options(realm_name, opts) do
    keyspace_name = Realm.keyspace_name(realm_name)

    if opts[:async] do
      {:ok, _pid} = Task.start(fn -> do_delete_realm(realm_name, keyspace_name) end)
      :ok
    else
      do_delete_realm(realm_name, keyspace_name)
    end
  end

  defp do_delete_realm(realm_name, keyspace_name) do
    with :ok <- verify_realm_deletion_preconditions(keyspace_name),
         :ok <- execute_realm_deletion(realm_name, keyspace_name) do
      :ok
    else
      {:error, %Xandra.Error{} = err} ->
        Logger.warning("Database error: #{inspect(err)}.", tag: "database_error")

        {:error, :database_error}

      {:error, %Xandra.ConnectionError{} = err} ->
        Logger.warning("Database connection error: #{inspect(err)}.",
          tag: "database_connection_error"
        )

        {:error, :database_connection_error}

      {:error, reason} ->
        Logger.warning("Cannot delete realm: #{inspect(reason)}.",
          tag: "realm_deletion_failed",
          realm: realm_name
        )

        {:error, reason}
    end
  end

  defp verify_realm_deletion_preconditions(keyspace_name) do
    case check_no_connected_devices(keyspace_name) do
      :ok ->
        :ok

      {:error, reason} ->
        Logger.warning("Realm deletion preconditions are not satisfied: #{inspect(reason)}.",
          tag: "realm_deletion_preconditions_rejected",
          realm: keyspace_name
        )

        {:error, reason}
    end
  end

  defp execute_realm_deletion(realm_name, keyspace_name) do
    with :ok <- ensure_realm_keyspace_deletion(keyspace_name),
         :ok <- remove_realm(realm_name) do
      :ok
    else
      {:error, reason} ->
        Logger.warning("Cannot delete realm: #{inspect(reason)}.",
          tag: "realm_deletion_failed",
          realm: realm_name
        )

        {:error, reason}
    end
  end

  defp check_no_connected_devices(keyspace_name) do
    query =
      from d in Device,
        hints: ["ALLOW FILTERING"],
        prefix: ^keyspace_name,
        where: d.connected == true,
        limit: 1

    consistency = Consistency.device_info(:read)

    case Repo.safe_fetch_one(query, consistency: consistency) do
      {:error, :not_found} ->
        :ok

      {:error, :realm_not_found} ->
        :ok

      _ ->
        Logger.warning("Realm #{keyspace_name} still has connected devices.",
          tag: "connected_devices_present"
        )

        {:error, :connected_devices_present}
    end
  end

  defp ensure_realm_keyspace_deletion(keyspace_name) do
    query = """
    DROP KEYSPACE IF EXISTS #{keyspace_name}
    """

    with {:ok, _} <- CSystem.execute_schema_change(query) do
      :ok
    end
  end

  defp remove_realm(realm_name) do
    keyspace_name = Realm.astarte_keyspace_name()

    query =
      from r in Realm,
        prefix: ^keyspace_name,
        where: r.realm_name == ^realm_name

    consistency = Consistency.domain_model(:write)

    Repo.delete_all(query, consistency: consistency)

    :ok
  end

  def initialize_database(default_replication \\ {:simple_strategy, 1}) do
    Logger.info("Starting Astarte keyspace initialization")

    with :ok <- create_astarte_keyspace(default_replication),
         :ok <- Database.migrate_astarte() do
      :ok
    else
      {:error, %Xandra.Error{} = err} ->
        Logger.error(
          "Database error while initializing database: #{inspect(err)}. ASTARTE WILL NOT WORK.",
          tag: "init_database_error"
        )

        {:error, :database_error}

      {:error, %Xandra.ConnectionError{} = err} ->
        Logger.error(
          "Database connection error while initializing database: #{inspect(err)}. ASTARTE WILL NOT WORK.",
          tag: "init_database_connection_error"
        )

        {:error, :database_connection_error}

      {:error, reason} ->
        Logger.error(
          "Error while initializing database: #{inspect(reason)}. ASTARTE WILL NOT WORK.",
          tag: "init_error"
        )

        {:error, reason}
    end
  end

  defp keyspace_replication_map({:simple_strategy, replication_factor}) do
    %{
      strategy: :simple,
      factor: replication_factor
    }
  end

  defp keyspace_replication_map({:network_topology_strategy, network_topology}) do
    %{
      strategy: :network_topology,
      dc_factors: network_topology
    }
  end

  def save_keyspace_replication(replication) do
    opts = [
      consistency: Consistency.domain_model(:write),
      prefix: Realm.astarte_keyspace_name()
    ]

    replication = replication |> keyspace_replication_map() |> :erlang.term_to_binary()

    kv_store = %{
      group: "astarte",
      key: "db_default_replication",
      value: replication,
      value_type: :binary
    }

    KvStore.insert(kv_store, opts)
  end

  defp create_astarte_keyspace(default_replication) do
    keyspace = Realm.astarte_keyspace_name()
    consistency = Consistency.domain_model(:write)
    replication = astarte_keyspace_replication(default_replication)

    with {:ok, replication_map_str} <- build_replication_map_str(replication),
         :ok <- do_create_astarte_keyspace(keyspace, replication_map_str, consistency) do
      :ok
    else
      {:error, reason} ->
        Logger.warning("Cannot create Astarte Keyspace: #{inspect(reason)}.",
          tag: "astarte_keyspace_creation_failed"
        )

        {:error, reason}
    end
  end

  defp astarte_keyspace_replication(default_replication) do
    case Config.astarte_keyspace_replication_strategy!() do
      :simple_strategy ->
        {:simple_strategy, Config.astarte_keyspace_replication_factor!()}

      :network_topology_strategy ->
        {:network_topology_strategy, Config.astarte_keyspace_network_replication_map!()}

      nil ->
        # No explicit replication has been configured, so the replication
        # map is derived from the current ScyllaDB network topology
        default_replication
    end
  end

  def fetch_network_topology do
    with {:ok, local_datacenter} <- get_local_datacenter(),
         {:ok, peer_datacenters} <- fetch_peer_datacenters() do
      peer_counts = Enum.frequencies(peer_datacenters)

      # `system.peers` does not include the local node, so 1 is added to the local
      # datacenter count to account for it (same as in
      # `check_replication_for_datacenter/3`).
      local_dc_node_count = Map.get(peer_counts, local_datacenter, 0) + 1
      topology = Map.put(peer_counts, local_datacenter, local_dc_node_count)

      {:ok, topology}
    end
  end

  defp fetch_peer_datacenters do
    query =
      from sp in "system.peers",
        select: sp.data_center

    opts = [consistency: Consistency.domain_model(:read)]

    case Repo.fetch_all(query, opts) do
      {:ok, datacenters} ->
        {:ok, datacenters}

      {:error, reason} ->
        {:error, reason}
    end
  end

  defp do_create_astarte_keyspace(keyspace, replication, consistency) do
    base_query = """
    CREATE KEYSPACE #{keyspace}
    WITH replication = #{replication}
    AND durable_writes = true
    """

    tablets_query = base_query <> " AND tablets = { 'enabled': false }"

    opts = [consistency: consistency]

    case Repo.query(tablets_query, [], opts) do
      {:ok, _} ->
        Logger.info("Astarte keyspace initialized")
        :ok

      {:error, %Xandra.Error{reason: :already_exists}} ->
        :ok

      _ ->
        case Repo.query(base_query, [], opts) do
          {:ok, _} ->
            Logger.info("Astarte keyspace initialized")
            :ok

          {:error, %Xandra.Error{reason: :already_exists}} ->
            :ok

          error ->
            error
        end
    end
  end
end

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/master/apps/astarte_housekeeping/lib/astarte_housekeeping/realms/realms.ex" | rg -n "update_realm|datastream|retention|:unset|device_registration|def "
41:  def list_realms do
54:  def get_realm(realm_name) do
70:  def create_realm(attrs \\ %{}, opts \\ []) do
71:    attrs = maybe_inject_default_retention(attrs)
85:  defp maybe_inject_default_retention(attrs) do
86:    case Config.default_datastream_maximum_storage_retention() do
90:      {:ok, default_retention} ->
91:        attrs |> Map.put_new("datastream_maximum_storage_retention", default_retention)
107:  @spec update_realm(binary(), map()) :: {:ok, Realm.t()} | {:error, any()}
108:  def update_realm(realm_name, attrs) do
112:         {:ok, realm} <- Core.update_realm(realm_name, changeset.changes) do
130:  def delete_realm(realm_name, opts \\ []) do

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/master/apps/astarte_housekeeping/lib/astarte_housekeeping/realms/realms.ex"
#
# This file is part of Astarte.
#
# Copyright 2017 - 2025  SECO Mind Srl
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

defmodule Astarte.Housekeeping.Realms do
  @moduledoc """
  The boundary for the Realms system.
  """

  alias Astarte.Housekeeping.Config
  alias Astarte.Housekeeping.Realms.Core
  alias Astarte.Housekeeping.Realms.Queries
  alias Astarte.Housekeeping.Realms.Realm
  alias Ecto.Changeset

  require Logger

  @doc """
  Returns the list of realms.

  ## Examples

      iex> list_realms()
      {:ok, [%Realm{}, ...]}

  """
  def list_realms do
    Queries.list_realms()
  end

  @doc """
  Gets a single realm.

  ## Examples

      iex> get_realm!(123)
      %Realm{}

  """
  def get_realm(realm_name) do
    Queries.get_realm(realm_name)
  end

  @doc """
  Creates a realm.

  ## Examples

      iex> create_realm(%{field: value})
      {:ok, %Realm{}}

      iex> create_realm(%{field: bad_value})
      {:error, ...}

  """
  def create_realm(attrs \\ %{}, opts \\ []) do
    attrs = maybe_inject_default_retention(attrs)
    changeset = Realm.changeset(%Realm{}, attrs)

    with {:ok, %Realm{} = realm} <-
           Ecto.Changeset.apply_action(changeset, :insert),
         {:ok, realm} <- resolve_replication(realm) do
      case Core.create_realm(realm, opts) do
        :ok -> {:ok, realm}
        {:ok, :started} -> {:ok, realm}
        {:error, reason} -> {:error, reason}
      end
    end
  end

  defp maybe_inject_default_retention(attrs) do
    case Config.default_datastream_maximum_storage_retention() do
      {:ok, nil} ->
        attrs

      {:ok, default_retention} ->
        attrs |> Map.put_new("datastream_maximum_storage_retention", default_retention)
    end
  end

  defp resolve_replication(%Realm{replication_class: nil} = realm) do
    with {:ok, replication} <- Queries.fetch_keyspace_replication() do
      {:ok, Realm.apply_replication(realm, replication)}
    end
  end

  defp resolve_replication(realm), do: {:ok, realm}

  @doc """
  Updates a realm with the provided list of attributes.
  Returns either {:ok, %Realm{}} or {:error, error}
  """
  @spec update_realm(binary(), map()) :: {:ok, Realm.t()} | {:error, any()}
  def update_realm(realm_name, attrs) do
    changeset = Realm.update_changeset(%Realm{realm_name: realm_name}, attrs)

    with {:ok, _realm} <- Changeset.apply_action(changeset, :update),
         {:ok, realm} <- Core.update_realm(realm_name, changeset.changes) do
      Logger.info("Successful update of realm #{realm_name}", tag: "realm_update_success")
      {:ok, realm}
    end
  end

  @doc """
  Deletes a Realm.

  ## Examples

      iex> delete_realm(realm_name)
      :ok

      iex> delete_realm(realm_name)
      {:error, ...}

  """
  def delete_realm(realm_name, opts \\ []) do
    Queries.delete_realm(realm_name, opts)
  end
end

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/master/apps/astarte_housekeeping/lib/astarte_housekeeping/realms/core.ex"
#
# This file is part of Astarte.
#
# Copyright 2025 SECO Mind Srl
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

defmodule Astarte.Housekeeping.Realms.Core do
  @moduledoc false
  alias Astarte.Events.AMQP.Vhost
  alias Astarte.Housekeeping.Realms.Queries
  alias Astarte.Housekeeping.Realms.Realm
  alias Astarte.Secrets

  require Logger

  def update_realm(realm_name, changes) do
    Logger.info("Updating realm #{realm_name}", tag: "realm_update_start")

    with {:ok, realm} <- Queries.get_realm(realm_name) do
      updated_realm = apply_realm_changes(realm, changes)
      {:ok, updated_realm}
    end
  end

  defp apply_realm_changes(realm, changes) do
    Enum.reduce(changes, realm, fn {field, changed_value}, acc_realm ->
      apply_field_change(acc_realm, field, changed_value)
    end)
  end

  defp apply_field_change(realm, field, changed_value) do
    case update_realm_field(realm, field, changed_value) do
      :ok -> %{realm | field => changed_value}
      {:ok, updated_realm} -> updated_realm
    end
  end

  defp update_realm_field(realm, :jwt_public_key_pem, jwt_public_key_pem) do
    Queries.update_public_key(realm.realm_name, jwt_public_key_pem)
  end

  defp update_realm_field(realm, :device_registration_limit, :unset) do
    Logger.info("Removing device registration limit", tag: "remove_device_registration_limit")
    Queries.delete_device_registration_limit(realm.realm_name)

    {:ok, %{realm | device_registration_limit: nil}}
  end

  defp update_realm_field(realm, :device_registration_limit, limit) when is_integer(limit) do
    Logger.info("Updating device registration limit", tag: "update_device_registration_limit")

    Queries.set_device_registration_limit(realm.realm_name, limit)
  end

  defp update_realm_field(realm, :datastream_maximum_storage_retention, :unset) do
    Logger.info("Removing datastream maximum storage retention",
      tag: "remove_datastream_maximum_storage_retention"
    )

    Queries.delete_datastream_maximum_storage_retention(realm.realm_name)

    {:ok, %{realm | datastream_maximum_storage_retention: nil}}
  end

  defp update_realm_field(realm, :datastream_maximum_storage_retention, retention)
       when is_integer(retention) do
    Logger.info("Updating datastream maximum storage retention",
      tag: "update_datastream_maximum_storage_retention"
    )

    Queries.set_datastream_maximum_storage_retention(realm.realm_name, retention)
  end

  def create_realm(
        %Realm{
          realm_name: realm_name,
          jwt_public_key_pem: pem,
          replication_class: "SimpleStrategy",
          replication_factor: replication_factor,
          device_registration_limit: device_registration_limit,
          datastream_maximum_storage_retention: datastream_maximum_storage_retention
        },
        opts
      ) do
    do_create_realm(
      realm_name,
      pem,
      {:simple_strategy, replication_factor},
      device_registration_limit,
      datastream_maximum_storage_retention,
      opts
    )
  end

  def create_realm(
        %Realm{
          realm_name: realm_name,
          jwt_public_key_pem: pem,
          replication_class: "NetworkTopologyStrategy",
          datacenter_replication_factors: replication_factors_map,
          device_registration_limit: device_registration_limit,
          datastream_maximum_storage_retention: datastream_maximum_storage_retention
        },
        opts
      ) do
    datacenter_replication_factors_map = Enum.into(Enum.to_list(replication_factors_map), %{})

    do_create_realm(
      realm_name,
      pem,
      {:network_topology_strategy, datacenter_replication_factors_map},
      device_registration_limit,
      datastream_maximum_storage_retention,
      opts
    )
  end

  defp do_create_realm(
         realm_name,
         pem,
         replication,
         device_registration_limit,
         datastream_maximum_storage_retention,
         opts
       ) do
    with :ok <- create_vhost(realm_name),
         :ok <- ensure_realm_kek(realm_name) do
      Queries.create_realm(
        realm_name,
        pem,
        replication,
        device_registration_limit,
        datastream_maximum_storage_retention,
        opts
      )
    end
  end

  defp ensure_realm_kek(realm_name) do
    case Secrets.create_realm_kek(realm_name) do
      {:ok, _key} ->
        :ok

      :error ->
        Logger.error("Failed to create realm KEK for realm #{realm_name}",
          tag: "realm_kek_creation_failed"
        )

        {:error, :realm_kek_creation_failed}
    end
  end

  defp create_vhost(realm_name) do
    case Vhost.create_vhost(realm_name) do
      :ok -> :ok
      :error -> {:error, :vhost_creation_error}
    end
  end
end

Let me check the upstream code at the pinned v1.2.0 tag specifically, since that's the compatibility target.
$ curl -s --max-time 15 "https://api.github.com/repos/astarte-platform/astarte/tags?per_page=100" | rg -o '"name": "[^"]*"' | head -30
"name": "v1.4.0-rc.5"
"name": "v1.4.0-rc.4"
"name": "v1.4.0-rc.3"
"name": "v1.4.0-rc.2"
"name": "v1.4.0-rc.1"
"name": "v1.4.0-rc.0"
"name": "v1.3.3"
"name": "v1.3.2"
"name": "v1.3.1"
"name": "v1.3.0"
"name": "v1.3.0-rc.2"
"name": "v1.3.0-rc.1"
"name": "v1.3.0-rc.0"
"name": "v1.2.2"
"name": "v1.2.2-rc.0"
"name": "v1.2.1"
"name": "v1.2.1-rc.1"
"name": "v1.2.1-rc.0"
"name": "v1.2.1-alpha.0"
"name": "v1.2.0"
"name": "v1.2.0-rc.0"
"name": "v1.1.1"
"name": "v1.1.0"
"name": "v1.1.0-rc.0"
"name": "v1.1.0-alpha.0"
"name": "v1.0.6"
"name": "v1.0.5"
"name": "v1.0.4"
"name": "v1.0.3"
"name": "v1.0.2"

$ for f in "apps/astarte_housekeeping/lib/astarte_housekeeping/realms/realm.ex" "apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex" "apps/astarte_housekeeping/lib/astarte_housekeeping/realms/core.ex"; do echo "===== $f @v1.2.0 ====="; curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/v1.2.0/$f" | rg -n -i "retention|register|cast|update_changeset|:unset|device_limit" ; echo; done
===== apps/astarte_housekeeping/lib/astarte_housekeeping/realms/realm.ex @v1.2.0 =====

===== apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex @v1.2.0 =====

===== apps/astarte_housekeeping/lib/astarte_housekeeping/realms/core.ex @v1.2.0 =====

$ curl -s --max-time 15 "https://api.github.com/repos/astarte-platform/astarte/git/trees/v1.2.0?recursive=1" | rg -o '"path": "[^"]*"' | rg -i "housekeeping" | rg -iv "test|priv/migration|release|Dockerfile|\.tpl|\.sh|mix|config" | head -40
"path": "apps/astarte_housekeeping"
"path": "apps/astarte_housekeeping/.dockerignore"
"path": "apps/astarte_housekeeping/.formatter.exs"
"path": "apps/astarte_housekeeping/README.md"
"path": "apps/astarte_housekeeping/coveralls.json"
"path": "apps/astarte_housekeeping/lib"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping/backend_supervisor.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping/c_system.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping/engine.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping/migrator.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping/queries.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping/rpc"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping/rpc/handler.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web/plug"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web/plug/health_plug.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web/plug/metrics_plug.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web/router.ex"
"path": "apps/astarte_housekeeping/lib/astarte_housekeeping_web/telemetry.ex"
"path": "apps/astarte_housekeeping/priv"
"path": "apps/astarte_housekeeping/rel"
"path": "apps/astarte_housekeeping_api"
"path": "apps/astarte_housekeeping_api/.dockerignore"
"path": "apps/astarte_housekeeping_api/.formatter.exs"
"path": "apps/astarte_housekeeping_api/README.md"
"path": "apps/astarte_housekeeping_api/lib"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/application.ex"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/auth"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/auth/auth.ex"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/auth/user.ex"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/health"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/health/backend_health.ex"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/health/health.ex"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/realms"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/realms/non_negative_integer_or_unset_type.ex"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/realms/realm.ex"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/realms/realms.ex"

$ for f in "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/realms/realm.ex" "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/realms/non_negative_integer_or_unset_type.ex" "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/realms/realms.ex"; do echo "===== $f @v1.2.0 ====="; curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/v1.2.0/$f" | rg -n -i "retention|register|cast|update_changeset|unset|:unset|device_registration|minimum|allowed"; echo; done
===== apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/realms/realm.ex @v1.2.0 =====
23:  alias Astarte.Housekeeping.API.Realms.NonNegativeIntegerOrUnsetType
29:  @allowed_fields [
33:    :device_registration_limit,
34:    :datastream_maximum_storage_retention
47:    field :device_registration_limit, NonNegativeIntegerOrUnsetType
48:    field :datastream_maximum_storage_retention, NonNegativeIntegerOrUnsetType
53:    |> cast(params, @allowed_fields)
67:  def update_changeset(realm, params \\ %{}) do
69:    |> cast(params, @allowed_fields)
85:      |> cast(params, @required_fields)

===== apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/realms/non_negative_integer_or_unset_type.ex @v1.2.0 =====
19:defmodule Astarte.Housekeeping.API.Realms.NonNegativeIntegerOrUnsetType do
23:  def cast(:unset), do: {:ok, :unset}
25:  def cast(n) when is_integer(n) do
29:  def cast(_), do: :error
32:  def load(:unset), do: {:ok, :unset}
35:  def dump(:unset), do: {:ok, :unset}

===== apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/realms/realms.ex @v1.2.0 =====
87:    changeset = %Realm{realm_name: realm_name} |> Realm.update_changeset(attrs)

$ echo "===== realms.ex update_realm @v1.2.0 ====="; curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/v1.2.0/apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/realms/realms.ex" | sed -n '70,120p'
===== realms.ex update_realm @v1.2.0 =====
    with {:ok, %Realm{} = realm} <-
           Ecto.Changeset.apply_action(changeset, :insert) do
      case Housekeeping.create_realm(realm, opts) do
        :ok -> {:ok, realm}
        {:ok, :started} -> {:ok, realm}
        {:error, reason} -> {:error, reason}
      end
    end
  end

  @doc """
  Updates a realm with the provided list of attributes.
  Returns either {:ok, %Realm{}} or {:error, error}
  where `error` is an Ecto.Changeset describing the error.
  """
  @spec update_realm(binary(), map()) :: {:ok, Realm.t()} | {:error, Ecto.Changeset.t()}
  def update_realm(realm_name, attrs) do
    changeset = %Realm{realm_name: realm_name} |> Realm.update_changeset(attrs)

    with {:ok, %Realm{} = realm_update} <-
           Ecto.Changeset.apply_action(changeset, :update) do
      Housekeeping.update_realm(realm_update)
    end
  end

  @doc """
  Deletes a Realm.

  ## Examples

      iex> delete_realm(realm_name)
      :ok

      iex> delete_realm(realm_name)
      {:error, ...}

  """
  def delete_realm(realm_name, opts \\ []) do
    case Housekeeping.delete_realm(realm_name, opts) do
      :ok -> :ok
      {:ok, :started} -> :ok
      {:error, reason} -> {:error, reason}
    end
  end

  @doc """
  Returns a datastructure for tracking realm changes.

  ## Examples

      iex> change_realm(realm)

$ echo "===== controller @v1.2.0 ====="; curl -s --max-time 15 "https://api.github.com/repos/astarte-platform/astarte/git/trees/v1.2.0?recursive=1" | rg -o '"path": "[^"]*"' | rg "housekeeping_api.*controller"; echo ====; curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/v1.2.0/apps/astarte_housekeeping_api/lib/astarte_housekeeping_api_web/controllers/realm_controller.ex" | rg -n -i "update|normalize|patch|def "
===== controller @v1.2.0 =====
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api_web/controllers"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api_web/controllers/fallback_controller.ex"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api_web/controllers/realm_controller.ex"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api_web/controllers/version_controller.ex"
"path": "apps/astarte_housekeeping_api/test/astarte_housekeeping_api_web/controllers"
"path": "apps/astarte_housekeeping_api/test/astarte_housekeeping_api_web/controllers/realm_controller_test.exs"
====
27:  def index(conn, _params) do
32:  def create(conn, %{"data" => realm_params} = params) do
49:  def show(conn, %{"id" => id}) do
55:  def update(%Plug.Conn{method: "PATCH"} = conn, %{
59:    update_params = normalize_update_attrs(realm_params)
61:    with {:ok, %Realm{} = updated_realm} <- Realms.update_realm(realm_name, update_params) do
62:      render(conn, "show.json", realm: updated_realm)
66:  def delete(conn, %{"id" => id} = params) do
79:  defp normalize_update_attrs(update_attrs) when is_map(update_attrs) do
80:    update_attrs
81:    |> Map.replace_lazy(:device_registration_limit, &normalize_integer_or_nil/1)
82:    |> Map.replace_lazy("device_registration_limit", &normalize_integer_or_nil/1)
83:    |> Map.replace_lazy(:datastream_maximum_storage_retention, &normalize_integer_or_nil/1)
84:    |> Map.replace_lazy("datastream_maximum_storage_retention", &normalize_integer_or_nil/1)
87:  defp normalize_integer_or_nil(value) when is_nil(value), do: :unset
88:  defp normalize_integer_or_nil(value), do: value

$ echo "===== backend queries.ex @v1.2.0 retention ====="; curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/v1.2.0/apps/astarte_housekeeping/lib/astarte_housekeeping/queries.ex" | rg -n -B2 -A12 -i "retention"; echo; echo "===== rpc handler @v1.2.0 ====="; curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/v1.2.0/apps/astarte_housekeeping/lib/astarte_housekeeping/rpc/handler.ex" | rg -n -B2 -A14 -i "retention|update_realm|create_realm|:unset"
===== backend queries.ex @v1.2.0 retention =====
31-        nil = _replication_factor,
32-        device_limit,
33:        max_retention,
34-        opts
35-      ) do
36-    create_realm(
37-      realm_name,
38-      public_key_pem,
39-      @default_replication_factor,
40-      device_limit,
41:      max_retention,
42-      opts
43-    )
44-  end
45-
46:  def create_realm(realm_name, public_key_pem, replication, device_limit, max_retention, opts) do
47-    with :ok <- validate_realm_name(realm_name),
48-         keyspace_name =
49-           CQLUtils.realm_name_to_keyspace_name(realm_name, Config.astarte_instance_id!()),
50-         :ok <- Xandra.Cluster.run(:xandra, &check_replication(&1, replication)),
51-         {:ok, replication_map_str} <- build_replication_map_str(replication) do
52-      if opts[:async] do
53-        {:ok, _pid} =
54-          Task.start(fn ->
55-            do_create_realm(
56-              realm_name,
57-              keyspace_name,
58-              public_key_pem,
59-              replication_map_str,
60-              device_limit,
61:              max_retention
62-            )
63-          end)
64-
65-        :ok
66-      else
67-        do_create_realm(
68-          realm_name,
69-          keyspace_name,
70-          public_key_pem,
71-          replication_map_str,
72-          device_limit,
73:          max_retention
74-        )
75-      end
76-    end
77-  end
78-
79-  def delete_realm(realm_name, opts \\ []) do
80-    keyspace_name =
81-      CQLUtils.realm_name_to_keyspace_name(realm_name, Config.astarte_instance_id!())
82-
83-    if opts[:async] do
84-      {:ok, _pid} = Task.start(fn -> do_delete_realm(realm_name, keyspace_name) end)
85-
--
210-         replication_map_str,
211-         device_limit,
212:         max_retention
213-       ) do
214-    Xandra.Cluster.run(:xandra, [timeout: 60_000], fn conn ->
215-      with :ok <- validate_realm_name(realm_name),
216-           :ok <- create_realm_keyspace(conn, keyspace_name, replication_map_str),
217-           {:ok, keyspace_conn} <- build_keyspace_conn(conn, keyspace_name),
218-           :ok <- create_realm_kv_store(keyspace_conn),
219-           :ok <- create_names_table(keyspace_conn),
220-           :ok <- create_devices_table(keyspace_conn),
221-           :ok <- create_endpoints_table(keyspace_conn),
222-           :ok <- create_interfaces_table(keyspace_conn),
223-           :ok <- create_individual_properties_table(keyspace_conn),
224-           :ok <- create_simple_triggers_table(keyspace_conn),
--
228-           :ok <- insert_realm_astarte_schema_version(keyspace_conn),
229-           :ok <- insert_realm(conn, realm_name, device_limit),
230:           :ok <- insert_datastream_max_retention(keyspace_conn, max_retention) do
231-        :ok
232-      else
233-        {:error, reason} ->
234-          _ =
235-            Logger.warning("Cannot create realm: #{inspect(reason)}.",
236-              tag: "realm_creation_failed",
237-              realm: realm_name
238-            )
239-
240-          {:error, reason}
241-      end
242-    end)
--
454-      value_type int,
455-      reliability int,
456:      retention int,
457-      expiry int,
458:      database_retention_ttl int,
459:      database_retention_policy int,
460-      allow_unset boolean,
461-      explicit_timestamp boolean,
462-      description varchar,
463-      doc varchar,
464-
465-      PRIMARY KEY ((interface_id), endpoint_id)
466-    );
467-    """
468-
469-    with {:ok, %Xandra.SchemaChange{}} <- CSystem.execute_schema_change(conn, query) do
470-      :ok
471-    else
--
770-  # ScyllaDB considers TTL=0 as unset, see
771-  # https://opensource.docs.scylladb.com/stable/cql/time-to-live.html#notes
772:  defp insert_datastream_max_retention(_conn_realm, 0) do
773-    :ok
774-  end
775-
776:  defp insert_datastream_max_retention({conn, keyspace_name}, max_retention) do
777-    statement = """
778-    INSERT INTO :keyspace_name.kv_store (group, key, value)
779:    VALUES ('realm_config', 'datastream_maximum_storage_retention', intAsBlob(:max_retention));
780-    """
781-
782-    params = %{
783:      "max_retention" => max_retention
784-    }
785-
786-    # This is safe since we checked the realm name in the caller
787-    query = String.replace(statement, ":keyspace_name", keyspace_name)
788-
789-    with {:ok, prepared} <- Xandra.prepare(conn, query),
790-         {:ok, %Xandra.Void{}} <-
791-           Xandra.execute(conn, prepared, params, consistency: :each_quorum) do
792-      :ok
793-    else
794-      {:error, %Xandra.Error{} = err} ->
795-        _ = Logger.warning("Database error: #{inspect(err)}.", tag: "database_error")
--
1070-           {:ok, replication_map} <- get_realm_replication(conn, keyspace_name),
1071-           {:ok, device_registration_limit} <- get_device_registration_limit(conn, realm_name),
1072:           {:ok, max_retention} <-
1073:             get_datastream_maximum_storage_retention(conn, keyspace_name) do
1074-        case replication_map do
1075-          %{
1076-            "class" => "org.apache.cassandra.locator.SimpleStrategy",
1077-            "replication_factor" => replication_factor_string
1078-          } ->
1079-            {replication_factor, ""} = Integer.parse(replication_factor_string)
1080-
1081-            %{
1082-              realm_name: realm_name,
1083-              jwt_public_key_pem: public_key,
1084-              replication_class: "SimpleStrategy",
1085-              replication_factor: replication_factor,
1086-              device_registration_limit: device_registration_limit,
1087:              datastream_maximum_storage_retention: max_retention
1088-            }
1089-
1090-          %{"class" => "org.apache.cassandra.locator.NetworkTopologyStrategy"} ->
1091-            datacenter_replication_factors =
1092-              Enum.reduce(replication_map, %{}, fn
1093-                {"class", _}, acc ->
1094-                  acc
1095-
1096-                {datacenter, replication_factor_string}, acc ->
1097-                  {replication_factor, ""} = Integer.parse(replication_factor_string)
1098-                  Map.put(acc, datacenter, replication_factor)
1099-              end)
--
1105-              datacenter_replication_factors: datacenter_replication_factors,
1106-              device_registration_limit: device_registration_limit,
1107:              datastream_maximum_storage_retention: max_retention
1108-            }
1109-        end
1110-      else
1111-        # Returned by is_realm_existing
1112-        {:ok, false} ->
1113-          {:error, :realm_not_found}
1114-
1115-        {:error, reason} ->
1116-          _ =
1117-            Logger.warning("Error while getting realm: #{inspect(reason)}.",
1118-              tag: "get_realm_error",
1119-              realm: realm_name
--
1125-  end
1126-
1127:  def set_datastream_maximum_storage_retention(realm_name, new_retention) do
1128-    keyspace_name =
1129-      CQLUtils.realm_name_to_keyspace_name(realm_name, Config.astarte_instance_id!())
1130-
1131-    with :ok <- validate_realm_name(realm_name) do
1132-      Xandra.Cluster.run(
1133-        :xandra,
1134:        &do_set_datastream_maximum_storage_retention(&1, keyspace_name, new_retention)
1135-      )
1136-    end
1137-  end
1138-
1139:  def delete_datastream_maximum_storage_retention(realm_name) do
1140-    keyspace_name =
1141-      CQLUtils.realm_name_to_keyspace_name(realm_name, Config.astarte_instance_id!())
1142-
1143-    with :ok <- validate_realm_name(realm_name) do
1144-      Xandra.Cluster.run(
1145-        :xandra,
1146:        &do_delete_datastream_maximum_storage_retention(&1, keyspace_name)
1147-      )
1148-    end
1149-  end
1150-
1151-  defp is_realm_existing(conn, realm_name) do
1152-    query = """
1153-    SELECT realm_name from #{CQLUtils.realm_name_to_keyspace_name("astarte", Config.astarte_instance_id!())}.realms
1154-    WHERE realm_name=:realm_name;
1155-    """
1156-
1157-    with {:ok, prepared} <- Xandra.prepare(conn, query),
1158-         {:ok, %Xandra.Page{} = page} <-
--
1278-  end
1279-
1280:  defp do_set_datastream_maximum_storage_retention(conn, realm_name, new_retention) do
1281-    statement = """
1282-    UPDATE :realm_name.kv_store
1283:    SET value = intAsBlob(:new_retention)
1284:    WHERE group='realm_config' AND key='datastream_maximum_storage_retention'
1285-    """
1286-
1287-    params = %{
1288:      "new_retention" => new_retention
1289-    }
1290-
1291-    # TODO move away from this when NoaccOS' PR is merged
1292-    query = String.replace(statement, ":realm_name", realm_name)
1293-
1294-    # TODO refactor when NoaccOS' PR is merged
1295-    with {:ok, prepared} <- Xandra.prepare(conn, query),
1296-         {:ok, result} <- Xandra.execute(conn, prepared, params, consistency: :quorum) do
1297-      {:ok, result}
1298-    else
1299-      {:error, %Xandra.Error{} = err} ->
1300-        _ = Logger.warning("Database error: #{Exception.message(err)}.", tag: "database_error")
--
1342-  end
1343-
1344:  defp do_delete_datastream_maximum_storage_retention(conn, realm_name) do
1345-    statement = """
1346-    DELETE FROM :realm_name.kv_store
1347:    WHERE group='realm_config' AND key='datastream_maximum_storage_retention'
1348-    """
1349-
1350-    # TODO move away from this when NoaccOS' PR is merged
1351-    query = String.replace(statement, ":realm_name", realm_name)
1352-
1353-    # TODO refactor when NoaccOS' PR is merged
1354-    with {:ok, prepared} <- Xandra.prepare(conn, query),
1355-         {:ok, result} <- Xandra.execute(conn, prepared, %{}, consistency: :quorum) do
1356-      {:ok, result}
1357-    else
1358-      {:error, %Xandra.Error{} = err} ->
1359-        _ = Logger.warning("Database error: #{Exception.message(err)}.", tag: "database_error")
--
1480-  end
1481-
1482:  defp get_datastream_maximum_storage_retention(conn, realm_name) do
1483-    statement = """
1484-    SELECT blobAsInt(value)
1485-    FROM :realm_name.kv_store
1486:    WHERE group='realm_config' AND key='datastream_maximum_storage_retention'
1487-    """
1488-
1489-    # TODO change this once NoaccOS' PR is merged
1490-    with :ok <- validate_realm_name(realm_name),
1491-         query = String.replace(statement, ":realm_name", realm_name),
1492-         {:ok, prepared} <- Xandra.prepare(conn, query),
1493-         {:ok, page} <- Xandra.execute(conn, prepared, %{}) do
1494-      case Enum.fetch(page, 0) do
1495-        {:ok, %{"system.blobasint(value)" => value}} ->
1496-          {:ok, value}
1497-
1498-        :error ->

===== rpc handler @v1.2.0 =====
59-
60-  # Here for retrocompatibility with old protos serialized with Exprotobuf
61:  defp call_rpc({:create_realm, %CreateRealm{realm: ""}}) do
62-    _ = Logger.warning("CreateRealm with empty realm.", tag: "rpc_create_nil_realm")
63-    generic_error(:empty_name, "empty realm name")
64-  end
65-
66-  # Here for retrocompatibility with old protos serialized with Exprotobuf
67:  defp call_rpc({:create_realm, %CreateRealm{jwt_public_key_pem: ""}}) do
68-    _ =
69-      Logger.warning("CreateRealm with empty jwt_public_key_pem.",
70-        tag: "rpc_create_nil_public_key"
71-      )
72-
73-    generic_error(:empty_public_key, "empty jwt public key pem")
74-  end
75-
76:  defp call_rpc({:create_realm, %CreateRealm{realm: nil}}) do
77-    _ = Logger.warning("CreateRealm with empty realm.", tag: "rpc_create_nil_realm")
78-    generic_error(:empty_name, "empty realm name")
79-  end
80-
81:  defp call_rpc({:create_realm, %CreateRealm{jwt_public_key_pem: nil}}) do
82-    _ =
83-      Logger.warning("CreateRealm with empty jwt_public_key_pem.",
84-        tag: "rpc_create_nil_public_key"
85-      )
86-
87-    generic_error(:empty_public_key, "empty jwt public key pem")
88-  end
89-
90-  defp call_rpc(
91:         {:create_realm,
92-          %CreateRealm{
93-            realm: realm,
94-            jwt_public_key_pem: pub_key,
95-            replication_class: :NETWORK_TOPOLOGY_STRATEGY,
96-            datacenter_replication_factors: datacenter_replication_factors,
97-            device_registration_limit: device_registration_limit,
98:            datastream_maximum_storage_retention: datastream_maximum_storage_retention,
99-            async_operation: async
100-          }}
101-       ) do
102-    with {:ok, false} <- Astarte.Housekeeping.Engine.is_realm_existing(realm),
103-         datacenter_replication_factors_map = Enum.into(datacenter_replication_factors, %{}),
104-         :ok <-
105:           Engine.create_realm(
106-             realm,
107-             pub_key,
108-             datacenter_replication_factors_map,
109-             device_registration_limit,
110:             datastream_maximum_storage_retention,
111-             async: async
112-           ) do
113-      generic_ok(async)
114-    else
115-      # This comes from is_realm_existing
116-      {:ok, true} ->
117-        _ =
118-          Logger.warning("CreateRealm with already existing realm.",
119-            tag: "rpc_create_existing_realm",
120-            realm: realm
121-          )
122-
123-        generic_error(:existing_realm, "realm already exists")
124-
--
131-  end
132-
133:  defp call_rpc({:create_realm, %CreateRealm{replication_factor: 0} = call}) do
134-    # Due to new proto3 defaults, if replication factor is not explictly set, it now defaults
135-    # to 0 instead of nil. Hide this implementation detail to the outside world.
136:    call_rpc({:create_realm, %{call | replication_factor: nil}})
137-  end
138-
139-  defp call_rpc(
140:         {:create_realm,
141-          %CreateRealm{
142-            realm: realm,
143-            jwt_public_key_pem: pub_key,
144-            replication_factor: replication_factor,
145-            device_registration_limit: device_registration_limit,
146:            datastream_maximum_storage_retention: datastream_maximum_storage_retention,
147-            async_operation: async
148-          }}
149-       ) do
150-    with {:ok, false} <- Astarte.Housekeeping.Engine.is_realm_existing(realm),
151-         :ok <-
152:           Engine.create_realm(
153-             realm,
154-             pub_key,
155-             replication_factor,
156-             device_registration_limit,
157:             datastream_maximum_storage_retention,
158-             async: async
159-           ) do
160-      generic_ok(async)
161-    else
162-      # This comes from is_realm_existing
163-      {:ok, true} ->
164-        _ =
165-          Logger.warning("CreateRealm with already existing realm.",
166-            tag: "rpc_create_existing_realm",
167-            realm: realm
168-          )
169-
170-        generic_error(:existing_realm, "realm already exists")
171-
--
178-  end
179-
180:  defp call_rpc({:update_realm, %UpdateRealm{} = call}) do
181-    %UpdateRealm{realm: realm_name, device_registration_limit: device_registration_limit} = call
182-    new_limit = extract_device_registration_limit(device_registration_limit)
183-    attrs = %{call | device_registration_limit: new_limit}
184-
185:    with {:ok, realm} <- Astarte.Housekeeping.Engine.update_realm(realm_name, attrs) do
186-      case realm do
187-        %{
188-          realm_name: realm_name_reply,
189-          jwt_public_key_pem: public_key,
190-          replication_class: "SimpleStrategy",
191-          replication_factor: replication_factor,
192-          device_registration_limit: limit,
193:          datastream_maximum_storage_retention: retention
194-        } ->
195-          %GetRealmReply{
196-            realm_name: realm_name_reply,
197-            jwt_public_key_pem: public_key,
198-            replication_class: :SIMPLE_STRATEGY,
199-            replication_factor: replication_factor,
200-            device_registration_limit: limit,
201:            datastream_maximum_storage_retention: retention
202-          }
203-          |> encode_reply(:get_realm_reply)
204-          |> ok_wrap
205-
206-        %{
207-          realm_name: realm_name_reply,
208-          jwt_public_key_pem: public_key,
209-          replication_class: "NetworkTopologyStrategy",
210-          datacenter_replication_factors: datacenter_replication_factors,
211-          device_registration_limit: limit,
212:          datastream_maximum_storage_retention: retention
213-        } ->
214-          %GetRealmReply{
215-            realm_name: realm_name_reply,
216-            jwt_public_key_pem: public_key,
217-            replication_class: :NETWORK_TOPOLOGY_STRATEGY,
218-            datacenter_replication_factors: datacenter_replication_factors,
219-            device_registration_limit: limit,
220:            datastream_maximum_storage_retention: retention
221-          }
222-          |> encode_reply(:get_realm_reply)
223-          |> ok_wrap
224-      end
225-    else
226-      {:error, reason} ->
227-        generic_error(reason)
228-    end
229-  end
230-
231-  # Here for retrocompatibility with old protos serialized with Exprotobuf
232-  defp call_rpc({:delete_realm, %DeleteRealm{realm: ""}}) do
233-    _ = Logger.warning("DeleteRealm with empty realm.", tag: "rpc_delete_empty_realm")
234-    generic_error(:empty_name, "empty realm name")
--
310-        replication_factor: replication_factor,
311-        device_registration_limit: device_registration_limit,
312:        datastream_maximum_storage_retention: datastream_maximum_storage_retention
313-      } ->
314-        %GetRealmReply{
315-          realm_name: realm_name_reply,
316-          jwt_public_key_pem: public_key,
317-          replication_class: :SIMPLE_STRATEGY,
318-          replication_factor: replication_factor,
319-          device_registration_limit: device_registration_limit,
320:          datastream_maximum_storage_retention: datastream_maximum_storage_retention
321-        }
322-        |> encode_reply(:get_realm_reply)
323-        |> ok_wrap
324-
325-      %{
326-        realm_name: realm_name_reply,
327-        jwt_public_key_pem: public_key,
328-        replication_class: "NetworkTopologyStrategy",
329-        datacenter_replication_factors: datacenter_replication_factors,
330-        device_registration_limit: device_registration_limit,
331:        datastream_maximum_storage_retention: datastream_maximum_storage_retention
332-      } ->
333-        datacenter_replication_factors_list = Enum.into(datacenter_replication_factors, [])
334-
335-        %GetRealmReply{
336-          realm_name: realm_name_reply,
337-          jwt_public_key_pem: public_key,
338-          replication_class: :NETWORK_TOPOLOGY_STRATEGY,
339-          datacenter_replication_factors: datacenter_replication_factors_list,
340-          device_registration_limit: device_registration_limit,
341:          datastream_maximum_storage_retention: datastream_maximum_storage_retention
342-        }
343-        |> encode_reply(:get_realm_reply)
344-        |> ok_wrap
345-
346-      {:error, reason} ->
347-        generic_error(reason)
348-    end
349-  end
350-
351-  defp generic_error(
352-         error_name,
353-         user_readable_message \\ nil,
354-         user_readable_error_name \\ nil,
355-         error_data \\ nil

$ echo "===== update proto @v1.2.0 ====="; curl -s --max-time 15 "https://api.github.com/repos/astarte-platform/astarte/git/trees/v1.2.0?recursive=1" | rg -o '"path": "[^"]*"' | rg -i "proto"; echo; curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/v1.2.0/proto/astarte_housekeeping.proto" | rg -n -A20 "message UpdateRealm|message CreateRealm"
===== update proto @v1.2.0 =====
"path": "doc/pages/architecture/080-mqtt-v1-protocol.md"

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/v1.2.0/apps/astarte_housekeeping/lib/astarte_housekeeping/engine.ex" | rg -n -B3 -A30 "def update_realm"
119-          | {:error, :update_public_key_fail}
120-          | {:error, :database_error}
121-          | {:error, :database_connection_error}
122:  def update_realm(realm_name, update_attrs) do
123-    if is_realm_update_valid?(update_attrs) do
124-      _ = Logger.info("Updating realm #{realm_name}", tag: "realm_update_start")
125-      do_update_realm(realm_name, update_attrs)
126-    else
127-      _ =
128-        Logger.warning("Rejecting update for realm #{realm_name}",
129-          tag: "invalid_update_parameters"
130-        )
131-
132-      {:error, :invalid_update_parameters}
133-    end
134-  end
135-
136-  defp do_update_realm(realm_name, update_attrs) do
137-    %{
138-      jwt_public_key_pem: new_jwt_public_key_pem,
139-      device_registration_limit: new_limit,
140-      datastream_maximum_storage_retention: new_retention
141-    } = update_attrs
142-
143-    with realm when is_map(realm) <- Queries.get_realm(realm_name),
144-         :ok <- update_jwt_public_key_pem(realm_name, new_jwt_public_key_pem),
145-         :ok <- update_device_registration_limit(realm_name, new_limit),
146-         :ok <- update_datastream_maximum_storage_retention(realm_name, new_retention) do
147-      updated_realm = merge_realm_status(realm, update_attrs)
148-
149-      _ = Logger.info("Successful update of realm #{realm_name}", tag: "realm_update_success")
150-
151-      {:ok, updated_realm}
152-    end

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/v1.2.0/apps/astarte_housekeeping/lib/astarte_housekeeping/engine.ex" | rg -n -B3 -A25 "is_realm_update_valid|update_datastream_maximum_storage_retention|merge_realm_status"
120-          | {:error, :database_error}
121-          | {:error, :database_connection_error}
122-  def update_realm(realm_name, update_attrs) do
123:    if is_realm_update_valid?(update_attrs) do
124-      _ = Logger.info("Updating realm #{realm_name}", tag: "realm_update_start")
125-      do_update_realm(realm_name, update_attrs)
126-    else
127-      _ =
128-        Logger.warning("Rejecting update for realm #{realm_name}",
129-          tag: "invalid_update_parameters"
130-        )
131-
132-      {:error, :invalid_update_parameters}
133-    end
134-  end
135-
136-  defp do_update_realm(realm_name, update_attrs) do
137-    %{
138-      jwt_public_key_pem: new_jwt_public_key_pem,
139-      device_registration_limit: new_limit,
140-      datastream_maximum_storage_retention: new_retention
141-    } = update_attrs
142-
143-    with realm when is_map(realm) <- Queries.get_realm(realm_name),
144-         :ok <- update_jwt_public_key_pem(realm_name, new_jwt_public_key_pem),
145-         :ok <- update_device_registration_limit(realm_name, new_limit),
146:         :ok <- update_datastream_maximum_storage_retention(realm_name, new_retention) do
147:      updated_realm = merge_realm_status(realm, update_attrs)
148-
149-      _ = Logger.info("Successful update of realm #{realm_name}", tag: "realm_update_success")
150-
151-      {:ok, updated_realm}
152-    end
153-  end
154-
155-  defp update_jwt_public_key_pem(_realm_name, nil), do: :ok
156-
157-  defp update_jwt_public_key_pem(realm_name, new_jwt_public_key_pem) do
158-    case Queries.update_public_key(realm_name, new_jwt_public_key_pem) do
159-      {:ok, %Xandra.Void{}} ->
160-        :ok
161-
162-      {:error, reason} ->
163-        _ =
164-          Logger.warning(
165-            "Cannot update JWT public key for realm #{realm_name}, error #{inspect(reason)}",
166-            tag: "update_public_key_fail"
167-          )
168-
169-        {:error, :update_public_key_fail}
170-    end
171-  end
172-
--
184-    :ok
185-  end
186-
187:  defp update_datastream_maximum_storage_retention(realm_name, new_retention)
188-       when is_integer(new_retention) and new_retention != 0 do
189-    _ =
190-      Logger.info("Updating datastream maximum storage retention",
191:        tag: "update_datastream_maximum_storage_retention"
192-      )
193-
194-    set_datastream_maximum_storage_retention(realm_name, new_retention)
195-  end
196-
197-  # ScyllaDB considers TTL=0 as unset, see
198-  # https://opensource.docs.scylladb.com/stable/cql/time-to-live.html#notes
199:  defp update_datastream_maximum_storage_retention(realm_name, 0) do
200-    _ =
201-      Logger.info("Removing datastream maximum storage retention",
202-        tag: "remove_datastream_maximum_storage_retention"
203-      )
204-
205-    delete_datastream_maximum_storage_retention(realm_name)
206-  end
207-
208:  defp update_datastream_maximum_storage_retention(_realm_name, nil) do
209-    :ok
210-  end
211-
212-  defp delete_device_registration_limit(realm_name) do
213-    case Queries.delete_device_registration_limit(realm_name) do
214-      {:ok, %Xandra.Void{}} ->
215-        :ok
216-
217-      {:error, reason} ->
218-        _ =
219-          Logger.warning(
220-            "Cannot delete device registration limit, error #{inspect(reason)}",
221-            tag: "delete_device_registration_limit"
222-          )
223-
224-        {:error, :delete_device_registration_limit_fail}
225-    end
226-  end
227-
228-  defp delete_datastream_maximum_storage_retention(realm_name) do
229-    case Queries.delete_datastream_maximum_storage_retention(realm_name) do
230-      {:ok, %Xandra.Void{}} ->
231-        :ok
232-
233-      {:error, reason} ->
--
274-    end
275-  end
276-
277:  defp is_realm_update_valid?(update_attrs) do
278-    # TODO from ScyllaDB >= 5.3, replication can be altered
279-    update_valid? =
280-      update_attrs.replication_factor == nil && update_attrs.replication_class == nil &&
281-        update_attrs.datacenter_replication_factors == %{}
282-
283-    unless update_valid? do
284-      _ =
285-        Logger.warning("Trying to update replication values for realm",
286-          tag: "invalid_replication_value_update"
287-        )
288-    end
289-
290-    update_valid?
291-  end
292-
293-  # If device_registration_limit is nil, it means it was not present in the update
294-  # payload. Act accordingly.
295:  defp merge_realm_status(old_realm, %{device_registration_limit: nil} = update_attrs) do
296-    attrs = Map.delete(update_attrs, :device_registration_limit)
297:    merge_realm_status(old_realm, attrs)
298-  end
299-
300:  defp merge_realm_status(old_realm, %{device_registration_limit: :remove_limit} = update_attrs) do
301-    attrs = Map.delete(update_attrs, :device_registration_limit)
302:    merge_realm_status(%{old_realm | device_registration_limit: nil}, attrs)
303-  end
304-
305-  # If datastream_maximum_storage_retention is nil, it means it was not present in the update
306-  # payload. Act accordingly.
307:  defp merge_realm_status(old_realm, %{datastream_maximum_storage_retention: nil} = update_attrs) do
308-    attrs = Map.delete(update_attrs, :datastream_maximum_storage_retention)
309:    merge_realm_status(old_realm, attrs)
310-  end
311-
312-  # Differently from device_registration_limit, datastream_maximum_storage_retention
313-  # uses 0 to unset the value. This follows the convention used by ScyllaDB/Cassandra.
314:  defp merge_realm_status(
315-         old_realm,
316-         %{datastream_maximum_storage_retention: 0} = update_attrs
317-       ) do
318-    attrs = Map.delete(update_attrs, :datastream_maximum_storage_retention)
319:    merge_realm_status(%{old_realm | datastream_maximum_storage_retention: nil}, attrs)
320-  end
321-
322:  defp merge_realm_status(old_realm, update_attrs) do
323-    Map.merge(old_realm, update_attrs, fn _k, old_value, new_value ->
324-      if new_value == nil, do: old_value, else: new_value
325-    end)
326-  end
327-end

$ curl -s --max-time 15 "https://api.github.com/repos/astarte-platform/astarte/git/trees/v1.2.0?recursive=1" | rg -o '"path": "[^"]*"' | rg -i "housekeeping_api.*rpc|protobuf|\.pb\.ex" | head
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/rpc"
"path": "apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/rpc/housekeeping.ex"

$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/v1.2.0/apps/astarte_housekeeping_api/lib/astarte_housekeeping_api/rpc/housekeeping.ex" | rg -n -B4 -A30 "update_realm|CreateRealm|UpdateRealm"
18-
19-defmodule Astarte.Housekeeping.API.RPC.Housekeeping do
20-  alias Astarte.RPC.Protocol.Housekeeping.{
21-    Call,
22:    CreateRealm,
23-    DeleteRealm,
24-    DoesRealmExist,
25-    DoesRealmExistReply,
26-    GenericErrorReply,
27-    GenericOkReply,
28-    GetHealth,
29-    GetHealthReply,
30-    GetRealm,
31-    GetRealmReply,
32-    GetRealmsList,
33-    GetRealmsListReply,
34-    RemoveLimit,
35-    Reply,
36-    SetLimit,
37:    UpdateRealm
38-  }
39-
40-  alias Astarte.Housekeeping.API.Config
41-  alias Astarte.Housekeeping.API.Realms.Realm
42-
43-  @rpc_client Config.rpc_client!()
44-  @destination Astarte.RPC.Protocol.Housekeeping.amqp_queue()
45-
46-  def create_realm(
47-        %Realm{
48-          realm_name: realm_name,
49-          jwt_public_key_pem: pem,
50-          replication_class: "SimpleStrategy",
51-          replication_factor: replication_factor,
52-          device_registration_limit: device_registration_limit,
53-          datastream_maximum_storage_retention: datastream_maximum_storage_retention
54-        },
55-        opts
56-      ) do
57:    %CreateRealm{
58-      realm: realm_name,
59-      async_operation: Keyword.get(opts, :async_operation, true),
60-      jwt_public_key_pem: pem,
61-      replication_class: :SIMPLE_STRATEGY,
62-      replication_factor: replication_factor,
63-      device_registration_limit: device_registration_limit,
64-      datastream_maximum_storage_retention: datastream_maximum_storage_retention
65-    }
66-    |> encode_call(:create_realm)
67-    |> @rpc_client.rpc_call(@destination, Config.rpc_timeout!())
68-    |> decode_reply()
69-    |> extract_reply()
70-  end
71-
72-  def create_realm(
73-        %Realm{
74-          realm_name: realm_name,
75-          jwt_public_key_pem: pem,
76-          replication_class: "NetworkTopologyStrategy",
77-          datacenter_replication_factors: replication_factors_map,
78-          device_registration_limit: device_registration_limit,
79-          datastream_maximum_storage_retention: datastream_maximum_storage_retention
80-        },
81-        opts
82-      ) do
83:    %CreateRealm{
84-      realm: realm_name,
85-      async_operation: Keyword.get(opts, :async_operation, true),
86-      jwt_public_key_pem: pem,
87-      replication_class: :NETWORK_TOPOLOGY_STRATEGY,
88-      datacenter_replication_factors: Enum.to_list(replication_factors_map),
89-      device_registration_limit: device_registration_limit,
90-      datastream_maximum_storage_retention: datastream_maximum_storage_retention
91-    }
92-    |> encode_call(:create_realm)
93-    |> @rpc_client.rpc_call(@destination)
94-    |> decode_reply()
95-    |> extract_reply()
96-  end
97-
98:  def update_realm(%Realm{
99-        realm_name: realm_name,
100-        jwt_public_key_pem: pem,
101-        replication_class: replication_class,
102-        replication_factor: replication_factor,
103-        datacenter_replication_factors: replication_factors_map,
104-        device_registration_limit: limit,
105-        datastream_maximum_storage_retention: max_retention
106-      }) do
107:    %UpdateRealm{
108-      realm: realm_name,
109-      jwt_public_key_pem: pem,
110-      replication_class: replication_class_to_atom(replication_class),
111-      replication_factor: replication_factor,
112-      datacenter_replication_factors: replication_factors_map,
113-      device_registration_limit: device_registration_limit_to_proto(limit),
114-      datastream_maximum_storage_retention:
115-        datastream_maximum_storage_retention_to_proto(max_retention)
116-    }
117:    |> encode_call(:update_realm)
118-    |> @rpc_client.rpc_call(@destination, Config.rpc_timeout!())
119-    |> decode_reply()
120-    |> extract_reply()
121-  end
122-
123-  def list_realms do
124-    %GetRealmsList{}
125-    |> encode_call(:get_realms_list)
126-    |> @rpc_client.rpc_call(@destination, Config.rpc_timeout!())
127-    |> decode_reply()
128-    |> extract_reply()
129-  end
130-
131-  def delete_realm(realm_name, opts) do
132-    async_operation = Keyword.get(opts, :async_operation, true)
133-
134-    %DeleteRealm{realm: realm_name, async_operation: async_operation}
135-    |> encode_call(:delete_realm)
136-    |> @rpc_client.rpc_call(@destination, Config.rpc_timeout!())
137-    |> decode_reply()
138-    |> extract_reply()
139-  end
140-
141-  def get_health do
142-    %GetHealth{}
143-    |> encode_call(:get_health)
144-    |> @rpc_client.rpc_call(@destination, Config.rpc_timeout!())
145-    |> decode_reply()
146-    |> extract_reply()
147-  end

Let me check the phase handoffs for the recorded upstream measurement of PATCH retention behaviour.
$ rg -l -i "retention" docs/handoff/ docs/COMPATIBILITY.md 2>/dev/null; echo ---; rg -n -i "retention|zero|0.*unset|unset.*0|PATCH" docs/handoff/phase-11-handoff.md 2>/dev/null | head -40
docs/handoff/phase-1-memory.md
---

$ rg -l -i "housekeep|trickle|zero.*retention|datastream_maximum" docs/ internal/ | head; echo ---; rg -n -i "retention" docs/COMPATIBILITY.md | head
docs/api/astarte_housekeeping_api.yaml
docs/OPERATIONS.md
docs/COMPATIBILITY.md
docs/api/astrate_native_api.yaml
docs/site/pairing-and-security.md
docs/api/astarte_realm_management_api.yaml
docs/site/troubleshooting.md
docs/ROADMAP.md
docs/site/configuration-reference.md
docs/swagger-ui/index.html
---

$ rg -n -i "retention" docs/OPERATIONS.md | head; echo ---; rg -n "retention" internal/store/*.go | head -20
---
internal/store/retention.go:9:// datastream_maximum_storage_retention (#72): datastream rows older than the
internal/store/retention.go:10:// ceiling are deleted regardless of their interface's own retention policy —
internal/store/retention.go:18:		SELECT id, datastream_maximum_storage_retention FROM realms
internal/store/retention.go:19:		WHERE datastream_maximum_storage_retention IS NOT NULL
internal/store/retention.go:20:		  AND datastream_maximum_storage_retention > 0`)
internal/store/retention.go:22:		return fmt.Errorf("store: listing retention ceilings: %w", err)
internal/store/retention.go:33:			return fmt.Errorf("store: scanning retention ceiling: %w", err)
internal/store/retention.go:39:		return fmt.Errorf("store: listing retention ceilings: %w", err)
internal/store/retention.go:49:				firstErr = fmt.Errorf("store: sweeping %s of realm id %d past its retention ceiling: %w", table, c.id, err)
internal/store/retention_test.go:13:// TestEnforceRealmRetentionCeilings covers the realm-level retention sweep
internal/store/interfaces_test.go:72:			valueType, reliability, retention, policy string
internal/store/interfaces_test.go:76:			SELECT value_type, reliability, retention, database_retention_policy, explicit_timestamp, allow_unset
internal/store/interfaces_test.go:78:			Scan(&valueType, &reliability, &retention, &policy, &explicitTS, &allowUnset)
internal/store/interfaces_test.go:82:		if valueType != "double" || reliability != "guaranteed" || retention != "discard" ||
internal/store/interfaces_test.go:85:				valueType, reliability, retention, policy, explicitTS, allowUnset)
internal/store/datastreams_test.go:715:			{"endpoint": "/ttl", "type": "double", "database_retention_policy": "use_ttl", "database_retention_ttl": 3600},
internal/store/datastreams_test.go:727:			{"endpoint": "/%{id}/a", "type": "double", "database_retention_policy": "use_ttl", "database_retention_ttl": 3600},
internal/store/datastreams_test.go:728:			{"endpoint": "/%{id}/b", "type": "double", "database_retention_policy": "use_ttl", "database_retention_ttl": 3600}
internal/store/datastreams_test.go:786:// (docs/DESIGN.md §2.5): applying it registers a policy_retention job on both
internal/store/datastreams_test.go:795:			WHERE proc_name = 'policy_retention'

$ ls internal/auth/; echo ---; rg -n "func NewMiddleware|func (m \*Middleware) RequireStatic|instanceKeys|RequireStatic" internal/auth/*.go | head
cache.go
cache_test.go
claims.go
claims_test.go
jwt.go
jwt_test.go
middleware.go
middleware_test.go
testdata
---
internal/auth/middleware.go:40:func NewMiddleware(keys KeySource) *Middleware {
internal/auth/middleware.go:95:// RequireStatic guards an instance-level route (Housekeeping) with a fixed
internal/auth/middleware.go:98:func (m *Middleware) RequireStatic(claim Claim, keysPEM []string) func(http.Handler) http.Handler {
internal/auth/middleware_test.go:46:		m.RequireStatic(ClaimHousekeeping, hkKeys)(okHandler))

Let me check the auth middleware to see if a container-free PATCH test is feasible, and look at the store realm handling.
$ rg -n "DatastreamMaximumStorageRetention" internal/store/ | rg -v _test | head
internal/store/realms.go:23:	DatastreamMaximumStorageRetention *int64
internal/store/realms.go:35:	DatastreamMaximumStorageRetention *int64
internal/store/realms.go:62:		nr.Name, keys, nr.CACertificatePEM, nr.CAPrivateKeySealed, nr.DeviceRegistrationLimit, nr.DatastreamMaximumStorageRetention)
internal/store/realms.go:246:		&r.DatastreamMaximumStorageRetention, &r.CreatedAt); err != nil {

→ Read internal/auth/middleware.go
Feasibility confirmed — `RequireStatic` never touches the key source, so a container-free HTTP rig is possible, but the PATCH side is unchanged, so no new test there. Now the analysis is settled. Let me verify my upstream reading one final way — check the repo's own realm-management doc convention and the retention sweep, which both treat 0 as unset:

$ rg -n "0 means unset|> 0" docs/api/astarte_realm_management_api.yaml internal/store/retention.go
docs/api/astarte_realm_management_api.yaml:526:                    description: Maximum datastream storage retention in seconds. 0 means unset.
internal/store/retention.go:20:		  AND datastream_maximum_storage_retention > 0`)

The evidence is conclusive. Summary of what I found before touching code:

- **Upstream v1.2.0** (the pinned compatibility target): engine.ex folds retention `0` to delete/unset — `update_datastream_maximum_storage_retention(realm, 0)` calls `delete_datastream_maximum_storage_retention`, `insert_datastream_max_retention(_, 0)` is a no-op, and `merge_realm_status` maps `0 → nil` with the comment "datastream_maximum_storage_retention uses 0 to unset the value... follows the convention used by ScyllaDB/Cassandra". The changeset type accepts `n >= 0`, so 0 is valid, not an error.
- **Astrate itself** treats 0 as unset everywhere else: `internal/store/retention.go:20` sweeps only `> 0`, and `docs/api/astarte_realm_management_api.yaml:526` documents "0 means unset".

So upstream parity is neither of the task's two literal branches (0-invalid-422, or 0-round-trips-as-literal). It is: **0 is valid and means unset**, which the PATCH side already implements (`http.go:217`). The asymmetry is only in POST, which persists literal 0. I'll fix POST to fold 0 → nil.
← Edit internal/housekeeping/service.go
Index: /root/astrate-mule/internal/housekeeping/service.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service.go
+++ /root/astrate-mule/internal/housekeeping/service.go
@@ -147,8 +147,16 @@
 	return nil, fmt.Errorf("%w: datastream_maximum_storage_retention must be non-negative", ErrValidation)
 }
 if retention == nil {
 	retention = s.defaultRetention
+} else if *retention == 0 {
+	// 0 ≡ unset, upstream parity (measured Astarte v1.2.0): the engine
+	// folds retention 0 to nil at the DB layer ("ScyllaDB considers
+	// TTL=0 as unset"), PATCH maps 0 to ClearRetention, and the store's
+	// retention sweep only treats values > 0 as ceilings. An explicit 0
+	// is still explicit: it must not pick up the configured default
+	// (#73), so fold before the nil-injection above.
+	retention = nil
 }
 
 realmCA, err := ca.Generate(name, 0)
 if err != nil {

← Edit internal/housekeeping/service_test.go
Index: /root/astrate-mule/internal/housekeeping/service_test.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service_test.go
+++ /root/astrate-mule/internal/housekeeping/service_test.go
@@ -177,8 +177,61 @@
 		}
 	})
 }
 
+// TestCreateRealmZeroRetentionUnsets pins the upstream-parity behaviour of
+// datastream_maximum_storage_retention: 0 is a valid value but means unset,
+// exactly like PATCH's 0→ClearRetention — upstream's engine folds retention 0
+// to nil at the DB layer ("ScyllaDB considers TTL=0 as unset", Astarte
+// v1.2.0). An explicit 0 must not pick up the configured default (#73).
+func TestCreateRealmZeroRetentionUnsets(t *testing.T) {
+	ctx := context.Background()
+	zero := int64(0)
+	def := int64(3600)
+
+	t.Run("explicit zero stores nil", func(t *testing.T) {
+		st := &fakeStore{}
+		svc := NewService(st, &fakeSealer{}, nil, nil)
+
+		rv, err := svc.CreateRealm(ctx, "zerorealm", "key", nil, &zero)
+		if err != nil {
+			t.Fatalf("CreateRealm: %v", err)
+		}
+		if got := st.created[0].DatastreamMaximumStorageRetention; got != nil {
+			t.Errorf("store received retention = %v, want nil (unset)", got)
+		}
+		if rv.DatastreamMaximumStorageRetention != nil {
+			t.Errorf("view retention = %v, want nil", *rv.DatastreamMaximumStorageRetention)
+		}
+	})
+
+	t.Run("explicit zero beats configured default", func(t *testing.T) {
+		st := &fakeStore{}
+		svc := NewService(st, &fakeSealer{}, nil, nil).
+			WithDefaultDatastreamMaximumStorageRetention(&def)
+
+		if _, err := svc.CreateRealm(ctx, "zerorealm2", "key", nil, &zero); err != nil {
+			t.Fatalf("CreateRealm: %v", err)
+		}
+		if got := st.created[0].DatastreamMaximumStorageRetention; got != nil {
+			t.Errorf("store received retention = %v, want nil, not the default %d", got, def)
+		}
+	})
+
+	t.Run("positive retention is stored", func(t *testing.T) {
+		st := &fakeStore{}
+		svc := NewService(st, &fakeSealer{}, nil, nil)
+
+		pos := int64(60)
+		if _, err := svc.CreateRealm(ctx, "zerorealm3", "key", nil, &pos); err != nil {
+			t.Fatalf("CreateRealm: %v", err)
+		}
+		if got := st.created[0].DatastreamMaximumStorageRetention; got == nil || *got != pos {
+			t.Errorf("store received retention = %v, want 60", got)
+		}
+	})
+}
+
 // TestDeleteRealmGating covers #75 at the service level: the deletion-disabled
 // flag answers ErrDeletionDisabled before any store access, and a realm with
 // connected devices answers ErrConnectedDevicesPresent without deleting.
 func TestDeleteRealmGating(t *testing.T) {

Now run the gates:
$ gofmt -l . ; echo "gofmt-exit=$?"
gofmt-exit=0

$ go vet ./internal/housekeeping/... 2>&1 | tail -20; echo "vet-exit=$?"
vet-exit=0

$ go test ./internal/housekeeping/... 2>&1 | tail -20; echo "test-exit=$?"
ok  	github.com/astrate-platform/astrate/internal/housekeeping	0.041s
test-exit=0

← Edit internal/housekeeping/service.go
Index: /root/astrate-mule/internal/housekeeping/service.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service.go
+++ /root/astrate-mule/internal/housekeeping/service.go
@@ -147,16 +147,8 @@
 	return nil, fmt.Errorf("%w: datastream_maximum_storage_retention must be non-negative", ErrValidation)
 }
 if retention == nil {
 	retention = s.defaultRetention
-} else if *retention == 0 {
-	// 0 ≡ unset, upstream parity (measured Astarte v1.2.0): the engine
-	// folds retention 0 to nil at the DB layer ("ScyllaDB considers
-	// TTL=0 as unset"), PATCH maps 0 to ClearRetention, and the store's
-	// retention sweep only treats values > 0 as ceilings. An explicit 0
-	// is still explicit: it must not pick up the configured default
-	// (#73), so fold before the nil-injection above.
-	retention = nil
 }
 
 realmCA, err := ca.Generate(name, 0)
 if err != nil {

Tests pass. Now the acceptance check — my new test must fail without the implementation change. I'll temporarily restore the old logic, run just the new test, then put my change back.
$ go test -run TestCreateRealmZeroRetentionUnsets ./internal/housekeeping/... 2>&1 | tail -25; echo "exit=$?"
--- FAIL: TestCreateRealmZeroRetentionUnsets (0.01s)
    --- FAIL: TestCreateRealmZeroRetentionUnsets/explicit_zero_stores_nil (0.00s)
        service_test.go:200: store received retention = 0x7f489245d0, want nil (unset)
        service_test.go:203: view retention = 0, want nil
    --- FAIL: TestCreateRealmZeroRetentionUnsets/explicit_zero_beats_configured_default (0.00s)
        service_test.go:216: store received retention = 0x7f489245d0, want nil, not the default 3600
FAIL
FAIL	github.com/astrate-platform/astrate/internal/housekeeping	0.034s
FAIL
exit=0

← Edit internal/housekeeping/service.go
Index: /root/astrate-mule/internal/housekeeping/service.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service.go
+++ /root/astrate-mule/internal/housekeeping/service.go
@@ -147,8 +147,16 @@
 	return nil, fmt.Errorf("%w: datastream_maximum_storage_retention must be non-negative", ErrValidation)
 }
 if retention == nil {
 	retention = s.defaultRetention
+} else if *retention == 0 {
+	// 0 ≡ unset, upstream parity (measured Astarte v1.2.0): the engine
+	// folds retention 0 to nil at the DB layer ("ScyllaDB considers
+	// TTL=0 as unset"), PATCH maps 0 to ClearRetention, and the store's
+	// retention sweep only treats values > 0 as ceilings. An explicit 0
+	// is still explicit: it must not pick up the configured default
+	// (#73), so fold before the nil-injection above.
+	retention = nil
 }
 
 realmCA, err := ca.Generate(name, 0)
 if err != nil {

The new test fails on the old code as required. Restoring the fix:
$ gofmt -l . ; go vet ./... 2>&1 | tail -5; go test ./... 2>&1 | tail -40; echo "exit=${PIPESTATUS:-$?}"
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
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
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)
exit=0

Full gate green. Let me force an uncached run of the changed package and try a race build:
$ go test -count=1 ./internal/housekeeping/... 2>&1 | tail -3; echo "---"; go test -race -count=1 ./internal/housekeeping/... 2>&1 | tail -5; echo "race-exit=$?"
ok  	github.com/astrate-platform/astrate/internal/housekeeping	0.041s
---
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/housekeeping	0.026s
FAIL
race-exit=0

$ git diff --stat; echo ---; git diff
internal/housekeeping/service.go      |  8 ++++++
 internal/housekeeping/service_test.go | 53 +++++++++++++++++++++++++++++++++++
 2 files changed, 61 insertions(+)
---
diff --git a/internal/housekeeping/service.go b/internal/housekeeping/service.go
index 5ebb832..43b5aad 100644
--- a/internal/housekeeping/service.go
+++ b/internal/housekeeping/service.go
@@ -148,6 +148,14 @@ func (s *Service) CreateRealm(ctx context.Context, name, jwtPublicKeyPEM string,
 	}
 	if retention == nil {
 		retention = s.defaultRetention
+	} else if *retention == 0 {
+		// 0 ≡ unset, upstream parity (measured Astarte v1.2.0): the engine
+		// folds retention 0 to nil at the DB layer ("ScyllaDB considers
+		// TTL=0 as unset"), PATCH maps 0 to ClearRetention, and the store's
+		// retention sweep only treats values > 0 as ceilings. An explicit 0
+		// is still explicit: it must not pick up the configured default
+		// (#73), so fold before the nil-injection above.
+		retention = nil
 	}
 
 	realmCA, err := ca.Generate(name, 0)
diff --git a/internal/housekeeping/service_test.go b/internal/housekeeping/service_test.go
index a1166c9..a3a6023 100644
--- a/internal/housekeeping/service_test.go
+++ b/internal/housekeeping/service_test.go
@@ -178,6 +178,59 @@ func TestCreateRealmDefaultRetention(t *testing.T) {
 	})
 }
 
+// TestCreateRealmZeroRetentionUnsets pins the upstream-parity behaviour of
+// datastream_maximum_storage_retention: 0 is a valid value but means unset,
+// exactly like PATCH's 0→ClearRetention — upstream's engine folds retention 0
+// to nil at the DB layer ("ScyllaDB considers TTL=0 as unset", Astarte
+// v1.2.0). An explicit 0 must not pick up the configured default (#73).
+func TestCreateRealmZeroRetentionUnsets(t *testing.T) {
+	ctx := context.Background()
+	zero := int64(0)
+	def := int64(3600)
+
+	t.Run("explicit zero stores nil", func(t *testing.T) {
+		st := &fakeStore{}
+		svc := NewService(st, &fakeSealer{}, nil, nil)
+
+		rv, err := svc.CreateRealm(ctx, "zerorealm", "key", nil, &zero)
+		if err != nil {
+			t.Fatalf("CreateRealm: %v", err)
+		}
+		if got := st.created[0].DatastreamMaximumStorageRetention; got != nil {
+			t.Errorf("store received retention = %v, want nil (unset)", got)
+		}
+		if rv.DatastreamMaximumStorageRetention != nil {
+			t.Errorf("view retention = %v, want nil", *rv.DatastreamMaximumStorageRetention)
+		}
+	})
+
+	t.Run("explicit zero beats configured default", func(t *testing.T) {
+		st := &fakeStore{}
+		svc := NewService(st, &fakeSealer{}, nil, nil).
+			WithDefaultDatastreamMaximumStorageRetention(&def)
+
+		if _, err := svc.CreateRealm(ctx, "zerorealm2", "key", nil, &zero); err != nil {
+			t.Fatalf("CreateRealm: %v", err)
+		}
+		if got := st.created[0].DatastreamMaximumStorageRetention; got != nil {
+			t.Errorf("store received retention = %v, want nil, not the default %d", got, def)
+		}
+	})
+
+	t.Run("positive retention is stored", func(t *testing.T) {
+		st := &fakeStore{}
+		svc := NewService(st, &fakeSealer{}, nil, nil)
+
+		pos := int64(60)
+		if _, err := svc.CreateRealm(ctx, "zerorealm3", "key", nil, &pos); err != nil {
+			t.Fatalf("CreateRealm: %v", err)
+		}
+		if got := st.created[0].DatastreamMaximumStorageRetention; got == nil || *got != pos {
+			t.Errorf("store received retention = %v, want 60", got)
+		}
+	})
+}
+
 // TestDeleteRealmGating covers #75 at the service level: the deletion-disabled
 // flag answers ErrDeletionDisabled before any store access, and a realm with
 // connected devices answers ErrConnectedDevicesPresent without deleting.

Change verified. Let me check one adjacent fact worth capturing as a follow-up — whether upstream folds `device_registration_limit` 0 the same way:
$ curl -s --max-time 15 "https://raw.githubusercontent.com/astarte-platform/astarte/v1.2.0/apps/astarte_housekeeping/lib/astarte_housekeeping/engine.ex" | rg -n -B2 -A10 "defp update_device_registration_limit"
171-  end
172-
173:  defp update_device_registration_limit(realm_name, new_limit) when is_integer(new_limit) do
174-    _ = Logger.info("Updating device registration limit", tag: "update_device_registration_limit")
175-    set_device_registration_limit(realm_name, new_limit)
176-  end
177-
178:  defp update_device_registration_limit(realm_name, :remove_limit) do
179-    _ = Logger.info("Removing device registration limit", tag: "remove_device_registration_limit")
180-    delete_device_registration_limit(realm_name)
181-  end
182-
183:  defp update_device_registration_limit(_realm_name, nil) do
184-    :ok
185-  end
186-
187-  defp update_datastream_maximum_storage_retention(realm_name, new_retention)
188-       when is_integer(new_retention) and new_retention != 0 do
189-    _ =
190-      Logger.info("Updating datastream maximum storage retention",
191-        tag: "update_datastream_maximum_storage_retention"
192-      )
193-

`device_registration_limit` has the same POST-0 asymmetry in astrate (persists literal 0, upstream folds to nil), but its PATCH already matches. Worth a follow-up task line. Let me check todo.md:
→ Read .mule/todo.md
✗ Edit .mule/todo.md failed
Error: Could not find oldString in the file. It must match exactly, including whitespace, indentation, and line endings.
One follow-up is warranted: `device_registration_limit` has the exact same POST asymmetry (parallel evidence from the same upstream code). Appending a task line:

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -72,8 +72,10 @@
     bash tools/mule.sh review
 
 ## Queue
 
+- [ ] hk-zero-reglimit-create-asymmetry: POST /realms persists device_registration_limit 0 as a literal 0 (service.go:143 only rejects < 0) while upstream Astarte v1.2.0 folds 0 to nil at create (`insert_realm` does `if device_limit == 0, do: nil`, so GET returns null); PATCH 0 already matches upstream (sets 0, http.go:208-214). Fold 0 → nil in CreateRealm for parity, same pattern as the retention fix, and pin with a container-free test.
+
 - [x] flowapi-autorestart-terminal-failure: `restartWithBackoff` (internal/flowapi/service.go:585-621) retries forever on a permanent rebuild failure — a flow whose pipeline was deleted (service.go:360-363 wraps it as ErrValidation) or whose config/definition no longer resolves keeps marking itself failed and rebuilding every ≤30s with no exit short of flow deletion. Classify terminal vs transient (ErrValidation → permanent; store/manager errors → transient) via a pure helper and stop the loop on terminal failures; add a container-free unit test on the helper. [auto]
 - [ ] broker-external-bus-intake [legion]: implement the second Intake implementation named in the TODO at internal/broker/intake.go:56-62 — an external-bus-backed (e.g. NATS JetStream) Intake for multi-instance deployment / restart survival, reproducing the embedded broker's per-device ordering, deferred-ack backpressure and QoS 0 drop semantics; an implementation of this interface needs containerised integration tests, so it runs on the Legion Go.
 - [x] store-todo-lttb-stale: correct the stale TODO at internal/store/store.go:139-143 — it claims Downsample should switch to timescaledb_toolkit lttb() when the extension is present, but that path already exists as DownsampleLTTB (internal/store/datastreams.go:375), caller-selected via HasToolkitLTTB (store.go:154) in the AppEngine data path (internal/appengine/data.go:176-177); the comment now misleads a reader into thinking the switch is unbuilt. Update or delete it.
 - [ ] hygiene-govulncheck [legion]: on the Legion Go, run govulncheck ./... on a fresh ~/astrate clone and file a task line for any REACHABLE advisory it reports (name the CVE and the call path, per the hygiene recipe); this box cannot build it — go install golang.org/x/vuln/cmd/govulncheck@latest is OOM-killed here (3.7GB, 4 cores), so the recipe's highest-priority check never runs unless done there.

Final verification of the complete working-tree state:
$ gofmt -l . && go test -count=1 ./internal/housekeeping/... 2>&1 | tail -2; echo ---; git diff --stat
ok  	github.com/astrate-platform/astrate/internal/housekeeping	0.040s
---
 .mule/todo.md                         |  2 ++
 internal/housekeeping/service.go      |  8 ++++++
 internal/housekeeping/service_test.go | 53 +++++++++++++++++++++++++++++++++++
 3 files changed, 63 insertions(+)

Done. Report follows.

**One-line meaning:** The survey flagged that `PATCH` treats retention 0 as clear while `POST` persists it literal; I pinned the upstream-parity behaviour — which, after reading upstream sources, is not either branch the task guessed.

```
Done: POST /realms now folds explicit datastream_maximum_storage_retention 0 to nil (unset), matching upstream v1.2.0 and the existing PATCH 0→ClearRetention; pinned container-free.
Files: internal/housekeeping/service.go, internal/housekeeping/service_test.go, .mule/todo.md
Verified: gofmt -l (clean); go vet ./... (pass); go test ./... and go test -count=1 ./internal/housekeeping/... (pass) -> pass; new test run against reverted old logic -> FAIL (accepted); go test -race not runnable here (39-bit VMA kernel, documented Pi limitation; no goroutines/locks touched) -> n/a
Unsure: I did NOT follow either suggested branch, because both contradict measured upstream. Task said "0 invalid (recommended)"→422, or "0 valid"→PATCH round-trips literal 0. Read upstream Astarte v1.2.0 (the pinned target) instead: NonNegativeIntegerOrUnsetType accepts n>=0 (0 is valid, never 422); engine.ex folds retention 0 to delete/unset ("ScyllaDB considers TTL=0 as unset", merge_realm_status maps 0→nil); POST 0 inserts nothing → GET null. Astrate's own store sweep (retention.go:20) and RM spec (yaml:526) already treat 0 as unset. So parity = "0 valid but ≡ null", which PATCH (http.go:217) already implements — the real asymmetry was only POST persisting literal 0. I fixed that: CreateRealm folds 0→nil before default injection, so explicit 0 does not grab the #73 default (matching upstream put_new semantics). Tests fail on old code (verified by temporary revert). Existing integration test ZeroRetentionUnsets (http_test.go) already pins PATCH 0→clear and is untouched. If you wanted strict branch A/B instead, say so.
Follow-ups: hk-zero-reglimit-create-asymmetry (device_registration_limit has the same POST-0 gap; PATCH already matches)
```
