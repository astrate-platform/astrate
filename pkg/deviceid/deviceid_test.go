package deviceid

import (
	"errors"
	"strings"
	"testing"
)

// Known (UUID, base64url) vectors cross-checked against the official
// astartectl binary, v1.0.1 (go install
// github.com/astarte-platform/astartectl@v1.0.1, run 2026-06-11):
//
//	$ astartectl utils device-id from-uuid 753ea14b-65bd-4d3e-8b12-73f6e64b3f96
//	dT6hS2W9TT6LEnP25ks_lg
//	$ astartectl utils device-id from-uuid 00000000-0000-0000-0000-000000000000
//	AAAAAAAAAAAAAAAAAAAAAA
//	$ astartectl utils device-id from-uuid f1e9a3b2-0c4d-4e5f-8a7b-6c5d4e3f2a1b
//	8emjsgxNTl-Ke2xdTj8qGw
//	$ astartectl utils device-id from-uuid 84b277c4-fdc6-4ab1-9d1b-9e5a47a09b6f
//	hLJ3xP3GSrGdG55aR6Cbbw
//	$ astartectl utils device-id to-uuid AJInS0w3VpWpuOqkXhgZdA
//	0092274b-4c37-5695-a9b8-eaa45e181974
var knownVectors = []struct {
	uuid string
	b64  string
}{
	{"753ea14b-65bd-4d3e-8b12-73f6e64b3f96", "dT6hS2W9TT6LEnP25ks_lg"},
	{"00000000-0000-0000-0000-000000000000", "AAAAAAAAAAAAAAAAAAAAAA"},
	{"f1e9a3b2-0c4d-4e5f-8a7b-6c5d4e3f2a1b", "8emjsgxNTl-Ke2xdTj8qGw"},
	{"84b277c4-fdc6-4ab1-9d1b-9e5a47a09b6f", "hLJ3xP3GSrGdG55aR6Cbbw"},
}

func TestKnownVectors(t *testing.T) {
	for _, v := range knownVectors {
		id, err := FromUUID(v.uuid)
		if err != nil {
			t.Fatalf("FromUUID(%q): %v", v.uuid, err)
		}
		if got := id.String(); got != v.b64 {
			t.Errorf("FromUUID(%q).String() = %q, want %q", v.uuid, got, v.b64)
		}
		parsed, err := Parse(v.b64)
		if err != nil {
			t.Fatalf("Parse(%q): %v", v.b64, err)
		}
		if parsed != id {
			t.Errorf("Parse(%q) = %v, want %v", v.b64, parsed, id)
		}
		if got := parsed.UUID(); got != v.uuid {
			t.Errorf("Parse(%q).UUID() = %q, want %q", v.b64, got, v.uuid)
		}
	}
}

func TestFromUUIDUppercase(t *testing.T) {
	lower, err := FromUUID("753ea14b-65bd-4d3e-8b12-73f6e64b3f96")
	if err != nil {
		t.Fatal(err)
	}
	upper, err := FromUUID("753EA14B-65BD-4D3E-8B12-73F6E64B3F96")
	if err != nil {
		t.Fatalf("FromUUID(uppercase): %v", err)
	}
	if lower != upper {
		t.Errorf("uppercase UUID parsed to %v, want %v", upper, lower)
	}
	// Output is always lowercase canonical.
	if got := upper.UUID(); got != "753ea14b-65bd-4d3e-8b12-73f6e64b3f96" {
		t.Errorf("UUID() = %q, want lowercase canonical", got)
	}
}

// FromNamespace vectors generated with astartectl v1.0.1 and independently
// cross-checked with CPython's uuid.uuid5 (both run 2026-06-11):
//
//	$ astartectl utils device-id compute-from-string \
//	      f79ad91f-c638-4889-ae74-9d001a3b4cf8 myidentifierdata
//	AJInS0w3VpWpuOqkXhgZdA
//	$ astartectl utils device-id compute-from-string \
//	      f79ad91f-c638-4889-ae74-9d001a3b4cf8 device_1
//	SEZYqHQXVfmUTY6AnlDXcA
//	$ astartectl utils device-id compute-from-string \
//	      f79ad91f-c638-4889-ae74-9d001a3b4cf8 ""
//	cUif8SUBU9a4x8oWAkZL4w
//
//	>>> uuid.uuid5(uuid.UUID('f79ad91f-c638-4889-ae74-9d001a3b4cf8'), 'myidentifierdata')
//	UUID('0092274b-4c37-5695-a9b8-eaa45e181974')
func TestFromNamespace(t *testing.T) {
	ns, err := FromUUID("f79ad91f-c638-4889-ae74-9d001a3b4cf8")
	if err != nil {
		t.Fatal(err)
	}
	vectors := []struct {
		payload string
		b64     string
		uuid    string
	}{
		{"myidentifierdata", "AJInS0w3VpWpuOqkXhgZdA", "0092274b-4c37-5695-a9b8-eaa45e181974"},
		{"device_1", "SEZYqHQXVfmUTY6AnlDXcA", "484658a8-7417-55f9-944d-8e809e50d770"},
		{"", "cUif8SUBU9a4x8oWAkZL4w", "71489ff1-2501-53d6-b8c7-ca1602464be3"},
	}
	for _, v := range vectors {
		id := FromNamespace(ns, v.payload)
		if got := id.String(); got != v.b64 {
			t.Errorf("FromNamespace(ns, %q) = %q, want %q", v.payload, got, v.b64)
		}
		if got := id.UUID(); got != v.uuid {
			t.Errorf("FromNamespace(ns, %q).UUID() = %q, want %q", v.payload, got, v.uuid)
		}
	}
}

