slug: recipe-milestones
verdict: proposed
at:  a3df2d0
ran: 2026-09-26T11:31:38Z on DietPi in 197s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md

$ ls -la /root/astrate-mule/.mule/ 2>&1
total 8052
drwxr-xr-x  7 root root    4096 Sep 26 13:27 .
drwxr-xr-x 17 root root    4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       0 Sep 26 13:05 .alarmed
-rw-r--r--  1 root root      13 Sep 26 13:28 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 25 22:19 .heartbeat
-rw-r--r--  1 root root       0 Sep 26 13:28 .last-output
-rw-r--r--  1 root root       0 Sep 26 13:28 .lock
-rw-r--r--  1 root root      11 Sep 26 13:28 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 7979047 Sep 26 13:28 cron.log
drwxr-xr-x  2 root root    4096 Sep 25 22:45 failed
-rw-r--r--  1 root root   67672 Sep 25 21:25 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root   10124 Sep 25 22:45 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 25 19:46 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   92699 Sep 25 22:45 todo.md

v2.0 is DONE, so v3.0 (upstream Astarte feature parity, 1.2.2 → 1.3/1.4) is the target. Gathering the inputs the recipe names.
$ cd /root/astrate-mule && ls .mule/research/ 2>&1; echo "=== docs ==="; ls docs/ 2>&1; echo "=== UPSTREAM-EXPERIMENTAL ==="; wc -l docs/UPSTREAM-EXPERIMENTAL.md 2>&1
ls: cannot access '.mule/research/': No such file or directory
=== docs ===
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
=== UPSTREAM-EXPERIMENTAL ===
25 docs/UPSTREAM-EXPERIMENTAL.md

$ cd /root/astrate-mule && gh issue list --label milestone-v3.0 --state all --limit 50 2>&1 | head -60; echo "=== ALL LABELS ==="; gh label list --limit 60 2>&1 | head -60
=== ALL LABELS ===
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

$ cd /root/astrate-mule && find . -path ./.git -prune -o -name "*upstream-parity*" -print 2>&1; echo "=== research dirs ==="; find . -path ./.git -prune -o -type d -name "research" -print 2>&1; echo "=== handoff ==="; ls docs/handoff/ 2>&1
=== research dirs ===
=== handoff ===
README.md
flow-design-a-named-durable-flows-2026-07-29.md
flow-design-b-container-block-2026-07-29.md
flow-lifecycle-followups-design-2026-08-22.md
flow-parity-audit-2026-07-29.md
flow-v2-decisions-2026-07-29.md
phase-1-handoff.md
phase-1-memory.md
phase-10-handoff.md
phase-10-memory.md
phase-11-handoff.md
phase-2-handoff.md
phase-2-memory.md
phase-3-handoff.md
phase-3-memory.md
phase-4-handoff.md
phase-4-memory.md
phase-5-handoff.md
phase-5-memory.md
phase-6-handoff.md
phase-6-memory.md
phase-7-memory.md
phase-8-handoff.md
phase-8-memory.md
phase-9-handoff.md
phase-9-memory.md
pokemon-agent-handoff.md
pokemon-agent-memory.md
session-2026-07-29-flow-factory-handoff.md
session-2026-07-29-handoff.md

