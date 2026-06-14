package engine

import (
	"strings"
	"testing"
	"time"

	promtest "github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// TestParseIntrospection covers the docs/DESIGN.md §3.3 introspection
// format: `;`-separated `name:major:minor` triples, empty = no interfaces,
// any malformed entry rejecting the whole payload.
func TestParseIntrospection(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    map[string]store.InterfaceVersion
		wantErr bool
	}{
		{
			name: "two interfaces",
			in:   "com.ex.Sensors:1:0;com.ex.Geo:0:1",
			want: map[string]store.InterfaceVersion{
				"com.ex.Sensors": {Major: 1, Minor: 0},
				"com.ex.Geo":     {Major: 0, Minor: 1},
			},
		},
		{
			name: "single interface",
			in:   "com.ex.Solo:2:7",
			want: map[string]store.InterfaceVersion{"com.ex.Solo": {Major: 2, Minor: 7}},
		},
		{
			name: "empty payload is an empty introspection",
			in:   "",
			want: map[string]store.InterfaceVersion{},
		},
		{
			name: "duplicate name keeps the last entry",
			in:   "com.ex.A:1:0;com.ex.A:2:1",
			want: map[string]store.InterfaceVersion{"com.ex.A": {Major: 2, Minor: 1}},
		},
		{name: "missing minor", in: "com.ex.A:1", wantErr: true},
		{name: "missing versions", in: "com.ex.A", wantErr: true},
		{name: "empty name", in: ":1:0", wantErr: true},
		{name: "non-numeric major", in: "com.ex.A:x:0", wantErr: true},
		{name: "negative major", in: "com.ex.A:-1:0", wantErr: true},
		{name: "non-numeric minor", in: "com.ex.A:1:y", wantErr: true},
		{name: "trailing semicolon", in: "com.ex.A:1:0;", wantErr: true},
		{name: "wildcard metacharacter in name", in: "com.ex.A#:1:0", wantErr: true},
		{name: "slash in name", in: "com/ex:1:0", wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseIntrospection(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("parseIntrospection(%q) accepted, want error", tc.in)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseIntrospection(%q): %v", tc.in, err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("parseIntrospection(%q) = %v, want %v", tc.in, got, tc.want)
			}
			for name, v := range tc.want {
				if got[name] != v {
					t.Errorf("entry %q = %+v, want %+v", name, got[name], v)
				}
			}
		})
	}
}

// TestIntrospectionHandler drives the wired handler end to end: persistence,
// device-cache update (the §2.6 step-1 gate follows immediately), and the
// broker ACL refresh.
func TestIntrospectionHandler(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})

	// Shrink the introspection to Minimal only.
	ack := &ackCounter{}
	rig.handle(deviceMsg("", "", 2, []byte("com.astrate.test.Minimal:0:1"), ack))
	if !ack.acked() {
		t.Fatal("introspection publish not acknowledged")
	}

	// Persisted.
	fs.mu.Lock()
	dev := fs.devices[deviceKey{realm: realmAlpha, id: devAlpha}]
	persisted := len(dev.Introspection)
	_, hasMinimal := dev.Introspection["com.astrate.test.Minimal"]
	fs.mu.Unlock()
	if persisted != 1 || !hasMinimal {
		t.Errorf("persisted introspection: %d entries (minimal=%t), want exactly Minimal", persisted, hasMinimal)
	}

	// Broker ACL state refreshed.
	if n := port.refreshCount(); n != 1 {
		t.Errorf("broker introspection refreshes: %d, want 1", n)
	}

	// The cached gate now rejects a previously declared interface...
	ack = &ackCounter{}
	rig.handle(deviceMsg("com.astrate.test.AllScalarTypes", "/double", 1, enc(t, 1.5, nil, payload.FormatBSON), ack))
	if got := promtest.ToFloat64(rig.e.met.rejects.WithLabelValues(reasonInterfaceNotDeclared)); got != 1 {
		t.Errorf("rejects[%s] = %v, want 1", reasonInterfaceNotDeclared, got)
	}

	// ...and re-declaring restores it without any store round trip beyond
	// the introspection update itself.
	ack = &ackCounter{}
	rig.handle(deviceMsg("", "", 2, []byte("com.astrate.test.Minimal:0:1;com.astrate.test.AllScalarTypes:1:0"), ack))
	if !ack.acked() {
		t.Fatal("second introspection publish not acknowledged")
	}
	ack = &ackCounter{}
	rig.handle(deviceMsg("com.astrate.test.AllScalarTypes", "/double", 1, enc(t, 1.5, nil, payload.FormatBSON), ack))
	rig.e.flushShard(t.Context(), rig.sh)
	if !ack.acked() {
		t.Error("publish on re-declared interface not accepted")
	}
}

// TestIntrospectionRejects: oversized and malformed payloads are rejected
// under reasonIntrospectionInvalid and consumed.
func TestIntrospectionRejects(t *testing.T) {
	rig, fs, _ := newWiredRig(t, Config{})

	ack := &ackCounter{}
	rig.handle(deviceMsg("", "", 2, []byte(strings.Repeat("x", maxIntrospectionBytes+1)), ack))
	if !ack.acked() {
		t.Error("oversized introspection not consumed")
	}

	ack = &ackCounter{}
	rig.handle(deviceMsg("", "", 2, []byte("com.ex.A:nope:1"), ack))
	if !ack.acked() {
		t.Error("malformed introspection not consumed")
	}

	if got := promtest.ToFloat64(rig.e.met.rejects.WithLabelValues(reasonIntrospectionInvalid)); got != 2 {
		t.Errorf("rejects[%s] = %v, want 2", reasonIntrospectionInvalid, got)
	}

	// Nothing was persisted.
	fs.mu.Lock()
	dev := fs.devices[deviceKey{realm: realmAlpha, id: devAlpha}]
	entries := len(dev.Introspection)
	fs.mu.Unlock()
	if entries != 8 {
		t.Errorf("introspection mutated by rejected payloads: %d entries", entries)
	}
}

// TestIntrospectionParking: a transient UpdateIntrospection failure parks
// the shard (no ack) and the retry persists.
func TestIntrospectionParking(t *testing.T) {
	rig, fs, _ := newWiredRig(t, Config{})
	fs.mu.Lock()
	fs.updateIntroErrs = []error{errTransient}
	fs.mu.Unlock()

	ack := &ackCounter{}
	start := time.Now()
	rig.handle(deviceMsg("", "", 2, []byte("com.astrate.test.Minimal:0:1"), ack))
	if elapsed := time.Since(start); elapsed < parkBackoffStart {
		t.Errorf("handler returned after %s, want >= %s of parking", elapsed, parkBackoffStart)
	}
	if !ack.acked() {
		t.Fatal("introspection not acknowledged after successful retry")
	}
	if ack.n.Load() != 1 {
		t.Errorf("ack fired %d times, want 1", ack.n.Load())
	}

	fs.mu.Lock()
	entries := len(fs.devices[deviceKey{realm: realmAlpha, id: devAlpha}].Introspection)
	fs.mu.Unlock()
	if entries != 1 {
		t.Errorf("retry did not persist: %d entries", entries)
	}
}
