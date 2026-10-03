package appengine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astrate-platform/astrate/internal/engine"
	"github.com/astrate-platform/astrate/pkg/payload"
)

func TestWriteErrorTaxonomy(t *testing.T) {
	a := &API{}
	cases := []struct {
		name   string
		err    error
		status int
		detail string
		keys   []string
	}{
		{"device-owned write", fmt.Errorf("%w: context", engine.ErrNotServerOwned),
			http.StatusMethodNotAllowed, "Cannot write to device owned resource", nil},
		{"unset on server-owned datastream", fmt.Errorf("%w: context", engine.ErrNotAProperty),
			http.StatusMethodNotAllowed, "Cannot write to read-only resource", nil},
		{"unknown interface", fmt.Errorf("%w: context", engine.ErrInterfaceNotFound),
			http.StatusNotFound, "Interface not found in device introspection", nil},
		{"unknown endpoint", fmt.Errorf("%w: context", engine.ErrPathNotFound),
			http.StatusBadRequest, "Endpoint not found", nil},
		{"read-path path-not-found collision guard", ErrPathNotFound,
			http.StatusNotFound, "Path not found", nil},
		{"object write with an undeclared key answers 400 naming it",
			&payload.RejectError{Reason: payload.ReasonBadObject,
				Detail:         `object keys [ghost] match no declared object leaf`,
				UnexpectedKeys: []string{"ghost"}},
			http.StatusBadRequest, "Unexpected object key", []string{"ghost"}},
		{"object write with several undeclared keys lists them sorted",
			fmt.Errorf("engine: %w", &payload.RejectError{Reason: payload.ReasonBadObject,
				Detail:         "object keys [a b] match no declared object leaf",
				UnexpectedKeys: []string{"a", "b"}}),
			http.StatusBadRequest, "Unexpected object key", []string{"a", "b"}},
		{"object write with no key list stays the canonical 400",
			&payload.RejectError{Reason: payload.ReasonBadObject, Detail: "object-aggregation document is empty"},
			http.StatusBadRequest, "Bad request", nil},
		{"object write omitting a required key answers 422",
			&payload.RejectError{Reason: payload.ReasonMissingRequired,
				Detail: "object-aggregation document is missing required key(s) lat"},
			http.StatusUnprocessableEntity, "Missing required mapping key", nil},
		{"unknown cause stays 500", errors.New("boom"),
			http.StatusInternalServerError, "Internal server error", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			a.writeError(rec, tc.err)
			if rec.Code != tc.status {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tc.status, rec.Body.String())
			}
			var body struct {
				Errors struct {
					Detail string   `json:"detail"`
					Keys   []string `json:"unexpected_keys"`
				} `json:"errors"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("body %q: %v", rec.Body.String(), err)
			}
			if body.Errors.Detail != tc.detail {
				t.Errorf("detail = %q, want %q", body.Errors.Detail, tc.detail)
			}
			if len(body.Errors.Keys) != len(tc.keys) {
				t.Fatalf("unexpected_keys = %v, want %v", body.Errors.Keys, tc.keys)
			}
			for i, k := range tc.keys {
				if body.Errors.Keys[i] != k {
					t.Errorf("unexpected_keys[%d] = %q, want %q", i, body.Errors.Keys[i], k)
				}
			}
		})
	}
}

// TestWriteErrorBadObjectEnvelope freezes the exact bytes upstream emits for
// an object write carrying an undeclared key (master b6d46ad4, #2237):
// {"errors":{"detail":"Unexpected object key","unexpected_keys":["ghost"]}}.
func TestWriteErrorBadObjectEnvelope(t *testing.T) {
	a := &API{}
	rec := httptest.NewRecorder()
	a.writeError(rec, &payload.RejectError{
		Reason:         payload.ReasonBadObject,
		Detail:         "object keys [ghost] match no declared object leaf",
		UnexpectedKeys: []string{"ghost"},
	})
	want := `{"errors":{"detail":"Unexpected object key","unexpected_keys":["ghost"]}}`
	if got := rec.Body.String(); got != want {
		t.Errorf("body = %s, want %s", got, want)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q", ct)
	}
}
