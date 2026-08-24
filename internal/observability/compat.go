package observability

import (
	"encoding/json"
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

// VersionHandler answers GET .../version with the upstream envelope
// {"data":"<version>"} (measured upstream 1.2.0, verify-versions.json).
func VersionHandler(version string) http.HandlerFunc {
	body, err := json.Marshal(map[string]string{"data": version})
	if err != nil {
		body = []byte(`{"data":""}`)
	}
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
}

// MountVersionCompat registers GET /{service}/version for one service.
func MountVersionCompat(mux *http.ServeMux, service, version string) {
	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
}
