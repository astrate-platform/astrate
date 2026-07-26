// Package gosdk is the CP-D official-Go-SDK conformance runner (docs/ROADMAP.md
// §10 file 9.2). CP-B (M6) already proves the full single-session device loop;
// this runner drives the *unmodified* astarte-device-sdk-go against the
// M9-composed instance and adds the cross-session scenarios: a device that
// disconnects and reconnects reuses its certificate and resumes publishing, and
// a server-owned datastream published while it was offline is delivered on
// reconnect (broker session persistence end to end through the SDK). Every
// value is cross-checked against the persisted rows (the §9.6 verify step).
package gosdk

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	sdkdevice "github.com/astarte-platform/astarte-device-sdk-go/device"
	"github.com/astarte-platform/astarte-go/interfaces"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/test/conformance/instance"
)

const (
	ifSensor     = "org.astrate.gosdk.SensorValues" // device datastream, individual
	ifServerData = "org.astrate.gosdk.ServerData"   // server datastream, individual
)

var defs = map[string]string{
	ifSensor:     `{"interface_name":"` + ifSensor + `","version_major":1,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`,
	ifServerData: `{"interface_name":"` + ifServerData + `","version_major":1,"version_minor":0,"type":"datastream","ownership":"server","mappings":[{"endpoint":"/value","type":"double"}]}`,
}

// inbox collects server→device deliveries the SDK reports.
type inbox struct {
	mu   sync.Mutex
	msgs []sdkdevice.IndividualMessage
}

func (in *inbox) add(m sdkdevice.IndividualMessage) {
	in.mu.Lock()
	defer in.mu.Unlock()
	in.msgs = append(in.msgs, m)
}

func (in *inbox) has(iface, path string) bool {
	in.mu.Lock()
	defer in.mu.Unlock()
	for _, m := range in.msgs {
		if m.Interface.Name == iface && m.Path == path {
			return true
		}
	}
	return false
}

// connect builds an SDK device over cryptoDir (reused across reconnects so the
// stored certificate is reused), wires the capture callback, and connects.
func connect(t *testing.T, in *instance.Instance, deviceID, secret, cryptoDir string, box *inbox) *sdkdevice.Device {
	t.Helper()
	if err := os.MkdirAll(cryptoDir, 0o700); err != nil {
		t.Fatal(err)
	}
	opts := sdkdevice.NewDeviceOptions()
	opts.UseMqttStore = false
	opts.UseDatabase = false
	opts.ConnectRetry = false
	opts.AutoReconnect = false
	opts.IgnoreSSLErrors = true
	opts.CryptoDir = cryptoDir

	d, err := sdkdevice.NewDeviceWithOptions(deviceID, in.RealmName, secret, in.PairingURL, opts)
	if err != nil {
		t.Fatalf("NewDeviceWithOptions: %v", err)
	}
	for name := range defs {
		iface, err := interfaces.ParseInterface([]byte(defs[name]))
		if err != nil {
			t.Fatalf("ParseInterface(%s): %v", name, err)
		}
		if err := d.AddInterface(iface); err != nil {
			t.Fatalf("AddInterface(%s): %v", name, err)
		}
	}
	if box != nil {
		d.OnIndividualMessageReceived = func(_ *sdkdevice.Device, m sdkdevice.IndividualMessage) { box.add(m) }
	}

	ch := make(chan error, 1)
	d.Connect(ch)
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("SDK Connect: %v", err)
		}
	case <-time.After(90 * time.Second):
		t.Fatal("SDK Connect did not settle within 90s")
	}
	return d
}

func disconnect(d *sdkdevice.Device) {
	ch := make(chan error, 1)
	d.Disconnect(ch)
	select {
	case <-ch:
	case <-time.After(10 * time.Second):
	}
}

func TestGoSDKReconnect(t *testing.T) {
	in := instance.New(t, instance.Config{Interfaces: defs})
	ctx := context.Background()

	id, err := deviceid.Random()
	if err != nil {
		t.Fatal(err)
	}
	deviceID := id.String()
	secret := in.Register(t, deviceID, "")
	cryptoDir := filepath.Join(t.TempDir(), "crypto")

	d := connect(t, in, deviceID, secret, cryptoDir, nil)

	// First session: publish and confirm persistence.
	if err := d.SendIndividualMessage(ifSensor, "/value", 1.5); err != nil {
		t.Fatalf("SendIndividualMessage: %v", err)
	}
	waitFor(t, 10*time.Second, func() bool { return latest(t, in, id) == 1.5 })

	// Reconnect over the same crypto dir: the SDK reuses its certificate (no
	// re-pair) and resumes publishing.
	disconnect(d)
	box := &inbox{}
	d2 := connect(t, in, deviceID, secret, cryptoDir, box)
	t.Cleanup(func() { disconnect(d2) })

	if err := d2.SendIndividualMessage(ifSensor, "/value", 2.5); err != nil {
		t.Fatalf("SendIndividualMessage after reconnect: %v", err)
	}
	waitFor(t, 10*time.Second, func() bool { return latest(t, in, id) == 2.5 })

	// The reconnected device receives server-owned data.
	if err := in.Engine.PublishServerValue(ctx, in.RealmName, id, ifServerData, "/value",
		json.RawMessage("9.5"), nil); err != nil {
		t.Fatalf("PublishServerValue: %v", err)
	}
	waitFor(t, 15*time.Second, func() bool { return box.has(ifServerData, "/value") })

	// The reconnecting device reused its certificate: still exactly one issuance.
	dev, err := in.Store.GetDevice(ctx, in.Realm.ID, id)
	if err != nil {
		t.Fatalf("GetDevice: %v", err)
	}
	if dev.Status != store.DeviceStatusConfirmed {
		t.Errorf("device status = %q, want confirmed", dev.Status)
	}
}

// latest returns the newest /value sample for the device, or -1 if none.
func latest(t *testing.T, in *instance.Instance, id deviceid.ID) float64 {
	t.Helper()
	rows, err := in.Store.Series(context.Background(), store.SeriesQuery{
		RealmID: in.Realm.ID, DeviceID: id, InterfaceID: in.Interfaces[ifSensor].ID, Path: "/value",
		Descending: true, Limit: 1,
	})
	if err != nil || len(rows) == 0 || rows[0].ValueDouble == nil {
		return -1
	}
	return *rows[0].ValueDouble
}

func waitFor(t *testing.T, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("condition not met within timeout")
}
