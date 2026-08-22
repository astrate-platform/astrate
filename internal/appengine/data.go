package appengine

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
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

// maxSafeJSONInteger is upstream's safe-longinteger bound (JavaScript's
// Number.MAX_SAFE_INTEGER): with allow_safe_bigintegers a longinteger leaf
// stays a number only while |v| ≤ this, else it degrades to a decimal string.
const maxSafeJSONInteger int64 = 9007199254740991

// ErrPathNotFound marks an object-aggregated series whose concrete path holds
// no rows at all (upstream answers 404 "Path not found"; the interface root
// renders {} instead).
var ErrPathNotFound = errors.New("appengine: path not found")

// QueryOpts are the datastream query parameters (upstream since/since_after/
// to/limit/downsample_to plus the rendering options). The zero value reads
// the whole series ascending and renders structured.
//
// DownsamplePoints is the request's `downsample_to` value — a target *point
// count*, which is upstream's semantics. The bucket width is derived locally
// inside datastreamData once the series' own time span is known.
//
// Format is one of structured/table/disjoint_tables ("structured" when the
// request omitted it). AllowBigIntegers/AllowSafeBigIntegers are nil when the
// request did not name them; they affect object-document rendering only.
type QueryOpts struct {
	Since                *time.Time
	SinceAfter           *time.Time
	To                   *time.Time
	Limit                int
	Descending           bool
	DownsamplePoints     int
	Format               string
	DownsampleKey        string
	RetrieveMetadata     bool
	AllowBigIntegers     *bool
	AllowSafeBigIntegers *bool
}

// Tabular carries a format=table payload plus its metadata object; serveData
// detects it and renders {data, metadata} instead of bare data.
type Tabular struct {
	Data     any            `json:"-"`
	Metadata map[string]any `json:"metadata"`
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
		// A property has one current value, not a series over time. Ignoring the
		// parameter here would answer 200 with the full snapshot a client asked
		// to have reduced — the same silent lie datastreamData refuses below.
		if opts.DownsamplePoints > 0 {
			return nil, fmt.Errorf("%w: downsample_to is not supported on properties interfaces", ErrValidation)
		}
		return s.propertiesData(ctx, r, path)
	}
	return s.datastreamData(ctx, r, path, opts)
}

