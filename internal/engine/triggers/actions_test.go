package triggers

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	promtest "github.com/prometheus/client_golang/prometheus/testutil"
)

// newTestExecutor builds an executor with a private registry, a discarding
// logger, and sub-millisecond backoff so retry tests stay fast. The executor
// is Closed on cleanup (Close is idempotent, so tests may also call it).
func newTestExecutor(t *testing.T, cfg ExecutorConfig) *Executor {
	t.Helper()
	cfg.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	if cfg.BackoffStart == 0 {
		cfg.BackoffStart = time.Millisecond
	}
	if cfg.BackoffCap == 0 {
		cfg.BackoffCap = 5 * time.Millisecond
	}
	if cfg.Registerer == nil {
		cfg.Registerer = prometheus.NewRegistry()
	}
	x := NewExecutor(cfg)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := x.Close(ctx); err != nil {
			t.Errorf("executor close: %v", err)
		}
	})
	return x
}

// testDelivery builds a delivery whose action is the given action.
func testDelivery(action *Action) Delivery {
	return Delivery{
		Realm:   "testrealm",
		Trigger: &Trigger{Name: "hook", Action: action},
		Event: SimpleEvent{
			Timestamp:   time.Date(2026, 6, 12, 10, 0, 0, 0, time.UTC),
			DeviceID:    "f0VMRgIBAQAAAAAAAAAAAA",
			TriggerName: "hook",
			Event:       NewIncomingDataEvent("com.ex.Sensors", "/v", 1.0),
		},
	}
}

// outcome reads one outcome counter.
func outcome(x *Executor, label string) float64 {
	return promtest.ToFloat64(x.outcomes.WithLabelValues(label))
}

// eventually polls get until it equals want or the deadline passes.
func eventually(t *testing.T, want float64, get func() float64) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if get() == want {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("value did not reach %v (last %v)", want, get())
}

// TestParseAction validates the action-object parser (docs/ROADMAP.md §7.2
// file 6.12): http_url+method, the legacy http_post_url, custom (non-HTTP)
// actions routed to the Forwarder, and the rejection / unsupported paths.
func TestParseAction(t *testing.T) {
	t.Run("http_url+method lowercases on the wire, uppercases for net/http", func(t *testing.T) {
		a, unsupported, err := parseAction([]byte(`{"http_url":"https://x/h","http_method":"post"}`))
		if err != nil {
			t.Fatal(err)
		}
		if a.Method != http.MethodPost || a.URL != "https://x/h" {
			t.Errorf("got %q %q", a.Method, a.URL)
		}
		if a.Custom != nil || len(unsupported) != 0 {
			t.Errorf("custom=%s unsupported=%v", a.Custom, unsupported)
		}
	})

	t.Run("legacy http_post_url implies POST", func(t *testing.T) {
		a, _, err := parseAction([]byte(`{"http_post_url":"https://y/legacy"}`))
		if err != nil {
			t.Fatal(err)
		}
		if a.Method != http.MethodPost || a.URL != "https://y/legacy" {
			t.Errorf("got %q %q", a.Method, a.URL)
		}
	})

	t.Run("static headers and ignore_ssl_errors carry through", func(t *testing.T) {
		a, _, err := parseAction([]byte(`{"http_url":"https://x","http_method":"put",` +
			`"http_static_headers":{"X-Foo":"bar"},"ignore_ssl_errors":true}`))
		if err != nil {
			t.Fatal(err)
		}
		if a.Method != http.MethodPut || a.StaticHeaders["X-Foo"] != "bar" || !a.IgnoreSSLErrors {
			t.Errorf("got %+v", a)
		}
	})

	t.Run("unknown http_method rejected", func(t *testing.T) {
		if _, _, err := parseAction([]byte(`{"http_url":"https://x","http_method":"frobnicate"}`)); err == nil {
			t.Fatal("want error for unknown method")
		}
	})

	t.Run("non-HTTP action kept verbatim for the forwarder", func(t *testing.T) {
		raw := `{"amqp_exchange":"events","amqp_routing_key":"r"}`
		a, _, err := parseAction([]byte(raw))
		if err != nil {
			t.Fatal(err)
		}
		if a.Method != "" || a.URL != "" || string(a.Custom) != raw {
			t.Errorf("got %+v", a)
		}
	})

	t.Run("mustache template noted as unsupported, default envelope still sent", func(t *testing.T) {
		a, unsupported, err := parseAction([]byte(
			`{"http_url":"https://x","http_method":"post","template":"{{ value }}","template_type":"mustache"}`))
		if err != nil {
			t.Fatal(err)
		}
		if a.URL == "" || len(unsupported) != 1 {
			t.Errorf("action=%+v unsupported=%v", a, unsupported)
		}
	})

	t.Run("missing action rejected", func(t *testing.T) {
		if _, _, err := parseAction(nil); err == nil {
			t.Fatal("want error for missing action")
		}
	})
}

