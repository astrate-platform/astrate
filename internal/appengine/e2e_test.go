//go:build integration && e2e

package appengine

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/golang-jwt/jwt/v5"

	"github.com/astrate-platform/astrate/internal/appengine/channels"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/engine"
	"github.com/astrate-platform/astrate/internal/pairing"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/payload"
)

const (
	cdServerData = "org.astrate.cd.ServerData" // server datastream
	cdDeviceData = "org.astrate.cd.DeviceData" // device datastream
)

var cdDefs = map[string]string{
	cdServerData: `{"interface_name":"org.astrate.cd.ServerData","version_major":1,"version_minor":0,"type":"datastream","ownership":"server","mappings":[{"endpoint":"/value","type":"double"}]}`,
	cdDeviceData: `{"interface_name":"org.astrate.cd.DeviceData","version_major":1,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`,
}

// cdEnv composes the cross-domain stack: store + engine + broker + pairing +
// the AppEngine service wired to the real engine as its ServerData port.
type cdEnv struct {
	st     *store.Store
	svc    *Service
	pairer *pairing.Service
	broker *broker.Broker
	engine *engine.Engine
	realm  *store.Realm
	roots  *x509.CertPool
	sslURL string
	jwtKey *rsa.PrivateKey
}

// cdDevice is a connected test device plus the material to reconnect it.
type cdDevice struct {
	*testutil.AstarteDevice
	id     deviceid.ID
	tlsCfg *tls.Config
	secret string
}

func newCDEnv(t *testing.T) (*cdEnv, func(t *testing.T) *cdDevice) {
	t.Helper()
	ctx := context.Background()

	pool := testutil.StartTimescale(t)
	st, err := store.New(ctx, pool.Config().ConnString())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(st.Close)

	master := make([]byte, store.MasterKeySize)
	_, _ = rand.Read(master)
	sealer, err := store.NewKeySealer(master)
	if err != nil {
		t.Fatal(err)
	}
	var suffix [4]byte
	_, _ = rand.Read(suffix[:])
	realmName := "cd" + hex.EncodeToString(suffix[:])
	certPEM, sealedKey, err := pairing.ProvisionCA(realmName, 0, sealer)
	if err != nil {
		t.Fatalf("ProvisionCA: %v", err)
	}
	// The Channels socket verifies its token against the realm's JWT keys, so
	// the cross-domain realm carries a real key pair.
	jwtKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKIXPublicKey(&jwtKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	jwtPub := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))

	realm, err := st.CreateRealm(ctx, store.NewRealm{
		Name: realmName, JWTPublicKeysPEM: []string{jwtPub},
		CACertificatePEM: certPEM, CAPrivateKeySealed: sealedKey,
	})
	if err != nil {
		t.Fatalf("CreateRealm: %v", err)
	}
	for _, def := range cdDefs {
		if _, err := st.InstallInterface(ctx, realm.ID, []byte(def)); err != nil {
			t.Fatalf("InstallInterface: %v", err)
		}
	}

	e, err := engine.New(st, nil, engine.Config{Shards: 2, BatchMaxRows: 4, BatchMaxWait: 20 * time.Millisecond, Logger: cdLogger()})
	if err != nil {
		t.Fatalf("engine.New: %v", err)
	}
	serverCert, roots := testutil.ServerTLSCert(t)
	b, err := broker.New(ctx, broker.Config{
		TLSAddr: "127.0.0.1:0", ServerTLSCert: serverCert,
		SessionStorePath: filepath.Join(t.TempDir(), "sessions.db"), Logger: cdLogger(),
	}, st, e, e)
	if err != nil {
		t.Fatalf("broker.New: %v", err)
	}
	e.AttachBroker(engine.AdaptBroker(b))
	runCtx, cancel := context.WithCancel(ctx)
	if err := e.Start(runCtx); err != nil {
		cancel()
		t.Fatalf("engine.Start: %v", err)
	}
	if err := b.Start(); err != nil {
		cancel()
		t.Fatalf("broker.Start: %v", err)
	}
	t.Cleanup(func() {
		_ = b.Close()
		dctx, dcancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer dcancel()
		_ = e.Drain(dctx)
		cancel()
	})

	pairer := pairing.New(st, sealer, pairing.Config{BrokerURL: "mqtts://" + b.TLSAddr()})
	env := &cdEnv{
		st: st, svc: NewService(st, e, cdLogger()), pairer: pairer, broker: b, engine: e,
		realm: realm, roots: roots, sslURL: "ssl://" + b.TLSAddr(), jwtKey: jwtKey,
	}

	newDevice := func(t *testing.T) *cdDevice {
		t.Helper()
		id, _ := deviceid.Random()
		secret, err := pairer.Register(ctx, realmName, id.String(), "")
		if err != nil {
			t.Fatalf("Register: %v", err)
		}
		key, csr := testutil.DeviceCSR(t)
		crt, err := pairer.Credentials(ctx, realmName, id.String(), secret, csr, netip.MustParseAddr("127.0.0.1"))
		if err != nil {
			t.Fatalf("Credentials: %v", err)
		}
		tlsCfg := testutil.DeviceTLSConfig(t, crt, key, roots)
		dev := testutil.ConnectAstarteDevice(t, env.sslURL, realmName, id, tlsCfg, true)
		t.Cleanup(dev.Disconnect)
		dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{
			cdServerData: {1, 0}, cdDeviceData: {1, 0},
		}))
		waitCond(t, 5*time.Second, func() bool {
			d, err := st.GetDevice(ctx, realm.ID, id)
			return err == nil && len(d.Introspection) == 2
		})
		return &cdDevice{AstarteDevice: dev, id: id, tlsCfg: tlsCfg, secret: secret}
	}

	return env, newDevice
}

