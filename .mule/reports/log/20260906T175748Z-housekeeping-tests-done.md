slug: housekeeping-tests
verdict: done
at:  7583a7d
ran: 2026-09-06T17:57:48Z on DietPi in 455s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
Task: Add container-free unit tests for `internal/housekeeping` — cover `CreateRealm` validation and `DeleteRealm` gating with httptest-style fakes.
→ Read internal/housekeeping
→ Read internal/housekeeping/service.go
→ Read internal/housekeeping/http_test.go
→ Read internal/store
→ Read internal/store/store.go
→ Read internal/store/crypto.go
✱ Grep "DeviceStats" 11 matches
✱ Grep "CreateRealm\(" 67 matches
→ Read internal/store/realms.go
→ Read internal/store/devices.go [offset=440, limit=40]
→ Read internal/housekeeping/http.go
✱ Grep "NewService\(" in internal/housekeeping · 2 matches
✱ Grep "Service\{|housekeeping\.NewService|NewService\(st" 31 matches
Good — `flowapi` already uses function/interface ports (`IngestFunc`), and `realm` uses an `Invalidator` interface. Now let me check `ca.Generate` since a default-retention test hits it.
✱ Grep "func Generate" in internal/pairing · 1 match
→ Read internal/pairing/ca/ca.go [limit=60, offset=55]
Now I have the full picture. `Service` holds concrete `*store.Store` and `*store.KeySealer`, neither fakeable without a DB. The codebase already ports the store behind interfaces (`realm.Invalidator`, `appengine.ServerData`, `flowapi.IngestFunc`), so I'll do the same here: replace the two concrete fields with narrow ports, keep `NewService` callers source-compatible, then add the fakes and tests.

← Edit internal/housekeeping/service.go
Index: /root/astrate-mule/internal/housekeeping/service.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service.go
+++ /root/astrate-mule/internal/housekeeping/service.go
@@ -47,8 +47,24 @@
 type Reloader interface {
 	ReloadRealms(ctx context.Context) error
 }
 
