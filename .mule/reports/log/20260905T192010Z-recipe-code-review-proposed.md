slug: recipe-code-review
verdict: proposed
at:  f407d45
ran: 2026-09-05T19:20:10Z on DietPi in 120s


> build · big-pickle

$ cat .mule/MULE.md 2>/dev/null || echo "NO MULE.md"
# Mule rules

You are the mule for Astrate: a Go implementation of the Astarte device-management platform,
wire-compatible with upstream Astarte. You work alone, one task at a time, on a branch. A
human reads your commits later.

These rules are one instance of the "unattended-agent safety envelope" pattern catalogued in
the system-wide `scaffold` skill — read that if setting up equivalent safety rules
for a different project. This project auto-commits on a green gate (`go vet`/`go test`/lint)
specifically because Go has a mechanical truth-check; a project without one (e.g. g-mind's
vault prose) should never let its wrapper commit unattended — see `scaffold` for why
that's the same rule, not a contradiction.

Everything you write — code, comments, task lines, notes — is in **English**.

## The one rule that matters

**Your useful context is about 30k tokens, not 200k.** After that your output gets worse
without you noticing. So:

- Read only the files the task names, plus what you must to be correct. Not the whole package.
- Use `rg` with narrow patterns instead of reading files to find things.
- When you have a working solution, stop. Do not go looking for more.
- If the task turns out to be bigger than it looked, **do the first independently useful part
  and add the rest to `.mule/todo.md` as new task lines.** Splitting is success, not failure.

You get a fresh session for every task. Nothing carries over except what you write to disk.

## Standing instructions from Giulio

These outrank any recipe. They are here because a recipe would otherwise keep re-deriving
them and putting the same proposal back in front of him.

- **Issue #1 is never to be raised.** It stays open permanently. Do not close it, do not
  triage it, do not propose it, do not write it into `for-giulio.md`. Any recipe that
  surfaces "stale issues" skips it silently.
- **Dependency bumps are decided at milestone boundaries, not opportunistically.** Re-run
  the `go list -m -u` sweep when a milestone tag or `APICompatVersion` moves, and propose a
  bump only for a fix Astrate actually needs. `go.mod` is on the never-touch list either way.

## Where you are running

Usually a Raspberry Pi (4 cores, 3.7GB, no Docker), on a timer, with nobody watching. Two
consequences you must internalise:

- **The race detector does not work there** (39-bit VMA kernel). Your gate is
  `go vet ./... && go test ./...`. So a concurrency bug will pass your gate. If a task asks
  you to touch goroutines, channels or locks, do it *and say loudly in your report that the
  change is unverified for races* — that report is the only warning anyone gets.
- **No Docker, so no database.** Anything needing a live DB, a broker, or a full stack is a
  `[legion]` task, not yours.

The **Legion Go** (`ssh legion` — 16 cores, 12GB, Docker with the whole Astarte stack up, Go
1.26.5 at `~/.local/go`) is the muscle. A task line tagged `[legion]` needs it and is skipped
automatically when it is asleep. It is the **only** machine here that can run `go test -race`
(~40s), so that is where race coverage comes from. Details in `.mule/recipes/legion-go.md` —
read that before any `[legion]` task.

A task line tagged `[readonly]` is a standing check: run it, report what you found, and
**change nothing**. Finding nothing wrong is the good outcome — do not invent a change to
have something to show. Your report is the whole deliverable and is kept verbatim.

## What you may and may not do

- Change code, tests, and files under `bench/` and `tools/` as the task requires.
- **Never touch git.** No commit, branch, checkout, stash, rebase. The script commits your
  work if the gates pass, and reverts it if they don't.
- **Never touch** `docs/DESIGN.md`, `docs/ROADMAP.md`, `migrations/`, `.github/`, `go.mod`,
  `go.sum`, `Dockerfile`, `docker-compose.yml`, `.mule/milestones.md`, or anything under
  `.trickle/`. These carry decisions or Giulio's voice. If a task seems to need one, say so
  and stop.
- Never invent an API, a config key, or a behaviour. Read the code or run it.
- Never weaken or delete a test to make a gate pass. A failing gate is a real result — report it.
- Never write to a path outside this repository.

## How to do a task

1. Read the task. Say in one line what you think it means. If it is ambiguous in a way that
   changes the code, **stop and write the ambiguity into your report** instead of guessing.
