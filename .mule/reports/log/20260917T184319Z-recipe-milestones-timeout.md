slug: recipe-milestones
verdict: timeout
at:  d9a534b
ran: 2026-09-17T18:43:19Z on DietPi in 1200s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
→ Read .mule/todo.md
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md

$ gh auth status 2>&1; echo "---"; gh issue list --label milestone-3.0 --state all --limit 50 2>&1
github.com
  ✓ Logged in to github.com account atsetilam (/root/.config/gh/hosts.yml)
  - Active account: true
  - Git operations protocol: https
  - Token: gho_************************************
  - Token scopes: 'gist', 'read:org', 'repo'
---

✱ Glob ".mule/**/*upstream-parity*" 0 matches

$ gh issue list --label upstream-parity --state all --limit 80 2>&1 | head -100
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

$ find . -iname '*upstream-parity*' -not -path './.git/*' 2>/dev/null; echo "---research dir---"; ls -la .mule/research/ 2>&1; echo "---for-giulio---"; ls -la .mule/ 2>&1
---research dir---
ls: cannot access '.mule/research/': No such file or directory
---for-giulio---
total 4668
drwxr-xr-x  7 root root    4096 Sep 17 20:22 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      13 Sep 17 20:23 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 17 13:45 .heartbeat
-rw-r--r--  1 root root    4096 Sep 17 20:24 .last-output
-rw-r--r--  1 root root       0 Sep 17 20:23 .lock
-rw-r--r--  1 root root      11 Sep 17 20:23 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 4582873 Sep 17 20:24 cron.log
drwxr-xr-x  2 root root    4096 Sep 15 19:13 failed
-rw-r--r--  1 root root   36409 Sep 16 20:39 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    7213 Sep 17 13:45 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 16 20:48 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   61140 Sep 17 13:45 todo.md

$ git branch -a 2>&1 | head; echo "---"; git log --oneline -5 2>&1; echo "---"; git log --all --oneline -- '.mule/research/upstream-parity-2026-08-22.md' 2>&1 | head
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
---
d9a534b mule: recipe hygiene timed out
e33ab46 mule: log docs-sync-native-version-value
08048b3 mule: docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
30665a5 mule: log docs-sync-pairing-version-value
c69a164 mule: docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
---

$ git log --all --oneline --name-only 2>/dev/null | rg -i 'upstream-parity|research' | head -20; echo "---branches with research---"; git ls-tree -r origin/mule/research --name-only 2>/dev/null | rg -i 'upstream-parity|research' | head -20
203af1c mule: hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto]
015a706 mule: hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto]
c867b2c mule: docs-sync-pairing-version-endpoint [auto]: add the undocumented `GET /pairing/v1/{realm}/version` route to docs/api/astarte_pairing_api.yaml — it exists in code (cmd/astrate/main.go:398, unauthenticated, upstream-parity per main.go:387-391) and answers 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38); mirror the RM spec's getVersion op (astarte_realm_management_api.yaml:532) without its a_rma security. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
2f6176c mule: docs-sync-pairing-status-enum [auto]: fix the `PairingInfo.status` enum in docs/api/astarte_pairing_api.yaml from `[confirmed, pending, denied, expired]` (yaml:412) to `[confirmed, pending, inhibited]` — the handler never emits `denied` or `expired`: `service.Info` maps only `DeviceStatusInhibited`→"inhibited", `DeviceStatusConfirmed`→"confirmed", and everything else→"pending" (internal/pairing/service.go:296-299, statuses defined in store/devices.go:19-27; upstream-parity per the comment at service.go:282-285). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/research/triaged.md
.mule/research/issues-2026-09-04.md
.mule/research/survey-2026-09-04.md
.mule/research/log.md
.mule/research/log.md
.mule/research/survey-2026-08-31.md
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
---branches with research---
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

$ git show origin/mule/research:.mule/research/upstream-parity-2026-08-22.md > /tmp/upstream-parity-2026-08-22.md 2>&1 && echo FOUND branch || echo "not on research branch"; ls -la .mule/research.md 2>/dev/null; find . -name 'upstream-parity*' -not -path './.git/*' 2>/dev/null
not on research branch

$ git ls-tree -r origin/mule/research --name-only 2>/dev/null | rg -i 'upstream-parity|research' ; echo "---"; git log --all --oneline -- '.mule/research/upstream-parity-2026-08-22.md' 2>/dev/null
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
.mule/research/survey-2026-08-31.md
.mule/research/survey-2026-09-04.md
.mule/research/triaged.md
---

$ git log --all --oneline -- '.mule/research/*upstream*' | head; echo "---all branches containing upstream-parity file---"; for b in $(git branch -a | tr -d ' *'); do if git ls-tree -r $b --name-only 2>/dev/null | rg -q 'upstream-parity-2026-08-22'; then echo "in $b"; git ls-tree -r $b --name-only | rg 'upstream-parity-2026-08-22'; fi; done
---all branches containing upstream-parity file---

✗ Read .mule/research/issues-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/issues-2026-08-22.md

$ cat .mule/research/log.md | rg -i 'pare|1\.3|1\.4|retention|alias' | tail -40
cat: .mule/research/log.md: No such file or directory

$ ls -la .mule/research/ | head -40; echo "===for-giulio==="; ls -la .mule/for-giulio* 2>/dev/null
ls: cannot access '.mule/research/': No such file or directory
===for-giulio===
-rw-r--r-- 1 root root 36409 Sep 16 20:39 .mule/for-giulio.md

