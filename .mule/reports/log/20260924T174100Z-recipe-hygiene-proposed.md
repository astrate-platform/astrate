slug: recipe-hygiene
verdict: proposed
at:  fab3c51
ran: 2026-09-24T17:41:00Z on DietPi in 583s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md
Proposal job. Let me run the recipe's checks.

$ govulncheck ./... 2>/dev/null || echo "govulncheck not installed"; which govulncheck
govulncheck not installed

$ command -v golangci-lint && golangci-lint run ./... 2>&1 | tail -30 || echo "golangci-lint not found"
golangci-lint not found

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ 2>/dev/null | head -40
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):

$ go list -m -u all 2>/dev/null | rg '\[' | head -20
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

→ Read internal/broker/intake.go [offset=40, limit=40]
Found one real TODO. Let me read it, and check the test-skip line, plus try to get govulncheck/golangci-lint.

$ ls /usr/local/bin /usr/bin ~/go/bin 2>/dev/null | rg -i 'golangci|govuln' ; echo "---"; go env GOFLAGS GOPATH
golangci-lint
govulncheck
---

/root/go

$ go test ./... 2>&1 | rg -i 'skip|no test files' | head -20
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]

$ ls -la /usr/local/bin/golangci-lint /usr/local/bin/govulncheck $HOME/go/bin/golangci-lint $HOME/go/bin/govulncheck 2>&1
ls: cannot access '/usr/local/bin/golangci-lint': No such file or directory
ls: cannot access '/usr/local/bin/govulncheck': No such file or directory
-rwxr-xr-x 1 root root 65622584 Jul 28 23:52 /root/go/bin/golangci-lint
-rwxr-xr-x 1 root root 13997377 Sep 15 20:20 /root/go/bin/govulncheck

$ PATH=$HOME/go/bin:$PATH golangci-lint run ./... 2>&1 | tail -30
0 issues.

$ PATH=$HOME/go/bin:$PATH govulncheck ./... 2>&1 | tail -40
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

$ PATH=$HOME/go/bin:$PATH govulncheck ./... 2>&1 | head -60
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

govulncheck: 4 reachable vulns (3 modules). golangci-lint: 0 issues. No skips. One real TODO. Let me check the existing todo.md and for-giulio.md formats.
$ ls .mule/; echo ---; tail -40 .mule/for-giulio.md 2>/dev/null
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
recipes
reports
reviews
tasks
todo.md
---
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
- **The mule has been idle 17h.** Filed by the dead-man's switch; see journalctl on the Pi.

→ Read .mule/todo.md
→ Read .mule/for-giulio.md [limit=60]
Let me check the for-giulio.md head for format, and the tail of todo.md to know where to append. Note line 86 already covers the intake.go TODO — let me verify I'm not duplicating it.

$ rg -n 'golang.org/x/crypto|golang.org/x/text|moby/go-archive' go.mod
21:	golang.org/x/crypto v0.53.0 // bcrypt cost 10 for credentials secrets (§4.1)
63:	github.com/moby/go-archive v0.2.0 // indirect
97:	golang.org/x/text v0.38.0 // indirect

Now let me check whether these vulnerable modules are direct deps and what uses them, and confirm current vs fixed versions.
$ rg -n 'golang.org/x/crypto|golang.org/x/text|moby/go-archive' go.sum | head
106:github.com/moby/go-archive v0.2.0 h1:zg5QDUM2mi0JIM9fdQZWC7U8+2ZfixfTYoHL7rWUcP8=
107:github.com/moby/go-archive v0.2.0/go.mod h1:mNeivT14o8xU+5q1YnNrkQVpK+dnNe/K6fHqnTg4qPU=
204:golang.org/x/crypto v0.53.0 h1:QZ4Muo8THX6CizN2vPPd5fBGHyogrdK9fG4wLPFUsto=
205:golang.org/x/crypto v0.53.0/go.mod h1:DNLU434OwVakk9PzuwV8w62mAJpRJL3vsgcfp4Qnsio=
217:golang.org/x/text v0.38.0 h1:sXmwo9DwP3OK9EZ7PqAdaooSGozfl/3a6/xJcbzPRhE=
218:golang.org/x/text v0.38.0/go.mod h1:YXZt3QhHUKYT53r2lLKFIVi6Ao1jdzrTR/KQ09qyxF4=

