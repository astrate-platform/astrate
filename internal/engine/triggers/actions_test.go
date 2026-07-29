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

	t.Run("mustache template parsed and not marked unsupported", func(t *testing.T) {
		a, unsupported, err := parseAction([]byte(
			`{"http_url":"https://x","http_method":"post","template":"{{ value }}","template_type":"mustache"}`))
		if err != nil {
			t.Fatal(err)
		}
		if a.URL == "" || a.Template != "{{ value }}" || len(unsupported) != 0 {
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

// TestWebhookMustacheTemplateRendered: a mustache action renders the body
// from realm/device/trigger/event fields instead of the default envelope.
func TestWebhookMustacheTemplateRendered(t *testing.T) {
	var gotBody []byte
	done := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
		close(done)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1})
	d := testDelivery(&Action{
		Method: http.MethodPost, URL: srv.URL,
		Template:     "{{realm}}/{{device_id}}/{{interface}}{{path}}={{value}}",
		TemplateType: "mustache",
	})
	x.Enqueue(d)

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("webhook never received")
	}
	eventually(t, 1, func() float64 { return outcome(x, outcomeDelivered) })

	const want = "testrealm/f0VMRgIBAQAAAAAAAAAAAA/com.ex.Sensors/v=1"
	if string(gotBody) != want {
		t.Errorf("body = %s, want %s", gotBody, want)
	}
}

// TestWebhookMustacheTemplateMalformedFallsBack: an unclosed tag fails to
// render and the delivery falls back to the default JSON envelope rather
// than crashing the dispatcher.
func TestWebhookMustacheTemplateMalformedFallsBack(t *testing.T) {
	var gotBody []byte
	done := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
		close(done)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1})
	d := testDelivery(&Action{
		Method: http.MethodPost, URL: srv.URL,
		Template:     "{{unclosed",
		TemplateType: "mustache",
	})
	x.Enqueue(d)

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("webhook never received")
	}
	eventually(t, 1, func() float64 { return outcome(x, outcomeDelivered) })

	wantBody, _ := json.Marshal(d.Event)
	if string(gotBody) != string(wantBody) {
		t.Errorf("body = %s, want default envelope %s", gotBody, wantBody)
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

// mustCompilePolicy compiles a policy JSON document or fails the test.
func mustCompilePolicy(t *testing.T, def string) *Policy {
	t.Helper()
	p, err := CompilePolicy([]byte(def))
	if err != nil {
		t.Fatalf("CompilePolicy: %v", err)
	}
	return p
}

// TestPolicyRetriesClientError: a policy with client_error=retry retries a
// 404 exactly retry_times times then fails — the opposite of the default rule
// (4xx = permanent), proving the policy is in charge.
func TestPolicyRetriesClientError(t *testing.T) {
	policy := mustCompilePolicy(t, `{
		"name":"retry-4xx",
		"error_handlers":[{"on":"client_error","strategy":"retry"}],
		"retry_times":2,
		"maximum_capacity":1
	}`)

	var calls atomic.Int32
	reqs := make(chan struct{}, 8)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusNotFound)
		reqs <- struct{}{}
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1})
	d := testDelivery(&Action{Method: http.MethodPost, URL: srv.URL})
	d.Trigger.AttachPolicy(policy)
	x.Enqueue(d)

	for i := 0; i < 3; i++ { // 1 initial + 2 retries = 3 attempts
		select {
		case <-reqs:
		case <-time.After(5 * time.Second):
			t.Fatalf("only %d of 3 attempts arrived", i)
		}
	}
	eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
	if got := promtest.ToFloat64(x.retries); got != 2 {
		t.Errorf("retries = %v, want 2", got)
	}
	if got := calls.Load(); got != 3 {
		t.Errorf("attempts = %d, want 3", got)
	}
}

// TestPolicyDiscardsServerError: a policy with server_error=discard gives up
// on a 500 after exactly one attempt.
func TestPolicyDiscardsServerError(t *testing.T) {
	policy := mustCompilePolicy(t, `{
		"name":"discard-5xx",
		"error_handlers":[{"on":"server_error","strategy":"discard"}],
		"maximum_capacity":1
	}`)

	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1, MaxAttempts: 5})
	d := testDelivery(&Action{Method: http.MethodPost, URL: srv.URL})
	d.Trigger.AttachPolicy(policy)
	x.Enqueue(d)

	eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
	time.Sleep(20 * time.Millisecond)
	if got := calls.Load(); got != 1 {
		t.Errorf("attempts = %d, want 1 (policy says discard)", got)
	}
	if got := promtest.ToFloat64(x.retries); got != 0 {
		t.Errorf("retries = %v, want 0", got)
	}
}

