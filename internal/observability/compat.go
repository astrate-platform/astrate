package observability

import (
	"context"
	"encoding/json"
	"net/http"
)

// MountServiceCompat registers the upstream-parity unauthenticated
// GET /{service}/health endpoint. The Astarte Dashboard polls one per service
// (appengine, realmmanagement, pairing) for its API status indicators; only
// the status code matters to it, but the body mirrors upstream's envelope.
//
// Deliberate deviation 18 (docs/COMPATIBILITY.md): upstream answers a static 200
// here, so the indicator it feeds cannot ever go red. Astrate runs check — the
// same dependency probe that backs /astrate/v1/readiness — and answers 503
// when it fails, which is strictly safer for whoever is reading the Dashboard.
// A nil check keeps upstream's static 200.
func MountServiceCompat(mux *http.ServeMux, service string, check Check) {
	mux.HandleFunc("GET /"+service+"/health", func(w http.ResponseWriter, r *http.Request) {
		if check != nil {
			ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
			defer cancel()
			if err := check(ctx); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte(`{"data":{"status":"unhealthy"}}`))
				return
			}
		}
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