// TestAppEngineCrossDomain is the M7b T3 suite (docs/ROADMAP.md §8.3): the
// AppEngine surface driving the live broker+engine through to a real device.
func TestAppEngineCrossDomain(t *testing.T) {
	env, newDevice := newCDEnv(t)
	ctx := context.Background()

	t.Run("ServerOwnedPutReachesDevice", func(t *testing.T) {
		dev := newDevice(t)
		if err := env.svc.PublishData(ctx, env.realm.Name, dev.id.String(), cdServerData, "/value",
			json.RawMessage("4.2"), nil); err != nil {
			t.Fatalf("AppEngine PublishData: %v", err)
		}
		// Delivery proof (value correctness is covered by the T2 + engine e2e).
		dev.WaitForTopic(t, 5*time.Second, dev.Base()+"/"+cdServerData+"/value")
	})

	t.Run("InhibitBlocksReconnect", func(t *testing.T) {
		dev := newDevice(t)
		if _, err := env.svc.PatchDevice(ctx, env.realm.Name, dev.id.String(),
			DevicePatch{CredentialsInhibited: ptrBool(true)}); err != nil {
			t.Fatalf("PatchDevice inhibit: %v", err)
		}
		waitCond(t, 5*time.Second, func() bool {
			d, err := env.st.GetDevice(ctx, env.realm.ID, dev.id)
			return err == nil && d.Status == store.DeviceStatusInhibited
		})

		// The broker must reject a fresh CONNECT from the inhibited device.
		dev.Disconnect()
		cn := env.realm.Name + "/" + dev.id.String()
		if _, _, err := testutil.MQTTTryConnect(t, env.sslURL, cn, true, dev.tlsCfg); err == nil {
			t.Error("inhibited device accepted by the broker after AppEngine PATCH")
		}

		// Pairing must also refuse new credentials for the inhibited device,
		// even with the correct secret.
		_, csr := testutil.DeviceCSR(t)
		if _, err := env.pairer.Credentials(ctx, env.realm.Name, dev.id.String(), dev.secret, csr,
			netip.MustParseAddr("127.0.0.1")); err == nil {
			t.Error("pairing issued credentials for an inhibited device")
		}
	})
}