// TestPolicyExplicitStatusRetries503Discards500: an explicit status list
// {503} retries 503 but discards 500.
func TestPolicyExplicitStatusRetries503Discards500(t *testing.T) {
	policy := mustCompilePolicy(t, `{
		"name":"retry-503-only",
		"error_handlers":[{"on":[503],"strategy":"retry"}],
		"retry_times":1,
		"maximum_capacity":1
	}`)

	var calls atomic.Int32
	reqs := make(chan struct{}, 8)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := calls.Add(1)
		if n == 1 {
			w.WriteHeader(http.StatusServiceUnavailable) // first: 503 → retry
		} else {
			w.WriteHeader(http.StatusInternalServerError) // second: 500 → discard
		}
		reqs <- struct{}{}
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1, MaxAttempts: 5})
	d := testDelivery(&Action{Method: http.MethodPost, URL: srv.URL})
	d.Trigger.AttachPolicy(policy)
	x.Enqueue(d)

	for i := 0; i < 2; i++ {
		select {
		case <-reqs:
		case <-time.After(5 * time.Second):
			t.Fatalf("only %d of 2 attempts arrived", i)
		}
	}
	eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
	if got := promtest.ToFloat64(x.retries); got != 1 {
		t.Errorf("retries = %v, want 1", got)
	}
	if got := calls.Load(); got != 2 {
		t.Errorf("attempts = %d, want 2", got)
	}
}

// TestPolicyTransportFailureFollowsServerError: a transport failure (closed
// port) is treated as a server error by the policy — it matches server_error
// handlers and is NOT matched by a policy whose only handler is an explicit
// status-code list like [500].
func TestPolicyTransportFailureFollowsServerError(t *testing.T) {
	policy := mustCompilePolicy(t, `{
		"name":"explicit-500-only",
		"error_handlers":[{"on":[500],"strategy":"retry"}],
		"retry_times":2,
		"maximum_capacity":1
	}`)

	// Point at a closed port to force a transport error.
	closedSrv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	closedSrv.Close() // close immediately so port is unreachable

	x := newTestExecutor(t, ExecutorConfig{Workers: 1, MaxAttempts: 5})
	d := testDelivery(&Action{Method: http.MethodPost, URL: closedSrv.URL})
	d.Trigger.AttachPolicy(policy)
	x.Enqueue(d)

	eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
	// Transport error (status 0) does not match the explicit [500] handler,
	// so no handler claims it → StrategyDiscard → exactly one attempt.
	if got := promtest.ToFloat64(x.retries); got != 0 {
		t.Errorf("retries = %v, want 0", got)
	}
}

// TestNoPolicyFollowsDefaultBehaviour: a trigger with no attached policy
// behaves exactly as the pre-policy hardcoded rules.
func TestNoPolicyFollowsDefaultBehaviour(t *testing.T) {
	reqs := make(chan struct{}, 8)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway) // 502 → retry under default
		reqs <- struct{}{}
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1, MaxAttempts: 3})
	x.Enqueue(testDelivery(&Action{Method: http.MethodPost, URL: srv.URL}))

	for i := 0; i < 3; i++ {
		select {
		case <-reqs:
		case <-time.After(5 * time.Second):
			t.Fatalf("only %d of 3 attempts arrived", i)
		}
	}
	eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
	if got := promtest.ToFloat64(x.retries); got != 2 {
		t.Errorf("retries = %v, want 2", got)
	}
}

// TestPolicyTTExpiresDelivery: a policy with event_ttl=1s causes a queued
// delivery to be dropped as expired before any attempt when the worker is
// stalled past the TTL.
func TestPolicyTTExpiresDelivery(t *testing.T) {
	policy := mustCompilePolicy(t, `{
		"name":"ttl-1s",
		"error_handlers":[{"on":"any_error","strategy":"retry"}],
		"retry_times":1,
		"maximum_capacity":2,
		"event_ttl":1
	}`)

	var calls atomic.Int32
	started := make(chan struct{})
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		started <- struct{}{}
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1})

	// First delivery stalls the sole worker.
	d1 := testDelivery(&Action{Method: http.MethodPost, URL: srv.URL})
	d1.Trigger.AttachPolicy(policy)
	x.Enqueue(d1)
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("worker never started the first delivery")
	}

	// Second delivery sits in the queue past the 1-second TTL.
	d2 := testDelivery(&Action{Method: http.MethodPost, URL: srv.URL})
	d2.Trigger.AttachPolicy(policy)
	x.Enqueue(d2)

	// Wait for the TTL to expire.
	time.Sleep(1100 * time.Millisecond)

	// Release the stall so the worker can process the second delivery.
	close(release)

	eventually(t, 1, func() float64 { return outcome(x, outcomeExpired) })
	if got := outcome(x, outcomeDelivered); got != 1 {
		t.Errorf("delivered = %v, want 1", got)
	}
	if got := calls.Load(); got != 1 {
		t.Errorf("attempts = %d, want 1 (expired delivery must not hit endpoint)", got)
	}
}

