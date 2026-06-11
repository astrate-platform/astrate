package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// Device status values (devices.status column, docs/DESIGN.md §2.2).
const (
	// DeviceStatusRegistered marks a device that has a credentials secret
	// but has never requested credentials.
	DeviceStatusRegistered = "registered"
	// DeviceStatusConfirmed marks a device that has requested credentials
	// at least once.
	DeviceStatusConfirmed = "confirmed"
	// DeviceStatusInhibited marks a device blocked from new credentials and
	// new connections (credentials_inhibited parity).
	DeviceStatusInhibited = "inhibited"
)

// InterfaceVersion is one introspection entry value:
// {"iface.Name": {"major": 1, "minor": 2}}.
type InterfaceVersion struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

// Device is one devices row (docs/DESIGN.md §2.2).
type Device struct {
	ID                       deviceid.ID
	RealmID                  int16
	CredentialsSecretHash    string
	Status                   string
	Introspection            map[string]InterfaceVersion
	OldIntrospection         map[string]InterfaceVersion
	Aliases                  map[string]string
	Attributes               map[string]string
	CertSerial               *string
	CertAKI                  *string
	FirstRegistration        time.Time
	FirstCredentialsRequest  *time.Time
	LastCredentialsRequestIP *netip.Addr
	LastConnection           *time.Time
	LastDisconnection        *time.Time
	LastSeenIP               *netip.Addr
	Connected                bool
	TotalReceivedMsgs        int64
	TotalReceivedBytes       int64
	PayloadFormatHint        string
}

const deviceColumns = `id, realm_id, credentials_secret_hash, status, introspection, old_introspection,
	aliases, attributes, cert_serial, cert_aki, first_registration, first_credentials_request,
	last_credentials_request_ip, last_connection, last_disconnection, last_seen_ip, connected,
	total_received_msgs, total_received_bytes, payload_format_hint`

// uuidParam adapts a device ID to the explicit pgx uuid type.
func uuidParam(id deviceid.ID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

// RegisterDevice inserts a device row, or rotates the credentials secret of
// an existing one that has not yet requested credentials (docs/DESIGN.md
// §4.4 flow A). Once the device has requested credentials, re-registration
// fails with ErrDeviceAlreadyConfirmed.
func (s *Store) RegisterDevice(ctx context.Context, realmID int16, id deviceid.ID, secretHash string) error {
	tag, err := s.pool.Exec(ctx, `
		INSERT INTO devices (id, realm_id, credentials_secret_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (realm_id, id) DO UPDATE
			SET credentials_secret_hash = EXCLUDED.credentials_secret_hash,
			    status = 'registered'
			WHERE devices.first_credentials_request IS NULL`,
		uuidParam(id), realmID, secretHash)
	if err != nil {
		return fmt.Errorf("store: registering device %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: device %s", ErrDeviceAlreadyConfirmed, id)
	}
	return nil
}

// UnregisterDevice makes a device registrable again without losing its data
// (DELETE /agent/devices parity, docs/DESIGN.md §4.4): the credentials
// secret and certificate trail are cleared, the row and all stored data stay.
func (s *Store) UnregisterDevice(ctx context.Context, realmID int16, id deviceid.ID) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE devices
		SET credentials_secret_hash = '', status = 'registered',
		    first_credentials_request = NULL, cert_serial = NULL, cert_aki = NULL
		WHERE realm_id = $1 AND id = $2`,
		realmID, uuidParam(id))
	if err != nil {
		return fmt.Errorf("store: unregistering device %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: device %s", ErrNotFound, id)
	}
	return nil
}

// GetDevice fetches one device.
func (s *Store) GetDevice(ctx context.Context, realmID int16, id deviceid.ID) (*Device, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT `+deviceColumns+` FROM devices WHERE realm_id = $1 AND id = $2`,
		realmID, uuidParam(id))
	d, err := scanDevice(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: device %s", ErrNotFound, id)
	}
	if err != nil {
		return nil, fmt.Errorf("store: getting device %s: %w", id, err)
	}
	return d, nil
}

