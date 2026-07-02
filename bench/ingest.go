package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ingestResult is the JSON document one ingest run produces.
type ingestResult struct {
	Target        string  `json:"target"`
	Realm         string  `json:"realm"`
	Devices       int     `json:"devices"`
	RatePerDevice float64 `json:"rate_per_device"`
	Payload       string  `json:"payload"`
	DurationS     float64 `json:"duration_s"`

	Sent         int64   `json:"sent"`
	Acked        int64   `json:"acked"`
	PublishFails int64   `json:"publish_fails"`
	Behind       int64   `json:"behind_ticks"` // ticks skipped because the previous PUBACK was still pending
	AchievedRate float64 `json:"achieved_rate_total"`

	Connect summary `json:"connack_latency"`
	Puback  summary `json:"puback_latency"` // see README: NOT comparable across platforms
	E2E     summary `json:"e2e_visibility_latency"`

	LossExpected int64 `json:"loss_check_expected"`
	LossFound    int64 `json:"loss_check_found"`
	LossMissing  int64 `json:"loss_check_missing"`
	LossSkipped  bool  `json:"loss_check_skipped"`
}

// cmdIngest connects N workload devices, publishes at a fixed per-device rate
// for the run duration, and measures CONNACK latency, PUBACK latency,
// publish→API-visible e2e latency (via dedicated probe devices), and loss
// (row count vs sent count, verified through the AppEngine API).
func cmdIngest(args []string) error {
	fs := flag.NewFlagSet("ingest", flag.ExitOnError)
	statePath := fs.String("state", "bench-state.json", "state file from provision")
	devices := fs.Int("devices", 100, "workload devices to run")
	rate := fs.Float64("rate", 1.0, "messages per second per device")
	duration := fs.Duration("duration", 60*time.Second, "publish duration")
	payload := fs.String("payload", "individual", "payload shape: individual | object")
	probeInterval := fs.Duration("probe-poll", 25*time.Millisecond, "e2e probe poll interval (quantizes e2e latency)")
	connConc := fs.Int("connect-concurrency", 32, "parallel device connects during setup")
	grace := fs.Duration("grace", 10*time.Second, "settle time between last publish and the loss check")
	brokerOverride := fs.String("broker-url", "", "override the stored broker URL")
	useRSA := fs.Bool("rsa", false, "use RSA-2048 device keys instead of ECDSA P-256")
	noLoss := fs.Bool("no-loss-check", false, "skip the end-of-run loss verification")
	outDir := fs.String("out", "results", "directory for the result JSON ('' to disable)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *payload != "individual" && *payload != "object" {
		return fmt.Errorf("-payload must be individual or object")
	}

	st, err := loadState(*statePath)
	if err != nil {
		return err
	}
	realmKey, err := parsePrivateKeyPEM([]byte(st.RealmKeyPEM))
	if err != nil {
		return fmt.Errorf("realm key from state: %w", err)
	}
	probes, workers, err := st.workload(*devices)
	if err != nil {
		return err
	}
	broker := st.BrokerURL
	if *brokerOverride != "" {
		broker = *brokerOverride
	}

	c := newClient(30 * time.Second)
	introspection := map[string]string{ifaceIndividual: "1:0", ifaceObject: "1:0"}

	// --- setup: identities + connects (bounded concurrency) ---------------
	all := append(append([]Device(nil), probes...), workers...)
	conns := make([]*mqttDevice, len(all))
	connHist := &histogram{}
	var setupErr error
	{
		var wg sync.WaitGroup
		var mu sync.Mutex
		sem := make(chan struct{}, *connConc)
		for i, dev := range all {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				d, err := newIdentity(c, st.Endpoints, st.Realm, dev, *useRSA, st.TLSSkipVerify)
				if err == nil {
					var lat time.Duration
					lat, err = d.connect(broker, 30*time.Second, introspection)
					if err == nil {
						connHist.add(lat)
						conns[i] = d
					}
				}
				if err != nil {
					mu.Lock()
					if setupErr == nil {
						setupErr = err
					}
					mu.Unlock()
				}
			}()
		}
		wg.Wait()
	}
	defer func() {
		for _, d := range conns {
			if d != nil {
				d.disconnect()
			}
		}
	}()
	if setupErr != nil {
		return fmt.Errorf("device setup: %w", setupErr)
	}
	probeConns, workerConns := conns[:len(probes)], conns[len(probes):]
	fmt.Printf("connected %d devices (%s), running %v at %.2f msg/s/device…\n",
		len(conns), connHist.summarize(), *duration, *rate)

	// --- run ----------------------------------------------------------------
	runStart := time.Now()
	ctx, cancel := context.WithDeadline(context.Background(), runStart.Add(*duration))
	defer cancel()

	var sent, acked, fails, behind atomic.Int64
	pubHist := &histogram{}
	sentPerDevice := make([]int64, len(workerConns))

	var wg sync.WaitGroup
	interval := time.Duration(float64(time.Second) / *rate)
	for wi, d := range workerConns {
		wg.Add(1)
		go func() {
			defer wg.Done()
			next := time.Now()
			var seq float64
			for {
				next = next.Add(interval)
				now := time.Now()
				if wait := next.Sub(now); wait > 0 {
					select {
					case <-ctx.Done():
						return
					case <-time.After(wait):
					}
				} else if -wait > interval {
					// The previous PUBACK held us past a full tick: the
					// closed loop is falling behind the requested rate.
					behind.Add(1)
					next = now
				}
				if ctx.Err() != nil {
					return
				}

				seq++
				ts := time.Now()
				var tok interface {
					WaitTimeout(time.Duration) bool
					Error() error
				}
				var err error
				if *payload == "object" {
					tok, err = d.publishObject(ifaceObject, "/data",
						bson.M{"temp": 20 + seq/1000, "hum": 40.0, "status": "ok"}, ts)
				} else {
					tok, err = d.publishDouble(ifaceIndividual, "/value", seq, ts)
				}
				if err != nil {
					fails.Add(1)
					continue
				}
				sent.Add(1)
				pubStart := time.Now()
				if !tok.WaitTimeout(30*time.Second) || tok.Error() != nil {
					fails.Add(1)
					continue
				}
				pubHist.add(time.Since(pubStart))
				acked.Add(1)
				sentPerDevice[wi]++
			}
		}()
	}

	// --- e2e probes: one marker in flight per probe device ------------------
	e2eHist := &histogram{}
	aeaToken, err := mintJWT(realmKey, "a_aea")
	if err != nil {
		return err
	}
	probeSent := make([]int64, len(probeConns))
	for pi, d := range probeConns {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ctx.Err() == nil {
				cycle := time.Now()
				marker := float64(cycle.UnixNano())
				tok, err := d.publishDouble(ifaceIndividual, "/value", marker, cycle)
				if err != nil || !tok.WaitTimeout(30*time.Second) || tok.Error() != nil {
					continue
				}
				probeSent[pi]++
				deadline := time.Now().Add(15 * time.Second)
				for time.Now().Before(deadline) && ctx.Err() == nil {
					rows, err := c.getSamples(st.Endpoints, aeaToken, st.Realm, d.id,
						ifaceIndividual, "/value", url.Values{"limit": {"1"}})
					if err == nil && len(rows) == 1 {
						if v, perr := rows[0].float64Value(); perr == nil && v == marker {
							e2eHist.add(time.Since(cycle))
							break
						}
					}
					time.Sleep(*probeInterval)
				}
				// One marker cycle per second keeps the probe load constant.
				if rest := time.Second - time.Since(cycle); rest > 0 {
					select {
					case <-ctx.Done():
					case <-time.After(rest):
					}
				}
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(runStart)

	// --- loss check ----------------------------------------------------------
	res := ingestResult{
		Target: st.BaseURL, Realm: st.Realm, Devices: len(workerConns),
		RatePerDevice: *rate, Payload: *payload, DurationS: elapsed.Seconds(),
		Sent: sent.Load(), Acked: acked.Load(), PublishFails: fails.Load(),
		Behind:       behind.Load(),
		AchievedRate: float64(acked.Load()) / elapsed.Seconds(),
		Connect:      connHist.summarize(),
		Puback:       pubHist.summarize(),
		E2E:          e2eHist.summarize(),
		LossSkipped:  *noLoss,
	}
	if !*noLoss {
		fmt.Printf("run done (%d acked), waiting %v before the loss check…\n", res.Acked, *grace)
		time.Sleep(*grace)
		// Object aggregates are stored one row per object at the object path;
		// counting must query /data, not a leaf.
		iface, path := ifaceIndividual, "/value"
		if *payload == "object" {
			iface, path = ifaceObject, "/data"
		}
		since := runStart.Add(-2 * time.Second)
		for wi, d := range workerConns {
			found, err := c.countSamples(st.Endpoints, aeaToken, st.Realm, d.id, iface, path, since, sentPerDevice[wi])
			if err != nil {
				return fmt.Errorf("loss check on %s: %w", d.id, err)
			}
			res.LossExpected += sentPerDevice[wi]
			res.LossFound += found
		}
		// Probe markers are rows too; count them against the same window.
		for pi, d := range probeConns {
			found, err := c.countSamples(st.Endpoints, aeaToken, st.Realm, d.id, ifaceIndividual, "/value", since, probeSent[pi])
			if err != nil {
				return fmt.Errorf("loss check on probe %s: %w", d.id, err)
			}
			res.LossExpected += probeSent[pi]
			res.LossFound += found
		}
		res.LossMissing = res.LossExpected - res.LossFound
	}

	// --- report ---------------------------------------------------------------
	fmt.Printf(`
target      %s (realm %s)
workload    %d devices × %.2f msg/s, payload=%s, %.1fs
sent/acked  %d / %d (fails %d, behind-ticks %d) → %.1f msg/s aggregate
connack     %s
puback      %s   (semantics differ per platform — see README)
e2e         %s   (poll quantization ±%v)
`, res.Target, res.Realm, res.Devices, res.RatePerDevice, res.Payload, res.DurationS,
		res.Sent, res.Acked, res.PublishFails, res.Behind, res.AchievedRate,
		res.Connect, res.Puback, res.E2E, *probeInterval)
	if !*noLoss {
		fmt.Printf("loss        %d expected, %d found, %d missing\n", res.LossExpected, res.LossFound, res.LossMissing)
	}
	if path, err := writeResult(*outDir, "ingest", res); err != nil {
		return err
	} else if path != "" {
		fmt.Printf("result → %s\n", path)
	}
	return nil
}

// countSamples counts rows on a concrete path newer than since. Single
// request when the expectation fits one page; otherwise windowed advance on
// since_after (ascending — Astrate honors sort=ascending, Astarte's since
// pagination is ascending already and ignores the extra parameter).
func (c *client) countSamples(ep Endpoints, token, realm, device, iface, path string, since time.Time, expected int64) (int64, error) {
	const page = 9000
	if expected <= page {
		rows, err := c.getSamples(ep, token, realm, device, iface, path, url.Values{
			"since": {since.UTC().Format(time.RFC3339Nano)},
			"limit": {strconv.FormatInt(expected+100, 10)},
			"sort":  {"ascending"},
		})
		if err != nil {
			return 0, err
		}
		return int64(len(rows)), nil
	}
	var total int64
	cursor := url.Values{
		"since": {since.UTC().Format(time.RFC3339Nano)},
		"limit": {strconv.Itoa(page)},
		"sort":  {"ascending"},
	}
	for {
		rows, err := c.getSamples(ep, token, realm, device, iface, path, cursor)
		if err != nil {
			return total, err
		}
		total += int64(len(rows))
		if len(rows) < page {
			return total, nil
		}
		last := rows[len(rows)-1].Timestamp
		cursor = url.Values{
			"since_after": {last},
			"limit":       {strconv.Itoa(page)},
			"sort":        {"ascending"},
		}
	}
}
