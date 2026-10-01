slug: recipe-astarte-upstream
verdict: proposed
at:  41a3674
ran: 2026-10-01T20:36:19Z on DietPi in 361s


> build · big-pickle

I'll start by reading the required context files.
→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"' 2>&1 | head -20
v1.4.0-rc.6	2026-09-30T07:28:27Z
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.4") | .body' 2>&1
### Fixed

- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
  uncollected binaries over time by forcing more frequent full sweep garbage collections
  on them.
- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.6") | .body' 2>&1
## Changes since v1.4.0-rc.5

Includes changes from v1.3.4

### Added

- [fdo] Add `POST /v1/:realm_name/fdo/ownership_vouchers/:guid/to0` to run TO0 again for an
  ownership voucher, refreshing how long the rendezvous server keeps serving its registration.
  The resulting expiry is reported when listing the ownership vouchers of a realm.
- [fdo] Add ownership voucher deletion

### Changed

- [fdo] Rename `ASTARTE_BASE_URL_DOMAIN` to `ASTARTE_BASE_URL_HOST`. Its value may now be
  either a domain name or an IP address
- [fdo] Deleting a device also deletes its ownership voucher.
- [fdo] Allow specifying device_id on ownership voucher upload
- [fdo] The device is now immediately registered on ownership voucher upload
- [astarte_data_updater_plant] Improve RPC server reliability

### Fixed

- [dashboard] Add unknown status when device messages aren't properly consumed
- [fdo] Accept single entry x5chain certificates on voucher upload
- [fdo] Invalid vouchers are no longer stored on the database
- [fdo] Ensure users cannot delete ownership vouchers belonging to other realms

$ cd /root/astrate-mule && rg -n '1\.4|1\.3\.4|rc\.[0-9]' docs/COMPATIBILITY.md docs/UPSTREAM-EXPERIMENTAL.md 2>&1 | head -40
docs/UPSTREAM-EXPERIMENTAL.md:24:| Required + encrypted interface mapping fields | #67 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
docs/UPSTREAM-EXPERIMENTAL.md:25:| `async_operation=false` on mutating RM endpoints | #68 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
docs/COMPATIBILITY.md:335:17. **Always synchronous where upstream 1.4 defaults to asynchronous** —
docs/COMPATIBILITY.md:336:    upstream 1.4 runs realm create/delete (Housekeeping), interface
docs/COMPATIBILITY.md:348:    `docs/UPSTREAM-EXPERIMENTAL.md` (#68) until upstream 1.4 is final.

$ ls -la /root/astrate-mule/.mule/ && wc -l /root/astrate-mule/.mule/todo.md /root/astrate-mule/.mule/for-giulio.md 2>&1
total 9636
drwxr-xr-x  7 root root    4096 Sep 29 12:49 .
drwxr-xr-x 17 root root    4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       0 Sep 29 12:49 .alarmed
-rw-r--r--  1 root root      14 Oct  1 22:30 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 28 20:01 .heartbeat
-rw-r--r--  1 root root       0 Oct  1 22:30 .last-output
-rw-r--r--  1 root root       0 Oct  1 22:30 .lock
-rw-r--r--  1 root root      17 Oct  1 22:30 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 9571096 Oct  1 22:30 cron.log
drwxr-xr-x  2 root root    4096 Sep 27 20:23 failed
-rw-r--r--  1 root root   82262 Oct  1 22:09 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root   11657 Sep 28 20:01 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 26 19:42 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root  109116 Sep 28 20:01 todo.md
   246 /root/astrate-mule/.mule/todo.md
   478 /root/astrate-mule/.mule/for-giulio.md
   724 total

$ cd /root/astrate-mule && rg -n 'compat-note|probe-|astarte-upstream|upstream' .mule/todo.md | tail -40
1:- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
3:- [x] compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).
65:`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
81:- [x] docs-sync-rm-triggers-422-nested-envelope [auto]: document the nested changeset envelope on `POST /realmmanagement/v1/{realm}/triggers` in docs/api/astarte_realm_management_api.yaml — when trigger compile fails, `createTrigger` answers 422 with the upstream nested body `{"errors":{"action":{...},"simple_triggers":[...]}}` (internal/realm/http.go:252-258, triggerErrorBody at http.go:268-285, issues #63/#70), but the spec's createTrigger 422 (yaml:344-345) refs only the flat `ValidationError` component (`{"errors":{"detail":...}}`) — which the handler also emits for plain ErrValidation (http.go:259-260, e.g. "trigger requires a name", service.go:429). Model the response as oneOf across the flat ErrorDetail and the nested envelope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
83:- [!] hk-zero-reglimit-create-asymmetry: POST /realms persists device_registration_limit 0 as a literal 0 (service.go:143 only rejects < 0) while upstream Astarte v1.2.0 folds 0 to nil at create (`insert_realm` does `if device_limit == 0, do: nil`, so GET returns null); PATCH 0 already matches upstream (sets 0, http.go:208-214). Fold 0 → nil in CreateRealm for parity, same pattern as the retention fix, and pin with a container-free test. — BLOCKED: opencode exited 1
112:- [x] purge-properties-compression-capability [auto]: upstream v1.3.0 adds a `purge_properties_compression_format` device capability (`zlib`|`plaintext`, default `zlib`) — a wire-visible capability value. Check whether Astrate's capabilities handling (internal/broker, the `<realm>/<device_id>/capabilities` topic, issue #16) needs to recognize/honour it, or whether zlib-only is already the deliberate default; propose the change or note why not needed.
113:- [!] empty-introspection-verification [auto]: upstream v1.3.0 changed "allow devices with empty introspection" — verify whether Astrate's device connection/introspection handling currently rejects an empty introspection string where upstream now accepts it, and propose a fix if so. — BLOCKED: wrote nothing
114:- [!] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages" — investigate only: does Astrate have an analogous delay between trigger install and first delivery? Report, do not patch. — BLOCKED: wrote nothing
115:- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
119:- [!] probe-property-resend-encoding [auto]: upstream v1.3.0 fixed outbound server-property values sent to a device on connect/emptyCache (commit 522ccf4f — a raw stored binaryblob is now wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding, so it ships as a BSON binary, not a string). Investigate only: does Astrate's resendServerProperties/rehydrate path (internal/engine/control.go:144, pkg/payload) emit binaryblob → BSON subtype-0 binary and datetime → UTC datetime exactly like upstream 1.3? Report, do not patch. — BLOCKED: wrote nothing
120:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
126:- [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
134:- [!] appengine-snapshot-ignores-query-params [auto]: in `internal/appengine/data.go:138-148`, the individual-datastream interface-root snapshot branch ignores `Since`/`SinceAfter`/`To`/`Limit`/`Descending` even though the root GET is documented with all of them (docs/api/astarte_appengine_api.yaml:577-582) — the same function refuses incompatible params elsewhere (downsample 422, data.go:118-120/165-167), so this is a silent drop. Probe upstream on the Legion Go to decide honour-vs-refuse, then either apply the window/limit or answer an explicit 422 and stop advertising the params on the root GET; add a T1 test in data_test.go/query_opts_test.go pinning the chosen behaviour. — BLOCKED: wrote nothing
137:- [x] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
138:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
147:- [x] pairing-bearer-secret-test [auto]: add a container-free unit test for `bearerSecret` (internal/pairing/http.go:343) locking the upstream-verbatim permissive match `~r/bearer\:?\s+(.*)$/i` — the colon form `Bearer: <secret>`, the case variants (`bearer `, `BEARER `), and the reject paths (no whitespace before value, non-bearer scheme, empty value). Only the plain `Bearer <secret>` form is exercised today (integration http_test.go); a regression to case-sensitive or drop-colon matching would pass every existing test. Pure helper in a non-integration file; no Docker.
155:- [x] docs-sync-pairing-status-enum [auto]: fix the `PairingInfo.status` enum in docs/api/astarte_pairing_api.yaml from `[confirmed, pending, denied, expired]` (yaml:412) to `[confirmed, pending, inhibited]` — the handler never emits `denied` or `expired`: `service.Info` maps only `DeviceStatusInhibited`→"inhibited", `DeviceStatusConfirmed`→"confirmed", and everything else→"pending" (internal/pairing/service.go:296-299, statuses defined in store/devices.go:19-27; upstream-parity per the comment at service.go:282-285). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
156:- [x] docs-sync-pairing-version-endpoint [auto]: add the undocumented `GET /pairing/v1/{realm}/version` route to docs/api/astarte_pairing_api.yaml — it exists in code (cmd/astrate/main.go:398, unauthenticated, upstream-parity per main.go:387-391) and answers 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38); mirror the RM spec's getVersion op (astarte_realm_management_api.yaml:532) without its a_rma security. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
177:- [x] compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.
181:- [x] auth-iat-not-required-test [auto]: pin in internal/auth/jwt_test.go that a token with a *future* `iat` still verifies — the "`iat` is not required — upstream parity" rule (jwt.go:104-105, no `jwt.WithIssuedAt` in the parser, jwt.go:115-118) has zero coverage, so a regression that starts validating `iat` or `WithIssuedAt` would silently reject tokens from a skewed-clock issuer and pass every existing suite. Sign RS256 with `iat` = now+1h and no `exp`/`nbf`, assert `Verify` succeeds.
184:- [x] hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto] — already resolved and re-queued by the refill: commit 015a706 chose the upstream-parity behaviour (0 ≡ unset, not invalid, not a literal 0) and pinned it — create folds 0 → nil with the container-free TestCreateRealmZeroRetentionUnsets (internal/housekeeping/service_test.go:186) and PATCH 0 → ClearRetention in the integration suite (internal/housekeeping/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).
194:- [x] docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
197:- [x] channels-group-watch-membership-scope: group watches authorize against `WATCH::groups/<name>/…` (ws.go:178-205, conformance channels.json) but `Room.dispatch` ever only filters by `watchEntry.deviceID` (empty for a group watch — room.go:36-40, 196-200) and the trigger compiles with no group/device scope (triggers/match.go:374-396, because upstream carries group_name at the payload top level), so a group-scoped token receives every realm event matching interface/path/on, members or not. Plumb a group-membership resolver into internal/appengine/channels (it knows only `Bus` today), resolve membership at Watch time and store it in `watchEntry`, filter deliveries by it. Add a container-free test with a fake resolver: watch group `probe`, publish an event from a member and a non-member, assert only the member's arrives. Verify the member-filter rule against upstream on the Legion if disputed; the boundary leak stands regardless.
199:- [!] channels-rejoin-authz-mismatch: `handleJoin`'s comment says upstream authorizes every join including a rejoin (ws.go:276-277), but the rejoin branch (ws.go:284-291) replies OK without calling `s.tok.AuthorizesChannel(auth.VerbJoin, name)` — behaviourally identical today (a token is immutable for the session) but the comment lies about what the code does. Either move the authz check to cover the rejoin branch too, or fix the comment; mechanical, no behaviour change, choose one. — BLOCKED: wrote nothing
200:- [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a locally-run `elixir -e` with `Base.url_decode64!`) with the device ID `"AAAAAAAAAAAAAAAAAAAAAB"` — nonzero unused bits in the 22nd base64url char — and report whether upstream accepts it: `pkg/deviceid/deviceid.go:34-37` uses `base64.RawURLEncoding.Strict()` and deviceid_test.go:176 asserts rejection while the comment claims parity with `Elixir Base.url_decode64!(padding: false)`, which decodes and discards those bits rather than erroring. Astrate may therefore be stricter than upstream on the same wire form (a device registered with such an id on upstream would be rejected here). If upstream rejects it identically, the strict encoding is verified parity and the task is done; if upstream accepts it, escalate the relax-vs-strict call to `.mule/for-giulio.md` with the measurement. Probe first, no code change either way.
203:- [ ] store-devices-inhibit-re-register [legion] [auto]: `RegisterDevice` (internal/store/devices.go:75-91) silently clears the inhibit flag — its `ON CONFLICT ... DO UPDATE SET status = 'registered'` fires for any device with `first_credentials_request IS NULL`, including one the admin inhibited via `SetDeviceInhibited` (devices.go:251-268), which sets `status='inhibited'` on unconfirmed devices too; §5.3 says an inhibited device blocks new credentials and connections, so the re-registration re-opens it. Preserve `'inhibited'` in the SET (e.g. `status = CASE WHEN devices.status = 'inhibited' THEN 'inhibited' ELSE 'registered' END`) and add a Lifecycle case in internal/store/devices_test.go: inhibit an unconfirmed device, re-register, assert status stays inhibited and the secret still rotates. Verify the assert against upstream's register-not-touching-inhibit on the Legion while the integration suite runs.
219:- [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
227:- [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
228:- [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
230:- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
235:- [x] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.
241:- [x] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:556) and `InternalServerError` shows `detail: Internal Server Error` (yaml:638), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37, golden fixture pkg/astarteapi/testdata/error_bad_request.json) and `"Internal server error"` (envelope.go:58, testdata/error_internal.json), both measured on upstream 1.2.0 (envelope.go:33-36, 55-57). The other five detail examples on this file (`Not Found`, `Device not found`, `Unauthorized`, `Forbidden`, `Too Many Requests`) are already correct — change only those two. The identical wrong pair is in the other four specs (appengine 1732/1792, housekeeping 362/385/477, realm_management 1265/1348, native 466); leave those for a run of their surface. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
243:- [x] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` path parameter (docs/api/astarte_pairing_api.yaml:360) and `RegisterRequest.hw_id` (yaml:373) — but `deviceid.Parse` accepts only the 22-character **unpadded base64url** form, in strict mode, and rejects the standard alphabet's `+` and `/` and any `=` padding with an error wrapping `ErrInvalid` (pkg/deviceid/deviceid.go:26, 37, 39-60; upstream parity cited at 36). Measured: with the same `encoding` the parser uses, `dT6hS2W9TT6LEnP25ks+lg` and `.../lg` fail at byte 19 and the padded 23-character form fails at byte 21, so a client generated from the current spec can send a well-formed standard-base64 ID and get 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}` (internal/pairing/http.go:307-309). State the exact encoding and the 22-character length in both descriptions and add a `pattern`/`minLength`/`maxLength` to `hw_id`. The existing `hw_id: dT6hS2W9TT6LEnP25ks_lg` example (yaml:123) is valid and stays. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
246:- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `base64.RawURLEncoding.Strict()` (pkg/deviceid/deviceid.go:37), so of the 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with `ErrInvalid` while `..._lA` is accepted, and the pairing spec's `hw_id` pattern `^[A-Za-z0-9_-]{22}$` (added by docs-sync-pairing-deviceid-base64url) deliberately does not encode the rule because a pattern is a client-side contract and upstream's leniency is unmeasured. The docstring justifies strictness as byte-identical round-tripping, which is a reason, not a parity claim, and docs/COMPATIBILITY.md records no entry for it. Measure upstream (Elixir `Base.url_decode64!(padding: false)`, which does not check those bits) on the Legion Go, then either record the divergence in docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.

