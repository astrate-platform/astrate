slug: recipe-milestones
verdict: proposed
at:  f7036c5
ran: 2026-09-18T11:00:39Z on DietPi in 379s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
Recipe read. Milestone to work: **v3.0** (first non-DONE), scope decided 2026-08-22. Now investigating per step 3.

$ gh issue list --label milestone-3.0 --state all --limit 50 2>&1; echo "---"; gh label list --limit 100 2>&1 | grep -i milestone
---
milestone-2.0	astarte-flow feature parity milestone	#0E8A16
milestone-4.0		#5319e7
milestone-3.0		#5319e7

✱ Glob "**/upstream-parity*.md" 0 matches

$ gh issue list --label upstream-parity --state all --limit 100 --json number,title,state,labels 2>&1
[{"labels":[{"id":"LA_kwDORmfsJs8AAAACte8Saw","name":"mule-review","description":"Mule pushed a change for this; waiting for a human or a strong model to review it","color":"0e8a16"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":93,"state":"OPEN","title":"In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"},{"id":"LA_kwDORmfsJs8AAAACxx3GJg","name":"upstream-experimental","description":"Adopted from an upstream experimental/unstable feature — review each upstream release","color":"d4c5f9"}],"number":92,"state":"OPEN","title":"keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte8SoA","name":"mule-blocked","description":"The mule tried this and could not do it; see its comment","color":"d93f0b"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":91,"state":"CLOSED","title":"Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":89,"state":"CLOSED","title":"Dashboard flow-block schema mismatch (split_map/virtual pools hardcoded; null_sink/log_sink unknown)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":88,"state":"CLOSED","title":"Flow auth: support a_f JWT claim"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":87,"state":"CLOSED","title":"Flow block: lua_map — needs embedded Lua runtime (parked)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":86,"state":"CLOSED","title":"Flow: pipeline source DSL — keep DAG-JSON as documented deviation?"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":85,"state":"CLOSED","title":"Flow API: user-defined composite blocks"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":84,"state":"CLOSED","title":"Flow blocks: virtual_device_pool / dynamic_virtual_device_pool"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":83,"state":"CLOSED","title":"Flow blocks: mqtt_source/mqtt_sink (+modbus_tcp_source?) demand-driven"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":82,"state":"CLOSED","title":"Flow blocks: http_source/http_sink (demand-driven)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":81,"state":"CLOSED","title":"Flow block: json_path_map"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":80,"state":"CLOSED","title":"Flow blocks: pure-transform set (to_json, update_metadata, split_map, random_source, sort)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":79,"state":"CLOSED","title":"Verify registration-limit-reached HTTP status vs upstream"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACtvELrg","name":"milestone-4.0","description":"","color":"5319e7"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":78,"state":"OPEN","title":"FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":77,"state":"CLOSED","title":"Verify per-service version endpoints (GET /version, GET /v1/{realm}/version) served everywhere"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":76,"state":"CLOSED","title":"Housekeeping: GET /v1/realm-defaults/replication — decide reject/deviate (Cassandra-shaped)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":75,"state":"CLOSED","title":"Housekeeping: decide realm-deletion gating/preconditions vs always-sync deviation"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":74,"state":"CLOSED","title":"Housekeeping: PATCH /v1/realms/{realm} (jwt key, registration limit, retention; null=unset)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":73,"state":"CLOSED","title":"Housekeeping: default datastream retention injection env var (upstream 1.4)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":72,"state":"CLOSED","title":"Realms: datastream_maximum_storage_retention ceiling (create/patch/enforce)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":71,"state":"CLOSED","title":"Pairing: realm-scoped health check GET /v1/{realm}/health (upstream 1.3)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":70,"state":"CLOSED","title":"Triggers: audit wildcard semantics (interface_name '*', match_path '/*' forcing rules)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":69,"state":"CLOSED","title":"Verify unknown-realm HTTP status on RM endpoints against upstream"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACte8SoA","name":"mule-blocked","description":"The mule tried this and could not do it; see its comment","color":"d93f0b"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"},{"id":"LA_kwDORmfsJs8AAAACxx3GJg","name":"upstream-experimental","description":"Adopted from an upstream experimental/unstable feature — review each upstream release","color":"d4c5f9"}],"number":68,"state":"CLOSED","title":"Decide async_operation=false params vs documented always-sync deviation"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"},{"id":"LA_kwDORmfsJs8AAAACxx3GJg","name":"upstream-experimental","description":"Adopted from an upstream experimental/unstable feature — review each upstream release","color":"d4c5f9"}],"number":67,"state":"CLOSED","title":"Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":66,"state":"CLOSED","title":"Realm Management: detailed=true interface listing with full mappings (upstream 1.4)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":65,"state":"CLOSED","title":"Policies: handler-overlap rejection + retry_times coupling + prefetch_count"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":64,"state":"CLOSED","title":"Triggers: decide AMQP action behavior (validate-reject vs NATS-forward deviation)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":63,"state":"CLOSED","title":"Triggers: HTTP action validation limits (URL/method/header blocklist/template size)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":62,"state":"CLOSED","title":"Realm Management: audit install/update/delete error codes and statuses"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":61,"state":"CLOSED","title":"Realm Management: audit interface/mapping validation matrix against astarte_core"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":60,"state":"CLOSED","title":"Realm Management: GET config/datastream_maximum_storage_retention (since upstream 1.2.0)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":59,"state":"CLOSED","title":"AppEngine: group create-body validation + UUID-v1 from_token for group device listing"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":58,"state":"CLOSED","title":"AppEngine: PATCH requires Content-Type application/merge-patch+json"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":57,"state":"CLOSED","title":"AppEngine: audit server-write error taxonomy against upstream"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":56,"state":"CLOSED","title":"AppEngine: object-aggregate response formats table and disjoint_tables"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":55,"state":"CLOSED","title":"AppEngine: complete data-query parameter set"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":54,"state":"CLOSED","title":"AppEngine: GET /groups/{g} show-group endpoint"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":53,"state":"CLOSED","title":"AppEngine: group-scoped device PATCH + interface data access under groups/{g}/devices/{d}"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":52,"state":"CLOSED","title":"AppEngine: devices-by-alias full mirror (PATCH via alias + interface data access)"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":51,"state":"CLOSED","title":"keyAgreement control topic + encrypted mappings (upstream 1.4 experimental) — investigate/park"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":50,"state":"CLOSED","title":"Broker/engine hardening: discard malformed introspection/purge/capabilities instead of crashing"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":49,"state":"CLOSED","title":"MQTT v1: binaryblob properties re-send as BSON subtype 0 in consumer/properties"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":48,"state":"CLOSED","title":"MQTT v1: tolerate devices connecting with an empty introspection"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":47,"state":"CLOSED","title":"MQTT v1: accept device capabilities incl. purge_properties_compression_format (plaintext purge)"}]

