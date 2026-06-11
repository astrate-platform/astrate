package payload

import (
	"encoding/binary"
	"math"
)

// DetectFormat classifies a raw inbound payload, implementing the exact
// docs/DESIGN.md §3.5.2 sniffing algorithm:
//
//	if len(p) == 0                      → FormatEmpty (property unset)
//	else if len(p) >= 5
//	     && int32LE(p[0:4]) == len(p)   → FormatBSON (self-describing prefix)
//	     && p[len(p)-1] == 0x00
//	else if first non-WS byte == '{'    → FormatJSON
//	else                                → FormatInvalid
//
// The two structural branches cannot collide: valid JSON text contains no
// NUL byte, so it can never satisfy the BSON terminator condition, while a
// BSON document whose first byte happens to be '{' (a 123-byte document) is
// claimed by the BSON branch first. DetectFormat never allocates.
func DetectFormat(p []byte) Format {
	if len(p) == 0 {
		return FormatEmpty
	}
	if len(p) >= 5 && len(p) <= math.MaxInt32 &&
		binary.LittleEndian.Uint32(p[0:4]) == uint32(len(p)) && //nolint:gosec // len bounded by MaxInt32 above
		p[len(p)-1] == 0x00 {
		return FormatBSON
	}
	for _, c := range p {
		switch c {
		case ' ', '\t', '\r', '\n':
			continue
		case '{':
			return FormatJSON
		}
		break
	}
	return FormatInvalid
}
