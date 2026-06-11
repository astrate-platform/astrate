package pairing

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/netip"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/astrate-platform/astrate/internal/pairing/ca"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// fakeStore is an in-memory pairing.Store implementing the same contracts
// as the real repositories (rotation rules, sentinel errors).
type fakeStore struct {
	realms  map[string]*store.Realm
	devices map[int16]map[deviceid.ID]*store.Device
	hints   map[deviceid.ID]string
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		realms:  map[string]*store.Realm{},
		devices: map[int16]map[deviceid.ID]*store.Device{},
		hints:   map[deviceid.ID]string{},
	}
}

func (f *fakeStore) addRealm(r *store.Realm) *store.Realm {
	f.realms[r.Name] = r
	f.devices[r.ID] = map[deviceid.ID]*store.Device{}
	return r
}

func (f *fakeStore) GetRealmByName(_ context.Context, name string) (*store.Realm, error) {
	r, ok := f.realms[name]
	if !ok {
		return nil, fmt.Errorf("%w: realm %q", store.ErrNotFound, name)
	}
	return r, nil
}

func (f *fakeStore) RegisterDevice(_ context.Context, realmID int16, id deviceid.ID, secretHash string) error {
	devs := f.devices[realmID]
	if d, ok := devs[id]; ok {
		if d.FirstCredentialsRequest != nil {
			return fmt.Errorf("%w: device %s", store.ErrDeviceAlreadyConfirmed, id)
		}
		d.CredentialsSecretHash = secretHash
		d.Status = store.DeviceStatusRegistered
		return nil
	}
	devs[id] = &store.Device{
		ID:                    id,
		RealmID:               realmID,
		CredentialsSecretHash: secretHash,
		Status:                store.DeviceStatusRegistered,
		PayloadFormatHint:     "bson",
	}
	return nil
}

func (f *fakeStore) UnregisterDevice(_ context.Context, realmID int16, id deviceid.ID) error {
	d, ok := f.devices[realmID][id]
	if !ok {
		return fmt.Errorf("%w: device %s", store.ErrNotFound, id)
	}
	d.CredentialsSecretHash = ""
	d.Status = store.DeviceStatusRegistered
	d.FirstCredentialsRequest = nil
	d.CertSerial, d.CertAKI = nil, nil
	return nil
}

func (f *fakeStore) GetDevice(_ context.Context, realmID int16, id deviceid.ID) (*store.Device, error) {
	d, ok := f.devices[realmID][id]
	if !ok {
		return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, id)
	}
	cp := *d
	return &cp, nil
}

func (f *fakeStore) SetDeviceCredentials(_ context.Context, realmID int16, id deviceid.ID, certSerial, certAKI string, requestIP netip.Addr) error {
	d, ok := f.devices[realmID][id]
	if !ok {
		return fmt.Errorf("%w: device %s", store.ErrNotFound, id)
	}
	now := time.Now()
	if d.FirstCredentialsRequest == nil {
		d.FirstCredentialsRequest = &now
	}
	d.CertSerial, d.CertAKI = &certSerial, &certAKI
	d.LastCredentialsRequestIP = &requestIP
	d.Status = store.DeviceStatusConfirmed
	return nil
}

func (f *fakeStore) SetPayloadFormatHint(_ context.Context, _ int16, id deviceid.ID, hint string) error {
	f.hints[id] = hint
	return nil
}

func (f *fakeStore) CountDevices(_ context.Context, realmID int16) (int64, error) {
	return int64(len(f.devices[realmID])), nil
}

// --- test fixtures ---------------------------------------------------------

func newSealer(t *testing.T) *store.KeySealer {
	t.Helper()
	key := make([]byte, store.MasterKeySize)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	sealer, err := store.NewKeySealer(key)
	if err != nil {
		t.Fatal(err)
	}
	return sealer
}

// newServiceFixture builds a fake-backed Service with a provisioned realm CA.
func newServiceFixture(t *testing.T, cfg Config) (*Service, *fakeStore, *store.Realm) {
	t.Helper()
	sealer := newSealer(t)
	certPEM, sealed, err := ProvisionCA("test", 0, sealer)
	if err != nil {
		t.Fatalf("ProvisionCA: %v", err)
	}
	fs := newFakeStore()
	realm := fs.addRealm(&store.Realm{
		ID:                 1,
		Name:               "test",
		CACertificatePEM:   certPEM,
		CAPrivateKeySealed: sealed,
	})
	return New(fs, sealer, cfg), fs, realm
}

