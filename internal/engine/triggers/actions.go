package triggers

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Action is a trigger's parsed delivery action (docs/ROADMAP.md §7.2 file
// 6.12): an HTTP webhook (upstream "http_url"+"http_method", or the legacy
// "http_post_url"), or a custom action routed to the Forwarder extension
// point (upstream AMQP actions land there, docs/DESIGN.md §1.1).
type Action struct {
	// Method is the upper-cased HTTP method.
	Method string
	// URL is the webhook endpoint.
	URL string
	// StaticHeaders are the action's extra request headers.
	StaticHeaders map[string]string
	// IgnoreSSLErrors disables server-certificate verification for this
	// action's requests (upstream "ignore_ssl_errors").
	IgnoreSSLErrors bool
	// Custom is the raw action object of a non-HTTP action, delivered
	// through the Forwarder extension point; nil for HTTP actions.
	Custom json.RawMessage
}

// httpAction is the upstream HTTP action JSON shape.
type httpAction struct {
	HTTPURL           string            `json:"http_url"`
	HTTPMethod        string            `json:"http_method"`
	HTTPPostURL       string            `json:"http_post_url"` // pre-1.1 legacy: implies POST
	HTTPStaticHeaders map[string]string `json:"http_static_headers"`
	IgnoreSSLErrors   bool              `json:"ignore_ssl_errors"`
	Template          string            `json:"template"`
	TemplateType      string            `json:"template_type"`
}

// httpMethods is the upstream-accepted method set (lowercase on the wire).
var httpMethods = map[string]bool{
	"delete": true, "get": true, "head": true, "options": true,
	"patch": true, "post": true, "put": true,
}

// parseAction validates one action object. It returns the parsed action
// plus the list of accepted-but-not-evaluated features (Mustache payload
// templates: the default JSON envelope is sent instead).
func parseAction(raw json.RawMessage) (*Action, []string, error) {
	if len(raw) == 0 {
		return nil, nil, fmt.Errorf("missing action")
	}
	var h httpAction
	if err := json.Unmarshal(raw, &h); err != nil {
		return nil, nil, fmt.Errorf("action does not parse: %w", err)
	}

	var unsupported []string
	if h.Template != "" || h.TemplateType != "" {
		unsupported = append(unsupported, "mustache payload template (default JSON envelope sent)")
	}

	switch {
	case h.HTTPPostURL != "":
		return &Action{
			Method: http.MethodPost, URL: h.HTTPPostURL,
			StaticHeaders: h.HTTPStaticHeaders, IgnoreSSLErrors: h.IgnoreSSLErrors,
		}, unsupported, nil
	case h.HTTPURL != "":
		if !httpMethods[h.HTTPMethod] {
			return nil, nil, fmt.Errorf("unsupported http_method %q", h.HTTPMethod)
		}
		return &Action{
			// Methods arrive lowercase on the wire (upstream convention) but
			// HTTP methods are case-sensitive uppercase tokens.
			Method: strings.ToUpper(h.HTTPMethod), URL: h.HTTPURL,
			StaticHeaders: h.HTTPStaticHeaders, IgnoreSSLErrors: h.IgnoreSSLErrors,
		}, unsupported, nil
	default:
		// Not an HTTP action (e.g. upstream AMQP): keep it verbatim for the
		// Forwarder extension point.
		return &Action{Custom: raw}, unsupported, nil
	}
}

// Forwarder is the extension point for non-HTTP trigger actions
// (docs/DESIGN.md §1.1: "AMQP action replaced by optional NATS/HTTP
// forwarding"). The executor hands it every matched event whose action is
// not an HTTP webhook.
//
// TODO(extension): provide NATS and HTTP-bus Forwarder implementations; the
// default (nil) logs the event and counts it as skipped, which is the
// designed v1 behaviour, not a gap in this code path.
type Forwarder interface {
	// Forward delivers one rendered event for a custom action.
	Forward(ctx context.Context, realm, trigger string, action json.RawMessage, event []byte) error
}

// Delivery is one matched (trigger, event) pair queued for execution.
type Delivery struct {
	// Realm is the tenant (rides in the Astarte-Realm header).
	Realm string
	// Trigger is the matched trigger.
	Trigger *Trigger
	// Event is the rendered envelope.
	Event SimpleEvent

	// enqueuedAt is stamped by Enqueue and used for TTL checks.
	enqueuedAt time.Time
}

