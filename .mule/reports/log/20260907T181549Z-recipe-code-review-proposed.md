slug: recipe-code-review
verdict: proposed
at:  ab49dcc
ran: 2026-09-07T18:15:49Z on DietPi in 153s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
Now let me check recent reviews and the log.
→ Read .mule/log.md

$ ls /root/astrate-mule/.mule/reviews/ 2>/dev/null || echo "no reviews dir"
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

→ Read .mule/todo.md
The recipe says to check `.mule/log.md` and `.mule/reviews/` for which areas were done recently. I need to pick one area/package that was NOT reviewed recently. Recent reviews: `broker-2026-09-04`, `flow-2026-09-05`, `engine-2026-09-06`. Also the log shows many docs-sync and other tasks.

Let me look at the recent reviews to see what was covered.
→ Read .mule/reviews/engine-2026-09-06.md
→ Read .mule/reviews/flow-2026-09-05.md
→ Read .mule/reviews/broker-2026-09-04.md
So reviewed recently: broker (09-04), flow (09-05), engine (09-06). Today is 09-07. I need to pick an area NOT reviewed recently. The rotation covered: store, realm, flow, engine, pairing, broker per the broker review note. Not yet reviewed here: `appengine`, `auth`, `config`, `housekeeping`, `httpx`, `observability`, `pairing` (mentioned as covered in rotation), `realm` (covered), `swagger`, `testutil`, and `pkg/*`.

Let me check sizes to pick a good un-reviewed area. Let me look at the map.

