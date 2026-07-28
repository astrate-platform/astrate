// Package flow implements stream-based message routing through a block graph,
// mirroring Astarte Flow's core concepts: messages belong to streams (identified
// by key), streams are processed in order within a lane, and different streams
// may interleave across lanes.
package flow

// DataType enumerates the wire-level value types Astarte Flow supports.
type DataType uint8

const (
	TypeInteger DataType = iota
	TypeReal
	TypeBoolean
	TypeDatetime
	TypeBinary
	TypeString
)

// FlowMessage is one unit of data flowing through a block graph. Every message
// carries a Key that identifies its stream; messages with the same key are
// processed in submission order by the same lane (consistent hashing).
type FlowMessage struct {
	// Key identifies the stream this message belongs to. It must be non-empty.
	Key string
	// Metadata is an optional string→string map carried alongside the payload.
	Metadata map[string]string
	// Type is the base data type of the payload.
	Type DataType
	// Subtype is an optional MIME hint (meaningful when Type is TypeBinary).
	Subtype string
	// Timestamp is the event-time in microseconds since Unix epoch.
	Timestamp int64
	// Data is the payload; its concrete type must match Type (int64 for
	// TypeInteger, float64 for TypeReal, bool for TypeBoolean, time.Time for
	// TypeDatetime, []byte for TypeBinary, string for TypeString).
	Data any
}
