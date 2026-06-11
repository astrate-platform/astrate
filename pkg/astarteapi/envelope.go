// Package astarteapi implements the JSON envelope conventions shared by every
// Astarte-compatible REST surface: success bodies are wrapped as
// {"data": ...}, error bodies as {"errors": {"detail": "..."}} (or the
// field-keyed changeset shape used by 422 validation failures), and request
// bodies arrive wrapped as {"data": ...}.
//
// These exact bytes are parsed by astartectl and the official device SDKs, so
// the golden fixtures under testdata/ are wire-frozen: changing any envelope
// produced here is a compatibility break, not a refactor.
package astarteapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ContentType is the Content-Type header value set on every envelope,
// matching what upstream Astarte's Phoenix endpoints emit.
const ContentType = "application/json; charset=utf-8"

// DefaultMaxBodyBytes is the request body size cap callers are expected to
// pass to DecodeData unless an endpoint has a documented reason to differ
// (interface uploads, for example, may need more than a pairing request).
const DefaultMaxBodyBytes int64 = 1 << 20 // 1 MiB

// Canonical upstream error detail strings. These are frozen: SDK and
// astartectl error paths match on them.
const (
	// DetailBadRequest is the canonical 400 detail.
	DetailBadRequest = "Bad Request"
	// DetailUnauthorized is the canonical 401 detail.
	DetailUnauthorized = "Unauthorized"
	// DetailForbidden is the canonical 403 detail.
	DetailForbidden = "Forbidden"
	// DetailNotFound is the canonical generic 404 detail.
	DetailNotFound = "Not Found"
	// DetailDeviceNotFound is the canonical 404 detail for unknown devices
	// (upstream AppEngine/Pairing shape).
	DetailDeviceNotFound = "Device not found"
	// DetailInternalServerError is the canonical 500 detail.
	DetailInternalServerError = "Internal Server Error"
)

// ErrMissingData is wrapped by DecodeData when the request body has no
// "data" key (or it is JSON null) — upstream rejects such bodies uniformly.
var ErrMissingData = errors.New(`missing "data" key in request body`)

// ErrBodyTooLarge is wrapped by DecodeData when the request body exceeds the
// caller-supplied size cap.
var ErrBodyTooLarge = errors.New("request body too large")

// dataEnvelope is the success wrapper: {"data": ...}.
type dataEnvelope struct {
	Data any `json:"data"`
}

// detailEnvelope is the error wrapper: {"errors": {"detail": "..."}}.
type detailEnvelope struct {
	Errors detailBody `json:"errors"`
}

type detailBody struct {
	Detail string `json:"detail"`
}

// fieldsEnvelope is the changeset-style error wrapper used by 422 validation
// failures: {"errors": {"<field>": ["<message>", ...], ...}}.
type fieldsEnvelope struct {
	Errors map[string][]string `json:"errors"`
}

// marshal encodes v without HTML escaping (upstream does not escape '<' or
// '&' in details or PEM payloads) and without a trailing newline, so envelope
// bytes are stable for golden comparison.
func marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	b := buf.Bytes()
	return b[:len(b)-1], nil // json.Encoder always appends '\n'
}

// write emits one envelope with the given status. If v cannot be marshalled
// the response degrades to the canonical 500 envelope and the marshalling
// error is returned.
func write(w http.ResponseWriter, status int, v any) error {
	body, err := marshal(v)
	if err != nil {
		// The 500 envelope is a static struct; it cannot itself fail.
		body, _ = marshal(detailEnvelope{Errors: detailBody{Detail: DetailInternalServerError}})
		status = http.StatusInternalServerError
		w.Header().Set("Content-Type", ContentType)
		w.WriteHeader(status)
		_, _ = w.Write(body)
		return fmt.Errorf("astarteapi: marshalling response envelope: %w", err)
	}
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("astarteapi: writing response envelope: %w", err)
	}
	return nil
}

// WriteData writes {"data": v} with the given status code.
func WriteData(w http.ResponseWriter, status int, v any) error {
	return write(w, status, dataEnvelope{Data: v})
}

// WriteError writes {"errors": {"detail": detail}} with the given status
// code. Use the canonical constructors below for the frozen upstream shapes;
// use WriteError directly for endpoint-specific details.
func WriteError(w http.ResponseWriter, status int, detail string) error {
	return write(w, status, detailEnvelope{Errors: detailBody{Detail: detail}})
}

// WriteFieldErrors writes the Phoenix-changeset-shaped error envelope
// {"errors": {"<field>": ["<message>", ...]}} used by upstream 422 validation
// failures (for example {"errors": {"hw_id": ["is invalid"]}}). Keys are
// emitted in sorted order (Go map marshalling), which matches the
// deterministic bodies upstream produces for single-field failures.
func WriteFieldErrors(w http.ResponseWriter, status int, fields map[string][]string) error {
	return write(w, status, fieldsEnvelope{Errors: fields})
}

// WriteBadRequest writes the canonical 400 envelope.
func WriteBadRequest(w http.ResponseWriter) error {
	return WriteError(w, http.StatusBadRequest, DetailBadRequest)
}

// WriteUnauthorized writes the canonical 401 envelope.
func WriteUnauthorized(w http.ResponseWriter) error {
	return WriteError(w, http.StatusUnauthorized, DetailUnauthorized)
}

// WriteForbidden writes the canonical 403 envelope.
func WriteForbidden(w http.ResponseWriter) error {
	return WriteError(w, http.StatusForbidden, DetailForbidden)
}

// WriteNotFound writes the canonical generic 404 envelope.
func WriteNotFound(w http.ResponseWriter) error {
	return WriteError(w, http.StatusNotFound, DetailNotFound)
}

// WriteDeviceNotFound writes the canonical 404 "Device not found" envelope.
func WriteDeviceNotFound(w http.ResponseWriter) error {
	return WriteError(w, http.StatusNotFound, DetailDeviceNotFound)
}

// WriteInternalServerError writes the canonical 500 envelope.
func WriteInternalServerError(w http.ResponseWriter) error {
	return WriteError(w, http.StatusInternalServerError, DetailInternalServerError)
}

// DecodeData reads at most maxBytes bytes from r, unwraps the mandatory
// {"data": ...} request envelope, and unmarshals the "data" value into dst.
//
// It fails with an error wrapping ErrBodyTooLarge when the body exceeds
// maxBytes, with one wrapping ErrMissingData when the "data" key is absent or
// null, and with the underlying JSON error for malformed bodies (including
// trailing garbage after the top-level value). Sibling keys next to "data"
// are ignored, matching upstream parameter handling.
func DecodeData(r io.Reader, maxBytes int64, dst any) error {
	body, err := io.ReadAll(io.LimitReader(r, maxBytes+1))
	if err != nil {
		return fmt.Errorf("astarteapi: reading request body: %w", err)
	}
	if int64(len(body)) > maxBytes {
		return fmt.Errorf("astarteapi: %w: exceeds %d bytes", ErrBodyTooLarge, maxBytes)
	}
	var env struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("astarteapi: decoding request envelope: %w", err)
	}
	if len(env.Data) == 0 || bytes.Equal(env.Data, []byte("null")) {
		return fmt.Errorf("astarteapi: %w", ErrMissingData)
	}
	if err := json.Unmarshal(env.Data, dst); err != nil {
		return fmt.Errorf("astarteapi: decoding request data: %w", err)
	}
	return nil
}
