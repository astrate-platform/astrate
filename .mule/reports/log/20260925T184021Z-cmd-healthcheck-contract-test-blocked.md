slug: cmd-healthcheck-contract-test
verdict: blocked
at:  93da5b6
ran: 2026-09-25T18:40:21Z on DietPi in 199s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/task.md
→ Read .mule/MULE.md
→ Read cmd/astrate/main.go [offset=500, limit=80]

$ ls cmd/astrate/ && sed -n '30,40p' Dockerfile
devcert_test.go
forward_nats_test.go
forward_test.go
main.go
main_test.go
newnats_default.go
newnats_nats.go
shutdown_test.go
version_test.go
# HTTP API, mTLS MQTT, and (insecure_dev_mode only) plaintext MQTT.
EXPOSE 8080 8883 1883
# Only writable path; the rest of the FS can be mounted read-only.
VOLUME ["/var/lib/astrate"]

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=5 \
    CMD ["/astrate", "-healthcheck"]

ENTRYPOINT ["/astrate"]

$ wc -l cmd/astrate/main_test.go && sed -n '1,60p' cmd/astrate/main_test.go && echo ---- && grep -rn "healthcheck" cmd/astrate/ --include=*.go
265 cmd/astrate/main_test.go
//go:build integration

package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"

	"github.com/astrate-platform/astrate/internal/config"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

const composeDSN = "postgres://astrate:astrate@127.0.0.1:5432/astrate?sslmode=disable"

