package flow_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/flow"
)

func TestSubstituteConfig_Replace(t *testing.T) {
	def := []byte(`{
		"blocks": [{
			"name": "http",
			"block_type": "http_sink",
			"config": {
				"url": "https://hooks.example/${config.tenant}/in",
				"token": "${config.token}",
				"retries": 3
			}
		}],
		"connections": []
	}`)
	out, err := flow.SubstituteConfig(def, map[string]any{
		"tenant": "acme",
		"token":  "secret",
	})
	if err != nil {
		t.Fatalf("SubstituteConfig: %v", err)
	}
	var root map[string]any
	if err := json.Unmarshal(out, &root); err != nil {
		t.Fatal(err)
	}
	blocks := root["blocks"].([]any)
	cfg := blocks[0].(map[string]any)["config"].(map[string]any)
	if cfg["url"] != "https://hooks.example/acme/in" {
		t.Errorf("url = %v", cfg["url"])
	}
	if cfg["token"] != "secret" {
		t.Errorf("token = %v", cfg["token"])
	}
	// Non-string leaves unchanged (JSON numbers stay numbers).
	if cfg["retries"].(float64) != 3 {
		t.Errorf("retries = %v", cfg["retries"])
	}
}

func TestSubstituteConfig_MissingKey(t *testing.T) {
	def := []byte(`{"blocks":[{"name":"b","config":{"x":"${config.missing}"}}]}`)
	_, err := flow.SubstituteConfig(def, map[string]any{"other": "v"})
	if err == nil || !strings.Contains(err.Error(), "missing config key") {
		t.Fatalf("err = %v, want missing key", err)
	}
}

func TestSubstituteConfig_NoPlaceholders(t *testing.T) {
	def := []byte(`{"blocks":[{"name":"b","config":{"x":"plain"}}],"connections":[]}`)
	out, err := flow.SubstituteConfig(def, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"x":"plain"`) {
		t.Errorf("out = %s", out)
	}
}

func TestSubstituteConfig_StringableNumberBool(t *testing.T) {
	def := []byte(`{"blocks":[{"config":{"n":"${config.n}","b":"${config.b}"}}]}`)
	out, err := flow.SubstituteConfig(def, map[string]any{"n": float64(42), "b": true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"n":"42"`) || !strings.Contains(string(out), `"b":"true"`) {
		t.Errorf("out = %s", out)
	}
}

func TestSubstituteConfig_ObjectValueRejected(t *testing.T) {
	def := []byte(`{"blocks":[{"config":{"x":"${config.obj}"}}]}`)
	_, err := flow.SubstituteConfig(def, map[string]any{"obj": map[string]any{"a": 1}})
	if err == nil || !strings.Contains(err.Error(), "stringable") {
		t.Fatalf("err = %v, want stringable reject", err)
	}
}
