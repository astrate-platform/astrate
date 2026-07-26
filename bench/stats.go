package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// histogram collects latency observations; percentiles are computed exactly
// over the sorted sample (bench runs stay well under memory-relevant sizes).
type histogram struct {
	mu sync.Mutex
	v  []time.Duration
}

func (h *histogram) add(d time.Duration) {
	h.mu.Lock()
	h.v = append(h.v, d)
	h.mu.Unlock()
}

// summary is the JSON-serializable percentile digest of one histogram.
type summary struct {
	Count int     `json:"count"`
	P50ms float64 `json:"p50_ms"`
	P95ms float64 `json:"p95_ms"`
	P99ms float64 `json:"p99_ms"`
	MaxMs float64 `json:"max_ms"`
}

func (h *histogram) summarize() summary {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.v) == 0 {
		return summary{}
	}
	sorted := append([]time.Duration(nil), h.v...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	pct := func(p float64) float64 {
		idx := int(p * float64(len(sorted)-1))
		return float64(sorted[idx]) / float64(time.Millisecond)
	}
	return summary{
		Count: len(sorted),
		P50ms: pct(0.50),
		P95ms: pct(0.95),
		P99ms: pct(0.99),
		MaxMs: float64(sorted[len(sorted)-1]) / float64(time.Millisecond),
	}
}

func (s summary) String() string {
	if s.Count == 0 {
		return "n=0"
	}
	return fmt.Sprintf("n=%d p50=%.1fms p95=%.1fms p99=%.1fms max=%.1fms",
		s.Count, s.P50ms, s.P95ms, s.P99ms, s.MaxMs)
}

// writeResult writes one result document into outDir (created if needed) as
// <name>-<timestamp>.json and returns the path.
func writeResult(outDir, name string, doc any) (string, error) {
	if outDir == "" {
		return "", nil
	}
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		return "", err
	}
	path := filepath.Join(outDir, fmt.Sprintf("%s-%s.json", name, time.Now().Format("20060102-150405")))
	raw, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil { //nolint:gosec // benchmark results, not secrets
		return "", err
	}
	return path, nil
}
