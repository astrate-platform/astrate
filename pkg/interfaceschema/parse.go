package interfaceschema

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Structural limits, matching upstream Astarte (astarte_core).
const (
	// MaxNameLength is the maximum interface name length.
	MaxNameLength = 128
	// MaxMappings is the maximum number of mappings per interface.
	MaxMappings = 1024
	// MaxEndpointDepth is the maximum number of endpoint levels.
	MaxEndpointDepth = 64
	// MaxDescriptionLength is the maximum description length at interface and
	// mapping level (probe-verified against upstream 1.2.0, issue #61).
	MaxDescriptionLength = 1000
	// MaxDocLength is the maximum doc length at interface and mapping level
	// (probe-verified against upstream 1.2.0, issue #61).
	MaxDocLength = 100000
	// MinDatabaseRetentionTTL is the smallest legal database_retention_ttl,
	// in seconds (upstream rejects anything below; issue #61).
	MinDatabaseRetentionTTL = 60
	// MaxDatabaseRetentionTTL is the exclusive upper bound for
	// database_retention_ttl, in seconds: twenty years (issue #61). A TTL of
	// exactly this value is rejected with "must be less than".
	MaxDatabaseRetentionTTL = 630720000
)

// ErrInvalid is wrapped by every ParseInterface failure, so callers can
// classify rejection with errors.Is regardless of the specific rule violated.
var ErrInvalid = errors.New("invalid interface")

// interfaceNameRe is upstream Astarte's interface-name pattern: a reverse
// domain name whose first and last labels start with a letter and contain
// only alphanumerics, with optional hyphenated middle labels.
var interfaceNameRe = regexp.MustCompile(
	`^([a-zA-Z][a-zA-Z0-9]*\.([a-zA-Z0-9][a-zA-Z0-9-]*\.)*)?[a-zA-Z][a-zA-Z0-9]*$`)

// Interface is a parsed and validated Astarte interface definition.
// Instances produced by ParseInterface are structurally valid; hand-built
// instances should be validated by round-tripping through ParseInterface.
type Interface struct {
	// Name is the reverse-domain interface name (≤ 128 characters).
	Name string
	// Major is the major version; majors coexist as distinct interfaces.
	Major int
	// Minor is the minor version; bumps must be additive (CheckMinorUpgrade).
	Minor int
	// Type discriminates datastream from properties.
	Type InterfaceType
	// Ownership states which side publishes on this interface.
	Ownership Ownership
	// Aggregation is individual or object (defaults to individual).
	Aggregation Aggregation
	// Description is the optional human-readable summary.
	Description string
	// Doc is the optional long-form documentation.
	Doc string
	// Mappings are the declared endpoints (1 to 1024 entries).
	Mappings []Mapping
}

// Mapping is one parsed endpoint declaration of an Interface.
type Mapping struct {
	// Endpoint is the declared path pattern, e.g. "/%{sensor_id}/value".
	Endpoint string
	// Type is the value type of data published on this endpoint.
	Type ValueType
	// Reliability is the delivery guarantee (datastream only).
	Reliability Reliability
	// Retention states what happens to undeliverable values (datastream only).
	Retention Retention
	// Expiry is the retention expiry in seconds; 0 means never (datastream only).
	Expiry int64
	// DatabaseRetentionPolicy states whether stored values get a TTL
	// (datastream only).
	DatabaseRetentionPolicy DatabaseRetentionPolicy
	// DatabaseRetentionTTL is the database TTL in seconds; set if and only
	// if DatabaseRetentionPolicy is UseTTL (datastream only).
	DatabaseRetentionTTL int64
	// AllowUnset permits unsetting the property (properties only).
	AllowUnset bool
	// ExplicitTimestamp states whether publishes carry their own timestamp
	// (datastream only).
	ExplicitTimestamp bool
	// Required states whether object documents must carry this key;
	// enforced at runtime only for object aggregation (upstream 1.4
	// experimental, issue #67).
	Required bool
	// Encrypted records the declared encryption flag; parsed and stored,
	// no runtime effect yet (upstream 1.4 experimental, issue #67).
	Encrypted bool
	// Description is the optional human-readable summary.
	Description string
	// Doc is the optional long-form documentation.
	Doc string
}