// Executor defaults.
const (
	// DefaultWorkers is the default delivery worker count.
	DefaultWorkers = 4
	// DefaultQueueSize is the default delivery queue capacity.
	DefaultQueueSize = 256
	// DefaultMaxAttempts bounds delivery attempts per event.
	DefaultMaxAttempts = 5
	// DefaultBackoffStart is the first retry delay.
	DefaultBackoffStart = 250 * time.Millisecond
	// DefaultBackoffCap bounds the exponential retry delay.
	DefaultBackoffCap = 5 * time.Second
	// DefaultRequestTimeout bounds one webhook request.
	DefaultRequestTimeout = 10 * time.Second
)

// Delivery outcome labels.
const (
	outcomeDelivered = "delivered"
	outcomeFailed    = "failed"
	outcomeDropped   = "dropped"
	outcomeExpired   = "expired"
	outcomeForwarded = "forwarded"
	outcomeSkipped   = "skipped"
)

// ExecutorConfig carries the executor's knobs; zero values select the
// defaults above.
type ExecutorConfig struct {
	// Workers is the number of delivery goroutines.
	Workers int
	// QueueSize is the bounded delivery queue capacity; a full queue drops
	// the event with a metric (triggers must never backpressure ingestion,
	// docs/DESIGN.md §1.4).
	QueueSize int
	// MaxAttempts bounds attempts per delivery (1 = no retries).
	MaxAttempts int
	// BackoffStart and BackoffCap shape the exponential retry delay.
	BackoffStart time.Duration
	BackoffCap   time.Duration
	// RequestTimeout bounds each webhook request.
	RequestTimeout time.Duration
	// Forwarder handles non-HTTP actions; nil logs and skips them.
	Forwarder Forwarder
	// Registerer receives the executor's collectors; nil leaves them
	// unregistered (they still work, which tests rely on).
	Registerer prometheus.Registerer
	// Logger receives delivery logs (default slog.Default()).
	Logger *slog.Logger
}

// Executor runs trigger actions asynchronously: a bounded queue feeding
// worker goroutines, HTTP webhooks with exponential-backoff retry, and
// delivery-outcome metrics (docs/ROADMAP.md §7.2 file 6.12).
type Executor struct {
	cfg      ExecutorConfig
	log      *slog.Logger
	ch       chan Delivery
	wg       sync.WaitGroup
	closing  chan struct{}
	stopOnce sync.Once

	client         *http.Client
	insecureClient *http.Client

	outcomes *prometheus.CounterVec
	retries  prometheus.Counter

	policyMu       sync.Mutex
	policyInflight map[string]int

	// sendMu serialises senders against the channel close in Close.
	sendMu sync.RWMutex
}

// NewExecutor builds and starts an executor.
func NewExecutor(cfg ExecutorConfig) *Executor {
	if cfg.Workers <= 0 {
		cfg.Workers = DefaultWorkers
	}
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = DefaultQueueSize
	}
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = DefaultMaxAttempts
	}
	if cfg.BackoffStart <= 0 {
		cfg.BackoffStart = DefaultBackoffStart
	}
	if cfg.BackoffCap <= 0 {
		cfg.BackoffCap = DefaultBackoffCap
	}
	if cfg.RequestTimeout <= 0 {
		cfg.RequestTimeout = DefaultRequestTimeout
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	x := &Executor{
		cfg:     cfg,
		log:     cfg.Logger.With("component", "triggers"),
		ch:      make(chan Delivery, cfg.QueueSize),
		closing: make(chan struct{}),
		client:  &http.Client{Timeout: cfg.RequestTimeout},
		insecureClient: &http.Client{
			Timeout: cfg.RequestTimeout,
			Transport: &http.Transport{
				// Explicitly requested per action via ignore_ssl_errors
				// (upstream parity).
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402 -- opt-in upstream-parity action flag
			},
		},
		outcomes: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "astrate_engine_trigger_deliveries_total",
			Help: "Trigger action deliveries by outcome (docs/DESIGN.md §5.2).",
		}, []string{"outcome"}),
		retries: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_engine_trigger_retries_total",
			Help: "Failed webhook attempts that were retried.",
		}),
		policyInflight: make(map[string]int),
	}
	for _, o := range []string{outcomeDelivered, outcomeFailed, outcomeDropped, outcomeExpired, outcomeForwarded, outcomeSkipped} {
		x.outcomes.WithLabelValues(o)
	}
	if cfg.Registerer != nil {
		cfg.Registerer.MustRegister(x.outcomes, x.retries)
	}

	x.wg.Add(cfg.Workers)
	for range cfg.Workers {
		go x.worker()
	}
	return x
}