$ rg -n 'ssh|norm\.' internal/ --glob '!**/*_test.go' | rg -v 'context|TODO' | head -20
(no output)

→ Read internal/store/notify.go
The intake.go TODO is already an open task line (todo.md:86), so no new task there. Now let me verify the production reachability of the x/text path for the report.

$ go mod why golang.org/x/text; echo ---; go list -deps ./internal/... 2>/dev/null | rg -c 'golang.org/x/text' ; rg -rn 'x/text' $(go env GOMODCACHE)/github.com/jackc/pgx*/@v*/go.mod 2>/dev/null | head
# golang.org/x/text
(main module does not need package golang.org/x/text)
---
24

$ go mod why -m golang.org/x/text 2>/dev/null; echo ---; rg -l 'golang.org/x/text' $(go env GOMODCACHE)/github.com/jackc/pgx/v5@v5*/ 2>/dev/null | head
# golang.org/x/text
github.com/astrate-platform/astrate/internal/flowapi
github.com/santhosh-tekuri/jsonschema/v6
golang.org/x/text/language
---
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.7.5/go.sum
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.10.0/go.sum
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.7.5/pgconn/auth_scram.go
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.10.0/CHANGELOG.md
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.7.5/go.mod
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.7.5/CHANGELOG.md
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.10.0/pgconn/auth_scram.go
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.10.0/go.mod

