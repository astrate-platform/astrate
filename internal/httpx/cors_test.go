package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func corsHandler(origins []string) http.Handler {
	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot) // sentinel: the inner handler ran
	})
	return CORS(origins)(inner)
}

func doReq(t *testing.T, h http.Handler, method, origin, acrm string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, "/appengine/v1/test/devices", nil)
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	if acrm != "" {
		req.Header.Set("Access-Control-Request-Method", acrm)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestCORSPreflight(t *testing.T) {
	h := corsHandler([]string{"http://localhost:4040"})
	rec := doReq(t, h, http.MethodOptions, "http://localhost:4040", "GET")

	if rec.Code != http.StatusNoContent {
		t.Fatalf("preflight status = %d, want 204", rec.Code)
	}
	for header, want := range map[string]string{
		"Access-Control-Allow-Origin":  "http://localhost:4040",
		"Access-Control-Allow-Methods": "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Authorization, Content-Type",
		"Access-Control-Max-Age":       "600",
		"Vary":                         "Origin",
	} {
		if got := rec.Header().Get(header); got != want {
			t.Errorf("%s = %q, want %q", header, got, want)
		}
	}
}

func TestCORSActualRequestStamped(t *testing.T) {
	h := corsHandler([]string{"http://localhost:4040"})
	rec := doReq(t, h, http.MethodGet, "http://localhost:4040", "")

	if rec.Code != http.StatusTeapot {
		t.Fatalf("actual request did not reach the inner handler (status %d)", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:4040" {
		t.Errorf("Allow-Origin = %q", got)
	}
}

func TestCORSDisallowedOriginPassesThroughUnstamped(t *testing.T) {
	h := corsHandler([]string{"http://localhost:4040"})
	rec := doReq(t, h, http.MethodGet, "http://evil.example", "")

	if rec.Code != http.StatusTeapot {
		t.Fatalf("disallowed origin blocked the request (status %d) — the browser enforces, not us", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Allow-Origin stamped for disallowed origin: %q", got)
	}
	// A preflight from a disallowed origin reaches the mux (which 404s/405s it);
	// no CORS approval is expressed either way.
	rec = doReq(t, h, http.MethodOptions, "http://evil.example", "GET")
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("preflight Allow-Origin stamped for disallowed origin: %q", got)
	}
}

func TestCORSWildcard(t *testing.T) {
	h := corsHandler([]string{"*"})
	rec := doReq(t, h, http.MethodOptions, "http://anything.example", "POST")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("wildcard preflight status = %d", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Allow-Origin = %q, want *", got)
	}
	if got := rec.Header().Get("Vary"); got != "" {
		t.Errorf("wildcard response should not vary on Origin, got %q", got)
	}
}

func TestCORSNoOriginUntouched(t *testing.T) {
	h := corsHandler([]string{"http://localhost:4040"})
	rec := doReq(t, h, http.MethodGet, "", "")
	if rec.Code != http.StatusTeapot {
		t.Fatalf("non-CORS request blocked (status %d)", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Allow-Origin stamped without Origin header: %q", got)
	}
}

// OPTIONS without Access-Control-Request-Method is not a preflight and must
// reach the mux.
func TestCORSPlainOptionsNotPreflight(t *testing.T) {
	h := corsHandler([]string{"http://localhost:4040"})
	rec := doReq(t, h, http.MethodOptions, "http://localhost:4040", "")
	if rec.Code != http.StatusTeapot {
		t.Fatalf("plain OPTIONS swallowed by CORS middleware (status %d)", rec.Code)
	}
}
