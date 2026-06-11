package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Trigger is one installed trigger (docs/DESIGN.md §2.2): the raw Astarte
// trigger JSON (simple_triggers + action), executed by the engine. Matching
// upstream Realm Management, triggers have no update operation —
// reconfiguration is delete + reinstall.
type Trigger struct {
	ID         int64
	RealmID    int16
	Name       string
	Definition []byte
}

// CreateTrigger installs a trigger; a duplicate name yields ErrAlreadyExists.
func (s *Store) CreateTrigger(ctx context.Context, realmID int16, name string, definition []byte) (*Trigger, error) {
	t := Trigger{RealmID: realmID, Name: name, Definition: definition}
	err := s.pool.QueryRow(ctx,
		`INSERT INTO triggers (realm_id, name, definition) VALUES ($1, $2, $3) RETURNING id`,
		realmID, name, definition).Scan(&t.ID)
	if pgErrCode(err) == pgCodeUniqueViolation {
		return nil, fmt.Errorf("%w: trigger %q", ErrAlreadyExists, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: creating trigger %q: %w", name, err)
	}
	return &t, nil
}

// GetTrigger fetches one trigger by name.
func (s *Store) GetTrigger(ctx context.Context, realmID int16, name string) (*Trigger, error) {
	var t Trigger
	err := s.pool.QueryRow(ctx,
		`SELECT id, realm_id, name, definition FROM triggers WHERE realm_id = $1 AND name = $2`,
		realmID, name).Scan(&t.ID, &t.RealmID, &t.Name, &t.Definition)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: trigger %q", ErrNotFound, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: getting trigger %q: %w", name, err)
	}
	return &t, nil
}

// DeleteTrigger removes a trigger.
func (s *Store) DeleteTrigger(ctx context.Context, realmID int16, name string) error {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM triggers WHERE realm_id = $1 AND name = $2`, realmID, name)
	if err != nil {
		return fmt.Errorf("store: deleting trigger %q: %w", name, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: trigger %q", ErrNotFound, name)
	}
	return nil
}

// ListTriggers returns every trigger of a realm ordered by name (the
// engine's trigger-cache load).
func (s *Store) ListTriggers(ctx context.Context, realmID int16) ([]Trigger, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, realm_id, name, definition FROM triggers WHERE realm_id = $1 ORDER BY name`,
		realmID)
	if err != nil {
		return nil, fmt.Errorf("store: listing triggers: %w", err)
	}
	defer rows.Close()

	var out []Trigger
	for rows.Next() {
		var t Trigger
		if err := rows.Scan(&t.ID, &t.RealmID, &t.Name, &t.Definition); err != nil {
			return nil, fmt.Errorf("store: scanning trigger: %w", err)
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: listing triggers: %w", err)
	}
	return out, nil
}
