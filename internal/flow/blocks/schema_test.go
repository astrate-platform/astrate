package blocks

import (
	"encoding/json"
	"sort"
	"testing"
)

func TestBuiltinSchemas_Complete(t *testing.T) {
	for typ := range builtinInfo {
		if _, ok := builtinSchemas[typ]; !ok {
			t.Errorf("builtinSchemas missing entry for %q (builtinInfo has docs)", typ)
		}
	}
}

func TestBuiltinSchemas_WellFormed(t *testing.T) {
	type property struct {
		Type        string `json:"type"`
		Description string `json:"description"`
	}
	var schema struct {
		Type       string              `json:"type"`
		Properties map[string]property `json:"properties"`
		Required   []string            `json:"required"`
	}
	keys := make([]string, 0, len(builtinSchemas))
	for k := range builtinSchemas {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, typ := range keys {
		if err := json.Unmarshal([]byte(builtinSchemas[typ]), &schema); err != nil {
			t.Errorf("%s: unmarshal: %v", typ, err)
			continue
		}
		if schema.Type != "object" {
			t.Errorf("%s: type = %q, want \"object\"", typ, schema.Type)
		}
		if schema.Properties == nil {
			t.Errorf("%s: properties missing", typ)
			continue
		}
		for _, req := range schema.Required {
			if _, ok := schema.Properties[req]; !ok {
				t.Errorf("%s: required key %q not in properties", typ, req)
			}
		}
	}
}

func TestBuiltinSchemas_SpotRows(t *testing.T) {
	var container struct {
		Properties map[string]struct {
			Type        string `json:"type"`
			Description string `json:"description"`
		} `json:"properties"`
		Required []string `json:"required"`
	}
	if err := json.Unmarshal([]byte(builtinSchemas[TypeContainer]), &container); err != nil {
		t.Fatalf("container: %v", err)
	}
	if len(container.Required) != 1 || container.Required[0] != "image" {
		t.Errorf("container required = %v, want [image]", container.Required)
	}

	var randomSource struct {
		Properties map[string]struct {
			Type string `json:"type"`
		} `json:"properties"`
	}
	if err := json.Unmarshal([]byte(builtinSchemas[TypeRandomSource]), &randomSource); err != nil {
		t.Fatalf("random_source: %v", err)
	}
	if got := randomSource.Properties["type"].Type; got != "string" {
		t.Errorf("random_source type prop = %q, want string", got)
	}

	var httpSource struct {
		Properties map[string]struct {
			Type string `json:"type"`
		} `json:"properties"`
	}
	if err := json.Unmarshal([]byte(builtinSchemas[TypeHTTPSource]), &httpSource); err != nil {
		t.Fatalf("http_source: %v", err)
	}
	if got := httpSource.Properties["urls"].Type; got != "array" {
		t.Errorf("http_source urls prop = %q, want array", got)
	}

	var nullSink struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if err := json.Unmarshal([]byte(builtinSchemas[TypeNullSink]), &nullSink); err != nil {
		t.Fatalf("null_sink: %v", err)
	}
	if len(nullSink.Properties) != 0 {
		t.Errorf("null_sink properties = %d entries, want 0", len(nullSink.Properties))
	}
}

func TestLookupInfo_ConfigSchema(t *testing.T) {
	info, ok := LookupInfo(TypeLogSink)
	if !ok {
		t.Fatal("LookupInfo(log_sink) not found")
	}
	if info.ConfigSchema == nil {
		t.Error("LookupInfo(log_sink).ConfigSchema is nil, want merged schema")
	}
}
