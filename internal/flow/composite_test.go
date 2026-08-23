package flow_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/flow"
)

type compNode struct {
	Name      string         `json:"name"`
	BlockType string         `json:"block_type"`
	Config    map[string]any `json:"config,omitempty"`
}

type compConn struct {
	From     string `json:"from"`
	FromPort string `json:"from_port,omitempty"`
	To       string `json:"to"`
	ToPort   string `json:"to_port,omitempty"`
}

type compGraph struct {
	Blocks      []compNode `json:"blocks"`
	Connections []compConn `json:"connections"`
}

func parseGraph(t *testing.T, def []byte) compGraph {
	t.Helper()
	var g compGraph
	if err := json.Unmarshal(def, &g); err != nil {
		t.Fatalf("expanded definition does not parse: %v", err)
	}
	return g
}

func resolverOf(store map[string]*flow.UserBlock) func(string) (*flow.UserBlock, error) {
	return func(name string) (*flow.UserBlock, error) {
		ub, ok := store[name]
		if !ok {
			return nil, nil
		}
		return ub, nil
	}
}

func nodeNames(g compGraph) []string {
	out := make([]string, 0, len(g.Blocks))
	for _, n := range g.Blocks {
		out = append(out, n.Name)
	}
	return out
}

func edgeKeys(g compGraph) []string {
	out := make([]string, 0, len(g.Connections))
	for _, c := range g.Connections {
		k := c.From
		if c.FromPort != "" {
			k += "[" + c.FromPort + "]"
		}
		k += "->" + c.To
		if c.ToPort != "" {
			k += "[" + c.ToPort + "]"
		}
		out = append(out, k)
	}
	return out
}

func assertSameStrings(t *testing.T, rule string, got, want []string) {
	t.Helper()
	gotSorted := append([]string(nil), got...)
	sort.Strings(gotSorted)
	wantSorted := append([]string(nil), want...)
	sort.Strings(wantSorted)
	if !reflect.DeepEqual(gotSorted, wantSorted) {
		t.Errorf("%s: got %v, want %v", rule, gotSorted, wantSorted)
	}
}

func TestExpandComposites_NoCompositePassthrough(t *testing.T) {
	def := []byte(`{
		"blocks": [
			{"name": "raw", "block_type": "raw"},
			{"name": "sink", "block_type": "log_sink", "config": {"level": "info", "retries": 3}}
		],
		"connections": [
			{"from": "raw", "from_port": "p", "to": "sink", "to_port": "in"}
		]
	}`)
	out, err := flow.ExpandComposites(def, resolverOf(nil))
	if err != nil {
		t.Fatalf("passthrough must not error: %v", err)
	}
	var inG, outG compGraph
	if err := json.Unmarshal(def, &inG); err != nil {
		t.Fatal(err)
	}
	outG = parseGraph(t, out)
	if !reflect.DeepEqual(inG, outG) {
		t.Errorf("no-composite passthrough must keep parsed graph identical: got %+v, want %+v", outG, inG)
	}
}

func TestExpandComposites_SimpleSplice(t *testing.T) {
	def := []byte(`{
		"blocks": [
			{"name": "raw", "block_type": "raw"},
			{"name": "comp", "block_type": "my_comp", "config": {"k": "v"}},
			{"name": "log_sink", "block_type": "log_sink"}
		],
		"connections": [
			{"from": "raw", "to": "comp"},
			{"from": "comp", "to": "log_sink"}
		]
	}`)
	store := map[string]*flow.UserBlock{
		"my_comp": {
			Name:      "my_comp",
			BlockType: "producer_consumer",
			Source:    []byte(`{"blocks":[{"name":"filter","block_type":"filter"},{"name":"map","block_type":"map"}],"connections":[{"from":"filter","to":"map"}]}`),
		},
	}
	out, err := flow.ExpandComposites(def, resolverOf(store))
	if err != nil {
		t.Fatalf("simple splice must not error: %v", err)
	}
	g := parseGraph(t, out)
	assertSameStrings(t,
		"composite must be replaced by its prefixed sub-chain nodes",
		nodeNames(g),
		[]string{"raw", "comp.filter", "comp.map", "log_sink"})
	assertSameStrings(t,
		"outer edges must re-attach to first/last inner nodes and the inner edge must survive",
		edgeKeys(g),
		[]string{"raw->comp.filter", "comp.filter->comp.map", "comp.map->log_sink"})
}

