slug: recipe-milestones
verdict: proposed
at:  9195e7f
ran: 2026-09-20T20:30:30Z on DietPi in 193s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
The first non-DONE milestone is **v3.0** (upstream Astarte parity). Its status names the first job explicitly: triage #47–#89 into an ordered plan, file sub-issues where work splits, escalate the "decide" set. Let me gather the reference and current state.

$ wc -l .mule/research/upstream-parity-2026-08-22.md docs/UPSTREAM-EXPERIMENTAL.md .mule/todo.md .mule/for-giulio.md 2>/dev/null; ls .mule/tasks/ 2>/dev/null
25 docs/UPSTREAM-EXPERIMENTAL.md
   202 .mule/todo.md
   324 .mule/for-giulio.md
   551 total
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md

$ gh issue list --label milestone-v3.0 --state all --limit 100
(no output)

$ gh issue list --label upstream-parity --state all --limit 100
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
91	CLOSED	Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health	mule-blocked, upstream-parity	2026-09-04T18:34:38Z
89	CLOSED	Dashboard flow-block schema mismatch (split_map/virtual pools hardcoded; null_sink/log_sink unknown)	enhancement, upstream-parity	2026-08-23T15:08:58Z
88	CLOSED	Flow auth: support a_f JWT claim	enhancement, upstream-parity	2026-08-23T15:08:59Z
87	CLOSED	Flow block: lua_map — needs embedded Lua runtime (parked)	enhancement, upstream-parity	2026-09-04T19:20:52Z
86	CLOSED	Flow: pipeline source DSL — keep DAG-JSON as documented deviation?	enhancement, upstream-parity	2026-08-23T11:57:02Z
85	CLOSED	Flow API: user-defined composite blocks	enhancement, upstream-parity	2026-08-23T21:47:32Z
84	CLOSED	Flow blocks: virtual_device_pool / dynamic_virtual_device_pool	enhancement, upstream-parity	2026-08-25T14:56:53Z
83	CLOSED	Flow blocks: mqtt_source/mqtt_sink (+modbus_tcp_source?) demand-driven	enhancement, upstream-parity	2026-08-23T15:09:03Z
82	CLOSED	Flow blocks: http_source/http_sink (demand-driven)	enhancement, upstream-parity	2026-08-23T15:09:04Z
81	CLOSED	Flow block: json_path_map	enhancement, upstream-parity	2026-08-23T15:09:06Z
80	CLOSED	Flow blocks: pure-transform set (to_json, update_metadata, split_map, random_source, sort)	enhancement, upstream-parity	2026-08-23T15:09:07Z
79	CLOSED	Verify registration-limit-reached HTTP status vs upstream	enhancement, upstream-parity	2026-08-24T13:21:57Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
77	CLOSED	Verify per-service version endpoints (GET /version, GET /v1/{realm}/version) served everywhere	enhancement, upstream-parity	2026-08-24T13:23:44Z
76	CLOSED	Housekeeping: GET /v1/realm-defaults/replication — decide reject/deviate (Cassandra-shaped)	enhancement, upstream-parity	2026-08-22T19:12:29Z
75	CLOSED	Housekeeping: decide realm-deletion gating/preconditions vs always-sync deviation	enhancement, upstream-parity	2026-08-23T15:09:10Z
74	CLOSED	Housekeeping: PATCH /v1/realms/{realm} (jwt key, registration limit, retention; null=unset)	enhancement, upstream-parity	2026-08-23T15:09:12Z
73	CLOSED	Housekeeping: default datastream retention injection env var (upstream 1.4)	enhancement, upstream-parity	2026-08-23T15:09:13Z
72	CLOSED	Realms: datastream_maximum_storage_retention ceiling (create/patch/enforce)	enhancement, upstream-parity	2026-08-23T15:09:15Z
71	CLOSED	Pairing: realm-scoped health check GET /v1/{realm}/health (upstream 1.3)	enhancement, upstream-parity	2026-08-23T15:09:16Z
70	CLOSED	Triggers: audit wildcard semantics (interface_name '*', match_path '/*' forcing rules)	enhancement, upstream-parity	2026-08-23T15:09:18Z
69	CLOSED	Verify unknown-realm HTTP status on RM endpoints against upstream	enhancement, upstream-parity	2026-08-24T13:23:42Z
68	CLOSED	Decide async_operation=false params vs documented always-sync deviation	enhancement, mule-blocked, upstream-parity, upstream-experimental	2026-09-04T18:34:35Z
67	CLOSED	Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4)	enhancement, upstream-parity, upstream-experimental	2026-09-04T18:34:19Z
66	CLOSED	Realm Management: detailed=true interface listing with full mappings (upstream 1.4)	enhancement, upstream-parity	2026-08-23T15:09:20Z
65	CLOSED	Policies: handler-overlap rejection + retry_times coupling + prefetch_count	enhancement, upstream-parity	2026-08-23T15:09:22Z
64	CLOSED	Triggers: decide AMQP action behavior (validate-reject vs NATS-forward deviation)	enhancement, upstream-parity	2026-08-23T15:09:23Z
63	CLOSED	Triggers: HTTP action validation limits (URL/method/header blocklist/template size)	enhancement, upstream-parity	2026-08-23T15:09:25Z
62	CLOSED	Realm Management: audit install/update/delete error codes and statuses	enhancement, upstream-parity	2026-08-23T15:09:26Z
61	CLOSED	Realm Management: audit interface/mapping validation matrix against astarte_core	enhancement, upstream-parity	2026-08-23T15:09:28Z
60	CLOSED	Realm Management: GET config/datastream_maximum_storage_retention (since upstream 1.2.0)	enhancement, upstream-parity	2026-08-22T03:23:47Z
59	CLOSED	AppEngine: group create-body validation + UUID-v1 from_token for group device listing	enhancement, upstream-parity	2026-08-22T10:35:18Z
58	CLOSED	AppEngine: PATCH requires Content-Type application/merge-patch+json	enhancement, upstream-parity	2026-08-22T03:15:59Z
57	CLOSED	AppEngine: audit server-write error taxonomy against upstream	enhancement, upstream-parity	2026-08-24T13:24:06Z
56	CLOSED	AppEngine: object-aggregate response formats table and disjoint_tables	enhancement, upstream-parity	2026-08-23T15:09:31Z
55	CLOSED	AppEngine: complete data-query parameter set	enhancement, upstream-parity	2026-08-23T15:09:33Z
54	CLOSED	AppEngine: GET /groups/{g} show-group endpoint	enhancement, upstream-parity	2026-08-22T10:35:16Z
53	CLOSED	AppEngine: group-scoped device PATCH + interface data access under groups/{g}/devices/{d}	enhancement, upstream-parity	2026-08-22T09:18:23Z
52	CLOSED	AppEngine: devices-by-alias full mirror (PATCH via alias + interface data access)	enhancement, upstream-parity	2026-08-22T09:18:20Z
51	CLOSED	keyAgreement control topic + encrypted mappings (upstream 1.4 experimental) — investigate/park	enhancement, upstream-parity	2026-08-22T08:48:28Z
50	CLOSED	Broker/engine hardening: discard malformed introspection/purge/capabilities instead of crashing	enhancement, upstream-parity	2026-08-22T04:04:24Z
49	CLOSED	MQTT v1: binaryblob properties re-send as BSON subtype 0 in consumer/properties	enhancement, upstream-parity	2026-08-22T04:04:23Z
48	CLOSED	MQTT v1: tolerate devices connecting with an empty introspection	enhancement, upstream-parity	2026-08-22T04:04:11Z
47	CLOSED	MQTT v1: accept device capabilities incl. purge_properties_compression_format (plaintext purge)	enhancement, upstream-parity	2026-08-22T04:04:09Z

✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
→ Read docs/UPSTREAM-EXPERIMENTAL.md
Most of #47–#89 are already closed. Let me read the reference research doc, the standing UPSTREAM-EXPERIMENTAL item, and check the remaining labels/issues.

$ gh issue list --label milestone-4.0 --state all --limit 50; echo "---90---"; gh issue view 90 --json number,title,state,labels,body
90	CLOSED	milestone-4.0: decide on bumping APICompatVersion from "1.2.2" to "1.3.0-rc.0" (or making it configurable)	enhancement, milestone-4.0	2026-08-23T21:06:34Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
28	CLOSED	milestone-4.0: investigate Edgehog compatibility with Astrate (readonly)	mule, readonly, milestone-4.0	2026-08-23T15:48:58Z
---90---
{"body":"Reminder of a pending decision parked during the upstream-parity investigation (#47–#89). **Do not raise the version now** — this is a milestone-4.0 item.\n\n## Context\n\n`realm.APICompatVersion` (`internal/realm/service.go`) is currently hardcoded to `\"1.2.2\"` — the emulated upstream API level Astrate declares.\n\nEdgehog 0.13.1 gates just-in-time provisioning triggers (`device_registered`, `device_deletion_started`) behind a Realm Management version >= 1.3. Astrate's engine already supports those trigger kinds end-to-end — verified in #28 (Edgehog compatibility investigation) — but Edgehog refuses to use them because of the declared version string.\n\n## What this issue must cover when picked up\n\n1. **Audit first**: what else does declaring \"1.3\" imply for other consumers?\n   - `astarte-dashboard` is currently validated against 1.2.2 behavior (see #13, error-name bundle comparison). Check whether any dashboard code path feature-gates on RM version besides the policies UI gate (>= 1.1.1).\n   - Sweep the parity findings (#47–#89): several are \"1.3/1.4 features\" we deliberately do NOT implement yet — does declaring 1.3 overpromise any of them?\n   - Check device SDKs and `astartectl` for version-gated behavior.\n2. **Then decide one of**:\n   - bump the constant to `\"1.3.0-rc.0\"` (upstream's own pre-release numbering for that line);\n   - or make it configuration (`internal/config`) so operators can declare the level their consumers expect.\n3. Update `docs/COMPATIBILITY.md` deviation #10 wording accordingly (propose, don't edit directly).\n\n## Outcome if done today (why it's parked)\n\nBumping without the audit risks silently promising surfaces we haven't verified; not bumping keeps Edgehog's JIT provisioning disabled even though the engine supports it. Milestone-4.0 scope decision.","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACtvELrg","name":"milestone-4.0","description":"","color":"5319e7"}],"number":90,"state":"CLOSED","title":"milestone-4.0: decide on bumping APICompatVersion from \"1.2.2\" to \"1.3.0-rc.0\" (or making it configurable)"}

