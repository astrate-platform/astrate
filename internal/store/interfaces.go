package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// StoredInterface is one installed interface: the raw JSON definition (the
// source of truth, docs/DESIGN.md §2.2) plus the routing-critical generated
// columns and the stable endpoint-pattern → endpoint-ID map the schema
// compiler stamps into CompiledMapping rows. It implements
// interfaceschema.EndpointIDResolver.
type StoredInterface struct {
	ID          int64
	RealmID     int16
	Name        string
	Major       int
	Minor       int
	Type        interfaceschema.InterfaceType
	Ownership   interfaceschema.Ownership
	Aggregation interfaceschema.Aggregation
	Definition  []byte
	// Endpoints maps the declared endpoint pattern (e.g. "/%{sensor_id}/value")
	// to its endpoints.id. IDs are stable for the lifetime of the
	// (realm, name, major) interface, across minor updates and reloads.
	Endpoints map[string]int64
}

// ResolveInterface implements interfaceschema.EndpointIDResolver.
func (si *StoredInterface) ResolveInterface(name string, major int) (int64, error) {
	if name != si.Name || major != si.Major {
		return 0, fmt.Errorf("store: resolver holds %s v%d, asked for %s v%d", si.Name, si.Major, name, major)
	}
	return si.ID, nil
}

// ResolveEndpoint implements interfaceschema.EndpointIDResolver.
func (si *StoredInterface) ResolveEndpoint(endpoint string) (int64, error) {
	id, ok := si.Endpoints[endpoint]
	if !ok {
		return 0, fmt.Errorf("store: interface %s v%d has no endpoint %q", si.Name, si.Major, endpoint)
	}
	return id, nil
}

