// Package forward provides Forwarder implementations that deliver non-HTTP
// trigger actions to downstream systems. This package is an extension seam;
// callers wire it in via the triggers.Forwarder interface.
package forward

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Config is the HTTP bus forwarder's configuration.
type Config struct {
	// URL is the bus endpoint every custom action is posted to. Required.
	URL string
	// Method defaults to POST when empty.
	Method string
	// StaticHeaders are extra request headers, applied after the fixed ones.
	StaticHeaders map[string]string
	// Client sends the requests; nil selects http.DefaultClient.
	Client *http.Client
}

// body is the envelope posted to the bus for every forwarded event.
type body struct {
	Realm   string          `json:"realm"`
	Trigger string          `json:"trigger"`
	Action  json.RawMessage `json:"action"`
	Event   json.RawMessage `json:"event"`
}

// HTTP posts every custom trigger action to a single bus endpoint.
type HTTP struct {
	url    string
	method string
	static map[string]string
	client *http.Client
}

// New builds an HTTP bus forwarder. It returns an error when URL is empty or
// does not parse, or when Method is not a valid HTTP method token.
func New(cfg Config) (*HTTP, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("forward: URL must not be empty")
	}
	// The forwarder is built at boot from [triggers.forward], so an unusable
	// endpoint must fail here rather than surface once per delivery. That
	// means an absolute http(s) URL, not merely something url.Parse
	// tolerates — a bare "nope" parses cleanly as a relative path.
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("forward: invalid URL %q: %w", cfg.URL, err)
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, fmt.Errorf("forward: URL %q must be an absolute http or https URL", cfg.URL)
	}
	m := cfg.Method
	if m == "" {
		m = http.MethodPost
	} else {
		if !validMethod(m) {
			return nil, fmt.Errorf("forward: invalid method %q", m)
		}
	}
	c := cfg.Client
	if c == nil {
		c = http.DefaultClient
	}
	return &HTTP{
		url:    cfg.URL,
		method: m,
		static: cfg.StaticHeaders,
		client: c,
	}, nil
}

// validMethod returns true when s is a token per RFC 7230 §3.2.6 that
// is also a recognised HTTP method. We accept any single-token string
// so that user-configured methods like PUT or PATCH are allowed without
// a whitelist.
func validMethod(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case 'A' <= c && c <= 'Z', 'a' <= c && c <= 'z', '0' <= c && c <= '9':
			continue
		case c == '!', c == '#', c == '$', c == '%', c == '&', c == '\'',
			c == '*', c == '+', c == '-', c == '.', c == '^', c == '_',
			c == '`', c == '|', c == '~':
			continue
		default:
			return false
		}
	}
	return true
}

// Forward implements the triggers.Forwarder contract.
func (h *HTTP) Forward(ctx context.Context, realm, trigger string, action json.RawMessage, event []byte) error {
	// A nil json.RawMessage marshals to null on its own, but a non-nil empty
	// one marshals to nothing at all and fails the whole envelope — so both
	// empty forms are normalised to nil here.
	var a json.RawMessage
	if len(action) > 0 {
		a = action
	}
	var e json.RawMessage
	if len(event) > 0 {
		e = event
	}
	b := body{
		Realm:   realm,
		Trigger: trigger,
		Action:  a,
		Event:   e,
	}
	raw, err := json.Marshal(b)
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
		return fmt.Errorf("forward: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return fmt.Errorf("forward: drain: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("forward: status %d", resp.StatusCode)
	}
	return nil
}