Most of #47–#89 are closed. Open `upstream-parity` issues: #92 (keyAgreement decision), #93 (mule-review), #78 (parked, milestone-4.0). The research doc isn't where milestones.md named it — let me find it and read the experimental chart.
$ ls -la .mule/research/ 2>&1; echo "---"; ls .mule/ 2>&1
ls: cannot access '.mule/research/': No such file or directory
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
tasks
todo.md

→ Read docs/UPSTREAM-EXPERIMENTAL.md
`.mule/research/` doesn't exist — the source doc named in milestones.md is gone. Let me search for it and check for any mirror, plus look at the open #92 body and the APICompatVersion decision #90.

$ find . -iname '*parity*' -not -path './.git/*' 2>/dev/null; echo "---git log---"; git log --oneline -5 2>&1
./internal/appengine/formats_parity_test.go
./internal/appengine/groups_parity_test.go
./internal/engine/mqttv1_parity_test.go
./docs/handoff/flow-parity-audit-2026-07-29.md
---git log---
f7036c5 mule: recipe hygiene timed out
bb8bd4e mule: log docs-sync-rm-triggers-422-nested-envelope
8bfedd6 mule: docs-sync-rm-triggers-422-nested-envelope [auto]: document the nested changeset envelope on `POST /realmmanagement/v1/{realm}/triggers` in docs/api/astarte_realm_management_api.yaml — when trigger compile fails, `createTrigger` answers 422 with the upstream nested body `{"errors":{"action":{...},"simple_triggers":[...]}}` (internal/realm/http.go:252-258, triggerErrorBody at http.go:268-285, issues #63/#70), but the spec's createTrigger 422 (yaml:344-345) refs only the flat `ValidationError` component (`{"errors":{"detail":...}}`) — which the handler also emits for plain ErrValidation (http.go:259-260, e.g. "trigger requires a name", service.go:429). Model the response as oneOf across the flat ErrorDetail and the nested envelope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
4807952 mule: log docs-sync-rm-policies-delete-422
c0521bc mule: docs-sync-rm-policies-delete-422 [auto]: add the missing `422` response to `DELETE /realmmanagement/v1/{realm}/policies/{name}` in docs/api/astarte_realm_management_api.yaml — `Service.DeletePolicy` refuses to delete a policy still referenced by a trigger (`"policy %q is still used by trigger %q"`, internal/realm/service.go:573-574, wrapped as ErrValidation) and `writeError` answers 422 for it (internal/realm/http.go:410-411, detail stripped by validationDetail), but the spec lists only 204/401/404/500 for the delete op (yaml:695-713). Add the 422 ValidationError response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ gh issue view 92 --json number,title,state,body,labels 2>&1; echo "===90==="; gh issue view 90 --json number,title,state,body,labels 2>&1
{"body":"Issue #51 was closed on 2026-08-22 with: *\"Full keyAgreement/encrypted-mappings support stays parked until the upstream 1.4 experimental spec stabilizes — reopen or file fresh when it does.\"* Filing fresh: that condition has now partly fired.\n\n**What changed.** Upstream commit `d084308` (\"docs: add encrypted endpoints key agreement (#2067)\", 2026-08-31) added `doc/pages/architecture/082-key_agreement_protocol.md` — a 267-line full wire protocol — and rewrote the Key Agreement section of `080-mqtt-v1-protocol.md`.\n\n**Verified against upstream `master` today, not taken from the survey:**\n\n- `082-key_agreement_protocol.md` exists and is 13,832 bytes.\n- The sentences the old parking rested on — \"reserved and routed correctly, but the handshake protocol is not yet implemented\" and \"accepted and acknowledged by Astarte, but no response is sent until the feature is fully implemented\" — are **gone from both 080 and 082** (grep count 0 in each).\n- 082 now states the handshake \"is implemented on the Astarte side (`astarte_data_updater_plant`)\", with only device-side session-key usage still described as liable to evolve.\n\n**What the spec now pins down**, where before there was only a sketch: five QoS-2 topics `control/keyAgreement/0..4` (InitExchange, ExchangeResp, SecretHash, HashOk, ExchangeFailed — the 080 table moved from a single `control/keyAgreement` at QoS 1); CBOR message bodies with CDDL given; `alg` 0 = `ECDH_P256-HKDF_SHA256-AES_256_GCM`; a CBOR-wrapped COSE_Key pubkey plus a 32-byte HkdfSalt; session-scoped keys never reused after reconnect or clean reset; AES-256-GCM whole-document ciphertext; and on \"no key established yet\" Astarte logs and **discards** rather than disconnecting. `ExchangeFailed` error codes are enumerated (0 internal, 1 invalid argument, 2 hash mismatch, 3 unprocessable entity).\n\n**The open question is scope, not feasibility, so this is not a mule task.** The spec is published and implemented upstream, but still ships only in `v1.4.0-rc.5` — `v1.3.3` remains the newest stable tag, and Astrate targets 1.2.2. So \"stabilized\" is true of the document and false of the release. Three ways to go:\n\n1. Implement now against the rc spec — largest surface in the parity backlog (CBOR codec, X25519/P-256, HKDF, AES-256-GCM, a 5-state handshake machine, shared-secret persistence, five new error names), and it would be built against a spec that may still move.\n2. Keep parking until `v1.4.0` is a stable tag, and re-file then.\n3. Do the narrow piece only — see the separate mule issue for the stale ACL comment, which is worth fixing regardless of which of these is chosen.\n\nRecorded in `.mule/for-giulio.md`. Survey source: `survey-2026-08-31.md` on `mule/research`.","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"},{"id":"LA_kwDORmfsJs8AAAACxx3GJg","name":"upstream-experimental","description":"Adopted from an upstream experimental/unstable feature — review each upstream release","color":"d4c5f9"}],"number":92,"state":"OPEN","title":"keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds"}
===90===
{"body":"Reminder of a pending decision parked during the upstream-parity investigation (#47–#89). **Do not raise the version now** — this is a milestone-4.0 item.\n\n## Context\n\n`realm.APICompatVersion` (`internal/realm/service.go`) is currently hardcoded to `\"1.2.2\"` — the emulated upstream API level Astrate declares.\n\nEdgehog 0.13.1 gates just-in-time provisioning triggers (`device_registered`, `device_deletion_started`) behind a Realm Management version >= 1.3. Astrate's engine already supports those trigger kinds end-to-end — verified in #28 (Edgehog compatibility investigation) — but Edgehog refuses to use them because of the declared version string.\n\n## What this issue must cover when picked up\n\n1. **Audit first**: what else does declaring \"1.3\" imply for other consumers?\n   - `astarte-dashboard` is currently validated against 1.2.2 behavior (see #13, error-name bundle comparison). Check whether any dashboard code path feature-gates on RM version besides the policies UI gate (>= 1.1.1).\n   - Sweep the parity findings (#47–#89): several are \"1.3/1.4 features\" we deliberately do NOT implement yet — does declaring 1.3 overpromise any of them?\n   - Check device SDKs and `astartectl` for version-gated behavior.\n2. **Then decide one of**:\n   - bump the constant to `\"1.3.0-rc.0\"` (upstream's own pre-release numbering for that line);\n   - or make it configuration (`internal/config`) so operators can declare the level their consumers expect.\n3. Update `docs/COMPATIBILITY.md` deviation #10 wording accordingly (propose, don't edit directly).\n\n## Outcome if done today (why it's parked)\n\nBumping without the audit risks silently promising surfaces we haven't verified; not bumping keeps Edgehog's JIT provisioning disabled even though the engine supports it. Milestone-4.0 scope decision.","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACtvELrg","name":"milestone-4.0","description":"","color":"5319e7"}],"number":90,"state":"CLOSED","title":"milestone-4.0: decide on bumping APICompatVersion from \"1.2.2\" to \"1.3.0-rc.0\" (or making it configurable)"}

