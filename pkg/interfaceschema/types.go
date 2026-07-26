// Package interfaceschema parses, validates, and compiles Astarte interface
// definitions — the dynamic typed schemas (datastreams and properties) that
// devices and applications exchange data through.
//
// The package mirrors upstream Astarte semantics (docs/DESIGN.md §2.6):
// ParseInterface performs a strict decode plus full structural validation of
// the Interface JSON, Compile turns a validated definition into the hot-path
// CompiledInterface (endpoint trie + object leaves), and CheckMinorUpgrade
// enforces the additive-only minor-version compatibility rule that Realm
// Management applies on interface updates.
package interfaceschema

import (
	"encoding/json"
	"fmt"
)

// InterfaceType discriminates datastream interfaces (time-ordered values)
// from properties interfaces (retained key/value state).
type InterfaceType uint8

// InterfaceType values (wire strings: "datastream", "properties").
const (
	// Datastream is a stream of timestamped values.
	Datastream InterfaceType = iota + 1
	// Properties is retained, settable/unsettable key/value state.
	Properties
)

// ParseInterfaceType parses the wire form of an interface type.
func ParseInterfaceType(s string) (InterfaceType, error) {
	switch s {
	case "datastream":
		return Datastream, nil
	case "properties":
		return Properties, nil
	default:
		return 0, fmt.Errorf("unknown interface type %q", s)
	}
}

// String returns the wire form.
func (t InterfaceType) String() string {
	switch t {
	case Datastream:
		return "datastream"
	case Properties:
		return "properties"
	default:
		return fmt.Sprintf("InterfaceType(%d)", uint8(t))
	}
}

// UnmarshalJSON decodes the wire form, rejecting unknown values.
func (t *InterfaceType) UnmarshalJSON(b []byte) error { return unmarshalEnum(b, t, ParseInterfaceType) }

// MarshalJSON encodes the wire form.
func (t InterfaceType) MarshalJSON() ([]byte, error) { return marshalEnum(t, Datastream, Properties) }

// Ownership states which side of the platform writes an interface: the
// device or the server (AppEngine and triggers).
type Ownership uint8

// Ownership values (wire strings: "device", "server").
const (
	// OwnershipDevice marks device-published interfaces.
	OwnershipDevice Ownership = iota + 1
	// OwnershipServer marks server-published interfaces.
	OwnershipServer
)

// ParseOwnership parses the wire form of an ownership.
func ParseOwnership(s string) (Ownership, error) {
	switch s {
	case "device":
		return OwnershipDevice, nil
	case "server":
		return OwnershipServer, nil
	default:
		return 0, fmt.Errorf("unknown ownership %q", s)
	}
}

// String returns the wire form.
func (o Ownership) String() string {
	switch o {
	case OwnershipDevice:
		return "device"
	case OwnershipServer:
		return "server"
	default:
		return fmt.Sprintf("Ownership(%d)", uint8(o))
	}
}

// UnmarshalJSON decodes the wire form, rejecting unknown values.
func (o *Ownership) UnmarshalJSON(b []byte) error { return unmarshalEnum(b, o, ParseOwnership) }

// MarshalJSON encodes the wire form.
func (o Ownership) MarshalJSON() ([]byte, error) {
	return marshalEnum(o, OwnershipDevice, OwnershipServer)
}

// Aggregation states whether each mapping is sent individually (one value
// per publish, full endpoint path) or as one object document of last-level
// keys published on the common path prefix.
type Aggregation uint8

// Aggregation values (wire strings: "individual", "object").
const (
	// AggregationIndividual sends one value per endpoint publish (default).
	AggregationIndividual Aggregation = iota + 1
	// AggregationObject sends all last-level values as a single document.
	AggregationObject
)

// ParseAggregation parses the wire form of an aggregation.
func ParseAggregation(s string) (Aggregation, error) {
	switch s {
	case "individual":
		return AggregationIndividual, nil
	case "object":
		return AggregationObject, nil
	default:
		return 0, fmt.Errorf("unknown aggregation %q", s)
	}
}

// String returns the wire form.
func (a Aggregation) String() string {
	switch a {
	case AggregationIndividual:
		return "individual"
	case AggregationObject:
		return "object"
	default:
		return fmt.Sprintf("Aggregation(%d)", uint8(a))
	}
}

// UnmarshalJSON decodes the wire form, rejecting unknown values.
func (a *Aggregation) UnmarshalJSON(b []byte) error { return unmarshalEnum(b, a, ParseAggregation) }