$ cd /root/astrate-mule && echo "=== milestone-v3.0 all ==="; gh issue list --label milestone-v3.0 --state all --limit 50 2>&1; echo "=== count ==="; gh issue list --label milestone-v3.0 --state all --limit 50 --json number 2>&1 | head -5; echo "=== upstream-parity open ==="; gh issue list --label upstream-parity --state all --limit 60 --json number,title,state --limit 60 2>&1 | head -200
=== milestone-v3.0 all ===
=== count ===
[]
=== upstream-parity open ===
[{"number":93,"state":"OPEN","title":"In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged"},{"number":92,"state":"OPEN","title":"keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds"},{"number":91,"state":"CLOSED","title":"Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health"},{"number":89,"state":"CLOSED","title":"Dashboard flow-block schema mismatch (split_map/virtual pools hardcoded; null_sink/log_sink unknown)"},{"number":88,"state":"CLOSED","title":"Flow auth: support a_f JWT claim"},{"number":87,"state":"CLOSED","title":"Flow block: lua_map — needs embedded Lua runtime (parked)"},{"number":86,"state":"CLOSED","title":"Flow: pipeline source DSL — keep DAG-JSON as documented deviation?"},{"number":85,"state":"CLOSED","title":"Flow API: user-defined composite blocks"},{"number":84,"state":"CLOSED","title":"Flow blocks: virtual_device_pool / dynamic_virtual_device_pool"},{"number":83,"state":"CLOSED","title":"Flow blocks: mqtt_source/mqtt_sink (+modbus_tcp_source?) demand-driven"},{"number":82,"state":"CLOSED","title":"Flow blocks: http_source/http_sink (demand-driven)"},{"number":81,"state":"CLOSED","title":"Flow block: json_path_map"},{"number":80,"state":"CLOSED","title":"Flow blocks: pure-transform set (to_json, update_metadata, split_map, random_source, sort)"},{"number":79,"state":"CLOSED","title":"Verify registration-limit-reached HTTP status vs upstream"},{"number":78,"state":"OPEN","title":"FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)"},{"number":77,"state":"CLOSED","title":"Verify per-service version endpoints (GET /version, GET /v1/{realm}/version) served everywhere"},{"number":76,"state":"CLOSED","title":"Housekeeping: GET /v1/realm-defaults/replication — decide reject/deviate (Cassandra-shaped)"},{"number":75,"state":"CLOSED","title":"Housekeeping: decide realm-deletion gating/preconditions vs always-sync deviation"},{"number":74,"state":"CLOSED","title":"Housekeeping: PATCH /v1/realms/{realm} (jwt key, registration limit, retention; null=unset)"},{"number":73,"state":"CLOSED","title":"Housekeeping: default datastream retention injection env var (upstream 1.4)"},{"number":72,"state":"CLOSED","title":"Realms: datastream_maximum_storage_retention ceiling (create/patch/enforce)"},{"number":71,"state":"CLOSED","title":"Pairing: realm-scoped health check GET /v1/{realm}/health (upstream 1.3)"},{"number":70,"state":"CLOSED","title":"Triggers: audit wildcard semantics (interface_name '*', match_path '/*' forcing rules)"},{"number":69,"state":"CLOSED","title":"Verify unknown-realm HTTP status on RM endpoints against upstream"},{"number":68,"state":"CLOSED","title":"Decide async_operation=false params vs documented always-sync deviation"},{"number":67,"state":"CLOSED","title":"Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4)"},{"number":66,"state":"CLOSED","title":"Realm Management: detailed=true interface listing with full mappings (upstream 1.4)"},{"number":65,"state":"CLOSED","title":"Policies: handler-overlap rejection + retry_times coupling + prefetch_count"},{"number":64,"state":"CLOSED","title":"Triggers: decide AMQP action behavior (validate-reject vs NATS-forward deviation)"},{"number":63,"state":"CLOSED","title":"Triggers: HTTP action validation limits (URL/method/header blocklist/template size)"},{"number":62,"state":"CLOSED","title":"Realm Management: audit install/update/delete error codes and statuses"},{"number":61,"state":"CLOSED","title":"Realm Management: audit interface/mapping validation matrix against astarte_core"},{"number":60,"state":"CLOSED","title":"Realm Management: GET config/datastream_maximum_storage_retention (since upstream 1.2.0)"},{"number":59,"state":"CLOSED","title":"AppEngine: group create-body validation + UUID-v1 from_token for group device listing"},{"number":58,"state":"CLOSED","title":"AppEngine: PATCH requires Content-Type application/merge-patch+json"},{"number":57,"state":"CLOSED","title":"AppEngine: audit server-write error taxonomy against upstream"},{"number":56,"state":"CLOSED","title":"AppEngine: object-aggregate response formats table and disjoint_tables"},{"number":55,"state":"CLOSED","title":"AppEngine: complete data-query parameter set"},{"number":54,"state":"CLOSED","title":"AppEngine: GET /groups/{g} show-group endpoint"},{"number":53,"state":"CLOSED","title":"AppEngine: group-scoped device PATCH + interface data access under groups/{g}/devices/{d}"},{"number":52,"state":"CLOSED","title":"AppEngine: devices-by-alias full mirror (PATCH via alias + interface data access)"},{"number":51,"state":"CLOSED","title":"keyAgreement control topic + encrypted mappings (upstream 1.4 experimental) — investigate/park"},{"number":50,"state":"CLOSED","title":"Broker/engine hardening: discard malformed introspection/purge/capabilities instead of crashing"},{"number":49,"state":"CLOSED","title":"MQTT v1: binaryblob properties re-send as BSON subtype 0 in consumer/properties"},{"number":48,"state":"CLOSED","title":"MQTT v1: tolerate devices connecting with an empty introspection"},{"number":47,"state":"CLOSED","title":"MQTT v1: accept device capabilities incl. purge_properties_compression_format (plaintext purge)"}]

