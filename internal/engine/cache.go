// Package engine is the heart of Astrate (docs/ROADMAP.md §7): it consumes
// every accepted device PUBLISH from the broker intake, validates it against
// the realm's compiled interface schemas (docs/DESIGN.md §2.6), and persists
// it through per-shard micro-batches with strict per-device ordering and
// ack-after-commit semantics (docs/DESIGN.md §1.4, §5.3).
//
// M6a (this milestone slice) implements the pipeline and the data path:
// compiled-interface and device caches, the sharded router, topic
// classification, the §2.6 validation pipeline, and batched persistence.
// M6b adds the control-channel handlers, server-owned data publishing,
// triggers, and the live stream — they attach to the seams the Engine struct
// already exposes (handler fields and the afterCommit hook).
package engine

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// Store is the engine's persistence port (hexagonal-lite, docs/DESIGN.md
// §1.3): the subset of *store.Store the M6a pipeline needs. The interface is
// defined here, on the consumer side, so tests substitute fakes and M8 wiring
// is a plain assignment.
type Store interface {
	// ListRealms feeds the realm name/ID resolution in the schema cache.
	ListRealms(ctx context.Context) ([]store.Realm, error)
	// LoadRealmInterfaces returns every installed interface of a realm for
	// schema-snapshot rebuilds.
	LoadRealmInterfaces(ctx context.Context, realmID int16) ([]*store.StoredInterface, error)
	// Listen subscribes to a NOTIFY channel (cache invalidation,
	// docs/DESIGN.md §2.6).
	Listen(ctx context.Context, channel string) (<-chan store.Notification, error)
	// GetDevice loads a device's introspection and payload-format hint.
	GetDevice(ctx context.Context, realmID int16, id deviceid.ID) (*store.Device, error)
	// SetPayloadFormatHint persists the sticky payload-format flip
	// (docs/DESIGN.md §3.5.4).
	SetPayloadFormatHint(ctx context.Context, realmID int16, id deviceid.ID, hint string) error
	// AppendDatastreams commits one micro-batch in a single transaction.
	AppendDatastreams(ctx context.Context, batch store.DatastreamBatch) error
	// UpsertProperty applies a property set (last-value-wins).
	UpsertProperty(ctx context.Context, p store.Property) error
	// UnsetProperty applies a property unset (row delete).
	UnsetProperty(ctx context.Context, realmID int16, deviceID deviceid.ID, interfaceID int64, path string) (bool, error)
}

// fullReloadDebounce rate-limits snapshot self-heal reloads triggered by
// realm-lookup misses, so a flood of messages for an unknown realm cannot
// hammer the database (same defense as the broker's introspection reload).
const fullReloadDebounce = time.Second

// realmSchema is one realm's compiled-interface set inside a snapshot. It is
// immutable after construction: rebuilds replace the whole *realmSchema.
type realmSchema struct {
	id   int16
	name string
	// interfaces is name → major → compiled interface (docs/ROADMAP.md §7.1
	// file 6.1).
	interfaces map[string]map[int]*interfaceschema.CompiledInterface
}

// iface resolves one name:major pair, or nil.
func (r *realmSchema) iface(name string, major int) *interfaceschema.CompiledInterface {
	if r == nil {
		return nil
	}
	return r.interfaces[name][major]
}

// schemaSnapshot is the immutable compiled-interface view shared by every
// shard. Readers never lock: they load the current snapshot pointer and use
// it for the whole message (docs/DESIGN.md §2.6 "Cache & invalidation").
type schemaSnapshot struct {
	byID   map[int16]*realmSchema
	byName map[string]*realmSchema
}

// schemaCache owns the snapshot pointer and its copy-on-write rebuilds.
type schemaCache struct {
	st  Store
	log *slog.Logger

	snap atomic.Pointer[schemaSnapshot]
	// reloadMu serializes rebuilds so concurrent invalidations cannot
	// interleave their read-modify-swap sequences.
	reloadMu sync.Mutex
	// lastFullReload is the wall-clock nanos of the last self-heal full
	// reload (fullReloadDebounce).
	lastFullReload atomic.Int64
}