// GetDeviceByAlias fetches a device by alias value (any alias tag). If
// several devices share an alias the lowest device ID wins; the API layer is
// expected to keep aliases unique per realm.
func (s *Store) GetDeviceByAlias(ctx context.Context, realmID int16, alias string) (*Device, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT `+deviceColumns+`
		FROM devices
		WHERE realm_id = $1
		  AND jsonb_path_exists(aliases, '$.* ? (@ == $alias)', jsonb_build_object('alias', $2::text))
		ORDER BY id
		LIMIT 1`,
		realmID, alias)
	d, err := scanDevice(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: device alias %q", ErrNotFound, alias)
	}
	if err != nil {
		return nil, fmt.Errorf("store: getting device by alias %q: %w", alias, err)
	}
	return d, nil
}

// ListDevices returns up to limit devices ordered by ID, starting after the
// optional keyset cursor. Keyset pagination keeps cursors stable while
// devices are inserted concurrently.
func (s *Store) ListDevices(ctx context.Context, realmID int16, after *deviceid.ID, limit int) ([]Device, error) {
	var (
		rows pgx.Rows
		err  error
	)
	if after != nil {
		rows, err = s.pool.Query(ctx,
			`SELECT `+deviceColumns+` FROM devices WHERE realm_id = $1 AND id > $2 ORDER BY id LIMIT $3`,
			realmID, uuidParam(*after), limit)
	} else {
		rows, err = s.pool.Query(ctx,
			`SELECT `+deviceColumns+` FROM devices WHERE realm_id = $1 ORDER BY id LIMIT $2`,
			realmID, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("store: listing devices: %w", err)
	}
	defer rows.Close()

	var out []Device
	for rows.Next() {
		d, err := scanDevice(rows)
		if err != nil {
			return nil, fmt.Errorf("store: scanning device: %w", err)
		}
		out = append(out, *d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: listing devices: %w", err)
	}
	return out, nil
}

// CountDevices returns the number of devices in a realm (used to enforce
// device_registration_limit in the pairing layer).
func (s *Store) CountDevices(ctx context.Context, realmID int16) (int64, error) {
	var n int64
	if err := s.pool.QueryRow(ctx,
		`SELECT count(*) FROM devices WHERE realm_id = $1`, realmID).Scan(&n); err != nil {
		return 0, fmt.Errorf("store: counting devices: %w", err)
	}
	return n, nil
}

// UpdateIntrospection replaces the device's introspection. Every
// (name, major) pair present in the current introspection but missing from
// the new one is merged into old_introspection (docs/ROADMAP.md §3.1 file
// 2.10) and returned, so the engine can react (cache eviction, property
// cleanup).
func (s *Store) UpdateIntrospection(ctx context.Context, realmID int16, id deviceid.ID, intro map[string]InterfaceVersion) (removed map[string]InterfaceVersion, err error) {
	if intro == nil {
		intro = map[string]InterfaceVersion{}
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("store: beginning introspection update: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var current, old map[string]InterfaceVersion
	err = tx.QueryRow(ctx, `
		SELECT introspection, old_introspection FROM devices
		WHERE realm_id = $1 AND id = $2
		FOR UPDATE`,
		realmID, uuidParam(id)).Scan(&current, &old)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: device %s", ErrNotFound, id)
	}
	if err != nil {
		return nil, fmt.Errorf("store: locking device %s: %w", id, err)
	}

	removed = map[string]InterfaceVersion{}
	if old == nil {
		old = map[string]InterfaceVersion{}
	}
	for name, v := range current {
		if nv, ok := intro[name]; !ok || nv.Major != v.Major {
			removed[name] = v
			old[name] = v
		}
	}

	if _, err := tx.Exec(ctx, `
		UPDATE devices SET introspection = $3, old_introspection = $4
		WHERE realm_id = $1 AND id = $2`,
		realmID, uuidParam(id), intro, old); err != nil {
		return nil, fmt.Errorf("store: updating introspection of %s: %w", id, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("store: committing introspection update: %w", err)
	}
	return removed, nil
}

// SetDeviceInhibited toggles the inhibit flag (credentials_inhibited PATCH
// parity): true forces status 'inhibited'; false restores 'confirmed' or
// 'registered' depending on whether credentials were ever requested.
func (s *Store) SetDeviceInhibited(ctx context.Context, realmID int16, id deviceid.ID, inhibited bool) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE devices
		SET status = CASE
			WHEN $3 THEN 'inhibited'
			WHEN first_credentials_request IS NULL THEN 'registered'
			ELSE 'confirmed'
		END
		WHERE realm_id = $1 AND id = $2`,
		realmID, uuidParam(id), inhibited)
	if err != nil {
		return fmt.Errorf("store: setting inhibit flag of %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: device %s", ErrNotFound, id)
	}
	return nil
}

// SetDeviceCredentials stamps a freshly issued client certificate on the
// device row (serial + authority key identifier, docs/DESIGN.md §4.3
// latest-serial enforcement) and records the credentials request.
func (s *Store) SetDeviceCredentials(ctx context.Context, realmID int16, id deviceid.ID, certSerial, certAKI string, requestIP netip.Addr) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE devices
		SET cert_serial = $3, cert_aki = $4,
		    first_credentials_request = coalesce(first_credentials_request, now()),
		    last_credentials_request_ip = $5,
		    status = 'confirmed'
		WHERE realm_id = $1 AND id = $2`,
		realmID, uuidParam(id), certSerial, certAKI, requestIP)
	if err != nil {
		return fmt.Errorf("store: stamping credentials of %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: device %s", ErrNotFound, id)
	}
	return nil
}

// SetDeviceConnected records a broker connection.
func (s *Store) SetDeviceConnected(ctx context.Context, realmID int16, id deviceid.ID, at time.Time, ip netip.Addr) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE devices
		SET connected = true, last_connection = $3, last_seen_ip = $4
		WHERE realm_id = $1 AND id = $2`,
		realmID, uuidParam(id), at, ip)
	if err != nil {
		return fmt.Errorf("store: marking %s connected: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: device %s", ErrNotFound, id)
	}
	return nil
}