func TestExpandComposites_BoundaryWiring(t *testing.T) {
	store := map[string]*flow.UserBlock{
		"my_comp": {
			Name:      "my_comp",
			BlockType: "producer_consumer",
			Source:    []byte(`{"blocks":[{"name":"filter","block_type":"filter"},{"name":"map","block_type":"map"}],"connections":[{"from":"filter","to":"map"}]}`),
		},
	}

	t.Run("composite last only incoming", func(t *testing.T) {
		def := []byte(`{
			"blocks": [
				{"name": "raw", "block_type": "raw"},
				{"name": "comp", "block_type": "my_comp"}
			],
			"connections": [{"from": "raw", "to": "comp"}]
		}`)
		out, err := flow.ExpandComposites(def, resolverOf(store))
		if err != nil {
			t.Fatalf("expand: %v", err)
		}
		g := parseGraph(t, out)
		assertSameStrings(t,
			"incoming outer edge must attach to the sub-chain's first node when the composite is last",
			edgeKeys(g),
			[]string{"raw->comp.filter", "comp.filter->comp.map"})
	})

	t.Run("composite first only outgoing", func(t *testing.T) {
		def := []byte(`{
			"blocks": [
				{"name": "comp", "block_type": "my_comp"},
				{"name": "log_sink", "block_type": "log_sink"}
			],
			"connections": [{"from": "comp", "to": "log_sink"}]
		}`)
		out, err := flow.ExpandComposites(def, resolverOf(store))
		if err != nil {
			t.Fatalf("expand: %v", err)
		}
		g := parseGraph(t, out)
		assertSameStrings(t,
			"outgoing outer edge must leave from the sub-chain's last node when the composite is first",
			edgeKeys(g),
			[]string{"comp.filter->comp.map", "comp.map->log_sink"})
	})
}

func TestExpandComposites_Params(t *testing.T) {
	source := []byte(`{"blocks":[{"name":"filter","block_type":"filter","config":{"threshold":"${config.threshold}"}}],"connections":[]}`)

	t.Run("supplied key is substituted into expanded config", func(t *testing.T) {
		def := []byte(`{
			"blocks": [{"name": "comp", "block_type": "thresh", "config": {"threshold": "0.7"}}],
			"connections": []
		}`)
		store := map[string]*flow.UserBlock{
			"thresh": {Name: "thresh", BlockType: "consumer", Source: source},
		}
		out, err := flow.ExpandComposites(def, resolverOf(store))
		if err != nil {
			t.Fatalf("expand with params must not error: %v", err)
		}
		g := parseGraph(t, out)
		for _, n := range g.Blocks {
			if n.Name == "comp.filter" && n.Config["threshold"] != "0.7" {
				t.Errorf("substituted value must survive into the expanded node's config: threshold = %v", n.Config["threshold"])
			}
		}
	})

	t.Run("missing key fails loudly", func(t *testing.T) {
		def := []byte(`{
			"blocks": [{"name": "comp", "block_type": "thresh", "config": {}}],
			"connections": []
		}`)
		store := map[string]*flow.UserBlock{
			"thresh": {Name: "thresh", BlockType: "consumer", Source: source},
		}
		_, err := flow.ExpandComposites(def, resolverOf(store))
		if err == nil || !strings.Contains(err.Error(), "missing config key") {
			t.Errorf("missing param must fail loudly, got %v", err)
		}
	})
}

func TestExpandComposites_NestedTwoLevels(t *testing.T) {
	def := []byte(`{
		"blocks": [
			{"name": "raw", "block_type": "raw"},
			{"name": "a", "block_type": "comp_a"},
			{"name": "log_sink", "block_type": "log_sink"}
		],
		"connections": [
			{"from": "raw", "to": "a"},
			{"from": "a", "to": "log_sink"}
		]
	}`)
	store := map[string]*flow.UserBlock{
		"comp_a": {
			Name:      "comp_a",
			BlockType: "producer_consumer",
			Source:    []byte(`{"blocks":[{"name":"b","block_type":"comp_b"}],"connections":[]}`),
		},
		"comp_b": {
			Name:      "comp_b",
			BlockType: "producer",
			Source:    []byte(`{"blocks":[{"name":"block","block_type":"virtual_device_pool"}],"connections":[]}`),
		},
	}
	out, err := flow.ExpandComposites(def, resolverOf(store))
	if err != nil {
		t.Fatalf("nested expansion must not error: %v", err)
	}
	g := parseGraph(t, out)
	assertSameStrings(t,
		"nested composite names must nest as outer.inner.block",
		nodeNames(g),
		[]string{"raw", "a.b.block", "log_sink"})
	assertSameStrings(t,
		"edges must chain through both nesting levels",
		edgeKeys(g),
		[]string{"raw->a.b.block", "a.b.block->log_sink"})
}

func TestExpandComposites_CycleRejected(t *testing.T) {
	t.Run("mutual reference A<->B", func(t *testing.T) {
		def := []byte(`{"blocks":[{"name":"x","block_type":"comp_a"}],"connections":[]}`)
		store := map[string]*flow.UserBlock{
			"comp_a": {Name: "comp_a", Source: []byte(`{"blocks":[{"name":"n","block_type":"comp_b"}],"connections":[]}`)},
			"comp_b": {Name: "comp_b", Source: []byte(`{"blocks":[{"name":"m","block_type":"comp_a"}],"connections":[]}`)},
		}
		_, err := flow.ExpandComposites(def, resolverOf(store))
		if err == nil || !strings.Contains(err.Error(), "cycle") {
			t.Errorf("composite cycle must be reported as cycle error, got %v", err)
		}
	})

	t.Run("self-reference A->A", func(t *testing.T) {
		def := []byte(`{"blocks":[{"name":"x","block_type":"comp_s"}],"connections":[]}`)
		store := map[string]*flow.UserBlock{
			"comp_s": {Name: "comp_s", Source: []byte(`{"blocks":[{"name":"n","block_type":"comp_s"}],"connections":[]}`)},
		}
		_, err := flow.ExpandComposites(def, resolverOf(store))
		if err == nil || !strings.Contains(err.Error(), "cycle") {
			t.Errorf("self-referencing composite must be reported as cycle error, got %v", err)
		}
	})
}

