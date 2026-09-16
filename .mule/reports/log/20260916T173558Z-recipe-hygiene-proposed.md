slug: recipe-hygiene
verdict: proposed
at:  7eaf6f1
ran: 2026-09-16T17:35:58Z on DietPi in 758s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md
Now I'll run the hygiene checks. Let me start with the ones that are quick.

$ which govulncheck && govulncheck ./... 2>/dev/null || echo "govulncheck not installed"
govulncheck not installed

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ | head -30
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},

$ go list -m -u all 2>/dev/null | rg '\[' | head -20
cel.dev/expr v0.24.0 [v0.25.3]
cloud.google.com/go v0.121.6 [v0.123.0]
cloud.google.com/go/auth v0.16.4 [v0.23.2]
cloud.google.com/go/compute/metadata v0.8.0 [v0.9.0]
cloud.google.com/go/iam v1.5.2 [v1.13.0]
cloud.google.com/go/longrunning v0.6.7 [v1.2.0]
cloud.google.com/go/monitoring v1.24.2 [v1.30.0]
cloud.google.com/go/spanner v1.85.0 [v1.95.1]
cloud.google.com/go/storage v1.56.0 [v1.67.1]
github.com/99designs/keyring v1.2.1 [v1.2.2]
github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0 [v1.23.1]
github.com/Azure/azure-sdk-for-go/sdk/internal v1.1.2 [v1.12.0]
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0 [v1.8.1]
github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c [v0.0.0-20260914234754-15785f84ba1c]
github.com/Azure/go-autorest/autorest/adal v0.9.16 (retracted) [v0.9.24] (deprecated)
github.com/Azure/go-autorest/autorest/date v0.3.0 [v0.3.1] (deprecated)
github.com/Azure/go-autorest/logger v0.2.1 [v0.2.2] (deprecated)
github.com/Azure/go-autorest/tracing v0.6.0 [v0.6.1] (deprecated)
github.com/ClickHouse/clickhouse-go v1.4.3 [v1.5.4]
github.com/DataDog/zstd v1.4.5 [v1.5.7]

$ go test ./... 2>&1 | rg -i 'skip|no test files' | head -20
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]

$ which golangci-lint && golangci-lint run ./... 2>&1 | tail -30 || echo "golangci-lint not installed"
golangci-lint not installed

→ Read internal/broker/intake.go [offset=40, limit=30]
→ Read internal/flow/blocks/httpblocks.go [offset=138, limit=18]

