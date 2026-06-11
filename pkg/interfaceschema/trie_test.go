package interfaceschema

import (
	"strings"
	"testing"
	"time"
)

// buildTestTrie declares a representative endpoint set:
//
//	/exact
//	/a/b            (exact sibling of the parametric branch below)
//	/a/%{p}         (leaf and interior node at once: see /a/%{p}/deep)
//	/a/%{p}/deep
//	/%{q}/c
//	/obj/%{id}/value
type testTrie struct {
	trie     *EndpointTrie
	mappings map[string]*CompiledMapping
}

func buildTestTrie(t testing.TB) testTrie {
	t.Helper()
	endpoints := []string{
		"/exact",
		"/a/b",
		"/a/%{p}",
		"/a/%{p}/deep",
		"/%{q}/c",
		"/obj/%{id}/value",
	}
	trie := NewEndpointTrie()
	mappings := make(map[string]*CompiledMapping, len(endpoints))
	for i, ep := range endpoints {
		m := &CompiledMapping{EndpointID: int64(i + 1), ValueType: Double}
		mappings[ep] = m
		if err := trie.Add(ep, m); err != nil {
			t.Fatalf("Add(%q): %v", ep, err)
		}
	}
	return testTrie{trie: trie, mappings: mappings}
}

func TestMatchTable(t *testing.T) {
	tt := buildTestTrie(t)
	longSeg := strings.Repeat("x", 256)
	overlongSeg := strings.Repeat("x", 257)

	cases := []struct {
		name string
		path string
		want string // endpoint key in tt.mappings; "" = no match
	}{
		{"exact literal", "/exact", "/exact"},
		{"exact beats parametric", "/a/b", "/a/b"},
		{"parametric fallback", "/a/z", "/a/%{p}"},
		{"parametric numeric value", "/a/42", "/a/%{p}"},
		{"deep under parametric", "/a/anything/deep", "/a/%{p}/deep"},
		{"root parametric", "/sensor1/c", "/%{q}/c"},
		{"mixed literal parametric literal", "/obj/dev-7/value", "/obj/%{id}/value"},
		{"placeholder at max length", "/a/" + longSeg, "/a/%{p}"},
		{"placeholder value with dots and colons", "/a/192.168.0.1:8883", "/a/%{p}"},

		{"miss unknown root", "/nope", ""},
		{"miss too deep", "/exact/extra", ""},
		{"miss too shallow", "/obj/dev-7", ""},
		{"miss interior without leaf", "/obj", ""},
		{"miss wrong leaf", "/obj/dev-7/values", ""},
		{"empty path", "", ""},
		{"bare slash", "/", ""},
		{"not rooted", "a/b", ""},
		{"empty middle segment", "/a//deep", ""},
		{"trailing slash", "/a/b/", ""},
		{"overlong placeholder value", "/a/" + overlongSeg, ""},
		{"plus in placeholder value", "/a/x+y", ""},
		{"hash in placeholder value", "/a/x#y", ""},
		{"wildcard alone", "/a/+", ""},
		{"hash under root param", "/#/c", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tt.trie.Match(tc.path)
			if tc.want == "" {
				if ok || got != nil {
					t.Fatalf("Match(%q) = %+v, want miss", tc.path, got)
				}
				return
			}
			if !ok || got == nil {
				t.Fatalf("Match(%q) missed, want %s", tc.path, tc.want)
			}
			if want := tt.mappings[tc.want]; got != want {
				t.Errorf("Match(%q) = EndpointID %d, want %s (EndpointID %d)",
					tc.path, got.EndpointID, tc.want, want.EndpointID)
			}
		})
	}
}

// TestMatchBacktracking documents that a dead-ended exact branch falls back
// to the parametric sibling: /a/b is a declared leaf, but /a/b/deep is only
// reachable through /a/%{p}/deep.
func TestMatchBacktracking(t *testing.T) {
	tt := buildTestTrie(t)
	got, ok := tt.trie.Match("/a/b/deep")
	if !ok {
		t.Fatal("Match(/a/b/deep) missed, want /a/%{p}/deep via backtracking")
	}
	if want := tt.mappings["/a/%{p}/deep"]; got != want {
		t.Errorf("Match(/a/b/deep) = EndpointID %d, want %d", got.EndpointID, want.EndpointID)
	}
}

func TestAddRejections(t *testing.T) {
	trie := NewEndpointTrie()
	m := &CompiledMapping{ValueType: Double}
	if err := trie.Add("/a/%{p}/x", m); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name     string
		endpoint string
		wantSub  string
	}{
		{"duplicate", "/a/%{p}/x", "duplicate endpoint"},
		{"placeholder rename", "/a/%{other}/y", "conflicting placeholder"},
		{"invalid syntax", "a/b", "must start with '/'"},
		{"empty segment", "/a//b", "empty segment"},
		{"partial placeholder", "/a%{p}/x", "whole segment"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := trie.Add(tc.endpoint, m)
			if err == nil {
				t.Fatalf("Add(%q) succeeded, want error containing %q", tc.endpoint, tc.wantSub)
			}
			if !strings.Contains(err.Error(), tc.wantSub) {
				t.Errorf("Add(%q) error %q does not contain %q", tc.endpoint, err.Error(), tc.wantSub)
			}
		})
	}

	t.Run("nil mapping", func(t *testing.T) {
		if err := trie.Add("/fresh", nil); err == nil {
			t.Error("Add(nil mapping) succeeded, want error")
		}
	})

	t.Run("same placeholder name merges", func(t *testing.T) {
		if err := trie.Add("/a/%{p}/y", m); err != nil {
			t.Errorf("Add with same placeholder name: %v", err)
		}
	})
}

func TestNilAndEmptyTrie(t *testing.T) {
	var nilTrie *EndpointTrie
	if _, ok := nilTrie.Match("/a"); ok {
		t.Error("nil trie matched")
	}
	if _, ok := NewEndpointTrie().Match("/a"); ok {
		t.Error("empty trie matched")
	}
}

// TestMatchZeroAllocs is a hard CI assertion of the §2.6 hot-path contract:
// Match performs zero heap allocations, hit or miss, shallow or deep.
func TestMatchZeroAllocs(t *testing.T) {
	tt := buildTestTrie(t)
	paths := []string{
		"/exact",
		"/a/b",
		"/a/sensor-with-a-long-parametric-value-0123456789",
		"/a/anything/deep",
		"/obj/dev-7/value",
		"/nope/miss/entirely",
		"/a/x+y", // hygiene-rejected parametric value
	}
	allocs := testing.AllocsPerRun(1000, func() {
		for _, p := range paths {
			tt.trie.Match(p)
		}
	})
	if allocs != 0 {
		t.Errorf("Match allocates: %v allocs/run, want 0", allocs)
	}
}

func BenchmarkMatch(b *testing.B) {
	tt := buildTestTrie(b)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		tt.trie.Match("/obj/dev-7/value")
	}
}

func BenchmarkMatchMiss(b *testing.B) {
	tt := buildTestTrie(b)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		tt.trie.Match("/obj/dev-7/values")
	}
}

// Compile-time guard: CompiledMapping stays the §2.6 shape.
var _ = CompiledMapping{
	EndpointID:        0,
	ValueType:         Double,
	Reliability:       0,
	Retention:         RetentionDiscard,
	Expiry:            time.Duration(0),
	ExplicitTimestamp: false,
	AllowUnset:        false,
	DBRetentionTTL:    time.Duration(0),
}
