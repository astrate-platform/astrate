slug: recipe-astarte-upstream
verdict: proposed
at:  24ee222
ran: 2026-10-03T18:34:45Z on DietPi in 516s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"' 2>&1 | head -20
v1.4.0-rc.6	2026-09-30T07:28:27Z
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z

$ ls /root/astrate-mule/docs/ && echo "---TODO---" && wc -l /root/astrate-mule/.mule/todo.md && echo "---MULE---" && ls /root/astrate-mule/.mule/
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
UPSTREAM-EXPERIMENTAL.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui
---TODO---
263 /root/astrate-mule/.mule/todo.md
---MULE---
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
tasks
todo.md

$ cd /root/astrate-mule && rg -n "1\.4|1\.3\.|upstream version|targets" docs/COMPATIBILITY.md | head -30
335:17. **Always synchronous where upstream 1.4 defaults to asynchronous** —
336:    upstream 1.4 runs realm create/delete (Housekeeping), interface
348:    `docs/UPSTREAM-EXPERIMENTAL.md` (#68) until upstream 1.4 is final.

→ Read .mule/todo.md
→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ cd /root/astrate-mule && rg -n "rc\.|1\.4|compat-note|upstream" .mule/todo.md | tail -40 && echo "=== FOR-GIULIO tail ===" && tail -60 .mule/for-giulio.md && echo "=== wc ===" && wc -l .mule/for-giulio.md
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
124:- [x] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
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
247:- [x] fdo-rc6-scope-delta-for-giulio [auto]: upstream v1.4.0-rc.6 (2026-09-30, 27 commits, 154 files, almost all FDO) grows the onboarding surface issue #78 is scoped against, and none of it is in #78's body (written against upstream 1.3 / rc.5): new `POST /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}/to0` re-runs TO0 and returns the refreshed rendezvous expiry (`{"data":{"expiry":…}}`, ownership_voucher.ex:63-73; 409 `device_already_onboarded` via fallback_controller.ex:120-125); voucher deletion, which revokes the rendezvous registration first and *keeps* the voucher if that fails (`:rendezvous_revocation_failed` → 500) and refuses another realm's voucher through `ensure_voucher_in_realm` (ownership_voucher.ex:40-60 — now expressible because the voucher row carries its realm); `device_id` accepted on voucher upload; the device registered immediately on upload rather than mid-TO2 (`add_unconfirmed_credentials`, device.ex:63-85); single-entry x5chain certs accepted; invalid vouchers no longer stored; vouchers and to2-sessions moved from the per-realm keyspace to the global `astarte` one with `fdo_guid` added to the device row (housekeeping/realms/queries.ex:784, 1129-1133); and RM device deletion now cascading to the voucher (device_removal/core.ex). Write the delta into `.mule/for-giulio.md` as a scope update for #78 and say which of these a client could actually observe — do not edit or comment on the issue (GitHub state is read-only for the mule) and do not implement anything.
248:- [x] compat-note-v1.4-rc.6 [auto]: record in `.mule/for-giulio.md` that v1.4.0-rc.6 is wire-inert for Astrate, as a delta on the open v1.3.2/v1.3.4 wording proposal rather than a re-derivation. Measured, the whole non-FDO delta is inert: DUP RPC-server reliability is OTP supervisor/GC tuning; the Dashboard's new `unknown_status` is computed client-side from `last_connection && !last_disconnection && !isConnected` (models/Device/index.ts:213-215) and both timestamps are already on Astrate's device shape (internal/appengine/service.go:118), so a device whose disconnection was never recorded reads as "Unknown status" rather than "Disconnected"; the `on_device_registered` dispatch moved from `Agent.register_device` to `Engine.register_device` (pairing/engine.ex:157-165, 186-193) so it also covers the FDO path, but the Agent path is byte-for-byte the same rule and Astrate already fires it on every successful registration (internal/pairing/service.go:221-223); the `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"); and `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma. So the proposal's "v1.4.0 is still rc.5-only" reference becomes rc.6, and its "FDO authentication (pairing, disabled by default)" line understates the 1.4 line, which ships a full ownership-voucher surface Astrate does not emulate (#78, parked). Do not edit docs/COMPATIBILITY.md and do not restate the standing proposal.
255:- [!] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy alias fields to `InterfaceDefinition` and `InterfaceMapping` in docs/api/astarte_realm_management_api.yaml — both declare `additionalProperties: false` (yaml:935, 972) yet omit interface-level `quality` (alias of `ownership`) and `aggregate` (alias of `aggregation`) and mapping-level `path` (alias of `endpoint`), all three of which the parser accepts and canonicalises: decoded into dedicated fields (pkg/interfaceschema/parse.go:116-122, 136-138), normalised before any rule runs (parse.go:212-236), and re-encoded canonically by `ParseInterfaceCanonical` so the store never sees them (parse.go:165-170, 307-311). A spec-conformant client — or any generated client with this schema baked in — therefore rejects a valid upstream-shaped install/update with a spurious 422 while the service would have accepted and stored it as `ownership`/`aggregation`/`endpoint`. Add the three as documented deprecated aliases (enums mirroring ownership/aggregation, `path` a plain string) or relax `additionalProperties`, and record the three rules the code enforces: `ownership` + `quality` is a violation (parse.go:218) and `aggregation` + `aggregate` is a violation (parse.go:227), while `endpoint` + `path` silently keeps `endpoint` (parse.go:307-311). None of the three may appear in a `required` list. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
=== FOR-GIULIO tail ===
  (`internal/engine/triggers/match.go:11-12`). Decision deferred, tied to issue #17
  (group-WATCH-path reconciliation, trickle work, not mule): whatever group-membership
  mechanism comes out of that phase should also report the perf cost for this decision —
  noted in a comment on #17 so it isn't benchmarked twice. (Same survey, source 4.)

---

- ~~The Pi cannot run the race detector~~ — **resolved 2026-07-27** by installing Go 1.26.5
  as a userland toolchain on the Legion Go (`~/.local/go`, no root, `rm -rf` to undo). The
  Pi still cannot run `-race` (39-bit VMA kernel vs the 48 ThreadSanitizer needs), so its
  gate remains `go vet ./... && go test ./...` — but race coverage now exists on the Legion
  Go, where the full suite runs clean in ~40s on 16 cores. The standing `race-check` task is
  the concurrency gate. Concurrency work is queueable again, provided the race-check runs
  after it.
- ~~golangci-lint is not installed on the Pi~~ — **resolved 2026-07-28**: installed v2.12.2
  via `go install` (the prebuilt-binary installer's published sha256 for linux/arm64 did not
  verify, reproducibly, so built from source instead). `mule.service` now sets
  `MULE_LINT_CMD=golangci-lint run ./...` and has `/root/go/bin` on `PATH`; the lint gate
  runs starting with the next tick.
- ~~`/root/astrate` on the Pi has uncommitted work~~ — **resolved 2026-07-27** with the new
  `tools/reconcile.sh`: rescued onto `origin/wip/DietPi-20260727T171543Z` (pushed, not
  reviewed — read the diff before merging anything from it) and `/root/astrate` is now clean
  on `main`. Also had to set `commit.gpgsign false` locally on that clone first — it had
  signing on with no working gpg-agent for a non-interactive session, same class of problem
  `/root/astrate-mule` already solved the same way (your call, confirmed 2026-07-27: align it
  with the existing unattended-Pi-clone convention rather than fix gpg-agent or bypass signing
  silently inside the script).

## 2026-08-23 — FDO promoted to milestone-4.0 candidate (Giulio's decision, recorded)

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
=== wc ===
521 .mule/for-giulio.md

→ Read .mule/todo.md [offset=174]
Prior run already handled rc.6 (todo.md:247-248 are ticked). Verifying nothing is newer before concluding.

$ cd /root/astrate-mule && rg -n "rc\.6|rc\.5|1\.3\.4|1\.3\.2" .mule/for-giulio.md | head -30
22:  **Verified unchanged, no new gaps.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition, frozen by you on 2026-09-04 ("wait for a stable v1.4.0", an `-rc.N` does not satisfy it), is still unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Of #47–#89 **all are closed**; the open non-alarm set is unchanged — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** (untouched per standing instruction). rc.6 adds no v3.0 gap: it is already audited in the two entries below, and I re-measured its non-FDO half myself (`compare/v1.4.0-rc.5...v1.4.0-rc.6`, 27 commits) — outside FDO the only library changes are a `# coveralls-ignore-stop` line in `apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex:43` and `realms/queries.ex` *removing* the ownership_voucher / to2_sessions / unconfirmed_devices table creation (moved into migrations). No Realm Management, AppEngine or Pairing wire surface moved.
32:- **Delta on the two open COMPATIBILITY.md wording entries further down this file (the v1.3.4 one folding in v1.3.3, and the original v1.3.2 one): upstream `v1.4.0-rc.6` is wire-inert for Astrate, so nothing in that proposal is re-derived — only its version reference and one line of its "not yet emulated" list change.** Not applied: `docs/COMPATIBILITY.md` is untouched (still the 1.2.2 target, deviation 18's "which upstream 404s" at `docs/COMPATIBILITY.md:360`, `APICompatVersion = "1.2.2"` at `internal/realm/service.go:588`), and the FDO half of the rc.6 delta already has its own entry below rather than being folded in here.
36:  - **The two edits the proposal needs, both yours to approve:** (1) the 2026-09-18 entry's parenthetical "(v1.4.0 is still rc.5-only)" becomes "(v1.4.0 is still rc.6-only)" — re-verified this session that rc.6 is the newest 1.4 prerelease and there is still **no stable v1.4.0** (`gh api repos/astarte-platform/astarte/releases`); (2) its "**FDO authentication** (pairing, disabled by default)" item, and the same "experimental FDO pairing auth" phrase in the proposed doc sentence, understates the 1.4 line, which ships a full ownership-voucher surface (upload with mandatory `hw_id`, deletion, re-TO0, per-voucher `expiry`, a spec-mandated realm-less TO2 path) that Astrate does not emulate — #78, deliberately parked at `.mule/milestones.md:124`. "FDO ownership vouchers (1.4 line, not stable upstream)" says it without claiming the rc.5 shape.
37:  Unverified: the upstream side of the item list is the upstream-watch recipe's `v1.4.0-rc.5...v1.4.0-rc.6` comparison — the same measurement as the FDO entry below, taken there and not re-derived here (a raw fetch of `models/Device/index.ts` at the rc.6 tag did not resolve from this box). I re-verified only the Astrate-side claims above and the release list. No `.mule/todo.md` lines: the deliverable is the delta, and the proposal itself stays open.
41:- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published 2026-09-30T07:28Z, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and the one path #78's body states as fact no longer exists upstream. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix") deletes the `realm_name` path parameter from all seven TO2 `operation` specs in `fdo_onboarding_controller.ex`. **This breaking change is not in the rc.6 CHANGELOG entry** (which lists only the to0 route, voucher deletion, the `ASTARTE_BASE_URL_*` rename, device-deletion cascade, `device_id` on upload, immediate registration, the two x5chain/voucher-store fixes and the cross-realm delete fix) — so anyone tracking upstream by its changelog, including #78's body, misses it. The realm now comes from the voucher row: the global-keyspace query `fetch_device_id_and_ownership_voucher_and_realm/1` returns `{realm, device_id, voucher_data}` and the `:fdo` plug fetches the session by GUID alone (`plug/fdo_session.ex`: `Session.fetch(guid)`); `RvTO2Addr` no longer carries the realm in the hostname either ("the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure"). **Yours:** #78's acceptance is "a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end" — against rc.6 that SDK posts to `/fdo/101/msg/{60..70}` with no realm, so mirroring rc.5's shape would mean mirroring a path upstream has deleted. No action while parked; the body is just stale.
49:  - **Config rename, operator-visible:** `ASTARTE_BASE_URL_DOMAIN` → `ASTARTE_BASE_URL_HOST`, and the value may now be an IP address as well as a domain (`libs/astarte_fdo/lib/config.ex`, new typed `BaseURLHost` module). A deployment carrying the rc.5 variable name stops booting. This one matters for Astrate even parked, because it is the same key our config reference would have to name.
51:- **Not client-observable, so not part of the decision:** the keyspace move itself — vouchers and TO2 sessions now live in the global `astarte` keyspace (`Realm.astarte_keyspace_name()` throughout `fdo/queries.ex`) with a `realm` column on the voucher row and `ALLOW FILTERING` on the list query, so **realm isolation is now application-level rather than keyspace-level** (which is what makes `ensure_voucher_in_realm`-style cross-realm checks necessary, `b3bf731` "ensure users cannot delete ownership vouchers belonging to other realms"); `fdo_guid` added to the device row (migration `realm/0025`) so device removal can find the voucher; the `unconfirmed_devices` drop (`realm/0020`, schema file deleted, `remove_device_from_unconfirmed_devices!` removed); the dashboard voucher page (`astarte-dashboard/src/FdoVoucherPage.tsx`, `hooks/useFdo.ts`, a new cypress spec); and the `data_updater_plant` RPC/GC fixes plus the `v1.3.4` forward-ports.
52:- **One premise in #78's own comment is now false.** It says "upstream ships almost no public FDO documentation (code-level only); ours must exceed that." rc.6 adds two documented sections: `doc/pages/user/035-register_device.md` gained `## FIDO Device Onboard` (Owner Keys, Ownership Vouchers) and `## Credentials Secret Lifecycle`, and `doc/pages/architecture/050-pairing_mechanism.md` gained `## FIDO Device Onboarding` (TO0, TO2, Owner Keys management) — commits `162062e`, `d1c39fb`. "Docs must exceed upstream" is still the right requirement, but the gap is far smaller than #78 assumes, and those two pages are now the cheapest available scope input.
53:- **Where Astrate stands, for the size estimate:** it has **no FDO surface at all** — the only FDO reference in the tree is the comment on the realm health probe (`internal/pairing/http.go:79`, `internal/pairing/service.go:368`), and `.mule/milestones.md:124` still lists #78 as a deliberately parked item. So none of this is a delta to partially-built code; every item above is net-new scope for a "Size L" issue that is blocked on "who runs a rendezvous server?". The rc.6 delta does not change that call, but it adds four things any eventual implementation must decide: a spec-mandated device-side path with no realm (so any deviation is ours alone), a mandatory `hw_id` in the upload body (so the owner must supply the device id *before* the device has onboarded), an upload that creates a visible realm device as a side effect, and a per-voucher `expiry` the operator has to poll to know when the rendezvous stops serving.
54:  Unverified: read-only — every line above comes from the compare API and the rc.6-tagged sources, and the 409/500 mappings come from `FallbackController` source rather than a live call (no Docker on this box). Deliberately no `.mule/todo.md` lines: the deliverable is the scope update, and #78 stays parked until you answer the rendezvous question.
78:- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Status: not started", which is no longer true.** Release sweep (`gh api repos/astarte-platform/astarte/releases`, run this session) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.5` (2026-08-20, prerelease) — byte-identical to what the 2026-09-24 run saw, so nothing shipped upstream in the two days since. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Open non-alarm issues remain exactly **#93** (aclhook comment, `mule-review`, its own recipe path), **#92** (keyAgreement, parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** (FDO, `milestone-4.0`), **#1** (untouched per standing instruction). No new gaps, **no issues filed, no `.mule/todo.md` task lines added** — there is nothing machine-checkable left, and filing anything here would duplicate #92/#93/#78.
153:- **Milestone recipe run, 2026-09-24: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-23 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
161:- **Milestone recipe run, 2026-09-23: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-21 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
177:- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
185:- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
197:- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
201:- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
209:- **COMPATIBILITY.md wording update for upstream v1.3.4 (newest stable, 2026-09-18; v1.4.0 is still rc.5-only) — folds v1.3.3 and v1.3.4 into the open v1.3.2 proposal below; both are maintenance-only and wire-inert, so that proposal stands unchanged.** v1.3.3 (2026-08-07) was an empty-body patch and v1.3.4 (2026-09-18, `astarte_data_updater_plant` only — more frequent fullsweep GC on AMQPDataConsumer processes plus an RPC-availability fix; release body via `gh api repos/astarte-platform/astarte/releases` `v1.3.4`) introduce **no wire/API surface change**, so the v1.3.2 wording proposal stands complete: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"). The only delta a v1.3.4-aware doc carries is the version reference: the proposed §Infrastructure-differences sentence's "until the milestone that adopts v1.3.2 as the target" reads "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target". Raw: [v1.3.4](https://github.com/astarte-platform/astarte/releases/tag/v1.3.4).
217:- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) shows a new **stable v1.3.4** (2026-09-18) — maintenance only (data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change), so no new gap for Astrate — and still **no stable v1.4.0**: newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
229:- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-13 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
233:- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-12 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
237:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
245:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
247:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
282:  on upstream v1.4 final (currently rc.5). `APICompatVersion` stays 1.2.2 per #90's frozen
347:  fixed the stale ACL comment that cited them). But the spec ships only in `v1.4.0-rc.5` —
362:- **COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable, 2026-07-14; v1.4.0 is still rc-only).** Astrate's doc and `APICompatVersion` still target upstream **1.2.2** (`internal/realm/service.go:588`); v1.3.0 (2026-05-06) introduced wire-surface changes Astrate does not yet emulate, so this is a decision — adopt v1.3.2 as the compatibility target (then update the doc + bump `APICompatVersion` together, per the bump rule) or keep 1.2.2 and add a "not yet emulated" note. Wire-relevant v1.3.0 deltas (release notes): **MQTT v1 capabilities** incl. `purge_properties_compression_format` (plaintext vs zlib purge — touches the `emptyCache`/`producer,properties` contract COMPATIBILITY.md deviation 1 documents); **empty introspection now allowed**; **device registration triggers** (pairing) and **device deletion started/completed triggers** (RM — the latter two already exist as Astrate deviation 9 emits both around the synchronous delete); **FDO authentication** (pairing, disabled by default); **realm-scoped health** — upstream v1.3 added `GET /pairing/v1/{realm}/health`, which Astrate already serves (`internal/pairing/http.go:78-81`, comment already says "upstream 1.3+"), so deviation 18's wording ("which upstream 404s") is now false against 1.3 and the note should be reworded either way. Proposed doc wording (for your approval, edit to taste): in §Infrastructure differences add a sentence — *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl.
364:plaintext `purge_properties_compression_format`, empty-introspection allowance, device registration/deletion triggers, experimental FDO pairing auth) are not yet emulated and are out of scope until the milestone that adopts v1.3.2 as the target."* — and reword deviation 18's realm-health note from "which upstream 404s" to "added by upstream v1.3 (Astrate serves it against a 1.2.2 target; kept, matching behavior)". Raw upstream changes: [v1.3.0](https://github.com/astarte-platform/astarte/releases/tag/v1.3.0), [v1.3.2](https://github.com/astarte-platform/astarte/releases/tag/v1.3.2).

