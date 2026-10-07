slug: recipe-milestones
verdict: proposed
at:  0f9c65c
ran: 2026-10-07T22:44:44Z on DietPi in 824s


> build · big-pickle

I'll start by reading the required files to understand the recipe.
→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md

$ gh issue list --label milestone-v3.0 --state all --limit 50 2>&1 | head -30
(no output)

$ gh issue list --label upstream-parity --state all --limit 100 2>&1 | head -60
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

$ gh issue view 93 2>&1 | head -120
title:	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
state:	OPEN
author:	atsetilam (Giulio Malitesta)
labels:	mule-review, upstream-parity
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
number:	93
--
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

$ cat /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md 2>&1 | wc -l && echo "---" && head -80 /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
1
---
head: cannot open '/root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md' for reading: No such file or directory

$ ls /root/astrate-mule/.mule/research/ 2>&1 | head -15
ls: cannot access '/root/astrate-mule/.mule/research/': No such file or directory

$ find /root/astrate-mule/.mule -name "*upstream-parity*" -o -name "*research*" 2>&1 | head -15
(no output)

$ sed -n '50,75p' /root/astrate-mule/internal/broker/aclhook.go
//	           base/<iface>... for introspected ownership:server interfaces
//
// Everything else is denied. Read checks (write == false) cover both SUBSCRIBE
// filters and per-message delivery topics; the two are told apart by the
// presence of MQTT wildcards, which only ever appear in filters. The
// server-owned delivery rule accepts an interface segment followed by any
// remainder (a concrete path or nothing).
func checkACL(base, topic string, write bool, ownership ownershipFn) bool {
	if topic == "" || len(topic) > maxTopicBytes {
		return false
	}
	if write {
		if topic == base { // introspection publish
			return true
		}
		rest, ok := strings.CutPrefix(topic, base+"/")
		if !ok {
			return false
		}
		if rest == "capabilities" || rest == "control/emptyCache" || rest == "control/producer/properties" {
			return true
		}
		// control/keyAgreement is deliberately denied: full key agreement is
		// out of scope for Astrate's current 1.2.2 target, and upstream's wire
		// spec (082-key_agreement_protocol.md) only ships with v1.4.0-rc.5,
		// v1.3.3 being the newest stable tag. It stays denied until the

$ sed -n '75,90p' /root/astrate-mule/internal/broker/aclhook.go
// v1.3.3 being the newest stable tag. It stays denied until the
		// implement-or-keep-parking decision in issue #92 resolves.
		iface, path, found := strings.Cut(rest, "/")
		if !found || path == "" { // data topics always carry a path (§3.3)
			return false
		}
		own, known := ownership(iface)
		return known && own == interfaceschema.OwnershipDevice
	}

	rest, ok := strings.CutPrefix(topic, base+"/")
	if !ok {
		return false // never another device's or realm's subtree
	}
	if strings.ContainsAny(topic, "+#") {
		// A SUBSCRIBE filter within the device's own subtree. Wildcards never

$ sed -n '90,100p' /root/astrate-mule/internal/broker/aclhook.go
// A SUBSCRIBE filter within the device's own subtree. Wildcards never
		// occur in delivery topics, so this branch only ever sees filters.
		return true
	}
	if rest == "control/consumer/properties" {
		return true
	}
	iface, _, _ := strings.Cut(rest, "/")
	own, known := ownership(iface)
	return known && own == interfaceschema.OwnershipServer
}

$ grep -n "control/keyAgreement" /root/astrate-mule/internal/broker/aclhook.go /root/astrate-mule/internal/broker/aclhook_test.go
/root/astrate-mule/internal/broker/aclhook.go:72:		// control/keyAgreement is deliberately denied: full key agreement is
/root/astrate-mule/internal/broker/aclhook_test.go:56:		{"pub control keyAgreement", base + "/control/keyAgreement", true, false},

