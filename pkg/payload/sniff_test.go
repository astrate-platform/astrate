package payload

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// TestDetectFormat covers each branch of the §3.5.2 algorithm.
func TestDetectFormat(t *testing.T) {
	emptyBSONDoc := []byte{0x05, 0x00, 0x00, 0x00, 0x00}
	cases := []struct {
		name string
		in   []byte
		want Format
	}{
		{"empty", nil, FormatEmpty},
		{"empty slice", []byte{}, FormatEmpty},
		{"minimal bson", emptyBSONDoc, FormatBSON},
		{"json envelope", []byte(`{"v":22.5}`), FormatJSON},
		{"json leading space", []byte(`   {"v":1}`), FormatJSON},
		{"json leading tab newline", []byte("\t\r\n{\"v\":true}"), FormatJSON},
		{"json garbage body still sniffs json", []byte(`{not json`), FormatJSON},
		{"bare number shorthand", []byte(`22.5`), FormatInvalid},
		{"bare string", []byte(`"hi"`), FormatInvalid},
		{"bare array", []byte(`[1,2]`), FormatInvalid},
		{"whitespace only", []byte("   \n"), FormatInvalid},
		{"bson prefix wrong length", []byte{0x06, 0x00, 0x00, 0x00, 0x00}, FormatInvalid},
		{"bson missing terminator", []byte{0x05, 0x00, 0x00, 0x00, 0x01}, FormatInvalid},
		{"four bytes never bson", []byte{0x04, 0x00, 0x00, 0x00}, FormatInvalid},
		{"random binary", []byte{0xde, 0xad, 0xbe, 0xef, 0x99}, FormatInvalid},
	}
	for _, tc := range cases {
		if got := DetectFormat(tc.in); got != tc.want {
			t.Errorf("%s: DetectFormat = %v; want %v", tc.name, got, tc.want)
		}
	}
}

// TestSniffNoCollision generates corpora of structurally valid documents of
// both kinds and asserts the sniff never crosses over (ROADMAP §2.4 gate):
// valid JSON never classifies as BSON, valid BSON never classifies as JSON.
func TestSniffNoCollision(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 7))

	// JSON corpus: envelopes of every shape, padded to every length from 5
	// to 600 bytes so the int32LE(p[0:4]) == len(p) condition is exercised
	// across the whole small-length range (including 123 = '{').
	jsonDocs := [][]byte{
		[]byte(`{"v":22.5}`),
		[]byte(`{"v":22.5,"t":"2026-06-10T12:34:56.789Z"}`),
		[]byte(`{"v":{"lat":45.0,"lon":11.3},"t":1718000000000}`),
		[]byte(`{"v":[1,2,3]}`),
		[]byte(`{"v":"` + strings.Repeat("a", 500) + `"}`),
	}
	for target := 5; target <= 600; target++ {
		pad := target - len(`{"v":""}`)
		if pad < 0 {
			continue
		}
		doc := []byte(`{"v":"` + strings.Repeat("x", pad) + `"}`)
		jsonDocs = append(jsonDocs, doc)
		// Same length reached via leading whitespace instead.
		ws := append([]byte(strings.Repeat(" ", 4)), []byte(`{"v":1}`)...)
		jsonDocs = append(jsonDocs, ws)
	}
	for _, doc := range jsonDocs {
		if !json.Valid(doc) {
			t.Fatalf("corpus bug: %q is not valid JSON", doc)
		}
		if got := DetectFormat(doc); got != FormatJSON {
			t.Errorf("valid JSON of %d bytes classified as %v: %q", len(doc), got, doc)
		}
	}

	// BSON corpus: valid documents of many sizes, built with the package
	// encoder (validated below), including the 123-byte document whose
	// first byte is '{'.
	for n := 0; n <= 256; n++ {
		v := strings.Repeat("b", n)
		doc, err := encodeBSON(v, nil)
		if err != nil {
			t.Fatalf("encodeBSON: %v", err)
		}
		if err := bson.Raw(doc).Validate(); err != nil {
			t.Fatalf("corpus bug: invalid BSON: %v", err)
		}
		if got := DetectFormat(doc); got != FormatBSON {
			t.Errorf("valid BSON of %d bytes classified as %v", len(doc), got)
		}
		if len(doc) == 123 && doc[0] != '{' {
			t.Error("expected the 123-byte document to start with '{'")
		}
	}

	// Random envelopes with random timestamps, both formats.
	for range 500 {
		ts := time.UnixMilli(rng.Int64N(maxDateTimeMillis)).UTC()
		val := rng.Float64() * 1e6
		doc, err := encodeBSON(val, &ts)
		if err != nil {
			t.Fatalf("encodeBSON: %v", err)
		}
		if got := DetectFormat(doc); got != FormatBSON {
			t.Errorf("random BSON classified as %v", got)
		}
		jdoc := fmt.Appendf(nil, `{"v":%g,"t":%d}`, val, ts.UnixMilli())
		if !json.Valid(jdoc) {
			t.Fatalf("corpus bug: %q", jdoc)
		}
		if got := DetectFormat(jdoc); got != FormatJSON {
			t.Errorf("random JSON classified as %v: %q", got, jdoc)
		}
	}
}