Important: the source doc `.mule/research/upstream-parity-2026-08-22.md` is **not in the repo**, and #47–#89 are nearly all CLOSED. Let me establish the real state before concluding.
→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ cd /root/astrate-mule && echo "=== #90 ==="; gh issue view 90 2>&1 | head -60; echo; echo "=== #92 ==="; gh issue view 92 2>&1 | head -50
=== #90 ===
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

=== #92 ===
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

$ cd /root/astrate-mule && echo "=== APICompatVersion ==="; rg -n "APICompatVersion" --glob '!*.log' -g '!.git' . 2>&1 | head -20; echo; echo "=== for-giulio v3.0 mentions ==="; rg -n "v3\.0|milestone.3\.0" .mule/for-giulio.md 2>&1 | tail -30
=== APICompatVersion ===
./cmd/astrate/main.go:437:		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(realm.APICompatVersion)))
./cmd/astrate/main.go:438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
./cmd/astrate/version_test.go:25:	want := `{"data":"` + realm.APICompatVersion + `"}`
./docs/COMPATIBILITY.md:253:    reports `realm.APICompatVersion` (currently `1.2.2`), a compatibility
./docs/UPSTREAM-EXPERIMENTAL.md:13:`APICompatVersion` (internal/realm/service.go) must never claim a level higher
./internal/realm/http.go:119:// feature UI on it — see APICompatVersion).
./internal/realm/http.go:121:	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
./internal/realm/dashboard_compat_test.go:35:		if v != APICompatVersion {
./internal/realm/dashboard_compat_test.go:36:			t.Errorf("version = %q, want %q", v, APICompatVersion)
./internal/realm/service.go:580:// APICompatVersion is the upstream Realm Management API level Astrate
./internal/realm/service.go:588:const APICompatVersion = "1.2.2"

=== for-giulio v3.0 mentions ===
81:- **Milestone recipe run, 2026-09-24: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-23 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
89:- **Milestone recipe run, 2026-09-23: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-21 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
105:- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
113:- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
125:- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
129:- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
145:- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) shows a new **stable v1.3.4** (2026-09-18) — maintenance only (data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change), so no new gap for Astrate — and still **no stable v1.4.0**: newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
157:- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-13 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
161:- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-12 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
165:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
173:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
175:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
205:- **Milestone recipe run, 2026-09-09: v3.0 unchanged since 2026-09-06 — still waiting on
206:  upstream gates.** `milestone-3.0` label remains empty. The three open `upstream-parity`
226:- **Milestone recipe run, 2026-09-06: v3.0 (upstream 1.2.2 → 1.3/1.4) has no open
228:  `milestone-3.0` label is empty (this milestone's work lives under `upstream-parity`).** All
240:  only) — until then v3.0 has no open implementable work. This run filed **no issues** and
243:- **`milestones.md` v3.0 names a source doc that does not exist: `.mule/research/
333:  **Decided 2026-09-04: the tag is `v0.2.0`, cut on the same day.** The `v2.0`/`v3.0` names
357:  schema and keys) is a multi-session project. Parked until the v3.0 queue clears and this
425:When v3.0 is marked DONE, the v4.0 section of `.mule/milestones.md` should be

$ cd /root/astrate-mule && echo "=== #90 comment ==="; gh issue view 90 --comments 2>&1 | tail -30; echo; echo "=== #67 closure ==="; gh issue view 67 --comments 2>&1 | tail -25; echo; echo "=== #68 closure ==="; gh issue view 68 --comments 2>&1 | tail -25
=== #90 comment ===
author:	atsetilam
association:	member
edited:	false
status:	none
--
Decision frozen 2026-08-23: **do not bump now, and never make it configurable.**

- `APICompatVersion` stays `1.2.2` until Astrate implements the FULL 1.3 surface; it is bumped only in the same change that completes that surface, after reconciling every row in docs/UPSTREAM-EXPERIMENTAL.md tagged with the level.
- Configurability rejected: an operator raising the declared level would make clients believe features exist that don't — the costliest lie available.
- Enforcement is mechanical, not memory-based: the bump rule lives inline on the constant (internal/realm/service.go), as the standing checklist item in .mule/recipes/milestones.md, and in docs/UPSTREAM-EXPERIMENTAL.md's version-declaration rule.
--

=== #67 closure ===
author:	atsetilam
association:	member
edited:	false
status:	none
--
Not delegated to the mule (assessed 2026-08-31).

Both code phases have landed (4a schema layer \`9c4d411\`, 4b runtime enforcement of \`required\` \`7a244b2\`). What is left is the docs closeout — the \`docs/COMPATIBILITY.md\` row explaining why astrate’s \`missing_required\` reason is translated to \`unexpected_object_key\`, and the \`docs/UPSTREAM-EXPERIMENTAL.md\` Adopted dates — which is prose the mule does not write, and it lands together with #68's closeout in fase 4c. Keeping this issue out of the mule queue on purpose.

\`encrypted\` runtime behaviour remains tied to keyAgreement (#51) and is out of this issue’s scope.
--
author:	atsetilam
association:	member
edited:	false
status:	none
--
Decided and shipped on 2026-08-23, before this issue was revisited: `required` and `encrypted` are parsed, stored, carried through interface compatibility checks, and `required` is enforced at runtime — an object-aggregated document missing a required key is rejected (`missing_required`, translated to the closest 1.2.2 error_name for the dashboard's closed set). `encrypted` is parsed and stored only, pending the keyAgreement work.

Tracked in `docs/UPSTREAM-EXPERIMENTAL.md` for reconciliation when upstream 1.4 is final (promoted → keep and drop the row; removed → deprecate in the same milestone).
--

=== #68 closure ===
Design is already frozen in \`.trickle/plans/MASTER-HANDOFF.md\` fase 4c and is not up for relitigation: Astrate stays always-sync, and \`async_operation=false\` is accepted and ignored on the mutating endpoints upstream exposes it on — housekeeping realm create/delete; realm-management interface install/update/delete and trigger/policy delete (policy routes are mounted in \`internal/realm/http.go\`). An unparseable or \`true\` value must not change behaviour either.

**In scope for the mule:** the code and its tests only.
**Out of scope:** \`docs/COMPATIBILITY.md\` (deviation 17), \`docs/UPSTREAM-EXPERIMENTAL.md\`, and closing this issue — the docs closeout stays with the architect and lands once, together with #67.
--
author:	atsetilam
association:	member
edited:	false
status:	none
--
The mule could not do this: **gates failed**

Taken out of its queue so it does not retry the same failure every half hour. Re-label `mule` to queue it again, ideally after making the task smaller or the requirement clearer.
--
author:	atsetilam
association:	member
edited:	false
status:	none
--
Decision: **Astrate stays always-synchronous**, recorded as deliberate deviation 17 rather than matched.

Every one of the endpoints upstream 1.4 runs asynchronously — realm create/delete, interface install/update/delete, delivery-policy delete — completes before Astrate answers. That is a strictly stronger guarantee than upstream's default, so a client written against upstream never sees a regression.

`?async_operation=false` was already accepted and ignored, which is the wanted behaviour, but nothing recorded or enforced it. As of 34b7da7 it is pinned by `TestRealmManagementAsyncOperationParam` and `TestHousekeepingAsyncOperationParam` (both parameter values, all six endpoints, effect visible on the next read), so a later tightening of query validation cannot silently break upstream clients. Documented in `docs/COMPATIBILITY.md` and on the public compatibility page; the `docs/UPSTREAM-EXPERIMENTAL.md` row stays until upstream 1.4 is final.
--

→ Read .mule/for-giulio.md [limit=90]
This recipe has run ~15 times already with a consistent finding. Let me verify the state fresh rather than copy the prior run.

$ cd /root/astrate-mule && echo "=== upstream releases ==="; gh api repos/astarte-platform/astarte/releases --jq '.[] | "\(.tag_name)\t\(.published_at)\tprerelease=\(.prerelease)\t draft=\(.draft)"' 2>&1 | head -12; echo; echo "=== upstream tags (recent) ==="; gh api repos/astarte-platform/astarte/tags --jq '.[].name' 2>&1 | head -12
=== upstream releases ===
v1.3.4	2026-09-18T09:42:57Z	prerelease=false	 draft=false
v1.4.0-rc.5	2026-08-20T12:28:04Z	prerelease=true	 draft=false
v1.3.3	2026-08-07T14:33:58Z	prerelease=false	 draft=false
v1.4.0-rc.4	2026-08-06T13:40:20Z	prerelease=true	 draft=false
v1.4.0-rc.3	2026-07-31T10:48:39Z	prerelease=true	 draft=false
v1.4.0-rc.2	2026-07-14T10:16:22Z	prerelease=true	 draft=false
v1.4.0-rc.1	2026-07-14T08:53:27Z	prerelease=true	 draft=false
v1.3.2	2026-07-14T10:07:26Z	prerelease=false	 draft=false
v1.3.1	2026-07-13T11:50:05Z	prerelease=false	 draft=false
v1.3.0	2026-05-06T15:30:50Z	prerelease=false	 draft=false
v1.2.2	2026-04-27T16:39:28Z	prerelease=false	 draft=false
v1.4.0-rc.0	2026-04-08T12:32:33Z	prerelease=true	 draft=false

=== upstream tags (recent) ===
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
v1.4.0-rc.0
v1.3.4
v1.3.3
v1.3.2
v1.3.1
v1.3.0
v1.3.0-rc.2

$ cd /root/astrate-mule && echo "=== recent commits on main ==="; git log --oneline -12 2>&1; echo; echo "=== pending todo milestone lines ==="; rg -n "^- \[ \] milestone" .mule/todo.md 2>&1 | head -20; echo "=== todo unchecked count ==="; rg -c "^- \[ \]" .mule/todo.md 2>&1
=== recent commits on main ===
a3df2d0 mule: recipe hygiene timed out
85b5697 mule: blocked docs-sync-rm-delete-device-async-operation-param
9f20120 mule: log docs-sync-rm-delete-device-async-operation-param
f965f39 mule: docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
7d78bd3 mule: log docs-sync-rm-async-operation-param
3f0eca8 mule: docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
17bfa74 mule: log docs-sync-hk-retention-zero-is-unset
c2ad442 mule: docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
57501a7 mule: log docs-sync-hk-async-operation-param
ed1cd72 mule: docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
19f4788 mule: refill the queue
ee48763 mule: recipe docs-sync ran (335s)

=== pending todo milestone lines ===
=== todo unchecked count ===
12

$ cd /root/astrate-mule && echo "=== all open, non-alarm ==="; gh issue list --state open --limit 60 --json number,title,labels,updatedAt --jq '.[] | select((.labels|map(.name)|index("mule-alarm")) == null) | "\(.number)\t\(.updatedAt)\t\(.labels|map(.name)|join(","))\t\(.title[0:70])"' 2>&1; echo; echo "=== #93 state ==="; gh issue view 93 --json state,labels,updatedAt 2>&1
=== all open, non-alarm ===
93	2026-09-04T19:54:58Z	mule-review,upstream-parity	In internal/broker/aclhook.go the control/keyAgreement deny comment (l
92	2026-09-04T19:52:21Z	enhancement,upstream-parity,upstream-experimental	keyAgreement: upstream published the full wire spec (082) — decide whe
78	2026-08-23T16:16:08Z	enhancement,milestone-4.0,upstream-parity	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 ca
1	2026-03-13T21:55:11Z	wontfix	Provide an Open Source IoT Platform unironically, unironically

=== #93 state ===
{"labels":[{"id":"LA_kwDORmfsJs8AAAACte8Saw","name":"mule-review","description":"Mule pushed a change for this; waiting for a human or a strong model to review it","color":"0e8a16"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"}],"state":"OPEN","updatedAt":"2026-09-04T19:54:58Z"}

State confirmed unchanged. Checking whether prior runs already raised the two things I think are actually actionable, so I don't duplicate.
$ cd /root/astrate-mule && echo "=== count of milestone-run entries ==="; rg -c "Milestone recipe run" .mule/for-giulio.md; echo "=== distinct dates ==="; rg -o "Milestone recipe run, [0-9-]+" .mule/for-giulio.md | sort -u; echo; echo "=== was 'Status: not started' staleness ever raised? ==="; rg -n "not started" .mule/for-giulio.md | head; echo "=== was the repetition/pruning point ever raised? ==="; rg -n "queue, not a log|not a log|prune|sixteenth|repeat" .mule/for-giulio.md | head
=== count of milestone-run entries ===
14
=== distinct dates ===
Milestone recipe run, 2026-09-06
Milestone recipe run, 2026-09-09
Milestone recipe run, 2026-09-10
Milestone recipe run, 2026-09-11
Milestone recipe run, 2026-09-12
Milestone recipe run, 2026-09-13
Milestone recipe run, 2026-09-15
Milestone recipe run, 2026-09-18
Milestone recipe run, 2026-09-19
Milestone recipe run, 2026-09-20
Milestone recipe run, 2026-09-21
Milestone recipe run, 2026-09-23
Milestone recipe run, 2026-09-24

=== was 'Status: not started' staleness ever raised? ===
=== was the repetition/pruning point ever raised? ===
9:line once you have dealt with it — this file is a queue, not a log.

Verified fresh. No new gaps, so nothing to file and nothing to queue. The one genuinely new observation is that `milestones.md`'s own v3.0 status line is now materially false. Writing the escalation entry.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,15 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Status: not started", which is no longer true.** Release sweep (`gh api repos/astarte-platform/astarte/releases`, run this session) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.5` (2026-08-20, prerelease) — byte-identical to what the 2026-09-24 run saw, so nothing shipped upstream in the two days since. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Open non-alarm issues remain exactly **#93** (aclhook comment, `mule-review`, its own recipe path), **#92** (keyAgreement, parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** (FDO, `milestone-4.0`), **#1** (untouched per standing instruction). No new gaps, **no issues filed, no `.mule/todo.md` task lines added** — there is nothing machine-checkable left, and filing anything here would duplicate #92/#93/#78.
+  **New, and the reason this run is not a copy of the last thirteen: `.mule/milestones.md:127-128` still reads "Status: not started. First recipe job: triage #47–#89 into an ordered plan (which are audits vs features vs decisions), file sub-issues where work splits, escalate the 'decide' set in one batch."** All of that happened: of #47–#89 every issue is closed except the deliberately-parked #78. The section describes as future work a triage that is finished, and it is the exact text every future run of this recipe reads to pick its target — so the recipe has been re-deriving "the backlog is done, the gate is upstream" from scratch on a timer for three weeks. Suggest rewriting those two lines to state the real position: *1.3 surface delivered (#47–#89 closed); milestone blocked on upstream v1.4.0 going stable; final phase = reconcile the two `docs/UPSTREAM-EXPERIMENTAL.md` rows, answer #92, bump `APICompatVersion`.* That also makes the blocking condition greppable instead of inferred.
+  **The recipe's standing item is satisfied, and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still has exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental* with reconcile trigger "upstream 1.4 final". v1.4 is not final, so both rows stay as written: nothing to promote, nothing to deprecate, and the register is accurate. `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision; the milestone-bump issue the recipe asks for **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": the milestone's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is done, and the 1.4 half is unresolved (#92 open, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut. Also still true and already escalated on 2026-09-06, so not repeated: the section's source document `.mule/research/upstream-parity-2026-08-22.md` does not exist in the repo.
+  **Standing ask, unchanged, for whenever a stable v1.4.0 ships:** answer #92, reconcile the two experimental rows, run the final-phase `APICompatVersion` bump (and update `docs/COMPATIBILITY.md` deviation #10 to match), then cut the v3.0 tag. One process note while this stays blocked: this is the fifteenth milestone-run entry in this file (fourteen prior, 2026-09-06 → 2026-09-24) and the header asks for a queue rather than a log — worth deleting the superseded duplicates whenever you next prune, or the signal-to-noise here keeps falling.
+
+---
+
 ## 2026-09-25 — docs-sync run, surface: **housekeeping** (the least-covered of the five)
 
 Nothing filed, no `gh issue create` in this recipe. Two `docs/api/` proposals appended to
 `.mule/todo.md` (`docs-sync-hk-async-operation-param`, `docs-sync-hk-retention-zero-is-unset`);

$ cd /root/astrate-mule && sed -n '1,20p' .mule/for-giulio.md | cut -c1-180; echo "..."; echo "=== structure check: separators and headings ==="; rg -n "^---$|^## " .mule/for-giulio.md | head -12; echo "=== git status (must be no git ops, just confirm what changed) ==="; git status --porcelain 2>&1
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Stat
  **New, and the reason this run is not a copy of the last thirteen: `.mule/milestones.md:127-128` still reads "Status: not started. First recipe job: triage #47–#89 into an orde
  **The recipe's standing item is satisfied, and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still has exactly its two rows — #67 required/encrypted mapping fields an
  **Standing ask, unchanged, for whenever a stable v1.4.0 ships:** answer #92, reconcile the two experimental rows, run the final-phase `APICompatVersion` bump (and update `docs/CO

---

## 2026-09-25 — docs-sync run, surface: **housekeeping** (the least-covered of the five)
...
=== structure check: separators and headings ===
11:---
18:---
20:## 2026-09-25 — docs-sync run, surface: **housekeeping** (the least-covered of the five)
82:---
86:---
90:---
94:---
98:---
102:---
106:---
110:---
114:---
=== git status (must be no git ops, just confirm what changed) ===
 M .mule/for-giulio.md

Entry placed newest-at-top; only `.mule/for-giulio.md` is modified. No issues filed (none warranted), no task lines appended (no new gap), no git operations.

**What the recipe resolved to:** v2.0 is `DONE`, so the target is v3.0 (upstream parity 1.2.2 → 1.3/1.4). I re-verified the state this session rather than trusting the prior run: upstream still has no stable v1.4.0 (newest stable v1.3.4, newest overall v1.4.0-rc.5), the `milestone-3.0` label is still empty because the ledger lives under `upstream-parity`, and the only open non-alarm issues remain #92, #93, #78, #1. No new gaps, so nothing to file and nothing to queue.

**The one new finding:** `.mule/milestones.md:127-128` still says v3.0 is "Status: not started. First recipe job: triage #47–#89…". That triage is finished — every issue in #47–#89 is closed except the deliberately parked #78. So the file every future run reads to pick its target is materially wrong, which is why this recipe has been re-deriving the same conclusion for three weeks. I proposed the corrected wording rather than editing it, since that file is Giulio's.

Two things I deliberately did **not** do, both deliberate per the recipe and precedent:
- Did not propose "milestone complete, cut the tag" — the scope names 1.4, and the 1.4 half is unresolved (#92 open, both `UPSTREAM-EXPERIMENTAL.md` rows unreconciled, `APICompatVersion` still `1.2.2`).
- Did not file the standing `APICompatVersion` bump issue — **#90 already exists** (closed with Giulio's frozen decision). Filing again would be the duplicate the recipe explicitly warns about.

One caveat worth your attention: this is now the 15th near-identical milestone entry in `for-giulio.md`, whose own header asks for "a queue, not a log". The repetition is diluting the file. Pruning the superseded duplicates would help, but I left them alone since deleting Giulio's escalation queue isn't my call.

```
Done: Milestone recipe run — v3.0 re-verified unchanged and still blocked on a stable upstream v1.4.0; no issues filed, no task lines added; escalated the one new finding (milestones.md:127-128 still claims "Status: not started" for a triage that is complete) to .mule/for-giulio.md
Files: .mule/for-giulio.md (one entry added, newest at top)
Verified: gh api repos/astarte-platform/astarte/releases (no stable v1.4.0; newest stable v1.3.4 2026-09-18, newest overall v1.4.0-rc.5 2026-08-20) -> pass; gh issue list --label milestone-3.0 --state all (empty) -> pass; gh issue list --state open minus mule-alarm (#93, #92, #78, #1) -> pass; docs/UPSTREAM-EXPERIMENTAL.md read, 2 rows both 1.4-experimental, neither reconcilable until 1.4 final -> pass; rg APICompatVersion (still "1.2.2" at internal/realm/service.go:588) -> pass. No build gate run: no code was changed.
Unsure: whether #93's pushed commit 8c61268 is merged (prior runs note it is on mule/queue, not origin/main) — I did not chase it, it is #93's own recipe path, not this milestone's
Follow-ups: none appended to .mule/todo.md — deliberately, there is no new machine-checkable gap; the only proposal is the milestones.md:127-128 rewrite in for-giulio.md, which is Giulio's file to edit
```
