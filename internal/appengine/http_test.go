//go:build integration

package appengine

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
	"net/netip"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

const (
	aeSensors = "com.ex.M7b.Sensors"
	aeConf    = "com.ex.M7b.Conf"
	sensorDef = `{"interface_name":"com.ex.M7b.Sensors","version_major":1,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"},{"endpoint":"/big","type":"longinteger"}]}`
	confDef   = `{"interface_name":"com.ex.M7b.Conf","version_major":1,"version_minor":0,"type":"properties","ownership":"device","mappings":[{"endpoint":"/%{k}","type":"string","allow_unset":true}]}`
)

// fakeServerData records the server-owned writes the AppEngine forwards.
type fakeServerData struct {
	mu     sync.Mutex
	pubs   []string
	unsets []string
}

func (f *fakeServerData) PublishServerValue(_ context.Context, realm string, id deviceid.ID, iface, path string, value json.RawMessage, _ *time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.pubs = append(f.pubs, realm+"|"+id.String()+"|"+iface+path+"|"+string(value))
	return nil
}

func (f *fakeServerData) UnsetServerProperty(_ context.Context, realm string, id deviceid.ID, iface, path string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.unsets = append(f.unsets, realm+"|"+id.String()+"|"+iface+path)
	return nil
}

type rig struct {
	st      *store.Store
	mux     *http.ServeMux
	sd      *fakeServerData
	realm   string
	realmID int16
	dev     deviceid.ID
	token   string
	wrong   string
	t2, t3  time.Time
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

	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	var suffix [4]byte
	_, _ = rand.Read(suffix[:])
	realmName := "ae" + hex.EncodeToString(suffix[:])
	realm, err := st.CreateRealm(ctx, store.NewRealm{
		Name: realmName, JWTPublicKeysPEM: []string{pubPEM(t, &key.PublicKey)},
		CACertificatePEM: "ca", CAPrivateKeySealed: []byte("k"),
	})
	if err != nil {
		t.Fatalf("CreateRealm: %v", err)
	}
	sensors, err := st.InstallInterface(ctx, realm.ID, []byte(sensorDef))
	if err != nil {
		t.Fatal(err)
	}
	conf, err := st.InstallInterface(ctx, realm.ID, []byte(confDef))
	if err != nil {
		t.Fatal(err)
	}

	dev, _ := deviceid.Random()
	if err := st.RegisterDevice(ctx, realm.ID, dev, "h"); err != nil {
		t.Fatal(err)
	}
	if _, err := st.UpdateIntrospection(ctx, realm.ID, dev, map[string]store.InterfaceVersion{
		aeSensors: {Major: 1, Minor: 0}, aeConf: {Major: 1, Minor: 0},
	}); err != nil {
		t.Fatal(err)
	}
	if err := st.SetDeviceConnected(ctx, realm.ID, dev, time.Now().UTC(), netip.MustParseAddr("127.0.0.1")); err != nil {
		t.Fatal(err)
	}
	if err := st.PatchDeviceAliases(ctx, realm.ID, dev, map[string]*string{"label": ptr("sensor-1")}); err != nil {
		t.Fatal(err)
	}

	// Seed three /value samples and one huge /big longinteger.
	t1 := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	t2 := t1.Add(time.Minute)
	t3 := t2.Add(time.Minute)
	rows := []store.IndividualRow{
		dsRow(realm.ID, dev, sensors, "/value", t1, 1.0),
		dsRow(realm.ID, dev, sensors, "/value", t2, 2.0),
		dsRow(realm.ID, dev, sensors, "/value", t3, 3.0),
	}
	big := int64(1) << 60
	rows = append(rows, store.IndividualRow{
		RealmID: realm.ID, DeviceID: dev, InterfaceID: sensors.ID, EndpointID: sensors.Endpoints["/big"],
		Path: "/big", TS: t3, ReceptionTS: t3, ValueLonginteger: &big,
	})
	if err := st.AppendDatastreams(ctx, store.DatastreamBatch{Individual: rows}); err != nil {
		t.Fatal(err)
	}

	// Seed two properties.
	for path, val := range map[string]string{"/alpha": `"a"`, "/beta": `"b"`} {
		if err := st.UpsertProperty(ctx, store.Property{
			RealmID: realm.ID, DeviceID: dev, InterfaceID: conf.ID, EndpointID: conf.Endpoints["/%{k}"],
			Path: path, Value: json.RawMessage(val), ValueType: interfaceschema.String, SetAt: t1,
		}); err != nil {
			t.Fatal(err)
		}
	}

	sd := &fakeServerData{}
	mux := http.NewServeMux()
	NewAPI(NewService(st, sd, discardLogger()), auth.NewMiddleware(st)).Mount(mux)

	return &rig{
		st: st, mux: mux, sd: sd, realm: realmName, realmID: realm.ID, dev: dev, t2: t2, t3: t3,
		token: mintToken(t, key, jwt.MapClaims{"a_aea": []string{".*::.*"}}),
		wrong: mintToken(t, key, jwt.MapClaims{"a_rma": []string{".*::.*"}}),
	}
}

