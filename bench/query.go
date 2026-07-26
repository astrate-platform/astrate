package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"
)

// queryResult is the JSON document one query run produces.
type queryResult struct {
	Target      string             `json:"target"`
	Devices     int                `json:"devices"`
	PerType     int                `json:"iterations_per_type"`
	Concurrency int                `json:"concurrency"`
	Errors      int64              `json:"errors"`
	Latency     map[string]summary `json:"latency"`
}

// cmdQuery runs a canned AppEngine read mix against previously ingested data
// and reports latency percentiles per query type. Run it after `bench
// ingest` so the series are non-empty.
func cmdQuery(args []string) error {
	fs := flag.NewFlagSet("query", flag.ExitOnError)
	statePath := fs.String("state", "bench-state.json", "state file from provision")
	devices := fs.Int("devices", 100, "devices to spread queries over")
	n := fs.Int("n", 200, "iterations per query type")
	concurrency := fs.Int("concurrency", 8, "parallel requests")
	outDir := fs.String("out", "results", "directory for the result JSON ('' to disable)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	st, err := loadState(*statePath)
	if err != nil {
		return err
	}
	realmKey, err := parsePrivateKeyPEM([]byte(st.RealmKeyPEM))
	if err != nil {
		return err
	}
	_, workers, err := st.workload(*devices)
	if err != nil {
		return err
	}
	token, err := mintJWT(realmKey, "a_aea")
	if err != nil {
		return err
	}
	c := newClient(30 * time.Second)

	base := func(dev Device) string {
		return st.Endpoints.AppEngine + "/v1/" + st.Realm + "/devices/" + dev.ID + "/interfaces/" + ifaceIndividual
	}
	types := map[string]func(Device) string{
		// Latest sample of one path: the dashboard/last-value pattern.
		"latest": func(d Device) string {
			return base(d) + "/value?" + url.Values{"limit": {"1"}}.Encode()
		},
		// One-hour window, bounded page: the chart-fetch pattern.
		"window_1h": func(d Device) string {
			return base(d) + "/value?" + url.Values{
				"since": {time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)},
				"limit": {"1000"},
			}.Encode()
		},
		// Interface-root snapshot: the device-overview pattern (nested tree).
		"snapshot": func(d Device) string {
			return base(d)
		},
	}

	res := queryResult{
		Target: st.BaseURL, Devices: len(workers), PerType: *n,
		Concurrency: *concurrency, Latency: map[string]summary{},
	}
	var errs int64
	var errMu sync.Mutex

	for name, mkURL := range types {
		hist := &histogram{}
		jobs := make(chan Device)
		var wg sync.WaitGroup
		for range *concurrency {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for dev := range jobs {
					start := time.Now()
					if err := c.getRaw(mkURL(dev), token); err != nil {
						errMu.Lock()
						errs++
						errMu.Unlock()
						continue
					}
					hist.add(time.Since(start))
				}
			}()
		}
		for range *n {
			jobs <- workers[rand.Intn(len(workers))] //nolint:gosec // benchmark scheduling, not crypto
		}
		close(jobs)
		wg.Wait()
		res.Latency[name] = hist.summarize()
	}
	res.Errors = errs

	fmt.Printf("target      %s (realm %s, %d devices, %d/type, concurrency %d)\n",
		res.Target, st.Realm, res.Devices, res.PerType, res.Concurrency)
	for name, s := range res.Latency {
		fmt.Printf("%-11s %s\n", name, s)
	}
	if errs > 0 {
		fmt.Printf("errors      %d\n", errs)
	}
	if path, err := writeResult(*outDir, "query", res); err != nil {
		return err
	} else if path != "" {
		fmt.Printf("result → %s\n", path)
	}
	return nil
}
