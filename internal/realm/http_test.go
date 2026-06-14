//go:build integration

package realm

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
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
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// Interface fixtures (one name, evolving versions, plus a deletable draft).
const (
	rmIface  = "com.ex.M7a.Sensors"
	rmDraft  = "com.ex.M7a.Draft"
	ifaceV1  = `{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`
	ifaceV1b = `{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":1,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"},{"endpoint":"/count","type":"integer"}]}`
	ifaceV1x = `{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":2,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"integer"}]}`
	ifaceV2  = `{"interface_name":"com.ex.M7a.Sensors","version_major":2,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"string"}]}`
	draftV0  = `{"interface_name":"com.ex.M7a.Draft","version_major":0,"version_minor":1,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/x","type":"double"}]}`
)

const triggerJSON = `{"name":"on_value","action":{"http_url":"https://example.com/hook","http_method":"post"},` +
	`"simple_triggers":[{"type":"data_trigger","on":"incoming_data","interface_name":"com.ex.M7a.Sensors","interface_major":1,"match_path":"/value","value_match_operator":"*"}]}`

type rig struct {
	st       *store.Store
	mux      *http.ServeMux
	realm    string
	realmID  int16
	rmaToken string // valid a_rma
	wrongTok string // valid token, no a_rma
	jwtKey   *rsa.PrivateKey
	otherPub string
}

func newRig(t *testing.T) *rig {
	t.Helper()
	ctx := context.Background()
	pool := testutil.StartTimescale(t)
	st, err := store.New(ctx, pool.Config().ConnString())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(st.Close)

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	other, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	var suffix [4]byte
	_, _ = rand.Read(suffix[:])
	realmName := "rm" + hex.EncodeToString(suffix[:])
	realm, err := st.CreateRealm(ctx, store.NewRealm{
		Name:               realmName,
		JWTPublicKeysPEM:   []string{pubPEM(t, &key.PublicKey)},
		CACertificatePEM:   "test-ca",
		CAPrivateKeySealed: []byte("sealed"),
	})
	if err != nil {
		t.Fatalf("CreateRealm: %v", err)
	}

	mux := http.NewServeMux()
	NewAPI(NewService(st, nil, discardLogger()), auth.NewMiddleware(st)).Mount(mux)

	return &rig{
		st: st, mux: mux, realm: realmName, realmID: realm.ID, jwtKey: key,
		otherPub: pubPEM(t, &other.PublicKey),
		rmaToken: mintToken(t, key, jwt.MapClaims{"a_rma": []string{".*::.*"}}),
		wrongTok: mintToken(t, key, jwt.MapClaims{"a_aea": []string{".*::.*"}}),
	}
}

// req drives one authenticated request; rawBody (if non-empty) is wrapped in
// the {"data": ...} envelope.
func (r *rig) req(t *testing.T, method, path, rawBody, token string) *httptest.ResponseRecorder {
	t.Helper()
	var body io.Reader
	if rawBody != "" {
		body = strings.NewReader(`{"data":` + rawBody + `}`)
	}
	httpReq := httptest.NewRequest(method, "/realmmanagement/v1/"+r.realm+path, body)
	if token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	r.mux.ServeHTTP(rec, httpReq)
	return rec
}

