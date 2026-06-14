package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// serverStreamDef is an inline server-owned datastream with non-default
// reliability and expiry, for the QoS/expiry delivery assertions.
const serverStreamDef = `{
	"interface_name": "com.astrate.test.ServerStream",
	"version_major": 1, "version_minor": 0,
	"type": "datastream", "ownership": "server",
	"mappings": [{"endpoint": "/cmd", "type": "string", "reliability": "guaranteed", "expiry": 60}]
}`

// ifaceServerStream is serverStreamDef's synthetic storage ID.
const ifaceServerStream = int64(17)

// addServerStream installs serverStreamDef mid-run through the public
// invalidation callback.
func addServerStream(t *testing.T, rig *pipelineRig, fs *fakeStore) {
	t.Helper()
	fs.addInterface(realmAlphaID, storedInterface(t, realmAlphaID, ifaceServerStream, []byte(serverStreamDef)))
	if err := rig.e.RefreshInterfaces(context.Background(), realmAlphaID); err != nil {
		t.Fatalf("RefreshInterfaces: %v", err)
	}
}

// TestPublishServerProperty: validate → upsert → retained QoS 2 publish in
// the device's hint format (docs/ROADMAP.md §7.2 file 6.9).
func TestPublishServerProperty(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})
	ctx := context.Background()
	const iface = "com.astrate.test.ServerProperties"

	if err := rig.e.PublishServerValue(ctx, realmAlpha, devAlpha, iface, "/limits/maxConnections",
		json.RawMessage("42"), nil); err != nil {
		t.Fatalf("PublishServerValue: %v", err)
	}

	fs.mu.Lock()
	upserts := len(fs.upserts)
	row := fs.upserts[0]
	fs.mu.Unlock()
	if upserts != 1 {
		t.Fatalf("upserts: %d, want 1", upserts)
	}
	if row.InterfaceID != ifaceServerProps || row.Path != "/limits/maxConnections" || string(row.Value) != "42" {
		t.Errorf("upserted row: %+v", row)
	}

	pubs := port.publishedTo(realmAlpha + "/" + devAlpha.String() + "/" + iface + "/limits/maxConnections")
	if len(pubs) != 1 {
		t.Fatalf("publishes: %+v", port.published())
	}
	want, err := payload.Encode(int32(42), nil, payload.FormatBSON)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pubs[0].payload, want) {
		t.Errorf("wire payload: %x, want %x", pubs[0].payload, want)
	}
	if pubs[0].qos != 2 || !pubs[0].retain || pubs[0].expiry != 0 {
		t.Errorf("property publish qos=%d retain=%t expiry=%s, want 2/true/0",
			pubs[0].qos, pubs[0].retain, pubs[0].expiry)
	}
}

// TestPublishServerPropertyJSONHint: a JSON-profile device receives server
// properties as JSON documents (docs/DESIGN.md §3.5.4).
func TestPublishServerPropertyJSONHint(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})
	devJSON := deviceid.ID{0xb2, 0x71, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	fs.addDevice(realmAlpha, realmAlphaID, devJSON, map[string]store.InterfaceVersion{
		"com.astrate.test.ServerProperties": {Major: 1, Minor: 2},
	}, hintJSON)

	if err := rig.e.PublishServerValue(context.Background(), realmAlpha, devJSON,
		"com.astrate.test.ServerProperties", "/identity/displayName",
		json.RawMessage(`"panel-9"`), nil); err != nil {
		t.Fatalf("PublishServerValue: %v", err)
	}

	pubs := port.published()
	if len(pubs) != 1 {
		t.Fatalf("publishes: %+v", pubs)
	}
	want, err := payload.Encode("panel-9", nil, payload.FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pubs[0].payload, want) {
		t.Errorf("wire payload: %s, want %s", pubs[0].payload, want)
	}
}

