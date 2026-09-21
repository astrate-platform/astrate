slug: broker-offlineacl-entry-eviction
verdict: done
at:  72b2686
ran: 2026-09-21T20:10:44Z on DietPi in 406s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/broker/aclhook.go
✱ Grep "lifecycleHook" 7 matches
✱ Grep "\.now\b|now func\(\) time\.Time" 18 matches
→ Read internal/broker/lifecycle.go
✱ Glob "internal/broker/*_test.go" 5 matches
→ Read internal/broker/aclhook_test.go
✱ Grep "newOfflineACL|offlineACL\{|offlineACLCacheTTL" 6 matches
→ Read internal/broker/broker.go [offset=140, limit=30]
✱ Grep "realmPools|newRealmPools|fakeStore|mockStore|type .*Store" 92 matches
→ Read internal/broker/broker_test.go [limit=160]
✱ Grep "type realmPools|func .*realmPools|func .*Lookup|type realmPool struct" 10 matches
→ Read internal/broker/authhook.go [limit=80, offset=55]
✱ Grep "func ParseCN|func loadOwnership|type Device struct" 4 matches
→ Read internal/broker/identity.go
→ Read internal/broker/broker_test.go [limit=80, offset=205]
→ Read internal/broker/authhook.go [limit=30, offset=255]
✱ Grep "func Random|func Parse|func Decode|func Encode|func New|func FromString" in pkg/deviceid · 2 matches
→ Read pkg/deviceid/deviceid.go
✱ Grep "func discardLogger|discardLogger\(" 31 matches
→ Read internal/broker/authhook.go [limit=54]

# Todos
[•] Add injectable now + lazy stale-entry eviction to offlineACL in aclhook.go
[ ] Add container-free unit test in aclhook_test.go
[ ] Run gate (go test ./..., gofmt -l .)

Now I understand the task and the code. Let me make the change to `internal/broker/aclhook.go`.
← Edit internal/broker/aclhook.go
Index: /root/astrate-mule/internal/broker/aclhook.go
===================================================================
--- /root/astrate-mule/internal/broker/aclhook.go
+++ /root/astrate-mule/internal/broker/aclhook.go
@@ -20,8 +20,15 @@
 	// offlineACLCacheTTL bounds how long a disconnected device's
 	// introspection-derived read permissions are cached for offline-queue
 	// delivery checks.
 	offlineACLCacheTTL = 10 * time.Second
+
+	// offlineACLEvictTTL is how old a cached entry may grow before it is
+	// reaped on the next access. A multiple of offlineACLCacheTTL: an entry
+	// kept alive by a steady stream of offline deliveries is re-stamped on
+	// every miss and survives; one untouched past this is a device that has
+	// gone away, and reaping it keeps the map bounded across device churn.
+	offlineACLEvictTTL = 2 * offlineACLCacheTTL
 )
 
 // ownershipFn resolves an interface name from the device's introspection to
 // its declared ownership. The second result is false for interfaces the

← Edit internal/broker/aclhook.go
Index: /root/astrate-mule/internal/broker/aclhook.go
===================================================================
--- /root/astrate-mule/internal/broker/aclhook.go
+++ /root/astrate-mule/internal/broker/aclhook.go
@@ -107,8 +107,9 @@
 type offlineACL struct {
 	st    Store
 	pools *realmPools
 	log   *slog.Logger
+	now   func() time.Time
 
 	mu      sync.Mutex
 	entries map[string]*offlineEntry
 }

← Edit internal/broker/aclhook.go
Index: /root/astrate-mule/internal/broker/aclhook.go
===================================================================
--- /root/astrate-mule/internal/broker/aclhook.go
+++ /root/astrate-mule/internal/broker/aclhook.go
@@ -119,25 +119,30 @@
 	loadedAt  time.Time
 }
 
 func newOfflineACL(st Store, pools *realmPools, log *slog.Logger) *offlineACL {
-	return &offlineACL{st: st, pools: pools, log: log, entries: map[string]*offlineEntry{}}
+	return &offlineACL{st: st, pools: pools, log: log, now: time.Now, entries: map[string]*offlineEntry{}}
 }
 
 // ownershipOf resolves iface ownership for the device identified by cn,
 // loading (and caching) its introspection from the store when needed.
 // Lookup failures cache as an empty map, which denies — the safe posture.
 func (o *offlineACL) ownershipOf(cn, iface string) (interfaceschema.Ownership, bool) {
+	now := o.now()
 	o.mu.Lock()
 	e := o.entries[cn]
-	if e != nil && time.Since(e.loadedAt) < offlineACLCacheTTL {
+	if e != nil && now.Sub(e.loadedAt) < offlineACLCacheTTL {
 		o.mu.Unlock()
 		own, ok := e.ownership[iface]
 		return own, ok
 	}
+	// Lazy eviction on the inserting access: drop entries untouched for
+	// offlineACLEvictTTL before claiming a slot, so a long-lived broker does
+	// not grow the map unboundedly across device churn.
+	o.evictStale(now)
 	// Claim the slot before the slow load so concurrent deliveries to the
 	// same device do not stampede the store.
-	e = &offlineEntry{ownership: map[string]interfaceschema.Ownership{}, loadedAt: time.Now()}
+	e = &offlineEntry{ownership: map[string]interfaceschema.Ownership{}, loadedAt: now}
 	o.entries[cn] = e
 	o.mu.Unlock()
 
 	ownership := o.load(cn)
@@ -149,8 +154,20 @@
 	own, ok := ownership[iface]
 	return own, ok
 }
 
+// evictStale drops entries whose loadedAt predates offlineACLEvictTTL. The
+// caller holds o.mu. It is only invoked from ownershipOf on the miss path —
+// no background reaper — which keeps the map bounded without a per-access
+// scan on every hit.
+func (o *offlineACL) evictStale(now time.Time) {
+	for cn, e := range o.entries {
+		if e.loadedAt.Before(now.Add(-offlineACLEvictTTL)) {
+			delete(o.entries, cn)
+		}
+	}
+}
+
 func (o *offlineACL) load(cn string) map[string]interfaceschema.Ownership {
 	identity, err := ParseCN(cn)
 	if err != nil {
 		return map[string]interfaceschema.Ownership{}

Now the test. Let me append it to `aclhook_test.go`.
← Edit internal/broker/aclhook_test.go
Index: /root/astrate-mule/internal/broker/aclhook_test.go
===================================================================
--- /root/astrate-mule/internal/broker/aclhook_test.go
+++ /root/astrate-mule/internal/broker/aclhook_test.go
@@ -1,10 +1,13 @@
 package broker
 
 import (
+	"fmt"
 	"strings"
 	"testing"
+	"time"
 
+	"github.com/astrate-platform/astrate/pkg/deviceid"
 	"github.com/astrate-platform/astrate/pkg/interfaceschema"
 )
 
 // TestCheckACLMatrix is the exhaustive §3.2 table test (docs/ROADMAP.md §6):

← Edit internal/broker/aclhook_test.go
Index: /root/astrate-mule/internal/broker/aclhook_test.go
===================================================================
--- /root/astrate-mule/internal/broker/aclhook_test.go
+++ /root/astrate-mule/internal/broker/aclhook_test.go
@@ -146,4 +146,50 @@
 			t.Errorf("checkACL(%q, write=%v) = true, want false", tc.topic, tc.write)
 		}
 	}
 }
