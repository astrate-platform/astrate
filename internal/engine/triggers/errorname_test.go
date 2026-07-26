package triggers

import (
	"testing"
)

func TestUpstreamErrorNameMapping(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"realm_unknown", "device_session_not_found"},
		{"device_unknown", "device_session_not_found"},
		{"malformed_topic", "invalid_path"},
		{"interface_not_in_introspection", "interface_loading_failed"},
		{"interface_not_installed", "interface_loading_failed"},
		{"ownership_violation", "write_on_server_owned_interface"},
		{"unexpected_path", "mapping_not_found"},
		{"introspection_invalid", "invalid_introspection"},
		{"control_unknown", "unexpected_control_message"},
		{"control_payload_invalid", "unexpected_control_message"},
		{"too_large", "value_size_exceeded"},
		{"unknown_format", "undecodable_bson_payload"},
		{"malformed", "undecodable_bson_payload"},
		{"no_value", "unexpected_value_type"},
		{"bad_timestamp", "unexpected_value_type"},
		{"type_mismatch", "unexpected_value_type"},
		{"value_too_large", "value_size_exceeded"},
		{"bad_object", "unexpected_object_key"},
		{"unset_not_allowed", "unexpected_value_type"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := UpstreamErrorName(tt.input)
			if got != tt.expected {
				t.Errorf("UpstreamErrorName(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestUpstreamErrorNamePassthrough(t *testing.T) {
	for _, name := range UpstreamErrorNames() {
		t.Run(name, func(t *testing.T) {
			got := UpstreamErrorName(name)
			if got != name {
				t.Errorf("UpstreamErrorName(%q) = %q, want %q (passthrough)", name, got, name)
			}
		})
	}
}

func TestUpstreamErrorNameClosedSetInvariant(t *testing.T) {
	fallback := "interface_loading_failed"
	unknowns := []string{"", "RejectReason(42)", "totally_made_up", "Device_Unknown", " device_unknown"}
	set := make(map[string]bool)
	for _, n := range UpstreamErrorNames() {
		set[n] = true
	}
	for _, input := range unknowns {
		t.Run(input, func(t *testing.T) {
			got := UpstreamErrorName(input)
			if !set[got] {
				t.Errorf("UpstreamErrorName(%q) = %q is not in UpstreamErrorNames()", input, got)
			}
			if got != fallback {
				t.Errorf("UpstreamErrorName(%q) = %q, want %q", input, got, fallback)
			}
		})
	}
}

func TestNewDeviceErrorEventWiring(t *testing.T) {
	t.Run("translated", func(t *testing.T) {
		evt := NewDeviceErrorEvent("unknown_format", map[string]string{"detail": "x"})
		if evt.ErrorName != "undecodable_bson_payload" {
			t.Errorf("ErrorName = %q, want undecodable_bson_payload", evt.ErrorName)
		}
		if evt.Metadata["detail"] != "x" {
			t.Errorf("Metadata[detail] = %q, want x", evt.Metadata["detail"])
		}
		if evt.Metadata["astrate_reason"] != "unknown_format" {
			t.Errorf("Metadata[astrate_reason] = %q, want unknown_format", evt.Metadata["astrate_reason"])
		}
	})
	t.Run("untranslated", func(t *testing.T) {
		evt := NewDeviceErrorEvent("mapping_not_found", nil)
		if evt.ErrorName != "mapping_not_found" {
			t.Errorf("ErrorName = %q, want mapping_not_found", evt.ErrorName)
		}
		if _, ok := evt.Metadata["astrate_reason"]; ok {
			t.Errorf("unexpected astrate_reason key in Metadata")
		}
		if evt.Metadata == nil {
			t.Errorf("Metadata is nil, want non-nil empty map")
		}
	})
}