// TestPolicyCapacityDrops: with maximum_capacity=1, a second Enqueue for the
// same policy is dropped while a delivery for a trigger with no policy still
// queues.
func TestPolicyCapacityDrops(t *testing.T) {
	policy := mustCompilePolicy(t, `{
		"name":"cap-1",
		"error_handlers":[{"on":"any_error","strategy":"discard"}],
		"maximum_capacity":1
	}`)

	var reqCount atomic.Int32
	started := make(chan struct{}, 2)
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		reqCount.Add(1)
		started <- struct{}{}
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1})

	// First delivery occupies the sole worker.
	d1 := testDelivery(&Action{Method: http.MethodPost, URL: srv.URL})
	d1.Trigger.AttachPolicy(policy)
	x.Enqueue(d1)
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("worker never started the first delivery")
	}

	// Second delivery for the same policy must be dropped (capacity full).
	d2 := testDelivery(&Action{Method: http.MethodPost, URL: srv.URL})
	d2.Trigger.AttachPolicy(policy)
	x.Enqueue(d2)

	// A delivery with no policy still queues (unaffected).
	d3 := testDelivery(&Action{Method: http.MethodPost, URL: srv.URL})
	x.Enqueue(d3)

	if got := outcome(x, outcomeDropped); got != 1 {
		t.Errorf("dropped = %v, want 1", got)
	}

	close(release)

	// Wait for both in-flight deliveries to complete.
	deadline := time.After(5 * time.Second)
	for reqCount.Load() < 2 {
		select {
		case <-deadline:
			t.Fatalf("only %d requests arrived, want 2", reqCount.Load())
		case <-time.After(time.Millisecond):
		}
	}

	if got := outcome(x, outcomeDropped); got != 1 {
		t.Errorf("dropped = %v, want 1", got)
	}
}

// TestPolicyCapacityResetsAfterDelivery: after the in-flight delivery
// completes, the per-policy counter returns to zero so a subsequent Enqueue
// for that policy succeeds.
func TestPolicyCapacityResetsAfterDelivery(t *testing.T) {
	policy := mustCompilePolicy(t, `{
		"name":"cap-reset",
		"error_handlers":[{"on":"any_error","strategy":"discard"}],
		"maximum_capacity":1
	}`)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1})

	// First delivery succeeds immediately.
	d1 := testDelivery(&Action{Method: http.MethodPost, URL: srv.URL})
	d1.Trigger.AttachPolicy(policy)
	x.Enqueue(d1)
	eventually(t, 1, func() float64 { return outcome(x, outcomeDelivered) })

	// Capacity counter is back to zero; a second Enqueue must succeed.
	d2 := testDelivery(&Action{Method: http.MethodPost, URL: srv.URL})
	d2.Trigger.AttachPolicy(policy)
	x.Enqueue(d2)
	eventually(t, 2, func() float64 { return outcome(x, outcomeDelivered) })
}

// TestEnqueueRacingCloseDoesNotPanic pins the shutdown race: Enqueue checks
// x.closing and then sends, and without the send lock a caller sitting
// between those two steps sends on a channel Close has already closed, which
// panics and takes the process down. Run under -race with many senders so a
// regression surfaces rather than lurking.
func TestEnqueueRacingCloseDoesNotPanic(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 2, QueueSize: 4})
	action := &Action{Method: http.MethodPost, URL: srv.URL}

	var wg sync.WaitGroup
	start := make(chan struct{})
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for range 50 {
				x.Enqueue(testDelivery(action))
			}
		}()
	}
	close(start)

	// Close while the senders are mid-flight; every Enqueue that loses the
	// race must drop its delivery, not panic.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := x.Close(ctx); err != nil {
		t.Fatalf("close during concurrent enqueue: %v", err)
	}
	wg.Wait()
}