// TestChannelsLiveEvents is phase 06's T3 (docs/COMPATIBILITY.md deviation 1):
// the Phoenix Channels socket mounted as the binary mounts it, driven with the
// exact frames the Dashboard's Device Live Events card sends, over a real
// device publishing through the real broker and engine.
func TestChannelsLiveEvents(t *testing.T) {
	env, newDevice := newCDEnv(t)
	dev := newDevice(t)

	mux := http.NewServeMux()
	channels.NewAPI(env.engine.Bus(), env.st).Mount(mux)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	claims := jwt.MapClaims{"a_ch": []string{"JOIN::.*", "WATCH::.*"}, "exp": time.Now().Add(time.Hour).Unix()}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(env.jwtKey)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	url := strings.Replace(srv.URL, "http", "ws", 1) +
		"/appengine/v1/socket/websocket?vsn=2.0.0&realm=" + env.realm.Name + "&token=" + token
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer conn.CloseNow()

	// The card's own frames: join rooms:<realm>:<name>, then one watch per
	// trigger. Only the incoming_data watch is exercised end to end here.
	topic := "rooms:" + env.realm.Name + ":dashboard_" + dev.id.String() + "_4242"
	write := func(frame string) {
		t.Helper()
		if err := conn.Write(ctx, websocket.MessageText, []byte(frame)); err != nil {
			t.Fatalf("write %s: %v", frame, err)
		}
	}
	expectOK := func(what string) {
		t.Helper()
		_, raw, err := conn.Read(ctx)
		if err != nil {
			t.Fatalf("read reply to %s: %v", what, err)
		}
		var f [5]json.RawMessage
		if err := json.Unmarshal(raw, &f); err != nil {
			t.Fatalf("decode reply to %s (%s): %v", what, raw, err)
		}
		var payload struct {
			Status string `json:"status"`
		}
		if err := json.Unmarshal(f[4], &payload); err != nil || payload.Status != "ok" {
			t.Fatalf("reply to %s: %s", what, raw)
		}
	}

	write(`["1","1","` + topic + `","phx_join",{}]`)
	expectOK("phx_join")

	watch := `{"name":"data","device_id":"` + dev.id.String() + `","simple_trigger":` +
		`{"type":"data_trigger","on":"incoming_data","interface_name":"*",` +
		`"value_match_operator":"*","match_path":"/*"}}`
	write(`["1","2","` + topic + `","watch",` + watch + `]`)
	expectOK("watch")

	// Give the room's bus subscription a moment before the device publishes;
	// the bus drops events with no subscriber rather than queueing them.
	time.Sleep(200 * time.Millisecond)
	dev.PublishValue(t, cdDeviceData, "/value", 4.2, nil, payload.FormatBSON, 1)

	// The card decodes new_event by trying its validators in turn, so the
	// envelope's field names are what must match, not merely the value.
	for {
		_, raw, err := conn.Read(ctx)
		if err != nil {
			t.Fatalf("read new_event: %v", err)
		}
		var f [5]json.RawMessage
		if err := json.Unmarshal(raw, &f); err != nil {
			t.Fatalf("decode frame %s: %v", raw, err)
		}
		var event string
		_ = json.Unmarshal(f[3], &event)
		if event != "new_event" {
			continue
		}
		var got struct {
			DeviceID  string `json:"device_id"`
			Timestamp string `json:"timestamp"`
			Event     struct {
				Type      string  `json:"type"`
				Interface string  `json:"interface"`
				Path      string  `json:"path"`
				Value     float64 `json:"value"`
			} `json:"event"`
		}
		if err := json.Unmarshal(f[4], &got); err != nil {
			t.Fatalf("decode new_event payload %s: %v", f[4], err)
		}
		if got.DeviceID != dev.id.String() || got.Timestamp == "" {
			t.Errorf("new_event envelope: %s", f[4])
		}
		if got.Event.Type != "incoming_data" || got.Event.Interface != cdDeviceData ||
			got.Event.Path != "/value" || got.Event.Value != 4.2 {
			t.Errorf("new_event body: %s", f[4])
		}
		return
	}
}

func cdLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func ptrBool(b bool) *bool { return &b }

func waitCond(t *testing.T, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("condition not met within timeout")
}