$ git ls-tree -r origin/mule/research --name-only | rg -i 'upstream-parity'; echo "---whoami/what repo---"; git remote -v; git config user.name
---whoami/what repo---
origin	git@github.com:astrate-platform/astrate.git (fetch)
origin	git@github.com:astrate-platform/astrate.git (push)
Giulio Malitesta

$ git ls-tree -r origin/mule/research --name-only | head -40
.dockerignore
.github/workflows/ci.yml
.github/workflows/docs.yml
.gitignore
.golangci.yml
.mule/.gitignore
.mule/MULE.md
.mule/config
.mule/for-giulio.md
.mule/hosts.example
.mule/log.md
.mule/recipes/astarte-upstream.md
.mule/recipes/benchmarks.md
.mule/recipes/code-review.md
.mule/recipes/cross-project-survey.md
.mule/recipes/docs-sync.md
.mule/recipes/github-issues.md
.mule/recipes/hygiene.md
.mule/recipes/legion-go.md
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
.mule/research/survey-2026-07-28.md

$ gh issue list --label upstream-parity --state all --limit 200 2>&1 | wc -l; echo "===open upstream-parity==="; gh issue list --label upstream-parity --state open --limit 100 2>&1; echo "===milestone-3.0 label==="; gh issue list --label milestone-3.0 --state all --limit 100 2>&1
46
===open upstream-parity===
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
===milestone-3.0 label===

$ gh issue list --label upstream-parity --state all --limit 100 2>&1 | sort
47	CLOSED	MQTT v1: accept device capabilities incl. purge_properties_compression_format (plaintext purge)	enhancement, upstream-parity	2026-08-22T04:04:09Z
48	CLOSED	MQTT v1: tolerate devices connecting with an empty introspection	enhancement, upstream-parity	2026-08-22T04:04:11Z
49	CLOSED	MQTT v1: binaryblob properties re-send as BSON subtype 0 in consumer/properties	enhancement, upstream-parity	2026-08-22T04:04:23Z
50	CLOSED	Broker/engine hardening: discard malformed introspection/purge/capabilities instead of crashing	enhancement, upstream-parity	2026-08-22T04:04:24Z
51	CLOSED	keyAgreement control topic + encrypted mappings (upstream 1.4 experimental) — investigate/park	enhancement, upstream-parity	2026-08-22T08:48:28Z
52	CLOSED	AppEngine: devices-by-alias full mirror (PATCH via alias + interface data access)	enhancement, upstream-parity	2026-08-22T09:18:20Z
53	CLOSED	AppEngine: group-scoped device PATCH + interface data access under groups/{g}/devices/{d}	enhancement, upstream-parity	2026-08-22T09:18:23Z
54	CLOSED	AppEngine: GET /groups/{g} show-group endpoint	enhancement, upstream-parity	2026-08-22T10:35:16Z
55	CLOSED	AppEngine: complete data-query parameter set	enhancement, upstream-parity	2026-08-23T15:09:33Z
56	CLOSED	AppEngine: object-aggregate response formats table and disjoint_tables	enhancement, upstream-parity	2026-08-23T15:09:31Z
57	CLOSED	AppEngine: audit server-write error taxonomy against upstream	enhancement, upstream-parity	2026-08-24T13:24:06Z
58	CLOSED	AppEngine: PATCH requires Content-Type application/merge-patch+json	enhancement, upstream-parity	2026-08-22T03:15:59Z
59	CLOSED	AppEngine: group create-body validation + UUID-v1 from_token for group device listing	enhancement, upstream-parity	2026-08-22T10:35:18Z
60	CLOSED	Realm Management: GET config/datastream_maximum_storage_retention (since upstream 1.2.0)	enhancement, upstream-parity	2026-08-22T03:23:47Z
61	CLOSED	Realm Management: audit interface/mapping validation matrix against astarte_core	enhancement, upstream-parity	2026-08-23T15:09:28Z
62	CLOSED	Realm Management: audit install/update/delete error codes and statuses	enhancement, upstream-parity	2026-08-23T15:09:26Z
63	CLOSED	Triggers: HTTP action validation limits (URL/method/header blocklist/template size)	enhancement, upstream-parity	2026-08-23T15:09:25Z
64	CLOSED	Triggers: decide AMQP action behavior (validate-reject vs NATS-forward deviation)	enhancement, upstream-parity	2026-08-23T15:09:23Z
65	CLOSED	Policies: handler-overlap rejection + retry_times coupling + prefetch_count	enhancement, upstream-parity	2026-08-23T15:09:22Z
66	CLOSED	Realm Management: detailed=true interface listing with full mappings (upstream 1.4)	enhancement, upstream-parity	2026-08-23T15:09:20Z
67	CLOSED	Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4)	enhancement, upstream-parity, upstream-experimental	2026-09-04T18:34:19Z
68	CLOSED	Decide async_operation=false params vs documented always-sync deviation	enhancement, mule-blocked, upstream-parity, upstream-experimental	2026-09-04T18:34:35Z
69	CLOSED	Verify unknown-realm HTTP status on RM endpoints against upstream	enhancement, upstream-parity	2026-08-24T13:23:42Z
70	CLOSED	Triggers: audit wildcard semantics (interface_name '*', match_path '/*' forcing rules)	enhancement, upstream-parity	2026-08-23T15:09:18Z
71	CLOSED	Pairing: realm-scoped health check GET /v1/{realm}/health (upstream 1.3)	enhancement, upstream-parity	2026-08-23T15:09:16Z
72	CLOSED	Realms: datastream_maximum_storage_retention ceiling (create/patch/enforce)	enhancement, upstream-parity	2026-08-23T15:09:15Z
73	CLOSED	Housekeeping: default datastream retention injection env var (upstream 1.4)	enhancement, upstream-parity	2026-08-23T15:09:13Z
74	CLOSED	Housekeeping: PATCH /v1/realms/{realm} (jwt key, registration limit, retention; null=unset)	enhancement, upstream-parity	2026-08-23T15:09:12Z
75	CLOSED	Housekeeping: decide realm-deletion gating/preconditions vs always-sync deviation	enhancement, upstream-parity	2026-08-23T15:09:10Z
76	CLOSED	Housekeeping: GET /v1/realm-defaults/replication — decide reject/deviate (Cassandra-shaped)	enhancement, upstream-parity	2026-08-22T19:12:29Z
77	CLOSED	Verify per-service version endpoints (GET /version, GET /v1/{realm}/version) served everywhere	enhancement, upstream-parity	2026-08-24T13:23:44Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
79	CLOSED	Verify registration-limit-reached HTTP status vs upstream	enhancement, upstream-parity	2026-08-24T13:21:57Z
80	CLOSED	Flow blocks: pure-transform set (to_json, update_metadata, split_map, random_source, sort)	enhancement, upstream-parity	2026-08-23T15:09:07Z
81	CLOSED	Flow block: json_path_map	enhancement, upstream-parity	2026-08-23T15:09:06Z
82	CLOSED	Flow blocks: http_source/http_sink (demand-driven)	enhancement, upstream-parity	2026-08-23T15:09:04Z
83	CLOSED	Flow blocks: mqtt_source/mqtt_sink (+modbus_tcp_source?) demand-driven	enhancement, upstream-parity	2026-08-23T15:09:03Z
84	CLOSED	Flow blocks: virtual_device_pool / dynamic_virtual_device_pool	enhancement, upstream-parity	2026-08-25T14:56:53Z
85	CLOSED	Flow API: user-defined composite blocks	enhancement, upstream-parity	2026-08-23T21:47:32Z
86	CLOSED	Flow: pipeline source DSL — keep DAG-JSON as documented deviation?	enhancement, upstream-parity	2026-08-23T11:57:02Z
87	CLOSED	Flow block: lua_map — needs embedded Lua runtime (parked)	enhancement, upstream-parity	2026-09-04T19:20:52Z
88	CLOSED	Flow auth: support a_f JWT claim	enhancement, upstream-parity	2026-08-23T15:08:59Z
89	CLOSED	Dashboard flow-block schema mismatch (split_map/virtual pools hardcoded; null_sink/log_sink unknown)	enhancement, upstream-parity	2026-08-23T15:08:58Z
91	CLOSED	Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health	mule-blocked, upstream-parity	2026-09-04T18:34:38Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z

