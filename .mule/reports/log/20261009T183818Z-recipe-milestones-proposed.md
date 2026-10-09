slug: recipe-milestones
verdict: proposed
at:  cf5d7b4
ran: 2026-10-09T18:38:18Z on DietPi in 181s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/todo.md
→ Read .mule/milestones.md
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md

$ gh issue list --label milestone-3.0 --state all --limit 50; echo "---LABELS---"; gh label list --limit 100 | grep -i milestone
---LABELS---
milestone-2.0	astarte-flow feature parity milestone	#0E8A16
milestone-4.0		#5319e7
milestone-3.0		#5319e7

→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ ls -la .mule/research/ 2>&1; echo "---"; find .mule -name '*upstream-parity*' 2>&1; echo "---ISSUES upstream-parity---"; gh issue list --label upstream-parity --state all --limit 60
ls: cannot access '.mule/research/': No such file or directory
---
---ISSUES upstream-parity---
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

$ ls -la .mule/for-giulio.md 2>&1; echo "=== tail for-giulio ==="; tail -60 .mule/for-giulio.md 2>&1
-rw-r--r-- 1 root root 161185 Oct  9 13:28 .mule/for-giulio.md
=== tail for-giulio ===
(cmd/astrate/main.go:531) — both real, both yours to call, neither has a test worth
writing. See `.mule/reviews/cmd-astrate-2026-09-25.md`.
- **The mule has been idle 16h.** Filed by the dead-man's switch; see journalctl on the Pi.

## 2026-10-04 — three housekeeping drifts in `docs/site/` the recipe is not allowed to touch

Ran the docs-sync recipe over the housekeeping surface (code routes vs
`docs/api/astarte_housekeeping_api.yaml`). The spec half is queued as four lines in
`.mule/todo.md`; the prose half needs you, because `docs/site/` is on the never-touch list.

**1. `docs/site/housekeeping-api.md:30` still talks about Cassandra.**

> Cassandra-specific fields (`replication_class`, `replication_factor`, etc.) from upstream are accepted but ignored.

Astrate stores PostgreSQL — `internal/store` on pgx, schema in `migrations/*.sql`; there is no
Cassandra anywhere in the tree. The wire shape is four fields (`internal/housekeeping/http.go:46-51`),
and the comment right above it says the opposite of "accepted": *"Astrate omits the
Cassandra-specific fields (replication factor/class) upstream carries"*
(`internal/housekeeping/http.go:44-45`). From a client's side the two readings are the same
(the request succeeds, the fields do nothing), so this is vocabulary rather than behaviour —
but it is the kind of sentence that sends an operator hunting for a Cassandra setting that
cannot exist.

**2. `PATCH /housekeeping/v1/realms/{realm}` is missing from the page entirely.**
`docs/site/housekeeping-api.md` goes Create → List → Get → Delete (lines 10-50), so the route
`docs-sync-hk-patch-endpoint` added and documented in the spec
(`docs/api/astarte_housekeeping_api.yaml:142-192`) has no prose at all. The same page also
predates `datastream_maximum_storage_retention`, which the spec has carried on both `Realm` and
`RealmCreate` since `docs-sync-hk-retention-field`, so its create example (lines 12-21) shows a
three-field body that is not the whole body.

**3. Two `[housekeeping]` config keys exist; `configuration-reference.md` does not list them.**
The table at `docs/site/configuration-reference.md:79-84` has exactly two rows and the section
heading calls the block "Instance-admin keys" — but `HousekeepingConfig`
(`internal/config/config.go:109-125`) has four fields:

- `default_datastream_maximum_storage_retention` (`*int64`, env
  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`, validated as a
  non-negative integer at `config.go:286-291`) — injected into a realm at creation when the
  caller omits the field (`internal/housekeeping/service.go:149-150`). An operator who wants it
  has no way to learn it exists.
- `realm_deletion_disabled` (bool, env `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`, a
  fail-loud boolean accepting only `1|true|TRUE|True|0|false|""` and refusing to load on anything
  else, `config.go:296-306`) — this is the switch behind the spec's documented
  `405 "Realm deletion disabled"` on `DELETE /housekeeping/v1/realms/{realm}`, so it changes a
  status code an API client sees.

Both belong in that table next to the two key lists, and the heading wants widening.

Checked and *not* escalated, for the record. The config-key sweep
(`rg -o '\bASTRATE_[A-Z_]+' -N internal/ | sort -u`) found nothing else stale: the two
`jwt_public_key*` keys the page quotes at lines 69-71 are real (`config.go:112-113`), and the
`master_key_file` row's `ASTRATE_MASTER_KEY`/`ASTRATE_MASTER_KEY_FILE` fallbacks are real too
(`internal/store/crypto.go:24-27`). And the one thing that looked like a live bug is not one:
`view()` returns only `JWTPublicKeysPEM[0]` (`internal/housekeeping/service.go:277-281`) while
the realm-management surface joins *all* keys (`internal/realm/service.go:627`) — but every
writer stores a single-element array (create at `internal/store/realms.go:54-61`, PATCH at
`internal/store/realms.go:165`, `putAuthConfig` at `internal/realm/service.go:636`), so a
two-key realm is unreachable and there is nothing to document. Worth noting because it reads
like a bug on first inspection.

$ ls -la .mule/ 2>&1; echo "=== git status ==="; git status --short 2>&1 | head -40; echo "=== branch ==="; git branch --show-current 2>&1; echo "=== git log research file ==="; git log --oneline --all -- .mule/research/upstream-parity-2026-08-22.md 2>&1 | head
total 16348
drwxr-xr-x  7 root root     4096 Oct  9 20:07 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       13 Oct  9 20:35 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  9 19:51 .heartbeat
-rw-r--r--  1 root root     4096 Oct  9 20:36 .last-output
-rw-r--r--  1 root root        0 Oct  9 20:35 .lock
-rw-r--r--  1 root root       11 Oct  9 20:35 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 16291240 Oct  9 20:36 cron.log
drwxr-xr-x  2 root root     4096 Oct  9 20:07 failed
-rw-r--r--  1 root root   161185 Oct  9 13:28 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    16266 Oct  9 20:07 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  8 19:28 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   175392 Oct  9 20:07 todo.md
=== git status ===
=== branch ===
mule/queue
=== git log research file ===

$ cat .mule/.gitignore; echo "=== tracked .mule files ==="; git ls-files .mule/ | head -40; echo "=== research in git all history ==="; git log --all --oneline -- '**/upstream-parity*' 2>&1 | head; echo "=== find any parity file ==="; find / -name '*upstream-parity*' 2>/dev/null | head
# The mule's own working files. None of these belong in git, and — more sharply — any of
# them that is not listed here makes the tree dirty, which makes `tick` refuse to run and so
# turns the timer into a permanent silent no-op.
failed/
task.md
.timeout
.lock
.budget
.bench-linux-amd64
.last-output
.rotation
.heartbeat
.alarmed
.impl.patch
hosts
research/
=== tracked .mule files ===
.mule/.gitignore
.mule/MULE.md
.mule/config
.mule/for-giulio.md
.mule/hosts.example
.mule/log.md
.mule/milestones.md
.mule/recipes/astarte-upstream.md
.mule/recipes/benchmarks.md
.mule/recipes/code-review.md
.mule/recipes/cross-project-survey.md
.mule/recipes/docs-sync.md
.mule/recipes/github-issues.md
.mule/recipes/hygiene.md
.mule/recipes/legion-go.md
.mule/recipes/milestones.md
.mule/reports/hygiene.md
.mule/reports/log/20260904T190911Z-purge-properties-compression-capabilityauto-done.md
.mule/reports/log/20260904T195018Z-broker-acl-coldstart-introspection-miss-done.md
.mule/reports/log/20260904T195529Z-issue-93-done.md
.mule/reports/log/20260904T195814Z-broker-disconnect-device-zombie-session-blocked.md
.mule/reports/log/20260904T203003Z-broker-offline-acl-tests-blocked.md
.mule/reports/log/20260904T210532Z-broker-onconnect-doc-comment-blocked.md
.mule/reports/log/20260905T110211Z-empty-introspection-verification-blocked.md
.mule/reports/log/20260905T111428Z-probe-trigger-install-notification-delay-blocked.md
.mule/reports/log/20260905T112539Z-compat-note-v132-done.md
.mule/reports/log/20260905T165512Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260905T171755Z-docs-sync-pairing-health-path-done.md
.mule/reports/log/20260905T172234Z-docs-sync-pairing-register-404-done.md
.mule/reports/log/20260905T173657Z-recipe-hygiene-proposed.md
.mule/reports/log/20260905T174528Z-swagger-httptest-coverage-done.md
.mule/reports/log/20260905T181053Z-recipe-milestones-proposed.md
.mule/reports/log/20260905T181219Z-recipe-github-issues-proposed.md
.mule/reports/log/20260905T184625Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260905T184903Z-probe-property-resend-encoding-blocked.md
.mule/reports/log/20260905T190855Z-compat-note-v133-blocked.md
.mule/reports/log/20260905T192010Z-recipe-code-review-proposed.md
.mule/reports/log/20260905T192911Z-flow-validate-source-sink-blocked.md
.mule/reports/log/20260905T201529Z-flow-validate-dead-source-sink-recompute-done.md
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md
=== research in git all history ===
=== find any parity file ===

$ ls -la .mule/reports/log/ | tail -20; echo "=== reports dir ==="; ls -R .mule/reports/ | head -40; echo "=== tasks ==="; ls -la .mule/tasks/
-rw-r--r-- 1 root root  12117 Oct  8 19:37 20261008T173721Z-testutil-wait-since-cursor-done.md
-rw-r--r-- 1 root root   3972 Oct  8 19:47 20261008T174752Z-testutil-control-frame-test-blocked.md
-rw-r--r-- 1 root root   6525 Oct  8 19:58 20261008T175838Z-engine-introspection-producer-roundtrip-done.md
-rw-r--r-- 1 root root 134121 Oct  8 20:26 20261008T182643Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root   3621 Oct  8 20:31 20261008T183153Z-docs-sync-pairing-register-404-unreachable-done.md
-rw-r--r-- 1 root root  14077 Oct  8 21:13 20261008T191310Z-docs-sync-pairing-body-cap-400-blocked.md
-rw-r--r-- 1 root root  23561 Oct  8 21:28 20261008T192820Z-docs-sync-pairing-realm-name-pattern-done.md
-rw-r--r-- 1 root root  27083 Oct  8 21:36 20261008T193604Z-docs-sync-pairing-deviceid-path-param-done.md
-rw-r--r-- 1 root root  99670 Oct  8 22:06 20261008T200623Z-recipe-hygiene-timeout.md
-rw-r--r-- 1 root root 186642 Oct  8 22:21 20261008T202151Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  61863 Oct  8 22:38 20261008T203843Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  89855 Oct  8 22:42 20261008T204257Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  67075 Oct  9 12:58 20261009T105827Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root   4761 Oct  9 13:29 20261009T112852Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root 222620 Oct  9 13:29 20261009T112951Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root  47057 Oct  9 19:11 20261009T171147Z-docs-sync-ae-data-envelope-metadata-done.md
-rw-r--r-- 1 root root  78983 Oct  9 19:43 20261009T174328Z-docs-sync-ae-data-write-object-errors-done.md
-rw-r--r-- 1 root root  25806 Oct  9 19:51 20261009T175114Z-docs-sync-ae-group-create-devices-required-done.md
-rw-r--r-- 1 root root  93256 Oct  9 20:07 20261009T180752Z-docs-sync-ae-notfound-component-blocked.md
-rw-r--r-- 1 root root  43667 Oct  9 20:29 20261009T182948Z-recipe-hygiene-proposed.md
=== reports dir ===
.mule/reports/:
hygiene.md
log

.mule/reports/log:
20260904T190911Z-purge-properties-compression-capabilityauto-done.md
20260904T195018Z-broker-acl-coldstart-introspection-miss-done.md
20260904T195529Z-issue-93-done.md
20260904T195814Z-broker-disconnect-device-zombie-session-blocked.md
20260904T203003Z-broker-offline-acl-tests-blocked.md
20260904T210532Z-broker-onconnect-doc-comment-blocked.md
20260905T110211Z-empty-introspection-verification-blocked.md
20260905T111428Z-probe-trigger-install-notification-delay-blocked.md
20260905T112539Z-compat-note-v132-done.md
20260905T165512Z-recipe-docs-sync-proposed.md
20260905T171755Z-docs-sync-pairing-health-path-done.md
20260905T172234Z-docs-sync-pairing-register-404-done.md
20260905T173657Z-recipe-hygiene-proposed.md
20260905T174528Z-swagger-httptest-coverage-done.md
20260905T181053Z-recipe-milestones-proposed.md
20260905T181219Z-recipe-github-issues-proposed.md
20260905T184625Z-recipe-astarte-upstream-proposed.md
20260905T184903Z-probe-property-resend-encoding-blocked.md
20260905T190855Z-compat-note-v133-blocked.md
20260905T192010Z-recipe-code-review-proposed.md
20260905T192911Z-flow-validate-source-sink-blocked.md
20260905T201529Z-flow-validate-dead-source-sink-recompute-done.md
20260906T112006Z-recipe-docs-sync-proposed.md
20260906T113143Z-docs-sync-rm-datastream-retention-endpoint-done.md
20260906T171045Z-docs-sync-rm-interfaces-detailed-param-done.md
20260906T174031Z-recipe-hygiene-proposed.md
20260906T175748Z-housekeeping-tests-done.md
20260906T180442Z-recipe-milestones-proposed.md
20260906T181352Z-recipe-github-issues-proposed.md
20260906T183102Z-recipe-astarte-upstream-proposed.md
20260906T183952Z-probe-props-resend-error-triggers-blocked.md
20260906T185735Z-recipe-code-review-proposed.md
20260906T195827Z-server-data-trigger-bus-done.md
20260906T200431Z-recipe-docs-sync-proposed.md
20260906T201002Z-docs-sync-appengine-by-alias-endpoints-done.md
=== tasks ===
total 60
drwxr-xr-x 2 root root 4096 Sep  7 21:11 .
drwxr-xr-x 7 root root 4096 Oct  9 20:07 ..
-rw-r--r-- 1 root root 3293 Sep  7 21:11 docs-sync-appengine-data-output-params.md
-rw-r--r-- 1 root root 5539 Jul 27 22:18 issue-16.md
-rw-r--r-- 1 root root 7510 Sep  4 20:46 issue-23.md
-rw-r--r-- 1 root root 6434 Sep  4 20:46 issue-24.md
-rw-r--r-- 1 root root 5665 Sep  4 20:46 issue-25.md
-rw-r--r-- 1 root root 6009 Sep  4 20:46 issue-26.md
-rw-r--r-- 1 root root 5683 Sep  4 20:46 issue-27.md

→ Read .mule/reports/log/20260905T181053Z-recipe-milestones-proposed.md

$ head -20 .mule/reports/log/20261008T202151Z-recipe-milestones-proposed.md; echo "=====TAIL====="; tail -150 .mule/reports/log/20261008T202151Z-recipe-milestones-proposed.md
slug: recipe-milestones
verdict: proposed
at:  c1b43cd
ran: 2026-10-08T20:21:51Z on DietPi in 574s

> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
v2.0 is DONE; **v3.0 is the current target**. Scope is decided (not TBD), so I investigate per step 3.

$ ls -la .mule/ .mule/research/ .mule/recipes/ 2>&1; wc -l .mule/research/upstream-parity-2026-08-22.md .mule/todo.md .mule/for-giulio.md docs/UPSTREAM-EXPERIMENTAL.md 2>&1
ls: cannot access '.mule/research/': No such file or directory
.mule/:
total 15368
drwxr-xr-x  7 root root     4096 Oct  8 22:06 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       14 Oct  8 22:12 .budget
=====TAIL=====

$ grep -n "171c1bd7\|vmq_plugin\|route_injection" .mule/reports/log/20261008T113945Z-recipe-astarte-upstream-proposed.md | head; echo "=== changelog diff ==="; gh api repos/astarte-platform/astarte/commits/171c1bd7 --jq '.files[] | select(.filename=="CHANGELOG.md") | .patch' 2>&1 | head -30
107:171c1bd7	2026-10-06T15:57:45Z	chore: forward-port astarte_vmq_plugin from release-1.4
125:$ rg -n "171c1bd7|vmq_plugin|2241|2252|2244|generators" /root/astrate-mule/.mule/todo.md /root/astrate-mule/.mule/for-giulio.md | head -30
134:$ rg -n "vmq|171c1bd7" /root/astrate-mule/.mule/ | head -20; rg -n "Milestone recipe run, 2026-10-0[678]" /root/astrate-mule/.mule/for-giulio.md
135:/root/astrate-mule/.mule/reports/log/20260920T173118Z-recipe-code-review-proposed.md:444:57:| **VerneMQ** + `astarte_vmq_plugin` | MQTT broker, mTLS termination, ACLs, bridging publishes onto AMQP | `internal/broker` | Embedded `mochi-mqtt/server` v2 with auth/ACL hooks; publishes flow into the engine through a Go channel, not AMQP. |
136:/root/astrate-mule/.mule/reports/log/20260920T173118Z-recipe-code-review-proposed.md:455:494:  `device_connected`/`device_disconnected` triggers — the work `astarte_vmq_plugin` does
137:/root/astrate-mule/.mule/reports/log/20261003T193009Z-appengine-payload-reason-status-map-blocked.md:332:/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/vmq_plugin.ex:43:  def publish(topic, payload, qos)
138:/root/astrate-mule/.mule/reports/log/20261003T193009Z-appengine-payload-reason-status-map-blocked.md:333:/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/vmq_plugin.ex:44:      when is_binary(topic) and is_binary(payload) and is_integer(qos) and qos >= 0 and qos <= 2 do
139:/root/astrate-mule/.mule/reports/log/20261003T193009Z-appengine-payload-reason-status-map-blocked.md:334:/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/vmq_plugin.ex:50:        payload: payload,
140:/root/astrate-mule/.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:284:171c1bd7	2026-10-06T15:57:45Z	chore: forward-port astarte_vmq_plugin from release-1.4
141:/root/astrate-mule/.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:289:$ gh api repos/astarte-platform/astarte/commits/171c1bd7 --jq '.files[] | "\(.status)\t+\(.additions)/-\(.deletions)\t\(.filename)"' | head -40; echo "=== msg ==="; gh api repos/astarte-platform/astarte/commits/171c1bd7 --jq '.commit.message'
=== changelog diff ===
@@ -10,6 +10,11 @@ and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.
 ### Added
 
 - [astarte_data_updater_plant] Generate `device_empty_cache_received` trigger events when a device sends `/emptyCache`
+- [astarte_vmq_plugin] Expose a new /healtz endpoint to check whether the AMQP connection toward RabbitMQ instance is ok and all the expected Mississippi EventsProducer workers are running.
+
+### Changed
+
+- [astarte_vmq_plugin] Pin vernemq to 2.2.0
 
 ## [1.4.0-rc.6] - 2026-09-28
 
@@ -30,6 +35,8 @@ Includes changes from v1.3.4
 - [fdo] Allow specifying device_id on ownership voucher upload
 - [fdo] The device is now immediately registered on ownership voucher upload
 - [astarte_data_updater_plant] Improve RPC server reliability
+- [astarte_vmq_plugin] Shard AMQP producers
+- [astarte_vmq_plugin] Increase RPC server reliability
 
 ### Fixed
 
@@ -68,10 +75,12 @@ Includes changes from v1.3.4
 ### Fixed
 
 - [astarte_data_updater_plant] Ensure RPC server is always available to clients. Resolved the issue where a temporary disconnection and reconnection of data_updater_plant to the cluster would make the RPC server inaccessible.
+- [astarte_vmq_plugin] Ensure RPC server is always available to clients. Resolved the issue where a temporary disconnection and reconnection of vmq_plugin to the cluster would make the RPC server inaccessible.
 
 ### Changed
 
 - [astarte_data_updater_plant] Ensure memory is properly garbage collected

$ grep -c "Step 5 deliberately not taken\|step 5" .mule/for-giulio.md; grep -n "Step 5\|step 5" .mule/for-giulio.md | head; echo "=== label escalation? ==="; grep -n "milestone-3.0 label\|milestone-v3.0\|ledger lives under" .mule/for-giulio.md | head
2
44:  **Step 5 deliberately not taken**, same reason as the 21 runs before: the section's declared scope is **1.2.2 → 1.3/1.4** and only the 1.3 half is delivered. The 1.4 half stays open — #92 unanswered, both register rows unreconciled, `APICompatVersion` unbumped — so the tag must not be cut. Your call when a stable v1.4.0 ships: answer #92, reconcile the register, run the final-phase bump, then cut the tag.
61:  **Step 5 deliberately not taken.** I am *not* proposing "milestone v3.0 looks complete, verify and cut the tag", for the same reason as the twenty runs before this one: the section's declared scope is **1.2.2 → 1.3/1.4** and only the 1.3 half is delivered (now including v1.3.5, which is inert anyway). The 1.4 half is open — #92 unanswered, both experimental rows unreconciled, `APICompatVersion` not bumped — so the tag must not be cut. Your call when a stable v1.4.0 ships: answer #92, settle the register rows, run the final-phase bump, then cut the tag.
=== label escalation? ===
41:- **Milestone recipe run, 2026-10-06: v3.0 verified unchanged — still gated on upstream — and the master delta since the 2026-10-05 run is a closed set of 13 commits, every one wire-inert for Astrate. No issue filed, no `.mule/todo.md` line queued.** Release sweep is byte-identical to yesterday's (`gh api repos/astarte-platform/astarte/releases`): **no stable v1.4.0** — newest stable `v1.3.5` (2026-10-05), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition (your 2026-09-04 "wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) stays unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** — this milestone's whole ledger lives under `upstream-parity`, where the open set is unchanged: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`); **#47–#89** otherwise all closed. Non-alarm open issues total exactly four (plus the 25 `mule-alarm` #94–#118): the three above and **#1** (untouched per standing instruction). `APICompatVersion` still `"1.2.2"` (`internal/realm/service.go:588`, re-checked).
59:  **#92's parking condition is still unmet, now for the twenty-first run.** Your 2026-09-04 decision ("wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) — newest stable is `v1.3.5` (today), newest overall remains `v1.4.0-rc.6` (2026-09-30, prerelease). The post-rc.6 master delta has grown by exactly one commit: `731bcddc` "chore: forward port release-1.4" (2026-10-05T09:34Z), and it is the housekeeping/FDO migration relocation the 2026-10-03 run already measured inside rc.5→rc.6 (six `astarte_housekeeping/priv/migrations/{astarte,realm}/*.sql` files removed, table creation moved into migrations) plus a DUP-internal `apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0011_remove_replacement_data.ex` (+31) and `database.ex` (+2/-1). No Realm Management, AppEngine, Pairing or Trigger wire surface. **The milestone ledger is unchanged:** `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** (this milestone's whole ledger lives under `upstream-parity`, as the 2026-09-10 entry established), of #47–#89 **all remain closed**, and the authoritative open set is 29 issues — 25 alarms (#94–#118, one new today: **#118**, that recipe's business) and exactly four non-alarm: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** (untouched per standing instruction). `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`, re-checked), per **#90's** frozen decision.
80:  **Gate re-verified, unchanged.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — nothing shipped upstream since the 2026-10-03 run, so #92's parking condition (your 2026-09-04 decision) is still unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** (this milestone's ledger lives under `upstream-parity`); of #47–#89 **all remain closed**; open non-alarm set unchanged — **#93** (`mule-review`, own recipe path), **#92** (parked), **#78** (`milestone-4.0`), **#1** (untouched). The post-rc.6 master delta is now a closed set of **8 commits** (`gh api .../commits -f since=2026-09-30T07:28:27Z`), of which exactly two carry wire surface and **both already have owners in the queue**: `b6d46ad4` (#2237 unexpected object keys) → done as `appengine-unexpected-object-key`, and `958bb9fe` (#2232 `device_empty_cache_received`) → queued as `device-empty-cache-received-trigger` (currently `[!]` blocked). Neither is from stable 1.3 or 1.4, so neither is a v3.0 gate.
101:  **Verified unchanged, no new gaps.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition, frozen by you on 2026-09-04 ("wait for a stable v1.4.0", an `-rc.N` does not satisfy it), is still unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Of #47–#89 **all are closed**; the open non-alarm set is unchanged — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** (untouched per standing instruction). rc.6 adds no v3.0 gap: it is already audited in the two entries below, and I re-measured its non-FDO half myself (`compare/v1.4.0-rc.5...v1.4.0-rc.6`, 27 commits) — outside FDO the only library changes are a `# coveralls-ignore-stop` line in `apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex:43` and `realms/queries.ex` *removing* the ownership_voucher / to2_sessions / unconfirmed_devices table creation (moved into migrations). No Realm Management, AppEngine or Pairing wire surface moved.
157:- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Status: not started", which is no longer true.** Release sweep (`gh api repos/astarte-platform/astarte/releases`, run this session) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.5` (2026-08-20, prerelease) — byte-identical to what the 2026-09-24 run saw, so nothing shipped upstream in the two days since. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Open non-alarm issues remain exactly **#93** (aclhook comment, `mule-review`, its own recipe path), **#92** (keyAgreement, parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** (FDO, `milestone-4.0`), **#1** (untouched per standing instruction). No new gaps, **no issues filed, no `.mule/todo.md` task lines added** — there is nothing machine-checkable left, and filing anything here would duplicate #92/#93/#78.
232:- **Milestone recipe run, 2026-09-24: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-23 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
240:- **Milestone recipe run, 2026-09-23: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-21 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
256:- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
264:- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
276:- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.

