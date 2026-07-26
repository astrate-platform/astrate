package channels

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

	"github.com/coder/websocket"
	"github.com/golang-jwt/jwt/v5"

	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/store"
)

const testRealmName = "testrealm"

type fakeRealmKeys struct {
	realm *store.Realm
}

func (k fakeRealmKeys) GetRealmByName(_ context.Context, name string) (*store.Realm, error) {
	if name == k.realm.Name {
		return k.realm, nil
	}
	return nil, store.ErrNotFound
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

func setupTestAPI(t *testing.T) (*API, *rsa.PrivateKey) {
	t.Helper()
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	der, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
	pub := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
	realm := &store.Realm{Name: testRealmName, JWTPublicKeysPEM: []string{pub}}
	bus := &fakeBus{ch: make(chan stream.Event)}
	api := NewAPI(bus, fakeRealmKeys{realm})
	return api, key
}

func TestPreUpgrade_MissingToken(t *testing.T) {
	api, _ := setupTestAPI(t)
	mux := http.NewServeMux()
	api.Mount(mux)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?realm=testrealm", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want 401", rec.Code)
	}
}

func TestPreUpgrade_MissingRealm(t *testing.T) {
	api, key := setupTestAPI(t)
	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{".*::.*"}})
	mux := http.NewServeMux()
	api.Mount(mux)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?token="+token, nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want 401", rec.Code)
	}
}

func TestPreUpgrade_UnknownRealm(t *testing.T) {
	api, key := setupTestAPI(t)
	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{".*::.*"}})
	mux := http.NewServeMux()
	api.Mount(mux)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?realm=nosuchrealm&token="+token, nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want 401", rec.Code)
	}
}

func TestPreUpgrade_InvalidToken(t *testing.T) {
	api, _ := setupTestAPI(t)
	mux := http.NewServeMux()
	api.Mount(mux)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?realm=testrealm&token=not.a.jwt", nil)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("got %d, want 401", rec.Code)
	}
}

func TestPreUpgrade_HappyPath(t *testing.T) {
	api, key := setupTestAPI(t)
	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{".*::.*"}})
	mux := http.NewServeMux()
	api.Mount(mux)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/appengine/v1/socket/websocket?realm=testrealm&token="+token, nil)
	mux.ServeHTTP(rec, req)
	if rec.Code == http.StatusUnauthorized {
		t.Errorf("happy path got 401, should have upgraded")
	}
}

func dialSession(t *testing.T, mux *http.ServeMux, token string) (*websocket.Conn, context.CancelFunc) {
	t.Helper()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	url := strings.Replace(srv.URL, "http", "ws", 1) +
		"/appengine/v1/socket/websocket?vsn=2.0.0&realm=" + testRealmName + "&token=" + token
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		cancel()
		t.Fatalf("Dial: %v", err)
	}
	return conn, cancel
}

func TestWireSession_Heartbeat(t *testing.T) {
	api, key := setupTestAPI(t)
	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{".*::.*"}})
	mux := http.NewServeMux()
	api.Mount(mux)
	conn, cancel := dialSession(t, mux, token)
	defer cancel()
	defer conn.CloseNow()

	ctx := context.Background()
	if err := conn.Write(ctx, websocket.MessageText, []byte(`[null,"1","phoenix","heartbeat",{}]`)); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, raw, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	// Check wire format: join_ref (element 0) must be null, not "".
	var arr []json.RawMessage
	if err := json.Unmarshal(raw, &arr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(arr) != 5 {
		t.Fatalf("got %d elements, want 5", len(arr))
	}
	if string(arr[0]) != "null" {
		t.Errorf("join_ref = %s, want null", string(arr[0]))
	}
	if string(arr[1]) != `"1"` {
		t.Errorf("ref = %s, want \"1\"", string(arr[1]))
	}

	// Check decoded Frame.
	var f Frame
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatalf("decode frame: %v", err)
	}
	if f.Event != EventPhxReply {
		t.Errorf("event = %s, want %s", f.Event, EventPhxReply)
	}
	var payload struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(f.Payload, &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if payload.Status != "ok" {
		t.Errorf("status = %s, want ok", payload.Status)
	}
}

