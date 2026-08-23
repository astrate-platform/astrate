package blocks

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
)

// HTTP block types (issue #82, astarte_flow parity): http_source polls GET
// URLs round-robin, http_sink POSTs message payloads to a URL. Both use
// net/http only and talk to whatever endpoints the operator configures.
const (
	TypeHTTPSource = "http_source"
	TypeHTTPSink   = "http_sink"
)

// positiveMillis reads key from config as a millisecond duration, applying
// def when the key is absent or null. Shared by both HTTP blocks so their
// validation messages cannot drift.
func positiveMillis(config map[string]any, key string, def time.Duration) (time.Duration, error) {
	v, ok := config[key]
	if !ok || v == nil {
		return def, nil
	}
	n, err := numAsInt64(v)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("%s must be positive", key)
	}
	return time.Duration(n) * time.Millisecond, nil
}

// mediaType strips parameters ("; charset=...") and surrounding spaces from
// a Content-Type header value; "" stays "".
func mediaType(contentType string) string {
	if i := strings.IndexByte(contentType, ';'); i >= 0 {
		contentType = contentType[:i]
	}
	return strings.TrimSpace(contentType)
}

// HTTPSource constructs a Source that polls configured GET URLs round-robin,
// emitting one binary message per response.
//
// Config keys:
//   - urls (required string array): at least one URL to poll; empty-string
//     entries are rejected
//   - interval_ms (int, default 1000): wait before every Emit (including the
//     first); must be positive
//   - timeout_ms (int, default 5000): per-request timeout; must be positive
func HTTPSource(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
	cfg, err := parseHTTPSourceConfig(config)
	if err != nil {
		return nil, fmt.Errorf("http_source: %w", err)
	}
	return &httpSource{
		name:     name,
		urls:     cfg.urls,
		interval: cfg.interval,
		timeout:  cfg.timeout,
		client:   &http.Client{},
	}, nil
}

type httpSourceConfig struct {
	urls     []string
	interval time.Duration
	timeout  time.Duration
}

func parseHTTPSourceConfig(config map[string]any) (httpSourceConfig, error) {
	cfg := httpSourceConfig{interval: time.Second, timeout: 5 * time.Second}

	urls, err := parseURLList(config["urls"])
	if err != nil {
		return cfg, err
	}
	for _, u := range urls {
		if u == "" {
			return cfg, fmt.Errorf("at least one url is required")
		}
	}
	if len(urls) == 0 {
		return cfg, fmt.Errorf("at least one url is required")
	}
	cfg.urls = urls

	if cfg.interval, err = positiveMillis(config, "interval_ms", cfg.interval); err != nil {
		return cfg, err
	}
	if cfg.timeout, err = positiveMillis(config, "timeout_ms", cfg.timeout); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// parseURLList coerces the urls config entry (a []string or a JSON-decoded
// []any of strings) into a []string.
func parseURLList(v any) ([]string, error) {
	switch list := v.(type) {
	case nil:
		return nil, nil
	case []string:
		out := make([]string, len(list))
		copy(out, list)
		return out, nil
	case []any:
		out := make([]string, 0, len(list))
		for _, item := range list {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("urls must contain only strings")
			}
			out = append(out, s)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("urls must be a list of strings")
	}
}

type httpSource struct {
	name     string
	urls     []string
	interval time.Duration
	timeout  time.Duration
	client   *http.Client
	next     atomic.Int64
}

var (
	_ flow.Block  = (*httpSource)(nil)
	_ flow.Source = (*httpSource)(nil)
)

func (s *httpSource) Name() string { return s.name }

// Process implements flow.Block for the non-pump path: it performs one fetch
// immediately without waiting for the interval.
func (s *httpSource) Process(_ *flow.Message) ([]*flow.Message, error) {
	msg, err := s.fetch(context.TODO())
	if err != nil {
		return nil, err
	}
	return []*flow.Message{msg}, nil
}

// Emit implements flow.Source: it waits one interval (including on the first
// call), then polls the next URL in round-robin order.
func (s *httpSource) Emit(ctx context.Context) ([]*flow.Message, error) {
	timer := time.NewTimer(s.interval)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
	}
	msg, err := s.fetch(ctx)
	if err != nil {
		return nil, err
	}
	return []*flow.Message{msg}, nil
}

