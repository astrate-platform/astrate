package triggers

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"unicode/utf8"
)

// Upstream http_action.ex v1.2.0 limits, probe-frozen (issue #63).
const (
	minURLRunes      = 8
	maxURLRunes      = 8192
	maxHeadersBytes  = 8192 // Σ(byte_size(name)+2+byte_size(value)) must be < this
	maxTemplateRunes = 1048576
)

// blockedHeaderNames is upstream's hop-by-hop/sensitive header blocklist,
// compared against strings.ToLower(strings.TrimSpace(name)).
var blockedHeaderNames = map[string]bool{
	"connection":           true,
	"content-length":       true,
	"date":                 true,
	"host":                 true,
	"te":                   true,
	"upgrade":              true,
	"x-forwarded-for":      true,
	"x-forwarded-host":     true,
	"x-forwarded-proto":    true,
	"sec-websocket-accept": true,
	"proxy-authorization":  true,
	"proxy-authenticate":   true,
}

// FieldErrors carries an upstream-shaped field→messages map for one part of a
// trigger definition. errors.Is/As-friendly so the realm layer can render the
// nested changeset envelope without string parsing.
type FieldErrors struct {
	Part   string // e.g. "action"
	Fields map[string][]string
}

// Error renders a stable human form ("action: field=message, ..." with fields
// sorted) so logs and plain-error consumers stay deterministic.
func (e *FieldErrors) Error() string {
	fields := make([]string, 0, len(e.Fields))
	for f := range e.Fields {
		fields = append(fields, f)
	}
	sort.Strings(fields)
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		parts = append(parts, fmt.Sprintf("%s=%s", f, strings.Join(e.Fields[f], ", ")))
	}
	return e.Part + ": " + strings.Join(parts, ", ")
}

// newFieldErrors starts an empty accumulation for the "action" part.
func newFieldErrors() *FieldErrors {
	return &FieldErrors{Part: "action", Fields: map[string][]string{}}
}

// add appends one message for field, preserving accumulation order.
func (e *FieldErrors) add(field, msg string) {
	e.Fields[field] = append(e.Fields[field], msg)
}

// validateURLField applies upstream's URL rules to one URL-typed field
// (http_url, or the legacy http_post_url): length messages first, then the
// format message appended after them on the same field.
func validateURLField(e *FieldErrors, field, raw string) {
	n := utf8.RuneCountInString(raw)
	switch {
	case n < minURLRunes:
		e.add(field, fmt.Sprintf("should be at least %d character(s)", minURLRunes))
	case n > maxURLRunes:
		e.add(field, fmt.Sprintf("should be at most %d character(s)", maxURLRunes))
	}
	if !isHTTPURL(raw) {
		e.add(field, "must be a valid http(s) URL")
	}
}

// isHTTPURL reports whether raw parses as an absolute URL with an http(s)
// scheme and a non-empty host (upstream's scheme+host check).
func isHTTPURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

// validateStaticHeaders applies the modern-branch header checks: the
// blocklist first, then the total-size ceiling; both can fail together.
func validateStaticHeaders(e *FieldErrors, headers map[string]string) {
	for name := range headers {
		blocked := blockedHeaderNames[strings.ToLower(strings.TrimSpace(name))]
		if blocked {
			e.add("http_static_headers", "must contain only allowed http headers")
			break
		}
	}
	total := 0
	for name, value := range headers {
		total += len(name) + 2 + len(value)
	}
	if total >= maxHeadersBytes {
		e.add("http_static_headers", fmt.Sprintf("headers total size must be lower than %d", maxHeadersBytes))
	}
}

// validateTemplate enforces the template rune-length cap regardless of
// template_type (upstream never validates template_type itself).
func validateTemplate(e *FieldErrors, template string) {
	if utf8.RuneCountInString(template) > maxTemplateRunes {
		e.add("template", fmt.Sprintf("should be at most %d character(s)", maxTemplateRunes))
	}
}