// Enqueue queues one delivery without blocking; a full queue drops it with
// a metric and a log line.
//
// The read lock spans the closing check and the send: without it a caller
// that passed the check could still be holding an unsent delivery when Close
// closes the channel underneath it, and the send would panic.
func (x *Executor) Enqueue(d Delivery) {
	x.sendMu.RLock()
	defer x.sendMu.RUnlock()

	select {
	case <-x.closing:
		x.outcomes.WithLabelValues(outcomeDropped).Inc()
		return
	default:
	}

	d.enqueuedAt = time.Now()

	if pol := d.Trigger.Policy(); pol != nil {
		x.policyMu.Lock()
		if x.policyInflight[pol.Name] >= pol.MaximumCapacity {
			x.policyMu.Unlock()
			x.outcomes.WithLabelValues(outcomeDropped).Inc()
			x.log.Warn("trigger delivery policy capacity full; event dropped",
				"realm", d.Realm, "trigger", d.Trigger.Name, "policy", pol.Name)
			return
		}
		x.policyInflight[pol.Name]++
		x.policyMu.Unlock()
	}

	select {
	case x.ch <- d:
	default:
		if pol := d.Trigger.Policy(); pol != nil {
			x.policyMu.Lock()
			x.policyInflight[pol.Name]--
			x.policyMu.Unlock()
		}
		x.outcomes.WithLabelValues(outcomeDropped).Inc()
		x.log.Warn("trigger delivery queue full; event dropped",
			"realm", d.Realm, "trigger", d.Trigger.Name, "device", d.Event.DeviceID)
	}
}

// Close stops accepting deliveries, lets the workers drain the queue, and
// waits for them bounded by ctx. In-flight retry sleeps abort immediately.
func (x *Executor) Close(ctx context.Context) error {
	x.stopOnce.Do(func() {
		close(x.closing)
		// Take the write lock before closing: it can only be acquired once
		// every Enqueue in flight has finished its send, and every later one
		// sees the closed x.closing and drops instead.
		x.sendMu.Lock()
		close(x.ch)
		x.sendMu.Unlock()
	})
	done := make(chan struct{})
	go func() {
		x.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("triggers: executor close interrupted: %w", ctx.Err())
	}
}

// worker drains the delivery queue.
func (x *Executor) worker() {
	defer x.wg.Done()
	for d := range x.ch {
		x.deliver(d)
	}
}

// deliver executes one delivery end to end.
func (x *Executor) deliver(d Delivery) {
	if pol := d.Trigger.Policy(); pol != nil {
		// Balances the in-flight count Enqueue took for this policy, on
		// every path out of deliver — expiry, marshal failure, delivery.
		defer func() {
			x.policyMu.Lock()
			x.policyInflight[pol.Name]--
			x.policyMu.Unlock()
		}()

		if pol.EventTTL > 0 && time.Since(d.enqueuedAt) > pol.EventTTL {
			x.outcomes.WithLabelValues(outcomeExpired).Inc()
			x.log.Info("trigger delivery expired",
				"realm", d.Realm, "trigger", d.Trigger.Name, "policy", pol.Name,
				"age", time.Since(d.enqueuedAt), "event_ttl", pol.EventTTL)
			return
		}
	}

	body, err := json.Marshal(d.Event)
	if err != nil {
		x.outcomes.WithLabelValues(outcomeFailed).Inc()
		x.log.Error("trigger event does not marshal",
			"realm", d.Realm, "trigger", d.Trigger.Name, "err", err)
		return
	}

	if d.Trigger.Action.Custom != nil {
		x.forward(d, body)
		return
	}
	x.webhook(d, body)
}

// forward hands a non-HTTP action to the Forwarder extension point.
func (x *Executor) forward(d Delivery, body []byte) {
	if x.cfg.Forwarder == nil {
		x.outcomes.WithLabelValues(outcomeSkipped).Inc()
		x.log.Info("no forwarder configured; custom trigger action skipped",
			"realm", d.Realm, "trigger", d.Trigger.Name)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), x.cfg.RequestTimeout)
	defer cancel()
	if err := x.cfg.Forwarder.Forward(ctx, d.Realm, d.Trigger.Name, d.Trigger.Action.Custom, body); err != nil {
		x.outcomes.WithLabelValues(outcomeFailed).Inc()
		x.log.Warn("custom trigger action failed",
			"realm", d.Realm, "trigger", d.Trigger.Name, "err", err)
		return
	}
	x.outcomes.WithLabelValues(outcomeForwarded).Inc()
}

