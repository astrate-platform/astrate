package interfaceschema_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// readFixtures returns name → bytes for every .json file in dir.
func readFixtures(t testing.TB, dir string) map[string][]byte {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		t.Fatalf("globbing %s: %v", dir, err)
	}
	if len(paths) == 0 {
		t.Fatalf("no fixtures in %s", dir)
	}
	out := make(map[string][]byte, len(paths))
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("reading %s: %v", p, err)
		}
		out[filepath.Base(p)] = b
	}
	return out
}

func TestValidFixturesParse(t *testing.T) {
	fixtures := readFixtures(t, filepath.Join("testdata", "valid"))
	if len(fixtures) < 12 {
		t.Fatalf("expected at least 12 valid fixtures, found %d", len(fixtures))
	}
	for name, data := range fixtures {
		t.Run(name, func(t *testing.T) {
			iface, err := interfaceschema.ParseInterface(data)
			if err != nil {
				t.Fatalf("ParseInterface: %v", err)
			}
			want := strings.TrimSuffix(name, ".json")
			if iface.Name != want {
				t.Errorf("Name = %q, want %q (fixture filename)", iface.Name, want)
			}
		})
	}
}

// TestValidFixturesCoverEveryValueType guards the fixture set itself: all 14
// value types must appear across the valid fixtures.
func TestValidFixturesCoverEveryValueType(t *testing.T) {
	seen := make(map[interfaceschema.ValueType]bool)
	for _, data := range readFixtures(t, filepath.Join("testdata", "valid")) {
		iface, err := interfaceschema.ParseInterface(data)
		if err != nil {
			continue // reported by TestValidFixturesParse
		}
		for _, m := range iface.Mappings {
			seen[m.Type] = true
		}
	}
	for vt := interfaceschema.Double; vt <= interfaceschema.DateTimeArray; vt++ {
		if !seen[vt] {
			t.Errorf("no valid fixture declares a %s mapping", vt)
		}
	}
}

func TestParseVendoredUpstreamFields(t *testing.T) {
	fixtures := readFixtures(t, filepath.Join("testdata", "valid"))

	t.Run("Geolocation object aggregation", func(t *testing.T) {
		iface, err := interfaceschema.ParseInterface(fixtures["org.astarte-platform.genericsensors.Geolocation.json"])
		if err != nil {
			t.Fatal(err)
		}
		if iface.Type != interfaceschema.Datastream {
			t.Errorf("Type = %v, want datastream", iface.Type)
		}
		if iface.Aggregation != interfaceschema.AggregationObject {
			t.Errorf("Aggregation = %v, want object", iface.Aggregation)
		}
		if iface.Major != 1 || iface.Minor != 0 {
			t.Errorf("version = %d.%d, want 1.0", iface.Major, iface.Minor)
		}
		if len(iface.Mappings) != 7 {
			t.Fatalf("len(Mappings) = %d, want 7", len(iface.Mappings))
		}
		for _, m := range iface.Mappings {
			if !m.ExplicitTimestamp {
				t.Errorf("mapping %s: ExplicitTimestamp = false, want true", m.Endpoint)
			}
			if m.Type != interfaceschema.Double {
				t.Errorf("mapping %s: Type = %v, want double", m.Endpoint, m.Type)
			}
		}
	})

	t.Run("ServerCommands retention", func(t *testing.T) {
		iface, err := interfaceschema.ParseInterface(fixtures["org.astarte-platform.genericcommands.ServerCommands.json"])
		if err != nil {
			t.Fatal(err)
		}
		if iface.Ownership != interfaceschema.OwnershipServer {
			t.Errorf("Ownership = %v, want server", iface.Ownership)
		}
		m := iface.Mappings[0]
		if m.DatabaseRetentionPolicy != interfaceschema.UseTTL {
			t.Errorf("DatabaseRetentionPolicy = %v, want use_ttl", m.DatabaseRetentionPolicy)
		}
		if m.DatabaseRetentionTTL != 86400 {
			t.Errorf("DatabaseRetentionTTL = %d, want 86400", m.DatabaseRetentionTTL)
		}
		if m.Type != interfaceschema.String {
			t.Errorf("Type = %v, want string", m.Type)
		}
	})

	t.Run("SamplingRate allow_unset", func(t *testing.T) {
		iface, err := interfaceschema.ParseInterface(fixtures["org.astarte-platform.genericsensors.SamplingRate.json"])
		if err != nil {
			t.Fatal(err)
		}
		if iface.Type != interfaceschema.Properties {
			t.Errorf("Type = %v, want properties", iface.Type)
		}
		for _, m := range iface.Mappings {
			if !m.AllowUnset {
				t.Errorf("mapping %s: AllowUnset = false, want true", m.Endpoint)
			}
		}
	})

	t.Run("Values parametric endpoint", func(t *testing.T) {
		iface, err := interfaceschema.ParseInterface(fixtures["org.astarte-platform.genericsensors.Values.json"])
		if err != nil {
			t.Fatal(err)
		}
		if got := iface.Mappings[0].Endpoint; got != "/%{sensor_id}/value" {
			t.Errorf("Endpoint = %q", got)
		}
		if got := iface.Mappings[0].Reliability; got != interfaceschema.ReliabilityUnreliable {
			t.Errorf("Reliability = %v, want default unreliable", got)
		}
	})
}