func TestWireSession_JoinOwnRealm(t *testing.T) {
	api, key := setupTestAPI(t)
	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{".*::.*"}})
	mux := http.NewServeMux()
	api.Mount(mux)
	conn, cancel := dialSession(t, mux, token)
	defer cancel()
	defer conn.CloseNow()

	ctx := context.Background()
	topic := "rooms:" + testRealmName + ":dashboard_1"
	if err := conn.Write(ctx, websocket.MessageText, []byte(`["1","1","`+topic+`","phx_join",{}]`)); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, raw, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var f Frame
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if f.Event != EventPhxReply {
		t.Errorf("event = %s, want %s", f.Event, EventPhxReply)
	}
	var payload struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(f.Payload, &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if payload.Status != "ok" {
		t.Errorf("status = %s, want ok", payload.Status)
	}
	if api.reg.Rooms() != 1 {
		t.Errorf("Rooms() = %d, want 1", api.reg.Rooms())
	}
}

func TestWireSession_JoinRejected(t *testing.T) {
	api, key := setupTestAPI(t)
	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{".*::.*"}})
	mux := http.NewServeMux()
	api.Mount(mux)
	conn, cancel := dialSession(t, mux, token)
	defer cancel()
	defer conn.CloseNow()

	ctx := context.Background()

	cases := []struct {
		name  string
		topic string
	}{
		{"other realm", "rooms:otherrealm:dashboard_1"},
		{"no room name", "rooms:" + testRealmName},
		{"empty room name", "rooms:" + testRealmName + ":"},
		{"wrong prefix", "lobby:" + testRealmName + ":dashboard_1"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			roomsBefore := api.reg.Rooms()
			msg := `["1","1","` + c.topic + `","phx_join",{}]`
			if err := conn.Write(ctx, websocket.MessageText, []byte(msg)); err != nil {
				t.Fatalf("write: %v", err)
			}

			_, raw, err := conn.Read(ctx)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			var f Frame
			if err := json.Unmarshal(raw, &f); err != nil {
				t.Fatalf("decode: %v", err)
			}
			var payload struct {
				Status string `json:"status"`
			}
			if err := json.Unmarshal(f.Payload, &payload); err != nil {
				t.Fatalf("decode payload: %v", err)
			}
			if payload.Status != "error" {
				t.Errorf("status = %s, want error", payload.Status)
			}
			// This row must be refused by its own rule, not merely produce
			// some error: no room may exist for the rejected topic.
			if api.reg.Rooms() != roomsBefore {
				t.Errorf("Rooms() = %d, want %d (%s must create no room)",
					api.reg.Rooms(), roomsBefore, c.name)
			}
		})
	}
}

func TestWireSession_Leave(t *testing.T) {
	api, key := setupTestAPI(t)
	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{".*::.*"}})
	mux := http.NewServeMux()
	api.Mount(mux)
	conn, cancel := dialSession(t, mux, token)
	defer cancel()
	defer conn.CloseNow()

	ctx := context.Background()
	topic := "rooms:" + testRealmName + ":dashboard_1"

	// Join first.
	if err := conn.Write(ctx, websocket.MessageText, []byte(`["1","1","`+topic+`","phx_join",{}]`)); err != nil {
		t.Fatalf("write join: %v", err)
	}
	if _, _, err := conn.Read(ctx); err != nil {
		t.Fatalf("read join reply: %v", err)
	}
	if api.reg.Rooms() != 1 {
		t.Fatalf("after join, Rooms() = %d, want 1", api.reg.Rooms())
	}

	// Leave.
	if err := conn.Write(ctx, websocket.MessageText, []byte(`["2","2","`+topic+`","phx_leave",{}]`)); err != nil {
		t.Fatalf("write leave: %v", err)
	}
	_, raw, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read leave reply: %v", err)
	}
	var f Frame
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatalf("decode: %v", err)
	}
	var payload struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(f.Payload, &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if payload.Status != "ok" {
		t.Errorf("status = %s, want ok", payload.Status)
	}
	if api.reg.Rooms() != 0 {
		t.Errorf("Rooms() = %d, want 0 after leave", api.reg.Rooms())
	}
}