func TestExpandComposites_NestingDepthCapped(t *testing.T) {
	store := map[string]*flow.UserBlock{}
	for i := 0; i < 40; i++ {
		next := "deep_" + strconv.Itoa(i+1)
		store["deep_"+strconv.Itoa(i)] = &flow.UserBlock{
			Name:   "deep_" + strconv.Itoa(i),
			Source: []byte(`{"blocks":[{"name":"n","block_type":"` + next + `"}],"connections":[]}`),
		}
	}
	def := []byte(`{"blocks":[{"name":"x","block_type":"deep_0"}],"connections":[]}`)
	_, err := flow.ExpandComposites(def, resolverOf(store))
	if err == nil || !strings.Contains(err.Error(), "nesting too deep") {
		t.Errorf("nesting beyond the depth cap must be rejected, got %v", err)
	}
}

func TestExpandComposites_ResolverErrorPropagates(t *testing.T) {
	sentinel := errors.New("store unavailable")
	def := []byte(`{"blocks":[{"name":"x","block_type":"whatever"}],"connections":[]}`)
	_, err := flow.ExpandComposites(def, func(string) (*flow.UserBlock, error) {
		return nil, sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Errorf("resolver store errors must propagate via errors.Is, got %v", err)
	}
}

func TestExpandComposites_InnerStructuralRejection(t *testing.T) {
	t.Run("connection to nonexistent block", func(t *testing.T) {
		def := []byte(`{"blocks":[{"name":"comp","block_type":"broken"}],"connections":[]}`)
		store := map[string]*flow.UserBlock{
			"broken": {
				Name:   "broken",
				Source: []byte(`{"blocks":[{"name":"f","block_type":"filter"}],"connections":[{"from":"f","to":"ghost"}]}`),
			},
		}
		_, err := flow.ExpandComposites(def, resolverOf(store))
		if err == nil || !strings.Contains(err.Error(), "unknown block") {
			t.Errorf("inner connection to a nonexistent block must be rejected, got %v", err)
		}
	})

	t.Run("empty blocks", func(t *testing.T) {
		def := []byte(`{"blocks":[{"name":"comp","block_type":"hollow"}],"connections":[]}`)
		store := map[string]*flow.UserBlock{
			"hollow": {Name: "hollow", Source: []byte(`{"blocks":[],"connections":[]}`)},
		}
		_, err := flow.ExpandComposites(def, resolverOf(store))
		if err == nil || !strings.Contains(err.Error(), "empty definition") {
			t.Errorf("composite with empty definition must be rejected, got %v", err)
		}
	})

	t.Run("duplicate inner names", func(t *testing.T) {
		def := []byte(`{"blocks":[{"name":"comp","block_type":"twins"}],"connections":[]}`)
		store := map[string]*flow.UserBlock{
			"twins": {
				Name:   "twins",
				Source: []byte(`{"blocks":[{"name":"f","block_type":"filter"},{"name":"f","block_type":"map"}],"connections":[]}`),
			},
		}
		_, err := flow.ExpandComposites(def, resolverOf(store))
		if err == nil || !strings.Contains(err.Error(), "duplicate") {
			t.Errorf("duplicate inner block names must be rejected, got %v", err)
		}
	})
}

func TestExpandComposites_PortsPreservedOnRetarget(t *testing.T) {
	def := []byte(`{
		"blocks": [
			{"name": "raw", "block_type": "raw"},
			{"name": "comp", "block_type": "my_comp"},
			{"name": "log_sink", "block_type": "log_sink"}
		],
		"connections": [
			{"from": "raw", "from_port": "p", "to": "comp"},
			{"from": "comp", "to": "log_sink", "to_port": "q"}
		]
	}`)
	store := map[string]*flow.UserBlock{
		"my_comp": {
			Name:      "my_comp",
			BlockType: "producer_consumer",
			Source:    []byte(`{"blocks":[{"name":"filter","block_type":"filter"},{"name":"map","block_type":"map"}],"connections":[{"from":"filter","to":"map"}]}`),
		},
	}
	out, err := flow.ExpandComposites(def, resolverOf(store))
	if err != nil {
		t.Fatalf("expand: %v", err)
	}
	g := parseGraph(t, out)
	assertSameStrings(t,
		"ports on retargeted outer edges must be preserved",
		edgeKeys(g),
		[]string{"raw[p]->comp.filter", "comp.filter->comp.map", "comp.map->log_sink[q]"})
}
