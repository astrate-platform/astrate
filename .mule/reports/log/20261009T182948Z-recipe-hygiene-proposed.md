slug: recipe-hygiene
verdict: proposed
at:  67809c0
ran: 2026-10-09T18:29:48Z on DietPi in 811s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md

$ which govulncheck golangci-lint 2>/dev/null; ls /Users/atsetilam/go/bin/golangci-lint 2>/dev/null; echo "---"; go version
---
go version go1.26.7 linux/arm64

$ go list -m -u all 2>/dev/null | rg '\[' | head -20; echo "---EXIT $?---"
cel.dev/expr v0.24.0 [v0.25.3]
cloud.google.com/go v0.121.6 [v0.123.0]
cloud.google.com/go/auth v0.16.4 [v0.24.1]
cloud.google.com/go/auth/oauth2adapt v0.2.8 [v0.3.0]
cloud.google.com/go/compute/metadata v0.8.0 [v0.10.0]
cloud.google.com/go/iam v1.5.2 [v1.14.0]
cloud.google.com/go/longrunning v0.6.7 [v1.3.0]
cloud.google.com/go/monitoring v1.24.2 [v1.31.0]
cloud.google.com/go/spanner v1.85.0 [v1.96.1]
cloud.google.com/go/storage v1.56.0 [v1.69.0]
github.com/99designs/keyring v1.2.1 [v1.2.2]
github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0 [v1.23.3]
github.com/Azure/azure-sdk-for-go/sdk/internal v1.1.2 [v1.13.0]
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0 [v1.8.2]
github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c [v0.0.0-20261008234032-65faa4be4f89]
github.com/Azure/go-autorest/autorest/adal v0.9.16 (retracted) [v0.9.24] (deprecated)
github.com/Azure/go-autorest/autorest/date v0.3.0 [v0.3.1] (deprecated)
github.com/Azure/go-autorest/logger v0.2.1 [v0.2.2] (deprecated)
github.com/Azure/go-autorest/tracing v0.6.0 [v0.6.1] (deprecated)
github.com/ClickHouse/clickhouse-go v1.4.3 [v1.5.4]
---EXIT 0---

$ govulncheck ./... 2>/dev/null || echo "govulncheck not installed"
govulncheck not installed

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ 2>/dev/null | head -40; echo "---EXIT $?---"
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},
---EXIT 0---

→ Read internal/broker/intake.go [offset=40, limit=40]

$ go test ./... 2>&1 | rg -i 'skip|no test files' | head -30; echo "---EXIT ${PIPESTATUS[0]}---"
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]
/bin/sh: 1: Bad substitution