$ grep -rn "milestone-v3.0\|milestone-3.0" .mule/todo.md .mule/recipes/ .mule/MULE.md | head; echo "=== labels on repo ==="; gh label list --limit 50 2>&1 | head -30
=== labels on repo ===
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

$ echo "--- retention ceiling (#72/#60) ---"; grep -rn "datastream_maximum_storage_retention" internal/ --include=*.go -l | head; echo "--- detailed=true listing (#66) ---"; grep -rn "detailed" internal/realm/*.go | head -5; echo "--- merge-patch (#58) ---"; grep -rn "merge-patch+json" internal/appengine/*.go | head -3; echo "--- group show (#54) ---"; grep -rn "groups/{group" internal/appengine/*.go | head -3; echo "--- APICompat ---"; grep -n 'APICompatVersion = ' internal/realm/service.go
--- retention ceiling (#72/#60) ---
internal/realm/service.go
internal/realm/http_test.go
internal/realm/http.go
internal/store/realms.go
internal/store/retention.go
internal/housekeeping/service.go
internal/housekeeping/http_test.go
internal/housekeeping/service_test.go
internal/housekeeping/http.go
internal/swagger/swagger_test.go
--- detailed=true listing (#66) ---
internal/realm/detailed.go:12:// detailed-listing contract promises (reliability unreliable, retention
internal/realm/detailed_required_encrypted_test.go:7:// TestRenderDetailedInterfaceRequiredEncrypted pins that the detailed
internal/realm/detailed_test.go:11:// ?detailed=true serves one fully materialised 1.4-style document per
internal/realm/detailed_test.go:27:	rec := r.req(t, http.MethodGet, "/interfaces?detailed=true", "", r.rmaToken)
internal/realm/detailed_test.go:29:		t.Fatalf("detailed listing: got %d, want 200 (%s)", rec.Code, rec.Body)
--- merge-patch (#58) ---
internal/appengine/http.go:210:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:231:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:522:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
--- group show (#54) ---
internal/appengine/http.go:70:	mux.Handle("GET "+base+"/groups/{group}", h(a.getGroup))
internal/appengine/http.go:71:	mux.Handle("GET "+base+"/groups/{group}/devices", h(a.listGroupDevices))
internal/appengine/http.go:72:	mux.Handle("POST "+base+"/groups/{group}/devices", h(a.addGroupDevice))
--- APICompat ---
588:const APICompatVersion = "1.2.2"

$ grep -rn "async_operation" internal/ --include=*.go | head -5; echo "--- required/encrypted mapping ---"; grep -rn "encrypted\|\"required\"" internal/interfaceschema/parse.go | head -8; echo "--- 082 keyAgreement deny ---"; grep -rn "keyAgreement" internal/broker/aclhook.go | head -3
internal/realm/async_operation_test.go:17:// with `?async_operation=false`. Astrate is always synchronous, so the
internal/realm/async_operation_test.go:24:	for i, q := range []string{"?async_operation=false", "?async_operation=true"} {
internal/housekeeping/async_operation_test.go:12:// asynchronously unless the caller opts out with `?async_operation=false`.
internal/housekeeping/async_operation_test.go:19:	for _, q := range []string{"?async_operation=false", "?async_operation=true"} {
internal/swagger/swagger_test.go:200:// spec tells clients the `?async_operation` parameter exists. Astrate accepts
--- required/encrypted mapping ---
--- 082 keyAgreement deny ---
72:		// control/keyAgreement is deliberately denied: full key agreement is
grep: internal/interfaceschema/parse.go: No such file or directory

$ git status --porcelain | head; echo "exit=$?"; wc -l .mule/todo.md .mule/for-giulio.md
exit=0
   298 .mule/todo.md
   657 .mule/for-giulio.md
   955 total

**Recipe carried out: `.mule/recipes/milestones.md`, target = v3.0** (v2.0 is `DONE`; scope is decided, so step 2 does not apply).

Step 3 investigation, measured this run:

- **Reference**: `.mule/research/upstream-parity-2026-08-22.md` (`.mule/milestones.md:114`) does not exist — `.mule/.gitignore` has a bare `research/`, so it can never be in the tree. Already escalated twice (`.mule/for-giulio.md:103`, `:395`), not re-escalated. I used the milestone's other named reference instead: the backlog itself.
- **Ledger**: `gh issue list --label milestone-v3.0 --state all` → empty (label `milestone-3.0` exists, carries nothing). Under `upstream-parity`: **#47–#89 all closed, none `not_planned`**; open = **#92** keyAgreement (parked on stable v1.4.0), **#78** FDO (`milestone-4.0`, parked), **#93** aclhook comment (`mule-review`). #90 closed with its frozen decision.
- **Release sweep**: no stable v1.4.0 — newest stable `v1.3.5` (2026-10-05), newest overall prerelease `v1.4.0-rc.6`. #92's parking condition unmet.
- **Standing item**: `docs/UPSTREAM-EXPERIMENTAL.md` still exactly its two rows (#67, #68), both `1.4 experimental`, trigger "upstream 1.4 final" — not final, so nothing to promote/deprecate. `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`).
- **Code coverage spot-check** (re-verified, not trusted): retention ceiling (`internal/store/retention.go`), `?detailed=true` (`internal/realm/detailed.go` + tests), `merge-patch+json` (`internal/appengine/http.go:210,231,522`), `GET /groups/{group}` (`http.go:70`), `async_operation` accepted-and-sync with tests, keyAgreement deny comment already rewritten per #93. Upstream master's only unmeasured commit (`171c1bd7`) is VerneMQ `/healtz` + producer sharding — infrastructure, no Astrate wire surface; already measured by yesterday's astarte-upstream run.

**Step 4 buckets**: nothing lands in *file now* — every remaining item is already an open issue (#92, #78, #93) or a closed decided one (#90), so filing would duplicate. Nothing new lands in *needs a decision* either — keyAgreement, the register rows, the tag cut and the wording/status fixes in `.mule/milestones.md` are all already on the queue.

**Step 5 deliberately not taken**, as in the 20 runs before: the label check is empty, but the milestone's scope is 1.2.2 → 1.3/**1.4** and the 1.4 half is open (#92 unanswered, both rows unreconciled, version unbumped) — proposing "looks complete, cut the tag" would be false.

```
Done: milestones recipe run for v3.0 — no new gaps; filed no issues and proposed no todo.md lines (all candidates already open or decided)
Files: none (working tree unchanged)
Verified: gh issue list --label milestone-{3.0,upstream-parity} --state all; gh api releases/tags; grep -n 'APICompatVersion =' internal/realm/service.go; coverage greps for #54/#58/#66/#72 -> pass, nothing new found
Unsure: the missing source doc means the 1.3/1.4 surface list was taken from the #47–#89 backlog rather than re-derived from upstream; "1.3 surface fully delivered" is prior runs' conclusion, spot-checked not exhaustively re-audited
Follow-ups: none
```

$ echo "=== releases ==="; gh api repos/astarte-platform/astarte/releases --jq '.[] | "\(.tag_name)\t\(.prerelease)\t\(.published_at)"' 2>&1 | head -15; echo "=== tags ==="; gh api repos/astarte-platform/astarte/tags --jq '.[].name' 2>&1 | head -10
=== releases ===
v1.3.5	false	2026-10-05T14:11:44Z
v1.4.0-rc.6	true	2026-09-30T07:28:27Z
v1.3.4	false	2026-09-18T09:42:57Z
v1.4.0-rc.5	true	2026-08-20T12:28:04Z
v1.3.3	false	2026-08-07T14:33:58Z
v1.4.0-rc.4	true	2026-08-06T13:40:20Z
v1.4.0-rc.3	true	2026-07-31T10:48:39Z
v1.4.0-rc.2	true	2026-07-14T10:16:22Z
v1.4.0-rc.1	true	2026-07-14T08:53:27Z
v1.3.2	false	2026-07-14T10:07:26Z
v1.3.1	false	2026-07-13T11:50:05Z
v1.3.0	false	2026-05-06T15:30:50Z
v1.2.2	false	2026-04-27T16:39:28Z
v1.4.0-rc.0	true	2026-04-08T12:32:33Z
v1.3.0-rc.2	true	2026-04-08T16:34:59Z
=== tags ===
v1.4.0-rc.6
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
v1.4.0-rc.0
v1.3.5
v1.3.4
v1.3.3

$ echo "=== milestone-3.0 ==="; gh issue list --label milestone-3.0 --state all --limit 50; echo "=== open upstream-parity ==="; gh issue list --label upstream-parity --state open --limit 50; echo "=== open milestone-4.0 ==="; gh issue list --label milestone-4.0 --state open --limit 20
=== milestone-3.0 ===
=== open upstream-parity ===
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
=== open milestone-4.0 ===
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z

→ Read .mule/for-giulio.md [limit=70]

$ echo "=== APICompatVersion ==="; rg -n 'APICompatVersion = ' internal/realm/service.go; echo "=== register rows ==="; rg -n '1.4 experimental' docs/UPSTREAM-EXPERIMENTAL.md; echo "=== issue 93 ancestor check ==="; git merge-base --is-ancestor 24ad5b8 origin/main 2>&1 && echo "on main" || echo "NOT on main"; echo "=== divergence ==="; git rev-list --left-right --count origin/main...HEAD 2>&1
=== APICompatVersion ===
588:const APICompatVersion = "1.2.2"
=== register rows ===
24:| Required + encrypted interface mapping fields | #67 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
25:| `async_operation=false` on mutating RM endpoints | #68 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
=== issue 93 ancestor check ===
NOT on main
=== divergence ===
4	746

$ echo "=== master commits since rc.6 (2026-09-30) ==="; gh api "repos/astarte-platform/astarte/commits?since=2026-09-30T07:28:27Z&per_page=100" --jq '.[] | "\(.sha[0:8])\t\(.commit.committer.date)\t\(.commit.message | split("\n")[0])"' 2>&1 | head -40
=== master commits since rc.6 (2026-09-30) ===
171c1bd7	2026-10-07T13:45:50Z	chore: forward-port astarte_vmq_plugin from release-1.4
1da04832	2026-10-06T11:01:17Z	chore: forward port release-1.4
a8a26d0f	2026-10-06T10:59:25Z	chore: forward port release-1.3
d85e0ed0	2026-10-06T10:02:32Z	docs: sync FDO workflow to the current behavior (#2241)
a6bd2c21	2026-10-06T09:59:16Z	fix(fdo): reject owner key names not accepted by OpenBao (#2252)
a51ab2e1	2026-10-06T07:21:18Z	refactor: Update Umbrella checks (#2244)
7628d2ed	2026-10-05T15:54:31Z	test(generators): validate interfaces with Core changesets
ed1e99f4	2026-10-05T15:54:31Z	test(generators): validate mappings with Core changesets
bb17455a	2026-10-05T15:54:31Z	test(generators): validate policy handlers with Core changesets
ec020bf8	2026-10-05T15:54:31Z	fix(generators): generate valid trigger policy names
8f0f5ef1	2026-10-05T14:09:21Z	chore: prepare v1.3.5 (#2245)
c4c63ebe	2026-10-05T09:55:53Z	chore(rm): bump xandra to 0.19.4 (#2246)
731bcddc	2026-10-05T09:34:08Z	chore: forward port release-1.4
4c4222cf	2026-10-05T09:02:23Z	chore: update horde
150e15d2	2026-10-02T13:03:27Z	build(deps): update horde and mississippi (#2240)
61b2a450	2026-10-02T07:52:45Z	chore(housekeeping): remove unused fields from to2 sessions (#2238)
b6d46ad4	2026-10-02T07:49:35Z	feat(appengine): include unexpected keys in object key errors (#2237)
5fc16d40	2026-10-01T15:04:48Z	chore: forward port release-1.4
958bb9fe	2026-10-01T15:03:27Z	feat: implement the device_empty_cache_received trigger (#2232)
b0d1ca71	2026-10-01T15:02:14Z	build: allow credo and format from umbrella root (#2233)
90efd5ee	2026-10-01T15:01:39Z	chore: remove unused labeler config (#2235)
acb5b379	2026-10-01T10:45:16Z	fleet: update libs paths (#2231)
b3c13d98	2026-10-01T10:44:59Z	ci: update paths in Release please (#2230)
5ae2b936	2026-09-30T07:54:22Z	refactor: Update libs paths in tools (#2229)

I read `MULE.md`, the `milestones.md` recipe, `.mule/milestones.md` (current target = **v3.0**; v2.0 is DONE), and carried out the recipe's step-3 investigation. Everything below was re-measured this run, not inherited:

- **Reference**: `.mule/research/upstream-parity-2026-08-22.md` is still absent (`.mule/.gitignore` has a bare `research/`, so it can never be in the tree) — already escalated twice, so not re-escalated. Used the milestone's other named reference, the backlog itself.
- **Ledger**: `gh issue list --label milestone-3.0 --state all` → empty (label exists, carries nothing; this milestone's ledger lives under `upstream-parity`). Open there: **#93** (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`). `#47–#89` otherwise all closed.
- **Release gate**: no stable v1.4.0 — newest stable `v1.3.5` (2026-10-05), newest overall `v1.4.0-rc.6` (prerelease). #92's parking condition (Giulio's "wait for a stable v1.4.0") still unmet.
- **Standing item**: `docs/UPSTREAM-EXPERIMENTAL.md` still exactly two rows (#67, #68), both `1.4 experimental`, trigger "upstream 1.4 final" — not final, nothing to reconcile. `APICompatVersion` still `"1.2.2"` (`internal/realm/service.go:588`).
- **Upstream master delta**: newest commit is still `171c1bd7` (2026-10-07, VerneMQ `/healtz` forward-port) — infrastructure, no Astrate wire surface; nothing new since the 2026-10-08 run.

No machine-checkable gap, no design decision un-escalated, and every remaining item already has an open owner — so filing an issue or queuing a task line would only duplicate #92/#93/#78. I also did not propose the step-5 "looks complete, cut the tag" line: the section's declared scope is 1.2.2 → 1.3/**1.4**, and the 1.4 half (keyAgreement, both register rows, the version bump) is unresolved.

```
Done: milestones recipe run for v3.0 — verified unchanged, still gated on a stable upstream v1.4.0; no new gaps, nothing filed, no task lines added
Files: none (working tree unchanged)
Verified: gh issue list --label milestone-3.0 / upstream-parity --state all; gh api repos/astarte-platform/astarte/releases (newest stable v1.3.5) and /commits since rc.6; rg 'APICompatVersion =' internal/realm/service.go; rg '1.4 experimental' docs/UPSTREAM-EXPERIMENTAL.md -> all consistent with prior runs, nothing new
Unsure: the missing source doc means the 1.3/1.4 surface list is taken from the #47–#89 backlog rather than re-derived from upstream ("1.3 fully delivered" is prior runs' conclusion, spot-checked here)
Follow-ups: none
```