// TestPublishServerDatastream: datastream values persist as a single typed
// row and ship with the mapping's QoS and expiry, timestamp on the wire.
func TestPublishServerDatastream(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})
	addServerStream(t, rig, fs)
	ts := time.Date(2026, 6, 2, 8, 30, 0, 0, time.UTC)

	if err := rig.e.PublishServerValue(context.Background(), realmAlpha, devAlpha,
		"com.astrate.test.ServerStream", "/cmd", json.RawMessage(`"reboot"`), &ts); err != nil {
		t.Fatalf("PublishServerValue: %v", err)
	}

	rows := fs.individualRows()
	if len(rows) != 1 {
		t.Fatalf("individual rows: %d, want 1", len(rows))
	}
	if rows[0].InterfaceID != ifaceServerStream || rows[0].ValueString == nil || *rows[0].ValueString != "reboot" {
		t.Errorf("row: %+v", rows[0])
	}
	if !rows[0].TS.Equal(ts) {
		t.Errorf("row ts %s, want %s", rows[0].TS, ts)
	}

	pubs := port.publishedTo(realmAlpha + "/" + devAlpha.String() + "/com.astrate.test.ServerStream/cmd")
	if len(pubs) != 1 {
		t.Fatalf("publishes: %+v", port.published())
	}
	want, err := payload.Encode("reboot", &ts, payload.FormatBSON)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pubs[0].payload, want) {
		t.Errorf("wire payload: %x, want %x", pubs[0].payload, want)
	}
	if pubs[0].qos != 1 { // reliability: guaranteed
		t.Errorf("qos = %d, want 1", pubs[0].qos)
	}
	if pubs[0].retain {
		t.Error("datastream publish retained")
	}
	if pubs[0].expiry != 60*time.Second {
		t.Errorf("expiry = %s, want 60s", pubs[0].expiry)
	}
}

// TestPublishServerValueErrors: every docs/ROADMAP.md §7.2 file 6.9
// validation failure surfaces as its sentinel (or a payload rejection).
func TestPublishServerValueErrors(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})
	ctx := context.Background()
	const propsIface = "com.astrate.test.ServerProperties"
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
			iface: propsIface, path: "/limits/maxConnections", value: "42", wantErr: ErrRealmUnknown},
		{name: "device unknown", realm: realmAlpha, device: ghost,
			iface: propsIface, path: "/limits/maxConnections", value: "42", wantErr: store.ErrNotFound},
		{name: "interface not installed", realm: realmAlpha, device: devAlpha,
			iface: "com.example.Nope", path: "/x", value: "42", wantErr: ErrInterfaceNotFound},
		{name: "device-owned interface", realm: realmAlpha, device: devAlpha,
			iface: "com.astrate.test.Minimal", path: "/value", value: "1.5", wantErr: ErrNotServerOwned},
		{name: "path not found", realm: realmAlpha, device: devAlpha,
			iface: propsIface, path: "/limits/nope", value: "42", wantErr: ErrPathNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := rig.e.PublishServerValue(ctx, tc.realm, tc.device, tc.iface, tc.path,
				json.RawMessage(tc.value), nil)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("err = %v, want %v", err, tc.wantErr)
			}
		})
	}

	// Type mismatch — and JSON null, which the §3.5.3 profile rejects as a
	// value (unset is the empty payload / DELETE path) — surface as payload
	// rejections.
	for _, bad := range []string{`"not a number"`, "null"} {
		err := rig.e.PublishServerValue(ctx, realmAlpha, devAlpha, propsIface,
			"/limits/maxConnections", json.RawMessage(bad), nil)
		if payload.ReasonOf(err) == payload.ReasonNone {
			t.Errorf("value %s: err = %v, want a payload rejection", bad, err)
		}
	}

	// Nothing was persisted or published by the failed attempts.
	fs.mu.Lock()
	upserts := len(fs.upserts)
	fs.mu.Unlock()
	if upserts != 0 || fs.batchCount() != 0 || len(port.published()) != 0 {
		t.Errorf("failed publishes left traces: %d upserts, %d batches, %d publishes",
			upserts, fs.batchCount(), len(port.published()))
	}
}

