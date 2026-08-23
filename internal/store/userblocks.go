package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// UserBlock is one stored per-realm user-defined composite block (issue #85,
// astarte_flow parity): a named producer/consumer body inlined at flow start
// by the flow engine. ConfigSchema is an optional JSON Schema for the block's
// params and is nil when the column is NULL.
type UserBlock struct {
	ID           int64
	RealmID      int16
	Name         string
	BlockType    string // producer, consumer or producer_consumer (DB CHECK)
	Source       []byte
	ConfigSchema []byte // nil when the column is NULL
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// CreateUserBlock stores a user block; a duplicate name yields
// ErrAlreadyExists. An invalid BlockType is rejected by the DB CHECK and
// surfaces as a plain wrapped error.
func (s *Store) CreateUserBlock(ctx context.Context, realmID int16, ub *UserBlock) (*UserBlock, error) {
	b := UserBlock{
		RealmID:      realmID,
		Name:         ub.Name,
		BlockType:    ub.BlockType,
		Source:       ub.Source,
		ConfigSchema: ub.ConfigSchema,
	}
	err := s.pool.QueryRow(ctx,
		`INSERT INTO user_blocks (realm_id, name, block_type, source, config_schema)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at, updated_at`,
		realmID, ub.Name, ub.BlockType, ub.Source, ub.ConfigSchema).
		Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
	if pgErrCode(err) == pgCodeUniqueViolation {
		return nil, fmt.Errorf("%w: user block %q", ErrAlreadyExists, ub.Name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: creating user block %q: %w", ub.Name, err)
	}
	return &b, nil
}

// GetUserBlock fetches one user block by name.
func (s *Store) GetUserBlock(ctx context.Context, realmID int16, name string) (*UserBlock, error) {
	var b UserBlock
	err := s.pool.QueryRow(ctx,
		`SELECT id, realm_id, name, block_type, source, config_schema, created_at, updated_at
		 FROM user_blocks WHERE realm_id = $1 AND name = $2`,
		realmID, name).
		Scan(&b.ID, &b.RealmID, &b.Name, &b.BlockType, &b.Source, &b.ConfigSchema, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: user block %q", ErrNotFound, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: getting user block %q: %w", name, err)
	}
	if len(b.ConfigSchema) == 0 {
		b.ConfigSchema = nil
	}
	return &b, nil
}

// UpdateUserBlock replaces a user block's type, source and schema; an unknown
// name yields ErrNotFound.
func (s *Store) UpdateUserBlock(ctx context.Context, realmID int16, name, blockType string, source, configSchema []byte) (*UserBlock, error) {
	var b UserBlock
	err := s.pool.QueryRow(ctx,
		`UPDATE user_blocks SET block_type = $3, source = $4, config_schema = $5, updated_at = now()
		 WHERE realm_id = $1 AND name = $2
		 RETURNING id, realm_id, name, block_type, source, config_schema, created_at, updated_at`,
		realmID, name, blockType, source, configSchema).
		Scan(&b.ID, &b.RealmID, &b.Name, &b.BlockType, &b.Source, &b.ConfigSchema, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: user block %q", ErrNotFound, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: updating user block %q: %w", name, err)
	}
	if len(b.ConfigSchema) == 0 {
		b.ConfigSchema = nil
	}
	return &b, nil
}

// DeleteUserBlock removes a user block; an unknown name yields ErrNotFound.
func (s *Store) DeleteUserBlock(ctx context.Context, realmID int16, name string) error {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM user_blocks WHERE realm_id = $1 AND name = $2`, realmID, name)
	if err != nil {
		return fmt.Errorf("store: deleting user block %q: %w", name, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: user block %q", ErrNotFound, name)
	}
	return nil
}

// ListUserBlocks returns every user block of a realm ordered by name.
func (s *Store) ListUserBlocks(ctx context.Context, realmID int16) ([]UserBlock, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, realm_id, name, block_type, source, config_schema, created_at, updated_at
		 FROM user_blocks WHERE realm_id = $1 ORDER BY name`,
		realmID)
	if err != nil {
		return nil, fmt.Errorf("store: listing user blocks: %w", err)
	}
	defer rows.Close()

	var out []UserBlock
	for rows.Next() {
		var b UserBlock
		if err := rows.Scan(&b.ID, &b.RealmID, &b.Name, &b.BlockType, &b.Source, &b.ConfigSchema, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("store: scanning user block: %w", err)
		}
		if len(b.ConfigSchema) == 0 {
			b.ConfigSchema = nil
		}
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: listing user blocks: %w", err)
	}
	return out, nil
}
