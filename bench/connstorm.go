package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

// connstormResult is the JSON document one connstorm run produces.
type connstormResult struct {
	Target      string  `json:"target"`
	Devices     int     `json:"devices"`
	Concurrency int     `json:"concurrency"`
	CertPhaseS  float64 `json:"cert_issuance_phase_s"` // pairing CA throughput, reported separately
	StormPhaseS float64 `json:"connect_phase_s"`
	ConnectRate float64 `json:"connects_per_s"`
	Failures    int64   `json:"failures"`
	Connack     summary `json:"connack_latency"`
}

// cmdConnstorm mass-connects pre-registered devices over mTLS (including the
// introspection handshake) and times it. Certificate issuance runs first,
// untimed against the storm, and is reported as its own phase — it exercises
// the pairing CA rather than the broker.
func cmdConnstorm(args []string) error {
	fs := flag.NewFlagSet("connstorm", flag.ExitOnError)
	statePath := fs.String("state", "bench-state.json", "state file from provision")
	devices := fs.Int("devices", 200, "devices to connect")
	concurrency := fs.Int("concurrency", 64, "parallel connects")
	timeout := fs.Duration("timeout", 30*time.Second, "per-connect timeout")
	hold := fs.Duration("hold", 5*time.Second, "how long to hold all sessions open before disconnecting")
	brokerOverride := fs.String("broker-url", "", "override the stored broker URL")
	useRSA := fs.Bool("rsa", false, "use RSA-2048 device keys instead of ECDSA P-256")
	outDir := fs.String("out", "results", "directory for the result JSON ('' to disable)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	st, err := loadState(*statePath)
	if err != nil {
		return err
	}
	probes, workers, err := st.workload(*devices - st.Probes)
	if err != nil {
		return err
	}
	targets := append(append([]Device(nil), probes...), workers...)
	broker := st.BrokerURL
	if *brokerOverride != "" {
		broker = *brokerOverride
	}
	c := newClient(30 * time.Second)
	introspection := map[string]string{ifaceIndividual: "1:0", ifaceObject: "1:0"}

	// Phase 1: certificate issuance (pairing CA), bounded by the same
	// concurrency but timed separately.
	fmt.Printf("issuing %d certificates…\n", len(targets))
	idents := make([]*mqttDevice, len(targets))
	certStart := time.Now()
	if err := forEachLimited(targets, *concurrency, func(i int, dev Device) error {
		d, err := newIdentity(c, st.Endpoints, st.Realm, dev, *useRSA, st.TLSSkipVerify)
		idents[i] = d
		return err
	}); err != nil {
		return fmt.Errorf("certificate phase: %w", err)
	}
	certPhase := time.Since(certStart)

	// Phase 2: the storm.
	fmt.Printf("connecting %d devices with concurrency %d…\n", len(idents), *concurrency)
	connHist := &histogram{}
	var failures int64
	var mu sync.Mutex
	stormStart := time.Now()
	_ = forEachLimited(targets, *concurrency, func(i int, _ Device) error {
		lat, err := idents[i].connect(broker, *timeout, introspection)
		if err != nil {
			mu.Lock()
			failures++
			mu.Unlock()
			return nil // keep storming; failures are the metric, not an abort
		}
		connHist.add(lat)
		return nil
	})
	stormPhase := time.Since(stormStart)

	time.Sleep(*hold)
	for _, d := range idents {
		if d != nil {
			d.disconnect()
		}
	}

	res := connstormResult{
		Target: st.BaseURL, Devices: len(idents), Concurrency: *concurrency,
		CertPhaseS: certPhase.Seconds(), StormPhaseS: stormPhase.Seconds(),
		ConnectRate: float64(int64(len(idents))-failures) / stormPhase.Seconds(),
		Failures:    failures, Connack: connHist.summarize(),
	}
	fmt.Printf(`
target      %s
devices     %d (concurrency %d)
cert phase  %.2fs (%.1f certs/s)
storm       %.2fs → %.1f connects/s, %d failures
connack     %s
`, res.Target, res.Devices, res.Concurrency,
		res.CertPhaseS, float64(res.Devices)/res.CertPhaseS,
		res.StormPhaseS, res.ConnectRate, res.Failures, res.Connack)
	if path, err := writeResult(*outDir, "connstorm", res); err != nil {
		return err
	} else if path != "" {
		fmt.Printf("result → %s\n", path)
	}
	return nil
}

// forEachLimited runs fn over items with at most limit in flight and returns
// the first error.
func forEachLimited(items []Device, limit int, fn func(int, Device) error) error {
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		firstErr error
	)
	sem := make(chan struct{}, limit)
	for i, it := range items {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			if err := fn(i, it); err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return firstErr
}