// newSchemaCache returns a cache holding an empty snapshot; call loadAll
// before serving traffic.
func newSchemaCache(st Store, log *slog.Logger) *schemaCache {
	c := &schemaCache{st: st, log: log}
	c.snap.Store(&schemaSnapshot{
		byID:   map[int16]*realmSchema{},
		byName: map[string]*realmSchema{},
	})
	return c
}

// realm resolves a realm by name from the current snapshot.
func (c *schemaCache) realm(name string) *realmSchema {
	return c.snap.Load().byName[name]
}

// realmOrReload resolves a realm by name, attempting one debounced full
// reload on a miss (self-heal: a realm created after startup becomes visible
// without a restart even if no interface notification ever fires).
func (c *schemaCache) realmOrReload(ctx context.Context, name string) *realmSchema {
	if r := c.realm(name); r != nil {
		return r
	}
	last := c.lastFullReload.Load()
	if time.Now().UnixNano()-last < int64(fullReloadDebounce) {
		return nil
	}
	if !c.lastFullReload.CompareAndSwap(last, time.Now().UnixNano()) {
		return nil // another shard is already reloading
	}
	if err := c.loadAll(ctx); err != nil {
		c.log.Error("schema cache self-heal reload failed", "realm", name, "err", err)
		return nil
	}
	return c.realm(name)
}

// loadAll rebuilds the whole snapshot from the store (startup, unknown-realm
// invalidation payloads, self-heal).
func (c *schemaCache) loadAll(ctx context.Context) error {
	c.reloadMu.Lock()
	defer c.reloadMu.Unlock()

	realms, err := c.st.ListRealms(ctx)
	if err != nil {
		return fmt.Errorf("engine: listing realms: %w", err)
	}
	next := &schemaSnapshot{
		byID:   make(map[int16]*realmSchema, len(realms)),
		byName: make(map[string]*realmSchema, len(realms)),
	}
	for i := range realms {
		rs, err := c.buildRealm(ctx, realms[i].ID, realms[i].Name)
		if err != nil {
			return err
		}
		next.byID[rs.id] = rs
		next.byName[rs.name] = rs
	}
	c.snap.Store(next)
	return nil
}

// reloadRealm rebuilds a single realm's entry copy-on-write: the maps are
// cloned, one entry is replaced, and the new snapshot is swapped in. A realm
// ID missing from the snapshot falls back to a full reload (it may be a
// newly created realm).
func (c *schemaCache) reloadRealm(ctx context.Context, realmID int16) error {
	c.reloadMu.Lock()
	cur := c.snap.Load()
	old, known := cur.byID[realmID]
	if !known {
		c.reloadMu.Unlock()
		return c.loadAll(ctx)
	}
	rs, err := c.buildRealm(ctx, realmID, old.name)
	if err != nil {
		c.reloadMu.Unlock()
		return err
	}
	next := &schemaSnapshot{
		byID:   make(map[int16]*realmSchema, len(cur.byID)),
		byName: make(map[string]*realmSchema, len(cur.byName)),
	}
	for id, r := range cur.byID {
		next.byID[id] = r
	}
	for name, r := range cur.byName {
		next.byName[name] = r
	}
	next.byID[realmID] = rs
	next.byName[rs.name] = rs
	c.snap.Store(next)
	c.reloadMu.Unlock()
	return nil
}

// buildRealm loads and compiles every interface of one realm.
func (c *schemaCache) buildRealm(ctx context.Context, realmID int16, name string) (*realmSchema, error) {
	stored, err := c.st.LoadRealmInterfaces(ctx, realmID)
	if err != nil {
		return nil, fmt.Errorf("engine: loading interfaces of realm %s: %w", name, err)
	}
	rs := &realmSchema{
		id:         realmID,
		name:       name,
		interfaces: make(map[string]map[int]*interfaceschema.CompiledInterface, len(stored)),
	}
	for _, si := range stored {
		iface, err := interfaceschema.ParseInterface(si.Definition)
		if err != nil {
			// A stored definition that no longer parses is a bug, not a
			// reason to take the whole realm down: skip it loudly.
			c.log.Error("stored interface definition does not parse",
				"realm", name, "interface", si.Name, "major", si.Major, "err", err)
			continue
		}
		ci, err := interfaceschema.Compile(iface, si)
		if err != nil {
			c.log.Error("stored interface does not compile",
				"realm", name, "interface", si.Name, "major", si.Major, "err", err)
			continue
		}
		byMajor := rs.interfaces[ci.Name]
		if byMajor == nil {
			byMajor = make(map[int]*interfaceschema.CompiledInterface, 1)
			rs.interfaces[ci.Name] = byMajor
		}
		byMajor[ci.Major] = ci
	}
	return rs, nil
}