// MarshalJSON encodes the wire form.
func (a Aggregation) MarshalJSON() ([]byte, error) {
	return marshalEnum(a, AggregationIndividual, AggregationObject)
}

// Reliability is the delivery guarantee of a datastream mapping. Its numeric
// value equals the MQTT QoS level it maps to.
type Reliability uint8

// Reliability values (wire strings: "unreliable", "guaranteed", "unique").
const (
	// ReliabilityUnreliable is at-most-once delivery (QoS 0, default).
	ReliabilityUnreliable Reliability = iota
	// ReliabilityGuaranteed is at-least-once delivery (QoS 1).
	ReliabilityGuaranteed
	// ReliabilityUnique is exactly-once delivery (QoS 2).
	ReliabilityUnique
)

// ParseReliability parses the wire form of a reliability.
func ParseReliability(s string) (Reliability, error) {
	switch s {
	case "unreliable":
		return ReliabilityUnreliable, nil
	case "guaranteed":
		return ReliabilityGuaranteed, nil
	case "unique":
		return ReliabilityUnique, nil
	default:
		return 0, fmt.Errorf("unknown reliability %q", s)
	}
}

// String returns the wire form.
func (r Reliability) String() string {
	switch r {
	case ReliabilityUnreliable:
		return "unreliable"
	case ReliabilityGuaranteed:
		return "guaranteed"
	case ReliabilityUnique:
		return "unique"
	default:
		return fmt.Sprintf("Reliability(%d)", uint8(r))
	}
}

// QoS returns the MQTT QoS byte this reliability maps to
// (unreliable→0, guaranteed→1, unique→2).
func (r Reliability) QoS() byte { return byte(r) }

// UnmarshalJSON decodes the wire form, rejecting unknown values.
func (r *Reliability) UnmarshalJSON(b []byte) error { return unmarshalEnum(b, r, ParseReliability) }

// MarshalJSON encodes the wire form.
func (r Reliability) MarshalJSON() ([]byte, error) {
	return marshalEnum(r, ReliabilityUnreliable, ReliabilityUnique)
}

// Retention states what happens to datastream values that cannot be
// delivered immediately.
type Retention uint8

// Retention values (wire strings: "discard", "volatile", "stored").
const (
	// RetentionDiscard drops undeliverable values (default).
	RetentionDiscard Retention = iota
	// RetentionVolatile keeps undeliverable values in memory.
	RetentionVolatile
	// RetentionStored keeps undeliverable values on disk.
	RetentionStored
)

// ParseRetention parses the wire form of a retention.
func ParseRetention(s string) (Retention, error) {
	switch s {
	case "discard":
		return RetentionDiscard, nil
	case "volatile":
		return RetentionVolatile, nil
	case "stored":
		return RetentionStored, nil
	default:
		return 0, fmt.Errorf("unknown retention %q", s)
	}
}

// String returns the wire form.
func (r Retention) String() string {
	switch r {
	case RetentionDiscard:
		return "discard"
	case RetentionVolatile:
		return "volatile"
	case RetentionStored:
		return "stored"
	default:
		return fmt.Sprintf("Retention(%d)", uint8(r))
	}
}

// UnmarshalJSON decodes the wire form, rejecting unknown values.
func (r *Retention) UnmarshalJSON(b []byte) error { return unmarshalEnum(b, r, ParseRetention) }

// MarshalJSON encodes the wire form.
func (r Retention) MarshalJSON() ([]byte, error) {
	return marshalEnum(r, RetentionDiscard, RetentionStored)
}

// DatabaseRetentionPolicy states whether stored datastream values expire
// from the database.
type DatabaseRetentionPolicy uint8

// DatabaseRetentionPolicy values (wire strings: "no_ttl", "use_ttl").
const (
	// NoTTL keeps values forever (default).
	NoTTL DatabaseRetentionPolicy = iota
	// UseTTL expires values after database_retention_ttl seconds.
	UseTTL
)

// ParseDatabaseRetentionPolicy parses the wire form of a database retention
// policy.
func ParseDatabaseRetentionPolicy(s string) (DatabaseRetentionPolicy, error) {
	switch s {
	case "no_ttl":
		return NoTTL, nil
	case "use_ttl":
		return UseTTL, nil
	default:
		return 0, fmt.Errorf("unknown database_retention_policy %q", s)
	}
}