// TestWebhookDelivered: a 2xx response is a successful delivery, and the
// request carries the rendered event JSON, POST, the Content-Type, the
// Astarte-Realm header, and any static headers.
func TestWebhookDelivered(t *testing.T) {
	var (
		gotMethod, gotCT, gotRealm, gotFoo string
		gotBody                            []byte
	)
	done := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotCT = r.Header.Get("Content-Type")
		gotRealm = r.Header.Get("Astarte-Realm")
		gotFoo = r.Header.Get("X-Foo")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
		close(done)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1})
	d := testDelivery(&Action{
		Method: http.MethodPost, URL: srv.URL,
		StaticHeaders: map[string]string{"X-Foo": "bar"},
	})
	x.Enqueue(d)

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("webhook never received")
	}
	eventually(t, 1, func() float64 { return outcome(x, outcomeDelivered) })

	if gotMethod != http.MethodPost {
		t.Errorf("method %q", gotMethod)
	}
	if gotCT != "application/json" {
		t.Errorf("content-type %q", gotCT)
	}
	if gotRealm != "testrealm" {
		t.Errorf("astarte-realm %q", gotRealm)
	}
	if gotFoo != "bar" {
		t.Errorf("x-foo %q", gotFoo)
	}
	wantBody, _ := json.Marshal(d.Event)
	if string(gotBody) != string(wantBody) {
		t.Errorf("body = %s, want %s", gotBody, wantBody)
	}
}

// TestWebhookRetriesThenSucceeds is the docs/ROADMAP.md §7.3 case: a 500
// followed by a 200 retries once and then delivers.
func TestWebhookRetriesThenSucceeds(t *testing.T) {
	var calls atomic.Int32
	reqs := make(chan struct{}, 4)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		reqs <- struct{}{}
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1, MaxAttempts: 5})
	x.Enqueue(testDelivery(&Action{Method: http.MethodPost, URL: srv.URL}))

	for i := 0; i < 2; i++ {
		select {
		case <-reqs:
		case <-time.After(5 * time.Second):
			t.Fatalf("only %d webhook requests arrived", i)
		}
	}
	eventually(t, 1, func() float64 { return outcome(x, outcomeDelivered) })
	if got := promtest.ToFloat64(x.retries); got != 1 {
		t.Errorf("retries = %v, want 1", got)
	}
	if got := outcome(x, outcomeFailed); got != 0 {
		t.Errorf("failed = %v, want 0", got)
	}
}

// TestWebhook4xxIsPermanent: a 4xx response is a permanent refusal — no
// retry, counted as failed.
func TestWebhook4xxIsPermanent(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1, MaxAttempts: 5})
	x.Enqueue(testDelivery(&Action{Method: http.MethodPost, URL: srv.URL}))

	eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
	// Give any (erroneous) retry a chance to fire before asserting call count.
	time.Sleep(20 * time.Millisecond)
	if got := calls.Load(); got != 1 {
		t.Errorf("attempts = %d, want 1 (4xx must not retry)", got)
	}
	if got := promtest.ToFloat64(x.retries); got != 0 {
		t.Errorf("retries = %v, want 0", got)
	}
}

// TestWebhookFailsAfterMaxAttempts: persistent 5xx exhausts the attempt
// budget and is counted as failed, with MaxAttempts-1 retries.
func TestWebhookFailsAfterMaxAttempts(t *testing.T) {
	reqs := make(chan struct{}, 8)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		reqs <- struct{}{}
	}))
	defer srv.Close()

	const maxAttempts = 3
	x := newTestExecutor(t, ExecutorConfig{Workers: 1, MaxAttempts: maxAttempts})
	x.Enqueue(testDelivery(&Action{Method: http.MethodPost, URL: srv.URL}))

	for i := 0; i < maxAttempts; i++ {
		select {
		case <-reqs:
		case <-time.After(5 * time.Second):
			t.Fatalf("only %d of %d attempts arrived", i, maxAttempts)
		}
	}
	eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
	if got := promtest.ToFloat64(x.retries); got != maxAttempts-1 {
		t.Errorf("retries = %v, want %d", got, maxAttempts-1)
	}
}

// recordingForwarder captures the last Forward call and returns a fixed error.
type recordingForwarder struct {
	mu      sync.Mutex
	called  int
	realm   string
	trigger string
	action  json.RawMessage
	event   []byte
	err     error
}

func (f *recordingForwarder) Forward(_ context.Context, realm, trigger string, action json.RawMessage, event []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.called++
	f.realm, f.trigger, f.action, f.event = realm, trigger, action, event
	return f.err
}

