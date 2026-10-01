slug: fdo-rc6-scope-delta-for-giulio
verdict: done
at:  41af3c4
ran: 2026-10-01T20:46:04Z on DietPi in 479s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ ls -a /root/astrate-mule && echo "---" && ls /root/astrate-mule/.mule
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
astrate
bench
cmd
deploy
docker-compose.yml
docs
examples
flow-container-echo
go.mod
go.sum
internal
migrations
opencode.json
pkg
site-dist
test
tools
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
recipes
reports
reviews
task.md
tasks
todo.md

→ Read .mule/for-giulio.md

$ cd /root/astrate-mule && rg -il 'fdo|ownership_voucher|to0' --glob '!.git' | head -50; echo "--- ls .mule/research .trickle 2>/dev/null"; ls -a .trickle 2>/dev/null | head
docs/api/astarte_pairing_api.yaml
docs/COMPATIBILITY.md
internal/pairing/service.go
internal/pairing/http.go
--- ls .mule/research .trickle 2>/dev/null

$ which gh git go elixir 2>&1; echo "---"; gh auth status 2>&1 | head -5
/usr/bin/gh
/usr/bin/git
/usr/bin/go
---
github.com
  ✓ Logged in to github.com account atsetilam (/root/.config/gh/hosts.yml)
  - Active account: true
  - Git operations protocol: https
  - Token: gho_************************************

$ cd /root/astrate-mule && rg -n 'fdo|FDO' internal/pairing/http.go internal/pairing/service.go docs/COMPATIBILITY.md | head -40
internal/pairing/http.go:79:	// FDO-flow devices probe it before they hold credentials. No rate
internal/pairing/service.go:368:// design — FDO-flow devices probe it before they hold credentials).
docs/COMPATIBILITY.md:361:    answers 404 and an unhealthy database 503, and FDO-flow devices can probe

$ cd /root/astrate-mule && timeout 60 gh issue view 78 2>&1 | head -80
title:	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
state:	OPEN
author:	atsetilam (Giulio Malitesta)
labels:	enhancement, milestone-4.0, upstream-parity
comments:	1
assignees:	
projects:	
milestone:	
issue-type:	
parent:	
sub-issues:	
sub-issues-completed:	
blocked-by:	
blocking:	
number:	78
--
Upstream 1.3 experimental, default-on in 1.4: TO2 CBOR endpoints under `/v1/{realm}/fdo/101/msg/{60..70}` (auto-registers the device mid-flow, hands over credentials_secret), owner side `POST/GET .../fdo/owner_keys` + ownership vouchers with TO0 claim, secrets in Vault/OpenBao, rendezvous server required.

Large surface; needs a product call (who runs a rendezvous server?) before any implementation. Parked.
Size L.\n\n---\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).

$ cd /root/astrate-mule && timeout 60 gh issue view 78 --comments 2>&1 | tail -20
## Verified facts (2026-08-23, upstream master)

- Upstream integrates FDO **inside the Pairing service**: `fdo_onboarding_controller.ex`, session plugs (`setup_fdo`, `fdo_session`), CBOR codec + voucher queries in `astarte_data_access`, and two dedicated libs `astarte_fdo` / `astarte_fdo_core` (SECO Mind, 2025; bugfixes through 2026-05). Not an optional sidecar.
- Device side is the official SDK `astarte-device-fdo-rust` (crate on docs.rs, pushed 2026-08-21) — actively maintained.
- Upstream does NOT implement manufacturing / rendezvous services either; those come from the FIDO Alliance reference stack (`fdo-rs/fido-device-onboard-rs`, used by Red Hat/Fedora IoT).

## Scope decision

Astrate implements **only the last mile**, mirroring upstream's cut:
- owner-side protocol surface in our pairing service (TO1 redirect + TO2 onboarding, CBOR wire format, ownership-voucher storage, owner key management)
- reuse the existing open-source FDO ecosystem for manufacturing/rendezvous — never reimplement it.

**Acceptance:** a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end and lands as a provisioned realm device.

**Docs are a first-class deliverable:** upstream ships almost no public FDO documentation (code-level only); ours must exceed that — operator guide covering the full chain (manufacturing → rendezvous → owner) with Astrate-specific setup.

## Next step

Investigation phase before any implementation: read upstream's TO2 handling end-to-end (`libs/astarte_fdo*", pairing controllers/plugs, `fdo/queries`), inventory the exact endpoints, credential/voucher schema and key material required, and measure what our pairing service lacks. Output feeds the v4.0 scope decision in `.mule/milestones.md`.
--

$ cd /root/astrate-mule && timeout 120 gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6 --jq '{ahead:.ahead_by, behind:.behind_by, total:.total_commits, files:(.files|length), commits:[.commits[]|.sha[0:7]+" "+(.commit.message|split("\n")[0])]}' 2>&1 | head -60
{"ahead":27,"behind":0,"commits":["d1c39fb doc: improve FDO architecture documentation (#2153)","162062e docs: improve FDO user guide","be04b70 feat(fdo): add ownership voucher deletion (#2156)","682ab62 build(e2e): fix docker build (#2185)","7714c6b chore(fdo): allow configuring an IP address for FDO URL composition (#2177)","7790c74 fix(dup): prevent AMQP consumer process to accumulate uncollected bin… (#2194)","37d8fd2 chore!(dup): remove fullsweep_after configuration (#2198)","f8da2b9 ci: use upstream astarte-device-fdo-rust (#2199)","e0b7161 feat(fdo): Allow specifying device_id on ownership voucher upload (#2196)","5ce1fd0 fix(dashboard): Add unknow status when device messgaes aren't proporly consumed (#2200)","d52f3c6 feat(dashboard): add device id and initial introspection to voucher upload  (#2202)","d05a3d7 chore(fdo): move voucher and session in astarte keyspace (#2168)","16befa3 fix(dup): ensure RPC server is always available to clients (#2203)","a5a0114 fix(fdo): Accept single entry x5chain certificates on voucher upload (#2206)","76b2153 chore: prepare v1.3.4 (#2208)","57eb1fe chore: forward-port release-1.3 into release-1.4","39d1acf fix(fdo): reject invalid ownership vouchers before storing them (#2209)","a6ef02a fix(dup): ensure RPC server restarts after name conflicts (#2205)","895821a chore(fdo): improve validations in queries to upload/delete ownership vouchers","b3bf731 fix(pairing): ensure users cannot delete ownership vouchers belonging to other realms","7cb1348 fix(fdo): Ensure ownership voucher cleanup on device deletion (#2211)","8d47574 chore: update horde (#2224)","f5f337e chore: forward port release-1.3","be403b4 feat(fdo): Add `POST /v1/:realm_name/fdo/ownership_vouchers/:guid/to0` to run TO0 again for an ownership voucher, refreshing how long the rendezvous server keeps serving its registration. The resulting expiry is reported when listing the ownership vouchers of a realm (#2220)","6c8d4b8 refactor(fdo)!: remove realm_name requirement from FDO onboarding routes (#2217)","2c2d5aa build(docker): add https routes to astarte services (#2223)","20df0d3 chore: prepare v1.4.0-rc.6 (#2227)"],"files":153,"total":27}

$ cd /root/astrate-mule && timeout 180 gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6 --jq '.files[] | .status + " " + .filename' 2>&1 | grep -Ei 'fdo|ownership|device_removal|voucher' 
added apps/astarte_housekeeping/priv/migrations/astarte/0006_create_ownership_vouchers_table.sql
added apps/astarte_housekeeping/priv/migrations/astarte/0007_add_expiry_to_ownership_vouchers.sql
added apps/astarte_housekeeping/priv/migrations/realm/0021_add_device_id_to_ownership_voucher.sql
added apps/astarte_housekeeping/priv/migrations/realm/0022_delete_vouchers.sql
added apps/astarte_housekeeping/priv/migrations/realm/0023_drop_ownership_voucher_table.sql
added apps/astarte_housekeeping/priv/migrations/realm/0025_add_fdo_guid_to_devices.sql
modified apps/astarte_pairing/lib/astarte_pairing_web/api_spec/schemas/ownership_voucher.ex
modified apps/astarte_pairing/lib/astarte_pairing_web/controllers/fdo_onboarding_controller.ex
modified apps/astarte_pairing/lib/astarte_pairing_web/controllers/ownership_voucher_controller.ex
modified apps/astarte_pairing/lib/astarte_pairing_web/plug/fdo_session.ex
modified apps/astarte_pairing/lib/astarte_pairing_web/views/ownership_voucher_view.ex
modified apps/astarte_pairing/test/astarte_pairing_web/controllers/fdo_onboarding_controller_test.exs
modified apps/astarte_pairing/test/astarte_pairing_web/controllers/ownership_voucher_controller_test.exs
modified apps/astarte_pairing/test/support/cases/fdo_session.ex
modified apps/astarte_pairing/test/support/helpers/fdo.ex
modified apps/astarte_realm_management/lib/astarte_realm_management/device_removal/core.ex
modified apps/astarte_realm_management/lib/astarte_realm_management/device_removal/device_remover.ex
modified apps/astarte_realm_management/lib/astarte_realm_management/device_removal/queries.ex
modified apps/astarte_realm_management/lib/astarte_realm_management/device_removal/scheduler.ex
modified apps/astarte_realm_management/test/astarte_realm_management/device_removal/device_remover_test.exs
modified apps/astarte_realm_management/test/astarte_realm_management/device_removal/scheduler_test.exs
added astarte-dashboard/cypress/e2e/fdo_vouchers_page.cy.js
modified astarte-dashboard/src/FdoVoucherPage.tsx
modified astarte-dashboard/src/FdoVouchersPage.tsx
modified astarte-dashboard/src/hooks/useFdo.ts
modified libs/astarte_data_access/lib/astarte_data_access/fdo/ownership_voucher.ex
modified libs/astarte_data_access/lib/astarte_data_access/fdo/queries.ex
modified libs/astarte_data_access/lib/astarte_data_access/fdo/to2_session.ex
removed libs/astarte_data_access/test/fdo/ownership_voucher/ownership_voucher_test.exs
modified libs/astarte_data_access/test/fdo/queries_test.exs
modified libs/astarte_fdo/config/test.exs
modified libs/astarte_fdo/lib/config.ex
added libs/astarte_fdo/lib/config/base_url_host.ex
modified libs/astarte_fdo/lib/owner_onboarding.ex
modified libs/astarte_fdo/lib/owner_onboarding/session.ex
modified libs/astarte_fdo/lib/ownership_voucher/load_request.ex
modified libs/astarte_fdo/lib/ownership_voucher/ownership_voucher.ex
modified libs/astarte_fdo/lib/service_info.ex
modified libs/astarte_fdo/lib/to0.ex
modified libs/astarte_fdo/mix.lock
added libs/astarte_fdo/test/astarte_fdo/config/base_url_host_test.exs
modified libs/astarte_fdo/test/astarte_fdo/config_test.exs
modified libs/astarte_fdo/test/astarte_fdo/onboarding/done_test.exs
modified libs/astarte_fdo/test/astarte_fdo/onboarding/prove_device_test.exs
modified libs/astarte_fdo/test/astarte_fdo/owner_onboarding/owner_onboarding_test.exs
modified libs/astarte_fdo/test/astarte_fdo/owner_onboarding/session_test.exs
modified libs/astarte_fdo/test/astarte_fdo/ownership_voucher/load_request_test.exs
modified libs/astarte_fdo/test/astarte_fdo/ownership_voucher_test.exs
modified libs/astarte_fdo/test/astarte_fdo/rendezvous/core_test.exs
modified libs/astarte_fdo/test/astarte_fdo/rendezvous/rv_to2_addr_test.exs
modified libs/astarte_fdo/test/astarte_fdo/service_info_test.exs
modified libs/astarte_fdo/test/astarte_fdo/to0_test.exs
modified libs/astarte_fdo/test/support/cases/fdo_session.ex
modified libs/astarte_fdo/test/support/helpers/database.ex
modified libs/astarte_fdo/test/support/helpers/fdo.ex
modified libs/astarte_fdo/test/test_helper.exs
modified libs/astarte_fdo_core/lib/rendezvous/rv_to2_addr.ex
modified libs/astarte_fdo_core/mix.exs
modified libs/astarte_fdo_core/mix.lock
added libs/astarte_fdo_core/test/astarte_fdo_core/fdo/rendezvous/rv_to2_addr_test.exs

$ cd /tmp/opencode && timeout 180 gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6 --jq '.files[] | select(.filename|test("router.ex|mix.exs$|\\.md$")) | .status + " " + .filename' 2>&1
modified CHANGELOG.md
modified apps/astarte_appengine_api/mix.exs
modified apps/astarte_data_updater_plant/mix.exs
modified apps/astarte_housekeeping/mix.exs
modified apps/astarte_pairing/lib/astarte_pairing_web/router.ex
modified apps/astarte_pairing/mix.exs
modified apps/astarte_realm_management/mix.exs
modified apps/astarte_trigger_engine/mix.exs
modified doc/mix.exs
modified doc/pages/architecture/050-pairing_mechanism.md
modified doc/pages/user/035-register_device.md
modified libs/astarte_config/mix.exs
modified libs/astarte_data_access/mix.exs
modified libs/astarte_events/mix.exs
modified libs/astarte_fdo_core/mix.exs
modified libs/astarte_generators/mix.exs
modified libs/astarte_rpc/mix.exs
modified libs/astarte_secrets/mix.exs
modified tools/astarte_device_fleet_simulator/mix.exs
modified tools/astarte_e2e/mix.exs
modified tools/astarte_export/mix.exs
modified tools/astarte_import/mix.exs

$ cd /tmp/opencode && timeout 180 gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6 --jq '.files[] | select(.filename=="apps/astarte_pairing/lib/astarte_pairing_web/router.ex") | .patch' 2>&1
@@ -36,8 +36,6 @@ defmodule Astarte.PairingWeb.Router do
   pipeline :fdo do
     plug :accepts, ["cbor"]
     plug :put_view, Astarte.PairingWeb.FDOView
-    plug Astarte.PairingWeb.Plug.LogRealm
-    plug Astarte.PairingWeb.Plug.VerifyRealmExists
     plug Astarte.PairingWeb.Plug.SetupFDO
   end
 
@@ -49,7 +47,8 @@ defmodule Astarte.PairingWeb.Router do
     plug Astarte.PairingWeb.Plug.DecryptAndVerify
   end
 
-  scope "/v1/:realm_name/fdo/101", Astarte.PairingWeb do
+  # URL mandated by FDO specification
+  scope "/fdo/101", Astarte.PairingWeb do
     pipe_through :fdo
 
     post "/msg/60", FDOOnboardingController, :hello_device
@@ -101,6 +100,8 @@ defmodule Astarte.PairingWeb.Router do
       post "/owner_keys_for_voucher", OwnershipVoucherController, :owner_keys_for_voucher
       get "/ownership_vouchers", OwnershipVoucherController, :list_ownership_vouchers
       post "/ownership_vouchers", OwnershipVoucherController, :register
