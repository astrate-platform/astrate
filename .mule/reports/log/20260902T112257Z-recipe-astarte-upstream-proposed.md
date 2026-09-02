slug: recipe-astarte-upstream
verdict: proposed
at:  c465e64
ran: 2026-09-02T11:22:57Z on DietPi in 104s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
I'll read the compatibility file first, then check upstream releases.
→ Read docs/COMPATIBILITY.md
The compatibility doc targets upstream v1.2.2 (realm API compat version 1.2.2, dashboard v1.2.2, measured against v1.2.0). But let me note that COMPATIBILITY.md doesn't explicitly state a release tag target in the version line. Let me check the releases.

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

The newest stable tag is `v1.3.3`, plus RC candidates for v1.4.0. The COMPATIBILITY.md names `1.2.2` (realm API compat version). Astrate targets upstream v1.2.x. Upstream is ahead — v1.3.3 is the newest stable release.

Let me check git history to see if this has been reviewed before (to understand what the last tracked tag was), and look at the release notes.
$ rg -n -i "upstream|v1\.2|v1\.3|compat" /root/astrate-mule/.mule/todo.md 2>/dev/null | tail -40
65:`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
97:- [!] control-producer-properties-compression: in `internal/engine/control.go`, accept a plaintext (and the 4-zero-byte empty-frame) device→server `producer/properties` purge list when the device declared `purge_properties_compression_format: plaintext` — `inflateProperties` currently rejects every non-zlib payload while upstream's `control_handler.ex` decodes per-capability **Approved by Giulio 2026-08-31**: yes, devices that cannot compress must be able to talk to us, following upstream. — BLOCKED: gates failed
98:- [!] probe-interface-default-values: does `GET /realmmanagement/v1/<realm>/interfaces/<name>/<major>` return the same defaulted mapping-parameter values upstream reports after 3f0b864 ("Properly show default values for all mapping parameters")? report, do not patch [auto] — BLOCKED: gates failed
99:- [!] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto] — BLOCKED: gates failed
100:- [!] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
101:- [!] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto] — BLOCKED: gates failed
102:- [!] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
103:- [!] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto] — BLOCKED: gates failed
104:- [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
105:- [!] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto] — BLOCKED: gates failed
108:- [!] issue-68-async-operation-accepted: on the mutating endpoints upstream surfaces `async_operation` on (housekeeping realm create/delete; realm-management interface install/update/delete, trigger/policy delete in internal/realm/http.go), accept and ignore `?async_operation=false`, with unparseable/`true` values also not changing behaviour (Astrate stays always-sync); tests cover the flag parsing (closes #68) — BLOCKED: gates failed
109:- [!] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto] — BLOCKED: gates failed
110:- [!] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
111:- [!] realm-policy-list-sorted: in `internal/realm/service.go`, sort the names returned by `ListPolicies` the way `ListTriggers`/`ListInterfaces`/`ListInterfaceMajors` already do (it is the only list method without `sort.Strings`), and extend `TestDashboardCompat.Policies` (or a new subtest in `http_test.go`) to install several policies in non-alphabetical order and assert `GET /policies` comes back sorted. [auto] — BLOCKED: gates failed

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'
unexpected end of JSON input

$ cd /root/astrate-mule && git log --oneline -15 -- .mule/todo.md 2>/dev/null
689fb1d mule: blocked realm-interface-lookup-404
7f29b3c mule: blocked realm-policy-list-sorted
4dbba24 mule: refill the queue
d67d8ea mule: blocked compat-note-v140-rc5
54b406b mule: blocked probe-object-validation-selected-interface
4cd1e88 mule: refill the queue
f4ad882 mule: blocked issue-68-async-operation-accepted
e0f8823 mule: blocked issue-91-pairing-health-serve-root
733f965 mule: refill the queue
5006f6e mule: blocked docs-sync-pairing-status-enum
49fcd62 mule: refill the queue
19ad0fe mule: blocked probe-properties-on-connect-encoding
8a983a3 mule: blocked probe-binaryblob-validation
25560b4 mule: blocked probe-mqtt-capabilities-declaration
21f74d8 mule: blocked compat-note-v133

The todo already has many entries related to upstream — including `probe-emptycache-resend-device-error` (v1.3.3), `compat-note-v1.3.3`, `compat-note-v1.4.0-rc.3`, and `compat-note-v1.4.0-rc.5`. These appear to have been proposed in previous runs. Since v1.3.3 has already been noted and v1.4.0-rc.5 too, but those tasks are marked `[!]` (blocked). 

The newest release is v1.4.0-rc.5 (2026-08-20), but v1.4.0-rc.5 has already been probed. Wait — let me check what the latest state is. The newest stable is v1.3.3. But there's already a `compat-note-v1.3.3` line and probes for v1.3.3 features. So this run's job is to see if there's anything new.

The releases listed are v1.4.0-rc.5, v1.3.3, v1.4.0-rc.4, v1.4.0-rc.3, v1.4.0-rc.2 (top 5). The newest is v1.4.0-rc.5 but the newest **stable** is v1.3.3. Let me check whether v1.4.0-rc.5 content has already been fully reviewed, or if v1.4.0-rc.5 is newer than what's been examined.

The existing todo has `probe-binaryblob-validation ... v1.4.0-rc.5's corrected validator` and `probe-object-validation-selected-interface ... v1.4.0-rc.5 #2141` and `compat-note-v1.4.0-rc.5`. So v1.4.0-rc.5 has already been reviewed and proposed.

