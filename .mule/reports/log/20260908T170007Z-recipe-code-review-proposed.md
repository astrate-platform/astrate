slug: recipe-code-review
verdict: proposed
at:  66054a7
ran: 2026-09-08T17:00:07Z on DietPi in 170s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md

$ ls /root/astrate-mule/.mule/reviews/ 2>/dev/null || echo "no reviews dir yet"
appengine-2026-09-07.md
broker-2026-09-04.md
engine-2026-09-06.md
flow-2026-09-05.md

$ ls /root/astrate-mule/internal/ /root/astrate-mule/pkg/ 2>/dev/null
/root/astrate-mule/internal/:
appengine
auth
broker
config
engine
flow
flowapi
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil

/root/astrate-mule/pkg/:
astarteapi
deviceid
interfaceschema
payload

$ for pkg in auth config flowapi housekeeping httpx observability pairing realm store swagger testutil; do echo -n "$pkg: "; rg -c '' --glob "internal/$pkg/**/*.go" 2>/dev/null | awk -F: '{s+=$2} END{print s+0}'; done; for pkg in astarteapi deviceid interfaceschema payload; do echo -n "pkg/$pkg: "; rg -c '' --glob "pkg/$pkg/**/*.go" 2>/dev/null | awk -F: '{s+=$2} END{print s+0}'; done
auth: 1532
config: 818
flowapi: 2230
housekeeping: 1311
httpx: 357
observability: 398
pairing: 2849
realm: 2573
store: 6677
swagger: 171
testutil: 697
pkg/astarteapi: 485
pkg/deviceid: 358
pkg/interfaceschema: 3176
pkg/payload: 2673

$ rg -c '' --glob 'internal/store/**/*.go' | sort -t: -k2 -rn | head -15
internal/store/datastreams_test.go:783
internal/store/devices.go:485
internal/store/datastreams.go:434
internal/store/interfaces.go:398
internal/store/devices_test.go:382
internal/store/realms_test.go:355
internal/store/store_test.go:274
internal/store/realms.go:250
internal/store/interfaces_test.go:249
internal/store/groups.go:241
internal/store/store.go:207
internal/store/flows.go:207
internal/store/crypto_test.go:207
internal/store/pipelines.go:193
internal/store/properties.go:183

→ Read internal/store/devices.go
→ Read internal/store/datastreams.go
Let me read the largest implementation files in the store package.
→ Read internal/store/interfaces.go
→ Read internal/store/realms.go
→ Read internal/store/groups.go
→ Read internal/store/store.go
→ Read internal/store/flows.go
→ Read internal/store/pipelines.go
→ Read internal/store/properties.go
→ Read internal/store/devices_test.go
→ Read internal/store/datastreams_test.go
Now let me read the test files.
→ Read internal/store/interfaces_test.go
→ Read internal/store/store_test.go
I've now read the core store files and tests. Let me also read the remaining store files (properties_test, crypto, notify, triggers, etc.) to get a full picture of test coverage and identify what might be missing.

Let me check which files have been reviewed recently to avoid duplication, and confirm I've picked a package (store) not recently reviewed.