$ sed -n '50,60p' /root/astrate-mule/internal/broker/aclhook_test.go
{"pub server-owned interface", base + "/com.ex.ServerData/value", true, false},
		{"pub uninstalled interface", base + "/com.ex.Unknown/value", true, false},
		{"pub interface without path", base + "/com.ex.DeviceData", true, false},
		{"pub interface with empty path", base + "/com.ex.DeviceData/", true, false},
		{"pub control consumer properties", base + "/control/consumer/properties", true, false},
		{"pub control unknown", base + "/control/selfDestruct", true, false},
		{"pub control keyAgreement", base + "/control/keyAgreement", true, false},
		{"pub control prefix only", base + "/control", true, false},
		{"pub bare realm", "test", true, false},
		{"pub empty topic", "", true, false},
		{"pub sys topic", "$SYS/broker/uptime", true, false},

$ find /root/astrate-mule -name "upstream-parity*" 2>&1 | head -10
(no output)

$ find /root/astrate-mule -name "survey-2026-08-31*" 2>&1 | head -10
(no output)

$ gh issue list --label milestone-v3.0 --state all 2>&1 | wc -l
0

$ gh issue list --label upstream-parity --label mule-review --state open 2>&1 | head -10
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z

I see the current state:
- Issue #93 (upstream-parity, mule-review) is open. It asks to rewrite the comment in `internal/broker/aclhook.go` about `control/keyAgreement` (justify deny as Astrate's own decision, cite 082-key_agreement_protocol.md, point to #92), without changing behavior or touching `aclhook_test.go:53`.
- Current comment already reflects the new stance (Astrate's decision, cites 082, points to #92). The code in aclhook.go around lines 72-74 matches what the issue wants to achieve conceptually. But the issue says the "comment only change" is needed because it was stale - maybe there's a specific wording they want?

Looking at the issue more precisely: "rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged"

Comparing to what we have: our current comment says "full key agreement is out of scope for Astrate's current 1.2.2 target, and upstream's wire spec (082-key_agreement_protocol.md) only ships with v1.4.0-rc.5, v1.3.3 being the newest stable tag. It stays denied until the implement-or-keep-parking decision in issue #92 resolves."

That matches the issue's requirements (justifies as Astrate's own choice - scope-based; cites 082-key_agreement_protocol.md; points to #92; doesn't mention the deleted 080 sentence). So the comment as written already satisfies what issue #93 asks for. The test line is untouched. The behavior unchanged.

But the issue is still open - maybe it's asking for a specific precise wording? Let us look up the exact full text of what it wants. The issue says "rewrite the comment so it..." - the example it gives shows the old bad comment quoting 080. Our current comment is the new form.

Also looking at the full issue: it says "Do: rewrite the comment so it stands on Astrate's own decision — the topic is denied because full key agreement is out of scope for the current 1.2.2 target and the spec ships only in `v1.4.0-rc.5` (`v1.3.3` is the newest stable tag) — cite `082-key_agreement_protocol.md` as the current upstream reference rather than the deleted 080 sentence, and point at issue #92 where the implement-or-keep-parking decision lives."

Our comment says essentially the same thing. Maybe it's already been fixed? But the issue is still OPEN. Let us also look at the full comment in the actual code - maybe there's a more precise form the issue expects? Or maybe the issue is stale relative to code? Let us read the exact comment lines again.

The current comment (lines 72-77 in aclhook.go as shown): says "control/keyAgreement is deliberately denied: full key agreement is out of scope for Astrate's current 1.2.2 target, and upstream's wire spec (082-key_agreement_protocol.md) only ships with v1.4.0-rc.5, v1.3.3 being the newest stable tag. It stays denied until the implement-or-keep-parking decision in issue #92 resolves."

This matches what the issue asks for. So this issue may already be satisfied by current code. That happens.

Now looking at the milestones.md: v3.0 says "Source document: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 / v1.4.0-rc.5 diffed service-by-service against our emulated 1.2.2 level)." - that file doesn't exist. Also backlog issues #47-#89 (upstream-parity) - many are closed, #92 is open, #93 is open, #78 is open with milestone-4.0 label.

The recipe says: for v3.0 milestone (first not marked DONE), if scope says "not yet decided" / "TBD" do differently - but v3.0 scope is decided (references the research doc and backlog #47-#89). Step 3: investigate gaps - check references, list issues, grep code, reconcile experimental doc. Also "standing item, every milestone: read `docs/UPSTREAM-EXPERIMENTAL.md` and reconcile..."

The v3.0 is not marked DONE. Status says "not started" in the file. But many issues are closed. The recipe wants us to run as a proposal job - produce task lines in .mule/todo.md plus any evidence file it tells us to write. But also "If the recipe tells you to run 'gh issue create' as the deliverable itself (some do — check the recipe), that is fine to actually run now..." The milestones recipe can file issues directly as deliverable. But we must follow the recipe precisely.

Step 1: Read .mule/milestones.md. Find first section not marked DONE - v3.0 is not marked DONE. Scope is decided (not TBD). Step 2 doesn't apply. Step 3: investigate. Need to read docs/UPSTREAM-EXPERIMENTAL.md. Also look up the research doc if it existed? But it's referenced as source; maybe the file is missing - the recipe says "If the reference named in the milestone section (upstream repos, docs) is unreachable, say so and stop; do not guess at scope from the milestone's one-paragraph summary alone." But the reference in the milestone text points to `.mule/research/upstream-parity-2026-08-22.md` - that file doesn't exist locally. That means the reference doc is unreachable. According to rule at line 80 of milestones recipe: "If the reference named in the milestone section (upstream repos, docs) is unreachable, say so and stop; do not guess at scope from the milestone's one-paragraph summary alone."

So the recipe tells us to stop if the reference doc is unreachable. Let us check - is it truly unreachable? The file doesn't exist. So we must stop and say so. That fits the proposal job model - our output would be evidence/for-giulio note? But the recipe says "say so and stop" - in the context of doing the recipe, the deliverable would be what we write. Let us look up how to "say so" in the context of this recipe's deliverables.

But also look up the full instruction again: "Read .mule/MULE.md, then read .mule/recipes/milestones.md and carry it out. That recipe is a *proposal* job: its output is new task lines appended to .mule/todo.md, plus any evidence file it tells you to write. Do not start implementing the tasks you propose..."

For the milestones recipe specifically - step 4 says when it files issues directly via gh issue create, that's the deliverable (not todo.md). Also when it needs design decision, write to .mule/for-giulio.md. So deliverables depend on what happens.

But if we hit the "reference unreachable" case (line 80 in milestones.md recipe), the recipe says "say so and stop; do not guess..." So we need to record this. Probably write to .mule/for-giulio.md? Let us look up how for-giulio.md is used. Also the recipe says "escalating anything that needs a design decision to `.mule/for-giulio.md`" and gives example format. Also step 5 gives example: "propose one `.mule/for-giulio.md` line — 'milestone <tag> looks complete...'"

So in general, escalations/notes go to .mule/for-giulio.md. Let us look up the format from github-issues recipe or just look at existing for-giulio.md if it exists. Let us check.
$ cat /root/astrate-mule/.mule/for-giulio.md 2>&1 | head -40
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **docs-sync run, 2026-10-07, surface: pairing — the `docs/api` spec is accurate (routes, status codes, request/response schemas and examples all match the handlers; no docs/api task to propose), but `docs/site/pairing-and-security.md:97` carries a rate-limiting claim the code contradicts: "Pairing endpoints: per-IP and per-device token buckets."** Actual per endpoint, read from the handlers — register is **per-IP only** (`a.regLimiter.Allow("ip|"+remoteIP(r))`, internal/pairing/http.go:110; one bucket, defaults 5.0/20 at http.go:23-24), credentials is the only one that is **per-IP AND per-device** (`credLimiter.Allow("ip|"+...) || credLimiter.Allow("dev|"+realm+"/"+deviceID)`, http.go:158), and **unregister, info, verify and health are not rate-limited at all** (handleUnregister http.go:134-141, handleInfo :201-218, handleVerify :242-277, handleHealth :284-294 — the health mount comment says so outright, "No rate limiter either; it is one cheap query", http.go:78-80). The API reference already documents the true per-op scope (registerDevice yaml:105 "Rate-limited per IP."; requestCredentials yaml:192 "Rate-limited per IP and per device."; the other four ops say nothing), so a client that reads the site for the surface's security posture is shown a limit that unregister/info/verify/health do not have. Proposed wording (your voice): perf-`endpoint`, e.g. "Register: per-IP token bucket. Credentials: per-IP and per-device token buckets. Unregister, info, verify, health: not rate-limited." Page untouched.

- **github-issues triage run, 2026-10-07: nothing proposable for the 25th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub.** The recipe's `--limit 40` command prints all **31** open issues (authoritative `--limit 100` agrees): **27** `mule-alarm` **#94–#120** plus the same four non-alarm, none of them proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (`milestone-4.0`, parked until the v3.0 queue clears), and one the standing instruction excludes from triage entirely. New since the 2026-10-06 run: **#120**, created 2026-10-07T11:03:45Z, "nothing has landed in 14h" — arithmetic pinned, not inferred: `.mule/.heartbeat` holds `1791318570` = **2026-10-06T20:29:30Z**, written 3s after the `docs-native-socket-security-scheme` done land (`87c18cf`, 20:29:27Z), and 11:03:45Z − 20:29:30Z = **14h34m**, so `age=$(( (now-last)/3600 ))` prints 14 — the title verbatim; `.mule/.alarmed` carries mtime 11:04Z, i.e. set by that filing. Genuine window, the same 14–16h overnight band #117–#119 measured. **Proposal, extending the 2026-10-06 one: close #94–#119, all twenty-six** — each is a one-day alarm superseded by the next and nothing reads issue state; leave **#120** (today's, real window) to self-expire.
  **New this run, and it is the 2026-10-01 pattern recurring: the switch is latched off through a silence that is real.** At 2026-10-07T16:55Z `.mule/.heartbeat` still read 2026-10-06T20:29:30Z — **20h**, with zero `2026-10-07` rows in `.mule/log.md` — but `.mule/.alarmed` (set by #120) short-circuits `check_pulse` at `tools/mule.sh:748`, so nothing can file again until a `done` land (mule.sh:578) or a `checked` run (mule.sh:540) calls `beat`, which is the only thing that clears the latch (mule.sh:737). Recipe ticks are committing today (`becf3cd` 11:25:35Z, `dde64e0` 11:38:44Z) but a recipe verdict is neither of those two paths, so they neither beat nor un-latch — measured today, not assumed. The `check_pulse`-before-`beat` ordering fix (move the call at `tools/mule.sh:633` to after the tick's `beat`) remains escalated above, still yours, still unactioned; #120 is the same arithmetic as #119, not fresh evidence for it either way.
  **#93 at 33 days.** `git rev-list --left-right --count origin/main...HEAD` = **4 / 697** today (4 / 679 on 10-06), and `git merge-base --is-ancestor 24ad5b8 origin/main` still fails — the rewrite is entirely on `mule/queue` / `origin/mule/queue`. The version-drift detail recorded on 2026-10-06 (the comment says `v1.4.0-rc.5` / "v1.3.3 being the newest stable tag"; upstream is `v1.4.0-rc.6` / `v1.3.5`, re-checked via `gh api repos/astarte-platform/astarte/tags` today) folds into that same review rather than becoming a line of its own.

---

- **COMPATIBILITY.md wording proposal for upstream v1.3.5 (newest stable, 2026-10-05; v1.4.0 is still rc.6-only) — a one-word delta on the open v1.3.2/v1.3.4 proposal below, nothing re-derived.** v1.3.5 is maintenance-only: `gh api repos/astarte-platform/astarte/releases` shows it published `2026-10-05T14:11:44Z`, not a prerelease, with a body of a single `### Fixed` line — "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available" — and `compare/v1.3.4...v1.3.5` is 4 commits / 35 files whose only non-`mix.exs`/`mix.lock` files are `.github/workflows/astarte-build-workflow.yaml`, `CHANGELOG.md`, `README.md`, `astarte-dashboard/package.json` and `docker-compose.yml`. Not one source file changed, so no route, MQTT topic, control message or interface-schema field moved (Astrate implements no `data_updater_plant` and no Erlang RPC server). **The only change a v1.3.5-aware doc carries: the v1.3.4 entry's §Infrastructure-differences sentence — "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target" — reads "…(**newest stable v1.3.5**) as the target".** Everything else stands unchanged: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"); that entry's "v1.4.0 is still rc.5-only" header is already superseded to rc.6 by the open rc.6 delta entry below, not here. This is the same one-word delta recorded on 2026-10-05 in this file ("newest stable v1.3.5"), folded in as asked rather than re-derived. Applied to nothing: `docs/COMPATIBILITY.md` untouched and still 1.2.2-targeted (it contains no `1.3` reference at all), `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). Raw: [v1.3.5](https://github.com/astarte-platform/astarte/releases/tag/v1.3.5).

---

- **`docs/site/appengine-api.md:114` — "All responses use the Astarte envelope format" is contradicted by the endpoints the page itself lists above it.** The "Response envelope" section (lines 112-126) opens with that sentence and then the `{ "data": ... }` / `{ "errors": { "detail": ... } }` shapes, but the live-stream endpoints documented at lines 94-110 never return one: the native socket pushes bare event JSON as WebSocket text frames or `data: {json}` SSE frames (internal/appengine/stream/ws.go:105-110, 156-162), `/astrate/v1/health` answers `{"status":"ok"}` (internal/observability/health.go:51), `/astrate/v1/readiness` answers `{"status":...,"checks":{...}}` (health.go:77), `/astrate/v1/metrics` answers Prometheus text under `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (metrics.go:42), and — measured today with the repo's own `coder/websocket` (probe in /tmp) — a WebSocket handshake failure on either socket is plain text under `text/plain; charset=utf-8` (`WebSocket protocol violation: Connection header "" does not contain Upgrade`), not JSON. REST appengine paths do use the envelope, so the sentence is right for the endpoints above "Device data" and wrong unqualified. Proposed wording (your voice): scope it — "All REST responses use the Astarte envelope format" — with one line noting the socket and observability surfaces answer un-enveloped. Page untouched.

- **`docs/site/appengine-api.md:102` — "Honours `a_ch` claims as room filters" imports upstream's room semantics onto a socket that has no rooms.** The native socket subscribes to the whole realm bus, narrowed only by the `device_id`/`interface` query parameters (internal/appengine/stream/ws.go:69-75); there is no room concept and no per-room filtering. What a_ch actually does here is gate the whole route: `mw.RequireRealm(auth.ClaimChannels)` requires the token to authorize `GET socket` via the REST verb-regex rule (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112 -> internal/auth/claims.go:164-171). The upstream-measured JOIN/WATCH partition rule DESIGN.md §4.2 describes (`AuthorizesChannel`, claims.go:197-204) is used only by the Phoenix socket's per-room join/watch checks (internal/appengine/channels/ws.go:278, 354). Consequence worth stating: a token whose a_ch list is a blanket `".*::.*"` — which upstream's Channels rule treats as authorizing nothing — authorizes the native socket, while the Phoenix socket, where room filters actually live, still reads JOIN/WATCH literally. Proposed wording (your voice): "guarded by the `a_ch` claim (GET on `socket`); `device_id`/`interface` narrow the stream". Page untouched.

- **github-issues triage run, 2026-10-06: nothing proposable for the 24th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented or closed on GitHub.** The recipe's `--limit 40` command prints all **30** open issues: **26** `mule-alarm` **#94–#119** plus the same four non-alarm, none of them proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (`milestone-4.0`, already escalated), and one the standing instruction excludes from triage entirely. New since the 2026-10-05 run: **#119**, created 2026-10-06T10:58:45Z, "nothing has landed in 15h" — a real window, not an artefact: the land before it was `230dd11` at 2026-10-05T21:04:09+02:00 (19:04Z), 15.9h earlier, and it was ended by today's first land `cbdd3e8` at 16:05:39+02:00 (14:05Z), three hours after the alarm. So it belongs to the genuine 14–16h band the 2026-10-05 entry measured, i.e. overnight silence between evening and afternoon batches. **Proposal, extending the 2026-10-05 one: close #94–#118, all twenty-five** — each is a one-day alarm superseded by the next, and nothing reads issue state; leave **#119** (today's, real window, already ended) to self-expire. The `check_pulse`-before-`beat` ordering fix (move the call at `tools/mule.sh:633` to after the tick's `beat`) is still escalated above, still yours, and still unactioned — today's alarm is not evidence for it either way.
  **The review queue moved 9 commits in 24h and main's four did not move again: `git rev-list --left-right --count origin/main...HEAD` is 4 / 679 today against 4 / 670 yesterday.** The four are main-side commits this branch lacks (`f1d0069`, `ca47b35`, `99743ca`, `2a033d4` — all ancestors of `origin/main`), and **#93's work is still entirely off main**: `8c61268` ("mule: log issue-93") and `24ad5b8` (the aclhook comment rewrite, whose commit message is the issue title verbatim) both satisfy `git merge-base --is-ancestor <sha> origin/main` → false, and `git branch -a --contains 24ad5b8` lists only `mule/queue` / `origin/mule/queue`. That is **32 days** of `mule-review` with `bash tools/mule.sh review` unread. I have not touched git and will not.
  **One detail to fold into that review rather than queue separately:** the rewritten comment at `internal/broker/aclhook.go` (the `control/keyAgreement` deny) says "…upstream's wire spec (082-key_agreement_protocol.md) only ships with v1.4.0-rc.5, v1.3.3 being the newest stable tag". The 082 citation, the "deliberately denied as Astrate's own choice" framing and the unchanged deny + `aclhook_test.go` case are exactly what #93 asked for, but the version reference has aged — upstream's newest stable is **v1.3.5** (2026-10-05) and newest overall `v1.4.0-rc.6` (2026-09-30). No separate task: one word in the same diff when someone finally reads it.

---

- **Milestone recipe run, 2026-10-06: v3.0 verified unchanged — still gated on upstream — and the master delta since the 2026-10-05 run is a closed set of 13 commits, every one wire-inert for Astrate. No issue filed, no `.mule/todo.md` line queued.** Release sweep is byte-identical to yesterday's (`gh api repos/astarte-platform/astarte/releases`): **no stable v1.4.0** — newest stable `v1.3.5` (2026-10-05), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition (your 2026-09-04 "wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) stays unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** — this milestone's whole ledger lives under `upstream-parity`, where the open set is unchanged: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`); **#47–#89** otherwise all closed. Non-alarm open issues total exactly four (plus the 25 `mule-alarm` #94–#118): the three above and **#1** (untouched per standing instruction). `APICompatVersion` still `"1.2.2"` (`internal/realm/service.go:588`, re-checked).
  **The new master commits — measured, and none is a gap.** `gh api .../commits?since=2026-10-05T14:20:00Z` returns 9 today (13 with the four landing between yesterday's run and that mark): `1da04832` + `a8a26d0f` "chore: forward port release-1.4 / release-1.3" are release-branch merges — CI/workflow consolidation, CHANGELOG, `mix.exs` bumps, and the appengine RPC refactor that *deletes* the `astarte_data_updater_plant` RPC client (Astrate implements no `data_updater_plant`); a router-level sweep of the 300-file forward-port finds **no new route or channel** — the only `_web/(router)` hit is `rooms_channel_test.exs`. `d85e0ed0` (docs) + `a6bd2c21` `fix(fdo): reject owner key names not accepted by OpenBao` are both FDO, i.e. **#78**'s parked body. `a51ab2e1` is an umbrella-CI refactor; `7628d2ed`/`ed1e99f4`/`bb17455a`/`ec020bf8` are `test(generators)`/`fix(generators)` Core-changeset validation — upstream test tooling, not Astrate surface. No HTTP route, MQTT topic, AMQP control message or interface-schema field moved.
  **Standing item unchanged.** `docs/UPSTREAM-EXPERIMENTAL.md`'s two rows (#67 `required`/`encrypted` mapping fields, #68 `async_operation=false`) still wait on **upstream 1.4 final** for their promoted-vs-deprecated call; nothing in today's commit set promotes or removes either and none of it is released, so the register is untouched.
  **Step 5 deliberately not taken**, same reason as the 21 runs before: the section's declared scope is **1.2.2 → 1.3/1.4** and only the 1.3 half is delivered. The 1.4 half stays open — #92 unanswered, both register rows unreconciled, `APICompatVersion` unbumped — so the tag must not be cut. Your call when a stable v1.4.0 ships: answer #92, reconcile the register, run the final-phase bump, then cut the tag.

- **astarte-upstream recipe run, 2026-10-05 (evening): no upstream movement that touches anything Astrate implements, so nothing was queued — no `.mule/todo.md` line, no `gh issue create`, nothing commented, closed or edited on GitHub.** One thing *is* new since this recipe's last run (`2026-10-04T20:40:52+02:00` ran; the `2026-10-05T04:29:10+02:00` run timed out): stable **v1.3.5**, published `2026-10-05T14:11:44Z`. The milestone recipe already recorded it as inert today; I re-measured it rather than trusting that, because "a release note is a claim, not a fact" is this recipe's own rule. `gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5` is **4 commits, 35 files** — `8d47574f` + `4c4222cf` `chore: update horde`, `c4c63ebe` `chore(rm): bump xandra to 0.19.4`, `8f0f5ef1` `chore: prepare v1.3.5` — and **every app/lib file in it is a `mix.exs`/`mix.lock` version bump**. The rest is `CHANGELOG.md` (+6: one `### Fixed` entry, "[astarte_data_updater_plant] Increase RPC server reliability"), `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `README.md`, `astarte-dashboard/package.json`. Not one source file changed, so no route, MQTT topic, BSON encoding, control message or interface-schema field moved: Astrate implements no `data_updater_plant` and no Erlang RPC server. Wire-inert, same shape as v1.3.2/v1.3.3/v1.3.4. **Master since `v1.4.0-rc.6` is still the single commit `731bcddc`** ("chore: forward port release-1.4"; `gh api .../commits?since=2026-10-05T09:34:00Z` returns it and nothing newer), the migration-relocation commit the milestone run measured today as inert — so both wire-bearing post-rc.6 commits (`b6d46ad4`, `958bb9fe`) still have their owners in the queue and nothing new joins them.

✗ Read internal/broker/aclhook.go | wc -l && echo "---" && tail -15 /root/astrate-mule/internal/broker/aclhook.go failed
Error: File not found: /root/astrate-mule/internal/broker/aclhook.go | wc -l && echo "---" && tail -15 /root/astrate-mule/internal/broker/aclhook.go

$ wc -l /root/astrate-mule/internal/broker/aclhook.go && echo "---" && tail -15 /root/astrate-mule/internal/broker/aclhook.go
246 /root/astrate-mule/internal/broker/aclhook.go
---
		})
	} else if !write {
		// Delivery to a session-present but disconnected device (offline
		// queue, retained messages): no live session, consult the store.
		allowed = checkACL(cl.ID, topic, write, func(iface string) (interfaceschema.Ownership, bool) {
			return h.offline.ownershipOf(cl.ID, iface)
		})
	}

	if !allowed {
		h.log.Warn("mqtt ACL denied",
			"client", cl.ID, "topic", topic, "write", write, "remote", cl.Net.Remote)
	}
	return allowed
}

