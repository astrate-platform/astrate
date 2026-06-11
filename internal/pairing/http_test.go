//go:build integration

package pairing

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/pairing/ca"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// httpFixture is the full wired stack: real store, real CA, real JWT
// middleware, httptest server.
type httpFixture struct {
	st        *store.Store
	sealer    *store.KeySealer
	realm     *store.Realm
	realmName string
	server    *httptest.Server
	jwtKey    *rsa.PrivateKey
	svc       *Service
}

// newHTTPFixture boots the stack against the shared TimescaleDB. cfg/apiCfg
// zero-values select production defaults.
func newHTTPFixture(t *testing.T, cfg Config, apiCfg APIConfig) *httpFixture {
	t.Helper()
	ctx := context.Background()

	pool := testutil.StartTimescale(t)
	st, err := store.New(ctx, pool.Config().ConnString())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(st.Close)

	key := make([]byte, store.MasterKeySize)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	sealer, err := store.NewKeySealer(key)
	if err != nil {
		t.Fatal(err)
	}

	jwtKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	pubDER, err := x509.MarshalPKIXPublicKey(&jwtKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pubPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))

	realmName := "p" + strconv.FormatInt(time.Now().UnixNano(), 36)
	caPEM, sealedKey, err := ProvisionCA(realmName, 0, sealer)
	if err != nil {
		t.Fatalf("ProvisionCA: %v", err)
	}
	realm, err := st.CreateRealm(ctx, store.NewRealm{
		Name:               realmName,
		JWTPublicKeysPEM:   []string{pubPEM},
		CACertificatePEM:   caPEM,
		CAPrivateKeySealed: sealedKey,
	})
	if err != nil {
		t.Fatalf("CreateRealm: %v", err)
	}

	if cfg.BrokerURL == "" {
		cfg.BrokerURL = "mqtts://broker.test.example:8883"
	}
	if cfg.Version == "" {
		cfg.Version = "0.1.0-test"
	}
	svc := New(st, sealer, cfg)
	api := NewAPI(svc, auth.NewMiddleware(st), apiCfg)
	mux := http.NewServeMux()
	api.Mount(mux)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return &httpFixture{
		st: st, sealer: sealer, realm: realm, realmName: realmName,
		server: server, jwtKey: jwtKey, svc: svc,
	}
}

// agentJWT signs a realm JWT with the given a_pa authorization strings (nil
// means catch-all).
func (f *httpFixture) agentJWT(t *testing.T, claim auth.Claim, grants []string) string {
	t.Helper()
	if grants == nil {
		grants = []string{".*::.*"}
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		string(claim): grants,
		"exp":         time.Now().Add(time.Hour).Unix(),
	}).SignedString(f.jwtKey)
	if err != nil {
		t.Fatal(err)
	}
	return tok
}