Sources narrowed: the research doc lives on `mule/research` (branch, not this checkout), but the actual v3.0 backlog is tracked on GitHub and is nearly all closed. Let me check what happened to closed #90, the current APICompatVersion, upstream 1.4 final status, and the current for-giulio state.
$ rg -n "APICompatVersion" internal/ docs/ --include:'*.go' --include:'*.md' 2>/dev/null | head -20
(no output)

→ Read .mule/for-giulio.md

$ gh issue view 90 --json closedAt,closingReviewRequest --jq '{closedAt, closingReviewRequest}' 2>&1; echo "===events==="; gh api repos/astrate-platform/astrate/issues/90/events --jq '.[] | {event, actor: .actor.login, created_at}' 2>&1 | tail -20
Unknown JSON field: "closingReviewRequest"
Available fields:
  assignees
  author
  blockedBy
  blocking
  body
  closed
  closedAt
  closedByPullRequestsReferences
  comments
  createdAt
  id
  isPinned
  issueType
  labels
  milestone
  number
  parent
  projectCards
  projectItems
  reactionGroups
  state
  stateReason
  subIssues
  subIssuesSummary
  title
  updatedAt
  url
===events===
{"actor":"atsetilam","created_at":"2026-08-22T02:28:21Z","event":"labeled"}
{"actor":"atsetilam","created_at":"2026-08-22T02:28:21Z","event":"labeled"}
{"actor":"atsetilam","created_at":"2026-08-22T02:48:19Z","event":"referenced"}
{"actor":"atsetilam","created_at":"2026-08-23T15:50:31Z","event":"referenced"}
{"actor":"atsetilam","created_at":"2026-08-23T21:06:35Z","event":"closed"}

