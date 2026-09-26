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
	"strings"
)

const (
	// maxErrorBody is how much of a non-2xx response body is quoted back in
	// the error. Enough for a JSON error document from a bus, small enough
	// that a hostile or broken endpoint cannot flood the logs.
	maxErrorBody = 512
	// maxDrain caps the body read purely to make the connection reusable,
	// matching the webhook path in internal/engine/triggers/actions.go.
	maxDrain = 1 << 20
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

// HTTP posts every custom trigger action to a single bus endpoint.
type HTTP struct {
	url    string
	method string
	static map[string]string
	client *http.Client
}

// New builds an HTTP bus forwarder. It returns an error when URL is empty or
// does not parse, when Method is not a valid HTTP method token, or when a
// StaticHeaders name is not a header field name or a value is not a legal
// field value.
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
		if !validToken(m) {
			return nil, fmt.Errorf("forward: invalid method %q", m)
		}
	}
	// Static headers are applied with Header.Set per delivery, and net/http
	// only checks field names and values when it writes the request — so a
	// typo in triggers.forward.static_headers would boot clean and then fail
	// every custom action. Same rule as the URL above: reject it here.
	for k, v := range cfg.StaticHeaders {
		if !validToken(k) {
			return nil, fmt.Errorf("forward: invalid static header name %q", k)
		}
		if !validHeaderValue(v) {
			return nil, fmt.Errorf("forward: invalid value for static header %q: %q", k, v)
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

// validHeaderValue reports whether v can be written as a header field value.
// It mirrors net/http's own write-time check (httpguts.ValidHeaderFieldValue),
// measured rather than assumed: SP, HTAB, obs-text and the empty string are
// accepted by a real request, while CR, LF and the other control bytes are
// refused with `net/http: invalid header field value`.
func validHeaderValue(v string) bool {
	for i := 0; i < len(v); i++ {
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
		return fmt.Errorf("forward: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	ok := resp.StatusCode >= 200 && resp.StatusCode <= 299
	// A non-2xx needs its body. "forward: status 500" tells an operator
	// nothing they can act on, and the bus's own explanation ("{"error":
	// "unknown realm"}") is the whole difference between a fixable failure and
	// a mystery — so read a bounded prefix of it for the error, in the same
	// shape as the container http bridge (blocks/container/httpbridge.go).
	// Bounded because the endpoint is not us: an unbounded read would pull an
	// arbitrarily large body into memory and into a log line.
	var snippet string
	if !ok {
		// The read error is deliberately dropped: a truncated or unreadable
		// body still has to surface as a status error, not as a read error
		// that hides the status.
		b, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
		snippet = strings.TrimSpace(string(b))
		if snippet != "" && len(b) == maxErrorBody {
			snippet += "…"
		}
	}
	// Drain so the connection is reusable. Bounded for the same reason, and to
	// match the sibling request path in this codebase, which caps its drain at
	// the same 1 MiB (triggers/actions.go): two near-identical request paths
	// should not drift on the line that decides how much of a peer's body we
	// are willing to hold.
	if _, err := io.Copy(io.Discard, io.LimitReader(resp.Body, maxDrain)); err != nil {
		return fmt.Errorf("forward: drain: %w", err)
	}
	if !ok {
		if snippet == "" {
			return fmt.Errorf("forward: status %d", resp.StatusCode)
		}
		return fmt.Errorf("forward: status %d: %s", resp.StatusCode, snippet)
	}
	return nil
}
