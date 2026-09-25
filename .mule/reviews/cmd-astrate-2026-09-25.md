# Code review — cmd/astrate — 2026-09-25

Area: the composition root. Every `internal/*` and `pkg/*` package had been reviewed within
the last 19 days (`.mule/reviews/`, oldest `engine` 2026-09-06), but `cmd/astrate` — the one
package that wires all of them together and owns the shutdown order — had never been the
area of a pass. Seven other reviews cite `cmd/astrate/main.go` in passing; none read it.

Measurement first (no Docker): `go build -tags nats ./cmd/...` and `go vet -tags nats ./cmd/...`
both exit 0 here, so the build-tag split compiles today. `go test -tags nats ./cmd/...` fails
only on `TestForwarderNATSBuilt` (testcontainers), as expected.

## What I read

- `cmd/astrate/main.go` (all 581 lines): flag/version/healthcheck entry (58-88), `run` boot
  order (93-243), `newForwarder` (253-281), `newBroker` (285-315), `selfSignedDevCert`
  (320-345), `mountAPIs` (350-424), `mountRealmVersion` (432-436), `brokerReadiness`
  (440-453), `autoProvisionRealm` (457-473), `shutdown`/`drainEngine` (478-510), `loadSealer`
  (515-526), `runHealthcheck` (531-558), `newLogger` (561-580).
- `cmd/astrate/newnats_default.go`, `newnats_nats.go` — the `-tags nats` split.
- Tests: `forward_test.go` (79 lines, container-free), `version_test.go` (29, container-free),
  `main_test.go` (265, `//go:build integration` — needs the DB, so **no part of `run` is
  covered by a test that runs on this box**), `forward_nats_test.go` (build tag + containers).
- Contracts `main.go` asserts, read to check them: `engine.New`/`Engine.Drain`
  (internal/engine/engine.go:58), `flowapi.Service.RehydrateAutoRestart` (service.go:325-352)
  and `MarkRunningFlowsStopped` (service.go:885-901), `flow.Manager.Shutdown`/`StopFlow`
  (internal/flow/flow.go:312-332, 248-273), `store.ListAutoRestartFlows` (flows.go:116-126),
  `store.Store.Close` (store.go:188), `config.validate` (config.go:311-371), the Makefile
  targets, and the Dockerfile/compose healthcheck wiring.

Read-and-correct, no task: the boot order does honour its comments — the HTTP listener is
only opened at 221-230, after `mountAPIs`, `autoProvisionRealm` (208) and
`RehydrateAutoRestart` (216), so no traffic is accepted before the durable flows are back; the
retention-ceiling sweep (193-206) runs once at boot and then hourly and exits on `ctx.Done`;
`newForwarder`'s disabled branch really is an untyped nil (pinned by `forward_test.go`).

## What I found

### 1. Boot restarts every `auto_restart` flow regardless of its last status — and the shutdown mark that says otherwise is never read

`RehydrateAutoRestart` is called before the listener opens (main.go:216) and restarts whatever
`ListAutoRestartFlows` returns. That query filters on **`f.auto_restart = true` only**
(internal/store/flows.go:125) — it never looks at `f.status`. `auto_restart` is written in
exactly one place, the create path (flowapi/http.go:141-146 → service.go:318 → the INSERT at
flows.go:63); nothing ever clears it, not the stop path, not anything else. So a flow created
with `auto_restart: true` and then stopped through the API comes back at the next boot.

The two halves of `shutdown` disagree about this. main.go:497 calls
`MarkRunningFlowsStopped`, whose entire job is to write `status = "stopped"` for every flow
that was running or stopped (flowapi/service.go:885-901) so the next boot does not resume
them — and the next boot's query ignores `status`, so those writes are dead. Meanwhile the
doc comment on the boot path says the opposite of the shutdown path: "RehydrateAutoRestart
starts **every** durable flow with auto_restart=true" (service.go:325). One of the two is
wrong; both readings are self-consistent, so this needs the intent stated, not guessed.

Operator-visible either way: stop a durable flow, restart astrate, and it is running again.

### 2. One 30s budget is shared by three sequential drain stages; only the engine gets a fresh one

`shutdown` creates a single `sctx` with `shutdownTimeout` (main.go:479) and spends it on
`srv.Shutdown` (484) → `b.Close` (489) → `flowSvc.Manager().Shutdown(sctx)` (494) →
`MarkRunningFlowsStopped(sctx)` (497). `drainEngine` then opens a **second** 30s context
(504). So the constant's own comment — "bounds the whole graceful drain" (main.go:55-56) — is
wrong in both directions: the total can reach ~60s, and conversely a single slow in-flight HTTP
request (a long AppEngine downsample read; the server sets no per-request deadline, 221) burns
the whole first budget and the flow stage runs on an already-expired context. When it does:
`Manager.Shutdown` → `StopFlow` returns at `f.router.Drain(ctx)` (flow.go:265-267) **before**
step 3 releases block resources (270-273), so blocks are never `Stop`ped and in-flight lane
work is dropped, and `MarkRunningFlowsStopped`'s `s.realmID(ctx, …)` fails and `continue`s
past every flow (service.go:895-897). Both failures surface only as one `log.Warn` line
(main.go:495). DESIGN §5.3 says "shards drain (bounded by timeout), batches flush" and gives no
number, so which stage deserves its own budget is a decision, not a derivation.