The research doc isn't at that path. Let me find it.
✱ Glob ".mule/research/**/*" 0 matches

$ gh issue view 92 --json number,title,state,labels,body,comments 2>/dev/null | head -100
{"body":"Issue #51 was closed on 2026-08-22 with: *\"Full keyAgreement/encrypted-mappings support stays parked until the upstream 1.4 experimental spec stabilizes — reopen or file fresh when it does.\"* Filing fresh: that condition has now partly fired.\n\n**What changed.** Upstream commit `d084308` (\"docs: add encrypted endpoints key agreement (#2067)\", 2026-08-31) added `doc/pages/architecture/082-key_agreement_protocol.md` — a 267-line full wire protocol — and rewrote the Key Agreement section of `080-mqtt-v1-protocol.md`.\n\n**Verified against upstream `master` today, not taken from the survey:**\n\n- `082-key_agreement_protocol.md` exists and is 13,832 bytes.\n- The sentences the old parking rested on — \"reserved and routed correctly, but the handshake protocol is not yet implemented\" and \"accepted and acknowledged by Astarte, but no response is sent until the feature is fully implemented\" — are **gone from both 080 and 082** (grep count 0 in each).\n- 082 now states the handshake \"is implemented on the Astarte side (`astarte_data_updater_plant`)\", with only device-side session-key usage still described as liable to evolve.\n\n**What the spec now pins down**, where before there was only a sketch: five QoS-2 topics `control/keyAgreement/0..4` (InitExchange, ExchangeResp, SecretHash, HashOk, ExchangeFailed — the 080 table moved from a single `control/keyAgreement` at QoS 1); CBOR message bodies with CDDL given; `alg` 0 = `ECDH_P256-HKDF_SHA256-AES_256_GCM`; a CBOR-wrapped COSE_Key pubkey plus a 32-byte HkdfSalt; session-scoped keys never reused after reconnect or clean reset; AES-256-GCM whole-document ciphertext; and on \"no key established yet\" Astarte logs and **discards** rather than disconnecting. `ExchangeFailed` error codes are enumerated (0 internal, 1 invalid argument, 2 hash mismatch, 3 unprocessable entity).\n\n**The open question is scope, not feasibility, so this is not a mule task.** The spec is published and implemented upstream, but still ships only in `v1.4.0-rc.5` — `v1.3.3` remains the newest stable tag, and Astrate targets 1.2.2. So \"stabilized\" is true of the document and false of the release. Three ways to go:\n\n1. Implement now against the rc spec — largest surface in the parity backlog (CBOR codec, X25519/P-256, HKDF, AES-256-GCM, a 5-state handshake machine, shared-secret persistence, five new error names), and it would be built against a spec that may still move.\n2. Keep parking until `v1.4.0` is a stable tag, and re-file then.\n3. Do the narrow piece only — see the separate mule issue for the stale ACL comment, which is worth fixing regardless of which of these is chosen.\n\nRecorded in `.mule/for-giulio.md`. Survey source: `survey-2026-08-31.md` on `mule/research`.","comments":[{"id":"IC_kwDORmfsJs8AAAABSo0f8A","author":{"login":"atsetilam"},"authorAssociation":"MEMBER","body":"Parked by Giulio on 2026-09-04: **wait for a stable \\`v1.4.0\\`.** The spec is complete but ships only in \\`v1.4.0-rc.5\\`, and building against an rc risks rebuilding it.\n\nSo this issue is not waiting on anyone to remember it. It has a row in [`.mule/waiting-on.md`](https://github.com/astrate-platform/astrate/blob/main/.mule/waiting-on.md), and the weekly `mule-upstream-watch` job now reads that file as its first step every run: when a stable `v1.4.0` tag exists upstream, it escalates to `.mule/for-giulio.md` naming this issue. An `-rc.N` does not satisfy the row.\n\nThat mechanism was added because it was missing here: #51 was parked on the same condition, the condition fired on 2026-08-31, and nothing connected the new upstream tags to the parked issue — it took a hand triage four days later to catch it.","createdAt":"2026-09-04T19:52:21Z","includesCreatedEdit":false,"isMinimized":false,"minimizedReason":"","reactionGroups":[],"url":"https://github.com/astrate-platform/astrate/issues/92#issuecomment-5545730032","viewerDidAuthor":true}],"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"},{"id":"LA_kwDORmfsJs8AAAACxx3GJg","name":"upstream-experimental","description":"Adopted from an upstream experimental/unstable feature — review each upstream release","color":"d4c5f9"}],"number":92,"state":"OPEN","title":"keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds"}

#47–#89 are nearly all closed; #92 and #78 are parked decisions, #93 is an open task. The milestone's named source doc is missing from this repo. Let me check the repo layout, the queue, and the escalation file.
$ git remote -v; echo "---branch---"; git branch --show-current; echo "---waiting-on---"; ls -la .mule/waiting-on.md 2>/dev/null; echo "---research dir---"; ls -la .mule/research/ 2>/dev/null; echo "---find parity docs---"; rg -l "upstream-parity" --glob '*.md' .mule/ docs/ 2>/dev/null
origin	git@github.com:astrate-platform/astrate.git (fetch)
origin	git@github.com:astrate-platform/astrate.git (push)
---branch---
mule/queue
---waiting-on---
---research dir---
---find parity docs---
.mule/todo.md
.mule/milestones.md
.mule/reviews/housekeeping-2026-09-15.md
.mule/reviews/observability-2026-09-16.md
.mule/for-giulio.md
.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260909T230411Z-recipe-milestones-proposed.md
.mule/reports/log/20260911T185303Z-recipe-milestones-proposed.md
.mule/reports/log/20260918T184115Z-recipe-hygiene-timeout.md
.mule/reports/log/20260911T191541Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260907T171928Z-recipe-milestones-proposed.md
.mule/reports/log/20260913T170550Z-recipe-github-issues-proposed.md
.mule/reports/log/20260908T112032Z-recipe-github-issues-proposed.md
.mule/reports/log/20260915T113026Z-recipe-code-review-proposed.md
.mule/reports/log/20260906T180442Z-recipe-milestones-proposed.md
.mule/reports/log/20260920T185956Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260918T193641Z-recipe-milestones-proposed.md
.mule/reports/log/20260907T204751Z-recipe-milestones-proposed.md
.mule/reports/log/20260913T203209Z-recipe-milestones-proposed.md
.mule/reports/log/20260910T184706Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260916T115527Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260915T173739Z-hk-zero-retention-asymmetry-done.md
.mule/reports/log/20260915T170704Z-hk-zero-retention-asymmetry-done.md
.mule/reports/log/20260916T202828Z-observability-health-content-type-test-done.md
.mule/reports/log/20260915T180720Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260911T194026Z-recipe-code-review-proposed.md
.mule/reports/log/20260917T111017Z-docs-sync-pairing-version-value-done.md
.mule/reports/log/20260904T195529Z-issue-93-done.md
.mule/reports/log/20260914T234550Z-recipe-milestones-proposed.md
.mule/reports/log/20260918T114100Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260917T184319Z-recipe-milestones-timeout.md
.mule/reports/log/20260915T185804Z-recipe-github-issues-proposed.md
.mule/reports/log/20260913T203830Z-recipe-github-issues-proposed.md
.mule/reports/log/20260905T181219Z-recipe-github-issues-proposed.md
.mule/reports/log/20260918T113442Z-recipe-github-issues-proposed.md
.mule/reports/log/20260913T203800Z-recipe-github-issues-proposed.md
.mule/reports/log/20260905T181053Z-recipe-milestones-proposed.md
.mule/reports/log/20260910T200120Z-recipe-github-issues-proposed.md
.mule/reports/log/20260908T193127Z-recipe-milestones-proposed.md
.mule/reports/log/20260918T110039Z-recipe-milestones-proposed.md
.mule/reports/log/20260911T190417Z-recipe-github-issues-proposed.md
.mule/reports/log/20260909T202852Z-docs-sync-native-compat-version-endpoints-done.md
.mule/reports/log/20260915T193626Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260916T183922Z-recipe-github-issues-proposed.md
.mule/reports/log/20260908T193321Z-recipe-github-issues-proposed.md
.mule/reports/log/20260914T114442Z-auth-iat-not-required-test-done.md
.mule/reports/log/20260917T190323Z-recipe-github-issues-proposed.md
.mule/reports/log/20260912T172901Z-recipe-github-issues-proposed.md
.mule/reports/log/20260910T110956Z-recipe-github-issues-proposed.md
.mule/reports/log/20260919T211455Z-recipe-milestones-proposed.md
.mule/reports/log/20260916T175801Z-recipe-milestones-proposed.md
.mule/reports/log/20260917T190703Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260920T111057Z-recipe-github-issues-proposed.md
.mule/reports/log/20260916T184218Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260910T170149Z-recipe-code-review-proposed.md
.mule/reports/log/20260909T175718Z-recipe-milestones-proposed.md
.mule/reports/log/20260915T201017Z-recipe-code-review-proposed.md
.mule/reports/log/20260916T184848Z-recipe-code-review-proposed.md
.mule/reports/log/20260912T171256Z-recipe-milestones-proposed.md
.mule/reports/log/20260910T192802Z-recipe-hygiene-proposed.md
.mule/reports/log/20260911T110803Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260918T194216Z-recipe-github-issues-proposed.md
.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260910T195304Z-recipe-milestones-proposed.md
.mule/reports/log/20260915T105740Z-recipe-github-issues-proposed.md
.mule/reports/log/20260908T193026Z-recipe-milestones-proposed.md
.mule/reports/log/20260906T181352Z-recipe-github-issues-proposed.md
.mule/reports/log/20260915T184406Z-recipe-milestones-proposed.md
.mule/reports/log/20260907T174812Z-recipe-github-issues-proposed.md