Prior runs verify v3.0 is gated on upstream stable v1.4.0 (not yet shipped as of 2026-09-15). Let me re-check that gate and today's open issue set.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | {tag_name, prerelease, published_at}' 2>&1 | head -30
{"prerelease":false,"published_at":"2026-09-18T09:42:57Z","tag_name":"v1.3.4"}
{"prerelease":true,"published_at":"2026-08-20T12:28:04Z","tag_name":"v1.4.0-rc.5"}
{"prerelease":false,"published_at":"2026-08-07T14:33:58Z","tag_name":"v1.3.3"}
{"prerelease":true,"published_at":"2026-08-06T13:40:20Z","tag_name":"v1.4.0-rc.4"}
{"prerelease":true,"published_at":"2026-07-31T10:48:39Z","tag_name":"v1.4.0-rc.3"}
{"prerelease":true,"published_at":"2026-07-14T10:16:22Z","tag_name":"v1.4.0-rc.2"}
{"prerelease":true,"published_at":"2026-07-14T08:53:27Z","tag_name":"v1.4.0-rc.1"}
{"prerelease":false,"published_at":"2026-07-14T10:07:26Z","tag_name":"v1.3.2"}
{"prerelease":false,"published_at":"2026-07-13T11:50:05Z","tag_name":"v1.3.1"}
{"prerelease":false,"published_at":"2026-05-06T15:30:50Z","tag_name":"v1.3.0"}
{"prerelease":false,"published_at":"2026-04-27T16:39:28Z","tag_name":"v1.2.2"}
{"prerelease":true,"published_at":"2026-04-08T12:32:33Z","tag_name":"v1.4.0-rc.0"}
{"prerelease":true,"published_at":"2026-04-08T16:34:59Z","tag_name":"v1.3.0-rc.2"}
{"prerelease":true,"published_at":"2026-04-08T16:31:03Z","tag_name":"v1.2.2-rc.0"}
{"prerelease":false,"published_at":"2026-03-12T16:11:20Z","tag_name":"v1.2.1"}
{"prerelease":true,"published_at":"2026-02-13T10:36:20Z","tag_name":"v1.2.1-rc.1"}
{"prerelease":true,"published_at":"2026-01-26T15:58:33Z","tag_name":"v1.3.0-rc.1"}
{"prerelease":true,"published_at":"2025-11-21T14:58:46Z","tag_name":"v1.3.0-rc.0"}
{"prerelease":true,"published_at":"2025-08-27T07:03:28Z","tag_name":"v1.2.1-rc.0"}
{"prerelease":true,"published_at":"2025-04-10T15:31:04Z","tag_name":"v1.2.1-alpha.0"}
{"prerelease":false,"published_at":"2024-07-02T14:09:07Z","tag_name":"v1.2.0"}
{"prerelease":true,"published_at":"2024-06-11T13:02:33Z","tag_name":"v1.2.0-rc.0"}
{"prerelease":false,"published_at":"2024-04-23T21:15:57Z","tag_name":"v1.0.6"}
{"prerelease":false,"published_at":"2023-11-15T15:29:57Z","tag_name":"v1.1.1"}
{"prerelease":false,"published_at":"2023-09-26T10:27:57Z","tag_name":"v1.0.5"}
{"prerelease":false,"published_at":"2023-06-20T18:59:09Z","tag_name":"v1.1.0"}
{"prerelease":true,"published_at":"2023-06-09T15:49:44Z","tag_name":"v1.1.0-rc.0"}
{"prerelease":true,"published_at":"2022-11-24T16:38:37Z","tag_name":"v1.1.0-alpha.0"}
{"prerelease":false,"published_at":"2022-10-25T16:50:21Z","tag_name":"v1.0.4"}
{"prerelease":false,"published_at":"2022-07-05T13:57:17Z","tag_name":"v1.0.3"}