### 3. `selfSignedDevCert` — the dev-mode broker identity is entirely unpinned

main.go:320-345 mints the throwaway mTLS server identity that `insecure_dev_mode` depends on,
and no test touches it (the only path to it is the integration boot suite). Its rules are all
load-bearing for dev mode and all silently breakable: `NotBefore: now-1h` (331 — `now` would
reject any device whose clock is behind), `NotAfter: now+365d` (332), `DNSNames: localhost`
plus `IPAddresses: 127.0.0.1, ::1` (336-337 — dropping these makes a device that dials
`mqtts://127.0.0.1:8883` fail hostname verification), `ExtKeyUsage: serverAuth` (334), and
self-signed `tmpl, tmpl` (340). Container-free: parse the returned leaf with
`x509.ParseCertificate` and assert the fields.

### 4. `runHealthcheck` — the container self-probe contract has no test

main.go:531-558 is the whole `HEALTHCHECK CMD ["/astrate", "-healthcheck"]` (Dockerfile:35-36)
contract, and it is uncovered: it reads only `ASTRATE_HTTP_ADDR` (532), defaults to `:8080`
→ `127.0.0.1:8080` (534-538), bounds itself at 3s (539), probes `/astrate/v1/readiness` (543) —
not `/health` — and returns 0 only on 200 (554). Container-free with `httptest` plus
`t.Setenv("ASTRATE_HTTP_ADDR", …)`: 200 → 0, 503 → 1, and a server that records the path
proving it is the readiness endpoint.

## What I decided NOT to propose, and why

- **The `-tags nats` half of the tree is compiled by no gate.** `newnats_nats.go`,
  `internal/engine/forward/nats.go` and `forward/nats_test.go` are all `//go:build nats`, and
  no Makefile target, the mule gate, or any test tier passes that tag (Makefile:23, 38, 42, 46)
  — so the entire NATS forwarding feature is invisible to CI, and `config.validate` accepts
  `triggers.forward.kind = "nats"` (config.go:339-343) on a binary that then fails at boot with
  the build-tag message. I measured that it compiles today (`go build -tags nats ./cmd/...` →
  0), so there is no bug to fix — the problem is that nothing keeps it true. Adding the tag to
  the gate is a CI-policy call and the obvious "test" for it (a test that shells out to
  `go build`) is self-referential: the runner strips the implementation and the guard with it,
  so it would never fail. Not a queue line; one line in `.mule/for-giulio.md`.
- **Hijacked WebSocket connections are not drained, and `st.Close()` is unbounded.**
  `srv.Shutdown` does not wait for hijacked conns, and both sockets hijack
  (`websocket.Accept`, appengine/stream/ws.go:87 and the Phoenix equivalent), so a client
  holding a socket at SIGTERM is still inside its handler while the drain runs; the store then
  closes with `pool.Close()` (store.go:188) after "shutdown complete" is already logged
  (main.go:500). The stream handler only reads the bus and writes frames, so the closed pool
  does not actually break it, and the process is exiting anyway — the truncation is one frame.
  Fixing it properly means tracking hijacked conns or giving the sockets a shutdown signal:
  a design decision touching three packages, not a tick.
- **`f.pumpWG.Wait()` (flow.go:262) is not context-bounded**, so a block whose `Process`
  never returns hangs the drain past any budget, inside a goroutine `shutdown` does not
  control. Real, but the fix (pass a context into `Block.Process`, or wait-with-timeout and
  log) changes the block interface — larger than one task, and it belongs with the §5.3
  budget decision in finding 2 rather than beside it. Noted so the next pass does not
  re-derive it.
- **`http.Server` sets only `ReadHeaderTimeout` (221)** — no `IdleTimeout`, so an idle
  keep-alive connection is never reaped, and with no `ReadTimeout` a client can send headers
  then trickle the body forever. Genuine hardening, and `IdleTimeout` is safe to add (it never
  applies to an in-flight request or a hijacked conn), but it has no behavioural test worth
  writing and no wrong behaviour to point at: it is a tuning decision for whoever owns the
  HTTP surface, not a defect.
- **`newLogger`'s level/format switch (561-580)** is untested, but `config.validate` rejects
  any level outside debug|info|warn|error (config.go:379-382) and any format but text/json, so
  the only untested lines are the unreachable defaults. A test would assert what
  `config.validate` already guarantees.
- **The healthcheck is blind to the config file** — it never calls `config.Load`, so an
  operator who sets `http.addr` in the TOML and not in the env gets a healthcheck probing
  127.0.0.1:8080 and a container that flaps unhealthy. Real trap, but the comment at 530
  documents the env-only rule deliberately, the shipped compose deployment does not set
  `http.addr` (so it works), and the fix is either "read the config" (which needs the flag
  plumbed into `runHealthcheck`) or "document it in the compose file" — Giulio's call, not a
  silent patch.
- **`loadSealer` mutates the process environment** (`os.Setenv`, main.go:517) to hand the
  master-key file to the store's env-based loader. Correct and commented; a reviewer of
  `internal/store` already looked at the loader side.
