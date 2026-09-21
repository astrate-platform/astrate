package store

import (
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

func TestValidatePipelineGraph_InvalidBlockRef(t *testing.T) {
	def := []byte(`{
		"blocks": [{"name": "a"}, {"name": "b"}],
		"connections": [{"from": "a", "to": "ghost"}]
	}`)
	err := validatePipelineGraph(def)
	if err == nil || errors.Is(err, ErrPipelineCyclic) {
		t.Errorf("validatePipelineGraph(invalid ref) = %v, want a non-cyclic error", err)
	}
}

func TestValidatePipelineGraph_DuplicateBlockName(t *testing.T) {
	def := []byte(`{
		"blocks": [{"name": "a"}, {"name": "a"}],
		"connections": []
	}`)
	err := validatePipelineGraph(def)
	if err == nil || errors.Is(err, ErrPipelineCyclic) {
		t.Errorf("validatePipelineGraph(duplicate name) = %v, want a non-cyclic error", err)
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