func (r *rig) req(t *testing.T, method, path, rawBody, token string) *httptest.ResponseRecorder {
	t.Helper()
	var body io.Reader
	if rawBody != "" {
		body = strings.NewReader(`{"data":` + rawBody + `}`)
	}
	hr := httptest.NewRequest(method, "/appengine/v1/"+r.realm+path, body)
	if token != "" {
		hr.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	r.mux.ServeHTTP(rec, hr)
	return rec
}

func (r *rig) dpath(suffix string) string {
	return "/devices/" + r.dev.String() + suffix
}

func TestAppEngine(t *testing.T) {
	r := newRig(t)

	t.Run("Auth", func(t *testing.T) {
		if rec := r.req(t, http.MethodGet, "/devices", "", ""); rec.Code != http.StatusUnauthorized {
			t.Errorf("no token: got %d, want 401", rec.Code)
		}
		if rec := r.req(t, http.MethodGet, "/devices", "", r.wrong); rec.Code != http.StatusForbidden {
			t.Errorf("wrong claim: got %d, want 403", rec.Code)
		}
	})

	t.Run("DeviceStatus", func(t *testing.T) {
		var ds DeviceStatus
		decodeData(t, r.req(t, http.MethodGet, r.dpath(""), "", r.token), &ds)
		if !ds.Connected || ds.Aliases["label"] != "sensor-1" {
			t.Errorf("status = %+v", ds)
		}
		if _, ok := ds.Introspection[aeSensors]; !ok {
			t.Errorf("introspection missing %s", aeSensors)
		}
		if rec := r.req(t, http.MethodGet, "/devices/"+unknownID(t), "", r.token); rec.Code != http.StatusNotFound {
			t.Errorf("unknown device: got %d, want 404", rec.Code)
		}
	})

	t.Run("DatastreamQueryBoundaries", func(t *testing.T) {
		var all []Sample
		decodeData(t, r.req(t, http.MethodGet, r.dpath("/interfaces/"+aeSensors+"/value"), "", r.token), &all)
		if len(all) != 3 {
			t.Fatalf("all samples = %d, want 3", len(all))
		}
		// Default ordering is descending (newest first).
		if all[0].Timestamp.Before(all[1].Timestamp) {
			t.Errorf("not descending: %v", all)
		}
		// since is inclusive, since_after is exclusive.
		var sinceT2, afterT2 []Sample
		decodeData(t, r.req(t, http.MethodGet, r.dpath("/interfaces/"+aeSensors+"/value")+"?since="+iso(r.t2), "", r.token), &sinceT2)
		decodeData(t, r.req(t, http.MethodGet, r.dpath("/interfaces/"+aeSensors+"/value")+"?since_after="+iso(r.t2), "", r.token), &afterT2)
		if len(sinceT2) != 2 || len(afterT2) != 1 {
			t.Errorf("since=%d (want 2), since_after=%d (want 1)", len(sinceT2), len(afterT2))
		}
		var limited []Sample
		decodeData(t, r.req(t, http.MethodGet, r.dpath("/interfaces/"+aeSensors+"/value")+"?limit=1", "", r.token), &limited)
		if len(limited) != 1 {
			t.Errorf("limit=1 returned %d", len(limited))
		}
	})

	t.Run("LongIntegerAsString", func(t *testing.T) {
		var samples []Sample
		decodeData(t, r.req(t, http.MethodGet, r.dpath("/interfaces/"+aeSensors+"/big"), "", r.token), &samples)
		if len(samples) != 1 {
			t.Fatalf("samples = %d", len(samples))
		}
		if got, ok := samples[0].Value.(string); !ok || got != "1152921504606846976" {
			t.Errorf("longinteger value = %v (%T), want decimal string", samples[0].Value, samples[0].Value)
		}
	})

	t.Run("PropertyTree", func(t *testing.T) {
		var tree map[string]json.RawMessage
		decodeData(t, r.req(t, http.MethodGet, r.dpath("/interfaces/"+aeConf), "", r.token), &tree)
		if string(tree["/alpha"]) != `"a"` || string(tree["/beta"]) != `"b"` {
			t.Errorf("property tree = %v", tree)
		}
		var one json.RawMessage
		decodeData(t, r.req(t, http.MethodGet, r.dpath("/interfaces/"+aeConf+"/alpha"), "", r.token), &one)
		if string(one) != `"a"` {
			t.Errorf("property /alpha = %s", one)
		}
	})

	t.Run("ServerOwnedWrite", func(t *testing.T) {
		if rec := r.req(t, http.MethodPut, r.dpath("/interfaces/com.ex.M7b.Srv/value"), "7.5", r.token); rec.Code != http.StatusOK {
			t.Fatalf("put: got %d, want 200 (%s)", rec.Code, rec.Body)
		}
		if rec := r.req(t, http.MethodDelete, r.dpath("/interfaces/com.ex.M7b.Srv/value"), "", r.token); rec.Code != http.StatusNoContent {
			t.Fatalf("delete: got %d, want 204", rec.Code)
		}
		r.sd.mu.Lock()
		defer r.sd.mu.Unlock()
		if len(r.sd.pubs) != 1 || !strings.HasSuffix(r.sd.pubs[0], "|com.ex.M7b.Srv/value|7.5") {
			t.Errorf("captured pubs = %v", r.sd.pubs)
		}
		if len(r.sd.unsets) != 1 {
			t.Errorf("captured unsets = %v", r.sd.unsets)
		}
	})

	t.Run("PatchDevice", func(t *testing.T) {
		var ds DeviceStatus
		decodeData(t, r.req(t, http.MethodPatch, r.dpath(""), `{"aliases":{"label":"renamed"}}`, r.token), &ds)
		if ds.Aliases["label"] != "renamed" {
			t.Errorf("alias after patch = %v", ds.Aliases)
		}
	})

	t.Run("Groups", func(t *testing.T) {
		body := `{"group_name":"g1","devices":[` + jsonStr(r.dev.String()) + `]}`
		if rec := r.req(t, http.MethodPost, "/groups", body, r.token); rec.Code != http.StatusCreated {
			t.Fatalf("create group: got %d, want 201 (%s)", rec.Code, rec.Body)
		}
		var names []string
		decodeData(t, r.req(t, http.MethodGet, "/groups", "", r.token), &names)
		if !contains(names, "g1") {
			t.Errorf("groups = %v", names)
		}
		var devs []string
		decodeData(t, r.req(t, http.MethodGet, "/groups/g1/devices", "", r.token), &devs)
		if !contains(devs, r.dev.String()) {
			t.Errorf("group devices = %v", devs)
		}
		if rec := r.req(t, http.MethodDelete, "/groups/g1/devices/"+r.dev.String(), "", r.token); rec.Code != http.StatusNoContent {
			t.Errorf("remove from group: got %d, want 204", rec.Code)
		}
	})
}

// --- helpers ----------------------------------------------------------------

func dsRow(rid int16, dev deviceid.ID, si *store.StoredInterface, path string, ts time.Time, v float64) store.IndividualRow {
	val := v
	return store.IndividualRow{
		RealmID: rid, DeviceID: dev, InterfaceID: si.ID, EndpointID: si.Endpoints[path],
		Path: path, TS: ts, ReceptionTS: ts, ValueDouble: &val,
	}
}

func ptr[T any](v T) *T { return &v }

func iso(t time.Time) string { return t.UTC().Format(time.RFC3339) }

func unknownID(t *testing.T) string {
	t.Helper()
	id, _ := deviceid.Random()
	return id.String()
}

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