// SetDeviceDisconnected records a broker disconnection.
func (s *Store) SetDeviceDisconnected(ctx context.Context, realmID int16, id deviceid.ID, at time.Time) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE devices
		SET connected = false, last_disconnection = $3
		WHERE realm_id = $1 AND id = $2`,
		realmID, uuidParam(id), at)
	if err != nil {
		return fmt.Errorf("store: marking %s disconnected: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: device %s", ErrNotFound, id)
	}
	return nil
}

// AddDeviceStats increments the received message/byte counters.
func (s *Store) AddDeviceStats(ctx context.Context, realmID int16, id deviceid.ID, msgs, bytes int64) error {
	tag, err := s.pool.Exec(ctx, `
		UPDATE devices
		SET total_received_msgs = total_received_msgs + $3,
		    total_received_bytes = total_received_bytes + $4
		WHERE realm_id = $1 AND id = $2`,
		realmID, uuidParam(id), msgs, bytes)
	if err != nil {
		return fmt.Errorf("store: updating stats of %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: device %s", ErrNotFound, id)
	}
	return nil
}

// PatchDeviceAliases merges patch into the alias map: non-nil values
// add/replace the tag, nil values remove it (JSON Merge Patch semantics).
func (s *Store) PatchDeviceAliases(ctx context.Context, realmID int16, id deviceid.ID, patch map[string]*string) error {
	return s.patchDeviceJSONB(ctx, "aliases", realmID, id, patch)
}

// PatchDeviceAttributes merges patch into the attributes map with the same
// semantics as PatchDeviceAliases.
func (s *Store) PatchDeviceAttributes(ctx context.Context, realmID int16, id deviceid.ID, patch map[string]*string) error {
	return s.patchDeviceJSONB(ctx, "attributes", realmID, id, patch)
}

// patchDeviceJSONB applies merge-patch semantics to one of the device's
// text→text jsonb columns. column is always a compile-time constant.
func (s *Store) patchDeviceJSONB(ctx context.Context, column string, realmID int16, id deviceid.ID, patch map[string]*string) error {
	removals := []string{}
	additions := map[string]string{}
	for k, v := range patch {
		if v == nil {
			removals = append(removals, k)
		} else {
			additions[k] = *v
		}
	}
	add, err := json.Marshal(additions)
	if err != nil {
		return fmt.Errorf("store: encoding %s patch: %w", column, err)
	}

	tag, err := s.pool.Exec(ctx, `
		UPDATE devices
		SET `+column+` = (`+column+` - $3::text[]) || $4::jsonb
		WHERE realm_id = $1 AND id = $2`,
		realmID, uuidParam(id), removals, add)
	if err != nil {
		return fmt.Errorf("store: patching %s of %s: %w", column, id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: device %s", ErrNotFound, id)
	}
	return nil
}

// SetPayloadFormatHint flips the device's preferred outbound payload format
// ("bson" or "json", docs/DESIGN.md §3.5.4).
func (s *Store) SetPayloadFormatHint(ctx context.Context, realmID int16, id deviceid.ID, hint string) error {
	if hint != "bson" && hint != "json" {
		return fmt.Errorf("store: invalid payload format hint %q", hint)
	}
	tag, err := s.pool.Exec(ctx, `
		UPDATE devices SET payload_format_hint = $3 WHERE realm_id = $1 AND id = $2`,
		realmID, uuidParam(id), hint)
	if err != nil {
		return fmt.Errorf("store: setting payload format hint of %s: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: device %s", ErrNotFound, id)
	}
	return nil
}

// scanDevice reads one device row in deviceColumns order.
func scanDevice(row pgx.Row) (*Device, error) {
	var (
		d Device
		u pgtype.UUID
	)
	if err := row.Scan(&u, &d.RealmID, &d.CredentialsSecretHash, &d.Status,
		&d.Introspection, &d.OldIntrospection, &d.Aliases, &d.Attributes,
		&d.CertSerial, &d.CertAKI, &d.FirstRegistration, &d.FirstCredentialsRequest,
		&d.LastCredentialsRequestIP, &d.LastConnection, &d.LastDisconnection,
		&d.LastSeenIP, &d.Connected, &d.TotalReceivedMsgs, &d.TotalReceivedBytes,
		&d.PayloadFormatHint); err != nil {
		return nil, err
	}
	d.ID = deviceid.ID(u.Bytes)
	return &d, nil
}