// interfaceJSON is the strict wire decoding target; pointers detect absence.
// The omitempty tags only matter on re-encode: after a legacy-alias
// normalisation the same struct is marshalled back as the canonical document
// (ParseInterfaceCanonical), where absent optionals must stay absent.
type interfaceJSON struct {
	InterfaceName *string        `json:"interface_name"`
	VersionMajor  *int           `json:"version_major"`
	VersionMinor  *int           `json:"version_minor"`
	Type          *InterfaceType `json:"type"`
	Ownership     *Ownership     `json:"ownership,omitempty"`
	// Quality is the upstream legacy alias of ownership (issue #61): accepted
	// alone, a violation next to it, never stored — the canonical value lands
	// in Ownership during normalisation.
	Quality *Ownership `json:"quality,omitempty"`
	// Aggregate is the upstream legacy alias of aggregation (issue #61).
	Aggregate   *bool         `json:"aggregate,omitempty"`
	Aggregation *Aggregation  `json:"aggregation,omitempty"`
	Description string        `json:"description,omitempty"`
	Doc         string        `json:"doc,omitempty"`
	Mappings    []mappingJSON `json:"mappings"`

	// aliases records whether any legacy alias field was consumed, so
	// canonical re-encoding is produced only when the wire form differs from
	// the canonical one. Not part of the wire format.
	aliases bool `json:"-"`
}

// mappingJSON is the strict wire decoding target for one mapping.
type mappingJSON struct {
	Endpoint *string `json:"endpoint,omitempty"`
	// Path is the upstream legacy alias of endpoint (issue #61): accepted in
	// place of it, silently ignored next to it (endpoint wins), never stored.
	Path                    *string                  `json:"path,omitempty"`
	Type                    *ValueType               `json:"type"`
	Reliability             *Reliability             `json:"reliability,omitempty"`
	Retention               *Retention               `json:"retention,omitempty"`
	Expiry                  *int64                   `json:"expiry,omitempty"`
	DatabaseRetentionPolicy *DatabaseRetentionPolicy `json:"database_retention_policy,omitempty"`
	DatabaseRetentionTTL    *int64                   `json:"database_retention_ttl,omitempty"`
	AllowUnset              *bool                    `json:"allow_unset,omitempty"`
	ExplicitTimestamp       *bool                    `json:"explicit_timestamp,omitempty"`
	Required                *bool                    `json:"required,omitempty"`
	Encrypted               *bool                    `json:"encrypted,omitempty"`
	Description             string                   `json:"description,omitempty"`
	Doc                     string                   `json:"doc,omitempty"`
}

// ParseInterface strictly decodes and validates an Astarte interface JSON
// document. Unknown fields, malformed values, and every structural rule
// violation (name syntax, versioning, endpoint syntax and uniqueness,
// aggregation and per-type field constraints) are rejected with an error
// wrapping ErrInvalid. Rejections for the rules whose upstream wire shape was
// probe-verified (issue #61) additionally carry a *ViolationsError in the
// error chain.
func ParseInterface(data []byte) (*Interface, error) {
	iface, _, err := ParseInterfaceCanonical(data)
	return iface, err
}

// ParseInterfaceCanonical parses like ParseInterface and additionally returns
// the canonical re-encoding of the document when any upstream legacy alias
// field (interface quality/aggregate, mapping path — issue #61) was consumed
// during normalisation: canon == nil means data already uses canonical field
// names only and can be stored as-is. Storing canon instead of data is what
// makes GET render the canonical form the way upstream does.
func ParseInterfaceCanonical(data []byte) (*Interface, []byte, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var raw interfaceJSON
	if err := dec.Decode(&raw); err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	if _, err := dec.Token(); err != io.EOF {
		return nil, nil, fmt.Errorf("%w: trailing data after interface document", ErrInvalid)
	}

	iface, err := buildInterface(&raw)
	if err != nil {
		var ve *ViolationsError
		if errors.As(err, &ve) {
			return nil, nil, err // the combined value already wraps ErrInvalid
		}
		return nil, nil, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	if !raw.aliases {
		return iface, nil, nil
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(&raw); err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrInvalid, err)
	}
	canon := buf.Bytes()
	return iface, canon[:len(canon)-1], nil // json.Encoder appends '\n'
}

