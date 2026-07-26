package channels

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestEncodeJoin(t *testing.T) {
	f := Frame{
		JoinRef: strPtr("1"),
		Ref:     strPtr("1"),
		Topic:   "rooms:test:dashboard_abc_42",
		Event:   EventPhxJoin,
	}
	b, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got := string(b)
	want := `["1","1","rooms:test:dashboard_abc_42","phx_join",{}]`
	if got != want {
		t.Errorf("EncodeJoin:\n  got  %s\n  want %s", got, want)
	}
}

func TestEncodeHeartbeat(t *testing.T) {
	f := Frame{
		JoinRef: nil,
		Ref:     strPtr("3"),
		Topic:   TopicHeartbeat,
		Event:   EventHeartbeat,
	}
	b, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got := string(b)
	want := `[null,"3","phoenix","heartbeat",{}]`
	if got != want {
		t.Errorf("EncodeHeartbeat:\n  got  %s\n  want %s", got, want)
	}
}

func TestEncodeServerPush(t *testing.T) {
	f := Frame{
		JoinRef: nil,
		Ref:     nil,
		Topic:   "rooms:test:dashboard_abc_42",
		Event:   "new_event",
		Payload: json.RawMessage(`{"device_id":"x"}`),
	}
	b, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got := string(b)
	want := `[null,null,"rooms:test:dashboard_abc_42","new_event",{"device_id":"x"}]`
	if got != want {
		t.Errorf("EncodeServerPush:\n  got  %s\n  want %s", got, want)
	}
}

func TestEncodeOKReply(t *testing.T) {
	in := Frame{
		JoinRef: strPtr("1"),
		Ref:     strPtr("1"),
		Topic:   "rooms:test:dashboard_abc_42",
		Event:   EventPhxJoin,
	}
	reply, err := OK(in, nil)
	if err != nil {
		t.Fatalf("OK: %v", err)
	}
	b, err := json.Marshal(reply)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got := string(b)
	want := `["1","1","rooms:test:dashboard_abc_42","phx_reply",{"status":"ok","response":{}}]`
	if got != want {
		t.Errorf("EncodeOKReply:\n  got  %s\n  want %s", got, want)
	}
}

func TestEncodeErrReply(t *testing.T) {
	in := Frame{
		JoinRef: strPtr("1"),
		Ref:     strPtr("1"),
		Topic:   "rooms:test:dashboard_abc_42",
		Event:   EventPhxJoin,
	}
	reply, err := Err(in, "unauthorized")
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	b, err := json.Marshal(reply)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got := string(b)
	want := `["1","1","rooms:test:dashboard_abc_42","phx_reply",{"status":"error","response":{"reason":"unauthorized"}}]`
	if got != want {
		t.Errorf("EncodeErrReply:\n  got  %s\n  want %s", got, want)
	}
}

func TestDecodeRoundTrip(t *testing.T) {
	goldens := []struct {
		name string
		json string
	}{
		{"join", `["1","1","rooms:test:dashboard_abc_42","phx_join",{}]`},
		{"heartbeat", `[null,"3","phoenix","heartbeat",{}]`},
		{"server_push", `[null,null,"rooms:test:dashboard_abc_42","new_event",{"device_id":"x"}]`},
	}
	for _, g := range goldens {
		t.Run(g.name, func(t *testing.T) {
			var f Frame
			if err := json.Unmarshal([]byte(g.json), &f); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}
			b, err := json.Marshal(f)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			got := string(b)
			if got != g.json {
				t.Errorf("round trip %s:\n  got  %s\n  want %s", g.name, got, g.json)
			}
		})
	}
}

func TestDecodeNullRefsArePreserved(t *testing.T) {
	heartbeat := `[null,"3","phoenix","heartbeat",{}]`
	var f Frame
	if err := json.Unmarshal([]byte(heartbeat), &f); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if f.JoinRef != nil {
		t.Errorf("JoinRef should be nil, got %q", *f.JoinRef)
	}
	if f.Ref == nil {
		t.Fatal("Ref should not be nil")
	}
	if *f.Ref != "3" {
		t.Errorf("Ref: got %q, want %q", *f.Ref, "3")
	}

	// An empty-string Ref must encode as "" and not null.
	emptyRef := Frame{
		Ref:   strPtr(""),
		Topic: "x",
		Event: "y",
	}
	b, err := json.Marshal(emptyRef)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	// Parse the raw array and check that the second element (Ref) is "", not null.
	var arr [5]json.RawMessage
	if err := json.Unmarshal(b, &arr); err != nil {
		t.Fatalf("unmarshal for check: %v", err)
	}
	if string(arr[1]) != `""` {
		t.Errorf("empty Ref element: got %s, want \"\"", string(arr[1]))
	}
}

func TestDecodeRejectsMalformed(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"json_object", `{"topic":"x"}`},
		{"four_elements", `[1,2,3,4]`},
		{"six_elements", `[1,2,3,4,5,6]`},
		{"empty_array", `[]`},
		{"not_json", `hello`},
		{"four_elements_well_typed", `["1","1","rooms:test:x","phx_join"]`},
		{"six_elements_well_typed", `["1","1","rooms:test:x","phx_join",{},"extra"]`},
		{"non_string_ref", `["1",42,"rooms:test:x","phx_join",{}]`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var f Frame
			err := json.Unmarshal([]byte(c.input), &f)
			if err == nil {
				t.Errorf("expected error for %s, got nil (frame: %+v)", c.name, f)
			}
		})
	}
}

func TestDecodeRejectsWrongLengthBeforeFields(t *testing.T) {
	input := `["1","1","rooms:test:x","phx_join"]`
	var f Frame
	err := json.Unmarshal([]byte(input), &f)
	if err == nil {
		t.Fatal("expected error for four_elements_well_typed, got nil")
	}
	if got := err.Error(); !strings.Contains(got, "4") {
		t.Errorf("error should mention length 4, got: %s", got)
	}
}

func strPtr(s string) *string { return &s }

func init() {
	// Ensure fmt is used to avoid import error.
	_ = fmt.Sprintf
}
