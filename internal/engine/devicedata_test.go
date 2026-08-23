package engine

// Device-owned ingest (issue #84): PublishDeviceValue must land data
// exactly like a real device's data lands — storage rows visible to
// AppEngine, zero MQTT traffic.

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// deviceStreamDef mirrors serverStreamDef as a device-owned datastream:
// same mapping shape, ownership flipped.
const deviceStreamDef = `{
	"interface_name": "com.astrate.test.DeviceStream",
	"version_major": 1, "version_minor": 0,
	"type": "datastream", "ownership": "device",
	"mappings": [{"endpoint": "/value", "type": "string", "reliability": "guaranteed", "expiry": 60}]
}`

// ifaceDeviceStream is deviceStreamDef's synthetic storage ID.
const ifaceDeviceStream = int64(19)

// devicePropsDef is an inline device-owned property with allow_unset.
const devicePropsDef = `{
	"interface_name": "com.astrate.test.DeviceProperties",
	"version_major": 1, "version_minor": 0,
	"type": "properties", "ownership": "device",
	"mappings": [{"endpoint": "/state", "type": "string", "allow_unset": true}]
}`

// ifaceDeviceProps is devicePropsDef's synthetic storage ID.
const ifaceDeviceProps = int64(20)

// addDeviceStream installs deviceStreamDef mid-run through the public
// invalidation callback.
func addDeviceStream(t *testing.T, rig *pipelineRig, fs *fakeStore) {
	t.Helper()
	fs.addInterface(realmAlphaID, storedInterface(t, realmAlphaID, ifaceDeviceStream, []byte(deviceStreamDef)))
	if err := rig.e.RefreshInterfaces(context.Background(), realmAlphaID); err != nil {
		t.Fatalf("RefreshInterfaces: %v", err)
	}
}

// addDeviceProps installs devicePropsDef mid-run through the public
// invalidation callback.
func addDeviceProps(t *testing.T, rig *pipelineRig, fs *fakeStore) {
	t.Helper()
	fs.addInterface(realmAlphaID, storedInterface(t, realmAlphaID, ifaceDeviceProps, []byte(devicePropsDef)))
	if err := rig.e.RefreshInterfaces(context.Background(), realmAlphaID); err != nil {
		t.Fatalf("RefreshInterfaces: %v", err)
	}
}

// TestPublishDeviceDatastream: a device-owned value persists as one typed
// row with the explicit ts honoured — and produces no MQTT traffic at all,
// which is the virtual-device contract.
func TestPublishDeviceDatastream(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})
	addDeviceStream(t, rig, fs)
	ts := time.Date(2026, 6, 2, 9, 15, 0, 0, time.UTC)

	if err := rig.e.PublishDeviceValue(context.Background(), realmAlpha, devAlpha,
		"com.astrate.test.DeviceStream", "/value", json.RawMessage(`"hello"`), &ts); err != nil {
		t.Fatalf("PublishDeviceValue: %v", err)
	}

	rows := fs.individualRows()
	if len(rows) != 1 {
		t.Fatalf("individual rows: %d, want 1", len(rows))
	}
	if rows[0].InterfaceID != ifaceDeviceStream || rows[0].ValueString == nil || *rows[0].ValueString != "hello" {
		t.Errorf("row: %+v", rows[0])
	}
	if !rows[0].TS.Equal(ts) {
		t.Errorf("row ts %s, want %s", rows[0].TS, ts)
	}
	if pubs := port.published(); len(pubs) != 0 {
		t.Errorf("device-owned ingest published to the broker: %+v", pubs)
	}
}

