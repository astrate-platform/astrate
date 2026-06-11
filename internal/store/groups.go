package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// Group is one device group (docs/DESIGN.md §2.2). Membership rows carry the
// composite (realm_id, device_id) foreign key so a device removal cascades
// out of its groups automatically.
type Group struct {
	ID      int64
	RealmID int16
	Name    string
}

// CreateGroup inserts a group; a duplicate name yields ErrAlreadyExists.
func (s *Store) CreateGroup(ctx context.Context, realmID int16, name string) (*Group, error) {
	var g Group
	g.RealmID = realmID
	g.Name = name
	err := s.pool.QueryRow(ctx,
		`INSERT INTO groups (realm_id, name) VALUES ($1, $2) RETURNING id`,
		realmID, name).Scan(&g.ID)
	if pgErrCode(err) == pgCodeUniqueViolation {
		return nil, fmt.Errorf("%w: group %q", ErrAlreadyExists, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: creating group %q: %w", name, err)
	}
	return &g, nil
}

// GetGroupByName fetches one group.
func (s *Store) GetGroupByName(ctx context.Context, realmID int16, name string) (*Group, error) {
	var g Group
	err := s.pool.QueryRow(ctx,
		`SELECT id, realm_id, name FROM groups WHERE realm_id = $1 AND name = $2`,
		realmID, name).Scan(&g.ID, &g.RealmID, &g.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: group %q", ErrNotFound, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: getting group %q: %w", name, err)
	}
	return &g, nil
}

// DeleteGroup removes a group; membership rows cascade.
func (s *Store) DeleteGroup(ctx context.Context, realmID int16, name string) error {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM groups WHERE realm_id = $1 AND name = $2`, realmID, name)
	if err != nil {
		return fmt.Errorf("store: deleting group %q: %w", name, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: group %q", ErrNotFound, name)
	}
	return nil
}

// ListGroups returns every group of a realm ordered by name.
func (s *Store) ListGroups(ctx context.Context, realmID int16) ([]Group, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, realm_id, name FROM groups WHERE realm_id = $1 ORDER BY name`, realmID)
	if err != nil {
		return nil, fmt.Errorf("store: listing groups: %w", err)
	}
	defer rows.Close()

	var out []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.RealmID, &g.Name); err != nil {
			return nil, fmt.Errorf("store: scanning group: %w", err)
		}
		out = append(out, g)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: listing groups: %w", err)
	}
	return out, nil
}

// AddGroupDevice adds a device to a group. The INSERT...SELECT guard ensures
// the group belongs to realmID, and the composite foreign key validates the
// device; an unknown group or device yields ErrNotFound, an existing
// membership ErrAlreadyExists.
func (s *Store) AddGroupDevice(ctx context.Context, groupID int64, realmID int16, deviceID deviceid.ID) error {
	tag, err := s.pool.Exec(ctx, `
		INSERT INTO group_devices (group_id, realm_id, device_id)
		SELECT g.id, $2, $3 FROM groups g WHERE g.id = $1 AND g.realm_id = $2`,
		groupID, realmID, uuidParam(deviceID))
	switch pgErrCode(err) {
	case pgCodeUniqueViolation:
		return fmt.Errorf("%w: device %s in group %d", ErrAlreadyExists, deviceID, groupID)
	case pgCodeFKViolation:
		return fmt.Errorf("%w: device %s", ErrNotFound, deviceID)
	}
	if err != nil {
		return fmt.Errorf("store: adding device %s to group %d: %w", deviceID, groupID, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: group %d", ErrNotFound, groupID)
	}
	return nil
}

// RemoveGroupDevice removes a device from a group.
func (s *Store) RemoveGroupDevice(ctx context.Context, groupID int64, realmID int16, deviceID deviceid.ID) error {
	tag, err := s.pool.Exec(ctx, `
		DELETE FROM group_devices WHERE group_id = $1 AND realm_id = $2 AND device_id = $3`,
		groupID, realmID, uuidParam(deviceID))
	if err != nil {
		return fmt.Errorf("store: removing device %s from group %d: %w", deviceID, groupID, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: device %s in group %d", ErrNotFound, deviceID, groupID)
	}
	return nil
}

// ListGroupDevices returns the IDs of every device in a group, ordered.
func (s *Store) ListGroupDevices(ctx context.Context, groupID int64) ([]deviceid.ID, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT device_id FROM group_devices WHERE group_id = $1 ORDER BY device_id`, groupID)
	if err != nil {
		return nil, fmt.Errorf("store: listing devices of group %d: %w", groupID, err)
	}
	defer rows.Close()

	var out []deviceid.ID
	for rows.Next() {
		var u pgtype.UUID
		if err := rows.Scan(&u); err != nil {
			return nil, fmt.Errorf("store: scanning group device: %w", err)
		}
		out = append(out, deviceid.ID(u.Bytes))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: listing devices of group %d: %w", groupID, err)
	}
	return out, nil
}

// ListDeviceGroups returns the names of every group the device belongs to.
func (s *Store) ListDeviceGroups(ctx context.Context, realmID int16, deviceID deviceid.ID) ([]string, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT g.name FROM groups g
		JOIN group_devices gd ON gd.group_id = g.id
		WHERE gd.realm_id = $1 AND gd.device_id = $2
		ORDER BY g.name`,
		realmID, uuidParam(deviceID))
	if err != nil {
		return nil, fmt.Errorf("store: listing groups of device %s: %w", deviceID, err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("store: scanning device group: %w", err)
		}
		out = append(out, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: listing groups of device %s: %w", deviceID, err)
	}
	return out, nil
}
