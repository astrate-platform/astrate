slug: recipe-code-review
verdict: proposed
at:  d67d8ea
ran: 2026-09-01T18:47:04Z on DietPi in 168s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
→ Read .mule/reviews
→ Read internal
→ Read pkg
→ Read .mule/log.md

$ rg -c '' --glob 'pkg/**/*.go' | sort -t: -k2 -rn | head -20
pkg/payload/payload_test.go:787
pkg/interfaceschema/parse.go:501
pkg/interfaceschema/types.go:406
pkg/payload/json.go:396
pkg/interfaceschema/parse_test.go:352
pkg/payload/value.go:349
pkg/payload/bson.go:343
pkg/interfaceschema/compat_test.go:250
pkg/interfaceschema/compile_test.go:235
pkg/deviceid/deviceid_test.go:226
pkg/astarteapi/envelope_test.go:222
pkg/payload/value_test.go:221
pkg/astarteapi/envelope.go:219
pkg/interfaceschema/trie_test.go:216
pkg/payload/payload.go:199
pkg/interfaceschema/trie.go:155
pkg/deviceid/deviceid.go:132
pkg/payload/fuzz_test.go:122
pkg/payload/sniff_test.go:120
pkg/interfaceschema/compile.go:120

$ rg -c '' --glob 'internal/**/*.go' | sort -t: -k2 -rn | head -20
internal/appengine/channels/ws_test.go:916
internal/engine/triggers/actions_test.go:805
internal/store/datastreams_test.go:783
internal/engine/fakestore_test.go:769
internal/pairing/service_test.go:632
internal/broker/broker_e2e_test.go:618
internal/pairing/http_test.go:614
internal/engine/triggers/match.go:595
internal/broker/sessionstore.go:533
internal/appengine/http_test.go:517
internal/engine/triggers/actions.go:514
internal/engine/cache.go:496
internal/realm/service.go:491
internal/engine/router.go:490
internal/broker/broker_test.go:475
internal/store/devices.go:462
internal/engine/data_test.go:461
internal/store/devices_test.go:446
internal/engine/batch_test.go:446
internal/broker/authhook.go:440

→ Read .mule/reviews/store-2026-07-27.md

$ rg -c '' --glob 'internal/engine/**/*.go' | awk -F: '{sum+=$2} END{print "engine total:", sum}' && rg -c '' --glob 'internal/broker/**/*.go' | awk -F: '{sum+=$2} END{print "broker total:", sum}' && rg -c '' --glob 'internal/pairing/**/*.go' | awk -F: '{sum+=$2} END{print "pairing total:", sum}' && rg -c '' --glob 'pkg/interfaceschema/**/*.go' | awk -F: '{sum+=$2} END{print "interfaceschema total:", sum}' && rg -c '' --glob 'pkg/payload/**/*.go' | awk -F: '{sum+=$2} END{print "payload total:", sum}'
engine total: 12413
broker total: 3441
pairing total: 2777
interfaceschema total: 2345
payload total: 2578

$ ls internal/engine/*.go | sort
internal/engine/batch.go
internal/engine/batch_test.go
internal/engine/bus_events_test.go
internal/engine/cache.go
internal/engine/cache_test.go
internal/engine/capabilities.go
internal/engine/capabilities_test.go
internal/engine/control.go
internal/engine/control_test.go
internal/engine/data.go
internal/engine/data_test.go
internal/engine/engine.go
internal/engine/engine_e2e_test.go
internal/engine/engine_integration_test.go
internal/engine/fakestore_test.go
internal/engine/introspection.go
internal/engine/introspection_test.go
internal/engine/router.go
internal/engine/router_test.go
internal/engine/serverdata.go
internal/engine/serverdata_test.go
internal/engine/topics.go
internal/engine/topics_test.go

$ ls internal/engine/
batch.go
batch_test.go
bus_events_test.go
cache.go
cache_test.go
capabilities.go
capabilities_test.go
control.go
control_test.go
data.go
data_test.go
engine.go
engine_e2e_test.go
engine_integration_test.go
fakestore_test.go
forward
introspection.go
introspection_test.go
router.go
router_test.go
serverdata.go
serverdata_test.go
stream
testdata
topics.go
topics_test.go
triggers

$ rg -c '' --glob 'internal/engine/triggers/**/*.go' | awk -F: '{sum+=$2} END{print "triggers total:", sum}' && rg -c '' --glob 'internal/housekeeping/**/*.go' | awk -F: '{sum+=$2} END{print "housekeeping total:", sum}' && rg -c '' --glob 'internal/realm/**/*.go' | awk -F: '{sum+=$2} END{print "realm total:", sum}'
triggers total: 3599
housekeeping total: 500
realm total: 1386