$ git ls-tree -r origin/mule/research --name-only | rg -i 'upstream-parity' ; echo "===research README==="; git show origin/mule/research:.mule/research/README.md 2>&1 | head -40
===research README===
# Branch `mule/research`

Output of the daily cross-project survey (`tools/mule-survey.sh`, `.mule/recipes/cross-project-survey.md`),
kept separate from `mule/queue` deliberately — see the commit that introduced this branch for why.

Each day this branch gets either:

- a new `.mule/research/survey-<date>.md` + `.mule/research/issues-<date>.md`, or
- nothing (a one-line note in `.mule/research/log.md`), when nothing material changed since
  the last survey.

**Nothing here is a queued task.** `issues-<date>.md` candidates need a human or a strong
model to triage before anything becomes a GitHub issue labelled `mule` — some of them may
look machine-checkable and not be (see `.mule/for-giulio.md` and issues #17-#19 on the main
repo for the shape of that trap: recording against a live upstream needs an architect, not an
unsupervised agent with no diff review).

The worked example this recipe was distilled from — read before triaging a new day's
output — is `survey-2026-07-27.md` / `issues-2026-07-27.md`, right here on this branch, from
the one-off run that produced GitHub issues #12-#19 and four entries in `.mule/for-giulio.md`
on `main`. Triage is done with the `mule-triage` skill; `triaged.md` in this directory tracks
which dates have already been handled.

$ echo "=== OPEN upstream-parity issues ==="; gh issue list --label upstream-parity --state open --limit 100 2>&1; echo; echo "=== issue #90 ==="; gh issue view 90 --json number,title,state,labels,body 2>&1 | python3 -m json.tool 2>/dev/null | head -60
=== OPEN upstream-parity issues ===
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z

=== issue #90 ===
{
    "body": "Reminder of a pending decision parked during the upstream-parity investigation (#47\u2013#89). **Do not raise the version now** \u2014 this is a milestone-4.0 item.\n\n## Context\n\n`realm.APICompatVersion` (`internal/realm/service.go`) is currently hardcoded to `\"1.2.2\"` \u2014 the emulated upstream API level Astrate declares.\n\nEdgehog 0.13.1 gates just-in-time provisioning triggers (`device_registered`, `device_deletion_started`) behind a Realm Management version >= 1.3. Astrate's engine already supports those trigger kinds end-to-end \u2014 verified in #28 (Edgehog compatibility investigation) \u2014 but Edgehog refuses to use them because of the declared version string.\n\n## What this issue must cover when picked up\n\n1. **Audit first**: what else does declaring \"1.3\" imply for other consumers?\n   - `astarte-dashboard` is currently validated against 1.2.2 behavior (see #13, error-name bundle comparison). Check whether any dashboard code path feature-gates on RM version besides the policies UI gate (>= 1.1.1).\n   - Sweep the parity findings (#47\u2013#89): several are \"1.3/1.4 features\" we deliberately do NOT implement yet \u2014 does declaring 1.3 overpromise any of them?\n   - Check device SDKs and `astartectl` for version-gated behavior.\n2. **Then decide one of**:\n   - bump the constant to `\"1.3.0-rc.0\"` (upstream's own pre-release numbering for that line);\n   - or make it configuration (`internal/config`) so operators can declare the level their consumers expect.\n3. Update `docs/COMPATIBILITY.md` deviation #10 wording accordingly (propose, don't edit directly).\n\n## Outcome if done today (why it's parked)\n\nBumping without the audit risks silently promising surfaces we haven't verified; not bumping keeps Edgehog's JIT provisioning disabled even though the engine supports it. Milestone-4.0 scope decision.",
    "labels": [
        {
            "id": "LA_kwDORmfsJs8AAAACbTsJIQ",
            "name": "enhancement",
            "description": "New feature or request",
            "color": "a2eeef"
        },
        {
            "id": "LA_kwDORmfsJs8AAAACtvELrg",
            "name": "milestone-4.0",
            "description": "",
            "color": "5319e7"
        }
    ],
    "number": 90,
    "state": "CLOSED",
    "title": "milestone-4.0: decide on bumping APICompatVersion from \"1.2.2\" to \"1.3.0-rc.0\" (or making it configurable)"
}

