slug: testutil-control-frame-test
verdict: blocked
at:  80ae7bc
ran: 2026-10-08T17:47:52Z on DietPi in 102s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/testutil/astartedevice.go
→ Read internal/engine/control.go [offset=260, limit=80]
→ Read internal/engine/control_test.go

$ rg -n 'astra.*internal/' internal/testutil/*.go | rg -v 'astartedevice.go' ; ls internal/testutil; rg -n 'internal/testutil' internal/engine/*.go
astartedevice.go
astartedevice_test.go
golden.go
golden_test.go
mqttclient.go
pg.go
pg_test.go
internal/engine/engine_e2e_test.go:25:	"github.com/astrate-platform/astrate/internal/testutil"
internal/engine/engine_integration_test.go:16:	"github.com/astrate-platform/astrate/internal/testutil"

→ Read internal/testutil/astartedevice_test.go

$ head -20 internal/testutil/pg_test.go internal/testutil/golden_test.go internal/testutil/golden.go internal/testutil/pg.go internal/testutil/mqttclient.go | head -80; go version; go env GOOS GOARCH
==> internal/testutil/pg_test.go <==
//go:build integration

package testutil

import (
	"context"
	"testing"
	"time"
)

// TestStartTimescale is the M0 T2 gate (docs/ROADMAP.md §1.2 file 0.7): the
// production-parity container boots and the timescaledb extension is present.
func TestStartTimescale(t *testing.T) {
	pool := StartTimescale(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var version string
	err := pool.QueryRow(ctx,

==> internal/testutil/golden_test.go <==
package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// recordingTB captures Errorf/Fatalf calls so the helpers under test can be
// observed failing without failing the real test.
type recordingTB struct {
	testing.TB
	errors []string
	fatals []string
}

func (r *recordingTB) Helper() {}

==> internal/testutil/golden.go <==
package testutil

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

// updateGolden rewrites golden files instead of comparing against them:
//
//	go test ./... -update
var updateGolden = flag.Bool("update", false, "rewrite golden files with the received content")

// Golden compares got against the golden file testdata/<name> (relative to the
// calling package). With -update the file is (re)written instead and the
// comparison always passes. Golden fixtures are wire-frozen bytes (envelopes,
// payload vectors), so the comparison is exact — byte for byte.
func Golden(t testing.TB, name string, got []byte) {

==> internal/testutil/pg.go <==
// Package testutil provides the shared test harness used across Astrate's
// verification tiers (docs/ROADMAP.md §0.2): the T2 TimescaleDB container
// helper and the golden-file comparison helper.
package testutil

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
go version go1.26.7 linux/arm64
linux
arm64

Now let me check the other testutil files for build tags/DB needs, and probe the exact wire bytes independently (Python zlib vs Go) per MULE.md's verify-with-a-probe rule:
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ mkdir -p /tmp/mule-probe && cat > /tmp/mule-probe/probe.py <<'EOF'
import zlib, binascii, struct
entries = ["com.astrate.test.PropertyArrays/config/thresholds","com.astrate.test.PropertyArrays/config/labels","com.example.Other/a/b"]
plain = ";".join(entries).encode()
frame = struct.pack(">I", len(plain)) + zlib.compress(plain)
print("multi:", binascii.hexlify(frame).decode())
print("len plain:", len(plain))
plain2 = b"a;b;c"
frame2 = struct.pack(">I", len(plain2)) + zlib.compress(plain2)
print("small:", binascii.hexlify(frame2).decode())
plain3 = b""
frame3 = struct.pack(">I", 0) + zlib.compress(plain3)
print("empty:", binascii.hexlify(frame3).decode())
EOF
python3 /tmp/mule-probe/probe.py failed
Error: The user rejected permission to use this specific tool call.