func TestParseDefaults(t *testing.T) {
	iface, err := interfaceschema.ParseInterface([]byte(`{
		"interface_name": "com.astrate.test.Defaults",
		"version_major": 0,
		"version_minor": 1,
		"type": "datastream",
		"ownership": "device",
		"mappings": [{"endpoint": "/v", "type": "integer"}]
	}`))
	if err != nil {
		t.Fatal(err)
	}
	if iface.Aggregation != interfaceschema.AggregationIndividual {
		t.Errorf("Aggregation = %v, want individual default", iface.Aggregation)
	}
	m := iface.Mappings[0]
	if m.Reliability != interfaceschema.ReliabilityUnreliable {
		t.Errorf("Reliability = %v, want unreliable default", m.Reliability)
	}
	if m.Reliability.QoS() != 0 {
		t.Errorf("QoS = %d, want 0", m.Reliability.QoS())
	}
	if m.Retention != interfaceschema.RetentionDiscard {
		t.Errorf("Retention = %v, want discard default", m.Retention)
	}
	if m.DatabaseRetentionPolicy != interfaceschema.NoTTL {
		t.Errorf("DatabaseRetentionPolicy = %v, want no_ttl default", m.DatabaseRetentionPolicy)
	}
	if m.Expiry != 0 || m.DatabaseRetentionTTL != 0 || m.ExplicitTimestamp || m.AllowUnset {
		t.Errorf("non-zero datastream defaults: %+v", m)
	}
}

func TestInvalidFixturesRejected(t *testing.T) {
	manifestBytes, err := os.ReadFile(filepath.Join("testdata", "invalid", "manifest.json"))
	if err != nil {
		t.Fatalf("reading manifest: %v", err)
	}
	var manifest map[string]string
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatalf("decoding manifest: %v", err)
	}
	if len(manifest) < 15 {
		t.Fatalf("expected at least 15 invalid fixtures in manifest, found %d", len(manifest))
	}

	fixtures := readFixtures(t, filepath.Join("testdata", "invalid"))
	delete(fixtures, "manifest.json")

	// The manifest and the fixture directory must stay in lockstep.
	for name := range manifest {
		if _, ok := fixtures[name]; !ok {
			t.Errorf("manifest entry %q has no fixture file", name)
		}
	}
	for name := range fixtures {
		if _, ok := manifest[name]; !ok {
			t.Errorf("fixture %q is missing from manifest.json", name)
		}
	}

	for name, wantSubstr := range manifest {
		data, ok := fixtures[name]
		if !ok {
			continue
		}
		t.Run(name, func(t *testing.T) {
			_, err := interfaceschema.ParseInterface(data)
			if err == nil {
				t.Fatalf("ParseInterface accepted %s, want rejection containing %q", name, wantSubstr)
			}
			if !errors.Is(err, interfaceschema.ErrInvalid) {
				t.Errorf("error %v does not wrap ErrInvalid", err)
			}
			if !strings.Contains(err.Error(), wantSubstr) {
				t.Errorf("error %q does not contain %q", err.Error(), wantSubstr)
			}
		})
	}
}