$ for f in /root/astrate-mule/.mule/reviews/*.md; do echo "=== $f ==="; head -20 "$f"; done
=== /root/astrate-mule/.mule/reviews/appengine-2026-09-07.md ===
# Code review: internal/appengine — 2026-09-07

Area chosen because recent reviews were `broker` (2026-09-04), `flow`
(2026-09-05), `engine` (2026-09-06); `appengine` is the largest package with no
review file in `.mule/reviews/` (15 files, ~7200 lines). I read the four
non-test source files largest-first plus the unit-tested seams.

## What I read

- `internal/appengine/service.go` (846) — device list/status/patch, groups,
  by-alias/by-group mirrors, applyPatch validation order, group token
  pagination
- `internal/appengine/http.go` (696) — route table, parseQueryOpts, writeError
  taxonomy, links builders
- `internal/appengine/data.go` (514) — resolve/GetData, datastreamData
  (snapshot / series / downsample / object), bigint reshaping, properties
  snapshot
- `internal/appengine/downsample.go` (29) — bucketFor
- tests: `query_opts_test.go`, `downsample_test.go`, `links_test.go`,
  `groups_parity_test.go`, `http_test.go` (grep)
=== /root/astrate-mule/.mule/reviews/broker-2026-09-04.md ===
# Code review: internal/broker — 2026-09-04

## Files read

- `internal/broker/broker.go` (343 lines) — Broker assembly, Close quiesce, intakeHook (deferred ack)
- `internal/broker/authhook.go` (440 lines) — realm CA pools, deviceSession, sessionRegistry, OnPacketRead CN rewrite, OnConnectAuthenticate
- `internal/broker/aclhook.go` (209 lines) — §3.2 checkACL matrix, offlineACL, OnACLCheck
- `internal/broker/sessionstore.go` (533 lines) — bbolt storage hook, inflight Seq/expiry, restore path
- `internal/broker/lifecycle.go` (90 lines) — connect/disconnect rows + LifecycleSink
- `internal/broker/intake.go` (92 lines) — InboundMessage/Intake/Lifecycle contracts
- `internal/broker/publisher.go` (62 lines) — inline Publisher facade
- `internal/broker/identity.go` (104 lines) — ParseCN, SplitTopic, realm-name shape
- Tests: `broker_test.go`, `broker_e2e_test.go`, `aclhook_test.go`, `sessionstore_test.go`, `identity_test.go`
- Also: mochi v2.7.9 `server.go` (`processPublish` ACL-deny semantics, `DisconnectClient`, `expire` derivation, inline `pk.Expiry` computation)

This is the first review of `broker`; the rotation had covered store, realm, flow, engine, pairing.

## What I read in detail

I focused on the three seams where mochi's semantics meet Astrate's rules and a wrong guess would cost data or a security boundary: (a) the ACL cold-cache miss path in `OnACLCheck`, (b) `DisconnectDevice` on device deletion vs the persistent bbolt session, (c) the offline-delivery ACL. I verified mochi internals directly for the claims below rather than positing them.
=== /root/astrate-mule/.mule/reviews/engine-2026-09-06.md ===
# Code review — internal/engine — 2026-09-06

Area chosen because `broker` (2026-09-04) and `flow` (2026-09-05) were recent reviews and `engine` is the largest un-reviewed package (48 files).

## What I read

- `internal/engine/router.go` (engine struct, sharding, Submit, Drain, metrics)
- `internal/engine/engine.go` (fireCommitted / fireData / fireDataChanges / fireDevice / handleLifecycle)
- `internal/engine/data.go` (validation pipeline, PersistOp, trackPrevious)
- `internal/engine/batch.go` (micro-batcher, row conversion, encodeValueJSON/jsons constraints)
- `internal/engine/cache.go` (schema + device caches, buildRealm trigger/policy load)
- `internal/engine/control.go` (emptyCache / producer+consumer properties, zip-bomb guard)
- `internal/engine/introspection.go`, `topics.go`, `capabilities.go`, `serverdata.go`
- `internal/engine/triggers/{match,actions,events,policy}.go`
- `internal/engine/stream/bus.go`
- tests: `serverdata_test.go`, `devicedata_test.go`, `triggers_change_test.go`

## What I found (worth proposing)

1. **Server- and device-owned publish paths never fire data triggers nor fan out to the live bus.** `publishAsOwner` (internal/engine/serverdata.go:94) persists the op (`UpsertProperty` / `AppendDatastreams`) and (for `PublishServerValue`) publishes to the broker, but never calls `fireData`/`fireCommitted` and never does `e.bus.Publish`. The broker-driven ingress path gets both through `finalize` → `afterCommit` (batch.go:186) → `fireCommitted` (engine.go:140). So:
=== /root/astrate-mule/.mule/reviews/flow-2026-09-05.md ===
# Code review: internal/flow — 2026-09-05

## Files read

- `internal/flow/message.go` (364 lines) — Message JSON wire encode/decode, scalar + map typed-data
- `internal/flow/flow.go` (332 lines) — Flow/Manager lifecycle, StartFlow, StopFlow, Shutdown, source pump
- `internal/flow/router.go` (284 lines) — lane-based Router, Submit QoS/overflow policies, Drain, FNV-1a laneOf
- `internal/flow/factory.go` (249 lines) — block Registry, Instantiate topo order, stopAll cleanup, ParseDefinition
- `internal/flow/pipeline.go` (165 lines) — Pipeline structure, Validate, MarshalJSON
- `internal/flow/composite.go` (156 lines) — UserBlock composite expansion (#85), nesting depth cap, cycle detection
- `internal/flow/substitute.go` (126 lines) — ${config.key} placeholder substitution
- `internal/flow/block.go` (113 lines) — Block/Source/Stopper interfaces + func adapters
- `internal/flow/graph.go` (64 lines) — linear BlockGraph, Run, Sources
- Tests: `pipeline_test.go`, plus targeted greps of router/flow/composite/factory/message tests

This is the first review of `flow`; the rotation checked broker last (2026-09-04).

## What I read in detail

I went through every non-test source file, largest first, and paid closest attention to the

✱ Glob "internal/store/*.go" in . · 32 matches
→ Read internal/store/triggers.go
→ Read internal/store/crypto.go
→ Read internal/store/notify.go
→ Read internal/store/policies.go
Recent reviews covered broker (09-04), flow (09-05), engine (09-06), appengine (09-07). The broker review mentions "the rotation had covered store, realm, flow, engine, pairing" in the past, but there's no `store-*.md` review file present. Store is the largest unreviewed package. Let me read the remaining store files that I haven't yet (crypto, notify, triggers, policies, userblocks, retention, and the remaining tests) to complete the review.
→ Read internal/store/userblocks.go
→ Read internal/store/retention.go
→ Read internal/store/crypto_test.go
I now have a comprehensive view of the store package. Let me review the remaining test files to know what's already covered, and look for interesting gaps.

Let me look at the areas where I found potential issues. Let me re-examine a few specific findings:

1. **`UpdateIntrospection` doesn't Notify** — In `interfaces.go`, `InstallInterface`/`UpdateInterface`/`DeleteInterface` change interface data, but I don't see `NotifyInterfacesChanged` being called anywhere. Let me check if these methods notify the engine.

2. **`patchDeviceJSONB` with empty removals → `$3::text[]` empty array** — potential edge.

3. **`ListDevices` unbounded** — it has a limit param, good.

Let me check the crypto_test and pipelines_validate_test for coverage, and look for any obvious untested rules.
→ Read internal/store/pipelines_validate_test.go
→ Read internal/store/properties_test.go
I now have a thorough read of the store package. Let me identify the strongest candidates to propose.

Let me re-examine a few specific things I noticed:

**Finding A — `SetPayloadFormatHint` validation vs DB**: Fine, validated in Go.

**Finding B — `UpdateIntrospection` transaction + removed returning**: covered by tests.

Let me look more carefully at potential real findings. Let me check `patchDeviceJSONB` — when `patch` is empty (both removals and additions empty), it still runs an UPDATE with `(col - '{}'::text[]) || '{}'::jsonb`. That's harmless. Not worth proposing.

Let me look at `Downsample`/`DownsampleLTTB` argument ordering in the SQL: In `Downsample`, the `bucket` param is appended AFTER the `where()` args. `where()` produces args `[realm, uuid, iface, path]` and `fmt.Fprintf("... make_interval(secs => $%d) ... %s ...", len(args), where)` → uses `$5` for bucket since len(args)=4 before appending. Correct. Covered by tests.

**Finding C — `SeriesQuery.where()` and `tail()`** — `tail()` appends `LIMIT` directly as integer. Fine.

Let me look at the `UpdateFlowRuntime` — `started_at = COALESCE($6, started_at)`. When starting a flow, caller passes startedAt; when stopping, it keeps existing. That's a design choice.

Now, the notable potential finding: **`DeleteDevice` and `DeleteRealm` sweep only `individual_datastreams` and `object_datastreams` hypertables. But `properties` cascade via FK — good. Groups via FK — good.**

Let me check the strongest candidate I actually found. Let me reconsider whether `UpdateIntrospection`'s removed list should NOTIFY. Let me check how `NotifyInterfacesChanged` is used — it's for interface CRUD. Let me grep who calls InstallInterface/UpdateInterface/DeleteInterface and whether they notify.

Actually, the recipe says to propose only what I found by reading, with file+line, in the priority order. Let me look for genuinely wrong/unguarded behavior.

Let me reconsider `patchDeviceJSONB`: the `removals` argument `$3::text[]`. If `patch` contains only additions (no removals), `removals` is an empty slice `[]string{}`. Postgres `- '{}'::text[]` removes nothing — fine. If only removals, `additions` is `{}` jsonb, `|| '{}'` no-op. Fine.

**Finding D — `PurgeDeviceOwnedExcept` with empty keep** — covered in test.

Let me focus on the one thing that stands out as a genuine gap in test coverage of an existing rule. Let me look at `ApplyGlobalRetention` and `EnforceRealmRetentionCeilings` — these are retention-related. `testRetention` covers ApplyGlobalRetention. Is there a test for `EnforceRealmRetentionCeilings` (#72)? Let me check.

$ cd /root/astrate-mule && rg -l "EnforceRealmRetentionCeilings" --glob '*.go' && echo "---notify callers---" && rg -n "NotifyInterfacesChanged" --glob '*.go' && echo "---SortIndices algo ret---" && rg -n "EnforceRealmRetention|ApplyGlobalRetention" --glob '*.go'
cmd/astrate/main.go
internal/store/retention.go
internal/store/retention_test.go
---notify callers---
internal/store/notify_test.go:41:				if err := s.NotifyInterfacesChanged(context.Background(), realmID); err != nil {
internal/store/notify_test.go:42:					t.Fatalf("NotifyInterfacesChanged: %v", err)
internal/store/notify.go:31:// NotifyInterfacesChanged emits a ChannelInterfaces notification for realmID.
internal/store/notify.go:32:func (s *Store) NotifyInterfacesChanged(ctx context.Context, realmID int16) error {
internal/engine/engine_integration_test.go:388:	if err := rig.raw.NotifyInterfacesChanged(ctx, rig.realmID); err != nil {
internal/engine/engine_integration_test.go:389:		t.Fatalf("NotifyInterfacesChanged: %v", err)
internal/realm/service.go:665:	if err := s.st.NotifyInterfacesChanged(ctx, rid); err != nil {
internal/realm/service.go:678:	if err := s.st.NotifyInterfacesChanged(ctx, rid); err != nil {
---SortIndices algo ret---
cmd/astrate/main.go:101:		if err := st.ApplyGlobalRetention(ctx, d); err != nil {
cmd/astrate/main.go:197:			if err := st.EnforceRealmRetentionCeilings(ctx); err != nil {
internal/store/store.go:168:// ApplyGlobalRetention sets a global drop-chunks retention policy on both
internal/store/store.go:173:func (s *Store) ApplyGlobalRetention(ctx context.Context, d time.Duration) error {
internal/store/datastreams_test.go:762:	if err := s.ApplyGlobalRetention(ctx, 90*24*time.Hour); err != nil {
internal/store/datastreams_test.go:763:		t.Fatalf("ApplyGlobalRetention(set): %v", err)
internal/store/datastreams_test.go:770:	if err := s.ApplyGlobalRetention(ctx, 30*24*time.Hour); err != nil {
internal/store/datastreams_test.go:771:		t.Fatalf("ApplyGlobalRetention(re-set): %v", err)
internal/store/datastreams_test.go:777:	if err := s.ApplyGlobalRetention(ctx, 0); err != nil {
internal/store/datastreams_test.go:778:		t.Fatalf("ApplyGlobalRetention(clear): %v", err)
internal/store/retention_test.go:13:// TestEnforceRealmRetentionCeilings covers the realm-level retention sweep
internal/store/retention_test.go:17:func TestEnforceRealmRetentionCeilings(t *testing.T) {
internal/store/retention_test.go:93:	if err := s.EnforceRealmRetentionCeilings(ctx); err != nil {
internal/store/retention_test.go:94:		t.Fatalf("EnforceRealmRetentionCeilings: %v", err)
internal/store/retention_test.go:117:	if err := s.EnforceRealmRetentionCeilings(ctx); err != nil {
internal/store/retention_test.go:118:		t.Fatalf("EnforceRealmRetentionCeilings(re-run): %v", err)
internal/store/retention.go:8:// EnforceRealmRetentionCeilings applies every realm's
internal/store/retention.go:16:func (s *Store) EnforceRealmRetentionCeilings(ctx context.Context) error {

`EnforceRealmRetentionCeilings` is tested. `NotifyInterfacesChanged` is called by the realm service on interface CRUD, not by the store methods directly — that's an architectural choice, fine.

Now let me look at one concrete thing I want to verify: In `realms.go` `UpdateRealm`, the retention handling. There's a subtle rule:
```go
if p.ClearRetention || (p.PatchRetention && p.SetRetention <= 0) {
    add("datastream_maximum_storage_retention", nil)
} else if p.PatchRetention {
```
Note `ClearRetention` alone clears but if `PatchRetention` is false and `ClearRetention` is false, nothing is added for retention. But there's a subtle precedence: if `ClearRetention` is true AND `PatchRetention` is true with a positive value, Clear wins (correct per doc "Clear beats Set"). Fine.

Let me look at the `patchDeviceJSONB` empty case more concretely with a proposal candidate. Actually, let me look at a genuine concern I noticed: **`RegisterDevice` ON CONFLICT — when the device already exists and `first_credentials_request IS NULL`, it rotates the secret and sets status='registered'. But it does NOT clear `cert_serial`/`cert_aki`. However, if first_credentials_request is NULL, cert would be from a previous unregister... Actually `UnregisterDevice` clears those. Fine.**

Let me consider the strongest genuinely-found issue. I noticed in `flows.go` `UpdateFlowRuntime`: `error_message` and `failed_block` semantics. Not an issue.

Let me look at what's genuinely most valuable. Let me reconsider **`SetPayloadFormatHint`** — there is no test that GetDevice reads the hint default "bson" (tested) and the flip (tested) and invalid (tested).

Now the strongest candidate worth proposing (missing test on existing rule): Looking at `devices.go` line 365/371 — `patchDeviceJSONB` applies merge-patch. There IS a test (`AliasesAndAttributes`) covering replace and remove. 

Let me look for a rule that has NO test. `AliasValuesTaken` — is it tested? Let me grep.
$ cd /root/astrate-mule && echo "=== AliasValuesTaken ===" && rg -n "AliasValuesTaken" --glob '*.go'; echo "=== SetDeviceDisconnected test ===" && rg -n "SetDeviceDisconnected|AddDeviceStats|DeviceStats" internal/store/*_test.go; echo "=== LatestIndividual ===" && rg -n "LatestIndividual" --glob '*.go'
=== AliasValuesTaken ===
internal/store/devices.go:340:// AliasValuesTaken reports whether any device other than id in the realm
internal/store/devices.go:343:func (s *Store) AliasValuesTaken(ctx context.Context, realmID int16, id deviceid.ID, values []string) (bool, error) {
internal/appengine/service.go:452:	return s.st.AliasValuesTaken(ctx, rid, id, values)
=== SetDeviceDisconnected test ===
internal/store/devices_test.go:97:		if err := s.AddDeviceStats(ctx, realm.ID, id, 3, 1024); err != nil {
internal/store/devices_test.go:100:		if err := s.AddDeviceStats(ctx, realm.ID, id, 2, 512); err != nil {
internal/store/devices_test.go:116:		if err := s.SetDeviceDisconnected(ctx, realm.ID, id, at.Add(time.Hour)); err != nil {
internal/store/devices_test.go:332:		total, connected, err := s.DeviceStats(ctx, realm.ID)
internal/store/devices_test.go:334:			t.Fatalf("DeviceStats: %v", err)
internal/store/devices_test.go:376:		total, connected, err = s.DeviceStats(ctx, realm.ID)
=== LatestIndividual ===
internal/engine/data.go:432:		r, err := e.st.LatestIndividual(ctx, op.RealmID, op.DeviceID, op.Interface.ID, op.Path)
internal/engine/fakestore_test.go:379:// LatestIndividual reads the newest sample of one series.
internal/engine/fakestore_test.go:380:func (f *fakeStore) LatestIndividual(_ context.Context, realmID int16, deviceID deviceid.ID, interfaceID int64, path string) (*store.IndividualRow, error) {
internal/engine/cache.go:56:	// LatestIndividual reads the newest sample of one datastream series —
internal/engine/cache.go:58:	LatestIndividual(ctx context.Context, realmID int16, deviceID deviceid.ID, interfaceID int64, path string) (*store.IndividualRow, error)
internal/store/datastreams.go:243:// LatestIndividual reads the newest sample of one individual-datastream
internal/store/datastreams.go:247:func (s *Store) LatestIndividual(ctx context.Context, realmID int16, deviceID deviceid.ID, interfaceID int64, path string) (*IndividualRow, error) {

I found two genuinely uncovered rules worth proposing:

1. **`AliasValuesTaken` (devices.go:343) has no test.** It's a real DB-backed rule (the `alias_already_in_use` ownership check) with several branches: no other device has the value → false; another device carries one of the values → true; the same device (excluded via `d.id <> $2`) with the value → false (must not count itself). This is a fully machine-checkable missing-test rule. This is high-value category #2 (missing test on a rule that already exists).

2. **`LatestIndividual` (datastreams.go:247) has no direct store test.** It's a real rule (previous-value lookup for value_change triggers) with error branch (ErrNotFound when series empty) plus "ORDER BY ts DESC returns the newest" and the value type returned. It's exercised indirectly through the engine, but the ordering-tiebreak behavior (two rows with same series, one newer — which one wins) and the ErrNotFound for an empty series are not asserted at the store level.

Let me verify these aren't tested elsewhere in store tests, and check the engine tests don't already assert them (they use a fakeStore anyway, so no real coverage).

Let me check whether `AliasValuesTaken` and `LatestIndividual` appear in any store test.
$ cd /root/astrate-mule && echo "=== AliasValuesTaken in any test ===" && rg -c "AliasValuesTaken" internal/store/*_test.go 2>/dev/null; echo "=== LatestIndividual in store tests ===" && rg -c "LatestIndividual" internal/store/*_test.go 2>/dev/null; echo "(empty = not tested at store level)"
=== AliasValuesTaken in any test ===
=== LatestIndividual in store tests ===
(empty = not tested at store level)

Confirmed: neither `AliasValuesTaken` nor `LatestIndividual` is tested at the store level. These are the two strongest candidates (category #2 — missing tests on rules that already exist, fully machine-checkable).

Let me also double-check `SetDeviceDisconnected`/`ListDevices` are covered (they are), and decide on a small third candidate. Let me consider whether there's anything worth proposing in the "wrong/unguarded" category. 

Re-examining `patchDeviceJSONB` (devices.go:377-404): when the patch is empty (no removals, no additions), it still issues an UPDATE with no effective change, but that's harmless and the caller (API) wouldn't send an empty patch. Not worth proposing.

One more: `Downsample`/`DownsampleLTTB` share the same `SeriesQuery.where()` which handles Since/SinceAfter/To — covered by tests.

Let me look at the `SeriesQuery.where()` combined with `Downsample`'s `$%d` — with multiple filters, the `$%d` for bucket still points correctly since args appended after where. Fine.

I'm confident in my two main findings. Let me write the review file and the todo entries. Today's date is 2026-09-08.

Let me check the review naming convention and todo format. Review file: `.mule/reviews/store-2026-09-08.md`. Todo lines follow `- [ ] <slug>: <outcome in one line naming file>`.

Let me also verify whether there are other untested store exports quickly to round out the review's "what I decided not to propose" section.
→ Read .mule/todo.md
I have all the context I need. None of my proposed areas were recently reviewed. Note the todo.md queue doesn't currently have a `store-*` review, and store is the largest un-reviewed package (other reviews covered broker, flow, engine, appengine).

Now let me look at how the recipe handles queue ordering (the note says "mule/queue is authoritative"), but my instructions say I should append new task lines to `.mule/todo.md`. Since this is a proposal job, I'll append to the file as instructed by the parent task, keeping it consistent with the existing `[auto]` convention and naming files.

Let me finalize my findings. My strongest two:

1. **`AliasValuesTaken` (devices.go:343) — no test.** Fully machine-checkable missing-test rule. Branches: other-device-owns-value → true; no-device → false; self-owns-value (excluded) → false.

2. **`LatestIndividual` (datastreams.go:247) — no store test.** Procured only via engine's fake store. No direct assertion of: newest-wins ordering, correct value column returned, ErrNotFound on empty series.

Both are category #2 (missing test on existing rule) which is the best kind. I decided not to propose anything in the wrong/unguarded category — nothing I read was genuinely wrong or unguarded; the package is well-written and thoroughly tested at the store level.

Let me write the review file and append the two task lines.

← Write .mule/reviews/store-2026-09-08.md
Wrote file successfully.

Now let me append the two task lines to `.mule/todo.md`. I'll add them at the end of the queue, following the existing line format and `[auto]` convention (these came from a proposal recipe run).
→ Read .mule/todo.md [offset=122, limit=2]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -120,4 +120,7 @@
 - [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
 - [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
 - [x] detailed-listing-required-encrypted-flags: emit the `required`/`encrypted` mapping flags in `writeDetailedMapping` (internal/realm/detailed.go) so `GET /realmmanagement/v1/{realm}/interfaces?detailed=true` round-trips them — the detailed listing currently materializes a fixed field set that drops both flags even when a datastream mapping declares one, while the plain GET preserves them via the stored definition. Datastream branch only (properties reject both at parse); add a T1 assertion in internal/realm/detailed_test.go.
+
+- [ ] store-alias-values-taken-test: add a store test for `AliasValuesTaken` (internal/store/devices.go:343, the `alias_already_in_use` ownership check behind PATCH /devices) pinning the three branches that currently have no direct coverage — no other device in the realm carries any of the values → false; another device carries one → true; the same device (excluded via `d.id <> $2`) carries the value → false (must not count itself). Add to internal/store/devices_test.go. [auto]
+- [ ] store-latest-individual-test: add a store test for `LatestIndividual` (internal/store/datastreams.go:247, the previous-value lookup feeding value_change/path_created triggers) — only exercised through the engine's fakeStore today; assert (a) of two samples on one series, `ORDER BY ts DESC LIMIT 1` returns the newest with the right value column, (b) an empty series yields ErrNotFound. Add to internal/store/datastreams_test.go. [auto]

Done. The recipe is a proposal job, so I wrote no code and didn't touch git.

```
Done: Reviewed internal/store (largest un-reviewed package) — found two missing-test gaps, proposed them as queue lines
Files: .mule/reviews/store-2026-09-08.md, .mule/todo.md
Verified: rg coverage searches for AliasValuesTaken / LatestIndividual across internal/store tests -> both absent
Unsure: nothing
Follow-ups: 2 task lines appended to .mule/todo.md ([auto] store-alias-values-taken-test, store-latest-individual-test)
```

Full deliverables:
- `.mule/reviews/store-2026-09-08.md` — what I read, what I propose, and what I decided not to propose and why, so the next run doesn't repeat it.
- Two task lines appended to `.mule/todo.md`:
  - `store-alias-values-taken-test` — test `AliasValuesTaken` (devices.go:343): no other device → false, another device carries one → true, same device (excluded) → false.
  - `store-latest-individual-test` — test `LatestIndividual` (datastreams.go:247): newest-wins ordering + right value column, and ErrNotFound on empty series.

No `gh issue create` was involved — the recipe only calls for the review file + todo lines, so nothing was filed.