// dialDSN finds a reachable database: ASTRATE_TEST_DSN first, then compose.
func dialDSN(t *testing.T) string {
	t.Helper()
	for _, dsn := range []string{envOr("ASTRATE_TEST_DSN", ""), composeDSN} {
		if dsn == "" {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		conn, err := pgx.Connect(ctx, dsn)
		cancel()
		if err == nil {
			_ = conn.Close(context.Background())
			return dsn
		}
	}
	t.Skip("the boot suite needs a database: run `make up` or set ASTRATE_TEST_DSN")
	return ""
}

func envOr(k, def string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return def
}

// freeAddr returns a currently-free loopback address (small TOCTOU window,
// acceptable in a test).
func freeAddr(t *testing.T) string {
----
cmd/astrate/main.go:64:	healthcheck := flag.Bool("healthcheck", false, "probe the local readiness endpoint and exit (for container HEALTHCHECK)")
cmd/astrate/main.go:71:	if *healthcheck {
cmd/astrate/main.go:538:// `astrate -healthcheck`. The HTTP address comes from ASTRATE_HTTP_ADDR.

$ for f in cmd/astrate/shutdown_test.go cmd/astrate/version_test.go cmd/astrate/forward_test.go; do echo "== $f"; sed -n '1,12p' $f; done
== cmd/astrate/shutdown_test.go
package main

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"
)

// TestDrainStageBudgetIsPerStage pins that every graceful-drain stage gets its
// own shutdownTimeout budget: a stage that burns all of it (a slow in-flight
== cmd/astrate/version_test.go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/realm"
)

// TestRealmScopedVersionReportsAPICompatLevel proves the realm-scoped version
== cmd/astrate/forward_test.go
package main

import (
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/config"
)

// quietLogger keeps newForwarder's info line out of the test output.

Now I'll write the container-free test file:
← Write cmd/astrate/healthcheck_test.go
Wrote file successfully.

← Write cmd/astrate/healthcheck_test.go
Wrote file successfully.

$ gofmt -l cmd/astrate/ && gofmt -w cmd/astrate/healthcheck_test.go && go vet ./cmd/astrate/ && go test -run 'TestRunHealthcheck' -v ./cmd/astrate/ 2>&1 | tail -30
=== RUN   TestRunHealthcheckProbesReadiness
    healthcheck_test.go:57: runHealthcheck() = 1, want 0 on a 200 readiness
--- FAIL: TestRunHealthcheckProbesReadiness (0.00s)
=== RUN   TestRunHealthcheckExitsZeroOnlyOn200
=== RUN   TestRunHealthcheckExitsZeroOnlyOn200/ok
    healthcheck_test.go:88: runHealthcheck() on readiness 200 = 1, want 0
    healthcheck_test.go:91: probed [/astrate/v1/readiness], want [/astarte/v1/readiness]
=== RUN   TestRunHealthcheckExitsZeroOnlyOn200/service_unavailable
    healthcheck_test.go:91: probed [/astrate/v1/readiness], want [/astarte/v1/readiness]
=== RUN   TestRunHealthcheckExitsZeroOnlyOn200/not_found
    healthcheck_test.go:91: probed [/astrate/v1/readiness], want [/astarte/v1/readiness]
=== RUN   TestRunHealthcheckExitsZeroOnlyOn200/server_error
    healthcheck_test.go:91: probed [/astrate/v1/readiness], want [/astarte/v1/readiness]
--- FAIL: TestRunHealthcheckExitsZeroOnlyOn200 (0.01s)
    --- FAIL: TestRunHealthcheckExitsZeroOnlyOn200/ok (0.00s)
    --- FAIL: TestRunHealthcheckExitsZeroOnlyOn200/service_unavailable (0.00s)
    --- FAIL: TestRunHealthcheckExitsZeroOnlyOn200/not_found (0.00s)
    --- FAIL: TestRunHealthcheckExitsZeroOnlyOn200/server_error (0.00s)
=== RUN   TestRunHealthcheckMapsBarePortToLoopback
    healthcheck_test.go:110: runHealthcheck() on ":39179" = 1, want 0 (bare port must map to 127.0.0.1)
--- FAIL: TestRunHealthcheckMapsBarePortToLoopback (0.00s)
=== RUN   TestRunHealthcheckFailsWhenNothingAnswers
--- PASS: TestRunHealthcheckFailsWhenNothingAnswers (0.00s)
FAIL
FAIL	github.com/astrate-platform/astrate/cmd/astrate	0.474s
FAIL

$ cd cmd/astrate && grep -n 'astrate/v1/readiness' healthcheck_test.go | od -c | grep -A2 -B2 'r e a d' | head -20
(no output)

$ grep -n 'readiness' cmd/astrate/healthcheck_test.go | cat -A | sed -n '1,40p'
11:// readinessRecorder answers only /astarte/v1/readiness with status and records$
15:type readinessRecorder struct {$
21:func (rec *readinessRecorder) ServeHTTP(w http.ResponseWriter, r *http.Request) {$
25:^Iif r.URL.Path != "/astarte/v1/readiness" {$
32:func (rec *readinessRecorder) probed() []string {$
38:// probeAddr starts a readiness endpoint answering with status and returns its$
40:func probeAddr(t *testing.T, status int) (string, *readinessRecorder) {$
42:^Irec := &readinessRecorder{status: status}$
49:// ["/astrate", "-healthcheck"]` contract probes: /astrate/v1/readiness, and$
57:^I^It.Fatalf("runHealthcheck() = %d, want 0 on a 200 readiness", code)$
63:^Iif paths[0] != "/astarte/v1/readiness" {$
64:^I^It.Errorf("probed %q, want %q", paths[0], "/astarte/v1/readiness")$
70:// readiness (a reachable but unready instance) marks the container unhealthy.$
88:^I^I^I^It.Errorf("runHealthcheck() on readiness %d = %d, want %d", tt.status, code, tt.want)$
90:^I^I^Iif paths := rec.probed(); len(paths) != 1 || paths[0] != "/astarte/v1/readiness" {$
91:^I^I^I^It.Errorf("probed %v, want [/astarte/v1/readiness]", paths)$

$ go test -run 'TestRunHealthcheckProbesReadiness|TestRunHealthcheckExitsZeroOnlyOn200' -v ./cmd/astrate/ 2>&1 | head -40; echo "=== env proxies:"; env | grep -i proxy
=== RUN   TestRunHealthcheckProbesReadiness
    healthcheck_test.go:57: runHealthcheck() = 1, want 0 on a 200 readiness
--- FAIL: TestRunHealthcheckProbesReadiness (0.01s)
=== RUN   TestRunHealthcheckExitsZeroOnlyOn200
=== RUN   TestRunHealthcheckExitsZeroOnlyOn200/ok
    healthcheck_test.go:88: runHealthcheck() on readiness 200 = 1, want 0
    healthcheck_test.go:91: probed [/astrate/v1/readiness], want [/astarte/v1/readiness]
=== RUN   TestRunHealthcheckExitsZeroOnlyOn200/service_unavailable
    healthcheck_test.go:91: probed [/astrate/v1/readiness], want [/astarte/v1/readiness]
=== RUN   TestRunHealthcheckExitsZeroOnlyOn200/not_found
    healthcheck_test.go:91: probed [/astrate/v1/readiness], want [/astarte/v1/readiness]
=== RUN   TestRunHealthcheckExitsZeroOnlyOn200/server_error
    healthcheck_test.go:91: probed [/astrate/v1/readiness], want [/astarte/v1/readiness]
--- FAIL: TestRunHealthcheckExitsZeroOnlyOn200 (0.02s)
    --- FAIL: TestRunHealthcheckExitsZeroOnlyOn200/ok (0.00s)
    --- FAIL: TestRunHealthcheckExitsZeroOnlyOn200/service_unavailable (0.00s)
    --- FAIL: TestRunHealthcheckExitsZeroOnlyOn200/not_found (0.00s)
    --- FAIL: TestRunHealthcheckExitsZeroOnlyOn200/server_error (0.00s)
FAIL
FAIL	github.com/astrate-platform/astrate/cmd/astrate	0.478s
FAIL
=== env proxies:

$ grep -nP '[^\x00-\x7F]' cmd/astrate/main.go | head; echo "--- non-ascii count above"; sed -n '551p' cmd/astrate/main.go | od -c | head -6
1:// Command astrate is the single Astrate binary (docs/ROADMAP.md §9 file 8.4):
56:// (docs/DESIGN.md §5.3). It is a var so tests can shrink it; a stage that
94:// it in the §5.3 order: HTTP → broker → flow manager → engine → store. It is
193:	// Realm datastream retention ceilings (#72): sweep hourly — and once
194:	// right now, before the first tick — capping even no_ttl interfaces
249:// not HTTP webhooks (docs/DESIGN.md §1.1). The empty kind is the default and
253:// The nil return is deliberately a nil *interface*, not a typed nil pointer —
393:	// Pairing's is public — both measured on upstream 1.2.0
433:// (COMPATIBILITY.md deviation 10) — not Astrate's build version. RM's
442:// it is accepted (docs/DESIGN.md §5.2 readiness broker check).
--- non-ascii count above
0000000  \t   r   e   q   ,       e   r   r       :   =       h   t   t
0000020   p   .   N   e   w   R   e   q   u   e   s   t   W   i   t   h
0000040   C   o   n   t   e   x   t   (   c   t   x   ,       h   t   t
0000060   p   .   M   e   t   h   o   d   G   e   t   ,       "   h   t
0000100   t   p   :   /   /   "   +   a   d   d   r   +   "   /   a   s
0000120   t   r   a   t   e   /   v   1   /   r   e   a   d   i   n   e

! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cat > /tmp/zz_debug_test.go <<'EOF'
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestZZDebug(t *testing.T) {
	var got string
	var status int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.Path
		w.WriteHeader(status)
	}))
	defer srv.Close()
	status = http.StatusOK
	t.Setenv("ASTRATE_HTTP_ADDR", srv.Listener.Addr().String())
	code := runHealthcheck()
	t.Logf("code=%d got=%q len=%d want=%q equal=%v", code, got, len(got), "/astrate/v1/readiness", got == "/astrate/v1/readiness")
}
EOF
cp /tmp/zz_debug_test.go cmd/astrate/zz_debug_test.go && go test -run TestZZDebug -v ./cmd/astrate/ 2>&1 | head -20 failed
Error: The user rejected permission to use this specific tool call.
