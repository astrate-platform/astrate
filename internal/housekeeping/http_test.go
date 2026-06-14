//go:build integration

package housekeeping

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/pairing"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

type hkRig struct {
	st       *store.Store
	sealer   *store.KeySealer
	mux      *http.ServeMux
	haToken  string // valid a_ha
	wrongTok string // valid token, no a_ha
	realmKey string // a JWT public key PEM realms are created with
}

func newHKRig(t *testing.T) *hkRig {
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

	instanceKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	realmKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	instancePub := pubPEM(t, &instanceKey.PublicKey)

	mux := http.NewServeMux()
	NewAPI(NewService(st, sealer, nil, discardLogger()), auth.NewMiddleware(st), []string{instancePub}).Mount(mux)

	return &hkRig{
		st: st, sealer: sealer, mux: mux,
		realmKey: pubPEM(t, &realmKey.PublicKey),
		haToken:  mintToken(t, instanceKey, jwt.MapClaims{"a_ha": []string{".*::.*"}}),
		wrongTok: mintToken(t, instanceKey, jwt.MapClaims{"a_aea": []string{".*::.*"}}),
	}
}

func (r *hkRig) req(t *testing.T, method, path, rawBody, token string) *httptest.ResponseRecorder {
	t.Helper()
	var body io.Reader
	if rawBody != "" {
		body = strings.NewReader(`{"data":` + rawBody + `}`)
	}
	httpReq := httptest.NewRequest(method, "/housekeeping/v1"+path, body)
	if token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	r.mux.ServeHTTP(rec, httpReq)
	return rec
}

func TestHousekeeping(t *testing.T) {
	r := newHKRig(t)
	realmName := "hk" + randSuffix(t)
	createBody := `{"realm_name":` + jsonStr(realmName) + `,"jwt_public_key_pem":` + jsonStr(r.realmKey) + `}`

	t.Run("Auth401WithoutToken", func(t *testing.T) {
		if rec := r.req(t, http.MethodGet, "/realms", "", ""); rec.Code != http.StatusUnauthorized {
			t.Errorf("no token: got %d, want 401", rec.Code)
		}
	})
	t.Run("Auth403WrongClaim", func(t *testing.T) {
		if rec := r.req(t, http.MethodGet, "/realms", "", r.wrongTok); rec.Code != http.StatusForbidden {
			t.Errorf("wrong claim: got %d, want 403", rec.Code)
		}
	})

	t.Run("CreateRealm", func(t *testing.T) {
		rec := r.req(t, http.MethodPost, "/realms", createBody, r.haToken)
		if rec.Code != http.StatusCreated {
			t.Fatalf("create: got %d, want 201 (%s)", rec.Code, rec.Body)
		}
		var body realmBody
		decodeData(t, rec, &body)
		if body.RealmName != realmName || body.JWTPublicKeyPEM != r.realmKey {
			t.Errorf("create body = %+v", body)
		}
		// Duplicate → 409.
		if rec := r.req(t, http.MethodPost, "/realms", createBody, r.haToken); rec.Code != http.StatusConflict {
			t.Errorf("duplicate create: got %d, want 409", rec.Code)
		}
	})

	t.Run("MissingJWTKey422", func(t *testing.T) {
		body := `{"realm_name":` + jsonStr("hk"+randSuffix(t)) + `}`
		if rec := r.req(t, http.MethodPost, "/realms", body, r.haToken); rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("missing jwt key: got %d, want 422 (%s)", rec.Code, rec.Body)
		}
	})

	t.Run("GetAndList", func(t *testing.T) {
		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusOK {
			t.Errorf("get realm: got %d, want 200", rec.Code)
		}
		if rec := r.req(t, http.MethodGet, "/realms/nope"+randSuffix(t), "", r.haToken); rec.Code != http.StatusNotFound {
			t.Errorf("get unknown realm: got %d, want 404", rec.Code)
		}
		var names []string
		decodeData(t, r.req(t, http.MethodGet, "/realms", "", r.haToken), &names)
		if !contains(names, realmName) {
			t.Errorf("realm list = %v, want to contain %s", names, realmName)
		}
	})

	// The realm housekeeping just created must immediately serve pairing: its
	// CA exists and is usable, so a device registers (docs/ROADMAP.md §8.3
	// cross-domain check).
	t.Run("CreatedRealmServesPairing", func(t *testing.T) {
		svc := pairing.New(r.st, r.sealer, pairing.Config{BrokerURL: "mqtts://localhost:8883"})
		dev, _ := deviceid.Random()
		secret, err := svc.Register(context.Background(), realmName, dev.String(), "")
		if err != nil {
			t.Fatalf("pairing register against housekeeping-created realm: %v", err)
		}
		if len(secret) != 44 {
			t.Errorf("credentials secret length = %d, want 44", len(secret))
		}
	})

	t.Run("DeleteRealm", func(t *testing.T) {
		if rec := r.req(t, http.MethodDelete, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusNoContent {
			t.Fatalf("delete realm: got %d, want 204", rec.Code)
		}
		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusNotFound {
			t.Errorf("get deleted realm: got %d, want 404", rec.Code)
		}
	})
}

// --- helpers ----------------------------------------------------------------

func discardLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func mintToken(t *testing.T, key *rsa.PrivateKey, claims jwt.MapClaims) string {
	t.Helper()
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	s, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

func pubPEM(t *testing.T, pub *rsa.PublicKey) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
}

func decodeData(t *testing.T, rec *httptest.ResponseRecorder, dst any) {
	t.Helper()
	if rec.Code/100 != 2 {
		t.Fatalf("non-2xx response %d: %s", rec.Code, rec.Body)
	}
	var env struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode envelope: %v (%s)", err, rec.Body)
	}
	if err := json.Unmarshal(env.Data, dst); err != nil {
		t.Fatalf("decode data: %v (%s)", err, env.Data)
	}
}

func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

func jsonStr(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func randSuffix(t *testing.T) string {
	t.Helper()
	var b [4]byte
	_, _ = rand.Read(b[:])
	return string([]byte{
		'a' + b[0]%26, 'a' + b[1]%26, 'a' + b[2]%26, 'a' + b[3]%26,
	})
}
