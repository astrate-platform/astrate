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
		"SELECT extversion FROM pg_extension WHERE extname = 'timescaledb'",
	).Scan(&version)
	if err != nil {
		t.Fatalf("querying timescaledb extension: %v", err)
	}
	if version == "" {
		t.Fatal("timescaledb extension version is empty")
	}
	t.Logf("timescaledb extension version %s on %s", version, TimescaleImage)

	// The pool must be usable for ordinary work, not just the ready ping.
	var one int
	if err := pool.QueryRow(ctx, "SELECT 1").Scan(&one); err != nil || one != 1 {
		t.Fatalf("SELECT 1 = %d, err %v", one, err)
	}
}