// buildInterface applies defaults and every validation rule to the raw
// decoded document. Rule violations whose upstream wire envelope was
// probe-verified (issue #61) accumulate in a ViolationsError and abort the
// parse as a combined ErrInvalid + *ViolationsError value once every mapping
// has been examined, so one response can report all offending fields; every
// other rule keeps its plain immediate rejection.
func buildInterface(raw *interfaceJSON) (*Interface, error) {
	vc := &ViolationsError{}

	// Legacy aliases (issue #61): normalise before any rule runs. quality is
	// accepted only without ownership (both → violation on ownership);
	// aggregate only without aggregation (both → violation on aggregation).
	if raw.Quality != nil {
		raw.aliases = true
		if raw.Ownership != nil {
			vc.add("ownership", "ownership and quality are mutually exclusive")
		} else {
			raw.Ownership = raw.Quality
		}
		raw.Quality = nil
	}
	if raw.Aggregate != nil {
		raw.aliases = true
		if raw.Aggregation != nil {
			vc.add("aggregation", "aggregation and aggregate are mutually exclusive")
		} else {
			agg := AggregationIndividual
			if *raw.Aggregate {
				agg = AggregationObject
			}
			raw.Aggregation = &agg
		}
		raw.Aggregate = nil
	}

	switch {
	case raw.InterfaceName == nil:
		return nil, errors.New(`missing "interface_name"`)
	case raw.VersionMajor == nil:
		return nil, errors.New(`missing "version_major"`)
	case raw.VersionMinor == nil:
		return nil, errors.New(`missing "version_minor"`)
	case raw.Type == nil:
		return nil, errors.New(`missing "type"`)
	case raw.Ownership == nil:
		return nil, errors.New(`missing "ownership"`)
	}

	name := *raw.InterfaceName
	if len(name) > MaxNameLength {
		return nil, fmt.Errorf("interface name exceeds %d characters", MaxNameLength)
	}
	if !interfaceNameRe.MatchString(name) {
		return nil, fmt.Errorf("invalid interface name %q", name)
	}

	major, minor := *raw.VersionMajor, *raw.VersionMinor
	if major < 0 {
		return nil, errors.New("version_major must be >= 0")
	}
	if minor < 0 {
		return nil, errors.New("version_minor must be >= 0")
	}
	if major == 0 && minor == 0 {
		return nil, errors.New("version_major and version_minor cannot both be 0")
	}

	iface := &Interface{
		Name:        name,
		Major:       major,
		Minor:       minor,
		Type:        *raw.Type,
		Ownership:   *raw.Ownership,
		Aggregation: AggregationIndividual,
		Description: raw.Description,
		Doc:         raw.Doc,
	}
	if raw.Aggregation != nil {
		iface.Aggregation = *raw.Aggregation
	}
	if iface.Type == Properties && iface.Aggregation == AggregationObject {
		return nil, errors.New("properties interfaces cannot use object aggregation")
	}

	if len(raw.Mappings) == 0 {
		return nil, errors.New("interface must declare at least one mapping")
	}
	if len(raw.Mappings) > MaxMappings {
		return nil, fmt.Errorf("interface declares %d mappings, maximum is %d", len(raw.Mappings), MaxMappings)
	}

	// Interface-level length bounds (#61), after alias normalisation.
	vc.MappingCount = len(raw.Mappings)
	if len(raw.Description) > MaxDescriptionLength {
		vc.add("description", tooLongMessage(MaxDescriptionLength))
	}
	if len(raw.Doc) > MaxDocLength {
		vc.add("doc", tooLongMessage(MaxDocLength))
	}

	endpoints := make([][]endpointSegment, 0, len(raw.Mappings))
	for i := range raw.Mappings {
		mj := &raw.Mappings[i]
		if mj.Path != nil {
			// Mapping path alias (issue #61): endpoint wins when both are
			// present (measured upstream behaviour); otherwise path becomes
			// the endpoint.
			raw.aliases = true
			if mj.Endpoint == nil {
				mj.Endpoint = mj.Path
			}
			mj.Path = nil
		}
		m, segs, err := buildMapping(mj, iface.Type, i, vc)
		if err != nil {
			return nil, err
		}
		iface.Mappings = append(iface.Mappings, m)
		endpoints = append(endpoints, segs)
	}

	// Abort on any collected violation before the structural checks below:
	// their plain messages would shadow the field-shaped ones.
	if err := vc.errOrNil(); err != nil {
		return nil, err
	}

	if err := checkEndpointUniqueness(iface, endpoints); err != nil {
		return nil, err
	}
	if iface.Aggregation == AggregationObject {
		if err := checkObjectAggregation(iface, endpoints); err != nil {
			return nil, err
		}
	}
	return iface, nil
}