+
+// TestOfflineACLStaleEntryEviction covers the lazy eviction of offline ACL
+// cache entries: entries untouched for offlineACLEvictTTL are reaped on the
+// next inserting access, so the map stays bounded across device churn
+// instead of growing one entry per device ever ACL-checked offline.
+func TestOfflineACLStaleEntryEviction(t *testing.T) {
+	st := newFakeStore()
+	pools := &realmPools{byName: map[string]*realmCA{"test": {id: 1, name: "test"}}}
+	acl := newOfflineACL(st, pools, discardLogger())
+
+	cn := func(i int) string {
+		// Deterministic distinct device IDs valid for ParseCN, so a seeded
+		// entry can also be resolved through the store path.
+		id := deviceid.FromNamespace(deviceid.ID{1}, fmt.Sprintf("dev-%d", i))
+		return Identity{Realm: "test", DeviceID: id}.CN()
+	}
+
+	const seeded = 1000
+	base := time.Unix(1_000_000, 0)
+	now := base
+	acl.now = func() time.Time { return now }
+	for i := 0; i < seeded; i++ {
+		acl.entries[cn(i)] = &offlineEntry{
+			ownership: map[string]interfaceschema.Ownership{
+				"com.ex.ServerData": interfaceschema.OwnershipServer,
+			},
+			loadedAt: base,
+		}
+	}
+
+	// Advance the clock well past several TTLs: every seeded entry is both
+	// TTL-stale and past the eviction age.
+	now = now.Add(5 * offlineACLCacheTTL)
+
+	// Accessing a device that was never cached must evict the stale seeds,
+	// leaving at most the freshly-inserted entry rather than 1001.
+	target := cn(seeded)
+	acl.ownershipOf(target, "com.ex.ServerData")
+
+	if got := len(acl.entries); got > 1 {
+		t.Errorf("offline ACL cache has %d entries after lazy eviction, want at most 1 (the accessed entry)", got)
+	}
+	if _, ok := acl.entries[target]; !ok {
+		t.Errorf("accessed entry %q missing from the cache after access", target)
+	}
+}

# Todos
[✓] Add injectable now + lazy stale-entry eviction to offlineACL in aclhook.go
[✓] Add container-free unit test in aclhook_test.go
[•] Run gate (go test ./..., gofmt -l .)

$ gofmt -l internal/broker/aclhook.go internal/broker/aclhook_test.go && go test ./internal/broker/ 2>&1 | tail -30
ok  	github.com/astrate-platform/astrate/internal/broker	3.127s

$ go test ./internal/broker/ -run TestOfflineACLStaleEntryEviction -v 2>&1 | tail -8
=== RUN   TestOfflineACLStaleEntryEviction
--- PASS: TestOfflineACLStaleEntryEviction (0.01s)
PASS
ok  	github.com/astrate-platform/astrate/internal/broker	0.135s

Now the full gate as MULE.md requires.
$ go vet ./... && go test ./... 2>&1 | tail -40
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	1.768s
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

$ git diff --stat && gofmt -l .
internal/broker/aclhook.go      | 31 +++++++++++++++++++++++---
 internal/broker/aclhook_test.go | 49 +++++++++++++++++++++++++++++++++++++++++
 2 files changed, 77 insertions(+), 3 deletions(-)

# Todos
[✓] Add injectable now + lazy stale-entry eviction to offlineACL in aclhook.go
[✓] Add container-free unit test in aclhook_test.go
[✓] Run gate (go test ./..., gofmt -l .)

Done: `offlineACL` now lazily evicts entries untouched for `2*offlineACLCacheTTL` on the insert path, with an injectable `now` clock and a container-free bounded-map test.
Files: internal/broker/aclhook.go, internal/broker/aclhook_test.go
Verified: `go vet ./...` -> pass; `go test ./...` -> pass (no `-race`, we're on the container-free machine where it can't run); `gofmt -l .` -> pass (no output); focused run of the new test -> pass
Unsure: nothing — I chose the "drop stale on insert" variant (evict before claiming the slot, so the fresh entry always survives); map stays bounded by the count of distinct CNs touched within the evict window
Follow-ups: none
