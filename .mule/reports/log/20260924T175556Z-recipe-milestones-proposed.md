slug: recipe-milestones
verdict: proposed
at:  94f7284
ran: 2026-09-24T17:55:56Z on DietPi in 94s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
The current target is **v3.0 — upstream Astarte feature parity**. Scope is decided (not TBD), so I investigate and triage. Let me gather the reference material and existing issue state.
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md

$ gh label list --limit 100
bug	Something isn't working	#d73a4a
documentation	Improvements or additions to documentation	#0075ca
duplicate	This issue or pull request already exists	#cfd3d7
enhancement	New feature or request	#a2eeef
good first issue	Good for newcomers	#7057ff
help wanted	Extra attention is needed	#008672
invalid	This doesn't seem right	#e4e669
question	Further information is requested	#d876e3
wontfix	This will not be worked on	#ffffff
mule	Queued for the autonomous mule to pick up	#5319e7
mule-alarm	Filed by the mule's dead-man switch: it has gone quiet	#b60205
mule-review	Mule pushed a change for this; waiting for a human or a strong model to review it	#0e8a16
mule-blocked	The mule tried this and could not do it; see its comment	#d93f0b
legion	Needs the Legion Go (Docker, race detector, benchmarks)	#fbca04
readonly	A standing check: verifies something, changes nothing, never closes	#c5def5
milestone-2.0	astarte-flow feature parity milestone	#0E8A16
milestone-4.0		#5319e7
milestone-3.0		#5319e7
upstream-parity	Feature parity with upstream astarte-platform	#0052cc
upstream-experimental	Adopted from an upstream experimental/unstable feature — review each upstream release	#d4c5f9

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

✱ Glob ".mule/research/*.md" 0 matches
The research doc isn't at that path. Let me find it, plus list the milestone-3.0 labeled open issues.

$ gh issue list --label milestone-3.0 --state all --limit 100
(no output)

$ ls -la .mule/ && echo "---" && ls -la .mule/research 2>/dev/null || echo "no research dir"
total 6900
drwxr-xr-x  7 root root    4096 Sep 24 18:56 .
drwxr-xr-x 17 root root    4096 Sep 23 19:04 ..
-rw-r--r--  1 root root      13 Sep 24 19:54 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 24 18:56 .heartbeat
-rw-r--r--  1 root root    8192 Sep 24 19:55 .last-output
-rw-r--r--  1 root root       0 Sep 24 19:54 .lock
-rw-r--r--  1 root root      11 Sep 24 19:54 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 6812154 Sep 24 19:55 cron.log
drwxr-xr-x  2 root root    4096 Sep 23 21:33 failed
-rw-r--r--  1 root root   56560 Sep 23 19:47 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    9307 Sep 24 18:56 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 23 20:39 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   82583 Sep 24 18:56 todo.md
---
no research dir

✱ Glob "**/upstream-parity*" 0 matches

$ git remote -v; echo "---"; ls /root/
origin	git@github.com:astrate-platform/astrate.git (fetch)
origin	git@github.com:astrate-platform/astrate.git (push)
---
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

The research doc referenced by the milestone doesn't exist in this repo. Let me check whether it lives elsewhere, including the astrate repo.
$ git log --oneline -15
94f7284 mule: recipe hygiene ran (583s)
fab3c51 mule: log docs-sync-rm-auth-403
9c05370 mule: docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
5d57685 mule: log docs-sync-rm-validation-example-prefix
02e269b mule: docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "realm: validation failed: interface definition is invalid"`, but `validationDetail` strips the `realm: validation failed: ` prefix on the wire (http.go:562-568) and the parser's real messages are `invalid interface: ...` (pkg/interfaceschema/parse.go:38, wrapped at service.go:163); the example should show the stripped, real message. Same class as the already-fixed docs-sync-hk-validation-example (housekeeping http.go:251-258). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
3df9563 mule: log docs-sync-rm-interface-422-shapes
93e784e mule: docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` (yaml:231-232) in docs/api/astarte_realm_management_api.yaml as a oneOf — both handlers go through `writeInterfaceError`, which answers three distinct 422 bodies: the flat ErrorDetail (`ErrValidation` via `validationDetail`, http.go:392-393), the nested violations changeset envelope (`writeViolations`, http.go:358-359 + 434-533, e.g. `{"errors":{"description":["should be at most 1000 character(s)"]}}` and the aligned full-length `mappings` array, http_test.go:480-501), and the named FieldErrors envelope for `ErrMaximumDatabaseRetentionExceeded` (`{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, http.go:355-357, http_test.go:282-320); today only the flat ValidationError is referenced. Mirror the createTrigger oneOf pattern (yaml:344-368). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
d4c11ca mule: log docs-sync-rm-put-interface-409
7256af9 mule: docs-sync-rm-put-interface-409 [auto]: add the missing `409` Conflict response to `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` in docs/api/astarte_realm_management_api.yaml (yaml:222-234 lists only 204/400/401/404/422/500) — `updateInterface` answers 409 for every url/body disagreement and minor-upgrade incompatibility: `ErrNameMismatch` (service.go:204 → "Interface name doesn't match the one in the interface json", http.go:367-369), `ErrMajorMismatch` (service.go:207 → http.go:370-372), and the `CheckMinorUpgrade` sentinels `ErrMinorNotIncreased`/`ErrDowngradeNotAllowed`/`ErrMissingEndpoints`/`ErrIncompatibleEndpointChange` (service.go:223-226 → http.go:373-384) plus `store.ErrAlreadyExists` (http.go:385-386) — the exact set `TestRealmManagementErrorCodes` pins as StatusConflict (http_test.go:417-437). Reuse the `Conflict` response component (already used on POST, yaml:131-132) and add name-mismatch/major-mismatch examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
63eae06 mule: refill the queue
50da4bf mule: recipe docs-sync ran (779s)
55bba71 mule: log flow-filter-key-contains-test
3875ba6 mule: flow-filter-key-contains-test [auto]: `Filter`'s `key_contains` rule (`strings.Contains`, internal/flow/blocks/transform.go:69) has zero test coverage — `transform_test.go` only exercises `key_prefix`, `type`, unknown-type and the at-least-one construct rule, so a regression silently dropping the Contains check would pass the whole suite. Add substring-match, substring-absent, and prefix+contains combined rows to the filter tests. Container-free.
83d0124 mule: blocked flow-randomsource-span-overflow
2e1c216 mule: log flow-sort-bounded-buffer

$ find /root/astrate -name "*parity*" 2>/dev/null; echo "---"; ls /root/astrate/.mule 2>/dev/null || echo "no .mule in astrate"; echo "---"; ls /root/astrate/docs 2>/dev/null | head -30
---
MULE.md
config
for-giulio.md
hosts.example
log.md
recipes
tasks
todo.md
---
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

! permission requested: external_directory (/root/astrate/.mule/recipes/*); auto-rejecting
✗ git log --all --oneline -- '.mule/research/' | head; echo "---"; ls /root/astrate/.mule/tasks 2>/dev/null; echo "---"; cat /root/astrate/.mule/recipes/* 2>/dev/null | head -5; ls /root/astrate/.mule/recipes 2>/dev/null failed
Error: The user rejected permission to use this specific tool call.
