package broker

import (
	"context"
	"log/slog"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

// lifecycleHook turns broker connection events into devices-row updates and
// LifecycleSink events (docs/ROADMAP.md §6 file 5.6) — the work
// astarte_vmq_plugin does in upstream Astarte (docs/DESIGN.md §3.1).
//
// It acts only on sessions the auth hook admitted in this process: the
// synthetic disconnects mochi emits while restoring persisted sessions at
// boot, and the old connection's teardown after a session takeover, touch
// neither the database nor the sink.
type lifecycleHook struct {
	mqtt.HookBase
	st       Store
	sink     LifecycleSink
	registry *sessionRegistry
	log      *slog.Logger
	now      func() time.Time
}

// ID implements mqtt.Hook.
func (h *lifecycleHook) ID() string { return "astrate-lifecycle" }

// Provides implements mqtt.Hook.
func (h *lifecycleHook) Provides(b byte) bool {
	return b == mqtt.OnSessionEstablished || b == mqtt.OnDisconnect
}

// OnSessionEstablished records the connection on the device row and emits
// device_connected.
func (h *lifecycleHook) OnSessionEstablished(cl *mqtt.Client, _ packets.Packet) {
	sess := h.registry.get(cl.ID)
	if sess == nil || sess.client != cl {
		return
	}
	at := h.now()

	ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
	defer cancel()
	if err := h.st.SetDeviceConnected(ctx, sess.realmID, sess.identity.DeviceID, at, sess.remote); err != nil {
		h.log.Error("recording device connection", "client", cl.ID, "error", err)
	}

	if h.sink != nil {
		h.sink.OnLifecycleEvent(LifecycleEvent{
			Type:     EventDeviceConnected,
			Realm:    sess.identity.Realm,
			DeviceID: sess.identity.DeviceID,
			RemoteIP: sess.remote,
			At:       at,
		})
	}
}

// OnDisconnect records the disconnection and emits device_disconnected.
// Session takeovers are skipped: the device is still connected, on the new
// channel, and the registry entry already belongs to it.
func (h *lifecycleHook) OnDisconnect(cl *mqtt.Client, _ error, _ bool) {
	if cl.IsTakenOver() {
		return
	}
	sess := h.registry.removeIfOwner(cl.ID, cl)
	if sess == nil {
		return
	}
	at := h.now()

	ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
	defer cancel()
	if err := h.st.SetDeviceDisconnected(ctx, sess.realmID, sess.identity.DeviceID, at); err != nil {
		h.log.Error("recording device disconnection", "client", cl.ID, "error", err)
	}

	if h.sink != nil {
		h.sink.OnLifecycleEvent(LifecycleEvent{
			Type:     EventDeviceDisconnected,
			Realm:    sess.identity.Realm,
			DeviceID: sess.identity.DeviceID,
			At:       at,
		})
	}
}
