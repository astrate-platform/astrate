package engine

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	promtest "github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// TestOrderingSingleDevice is the docs/ROADMAP.md §7.3 T1 ordering property
// test: N messages of one device through a multi-shard engine with a
// synthetic slow store land persisted in publish order.
func TestOrderingSingleDevice(t *testing.T) {
	fs := newFakeStore()
	seedAlpha(t, fs)
	fs.mu.Lock()
	fs.appendDelay = 2 * time.Millisecond // synthetic slow store
	fs.mu.Unlock()

	e := startTestEngine(t, fs, Config{
		Shards: 4, ShardQueue: 1024, BatchMaxRows: 8, BatchMaxWait: 2 * time.Millisecond,
	})

	const n = 300
	acks := make([]*ackCounter, n)
	for i := range n {
		acks[i] = &ackCounter{}
		body := enc(t, float64(i), nil, payload.FormatBSON)
		e.Submit(deviceMsg("com.astrate.test.Minimal", "/value", 0, body, acks[i]))
	}

	waitFor(t, 15*time.Second, "all rows persisted", func() bool {
		return len(fs.individualRows()) == n
	})
	rows := fs.individualRows()
	for i, r := range rows {
		if r.ValueDouble == nil || *r.ValueDouble != float64(i) {
			t.Fatalf("row %d out of order: %+v", i, r)
		}
		if r.InterfaceID != ifaceMinimal || r.Path != "/value" {
			t.Fatalf("row %d misrouted: %+v", i, r)
		}
	}
	for i, a := range acks {
		if !a.acked() {
			t.Fatalf("message %d never acknowledged", i)
		}
	}
}

// TestCrossDeviceParallelism: devices on different shards flush
// concurrently (docs/ROADMAP.md §7.3 "cross-device parallelism observed").
func TestCrossDeviceParallelism(t *testing.T) {
	const shards = 4
	devB := devAlpha
	for shardOf(devB[:], shards) == shardOf(devAlpha[:], shards) {
		devB[15]++
	}

	fs := newFakeStore()
	seedAlpha(t, fs)
	fs.addDevice(realmAlpha, realmAlphaID, devB, map[string]store.InterfaceVersion{
		"com.astrate.test.Minimal": {Major: 0, Minor: 1},
	}, "")
	fs.mu.Lock()
	fs.appendDelay = 20 * time.Millisecond
	fs.mu.Unlock()

	e := startTestEngine(t, fs, Config{Shards: shards, BatchMaxRows: 1})

	const perDevice = 5
	ack := &ackCounter{}
	for i := range perDevice {
		body := enc(t, float64(i), nil, payload.FormatBSON)
		e.Submit(deviceMsgFor(devAlpha, "com.astrate.test.Minimal", "/value", 0, body, ack))
		e.Submit(deviceMsgFor(devB, "com.astrate.test.Minimal", "/value", 0, body, ack))
	}
	waitFor(t, 10*time.Second, "all rows persisted", func() bool {
		return len(fs.individualRows()) == 2*perDevice
	})
	if fs.maxInFlightNow() < 2 {
		t.Errorf("max concurrent flushes = %d, want >= 2 (no cross-device parallelism)", fs.maxInFlightNow())
	}
}

// TestBackpressureQoS1Blocks: a filled shard blocks QoS >= 1 submits — the
// deferred-ack backpressure of docs/DESIGN.md §1.4 — and releases them once
// the store recovers, acknowledgments strictly after commit.
func TestBackpressureQoS1Blocks(t *testing.T) {
	fs := newFakeStore()
	seedAlpha(t, fs)
	gate := make(chan struct{})
	fs.mu.Lock()
	fs.appendGate = gate
	fs.mu.Unlock()

	e := startTestEngine(t, fs, Config{
		Shards: 1, ShardQueue: 2, BatchMaxRows: 1, BatchMaxWait: time.Hour,
	})
	var gateOnce sync.Once
	openGate := func() { gateOnce.Do(func() { close(gate) }) }
	t.Cleanup(openGate) // runs before the engine drain cleanup (LIFO)

	const n = 5
	acks := make([]*ackCounter, n)
	var progress atomic.Int32
	go func() {
		for i := range n {
			acks[i] = &ackCounter{}
			body := enc(t, float64(i), nil, payload.FormatBSON)
			e.Submit(deviceMsg("com.astrate.test.Minimal", "/value", 1, body, acks[i]))
			progress.Store(int32(i + 1))
		}
	}()

	// msg0 is taken by the shard and its flush blocks on the gate; msgs 1-2
	// fill the queue; the 4th submit must block.
	waitFor(t, 5*time.Second, "shard blocked in flush", func() bool {
		return fs.inFlightNow() == 1 && len(e.shards[0].ch) == 2 && progress.Load() == 3
	})
	time.Sleep(100 * time.Millisecond)
	if got := progress.Load(); got != 3 {
		t.Fatalf("submit progressed to %d with a full shard, want blocked at 3", got)
	}
	for i := range 3 {
		if acks[i].acked() {
			t.Fatalf("message %d acknowledged before its batch committed", i)
		}
	}

	openGate()
	waitFor(t, 5*time.Second, "all messages persisted and acked", func() bool {
		if len(fs.individualRows()) != n {
			return false
		}
		for i := range n {
			if acks[i] == nil || !acks[i].acked() {
				return false
			}
		}
		return true
	})
	rows := fs.individualRows()
	for i, r := range rows {
		if *r.ValueDouble != float64(i) {
			t.Fatalf("row %d out of order after backpressure release: %v", i, *r.ValueDouble)
		}
	}
}

