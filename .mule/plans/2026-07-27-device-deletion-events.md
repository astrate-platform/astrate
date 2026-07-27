# Plan: device-deletion-events (issue-21)

## Task
Emit `device_deletion_started` + `device_deletion_finished` events around the
synchronous `Service.DeleteDevice`.

## Interpretation
The stream bus (`internal/engine/stream/bus.go`) defines lifecycle event kinds
(`device_connected`, `device_disconnected`, `device_error`). Two new kinds must
be added for deletion lifecycle. The events must fire before and after the actual
store deletion in `Service.DeleteDevice` (`internal/realm/service.go:76`).

## Approach

### 1. New event kinds — `internal/engine/stream/bus.go`
Add two constants:
```go
KindDeviceDeletionStarted  = "device_deletion_started"
KindDeviceDeletionFinished = "device_deletion_finished"
```

### 2. Emitter seam — `internal/realm/service.go`
The realm package can't import `engine` (it already imports `engine/triggers`,
which would cycle). Define a callback interface following the `Disconnecter`
pattern:

```go
type DeletionEmitter interface {
    EmitDeviceDeletionStart(realm, deviceID string)
    EmitDeviceDeletionFinish(realm, deviceID string)
}
```

Add to `Service`:
- `emitter DeletionEmitter` field
- `WithDeletionEmitter(DeletionEmitter) *Service` setter (nil-safe, returns self)

Modify `DeleteDevice`:
```go
func (s *Service) DeleteDevice(ctx context.Context, realm, deviceID string) error {
    rid, err := s.realmID(ctx, realm)
    if err != nil { return err }
    id, err := deviceid.Parse(deviceID)
    if err != nil { return fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID) }

    if s.emitter != nil {
        s.emitter.EmitDeviceDeletionStart(realm, deviceID)
    }
    if s.disc != nil {
        s.disc.DisconnectDevice(realm, id)
    }
    if err := s.st.DeleteDevice(ctx, rid, id); err != nil {
        if s.emitter != nil {
            s.emitter.EmitDeviceDeletionFinish(realm, deviceID)
        }
        return err
    }
    if s.emitter != nil {
        s.emitter.EmitDeviceDeletionFinish(realm, deviceID)
    }
    return nil
}
```

Key: `device_deletion_finished` fires even on error — a started deletion is
always finished (upstream parity: the event signals lifecycle completion, not
success).

### 3. Engine adapter — `cmd/astrate/main.go`
Define a small adapter (unexported, local to main) that publishes to the bus:
```go
type busEmitter struct{ bus *stream.Bus }

func (e *busEmitter) EmitDeviceDeletionStart(realm, deviceID string) {
    e.bus.Publish(stream.Event{
        Kind: stream.KindDeviceDeletionStarted, Realm: realm, DeviceID: deviceID,
        Timestamp: time.Now(),
    })
}
func (e *busEmitter) EmitDeviceDeletionFinish(realm, deviceID string) {
    e.bus.Publish(stream.Event{
        Kind: stream.KindDeviceDeletionFinished, Realm: realm, DeviceID: deviceID,
        Timestamp: time.Now(),
    })
}
```

Wire it:
```go
realm.NewAPI(
    realm.NewService(st, e, log).
        WithDisconnecter(b).
        WithDeletionEmitter(&busEmitter{bus: e.Bus()}),
    mw,
).Mount(mux)
```

### 4. Test — `internal/realm/dashboard_compat_test.go`
Extend the existing `TestDashboardCompat/DeviceDeletion` subtest (or add a new
subtest `DeviceDeletionEvents`). Define a `fakeDeletionEmitter` that records
calls, wire it via `WithDeletionEmitter`, and assert:
- `EmitDeviceDeletionStart` was called before the store delete
- `EmitDeviceDeletionFinish` was called after
- Both received the correct realm and device ID

This is an integration test (`//go:build integration`) because
`Service.DeleteDevice` needs a live database. The Pi gate
(`go test ./...`) won't run it; it's for Legion.

**Test that fails without the change:** The fake emitter's call log is
verified after `DeleteDevice`. Without the emitter calls in `DeleteDevice`,
the log is empty and the assertions fail.

## Files changed
1. `internal/engine/stream/bus.go` — 2 new constants
2. `internal/realm/service.go` — interface + field + setter + calls in DeleteDevice
3. `cmd/astrate/main.go` — adapter type + wiring
4. `internal/realm/dashboard_compat_test.go` — test subtest

## Verification
```bash
go vet ./...
go test ./...          # gate on Pi (no integration tag)
```
Integration test (Legion):
```bash
go test -tags integration ./internal/realm/...
```

## Unsure
- Whether `device_deletion_finished` should fire on error or only on success.
  I chose "always" (started ⟹ finished) for upstream parity. If it should be
  success-only, the error-path `EmitDeviceDeletionFinish` call should be removed.
