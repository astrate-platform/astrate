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
	"net/netip"
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

// hkRigOpts selects the optional service knobs a rig wires.
type hkRigOpts struct {
	defaultRetention *int64
	deletionDisabled bool
}

func newHKRig(t *testing.T) *hkRig {
	t.Helper()
	return newHKRigOpts(t, hkRigOpts{})
}

func newHKRigWithDefault(t *testing.T, defaultRetention *int64) *hkRig {
	t.Helper()
	return newHKRigOpts(t, hkRigOpts{defaultRetention: defaultRetention})
}

func newHKRigWithDeletionDisabled(t *testing.T) *hkRig {
	t.Helper()
	return newHKRigOpts(t, hkRigOpts{deletionDisabled: true})
}

func newHKRigOpts(t *testing.T, opts hkRigOpts) *hkRig {
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
	NewAPI(NewService(st, sealer, nil, discardLogger()).
		WithDefaultDatastreamMaximumStorageRetention(opts.defaultRetention).
		WithRealmDeletionDisabled(opts.deletionDisabled),
		auth.NewMiddleware(st), []string{instancePub}).Mount(mux)

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

// TestHousekeepingPatchRealm drives #74 against the live wire shapes measured
// on upstream Astarte v1.2.0 (see .trickle-phase.md): every happy path answers
// 200 with an EMPTY body and its effect is re-read via GET.
func TestHousekeepingPatchRealm(t *testing.T) {
	r := newHKRig(t)
	realmName := "pk" + randSuffix(t)
	createBody := `{"realm_name":` + jsonStr(realmName) + `,"jwt_public_key_pem":` + jsonStr(r.realmKey) + `}`
	if rec := r.req(t, http.MethodPost, "/realms", createBody, r.haToken); rec.Code != http.StatusCreated {
		t.Fatalf("create: got %d, want 201 (%s)", rec.Code, rec.Body)
	}

	patch := func(t *testing.T, rawBody string) *httptest.ResponseRecorder {
		t.Helper()
		return r.req(t, http.MethodPatch, "/realms/"+realmName, rawBody, r.haToken)
	}
	getBody := func(t *testing.T) realmBody {
		t.Helper()
		var body realmBody
		decodeData(t, r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken), &body)
		return body
	}
	wantFieldErrors := func(t *testing.T, rawBody, want string) {
		t.Helper()
		rec := patch(t, rawBody)
		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("patch %s: got %d, want 422", rawBody, rec.Code)
		}
		if got := rec.Body.String(); got != want {
			t.Errorf("patch %s: body = %s, want %s", rawBody, got, want)
		}
	}

	t.Run("Auth401WithoutToken", func(t *testing.T) {
		if rec := r.req(t, http.MethodPatch, "/realms/"+realmName, `{"device_registration_limit":7}`, ""); rec.Code != http.StatusUnauthorized {
			t.Errorf("no token: got %d, want 401", rec.Code)
		}
	})

	t.Run("PatchRetention", func(t *testing.T) {
		if rec := patch(t, `{"datastream_maximum_storage_retention":86400}`); rec.Code != http.StatusOK || rec.Body.Len() != 0 {
			t.Fatalf("patch retention: got %d %q, want 200 empty", rec.Code, rec.Body)
		}
		if body := getBody(t); body.DatastreamMaximumStorageRetention == nil || *body.DatastreamMaximumStorageRetention != 86400 {
			t.Errorf("retention after patch = %v, want 86400", body.DatastreamMaximumStorageRetention)
		}
	})

	t.Run("PatchRegistrationLimit", func(t *testing.T) {
		if rec := patch(t, `{"device_registration_limit":7}`); rec.Code != http.StatusOK || rec.Body.Len() != 0 {
			t.Fatalf("patch limit: got %d %q, want 200 empty", rec.Code, rec.Body)
		}
		body := getBody(t)
		if body.DeviceRegistrationLimit == nil || *body.DeviceRegistrationLimit != 7 {
			t.Errorf("limit after patch = %v, want 7", body.DeviceRegistrationLimit)
		}
	})

	t.Run("PatchJWTPublicKeyPEM", func(t *testing.T) {
		newKey := pubPEM(t, mustRSAKey(t))
		pemJSON, _ := json.Marshal(newKey)
		if rec := patch(t, `{"jwt_public_key_pem":`+string(pemJSON)+`}`); rec.Code != http.StatusOK || rec.Body.Len() != 0 {
			t.Fatalf("patch pem: got %d %q, want 200 empty", rec.Code, rec.Body)
		}
		if body := getBody(t); body.JWTPublicKeyPEM != newKey {
			t.Error("pem after patch mismatch")
		}
	})

	t.Run("NullUnsets", func(t *testing.T) {
		if rec := patch(t, `{"device_registration_limit":null,"datastream_maximum_storage_retention":null}`); rec.Code != http.StatusOK || rec.Body.Len() != 0 {
			t.Fatalf("patch nulls: got %d %q, want 200 empty", rec.Code, rec.Body)
		}
		body := getBody(t)
		if body.DeviceRegistrationLimit != nil || body.DatastreamMaximumStorageRetention != nil {
			t.Errorf("fields after null patch = %v/%v, want nil/nil", body.DeviceRegistrationLimit, body.DatastreamMaximumStorageRetention)
		}
	})

	t.Run("ZeroRetentionUnsets", func(t *testing.T) {
		if rec := patch(t, `{"datastream_maximum_storage_retention":60}`); rec.Code != http.StatusOK {
			t.Fatalf("set retention: got %d (%s)", rec.Code, rec.Body)
		}
		if rec := patch(t, `{"datastream_maximum_storage_retention":0}`); rec.Code != http.StatusOK || rec.Body.Len() != 0 {
			t.Fatalf("zero-set retention: got %d %q, want 200 empty", rec.Code, rec.Body)
		}
		if body := getBody(t); body.DatastreamMaximumStorageRetention != nil {
			t.Errorf("retention after zero patch = %v, want nil", *body.DatastreamMaximumStorageRetention)
		}
	})

	t.Run("Rejections", func(t *testing.T) {
		wantFieldErrors(t, `{"datastream_maximum_storage_retention":-5}`,
			`{"errors":{"datastream_maximum_storage_retention":["is invalid"]}}`)
		wantFieldErrors(t, `{"device_registration_limit":-1}`,
			`{"errors":{"device_registration_limit":["is invalid"]}}`)
		wantFieldErrors(t, `{"jwt_public_key_pem":""}`,
			`{"errors":{"jwt_public_key_pem":["can't be blank"]}}`)
		wantFieldErrors(t, `{"replication_factor":3}`,
			`{"errors":{"error_name":["invalid_update_parameters"]}}`)
		// Rejected patches must not have mutated anything.
		if body := getBody(t); body.JWTPublicKeyPEM == "" {
			t.Error("realm lost its jwt key after rejected patches")
		}
	})

	t.Run("UnknownRealm404", func(t *testing.T) {
		rec := r.req(t, http.MethodPatch, "/realms/nope"+randSuffix(t), `{"device_registration_limit":7}`, r.haToken)
		if rec.Code != http.StatusNotFound {
			t.Errorf("unknown realm: got %d, want 404", rec.Code)
		}
	})
}