$ echo "=== milestone-3.0 label, all ==="; gh issue list --label milestone-3.0 --state all --limit 100 2>&1; echo; echo "=== all open issues repo-wide ==="; gh issue list --state open --limit 100 2>&1
=== milestone-3.0 label, all ===

=== all open issues repo-wide ===
105	OPEN	mule: nothing has landed in 14h	mule-alarm	2026-09-17T10:51:07Z
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

$ cat .mule/research/triaged.md; echo "===log tail==="; tail -30 .mule/research/log.md
cat: .mule/research/triaged.md: No such file or directory
tail: cannot open '.mule/research/log.md' for reading: No such file or directory
===log tail===

$ git show origin/mule/research:.mule/research/upstream-parity-2026-08-22.md 2>&1 | head -5; echo "===issues-2026-08-22==="; git show origin/mule/research:.mule/research/issues-2026-08-22.md 2>&1 | head -30
fatal: path '.mule/research/upstream-parity-2026-08-22.md' does not exist in 'origin/mule/research'
===issues-2026-08-22===
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

$ for n in $(seq 47 89); do gh issue view $n --json number,state,title,labels -q '.number+"\t"+.state+"\t"+(.labels|map(.name)|join(","))+"\t"+.title' 2>/dev/null; done | column -t -s $'\t' | head -80
(no output)

$ for n in $(seq 89 47) 91 92 93; do gh issue view $n --json number,state --jq '"#\(.number) \(.state)"' 2>&1; done | paste -sd' ' - | fold
#91 CLOSED #92 OPEN #93 OPEN

$ for n in 48 60 47 89 51; do echo "== $n =="; gh issue view $n --json number,state,title --jq '"\(.number) \(.state) \(.title)"' 2>&1; done
== 48 ==
48 CLOSED MQTT v1: tolerate devices connecting with an empty introspection
== 60 ==
60 CLOSED Realm Management: GET config/datastream_maximum_storage_retention (since upstream 1.2.0)
== 47 ==
47 CLOSED MQTT v1: accept device capabilities incl. purge_properties_compression_format (plaintext purge)
== 89 ==
89 CLOSED Dashboard flow-block schema mismatch (split_map/virtual pools hardcoded; null_sink/log_sink unknown)
== 51 ==
51 CLOSED keyAgreement control topic + encrypted mappings (upstream 1.4 experimental) — investigate/park

$ git ls-tree -r origin/mule/research --name-only | rg -i '2026-08-22|parity'; echo "---check astrate repo research branch survey---"; git cat-file -e origin/mule/research:.mule/research/survey-2026-08-22.md 2>&1 && echo "survey-2026-08-22 EXISTS" || echo "no survey-2026-08-22"
.mule/research/issues-2026-08-22.md
.mule/research/survey-2026-08-22.md
---check astrate repo research branch survey---
survey-2026-08-22 EXISTS

$ gh label list --limit 100 2>&1 | rg -i 'milestone|parity|decide|mule' 
mule	Queued for the autonomous mule to pick up	#5319e7
mule-alarm	Filed by the mule's dead-man switch: it has gone quiet	#b60205
mule-review	Mule pushed a change for this; waiting for a human or a strong model to review it	#0e8a16
mule-blocked	The mule tried this and could not do it; see its comment	#d93f0b
milestone-2.0	astarte-flow feature parity milestone	#0E8A16
milestone-4.0		#5319e7
milestone-3.0		#5319e7
upstream-parity	Feature parity with upstream astarte-platform	#0052cc

$ echo "=== milestone-3.0 state all ==="; gh issue list --label milestone-3.0 --state all --limit 100 2>&1
=== milestone-3.0 state all ===

