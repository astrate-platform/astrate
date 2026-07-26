package broker

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

const (
	// maxTopicBytes bounds accepted topic/filter lengths
	// (docs/DESIGN.md §4.5 input bounds).
	maxTopicBytes = 512

	// offlineACLCacheTTL bounds how long a disconnected device's
	// introspection-derived read permissions are cached for offline-queue
	// delivery checks.
	offlineACLCacheTTL = 10 * time.Second
)

// ownershipFn resolves an interface name from the device's introspection to
// its declared ownership. The second result is false for interfaces the
// device has not introspected (or that are not installed in the realm).
type ownershipFn func(iface string) (interfaceschema.Ownership, bool)

// checkACL is the pure §3.2 ACL matrix decision for a device whose base
// topic is base ("<realm>/<device_id>"):
//
//	PUBLISH:   base | base/control/emptyCache | base/control/producer/properties
//	           | base/<iface><path> for introspected ownership:device interfaces
//	SUBSCRIBE: any filter within the device's own subtree (base/...) — a
//	           tolerated superset; harmless because a device can only ever
//	           match its own topics, and the official SDKs subscribe to
//	           base/<server-iface>/# *before* sending the introspection that
//	           would prove ownership (CP-B). Per-message delivery is gated
//	           independently.
//	DELIVERY:  base/control/consumer/properties | concrete
//	           base/<iface>... for introspected ownership:server interfaces
//
// Everything else is denied. Read checks (write == false) cover both SUBSCRIBE
// filters and per-message delivery topics; the two are told apart by the
// presence of MQTT wildcards, which only ever appear in filters. The
// server-owned delivery rule accepts an interface segment followed by any
// remainder (a concrete path or nothing).
func checkACL(base, topic string, write bool, ownership ownershipFn) bool {
	if topic == "" || len(topic) > maxTopicBytes {
		return false
	}
	if write {
		if topic == base { // introspection publish
			return true
		}
		rest, ok := strings.CutPrefix(topic, base+"/")
		if !ok {
			return false
		}
		if rest == "control/emptyCache" || rest == "control/producer/properties" {
			return true
		}
		iface, path, found := strings.Cut(rest, "/")
		if !found || path == "" { // data topics always carry a path (§3.3)
			return false
		}
		own, known := ownership(iface)
		return known && own == interfaceschema.OwnershipDevice
	}

	rest, ok := strings.CutPrefix(topic, base+"/")
	if !ok {
		return false // never another device's or realm's subtree
	}
	if strings.ContainsAny(topic, "+#") {
		// A SUBSCRIBE filter within the device's own subtree. Wildcards never
		// occur in delivery topics, so this branch only ever sees filters.
		return true
	}
	if rest == "control/consumer/properties" {
		return true
	}
	iface, _, _ := strings.Cut(rest, "/")
	own, known := ownership(iface)
	return known && own == interfaceschema.OwnershipServer
}

// offlineACL answers read-side ACL checks for devices that hold a persistent
// session but are not currently connected: when the engine publishes
// server-owned data to an offline device, mochi still runs the delivery ACL
// check, and there is no live deviceSession to consult. Permissions are
// rebuilt from the store and cached briefly.
type offlineACL struct {
	st    Store
	pools *realmPools
	log   *slog.Logger

	mu      sync.Mutex
	entries map[string]*offlineEntry
}

type offlineEntry struct {
	ownership map[string]interfaceschema.Ownership
	loadedAt  time.Time
}

func newOfflineACL(st Store, pools *realmPools, log *slog.Logger) *offlineACL {
	return &offlineACL{st: st, pools: pools, log: log, entries: map[string]*offlineEntry{}}
}

// ownershipOf resolves iface ownership for the device identified by cn,
// loading (and caching) its introspection from the store when needed.
// Lookup failures cache as an empty map, which denies — the safe posture.
func (o *offlineACL) ownershipOf(cn, iface string) (interfaceschema.Ownership, bool) {
	o.mu.Lock()
	e := o.entries[cn]
	if e != nil && time.Since(e.loadedAt) < offlineACLCacheTTL {
		o.mu.Unlock()
		own, ok := e.ownership[iface]
		return own, ok
	}
	// Claim the slot before the slow load so concurrent deliveries to the
	// same device do not stampede the store.
	e = &offlineEntry{ownership: map[string]interfaceschema.Ownership{}, loadedAt: time.Now()}
	o.entries[cn] = e
	o.mu.Unlock()

	ownership := o.load(cn)

	o.mu.Lock()
	e.ownership = ownership
	o.mu.Unlock()

	own, ok := ownership[iface]
	return own, ok
}

func (o *offlineACL) load(cn string) map[string]interfaceschema.Ownership {
	identity, err := ParseCN(cn)
	if err != nil {
		return map[string]interfaceschema.Ownership{}
	}
	rc, ok := o.pools.Lookup(identity.Realm)
	if !ok {
		return map[string]interfaceschema.Ownership{}
	}
	ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
	defer cancel()
	dev, err := o.st.GetDevice(ctx, rc.id, identity.DeviceID)
	if err != nil {
		o.log.Debug("offline ACL device lookup failed", "client", cn, "error", err)
		return map[string]interfaceschema.Ownership{}
	}
	return loadOwnership(ctx, o.st, rc.id, dev.Introspection, o.log)
}

// aclHook enforces the §3.2 matrix on every SUBSCRIBE filter, device
// PUBLISH, and per-message delivery (docs/ROADMAP.md §6 file 5.4). Inline
// (server-side) clients bypass it by design.
type aclHook struct {
	mqtt.HookBase
	st       Store
	registry *sessionRegistry
	offline  *offlineACL
	log      *slog.Logger
}

// ID implements mqtt.Hook.
func (h *aclHook) ID() string { return "astrate-acl" }

// Provides implements mqtt.Hook.
func (h *aclHook) Provides(b byte) bool { return b == mqtt.OnACLCheck }

// OnACLCheck implements mqtt.Hook.
func (h *aclHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	if cl.Net.Inline {
		return true // engine/AppEngine publishes bypass ACLs (§3.2)
	}

	var allowed bool
	if sess := h.registry.get(cl.ID); sess != nil && sess.client == cl {
		allowed = checkACL(sess.identity.BaseTopic(), topic, write, func(iface string) (interfaceschema.Ownership, bool) {
			own, ok := sess.ownershipOf(iface)
			if !ok {
				// The interface may have been introspected after connect
				// (the engine updates the row); reload, debounced.
				ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
				defer cancel()
				sess.refreshIfStale(ctx, h.st, h.log)
				own, ok = sess.ownershipOf(iface)
			}
			return own, ok
		})
	} else if !write {
		// Delivery to a session-present but disconnected device (offline
		// queue, retained messages): no live session, consult the store.
		allowed = checkACL(cl.ID, topic, write, func(iface string) (interfaceschema.Ownership, bool) {
			return h.offline.ownershipOf(cl.ID, iface)
		})
	}

	if !allowed {
		h.log.Warn("mqtt ACL denied",
			"client", cl.ID, "topic", topic, "write", write, "remote", cl.Net.Remote)
	}
	return allowed
}