func TestParseRejectsNonObject(t *testing.T) {
	for _, in := range []string{"", "null", "[]", `"interface"`, "42", "{", `{"interface_name": 7}`} {
		if _, err := interfaceschema.ParseInterface([]byte(in)); err == nil {
			t.Errorf("ParseInterface(%q) succeeded, want error", in)
		} else if !errors.Is(err, interfaceschema.ErrInvalid) {
			t.Errorf("ParseInterface(%q) error %v does not wrap ErrInvalid", in, err)
		}
	}
}

func TestEnumRoundTrips(t *testing.T) {
	// Every enum's wire string must survive JSON marshal → unmarshal.
	check := func(t *testing.T, wire string, v json.Marshaler, parseBack func([]byte) (string, error)) {
		t.Helper()
		b, err := v.MarshalJSON()
		if err != nil {
			t.Fatalf("marshal %s: %v", wire, err)
		}
		got, err := parseBack(b)
		if err != nil {
			t.Fatalf("unmarshal %s: %v", string(b), err)
		}
		if got != wire {
			t.Errorf("round-trip %q -> %q", wire, got)
		}
	}
	for _, v := range []interfaceschema.InterfaceType{interfaceschema.Datastream, interfaceschema.Properties} {
		check(t, v.String(), v, func(b []byte) (string, error) {
			var out interfaceschema.InterfaceType
			err := json.Unmarshal(b, &out)
			return out.String(), err
		})
	}
	for _, v := range []interfaceschema.Reliability{
		interfaceschema.ReliabilityUnreliable, interfaceschema.ReliabilityGuaranteed, interfaceschema.ReliabilityUnique,
	} {
		check(t, v.String(), v, func(b []byte) (string, error) {
			var out interfaceschema.Reliability
			err := json.Unmarshal(b, &out)
			return out.String(), err
		})
	}
	for vt := interfaceschema.Double; vt <= interfaceschema.DateTimeArray; vt++ {
		check(t, vt.String(), vt, func(b []byte) (string, error) {
			var out interfaceschema.ValueType
			err := json.Unmarshal(b, &out)
			return out.String(), err
		})
	}
}

func TestReliabilityQoS(t *testing.T) {
	cases := []struct {
		r    interfaceschema.Reliability
		want byte
	}{
		{interfaceschema.ReliabilityUnreliable, 0},
		{interfaceschema.ReliabilityGuaranteed, 1},
		{interfaceschema.ReliabilityUnique, 2},
	}
	for _, tc := range cases {
		if got := tc.r.QoS(); got != tc.want {
			t.Errorf("%v.QoS() = %d, want %d", tc.r, got, tc.want)
		}
	}
}

func TestValueTypeElem(t *testing.T) {
	cases := []struct {
		in, want interfaceschema.ValueType
		isArray  bool
	}{
		{interfaceschema.Double, interfaceschema.Double, false},
		{interfaceschema.DoubleArray, interfaceschema.Double, true},
		{interfaceschema.LongIntegerArray, interfaceschema.LongInteger, true},
		{interfaceschema.DateTimeArray, interfaceschema.DateTime, true},
		{interfaceschema.BinaryBlob, interfaceschema.BinaryBlob, false},
	}
	for _, tc := range cases {
		if got := tc.in.IsArray(); got != tc.isArray {
			t.Errorf("%v.IsArray() = %t, want %t", tc.in, got, tc.isArray)
		}
		if got := tc.in.Elem(); got != tc.want {
			t.Errorf("%v.Elem() = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func FuzzParseInterface(f *testing.F) {
	for _, dir := range []string{"valid", "invalid"} {
		paths, err := filepath.Glob(filepath.Join("testdata", dir, "*.json"))
		if err != nil {
			f.Fatal(err)
		}
		for _, p := range paths {
			b, err := os.ReadFile(p)
			if err != nil {
				f.Fatal(err)
			}
			f.Add(b)
		}
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		iface, err := interfaceschema.ParseInterface(data)
		if err != nil {
			if !errors.Is(err, interfaceschema.ErrInvalid) {
				t.Errorf("rejection %v does not wrap ErrInvalid", err)
			}
			return
		}
		// Accepted documents are structurally valid by contract.
		if iface.Name == "" || len(iface.Mappings) == 0 {
			t.Errorf("accepted interface with empty name or mappings: %+v", iface)
		}
		if iface.Major == 0 && iface.Minor == 0 {
			t.Error("accepted interface with version 0.0")
		}
	})
}