$ gh issue list --label upstream-parity --state all --limit 100 2>&1 | awk -F'\t' '{print $1, $2, $4}' | sort -n
47 CLOSED enhancement, upstream-parity
48 CLOSED enhancement, upstream-parity
49 CLOSED enhancement, upstream-parity
50 CLOSED enhancement, upstream-parity
51 CLOSED enhancement, upstream-parity
52 CLOSED enhancement, upstream-parity
53 CLOSED enhancement, upstream-parity
54 CLOSED enhancement, upstream-parity
55 CLOSED enhancement, upstream-parity
56 CLOSED enhancement, upstream-parity
57 CLOSED enhancement, upstream-parity
58 CLOSED enhancement, upstream-parity
59 CLOSED enhancement, upstream-parity
60 CLOSED enhancement, upstream-parity
61 CLOSED enhancement, upstream-parity
62 CLOSED enhancement, upstream-parity
63 CLOSED enhancement, upstream-parity
64 CLOSED enhancement, upstream-parity
65 CLOSED enhancement, upstream-parity
66 CLOSED enhancement, upstream-parity
67 CLOSED enhancement, upstream-parity, upstream-experimental
68 CLOSED enhancement, mule-blocked, upstream-parity, upstream-experimental
69 CLOSED enhancement, upstream-parity
70 CLOSED enhancement, upstream-parity
71 CLOSED enhancement, upstream-parity
72 CLOSED enhancement, upstream-parity
73 CLOSED enhancement, upstream-parity
74 CLOSED enhancement, upstream-parity
75 CLOSED enhancement, upstream-parity
76 CLOSED enhancement, upstream-parity
77 CLOSED enhancement, upstream-parity
78 OPEN enhancement, milestone-4.0, upstream-parity
79 CLOSED enhancement, upstream-parity
80 CLOSED enhancement, upstream-parity
81 CLOSED enhancement, upstream-parity
82 CLOSED enhancement, upstream-parity
83 CLOSED enhancement, upstream-parity
84 CLOSED enhancement, upstream-parity
85 CLOSED enhancement, upstream-parity
86 CLOSED enhancement, upstream-parity
87 CLOSED enhancement, upstream-parity
88 CLOSED enhancement, upstream-parity
89 CLOSED enhancement, upstream-parity
91 CLOSED mule-blocked, upstream-parity
92 OPEN enhancement, upstream-parity, upstream-experimental
93 OPEN mule-review, upstream-parity

$ git ls-tree -r origin/mule/research --name-only | rg -i 'parity|22' ; echo "===README triaged==="; git show origin/mule/research:.mule/research/README.md 2>&1 | rg -i 'upstream-parity|parity'; echo "===search all branches for the named doc==="; for b in $(git for-each-ref --format='%(refname)' refs/remotes); do if git ls-tree -r $b --name-only 2>/dev/null | rg -q 'upstream-parity-2026-08-22'; then echo "FOUND in $b"; fi; done
.mule/research/issues-2026-08-22.md
.mule/research/survey-2026-08-22.md
internal/pairing/testdata/http/register_422_already.json
internal/pairing/testdata/http/register_422_bad_hwid.json
internal/pairing/testdata/http/register_422_blank_hwid.json
internal/pairing/testdata/http/register_422_limit.json
pkg/interfaceschema/testdata/invalid/22-ttl-without-policy.json
===README triaged===
===search all branches for the named doc===

$ gh issue view 1 --json body --jq .body 2>/dev/null | head; echo "===read the research branch 08-22 issues doc==="; git show origin/mule/research:.mule/research/issues-2026-08-22.md
As of the day this issue was opened, there's plenty of choices when it comes to developing IoT projects ironically, but there's a lack of concrete resources when one has to break through the layers of irony and face production.

We have to make sure developers can rely on a non-ironic solution.
===read the research branch 08-22 issues doc===
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
`keyAgreement/0`–`/4`, InitExchange/ExchangeResp/SecretHash/HashOk/ExchangeFailed,
X25519 suite 1 / ECDH P-256 suite 0, HKDF-derived 32-byte AES-GCM key persisted to
`devices.shared_secret` by `0da8336`) even though its docs still claim otherwise. The
deny stays correct for Astrate for now, but fix the comment's citation when this
lands. Work items unchanged: (a) allow the topic in the broker ACL for the device's
own base, (b) decode the 5 CBOR message types, (c) implement the state machine,
(d) derive and persist the shared secret, (e) emit upstream error names on failure
(`key_agreement_error`, `key_agreement_transition_error`, `secret_hash_error`,
`hash_ok_error`, `exchange_failed_error`). Verification: T3 handshake test and an
invalid-InitExchange test asserting `key_agreement_error`. Prerequisite for
`encrypted-endpoints`.

---

### encrypted-endpoints
Title: Implement encrypted endpoints (DEK/KEK at-rest encryption, COSE Encrypt0 transport both directions, parameter-aware)
Scope: mule-spec
Body:
Carried forward (re-grep-confirmed: no DEK/KEK/COSE in `internal/`; only an unrelated
CA-key comment at `internal/store/crypto.go:15` and JWT test doubles match).
Upstream's two layers unchanged: (1) at-rest DEK/KEK wrap via Astarte Secrets
(`encryptedblob_value`/`encrypted_dek`, per-device shared secret under realm KEK);
(2) COSE Encrypt0 transport on both directions using the key-agreement shared secret
(12-byte random IV + AES-GCM, mapping endpoint parameter-substituted as context,
whole-payload if any mapping is encrypted); read path encrypts (`device.ex`),
DUP write path decrypts before BSON decode (`maybe_decrypt_payload`). Failure modes:
write side `key_agreement_error`/`decryption_error`, read side
`device_not_ready_for_encryption`. Note upstream #1999 (08-21) only made such
interfaces installable — it added no crypto — so this candidate remains the real
work item behind yesterday's attribute plumbing. Fix/verification unchanged: columns,
secrets library, COSE Encrypt0 with endpoint context, encrypt-on-write/
decrypt-before-decode/decrypt-on-read; T3 test incl. parametrized path, byte-identical
round-trip fixture, named-error assertions. Follows `key-agreement-protocol`.

---