// request performs an HTTP call and returns status + body.
func (f *httpFixture) request(t *testing.T, method, path, bearer string, body any) (int, []byte) {
	t.Helper()
	var reader io.Reader
	if body != nil {
		enc, err := json.Marshal(map[string]any{"data": body})
		if err != nil {
			t.Fatal(err)
		}
		reader = bytes.NewReader(enc)
	}
	req, err := http.NewRequest(method, f.server.URL+path, reader)
	if err != nil {
		t.Fatal(err)
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := f.server.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, got
}

// normalize replaces dynamic body fields (secrets, PEMs, timestamps) with
// stable placeholders so bodies golden-compare byte-for-byte.
var normalizers = []struct {
	re   *regexp.Regexp
	repl string
}{
	{regexp.MustCompile(`"credentials_secret":"[^"]*"`), `"credentials_secret":"<SECRET>"`},
	{regexp.MustCompile(`"client_crt":"[^"]*"`), `"client_crt":"<PEM>"`},
	{regexp.MustCompile(`"ca_crt":"[^"]*"`), `"ca_crt":"<PEM>"`},
	{regexp.MustCompile(`"timestamp":"[^"]*"`), `"timestamp":"<TS>"`},
	{regexp.MustCompile(`"until":"[^"]*"`), `"until":"<TS>"`},
}

func normalize(body []byte) []byte {
	for _, n := range normalizers {
		body = n.re.ReplaceAll(body, []byte(n.repl))
	}
	return body
}

// dataField unwraps one string field from a {"data": {...}} body.
func dataField(t *testing.T, body []byte, field string) string {
	t.Helper()
	var env struct {
		Data map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		t.Fatalf("decoding %s: %v (body %s)", field, err, body)
	}
	var s string
	if err := json.Unmarshal(env.Data[field], &s); err != nil {
		t.Fatalf("decoding %s: %v (body %s)", field, err, body)
	}
	return s
}

const elixirDateTimeLayout = "2006-01-02 15:04:05.000Z"

// TestPairingHTTP is the umbrella T2 suite (docs/ROADMAP.md §5): golden
// HTTP flows A–C against a live store. Sub-tests share fixture state and
// run in order.
func TestPairingHTTP(t *testing.T) {
	f := newHTTPFixture(t, Config{EnforceLatestCert: true}, APIConfig{})
	agentToken := f.agentJWT(t, auth.ClaimPairing, nil)
	base := "/pairing/v1/" + f.realmName

	hwID := mustRandomDeviceID(t)
	var (
		secret    string
		clientCrt string
	)

	// --- flow A ---------------------------------------------------------

	t.Run("FlowA_Register201", func(t *testing.T) {
		status, body := f.request(t, "POST", base+"/agent/devices", agentToken,
			map[string]string{"hw_id": hwID})
		if status != http.StatusCreated {
			t.Fatalf("status: got %d, want 201 (body %s)", status, body)
		}
		secret = dataField(t, body, "credentials_secret")
		if len(secret) != 44 {
			t.Errorf("secret length: got %d, want 44", len(secret))
		}
		testutil.Golden(t, "http/register_201.json", normalize(body))
	})

	t.Run("FlowA_RotateBeforeCredentials", func(t *testing.T) {
		status, body := f.request(t, "POST", base+"/agent/devices", agentToken,
			map[string]string{"hw_id": hwID})
		if status != http.StatusCreated {
			t.Fatalf("status: got %d, want 201 (body %s)", status, body)
		}
		rotated := dataField(t, body, "credentials_secret")
		if rotated == secret {
			t.Error("pre-credentials re-registration must rotate the secret")
		}
		secret = rotated
	})

	t.Run("FlowA_BadHwID422", func(t *testing.T) {
		for _, bad := range []string{
			"h4-Dx_RYTU-RbpDOTabhR",   // 21 characters
			"h4-Dx_RYTU-RbpDOTabhRg=", // padded
			"h4+Dx/RYTU+RbpDOTabhRg",  // non-url alphabet
		} {
			status, body := f.request(t, "POST", base+"/agent/devices", agentToken,
				map[string]string{"hw_id": bad})
			if status != http.StatusUnprocessableEntity {
				t.Fatalf("hw_id %q: got %d, want 422 (body %s)", bad, status, body)
			}
			testutil.Golden(t, "http/register_422_bad_hwid.json", body)
		}

		status, body := f.request(t, "POST", base+"/agent/devices", agentToken,
			map[string]string{})
		if status != http.StatusUnprocessableEntity {
			t.Fatalf("blank hw_id: got %d, want 422 (body %s)", status, body)
		}
		testutil.Golden(t, "http/register_422_blank_hwid.json", body)
	})

	t.Run("FlowA_AuthRequired", func(t *testing.T) {
		status, body := f.request(t, "POST", base+"/agent/devices", "",
			map[string]string{"hw_id": hwID})
		if status != http.StatusUnauthorized {
			t.Fatalf("no token: got %d, want 401", status)
		}
		testutil.Golden(t, "http/envelope_401.json", body)

		wrongClaim := f.agentJWT(t, auth.ClaimAppEngine, nil)
		status, body = f.request(t, "POST", base+"/agent/devices", wrongClaim,
			map[string]string{"hw_id": hwID})
		if status != http.StatusForbidden {
			t.Fatalf("wrong claim: got %d, want 403", status)
		}
		testutil.Golden(t, "http/envelope_403.json", body)
	})

	// --- flow B ---------------------------------------------------------

	t.Run("FlowB_Credentials201", func(t *testing.T) {
		csr := deviceCSR(t)
		status, body := f.request(t, "POST",
			base+"/devices/"+hwID+"/protocols/astarte_mqtt_v1/credentials",
			secret, map[string]string{"csr": csr})
		if status != http.StatusCreated {
			t.Fatalf("status: got %d, want 201 (body %s)", status, body)
		}
		testutil.Golden(t, "http/credentials_201.json", normalize(body))

		clientCrt = dataField(t, body, "client_crt")
		cert, err := ca.ParseCertificatePEM(clientCrt)
		if err != nil {
			t.Fatalf("client_crt does not parse: %v", err)
		}
		if got, want := cert.Subject.CommonName, f.realmName+"/"+hwID; got != want {
			t.Errorf("CN: got %q, want %q", got, want)
		}

		// Chain check 1: crypto/x509 against the realm CA from flow C.
		roots := x509.NewCertPool()
		if !roots.AppendCertsFromPEM([]byte(f.realm.CACertificatePEM)) {
			t.Fatal("adding realm CA to pool")
		}
		if _, err := cert.Verify(x509.VerifyOptions{
			Roots:     roots,
			KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		}); err != nil {
			t.Errorf("client_crt does not chain to the realm CA: %v", err)
		}

		// Chain check 2: openssl verify exec smoke.
		opensslVerify(t, f.realm.CACertificatePEM, clientCrt)
	})

	t.Run("FlowA_RegisterAfterCredentials422", func(t *testing.T) {
		status, body := f.request(t, "POST", base+"/agent/devices", agentToken,
			map[string]string{"hw_id": hwID})
		if status != http.StatusUnprocessableEntity {
			t.Fatalf("status: got %d, want 422 (body %s)", status, body)
		}
		testutil.Golden(t, "http/register_422_already.json", body)
	})

	t.Run("FlowB_UniformUnauthorized", func(t *testing.T) {
		csr := deviceCSR(t)

		startWrong := time.Now()
		status1, body1 := f.request(t, "POST",
			base+"/devices/"+hwID+"/protocols/astarte_mqtt_v1/credentials",
			"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", map[string]string{"csr": csr})
		wrongDur := time.Since(startWrong)

		unknown := mustRandomDeviceID(t)
		startUnknown := time.Now()
		status2, body2 := f.request(t, "POST",
			base+"/devices/"+unknown+"/protocols/astarte_mqtt_v1/credentials",
			"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", map[string]string{"csr": csr})
		unknownDur := time.Since(startUnknown)

		if status1 != http.StatusUnauthorized || status2 != http.StatusUnauthorized {
			t.Fatalf("statuses: got %d/%d, want 401/401", status1, status2)
		}
		if !bytes.Equal(body1, body2) {
			t.Errorf("uniform error bodies differ: %s vs %s", body1, body2)
		}
		testutil.Golden(t, "http/envelope_401.json", body1)

		// Same timing class: both paths burn one bcrypt comparison
		// (cost 10 ≈ tens of ms; 5ms is a generous lower bound that
		// catches an accidental early return).
		const minWork = 5 * time.Millisecond
		if wrongDur < minWork || unknownDur < minWork {
			t.Errorf("auth failure too fast (wrong=%v unknown=%v): bcrypt work missing", wrongDur, unknownDur)
		}
	})

	t.Run("FlowB_Inhibited403", func(t *testing.T) {
		id, err := deviceid.Parse(hwID)
		if err != nil {
			t.Fatal(err)
		}
		if err := f.st.SetDeviceInhibited(context.Background(), f.realm.ID, id, true); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			_ = f.st.SetDeviceInhibited(context.Background(), f.realm.ID, id, false)
		})

		status, body := f.request(t, "POST",
			base+"/devices/"+hwID+"/protocols/astarte_mqtt_v1/credentials",
			secret, map[string]string{"csr": deviceCSR(t)})
		if status != http.StatusForbidden {
			t.Fatalf("status: got %d, want 403 (body %s)", status, body)
		}
		testutil.Golden(t, "http/envelope_403.json", body)
	})

	// --- flow C ---------------------------------------------------------

	t.Run("FlowC_Info200", func(t *testing.T) {
		status, body := f.request(t, "GET", base+"/devices/"+hwID, secret, nil)
		if status != http.StatusOK {
			t.Fatalf("status: got %d, want 200 (body %s)", status, body)
		}
		var env struct {
			Data struct {
				Protocols struct {
					AstarteMQTTV1 struct {
						BrokerURL string `json:"broker_url"`
						CACrt     string `json:"ca_crt"`
					} `json:"astarte_mqtt_v1"`
				} `json:"protocols"`
				Status  string `json:"status"`
				Version string `json:"version"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &env); err != nil {
			t.Fatalf("decoding info: %v", err)
		}
		if env.Data.Status != "confirmed" {
			t.Errorf("status: got %q, want confirmed", env.Data.Status)
		}
		if env.Data.Protocols.AstarteMQTTV1.BrokerURL != "mqtts://broker.test.example:8883" {
			t.Errorf("broker_url: got %q", env.Data.Protocols.AstarteMQTTV1.BrokerURL)
		}
		if env.Data.Protocols.AstarteMQTTV1.CACrt != f.realm.CACertificatePEM {
			t.Error("ca_crt does not match the realm CA")
		}
		testutil.Golden(t, "http/info_200.json", normalize(body))
	})

	t.Run("FlowC_VerifyValid", func(t *testing.T) {
		status, body := f.request(t, "POST",
			base+"/devices/"+hwID+"/protocols/astarte_mqtt_v1/credentials/verify",
			secret, map[string]string{"client_crt": clientCrt})
		if status != http.StatusOK {
			t.Fatalf("status: got %d, want 200 (body %s)", status, body)
		}
		var env struct {
			Data struct {
				Valid     bool   `json:"valid"`
				Timestamp string `json:"timestamp"`
				Until     string `json:"until"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &env); err != nil {
			t.Fatal(err)
		}
		if !env.Data.Valid {
			t.Fatalf("verify: got %s, want valid", body)
		}
		until, err := time.Parse(elixirDateTimeLayout, env.Data.Until)
		if err != nil {
			t.Fatalf("until %q does not parse with the upstream layout: %v", env.Data.Until, err)
		}
		cert, _ := ca.ParseCertificatePEM(clientCrt)
		if until.Unix() != cert.NotAfter.Unix() {
			t.Errorf("until: got %v, want certificate NotAfter %v", until, cert.NotAfter)
		}
		if _, err := time.Parse(elixirDateTimeLayout, env.Data.Timestamp); err != nil {
			t.Errorf("timestamp %q does not parse with the upstream layout: %v", env.Data.Timestamp, err)
		}
		testutil.Golden(t, "http/verify_200_valid.json", normalize(body))
	})

	t.Run("FlowC_VerifyExpired", func(t *testing.T) {
		// Issue through a second service sharing the store but with a 1s
		// TTL, wait out the certificate, then verify it on the main API.
		shortSvc := New(f.st, f.sealer, Config{CertTTL: time.Second, EnforceLatestCert: true})
		shortCrt, err := shortSvc.Credentials(context.Background(), f.realmName, hwID, secret,
			deviceCSR(t), testIP)
		if err != nil {
			t.Fatalf("issuing short-TTL certificate: %v", err)
		}
		time.Sleep(1500 * time.Millisecond)

		status, body := f.request(t, "POST",
			base+"/devices/"+hwID+"/protocols/astarte_mqtt_v1/credentials/verify",
			secret, map[string]string{"client_crt": shortCrt})
		if status != http.StatusOK {
			t.Fatalf("status: got %d, want 200 (body %s)", status, body)
		}
		testutil.Golden(t, "http/verify_200_expired.json", normalize(body))
	})

	t.Run("FlowC_VerifyForeignCA", func(t *testing.T) {
		foreign, err := ca.Generate(f.realmName, 0)
		if err != nil {
			t.Fatal(err)
		}
		foreignCrt, _, _, err := foreign.SignCSR(deviceCSR(t), f.realmName, hwID, time.Hour)
		if err != nil {
			t.Fatal(err)
		}
		status, body := f.request(t, "POST",
			base+"/devices/"+hwID+"/protocols/astarte_mqtt_v1/credentials/verify",
			secret, map[string]string{"client_crt": foreignCrt})
		if status != http.StatusOK {
			t.Fatalf("status: got %d, want 200 (body %s)", status, body)
		}
		testutil.Golden(t, "http/verify_200_invalid.json", normalize(body))
	})

	t.Run("FlowC_VerifyRevokedAfterRotation", func(t *testing.T) {
		// The short-TTL issuance in VerifyExpired rotated the recorded
		// serial, so the original (still in-window) certificate is now
		// superseded → REVOKED under enforcement.
		status, body := f.request(t, "POST",
			base+"/devices/"+hwID+"/protocols/astarte_mqtt_v1/credentials/verify",
			secret, map[string]string{"client_crt": clientCrt})
		if status != http.StatusOK {
			t.Fatalf("status: got %d, want 200 (body %s)", status, body)
		}
		testutil.Golden(t, "http/verify_200_revoked.json", normalize(body))
	})

	// --- unregister / re-register ----------------------------------------

	t.Run("FlowA_Unregister", func(t *testing.T) {
		status, body := f.request(t, "DELETE", base+"/agent/devices/"+hwID, agentToken, nil)
		if status != http.StatusNoContent {
			t.Fatalf("status: got %d, want 204 (body %s)", status, body)
		}
		if len(body) != 0 {
			t.Errorf("204 body must be empty, got %q", body)
		}

		// Old secret is dead.
		status, _ = f.request(t, "POST",
			base+"/devices/"+hwID+"/protocols/astarte_mqtt_v1/credentials",
			secret, map[string]string{"csr": deviceCSR(t)})
		if status != http.StatusUnauthorized {
			t.Fatalf("old secret after unregister: got %d, want 401", status)
		}

		// The device is registrable again (its row and data survived).
		status, body = f.request(t, "POST", base+"/agent/devices", agentToken,
			map[string]string{"hw_id": hwID})
		if status != http.StatusCreated {
			t.Fatalf("re-register after unregister: got %d (body %s)", status, body)
		}

		// Unknown device → upstream 404 shape.
		status, body = f.request(t, "DELETE", base+"/agent/devices/"+mustRandomDeviceID(t), agentToken, nil)
		if status != http.StatusNotFound {
			t.Fatalf("unknown unregister: got %d, want 404 (body %s)", status, body)
		}
		testutil.Golden(t, "http/device_404.json", body)
	})

	t.Run("RegistrationLimit422", func(t *testing.T) {
		limit := int32(0)
		limited, err := f.st.CreateRealm(context.Background(), store.NewRealm{
			Name:                    f.realmName + "l",
			JWTPublicKeysPEM:        f.realm.JWTPublicKeysPEM,
			CACertificatePEM:        f.realm.CACertificatePEM,
			CAPrivateKeySealed:      f.realm.CAPrivateKeySealed,
			DeviceRegistrationLimit: &limit,
		})
		if err != nil {
			t.Fatal(err)
		}
		status, body := f.request(t, "POST", "/pairing/v1/"+limited.Name+"/agent/devices",
			agentToken, map[string]string{"hw_id": mustRandomDeviceID(t)})
		if status != http.StatusUnprocessableEntity {
			t.Fatalf("status: got %d, want 422 (body %s)", status, body)
		}
		testutil.Golden(t, "http/register_422_limit.json", body)
	})
}

