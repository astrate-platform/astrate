package forward

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

// Assert at compile time that *HTTP satisfies the triggers.Forwarder contract
// without importing the triggers package (which would create an import cycle).
type forwarder interface {
	Forward(ctx context.Context, realm, trigger string, action json.RawMessage, event []byte) error
}

var _ forwarder = (*HTTP)(nil)

// bodyShape is used to unmarshal the forwarded envelope for inspection.
type bodyShape struct {
	Realm   string          `json:"realm"`
	Trigger string          `json:"trigger"`
	Action  json.RawMessage `json:"action"`
	Event   json.RawMessage `json:"event"`
}

func TestHappyPath(t *testing.T) {
	const realm = "test-realm"
	const trigger = "test-trigger"
	action := json.RawMessage(`{"amqp_exchange":"x"}`)
	event := json.RawMessage(`{"device_id":"abc"}`)

	var gotMethod, gotCT, gotRealm, gotTrigger string
	var gotAction, gotEvent json.RawMessage

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotCT = r.Header.Get("Content-Type")
		gotRealm = r.Header.Get("Astarte-Realm")
		gotTrigger = r.Header.Get("Astrate-Trigger-Name")
		raw, _ := io.ReadAll(r.Body)
		var b bodyShape
		if err := json.Unmarshal(raw, &b); err != nil {
			t.Fatalf("server: unmarshal body: %v", err)
		}
		gotAction = b.Action
		gotEvent = b.Event
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f, err := New(Config{URL: srv.URL})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	err = f.Forward(context.Background(), realm, trigger, action, event)
	if err != nil {
		t.Fatalf("Forward: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotCT != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", gotCT)
	}
	if gotRealm != realm {
		t.Errorf("Astarte-Realm = %q, want %q", gotRealm, realm)
	}
	if gotTrigger != trigger {
		t.Errorf("Astrate-Trigger-Name = %q, want %q", gotTrigger, trigger)
	}
	if string(gotAction) != string(action) {
		t.Errorf("action = %s, want %s", gotAction, action)
	}
	if string(gotEvent) != string(event) {
		t.Errorf("event = %s, want %s", gotEvent, event)
	}
}

func TestStaticHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if v := r.Header.Get("X-Custom"); v != "custom-val" {
			t.Errorf("X-Custom = %q, want custom-val", v)
		}
		if v := r.Header.Get("X-Other"); v != "other-val" {
			t.Errorf("X-Other = %q, want other-val", v)
		}
		// Fixed headers must still be present.
		if v := r.Header.Get("Content-Type"); v != "application/json" {
			t.Errorf("Content-Type displaced: %q", v)
		}
		if v := r.Header.Get("Astarte-Realm"); v != "r" {
			t.Errorf("Astarte-Realm displaced: %q", v)
		}
		if v := r.Header.Get("Astrate-Trigger-Name"); v != "t" {
			t.Errorf("Astrate-Trigger-Name displaced: %q", v)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f, err := New(Config{
		URL:           srv.URL,
		StaticHeaders: map[string]string{"X-Custom": "custom-val", "X-Other": "other-val"},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := f.Forward(context.Background(), "r", "t", json.RawMessage(`{}`), []byte(`{}`)); err != nil {
		t.Fatalf("Forward: %v", err)
	}
}

// TestStaticHeadersOverrideFixed pins the precedence rule the Config comment
// states ("applied after the fixed ones") and Forward implements: a static
// header whose name collides with one of the three fixed headers wins, because
// the static loop runs last. TestStaticHeaders only covers the non-colliding
// case, so hoisting the loop above the three fixed Set calls — or filtering the
// fixed names out of h.static — would leave the suite green while inverting
// what the bus receives.
//
// The keys are lowercase on purpose: Header.Set canonicalises, so a TOML
// operator writing `astarte-realm` collides with the fixed Astarte-Realm even
// though the strings differ. The realm header is the interesting one — a bus
// filters and routes on it, so "spoofed" is the value that arrives.
func TestStaticHeadersOverrideFixed(t *testing.T) {
	const want = "spoofed"
	var gotContentType, gotRealm, gotTrigger string
	var gotRealmValues []string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		gotRealm = r.Header.Get("Astarte-Realm")
		gotRealmValues = r.Header.Values("Astarte-Realm")
		gotTrigger = r.Header.Get("Astrate-Trigger-Name")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f, err := New(Config{
		URL: srv.URL,
		StaticHeaders: map[string]string{
			"content-type":         "application/x-ndjson",
			"astarte-realm":        want,
			"astrate-trigger-name": "other-trigger",
		},
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := f.Forward(context.Background(), "r", "t", json.RawMessage(`{}`), []byte(`{}`)); err != nil {
		t.Fatalf("Forward: %v", err)
	}

	if gotRealm != want {
		t.Errorf("Astarte-Realm = %q, want the static value %q", gotRealm, want)
	}
	if len(gotRealmValues) != 1 {
		t.Errorf("Astarte-Realm sent %d time(s) (%q), want a single Set value", len(gotRealmValues), gotRealmValues)
	}
	if gotContentType != "application/x-ndjson" {
		t.Errorf("Content-Type = %q, want the static value", gotContentType)
	}
	if gotTrigger != "other-trigger" {
		t.Errorf("Astrate-Trigger-Name = %q, want the static value", gotTrigger)
	}
}

func TestMethodHonoured(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %q, want PUT", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f, err := New(Config{URL: srv.URL, Method: "PUT"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := f.Forward(context.Background(), "r", "t", json.RawMessage(`{}`), []byte(`{}`)); err != nil {
		t.Fatalf("Forward: %v", err)
	}
}

func TestNilActionAndEvent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var b bodyShape
		if err := json.Unmarshal(raw, &b); err != nil {
			t.Fatalf("server: unmarshal body: %v", err)
		}
		if string(b.Action) != "null" {
			t.Errorf("action = %s, want null", b.Action)
		}
		if string(b.Event) != "null" {
			t.Errorf("event = %s, want null", b.Event)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f, err := New(Config{URL: srv.URL})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := f.Forward(context.Background(), "r", "t", nil, nil); err != nil {
		t.Fatalf("Forward: %v", err)
	}
}

func TestStatusTable(t *testing.T) {
	tests := []struct {
		code int
		want bool // true = nil error expected
	}{
		{200, true},
		{204, true},
		{299, true},
		{300, false},
		{400, false},
		{404, false},
		{500, false},
		{503, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.code), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.code)
			}))
			defer srv.Close()

			f, err := New(Config{URL: srv.URL})
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			err = f.Forward(context.Background(), "r", "t", json.RawMessage(`{}`), []byte(`{}`))
			if tt.want && err != nil {
				t.Errorf("expected nil error, got %v", err)
			}
			if !tt.want && err == nil {
				t.Errorf("expected non-nil error for status %d", tt.code)
			}
			if !tt.want && err != nil {
				want := fmt.Sprintf("%d", tt.code)
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q does not mention status %d", err.Error(), tt.code)
				}
			}
		})
	}
}

// TestEmptyNonNilActionAndEvent pins the normalisation in Forward that a nil
// json.RawMessage does not need: a non-nil, zero-length action or event
// marshals to nothing and would corrupt the whole envelope. Removing the
// len() guards leaves TestNilActionAndEvent green, so this is the row that
// actually binds the rule.
// TestStatusErrorCarriesBody pins the half of a failed delivery an operator
// needs. A bare "forward: status 500" survives the bus's own explanation
// nowhere — the drain at the end of Forward throws the body away — so a 500
// from the bus is unactionable. This is the row that breaks if the body prefix
// is dropped from the error: TestStatusTable only pins that the status code is
// mentioned, and it writes no body.
func TestStatusErrorCarriesBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	}))
	defer srv.Close()

	f, err := New(Config{URL: srv.URL})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	err = f.Forward(context.Background(), "r", "t", json.RawMessage(`{}`), []byte(`{}`))
	if err == nil {
		t.Fatal("expected non-nil error for status 500")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error %q does not mention the status", err.Error())
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Errorf("error %q does not mention the response body", err.Error())
	}
}