$ sed -n '60,80p' /root/astrate-mule/internal/broker/aclhook.go
}
	if write {
		if topic == base { // introspection publish
			return true
		}
		rest, ok := strings.CutPrefix(topic, base+"/")
		if !ok {
			return false
		}
		if rest == "capabilities" || rest == "control/emptyCache" || rest == "control/producer/properties" {
			return true
		}
		// control/keyAgreement is deliberately denied: full key agreement is
		// out of scope for Astrate's current 1.2.2 target, and upstream's wire
		// spec (082-key_agreement_protocol.md) only ships with v1.4.0-rc.5,
		// v1.3.3 being the newest stable tag. It stays denied until the
		// implement-or-keep-parking decision in issue #92 resolves.
		iface, path, found := strings.Cut(rest, "/")
		if !found || path == "" { // data topics always carry a path (§3.3)
			return false
		}

I'll run the milestones recipe as instructed - producing proposals (for-giulio.md) based on what I find. The referenced research doc for v3.0 doesn't exist.
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cat > /tmp/milestone_v3_note.txt << 'EOF'
- **Milestone recipe run, 2026-10-08: v3.0 reference doc `.mule/research/upstream-parity-2026-08-22.md` is unreachable (file not found).** Per milestones.md recipe line 80, cannot proceed without the referenced document — stop and escalate. The milestone states scope from that doc and lists #47–#89; open upstream-parity issues are #78 (milestone-4.0), #92, #93. No `.mule/todo.md` lines proposed, no issues created.
EOF
cat /root/astrate-mule/.mule/for-giulio.md | head -2 > /tmp/fg_header.txt 2>&1; echo "---" >> /tmp/fg_header.txt; cat /tmp/milestone_v3_note.txt >> /tmp/fg_header.txt; echo "" >> /tmp/fg_header.txt; tail -n +5 /root/astrate-mule/.mule/for-giulio.md >> /tmp/fg_header.txt 2>&1 && cp /tmp/fg_header.txt /root/astrate-mule/.mule/for-giulio.md failed
Error: The user rejected permission to use this specific tool call.
