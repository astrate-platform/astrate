package engine

import (
	"context"
	"sync"
	"testing"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// TestSchemaCacheLoadAndLookup covers the snapshot build and the hot-path
// lookups (docs/ROADMAP.md §7.1 file 6.1).
func TestSchemaCacheLoadAndLookup(t *testing.T) {
	ctx := context.Background()
	fs := newFakeStore()
	fs.addRealm(1, "alpha")
	fs.addRealm(2, "beta")
	fs.addInterface(1, fixtureStored(t, 1, 10, "com.astrate.test.AllScalarTypes.json"))
	fs.addInterface(1, fixtureStored(t, 1, 11, "com.astrate.test.ObjectFlat.json"))
	fs.addInterface(2, fixtureStored(t, 2, 20, "com.astrate.test.Minimal.json"))

	c := newSchemaCache(fs, discardLogger())
	if err := c.loadAll(ctx); err != nil {
		t.Fatalf("loadAll: %v", err)
	}

	alpha := c.realm("alpha")
	if alpha == nil || alpha.id != 1 {
		t.Fatalf("realm alpha not resolved: %+v", alpha)
	}
	ci := alpha.iface("com.astrate.test.AllScalarTypes", 1)
	if ci == nil || ci.ID != 10 {
		t.Fatalf("AllScalarTypes v1 not compiled: %+v", ci)
	}
	if m, ok := ci.Trie.Match("/double"); !ok || m.EndpointID != 10*1000 {
		t.Errorf("trie match /double: ok=%v mapping=%+v", ok, m)
	}
	if alpha.iface("com.astrate.test.AllScalarTypes", 2) != nil {
		t.Error("wrong major resolved")
	}
	if alpha.iface("com.astrate.test.Minimal", 0) != nil {
		t.Error("beta's interface leaked into alpha")
	}
	if c.realm("gamma") != nil {
		t.Error("unknown realm resolved")
	}
	if obj := alpha.iface("com.astrate.test.ObjectFlat", 1); obj == nil || obj.ObjectLeaves["latitude"] == nil {
		t.Errorf("object leaves not compiled: %+v", obj)
	}
}

// TestSchemaCacheReloadRealm covers the copy-on-write single-realm rebuild:
// the reloaded realm changes, untouched realms keep their pointers.
func TestSchemaCacheReloadRealm(t *testing.T) {
	ctx := context.Background()
	fs := newFakeStore()
	fs.addRealm(1, "alpha")
	fs.addRealm(2, "beta")
	fs.addInterface(2, fixtureStored(t, 2, 20, "com.astrate.test.Minimal.json"))

	c := newSchemaCache(fs, discardLogger())
	if err := c.loadAll(ctx); err != nil {
		t.Fatalf("loadAll: %v", err)
	}
	betaBefore := c.realm("beta")

	if c.realm("alpha").iface("com.astrate.test.AllScalarTypes", 1) != nil {
		t.Fatal("interface resolvable before install")
	}
	fs.addInterface(1, fixtureStored(t, 1, 10, "com.astrate.test.AllScalarTypes.json"))
	if err := c.reloadRealm(ctx, 1); err != nil {
		t.Fatalf("reloadRealm: %v", err)
	}
	if c.realm("alpha").iface("com.astrate.test.AllScalarTypes", 1) == nil {
		t.Error("interface not resolvable after reloadRealm")
	}
	if c.realm("beta") != betaBefore {
		t.Error("untouched realm was rebuilt (copy-on-write violated)")
	}
}

// TestSchemaCacheReloadUnknownRealm covers the full-reload fallback for
// invalidations naming a realm the snapshot has never seen.
func TestSchemaCacheReloadUnknownRealm(t *testing.T) {
	ctx := context.Background()
	fs := newFakeStore()
	fs.addRealm(1, "alpha")

	c := newSchemaCache(fs, discardLogger())
	if err := c.loadAll(ctx); err != nil {
		t.Fatalf("loadAll: %v", err)
	}

	fs.addRealm(7, "gamma")
	fs.addInterface(7, fixtureStored(t, 7, 70, "com.astrate.test.Minimal.json"))
	if err := c.reloadRealm(ctx, 7); err != nil {
		t.Fatalf("reloadRealm(unknown): %v", err)
	}
	if c.realm("gamma").iface("com.astrate.test.Minimal", 0) == nil {
		t.Error("new realm not visible after fallback full reload")
	}
}

// TestSchemaCacheSelfHeal covers realmOrReload: a miss triggers one debounced
// full reload; rapid follow-up misses do not hammer the store.
func TestSchemaCacheSelfHeal(t *testing.T) {
	ctx := context.Background()
	fs := newFakeStore()
	c := newSchemaCache(fs, discardLogger())
	if err := c.loadAll(ctx); err != nil {
		t.Fatalf("loadAll: %v", err)
	}
	fs.mu.Lock()
	calls := fs.listRealmsCalls
	fs.mu.Unlock()

	fs.addRealm(1, "alpha")
	if r := c.realmOrReload(ctx, "alpha"); r == nil || r.id != 1 {
		t.Fatalf("self-heal reload did not resolve the realm: %+v", r)
	}
	if c.realmOrReload(ctx, "missing") != nil {
		t.Error("missing realm resolved")
	}
	fs.mu.Lock()
	delta := fs.listRealmsCalls - calls
	fs.mu.Unlock()
	if delta != 1 {
		t.Errorf("ListRealms called %d times within the debounce window, want 1", delta)
	}
}

// TestSchemaCacheBrokenDefinitions: a stored definition that fails to parse
// or compile is skipped loudly without sinking the realm.
func TestSchemaCacheBrokenDefinitions(t *testing.T) {
	ctx := context.Background()
	fs := newFakeStore()
	fs.addRealm(1, "alpha")
	fs.addInterface(1, &store.StoredInterface{
		ID: 9, RealmID: 1, Name: "com.broken.Parse", Major: 1,
		Definition: []byte(`{"not": "an interface"}`),
	})
	// Parses, but the resolver has no endpoint IDs: Compile fails.
	noEndpoints := fixtureStored(t, 1, 8, "com.astrate.test.Minimal.json")
	noEndpoints.Endpoints = map[string]int64{}
	fs.addInterface(1, noEndpoints)
	fs.addInterface(1, fixtureStored(t, 1, 10, "com.astrate.test.AllScalarTypes.json"))

	c := newSchemaCache(fs, discardLogger())
	if err := c.loadAll(ctx); err != nil {
		t.Fatalf("loadAll: %v", err)
	}
	alpha := c.realm("alpha")
	if alpha.iface("com.broken.Parse", 1) != nil || alpha.iface("com.astrate.test.Minimal", 0) != nil {
		t.Error("broken definitions were compiled")
	}
	if alpha.iface("com.astrate.test.AllScalarTypes", 1) == nil {
		t.Error("healthy interface was lost alongside broken ones")
	}
}

// TestSchemaCacheSnapshotRace exercises concurrent snapshot reads against
// copy-on-write rebuilds; the -race gate (docs/ROADMAP.md §7.3) is the
// assertion.
func TestSchemaCacheSnapshotRace(t *testing.T) {
	ctx := context.Background()
	fs := newFakeStore()
	fs.addRealm(1, "alpha")
	fs.addInterface(1, fixtureStored(t, 1, 10, "com.astrate.test.AllScalarTypes.json"))

	c := newSchemaCache(fs, discardLogger())
	if err := c.loadAll(ctx); err != nil {
		t.Fatalf("loadAll: %v", err)
	}

	stop := make(chan struct{})
	var wg sync.WaitGroup
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				if r := c.realm("alpha"); r != nil {
					if ci := r.iface("com.astrate.test.AllScalarTypes", 1); ci != nil {
						ci.Trie.Match("/double")
					}
				}
			}
		}()
	}
	for i := range 200 {
		if i%10 == 0 {
			if err := c.loadAll(ctx); err != nil {
				t.Errorf("loadAll under load: %v", err)
			}
			continue
		}
		if err := c.reloadRealm(ctx, 1); err != nil {
			t.Errorf("reloadRealm under load: %v", err)
		}
	}
	close(stop)
	wg.Wait()
}

