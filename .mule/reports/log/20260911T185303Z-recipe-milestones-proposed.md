slug: recipe-milestones
verdict: proposed
at:  20af0ea
ran: 2026-09-11T18:53:03Z on DietPi in 167s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
This is the milestone recipe — first non-DONE section is **v3.0 (upstream parity)**. Its status names a first job: triage #47–#89, escalate the "decide" batch, file sub-issues where work splits. Let me gather the sources.
→ Read docs/UPSTREAM-EXPERIMENTAL.md

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

$ gh issue list --label upstream-parity --state all --limit 60 --json number,title,labels,state --jq '.[] | [.number, .state, ([.labels[].name]|join(",")), .title] | @tsv'
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

✱ Glob "**/upstream-parity*.md" 0 matches
The research doc isn't where the milestone says. Let me find it and check the actual state of the backlog (#47–#89 look mostly CLOSED), plus the `milestone-3.0` label.
→ Read .mule/todo.md

$ gh issue list --label milestone-3.0 --state all --limit 50 --json number,title,labels,state --jq '.[] | [.number, .state, ([.labels[].name]|join(",")), .title] | @tsv'
(no output)

$ ls -la .mule/research/ 2>/dev/null; echo "---"; ls .mule/; echo "---"; ls docs/handoff/ 2>/dev/null
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
---
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

Most of the #47–#89 backlog is already CLOSED, `milestone-3.0` has no issues at all, and the research doc is missing from its stated path. Let me inspect the research dir, the open backlog items, and the current for-giulio escalations.
→ Read .mule/for-giulio.md

