// Package load is the T5 load/footprint smoke (docs/ROADMAP.md §10 file 9.7).
// It composes the in-process Astrate stack and drives concurrent device traffic,
// asserting the steady-state ingestion guarantees: every acknowledged QoS 1
// publish is persisted (zero loss — PUBACK rides behind the commit), no
// well-formed message is rejected, and end-to-end publish→committed latency
// stays under budget. The scale is env-tunable (LOAD_DEVICES, LOAD_MESSAGES,
// LOAD_P99_MS): the default is a fast CI smoke; the nightly run sets it to the
// §5.4 target of 1000 devices @ 1 msg/s for 15 minutes against the binary, where
// the RSS/Postgres budgets are enforced against the isolated process.
package load

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/payload"
	"github.com/astrate-platform/astrate/test/conformance/instance"
)

const sensor = "org.astrate.load.Sensor"

var defs = map[string]string{
	sensor: `{"interface_name":"` + sensor + `","version_major":1,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`,
}

func TestLoadSmoke(t *testing.T) {
	devices := envInt("LOAD_DEVICES", 10)
	perDevice := envInt("LOAD_MESSAGES", 20)

	reg := prometheus.NewRegistry()
	in := instance.New(t, instance.Config{Interfaces: defs, Registerer: reg})
	ctx := context.Background()
	iface := in.Interfaces[sensor]

	// Connect and introspect every device (sequential — uses t.Fatal).
	devs := make([]*testutil.AstarteDevice, devices)
	ids := make([]deviceid.ID, devices)
	for i := range devices {
		d, id := in.NewDevice(t, "")
		d.PublishIntrospection(t, testutil.Introspection(map[string][2]int{sensor: {1, 0}}))
		devs[i], ids[i] = d, id
	}
	for _, id := range ids {
		waitFor(t, 15*time.Second, func() bool {
			dev, err := in.Store.GetDevice(ctx, in.Realm.ID, id)
			return err == nil && len(dev.Introspection) == 1
		})
	}

	// Publish concurrently; a QoS 1 publish is PUBACK'd only after the batch
	// commits, so each publish's round-trip time is its end-to-end latency.
	type result struct {
		lat time.Duration
		err error
	}
	results := make(chan result, devices*perDevice)
	var wg sync.WaitGroup
	startAll := time.Now()
	for i := range devices {
		wg.Add(1)
		go func(d *testutil.AstarteDevice) {
			defer wg.Done()
			topic := d.DataTopic(sensor, "/value")
			for j := range perDevice {
				ts := time.Now().UTC()
				body, err := payload.Encode(float64(j), &ts, payload.FormatBSON)
				if err != nil {
					results <- result{err: err}
					continue
				}
				start := time.Now()
				tok := d.Client.Publish(topic, 1, false, body)
				switch {
				case !tok.WaitTimeout(30 * time.Second):
					results <- result{err: fmt.Errorf("publish timed out")}
				case tok.Error() != nil:
					results <- result{err: tok.Error()}
				default:
					results <- result{lat: time.Since(start)}
				}
			}
		}(devs[i])
	}
	wg.Wait()
	close(results)
	elapsed := time.Since(startAll)

	var lats []time.Duration
	var failures int
	for r := range results {
		if r.err != nil {
			failures++
			continue
		}
		lats = append(lats, r.lat)
	}
	if failures > 0 {
		t.Fatalf("%d/%d publishes failed", failures, devices*perDevice)
	}

	total := devices * perDevice
	// Zero loss: every acknowledged QoS 1 publish is committed.
	waitFor(t, 30*time.Second, func() bool { return totalRows(t, in, iface.ID, ids) == total })

	// Zero validation rejects for well-formed traffic.
	if rj := counterTotal(t, reg, "astrate_engine_rejects_total"); rj != 0 {
		t.Errorf("validation rejects = %v, want 0", rj)
	}

	p99 := percentile(lats)
	t.Logf("load: %d devices x %d msgs = %d in %v (%.0f msg/s), p99 publish->commit = %v",
		devices, perDevice, total, elapsed.Round(time.Millisecond),
		float64(total)/elapsed.Seconds(), p99.Round(time.Millisecond))

	// The §5.4 latency budget is meaningful only on an isolated host, so it is
	// asserted on demand (the nightly run sets LOAD_P99_MS=250 against the
	// binary). A shared CI/dev host is too contended for a stable bound, so the
	// smoke leaves the correctness guarantees above as the gate and only logs
	// the latency.
	if v := os.Getenv("LOAD_P99_MS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && p99 > time.Duration(n)*time.Millisecond {
			t.Errorf("ingest p99 = %v, want < %dms", p99, n)
		}
	}
}

func totalRows(t *testing.T, in *instance.Instance, ifaceID int64, ids []deviceid.ID) int {
	t.Helper()
	n := 0
	for _, id := range ids {
		rows, err := in.Store.Series(context.Background(), store.SeriesQuery{
			RealmID: in.Realm.ID, DeviceID: id, InterfaceID: ifaceID, Path: "/value",
		})
		if err != nil {
			t.Fatalf("Series: %v", err)
		}
		n += len(rows)
	}
	return n
}

// counterTotal sums a counter (across labels) from the registry.
func counterTotal(t *testing.T, reg *prometheus.Registry, name string) float64 {
	t.Helper()
	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("Gather: %v", err)
	}
	var total float64
	for _, fam := range families {
		if fam.GetName() != name {
			continue
		}
		for _, m := range fam.GetMetric() {
			if c := m.GetCounter(); c != nil {
				total += c.GetValue()
			}
		}
	}
	return total
}

func percentile(lats []time.Duration) time.Duration {
	if len(lats) == 0 {
		return 0
	}
	sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })
	idx := int(float64(len(lats)) * 0.99)
	if idx >= len(lats) {
		idx = len(lats) - 1
	}
	return lats[idx]
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}

func waitFor(t *testing.T, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("condition not met within timeout")
}
