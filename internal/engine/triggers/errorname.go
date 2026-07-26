package triggers

// upstreamErrorNames is the closed set astarte-dashboard 1.2.2 validates
// error_name against (its bundle's `nke` array). An event whose error_name is
// outside this set fails the client's validator and is discarded whole.
var upstreamErrorNames = []string{
	"write_on_server_owned_interface",
	"invalid_interface",
	"invalid_path",
	"mapping_not_found",
	"interface_loading_failed",
	"ambiguous_path",
	"undecodable_bson_payload",
	"unexpected_value_type",
	"value_size_exceeded",
	"unexpected_object_key",
	"invalid_introspection",
	"unexpected_control_message",
	"device_session_not_found",
	"resend_interface_properties_failed",
	"empty_cache_error",
}

// upstreamFallback is what an unrecognised reason becomes: upstream's most
// generic "the platform could not process this". The true reason is never lost
// — NewDeviceErrorEvent keeps it under metadata["astrate_reason"].
const upstreamFallback = "interface_loading_failed"

var upstreamSet = func() map[string]bool {
	s := make(map[string]bool, len(upstreamErrorNames))
	for _, n := range upstreamErrorNames {
		s[n] = true
	}
	return s
}()

// UpstreamErrorNames returns a copy of the closed set, so a caller cannot
// mutate the package's own slice.
func UpstreamErrorNames() []string {
	return append([]string(nil), upstreamErrorNames...)
}

var astrateToUpstream = map[string]string{
	"realm_unknown":                  "device_session_not_found",
	"device_unknown":                 "device_session_not_found",
	"malformed_topic":                "invalid_path",
	"interface_not_in_introspection": "invalid_interface",
	"interface_not_installed":        "interface_loading_failed",
	"ownership_violation":            "write_on_server_owned_interface",
	"unexpected_path":                "mapping_not_found",
	"introspection_invalid":          "invalid_introspection",
	"control_unknown":                "unexpected_control_message",
	"control_payload_invalid":        "unexpected_control_message",
	"too_large":                      "value_size_exceeded",
	"unknown_format":                 "undecodable_bson_payload",
	"malformed":                      "undecodable_bson_payload",
	"no_value":                       "undecodable_bson_payload",
	"bad_timestamp":                  "unexpected_value_type",
	"type_mismatch":                  "unexpected_value_type",
	"value_too_large":                "value_size_exceeded",
	"bad_object":                     "unexpected_object_key",
	"unset_not_allowed":              "unexpected_value_type",
}

// UpstreamErrorName maps one of Astrate's reject-reason labels to the upstream
// error_name enum value. If the input is already an upstream enum value it is
// returned unchanged. Any other input maps to interface_loading_failed.
func UpstreamErrorName(reason string) string {
	if upstreamSet[reason] {
		return reason
	}
	if m, ok := astrateToUpstream[reason]; ok {
		return m
	}
	return upstreamFallback
}