// datastreamData reads a datastream endpoint. For an object-aggregated
// interface it returns the stored JSON document per sample; for individual it
// re-encodes the typed value per §2.3. A downsample_to opt reduces a numeric
// individual series to bucket averages. The requested format shapes both
// branches, except where noted: the snapshot view and the downsampled view
// ignore it (documented deviation — upstream has no equivalent of either).
func (s *Service) datastreamData(ctx context.Context, r *resolved, path string, opts QueryOpts) (any, error) {
	// Interface-root query (no path) on an individual datastream: the upstream
	// "data-snapshot" view — the latest sample for every endpoint, rendered as
	// a nested {segment: {... : {value, timestamp}}} tree (astarte-go walks it
	// via parseDatastreamMap, keyed on the "value" leaf field). Format is
	// ignored here: there is no upstream table rendering of this tree to copy.
	if path == "" && opts.DownsamplePoints == 0 && r.iface.Aggregation != interfaceschema.AggregationObject {
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

	if opts.DownsamplePoints > 0 {
		if r.iface.Aggregation == interfaceschema.AggregationObject {
			return nil, fmt.Errorf("%w: downsample_to is not supported on object-aggregated interfaces", ErrValidation)
		}
		// An interface-root query is the snapshot view — the latest sample of
		// every endpoint — not a series, so there is nothing to reduce. Without
		// this the query falls through to a SeriesQuery whose empty path matches
		// no row and answers 200 + [], the same silent lie the object case above
		// exists to prevent.
		if path == "" {
			return nil, fmt.Errorf("%w: downsample_to requires an endpoint path", ErrValidation)
		}
		first, last, ok, err := s.st.SeriesSpan(ctx, q)
		if err != nil {
			return nil, err
		}
		if !ok {
			return []Sample{}, nil
		}
		// lttb rejects resolution < 3, but downsample_to permits 2 — fall back to bucket path.
		if s.st.HasToolkitLTTB() && opts.DownsamplePoints >= 3 {
			points, err := s.st.DownsampleLTTB(ctx, q, opts.DownsamplePoints)
			if err != nil {
				return nil, err
			}
			out := make([]Sample, len(points))
			for i := range points {
				out[i] = Sample{Value: points[i].Value, Timestamp: points[i].Bucket}
			}
			return out, nil
		}
		bucket := bucketFor(last.Sub(first), opts.DownsamplePoints)
		if bucket <= 0 {
			// A zero-span series (one sample or simultaneous samples) produces
			// a zero bucket; use the finest bucket that still collapses
			// simultaneous samples into one point.
			bucket = time.Microsecond
		}
		points, err := s.st.Downsample(ctx, q, bucket)
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
		// Upstream pack_result on an EMPTY object series: the interface root
		// renders {} while a concrete path is a 404 "Path not found" (probed
		// live; the root's {} is also what keeps 1.2.x disjoint_tables from
		// crashing there).
		if len(rows) == 0 {
			if path == "" {
				return map[string]any{}, nil
			}
			return nil, fmt.Errorf("%w: %s", ErrPathNotFound, path)
		}
		docs := make([]map[string]any, len(rows))
		stamps := make([]time.Time, len(rows))
		for i := range rows {
			doc, err := unmarshalDoc(rows[i].Value)
			if err != nil {
				return nil, err
			}
			docs[i] = doc
			stamps[i] = rows[i].TS
		}
		if opts.AllowBigIntegers != nil || opts.AllowSafeBigIntegers != nil {
			colTypes, err := objectColumnTypes(r.iface.Definition)
			if err != nil {
				return nil, err
			}
			applyBigIntegerOpts(docs, colTypes, opts)
		}
		return renderObject(docs, stamps, opts.Format), nil
	}

	rows, err := s.st.Series(ctx, q)
	if err != nil {
		return nil, err
	}
	return renderIndividual(rows, path, opts.Format), nil
}

// renderIndividual shapes one individual series per the requested format.
// structured is the [{timestamp,value}] envelope Astrate has always emitted;
// table emits [ts,value] pairs (timestamp FIRST) with column metadata keyed on
// the path's last segment; disjoint_tables pivots to {"value": [[value,ts]]} —
// value BEFORE timestamp, upstream's pair order. longinteger leaves stay
// decimal strings in every format here: upstream always formats individual
// samples with no bigint options.
func renderIndividual(rows []store.IndividualRow, path, format string) any {
	switch format {
	case "table":
		name := lastSegment(path)
		pairs := make([][]any, len(rows))
		for i := range rows {
			pairs[i] = []any{rows[i].TS, individualValue(&rows[i])}
		}
		return &Tabular{
			Data: pairs,
			Metadata: map[string]any{
				"columns":      map[string]int{"timestamp": 0, name: 1},
				"table_header": []string{"timestamp", name},
			},
		}
	case "disjoint_tables":
		pairs := make([][]any, len(rows))
		for i := range rows {
			pairs[i] = []any{individualValue(&rows[i]), rows[i].TS}
		}
		return map[string]any{"value": pairs}
	default:
		out := make([]Sample, len(rows))
		for i := range rows {
			out[i] = Sample{Value: individualValue(&rows[i]), Timestamp: rows[i].TS}
		}
		return out
	}
}

// renderObject shapes object-aggregated documents per the requested format.
// docs carry json.Number leaves so unrendered values re-marshal exactly as
// stored. structured flattens each document with its timestamp merged in
// ([{<name>:v,...,timestamp:t},...]); table aligns rows to the sorted union of
// the FIRST document's keys and "timestamp" (missing keys render null);
// disjoint_tables pivots to one [value,timestamp] pair list per column of the
// first document.
func renderObject(docs []map[string]any, stamps []time.Time, format string) any {
	switch format {
	case "table":
		header := make([]string, 0, len(docs[0])+1)
		for k := range docs[0] {
			header = append(header, k)
		}
		if !slices.Contains(header, "timestamp") {
			header = append(header, "timestamp")
		}
		slices.Sort(header)
		columns := make(map[string]int, len(header))
		for i, name := range header {
			columns[name] = i
		}
		rows := make([][]any, len(docs))
		for i, doc := range docs {
			row := make([]any, len(header))
			for j, name := range header {
				if name == "timestamp" {
					row[j] = stamps[i]
					continue
				}
				row[j] = doc[name]
			}
			rows[i] = row
		}
		return &Tabular{
			Data:     rows,
			Metadata: map[string]any{"columns": columns, "table_header": header},
		}
	case "disjoint_tables":
		out := make(map[string]any, len(docs[0]))
		for name := range docs[0] {
			pairs := make([][]any, len(docs))
			for i, doc := range docs {
				pairs[i] = []any{doc[name], stamps[i]}
			}
			out[name] = pairs
		}
		return out
	default:
		out := make([]map[string]any, len(docs))
		for i, doc := range docs {
			m := make(map[string]any, len(doc)+1)
			for k, v := range doc {
				m[k] = v
			}
			m["timestamp"] = stamps[i]
			out[i] = m
		}
		return out
	}
}

// lastSegment returns the text after the final "/" of an endpoint or path.
func lastSegment(p string) string {
	if i := strings.LastIndexByte(p, '/'); i >= 0 {
		return p[i+1:]
	}
	return p
}

// unmarshalDoc decodes one stored object document, keeping numbers as
// json.Number so longinteger leaves survive re-marshalling byte-exactly.
func unmarshalDoc(raw []byte) (map[string]any, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var doc map[string]any
	if err := dec.Decode(&doc); err != nil {
		return nil, fmt.Errorf("appengine: decoding stored object document: %w", err)
	}
	return doc, nil
}

// objectColumnTypes maps each mapping's LAST segment (the object-document key)
// to its declared value type, re-parsed from the stored definition.
func objectColumnTypes(definition []byte) (map[string]interfaceschema.ValueType, error) {
	iface, err := interfaceschema.ParseInterface(definition)
	if err != nil {
		return nil, fmt.Errorf("appengine: re-parsing stored interface definition: %w", err)
	}
	out := make(map[string]interfaceschema.ValueType, len(iface.Mappings))
	for i := range iface.Mappings {
		out[lastSegment(iface.Mappings[i].Endpoint)] = iface.Mappings[i].Type
	}
	return out, nil
}

// applyBigIntegerOpts rewrites longinteger leaves across the documents per the
// request's bigint flags (upstream fetch_biginteger_opts_or_default +
// AstarteValue.to_json_friendly): allow_bigintegers=false forces decimal
// strings; allow_safe_bigintegers=true keeps only |v| ≤ maxSafeJSONInteger as
// numbers; neither flag set leaves the stored representation untouched.
// Object results only — individual samples never consult these options.
func applyBigIntegerOpts(docs []map[string]any, colTypes map[string]interfaceschema.ValueType, opts QueryOpts) {
	for _, doc := range docs {
		for name, v := range doc {
			if colTypes[name] != interfaceschema.LongInteger {
				continue
			}
			n, ok := v.(json.Number)
			if !ok {
				continue
			}
			i, err := strconv.ParseInt(n.String(), 10, 64)
			if err != nil {
				// Beyond int64 there is no numeric JSON rendering at all.
				doc[name] = n.String()
				continue
			}
			switch {
			case opts.AllowBigIntegers != nil && !*opts.AllowBigIntegers:
				doc[name] = n.String()
			case opts.AllowSafeBigIntegers != nil && *opts.AllowSafeBigIntegers:
				if i > maxSafeJSONInteger || i < -maxSafeJSONInteger {
					doc[name] = n.String()
				}
			}
		}
	}
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
