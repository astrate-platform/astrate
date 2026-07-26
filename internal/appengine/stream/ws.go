// Package stream is the AppEngine live event socket (docs/DESIGN.md §3.7,
// §1.1 deviation; ROADMAP §8.2 file 7.9): a WebSocket (with an SSE fallback)
// at /astrate/v1/{realm}/socket fed by the engine's in-process fan-out bus.
// Subscribers receive the realm's committed data and device lifecycle events
// as JSON; an a_ch JWT guards the endpoint, and device_id/interface query
// parameters narrow the room.
//
// This is a deliberate departure from upstream Astarte's Phoenix Channels
// wire protocol (docs/DESIGN.md §1.1): the surface is a plain WebSocket/SSE,
// not phx_join frames, so it is documented as Astrate-specific.
package stream

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/coder/websocket"

	"github.com/astrate-platform/astrate/internal/auth"
	engstream "github.com/astrate-platform/astrate/internal/engine/stream"
)

// subscribeBuffer is the per-connection bus buffer; a slow socket drops
// events at the bus (docs/DESIGN.md §1.4), never blocking ingestion.
const subscribeBuffer = 256

// Bus is the subset of the engine fan-out bus the socket consumes.
type Bus interface {
	Subscribe(realm string, f engstream.Filter, buffer int) (<-chan engstream.Event, func())
}

// API is the live-socket HTTP surface, guarded by an a_ch realm JWT.
type API struct {
	bus     Bus
	require func(http.Handler) http.Handler
}

// NewAPI wires the socket to the engine bus and the a_ch middleware.
func NewAPI(bus Bus, mw *auth.Middleware) *API {
	return &API{bus: bus, require: mw.RequireRealm(auth.ClaimChannels)}
}

// Mount registers the socket route on mux.
func (a *API) Mount(mux *http.ServeMux) {
	mux.Handle("GET /astrate/v1/{realm}/socket", a.require(http.HandlerFunc(a.handle)))
}

// wireEvent is the JSON shape pushed to subscribers.
type wireEvent struct {
	Event     string    `json:"event"`
	Realm     string    `json:"realm"`
	DeviceID  string    `json:"device_id"`
	Interface string    `json:"interface,omitempty"`
	Path      string    `json:"path,omitempty"`
	Value     any       `json:"value,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func toWire(ev *engstream.Event) wireEvent {
	return wireEvent{
		Event: ev.Kind, Realm: ev.Realm, DeviceID: ev.DeviceID,
		Interface: ev.Interface, Path: ev.Path, Value: ev.Value, Timestamp: ev.Timestamp,
	}
}

// handle upgrades to WebSocket, or streams SSE when the client asks for
// text/event-stream.
func (a *API) handle(w http.ResponseWriter, r *http.Request) {
	realm := r.PathValue("realm")
	filter := engstream.Filter{
		DeviceID:  r.URL.Query().Get("device_id"),
		Interface: r.URL.Query().Get("interface"),
	}
	events, cancel := a.bus.Subscribe(realm, filter, subscribeBuffer)
	defer cancel()

	if wantsSSE(r) {
		a.serveSSE(w, r, events)
		return
	}
	a.serveWebSocket(w, r, events)
}

// serveWebSocket pumps bus events to the client until either side closes.
func (a *API) serveWebSocket(w http.ResponseWriter, r *http.Request, events <-chan engstream.Event) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return // Accept already wrote the failure
	}
	defer func() { _ = conn.CloseNow() }()
	// CloseRead drains client frames and gives us a context that cancels when
	// the client disconnects (the socket is server-push only).
	ctx := conn.CloseRead(r.Context())

	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-events:
			if !ok {
				_ = conn.Close(websocket.StatusGoingAway, "bus closed")
				return
			}
			payload, err := json.Marshal(toWire(&ev))
			if err != nil {
				continue
			}
			if err := conn.Write(ctx, websocket.MessageText, payload); err != nil {
				return
			}
		}
	}
}

// serveSSE streams bus events as Server-Sent Events.
func (a *API) serveSSE(w http.ResponseWriter, r *http.Request, events <-chan engstream.Event) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-events:
			if !ok {
				return
			}
			payload, err := json.Marshal(toWire(&ev))
			if err != nil {
				continue
			}
			if _, err := w.Write(sseFrame(payload)); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

// wantsSSE reports whether the client requested an SSE stream.
func wantsSSE(r *http.Request) bool {
	return r.Header.Get("Accept") == "text/event-stream" || r.URL.Query().Get("transport") == "sse"
}

// sseFrame wraps a JSON payload in an SSE "data:" frame.
func sseFrame(payload []byte) []byte {
	out := make([]byte, 0, len(payload)+8)
	out = append(out, "data: "...)
	out = append(out, payload...)
	out = append(out, '\n', '\n')
	return out
}
