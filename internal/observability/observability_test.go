package observability

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetricsExposesGauges(t *testing.T) {
	m := NewMetrics()
	m.RegisterBrokerSessions(func() float64 { return 3 })
	m.RegisterDBPool(func() DBPoolStats {
		return DBPoolStats{AcquiredConns: 2, IdleConns: 5, TotalConns: 7, MaxConns: 10}
	})

	rec := httptest.NewRecorder()
	m.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/metrics", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("scrape status = %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		"astrate_broker_sessions 3",
		"astrate_db_pool_acquired_conns 2",
		"astrate_db_pool_max_conns 10",
		"go_goroutines", // runtime collector present
	} {
		if !strings.Contains(body, want) {
			t.Errorf("scrape missing %q\n%s", want, body)
		}
	}
}

func TestReadiness(t *testing.T) {
	m := NewMetrics()
	h := NewHealth(m.Handler())

	mux := http.NewServeMux()
	h.AddReadiness("database", func(context.Context) error { return nil })
	failing := errors.New("down")
	h.AddReadiness("broker", func(context.Context) error { return failing })
	h.Mount(mux)

	// Liveness is always 200.
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/health", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("health = %d, want 200", rec.Code)
	}

	// One failing dependency makes readiness 503.
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/readiness", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("readiness = %d, want 503", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "down") {
		t.Errorf("readiness body should name the failing check: %s", rec.Body)
	}
}

func TestReadinessAllOK(t *testing.T) {
	h := NewHealth(NewMetrics().Handler())
	h.AddReadiness("database", func(context.Context) error { return nil })
	mux := http.NewServeMux()
	h.Mount(mux)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/readiness", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("readiness = %d, want 200", rec.Code)
	}
}