→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ cd /root/astrate-mule && rg -n '^- \[ \]' .mule/todo.md | tail -30 && echo "---- total lines" && wc -l .mule/todo.md
86:- [ ] broker-external-bus-intake [legion]: implement the second Intake implementation named in the TODO at internal/broker/intake.go:56-62 — an external-bus-backed (e.g. NATS JetStream) Intake for multi-instance deployment / restart survival, reproducing the embedded broker's per-device ordering, deferred-ack backpressure and QoS 0 drop semantics; an implementation of this interface needs containerised integration tests, so it runs on the Legion Go.
88:- [ ] hygiene-govulncheck [legion]: on the Legion Go, run govulncheck ./... on a fresh ~/astrate clone and file a task line for any REACHABLE advisory it reports (name the CVE and the call path, per the hygiene recipe); this box cannot build it — go install golang.org/x/vuln/cmd/govulncheck@latest is OOM-killed here (3.7GB, 4 cores), so the recipe's highest-priority check never runs unless done there.
89:- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.go:560) with no stop signal, so a block-death racing process shutdown (cmd/astrate/main.go:482-488 drains the Manager but never cancels the goroutine) can rebuild the flow via `mgr.StartFlow` on a background pump (internal/flow/flow.go:173-175) after `Manager.Shutdown` and flip the durable status back to running. Thread a stop channel/context through the Service, checked in the loop's sleep, cancelled on shutdown; verification needs a `[legion]` integration or timing-based test. [auto]
101:- [ ] race-check-store: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/store/... ./internal/housekeeping/... ./migrations/...`. Report any failure to .mule/for-giulio.md with the full race report. Split out of the former single `race-check` line, which timed out running the whole tree at once. [legion] [readonly]
102:- [ ] race-check-engine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/engine/... ./internal/broker/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
103:- [ ] race-check-flow: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/flow/... ./internal/realm/... ./internal/pairing/... ./internal/auth/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
104:- [ ] race-check-appengine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/appengine/... ./internal/observability/... ./internal/httpx/... ./internal/config/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
105:- [ ] race-check-pkg: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./pkg/... ./cmd/... ./internal/testutil/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
200:- [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a locally-run `elixir -e` with `Base.url_decode64!`) with the device ID `"AAAAAAAAAAAAAAAAAAAAAB"` — nonzero unused bits in the 22nd base64url char — and report whether upstream accepts it: `pkg/deviceid/deviceid.go:34-37` uses `base64.RawURLEncoding.Strict()` and deviceid_test.go:176 asserts rejection while the comment claims parity with `Elixir Base.url_decode64!(padding: false)`, which decodes and discards those bits rather than erroring. Astrate may therefore be stricter than upstream on the same wire form (a device registered with such an id on upstream would be rejected here). If upstream rejects it identically, the strict encoding is verified parity and the task is done; if upstream accepts it, escalate the relax-vs-strict call to `.mule/for-giulio.md` with the measurement. Probe first, no code change either way.
203:- [ ] store-devices-inhibit-re-register [legion] [auto]: `RegisterDevice` (internal/store/devices.go:75-91) silently clears the inhibit flag — its `ON CONFLICT ... DO UPDATE SET status = 'registered'` fires for any device with `first_credentials_request IS NULL`, including one the admin inhibited via `SetDeviceInhibited` (devices.go:251-268), which sets `status='inhibited'` on unconfirmed devices too; §5.3 says an inhibited device blocks new credentials and connections, so the re-registration re-opens it. Preserve `'inhibited'` in the SET (e.g. `status = CASE WHEN devices.status = 'inhibited' THEN 'inhibited' ELSE 'registered' END`) and add a Lifecycle case in internal/store/devices_test.go: inhibit an unconfirmed device, re-register, assert status stays inhibited and the secret still rotates. Verify the assert against upstream's register-not-touching-inhibit on the Legion while the integration suite runs.
205:- [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (internal/store/devices.go:130-131, "if several devices share an alias the lowest device ID wins", `ORDER BY id LIMIT 1`) — no test drives it today; add a case in internal/store/devices_test.go that gives two devices the same alias and asserts the lookup resolves the lower ID. Needs the integration DB.
221:- [ ] flow-boot-resumes-stopped-durable-flows [legion] [auto]: RehydrateAutoRestart (cmd/astrate/main.go:216) restarts every flow with auto_restart=true whatever its last status — `ListAutoRestartFlows` filters on `f.auto_restart = true` alone (internal/store/flows.go:125) and nothing ever clears the flag, which is written once at create (flowapi/http.go:141-146 -> service.go:318) and never on stop — so a durable flow stopped through the API is running again after the next boot, while the shutdown mark that says otherwise (MarkRunningFlowsStopped, main.go:497 -> internal/flowapi/service.go:885-901, which writes status="stopped") is never read by that query. The two halves disagree: service.go:325 documents "every durable flow", the shutdown call assumes "only what was running". Say which is the intent, then make them agree — either add `AND f.status = 'running'` to flows.go:125 (making the shutdown mark load-bearing) or delete the mark as dead code — and add a case in internal/store/flows_test.go pinning the chosen rule (a stopped auto_restart row must, or must not, come back). Needs the DB.
246:- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `base64.RawURLEncoding.Strict()` (pkg/deviceid/deviceid.go:37), so of the 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with `ErrInvalid` while `..._lA` is accepted, and the pairing spec's `hw_id` pattern `^[A-Za-z0-9_-]{22}$` (added by docs-sync-pairing-deviceid-base64url) deliberately does not encode the rule because a pattern is a client-side contract and upstream's leniency is unmeasured. The docstring justifies strictness as byte-identical round-tripping, which is a reason, not a parity claim, and docs/COMPATIBILITY.md records no entry for it. Measure upstream (Elixir `Base.url_decode64!(padding: false)`, which does not check those bits) on the Legion Go, then either record the divergence in docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.
---- total lines
246 .mule/todo.md

$ cd /root/astrate-mule && rg -in 'voucher|to0|fdo|rc\.6|rc\.5|rendezvous' .mule/todo.md | head -30
120:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
138:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
177:- [x] compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.

$ cd /root/astrate-mule && rg -il 'voucher|to0|rendezvous' --glob '!.git' . | head -40
(no output)

$ cd /root/astrate-mule && rg -in 'fdo' --glob '*.go' --glob '*.yaml' --glob '*.yml' --glob '*.md' -l | head -20
docs/api/astarte_pairing_api.yaml
docs/COMPATIBILITY.md
internal/pairing/http.go
internal/pairing/service.go

$ cd /root/astrate-mule && rg -in 'fdo' -g '!.git' | head -30
docs/api/astarte_pairing_api.yaml:33:        and unrate-limited — FDO-flow devices probe it before they hold
internal/pairing/http.go:79:	// FDO-flow devices probe it before they hold credentials. No rate
docs/COMPATIBILITY.md:361:    answers 404 and an unhealthy database 503, and FDO-flow devices can probe
internal/pairing/service.go:368:// design — FDO-flow devices probe it before they hold credentials).

$ cd /root/astrate-mule && gh api 'repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6' --jq '{ahead: .ahead_by, files: [.files[].filename]}' 2>&1 | head -60
{"ahead":27,"files":[".env",".github/workflows/astarte-build-workflow.yaml",".github/workflows/astarte-end-to-end-test-workflow.yaml",".gitignore",".typos.toml","CHANGELOG.md","apps/astarte_appengine_api/mix.exs","apps/astarte_appengine_api/mix.lock","apps/astarte_appengine_api/test/support/helpers/database.ex","apps/astarte_appengine_api/test/support/helpers/database_v2.ex","apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/rpc/server.ex","apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/rpc/supervisor.ex","apps/astarte_data_updater_plant/mix.exs","apps/astarte_data_updater_plant/mix.lock","apps/astarte_data_updater_plant/test/support/database_test_helper.ex","apps/astarte_data_updater_plant/test/support/helpers/database.ex","apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex","apps/astarte_housekeeping/lib/astarte_housekeeping/realms/queries.ex","apps/astarte_housekeeping/mix.exs","apps/astarte_housekeeping/mix.lock","apps/astarte_housekeeping/priv/migrations/astarte/0004_create_session_key_type.sql","apps/astarte_housekeeping/priv/migrations/astarte/0005_create_to2_sessions_table.sql","apps/astarte_housekeeping/priv/migrations/astarte/0006_create_ownership_vouchers_table.sql","apps/astarte_housekeeping/priv/migrations/astarte/0007_add_expiry_to_ownership_vouchers.sql","apps/astarte_housekeeping/priv/migrations/realm/0020_drop_unconfirmed_devices_table.sql","apps/astarte_housekeeping/priv/migrations/realm/0021_add_device_id_to_ownership_voucher.sql","apps/astarte_housekeeping/priv/migrations/realm/0022_delete_vouchers.sql","apps/astarte_housekeeping/priv/migrations/realm/0023_drop_ownership_voucher_table.sql","apps/astarte_housekeeping/priv/migrations/realm/0024_drop_to2_session_table.sql","apps/astarte_housekeeping/priv/migrations/realm/0025_add_fdo_guid_to_devices.sql","apps/astarte_housekeeping/test/support/helpers/database.ex","apps/astarte_pairing/config/dev.exs","apps/astarte_pairing/config/test.exs","apps/astarte_pairing/lib/astarte_pairing/agent/agent.ex","apps/astarte_pairing/lib/astarte_pairing/agent/device_registration_request.ex","apps/astarte_pairing/lib/astarte_pairing/engine.ex","apps/astarte_pairing/lib/astarte_pairing/queries.ex","apps/astarte_pairing/lib/astarte_pairing_web/api_spec/schemas/agent.ex","apps/astarte_pairing/lib/astarte_pairing_web/api_spec/schemas/ownership_voucher.ex","apps/astarte_pairing/lib/astarte_pairing_web/controllers/fallback_controller.ex","apps/astarte_pairing/lib/astarte_pairing_web/controllers/fdo_onboarding_controller.ex","apps/astarte_pairing/lib/astarte_pairing_web/controllers/ownership_voucher_controller.ex","apps/astarte_pairing/lib/astarte_pairing_web/plug/fdo_session.ex","apps/astarte_pairing/lib/astarte_pairing_web/router.ex","apps/astarte_pairing/lib/astarte_pairing_web/views/error_view.ex","apps/astarte_pairing/lib/astarte_pairing_web/views/ownership_voucher_view.ex","apps/astarte_pairing/mix.exs","apps/astarte_pairing/mix.lock","apps/astarte_pairing/test/astarte_pairing/certs/engine_test.exs","apps/astarte_pairing/test/astarte_pairing_web/controllers/fdo_onboarding_controller_test.exs","apps/astarte_pairing/test/astarte_pairing_web/controllers/ownership_voucher_controller_test.exs","apps/astarte_pairing/test/support/cases/device.ex","apps/astarte_pairing/test/support/cases/fdo_session.ex","apps/astarte_pairing/test/support/helpers/database.ex","apps/astarte_pairing/test/support/helpers/device.ex","apps/astarte_pairing/test/support/helpers/fdo.ex","apps/astarte_pairing/test/test_helper.exs","apps/astarte_realm_management/lib/astarte_realm_management/application.ex","apps/astarte_realm_management/lib/astarte_realm_management/device_removal/core.ex","apps/astarte_realm_management/lib/astarte_realm_management/device_removal/device_remover.ex","apps/astarte_realm_management/lib/astarte_realm_management/device_removal/queries.ex","apps/astarte_realm_management/lib/astarte_realm_management/device_removal/scheduler.ex","apps/astarte_realm_management/mix.exs","apps/astarte_realm_management/mix.lock","apps/astarte_realm_management/test/astarte_realm_management/device_removal/device_remover_test.exs","apps/astarte_realm_management/test/astarte_realm_management/device_removal/scheduler_test.exs","apps/astarte_realm_management/test/support/helpers/database.ex","apps/astarte_trigger_engine/mix.exs","apps/astarte_trigger_engine/mix.lock","apps/astarte_trigger_engine/test/support/helpers/database.ex","astarte-dashboard/cypress/e2e/fdo_vouchers_page.cy.js","astarte-dashboard/package.json","astarte-dashboard/src/DeviceStatusPage/DeviceInfoCard.tsx","astarte-dashboard/src/DevicesPage.tsx","astarte-dashboard/src/FdoVoucherPage.tsx","astarte-dashboard/src/FdoVouchersPage.tsx","astarte-dashboard/src/RegisterDevicePage.tsx","astarte-dashboard/src/astarte-client/client.ts","astarte-dashboard/src/astarte-client/models/Device/index.ts","astarte-dashboard/src/components/Icon.tsx","astarte-dashboard/src/components/IntrospectionTable.tsx","astarte-dashboard/src/hooks/useFdo.ts","compose/traefik/traefik.yaml","doc/mix.exs","doc/pages/architecture/050-pairing_mechanism.md","doc/pages/user/035-register_device.md","docker-compose.yml","libs/astarte_adapters/mix.lock","libs/astarte_config/mix.exs","libs/astarte_config/mix.lock","libs/astarte_data_access/lib/astarte_data_access/device.ex","libs/astarte_data_access/lib/astarte_data_access/device/unconfirmed_device.ex","libs/astarte_data_access/lib/astarte_data_access/devices/device.ex","libs/astarte_data_access/lib/astarte_data_access/fdo/ownership_voucher.ex","libs/astarte_data_access/lib/astarte_data_access/fdo/queries.ex","libs/astarte_data_access/lib/astarte_data_access/fdo/to2_session.ex","libs/astarte_data_access/mix.exs","libs/astarte_data_access/mix.lock","libs/astarte_data_access/test/device_test.exs","libs/astarte_data_access/test/fdo/ownership_voucher/ownership_voucher_test.exs","libs/astarte_data_access/test/fdo/queries_test.exs","libs/astarte_data_access/test/support/database_test_helper.exs","libs/astarte_events/mix.exs","libs/astarte_events/mix.lock","libs/astarte_events/test/support/helpers/database_test_helper.ex","libs/astarte_fdo/config/test.exs","libs/astarte_fdo/lib/config.ex","libs/astarte_fdo/lib/config/base_url_host.ex","libs/astarte_fdo/lib/owner_onboarding.ex","libs/astarte_fdo/lib/owner_onboarding/session.ex","libs/astarte_fdo/lib/ownership_voucher/load_request.ex","libs/astarte_fdo/lib/ownership_voucher/ownership_voucher.ex","libs/astarte_fdo/lib/service_info.ex","libs/astarte_fdo/lib/to0.ex","libs/astarte_fdo/mix.lock","libs/astarte_fdo/test/astarte_fdo/config/base_url_host_test.exs","libs/astarte_fdo/test/astarte_fdo/config_test.exs","libs/astarte_fdo/test/astarte_fdo/onboarding/done_test.exs","libs/astarte_fdo/test/astarte_fdo/onboarding/prove_device_test.exs","libs/astarte_fdo/test/astarte_fdo/owner_onboarding/owner_onboarding_test.exs","libs/astarte_fdo/test/astarte_fdo/owner_onboarding/session_test.exs","libs/astarte_fdo/test/astarte_fdo/ownership_voucher/load_request_test.exs","libs/astarte_fdo/test/astarte_fdo/ownership_voucher_test.exs","libs/astarte_fdo/test/astarte_fdo/rendezvous/core_test.exs","libs/astarte_fdo/test/astarte_fdo/rendezvous/rv_to2_addr_test.exs","libs/astarte_fdo/test/astarte_fdo/service_info_test.exs","libs/astarte_fdo/test/astarte_fdo/to0_test.exs","libs/astarte_fdo/test/support/cases/fdo_session.ex","libs/astarte_fdo/test/support/helpers/database.ex","libs/astarte_fdo/test/support/helpers/fdo.ex","libs/astarte_fdo/test/test_helper.exs","libs/astarte_fdo_core/lib/rendezvous/rv_to2_addr.ex","libs/astarte_fdo_core/mix.exs","libs/astarte_fdo_core/mix.lock","libs/astarte_fdo_core/test/astarte_fdo_core/fdo/rendezvous/rv_to2_addr_test.exs","libs/astarte_generators/mix.exs","libs/astarte_generators/mix.lock","libs/astarte_rpc/mix.exs","libs/astarte_rpc/mix.lock","libs/astarte_rpc/test/support/helpers/database.ex","libs/astarte_secrets/lib/astarte_secrets/astarte_secrets.ex","libs/astarte_secrets/mix.exs","libs/astarte_secrets/mix.lock","libs/astarte_test_suite/lib/astarte/test_suite/helpers/realm.ex","libs/astarte_test_suite/mix.lock","tools/astarte_device_fleet_simulator/mix.exs","tools/astarte_e2e/Dockerfile","tools/astarte_e2e/mix.exs","tools/astarte_e2e/mix.lock","tools/astarte_export/mix.exs","tools/astarte_export/mix.lock","tools/astarte_import/mix.exs","tools/astarte_import/mix.lock"]}

$ cd /root/astrate-mule && gh api 'repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6' --jq '.files[] | select(.filename=="apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex" or .filename=="apps/astarte_housekeeping/lib/astarte_housekeeping/realms/queries.ex") | "=== \(.filename)\n\(.patch)"' 2>&1 | head -120
=== apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex
@@ -40,4 +40,5 @@ defmodule Astarte.Housekeeping.Realms.NonNegativeIntegerOrUnsetType do
   def dump(:unset), do: {:ok, :unset}
   def dump(n) when is_integer(n), do: {:ok, n}
   def dump(_), do: :error
+  # coveralls-ignore-stop
 end
=== apps/astarte_housekeeping/lib/astarte_housekeeping/realms/queries.ex
@@ -310,9 +310,6 @@ defmodule Astarte.Housekeeping.Realms.Queries do
          :ok <- create_simple_triggers_table(keyspace_name),
          :ok <- create_grouped_devices_table(keyspace_name),
          :ok <- create_deletion_in_progress_table(keyspace_name),
-         :ok <- create_unconfirmed_devices_table(keyspace_name),
-         :ok <- create_ownership_vouchers_table(keyspace_name),
-         :ok <- create_to2_sessions_table(keyspace_name),
          :ok <- insert_realm_public_key(keyspace_name, public_key_pem),
          :ok <- insert_realm_astarte_schema_version(keyspace_name),
          :ok <- insert_realm(realm_name, device_limit),
@@ -650,42 +647,6 @@ defmodule Astarte.Housekeeping.Realms.Queries do
     end
   end
 
-  defp create_unconfirmed_devices_table(keyspace_name) do
-    query = """
-    CREATE TABLE #{keyspace_name}.unconfirmed_devices (
-      device_id uuid,
-      created_at timestamp,
-      PRIMARY KEY (device_id)
-    );
-    """
-
-    with {:ok, %{rows: nil, num_rows: 1}} <- CSystem.execute_schema_change(query) do
-      :ok
-    end
-  end
-
-  defp create_ownership_vouchers_table(keyspace_name) do
-    query = """
-    CREATE TABLE #{keyspace_name}.ownership_vouchers (
-      guid blob,
-      voucher_data blob,
-      output_voucher blob,
-      replacement_guid blob,
-      replacement_rendezvous_info blob,
-      replacement_public_key blob,
-      key_name varchar,
-      key_algorithm int,
-      user_id blob,
-      status int,
-      PRIMARY KEY (guid)
-    );
-    """
-
-    with {:ok, %{rows: nil, num_rows: 1}} <- CSystem.execute_schema_change(query) do
-      :ok
-    end
-  end
-
   defp create_session_key_type(keyspace_name) do
     query = """
     CREATE TYPE #{keyspace_name}.session_key (
@@ -699,40 +660,6 @@ defmodule Astarte.Housekeeping.Realms.Queries do
     end
   end
 
-  defp create_to2_sessions_table(keyspace_name) do
-    query = """
-    CREATE TABLE #{keyspace_name}.to2_sessions (
-      guid blob,
-      device_id uuid,
-      hmac blob,
-      nonce blob,
-      sig_type int,
-      epid_group blob,
-      device_public_key blob,
-      prove_dv_nonce blob,
-      setup_dv_nonce blob,
-      kex_suite_name ascii,
-      cipher_suite_name int,
-      max_owner_service_info_size int,
-      owner_random blob,
-      secret blob,
-      sevk session_key,
-      svk session_key,
-      sek session_key,
-      device_service_info map<tuple<text, text>, blob>,
-      owner_service_info list<blob>,
-      last_chunk_sent int,
-      replacement_hmac blob,
-      PRIMARY KEY (guid)
-    )
-    WITH default_time_to_live = 7200;
-    """
-
-    with {:ok, %{rows: nil, num_rows: 1}} <- CSystem.execute_schema_change(query) do
-      :ok
-    end
-  end
-
   defp create_grouped_devices_table(keyspace_name) do
     query = """
     CREATE TABLE #{keyspace_name}.grouped_devices (
@@ -857,6 +784,7 @@ defmodule Astarte.Housekeeping.Realms.Queries do
       attributes map<varchar, varchar>,
       groups map<varchar, timeuuid>,
       capabilities capabilities,
+      fdo_guid blob,
 
       PRIMARY KEY (device_id)
     );
@@ -1201,6 +1129,9 @@ defmodule Astarte.Housekeeping.Realms.Queries do
     with :ok <- create_astarte_keyspace(default_replication),
          :ok <- create_realms_table(),
          :ok <- create_astarte_kv_store(),
+         :ok <- create_astarte_session_key_type(),
+         :ok <- create_astarte_to2_sessions_table(),
+         :ok <- create_astarte_ownership_vouchers_table(),
          :ok <- insert_astarte_schema_version(),
          keyspace_replication_map = keyspace_replication_map(default_replication),
          :ok <- save_keyspace_replication(keyspace_replication_map) do

$ cd /root/astrate-mule && gh api 'repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6' --jq '.files[] | select(.filename|test("astarte-dashboard")) | "=== \(.filename)\n\(.patch)"' 2>&1 | head -150
=== astarte-dashboard/cypress/e2e/fdo_vouchers_page.cy.js
@@ -0,0 +1,125 @@
+const hoursFromNow = (hours) => new Date(Date.now() + hours * 3600 * 1000).toISOString();
+
+// Built once per test: the rendered expiry is compared against these exact
+// timestamps, so they must not be recomputed between the stub and the assertion.
+const buildVouchers = () => ({
+  active: {
+    guid: '11111111-1111-1111-1111-111111111111',
+    status: 'created',
+    input_voucher: '-----BEGIN OWNERSHIP VOUCHER-----\nQUJD\n-----END OWNERSHIP VOUCHER-----\n',
+    output_voucher: null,
+    output_guid: null,
+    expiry: hoursFromNow(1),
+  },
+  expired: {
+    guid: '22222222-2222-2222-2222-222222222222',
+    status: 'created',
+    input_voucher: null,
+    output_voucher: null,
+    output_guid: null,
+    expiry: hoursFromNow(-1),
+  },
+  claimed: {
+    guid: '33333333-3333-3333-3333-333333333333',
+    status: 'claimed',
+    input_voucher: null,
+    output_voucher: null,
+    output_guid: null,
+    expiry: null,
+  },
+});
+
+const voucherRow = (guid) => cy.contains('tbody tr', guid);
+
+describe('FDO vouchers page tests', () => {
+  context('no access before login', () => {
+    it('redirects to login', function () {
+      cy.visit('/fdo-vouchers');
+      cy.location('pathname').should('eq', '/login');
+    });
+  });
+
+  context('authenticated', () => {
+    beforeEach(function () {
+      const vouchers = buildVouchers();
+      cy.wrap(vouchers).as('vouchers');
+
+      cy.fixture('realm').then((realm) => {
+        cy.intercept('GET', `/pairing/v1/${realm.name}/fdo/ownership_vouchers`, {
+          body: { data: [vouchers.active, vouchers.expired, vouchers.claimed] },
+        }).as('vouchersRequest');
+        cy.login();
+        cy.visit('/fdo-vouchers');
+        cy.wait('@vouchersRequest');
+      });
+    });
+
+    it('lists the rendezvous expiry of every voucher', function () {
+      const { active } = this.vouchers;
+
+      cy.get('h2').contains('FDO Ownership Vouchers');
+      cy.get('thead th').eq(2).contains('Rendezvous expiry');
+
+      voucherRow(active.guid).within(() => {
+        cy.contains(new Date(active.expiry).toLocaleString());
+        cy.contains('Expired').should('not.exist');
+      });
+    });
+
+    it('flags a voucher whose registration is no longer served', function () {
+      voucherRow(this.vouchers.expired.guid).contains('Expired');
+    });
+
+    it('shows no expiry for a voucher whose device completed Device Onboard', function () {
+      voucherRow(this.vouchers.claimed.guid).within(() => {
+        cy.get('td').eq(2).should('contain', '\u2014');
+        cy.contains('Re-run TO0').should('not.exist');
+      });
+    });
+
+    it('re-runs TO0 and reports the refreshed expiry', function () {
+      const { expired } = this.vouchers;
+      const refreshedExpiry = hoursFromNow(2);
+
+      cy.fixture('realm').then((realm) => {
+        cy.intercept(
+          'POST',
+          `/pairing/v1/${realm.name}/fdo/ownership_vouchers/${expired.guid}/to0`,
+          { body: { data: { expiry: refreshedExpiry } } },
+        ).as('runTo0');
+
+        voucherRow(expired.guid).contains('Re-run TO0').click();
+        cy.wait('@runTo0');
+
+        cy.get('.alert-success').should(
+          'contain',
+          `The registration is now served until ${new Date(refreshedExpiry).toLocaleString()}`,
+        );
+      });
+    });
+
+    it('explains a TO0 rejected because Device Onboard already completed', function () {
+      const { active } = this.vouchers;
+
+      cy.fixture('realm').then((realm) => {
+        cy.intercept(
+          'POST',
+          `/pairing/v1/${realm.name}/fdo/ownership_vouchers/${active.guid}/to0`,
+          {
+            statusCode: 409,
+            body: {
+              errors: { detail: 'Device Onboard has already completed for the voucher' },
+            },
+          },
+        ).as('runTo0');
+
+        voucherRow(active.guid).contains('Re-run TO0').click();
+        cy.wait('@runTo0');
+
+        cy.get('.alert-danger')
+          .should('contain', 'Device Onboard has already completed')
+          .and('contain', active.guid);
+      });
+    });
+  });
+});
=== astarte-dashboard/package.json
@@ -1,7 +1,7 @@
 {
   "type": "module",
   "name": "astarte-dashboard",
-  "version": "1.4.0-rc.5",
+  "version": "1.4.0-rc.6",
   "description": "Astarte dashboard",
   "keywords": ["astarte", "dashboard", "iot"],
   "contributors": [
=== astarte-dashboard/src/DeviceStatusPage/DeviceInfoCard.tsx
@@ -48,6 +48,11 @@ const DeviceStatus = ({ status }: DeviceStatusProps): React.ReactElement => {
       icon = 'statusInDeletion' as const;
       break;
 
+    case 'unknown_status':
+      statusString = 'Unknown status';
+      icon = 'statusUnknown' as const;
+      break;
+
     case 'never_connected':
     default:
       statusString = 'Never connected';

$ cd /root/astrate-mule && gh api 'repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6' --jq '.files[] | select(.filename|test("astarte-dashboard/src")) | "=== \(.filename)\n\(.patch)"' 2>&1 | tail -120
+        <td colSpan={3}></td>
+      )}
+    </tr>
+  );
+};
+
+interface IntrospectionTableProps {
+  interfaces: Map<AstarteInterfaceDescriptor['name'], AstarteInterfaceDescriptor>;
+  onAddInterface: (interfaceDescriptor: AstarteInterfaceDescriptor) => void;
+  onRemoveInterface: (interfaceDescriptor: AstarteInterfaceDescriptor) => void;
+}
+
+const IntrospectionTable = ({
+  interfaces,
+  onAddInterface,
+  onRemoveInterface,
+}: IntrospectionTableProps): React.ReactElement => (
+  <Table className="mb-4" responsive>
+    <thead>
+      <tr>
+        <th>Interface name</th>
+        <th>Major</th>
+        <th>Minor</th>
+        <th className="action-column"> </th>
+      </tr>
+    </thead>
+    <tbody>
+      {Array.from(interfaces).map(([key, interfaceDescriptor]) => (
+        <InterfaceIntrospectionRow
+          key={key}
+          interfaceDescriptor={interfaceDescriptor}
+          onRemove={() => onRemoveInterface(interfaceDescriptor)}
+        />
+      ))}
+      <IntrospectionControlRow onAddInterface={onAddInterface} interfaces={interfaces} />
+    </tbody>
+  </Table>
+);
+
+export default IntrospectionTable;
=== astarte-dashboard/src/hooks/useFdo.ts
@@ -1,4 +1,5 @@
 import { useState, useCallback } from 'react';
+import type { AstarteDevice, AstarteInterfaceDescriptor } from 'astarte-client';
 import { useAstarte } from '../AstarteManager';
 
 interface UploadState {
@@ -9,12 +10,16 @@ interface UploadState {
 export const useFdo = () => {
   const { client } = useAstarte();
   const [state, setState] = useState<UploadState>({ status: 'idle', error: null });
+  const [deleteState, setDeleteState] = useState<UploadState>({ status: 'idle', error: null });
+  const [to0State, setTo0State] = useState<UploadState>({ status: 'idle', error: null });
 
   const uploadVoucher = useCallback(
     async (
+      hwId: AstarteDevice['id'],
       keyName: string,
       voucherText: string,
       options?: {
+        initialIntrospection?: { [interfaceName: string]: AstarteInterfaceDescriptor };
         keyAlgorithm?: string;
         replacementGuid?: string;
         replacementRvInfo?: string;
@@ -24,7 +29,7 @@ export const useFdo = () => {
       setState({ status: 'loading', error: null });
 
       try {
-        const response = await client.uploadFdoVoucher(keyName, voucherText, options);
+        const response = await client.uploadFdoVoucher(hwId, keyName, voucherText, options);
         setState({ status: 'success', error: null });
         return response;
       } catch (err: any) {
@@ -35,9 +40,46 @@ export const useFdo = () => {
     [client],
   );
 
+  const deleteVoucher = useCallback(
+    async (guid: string) => {
+      setDeleteState({ status: 'loading', error: null });
+
+      try {
+        await client.deleteFdoVoucher(guid);
+        setDeleteState({ status: 'success', error: null });
+      } catch (err: any) {
+        setDeleteState({ status: 'error', error: err });
+        throw err;
+      }
+    },
+    [client],
+  );
+
+  const runTo0 = useCallback(
+    async (guid: string) => {
+      setTo0State({ status: 'loading', error: null });
+
+      try {
+        const { expiry } = await client.runFdoVoucherTo0(guid);
+        setTo0State({ status: 'success', error: null });
+        return expiry;
+      } catch (err: any) {
+        setTo0State({ status: 'error', error: err });
+        throw err;
+      }
+    },
+    [client],
+  );
+
   return {
     uploadVoucher,
     status: state.status,
     error: state.error,
+    deleteVoucher,
+    deleteStatus: deleteState.status,
+    deleteError: deleteState.error,
+    runTo0,
+    to0Status: to0State.status,
+    to0Error: to0State.error,
   };
 };

$ cd /root/astrate-mule && gh api 'repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6' --jq '.files[] | select(.filename|test("astarte_data_access/lib/astarte_data_access/device|devices/device.ex|unconfirmed_device.ex")) | "=== \(.filename)\n\(.patch)"' 2>&1 | head -140
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

$ cd /root/astrate-mule && gh api 'repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6' --jq '.files[] | select(.filename|test("models/Device/index.ts|DevicesPage.tsx")) | "=== \(.filename)\n\(.patch)"' 2>&1 | head -100
=== astarte-dashboard/src/DevicesPage.tsx
@@ -119,11 +119,16 @@ const DeviceRow = ({ device, filters }: DeviceRowProps): React.ReactElement => {
     iconTooltip = 'Connected';
     statusLabel = 'Connected';
     lastEvent = `Connected on ${(device.lastConnection as Date).toLocaleString()}`;
-  } else if (device.lastConnection) {
+  } else if (device.lastConnection && device.lastDisconnection) {
     icon = 'statusDisconnected' as const;
     iconTooltip = 'Disconnected';
     statusLabel = 'Disconnected';
     lastEvent = `Disconnected on ${(device.lastDisconnection as Date).toLocaleString()}`;
+  } else if (device.lastConnection) {
+    icon = 'statusUnknown' as const;
+    iconTooltip = 'Unknown';
+    statusLabel = 'Unknown status';
+    lastEvent = `Connected on ${(device.lastConnection as Date).toLocaleString()}`;
   } else {
     icon = 'statusNeverConnected' as const;
     iconTooltip = 'Never connected';
=== astarte-dashboard/src/astarte-client/models/Device/index.ts
@@ -129,7 +129,12 @@ const astarteDeviceObjectSchema = yup
   })
   .required();
 
-type AstarteDeviceStatus = 'never_connected' | 'connected' | 'disconnected' | 'in_deletion';
+type AstarteDeviceStatus =
+  | 'never_connected'
+  | 'connected'
+  | 'disconnected'
+  | 'in_deletion'
+  | 'unknown_status';
 
 export class AstarteDevice {
   id: string;
@@ -205,6 +210,9 @@ export class AstarteDevice {
     if (this.lastConnection == null) {
       return 'never_connected';
     }
+    if (this.lastConnection && !this.lastDisconnection && !this.isConnected) {
+      return 'unknown_status';
+    }
     return this.isConnected ? 'connected' : 'disconnected';
   }
 }

$ cd /root/astrate-mule && rg -n 'last_disconnection|last_connection' docs/api/astarte_appengine_api.yaml | head -20
79:                        last_connection: "2024-05-30 13:49:57.045Z"
81:                        last_disconnection: null
1539:        last_connection:
1547:        last_disconnection:

$ cd /root/astrate-mule && rg -n 'last_disconnection|last_disconnection|last_connection|lastConnection' --glob '*.go' --glob '*.yaml' -l | head -20
docs/api/astarte_appengine_api.yaml
internal/broker/broker_e2e_test.go
internal/store/devices.go
internal/appengine/service.go

$ cd /root/astrate-mule && rg -n 'last_disconnection' internal/appengine/service.go internal/store/devices.go docs/api/astarte_appengine_api.yaml
internal/appengine/service.go:118:	LastDisconnection    *time.Time                    `json:"last_disconnection"`
internal/store/devices.go:63:	last_credentials_request_ip, last_connection, last_disconnection, last_seen_ip, connected,
internal/store/devices.go:311:		SET connected = false, last_disconnection = $3
docs/api/astarte_appengine_api.yaml:81:                        last_disconnection: null
docs/api/astarte_appengine_api.yaml:1547:        last_disconnection:

$ cd /root/astrate-mule && gh api 'repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6' --jq '.files[] | select(.filename|test("astarte_pairing_web/router.ex|agent/agent.ex|agent/device_registration_request.ex|views/error_view.ex|fallback_controller.ex")) | "=== \(.filename)\n\(.patch)"' 2>&1 | head -160
=== apps/astarte_pairing/lib/astarte_pairing/agent/agent.ex
@@ -26,31 +26,15 @@ defmodule Astarte.Pairing.Agent do
   alias Astarte.Pairing.Agent.DeviceRegistrationResponse
   alias Astarte.Pairing.Engine
 
-  alias Astarte.Core.Triggers.SimpleEvents.DeviceRegisteredEvent
-  alias Astarte.Events.Triggers
-  alias Astarte.Events.TriggersHandler
-
   def register_device(realm, attrs \\ %{}) do
     changeset =
       %DeviceRegistrationRequest{}
       |> DeviceRegistrationRequest.changeset(attrs)
 
-    with {:ok,
-          %DeviceRegistrationRequest{hw_id: hw_id, initial_introspection: initial_introspection}} <-
-           Ecto.Changeset.apply_action(changeset, :insert),
-         initial_introspection =
-           Enum.map(initial_introspection, fn {interface_name,
-                                               %{"major" => major, "minor" => minor}} ->
-             %{
-               interface_name: interface_name,
-               major_version: major,
-               minor_version: minor
-             }
-           end),
-         {:ok, credentials_secret} <-
-           Engine.register_device(realm, hw_id, initial_introspection: initial_introspection) do
-      dispatch_device_registration_trigger(realm, hw_id)
-
+    with {:ok, request} <- Ecto.Changeset.apply_action(changeset, :insert),
+         %{hw_id: hw_id, initial_introspection: initial_introspection} = request,
+         opts = [initial_introspection: initial_introspection],
+         {:ok, credentials_secret} <- Engine.register_device(realm, hw_id, opts) do
       {:ok, %DeviceRegistrationResponse{credentials_secret: credentials_secret}}
     end
   end
@@ -60,29 +44,4 @@ defmodule Astarte.Pairing.Agent do
       Engine.unregister_device(realm, device_id)
     end
   end
-
-  defp dispatch_device_registration_trigger(realm_name, hw_id) do
-    timestamp = DateTime.utc_now() |> DateTime.to_unix(:millisecond)
-    event_key = :on_device_registered
-    event_type = :device_registered_event
-    {:ok, device_id} = Device.decode_device_id(hw_id, allow_extended_id: true)
-
-    Triggers.find_device_trigger_targets(realm_name, device_id, event_key)
-    |> dispatch_all(realm_name, hw_id, timestamp, event_type, %DeviceRegisteredEvent{})
-  end
-
-  defp dispatch_all(targets, realm_name, device_id, timestamp, event_type, event) do
-    targets
-    |> Enum.map(fn {target, policy} ->
-      TriggersHandler.dispatch_event(
-        event,
-        event_type,
-        target,
-        realm_name,
-        device_id,
-        timestamp,
-        policy
-      )
-    end)
-  end
 end
=== apps/astarte_pairing/lib/astarte_pairing/agent/device_registration_request.ex
@@ -22,7 +22,7 @@ defmodule Astarte.Pairing.Agent.DeviceRegistrationRequest do
 
   import Ecto.Changeset
 
-  alias Astarte.Core.Device
+  alias Astarte.FDO.OwnershipVoucher.LoadRequest
   alias Astarte.Pairing.Agent.DeviceRegistrationRequest
 
   @primary_key false
@@ -36,36 +36,7 @@ defmodule Astarte.Pairing.Agent.DeviceRegistrationRequest do
     request
     |> cast(attrs, [:hw_id, :initial_introspection])
     |> validate_required([:hw_id])
-    |> validate_hw_id(:hw_id)
-    |> validate_change(:initial_introspection, &validate_introspection/2)
-  end
-
-  defp validate_hw_id(changeset, field) do
-    with {:ok, hw_id} <- fetch_change(changeset, field),
-         {:ok, _decoded_id} <- Device.decode_device_id(hw_id, allow_extended_id: true) do
-      changeset
-    else
-      # No hw_id, already handled
-      :error ->
-        changeset
-
-      _ ->
-        add_error(changeset, field, "is not a valid base64 encoded 128 bits id")
-    end
-  end
-
-  defp validate_introspection(field, introspection) when is_map(introspection) do
-    Enum.reduce(introspection, [], fn
-      {interface_name, %{"major" => major, "minor" => minor}}, acc
-      when is_integer(major) and is_integer(minor) ->
-        if major < 0 or minor < 0 do
-          [{field, "has negative versions in interface #{interface_name}"}]
-        else
-          acc
-        end
-
-      {interface_name, _}, acc ->
-        [{field, "has invalid format for interface #{interface_name}"} | acc]
-    end)
+    |> LoadRequest.validate_hw_id(:hw_id)
+    |> validate_change(:initial_introspection, &LoadRequest.validate_introspection/2)
   end
 end
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
=== apps/astarte_pairing/lib/astarte_pairing_web/router.ex
@@ -36,8 +36,6 @@ defmodule Astarte.PairingWeb.Router do
   pipeline :fdo do
     plug :accepts, ["cbor"]
     plug :put_view, Astarte.PairingWeb.FDOView
-    plug Astarte.PairingWeb.Plug.LogRealm
-    plug Astarte.PairingWeb.Plug.VerifyRealmExists
     plug Astarte.PairingWeb.Plug.SetupFDO
   end

$ cd /root/astrate-mule && rg -rn 'device_registered|DeviceRegistered|on_device_registered' --glob '*.go' -l | head
cmd/astrate/main.go
internal/flow/blocks/virtualdevicepool/virtualdevicepool_test.go
internal/flow/blocks/virtualdevicepool/virtualdevicepool.go
internal/flow/factory.go
internal/engine/stream/bus.go
internal/engine/engine.go
internal/engine/triggers/match.go
internal/engine/triggers/events_test.go
internal/engine/triggers/events.go
internal/pairing/service.go

$ cd /root/astrate-mule && gh search code --owner astarte-platform 'on_device_registered' --limit 20 --json path,repository 2>&1 | head -40
[{"path":"apps/astarte_pairing/lib/astarte_pairing/engine.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"}},{"path":"apps/astarte_events/lib/astarte_events/triggers/core.ex","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"}},{"path":"apps/astarte_events/test/astarte_events/triggers/core_test.exs","repository":{"id":"MDEwOlJlcG9zaXRvcnkxMjc0MzAxNDE=","isFork":false,"isPrivate":false,"nameWithOwner":"astarte-platform/astarte","url":"https://github.com/astarte-platform/astarte"}}]

$ cd /root/astrate-mule && for ref in v1.4.0-rc.5 v1.4.0-rc.6; do echo "=== $ref"; gh api "repos/astarte-platform/astarte/contents/apps/astarte_pairing/lib/astarte_pairing/engine.ex?ref=$ref" --jq '.content' | base64 -d | rg -n 'device_registered|Triggers|DeviceRegisteredEvent|def register_device|def unregister' | head -20; done
=== v1.4.0-rc.5
125:  def register_device(realm, hardware_id, opts \\ []) do
184:  def unregister_device(realm, encoded_device_id) do
=== v1.4.0-rc.6
25:  alias Astarte.Core.Triggers.SimpleEvents.DeviceRegisteredEvent
27:  alias Astarte.Events.Triggers
28:  alias Astarte.Events.TriggersHandler
130:  def register_device(realm, hardware_id, opts \\ []) do
188:    event_key = :on_device_registered
189:    event_type = :device_registered_event
192:    Triggers.find_device_trigger_targets(realm_name, device_id, event_key)
193:    |> dispatch_all(realm_name, hw_id, timestamp, event_type, %DeviceRegisteredEvent{})
199:      TriggersHandler.dispatch_event(
247:  def unregister_device(realm, encoded_device_id) do

$ cd /root/astrate-mule && gh api "repos/astarte-platform/astarte/contents/apps/astarte_pairing/lib/astarte_pairing/engine.ex?ref=v1.4.0-rc.6" --jq '.content' | base64 -d | sed -n '125,215p'
end
  end

  @spec register_device(String.t(), Device.encoded_device_id(), keyword()) ::
          {:ok, nil} | {:ok, String.t()} | {:error, term()}
  def register_device(realm, hardware_id, opts \\ []) do
    Logger.debug(
      "register_device request for device #{inspect(hardware_id)} in realm #{inspect(realm)}"
    )

    :telemetry.execute([:astarte, :pairing, :register_new_device], %{}, %{realm: realm})

    opts =
      Keyword.update(
        opts,
        :initial_introspection,
        [],
        &Enum.map(&1, fn {interface_name, %{"major" => major, "minor" => minor}} ->
          %{
            interface_name: interface_name,
            major_version: major,
            minor_version: minor
          }
        end)
      )

    with_credentials? = Keyword.get(opts, :with_credentials?, true)

    with {:ok, device_id} <- Device.decode_device_id(hardware_id, allow_extended_id: true),
         :ok <- verify_can_register_device(realm, device_id),
         {credentials_secret, secret_hash} = generate_credentials_secret(with_credentials?),
         {:ok, _device} <-
           DataAccessDevice.register(realm, device_id, hardware_id, secret_hash, opts) do
      dispatch_device_registration_trigger(realm, hardware_id)

      {:ok, credentials_secret}
    else
      {:error, :shutdown} ->
        {:error, :realm_not_found}

      {:error, reason} ->
        {:error, reason}
    end
  end

  def add_unconfirmed_credentials(realm_name, device_id) do
    {credentials_secret, secret_hash} = generate_credentials_secret(true)

    with :ok <- DataAccessDevice.add_unconfirmed_credentials(realm_name, device_id, secret_hash) do
      {:ok, credentials_secret}
    end
  end

  defp generate_credentials_secret(false = _with_credentials?), do: {nil, nil}

  defp generate_credentials_secret(true = _with_credentials?) do
    credentials_secret = CredentialsSecret.generate()
    secret_hash = CredentialsSecret.hash(credentials_secret)
    {credentials_secret, secret_hash}
  end

  defp dispatch_device_registration_trigger(realm_name, hw_id) do
    timestamp = DateTime.utc_now() |> DateTime.to_unix(:millisecond)
    event_key = :on_device_registered
    event_type = :device_registered_event
    {:ok, device_id} = Device.decode_device_id(hw_id, allow_extended_id: true)

    Triggers.find_device_trigger_targets(realm_name, device_id, event_key)
    |> dispatch_all(realm_name, hw_id, timestamp, event_type, %DeviceRegisteredEvent{})
  end

  defp dispatch_all(targets, realm_name, device_id, timestamp, event_type, event) do
    targets
    |> Enum.map(fn {target, policy} ->
      TriggersHandler.dispatch_event(
        event,
        event_type,
        target,
        realm_name,
        device_id,
        timestamp,
        policy
      )
    end)
  end

  defp verify_can_register_device(realm_name, device_id) do
    if Queries.check_already_registered_device(realm_name, device_id) do
      # An already existing device should always be able to retrieve a new credentials secret
      :ok
    else

$ cd /root/astrate-mule && rg -n 'device_registered|DeviceRegistered' internal/pairing/service.go internal/engine/triggers/events.go | head -20
internal/pairing/service.go:123:// engine can emit a DeviceRegisteredEvent on the trigger bus. The callback
internal/engine/triggers/events.go:181:// DeviceRegisteredEvent is the device_registered event body.
internal/engine/triggers/events.go:182:type DeviceRegisteredEvent struct {
internal/engine/triggers/events.go:183:	// Type is always "device_registered".
internal/engine/triggers/events.go:187:// NewDeviceRegisteredEvent builds a device_registered event body.
internal/engine/triggers/events.go:188:func NewDeviceRegisteredEvent() DeviceRegisteredEvent {
internal/engine/triggers/events.go:189:	return DeviceRegisteredEvent{Type: OnDeviceRegistered}

$ cd /root/astrate-mule && gh api "repos/astarte-platform/astarte/contents/libs/astarte_fdo/lib/ownership_voucher/load_request.ex?ref=v1.4.0-rc.6" --jq '.content' | base64 -d | rg -n 'add_error|def validate'
114:      {:ok, _old_voucher} -> add_error(changeset, voucher, "guid has already been claimed")
127:      {:ok, _device} -> add_error(changeset, device_id_field, "already exists")
140:  def validate_hw_id(changeset, _field) when not changeset.valid?, do: changeset
142:  def validate_hw_id(changeset, field) do
155:        {:error, add_error(changeset, field, "is not a valid base64 encoded 128 bits id")}
159:  def validate_introspection(field, introspection) when is_map(introspection) do
200:          :error -> add_error(changeset, :replacement_guid, "is not valid base64")
216:            add_error(
233:          :error -> add_error(changeset, :replacement_public_key, "is not a valid PEM public key")
263:        add_error(changeset, :ownership_voucher, "is not a valid ownership voucher")
279:          add_error(
287:        add_error(changeset, :ownership_voucher, "has an unsupported key type")
300:      {:error, _} -> add_error(changeset, :key_name, "does not exist in secrets store")
332:      add_error(

$ cd /root/astrate-mule && rg -n 'DeviceRegistered|onRegistered|RegisterCallback|cb ' internal/pairing/service.go | head -20 && echo '---' && sed -n '105,150p' internal/pairing/service.go
123:// engine can emit a DeviceRegisteredEvent on the trigger bus. The callback
---
	// (e.g. "mqtts://host:8883").
	BrokerURL string
	// CertTTL is the client certificate validity; zero selects
	// ca.DefaultCertTTL (30 days).
	CertTTL time.Duration
	// EnforceLatestCert enables the always-online-CRL behaviour
	// (docs/DESIGN.md §4.3): verify reports REVOKED for certificates whose
	// serial differs from the device's latest issuance.
	EnforceLatestCert bool
	// Version is reported by the info endpoint; empty selects
	// DefaultVersion.
	Version string
	// BcryptCost hashes credentials secrets; zero selects
	// bcrypt.DefaultCost (10, docs/DESIGN.md §4.1).
	BcryptCost int
}

// OnRegisteredFunc is called after a successful device registration so the
// engine can emit a DeviceRegisteredEvent on the trigger bus. The callback
// receives the realm name, encoded device ID, and the registration instant.
type OnRegisteredFunc func(realmName string, deviceID string, at time.Time)

// Service implements pairing flows A–C over a Store and the per-realm CA.
type Service struct {
	st     Store
	sealer *store.KeySealer
	cfg    Config
	now    func() time.Time
	// OnRegistered, when non-nil, is called after a successful registration.
	OnRegistered OnRegisteredFunc
}

// New builds a pairing Service. The sealer opens realms' AES-GCM-sealed CA
// private keys (docs/DESIGN.md §4.3); it is owned by the caller and shared
// with housekeeping (which seals new realm CAs).
func New(st Store, sealer *store.KeySealer, cfg Config) *Service {
	if cfg.Version == "" {
		cfg.Version = DefaultVersion
	}
	if cfg.BcryptCost == 0 {
		cfg.BcryptCost = bcrypt.DefaultCost
	}
	if cfg.CertTTL <= 0 {
		cfg.CertTTL = ca.DefaultCertTTL
	}
	return &Service{st: st, sealer: sealer, cfg: cfg, now: time.Now}

$ cd /root/astrate-mule && rg -n 'OnRegistered' -A6 internal/pairing/service.go | head -40
122:// OnRegisteredFunc is called after a successful device registration so the
123-// engine can emit a DeviceRegisteredEvent on the trigger bus. The callback
124-// receives the realm name, encoded device ID, and the registration instant.
125:type OnRegisteredFunc func(realmName string, deviceID string, at time.Time)
126-
127-// Service implements pairing flows A–C over a Store and the per-realm CA.
128-type Service struct {
129-	st     Store
130-	sealer *store.KeySealer
131-	cfg    Config
--
133:	// OnRegistered, when non-nil, is called after a successful registration.
134:	OnRegistered OnRegisteredFunc
135-}
136-
137-// New builds a pairing Service. The sealer opens realms' AES-GCM-sealed CA
138-// private keys (docs/DESIGN.md §4.3); it is owned by the caller and shared
139-// with housekeeping (which seals new realm CAs).
140-func New(st Store, sealer *store.KeySealer, cfg Config) *Service {
--
221:	if s.OnRegistered != nil {
222:		s.OnRegistered(realmName, hwID, time.Now())
223-	}
224-	return secret, nil
225-}
226-
227-// Unregister implements the flow A DELETE: the device becomes registrable
228-// again, its data is kept (store.UnregisterDevice clears only the

$ cd /root/astrate-mule && sed -n '185,225p' internal/pairing/service.go
return "", err
	}

	// Upstream parity: the limit gates every registration attempt, secret
	// rotations for existing devices included.
	if realm.DeviceRegistrationLimit != nil {
		n, err := s.st.CountDevices(ctx, realm.ID)
		if err != nil {
			return "", err
		}
		if n >= int64(*realm.DeviceRegistrationLimit) {
			return "", ErrRegistrationLimitReached
		}
	}

	secret, err := generateSecret()
	if err != nil {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), s.cfg.BcryptCost)
	if err != nil {
		return "", fmt.Errorf("pairing: hashing credentials secret: %w", err)
	}

	if err := s.st.RegisterDevice(ctx, realm.ID, id, string(hash)); err != nil {
		if errors.Is(err, store.ErrDeviceAlreadyConfirmed) {
			return "", fmt.Errorf("%w: %s", ErrAlreadyRegistered, hwID)
		}
		return "", err
	}

	if initialFormat != "" {
		if err := s.st.SetPayloadFormatHint(ctx, realm.ID, id, initialFormat); err != nil {
			return "", err
		}
	}
	if s.OnRegistered != nil {
		s.OnRegistered(realmName, hwID, time.Now())
	}
	return secret, nil
}

$ cd /root/astrate-mule && rg -n 'compat-note|COMPATIBILITY' .mule/for-giulio.md | head -20
38:  **Standing ask, unchanged, for whenever a stable v1.4.0 ships:** answer #92, reconcile the two experimental rows, run the final-phase `APICompatVersion` bump (and update `docs/COMPATIBILITY.md` deviation #10 to match), then cut the v3.0 tag. One process note while this stays blocked: this is the fifteenth milestone-run entry in this file (fourteen prior, 2026-09-06 → 2026-09-24) and the header asks for a queue rather than a log — worth deleting the superseded duplicates whenever you next prune, or the signal-to-noise here keeps falling.
146:- **COMPATIBILITY.md §8 wording proposal: trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (Forwarder) actions are single-shot by design, receiving only `maximum_capacity`/`event_ttl`.** Propose adding to the "Trigger delivery policies" deviation: *"Policy `error_handlers` (retry/discard) and the `retry_times` attempt cap govern HTTP webhook actions only: a custom action delivered through the Forwarder seam is single-shot by design, receiving from the policy solely its `maximum_capacity` in-flight bound and its `event_ttl` staleness drop, with a forward failure counted once and never retried."* Verified in `internal/engine/triggers/actions.go`: `Action.Custom` routes to `forward` (actions.go:435-439), which never consults `policy.Decide` — the Forwarder seam returns only an error, no HTTP status (actions.go:161-167, doc comment at actions.go:471-479) — while `Enqueue` bounds `maximum_capacity` on every delivery (actions.go:338-349) and `deliver` drops past `event_ttl` before the custom/webhook branch (actions.go:400-416). The existing warn-on-single-shot line already tells the operator when a policy's `retry_times` cannot apply to the forward path (actions.go:491-497).
166:- **COMPATIBILITY.md wording update for upstream v1.3.4 (newest stable, 2026-09-18; v1.4.0 is still rc.5-only) — folds v1.3.3 and v1.3.4 into the open v1.3.2 proposal below; both are maintenance-only and wire-inert, so that proposal stands unchanged.** v1.3.3 (2026-08-07) was an empty-body patch and v1.3.4 (2026-09-18, `astarte_data_updater_plant` only — more frequent fullsweep GC on AMQPDataConsumer processes plus an RPC-availability fix; release body via `gh api repos/astarte-platform/astarte/releases` `v1.3.4`) introduce **no wire/API surface change**, so the v1.3.2 wording proposal stands complete: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"). The only delta a v1.3.4-aware doc carries is the version reference: the proposed §Infrastructure-differences sentence's "until the milestone that adopts v1.3.2 as the target" reads "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target". Raw: [v1.3.4](https://github.com/astarte-platform/astarte/releases/tag/v1.3.4).
319:- **COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable, 2026-07-14; v1.4.0 is still rc-only).** Astrate's doc and `APICompatVersion` still target upstream **1.2.2** (`internal/realm/service.go:588`); v1.3.0 (2026-05-06) introduced wire-surface changes Astrate does not yet emulate, so this is a decision — adopt v1.3.2 as the compatibility target (then update the doc + bump `APICompatVersion` together, per the bump rule) or keep 1.2.2 and add a "not yet emulated" note. Wire-relevant v1.3.0 deltas (release notes): **MQTT v1 capabilities** incl. `purge_properties_compression_format` (plaintext vs zlib purge — touches the `emptyCache`/`producer,properties` contract COMPATIBILITY.md deviation 1 documents); **empty introspection now allowed**; **device registration triggers** (pairing) and **device deletion started/completed triggers** (RM — the latter two already exist as Astrate deviation 9 emits both around the synchronous delete); **FDO authentication** (pairing, disabled by default); **realm-scoped health** — upstream v1.3 added `GET /pairing/v1/{realm}/health`, which Astrate already serves (`internal/pairing/http.go:78-81`, comment already says "upstream 1.3+"), so deviation 18's wording ("which upstream 404s") is now false against 1.3 and the note should be reworded either way. Proposed doc wording (for your approval, edit to taste): in §Infrastructure differences add a sentence — *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl.

→ Read .mule/config

$ cd /root/astrate-mule && gh api 'repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6' --jq '.files[] | select(.filename|test("realm_management/device_removal/(core|device_remover|scheduler|queries).ex")) | "=== \(.filename)\n\(.patch)"' 2>&1 | rg -v '^[-+ ].*voucher' | head -70
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
+
+  N.B.: only the database row is removed.
+  """
+
+    _ =
+      Logger.info(
+        realm: realm_name,
+      )
+
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
=== apps/astarte_realm_management/lib/astarte_realm_management/device_removal/device_remover.ex
@@ -27,20 +27,25 @@ defmodule Astarte.RealmManagement.DeviceRemoval.DeviceRemover do
   use Task
   require Logger
   alias Astarte.Core.Device
+  alias Astarte.DataAccess.Device, as: DeviceQueries
   alias Astarte.RealmManagement.DeviceRemoval.Core
 
   @spec run(%{:device_id => <<_::128>>, :realm_name => binary()}) :: :ok | no_return()
   def run(%{realm_name: realm_name, device_id: device_id}) do
     encoded_device_id = Device.encode_device_id(device_id)
     _ = Logger.info("Starting to remove device #{encoded_device_id}", tag: "device_delete_start")
 
-    Core.delete_individual_datastreams!(realm_name, device_id)
-    Core.delete_individual_properties!(realm_name, device_id)
-    Core.delete_object_datastream!(realm_name, device_id)
-    Core.delete_aliases!(realm_name, device_id)
-    Core.delete_groups!(realm_name, device_id)
-    Core.delete_kv_store_entries!(realm_name, encoded_device_id)
-    Core.delete_device!(realm_name, device_id)
+    with {:ok, device} <- DeviceQueries.fetch(realm_name, device_id) do
+      Core.delete_individual_datastreams!(realm_name, device_id)
+      Core.delete_individual_properties!(realm_name, device_id)
+      Core.delete_object_datastream!(realm_name, device_id)
+      Core.delete_aliases!(realm_name, device_id)
+      Core.delete_groups!(realm_name, device_id)
+      Core.delete_kv_store_entries!(realm_name, encoded_device_id)
+      Core.delete_device!(realm_name, device_id)

$ cd /root/astrate-mule && gh issue list --limit 100 --state open --json number,title --jq '.[] | "\(.number)\t\(.title)"' 2>&1 | rg -i 'fdo|voucher|1\.4|compat' | head -20
78	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)

$ cd /root/astrate-mule && gh issue list --limit 200 --state all --json number,title,state --jq '.[] | select(.title|test("FDO|voucher";"i")) | "\(.number)\t\(.state)\t\(.title)"' 2>&1 | head -20
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)

$ cd /root/astrate-mule && gh issue view 78 --json number,title,body,labels --jq '"\(.number) \(.title)\nlabels: \(.labels|map(.name)|join(","))\n\n\(.body)"' 2>&1 | head -80
78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
labels: enhancement,milestone-4.0,upstream-parity

Upstream 1.3 experimental, default-on in 1.4: TO2 CBOR endpoints under `/v1/{realm}/fdo/101/msg/{60..70}` (auto-registers the device mid-flow, hands over credentials_secret), owner side `POST/GET .../fdo/owner_keys` + ownership vouchers with TO0 claim, secrets in Vault/OpenBao, rendezvous server required.

Large surface; needs a product call (who runs a rendezvous server?) before any implementation. Parked.
Size L.\n\n---\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).

$ cd /root/astrate-mule && gh api "repos/astarte-platform/astarte/contents/libs/astarte_fdo/lib/ownership_voucher/ownership_voucher.ex?ref=v1.4.0-rc.6" --jq '.content' | base64 -d | rg -n 'def delete|realm_name|get_by_guid|where' | head -25
36:  def list(realm_name) do
37:    Queries.list_ownership_vouchers(realm_name)
47:  def delete(realm, guid) do
103:  defp register_on_rendezvous(realm_name, guid, voucher_cbor) do
105:         {:ok, owner_key} <- Secrets.get_key_for_guid(realm_name, guid) do
129:  defp ensure_device_voucher_deletion(realm_name, device_id) do
132:    case RealmManagement.delete_device(realm_name, encoded_device_id) do
139:  defp revoke_rendezvous_registration(realm_name, guid, voucher_cbor) do
141:         {:ok, owner_key} <- Secrets.get_key_for_guid(realm_name, guid),
169:    with {:ok, {realm_name, device_id, ownership_voucher_cbor}} <-
172:      {:ok, {realm_name, device_id, ownership_voucher}}

$ cd /root/astrate-mule && gh api "repos/astarte-platform/astarte/contents/libs/astarte_fdo/lib/ownership_voucher/ownership_voucher.ex?ref=v1.4.0-rc.6" --jq '.content' | base64 -d | sed -n '40,70p'
@doc """
  Deletes an ownership voucher.

  The corresponding registration on the FDO rendezvous server is revoked first;
  if that fails, the voucher is not deleted.
  """
  @spec delete(String.t(), binary()) :: :ok | {:error, term()}
  def delete(realm, guid) do
    with {:ok, ownership_voucher} <- Queries.fetch_ownership_voucher(guid),
         :ok <- ensure_voucher_in_realm(ownership_voucher.realm, realm),
         :ok <-
           revoke_rendezvous_registration(
             ownership_voucher.realm,
             guid,
             ownership_voucher.voucher_data
           ),
         :ok <-
           ensure_device_voucher_deletion(ownership_voucher.realm, ownership_voucher.device_id) do
      Queries.delete_ownership_voucher(guid)
    end
  end

  @doc """
  Re-runs TO0 with the rendezvous server for an ownership voucher whose device
  has not completed Device Onboard yet, refreshing its expiry.

  Returns the new expiry on success.
  """
  @spec run_to0(String.t(), binary()) :: {:ok, DateTime.t()} | {:error, term()}
  def run_to0(realm, guid) do
    with {:ok, ownership_voucher} <- Queries.fetch_ownership_voucher(guid),

$ cd /root/astrate-mule && tail -c 900 .mule/todo.md | cat -A | tail -5 | cut -c1-200
e 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with

$ cd /root/astrate-mule && tail -1 .mule/todo.md | tail -c 120 | cat -A
n docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.$

$ cd /root/astrate-mule && rg -c 'with the canonical-form check moved to `String`' .mule/todo.md
1

→ Read .mule/todo.md [offset=246, limit=1]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -243,4 +243,6 @@
 - [x] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` path parameter (docs/api/astarte_pairing_api.yaml:360) and `RegisterRequest.hw_id` (yaml:373) — but `deviceid.Parse` accepts only the 22-character **unpadded base64url** form, in strict mode, and rejects the standard alphabet's `+` and `/` and any `=` padding with an error wrapping `ErrInvalid` (pkg/deviceid/deviceid.go:26, 37, 39-60; upstream parity cited at 36). Measured: with the same `encoding` the parser uses, `dT6hS2W9TT6LEnP25ks+lg` and `.../lg` fail at byte 19 and the padded 23-character form fails at byte 21, so a client generated from the current spec can send a well-formed standard-base64 ID and get 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}` (internal/pairing/http.go:307-309). State the exact encoding and the 22-character length in both descriptions and add a `pattern`/`minLength`/`maxLength` to `hw_id`. The existing `hw_id: dT6hS2W9TT6LEnP25ks_lg` example (yaml:123) is valid and stays. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in docs/api/astarte_pairing_api.yaml (yaml:374-376), an unconstrained string described only as "Optional Astrate extension specifying the device's initial payload format" — so a generated client can send any string, while `Service.Register` accepts only `"bson"` or `"json"` and rejects everything else with 422 `{"errors":{"initial_payload_format":["is invalid"]}}` (internal/pairing/service.go:179-181 -> internal/pairing/http.go:310-312; the allowed set is stated at service.go:173). Add `enum: [bson, json]` plus a note that absent/empty means unset, and add the three missing `ValidationError` examples while you are in the component: `initial_payload_format: ["is invalid"]` (http.go:310-312), `csr: ["can't be blank"]` (http.go:173-177) and `client_crt: ["can't be blank"]` (http.go:253-257) — the component currently shows only the four hw_id/error_name rows (yaml:594-618) and the `csr`/`client_crt` blank-field cases are what a device actually hits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` description in docs/api/astarte_pairing_api.yaml (yaml:153), "Removes a device from the realm. Returns 204 on success.", which reads as a device-and-data deletion. The row and every byte of its stored data survive: `store.UnregisterDevice` only clears the credential trail (`credentials_secret_hash = ''`, `first_credentials_request = NULL`, `cert_serial = NULL`, `cert_aki = NULL`) and flips `status` back to `'registered'` (internal/store/devices.go:93-110), and the service docstring says so outright ("its data is kept", internal/pairing/service.go:227-229). State the real semantics: the device becomes registrable again, its interfaces, datastream data and group memberships are retained, and a second DELETE is a 204 no-op (the row is still there, so RowsAffected is 1, devices.go:106-108) — not the documented 404. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `base64.RawURLEncoding.Strict()` (pkg/deviceid/deviceid.go:37), so of the 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with `ErrInvalid` while `..._lA` is accepted, and the pairing spec's `hw_id` pattern `^[A-Za-z0-9_-]{22}$` (added by docs-sync-pairing-deviceid-base64url) deliberately does not encode the rule because a pattern is a client-side contract and upstream's leniency is unmeasured. The docstring justifies strictness as byte-identical round-tripping, which is a reason, not a parity claim, and docs/COMPATIBILITY.md records no entry for it. Measure upstream (Elixir `Base.url_decode64!(padding: false)`, which does not check those bits) on the Legion Go, then either record the divergence in docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.
