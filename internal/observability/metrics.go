// Package observability owns Astrate's Prometheus registry and the
// health/readiness/metrics HTTP surface under /astrate/v1 (docs/DESIGN.md
// §5.2). The engine and trigger executor self-register their collectors
// through Registerer(); the broker-session and DB-pool gauges are attached
// from caller-supplied closures, so this package imports neither the broker
// nor the store and stays at the top of the dependency order.
package observability

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// metricsNamespace prefixes every Astrate-defined series (astrate_*).
const metricsNamespace = "astrate"

// Metrics owns the process-wide Prometheus registry.
type Metrics struct {
	reg *prometheus.Registry
}

// NewMetrics builds a registry preloaded with the Go runtime and process
// collectors (go_*, process_* — the RSS/GC budget of §5.4 reads them).
func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	return &Metrics{reg: reg}
}

// Registerer is the engine's Prometheus.Registerer (engine + executor + stream
// bus collectors attach here).
func (m *Metrics) Registerer() prometheus.Registerer { return m.reg }

// Handler serves the registry in the Prometheus text exposition format.
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})
}

// DBPoolStats is the subset of pgxpool.Stat exposed as gauges; the caller maps
// store.Stat() onto it (keeping pgx out of this package).
type DBPoolStats struct {
	AcquiredConns int32
	IdleConns     int32
	TotalConns    int32
	MaxConns      int32
}

// RegisterBrokerSessions publishes astrate_broker_sessions, read on scrape from
// fn (broker.SessionCount).
func (m *Metrics) RegisterBrokerSessions(fn func() float64) {
	m.reg.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Subsystem: "broker",
		Name:      "sessions",
		Help:      "Live authenticated MQTT device sessions.",
	}, fn))
}

// RegisterDBPool publishes the astrate_db_pool_* gauges, read on scrape from fn
// (a snapshot of store.Stat()).
func (m *Metrics) RegisterDBPool(fn func() DBPoolStats) {
	gauge := func(name, help string, pick func(DBPoolStats) float64) {
		m.reg.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Subsystem: "db_pool",
			Name:      name,
			Help:      help,
		}, func() float64 { return pick(fn()) }))
	}
	gauge("acquired_conns", "Connections currently in use.", func(s DBPoolStats) float64 { return float64(s.AcquiredConns) })
	gauge("idle_conns", "Idle connections in the pool.", func(s DBPoolStats) float64 { return float64(s.IdleConns) })
	gauge("total_conns", "Total connections in the pool.", func(s DBPoolStats) float64 { return float64(s.TotalConns) })
	gauge("max_conns", "Configured maximum pool size.", func(s DBPoolStats) float64 { return float64(s.MaxConns) })
}
