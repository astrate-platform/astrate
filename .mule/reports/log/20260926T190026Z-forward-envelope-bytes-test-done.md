slug: forward-envelope-bytes-test
verdict: done
at:  ed869b6
ran: 2026-09-26T19:00:26Z on DietPi in 213s


> build · big-pickle

I'll start by reading the two required files.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/engine/forward/envelope.go

$ rg -n marshalEnvelope /root/astrate-mule --glob '*.go'; ls /root/astrate-mule/internal/engine/forward/
/root/astrate-mule/internal/engine/forward/envelope.go:16:// marshalEnvelope builds the envelope's JSON bytes. A nil json.RawMessage
/root/astrate-mule/internal/engine/forward/envelope.go:19:func marshalEnvelope(realm, trigger string, action, event []byte) ([]byte, error) {
/root/astrate-mule/internal/engine/forward/nats.go:45:	payload, err := marshalEnvelope(realm, trigger, action, event)
/root/astrate-mule/internal/engine/forward/http.go:140:	raw, err := marshalEnvelope(realm, trigger, action, event)
envelope.go
http.go
http_test.go
nats.go
nats_test.go

→ Read internal/engine/forward/http_test.go [limit=60]
→ Read internal/engine/forward/nats_test.go [limit=60]

$ rg -n "range tests|for _, tc := range|for _, tt := range|tests := \[\]struct" /root/astrate-mule/internal/engine/forward/*.go | head -20
/root/astrate-mule/internal/engine/forward/http_test.go:220:	tests := []struct {
/root/astrate-mule/internal/engine/forward/http_test.go:234:	for _, tt := range tests {

→ Read internal/engine/forward/http_test.go [offset=205, limit=55]

$ mkdir -p /tmp/opencode/envprobe && cat > /tmp/opencode/envprobe/main.go <<'EOF'
package main

import (
	"encoding/json"
	"fmt"
)

type envelope struct {
	Realm   string          `json:"realm"`
	Trigger string          `json:"trigger"`
	Action  json.RawMessage `json:"action"`
	Event   json.RawMessage `json:"event"`
}

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

func main() {
	cases := []struct {
		r, tg string
		a, e  []byte
	}{
		{"r", "t", nil, nil},
		{"r", "t", []byte{}, []byte{}},
		{"r", "t", nil, []byte{}},
		{"r", "t", []byte(`{"amqp_exchange":"x"}`), []byte(`{"device_id":"abc"}`)},
		{"r", "t", []byte(`{ "a" : 1 }`), []byte(`null`)},
		{"r&t", "t", []byte(`"<b>"`), []byte(`"a&b"`)},
	}
	for i, c := range cases {
		b, err := marshalEnvelope(c.r, c.tg, c.a, c.e)
		fmt.Printf("%d: %s err=%v\n", i, b, err)
	}
}
EOF
cd /tmp/opencode/envprobe && go run main.go
0: {"realm":"r","trigger":"t","action":null,"event":null} err=<nil>
1: {"realm":"r","trigger":"t","action":null,"event":null} err=<nil>
2: {"realm":"r","trigger":"t","action":null,"event":null} err=<nil>
3: {"realm":"r","trigger":"t","action":{"amqp_exchange":"x"},"event":{"device_id":"abc"}} err=<nil>
4: {"realm":"r","trigger":"t","action":{"a":1},"event":null} err=<nil>
5: {"realm":"r\u0026t","trigger":"t","action":"\u003cb\u003e","event":"a\u0026b"} err=<nil>

$ sed -n 90,110p /root/astrate-mule/internal/engine/forward/http.go; echo ---; sed -n 130,155p /root/astrate-mule/internal/engine/forward/http.go; echo ---; sed -n 38,55p /root/astrate-mule/internal/engine/forward/nats.go
return &HTTP{
		url:    cfg.URL,
		method: m,
		static: cfg.StaticHeaders,
		client: c,
	}, nil
}

// validToken returns true when s is a non-empty token per RFC 7230 §3.2.6.
// That is the grammar of the method (RFC 7230 §4.1) and, identically, of the
// header field name (RFC 7230 §3.2), so both are checked with this one
// predicate. We accept any single-token string so that user-configured
// methods like PUT or PATCH are allowed without a whitelist.
func validToken(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case 'A' <= c && c <= 'Z', 'a' <= c && c <= 'z', '0' <= c && c <= '9':
---
		c := v[i]
		if (c < 0x20 && c != '\t') || c == 0x7f {
			return false
		}
	}
	return true
}

// Forward implements the triggers.Forwarder contract.
func (h *HTTP) Forward(ctx context.Context, realm, trigger string, action json.RawMessage, event []byte) error {
	raw, err := marshalEnvelope(realm, trigger, action, event)
	if err != nil {
		return fmt.Errorf("forward: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, h.method, h.url, bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("forward: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Astarte-Realm", realm)
	req.Header.Set("Astrate-Trigger-Name", trigger)
	for k, v := range h.static {
		req.Header.Set(k, v)
	}
	resp, err := h.client.Do(req)
	if err != nil {
---
		return nil, fmt.Errorf("forward: %w", err)
	}
	return &NATS{conn: conn, subject: cfg.Subject}, nil
}

// Forward implements the triggers.Forwarder contract.
func (n *NATS) Forward(_ context.Context, realm, trigger string, action json.RawMessage, event []byte) error {
	payload, err := marshalEnvelope(realm, trigger, action, event)
	if err != nil {
		return fmt.Errorf("forward: marshal: %w", err)
	}
	if err := n.conn.Publish(n.subject, payload); err != nil {
		return fmt.Errorf("forward: %w", err)
	}
	return nil
}

← Write internal/engine/forward/envelope_test.go
Wrote file successfully.

$ cd /root/astrate-mule && gofmt -w internal/engine/forward/envelope_test.go && gofmt -l . && go test ./internal/engine/forward/ -run TestMarshalEnvelopeBytes -v 2>&1 | tail -25
=== RUN   TestMarshalEnvelopeBytes
=== RUN   TestMarshalEnvelopeBytes/nil_pair
=== RUN   TestMarshalEnvelopeBytes/empty_non-nil_pair
=== RUN   TestMarshalEnvelopeBytes/empty_action_with_event
=== RUN   TestMarshalEnvelopeBytes/valued_pair
=== RUN   TestMarshalEnvelopeBytes/padded_pair_is_compacted
--- PASS: TestMarshalEnvelopeBytes (0.00s)
    --- PASS: TestMarshalEnvelopeBytes/nil_pair (0.00s)
    --- PASS: TestMarshalEnvelopeBytes/empty_non-nil_pair (0.00s)
    --- PASS: TestMarshalEnvelopeBytes/empty_action_with_event (0.00s)
    --- PASS: TestMarshalEnvelopeBytes/valued_pair (0.00s)
    --- PASS: TestMarshalEnvelopeBytes/padded_pair_is_compacted (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/engine/forward	0.019s

$ cd /root/astrate-mule && cp internal/engine/forward/envelope.go /tmp/opencode/envelope.go.orig && sed -i 's/json:"action"/json:"Action"/' internal/engine/forward/envelope.go && go test ./internal/engine/forward/ -run TestMarshalEnvelopeBytes 2>&1 | tail -6; echo "=== mutation 2: drop normalisation ==="; cp /tmp/opencode/envelope.go.orig internal/engine/forward/envelope.go && python3 - <<'EOF'
p='internal/engine/forward/envelope.go'
s=open(p).read()
s=s.replace("\tif len(action) > 0 {\n\t\ta = action\n\t}\n","\ta = action\n").replace("\tif len(event) > 0 {\n\t\te = event\n\t}\n","\te = event\n")
open(p,'w').write(s)
EOF
go test ./internal/engine/forward/ -run TestMarshalEnvelopeBytes 2>&1 | tail -8; cp /tmp/opencode/envelope.go.orig internal/engine/forward/envelope.go; diff /tmp/opencode/envelope.go.orig internal/engine/forward/envelope.go && echo "restored-clean"
envelope_test.go:70: envelope bytes
             got: {"realm":"r","trigger":"t","Action":{"amqp_exchange":"x"},"event":{}}
            want: {"realm":"r","trigger":"t","action":{"amqp_exchange":"x"},"event":{}}
FAIL
FAIL	github.com/astrate-platform/astrate/internal/engine/forward	0.019s
FAIL
=== mutation 2: drop normalisation ===
--- FAIL: TestMarshalEnvelopeBytes (0.00s)
    --- FAIL: TestMarshalEnvelopeBytes/empty_non-nil_pair (0.00s)
        envelope_test.go:67: marshalEnvelope: json: error calling MarshalJSON for type json.RawMessage: unexpected end of JSON input
    --- FAIL: TestMarshalEnvelopeBytes/empty_action_with_event (0.00s)
        envelope_test.go:67: marshalEnvelope: json: error calling MarshalJSON for type json.RawMessage: unexpected end of JSON input
FAIL
FAIL	github.com/astrate-platform/astrate/internal/engine/forward	0.030s
FAIL
restored-clean

$ cd /root/astrate-mule && go vet ./... && echo "VET OK" && go test -race ./... 2>&1 | grep -v "^ok\|no test files" | head -20; echo "race-exit:$?"
VET OK
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/cmd/astrate	0.083s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/examples/flow-container-echo	0.122s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine	0.058s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine/channels	0.065s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine/stream	0.084s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/auth	0.073s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
race-exit:0

$ cd /root/astrate-mule && go test ./... 2>&1 | grep -v "no test files" | grep -v "^ok" | head -20; echo "plain-exit:${PIPESTATUS:-see-above}"; go test ./... > /tmp/opencode/plain.log 2>&1; echo "exit=$?"; grep -c "^ok" /tmp/opencode/plain.log; gofmt -l .; echo "gofmt-clean"
plain-exit:see-above
exit=0
31
gofmt-clean

$ cd /root/astrate-mule && git status --porcelain && rm -rf /tmp/opencode/envprobe /tmp/opencode/envelope.go.orig /tmp/opencode/plain.log
?? internal/engine/forward/envelope_test.go

```
Done: Added internal/engine/forward/envelope_test.go — a table pinning marshalEnvelope's exact output bytes for (nil,nil), the empty-non-nil pair from envelope.go:16-18, a mixed empty/valued pair, a valued pair, and a padded pair (compaction), with no HTTP server and no containers.
Files: internal/engine/forward/envelope_test.go (new; no source files changed)
Verified: go vet ./... && go test ./... -> pass (31 pkgs ok, exit 0); gofmt -l . -> clean; go test -race ./... -> FAILS everywhere with "ThreadSanitizer: unsupported VMA range (39-bit)" — this is the Pi, race builds are unusable here, as MULE.md warns. Mutation-checked: renaming the action json tag to "Action" fails the table, and dropping the empty-normalisation fails the two empty subtests, so the test does break without the implementation. No concurrency touched.
Unsure: I deliberately did not pin Go's HTML escaping (json.Marshal turns `&`/`<`/`>` into \u0026 etc., so a realm or action containing them is not forwarded verbatim despite DESIGN.md decision 10). Whether that escaping is intended is a design question, not a test question, so I left it unasserted rather than cementing behaviour that may contradict the design.
Follow-ups: none
```
