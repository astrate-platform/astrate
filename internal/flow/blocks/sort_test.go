package blocks_test

import (
	"testing"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
)

func sortBlock(t *testing.T, config map[string]any) flow.Block {
	t.Helper()
	b, err := blocks.Sort("s", config, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func intMsg(key string, ts int64) *flow.Message {
	return &flow.Message{Key: key, Timestamp: ts, Type: flow.TypeInteger, Data: ts}
}

func TestSort_FlushesAscendingOnArrival(t *testing.T) {
	b := sortBlock(t, map[string]any{"window_ms": 10})

	// Two arrivals within the window of each other: nothing flushed yet.
	out, err := b.Process(intMsg("k", 200_000))
	if err != nil || len(out) != 0 {
		t.Fatalf("first: %v %v", out, err)
	}
	out, err = b.Process(intMsg("k", 205_000))
	if err != nil || len(out) != 0 {
		t.Fatalf("second (within window): %v %v", out, err)
	}

	// A much later arrival triggers the flush of both earlier messages,
	// ascending; the newest message itself stays buffered.
	out, err = b.Process(intMsg("k", 1_000_000))
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("flushed %d messages, want 2", len(out))
	}
	if out[0].Timestamp != 200_000 || out[1].Timestamp != 205_000 {
		t.Errorf("flush order = %d, %d", out[0].Timestamp, out[1].Timestamp)
	}

	// The newest message is never emitted early.
	out, _ = b.Process(nil)
	if len(out) != 0 {
		t.Errorf("nil produced %v", out)
	}

	// Out-of-order arrival inside the window edge is held, then flushed
	// by a newer timestamp.
	if out, _ = b.Process(intMsg("k", 995_000)); len(out) != 0 {
		t.Fatalf("premature flush: %v", out)
	}
	out, _ = b.Process(intMsg("k", 5_000_000))
	if len(out) != 2 || out[0].Timestamp != 995_000 || out[1].Timestamp != 1_000_000 {
		t.Errorf("second flush = %v", out)
	}
}

func TestSort_Dedup(t *testing.T) {
	b := sortBlock(t, map[string]any{"window_ms": 10, "dedup": true})
	same := func() *flow.Message {
		m := intMsg("k", 100_000)
		m.Metadata = map[string]string{"a": "b"}
		return m
	}
	if _, err := b.Process(same()); err != nil {
		t.Fatal(err)
	}
	// Identical wire encoding → dropped.
	out, err := b.Process(same())
	if err != nil || len(out) != 0 {
		t.Fatalf("duplicate not dropped: %v %v", out, err)
	}
	// A differing message is kept.
	differing := same()
	differing.Timestamp = 105_000
	out, err = b.Process(differing)
	if err != nil || len(out) != 0 {
		t.Fatalf("differing dropped or errored: %v %v", out, err)
	}
	// Flush both buffered messages via a later arrival.
	out, _ = b.Process(intMsg("k", 1_000_000))
	if len(out) != 2 || out[0].Timestamp != 100_000 || out[1].Timestamp != 105_000 {
		t.Fatalf("flush = %v", out)
	}
	// Same payload after leaving the buffer is accepted again (dedup only
	// compares against currently-buffered entries); it flushes immediately
	// because it trails the newest buffered timestamp by far.
	out, err = b.Process(same())
	if err != nil || len(out) != 1 || out[0].Timestamp != 100_000 {
		t.Fatalf("re-arrival after eviction: %v %v", out, err)
	}
}

func TestSort_NilPassthrough(t *testing.T) {
	b := sortBlock(t, nil)
	out, err := b.Process(nil)
	if err != nil || len(out) != 0 {
		t.Fatalf("nil: %v %v", out, err)
	}
}

func TestSort_ConstructRejections(t *testing.T) {
	const wantErr = "sort: window_ms must be non-negative"
	tests := []struct {
		name string
		bad  map[string]any
		good map[string]any
	}{
		{
			name: "negative window",
			bad:  map[string]any{"window_ms": -1},
			good: map[string]any{"window_ms": 0},
		},
		{
			name: "non-numeric window",
			bad:  map[string]any{"window_ms": "soon"},
			good: map[string]any{"window_ms": 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+"/reject", func(t *testing.T) {
			_, err := blocks.Sort("s", tt.bad, flow.Deps{})
			if err == nil || err.Error() != wantErr {
				t.Fatalf("err = %v, want %q", err, wantErr)
			}
		})
		t.Run(tt.name+"/accept", func(t *testing.T) {
			sortBlock(t, tt.good)
		})
	}
}

func TestSort_EqualTimestampsKeepArrivalOrder(t *testing.T) {
	b := sortBlock(t, map[string]any{"window_ms": 10})
	first := intMsg("k", 500_000)
	second := &flow.Message{Key: "k", Timestamp: 500_000, Type: flow.TypeString, Data: "later"}
	if _, err := b.Process(first); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Process(second); err != nil {
		t.Fatal(err)
	}
	out, _ := b.Process(intMsg("k", 5_000_000))
	if len(out) != 2 || out[0] != first || out[1] != second {
		t.Errorf("equal-ts order broken: %v", out)
	}
}
