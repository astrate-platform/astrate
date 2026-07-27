package channels

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

// writeTimeout bounds a single frame write, so one stalled viewer cannot hold
// the session's write mutex indefinitely.
const writeTimeout = 5 * time.Second

// Event-name constants for the watch/unwatch/new_event cycle.
const (
	EventWatch    = "watch"
	EventUnwatch  = "unwatch"
	EventNewEvent = "new_event"
)

// API serves the Phoenix V2 WebSocket endpoint.
type API struct {
	reg   *Registry
	keys  auth.KeySource
	cache *auth.Cache
}

// NewAPI creates an API backed by the given bus and key source.
func NewAPI(bus Bus, keys auth.KeySource) *API {
	return &API{
		reg:   NewRegistry(bus),
		keys:  keys,
		cache: auth.NewCache(auth.DefaultCacheSize),
	}
}

// Mount registers the WebSocket handler on the given mux.
func (a *API) Mount(mux *http.ServeMux) {
	mux.Handle("GET /appengine/v1/socket/websocket", http.HandlerFunc(a.handle))
}

// handle upgrades the HTTP connection and runs the session loop.
func (a *API) handle(w http.ResponseWriter, r *http.Request) {
	realm := r.URL.Query().Get("realm")
	token := r.URL.Query().Get("token")

	if realm == "" || token == "" {
		_ = astarteapi.WriteUnauthorized(w)
		return
	}

	row, err := a.keys.GetRealmByName(r.Context(), realm)
	switch {
	case errors.Is(err, store.ErrNotFound):
		_ = astarteapi.WriteUnauthorized(w)
		return
	case err != nil:
		_ = astarteapi.WriteInternalServerError(w)
		return
	}

	tok, err := a.cache.Verify(token, row.JWTPublicKeysPEM)
	if err != nil {
		_ = astarteapi.WriteUnauthorized(w)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return
	}

	s := &session{
		conn:  conn,
		ctx:   r.Context(),
		realm: realm,
		tok:   tok,
		reg:   a.reg,
		rooms: make(map[string]*joined),
	}
	s.loop(r.Context())
}

// session holds per-connection state for the WebSocket session loop.
type session struct {
	conn    *websocket.Conn
	ctx     context.Context
	realm   string
	tok     *auth.Token
	reg     *Registry
	mu      sync.Mutex
	roomsMu sync.Mutex
	rooms   map[string]*joined
}

// joined tracks a topic this session has joined.
type joined struct {
	room    *Room
	member  *Member
	joinRef *string
}

// triggerShape decodes the fields needed to derive the WATCH authorization path.
type triggerShape struct {
	Type          string `json:"type"`
	InterfaceName string `json:"interface_name"`
	MatchPath     string `json:"match_path"`
	GroupName     string `json:"group_name"`
	DeviceID      string `json:"device_id"`
}

// errUnauthorizedTrigger is a device trigger that names no device, or names a
// different one than the request. Upstream refuses these — see watchAuthPath.
var errUnauthorizedTrigger = errors.New("device trigger does not name the request's device")

// watchAuthPath derives the authorization path from a WatchRequest, in the
// shape upstream matches a_ch WATCH claims against. The shapes are recorded in
// test/conformance/upstream/channels.json:
//
//	data trigger        <device_id>/<interface_name><match_path>
//	device trigger      <device_id>
//	group data trigger  groups/<name>/<interface_name><match_path>
//	group device trigger groups/<name>
//
// The match path carries its own leading slash, so it is concatenated without a
// separator. That detail is measured, not inferred: a claim written as
// "<device>/<interface>" — the shape Astrate used to build — is refused by
// upstream for a trigger on /value, while "<device>/<interface>/value" is
// accepted. Only the group shapes are unmeasured, the realm having no group;
// they come from the same upstream function whose device shapes the recording
// just confirmed.
//
// A device trigger must name its device inside simple_trigger, and it must be
// the request's device: upstream refuses a device_id present only at the
// payload's top level (where the AppEngine REST API puts it) and refuses the
// wildcard "*", both with the reason "unauthorized". Astrate used to fall back
// to the top-level device_id and accept both.
func watchAuthPath(req WatchRequest) (string, error) {
	var ts triggerShape
	if err := json.Unmarshal(req.SimpleTrigger, &ts); err != nil {
		return "", err
	}
	if ts.Type == "data_trigger" {
		if ts.GroupName != "" {
			return "groups/" + ts.GroupName + "/" + ts.InterfaceName + ts.MatchPath, nil
		}
		return req.DeviceID + "/" + ts.InterfaceName + ts.MatchPath, nil
	}
	if ts.GroupName != "" {
		return "groups/" + ts.GroupName, nil
	}
	if ts.DeviceID == "" || ts.DeviceID != req.DeviceID {
		return "", errUnauthorizedTrigger
	}
	return ts.DeviceID, nil
}

// loop reads frames from the connection and dispatches them.
func (s *session) loop(ctx context.Context) {
	defer s.teardown()

	for {
		typ, msg, err := s.conn.Read(ctx)
		if err != nil {
			return
		}
		if typ != websocket.MessageText {
			continue
		}

		var f Frame
		if err := json.Unmarshal(msg, &f); err != nil {
			continue
		}

		switch {
		case f.Topic == TopicHeartbeat && f.Event == EventHeartbeat:
			s.handleHeartbeat(f)

		case f.Event == EventPhxJoin:
			s.handleJoin(f)

		case f.Event == EventPhxLeave:
			s.handleLeave(f)

		case f.Event == EventWatch:
			s.handleWatch(f)

		case f.Event == EventUnwatch:
			s.handleUnwatch(f)

		default:
			if f.Ref != nil {
				s.sendErr(f, "unmatched topic")
			}
		}
	}
}