Since this run's job is to watch upstream — and v1.3.3 and v1.4.0-rc.5 have both already been covered in previous runs — the question is whether there's any newer movement. The newest tag here is v1.4.0-rc.5 which is already covered. There's nothing newer than what's already been examined.

Wait, but let me reconsider. The recipe says "Compare the newest tag to the version named in docs/COMPATIBILITY.md." The compat doc targets v1.2.x (realm API compat 1.2.2). Upstream is ahead (v1.3.3 stable, v1.4.0-rc.5 newest). But the todo already has proposals covering v1.3.3 and v1.4.0-rc.5. 

So the situation: the prior runs already proposed tasks for these versions. Today's run should check if there's anything *newer* than what's already been reviewed. The newest tag is v1.4.0-rc.5 (2026-08-20), which is already covered.

Let me verify whether these existing lines came from previous runs by checking git log for todo.md, and confirm there's nothing genuinely new. Actually, let me check whether the newest release (v1.4.0-rc.5) has already been fully processed. Since it's already in todo.md as proposed lines, the current run has nothing new to add.

But wait — I should double check. The recipe says the expected good outcome is "If they match, say so and stop." They don't match (upstream ahead). But upstream movement has already been captured. The rules say "If nothing changed since last time, append nothing to the queue."

Let me verify there's no even-newer activity. Let me check if v1.4.0-rc.5 bodies were already reviewed, and confirm the existing todo proposals for these versions. Let me check the git log to confirm these lines were added recently.
$ cd /root/astrate-mule && rg -n "compat-note-v1\.4\.0-rc\.5|probe-binaryblob|probe-object-validation-selected" .mule/todo.md
104:- [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
109:- [!] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto] — BLOCKED: gates failed
110:- [!] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed

The upstream v1.3.3 and v1.4.0-rc.5 changes have already been reviewed and proposed in prior runs (they're blocked in the todo queue, but the *proposals* were already made). The newest release tag in the list is v1.4.0-rc.5 (2026-08-20), which is already fully covered by existing task lines (`compat-note-v1.4.0-rc.5`, `probe-binaryblob-validation`, `probe-object-validation-selected-interface`).

Let me confirm there is no release newer than v1.4.0-rc.5 that isn't yet reviewed. The top-5 list I got is: v1.4.0-rc.5, v1.3.3, v1.4.0-rc.4, v1.4.0-rc.3, v1.4.0-rc.2. Nothing newer than v1.4.0-rc.5 exists. And v1.4.0-rc.5 is already reviewed.

So there's **no new upstream movement** since the last run. According to the recipe: "If nothing changed since last time, append nothing to the queue and say 'no upstream movement since <tag>'."

The last fully-reviewed/newest tag is v1.4.0-rc.5. I should append nothing to the queue.

Let me verify I haven't missed any release between the top-5 and what's already been processed. Let me check that v1.4.0-rc.5 is indeed the newest and that nothing came after it.

The dates confirm: v1.4.0-rc.5 = 2026-08-20 is the newest. Today is 2026-09-02. Nothing newer exists.

No upstream movement since v1.4.0-rc.5. I append nothing.

Done: No upstream movement since v1.4.0-rc.5 (2026-08-20); nothing appended.
Files: none
Verified: `gh api repos/astarte-platform/astarte/releases` -> the newest tag is v1.4.0-rc.5, already covered by existing todo lines (compat-note-v1.4.0-rc.5, probe-binaryblob-validation, probe-object-validation-selected-interface)
Unsure: nothing
Follow-ups: none

Wait — let me reconsider. The newest **stable** release is v1.3.3. The newest **tag** (including RC) is v1.4.0-rc.5. The recipe says "Compare the newest tag to the version named in docs/COMPATIBILITY.md." Both v1.3.3 and v1.4.0-rc.5 have already been proposed in the existing queue. So there's genuinely nothing new to add.

Let me double check there's no task in the queue already for the *latest* aspect I might be missing. Actually the relevant point is: has this exact run already happened, i.e. is v1.4.0-rc.5 the latest reviewed? Given the todo already carries `compat-note-v1.4.0-rc.5`, that version was already processed. Today's run sees no newer tag than v1.4.0-rc.5. So nothing new.

The recipe's expected good outcome is exactly this: "no upstream movement since v1.4.0-rc.5". I append nothing. Report.
