package engine

import (
	"context"
	"net/netip"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// busRig is a pipeline rig with the live bus wired up. newPipelineRig builds
// its engine through newEngine, which — unlike New — leaves bus, afterCommit
// and onLifecycle unset, so a test that wants bus events must supply them.
type busRig struct {
	*pipelineRig
	events <-chan stream.Event
}

func newBusRig(t *testing.T) *busRig {
	t.Helper()
	rig, _ := newPipelineRig(t, Config{})
	rig.e.bus = stream.New(nil)
	rig.e.afterCommit = rig.e.fireCommitted
	rig.e.broker = &fakePort{} // handleLifecycle's connect case publishes consumer/properties
	t.Cleanup(func() {
		rig.e.bg.Wait()
		rig.e.bus.Close()
	})
	events, cancel := rig.e.bus.Subscribe(realmAlpha, stream.Filter{}, 8)
	t.Cleanup(cancel)
	return &busRig{pipelineRig: rig, events: events}
}

// next returns the next bus event, failing the test if none arrives.
func (r *busRig) next(t *testing.T) stream.Event {
	t.Helper()
	select {
	case ev, ok := <-r.events:
		if !ok {
			t.Fatal("bus closed before an event arrived")
		}
		return ev
	case <-time.After(2 * time.Second):
		t.Fatal("no bus event within 2s")
		return stream.Event{}
	}
}

// TestBusEventDetail pins the fields a Channels viewer rebuilds upstream
// trigger-event bodies from: they are all available at the publish site, and
// each one was previously dropped there.
func TestBusEventDetail(t *testing.T) {
	t.Run("IncomingDataCarriesMajor", func(t *testing.T) {
		rig := newBusRig(t)
		ack := &ackCounter{}
		rig.handle(deviceMsg("com.astrate.test.AllScalarTypes", "/double", 0,
			enc(t, 22.5, nil, payload.FormatBSON), ack))
		// afterCommit fires from the batch, not from handle.
		rig.sh.batch.flush(context.Background())

		ev := rig.next(t)
		if ev.Kind != stream.KindIncomingData {
			t.Fatalf("Kind = %q, want %q", ev.Kind, stream.KindIncomingData)
		}
		if ev.Interface != "com.astrate.test.AllScalarTypes" {
			t.Errorf("Interface = %q", ev.Interface)
		}
		if ev.InterfaceMajor != 1 {
			t.Errorf("InterfaceMajor = %d, want 1", ev.InterfaceMajor)
		}
	})

	t.Run("DeviceConnectedCarriesIP", func(t *testing.T) {
		rig := newBusRig(t)
		at := time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC)
		rig.e.handleLifecycle(broker.LifecycleEvent{
			Type: broker.EventDeviceConnected, Realm: realmAlpha, DeviceID: devAlpha,
			RemoteIP: netip.MustParseAddr("203.0.113.7"), At: at,
		})

		ev := rig.next(t)
		if ev.Kind != stream.KindDeviceConnected {
			t.Fatalf("Kind = %q, want %q", ev.Kind, stream.KindDeviceConnected)
		}
		if ev.IP != "203.0.113.7" {
			t.Errorf("IP = %q, want 203.0.113.7", ev.IP)
		}
	})

	t.Run("DeviceErrorIsPublished", func(t *testing.T) {
		rig := newBusRig(t)
		at := time.Date(2026, 7, 1, 9, 5, 0, 0, time.UTC)
		rig.e.fireDeviceError(broker.InboundMessage{
			Realm: realmAlpha, DeviceID: devAlpha, ReceivedAt: at,
		}, "interface_not_declared", "com.nope.Iface")

		ev := rig.next(t)
		if ev.Kind != stream.KindDeviceError {
			t.Fatalf("Kind = %q, want %q", ev.Kind, stream.KindDeviceError)
		}
		if ev.ErrorName != "interface_not_declared" {
			t.Errorf("ErrorName = %q", ev.ErrorName)
		}
		if got := ev.ErrorMetadata["detail"]; got != "com.nope.Iface" {
			t.Errorf("metadata detail = %q", got)
		}
		if !ev.Timestamp.Equal(at) {
			t.Errorf("Timestamp = %v, want %v", ev.Timestamp, at)
		}
	})

	t.Run("DeviceErrorPublishesWithoutAnySchema", func(t *testing.T) {
		// The invariant the change exists for: the trigger path returns early
		// for a realm the schema cache does not know, and the bus publish must
		// escape that early return. The seeded realm declares no triggers
		// either, so both halves of the old guard are covered.
		rig := newBusRig(t)
		if rs := rig.e.schemas.realm(realmAlpha); rs != nil && len(rs.triggers) != 0 {
			t.Fatalf("rig realm unexpectedly declares %d triggers", len(rs.triggers))
		}
		events, cancel := rig.e.bus.Subscribe("ghost", stream.Filter{}, 4)
		defer cancel()

		rig.e.fireDeviceError(broker.InboundMessage{
			Realm: "ghost", DeviceID: devAlpha, ReceivedAt: time.Now().UTC(),
		}, "realm_unknown", "ghost")

		select {
		case ev := <-events:
			if ev.Kind != stream.KindDeviceError || ev.ErrorName != "realm_unknown" {
				t.Errorf("event = %+v", ev)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("device_error for an unknown realm never reached the bus")
		}
	})
}
