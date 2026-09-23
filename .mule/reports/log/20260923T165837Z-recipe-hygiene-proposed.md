slug: recipe-hygiene
verdict: proposed
at:  b8fbc71
ran: 2026-09-23T16:58:37Z on DietPi in 619s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md

$ go version; which go govulncheck golangci-lint 2>/dev/null; git -C /root/astrate-mule status --short --branch 2>/dev/null | head; uname -a
go version go1.26.7 linux/arm64
/usr/bin/go
## mule/queue
Linux DietPi 6.18.39+rpt-rpi-v8 #1 SMP PREEMPT Debian 1:6.18.39-1+rpt1 (2026-07-29) aarch64 GNU/Linux

→ Read .mule/todo.md
This is the Pi (mule/queue branch). Running the hygiene recipe checks now.

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ | head -40
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},

$ go list -m -u all 2>/dev/null | rg '\[' | head -20; echo "---govulncheck---"; govulncheck ./... 2>/dev/null || echo "govulncheck not installed"; echo "---golangci---"; ls /Users/atsetilam/go/bin/golangci-lint 2>/dev/null; command -v golangci-lint 2>/dev/null || echo "golangci-lint not on PATH"
cel.dev/expr v0.24.0 [v0.25.3]
cloud.google.com/go v0.121.6 [v0.123.0]
cloud.google.com/go/auth v0.16.4 [v0.23.3]
cloud.google.com/go/compute/metadata v0.8.0 [v0.9.1]
cloud.google.com/go/iam v1.5.2 [v1.13.0]
cloud.google.com/go/longrunning v0.6.7 [v1.2.0]
cloud.google.com/go/monitoring v1.24.2 [v1.30.0]
cloud.google.com/go/spanner v1.85.0 [v1.95.1]
cloud.google.com/go/storage v1.56.0 [v1.68.0]
github.com/99designs/keyring v1.2.1 [v1.2.2]
github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0 [v1.23.1]
github.com/Azure/azure-sdk-for-go/sdk/internal v1.1.2 [v1.12.0]
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0 [v1.8.1]
github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c [v0.0.0-20260917205352-e937bb47801a]
github.com/Azure/go-autorest/autorest/adal v0.9.16 (retracted) [v0.9.24] (deprecated)
github.com/Azure/go-autorest/autorest/date v0.3.0 [v0.3.1] (deprecated)
github.com/Azure/go-autorest/logger v0.2.1 [v0.2.2] (deprecated)
github.com/Azure/go-autorest/tracing v0.6.0 [v0.6.1] (deprecated)
github.com/ClickHouse/clickhouse-go v1.4.3 [v1.5.4]
github.com/DataDog/zstd v1.4.5 [v1.5.7]
---govulncheck---
govulncheck not installed
---golangci---
golangci-lint not on PATH

$ ls .golangci* .golangci.yml 2>/dev/null; ls tools/; rg -n 'golangci|lint' tools/mule.sh Makefile 2>/dev/null | head -20
.golangci.yml
.golangci.yml
bsoncapture
mule-plan-day.sh
mule-setup-pi.sh
mule-survey.sh
mule.sh
reconcile.sh
tools/mule.sh:213:# over five weeks by ten pre-existing lint errors in the checkout's own base — the gate lints
tools/mule.sh:221:# errors, test failures and golangci-lint findings all carry a file:line prefix.
tools/mule.sh:244:    note "gate: lint"
tools/mule.sh:246:      || { gate_reason "lint failed"; bad "$GATE_REASON"; rc=1; }
tools/mule.sh:306:      && ok "lint baseline is clean" \
tools/mule.sh:307:      || bad "lint baseline is not clean (clear the linter cache before believing this)"
tools/mule.sh:886:   6  hygiene    dependency, CI and lint upkeep
Makefile:4:#   make lint                golangci-lint (config: .golangci.yml)
Makefile:19:.PHONY: build lint test test-integration test-e2e test-conformance tools up down clean
Makefile:30:## lint: run the pinned golangci-lint (install with `make tools`).
Makefile:31:lint:
Makefile:32:	@command -v golangci-lint >/dev/null 2>&1 || { \
Makefile:33:		echo "golangci-lint not found; run 'make tools' (installs $(GOLANGCI_LINT_VERSION))"; exit 1; }
Makefile:34:	golangci-lint run ./...
Makefile:58:	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

$ go test ./... 2>&1 | rg -i 'skip|no test files' | head -30
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]

