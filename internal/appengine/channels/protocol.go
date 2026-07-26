// Package channels implements the upstream Phoenix Channels V2 wire protocol
// served at /appengine/v1/socket/websocket. It is the compatibility counterpart
// to the Astrate-native socket in internal/appengine/stream.
package channels

import (
	"encoding/json"
	"fmt"
)

// Event-name constants used by the Phoenix V2 protocol.
const (
	EventPhxJoin   = "phx_join"
	EventPhxReply  = "phx_reply"
	EventPhxLeave  = "phx_leave"
	EventPhxClose  = "phx_close"
	EventPhxError  = "phx_error"
	EventHeartbeat = "heartbeat"

	// TopicHeartbeat is the reserved topic for heartbeat messages.
	TopicHeartbeat = "phoenix"
)

// Frame is one Phoenix V2 message on the wire, encoded as a five-element JSON
// array: [join_ref, ref, topic, event, payload]. A nil JoinRef or Ref is
// marshalled as JSON null; the two are distinguished from the empty string.
type Frame struct {
	JoinRef *string
	Ref     *string
	Topic   string
	Event   string
	Payload json.RawMessage
}

// MarshalJSON renders the frame as a five-element JSON array. A nil JoinRef or
// Ref becomes JSON null. A nil or zero-length Payload becomes {}.
func (f Frame) MarshalJSON() ([]byte, error) {
	payload := json.RawMessage(`{}`)
	if len(f.Payload) > 0 {
		payload = f.Payload
	}
	return json.Marshal([5]interface{}{f.JoinRef, f.Ref, f.Topic, f.Event, payload})
}

// UnmarshalJSON parses a five-element JSON array back into the frame. It
// rejects anything that is not a JSON array of exactly five elements.
func (f *Frame) UnmarshalJSON(b []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("channels.Frame: invalid JSON: %w", err)
	}
	if len(raw) != 5 {
		return fmt.Errorf("channels.Frame: expected a 5-element JSON array, got %d elements", len(raw))
	}
	// Decode join_ref and ref as nullable strings.
	jr, err := unmarshalNullString(raw[0])
	if err != nil {
		return fmt.Errorf("channels.Frame: join_ref: %w", err)
	}
	f.JoinRef = jr
	rr, err := unmarshalNullString(raw[1])
	if err != nil {
		return fmt.Errorf("channels.Frame: ref: %w", err)
	}
	f.Ref = rr
	if err := json.Unmarshal(raw[2], &f.Topic); err != nil {
		return fmt.Errorf("channels.Frame: topic: %w", err)
	}
	if err := json.Unmarshal(raw[3], &f.Event); err != nil {
		return fmt.Errorf("channels.Frame: event: %w", err)
	}
	f.Payload = json.RawMessage(raw[4])
	return nil
}

// unmarshalNullString decodes a JSON value that is either a string or null,
// returning *string. Null maps to nil; the empty string maps to a pointer to "".
func unmarshalNullString(raw json.RawMessage) (*string, error) {
	s := string(raw)
	if s == "null" {
		return nil, nil
	}
	var out string
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("expected a string or null, got %s", raw)
	}
	return &out, nil
}

// Reply builds a phx_reply answer to in. The JoinRef and Ref are echoed from in
// so the client can route the reply to its caller. The payload is
// {"status":"<status>","response":<response>}. A nil response renders as {}.
func Reply(in Frame, status string, response any) (Frame, error) {
	respJSON := json.RawMessage(`{}`)
	if response != nil {
		b, err := json.Marshal(response)
		if err != nil {
			return Frame{}, fmt.Errorf("channels.Reply: marshal response: %w", err)
		}
		respJSON = b
	}
	type replyPayload struct {
		Status   string          `json:"status"`
		Response json.RawMessage `json:"response"`
	}
	payload, err := json.Marshal(replyPayload{Status: status, Response: respJSON})
	if err != nil {
		return Frame{}, fmt.Errorf("channels.Reply: marshal payload: %w", err)
	}
	return Frame{
		JoinRef: in.JoinRef,
		Ref:     in.Ref,
		Topic:   in.Topic,
		Event:   EventPhxReply,
		Payload: payload,
	}, nil
}

// OK is a convenience wrapper over Reply that sends a "ok" status.
func OK(in Frame, response any) (Frame, error) {
	return Reply(in, "ok", response)
}

// Err is a convenience wrapper over Reply that sends an "error" status with a
// reason field.
func Err(in Frame, reason string) (Frame, error) {
	return Reply(in, "error", map[string]string{"reason": reason})
}
