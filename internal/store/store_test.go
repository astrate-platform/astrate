//go:build integration

package store

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/migrations"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// TestStore is the umbrella T2 suite: one TimescaleDB container (or the
// ASTRATE_TEST_DSN database), one migrated Store, and one sub-suite per
// repository file (docs/ROADMAP.md §3.1 file 2.16). Sub-suites isolate
// through per-run-unique realm names, so the suite is rerunnable against a
// reused database. MigrationCycle tears the schema down to version 1 and
// back up, so it must stay last.
func TestStore(t *testing.T) {
	pool := testutil.StartTimescale(t)
	dsn := pool.Config().ConnString()
	ctx := context.Background()

	s, err := New(ctx, dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(s.Close)

	if err := s.Health(ctx); err != nil {
		t.Fatalf("Health: %v", err)
	}
	t.Logf("timescaledb_toolkit available: %v", s.HasToolkitLTTB())

	t.Run("Migrations", func(t *testing.T) { testMigrations(t, s) })
	t.Run("Realms", func(t *testing.T) { testRealms(t, s) })
	t.Run("Interfaces", func(t *testing.T) { testInterfaces(t, s) })
	t.Run("Devices", func(t *testing.T) { testDevices(t, s) })
	t.Run("Properties", func(t *testing.T) { testProperties(t, s) })
	t.Run("Datastreams", func(t *testing.T) { testDatastreams(t, s) })
	t.Run("Retention", func(t *testing.T) { testRetention(t, s) })
	t.Run("TTLJob", func(t *testing.T) { testTTLJob(t, s) })
	t.Run("Groups", func(t *testing.T) { testGroups(t, s) })
	t.Run("Triggers", func(t *testing.T) { testTriggers(t, s) })
	t.Run("Notify", func(t *testing.T) { testNotify(t, s) })
	t.Run("MigrationCycle", func(t *testing.T) { testMigrationCycle(t, s) })
}

// testMigrations holds the catalog assertions of the M2 gate
// (docs/ROADMAP.md §3.2): schema at head, hypertables present, compression
// configured exactly as docs/DESIGN.md §2.4, policies and TTL job registered.
func testMigrations(t *testing.T, s *Store) {
	ctx := context.Background()

	var (
		version int64
		dirty   bool
	)
	if err := s.pool.QueryRow(ctx, `SELECT version, dirty FROM schema_migrations`).Scan(&version, &dirty); err != nil {
		t.Fatalf("reading schema_migrations: %v", err)
	}
	if version != 5 || dirty {
		t.Fatalf("schema_migrations: got version=%d dirty=%v, want version=5 dirty=false", version, dirty)
	}

	hypertables := []string{"individual_datastreams", "object_datastreams"}

	var n int
	if err := s.pool.QueryRow(ctx, `
		SELECT count(*) FROM timescaledb_information.hypertables
		WHERE hypertable_name = ANY($1)`, hypertables).Scan(&n); err != nil {
		t.Fatalf("querying hypertables: %v", err)
	}
	if n != 2 {
		t.Errorf("hypertables: got %d of %v, want 2", n, hypertables)
	}

	for _, ht := range hypertables {
		var interval string
		if err := s.pool.QueryRow(ctx, `
			SELECT time_interval::text FROM timescaledb_information.dimensions
			WHERE hypertable_name = $1`, ht).Scan(&interval); err != nil {
			t.Fatalf("querying %s dimension: %v", ht, err)
		}
		if interval != "7 days" {
			t.Errorf("%s chunk interval: got %q, want %q", ht, interval, "7 days")
		}

		var segmentby, orderby string
		if err := s.pool.QueryRow(ctx, `
			SELECT segmentby, orderby FROM timescaledb_information.hypertable_compression_settings
			WHERE hypertable = $1::regclass`, ht).Scan(&segmentby, &orderby); err != nil {
			t.Fatalf("querying %s compression settings: %v", ht, err)
		}
		if got, want := normalizeIdentList(segmentby), "realm_id,device_id,interface_id,path"; got != want {
			t.Errorf("%s compress_segmentby: got %q (normalized %q), want %q", ht, segmentby, got, want)
		}
		if got, want := normalizeIdentList(orderby), "ts desc"; got != want {
			t.Errorf("%s compress_orderby: got %q (normalized %q), want %q", ht, orderby, got, want)
		}

		var policies int
		if err := s.pool.QueryRow(ctx, `
			SELECT count(*) FROM timescaledb_information.jobs
			WHERE proc_name = 'policy_compression' AND hypertable_name = $1`, ht).Scan(&policies); err != nil {
			t.Fatalf("querying %s compression policy: %v", ht, err)
		}
		if policies != 1 {
			t.Errorf("%s compression policies: got %d, want 1", ht, policies)
		}
	}

	var ttlSchedule string
	err := s.pool.QueryRow(ctx, `
		SELECT schedule_interval::text FROM timescaledb_information.jobs
		WHERE proc_name = 'astrate_apply_endpoint_ttl'`).Scan(&ttlSchedule)
	if err != nil {
		t.Fatalf("querying TTL job: %v", err)
	}
	if ttlSchedule != "01:00:00" {
		t.Errorf("TTL job schedule: got %q, want %q", ttlSchedule, "01:00:00")
	}

	for _, idx := range []string{"devices_aliases_gin", "ids_series_idx", "ods_series_idx"} {
		var exists bool
		if err := s.pool.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM pg_indexes WHERE schemaname = 'public' AND indexname = $1)`,
			idx).Scan(&exists); err != nil {
			t.Fatalf("querying index %s: %v", idx, err)
		}
		if !exists {
			t.Errorf("index %s missing", idx)
		}
	}
}

// testMigrationCycle exercises the down migrations: head → version 1 → head.
// Destructive (drops every table), hence last in the umbrella.
func testMigrationCycle(t *testing.T, s *Store) {
	db := stdlib.OpenDBFromPool(s.pool)
	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		t.Fatalf("iofs: %v", err)
	}
	driver, err := migratepgx.WithInstance(db, &migratepgx.Config{})
	if err != nil {
		t.Fatalf("migrate driver: %v", err)
	}
	m, err := migrate.NewWithInstance("iofs", src, "pgx/v5", driver)
	if err != nil {
		t.Fatalf("migrator: %v", err)
	}
	defer m.Close()

	if err := m.Migrate(1); err != nil {
		t.Fatalf("migrating down to version 1: %v", err)
	}
	if err := m.Up(); err != nil {
		t.Fatalf("migrating back up: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		t.Fatalf("reading version: %v", err)
	}
	if version != 5 || dirty {
		t.Fatalf("after cycle: got version=%d dirty=%v, want version=5 dirty=false", version, dirty)
	}
}

// normalizeIdentList lowercases and strips spaces/quotes so catalog-rendered
// identifier lists compare against the migration's intent.
func normalizeIdentList(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, `"`, "")
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.Join(strings.Fields(parts[i]), " ")
	}
	return strings.Join(parts, ",")
}