$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/tags --jq '.[0:8][] | "\(.name)\t\(.commit.sha[0:8])"' 2>&1 | head -12; echo "=== prereleases/drafts ==="; gh api 'repos/astarte-platform/astarte/releases?per_page=10' --jq '.[] | "\(.tag_name)\t\(.prerelease)\t\(.draft)\t\(.published_at)"' 2>&1 | head -12
v1.4.0-rc.6	20df0d32
v1.4.0-rc.5	3484096d
v1.4.0-rc.4	171b83c4
v1.4.0-rc.3	bf9437ca
v1.4.0-rc.2	08d7d5e6
v1.4.0-rc.1	c5fd781a
v1.4.0-rc.0	c89cda80
v1.3.4	76b21535
=== prereleases/drafts ===
v1.4.0-rc.6	true	false	2026-09-30T07:28:27Z
v1.3.4	false	false	2026-09-18T09:42:57Z
v1.4.0-rc.5	true	false	2026-08-20T12:28:04Z
v1.3.3	false	false	2026-08-07T14:33:58Z
v1.4.0-rc.4	true	false	2026-08-06T13:40:20Z
v1.4.0-rc.3	true	false	2026-07-31T10:48:39Z
v1.4.0-rc.2	true	false	2026-07-14T10:16:22Z
v1.4.0-rc.1	true	false	2026-07-14T08:53:27Z
v1.3.2	false	false	2026-07-14T10:07:26Z
v1.3.1	false	false	2026-07-13T11:50:05Z