// TestForwarderSkippedWhenUnset: a custom action with no Forwarder is skipped
// (the designed v1 default, not a failure).
func TestForwarderSkippedWhenUnset(t *testing.T) {
	x := newTestExecutor(t, ExecutorConfig{Workers: 1})
	x.Enqueue(testDelivery(&Action{Custom: json.RawMessage(`{"amqp_exchange":"e"}`)}))
	eventually(t, 1, func() float64 { return outcome(x, outcomeSkipped) })
	if got := outcome(x, outcomeForwarded); got != 0 {
		t.Errorf("forwarded = %v, want 0", got)
	}
}

// TestForwarderForwards: a custom action with a Forwarder is handed the realm,
// trigger name, raw action, and rendered event, and counted as forwarded.
func TestForwarderForwards(t *testing.T) {
	fwd := &recordingForwarder{}
	x := newTestExecutor(t, ExecutorConfig{Workers: 1, Forwarder: fwd})
	raw := json.RawMessage(`{"amqp_exchange":"events"}`)
	d := testDelivery(&Action{Custom: raw})
	x.Enqueue(d)

	eventually(t, 1, func() float64 { return outcome(x, outcomeForwarded) })
	fwd.mu.Lock()
	defer fwd.mu.Unlock()
	if fwd.realm != "testrealm" || fwd.trigger != "hook" || string(fwd.action) != string(raw) {
		t.Errorf("forward args: realm=%q trigger=%q action=%s", fwd.realm, fwd.trigger, fwd.action)
	}
	wantEvent, _ := json.Marshal(d.Event)
	if string(fwd.event) != string(wantEvent) {
		t.Errorf("forward event = %s, want %s", fwd.event, wantEvent)
	}
}

// TestForwarderError: a Forwarder error counts as failed.
func TestForwarderError(t *testing.T) {
	fwd := &recordingForwarder{err: io.ErrUnexpectedEOF}
	x := newTestExecutor(t, ExecutorConfig{Workers: 1, Forwarder: fwd})
	x.Enqueue(testDelivery(&Action{Custom: json.RawMessage(`{"amqp_exchange":"e"}`)}))
	eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
}

// TestEnqueueAfterCloseDrops: enqueueing after Close drops with a metric and
// never panics on the closed channel.
func TestEnqueueAfterCloseDrops(t *testing.T) {
	x := newTestExecutor(t, ExecutorConfig{Workers: 1})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := x.Close(ctx); err != nil {
		t.Fatalf("close: %v", err)
	}
	x.Enqueue(testDelivery(&Action{Method: http.MethodPost, URL: "https://unused"}))
	if got := outcome(x, outcomeDropped); got != 1 {
		t.Errorf("dropped = %v, want 1", got)
	}
}

// TestQueueFullDrops: when every worker is busy and the bounded queue is full,
// further enqueues drop with a metric (triggers never backpressure ingestion,
// docs/DESIGN.md §1.4).
func TestQueueFullDrops(t *testing.T) {
	started := make(chan struct{}, 8)
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		started <- struct{}{}
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	defer close(release)

	x := newTestExecutor(t, ExecutorConfig{Workers: 1, QueueSize: 1})
	a := &Action{Method: http.MethodPost, URL: srv.URL}

	x.Enqueue(testDelivery(a)) // dequeued by the sole worker, now blocked in the handler
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("worker never started the first delivery")
	}
	x.Enqueue(testDelivery(a)) // fills the 1-slot queue
	x.Enqueue(testDelivery(a)) // queue full -> dropped

	if got := outcome(x, outcomeDropped); got != 1 {
		t.Errorf("dropped = %v, want 1", got)
	}
}

// TestCloseDrainsQueue: Close lets queued deliveries finish before returning.
func TestCloseDrainsQueue(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 2})
	const n = 6
	for range n {
		x.Enqueue(testDelivery(&Action{Method: http.MethodPost, URL: srv.URL}))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := x.Close(ctx); err != nil {
		t.Fatalf("close: %v", err)
	}
	if got := calls.Load(); got != n {
		t.Errorf("delivered %d of %d before close returned", got, n)
	}
	if got := outcome(x, outcomeDelivered); got != n {
		t.Errorf("delivered metric = %v, want %d", got, n)
	}
}

// TestIgnoreSSLErrors: a self-signed TLS endpoint is reachable only when the
// action opts into ignore_ssl_errors.
func TestIgnoreSSLErrors(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	t.Run("trusted when opted in", func(t *testing.T) {
		x := newTestExecutor(t, ExecutorConfig{Workers: 1, MaxAttempts: 1})
		x.Enqueue(testDelivery(&Action{Method: http.MethodPost, URL: srv.URL, IgnoreSSLErrors: true}))
		eventually(t, 1, func() float64 { return outcome(x, outcomeDelivered) })
	})

	t.Run("rejected without opt-in", func(t *testing.T) {
		x := newTestExecutor(t, ExecutorConfig{Workers: 1, MaxAttempts: 1})
		x.Enqueue(testDelivery(&Action{Method: http.MethodPost, URL: srv.URL}))
		eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
	})
}
