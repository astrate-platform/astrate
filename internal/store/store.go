// Package store is Astrate's single PostgreSQL/TimescaleDB access layer
// (docs/DESIGN.md §1.3): every domain package reads and writes the database
// through it, and it imports none of them. It owns connection pooling, the
// embedded schema migrations (docs/DESIGN.md §2.2–2.5), CA-key sealing, and
// one repository file per aggregate (realms, interfaces, devices, properties,
// datastreams, groups, triggers) plus LISTEN/NOTIFY plumbing.
package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/astrate-platform/astrate/migrations"
)

// defaultPoolMaxConns bounds the pool when the DSN does not specify
// pool_max_conns. The server is tuned to max_connections=50
// (docs/DESIGN.md §5.4); 20 leaves headroom for LISTEN connections
// (Listen dials outside the pool), an optional hot-standby instance,
// and operator sessions.
const defaultPoolMaxConns = 20

// Sentinel errors shared by every repository in this package. Repositories
// wrap them with context; callers test with errors.Is.
var (
	// ErrNotFound reports that the addressed row does not exist.
	ErrNotFound = errors.New("store: not found")
	// ErrAlreadyExists reports a uniqueness conflict (realm name, interface
	// name+major, group name, trigger name, ...).
	ErrAlreadyExists = errors.New("store: already exists")
	// ErrInvalidRealmName reports a realm name rejected by the schema's
	// CHECK constraint (must match ^[a-z][a-z0-9]*$).
	ErrInvalidRealmName = errors.New("store: invalid realm name")
	// ErrDeviceAlreadyConfirmed reports a re-registration attempt for a
	// device that has already requested credentials (docs/DESIGN.md §4.4
	// flow A: 422 conflict parity).
	ErrDeviceAlreadyConfirmed = errors.New("store: device has already requested credentials")
	// ErrInterfaceInUse reports an interface delete blocked because some
	// device still declares it in its introspection.
	ErrInterfaceInUse = errors.New("store: interface is referenced by device introspection")
	// ErrInterfaceMajorNotZero reports an interface delete blocked by the
	// upstream draining rule: only major version 0 interfaces are deletable.
	ErrInterfaceMajorNotZero = errors.New("store: only major version 0 interfaces can be deleted")
)

// Store is the shared database handle: a bounded pgx pool over an
// already-migrated schema. All repository methods hang off it; the zero
// value is not usable, construct with New.
type Store struct {
	pool *pgxpool.Pool

	// hasToolkit records whether the timescaledb_toolkit extension was
	// present at startup (docs/DESIGN.md §2.5 capability probe).
	hasToolkit bool
}

// New connects to dsn, applies any pending embedded migrations, probes
// optional capabilities, and returns a ready Store. The pool is capped at
// defaultPoolMaxConns unless the DSN carries an explicit pool_max_conns.
func New(ctx context.Context, dsn string) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("store: parsing DSN: %w", err)
	}
	if !strings.Contains(dsn, "pool_max_conns") {
		cfg.MaxConns = defaultPoolMaxConns
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("store: creating pool: %w", err)
	}
	s := &Store{pool: pool}

	if err := s.runMigrations(); err != nil {
		pool.Close()
		return nil, err
	}
	if err := s.probeCapabilities(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return s, nil
}

// runMigrations applies the embedded migrations up to head. It is a no-op on
// an up-to-date schema, so every startup (and every test) can call it.
func (s *Store) runMigrations() error {
	// Closing this *sql.DB releases its connections back to the pool; it
	// does not close the pool itself (pgx stdlib.OpenDBFromPool contract).
	db := stdlib.OpenDBFromPool(s.pool)

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		_ = db.Close()
		return fmt.Errorf("store: loading embedded migrations: %w", err)
	}
	driver, err := migratepgx.WithInstance(db, &migratepgx.Config{})
	if err != nil {
		_ = db.Close()
		return fmt.Errorf("store: preparing migration driver: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", src, "pgx/v5", driver)
	if err != nil {
		_ = db.Close()
		return fmt.Errorf("store: preparing migrator: %w", err)
	}

	upErr := m.Up()
	srcErr, dbErr := m.Close() // also closes db via the driver
	if upErr != nil && !errors.Is(upErr, migrate.ErrNoChange) {
		return fmt.Errorf("store: applying migrations: %w", upErr)
	}
	if srcErr != nil {
		return fmt.Errorf("store: closing migration source: %w", srcErr)
	}
	if dbErr != nil {
		return fmt.Errorf("store: closing migration database handle: %w", dbErr)
	}
	return nil
}

// probeCapabilities detects optional database features once at startup.
//
// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5):
// when timescaledb_toolkit is present, Downsample should switch from the
// time_bucket+avg default to toolkit lttb() downsampling. The probe already
// records availability in s.hasToolkit; the time_bucket path in
// datastreams.go is the always-working default.
func (s *Store) probeCapabilities(ctx context.Context) error {
	const q = `SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'timescaledb_toolkit')`
	if err := s.pool.QueryRow(ctx, q).Scan(&s.hasToolkit); err != nil {
		return fmt.Errorf("store: probing timescaledb_toolkit: %w", err)
	}
	return nil
}

// HasToolkitLTTB reports whether the timescaledb_toolkit extension (and so
// lttb-based downsampling) was available at startup.
func (s *Store) HasToolkitLTTB() bool { return s.hasToolkit }

// Health verifies database liveness (readiness probe backend).
func (s *Store) Health(ctx context.Context) error {
	if err := s.pool.Ping(ctx); err != nil {
		return fmt.Errorf("store: health check: %w", err)
	}
	return nil
}

// Close releases the connection pool. The Store is unusable afterwards.
func (s *Store) Close() { s.pool.Close() }

// pgErrCode extracts the PostgreSQL error code from err, or "" if err is not
// a server-reported error.
func pgErrCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

const (
	pgCodeUniqueViolation = "23505"
	pgCodeCheckViolation  = "23514"
	pgCodeFKViolation     = "23503"
)