$ sed -n '235,255p' tools/mule.sh; echo '---'; rg -n 'GOLANGCI_LINT_VERSION|golangci-lint' Makefile
: > "$GATE_LOG"
  GATE_REASON=""
  if [ -n "${MULE_FIX_CMD:-}" ]; then ( cd "$REPO" && eval "$MULE_FIX_CMD" ) >/dev/null 2>&1; fi
  if [ -n "${MULE_TEST_CMD:-}" ]; then
    note "gate: $MULE_TEST_CMD"
    ( cd "$REPO" && eval "$MULE_TEST_CMD" ) >>"$GATE_LOG" 2>&1 \
      || { gate_reason "tests failed"; bad "$GATE_REASON"; rc=1; }
  fi
  if [ "$rc" = 0 ] && [ -n "${MULE_LINT_CMD:-}" ]; then
    note "gate: lint"
    ( cd "$REPO" && eval "$MULE_LINT_CMD" ) >>"$GATE_LOG" 2>&1 \
      || { gate_reason "lint failed"; bad "$GATE_REASON"; rc=1; }
  fi
  if [ "$rc" = 0 ] && [ -z "${MULE_NO_PROOF_GATE:-}" ]; then
    proof_gate || rc=1
    # proof_gate marks untracked files intent-to-add so diff can see them, which leaves index
    # entries behind on every exit path. Nothing is deliberately staged during a task — the
    # runner stages after the gates — so clearing the index here is safe, and not clearing it
    # leaves a phantom deletion that the revert path would act on.
    git -C "$REPO" reset -q >/dev/null 2>&1 || true
  fi
---
4:#   make lint                golangci-lint (config: .golangci.yml)
12:GOLANGCI_LINT_VERSION  ?= v2.12.2
30:## lint: run the pinned golangci-lint (install with `make tools`).
32:	@command -v golangci-lint >/dev/null 2>&1 || { \
33:		echo "golangci-lint not found; run 'make tools' (installs $(GOLANGCI_LINT_VERSION))"; exit 1; }
34:	golangci-lint run ./...
58:	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

