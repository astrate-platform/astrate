package interfaceschema_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// stubResolver maps endpoints to sequential IDs and records the interface
// lookup it served.
type stubResolver struct {
	interfaceID int64
	endpointIDs map[string]int64
	failOn      string
}

func (r *stubResolver) ResolveInterface(string, int) (int64, error) {
	return r.interfaceID, nil
}

func (r *stubResolver) ResolveEndpoint(endpoint string) (int64, error) {
	if endpoint == r.failOn {
		return 0, errors.New("stub resolver failure")
	}
	id, ok := r.endpointIDs[endpoint]
	if !ok {
		return 0, fmt.Errorf("unknown endpoint %q", endpoint)
	}
	return id, nil
}

func parseFixture(t testing.TB, name string) *interfaceschema.Interface {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "valid", name))
	if err != nil {
		t.Fatal(err)
	}
	iface, err := interfaceschema.ParseInterface(data)
	if err != nil {
		t.Fatalf("ParseInterface(%s): %v", name, err)
	}
	return iface
}

func TestCompileIndividual(t *testing.T) {
	iface := parseFixture(t, "org.astarte-platform.genericsensors.Values.json")
	ids := &stubResolver{
		interfaceID: 42,
		endpointIDs: map[string]int64{"/%{sensor_id}/value": 7},
	}
	ci, err := interfaceschema.Compile(iface, ids)
	if err != nil {
		t.Fatal(err)
	}
	if ci.ID != 42 {
		t.Errorf("ID = %d, want 42", ci.ID)
	}
	if ci.Name != "org.astarte-platform.genericsensors.Values" || ci.Major != 1 || ci.Minor != 0 {
		t.Errorf("identity = %s v%d.%d", ci.Name, ci.Major, ci.Minor)
	}
	if ci.Type != interfaceschema.Datastream || ci.Ownership != interfaceschema.OwnershipDevice {
		t.Errorf("type/ownership = %v/%v", ci.Type, ci.Ownership)
	}
	if ci.Aggregation != interfaceschema.AggregationIndividual {
		t.Errorf("Aggregation = %v, want individual", ci.Aggregation)
	}
	if ci.ObjectLeaves != nil {
		t.Errorf("ObjectLeaves = %v, want nil for individual aggregation", ci.ObjectLeaves)
	}

	m, ok := ci.Trie.Match("/sensor1/value")
	if !ok {
		t.Fatal("Match(/sensor1/value) missed")
	}
	if m.EndpointID != 7 {
		t.Errorf("EndpointID = %d, want 7", m.EndpointID)
	}
	if m.ValueType != interfaceschema.Double {
		t.Errorf("ValueType = %v, want double", m.ValueType)
	}
	if m.Reliability != 0 {
		t.Errorf("Reliability(QoS) = %d, want 0", m.Reliability)
	}
	if !m.ExplicitTimestamp {
		t.Error("ExplicitTimestamp = false, want true")
	}
	if _, ok := ci.Trie.Match("/value"); ok {
		t.Error("Match(/value) hit, want miss (wrong depth)")
	}
	if _, ok := ci.Trie.Match("/sensor1/value/x"); ok {
		t.Error("Match(/sensor1/value/x) hit, want miss (too deep)")
	}
}

func TestCompileObjectAggregation(t *testing.T) {
	iface := parseFixture(t, "org.astarte-platform.genericsensors.Geolocation.json")
	ci, err := interfaceschema.Compile(iface, nil)
	if err != nil {
		t.Fatal(err)
	}
	if ci.ID != 0 {
		t.Errorf("ID = %d, want 0 with nil resolver", ci.ID)
	}
	wantLeaves := []string{"latitude", "longitude", "altitude", "accuracy", "altitudeAccuracy", "heading", "speed"}
	if len(ci.ObjectLeaves) != len(wantLeaves) {
		t.Fatalf("len(ObjectLeaves) = %d, want %d", len(ci.ObjectLeaves), len(wantLeaves))
	}
	for _, leaf := range wantLeaves {
		m, ok := ci.ObjectLeaves[leaf]
		if !ok {
			t.Errorf("ObjectLeaves missing %q", leaf)
			continue
		}
		if m.ValueType != interfaceschema.Double || !m.ExplicitTimestamp {
			t.Errorf("leaf %q: ValueType=%v ExplicitTimestamp=%t", leaf, m.ValueType, m.ExplicitTimestamp)
		}
	}
	// Full-path matching still works for object interfaces (AppEngine reads).
	m, ok := ci.Trie.Match("/gps0/latitude")
	if !ok {
		t.Fatal("Match(/gps0/latitude) missed")
	}
	if m != ci.ObjectLeaves["latitude"] {
		t.Error("trie leaf and ObjectLeaves entry are not the same mapping")
	}
}

