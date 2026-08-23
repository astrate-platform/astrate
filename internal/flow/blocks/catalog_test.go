package blocks_test

import (
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
)

func TestDefaultRegistry_Types(t *testing.T) {
	reg := blocks.DefaultRegistry()
	for _, want := range []string{
		blocks.TypeAstarteSource,
		blocks.TypeFilter,
		blocks.TypeMap,
		blocks.TypeToJSON,
		blocks.TypeUpdateMetadata,
		blocks.TypeSplitMap,
		blocks.TypeRandomSource,
		blocks.TypeSort,
		blocks.TypeJSONPathMap,
		blocks.TypeContainer,
		blocks.TypeLogSink,
		blocks.TypeNullSink,
	} {
		if !reg.Has(want) {
			t.Errorf("missing %q", want)
		}
		if _, ok := blocks.LookupInfo(want); !ok {
			t.Errorf("LookupInfo missing docs for %q", want)
		}
	}
	infos := blocks.InfoForTypes(reg.Types())
	if len(infos) != len(reg.Types()) {
		t.Fatalf("InfoForTypes len %d, types %d", len(infos), len(reg.Types()))
	}
}

func TestAstarteSource_RequiresBus(t *testing.T) {
	_, err := blocks.AstarteSource("src", map[string]any{"realm": "r"}, flow.Deps{})
	if err == nil || !strings.Contains(err.Error(), "bus") {
		t.Fatalf("err = %v, want bus required", err)
	}
}

func TestAstarteSource_RequiresRealm(t *testing.T) {
	bus := stream.New(nil)
	_, err := blocks.AstarteSource("src", nil, flow.Deps{Bus: bus})
	if err == nil || !strings.Contains(err.Error(), "realm") {
		t.Fatalf("err = %v, want realm required", err)
	}
}

func TestAstarteSource_UsesDepsRealm(t *testing.T) {
	bus := stream.New(nil)
	b, err := blocks.AstarteSource("src", nil, flow.Deps{Bus: bus, Realm: "acme"})
	if err != nil {
		t.Fatalf("AstarteSource: %v", err)
	}
	if b.Name() != "src" {
		t.Errorf("Name = %q", b.Name())
	}
	if s, ok := b.(flow.Stopper); ok {
		s.Stop()
	} else {
		t.Fatal("expected Stopper")
	}
}

func TestNullAndLogSink(t *testing.T) {
	n, err := blocks.NullSink("n", nil, flow.Deps{})
	if err != nil || n.Name() != "n" {
		t.Fatalf("NullSink: %v name=%q", err, nameOr(n))
	}
	l, err := blocks.LogSink("l", nil, flow.Deps{})
	if err != nil || l.Name() != "l" {
		t.Fatalf("LogSink: %v name=%q", err, nameOr(l))
	}
	if _, err := n.Process(&flow.Message{Key: "k", Type: flow.TypeString, Data: "x"}); err != nil {
		t.Fatal(err)
	}
	if _, err := l.Process(&flow.Message{Key: "k", Type: flow.TypeString, Data: "x"}); err != nil {
		t.Fatal(err)
	}
}

func nameOr(b flow.Block) string {
	if b == nil {
		return "<nil>"
	}
	return b.Name()
}