$ cd /root/astrate-mule && tail -25 .mule/log.md
| 2026-09-27 | docs-sync-ae-add-group-device-422 | done | 274s | dff3be8 |
| 2026-09-27 | docs-sync-ae-delete-data-400 | done | 330s | 66bc366 |
| 2026-09-27 | docs-sync-ae-list-devices-422 | transient | 344s |  > build · big-pickle  Error: Error from provider (Console): Rate limit exceeded. Please try again later.  |
| 2026-09-27 | docs-sync-ae-list-devices-422 | transient | 248s |  > build · big-pickle  Error: Error from provider (Console): Rate limit exceeded. Please try again later.  |
| 2026-09-27 | docs-sync-ae-list-devices-422 | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-09-28 | docs-sync-pairing-error-example-capitalisation | done | 166s | 7ed78cb |
| 2026-09-28 | docs-sync-pairing-info-version-example | done | 333s | d9ec98f |
| 2026-09-28 | docs-sync-pairing-deviceid-base64url | done | 480s | 32693cf |
| 2026-09-28 | docs-sync-pairing-initial-payload-format-enum | done | 95s | 1f6a6c2 |
| 2026-09-28 | docs-sync-pairing-unregister-description | done | 261s | 8ffb39e |
| 2026-10-01 | fdo-rc6-scope-delta-for-giulio | done | 479s | 41af3c4 |
| 2026-10-02 | compat-note-v14-rc6 | done | 398s | 0215f97 |
| 2026-10-02 | container-stop-deadline | done | 410s | f9af73b |
| 2026-10-02 | container-timeout-bounds | blocked | 1201s | TIMEOUT after 1201s — task too big, split it |
| 2026-10-02 | container-parseconfig-rules-test | done | 274s | 32c85af |
| 2026-10-02 | container-response-cap-test | blocked | 271s | lint failed: internal/flow/blocks/container/block_test.go:500:2: redefines-builtin-id: redefinition of the built-in function cap (rev |
| 2026-10-02 | docs-sync-rm-error-example-capitalisation | done | 751s | 1d1e6f4 |
| 2026-10-02 | docs-sync-rm-legacy-alias-fields | blocked | 329s | wrote nothing |
| 2026-10-02 | docs-sync-rm-validationerror-example | done | 732s | 8a09a1a |
| 2026-10-02 | docs-sync-rm-deviceid-param | done | 303s | 2d1a1c1 |
| 2026-10-02 | docs-sync-rm-update-interface-body | blocked | 1201s | TIMEOUT after 1201s — task too big, split it |
| 2026-10-03 | docs-sync-ae-error-example-capitalisation | done | 298s | f3bcc6d |
| 2026-10-03 | docs-sync-hk-error-example-capitalisation | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-10-03 | docs-sync-native-error-example-capitalisation | done | 284s | 7cc237f |
| 2026-10-03 | docs-sync-hk-error-detail-examples-split | done | 114s | 23006d4 |

$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte --jq '"default_branch=\(.default_branch)"' && gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.6...master --jq '"ahead_by=\(.ahead_by) behind_by=\(.behind_by) status=\(.status)"' 2>&1 | head -3
default_branch=master
ahead_by=1495 behind_by=0 status=ahead

$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/commits/master --jq '"\(.sha[0:8]) \(.commit.author.date) \(.commit.message | split("\n")[0])"' && gh api repos/astarte-platform/astarte/commits/20df0d32 --jq '"tag rc.6 commit: \(.commit.author.date) \(.commit.message | split("\n")[0])"'
b6d46ad4 2026-10-02T07:49:35Z feat(appengine): include unexpected keys in object key errors (#2237)
tag rc.6 commit: 2026-09-29T10:56:18Z chore: prepare v1.4.0-rc.6 (#2227)

$ cd /root/astrate-mule && gh api 'repos/astarte-platform/astarte/commits?sha=master&since=2026-09-29T12:00:00Z&per_page=100' --jq 'length as $n | "count(first page)=\($n)", (.[] | "\(.commit.author.date[0:10]) \(.sha[0:8]) \(.commit.message | split("\n")[0])")' 2>&1 | head -50
count(first page)=9
2026-10-02 b6d46ad4 feat(appengine): include unexpected keys in object key errors (#2237)
2026-09-30 5fc16d40 chore: forward port release-1.4
2026-10-01 958bb9fe feat: implement the device_empty_cache_received trigger (#2232)
2026-10-01 b0d1ca71 build: allow credo and format from umbrella root (#2233)
2026-10-01 90efd5ee chore: remove unused labeler config (#2235)
2026-10-01 acb5b379 fleet: update libs paths (#2231)
2026-10-01 b3c13d98 ci: update paths in Release please (#2230)
2026-09-30 5ae2b936 refactor: Update libs paths in tools (#2229)
2026-09-30 34012b2c fix(ci): Update codecov libs paths (#2228)

Two new upstream commits on master since the rc.6 cut, both wire-visible candidates. Reading the actual diffs.
$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/commits/958bb9fe --jq '"FILES:", (.files[] | "\(.filename) +\(.additions)/-\(.deletions)")' 2>&1 | head -30; echo "=== 2237 ==="; gh api repos/astarte-platform/astarte/commits/b6d46ad4 --jq '"FILES:", (.files[] | "\(.filename) +\(.additions)/-\(.deletions)")' 2>&1 | head -20
FILES:
CHANGELOG.md +4/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_empty_cache_received_event.pb.ex +5/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/device_empty_cache_received_event.proto +4/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/encoder.ex +11/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/simple_event.pb.ex +5/-0
apps/astarte_core/lib/astarte_core/triggers/simple_events/simple_event.proto +2/-0
apps/astarte_core/test/astarte_core/triggers/simple_events_encoder_test.exs +12/-0
apps/astarte_core/test/astarte_core/triggers/simple_events_test.exs +10/-0
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex +2/-0
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/trigger.ex +13/-0
apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/triggers_handler.ex +19/-0
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/control_handler_test.exs +52/-0
apps/astarte_data_updater_plant/test/astarte_data_updater_plant/triggers_handler_test.exs +44/-0
apps/astarte_events/lib/astarte_events/triggers/core.ex +0/-1
doc/pages/architecture/060-triggers.md +11/-0
mix.lock +1/-1
=== 2237 ===
FILES:
apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex +9/-1
apps/astarte_appengine_api/lib/astarte_appengine_api/device/interface_value.ex +14/-1
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/fallback_controller.ex +2/-2
apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/interface_values_by_group_controller.ex +12/-37
apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex +2/-2
apps/astarte_appengine_api/test/astarte_appengine_api/device/device_test.exs +1/-1
apps/astarte_appengine_api/test/astarte_appengine_api/device/interface_value_test.exs +2/-1
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_values_by_device_alias_controller_test.exs +4/-1
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_values_by_group_controller_test.exs +4/-1
apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_values_controller_test.exs +4/-1

$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/commits/958bb9fe --jq '.files[] | select(.filename | test("encoder.ex|device_empty_cache_received_event.proto|simple_event.proto|core/trigger.ex|triggers_handler.ex|060-triggers.md")) | "### \(.filename)\n\(.patch)"' 2>&1 | head -160
### apps/astarte_core/lib/astarte_core/triggers/simple_events/device_empty_cache_received_event.proto
@@ -0,0 +1,4 @@
+syntax = "proto3";
+
+message DeviceEmptyCacheReceivedEvent {
+}
### apps/astarte_core/lib/astarte_core/triggers/simple_events/encoder.ex
@@ -51,6 +51,17 @@ defmodule Astarte.Core.Triggers.SimpleEvents.Encoder do
     end
   end
 
+  defimpl Jason.Encoder, for: SimpleEvents.DeviceEmptyCacheReceivedEvent do
+    alias Astarte.Core.Triggers.SimpleEvents.DeviceEmptyCacheReceivedEvent
+
+    def encode(%DeviceEmptyCacheReceivedEvent{}, opts) do
+      %{
+        "type" => "device_empty_cache_received"
+      }
+      |> Jason.Encoder.encode(opts)
+    end
+  end
+
   defimpl Jason.Encoder, for: SimpleEvents.DeviceErrorEvent do
     alias Astarte.Core.Triggers.SimpleEvents.DeviceErrorEvent
 
### apps/astarte_core/lib/astarte_core/triggers/simple_events/simple_event.proto
@@ -20,6 +20,7 @@ syntax = "proto3";
 
 import "lib/astarte_core/triggers/simple_events/device_connected_event.proto";
 import "lib/astarte_core/triggers/simple_events/device_disconnected_event.proto";
+import "lib/astarte_core/triggers/simple_events/device_empty_cache_received_event.proto";
 import "lib/astarte_core/triggers/simple_events/device_registered_event.proto";
 import "lib/astarte_core/triggers/simple_events/device_deletion_started_event.proto";
 import "lib/astarte_core/triggers/simple_events/device_deletion_finished_event.proto";
@@ -62,5 +63,6 @@ message SimpleEvent {
         DeviceRegisteredEvent device_registered_event = 20;
         DeviceDeletionStartedEvent device_deletion_started_event = 21;
         DeviceDeletionFinishedEvent device_deletion_finished_event = 22;
+        DeviceEmptyCacheReceivedEvent device_empty_cache_received_event = 23;
     }
 }
### apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/trigger.ex
@@ -62,6 +62,19 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.Trigger do
     end
   end
 
+  def execute_empty_cache_received_triggers(state, timestamp) do
+    timestamp_ms = div(timestamp, 10_000)
+
+    TriggersHandler.device_empty_cache_received(
+      state.realm,
+      state.device_id,
+      state.groups,
+      timestamp_ms
+    )
+
+    :ok
+  end
+
   def execute_device_error_triggers(state, error_name, error_metadata \\ %{}, timestamp) do
     timestamp_ms = div(timestamp, 10_000)
 
### apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/triggers_handler.ex
@@ -25,6 +25,7 @@ defmodule Astarte.DataUpdaterPlant.TriggersHandler do
   alias Astarte.Core.Triggers.SimpleEvents.{
     DeviceConnectedEvent,
     DeviceDisconnectedEvent,
+    DeviceEmptyCacheReceivedEvent,
     DeviceErrorEvent,
     IncomingDataEvent,
     IncomingIntrospectionEvent,
@@ -83,6 +84,24 @@ defmodule Astarte.DataUpdaterPlant.TriggersHandler do
     end)
   end
 
+  def device_empty_cache_received(realm, device_id, groups, timestamp) do
+    event = %DeviceEmptyCacheReceivedEvent{}
+    hw_id = Device.encode_device_id(device_id)
+
+    Triggers.find_device_trigger_targets(realm, device_id, groups, :on_empty_cache_received)
+    |> execute_all_ok(fn {target, policy} ->
+      dispatch_event_with_telemetry(
+        event,
+        :device_empty_cache_received_event,
+        target,
+        realm,
+        hw_id,
+        timestamp,
+        policy
+      )
+    end)
+  end
+
   def device_error(
         realm,
         device_id,
### doc/pages/architecture/060-triggers.md
@@ -104,6 +104,9 @@ This is the generic representation of a Device Trigger:
 
 - `device_connected`: triggered when a device connects to its transport.
 - `device_disconnected`: triggered when a device disconnects from its transport.
+- `device_empty_cache_received`: triggered when a device sends an `/emptyCache` message and Astarte
+  has finished processing it, i.e. after the device properties have been resent and the pending
+  empty cache flag has been cleared. It is not triggered if that processing fails.
 - `device_error`: triggered when data from a device causes an error.
 - `device_registered`: triggered when a device is registered to Astarte.
 - `device_deletion_started`: triggered when the deletion of a device is started.
@@ -285,6 +288,14 @@ Additionally, the realm that originated the trigger is available in the request
 }
 ```
 
+###### DeviceEmptyCacheReceivedEvent
+
+```json
+{
+  "type": "device_empty_cache_received"
+}
+```
+
 ###### DeviceErrorEvent
 
 ```json

$ cd /root/astrate-mule && echo "=== empty_cache in astrate ===" && rg -n "empty_cache_received|emptyCache|EmptyCache" --glob '!**/*_test.go' -l | head -20 && echo "=== on_device_connected family ===" && rg -n "on_device_connected|on_device_disconnected|on_device_registered|on_device_deletion" --glob '!**/*_test.go' | head -20
=== empty_cache in astrate ===
docs/api/astarte_realm_management_api.yaml
docs/ROADMAP.md
docs/COMPATIBILITY.md
docs/site/pairing-and-security.md
docs/site/payload-formats.md
docs/site/json-payload-profile.md
docs/site/mqtt-protocol.md
docs/JSON-PAYLOAD-PROFILE.md
internal/engine/data.go
internal/engine/control.go
internal/engine/triggers/match.go
internal/engine/cache.go
docs/DESIGN.md
bench/device.go
internal/engine/topics.go
internal/engine/router.go
internal/broker/aclhook.go
internal/testutil/astartedevice.go
test/conformance/upstream/recordchannels/main.go
test/conformance/upstream/README.md
=== on_device_connected family ===

$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/commits/b6d46ad4 --jq '.files[] | select(.filename | test("device.ex|interface_value.ex|fallback_controller.ex|error_view.ex|interface_values_by_group_controller.ex$")) | "### \(.filename)\n\(.patch)"' 2>&1 | head -140
### apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex
@@ -565,6 +565,13 @@ defmodule Astarte.AppEngine.API.Device do
         Logger.warning("Unexpected value type.", tag: "unexpected_value_type")
         {:error, :unexpected_value_type, expected: value_type}
 
+      {:error, :unexpected_object_key, keys: keys} ->
+        Logger.warning("Unexpected object keys #{inspect(keys)} in object interface update.",
+          tag: "unexpected_object_key"
+        )
+
+        {:error, :unexpected_object_key, keys: keys}
+
       {:error, :invalid_object_aggregation_path} ->
         Logger.warning("Error while trying to publish on path for object aggregated interface.",
           tag: "invalid_object_aggregation_path"
@@ -875,7 +882,8 @@ defmodule Astarte.AppEngine.API.Device do
           {:halt, {:error, reason}}
 
         :error ->
-          {:halt, {:error, :unexpected_object_key}}
+          keys = InterfaceValue.unexpected_keys(mappings_by_key, object)
+          {:halt, {:error, :unexpected_object_key, keys: keys}}
       end
     end)
   end
### apps/astarte_appengine_api/lib/astarte_appengine_api/device/interface_value.ex
@@ -33,7 +33,7 @@ defmodule Astarte.AppEngine.API.Device.InterfaceValue do
           {:halt, {:error, reason, expected}}
 
         :error ->
-          {:halt, {:error, :unexpected_object_key}}
+          {:halt, {:error, :unexpected_object_key, keys: unexpected_keys(expected_types, object)}}
       end
     end)
   end
@@ -148,6 +148,19 @@ defmodule Astarte.AppEngine.API.Device.InterfaceValue do
     {:ok, anyvalue}
   end
 
+  @doc """
+  Return the keys of `object` that have no counterpart in `expected`, sorted.
+
+  Only the keys of `expected` are looked at, so it accepts both an expected
+  types map and a mappings map.
+  """
+  def unexpected_keys(expected, object) do
+    object
+    |> Map.keys()
+    |> Enum.reject(&Map.has_key?(expected, &1))
+    |> Enum.sort()
+  end
+
   defp map_while_ok(values, fun) when is_list(values) do
     result =
       Enum.reduce_while(values, {:ok, []}, fn value, {:ok, acc} ->
### apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/fallback_controller.ex
@@ -198,11 +198,11 @@ defmodule Astarte.AppEngine.APIWeb.FallbackController do
     |> render(:"422_invalid_attributes")
   end
 
-  def call(conn, {:error, :unexpected_object_key}) do
+  def call(conn, {:error, :unexpected_object_key, keys: keys}) do
     conn
     |> put_status(:bad_request)
     |> put_view(Astarte.AppEngine.APIWeb.ErrorView)
-    |> render(:"422_unexpected_object_key")
+    |> render(:"422_unexpected_object_key", keys: keys)
   end
 
   def call(conn, {:error, :unset_not_allowed}) do
### apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/interface_values_by_group_controller.ex
@@ -205,18 +205,11 @@ defmodule Astarte.AppEngine.APIWeb.InterfaceValuesByGroupController do
         "group_name" => group_name,
         "device_id" => device_id
       }) do
-    with {:ok, true} <- Groups.check_device_in_group(realm_name, group_name, device_id),
+    with :ok <- ensure_device_in_group(realm_name, group_name, device_id),
          {:ok, interfaces} <- Device.list_interfaces(realm_name, device_id) do
       conn
       |> put_view(InterfaceValuesView)
       |> render("index.json", interfaces: interfaces)
-    else
-      {:ok, false} ->
-        {:error, :device_not_found}
-
-      {:error, reason} ->
-        # To FallbackController
-        {:error, reason}
     end
   end
 
@@ -238,7 +231,7 @@ defmodule Astarte.AppEngine.APIWeb.InterfaceValuesByGroupController do
            "path" => path
          } = parameters
        ) do
-    with {:ok, true} <- Groups.check_device_in_group(realm_name, group_name, device_id),
+    with :ok <- ensure_device_in_group(realm_name, group_name, device_id),
          {:ok, %InterfaceValues{} = interface_values} <-
            Device.get_interface_values!(
              realm_name,
@@ -250,13 +243,6 @@ defmodule Astarte.AppEngine.APIWeb.InterfaceValuesByGroupController do
       conn
       |> put_view(InterfaceValuesView)
       |> render("show.json", interface_values: interface_values)
-    else
-      {:ok, false} ->
-        {:error, :device_not_found}
-
-      {:error, reason} ->
-        # To FallbackController
-        {:error, reason}
     end
   end
 
@@ -269,7 +255,7 @@ defmodule Astarte.AppEngine.APIWeb.InterfaceValuesByGroupController do
            "interface" => interface
          } = parameters
        ) do
-    with {:ok, true} <- Groups.check_device_in_group(realm_name, group_name, device_id),
+    with :ok <- ensure_device_in_group(realm_name, group_name, device_id),
          {:ok, %InterfaceValues{} = interface_values} <-
            Device.get_interface_values!(
              realm_name,
@@ -280,13 +266,6 @@ defmodule Astarte.AppEngine.APIWeb.InterfaceValuesByGroupController do
       conn
       |> put_view(InterfaceValuesView)
       |> render("show.json", interface_values: interface_values)
-    else
-      {:ok, false} ->
-        {:error, :device_not_found}
-
-      {:error, reason} ->
-        # To FallbackController
-        {:error, reason}
     end
   end
 
@@ -301,7 +280,7 @@ defmodule Astarte.AppEngine.APIWeb.InterfaceValuesByGroupController do
           "data" => value
         } = parameters

$ cd /root/astrate-mule && echo "=== device trigger names ===" && rg -n "device_connected|device_disconnected|device_registered|device_deletion_started" internal/ pkg/ --glob '!**/*_test.go' | head -25 && echo "=== unexpected_object_key ===" && rg -n "unexpected_object_key|unexpected keys|UnexpectedObjectKey" internal/ pkg/ | head -20
=== device trigger names ===
internal/engine/stream/bus.go:22:	KindDeviceRegistered = "device_registered"
internal/engine/stream/bus.go:24:	KindDeviceConnected = "device_connected"
internal/engine/stream/bus.go:26:	KindDeviceDisconnected = "device_disconnected"
internal/engine/stream/bus.go:29:	KindDeviceDeletionStarted = "device_deletion_started"
internal/engine/stream/bus.go:53:	// IP is the peer address of a device_connected event, empty otherwise.
internal/engine/engine.go:328:// eviction): device_connected / device_disconnected triggers, live bus
internal/engine/engine.go:375:// device registration. It fires device_registered triggers and publishes the
internal/engine/engine.go:376:// event on the live bus — the same path device_connected already uses.
internal/engine/engine.go:395:// before a synchronous device delete. It fires device_deletion_started
internal/engine/triggers/testdata/device_connected.json:1:{"timestamp":"2026-06-12T10:00:00.123Z","device_id":"f0VMRgIBAQAAAAAAAAAAAA","event":{"type":"device_connected","device_ip_address":"203.0.113.89"},"trigger_name":"example_trigger"}
internal/engine/triggers/testdata/device_registered.json:1:{"timestamp":"2026-06-12T10:00:00.123Z","device_id":"f0VMRgIBAQAAAAAAAAAAAA","event":{"type":"device_registered"},"trigger_name":"example_trigger"}
internal/engine/triggers/testdata/device_deletion_started.json:1:{"timestamp":"2026-06-12T10:00:00.123Z","device_id":"f0VMRgIBAQAAAAAAAAAAAA","event":{"type":"device_deletion_started"},"trigger_name":"example_trigger"}
internal/engine/triggers/testdata/device_disconnected.json:1:{"timestamp":"2026-06-12T10:00:00.123Z","device_id":"f0VMRgIBAQAAAAAAAAAAAA","event":{"type":"device_disconnected"},"trigger_name":"example_trigger"}
internal/engine/triggers/match.go:50:	OnDeviceRegistered = "device_registered"
internal/engine/triggers/match.go:52:	OnDeviceConnected = "device_connected"
internal/engine/triggers/match.go:54:	OnDeviceDisconnected = "device_disconnected"
internal/engine/triggers/match.go:57:	OnDeviceDeletionStarted = "device_deletion_started"
internal/engine/triggers/events.go:181:// DeviceRegisteredEvent is the device_registered event body.
internal/engine/triggers/events.go:183:	// Type is always "device_registered".
internal/engine/triggers/events.go:187:// NewDeviceRegisteredEvent builds a device_registered event body.
internal/engine/triggers/events.go:192:// DeviceConnectedEvent is the device_connected event body.
internal/engine/triggers/events.go:194:	// Type is always "device_connected".
internal/engine/triggers/events.go:200:// NewDeviceConnectedEvent builds a device_connected event body.
internal/engine/triggers/events.go:205:// DeviceDisconnectedEvent is the device_disconnected event body.
internal/engine/triggers/events.go:207:	// Type is always "device_disconnected".
=== unexpected_object_key ===
internal/engine/triggers/errorname.go:16:	"unexpected_object_key",
internal/engine/triggers/errorname.go:61:	"bad_object":                     "unexpected_object_key",
internal/engine/triggers/errorname.go:64:	"missing_required":  "unexpected_object_key",
internal/engine/triggers/errorname_test.go:29:		{"bad_object", "unexpected_object_key"},
internal/engine/triggers/errorname_test.go:30:		{"missing_required", "unexpected_object_key"},

$ cd /root/astrate-mule && echo "=== match.go 40-80 ===" && sed -n '40,80p' internal/engine/triggers/match.go && echo "=== where emptyCache handled ===" && rg -n "EmptyCache|emptyCache|empty_cache" internal/engine/*.go | head -25
=== match.go 40-80 ===
	OnValueChangeApplied = "value_change_applied"
	// OnPathCreated fires on the first value of a path (previous missing,
	// new value present — upstream Core.Trigger post-change semantics).
	OnPathCreated = "path_created"
	// OnPathRemoved fires when an existing property is unset (previous
	// present, new value absent).
	OnPathRemoved = "path_removed"
	// OnValueStored fires after an accepted individual-datastream insert.
	OnValueStored = "value_stored"
	// OnDeviceRegistered fires when a device is registered via the Pairing API.
	OnDeviceRegistered = "device_registered"
	// OnDeviceConnected fires when a device session is established.
	OnDeviceConnected = "device_connected"
	// OnDeviceDisconnected fires when a device connection ends.
	OnDeviceDisconnected = "device_disconnected"
	// OnDeviceDeletionStarted fires immediately before a device is deleted
	// (upstream: start of async deletion; Astrate: before the sync delete).
	OnDeviceDeletionStarted = "device_deletion_started"
	// OnDeviceDeletionFinished fires immediately after a device is deleted
	// (upstream: end of async deletion; Astrate: after the sync delete).
	OnDeviceDeletionFinished = "device_deletion_finished"
	// OnDeviceEmptyCacheReceived fires on control/emptyCache (accepted, not
	// evaluated in v1: upstream defines no SimpleEvent variant for it).
	OnDeviceEmptyCacheReceived = "device_empty_cache_received"
	// OnDeviceError fires when the validation pipeline rejects a message.
	OnDeviceError = "device_error"
	// OnIncomingIntrospection fires on every introspection publish.
	OnIncomingIntrospection = "incoming_introspection"
	// OnInterfaceAdded fires when an introspection declares a new
	// name:major pair.
	OnInterfaceAdded = "interface_added"
	// OnInterfaceRemoved fires when an introspection drops a name:major pair.
	OnInterfaceRemoved = "interface_removed"
	// OnInterfaceMinorUpdated fires on a minor bump (accepted, not evaluated
	// in v1).
	OnInterfaceMinorUpdated = "interface_minor_updated"
)

// anyToken is the wildcard accepted for device IDs, interface names, and
// value operators.
const anyToken = "*"
=== where emptyCache handled ===
internal/engine/topics_test.go:25:		{"control", "alpha/" + dev + "/control/emptyCache", "control/emptyCache", true},
internal/engine/topics_test.go:54:		{"control/emptyCache", kindControl, "emptyCache"},
internal/engine/control_test.go:131:// TestEmptyCache: the device receives every server-owned property on its
internal/engine/control_test.go:134:func TestEmptyCache(t *testing.T) {
internal/engine/control_test.go:153:	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
internal/engine/control_test.go:155:		t.Fatal("emptyCache not acknowledged")
internal/engine/control_test.go:212:		t.Error("emptyCache did not arm the hint reset")
internal/engine/control_test.go:216:// TestEmptyCacheJSONHint: a device flipped to the JSON profile receives its
internal/engine/control_test.go:217:// emptyCache resends as JSON documents; the control frame stays zlib
internal/engine/control_test.go:219:func TestEmptyCacheJSONHint(t *testing.T) {
internal/engine/control_test.go:232:	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
internal/engine/control_test.go:234:		t.Fatal("emptyCache not acknowledged")
internal/engine/engine_e2e_test.go:258:	t.Run("EmptyCacheResync", func(t *testing.T) { testE2EEmptyCache(t, env) })
internal/engine/engine_e2e_test.go:349:// testE2EEmptyCache: after a server-owned property is set, control/emptyCache
internal/engine/engine_e2e_test.go:353:func testE2EEmptyCache(t *testing.T, env *e2eEnv) {
internal/engine/engine_e2e_test.go:364:	dev.EmptyCache(t)
internal/engine/data.go:506:// after an emptyCache armed the reset, the next BSON data payload flips it
internal/engine/capabilities_test.go:173:	// emptyCache triggers consumer/properties send.
internal/engine/capabilities_test.go:175:	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
internal/engine/capabilities_test.go:177:		t.Fatal("emptyCache not acknowledged")
internal/engine/capabilities_test.go:208:	// emptyCache without any capability set → zlib (default).
internal/engine/capabilities_test.go:210:	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
internal/engine/capabilities_test.go:212:		t.Fatal("emptyCache not acknowledged")
internal/engine/data_test.go:328:// emptyCache reset followed by a BSON payload.
internal/engine/data_test.go:362:		t.Error("hint flipped back to bson without an emptyCache reset")

$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/commits/958bb9fe --jq '.files[] | select(.filename | test("control_handler.ex|astarte_events/lib")) | "### \(.filename)\n\(.patch)"' 2>&1 | head -40; echo "=== 422_unexpected_object_key view ==="; gh api repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/422_unexpected_object_key.json.eex?ref=b6d46ad4 --jq '.content' 2>&1 | head -5
### apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/control_handler.ex
@@ -158,6 +158,8 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.ControlHandler do
     with :ok <- send_control_consumer_properties(state, timestamp),
          {:ok, state} <- resend_all_properties(state, timestamp),
          :ok <- set_pending_empty_cache(state, timestamp) do
+      :ok = Core.Trigger.execute_empty_cache_received_triggers(state, timestamp)
+
       :telemetry.execute(
         [:astarte, :data_updater_plant, :data_updater, :processed_empty_cache],
         %{},
### apps/astarte_events/lib/astarte_events/triggers/core.ex
@@ -187,7 +187,6 @@ defmodule Astarte.Events.Triggers.Core do
     end
   end
 
-  # TODO: implement on_empty_cache_received
   def load_trigger(realm_name, {:device_trigger, proto_buf_device_trigger}, trigger_target, state) do
     device_triggers = state.device_triggers
 
=== 422_unexpected_object_key view ===
{"message":"Not Found","documentation_url":"https://docs.github.com/rest/repos/contents#get-repository-content","status":"404"}gh: Not Found (HTTP 404)

$ cd /root/astrate-mule && echo "=== astrate object key rejection ===" && rg -n "ReasonBadObject|bad_object|missing_required|unknown key|unexpected key" internal/ pkg/ --glob '!**/*_test.go' | head -20 && echo "=== 400 body for object write ===" && rg -n "BadObject|bad_object" internal/appengine/*.go internal/engine/*.go | head -15
=== astrate object key rejection ===
pkg/payload/value.go:143:	// ReasonBadObject: object-aggregation shape violation (not a
pkg/payload/value.go:145:	ReasonBadObject
pkg/payload/value.go:157:	"type_mismatch", "value_too_large", "bad_object", "unset_not_allowed",
pkg/payload/value.go:158:	"missing_required",
pkg/payload/json.go:203:		return nil, rejectf(ReasonBadObject, "JSON value %s where an object-aggregation document was expected", clip(raw))
pkg/payload/json.go:210:		return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
pkg/payload/json.go:216:			return nil, rejectf(ReasonBadObject, "key %q matches no declared object leaf", key)
pkg/payload/json.go:365:			return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
pkg/payload/json.go:370:				return nil, rejectf(ReasonBadObject, "key %q: nested objects are not encodable", k)
internal/engine/triggers/errorname.go:61:	"bad_object":                     "unexpected_object_key",
internal/engine/triggers/errorname.go:62:	// Upstream 1.4 emits missing_required_mapping, outside dashboard
internal/engine/triggers/errorname.go:64:	"missing_required":  "unexpected_object_key",
pkg/payload/bson.go:140:		return nil, rejectf(ReasonBadObject, "BSON %s where an object-aggregation document was expected", rv.Type)
pkg/payload/bson.go:147:		return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
pkg/payload/bson.go:157:			return nil, rejectf(ReasonBadObject, "key %q matches no declared object leaf", key)
pkg/payload/bson.go:160:			return nil, rejectf(ReasonBadObject, "duplicate object key %q", key)
pkg/payload/bson.go:319:		return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
pkg/payload/bson.go:331:			return nil, rejectf(ReasonBadObject, "invalid object key %q", k)
pkg/payload/bson.go:334:			return nil, rejectf(ReasonBadObject, "key %q: nested objects are not encodable", k)
=== 400 body for object write ===
internal/engine/data_test.go:186:			reason: payload.ReasonBadObject.String(),

$ cd /root/astrate-mule && gh api 'repos/astarte-platform/astarte/contents/apps/astarte_appengine_api/lib/astarte_appengine_api_web/views?ref=b6d46ad4' --jq '.[].name' 2>&1 | head -30
changeset_view.ex
device_status_by_group_view.ex
device_status_view.ex
error_helpers.ex
error_view.ex
groups_view.ex
interface_values_view.ex
stats_view.ex
version_view.ex
watch_request_view.ex

$ cd /root/astrate-mule && rg -n "ReasonBadObject|reasonFor|ValidateValue|WriteError" internal/appengine/data.go | head -20; echo "=== write validation entry ==="; rg -n "func .*Write|validate" internal/engine/serverdata.go | head -20
=== write validation entry ===
38:// PublishServerValue validates, persists, and delivers one server-owned
41:// validated against the mapping exactly like an inbound JSON-profile
66:// PublishDeviceValue validates and persists one device-owned value without
88:// publishAsOwner is the shared validate-and-persist core of the publish

$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/commits/b6d46ad4 --jq '.files[] | select(.filename | test("error_view.ex")) | .patch' 2>&1 | head -30; echo "=== ensure_device_in_group ==="; gh api repos/astarte-platform/astarte/commits/b6d46ad4 --jq '.files[] | select(.filename | test("interface_values_by_group_controller.ex$")) | .patch' 2>&1 | tail -30
@@ -98,8 +98,8 @@ defmodule Astarte.AppEngine.APIWeb.ErrorView do
     %{errors: %{detail: "Invalid attributes"}}
   end
 
-  def render("422_unexpected_object_key.json", _assigns) do
-    %{errors: %{detail: "Unexpected object key"}}
+  def render("422_unexpected_object_key.json", %{keys: keys} = _assigns) do
+    %{errors: %{detail: "Unexpected object key", unexpected_keys: keys}}
   end
 
   def render("422_missing_required_mapping.json", _assigns) do
=== ensure_device_in_group ===
-
-      {:error, reason} ->
-        {:error, reason}
     end
   end
 
@@ -330,15 +303,17 @@ defmodule Astarte.AppEngine.APIWeb.InterfaceValuesByGroupController do
         "interface" => interface,
         "path" => path
       }) do
-    with {:ok, true} <- Groups.check_device_in_group(realm_name, group_name, device_id),
+    with :ok <- ensure_device_in_group(realm_name, group_name, device_id),
          :ok <- Device.delete_interface_values(realm_name, device_id, interface, path) do
       send_resp(conn, :no_content, "")
-    else
-      {:ok, false} ->
-        {:error, :device_not_found}
+    end
+  end
 
-      {:error, reason} ->
-        {:error, reason}
+  defp ensure_device_in_group(realm_name, group_name, device_id) do
+    case Groups.check_device_in_group(realm_name, group_name, device_id) do
+      {:ok, true} -> :ok
+      {:ok, false} -> {:error, :device_not_found}
+      {:error, reason} -> {:error, reason}
     end
   end
 end

$ cd /root/astrate-mule && rg -n "Validate|Reason|StatusBadRequest|WriteBadRequest" internal/appengine/data.go | head -25; echo "=== serverdata validate ==="; sed -n '88,140p' internal/engine/serverdata.go
=== serverdata validate ===
// publishAsOwner is the shared validate-and-persist core of the publish
// paths (docs/ROADMAP.md §7.2 file 6.9): realm, device, and interface
// resolution with the ownership gate parameterised, §2.6 payload
// validation, and the persist switch. ownership selects which side may
// write and which sentinel an ownership mismatch carries. On success the
// op persists immediately and then runs the same post-commit observers as
// the broker ingress path — previous-value capture for the change-derived
// triggers, then fireData (data triggers + live bus), so a server- or
// device-owned write is observed exactly like a real device's data
// (devicedata_test.go:3-4). Wire delivery stays with the caller.
func (e *Engine) publishAsOwner(ctx context.Context, realm string, id deviceid.ID,
	ifaceName, path string, value json.RawMessage, ts *time.Time,
	ownership interfaceschema.Ownership) (*ownerDelivery, error) {
	notOwned := ErrNotDeviceOwned
	if ownership == interfaceschema.OwnershipServer {
		notOwned = ErrNotServerOwned
	}
	rs := e.schemas.realmOrReload(ctx, realm)
	if rs == nil {
		return nil, fmt.Errorf("%w: %s", ErrRealmUnknown, realm)
	}
	view, err := e.deviceView(ctx, rs, id)
	if err != nil {
		return nil, err
	}
	ci := resolveServerInterface(rs, view, ifaceName)
	if ci == nil {
		return nil, fmt.Errorf("%w: %s", ErrInterfaceNotFound, ifaceName)
	}
	if ci.Ownership != ownership {
		return nil, fmt.Errorf("%w: %s", notOwned, ifaceName)
	}

	effTS := time.Now().UTC()
	if ts != nil {
		effTS = ts.UTC()
	}
	envelope, err := serverEnvelope(value, effTS)
	if err != nil {
		return nil, err
	}

	var (
		dp      payload.DecodedPayload
		mapping *interfaceschema.CompiledMapping
	)
	dec := payload.Decoder{MaxSize: len(envelope)}
	if ci.Aggregation == interfaceschema.AggregationObject {
		if !objectPathOK(ci, path) {
			return nil, fmt.Errorf("%w: %q is not an aggregation prefix of %s", ErrPathNotFound, path, ifaceName)
		}
		dp, err = dec.Object(envelope, ci.ObjectLeaves)
		mapping = anyObjectLeaf(ci)

$ cd /root/astrate-mule && rg -n "payload.Reason|Reason\)" internal/engine/*.go internal/appengine/*.go --glob '!**/*_test.go' | grep -v "^internal/engine/data.go" | head -25
internal/engine/data_test.go:134:			reason: payload.ReasonTooLarge.String(),
internal/engine/data_test.go:141:			reason: payload.ReasonUnknownFormat.String(),
internal/engine/data_test.go:148:			reason: payload.ReasonMalformed.String(),
internal/engine/data_test.go:155:			reason: payload.ReasonNoValue.String(),
internal/engine/data_test.go:163:			reason: payload.ReasonBadTimestamp.String(),
internal/engine/data_test.go:170:			reason: payload.ReasonTypeMismatch.String(),
internal/engine/data_test.go:178:			reason: payload.ReasonValueTooLarge.String(),
internal/engine/data_test.go:186:			reason: payload.ReasonBadObject.String(),
internal/engine/data_test.go:193:			reason: payload.ReasonUnsetNotAllowed.String(),
internal/engine/data_test.go:200:			reason: payload.ReasonUnsetNotAllowed.String(),
internal/engine/serverdata_test.go:196:		if payload.ReasonOf(err) == payload.ReasonNone {
internal/engine/devicedata_test.go:164:		if got := payload.ReasonOf(err); got != payload.ReasonTypeMismatch {

$ cd /root/astrate-mule && rg -n "ReasonOf" --glob '!**/*_test.go' | head -15; echo "=== appengine writeError sentinels ==="; rg -n "ErrBadRequest|WriteBadRequest|StatusBadRequest|Invalid value|ErrValue" internal/appengine/http.go internal/engine/errors.go 2>/dev/null | head -20
pkg/payload/value.go:197:// ReasonOf extracts the RejectReason from err, or ReasonNone if err is nil
pkg/payload/value.go:199:func ReasonOf(err error) RejectReason {
internal/engine/data.go:289:		reason := payload.ReasonOf(err)
=== appengine writeError sentinels ===
internal/appengine/http.go:215:		_ = astarteapi.WriteBadRequest(w)
internal/appengine/http.go:231:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Bad request")
internal/appengine/http.go:236:		_ = astarteapi.WriteBadRequest(w)
internal/appengine/http.go:295:		_ = astarteapi.WriteBadRequest(w)
internal/appengine/http.go:324:		_ = astarteapi.WriteBadRequest(w)
internal/appengine/http.go:343:		_ = astarteapi.WriteBadRequest(w)
internal/appengine/http.go:401:		_ = astarteapi.WriteBadRequest(w)
internal/appengine/http.go:489:		_ = astarteapi.WriteBadRequest(w)
internal/appengine/http.go:527:		_ = astarteapi.WriteBadRequest(w)
internal/appengine/http.go:649:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Invalid alias")
internal/appengine/http.go:653:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Alias tag not found")
internal/appengine/http.go:655:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Invalid attributes")
internal/appengine/http.go:680:		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Endpoint not found")

$ cd /root/astrate-mule && sed -n '640,690p' internal/appengine/http.go; echo "=== putData ==="; sed -n '286,312p' internal/appengine/http.go
}

// writeError maps service/store errors onto upstream-shaped responses.
func (a *API) writeError(w http.ResponseWriter, err error) {
	var fe FieldErrors
	switch {
	case errors.As(err, &fe):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity, fe)
	case errors.Is(err, ErrInvalidAlias):
		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Invalid alias")
	case errors.Is(err, ErrAliasAlreadyInUse):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Alias already in use")
	case errors.Is(err, ErrAliasTagNotFound):
		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Alias tag not found")
	case errors.Is(err, ErrInvalidAttributes):
		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Invalid attributes")
	case errors.Is(err, ErrAttributeKeyNotFound):
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Attribute key not found")
	case errors.Is(err, ErrValidation):
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
	case errors.Is(err, ErrGroupAlreadyExists):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Group already exists")
	case errors.Is(err, ErrDeviceAlreadyInGroup):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Device already in group")
	case errors.Is(err, store.ErrAlreadyExists):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Already exists")
	case errors.Is(err, ErrGroupNotFound):
		_ = astarteapi.WriteError(w, http.StatusNotFound, "Group not found")
	case errors.Is(err, ErrPathNotFound):
		_ = astarteapi.WriteError(w, http.StatusNotFound, "Path not found")
	// Engine write sentinels must precede store.ErrNotFound: wrapped
	// ErrInterfaceNotFound errors are distinct from store's, and Go's switch
	// takes the first match (issue #57, measured taxonomy).
	case errors.Is(err, engine.ErrNotServerOwned):
		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Cannot write to device owned resource")
	case errors.Is(err, engine.ErrNotAProperty):
		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Cannot write to read-only resource")
	case errors.Is(err, engine.ErrInterfaceNotFound):
		_ = astarteapi.WriteError(w, http.StatusNotFound, "Interface not found in device introspection")
	case errors.Is(err, engine.ErrPathNotFound):
		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Endpoint not found")
	case errors.Is(err, store.ErrNotFound):
		_ = astarteapi.WriteDeviceNotFound(w)
	default:
		_ = astarteapi.WriteInternalServerError(w)
	}
}

// validationDetail strips the ErrValidation prefix for the response detail.
func validationDetail(err error) string {
	const prefix = "appengine: validation failed: "
=== putData ===
	a.serveData(w, r, func(opts QueryOpts) (any, error) {
		return a.svc.GetDataInGroup(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device"),
			r.PathValue("interface"), pathParam(r), opts)
	})
}

func (a *API) putData(w http.ResponseWriter, r *http.Request) {
	var value json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &value); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if len(value) > maxValueBytes {
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
		return
	}
	err := a.svc.PublishData(r.Context(), r.PathValue("realm"), r.PathValue("device"),
		r.PathValue("interface"), pathParam(r), value, nil)
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) deleteData(w http.ResponseWriter, r *http.Request) {
	err := a.svc.UnsetProperty(r.Context(), r.PathValue("realm"), r.PathValue("device"),

$ cd /root/astrate-mule && rg -n -i "unexpected object key|unexpected_object_key|object key" .mule/todo.md .mule/tasks/*.md 2>/dev/null | head -10; echo "=== empty_cache trigger in queue ==="; rg -n "empty_cache_received" .mule/todo.md docs/COMPATIBILITY.md | head -10
.mule/todo.md:1:- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
=== empty_cache trigger in queue ===

Now verifying the Astrate-side anchors my task lines will name.
$ cd /root/astrate-mule && rg -n "func .*[Ee]mptyCache|emptyCache" internal/engine/control.go internal/engine/cache.go | head -20; echo "=== match.go device trigger set / where events fire ==="; rg -n "OnDeviceConnected|OnDeviceDeletionFinished" internal/engine/*.go --glob '!**/*_test.go' | head -12
internal/engine/control.go:21:	controlEmptyCache         = "emptyCache"
internal/engine/control.go:52:// handleEmptyCache implements `control/emptyCache` (docs/DESIGN.md §3.3):
internal/engine/control.go:58:func (e *Engine) handleEmptyCache(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
internal/engine/control.go:66:	if !e.retryStore(ctx, m, "emptyCache resync", func() error {
internal/engine/control.go:140:// device on its data topic (QoS 2, the emptyCache resync of docs/DESIGN.md
internal/engine/control.go:185:// property paths (docs/DESIGN.md §3.4). It is sent after emptyCache, after
internal/engine/cache.go:66:	// (emptyCache resend + `consumer/properties` payload, docs/DESIGN.md §3.4).
internal/engine/cache.go:360:	// resetHintOnBSON is armed by the emptyCache control handler (M6b): the
internal/engine/cache.go:395:// armHintReset arms the §3.5.4 emptyCache rule: the next BSON data payload
internal/engine/cache.go:397:// calls it while processing `control/emptyCache`.
=== match.go device trigger set / where events fire ===
internal/engine/engine.go:348:			triggers.DeviceEvent{DeviceID: deviceID, On: triggers.OnDeviceConnected},
internal/engine/engine.go:409:		triggers.OnDeviceDeletionFinished, stream.KindDeviceDeletionFinished,

$ cd /root/astrate-mule && sed -n '52,100p' internal/engine/control.go; echo "=== engine.go 328-360 (device_connected pattern) ==="; sed -n '326,356p' internal/engine/engine.go
// handleEmptyCache implements `control/emptyCache` (docs/DESIGN.md §3.3):
// the device lost its local cache, so Astrate re-sends every server-owned
// property on its data topic (QoS 2, format per the device hint) and then
// publishes the `consumer/properties` purge message. It also arms the
// §3.5.4 hint reset: if the device's next data payload is BSON, the sticky
// json hint flips back.
func (e *Engine) handleEmptyCache(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
	dev, ok := e.deviceState(ctx, m, realm)
	if !ok {
		return
	}
	dev.armHintReset()

	format := formatForHint(dev.hint())
	if !e.retryStore(ctx, m, "emptyCache resync", func() error {
		if err := e.resendServerProperties(ctx, realm, m.DeviceID, format); err != nil {
			return err
		}
		return e.sendConsumerProperties(ctx, realm, m.DeviceID)
	}) {
		return
	}
	m.Ack()
}

// handleProducerProperties implements `control/producer/properties`
// (docs/DESIGN.md §3.3): the payload is the exhaustive zlib-compressed list
// of device-owned properties the device still holds, and every device-owned
// row not in it is purged — the device is the source of truth for its own
// properties. Malformed payloads (bad frame, zip bomb) are rejected and
// consumed.
func (e *Engine) handleProducerProperties(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
	entries, err := inflateProperties(m.Payload)
	if err != nil {
		e.reject(m, reasonControlInvalid, err.Error())
		return
	}
	dev, ok := e.deviceState(ctx, m, realm)
	if !ok {
		return
	}

	keep := e.resolvePropertyRefs(realm, dev, entries)
	if !e.retryStore(ctx, m, "producer/properties purge", func() error {
		purged, err := e.st.PurgeDeviceOwnedExcept(ctx, realm.id, m.DeviceID, keep)
		if err == nil && purged > 0 {
			e.log.Info("purged device-owned properties",
				"realm", m.Realm, "device", m.DeviceID.String(), "purged", purged, "kept", len(keep))
		}
=== engine.go 328-360 (device_connected pattern) ===

// handleLifecycle is the onLifecycle observer (after the built-in cache
// eviction): device_connected / device_disconnected triggers, live bus
// events, and — on connect — the asynchronous `consumer/properties` purge
// send. A (re)connecting device always receives the full current truth, a
// safe superset of the session_present=0 requirement of docs/DESIGN.md §3.4
// (the purge list is idempotent). The send runs off the broker's connection
// goroutine: LifecycleSink implementations must not block.
func (e *Engine) handleLifecycle(ev broker.LifecycleEvent) {
	rs := e.schemas.realm(ev.Realm)
	if rs == nil {
		return
	}
	deviceID := ev.DeviceID.String()

	switch ev.Type {
	case broker.EventDeviceConnected:
		ip := ""
		if ev.RemoteIP.IsValid() {
			ip = ev.RemoteIP.String()
		}
		e.fireDevice(rs, ev.DeviceID, ev.At,
			triggers.DeviceEvent{DeviceID: deviceID, On: triggers.OnDeviceConnected},
			triggers.NewDeviceConnectedEvent(ip))
		e.bus.Publish(stream.Event{
			Kind: stream.KindDeviceConnected, Realm: rs.name, DeviceID: deviceID, Timestamp: ev.At, IP: ip,
		})

		e.bg.Add(1)
		go func() {
			defer e.bg.Done()

$ cd /root/astrate-mule && rg -n "OnDeviceEmptyCacheReceived|OnIncomingIntrospection" internal/ pkg/ | head -10; echo "=== known-on registry ==="; rg -n "OnIncomingIntrospection|knownOn|deviceOn|validOn" internal/engine/triggers/*.go --glob '!**/*_test.go' | head -15
internal/engine/triggers/match.go:61:	// OnDeviceEmptyCacheReceived fires on control/emptyCache (accepted, not
internal/engine/triggers/match.go:63:	OnDeviceEmptyCacheReceived = "device_empty_cache_received"
internal/engine/triggers/match.go:66:	// OnIncomingIntrospection fires on every introspection publish.
internal/engine/triggers/match.go:67:	OnIncomingIntrospection = "incoming_introspection"
internal/engine/triggers/match.go:123:	OnDeviceEmptyCacheReceived: false,
internal/engine/triggers/match.go:125:	OnIncomingIntrospection:    true,
internal/engine/triggers/match.go:445:	if c.On == OnIncomingIntrospection && hasName {
internal/engine/introspection.go:77:		triggers.DeviceEvent{DeviceID: m.DeviceID.String(), On: triggers.OnIncomingIntrospection},
internal/engine/triggers/events.go:280:	return IncomingIntrospectionEvent{Type: OnIncomingIntrospection, Introspection: introspection}
=== known-on registry ===
internal/engine/triggers/events.go:280:	return IncomingIntrospectionEvent{Type: OnIncomingIntrospection, Introspection: introspection}
internal/engine/triggers/match.go:66:	// OnIncomingIntrospection fires on every introspection publish.
internal/engine/triggers/match.go:67:	OnIncomingIntrospection = "incoming_introspection"
internal/engine/triggers/match.go:115:// deviceOns enumerates the device_trigger conditions; the value records
internal/engine/triggers/match.go:117:var deviceOns = map[string]bool{
internal/engine/triggers/match.go:125:	OnIncomingIntrospection:    true,
internal/engine/triggers/match.go:406:	evaluated, known := deviceOns[c.On]
internal/engine/triggers/match.go:445:	if c.On == OnIncomingIntrospection && hasName {

$ cd /root/astrate-mule && rg -n "^\s*add\)|add\)" tools/mule.sh | head -5; rg -n "cmd_add|^add\)" -A 12 tools/mule.sh | head -30
973:  add)       shift; cmd_add "$@";;
421:cmd_add() {
422-  [ $# -gt 0 ] || die 'add needs a title: mule.sh add "slug: what to do"'
423-  printf -- '- [ ] %s\n' "$*" >> "$TODO"
424-  ok "queued: $*"
425-}
426-
427-# Full-transcript history written by save_report — filenames sort chronologically, so `ls`
428-# is the ordering. No args: list. A number: print that many full transcripts, newest first.
429-cmd_reports() {
430-  local dir="$MULE/reports/log" n="${1:-}"
431-  [ -d "$dir" ] || { note "no reports yet"; return 0; }
432-  if [ -z "$n" ]; then
433-    ls -1 "$dir" | sort | tac
--
973:  add)       shift; cmd_add "$@";;
974-  menu)      shift; cmd_menu "$@";;
975-  recipe)    shift; cmd_recipe "$@";;
976-  refill)    shift; cmd_refill "$@";;
977-  review)    shift; cmd_review "$@";;
978-  budget)    shift; cmd_budget "$@";;
979-  legion)    shift; cmd_legion "$@";;
980-  tick)      shift; cmd_tick "$@";;
981-  revert)    shift; cmd_revert "$@";;
982-  reports)   shift; cmd_reports "$@";;
983-  *) sed -n '2,19p' "$0" | sed 's/^# \{0,1\}//'; exit 1;;
984-esac

Both diffs read and the Astrate-side anchors verified. Appending the two proposals via the queue tool.
$ cd /root/astrate-mule && tools/mule.sh add 'device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implements `device_empty_cache_received` end to end — `SimpleEvent` field 23, an empty `DeviceEmptyCacheReceivedEvent` whose encoder emits exactly `{"type":"device_empty_cache_received"}` (astarte_core encoder.ex:54-62), fired by DUP at `control_handler.ex:161` *inside* the `with` that runs `send_control_consumer_properties` → `resend_all_properties` → `set_pending_empty_cache`, so it fires only when all three succeed and never on failure, and the commit deletes the `# TODO: implement on_empty_cache_received` from `astarte_events/lib/astarte_events/triggers/core.ex`. Astrate already accepts the trigger but never evaluates it: `OnDeviceEmptyCacheReceived` sits in `deviceOns` with value `false` (internal/engine/triggers/match.go:123, the `evaluated` flag read at match.go:406) and the comment'"'"'s stated reason ("upstream defines no SimpleEvent variant for it") is now false. Add the event body constructor to internal/engine/triggers/events.go (the payload is `{"type":"device_empty_cache_received"}` and nothing else — no device id, no metadata), flip the `deviceOns` value to true, and fire it from `handleEmptyCache` (internal/engine/control.go:58) with the same idiom `handleLifecycle` uses for `device_connected` (`e.fireDevice` + `stream.Publish`, internal/engine/engine.go:346-353), at the point where `m.Ack()` runs — i.e. after `resendServerProperties` and `sendConsumerProperties` have both succeeded inside `retryStore`. Add a T1 test that a trigger on it fires exactly once per accepted emptyCache and does not fire when the resend fails.'
ok queued: device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implements `device_empty_cache_received` end to end — `SimpleEvent` field 23, an empty `DeviceEmptyCacheReceivedEvent` whose encoder emits exactly `{"type":"device_empty_cache_received"}` (astarte_core encoder.ex:54-62), fired by DUP at `control_handler.ex:161` *inside* the `with` that runs `send_control_consumer_properties` → `resend_all_properties` → `set_pending_empty_cache`, so it fires only when all three succeed and never on failure, and the commit deletes the `# TODO: implement on_empty_cache_received` from `astarte_events/lib/astarte_events/triggers/core.ex`. Astrate already accepts the trigger but never evaluates it: `OnDeviceEmptyCacheReceived` sits in `deviceOns` with value `false` (internal/engine/triggers/match.go:123, the `evaluated` flag read at match.go:406) and the comment's stated reason ("upstream defines no SimpleEvent variant for it") is now false. Add the event body constructor to internal/engine/triggers/events.go (the payload is `{"type":"device_empty_cache_received"}` and nothing else — no device id, no metadata), flip the `deviceOns` value to true, and fire it from `handleEmptyCache` (internal/engine/control.go:58) with the same idiom `handleLifecycle` uses for `device_connected` (`e.fireDevice` + `stream.Publish`, internal/engine/engine.go:346-353), at the point where `m.Ack()` runs — i.e. after `resendServerProperties` and `sendConsumerProperties` have both succeeded inside `retryStore`. Add a T1 test that a trigger on it fires exactly once per accepted emptyCache and does not fire when the resend fails.

$ cd /root/astrate-mule && tools/mule.sh add 'appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.'
ok queued: appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.

$ cd /root/astrate-mule && tail -4 .mule/todo.md | cut -c1-120 && echo "---" && git status --short && echo "=== gate (container-free) ===" && gofmt -l . | head -5 && go vet ./internal/engine/... 2>&1 | tail -3
- [x] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_n
- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-c
- [ ] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no 
- [ ] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an ob
---
 M .mule/todo.md
=== gate (container-free) ===

Done: carried out `.mule/recipes/astarte-upstream.md` — no release movement since `v1.4.0-rc.6` (newest tag, 2026-09-30, already audited by the 2026-10-01/02 runs), but upstream master moved 9 commits past the rc.6 cut with two wire-visible changes; read both diffs and appended 2 task lines.
Files: `.mule/todo.md` (2 lines appended, nothing else touched)
Verified: `gh api repos/astarte-platform/astarte/releases` (newest `v1.4.0-rc.6`, prerelease, 2026-09-30; newest stable `v1.3.4`) → unchanged since the last two runs; `commits?sha=master&since=2026-09-29` → 9 commits, 2 wire-visible (#2232, #2237), 7 CI/build/chore; both diffs read (encoder/simple_event/control_handler for #2232; device.ex/interface_value.ex/fallback_controller.ex/error_view.ex for #2237); Astrate side confirmed by grep — `deviceOns[OnDeviceEmptyCacheReceived]=false` at match.go:123, no `payload.Reason` mapping in `writeError` at appengine/http.go:643 → object-shape rejects fall to 500; `gofmt -l .` and `go vet ./internal/engine/...` → clean (no code changed, so this is a formality).
Unsure: whether `ReasonMissingRequired` has the same 500-instead-of-422 hole as the object key (upstream renders a distinct 422 `missing_required_mapping`) — left as a "report, don't guess" clause inside the task line rather than asserted. Also: I did **not** append the recipe's mandatory `compat-note-<tag>` line, because the tag is unchanged and `compat-note-v1.4-rc.6` already exists and completed on 2026-10-02 (commit `0215f97`) with its delta recorded in `for-giulio.md` — a second line would be the duplicate your instructions forbid. Both new lines state explicitly that their upstream changes are master-only and in no tag. One thing the recipe does not cover, worth your call: two new commits are sitting on master untagged, so a release-driven recipe structurally cannot see them until rc.7 ships.
Follow-ups: `device-empty-cache-received-trigger`, `appengine-unexpected-object-key` (both `[auto]`, appended via `tools/mule.sh add`; no `gh issue create`, no git operations)
