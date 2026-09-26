package forward

import (
	"testing"
)

// TestMarshalEnvelopeBytes pins the exact wire bytes of an envelope. The
// other suites in this package unmarshal into bodyShape, which cannot see key
// order or the literal bytes, so a renamed or reordered field would pass them
// unchanged. These expectations are therefore compared as strings.
func TestMarshalEnvelopeBytes(t *testing.T) {
	tests := []struct {
		name           string
		realm, trigger string
		action, event  []byte
		want           string
	}{
		{
			name:    "nil pair",
			realm:   "r",
			trigger: "t",
			want:    `{"realm":"r","trigger":"t","action":null,"event":null}`,
		},
		{
			// The form the doc comment at envelope.go:16-18 is about: a
			// non-nil but empty RawMessage marshals to nothing and would
			// corrupt the envelope, so it must come out as null too.
			name:    "empty non-nil pair",
			realm:   "r",
			trigger: "t",
			action:  []byte{},
			event:   []byte{},
			want:    `{"realm":"r","trigger":"t","action":null,"event":null}`,
		},
		{
			name:    "empty action with event",
			realm:   "r",
			trigger: "t",
			action:  []byte{},
			event:   []byte(`{"device_id":"abc"}`),
			want:    `{"realm":"r","trigger":"t","action":null,"event":{"device_id":"abc"}}`,
		},
		{
			name:    "valued pair",
			realm:   "r",
			trigger: "t",
			action:  []byte(`{"amqp_exchange":"x"}`),
			event:   []byte(`{"device_id":"abc"}`),
			want:    `{"realm":"r","trigger":"t","action":{"amqp_exchange":"x"},"event":{"device_id":"abc"}}`,
		},
		{
			// json.Marshal compacts embedded RawMessage, so the padding a
			// caller happens to pass must not reach the bus.
			name:    "padded pair is compacted",
			realm:   "r",
			trigger: "t",
			action:  []byte("{\n  \"amqp_exchange\" : \"x\"\n}"),
			event:   []byte(`{ }`),
			want:    `{"realm":"r","trigger":"t","action":{"amqp_exchange":"x"},"event":{}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := marshalEnvelope(tt.realm, tt.trigger, tt.action, tt.event)
			if err != nil {
				t.Fatalf("marshalEnvelope: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("envelope bytes\n got: %s\nwant: %s", got, tt.want)
			}
		})
	}
}