// TestPairingHTTPRateLimit exercises 429 behaviour on a dedicated fixture
// with a tiny credentials budget.
func TestPairingHTTPRateLimit(t *testing.T) {
	f := newHTTPFixture(t, Config{}, APIConfig{
		CredentialsRate:  0.0001, // effectively no refill within the test
		CredentialsBurst: 2,
	})
	base := "/pairing/v1/" + f.realmName
	path := base + "/devices/" + mustRandomDeviceID(t) + "/protocols/astarte_mqtt_v1/credentials"

	for i := 0; i < 2; i++ {
		status, _ := f.request(t, "POST", path, "irrelevant", map[string]string{"csr": "x"})
		if status == http.StatusTooManyRequests {
			t.Fatalf("request %d within burst must not be rate limited", i)
		}
	}
	status, body := f.request(t, "POST", path, "irrelevant", map[string]string{"csr": "x"})
	if status != http.StatusTooManyRequests {
		t.Fatalf("status: got %d, want 429 (body %s)", status, body)
	}
	testutil.Golden(t, "http/envelope_429.json", body)
}

// mustRandomDeviceID returns a fresh random device ID wire string.
func mustRandomDeviceID(t *testing.T) string {
	t.Helper()
	id, err := deviceid.Random()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

// opensslVerify shells out to `openssl verify` as an independent
// implementation cross-check of the issued chain.
func opensslVerify(t *testing.T, caPEM, certPEM string) {
	t.Helper()
	openssl, err := exec.LookPath("openssl")
	if err != nil {
		t.Skip("openssl not found in PATH; chain already verified with crypto/x509")
	}

	dir := t.TempDir()
	caPath := filepath.Join(dir, "ca.pem")
	certPath := filepath.Join(dir, "client.pem")
	if err := os.WriteFile(caPath, []byte(caPEM), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(certPath, []byte(certPEM), 0o600); err != nil {
		t.Fatal(err)
	}

	out, err := exec.Command(openssl, "verify", "-CAfile", caPath, certPath).CombinedOutput()
	if err != nil {
		t.Fatalf("openssl verify failed: %v\n%s", err, out)
	}
	if want := fmt.Sprintf("%s: OK", certPath); !bytes.Contains(out, []byte(want)) {
		t.Fatalf("openssl verify output: %s", out)
	}
}
