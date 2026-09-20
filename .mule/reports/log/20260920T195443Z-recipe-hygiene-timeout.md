slug: recipe-hygiene
verdict: timeout
at:  395d8bb
ran: 2026-09-20T19:54:45Z on DietPi in 1286s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md

$ govulncheck ./... 2>/dev/null || echo "govulncheck not installed"
govulncheck not installed

$ ls .mule/ recipes 2>/dev/null; ls .mule/ 2>/dev/null
.mule/:
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

$ command -v golangci-lint && golangci-lint version || echo "golangci-lint not on PATH"
golangci-lint not on PATH

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ | head -30
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):

$ go test ./... 2>&1 | rg -i 'skip|no test files' | head -20
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]

→ Read internal/broker/intake.go [offset=50, limit=15]

$ ls -a .golangci.yml .golangci.yaml 2>/dev/null; ls .mule/config 2>/dev/null; cat .mule/config 2>/dev/null
.golangci.yml
.mule/config
# mule config for this repo. Sourced as bash; env vars override it.
#
# This file is read on the Mac AND on the Pi, which do not have the same toolchain, so it
# detects rather than hardcodes. Measured facts behind the detection:
#
#   Mac (arm64 darwin) : go test -race ./...  green, golangci-lint pinned in ~/go/bin
#   Pi  (arm64 linux)  : go test ./...        green in ~3m over 20 packages
#                        go test -race        IMPOSSIBLE — see below
#
# **The race detector cannot run on the Pi.** ThreadSanitizer needs a 48-bit VMA and this
# kernel is built with 39 ("FATAL: ThreadSanitizer: unsupported VMA range / Found 39 -
# Supported 48"). That is a kernel config, not a flake, and no amount of retrying changes it.
# So the unattended gate is weaker than the Mac's by construction: a data race passes on the
# Pi. Race coverage lives on the Legion Go instead (Go 1.26.5 at ~/.local/go, full suite
# green in ~40s on 16 cores) via the standing [legion] race-check task — which is therefore
# not optional housekeeping but the concurrency gate itself.

_gcl="$(command -v golangci-lint || echo "$HOME/go/bin/golangci-lint")"

if go test -race -run XXX_nonexistent ./pkg/payload >/dev/null 2>&1; then
  MULE_TEST_CMD="go test -race ./..."
else
  # No race detector here. `go vet` is not a replacement, but it is free and it catches a
  # different class of real defect, so it is not nothing.
  MULE_TEST_CMD="go vet ./... && go test ./..."
fi

if [ -x "$_gcl" ]; then
  MULE_LINT_CMD="$_gcl run ./..."
  MULE_FIX_CMD="gofmt -w . && $_gcl run --fix ./... >/dev/null 2>&1 || true"
else
  MULE_LINT_CMD=""            # absent on the Pi; gofmt still runs
  MULE_FIX_CMD="gofmt -w ."
fi

# Paths the mule must never modify: frozen spec, schema decisions, Giulio's voice, build pins.
MULE_NEVER="docs/DESIGN.md docs/ROADMAP.md docs/site/* migrations/* .github/* deploy/* go.mod go.sum Dockerfile docker-compose.yml .trickle/*"

# The mule never commits to main.
MULE_BRANCH="mule/queue"

# Per-task budget in seconds. A task needing more than this was specced too big. Generous
# here because the Pi is slow: the gate alone is 3 minutes of it.
MULE_TIMEOUT="1200"

# How many tasks the timer may run per day. The provider is free; do not hammer it. At one
# task per 30-minute tick the ceiling would be 48, so this is a deliberate ~3x reduction.
MULE_DAILY_MAX="16"

MULE_AGENT="build"
# Pin the model explicitly. Leaving this empty lets opencode resolve the agent's default,
# and since 2026-08-27 that resolution answers "UnknownError: Unexpected server error" every
# single time — it killed four nights of the daily survey. `--model opencode/big-pickle`
# succeeds on the same box, same agent, same second.
MULE_MODEL="opencode/big-pickle"

$ ls ~/go/bin 2>/dev/null; ls /usr/local/bin | rg -i 'golangci|govuln|staticcheck' 2>/dev/null; echo "---"; find / -name 'golangci-lint' -type f 2>/dev/null | head -5
goimports
golangci-lint
govulncheck
---
/root/go/bin/golangci-lint
/root/go/pkg/mod/github.com/docker/docker@v28.3.3+incompatible/hack/validate/golangci-lint

golangci-lint and govulncheck exist at `~/go/bin` (Pi home, not the Mac path in the recipe). Running both:
$ ~/go/bin/golangci-lint version && ~/go/bin/golangci-lint run ./... 2>&1 | tail -30; echo "EXIT: $?"
golangci-lint has version 2.12.2 built with go1.26.5 from (unknown, modified: ?, mod sum: "h1:7+d1uY0bq1MU2UV3R5pW5Q7QWdcoq4naMRXM+gsJKrs=") on (unknown)
0 issues.
EXIT: 0