+      delete "/ownership_vouchers/:guid", OwnershipVoucherController, :delete_ownership_voucher
+      post "/ownership_vouchers/:guid/to0", OwnershipVoucherController, :run_to0
     end
   end

$ cd /tmp/opencode && timeout 180 gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6 --jq '.files[] | select(.filename|test("controllers/ownership_voucher_controller.ex$|views/ownership_voucher_view.ex$|api_spec/schemas/ownership_voucher.ex$")) | "=== " + .filename + "\n" + (.patch // "NO PATCH")' 2>&1
=== apps/astarte_pairing/lib/astarte_pairing_web/api_spec/schemas/ownership_voucher.ex
@@ -24,12 +24,16 @@ defmodule Astarte.PairingWeb.ApiSpec.Schemas.OwnershipVoucher do
     @moduledoc false
     require OpenApiSpex
 
+    alias Astarte.PairingWeb.ApiSpec.Schemas.Agent.DeviceRegistrationRequest
+
     OpenApiSpex.schema(%{
       type: :object,
       properties: %{
         data: %Schema{
           type: :object,
           properties: %{
+            hw_id: DeviceRegistrationRequest.hw_id(),
+            initial_introspection: DeviceRegistrationRequest.initial_introspection(),
             ownership_voucher: %Schema{
               type: :string,
               description:
@@ -61,7 +65,7 @@ defmodule Astarte.PairingWeb.ApiSpec.Schemas.OwnershipVoucher do
               description: "Optional PEM-encoded replacement public key."
             }
           },
-          required: [:ownership_voucher, :key_name, :key_algorithm]
+          required: [:hw_id, :ownership_voucher, :key_name, :key_algorithm]
         }
       },
       required: [:data]
@@ -99,11 +103,41 @@ defmodule Astarte.PairingWeb.ApiSpec.Schemas.OwnershipVoucher do
                 type: :string,
                 nullable: true,
                 description: "The PEM-encoded output ownership voucher, if any."
+              },
+              expiry: %Schema{
+                type: :string,
+                format: :"date-time",
+                nullable: true,
+                description:
+                  "The instant at which the rendezvous server stops serving the registration " <>
+                    "made during TO0, if the voucher was ever registered."
               }
             }
           }
         }
       }
     })
   end
+
+  defmodule TO0Response do
+    @moduledoc false
+    require OpenApiSpex
+
+    OpenApiSpex.schema(%{
+      type: :object,
+      properties: %{
+        data: %Schema{
+          type: :object,
+          properties: %{
+            expiry: %Schema{
+              type: :string,
+              format: :"date-time",
+              description:
+                "The instant at which the rendezvous server stops serving the new registration."
+            }
+          }
+        }
+      }
+    })
+  end
 end
=== apps/astarte_pairing/lib/astarte_pairing_web/controllers/ownership_voucher_controller.ex
@@ -22,14 +22,21 @@ defmodule Astarte.PairingWeb.OwnershipVoucherController do
 
   alias Astarte.FDO.OwnershipVoucher
   alias Astarte.FDO.OwnershipVoucher.LoadRequest
-  alias Astarte.FDO.TO0
+  alias Astarte.Pairing.Engine
+  alias Astarte.PairingWeb.ApiSpec.Schemas.Errors
   alias Astarte.PairingWeb.ApiSpec.Schemas.OwnershipVoucher, as: OVApiSpec
   alias Astarte.PairingWeb.OwnershipVoucherView
   alias Astarte.Secrets.Core, as: SecretsCore
+  alias OpenApiSpex.MediaType
+  alias OpenApiSpex.Response
   alias OpenApiSpex.Schema
 
   action_fallback Astarte.PairingWeb.FallbackController
 
+  # `action_fallback` renders from the conn as it entered the action, so this
+  # has to be a plug for the error view to see the GUID
+  plug :assign_fdo_guid when action in [:run_to0]
+
   tags ["fdo"]
 
   operation :register,
@@ -51,6 +58,19 @@ defmodule Astarte.PairingWeb.OwnershipVoucherController do
       ok: {"Ownership voucher registered successfully", nil, nil},
       bad_request: {"Invalid request body", nil, nil},
       unauthorized: {"Unauthorized", nil, nil},
+      forbidden: %Response{
+        description: "Forbidden or Authorization path not matched",
+        content: %{
+          "application/json" => %MediaType{
+            schema: %Schema{
+              oneOf: [
+                Errors.ForbiddenResponse,
+                Errors.AuthorizationPathNotMatchedResponse
+              ]
+            }
+          }
+        }
+      },
       not_found: {"Realm not found", nil, nil},
       internal_server_error: {"Internal server error", nil, nil}
     ]
@@ -75,6 +95,64 @@ defmodule Astarte.PairingWeb.OwnershipVoucherController do
       internal_server_error: {"Internal server error", nil, nil}
     ]
 
+  operation :delete_ownership_voucher,
+    summary: "Delete an ownership voucher",
+    description:
+      "Deletes an ownership voucher. Revokes the corresponding registration on the FDO " <>
+        "rendezvous server first; if the revocation fails, the voucher is not deleted.",
+    operation_id: "deleteOwnershipVoucher",
+    security: [%{"JWT" => []}],
+    parameters: [
+      realm_name: [
+        in: :path,
+        description: "Name of the realm.",
+        type: :string,
+        required: true
+      ],
+      guid: [
+        in: :path,
+        description: "GUID of the ownership voucher to delete.",
+        type: :string,
+        required: true
+      ]
+    ],
+    responses: [
+      no_content: {"Ownership voucher deleted successfully", nil, nil},
+      unauthorized: {"Unauthorized", nil, nil},
+      not_found: {"Ownership voucher not found", nil, nil},
+      internal_server_error: {"Internal server error", nil, nil}
+    ]
+
+  operation :run_to0,
+    summary: "Re-run TO0 for an ownership voucher",
+    description:
+      "Registers the ownership voucher on the FDO rendezvous server again and refreshes its " <>
+        "expiry. Only vouchers whose device has not completed Device Onboard yet can be " <>
+        "re-registered.",
+    operation_id: "runOwnershipVoucherTO0",
+    security: [%{"JWT" => []}],
+    parameters: [
+      realm_name: [
+        in: :path,
+        description: "Name of the realm.",
+        type: :string,
+        required: true
+      ],
+      guid: [
+        in: :path,
+        description: "GUID of the ownership voucher to re-register.",
+        type: :string,
+        required: true
+      ]
+    ],
+    responses: [
+      ok: {"TO0 completed successfully", "application/json", OVApiSpec.TO0Response},
+      unauthorized: {"Unauthorized", nil, nil},
+      not_found: {"Ownership voucher not found", nil, nil},
+      conflict: {"Device Onboard has already completed for the ownership voucher", nil, nil},
+      internal_server_error: {"Internal server error", nil, nil}
+    ]
+
   operation :owner_keys_for_voucher,
     summary: "List owner keys compatible with an ownership voucher",
     description:
@@ -139,22 +217,19 @@ defmodule Astarte.PairingWeb.OwnershipVoucherController do
     with {:ok, req} <-
            LoadRequest.changeset(%LoadRequest{}, Map.put(data, "realm_name", realm_name))
            |> Ecto.Changeset.apply_action(:insert),
