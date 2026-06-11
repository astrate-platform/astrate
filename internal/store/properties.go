package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// Property is one last-value-wins property row (docs/DESIGN.md §2.3). Value
// is the jsonb rendering; ValueType is the endpoint's declared type, kept on
// the row so the API layer re-encodes precisely (longinteger as decimal
// string, binaryblob as base64, datetime as RFC 3339).
type Property struct {
	RealmID     int16
	DeviceID    deviceid.ID
	InterfaceID int64
	EndpointID  int64
	Path        string
	Value       []byte
	ValueType   interfaceschema.ValueType
	SetAt       time.Time
}

// PropertyRef addresses one property of a device: the installed interface
// plus the concrete path.
type PropertyRef struct {
	InterfaceID int64
	Path        string
}

const propertyColumns = `realm_id, device_id, interface_id, endpoint_id, path, value, value_type, set_at`

// propertyColumnsQualified disambiguates the column list when properties is
// joined against interfaces (both carry realm_id).
const propertyColumnsQualified = `p.realm_id, p.device_id, p.interface_id, p.endpoint_id, p.path, p.value, p.value_type, p.set_at`

// UpsertProperty inserts or replaces a property value (last-value-wins).
func (s *Store) UpsertProperty(ctx context.Context, p Property) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO properties (realm_id, device_id, interface_id, endpoint_id, path, value, value_type, set_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, now())
		ON CONFLICT (realm_id, device_id, interface_id, path) DO UPDATE
			SET endpoint_id = EXCLUDED.endpoint_id,
			    value = EXCLUDED.value,
			    value_type = EXCLUDED.value_type,
			    set_at = now()`,
		p.RealmID, uuidParam(p.DeviceID), p.InterfaceID, p.EndpointID, p.Path, p.Value, p.ValueType.String())
	if err != nil {
		return fmt.Errorf("store: upserting property %s: %w", p.Path, err)
	}
	return nil
}

// UnsetProperty deletes a property row (empty payload with allow_unset,
// docs/DESIGN.md §2.3). It reports whether a row existed.
func (s *Store) UnsetProperty(ctx context.Context, realmID int16, deviceID deviceid.ID, interfaceID int64, path string) (bool, error) {
	tag, err := s.pool.Exec(ctx, `
		DELETE FROM properties
		WHERE realm_id = $1 AND device_id = $2 AND interface_id = $3 AND path = $4`,
		realmID, uuidParam(deviceID), interfaceID, path)
	if err != nil {
		return false, fmt.Errorf("store: unsetting property %s: %w", path, err)
	}
	return tag.RowsAffected() > 0, nil
}

// GetProperty fetches one property.
func (s *Store) GetProperty(ctx context.Context, realmID int16, deviceID deviceid.ID, interfaceID int64, path string) (*Property, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT `+propertyColumns+` FROM properties
		WHERE realm_id = $1 AND device_id = $2 AND interface_id = $3 AND path = $4`,
		realmID, uuidParam(deviceID), interfaceID, path)
	p, err := scanProperty(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: property %s", ErrNotFound, path)
	}
	if err != nil {
		return nil, fmt.Errorf("store: getting property %s: %w", path, err)
	}
	return p, nil
}

// ListProperties returns all properties of one device interface, ordered by
// path.
func (s *Store) ListProperties(ctx context.Context, realmID int16, deviceID deviceid.ID, interfaceID int64) ([]Property, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT `+propertyColumns+` FROM properties
		WHERE realm_id = $1 AND device_id = $2 AND interface_id = $3
		ORDER BY path`,
		realmID, uuidParam(deviceID), interfaceID)
	if err != nil {
		return nil, fmt.Errorf("store: listing properties: %w", err)
	}
	return collectProperties(rows)
}

// PurgeDeviceOwnedExcept deletes every device-owned property of the device
// that is not in keep — the `/control/producer/properties` resync
// (docs/DESIGN.md §3.3): the device sends the exhaustive list of properties
// it still holds and the server drops the complement. Server-owned
// properties are never touched. It returns the number of rows deleted.
func (s *Store) PurgeDeviceOwnedExcept(ctx context.Context, realmID int16, deviceID deviceid.ID, keep []PropertyRef) (int64, error) {
	ifaceIDs := make([]int64, len(keep))
	paths := make([]string, len(keep))
	for i, ref := range keep {
		ifaceIDs[i] = ref.InterfaceID
		paths[i] = ref.Path
	}
	tag, err := s.pool.Exec(ctx, `
		DELETE FROM properties p
		USING interfaces i
		WHERE i.id = p.interface_id
		  AND i.ownership = 'device'
		  AND p.realm_id = $1 AND p.device_id = $2
		  AND NOT EXISTS (
			SELECT 1 FROM unnest($3::bigint[], $4::text[]) AS k(interface_id, path)
			WHERE k.interface_id = p.interface_id AND k.path = p.path
		  )`,
		realmID, uuidParam(deviceID), ifaceIDs, paths)
	if err != nil {
		return 0, fmt.Errorf("store: purging device-owned properties: %w", err)
	}
	return tag.RowsAffected(), nil
}

// ListServerOwnedProperties returns every property of the device whose
// interface is server-owned — the payload of the
// `/control/consumer/properties` resync message (docs/DESIGN.md §3.4).
func (s *Store) ListServerOwnedProperties(ctx context.Context, realmID int16, deviceID deviceid.ID) ([]Property, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT `+propertyColumnsQualified+` FROM properties p
		JOIN interfaces i ON i.id = p.interface_id
		WHERE i.ownership = 'server' AND p.realm_id = $1 AND p.device_id = $2
		ORDER BY p.interface_id, p.path`,
		realmID, uuidParam(deviceID))
	if err != nil {
		return nil, fmt.Errorf("store: listing server-owned properties: %w", err)
	}
	return collectProperties(rows)
}

// collectProperties drains rows into a slice, closing them.
func collectProperties(rows pgx.Rows) ([]Property, error) {
	defer rows.Close()
	var out []Property
	for rows.Next() {
		p, err := scanProperty(rows)
		if err != nil {
			return nil, fmt.Errorf("store: scanning property: %w", err)
		}
		out = append(out, *p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: reading properties: %w", err)
	}
	return out, nil
}

// scanProperty reads one property row in propertyColumns order.
func scanProperty(row pgx.Row) (*Property, error) {
	var (
		p  Property
		u  pgtype.UUID
		vt string
	)
	if err := row.Scan(&p.RealmID, &u, &p.InterfaceID, &p.EndpointID, &p.Path,
		&p.Value, &vt, &p.SetAt); err != nil {
		return nil, err
	}
	p.DeviceID = deviceid.ID(u.Bytes)
	var err error
	if p.ValueType, err = interfaceschema.ParseValueType(vt); err != nil {
		return nil, fmt.Errorf("store: property %s: %w", p.Path, err)
	}
	return &p, nil
}