$ rg -n 't\.Skip|testing\.Short|SkipNow' --glob '*_test.go' internal/ pkg/ cmd/ 2>/dev/null; echo "---EXIT $?---"
cmd/astrate/main_test.go:47:	t.Skip("the boot suite needs a database: run `make up` or set ASTRATE_TEST_DSN")
cmd/astrate/forward_test.go:68:		t.Skip("built with -tags nats; see forward_nats_test.go")
internal/pairing/http_test.go:629:		t.Skip("openssl not found in PATH; chain already verified with crypto/x509")
internal/pairing/ca/ca_test.go:292:	if testing.Short() {
internal/pairing/ca/ca_test.go:293:		t.Skip("10k issuance draw skipped in -short mode")
internal/store/datastreams_test.go:420:			t.Skip("timescaledb_toolkit not installed; lttb path unavailable")
---EXIT 0---

$ for d in $(go list ./... 2>/dev/null | rg 'astrate/(internal|pkg|cmd)'); do dir=$(go list -f '{{.Dir}}' "$d" 2>/dev/null); if ! ls "$dir"/*_test.go >/dev/null 2>&1; then echo "NO_TEST: $d"; fi; done
(no output)

$ ls -la .golangci* 2>/dev/null; echo "---"; rg -n 'golangci|govulncheck|traceguard' .github/ Makefile 2>/dev/null | head -30; echo "---GOVULN TRY---"; timeout 90 go run golang.org/x/vuln/cmd/govulncheck@latest ./... 2>&1 | tail -40
-rw-r--r-- 1 root root 1067 Jul 27 05:10 .golangci.yml
---
Makefile:4:#   make lint                golangci-lint (config: .golangci.yml)
Makefile:30:## lint: run the pinned golangci-lint (install with `make tools`).
Makefile:32:	@command -v golangci-lint >/dev/null 2>&1 || { \
Makefile:33:		echo "golangci-lint not found; run 'make tools' (installs $(GOLANGCI_LINT_VERSION))"; exit 1; }
Makefile:34:	golangci-lint run ./...
Makefile:58:	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
.github/workflows/ci.yml:30:      - uses: golangci/golangci-lint-action@v8
---GOVULN TRY---
signal: killed

$ free -h; nproc; echo "---"; sed -n '25,62p' Makefile; echo "---"; cat .golangci.yml
total        used        free      shared  buff/cache   available
Mem:           3.7Gi       2.3Gi       1.1Gi       1.2Gi       1.6Gi       1.5Gi
Swap:             0B          0B          0B
4
---
		mkdir -p $(DIST); \
		CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(DIST)/astrate ./cmd/astrate; \
		echo "built $(DIST)/astrate"; \
	fi

## lint: run the pinned golangci-lint (install with `make tools`).
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not found; run 'make tools' (installs $(GOLANGCI_LINT_VERSION))"; exit 1; }
	golangci-lint run ./...

## test: T1 — pure unit tests, race detector on.
test:
	$(GO) test -race ./...

## test-integration: T2 — requires a Docker daemon (or ASTRATE_TEST_DSN to reuse a database).
test-integration:
	$(GO) test -race -count=1 -tags integration ./...

## test-e2e: T3 — component tests behind the e2e build tag (suites land from M5 onward).
test-e2e:
	$(GO) test -race -count=1 -tags "integration e2e" ./...

## test-conformance: T4 — official-SDK harness; lives in test/ with its own go.mod (M9).
test-conformance:
	@if [ -f test/conformance/go.mod ]; then \
		cd test/conformance && $(GO) test -count=1 ./...; \
	else \
		echo "conformance harness lands in M9 (docs/ROADMAP.md §10); nothing to run yet"; \
	fi

## tools: install the pinned developer toolchain.
tools:
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

## up/down: local TimescaleDB for development and DSN-reuse test runs.
up:
	docker compose up -d
---
# golangci-lint v2 configuration (pinned v2.12.2 — see Makefile GOLANGCI_LINT_VERSION).
# Linter set per docs/ROADMAP.md §1.2 (file 0.4): govet, staticcheck, errcheck, gosec, revive.
version: "2"

run:
  build-tags:
    # Lint everything, including T2/T3 suites that are tag-gated out of `go test ./...`.
    - integration
    - e2e

linters:
  default: none
  enable:
    - errcheck
    - gosec
    - govet
    - revive
    - staticcheck
  settings:
    gosec:
      excludes:
        # G104 duplicates errcheck.
        - G104
  exclusions:
    generated: lax
    rules:
      # Tests: hardening lints and strict error checking add noise, not safety.
      - path: _test\.go
        linters:
          - gosec
      - path: _test\.go
        linters:
          - errcheck
      # Test utilities are test-scoped code shipped outside _test.go files.
      - path: internal/testutil/
        linters:
          - gosec

formatters:
  enable:
    - gofmt
    - goimports
  settings:
    goimports:
      local-prefixes:
        - github.com/astrate-platform/astrate

$ ls -la .mule/; echo "---TAGS---"; git tag --list 2>/dev/null | tail -20; echo "---MILESTONES---"; sed -n '1,60p' .mule/milestones.md 2>/dev/null
total 16296
drwxr-xr-x  7 root root     4096 Oct  9 20:07 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       13 Oct  9 20:16 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  9 19:51 .heartbeat
-rw-r--r--  1 root root     4096 Oct  9 20:18 .last-output
-rw-r--r--  1 root root        0 Oct  9 20:16 .lock
-rw-r--r--  1 root root        8 Oct  9 20:16 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 16238012 Oct  9 20:18 cron.log
drwxr-xr-x  2 root root     4096 Oct  9 20:07 failed
-rw-r--r--  1 root root   161185 Oct  9 13:28 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    16266 Oct  9 20:07 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  8 19:28 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   175392 Oct  9 20:07 todo.md
---TAGS---
mule-queue-archive-2026-09-04
v0.1.0
v0.2.0
---MILESTONES---
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

Status: **DONE** (2026-07-29), tagged `v0.2.0` on 2026-09-04 — the milestone names
(`v2.0`, `v3.0`) are milestone names, not release versions; the project is pre-1.0 and the
version number keeps its own line. Runtime, factory, catalog (incl. filter/map), process
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

$ rg -n 'APICompatVersion' --glob '*.go' . | head; echo "---FORGIULIO---"; cat .mule/for-giulio.md 2>/dev/null | head -60
./cmd/astrate/main.go:437:		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(realm.APICompatVersion)))
./cmd/astrate/main.go:438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
./cmd/astrate/version_test.go:25:	want := `{"data":"` + realm.APICompatVersion + `"}`
./internal/realm/http.go:119:// feature UI on it — see APICompatVersion).
./internal/realm/http.go:121:	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
./internal/realm/dashboard_compat_test.go:35:		if v != APICompatVersion {
./internal/realm/dashboard_compat_test.go:36:			t.Errorf("version = %q, want %q", v, APICompatVersion)
./internal/realm/service.go:580:// APICompatVersion is the upstream Realm Management API level Astrate
./internal/realm/service.go:588:const APICompatVersion = "1.2.2"
---FORGIULIO---
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **docs-sync run, 2026-10-09, surface: appengine — `docs/site/appengine-api.md` mislabels the `downsample_to` datastream query parameter as a duration.** The datastream-query snippet writes it as `&downsample_to=<bucket_duration>` (line 60) and the note below says "`downsample_to` maps onto Timescale `time_bucket()`" (line 65), but the code reads it as a point COUNT, not a time span: `parseQueryOpts` parses an integer and requires it `> 2`, storing it in `opts.DownsamplePoints` (internal/appengine/http.go:588-598), and the bucket interval is computed *from* that count downstream (`bucketFor(last.Sub(first), opts.DownsamplePoints)` then `s.Downsample(ctx, q, bucket)`, internal/appengine/data.go:187-196) — which is what the OpenAPI already says ("Downsample to this many data points (must be > 2)", docs/api/astarte_appengine_api.yaml:1509-1514). A reader following the page sends a duration like `1h` and gets a 422 (`downsample_to is invalid`, http.go:591-592). Proposed wording (your voice): `&downsample_to=<max_points>` and "`downsample_to` is the target number of points (must be greater than 2); the server chooses the Timescale `time_bucket()` interval." Page untouched.

- **github-issues triage run, 2026-10-08 (evening, 20:37Z): nothing proposable for the 27th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub.** The recipe's `--limit 40` command prints all **32** open issues: still **28** `mule-alarm` **#94–#121** (no #122 — nothing new filed since this morning's run) plus the same four non-alarm, none proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), one excluded by standing instruction. **The alarm latch is working again and #121's window has ended**: `.mule/.alarmed` is gone, `.mule/.heartbeat` = `1791488167` = **2026-10-08T19:36:07Z**, and `.mule/log.md` now carries four `2026-10-08` `done` rows the morning run could not see (`d070345`, `cfa6e76`, `ed301fb`, `1cda0fb`, plus one 1200s timeout) — the very land that clears the latch (mule.sh:578→737), so #121 needs no action either. **Proposal unchanged from earlier today: close #94–#120, all twenty-seven** — each is a one-day alarm superseded by the next and nothing reads issue state; leave #121 to self-expire. Re-verified today, not inherited: **#92's parking condition still unmet** (`gh api .../tags` → newest `v1.4.0-rc.6`, `/releases` → newest stable `v1.3.5`; an `-rc.N` does not satisfy "wait for a stable v1.4.0"), and **#93's rewrite still entirely off main** — `git branch -a --contains 24ad5b8` lists only `mule/queue` / `origin/mule/queue`, `origin/main...HEAD` = 4 / **732** (4 / 711 this morning).

---

- **github-issues triage run, 2026-10-08: nothing proposable for the 26th run — no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub.** The recipe's `--limit 40` command prints all **32** open issues (fits the limit, no silent drop): **28** `mule-alarm` **#94–#121** plus the same four non-alarm, none of them proposable — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (`milestone-4.0`, already escalated), and one the standing instruction excludes from triage entirely. New since the 2026-10-07 run: **#121**, created 2026-10-08T11:02:45Z, "nothing has landed in 16h" — arithmetic pinned, not inferred: `.mule/.heartbeat` holds `1791398840` = **2026-10-07T18:47:20Z**, written 3s after the last `done` land (`322242a` pairing-unregister-bad-id-test, 2026-10-07T20:47:17+02:00), and 11:02:45Z − 18:47:20Z = **16h15m**, so `age=$(( (now-last)/3600 ))` prints 16 — the title verbatim; `.mule/.alarmed` carries mtime **11:03:05Z**, i.e. set by that filing. Genuine window, the same overnight band #117–#120 measured, and it is still open (heartbeat unmoved at 11:09Z, six minutes on). **Proposal, extending the 2026-10-07 one: close #94–#120, all twenty-seven** — each is a one-day alarm superseded by the next and nothing reads issue state; leave **#121** (today's, real window) to self-expire.
  **New this run: nothing has landed *today* at all, and the latch is already set.** `.mule/log.md` has **zero** `2026-10-08` rows (6 / 6 / 3 rows on 10-05 / 10-06 / 10-07) and HEAD's newest commit is `0d2e197` "recipe milestones ran" at 2026-10-07T22:44:44Z — **~12h before the alarm**, so no commit at all since. The timer itself is alive (it filed #121 at 11:02:45Z), so this is the same shape the 2026-10-07 entry measured, one day further: `.mule/.alarmed` short-circuits `check_pulse` at `tools/mule.sh:748` (set at :749 by this filing), and only a `done` land (mule.sh:578) or a `checked` run (mule.sh:540) calls `beat` (mule.sh:737), which is the only thing that clears it — so the dead-man's switch is now off through today's silence, exactly as it was through yesterday's. The `check_pulse`-before-`beat` ordering fix (move the call at `tools/mule.sh:633` to after the tick's `beat`) remains escalated above, still yours, still unactioned, and #121 is the same arithmetic as #120 — not fresh evidence for it either way. Queue context: `tools/mule.sh status` reports **14 open / 37 blocked**, so the next `done` (the only un-latch) may be a while coming.
  **#93 at 34 days.** `git rev-list --left-right --count origin/main...HEAD` = **4 / 711** today (4 / 697 on 10-07), and `git merge-base --is-ancestor 24ad5b8 origin/main` still fails — `git branch -a --contains 24ad5b8` lists only `mule/queue` / `origin/mule/queue`, so the rewrite is entirely off main. The version-drift detail recorded on 2026-10-06/07 is unchanged and re-verified today (`gh api repos/astarte-platform/astarte/tags` → `v1.4.0-rc.6`; `/releases` → `v1.3.5` stable 2026-10-05): the in-tree comment at the `control/keyAgreement` deny (internal/broker/aclhook.go) still says "v1.4.0-rc.5 … v1.3.3 being the newest stable tag". It folds into that one review, not a line of its own. **#92's parking condition re-verified unmet** by the same sweep: newest stable is `v1.3.5`, newest overall is the prerelease `v1.4.0-rc.6` — an `-rc.N` does not satisfy your "wait for a stable v1.4.0".

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
  **The parked-decisions row for #92 is still unmet, re-checked against the tag list rather than the release list** (this host has no `.mule/waiting-on.md` — absent from `HEAD`, already reported; I read it out of `origin/main` read-only): `gh api repos/astarte-platform/astarte/tags?per_page=100` gives newest stable **v1.3.5** and newest overall **v1.4.0-rc.6**, with `v1.4.0-rc.0`…`rc.6` and no stable `v1.4.0`. The row asks for a **stable** tag, which an `-rc.N` does not satisfy, so the twenty-second run ends the same way as the twenty-one before it.
  **The one delta to your open COMPATIBILITY.md wording proposal (the v1.3.2 one, folding v1.3.3/v1.3.4, further down this file): its version-reference sentence should read "newest stable v1.3.5" where it currently says "v1.3.4". Nothing else in it changes** — v1.3.5 adds no capability, so the same decision and the same "not yet emulated" list stand. Not applied: `docs/COMPATIBILITY.md` is untouched and still targets **1.2.2**, `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). I deliberately did **not** queue a fresh `compat-note-v1.3.5` line: the proposal is already open here, and a second line would re-derive it for the next run to find.

- **github-issues triage run, 2026-10-05: nothing proposable for the 23rd run, and the alarm pile's arithmetic is now fully pinned — 27 gaps over 8h since 2026-08-20 and **not one of them between 8h and 13h**, so the threshold cannot be a number. Two things are new: today's #118 was silenced by a land 5 minutes later, and #93's commit has now been unreviewed for 31 days.** 29 open issues (authoritative `--limit 100`: 25 alarms **#94–#118**, one new today, plus exactly four non-alarm — **#93** aclhook comment `mule-review`, **#92** keyAgreement parked, **#78** FDO `milestone-4.0`, **#1** untouched per standing instruction). Note the recipe's own `--limit 40` command printed **all 29** today, so the "silently dropped #94" parenthetical in the 2026-10-04 entry did not reproduce — treat the limit as adequate at this pile size. Zero machine-checkable candidates, so **no `.mule/todo.md` lines added, no `gh issue create`, and nothing commented, closed or edited on GitHub**.
  **My recount reproduces yesterday's numbers rather than replacing them, which is the useful part: the distribution is bimodal with an empty middle.** Over all **139** `mule: log` lands since 2026-08-20 (138 gaps): median **0.8h**, and the counts `> 8h` = `> 12h` = `> 13h` = **27** — so no gap anywhere in the record falls in (8h, 13h). Each of the 25 open alarms therefore measured a real 14.1–16.2h window (the only outlier is **#114** at 8.6h, filed 2026-09-27T20:09Z off the 11:35Z land that same day, which is the bimodality showing its other edge: a 20:00Z check before the evening batch starts). The 2026-10-04 entry's conclusion holds and hardens — a threshold in [8h, 13h) fires on every healthy day by construction; the real stall in the record (2026-09-28 → 10-01) is **74.7h**. **The proposed fix is unchanged and remains yours: move the `check_pulse` call from `tools/mule.sh:633` to after the tick's `beat`, so the switch cannot file against the work it is about to watch land.** Today's instance is the cleanest yet: **#118 filed 10:54:06Z** measuring a genuine 14.8h back to the 2026-10-04T20:06Z land, and the next land, `store-validatepipelinegraph-error-branches`, committed **10:59:49Z — 5 minutes later** (#117 was 23). Six lands today (10:59Z → 19:34Z), five `done` in `.mule/log.md`.
  **The 2026-10-01 entry's "latched off for three days" diagnosis is now superseded, and `.alarmed` explains how: `beat()` at `tools/mule.sh:737` does `rm -f "$MULE/.alarmed"` on every land, so one land per day is enough to re-arm the alarm every single morning.** `.mule/.alarmed` does not exist now, the switch re-armed on **2026-10-02**, and it has filed four more (#115–#118) since. The queue has not been idle for a month — my `git log --since=2026-08-20` count is 139 lands. Do not act on the September idle-streak framing any more; it described four real days and then the branch simply started moving again.
  **Proposal, extending the 2026-09-25, 09-26, 10-01, 10-03 and 10-04 ones: close #94–#117, all twenty-four.** Each is a one-day alarm superseded by the next, and as established on 2026-09-25 nothing reads issue state, so closing stays inert — but 24 open "the mule is dead" issues whose bodies point at a timer that does not exist is the noise that buried a real three-day stall in September. Leave **#118** (today's) open: like #117 it is the one whose window was unambiguously real.
  **New, and it is the same root cause as yesterday's merge finding: the review queue is not draining, so #93 has been waiting 31 days.** `8c61268` ("mule: log issue-93", aclhook comment rewrite) landed **2026-09-04T21:55:29+02:00**, is still **not** an ancestor of `origin/main`, and `git branch -a --contains 8c61268` lists only `mule/queue` / `origin/mule/queue` — a month of `mule-review` with no reviewer, while `bash tools/mule.sh review` sits unread. Divergence is now **4 / 670** (`git rev-list --left-right --count origin/main...HEAD`), the same 4 unpushed-here commits as yesterday (`f1d0069`, `ca47b35`, `99743ca`, `2a033d4`, all 2026-09-04 21:40–21:52) with `git merge-base --is-ancestor f1d0069 HEAD` → false. Yesterday measured this at 4 / 643: **this branch moved 27 commits in 24h and main's four did not move at all**, i.e. the review step has now missed a full day of landings on top of the month it has already missed. That single merge is what retires the pile's misleading bodies *and* unblocks #93; I have not touched git and will not.

$ ls -la "$(go env GOPATH)/bin" 2>/dev/null; echo "---GOBIN $(go env GOBIN)---"; ls /root/go/bin 2>/dev/null
total 85260
drwxr-xr-x 2 root root     4096 Sep 15 20:20 .
drwxr-xr-x 4 root root     4096 Jul 28 23:52 ..
-rwxr-xr-x 1 root root  7668042 Sep  2 21:44 goimports
-rwxr-xr-x 1 root root 65622584 Jul 28 23:52 golangci-lint
-rwxr-xr-x 1 root root 13997377 Sep 15 20:20 govulncheck
---GOBIN ---
goimports
golangci-lint
govulncheck

$ go test ./... 2>&1 | rg -v '^ok |no test files' | head -40; echo "===DONE==="
===DONE===

$ /root/go/bin/golangci-lint version 2>&1 | head -3; echo "---RUN---"; timeout 600 /root/go/bin/golangci-lint run ./... 2>&1 | tail -60
golangci-lint has version 2.12.2 built with go1.26.5 from (unknown, modified: ?, mod sum: "h1:7+d1uY0bq1MU2UV3R5pW5Q7QWdcoq4naMRXM+gsJKrs=") on (unknown)
---RUN---
0 issues.

$ /root/go/bin/govulncheck -version 2>&1 | head -5; echo "---"; timeout 600 /root/go/bin/govulncheck ./... 2>&1 | tail -60
Go: go1.26.7
Scanner: govulncheck@v1.8.0
DB: https://vuln.go.dev
DB updated: 2026-10-08 22:31:09 +0000 UTC

---

! permission requested: external_directory (/tmp/*); auto-rejecting
✗ /root/go/bin/govulncheck ./... > /tmp/govuln.out 2>/tmp/govuln.err; echo "EXIT=$?"; echo "---STDOUT---"; cat /tmp/govuln.out; echo "---STDERR---"; cat /tmp/govuln.err failed
Error: The user rejected permission to use this specific tool call.