// teardown leaves every room the session joined.
func (s *session) teardown() {
	s.roomsMu.Lock()
	for _, j := range s.rooms {
		j.member.Leave()
	}
	s.rooms = make(map[string]*joined)
	s.roomsMu.Unlock()
	_ = s.conn.CloseNow()
}

func (s *session) handleHeartbeat(f Frame) {
	rep, _ := OK(f, nil)
	s.writeFrame(rep)
}

func (s *session) handleJoin(f Frame) {
	prefix := "rooms:" + s.realm + ":"
	if !strings.HasPrefix(f.Topic, prefix) {
		s.sendErr(f, "unauthorized")
		return
	}
	name := strings.TrimPrefix(f.Topic, prefix)
	if name == "" {
		s.sendErr(f, "unauthorized")
		return
	}
	// Upstream authorizes every join, including a rejoin, against the room
	// name with the realm stripped off.
	if !s.tok.AuthorizesChannel(auth.VerbJoin, name) {
		s.sendErr(f, "unauthorized")
		return
	}

	s.roomsMu.Lock()
	if j, ok := s.rooms[f.Topic]; ok {
		// A rejoin on the same socket carries a fresh join_ref; server pushes
		// must be tagged with the current one.
		j.joinRef = f.Ref
		s.roomsMu.Unlock()
		rep, _ := OK(f, nil)
		s.writeFrame(rep)
		return
	}

	rm := s.reg.Join(s.realm, f.Topic)
	m := rm.AddMember(0)
	s.rooms[f.Topic] = &joined{room: rm, member: m, joinRef: f.Ref}
	s.roomsMu.Unlock()

	go s.pumpEvents(f.Topic)

	rep, _ := OK(f, nil)
	s.writeFrame(rep)
}

func (s *session) handleLeave(f Frame) {
	s.roomsMu.Lock()
	j, ok := s.rooms[f.Topic]
	if !ok {
		s.roomsMu.Unlock()
		s.sendErr(f, "unmatched topic")
		return
	}
	j.member.Leave()
	delete(s.rooms, f.Topic)
	s.roomsMu.Unlock()
	rep, _ := OK(f, nil)
	s.writeFrame(rep)
}

func (s *session) handleWatch(f Frame) {
	s.roomsMu.Lock()
	j, ok := s.rooms[f.Topic]
	s.roomsMu.Unlock()
	if !ok {
		s.sendErr(f, "unmatched topic")
		return
	}

	var req WatchRequest
	if err := json.Unmarshal(f.Payload, &req); err != nil {
		s.sendErr(f, "invalid payload")
		return
	}

	path, err := watchAuthPath(req)
	if errors.Is(err, errUnauthorizedTrigger) {
		// Upstream's reason for this is "unauthorized", not a payload error,
		// even though the cause is the payload's shape.
		s.sendErr(f, "unauthorized")
		return
	}
	if err != nil {
		s.sendErr(f, "invalid simple_trigger")
		return
	}
	if !s.tok.AuthorizesChannel(auth.VerbWatch, path) {
		s.sendErr(f, "unauthorized")
		return
	}

	if err := j.room.Watch(req); err != nil {
		s.sendErr(f, "invalid trigger")
		return
	}
	rep, _ := OK(f, nil)
	s.writeFrame(rep)
}

func (s *session) handleUnwatch(f Frame) {
	s.roomsMu.Lock()
	j, ok := s.rooms[f.Topic]
	s.roomsMu.Unlock()
	if !ok {
		s.sendErr(f, "unmatched topic")
		return
	}

	var payload struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(f.Payload, &payload)
	j.room.Unwatch(payload.Name)
	rep, _ := OK(f, nil)
	s.writeFrame(rep)
}

// pumpEvents ranges over member events and pushes new_event frames to the client.
func (s *session) pumpEvents(topic string) {
	s.roomsMu.Lock()
	j, ok := s.rooms[topic]
	if !ok {
		s.roomsMu.Unlock()
		return
	}
	events := j.member.Events()
	s.roomsMu.Unlock()

	for ev := range events {
		s.roomsMu.Lock()
		j, ok := s.rooms[topic]
		if !ok {
			s.roomsMu.Unlock()
			return
		}
		joinRef := j.joinRef
		s.roomsMu.Unlock()

		b, err := json.Marshal(ev)
		if err != nil {
			continue
		}
		s.writeFrame(Frame{
			JoinRef: joinRef,
			Ref:     nil,
			Topic:   topic,
			Event:   EventNewEvent,
			Payload: b,
		})
	}
}

func (s *session) sendErr(f Frame, reason string) {
	rep, _ := Err(f, reason)
	s.writeFrame(rep)
}

// writeFrame serialises and sends a frame under mu so that two goroutines
// never corrupt the WebSocket stream.
//
// The write is bounded by writeTimeout rather than running on a background
// context: a viewer that has stopped reading must not wedge this mutex, which
// the event pump also takes (docs/DESIGN.md §1.4 — a slow viewer drops frames,
// it never backpressures anyone else).
func (s *session) writeFrame(f Frame) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, err := json.Marshal(f)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(s.ctx, writeTimeout)
	defer cancel()
	_ = s.conn.Write(ctx, websocket.MessageText, b)
}