func TestWireSession_LeaveWithoutJoin(t *testing.T) {
	api, key := setupTestAPI(t)
	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{".*::.*"}})
	mux := http.NewServeMux()
	api.Mount(mux)
	conn, cancel := dialSession(t, mux, token)
	defer cancel()
	defer conn.CloseNow()

	ctx := context.Background()
	topic := "rooms:" + testRealmName + ":dashboard_1"
	if err := conn.Write(ctx, websocket.MessageText, []byte(`["1","1","`+topic+`","phx_leave",{}]`)); err != nil {
		t.Fatalf("write: %v", err)
	}
	_, raw, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var f Frame
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatalf("decode: %v", err)
	}
	var payload struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(f.Payload, &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if payload.Status != "error" {
		t.Errorf("status = %s, want error", payload.Status)
	}
}

func setupTestAPIWithBus(t *testing.T) (*API, *rsa.PrivateKey, *fakeBus) {
	t.Helper()
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	der, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
	pub := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
	realm := &store.Realm{Name: testRealmName, JWTPublicKeysPEM: []string{pub}}
	bus := &fakeBus{ch: make(chan stream.Event, 16)}
	api := NewAPI(bus, fakeRealmKeys{realm})
	return api, key, bus
}

func joinTopic(t *testing.T, conn *websocket.Conn, topic, ref string) {
	t.Helper()
	msg := `["` + ref + `","` + ref + `","` + topic + `","phx_join",{}]`
	if err := conn.Write(context.Background(), websocket.MessageText, []byte(msg)); err != nil {
		t.Fatalf("write join: %v", err)
	}
}

func readFramePayloadStatus(t *testing.T, conn *websocket.Conn) string {
	t.Helper()
	_, raw, err := conn.Read(context.Background())
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var f Frame
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatalf("decode: %v", err)
	}
	var payload struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(f.Payload, &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	return payload.Status
}

func TestWireSession_Watch_Accepted(t *testing.T) {
	api, key, _ := setupTestAPIWithBus(t)
	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{".*::.*"}})
	mux := http.NewServeMux()
	api.Mount(mux)
	conn, cancel := dialSession(t, mux, token)
	defer cancel()
	defer conn.CloseNow()

	ctx := context.Background()
	topic := "rooms:" + testRealmName + ":dashboard_1"
	joinTopic(t, conn, topic, "1")
	if _, _, err := conn.Read(ctx); err != nil {
		t.Fatalf("read join reply: %v", err)
	}

	watchPayload := `{"name":"connectiontrigger-D","device_id":"` + validDeviceIDA + `","simple_trigger":{"type":"device_trigger","on":"device_connected","device_id":"` + validDeviceIDA + `"}}`
	msg := `["1","2","` + topic + `","watch",` + watchPayload + `]`
	if err := conn.Write(ctx, websocket.MessageText, []byte(msg)); err != nil {
		t.Fatalf("write watch: %v", err)
	}

	status := readFramePayloadStatus(t, conn)
	if status != "ok" {
		t.Errorf("watch status = %s, want ok", status)
	}

	rm := api.reg.Join(testRealmName, topic)
	if rm.Watches() != 1 {
		t.Errorf("Watches() = %d, want 1", rm.Watches())
	}
}

