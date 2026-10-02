# Review — internal/flow/blocks/container, 2026-10-02

Rotation check first: every `internal/` and `pkg/` package has a review file except
`internal/testutil` (test helpers, no logic worth proposing on). The 2026-09-23 flow
pass read `container/block.go` and `container/httpbridge.go` but **not `container/docker.go`**
(211 lines, the whole Docker-CLI surface: `CLIRunner`, `cliInstance`, `Spec`, `Waiter`,
`encodeFlowConfigJSON`) and drew no findings from this package at all. This pass is that
missing half.

## What I read

- `internal/flow/blocks/container/block.go` (all 320) — `Config`, `New` (Start, bridge
  wiring, `WaitReady`, death-watch goroutine), `watchExit`, `parseConfig`, `asInt`,
  `dockerName`, `Process`, `Stop`.
- `internal/flow/blocks/container/docker.go` (all 211) — `Spec`/`Instance`/`Runner`/
  `Waiter`, `CLIRunner.Start`/`run`/`hostPort`, `cliInstance.{BaseURL,ID,Stop,Wait}`,
  `encodeFlowConfigJSON`. **First review of this file.**
- `internal/flow/blocks/container/httpbridge.go` (all 178) — `Bridge`, `WaitReady`,
  `RoundTrip`, the body cap.
- `internal/flow/blocks/container/block_test.go` (all 424) — `fakeInstance`/`fakeRunner`
  and 11 tests.
- Callers and neighbours, only where needed: `internal/flowapi/service.go:380-397`
  (`flowDeps`, `startFlowInstance`), `:529-562` (`onBlockFatal`), `internal/flow/factory.go:100-110`
  (`Instantiate`), `internal/flow/router.go:213-231` (`processOne`), `internal/flow/message.go:142-190`
  (`UnmarshalJSON`), `internal/flow/blocks/info.go:112-117` + `internal/flow/blocks/schema.go:19`
  (the operator-facing catalog entry for `container`),
  `docs/handoff/flow-design-b-container-block-2026-07-29.md` (§6 must/must-not, §6 MVP items,
  §7 security, §9 plan B1–B5).
- Coverage, measured: `go test ./internal/flow/blocks/container/ -coverprofile` → 71.2%
  overall; the per-block `count=0` lines below come from that profile, not from reading.
- Probes: a throwaway program for the duration arithmetic, and a temporary in-package test
  that drove the package's own `CLIRunner.Run` hook to observe the context each `docker`
  invocation receives. Both deleted afterwards; the working tree is unchanged.

## Findings, in value order

### 1. Both orphan-cleanup `docker rm -f` calls run with no deadline

`cliInstance.Stop` injects a 15s bound **only when `ctx == nil`** (docker.go:171-175).
`context.Background()` is not nil, so both real callers defeat that guard:

- docker.go:122 — the mapping-failure cleanup inside `Start`, whose own comment claims
  "Best-effort cleanup so we do not leave orphans on mapping failure";
- block.go:121 — the not-ready cleanup inside `New` (`_ = inst.Stop(context.Background())`).

`Block.Stop` gets it right two hundred lines away, building a 15s timeout explicitly
(block.go:310-311), which is why the guard reads as intended-but-misfired rather than
deliberate.

Measured with the package's `Run` hook: on the hostPort-failure path and on
`inst.Stop(context.Background())`, the `rm` invocation arrives with `ctx.Deadline()`
unset (false in every case; the `run`/`port` calls that *do* inherit the caller's 60s
ctx show the hook is reporting truthfully). What I could **not** measure is the real
hang: this host has no Docker, so "a wedged daemon makes `docker rm -f` block forever
inside `exec.CommandContext`, holding flow instantiation" is reasoning from the
exec/cancel semantics, not a measurement.

Both sites are also invisible to tests: docker.go:120-124 (the cleanup never runs in any
test) and docker.go:171-175 (count=0 in the profile).

### 2. `timeout_ms` / `ready_timeout_ms` are unbounded above and the ms→Duration multiply wraps

`parseConfig` checks only `n > 0` for both (block.go:205-219) and then does
`time.Duration(n) * time.Millisecond` with no bound and no overflow check. Measured:

| config value | duration `parseConfig` computes | what `New` does |
|---|---|---|
| `ready_timeout_ms: 18446744073709` | `-551.616µs` | returns `container "c": container: not ready after wait` **without making one probe** — the wrapped negative deadline fails `time.Now().Before(deadline)` (httpbridge.go:64) so the post-loop error (httpbridge.go:94-97) is what fires |
| `ready_timeout_ms: 9223372036854` | `2562047h47m16s` (~292 y) | **never returns** — still polling `/healthz` when I killed the probe at 3s |

The two consumers of a non-positive duration also disagree: `Bridge.timeout()`
(httpbridge.go:39-44) reads it as "use the 5s default", while `net/http` reads a
non-positive `Client.Timeout` as *no deadline at all* (measured: a client with
`Timeout: -551µs` against a 3-second handler came back `err=<nil>` after 3s). Today
`RoundTrip`'s own `context.WithTimeout(ctx, b.timeout())` (httpbridge.go:118) rescues
the per-message path, so the observable damage is concentrated in the ready-wait: a
fat-fingered value hangs a flow start forever, and a wrapped one fails with a message
that names neither the config key nor the value.