+- [ ] fdo-rc6-scope-delta-for-giulio [auto]: upstream v1.4.0-rc.6 (2026-09-30, 27 commits, 154 files, almost all FDO) grows the onboarding surface issue #78 is scoped against, and none of it is in #78's body (written against upstream 1.3 / rc.5): new `POST /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}/to0` re-runs TO0 and returns the refreshed rendezvous expiry (`{"data":{"expiry":…}}`, ownership_voucher.ex:63-73; 409 `device_already_onboarded` via fallback_controller.ex:120-125); voucher deletion, which revokes the rendezvous registration first and *keeps* the voucher if that fails (`:rendezvous_revocation_failed` → 500) and refuses another realm's voucher through `ensure_voucher_in_realm` (ownership_voucher.ex:40-60 — now expressible because the voucher row carries its realm); `device_id` accepted on voucher upload; the device registered immediately on upload rather than mid-TO2 (`add_unconfirmed_credentials`, device.ex:63-85); single-entry x5chain certs accepted; invalid vouchers no longer stored; vouchers and to2-sessions moved from the per-realm keyspace to the global `astarte` one with `fdo_guid` added to the device row (housekeeping/realms/queries.ex:784, 1129-1133); and RM device deletion now cascading to the voucher (device_removal/core.ex). Write the delta into `.mule/for-giulio.md` as a scope update for #78 and say which of these a client could actually observe — do not edit or comment on the issue (GitHub state is read-only for the mule) and do not implement anything.
+- [ ] compat-note-v1.4-rc.6 [auto]: record in `.mule/for-giulio.md` that v1.4.0-rc.6 is wire-inert for Astrate, as a delta on the open v1.3.2/v1.3.4 wording proposal rather than a re-derivation. Measured, the whole non-FDO delta is inert: DUP RPC-server reliability is OTP supervisor/GC tuning; the Dashboard's new `unknown_status` is computed client-side from `last_connection && !last_disconnection && !isConnected` (models/Device/index.ts:213-215) and both timestamps are already on Astrate's device shape (internal/appengine/service.go:118), so a device whose disconnection was never recorded reads as "Unknown status" rather than "Disconnected"; the `on_device_registered` dispatch moved from `Agent.register_device` to `Engine.register_device` (pairing/engine.ex:157-165, 186-193) so it also covers the FDO path, but the Agent path is byte-for-byte the same rule and Astrate already fires it on every successful registration (internal/pairing/service.go:221-223); the `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"); and `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma. So the proposal's "v1.4.0 is still rc.5-only" reference becomes rc.6, and its "FDO authentication (pairing, disabled by default)" line understates the 1.4 line, which ships a full ownership-voucher surface Astrate does not emulate (#78, parked). Do not edit docs/COMPATIBILITY.md and do not restate the standing proposal.