// TestBackpressureQoS0Drops: a full shard drops QoS 0 messages with a metric
// instead of blocking (docs/DESIGN.md §1.4).
func TestBackpressureQoS0Drops(t *testing.T) {
	fs := newFakeStore()
	seedAlpha(t, fs)
	gate := make(chan struct{})
	fs.mu.Lock()
	fs.appendGate = gate
	fs.mu.Unlock()

	e := startTestEngine(t, fs, Config{
		Shards: 1, ShardQueue: 2, BatchMaxRows: 1, BatchMaxWait: time.Hour,
	})
	var gateOnce sync.Once
	openGate := func() { gateOnce.Do(func() { close(gate) }) }
	t.Cleanup(openGate)

	ack := &ackCounter{}
	submit := func(v float64) {
		e.Submit(deviceMsg("com.astrate.test.Minimal", "/value", 0, enc(t, v, nil, payload.FormatBSON), ack))
	}
	submit(0)
	waitFor(t, 5*time.Second, "shard blocked in flush", func() bool {
		return fs.inFlightNow() == 1
	})
	submit(1)
	submit(2) // queue now full
	for i := range 3 {
		submit(float64(10 + i)) // dropped
	}
	if got := promtest.ToFloat64(e.met.qos0Drops); got != 3 {
		t.Errorf("qos0 drop counter = %v, want 3", got)
	}

	openGate()
	waitFor(t, 5*time.Second, "accepted rows persisted", func() bool {
		return len(fs.individualRows()) == 3
	})
	time.Sleep(50 * time.Millisecond)
	if got := len(fs.individualRows()); got != 3 {
		t.Errorf("%d rows persisted, want 3 (dropped messages must not appear)", got)
	}
}