$ gh issue view 92 --json number,title,state,labels,body,comments --jq '"#\(.number) \(.state) [\([.labels[].name]|join(","))]\n\(.title)\n---\n\(.body)\n---COMMENTS---\n\([.comments[].body]|join("\n---\n"))"'
#92 OPEN [enhancement,upstream-parity,upstream-experimental]
keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds
---
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
---COMMENTS---
Parked by Giulio on 2026-09-04: **wait for a stable \`v1.4.0\`.** The spec is complete but ships only in \`v1.4.0-rc.5\`, and building against an rc risks rebuilding it.

So this issue is not waiting on anyone to remember it. It has a row in [`.mule/waiting-on.md`](https://github.com/astrate-platform/astrate/blob/main/.mule/waiting-on.md), and the weekly `mule-upstream-watch` job now reads that file as its first step every run: when a stable `v1.4.0` tag exists upstream, it escalates to `.mule/for-giulio.md` naming this issue. An `-rc.N` does not satisfy the row.

That mechanism was added because it was missing here: #51 was parked on the same condition, the condition fired on 2026-08-31, and nothing connected the new upstream tags to the parked issue — it took a hand triage four days later to catch it.

$ gh issue view 93 --json number,title,state,labels,body,comments --jq '"#\(.number) \(.state) [\([.labels[].name]|join(","))]\n\(.title)\n---\n\(.body)\n---COMMENTS---\n\([.comments[].body]|join("\n---\n"))"'
#93 OPEN [mule-review,upstream-parity]
In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
---
Comment-only change. The deny itself stays; only its stated justification is stale.

`internal/broker/aclhook.go` currently reads:

```go
// control/keyAgreement (080-mqtt-v1-protocol.md:172-176) is deliberately
// unsupported: upstream itself documents the handshake as "reserved and
// routed correctly, but not yet implemented", so falling through to the
// deny path below matches upstream's own status rather than being a gap.
```

That reasoning no longer holds on either half. Verified against upstream `master`:

- The quoted sentence, and the companion "accepted and acknowledged by Astarte, but no response is sent until the feature is fully implemented", are **gone from both** `080-mqtt-v1-protocol.md` and the new `082-key_agreement_protocol.md` (grep count 0 in each).
- `082-key_agreement_protocol.md` (13,832 bytes, added by upstream `d084308` on 2026-08-31) states the handshake "is implemented on the Astarte side (`astarte_data_updater_plant`)".
- The 080 topic table row also changed shape: a single `control/keyAgreement` at QoS 1 became `control/keyAgreement/<message id>` at QoS 2.

So the comment now tells a reader something false about upstream, and a future reader could reasonably "fix the gap" on the strength of it.

**Do:** rewrite the comment so it stands on Astrate's own decision — the topic is denied because full key agreement is out of scope for the current 1.2.2 target and the spec ships only in `v1.4.0-rc.5` (`v1.3.3` is the newest stable tag) — cite `082-key_agreement_protocol.md` as the current upstream reference rather than the deleted 080 sentence, and point at issue #92 where the implement-or-keep-parking decision lives.

**Do not:** change the ACL behaviour, and do not touch the expectation at `internal/broker/aclhook_test.go:53` (`pub control keyAgreement` still denies).

**Acceptance:** `go build ./...` and `go test ./internal/broker/...` pass unchanged; the comment no longer contains the phrase "not yet implemented"; `aclhook_test.go:53` is untouched.

Survey source: `survey-2026-08-31.md` on `mule/research`.
---COMMENTS---
The mule pushed `8c61268` to `mule/queue` for this.

Unreviewed: the gates passed — it compiles, the tests pass, a new test was shown to fail without the change, and no frozen file was touched. That is not the same as the change being worth having. Nothing merges to `main` until someone reads it:

    bash tools/mule.sh review

Left open on purpose. Whether this actually resolves the issue is your call, not mine.

This is a recurring run of the same proposal recipe — prior runs (2026-09-06/09/10) already triaged the #47–#89 backlog to CLOSED and escalated the open decision set. My job now is to verify the state is still current (mainly: has upstream shipped stable v1.4.0 yet?) and record a one-line finding. Let me re-check upstream and the open set.
$ gh api repos/astarte-platform/astarte/releases --jq '.[].tag_name' | head -20
v1.4.0-rc.5
v1.3.3
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
v1.3.2
v1.3.1
v1.3.0
v1.2.2
v1.4.0-rc.0
v1.3.0-rc.2
v1.2.2-rc.0
v1.2.1
v1.2.1-rc.1
v1.3.0-rc.1
v1.3.0-rc.0
v1.2.1-rc.0
v1.2.1-alpha.0
v1.2.0

$ echo "OPEN upstream-parity:"; gh issue list --label upstream-parity --state open --limit 20 --json number,title --jq '.[] | "#\(.number) \(.title)"'; echo "OPEN milestone-3.0:"; gh issue list --label milestone-3.0 --state open --limit 20 --json number,title --jq 'length'
OPEN upstream-parity:
#93 In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
#92 keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds
#78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
OPEN milestone-3.0:
0

State confirmed identical to the 2026-09-10 run: upstream newest stable is still `v1.3.3`, no stable `v1.4.0` (newest is `rc.5`); `milestone-3.0` is empty; the only open `upstream-parity` items are the already-escalated #92 (decision) and #93 (mule-review), plus #78 (milestone-4.0 FDO, parked). No new gaps, no issues to file, and no "cut the tag" proposal — the 1.4 gate is still open. The deliverable is one for-giulio line.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,10 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+
 - **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
 
 - **Extension point: external-bus Intake (NATS JetStream), 2026-09-10.** `internal/broker/intake.go:56-62` carries a `TODO(extension point)` for a second `Intake` implementation backed by a durable external bus (NATS JetStream) for multi-instance deployment or restart survival. The frozen design decision is in `docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4`; the consumer must reproduce per-device ordering, deferred-ack backpressure, and QoS 0 drop semantics. Not tracked by any open issue. Your call: file a milestone issue or park until the roadmap calls for it.
 

Done: v3.0 milestone verified unchanged, one `for-giulio.md` line recorded, nothing filed.
Files: .mule/for-giulio.md
Verified: `gh api repos/astarte-platform/astarte/releases` (no stable v1.4.0), `gh issue list --label milestone-3.0` (0), `gh issue list --label upstream-parity --state open` (#92/#93/#78) -> no new state vs 2026-09-10; no gate needed
Unsure: nothing
Follow-ups: none (no gaps to file, no "cut the tag" proposal — 1.4 gate open)
