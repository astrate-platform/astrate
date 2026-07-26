package observability

import (
	"net/http"
)

// MountServiceCompat registers the upstream-parity unauthenticated
// GET /{service}/health endpoint. The Astarte Dashboard polls one per service
// (appengine, realmmanagement, pairing) for its API status indicators; only
// the 2xx matters to it, but the body mirrors upstream's envelope.
func MountServiceCompat(mux *http.ServeMux, service string) {
	mux.HandleFunc("GET /"+service+"/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"status":"ok"}}`))
	})
}