$ cd /root/astrate-mule && wc -l .mule/todo.md && rg -n '^- \[ \] (fdo-rc6|compat-note-v1.4-rc.6)' .mule/todo.md | cut -c1-80 && git status --short
248 .mule/todo.md
247:- [ ] fdo-rc6-scope-delta-for-giulio [auto]: upstream v1.4.0-rc.6 (2026-09-3
248:- [ ] compat-note-v1.4-rc.6 [auto]: record in `.mule/for-giulio.md` that v1.
 M .mule/todo.md

Done: upstream is ahead (v1.4.0-rc.6, 2026-09-30 vs the 1.2.2 target in docs/COMPATIBILITY.md); read the rc.5→rc.6 diff and proposed two task lines — the FDO scope growth against issue #78, and the rc.6 delta on the standing v1.3.x wording proposal — because rc.6 is otherwise wire-inert for Astrate.
Files: .mule/todo.md (2 lines appended, nothing else touched; no git, no `gh issue create` — the recipe does not ask for one, so no evidence file either)
Verified: `gh api repos/astarte-platform/astarte/releases`, `.../compare/v1.4.0-rc.5...v1.4.0-rc.6` (+ per-file patches for housekeeping/pairing/RM/dashboard), `gh search code --owner astarte-platform on_device_registered` -> all four non-FDO claims checked against the diff and found already-parity or inert: DUP RPC reliability is OTP tuning; the Dashboard's `unknown_status` is client-side from `last_connection`/`last_disconnection` (Astrate returns both, internal/appengine/service.go:118); the `on_device_registered` dispatch moved Agent→Engine (engine.ex:157-165) so the Agent rule is unchanged and Astrate already fires it (internal/pairing/service.go:221-223); the `hw_id` error text is byte-identical after moving into `astarte_fdo`'s `LoadRequest`
Unsure: whether the rc.6 FDO additions should grow #78's scope or confirm it stays parked — that is the product call #78 already waits on, so both lines only put the delta in front of Giulio
Follow-ups: fdo-rc6-scope-delta-for-giulio, compat-note-v1.4-rc.6
