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
	sensors *store.StoredInterface
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
		st: st, mux: mux, sd: sd, realm: realmName, realmID: realm.ID, dev: dev, sensors: sensors, t2: t2, t3: t3,
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

	t.Run("DownsampleTo", func(t *testing.T) {
		// A hundred-minute ramp on its own device, so the three samples the rig
		// seeds cannot be mistaken for the downsampled result.
		ctx := context.Background()
		dev, _ := deviceid.Random()
		if err := r.st.RegisterDevice(ctx, r.realmID, dev, "h"); err != nil {
			t.Fatal(err)
		}
		const aeObj = "com.ex.M7b.Agg"
		objDef := `{"interface_name":"` + aeObj + `","version_major":1,"version_minor":0,` +
			`"type":"datastream","ownership":"device","aggregation":"object",` +
			`"mappings":[{"endpoint":"/%{id}/temp","type":"double"}]}`
		if _, err := r.st.InstallInterface(ctx, r.realmID, []byte(objDef)); err != nil {
			t.Fatal(err)
		}
		if _, err := r.st.UpdateIntrospection(ctx, r.realmID, dev, map[string]store.InterfaceVersion{
			aeSensors: {Major: 1, Minor: 0}, aeObj: {Major: 1, Minor: 0},
		}); err != nil {
			t.Fatal(err)
		}
		base := time.Date(2026, 6, 2, 8, 0, 0, 0, time.UTC)
		const samples = 100
		var rows []store.IndividualRow
		for i := range samples {
			rows = append(rows, dsRow(r.realmID, dev, r.sensors, "/value",
				base.Add(time.Duration(i)*time.Minute), float64(i)))
		}
		if err := r.st.AppendDatastreams(ctx, store.DatastreamBatch{Individual: rows}); err != nil {
			t.Fatal(err)
		}
		last := base.Add((samples - 1) * time.Minute)
		path := "/devices/" + dev.String() + "/interfaces/" + aeSensors + "/value"

		// The whole point of the feature: 100 stored samples, 10 asked for, and
		// the bucket derived from the series' own span because the request names
		// no window. Epoch-aligned buckets can straddle one extra boundary, so
		// N±1 is the honest tolerance.
		var got []Sample
		decodeData(t, r.req(t, http.MethodGet, path+"?downsample_to=10&sort=ascending", "", r.token), &got)
		if len(got) < 9 || len(got) > 11 {
			t.Fatalf("downsample_to=10 returned %d points, want 10±1", len(got))
		}
		// The reduced series must still span the real one: no window was given,
		// so a bucket derived from anything narrower would clip the tail.
		if got[0].Timestamp.After(base) {
			t.Errorf("first bucket %v starts after the first sample %v", got[0].Timestamp, base)
		}
		if tail := got[len(got)-1].Timestamp; tail.Before(last.Add(-10 * time.Minute)) {
			t.Errorf("last bucket %v does not reach the last sample %v", tail, last)
		}
		// Averages of a monotonic ramp are themselves monotonic — the evidence
		// that these are aggregates of the series and not raw samples.
		for i := 1; i < len(got); i++ {
			prev, _ := got[i-1].Value.(float64)
			cur, _ := got[i].Value.(float64)
			if cur <= prev {
				t.Errorf("bucket %d value %v not above %v", i, cur, prev)
			}
		}

		// A finer request must yield strictly more points than a coarser one.
		var coarse []Sample
		decodeData(t, r.req(t, http.MethodGet, path+"?downsample_to=4&sort=ascending", "", r.token), &coarse)
		if len(coarse) >= len(got) {
			t.Errorf("downsample_to=4 returned %d points, downsample_to=10 returned %d", len(coarse), len(got))
		}

		// downsample_to is a point count greater than one. A malformed one is
		// rejected by the parser before any query runs — a 400, the same status
		// the other malformed query parameters get, not the 422 the service
		// layer reserves for well-formed requests that break a rule.
		for _, bad := range []string{"1", "0", "-3", "abc"} {
			if rec := r.req(t, http.MethodGet, path+"?downsample_to="+bad, "", r.token); rec.Code != http.StatusBadRequest {
				t.Errorf("downsample_to=%s: got %d, want 400", bad, rec.Code)
			}
		}

		// downsample_to=2 is legal (the parser only rejects < 2), but toolkit
		// lttb() refuses a resolution below 3 outright — "resolution must be
		// greater than 2". On a toolkit-bearing server the service must fall
		// back to the time_bucket path for it rather than surface a 500.
		var two []Sample
		rec := r.req(t, http.MethodGet, path+"?downsample_to=2&sort=ascending", "", r.token)
		if rec.Code != http.StatusOK {
			t.Fatalf("downsample_to=2: got %d, want 200 (body %s)", rec.Code, rec.Body.String())
		}
		decodeData(t, rec, &two)
		if len(two) == 0 {
			t.Error("downsample_to=2 returned no points")
		}

		// store.Downsample reads the individual-datastream table only, so an
		// object-aggregated interface cannot be downsampled. Saying so is the
		// point: without the check it would silently return an empty array.
		if rec := r.req(t, http.MethodGet,
			"/devices/"+dev.String()+"/interfaces/"+aeObj+"/12?downsample_to=5", "", r.token); rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("object aggregation: got %d, want 422", rec.Code)
		}

		// The interface root is a snapshot of every endpoint's latest sample, not
		// a series: there is no span to bucket, so asking to reduce it is a 422
		// rather than the empty array an empty-path series query would return.
		if rec := r.req(t, http.MethodGet,
			"/devices/"+dev.String()+"/interfaces/"+aeSensors+"?downsample_to=5", "", r.token); rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("interface-root downsample: got %d, want 422", rec.Code)
		}

		// Likewise a properties interface: one current value per path, no series
		// to reduce, at the root and at a concrete path alike.
		for _, p := range []string{"", "/alpha"} {
			if rec := r.req(t, http.MethodGet,
				r.dpath("/interfaces/"+aeConf+p)+"?downsample_to=5", "", r.token); rec.Code != http.StatusUnprocessableEntity {
				t.Errorf("properties downsample %q: got %d, want 422", p, rec.Code)
			}
		}

		// An endpoint with no samples downsamples to an empty array, not an error.
		var empty []Sample
		decodeData(t, r.req(t, http.MethodGet,
			"/devices/"+dev.String()+"/interfaces/"+aeSensors+"/big?downsample_to=5", "", r.token), &empty)
		if len(empty) != 0 {
			t.Errorf("empty series downsampled to %d points", len(empty))
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
		// Interface-root query is the nested upstream snapshot tree, not a flat
		// {path: value} map (astarte-go's parsePropertiesMap walks the nesting).
		var tree map[string]json.RawMessage
		decodeData(t, r.req(t, http.MethodGet, r.dpath("/interfaces/"+aeConf), "", r.token), &tree)
		if string(tree["alpha"]) != `"a"` || string(tree["beta"]) != `"b"` {
			t.Errorf("property tree = %v", tree)
		}
		var one json.RawMessage
		decodeData(t, r.req(t, http.MethodGet, r.dpath("/interfaces/"+aeConf+"/alpha"), "", r.token), &one)
		if string(one) != `"a"` {
			t.Errorf("property /alpha = %s", one)
		}
	})

	t.Run("DatastreamSnapshot", func(t *testing.T) {
		// Interface-root query on an individual datastream returns the latest
		// sample per endpoint as a nested {endpoint: {value, timestamp}} tree.
		var snap map[string]Sample
		decodeData(t, r.req(t, http.MethodGet, r.dpath("/interfaces/"+aeSensors), "", r.token), &snap)
		v, ok := snap["value"]
		if !ok {
			t.Fatalf("snapshot missing /value endpoint: %v", snap)
		}
		if got, ok := v.Value.(float64); !ok || got != 3.0 {
			t.Errorf("snapshot /value = %v (%T), want newest 3.0", v.Value, v.Value)
		}
		if !v.Timestamp.Equal(r.t3) {
			t.Errorf("snapshot /value timestamp = %v, want newest %v", v.Timestamp, r.t3)
		}
		big, ok := snap["big"]
		if !ok {
			t.Fatalf("snapshot missing /big endpoint: %v", snap)
		}
		if got, ok := big.Value.(string); !ok || got != "1152921504606846976" {
			t.Errorf("snapshot /big = %v (%T), want longinteger string", big.Value, big.Value)
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
		decodeData(t, r.reqCT(t, http.MethodPatch, r.dpath(""),
			`{"aliases":{"label":"renamed"}}`, r.token, "application/merge-patch+json"), &ds)
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