### capabilities-malformed-disconnect
Title: Match upstream's malformed-capabilities handling (force-disconnect + `malformed_capabilities_message` error_name)
Scope: mule-line
Body:
Carried forward (re-grep-confirmed against `origin/main`: `handleCapabilities` now at
`internal/engine/capabilities.go:29` (+2 lines from the lint import regroup in
`4c881d0`, logic identical) still rejects on decode/validate failure and keeps the
connection; `malformed_capabilities_message` still absent from `errorname.go`,
only `capabilities_payload_invalid` exists at `internal/engine/data.go:50`).
Upstream force-disconnects on malformed capabilities
(`080-mqtt-v1-protocol.md:125-128`) and fires device_error with
`error_name: "malformed_capabilities_message"` (`capabilities_handler.ex:74-81`);
Astrate surfaces fallback `interface_loading_failed` and lets the device retry
forever. The broker primitive exists unused (`Broker.DisconnectDevice`,
`broker.go:~260`). Fix: call `DisconnectDevice` on the malformed path via the
engine's broker port; map `capabilities_payload_invalid` →
`malformed_capabilities_message` in `errorname.go` (coordinate with
`upstream-error-names-new` so the mapping lands once). Verification: T2 test
asserting session termination + correct error_name. Secondary note unchanged:
upstream persists capabilities per device (`090-database.md:123-129`), Astrate keeps
them session-scoped — record as documented divergence if not persisted.

---

### upstream-error-names-new
Title: Add `decryption_error` and `malformed_capabilities_message` to the upstream error name set
Scope: mule-line
Body:
Carried forward (re-grep-confirmed: neither name appears anywhere in
`origin/main:internal/engine/triggers/errorname.go`; unmapped reasons fall back to
`interface_loading_failed` via `UpstreamErrorName`). Upstream emits
`decryption_error` on failed decryption of an incoming encrypted payload
(`data_handler.ex:453`) and `malformed_capabilities_message` on malformed capabilities
BSON (`capabilities_handler.ex:78`); neither is in the Dashboard v1.2.2 validated set,
so Channels would discard them either way — but parity means emitting what upstream
emits. Fix: add both to `upstreamErrorNames` plus `astrateToUpstream` entries for
Astrate's producers (`capabilities_payload_invalid` →
`malformed_capabilities_message` overlaps `capabilities-malformed-disconnect` — land
once, there; `decryption_error` has no producer until `encrypted-endpoints`).
Verification: extend `TestUpstreamErrorNameMapping` and the closed-set invariant test
in `errorname_test.go`.

---

### empty-cache-error-triggers
Title: Fire device_error triggers on `control/emptyCache` resend failure, matching upstream's error_name split
Scope: mule-spec
Body:
Carried forward (re-grep-confirmed against `origin/main`: `handleEmptyCache`
`internal/engine/control.go:58-76` wraps resend in `retryStore` and on permanent
failure returns without ack — no device_error trigger, no session end;
`resendServerProperties` skips unresolvable interfaces/unmatched paths/unrehydratable
values loudly with `continue`; `resend_interface_properties_failed`/
`empty_cache_error` stay mapped-but-never-emitted, `errorname.go:20-21`, only fixture
uses). Upstream fires a distinct device_error per failure class and ends the session:
`interface_loading_failed` on load failure, `resend_interface_properties_failed` when
sending fails, `empty_cache_error` (reason in `extra_error_metadata`) otherwise.
Design decision still required: Astrate's skip-loudly behaviour is deliberate
robustness vs upstream's fail+disconnect — adopt upstream semantics, diverge with a
COMPATIBILITY.md note, or hybrid (skip quietly per row, trigger on total failure).
Verification: T2/T3 test asserting the right error_name fires on emptyCache resend
failure.

---

### unset-not-allowed-mapping
Title: Fix `unset_not_allowed` error name mapping to match upstream
Scope: mule-line
Body:
Carried forward (re-grep-confirmed: `errorname.go:62` still maps
`"unset_not_allowed"` → `"unexpected_value_type"`, wrong value still pinned by
`errorname_test.go:30`). Upstream emits the literal error name when a device unsets a
property with `allow_unset: false` (`data_handler.ex:255-260`). Not in the Dashboard
validated set either way; emitted name should match upstream regardless. Fix:
pass-through mapping. Verification: T1 test asserting
`UpstreamErrorName("unset_not_allowed") == "unset_not_allowed"`; update the pinning
test.

---

### bson-null-value-acceptance
Title: Verify whether a BSON-null individual value must be accepted (upstream `4789d047`)
Scope: mule-line
Body:
Carried forward, evidence stable (re-grep-confirmed: null-rejection pinned by
`payload_test.go:499` "null v rejected", `ReasonTypeMismatch`; no `bson.TypeNull`
case in the decoder). Upstream restored type validation such that a `nil` individual
value passes validation. Low confidence it's a real gap: official SDKs encode unset
as the empty payload, not `{v: null}`, and no upstream test drives a null document
through the data path — but the acceptance clause is deliberate. Fix, if verification
finds a sending path: accept `bson.TypeNull` for individual mappings as explicit nil
(or reject loudly with a documented divergence); otherwise add a one-line comment to
`bson.go` documenting the deliberate divergence. Verification: check upstream SDK
encode paths for `{v: null}` emission; record the finding either way.

---

**Deferred, not restated:** `group-scoped-triggers` (design-decision, carried since
07-27) stays parked pending bench group provisioning — nothing changed today
(`match.go:318,354` untouched); it shares its blocker with issue #17.

**Already filed, unchanged, not restated (verify before re-filing anything):**
#17 group WATCH paths (still unmeasured, `COMPATIBILITY.md` deviation 1), #18
`write_on_server_owned_interface` and #19 `value_size_exceeded` verifications (still
unverified guesses, `COMPATIBILITY.md:122`). **#10 needs human re-check** — its
premise ("external-bus intake TODO at store.go:135") does not hold on main; see
survey correction B before acting on it either way.