// InstallInterface validates definition and inserts the interface row plus
// one endpoints row per mapping, all in one transaction. The generated
// columns (name, versions, type, ownership, aggregation) derive from the
// stored JSON itself. A (realm, name, major) duplicate yields
// ErrAlreadyExists.
func (s *Store) InstallInterface(ctx context.Context, realmID int16, definition []byte) (*StoredInterface, error) {
	iface, err := interfaceschema.ParseInterface(definition)
	if err != nil {
		return nil, fmt.Errorf("store: installing interface: %w", err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("store: beginning interface install: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var id int64
	err = tx.QueryRow(ctx,
		`INSERT INTO interfaces (realm_id, definition) VALUES ($1, $2) RETURNING id`,
		realmID, definition).Scan(&id)
	if pgErrCode(err) == pgCodeUniqueViolation {
		return nil, fmt.Errorf("%w: interface %s v%d", ErrAlreadyExists, iface.Name, iface.Major)
	}
	if err != nil {
		return nil, fmt.Errorf("store: inserting interface %s v%d: %w", iface.Name, iface.Major, err)
	}

	endpoints := make(map[string]int64, len(iface.Mappings))
	for i := range iface.Mappings {
		epID, err := insertEndpoint(ctx, tx, id, &iface.Mappings[i])
		if err != nil {
			return nil, err
		}
		endpoints[iface.Mappings[i].Endpoint] = epID
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("store: committing interface install: %w", err)
	}
	return &StoredInterface{
		ID:          id,
		RealmID:     realmID,
		Name:        iface.Name,
		Major:       iface.Major,
		Minor:       iface.Minor,
		Type:        iface.Type,
		Ownership:   iface.Ownership,
		Aggregation: iface.Aggregation,
		Definition:  definition,
		Endpoints:   endpoints,
	}, nil
}

// UpdateInterface applies a minor-version bump: it replaces the stored
// definition of the existing (realm, name, major) row and inserts endpoints
// for newly added mappings only — existing endpoint rows are never touched,
// keeping their IDs stable. The semantic compatibility check
// (interfaceschema.CheckMinorUpgrade) is the caller's responsibility; this
// method still refuses a definition that drops an existing endpoint.
func (s *Store) UpdateInterface(ctx context.Context, realmID int16, definition []byte) (*StoredInterface, error) {
	iface, err := interfaceschema.ParseInterface(definition)
	if err != nil {
		return nil, fmt.Errorf("store: updating interface: %w", err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("store: beginning interface update: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var id int64
	err = tx.QueryRow(ctx, `
		SELECT id FROM interfaces
		WHERE realm_id = $1 AND name = $2 AND major_version = $3
		FOR UPDATE`,
		realmID, iface.Name, iface.Major).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: interface %s v%d", ErrNotFound, iface.Name, iface.Major)
	}
	if err != nil {
		return nil, fmt.Errorf("store: locking interface %s v%d: %w", iface.Name, iface.Major, err)
	}

	existing, err := loadEndpointIDs(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	declared := make(map[string]bool, len(iface.Mappings))
	for i := range iface.Mappings {
		declared[iface.Mappings[i].Endpoint] = true
	}
	for endpoint := range existing {
		if !declared[endpoint] {
			return nil, fmt.Errorf("store: updating interface %s v%d: definition drops existing endpoint %q (minor updates are additive only)",
				iface.Name, iface.Major, endpoint)
		}
	}

	if _, err := tx.Exec(ctx, `UPDATE interfaces SET definition = $2 WHERE id = $1`, id, definition); err != nil {
		return nil, fmt.Errorf("store: updating interface %s v%d definition: %w", iface.Name, iface.Major, err)
	}

	for i := range iface.Mappings {
		m := &iface.Mappings[i]
		if _, ok := existing[m.Endpoint]; ok {
			continue
		}
		epID, err := insertEndpoint(ctx, tx, id, m)
		if err != nil {
			return nil, err
		}
		existing[m.Endpoint] = epID
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("store: committing interface update: %w", err)
	}
	return &StoredInterface{
		ID:          id,
		RealmID:     realmID,
		Name:        iface.Name,
		Major:       iface.Major,
		Minor:       iface.Minor,
		Type:        iface.Type,
		Ownership:   iface.Ownership,
		Aggregation: iface.Aggregation,
		Definition:  definition,
		Endpoints:   existing,
	}, nil
}

// DeleteInterface removes an installed interface, enforcing the upstream
// draining rules: only major version 0 is deletable, and only while no
// device declares it in its introspection. Properties cascade through their
// foreign key; datastream rows (no foreign key, docs/DESIGN.md §2.4) are
// swept in the same transaction.
func (s *Store) DeleteInterface(ctx context.Context, realmID int16, name string, major int) error {
	if major != 0 {
		return fmt.Errorf("%w: interface %s v%d", ErrInterfaceMajorNotZero, name, major)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("store: beginning interface delete: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var id int64
	err = tx.QueryRow(ctx, `
		SELECT id FROM interfaces
		WHERE realm_id = $1 AND name = $2 AND major_version = $3
		FOR UPDATE`,
		realmID, name, major).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: interface %s v%d", ErrNotFound, name, major)
	}
	if err != nil {
		return fmt.Errorf("store: locking interface %s v%d: %w", name, major, err)
	}

	var inUse bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM devices
			WHERE realm_id = $1 AND (introspection->$2->>'major')::int = $3
		)`,
		realmID, name, major).Scan(&inUse)
	if err != nil {
		return fmt.Errorf("store: checking introspection references of %s v%d: %w", name, major, err)
	}
	if inUse {
		return fmt.Errorf("%w: interface %s v%d", ErrInterfaceInUse, name, major)
	}

	for _, table := range []string{"individual_datastreams", "object_datastreams"} {
		if _, err := tx.Exec(ctx, `DELETE FROM `+table+` WHERE realm_id = $1 AND interface_id = $2`, realmID, id); err != nil {
			return fmt.Errorf("store: deleting %s of interface %s v%d: %w", table, name, major, err)
		}
	}
	if _, err := tx.Exec(ctx, `DELETE FROM interfaces WHERE id = $1`, id); err != nil {
		return fmt.Errorf("store: deleting interface %s v%d: %w", name, major, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("store: committing interface delete: %w", err)
	}
	return nil
}

// GetInterface loads one installed interface with its endpoint IDs.
func (s *Store) GetInterface(ctx context.Context, realmID int16, name string, major int) (*StoredInterface, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT id, realm_id, name, major_version, minor_version, type, ownership, aggregation, definition
		FROM interfaces
		WHERE realm_id = $1 AND name = $2 AND major_version = $3`,
		realmID, name, major)
	si, err := scanStoredInterface(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: interface %s v%d", ErrNotFound, name, major)
	}
	if err != nil {
		return nil, fmt.Errorf("store: getting interface %s v%d: %w", name, major, err)
	}
	si.Endpoints, err = loadEndpointIDs(ctx, s.pool, si.ID)
	if err != nil {
		return nil, err
	}
	return si, nil
}

// LoadRealmInterfaces loads every interface of a realm with endpoint IDs —
// the input the engine's schema-compiler cache rebuilds from
// (docs/ROADMAP.md §3.1 file 2.9 "LoadRealm").
func (s *Store) LoadRealmInterfaces(ctx context.Context, realmID int16) ([]*StoredInterface, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, realm_id, name, major_version, minor_version, type, ownership, aggregation, definition
		FROM interfaces
		WHERE realm_id = $1
		ORDER BY name, major_version`, realmID)
	if err != nil {
		return nil, fmt.Errorf("store: loading realm %d interfaces: %w", realmID, err)
	}
	defer rows.Close()

	var (
		out  []*StoredInterface
		byID = map[int64]*StoredInterface{}
	)
	for rows.Next() {
		si, err := scanStoredInterface(rows)
		if err != nil {
			return nil, fmt.Errorf("store: scanning interface: %w", err)
		}
		si.Endpoints = map[string]int64{}
		out = append(out, si)
		byID[si.ID] = si
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: loading realm %d interfaces: %w", realmID, err)
	}
	rows.Close()

	epRows, err := s.pool.Query(ctx, `
		SELECT e.interface_id, e.endpoint, e.id
		FROM endpoints e
		JOIN interfaces i ON i.id = e.interface_id
		WHERE i.realm_id = $1`, realmID)
	if err != nil {
		return nil, fmt.Errorf("store: loading realm %d endpoints: %w", realmID, err)
	}
	defer epRows.Close()
	for epRows.Next() {
		var (
			ifaceID  int64
			endpoint string
			epID     int64
		)
		if err := epRows.Scan(&ifaceID, &endpoint, &epID); err != nil {
			return nil, fmt.Errorf("store: scanning endpoint: %w", err)
		}
		if si, ok := byID[ifaceID]; ok {
			si.Endpoints[endpoint] = epID
		}
	}
	if err := epRows.Err(); err != nil {
		return nil, fmt.Errorf("store: loading realm %d endpoints: %w", realmID, err)
	}
	return out, nil
}

// queryRower is the subset of pgx.Tx / pgxpool.Pool the helpers below need.
type queryRower interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// insertEndpoint inserts one endpoints row for mapping m and returns its ID.
func insertEndpoint(ctx context.Context, q queryRower, interfaceID int64, m *interfaceschema.Mapping) (int64, error) {
	var ttl *int64
	if m.DatabaseRetentionPolicy == interfaceschema.UseTTL {
		t := m.DatabaseRetentionTTL
		ttl = &t
	}
	var id int64
	err := q.QueryRow(ctx, `
		INSERT INTO endpoints (interface_id, endpoint, value_type, reliability, retention,
		                       expiry, database_retention_policy, database_retention_ttl,
		                       explicit_timestamp, allow_unset)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`,
		interfaceID, m.Endpoint, m.Type.String(), m.Reliability.String(), m.Retention.String(),
		m.Expiry, m.DatabaseRetentionPolicy.String(), ttl,
		m.ExplicitTimestamp, m.AllowUnset).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("store: inserting endpoint %q: %w", m.Endpoint, err)
	}
	return id, nil
}

// loadEndpointIDs returns the endpoint-pattern → ID map of one interface.
func loadEndpointIDs(ctx context.Context, q queryRower, interfaceID int64) (map[string]int64, error) {
	rows, err := q.Query(ctx, `SELECT endpoint, id FROM endpoints WHERE interface_id = $1`, interfaceID)
	if err != nil {
		return nil, fmt.Errorf("store: loading endpoints of interface %d: %w", interfaceID, err)
	}
	defer rows.Close()

	out := map[string]int64{}
	for rows.Next() {
		var (
			endpoint string
			id       int64
		)
		if err := rows.Scan(&endpoint, &id); err != nil {
			return nil, fmt.Errorf("store: scanning endpoint: %w", err)
		}
		out[endpoint] = id
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: loading endpoints of interface %d: %w", interfaceID, err)
	}
	return out, nil
}

// scanStoredInterface reads one interfaces row (without endpoints).
func scanStoredInterface(row pgx.Row) (*StoredInterface, error) {
	var (
		si                          StoredInterface
		typ, ownership, aggregation string
	)
	if err := row.Scan(&si.ID, &si.RealmID, &si.Name, &si.Major, &si.Minor,
		&typ, &ownership, &aggregation, &si.Definition); err != nil {
		return nil, err
	}
	var err error
	if si.Type, err = interfaceschema.ParseInterfaceType(typ); err != nil {
		return nil, fmt.Errorf("store: interface %s v%d: %w", si.Name, si.Major, err)
	}
	if si.Ownership, err = interfaceschema.ParseOwnership(ownership); err != nil {
		return nil, fmt.Errorf("store: interface %s v%d: %w", si.Name, si.Major, err)
	}
	if si.Aggregation, err = interfaceschema.ParseAggregation(aggregation); err != nil {
		return nil, fmt.Errorf("store: interface %s v%d: %w", si.Name, si.Major, err)
	}
	return &si, nil
}
