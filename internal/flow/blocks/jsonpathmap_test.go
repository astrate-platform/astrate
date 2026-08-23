package blocks_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
)

func newJSONPathMap(t *testing.T, config map[string]any) flow.Block {
	t.Helper()
	b, err := blocks.JSONPathMap("jpm", config, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestJSONPathMap_ConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]any
		wantErr bool
	}{
		{"missing template", nil, true},
		{"empty template", map[string]any{"template": ""}, true},
		{"acceptance twin", map[string]any{"template": `{"v": $MESSAGE}`}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := blocks.JSONPathMap("jpm", tt.config, flow.Deps{})
			if tt.wantErr {
				if err == nil || !strings.Contains(err.Error(), "template is required") {
					t.Fatalf("err = %v, want template is required", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestJSONPathMap_HappyPath(t *testing.T) {
	b := newJSONPathMap(t, map[string]any{
		"template": `{"n": $MESSAGE.count, "who": "$METADATA.user", "src": $MESSAGE}`,
	})
	in := &flow.Message{
		Key:        "dev1",
		Timestamp:  1700000000000000,
		Metadata:   map[string]string{"user": "alice"},
		Type:       flow.TypeMap,
		Data:       map[string]any{"count": int64(3), "tags": []any{"a", "b"}},
		FieldTypes: map[string]flow.DataType{"count": flow.TypeInteger},
	}
	out, err := b.Process(in)
	if err != nil || len(out) != 1 {
		t.Fatalf("Process: %v %v", out, err)
	}
	got := out[0]
	data, ok := got.Data.(map[string]any)
	if !ok {
		t.Fatalf("Data = %T, want map[string]any", got.Data)
	}
	if n := data["n"]; n != float64(3) {
		t.Errorf("n = %v (%T), want float64(3)", n, n)
	}
	if who := data["who"]; who != "alice" {
		t.Errorf("who = %v, want \"alice\"", who)
	}
	wantSrc := `{"count":3,"tags":["a","b"]}`
	if src := data["src"]; src != wantSrc {
		t.Errorf("src = %q, want %q", src, wantSrc)
	}
	wantTypes := map[string]flow.DataType{
		"n":   flow.TypeReal,
		"who": flow.TypeString,
		"src": flow.TypeString,
	}
	if !reflect.DeepEqual(got.FieldTypes, wantTypes) {
		t.Errorf("FieldTypes = %v, want %v", got.FieldTypes, wantTypes)
	}
	if got.Type != flow.TypeMap || got.FieldSubtypes != nil {
		t.Errorf("type/subtype = %v/%v", got.Type, got.FieldSubtypes)
	}
	if got.Key != "dev1" || got.Timestamp != 1700000000000000 {
		t.Errorf("key/timestamp not inherited: %+v", got)
	}
	if got.Metadata["user"] != "alice" {
		t.Errorf("metadata not cloned: %v", got.Metadata)
	}
}

func TestJSONPathMap_DeepPath(t *testing.T) {
	b := newJSONPathMap(t, map[string]any{"template": `{"v": $MESSAGE.a.b}`})
	in := &flow.Message{
		Key:  "k",
		Type: flow.TypeMap,
		Data: map[string]any{"a": map[string]any{"b": int64(7)}},
	}
	out, err := b.Process(in)
	if err != nil || len(out) != 1 {
		t.Fatalf("Process: %v %v", out, err)
	}
	data := out[0].Data.(map[string]any)
	if v := data["v"]; v != float64(7) {
		t.Errorf("v = %v (%T), want float64(7)", v, v)
	}
	if out[0].FieldTypes["v"] != flow.TypeReal {
		t.Errorf("field type = %v, want real", out[0].FieldTypes["v"])
	}
}

func TestJSONPathMap_ResolveErrors(t *testing.T) {
	tests := []struct {
		name     string
		template string
		msg      *flow.Message
		wantErr  string
	}{
		{
			name:     "missing top-level key",
			template: `{"v": $MESSAGE.nope}`,
			msg: &flow.Message{
				Key:  "k",
				Type: flow.TypeMap,
				Data: map[string]any{"count": int64(1)},
			},
			wantErr: "cannot resolve $MESSAGE.nope",
		},
		{
			name:     "scalar mid-path",
			template: `{"v": $MESSAGE.s.x}`,
			msg: &flow.Message{
				Key:  "k",
				Type: flow.TypeMap,
				Data: map[string]any{"s": "txt"},
			},
			wantErr: "cannot resolve $MESSAGE.s.x",
		},
		{
			name:     "missing metadata name",
			template: `{"v": "$METADATA.absent"}`,
			msg: &flow.Message{
				Key:      "k",
				Metadata: map[string]string{"user": "alice"},
				Type:     flow.TypeMap,
				Data:     map[string]any{},
			},
			wantErr: "cannot resolve $METADATA.absent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newJSONPathMap(t, map[string]any{"template": tt.template})
			out, err := b.Process(tt.msg)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want containing %q", err, tt.wantErr)
			}
			if len(out) != 0 {
				t.Fatalf("outputs on error = %v, want none", out)
			}
		})
	}
}

func TestJSONPathMap_ScalarPayloadBareMessage(t *testing.T) {
	b := newJSONPathMap(t, map[string]any{"template": `{"v": $MESSAGE}`})
	in := &flow.Message{Key: "k", Type: flow.TypeInteger, Data: int64(5)}
	out, err := b.Process(in)
	if err != nil || len(out) != 1 {
		t.Fatalf("Process: %v %v", out, err)
	}
	data := out[0].Data.(map[string]any)
	if v := data["v"]; v != float64(5) {
		t.Errorf("v = %v (%T), want float64(5)", v, v)
	}
}

func TestJSONPathMap_RepeatedPlaceholder(t *testing.T) {
	b := newJSONPathMap(t, map[string]any{
		"template": `{"a": $MESSAGE.count, "b": $MESSAGE.count}`,
	})
	in := &flow.Message{
		Key:  "k",
		Type: flow.TypeMap,
		Data: map[string]any{"count": int64(3)},
	}
	out, err := b.Process(in)
	if err != nil || len(out) != 1 {
		t.Fatalf("Process: %v %v", out, err)
	}
	data := out[0].Data.(map[string]any)
	if data["a"] != float64(3) || data["b"] != float64(3) {
		t.Errorf("a/b = %v/%v, want both 3", data["a"], data["b"])
	}
}

func TestJSONPathMap_LiteralAdjacency(t *testing.T) {
	b := newJSONPathMap(t, map[string]any{"template": `{"v": "n=$MESSAGES"}`})
	in := &flow.Message{Key: "k", Type: flow.TypeInteger, Data: int64(5)}
	out, err := b.Process(in)
	if err != nil || len(out) != 1 {
		t.Fatalf("Process: %v %v", out, err)
	}
	data := out[0].Data.(map[string]any)
	if v := data["v"]; v != "n=5S" {
		t.Errorf("v = %q, want \"n=5S\" ($MESSAGE + literal S)", v)
	}
}

func TestJSONPathMap_InvalidResultingJSON(t *testing.T) {
	b := newJSONPathMap(t, map[string]any{"template": `{"v": "$MESSAGE.tags"}`})
	in := &flow.Message{
		Key:  "k",
		Type: flow.TypeMap,
		Data: map[string]any{"tags": []any{"a"}},
	}
	out, err := b.Process(in)
	if err == nil || !strings.Contains(err.Error(), "template did not produce a JSON object") {
		t.Fatalf("err = %v, want JSON object error", err)
	}
	if len(out) != 0 {
		t.Fatalf("outputs on error = %v, want none", out)
	}
}

func TestJSONPathMap_InputNotMutated_MetadataIndependent(t *testing.T) {
	b := newJSONPathMap(t, map[string]any{"template": `{"v": $MESSAGE.n}`})
	in := &flow.Message{
		Key:        "k",
		Timestamp:  42,
		Metadata:   map[string]string{"u": "alice"},
		Type:       flow.TypeMap,
		Data:       map[string]any{"n": int64(1)},
		FieldTypes: map[string]flow.DataType{"n": flow.TypeInteger},
	}
	out, err := b.Process(in)
	if err != nil || len(out) != 1 {
		t.Fatalf("Process: %v %v", out, err)
	}
	out[0].Metadata["u"] = "mutated"
	delete(out[0].Metadata, "u")
	if in.Metadata["u"] != "alice" || len(in.Metadata) != 1 {
		t.Errorf("input metadata changed via output: %v", in.Metadata)
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
	if in.Key != "k" || in.Timestamp != 42 {
		t.Errorf("input key/timestamp mutated: %+v", in)
	}
}

func TestJSONPathMap_NilPassthrough(t *testing.T) {
	b := newJSONPathMap(t, map[string]any{"template": `{"v": $MESSAGE}`})
	out, err := b.Process(nil)
	if err != nil || len(out) != 0 {
		t.Fatalf("nil: %v %v", out, err)
	}
}

func TestJSONPathMap_Registered(t *testing.T) {
	reg := blocks.DefaultRegistry()
	if !reg.Has(blocks.TypeJSONPathMap) {
		t.Error("DefaultRegistry missing json_path_map")
	}
	info, ok := blocks.LookupInfo(blocks.TypeJSONPathMap)
	if !ok {
		t.Fatal("LookupInfo missing json_path_map docs")
	}
	if info.Role != blocks.RoleTransform {
		t.Errorf("Role = %v, want transform", info.Role)
	}
}