// formatHint values mirror devices.payload_format_hint (docs/DESIGN.md
// §3.5.4; CHECK-constrained to bson|json in the schema).
const (
	hintBSON = "bson"
	hintJSON = "json"
)

// deviceKey addresses one device across realms.
type deviceKey struct {
	realm string
	id    deviceid.ID
}

// deviceState is the per-connected-device cache entry (docs/ROADMAP.md §7.1
// file 6.1): the introspection gate input plus the payload-format hint
// state. All fields after construction are guarded by mu — the owning shard
// is the only writer, but eviction and (in M6b) the AppEngine publish path
// read concurrently.
type deviceState struct {
	realmID int16

	mu sync.Mutex
	// introspection is the device's declared interface set (name →
	// major/minor), the §2.6 step 1 gate.
	introspection map[string]store.InterfaceVersion
	// formatHint is the current outbound payload format ("bson"|"json").
	formatHint string
	// resetHintOnBSON is armed by the emptyCache control handler (M6b): the
	// next BSON data payload flips the hint back to bson (docs/DESIGN.md
	// §3.5.4).
	resetHintOnBSON bool
}

// declares reports whether the device's introspection contains name, and at
// which version.
func (d *deviceState) declares(name string) (store.InterfaceVersion, bool) {
	d.mu.Lock()
	v, ok := d.introspection[name]
	d.mu.Unlock()
	return v, ok
}

// setIntrospection replaces the cached introspection (M6b's introspection
// handler updates it after persisting).
func (d *deviceState) setIntrospection(intro map[string]store.InterfaceVersion) {
	d.mu.Lock()
	d.introspection = intro
	d.mu.Unlock()
}

// hint returns the current payload-format hint.
func (d *deviceState) hint() string {
	d.mu.Lock()
	h := d.formatHint
	d.mu.Unlock()
	return h
}

// deviceCache holds per-connected-device state, loaded lazily on a device's
// first message and evicted on disconnect (docs/ROADMAP.md §7.1 file 6.1),
// so memory scales with connected devices, not registered ones.
type deviceCache struct {
	st Store

	mu sync.Mutex
	m  map[deviceKey]*deviceState
}

// newDeviceCache returns an empty device cache.
func newDeviceCache(st Store) *deviceCache {
	return &deviceCache{st: st, m: make(map[deviceKey]*deviceState)}
}

// get returns the cached state for the device, loading it from the store on
// a miss. Every message of a device is processed by the same shard, so loads
// for one device never race each other.
func (c *deviceCache) get(ctx context.Context, realm string, realmID int16, id deviceid.ID) (*deviceState, error) {
	key := deviceKey{realm: realm, id: id}
	c.mu.Lock()
	st, ok := c.m[key]
	c.mu.Unlock()
	if ok {
		return st, nil
	}

	dev, err := c.st.GetDevice(ctx, realmID, id)
	if err != nil {
		return nil, err
	}
	hint := dev.PayloadFormatHint
	if hint == "" {
		hint = hintBSON
	}
	st = &deviceState{
		realmID:       realmID,
		introspection: dev.Introspection,
		formatHint:    hint,
	}
	c.mu.Lock()
	c.m[key] = st
	c.mu.Unlock()
	return st, nil
}

// evict drops the device's entry; the next message reloads it from the
// store. Wired to device_disconnected lifecycle events.
func (c *deviceCache) evict(realm string, id deviceid.ID) {
	c.mu.Lock()
	delete(c.m, deviceKey{realm: realm, id: id})
	c.mu.Unlock()
}

// len reports the number of cached devices (metrics, tests).
func (c *deviceCache) len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.m)
}
