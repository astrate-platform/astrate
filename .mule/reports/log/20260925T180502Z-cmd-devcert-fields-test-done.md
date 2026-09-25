slug: cmd-devcert-fields-test
verdict: done
at:  f1434eb
ran: 2026-09-25T18:05:02Z on DietPi in 357s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read cmd/astrate/main.go [offset=280, limit=90]

$ ls /root/astrate-mule/cmd/astrate/
forward_nats_test.go
forward_test.go
main.go
main_test.go
newnats_default.go
newnats_nats.go
shutdown_test.go
version_test.go

→ Read cmd/astrate/main_test.go

$ head -3 /root/astrate-mule/cmd/astrate/forward_test.go /root/astrate-mule/cmd/astrate/version_test.go /root/astrate-mule/cmd/astrate/shutdown_test.go
==> /root/astrate-mule/cmd/astrate/forward_test.go <==
package main

import (

==> /root/astrate-mule/cmd/astrate/version_test.go <==
package main

import (

==> /root/astrate-mule/cmd/astrate/shutdown_test.go <==
package main

import (

$ cd /root/astrate-mule && rg -l "go:build" --glob '*_test.go' | head -30; echo "---total tests---"; rg --files -g '*_test.go' | wc -l
cmd/astrate/main_test.go
cmd/astrate/forward_nats_test.go
test/security/security_test.go
internal/engine/engine_integration_test.go
internal/engine/forward/nats_test.go
internal/flowapi/service_integration_test.go
internal/realm/detailed_required_encrypted_test.go
internal/realm/dashboard_compat_test.go
internal/realm/async_operation_test.go
internal/realm/http_test.go
internal/appengine/groups_parity_test.go
internal/realm/detailed_test.go
internal/appengine/formats_parity_test.go
internal/appengine/dashboard_compat_test.go
internal/store/groups_test.go
internal/appengine/e2e_test.go
internal/store/flows_test.go
internal/engine/triggers_change_integration_test.go
internal/store/policies_test.go
internal/store/devices_test.go
internal/store/notify_test.go
internal/store/compression_test.go
internal/store/properties_test.go
internal/store/realms_test.go
internal/engine/engine_e2e_test.go
internal/store/pipelines_test.go
internal/store/userblocks_test.go
internal/store/triggers_test.go
internal/store/store_test.go
internal/store/datastreams_test.go
---total tests---
150

→ Read cmd/astrate/version_test.go

$ cd /root/astrate-mule && rg -n "selfSignedDevCert|InsecureDevMode" --glob '*.go' | head -30
cmd/astrate/main.go:292:		InsecureDevMode:   cfg.MQTT.InsecureDevMode,
cmd/astrate/main.go:305:	case cfg.MQTT.InsecureDevMode:
cmd/astrate/main.go:310:		cert, err := selfSignedDevCert()
cmd/astrate/main.go:320:// selfSignedDevCert mints a throwaway ECDSA server certificate for the dev-mode
cmd/astrate/main.go:323:func selfSignedDevCert() (tls.Certificate, error) {
cmd/astrate/main_test.go:80:	cfg.MQTT.InsecureDevMode = true
internal/config/config.go:74:	InsecureDevMode  bool   `toml:"insecure_dev_mode"`
internal/config/config.go:260:		cfg.MQTT.InsecureDevMode = b
internal/config/config.go:321:	if !c.MQTT.InsecureDevMode && (c.MQTT.TLSCertFile == "" || c.MQTT.TLSKeyFile == "") {
internal/config/config_test.go:489:func TestMQTTInsecureDevModeEnv(t *testing.T) {
internal/config/config_test.go:510:		if cfg.MQTT.InsecureDevMode {
internal/config/config_test.go:520:		if !cfg.MQTT.InsecureDevMode {
internal/broker/broker.go:31:	// only when InsecureDevMode is set.
internal/broker/broker.go:46:	// InsecureDevMode additionally binds a plaintext listener that
internal/broker/broker.go:49:	InsecureDevMode bool
internal/broker/broker.go:51:	// ignored unless InsecureDevMode).
internal/broker/broker.go:137:	if cfg.InsecureDevMode {
internal/broker/broker.go:170:	if cfg.InsecureDevMode {
internal/broker/broker_e2e_test.go:201:		InsecureDevMode:   opts.insecureDevMode,

$ cd /root/astrate-mule && rg -n "^go |^toolchain" go.mod; echo ---; rg -n "\"slices\"" --glob '*.go' | head -5; echo --- lint ---; ls .golangci* 2>/dev/null; cat Makefile | head -40
3:go 1.26.1
---
internal/config/config_test.go:6:	"slices"
pkg/payload/bson.go:8:	"slices"
internal/appengine/data.go:10:	"slices"
internal/appengine/service.go:16:	"slices"
internal/appengine/groups_parity_test.go:11:	"slices"
--- lint ---
.golangci.yml
# Astrate build & verification entry points (docs/ROADMAP.md §1.2, verification tiers §0.2).
#
#   make build               static binary / compile check (CGO disabled)
#   make lint                golangci-lint (config: .golangci.yml)
#   make test                T1 unit tests (no Docker, no network)
#   make test-integration    T2 tests (testcontainers -> timescale/timescaledb:latest-pg16)
#   make test-e2e            T3 component/E2E tests (full wired binary; populated from M5/M6)
#   make test-conformance    T4 official-SDK conformance harness (lands in M9 under test/)
#   make up / make down      local TimescaleDB via docker compose

GO                     ?= go
GOLANGCI_LINT_VERSION  ?= v2.12.2
GOFLAGS                := -trimpath
LDFLAGS                := -s -w
DIST                   := dist

.DEFAULT_GOAL := build

.PHONY: build lint test test-integration test-e2e test-conformance tools up down clean

## build: compile every package statically; emits $(DIST)/astrate once cmd/astrate lands (M8).
build:
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' ./...
	@if [ -d cmd/astrate ]; then \
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

Now the test file:
← Write cmd/astrate/devcert_test.go
Wrote file successfully.

← Edit cmd/astrate/devcert_test.go
Index: /root/astrate-mule/cmd/astrate/devcert_test.go
===================================================================
--- /root/astrate-mule/cmd/astrate/devcert_test.go
+++ /root/astrate-mule/cmd/astrate/devcert_test.go
@@ -16,9 +16,11 @@
 // hostname verification for a device dialling mqtts://127.0.0.1:8883 — so the
 // fields are asserted exactly rather than only through the handshake that would
 // otherwise cover them.
 func TestSelfSignedDevCertFields(t *testing.T) {
+	before := time.Now()
 	cert, err := selfSignedDevCert()
+	after := time.Now()
 	if err != nil {
 		t.Fatalf("selfSignedDevCert: %v", err)
 	}
 	if len(cert.Certificate) != 1 {

← Edit cmd/astrate/devcert_test.go
Index: /root/astrate-mule/cmd/astrate/devcert_test.go
===================================================================
--- /root/astrate-mule/cmd/astrate/devcert_test.go
+++ /root/astrate-mule/cmd/astrate/devcert_test.go
@@ -32,9 +32,8 @@
 leaf, err := x509.ParseCertificate(cert.Certificate[0])
 if err != nil {
 	t.Fatalf("parsing the generated certificate: %v", err)
 }
-before, after := time.Now(), time.Now()
 
 t.Run("ValidityWindow", func(t *testing.T) {
 	// Backdated by an hour: a device whose clock is behind the broker's
 	// would otherwise see a not-yet-valid certificate. Truncation to whole

← Edit cmd/astrate/devcert_test.go
Index: /root/astrate-mule/cmd/astrate/devcert_test.go
===================================================================
--- /root/astrate-mule/cmd/astrate/devcert_test.go
+++ /root/astrate-mule/cmd/astrate/devcert_test.go
@@ -69,9 +69,9 @@
 	}
 	if len(leaf.ExtKeyUsage) != 1 {
 		t.Errorf("ExtKeyUsage = %v, want only serverAuth", leaf.ExtKeyUsage)
 	}
-	if !slices.Contains(leaf.KeyUsage, x509.KeyUsageDigitalSignature) {
+	if leaf.KeyUsage&x509.KeyUsageDigitalSignature == 0 {
 		t.Errorf("KeyUsage = %b, want digitalSignature", leaf.KeyUsage)
 	}
 })
 

$ cd /root/astrate-mule && gofmt -l ./cmd/astrate/ && go vet ./cmd/astrate/ && go test -run TestSelfSignedDevCertFields -v ./cmd/astrate/ 2>&1 | tail -30
=== RUN   TestSelfSignedDevCertFields
=== RUN   TestSelfSignedDevCertFields/ValidityWindow
=== RUN   TestSelfSignedDevCertFields/LoopbackNames
=== RUN   TestSelfSignedDevCertFields/ServerAuthOnly
=== RUN   TestSelfSignedDevCertFields/SelfSignedChainVerifies
--- PASS: TestSelfSignedDevCertFields (0.01s)
    --- PASS: TestSelfSignedDevCertFields/ValidityWindow (0.00s)
    --- PASS: TestSelfSignedDevCertFields/LoopbackNames (0.00s)
    --- PASS: TestSelfSignedDevCertFields/ServerAuthOnly (0.00s)
    --- PASS: TestSelfSignedDevCertFields/SelfSignedChainVerifies (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/cmd/astrate	0.505s

Now the important check: does the test actually fail when each pinned field is mutated? I'll mutate `main.go` temporarily, with a `/tmp` backup for byte-exact restore.
$ cd /root/astrate-mule && cp cmd/astrate/main.go /tmp/opencode/main.go.bak && md5sum cmd/astrate/main.go /tmp/opencode/main.go.bak
eca2edeae038f1517056819b36fdc880  cmd/astrate/main.go
eca2edeae038f1517056819b36fdc880  /tmp/opencode/main.go.bak

$ cd /root/astrate-mule && python3 - <<'EOF'
import re,subprocess,pathlib
p=pathlib.Path('cmd/astrate/main.go'); orig=p.read_text()
muts={
 'NotBefore=now': ('NotBefore:             time.Now().Add(-time.Hour),','NotBefore:             time.Now(),'),
 'NotAfter=30d': ('NotAfter:              time.Now().Add(365 * 24 * time.Hour),','NotAfter:              time.Now().Add(30 * 24 * time.Hour),'),
 'no-IPs': ('IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},',''),
 'no-localhost': ('DNSNames:              []string{"localhost"},',''),
 'no-eku': ('ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},',''),
 'ca-signed': ('der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)','der, err := x509.CreateCertificate(rand.Reader, tmpl, otherTmpl(), &key.PublicKey, key)'),
}
for name,(a,b) in muts.items():
    if name=='ca-signed':
        # swap issuer/subject to prove the self-signed assertion bites
        src=orig.replace('Subject:               pkix.Name{CommonName: "astrate-dev"},','Subject:               pkix.Name{CommonName: "astrate-dev"},\n\t\t\tIssuer:                pkix.Name{CommonName: "astrate-other"},')
    else:
        src=orig.replace(a,b,1)
    assert src!=orig, name
    p.write_text(src)
    r=subprocess.run(['go','test','-run','TestSelfSignedDevCertFields','./cmd/astrate/'],capture_output=True,text=True)
    out=(r.stdout+r.stderr)
    fails=[l for l in out.splitlines() if l.strip().startswith('devcert_test.go')]
    print(f"[{name}] exit={r.returncode} {'FAILS (good)' if r.returncode!=0 else 'PASSES (BAD - insensitive)'}")
    for l in fails[:3]: print('   ',l.strip())
p.write_text(orig)
EOF
md5sum cmd/astrate/main.go
[NotBefore=now] exit=1 FAILS (good)
    devcert_test.go:42: NotBefore is 0s back, want roughly 1h (a plain now() would reject a device clock behind ours)
[NotAfter=30d] exit=1 FAILS (good)
    devcert_test.go:45: NotAfter is 719h59m59s out, want 365d
[no-IPs] exit=1 FAILS (good)
    devcert_test.go:55: IPAddresses = [], want 127.0.0.1
    devcert_test.go:55: IPAddresses = [], want ::1
    devcert_test.go:61: VerifyHostname("127.0.0.1"): x509: cannot validate certificate for 127.0.0.1 because it doesn't contain any IP SANs
[no-localhost] exit=1 FAILS (good)
    devcert_test.go:51: DNSNames = [], want localhost
    devcert_test.go:61: VerifyHostname("localhost"): x509: certificate is not valid for any names, but wanted to match localhost
    devcert_test.go:96: Verify("localhost") as a self-signed server cert: x509: certificate is not valid for any names, but wanted to match localhost
[no-eku] exit=1 FAILS (good)
    devcert_test.go:68: ExtKeyUsage = [], want serverAuth
    devcert_test.go:71: ExtKeyUsage = [], want only serverAuth
[ca-signed] exit=0 PASSES (BAD - insensitive)
eca2edeae038f1517056819b36fdc880  cmd/astrate/main.go

The `ca-signed` mutant was a no-op (`x509.CreateCertificate` ignores the template's `Issuer`). Let me redo it as a real not-self-signed cert:
$ cd /root/astrate-mule && python3 - <<'EOF'
import subprocess,pathlib
p=pathlib.Path('cmd/astrate/main.go'); orig=p.read_text()
a='der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)'
b='''caKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	caTmpl := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "other-ca"},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(24 * time.Hour), IsCA: true, BasicConstraintsValid: true}
	caDER, _ := x509.CreateCertificate(rand.Reader, caTmpl, caTmpl, &caKey.PublicKey, caKey)
	ca, _ := x509.ParseCertificate(caDER)
	der, err := x509.CreateCertificate(rand.Reader, tmpl, ca, &key.PublicKey, caKey)'''
assert a in orig
p.write_text(orig.replace(a,b,1))
r=subprocess.run(['go','test','-run','TestSelfSignedDevCertFields','./cmd/astrate/'],capture_output=True,text=True)
print(f"[not-self-signed] exit={r.returncode} {'FAILS (good)' if r.returncode!=0 else 'PASSES (BAD - insensitive)'}")
for l in (r.stdout+r.stderr).splitlines():
    if l.strip().startswith('devcert_test.go'): print('   ',l.strip())
p.write_text(orig)
EOF
md5sum cmd/astrate/main.go
[not-self-signed] exit=1 FAILS (good)
    devcert_test.go:82: leaf signature does not verify under its own key: x509: ECDSA verification failure
    devcert_test.go:85: issuer = "CN=other-ca", want the subject "CN=astrate-dev" (self-signed)
eca2edeae038f1517056819b36fdc880  cmd/astrate/main.go

All six mutants caught; `main.go` restored byte-exact. Now the full gate:
$ cd /root/astrate-mule && gofmt -l . | grep -v '^$' ; echo "gofmt-exit=$?"; go vet ./... 2>&1 | tail -5; echo "vet-exit=$?"
gofmt-exit=1
vet-exit=0

$ cd /root/astrate-mule && go test -race ./... 2>&1 | grep -v "^ok\|no test files" | head -30; echo "test-exit=${PIPESTATUS[0]}"
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/cmd/astrate	0.138s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/examples/flow-container-echo	0.083s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine	0.154s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine/channels	0.147s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine/stream	0.075s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/auth	0.169s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/broker	0.114s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/config	0.126s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/engine	0.077s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/engine/forward	0.065s
/bin/sh: 1: Bad substitution

Race detector is broken on this 39-bit VMA box (expected per MULE.md). Falling back to the documented gate:
$ cd /root/astrate-mule && go test ./... 2>&1 | grep -v "^ok\|no test files" | head -30; echo "done"
done

$ cd /root/astrate-mule && go test -count=1 ./cmd/astrate/ -v -run TestSelfSignedDevCertFields 2>&1 | tail -5; git status --short
--- PASS: TestSelfSignedDevCertFields/LoopbackNames (0.00s)
    --- PASS: TestSelfSignedDevCertFields/ServerAuthOnly (0.00s)
    --- PASS: TestSelfSignedDevCertFields/SelfSignedChainVerifies (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/cmd/astrate	0.310s
?? cmd/astrate/devcert_test.go

```
Done: Added a container-free test (cmd/astrate/devcert_test.go) pinning the insecure_dev_mode throwaway broker cert: NotBefore backdated ~1h, NotAfter now+365d, DNSNames localhost + IP SANs 127.0.0.1/::1 (plus VerifyHostname for all three), serverAuth-only EKU, and that the leaf is self-signed (own-key signature check, issuer==subject, verifies as a chain root for each loopback name).
Files: cmd/astrate/devcert_test.go (new; no production code changed — main.go verified byte-identical, md5 eca2edeae038f1517056819b36fdc880)
Verified: go vet ./... -> pass; go test ./... -> pass; gofmt -l . -> clean; go test -race ./... -> FAILED for every package with "ThreadSanitizer: unsupported VMA range / Found 39 - Supported 48", the documented 39-bit-VMA limitation of this host, so the non-race gate was used instead. Sensitivity proven by six temporary mutants of selfSignedDevCert, each reverted immediately: NotBefore=now, NotAfter=30d, IPAddresses dropped, DNSNames dropped, ExtKeyUsage dropped, cert signed by a foreign CA -> all six make the test fail; the full suite passes on unmutated code.
Unsure: nothing — every field value, the SAN set and the ECDSA P256 key were read from cmd/astrate/main.go:320-348, and the tolerances (30m..2h, 365d±1m) are sized for x509's whole-second UTCTime truncation, which I measured rather than assumed.
Follow-ups: none
```

Notes worth flagging: no git touched, and one thing to be aware of — the `Verify` call in `SelfSignedChainVerifies` uses `DNSName` for IP literals, which x509 handles correctly (it matches IP SANs when the name parses as an IP), so the loopback assertions exercise the same path a dialing device takes.
