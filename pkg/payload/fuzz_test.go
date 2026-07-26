package payload

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// fuzzMappings is one mapping per value type (plus an allow_unset property
// shape), so every fuzz input is driven through every decode path.
var fuzzMappings = func() []*interfaceschema.CompiledMapping {
	out := make([]*interfaceschema.CompiledMapping, 0, 16)
	for vt := interfaceschema.Double; vt <= interfaceschema.DateTimeArray; vt++ {
		out = append(out, &interfaceschema.CompiledMapping{ValueType: vt})
	}
	out = append(out,
		&interfaceschema.CompiledMapping{ValueType: interfaceschema.Double, AllowUnset: true},
		&interfaceschema.CompiledMapping{ValueType: interfaceschema.DateTime, ExplicitTimestamp: true},
	)
	return out
}()

// FuzzDecode drives the full facade with arbitrary input across every value
// type and the object path (ROADMAP §2.4 gate): it must never panic, every
// failure must carry a typed RejectReason, and valid inputs of either format
// must never be misclassified by the sniffer.
func FuzzDecode(f *testing.F) {
	// Seed with every golden SDK vector...
	entries, err := os.ReadDir(filepath.Join("testdata", "bson"))
	if err != nil {
		f.Fatalf("reading seed vectors: %v", err)
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".hex") {
			continue
		}
		b, err := os.ReadFile(filepath.Join("testdata", "bson", e.Name()))
		if err != nil {
			f.Fatalf("reading %s: %v", e.Name(), err)
		}
		raw, err := hex.DecodeString(strings.TrimSpace(string(b)))
		if err != nil {
			f.Fatalf("decoding %s: %v", e.Name(), err)
		}
		f.Add(raw)
	}
	// ...JSON-profile documents...
	for _, s := range []string{
		`{"v":22.5}`,
		`{"v":22.5,"t":"2026-06-10T12:34:56.789Z"}`,
		`{"v":"9007199254740993","t":1781094896789}`,
		`{"v":{"lat":45.0,"lon":9.0,"samples":3,"ok":true}}`,
		`{"v":[1,2,3]}`,
		`{"v":"3q2+7w=="}`,
		`{"v":null}`,
		`{}`,
	} {
		f.Add([]byte(s))
	}
	// ...and structural edge cases.
	f.Add([]byte{})
	f.Add([]byte{0x05, 0x00, 0x00, 0x00, 0x00})
	f.Add([]byte{0xDE, 0xAD, 0xBE, 0xEF, 0x99})

	leaves := map[string]*interfaceschema.CompiledMapping{
		"lat":     {ValueType: interfaceschema.Double},
		"lon":     {ValueType: interfaceschema.Double},
		"samples": {ValueType: interfaceschema.Integer},
		"ok":      {ValueType: interfaceschema.Boolean},
	}

	f.Fuzz(func(t *testing.T, p []byte) {
		for _, m := range fuzzMappings {
			dec, err := Decode(p, m)
			if err != nil {
				if ReasonOf(err) == ReasonNone {
					t.Fatalf("Decode returned an untyped error: %v", err)
				}
			} else if dec.Format == FormatInvalid {
				t.Fatalf("Decode succeeded with FormatInvalid for %q", p)
			}
		}
		if dec, err := DecodeObject(p, leaves); err != nil {
			if ReasonOf(err) == ReasonNone {
				t.Fatalf("DecodeObject returned an untyped error: %v", err)
			}
		} else if dec.Format == FormatInvalid {
			t.Fatalf("DecodeObject succeeded with FormatInvalid for %q", p)
		}

		// Misclassification guards (§3.5.2): a structurally valid document
		// of either format must sniff as that format.
		if json.Valid(p) && firstNonWS(p) == '{' && DetectFormat(p) != FormatJSON {
			t.Fatalf("valid JSON document classified as %v: %q", DetectFormat(p), p)
		}
		if len(p) >= 5 && int64(binary.LittleEndian.Uint32(p[0:4])) == int64(len(p)) &&
			bson.Raw(p).Validate() == nil && DetectFormat(p) != FormatBSON {
			t.Fatalf("valid BSON document classified as %v: %x", DetectFormat(p), p)
		}
	})
}

// firstNonWS returns the first non-whitespace byte of p (JSON whitespace
// set), or 0 if none.
func firstNonWS(p []byte) byte {
	for _, c := range p {
		switch c {
		case ' ', '\t', '\r', '\n':
			continue
		}
		return c
	}
	return 0
}