// tooLongMessage is upstream's exact length-violation message for a bound.
func tooLongMessage(bound int) string {
	return fmt.Sprintf("should be at most %d character(s)", bound)
}

// buildMapping validates one raw mapping against the per-type field rules
// and returns it with defaults applied, plus its parsed endpoint segments.
// Violations whose upstream wire envelope was probe-verified (issue #61) are
// appended to vc (with the mapping's declared index) instead of aborting, so
// sibling mappings still get checked; plain errors return immediately.
func buildMapping(raw *mappingJSON, ifaceType InterfaceType, index int, vc *ViolationsError) (Mapping, []endpointSegment, error) {
	var m Mapping
	if raw.Endpoint == nil {
		return m, nil, errors.New(`mapping is missing "endpoint"`)
	}
	if raw.Type == nil {
		return m, nil, fmt.Errorf(`mapping %q is missing "type"`, *raw.Endpoint)
	}
	segs, err := splitEndpoint(*raw.Endpoint)
	if err != nil {
		return m, nil, err
	}

	m = Mapping{
		Endpoint:    *raw.Endpoint,
		Type:        *raw.Type,
		Description: raw.Description,
		Doc:         raw.Doc,
	}

	// Mapping-level length bounds (#61), both interface types.
	if len(raw.Description) > MaxDescriptionLength {
		vc.addMapping(index, "description", tooLongMessage(MaxDescriptionLength))
	}
	if len(raw.Doc) > MaxDocLength {
		vc.addMapping(index, "doc", tooLongMessage(MaxDocLength))
	}

	if ifaceType == Properties {
		// Datastream-only attributes are rejected outright when present:
		// strict decode, no silent dropping.
		switch {
		case raw.Reliability != nil:
			return m, nil, fmt.Errorf(`mapping %q: "reliability" is not allowed on properties`, m.Endpoint)
		case raw.Retention != nil:
			return m, nil, fmt.Errorf(`mapping %q: "retention" is not allowed on properties`, m.Endpoint)
		case raw.Expiry != nil:
			return m, nil, fmt.Errorf(`mapping %q: "expiry" is not allowed on properties`, m.Endpoint)
		case raw.ExplicitTimestamp != nil:
			return m, nil, fmt.Errorf(`mapping %q: "explicit_timestamp" is not allowed on properties`, m.Endpoint)
		case raw.DatabaseRetentionPolicy != nil:
			return m, nil, fmt.Errorf(`mapping %q: "database_retention_policy" is not allowed on properties`, m.Endpoint)
		case raw.DatabaseRetentionTTL != nil:
			return m, nil, fmt.Errorf(`mapping %q: "database_retention_ttl" is not allowed on properties`, m.Endpoint)
		case raw.Required != nil:
			return m, nil, fmt.Errorf(`mapping %q: "required" is not allowed on properties`, m.Endpoint)
		case raw.Encrypted != nil:
			return m, nil, fmt.Errorf(`mapping %q: "encrypted" is not allowed on properties`, m.Endpoint)
		}
		if raw.AllowUnset != nil {
			m.AllowUnset = *raw.AllowUnset
		}
		return m, segs, nil
	}

	// Datastream.
	if raw.AllowUnset != nil {
		return m, nil, fmt.Errorf(`mapping %q: "allow_unset" is only allowed on properties`, m.Endpoint)
	}
	if raw.Required != nil {
		m.Required = *raw.Required
	}
	if raw.Encrypted != nil {
		m.Encrypted = *raw.Encrypted
	}
	if raw.Reliability != nil {
		m.Reliability = *raw.Reliability
	}
	if raw.Retention != nil {
		m.Retention = *raw.Retention
	}
	if raw.Expiry != nil {
		if *raw.Expiry < 0 {
			return m, nil, fmt.Errorf("mapping %q: expiry must be >= 0", m.Endpoint)
		}
		m.Expiry = *raw.Expiry
	}
	if raw.DatabaseRetentionPolicy != nil {
		m.DatabaseRetentionPolicy = *raw.DatabaseRetentionPolicy
	}
	switch m.DatabaseRetentionPolicy {
	case UseTTL:
		if raw.DatabaseRetentionTTL == nil {
			// Plain error, not a violation: this envelope was not probed.
			return m, nil, fmt.Errorf(
				`mapping %q: database_retention_policy "use_ttl" requires database_retention_ttl >= 1`, m.Endpoint)
		}
		// Probe-verified band (issue #61): [60, 630720000), replacing the old
		// ">= 1" rule.
		ttl := *raw.DatabaseRetentionTTL
		switch {
		case ttl < MinDatabaseRetentionTTL:
			vc.addMapping(index, "database_retention_ttl",
				fmt.Sprintf("must be greater than or equal to %d", MinDatabaseRetentionTTL))
		case ttl >= MaxDatabaseRetentionTTL:
			vc.addMapping(index, "database_retention_ttl",
				fmt.Sprintf("must be less than %d", MaxDatabaseRetentionTTL))
		default:
			m.DatabaseRetentionTTL = ttl
		}
	case NoTTL:
		if raw.DatabaseRetentionTTL != nil {
			return m, nil, fmt.Errorf(
				`mapping %q: database_retention_ttl requires database_retention_policy "use_ttl"`, m.Endpoint)
		}
	}
	if raw.ExplicitTimestamp != nil {
		m.ExplicitTimestamp = *raw.ExplicitTimestamp
	}
	return m, segs, nil
}

