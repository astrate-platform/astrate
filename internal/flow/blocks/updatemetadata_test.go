package blocks_test

import (
	"testing"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
)

func TestUpdateMetadata_SetThenDelete(t *testing.T) {
	b, err := blocks.UpdateMetadata("um", map[string]any{
		"set_metadata":    map[string]any{"a": "1", "b": "2"},
		"delete_metadata": []any{"a", "gone"},
	}, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}

	in := &flow.Message{
		Key:       "k",
		Timestamp: 1234,
		Type:      flow.TypeString,
		Data:      "v",
		Metadata:  map[string]string{"keep": "k", "a": "old"},
	}
	out, err := b.Process(in)
	if err != nil || len(out) != 1 {
		t.Fatalf("Process: %v %v", out, err)
	}
	got := out[0]
	if _, ok := got.Metadata["a"]; ok {
		t.Errorf("delete should win over set for key a: %v", got.Metadata)
	}
	if got.Metadata["b"] != "2" {
		t.Errorf("set_metadata not applied: %v", got.Metadata)
	}
	if got.Metadata["keep"] != "k" {
		t.Errorf("existing keys must survive: %v", got.Metadata)
	}
	if got.Key != in.Key || got.Timestamp != 1234 || got.Type != flow.TypeString || got.Data != "v" {
		t.Errorf("payload/key/timestamp changed: %+v", got)
	}
	// Input must not be mutated.
	if in.Metadata["b"] != "" || in.Metadata["a"] != "old" {
		t.Errorf("input metadata mutated: %v", in.Metadata)
	}
}

func TestUpdateMetadata_NilPassthrough(t *testing.T) {
	b, err := blocks.UpdateMetadata("um", map[string]any{"set_metadata": map[string]string{"x": "y"}}, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := b.Process(nil)
	if err != nil || len(out) != 0 {
		t.Fatalf("nil: %v %v", out, err)
	}
}

func TestUpdateMetadata_ConstructRejections(t *testing.T) {
	const wantErr = "update_metadata: at least one of set_metadata, delete_metadata is required"
	tests := []struct {
		name string
		bad  map[string]any
		good map[string]any
	}{
		{
			name: "empty config",
			bad:  map[string]any{},
			good: map[string]any{"set_metadata": map[string]any{"k": "v"}},
		},
		{
			name: "nil config",
			bad:  nil,
			good: map[string]any{"delete_metadata": []any{"k"}},
		},
		{
			name: "empty values",
			bad:  map[string]any{"set_metadata": map[string]any{}, "delete_metadata": []any{}},
			good: map[string]any{"set_metadata": map[string]any{}, "delete_metadata": []string{"k"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+"/reject", func(t *testing.T) {
			_, err := blocks.UpdateMetadata("um", tt.bad, flow.Deps{})
			if err == nil || err.Error() != wantErr {
				t.Fatalf("err = %v, want %q", err, wantErr)
			}
		})
		t.Run(tt.name+"/accept", func(t *testing.T) {
			b, err := blocks.UpdateMetadata("um", tt.good, flow.Deps{})
			if err != nil {
				t.Fatalf("twin rejected: %v", err)
			}
			out, err := b.Process(&flow.Message{Key: "k", Type: flow.TypeString})
			if err != nil || len(out) != 1 {
				t.Fatalf("twin Process: %v %v", out, err)
			}
		})
	}
}
