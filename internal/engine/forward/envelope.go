package forward

import "encoding/json"

// envelope is the wire shape posted to a bus for every forwarded custom
// action (docs/DESIGN.md §1.1, decision 10: forward the action verbatim,
// never interpret it). Every Forwarder implementation in this package must
// produce byte-identical envelopes for the same inputs.
type envelope struct {
	Realm   string          `json:"realm"`
	Trigger string          `json:"trigger"`
	Action  json.RawMessage `json:"action"`
	Event   json.RawMessage `json:"event"`
}

// marshalEnvelope builds the envelope's JSON bytes. A nil json.RawMessage
// marshals to null on its own, but a non-nil empty one marshals to nothing
// and corrupts the whole envelope, so both empty forms are normalised to nil.
func marshalEnvelope(realm, trigger string, action, event []byte) ([]byte, error) {
	var a json.RawMessage
	if len(action) > 0 {
		a = action
	}
	var e json.RawMessage
	if len(event) > 0 {
		e = event
	}
	return json.Marshal(envelope{Realm: realm, Trigger: trigger, Action: a, Event: e})
}