// deviceCSR builds a fresh EC key + CSR for a device.
func deviceCSR(t *testing.T) string {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: "ignored"},
	}, key)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: der}))
}

func randomDeviceID(t *testing.T) string {
	t.Helper()
	id, err := deviceid.Random()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

var testIP = netip.MustParseAddr("192.0.2.10")

// --- flow A ----------------------------------------------------------------

func TestRegister(t *testing.T) {
	ctx := context.Background()
	svc, fs, _ := newServiceFixture(t, Config{})
	hwID := randomDeviceID(t)

	secret, err := svc.Register(ctx, "test", hwID, "")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if len(secret) != 44 {
		t.Errorf("secret length: got %d, want 44", len(secret))
	}
	raw, err := base64.StdEncoding.DecodeString(secret)
	if err != nil || len(raw) != 32 {
		t.Errorf("secret must be standard base64 of 32 bytes: err=%v len=%d", err, len(raw))
	}

	id, _ := deviceid.Parse(hwID)
	dev := fs.devices[1][id]
	if dev == nil {
		t.Fatal("device row missing after registration")
	}
	if bcrypt.CompareHashAndPassword([]byte(dev.CredentialsSecretHash), []byte(secret)) != nil {
		t.Error("stored hash does not match the returned secret")
	}

	// Pre-credentials re-registration rotates the secret.
	secret2, err := svc.Register(ctx, "test", hwID, "")
	if err != nil {
		t.Fatalf("re-Register: %v", err)
	}
	if secret2 == secret {
		t.Error("re-registration must rotate the secret")
	}
	if bcrypt.CompareHashAndPassword([]byte(fs.devices[1][id].CredentialsSecretHash), []byte(secret2)) != nil {
		t.Error("stored hash not rotated")
	}

	// After the first credentials request, registration conflicts.
	if err := fs.SetDeviceCredentials(ctx, 1, id, "1", "aa", testIP); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Register(ctx, "test", hwID, ""); !errors.Is(err, ErrAlreadyRegistered) {
		t.Errorf("post-credentials re-register: got %v, want ErrAlreadyRegistered", err)
	}
}

func TestRegisterInvalidHWID(t *testing.T) {
	svc, _, _ := newServiceFixture(t, Config{})
	for _, hwID := range []string{
		"",
		"tooshort",
		"h4-Dx_RYTU-RbpDOTabhR",    // 21 characters
		"h4-Dx_RYTU-RbpDOTabhRg=",  // padded
		"h4+Dx/RYTU+RbpDOTabhRg",   // standard alphabet, not url-safe
		"h4-Dx_RYTU-RbpDOTabhRgg2", // 24 characters
	} {
		if _, err := svc.Register(context.Background(), "test", hwID, ""); !errors.Is(err, ErrInvalidHWID) {
			t.Errorf("Register(%q): got %v, want ErrInvalidHWID", hwID, err)
		}
	}
}

func TestRegisterInitialPayloadFormat(t *testing.T) {
	ctx := context.Background()
	svc, fs, _ := newServiceFixture(t, Config{})

	hwID := randomDeviceID(t)
	if _, err := svc.Register(ctx, "test", hwID, "json"); err != nil {
		t.Fatalf("Register with json format: %v", err)
	}
	id, _ := deviceid.Parse(hwID)
	if fs.hints[id] != "json" {
		t.Errorf("payload format hint: got %q, want json", fs.hints[id])
	}

	// Default: no hint call.
	hwID2 := randomDeviceID(t)
	if _, err := svc.Register(ctx, "test", hwID2, ""); err != nil {
		t.Fatal(err)
	}
	id2, _ := deviceid.Parse(hwID2)
	if _, ok := fs.hints[id2]; ok {
		t.Error("no hint must be recorded without initial_payload_format")
	}

	if _, err := svc.Register(ctx, "test", randomDeviceID(t), "cbor"); !errors.Is(err, ErrInvalidPayloadFormat) {
		t.Errorf("invalid format: got %v, want ErrInvalidPayloadFormat", err)
	}
}

func TestRegisterLimit(t *testing.T) {
	ctx := context.Background()
	svc, fs, realm := newServiceFixture(t, Config{})
	limit := int32(1)
	realm.DeviceRegistrationLimit = &limit

	if _, err := svc.Register(ctx, "test", randomDeviceID(t), ""); err != nil {
		t.Fatalf("first Register: %v", err)
	}
	if _, err := svc.Register(ctx, "test", randomDeviceID(t), ""); !errors.Is(err, ErrRegistrationLimitReached) {
		t.Errorf("over-limit Register: got %v, want ErrRegistrationLimitReached", err)
	}
	if got := len(fs.devices[1]); got != 1 {
		t.Errorf("device count: got %d, want 1", got)
	}
}

func TestUnregister(t *testing.T) {
	ctx := context.Background()
	svc, fs, _ := newServiceFixture(t, Config{})
	hwID := randomDeviceID(t)

	if err := svc.Unregister(ctx, "test", hwID); !errors.Is(err, store.ErrNotFound) {
		t.Errorf("unknown device: got %v, want store.ErrNotFound", err)
	}

	secret, err := svc.Register(ctx, "test", hwID, "")
	if err != nil {
		t.Fatal(err)
	}
	id, _ := deviceid.Parse(hwID)
	if err := fs.SetDeviceCredentials(ctx, 1, id, "1", "aa", testIP); err != nil {
		t.Fatal(err)
	}

	if err := svc.Unregister(ctx, "test", hwID); err != nil {
		t.Fatalf("Unregister: %v", err)
	}
	// The device is re-registrable, the old secret unusable.
	if _, err := svc.Register(ctx, "test", hwID, ""); err != nil {
		t.Errorf("re-register after unregister: %v", err)
	}
	if _, err := svc.Credentials(ctx, "test", hwID, secret, deviceCSR(t), testIP); !errors.Is(err, ErrUnauthorized) {
		t.Errorf("old secret after unregister: got %v, want ErrUnauthorized", err)
	}
}

// --- flow B ----------------------------------------------------------------

func TestCredentials(t *testing.T) {
	ctx := context.Background()
	svc, fs, realm := newServiceFixture(t, Config{EnforceLatestCert: true})
	hwID := randomDeviceID(t)
	id, _ := deviceid.Parse(hwID)

	secret, err := svc.Register(ctx, "test", hwID, "")
	if err != nil {
		t.Fatal(err)
	}

	certPEM, err := svc.Credentials(ctx, "test", hwID, secret, deviceCSR(t), testIP)
	if err != nil {
		t.Fatalf("Credentials: %v", err)
	}
	cert, err := ca.ParseCertificatePEM(certPEM)
	if err != nil {
		t.Fatalf("issued certificate does not parse: %v", err)
	}
	if got, want := cert.Subject.CommonName, "test/"+hwID; got != want {
		t.Errorf("CN: got %q, want %q", got, want)
	}

	dev := fs.devices[1][id]
	if dev.Status != store.DeviceStatusConfirmed {
		t.Errorf("status: got %q, want confirmed", dev.Status)
	}
	if dev.FirstCredentialsRequest == nil {
		t.Error("first_credentials_request not stamped")
	}
	if dev.CertSerial == nil || *dev.CertSerial != cert.SerialNumber.String() {
		t.Errorf("recorded serial: got %v, want %s", dev.CertSerial, cert.SerialNumber.String())
	}
	if dev.LastCredentialsRequestIP == nil || *dev.LastCredentialsRequestIP != testIP {
		t.Errorf("recorded IP: got %v, want %v", dev.LastCredentialsRequestIP, testIP)
	}

	// The certificate chains to the realm CA.
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM([]byte(realm.CACertificatePEM)) {
		t.Fatal("appending realm CA")
	}
	if _, err := cert.Verify(x509.VerifyOptions{
		Roots:     roots,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}); err != nil {
		t.Errorf("certificate does not chain to the realm CA: %v", err)
	}
}

func TestCredentialsUniformUnauthorized(t *testing.T) {
	ctx := context.Background()
	svc, _, _ := newServiceFixture(t, Config{})
	hwID := randomDeviceID(t)
	secret, err := svc.Register(ctx, "test", hwID, "")
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string]func() error{
		"wrong secret": func() error {
			_, err := svc.Credentials(ctx, "test", hwID, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", deviceCSR(t), testIP)
			return err
		},
		"unknown device": func() error {
			_, err := svc.Credentials(ctx, "test", randomDeviceID(t), secret, deviceCSR(t), testIP)
			return err
		},
		"unknown realm": func() error {
			_, err := svc.Credentials(ctx, "nosuch", hwID, secret, deviceCSR(t), testIP)
			return err
		},
		"malformed device id": func() error {
			_, err := svc.Credentials(ctx, "test", "not-a-device-id", secret, deviceCSR(t), testIP)
			return err
		},
	}
	for name, call := range cases {
		if err := call(); !errors.Is(err, ErrUnauthorized) {
			t.Errorf("%s: got %v, want ErrUnauthorized", name, err)
		}
	}
}

func TestCredentialsInhibited(t *testing.T) {
	ctx := context.Background()
	svc, fs, _ := newServiceFixture(t, Config{})
	hwID := randomDeviceID(t)
	secret, err := svc.Register(ctx, "test", hwID, "")
	if err != nil {
		t.Fatal(err)
	}
	id, _ := deviceid.Parse(hwID)
	fs.devices[1][id].Status = store.DeviceStatusInhibited

	if _, err := svc.Credentials(ctx, "test", hwID, secret, deviceCSR(t), testIP); !errors.Is(err, ErrInhibited) {
		t.Errorf("inhibited: got %v, want ErrInhibited", err)
	}
}

func TestCredentialsInvalidCSR(t *testing.T) {
	ctx := context.Background()
	svc, _, _ := newServiceFixture(t, Config{})
	hwID := randomDeviceID(t)
	secret, err := svc.Register(ctx, "test", hwID, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Credentials(ctx, "test", hwID, secret, "garbage", testIP); !errors.Is(err, ErrInvalidCSR) {
		t.Errorf("garbage CSR: got %v, want ErrInvalidCSR", err)
	}
}

// --- flow C ----------------------------------------------------------------

func TestInfo(t *testing.T) {
	ctx := context.Background()
	svc, fs, realm := newServiceFixture(t, Config{BrokerURL: "mqtts://broker.example.com:8883", Version: "1.2.0"})
	hwID := randomDeviceID(t)
	id, _ := deviceid.Parse(hwID)
	secret, err := svc.Register(ctx, "test", hwID, "")
	if err != nil {
		t.Fatal(err)
	}

	info, err := svc.Info(ctx, "test", hwID, secret)
	if err != nil {
		t.Fatalf("Info: %v", err)
	}
	if info.Status != "pending" {
		t.Errorf("pre-credentials status: got %q, want pending", info.Status)
	}
	if info.BrokerURL != "mqtts://broker.example.com:8883" || info.Version != "1.2.0" {
		t.Errorf("info: %+v", info)
	}
	if info.CACertPEM != realm.CACertificatePEM {
		t.Error("info must carry the realm CA certificate")
	}

	if err := fs.SetDeviceCredentials(ctx, 1, id, "1", "aa", testIP); err != nil {
		t.Fatal(err)
	}
	if info, err = svc.Info(ctx, "test", hwID, secret); err != nil || info.Status != "confirmed" {
		t.Errorf("post-credentials status: got (%v, %v), want confirmed", info, err)
	}

	fs.devices[1][id].Status = store.DeviceStatusInhibited
	if info, err = svc.Info(ctx, "test", hwID, secret); err != nil || info.Status != "inhibited" {
		t.Errorf("inhibited status: got (%v, %v), want inhibited", info, err)
	}

	if _, err := svc.Info(ctx, "test", hwID, "wrong"); !errors.Is(err, ErrUnauthorized) {
		t.Errorf("wrong secret on info: got %v, want ErrUnauthorized", err)
	}
}

func TestVerifyCredentials(t *testing.T) {
	ctx := context.Background()
	svc, fs, _ := newServiceFixture(t, Config{EnforceLatestCert: true})
	hwID := randomDeviceID(t)
	id, _ := deviceid.Parse(hwID)
	secret, err := svc.Register(ctx, "test", hwID, "")
	if err != nil {
		t.Fatal(err)
	}
	certPEM, err := svc.Credentials(ctx, "test", hwID, secret, deviceCSR(t), testIP)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Valid", func(t *testing.T) {
		res, err := svc.VerifyCredentials(ctx, "test", hwID, secret, certPEM)
		if err != nil {
			t.Fatalf("VerifyCredentials: %v", err)
		}
		if !res.Valid || res.Cause != "" {
			t.Fatalf("result: %+v, want valid", res)
		}
		cert, _ := ca.ParseCertificatePEM(certPEM)
		if !res.Until.Equal(cert.NotAfter) {
			t.Errorf("until: got %v, want %v", res.Until, cert.NotAfter)
		}
	})

	t.Run("Expired", func(t *testing.T) {
		// Time-travel the service clock past the certificate expiry.
		svcExpired, _, _ := newServiceFixture(t, Config{})
		secretE, err := svcExpired.Register(ctx, "test", hwID, "")
		if err != nil {
			t.Fatal(err)
		}
		certE, err := svcExpired.Credentials(ctx, "test", hwID, secretE, deviceCSR(t), testIP)
		if err != nil {
			t.Fatal(err)
		}
		svcExpired.now = func() time.Time { return time.Now().Add(31 * 24 * time.Hour) }
		res, err := svcExpired.VerifyCredentials(ctx, "test", hwID, secretE, certE)
		if err != nil {
			t.Fatal(err)
		}
		if res.Valid || res.Cause != CauseExpired {
			t.Errorf("result: %+v, want cause EXPIRED", res)
		}
	})

	t.Run("ForeignCA", func(t *testing.T) {
		foreign, err := ca.Generate("test", 0)
		if err != nil {
			t.Fatal(err)
		}
		foreignCert, _, _, err := foreign.SignCSR(deviceCSR(t), "test", hwID, time.Hour)
		if err != nil {
			t.Fatal(err)
		}
		res, err := svc.VerifyCredentials(ctx, "test", hwID, secret, foreignCert)
		if err != nil {
			t.Fatal(err)
		}
		if res.Valid || res.Cause != CauseInvalid {
			t.Errorf("result: %+v, want cause INVALID", res)
		}
	})

	t.Run("Garbage", func(t *testing.T) {
		res, err := svc.VerifyCredentials(ctx, "test", hwID, secret, "garbage")
		if err != nil {
			t.Fatal(err)
		}
		if res.Valid || res.Cause != CauseInvalid {
			t.Errorf("result: %+v, want cause INVALID", res)
		}
	})

	t.Run("RevokedByRotation", func(t *testing.T) {
		// Issue a second certificate: the first one's serial is no longer
		// the latest → REVOKED under enforcement.
		if _, err := svc.Credentials(ctx, "test", hwID, secret, deviceCSR(t), testIP); err != nil {
			t.Fatal(err)
		}
		res, err := svc.VerifyCredentials(ctx, "test", hwID, secret, certPEM)
		if err != nil {
			t.Fatal(err)
		}
		if res.Valid || res.Cause != CauseRevoked {
			t.Errorf("result: %+v, want cause REVOKED", res)
		}

		// Without enforcement the still-in-window certificate stays valid.
		lax := New(fs, svc.sealer, Config{EnforceLatestCert: false})
		res, err = lax.VerifyCredentials(ctx, "test", hwID, secret, certPEM)
		if err != nil {
			t.Fatal(err)
		}
		if !res.Valid {
			t.Errorf("result without enforcement: %+v, want valid", res)
		}
	})

	t.Run("Inhibited", func(t *testing.T) {
		fs.devices[1][id].Status = store.DeviceStatusInhibited
		defer func() { fs.devices[1][id].Status = store.DeviceStatusConfirmed }()
		if _, err := svc.VerifyCredentials(ctx, "test", hwID, secret, certPEM); !errors.Is(err, ErrInhibited) {
			t.Errorf("inhibited verify: got %v, want ErrInhibited", err)
		}
	})
}