func TestWireSession_Watch_Rejected(t *testing.T) {
	api, key, _ := setupTestAPIWithBus(t)
	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{".*::.*"}})
	restrictiveToken := mintToken(t, key, jwt.MapClaims{"a_ch": []string{"WATCH::" + validDeviceIDA}})
	mux := http.NewServeMux()
	api.Mount(mux)

	topic := "rooms:" + testRealmName + ":dashboard_1"
	devATrigger := `{"name":"connectiontrigger-D","device_id":"` + validDeviceIDA + `","simple_trigger":{"type":"device_trigger","on":"device_connected","device_id":"` + validDeviceIDA + `"}}`
	devBTrigger := `{"name":"connectiontrigger-B","device_id":"` + validDeviceIDB + `","simple_trigger":{"type":"device_trigger","on":"device_connected","device_id":"` + validDeviceIDB + `"}}`

	cases := []struct {
		name    string
		token   string
		topic   string
		payload string
	}{
		{
			name:    "no name",
			token:   token,
			topic:   topic,
			payload: `{"name":"","device_id":"` + validDeviceIDA + `","simple_trigger":{"type":"device_trigger","on":"device_connected","device_id":"` + validDeviceIDA + `"}}`,
		},
		{
			name:    "no simple_trigger",
			token:   token,
			topic:   topic,
			payload: `{"name":"x","device_id":"` + validDeviceIDA + `"}`,
		},
		{
			name:    "not joined",
			token:   token,
			topic:   "rooms:" + testRealmName + ":never_joined",
			payload: devATrigger,
		},
		{
			name:    "unauthorized device",
			token:   restrictiveToken,
			topic:   topic,
			payload: devBTrigger,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			conn, cancel := dialSession(t, mux, c.token)
			defer cancel()
			defer conn.CloseNow()
			ctx := context.Background()

			// Join the topic if it's not the "not joined" case.
			if c.topic == topic {
				joinTopic(t, conn, topic, "1")
				if _, _, err := conn.Read(ctx); err != nil {
					t.Fatalf("read join reply: %v", err)
				}
			}

			rm := api.reg.Join(testRealmName, topic)
			watchesBefore := rm.Watches()

			msg := `["1","2","` + c.topic + `","watch",` + c.payload + `]`
			if err := conn.Write(ctx, websocket.MessageText, []byte(msg)); err != nil {
				t.Fatalf("write watch: %v", err)
			}

			status := readFramePayloadStatus(t, conn)
			if status != "error" {
				t.Errorf("status = %s, want error", status)
			}
			if rm.Watches() != watchesBefore {
				t.Errorf("Watches() = %d, want %d (must not increase)", rm.Watches(), watchesBefore)
			}
		})
	}

	// The "unauthorized device" subtest above uses a restrictive token that
	// should deny B but allow A. Verify A is accepted on the same token.
	conn, cancel := dialSession(t, mux, restrictiveToken)
	defer cancel()
	defer conn.CloseNow()
	ctx := context.Background()

	joinTopic(t, conn, topic, "1")
	if _, _, err := conn.Read(ctx); err != nil {
		t.Fatalf("read join reply: %v", err)
	}

	rm := api.reg.Join(testRealmName, topic)
	watchesBefore := rm.Watches()
	msg := `["1","2","` + topic + `","watch",` + devATrigger + `]`
	if err := conn.Write(ctx, websocket.MessageText, []byte(msg)); err != nil {
		t.Fatalf("write watch: %v", err)
	}
	status := readFramePayloadStatus(t, conn)
	if status != "ok" {
		t.Errorf("restrictive token watch for A: status = %s, want ok", status)
	}
	if rm.Watches() != watchesBefore+1 {
		t.Errorf("Watches() = %d, want %d after accepted watch", rm.Watches(), watchesBefore+1)
	}
}

