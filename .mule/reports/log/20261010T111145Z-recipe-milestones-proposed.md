slug: recipe-milestones
verdict: proposed
at:  69f3cff
ran: 2026-10-10T11:11:45Z on DietPi in 403s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
→ Read .mule/for-giulio.md
→ Read .mule/todo.md
I'll work through the milestones recipe. First, let me establish current state: the target milestone (v3.0), what's already filed, and the reference docs.
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ gh issue list --label milestone-3.0 --state all --limit 50 2>&1; echo "---UPSTREAM-PARITY---"; gh issue list --label upstream-parity --state all --limit 100 2>&1 | head -120; echo "---MILESTONE-4.0---"; gh issue list --label milestone-4.0 --state all --limit 50 2>&1
---UPSTREAM-PARITY---
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
---MILESTONE-4.0---
90	CLOSED	milestone-4.0: decide on bumping APICompatVersion from "1.2.2" to "1.3.0-rc.0" (or making it configurable)	enhancement, milestone-4.0	2026-08-23T21:06:34Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
28	CLOSED	milestone-4.0: investigate Edgehog compatibility with Astrate (readonly)	mule, readonly, milestone-4.0	2026-08-23T15:48:58Z

✱ Glob ".mule/research/**" 0 matches

$ gh issue view 90 --json number,title,state,closedAt,labels,body 2>&1 | head -60; echo "===92==="; gh issue view 92 --json number,title,state,labels,body 2>&1 | head -40
{"body":"Reminder of a pending decision parked during the upstream-parity investigation (#47–#89). **Do not raise the version now** — this is a milestone-4.0 item.\n\n## Context\n\n`realm.APICompatVersion` (`internal/realm/service.go`) is currently hardcoded to `\"1.2.2\"` — the emulated upstream API level Astrate declares.\n\nEdgehog 0.13.1 gates just-in-time provisioning triggers (`device_registered`, `device_deletion_started`) behind a Realm Management version >= 1.3. Astrate's engine already supports those trigger kinds end-to-end — verified in #28 (Edgehog compatibility investigation) — but Edgehog refuses to use them because of the declared version string.\n\n## What this issue must cover when picked up\n\n1. **Audit first**: what else does declaring \"1.3\" imply for other consumers?\n   - `astarte-dashboard` is currently validated against 1.2.2 behavior (see #13, error-name bundle comparison). Check whether any dashboard code path feature-gates on RM version besides the policies UI gate (>= 1.1.1).\n   - Sweep the parity findings (#47–#89): several are \"1.3/1.4 features\" we deliberately do NOT implement yet — does declaring 1.3 overpromise any of them?\n   - Check device SDKs and `astartectl` for version-gated behavior.\n2. **Then decide one of**:\n   - bump the constant to `\"1.3.0-rc.0\"` (upstream's own pre-release numbering for that line);\n   - or make it configuration (`internal/config`) so operators can declare the level their consumers expect.\n3. Update `docs/COMPATIBILITY.md` deviation #10 wording accordingly (propose, don't edit directly).\n\n## Outcome if done today (why it's parked)\n\nBumping without the audit risks silently promising surfaces we haven't verified; not bumping keeps Edgehog's JIT provisioning disabled even though the engine supports it. Milestone-4.0 scope decision.","closedAt":"2026-08-23T21:06:34Z","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACtvELrg","name":"milestone-4.0","description":"","color":"5319e7"}],"number":90,"state":"CLOSED","title":"milestone-4.0: decide on bumping APICompatVersion from \"1.2.2\" to \"1.3.0-rc.0\" (or making it configurable)"}
===92===
{"body":"Issue #51 was closed on 2026-08-22 with: *\"Full keyAgreement/encrypted-mappings support stays parked until the upstream 1.4 experimental spec stabilizes — reopen or file fresh when it does.\"* Filing fresh: that condition has now partly fired.\n\n**What changed.** Upstream commit `d084308` (\"docs: add encrypted endpoints key agreement (#2067)\", 2026-08-31) added `doc/pages/architecture/082-key_agreement_protocol.md` — a 267-line full wire protocol — and rewrote the Key Agreement section of `080-mqtt-v1-protocol.md`.\n\n**Verified against upstream `master` today, not taken from the survey:**\n\n- `082-key_agreement_protocol.md` exists and is 13,832 bytes.\n- The sentences the old parking rested on — \"reserved and routed correctly, but the handshake protocol is not yet implemented\" and \"accepted and acknowledged by Astarte, but no response is sent until the feature is fully implemented\" — are **gone from both 080 and 082** (grep count 0 in each).\n- 082 now states the handshake \"is implemented on the Astarte side (`astarte_data_updater_plant`)\", with only device-side session-key usage still described as liable to evolve.\n\n**What the spec now pins down**, where before there was only a sketch: five QoS-2 topics `control/keyAgreement/0..4` (InitExchange, ExchangeResp, SecretHash, HashOk, ExchangeFailed — the 080 table moved from a single `control/keyAgreement` at QoS 1); CBOR message bodies with CDDL given; `alg` 0 = `ECDH_P256-HKDF_SHA256-AES_256_GCM`; a CBOR-wrapped COSE_Key pubkey plus a 32-byte HkdfSalt; session-scoped keys never reused after reconnect or clean reset; AES-256-GCM whole-document ciphertext; and on \"no key established yet\" Astarte logs and **discards** rather than disconnecting. `ExchangeFailed` error codes are enumerated (0 internal, 1 invalid argument, 2 hash mismatch, 3 unprocessable entity).\n\n**The open question is scope, not feasibility, so this is not a mule task.** The spec is published and implemented upstream, but still ships only in `v1.4.0-rc.5` — `v1.3.3` remains the newest stable tag, and Astrate targets 1.2.2. So \"stabilized\" is true of the document and false of the release. Three ways to go:\n\n1. Implement now against the rc spec — largest surface in the parity backlog (CBOR codec, X25519/P-256, HKDF, AES-256-GCM, a 5-state handshake machine, shared-secret persistence, five new error names), and it would be built against a spec that may still move.\n2. Keep parking until `v1.4.0` is a stable tag, and re-file then.\n3. Do the narrow piece only — see the separate mule issue for the stale ACL comment, which is worth fixing regardless of which of these is chosen.\n\nRecorded in `.mule/for-giulio.md`. Survey source: `survey-2026-08-31.md` on `mule/research`.","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"},{"id":"LA_kwDORmfsJs8AAAACxx3GJg","name":"upstream-experimental","description":"Adopted from an upstream experimental/unstable feature — review each upstream release","color":"d4c5f9"}],"number":92,"state":"OPEN","title":"keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds"}

