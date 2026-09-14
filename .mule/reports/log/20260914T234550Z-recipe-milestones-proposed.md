slug: recipe-milestones
verdict: proposed
at:  469597f
ran: 2026-09-14T23:45:50Z on DietPi in 95s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
The current milestone is **v3.0** (upstream parity, backed by issues #47–#89). Let me read the research doc and survey the existing issues.

$ gh issue list --label upstream-parity --state all --limit 50 --json number,title,labels,state,body 2>&1 | head -c 20000
[{"body":"Comment-only change. The deny itself stays; only its stated justification is stale.\n\n`internal/broker/aclhook.go` currently reads:\n\n```go\n// control/keyAgreement (080-mqtt-v1-protocol.md:172-176) is deliberately\n// unsupported: upstream itself documents the handshake as \"reserved and\n// routed correctly, but not yet implemented\", so falling through to the\n// deny path below matches upstream's own status rather than being a gap.\n```\n\nThat reasoning no longer holds on either half. Verified against upstream `master`:\n\n- The quoted sentence, and the companion \"accepted and acknowledged by Astarte, but no response is sent until the feature is fully implemented\", are **gone from both** `080-mqtt-v1-protocol.md` and the new `082-key_agreement_protocol.md` (grep count 0 in each).\n- `082-key_agreement_protocol.md` (13,832 bytes, added by upstream `d084308` on 2026-08-31) states the handshake \"is implemented on the Astarte side (`astarte_data_updater_plant`)\".\n- The 080 topic table row also changed shape: a single `control/keyAgreement` at QoS 1 became `control/keyAgreement/<message id>` at QoS 2.\n\nSo the comment now tells a reader something false about upstream, and a future reader could reasonably \"fix the gap\" on the strength of it.\n\n**Do:** rewrite the comment so it stands on Astrate's own decision — the topic is denied because full key agreement is out of scope for the current 1.2.2 target and the spec ships only in `v1.4.0-rc.5` (`v1.3.3` is the newest stable tag) — cite `082-key_agreement_protocol.md` as the current upstream reference rather than the deleted 080 sentence, and point at issue #92 where the implement-or-keep-parking decision lives.\n\n**Do not:** change the ACL behaviour, and do not touch the expectation at `internal/broker/aclhook_test.go:53` (`pub control keyAgreement` still denies).\n\n**Acceptance:** `go build ./...` and `go test ./internal/broker/...` pass unchanged; the comment no longer contains the phrase \"not yet implemented\"; `aclhook_test.go:53` is untouched.\n\nSurvey source: `survey-2026-08-31.md` on `mule/research`.","labels":[{"id":"LA_kwDORmfsJs8AAAACte8Saw","name":"mule-review","description":"Mule pushed a change for this; waiting for a human or a strong model to review it","color":"0e8a16"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":93,"state":"OPEN","title":"In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged"},{"body":"Issue #51 was closed on 2026-08-22 with: *\"Full keyAgreement/encrypted-mappings support stays parked until the upstream 1.4 experimental spec stabilizes — reopen or file fresh when it does.\"* Filing fresh: that condition has now partly fired.\n\n**What changed.** Upstream commit `d084308` (\"docs: add encrypted endpoints key agreement (#2067)\", 2026-08-31) added `doc/pages/architecture/082-key_agreement_protocol.md` — a 267-line full wire protocol — and rewrote the Key Agreement section of `080-mqtt-v1-protocol.md`.\n\n**Verified against upstream `master` today, not taken from the survey:**\n\n- `082-key_agreement_protocol.md` exists and is 13,832 bytes.\n- The sentences the old parking rested on — \"reserved and routed correctly, but the handshake protocol is not yet implemented\" and \"accepted and acknowledged by Astarte, but no response is sent until the feature is fully implemented\" — are **gone from both 080 and 082** (grep count 0 in each).\n- 082 now states the handshake \"is implemented on the Astarte side (`astarte_data_updater_plant`)\", with only device-side session-key usage still described as liable to evolve.\n\n**What the spec now pins down**, where before there was only a sketch: five QoS-2 topics `control/keyAgreement/0..4` (InitExchange, ExchangeResp, SecretHash, HashOk, ExchangeFailed — the 080 table moved from a single `control/keyAgreement` at QoS 1); CBOR message bodies with CDDL given; `alg` 0 = `ECDH_P256-HKDF_SHA256-AES_256_GCM`; a CBOR-wrapped COSE_Key pubkey plus a 32-byte HkdfSalt; session-scoped keys never reused after reconnect or clean reset; AES-256-GCM whole-document ciphertext; and on \"no key established yet\" Astarte logs and **discards** rather than disconnecting. `ExchangeFailed` error codes are enumerated (0 internal, 1 invalid argument, 2 hash mismatch, 3 unprocessable entity).\n\n**The open question is scope, not feasibility, so this is not a mule task.** The spec is published and implemented upstream, but still ships only in `v1.4.0-rc.5` — `v1.3.3` remains the newest stable tag, and Astrate targets 1.2.2. So \"stabilized\" is true of the document and false of the release. Three ways to go:\n\n1. Implement now against the rc spec — largest surface in the parity backlog (CBOR codec, X25519/P-256, HKDF, AES-256-GCM, a 5-state handshake machine, shared-secret persistence, five new error names), and it would be built against a spec that may still move.\n2. Keep parking until `v1.4.0` is a stable tag, and re-file then.\n3. Do the narrow piece only — see the separate mule issue for the stale ACL comment, which is worth fixing regardless of which of these is chosen.\n\nRecorded in `.mule/for-giulio.md`. Survey source: `survey-2026-08-31.md` on `mule/research`.","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"},{"id":"LA_kwDORmfsJs8AAAACxx3GJg","name":"upstream-experimental","description":"Adopted from an upstream experimental/unstable feature — review each upstream release","color":"d4c5f9"}],"number":92,"state":"OPEN","title":"keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds"},{"body":"Found during the 2026-08-24 verify batch (#57/#69/#77/#79 probes): upstream Astarte 1.2.0 answers\n\n- `GET /pairing/health` → 200\n- `GET /pairing/v1/{realm}/health` → 404 Page not found (bench and scratch realms both)\n\nAstrate implements the unauthenticated health route at `GET /pairing/v1/{realm}/health` (landed in realmcfg-02, issue #71, commit 5111e3c) — the path shape does not exist upstream. The dashboard polls per-service health at the service root, so the current route may be invisible to it.\n\nDecide: move/duplicate the handler at `GET /pairing/health` (keeping or dropping the v1/{realm} variant), and re-check whatever probe led realmcfg-02 to the v1 shape — that session's stack had a port override in play, so the original measurement may have hit a different router.","labels":[{"id":"LA_kwDORmfsJs8AAAACte8SoA","name":"mule-blocked","description":"The mule tried this and could not do it; see its comment","color":"d93f0b"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":91,"state":"CLOSED","title":"Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health"},{"body":"Dashboard v1.3.3 ships schemas for container, split_map, virtual_device_pool, dynamic_virtual_device_pool as native blocks, and doesn't know Astrate's null_sink/log_sink — friction in the Flow UI pages. Options: serve block schemas from the blocks-discovery API and/or contribute upstream; at minimum document which blocks render cleanly.\nSize cosmetic-M.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":89,"state":"CLOSED","title":"Dashboard flow-block schema mismatch (split_map/virtual pools hardcoded; null_sink/log_sink unknown)"},{"body":"Upstream flow uses its own `a_f` authorization claim (GuardianAuthorizePath against {realm}/{resource}); Astrate's flow API rides on `a_rma`. Additive change to the claim model; check Dashboard token parsing expectations too.\nSize S, confirmed (flow guardian.ex).\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":88,"state":"CLOSED","title":"Flow auth: support a_f JWT claim"},{"body":"Upstream transforms messages with embedded Lua 5.2 scripts. Already parked out of v2.0; kept visible. Requires a Go-embedded Lua interpreter (gopher-lua) — nontrivial dependency decision.\nSize L, parked.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":87,"state":"CLOSED","title":"Flow block: lua_map — needs embedded Lua runtime (parked)"},{"body":"Upstream pipelines carry a custom DSL string (`source`) parsed by a lexer/parser plus a Draft-4 JSON Schema validating flow-level config. Astrate uses serialisable DAG JSON instead — arguably better, but a visible operator-facing difference. Decide and record as deviation (or plan a DSL front-end later).\nDecision issue; size L if ever implemented.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":86,"state":"CLOSED","title":"Flow: pipeline source DSL — keep DAG-JSON as documented deviation?"},{"body":"Upstream `POST /blocks` defines a reusable named sub-chain (name + pipeline-source string + type + schema), stored per-realm, cannot shadow defaults. Astrate has pipelines but no composite-block concept.\nSize M/L, confirmed (blocks.ex).\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":85,"state":"CLOSED","title":"Flow API: user-defined composite blocks"},{"body":"Consumer blocks writing to Astarte as one or more virtual devices, deriving device/interface/path from the message key; the dynamic variant auto-registers virtual devices on demand with persisted credentials. Genuinely useful for gateways/simulations. The dashboard hardcodes their schemas as native blocks.\nSize M/L, confirmed guides + dashboard definitions.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":84,"state":"CLOSED","title":"Flow blocks: virtual_device_pool / dynamic_virtual_device_pool"},{"body":"I/O-bound upstream blocks: subscribe/publish over MQTT; modbus polls TCP registers. Only worth doing when someone actually needs them — tracked so the gap is visible.\nSize M each.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":83,"state":"CLOSED","title":"Flow blocks: mqtt_source/mqtt_sink (+modbus_tcp_source?) demand-driven"},{"body":"Upstream has `http_source` (poll GET URLs round-robin) and `http_sink` (POST binary payload). http_sink was already parked out of v2.0 as demand-driven; keeping both tracked together here until demand shows up.\nSize M each, confirmed.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":82,"state":"CLOSED","title":"Flow blocks: http_source/http_sink (demand-driven)"},{"body":"Reshape JSON messages via a template with $MESSAGE/metadata path substitution. Pure transform.\nSize M, confirmed guide.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":81,"state":"CLOSED","title":"Flow block: json_path_map"},{"body":"Five easy native blocks toward catalog parity with astarte_flow (13 registered blocks upstream):\n- `to_json`: serialize payload to JSON binary (S)\n- `update_metadata`: add/remove/rewrite metadata keys (S)\n- `split_map`: split a map message into one message per key (S)\n- `random_source`: interval-emitted random int/real/bool (S)\n- `sort`: stateful reorderer for out-of-order streams, optional dedup (M)\nAll pure logic except random_source's timer. One issue each or batch — mule's choice.\nSize S-M total, confirmed catalog.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":80,"state":"CLOSED","title":"Flow blocks: pure-transform set (to_json, update_metadata, split_map, random_source, sort)"},{"body":"When the realm registration limit is hit, upstream's engine returns `:device_registration_limit_reached` but no error-view mapping clause was found — likely surfaces as an unmapped fallback **500**. If Astrate returns something saner, document as deliberate deviation.\nSize S verification, low confidence on upstream behavior.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":79,"state":"CLOSED","title":"Verify registration-limit-reached HTTP status vs upstream"},{"body":"Upstream 1.3 experimental, default-on in 1.4: TO2 CBOR endpoints under `/v1/{realm}/fdo/101/msg/{60..70}` (auto-registers the device mid-flow, hands over credentials_secret), owner side `POST/GET .../fdo/owner_keys` + ownership vouchers with TO0 claim, secrets in Vault/OpenBao, rendezvous server required.\n\nLarge surface; needs a product call (who runs a rendezvous server?) before any implementation. Parked.\nSize L.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACtvELrg","name":"milestone-4.0","description":"","color":"5319e7"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":78,"state":"OPEN","title":"FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)"},{"body":"Each upstream service exposes unauthenticated `GET /version` and realm-scoped `GET /v1/{realm}/version` returning `{\"data\": \"<version>\"}`. The dashboard hits these at login/startup on all four services (RM, AppEngine, pairing, flow). Only RM's realm-scoped endpoint is known-good in Astrate. Audit and fill gaps; keep the APICompatVersion convention consistent.\nSize S verification.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":77,"state":"CLOSED","title":"Verify per-service version endpoints (GET /version, GET /v1/{realm}/version) served everywhere"},{"body":"New in 1.4: returns the default keyspace replication for new realms ({replication_class, replication_factor} or network_topology map). Meaningless without Cassandra/Scylla keyspaces. Likely wontfix-with-note in COMPATIBILITY.md, but record the decision.\nSize S decision.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":76,"state":"CLOSED","title":"Housekeeping: GET /v1/realm-defaults/replication — decide reject/deviate (Cassandra-shaped)"},{"body":"Upstream: cluster flag off → DELETE realm returns **405** `realm_deletion_disabled`; connected devices present → **422** `connected_devices_present`; deletion is idempotent/resumable. Astrate deletes synchronously regardless of connections. Decide: match (cheap status codes) or document as deviation.\nSize S decision.\\n\\n---\\nInvestigation: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 vs Astrate target 1.2.2).","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"number":75,"state":"CLOSED","title":"Housekeeping: decide realm-deletion gating/preconditions vs always-sync deviation"},{"body":"Present well

✱ Glob ".mule/research/**/*.md" 0 matches
The research doc isn't at that path. Let me find it and get the full issue list.

$ gh issue list --label upstream-parity --state all --limit 50 --json number,title,labels,state > /tmp/opencode/upstream-parity-issues.json 2>&1; echo "exit=$?"; wc -c /tmp/opencode/upstream-parity-issues.json
exit=0
18781 /tmp/opencode/upstream-parity-issues.json

✱ Glob "**/upstream-parity*.md" 0 matches
✱ Glob ".mule/**/*.md" 0 matches
→ Read .mule
✗ Read .mule/research failed
Error: File not found: /root/astrate-mule/.mule/research
→ Read .mule/hosts
→ Read .mule/for-giulio.md
→ Read .mule/todo.md
The for-giulio.md history shows milestone runs through 2026-09-13. Today is 2026-09-15, so this run must re-verify. Let me check the repo layout and current state.

$ ls /root/ 2>&1; echo "---"; ls /root/astrate-mule 2>&1 | head -30; echo "---"; ls /root/astrate 2>&1 | head -30
Desktop
SurvivalManual
astrate
astrate-mule
astrate-survey
g-mind
go
govd
libheif
listello-buildtest
listello-mule
milestones-recipe.log
mule-crontab.bak
mule-crontab.paused
mule-failed-2026-09-04.tar.gz
mule-gate.log
pi.tail4db093.ts.net.crt
pi.tail4db093.ts.net.key
---
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
---
Dockerfile
LICENSE
Makefile
README.md
bench
cmd
deploy
docker-compose.yml
docs
go.mod
go.sum
internal
migrations
pkg
test
tools

$ gh issue list --label milestone-3.0 --state all --limit 50 --json number,title,state 2>&1; echo "exit=$?"
[]
exit=0

! permission requested: external_directory (/root/astrate/*); auto-rejecting
✗ ls /root/astrate/docs/UPSTREAM-EXPERIMENTAL.md /root/astrate-mule/docs/UPSTREAM-EXPERIMENTAL.md 2>&1; echo "---"; cd /root/astrate && git log --oneline -1 2>&1; cd /root/astrate-mule && git log --oneline -1 2>&1 failed
Error: The user rejected permission to use this specific tool call.

$ jq -r '.[] | [.number, .state, (.labels|map(.name)|join(","))] | @tsv' /tmp/opencode/upstream-parity-issues.json 2>/dev/null | sort -n | awk -F'\t' '{printf "%s\t%s\t%s\n", $1, $2, $3}'
47	CLOSED	enhancement,upstream-parity
48	CLOSED	enhancement,upstream-parity
49	CLOSED	enhancement,upstream-parity
50	CLOSED	enhancement,upstream-parity
51	CLOSED	enhancement,upstream-parity
52	CLOSED	enhancement,upstream-parity
53	CLOSED	enhancement,upstream-parity
54	CLOSED	enhancement,upstream-parity
55	CLOSED	enhancement,upstream-parity
56	CLOSED	enhancement,upstream-parity
57	CLOSED	enhancement,upstream-parity
58	CLOSED	enhancement,upstream-parity
59	CLOSED	enhancement,upstream-parity
60	CLOSED	enhancement,upstream-parity
61	CLOSED	enhancement,upstream-parity
62	CLOSED	enhancement,upstream-parity
63	CLOSED	enhancement,upstream-parity
64	CLOSED	enhancement,upstream-parity
65	CLOSED	enhancement,upstream-parity
66	CLOSED	enhancement,upstream-parity
67	CLOSED	enhancement,upstream-parity,upstream-experimental
68	CLOSED	enhancement,mule-blocked,upstream-parity,upstream-experimental
69	CLOSED	enhancement,upstream-parity
70	CLOSED	enhancement,upstream-parity
71	CLOSED	enhancement,upstream-parity
72	CLOSED	enhancement,upstream-parity
73	CLOSED	enhancement,upstream-parity
74	CLOSED	enhancement,upstream-parity
75	CLOSED	enhancement,upstream-parity
76	CLOSED	enhancement,upstream-parity
77	CLOSED	enhancement,upstream-parity
78	OPEN	enhancement,milestone-4.0,upstream-parity
79	CLOSED	enhancement,upstream-parity
80	CLOSED	enhancement,upstream-parity
81	CLOSED	enhancement,upstream-parity
82	CLOSED	enhancement,upstream-parity
83	CLOSED	enhancement,upstream-parity
84	CLOSED	enhancement,upstream-parity
85	CLOSED	enhancement,upstream-parity
86	CLOSED	enhancement,upstream-parity
87	CLOSED	enhancement,upstream-parity
88	CLOSED	enhancement,upstream-parity
89	CLOSED	enhancement,upstream-parity
91	CLOSED	mule-blocked,upstream-parity
92	OPEN	enhancement,upstream-parity,upstream-experimental
93	OPEN	mule-review,upstream-parity

$ gh api repos/astarte-platform/astarte/releases --paginate --jq '.[] | select(.prerelease==false) | [.tag_name, .published_at] | @tsv' 2>&1 | sort -V | tail -5; echo "---prereleases---"; gh api repos/astarte-platform/astarte/releases --paginate --jq '.[] | select(.prerelease==true) | [.tag_name, .published_at] | @tsv' 2>&1 | sort -V | tail -5
v1.2.2	2026-04-27T16:39:28Z
v1.3.0	2026-05-06T15:30:50Z
v1.3.1	2026-07-13T11:50:05Z
v1.3.2	2026-07-14T10:07:26Z
v1.3.3	2026-08-07T14:33:58Z
---prereleases---
v1.4.0-rc.1	2026-07-14T08:53:27Z
v1.4.0-rc.2	2026-07-14T10:16:22Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.5	2026-08-20T12:28:04Z

