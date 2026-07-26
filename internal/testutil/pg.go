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
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	// TimescaleImage is the production-parity database image (docs/DESIGN.md §5.4):
	// tests run against the exact image the docker-compose deployment ships.
	TimescaleImage = "timescale/timescaledb:latest-pg16"

	// EnvTestDSN names the env var that, when set, short-circuits container
	// startup and connects to an already-running database instead (e.g.
	// `make up` + ASTRATE_TEST_DSN=...). This is the fast local iteration
	// path; CI always boots a fresh container.
	EnvTestDSN = "ASTRATE_TEST_DSN"

	// poolMaxConns keeps a single test package well inside the tuned
	// max_connections=50 server budget (docs/DESIGN.md §5.4) even when several
	// T2 suites share one database via EnvTestDSN.
	poolMaxConns = 8

	// readyTimeout bounds how long we wait for the server to answer a ping
	// after the container reports started (crash-recovery/initdb headroom).
	readyTimeout = 60 * time.Second
)

// StartTimescale returns a pgxpool.Pool connected to a ready TimescaleDB
// instance. If EnvTestDSN is set it reuses that database; otherwise it boots a
// fresh TimescaleImage container. Teardown (pool close, container termination)
// is registered on t.Cleanup, so callers just use the pool.
func StartTimescale(t testing.TB) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	if dsn := os.Getenv(EnvTestDSN); dsn != "" {
		t.Logf("testutil: reusing database from %s", EnvTestDSN)
		return openPool(ctx, t, dsn)
	}

	ctr, err := tcpostgres.Run(ctx, TimescaleImage,
		tcpostgres.WithDatabase("astrate_test"),
		tcpostgres.WithUsername("astrate"),
		tcpostgres.WithPassword("astrate-test-password"),
		tcpostgres.BasicWaitStrategies(),
	)
	testcontainers.CleanupContainer(t, ctr)
	if err != nil {
		t.Fatalf("testutil: starting %s container: %v", TimescaleImage, err)
	}

	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("testutil: resolving container connection string: %v", err)
	}
	return openPool(ctx, t, dsn)
}

// openPool builds a bounded pgx pool for dsn and blocks until the server
// answers a ping (or readyTimeout elapses). Closing is hooked on t.Cleanup.
func openPool(ctx context.Context, t testing.TB, dsn string) *pgxpool.Pool {
	t.Helper()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("testutil: parsing DSN: %v", err)
	}
	cfg.MaxConns = poolMaxConns

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("testutil: creating pool: %v", err)
	}
	t.Cleanup(pool.Close)

	deadline := time.Now().Add(readyTimeout)
	for {
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err = pool.Ping(pingCtx)
		cancel()
		if err == nil {
			return pool
		}
		if time.Now().After(deadline) {
			t.Fatalf("testutil: database not ready after %s: %v", readyTimeout, err)
		}
		time.Sleep(250 * time.Millisecond)
	}
}