$ ls -la .mule/ 2>&1; echo "---research dir---"; ls -la .mule/research/ 2>&1; echo "---tasks dir---"; ls .mule/tasks/ 2>&1 | head
total 16848
drwxr-xr-x  7 root root     4096 Oct 10 12:50 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root        0 Oct 10 12:50 .alarmed
-rw-r--r--  1 root root       13 Oct 10 13:05 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  9 22:26 .heartbeat
-rw-r--r--  1 root root    12288 Oct 10 13:07 .last-output
-rw-r--r--  1 root root        0 Oct 10 13:04 .lock
-rw-r--r--  1 root root       11 Oct 10 13:05 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 16782607 Oct 10 13:07 cron.log
drwxr-xr-x  2 root root     4096 Oct  9 20:07 failed
-rw-r--r--  1 root root   163554 Oct  9 21:05 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    16410 Oct  9 22:26 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  9 21:54 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   177588 Oct  9 22:26 todo.md
---research dir---
ls: cannot access '.mule/research/': No such file or directory
---tasks dir---
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md

$ git branch -a 2>&1 | head -30; echo "---"; git ls-files '.mule/research/*' 2>&1 | head; echo "---search all refs---"; git log --all --oneline -- '.mule/research/upstream-parity-2026-08-22.md' 2>&1 | head; echo "---find any upstream-parity doc---"; git ls-files '.mule/*' 2>&1 | grep -i 'research\|parity\|milestone' | head -40
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
---
---search all refs---
---find any upstream-parity doc---
.mule/milestones.md
.mule/recipes/milestones.md
.mule/reports/log/20260905T181053Z-recipe-milestones-proposed.md
.mule/reports/log/20260906T180442Z-recipe-milestones-proposed.md
.mule/reports/log/20260907T171928Z-recipe-milestones-proposed.md
.mule/reports/log/20260907T204751Z-recipe-milestones-proposed.md
.mule/reports/log/20260908T193026Z-recipe-milestones-proposed.md
.mule/reports/log/20260908T193127Z-recipe-milestones-proposed.md
.mule/reports/log/20260909T175718Z-recipe-milestones-proposed.md
.mule/reports/log/20260909T230411Z-recipe-milestones-proposed.md
.mule/reports/log/20260910T195304Z-recipe-milestones-proposed.md
.mule/reports/log/20260911T185303Z-recipe-milestones-proposed.md
.mule/reports/log/20260912T171256Z-recipe-milestones-proposed.md
.mule/reports/log/20260913T165246Z-recipe-milestones-proposed.md
.mule/reports/log/20260913T203209Z-recipe-milestones-proposed.md
.mule/reports/log/20260914T234550Z-recipe-milestones-proposed.md
.mule/reports/log/20260915T184406Z-recipe-milestones-proposed.md
.mule/reports/log/20260916T175801Z-recipe-milestones-proposed.md
.mule/reports/log/20260917T184319Z-recipe-milestones-timeout.md
.mule/reports/log/20260918T110039Z-recipe-milestones-proposed.md
.mule/reports/log/20260918T193641Z-recipe-milestones-proposed.md
.mule/reports/log/20260919T211455Z-recipe-milestones-proposed.md
.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md
.mule/reports/log/20260921T183804Z-recipe-milestones-proposed.md
.mule/reports/log/20260922T123011Z-recipe-milestones-proposed.md
.mule/reports/log/20260922T190954Z-recipe-milestones-proposed.md
.mule/reports/log/20260922T202552Z-recipe-milestones-proposed.md
.mule/reports/log/20260923T174341Z-recipe-milestones-proposed.md
.mule/reports/log/20260924T175556Z-recipe-milestones-proposed.md
.mule/reports/log/20260924T192153Z-recipe-milestones-proposed.md
.mule/reports/log/20260925T115624Z-recipe-milestones-timeout.md
.mule/reports/log/20260926T113138Z-recipe-milestones-proposed.md
.mule/reports/log/20260927T191018Z-recipe-milestones-timeout.md
.mule/reports/log/20260928T200921Z-recipe-milestones-proposed.md
.mule/reports/log/20260929T110049Z-recipe-milestones-proposed.md
.mule/reports/log/20260929T175252Z-recipe-milestones-proposed.md
.mule/reports/log/20260929T192156Z-recipe-milestones-proposed.md
.mule/reports/log/20260930T165653Z-recipe-milestones-proposed.md
.mule/reports/log/20260930T184654Z-recipe-milestones-proposed.md
.mule/reports/log/20260930T194054Z-recipe-milestones-proposed.md

