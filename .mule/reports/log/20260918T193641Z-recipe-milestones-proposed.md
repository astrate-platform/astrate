slug: recipe-milestones
verdict: proposed
at:  eb95cdd
ran: 2026-09-18T19:36:42Z on DietPi in 176s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
The first non-DONE milestone is **v3.0**. Let me gather the reference material and current state.
→ Read .mule/for-giulio.md
→ Read .mule/todo.md
v2.0 is DONE; the current target is **v3.0** (upstream parity 1.2.2 → 1.3/1.4), scope decided. Let me investigate per the recipe.

$ gh issue list --label upstream-parity --state all --limit 60 --json number,title,state,labels --jq '.[] | "\(.number)\t\(.state)\t\([.labels[].name]|join(","))\t\(.title)"'
93	OPEN	mule-review,upstream-parity	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
92	OPEN	enhancement,upstream-parity,upstream-experimental	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds
91	CLOSED	mule-blocked,upstream-parity	Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health
89	CLOSED	enhancement,upstream-parity	Dashboard flow-block schema mismatch (split_map/virtual pools hardcoded; null_sink/log_sink unknown)
88	CLOSED	enhancement,upstream-parity	Flow auth: support a_f JWT claim
87	CLOSED	enhancement,upstream-parity	Flow block: lua_map — needs embedded Lua runtime (parked)
86	CLOSED	enhancement,upstream-parity	Flow: pipeline source DSL — keep DAG-JSON as documented deviation?
85	CLOSED	enhancement,upstream-parity	Flow API: user-defined composite blocks
84	CLOSED	enhancement,upstream-parity	Flow blocks: virtual_device_pool / dynamic_virtual_device_pool
83	CLOSED	enhancement,upstream-parity	Flow blocks: mqtt_source/mqtt_sink (+modbus_tcp_source?) demand-driven
82	CLOSED	enhancement,upstream-parity	Flow blocks: http_source/http_sink (demand-driven)
81	CLOSED	enhancement,upstream-parity	Flow block: json_path_map
80	CLOSED	enhancement,upstream-parity	Flow blocks: pure-transform set (to_json, update_metadata, split_map, random_source, sort)
79	CLOSED	enhancement,upstream-parity	Verify registration-limit-reached HTTP status vs upstream
78	OPEN	enhancement,milestone-4.0,upstream-parity	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
77	CLOSED	enhancement,upstream-parity	Verify per-service version endpoints (GET /version, GET /v1/{realm}/version) served everywhere
76	CLOSED	enhancement,upstream-parity	Housekeeping: GET /v1/realm-defaults/replication — decide reject/deviate (Cassandra-shaped)
75	CLOSED	enhancement,upstream-parity	Housekeeping: decide realm-deletion gating/preconditions vs always-sync deviation
74	CLOSED	enhancement,upstream-parity	Housekeeping: PATCH /v1/realms/{realm} (jwt key, registration limit, retention; null=unset)
73	CLOSED	enhancement,upstream-parity	Housekeeping: default datastream retention injection env var (upstream 1.4)
72	CLOSED	enhancement,upstream-parity	Realms: datastream_maximum_storage_retention ceiling (create/patch/enforce)
71	CLOSED	enhancement,upstream-parity	Pairing: realm-scoped health check GET /v1/{realm}/health (upstream 1.3)
70	CLOSED	enhancement,upstream-parity	Triggers: audit wildcard semantics (interface_name '*', match_path '/*' forcing rules)
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
