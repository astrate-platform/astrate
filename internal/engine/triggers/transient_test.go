package triggers

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCompileConditionDeviceTrigger(t *testing.T) {
	raw := json.RawMessage(`{"type":"device_trigger","on":"device_connected"}`)
	tg, err := CompileCondition("live-viewer", raw)
	if err != nil {
		t.Fatalf("CompileCondition returned error: %v", err)
	}
	if tg == nil {
		t.Fatal("expected non-nil trigger")
	}
	if !tg.MatchesDevice(DeviceEvent{DeviceID: "abc123", On: OnDeviceConnected}) {
		t.Error("expected MatchesDevice to report true")
	}
}

func TestCompileConditionDataTrigger(t *testing.T) {
	raw := json.RawMessage(`{
		"type":"data_trigger",
		"on":"incoming_data",
		"interface_name":"org.example.Sensors",
		"interface_major":1,
		"match_path":"/value",
		"value_match_operator":"*"
	}`)
	tg, err := CompileCondition("data-view", raw)
	if err != nil {
		t.Fatalf("CompileCondition returned error: %v", err)
	}
	if tg == nil {
		t.Fatal("expected non-nil trigger")
	}
	if !tg.MatchesData(DataEvent{
		DeviceID:  "dev1",
		On:        OnIncomingData,
		Interface: "org.example.Sensors",
		Major:     1,
		Path:      "/value",
	}) {
		t.Error("expected MatchesData true for matching event")
	}
	if tg.MatchesData(DataEvent{
		DeviceID:  "dev1",
		On:        OnIncomingData,
		Interface: "org.other.Interface",
		Major:     1,
		Path:      "/value",
	}) {
		t.Error("expected MatchesData false for non-matching interface")
	}
}

func TestCompileConditionActionIsNil(t *testing.T) {
	raw := json.RawMessage(`{"type":"device_trigger","on":"device_connected"}`)
	tg, err := CompileCondition("nil-action", raw)
	if err != nil {
		t.Fatalf("CompileCondition returned error: %v", err)
	}
	if tg.Action != nil {
		t.Errorf("expected Action == nil, got %v", tg.Action)
	}
}

func TestCompileConditionWildcardData(t *testing.T) {
	raw := json.RawMessage(`{
		"type":"data_trigger",
		"on":"incoming_data",
		"interface_name":"*",
		"match_path":"/*",
		"value_match_operator":"*"
	}`)
	tg, err := CompileCondition("watch-all", raw)
	if err != nil {
		t.Fatalf("CompileCondition returned error: %v", err)
	}
	if tg == nil {
		t.Fatal("expected non-nil trigger")
	}
	if !tg.MatchesData(DataEvent{
		DeviceID:  "dev1",
		On:        OnIncomingData,
		Interface: "org.any.Thing",
		Major:     1,
		Path:      "/some/path",
	}) {
		t.Error("expected wildcard to match any DataEvent")
	}
}

func TestCompileConditionInvalidInputs(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantMsg string
	}{
		{
			name:    "unknown type",
			raw:     `{"type":"bogus_trigger","on":"incoming_data"}`,
			wantMsg: `unknown trigger type "bogus_trigger"`,
		},
		{
			name:    "data_trigger missing interface_name",
			raw:     `{"type":"data_trigger","on":"incoming_data","match_path":"/p","value_match_operator":"*"}`,
			wantMsg: "data_trigger requires interface_name",
		},
		{
			name:    "data_trigger missing interface_major",
			raw:     `{"type":"data_trigger","on":"incoming_data","interface_name":"org.example.Sensors","match_path":"/p","value_match_operator":"*"}`,
			wantMsg: "data_trigger requires interface_major",
		},
		{
			name:    "unknown device on",
			raw:     `{"type":"device_trigger","on":"bogus_event"}`,
			wantMsg: `unknown device_trigger condition "bogus_event"`,
		},
		{
			name:    "not json",
			raw:     `not json at all`,
			wantMsg: "does not parse",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CompileCondition("bad-trigger", []byte(tt.raw))
			if err == nil {
				t.Fatal("expected non-nil error")
			}
			if !strings.Contains(err.Error(), tt.wantMsg) {
				t.Errorf("error %q does not contain %q", err.Error(), tt.wantMsg)
			}
			if !strings.Contains(err.Error(), "bad-trigger") {
				t.Errorf("error %q does not mention trigger name", err.Error())
			}
		})
	}
}