$ ls internal/realm/
dashboard_compat_test.go
http.go
http_test.go
policy.go
policy_test.go
service.go

`internal/realm` is a good candidate: not in the recent review (only `store` was done), a manageable size (~1386 lines + tests), not yet reviewed. Let me go deep there.
$ rg -c '' internal/realm/*.go | sort -t: -k2 -rn
internal/realm/service.go:491
internal/realm/http.go:318
internal/realm/http_test.go:301
internal/realm/dashboard_compat_test.go:189
internal/realm/policy_test.go:70
internal/realm/policy.go:17

→ Read internal/realm/service.go
Now let me read the policy validation and the http.go, plus the tests.
→ Read internal/realm/policy.go
→ Read internal/realm/http.go
Now let me read the tests to understand what's covered.
→ Read internal/realm/http_test.go
→ Read internal/realm/policy_test.go
Now let me read the dashboard_compat_test.go to understand what's covered there.
→ Read internal/realm/dashboard_compat_test.go
Now I have a full picture of the `internal/realm` package. Let me analyze carefully for genuine issues to propose.

Let me look at a few things more carefully:

**1. `DeletePolicy` scans triggers by name match, but ignores unmarshal errors (`continue`).** In `service.go:392-399`, `DeletePolicy` unmarshals each trigger's definition into a `triggerPolicyRef`; if the unmarshal fails it `continue`s (treats it as not referencing). A corrupt trigger definition (which shouldn't happen since `CreateTrigger` always validates through `triggers.Compile`) would let a policy deletion proceed even if the trigger referenced it. This is defensive but the `continue` silently swallows. However, since `CreateTrigger` enforces parseability, this path is unreachable in practice. Low value.

**2. `DeleteInterface` doesn't check in its own invalidation** — it calls `interfacesChanged` correctly.

**3. Missing test: `UpdateInterface` with a stored interface that fails to parse.** In `service.go:140-143`, there's an error branch `"realm: stored interface %s v%d does not parse"` that can never be reached because `InstallInterface`/`UpdateInterface` always validate. Defensive.

**4. The `DeleteInterface` route (`DELETE /interfaces/{name}/{major}`)** — let me check: the `deleteDevice` path and others use major. Actually `DeleteInterface` is at http.go mount line 40. Fine.

**5. Duplicate validation in `CreateTrigger` and `Compile`.** `triggers.Compile` returns `ct.PolicyName`. Then it re-validates policy exists. OK.

**6. Missing test — device registration limit behavior.** The `config/device_registration_limit` GET returns nil for unlimited. The dashboard test checks `{"data":null}`. But there's no test setting/modifying the limit since there's no PUT for it. Not much to test.

**7. `ListPolicies` vs `ListTriggers` ordering.** Both sort. `ListPolicies` doesn't sort (line 374-377 in service.go) — it just copies in store order! Compare to `ListTriggers` (line 319-324) which does `sort.Strings(names)`. And `ListInterfaces` sorts, `ListInterfaceMajors` sorts. But **`ListPolicies` does NOT sort** its names. This is an inconsistency. Let me verify: 

service.go ListTriggers:
```go
names := make([]string, 0, len(trs))
for i := range trs {
    names = append(names, trs[i].Name)
}
sort.Strings(names)
return names, nil
```

service.go ListPolicies:
```go
names := make([]string, len(ps))
for i := range ps {
    names[i] = ps[i].Name
}
return names, nil
```

Indeed — **`ListPolicies` is the only list method that does not sort its output.** This is a real inconsistency. All sibling methods (`ListInterfaces`, `ListInterfaceMajors`, `ListTriggers`) sort their output for stable ordering. Whether the store returns policies in insertion order or name order, the API response order is undefined/unsorted while every sibling is sorted. This is a small clarity/consistency bug worth proposing — a mechanical change with an existing test (`TestDashboardCompat.Policies` checks only a single-element list), so no test currently asserts on ordering.

Let me check whether the store orders ListTriggerPolicies.

Let me check the store's ListTriggerPolicies to see if it's already deterministic.

$ rg -n "ListTriggerPolicies|func.*ListTrigger" internal/store/*.go
internal/store/triggers.go:67:func (s *Store) ListTriggers(ctx context.Context, realmID int16) ([]Trigger, error) {
internal/store/policies_test.go:56:	list, err := s.ListTriggerPolicies(ctx, realm.ID)
internal/store/policies_test.go:58:		t.Fatalf("ListTriggerPolicies: %v", err)
internal/store/policies.go:66:// ListTriggerPolicies returns every policy of a realm ordered by name.
internal/store/policies.go:67:func (s *Store) ListTriggerPolicies(ctx context.Context, realmID int16) ([]TriggerPolicy, error) {

The store's `ListTriggerPolicies` says "ordered by name" (policies.go:66). So the HTTP response *is* already deterministic/stable because the store orders by name. That means `ListPolicies` doesn't actually need its own sort — the ordering comes from the store. So it's not a bug, just a stylistic asymmetry. That's weak — the recipe says don't propose renames/reformatting for their own sake.

Let me look for something more substantive. Let me re-examine the `DeleteDevice` error handling and the `SetAuthKey`.

**`SetAuthKey`** (service.go:437-445): sets a single-element key list. But the doc comment on `GetAuthKey` says "Astrate stores a list for rotation". `SetAuthKey` replaces the whole set with exactly one key, discarding any rotation window. That's a potential real behavior concern — but the HTTP test does rotate to a 2-key set by... wait, `SetAuthKey` only ever stores `[]string{keyPEM}` — a single element. But the HTTP test at http_test.go:238-239 does a PUT with a *concatenated* 2-key PEM and expects it back unchanged.

Let me check: `SetAuthKey` stores `[]string{keyPEM}` where keyPEM is the concatenated string including `\n`. Then `joinPEM` with a single element returns it verbatim. So the concatenation is preserved as ONE element, and `GET` returns it as-is. Interesting — the comment in `GetAuthKey` says the verifier "already splits into individual keys". So storing concatenated as a single string works because neither `Set` splits it nor `Get` splits it — it round-trips.

But wait — the whole purpose of storing a *list* of keys is `SetRealmJWTPublicKeys(ctx, realm, []string{keyPEM})` passing a single concatenated string. So the "list" storage capability (for rotation) is never populated from this endpoint — `SetAuthKey` always stores exactly one (concatenated) string. That's consistent with upstream (PUT replaces the whole set). Not a bug.

Let me reconsider the strongest candidates:

**Candidate A — `DeletePolicy` silently ignores triggers whose definition can't be unmarshalled (`continue`).** A trigger stored in the DB that can't be parsed as `{"policy": ...}` will be skipped, so a policy it references could be deleted anyway, leaving a dangling reference. But `CreateTrigger` always validates via `triggers.Compile`, so stored triggers always have a policy field parseable by `triggerPolicyRef`. Also `triggers.Compile` might read `policy` from a different path... Let me check whether the policy could legitimately not appear at the top level of the JSON. Actually `triggerPolicyRef` only reads top-level `policy`. If `triggers.Compile` accepts a trigger where the policy name comes from a field not at top-level or the JSON is exactly `{"policy": "..."}` — then unmarshal of `triggerPolicyRef` would succeed. Fine.

This is defensive code with an effectively-unreachable branch. Weak.

**Candidate B — Missing test for the `error_handlers` deletion ordering / Multi.** Weak.

Let me look at whether there is a genuine wrong/unguarded behaviour. Let me reconsider `DeleteDevice`:

```go
id, err := deviceid.Parse(deviceID)
if err != nil {
    return fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
}
```

This maps a malformed device ID to `store.ErrNotFound` → 404 "Device not found". The test at dashboard_compat_test.go:126-134 explicitly asserts `DELETE /devices/not-an-id` returns the device-404 envelope. So this is deliberately matching upstream. Not a bug.

Let me look at the `realmID` resolution: `DeleteDevice`, `InstallInterface`, etc. all call `s.realmID(ctx, realm)` which returns `store.ErrNotFound` for an unknown realm → mapped in `writeError` to 404. But the auth middleware normally rejects unknown realms first (401/403). The HTTP comment even notes this. Fine.

Now let me think about **the `policy` revalidation split**. In `CreateTrigger`:
```go
ct, err := triggers.Compile(tn.Name, def)
...
if ct.PolicyName != "" {
    if _, err := s.st.GetTriggerPolicy(ctx, rid, ct.PolicyName); err != nil {
        if errors.Is(err, store.ErrNotFound) {
            return nil, fmt.Errorf("%w: policy %q does not exist in realm %q", ErrValidation, ct.PolicyName, realm)
        }
        return nil, err
    }
}
```
This is the referential integrity check, well covered by `TestPolicyReferentialIntegrity`.

Let me re-focus. The recipe's highest-value categories: (1) wrong/unguarded behaviour, (2) missing test on existing rule.

Let me think about the **`ListPolicies` asymmetry** more. Even though the store orders by name today, the *service* contract implicitly promises nothing. But siblings sort. If a future store change reorders, `ListTriggers` stays stable (service sorts) while `ListPolicies` silently becomes unstable. It's a one-line mechanical change that aligns behaviour and is machine-checkable. This is a legitimate clarity/consistency proposal — it's not a rename or reformat, it's an unguarded ordering contract that siblings enforce. But the test assertion is weak if ordering isn't actually tested. To make it a real machine-checkable task, the task should add a sorting test with multiple out-of-order policies. That would be the "missing test" angle: no test asserts `ListPolicies` returns sorted names, while `ListTriggers`/`ListInterfaces` do. 

Actually this is the strongest *proposal*: a missing test that `GET /policies` returns names sorted, with multiple policies installed out of order. Currently the only test installs one policy. And there's a real asymmetry: `ListPolicies` doesn't sort while every sibling does. A task to add `sort.Strings` + a multi-policy ordering test is clean and machine-checkable.

Let me also reconsider **`DeletePolicy` unmarshal `continue`** once more — is there a genuine gap? Consider a trigger whose `policy` field references the policy but where the definition JSON, after `triggers.Compile` normalization and storage, differs from what `triggerPolicyRef` reads. Both read the JSON `"policy"` top-level key. No gap.

Let me also check `maxBodyBytes` cap consistency — policies/triggers/interfaces all capped at 1MiB. Fine.

Let me consider one more genuine issue: **`CreatePolicy` validates but does not invalidate the engine, with a long comment explaining why** (lines 327-335). The comment reasoning: a policy is unreferenced both at create and delete, so no snapshot effect. This is a deliberate design decision documented. Not a bug. Good.

Now, one more: the `APICompatVersion` is `"1.2.2"` but the dashboard_test comment says "Dashboard v1.2.2 requires". The trigger delivery policies page requires >= 1.1.1. Fine, no issue.

Let me settle on proposals. Given I want genuinely valuable, machine-checkable, at most 5 items. Let me also double check the `writeError` ordering — a `store.ErrNotFound` wrapped in `ErrValidation`? No. Order: ErrValidation, ErrAlreadyExists, ErrInterfaceMajorNotZero, ErrInterfaceInUse, ErrNotFound. If an error wraps both... they don't overlap. Fine.

Let me also verify one thing about `DeleteInterface` — the `majorParam` negative check vs store. Fine.

I think the honest finding is: this package is well-tested and the main genuinely-actionable items are:
1. **`ListPolicies` doesn't sort** while every sibling list method does — missing test for sorted policy list output (multi-policy ordering). This is both a consistency gap and a missing machine-checkable rule.
2. Possibly the `DeletePolicy` silent-`continue` — but I decided it's defensively unreachable.

Let me look a bit more for a "wrong/unguarded" behaviour. Let me re-read `interfacesChanged`/`triggersChanged` — both log and swallow errors after commit. Intended (self-heal). Fine.

Let me double-check the `GetInterface` major for negative: `GetInterface` (service) takes `major` already parsed by `majorParam` which rejects negatives. Fine.

One more real check: **`installInterface` with an empty/synthetic major** — `CreateTrigger` requires name; does `InstallInterface` require the interface to have major 0 for draft vs major >= 1? That's in `interfaceschema.ParseInterface`. Out of scope.

Actually, let me reconsider the **`DeletePolicy` unmarshal-ignore** once more from a different angle: even if unreachable today, the recipe's category 1 is "unguarded". But the comment's whole referential-integrity design (TestPolicyReferentialIntegrity) exists because the checks are real. The `continue` means: if the definition is corrupt, we can't tell if it references the policy, and we proceed with deletion. The two-source-of-truth (definition JSON policy field vs compiled PolicyName) is the actual design smell. But since it only goes through `triggers.Compile` on create, a corrupt stored trigger is impossible via the API. I'll mention it as deliberately-not-proposed.

Let me settle. I'll propose:
1. `realm-listpolicies-sort`: add `sort.Strings` to `ListPolicies` for consistency with its three siblings, plus a test installing multiple policies out of order and asserting `GET /policies` returns them sorted.

That's genuinely one solid machine-checkable task. Let me look for 1-2 more solid ones to make the review worthwhile but I don't want to pad.

Let me reconsider the `UpdateInterface` stored-parse failure branch (lines 140-143): it's another defensive unreachable-path but represents a robustness concern: if for any reason a stored definition is corrupt (e.g. manual DB edit), `UpdateInterface` returns an error instead of doing anything worse — that's *good* handling, not a bug.

Let me think about whether there's a genuine **data race** or **context-ignored** issue. `Service` methods all use ctx. The `Disconnecter` interface `DisconnectDevice` takes no ctx — it's meant for a live session kick, fine.

I think honestly this package is in good shape. Let me check one possibly-missed thing: does `CreateTrigger` / `DeleteTrigger` handle the policy invalidation? The comment at 327-335 says trigger mutations refresh via triggersChanged. OK.

Let me finalise with 2 tightly-scoped proposals:
1. `realm-listpolicies-sorted`: sort ListPolicies + ordering test.
2. Maybe a second: verify there's no test for `ListInterfaceMajors` unknown-name → 404 (service.go:211-213). The HTTP test `ListAndGet` only tests majors with an existing name; no test asserts the 404 path for `ListInterfaceMajors`. That's a small missing test. But is it a "rule that already exists"? Yes — the `ErrNotFound` when no major. It's machine-checkable and cheap. Good second item.

Actually let me verify that 404 path isn't tested.
$ rg -n "interface_names|/interfaces/|ListInterfaceMajors|ErrNotFound" internal/realm/*_test.go
internal/realm/http_test.go:147:		decodeData(t, r.req(t, http.MethodGet, "/interfaces/"+rmIface, "", r.rmaToken), &majors)
internal/realm/http_test.go:151:		if rec := r.req(t, http.MethodGet, "/interfaces/"+rmIface+"/1", "", r.rmaToken); rec.Code != http.StatusOK {
internal/realm/http_test.go:157:		if rec := r.req(t, http.MethodPut, "/interfaces/"+rmIface+"/1", ifaceV1b, r.rmaToken); rec.Code != http.StatusNoContent {
internal/realm/http_test.go:163:		if rec := r.req(t, http.MethodPut, "/interfaces/"+rmIface+"/1", ifaceV1x, r.rmaToken); rec.Code != http.StatusUnprocessableEntity {
internal/realm/http_test.go:173:		decodeData(t, r.req(t, http.MethodGet, "/interfaces/"+rmIface, "", r.rmaToken), &majors)
internal/realm/http_test.go:182:		if rec := r.req(t, http.MethodDelete, "/interfaces/"+rmIface+"/1", "", r.rmaToken); rec.Code != http.StatusUnprocessableEntity {
internal/realm/http_test.go:197:		if rec := r.req(t, http.MethodDelete, "/interfaces/"+rmDraft+"/0", "", r.rmaToken); rec.Code != http.StatusUnprocessableEntity {
internal/realm/http_test.go:204:		if rec := r.req(t, http.MethodDelete, "/interfaces/"+rmDraft+"/0", "", r.rmaToken); rec.Code != http.StatusNoContent {

Confirmed: `GET /interfaces/{name}` (unknown name → 404) is never tested, and `GET /interfaces/{name}/{major}` (unknown major → 404) via `ListInterfaceMajors`/`GetInterface` 404 path also isn't asserted. The `MajorCoexistence` test only checks the happy path.

Let me verify the recipe wants max 5 task lines and that these are worth it. Let me also double-check `ListInterfaces` unknown-name behavior: `ListInterfaces` returns empty list for unknown interface names (no ErrNotFound) since it lists distinct names from load. `ListInterfaceMajors` is the one with the explicit ErrNotFound rule.

Let me settle on my proposals. Let me look at the `go vet`/existing helpers for the test to confirm ordering assertions are feasible. The test uses `contains` not ordering. A sorted test would need an ordering assertion — easy.

Let me finalize the review. I have two solid, machine-checkable, consistency/coverage items:

1. **`ListPolicies` unsorted + no ordering test.** Real inconsistency (only list method without `sort.Strings`), machine-checkable by installing multiple policies out of order and asserting the `GET /policies` response order.

2. **Missing 404 tests for interface lookups.** `ListInterfaceMajors` extra-major / unknown-name → 404 (`service.go:211-213`), and `GetInterface` unknown major → 404. No test asserts these; only happy paths and delete paths are covered.

Both are legitimate. I'll also document what I decided not to propose.

Let me write the review file and append task lines. Today is 2026-09-01.

Let me check the todo.md first to match format.
→ Read .mule/todo.md [limit=40]
→ Read .mule/todo.md [offset=40, limit=70]
I see the format for `[auto]` proposal lines. Note the note at line 21-24: "mule/queue is authoritative for this file... editing the queue on main resurrects completed tasks" — but the recipe for code-review explicitly says "append at most five task lines to `.mule/todo.md`". This is the code-review recipe's deliverable, so I follow the recipe's instruction to append to todo.md.

Now let me write the review file and append the task lines. My chosen proposals:

1. `realm-policy-list-sorted`: `ListPolicies` is the only realm list method without `sort.Strings` (service.go:379), and no test asserts `GET /policies` returns names sorted — install >1 policy out of order and assert sorted order (consistent with ListTriggers/ListInterfaces).

2. `realm-interface-lookup-404`: no test asserts the 404 path for `GET /interfaces/{name}` unknown name / `GET /interfaces/{name}/{major}` unknown major (service.go ListInterfaceMajors:211-213, GetInterface) — add coverage.

Both machine-checkable, both missing-test/consistency category. Let me write the review file.

← Write .mule/reviews/realm-2026-09-01.md
Wrote file successfully.

Now append the task lines to todo.md:
→ Read .mule/todo.md [offset=109]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -107,4 +107,6 @@
 - [!] issue-91-pairing-health-serve-root: in internal/pairing/http.go, add the unauthenticated `GET /pairing/health` route (no realm segment) sharing the handler/payload of the existing `GET /pairing/v1/{realm}/health`, keeping the v1 route as-is; tests for both paths (closes #91) — BLOCKED: gates failed
 - [!] issue-68-async-operation-accepted: on the mutating endpoints upstream surfaces `async_operation` on (housekeeping realm create/delete; realm-management interface install/update/delete, trigger/policy delete in internal/realm/http.go), accept and ignore `?async_operation=false`, with unparseable/`true` values also not changing behaviour (Astrate stays always-sync); tests cover the flag parsing (closes #68) — BLOCKED: gates failed
 - [!] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto] — BLOCKED: gates failed
 - [!] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
+- [ ] realm-policy-list-sorted: in `internal/realm/service.go`, sort the names returned by `ListPolicies` the way `ListTriggers`/`ListInterfaces`/`ListInterfaceMajors` already do (it is the only list method without `sort.Strings`), and extend `TestDashboardCompat.Policies` (or a new subtest in `http_test.go`) to install several policies in non-alphabetical order and assert `GET /policies` comes back sorted. [auto]
+- [ ] realm-interface-lookup-404: in `internal/realm/http_test.go` `TestRealmManagement`, assert the 404 paths that nothing currently covers — `GET /interfaces/<nonexistent-name>` and `GET /interfaces/<name>/<no-such-major>` (service.go `ListInterfaceMajors`:211-213 and `GetInterface`, mapped to 404 via writeError). A rule with no test currently. [auto]

Both task lines appended, review file written. No code changes, no git operations — as required by the proposal recipe. The recipe didn't ask for any `gh issue create`, so no issue filing.

The package `internal/realm` was well within budget (small), not in the previous review list, and I read it entire. Two genuinely machine-checkable items proposed (a consistency/missing-sort gap and a missing 404 test); several plausible candidates deliberately NOT proposed with rationale so the next run doesn't repeat the work.

Done: Proposed sorted `ListPolicies` + interface-lookup-404 tests after a full read of `internal/realm`; wrote review file and 2 task lines, no code/git touched.
Files: .mule/reviews/realm-2026-09-01.md, .mule/todo.md
Verified: n/a — review only, no code to gate
Unsure: `ListPolicies` ordering is currently stable only via the store's ORDER BY; if that store change surprises, the task's test pins it
Follow-ups: realm-policy-list-sorted, realm-interface-lookup-404 (both `[auto]`)
