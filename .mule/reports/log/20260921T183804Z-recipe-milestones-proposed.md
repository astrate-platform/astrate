slug: recipe-milestones
verdict: proposed
at:  8311645
ran: 2026-09-21T18:38:04Z on DietPi in 342s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
The first non-DONE milestone is **v3.0** (upstream parity). Its scope is decided, so I go to step 3: investigate. Let me check the reference research doc, the standing EXPERIMENTAL reconciliation file, and what's already filed.

$ wc -l .mule/research/upstream-parity-2026-08-22.md docs/UPSTREAM-EXPERIMENTAL.md .mule/recipes/astarte-upstream.md 2>&1; ls .mule/research/ 2>&1
wc: .mule/research/upstream-parity-2026-08-22.md: No such file or directory
  25 docs/UPSTREAM-EXPERIMENTAL.md
  62 .mule/recipes/astarte-upstream.md
  87 total
ls: cannot access '.mule/research/': No such file or directory

$ gh issue list --label upstream-parity --state all --limit 60 2>&1; echo ---; gh issue list --label milestone-3.0 --state all --limit 60 2>&1
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
---

$ ls -la .mule/ 2>&1; echo ---; ls .mule/research 2>&1; echo ---; git log --oneline -15 2>&1
total 5880
drwxr-xr-x  7 root root    4096 Sep 21 20:20 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      13 Sep 21 20:32 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 21 19:35 .heartbeat
-rw-r--r--  1 root root    4096 Sep 21 20:33 .last-output
-rw-r--r--  1 root root       0 Sep 21 20:32 .lock
-rw-r--r--  1 root root      11 Sep 21 20:32 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 5791352 Sep 21 20:33 cron.log
drwxr-xr-x  2 root root    4096 Sep 19 12:52 failed
-rw-r--r--  1 root root   50415 Sep 20 22:42 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    8227 Sep 21 19:35 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 21 13:32 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   71911 Sep 21 19:35 todo.md
---
ls: cannot access '.mule/research': No such file or directory
---
8311645 mule: recipe hygiene timed out
80d2aaf mule: recipe docs-sync ran (328s)
b200577 mule: log store-pipelines-empty-name-zero-blocks-test
82fc0a4 mule: store-pipelines-empty-name-zero-blocks-test [auto]: extend internal/store/pipelines_validate_test.go with the two untested error branches of `validatePipelineGraph` (internal/store/pipelines.go:51-64) — a definition with zero blocks and one whose block has an empty name must both be rejected. Pure helper, no DB; the existing suite covers only acyclic/cyclic/unknown-ref/duplicate-name.
d4baf48 mule: refill the queue
5bf6606 mule: recipe code-review ran (480s)
7568e53 mule: recipe astarte-upstream ran (110s)
7aa7f65 mule: recipe astarte-upstream ran (99s)
7c9809b mule: refill the queue
b82e0ec mule: recipe github-issues ran (186s)
2e34867 mule: refill the queue
fc5f492 mule: recipe milestones ran (193s)
9195e7f mule: recipe hygiene timed out
395d8bb mule: log docs-sync-native-socket-missing-403-500
f9d8b27 mule: docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

→ Read docs/UPSTREAM-EXPERIMENTAL.md
The research doc path doesn't exist here, and a large set of `upstream-parity` issues are already filed (mostly closed). Let me check the `.mule` tree and the standing files before deciding anything.
→ Read .mule/recipes/astarte-upstream.md
The research doc referenced by the v3.0 milestone (`.mule/research/upstream-parity-2026-08-22.md`) is missing here, yet a large set of `upstream-parity` issues already exist. Let me dig into the actual state before deciding anything.

$ git remote -v; echo ---; git branch -a 2>&1 | head -30; echo ---; find . -name '*upstream-parity*' -not -path './.git/*' 2>/dev/null; find . -path '*research*' -not -path './.git/*' 2>/dev/null | head
origin	git@github.com:astrate-platform/astrate.git (fetch)
origin	git@github.com:astrate-platform/astrate.git (push)
---
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