// TestDeviceCache covers lazy load, default hint, single round-trip reuse,
// and disconnect eviction.
func TestDeviceCache(t *testing.T) {
	ctx := context.Background()
	id := deviceid.ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	fs := newFakeStore()
	fs.addDevice("alpha", 1, id, map[string]store.InterfaceVersion{
		"com.astrate.test.Minimal": {Major: 0, Minor: 1},
	}, "")

	c := newDeviceCache(fs)
	st, err := c.get(ctx, "alpha", 1, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if v, ok := st.declares("com.astrate.test.Minimal"); !ok || v.Major != 0 {
		t.Errorf("introspection not loaded: %+v ok=%v", v, ok)
	}
	if _, ok := st.declares("com.other.Iface"); ok {
		t.Error("undeclared interface reported as declared")
	}
	if st.hint() != hintBSON {
		t.Errorf("empty stored hint defaulted to %q, want bson", st.hint())
	}

	again, err := c.get(ctx, "alpha", 1, id)
	if err != nil {
		t.Fatalf("get (cached): %v", err)
	}
	if again != st {
		t.Error("second get did not return the cached entry")
	}
	fs.mu.Lock()
	calls := fs.getDeviceCalls
	fs.mu.Unlock()
	if calls != 1 {
		t.Errorf("GetDevice called %d times, want 1", calls)
	}

	c.evict("alpha", id)
	if c.len() != 0 {
		t.Fatalf("cache holds %d entries after evict", c.len())
	}
	if _, err := c.get(ctx, "alpha", 1, id); err != nil {
		t.Fatalf("get after evict: %v", err)
	}
	fs.mu.Lock()
	calls = fs.getDeviceCalls
	fs.mu.Unlock()
	if calls != 2 {
		t.Errorf("GetDevice called %d times after evict, want 2", calls)
	}
}

// TestDeviceCacheLoadError: store failures surface to the caller and nothing
// is cached.
func TestDeviceCacheLoadError(t *testing.T) {
	ctx := context.Background()
	id := deviceid.ID{1}
	fs := newFakeStore()
	fs.addDevice("alpha", 1, id, nil, "json")
	fs.mu.Lock()
	fs.getDeviceErrs = []error{errTransient}
	fs.mu.Unlock()

	c := newDeviceCache(fs)
	if _, err := c.get(ctx, "alpha", 1, id); err == nil {
		t.Fatal("transient store error swallowed")
	}
	if c.len() != 0 {
		t.Error("failed load left a cache entry")
	}
	st, err := c.get(ctx, "alpha", 1, id)
	if err != nil {
		t.Fatalf("get after transient error: %v", err)
	}
	if st.hint() != hintJSON {
		t.Errorf("stored hint %q not honoured", st.hint())
	}
}