$ git cat-file -s origin/mule/research:.mule/research/upstream-parity-2026-08-22.md 2>&1; echo "---list research dir on that branch---"; git ls-tree origin/mule/research:.mule/research/ 2>&1; echo "---recent milestone reports---"; ls -t .mule/reports/log/*milestones* 2>&1 | head -6
fatal: path '.mule/research/upstream-parity-2026-08-22.md' does not exist in 'origin/mule/research'
---list research dir on that branch---
100644 blob 8e9d610f99cfbfb0940c6730cbb9a8b434b99dbc	README.md
100644 blob 118b479a124bc5a7717e1bf45b155c98a90282dd	issues-2026-07-27.md
100644 blob 9a61876cef66e024c1bf8cbec37a399ccbded045	issues-2026-07-28.md
100644 blob b4caf34f968c87ef6b7ac43f82f17207e4281cd8	issues-2026-07-30.md
100644 blob 3bc1f2c43c133ccdad71429871a9d5a170e292d6	issues-2026-08-01.md
100644 blob fe381f0b14c47b4ec724bc2a613c74f0d62c705b	issues-2026-08-06.md
100644 blob 15c3d3f538a0f4e8e5bc2b794ec7dc8140f9194c	issues-2026-08-07.md
100644 blob 6dccbbd40d1cbc33a3ba320f1ed329d98a75f6e5	issues-2026-08-08.md
100644 blob 6768ea33dd31618aecbc877375c2f68388da3448	issues-2026-08-13.md
100644 blob 8aeee514e507a06c0f6983f82421da035778b6a0	issues-2026-08-15.md
100644 blob 77c6a4afc1a3cc1d2870fe43b67dfb4dab0ea270	issues-2026-08-18.md
100644 blob 5c1ce408a56e9e06a8c5ad0ddc1e99d88edce0a3	issues-2026-08-21.md
100644 blob 0edd551c5eca6f6ca89126ea1f80be58cf09979a	issues-2026-08-22.md
100644 blob 72452b90ce21503d292e3b28f5beb8b3442df237	issues-2026-08-23.md
100644 blob 66bc4a26da9300ee8a0a64b8d8457bd4f68fe950	issues-2026-08-24.md
100644 blob 0769696b8c1b9b5b8d5fefafdf864d2e3419cb2e	issues-2026-08-25.md
100644 blob 6bd49012a7b0623448a4c611402160a478d67df5	issues-2026-08-26.md
100644 blob 3d93956614bbae53fcabc76a8b94a3a45cd0a0cf	issues-2026-09-04.md
100644 blob 32617005335a21bfb79c5ff58bc9a2e5803df607	log.md
100644 blob fdf2c0d3ba1a8277a5964767d976b4761f923b33	survey-2026-07-27.md
100644 blob 95e854d8ed99979b4a2dbf32870b52a18ba3fdc4	survey-2026-07-28.md
100644 blob 48376d7b4097de0a68de5e7d6f95bf6fb0495a36	survey-2026-07-30.md
100644 blob 689cf1df58260396c2052aab937efc301bcc710f	survey-2026-08-01.md
100644 blob 1c0997af2482002cb8d43b189d5a49c6d374b3b9	survey-2026-08-06.md
100644 blob 97cae62ccaeddadae27fdafe371915cb42d92f54	survey-2026-08-07.md
100644 blob edf791a7f6b1317bbf10133b164cc6e668690c76	survey-2026-08-08.md
100644 blob 0e83a661338971695b21e76521e94774bcafe345	survey-2026-08-13.md
100644 blob d9cbf59451aa43a9c562943898854d49e00196f1	survey-2026-08-15.md
100644 blob 7a8ee23575f1708be019a64bc6dd3888459d2419	survey-2026-08-18.md
100644 blob 35bd411df1b4e7eefcf0f3b0ca07fe000c26d655	survey-2026-08-21.md
100644 blob 7a7919e92471f51b58c01c7d5867885a9f897914	survey-2026-08-22.md
100644 blob 81960b18de2df20b74fed6e39938c62e3cb0a617	survey-2026-08-23.md
100644 blob 3729ecded047f0ae9bd7d246128d490a492cb13d	survey-2026-08-24.md
100644 blob be8614ca3594bab96396b284eb6a77eca53f8d46	survey-2026-08-25.md
100644 blob 7b0aa067e633ae5c5a03f1b9fc97520d21f408b3	survey-2026-08-26.md
100644 blob 30efa98f7761336032a7797bc51d5a894866436c	survey-2026-08-31.md
100644 blob 2e7b20f33af19d3fdb7aa2efb06032a90c1b15c4	survey-2026-09-04.md
100644 blob efdbcc1edc1118987ce9953a7d01ef145f74c0b4	triaged.md
---recent milestone reports---
.mule/reports/log/20261009T183818Z-recipe-milestones-proposed.md
.mule/reports/log/20261008T202151Z-recipe-milestones-proposed.md
.mule/reports/log/20261007T224444Z-recipe-milestones-proposed.md
.mule/reports/log/20261007T113844Z-recipe-milestones-proposed.md
.mule/reports/log/20261006T141129Z-recipe-milestones-proposed.md
.mule/reports/log/20261005T193452Z-recipe-milestones-proposed.md

→ Read .mule/reports/log/20261009T183818Z-recipe-milestones-proposed.md

$ echo "=== RELEASES ==="; gh api repos/astarte-platform/astarte/releases --jq '.[] | "\(.tag_name)\tprerelease=\(.prerelease)\t\(.published_at)"' 2>&1 | head -8; echo "=== TAGS ==="; gh api repos/astarte-platform/astarte/tags --jq '.[].name' 2>&1 | head -8; echo "=== master commits since 2026-10-08 ==="; gh api "repos/astarte-platform/astarte/commits?since=2026-10-08T00:00:00Z&per_page=40" --jq '.[] | "\(.sha[0:8])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"' 2>&1 | head -40
=== RELEASES ===
v1.3.5	prerelease=false	2026-10-05T14:11:44Z
v1.4.0-rc.6	prerelease=true	2026-09-30T07:28:27Z
v1.3.4	prerelease=false	2026-09-18T09:42:57Z
v1.4.0-rc.5	prerelease=true	2026-08-20T12:28:04Z
v1.3.3	prerelease=false	2026-08-07T14:33:58Z
v1.4.0-rc.4	prerelease=true	2026-08-06T13:40:20Z
v1.4.0-rc.3	prerelease=true	2026-07-31T10:48:39Z
v1.4.0-rc.2	prerelease=true	2026-07-14T10:16:22Z
=== TAGS ===
v1.4.0-rc.6
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
v1.4.0-rc.0
v1.3.5
=== master commits since 2026-10-08 ===

$ echo "=== milestone-3.0 all ==="; gh issue list --label milestone-3.0 --state all --limit 50 2>&1; echo "=== open upstream-parity ==="; gh issue list --label upstream-parity --state open --limit 50 2>&1; echo "=== open milestone-4.0 ==="; gh issue list --label milestone-4.0 --state open --limit 50 2>&1; echo "=== APICompatVersion ==="; grep -rn 'APICompatVersion = ' internal/realm/service.go 2>&1
=== milestone-3.0 all ===
=== open upstream-parity ===
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
=== open milestone-4.0 ===
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
=== APICompatVersion ===
588:const APICompatVersion = "1.2.2"

$ echo "=== latest master commits ==="; gh api "repos/astarte-platform/astarte/commits?per_page=10" --jq '.[] | "\(.sha[0:8])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"' 2>&1 | head -12; echo "=== for-giulio milestone entries recent ==="; grep -n "Milestone recipe run" .mule/for-giulio.md | head -8
=== latest master commits ===
171c1bd7	2026-10-06T15:57:45Z	chore: forward-port astarte_vmq_plugin from release-1.4
1da04832	2026-10-05T14:57:19Z	chore: forward port release-1.4
a8a26d0f	2026-10-05T14:40:35Z	chore: forward port release-1.3
d85e0ed0	2026-10-06T10:02:32Z	docs: sync FDO workflow to the current behavior (#2241)
a6bd2c21	2026-10-06T09:59:16Z	fix(fdo): reject owner key names not accepted by OpenBao (#2252)
a51ab2e1	2026-10-06T07:21:18Z	refactor: Update Umbrella checks (#2244)
7628d2ed	2026-10-05T15:49:18Z	test(generators): validate interfaces with Core changesets
ed1e99f4	2026-10-05T15:49:04Z	test(generators): validate mappings with Core changesets
bb17455a	2026-10-05T15:48:50Z	test(generators): validate policy handlers with Core changesets
ec020bf8	2026-10-05T15:48:31Z	fix(generators): generate valid trigger policy names
=== for-giulio milestone entries recent ===
50:- **Milestone recipe run, 2026-10-06: v3.0 verified unchanged — still gated on upstream — and the master delta since the 2026-10-05 run is a closed set of 13 commits, every one wire-inert for Astrate. No issue filed, no `.mule/todo.md` line queued.** Release sweep is byte-identical to yesterday's (`gh api repos/astarte-platform/astarte/releases`): **no stable v1.4.0** — newest stable `v1.3.5` (2026-10-05), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition (your 2026-09-04 "wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) stays unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** — this milestone's whole ledger lives under `upstream-parity`, where the open set is unchanged: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`); **#47–#89** otherwise all closed. Non-alarm open issues total exactly four (plus the 25 `mule-alarm` #94–#118): the three above and **#1** (untouched per standing instruction). `APICompatVersion` still `"1.2.2"` (`internal/realm/service.go:588`, re-checked).
67:- **Milestone recipe run, 2026-10-05: v3.0's gate is unchanged, but the release sweep finally has something new to swallow — upstream shipped stable v1.3.5 today, seven hours ago, and it is wire-inert for Astrate. No issue filed, no task line queued.**  **The new release, measured.** `gh api repos/astarte-platform/astarte/releases` now shows a stable **v1.3.5** at `2026-10-05T14:11:44Z` where every prior run of this recipe (2026-09-10 onward, 20+ entries in this file) saw `v1.3.4` (2026-09-18) as newest. `compare/v1.3.4...v1.3.5` is **4 commits, 35 files**: two `chore: update horde` (`8d47574f`, `4c4222cf`), `chore(rm): bump xandra to 0.19.4` (`c4c63ebe`), `chore: prepare v1.3.5` (`8f0f5ef1`). The file list is `mix.exs`/`mix.lock` version bumps across all apps and libs, `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `astarte-dashboard/package.json`, `doc/mix.exs`, `README.md` — plus a **single** CHANGELOG entry, `## [1.3.5] - 2026-10-05` / `### Fixed` / "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available". That is the whole delta and none of it is Astarte-facing: Astrate implements no data_updater_plant and no Erlang RPC server (the only `rpc` hit under `internal/` is a comment at `internal/pairing/http.go:300` about AMQP reply-code conflicts), so no HTTP route, MQTT topic, AMQP control message or schema field moved. Same shape as the v1.3.2/v1.3.3/v1.3.4 maintenance releases already recorded here, so it is **not a v3.0 gap** and there is nothing to file.
88:- **Milestone recipe run, 2026-10-04: the gate is unchanged as always, but the recipe's standing item finally produced a real finding — both rows of `docs/UPSTREAM-EXPERIMENTAL.md` are tagged "1.4 experimental" and neither of the two features upstream calls experimental. It is a wording/policy call on two files that are yours, so this run filed no issue and queued no task.**
109:- **Milestone recipe run, 2026-10-03: v3.0's gate is unchanged, and this run's one new finding is not a gap in Astrate's code — it is that this milestone's declared source document can never be committed, so every run re-derives the milestone from the issue bodies.**
166:- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Status: not started", which is no longer true.** Release sweep (`gh api repos/astarte-platform/astarte/releases`, run this session) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.5` (2026-08-20, prerelease) — byte-identical to what the 2026-09-24 run saw, so nothing shipped upstream in the two days since. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Open non-alarm issues remain exactly **#93** (aclhook comment, `mule-review`, its own recipe path), **#92** (keyAgreement, parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** (FDO, `milestone-4.0`), **#1** (untouched per standing instruction). No new gaps, **no issues filed, no `.mule/todo.md` task lines added** — there is nothing machine-checkable left, and filing anything here would duplicate #92/#93/#78.
241:- **Milestone recipe run, 2026-09-24: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-23 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
249:- **Milestone recipe run, 2026-09-23: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-21 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
265:- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

$ gh issue list --label upstream-parity --state closed --limit 60 --json number,stateReason --jq '.[] | select(.stateReason != "COMPLETED") | "\(.number) \(.stateReason)"' 2>&1; echo "=== count closed/notplanned done ==="; gh issue list --label upstream-parity --state closed --limit 60 --json number,stateReason --jq '[.[] | .stateReason] | group_by(.) | map({reason: .[0], n: length})' 2>&1; echo "=== grep looks complete in for-giulio/reports ==="; grep -rn "looks complete, verify and cut" .mule/for-giulio.md .mule/reports/ 2>&1 | head
=== count closed/notplanned done ===
[{"n":43,"reason":"COMPLETED"}]
=== grep looks complete in for-giulio/reports ===
.mule/for-giulio.md:70:  **Step 5 deliberately not taken.** I am *not* proposing "milestone v3.0 looks complete, verify and cut the tag", for the same reason as the twenty runs before this one: the section's declared scope is **1.2.2 → 1.3/1.4** and only the 1.3 half is delivered (now including v1.3.5, which is inert anyway). The 1.4 half is open — #92 unanswered, both experimental rows unreconciled, `APICompatVersion` not bumped — so the tag must not be cut. Your call when a stable v1.4.0 ships: answer #92, settle the register rows, run the final-phase bump, then cut the tag.
.mule/for-giulio.md:111:  **The recipe's standing item is satisfied and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental*, reconcile trigger "upstream 1.4 final". 1.4 is not final, so there is nothing to promote and nothing to deprecate, and the register is accurate as written. `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked this run) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision ("bumped only in the same change that completes that surface, after reconciling every row tagged with the level"). The milestone-bump issue the recipe asks me to file **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": v3.0's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is delivered, the 1.4 half is open (#92, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut.
.mule/for-giulio.md:168:  **The recipe's standing item is satisfied, and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still has exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental* with reconcile trigger "upstream 1.4 final". v1.4 is not final, so both rows stay as written: nothing to promote, nothing to deprecate, and the register is accurate. `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision; the milestone-bump issue the recipe asks for **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": the milestone's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is done, and the 1.4 half is unresolved (#92 open, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut. Also still true and already escalated on 2026-09-06, so not repeated: the section's source document `.mule/research/upstream-parity-2026-08-22.md` does not exist in the repo.
.mule/for-giulio.md:487:- **milestone 2.0 looks complete, verify and cut the tag** — all 11 `milestone-2.0` issues
.mule/reports/log/20261003T180922Z-recipe-github-issues-proposed.md:232:  **The recipe's standing item is satisfied and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental*, reconcile trigger "upstream 1.4 final". 1.4 is not final, so there is nothing to promote and nothing to deprecate, and the register is accurate as written. `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked this run) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision ("bumped only in the same change that completes that surface, after reconciling every row tagged with the level"). The milestone-bump issue the recipe asks me to file **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": v3.0's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is delivered, the 1.4 half is open (#92, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut.
.mule/reports/log/20261003T180922Z-recipe-github-issues-proposed.md:429:   **The recipe's standing item is satisfied and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental*, reconcile trigger "upstream 1.4 final". 1.4 is not final, so there is nothing to promote and nothing to deprecate, and the register is accurate as written. `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked this run) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision ("bumped only in the same change that completes that surface, after reconciling every row tagged with the level"). The milestone-bump issue the recipe asks me to file **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": v3.0's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is delivered, the 1.4 half is open (#92, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut.
.mule/reports/log/20261007T191848Z-recipe-docs-sync-proposed.md:698:94:  **The recipe's standing item is satisfied and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental*, reconcile trigger "upstream 1.4 final". 1.4 is not final, so there is nothing to promote and nothing to deprecate, and the register is accurate as written. `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked this run) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision ("bumped only in the same change that completes that surface, after reconciling every row tagged with the level"). The milestone-bump issue the recipe asks me to file **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": v3.0's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is delivered, the 1.4 half is open (#92, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut.
.mule/reports/log/20261007T191848Z-recipe-docs-sync-proposed.md:701:151:  **The recipe's standing item is satisfied, and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still has exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental* with reconcile trigger "upstream 1.4 final". v1.4 is not final, so both rows stay as written: nothing to promote, nothing to deprecate, and the register is accurate. `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision; the milestone-bump issue the recipe asks for **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": the milestone's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is done, and the 1.4 half is unresolved (#92 open, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut. Also still true and already escalated on 2026-09-06, so not repeated: the section's source document `.mule/research/upstream-parity-2026-08-22.md` does not exist in the repo.
.mule/reports/log/20261005T193452Z-recipe-milestones-proposed.md:337:+  **Step 5 deliberately not taken.** I am *not* proposing "milestone v3.0 looks complete, verify and cut the tag", for the same reason as the twenty runs before this one: the section's declared scope is **1.2.2 → 1.3/1.4** and only the 1.3 half is delivered (now including v1.3.5, which is inert anyway). The 1.4 half is open — #92 unanswered, both experimental rows unreconciled, `APICompatVersion` not bumped — so the tag must not be cut. Your call when a stable v1.4.0 ships: answer #92, settle the register rows, run the final-phase bump, then cut the tag.
.mule/reports/log/20261005T193452Z-recipe-milestones-proposed.md:363:  **Step 5 deliberately not taken.** I am *not* proposing "milestone v3.0 looks complete, verify and cut the tag", for the same reason as the twenty runs before

I've completed the recipe. Here's what I found.

**Target:** v2.0 is `DONE`; v3.0 (upstream parity 1.2.2 → 1.3/1.4) is the current milestone, and its scope is decided (not TBD), so step 2 doesn't apply — I investigated per step 3.

**Findings (measured this run):**
- `gh issue list --label milestone-3.0 --state all` → empty (this milestone's ledger lives under `upstream-parity`).
- Under `upstream-parity`: **#47–#89 all closed**, and all 43 closed issues are `state_reason=COMPLETED` (none `not_planned`, so no dropped gap). Open: **#92** keyAgreement (parked decision), **#93** aclhook comment (`mule-review`, its own path), **#78** FDO (`milestone-4.0`). #90 closed (frozen "don't bump now").
- Release sweep: **no stable v1.4.0** — newest stable `v1.3.5` (2026-10-05), newest overall `v1.4.0-rc.6` (prerelease). #92's parking condition unmet. Upstream `master` newest commit is still `171c1bd7` (2026-10-06), so nothing new since the last run.
- Standing item: `docs/UPSTREAM-EXPERIMENTAL.md` still has exactly its two rows (#67, #68), both gated on "upstream 1.4 final" — not final, so nothing to promote/deprecate. `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`).

**Buckets:** nothing is machine-checkable-and-unfiled; every remaining item is already an open issue (#92/#78/#93) or a closed decided one (#90), so filing would duplicate. No new design-decision gaps.

**Step 5 not taken**, consistent with the 21 prior runs: the milestone's declared scope is 1.2.2 → 1.3/**1.4**, the 1.4 register rows are unreconciled and #92 is parked — so "looks complete, cut the tag" would be false. The `.mule/research/upstream-parity-2026-08-22.md` source doc is absent (gitignored `research/`), already escalated; I used the #47–#89 backlog plus upstream releases as the reference.

```
Done: milestones recipe run for v3.0 — investigated, found no new gaps; filed no issues and appended no .mule/todo.md lines (all candidates already open, parked, or decided)
Files: none (working tree unchanged)
Verified: gh issue list --label milestone-3.0/upstream-parity/milestone-4.0 --state all; gh api releases/tags/commits; grep 'APICompatVersion =' internal/realm/service.go; read docs/UPSTREAM-EXPERIMENTAL.md -> pass, nothing new / nothing filed
Unsure: source doc .mule/research/upstream-parity-2026-08-22.md absent (gitignored), so the 1.3/1.4 surface came from the #47–#89 backlog rather than re-derived from upstream
Follow-ups: none
```
