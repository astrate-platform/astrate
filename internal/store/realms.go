package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Realm is one tenancy row (docs/DESIGN.md §1.5, §2.2). CAPrivateKeySealed is
// the AES-256-GCM box produced by KeySealer; the store never sees the
// plaintext key (sealing happens in the pairing/CA layer).
type Realm struct {
	ID                      int16
	Name                    string
	JWTPublicKeysPEM        []string
	CACertificatePEM        string
	CAPrivateKeySealed      []byte
	DeviceRegistrationLimit *int32
	CreatedAt               time.Time
}

// NewRealm carries the fields needed to create a realm. CA material is
// mandatory: a realm without a CA cannot pair devices.
type NewRealm struct {
	Name                    string
	JWTPublicKeysPEM        []string
	CACertificatePEM        string
	CAPrivateKeySealed      []byte
	DeviceRegistrationLimit *int32
}

const realmColumns = `id, name, jwt_public_keys, ca_certificate, ca_private_key, device_registration_limit, created_at`

// CreateRealm inserts the realm row together with its CA material in one
// transaction (docs/ROADMAP.md §3.1 file 2.8) and returns the stored realm.
func (s *Store) CreateRealm(ctx context.Context, nr NewRealm) (*Realm, error) {
	keys := nr.JWTPublicKeysPEM
	if keys == nil {
		keys = []string{}
	}
	row := s.pool.QueryRow(ctx, `
		INSERT INTO realms (name, jwt_public_keys, ca_certificate, ca_private_key, device_registration_limit)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+realmColumns,
		nr.Name, keys, nr.CACertificatePEM, nr.CAPrivateKeySealed, nr.DeviceRegistrationLimit)

	realm, err := scanRealm(row)
	switch pgErrCode(err) {
	case pgCodeUniqueViolation:
		return nil, fmt.Errorf("%w: realm %q", ErrAlreadyExists, nr.Name)
	case pgCodeCheckViolation:
		return nil, fmt.Errorf("%w: %q", ErrInvalidRealmName, nr.Name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: creating realm %q: %w", nr.Name, err)
	}
	return realm, nil
}

// GetRealm fetches a realm by ID.
func (s *Store) GetRealm(ctx context.Context, id int16) (*Realm, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+realmColumns+` FROM realms WHERE id = $1`, id)
	realm, err := scanRealm(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: realm id %d", ErrNotFound, id)
	}
	if err != nil {
		return nil, fmt.Errorf("store: getting realm id %d: %w", id, err)
	}
	return realm, nil
}

// GetRealmByName fetches a realm by name.
func (s *Store) GetRealmByName(ctx context.Context, name string) (*Realm, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+realmColumns+` FROM realms WHERE name = $1`, name)
	realm, err := scanRealm(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: realm %q", ErrNotFound, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: getting realm %q: %w", name, err)
	}
	return realm, nil
}

// ListRealms returns every realm ordered by name.
func (s *Store) ListRealms(ctx context.Context) ([]Realm, error) {
	rows, err := s.pool.Query(ctx, `SELECT `+realmColumns+` FROM realms ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("store: listing realms: %w", err)
	}
	defer rows.Close()

	var realms []Realm
	for rows.Next() {
		r, err := scanRealm(rows)
		if err != nil {
			return nil, fmt.Errorf("store: scanning realm: %w", err)
		}
		realms = append(realms, *r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: listing realms: %w", err)
	}
	return realms, nil
}

// SetRealmJWTPublicKeys replaces the realm's JWT public key set (PEM strings).
func (s *Store) SetRealmJWTPublicKeys(ctx context.Context, name string, keysPEM []string) error {
	if keysPEM == nil {
		keysPEM = []string{}
	}
	tag, err := s.pool.Exec(ctx,
		`UPDATE realms SET jwt_public_keys = $2 WHERE name = $1`, name, keysPEM)
	if err != nil {
		return fmt.Errorf("store: updating realm %q JWT keys: %w", name, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: realm %q", ErrNotFound, name)
	}
	return nil
}

// SetRealmCA replaces the realm's CA certificate and sealed private key
// (re-keying flow, docs/DESIGN.md §4.3).
func (s *Store) SetRealmCA(ctx context.Context, name, caCertPEM string, caKeySealed []byte) error {
	tag, err := s.pool.Exec(ctx,
		`UPDATE realms SET ca_certificate = $2, ca_private_key = $3 WHERE name = $1`,
		name, caCertPEM, caKeySealed)
	if err != nil {
		return fmt.Errorf("store: updating realm %q CA: %w", name, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: realm %q", ErrNotFound, name)
	}
	return nil
}

// DeleteRealm removes the realm and everything it owns in one transaction
// (docs/DESIGN.md §2.1: realm deletion is a transactional cascade).
// Metadata and properties cascade through foreign keys; the datastream
// hypertables have no foreign keys (§2.4) and are swept explicitly.
func (s *Store) DeleteRealm(ctx context.Context, name string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("store: beginning realm delete: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var realmID int16
	err = tx.QueryRow(ctx, `SELECT id FROM realms WHERE name = $1`, name).Scan(&realmID)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: realm %q", ErrNotFound, name)
	}
	if err != nil {
		return fmt.Errorf("store: resolving realm %q: %w", name, err)
	}

	for _, table := range []string{"individual_datastreams", "object_datastreams"} {
		if _, err := tx.Exec(ctx, `DELETE FROM `+table+` WHERE realm_id = $1`, realmID); err != nil {
			return fmt.Errorf("store: deleting %s of realm %q: %w", table, name, err)
		}
	}
	if _, err := tx.Exec(ctx, `DELETE FROM realms WHERE id = $1`, realmID); err != nil {
		return fmt.Errorf("store: deleting realm %q: %w", name, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("store: committing realm delete: %w", err)
	}
	return nil
}

// scanRealm reads one realm row in realmColumns order.
func scanRealm(row pgx.Row) (*Realm, error) {
	var r Realm
	if err := row.Scan(&r.ID, &r.Name, &r.JWTPublicKeysPEM, &r.CACertificatePEM,
		&r.CAPrivateKeySealed, &r.DeviceRegistrationLimit, &r.CreatedAt); err != nil {
		return nil, err
	}
	return &r, nil
}
