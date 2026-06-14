// Package pysdk is the CP-D official-Python-SDK conformance runner
// (docs/ROADMAP.md §10 file 9.3). It composes a live Astrate instance,
// registers a device, and drives the *unmodified* pinned
// astarte-device-sdk-python (requirements.txt) through the device loop by
// shelling out to runner.py — mirroring how the astartectl checkpoints invoke
// the pinned CLI. The persisted rows are cross-checked here.
//
// The SDK has platform-specific native dependencies, so this runner is gated on
// availability: it skips with instructions when python3 or the SDK is missing
// (the Linux CI/nightly job installs requirements.txt). Set ASTRATE_PYSDK_PYTHON
// to a venv interpreter that has the SDK to run it locally.
package pysdk

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/test/conformance/instance"
)

const (
	ifSensor = "org.astrate.pysdk.Sensor" // device datastream, individual
	ifConf   = "org.astrate.pysdk.Conf"   // device properties
)

var defs = map[string]string{
	ifSensor: `{"interface_name":"` + ifSensor + `","version_major":1,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`,
	ifConf:   `{"interface_name":"` + ifConf + `","version_major":1,"version_minor":0,"type":"properties","ownership":"device","mappings":[{"endpoint":"/%{k}","type":"string","allow_unset":true}]}`,
}

// pythonWithSDK resolves an interpreter that can import the Astarte SDK, or
// skips the test.
func pythonWithSDK(t *testing.T) string {
	t.Helper()
	py := os.Getenv("ASTRATE_PYSDK_PYTHON")
	if py == "" {
		py = "python3"
	}
	if _, err := exec.LookPath(py); err != nil {
		t.Skipf("pysdk: %q not found (set ASTRATE_PYSDK_PYTHON)", py)
	}
	if out, err := exec.Command(py, "-c", "import astarte.device").CombinedOutput(); err != nil {
		t.Skipf("pysdk: astarte-device-sdk not importable by %s — install test/conformance/pysdk/requirements.txt (%s)", py, out)
	}
	return py
}

func TestPythonSDK(t *testing.T) {
	py := pythonWithSDK(t)
	in := instance.New(t, instance.Config{Interfaces: defs})
	ctx := context.Background()

	id, err := deviceid.Random()
	if err != nil {
		t.Fatal(err)
	}
	secret := in.Register(t, id.String(), "")

	dir := t.TempDir()
	dsFile := filepath.Join(dir, "sensor.json")
	propFile := filepath.Join(dir, "conf.json")
	if err := os.WriteFile(dsFile, []byte(defs[ifSensor]), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(propFile, []byte(defs[ifConf]), 0o600); err != nil {
		t.Fatal(err)
	}

	_, file, _, _ := runtime.Caller(0)
	runner := filepath.Join(filepath.Dir(file), "runner.py")
	cmd := exec.Command(py, runner,
		"--pairing-url", in.PairingURL,
		"--realm", in.RealmName,
		"--device-id", id.String(),
		"--secret", secret,
		"--persistency-dir", filepath.Join(dir, "persistency"),
		"--datastream-interface", dsFile,
		"--properties-interface", propFile,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("python runner failed: %v\n%s", err, out)
	}

	// The SDK published an individual datastream and set a property.
	waitFor(t, 15*time.Second, func() bool {
		rows, err := in.Store.Series(ctx, store.SeriesQuery{
			RealmID: in.Realm.ID, DeviceID: id, InterfaceID: in.Interfaces[ifSensor].ID, Path: "/value",
		})
		return err == nil && len(rows) == 1 && rows[0].ValueDouble != nil && *rows[0].ValueDouble == 42.5
	})
	waitFor(t, 15*time.Second, func() bool {
		p, err := in.Store.GetProperty(ctx, in.Realm.ID, id, in.Interfaces[ifConf].ID, "/mode")
		return err == nil && string(p.Value) == `"eco"`
	})
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
