//go:build integration && e2e

package appengine

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/netip"
	"path/filepath"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/engine"
	"github.com/astrate-platform/astrate/internal/pairing"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
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
	realm  *store.Realm
	roots  *x509.CertPool
	sslURL string
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
	realm, err := st.CreateRealm(ctx, store.NewRealm{
		Name: realmName, CACertificatePEM: certPEM, CAPrivateKeySealed: sealedKey,
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
		st: st, svc: NewService(st, e, cdLogger()), pairer: pairer, broker: b,
		realm: realm, roots: roots, sslURL: "ssl://" + b.TLSAddr(),
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