$ ls -la .mule/research/ 2>&1; echo "=== in repo on this branch? ==="; ls .mule/ 2>&1; echo "=== in origin/mule/research tree? ==="; git ls-tree -r origin/mule/research --name-only | rg 'upstream-parity-2026-08-22|surveys-2026-08-22'
ls: cannot access '.mule/research/': No such file or directory
=== in repo on this branch? ===
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
=== in origin/mule/research tree? ===

$ echo "=== search all branches/trees for the doc ==="; for b in origin/main origin/mule/research origin/mule/queue; do git ls-tree -r $b --name-only 2>/dev/null; done | rg -i 'upstream-parity-2026-08-22|parity' ; echo "=== is 'upstream-parity-2026-08-22' in research README for clients ==="; git show origin/mule/research:.mule/research/README.md | rg -i 'upstream-parity|report|log' | head
=== search all branches/trees for the doc ===
docs/handoff/flow-parity-audit-2026-07-29.md
internal/appengine/formats_parity_test.go
internal/appengine/groups_parity_test.go
internal/engine/mqttv1_parity_test.go
docs/handoff/flow-parity-audit-2026-07-29.md
internal/appengine/formats_parity_test.go
internal/appengine/groups_parity_test.go
internal/engine/mqttv1_parity_test.go
=== is 'upstream-parity-2026-08-22' in research README for clients ===
- nothing (a one-line note in `.mule/research/log.md`), when nothing material changed since

$ git var GIT_AUTHOR_IDENT; echo "=== files in research branch root ==="; git ls-tree origin/mule/research --name-only | grep .mule; git ls-tree -r origin/mule/research --name-only | rg 'survey-2026-08-22'
Giulio Malitesta <atsetilam@gmail.com> 1789670143 +0200
=== files in research branch root ===
.mule
.mule/research/survey-2026-08-22.md

$ echo "=== log.md on research branch (2026-08-22 entry) ==="; git show origin/mule/research:.mule/research/log.md 2>&1 | rg -n -A5 '2026-08-22|upstream-parity' | head -40; echo "=== survey-2026-08-22 first lines ==="; git show origin/mule/research:.mule/research/survey-2026-08-22.md 2>&1 | head -20
=== log.md on research branch (2026-08-22 entry) ===
=== survey-2026-08-22 first lines ===
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

$ git show origin/mule/research:.mule/research/survey-2026-08-22.md | head -100
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

- `store.go:135` holds a different, real TODO: the Timescale toolkit `lttb` downsample
  extension point (`probeCapabilities`). The range `89145e6..4c881d0f` never touched
  `internal/store/`, so this was wrong for `89145e6` too.
- `git log -S 'external-bus' -- internal/store/store.go` over main returns **nothing** —
  such a TODO never existed there.
- Exactly one sanctioned `// TODO` exists on main (`store.go:135`). Of ROADMAP §0.1
  rule 3's three named extension points, NATS/HTTP trigger forwarding appears
  **implemented** (`internal/engine/forward/nats.go` + test), and the engine has a live
  fan-out bus (`internal/engine/engine.go:68,123`, M7b stream socket) — so the third
  point (external-bus intake, §1.4) may be resolved or reshaped rather than "blocked".
  **Human action:** re-check what GitHub issue #10 actually tracks before anyone acts on
  it; do not treat it as covered by a store TODO.

---

## Source 1 — Astrate codebase and docs (`origin/main` 4c881d0f)

### What the five new commits did (none change parity)

