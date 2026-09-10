slug: recipe-code-review
verdict: proposed
at:  f463039
ran: 2026-09-10T17:01:49Z on DietPi in 513s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
→ Read .mule/reviews
→ Read internal
→ Read .mule/log.md
→ Read pkg
→ Read .mule/todo.md

$ cd /root/astrate-mule && rg -c '' --glob 'internal/auth/**/*.go' 2>/dev/null | sort -t: -k2 -rn; echo "---"; rg -c '' --glob 'internal/flowapi/**/*.go' 2>/dev/null | sort -t: -k2 -rn; echo "---"; rg -c '' --glob 'internal/housekeeping/**/*.go' 2>/dev/null | sort -t: -k2 -rn; echo "---"; rg -c '' --glob 'internal/httpx/**/*.go' 2>/dev/null | sort -t: -k2 -rn; echo "---"; rg -c '' --glob 'internal/observability/**/*.go' 2>/dev/null | sort -t: -k2 -rn; echo "---"; rg -c '' --glob 'pkg/**/*.go' 2>/dev/null | sort -t: -k2 -rn; echo "---"; rg -c '' --glob 'internal/config/**/*.go' 2>/dev/null | sort -t: -k2 -rn
internal/auth/middleware_test.go:271
internal/auth/claims_test.go:249
internal/auth/jwt_test.go:248
internal/auth/claims.go:204
internal/auth/middleware.go:196
internal/auth/cache_test.go:144
internal/auth/jwt.go:139
internal/auth/cache.go:81
---
internal/flowapi/service.go:976
internal/flowapi/service_integration_test.go:334
internal/flowapi/http.go:304
internal/flowapi/userblocks.go:209
internal/flowapi/service_test.go:137
internal/flowapi/composite_wire_test.go:133
internal/flowapi/userblocks_test.go:97
internal/flowapi/http_test.go:40
---
internal/housekeeping/http_test.go:511
internal/housekeeping/service.go:280
internal/housekeeping/http.go:259
internal/housekeeping/service_test.go:224
internal/housekeeping/async_operation_test.go:37
---
internal/httpx/notfound_test.go:139
internal/httpx/cors_test.go:113
internal/httpx/notfound.go:56
internal/httpx/cors.go:49
---
internal/observability/observability_test.go:128
internal/observability/health.go:83
internal/observability/metrics.go:80
internal/observability/compat_version_test.go:55
internal/observability/compat.go:52
---
pkg/payload/payload_test.go:787
pkg/interfaceschema/parse.go:664
pkg/interfaceschema/types.go:406
pkg/payload/json.go:396
pkg/payload/value.go:353
pkg/interfaceschema/parse_test.go:352
pkg/payload/bson.go:343
pkg/interfaceschema/violations_test.go:309
pkg/interfaceschema/compat_test.go:304
pkg/astarteapi/envelope.go:249
pkg/astarteapi/envelope_test.go:236
pkg/interfaceschema/compile_test.go:235
pkg/deviceid/deviceid_test.go:226
pkg/payload/value_test.go:223
pkg/payload/payload.go:217
pkg/interfaceschema/trie_test.go:216
pkg/interfaceschema/required_encrypted_test.go:194
pkg/interfaceschema/trie.go:161
pkg/deviceid/deviceid.go:132
pkg/interfaceschema/compat.go:131
pkg/payload/fuzz_test.go:122
pkg/interfaceschema/compile.go:122
pkg/payload/sniff_test.go:120
pkg/interfaceschema/violations.go:82
pkg/payload/required_test.go:71
pkg/payload/sniff.go:41
---
internal/config/config.go:411
internal/config/config_test.go:407

→ Read internal/flowapi/service.go
I'll review `internal/flowapi` — it's sizeable (~2500 lines), unrevealed, and it's the HTTP/service layer for the flow engine. Let me read the main files.
→ Read internal/flowapi/http.go
→ Read internal/flowapi/userblocks.go
→ Read internal/flowapi/service_test.go
→ Read internal/flowapi/userblocks_test.go
→ Read internal/flowapi/http_test.go
→ Read internal/flowapi/composite_wire_test.go
→ Read internal/flowapi/service_integration_test.go

$ cd /root/astrate-mule && rg -c '' --glob 'internal/flow/*.go' | sort -t: -k2 -rn
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
internal/flow/composite.go:156
internal/flow/pipeline.go:155
internal/flow/substitute.go:126
internal/flow/block.go:113
internal/flow/substitute_test.go:85
internal/flow/graph.go:64

