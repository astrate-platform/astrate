package store

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestValidatePipelineGraph_Acyclic(t *testing.T) {
	def := []byte(`{
		"blocks": [{"name": "a"}, {"name": "b"}, {"name": "c"}],
		"connections": [{"from": "a", "to": "b"}, {"from": "b", "to": "c"}]
	}`)
	if err := validatePipelineGraph(def); err != nil {
		t.Errorf("validatePipelineGraph(acyclic) = %v, want nil", err)
	}
}

func TestValidatePipelineGraph_Cyclic(t *testing.T) {
	def := []byte(`{
		"blocks": [{"name": "a"}, {"name": "b"}, {"name": "c"}],
		"connections": [{"from": "a", "to": "b"}, {"from": "b", "to": "c"}, {"from": "c", "to": "a"}]
	}`)
	err := validatePipelineGraph(def)
	if !errors.Is(err, ErrPipelineCyclic) {
		t.Errorf("validatePipelineGraph(cyclic) = %v, want ErrPipelineCyclic", err)
	}
}

func TestValidatePipelineGraph_MalformedJSON(t *testing.T) {
	def := []byte(`{"blocks": [{"name": "a"}`)
	err := validatePipelineGraph(def)
	if err == nil || errors.Is(err, ErrPipelineCyclic) {
		t.Fatalf("validatePipelineGraph(malformed) = %v, want a non-cyclic error", err)
	}
	if !strings.Contains(err.Error(), "does not parse") {
		t.Errorf("validatePipelineGraph(malformed) = %v, want the parse error", err)
	}
	var syntaxErr *json.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Errorf("validatePipelineGraph(malformed) = %v, want it to wrap *json.SyntaxError", err)
	}
	// The parse must fail before the graph is inspected, so the zero-blocks
	// check cannot be what rejects this definition.
	if strings.Contains(err.Error(), "no blocks") {
		t.Errorf("validatePipelineGraph(malformed) = %v, want the parse error to win over the no-blocks error", err)
	}
}

func TestValidatePipelineGraph_InvalidBlockRef(t *testing.T) {
	tests := []struct {
		name    string
		def     string
		wantMsg string
	}{
		{
			name:    "unknown destination",
			def:     `{"blocks": [{"name": "a"}, {"name": "b"}], "connections": [{"from": "a", "to": "ghost"}]}`,
			wantMsg: `unknown block "a" -> "ghost"`,
		},
		{
			name:    "unknown source",
			def:     `{"blocks": [{"name": "a"}, {"name": "b"}], "connections": [{"from": "ghost", "to": "b"}]}`,
			wantMsg: `unknown block "ghost" -> "b"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePipelineGraph([]byte(tt.def))
			if err == nil || errors.Is(err, ErrPipelineCyclic) {
				t.Fatalf("validatePipelineGraph(invalid ref) = %v, want a non-cyclic error", err)
			}
			if !strings.Contains(err.Error(), tt.wantMsg) {
				t.Errorf("validatePipelineGraph(invalid ref) = %v, want it to mention %q", err, tt.wantMsg)
			}
		})
	}
}

func TestValidatePipelineGraph_DuplicateBlockName(t *testing.T) {
	def := []byte(`{
		"blocks": [{"name": "a"}, {"name": "a"}],
		"connections": []
	}`)
	err := validatePipelineGraph(def)
	if err == nil || errors.Is(err, ErrPipelineCyclic) {
		t.Fatalf("validatePipelineGraph(duplicate name) = %v, want a non-cyclic error", err)
	}
	if !strings.Contains(err.Error(), `duplicate block name "a"`) {
		t.Errorf("validatePipelineGraph(duplicate name) = %v, want the duplicate-name error", err)
	}
}

func TestValidatePipelineGraph_NoBlocks(t *testing.T) {
	def := []byte(`{
		"blocks": [],
		"connections": []
	}`)
	err := validatePipelineGraph(def)
	if err == nil || errors.Is(err, ErrPipelineCyclic) {
		t.Errorf("validatePipelineGraph(no blocks) = %v, want a non-cyclic error", err)
	}
	if !strings.Contains(err.Error(), "pipeline has no blocks") {
		t.Errorf("validatePipelineGraph(no blocks) = %v, want the no-blocks error", err)
	}
}

func TestValidatePipelineGraph_EmptyBlockName(t *testing.T) {
	def := []byte(`{
		"blocks": [{"name": ""}],
		"connections": []
	}`)
	err := validatePipelineGraph(def)
	if err == nil || errors.Is(err, ErrPipelineCyclic) {
		t.Errorf("validatePipelineGraph(empty name) = %v, want a non-cyclic error", err)
	}
	if !strings.Contains(err.Error(), "block with empty name") {
		t.Errorf("validatePipelineGraph(empty name) = %v, want the empty-name error", err)
	}
}