$ ls -la examples/flow-container-echo/ && wc -l examples/flow-container-echo/*.go 2>/dev/null
total 20
drwxr-xr-x 2 root root 4096 Sep  4 20:46 .
drwxr-xr-x 4 root root 4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root  370 Sep  4 20:46 Dockerfile
-rw-r--r-- 1 root root 1694 Sep  4 20:46 README.md
-rw-r--r-- 1 root root 2076 Sep  4 20:46 main.go
75 examples/flow-container-echo/main.go

→ Read internal/flow/blocks/astartesource/source.go [offset=48, limit=16]
govulncheck and golangci-lint aren't installed. Let me check the remaining context: the intake.go TODO, context.TODO usages, the untested example package, and the current todo/for-giulio files.
→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ ls /root/go/bin 2>/dev/null; echo "---"; cat .mule/config 2>/dev/null; echo "---"; echo $PATH
goimports
golangci-lint
govulncheck
---
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
---
/root/.opencode/bin:/root/.opencode/bin:/usr/local/bin:/usr/bin:/bin

Let me check the environment more carefully — for-giulio.md says golangci-lint and govulncheck are installed on the Pi at `/root/go/bin`, and lint findings are a proposal source.
$ rg -n 't\.Skip|skip\.' internal/ pkg/ cmd/ --glob '*_test.go' | head -20; echo "---"; rg -n 'SkipIf|testing.Short' internal/ pkg/ cmd/ | head -10
cmd/astrate/main_test.go:47:	t.Skip("the boot suite needs a database: run `make up` or set ASTRATE_TEST_DSN")
cmd/astrate/forward_test.go:68:		t.Skip("built with -tags nats; see forward_nats_test.go")
internal/pairing/http_test.go:629:		t.Skip("openssl not found in PATH; chain already verified with crypto/x509")
internal/store/datastreams_test.go:420:			t.Skip("timescaledb_toolkit not installed; lttb path unavailable")
internal/pairing/ca/ca_test.go:293:		t.Skip("10k issuance draw skipped in -short mode")
---
internal/pairing/ca/ca_test.go:292:	if testing.Short() {

$ cat > /tmp/notest.sh <<'EOF'
go list ./... 2>/dev/null | while read p; do
  if ! ls "$p"/*_test.go >/dev/null 2>&1; then echo "no tests: $p"; fi
done
EOF
bash /tmp/notest.sh
no tests: github.com/astrate-platform/astrate/cmd/astrate
no tests: github.com/astrate-platform/astrate/docs
no tests: github.com/astrate-platform/astrate/examples/flow-container-echo
no tests: github.com/astrate-platform/astrate/internal/appengine
no tests: github.com/astrate-platform/astrate/internal/appengine/channels
no tests: github.com/astrate-platform/astrate/internal/appengine/stream
no tests: github.com/astrate-platform/astrate/internal/auth
no tests: github.com/astrate-platform/astrate/internal/broker
no tests: github.com/astrate-platform/astrate/internal/config
no tests: github.com/astrate-platform/astrate/internal/engine
no tests: github.com/astrate-platform/astrate/internal/engine/forward
no tests: github.com/astrate-platform/astrate/internal/engine/stream
no tests: github.com/astrate-platform/astrate/internal/engine/triggers
no tests: github.com/astrate-platform/astrate/internal/flow
no tests: github.com/astrate-platform/astrate/internal/flow/blocks
no tests: github.com/astrate-platform/astrate/internal/flow/blocks/astartesource
no tests: github.com/astrate-platform/astrate/internal/flow/blocks/container
no tests: github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool
no tests: github.com/astrate-platform/astrate/internal/flowapi
no tests: github.com/astrate-platform/astrate/internal/housekeeping
no tests: github.com/astrate-platform/astrate/internal/httpx
no tests: github.com/astrate-platform/astrate/internal/observability
no tests: github.com/astrate-platform/astrate/internal/pairing
no tests: github.com/astrate-platform/astrate/internal/pairing/ca
no tests: github.com/astrate-platform/astrate/internal/realm
no tests: github.com/astrate-platform/astrate/internal/store
no tests: github.com/astrate-platform/astrate/internal/swagger
no tests: github.com/astrate-platform/astrate/internal/testutil
no tests: github.com/astrate-platform/astrate/migrations
no tests: github.com/astrate-platform/astrate/pkg/astarteapi
no tests: github.com/astrate-platform/astrate/pkg/deviceid
no tests: github.com/astrate-platform/astrate/pkg/interfaceschema
no tests: github.com/astrate-platform/astrate/pkg/payload

PATH doesn't include `/root/go/bin` but the tools are there. We're on the Pi. Let me run golangci-lint and govulncheck with the full paths.
$ /root/go/bin/golangci-lint run ./... 2>&1 | tail -30
0 issues.

$ timeout 280 /root/go/bin/govulncheck ./... 2>&1 | tail -40; echo "exit=$?"
exit=0

$ timeout 280 /root/go/bin/govulncheck -json ./... > /tmp/gvn.json 2>/tmp/gvn.err; echo "exit=$?"; echo "--- stdout size:"; wc -c /tmp/gvn.json; echo "--- stderr:"; head -c 2000 /tmp/gvn.err
exit=124
--- stdout size:
9797 /tmp/gvn.json
--- stderr:

govulncheck timed out at 280s mid-scan (exit 124) — this box is too slow, which is exactly why `hygiene-govulncheck [legion]` is already queued. Let me peek at whatever the partial scan recorded before giving up here.
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cat /tmp/gvn.err; echo "---"; python3 -c "
import json
d=json.load(open('/tmp/gvn.json'))
print('version:', d.get('Version'))
vm=d.get('Vulnerabilities')
print('vulns:', len(vm))
for v in vm:
    print(v['OSV'].get('id'), v['OSV'].get('aliases'), v['OSV'].get('summary','')[:90])
" 2>&1 | head -40 failed
Error: The user rejected permission to use this specific tool call.