2. Read the named files first. Grep for callers of anything you change.
3. Make the change. Match the surrounding style — this codebase is idiomatic Go with short
   receiver names, wrapped errors (`fmt.Errorf("...: %w", err)`), and table-driven tests.
4. **If you changed behaviour, a test must prove it.** Write one that fails without your
   change. If you cannot write such a test, say so explicitly in the report.
   This is checked, not trusted: the runner takes your implementation back out, keeps your
   tests, and requires them to break. A test that passes without your change is rejected and
   the whole task is thrown away. So do not write a test that asserts something already true
   in order to have something to show — it will not work, and you will have wasted the run.
5. Run the gate yourself before finishing: `go test -race ./...` and `gofmt -l .`
   (`go test ./...` without `-race` if the race build is slow, but say which you ran).
6. Report.

## Verifying against things you do not control

For anything involving upstream Astarte, a database, or an external dependency: **write a
twenty-line throwaway program that answers one question, and run it.** Reading the code and
guessing has repeatedly been wrong here — a live probe is cheaper than a wrong patch. Put
probes in `/tmp`, not in the repo. Report the measurement as a measurement.

## Report format

End every task with exactly this, and nothing longer:

```
Done: <one line>
Files: <paths>
Verified: <commands run> -> <pass/fail>
Unsure: <what you guessed, or "nothing">
Follow-ups: <task lines you appended to .mule/todo.md, or "none">
```

A stated doubt is worth more than a silent patch. Your diff gets read.

## Tasks that came from an issue

A task slugged `issue-<N>` came from a GitHub issue. Read it first — `gh issue view <N>` —
because the line in the queue is only its title, and the body is where the actual requirement
is. If the issue asks for more than one tick's worth, do the first independently useful part
and propose the rest as new lines, exactly as you would for any oversized task. Do not close
the issue; the script comments the sha on it for you.

## Tasks that file an issue instead of changing code

A task slugged `milestone-<tag>-issue-...` (from `.mule/recipes/milestones.md`) is done when
you have run the `gh issue create` command it names — that command *is* the deliverable,
there is no code diff to produce. Run exactly the command proposed (title, labels, body);
don't improvise the wording. This leaves the working tree clean, which is expected: it is
logged the same way a `[readonly]` probe is, not committed as a code change.

## Adding tasks

Append to `.mule/todo.md`, one line each, never more than ~8 at a time:

```
- [ ] <slug>: <the outcome, in one line — what should be true when it's done>
```

A good task line names the outcome and the file. "fix logging" is a bad line;
"trigger-log-level: demote the per-message trigger dispatch log in
internal/triggers/dispatch.go from Info to Debug" is a good one. If a task needs more than
one line to state, write `.mule/tasks/<slug>.md` with the detail and keep the line short.

$ cat .mule/recipes/code-review.md 2>/dev/null || echo "NO RECIPE"
# Recipe — codebase review

Read one area of the codebase and propose work. **This run writes no code.**

## Pick one area, not the codebase

Your context does not stretch to a whole-repo review, and a review that skims everything
finds nothing. Pick **one** package and go deep. Rotate through them across runs — check
`.mule/log.md` and `.mule/reviews/` for which areas were done recently, and pick one that
was not.

```sh
ls internal/ pkg/                                   # the map
rg -c '' --glob 'internal/<pkg>/**/*.go' | sort -t: -k2 -rn | head   # the size of it
```

Then read that package's files, largest first, and its tests.

## What is worth proposing

In descending order of value:

1. **A behaviour that is wrong or unguarded** — an unchecked error, a nil path, an
   unbounded slice grown from network input, a goroutine with no exit, a lock held across
   I/O, a context ignored. Propose with the file and line.
2. **A missing test on a rule that already exists.** Rules nothing asserts on rot silently.
   Best kind of task there is: fully machine-checkable.
3. **Performance, but only where you can name the workload.** "This allocates per message
   in the ingest hot path" is a proposal. "This could be faster" is not. Propose a
   `go test -bench` measurement task *before* an optimisation task, never the reverse.
4. **Clarity** — a function that does three things, a name that lies, duplicated logic in
   two packages. Propose only when the change is mechanical and the tests already cover the
   behaviour.
5. **New features.** Rarely. A feature is a design decision and it is not yours: write it to
   `.mule/for-giulio.md` as one line saying what and why, not to the queue.

## What not to propose