func TestCompileDurationsAndQoS(t *testing.T) {
	iface := parseFixture(t, "com.astrate.test.AllScalarTypes.json")
	ci, err := interfaceschema.Compile(iface, nil)
	if err != nil {
		t.Fatal(err)
	}

	str, ok := ci.Trie.Match("/string")
	if !ok {
		t.Fatal("Match(/string) missed")
	}
	if str.Expiry != 3600*time.Second {
		t.Errorf("Expiry = %v, want 1h", str.Expiry)
	}
	if str.DBRetentionTTL != 31536000*time.Second {
		t.Errorf("DBRetentionTTL = %v, want 365d", str.DBRetentionTTL)
	}
	if str.Retention != interfaceschema.RetentionStored {
		t.Errorf("Retention = %v, want stored", str.Retention)
	}

	integer, _ := ci.Trie.Match("/integer")
	if integer == nil || integer.Reliability != 1 {
		t.Errorf("integer Reliability = %+v, want QoS 1", integer)
	}
	boolean, _ := ci.Trie.Match("/boolean")
	if boolean == nil || boolean.Reliability != 2 {
		t.Errorf("boolean Reliability = %+v, want QoS 2", boolean)
	}
	dbl, _ := ci.Trie.Match("/double")
	if dbl == nil || dbl.Expiry != 0 || dbl.DBRetentionTTL != 0 {
		t.Errorf("double durations = %+v, want zero", dbl)
	}
}

func TestCompileEveryValidFixture(t *testing.T) {
	for name, data := range readFixtures(t, filepath.Join("testdata", "valid")) {
		t.Run(name, func(t *testing.T) {
			iface, err := interfaceschema.ParseInterface(data)
			if err != nil {
				t.Fatal(err)
			}
			ci, err := interfaceschema.Compile(iface, nil)
			if err != nil {
				t.Fatalf("Compile: %v", err)
			}
			// Every declared endpoint must be reachable through the trie
			// with its placeholders substituted.
			for _, m := range iface.Mappings {
				path := concretePath(m.Endpoint)
				if _, ok := ci.Trie.Match(path); !ok {
					t.Errorf("declared endpoint %q unreachable as %q", m.Endpoint, path)
				}
			}
		})
	}
}

// concretePath substitutes every %{placeholder} with a literal value.
func concretePath(endpoint string) string {
	out := ""
	for _, seg := range splitSegments(endpoint) {
		if len(seg) > 1 && seg[0] == '%' {
			seg = "concrete0"
		}
		out += "/" + seg
	}
	return out
}

// splitSegments is a test-local naive splitter (endpoints here are valid).
func splitSegments(endpoint string) []string {
	var segs []string
	rest := endpoint[1:]
	for len(rest) > 0 {
		i := 0
		for i < len(rest) && rest[i] != '/' {
			i++
		}
		segs = append(segs, rest[:i])
		if i == len(rest) {
			break
		}
		rest = rest[i+1:]
	}
	return segs
}

func TestCompileResolverErrors(t *testing.T) {
	iface := parseFixture(t, "com.astrate.test.AllScalarTypes.json")
	ids := &stubResolver{failOn: "/boolean", endpointIDs: map[string]int64{
		"/double": 1, "/integer": 2, "/longinteger": 4, "/string": 5, "/binaryblob": 6, "/datetime": 7,
	}}
	if _, err := interfaceschema.Compile(iface, ids); err == nil {
		t.Error("Compile with failing resolver succeeded, want error")
	}
}

func TestCompileNil(t *testing.T) {
	if _, err := interfaceschema.Compile(nil, nil); err == nil {
		t.Error("Compile(nil) succeeded, want error")
	}
}
