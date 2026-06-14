package store

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// IndividualRow is one individual_datastreams row (docs/DESIGN.md §2.4):
// exactly one Value* field must be set, matching the endpoint's declared
// type. ValueArray and the object row Value carry pre-encoded JSON.
type IndividualRow struct {
	RealmID     int16
	DeviceID    deviceid.ID
	InterfaceID int64
	EndpointID  int64
	Path        string
	TS          time.Time
	ReceptionTS time.Time

	ValueDouble      *float64
	ValueInteger     *int32
	ValueLonginteger *int64
	ValueBoolean     *bool
	ValueString      *string
	ValueBinaryblob  []byte
	ValueDatetime    *time.Time
	ValueArray       []byte
}

// validate enforces the exactly-one-value-column invariant.
func (r *IndividualRow) validate() error {
	n := 0
	if r.ValueDouble != nil {
		n++
	}
	if r.ValueInteger != nil {
		n++
	}
	if r.ValueLonginteger != nil {
		n++
	}
	if r.ValueBoolean != nil {
		n++
	}
	if r.ValueString != nil {
		n++
	}
	if r.ValueBinaryblob != nil {
		n++
	}
	if r.ValueDatetime != nil {
		n++
	}
	if r.ValueArray != nil {
		n++
	}
	if n != 1 {
		return fmt.Errorf("store: datastream row %s@%s sets %d value columns, want exactly 1", r.Path, r.TS.Format(time.RFC3339), n)
	}
	return nil
}

// ObjectRow is one object_datastreams row: an object-aggregated publish on a
// path prefix, with the last-level keys as one JSON document.
type ObjectRow struct {
	RealmID     int16
	DeviceID    deviceid.ID
	InterfaceID int64
	Path        string
	TS          time.Time
	ReceptionTS time.Time
	Value       []byte
}

// DatastreamBatch is one persistence flush: the engine's per-shard
// micro-batches mix individual and object rows (docs/DESIGN.md §1.4).
type DatastreamBatch struct {
	Individual []IndividualRow
	Objects    []ObjectRow
}

var (
	individualColumns = []string{
		"realm_id", "device_id", "interface_id", "endpoint_id", "path", "ts", "reception_ts",
		"value_double", "value_integer", "value_longinteger", "value_boolean",
		"value_string", "value_binaryblob", "value_datetime", "value_array",
	}
	objectColumns = []string{
		"realm_id", "device_id", "interface_id", "path", "ts", "reception_ts", "value",
	}
)