$ echo "=== milestone-3.0 open ==="; gh issue list --label milestone-3.0 --state open 2>&1; echo "=== upstream-parity open ==="; gh issue list --label upstream-parity --state open 2>&1
=== milestone-3.0 open ===
=== upstream-parity open ===
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z

New stable today: **v1.3.4** (patch, 2026-09-18). v1.4.0 is still rc.5. Let me check what's in v1.3.4 and whether any EXPERIMENTAL rows are affected.
$ rg -n "APICompatVersion" . --include='*.go' --include='*.md' 2>/dev/null | head -20
(no output)

$ gh api repos/astarte-platform/astarte/releases/tags/v1.3.4 --jq '.body' 2>&1 | head -80
### Fixed

- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
  uncollected binaries over time by forcing more frequent full sweep garbage collections
  on them.
- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.

v1.3.4 is a pure ops patch — no wire impact. Let me check the APICompatVersion constant's actual location.
$ rg -n "CompatVersion|compatVersion|apiCompat" internal/ 2>/dev/null | head; echo "---"; rg -rn "1\.2\.2" internal/realm/ 2>/dev/null | head
internal/realm/http.go:119:// feature UI on it — see APICompatVersion).
internal/realm/http.go:121:	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
internal/realm/dashboard_compat_test.go:35:		if v != APICompatVersion {
internal/realm/dashboard_compat_test.go:36:			t.Errorf("version = %q, want %q", v, APICompatVersion)
internal/realm/service.go:580:// APICompatVersion is the upstream Realm Management API level Astrate
internal/realm/service.go:588:const APICompatVersion = "1.2.2"
---
internal/realm/dashboard_compat_test.go:// Dashboard vn requires: version, device_registration_limit, delivery
internal/realm/service.go:const APICompatVersion = "n"

Let me confirm the full open-issue set hasn't changed since the last run and nothing new appeared around this milestone.
$ gh issue list --state open --limit 50 --json number,title,labels --jq '.[] | "\(.number)\t\(.labels|map(.name)|join(","))\t\(.title)"' 2>&1
106	mule-alarm	mule: nothing has landed in 14h
105	mule-alarm	mule: nothing has landed in 14h
104	mule-alarm	mule: nothing has landed in 14h
103	mule-alarm	mule: nothing has landed in 16h
102	mule-alarm	mule: nothing has landed in 14h
101	mule-alarm	mule: nothing has landed in 14h
100	mule-alarm	mule: nothing has landed in 15h
99	mule-alarm	mule: nothing has landed in 14h
98	mule-alarm	mule: nothing has landed in 16h
97	mule-alarm	mule: nothing has landed in 15h
96	mule-alarm	mule: nothing has landed in 14h
95	mule-alarm	mule: nothing has landed in 15h
94	mule-alarm	mule: nothing has landed in 15h
93	mule-review,upstream-parity	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
92	enhancement,upstream-parity,upstream-experimental	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds
78	enhancement,milestone-4.0,upstream-parity	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
1	wontfix	Provide an Open Source IoT Platform unironically, unironically

The open set is unchanged — no new gaps. The only new fact since the 2026-09-15 run is today's v1.3.4 patch (GC/RPC-only, no wire impact). Let me confirm todo.md has nothing I'd duplicate.
→ Read .mule/todo.md
The open set is unchanged; the only new upstream fact is today's `v1.3.4` patch, which is maintenance-only (data_updater_plant GC/RPC), no wire impact. Let me check the experimental register once more and record the run outcome in for-giulio.md.

The register (`docs/UPSTREAM-EXPERIMENTAL.md`) still has both rows waiting on **v1.4 final** — v1.3.4 doesn't trigger them. Verification complete; the milestone remains gated on a still-unshipped stable v1.4.0. Recording the run in for-giulio.md (newest at top):

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) shows a new **stable v1.3.4** (2026-09-18) — maintenance only (data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change), so no new gap for Astrate — and still **no stable v1.4.0**: newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+
+---
+
 - **docs-sync realm-management, 2026-09-17: `docs/site/realm-management-api.md` documents an "Update trigger" endpoint Astrate does not serve.** Lines 81-86 show `PUT /realmmanagement/v1/<realm>/triggers/<name>` with a `{ "data": <updated trigger JSON> }` body, but no such route is registered — internal/realm/http.go:44-47 wires only GET/POST `/triggers` and GET/DELETE `/triggers/{name}` (an update op is deliberately absent; not even a stub or TODO exists). Site prose is yours — drop the "Update trigger" section or confirm the route is planned.
 
 ---
 

Run complete. The v3.0 milestone record is updated in `.mule/for-giulio.md` with today's verification.

**Summary:**
- New upstream `v1.3.4` shipped today (maintenance patch, no wire/API impact) — still **no stable `v1.4.0`** (newest: `v1.4.0-rc.5`)
- `milestone-3.0` label empty; open `upstream-parity` issues unchanged: **#92** (keyAgreement decision, parked on stable v1.4.0), **#93** (mule-review, its own path), **#78** (FDO, milestone-4.0)
- UPSTREAM-EXPERIMENTAL rows #67/#68 still wait on v1.4 final
- `APICompatVersion` stays `1.2.2` per #90's frozen decision
- No new gaps, no issues filed, no task lines proposed — deliberately not proposing "complete, cut the tag" (milestone targets 1.3/**1.4**, gate unresolved)