// TestHousekeepingRetentionDefault covers #73: HOUSEKEEPING_DEFAULT_...
// injects at creation ONLY when the caller omits the field.
func TestHousekeepingRetentionDefault(t *testing.T) {
	def := int64(3600)
	r := newHKRigWithDefault(t, &def)

	create := func(t *testing.T, name string, retentionJSON string) {
		t.Helper()
		body := `{"realm_name":` + jsonStr(name) + `,"jwt_public_key_pem":` + jsonStr(r.realmKey) + retentionJSON + `}`
		if rec := r.req(t, http.MethodPost, "/realms", body, r.haToken); rec.Code != http.StatusCreated {
			t.Fatalf("create %s: got %d, want 201 (%s)", name, rec.Code, rec.Body)
		}
	}
	retentionOf := func(t *testing.T, name string) *int64 {
		t.Helper()
		var body realmBody
		decodeData(t, r.req(t, http.MethodGet, "/realms/"+name, "", r.haToken), &body)
		return body.DatastreamMaximumStorageRetention
	}

	t.Run("ExplicitValueBeatsDefault", func(t *testing.T) {
		name := "pd" + randSuffix(t)
		create(t, name, `,"datastream_maximum_storage_retention":60`)
		if got := retentionOf(t, name); got == nil || *got != 60 {
			t.Errorf("explicit retention = %v, want 60", got)
		}
	})
	t.Run("OmittedInjectsDefault", func(t *testing.T) {
		name := "pi" + randSuffix(t)
		create(t, name, "")
		if got := retentionOf(t, name); got == nil || *got != def {
			t.Errorf("injected retention = %v, want %d", got, def)
		}
	})
	t.Run("NoDefaultStaysNull", func(t *testing.T) {
		r := newHKRig(t)
		name := "pn" + randSuffix(t)
		body := `{"realm_name":` + jsonStr(name) + `,"jwt_public_key_pem":` + jsonStr(r.realmKey) + `}`
		if rec := r.req(t, http.MethodPost, "/realms", body, r.haToken); rec.Code != http.StatusCreated {
			t.Fatalf("create: got %d, want 201 (%s)", rec.Code, rec.Body)
		}
		var rb realmBody
		decodeData(t, r.req(t, http.MethodGet, "/realms/"+name, "", r.haToken), &rb)
		if rb.DatastreamMaximumStorageRetention != nil {
			t.Errorf("retention without default = %v, want nil", *rb.DatastreamMaximumStorageRetention)
		}
	})
}