// TestPublishDeviceProperty: a device-owned property upserts via the
// property store like any property write, again with zero broker traffic.
func TestPublishDeviceProperty(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})
	addDeviceProps(t, rig, fs)

	if err := rig.e.PublishDeviceValue(context.Background(), realmAlpha, devAlpha,
		"com.astrate.test.DeviceProperties", "/state", json.RawMessage(`"on"`), nil); err != nil {
		t.Fatalf("PublishDeviceValue: %v", err)
	}

	fs.mu.Lock()
	upserts := len(fs.upserts)
	row := fs.upserts[0]
	fs.mu.Unlock()
	if upserts != 1 {
		t.Fatalf("upserts: %d, want 1", upserts)
	}
	if row.InterfaceID != ifaceDeviceProps || row.Path != "/state" || string(row.Value) != `"on"` {
		t.Errorf("upserted row: %+v", row)
	}
	if pubs := port.published(); len(pubs) != 0 {
		t.Errorf("device-owned property publish hit the broker: %+v", pubs)
	}
}

// TestPublishDeviceValueErrors: every refusal pairs with the acceptance
// tests above, so a blanket refusal could not pass them; each case names
// its rule in the assertion. Server-owned names are refused with
// ErrNotDeviceOwned — the flip side of PublishServerValue's gate.
func TestPublishDeviceValueErrors(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})
	addDeviceStream(t, rig, fs)
	addServerStream(t, rig, fs)
	ctx := context.Background()
	const stream = "com.astrate.test.DeviceStream"
	ghost := deviceid.ID{0x01, 0x02}

	cases := []struct {
		name    string
		realm   string
		device  deviceid.ID
		iface   string
		path    string
		value   string
		wantErr error
	}{
		{name: "realm unknown", realm: "ghost", device: devAlpha,
			iface: stream, path: "/value", value: `"x"`, wantErr: ErrRealmUnknown},
		{name: "device unknown", realm: realmAlpha, device: ghost,
			iface: stream, path: "/value", value: `"x"`, wantErr: store.ErrNotFound},
		{name: "interface not installed", realm: realmAlpha, device: devAlpha,
			iface: "com.example.Nope", path: "/value", value: `"x"`, wantErr: ErrInterfaceNotFound},
		{name: "server-owned interface", realm: realmAlpha, device: devAlpha,
			iface: "com.astrate.test.ServerStream", path: "/cmd", value: `"x"`, wantErr: ErrNotDeviceOwned},
		{name: "path not found", realm: realmAlpha, device: devAlpha,
			iface: stream, path: "/nope", value: `"x"`, wantErr: ErrPathNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := rig.e.PublishDeviceValue(ctx, tc.realm, tc.device, tc.iface, tc.path,
				json.RawMessage(tc.value), nil)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("err = %v, want %v", err, tc.wantErr)
			}
		})
	}

	// JSON null is never a value (§3.5.3): it surfaces as the same payload
	// rejection an inbound device payload would get, not as a sentinel.
	t.Run("null value", func(t *testing.T) {
		err := rig.e.PublishDeviceValue(ctx, realmAlpha, devAlpha, stream, "/value",
			json.RawMessage("null"), nil)
		if got := payload.ReasonOf(err); got != payload.ReasonTypeMismatch {
			t.Errorf("err = %v (%s), want the type_mismatch payload rejection", err, got)
		}
	})
	// The empty value is refused as an unset: only the empty-payload unset
	// of a real device may clear state, through its own paths.
	t.Run("empty value", func(t *testing.T) {
		err := rig.e.PublishDeviceValue(ctx, realmAlpha, devAlpha, stream, "/value",
			json.RawMessage(""), nil)
		if !errors.Is(err, ErrUnsetNotAllowed) {
			t.Errorf("err = %v, want %v", err, ErrUnsetNotAllowed)
		}
	})

	// Nothing was persisted or published by the failed attempts.
	fs.mu.Lock()
	upserts := len(fs.upserts)
	fs.mu.Unlock()
	if upserts != 0 || fs.batchCount() != 0 || len(port.published()) != 0 {
		t.Errorf("failed publishes left traces: %d upserts, %d batches, %d publishes",
			upserts, fs.batchCount(), len(port.published()))
	}
}
