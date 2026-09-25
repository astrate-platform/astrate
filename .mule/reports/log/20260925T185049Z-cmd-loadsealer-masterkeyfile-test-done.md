slug: cmd-loadsealer-masterkeyfile-test
verdict: done
at:  a6711f9
ran: 2026-09-25T18:50:49Z on DietPi in 221s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read cmd/astrate/main.go [offset=470, limit=80]

$ ls /root/astrate-mule/cmd/astrate/
devcert_test.go
forward_nats_test.go
forward_test.go
main.go
main_test.go
newnats_default.go
newnats_nats.go
shutdown_test.go
version_test.go

$ cd /root/astrate-mule && rg -n "MasterKeyFile|MasterKey" internal/config/*.go | head -40
internal/config/config.go:134:// private keys (§4.3). When MasterKeyFile is empty, the store falls back to
internal/config/config.go:137:	MasterKeyFile string `toml:"master_key_file"`
internal/config/config.go:248:	str("ASTRATE_SECURITY_MASTER_KEY_FILE", &cfg.Security.MasterKeyFile)

$ cd /root/astrate-mule && rg -n "EnvMasterKey|NewKeySealerFromEnv" --type go | head -50
cmd/astrate/main.go:525:		if err := os.Setenv(store.EnvMasterKeyFile, cfg.Security.MasterKeyFile); err != nil {
cmd/astrate/main.go:529:	sealer, err := store.NewKeySealerFromEnv()
internal/store/crypto.go:22:	// EnvMasterKey names the env var holding the master key itself,
internal/store/crypto.go:24:	EnvMasterKey = "ASTRATE_MASTER_KEY"
internal/store/crypto.go:25:	// EnvMasterKeyFile names the env var holding a path to a file that
internal/store/crypto.go:27:	EnvMasterKeyFile = "ASTRATE_MASTER_KEY_FILE"
internal/store/crypto.go:33:var ErrNoMasterKey = fmt.Errorf("store: neither %s nor %s is set", EnvMasterKey, EnvMasterKeyFile)
internal/store/crypto.go:58:// NewKeySealerFromEnv builds a KeySealer from the environment references
internal/store/crypto.go:59:// (EnvMasterKey first, then EnvMasterKeyFile).
internal/store/crypto.go:60:func NewKeySealerFromEnv() (*KeySealer, error) {
internal/store/crypto.go:68:// LoadMasterKey resolves the master key from the environment: EnvMasterKey
internal/store/crypto.go:69:// (hex or base64 text) wins over EnvMasterKeyFile (raw 32 bytes, or hex or
internal/store/crypto.go:72:	if v := os.Getenv(EnvMasterKey); v != "" {
internal/store/crypto.go:75:			return nil, fmt.Errorf("store: %s: %w", EnvMasterKey, err)
internal/store/crypto.go:79:	if path := os.Getenv(EnvMasterKeyFile); path != "" {
internal/store/crypto.go:82:			return nil, fmt.Errorf("store: reading %s: %w", EnvMasterKeyFile, err)
internal/store/crypto.go:89:			return nil, fmt.Errorf("store: %s (%s): %w", EnvMasterKeyFile, path, err)
internal/store/crypto_test.go:126:	t.Setenv(EnvMasterKey, hex.EncodeToString(key))
internal/store/crypto_test.go:143:			t.Setenv(EnvMasterKey, enc.EncodeToString(key))
internal/store/crypto_test.go:163:		t.Setenv(EnvMasterKey, "")
internal/store/crypto_test.go:164:		t.Setenv(EnvMasterKeyFile, path)
internal/store/crypto_test.go:179:		t.Setenv(EnvMasterKey, "")
internal/store/crypto_test.go:180:		t.Setenv(EnvMasterKeyFile, path)
internal/store/crypto_test.go:192:	t.Setenv(EnvMasterKey, "")
internal/store/crypto_test.go:193:	t.Setenv(EnvMasterKeyFile, "")
internal/store/crypto_test.go:198:	t.Setenv(EnvMasterKey, "not-a-key")
internal/store/crypto_test.go:203:	t.Setenv(EnvMasterKey, base64.StdEncoding.EncodeToString(make([]byte, 16)))

→ Read internal/store/crypto.go [limit=110]
→ Read internal/store/crypto_test.go [offset=110, limit=100]
→ Read internal/store/crypto.go [offset=110]

$ cd /root/astrate-mule && rg -n "^func Test|^func " cmd/astrate/main_test.go | head -60
33:func dialDSN(t *testing.T) string {
51:func envOr(k, def string) string {
60:func freeAddr(t *testing.T) string {
73:func bootConfig(t *testing.T, dsn, httpAddr, sessPath, realmName, pub string) config.Config {
90:func TestBoot(t *testing.T) {
200:func waitReady(t *testing.T, base string) {
212:func get(t *testing.T, url string) (int, string) {
225:func mintAgentToken(t *testing.T, key *rsa.PrivateKey) string {
240:func registerDevice(t *testing.T, base, realm, deviceID, token string) string {

→ Read cmd/astrate/main_test.go [limit=100]

$ cd /root/astrate-mule && rg -n "^package|^import|build " cmd/astrate/shutdown_test.go cmd/astrate/version_test.go | head -20
cmd/astrate/shutdown_test.go:1:package main
cmd/astrate/shutdown_test.go:3:import (
cmd/astrate/version_test.go:1:package main
cmd/astrate/version_test.go:3:import (

→ Read internal/config/config.go [offset=125, limit=25]
→ Read cmd/astrate/shutdown_test.go [limit=45]

$ cd /root/astrate-mule && ls && cat .golangci.yml 2>/dev/null || cat .golangci.yaml 2>/dev/null
AGENTS.md
Dockerfile
LICENSE
Makefile
README.md
astrate
bench
cmd
deploy
docker-compose.yml
docs
examples
flow-container-echo
go.mod
go.sum
internal
migrations
opencode.json
pkg
site-dist
test
tools
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

$ cd /root/astrate-mule && rg -n "^lint:|^test:|^vet:|gofmt" Makefile | head -20
31:lint:
37:test:

→ Read Makefile [limit=55, offset=25]

← Write cmd/astrate/sealer_test.go
Wrote file successfully.

← Edit cmd/astrate/sealer_test.go
Index: /root/astrate-mule/cmd/astrate/sealer_test.go
===================================================================
--- /root/astrate-mule/cmd/astrate/sealer_test.go
+++ /root/astrate-mule/cmd/astrate/sealer_test.go
@@ -150,9 +150,9 @@
 }
 if ks != nil {
 	t.Errorf("loadSealer returned a sealer alongside error %v", err)
 }
-if !errors.Is(err, store.ErrNoMasterKey) && file == "" {
+if file == "" && !errors.Is(err, store.ErrNoMasterKey) {
 	t.Errorf("no-reference case: got %v, want it to wrap store.ErrNoMasterKey", err)
 }
 for _, want := range []string{
 	store.EnvMasterKey,

$ cd /root/astrate-mule && gofmt -l . && go vet ./cmd/astrate/ && go test -run 'TestLoadSealer' -v ./cmd/astrate/ 2>&1 | tail -40
=== RUN   TestLoadSealerReadsConfiguredKeyFile
--- PASS: TestLoadSealerReadsConfiguredKeyFile (0.00s)
=== RUN   TestLoadSealerPrefersInlineEnvKey
--- PASS: TestLoadSealerPrefersInlineEnvKey (0.00s)
=== RUN   TestLoadSealerMissingKeyNamesEveryReference
=== RUN   TestLoadSealerMissingKeyNamesEveryReference/no_reference_at_all
=== RUN   TestLoadSealerMissingKeyNamesEveryReference/configured_file_is_absent
--- PASS: TestLoadSealerMissingKeyNamesEveryReference (0.00s)
    --- PASS: TestLoadSealerMissingKeyNamesEveryReference/no_reference_at_all (0.00s)
    --- PASS: TestLoadSealerMissingKeyNamesEveryReference/configured_file_is_absent (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/cmd/astrate	0.447s

Now verifying the tests genuinely fail without the hand-off branch:
$ cd /root/astrate-mule && cp cmd/astrate/main.go /tmp/opencode/main.go.bak && python3 - <<'EOF'
import re
p='cmd/astrate/main.go'
s=open(p).read()
mut = s.replace('''	if cfg.Security.MasterKeyFile != "" {
		if err := os.Setenv(store.EnvMasterKeyFile, cfg.Security.MasterKeyFile); err != nil {
			return nil, err
		}
	}
''','''	_ = cfg.Security.MasterKeyFile
''')
assert mut != s
open(p,'w').write(mut)
EOF
go test -run 'TestLoadSealer' ./cmd/astrate/ 2>&1 | tail -30; cp /tmp/opencode/main.go.bak cmd/astrate/main.go && gofmt -l cmd/astrate/main.go; git diff --stat
--- FAIL: TestLoadSealerReadsConfiguredKeyFile (0.00s)
    sealer_test.go:98: loadSealer with security.master_key_file=/tmp/TestLoadSealerReadsConfiguredKeyFile2479270234/001/master.key: master key: store: neither ASTRATE_MASTER_KEY nor ASTRATE_MASTER_KEY_FILE is set (set ASTRATE_MASTER_KEY, ASTRATE_MASTER_KEY_FILE, or security.master_key_file)
FAIL
FAIL	github.com/astrate-platform/astrate/cmd/astrate	0.442s
FAIL

$ cd /root/astrate-mule && go vet ./... && go test -race ./... 2>&1 | grep -v "^ok\|no test files" | head -20; echo "EXIT_ABOVE"; gofmt -l .
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/cmd/astrate	0.071s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/examples/flow-container-echo	0.101s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine	0.050s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine/channels	0.080s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine/stream	0.066s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/auth	0.099s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
EXIT_ABOVE

$ cd /root/astrate-mule && go test ./... 2>&1 | grep -v "^ok\|no test files" | head -20; echo "---gofmt---"; gofmt -l .
---gofmt---

$ cd /root/astrate-mule && command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./cmd/astrate/ 2>&1 | tail -20 || echo "golangci-lint not installed"
golangci-lint not installed

```
Done: added container-free tests for loadSealer covering the security.master_key_file -> ASTRATE_MASTER_KEY_FILE hand-off, inline ASTRATE_MASTER_KEY precedence, and the error naming all three key references.
Files: cmd/astrate/sealer_test.go (new, no production code changed)
Verified: gofmt -l . (clean) -> pass; go vet ./... -> pass; go test ./... -> pass; go test -race ./... -> not runnable here (ThreadSanitizer: unsupported VMA range, the 39-bit-VMA limitation MULE.md documents, so the gate was the non-race build); mutation check: deleting the os.Setenv branch in main.go:525 makes TestLoadSealerReadsConfiguredKeyFile fail with "neither ASTRATE_MASTER_KEY nor ASTRATE_MASTER_KEY_FILE is set", main.go restored byte-identical afterwards
Unsure: nothing — golangci-lint is not installed on this machine, so only gofmt+vet ran; the file writes use 0o600 and t.Setenv both master-key vars up front so loadSealer's own os.Setenv is restored at cleanup
Follow-ups: none
```
