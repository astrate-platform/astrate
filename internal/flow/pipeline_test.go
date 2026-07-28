package flow

import (
	"encoding/json"
	"testing"
)

func TestPipeline_Validate(t *testing.T) {
	tests := []struct {
		name    string
		p       Pipeline
		wantErr bool
	}{
		{
			name: "valid linear pipeline",
			p: Pipeline{
				ID:   "pipe-1",
				Name: "simple",
				Blocks: []PipelineNode{
					{Name: "src", BlockType: "source"},
					{Name: "sink", BlockType: "sink"},
				},
				Connections: []Connection{
					{From: "src", To: "sink"},
				},
			},
		},
		{
			name: "valid diamond pipeline",
			p: Pipeline{
				ID:   "pipe-diamond",
				Name: "diamond",
				Blocks: []PipelineNode{
					{Name: "src", BlockType: "source"},
					{Name: "a", BlockType: "transform"},
					{Name: "b", BlockType: "transform"},
					{Name: "sink", BlockType: "sink"},
				},
				Connections: []Connection{
					{From: "src", To: "a"},
					{From: "src", To: "b"},
					{From: "a", To: "sink"},
					{From: "b", To: "sink"},
				},
			},
		},
		{
			name: "valid with ports",
			p: Pipeline{
				ID:   "pipe-ports",
				Name: "ports",
				Blocks: []PipelineNode{
					{Name: "src", BlockType: "source"},
					{Name: "proc", BlockType: "transform"},
					{Name: "out", BlockType: "sink"},
				},
				Connections: []Connection{
					{From: "src", FromPort: "main", To: "proc", ToPort: "input"},
					{From: "proc", FromPort: "output", To: "out", ToPort: "main"},
				},
			},
		},
		{
			name: "empty ID",
			p: Pipeline{
				Blocks: []PipelineNode{{Name: "a", BlockType: "x"}},
			},
			wantErr: true,
		},
		{
			name: "no blocks",
			p: Pipeline{
				ID: "empty",
			},
			wantErr: true,
		},
		{
			name: "duplicate block names",
			p: Pipeline{
				ID: "dup",
				Blocks: []PipelineNode{
					{Name: "a", BlockType: "x"},
					{Name: "a", BlockType: "y"},
				},
				Connections: []Connection{
					{From: "a", To: "a"},
				},
			},
			wantErr: true,
		},
		{
			name: "connection references missing block",
			p: Pipeline{
				ID: "bad-ref",
				Blocks: []PipelineNode{
					{Name: "src", BlockType: "source"},
				},
				Connections: []Connection{
					{From: "src", To: "missing"},
				},
			},
			wantErr: true,
		},
		{
			name: "self-loop",
			p: Pipeline{
				ID: "self-loop",
				Blocks: []PipelineNode{
					{Name: "a", BlockType: "x"},
				},
				Connections: []Connection{
					{From: "a", To: "a"},
				},
			},
			wantErr: true,
		},
		{
			name: "cycle",
			p: Pipeline{
				ID: "cycle",
				Blocks: []PipelineNode{
					{Name: "a", BlockType: "x"},
					{Name: "b", BlockType: "x"},
					{Name: "c", BlockType: "x"},
				},
				Connections: []Connection{
					{From: "a", To: "b"},
					{From: "b", To: "c"},
					{From: "c", To: "a"},
				},
			},
			wantErr: true,
		},
		{
			name: "no source (all blocks have incoming edges)",
			p: Pipeline{
				ID: "no-source",
				Blocks: []PipelineNode{
					{Name: "a", BlockType: "x"},
					{Name: "b", BlockType: "x"},
				},
				Connections: []Connection{
					{From: "a", To: "b"},
					{From: "b", To: "a"},
				},
			},
			wantErr: true,
		},
		{
			name: "no sink (all blocks have outgoing edges)",
			p: Pipeline{
				ID: "no-sink",
				Blocks: []PipelineNode{
					{Name: "a", BlockType: "x"},
					{Name: "b", BlockType: "x"},
				},
				Connections: []Connection{
					{From: "a", To: "b"},
					{From: "b", To: "a"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.p.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPipeline_JSONRoundTrip(t *testing.T) {
	p := Pipeline{
		ID:   "rt-pipe",
		Name: "round-trip",
		Blocks: []PipelineNode{
			{
				Name:      "src",
				BlockType: "generic_source",
				Config:    map[string]any{"url": "https://example.com"},
			},
			{
				Name:      "transform",
				BlockType: "lua",
				Config:    map[string]any{"script": "return msg"},
			},
			{
				Name:      "sink",
				BlockType: "json_out",
				Config:    map[string]any{"pretty": true},
			},
		},
		Connections: []Connection{
			{From: "src", FromPort: "out", To: "transform", ToPort: "in"},
			{From: "transform", To: "sink"},
		},
	}

	b, err := json.Marshal(&p)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got Pipeline
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.ID != p.ID {
		t.Errorf("ID = %q, want %q", got.ID, p.ID)
	}
	if got.Name != p.Name {
		t.Errorf("Name = %q, want %q", got.Name, p.Name)
	}
	if len(got.Blocks) != len(p.Blocks) {
		t.Fatalf("Blocks len = %d, want %d", len(got.Blocks), len(p.Blocks))
	}
	for i, b := range p.Blocks {
		gb := got.Blocks[i]
		if gb.Name != b.Name || gb.BlockType != b.BlockType {
			t.Errorf("Block[%d] = %+v, want %+v", i, gb, b)
		}
	}
	if len(got.Connections) != len(p.Connections) {
		t.Fatalf("Connections len = %d, want %d", len(got.Connections), len(p.Connections))
	}
	for i, c := range p.Connections {
		gc := got.Connections[i]
		if gc.From != c.From || gc.To != c.To || gc.FromPort != c.FromPort || gc.ToPort != c.ToPort {
			t.Errorf("Connection[%d] = %+v, want %+v", i, gc, c)
		}
	}
}

func TestPipeline_MarshalRejectsInvalidPipeline(t *testing.T) {
	p := Pipeline{
		ID: "bad",
		Blocks: []PipelineNode{
			{Name: "a", BlockType: "x"},
		},
		Connections: []Connection{
			{From: "a", To: "missing"},
		},
	}
	_, err := json.Marshal(&p)
	if err == nil {
		t.Fatal("expected marshal error for invalid pipeline")
	}
}

func TestPipeline_JSONContainsAllFields(t *testing.T) {
	p := Pipeline{
		ID:   "field-check",
		Name: "fields",
		Blocks: []PipelineNode{
			{Name: "src", BlockType: "source", Config: map[string]any{"k": "v"}},
			{Name: "out", BlockType: "sink"},
		},
		Connections: []Connection{
			{From: "src", FromPort: "out", To: "out", ToPort: "in"},
		},
	}
	b, err := json.Marshal(&p)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"id", "name", "blocks", "connections"} {
		if _, ok := raw[key]; !ok {
			t.Errorf("missing key %q in JSON", key)
		}
	}
}