// TestHousekeepingDeleteGating drives #75 against upstream's wire shapes:
// the cluster flag answers 405 verbatim, connected devices answer 422 with
// the connected_devices_present error_name, and everything else keeps the
// synchronous always-delete behavior.
func TestHousekeepingDeleteGating(t *testing.T) {
	create := func(t *testing.T, r *hkRig) string {
		t.Helper()
		name := "dg" + randSuffix(t)
		body := `{"realm_name":` + jsonStr(name) + `,"jwt_public_key_pem":` + jsonStr(r.realmKey) + `}`
		if rec := r.req(t, http.MethodPost, "/realms", body, r.haToken); rec.Code != http.StatusCreated {
			t.Fatalf("create %s: got %d, want 201 (%s)", name, rec.Code, rec.Body)
		}
		return name
	}
	registerDevice := func(t *testing.T, r *hkRig, realmName string) deviceid.ID {
		t.Helper()
		svc := pairing.New(r.st, r.sealer, pairing.Config{BrokerURL: "mqtts://localhost:8883"})
		dev, err := deviceid.Random()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := svc.Register(context.Background(), realmName, dev.String(), ""); err != nil {
			t.Fatalf("registering device into %s: %v", realmName, err)
		}
		return dev
	}

	t.Run("UnknownRealm404", func(t *testing.T) {
		r := newHKRig(t)
		if rec := r.req(t, http.MethodDelete, "/realms/nope"+randSuffix(t), "", r.haToken); rec.Code != http.StatusNotFound {
			t.Errorf("delete unknown realm: got %d, want 404", rec.Code)
		}
	})

	t.Run("ConnectedDevicesPresent422", func(t *testing.T) {
		r := newHKRig(t)
		realmName := create(t, r)
		dev := registerDevice(t, r, realmName)
		realm, err := r.st.GetRealmByName(context.Background(), realmName)
		if err != nil {
			t.Fatal(err)
		}
		if err := r.st.SetDeviceConnected(context.Background(), realm.ID, dev, time.Now(), netip.MustParseAddr("10.0.0.1")); err != nil {
			t.Fatal(err)
		}
		rec := r.req(t, http.MethodDelete, "/realms/"+realmName, "", r.haToken)
		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("delete with connected device: got %d, want 422 (%s)", rec.Code, rec.Body)
		}
		if got, want := rec.Body.String(), `{"errors":{"error_name":["connected_devices_present"]}}`; got != want {
			t.Errorf("body = %s, want %s", got, want)
		}
		// Gated, not deleted: the realm must still be there.
		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusOK {
			t.Errorf("realm survived the gated delete: got %d, want 200", rec.Code)
		}

		// Same realm after disconnection deletes normally.
		if err := r.st.SetDeviceDisconnected(context.Background(), realm.ID, dev, time.Now()); err != nil {
			t.Fatal(err)
		}
		if rec := r.req(t, http.MethodDelete, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusNoContent {
			t.Errorf("delete after disconnect: got %d, want 204 (%s)", rec.Code, rec.Body)
		}
		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusNotFound {
			t.Errorf("get deleted realm: got %d, want 404", rec.Code)
		}
	})

	t.Run("FlagSet405EvenWithoutDevices", func(t *testing.T) {
		r := newHKRigWithDeletionDisabled(t)
		realmName := create(t, r)
		rec := r.req(t, http.MethodDelete, "/realms/"+realmName, "", r.haToken)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("deletion disabled: got %d, want 405 (%s)", rec.Code, rec.Body)
		}
		if got, want := rec.Body.String(), `{"errors":{"detail":"Realm deletion disabled"}}`; got != want {
			t.Errorf("body = %s, want %s", got, want)
		}
		// The gate precedes existence checks: an unknown realm gets the same 405.
		if rec := r.req(t, http.MethodDelete, "/realms/nope"+randSuffix(t), "", r.haToken); rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("deletion disabled, unknown realm: got %d, want 405", rec.Code)
		}
	})
}

// --- helpers ----------------------------------------------------------------

func discardLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func mustRSAKey(t *testing.T) *rsa.PublicKey {
	t.Helper()
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	return &k.PublicKey
}

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