Let me check what in pgx's connect path actually hits `norm.Form`, and verify the testutil/golangci config exclusions.
$ grep -c 'golang.org/x/text' /dev/null; rg -ln 'x/text' $(go env GOMODCACHE)/github.com/jackc/pgx* 2>/dev/null | head; rg -n 'norm\.' $(go env GOMODCACHE)/github.com/jackc/pgx/v5@v5.8.0/pgconn/*.go 2>/dev/null | head
0
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.7.5/go.mod
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.7.5/CHANGELOG.md
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.7.5/go.sum
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.10.0/go.sum
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.7.5/pgconn/auth_scram.go
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.10.0/CHANGELOG.md
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.10.0/pgconn/auth_scram.go
/root/go/pkg/mod/github.com/jackc/pgx/v5@v5.10.0/go.mod

$ M=$(go env GOMODCACHE); for m in "golang.org/x/text@v0.39.0" "golang.org/x/crypto@v0.56.0" "github.com/moby/go-archive@v0.3.0"; do echo "== $m"; ls -d "$M/$m" 2>/dev/null || echo "  (not in cache)"; done
== golang.org/x/text@v0.39.0
  (not in cache)
== golang.org/x/crypto@v0.56.0
  (not in cache)
== github.com/moby/go-archive@v0.3.0
  (not in cache)

Only one proposable finding block: the 3 vulnerable modules. Let me check the changelog/breaking-change angle for each, since the bump proposal needs it.
$ rg -ln 'x/crypto/bcrypt|bcrypt' internal/ pkg/ cmd/ | head; echo ---; rg -n 'bcrypt' go.mod
internal/config/config.go
internal/config/config.example.toml
internal/pairing/service_test.go
internal/pairing/http_test.go
internal/pairing/service.go
---
21:	golang.org/x/crypto v0.53.0 // bcrypt cost 10 for credentials secrets (§4.1)

govulncheck found real issues. Let me check for-giulio.md for any existing mentions to avoid duplication.
$ rg -n -i 'govuln|vuln|CVE|GO-2026|x/text|x/crypto|go-archive' .mule/for-giulio.md | head -20
8:One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
210:- **govulncheck GO-2026-5970: reachable DoS in golang.org/x/text (infinite loop on invalid input, fixed in v0.39.0, available v0.41.0).** Astrate pins `x/text` indirect at v0.38.0 (go.mod:97) and pgx pulls it into production: `internal/store/notify.go:59` `store.Listen` → `pgx.ConnectConfig` → `unicode/norm.*`. This is the only govulncheck symbol finding that is not test-harness-only: GO-2026-6355/6354 (x/crypto/ssh deadlocked-channel DoS) and GO-2026-6253 (moby/go-archive tar path traversal) are reachable only through testcontainers in `internal/testutil/pg.go`, i.e. never in the deployed binary. `x/text` keeps API compatibility minor-to-minor and the modules Astrate exercises (`unicode/norm` via pgx, `text/language` via jsonschema) are unchanged, so this is a fix Astrate actually needs — the hygiene recipe's highest-priority category. Not a mule task (go.mod never-touch): your decision to bump ≥v0.39.0 now or fold into the next milestone-boundary sweep. Raw: https://pkg.go.dev/vuln/GO-2026-5970. (The 2026-09-04 dep sweep did not list x/text.)
226:  `golangci-lint` and `govulncheck` have been installed on the Pi the whole time — the two
241:  - `golang.org/x/crypto` v0.53.0 → v0.56.0 — x/crypto keeps API compatibility; used only for bcrypt in `internal/auth`.
242:  Note (corrected 2026-09-04): `govulncheck` and `golangci-lint` **are** installed on the Pi (`/root/go/bin`, since 2026-07-28 and 2026-09-01) and `.mule/config` finds them there, so both checks were available; the sweep that produced this list simply ran without invoking govulncheck.

GO-2026-5970 (x/text) is already escalated. Let me read that region to see what's been decided and whether x/crypto + go-archive were covered.
→ Read .mule/for-giulio.md [offset=195, limit=55]
govulncheck findings are all already escalated (for-giulio.md:210 covers GO-2026-5970 and the two test-harness-only ones; the bump decision is Giulio's pending). Let me verify the only two test-less packages contain no logic before concluding.

$ ls docs/ migrations/ 2>/dev/null; echo ---; ls docs/*.go migrations/*.go 2>/dev/null
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
---
docs/embed.go
migrations/migrations.go

All checks done. Findings require no new proposals — here's the assessment:

- **govulncheck** (highest-priority check): 4 reachable findings in 3 modules. **All already escalated** — for-giulio.md:210 covers GO-2026-5970 (x/text, the only production-reachable one — `store.Listen`→pgx→`unicode/norm`) as "fixed in v0.39.0" and explicitly notes the other two (x/crypto/ssh, moby/go-archive) are test-harness-only via testcontainers. I re-confirmed the exact same call paths and versions; no new info, so no duplicate escalation.
- **golangci-lint**: 0 issues.
- **TODO scan**: only `context.TODO()` uses (httpblocks.go:147, astartesource/source.go:56) plus one real TODO comment at internal/broker/intake.go:56 — that one names a real missing behaviour but is **already an open task line** (todo.md:86 `broker-external-bus-intake [legion]`). No new line.
- **Skipped tests**: none.
- **No-test-file packages**: only `docs/` and `migrations/` — both embed shims (embed.go), no logic, and `migrations/` is on the never-touch list anyway.
- **`go list -m -u`**: bump sweep is decided at milestone boundaries per MULE.md; the only bump Astrate actually needs (x/text CVE) is already pending Giulio's decision.

Nothing appended to todo.md, nothing added to for-giulio.md, no `gh issue create` (recipe never asks for one — bumps are a proposal, not a new-issue deliverable).

Done: hygiene sweep found nothing not already escalated; x/text GO-2026-5970 re-confirmed present and still pending Giulio's bump decision
Files: none changed
Verified: govulncheck ./... (4 reachable, all previously escalated); golangci-lint run ./... -> 0 issues; go test ./... (no skips); rg TODO scan -> intake.go TODO already queued
Unsure: nothing
Follow-ups: none
