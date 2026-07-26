package httpx

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNotFound(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /appengine/v1/{realm}/devices", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Test", "appengine")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("devices-ok"))
	})

	mux.HandleFunc("GET /housekeeping/v1/realms", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("realms-ok"))
	})

	h := NotFound(mux)

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
		wantCT     string
	}{
		{
			name:       "appengine_unmatched_gets_json_envelope",
			method:     "GET",
			path:       "/appengine/v1/nope",
			wantStatus: 404,
			wantBody:   `{"errors":{"detail":"Not found"}}`,
			wantCT:     "application/json; charset=utf-8",
		},
		{
			name:       "realmmanagement_unmatched_gets_json_envelope",
			method:     "GET",
			path:       "/realmmanagement/v1/nope",
			wantStatus: 404,
			wantBody:   `{"errors":{"detail":"Not found"}}`,
			wantCT:     "application/json; charset=utf-8",
		},
		{
			name:       "pairing_unmatched_gets_json_envelope",
			method:     "GET",
			path:       "/pairing/v1/nope",
			wantStatus: 404,
			wantBody:   `{"errors":{"detail":"Page not found"}}`,
			wantCT:     "application/json; charset=utf-8",
		},
		{
			name:       "housekeeping_unmatched_gets_json_envelope",
			method:     "GET",
			path:       "/housekeeping/v1/nope",
			wantStatus: 404,
			wantBody:   `{"errors":{"detail":"Page not found"}}`,
			wantCT:     "application/json; charset=utf-8",
		},
		{
			name:       "unknown_prefix_keeps_go_default",
			method:     "GET",
			path:       "/nope",
			wantStatus: 404,
			wantBody:   "404 page not found\n",
			wantCT:     "text/plain; charset=utf-8",
		},
		{
			name:       "matched_route_is_transparent",
			method:     "GET",
			path:       "/appengine/v1/someRealm/devices",
			wantStatus: 200,
			wantBody:   "devices-ok",
		},
		{
			name:       "method_mismatch_stays_405",
			method:     "POST",
			path:       "/appengine/v1/someRealm/devices",
			wantStatus: 405,
			wantBody:   "Method Not Allowed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("status = %d; want %d", rr.Code, tt.wantStatus)
			}

			body := rr.Body.String()
			if body != tt.wantBody {
				t.Errorf("body = %q; want %q", body, tt.wantBody)
			}

			if tt.wantCT != "" {
				ct := rr.Header().Get("Content-Type")
				if ct != tt.wantCT {
					t.Errorf("Content-Type = %q; want %q", ct, tt.wantCT)
				}
			}
		})
	}
}

func TestNotFound_TransparentBodyPassthrough(t *testing.T) {
	longBody := strings.Repeat("x", 64*1024)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /appengine/v1/{path...}", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Custom", "yes")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(longBody))
	})

	h := NotFound(mux)
	req := httptest.NewRequest("GET", "/appengine/v1/some/path", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("status = %d; want 200", rr.Code)
	}
	if rr.Body.String() != longBody {
		t.Errorf("body length = %d; want %d", rr.Body.Len(), len(longBody))
	}
	if rr.Header().Get("X-Custom") != "yes" {
		t.Error("X-Custom header not forwarded")
	}
}