// webhook POSTs (or whatever the action's method is) the event JSON and
// handles the delivery result.
//
// When the trigger has no attached policy (Policy() == nil), behaviour is the
// hardcoded default: 2xx/3xx are delivered, 4xx are permanently failed, 5xx
// and transport errors retry with exponential backoff up to cfg.MaxAttempts.
//
// When a policy is attached, it governs every retry/discard decision: after
// each attempt the HTTP status (or StatusTransport for transport errors) is
// passed to policy.Decide. StrategyRetry retries with the same exponential
// backoff; StrategyDiscard ends the delivery as failed. The attempt cap is
// policy.RetryTimes+1 (RetryTimes counts retries, not attempts). Success
// remains status in 200..399, checked before consulting the policy.
//
// A request that never produced a response is passed as StatusTransport and
// treated as a server error: upstream's policy schema speaks only in status
// codes, so this is Astrate's reading, recorded in docs/COMPATIBILITY.md.
// Every decision logs the reason that produced it.
func (x *Executor) webhook(d Delivery, body []byte) {
	policy := d.Trigger.Policy()
	backoff := x.cfg.BackoffStart
	maxAttempts := x.cfg.MaxAttempts
	if policy != nil {
		maxAttempts = policy.RetryTimes + 1
	}
	for attempt := 1; ; attempt++ {
		status, err := x.attempt(d, body)
		if err == nil && status < 400 {
			x.outcomes.WithLabelValues(outcomeDelivered).Inc()
			return
		}
		if policy != nil {
			effectiveStatus := status
			if err != nil {
				effectiveStatus = StatusTransport
			}
			dec := policy.Decide(effectiveStatus)
			if dec.Strategy == StrategyDiscard {
				x.outcomes.WithLabelValues(outcomeFailed).Inc()
				x.log.Warn("webhook delivery discarded by policy",
					"realm", d.Realm, "trigger", d.Trigger.Name,
					"status", status, "policy", policy.Name, "reason", dec.Reason)
				return
			}
			if attempt >= maxAttempts {
				x.outcomes.WithLabelValues(outcomeFailed).Inc()
				x.log.Warn("webhook delivery failed after final attempt",
					"realm", d.Realm, "trigger", d.Trigger.Name, "attempts", attempt,
					"status", status, "err", err,
					"policy", policy.Name, "reason", dec.Reason)
				return
			}
			x.log.Info("webhook delivery retrying",
				"realm", d.Realm, "trigger", d.Trigger.Name, "attempt", attempt,
				"status", status, "policy", policy.Name, "reason", dec.Reason)
			x.retries.Inc()
			select {
			case <-x.closing:
				x.outcomes.WithLabelValues(outcomeFailed).Inc()
				x.log.Warn("webhook delivery abandoned at shutdown",
					"realm", d.Realm, "trigger", d.Trigger.Name, "attempts", attempt,
					"policy", policy.Name, "reason", dec.Reason)
				return
			case <-time.After(backoff):
			}
			backoff = min(backoff*2, x.cfg.BackoffCap)
			continue
		}
		if err == nil && status < 500 {
			x.outcomes.WithLabelValues(outcomeFailed).Inc()
			x.log.Warn("webhook permanently refused",
				"realm", d.Realm, "trigger", d.Trigger.Name, "status", status)
			return
		}
		if attempt >= maxAttempts {
			x.outcomes.WithLabelValues(outcomeFailed).Inc()
			x.log.Warn("webhook delivery failed after final attempt",
				"realm", d.Realm, "trigger", d.Trigger.Name, "attempts", attempt,
				"status", status, "err", err)
			return
		}
		x.retries.Inc()
		select {
		case <-x.closing:
			x.outcomes.WithLabelValues(outcomeFailed).Inc()
			x.log.Warn("webhook delivery abandoned at shutdown",
				"realm", d.Realm, "trigger", d.Trigger.Name, "attempts", attempt)
			return
		case <-time.After(backoff):
		}
		backoff = min(backoff*2, x.cfg.BackoffCap)
	}
}

// attempt performs one webhook request, returning the response status (0 on
// transport errors).
func (x *Executor) attempt(d Delivery, body []byte) (int, error) {
	a := d.Trigger.Action
	ctx, cancel := context.WithTimeout(context.Background(), x.cfg.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, a.Method, a.URL, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Astarte-Realm", d.Realm)
	for k, v := range a.StaticHeaders {
		req.Header.Set(k, v)
	}

	client := x.client
	if a.IgnoreSSLErrors {
		client = x.insecureClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	// Drain so the connection is reusable; the body itself is irrelevant.
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
	_ = resp.Body.Close()
	return resp.StatusCode, nil
}