- `557a169` — stop excluding `docs/` from the Docker build context (.dockerignore fix).
- `763be9d` — test-only state sync in `TestRegistry_InstantiateLinear`.
- `801fc48` — fix(flow): stop already-constructed blocks when `Instantiate` fails
  partway (container leak on partial pipeline build; found scoping #42).
- `403059a` — `.mule/milestones.md` only: mark Flow v2.0 DONE, close #42, split remainders
  into #44/#45/#46. No changes to DESIGN/ROADMAP/COMPATIBILITY in the whole range.
- `4c881d0` — golangci-lint v2.12.2 sweep: renames in `internal/flow`
  (FlowMessage→Message etc.), doc comments, gofmt, FNV-1a byte-wise in `laneOf`,
  example hardening. In the engine proper this is **import-grouping only**
  (`capabilities.go` +2/-1, no logic), so all pre-existing citations hold with small
  line shifts.

### Standing gaps re-grep-verified today against `4c881d0f`

All rows below were individually confirmed still-open (grep/sed against the fresh clone):

- `unset_not_allowed` mis-map — `triggers/errorname.go:62` → `"unexpected_value_type"`;
  wrong value pinned by `errorname_test.go:30`.
- `decryption_error` / `malformed_capabilities_message` — absent from `errorname.go`.
- Malformed capabilities reject-not-disconnect — `handleCapabilities` now at
  `internal/engine/capabilities.go:29`; `Broker.DisconnectDevice` (`broker.go:~260`)
  still uncalled from that path.
- Key agreement ACL-denied — `aclhook.go:63-68` (see correction A for the citation fix).
- Encrypted-endpoint machinery absent — `grep -ri 'dek|kek|cose' internal/` matches only
  an unrelated CA-key comment (`internal/store/crypto.go:15`) and JWT test doubles.
- emptyCache resend failure silent — `handleEmptyCache` `control.go:58-76`: `retryStore`
  wraps resend, permanent failure returns without ack; no device_error trigger, no
  session end. `resend_interface_properties_failed`/`empty_cache_error` stay
  mapped-but-never-emitted (`errorname.go:20-21`).
- BSON-null individual value rejected — pinned by `payload_test.go:499` ("null v

$ git show origin/mule/research:.mule/research/survey-2026-08-22.md | wc -l; git show origin/mule/research:.mule/research/issues-2026-08-22.md | wc -l
169
193

$ git show origin/mule/research:.mule/research/survey-2026-08-22.md
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

- `store.go:135` holds a different, real TODO: the Timescale toolkit `lttb` downsample
  extension point (`probeCapabilities`). The range `89145e6..4c881d0f` never touched
  `internal/store/`, so this was wrong for `89145e6` too.
- `git log -S 'external-bus' -- internal/store/store.go` over main returns **nothing** —
  such a TODO never existed there.
- Exactly one sanctioned `// TODO` exists on main (`store.go:135`). Of ROADMAP §0.1
  rule 3's three named extension points, NATS/HTTP trigger forwarding appears
  **implemented** (`internal/engine/forward/nats.go` + test), and the engine has a live
  fan-out bus (`internal/engine/engine.go:68,123`, M7b stream socket) — so the third
  point (external-bus intake, §1.4) may be resolved or reshaped rather than "blocked".
  **Human action:** re-check what GitHub issue #10 actually tracks before anyone acts on
  it; do not treat it as covered by a store TODO.

---

## Source 1 — Astrate codebase and docs (`origin/main` 4c881d0f)

### What the five new commits did (none change parity)

- `557a169` — stop excluding `docs/` from the Docker build context (.dockerignore fix).
- `763be9d` — test-only state sync in `TestRegistry_InstantiateLinear`.
- `801fc48` — fix(flow): stop already-constructed blocks when `Instantiate` fails
  partway (container leak on partial pipeline build; found scoping #42).
- `403059a` — `.mule/milestones.md` only: mark Flow v2.0 DONE, close #42, split remainders
  into #44/#45/#46. No changes to DESIGN/ROADMAP/COMPATIBILITY in the whole range.
- `4c881d0` — golangci-lint v2.12.2 sweep: renames in `internal/flow`
  (FlowMessage→Message etc.), doc comments, gofmt, FNV-1a byte-wise in `laneOf`,
  example hardening. In the engine proper this is **import-grouping only**
  (`capabilities.go` +2/-1, no logic), so all pre-existing citations hold with small
  line shifts.

### Standing gaps re-grep-verified today against `4c881d0f`

All rows below were individually confirmed still-open (grep/sed against the fresh clone):

- `unset_not_allowed` mis-map — `triggers/errorname.go:62` → `"unexpected_value_type"`;
  wrong value pinned by `errorname_test.go:30`.
- `decryption_error` / `malformed_capabilities_message` — absent from `errorname.go`.
- Malformed capabilities reject-not-disconnect — `handleCapabilities` now at
  `internal/engine/capabilities.go:29`; `Broker.DisconnectDevice` (`broker.go:~260`)
  still uncalled from that path.
- Key agreement ACL-denied — `aclhook.go:63-68` (see correction A for the citation fix).
- Encrypted-endpoint machinery absent — `grep -ri 'dek|kek|cose' internal/` matches only
  an unrelated CA-key comment (`internal/store/crypto.go:15`) and JWT test doubles.
- emptyCache resend failure silent — `handleEmptyCache` `control.go:58-76`: `retryStore`
  wraps resend, permanent failure returns without ack; no device_error trigger, no
  session end. `resend_interface_properties_failed`/`empty_cache_error` stay
  mapped-but-never-emitted (`errorname.go:20-21`).
- BSON-null individual value rejected — pinned by `payload_test.go:499` ("null v
  rejected", `ReasonTypeMismatch`).
- Group-scoped triggers compile-but-never-match — `match.go:318,354`.

---

## Source 2 — mule queue

`origin/mule/queue` `4f6d95c` → `28f3adf5`. Delta: `83568dd`, `3a4b693`, `3a746d1`
(recipe-run bookkeeping) and `28f3adf` (mule cron-cwd fix + queue-state recovery —
tooling, not Astrate). Still **no new implementation commits** on the queue branch; its
three headline implementations remain superseded-on-main via `aba2d18` (see the 08-21
correction). `.mule/log.md` task table still has zero rows.

---

## Source 3 — upstream `astarte-platform`

### master delta (`739724991..372a2cbe`, 1 commit)

`372a2cb feat: add encrypted endpoints (#1999)` — "enable installation of interfaces
with encrypted endpoints". 16 files, +59/-4, **no crypto**:

- Realm Management API schema gains `encrypted: boolean, default: false` on mappings
  (`apps/astarte_realm_management/lib/astarte_realm_management_web/api_spec/schemas/
  mapping.ex:113-119`, description: "Defines if the endpoint is encrypted or not.").
- Persistence plumbing: `interfaces/queries.ex` insert gains an `encrypted` column
  binding; `libs/astarte_data_access/mappings.ex` selects it; adapters/interface struct
  gains the field; fetchdata queries adjusted.
- Dashboard UI + client models/DTO updated.

So the transport machinery (DUP `key_agreement/*`, COSE Encrypt0 paths in appengine/DUP)
predates this commit — what #1999 changes is the **installation boundary**: interfaces
carrying `"encrypted": true` on a mapping previously failed API validation; now they
install and persist. Note `CHANGELOG.md`'s Unreleased section is still empty (release-
please was added 08-19 but has not generated an entry yet), and no tag newer than
`v1.4.0-rc.5` exists.

### Release branches

No movement beyond rc.5/v1.3.3 (both analysed 08-21). Nothing new here.

---

## Source 4 — parity gap analysis (one row added)

| Area | Upstream behaviour | Astrate status (main 4c881d0f) | Reference |
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

