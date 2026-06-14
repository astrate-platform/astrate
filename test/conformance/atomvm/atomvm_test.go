// Package atomvm is the CP-D AtomVM-profile conformance runner (docs/ROADMAP.md
// §10 file 9.4, DESIGN §3.5): it drives Astrate with a device that speaks the
// documented plain-JSON payload profile instead of BSON — the AtomVM use case.
// It proves a JSON device is a first-class Astarte device: it registers with
// initial_payload_format "json", publishes JSON individual + object datastreams
// and properties (set + unset), and receives server-owned data JSON-encoded
// (payload_format_hint honoured), with every value cross-checked against the
// persisted rows.
package atomvm

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/payload"
	"github.com/astrate-platform/astrate/test/conformance/instance"
)

const (
	ifSensor     = "org.astrate.atomvm.Sensor"     // device datastream, individual
	ifCoords     = "org.astrate.atomvm.Coords"     // device datastream, object
	ifConf       = "org.astrate.atomvm.Conf"       // device properties
	ifServerData = "org.astrate.atomvm.ServerData" // server datastream, individual
)

var defs = map[string]string{
	ifSensor:     `{"interface_name":"` + ifSensor + `","version_major":1,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`,
	ifCoords:     `{"interface_name":"` + ifCoords + `","version_major":1,"version_minor":0,"type":"datastream","ownership":"device","aggregation":"object","mappings":[{"endpoint":"/coords/latitude","type":"double"},{"endpoint":"/coords/longitude","type":"double"}]}`,
	ifConf:       `{"interface_name":"` + ifConf + `","version_major":1,"version_minor":0,"type":"properties","ownership":"device","mappings":[{"endpoint":"/%{k}","type":"string","allow_unset":true}]}`,
	ifServerData: `{"interface_name":"` + ifServerData + `","version_major":1,"version_minor":0,"type":"datastream","ownership":"server","mappings":[{"endpoint":"/value","type":"double"}]}`,
}

func TestAtomVMJSONProfile(t *testing.T) {
	in := instance.New(t, instance.Config{Interfaces: defs})
	ctx := context.Background()

	// A JSON-only device: registered with the initial_payload_format extension
	// so server-owned messages arrive JSON-encoded from the first one.
	dev, id := in.NewDevice(t, "json")
	dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{
		ifSensor: {1, 0}, ifCoords: {1, 0}, ifConf: {1, 0}, ifServerData: {1, 0},
	}))
	waitFor(t, 10*time.Second, func() bool {
		d, err := in.Store.GetDevice(ctx, in.Realm.ID, id)
		return err == nil && len(d.Introspection) == 4
	})

	t.Run("IndividualDatastream", func(t *testing.T) {
		dev.PublishValue(t, ifSensor, "/value", 42.5, nil, payload.FormatJSON, 1)
		waitFor(t, 10*time.Second, func() bool {
			rows, err := in.Store.Series(ctx, store.SeriesQuery{
				RealmID: in.Realm.ID, DeviceID: id, InterfaceID: in.Interfaces[ifSensor].ID, Path: "/value",
			})
			return err == nil && len(rows) == 1 && rows[0].ValueDouble != nil && *rows[0].ValueDouble == 42.5
		})
	})

	t.Run("ObjectDatastream", func(t *testing.T) {
		dev.PublishValue(t, ifCoords, "/coords",
			map[string]any{"latitude": 45.07, "longitude": 7.69}, nil, payload.FormatJSON, 1)
		waitFor(t, 10*time.Second, func() bool {
			rows, err := in.Store.ObjectSeries(ctx, store.SeriesQuery{
				RealmID: in.Realm.ID, DeviceID: id, InterfaceID: in.Interfaces[ifCoords].ID, Path: "/coords",
			})
			return err == nil && len(rows) == 1 &&
				contains(string(rows[0].Value), "45.07") && contains(string(rows[0].Value), "7.69")
		})
	})

	t.Run("PropertySetUnset", func(t *testing.T) {
		ifaceID := in.Interfaces[ifConf].ID
		dev.PublishValue(t, ifConf, "/mode", "eco", nil, payload.FormatJSON, 1)
		waitFor(t, 10*time.Second, func() bool {
			p, err := in.Store.GetProperty(ctx, in.Realm.ID, id, ifaceID, "/mode")
			return err == nil && string(p.Value) == `"eco"`
		})
		// An empty payload unsets the property (JSON profile §5).
		dev.PublishRaw(t, ifConf, "/mode", nil, 1)
		waitFor(t, 10*time.Second, func() bool {
			_, err := in.Store.GetProperty(ctx, in.Realm.ID, id, ifaceID, "/mode")
			return errors.Is(err, store.ErrNotFound)
		})
	})

	t.Run("ServerOwnedArrivesJSON", func(t *testing.T) {
		if err := in.Engine.PublishServerValue(ctx, in.RealmName, id, ifServerData, "/value",
			json.RawMessage("9.5"), nil); err != nil {
			t.Fatalf("PublishServerValue: %v", err)
		}
		msg := dev.WaitForTopic(t, 15*time.Second, dev.Base()+"/"+ifServerData+"/value")
		if got := payload.DetectFormat(msg.Payload); got != payload.FormatJSON {
			t.Fatalf("server-owned payload format = %v, want JSON (%q)", got, msg.Payload)
		}
		var env struct {
			V float64 `json:"v"`
		}
		if err := json.Unmarshal(msg.Payload, &env); err != nil {
			t.Fatalf("server payload not JSON-decodable: %v (%q)", err, msg.Payload)
		}
		if env.V != 9.5 {
			t.Errorf("server-owned value = %v, want 9.5", env.V)
		}
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

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
