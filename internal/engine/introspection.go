package engine

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/engine/triggers"
	"github.com/astrate-platform/astrate/internal/store"
)

// maxIntrospectionBytes bounds the accepted introspection payload
// (docs/DESIGN.md §4.5 "introspection ≤ 64 KiB").
const maxIntrospectionBytes = 64 << 10

// handleIntrospection processes a device introspection publish — the bare
// "<realm>/<device_id>" topic (docs/DESIGN.md §3.3, docs/ROADMAP.md §7.2
// file 6.7) — on the device's shard goroutine:
//
//  1. parse the `;`-separated `name:major:minor` triples;
//  2. diff against the device's stored introspection;
//  3. persist the new set (the store merges dropped pairs into
//     old_introspection);
//  4. refresh the broker's ACL-relevant introspection cache;
//  5. fire incoming_introspection / interface_added / interface_removed
//     trigger events;
//  6. acknowledge.
//
// Transient store failures park the shard (deferred-ack backpressure,
// docs/DESIGN.md §1.4); malformed payloads are rejected and consumed.
func (e *Engine) handleIntrospection(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
	if len(m.Payload) > maxIntrospectionBytes {
		e.reject(m, reasonIntrospectionInvalid,
			fmt.Sprintf("%d byte introspection exceeds the %d byte cap", len(m.Payload), maxIntrospectionBytes))
		return
	}
	raw := string(m.Payload)
	intro, err := parseIntrospection(raw)
	if err != nil {
		e.reject(m, reasonIntrospectionInvalid, err.Error())
		return
	}

	dev, ok := e.deviceState(ctx, m, realm)
	if !ok {
		return // rejected or abandoned; deviceState handled the message
	}
	prev := dev.introspectionCopy()

	var removed map[string]store.InterfaceVersion
	if !e.retryStore(ctx, m, "introspection update", func() error {
		var err error
		removed, err = e.st.UpdateIntrospection(ctx, realm.id, m.DeviceID, intro)
		return err
	}) {
		return
	}
	dev.setIntrospection(intro)

	// The broker caches iface→ownership per session for its ACL checks; a
	// device declaring a new interface must become publishable immediately
	// (docs/ROADMAP.md §7.2 file 6.7). Failures degrade to the broker's own
	// debounced self-heal reload, so they are logged, not fatal.
	if e.broker != nil {
		if err := e.broker.RefreshIntrospection(ctx, m.Realm, m.DeviceID); err != nil {
			e.log.Warn("broker introspection refresh failed",
				"realm", m.Realm, "device", m.DeviceID.String(), "err", err)
		}
	}

	at := m.ReceivedAt
	e.fireDevice(realm, m.DeviceID, at,
		triggers.DeviceEvent{DeviceID: m.DeviceID.String(), On: triggers.OnIncomingIntrospection},
		triggers.NewIncomingIntrospectionEvent(raw))
	for name, v := range intro {
		if pv, declared := prev[name]; !declared || pv.Major != v.Major {
			e.fireDevice(realm, m.DeviceID, at,
				triggers.DeviceEvent{DeviceID: m.DeviceID.String(), On: triggers.OnInterfaceAdded, Interface: name, Major: v.Major},
				triggers.NewInterfaceAddedEvent(name, v.Major, v.Minor))
		}
	}
	for name, v := range removed {
		e.fireDevice(realm, m.DeviceID, at,
			triggers.DeviceEvent{DeviceID: m.DeviceID.String(), On: triggers.OnInterfaceRemoved, Interface: name, Major: v.Major},
			triggers.NewInterfaceRemovedEvent(name, v.Major))
	}

	m.Ack()
}

// parseIntrospection parses the `;`-separated `name:major:minor` triples
// (docs/DESIGN.md §3.3). The empty payload is a valid empty introspection
// (the device declares no interfaces); any malformed entry rejects the whole
// message — partial introspections would desynchronize the ACL state.
func parseIntrospection(s string) (map[string]store.InterfaceVersion, error) {
	intro := map[string]store.InterfaceVersion{}
	if s == "" {
		return intro, nil
	}
	for _, entry := range strings.Split(s, ";") {
		name, versions, ok := strings.Cut(entry, ":")
		if !ok || name == "" {
			return nil, fmt.Errorf("malformed introspection entry %q", entry)
		}
		if strings.ContainsAny(name, "/+#") {
			return nil, fmt.Errorf("interface name %q contains topic metacharacters", name)
		}
		majorStr, minorStr, ok := strings.Cut(versions, ":")
		if !ok {
			return nil, fmt.Errorf("malformed introspection entry %q", entry)
		}
		major, err := strconv.Atoi(majorStr)
		if err != nil || major < 0 {
			return nil, fmt.Errorf("invalid major version in introspection entry %q", entry)
		}
		minor, err := strconv.Atoi(minorStr)
		if err != nil || minor < 0 {
			return nil, fmt.Errorf("invalid minor version in introspection entry %q", entry)
		}
		intro[name] = store.InterfaceVersion{Major: major, Minor: minor}
	}
	return intro, nil
}

// retryStore runs fn with shard-parking retries for transient store
// failures, mirroring the data path's DB-outage behaviour (docs/DESIGN.md
// §5.3): while parked the shard does not consume its channel, so the §1.4
// backpressure builds. It reports whether fn eventually succeeded; on every
// false return the message has been fully handled (rejected and consumed, or
// abandoned unacknowledged at shutdown).
func (e *Engine) retryStore(ctx context.Context, m broker.InboundMessage, what string, fn func() error) bool {
	backoff := parkBackoffStart
	for {
		err := fn()
		if err == nil {
			return true
		}
		if errors.Is(err, store.ErrNotFound) {
			e.reject(m, reasonDeviceUnknown, err.Error())
			return false
		}
		if isPermanentCommitError(err) {
			// The database permanently refuses the operation (integrity- or
			// data-class rejection): consume the message so the device is
			// not wedged forever, and be loud about it.
			e.met.internalErrors.Inc()
			e.log.Error(what+" permanently rejected by the database; dropping message",
				"realm", m.Realm, "device", m.DeviceID.String(), "topic", m.Topic, "err", err)
			m.Ack()
			return false
		}
		e.met.flushRetries.Inc()
		e.log.Warn(what+" failed; parking shard",
			"realm", m.Realm, "device", m.DeviceID.String(), "backoff", backoff, "err", err)
		select {
		case <-e.quit:
			return false
		case <-ctx.Done():
			return false
		case <-time.After(backoff):
		}
		backoff = min(backoff*2, parkBackoffCap)
	}
}