// TestStatusErrorBodyBounded covers the other half: the body quoted into the
// error is capped, so a peer that answers 500 with a megabyte of text cannot
// turn one log line into a megabyte. Without the cap this test still passes on
// content, so it asserts the bound itself.
func TestStatusErrorBodyBounded(t *testing.T) {
	const filler = "0123456789"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(strings.Repeat(filler, 4096)))
	}))
	defer srv.Close()

	f, err := New(Config{URL: srv.URL})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	err = f.Forward(context.Background(), "r", "t", json.RawMessage(`{}`), []byte(`{}`))
	if err == nil {
		t.Fatal("expected non-nil error for status 500")
	}
	// One "forward: status 500: " prefix plus the 512-byte cap, the trailing
	// ellipsis marking the cut, and slack for the numbers in the prefix.
	if n := len(err.Error()); n > 600 {
		t.Errorf("error is %d bytes, want the body prefix bounded to ~512", n)
	}
	if !strings.Contains(err.Error(), "…") {
		t.Errorf("error %q does not mark the body as truncated", err.Error())
	}
}

func TestEmptyNonNilActionAndEvent(t *testing.T) {
	var body []byte
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		body, _ = io.ReadAll(r.Body)
	}))
	defer srv.Close()

	f, err := New(Config{URL: srv.URL})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := f.Forward(context.Background(), "r", "t", json.RawMessage{}, []byte{}); err != nil {
		t.Fatalf("Forward: %v", err)
	}
	var b bodyShape
	if err := json.Unmarshal(body, &b); err != nil {
		t.Fatalf("envelope does not unmarshal: %v (body %q)", err, body)
	}
	if string(b.Action) != "null" {
		t.Errorf("action = %s, want null", b.Action)
	}
	if string(b.Event) != "null" {
		t.Errorf("event = %s, want null", b.Event)
	}
}

func TestTransportError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	closed := srv.URL
	srv.Close() // close immediately so connection fails

	f, err := New(Config{URL: closed})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	err = f.Forward(context.Background(), "r", "t", json.RawMessage(`{}`), []byte(`{}`))
	if err == nil {
		t.Fatal("expected non-nil error for closed server")
	}
}

func TestCancelledContext(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()

	f, err := New(Config{URL: srv.URL})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	err = f.Forward(ctx, "r", "t", json.RawMessage(`{}`), []byte(`{}`))
	if err == nil {
		t.Fatal("expected non-nil error for cancelled context")
	}
}

func TestNewRejectsBadConfig(t *testing.T) {
	var errCount atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		errCount.Add(1)
	}))
	defer srv.Close()

	bad := []struct {
		name string
		cfg  Config
	}{
		{"empty URL", Config{}},
		{"unparseable URL", Config{URL: "://nope"}},
		// A relative URL parses cleanly, so it must be rejected on the
		// absolute-http(s) rule specifically, not on a parse failure —
		// otherwise the misconfiguration surfaces once per delivery instead
		// of at boot.
		{"relative URL", Config{URL: "nope"}},
		{"no host", Config{URL: "http://"}},
		{"non-http scheme", Config{URL: "ftp://bus.example/x"}},
		{"bad method", Config{URL: srv.URL, Method: "GET NOPE"}},
		// A static header is applied per delivery, and net/http validates
		// names and values only at write time, so these are exactly the typos
		// that would boot clean and then fail every delivery. A space in a
		// value is legal, so the value cases must be about control bytes.
		{"bad static header name", Config{URL: srv.URL, StaticHeaders: map[string]string{"X Bad Name": "v"}}},
		{"static header name with colon", Config{URL: srv.URL, StaticHeaders: map[string]string{"X:Foo": "v"}}},
		{"static header value with newline", Config{URL: srv.URL, StaticHeaders: map[string]string{"X-Foo": "a\nb"}}},
		{"static header value with carriage return", Config{URL: srv.URL, StaticHeaders: map[string]string{"X-Foo": "a\rb"}}},
	}
	for _, b := range bad {
		t.Run(b.name, func(t *testing.T) {
			_, err := New(b.cfg)
			if err == nil {
				t.Error("expected non-nil error")
			}
		})
	}
	if n := errCount.Load(); n != 0 {
		t.Errorf("%d request(s) reached the endpoint for rejected configs, want 0", n)
	}

	// A valid config with the same URL must be accepted.
	f, err := New(Config{URL: srv.URL, StaticHeaders: map[string]string{"Authorization": "Bearer tok"}})
	if err != nil {
		t.Fatalf("New with valid config: %v", err)
	}
	if err := f.Forward(context.Background(), "r", "t", json.RawMessage(`{}`), []byte(`{}`)); err != nil {
		t.Fatalf("Forward: %v", err)
	}
	if n := errCount.Load(); n != 1 {
		t.Errorf("%d request(s) reached the endpoint after one Forward, want 1", n)
	}
}