$ ls migrations/ examples/flow-container-echo/ docs/ 2>/dev/null; echo '---migrations tests---'; rg -l '_test.go' migrations/ examples/ 2>/dev/null; echo '---'; ls migrations/*_test.go 2>/dev/null
docs/:
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
UPSTREAM-EXPERIMENTAL.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui

examples/flow-container-echo/:
Dockerfile
README.md
main.go

migrations/:
000001_extensions.down.sql
000001_extensions.up.sql
000002_metadata.down.sql
000002_metadata.up.sql
000003_properties.down.sql
000003_properties.up.sql
000004_datastreams.down.sql
000004_datastreams.up.sql
000005_jobs.down.sql
000005_jobs.up.sql
000006_policies.down.sql
000006_policies.up.sql
000007_toolkit.down.sql
000007_toolkit.up.sql
000008_pipelines.down.sql
000008_pipelines.up.sql
000009_flows.down.sql
000009_flows.up.sql
000010_flow_failed_block.down.sql
000010_flow_failed_block.up.sql
000011_realm_retention.down.sql
000011_realm_retention.up.sql
000012_user_blocks.down.sql
000012_user_blocks.up.sql
migrations.go
---migrations tests---
---

$ ls ~/go/bin /root/go/bin /usr/local/bin 2>/dev/null; echo '---'; rg -n 'golangci|MULE_LINT_CMD|lint' tools/mule-setup-pi.sh tools/mule.sh | head -25
/root/go/bin:
goimports
golangci-lint
govulncheck

/root/go/bin:
goimports
golangci-lint
govulncheck

/usr/local/bin:
heif-convert
heif-dec
heif-enc
heif-info
heif-thumbnailer
hermes
node
npm
npx
opencode-telegram
pihole
---
tools/mule.sh:213:# over five weeks by ten pre-existing lint errors in the checkout's own base — the gate lints
tools/mule.sh:221:# errors, test failures and golangci-lint findings all carry a file:line prefix.
tools/mule.sh:243:  if [ "$rc" = 0 ] && [ -n "${MULE_LINT_CMD:-}" ]; then
tools/mule.sh:244:    note "gate: lint"
tools/mule.sh:245:    ( cd "$REPO" && eval "$MULE_LINT_CMD" ) >>"$GATE_LOG" 2>&1 \
tools/mule.sh:246:      || { gate_reason "lint failed"; bad "$GATE_REASON"; rc=1; }
tools/mule.sh:304:  if [ -n "${MULE_LINT_CMD:-}" ]; then
tools/mule.sh:305:    ( cd "$REPO" && eval "$MULE_LINT_CMD" ) >/dev/null 2>&1 \
tools/mule.sh:306:      && ok "lint baseline is clean" \
tools/mule.sh:307:      || bad "lint baseline is not clean (clear the linter cache before believing this)"
tools/mule.sh:886:   6  hygiene    dependency, CI and lint upkeep

golangci-lint and govulncheck exist at /root/go/bin (not on PATH — the recipe's Mac path is wrong for this box). Running the real checks.
$ /root/go/bin/govulncheck ./... 2>&1 | tail -40
Example traces found:
      #1: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls ssh.NewClientConn

Vulnerability #3: GO-2026-6253
    moby/go-archive: Crafted tar archive can write outside the extraction
    directory in github.com/moby/go-archive
  More info: https://pkg.go.dev/vuln/GO-2026-6253
  Module: github.com/moby/go-archive
    Found in: github.com/moby/go-archive@v0.2.0
    Fixed in: github.com/moby/go-archive@v0.3.0
    Example traces found:
      #1: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls archive.TarWithOptions
      #2: internal/flow/pipeline.go:154:21: flow.Pipeline.MarshalJSON calls json.Marshal, which eventually calls archive.init
      #3: internal/testutil/pg.go:13:2: testutil.init calls testcontainers.init, which calls archive.init
      #4: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls compression.CompressStream
      #5: internal/flow/pipeline.go:154:21: flow.Pipeline.MarshalJSON calls json.Marshal, which eventually calls compression.init
      #6: internal/testutil/pg.go:13:2: testutil.init calls testcontainers.init, which eventually calls compression.init
      #7: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls compression.nopWriteCloser.Close
      #8: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls tarheader.FileInfoHeaderNoLookups
      #9: internal/testutil/pg.go:13:2: testutil.init calls testcontainers.init, which eventually calls tarheader.init
      #10: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls tarheader.nosysFileInfo.Gname
      #11: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls tarheader.nosysFileInfo.Sys
      #12: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls tarheader.nosysFileInfo.Uname

Vulnerability #4: GO-2026-5970
    Infinite loop on invalid input in golang.org/x/text
  More info: https://pkg.go.dev/vuln/GO-2026-5970
  Module: golang.org/x/text
    Found in: golang.org/x/text@v0.38.0
    Fixed in: golang.org/x/text@v0.39.0
    Example traces found:
      #1: internal/store/notify.go:59:34: store.Listen calls pgx.ConnectConfig, which eventually calls norm.Form.Properties
      #2: internal/store/notify.go:59:34: store.Listen calls pgx.ConnectConfig, which eventually calls norm.Form.Span
      #3: internal/store/notify.go:59:34: store.Listen calls pgx.ConnectConfig, which eventually calls norm.Form.Transform

Your code is affected by 4 vulnerabilities from 3 modules.
This scan also found 2 vulnerabilities in packages you import and 3
vulnerabilities in modules you require, but your code doesn't appear to call
these vulnerabilities.
Use '-show verbose' for more details.

$ /root/go/bin/golangci-lint run ./... 2>&1 | tail -40; echo "EXIT=$?"
0 issues.
EXIT=0

$ /root/go/bin/govulncheck ./... 2>&1 | head -60
=== Symbol Results ===

Vulnerability #1: GO-2026-6355
    Prevent DoS on deadlocked established channel in golang.org/x/crypto/ssh
  More info: https://pkg.go.dev/vuln/GO-2026-6355
  Module: golang.org/x/crypto
    Found in: golang.org/x/crypto@v0.53.0
    Fixed in: golang.org/x/crypto@v0.56.0
    Example traces found:
      #1: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls ssh.NewClientConn

Vulnerability #2: GO-2026-6354
    Prevent DoS on deadlocked undecided channel in golang.org/x/crypto/ssh
  More info: https://pkg.go.dev/vuln/GO-2026-6354
  Module: golang.org/x/crypto
    Found in: golang.org/x/crypto@v0.53.0
    Fixed in: golang.org/x/crypto@v0.56.0
    Example traces found:
      #1: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls ssh.NewClientConn

Vulnerability #3: GO-2026-6253
    moby/go-archive: Crafted tar archive can write outside the extraction
    directory in github.com/moby/go-archive
  More info: https://pkg.go.dev/vuln/GO-2026-6253
  Module: github.com/moby/go-archive
    Found in: github.com/moby/go-archive@v0.2.0
    Fixed in: github.com/moby/go-archive@v0.3.0
    Example traces found:
      #1: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls archive.TarWithOptions
      #2: internal/testutil/pg.go:13:2: testutil.init calls testcontainers.init, which calls archive.init
      #3: internal/flow/pipeline.go:154:21: flow.Pipeline.MarshalJSON calls json.Marshal, which eventually calls archive.init
      #4: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls compression.CompressStream
      #5: internal/testutil/pg.go:13:2: testutil.init calls testcontainers.init, which eventually calls compression.init
      #6: internal/flow/pipeline.go:154:21: flow.Pipeline.MarshalJSON calls json.Marshal, which eventually calls compression.init
      #7: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls compression.nopWriteCloser.Close
      #8: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls tarheader.FileInfoHeaderNoLookups
      #9: internal/testutil/pg.go:13:2: testutil.init calls testcontainers.init, which eventually calls tarheader.init
      #10: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls tarheader.nosysFileInfo.Gname
      #11: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls tarheader.nosysFileInfo.Sys
      #12: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls tarheader.nosysFileInfo.Uname

Vulnerability #4: GO-2026-5970
    Infinite loop on invalid input in golang.org/x/text
  More info: https://pkg.go.dev/vuln/GO-2026-5970
  Module: golang.org/x/text
    Found in: golang.org/x/text@v0.38.0
    Fixed in: golang.org/x/text@v0.39.0
    Example traces found:
      #1: internal/store/notify.go:59:34: store.Listen calls pgx.ConnectConfig, which eventually calls norm.Form.Properties
      #2: internal/store/notify.go:59:34: store.Listen calls pgx.ConnectConfig, which eventually calls norm.Form.Span
      #3: internal/store/notify.go:59:34: store.Listen calls pgx.ConnectConfig, which eventually calls norm.Form.Transform

Your code is affected by 4 vulnerabilities from 3 modules.
This scan also found 2 vulnerabilities in packages you import and 3
vulnerabilities in modules you require, but your code doesn't appear to call
these vulnerabilities.
Use '-show verbose' for more details.

% WebFetch https://pkg.go.dev/vuln/GO-2026-5970
% WebFetch https://pkg.go.dev/vuln/GO-2026-6253
% WebFetch https://pkg.go.dev/vuln/GO-2026-6355
% WebFetch https://pkg.go.dev/vuln/GO-2026-6354
govulncheck runs fine on this box now (installed at /root/go/bin). Found 4 reachable advisories. Lint is clean. Now mapping GO-IDs to CVEs and checking the no-test-file packages for logic.
→ Read migrations/migrations.go
→ Read .mule/for-giulio.md
→ Read examples/flow-container-echo/main.go
The vulns are already escalated in for-giulio.md (the x/text entry even names all four GO-IDs). Let me confirm that entry's exact scope and check the example's README for the contract wording.

$ rg -n 'GO-2026|govulncheck' .mule/for-giulio.md | head
202:- **govulncheck GO-2026-5970: reachable DoS in golang.org/x/text (infinite loop on invalid input, fixed in v0.39.0, available v0.41.0).** Astrate pins `x/text` indirect at v0.38.0 (go.mod:97) and pgx pulls it into production: `internal/store/notify.go:59` `store.Listen` → `pgx.ConnectConfig` → `unicode/norm.*`. This is the only govulncheck symbol finding that is not test-harness-only: GO-2026-6355/6354 (x/crypto/ssh deadlocked-channel DoS) and GO-2026-6253 (moby/go-archive tar path traversal) are reachable only through testcontainers in `internal/testutil/pg.go`, i.e. never in the deployed binary. `x/text` keeps API compatibility minor-to-minor and the modules Astrate exercises (`unicode/norm` via pgx, `text/language` via jsonschema) are unchanged, so this is a fix Astrate actually needs — the hygiene recipe's highest-priority category. Not a mule task (go.mod never-touch): your decision to bump ≥v0.39.0 now or fold into the next milestone-boundary sweep. Raw: https://pkg.go.dev/vuln/GO-2026-5970. (The 2026-09-04 dep sweep did not list x/text.)
218:  `golangci-lint` and `govulncheck` have been installed on the Pi the whole time — the two
234:  Note (corrected 2026-09-04): `govulncheck` and `golangci-lint` **are** installed on the Pi (`/root/go/bin`, since 2026-07-28 and 2026-09-01) and `.mule/config` finds them there, so both checks were available; the sweep that produced this list simply ran without invoking govulncheck.

✱ Glob "examples/**/*_test.go" 0 matches
→ Read examples/flow-container-echo/README.md
All four vulns are already escalated (for-giulio.md line 202, still undecided) — identical versions, nothing new. The only new proposal from this hygiene run is the no-test package. Let me check the tail of the queue before appending.
→ Read .mule/todo.md [offset=174]
All four vulns match the existing for-giulio escalation exactly (same GO-IDs, same versions, still undecided) — re-escalating would duplicate it. Lint, TODO and skip scans are clean. Only one package qualifies for the no-test-file rule. Appending that task line.