func TestRoundTripRandom(t *testing.T) {
	for range 1000 {
		id, err := Random()
		if err != nil {
			t.Fatalf("Random(): %v", err)
		}
		s := id.String()
		if len(s) != EncodedLen {
			t.Fatalf("String() length = %d, want %d", len(s), EncodedLen)
		}
		back, err := Parse(s)
		if err != nil {
			t.Fatalf("Parse(%q): %v", s, err)
		}
		if back != id {
			t.Fatalf("round-trip mismatch: %v -> %q -> %v", id, s, back)
		}
		// UUID round-trip too.
		fromUUID, err := FromUUID(id.UUID())
		if err != nil {
			t.Fatalf("FromUUID(%q): %v", id.UUID(), err)
		}
		if fromUUID != id {
			t.Fatalf("UUID round-trip mismatch: %v -> %q -> %v", id, id.UUID(), fromUUID)
		}
	}
}

func TestRandomIsV4Shaped(t *testing.T) {
	seen := make(map[ID]bool)
	for range 100 {
		id, err := Random()
		if err != nil {
			t.Fatal(err)
		}
		if version := id[6] >> 4; version != 4 {
			t.Fatalf("version nibble = %d, want 4 (id %v)", version, id)
		}
		if variant := id[8] >> 6; variant != 0b10 {
			t.Fatalf("variant bits = %b, want 10 (id %v)", variant, id)
		}
		if seen[id] {
			t.Fatalf("duplicate random ID %v", id)
		}
		seen[id] = true
	}
}

func TestParseRejections(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"empty", ""},
		{"21 chars", "dT6hS2W9TT6LEnP25ks_l"},
		{"23 chars", "dT6hS2W9TT6LEnP25ks_lgA"},
		{"padded 24 chars", "dT6hS2W9TT6LEnP25ks_lg=="},
		{"padding char within 22", "dT6hS2W9TT6LEnP25ks_g="},
		{"standard alphabet plus", "dT6hS2W9TT6LEnP25ks+lg"},
		{"standard alphabet slash", "dT6hS2W9TT6LEnP25ks/lg"},
		{"whitespace", "dT6hS2W9TT6LEnP25ks lg"},
		{"non-ascii 22 bytes", "dT6hS2W9TT6LEnP25ks_é"}, // é is 2 bytes: total 22
		{"non-ascii 22 runes", "éééééééééééééééééééééé"},
		{"non-canonical trailing bits", "AAAAAAAAAAAAAAAAAAAAAB"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Parse(tc.in); err == nil {
				t.Errorf("Parse(%q) accepted, want rejection", tc.in)
			} else if !errors.Is(err, ErrInvalid) {
				t.Errorf("Parse(%q) error %v does not wrap ErrInvalid", tc.in, err)
			}
		})
	}
}

func TestFromUUIDRejections(t *testing.T) {
	cases := []string{
		"",
		"753ea14b-65bd-4d3e-8b12-73f6e64b3f9",    // 35 chars
		"753ea14b-65bd-4d3e-8b12-73f6e64b3f966",  // 37 chars
		"753ea14b65bd4d3e8b1273f6e64b3f96",       // no hyphens
		"{753ea14b-65bd-4d3e-8b12-73f6e64b3f96}", // braces
		"753ea14b-65bd-4d3e-8b12-73f6e64b3g96",   // non-hex
		"753ea14b+65bd-4d3e-8b12-73f6e64b3f96",   // wrong separator
		"urn:uuid:753ea14b-65bd-4d3e-8b12-73f6",  // URN prefix
	}
	for _, in := range cases {
		if _, err := FromUUID(in); err == nil {
			t.Errorf("FromUUID(%q) accepted, want rejection", in)
		} else if !errors.Is(err, ErrInvalid) {
			t.Errorf("FromUUID(%q) error %v does not wrap ErrInvalid", in, err)
		}
	}
}

func FuzzParse(f *testing.F) {
	f.Add("dT6hS2W9TT6LEnP25ks_lg")
	f.Add("AAAAAAAAAAAAAAAAAAAAAA")
	f.Add("")
	f.Add("dT6hS2W9TT6LEnP25ks+lg")
	f.Add(strings.Repeat("=", 22))
	f.Add("AAAAAAAAAAAAAAAAAAAAAB")
	f.Fuzz(func(t *testing.T, s string) {
		id, err := Parse(s)
		if err != nil {
			return
		}
		// Accepted inputs must round-trip byte-identically.
		if got := id.String(); got != s {
			t.Errorf("accepted %q but round-trips as %q", s, got)
		}
	})
}