// endpointSegment is one parsed level of an endpoint pattern: either a
// literal name or a %{placeholder}.
type endpointSegment struct {
	literal string // segment text; placeholder name when parametric
	param   bool
}

// splitEndpoint parses and validates an endpoint pattern into its segments.
// Endpoints are '/'-rooted, 1–64 levels deep; each level is either a literal
// identifier ([A-Za-z_][A-Za-z0-9_]*) or a whole-segment %{placeholder}.
func splitEndpoint(endpoint string) ([]endpointSegment, error) {
	if endpoint == "" || endpoint[0] != '/' {
		return nil, fmt.Errorf("invalid endpoint %q: must start with '/'", endpoint)
	}
	var segs []endpointSegment
	rest := endpoint[1:]
	for {
		var raw string
		if i := strings.IndexByte(rest, '/'); i >= 0 {
			raw, rest = rest[:i], rest[i+1:]
		} else {
			raw, rest = rest, ""
		}
		seg, err := parseSegment(endpoint, raw)
		if err != nil {
			return nil, err
		}
		segs = append(segs, seg)
		if rest == "" && !endsWithSlash(endpoint) {
			break
		}
		if rest == "" {
			return nil, fmt.Errorf("invalid endpoint %q: empty segment", endpoint)
		}
	}
	if len(segs) > MaxEndpointDepth {
		return nil, fmt.Errorf("invalid endpoint %q: exceeds %d levels", endpoint, MaxEndpointDepth)
	}
	return segs, nil
}

// parseSegment validates one raw endpoint level.
func parseSegment(endpoint, raw string) (endpointSegment, error) {
	if raw == "" {
		return endpointSegment{}, fmt.Errorf("invalid endpoint %q: empty segment", endpoint)
	}
	if raw[0] == '%' {
		if len(raw) < 4 || raw[1] != '{' || raw[len(raw)-1] != '}' {
			return endpointSegment{}, fmt.Errorf("invalid endpoint %q: malformed placeholder %q", endpoint, raw)
		}
		name := raw[2 : len(raw)-1]
		if !validIdentifier(name) {
			return endpointSegment{}, fmt.Errorf("invalid endpoint %q: invalid placeholder name %q", endpoint, name)
		}
		return endpointSegment{literal: name, param: true}, nil
	}
	if !validIdentifier(raw) {
		if strings.ContainsAny(raw, "%{}") {
			return endpointSegment{}, fmt.Errorf(
				"invalid endpoint %q: placeholders must span a whole segment", endpoint)
		}
		return endpointSegment{}, fmt.Errorf("invalid endpoint %q: invalid segment %q", endpoint, raw)
	}
	return endpointSegment{literal: raw}, nil
}

// validIdentifier reports whether s matches [A-Za-z_][A-Za-z0-9_]*.
func validIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c == '_':
		case c >= '0' && c <= '9':
			if i == 0 {
				return false
			}
		default:
			return false
		}
	}
	return true
}

