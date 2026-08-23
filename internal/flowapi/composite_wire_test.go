package flowapi

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/flow"
)

func TestValidateCompositeConfigs(t *testing.T) {
	schema := json.RawMessage(`{
		"type": "object",
		"properties": {"threshold": {"type": "integer"}},
		"required": ["threshold"],
		"additionalProperties": false
	}`)
	stored := []*flow.UserBlock{
		{Name: "windowed_counter", BlockType: "producer_consumer", ConfigSchema: schema},
		{Name: "passthrough", BlockType: "consumer"},
	}

	t.Run("valid config passes", func(t *testing.T) {
		nodes := []configNode{
			{Name: "w", BlockType: "windowed_counter", Config: map[string]any{"threshold": float64(3)}},
			{Name: "p", BlockType: "passthrough", Config: map[string]any{"anything": true}},
		}
		if err := validateCompositeConfigs(nodes, stored); err != nil {
			t.Fatalf("valid configs rejected: %v", err)
		}
	})

	t.Run("wrong-typed param rejected naming the node", func(t *testing.T) {
		nodes := []configNode{
			{Name: "w", BlockType: "windowed_counter", Config: map[string]any{"threshold": "three"}},
		}
		err := validateCompositeConfigs(nodes, stored)
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("expected ErrValidation, got %v", err)
		}
		for _, want := range []string{`"w"`, `"windowed_counter"`} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not name %s", err, want)
			}
		}
	})

	t.Run("missing required key rejected", func(t *testing.T) {
		nodes := []configNode{
			{Name: "w", BlockType: "windowed_counter"},
		}
		if err := validateCompositeConfigs(nodes, stored); !errors.Is(err, ErrValidation) {
			t.Fatalf("expected ErrValidation for missing threshold, got %v", err)
		}
	})

	t.Run("block with empty schema skipped", func(t *testing.T) {
		nodes := []configNode{
			{Name: "p", BlockType: "passthrough", Config: map[string]any{"junk": true}},
		}
		if err := validateCompositeConfigs(nodes, stored); err != nil {
			t.Fatalf("schema-less composite should be skipped: %v", err)
		}
	})

	t.Run("unknown node type skipped", func(t *testing.T) {
		nodes := []configNode{
			{Name: "src", BlockType: "astarte_source", Config: map[string]any{"whatever": 1}},
		}
		if err := validateCompositeConfigs(nodes, stored); err != nil {
			t.Fatalf("built-in-typed node should be skipped: %v", err)
		}
	})

	t.Run("broken schema rejected naming the composite", func(t *testing.T) {
		broken := []*flow.UserBlock{
			{Name: "bad_comp", BlockType: "consumer", ConfigSchema: json.RawMessage(`{"type":"nope"}`)},
		}
		err := validateCompositeConfigs(nil, broken)
		if !errors.Is(err, ErrValidation) || !strings.Contains(err.Error(), `"bad_comp"`) {
			t.Fatalf("expected ErrValidation naming bad_comp, got %v", err)
		}
	})
}

func TestUnknownBlockTypeErr(t *testing.T) {
	builtins := []string{"astarte_source", "null_sink"}

	cases := []struct {
		name    string
		stored  []*flow.UserBlock
		wantSub string
	}{
		{
			name:    "no stored composites keeps today's message",
			stored:  nil,
			wantSub: `unknown block_type "nope" on block "src" (known: [astarte_source null_sink])`,
		},
		{
			name: "stored composites listed after known built-ins",
			stored: []*flow.UserBlock{
				{Name: "alpha"}, {Name: "beta"},
			},
			wantSub: `(known: [astarte_source null_sink]) (stored composites: [alpha beta])`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := unknownBlockTypeErr("nope", "src", builtins, tc.stored)
			if !errors.Is(err, ErrValidation) {
				t.Fatalf("expected ErrValidation, got %v", err)
			}
			if !strings.Contains(err.Error(), tc.wantSub) {
				t.Errorf("error %q missing %q", err, tc.wantSub)
			}
		})
	}
}

func TestDecodePipelineNodes(t *testing.T) {
	def := []byte(`{"blocks":[{"name":"w","block_type":"comp","config":{"k":"${config.v}"}},{"name":"s","block_type":"log_sink"}],"connections":[]}`)
	nodes, err := decodePipelineNodes(def)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 2 || nodes[0].Name != "w" || nodes[0].Config["k"] != "${config.v}" {
		t.Fatalf("unexpected nodes: %+v", nodes)
	}
	if _, err := decodePipelineNodes([]byte(`not json`)); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation for unparseable definition, got %v", err)
	}
}
