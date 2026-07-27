package engine

import (
	"context"
	"fmt"

	"github.com/astrate-platform/astrate/internal/broker"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Known capability keys (upstream MQTT v1 §2 — Device capabilities).
const (
	capPurgePropertiesCompressionFormat = "purge_properties_compression_format"
)

// Valid values for purge_properties_compression_format.
const (
	compressionZlib      = "zlib"
	compressionPlaintext = "plaintext"
)

// handleCapabilities parses the BSON capabilities payload published by the
// device on `<realm>/<device_id>/capabilities` and stores the recognised
// capabilities in the device's in-memory state. The only upstream capability
// today is purge_properties_compression_format, which controls whether the
// server sends consumer/properties payloads compressed (zlib, the default) or
// in plaintext.
func (e *Engine) handleCapabilities(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
	caps, err := decodeCapabilities(m.Payload)
	if err != nil {
		e.reject(m, reasonCapabilitiesInvalid, err.Error())
		return
	}
	for k, v := range caps {
		if err := validateCapability(k, v); err != nil {
			e.reject(m, reasonCapabilitiesInvalid, err.Error())
			return
		}
	}
	dev, ok := e.deviceState(ctx, m, realm)
	if !ok {
		return
	}
	if v, ok := caps[capPurgePropertiesCompressionFormat]; ok {
		dev.setPurgeCompression(v)
	}
	m.Ack()
}

// decodeCapabilities decodes a BSON capabilities document into a flat
// string→string map. Unknown keys are silently accepted (forward-compatible
// with future upstream capabilities). The document must be valid BSON.
func decodeCapabilities(p []byte) (map[string]string, error) {
	if len(p) == 0 {
		return nil, fmt.Errorf("empty capabilities payload")
	}
	var raw bson.M
	if err := bson.Unmarshal(p, &raw); err != nil {
		return nil, fmt.Errorf("capabilities payload is not valid BSON: %w", err)
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		if s, ok := v.(string); ok {
			out[k] = s
		}
	}
	return out, nil
}

// validateCapability checks that a capability value is within its allowed
// set. Returns a descriptive error for invalid values.
func validateCapability(key, value string) error {
	switch key {
	case capPurgePropertiesCompressionFormat:
		switch value {
		case compressionZlib, compressionPlaintext:
			return nil
		default:
			return fmt.Errorf("invalid value %q for %s (must be %q or %q)",
				value, key, compressionZlib, compressionPlaintext)
		}
	default:
		// Unknown capabilities are accepted (forward-compatible).
		return nil
	}
}

// purgeCompressionFor returns the consumer/properties compression format for
// the device: "zlib" or "plaintext". The store fallback covers the case
// where the device state was evicted and reloaded from the store (the
// capability is session-scoped and re-published on connect).
func purgeCompressionFor(dev *deviceState) string {
	if dev == nil {
		return compressionZlib
	}
	f := dev.purgeCompressionFormat()
	if f == "" {
		return compressionZlib
	}
	return f
}