// endsWithSlash reports whether the endpoint has a trailing '/', which is an
// empty final segment.
func endsWithSlash(endpoint string) bool {
	return len(endpoint) > 0 && endpoint[len(endpoint)-1] == '/'
}

// checkEndpointUniqueness rejects duplicate and conflicting endpoint pairs.
// Two same-depth endpoints conflict when every level overlaps (equal
// literals, or at least one side parametric): such a pair would make some
// concrete inbound path ambiguous. A pair that overlaps with no parametric
// disagreement at any level is a duplicate.
func checkEndpointUniqueness(iface *Interface, endpoints [][]endpointSegment) error {
	for i := range endpoints {
		for j := i + 1; j < len(endpoints); j++ {
			a, b := endpoints[i], endpoints[j]
			if len(a) != len(b) {
				continue
			}
			overlap, identical := true, true
			for k := range a {
				switch {
				case a[k].param && b[k].param:
					// Placeholder names are not semantic; both parametric is
					// the same match set.
				case a[k].param != b[k].param:
					identical = false
				case a[k].literal != b[k].literal:
					overlap, identical = false, false
				}
				if !overlap {
					break
				}
			}
			if identical {
				return fmt.Errorf("duplicate endpoint %q (also declared as %q)",
					iface.Mappings[j].Endpoint, iface.Mappings[i].Endpoint)
			}
			if overlap {
				return fmt.Errorf("conflicting endpoints %q and %q match overlapping paths",
					iface.Mappings[i].Endpoint, iface.Mappings[j].Endpoint)
			}
		}
	}
	return nil
}

// checkObjectAggregation enforces the object-aggregation shape: every
// mapping shares the same depth, the same prefix (all levels but the last),
// uniform datastream attributes, and a distinct literal last level (the
// object key).
func checkObjectAggregation(iface *Interface, endpoints [][]endpointSegment) error {
	first := endpoints[0]
	ref := iface.Mappings[0]
	for i, segs := range endpoints {
		m := iface.Mappings[i]
		if len(segs) != len(first) {
			return fmt.Errorf("object-aggregated mappings must all have the same depth: %q has %d levels, %q has %d",
				m.Endpoint, len(segs), ref.Endpoint, len(first))
		}
		for k := 0; k < len(segs)-1; k++ {
			if segs[k] != first[k] {
				return fmt.Errorf("object-aggregated mappings must share the same prefix: %q diverges from %q",
					m.Endpoint, ref.Endpoint)
			}
		}
		if segs[len(segs)-1].param {
			return fmt.Errorf("object-aggregated mapping %q must end in a literal segment", m.Endpoint)
		}
		if err := sameObjectAttributes(ref, m); err != nil {
			return err
		}
	}
	return nil
}

// sameObjectAttributes verifies attribute uniformity between two mappings of
// an object-aggregated interface (description/doc may differ).
func sameObjectAttributes(ref, m Mapping) error {
	attr := ""
	switch {
	case m.Reliability != ref.Reliability:
		attr = "reliability"
	case m.Retention != ref.Retention:
		attr = "retention"
	case m.Expiry != ref.Expiry:
		attr = "expiry"
	case m.ExplicitTimestamp != ref.ExplicitTimestamp:
		attr = "explicit_timestamp"
	case m.DatabaseRetentionPolicy != ref.DatabaseRetentionPolicy:
		attr = "database_retention_policy"
	case m.DatabaseRetentionTTL != ref.DatabaseRetentionTTL:
		attr = "database_retention_ttl"
	}
	if attr != "" {
		return fmt.Errorf("object-aggregated mappings must all share the same %s: %q differs from %q",
			attr, m.Endpoint, ref.Endpoint)
	}
	return nil
}

// normalizeEndpoint returns the endpoint with placeholder names erased
// (each parametric level becomes "%{}"), the canonical key for comparing
// endpoints across interface versions where placeholder renames are not
// semantic.
func normalizeEndpoint(segs []endpointSegment) string {
	var buf bytes.Buffer
	for _, s := range segs {
		buf.WriteByte('/')
		if s.param {
			buf.WriteString("%{}")
		} else {
			buf.WriteString(s.literal)
		}
	}
	return buf.String()
}