$ printf '%s\n' '- [ ] examples-echo-container-contract-test [auto]: add a container-free httptest suite for examples/flow-container-echo — the only no-test package in the repo that contains logic (`handleMessage`, main.go:49-75; `sanitizeLog`, main.go:40-47) — pinning (1) POST `/v1/message` echoes a valid JSON-object body verbatim with 200, answers 400 on invalid JSON and 204 on an empty body (the 1 MiB read cap), and (2) metadata `echo_drop` == "1" forces the 204 filter/drop path; the example is the canonical reference container contract (README.md, docs/handoff/flow-design-b-container-block-2026-07-29.md), so a regression silently misleads every container author who copies it.' >> .mule/todo.md && tail -c 400 .mule/todo.md
ge` echoes a valid JSON-object body verbatim with 200, answers 400 on invalid JSON and 204 on an empty body (the 1 MiB read cap), and (2) metadata `echo_drop` == "1" forces the 204 filter/drop path; the example is the canonical reference container contract (README.md, docs/handoff/flow-design-b-container-block-2026-07-29.md), so a regression silently misleads every container author who copies it.

$ tail -n 2 .mule/todo.md; git status --short
- [!] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine_api.yaml — `retrieve_metadata` and `downsample_key` are parsed for every data read by `parseQueryOpts` (internal/appengine/http.go:624 `retrieve_metadata`, 633 `downsample_key`) yet appear on none of the read endpoints: device-scoped GET `{interface}` (yaml:583-606) and `{interface}/{path}` (yaml:623-649), by-alias (yaml:324-368, 369-417), and in-group (yaml:1108-1132, 1149-1176) — the parameter block on each lists only since/since_after/to/limit/downsample_to/sort/format/allow_bigintegers/allow_safe_bigintegers (components DataSince..DataAllowSafeBigIntegers, yaml:1360-1439). Add the two `$ref`s (new DataRetrieveMetadata and DataDownsampleKey components, or inline — mirror the existing style) to all six, and update the DataSort-independent `sort` default note only if verified. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [ ] examples-echo-container-contract-test [auto]: add a container-free httptest suite for examples/flow-container-echo — the only no-test package in the repo that contains logic (`handleMessage`, main.go:49-75; `sanitizeLog`, main.go:40-47) — pinning (1) POST `/v1/message` echoes a valid JSON-object body verbatim with 200, answers 400 on invalid JSON and 204 on an empty body (the 1 MiB read cap), and (2) metadata `echo_drop` == "1" forces the 204 filter/drop path; the example is the canonical reference container contract (README.md, docs/handoff/flow-design-b-container-block-2026-07-29.md), so a regression silently misleads every container author who copies it.
 M .mule/todo.md

Hygiene run carried out on the Pi (govulncheck and golangci-lint both run here — `/root/go/bin`, not the recipe's Mac path).

Done: hygiene recipe run — govulncheck re-verified identical to the existing for-giulio escalation, lint/TODO/skip scans clean, one test-task proposed for the only no-test package with logic
Files: .mule/todo.md (one line appended)
Verified: go list -m -u all (many newer transitives, no proposal per milestone-boundary rule) -> pass; /root/go/bin/govulncheck ./... -> 4 reachable advisories, all already escalated; /root/go/bin/golangci-lint run ./... -> 0 issues; rg TODO/FIXME/XXX/HACK -> only context.TODO() + intake.go:56 (already tracked as broker-external-bus-intake); go test ./... -> no skips, no-test packages = docs/migrations (no logic) + examples/flow-container-echo
Unsure: nothing
Follow-ups: examples-echo-container-contract-test

Notes: the four advisories (GO-2026-5970 CVE-2026-56852 x/text v0.38.0→v0.39.0, plus the x/crypto/ssh pair and go-archive, all test-harness-only) exactly match for-giulio.md's open GO-2026-5970 entry — same GO-IDs, same versions, still undecided — so no duplicate was filed. The `hygiene-govulncheck [legion]` line in todo.md is stale: the Pi can run govulncheck (installed since 2026-07-28, per the 2026-09-04 for-giulio correction), which I did here; it cannot be a mule task to retire it (no code diff → proof gate), so it's reported rather than queued.