func TestWireSession_NewEvent(t *testing.T) {
	api, key, bus := setupTestAPIWithBus(t)
	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{".*::.*"}})
	mux := http.NewServeMux()
	api.Mount(mux)
	conn, cancel := dialSession(t, mux, token)
	defer cancel()
	defer conn.CloseNow()

	ctx := context.Background()
	topic := "rooms:" + testRealmName + ":dashboard_1"
	joinTopic(t, conn, topic, "1")

	// Read and discard the join reply.
	_, _, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read join reply: %v", err)
	}

	// Register a watch.
	watchPayload := `{"name":"connectiontrigger-D","device_id":"` + validDeviceIDA + `","simple_trigger":{"type":"device_trigger","on":"device_connected","device_id":"` + validDeviceIDA + `"}}`
	watchMsg := `["1","2","` + topic + `","watch",` + watchPayload + `]`
	if err := conn.Write(ctx, websocket.MessageText, []byte(watchMsg)); err != nil {
		t.Fatalf("write watch: %v", err)
	}

	// Read and discard the watch reply.
	_, _, err = conn.Read(ctx)
	if err != nil {
		t.Fatalf("read watch reply: %v", err)
	}

	// Push an event into the bus.
	now := time.Now().UTC().Truncate(time.Millisecond)
	bus.ch <- stream.Event{Kind: stream.KindDeviceConnected, DeviceID: validDeviceIDA, IP: "1.2.3.4", Timestamp: now}

	// Read the new_event frame.
	_, raw, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read new_event: %v", err)
	}

	// Parse as raw []json.RawMessage to check element-by-element.
	var arr []json.RawMessage
	if err := json.Unmarshal(raw, &arr); err != nil {
		t.Fatalf("unmarshal raw: %v", err)
	}
	if len(arr) != 5 {
		t.Fatalf("got %d elements, want 5", len(arr))
	}

	// Element 0: join_ref must be "1", not null.
	if string(arr[0]) != `"1"` {
		t.Errorf("join_ref = %s, want \"1\"", string(arr[0]))
	}

	// Element 1: ref must be null.
	if string(arr[1]) != "null" {
		t.Errorf("ref = %s, want null", string(arr[1]))
	}

	// Element 3: event must be "new_event".
	if string(arr[3]) != `"new_event"` {
		t.Errorf("event = %s, want \"new_event\"", string(arr[3]))
	}

	// Element 4: payload must decode to an object with non-empty device_id, timestamp, event.
	var evPayload struct {
		DeviceID  string `json:"device_id"`
		Timestamp string `json:"timestamp"`
		Event     any    `json:"event"`
	}
	if err := json.Unmarshal(arr[4], &evPayload); err != nil {
		t.Fatalf("decode event payload: %v", err)
	}
	if evPayload.DeviceID == "" {
		t.Error("device_id is empty")
	}
	if evPayload.Timestamp == "" {
		t.Error("timestamp is empty")
	}
	if evPayload.Event == nil {
		t.Error("event is nil")
	}
}

// TestWireSession_PumpRacesLeave drives bus events at a joined room while the
// session leaves it, so the pump goroutine and the read loop touch s.rooms and
// joinRef concurrently. The specced tests all push events with the room
// quiescent, which leaves the guarding of that map unexercised: without this,
// removing the pump's locking still passes the -race gate. The assertions are
// the race detector and not panicking; frame contents are covered elsewhere.
func TestWireSession_PumpRacesLeave(t *testing.T) {
	api, key, bus := setupTestAPIWithBus(t)
	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{".*::.*"}})
	mux := http.NewServeMux()
	api.Mount(mux)
	conn, cancel := dialSession(t, mux, token)
	defer cancel()
	defer conn.CloseNow()

	ctx := context.Background()
	topic := "rooms:" + testRealmName + ":dashboard_race"

	if err := conn.Write(ctx, websocket.MessageText, []byte(`["1","1","`+topic+`","phx_join",{}]`)); err != nil {
		t.Fatalf("write join: %v", err)
	}
	if _, _, err := conn.Read(ctx); err != nil {
		t.Fatalf("read join reply: %v", err)
	}
	watch := `["1","2","` + topic + `","watch",{"name":"connectiontrigger-D","device_id":"` +
		validDeviceIDA + `","simple_trigger":{"type":"device_trigger","on":"device_connected","device_id":"` +
		validDeviceIDA + `"}}]`
	if err := conn.Write(ctx, websocket.MessageText, []byte(watch)); err != nil {
		t.Fatalf("write watch: %v", err)
	}
	if _, _, err := conn.Read(ctx); err != nil {
		t.Fatalf("read watch reply: %v", err)
	}

	// Drain the socket so the pump is never the thing that blocks.
	readCtx, stopReading := context.WithCancel(ctx)
	defer stopReading()
	go func() {
		for {
			if _, _, err := conn.Read(readCtx); err != nil {
				return
			}
		}
	}()

	// Feed the bus while the leave lands. Sends are non-blocking: once the room
	// retires nothing consumes the channel, and this must not wedge the test.
	fed := make(chan struct{})
	go func() {
		defer close(fed)
		for i := 0; i < 200; i++ {
			select {
			case bus.ch <- stream.Event{
				Kind: stream.KindDeviceConnected, DeviceID: validDeviceIDA,
				IP: "1.2.3.4", Timestamp: time.Now().UTC(),
			}:
			default:
			}
		}
	}()

	if err := conn.Write(ctx, websocket.MessageText, []byte(`["1","3","`+topic+`","phx_leave",{}]`)); err != nil {
		t.Fatalf("write leave: %v", err)
	}
	<-fed
}
