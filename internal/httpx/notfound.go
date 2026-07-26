package httpx

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

// NotFound returns a handler that serves mux, replacing Go's plain-text
// "404 page not found" with the JSON error envelope upstream Astarte emits for
// an unmatched route under a known service prefix.
func NotFound(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pattern := mux.Handler(r)
		if pattern != "" {
			mux.ServeHTTP(w, r)
			return
		}

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, r)

		// Replay what the mux produced, byte for byte.
		replay := func() {
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(rec.Code)
			_, _ = rec.Body.WriteTo(w)
		}

		// A method mismatch also reaches here with an empty pattern, but the
		// mux answers it 405: only a genuine 404 gets the envelope.
		if rec.Code != http.StatusNotFound {
			replay()
			return
		}

		seg := strings.TrimPrefix(r.URL.Path, "/")
		seg, _, _ = strings.Cut(seg, "/")

		switch seg {
		case "appengine", "realmmanagement":
			_ = astarteapi.WriteError(w, http.StatusNotFound, astarteapi.DetailRouteNotFound)
		case "pairing":
			_ = astarteapi.WriteError(w, http.StatusNotFound, astarteapi.DetailPageNotFound)
		default:
			// Only the three segments above were observed upstream; every
			// other prefix keeps today's plain-text 404 deliberately.
			replay()
		}
	})
}