// --- shared helpers -------------------------------------------------------

var realmSeq int

// uniqueRealmName returns a schema-valid realm name unique per test run, so
// suites stay rerunnable against an ASTRATE_TEST_DSN database.
func uniqueRealmName(t *testing.T) string {
	t.Helper()
	realmSeq++
	return "t" + strconv.FormatInt(time.Now().UnixNano(), 36) + strconv.Itoa(realmSeq)
}

// mustCreateRealm creates a uniquely named realm with placeholder CA
// material and returns it.
func mustCreateRealm(t *testing.T, s *Store) *Realm {
	t.Helper()
	sealed := make([]byte, 48)
	if _, err := rand.Read(sealed); err != nil {
		t.Fatal(err)
	}
	realm, err := s.CreateRealm(context.Background(), NewRealm{
		Name:               uniqueRealmName(t),
		JWTPublicKeysPEM:   []string{"-----BEGIN PUBLIC KEY-----\nplaceholder\n-----END PUBLIC KEY-----\n"},
		CACertificatePEM:   "-----BEGIN CERTIFICATE-----\nplaceholder\n-----END CERTIFICATE-----\n",
		CAPrivateKeySealed: sealed,
	})
	if err != nil {
		t.Fatalf("CreateRealm: %v", err)
	}
	return realm
}

// mustRegisterDevice registers a random device in the realm and returns its ID.
func mustRegisterDevice(t *testing.T, s *Store, realmID int16) deviceid.ID {
	t.Helper()
	id, err := deviceid.Random()
	if err != nil {
		t.Fatalf("deviceid.Random: %v", err)
	}
	if err := s.RegisterDevice(context.Background(), realmID, id, fmt.Sprintf("$2a$10$fakehashfor%s", id)); err != nil {
		t.Fatalf("RegisterDevice: %v", err)
	}
	return id
}

// mustInstallInterface installs an interface definition and returns it.
func mustInstallInterface(t *testing.T, s *Store, realmID int16, definition string) *StoredInterface {
	t.Helper()
	si, err := s.InstallInterface(context.Background(), realmID, []byte(definition))
	if err != nil {
		t.Fatalf("InstallInterface: %v", err)
	}
	return si
}

// mustValueType parses a value-type wire string.
func mustValueType(t *testing.T, s string) interfaceschema.ValueType {
	t.Helper()
	vt, err := interfaceschema.ParseValueType(s)
	if err != nil {
		t.Fatalf("ParseValueType(%q): %v", s, err)
	}
	return vt
}