- Renames and reformatting for their own sake. The linter owns that.
- Anything in `docs/DESIGN.md` or `docs/ROADMAP.md` — that is the frozen spec, and if the
  code contradicts it, **the code is the bug**; say so in `.mule/for-giulio.md`.
- Refactors that touch more than two files. They cannot be verified in one task.
- Anything you found by pattern-matching rather than by reading. Say "I read
  `internal/x/y.go` lines 40-120" or do not propose it.

## Output

Write `.mule/reviews/<pkg>-<date>.md`: what you read, what you found, and what you decided
*not* to propose and why (that half is what stops the next run repeating your work). Then
append at most five task lines to `.mule/todo.md`, each naming its file.

Finding little in a well-written package is a correct result. Say so and stop.

$ ls internal/ pkg/ 2>/dev/null
internal/:
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

pkg/:
astarteapi
deviceid
interfaceschema
payload

$ cat .mule/log.md 2>/dev/null | tail -50; echo "=== reviews ==="; ls -la .mule/reviews/ 2>/dev/null
# Mule log

One row per task attempt, written by `tools/mule.sh`. This is the record of what the cheap
layer can actually do — read it before deciding whether a kind of task is worth delegating.

`secs` is the honest signal: a task that used most of its 900s budget was too big.

| date | task | outcome | secs | note |
| --- | --- | --- | --- | --- |
| 2026-09-04 | purge-properties-compression-capabilityauto | done | 86s | 6f7a3d6 |
| 2026-09-04 | broker-acl-coldstart-introspection-miss | done | 629s | d201db4 |
| 2026-09-04 | issue-93 | done | 108s | 24ad5b8 |
| 2026-09-04 | broker-disconnect-device-zombie-session | blocked | 117s | wrote nothing |
| 2026-09-04 | broker-offline-acl-tests | blocked | 208s | wrote nothing |
| 2026-09-04 | broker-onconnect-doc-comment | blocked | 438s | tests failed: --- FAIL: TestMQTTSink_Retained (0.02s) |
| 2026-09-05 | empty-introspection-verification | blocked | 106s | wrote nothing |
| 2026-09-05 | probe-trigger-install-notification-delay | blocked | 308s | wrote nothing |
| 2026-09-05 | compat-note-v132 | done | 73s | 4bc3e1a |
| 2026-09-05 | docs-sync-pairing-health-path | done | 186s | d73e225 |
| 2026-09-05 | docs-sync-pairing-register-404 | done | 106s | 44a7cae |
| 2026-09-05 | swagger-httptest-coverage | done | 267s | ba79b34 |
| 2026-09-05 | probe-property-resend-encoding | blocked | 100s | wrote nothing |
| 2026-09-05 | compat-note-v133 | blocked | 44s | wrote nothing |
=== reviews ===
total 16
drwxr-xr-x 2 root root 4096 Sep  4 20:46 .
drwxr-xr-x 7 root root 4096 Sep  5 21:08 ..
-rw-r--r-- 1 root root 8188 Sep  4 20:43 broker-2026-09-04.md

