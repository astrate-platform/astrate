package blocks_test

import (
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
)

func TestFilter_RequiresCondition(t *testing.T) {
	_, err := blocks.Filter("f", map[string]any{}, flow.Deps{})
	if err == nil || !strings.Contains(err.Error(), "at least one") {
		t.Fatalf("err = %v", err)
	}
}

func TestFilter_KeyPrefixAndMetadata(t *testing.T) {
	b, err := blocks.Filter("f", map[string]any{
		"key_prefix": "acme/",
		"metadata":   map[string]any{"kind": "data"},
	}, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}

	pass := &flow.FlowMessage{
		Key:      "acme/dev1",
		Type:     flow.TypeString,
		Data:     "x",
		Metadata: map[string]string{"kind": "data", "path": "/t"},
	}
	out, err := b.Process(pass)
	if err != nil || len(out) != 1 || out[0] != pass {
		t.Fatalf("pass: out=%v err=%v", out, err)
	}

	// Wrong prefix.
	drop, err := b.Process(&flow.FlowMessage{
		Key: "other/dev", Type: flow.TypeString, Data: "x",
		Metadata: map[string]string{"kind": "data"},
	})
	if err != nil || len(drop) != 0 {
		t.Fatalf("drop prefix: out=%v err=%v", drop, err)
	}

	// Wrong metadata.
	drop2, err := b.Process(&flow.FlowMessage{
		Key: "acme/dev", Type: flow.TypeString, Data: "x",
		Metadata: map[string]string{"kind": "connection"},
	})
	if err != nil || len(drop2) != 0 {
		t.Fatalf("drop metadata: out=%v err=%v", drop2, err)
	}
}

func TestFilter_Type(t *testing.T) {
	b, err := blocks.Filter("f", map[string]any{"type": "integer"}, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := b.Process(&flow.FlowMessage{Key: "k", Type: flow.TypeInteger, Data: int64(1)})
	if err != nil || len(out) != 1 {
		t.Fatalf("integer: %v %v", out, err)
	}
	out, err = b.Process(&flow.FlowMessage{Key: "k", Type: flow.TypeString, Data: "1"})
	if err != nil || len(out) != 0 {
		t.Fatalf("string: %v %v", out, err)
	}
}

func TestFilter_UnknownType(t *testing.T) {
	_, err := blocks.Filter("f", map[string]any{"type": "blob"}, flow.Deps{})
	if err == nil || !strings.Contains(err.Error(), "unknown type") {
		t.Fatalf("err = %v", err)
	}
}

func TestMap_RequiresOp(t *testing.T) {
	_, err := blocks.Map("m", nil, flow.Deps{})
	if err == nil || !strings.Contains(err.Error(), "at least one") {
		t.Fatalf("err = %v", err)
	}
}

func TestMap_KeyTemplateAndMetadata(t *testing.T) {
	b, err := blocks.Map("m", map[string]any{
		"key":             "stream/{metadata.interface}/{key}",
		"set_metadata":    map[string]string{"stage": "mapped"},
		"delete_metadata": []any{"path"},
	}, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}

	in := &flow.FlowMessage{
		Key:  "acme/dev1",
		Type: flow.TypeString,
		Data: "v",
		Metadata: map[string]string{
			"interface": "org.ex.Sensor",
			"path":      "/temp",
			"kind":      "data",
		},
	}
	out, err := b.Process(in)
	if err != nil || len(out) != 1 {
		t.Fatalf("Process: %v %v", out, err)
	}
	got := out[0]
	if got.Key != "stream/org.ex.Sensor/acme/dev1" {
		t.Errorf("Key = %q", got.Key)
	}
	if got.Metadata["stage"] != "mapped" {
		t.Errorf("stage = %q", got.Metadata["stage"])
	}
	if _, ok := got.Metadata["path"]; ok {
		t.Errorf("path should be deleted: %v", got.Metadata)
	}
	if got.Metadata["kind"] != "data" {
		t.Errorf("kind should be preserved")
	}
	// Input must not be mutated.
	if in.Key != "acme/dev1" {
		t.Errorf("input key mutated: %q", in.Key)
	}
	if _, ok := in.Metadata["stage"]; ok {
		t.Errorf("input metadata mutated: %v", in.Metadata)
	}
	if in.Metadata["path"] != "/temp" {
		t.Errorf("input path should remain: %v", in.Metadata)
	}
	// Payload untouched.
	if got.Type != flow.TypeString || got.Data != "v" {
		t.Errorf("payload changed: type=%v data=%v", got.Type, got.Data)
	}
}

func TestMap_MissingMetadataPlaceholder(t *testing.T) {
	b, err := blocks.Map("m", map[string]any{"key": "x/{metadata.missing}/y"}, flow.Deps{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := b.Process(&flow.FlowMessage{Key: "k", Type: flow.TypeString, Data: ""})
	if err != nil || len(out) != 1 {
		t.Fatalf("%v %v", out, err)
	}
	if out[0].Key != "x//y" {
		t.Errorf("Key = %q", out[0].Key)
	}
}

func TestFilterMap_PipelineInstantiate(t *testing.T) {
	reg := blocks.DefaultRegistry()
	p := &flow.Pipeline{
		ID: "pipe",
		Blocks: []flow.PipelineNode{
			{Name: "src", BlockType: blocks.TypeAstarteSource, Config: map[string]any{"realm": "r"}},
			{Name: "f", BlockType: blocks.TypeFilter, Config: map[string]any{"key_prefix": "r/"}},
			{Name: "m", BlockType: blocks.TypeMap, Config: map[string]any{"set_metadata": map[string]any{"ok": "1"}}},
			{Name: "sink", BlockType: blocks.TypeNullSink},
		},
		Connections: []flow.Connection{
			{From: "src", To: "f"},
			{From: "f", To: "m"},
			{From: "m", To: "sink"},
		},
	}
	bus := stream.New(nil)
	blks, err := reg.Instantiate(p, flow.Deps{
		Bus:   bus,
		Realm: "r",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(blks) != 4 {
		t.Fatalf("len=%d", len(blks))
	}
	for i, want := range []string{"src", "f", "m", "sink"} {
		if blks[i].Name() != want {
			t.Errorf("block[%d]=%q", i, blks[i].Name())
		}
	}
	if s, ok := blks[0].(flow.Stopper); ok {
		s.Stop()
	}
}
