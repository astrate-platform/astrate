package stream

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/astrate-platform/astrate/internal/auth"
	engstream "github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/store"
)

// closedBus yields one pre-loaded event then closes — the handler writes the
// event and returns, making the SSE assertion deterministic without sleeps.
type closedBus struct {
	ev      engstream.Event
	gotReal string
	gotF    engstream.Filter
}

func (b *closedBus) Subscribe(realm string, f engstream.Filter, _ int) (<-chan engstream.Event, func()) {
	b.gotReal, b.gotF = realm, f
	ch := make(chan engstream.Event, 1)
	ch <- b.ev
	close(ch)
	return ch, func() {}
}

// fakeKeys is a KeySource returning one realm's JWT keys.
type fakeKeys struct{ realm *store.Realm }

func (k fakeKeys) GetRealmByName(_ context.Context, name string) (*store.Realm, error) {
	if name == k.realm.Name {
		return k.realm, nil
	}
	return nil, store.ErrNotFound
}

func TestSocketSSE(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	der, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
	pub := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
	realm := &store.Realm{Name: "testrealm", JWTPublicKeysPEM: []string{pub}}

	bus := &closedBus{ev: engstream.Event{
		Kind: engstream.KindIncomingData, Realm: "testrealm", DeviceID: "dev1",
		Interface: "com.ex.S", Path: "/v", Value: 1.5, Timestamp: time.Unix(0, 0).UTC(),
	}}
	mux := http.NewServeMux()
	NewAPI(bus, auth.NewMiddleware(fakeKeys{realm})).Mount(mux)

	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{".*::.*"}})

	t.Run("Unauthorized", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/astrate/v1/testrealm/socket", nil)
		req.Header.Set("Accept", "text/event-stream")
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("no token: got %d, want 401", rec.Code)
		}
	})

	t.Run("StreamsEventAsSSE", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/astrate/v1/testrealm/socket?device_id=dev1&interface=com.ex.S", nil)
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Authorization", "Bearer "+token)
		mux.ServeHTTP(rec, req)

		if ct := rec.Header().Get("Content-Type"); ct != "text/event-stream" {
			t.Errorf("content-type = %q", ct)
		}
		body := rec.Body.String()
		if !strings.HasPrefix(body, "data: ") || !strings.HasSuffix(body, "\n\n") {
			t.Fatalf("not an SSE frame: %q", body)
		}
		var ev wireEvent
		if err := json.Unmarshal([]byte(strings.TrimSpace(strings.TrimPrefix(body, "data: "))), &ev); err != nil {
			t.Fatalf("frame payload: %v (%q)", err, body)
		}
		if ev.Event != engstream.KindIncomingData || ev.DeviceID != "dev1" || ev.Interface != "com.ex.S" {
			t.Errorf("event = %+v", ev)
		}
		// The room filter reached the bus.
		if bus.gotReal != "testrealm" || bus.gotF.DeviceID != "dev1" || bus.gotF.Interface != "com.ex.S" {
			t.Errorf("subscribe got realm=%q filter=%+v", bus.gotReal, bus.gotF)
		}
	})
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