+// realmStore is the slice of the store's realm repository the service needs.
+// *store.Store satisfies it; tests fake it to keep the service's unit tests
+// container-free.
+type realmStore interface {
+	CreateRealm(ctx context.Context, nr store.NewRealm) (*store.Realm, error)
+	GetRealmByName(ctx context.Context, name string) (*store.Realm, error)
+	DeviceStats(ctx context.Context, realmID int16) (total, connected int64, err error)
+	DeleteRealm(ctx context.Context, name string) error
+}
+
+// sealer protects freshly-minted CA private keys before they reach the store.
+// *store.KeySealer satisfies it; tests fake it.
+type sealer interface {
+	Seal(plaintext []byte) ([]byte, error)
+}
+
 // Service implements the Housekeeping business logic over the store, holding
 // the key sealer used to protect freshly-minted CA private keys.
 type Service struct {
 	st       *store.Store

← Edit internal/housekeeping/service.go
Index: /root/astrate-mule/internal/housekeeping/service.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service.go
+++ /root/astrate-mule/internal/housekeeping/service.go
@@ -66,10 +66,10 @@
 
 // Service implements the Housekeeping business logic over the store, holding
 // the key sealer used to protect freshly-minted CA private keys.
 type Service struct {
-	st       *store.Store
-	sealer   *store.KeySealer
+	st       realmStore
+	sealer   sealer
 	reloader Reloader
 	log      *slog.Logger
 
 	// defaultRetention is the realm-level datastream retention ceiling

← Edit internal/housekeeping/service.go
Index: /root/astrate-mule/internal/housekeeping/service.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service.go
+++ /root/astrate-mule/internal/housekeeping/service.go
@@ -85,9 +85,9 @@
 }
 
 // NewService builds the service. reloader (the broker) may be nil; log
 // defaults to slog.Default().
-func NewService(st *store.Store, sealer *store.KeySealer, reloader Reloader, log *slog.Logger) *Service {
+func NewService(st realmStore, sealer sealer, reloader Reloader, log *slog.Logger) *Service {
 	if log == nil {
 		log = slog.Default()
 	}
 	return &Service{st: st, sealer: sealer, reloader: reloader, log: log}

Now the test file:
← Write internal/housekeeping/service_test.go
Wrote file successfully.

The `storePtr` helper is awkward — let me replace it with plain comparisons.
← Edit internal/housekeeping/service_test.go
Index: /root/astrate-mule/internal/housekeeping/service_test.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service_test.go
+++ /root/astrate-mule/internal/housekeeping/service_test.go
@@ -133,13 +133,13 @@
 	}
 	if len(st.created) != 1 {
 		t.Fatalf("store.CreateRealm calls = %d, want 1", len(st.created))
 	}
-	if got := storePtr(st.created[0].DatastreamMaximumStorageRetention, def); !*got {
-		t.Errorf("store received retention = %v, want %d", st.created[0].DatastreamMaximumStorageRetention, def)
+	if got := st.created[0].DatastreamMaximumStorageRetention; got == nil || *got != def {
+		t.Errorf("store received retention = %v, want %d", got, def)
 	}
-	if got := storePtr(rv.DatastreamMaximumStorageRetention, def); !*got {
-		t.Errorf("view retention = %v, want %d", rv.DatastreamMaximumStorageRetention, def)
+	if got := rv.DatastreamMaximumStorageRetention; got == nil || *got != def {
+		t.Errorf("view retention = %v, want %d", got, def)
 	}
 })
 
 t.Run("explicit value beats default", func(t *testing.T) {

← Edit internal/housekeeping/service_test.go
Index: /root/astrate-mule/internal/housekeeping/service_test.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service_test.go
+++ /root/astrate-mule/internal/housekeeping/service_test.go
@@ -150,10 +150,10 @@
 
 	if _, err := svc.CreateRealm(ctx, "micro", "key", nil, &explicit); err != nil {
 		t.Fatalf("CreateRealm: %v", err)
 	}
-	if got := storePtr(st.created[0].DatastreamMaximumStorageRetention, explicit); !*got {
-		t.Errorf("store received retention = %v, want %d", st.created[0].DatastreamMaximumStorageRetention, explicit)
+	if got := st.created[0].DatastreamMaximumStorageRetention; got == nil || *got != explicit {
+		t.Errorf("store received retention = %v, want %d", got, explicit)
 	}
 })
 
 t.Run("no default stays nil", func(t *testing.T) {

← Edit internal/housekeeping/service_test.go
Index: /root/astrate-mule/internal/housekeeping/service_test.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service_test.go
+++ /root/astrate-mule/internal/housekeeping/service_test.go
@@ -211,20 +211,5 @@
 		if len(st.deleted) != 1 || st.deleted[0] != "calm" {
 			t.Errorf("store.DeleteRealm calls = %v, want [calm]", st.deleted)
 		}
 	})
-}
-
-// storePtr reports whether p holds v; it exists to keep the pointer-compare
-// t.Error calls above readable.
-func storePtr(p *int64, v int64) *bool {
-	if p == nil {
-		f := false
-		return &f
-	}
-	if *p == v {
-		t := true
-		return &t
-	}
-	t := false
-	return &t
 }
\ No newline at end of file

Now format and run the gate:
$ gofmt -l internal/housekeeping/ && go vet ./internal/housekeeping/ && go test ./internal/housekeeping/ -run 'TestCreateRealm|TestDeleteRealm' -v 2>&1 | tail -40
internal/housekeeping/service_test.go
# github.com/astrate-platform/astrate/internal/housekeeping
# [github.com/astrate-platform/astrate/internal/housekeeping]
vet: internal/housekeeping/service.go:204:14: s.st.UpdateRealm undefined (type realmStore has no field or method UpdateRealm)

`UpdateRealm` also goes through the store — add it to the port.
← Edit internal/housekeeping/service.go
Index: /root/astrate-mule/internal/housekeeping/service.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service.go
+++ /root/astrate-mule/internal/housekeeping/service.go
@@ -53,8 +53,9 @@
 // container-free.
 type realmStore interface {
 	CreateRealm(ctx context.Context, nr store.NewRealm) (*store.Realm, error)
 	GetRealmByName(ctx context.Context, name string) (*store.Realm, error)
+	UpdateRealm(ctx context.Context, name string, p store.RealmPatch) error
 	DeviceStats(ctx context.Context, realmID int16) (total, connected int64, err error)
 	DeleteRealm(ctx context.Context, name string) error
 }
 

$ gofmt -w internal/housekeeping/service_test.go && gofmt -l internal/housekeeping/ && go vet ./internal/housekeeping/ && go test ./internal/housekeeping/ -run 'TestCreateRealm|TestDeleteRealm' -v 2>&1 | tail -40
# github.com/astrate-platform/astrate/internal/housekeeping
# [github.com/astrate-platform/astrate/internal/housekeeping]
vet: internal/housekeeping/service_test.go:28:20: cannot use (*fakeStore)(nil) (value of type *fakeStore) as realmStore value in variable declaration: *fakeStore does not implement realmStore (missing method UpdateRealm)

The fake store needs the new method too.
← Edit internal/housekeeping/service_test.go
Index: /root/astrate-mule/internal/housekeeping/service_test.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service_test.go
+++ /root/astrate-mule/internal/housekeeping/service_test.go
@@ -50,8 +50,12 @@
 	}
 	return f.realm, nil
 }
 