→ Read .mule/todo.md
The todo queue is long. Let me check the most recent milestones-proposal reports and locate the missing research doc (maybe on another branch or in history).

$ for f in 20260919T211455Z-recipe-milestones-proposed.md 20260918T193641Z-recipe-milestones-proposed.md 20260918T110039Z-recipe-milestones-proposed.md 20260917T184319Z-recipe-milestones-timeout.md 20260916T175801Z-recipe-milestones-proposed.md; do echo "=== $f ==="; cat ".mule/reports/log/$f"; echo; done
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_0c0821fe9001ECfmTHk5GZDelG

|---|---|---|---|
| **encrypted mapping attribute** *(new)* | Interfaces with `encrypted: true` mappings install and persist (post-#1999) | **Rejected at parse time** — strict decoding (`DisallowUnknownFields`, `parse.go:121`; no `Encrypted` in `mappingJSON`, `parse.go:99-113`) | RM schema `mapping.ex:113-119` |
| **key agreement** | Handshake implemented server-side (despite stale doc note) — see correction A | Topic ACL-denied (`aclhook.go:63-68`); zero key machinery | `control_handler.ex:171+`; `0da8336` |
| **encrypted endpoints (transport)** | DEK/KEK at-rest wrap + COSE Encrypt0 both directions | Absent (no DEK/KEK/COSE in `internal/`) | `maybe_encrypt_value`; `device.ex` |
| **emptyCache resend-failure triggers** | Distinct device_error names + session end | Retries/skips loudly; never triggers, never disconnects (`control.go:58-76`) | `errorname.go:20-21` |
| **malformed capabilities** | Force-disconnect + `malformed_capabilities_message` | Rejects, keeps connection, fallback error name (`capabilities.go:29`) | `broker.go:260` unused here |
| **`unset_not_allowed`** | Emits literal name | Mis-mapped to `unexpected_value_type` (`errorname.go:62`) | `data_handler.ex:255-260` |
| **new upstream error names** | `decryption_error`, `malformed_capabilities_message` | Absent; fall back to `interface_loading_failed` | `data_handler.ex:453`; `capabilities_handler.ex:78` |
| **BSON-null value** | `nil` individual passes validation | Rejected, test-pinned (`payload_test.go:499`) | upstream `4789d047` |
| **group-scoped triggers** | `group_name` restricts trigger scope | Compiled, never matched (`match.go:318,354`) | `060-triggers.md:121-123` |
| **unverified error names** | `write_on_server_owned_interface`, `value_size_exceeded` | Mapped but unmeasured (`COMPATIBILITY.md:122`) | issues #18/#19 |
| **group WATCH paths** | `groups/<name>/...` authorization shapes | Not measured (bench realm has no group) | issue #17 |

Net: **one new gap** (attribute acceptance at install time), everything else standing
re-verified; corrections A/B refine two rows' provenance without changing their status.

---

## Source 5 — upstream `atomvm/AtomVM`

Unchanged (`0220c78`, same since 2026-07-13). Conclusion stands: AtomVM matters to
Astrate only as a target platform for the JSON payload profile (`docs/DESIGN.md` §3.5);
no AtomVM-specific work item.

$ echo "=== milestone recipe runs in mule log ==="; rg -n -i 'milestone' .mule/log.md | head -40; echo "=== for-giulio mentions of milestone/3.0 ==="; rg -n -i 'milestone.?3|v3\.0|milestone' .mule/for-giulio.md | head -40
=== milestone recipe runs in mule log ===
=== for-giulio mentions of milestone/3.0 ===
13:- **github-issues triage run, 2026-09-16: still nothing proposable — the mule-alarm pile is now 11 straight days (#94–#104).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), plus the daily alarm. New since the 2026-09-15 run: **#104**, today's alarm (created 09-16 11:12Z, "nothing has landed in 14h"); **#103** has expired into the pile (superseded by #104). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-15 one: close #94–#103** — each was a one-day event, superseded by next day's alarm, never actionable; leave #104 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
15:- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-13 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
17:- **github-issues triage run, 2026-09-15: still nothing proposable — the mule-alarm pile is now 10 straight days (#94–#103).** Fourteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-13 run: **#103**, the alarm created 2026-09-14 11:27Z ("nothing has landed in 16h"); **#102** has expired into the pile (superseded by #103). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-13 one: close #94–#102** — each was a one-day event, superseded by next day's alarm, never actionable; leave #103 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
19:- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-12 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
21:- **github-issues triage run, 2026-09-13: still nothing proposable — the mule-alarm pile is now 9 straight days (#94–#102).** Thirteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-12 run: **#102**, today's alarm (created 11:09Z, "nothing has landed in 14h"); **#101** has expired into the pile (superseded by #102). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-12 one: close #94–#101** — each was a one-day event, superseded by next day's alarm, never actionable; leave #102 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
23:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
27:- **github-issues triage run, 2026-09-12: still nothing proposable — the mule-alarm pile is now 8 straight days (#94–#101).** Twelve open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-11 run: **#101**, today's alarm (created 10:58Z, "nothing has landed in 14h"); **#100** has expired into the pile (superseded by #101). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-11 one: close #94–#100** — each was a one-day event, superseded by next day's alarm, never actionable; leave #101 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
29:- **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pile is now 7 straight days and worth a look.** Twelve open issues. Still no machine-checkable fix candidates, so no task lines: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-10 run: **#100**, today's alarm (created 10:55Z, "nothing has landed in 15h"). #99–#94 are the same `mule-alarm` (one per day 2026-09-05→09-10, ~11:00Z each, zero comments) — telemetry, not code issues, so never proposable. But seven consecutive alarms is past the "low-activity window" wording earlier runs used: the mule has landed no commit since ~2026-09-04/05, and `main`'s queue copy shows the top line as `- [ ]` while the landed work below is all `[!] BLOCKED` (`wrote nothing` / `tests failed`) — though `mule/queue` is the authoritative copy. **Proposal: close #94–#99** — each was a one-day event, superseded by next day's alarm, never actionable; leave #100 (live today) to self-expire. If the idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
31:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
33:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
35:- **Extension point: external-bus Intake (NATS JetStream), 2026-09-10.** `internal/broker/intake.go:56-62` carries a `TODO(extension point)` for a second `Intake` implementation backed by a durable external bus (NATS JetStream) for multi-instance deployment or restart survival. The frozen design decision is in `docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4`; the consumer must reproduce per-device ordering, deferred-ack backpressure, and QoS 0 drop semantics. Not tracked by any open issue. Your call: file a milestone issue or park until the roadmap calls for it.
37:- **Extension point: timescaledb_toolkit lttb downsampling, 2026-09-10.** `internal/store/store.go:139-143` carries a `TODO(extension point)` for switching the `Downsample` method from `time_bucket+avg` to the toolkit's `lttb()` when `timescaledb_toolkit` is present — the probe at `store.go:144-150` already records availability in `s.hasToolkit`, but `datastreams.go` always uses the default path. Design reference: `docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5`. Not tracked by any open issue. Your call: file a milestone issue or park.
60:  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated. #1
63:- **Milestone recipe run, 2026-09-09: v3.0 unchanged since 2026-09-06 — still waiting on
64:  upstream gates.** `milestone-3.0` label remains empty. The three open `upstream-parity`
66:  aclhook comment (pushed commit `8c61268`, mule-review), #78 FDO (milestone-4.0). Both
77:  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated below.
84:- **Milestone recipe run, 2026-09-06: v3.0 (upstream 1.2.2 → 1.3/1.4) has no open
85:  machine-checkable gaps left to file; the milestone now waits on two upstream gates, and the
86:  `milestone-3.0` label is empty (this milestone's work lives under `upstream-parity`).** All
92:  (milestone-4.0, escalated), **#93** aclhook comment rewrite (mule-review, pushed commit
98:  only) — until then v3.0 has no open implementable work. This run filed **no issues** and
101:- **`milestones.md` v3.0 names a source doc that does not exist: `.mule/research/
103:  investigation's *output* (issues #47–#89) is intact and reachable, so milestone work is
105:  milestone paragraph re-pointed) before v4.0 drafting leans on it.
118:  row. #78 is the milestone-4.0 FDO design/investigation, already escalated below. #1
139:  2026-09-05 — a prior escalation from the 2026-09-04 milestone run was lost in the queue
144:- **govulncheck GO-2026-5970: reachable DoS in golang.org/x/text (infinite loop on invalid input, fixed in v0.39.0, available v0.41.0).** Astrate pins `x/text` indirect at v0.38.0 (go.mod:97) and pgx pulls it into production: `internal/store/notify.go:59` `store.Listen` → `pgx.ConnectConfig` → `unicode/norm.*`. This is the only govulncheck symbol finding that is not test-harness-only: GO-2026-6355/6354 (x/crypto/ssh deadlocked-channel DoS) and GO-2026-6253 (moby/go-archive tar path traversal) are reachable only through testcontainers in `internal/testutil/pg.go`, i.e. never in the deployed binary. `x/text` keeps API compatibility minor-to-minor and the modules Astrate exercises (`unicode/norm` via pgx, `text/language` via jsonschema) are unchanged, so this is a fix Astrate actually needs — the hygiene recipe's highest-priority category. Not a mule task (go.mod never-touch): your decision to bump ≥v0.39.0 now or fold into the next milestone-boundary sweep. Raw: https://pkg.go.dev/vuln/GO-2026-5970. (The 2026-09-04 dep sweep did not list x/text.)
150:plaintext `purge_properties_compression_format`, empty-introspection allowance, device registration/deletion triggers, experimental FDO pairing auth) are not yet emulated and are out of scope until the milestone that adopts v1.3.2 as the target."* — and reword deviation 18's realm-health note from "which upstream 404s" to "added by upstream v1.3 (Astrate serves it against a 1.2.2 target; kept, matching behavior)". Raw upstream changes: [v1.3.0](https://github.com/astarte-platform/astarte/releases/tag/v1.3.0), [v1.3.2](https://github.com/astarte-platform/astarte/releases/tag/v1.3.2).
180:  milestone boundary (the point where `APICompatVersion` or a milestone tag moves), and bump
185:- **milestone 2.0 looks complete, verify and cut the tag** — all 11 `milestone-2.0` issues
189:  durable named flows all landed).   (Milestones recipe run, 2026-09-03.)
191:  **Decided 2026-09-04: the tag is `v0.2.0`, cut on the same day.** The `v2.0`/`v3.0` names
192:  are milestone names, not release versions — the project is still pre-1.0 and the version
193:  number keeps its own line rather than jumping two majors to match a milestone label.
213:- **#78 FDO device onboarding — milestone-4.0, investigation phase.** Too large for a single
215:  schema and keys) is a multi-session project. Parked until the v3.0 queue clears and this
216:  becomes the next milestone target.
276:## 2026-08-23 — FDO promoted to milestone-4.0 candidate (Giulio's decision, recorded)
283:When v3.0 is marked DONE, the v4.0 section of `.mule/milestones.md` should be

$ sed -n '107,145p' .mule/milestones.md
## v3.0 — upstream Astarte feature parity (1.2.2 → 1.3/1.4)

Scope decided by Giulio on 2026-08-22. The original v3.0 candidate (building the
remaining SECO suite pieces) was dropped for legal reasons, and Edgehog is already
fully compatible with Astrate (verified in #28) — so this milestone takes the body of
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

## v4.0+ — next

Placeholder. Scope undecided: candidates once v3.0 is DONE include the remaining
`upstream-parity` backlog items parked today, or whatever emerges from upstream
releases after v1.4 ships. The recipe should propose the next section's draft
(name + one-line scope) as a `for-giulio.md` entry once v3.0 is `DONE`, rather than
leaving this placeholder to rot.

Status: not started, scope undecided.

→ Read .mule/milestones.md [offset=105, limit=40]
Let me re-read the actual v3.0 section of the milestones file precisely, since it's the current target.

$ rg -n -A6 'upstream-parity-2026-08-22\.md.*names a source doc|names a source doc' .mule/for-giulio.md | head -30
101:- **`milestones.md` v3.0 names a source doc that does not exist: `.mule/research/
102-  upstream-parity-2026-08-22.md` is absent from `main` and from `origin/mule/research`.** The
103-  investigation's *output* (issues #47–#89) is intact and reachable, so milestone work is
104-  unaffected, but if that doc is your canonical 1.3/1.4 diff it needs restoring (or the
105-  milestone paragraph re-pointed) before v4.0 drafting leans on it.
106-
107-- **docs-sync realm-management: two ASTRATE_ config keys exist in code but are missing from

→ Read .mule/milestones.md [offset=105, limit=40]

$ echo "=== docs/UPSTREAM-EXPERIMENTAL.md head ==="; sed -n '1,50p' docs/UPSTREAM-EXPERIMENTAL.md 2>/dev/null || find . -iname 'UPSTREAM-EXPERIMENTAL.md' -not -path './.git/*'
=== docs/UPSTREAM-EXPERIMENTAL.md head ===
# Upstream experimental features register

Astarte's brain is the upstream Astarte team. When upstream ships a feature
marked experimental / unstable / rc-only, Astrate adopts it early (so the wire
behaviour is ready) but tracks it here until upstream promotes or deprecates
it. **Review trigger: every time a future parity milestone targets a newer
upstream version, check each row below against what that release did with the
feature** — promoted → keep and drop the row; deprecated → deprecate in the
same milestone, in sync.

## The version-declaration rule

`APICompatVersion` (internal/realm/service.go) must never claim a level higher
than what Astrate actually implements, and must never be operator-configurable.
It is bumped only in the same change that completes the full surface of that
upstream level. When that happens, every open row below tagged with that level
must be reconciled first — the register is the checklist, so nothing has to be
remembered.

## Rows

| Feature | Issue | Upstream level | Adopted | Reconcile when |
|---|---|---|---|---|
| Required + encrypted interface mapping fields | #67 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
| `async_operation=false` on mutating RM endpoints | #68 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |

=== 20260916T175801Z-recipe-milestones-proposed.md ===
slug: recipe-milestones
verdict: proposed
at:  70ca8df
ran: 2026-09-16T17:58:01Z on DietPi in 164s

> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
Milestone v3.0 is the first non-DONE section. Let me load the reference doc, the issue state, and the upstream recipe rules.

$ wc -l .mule/research/upstream-parity-2026-08-22.md docs/UPSTREAM-EXPERIMENTAL.md .mule/recipes/astarte-upstream.md .mule/todo.md .mule/for-giulio.md 2>/dev/null; ls .mule/research/
25 docs/UPSTREAM-EXPERIMENTAL.md
   62 .mule/recipes/astarte-upstream.md
  181 .mule/todo.md
  284 .mule/for-giulio.md
  552 total
ls: cannot access '.mule/research/': No such file or directory

$ gh issue list --label milestone-3.0 --state all --limit 50 2>&1; echo '---'; gh issue list --label upstream-parity --state all --limit 60 2>&1 | head -60
---
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
91	CLOSED	Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health	mule-blocked, upstream-parity	2026-09-04T18:34:38Z
89	CLOSED	Dashboard flow-block schema mismatch (split_map/virtual pools hardcoded; null_sink/log_sink unknown)	enhancement, upstream-parity	2026-08-23T15:08:58Z
88	CLOSED	Flow auth: support a_f JWT claim	enhancement, upstream-parity	2026-08-23T15:08:59Z
87	CLOSED	Flow block: lua_map — needs embedded Lua runtime (parked)	enhancement, upstream-parity	2026-09-04T19:20:52Z
86	CLOSED	Flow: pipeline source DSL — keep DAG-JSON as documented deviation?	enhancement, upstream-parity	2026-08-23T11:57:02Z
85	CLOSED	Flow API: user-defined composite blocks	enhancement, upstream-parity	2026-08-23T21:47:32Z
84	CLOSED	Flow blocks: virtual_device_pool / dynamic_virtual_device_pool	enhancement, upstream-parity	2026-08-25T14:56:53Z
83	CLOSED	Flow blocks: mqtt_source/mqtt_sink (+modbus_tcp_source?) demand-driven	enhancement, upstream-parity	2026-08-23T15:09:03Z
82	CLOSED	Flow blocks: http_source/http_sink (demand-driven)	enhancement, upstream-parity	2026-08-23T15:09:04Z
81	CLOSED	Flow block: json_path_map	enhancement, upstream-parity	2026-08-23T15:09:06Z
80	CLOSED	Flow blocks: pure-transform set (to_json, update_metadata, split_map, random_source, sort)	enhancement, upstream-parity	2026-08-23T15:09:07Z
79	CLOSED	Verify registration-limit-reached HTTP status vs upstream	enhancement, upstream-parity	2026-08-24T13:21:57Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
77	CLOSED	Verify per-service version endpoints (GET /version, GET /v1/{realm}/version) served everywhere	enhancement, upstream-parity	2026-08-24T13:23:44Z
76	CLOSED	Housekeeping: GET /v1/realm-defaults/replication — decide reject/deviate (Cassandra-shaped)	enhancement, upstream-parity	2026-08-22T19:12:29Z
75	CLOSED	Housekeeping: decide realm-deletion gating/preconditions vs always-sync deviation	enhancement, upstream-parity	2026-08-23T15:09:10Z
74	CLOSED	Housekeeping: PATCH /v1/realms/{realm} (jwt key, registration limit, retention; null=unset)	enhancement, upstream-parity	2026-08-23T15:09:12Z
73	CLOSED	Housekeeping: default datastream retention injection env var (upstream 1.4)	enhancement, upstream-parity	2026-08-23T15:09:13Z
72	CLOSED	Realms: datastream_maximum_storage_retention ceiling (create/patch/enforce)	enhancement, upstream-parity	2026-08-23T15:09:15Z
71	CLOSED	Pairing: realm-scoped health check GET /v1/{realm}/health (upstream 1.3)	enhancement, upstream-parity	2026-08-23T15:09:16Z
70	CLOSED	Triggers: audit wildcard semantics (interface_name '*', match_path '/*' forcing rules)	enhancement, upstream-parity	2026-08-23T15:09:18Z
69	CLOSED	Verify unknown-realm HTTP status on RM endpoints against upstream	enhancement, upstream-parity	2026-08-24T13:23:42Z
68	CLOSED	Decide async_operation=false params vs documented always-sync deviation	enhancement, mule-blocked, upstream-parity, upstream-experimental	2026-09-04T18:34:35Z
67	CLOSED	Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4)	enhancement, upstream-parity, upstream-experimental	2026-09-04T18:34:19Z
66	CLOSED	Realm Management: detailed=true interface listing with full mappings (upstream 1.4)	enhancement, upstream-parity	2026-08-23T15:09:20Z
65	CLOSED	Policies: handler-overlap rejection + retry_times coupling + prefetch_count	enhancement, upstream-parity	2026-08-23T15:09:22Z
64	CLOSED	Triggers: decide AMQP action behavior (validate-reject vs NATS-forward deviation)	enhancement, upstream-parity	2026-08-23T15:09:23Z
63	CLOSED	Triggers: HTTP action validation limits (URL/method/header blocklist/template size)	enhancement, upstream-parity	2026-08-23T15:09:25Z
62	CLOSED	Realm Management: audit install/update/delete error codes and statuses	enhancement, upstream-parity	2026-08-23T15:09:26Z
61	CLOSED	Realm Management: audit interface/mapping validation matrix against astarte_core	enhancement, upstream-parity	2026-08-23T15:09:28Z
60	CLOSED	Realm Management: GET config/datastream_maximum_storage_retention (since upstream 1.2.0)	enhancement, upstream-parity	2026-08-22T03:23:47Z
59	CLOSED	AppEngine: group create-body validation + UUID-v1 from_token for group device listing	enhancement, upstream-parity	2026-08-22T10:35:18Z
58	CLOSED	AppEngine: PATCH requires Content-Type application/merge-patch+json	enhancement, upstream-parity	2026-08-22T03:15:59Z
57	CLOSED	AppEngine: audit server-write error taxonomy against upstream	enhancement, upstream-parity	2026-08-24T13:24:06Z
56	CLOSED	AppEngine: object-aggregate response formats table and disjoint_tables	enhancement, upstream-parity	2026-08-23T15:09:31Z
55	CLOSED	AppEngine: complete data-query parameter set	enhancement, upstream-parity	2026-08-23T15:09:33Z
54	CLOSED	AppEngine: GET /groups/{g} show-group endpoint	enhancement, upstream-parity	2026-08-22T10:35:16Z
53	CLOSED	AppEngine: group-scoped device PATCH + interface data access under groups/{g}/devices/{d}	enhancement, upstream-parity	2026-08-22T09:18:23Z
52	CLOSED	AppEngine: devices-by-alias full mirror (PATCH via alias + interface data access)	enhancement, upstream-parity	2026-08-22T09:18:20Z
51	CLOSED	keyAgreement control topic + encrypted mappings (upstream 1.4 experimental) — investigate/park	enhancement, upstream-parity	2026-08-22T08:48:28Z
50	CLOSED	Broker/engine hardening: discard malformed introspection/purge/capabilities instead of crashing	enhancement, upstream-parity	2026-08-22T04:04:24Z
49	CLOSED	MQTT v1: binaryblob properties re-send as BSON subtype 0 in consumer/properties	enhancement, upstream-parity	2026-08-22T04:04:23Z
48	CLOSED	MQTT v1: tolerate devices connecting with an empty introspection	enhancement, upstream-parity	2026-08-22T04:04:11Z
47	CLOSED	MQTT v1: accept device capabilities incl. purge_properties_compression_format (plaintext purge)	enhancement, upstream-parity	2026-08-22T04:04:09Z

$ ls /root/ ; echo '---'; git -C /root/astrate-mule remote -v 2>/dev/null
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
origin	git@github.com:astrate-platform/astrate.git (fetch)
origin	git@github.com:astrate-platform/astrate.git (push)

$ gh issue list --label upstream-parity --state all --limit 60 2>&1 | tail -5; echo '--- total ---'; gh issue list --label upstream-parity --state all --limit 100 2>&1 | wc -l
51	CLOSED	keyAgreement control topic + encrypted mappings (upstream 1.4 experimental) — investigate/park	enhancement, upstream-parity	2026-08-22T08:48:28Z
50	CLOSED	Broker/engine hardening: discard malformed introspection/purge/capabilities instead of crashing	enhancement, upstream-parity	2026-08-22T04:04:24Z
49	CLOSED	MQTT v1: binaryblob properties re-send as BSON subtype 0 in consumer/properties	enhancement, upstream-parity	2026-08-22T04:04:23Z
48	CLOSED	MQTT v1: tolerate devices connecting with an empty introspection	enhancement, upstream-parity	2026-08-22T04:04:11Z
47	CLOSED	MQTT v1: accept device capabilities incl. purge_properties_compression_format (plaintext purge)	enhancement, upstream-parity	2026-08-22T04:04:09Z
--- total ---
46

$ git -C /root/astrate branch -a 2>/dev/null | head -20; echo '---'; ls /root/astrate/.mule/research/ 2>/dev/null; echo '---'; git -C /root/astrate rev-parse --abbrev-ref HEAD 2>/dev/null
* main
  wip/DietPi-20260727T171259Z
  wip/DietPi-20260727T171543Z
  remotes/origin/HEAD -> origin/main
  remotes/origin/feat/pokemon-agent
  remotes/origin/fix/blocked-mule-issues
  remotes/origin/main
  remotes/origin/mule/queue
  remotes/origin/mule/research
  remotes/origin/wip/DietPi-20260727T171543Z
  remotes/origin/worktree-m12-05-match-upstream
  remotes/origin/worktree-m12-06-ach-join
  remotes/origin/worktree-m12-06b-watch-path
---
---
main

$ gh issue list --label milestone-4.0 --state all --limit 20 2>&1; echo '--- open milestone-3x none? ---'; gh issue list --label milestone-3.0 --state all --limit 20 2>&1
90	CLOSED	milestone-4.0: decide on bumping APICompatVersion from "1.2.2" to "1.3.0-rc.0" (or making it configurable)	enhancement, milestone-4.0	2026-08-23T21:06:34Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
28	CLOSED	milestone-4.0: investigate Edgehog compatibility with Astrate (readonly)	mule, readonly, milestone-4.0	2026-08-23T15:48:58Z
--- open milestone-3x none? ---

$ git -C /root/astrate ls-tree --name-only origin/mule/research .mule/research/ 2>/dev/null; echo '--- worktree research fields ---'; git -C /root/astrate show origin/mule/research:.mule/research/ 2>/dev/null | head; echo '--- display name ---'; git -C /root/astrate show origin/mule/research --stat 2>/dev/null | head -5
.mule/research/README.md
.mule/research/issues-2026-07-27.md
.mule/research/issues-2026-07-28.md
.mule/research/issues-2026-07-30.md
.mule/research/issues-2026-08-01.md
.mule/research/issues-2026-08-06.md
.mule/research/issues-2026-08-07.md
.mule/research/issues-2026-08-08.md
.mule/research/issues-2026-08-13.md
.mule/research/issues-2026-08-15.md
.mule/research/issues-2026-08-18.md
.mule/research/issues-2026-08-21.md
.mule/research/issues-2026-08-22.md
.mule/research/issues-2026-08-23.md
.mule/research/issues-2026-08-24.md
.mule/research/issues-2026-08-25.md
.mule/research/issues-2026-08-26.md
.mule/research/log.md
.mule/research/survey-2026-07-27.md
.mule/research/survey-2026-07-28.md
.mule/research/survey-2026-07-30.md
.mule/research/survey-2026-08-01.md
.mule/research/survey-2026-08-06.md
.mule/research/survey-2026-08-07.md
.mule/research/survey-2026-08-08.md
.mule/research/survey-2026-08-13.md
.mule/research/survey-2026-08-15.md
.mule/research/survey-2026-08-18.md
.mule/research/survey-2026-08-21.md
.mule/research/survey-2026-08-22.md
.mule/research/survey-2026-08-23.md
.mule/research/survey-2026-08-24.md
.mule/research/survey-2026-08-25.md
.mule/research/survey-2026-08-26.md
.mule/research/triaged.md
--- worktree research fields ---
tree origin/mule/research:.mule/research/

README.md
issues-2026-07-27.md
issues-2026-07-28.md
issues-2026-07-30.md
issues-2026-08-01.md
issues-2026-08-06.md
issues-2026-08-07.md
issues-2026-08-08.md
--- display name ---
commit ce759e1870a69257df5b2486cd679f9c455a4be1
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 31 01:48:20 2026 +0200

    survey: recover 16 daily reports the broken commit check left staged (2026-07-28..2026-08-26), and take the tooling fix from main

$ ls /root/astrate-survey 2>/dev/null | head; echo '---'; find /root -maxdepth 3 -name 'upstream-parity*' 2>/dev/null
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
---

$ git -C /root/astrate show origin/mule/research:.mule/research/survey-2026-08-22.md 2>/dev/null | head -50; echo '---END HEAD---'; git -C /root/astrate show origin/mule/research:.mule/research/issues-2026-08-22.md 2>/dev/null | head -50
# Survey — cross-project gap analysis (2026-08-22)

Sources: Astrate codebase/docs (`main` moved `89145e6` → `4c881d0f`, read from a fresh
shallow https clone — this checkout's working tree remains stale and was **not** used),
mule queue history (`origin/mule/queue` `4f6d95c` → `28f3adf5`), upstream
`astarte-platform` (`master` moved `739724991` → `372a2cbe`; **no new tags** — latest
still `v1.3.3` and `v1.4.0-rc.5`), upstream `atomvm/AtomVM` (unchanged). Every claim
cites a specific path, commit, or doc section.

Incremental note: last survey was `survey-2026-08-21.md`. Since then:

- **Astrate's `main` moved for the first time since 07-29** — `89145e6` → `4c881d0f`,
  5 commits. All are Flow-v2.0 follow-ups, lint, and milestone bookkeeping; **no
  broker/engine/payload behaviour change** (details in Source 1). Two citation-relevant
  side effects: `capabilities.go` shifted +2 lines (import regrouping only), and the
  range exposed another stale claim in recent surveys (correction B below).
- **Upstream `astarte-platform` master moved** — one commit,
  `372a2cb feat: add encrypted endpoints (#1999)`. It is **schema/plumbing only**: it
  makes interfaces *with* encrypted endpoints installable by adding an optional boolean
  `encrypted` mapping attribute. It lands no crypto (analysis in Source 3). This
  creates **one new wire-visible gap for Astrate** (Source 4, new candidate).
- **`origin/mule/queue` moved** — `4f6d95c` → `28f3adf5`: three recipe-run bookkeeping
  entries plus `28f3adf` ("mule: fix ticks dying outside the repo (cron cwd), recover
  unsaved queue state") — mule infrastructure, no Astrate code.
- **AtomVM unchanged** — HEAD `0220c78` (2026-07-13), byte-identical to every run since;
  conclusion stands (Source 5).

## Corrections to prior surveys

### A. Upstream doc/code contradiction on keyAgreement (affects the ACL-deny rationale)

Astrate's deliberate ACL denial of `control/keyAgreement`
(`internal/broker/aclhook.go:63-68`) justifies itself by citing upstream's protocol doc:
*"upstream itself documents the handshake as 'reserved and routed correctly, but not yet
implemented'"*. That is still literally true of the **doc**
(`astarte-platform` `doc/pages/architecture/080-mqtt-v1-protocol.md:172-176`, unchanged)
but has been false of the **code** on master since before 07-29:
`data_updater_plant/.../control_handler.ex:171+` decodes and state-machines
`keyAgreement/0`–`/4` (InitExchange/ExchangeResp/SecretHash/HashOk/ExchangeFailed,
sending `exchange_failed` on transition errors), and `0da8336` ("chore(dup): store
shared secret during key agreement", #2101) persists the shared secret. Upstream's own
docs lag their implementation. The deny remains safe for Astrate today (it has zero key
machinery either way), but the comment's *citation* is misleading and
`key-agreement-protocol`'s body should point at the code, not the doc.

### B. "External-bus intake TODO at store.go:135" is a phantom — drop it

The 08-21 survey (Source 1 and the issues file's already-filed note) carried:
*"External-bus intake TODO still open — internal/store/store.go:135 (blocked on GitHub
issue #10)"*. On `origin/main` `4c881d0f`:
---END HEAD---
# Candidate issues — 2026-08-22

Distilled from `survey-2026-08-22.md` and carried forward from
`issues-2026-08-21.md`. Ordered highest-value-first, capped at 8. Every `mule-line`
and `mule-spec` candidate was re-grep-confirmed today as not already implemented
against `origin/main` `4c881d0f` (== https remote HEAD; moved from `89145e6` — the
range is flow/lint/milestones only, so all prior verifications carry over; line refs
updated where the lint commit shifted them). Candidates already filed as GitHub issues
are marked with their number so triage does not re-file them.

**Two corrections affecting this list:** (A) upstream's protocol-doc note claiming
key agreement "is not yet implemented" is stale — master's DUP has handled
`keyAgreement/0`–`/4` since before 07-29 (`control_handler.ex:171+`, shared secret
stored by `0da8336`) — so `key-agreement-protocol` now cites code, not doc.
(B) the previously-listed already-filed item "#10 external-bus intake design (TODO
store.go:135)" rests on a phantom TODO — see survey correction B; #10 needs a human
re-check before anyone treats it as actionable or resolved.

---

### encrypted-mapping-attribute
Title: Accept and persist the new `encrypted` mapping attribute instead of rejecting such interfaces at install time
Scope: mule-line
Body:
**New today**, caused by upstream `372a2cb` (#1999, landed on master 08-21): interfaces
whose mappings carry `"encrypted": true` are now valid upstream and install + persist
(RM schema `mapping.ex:113-119`, boolean default false; persisted via
`interfaces/queries.ex`). Astrate's strict parser rejects the whole interface at parse
time — unknown field (`pkg/interfaceschema/parse.go:121`
`dec.DisallowUnknownFields()`; `mappingJSON` at `parse.go:99-113` has no `Encrypted`),
error wrapping `ErrInvalid`. Net effect: interface documents that install fine on
upstream fail on Astrate — wire-visible at the installation boundary. Fix: add
`Encrypted *bool` to `mappingJSON`, thread it through compile/store like the other
optional attributes, and document that encryption itself is unsupported until
key-agreement lands (mirroring how `aclhook.go` documents keyAgreement) — accepting
the attribute must not imply Astrate can decrypt anything. Verification: T1 test that
an interface with `"encrypted": true` parses (and round-trips the flag), plus one
asserting the documented unsupported-encryption note.

---

### key-agreement-protocol
Title: Handle the `<realm>/<device_id>/control/keyAgreement` CBOR topic
Scope: mule-spec
Body:
Carried forward, **body updated** (re-grep-confirmed against `origin/main`: only the
ACL deny comment `internal/broker/aclhook.go:63-68` and its test `aclhook_test.go`;
zero key machinery anywhere in `internal/`). Correction A: the comment's citation of
the protocol doc's "not yet implemented" note is stale — upstream's DUP implements the
handshake (`data_updater_plant/.../control_handler.ex:171+`: 5 CBOR messages over

$ git -C /root/astrate log origin/mule/research --oneline --all -- '.mule/research/upstream-parity*' 2>/dev/null; echo '---'; git -C /root/astrate log origin/mule/research --oneline --all -- '.mule/research/upstream*' 2>/dev/null | head -10
---

$ find /root/astrate -path '*/.git' -prune -o -type f \( -name '*parity*' -o -name '*upstream*' \) -print 2>/dev/null | head -30; echo '--- docs ---'; ls /root/astrate/docs 2>/dev/null
/root/astrate/.mule/recipes/astarte-upstream.md
--- docs ---
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui

$ gh issue list --repo astrate-platform/astrate --state open --limit 100 2>/dev/null | head -60
104	OPEN	mule: nothing has landed in 14h	mule-alarm	2026-09-16T11:12:14Z
103	OPEN	mule: nothing has landed in 16h	mule-alarm	2026-09-14T11:27:08Z
102	OPEN	mule: nothing has landed in 14h	mule-alarm	2026-09-13T11:09:08Z
101	OPEN	mule: nothing has landed in 14h	mule-alarm	2026-09-12T10:58:09Z
100	OPEN	mule: nothing has landed in 15h	mule-alarm	2026-09-11T10:55:09Z
99	OPEN	mule: nothing has landed in 14h	mule-alarm	2026-09-10T11:08:07Z
98	OPEN	mule: nothing has landed in 16h	mule-alarm	2026-09-09T11:05:08Z
97	OPEN	mule: nothing has landed in 15h	mule-alarm	2026-09-08T11:19:07Z
96	OPEN	mule: nothing has landed in 14h	mule-alarm	2026-09-07T11:15:07Z
95	OPEN	mule: nothing has landed in 15h	mule-alarm	2026-09-06T11:17:07Z
94	OPEN	mule: nothing has landed in 15h	mule-alarm	2026-09-05T11:00:06Z
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
1	OPEN	Provide an Open Source IoT Platform unironically, unironically	wontfix	2026-03-13T21:55:11Z

! permission requested: external_directory (/root/astrate/docs/*); auto-rejecting
✗ Read /root/astrate/docs/UPSTREAM-EXPERIMENTAL.md failed
Error: The user rejected permission to use this specific tool call.

$ gh issue view 92 2>&1; echo '====90===='; gh issue view 90 2>&1
title:	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds
state:	OPEN
author:	atsetilam (Giulio Malitesta)
labels:	enhancement, upstream-experimental, upstream-parity
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
number:	92
--
Issue #51 was closed on 2026-08-22 with: *"Full keyAgreement/encrypted-mappings support stays parked until the upstream 1.4 experimental spec stabilizes — reopen or file fresh when it does."* Filing fresh: that condition has now partly fired.

**What changed.** Upstream commit `d084308` ("docs: add encrypted endpoints key agreement (#2067)", 2026-08-31) added `doc/pages/architecture/082-key_agreement_protocol.md` — a 267-line full wire protocol — and rewrote the Key Agreement section of `080-mqtt-v1-protocol.md`.

**Verified against upstream `master` today, not taken from the survey:**

- `082-key_agreement_protocol.md` exists and is 13,832 bytes.
- The sentences the old parking rested on — "reserved and routed correctly, but the handshake protocol is not yet implemented" and "accepted and acknowledged by Astarte, but no response is sent until the feature is fully implemented" — are **gone from both 080 and 082** (grep count 0 in each).
- 082 now states the handshake "is implemented on the Astarte side (`astarte_data_updater_plant`)", with only device-side session-key usage still described as liable to evolve.

**What the spec now pins down**, where before there was only a sketch: five QoS-2 topics `control/keyAgreement/0..4` (InitExchange, ExchangeResp, SecretHash, HashOk, ExchangeFailed — the 080 table moved from a single `control/keyAgreement` at QoS 1); CBOR message bodies with CDDL given; `alg` 0 = `ECDH_P256-HKDF_SHA256-AES_256_GCM`; a CBOR-wrapped COSE_Key pubkey plus a 32-byte HkdfSalt; session-scoped keys never reused after reconnect or clean reset; AES-256-GCM whole-document ciphertext; and on "no key established yet" Astarte logs and **discards** rather than disconnecting. `ExchangeFailed` error codes are enumerated (0 internal, 1 invalid argument, 2 hash mismatch, 3 unprocessable entity).

**The open question is scope, not feasibility, so this is not a mule task.** The spec is published and implemented upstream, but still ships only in `v1.4.0-rc.5` — `v1.3.3` remains the newest stable tag, and Astrate targets 1.2.2. So "stabilized" is true of the document and false of the release. Three ways to go:

1. Implement now against the rc spec — largest surface in the parity backlog (CBOR codec, X25519/P-256, HKDF, AES-256-GCM, a 5-state handshake machine, shared-secret persistence, five new error names), and it would be built against a spec that may still move.
2. Keep parking until `v1.4.0` is a stable tag, and re-file then.
3. Do the narrow piece only — see the separate mule issue for the stale ACL comment, which is worth fixing regardless of which of these is chosen.

Recorded in `.mule/for-giulio.md`. Survey source: `survey-2026-08-31.md` on `mule/research`.
====90====
title:	milestone-4.0: decide on bumping APICompatVersion from "1.2.2" to "1.3.0-rc.0" (or making it configurable)
state:	CLOSED
author:	atsetilam (Giulio Malitesta)
labels:	enhancement, milestone-4.0
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
number:	90
--
Reminder of a pending decision parked during the upstream-parity investigation (#47–#89). **Do not raise the version now** — this is a milestone-4.0 item.

## Context

`realm.APICompatVersion` (`internal/realm/service.go`) is currently hardcoded to `"1.2.2"` — the emulated upstream API level Astrate declares.

Edgehog 0.13.1 gates just-in-time provisioning triggers (`device_registered`, `device_deletion_started`) behind a Realm Management version >= 1.3. Astrate's engine already supports those trigger kinds end-to-end — verified in #28 (Edgehog compatibility investigation) — but Edgehog refuses to use them because of the declared version string.

## What this issue must cover when picked up

1. **Audit first**: what else does declaring "1.3" imply for other consumers?
   - `astarte-dashboard` is currently validated against 1.2.2 behavior (see #13, error-name bundle comparison). Check whether any dashboard code path feature-gates on RM version besides the policies UI gate (>= 1.1.1).
   - Sweep the parity findings (#47–#89): several are "1.3/1.4 features" we deliberately do NOT implement yet — does declaring 1.3 overpromise any of them?
   - Check device SDKs and `astartectl` for version-gated behavior.
2. **Then decide one of**:
   - bump the constant to `"1.3.0-rc.0"` (upstream's own pre-release numbering for that line);
   - or make it configuration (`internal/config`) so operators can declare the level their consumers expect.
3. Update `docs/COMPATIBILITY.md` deviation #10 wording accordingly (propose, don't edit directly).

## Outcome if done today (why it's parked)

Bumping without the audit risks silently promising surfaces we haven't verified; not bumping keeps Edgehog's JIT provisioning disabled even though the engine supports it. Milestone-4.0 scope decision.

$ git log --oneline --all -15 -- .mule/research/ 2>/dev/null; echo "---branches---"; git branch -a; echo "---finding file in history---"; git log --all --oneline --diff-filter=D -- '.mule/research/*' -3 2>/dev/null
916796b survey: triage the 2026-07-28..2026-09-04 backlog
6d8cae6 mule: daily cross-project survey, 2026-09-04
580faf1 mule: daily cross-project survey, 2026-09-02 (no material change)
8d44f59 mule: daily cross-project survey, 2026-09-01 (no material change)
1c618eb mule: daily cross-project survey, 2026-08-31
ce759e1 survey: recover 16 daily reports the broken commit check left staged (2026-07-28..2026-08-26), and take the tooling fix from main
f92d42a mule/research: commit the calibration survey, and mark it triaged
c6e0301 mule/research: seed the branch the daily survey timer will commit to
---branches---
  main
* mule/queue
  remotes/origin/HEAD -> origin/main
  remotes/origin/feat/pokemon-agent
  remotes/origin/fix/blocked-mule-issues
  remotes/origin/m12-06b-watch-path-onto-main
  remotes/origin/main
  remotes/origin/mule/queue
  remotes/origin/mule/research
  remotes/origin/wip/DietPi-20260727T171543Z
  remotes/origin/worktree-m12-05-match-upstream
  remotes/origin/worktree-m12-06-ach-join
  remotes/origin/worktree-m12-06b-watch-path
---finding file in history---

Much context already exists. Let me read the two most recent milestones-proposal reports directly (the earlier output got truncated), and check today's upstream release state — the v3.0 gate is a stable v1.4.0.
$ echo "=== 2026-09-19 ==="; cat .mule/reports/log/20260919T211455Z-recipe-milestones-proposed.md; echo; echo "=== 2026-09-18 19:36 ==="; cat .mule/reports/log/20260918T193641Z-recipe-milestones-proposed.md; echo; echo "=== 2026-09-18 11:00 ==="; cat .mule/reports/log/20260918T110039Z-recipe-milestones-proposed.md
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_0c08248d4001p4hQUN77VqqAOK

69	CLOSED	enhancement,upstream-parity	Verify unknown-realm HTTP status on RM endpoints against upstream
68	CLOSED	enhancement,mule-blocked,upstream-parity,upstream-experimental	Decide async_operation=false params vs documented always-sync deviation
67	CLOSED	enhancement,upstream-parity,upstream-experimental	Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4)
66	CLOSED	enhancement,upstream-parity	Realm Management: detailed=true interface listing with full mappings (upstream 1.4)
65	CLOSED	enhancement,upstream-parity	Policies: handler-overlap rejection + retry_times coupling + prefetch_count
64	CLOSED	enhancement,upstream-parity	Triggers: decide AMQP action behavior (validate-reject vs NATS-forward deviation)
63	CLOSED	enhancement,upstream-parity	Triggers: HTTP action validation limits (URL/method/header blocklist/template size)
62	CLOSED	enhancement,upstream-parity	Realm Management: audit install/update/delete error codes and statuses
61	CLOSED	enhancement,upstream-parity	Realm Management: audit interface/mapping validation matrix against astarte_core
60	CLOSED	enhancement,upstream-parity	Realm Management: GET config/datastream_maximum_storage_retention (since upstream 1.2.0)
59	CLOSED	enhancement,upstream-parity	AppEngine: group create-body validation + UUID-v1 from_token for group device listing
58	CLOSED	enhancement,upstream-parity	AppEngine: PATCH requires Content-Type application/merge-patch+json
57	CLOSED	enhancement,upstream-parity	AppEngine: audit server-write error taxonomy against upstream
56	CLOSED	enhancement,upstream-parity	AppEngine: object-aggregate response formats table and disjoint_tables
55	CLOSED	enhancement,upstream-parity	AppEngine: complete data-query parameter set
54	CLOSED	enhancement,upstream-parity	AppEngine: GET /groups/{g} show-group endpoint
53	CLOSED	enhancement,upstream-parity	AppEngine: group-scoped device PATCH + interface data access under groups/{g}/devices/{d}
52	CLOSED	enhancement,upstream-parity	AppEngine: devices-by-alias full mirror (PATCH via alias + interface data access)
51	CLOSED	enhancement,upstream-parity	keyAgreement control topic + encrypted mappings (upstream 1.4 experimental) — investigate/park
50	CLOSED	enhancement,upstream-parity	Broker/engine hardening: discard malformed introspection/purge/capabilities instead of crashing
49	CLOSED	enhancement,upstream-parity	MQTT v1: binaryblob properties re-send as BSON subtype 0 in consumer/properties
48	CLOSED	enhancement,upstream-parity	MQTT v1: tolerate devices connecting with an empty introspection
47	CLOSED	enhancement,upstream-parity	MQTT v1: accept device capabilities incl. purge_properties_compression_format (plaintext purge)

$ ls -la /root/astrate-mule/.mule/research/ 2>&1; echo "---"; ls /root/astrate-mule/.mule/
ls: cannot access '/root/astrate-mule/.mule/research/': No such file or directory
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

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:8][] | "\(.tag_name)\t\(.prerelease)\t\(.name)\t\(.published_at)"'
v1.3.4	false	v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	true	v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	false		2026-08-07T14:33:58Z
v1.4.0-rc.4	true	v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	true	v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	true	v1.4.0-rc.2	2026-07-14T10:16:22Z
v1.4.0-rc.1	true	v1.4.0-rc.1	2026-07-14T08:53:27Z
v1.3.2	false	v1.3.2	2026-07-14T10:07:26Z

$ gh issue list --label milestone-3.0 --state all --limit 50 --json number,title,state --jq '.[] | "\(.number)\t\(.state)\t\(.title)"' ; echo "---milestone-v3.0---"; gh issue list --label milestone-v3.0 --state all --limit 50 --json number,title,state --jq '.[] | "\(.number)\t\(.state)\t\(.title)"' ; echo "---end---"
---milestone-v3.0---
---end---

The ledger lives under `upstream-parity`. Open items: #92 (keyAgreement, escalated), #93 (aclhook, own recipe path), #78 (FDO, milestone-4.0). Let me run the standing-item check and confirm nothing changed since the last run.
→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ gh issue view 92 --json title,state,labels,body --jq '{title,state,labels:[.labels[].name]}' ; echo "---93---"; gh issue view 93 --json title,state,labels --jq '{title,state,labels:[.labels[].name]}'; echo "---78---"; gh issue view 78 --json title,state,labels --jq '{title,state,labels:[.labels[].name]}'
{"labels":["enhancement","upstream-parity","upstream-experimental"],"state":"OPEN","title":"keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds"}
---93---
{"labels":["mule-review","upstream-parity"],"state":"OPEN","title":"In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged"}
---78---
{"labels":["enhancement","milestone-4.0","upstream-parity"],"state":"OPEN","title":"FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)"}

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+
+---
+
 - **deviceid wire-parity concern, from the 2026-09-18 code review: Astrate is likely stricter than upstream on the wire form's unused bits.** `pkg/deviceid/deviceid.go:37` uses `base64.RawURLEncoding.Strict()`, so `Parse` rejects `"AAAAAAAAAAAAAAAAAAAAAB"` (unused low 4 bits of the 22nd char set; rejection pinned at deviceid_test.go:176), while the comment at deviceid.go:34-36 claims parity with `Elixir Base.url_decode64!(padding: false)` — which decodes the 128 payload bits and discards the unused ones rather than erroring. If so, a device registered with a non-canonical-trailing-bits id on upstream Astarte is accepted there but rejected by Astrate's `Parse` — a wire-compat break in a package whose purpose is the wire form. Not verified on the Pi (no Docker/Elixir); a `[legion]` probe (`deviceid-trailing-bits-upstream-probe`) is queued to measure it. Once confirmed, the call is yours: relax `Parse` to upstream's decode-and-ignore, or keep strict canonicality as the package's contract.
 
 ---
 

Standing-item check done — both UPSTREAM-EXPERIMENTAL rows (#67, #68) are still tagged "1.4 experimental", gated on v1.4 final. The ledger is consistent with the same-day earlier run: only #92/#93/#78 open, all escalated or on their own recipe path. Release sweep shows v1.4.0 still rc.5-only (newest stable v1.3.4, today, maintenance-only). No new machine-checkable gaps exist. Writing the run entry to for-giulio.md.
Now I follow the recipe's deliverable decision. Since there are no new gaps (the 1.3 surface is fully delivered, the two UPSTREAM-EXPERIMENTAL rows are gated on v1.4 final which is still rc.5-only, and the only open `upstream-parity` items are already escalated), no issues get filed and no task lines get appended — the escalated "your call when v1.4.0 ships" line is the standing decision request. Let me verify the file and git state.
$ git status --short; echo "---"; git log --oneline -3
M .mule/for-giulio.md
---
eb95cdd mule: recipe hygiene timed out
4c4cf5b mule: recipe docs-sync ran (421s)
f6dc460 mule: refill the queue

→ Read .mule/for-giulio.md [limit=20]
The milestone recipe is complete. The first non-DONE milestone is **v3.0** (upstream 1.2.2 → 1.3/1.4 parity); scope is decided, so I investigated rather than proposing candidates.

**Result:** no new gaps exist, so no `gh issue create` calls and no todo.md lines were produced.

- Release sweep: still **no stable v1.4.0** — newest stable is today's `v1.3.4` (maintenance-only, no wire change), newest overall `v1.4.0-rc.5`.
- Open `upstream-parity` items are exactly the already-escalated set: #92 (keyAgreement), #93 (aclhook, own recipe path), #78 (FDO, milestone-4.0).
- Standing item: both UPSTREAM-EXPERIMENTAL rows (#67, #68) still gated on v1.4 final; the 1.3 surface is fully delivered.
- `APICompatVersion` bump (#90) stays the final-phase item.

The one deliverable is the `for-giulio.md` escalation entry recording the run's state and your decision point (when stable v1.4.0 ships: answer #92, run the final-phase bump, cut the tag).

Done: v3.0 verified unchanged — still gated on upstream v1.4 final; no new gaps, nothing filed.
Files: .mule/for-giulio.md
Verified: gh issue list/search, gh api releases, docs/UPSTREAM-EXPERIMENTAL.md -> pass
Unsure: nothing
Follow-ups: none — deliberate, per the recipe the "cut the tag" line is not proposed while the v1.4 gate is unresolved

=== 2026-09-18 11:00 ===
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

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | "\(.tag_name) \(.prerelease) \(.published_at)"' 2>/dev/null | head -10
v1.3.4 false 2026-09-18T09:42:57Z
v1.4.0-rc.5 true 2026-08-20T12:28:04Z
v1.3.3 false 2026-08-07T14:33:58Z
v1.4.0-rc.4 true 2026-08-06T13:40:20Z
v1.4.0-rc.3 true 2026-07-31T10:48:39Z
v1.4.0-rc.2 true 2026-07-14T10:16:22Z
v1.4.0-rc.1 true 2026-07-14T08:53:27Z
v1.3.2 false 2026-07-14T10:07:26Z
v1.3.1 false 2026-07-13T11:50:05Z
v1.3.0 false 2026-05-06T15:30:50Z

The picture is clear from prior runs. Let me read the most recent (2026-09-19) report precisely so this run doesn't duplicate it, and confirm nothing's changed since.
→ Read .mule/for-giulio.md [limit=40]
→ Read .mule/reports/log/20260919T211455Z-recipe-milestones-proposed.md
Previous runs have this milestone well-characterized. This run (2026-09-20) needs a fresh verification since yesterday's run. Let me check the two gates: upstream releases (stable v1.4.0?) and the open issue set.

$ echo "=== upstream-parity open ==="; gh issue list --label upstream-parity --state open --json number,title --jq '.[] | "\(.number)\t\(.title)"'; echo "=== ALL open issues (non-labelfilter) newest 5 ==="; gh issue list --state open --limit 100 --json number,title,labels,createdAt --jq '.[] | "\(.number)\t\(.labels|map(.name)|join(","))\t\(.title)"' | sort -t $'\t' -k1 -n | tail -8
=== upstream-parity open ===
93	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
92	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds
78	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
=== ALL open issues (non-labelfilter) newest 5 ===
sort: multi-character tab ‘$\\t’

$ echo "=== newest release ==="; gh api repos/astarte-platform/astarte/releases --jq '.[0] | {tag_name, prerelease, published_at, name}'; echo "=== stable line ==="; gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.prerelease==false) | .tag_name' | sort -V | tail -3; echo "=== prerelease line ==="; gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.prerelease==true) | "\(.tag_name) \(.published_at)"' | sort -V | tail -3
=== newest release ===
{"name":"v1.3.4","prerelease":false,"published_at":"2026-09-18T09:42:57Z","tag_name":"v1.3.4"}
=== stable line ===
v1.3.2
v1.3.3
v1.3.4
=== prerelease line ===
v1.4.0-rc.3 2026-07-31T10:48:39Z
v1.4.0-rc.4 2026-08-06T13:40:20Z
v1.4.0-rc.5 2026-08-20T12:28:04Z

$ gh issue list --state open --limit 100 --json number,title,labels --jq '.[] | "\(.number)|\(.labels|map(.name)|join(","))|\(.title)"' | sort -t '|' -k1 -n | tail -10
98|mule-alarm|mule: nothing has landed in 16h
99|mule-alarm|mule: nothing has landed in 14h
100|mule-alarm|mule: nothing has landed in 15h
101|mule-alarm|mule: nothing has landed in 14h
102|mule-alarm|mule: nothing has landed in 14h
103|mule-alarm|mule: nothing has landed in 16h
104|mule-alarm|mule: nothing has landed in 14h
105|mule-alarm|mule: nothing has landed in 14h
106|mule-alarm|mule: nothing has landed in 14h
107|mule-alarm|mule: nothing has landed in 14h

Nothing has changed since yesterday's run: still no stable v1.4.0 upstream (newest stable `v1.3.4`, newest overall `v1.4.0-rc.5`), and the open `upstream-parity` set is unchanged (#92 parked decision, #93 mule-review on its own path, #78 milestone-4.0 FDO). Both UPSTREAM-EXPERIMENTAL rows remain gated on v1.4 final. Per the recipe's bucket rule the open items are all already escalated, so no new gap exists to file or queue.

Following the established pattern for this recipe's no-op runs, I'll record today's verification in `.mule/for-giulio.md` (newest at top).
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+
+---
+
 - **COMPATIBILITY.md §8 wording proposal: trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (Forwarder) actions are single-shot by design, receiving only `maximum_capacity`/`event_ttl`.** Propose adding to the "Trigger delivery policies" deviation: *"Policy `error_handlers` (retry/discard) and the `retry_times` attempt cap govern HTTP webhook actions only: a custom action delivered through the Forwarder seam is single-shot by design, receiving from the policy solely its `maximum_capacity` in-flight bound and its `event_ttl` staleness drop, with a forward failure counted once and never retried."* Verified in `internal/engine/triggers/actions.go`: `Action.Custom` routes to `forward` (actions.go:435-439), which never consults `policy.Decide` — the Forwarder seam returns only an error, no HTTP status (actions.go:161-167, doc comment at actions.go:471-479) — while `Enqueue` bounds `maximum_capacity` on every delivery (actions.go:338-349) and `deliver` drops past `event_ttl` before the custom/webhook branch (actions.go:400-416). The existing warn-on-single-shot line already tells the operator when a policy's `retry_times` cannot apply to the forward path (actions.go:491-497).
 
 ---
 

→ Read .mule/for-giulio.md [limit=24]
Recipe carried out for the active milestone **v3.0** (first non-DONE in `.mule/milestones.md`). Investigation complete, no new gaps found, deliverable recorded.

**What I did:** read MULE.md and the milestones recipe, identified v3.0 as the current target, investigated the reference and current state:
- Release sweep: still **no stable v1.4.0** — newest stable `v1.3.4` (2026-09-18, maintenance-only), newest overall `v1.4.0-rc.5`. Nothing shipped since the 2026-09-19 run.
- `gh issue list`: `milestone-3.0` label empty (ledger lives under `upstream-parity`); open items are exactly the already-escalated #92 (keyAgreement, parked on stable v1.4.0 per Giulio), #93 (aclhook, mule-review, own recipe path), #78 (FDO, milestone-4.0). #47–#89 rest all closed.
- Standing item: both UPSTREAM-EXPERIMENTAL rows (#67, #68) still wait on v1.4 final; `APICompatVersion` stays 1.2.2 per #90's frozen decision.
- Result: no machine-checkable gap exists → **no `gh issue create`**, **no todo.md task lines**. Deliberately not proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved.

Done: v3.0 verified unchanged — still gated on upstream v1.4 final; no new gaps, nothing filed
Files: .mule/for-giulio.md
Verified: gh api releases sweep + gh issue list (upstream-parity / milestone-3.0) + docs/UPSTREAM-EXPERIMENTAL.md -> pass (state unchanged)
Unsure: nothing
Follow-ups: none — the standing "when stable v1.4.0 ships: answer #92, run final-phase bump, cut v3.0 tag" escalation is re-recorded at the top of for-giulio.md
