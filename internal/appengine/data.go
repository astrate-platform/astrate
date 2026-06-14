package appengine

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// datetimeLayout renders datetime values in query results (the §2.3 / §3.5.3
// JSON profile form: RFC 3339, UTC, millisecond precision).
const datetimeLayout = "2006-01-02T15:04:05.000Z"

// QueryOpts are the datastream query parameters (upstream
// since/since_after/to/limit/downsample_to). The zero value reads the whole
// series ascending.
type QueryOpts struct {
	Since        *time.Time
	SinceAfter   *time.Time
	To           *time.Time
	Limit        int
	Descending   bool
	DownsampleTo *time.Duration
}

// Sample is one datastream point rendered for the wire.
type Sample struct {
	Value     any       `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

// resolved is a device interface resolved against the device's introspection.
type resolved struct {
	rid   int16
	id    deviceid.ID
	iface *store.StoredInterface
}

// resolve maps (realm, device, interface name) to the installed interface the
// device declares, returning store.ErrNotFound when the device, the
// introspection entry, or the interface is missing.
func (s *Service) resolve(ctx context.Context, realm, deviceID, ifaceName string) (*resolved, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	id, err := deviceid.Parse(deviceID)
	if err != nil {
		return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
	}
	d, err := s.st.GetDevice(ctx, rid, id)
	if err != nil {
		return nil, err
	}
	v, ok := d.Introspection[ifaceName]
	if !ok {
		return nil, fmt.Errorf("%w: interface %s not in device introspection", store.ErrNotFound, ifaceName)
	}
	si, err := s.st.GetInterface(ctx, rid, ifaceName, v.Major)
	if err != nil {
		return nil, err
	}
	return &resolved{rid: rid, id: id, iface: si}, nil
}

// GetData reads an interface endpoint (upstream GET .../interfaces/{iface}
// [/{path}]), dispatching on the interface type: a datastream series (with the
// query options) or a properties snapshot.
func (s *Service) GetData(ctx context.Context, realm, deviceID, ifaceName, path string, opts QueryOpts) (any, error) {
	r, err := s.resolve(ctx, realm, deviceID, ifaceName)
	if err != nil {
		return nil, err
	}
	if r.iface.Type == interfaceschema.Properties {
		return s.propertiesData(ctx, r, path)
	}
	return s.datastreamData(ctx, r, path, opts)
}

// datastreamData reads a datastream endpoint. For an object-aggregated
// interface it returns the stored JSON document per sample; for individual it
// re-encodes the typed value per §2.3. A downsample_to opt reduces a numeric
// individual series to bucket averages.
func (s *Service) datastreamData(ctx context.Context, r *resolved, path string, opts QueryOpts) (any, error) {
	// Interface-root query (no path) on an individual datastream: the upstream
	// "data-snapshot" view — the latest sample for every endpoint, rendered as
	// a nested {segment: {... : {value, timestamp}}} tree (astarte-go walks it
	// via parseDatastreamMap, keyed on the "value" leaf field).
	if path == "" && opts.DownsampleTo == nil && r.iface.Aggregation != interfaceschema.AggregationObject {
		rows, err := s.st.IndividualSnapshot(ctx, r.rid, r.id, r.iface.ID)
		if err != nil {
			return nil, err
		}
		leaves := make(map[string]any, len(rows))
		for i := range rows {
			leaves[rows[i].Path] = Sample{Value: individualValue(&rows[i]), Timestamp: rows[i].TS}
		}
		return nestTree(leaves), nil
	}

	q := store.SeriesQuery{
		RealmID: r.rid, DeviceID: r.id, InterfaceID: r.iface.ID, Path: path,
		Since: opts.Since, SinceAfter: opts.SinceAfter, To: opts.To,
		Limit: opts.Limit, Descending: opts.Descending,
	}

	if opts.DownsampleTo != nil {
		points, err := s.st.Downsample(ctx, q, *opts.DownsampleTo)
		if err != nil {
			return nil, err
		}
		out := make([]Sample, len(points))
		for i := range points {
			out[i] = Sample{Value: points[i].Value, Timestamp: points[i].Bucket}
		}
		return out, nil
	}

	if r.iface.Aggregation == interfaceschema.AggregationObject {
		rows, err := s.st.ObjectSeries(ctx, q)
		if err != nil {
			return nil, err
		}
		out := make([]objectSample, len(rows))
		for i := range rows {
			out[i] = objectSample{Value: json.RawMessage(rows[i].Value), Timestamp: rows[i].TS}
		}
		return out, nil
	}

	rows, err := s.st.Series(ctx, q)
	if err != nil {
		return nil, err
	}
	out := make([]Sample, len(rows))
	for i := range rows {
		out[i] = Sample{Value: individualValue(&rows[i]), Timestamp: rows[i].TS}
	}
	return out, nil
}

// objectSample is one object-aggregated datastream point.
type objectSample struct {
	Value     json.RawMessage `json:"value"`
	Timestamp time.Time       `json:"timestamp"`
}

// propertiesData reads a properties interface. With a concrete path it returns
// the single value; at the interface root it returns the nested {segment: {...:
// value}} tree of every set property (the upstream interface snapshot shape,
// which astarte-go flattens via parsePropertiesMap). The stored jsonb already
// carries the §2.3 rendering.
func (s *Service) propertiesData(ctx context.Context, r *resolved, path string) (any, error) {
	if path != "" {
		p, err := s.st.GetProperty(ctx, r.rid, r.id, r.iface.ID, path)
		if err != nil {
			return nil, err
		}
		return json.RawMessage(p.Value), nil
	}
	props, err := s.st.ListProperties(ctx, r.rid, r.id, r.iface.ID)
	if err != nil {
		return nil, err
	}
	leaves := make(map[string]any, len(props))
	for i := range props {
		leaves[props[i].Path] = json.RawMessage(props[i].Value)
	}
	return nestTree(leaves), nil
}

// nestTree expands a flat map of Astarte endpoint paths ("/a/b") into the
// nested JSON object an AppEngine interface-root query returns: each "/"
// segment becomes a level and the leaf value is placed at the full path. An
// empty input yields an empty object.
func nestTree(leaves map[string]any) map[string]any {
	root := map[string]any{}
	for p, leaf := range leaves {
		segs := strings.Split(strings.TrimPrefix(p, "/"), "/")
		m := root
		for i, seg := range segs {
			if i == len(segs)-1 {
				m[seg] = leaf
				break
			}
			child, ok := m[seg].(map[string]any)
			if !ok {
				child = map[string]any{}
				m[seg] = child
			}
			m = child
		}
	}
	return root
}

// PublishData writes a server-owned value (upstream PUT/POST
// .../interfaces/{iface}/{path}) through the engine. value is the unwrapped
// "data" JSON; ts is the optional explicit timestamp.
func (s *Service) PublishData(ctx context.Context, realm, deviceID, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
	if s.sd == nil {
		return fmt.Errorf("appengine: server-owned writes are disabled (no engine)")
	}
	id, err := deviceid.Parse(deviceID)
	if err != nil {
		return fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
	}
	return s.sd.PublishServerValue(ctx, realm, id, ifaceName, path, value, ts)
}

// UnsetProperty unsets a server-owned property (upstream DELETE
// .../interfaces/{iface}/{path}).
func (s *Service) UnsetProperty(ctx context.Context, realm, deviceID, ifaceName, path string) error {
	if s.sd == nil {
		return fmt.Errorf("appengine: server-owned writes are disabled (no engine)")
	}
	id, err := deviceid.Parse(deviceID)
	if err != nil {
		return fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
	}
	return s.sd.UnsetServerProperty(ctx, realm, id, ifaceName, path)
}

// individualValue re-encodes one individual-datastream row's typed value into
// its §2.3 JSON-friendly form: longinteger as a decimal string, binaryblob as
// base64, datetime as RFC 3339, arrays passed through as stored jsonb.
func individualValue(r *store.IndividualRow) any {
	switch {
	case r.ValueDouble != nil:
		return *r.ValueDouble
	case r.ValueInteger != nil:
		return *r.ValueInteger
	case r.ValueLonginteger != nil:
		return strconv.FormatInt(*r.ValueLonginteger, 10)
	case r.ValueBoolean != nil:
		return *r.ValueBoolean
	case r.ValueString != nil:
		return *r.ValueString
	case r.ValueBinaryblob != nil:
		return base64.StdEncoding.EncodeToString(r.ValueBinaryblob)
	case r.ValueDatetime != nil:
		return r.ValueDatetime.UTC().Format(datetimeLayout)
	case r.ValueArray != nil:
		return json.RawMessage(r.ValueArray)
	default:
		return nil
	}
}
