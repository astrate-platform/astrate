package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// TriggerPolicy is one stored trigger delivery policy (upstream 1.1 Realm
// Management surface): the raw policy JSON. The engine compiles it into its
// realm snapshot and the trigger executor honours it per delivery; the
// deviations from upstream are recorded in docs/COMPATIBILITY.md.
type TriggerPolicy struct {
	ID         int64
	RealmID    int16
	Name       string
	Definition []byte
}

// CreateTriggerPolicy installs a policy; a duplicate name yields
// ErrAlreadyExists.
func (s *Store) CreateTriggerPolicy(ctx context.Context, realmID int16, name string, definition []byte) (*TriggerPolicy, error) {
	p := TriggerPolicy{RealmID: realmID, Name: name, Definition: definition}
	err := s.pool.QueryRow(ctx,
		`INSERT INTO trigger_policies (realm_id, name, definition) VALUES ($1, $2, $3) RETURNING id`,
		realmID, name, definition).Scan(&p.ID)
	if pgErrCode(err) == pgCodeUniqueViolation {
		return nil, fmt.Errorf("%w: policy %q", ErrAlreadyExists, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: creating policy %q: %w", name, err)
	}
	return &p, nil
}

// GetTriggerPolicy fetches one policy by name.
func (s *Store) GetTriggerPolicy(ctx context.Context, realmID int16, name string) (*TriggerPolicy, error) {
	var p TriggerPolicy
	err := s.pool.QueryRow(ctx,
		`SELECT id, realm_id, name, definition FROM trigger_policies WHERE realm_id = $1 AND name = $2`,
		realmID, name).Scan(&p.ID, &p.RealmID, &p.Name, &p.Definition)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: policy %q", ErrNotFound, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: getting policy %q: %w", name, err)
	}
	return &p, nil
}

// DeleteTriggerPolicy removes a policy.
func (s *Store) DeleteTriggerPolicy(ctx context.Context, realmID int16, name string) error {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM trigger_policies WHERE realm_id = $1 AND name = $2`, realmID, name)
	if err != nil {
		return fmt.Errorf("store: deleting policy %q: %w", name, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: policy %q", ErrNotFound, name)
	}
	return nil
}

// ListTriggerPolicies returns every policy of a realm ordered by name.
func (s *Store) ListTriggerPolicies(ctx context.Context, realmID int16) ([]TriggerPolicy, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, realm_id, name, definition FROM trigger_policies WHERE realm_id = $1 ORDER BY name`,
		realmID)
	if err != nil {
		return nil, fmt.Errorf("store: listing policies: %w", err)
	}
	defer rows.Close()

	var out []TriggerPolicy
	for rows.Next() {
		var p TriggerPolicy
		if err := rows.Scan(&p.ID, &p.RealmID, &p.Name, &p.Definition); err != nil {
			return nil, fmt.Errorf("store: scanning policy: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: listing policies: %w", err)
	}
	return out, nil
}
