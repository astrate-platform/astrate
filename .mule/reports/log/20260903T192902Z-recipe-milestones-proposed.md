slug: recipe-milestones
verdict: proposed
at:  add4eb7
ran: 2026-09-03T19:29:02Z on DietPi in 464s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md

# Todos
[•] Investigate astarte_flow upstream capabilities from its repo/docs
[ ] Check what Astrate already covers (internal/, docs/)
[ ] Check existing milestone-2.0 issues
[ ] File issues or escalate decisions for each gap

% WebFetch https://github.com/astarte-platform/astarte_flow
% WebFetch https://github.com/astarte-platform/astarte_flow/blob/master/README.md
% WebFetch https://docs.astarte-platform.org/flow/snapshot/
% WebFetch https://github.com/astarte-platform/astarte_flow/tree/master/src
% WebFetch https://github.com/astarte-platform/astarte_flow/tree/master/lib
% WebFetch https://github.com/astarte-platform/astarte_flow/tree/master/guides
% WebFetch https://github.com/astarte-platform/astarte_flow/blob/master/src/pipelineparser.yrl
✗ WebFetch https://docs.astarte-platform.org/flow/snapshot/010-astarte_flow.html failed
Error: StatusCode: non 2xx status code (404 GET https://docs.astarte-platform.org/flow/snapshot/010-astarte_flow.html)
✗ WebFetch https://docs.astarte-platform.org/flow/snapshot/020-blocks.html failed
Error: StatusCode: non 2xx status code (404 GET https://docs.astarte-platform.org/flow/snapshot/020-blocks.html)
% WebFetch https://github.com/astarte-platform/astarte_flow/tree/master/lib/astarte_flow
% WebFetch https://github.com/astarte-platform/astarte_flow/tree/master/lib/astarte_flow/blocks
% WebFetch https://github.com/astarte-platform/astarte_flow/tree/master/lib/astarte_flow_web

# Todos
[✓] Investigate astarte_flow upstream capabilities from its repo/docs
[•] Check what Astrate already covers (internal/, docs/)
[ ] Check existing milestone-2.0 issues
[ ] File issues or escalate decisions for each gap

Good, I now have a clear picture of astarte_flow's block types and architecture. Let me now investigate Astrate's codebase and check existing issues.
• Explore Astrate flow coverage Explore Agent

$ cd /root/astrate-mule && gh issue list --label milestone-2.0 --state all --limit 50 2>/dev/null || echo "NO_ISSUES_FOUND"
43	CLOSED	flow-container-block: PoC → MVP (custom processing in containers)	enhancement, milestone-2.0	2026-07-29T18:54:52Z
42	CLOSED	flow-rehydrate-edge-cases: follow-ups after durable auto_restart MVP	enhancement, milestone-2.0	2026-07-29T20:39:19Z
41	CLOSED	flow-durable: persist flows + auto_restart rehydrate on process boot	enhancement, milestone-2.0	2026-07-29T18:54:49Z
40	CLOSED	flow-multi-instance: named flows + pipeline config (decided — design then implement)	enhancement, question, milestone-2.0	2026-07-29T18:54:47Z
39	CLOSED	flow-blocks-discovery: GET /flow/v1/{realm}/blocks for registered catalog types	enhancement, milestone-2.0	2026-07-29T18:54:44Z
37	CLOSED	flow-source-pump: poll AstarteSource and Stop on flow teardown	milestone-2.0	2026-07-29T14:13:12Z
27	CLOSED	flow-astarte-source: Astarte device events source block	mule-review, milestone-2.0	2026-07-29T07:45:22Z
26	CLOSED	flow-router: Stream-based message routing through the block graph	milestone-2.0	2026-07-29T07:33:30Z
25	CLOSED	flow-lifecycle: Flow start/stop lifecycle and runtime management	milestone-2.0	2026-07-29T07:33:28Z
24	CLOSED	flow-pipeline-store: Pipeline CRUD storage and graph validation	mule-review, milestone-2.0	2026-07-29T07:45:19Z
23	CLOSED	flow-models: FlowMessage wire format and Block/Pipeline/Flow data models	milestone-2.0	2026-07-29T07:33:26Z

✓ Explore Astrate flow coverage Explore Agent
→ Read internal/flow
Excellent. The explore agent found that `internal/flow/` has the core framework (Block, Pipeline, Router, Manager, FlowMessage) but **no HTTP endpoints, no concrete block implementations beyond the interface, and no broker integration**. All 11 milestone-2.0 issues are CLOSED.

Let me look more closely at what's actually in the codebase now vs what the closed issues covered.
✱ Grep "flow|pipeline|block" in internal/appengine · 4 matches
✱ Grep "flow|pipeline|block" in internal/realm · 1 match

$ cd /root/astrate-mule && git log --oneline -30
add4eb7 mule: blocked store-downsample-lttb
d2c1ccd mule: blocked housekeeping-unit-tests
79f9c1e mule: refill the queue
6874214 mule: recipe hygiene ran (184s)
6b7ade3 mule: refill the queue
6a69d53 mule: recipe docs-sync ran (229s)
da12834 mule: blocked object-path-ok-comment
eb3b223 mule: blocked consumer-properties-plaintext-fallback
19d552b mule: blocked serverdata-aftercommit-test
bde026d mule: refill the queue
5f27771 mule: recipe code-review ran (230s)
2d82565 mule: refill the queue
3debf0d mule: recipe astarte-upstream ran (105s)
4ed448a mule: refill the queue
f72e460 mule: recipe github-issues ran (103s)
c7deb30 mule: refill the queue
b1fd205 mule: recipe milestones ran (336s)
262f6f6 mule: blocked lint-flow-stutter
2cbc7d3 mule: blocked lint-pairing-test-param
eab8e24 mule: blocked lint-flow-findings
154ee25 mule: blocked lint-engine-goimports
6750a4b mule: refill the queue
59250ec mule: recipe hygiene ran (279s)
88ed968 mule: blocked docs-sync-appengine-socket
8e955e7 mule: blocked docs-sync-appengine-422
ca9fcad mule: refill the queue
7ff4280 mule: recipe docs-sync ran (316s)
4ab5680 mule: blocked flow-datawirescalar-fallthrough
a88cdb3 mule: blocked flow-unmarshal-error-tests
04ac1b8 mule: blocked flow-validate-deadcode

$ cd /root/astrate-mule && git branch -a | head -20
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

$ cd /root/astrate-mule && git log --oneline --all | grep -i flow | head -20
262f6f6 mule: blocked lint-flow-stutter
eab8e24 mule: blocked lint-flow-findings
4ab5680 mule: blocked flow-datawirescalar-fallthrough
a88cdb3 mule: blocked flow-unmarshal-error-tests
04ac1b8 mule: blocked flow-validate-deadcode
4c636be mule: blocked flow-setstatus-race
84ae2b6 flow: dynamic virtual_device_pool — first-seen auto-registration (#84)
14708db feat(flowapi): composites accepted and expanded in the flow lifecycle (issue #85, phase 85d)
4b6dfc7 feat(flowapi): user-block CRUD endpoints with schema validation (issue #85, phase 85c)
78e293e feat(flow): ExpandComposites — inline expansion of user blocks with cycle guard (issue #85, phase 85a)
6d60193 feat(flowapi): plumb the engine ingest path into flow Deps (issue #84, phase 84c)
ad55cab feat(flow/blocks): virtual_device_pool block over the engine ingest path (issue #84, phase 84b)
970ce71 feat(flow/blocks): register virtual_device_pool type (stub) and Ingest dep (issue #84, prep 84b)
abd1556 feat(flow/blocks): mqtt_source and mqtt_sink blocks (issue #83)
671982e feat(flow/blocks): serve config schemas from the blocks-discovery API (issue #89)
55bb717 feat(auth,flowapi): a_f JWT claim for the Flow API (issue #88)
4d456bd feat(flow/blocks): http_source and http_sink blocks (issue #82)
7f91d26 feat(flow/blocks): json_path_map block (issue #81)
4fa0e3d feat(flow/blocks): pure-transform block set — to_json, update_metadata, split_map, random_source, sort (issue #80)
fb3f335 feat(flow): block-death detection + auto-restart with backoff (#45 phase 1)

→ Read .mule/todo.md
→ Read docs/DESIGN.md [offset=55, limit=15]

$ cd /root/astrate-mule && gh issue list --label milestone-2.0 --state open --limit 50 2>/dev/null || echo "NO_OPEN_ISSUES"
(no output)

→ Read .mule/for-giulio.md

$ cd /root/astrate-mule && gh issue list --label milestone-2.0 --state all --limit 50 2>/dev/null | grep -v "^[0-9]" || true
(no output)

$ cd /root/astrate-mule && git show origin/main:.mule/milestones.md 2>/dev/null | head -60
# Milestones

Release-tag-gated goals. Each section names a body of work that must be true by the time
that tag is cut. `.mule/recipes/milestones.md` reads this file to find the **current
target** — the first section below not marked `DONE` — and works toward it: investigating,
filing GitHub issues (`milestone-<tag>` label, alongside `mule`) with sub-issues where the
work splits, and escalating anything that needs a design decision to `.mule/for-giulio.md`.

**This file is Giulio's.** The recipe reads it and may propose edits via
`.mule/for-giulio.md`, but never edits it directly — same rule as `docs/COMPATIBILITY.md`.
Mark a milestone `DONE` yourself once the tag is actually cut.

Order matters: milestones are worked **in order**, lowest tag first. Do not start v3.0
investigation while v2.0 has open, un-escalated gaps — say so and stop instead.

---

## v2.0 — astarte-flow feature parity

Reference: astarte_flow (https://github.com/astarte-platform/astarte_flow) — upstream's
Elixir "Flow" component. It lets ingested data get piped through a graph of processing
blocks (native blocks and containerised ones) before it lands in storage, with pipelines
described as reusable, parametrised graphs.

Scope for this milestone: whatever set of Flow's *capabilities* Astrate needs to expose the
same wire-visible behaviour and operator-facing concepts (pipelines, blocks, native vs.
containerised blocks) — not a port of the Elixir implementation. See
`.mule/recipes/astarte-upstream.md`'s rule: port the idea, restated in Go, never the code.

Status: **DONE** (2026-07-29). Runtime, factory, catalog (incl. filter/map), process
wiring, and `/flow/v1` API are on `main`. **Parity audit** + **product decisions**
recorded 2026-07-29 (`docs/handoff/flow-parity-audit-2026-07-29.md`,
`docs/handoff/flow-v2-decisions-2026-07-29.md`).

**Design A + B landed on `main`** (commit `89145e6`, 2026-07-29): durable flows +
auto_restart (**#41 closed**); named multi-instance + config (**#40 closed**) —
migration `000009`, store, `${config.*}`, API, boot rehydrate. Container block
PoC→MVP (**#43 closed**) — registered in the catalog, usable inside stored
pipelines/named flows; see `docs/handoff/flow-design-b-container-block-2026-07-29.md`.
Blocks discovery (**#39 closed**).
**Not a v2.0 gate:** native Lua / MQTT blocks.

**#42 closed 2026-07-29** (rehydrate edge cases). Triaged all seven candidates against a
live e2e Docker smoke test of #43: one real bug found and fixed on `main`
(`801fc48` — `Registry.Instantiate` leaked an already-started Docker container when a
later block's constructor failed mid-pipeline-build; mutation-tested regression added);
two candidates ("one flow fails at boot, others still start"; "pipeline deleted while
flows reference it") were already correct by design, no code change needed; three
("hot-reload a running pipeline" #44, "partial restart of failed blocks" #45, "update a
running flow's config" #46) split into their own demand-driven backlog issues per the
explicit v2.0 decisions doc; multi-process/HA managers stay out of scope indefinitely
(single-process design throughout).

### Landed (on main as of 2026-07-29)

| Piece | Where | Notes |
|---|---|---|
| FlowMessage wire format | `internal/flow/message.go` | `astarte_flow/message/v0.1` |
| Block / Source / Stopper / graph | `internal/flow/block.go`, `graph.go` | linear chain; Sources skipped in `Run` |
| Pipeline DAG model + validate | `internal/flow/pipeline.go` | serialisable description |

The `for-giulio.md` already has an entry (lines 45-58) from a previous milestones run that explains v2.0 is DONE on `main` and the `mule/queue` branch is stale. Let me verify the current state.
$ cd /root/astrate-mule && gh issue list --label milestone-2.0 --state all --limit 50 2>/dev/null | head -20
43	CLOSED	flow-container-block: PoC → MVP (custom processing in containers)	enhancement, milestone-2.0	2026-07-29T18:54:52Z
42	CLOSED	flow-rehydrate-edge-cases: follow-ups after durable auto_restart MVP	enhancement, milestone-2.0	2026-07-29T20:39:19Z
41	CLOSED	flow-durable: persist flows + auto_restart rehydrate on process boot	enhancement, milestone-2.0	2026-07-29T18:54:49Z
40	CLOSED	flow-multi-instance: named flows + pipeline config (decided — design then implement)	enhancement, question, milestone-2.0	2026-07-29T18:54:47Z
39	CLOSED	flow-blocks-discovery: GET /flow/v1/{realm}/blocks for registered catalog types	enhancement, milestone-2.0	2026-07-29T18:54:44Z
37	CLOSED	flow-source-pump: poll AstarteSource and Stop on flow teardown	milestone-2.0	2026-07-29T14:13:12Z
27	CLOSED	flow-astarte-source: Astarte device events source block	mule-review, milestone-2.0	2026-07-29T07:45:22Z
26	CLOSED	flow-router: Stream-based message routing through the block graph	milestone-2.0	2026-07-29T07:33:30Z
25	CLOSED	flow-lifecycle: Flow start/stop lifecycle and runtime management	milestone-2.0	2026-07-29T07:33:28Z
24	CLOSED	flow-pipeline-store: Pipeline CRUD storage and graph validation	mule-review, milestone-2.0	2026-07-29T07:45:19Z
23	CLOSED	flow-models: FlowMessage wire format and Block/Pipeline/Flow data models	milestone-2.0	2026-07-29T07:33:26Z

# Todos
[✓] Investigate astarte_flow upstream capabilities from its repo/docs
[✓] Check what Astrate already covers (internal/, docs/)
[✓] Check existing milestone-2.0 issues
[•] File issues or escalate decisions for each gap

All 11 `milestone-2.0` issues are CLOSED, no open issues exist, and no new gaps remain. The `for-giulio.md` already has a detailed entry (lines 45-58) from a previous run confirming v2.0 is DONE on `main`. Per recipe step 5, I'll add a concise re-confirmation line.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,17 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **milestone 2.0 looks complete, verify and cut the tag** — all 11 `milestone-2.0` issues
+  CLOSED (#23–#27, #37, #39–#43), no open issues, no new gaps after re-checking upstream
+  astarte_flow block catalog against `internal/flow/` + git log (MQTT/HTTP source/sink,
+  json_path_map, pure-transform set, virtual_device_pool, container block MVP, flow API,
+  durable named flows all landed). `mule/queue` branch is stale — see entry below.
+  (Milestones recipe run, 2026-09-03.)
+
+---
+
 - **v3.0 decide #67** (`upstream-parity`): `required` and `encrypted` mapping fields exist on upstream master but not at 1.2.2/1.3.0 (`required` = missing key at runtime → `missing_required_mapping` on object-aggregated mappings; `encrypted` ties into keyAgreement #51). Decide now: reject (strict 1.2.2), ignore-and-accept, or implement `required` early. Issue is OPEN, no decision on record, not yet escalated here.
 
 ---
 

# Todos
[✓] Investigate astarte_flow upstream capabilities from its repo/docs
[✓] Check what Astrate already covers (internal/, docs/)
[✓] Check existing milestone-2.0 issues
[✓] File issues or escalate decisions for each gap

Done: v2.0 milestone re-confirmed complete; added for-giulio.md line per recipe step 5
Files: `.mule/for-giulio.md`
Verified: `gh issue list --label milestone-2.0 --state open` → empty; `gh issue list --label milestone-2.0 --state all` → 11 CLOSED
Unsure: nothing
Follow-ups: none (v2.0 done, no new issues filed; branch needs rebuild for v3.0 target)