+func (f *fakeStore) UpdateRealm(_ context.Context, _ string, _ store.RealmPatch) error {
+	return f.createErr
+}
+
 func (f *fakeStore) DeviceStats(_ context.Context, _ int16) (total, connected int64, err error) {
 	if f.statsErr != nil {
 		return 0, 0, f.statsErr
 	}

← Edit internal/housekeeping/service_test.go
Index: /root/astrate-mule/internal/housekeeping/service_test.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service_test.go
+++ /root/astrate-mule/internal/housekeeping/service_test.go
@@ -17,8 +17,9 @@
 total, connected int64
 statsErr         error
 
 createErr error
+updateErr error
 deleteErr error
 
 created []store.NewRealm
 deleted []string

← Edit internal/housekeeping/service_test.go
Index: /root/astrate-mule/internal/housekeeping/service_test.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service_test.go
+++ /root/astrate-mule/internal/housekeeping/service_test.go
@@ -52,9 +52,9 @@
 	return f.realm, nil
 }
 
 func (f *fakeStore) UpdateRealm(_ context.Context, _ string, _ store.RealmPatch) error {
-	return f.createErr
+	return f.updateErr
 }
 
 func (f *fakeStore) DeviceStats(_ context.Context, _ int16) (total, connected int64, err error) {
 	if f.statsErr != nil {

$ gofmt -w internal/housekeeping/service_test.go && gofmt -l internal/housekeeping/ && go vet ./internal/housekeeping/ && go test ./internal/housekeeping/ -run 'TestCreateRealm|TestDeleteRealm' -v 2>&1 | tail -45
# github.com/astrate-platform/astrate/internal/housekeeping
# [github.com/astrate-platform/astrate/internal/housekeeping]
vet: internal/housekeeping/service.go:230:22: s.st.ListRealms undefined (type realmStore has no field or method ListRealms)

← Edit internal/housekeeping/service.go
Index: /root/astrate-mule/internal/housekeeping/service.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service.go
+++ /root/astrate-mule/internal/housekeeping/service.go
@@ -54,8 +54,9 @@
 type realmStore interface {
 	CreateRealm(ctx context.Context, nr store.NewRealm) (*store.Realm, error)
 	GetRealmByName(ctx context.Context, name string) (*store.Realm, error)
 	UpdateRealm(ctx context.Context, name string, p store.RealmPatch) error
+	ListRealms(ctx context.Context) ([]store.Realm, error)
 	DeviceStats(ctx context.Context, realmID int16) (total, connected int64, err error)
 	DeleteRealm(ctx context.Context, name string) error
 }
 

← Edit internal/housekeeping/service_test.go
Index: /root/astrate-mule/internal/housekeeping/service_test.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service_test.go
+++ /root/astrate-mule/internal/housekeeping/service_test.go
@@ -55,8 +55,12 @@
 func (f *fakeStore) UpdateRealm(_ context.Context, _ string, _ store.RealmPatch) error {
 	return f.updateErr
 }
 
+func (f *fakeStore) ListRealms(_ context.Context) ([]store.Realm, error) {
+	return nil, nil
+}
+
 func (f *fakeStore) DeviceStats(_ context.Context, _ int16) (total, connected int64, err error) {
 	if f.statsErr != nil {
 		return 0, 0, f.statsErr
 	}

$ gofmt -w internal/housekeeping/service_test.go && gofmt -l internal/housekeeping/ && go vet ./internal/housekeeping/ && go test ./internal/housekeeping/ -run 'TestCreateRealm|TestDeleteRealm' -v 2>&1 | tail -45
=== RUN   TestCreateRealmValidation
=== RUN   TestCreateRealmValidation/blank_name
=== RUN   TestCreateRealmValidation/blank_jwt_key
=== RUN   TestCreateRealmValidation/negative_registration_limit
=== RUN   TestCreateRealmValidation/negative_retention
--- PASS: TestCreateRealmValidation (0.00s)
    --- PASS: TestCreateRealmValidation/blank_name (0.00s)
    --- PASS: TestCreateRealmValidation/blank_jwt_key (0.00s)
    --- PASS: TestCreateRealmValidation/negative_registration_limit (0.00s)
    --- PASS: TestCreateRealmValidation/negative_retention (0.00s)
=== RUN   TestCreateRealmDefaultRetention
=== RUN   TestCreateRealmDefaultRetention/omitted_injects_default
=== RUN   TestCreateRealmDefaultRetention/explicit_value_beats_default
=== RUN   TestCreateRealmDefaultRetention/no_default_stays_nil
--- PASS: TestCreateRealmDefaultRetention (0.01s)
    --- PASS: TestCreateRealmDefaultRetention/omitted_injects_default (0.00s)
    --- PASS: TestCreateRealmDefaultRetention/explicit_value_beats_default (0.00s)
    --- PASS: TestCreateRealmDefaultRetention/no_default_stays_nil (0.00s)
=== RUN   TestDeleteRealmGating
=== RUN   TestDeleteRealmGating/deletion_disabled_short-circuits
=== RUN   TestDeleteRealmGating/connected_device_blocks_delete
=== RUN   TestDeleteRealmGating/no_connected_devices_deletes
--- PASS: TestDeleteRealmGating (0.00s)
    --- PASS: TestDeleteRealmGating/deletion_disabled_short-circuits (0.00s)
    --- PASS: TestDeleteRealmGating/connected_device_blocks_delete (0.00s)
    --- PASS: TestDeleteRealmGating/no_connected_devices_deletes (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/housekeeping	0.027s

All new tests pass, gofmt clean. Now the full gate — `go vet ./...` + `go test ./...`, plus a `-tags integration` vet so the integration-tagged `http_test.go` is typechecked against the new `NewService` signature.
$ go vet ./... 2>&1 | tail -5 && go test ./... 2>&1 | tail -15
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	0.033s
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

$ go vet --tags integration ./internal/housekeeping/ 2>&1 | tail -5; echo "vet tags exit: $?"
vet tags exit: 0

$ gofmt -l . && go test -race ./internal/housekeeping/ 2>&1 | tail -3
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/housekeeping	0.023s
FAIL

Race detector can't initialize on this kernel (39-bit VMA — same limitation as the Pi), and `gofmt -l .` is clean (no output before the FAIL). Both the normal gate commands passed. Done.

```
Done: Added container-free unit tests for internal/housekeeping CreateRealm validation + default-retention injection, and DeleteRealm gating, via package-private store/sealer ports
Files: internal/housekeeping/service.go, internal/housekeeping/service_test.go
Verified: go vet ./... -> pass; go test ./... -> pass; go vet --tags integration ./internal/housekeeping/ -> pass; gofmt -l . -> pass (clean); go test -race ./internal/housekeeping/ -> fail (kernel: Found 39 - Supported 48, race detector cannot initialize on this 39-bit-VMA host, like the Pi; no changd concurrency here)
Unsure: NewService's params were widened from *store.Store/*store.KeySealer to package-private realmStore/sealer ports (the only way to fake them without Docker); every existing caller still compiles as-is (cmd/astrate, cpc, security, http_test.go all typecheck). Also added two flip-side subcases beyond the literal list — explicit-value-beats-default for retention, and clean-delete-when-disconnected for the gate.
Follow-ups: none
```
