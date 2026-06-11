package astarteapi_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

// TestGoldenEnvelopes freezes the exact bytes of every envelope constructor.
// These are the bodies astartectl and the official SDK error paths parse;
// any diff here is a wire-compatibility break.
func TestGoldenEnvelopes(t *testing.T) {
	cases := []struct {
		name       string
		golden     string
		wantStatus int
		write      func(w http.ResponseWriter) error
	}{
		{
			name: "data object", golden: "data_object.json", wantStatus: 201,
			write: func(w http.ResponseWriter) error {
				return astarteapi.WriteData(w, 201, map[string]string{
					"credentials_secret": "TTkd0vDHgAs6FSyhmnnpfMZlDOdUWLNxNHbLeuTpegc=",
				})
			},
		},
		{
			name: "data array", golden: "data_array.json", wantStatus: 200,
			write: func(w http.ResponseWriter) error {
				return astarteapi.WriteData(w, 200, []string{"dT6hS2W9TT6LEnP25ks_lg"})
			},
		},
		{
			name: "data null", golden: "data_null.json", wantStatus: 200,
			write: func(w http.ResponseWriter) error {
				return astarteapi.WriteData(w, 200, nil)
			},
		},
		{
			name: "data no html escaping", golden: "data_pem.json", wantStatus: 201,
			write: func(w http.ResponseWriter) error {
				// PEM and URL payloads must round-trip without <-style escapes.
				return astarteapi.WriteData(w, 201, map[string]string{
					"client_crt": "-----BEGIN CERTIFICATE-----\nMIIB<&>\n-----END CERTIFICATE-----",
				})
			},
		},
		{
			name: "bad request", golden: "error_bad_request.json", wantStatus: 400,
			write: astarteapi.WriteBadRequest,
		},
		{
			name: "unauthorized", golden: "error_unauthorized.json", wantStatus: 401,
			write: astarteapi.WriteUnauthorized,
		},
		{
			name: "forbidden", golden: "error_forbidden.json", wantStatus: 403,
			write: astarteapi.WriteForbidden,
		},
		{
			name: "not found", golden: "error_not_found.json", wantStatus: 404,
			write: astarteapi.WriteNotFound,
		},
		{
			name: "device not found", golden: "error_device_not_found.json", wantStatus: 404,
			write: astarteapi.WriteDeviceNotFound,
		},
		{
			name: "internal server error", golden: "error_internal.json", wantStatus: 500,
			write: astarteapi.WriteInternalServerError,
		},
		{
			name: "custom detail", golden: "error_custom_detail.json", wantStatus: 422,
			write: func(w http.ResponseWriter) error {
				return astarteapi.WriteError(w, 422, "Invalid interface document")
			},
		},
		{
			name: "field errors", golden: "error_fields.json", wantStatus: 422,
			write: func(w http.ResponseWriter) error {
				return astarteapi.WriteFieldErrors(w, 422, map[string][]string{
					"hw_id": {"is invalid"},
				})
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			if err := tc.write(rec); err != nil {
				t.Fatalf("write: %v", err)
			}
			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if ct := rec.Header().Get("Content-Type"); ct != astarteapi.ContentType {
				t.Errorf("Content-Type = %q, want %q", ct, astarteapi.ContentType)
			}
			testutil.Golden(t, tc.golden, rec.Body.Bytes())
		})
	}
}

func TestWriteDataMarshalFailure(t *testing.T) {
	rec := httptest.NewRecorder()
	err := astarteapi.WriteData(rec, 200, func() {}) // funcs are not marshallable
	if err == nil {
		t.Fatal("WriteData(func) succeeded, want error")
	}
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
	testutil.Golden(t, "error_internal.json", rec.Body.Bytes())
}

func TestDecodeData(t *testing.T) {
	type registration struct {
		HWID string `json:"hw_id"`
	}

	t.Run("happy path", func(t *testing.T) {
		var dst registration
		body := `{"data": {"hw_id": "dT6hS2W9TT6LEnP25ks_lg"}}`
		if err := astarteapi.DecodeData(strings.NewReader(body), 1024, &dst); err != nil {
			t.Fatalf("DecodeData: %v", err)
		}
		if dst.HWID != "dT6hS2W9TT6LEnP25ks_lg" {
			t.Errorf("hw_id = %q", dst.HWID)
		}
	})

	t.Run("sibling keys ignored", func(t *testing.T) {
		var dst registration
		body := `{"meta": 1, "data": {"hw_id": "x"}, "other": [true]}`
		if err := astarteapi.DecodeData(strings.NewReader(body), 1024, &dst); err != nil {
			t.Fatalf("DecodeData: %v", err)
		}
		if dst.HWID != "x" {
			t.Errorf("hw_id = %q, want %q", dst.HWID, "x")
		}
	})

	t.Run("scalar data", func(t *testing.T) {
		var dst int
		if err := astarteapi.DecodeData(strings.NewReader(`{"data": 42}`), 1024, &dst); err != nil {
			t.Fatalf("DecodeData: %v", err)
		}
		if dst != 42 {
			t.Errorf("data = %d, want 42", dst)
		}
	})

	rejections := []struct {
		name    string
		body    string
		max     int64
		wantErr error
	}{
		{"missing data", `{"hw_id": "x"}`, 1024, astarteapi.ErrMissingData},
		{"null data", `{"data": null}`, 1024, astarteapi.ErrMissingData},
		{"empty object", `{}`, 1024, astarteapi.ErrMissingData},
		{"empty body", ``, 1024, nil},
		{"not json", `data=x`, 1024, nil},
		{"top-level array", `[{"data": 1}]`, 1024, nil},
		{"trailing garbage", `{"data": 1} {"data": 2}`, 1024, nil},
		{"oversized", `{"data": "` + strings.Repeat("a", 64) + `"}`, 32, astarteapi.ErrBodyTooLarge},
		{"oversized by one", `{"data": 11}`, 11, astarteapi.ErrBodyTooLarge},
	}
	for _, tc := range rejections {
		t.Run(tc.name, func(t *testing.T) {
			var dst any
			err := astarteapi.DecodeData(strings.NewReader(tc.body), tc.max, &dst)
			if err == nil {
				t.Fatalf("DecodeData(%q) succeeded, want error", tc.body)
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Errorf("error %v does not wrap %v", err, tc.wantErr)
			}
		})
	}

	t.Run("at limit exactly", func(t *testing.T) {
		body := `{"data": 11}` // 12 bytes
		var dst int
		if err := astarteapi.DecodeData(strings.NewReader(body), 12, &dst); err != nil {
			t.Fatalf("DecodeData at exact limit: %v", err)
		}
		if dst != 11 {
			t.Errorf("data = %d, want 11", dst)
		}
	})

	t.Run("type mismatch", func(t *testing.T) {
		var dst registration
		err := astarteapi.DecodeData(strings.NewReader(`{"data": [1, 2]}`), 1024, &dst)
		if err == nil {
			t.Fatal("DecodeData(array into struct) succeeded, want error")
		}
	})
}
