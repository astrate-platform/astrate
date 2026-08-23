package blocks_test

import (
	"bytes"
	"testing"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
)

// byKey indexes split outputs by their Key (map iteration order is random).
func byKey(msgs []*flow.Message) map[string]*flow.Message {
	m := make(map[string]*flow.Message, len(msgs))
	for _, msg := range msgs {
		m[msg.Key] = msg
	}
	return m
}

func TestSplitMap_FieldTypesHonored(t *testing.T) {
	b, err := blocks.SplitMap("sm", nil, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	in := &flow.Message{
		Key:       "dev1",
		Timestamp: 1700000000000000,
		Metadata:  map[string]string{"kind": "data"},
		Type:      flow.TypeMap,
		Data:      map[string]any{"n": int64(7), "s": "hello"},
		FieldTypes: map[string]flow.DataType{
			"n": flow.TypeInteger,
			"s": flow.TypeString,
		},
	}
	out, err := b.Process(in)
	if err != nil || len(out) != 2 {
		t.Fatalf("Process: %v %v", out, err)
	}
	msgs := byKey(out)
	n, ok := msgs["dev1/n"]
	if !ok {
		t.Fatalf("missing dev1/n: %v", msgs)
	}
	if n.Type != flow.TypeInteger || n.Data != int64(7) {
		t.Errorf("n: type=%v data=%v", n.Type, n.Data)
	}
	s := msgs["dev1/s"]
	if s == nil || s.Type != flow.TypeString || s.Data != "hello" {
		t.Errorf("s: %+v", s)
	}
	for k, m := range msgs {
		if m.Timestamp != 1700000000000000 {
			t.Errorf("%s: timestamp inherited", k)
		}
		if m.Metadata["kind"] != "data" {
			t.Errorf("%s: metadata not cloned", k)
		}
		if m.FieldTypes != nil || m.FieldSubtypes != nil {
			t.Errorf("%s: field maps should be nil", k)
		}
	}
}

func TestSplitMap_TypeInference(t *testing.T) {
	b, err := blocks.SplitMap("sm", nil, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	in := &flow.Message{
		Key:  "k",
		Type: flow.TypeMap,
		Data: map[string]any{
			"b": true,
			"i": int64(3),
			"r": 2.5,
			"s": "txt",
			"x": []any{1},
		},
	}
	out, err := b.Process(in)
	if err != nil || len(out) != 5 {
		t.Fatalf("Process: %v %v", out, err)
	}
	msgs := byKey(out)
	if m := msgs["k/b"]; m == nil || m.Type != flow.TypeBoolean || m.Data != true {
		t.Errorf("b: %+v", m)
	}
	if m := msgs["k/i"]; m == nil || m.Type != flow.TypeInteger || m.Data != int64(3) {
		t.Errorf("i: %+v", m)
	}
	if m := msgs["k/r"]; m == nil || m.Type != flow.TypeReal || m.Data != 2.5 {
		t.Errorf("r: %+v", m)
	}
	if m := msgs["k/s"]; m == nil || m.Type != flow.TypeString || m.Data != "txt" {
		t.Errorf("s: %+v", m)
	}
	// Fallback: neither bool nor int64 nor float64 → string via fmt.Sprint.
	if m := msgs["k/x"]; m == nil || m.Type != flow.TypeString || m.Data != "[1]" {
		t.Errorf("x: %+v", m)
	}
}

func TestSplitMap_BinaryFieldSubtype(t *testing.T) {
	b, err := blocks.SplitMap("sm", nil, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	in := &flow.Message{
		Key:           "k",
		Type:          flow.TypeMap,
		Data:          map[string]any{"blob": []byte{1, 2}},
		FieldTypes:    map[string]flow.DataType{"blob": flow.TypeBinary},
		FieldSubtypes: map[string]string{"blob": "application/octet-stream"},
	}
	out, err := b.Process(in)
	if err != nil || len(out) != 1 {
		t.Fatalf("Process: %v %v", out, err)
	}
	got := out[0]
	if got.Key != "k/blob" || got.Type != flow.TypeBinary {
		t.Errorf("key/type = %q/%v", got.Key, got.Type)
	}
	if got.Subtype != "application/octet-stream" {
		t.Errorf("Subtype = %q", got.Subtype)
	}
	if !bytes.Equal(got.Data.([]byte), []byte{1, 2}) {
		t.Errorf("Data = %v", got.Data)
	}
}

func TestSplitMap_MetadataClonedPerOutput(t *testing.T) {
	b, err := blocks.SplitMap("sm", nil, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	in := &flow.Message{
		Key:      "dev",
		Type:     flow.TypeMap,
		Metadata: map[string]string{"kind": "data"},
		Data:     map[string]any{"a": int64(1), "b": int64(2)},
	}
	out, err := b.Process(in)
	if err != nil || len(out) != 2 {
		t.Fatalf("Process: %v %v", out, err)
	}
	// Mutate output[0]'s Metadata after Process returns; output[1] must be
	// unaffected (each output carries its own clone).
	out[0].Metadata["kind"] = "mutated"
	delete(out[0].Metadata, "kind")
	if out[1].Metadata["kind"] != "data" || len(out[1].Metadata) != 1 {
		t.Errorf("output[1] metadata changed when output[0] was mutated: %v", out[1].Metadata)
	}
}

func TestSplitMap_NonMapDropped(t *testing.T) {
	b, err := blocks.SplitMap("sm", nil, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := b.Process(nil)
	if err != nil || len(out) != 0 {
		t.Fatalf("nil: %v %v", out, err)
	}
	out, err = b.Process(&flow.Message{Key: "k", Type: flow.TypeString, Data: "nope"})
	if err != nil || len(out) != 0 {
		t.Fatalf("string: %v %v", out, err)
	}
}

func TestSplitMap_CustomKeyTemplate(t *testing.T) {
	b, err := blocks.SplitMap("sm", map[string]any{"key_template": "{field}@{key}"}, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := b.Process(&flow.Message{
		Key:  "dev",
		Type: flow.TypeMap,
		Data: map[string]any{"a": int64(1), "b": int64(2)},
		FieldTypes: map[string]flow.DataType{
			"a": flow.TypeInteger,
			"b": flow.TypeInteger,
		},
	})
	if err != nil || len(out) != 2 {
		t.Fatalf("Process: %v %v", out, err)
	}
	msgs := byKey(out)
	if msgs["a@dev"] == nil || msgs["b@dev"] == nil {
		t.Errorf("keys = %v", msgs)
	}
}
