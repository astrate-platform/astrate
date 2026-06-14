package observability

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// readinessTimeout bounds the whole readiness probe so a wedged dependency
// can't hang the endpoint (and the orchestrator that polls it).
const readinessTimeout = 3 * time.Second

// Check reports a dependency's readiness; a nil error means ready. The store
// DB ping and a broker-listener check are the two §5.2 readiness backends.
type Check func(ctx context.Context) error

// Health is the Astrate-native operational surface (docs/DESIGN.md §5.2):
// liveness, readiness, and the Prometheus scrape endpoint, all under
// /astrate/v1 so they never collide with the upstream API namespace.
type Health struct {
	metrics http.Handler
	ready   []namedCheck
}

type namedCheck struct {
	name  string
	check Check
}

// NewHealth builds the surface around an existing metrics handler.
func NewHealth(metrics http.Handler) *Health {
	return &Health{metrics: metrics}
}

// AddReadiness registers a named dependency probed by /readiness.
func (h *Health) AddReadiness(name string, check Check) {
	h.ready = append(h.ready, namedCheck{name: name, check: check})
}

// Mount registers the endpoints on mux.
func (h *Health) Mount(mux *http.ServeMux) {
	mux.HandleFunc("GET /astrate/v1/health", h.handleHealth)
	mux.HandleFunc("GET /astrate/v1/readiness", h.handleReadiness)
	mux.Handle("GET /astrate/v1/metrics", h.metrics)
}

// handleHealth is liveness: the process is up and serving.
func (h *Health) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleReadiness runs every dependency check; any failure yields 503 with the
// per-check status so operators see which dependency is down.
func (h *Health) handleReadiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
	defer cancel()

	checks := make(map[string]string, len(h.ready))
	ready := true
	for _, nc := range h.ready {
		if err := nc.check(ctx); err != nil {
			ready = false
			checks[nc.name] = "error: " + err.Error()
		} else {
			checks[nc.name] = "ok"
		}
	}

	status := http.StatusOK
	overall := "ok"
	if !ready {
		status = http.StatusServiceUnavailable
		overall = "unavailable"
	}
	writeJSON(w, status, map[string]any{"status": overall, "checks": checks})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