-         :ok <-
-           OwnershipVoucher.save_voucher(realm_name, %{
-             voucher_data: req.cbor_ownership_voucher,
-             guid: req.device_guid,
-             key_name: req.key_name,
-             key_algorithm: req.key_algorithm,
-             replacement_guid: req.replacement_guid,
-             replacement_rendezvous_info: req.decoded_replacement_rendezvous_info,
-             replacement_public_key: req.decoded_replacement_public_key
-           }),
-         :ok <-
-           TO0.claim_ownership_voucher(
-             realm_name,
+         {:ok, expiry} <-
+           OwnershipVoucher.claim_on_rendezvous(
+             req.device_guid,
              req.decoded_ownership_voucher,
              req.extracted_owner_key
-           ) do
+           ),
+         :ok <- LoadRequest.store_voucher(req, expiry),
+         opts = [
+           initial_introspection: req.initial_introspection,
+           with_credentials?: false,
+           fdo_guid: req.device_guid
+         ],
+         {:ok, nil} <- Engine.register_device(realm_name, req.hw_id, opts) do
       json(conn, %{
         data: %{
           public_key: req.extracted_owner_key.public_pem,
@@ -175,6 +250,32 @@ defmodule Astarte.PairingWeb.OwnershipVoucherController do
     end
   end
 
+  @doc """
+  Deletes an ownership voucher.
+
+  Returns `204 No Content` on success, `404 Not Found` if the GUID is unknown (in the realm).
+  """
+  def delete_ownership_voucher(conn, %{"realm_name" => realm_name, "guid" => guid_str}) do
+    with {:ok, guid} <- decode_guid(guid_str),
+         :ok <- OwnershipVoucher.delete(realm_name, guid) do
+      send_resp(conn, :no_content, "")
+    end
+  end
+
+  @doc """
+  Re-runs TO0 with the rendezvous server for an ownership voucher, refreshing
+  its expiry.
+
+  Returns `200 OK` with the new expiry, `404 Not Found` if the GUID is unknown
+  (in the realm), `409 Conflict` if Device Onboard already completed for it.
+  """
+  def run_to0(conn, %{"realm_name" => realm_name, "guid" => guid_str}) do
+    with {:ok, guid} <- decode_guid(guid_str),
+         {:ok, expiry} <- OwnershipVoucher.run_to0(realm_name, guid) do
+      json(conn, %{data: %{expiry: expiry}})
+    end
+  end
+
   @doc """
   Returns the list of registered owner keys that are compatible with the
   given ownership voucher.
@@ -196,4 +297,15 @@ defmodule Astarte.PairingWeb.OwnershipVoucherController do
 
   defp ensure_ownership_voucher_parameter(_params),
     do: {:error, :missing_ownership_voucher}
+
+  defp assign_fdo_guid(conn, _opts) do
+    assign(conn, :fdo_guid, conn.params["guid"])
+  end
+
+  defp decode_guid(guid_str) do
+    case Ecto.UUID.dump(guid_str) do
+      {:ok, _} = ok -> ok
+      :error -> {:error, :not_found}
+    end
+  end
 end
=== apps/astarte_pairing/lib/astarte_pairing_web/views/ownership_voucher_view.ex
@@ -34,7 +34,8 @@ defmodule Astarte.PairingWeb.OwnershipVoucherView do
       replacement_guid: output_guid,
       voucher_data: input_voucher,
       output_voucher: output_voucher,
-      status: status
+      status: status,
+      expiry: expiry
     } = voucher
 
     guid = render_one(guid, OwnershipVoucherView, "guid.json", as: :guid)
@@ -47,7 +48,8 @@ defmodule Astarte.PairingWeb.OwnershipVoucherView do
       status: status,
       output_guid: output_guid,
       input_voucher: input_voucher,
-      output_voucher: output_voucher
+      output_voucher: output_voucher,
+      expiry: expiry
     }
   end

$ cd /tmp/opencode && timeout 180 gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6 --jq '.files[] | select(.filename|test("device\\.ex$|load_request.ex$|fdo/queries.ex$|fdo/ownership_voucher.ex$")) | "=== " + .filename + "\n" + (.patch // "NO PATCH")' 2>&1
=== apps/astarte_pairing/test/support/cases/device.ex
@@ -83,6 +83,8 @@ defmodule Astarte.Cases.Device do
         credentials_secret: confirmed_credentials_secret
     }
 
+    device_without_credentials = %{base_device(interfaces) | credentials_secret: nil}
+
     for interface <- interfaces, do: insert_interface_cleanly(realm_name, interface)
 
     interface_descriptors = get_interface_descriptors(realm_name, interfaces)
@@ -111,6 +113,7 @@ defmodule Astarte.Cases.Device do
 
     insert_device_cleanly(realm_name, confirmed_device, interfaces, confirmed_credentials_secret)
     insert_device_cleanly(realm_name, inhibited_device, interfaces, inhibited_credentials_secret)
+    insert_device_cleanly(realm_name, device_without_credentials, interfaces, nil)
 
     %{
       interfaces: interfaces,
@@ -119,6 +122,7 @@ defmodule Astarte.Cases.Device do
       unconfirmed_device2: unconfirmed_device2,
       confirmed_device: confirmed_device,
       inhibited_device: inhibited_device,
+      device_without_credentials: device_without_credentials,
       unregistered_device: unregistered_device,
       interface_descriptors: interface_descriptors,
       endpoints: endpoints
=== apps/astarte_pairing/test/support/helpers/device.ex
@@ -68,14 +68,14 @@ defmodule Astarte.Helpers.Device do
     introspection_minor = interfaces |> Map.new(&{&1.name, &1.minor_version})
     interfaces_bytes = Map.fetch!(device, :interfaces_bytes)
     interfaces_msgs = Map.fetch!(device, :interfaces_msgs)
-    secret = CredentialsSecret.hash(secret)
+    hash = if secret, do: CredentialsSecret.hash(secret)
 
     device_db_params = %{
       introspection: introspection,
       introspection_minor: introspection_minor,
       exchanged_bytes_by_interface: interfaces_bytes,
       exchanged_msgs_by_interface: interfaces_msgs,
-      credentials_secret: secret
+      credentials_secret: hash
     }
 
     device_db = struct(Device, Map.merge(device, device_db_params))
=== libs/astarte_data_access/lib/astarte_data_access/device.ex
@@ -23,13 +23,15 @@ defmodule Astarte.DataAccess.Device do
   require Logger
   alias Astarte.Core.Device, as: DeviceCore
   alias Astarte.DataAccess.Consistency
-  alias Astarte.DataAccess.Device.UnconfirmedDevice
   alias Astarte.DataAccess.Devices.Device
   alias Astarte.DataAccess.Realms.Realm
   alias Astarte.DataAccess.Repo
+  alias Ecto.Changeset
 
   import Ecto.Query
 
+  @ten_minutes_in_seconds 600
+
   @spec interface_version(String.t(), DeviceCore.device_id(), String.t()) ::
           {:ok, integer} | {:error, atom}
   def interface_version(realm, device_id, interface_name) do
@@ -58,6 +60,29 @@ defmodule Astarte.DataAccess.Device do
     {:error, :device_not_found}
   end
 
+  @spec add_unconfirmed_credentials(String.t(), DeviceCore.device_id(), String.t()) ::
+          :ok | {:error, term()}
+  def add_unconfirmed_credentials(realm_name, device_id, credentials_secret) do
+    opts = [
+      prefix: Realm.keyspace_name(realm_name),
+      consistency: Consistency.device_info(:write),
+      ttl: @ten_minutes_in_seconds,
+      allow_insert: false,
+      allow_stale: true
+    ]
+
+    device_result =
+      %Device{device_id: device_id}
+      |> Changeset.change(%{
+        credentials_secret: credentials_secret
+      })
+      |> Repo.update(opts)
+
+    with {:ok, _} <- device_result do
+      :ok
+    end
+  end
+
   def register(realm_name, device_id, extended_id, credentials_secret, opts \\ []) do
     case fetch(realm_name, device_id) do
       {:error, :device_not_found} ->
@@ -108,25 +133,25 @@ defmodule Astarte.DataAccess.Device do
       |> Keyword.get(:initial_introspection, [])
       |> build_initial_introspection_maps()
 
+    fdo_guid = Keyword.get(opts, :fdo_guid)
+
     keyspace_name = Realm.keyspace_name(realm_name)
     consistency = Consistency.device_info(:write)
-    unconfirmed? = Keyword.get(opts, :unconfirmed, false)
     opts = [prefix: keyspace_name, consistency: consistency]
 
-    with :ok <- register_unconfirmed_device(unconfirmed?, device_id, opts) do
-      %Device{
-        device_id: device_id,
-        first_registration: registration_timestamp,
-        credentials_secret: credentials_secret,
-        inhibit_credentials_request: false,
-        protocol_revision: 0,
-        total_received_bytes: 0,
-        total_received_msgs: 0,
-        introspection: introspection,
-        introspection_minor: introspection_minor
-      }
-      |> Repo.insert(opts)
-    end
+    %Device{
+      device_id: device_id,
+      first_registration: registration_timestamp,
+      credentials_secret: credentials_secret,
+      inhibit_credentials_request: false,
+      protocol_revision: 0,
+      total_received_bytes: 0,
+      total_received_msgs: 0,
+      introspection: introspection,
+      introspection_minor: introspection_minor,
+      fdo_guid: fdo_guid
+    }
+    |> Repo.insert(opts)
   end
 
   defp do_register_unconfirmed_device(
@@ -140,22 +165,28 @@ defmodule Astarte.DataAccess.Device do
       |> Keyword.get(:initial_introspection, [])
       |> build_initial_introspection_maps()
 
+    changes = %{
+      credentials_secret: credentials_secret,
+      inhibit_credentials_request: false,
+      protocol_revision: 0,
+      introspection: introspection,
+      introspection_minor: introspection_minor
+    }
+
+    # Don't overwrite an already assigned fdo_guid when none is given
+    changes =
+      case Keyword.get(opts, :fdo_guid) do
+        nil -> changes
+        fdo_guid -> Map.put(changes, :fdo_guid, fdo_guid)
+      end
+
     keyspace_name = Realm.keyspace_name(realm_name)
     consistency = Consistency.device_info(:write)
-    unconfirmed? = Keyword.get(opts, :unconfirmed, false)
     opts = [prefix: keyspace_name, consistency: consistency]
 
-    with :ok <- register_unconfirmed_device(unconfirmed?, device.device_id, opts) do
-      device
-      |> Ecto.Changeset.change(%{
-        credentials_secret: credentials_secret,
-        inhibit_credentials_request: false,
-        protocol_revision: 0,
-        introspection: introspection,
-        introspection_minor: introspection_minor
-      })
-      |> Repo.insert(opts)
-    end
+    device
+    |> Ecto.Changeset.change(changes)
+    |> Repo.insert(opts)
   end
 
   defp build_initial_introspection_maps(initial_introspection) do
@@ -170,38 +201,31 @@ defmodule Astarte.DataAccess.Device do
     end)
   end
 
-  defp register_unconfirmed_device(false, _device_id, _opts), do: :ok
-
-  defp register_unconfirmed_device(true, device_id, opts) do
-    opts =
-      opts
-      |> Keyword.put(:overwrite, false)
-      |> Keyword.put(:allow_stale, true)
-
-    result =
-      %UnconfirmedDevice{device_id: device_id, created_at: DateTime.utc_now()}
-      |> Repo.insert(opts)
-
-    with {:ok, _} <- result do
-      :ok
-    end
-  end
-
   @spec confirm(String.t(), DeviceCore.device_id()) ::
-          {:ok, Device.t()} | {:error, :device_not_found}
+          {:ok, Device.t()} | {:error, :device_not_found} | {:error, :expired_credentials}
   def confirm(realm_name, device_id) do
-    unconfirmed_device = %UnconfirmedDevice{device_id: device_id}
     keyspace = Realm.keyspace_name(realm_name)
-    delete_opts = [prefix: keyspace, consistency: Consistency.device_info(:write)]
+
+    insert_opts = [prefix: keyspace, consistency: Consistency.device_info(:write)]
 
     fetch_opts = [
       prefix: keyspace,
       consistency: Consistency.device_info(:read),
       error: :device_not_found
     ]
 
-    Repo.delete!(unconfirmed_device, delete_opts)
-    Repo.fetch(Device, device_id, fetch_opts)
+    with {:ok, device} <- Repo.fetch(Device, device_id, fetch_opts),
+         :ok <- check_not_expired_credentials(device) do
+      # Removes TTL from credentials_secret
+      Repo.insert(device, insert_opts)
+    end
+  end
+
+  defp check_not_expired_credentials(device) do
+    case device.credentials_secret do
+      nil -> {:error, :expired_credentials}
+      _ -> :ok
+    end
   end
 
   def fetch(realm_name, device_id) do
@@ -215,4 +239,141 @@ defmodule Astarte.DataAccess.Device do
       error: :device_not_found
     )
   end
+
+  def fetch_with_unconfirmed_status(realm_name, device_id) do
+    keyspace_name = Realm.keyspace_name(realm_name)
+
+    opts = [
+      prefix: keyspace_name,
+      consistency: Consistency.device_info(:read),
+      error: :device_not_found
+    ]
+
+    device_query =
+      from d in Device,
+        select: %{
+          device_id: d.device_id,
+          aliases: d.aliases,
+          attributes: d.attributes,
+          cert_aki: d.cert_aki,
+          cert_serial: d.cert_serial,
+          connected: d.connected,
+          credentials_secret: d.credentials_secret,
+          exchanged_bytes_by_interface: d.exchanged_bytes_by_interface,
+          exchanged_msgs_by_interface: d.exchanged_msgs_by_interface,
+          first_credentials_request: d.first_credentials_request,
+          first_registration: d.first_registration,
+          groups: d.groups,
+          fdo_guid: d.fdo_guid,
+          inhibit_credentials_request: d.inhibit_credentials_request,
+          introspection: d.introspection,
+          introspection_minor: d.introspection_minor,
+          last_connection: d.last_connection,
+          last_credentials_request_ip: d.last_credentials_request_ip,
+          last_disconnection: d.last_disconnection,
+          last_seen_ip: d.last_seen_ip,
+          old_introspection: d.old_introspection,
+          capabilities: d.capabilities,
+          pending_empty_cache: d.pending_empty_cache,
+          protocol_revision: d.protocol_revision,
+          total_received_bytes: d.total_received_bytes,
+          total_received_msgs: d.total_received_msgs,
+          credentials_secret_ttl: fragment("TTL(?)", d.credentials_secret)
+        }
+
+    with {:ok, device_params} <- Repo.fetch(device_query, device_id, opts) do
+      %{
+        device_id: device_id,
+        aliases: aliases,
+        attributes: attributes,
+        cert_aki: cert_aki,
+        cert_serial: cert_serial,
+        connected: connected,
+        credentials_secret: credentials_secret,
+        exchanged_bytes_by_interface: exchanged_bytes_by_interface,
+        exchanged_msgs_by_interface: exchanged_msgs_by_interface,
+        first_credentials_request: first_credentials_request,
+        first_registration: first_registration,
+        groups: groups,
+        fdo_guid: fdo_guid,
+        inhibit_credentials_request: inhibit_credentials_request,
+        introspection: introspection,
+        introspection_minor: introspection_minor,
+        last_connection: last_connection,
+        last_credentials_request_ip: last_credentials_request_ip,
+        last_disconnection: last_disconnection,
+        last_seen_ip: last_seen_ip,
+        old_introspection: old_introspection,
+        capabilities: capabilities,
+        pending_empty_cache: pending_empty_cache,
+        protocol_revision: protocol_revision,
+        total_received_bytes: total_received_bytes,
+        total_received_msgs: total_received_msgs,
+        credentials_secret_ttl: credentials_secret_ttl
+      } = device_params
+
+      confirmation_status =
+        case credentials_secret_ttl do
+          nil -> :confirmed
+          _ -> :unconfirmed
+        end
+
+      device =
+        %Device{
+          device_id: device_id,
+          aliases: aliases,
+          attributes: attributes,
+          cert_aki: cert_aki,
+          cert_serial: cert_serial,
+          connected: connected,
+          credentials_secret: credentials_secret,
+          exchanged_bytes_by_interface: exchanged_bytes_by_interface,
+          exchanged_msgs_by_interface: exchanged_msgs_by_interface,
+          first_credentials_request: first_credentials_request,
+          first_registration: first_registration,
+          groups: groups,
+          fdo_guid: fdo_guid,
+          inhibit_credentials_request: inhibit_credentials_request,
+          introspection: introspection,
+          introspection_minor: introspection_minor,
+          last_connection: last_connection,
+          last_credentials_request_ip: last_credentials_request_ip,
+          last_disconnection: last_disconnection,
+          last_seen_ip: last_seen_ip,
+          old_introspection: old_introspection,
+          capabilities: capabilities,
+          pending_empty_cache: pending_empty_cache,
+          protocol_revision: protocol_revision,
+          total_received_bytes: total_received_bytes,
+          total_received_msgs: total_received_msgs,
+          confirmation_status: confirmation_status
+        }
+
+      {:ok, device}
+    end
+  end
+
+  def unregister(realm_name, device_id) do
+    with {:ok, device} <- fetch(realm_name, device_id),
+         {:ok, _device} <- do_unregister_device(realm_name, device) do
+      :ok
+    else
+      {:error, reason} ->
+        Logger.warning("Unregister error: #{inspect(reason)}")
+        {:error, reason}
+    end
+  end
+
+  defp do_unregister_device(realm_name, %Device{} = device) do
+    keyspace_name = Realm.keyspace_name(realm_name)
+
+    consistency = Consistency.device_info(:write)
+
+    device
+    |> Changeset.change(
+      first_credentials_request: nil,
+      credentials_secret: nil
+    )
+    |> Repo.update(prefix: keyspace_name, consistency: consistency)
+  end
 end
=== libs/astarte_data_access/lib/astarte_data_access/device/unconfirmed_device.ex
@@ -1,32 +0,0 @@
-#
-# This file is part of Astarte.
-#
-# Copyright 2026 SECO Mind Srl
-#
-# Licensed under the Apache License, Version 2.0 (the "License");
-# you may not use this file except in compliance with the License.
-# You may obtain a copy of the License at
-#
-#    http://www.apache.org/licenses/LICENSE-2.0
-#
-# Unless required by applicable law or agreed to in writing, software
-# distributed under the License is distributed on an "AS IS" BASIS,
-# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-# See the License for the specific language governing permissions and
-# limitations under the License.
-#
-
-defmodule Astarte.DataAccess.Device.UnconfirmedDevice do
-  @moduledoc """
-  This module defines the Ecto schema for the `deletion_in_progress` table.
-  """
-  use TypedEctoSchema
-
-  alias Astarte.DataAccess.DateTime, as: DateTimeMs
-  alias Astarte.DataAccess.UUID
-
-  @primary_key {:device_id, UUID, autogenerate: false}
-  typed_schema "unconfirmed_devices" do
-    field :created_at, DateTimeMs
-  end
-end
=== libs/astarte_data_access/lib/astarte_data_access/devices/device.ex
@@ -48,6 +48,7 @@ defmodule Astarte.DataAccess.Devices.Device do
     field :first_credentials_request, DateTimeMs
     field :first_registration, DateTimeMs
     field :groups, Exandra.Map, key: :string, value: UUID
+    field :fdo_guid, :binary
     field :inhibit_credentials_request, :boolean
     field :introspection, Exandra.Map, key: :string, value: :integer
     field :introspection_minor, Exandra.Map, key: :string, value: :integer
@@ -68,5 +69,7 @@ defmodule Astarte.DataAccess.Devices.Device do
     field :protocol_revision, :integer
     field :total_received_bytes, :integer
     field :total_received_msgs, :integer
+
+    field :confirmation_status, Ecto.Enum, values: [:confirmed, :unconfirmed], virtual: true
   end
 end
=== libs/astarte_data_access/lib/astarte_data_access/fdo/ownership_voucher.ex
@@ -22,17 +22,17 @@ defmodule Astarte.DataAccess.FDO.OwnershipVoucher do
   """
   use TypedEctoSchema
 
-  import Ecto.Changeset
-
+  alias Astarte.DataAccess.DateTime, as: DateTimeMs
   alias Astarte.DataAccess.FDO.CBOR.Encoded, as: CBOREncoded
-  alias Astarte.DataAccess.FDO.OwnershipVoucher
   alias Astarte.FDO.Core.OwnershipVoucher.RendezvousInfo
   alias Astarte.FDO.Core.PublicKey
 
   @primary_key false
   typed_schema "ownership_vouchers" do
     field :guid, Astarte.DataAccess.UUID, primary_key: true
-    field :status, Ecto.Enum, values: [created: 0, claimed: 1]
+    field :device_id, Astarte.DataAccess.UUID
+    field :status, Ecto.Enum, values: [created: 0, claimed: 1], default: :created
+    field :realm, :string
     field :voucher_data, :binary
     field :output_voucher, :binary
     field :user_id, :binary
@@ -41,20 +41,6 @@ defmodule Astarte.DataAccess.FDO.OwnershipVoucher do
     field :replacement_guid, :binary
     field :replacement_rendezvous_info, CBOREncoded, using: RendezvousInfo
     field :replacement_public_key, CBOREncoded, using: PublicKey
-  end
-
-  @doc false
-  def changeset(%OwnershipVoucher{} = record, attrs) do
-    record
-    |> cast(attrs, [
-      :key_name,
-      :voucher_data,
-      :guid,
-      :key_algorithm,
-      :replacement_guid,
-      :replacement_rendezvous_info,
-      :replacement_public_key
-    ])
-    |> validate_required([:key_name, :key_algorithm, :voucher_data, :guid])
+    field :expiry, DateTimeMs
   end
 end
=== libs/astarte_data_access/lib/astarte_data_access/fdo/queries.ex
@@ -31,29 +31,42 @@ defmodule Astarte.DataAccess.FDO.Queries do
 
   require Logger
 
-  def get_ownership_voucher(realm_name, guid) do
-    keyspace_name = Realm.keyspace_name(realm_name)
+  def fetch_ownership_voucher(guid) do
+    keyspace_name = Realm.astarte_keyspace_name()
+    opts = [consistency: Consistency.domain_model(:read), prefix: keyspace_name]
+
+    Repo.fetch(OwnershipVoucher, guid, opts)
+  end
+
+  def fetch_device_id_and_ownership_voucher_and_realm(guid) do
+    keyspace_name = Realm.astarte_keyspace_name()
 
     query =
       from o in OwnershipVoucher,
         prefix: ^keyspace_name,
-        select: o.voucher_data
+        select: {o.realm, o.device_id, o.voucher_data}
 
     consistency = Consistency.domain_model(:read)
 
     Repo.fetch(query, guid, consistency: consistency)
   end
 
   def list_ownership_vouchers(realm_name) do
-    keyspace = Realm.keyspace_name(realm_name)
+    keyspace_name = Realm.astarte_keyspace_name()
+
+    query =
+      from o in OwnershipVoucher,
+        hints: ["ALLOW FILTERING"],
+        prefix: ^keyspace_name,
+        where: o.realm == ^realm_name
+
     consistency = Consistency.domain_model(:read)
-    opts = [consistency: consistency, prefix: keyspace]
 
-    Repo.fetch_all(OwnershipVoucher, opts)
+    Repo.fetch_all(query, consistency: consistency)
   end
 
-  def get_owner_key_params(realm_name, guid) do
-    keyspace_name = Realm.keyspace_name(realm_name)
+  def get_owner_key_params(guid) do
+    keyspace_name = Realm.astarte_keyspace_name()
 
     query =
       from OwnershipVoucher,
@@ -68,65 +81,105 @@ defmodule Astarte.DataAccess.FDO.Queries do
     end
   end
 
-  def get_replacement_data(realm_name, guid) do
-    keyspace = Realm.keyspace_name(realm_name)
+  def get_replacement_data(guid) do
+    keyspace = Realm.astarte_keyspace_name()
 
     fields = [:replacement_guid, :replacement_rendezvous_info, :replacement_public_key]
 
     query =
       from OwnershipVoucher,
+        prefix: ^keyspace,
         select: ^fields
 
     consistency = Consistency.domain_model(:read)
-    opts = [consistency: consistency, prefix: keyspace]
+    opts = [consistency: consistency]
 
     with {:ok, data} <- Repo.fetch(query, guid, opts) do
       result = Map.take(data, fields)
       {:ok, result}
     end
   end
 
-  def create_ownership_voucher(
-        realm_name,
-        attrs
-      ) do
-    keyspace_name = Realm.keyspace_name(realm_name)
+  @spec create_ownership_voucher(OwnershipVoucher.t()) :: :ok | {:error, term()}
+  def create_ownership_voucher(ownership_voucher) do
+    keyspace_name = Realm.astarte_keyspace_name()
+
+    # explicitly prevent updates to an already existing entry for the same GUID;
+    # expect a "stale entry error" if no rows were changed due to GUID already present
+    opts = [
+      prefix: keyspace_name,
+      consistency: Consistency.device_info(:write),
+      overwrite: false,
+      stale_error_field: :guid
+    ]
 
-    opts = [prefix: keyspace_name, consistency: Consistency.device_info(:write)]
+    case Repo.insert(ownership_voucher, opts) do
+      {:ok, _} ->
+        :ok
 
-    %OwnershipVoucher{status: :created}
-    |> OwnershipVoucher.changeset(attrs)
-    |> Repo.insert(opts)
+      {:error, %Ecto.Changeset{errors: [guid: {_, [stale: true]}]}} ->
+        {:error, :duplicated_voucher_guid}
+
+      {:error, changeset} ->
+        {:error, changeset}
+    end
   end
 
-  def delete_ownership_voucher(realm_name, guid) do
-    keyspace = Realm.keyspace_name(realm_name)
+  def delete_ownership_voucher(guid) do
+    keyspace = Realm.astarte_keyspace_name()
+    consistency = Consistency.device_info(:write)
+    opts = [prefix: keyspace, consistency: consistency]
 
     %OwnershipVoucher{
       guid: guid
     }
-    |> Repo.delete(prefix: keyspace)
+    |> Repo.delete(opts)
+    |> case do
+      {:ok, _voucher} -> :ok
+      error -> error
+    end
   end
 
-  def mark_voucher_as_claimed(realm_name, guid) do
-    keyspace = Realm.keyspace_name(realm_name)
+  @doc """
+  Marks an ownership voucher as claimed by its device.
+
+  The registration made on the rendezvous server during TO0 is consumed by TO2,
+  so the expiry is cleared along with the status change.
+  """
+  def mark_voucher_as_claimed(guid) do
+    keyspace = Realm.astarte_keyspace_name()
     consistency = Consistency.device_info(:write)
     opts = [prefix: keyspace, consistency: consistency]
 
     result =
       %OwnershipVoucher{guid: guid}
       |> Ecto.Changeset.change(status: :claimed)
+      # `change/2` skips values that already match the (empty) struct, so the
+      # expiry has to be forced in to actually be written as null
+      |> Ecto.Changeset.force_change(:expiry, nil)
+      |> Repo.update(opts)
+
+    with {:ok, _} <- result, do: :ok
+  end
+
+  def update_voucher_expiry(guid, expiry) do
+    keyspace = Realm.astarte_keyspace_name()
+    consistency = Consistency.device_info(:write)
+    opts = [prefix: keyspace, consistency: consistency]
+
+    result =
+      %OwnershipVoucher{guid: guid}
+      |> Ecto.Changeset.change(expiry: expiry)
       |> Repo.update(opts)
 
     with {:ok, _} <- result, do: :ok
   end
 
   def add_output_voucher(
-        realm_name,
         guid,
         new_voucher
       ) do
-    keyspace = Realm.keyspace_name(realm_name)
+    keyspace = Realm.astarte_keyspace_name()
     consistency = Consistency.device_info(:write)
     opts = [prefix: keyspace, consistency: consistency]
 
@@ -138,74 +191,72 @@ defmodule Astarte.DataAccess.FDO.Queries do
     with {:ok, _} <- result, do: :ok
   end
 
-  def store_session(realm_name, guid, session) do
-    keyspace = Realm.keyspace_name(realm_name)
+  def store_session(session) do
+    keyspace = Realm.astarte_keyspace_name()
     consistency = Consistency.device_info(:write)
     opts = [prefix: keyspace, consistency: consistency]
 
-    session = %{session | guid: guid}
-
     with {:ok, _} <- Repo.insert(session, opts) do
       :ok
     end
   end
 
-  def delete_session(realm_name, guid) do
-    keyspace = Realm.keyspace_name(realm_name)
+  def delete_session(guid) do
+    keyspace = Realm.astarte_keyspace_name()
     consistency = Consistency.device_info(:write)
     opts = [prefix: keyspace, consistency: consistency]
 
     Repo.delete(%TO2Session{guid: guid}, opts)
     :ok
   end
 
-  def add_session_max_owner_service_info_size(realm_name, guid, size) do
+  def add_session_max_owner_service_info_size(guid, size) do
     updates = [max_owner_service_info_size: size]
-    update_session(realm_name, guid, updates)
+    update_session(guid, updates)
   end
 
-  def add_session_secret(realm_name, guid, secret) do
+  def add_session_secret(guid, secret) do
     updates = [secret: secret]
-    update_session(realm_name, guid, updates)
+    update_session(guid, updates)
   end
 
-  def add_session_keys(realm_name, guid, sevk, svk, sek) do
+  def add_session_keys(guid, sevk, svk, sek) do
     updates = [sevk: sevk, svk: svk, sek: sek]
-    update_session(realm_name, guid, updates)
+    update_session(guid, updates)
   end
 
-  def session_add_setup_dv_nonce(realm_name, guid, setup_dv_nonce) do
+  def session_add_setup_dv_nonce(guid, setup_dv_nonce) do
     updates = [setup_dv_nonce: setup_dv_nonce]
-    update_session(realm_name, guid, updates)
+    update_session(guid, updates)
   end
 
-  def session_update_device_id(realm_name, guid, device_id) do
+  def session_update_device_id(guid, device_id) do
     updates = [device_id: device_id]
-    update_session(realm_name, guid, updates)
+    update_session(guid, updates)
   end
 
-  def session_add_device_service_info(realm_name, guid, service_info) do
+  def session_add_device_service_info(guid, service_info) do
     updates = [device_service_info: service_info]
-    update_session(realm_name, guid, updates)
+    update_session(guid, updates)
   end
 
-  def session_add_owner_service_info(realm_name, guid, owner_service_info) do
+  def session_add_owner_service_info(guid, owner_service_info) do
     updates = [owner_service_info: owner_service_info]
-    update_session(realm_name, guid, updates)
+    update_session(guid, updates)
   end
 
-  def session_update_last_chunk_sent(realm_name, guid, last_chunk) do
+  def session_update_last_chunk_sent(guid, last_chunk) do
     updates = [last_chunk_sent: last_chunk]
-    update_session(realm_name, guid, updates)
+    update_session(guid, updates)
   end
 
-  def session_add_replacement_hmac(realm_name, guid, hmac) do
+  def session_add_replacement_hmac(guid, hmac) do
     updates = [replacement_hmac: hmac]
-    update_session(realm_name, guid, updates)
+    update_session(guid, updates)
   end
 
-  defp update_session(realm_name, guid, updates) do
-    keyspace = Realm.keyspace_name(realm_name)
+  defp update_session(guid, updates) do
+    keyspace = Realm.astarte_keyspace_name()
     consistency = Consistency.device_info(:write)
     opts = [prefix: keyspace, consistency: consistency]
 
@@ -218,8 +269,8 @@ defmodule Astarte.DataAccess.FDO.Queries do
     end
   end
 
-  def fetch_session(realm_name, guid) do
-    keyspace = Realm.keyspace_name(realm_name)
+  def fetch_session(guid) do
+    keyspace = Realm.astarte_keyspace_name()
     consistency = Consistency.device_info(:read)
     opts = [prefix: keyspace, consistency: consistency]
     Repo.fetch(TO2Session, guid, opts)
=== libs/astarte_fdo/lib/ownership_voucher/load_request.ex
@@ -23,6 +23,10 @@ defmodule Astarte.FDO.OwnershipVoucher.LoadRequest do
 
   use TypedEctoSchema
 
+  alias Astarte.Core.Device
+  alias Astarte.DataAccess.Device, as: DeviceQueries
+  alias Astarte.DataAccess.FDO.OwnershipVoucher, as: OwnershipVoucherStruct
+  alias Astarte.DataAccess.FDO.Queries
   alias Astarte.FDO.Core.OwnershipVoucher
   alias Astarte.FDO.Core.OwnershipVoucher.Core, as: OVCore
   alias Astarte.FDO.Core.OwnershipVoucher.RendezvousInfo
@@ -37,6 +41,9 @@ defmodule Astarte.FDO.OwnershipVoucher.LoadRequest do
   import Ecto.Changeset
 
   typed_embedded_schema do
+    field :hw_id, :string
+    field :initial_introspection, :map, default: %{}
+    field :device_id, Astarte.DataAccess.UUID, virtual: true
     field :ownership_voucher, :string
     field :realm_name, :string
     field :key_name, :string
@@ -61,6 +68,8 @@ defmodule Astarte.FDO.OwnershipVoucher.LoadRequest do
   def changeset(%LoadRequest{} = request, params) do
     request
     |> cast(params, [
+      :hw_id,
+      :initial_introspection,
       :ownership_voucher,
       :realm_name,
       :key_name,
@@ -69,8 +78,12 @@ defmodule Astarte.FDO.OwnershipVoucher.LoadRequest do
       :replacement_public_key,
       :replacement_guid
     ])
-    |> validate_required([:ownership_voucher, :realm_name, :key_name, :key_algorithm])
+    |> validate_required([:hw_id, :ownership_voucher, :realm_name, :key_name, :key_algorithm])
+    |> put_device_id(:hw_id, :device_id)
     |> put_device_guid()
+    |> ensure_voucher_not_already_claimed(:device_guid, :ownership_voucher)
+    |> ensure_device_does_not_exist(:realm_name, :device_id)
+    |> validate_change(:initial_introspection, &validate_introspection/2)
     |> validate_key_algorithm_compatible()
     |> fetch_owner_key()
     |> verify_owner_key_matches()
@@ -90,6 +103,74 @@ defmodule Astarte.FDO.OwnershipVoucher.LoadRequest do
     |> decode_replacement_fields()
   end
 
+  defp ensure_voucher_not_already_claimed(%{valid?: false} = changeset, _, _), do: changeset
+
+  defp ensure_voucher_not_already_claimed(changeset, guid, voucher) do
+    # SAFETY: we only call this on valid vouchers
+    guid = fetch_field!(changeset, guid)
+
+    case Queries.fetch_ownership_voucher(guid) do
+      {:error, :not_found} -> changeset
+      {:ok, _old_voucher} -> add_error(changeset, voucher, "guid has already been claimed")
+    end
+  end
+
+  defp ensure_device_does_not_exist(%{valid?: false} = changeset, _, _), do: changeset
+
+  defp ensure_device_does_not_exist(changeset, realm_name_field, device_id_field) do
+    # SAFETY: we only call this on valid vouchers
+    realm_name = fetch_field!(changeset, realm_name_field)
+    device_id = fetch_field!(changeset, device_id_field)
+
+    case DeviceQueries.fetch(realm_name, device_id) do
+      {:error, :device_not_found} -> changeset
+      {:ok, _device} -> add_error(changeset, device_id_field, "already exists")
+    end
+  end
+
+  def put_device_id(%{valid?: false} = changeset, _encoded_field, _device_id_field), do: changeset
+
+  def put_device_id(changeset, encoded_field, device_id_field) do
+    case validate_hw_id_change(changeset, encoded_field) do
+      {:ok, device_id} -> changeset |> put_change(device_id_field, device_id)
+      {:error, changeset} -> changeset
+    end
+  end
+
+  def validate_hw_id(changeset, _field) when not changeset.valid?, do: changeset
+
+  def validate_hw_id(changeset, field) do
+    case validate_hw_id_change(changeset, field) do
+      {:ok, _device_id} -> changeset
+      {:error, changeset} -> changeset
+    end
+  end
+
+  defp validate_hw_id_change(changeset, field) do
+    with {:ok, hw_id} <- fetch_change(changeset, field),
+         {:ok, device_id} <- Device.decode_device_id(hw_id, allow_extended_id: true) do
+      {:ok, device_id}
+    else
+      _ ->
+        {:error, add_error(changeset, field, "is not a valid base64 encoded 128 bits id")}
+    end
+  end
+
+  def validate_introspection(field, introspection) when is_map(introspection) do
+    Enum.reduce(introspection, [], fn
+      {interface_name, %{"major" => major, "minor" => minor}}, acc
+      when is_integer(major) and is_integer(minor) ->
+        if major < 0 or minor < 0 do
+          [{field, "has negative versions in interface #{interface_name}"} | acc]
+        else
+          acc
+        end
+
+      {interface_name, _}, acc ->
+        [{field, "has invalid format for interface #{interface_name}"} | acc]
+    end)
+  end
+
   defp validate_replacement_rendezvous_info(changeset) do
     validate_change(changeset, :replacement_rendezvous_info, fn :replacement_rendezvous_info,
                                                                 b64_string ->
@@ -274,6 +355,15 @@ defmodule Astarte.FDO.OwnershipVoucher.LoadRequest do
     end
   end
 
+  defp public_keys_match?(:x5chain, cert_der, pem) when is_binary(cert_der) do
+    with {:ok, cert_spki_der} <- spki_der_from_cert(cert_der),
+         [{_, pem_spki_der, :not_encrypted}] <- :public_key.pem_decode(pem) do
+      cert_spki_der == pem_spki_der
+    else
+      _ -> false
+    end
+  end
+
   # :cosekey — decode CBOR map, decode PEM via OTP, compare key material.
   defp public_keys_match?(:cosekey, cosekey_cbor, pem) do
     with {:ok, cose_map, ""} <- CBOR.decode(cosekey_cbor),
@@ -359,4 +449,41 @@ defmodule Astarte.FDO.OwnershipVoucher.LoadRequest do
         :error
     end
   end
+
+  @doc """
+  Persists the ownership voucher described by `load_request`.
+
+  `expiry` is the instant at which the rendezvous server stops serving the
+  registration made during TO0.
+  """
+  @spec store_voucher(t(), DateTime.t()) :: :ok | {:error, term()}
+  def store_voucher(load_request, expiry) do
+    %LoadRequest{
+      realm_name: realm_name,
+      device_guid: guid,
+      device_id: device_id,
+      cbor_ownership_voucher: cbor_ownership_voucher,
+      key_name: key_name,
+      key_algorithm: key_algorithm,
+      replacement_guid: replacement_guid,
+      decoded_replacement_rendezvous_info: decoded_replacement_rendezvous_info,
+      decoded_replacement_public_key: decoded_replacement_public_key
+    } = load_request
+
+    ownership_voucher = %OwnershipVoucherStruct{
+      guid: guid,
+      realm: realm_name,
+      device_id: device_id,
+      status: :created,
+      voucher_data: cbor_ownership_voucher,
+      key_name: key_name,
+      key_algorithm: key_algorithm,
+      replacement_guid: replacement_guid,
+      replacement_rendezvous_info: decoded_replacement_rendezvous_info,
+      replacement_public_key: decoded_replacement_public_key,
+      expiry: expiry
+    }
+
+    Queries.create_ownership_voucher(ownership_voucher)
+  end
 end

$ cd /tmp/opencode && timeout 180 gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6 --jq '.files[] | select(.filename|test("CHANGELOG.md$|fallback_controller.ex$")) | "=== " + .filename + "\n" + (.patch // "NO PATCH")' 2>&1 | head -120
=== CHANGELOG.md
@@ -7,6 +7,33 @@ and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.
 
 ## Unreleased
 
+## [1.4.0-rc.6] - 2026-09-28
+
+Includes changes from v1.3.4
+
+### Added
+
+- [fdo] Add `POST /v1/:realm_name/fdo/ownership_vouchers/:guid/to0` to run TO0 again for an
+  ownership voucher, refreshing how long the rendezvous server keeps serving its registration.
+  The resulting expiry is reported when listing the ownership vouchers of a realm.
+- [fdo] Add ownership voucher deletion
+
+### Changed
+
+- [fdo] Rename `ASTARTE_BASE_URL_DOMAIN` to `ASTARTE_BASE_URL_HOST`. Its value may now be
+  either a domain name or an IP address
+- [fdo] Deleting a device also deletes its ownership voucher.
+- [fdo] Allow specifying device_id on ownership voucher upload
+- [fdo] The device is now immediately registered on ownership voucher upload
+- [astarte_data_updater_plant] Improve RPC server reliability
+
+### Fixed
+
+- [dashboard] Add unknown status when device messages aren't properly consumed
+- [fdo] Accept single entry x5chain certificates on voucher upload
+- [fdo] Invalid vouchers are no longer stored on the database
+- [fdo] Ensure users cannot delete ownership vouchers belonging to other realms
+
 ## [1.4.0-rc.5] - 2026-08-20
 
 ### Fixed
@@ -96,6 +123,15 @@ and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.
 
 - [astarte_data_updater_plant] Use mississippi consumer for data updater processes
 
+## [1.3.4] - 2026-09-17
+
+### Fixed
+
+- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
+  uncollected binaries over time by forcing more frequent full sweep garbage collections
+  on them.
+- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.
+
 ## [1.3.3] - 2026-08-07
 
 ### Fixed
=== apps/astarte_pairing/lib/astarte_pairing_web/controllers/fallback_controller.ex
@@ -102,6 +102,32 @@ defmodule Astarte.PairingWeb.FallbackController do
     |> render(:missing_ownership_voucher)
   end
 
+  # Deleting an ownership voucher requires first revoking its registration on
+  # the FDO rendezvous server and the voucher is not deleted if that fails
+  def call(conn, {:error, :rendezvous_revocation_failed}) do
+    conn
+    |> put_status(:internal_server_error)
+    |> put_view(Astarte.PairingWeb.ErrorView)
+    |> render(:"500")
+  end
+
+  # Both uploading a voucher and re-running TO0 register it on the FDO
+  # rendezvous server, and neither can succeed if that fails
+  def call(conn, {:error, :rendezvous_registration_failed}) do
+    conn
+    |> put_status(:internal_server_error)
+    |> put_view(Astarte.PairingWeb.ErrorView)
+    |> render(:"500")
+  end
+
+  # Once Device Onboard completes, the voucher is not registered on the rendezvous anymore
+  def call(conn, {:error, :device_already_onboarded}) do
+    conn
+    |> put_status(:conflict)
+    |> put_view(Astarte.PairingWeb.ErrorView)
+    |> render(:device_already_onboarded)
+  end
+
   # This is called when no JWT token is present
   def auth_error(conn, {:unauthenticated, :unauthenticated}, _opts) do
     conn

$ cd /tmp/opencode && timeout 180 gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6 --jq '.files[] | select(.filename|test("fdo_onboarding_controller.ex$|plug/fdo_session.ex$")) | "=== " + .filename + "\n" + (.patch // "NO PATCH")' 2>&1 | head -250
=== apps/astarte_pairing/lib/astarte_pairing_web/controllers/fdo_onboarding_controller.ex
@@ -20,7 +20,6 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
   use Astarte.PairingWeb, :controller
   use OpenApiSpex.ControllerSpecs
 
-  alias Astarte.Core.Device
   alias Astarte.FDO.Core.OwnerOnboarding.DeviceServiceInfo
   alias Astarte.FDO.Core.OwnerOnboarding.DeviceServiceInfoReady
   alias Astarte.FDO.OwnerOnboarding
@@ -41,14 +40,6 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
     summary: "FDO TO2 Hello Device",
     operation_id: "FDOHelloDevice",
     description: "Sets up new owner for proof of ownership.",
-    parameters: [
-      realm_name: [
-        in: :path,
-        description: "Name of the realm the device belongs to.",
-        type: :string,
-        required: true
-      ]
-    ],
     request_body: {
       "TO2 Hello Device",
       "application/cbor",
@@ -102,14 +93,6 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
     summary: "FDO TO2 GetOVNextEntry",
     operation_id: "FDOGetOVNextEntry",
     description: "Requests the next Ownership Voucher Entry.",
-    parameters: [
-      realm_name: [
-        in: :path,
-        description: "Name of the realm the device belongs to.",
-        type: :string,
-        required: true
-      ]
-    ],
     request_body: {
       "TO2 GetOVNextEntry",
       "application/cbor",
@@ -153,14 +136,6 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
     summary: "FDO TO2 ProveDevice",
     operation_id: "FDOProveDevice",
     description: "Proves the provenance of the Device to the new owner.",
-    parameters: [
-      realm_name: [
-        in: :path,
-        description: "Name of the realm the device belongs to.",
-        type: :string,
-        required: true
-      ]
-    ],
     request_body: {
       "TO2 ProveDeviceRequest",
       "application/cbor",
@@ -207,14 +182,6 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
     summary: "FDO TO2 Done",
     operation_id: "FDODone",
     description: "Indicates successful completion of the Transfer of Ownership.",
-    parameters: [
-      realm_name: [
-        in: :path,
-        description: "Name of the realm the device belongs to.",
-        type: :string,
-        required: true
-      ]
-    ],
     request_body: {
       "TO2 Done",
       "application/cbor",
@@ -229,7 +196,7 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
            format: :binary,
            description:
              "This message provides an opportunity for a final ACK after the Owner has invoked the System Info block to
-                          establish agent-to-server communications between the Device and its final Owner.",
+                           establish agent-to-server communications between the Device and its final Owner.",
            example: [
              "NonceTO2SetupDv"
            ]
@@ -259,15 +226,7 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
     operation_id: "FDOServiceInfoReady",
     description:
       "This message signals a state change between the authentication phase of the protocol and the provisioning
-                  phase (ServiceInfo) negotiation.",
-    parameters: [
-      realm_name: [
-        in: :path,
-        description: "Name of the realm the device belongs to.",
-        type: :string,
-        required: true
-      ]
-    ],
+                   phase (ServiceInfo) negotiation.",
     request_body: {
       "TO2 Service Info Ready",
       "application/cbor",
@@ -282,7 +241,7 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
            format: :binary,
            description:
              "This message responds to TO2.DeviceServiceInfoReady and indicates that the Owner Onboarding Service is
-                        ready to start ServiceInfo.",
+                         ready to start ServiceInfo.",
            example: [
              "maxDeviceServiceInfoSz"
            ]
@@ -312,15 +271,7 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
     operation_id: "FDOServiceInfoEnd",
     description:
       "Sends as many Device to Owner ServiceInfo entries as will conveniently fit into a message, based on protocol
-            and Device constraints.",
-    parameters: [
-      realm_name: [
-        in: :path,
-        description: "Name of the realm the device belongs to.",
-        type: :string,
-        required: true
-      ]
-    ],
+             and Device constraints.",
     request_body: {
       "CBOR Service Info End",
       "application/cbor",
@@ -335,7 +286,7 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
            format: :binary,
            description:
              "Sends as many Owner to Device ServiceInfo entries as will conveniently fit into a message, based on protocol
-                          and implementation constraints. This message is part of a loop with TO2.DeviceServiceInfo.",
+                           and implementation constraints. This message is part of a loop with TO2.DeviceServiceInfo.",
            example: [
              "IsMoreServiceInfo,",
              "IsDone",
@@ -362,39 +313,35 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
     ]
 
   def hello_device(conn, _params) do
-    realm_name = Map.fetch!(conn.params, "realm_name")
     cbor_hello_device = conn.assigns.cbor_body
 
     with {:ok, token, response_msg} <-
-           OwnerOnboarding.hello_device(realm_name, cbor_hello_device) do
+           OwnerOnboarding.hello_device(cbor_hello_device) do
       conn
       |> put_resp_header("authorization", token)
       |> render("default.cbor", %{cbor_response: response_msg})
     end
   end
 
   def ov_next_entry(conn, _params) do
-    realm_name = Map.fetch!(conn.params, "realm_name")
     cbor_body = conn.assigns.cbor_body
-
-    guid = conn.assigns.to2_session.guid
+    session = conn.assigns.to2_session
 
     with {:ok, response} <-
-           OwnerOnboarding.ov_next_entry(cbor_body, realm_name, guid) do
+           OwnerOnboarding.ov_next_entry(cbor_body, session.guid) do
       conn
       |> render("default.cbor", %{cbor_response: response})
     end
   end
 
   def prove_device(conn, _params) do
-    realm_name = Map.fetch!(conn.params, "realm_name")
     cbor_body = conn.assigns.cbor_body
+    session = conn.assigns.to2_session
 
     with {:ok, session, response} <-
            OwnerOnboarding.prove_device(
-             realm_name,
              cbor_body,
-             conn.assigns.to2_session
+             session
            ) do
       conn
       |> assign(:to2_session, session)
@@ -403,23 +350,22 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
   end
 
   def done(conn, _params) do
-    realm_name = Map.fetch!(conn.params, "realm_name")
     to2_session = conn.assigns.to2_session
 
-    with {:ok, response_msg} <- OwnerOnboarding.done(realm_name, to2_session, conn.assigns.body) do
+    with {:ok, response_msg} <-
+           OwnerOnboarding.done(to2_session, conn.assigns.body) do
       conn
       |> render("secure.cbor", %{cbor_response: response_msg})
     end
   end
 
   def service_info_start(conn, _params) do
-    realm_name = Map.fetch!(conn.params, "realm_name")
+    session = conn.assigns.to2_session
 
     with {:ok, device_service_info_ready} <- DeviceServiceInfoReady.decode(conn.assigns.body),
          {:ok, session, response} <-
            OwnerOnboarding.build_owner_service_info_ready(
-             realm_name,
-             conn.assigns.to2_session,
+             session,
              device_service_info_ready
            ) do
       conn
@@ -429,13 +375,12 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
   end
 
   def service_info_end(conn, _params) do
-    realm_name = Map.fetch!(conn.params, "realm_name")
+    session = conn.assigns.to2_session
 
     with {:ok, device_service_info} <- DeviceServiceInfo.decode(conn.assigns.body),
          {:ok, response} <-
            build_owner_service_info_response(
-             realm_name,
-             conn.assigns.to2_session,
+             session,
              device_service_info
            ) do
       conn
@@ -444,53 +389,34 @@ defmodule Astarte.PairingWeb.FDOOnboardingController do
   end
 
   defp build_owner_service_info_response(
-         realm_name,
          session,
          %DeviceServiceInfo{is_more_service_info: true} = device_service_info
        ) do
-    ServiceInfo.build_owner_service_info(realm_name, session, device_service_info)
+    ServiceInfo.build_owner_service_info(session, device_service_info)
   end
 
   defp build_owner_service_info_response(
-         realm_name,
          session,
          %DeviceServiceInfo{is_more_service_info: false, service_info: service_info} =
            device_service_info
        )
        when map_size(service_info) == 0 do
-    ServiceInfo.build_owner_service_info(realm_name, session, device_service_info)
+    ServiceInfo.build_owner_service_info(session, device_service_info)
   end
 
   defp build_owner_service_info_response(
-         realm_name,
          session,
          %DeviceServiceInfo{is_more_service_info: false, service_info: service_info}
        ) do
-    with {:ok, session} <-
-           Session.add_device_service_info(session, realm_name, service_info) do
-      decoded_service_info =

$ cd /tmp/opencode && timeout 180 gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6 --jq '.files[] | select(.filename|test("plug/fdo_session.ex$|device_removal/core.ex$|fdo/lib/config.ex$")) | "=== " + .filename + "\n" + (.patch // "NO PATCH")' 2>&1 | head -200
=== apps/astarte_pairing/lib/astarte_pairing_web/plug/fdo_session.ex
@@ -35,12 +35,10 @@ defmodule Astarte.PairingWeb.Plug.FDOSession do
   end
 
   def call(conn, _opts) do
-    realm_name = Map.fetch!(conn.path_params, "realm_name")
-
     with [token] <- get_req_header(conn, "authorization"),
          {:ok, guid, token_nonce} <-
            SessionToken.verify(token),
-         {:ok, session} <- Session.fetch(realm_name, guid),
+         {:ok, session} <- Session.fetch(guid),
          :ok <- verify_nonce(session.nonce, token_nonce) do
       conn
       |> put_resp_header("authorization", token)
=== apps/astarte_realm_management/lib/astarte_realm_management/device_removal/core.ex
@@ -24,10 +24,13 @@ defmodule Astarte.RealmManagement.DeviceRemoval.Core do
   """
   alias Astarte.Core.CQLUtils
   alias Astarte.Core.InterfaceDescriptor
+  alias Astarte.DataAccess.FDO.Queries, as: FDOQueries
   alias Astarte.DataAccess.Interface
   alias Astarte.RealmManagement.DeviceRemoval.Queries, as: DeviceRemovalQueries
   alias Astarte.RealmManagement.TriggersHandler
 
+  require Logger
+
   @doc """
   Deletes individual datastreams for a device in a realm.
   """
@@ -179,6 +182,24 @@ defmodule Astarte.RealmManagement.DeviceRemoval.Core do
     DeviceRemovalQueries.delete_kv_store_entry!(realm_name, group_name, key)
   end
 
+  @doc """
+  Deletes the ownership voucher bound to a device, if any.
+
+  N.B.: only the database row is removed.
+  """
+  def delete_ownership_voucher!(_realm_name, nil), do: :ok
+
+  def delete_ownership_voucher!(realm_name, guid) do
+    _ =
+      Logger.info(
+        "Deleting ownership voucher without revoking its rendezvous registration",
+        realm: realm_name,
+        tag: "fdo_voucher_deleted_without_revocation"
+      )
+
+    :ok = FDOQueries.delete_ownership_voucher(guid)
+  end
+
   @doc """
   Removes a device from the database.
   """
@@ -194,7 +215,6 @@ defmodule Astarte.RealmManagement.DeviceRemoval.Core do
       TriggersHandler.device_deletion_finished(realm_name, device_id, groups)
     end
 
-    DeviceRemovalQueries.remove_device_from_unconfirmed_devices!(realm_name, device_id)
     DeviceRemovalQueries.remove_device_from_deletion_in_progress!(realm_name, device_id)
   end
 end
=== libs/astarte_fdo/lib/config.ex
@@ -23,6 +23,7 @@ defmodule Astarte.FDO.Config do
 
   use Astarte.Config
 
+  alias Astarte.FDO.Config.BaseURLHost
   alias Astarte.FDO.Config.BaseURLProtocol
 
   @envdoc "The port the ingress is listening on, used for FDO authentication mechanism"
@@ -37,10 +38,13 @@ defmodule Astarte.FDO.Config do
     type: BaseURLProtocol,
     required: true
 
-  @envdoc "The astarte base domain, used for FDO authentication mechanism"
-  app_env :base_url_domain, :astarte_fdo, :base_url_domain,
-    os_env: "ASTARTE_BASE_URL_DOMAIN",
-    type: :binary,
+  @envdoc """
+  The astarte base URL host, used for FDO authentication mechanism. It can be
+  either a domain name or an IP address.
+  """
+  app_env :base_url_host, :astarte_fdo, :base_url_host,
+    os_env: "ASTARTE_BASE_URL_HOST",
+    type: BaseURLHost,
     required: true
 
   url_env :rendezvous, :astarte_fdo, :rendezvous, env_app: "PAIRING_FDO", default_port: 8041
@@ -57,9 +61,9 @@ defmodule Astarte.FDO.Config do
 
   def base_url! do
     protocol = __MODULE__.base_url_protocol!()
-    domain = __MODULE__.base_url_domain!()
+    host = __MODULE__.base_url_host!()
     port = __MODULE__.base_url_port!()
 
-    "#{protocol}://#{domain}:#{port}"
+    "#{protocol}://#{host}:#{port}"
   end
 end

$ cd /root/astrate-mule && sed -n '74,84p' internal/pairing/http.go && echo "=== svc" && sed -n '362,376p' internal/pairing/service.go
// Mount registers the pairing routes on mux. Paths are wire-frozen
// (docs/DESIGN.md §4.4): they are exactly what the official SDKs and
// astartectl call.
func (a *API) Mount(mux *http.ServeMux) {
	// Realm-scoped health (#71, upstream 1.3+): unauthenticated by design —
	// FDO-flow devices probe it before they hold credentials. No rate
	// limiter either; it is one cheap query.
	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
	mux.Handle("POST /pairing/v1/{realm}/agent/devices",
		a.requireAgent(http.HandlerFunc(a.handleRegister)))
	mux.Handle("DELETE /pairing/v1/{realm}/agent/devices/{deviceID}",
=== svc
	res.Until = until
	return res, nil
}

// Health reports whether the realm exists and the database behind it is
// healthy (GET /pairing/v1/{realm}/health, upstream 1.3+; unauthenticated by
// design — FDO-flow devices probe it before they hold credentials).
func (s *Service) Health(ctx context.Context, realm string) error {
	if _, err := s.st.GetRealmByName(ctx, realm); err != nil {
		return err // store.ErrNotFound → 404
	}
	return s.st.Health(ctx) // nil → 200; anything else → 503
}

// authenticateDevice resolves the realm and device and bcrypt-compares the

$ cd /tmp/opencode && for p in doc/pages/user/035-register_device.md doc/pages/architecture/050-pairing_mechanism.md; do echo "=== $p"; timeout 90 gh api "repos/astarte-platform/astarte/contents/$p?ref=v1.4.0-rc.6" --jq '.content' | base64 -d | grep -E '^#{1,4} ' | head -25; done
=== doc/pages/user/035-register_device.md
# Registering a Device
## Registration Agent
### On Board Agent
### 3rd Party Agent
#### "Local" or "Plant" Agents
#### "Remote" Agents
## FIDO Device Onboard
### Owner Keys
### Ownership Vouchers
## Credentials Secret Lifecycle
## Unregistering a device
=== doc/pages/architecture/050-pairing_mechanism.md
# Pairing Mechanism
## Authentication flow
## Credentials Secret vs. Transport Credentials
### Credentials Secret storage recommendations
## Using SSL Certificates as Transport Credentials
### Mutual SSL Authentication Flow
### Certificate Authority
### Certificate flow
## Agents
### On Board Agent
## FIDO Device Onboarding
### Transfer Ownership protocols
#### TO0 Protocol
#### TO2 protocol
### Owner Keys management
## Transport responsibility
## Pairing facilities

$ cd /root/astrate-mule && rg -n -i 'fdo|fido device|ownership voucher' docs/ .mule/milestones.md docs/api/astrate_pairing_api.yaml 2>/dev/null | head -30
.mule/milestones.md:124:(FDO #78, lua_map #87, pipeline DSL #86, keyAgreement #51) stay out unless Giulio
docs/COMPATIBILITY.md:361:    answers 404 and an unhealthy database 503, and FDO-flow devices can probe
docs/api/astarte_pairing_api.yaml:33:        and unrate-limited — FDO-flow devices probe it before they hold
docs/site/api/astarte_pairing_api.yaml:33:        and unrate-limited — FDO-flow devices probe it before they hold

$ cd /root/astrate-mule && sed -n '112,132p' .mule/milestones.md
work produced by the 2026-08-22 upstream-parity investigation instead:

- Source document: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 /
  v1.4.0-rc.5 diffed service-by-service against our emulated 1.2.2 level).
- Backlog: GitHub issues **#47–#89** (`upstream-parity` label), plus the parked
  **#90** APICompatVersion decision (milestone-4.0 label — pull into this milestone's
  final phase, after the 1.3-surface audit it requires).

Shape of the milestone: work through #47–#89 in the research doc's priority order
(retention ceiling → alias/group surfaces → query formats → wire capabilities →
validation/error-code fidelity), escalating every issue marked "decide" to
`.mule/for-giulio.md` instead of choosing unilaterally. Deliberately parked items
(FDO #78, lua_map #87, pipeline DSL #86, keyAgreement #51) stay out unless Giulio
pulls them in.

Status: not started. First recipe job: triage #47–#89 into an ordered plan (which are
audits vs features vs decisions), file sub-issues where work splits, escalate the
"decide" set in one batch.

---

$ cd /tmp/opencode && timeout 180 gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6 --jq '.files[] | select(.filename|test("migrations/")) | .status + " " + .filename' 2>&1; echo "=== release"; timeout 90 gh api repos/astarte-platform/astarte/releases/tags/v1.4.0-rc.6 --jq '{tag:.tag_name,published:.published_at,prerelease:.prerelease,target:.target_commitish}' 2>&1
added apps/astarte_housekeeping/priv/migrations/astarte/0004_create_session_key_type.sql
added apps/astarte_housekeeping/priv/migrations/astarte/0005_create_to2_sessions_table.sql
added apps/astarte_housekeeping/priv/migrations/astarte/0006_create_ownership_vouchers_table.sql
added apps/astarte_housekeeping/priv/migrations/astarte/0007_add_expiry_to_ownership_vouchers.sql
added apps/astarte_housekeeping/priv/migrations/realm/0020_drop_unconfirmed_devices_table.sql
added apps/astarte_housekeeping/priv/migrations/realm/0021_add_device_id_to_ownership_voucher.sql
added apps/astarte_housekeeping/priv/migrations/realm/0022_delete_vouchers.sql
added apps/astarte_housekeeping/priv/migrations/realm/0023_drop_ownership_voucher_table.sql
added apps/astarte_housekeeping/priv/migrations/realm/0024_drop_to2_session_table.sql
added apps/astarte_housekeeping/priv/migrations/realm/0025_add_fdo_guid_to_devices.sql
=== release
{"prerelease":true,"published":"2026-09-30T07:28:27Z","tag":"v1.4.0-rc.6","target":"master"}

$ cd /tmp/opencode && for s in 6c8d4b8 f8da2b9 7714c6b; do echo "=== $s"; timeout 90 gh api repos/astarte-platform/astarte/commits/$s --jq '.commit.message' 2>&1 | head -25; done
=== 6c8d4b8
refactor(fdo)!: remove realm_name requirement from FDO onboarding routes (#2217)

BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix, ensuring full compliance with FDO specification. 

Also, the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure.

Signed-off-by: Eddy Babetto <eddy.babetto@secomind.com>
Co-authored-by: Francesco Noacco <francesco.noacco@secomind.com>
Co-authored-by: frossq <francesco.asquini@secomind.com>
=== f8da2b9
ci: use upstream astarte-device-fdo-rust (#2199)

Signed-off-by: Francesco Noacco <francesco.noacco@secomind.com>
=== 7714c6b
chore(fdo): allow configuring an IP address for FDO URL composition (#2177)

Signed-off-by: Riccardo Nalgi <riccardo.nalgi@secomind.com>
Co-authored-by: Francesco Noacco <francesco.noacco@secomind.com>

Now I have the evidence. Writing the entry at the top of the escalation queue:
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,23 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Scope update for **#78** (FDO): upstream `v1.4.0-rc.6` (prerelease, published **2026-09-30T07:28Z**, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and **the one path #78's body states as fact no longer exists upstream**. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix") deletes the `realm_name` path parameter from all seven TO2 `operation` specs in `fdo_onboarding_controller.ex`. **This breaking change is not in the rc.6 CHANGELOG entry** (which lists only the to0 route, voucher deletion, the `ASTARTE_BASE_URL_*` rename, device-deletion cascade, `device_id` on upload, immediate registration, the two x5chain/voucher-store fixes and the cross-realm delete fix) — so anyone tracking upstream by its changelog, including #78's body, misses it. The realm now comes from the voucher row: the global-keyspace query `fetch_device_id_and_ownership_voucher_and_realm/1` returns `{realm, device_id, voucher_data}` and the `:fdo` plug fetches the session by GUID alone (`plug/fdo_session.ex`: `Session.fetch(guid)`); `RvTO2Addr` no longer carries the realm in the hostname either ("the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure"). **Yours:** #78's acceptance is "a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end" — against rc.6 that SDK posts to `/fdo/101/msg/{60..70}` with no realm, so mirroring rc.5's shape would mean mirroring a path upstream has deleted. No action while parked; the body is just stale.
+- **Client-observable delta, from that same compare — every item here is on the wire, not internal:**
+  - **Two new owner-side operations.** `DELETE /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}` → `204`, `404`, and `500 :rendezvous_revocation_failed` when the rendezvous revoke fails — **the voucher is kept** in that case (`FallbackController` + `OwnershipVoucherView`/`OwnershipVoucherController.delete_ownership_voucher/2`). `POST .../ownership_vouchers/{guid}/to0` re-runs TO0 against the rendezvous server → `200 {"data":{"expiry":…}}`, `404`, `500 :rendezvous_registration_failed`, and **`409 device_already_onboarded`** ("Once Device Onboard completes, the voucher is not registered on the rendezvous anymore"). A guid that is not a UUID is a `404`, not a `422` (`decode_guid/1` maps `Ecto.UUID.dump` failure to `:not_found`).
+  - **`GET .../ownership_vouchers` gains `expiry`** per voucher (new field in the list view; migration `astarte/0007_add_expiry_to_ownership_vouchers.sql`), cleared back to null when TO2 consumes the registration (`mark_voucher_as_claimed` force-changes `expiry` to nil). Additive, so old clients keep working.
+  - **The voucher-upload body is not backwards compatible.** `hw_id` is now **required** and `initial_introspection` is accepted (`api_spec/schemas/ownership_voucher.ex`: `required: [:hw_id, :ownership_voucher, :key_name, :key_algorithm]`), and `LoadRequest.changeset/2` adds three rejectable cases that surface as `400`: `guid has already been claimed`, device `already exists`, and invalid `initial_introspection` (negative or malformed versions). `403` (Forbidden / AuthorizationPathNotMatched) is also now documented on `register`. **Correction to the CHANGELOG wording "Allow specifying device_id on ownership voucher upload": on the wire the client sends `hw_id`, and `device_id` is a `virtual: true` Ecto field derived from it** via `Device.decode_device_id(hw_id, allow_extended_id: true)` — so it is the 128-bit device id, not an arbitrary one, and the new "already exists" check is a realm-device collision check. Astrate must derive it the same way (`pkg/deviceid`) or it will disagree on the id space.
+  - **Invalid vouchers are no longer stored** — the order in `register/2` inverted: rendezvous claim first (`OwnershipVoucher.claim_on_rendezvous/3` → `{:ok, expiry}`), DB write second (`LoadRequest.store_voucher/2`), `500 :rendezvous_registration_failed` if the claim fails. Combined with `fetch_ownership_voucher(guid)` returning `:not_found` for an absent row and `create_ownership_voucher/1` using `overwrite: false` + `stale_error_field: :guid` → `{:error, :duplicated_voucher_guid}`.
+  - **Single-entry x5chain certs accepted** (a PEM carrying one certificate instead of a chain), so vouchers previously rejected at upload now go through.
+  - **The device row now exists before the device has onboarded.** `register/2` ends with `Engine.register_device(realm_name, req.hw_id, initial_introspection: …, with_credentials?: false, fdo_guid: req.device_guid)`: uploading a voucher *creates the realm device* (visible in realm management and the dashboard, `connected: false`, no credentials secret), and the 10-minute unconfirmed window moved off the `unconfirmed_devices` table onto a Cassandra TTL on `devices.credentials_secret` (`DataAccess.Device.add_unconfirmed_credentials/3`, `@ten_minutes_in_seconds 600`; `confirm/2` now returns `{:error, :expired_credentials}` when the TTL lapsed). A device's credentials secret is added mid-TO2 rather than at upload. Re-uploading a voucher for an existing device is refused up front.
+  - **Config rename, operator-visible:** `ASTARTE_BASE_URL_DOMAIN` → `ASTARTE_BASE_URL_HOST`, and the value may now be an IP address as well as a domain (`libs/astarte_fdo/lib/config.ex`, new typed `BaseURLHost` module). A deployment carrying the rc.5 variable name stops booting. This one matters for Astrate even parked, because it is the same key our config reference would have to name.
+- **Two behaviours that look like one feature but are not**, and #78 should not treat them as one: owner-side voucher deletion revokes the rendezvous registration **first** and keeps the voucher if that fails, whereas RM device removal deletes the voucher row **without** revoking — `RealmManagement.DeviceRemoval.Core.delete_ownership_voucher!/2` logs `"Deleting ownership voucher without revoking its rendezvous registration"` (tag `fdo_voucher_deleted_without_revocation`) and its own doc says "N.B.: only the database row is removed." So removing a device can leave a live rendezvous registration pointing at a voucher Astrate/upstream no longer holds. Not a blocker, but it is the kind of asymmetry an operator guide must state.
+- **Not client-observable, so not part of the decision:** the keyspace move itself — vouchers and TO2 sessions now live in the global `astarte` keyspace (`Realm.astarte_keyspace_name()` throughout `fdo/queries.ex`) with a `realm` column on the voucher row and `ALLOW FILTERING` on the list query, so **realm isolation is now application-level rather than keyspace-level** (which is what makes `ensure_voucher_in_realm`-style cross-realm checks necessary, `b3bf731` "ensure users cannot delete ownership vouchers belonging to other realms"); `fdo_guid` added to the device row (migration `realm/0025`) so device removal can find the voucher; the `unconfirmed_devices` drop (`realm/0020`, schema file deleted, `remove_device_from_unconfirmed_devices!` removed); the dashboard voucher page (`astarte-dashboard/src/FdoVoucherPage.tsx`, `hooks/useFdo.ts`, a new cypress spec); and the `data_updater_plant` RPC/GC fixes plus the `v1.3.4` forward-ports.
+- **One premise in #78's own comment is now false.** It says "upstream ships almost no public FDO documentation (code-level only); ours must exceed that." rc.6 adds two documented sections: `doc/pages/user/035-register_device.md` gained `## FIDO Device Onboard` (Owner Keys, Ownership Vouchers) and `## Credentials Secret Lifecycle`, and `doc/pages/architecture/050-pairing_mechanism.md` gained `## FIDO Device Onboarding` (TO0, TO2, Owner Keys management) — commits `162062e`, `d1c39fb`. "Docs must exceed upstream" is still the right requirement, but the gap is far smaller than #78 assumes, and those two pages are now the cheapest available scope input.
+- **Where Astrate stands, for the size estimate:** it has **no FDO surface at all** — the only FDO reference in the tree is the comment on the realm health probe (`internal/pairing/http.go:79`, `internal/pairing/service.go:368`), and `.mule/milestones.md:124` still lists #78 as a deliberately parked item. So none of this is a delta to partially-built code; every item above is net-new scope for a "Size L" issue that is blocked on "who runs a rendezvous server?". The rc.6 delta does not change that call, but it adds four things any eventual implementation must decide: a spec-mandated device-side path with no realm (so any deviation is ours alone), a mandatory `hw_id` in the upload body (so the owner must supply the device id *before* the device has onboarded), an upload that creates a visible realm device as a side effect, and a per-voucher `expiry` the operator has to poll to know when the rendezvous stops serving.
+  Unverified: read-only — every line above comes from the compare API and the rc.6-tagged sources, and the 409/500 mappings come from `FallbackController` source rather than a live call (no Docker on this box). Deliberately no `.mule/todo.md` lines: the deliverable is the scope update, and #78 stays parked until you answer the rendezvous question.
+
 - **github-issues triage run, 2026-10-01: nothing proposable, and the 2026-09-26 reading of the alarm pile is now superseded — this silence is real, and the dead-man's switch is latched off so it cannot report it.** 25 open issues, the same set plus **#113** (09-27 11:07Z) and **#114** (09-27 20:09Z, "nothing has landed in 8h"). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction.
   **The alarm stopped firing because of a bug in the switch, not because the mule recovered.** Last heartbeat **2026-09-28T18:01Z** (`.mule/.heartbeat`), and the last row in `.mule/log.md` is 2026-09-28 (`docs-sync-pairing-unregister-description`, `8ffb39e`) — no task attempt in the ~74h since, while recipe ticks are still committing (today 19:46–21:41Z), so the timer is alive and only the task queue is silent. On **2026-09-29T10:49Z** the pulse did fire (16h) but took the no-`gh` branch: the last line of this file is the fallback text "The mule has been idle 16h" (`mule.sh:771-773`) and **no issue was filed** — #114 is the last alarm ever filed. That same call created `.mule/.alarmed`, which only `beat` clears (`mule.sh:737`), and `check_pulse` returns immediately when it exists (`mule.sh:748`), so the switch has been **latched off for three days** and will file nothing until some tick lands work. The longest silence in the log is the one the switch reports least. That also gives the 2026-09-25 "the body is wrong" note a consequence at last: "Close it — the mule reopens a new one if the silence continues" (`mule.sh:764`) is false by construction — one alarm per silence, never a second.
   **Two cheap causes on this box, one command each.** (a) The queue is 12-of-13 Legion Go: 13 `- [ ]` lines, and the only one runnable here is `flowapi-autorestart-shutdown-cancel` (todo.md:89) — if that is the task the runner is wedged on, it is the whole silence. I probed `ssh -o ConnectTimeout=5 legion` once and it **timed out: the Legion Go is off**, so per the standing rule I did not retry or scan; powering it on unblocks the other twelve. (b) The 09-29 fallback means `gh` was missing from `PATH` or `MULE_PUSH` was empty in that tick's environment, which is why the alarm you needed went to this file instead of GitHub.
   **Proposal, extending the 2026-09-26 one: close #94–#114**, all twenty-one. Each was a one-day false alarm superseded by the next, and as established on 2026-09-25 nothing reads issue state, so closing stays inert — but 21 open "the mule is dead" issues is the noise that hid a silence that turned out to be real. The threshold question (2026-09-26: raise `MULE_IDLE_ALARM_HOURS` to ~24, or pulse off `.mule/log.md` instead of the heartbeat) is still yours and still unanswered; knowing the failure mode now, I would add a third option — stamp the alarm time into `.alarmed` and re-file when the silence doubles, so a continuing outage reports more than once.

$ cd /root/astrate-mule && gofmt -l . 2>&1 | head; echo "gofmt-exit:$?"; timeout 600 go vet ./... 2>&1 | tail -5; echo "vet-exit:$?"
gofmt-exit:0
vet-exit:0

$ cd /root/astrate-mule && timeout 900 go test ./... 2>&1 | grep -Ev '^(ok|---|\?)' | head -20; echo "test-pipeline-exit:$?"
test-pipeline-exit:0

$ cd /root/astrate-mule && git status --short && sed -n '9,16p' .mule/for-giulio.md | cut -c1-140
M .mule/for-giulio.md
line once you have dealt with it — this file is a queue, not a log.

---

- **Scope update for **#78** (FDO): upstream `v1.4.0-rc.6` (prerelease, published **2026-09-30T07:28Z**, CHANGELOG header says 2026-09-28) g
- **Client-observable delta, from that same compare — every item here is on the wire, not internal:**
  - **Two new owner-side operations.** `DELETE /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}` → `204`, `404`, and `500 :rendezvous_rev
  - **`GET .../ownership_vouchers` gains `expiry`** per voucher (new field in the list view; migration `astarte/0007_add_expiry_to_ownership

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,9 +9,9 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
-- **Scope update for **#78** (FDO): upstream `v1.4.0-rc.6` (prerelease, published **2026-09-30T07:28Z**, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and **the one path #78's body states as fact no longer exists upstream**. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix") deletes the `realm_name` path parameter from all seven TO2 `operation` specs in `fdo_onboarding_controller.ex`. **This breaking change is not in the rc.6 CHANGELOG entry** (which lists only the to0 route, voucher deletion, the `ASTARTE_BASE_URL_*` rename, device-deletion cascade, `device_id` on upload, immediate registration, the two x5chain/voucher-store fixes and the cross-realm delete fix) — so anyone tracking upstream by its changelog, including #78's body, misses it. The realm now comes from the voucher row: the global-keyspace query `fetch_device_id_and_ownership_voucher_and_realm/1` returns `{realm, device_id, voucher_data}` and the `:fdo` plug fetches the session by GUID alone (`plug/fdo_session.ex`: `Session.fetch(guid)`); `RvTO2Addr` no longer carries the realm in the hostname either ("the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure"). **Yours:** #78's acceptance is "a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end" — against rc.6 that SDK posts to `/fdo/101/msg/{60..70}` with no realm, so mirroring rc.5's shape would mean mirroring a path upstream has deleted. No action while parked; the body is just stale.
+- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published **2026-09-30T07:28Z**, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and **the one path #78's body states as fact no longer exists upstream**. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix") deletes the `realm_name` path parameter from all seven TO2 `operation` specs in `fdo_onboarding_controller.ex`. **This breaking change is not in the rc.6 CHANGELOG entry** (which lists only the to0 route, voucher deletion, the `ASTARTE_BASE_URL_*` rename, device-deletion cascade, `device_id` on upload, immediate registration, the two x5chain/voucher-store fixes and the cross-realm delete fix) — so anyone tracking upstream by its changelog, including #78's body, misses it. The realm now comes from the voucher row: the global-keyspace query `fetch_device_id_and_ownership_voucher_and_realm/1` returns `{realm, device_id, voucher_data}` and the `:fdo` plug fetches the session by GUID alone (`plug/fdo_session.ex`: `Session.fetch(guid)`); `RvTO2Addr` no longer carries the realm in the hostname either ("the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure"). **Yours:** #78's acceptance is "a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end" — against rc.6 that SDK posts to `/fdo/101/msg/{60..70}` with no realm, so mirroring rc.5's shape would mean mirroring a path upstream has deleted. No action while parked; the body is just stale.
 - **Client-observable delta, from that same compare — every item here is on the wire, not internal:**
   - **Two new owner-side operations.** `DELETE /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}` → `204`, `404`, and `500 :rendezvous_revocation_failed` when the rendezvous revoke fails — **the voucher is kept** in that case (`FallbackController` + `OwnershipVoucherView`/`OwnershipVoucherController.delete_ownership_voucher/2`). `POST .../ownership_vouchers/{guid}/to0` re-runs TO0 against the rendezvous server → `200 {"data":{"expiry":…}}`, `404`, `500 :rendezvous_registration_failed`, and **`409 device_already_onboarded`** ("Once Device Onboard completes, the voucher is not registered on the rendezvous anymore"). A guid that is not a UUID is a `404`, not a `422` (`decode_guid/1` maps `Ecto.UUID.dump` failure to `:not_found`).
   - **`GET .../ownership_vouchers` gains `expiry`** per voucher (new field in the list view; migration `astarte/0007_add_expiry_to_ownership_vouchers.sql`), cleared back to null when TO2 consumes the registration (`mark_voucher_as_claimed` force-changes `expiry` to nil). Additive, so old clients keep working.
   - **The voucher-upload body is not backwards compatible.** `hw_id` is now **required** and `initial_introspection` is accepted (`api_spec/schemas/ownership_voucher.ex`: `required: [:hw_id, :ownership_voucher, :key_name, :key_algorithm]`), and `LoadRequest.changeset/2` adds three rejectable cases that surface as `400`: `guid has already been claimed`, device `already exists`, and invalid `initial_introspection` (negative or malformed versions). `403` (Forbidden / AuthorizationPathNotMatched) is also now documented on `register`. **Correction to the CHANGELOG wording "Allow specifying device_id on ownership voucher upload": on the wire the client sends `hw_id`, and `device_id` is a `virtual: true` Ecto field derived from it** via `Device.decode_device_id(hw_id, allow_extended_id: true)` — so it is the 128-bit device id, not an arbitrary one, and the new "already exists" check is a realm-device collision check. Astrate must derive it the same way (`pkg/deviceid`) or it will disagree on the id space.

$ cd /root/astrate-mule && rg -n '\*\*[^*]*\*\*[^*]*\*\*[^*]*\*\*' .mule/for-giulio.md | cut -c1-120 | head -5
13:- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published **2026-09-30T07:28Z**, CHANGELOG 
15:  - **Two new owner-side operations.** `DELETE /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}` → `204`, `404`, a
17:  - **The voucher-upload body is not backwards compatible.** `hw_id` is now **required** and `initial_introspection` 
22:- **Two behaviours that look like one feature but are not**, and #78 should not treat them as one: owner-side voucher
23:- **Not client-observable, so not part of the decision:** the keyspace move itself — vouchers and TO2 sessions now

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,9 +9,9 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
-- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published **2026-09-30T07:28Z**, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and **the one path #78's body states as fact no longer exists upstream**. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix") deletes the `realm_name` path parameter from all seven TO2 `operation` specs in `fdo_onboarding_controller.ex`. **This breaking change is not in the rc.6 CHANGELOG entry** (which lists only the to0 route, voucher deletion, the `ASTARTE_BASE_URL_*` rename, device-deletion cascade, `device_id` on upload, immediate registration, the two x5chain/voucher-store fixes and the cross-realm delete fix) — so anyone tracking upstream by its changelog, including #78's body, misses it. The realm now comes from the voucher row: the global-keyspace query `fetch_device_id_and_ownership_voucher_and_realm/1` returns `{realm, device_id, voucher_data}` and the `:fdo` plug fetches the session by GUID alone (`plug/fdo_session.ex`: `Session.fetch(guid)`); `RvTO2Addr` no longer carries the realm in the hostname either ("the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure"). **Yours:** #78's acceptance is "a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end" — against rc.6 that SDK posts to `/fdo/101/msg/{60..70}` with no realm, so mirroring rc.5's shape would mean mirroring a path upstream has deleted. No action while parked; the body is just stale.
+- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published 2026-09-30T07:28Z, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and **the one path #78's body states as fact no longer exists upstream**. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix") deletes the `realm_name` path parameter from all seven TO2 `operation` specs in `fdo_onboarding_controller.ex`. **This breaking change is not in the rc.6 CHANGELOG entry** (which lists only the to0 route, voucher deletion, the `ASTARTE_BASE_URL_*` rename, device-deletion cascade, `device_id` on upload, immediate registration, the two x5chain/voucher-store fixes and the cross-realm delete fix) — so anyone tracking upstream by its changelog, including #78's body, misses it. The realm now comes from the voucher row: the global-keyspace query `fetch_device_id_and_ownership_voucher_and_realm/1` returns `{realm, device_id, voucher_data}` and the `:fdo` plug fetches the session by GUID alone (`plug/fdo_session.ex`: `Session.fetch(guid)`); `RvTO2Addr` no longer carries the realm in the hostname either ("the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure"). **Yours:** #78's acceptance is "a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end" — against rc.6 that SDK posts to `/fdo/101/msg/{60..70}` with no realm, so mirroring rc.5's shape would mean mirroring a path upstream has deleted. No action while parked; the body is just stale.
 - **Client-observable delta, from that same compare — every item here is on the wire, not internal:**
   - **Two new owner-side operations.** `DELETE /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}` → `204`, `404`, and `500 :rendezvous_revocation_failed` when the rendezvous revoke fails — **the voucher is kept** in that case (`FallbackController` + `OwnershipVoucherView`/`OwnershipVoucherController.delete_ownership_voucher/2`). `POST .../ownership_vouchers/{guid}/to0` re-runs TO0 against the rendezvous server → `200 {"data":{"expiry":…}}`, `404`, `500 :rendezvous_registration_failed`, and **`409 device_already_onboarded`** ("Once Device Onboard completes, the voucher is not registered on the rendezvous anymore"). A guid that is not a UUID is a `404`, not a `422` (`decode_guid/1` maps `Ecto.UUID.dump` failure to `:not_found`).
   - **`GET .../ownership_vouchers` gains `expiry`** per voucher (new field in the list view; migration `astarte/0007_add_expiry_to_ownership_vouchers.sql`), cleared back to null when TO2 consumes the registration (`mark_voucher_as_claimed` force-changes `expiry` to nil). Additive, so old clients keep working.
   - **The voucher-upload body is not backwards compatible.** `hw_id` is now **required** and `initial_introspection` is accepted (`api_spec/schemas/ownership_voucher.ex`: `required: [:hw_id, :ownership_voucher, :key_name, :key_algorithm]`), and `LoadRequest.changeset/2` adds three rejectable cases that surface as `400`: `guid has already been claimed`, device `already exists`, and invalid `initial_introspection` (negative or malformed versions). `403` (Forbidden / AuthorizationPathNotMatched) is also now documented on `register`. **Correction to the CHANGELOG wording "Allow specifying device_id on ownership voucher upload": on the wire the client sends `hw_id`, and `device_id` is a `virtual: true` Ecto field derived from it** via `Device.decode_device_id(hw_id, allow_extended_id: true)` — so it is the 128-bit device id, not an arbitrary one, and the new "already exists" check is a realm-device collision check. Astrate must derive it the same way (`pkg/deviceid`) or it will disagree on the id space.

$ cd /root/astrate-mule && rg -c '' .mule/for-giulio.md && rg -n '\*\*[^*]*\*\*[^*]*\*\*[^*]*\*\*' .mule/for-giulio.md | cut -c1-90
493
13:- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published 202
15:  - **Two new owner-side operations.** `DELETE /pairing/v1/{realm}/fdo/ownership_vouche
17:  - **The voucher-upload body is not backwards compatible.** `hw_id` is now **required*
22:- **Two behaviours that look like one feature but are not**, and #78 should not treat t
23:- **Not client-observable, so not part of the decision:** the keyspace move itself — 
25:- **Where Astrate stands, for the size estimate:** it has **no FDO surface at all** —
28:- **github-issues triage run, 2026-10-01: nothing proposable, and the 2026-09-26 readin
29:  **The alarm stopped firing because of a bug in the switch, not because the mule recov
30:  **Two cheap causes on this box, one command each.** (a) The queue is 12-of-13 Legion 
35:  **The HTTP forwarder follows redirects, re-sending the realm header and the event bod
41:- **github-issues triage run, 2026-09-26: nothing proposable (as every run since 2026-0
43:  **What the alarm really measures, and why it fires daily by construction.** `check_pu
50:- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstre
52:  **The recipe's standing item is satisfied, and I did not re-file it.** `docs/UPSTREAM
121:- **github-issues triage run, 2026-09-25: still nothing proposable — the mule-alarm 
125:- **Milestone recipe run, 2026-09-24: v3.0 verified unchanged — still gated on upstr
129:- **github-issues triage run, 2026-09-24: still nothing proposable — the mule-alarm 
133:- **Milestone recipe run, 2026-09-23: v3.0 verified unchanged — still gated on upstr
137:- **github-issues triage run, 2026-09-23: still nothing proposable — the mule-alarm 
141:- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are 
145:- **github-issues triage run, 2026-09-21: still nothing proposable — the mule-alarm 
149:- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstr
153:- **github-issues triage re-run, 2026-09-20 (evening): nothing proposable, no change s
157:- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstr
165:- **github-issues triage run, 2026-09-20: still nothing proposable — the mule-alarm 
169:- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstr
173:- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — s
181:- **COMPATIBILITY.md wording update for upstream v1.3.4 (newest stable, 2026-09-18; v1
185:- **github-issues triage run, 2026-09-18: still nothing proposable — the mule-alarm 
189:- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstr
197:- **github-issues triage run, 2026-09-17: still nothing proposable — the mule-alarm 
199:- **github-issues triage run, 2026-09-16: still nothing proposable — the mule-alarm 
201:- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstr
203:- **github-issues triage run, 2026-09-15: still nothing proposable — the mule-alarm 
205:- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstr
207:- **github-issues triage run, 2026-09-13: still nothing proposable — the mule-alarm 
209:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstr
213:- **github-issues triage run, 2026-09-12: still nothing proposable — the mule-alarm 
215:- **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pi
217:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstr
219:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstr
227:- **docs-sync pairing, 2026-09-10: docs/site names the pre-credentials wire status `re
280:  only when upstream ships **v1.4 final** — not yet (still rc). Per **#90's frozen d
334:- **COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable, 2026-07-14; v1
396:- ~~**#87 `lua_map` — needs embedded Lua runtime, parked.**~~ — **closed 2026-09-0

$ cd /root/astrate-mule && sed -n '13p' .mule/for-giulio.md | cut -c1-1000
- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published 2026-09-30T07:28Z, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and **the one path #78's body states as fact no longer exists upstream**. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairi

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,9 +9,9 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
-- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published 2026-09-30T07:28Z, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and **the one path #78's body states as fact no longer exists upstream**. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix") deletes the `realm_name` path parameter from all seven TO2 `operation` specs in `fdo_onboarding_controller.ex`. **This breaking change is not in the rc.6 CHANGELOG entry** (which lists only the to0 route, voucher deletion, the `ASTARTE_BASE_URL_*` rename, device-deletion cascade, `device_id` on upload, immediate registration, the two x5chain/voucher-store fixes and the cross-realm delete fix) — so anyone tracking upstream by its changelog, including #78's body, misses it. The realm now comes from the voucher row: the global-keyspace query `fetch_device_id_and_ownership_voucher_and_realm/1` returns `{realm, device_id, voucher_data}` and the `:fdo` plug fetches the session by GUID alone (`plug/fdo_session.ex`: `Session.fetch(guid)`); `RvTO2Addr` no longer carries the realm in the hostname either ("the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure"). **Yours:** #78's acceptance is "a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end" — against rc.6 that SDK posts to `/fdo/101/msg/{60..70}` with no realm, so mirroring rc.5's shape would mean mirroring a path upstream has deleted. No action while parked; the body is just stale.
+- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published 2026-09-30T07:28Z, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and the one path #78's body states as fact no longer exists upstream. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix") deletes the `realm_name` path parameter from all seven TO2 `operation` specs in `fdo_onboarding_controller.ex`. **This breaking change is not in the rc.6 CHANGELOG entry** (which lists only the to0 route, voucher deletion, the `ASTARTE_BASE_URL_*` rename, device-deletion cascade, `device_id` on upload, immediate registration, the two x5chain/voucher-store fixes and the cross-realm delete fix) — so anyone tracking upstream by its changelog, including #78's body, misses it. The realm now comes from the voucher row: the global-keyspace query `fetch_device_id_and_ownership_voucher_and_realm/1` returns `{realm, device_id, voucher_data}` and the `:fdo` plug fetches the session by GUID alone (`plug/fdo_session.ex`: `Session.fetch(guid)`); `RvTO2Addr` no longer carries the realm in the hostname either ("the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure"). **Yours:** #78's acceptance is "a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end" — against rc.6 that SDK posts to `/fdo/101/msg/{60..70}` with no realm, so mirroring rc.5's shape would mean mirroring a path upstream has deleted. No action while parked; the body is just stale.
 - **Client-observable delta, from that same compare — every item here is on the wire, not internal:**
   - **Two new owner-side operations.** `DELETE /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}` → `204`, `404`, and `500 :rendezvous_revocation_failed` when the rendezvous revoke fails — **the voucher is kept** in that case (`FallbackController` + `OwnershipVoucherView`/`OwnershipVoucherController.delete_ownership_voucher/2`). `POST .../ownership_vouchers/{guid}/to0` re-runs TO0 against the rendezvous server → `200 {"data":{"expiry":…}}`, `404`, `500 :rendezvous_registration_failed`, and **`409 device_already_onboarded`** ("Once Device Onboard completes, the voucher is not registered on the rendezvous anymore"). A guid that is not a UUID is a `404`, not a `422` (`decode_guid/1` maps `Ecto.UUID.dump` failure to `:not_found`).
   - **`GET .../ownership_vouchers` gains `expiry`** per voucher (new field in the list view; migration `astarte/0007_add_expiry_to_ownership_vouchers.sql`), cleared back to null when TO2 consumes the registration (`mark_voucher_as_claimed` force-changes `expiry` to nil). Additive, so old clients keep working.
   - **The voucher-upload body is not backwards compatible.** `hw_id` is now **required** and `initial_introspection` is accepted (`api_spec/schemas/ownership_voucher.ex`: `required: [:hw_id, :ownership_voucher, :key_name, :key_algorithm]`), and `LoadRequest.changeset/2` adds three rejectable cases that surface as `400`: `guid has already been claimed`, device `already exists`, and invalid `initial_introspection` (negative or malformed versions). `403` (Forbidden / AuthorizationPathNotMatched) is also now documented on `register`. **Correction to the CHANGELOG wording "Allow specifying device_id on ownership voucher upload": on the wire the client sends `hw_id`, and `device_id` is a `virtual: true` Ecto field derived from it** via `Device.decode_device_id(hw_id, allow_extended_id: true)` — so it is the 128-bit device id, not an arbitrary one, and the new "already exists" check is a realm-device collision check. Astrate must derive it the same way (`pkg/deviceid`) or it will disagree on the id space.

$ cd /root/astrate-mule && python3 - <<'EOF'
import re
line = open('.mule/for-giulio.md').read().split('\n')[12]
# find bold spans
spans=[(m.start(),m.end()) for m in re.finditer(r'\*\*',line)]
print("markers:",len(spans))
# check nesting: print the lead up to the 2nd marker
end=spans[1][1]
print("LEAD:",line[:end][-60:])
print("lead contains inner ** ?", line[2:end-2].count('**'))
EOF
markers: 8
LEAD: ours, and this only changes how big the eventual build is.**
lead contains inner ** ? 1

$ cd /root/astrate-mule && python3 - <<'EOF'
import re
line = open('.mule/for-giulio.md').read().split('\n')[12]
m=[x for x in re.finditer(r'\*\*',line)]
end=m[1].end()
lead=line[2:end-2]
i=lead.find('**')
print(repr(lead[max(0,i-90):i+90]))
EOF
'**Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published 2026-09-'

$ cd /root/astrate-mule && sed -n '13,27p' .mule/for-giulio.md | cut -c1-100 | nl -ba
1	- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published 2026-09-30T07:28
     2	- **Client-observable delta, from that same compare — every item here is on the wire, not internal
     3	  - **Two new owner-side operations.** `DELETE /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}` �
     4	  - **`GET .../ownership_vouchers` gains `expiry`** per voucher (new field in the list view; migrati
     5	  - **The voucher-upload body is not backwards compatible.** `hw_id` is now **required** and `initia
     6	  - **Invalid vouchers are no longer stored** — the order in `register/2` inverted: rendezvous cla
     7	  - **Single-entry x5chain certs accepted** (a PEM carrying one certificate instead of a chain), so 
     8	  - **The device row now exists before the device has onboarded.** `register/2` ends with `Engine.re
     9	  - **Config rename, operator-visible:** `ASTARTE_BASE_URL_DOMAIN` → `ASTARTE_BASE_URL_HOST`, and 
    10	- **Two behaviours that look like one feature but are not**, and #78 should not treat them as one: o
    11	- **Not client-observable, so not part of the decision:** the keyspace move itself — vouchers and 
    12	- **One premise in #78's own comment is now false.** It says "upstream ships almost no public FDO do
    13	- **Where Astrate stands, for the size estimate:** it has **no FDO surface at all** — the only FDO
    14	  Unverified: read-only — every line above comes from the compare API and the rc.6-tagged sources,
    15

Done: wrote the `v1.4.0-rc.6` scope delta for issue #78 into `.mule/for-giulio.md` as a new top-of-queue entry — verified against the real compare API (27 commits, **153** files, not 154), flagging that the breaking `refactor(fdo)!` (device-side path is now `/fdo/101` with no realm, `6c8d4b8`) is **absent from rc.6's CHANGELOG** and that #78's body states a path upstream has deleted, then splitting the delta into wire-observable (DELETE + `POST .../to0` with 409 `device_already_onboarded`, `expiry` in the voucher list, now-**required** `hw_id` + `initial_introspection` on upload, rendezvous-claim-before-store, single-entry x5chain, device row created at upload with a 10-min TTL credentials secret, `ASTARTE_BASE_URL_DOMAIN`→`_HOST` rename) vs internal (global-keyspace move with `realm` column + `ALLOW FILTERING`, `fdo_guid`, `unconfirmed_devices` drop, dashboard), plus the delete-without-revoking asymmetry in device removal and the fact that upstream now documents FDO in two doc pages; no code, no issue edit, no todo lines.
Files: `.mule/for-giulio.md` (inserted entry at line 13; nothing else touched)
Verified: `gofmt -l .` -> clean; `go vet ./...` -> pass; `go test ./...` -> pass (no Go changed, documentation-only, so no behavioural test is possible or required); every upstream claim read off `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6` and the rc.6-tagged sources in this session
Unsure: the status-code mappings come from `FallbackController` source, not a live call (no Docker on this box); the task line's "`device_id` accepted on voucher upload" is imprecise — the wire field is `hw_id` (required) and `device_id` is a virtual field derived from it via `Device.decode_device_id(..., allow_extended_id: true)`, which is how I reported it; I did not audit the 8 `data_updater_plant`/forward-port commits beyond confirming they are wire-inert
Follow-ups: none
