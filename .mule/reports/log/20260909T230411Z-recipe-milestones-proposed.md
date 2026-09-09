slug: recipe-milestones
verdict: proposed
at:  e8ebeb6
ran: 2026-09-09T23:04:11Z on DietPi in 214s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
The first non-DONE milestone is **v3.0** (upstream Astarte parity), scope decided on 2026-08-22. Its first job: triage issues #47–#89. Let me investigate.

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

$ gh issue list --label milestone-v3.0 --state all --limit 50
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

✱ Glob ".mule/research/*" 0 matches

$ gh issue list --label milestone-3.0 --state all --limit 50
(no output)

$ find /root/astrate-mule -iname "*parity*" -o -iname "*research*" 2>/dev/null | rg -v "\.git/"
/root/astrate-mule/internal/appengine/formats_parity_test.go
/root/astrate-mule/internal/appengine/groups_parity_test.go
/root/astrate-mule/internal/engine/mqttv1_parity_test.go
/root/astrate-mule/docs/handoff/flow-parity-audit-2026-07-29.md

$ git -C /root/astrate-mule status && git -C /root/astrate-mule branch -a
On branch mule/queue
nothing to commit, working tree clean
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

The research doc isn't on disk. Let me check the git state and available branches.
$ git -C /root/astrate-mule log --oneline -15
e8ebeb6 mule: recipe hygiene timed out
979872f mule: log docs-sync-native-compat-version-endpoints
d148f16 mule: docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
7369a0d mule: log docs-sync-native-compat-health-503
656e21e mule: docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmmanagement/health`, `GET /pairing/health` — in docs/api/astrate_native_api.yaml, and widen the shared `DataEnvelopeHealthStatus.status` enum (yaml:301-303) from `[ok]` to `[ok, unhealthy]` — the handler answers 503 `{"data":{"status":"unhealthy"}}` when the readiness probe fails (internal/observability/compat.go:19-34, deviation 18), but the spec documents only 200 with `status: ok` for all three paths. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
d1188bd mule: refill the queue
71555dc mule: recipe docs-sync ran (282s)
304c126 mule: log realm-pure-helper-tests
39eb70a mule: realm-pure-helper-tests [auto]: add container-free unit tests for the four untested pure helpers in internal/realm. Create `internal/realm/helpers_test.go` (no build tag) covering: (1) `normaliseIfaceName` — lowercase, hyphen strip, mixed case+hyphens, empty string, already-normalized passthrough; (2) `joinPEM` — empty slice, single key, two keys separated by `\n`; (3) `triggerErrorBody` — empty action+simple_triggers, action-only, simple_triggers-only with index-aligned `{}` substitution, both present; (4) `majorParam` — valid positive, zero, negative (→404), non-numeric (→404). All tests use httptest.ResponseRecorder for majorParam; the rest are direct function calls. No Docker.
84fe5df mule: refill the queue
2340774 mule: recipe code-review ran (364s)
a701b05 mule: recipe astarte-upstream ran (239s)
cf14e25 mule: recipe github-issues ran (96s)
1d57030 mule: refill the queue
6e1ab2e mule: recipe milestones ran (115s)

$ git -C /root/astrate-mule ls-tree -r --name-only origin/mule/research | rg -i "parity-2026|research/"
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

$ git -C /root/astrate-mule cat-file -e origin/mule/research:.mule/research/issues-2026-08-22.md && git -C /root/astrate-mule show origin/mule/research:.mule/research/issues-2026-08-22.md
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

$ git -C /root/astrate-mule log --all --oneline --name-only --diff-filter=A | rg -i -B2 "parity-2026-08-22|research"
.mule/reports/log/20260904T110302Z-pairing-info-verify-ratelimit-blocked.md
6d8cae6 mule: daily cross-project survey, 2026-09-04
.mule/research/issues-2026-09-04.md
.mule/research/survey-2026-09-04.md
--
.mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md
1c618eb mule: daily cross-project survey, 2026-08-31
.mule/research/survey-2026-08-31.md
ce759e1 survey: recover 16 daily reports the broken commit check left staged (2026-07-28..2026-08-26), and take the tooling fix from main
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
--
internal/engine/capabilities.go
internal/engine/capabilities_test.go
f92d42a mule/research: commit the calibration survey, and mark it triaged
.mule/research/issues-2026-07-27.md
.mule/research/survey-2026-07-27.md
.mule/research/triaged.md
c6e0301 mule/research: seed the branch the daily survey timer will commit to
.mule/research/README.md

$ git -C /root/astrate-mule remote -v; ls /root/astrate-mule/.mule/reports/log/ 2>/dev/null | tail -20
origin	git@github.com:astrate-platform/astrate.git (fetch)
origin	git@github.com:astrate-platform/astrate.git (push)
20260908T193127Z-recipe-milestones-proposed.md
20260908T193321Z-recipe-github-issues-proposed.md
20260908T193614Z-recipe-astarte-upstream-noop.md
20260908T193814Z-recipe-astarte-upstream-proposed.md
20260909T110758Z-recipe-code-review-proposed.md
20260909T111504Z-pairing-bearer-secret-test-done.md
20260909T165402Z-pairing-remoteip-fallback-test-done.md
20260909T165735Z-recipe-docs-sync-proposed.md
20260909T170637Z-docs-sync-hk-patch-endpoint-done.md
20260909T171614Z-docs-sync-hk-retention-field-done.md
20260909T174959Z-recipe-hygiene-timeout.md
20260909T175718Z-recipe-milestones-proposed.md
20260909T181348Z-recipe-github-issues-proposed.md
20260909T184315Z-recipe-astarte-upstream-proposed.md
20260909T191420Z-recipe-code-review-proposed.md
20260909T195057Z-realm-pure-helper-tests-done.md
20260909T200459Z-recipe-docs-sync-proposed.md
20260909T201715Z-docs-sync-native-compat-health-503-done.md
20260909T202852Z-docs-sync-native-compat-version-endpoints-done.md
20260909T230006Z-recipe-hygiene-timeout.md

The research doc referenced by the milestone isn't in this repo. Let me check the previous milestones run report (there was a recent one, commit `6e1ab2e`) and the current `for-giulio.md`, plus the remote.
→ Read .mule/for-giulio.md
→ Read .mule/reports/log/20260909T175718Z-recipe-milestones-proposed.md
Let me check the previous milestones run report and the current live state of the open issues and upstream releases.

$ gh release view -R astarte-platform/astarte --json tagName,publishedAt,isDraft,isPrerelease --jq '.' 2>&1 | head -20
{"isDraft":false,"isPrerelease":false,"publishedAt":"2026-08-07T14:33:58Z","tagName":"v1.3.3"}

$ gh issue view 92 --json number,title,state,updatedAt,labels,comments --jq '{n:.number,t:.title,s:.state,u:.updatedAt,l:[.labels[].name]}' ; gh issue view 93 --json number,title,state,updatedAt,labels --jq '{n:.number,t:.title,s:.state,u:.updatedAt,l:[.labels[].name]}'
{"l":["enhancement","upstream-parity","upstream-experimental"],"n":92,"s":"OPEN","t":"keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds","u":"2026-09-04T19:52:21Z"}
{"l":["mule-review","upstream-parity"],"n":93,"s":"OPEN","t":"In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged","u":"2026-09-04T19:54:58Z"}

→ Read docs/UPSTREAM-EXPERIMENTAL.md
Upstream is still at **v1.3.3 stable** (no v1.4.0 yet) — the gate hasn't moved since yesterday's run. Let me do the standing UPSTREAM-EXPERIMENTAL reconciliation and reconfirm the open issue set rather than trusting prior claims.

$ gh issue list --label upstream-parity --state open --limit 50 --json number,title,labels --jq '.[] | {n:.number, l:[.labels[].name], t:.title}'
{"l":["mule-review","upstream-parity"],"n":93,"t":"In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged"}
{"l":["enhancement","upstream-parity","upstream-experimental"],"n":92,"t":"keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds"}
{"l":["enhancement","milestone-4.0","upstream-parity"],"n":78,"t":"FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)"}

$ gh issue list --label milestone-3.0 --state all --limit 50 --json number,title,state --jq '.'
[]

$ rg -n "APICompatVersion" /root/astrate-mule/internal/realm/service.go -A2 | head -12
580:// APICompatVersion is the upstream Realm Management API level Astrate
581-// emulates, reported by GET /v1/{realm}/version. The Astarte Dashboard
582-// feature-gates UI sections on it (the Trigger Delivery Policies page
--
588:const APICompatVersion = "1.2.2"
589-
590-// GetDeviceRegistrationLimit returns the realm's device registration limit

$ rg -ni "milestone|parity-2026|v3.0" /root/astrate-mule/.mule/todo.md | head -20
(no output)

Nothing has changed since yesterday's run. Let me confirm `APICompatVersion` and the todo state to close out the verification.
$ gh release list -R astarte-platform/astarte --limit 5 --json tagName,isPrerelease,publishedAt --jq '.[] | "\(.tagName) prerelease=\(.isPrerelease) \(.publishedAt)"'
v1.4.0-rc.5 prerelease=true 2026-08-20T12:28:04Z
v1.3.3 prerelease=false 2026-08-07T14:33:58Z
v1.4.0-rc.4 prerelease=true 2026-08-06T13:40:20Z
v1.4.0-rc.3 prerelease=true 2026-07-31T10:48:39Z
v1.4.0-rc.2 prerelease=true 2026-07-14T10:16:22Z

Verified. The state is identical to the 2026-09-09 run — nothing to file, propose, or escalate.

**Milestone carried out (v3.0 — upstream Astarte parity 1.2.2 → 1.3/1.4):**

- **Reference:** `.mule/research/upstream-parity-2026-08-22.md` is absent from the repo and from `origin/mule/research` — already escalated to `.mule/for-giulio.md:51-55` (2026-09-05), still open. The milestone's investigation *output* (#47–#89) survives, so milestone work is unaffected.
- **Already filed:** `gh issue list --label milestone-3.0 --state all` → empty (this milestone's work lives under `upstream-parity`). Open `upstream-parity` issues are exactly the three known escalation items: **#92** keyAgreement decision (escalated 2026-09-05, `for-giulio.md:74-91`), **#93** aclhook comment rewrite (`mule-review`, pushed commit awaiting review), **#78** FDO (milestone-4.0, out of v3.0 scope).
- **UPSTREAM-EXPERIMENTAL standing item:** rows #67 (required/encrypted mapping fields) and #68 (`async_operation=false`) both reconcile only when upstream **v1.4 final** ships — upstream latest stable is still **v1.3.3**; v1.4.0 is rc.5. No reconciliation action yet.
- **APICompatVersion** stays **1.2.2** (`internal/realm/service.go:588`) per #90's frozen decision — bump not due (full 1.3 surface not complete, experimental rows unreconciled).
- **Gap buckets:** no machine-checkable gaps → nothing to file (the recipe cap "file at most five" is irrelevant — zero candidates). No new design decisions → nothing to escalate. The "complete, cut the tag" proposal (`for-giulio.md` step 5) is **not** warranted: the milestone's open decision set (#92 keyAgreement gated on stable v1.4.0) is real and is Giulio's call.

No changes made to any file; adding a fourth near-identical for-giulio entry one day after the 2026-09-09 run would be noise (MULE.md: "do not invent a change to have something to show").

```
Done: Milestone recipe run on v3.0 — no new gaps, no issues filed, no task lines; still gated on upstream v1.4 final and Giulio's #92/#93/#78 decisions
Files: none
Verified: gh issue list --label milestone-3.0 (empty); gh issue list --label upstream-parity --state open (#92,#93,#78 unchanged); gh release -R astarte-platform/astarte (stable v1.3.3, 1.4.0-rc.5); APICompatVersion=1.2.2 in internal/realm/service.go:588 -> pass
Unsure: nothing
Follow-ups: none
```