→ Read docs/COMPATIBILITY.md
The repo is the astrate repo itself, and `.mule/research/` doesn't exist locally. Let me check the `mule/research` branch and current state.

$ git show origin/mule/research:.mule/research/upstream-parity-2026-08-22.md 2>&1 | head -5; echo "=== branch file listing ==="; git ls-tree -r origin/mule/research --name-only 2>&1 | rg -i 'research|parity' | head -20
fatal: path '.mule/research/upstream-parity-2026-08-22.md' does not exist in 'origin/mule/research'
=== branch file listing ===
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
.mule/research/issues-2026-09-04.md
.mule/research/log.md
.mule/research/survey-2026-07-27.md

$ git status --short 2>&1 | head; echo ---; git log --oneline --all --grep='milestones' -10 2>&1
---
fc5f492 mule: recipe milestones ran (193s)
f166917 mule: recipe milestones ran (124s)
b3c4bad mule: recipe milestones ran (176s)
8190c4e mule: recipe milestones ran (379s)
9d75bee mule: recipe milestones timed out
ed19432 mule: recipe milestones ran (164s)
3bdfa4a mule: recipe milestones ran (116s)
dcb1a41 mule: recipe milestones ran (95s)
36ae034 mule: recipe milestones ran (498s)
8ed093d mule: recipe milestones ran (95s)

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:6][] | "\(.tag_name)\t\(.published_at)"' 2>&1
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

The research doc path differs. Let me search all branches for it and check recent log/handoffs.
$ for b in main mule/research mule/queue; do echo "== $b =="; git ls-tree -r origin/$b --name-only 2>/dev/null | rg -i 'parity-2026|upstream-parity' ; done
== main ==
== mule/research ==
== mule/queue ==

→ Read .mule/log.md
The v3.0 reference doc path is missing, but the milestone's actual backlog is the `upstream-parity` GitHub issues (mostly closed now). Let me check the state of open issues, the for-giulio file, and what the last milestones run produced.

$ git show fc5f492 --stat 2>&1 | head -30; echo ===; git show fc5f492 -- .mule/todo.md .mule/for-giulio.md 2>&1 | head -120
commit fc5f4927a07332af5c6123702d65ca91d25a0228
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Sun Sep 20 22:30:30 2026 +0200

    mule: recipe milestones ran (193s)

 .../20260920T203030Z-recipe-milestones-proposed.md | 1196 ++++++++++++++++++++
 1 file changed, 1196 insertions(+)
===

$ gh issue list --label upstream-parity --state open 2>&1; echo "=== all with milestone-3.0/4.0 ==="; gh issue list --label milestone-4.0 --state all --limit 30 2>&1; echo "=== issue 92,93 ==="; gh issue view 92 2>&1 | head -40
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
=== all with milestone-3.0/4.0 ===
90	CLOSED	milestone-4.0: decide on bumping APICompatVersion from "1.2.2" to "1.3.0-rc.0" (or making it configurable)	enhancement, milestone-4.0	2026-08-23T21:06:34Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
28	CLOSED	milestone-4.0: investigate Edgehog compatibility with Astrate (readonly)	mule, readonly, milestone-4.0	2026-08-23T15:48:58Z
=== issue 92,93 ===
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