// String returns the wire form.
func (p DatabaseRetentionPolicy) String() string {
	switch p {
	case NoTTL:
		return "no_ttl"
	case UseTTL:
		return "use_ttl"
	default:
		return fmt.Sprintf("DatabaseRetentionPolicy(%d)", uint8(p))
	}
}

// UnmarshalJSON decodes the wire form, rejecting unknown values.
func (p *DatabaseRetentionPolicy) UnmarshalJSON(b []byte) error {
	return unmarshalEnum(b, p, ParseDatabaseRetentionPolicy)
}

// MarshalJSON encodes the wire form.
func (p DatabaseRetentionPolicy) MarshalJSON() ([]byte, error) { return marshalEnum(p, NoTTL, UseTTL) }

// ValueType is the Astarte mapping value type: seven scalars and their seven
// array counterparts. It drives both payload validation and BSON/JSON
// decoding (docs/DESIGN.md §2.6 step 5).
type ValueType uint8

// ValueType values. Wire strings are the lowercase names ("double",
// "doublearray", ...).
const (
	// Double is an IEEE 754 binary64 number.
	Double ValueType = iota + 1
	// Integer is a signed 32-bit integer.
	Integer
	// Boolean is true or false.
	Boolean
	// LongInteger is a signed 64-bit integer.
	LongInteger
	// String is a UTF-8 string (≤ 64 KiB).
	String
	// BinaryBlob is an arbitrary byte sequence (base64 in JSON payloads).
	BinaryBlob
	// DateTime is a UTC timestamp with millisecond precision.
	DateTime
	// DoubleArray is a homogeneous array of Double.
	DoubleArray
	// IntegerArray is a homogeneous array of Integer.
	IntegerArray
	// BooleanArray is a homogeneous array of Boolean.
	BooleanArray
	// LongIntegerArray is a homogeneous array of LongInteger.
	LongIntegerArray
	// StringArray is a homogeneous array of String.
	StringArray
	// BinaryBlobArray is a homogeneous array of BinaryBlob.
	BinaryBlobArray
	// DateTimeArray is a homogeneous array of DateTime.
	DateTimeArray
)

// valueTypeNames maps ValueType (offset by 1) to its wire string.
var valueTypeNames = [...]string{
	"double", "integer", "boolean", "longinteger", "string", "binaryblob", "datetime",
	"doublearray", "integerarray", "booleanarray", "longintegerarray", "stringarray",
	"binaryblobarray", "datetimearray",
}

// ParseValueType parses the wire form of a value type.
func ParseValueType(s string) (ValueType, error) {
	for i, name := range valueTypeNames {
		if s == name {
			return ValueType(i + 1), nil //nolint:gosec // i < 14, no overflow
		}
	}
	return 0, fmt.Errorf("unknown value type %q", s)
}

// String returns the wire form.
func (v ValueType) String() string {
	if v >= Double && v <= DateTimeArray {
		return valueTypeNames[v-1]
	}
	return fmt.Sprintf("ValueType(%d)", uint8(v))
}

// IsArray reports whether v is one of the seven array types.
func (v ValueType) IsArray() bool { return v >= DoubleArray && v <= DateTimeArray }

// Elem returns the scalar element type for array types and v itself for
// scalar types.
func (v ValueType) Elem() ValueType {
	if v.IsArray() {
		return v - DoubleArray + Double
	}
	return v
}

// UnmarshalJSON decodes the wire form, rejecting unknown values.
func (v *ValueType) UnmarshalJSON(b []byte) error { return unmarshalEnum(b, v, ParseValueType) }

// MarshalJSON encodes the wire form.
func (v ValueType) MarshalJSON() ([]byte, error) { return marshalEnum(v, Double, DateTimeArray) }

// enum constrains the generic JSON helpers to this package's uint8-backed
// enum types.
type enum interface {
	~uint8
	fmt.Stringer
}

// unmarshalEnum decodes a JSON string token and parses it with parse.
func unmarshalEnum[E enum](b []byte, dst *E, parse func(string) (E, error)) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("expected a string: %w", err)
	}
	v, err := parse(s)
	if err != nil {
		return err
	}
	*dst = v
	return nil
}

// marshalEnum encodes the wire string of v, rejecting out-of-range values.
func marshalEnum[E enum](v, lo, hi E) ([]byte, error) {
	if v < lo || v > hi {
		return nil, fmt.Errorf("cannot marshal invalid %s", v.String())
	}
	return json.Marshal(v.String())
}
