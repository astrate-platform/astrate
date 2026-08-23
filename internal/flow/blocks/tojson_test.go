package blocks_test

import (
	"bytes"
	"encoding/base64"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
)

func TestToJSON_MapPayload(t *testing.T) {
	b, err := blocks.ToJSON("j", nil, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	in := &flow.Message{
		Key:       "dev1/temp",
		Timestamp: 1700000000000000,
		Metadata:  map[string]string{"kind": "data"},
		Type:      flow.TypeMap,
		Data:      map[string]any{"n": int64(42), "s": "hi"},
		FieldTypes: map[string]flow.DataType{
			"n": flow.TypeInteger,
			"s": flow.TypeString,
		},
	}
	out, err := b.Process(in)
	if err != nil || len(out) != 1 {
		t.Fatalf("Process: %v %v", out, err)
	}
	got := out[0]
	if !bytes.Equal(got.Data.([]byte), []byte(`{"n":42,"s":"hi"}`)) {
		t.Errorf("Data = %q", got.Data)
	}
	if got.Type != flow.TypeBinary {
		t.Errorf("Type = %v, want binary", got.Type)
	}
	if got.Subtype != "application/json" {
		t.Errorf("Subtype = %q", got.Subtype)
	}
	if got.Key != in.Key || got.Timestamp != in.Timestamp {
		t.Errorf("key/timestamp not preserved: %+v", got)
	}
	if got.Metadata["kind"] != "data" {
		t.Errorf("metadata not preserved: %v", got.Metadata)
	}
	if got.FieldTypes != nil || got.FieldSubtypes != nil {
		t.Errorf("field maps should be cleared")
	}
}

func TestToJSON_DatetimeAndBinaryScalars(t *testing.T) {
	b, err := blocks.ToJSON("j", nil, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	ts := time.Date(2023, 11, 14, 22, 13, 20, 123456789, time.UTC)

	out, err := b.Process(&flow.Message{Key: "k", Type: flow.TypeDatetime, Data: ts})
	if err != nil || len(out) != 1 {
		t.Fatalf("datetime: %v %v", out, err)
	}
	if got := string(out[0].Data.([]byte)); got != `"2023-11-14T22:13:20.123456789Z"` {
		t.Errorf("datetime Data = %s", got)
	}

	raw := []byte{0xde, 0xad}
	out, err = b.Process(&flow.Message{Key: "k", Type: flow.TypeBinary, Subtype: "application/octet-stream", Data: raw})
	if err != nil || len(out) != 1 {
		t.Fatalf("binary: %v %v", out, err)
	}
	want := `"` + base64.StdEncoding.EncodeToString(raw) + `"`
	if got := string(out[0].Data.([]byte)); got != want {
		t.Errorf("binary Data = %s, want %s", got, want)
	}
	if out[0].Subtype != "application/json" {
		t.Errorf("binary subtype overwritten: %q", out[0].Subtype)
	}
}

func TestToJSON_NilPassthrough(t *testing.T) {
	b, err := blocks.ToJSON("j", nil, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := b.Process(nil)
	if err != nil || len(out) != 0 {
		t.Fatalf("nil: %v %v", out, err)
	}
}

func TestToJSON_InputNotMutated(t *testing.T) {
	b, err := blocks.ToJSON("j", nil, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	in := &flow.Message{
		Key:        "k",
		Type:       flow.TypeMap,
		Data:       map[string]any{"n": int64(1)},
		FieldTypes: map[string]flow.DataType{"n": flow.TypeInteger},
	}
	if _, err := b.Process(in); err != nil {
		t.Fatal(err)
	}
	if in.Type != flow.TypeMap {
		t.Errorf("input type mutated: %v", in.Type)
	}
	if _, ok := in.Data.(map[string]any); !ok {
		t.Errorf("input data mutated: %T", in.Data)
	}
	if in.FieldTypes["n"] != flow.TypeInteger {
		t.Errorf("input field types mutated: %v", in.FieldTypes)
	}
	if in.Subtype != "" {
		t.Errorf("input subtype mutated: %q", in.Subtype)
	}
}