Only `broker` was reviewed (2026-09-04). Let me check the sizes of the packages to pick one, preferably the largest or one not recently reviewed. `broker` was already reviewed. Let me find the size of each package to pick a large one.
$ for d in internal/* pkg/*; do n=$(rg -l '' "$d" --glob '*.go' 2>/dev/null | wc -l); echo "$n $d"; done | sort -rn
49 internal/flow
48 internal/engine
32 internal/store
23 internal/appengine
13 internal/broker
12 pkg/interfaceschema
10 pkg/payload
9 internal/realm
8 internal/pairing
8 internal/flowapi
8 internal/auth
6 internal/testutil
5 internal/observability
4 internal/httpx
4 internal/housekeeping
2 pkg/deviceid
2 pkg/astarteapi
2 internal/swagger
2 internal/config

`internal/flow` is the largest (49 files) and was never reviewed. Let me look at it.
$ rg -c '' --glob '*.go' internal/flow/*.go | sort -t: -k2 -rn | head -30
internal/flow/router_test.go:565
internal/flow/composite_test.go:391
internal/flow/message.go:364
internal/flow/message_test.go:348
internal/flow/flow.go:332
internal/flow/manager_test.go:318
internal/flow/router.go:284
internal/flow/pipeline_test.go:280
internal/flow/factory_test.go:274
internal/flow/factory.go:249
internal/flow/flow_test.go:198
internal/flow/pump_test.go:194
internal/flow/pipeline.go:165
internal/flow/composite.go:156
internal/flow/substitute.go:126
internal/flow/block.go:113
internal/flow/substitute_test.go:85
internal/flow/graph.go:64

$ ls internal/flow/*.go
internal/flow/block.go
internal/flow/composite.go
internal/flow/composite_test.go
internal/flow/factory.go
internal/flow/factory_test.go
internal/flow/flow.go
internal/flow/flow_test.go
internal/flow/graph.go
internal/flow/manager_test.go
internal/flow/message.go
internal/flow/message_test.go
internal/flow/pipeline.go
internal/flow/pipeline_test.go
internal/flow/pump_test.go
internal/flow/router.go
internal/flow/router_test.go
internal/flow/substitute.go
internal/flow/substitute_test.go

Let me read the main source files (largest first: message.go, flow.go, router.go, factory.go, pipeline.go).
→ Read internal/flow/message.go
→ Read internal/flow/flow.go
Let me read router.go and factory.go next.
→ Read internal/flow/router.go
→ Read internal/flow/factory.go
→ Read internal/flow/pipeline.go
→ Read internal/flow/block.go
→ Read internal/flow/composite.go
→ Read internal/flow/graph.go
Let me read the remaining source files: pipeline.go, composite.go, block.go, substitute.go, graph.go.
→ Read internal/flow/substitute.go
I've read all the source files. Now let me look at the tests to understand what's covered, focusing on what rules exist but aren't tested. Let me grep the test files for coverage insights and check the specific areas.

Let me look at a couple of observations I have:

1. **`manager_test.go`** (largest test) - I want to check StopFlow / concurrency coverage.

2. **`composite.go`** - `expandComposites` has a subtle issue I'd like to verify.

Let me think about the substantive findings before checking tests.

Let me look at observations:

**Finding A — `ExpandComposites` renames inner block names with `subPrefix + "name"`, but this can collide with blocks in the outer level.** Actually wait, more importantly: In `expandComposites`, when a composite block `node` at the top level is replaced, `outNodes` gets `expNodes` (renamed with prefix). The outer `conns` referencing the node are renamed to first/last. But what about the case where the *renamed* inner node name could collide with another block in the outer level? If the outer already has a block literally named e.g. `innerName` and a composite `comp` containing block `innerName`, then `comp.innerName` vs `innerName` won't collide. Prefixing guards it. OK.

**Finding B — `Pipeline.Validate` has clearly duplicated code (lines 117-147).** It computes `hasSource`/`hasSink` at lines 117-126, then *re-computes* from fresh maps `inDeg2`/`outDeg2` at lines 128-147, discarding the first computation entirely. The first block of source/sink detection (lines 117-126) is dead code — its results are immediately overwritten. This is a clarity issue: a block of dead/commented-out logic. It's mechanical to remove. This matches recipe category 4 (clarity, mechanical, tests cover behaviour).

**Finding C — timing/`Drain` with QoS0 OverflowDrop under drain.** Not obviously a bug.

**Finding D — `flow.go` StopFlow: `setStatus` is called without the lock in `StartFlow`?** Let me check. In `StartFlow`, `f.mu.Lock()` then `f.setStatus(FlowStatusRunning)`. setStatus itself doesn't lock. OK that's under lock. Good.

**Finding E - `Router.Submit` count of `droppedQoS0`/`droppedQoS1` and `onDrop`:** onDrop is set on flowMsg but never set anywhere in this package (declared but the struct field only ever nil). Actually `fm := &flowMsg{msg: msg, qos: qos}` — onDrop is never assigned. So the `if fm.onDrop != nil` checks are dead code. Minor clarity point.

Let me focus and check the biggest real finding: **the duplicate source/sink detection dead code in pipeline.go**. Let me verify it's truly dead. Lines 117-126 set `hasSource` and `hasSink` using `outDeg`/`inDeg` which were modified by the topo sort (zeroed). Then lines 128-147 recompute into new maps and overwrite `hasSource`/`hasSink`. So lines 117-126 result is thrown away. This is genuinely redundant code that the tests presumably cover (Validate's hasSource/hasSink error cases). Removing it is a mechanical, safe clarity fix.

Let me confirm there's a test for the hasSource/hasSink validation paths, and check other spots. Let me look at the test files to check coverage of the rules that exist.

$ rg -n "has no source|has no sink|no source block|no sink block" internal/flow/*_test.go
(no output)

✱ Grep "onDrop" in internal/flow · 5 matches
So: the "has no source" and "has no sink" validation error paths in `Validate()` have **no test at all**. That's a missing test on a rule that exists (recipe category 2 — the best kind). Let me confirm by searching the whole package tests for Validate tests.

Also `onDrop` is never set anywhere in the package — dead branches. Let me check the pipeline test coverage.

$ rg -n "Validate|Pipeline" internal/flow/*_test.go | head -40
internal/flow/pipeline_test.go:8:func TestPipeline_Validate(t *testing.T) {
internal/flow/pipeline_test.go:11:		p       Pipeline
internal/flow/pipeline_test.go:16:			p: Pipeline{
internal/flow/pipeline_test.go:19:				Blocks: []PipelineNode{
internal/flow/pipeline_test.go:30:			p: Pipeline{
internal/flow/pipeline_test.go:33:				Blocks: []PipelineNode{
internal/flow/pipeline_test.go:49:			p: Pipeline{
internal/flow/pipeline_test.go:52:				Blocks: []PipelineNode{
internal/flow/pipeline_test.go:65:			p: Pipeline{
internal/flow/pipeline_test.go:66:				Blocks: []PipelineNode{{Name: "a", BlockType: "x"}},
internal/flow/pipeline_test.go:72:			p: Pipeline{
internal/flow/pipeline_test.go:79:			p: Pipeline{
internal/flow/pipeline_test.go:81:				Blocks: []PipelineNode{
internal/flow/pipeline_test.go:93:			p: Pipeline{
internal/flow/pipeline_test.go:95:				Blocks: []PipelineNode{
internal/flow/pipeline_test.go:106:			p: Pipeline{
internal/flow/pipeline_test.go:108:				Blocks: []PipelineNode{
internal/flow/pipeline_test.go:119:			p: Pipeline{
internal/flow/pipeline_test.go:121:				Blocks: []PipelineNode{
internal/flow/pipeline_test.go:136:			p: Pipeline{
internal/flow/pipeline_test.go:138:				Blocks: []PipelineNode{
internal/flow/pipeline_test.go:151:			p: Pipeline{
internal/flow/pipeline_test.go:153:				Blocks: []PipelineNode{
internal/flow/pipeline_test.go:168:			err := tt.p.Validate()
internal/flow/pipeline_test.go:170:				t.Errorf("Validate() err = %v, wantErr %v", err, tt.wantErr)
internal/flow/pipeline_test.go:176:func TestPipeline_JSONRoundTrip(t *testing.T) {
internal/flow/pipeline_test.go:177:	p := Pipeline{
internal/flow/pipeline_test.go:180:		Blocks: []PipelineNode{
internal/flow/pipeline_test.go:208:	var got Pipeline
internal/flow/pipeline_test.go:239:func TestPipeline_MarshalRejectsInvalidPipeline(t *testing.T) {
internal/flow/pipeline_test.go:240:	p := Pipeline{
internal/flow/pipeline_test.go:242:		Blocks: []PipelineNode{
internal/flow/pipeline_test.go:255:func TestPipeline_JSONContainsAllFields(t *testing.T) {
internal/flow/pipeline_test.go:256:	p := Pipeline{
internal/flow/pipeline_test.go:259:		Blocks: []PipelineNode{
internal/flow/flow_test.go:40:					PipelineID: "trans-1",
internal/flow/flow_test.go:57:					PipelineID: "trans-2",
internal/flow/flow_test.go:87:			PipelineID: "trans-3",
internal/flow/flow_test.go:103:		PipelineID: "ts-test",
internal/flow/flow_test.go:126:		PipelineID: "stop-test",

→ Read internal/flow/pipeline_test.go [offset=1, limit=175]
I see. The "no source" and "no sink" tests exist, but **both use a cycle** (`a→b, b→a`), so they hit the "cycle" error before reaching the source/sink check. They never actually reach the `hasSource`/`hasSink` validation branch. So those error paths are effectively untested — the tests pass via the cycle error, not via the source/sink error they claim to assert. That's a genuine missing-test finding (a rule that exists but nothing asserts on it directly).

Now let me verify my dead-code finding in pipeline.go. Lines 117-126 compute `hasSource`/`hasSink` from the topo-sorted (zeroed) `inDeg`/`outDeg`, then lines 128-147 recompute from fresh `inDeg2`/`outDeg2` and overwrite them. So lines 117-126 are dead. That's a clarity finding.

Let me verify the build/tests pass baseline, then write the review file and todo lines. Let me run the gate quickly to confirm the package compiles.

$ go test ./internal/flow/ 2>&1 | tail -5; echo "---gofmt---"; gofmt -l internal/flow/
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
---gofmt---

→ Read .mule/reviews/broker-2026-09-04.md [limit=30]
Now let me verify my claim that the "no source"/"no sink" test cases never reach the intended branch. Let me check which error these actually produce — both my candidate test cases are cycles, so they hit the cycle error first. Let me confirm the dead code claim is airtight by checking line numbers again.

The "no source" and "no sink" redundant test cases both use `a→b, b→a` which is a cycle → `Validate` returns the cycle error at line 114, before ever reaching `hasSource`/`hasSink` at line 148. So the `hasSource`/`hasSink` branches have no reachable test. Good finding.

Let me also confirm the dead-code in pipeline.go and the `onDrop` dead field. Let me create the review file and append tasks now.

Let me check the date and existing review format for consistency.

$ cat .mule/todo.md 2>/dev/null | tail -20
- [ ] race-check-store: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/store/... ./internal/housekeeping/... ./migrations/...`. Report any failure to .mule/for-giulio.md with the full race report. Split out of the former single `race-check` line, which timed out running the whole tree at once. [legion] [readonly]
- [ ] race-check-engine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/engine/... ./internal/broker/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
- [ ] race-check-flow: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/flow/... ./internal/realm/... ./internal/pairing/... ./internal/auth/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
- [ ] race-check-appengine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/appengine/... ./internal/observability/... ./internal/httpx/... ./internal/config/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
- [ ] race-check-pkg: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./pkg/... ./cmd/... ./internal/testutil/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]

- [x] broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [approved 2026-09-04]
- [!] broker-disconnect-device-zombie-session: in `internal/broker/broker.go` `DisconnectDevice` (lines 260-266), force-closing the live connection leaves the persisted bbolt session (subscriptions + offline queue, keyed by CN) in place — mochi derives `expire` from the Clean flag, not the reason code (server.go:485-491), so `dropClient` (sessionstore.go:500) is never reached and a later device with the same CN resurrects the deleted device's session. Wire the session store handle into `Broker` and call `dropClient` from `DisconnectDevice`, and add a T1 test in `broker_test.go`: connect clean=false + subscribe, call `DisconnectDevice`, reconnect — assert `session_present` is false and the stale subscription is gone. [approved 2026-09-04] — BLOCKED: wrote nothing
- [!] broker-offline-acl-tests: in `internal/broker/aclhook_test.go`, unit-test the offline-delivery ACL — `offlineACL.ownershipOf` cache hit within TTL, TTL expiry triggering a reload, and a load failure caching as empty and denying (aclhook.go:116-138) — plus a T1 `OnACLCheck` check that the offline branch (delivery to a session-present but disconnected device, aclhook.go:196-202) consults the store. The rule currently has zero direct tests; `TestCheckACLMatrix` exercises only the pure `checkACL`. [approved 2026-09-04] — BLOCKED: wrote nothing
- [!] broker-onconnect-doc-comment: in `internal/broker/authhook.go:314`, restore the missing first line of the `OnConnect` doc comment — it opens mid-sentence ("of the Astarte MQTT v1 protocol, and mochi publishes them without a publish-side ACL check…") and the Will-clearing security rationale (a retained LWT would escape the §3.2 matrix) is only half-documented. [approved 2026-09-04] — BLOCKED: tests failed: --- FAIL: TestMQTTSink_Retained (0.02s)

- [x] purge-properties-compression-capability [auto]: upstream v1.3.0 adds a `purge_properties_compression_format` device capability (`zlib`|`plaintext`, default `zlib`) — a wire-visible capability value. Check whether Astrate's capabilities handling (internal/broker, the `<realm>/<device_id>/capabilities` topic, issue #16) needs to recognize/honour it, or whether zlib-only is already the deliberate default; propose the change or note why not needed.
- [!] empty-introspection-verification [auto]: upstream v1.3.0 changed "allow devices with empty introspection" — verify whether Astrate's device connection/introspection handling currently rejects an empty introspection string where upstream now accepts it, and propose a fix if so. — BLOCKED: wrote nothing
- [!] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages" — investigate only: does Astrate have an analogous delay between trigger install and first delivery? Report, do not patch. — BLOCKED: wrote nothing
- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
- [x] docs-sync-pairing-health-path [auto]: add the undocumented `GET /pairing/v1/{realm}/health` route to docs/api/astarte_pairing_api.yaml — it exists in code (internal/pairing/http.go:81) and returns 200 `{"data":{"status":"ok"}}` for an existing realm, 404 for an unknown realm, 503 when unhealthy, unauthenticated and unrate-limited. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-register-404 [auto]: document the missing `404` (DeviceNotFound) on `POST /pairing/v1/{realm}/agent/devices` in docs/api/astarte_pairing_api.yaml — the handler returns it when `GetRealmByName` fails for an unknown realm (internal/pairing/service.go:183-186 → http.go:326-327), but the spec lists only 201/400/401/403/422/429/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] swagger-httptest-coverage [auto]: add a container-free httptest suite for internal/swagger, which currently has no test file — cover two rules: (1) `Mount` wires `GET /swagger` to a 302→/swagger/index.html, serves the embedded UI under `/swagger/`, and serves the OpenAPI YAML specs under `/api/`; (2) `Specs` returns exactly the `.yaml` filenames from `docs.APIYAML` with no path prefix or dirs. httptest only, no Docker.
- [!] probe-property-resend-encoding [auto]: upstream v1.3.0 fixed outbound server-property values sent to a device on connect/emptyCache (commit 522ccf4f — a raw stored binaryblob is now wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding, so it ships as a BSON binary, not a string). Investigate only: does Astrate's resendServerProperties/rehydrate path (internal/engine/control.go:144, pkg/payload) emit binaryblob → BSON subtype-0 binary and datetime → UTC datetime exactly like upstream 1.3? Report, do not patch. — BLOCKED: wrote nothing
- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing

Let me note the todo.md format. Now let me write the review file. My findings:

1. **Missing test (category 2, best kind):** `Pipeline.Validate` has no direct test of the `hasSource` and `hasSink` error branches — the two existing test cases ("no source", "no sink") both use a cycle (`a→b, b→a`), so they hit the *cycle* error first and never reach the branch they claim to cover. This is a rule that exists but nothing asserts on directly.

2. **Dead/duplicated code (category 4, clarity, mechanical):** `pipeline.go:117-126` computes `hasSource`/`hasSink` from the topo-sorted (zeroed) `inDeg`/`outDeg`, then lines 128-147 recompute from fresh `inDeg2`/`outDeg2` and overwrite them. Lines 117-126 are dead.

Not proposing:
- The `onDrop` field on `flowMsg` is never set anywhere (nil-only), so `if fm.onDrop != nil` in `router.go:140,152` are dead branches — but removing it is a tiny refactor; I'll leave it as a note, not a task (it's a deliberate API extension point and not worth a queue line).
- The `topoOrder`/`Validate` both compute graph properties — duplicated logic between pipeline.go and factory.go, but refactoring across two files is out of scope per the recipe.
- No correctness bugs found in router/flow concurrency — the locking is well-reasoned and commented; I verified the Drain/Submit fencing is sound.

Let me write the review file.
← Write .mule/reviews/flow-2026-09-05.md
Wrote file successfully.

Wait, I need to reconsider finding 1. Let me re-examine whether there's actually a way to test the no-source branch that isn't a cycle. Looking at `Validate`:

- It first detects cycles (lines 100-115). If there's a cycle → returns.
- Then checks `hasSource`/`hasSink`.

A graph where every block has an incoming edge but is acyclic is impossible — a DAG must have at least one source. So the only way to reach `hasSource == false` is via a cycle, which short-circuits first. Similarly for sink.

This means in `Validate`, the `hasSource`/`hasSink` branches are **unreachable whenever the graph is acyclic**, and only reachable when it's cyclic — but then the cycle check fires first. So they are **genuinely dead branches**, not just untested. That's a stronger finding.

Wait — is the cycle check exhaustive before? The cycle detection uses one Kahn pass and compares `visited != len(p.Blocks)`. That catches all cycles. So yes: any graph reaching the source/sink check is acyclic, and any acyclic graph must have ≥1 source and ≥1 sink. So `hasSource` and `hasSink` are *always* true at that point → the two error branches at lines 148-153 are unreachable dead code, alongside the dead first computation at 117-126.

This strengthens the case: both the source/sink checking as a whole is dead. Let me reconsider my review content to reflect this accurately. The whole source/sink validation block (both the dead first computation AND the second computation) is unreachable because it's guarded by the cycle check that precedes it.

Actually, let me be careful — is the second computation (lines 128-147) also dead? Yes, because it's only reached for acyclic graphs, and acyclic graphs always have a source and sink. So the `if !hasSource`/`if !hasSink` error returns are unreachable. The whole source/sink block described in the doc comment "verify there is at least one source and one sink" is vacuous.

This is a cleaner, more valuable finding than my draft suggests. Let me update the review file to state this precisely, and adjust the proposed task: the source/sink validation cannot ever fire, so propose removing the dead source/sink checking AND the misleading dead test cases, OR (better-value framing) recognize the rule is vacuous. But removing validation that documents intent is a judgement call.

Let me reconsider. The recipe values: category 1 (wrong/unguarded), category 2 (missing test on a rule that exists), category 4 (clarity). 

The accurate finding: **`Pipeline.Validate` asserts "has no source"/"has no sink" errors that can never occur** — the graph-cycle check that precedes them guarantees every acyclic graph has a source and a sink, so these branches are unreachable; and the tests claiming to cover them actually exercise the cycle path. This is both dead code (category 4) and a misleading test (the tests assert an error but not the named one).

The cleanest, machine-checkable proposal — a missing-test framing: make the no-source/no-sink tests fail honestly, which they would, revealing the validation is vacuous. Actually the best deliverable is a task to make the source/sink validation *real* if the intent matters, or delete the vacuous block. Given the doc comment says "verify there is at least one source and one sink", and that's the documented rule, the highest-value task is a **correctness** check: does anything depend on source/sink validation? Since it can never fire, this is unguarded/dishonest validation.

Let me refine the review file to present this accurately and cleanly. Let me rewrite the findings section.
← Write .mule/reviews/flow-2026-09-05.md
Wrote file successfully.

Now let me append the task lines to todo.md. I'll propose at most 5 lines. The main one is the pipeline.Validate rule. Let me frame task lines.
→ Read .mule/todo.md [offset=1, limit=5]

$ cat >> .mule/todo.md <<'EOF'
- [ ] flow-validate-source-sink: `Pipeline.Validate`'s no-source/no-sink error branches (internal/flow/pipeline.go:148-153) are unreachable — the preceding Kahn cycle check (pipeline.go:100-115) guarantees any graph that gets past it is acyclic and therefore always has a source and a sink. The two test cases "no source (all blocks have incoming edges)" and "no sink (all blocks have outgoing edges)" (pipeline_test.go:134-163) both use a 2-cycle, so they pass via the cycle error, never the branch they name. Make the rule real (a cycle-independent formulation, e.g. reject a graph where a node has no path to a sink) or delete the vacuous check, and fix the two test cases to assert the actual error message so a regression to wrong-path coverage fails.
- [ ] flow-validate-dead-source-sink-recompute: in `Pipeline.Validate`, the first `hasSource`/`hasSink` computation (pipeline.go:117-126) is overwritten by the second (pipeline.go:128-147) and its result discarded — dead code (the comment at pipeline.go:127 documents it). Delete the first block; behaviour is unchanged. Tests: existing `TestPipeline_Validate` passes.
EOF
echo done
done

Let me verify the todo.md appends look clean and run the gate to confirm the package still builds (no code changed, so it should be identical).
$ go test ./internal/flow/ 2>&1 | tail -3; echo "---"; git status --short
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
---
 M .mule/todo.md
?? .mule/reviews/flow-2026-09-05.md

The recipe is a proposal job; I did not implement anything or touch git beyond editing the mule files (which is the deliverable). No `gh issue create` required — this recipe's output is review file + todo lines.

```
Done: Reviewed internal/flow (#85 composites, router/flow concurrency, pipeline validation), proposed 2 tasks
Files: .mule/reviews/flow-2026-09-05.md, .mule/todo.md
Verified: go test ./internal/flow/ -> pass; gofmt -l internal/flow/ -> clean
Unsure: nothing (foundings 1+2 confirmed against tests/code)
Follow-ups: flow-validate-source-sink, flow-validate-dead-source-sink-recompute
```
