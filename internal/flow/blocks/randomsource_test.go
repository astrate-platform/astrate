package blocks_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
)

func TestRandomSource_ConstructRejections(t *testing.T) {
	tests := []struct {
		name string
		bad  map[string]any
		good map[string]any
		want string
	}{
		{
			name: "missing type",
			bad:  map[string]any{"interval_ms": 5},
			good: map[string]any{"interval_ms": 5, "type": "integer"},
			want: `random_source: unknown type "" (want integer|real|boolean)`,
		},
		{
			name: "unknown type",
			bad:  map[string]any{"type": "blob"},
			good: map[string]any{"type": "real"},
			want: `random_source: unknown type "blob" (want integer|real|boolean)`,
		},
		{
			name: "zero interval",
			bad:  map[string]any{"type": "integer", "interval_ms": 0},
			good: map[string]any{"type": "integer", "interval_ms": 1},
			want: `random_source: interval_ms must be positive`,
		},
		{
			name: "negative interval",
			bad:  map[string]any{"type": "integer", "interval_ms": -5},
			good: map[string]any{"type": "integer", "interval_ms": 5},
			want: `random_source: interval_ms must be positive`,
		},
		{
			name: "min above max",
			bad:  map[string]any{"type": "integer", "min": 10, "max": 9},
			good: map[string]any{"type": "integer", "min": 9, "max": 10},
			want: `random_source: min must be <= max`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+"/reject", func(t *testing.T) {
			_, err := blocks.RandomSource("r", tt.bad, flow.Deps{})
			if err == nil || err.Error() != tt.want {
				t.Fatalf("err = %v, want %q", err, tt.want)
			}
		})
		t.Run(tt.name+"/accept", func(t *testing.T) {
			b, err := blocks.RandomSource("r", tt.good, flow.Deps{})
			if err != nil {
				t.Fatalf("twin rejected: %v", err)
			}
			if s, ok := b.(flow.Stopper); ok {
				s.Stop()
			}
		})
	}
}

func TestRandomSource_EmitsDeterministicIntegerWithinDeadline(t *testing.T) {
	b, err := blocks.RandomSource("r", map[string]any{
		"type":        "integer",
		"min":         7,
		"max":         7,
		"interval_ms": 5,
	}, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	src, ok := b.(flow.Source)
	if !ok {
		t.Fatal("block is not a Source")
	}
	stopper := b.(flow.Stopper)
	defer stopper.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	out, err := src.Emit(ctx)
	if err != nil {
		t.Fatalf("Emit: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("len = %d, want 1", len(out))
	}
	if got := out[0].Data; got != int64(7) {
		t.Errorf("Data = %v (%T), want int64(7)", got, got)
	}
	if out[0].Type != flow.TypeInteger {
		t.Errorf("Type = %v", out[0].Type)
	}
	if out[0].Key != "random" {
		t.Errorf("Key = %q", out[0].Key)
	}
	if out[0].Timestamp == 0 {
		t.Error("timestamp not set")
	}
}

func TestRandomSource_RealAndBooleanShapes(t *testing.T) {
	realB, err := blocks.RandomSource("r", map[string]any{
		"type": "real", "min": 2.0, "max": 3.0, "interval_ms": 1, "key": "k/real",
	}, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := realB.(flow.Source).Emit(context.Background())
	if err != nil || len(out) != 1 {
		t.Fatalf("real Emit: %v %v", out, err)
	}
	v, ok := out[0].Data.(float64)
	if !ok || v < 2.0 || v >= 3.0 {
		t.Errorf("real Data = %v (%T)", out[0].Data, out[0].Data)
	}
	if out[0].Key != "k/real" || out[0].Type != flow.TypeReal {
		t.Errorf("key/type = %q/%v", out[0].Key, out[0].Type)
	}

	boolB, err := blocks.RandomSource("b", map[string]any{
		"type": "boolean", "interval_ms": 1,
	}, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	out, err = boolB.(flow.Source).Emit(context.Background())
	if err != nil || len(out) != 1 {
		t.Fatalf("bool Emit: %v %v", out, err)
	}
	if _, ok := out[0].Data.(bool); !ok || out[0].Type != flow.TypeBoolean {
		t.Errorf("bool Data = %v (%T) type=%v", out[0].Data, out[0].Data, out[0].Type)
	}
}

func TestRandomSource_StopThenEmitReturnsCtxErr(t *testing.T) {
	b, err := blocks.RandomSource("r", map[string]any{"type": "integer", "interval_ms": 5000}, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	src := b.(flow.Source)
	b.(flow.Stopper).Stop()
	b.(flow.Stopper).Stop() // idempotent

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	done := make(chan struct{})
	var emitErr error
	go func() {
		defer close(done)
		_, emitErr = src.Emit(ctx)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Emit did not return after Stop + cancelled ctx")
	}
	if !errors.Is(emitErr, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", emitErr)
	}
	// Stopped source never emits again.
	if out, _ := b.Process(nil); len(out) != 0 {
		t.Fatalf("Process after Stop returned %d messages", len(out))
	}
}