// TestDrain: graceful shutdown processes queued messages, flushes the final
// batch with acknowledgment, and refuses (without acking) later submits
// (docs/DESIGN.md §5.3).
func TestDrain(t *testing.T) {
	fs := newFakeStore()
	seedAlpha(t, fs)
	// Batches flush only on drain: row cap and wait are both unreachable.
	e := startTestEngine(t, fs, Config{Shards: 2, BatchMaxRows: 1000, BatchMaxWait: time.Hour})

	const n = 3
	acks := make([]*ackCounter, n)
	for i := range n {
		acks[i] = &ackCounter{}
		body := enc(t, float64(i), nil, payload.FormatBSON)
		e.Submit(deviceMsg("com.astrate.test.Minimal", "/value", 1, body, acks[i]))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.drain(ctx); err != nil {
		t.Fatalf("drain: %v", err)
	}
	if got := len(fs.individualRows()); got != n {
		t.Fatalf("%d rows after drain, want %d", got, n)
	}
	for i := range n {
		if !acks[i].acked() {
			t.Errorf("message %d not acknowledged by drain flush", i)
		}
	}

	late := &ackCounter{}
	e.Submit(deviceMsg("com.astrate.test.Minimal", "/value", 1,
		enc(t, 9.0, nil, payload.FormatBSON), late))
	if late.acked() {
		t.Error("post-drain submit was acknowledged")
	}
	if got := promtest.ToFloat64(e.met.droppedShutdown); got != 1 {
		t.Errorf("shutdown drop counter = %v, want 1", got)
	}
	if err := e.drain(ctx); err != nil {
		t.Errorf("second drain not idempotent: %v", err)
	}
}

// TestShardPanicRecovery: a panic inside a shard is recovered and logged,
// and the restarted loop retries the pending batch (docs/ROADMAP.md §7.1
// file 6.2).
func TestShardPanicRecovery(t *testing.T) {
	fs := newFakeStore()
	seedAlpha(t, fs)
	fs.mu.Lock()
	fs.appendPanics = 1
	fs.mu.Unlock()

	e := startTestEngine(t, fs, Config{Shards: 1, BatchMaxRows: 1, BatchMaxWait: 5 * time.Millisecond})

	ack := &ackCounter{}
	e.Submit(deviceMsg("com.astrate.test.Minimal", "/value", 1,
		enc(t, 1.5, nil, payload.FormatBSON), ack))

	waitFor(t, 5*time.Second, "message persisted after panic restart", func() bool {
		return ack.acked() && len(fs.individualRows()) == 1
	})
	if got := promtest.ToFloat64(e.met.internalErrors); got < 1 {
		t.Errorf("internal error counter = %v, want >= 1 (recovered panic)", got)
	}

	// The shard must still be alive for subsequent traffic.
	again := &ackCounter{}
	e.Submit(deviceMsg("com.astrate.test.Minimal", "/value", 1,
		enc(t, 2.5, nil, payload.FormatBSON), again))
	waitFor(t, 5*time.Second, "post-panic message persisted", func() bool {
		return again.acked()
	})
}

// TestInvalidation: the compiled-interface snapshot follows both
// invalidation paths — LISTEN/NOTIFY and the in-process callback
// (docs/DESIGN.md §2.6, docs/ROADMAP.md §7.3).
func TestInvalidation(t *testing.T) {
	fs := newFakeStore()
	fs.addRealm(1, "alpha")
	e := startTestEngine(t, fs, Config{Shards: 1})

	if e.schemas.realm("alpha").iface("com.astrate.test.Minimal", 0) != nil {
		t.Fatal("interface resolvable before install")
	}

	// Path 1: LISTEN/NOTIFY.
	fs.addInterface(1, fixtureStored(t, 1, 15, "com.astrate.test.Minimal.json"))
	fs.notifyCh <- store.Notification{Channel: store.ChannelInterfaces, Payload: "1"}
	waitFor(t, 5*time.Second, "notify-driven snapshot rebuild", func() bool {
		return e.schemas.realm("alpha").iface("com.astrate.test.Minimal", 0) != nil
	})

	// Path 2: in-process callback, synchronous.
	fs.addInterface(1, fixtureStored(t, 1, 10, "com.astrate.test.AllScalarTypes.json"))
	if err := e.RefreshInterfaces(context.Background(), 1); err != nil {
		t.Fatalf("RefreshInterfaces: %v", err)
	}
	if e.schemas.realm("alpha").iface("com.astrate.test.AllScalarTypes", 1) == nil {
		t.Error("callback-driven rebuild missed the new interface")
	}

	// A malformed notification payload degrades to a full reload.
	fs.addRealm(2, "beta")
	fs.notifyCh <- store.Notification{Channel: store.ChannelInterfaces, Payload: "bogus"}
	waitFor(t, 5*time.Second, "full reload on malformed payload", func() bool {
		return e.schemas.realm("beta") != nil
	})
}

// TestLifecycleEviction: device disconnects evict the per-device cache; the
// M6b seam observes every event.
func TestLifecycleEviction(t *testing.T) {
	ctx := context.Background()
	fs := newFakeStore()
	seedAlpha(t, fs)
	e := newTestEngine(t, fs, Config{Shards: 1})

	var seen []broker.LifecycleEventType
	e.onLifecycle = func(ev broker.LifecycleEvent) { seen = append(seen, ev.Type) }

	if _, err := e.devices.get(ctx, realmAlpha, realmAlphaID, devAlpha); err != nil {
		t.Fatalf("priming device cache: %v", err)
	}
	e.OnLifecycleEvent(broker.LifecycleEvent{
		Type: broker.EventDeviceConnected, Realm: realmAlpha, DeviceID: devAlpha,
	})
	if e.devices.len() != 1 {
		t.Error("connect event evicted the device cache")
	}
	e.OnLifecycleEvent(broker.LifecycleEvent{
		Type: broker.EventDeviceDisconnected, Realm: realmAlpha, DeviceID: devAlpha,
	})
	if e.devices.len() != 0 {
		t.Error("disconnect event did not evict the device cache")
	}
	if len(seen) != 2 || seen[0] != broker.EventDeviceConnected || seen[1] != broker.EventDeviceDisconnected {
		t.Errorf("lifecycle seam saw %v", seen)
	}
}

// TestShardOf: routing is deterministic and in range.
func TestShardOf(t *testing.T) {
	for n := 1; n <= 32; n *= 2 {
		id := deviceid.ID{}
		for i := range 64 {
			id[i%16] ^= byte(i*31 + 7)
			s := shardOf(id[:], n)
			if s != shardOf(id[:], n) {
				t.Fatalf("shardOf not deterministic for %v", id)
			}
			if s < 0 || s >= n {
				t.Fatalf("shardOf(%v, %d) = %d out of range", id, n, s)
			}
		}
	}
}