func TestRealmManagement(t *testing.T) {
	r := newRig(t)

	t.Run("Auth401WithoutToken", func(t *testing.T) {
		if rec := r.req(t, http.MethodGet, "/interfaces", "", ""); rec.Code != http.StatusUnauthorized {
			t.Errorf("no token: got %d, want 401", rec.Code)
		}
	})
	t.Run("Auth403WrongClaim", func(t *testing.T) {
		if rec := r.req(t, http.MethodGet, "/interfaces", "", r.wrongTok); rec.Code != http.StatusForbidden {
			t.Errorf("wrong claim: got %d, want 403", rec.Code)
		}
	})

	t.Run("InstallInterface", func(t *testing.T) {
		if rec := r.req(t, http.MethodPost, "/interfaces", ifaceV1, r.rmaToken); rec.Code != http.StatusCreated {
			t.Fatalf("install: got %d, want 201 (%s)", rec.Code, rec.Body)
		}
		// Duplicate → 409.
		if rec := r.req(t, http.MethodPost, "/interfaces", ifaceV1, r.rmaToken); rec.Code != http.StatusConflict {
			t.Errorf("duplicate install: got %d, want 409", rec.Code)
		}
	})

	t.Run("ListAndGet", func(t *testing.T) {
		var got []string
		decodeData(t, r.req(t, http.MethodGet, "/interfaces", "", r.rmaToken), &got)
		if !contains(got, rmIface) {
			t.Errorf("list interfaces = %v, want to contain %s", got, rmIface)
		}
		var majors []int
		decodeData(t, r.req(t, http.MethodGet, "/interfaces/"+rmIface, "", r.rmaToken), &majors)
		if len(majors) != 1 || majors[0] != 1 {
			t.Errorf("majors = %v, want [1]", majors)
		}
		if rec := r.req(t, http.MethodGet, "/interfaces/"+rmIface+"/1", "", r.rmaToken); rec.Code != http.StatusOK {
			t.Errorf("get interface: got %d, want 200", rec.Code)
		}
	})

	t.Run("MinorUpgradeAccepted", func(t *testing.T) {
		if rec := r.req(t, http.MethodPut, "/interfaces/"+rmIface+"/1", ifaceV1b, r.rmaToken); rec.Code != http.StatusNoContent {
			t.Fatalf("additive minor upgrade: got %d, want 204 (%s)", rec.Code, rec.Body)
		}
	})
	t.Run("MappingMutationRejected", func(t *testing.T) {
		// Changing /value's type is not an additive upgrade (CheckMinorUpgrade).
		if rec := r.req(t, http.MethodPut, "/interfaces/"+rmIface+"/1", ifaceV1x, r.rmaToken); rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("mapping mutation: got %d, want 422 (%s)", rec.Code, rec.Body)
		}
	})

	t.Run("MajorCoexistence", func(t *testing.T) {
		if rec := r.req(t, http.MethodPost, "/interfaces", ifaceV2, r.rmaToken); rec.Code != http.StatusCreated {
			t.Fatalf("install major 2: got %d, want 201 (%s)", rec.Code, rec.Body)
		}
		var majors []int
		decodeData(t, r.req(t, http.MethodGet, "/interfaces/"+rmIface, "", r.rmaToken), &majors)
		if len(majors) != 2 || majors[0] != 1 || majors[1] != 2 {
			t.Errorf("majors after coexistence = %v, want [1 2]", majors)
		}
	})

	t.Run("DeleteRules", func(t *testing.T) {
		ctx := context.Background()
		// Major != 0 can't be deleted.
		if rec := r.req(t, http.MethodDelete, "/interfaces/"+rmIface+"/1", "", r.rmaToken); rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("delete major 1: got %d, want 422", rec.Code)
		}
		// Install a draft (major 0), reference it in an introspection → can't delete.
		if rec := r.req(t, http.MethodPost, "/interfaces", draftV0, r.rmaToken); rec.Code != http.StatusCreated {
			t.Fatalf("install draft: got %d (%s)", rec.Code, rec.Body)
		}
		dev, _ := deviceid.Random()
		if err := r.st.RegisterDevice(ctx, r.realmID, dev, "h"); err != nil {
			t.Fatal(err)
		}
		if _, err := r.st.UpdateIntrospection(ctx, r.realmID, dev,
			map[string]store.InterfaceVersion{rmDraft: {Major: 0, Minor: 1}}); err != nil {
			t.Fatal(err)
		}
		if rec := r.req(t, http.MethodDelete, "/interfaces/"+rmDraft+"/0", "", r.rmaToken); rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("delete introspected draft: got %d, want 422", rec.Code)
		}
		// Drop it from the introspection → now deletable.
		if _, err := r.st.UpdateIntrospection(ctx, r.realmID, dev, map[string]store.InterfaceVersion{}); err != nil {
			t.Fatal(err)
		}
		if rec := r.req(t, http.MethodDelete, "/interfaces/"+rmDraft+"/0", "", r.rmaToken); rec.Code != http.StatusNoContent {
			t.Errorf("delete unused draft: got %d, want 204", rec.Code)
		}
	})

	t.Run("Triggers", func(t *testing.T) {
		if rec := r.req(t, http.MethodPost, "/triggers", triggerJSON, r.rmaToken); rec.Code != http.StatusCreated {
			t.Fatalf("create trigger: got %d, want 201 (%s)", rec.Code, rec.Body)
		}
		if rec := r.req(t, http.MethodPost, "/triggers", triggerJSON, r.rmaToken); rec.Code != http.StatusConflict {
			t.Errorf("duplicate trigger: got %d, want 409", rec.Code)
		}
		if rec := r.req(t, http.MethodPost, "/triggers", `{"name":"bad","action":{},"simple_triggers":[]}`, r.rmaToken); rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("invalid trigger: got %d, want 422", rec.Code)
		}
		var names []string
		decodeData(t, r.req(t, http.MethodGet, "/triggers", "", r.rmaToken), &names)
		if !contains(names, "on_value") {
			t.Errorf("trigger list = %v, want to contain on_value", names)
		}
		if rec := r.req(t, http.MethodGet, "/triggers/on_value", "", r.rmaToken); rec.Code != http.StatusOK {
			t.Errorf("get trigger: got %d, want 200", rec.Code)
		}
		if rec := r.req(t, http.MethodDelete, "/triggers/on_value", "", r.rmaToken); rec.Code != http.StatusNoContent {
			t.Errorf("delete trigger: got %d, want 204", rec.Code)
		}
		if rec := r.req(t, http.MethodGet, "/triggers/on_value", "", r.rmaToken); rec.Code != http.StatusNotFound {
			t.Errorf("get deleted trigger: got %d, want 404", rec.Code)
		}
	})

	t.Run("ConfigAuth", func(t *testing.T) {
		// Rotate to a 2-key set that still includes the original, so the test
		// token keeps verifying.
		rotated := pubPEM(t, &r.jwtKey.PublicKey) + "\n" + r.otherPub
		if rec := r.req(t, http.MethodPut, "/config/auth", `{"jwt_public_key_pem":`+jsonStr(rotated)+`}`, r.rmaToken); rec.Code != http.StatusNoContent {
			t.Fatalf("put config/auth: got %d, want 204 (%s)", rec.Code, rec.Body)
		}
		var cfg authConfig
		decodeData(t, r.req(t, http.MethodGet, "/config/auth", "", r.rmaToken), &cfg)
		if cfg.JWTPublicKeyPEM != rotated {
			t.Errorf("config/auth key not rotated")
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