Let me check what the `flow` package's tests cover to avoid duplicating proposals.
$ cd /root/astrate-mule && ls internal/flow/*_test.go 2>/dev/null
internal/flow/composite_test.go
internal/flow/factory_test.go
internal/flow/flow_test.go
internal/flow/manager_test.go
internal/flow/message_test.go
internal/flow/pipeline_test.go
internal/flow/pump_test.go
internal/flow/router_test.go
internal/flow/substitute_test.go

$ cd /root/astrate-mule && rg -n "func (m \*Manager) Shutdown|func (m \*Manager) StartFlow|func (m \*Manager) StopFlow|func (m \*Manager) ListFlows" internal/flow/flow.go internal/flow/*.go
(no output)

$ cd /root/astrate-mule && rg -n "func ParseDefinition|func (m \*Manager) UnregisterFlow" internal/flow/
internal/flow/factory.go:214:func ParseDefinition(id, name string, definition []byte) (*Pipeline, error) {

Let me verify a couple of claims before finalising: how `Manager.Shutdown` behaves, and the `ParseDefinition` signature.
$ cd /root/astrate-mule && rg -n "func \(m \*Manager\) (Shutdown|StartFlow|StopFlow|ListFlows|UnregisterFlow)" internal/flow/
internal/flow/flow.go:144:func (m *Manager) StartFlow(ctx context.Context, cfg Config) (*Flow, error) {
internal/flow/flow.go:187:func (m *Manager) UnregisterFlow(instanceID string) {
internal/flow/flow.go:248:func (m *Manager) StopFlow(ctx context.Context, pipelineID string) error {
internal/flow/flow.go:299:func (m *Manager) ListFlows() []*Flow {
internal/flow/flow.go:312:func (m *Manager) Shutdown(ctx context.Context) error {

→ Read internal/flow/flow.go [offset=100, limit=230]

$ cd /root/astrate-mule && rg -n "MarkRunningFlowsStopped|flowapi|flowSvc|flowapiSvc|Service\)\.Manager|\.Manager\(\)" cmd/astrate/main.go
40:	"github.com/astrate-platform/astrate/internal/flowapi"
160:	flowSvc := flowapi.NewService(st, flowMgr, blocks.DefaultRegistry(), e.Bus(),
184:	handler, hkSvc, err := mountAPIs(cfg, st, e, b, sealer, metrics, flowSvc, pairer, log)
186:		shutdown(nil, b, flowSvc, e, log)
210:			shutdown(nil, b, flowSvc, e, log)
216:	if err := flowSvc.RehydrateAutoRestart(ctx); err != nil {
217:		shutdown(nil, b, flowSvc, e, log)
236:		shutdown(nil, b, flowSvc, e, log)
240:	shutdown(srv, b, flowSvc, e, log)
350:func mountAPIs(cfg config.Config, st *store.Store, e *engine.Engine, b *broker.Broker, sealer *store.KeySealer, metrics *observability.Metrics, flowSvc *flowapi.Service, pairer *pairing.Service, log *slog.Logger) (http.Handler, *housekeeping.Service, error) {
379:	flowapi.NewAPI(flowSvc, mw).Mount(mux)
466:// error before the HTTP server began serving. flowSvc may be nil when Flow
468:func shutdown(srv *http.Server, b *broker.Broker, flowSvc *flowapi.Service, e *engine.Engine, log *slog.Logger) {
482:	if flowSvc != nil {
484:		if err := flowSvc.Manager().Shutdown(sctx); err != nil {
487:		flowSvc.MarkRunningFlowsStopped(sctx)

→ Read cmd/astrate/main.go [offset=466, limit=40]

$ cd /root/astrate-mule && ls .mule/for-giulio.md 2>/dev/null && cat .mule/for-giulio.md
.mule/for-giulio.md
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **github-issues triage run, 2026-09-10: nothing proposable, nothing stale.** Nine open
  issues: #99–#94 are mule-alarms (idle queue), not code issues. #93 has a pushed commit
  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated. #1
  untouched per standing instruction. No task lines proposed, no for-giulio close suggestions.

- **Milestone recipe run, 2026-09-09: v3.0 unchanged since 2026-09-06 — still waiting on
  upstream gates.** `milestone-3.0` label remains empty. The three open `upstream-parity`
  issues are the same: #92 keyAgreement (upstream decision, gated on stable v1.4.0), #93
  aclhook comment (pushed commit `8c61268`, mule-review), #78 FDO (milestone-4.0). Both
  UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait
  on upstream v1.4 final (currently rc.5). `APICompatVersion` stays 1.2.2 per #90's frozen
  decision. No new gaps found, no issues filed, no task lines proposed.

---

- **github-issues triage run, 2026-09-08: nothing proposable, nothing stale.** All six open
  issues were already covered by another mechanism. #97/#96/#95/#94 are mule-alarms (idle
  queue), not code issues — expected during low-activity windows. #93 has a pushed commit
  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated below.
  #1 untouched per standing instruction. No task lines proposed, no for-giulio close
  suggestions. (Identical to the 2026-09-07 run — issue list unchanged except the new
  alarm #97.)

---

- **Milestone recipe run, 2026-09-06: v3.0 (upstream 1.2.2 → 1.3/1.4) has no open
  machine-checkable gaps left to file; the milestone now waits on two upstream gates, and the
  `milestone-3.0` label is empty (this milestone's work lives under `upstream-parity`).** All
  the 1.3-surface gaps from the #47–#89 backlog are delivered/closed — retention ceiling #72,
  alias/group #52–#54 #59, query formats #55 #56, wire capabilities #47–#49, validation /
  error-code fidelity #57 #61 #62 #79, per-service version #77, housekeeping #73–#76 — and the
  remaining open items are already escalated, not fileable gaps: **#92 keyAgreement** (1.4
  experimental decision, escalated 2026-09-05, gated on a stable upstream v1.4.0), **#78 FDO**
  (milestone-4.0, escalated), **#93** aclhook comment rewrite (mule-review, pushed commit
  `8c61268`). `docs/UPSTREAM-EXPERIMENTAL.md` rows #67/#68 are "1.4 experimental" and reconcile
  only when upstream ships **v1.4 final** — not yet (still rc). Per **#90's frozen decision**
  (2026-08-23), `APICompatVersion` stays 1.2.2 until the FULL 1.3 surface is complete *and*
  every UPSTREAM-EXPERIMENTAL row at that level is reconciled, so the bump is not yet due. Your
  call when v1.4.0 goes stable: answer #92 (implement against the rc, re-park, or narrow #93
  only) — until then v3.0 has no open implementable work. This run filed **no issues** and
  proposes **no** "complete, cut the tag" line: the open decision set is real and is yours.

- **`milestones.md` v3.0 names a source doc that does not exist: `.mule/research/
  upstream-parity-2026-08-22.md` is absent from `main` and from `origin/mule/research`.** The
  investigation's *output* (issues #47–#89) is intact and reachable, so milestone work is
  unaffected, but if that doc is your canonical 1.3/1.4 diff it needs restoring (or the
  milestone paragraph re-pointed) before v4.0 drafting leans on it.

- **docs-sync realm-management: two ASTRATE_ config keys exist in code but are missing from
  `docs/site/configuration-reference.md`** (site prose is yours, never edited by the mule).
  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` and
  `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` are consumed in `internal/config/config.go`
  (lines 270-271 and 288); neither appears in `docs/site/configuration-reference.md`.
  Add them or confirm they are meant to be internal-only.

- **github-issues triage run, 2026-09-06: nothing proposable, nothing stale.** All six open
  issues were already covered by another mechanism. #95 and #94 are mule-alarms (idle queue),
  not code issues — expected during low-activity windows. #93 has a pushed commit (`8c61268`)
  awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the waiting-on
  row. #78 is the milestone-4.0 FDO design/investigation, already escalated below. #1
  untouched per standing instruction. No task lines proposed, no for-giulio close suggestions.
  (Identical to 2026-09-05 run — issue list unchanged except new alarm #95.)

---

- **keyAgreement: the parking condition from #51 has fired — implement now, or wait for a
  stable v1.4.0?** (issue #92, `upstream-parity`/`upstream-experimental`). #51 closed
  2026-08-22 with "parked until the upstream 1.4 experimental spec stabilizes — reopen or
  file fresh when it does". The document side has now stabilized: upstream `d084308`
  (2026-08-31) published `082-key_agreement_protocol.md`, a full 267-line wire protocol
  (topics `control/keyAgreement/0..4` — InitExchange/ExchangeResp/SecretHash/HashOk/
  ExchangeFailed at QoS 2, CBOR bodies with CDDL, `alg` 0 ECDH_P256-HKDF_SHA256-AES_256_GCM,
  CBOR-wrapped COSE_Key + 32-byte HkdfSalt, session-scoped keys, enumerated `ExchangeFailed`
  codes), and deleted the "not yet implemented" sentences the old parking quoted (#93 already
  fixed the stale ACL comment that cited them). But the spec ships only in `v1.4.0-rc.5` —
  `v1.3.3` is still the newest stable tag and Astrate targets 1.2.2 — so the document has
  stabilized and the release has not. Implementing it is the largest surface in the parity
  backlog (CBOR codec, X25519/P-256, HKDF, AES-256-GCM, a 5-state handshake machine,
  shared-secret persistence, five new error names). Your call: build against the rc now,
  re-park until v1.4.0 is a stable tag, or take only the narrow #93 fix. (Escalated again
  2026-09-05 — a prior escalation from the 2026-09-04 milestone run was lost in the queue
  rebuild.)

---

- **govulncheck GO-2026-5970: reachable DoS in golang.org/x/text (infinite loop on invalid input, fixed in v0.39.0, available v0.41.0).** Astrate pins `x/text` indirect at v0.38.0 (go.mod:97) and pgx pulls it into production: `internal/store/notify.go:59` `store.Listen` → `pgx.ConnectConfig` → `unicode/norm.*`. This is the only govulncheck symbol finding that is not test-harness-only: GO-2026-6355/6354 (x/crypto/ssh deadlocked-channel DoS) and GO-2026-6253 (moby/go-archive tar path traversal) are reachable only through testcontainers in `internal/testutil/pg.go`, i.e. never in the deployed binary. `x/text` keeps API compatibility minor-to-minor and the modules Astrate exercises (`unicode/norm` via pgx, `text/language` via jsonschema) are unchanged, so this is a fix Astrate actually needs — the hygiene recipe's highest-priority category. Not a mule task (go.mod never-touch): your decision to bump ≥v0.39.0 now or fold into the next milestone-boundary sweep. Raw: https://pkg.go.dev/vuln/GO-2026-5970. (The 2026-09-04 dep sweep did not list x/text.)

---

- **COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable, 2026-07-14; v1.4.0 is still rc-only).** Astrate's doc and `APICompatVersion` still target upstream **1.2.2** (`internal/realm/service.go:588`); v1.3.0 (2026-05-06) introduced wire-surface changes Astrate does not yet emulate, so this is a decision — adopt v1.3.2 as the compatibility target (then update the doc + bump `APICompatVersion` together, per the bump rule) or keep 1.2.2 and add a "not yet emulated" note. Wire-relevant v1.3.0 deltas (release notes): **MQTT v1 capabilities** incl. `purge_properties_compression_format` (plaintext vs zlib purge — touches the `emptyCache`/`producer,properties` contract COMPATIBILITY.md deviation 1 documents); **empty introspection now allowed**; **device registration triggers** (pairing) and **device deletion started/completed triggers** (RM — the latter two already exist as Astrate deviation 9 emits both around the synchronous delete); **FDO authentication** (pairing, disabled by default); **realm-scoped health** — upstream v1.3 added `GET /pairing/v1/{realm}/health`, which Astrate already serves (`internal/pairing/http.go:78-81`, comment already says "upstream 1.3+"), so deviation 18's wording ("which upstream 404s") is now false against 1.3 and the note should be reworded either way. Proposed doc wording (for your approval, edit to taste): in §Infrastructure differences add a sentence — *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl.

plaintext `purge_properties_compression_format`, empty-introspection allowance, device registration/deletion triggers, experimental FDO pairing auth) are not yet emulated and are out of scope until the milestone that adopts v1.3.2 as the target."* — and reword deviation 18's realm-health note from "which upstream 404s" to "added by upstream v1.3 (Astrate serves it against a 1.2.2 target; kept, matching behavior)". Raw upstream changes: [v1.3.0](https://github.com/astarte-platform/astarte/releases/tag/v1.3.0), [v1.3.2](https://github.com/astarte-platform/astarte/releases/tag/v1.3.2).

---
- **The mule's 42 straight failures were ten lint errors in its own base, not a stale branch.**
  `mule/queue` stopped taking `main` on 2026-07-27 and drifted 120 commits behind, which is
  real and is why it is being rebuilt — but it is not what blocked the work. The lint gate runs
  `golangci-lint run ./...` over the whole repo, and the branch's own `internal/flow/*` code
  carried 10 findings (1 goimports, 1 gosec, 8 revive). So every task failed the gate no matter
  how good the change was, including the four tasks queued to fix those very findings: each
  fixed one and died on the other nine. Base tests and `go vet` pass on that checkout, and
  `golangci-lint` and `govulncheck` have been installed on the Pi the whole time — the two
  things that looked broken were not. `main` is lint-clean, so rebuilding from it clears the
  deadlock. `tools/mule.sh preflight` checks the lint baseline and would have said so on day
  one; nobody ran it between 2026-08-31 and 2026-09-04. (Diagnosed 2026-09-04.)

---

- **Dependency sweep corrected: direct (pinned) deps DO have newer versions** — the 2026-09-02 note said the `go list -m -u` sweep showed "only version-skew on transitive deps", but that run hit the recipe's `head -20` cutoff (all cloud/azure/transitive) and never reached the directly-required modules. Full sweep, 2026-09-04. None of these is a fix this repo *needs*, so no bump is proposed — recorded for the decision. Per module (current → available; breaking change; repo use):
  - `github.com/coder/websocket` v1.8.14 → v1.8.15 — no breaking (patch); used in `internal/appengine/stream/ws.go`, `channels/ws.go`; worth it only for the "transmit in single frame when compression enabled" fix + read-path alloc reduction.
  - `go.etcd.io/bbolt` v1.4.3 → v1.5.0 — bbolt's semver promises no API change between patch/minor, so additive-only; used in `internal/broker/sessionstore.go`; v1.5 adds a data-file size limit and panic-recovery hardening, nothing Astrate needs.
  - `go.mongodb.org/mongo-driver/v2` v2.6.0 → v2.8.2 — the 2.8.0 breaking changes are confined to Queryable Encryption string-query options (`options.Text()`→`String()`); Astrate uses only the raw BSON API (`pkg/payload/bson.go`, `internal/engine/capabilities.go`, `bench/`) and is unaffected.
  - `github.com/nats-io/nats.go` v1.52.0 → v1.53.1 — no breaking; the headline fixes (JetStream `resetOrderedConsumer` race, KV dot-rejection) are paths Astrate does not use — `internal/engine/forward/nats.go` is core NATS publish only.
  - `github.com/prometheus/client_golang` v1.23.2 → v1.24.1 — requires Go ≥1.25 (fine, repo is 1.26.1); the breaking `LabelNames`/remote-api renames don't touch repo usage (`prometheus`/`collectors`/`promhttp` in `internal/observability/metrics.go`, flow/engine metrics); would buy `Gather()` panic-recovery and opt-in `CoalesceGather` scrape-pile-up protection.
  - `github.com/testcontainers/testcontainers-go` v0.43.0 → v0.44.0 (modules/postgres v0.42.0, modules/nats v0.43.0) — breaking in `wait.ForSQL` (callback now takes `network.Port`) and `ImageProvider` (new `PullImageWithPlatform`); Astrate's `internal/testutil/pg.go` looks unaffected but it is test-only anyway.
  - `golang.org/x/crypto` v0.53.0 → v0.56.0 — x/crypto keeps API compatibility; used only for bcrypt in `internal/auth`.
  Note (corrected 2026-09-04): `govulncheck` and `golangci-lint` **are** installed on the Pi (`/root/go/bin`, since 2026-07-28 and 2026-09-01) and `.mule/config` finds them there, so both checks were available; the sweep that produced this list simply ran without invoking govulncheck.

  **Decided 2026-09-04: no bumps.** None of the seven fixes anything this repo has, and each
  one costs a full test run to land. The standing rule instead: re-run this sweep at every
  milestone boundary (the point where `APICompatVersion` or a milestone tag moves), and bump
  only what carries a fix Astrate actually needs. `go.mod` stays on the never-touch list.

---

- **milestone 2.0 looks complete, verify and cut the tag** — all 11 `milestone-2.0` issues
  CLOSED (#23–#27, #37, #39–#43), no open issues, no new gaps after re-checking upstream
  astarte_flow block catalog against `internal/flow/` + git log (MQTT/HTTP source/sink,
  json_path_map, pure-transform set, virtual_device_pool, container block MVP, flow API,
  durable named flows all landed).   (Milestones recipe run, 2026-09-03.)

  **Decided 2026-09-04: the tag is `v0.2.0`, cut on the same day.** The `v2.0`/`v3.0` names
  are milestone names, not release versions — the project is still pre-1.0 and the version
  number keeps its own line rather than jumping two majors to match a milestone label.

---
- ~~`docs/site/appengine-api.md` documents `GET` and `DELETE` on `/appengine/v1/<realm>/groups/<name>`~~
  — **resolved 2026-09-04, and the original note was half wrong.** `GET /groups/{group}` does
  exist (`internal/appengine/http.go:69`); only `DELETE` was absent, and it is absent from
  upstream's own spec too (`docs/api/astarte_appengine_api.yaml` has no `/groups/{group}`
  path at all). So the docs line was spurious, not the code. Line dropped.
---
- ~~`Router.Submit` TOCTOU on the `closed` flag~~ — **fixed 2026-09-04**, and none of the three
  options in the original note was taken: (b) was wrong (a `select` does not save a send on a
  closed channel — it panics either way). The fix is that `Drain` no longer closes the lane
  channels at all; it closes `quit` under a write lock while `Submit` holds the read lock
  across its send, so a lane can never be retired underneath an in-flight sender, and the lanes
  exit on `quit` after draining what is buffered. `TestRouter_SubmitParkedWhenDrainRuns`
  reproduces the old panic deterministically and passes on the fix.
---
- ~~**#87 `lua_map` — needs embedded Lua runtime, parked.**~~ — **closed 2026-09-04.**
  Embedding a Lua VM is not machine-checkable by the mule and Lua flow support is on no active
  roadmap. Reopen in a minute if that changes.
- **#78 FDO device onboarding — milestone-4.0, investigation phase.** Too large for a single
  mule task; the investigation work (reading upstream's TO2 handling, inventorying endpoints,
  schema and keys) is a multi-session project. Parked until the v3.0 queue clears and this
  becomes the next milestone target.
- **#1 is never to be raised again.** Giulio's standing instruction, 2026-09-04: it stays open
  permanently and is not a candidate for closing, triage, or a for-giulio entry. Do not
  propose it again.
---

- ~~**Flow v2.0: named multi-instance flows + pipeline config?**~~ — **decided
  2026-07-29: (b) named multi-instance + config.** Design then implement (#40).
  Doc: `docs/handoff/flow-v2-decisions-2026-07-29.md`.
- ~~**Flow v2.0: durable `flows` table vs in-memory only?**~~ — **decided
  2026-07-29: (b) durable records; `auto_restart` default true, optional never;
  process rehydrates on boot; fail loudly on bad pipeline/start.** Filed #41;
  edge-case follow-ups #42. Same design package as #40.
- ~~**Flow v2.0: containers in scope?**~~ — **decided 2026-07-29: yes, phased
  PoC → MVP (#43).** Native Lua/MQTT blocks are *not* a v2.0 gate (containers cover
  custom logic). Doc: `docs/handoff/flow-v2-decisions-2026-07-29.md`.
- ~~`device_deletion_started`/`device_deletion_finished` trigger events are not emitted~~ —
  **decided 2026-07-27: emit both, back-to-back, around the synchronous delete.** Filed as
  issue #21 (`mule`). (Cross-project survey, 2026-07-27,
  `.mule/research/survey-2026-07-27.md` source 4.)
- ~~Mustache trigger-action templates are accepted but not rendered~~ — **decided
  2026-07-27: implement it.** Guiding principle clarified: Astarte compatibility means
  SDK/wire compatibility, not minimum dependency count — Astrate is allowed to be a
  compatible *superset*. Library picked: `github.com/cbroglie/mustache`. Filed as issue #22
  (`mule`). (Same survey, source 4.)
- ~~`value_change`/`value_change_applied`/`path_created`/`path_removed`/`value_stored` trigger
  types compile but never fire~~ — **resolved: implemented in commit 6bd14a7 (2026-08-22)**,
  and the cost question answered by #20's bench on the Legion Go (2026-08-24 closeout):
  ~0.13 ms keyed read per message, −0.7% throughput even paid once per message over REST.
  Performance is not the constraint for these types (nor, a fortiori, for the group-scoped
  line below). (Same survey, source 4.)
- **Group-scoped triggers (`group_name` on device/data triggers) compile but never match**
  (`internal/engine/triggers/match.go:11-12`). Decision deferred, tied to issue #17
  (group-WATCH-path reconciliation, trickle work, not mule): whatever group-membership
  mechanism comes out of that phase should also report the perf cost for this decision —
  noted in a comment on #17 so it isn't benchmarked twice. (Same survey, source 4.)

---

- ~~The Pi cannot run the race detector~~ — **resolved 2026-07-27** by installing Go 1.26.5
  as a userland toolchain on the Legion Go (`~/.local/go`, no root, `rm -rf` to undo). The
  Pi still cannot run `-race` (39-bit VMA kernel vs the 48 ThreadSanitizer needs), so its
  gate remains `go vet ./... && go test ./...` — but race coverage now exists on the Legion
  Go, where the full suite runs clean in ~40s on 16 cores. The standing `race-check` task is
  the concurrency gate. Concurrency work is queueable again, provided the race-check runs
  after it.
- ~~golangci-lint is not installed on the Pi~~ — **resolved 2026-07-28**: installed v2.12.2
  via `go install` (the prebuilt-binary installer's published sha256 for linux/arm64 did not
  verify, reproducibly, so built from source instead). `mule.service` now sets
  `MULE_LINT_CMD=golangci-lint run ./...` and has `/root/go/bin` on `PATH`; the lint gate
  runs starting with the next tick.
- ~~`/root/astrate` on the Pi has uncommitted work~~ — **resolved 2026-07-27** with the new
  `tools/reconcile.sh`: rescued onto `origin/wip/DietPi-20260727T171543Z` (pushed, not
  reviewed — read the diff before merging anything from it) and `/root/astrate` is now clean
  on `main`. Also had to set `commit.gpgsign false` locally on that clone first — it had
  signing on with no working gpg-agent for a non-interactive session, same class of problem
  `/root/astrate-mule` already solved the same way (your call, confirmed 2026-07-27: align it
  with the existing unattended-Pi-clone convention rather than fix gpg-agent or bypass signing
  silently inside the script).

## 2026-08-23 — FDO promoted to milestone-4.0 candidate (Giulio's decision, recorded)

#78 is no longer parked: zero-touch onboarding is strategic for commercial
viability. Scope frozen on the issue: owner-side TO1/TO2 in our Pairing
service only (last mile, like upstream), reuse fdo-rs for
manufacturing/rendezvous, acceptance = official `astarte-device-fdo-rust`
SDK completes onboarding against Astrate, docs as a first-class deliverable.
When v3.0 is marked DONE, the v4.0 section of `.mule/milestones.md` should be
drafted with this investigation as its first item (issue #78 has the full
verified context).

I have a full picture now. Let me write the review file and append the task lines.
← Write .mule/reviews/flowapi-2026-09-10.md
Wrote file successfully.

← Edit .mule/reviews/flowapi-2026-09-10.md
Index: /root/astrate-mule/.mule/reviews/flowapi-2026-09-10.md
===================================================================
--- /root/astrate-mule/.mule/reviews/flowapi-2026-09-10.md
+++ /root/astrate-mule/.mule/reviews/flowapi-2026-09-10.md
@@ -8,9 +8,9 @@
 - `internal/flowapi/http.go` (304 lines, all) — /flow/v1 routing + error mapping.
 - `internal/flowapi/userblocks.go` (209 lines, all) — stored-composite CRUD.
 - Tests: `service_test.go`, `userblocks_test.go`, `http_test.go`,
   `composite_wire_test.go` (container-free), `service_integration_test.go`
-  (`integration` tag, Timescale/DB — Legua-only).
+  (`integration` tag, Timescale/DB — Legion-only).
 - Context: `internal/flow/flow.go` `Manager.StartFlow/Shutdown` (lines 144-183,
   312-329) to verify the shutdown claim; `cmd/astrate/main.go` shutdown path
   (lines 468-491).
 

Now the task lines and the for-giulio escalation:
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -1,6 +1,11 @@
 # Mule queue
 
+- [ ] flowapi-autorestart-terminal-failure: `restartWithBackoff` (internal/flowapi/service.go:585-621) retries forever on a permanent rebuild failure — a flow whose pipeline was deleted (service.go:360-363 wraps it as ErrValidation) or whose config/definition no longer resolves keeps marking itself failed and rebuilding every ≤30s with no exit short of flow deletion. Classify terminal vs transient (ErrValidation → permanent; store/manager errors → transient) via a pure helper and stop the loop on terminal failures; add a container-free unit test on the helper. [auto]
+- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.go:560) with no stop signal, so a block-death racing process shutdown (cmd/astrate/main.go:482-488 drains the Manager but never cancels the goroutine) can rebuild the flow via `mgr.StartFlow` on a background pump (internal/flow/flow.go:173-175) after `Manager.Shutdown` and flip the durable status back to running. Thread a stop channel/context through the Service, checked in the loop's sleep, cancelled on shutdown; verification needs a `[legion]` integration or timing-based test. [auto]
+- [ ] flowapi-validationdetail-test: add a container-free unit test for `validationDetail` (internal/flowapi/http.go:297-304) — the 422-detail prefix-strip rule currently has zero coverage. Assert (1) `"flowapi: validation failed: bad block"` strips to `"bad block"`, (2) a non-prefix message passes through unchanged, (3) a message exactly equal to the prefix is NOT stripped (the guard is `len(msg) > len(prefix)`). Pure helper; no store. [auto]
+- [ ] flowapi-autorestart-default-test: pin the wire rule "POST /flows without `auto_restart` defaults to true" (internal/flowapi/http.go:137-144) — currently untested, so a regression to default-false would pass every existing suite. Extract the default into a pure helper (missing → true, present → honored) and assert all three cases container-free. [auto]
+
 One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).
 The mule runs the topmost open line. Reorder freely; the order is the priority.
 
 If a line needs more detail than fits, write `.mule/tasks/<slug>.md` and keep the line short.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,20 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Flow code review, 2026-09-10: `DeletePipeline` semantics need your call.**
+  `DeletePipeline` (internal/flowapi/service.go:199-206) removes a pipeline
+  silently even when durable flows still reference it; running instances keep
+  the old graph, but their next reload fails ("pipeline not found") and a
+  block death after that starts the auto-restart loop retrying forever against
+  the missing pipeline. `UpdatePipeline` already reports `referencing_flows`
+  (service.go:162-182, issue #44) — should DELETE also refuse (409) while
+  flows reference it, or return the referencing names and allow? Design choice,
+  not a mule task.
+
+---
+
 - **github-issues triage run, 2026-09-10: nothing proposable, nothing stale.** Nine open
   issues: #99–#94 are mule-alarms (idle queue), not code issues. #93 has a pushed commit
   (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
   waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated. #1

→ Read .mule/todo.md [limit=12]
The lines landed above the explanatory block — they belong under the `## Queue` header. Let me move them.
→ Read .mule/todo.md [limit=18, offset=76]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -1,11 +1,6 @@
 # Mule queue
 
-- [ ] flowapi-autorestart-terminal-failure: `restartWithBackoff` (internal/flowapi/service.go:585-621) retries forever on a permanent rebuild failure — a flow whose pipeline was deleted (service.go:360-363 wraps it as ErrValidation) or whose config/definition no longer resolves keeps marking itself failed and rebuilding every ≤30s with no exit short of flow deletion. Classify terminal vs transient (ErrValidation → permanent; store/manager errors → transient) via a pure helper and stop the loop on terminal failures; add a container-free unit test on the helper. [auto]
-- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.go:560) with no stop signal, so a block-death racing process shutdown (cmd/astrate/main.go:482-488 drains the Manager but never cancels the goroutine) can rebuild the flow via `mgr.StartFlow` on a background pump (internal/flow/flow.go:173-175) after `Manager.Shutdown` and flip the durable status back to running. Thread a stop channel/context through the Service, checked in the loop's sleep, cancelled on shutdown; verification needs a `[legion]` integration or timing-based test. [auto]
-- [ ] flowapi-validationdetail-test: add a container-free unit test for `validationDetail` (internal/flowapi/http.go:297-304) — the 422-detail prefix-strip rule currently has zero coverage. Assert (1) `"flowapi: validation failed: bad block"` strips to `"bad block"`, (2) a non-prefix message passes through unchanged, (3) a message exactly equal to the prefix is NOT stripped (the guard is `len(msg) > len(prefix)`). Pure helper; no store. [auto]
-- [ ] flowapi-autorestart-default-test: pin the wire rule "POST /flows without `auto_restart` defaults to true" (internal/flowapi/http.go:137-144) — currently untested, so a regression to default-false would pass every existing suite. Extract the default into a pure helper (missing → true, present → honored) and assert all three cases container-free. [auto]
-
 One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).
 The mule runs the topmost open line. Reorder freely; the order is the priority.
 
 If a line needs more detail than fits, write `.mule/tasks/<slug>.md` and keep the line short.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -72,8 +72,12 @@
     bash tools/mule.sh review
 
 ## Queue
 
+- [ ] flowapi-autorestart-terminal-failure: `restartWithBackoff` (internal/flowapi/service.go:585-621) retries forever on a permanent rebuild failure — a flow whose pipeline was deleted (service.go:360-363 wraps it as ErrValidation) or whose config/definition no longer resolves keeps marking itself failed and rebuilding every ≤30s with no exit short of flow deletion. Classify terminal vs transient (ErrValidation → permanent; store/manager errors → transient) via a pure helper and stop the loop on terminal failures; add a container-free unit test on the helper. [auto]
+- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.go:560) with no stop signal, so a block-death racing process shutdown (cmd/astrate/main.go:482-488 drains the Manager but never cancels the goroutine) can rebuild the flow via `mgr.StartFlow` on a background pump (internal/flow/flow.go:173-175) after `Manager.Shutdown` and flip the durable status back to running. Thread a stop channel/context through the Service, checked in the loop's sleep, cancelled on shutdown; verification needs a `[legion]` integration or timing-based test. [auto]
+- [ ] flowapi-validationdetail-test: add a container-free unit test for `validationDetail` (internal/flowapi/http.go:297-304) — the 422-detail prefix-strip rule currently has zero coverage. Assert (1) `"flowapi: validation failed: bad block"` strips to `"bad block"`, (2) a non-prefix message passes through unchanged, (3) a message exactly equal to the prefix is NOT stripped (the guard is `len(msg) > len(prefix)`). Pure helper; no store. [auto]
+- [ ] flowapi-autorestart-default-test: pin the wire rule "POST /flows without `auto_restart` defaults to true" (internal/flowapi/http.go:137-144) — currently untested, so a regression to default-false would pass every existing suite. Extract the default into a pure helper (missing → true, present → honored) and assert all three cases container-free. [auto]
 - [x] bench-tiers: create the tiered benchmark definitions per .mule/recipes/benchmarks.md (first run only — this task builds the harness, it does not run it)
 - [~] bench-small-astrate: run `bench/scripts/run-tier.sh small astrate` against local Astrate, commit results (two runs minimum) [legion]
 - [~] bench-medium-astrate: run `bench/scripts/run-tier.sh medium astrate` against local Astrate, commit results (two runs minimum) [legion]
 - [~] bench-small-astarte: run `bench/scripts/run-tier.sh small astarte` against local Astarte, commit results (two runs minimum) [legion]

$ cd /root/astrate-mule && sed -n '1,12p' .mule/todo.md && echo "=====QUEUE=====" && sed -n '75,90p' .mule/todo.md
# Mule queue

One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).
The mule runs the topmost open line. Reorder freely; the order is the priority.

If a line needs more detail than fits, write `.mule/tasks/<slug>.md` and keep the line short.

    tools/mule.sh add "<slug>: <outcome>"
    tools/mule.sh loop

A line tagged `[legion]` needs the Legion Go and is skipped automatically while it is asleep;
the queue moves on to the next runnable line rather than stalling behind it.
=====QUEUE=====

- [ ] flowapi-autorestart-terminal-failure: `restartWithBackoff` (internal/flowapi/service.go:585-621) retries forever on a permanent rebuild failure — a flow whose pipeline was deleted (service.go:360-363 wraps it as ErrValidation) or whose config/definition no longer resolves keeps marking itself failed and rebuilding every ≤30s with no exit short of flow deletion. Classify terminal vs transient (ErrValidation → permanent; store/manager errors → transient) via a pure helper and stop the loop on terminal failures; add a container-free unit test on the helper. [auto]
- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.go:560) with no stop signal, so a block-death racing process shutdown (cmd/astrate/main.go:482-488 drains the Manager but never cancels the goroutine) can rebuild the flow via `mgr.StartFlow` on a background pump (internal/flow/flow.go:173-175) after `Manager.Shutdown` and flip the durable status back to running. Thread a stop channel/context through the Service, checked in the loop's sleep, cancelled on shutdown; verification needs a `[legion]` integration or timing-based test. [auto]
- [ ] flowapi-validationdetail-test: add a container-free unit test for `validationDetail` (internal/flowapi/http.go:297-304) — the 422-detail prefix-strip rule currently has zero coverage. Assert (1) `"flowapi: validation failed: bad block"` strips to `"bad block"`, (2) a non-prefix message passes through unchanged, (3) a message exactly equal to the prefix is NOT stripped (the guard is `len(msg) > len(prefix)`). Pure helper; no store. [auto]
- [ ] flowapi-autorestart-default-test: pin the wire rule "POST /flows without `auto_restart` defaults to true" (internal/flowapi/http.go:137-144) — currently untested, so a regression to default-false would pass every existing suite. Extract the default into a pure helper (missing → true, present → honored) and assert all three cases container-free. [auto]
- [x] bench-tiers: create the tiered benchmark definitions per .mule/recipes/benchmarks.md (first run only — this task builds the harness, it does not run it)
- [~] bench-small-astrate: run `bench/scripts/run-tier.sh small astrate` against local Astrate, commit results (two runs minimum) [legion]
- [~] bench-medium-astrate: run `bench/scripts/run-tier.sh medium astrate` against local Astrate, commit results (two runs minimum) [legion]
- [~] bench-small-astarte: run `bench/scripts/run-tier.sh small astarte` against local Astarte, commit results (two runs minimum) [legion]
- [~] bench-medium-astarte: run `bench/scripts/run-tier.sh medium astarte` against local Astarte, commit results (two runs minimum) [legion]
- [~] bench-big-astrate [legion]: run `bench/scripts/run-tier.sh big astrate` against Legion Go Astrate, commit results (two runs minimum)
- [~] bench-giant-astrate [legion]: run `bench/scripts/run-tier.sh giant astrate` against Legion Go Astrate, commit results (two runs minimum)
- [~] bench-big-astarte [legion]: run `bench/scripts/run-tier.sh big astarte` against Legion Go Astarte, commit results (two runs minimum)
- [~] bench-giant-astarte [legion]: run `bench/scripts/run-tier.sh giant astarte` against Legion Go Astarte, commit results (two runs minimum)
- [ ] race-check-store: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/store/... ./internal/housekeeping/... ./migrations/...`. Report any failure to .mule/for-giulio.md with the full race report. Split out of the former single `race-check` line, which timed out running the whole tree at once. [legion] [readonly]
- [ ] race-check-engine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/engine/... ./internal/broker/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]

→ Read .mule/for-giulio.md [limit=16]
Recipe complete. Reviewed `internal/flowapi` (first time — deepest unreviewed package), wrote the review file, appended four task lines, and escalated one design decision.

```
Done: code review of internal/flowapi — 4 proposals: auto-restart retries forever on permanent failure, restart goroutine unreachable on shutdown, validationDetail untested, auto_restart default untested
Files: .mule/reviews/flowapi-2026-09-10.md, .mule/todo.md, .mule/for-giulio.md
Verified: none run — proposal run, no code changed
Unsure: nothing
Follow-ups: flowapi-autorestart-terminal-failure, flowapi-autorestart-shutdown-cancel, flowapi-validationdetail-test, flowapi-autorestart-default-test (all [auto], in .mule/todo.md)
```

Highlights of what the review found:
- `restartWithBackoff` (service.go:585-621) has no terminal condition — a flow whose pipeline was deleted hot-loops every ≤30s forever on a dead config.
- The `go s.restartWithBackoff` (service.go:560) has no stop signal, so a block-death racing process shutdown can rebuild a flow post-`Manager.Shutdown`.
- `validationDetail` (http.go:297-304) and the `auto_restart: true` default (http.go:137-144) are untested wire rules.

`DeletePipeline`-references semantics (refuse vs warn) went to `.mule/for-giulio.md` — that's a product decision, not a mule task.
