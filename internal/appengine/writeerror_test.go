package appengine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astrate-platform/astrate/internal/engine"
)

func TestWriteErrorTaxonomy(t *testing.T) {
	a := &API{}
	cases := []struct {
		name   string
		err    error
		status int
		detail string
	}{
		{"device-owned write", fmt.Errorf("%w: context", engine.ErrNotServerOwned),
			http.StatusMethodNotAllowed, "Cannot write to device owned resource"},
		{"unset on server-owned datastream", fmt.Errorf("%w: context", engine.ErrNotAProperty),
			http.StatusMethodNotAllowed, "Cannot write to read-only resource"},
		{"unknown interface", fmt.Errorf("%w: context", engine.ErrInterfaceNotFound),
			http.StatusNotFound, "Interface not found in device introspection"},
		{"unknown endpoint", fmt.Errorf("%w: context", engine.ErrPathNotFound),
			http.StatusBadRequest, "Endpoint not found"},
		{"read-path path-not-found collision guard", ErrPathNotFound,
			http.StatusNotFound, "Path not found"},
		{"unknown cause stays 500", errors.New("boom"),
			http.StatusInternalServerError, "Internal server error"},
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
					Detail string `json:"detail"`
				} `json:"errors"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("body %q: %v", rec.Body.String(), err)
			}
			if body.Errors.Detail != tc.detail {
				t.Errorf("detail = %q, want %q", body.Errors.Detail, tc.detail)
			}
		})
	}
}
