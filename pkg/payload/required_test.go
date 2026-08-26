package payload

import (
	"errors"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// TestObjectRequiredKeys pins the upstream 1.4 `required` enforcement
// (issue #67): an object-aggregation document omitting a key whose mapping
// declares required is rejected with ReasonMissingRequired in both wire
// formats, while complete documents still decode. The acceptance rows guard
// against a blanket refusal of every object document passing the rejection
// rows.
func TestObjectRequiredKeys(t *testing.T) {
	leaves := map[string]*interfaceschema.CompiledMapping{
		"lat": {ValueType: interfaceschema.Double, Required: true},
		"lon": {ValueType: interfaceschema.Double},
	}

	// JSON acceptance.
	if _, err := DecodeObject([]byte(`{"v":{"lat":45.0,"lon":9.0}}`), leaves); err != nil {
		t.Fatalf("JSON acceptance: %v", err)
	}

	// BSON acceptance.
	bsonOK, err := Encode(map[string]Value{"lat": 45.0, "lon": 9.0}, nil, FormatBSON)
	if err != nil {
		t.Fatalf("Encode(BSON): %v", err)
	}
	if _, err := DecodeObject(bsonOK, leaves); err != nil {
		t.Fatalf("BSON acceptance: %v", err)
	}

	// JSON rejection: the exact reason, not merely any rejection.
	if _, err := DecodeObject([]byte(`{"v":{"lon":9.0}}`), leaves); ReasonOf(err) != ReasonMissingRequired {
		t.Fatalf("JSON rejection reason = %v (err %v); want %v", ReasonOf(err), err, ReasonMissingRequired)
	}

	// BSON rejection: same exact reason.
	bsonBad, err := Encode(map[string]Value{"lon": 9.0}, nil, FormatBSON)
	if err != nil {
		t.Fatalf("Encode(BSON): %v", err)
	}
	if _, err := DecodeObject(bsonBad, leaves); ReasonOf(err) != ReasonMissingRequired {
		t.Fatalf("BSON rejection reason = %v (err %v); want %v", ReasonOf(err), err, ReasonMissingRequired)
	}

	// Determinism: with several required keys absent, the detail mentions
	// all of them regardless of map iteration order (sorted, then joined).
	multiLeaves := map[string]*interfaceschema.CompiledMapping{
		"lat": {ValueType: interfaceschema.Double, Required: true},
		"alt": {ValueType: interfaceschema.Double, Required: true},
		"lon": {ValueType: interfaceschema.Double},
	}
	if _, err := DecodeObject([]byte(`{"v":{"lon":1.0}}`), multiLeaves); err != nil {
		var re *RejectError
		if !errors.As(err, &re) {
			t.Fatalf("multi-missing case returned %T; want *RejectError", err)
		}
		for _, key := range []string{"alt", "lat"} {
			if !strings.Contains(re.Detail, key) {
				t.Errorf("detail %q does not mention required key %q", re.Detail, key)
			}
		}
	} else {
		t.Fatal("multi-missing case decoded without error")
	}
}