$ ls -la .mule/reports/ .mule/reports/* 2>&1 | tail -20; echo ===; ls .mule/reports/20260929* 2>&1; ls .mule/reports | head
-rw-r--r-- 1 root root   3612 Sep 19 19:08 20260919T170808Z-recipe-hygiene-timeout.md
-rw-r--r-- 1 root root  18560 Sep 19 23:14 20260919T211455Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  13381 Sep 20 13:10 20260920T111057Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   6985 Sep 20 13:36 20260920T113611Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  35633 Sep 20 19:31 20260920T173118Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root   3454 Sep 20 20:11 20260920T181128Z-errorname-missing-required-test-done.md
-rw-r--r-- 1 root root  22144 Sep 20 20:49 20260920T184908Z-triggers-custom-action-policy-nodecide-done.md
-rw-r--r-- 1 root root   9688 Sep 20 20:55 20260920T185539Z-compat-note-custom-action-policy-boundary-done.md
-rw-r--r-- 1 root root  30206 Sep 20 20:59 20260920T185956Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root   3139 Sep 20 21:16 20260920T191642Z-docs-sync-native-metrics-example-fake-series-done.md
-rw-r--r-- 1 root root   7254 Sep 20 21:27 20260920T192702Z-docs-sync-native-socket-missing-403-500-done.md
-rw-r--r-- 1 root root   6408 Sep 20 21:54 20260920T195443Z-recipe-hygiene-timeout.md
-rw-r--r-- 1 root root 134480 Sep 20 22:30 20260920T203030Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  22583 Sep 20 22:42 20260920T204223Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root    218 Sep 21 04:06 20260921T020624Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   8924 Sep 21 12:57 20260921T105709Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  21111 Sep 21 13:33 20260921T113318Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root   2996 Sep 21 19:35 20260921T173503Z-store-pipelines-empty-name-zero-blocks-test-done.md
-rw-r--r-- 1 root root  29006 Sep 21 19:46 20260921T174646Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root   2624 Sep 21 20:20 20260921T182029Z-recipe-hygiene-timeout.md
===
ls: cannot access '.mule/reports/20260929*': No such file or directory
hygiene.md
log

✱ Grep "milestone|v3\.0|upstream-parity|#9[0-3]|keyAgreement" in .mule/for-giulio.md · 100 matches

$ git log -1 --format='%ci %s' fc5f492; echo ===; git status --short | head
2026-09-20 22:30:30 +0200 mule: recipe milestones ran (193s)
===

Let me read the latest milestone recipe report to see what the previous run concluded.
$ tail -60 .mule/reports/20260920T203030Z-recipe-milestones-proposed.md
tail: cannot open '.mule/reports/20260920T203030Z-recipe-milestones-proposed.md' for reading: No such file or directory

$ ls .mule/reports/ 2>/dev/null | rg milestones; echo ===; ls reports 2>/dev/null | rg milestones; echo ===; find /root/astrate-mule/.mule -name '*milestones-proposed*' -newermt '2026-09-20' 2>/dev/null
===
===
/root/astrate-mule/.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md

$ ls -la .mule/reports/log/ | rg 'milestones|github-issues|astarte-upstream' ; echo ===; wc -l .mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root 120837 Sep  5 20:10 20260905T181053Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  24739 Sep  5 20:12 20260905T181219Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root 152510 Sep  5 20:46 20260905T184625Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  63922 Sep  6 20:04 20260906T180442Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root   8576 Sep  6 20:13 20260906T181352Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  57962 Sep  6 20:31 20260906T183102Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root    218 Sep  7 04:03 20260907T020359Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   8038 Sep  7 19:19 20260907T171928Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  49686 Sep  7 19:48 20260907T174812Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  82600 Sep  7 19:55 20260907T175524Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  11829 Sep  7 22:47 20260907T204751Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  49792 Sep  8 13:20 20260908T112032Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  85176 Sep  8 13:34 20260908T113422Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   4886 Sep  8 21:31 20260908T193026Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  49296 Sep  8 21:31 20260908T193127Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  57250 Sep  8 21:33 20260908T193321Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   2013 Sep  8 21:38 20260908T193614Z-recipe-astarte-upstream-noop.md
-rw-r--r-- 1 root root  28135 Sep  8 21:38 20260908T193814Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   4409 Sep  9 19:57 20260909T175718Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root    215 Sep  9 20:13 20260909T181348Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  27766 Sep  9 20:43 20260909T184315Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  35463 Sep 10 01:04 20260909T230411Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  28297 Sep 10 13:09 20260910T110956Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   1602 Sep 10 13:25 20260910T112515Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  17437 Sep 10 13:25 20260910T112545Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  17483 Sep 10 21:53 20260910T195304Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  67965 Sep 10 22:01 20260910T200120Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   2160 Sep 10 22:29 20260910T202936Z-recipe-astarte-upstream-noop.md
-rw-r--r-- 1 root root  17656 Sep 10 22:29 20260910T202949Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  22808 Sep 11 20:53 20260911T185303Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  14010 Sep 11 21:04 20260911T190417Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  19257 Sep 11 21:15 20260911T191541Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  22280 Sep 12 19:12 20260912T171256Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  24726 Sep 12 19:29 20260912T172901Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  45876 Sep 12 20:07 20260912T180735Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root    212 Sep 13 18:52 20260913T165246Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  68245 Sep 13 19:05 20260913T170550Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   1971 Sep 13 19:13 20260913T171327Z-recipe-astarte-upstream-noop.md
-rw-r--r-- 1 root root  25213 Sep 13 19:13 20260913T171351Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  45997 Sep 13 22:32 20260913T203209Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root   2991 Sep 13 22:38 20260913T203800Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  54978 Sep 13 22:38 20260913T203830Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   3313 Sep 13 22:44 20260913T204415Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   2038 Sep 14 04:06 20260914T020552Z-recipe-astarte-upstream-noop.md
-rw-r--r-- 1 root root   3697 Sep 14 04:06 20260914T020608Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  25311 Sep 15 01:45 20260914T234550Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  26800 Sep 15 12:57 20260915T105740Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   2011 Sep 15 13:13 20260915T111354Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  14548 Sep 15 20:44 20260915T184406Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  78833 Sep 15 20:58 20260915T185804Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  73652 Sep 15 21:36 20260915T193626Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  27197 Sep 16 19:58 20260916T175801Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  50143 Sep 16 20:39 20260916T183922Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   2049 Sep 16 20:42 20260916T184159Z-recipe-astarte-upstream-noop.md
-rw-r--r-- 1 root root  32721 Sep 16 20:42 20260916T184218Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  97405 Sep 17 20:43 20260917T184319Z-recipe-milestones-timeout.md
-rw-r--r-- 1 root root  74737 Sep 17 21:03 20260917T190323Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   2026 Sep 17 21:06 20260917T190650Z-recipe-astarte-upstream-noop.md
-rw-r--r-- 1 root root  47648 Sep 17 21:07 20260917T190703Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  40531 Sep 18 13:00 20260918T110039Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  11055 Sep 18 13:34 20260918T113442Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   3271 Sep 18 13:40 20260918T113930Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  64748 Sep 18 13:41 20260918T114100Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  14495 Sep 18 21:36 20260918T193641Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  11462 Sep 18 21:42 20260918T194216Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  10735 Sep 18 21:47 20260918T194746Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  18560 Sep 19 23:14 20260919T211455Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  13381 Sep 20 13:10 20260920T111057Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   6985 Sep 20 13:36 20260920T113611Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root 134480 Sep 20 22:30 20260920T203030Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  22583 Sep 20 22:42 20260920T204223Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root    218 Sep 21 04:06 20260921T020624Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   8924 Sep 21 12:57 20260921T105709Z-recipe-astarte-upstream-proposed.md
===
1196 .mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md

→ Read .mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md [offset=1080, limit=116]
Now I have the picture. Let me run a fresh verification for today's run (the last milestones run was 2026-09-20; this repo has an established no-op pattern for v3.0, gated on upstream v1.4 final).

$ echo "=== releases ==="; gh api repos/astarte-platform/astarte/releases --jq '.[] | "\(.tag_name) \(.prerelease) \(.published_at)"' 2>&1 | head -8; echo "=== milestone-3.0 label ==="; gh issue list --label milestone-3.0 --state all --limit 20 2>&1; echo "=== upstream-parity open ==="; gh issue list --label upstream-parity --state open --json number,title --jq '.[] | "\(.number)\t\(.title)"' 2>&1
=== releases ===
v1.3.4 false 2026-09-18T09:42:57Z
v1.4.0-rc.5 true 2026-08-20T12:28:04Z
v1.3.3 false 2026-08-07T14:33:58Z
v1.4.0-rc.4 true 2026-08-06T13:40:20Z
v1.4.0-rc.3 true 2026-07-31T10:48:39Z
v1.4.0-rc.2 true 2026-07-14T10:16:22Z
v1.4.0-rc.1 true 2026-07-14T08:53:27Z
v1.3.2 false 2026-07-14T10:07:26Z
=== milestone-3.0 label ===
=== upstream-parity open ===
93	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
92	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds
78	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)

$ echo "=== upstream-parity ALL (count + newest numbers) ==="; gh issue list --label upstream-parity --state all --limit 100 --json number,state,createdAt --jq '.[] | "\(.number)\t\(.state)\t\(.createdAt)"' 2>&1 | sort -t $'\t' -k1 -n | tail -8; echo "=== total ==="; gh issue list --label upstream-parity --state all --limit 100 2>&1 | wc -l
=== upstream-parity ALL (count + newest numbers) ===
sort: multi-character tab ‘$\\t’
=== total ===
46

The state is unchanged since the 2026-09-20 run. Let me confirm the top of `for-giulio.md` to match the established format for the verification entry.
→ Read .mule/for-giulio.md [limit=30]
State confirmed unchanged since the 2026-09-20 run: no stable v1.4.0 (newest stable `v1.3.4`, newest overall `v1.4.0-rc.5`), open `upstream-parity` set is exactly #92/#93/#78, `milestone-3.0` label empty, both EXPERIMENTAL register rows still gated on v1.4 final. Per the recipe's bucket rule all open items are already escalated — no new gap, so nothing to file and no todo lines. Recording the verification (newest at top) as the deliverable.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+
+---
+
 - **github-issues triage re-run, 2026-09-20 (evening): nothing proposable, no change since this morning's 13:10 run.** Re-surveyed with the recipe command: the issue set is identical to the 2026-09-20 morning triage — **#94–#107** are still the daily `mule-alarm` noise (14 issues, zero comments, telemetry not code), **#93** aclhook comment rewrite is still mule-review with `8c61268` already pushed (its own recipe path), **#92** keyAgreement is still parked on a stable upstream v1.4.0 with its waiting-on.md row and the mule-upstream-watch escalation wired (an -rc does not satisfy it), **#78** FDO is still the milestone-4.0 design/investigation already escalated, and **#1** is untouched per standing instruction. No new issue has been filed since the morning run. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines proposed**. This morning's proposal stands unchanged: **close #94–#106** — each was a one-day event, superseded by the next day's alarm, never actionable — and leave **#107** (still the newest, no #108 yet) to self-expire. The queue idle streak since ~2026-09-04/05 persists; if that is not intentional the wedge is on the Pi (queue review = the fix, per the 13:10 line).
 
 ---
 

→ Read .mule/for-giulio.md [limit=8, offset=11]
Recipe complete. The active milestone is **v3.0** (first non-DONE). Its scope is decided (2026-08-22 parity investigation), so I investigated the gap:

- **Upstream gate:** no stable v1.4.0 yet — newest stable `v1.3.4` (2026-09-18, maintenance-only), newest overall `v1.4.0-rc.5` (2026-08-20). Nothing shipped since the 2026-09-20 run.
- **Already filed (`upstream-parity` ledger, since `milestone-3.0` label is unused):** open items are exactly #92 (keyAgreement decision, parked on stable v1.4.0), #93 (aclhook comment, mule-review with `8c61268` pushed, its own recipe path), #78 (FDO, milestone-4.0). The #47–#89 rest are closed.
- **Standing item:** `docs/UPSTREAM-EXPERIMENTAL.md` rows #67/#68 still wait on v1.4 final; `APICompatVersion` stays 1.2.2 per #90's frozen decision.
- **Bucket:** every open item is already escalated → no machine-checkable gap; nothing to file, no todo lines. Deliberately not proposing "cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved.

The reference doc path (`.mule/research/upstream-parity-2026-08-22.md`) is absent, but the milestone's actual source of record is the `upstream-parity` issue backlog, which is reachable and consistent — no need to stop.

```
Done: v3.0 verified unchanged — still gated on upstream v1.4 final; recorded in for-giulio.md
Files: .mule/for-giulio.md
Verified: gh api releases sweep + gh issue list (upstream-parity / milestone-3.0) + docs/UPSTREAM-EXPERIMENTAL.md read -> pass (state unchanged)
Unsure: nothing
Follow-ups: none
```