// AppendDatastreams persists a batch through binary COPY, both tables in one
// transaction (docs/DESIGN.md §5.3: the broker PUBACKs only after this
// commits). The hypertables have no unique constraints, so duplicate
// (series, ts) rows from at-least-once redelivery are tolerated by design.
func (s *Store) AppendDatastreams(ctx context.Context, batch DatastreamBatch) error {
	for i := range batch.Individual {
		if err := batch.Individual[i].validate(); err != nil {
			return err
		}
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("store: beginning datastream append: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if len(batch.Individual) > 0 {
		_, err := tx.CopyFrom(ctx, pgx.Identifier{"individual_datastreams"}, individualColumns,
			pgx.CopyFromSlice(len(batch.Individual), func(i int) ([]any, error) {
				r := &batch.Individual[i]
				return []any{
					r.RealmID, uuidParam(r.DeviceID), r.InterfaceID, r.EndpointID, r.Path,
					r.TS, r.ReceptionTS,
					r.ValueDouble, r.ValueInteger, r.ValueLonginteger, r.ValueBoolean,
					r.ValueString, r.ValueBinaryblob, r.ValueDatetime, r.ValueArray,
				}, nil
			}))
		if err != nil {
			return fmt.Errorf("store: copying individual datastreams: %w", err)
		}
	}
	if len(batch.Objects) > 0 {
		_, err := tx.CopyFrom(ctx, pgx.Identifier{"object_datastreams"}, objectColumns,
			pgx.CopyFromSlice(len(batch.Objects), func(i int) ([]any, error) {
				r := &batch.Objects[i]
				return []any{
					r.RealmID, uuidParam(r.DeviceID), r.InterfaceID, r.Path,
					r.TS, r.ReceptionTS, r.Value,
				}, nil
			}))
		if err != nil {
			return fmt.Errorf("store: copying object datastreams: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("store: committing datastream append: %w", err)
	}
	return nil
}

// SeriesQuery addresses one series (a concrete path of one device interface)
// and the time window to read. Boundary semantics follow AppEngine parity:
// Since is inclusive (ts >= Since), SinceAfter is exclusive (ts > SinceAfter),
// To is inclusive (ts <= To). Limit 0 means no limit; Descending flips the
// ts ordering (used for "latest N" queries).
type SeriesQuery struct {
	RealmID     int16
	DeviceID    deviceid.ID
	InterfaceID int64
	Path        string
	Since       *time.Time
	SinceAfter  *time.Time
	To          *time.Time
	Limit       int
	Descending  bool
}

// where renders the WHERE clause and arguments shared by the series queries.
func (q *SeriesQuery) where() (string, []any) {
	var sb strings.Builder
	args := []any{q.RealmID, uuidParam(q.DeviceID), q.InterfaceID, q.Path}
	sb.WriteString("WHERE realm_id = $1 AND device_id = $2 AND interface_id = $3 AND path = $4")
	for _, c := range []struct {
		op string
		t  *time.Time
	}{
		{">=", q.Since},
		{">", q.SinceAfter},
		{"<=", q.To},
	} {
		if c.t != nil {
			args = append(args, *c.t)
			fmt.Fprintf(&sb, " AND ts %s $%d", c.op, len(args))
		}
	}
	return sb.String(), args
}

// tail renders the ORDER BY / LIMIT suffix.
func (q *SeriesQuery) tail() string {
	order := " ORDER BY ts ASC"
	if q.Descending {
		order = " ORDER BY ts DESC"
	}
	if q.Limit > 0 {
		order += " LIMIT " + strconv.Itoa(q.Limit)
	}
	return order
}

// individualSelect is the column list shared by the individual-datastream
// read queries; scanIndividualRows consumes rows in this exact order.
const individualSelect = `realm_id, device_id, interface_id, endpoint_id, path, ts, reception_ts,
	       value_double, value_integer, value_longinteger, value_boolean,
	       value_string, value_binaryblob, value_datetime, value_array`

// scanIndividualRows drains rows selecting individualSelect into IndividualRows.
func scanIndividualRows(rows pgx.Rows) ([]IndividualRow, error) {
	defer rows.Close()
	var out []IndividualRow
	for rows.Next() {
		var (
			r IndividualRow
			u pgtype.UUID
		)
		if err := rows.Scan(&r.RealmID, &u, &r.InterfaceID, &r.EndpointID, &r.Path,
			&r.TS, &r.ReceptionTS,
			&r.ValueDouble, &r.ValueInteger, &r.ValueLonginteger, &r.ValueBoolean,
			&r.ValueString, &r.ValueBinaryblob, &r.ValueDatetime, &r.ValueArray); err != nil {
			return nil, fmt.Errorf("store: scanning datastream row: %w", err)
		}
		r.DeviceID = deviceid.ID(u.Bytes)
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: reading datastream rows: %w", err)
	}
	return out, nil
}

// Series reads one individual-datastream series.
func (s *Store) Series(ctx context.Context, q SeriesQuery) ([]IndividualRow, error) {
	where, args := q.where()
	rows, err := s.pool.Query(ctx,
		`SELECT `+individualSelect+` FROM individual_datastreams `+where+q.tail(), args...)
	if err != nil {
		return nil, fmt.Errorf("store: querying series %s: %w", q.Path, err)
	}
	return scanIndividualRows(rows)
}

// IndividualSnapshot returns the most recent sample for each distinct path of
// an individual-datastream interface — the AppEngine interface-root snapshot
// ("data-snapshot") view upstream renders as a nested tree. Paths that never
// received a sample are absent; the result is ordered by path.
func (s *Store) IndividualSnapshot(ctx context.Context, realmID int16, deviceID deviceid.ID, interfaceID int64) ([]IndividualRow, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT DISTINCT ON (path) `+individualSelect+`
		 FROM individual_datastreams
		 WHERE realm_id = $1 AND device_id = $2 AND interface_id = $3
		 ORDER BY path, ts DESC`,
		realmID, uuidParam(deviceID), interfaceID)
	if err != nil {
		return nil, fmt.Errorf("store: querying individual snapshot: %w", err)
	}
	return scanIndividualRows(rows)
}

// ObjectSeries reads one object-datastream series (the path is the
// aggregation prefix).
func (s *Store) ObjectSeries(ctx context.Context, q SeriesQuery) ([]ObjectRow, error) {
	where, args := q.where()
	rows, err := s.pool.Query(ctx, `
		SELECT realm_id, device_id, interface_id, path, ts, reception_ts, value
		FROM object_datastreams `+where+q.tail(), args...)
	if err != nil {
		return nil, fmt.Errorf("store: querying object series %s: %w", q.Path, err)
	}
	defer rows.Close()

	var out []ObjectRow
	for rows.Next() {
		var (
			r ObjectRow
			u pgtype.UUID
		)
		if err := rows.Scan(&r.RealmID, &u, &r.InterfaceID, &r.Path, &r.TS, &r.ReceptionTS, &r.Value); err != nil {
			return nil, fmt.Errorf("store: scanning object series row: %w", err)
		}
		r.DeviceID = deviceid.ID(u.Bytes)
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: reading object series %s: %w", q.Path, err)
	}
	return out, nil
}

// DownsamplePoint is one downsampled bucket: the bucket start time and the
// aggregated numeric value.
type DownsamplePoint struct {
	Bucket time.Time
	Value  float64
}

// Downsample reduces a numeric individual-datastream series to one averaged
// point per bucket via TimescaleDB time_bucket (AppEngine downsample_to,
// docs/DESIGN.md §2.5). Non-numeric rows in the window are ignored.
//
// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5):
// when s.hasToolkit is set (timescaledb_toolkit probed at startup), switch to
// toolkit lttb() for shape-preserving downsampling; this time_bucket+avg
// path is the always-available default and stays as the fallback.
func (s *Store) Downsample(ctx context.Context, q SeriesQuery, bucket time.Duration) ([]DownsamplePoint, error) {
	if bucket <= 0 {
		return nil, fmt.Errorf("store: downsample bucket must be positive, got %s", bucket)
	}
	where, args := q.where()
	args = append(args, bucket.Seconds())
	order := " ORDER BY bucket ASC"
	if q.Descending {
		order = " ORDER BY bucket DESC"
	}
	limit := ""
	if q.Limit > 0 {
		limit = " LIMIT " + strconv.Itoa(q.Limit)
	}

	rows, err := s.pool.Query(ctx, fmt.Sprintf(`
		SELECT time_bucket(make_interval(secs => $%d), ts) AS bucket,
		       avg(coalesce(value_double, value_integer::double precision, value_longinteger::double precision)) AS value
		FROM individual_datastreams
		%s
		  AND (value_double IS NOT NULL OR value_integer IS NOT NULL OR value_longinteger IS NOT NULL)
		GROUP BY bucket`+order+limit, len(args), where), args...)
	if err != nil {
		return nil, fmt.Errorf("store: downsampling series %s: %w", q.Path, err)
	}
	defer rows.Close()

	var out []DownsamplePoint
	for rows.Next() {
		var p DownsamplePoint
		if err := rows.Scan(&p.Bucket, &p.Value); err != nil {
			return nil, fmt.Errorf("store: scanning downsample row: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: reading downsample of %s: %w", q.Path, err)
	}
	return out, nil
}