// TestUnsetServerProperty: row delete, empty retained publish, and the
// consumer/properties purge follow-up (docs/DESIGN.md §3.4).
func TestUnsetServerProperty(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})
	ctx := context.Background()
	const iface = "com.astrate.test.ServerProperties"

	// Two properties; one gets unset.
	for path, val := range map[string]string{"/identity/displayName": `"p"`, "/limits/maxPayload": "9000"} {
		if err := rig.e.PublishServerValue(ctx, realmAlpha, devAlpha, iface, path,
			json.RawMessage(val), nil); err != nil {
			t.Fatalf("seed %s: %v", path, err)
		}
	}

	if err := rig.e.UnsetServerProperty(ctx, realmAlpha, devAlpha, iface, "/identity/displayName"); err != nil {
		t.Fatalf("UnsetServerProperty: %v", err)
	}

	fs.mu.Lock()
	unsets := len(fs.unsets)
	call := fs.unsets[0]
	fs.mu.Unlock()
	if unsets != 1 || call.interfaceID != ifaceServerProps || call.path != "/identity/displayName" {
		t.Errorf("unset call: %+v (count %d)", call, unsets)
	}

	base := realmAlpha + "/" + devAlpha.String()
	clears := port.publishedTo(base + "/" + iface + "/identity/displayName")
	last := clears[len(clears)-1]
	if len(last.payload) != 0 || !last.retain || last.qos != 2 {
		t.Errorf("clear publish: %d bytes retain=%t qos=%d, want empty/true/2",
			len(last.payload), last.retain, last.qos)
	}

	purges := port.publishedTo(base + "/control/consumer/properties")
	if len(purges) != 1 {
		t.Fatalf("purge messages after unset: %d, want 1", len(purges))
	}
	entries, err := inflateProperties(purges[0].payload)
	if err != nil {
		t.Fatalf("purge payload: %v", err)
	}
	if len(entries) != 1 || entries[0] != iface+"/limits/maxPayload" {
		t.Errorf("purge entries after unset: %v", entries)
	}

	// Guard rails.
	if err := rig.e.UnsetServerProperty(ctx, realmAlpha, devAlpha, iface, "/limits/maxConnections"); !errors.Is(err, ErrUnsetNotAllowed) {
		t.Errorf("unset without allow_unset: %v, want ErrUnsetNotAllowed", err)
	}
	addServerStream(t, rig, fs)
	if err := rig.e.UnsetServerProperty(ctx, realmAlpha, devAlpha, "com.astrate.test.ServerStream", "/cmd"); !errors.Is(err, ErrNotAProperty) {
		t.Errorf("unset on datastream: %v, want ErrNotAProperty", err)
	}
}

// TestServerObjectAggregate: object-aggregated server interfaces accept one
// document of last-level keys and persist a single object row.
func TestServerObjectAggregate(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})
	const def = `{
		"interface_name": "com.astrate.test.ServerObject",
		"version_major": 1, "version_minor": 0,
		"type": "datastream", "ownership": "server", "aggregation": "object",
		"mappings": [
			{"endpoint": "/setpoints/heating", "type": "double"},
			{"endpoint": "/setpoints/cooling", "type": "double"}
		]
	}`
	fs.addInterface(realmAlphaID, storedInterface(t, realmAlphaID, 18, []byte(def)))
	if err := rig.e.RefreshInterfaces(context.Background(), realmAlphaID); err != nil {
		t.Fatalf("RefreshInterfaces: %v", err)
	}

	if err := rig.e.PublishServerValue(context.Background(), realmAlpha, devAlpha,
		"com.astrate.test.ServerObject", "/setpoints",
		json.RawMessage(`{"heating": 21.5, "cooling": 24.0}`), nil); err != nil {
		t.Fatalf("PublishServerValue(object): %v", err)
	}

	rows := fs.objectRows()
	if len(rows) != 1 || rows[0].Path != "/setpoints" {
		t.Fatalf("object rows: %+v", rows)
	}
	var doc map[string]float64
	if err := json.Unmarshal(rows[0].Value, &doc); err != nil {
		t.Fatalf("object document: %v", err)
	}
	if doc["heating"] != 21.5 || doc["cooling"] != 24.0 {
		t.Errorf("object document values: %v", doc)
	}
	if len(port.published()) != 1 {
		t.Fatalf("publishes: %+v", port.published())
	}

	// Wrong prefix is a path error.
	err := rig.e.PublishServerValue(context.Background(), realmAlpha, devAlpha,
		"com.astrate.test.ServerObject", "/nope", json.RawMessage(`{"heating": 1.0}`), nil)
	if !errors.Is(err, ErrPathNotFound) {
		t.Errorf("bad prefix err = %v, want ErrPathNotFound", err)
	}
}