// fetch GETs the next URL in round-robin order and converts the response to
// one binary message keyed by the polled URL.
func (s *httpSource) fetch(ctx context.Context) (*flow.Message, error) {
	url := s.urls[int(s.next.Add(1)-1)%len(s.urls)]

	reqCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http_source: GET %s: %w", url, err)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http_source: GET %s: %w", url, err)
	}
	body, readErr := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if readErr != nil {
		return nil, fmt.Errorf("http_source: GET %s: %w", url, readErr)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http_source: GET %s: status %d", url, resp.StatusCode)
	}
	return &flow.Message{
		Key:       url,
		Type:      flow.TypeBinary,
		Subtype:   mediaType(resp.Header.Get("Content-Type")),
		Data:      body,
		Timestamp: time.Now().UnixMicro(),
	}, nil
}

// HTTPSink constructs a Sink that sends each message payload to a URL.
//
// Config keys:
//   - url (required string): target of every request
//   - method (string, default "POST"): any non-empty token allowed
//   - timeout_ms (int, default 5000): per-request timeout; must be positive
//   - headers (object string→string): set verbatim on each request
//
// Payload mapping: binary messages send their bytes with Subtype (or
// application/octet-stream) as Content-Type; strings send text/plain;
// everything else is JSON-encoded with application/json. Status >= 400 fails
// the message; 2xx/3xx succeed.
func HTTPSink(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
	cfg, client, err := parseHTTPSinkConfig(config)
	if err != nil {
		return nil, fmt.Errorf("http_sink: %w", err)
	}
	return flow.NewSinkBlock(name, func(msg *flow.Message) error {
		if msg == nil {
			return nil
		}
		return deliverHTTP(cfg, client, msg)
	}), nil
}

type httpSinkConfig struct {
	url     string
	method  string
	timeout time.Duration
	headers map[string]string
}

func parseHTTPSinkConfig(config map[string]any) (*httpSinkConfig, *http.Client, error) {
	cfg := &httpSinkConfig{
		method:  http.MethodPost,
		timeout: 5 * time.Second,
	}

	cfg.url = stringConfig(config, "url", "")
	if cfg.url == "" {
		return nil, nil, fmt.Errorf("url is required")
	}
	if m := stringConfig(config, "method", ""); m != "" {
		cfg.method = m
	}

	var err error
	if cfg.timeout, err = positiveMillis(config, "timeout_ms", cfg.timeout); err != nil {
		return nil, nil, err
	}

	raw, ok := config["headers"]
	if ok && raw != nil {
		headers, err := headerMap(raw)
		if err != nil {
			return nil, nil, err
		}
		cfg.headers = headers
	}
	return cfg, &http.Client{}, nil
}

// headerMap coerces the headers config entry into a string→string map.
func headerMap(v any) (map[string]string, error) {
	switch m := v.(type) {
	case map[string]string:
		out := make(map[string]string, len(m))
		for k, val := range m {
			out[k] = val
		}
		return out, nil
	case map[string]any:
		out := make(map[string]string, len(m))
		for k, val := range m {
			s, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("headers values must be strings")
			}
			out[k] = s
		}
		return out, nil
	default:
		return nil, fmt.Errorf("headers must be an object of string→string")
	}
}

// sinkPayload maps a message to request body bytes and a derived
// Content-Type (used only when the operator did not configure one).
func sinkPayload(msg *flow.Message) ([]byte, string, error) {
	switch msg.Type {
	case flow.TypeBinary:
		bs, ok := msg.Data.([]byte)
		if ok {
			ct := msg.Subtype
			if ct == "" {
				ct = "application/octet-stream"
			}
			return bs, ct, nil
		}
	case flow.TypeString:
		if s, ok := msg.Data.(string); ok {
			return []byte(s), "text/plain; charset=utf-8", nil
		}
	}
	body, err := toJSONBytes(msg)
	if err != nil {
		return nil, "", fmt.Errorf("http_sink: encode payload: %w", err)
	}
	return body, "application/json", nil
}

// deliverHTTP sends one message to the configured endpoint.
func deliverHTTP(cfg *httpSinkConfig, client *http.Client, msg *flow.Message) error {
	body, contentType, err := sinkPayload(msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, cfg.method, cfg.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http_sink: %s %s: %w", cfg.method, cfg.url, err)
	}
	for k, v := range cfg.headers {
		req.Header.Set(k, v)
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http_sink: %s %s: %w", cfg.method, cfg.url, err)
	}
	_, _ = io.Copy(io.Discard, resp.Body) // drain so the connection can be reused
	_ = resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("http_sink: %s %s: status %d", cfg.method, cfg.url, resp.StatusCode)
	}
	return nil
}