$ cd /root/astrate-mule && for p in appengine auth config housekeeping httpx observability pair realm store swagger; do n=$(ls internal/$p/*.go 2>/dev/null | wc -l); c=$(rg -c '' --glob "internal/$p/**/*.go" 2>/dev/null | awk -F: '{s+=$2} END {print s}'); echo "$p files=$n lines=$c"; done
appengine files=15 lines=7238
auth files=8 lines=1532
config files=2 lines=818
housekeeping files=5 lines=1311
httpx files=4 lines=357
observability files=5 lines=398
pair files=0 lines=
realm files=9 lines=2544
store files=32 lines=6677
swagger files=2 lines=171

The broker note said the rotation covered "store, realm, flow, engine, pairing" but none of those have review files here (only broker, flow, engine have review md files). The log shows lots of docs-sync and engine/probe work, but reviews specifically: broker, flow, engine.

To avoid re-reviewing what was just done, and pick a not-recently-reviewed area, `appengine` is the largest un-reviewed package (7238 lines, 15 files). The recent docs-sync tasks touched appengine's HTTP spec but not the code deeply. Let me also check `store` (6677 lines, 32 files) — though broker said store was covered in rotation, there's no review file for it in `.mule/reviews/`, so it might have been reviewed before this directory's records. 

Given the recent reviews are broker, flow, engine — and appengine is large and wasn't in the reviewed set with a file — I'll pick `appengine`. It's the biggest un-reviewed-by-file package. Let me look at its structure.
$ cd /root/astrate-mule && ls -la internal/appengine/*.go
-rw-r--r-- 1 root root  4598 Jul 27 05:10 internal/appengine/dashboard_compat_test.go
-rw-r--r-- 1 root root 17979 Sep  4 20:46 internal/appengine/data.go
-rw-r--r-- 1 root root   652 Jul 27 05:10 internal/appengine/downsample.go
-rw-r--r-- 1 root root  1855 Jul 27 05:10 internal/appengine/downsample_test.go
-rw-r--r-- 1 root root 11461 Jul 27 22:18 internal/appengine/e2e_test.go
-rw-r--r-- 1 root root 15321 Sep  4 20:46 internal/appengine/formats_parity_test.go
-rw-r--r-- 1 root root 10462 Sep  4 20:46 internal/appengine/groups_parity_test.go
-rw-r--r-- 1 root root 24908 Sep  4 20:46 internal/appengine/http.go
-rw-r--r-- 1 root root 18539 Sep  4 20:46 internal/appengine/http_test.go
-rw-r--r-- 1 root root  2166 Jul 27 05:10 internal/appengine/links_test.go
-rw-r--r-- 1 root root  9611 Sep  4 20:46 internal/appengine/mirror_test.go
-rw-r--r-- 1 root root  3819 Sep  4 20:46 internal/appengine/patch_test.go
-rw-r--r-- 1 root root  5826 Sep  4 20:46 internal/appengine/query_opts_test.go
-rw-r--r-- 1 root root 29319 Sep  4 20:46 internal/appengine/service.go
-rw-r--r-- 1 root root  1717 Sep  4 20:46 internal/appengine/writeerror_test.go

The core non-test files are `data.go`, `http.go`, `service.go`. Let me read them, largest first: service.go (29k), http.go (24.9k), data.go (18k). These are large. Given my ~30k token constraint, I'll read these three carefully.

Let me read service.go first.
→ Read internal/appengine/service.go
service.go is well-written and heavily documented. Now let me read http.go and data.go.
→ Read internal/appengine/http.go
http.go is well-written too. Now let me read data.go, the third core file.
→ Read internal/appengine/data.go
→ Read internal/appengine/query_opts_test.go
I've read the three core files of the appengine package. Let me look at the `downsample.go` and `links_test.go` briefly, and check the QueryOpts `Limit` handling. Let me consider a few things I noticed.

Let me note observations so far:

1. **`datastreamData` interface-root snapshot path** ignores `Limit`, `Since`, `SinceAfter`, `To`, and `Descending`. On line 138, when `path == "" && opts.DownsamplePoints == 0 && aggregation != object`, it calls `IndividualSnapshot` which returns the *latest sample for every endpoint* — ignoring the requested `since`/`to`/`limit` params. Is that upstream behaviour? This is a documented "data-snapshot view". But the comment encodes it as the design. This matches the documented deviation. The question is whether ignoring since/to/limit is intended. This is really upstream's behavior question, not a clear bug.

2. **`applyBigIntegerOpts`** — AllowBigIntegers *nil* means leave as stored. This is documented. Fine.

3. **`deviceListLinks` self link** — the `self` link echoes the request query. Fine.

Let me look at the `nestTree` function with `slices` import and `renderObject` `docs[0]` indexing. In `renderObject`, `docs[0]` is accessed — but this is only called after the empty-check (`len(rows) == 0` returns earlier). So docs is guaranteed non-empty. Good.

4. **`bodyGroup`/`createGroup`** — the 201 echo carries `body`, which includes GroupName and Devices. The comment says the echo carries the original body list. Fine.

Let me check `parseQueryOpts` — the `limit` param. There's a subtle issue: upstream `limit` — what's max? Not relevant.

Let me focus on potential real findings. Let me look at `parseGroupToken` / `groupTokenFor` and the `Note` about `int` offset. On 32-bit... wait the mule runs on 64-bit. But `parseGroupToken` returns `int(binary.BigEndian.Uint32(...))` which on a 32-bit platform could overflow. But this is minor and not on the runtime platform.

Let me look more carefully at potential real value. Let me check the tests coverage for these paths. Let me look at what `query_opts_test.go` and downstream tests cover, and check whether the individual datastream query's `Limit` is combined with `Descending` correctly in store.

Actually, let me reconsider. The strongest candidate from reading: In `datastreamData`, the interface-root snapshot (path == "") silently ignores `since`, `since_after`, `to`, and `limit` query params. The handler `parseQueryOpts` reads and validates these, so a client that asks `GET .../interfaces/{iface}?since=...&limit=50` gets the full snapshot tree with no filtering, silently. For the downsample case there's an explicit refusal ("requires an endpoint path"), and for properties there's explicit refusal of downsample. But for the snapshot view, since/to/limit are silently dropped. This mirrors what the code comments call "the silent lie datastreamData refuses below" — the codebase has a clear principle of refusing rather than silently ignoring query params. Yet the snapshot path silently ignores since/to/limit.

However — is this upstream's actual behavior? The snapshot view is described as a documented deviation ("there is no upstream table rendering of this tree to copy"). Upstream's `GET /devices/{id}/interfaces/{iface}` with no path is a real query that should respect since/limit. Hmm. This is worth verifying against upstream but I can't run it without the Legion Go. This is a plausible category-1 finding (silent ignore), but I should be careful.

Let me verify by checking the docs for this endpoint. Actually, let me check whether this is genuinely a "silent lie" — the codebase's own standard (visible in downsample handling) is: refuse rather than silently ignore. Let me check the upstream semantics. The AppEngine `GET /devices/{id}/interfaces/{iface}` converges to `shown limits` — it does apply since/limit in upstream. Actually in upstream, datastream root query does apply pagination limits.

This is getting into uncertain territory about upstream. Let me look for something more clearly machine-checkable — a missing test, which the recipe ranks #2 and is fully machine-checkable.

Let me look at `parseGroupToken`'s placeholder — is there test coverage for the group pagination? `links_test.go` exists. Let me check what test files cover group token and the `groupTokenFor` round-trip, and the bigint opts.

Actually, let me reconsider the strongest, clearest finding. Let me re-examine `renderObject` table format and `applyBigIntegerOpts` for a subtle bug: `applyBigIntegerOpts` only rewrites the value when a flag is explicitly set. When `AllowBigIntegers` is nil AND `AllowSafeBigIntegers` is nil, it leaves stored representation untouched — but the docs say "longinteger as a decimal string". Actually the object documents store numbers as json.Number via `UseNumber()`. So without flags they re-marshal as-is.

But wait — the condition at line 230 in `datastreamData`: `if opts.AllowBigIntegers != nil || opts.AllowSafeBigIntegers != nil` — it only calls objectColumnTypes (which re-parses the stored interface) when a bigint flag is set. That's an optimization that avoids re-parsing. Fine.

Let me look at the downsample note at line 189-193: `bucket = time.Microsecond` fallback. Fine.

Let me check the `parseGroupToken` returning `int` — actually let me check if there's a test for the token overrun, and reconsider the pagination. This is minor.

Let me look at a cleaner potential finding: **`listDevices`/`ListDevices` — the `details` pagination**. In `ListDevices`, the page.Next logic: `if len(devs) > limit { page.Next = devs[limit-1].ID.String() }`. This is cursor-based where the cursor is the `limit-1`-th element (the last returned), and it's *exclusive* (`after` cursor, documented as "the cursor (exclusive)"). But setting Next to `devs[limit-1]` — the last element of the returned page — and the next request passes it as `after` (exclusive), which would exclude that last element. So it returns limit elements per page, cursor = last element, and the next page starts *after* the last element — correctly returning the next `limit` elements. That's correct (keyset pagination by ID). OK.

Let me look at whether there's a test gap. Let me check `query_opts_test.go` and `downsample_test.go` and think about what's not covered.

Actually, let me step back and pick the clearest finding. The strongest candidates:

**A. Datastream interface-root snapshot silently ignores query params (since/to/limit).** This is consistent with the codebase's own anti-silent-ignore principle, but I can't verify upstream behavior without the Legion Go, and the code explicitly documents this as the "data-snapshot" design. Risk of false positive.

Let me instead verify against docs. Let me check `docs/` or the API yaml for this endpoint's params.

Given my token budget and the requirement to be concrete with file:line, let me examine the test coverage to find a clean missing-test finding. Let me look at the tests quickly.
→ Read internal/appengine/downsample_test.go
→ Read internal/appengine/links_test.go
Let me check the `formats_parity_test.go`, `groups_parity_test.go`, and `patch_test.go` to understand what's already covered, and look at http_test.go to understand the coverage level. Given my token budget, let me be targeted.

Let me focus on the strongest concrete finding I have. Let me reconsider the datastream root snapshot silently dropping query params. Let me check the docs to see whether this endpoint is documented with since/limit params, and think about whether this is worth proposing.

Actually, I have a cleaner, more provable candidate. Look at the interface-root snapshot path more carefully:

```go
if path == "" && opts.DownsamplePoints == 0 && r.iface.Aggregation != interfaceschema.AggregationObject {
    rows, err := s.st.IndividualSnapshot(ctx, r.rid, r.id, r.iface.ID)
```

This ignores `opts.Since`, `opts.To`, `opts.Limit`, `opts.Descending` — all query params the handler parsed and validated. So `GET /devices/{id}/interfaces/{iface}?since=...&to=...&limit=10` silently returns the full latest-snapshot with none of these applied. The codebase's own documented principle is "refuse rather than silently ignore" (the downsample-to-on-properties and downsample-to-requires-path rejections). This is a genuine candidate.

But it's a design decision that needs upstream confirmation, and the comment explicitly frames the snapshot as ignoring format (a documented deviation). Ignoring since/limit isn't documented though — only format is.

Let me verify the docs on this endpoint to reduce false-positive risk. Let me look at the appengine yaml for the interface GET.

$ cd /root/astrate-mule && rg -n "interfaces/\{interface\}" docs/api/astarte_appengine_api.yaml | head; echo "---"; rg -n "since|limit|downsample" docs/api/astarte_appengine_api.yaml | head -40
319:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
361:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
562:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
599:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
1067:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}:
1105:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}:
---
49:        - name: limit
70:                      self: "/v1/test/devices?details=false&limit=1"
86:                      self: "/v1/test/devices?details=true&limit=1"
87:                      next: "/v1/test/devices?details=true&from_token=dT6hS2W9TT6LEnP25ks_lg&limit=1"
876:        - name: limit
1304:      name: since
1309:      description: Return data points since this timestamp (RFC 3339).
1312:      name: since_after
1328:      name: limit
1335:      name: downsample_to

→ Read docs/api/astarte_appengine_api.yaml [offset=562, limit=80]
So the interface-root GET is documented as accepting `DataSince`, `DataSinceAfter`, `DataTo`, `DataLimit`, `DataDownsample`, `DataSort` — the full param set. So a client can legitimately request `GET .../interfaces/{interface}?since=...&limit=10`, but `datastreamData`'s snapshot path (path == "") silently ignores all of them. This is exactly the "silent lie" pattern the code itself refuses elsewhere (downsample on properties, downsample requiring a path, object series 404).

However — I should be careful. The snapshot path is only reached for an individual datastream. Upstream behavior for `GET /interfaces/{iface}` (no path) on an individual datastream... In upstream, this returns the latest snapshot. Whether upstream honors `since`/`limit` there is the question. The doc on this endpoint applies to both interface-root and path level.

Actually, given the documented deviation note (the code explicitly documents format being ignored on the snapshot), and that this is exactly the codebase's anti-silent-lie principle, this is a legit category-1 proposal. But it may warrant upstream verification on the Legion Go. The recipe about performance says propose measurement before optimisation; for correctness, I can propose the task that verifies + fixes.

Actually, let me reconsider whether the snapshot honoring since/to makes sense. The snapshot is "latest sample for every endpoint" — that's inherently not time-filterable (it's the latest value, not a window). Honoring `limit` and `to` for a single snapshot per endpoint is meaningless. So actually upstream likely just ignores these on the root-snapshot too, OR refuses them. The real question is whether it should refuse them (422) rather than silently returning. Given that a client passing `since` with no path isn't really asking something coherent, and the code already has a documented "snapshot ignores format" deviation, this is likely a deliberate design and NOT a bug. 

I think I'm overanalyzing this one. Let me look for something cleaner. Let me reconsider.

Let me look more carefully at `ListGroupDevices` details branch — the offset/limit params. When `details=true`, the doc says "details stays unpaginated" — it ignores the limit and from_token. That's documented. OK.

Let me reconsider a different angle. Let me look for the strongest machine-checkable finding: a missing test. Let me look at what the group token functions (`groupTokenFor`/`parseGroupToken`) have — is there a test? And the bigint reshaping (`applyBigIntegerOpts`)? Let me check the parity tests.

$ cd /root/astrate-mule && rg -n "groupTokenFor|parseGroupToken|applyBigIntegerOpts|objectColumnTypes|nestTree|lastSegment|individualValue" internal/appengine/*_test.go | head -40
internal/appengine/groups_parity_test.go:242:			} else if _, ok := parseGroupToken(got); !ok {
internal/appengine/groups_parity_test.go:271:	rec = r.req(t, http.MethodGet, "/groups/"+pg+"/devices?from_token="+groupTokenFor(0), "", r.token)
internal/appengine/groups_parity_test.go:273:		t.Errorf("from_token=groupTokenFor(0): got %d (%s), want 200", rec.Code, rec.Body)

So there's minimal direct unit coverage of `parseGroupToken`/`groupTokenFor` round-trip boundaries (only exercised via parity tests with the runtime store, which are integration). And `applyBigIntegerOpts`, `objectColumnTypes`, `nestTree`, `lastSegment`, `individualValue` have no direct unit tests I can find. But that's a weak proposal (these are heavily exercised via http_test/parity tests indirectly).

Let me refocus. Let me find the genuinely strongest finding. Let me re-examine `applyBigIntegerOpts` at line 351 vs line 230. There's an asymmetry I want to verify: line 230 gates on either flag being non-nil; but `allow_safe_bigintegers=true` with `allow_bigintegers=nil` — the `ApplyBigIntegerOpts` only rewrites when the safe flag makes value exceed the bound, and leaves the stored (json.Number) as-is otherwise. That's consistent with "neither flag set leaves the stored representation untouched" — but here one flag IS set (safe). The doc comment says that ambiguity but the behavior seems intentional.

I think I'm spending too long. The strongest, most defensible, machine-checkable finding so far, matching the codebase's own stated anti-silent-lie principle, is the **snapshot path silently ignoring the documented since/to/limit query params**. But as analyzed it's likely deliberate.

Let me reconsider — there's actually a cleaner issue I should verify: **`parseQueryOpts` accepts `downsample_to` on a properties interface at the HTTP layer**, but the refusal happens in `readResolved`/`GetData`. That's fine.

Let me look again at a concrete correctness issue: in `datastreamData`, the object-aggregated interface with `downsample_to>0` returns a 422 (line 158). But the `parseQueryOpts`/doc doesn't distinguish. Fine.

Let me instead verify one thing that could be a real bug: **`renderObject` `disjoint_tables` for an empty subset**. When `docs[0]` has keys and `len(docs)==1`, that's fine.

Let me look at http_test.go briefly to understand the overall coverage and see whether there's a dead/uncovered branch. Actually, let me reconsider the strongest find given my constraints and settle.

Given the recipe's ranking (behaviour-wrong #1, missing-test #2), and the difficulty of confirming upstream on the snapshot (needs Legion), let me identify the clearest missing-test finding that is fully machine-checkable and doesn't depend on upstream semantics.

Looking at `parseGroupToken` comment: "The zero fields are not demanded — upstream accepts any well-formed v1 uuid." And `groupTokenFor` writes offset into time_low (uint32). There's a subtle correctness point: **`parseGroupToken` returns `int`**, and on the offset from `groupTokenFor(offset)` the round-trip works for offsets < 2^32. But `ListGroupDevices` gets `offset` from the token, then calls `ListGroupDevicesPage(ctx, g.ID, offset, limit+1)` and produces `groupTokenFor(offset + limit)`. The `offset+limit` could overflow uint32 if offset near MaxUint32, but that can't happen realistically.

Let me stop exploring and settle on the best findings. Given the code is exceptionally well-written with strong test coverage, "finding little is a correct result." But I should report at least the strongest candidate.

Let me make a final decision on the snapshot query-params finding. I'll phrase it carefully as a category-1 candidate that needs upstream confirmation, propose it as a task with the "verify against upstream then refuse-or-honour" framing, similar to how engine review framed its serverdata finding ("confirm with Giulio before committing"). Actually — the recipe says feature/design decisions go to for-giulio.md, not the queue. Refusing silently-dropped params is arguably a design decision.

Hmm. Let me reconsider whether proposing is right. The code comment at data.go:161-164 explicitly flags the downsample-root case as "the same silent lie the object case above exists to prevent" — meaning the author cares deeply about not silently dropping. Yet since/to/limit on the snapshot root are silently dropped with no comment at all. That inconsistency within the same function is the strongest signal it's an oversight. I'll propose it.

Let me also do a final scan of a couple test files to check the http_test coverage for the snapshot path to confirm it's not already asserted. Let me grep quickly.
$ cd /root/astrate-mule && rg -n "IndividualSnapshot|snapshot|/interfaces/"+'"'+"/interfaces/" 2>/dev/null; echo "==="; rg -n "interfaces/\\\{interface\\\}" internal/appengine/http_test.go | head
examples/pokemon-agent/emulator-agent/emulator_agent/astrate_client.py:127:        """Publish an object-aggregated GameState snapshot."""
bench/query.go:69:		// Interface-root snapshot: the device-overview pattern (nested tree).
bench/query.go:70:		"snapshot": func(d Device) string {
bench/README.md:13:bench query       AppEngine read mix: latest / 1h window / snapshot
examples/pokemon-agent/emulator-agent/emulator_agent/main.py:318:                # First tick: publish initial party snapshot
examples/pokemon-agent/docs/DESIGN.md:149:snapshot + publish on every *changed* state (position, battle flag, or dialog text
examples/pokemon-agent/docs/DESIGN.md:181:- **Aggregation:** Object — all fields published atomically as a single snapshot
examples/pokemon-agent/docs/DESIGN.md:185:  always sees a consistent snapshot, never a partial one.
examples/pokemon-agent/docs/DECISIONS.md:23:- Astrate delivers the entire object atomically; orchestrator never sees half a snapshot
docs/site/interface-schema.md:96:On load, each interface is compiled into a `CompiledInterface` with an endpoint trie for O(depth) matching. The compiled cache is held behind an `atomic.Pointer` snapshot (copy-on-write; readers never lock). Invalidation fires on Realm Management CRUD via Postgres `NOTIFY`.
docs/site/pokemon-agent.md:149:snapshot + publish on every *changed* state (position, battle flag, or dialog text
docs/site/pokemon-agent.md:181:- **Aggregation:** Object — all fields published atomically as a single snapshot
docs/site/pokemon-agent.md:185:  always sees a consistent snapshot, never a partial one.
docs/site/deployment.md:102:- **PostgreSQL** (`pgdata`): all metadata, properties, datastreams. Back up with `pg_dump`/`pg_basebackup` or volume snapshots. This is the source of truth.
docs/site/appengine-api.md:50:Returns the properties snapshot tree for a device's interfaces.
docs/site/operations.md:118:  with `pg_dump`/`pg_basebackup` or volume snapshots. This is the source of
docs/DESIGN.md:450:`map[realmID]map[name]map[major]*CompiledInterface` behind an `atomic.Pointer` snapshot (copy-on-
docs/DESIGN.md:453:snapshot. NOTIFY is not needed for correctness in a single process (a direct in-process callback
docs/OPERATIONS.md:119:  with `pg_dump`/`pg_basebackup` or volume snapshots. This is the source of
docs/handoff/flow-v2-decisions-2026-07-29.md:49:| Config drift | Store the flow’s **config snapshot** with the row; do not silently rewrite it | Named multi-instance needs stable config per flow |
docs/handoff/flow-lifecycle-followups-design-2026-08-22.md:71:  pipeline + config snapshot are stored; the resolved graph lives in memory).
docs/handoff/flow-lifecycle-followups-design-2026-08-22.md:121:  2. Persist the new config snapshot on the durable row.
test/conformance/cpc/cpc_test.go:11://	appengine devices data-snapshot / get-samples return the value →
test/conformance/cpc/cpc_test.go:333:	// 5. AppEngine: astartectl must read the published value back. data-snapshot
test/conformance/cpc/cpc_test.go:334:	// exercises the nested interface-root snapshot shape; get-samples the
test/conformance/cpc/cpc_test.go:339:			"appengine", "devices", "data-snapshot", id.String(), valuesName,
test/conformance/cpc/cpc_test.go:343:			t.Fatalf("data-snapshot did not return the published sample:\n%s", out)
docs/ROADMAP.md:320:| 6.1 | `internal/engine/cache.go` | Compiled-interface snapshot `map[realmID]map[name]map[major]*CompiledInterface` behind `atomic.Pointer`, copy-on-write rebuild from `store.LoadRealm`; invalidation via in-process callback **and** `store.Listen("astrate_interfaces")` (§2.6); device-introspection cache (per-connected-device, evicted on disconnect) | 200 |
docs/ROADMAP.md:345:- **T3 (broker+engine+DB):** test device connects (M4 certs) → introspection → publishes individual BSON, individual JSON, object-aggregated BSON → rows correct; property set/unset; `emptyCache` ⇒ device receives server-owned properties + purge message (decode/inflate asserted); `producer/properties` purge removes the right rows; LISTEN/NOTIFY: interface installed mid-run becomes publishable without restart; cache snapshot swap is race-free under `-race`.
docs/ROADMAP.md:372:| 7.7 | `internal/appengine/data.go` | GET interface data: properties snapshot tree; datastream queries (`since/since_after/to/limit/downsample_to` → store Series/Downsample; value re-encoding per stored `value_type` — longinteger as string, binaryblob base64, datetime RFC 3339, §2.3); PUT/POST server-owned (→ `engine.ServerData`); DELETE property unset | 260 |
internal/flow/factory_test.go:20:	snapshot := func() []string {
internal/flow/factory_test.go:86:		final = snapshot()
internal/engine/router_test.go:280:// TestInvalidation: the compiled-interface snapshot follows both
internal/engine/router_test.go:295:	waitFor(t, 5*time.Second, "notify-driven snapshot rebuild", func() bool {
internal/flow/blocks/httpblocks_test.go:286:func (c *captureState) snapshot() []capturedRequest {
internal/flow/blocks/httpblocks_test.go:371:	got := state.snapshot()
internal/flow/blocks/httpblocks_test.go:413:	got := state.snapshot()
internal/flow/blocks/httpblocks_test.go:434:	got := state.snapshot()
internal/flow/blocks/httpblocks_test.go:465:	if got := state.snapshot(); len(got) != 0 {
internal/engine/data.go:21:	// reasonRealmUnknown: the message realm is not in the schema snapshot.
internal/engine/data.go:130:	// Previous-value snapshot for the change-derived trigger conditions,
internal/engine/data.go:134:	// marks that a snapshot was taken at all (untracked ops skip change
internal/engine/data.go:156:		e.reject(m, reasonRealmUnknown, "realm not in schema snapshot")
internal/engine/data.go:354:// trackPrevious captures the pre-write previous-value snapshot an op needs
internal/engine/data.go:386:		e.log.Error("previous-value snapshot skipped: value does not render",
internal/engine/data.go:446:// capturePrevious fills the previous-value snapshot of an op that persists
internal/flow/flow.go:298:// ListFlows returns a snapshot of all registered flows.
internal/engine/engine.go:86:// Start loads the schema snapshot, subscribes to interface-change
internal/engine/engine.go:130:// RefreshTriggers rebuilds a realm's compiled snapshot — triggers ride in
internal/engine/engine.go:131:// the same snapshot as interfaces, so this is the in-process invalidation
internal/engine/engine.go:221:// strictly post-persist against the accept-time snapshot:
internal/engine/router.go:195:// start loads the schema snapshot, subscribes to interface-change
internal/engine/serverdata.go:20:	// ErrRealmUnknown: the realm is not in the schema snapshot.
internal/engine/serverdata.go:249:	// Previous-value snapshot taken before the row is gone, so the
internal/engine/fakestore_test.go:671:// schema snapshot preloaded — data_test and batch_test drive handle/flush
internal/engine/fakestore_test.go:787:// with the schema snapshot preloaded. Tests drive handle synchronously.
internal/engine/batch.go:36:	// pendingPrev chains previous-value snapshots inside one unflushed batch
internal/engine/cache_test.go:13:// TestSchemaCacheLoadAndLookup covers the snapshot build and the hot-path
internal/engine/cache_test.go:85:// invalidations naming a realm the snapshot has never seen.
internal/engine/cache_test.go:163:// TestSchemaCacheSnapshotRace exercises concurrent snapshot reads against
internal/engine/cache_test.go:269:// snapshot still loads.
docs/handoff/flow-design-a-named-durable-flows-2026-07-29.md:33:Flow      — named running (or desired) instance of a pipeline + config snapshot.
docs/handoff/flow-design-a-named-durable-flows-2026-07-29.md:41:- **Config snapshot:** the `config` JSON stored on the flow row is what start uses.
internal/engine/cache.go:37:	// schema-snapshot rebuilds.
internal/engine/cache.go:72:	// snapshot can attach compiled policies to triggers at load time.
internal/engine/cache.go:76:// fullReloadDebounce rate-limits snapshot self-heal reloads triggered by
internal/engine/cache.go:81:// realmSchema is one realm's compiled-interface set inside a snapshot. It is
internal/engine/cache.go:131:// shard. Readers never lock: they load the current snapshot pointer and use
internal/engine/cache.go:138:// schemaCache owns the snapshot pointer and its copy-on-write rebuilds.
internal/engine/cache.go:152:// newSchemaCache returns a cache holding an empty snapshot; call loadAll
internal/engine/cache.go:163:// realm resolves a realm by name from the current snapshot.
internal/engine/cache.go:189:// loadAll rebuilds the whole snapshot from the store (startup, unknown-realm
internal/engine/cache.go:216:// cloned, one entry is replaced, and the new snapshot is swapped in. A realm
internal/engine/cache.go:217:// ID missing from the snapshot falls back to a full reload (it may be a
internal/engine/triggers/match.go:16:// semantics; the engine captures the previous-value snapshot pre-write and
internal/engine/triggers/match.go:36:	// post-commit against the accept-time previous-value snapshot).
internal/engine/triggers/match.go:105:// they need the engine's pre-write snapshot (upstream get_value_change_
internal/engine/triggers/change_test.go:126:// value_stored never need a snapshot.
internal/engine/engine_e2e_test.go:214:// live schema snapshot).
internal/engine/engine_e2e_test.go:219:		t.Fatalf("realm %s not in schema snapshot", env.realm.Name)
internal/engine/triggers/events.go:82:// snapshot; a missing previous renders as null.
internal/testutil/astartedevice.go:140:// Messages returns a snapshot of every server message captured so far.
internal/broker/broker_test.go:83:func (r *recorderSink) snapshot() []LifecycleEvent {
internal/broker/broker_test.go:354:		events := sink.snapshot()
internal/broker/authhook.go:57:// snapshot (docs/DESIGN.md §3.1). Reload rebuilds both from the store; M7's
internal/broker/authhook.go:107:// rebuildTLSLocked recomputes the handshake snapshot. Callers hold p.mu (or
internal/store/store.go:164:// Stat returns a snapshot of the connection pool statistics
internal/store/notify.go:14:// realm's compiled-interface snapshot. The payload is the realm ID in
internal/broker/broker_e2e_test.go:266:	events := sink.snapshot()
internal/broker/broker_e2e_test.go:603:	for _, e := range sink.snapshot() {
internal/broker/broker_e2e_test.go:614:	events := sink.snapshot()
internal/flowapi/http.go:186:// updateFlowConfig replaces a flow's config snapshot; a live flow is rebuilt
internal/store/flows_test.go:162:	// UpdateFlowConfig replaces the snapshot and touches nothing else.
internal/store/datastreams.go:266:// IndividualSnapshot returns the most recent sample for each distinct path of
internal/store/datastreams.go:267:// an individual-datastream interface — the AppEngine interface-root snapshot
internal/store/datastreams.go:268:// ("data-snapshot") view upstream renders as a nested tree. Paths that never
internal/store/datastreams.go:270:func (s *Store) IndividualSnapshot(ctx context.Context, realmID int16, deviceID deviceid.ID, interfaceID int64) ([]IndividualRow, error) {
internal/store/datastreams.go:278:		return nil, fmt.Errorf("store: querying individual snapshot: %w", err)
internal/store/policies.go:13:// realm snapshot and the trigger executor honours it per delivery; the
internal/flowapi/service.go:634:// UpdateFlowConfig replaces a flow's config snapshot (issue #46). The new
internal/store/flows.go:174:// UpdateFlowConfig replaces the config snapshot of a durable flow. The
internal/store/datastreams_test.go:121:	t.Run("IndividualSnapshot", func(t *testing.T) {
internal/store/datastreams_test.go:128:		// Two distinct paths, three samples each: the snapshot must return
internal/store/datastreams_test.go:152:		rows, err := s.IndividualSnapshot(ctx, realm.ID, device, si.ID)
internal/store/datastreams_test.go:154:			t.Fatalf("IndividualSnapshot: %v", err)
internal/store/datastreams_test.go:157:			t.Fatalf("snapshot rows = %d, want 2 (one per path)", len(rows))
internal/store/datastreams_test.go:164:			t.Errorf("/d snapshot = %+v, want value 2", latest["/d"])
internal/store/datastreams_test.go:167:			t.Errorf("/s snapshot = %+v, want value v2", latest["/s"])
internal/store/datastreams_test.go:170:			t.Errorf("/d snapshot ts = %v, want newest", latest["/d"].TS)
internal/realm/service.go:56:// realm mutation the service refreshes the engine's compiled snapshot so the
internal/realm/service.go:505:// The engine compiles policies into its realm snapshot and attaches them to
internal/realm/service.go:509:// policy enters or leaves a snapshot — already refreshes it. Relax either
internal/realm/service.go:661:// interfacesChanged emits the store NOTIFY and refreshes the engine snapshot
internal/realm/service.go:675:// triggersChanged refreshes the engine snapshot after a trigger mutation
internal/realm/service.go:676:// (triggers ride in the same realm snapshot as interfaces).
internal/observability/metrics.go:66:// (a snapshot of store.Stat()).
internal/appengine/data.go:108:// query options) or a properties snapshot.
internal/appengine/data.go:116:		// parameter here would answer 200 with the full snapshot a client asked
internal/appengine/data.go:130:// branches, except where noted: the snapshot view and the downsampled view
internal/appengine/data.go:134:	// "data-snapshot" view — the latest sample for every endpoint, rendered as
internal/appengine/data.go:139:		rows, err := s.st.IndividualSnapshot(ctx, r.rid, r.id, r.iface.ID)
internal/appengine/data.go:160:		// An interface-root query is the snapshot view — the latest sample of
internal/appengine/data.go:416:// value}} tree of every set property (the upstream interface snapshot shape,
internal/appengine/http_test.go:320:		// The interface root is a snapshot of every endpoint's latest sample, not
internal/appengine/http_test.go:358:		// Interface-root query is the nested upstream snapshot tree, not a flat
internal/appengine/http_test.go:379:			t.Fatalf("snapshot missing /value endpoint: %v", snap)
internal/appengine/http_test.go:382:			t.Errorf("snapshot /value = %v (%T), want newest 3.0", v.Value, v.Value)
internal/appengine/http_test.go:385:			t.Errorf("snapshot /value timestamp = %v, want newest %v", v.Timestamp, r.t3)
internal/appengine/http_test.go:389:			t.Fatalf("snapshot missing /big endpoint: %v", snap)
internal/appengine/http_test.go:392:			t.Errorf("snapshot /big = %v (%T), want longinteger string", big.Value, big.Value)
internal/appengine/service.go:715:		// See GetData: reducing a snapshot is refused, not silently ignored.
===
rg: regex parse error:
    (?:interfaces/\\{interface\\})
                     ^
error: repetition quantifier expects a valid decimal

I have enough. The snapshot behavior is well-tested at the individual-datastream level (http_test.go:320-392), and the code documents the snapshot as a deliberate design. The since/to/limit-on-snapshot question is a design subtlety I can't confirm against upstream without the Legion Go, and the code's extensive documentation suggests it's deliberate. Proposing it as a firm bug would be a false positive risk; I'll note it as a "decided not to propose" with reasoning, and flag the upstream-verification question in for-giulio.md as a potential silent-ignore worth confirming — actually, the recipe says feature/design items go to for-giulio.md. This is borderline.

Let me reconsider once more for the strongest clean finding. Let me re-examine `applyBigIntegerOpts` more carefully, as it has a real asymmetry worth checking.

Lines 386-412:
```go
switch {
case opts.AllowBigIntegers != nil && !*opts.AllowBigIntegers:
    doc[name] = n.String()
case opts.AllowSafeBigIntegers != nil && *opts.AllowSafeBigIntegers:
    if i > maxSafeJSONInteger || i < -maxSafeJSONInteger {
        doc[name] = n.String()
    }
}
```

And it's only invoked when `opts.AllowBigIntegers != nil || opts.AllowSafeBigIntegers != nil` (line 230). So if a client sends `allow_safe_bigintegers=false`, nothing is rewritten — values stay as stored (json.Number). Is "allow_safe_bigintegers=false" supposed to degrade all longintegers to strings? Per upstream, `allow_safe_bigintegers` semantics: "true keeps safe values as numbers". With it false, everything should become a string. But here `*opts.AllowSafeBigIntegers == false` → the case doesn't match → no rewrite → stays as json.Number. That could be a bug but requires upstream semantics knowledge. Too uncertain.

Let me settle. Given the extremely high quality of this code and test coverage, "finding little is a correct result." I'll do a final check of `downsample.go` (the one file I haven't read) and then write the review with modest, defensible findings.
→ Read internal/appengine/downsample.go
Now let me consolidate. I've read all four non-test source files of `internal/appengine`. Let me decide on the findings.

My strongest defensible candidate: **the interface-root snapshot silently ignores the query params the endpoint documents.** But I've concluded this is likely deliberate (the "data-snapshot" deviates from format, and it's the latest-value-per-endpoint so time-window/limit don't apply). It's a design-y question. Per the recipe, a design question belongs in for-giulio.md as one line, not the queue. Actually — the recipe says "A feature is a design decision and it is not yours: write it to .mule/for-giulio.md." But a *silent* ignore of documented parameters is arguably a correctness bug, not a feature. The codebase's own principle (refuse rather than silently ignore) is violated with no comment on this branch, unlike the downsample branches which explicitly refuse.

I'll propose this as a task (category-1, with the anti-silent-lie precedent in the same function), framed to first verify against upstream. This matches how the engine review handled the serverdata finding.

Let me also add the missing-test finding for `parseGroupToken`/`groupTokenFor` round-trip boundaries — fully machine-checkable, no upstream dependence. This is the #2 category. The functions only exist as pure functions but only get exercised through the pg-backed parity tests (integration). A pure unit round-trip test (including the v1-nibble rejection / non-v1 acceptance) is cheap and machine-checkable. That's a solid one.

Let me settle on these two task lines. Both name files and lines. Let me reconsider the snapshot one's strength — I want to make sure it's not a false positive. The recipe warns against pattern-matching rather than reading. I read the code directly. The finding is: within `datastreamData`, three branches explicitly refuse incompatible params, but the first branch (snapshot) drops since/since_after/to/limit/sort silently with no comment. That asymmetry is a real, observed inconsistency, not pattern-matching. I'll propose it but clearly tag it as needing upstream verification first.

Now let me write the review file and append the task lines. Also I should check whether to add a for-giulio.md note. The snapshot question is tied to the task proposal, so I'll keep it in the task rather than for-giulio.

Let me write `.mule/reviews/appengine-2026-09-07.md`.

← Write .mule/reviews/appengine-2026-09-07.md
Wrote file successfully.

Now let me append the task lines to `.mule/todo.md`. I need to add them tagged with `[auto]` per the recipe (tasks invented are tagged `[auto]`). Two tasks, each naming its file.

Let me append at the end of the queue.
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -114,4 +114,6 @@
 - [x] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups/{group}/devices is missing `from_token` and `limit` pagination params that code reads at http.go:439-441 (and returns via links.next), (2) POST /groups is missing 422 validation-error response (FieldErrors for blank group_name/empty devices, http.go:404-416), (3) PATCH /devices/{device} and PATCH /devices-by-alias/{alias} are missing 422 ("Attribute key not found" via ErrAttributeKeyNotFound, service.go:352-357 → http.go:656-657), and the device-scoped PATCH is also missing 409 ("Alias already in use", http.go:651) and its documented 500 covers the merge-patch content-type mismatch which code answers with an unmapped 500 too, (4) GET /devices/{device}/interfaces/{interface}/{path} is missing 422 (invalid since/since_after/to/limit/downsample_to/format params, http.go:554-639). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.
 - [x] docs-sync-appengine-group-patch-status [auto]: PATCH /groups/{group}/devices/{device} (internal/appengine/http.go:518-537, patchGroupDevice) shares applyPatch with the device-scoped PATCH and also answers 422 ("Attribute key not found", ErrAttributeKeyNotFound, service.go:352-357) and 409 ("Alias already in use", http.go:651), but docs/api/astarte_appengine_api.yaml lists only 200/400/401/404/500 for it — add the two missing responses.
 - [x] docs-sync-appengine-data-422-interface-level [auto]: the interface-level data GETs — GET /devices/{device}/interfaces/{interface}, GET /devices-by-alias/{alias}/interfaces/{interface}, GET /groups/{group}/devices/{device}/interfaces/{interface} — share parseQueryOpts (internal/appengine/http.go:554-639) and answer 422 FieldErrors for invalid since/since_after/to/limit/downsample_to/format params, but the spec documents only the path-level GET with 422 — add 422 to the interface-level GETs.
+- [ ] appengine-snapshot-ignores-query-params [auto]: in `internal/appengine/data.go:138-148`, the individual-datastream interface-root snapshot branch ignores `Since`/`SinceAfter`/`To`/`Limit`/`Descending` even though the root GET is documented with all of them (docs/api/astarte_appengine_api.yaml:577-582) — the same function refuses incompatible params elsewhere (downsample 422, data.go:118-120/165-167), so this is a silent drop. Probe upstream on the Legion Go to decide honour-vs-refuse, then either apply the window/limit or answer an explicit 422 and stop advertising the params on the root GET; add a T1 test in data_test.go/query_opts_test.go pinning the chosen behaviour.
+- [ ] appengine-group-token-roundtrip-test [auto]: add a container-free unit test for `parseGroupToken`/`groupTokenFor` (internal/appengine/service.go:602-623) — currently only exercised via the pg-backed groups_parity_test.go. Assert (1) round-trip `groupTokenFor(offset)` → `parseGroupToken` recovers offset, (2) a non-v1-nibble UUID is rejected, (3) a well-formed v1 UUID with zeroed time fields is accepted (parseGroupToken comment, service.go:614-616).

I've completed the code-review recipe. Let me summarize.

Note: I did not touch git, and this was a proposal job — no `gh issue create` was involved (the recipe's deliverable is task lines + a review file, which I wrote). I only appended task lines to `.mule/todo.md` and wrote the evidence file.

```
Done: Reviewed internal/appengine (read service.go, http.go, data.go, downsample.go + tests); proposed 2 tasks
Files: .mule/reviews/appengine-2026-09-07.md, .mule/todo.md
Verified: none needed (proposal job; no code changed, no build)
Unsure: whether the snapshot path should honour vs refuse the documented query params — needs upstream probe on the Legion; and applyBigIntegerOpts's safe=false semantics (left for a probe, not proposed)
Follow-ups: appended 2 tasks to .mule/todo.md:
  - appengine-snapshot-ignores-query-params [auto]
  - appengine-group-token-roundtrip-test [auto]
```

Picked `internal/appengine` (largest package with no prior review file), read the four non-test sources largest-first. Yield: (1) a category-1 finding — the interface-root snapshot silently drops the endpoint's documented `since/to/limit` params where the same function refuses incompatible params elsewhere, and (2) a category-2 missing-test on the group-token cursor round-trip. The rest was well-covered; I documented what I decided not to propose and why so the next run won't repeat it.
