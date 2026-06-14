package stream

import (
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	promtest "github.com/prometheus/client_golang/prometheus/testutil"
)

// dataEvent builds an incoming_data event for one realm/device/interface.
func dataEvent(realm, device, iface string) Event {
	return Event{
		Kind: KindIncomingData, Realm: realm, DeviceID: device,
		Interface: iface, Path: "/v", Value: 1.0, Timestamp: time.Unix(0, 0).UTC(),
	}
}

// recv reads one event from ch within a short deadline.
func recv(t *testing.T, ch <-chan Event) Event {
	t.Helper()
	select {
	case ev := <-ch:
		return ev
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for an event")
		return Event{}
	}
}

// expectEmpty asserts ch holds no buffered event (Publish is synchronous, so
// the channel state is settled once it returns).
func expectEmpty(t *testing.T, ch <-chan Event) {
	t.Helper()
	if n := len(ch); n != 0 {
		t.Fatalf("channel holds %d unexpected events", n)
	}
}

// TestPublishRoutesByRealm: a subscriber receives only its own realm's events.
func TestPublishRoutesByRealm(t *testing.T) {
	b := New(nil)
	defer b.Close()

	ch, cancel := b.Subscribe("realmA", Filter{}, 0)
	defer cancel()

	b.Publish(dataEvent("realmA", "dev1", "com.ex.S"))
	b.Publish(dataEvent("realmB", "dev1", "com.ex.S")) // other realm: dropped

	if ev := recv(t, ch); ev.Realm != "realmA" || ev.DeviceID != "dev1" {
		t.Errorf("got %+v", ev)
	}
	expectEmpty(t, ch)
}

// TestFilter: device and interface filters narrow a subscription; lifecycle
// events (no interface) survive only an interface-less filter.
func TestFilter(t *testing.T) {
	t.Run("by device", func(t *testing.T) {
		b := New(nil)
		defer b.Close()
		ch, cancel := b.Subscribe("r", Filter{DeviceID: "keep"}, 0)
		defer cancel()
		b.Publish(dataEvent("r", "drop", "com.ex.S"))
		b.Publish(dataEvent("r", "keep", "com.ex.S"))
		if ev := recv(t, ch); ev.DeviceID != "keep" {
			t.Errorf("got device %q", ev.DeviceID)
		}
		expectEmpty(t, ch)
	})

	t.Run("by interface", func(t *testing.T) {
		b := New(nil)
		defer b.Close()
		ch, cancel := b.Subscribe("r", Filter{Interface: "com.ex.Keep"}, 0)
		defer cancel()
		b.Publish(dataEvent("r", "d", "com.ex.Drop"))
		b.Publish(dataEvent("r", "d", "com.ex.Keep"))
		if ev := recv(t, ch); ev.Interface != "com.ex.Keep" {
			t.Errorf("got interface %q", ev.Interface)
		}
		expectEmpty(t, ch)
	})

	t.Run("interface filter drops lifecycle events", func(t *testing.T) {
		b := New(nil)
		defer b.Close()
		ch, cancel := b.Subscribe("r", Filter{Interface: "com.ex.S"}, 0)
		defer cancel()
		b.Publish(Event{Kind: KindDeviceConnected, Realm: "r", DeviceID: "d"})
		expectEmpty(t, ch)
	})

	t.Run("empty filter keeps lifecycle events", func(t *testing.T) {
		b := New(nil)
		defer b.Close()
		ch, cancel := b.Subscribe("r", Filter{}, 0)
		defer cancel()
		b.Publish(Event{Kind: KindDeviceConnected, Realm: "r", DeviceID: "d"})
		if ev := recv(t, ch); ev.Kind != KindDeviceConnected {
			t.Errorf("got kind %q", ev.Kind)
		}
	})
}

// TestSlowConsumerDrops: a full subscriber channel drops further events with a
// metric, never blocking the publisher (docs/DESIGN.md §1.4).
func TestSlowConsumerDrops(t *testing.T) {
	reg := prometheus.NewRegistry()
	b := New(reg)
	defer b.Close()

	ch, cancel := b.Subscribe("r", Filter{}, 1) // buffer of exactly 1
	defer cancel()

	b.Publish(dataEvent("r", "d", "i")) // buffered
	b.Publish(dataEvent("r", "d", "i")) // channel full -> dropped
	b.Publish(dataEvent("r", "d", "i")) // dropped

	if got := promtest.ToFloat64(b.dropped); got != 2 {
		t.Errorf("dropped = %v, want 2", got)
	}
	if ev := recv(t, ch); ev.Realm != "r" { // the one buffered event survives
		t.Errorf("got %+v", ev)
	}
	expectEmpty(t, ch)
}

// TestCancelUnsubscribes: cancel removes the subscriber, closes its channel,
// and is safe to call more than once.
func TestCancelUnsubscribes(t *testing.T) {
	b := New(nil)
	defer b.Close()

	ch, cancel := b.Subscribe("r", Filter{}, 0)
	if got := b.Subscribers(); got != 1 {
		t.Fatalf("subscribers = %d, want 1", got)
	}

	cancel()
	cancel() // idempotent

	if got := b.Subscribers(); got != 0 {
		t.Errorf("subscribers after cancel = %d, want 0", got)
	}
	if _, ok := <-ch; ok {
		t.Error("channel should be closed after cancel")
	}
	// A publish after cancel must not panic on the closed channel.
	b.Publish(dataEvent("r", "d", "i"))
}

// TestCloseClosesAllSubscribers: Close closes every subscriber channel; later
// publishes are no-ops and a post-close cancel does not double-close.
func TestCloseClosesAllSubscribers(t *testing.T) {
	b := New(nil)
	ch1, cancel1 := b.Subscribe("r", Filter{}, 0)
	ch2, _ := b.Subscribe("r", Filter{}, 0)

	b.Close()
	b.Close() // idempotent

	if _, ok := <-ch1; ok {
		t.Error("ch1 not closed by Close")
	}
	if _, ok := <-ch2; ok {
		t.Error("ch2 not closed by Close")
	}
	cancel1() // cancel after Close must not panic (once-guarded)
	b.Publish(dataEvent("r", "d", "i"))
}

// TestSubscribeAfterClose returns an already-closed channel and a no-op cancel.
func TestSubscribeAfterClose(t *testing.T) {
	b := New(nil)
	b.Close()

	ch, cancel := b.Subscribe("r", Filter{}, 0)
	if _, ok := <-ch; ok {
		t.Error("channel from a closed bus should be closed")
	}
	cancel() // no-op, must not panic
	if got := b.Subscribers(); got != 0 {
		t.Errorf("subscribers = %d, want 0", got)
	}
}

// TestConcurrentPublishSubscribe exercises the lock discipline under -race:
// publishers, subscribers, and cancels run concurrently without data races.
func TestConcurrentPublishSubscribe(t *testing.T) {
	b := New(nil)
	defer b.Close()

	var wg sync.WaitGroup
	stop := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				b.Publish(dataEvent("r", "d", "i"))
			}
		}
	}()

	wg.Add(8)
	for range 8 {
		go func() {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				ch, cancel := b.Subscribe("r", Filter{}, 4)
				select {
				case <-ch:
				default:
				}
				cancel()
			}
		}()
	}

	time.Sleep(50 * time.Millisecond)
	close(stop)
	wg.Wait()

	if got := b.Subscribers(); got != 0 {
		t.Errorf("subscribers after all cancels = %d, want 0", got)
	}
}
