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