Nothing bounds it on the wire either: `builtinSchemas` advertises both properties as bare
`"type":"integer"` with no `maximum` (internal/flow/blocks/schema.go:19), and
`Info.Config` documents only the defaults (internal/flow/blocks/info.go:116).

### 3. The whole `parseConfig` validation surface except `image` is untested

`block.go:192-193` (`config` must be a JSON object), `:197-202` (the entire `port` 1–65535
rule), `:205-210` and `:215-217` (both timeouts must be positive) are all `count=0` in
the coverage profile. The only config case in the suite is
`TestConstructor_RequiresImage` (block_test.go:99). The port rule is the one worth
noticing: it is what keeps `127.0.0.1::99999` out of the `docker run` argv
(docker.go:93), and deleting it leaves every test green.

### 4. `RoundTrip`'s 1 MiB response cap is never exercised

`io.LimitReader(resp.Body, maxBody+1)` plus the explicit rejection (httpbridge.go:133-140):
lines 136, 138.16 and 138.40 are all uncovered. This is the only bound between a
tenant-supplied container and process memory on the message path, and it is not even
mentioned in the operator-facing catalog entry (info.go:116) or the schema — an author
whose container emits 2 MiB finds out from a runtime error, not from discovery.

### 5. `[legion]` Nothing reaps Astrate-labelled containers left behind by a crash

`CLIRunner.Start` sets `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` labels
(block.go:84-93) precisely so containers are findable, passes no `--rm` and no restart
policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever looks
for them: `rg 'docker (ps|rm|run)'` hits only this package, the example README and the
handoff doc. The boot sweep is a known deferral — handoff §6 MVP item 5 and step **B5
"Optional: orphan cleanup on boot"** (docs/handoff/flow-design-b-container-block-2026-07-29.md:199,239)
— so this is not an oversight. The failure mode is a crash loop though: every astrate
restart leaves the previous run's containers up, holding loopback-published ports and
RAM, with no code path that ever notices.

## What I decided *not* to propose, and why

- **`New` drops the caller's context entirely** (block.go:98 and 118 build both the 60s
  start context and the ready context from `context.Background()`), so a cancelled or
  timed-out flow start keeps its container and keeps polling `/healthz`; nothing outside
  the package can cancel it. This is arguably the strongest thing I found and I am still
  not proposing it: closing it means threading a context through `Registry.Instantiate`
  (internal/flow/factory.go:105), its production caller `resolveAndBuild`
  (internal/flowapi/service.go:380) and every constructor in `internal/flow/blocks/` —
  far past the two-file rule. Recorded so the next run does not re-derive it; finding 2's
  bound removes the worst symptom of it.
- **`Process` logs Error for every failed message** (block.go:281) while `processOne`
  logs the same error and increments `blockErrors` (internal/flow/router.go:226-229), so
  a container that starts 500ing emits two ERROR lines per message indefinitely. Not a
  behaviour bug but a log-level judgement call, and pinning it needs either logger
  injection into `Block` (`log` is hard-wired to `slog.Default()` at block.go:128) or a
  test that swaps the global default logger.
- **`hostPort` returns whatever follows the last colon on the first non-empty line**
  (docker.go:143-155), with no numeric check and no loopback assertion. Not proposed:
  the publish flag is hard-coded `127.0.0.1::port` (docker.go:93), so the daemon cannot
  answer with a routable address, and a non-numeric port produces a URL that
  `http.NewRequestWithContext` rejects. The blank-line/multi-line branches are untested,
  but a test would only pin docker's output format.
- **`dockerName`'s uniqueness suffix** (block.go:261, `time.Now().UnixNano()%1e9`) —
  two calls can only collide inside the same nanosecond, and duplicate block names in one
  pipeline are rejected before this runs. Nothing to fix.
- **`Run`'s unbounded `bytes.Buffer` for stdout/stderr** (docker.go:73-75) and the
  `-e ASTRATE_FLOW_CONFIG=<json>` argv (docker.go:99, E2BIG on a very large nested
  config) — both inputs Astrate itself produces, both with an existing error path, neither
  attacker-controlled.
- **`Process`'s `bridge == nil` guard** (block.go:274) is dead: `bridge` is set in `New`
  and never cleared (`Stop` nils `inst`, not `bridge`), and `stopped` is always set first.
  One-line clarity, worth folding into whatever else touches that function.
- **No metrics, and `Bridge.MaxBodyBytes` is unreachable from config** — both are MVP
  items in the handoff (§6 items 3 and 4, docs/handoff/flow-design-b-container-block-2026-07-29.md:198-199)
  and both are feature-shaped (metric names and a dashboard decision; a new config key).
  They went to `.mule/for-giulio.md` instead of the queue. The design handoff's §7
  "Container runs with Docker defaults ... document that operators trust the image" is
  satisfied by the current `Info.Summary`, so the missing image allowlist is documented
  behaviour rather than a gap — not proposed.

## Task lines appended

`container-stop-deadline`, `container-timeout-bounds`, `container-parseconfig-rules-test`,
`container-response-cap-test`, `container-orphan-boot-cleanup`.
