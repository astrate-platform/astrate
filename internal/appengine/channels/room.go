package channels

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/engine/triggers"
)

// DefaultMailbox is the mailbox capacity used when AddMember is given a non-positive buffer.
const DefaultMailbox = 64

// Bus is the subset of *stream.Bus a room needs. *stream.Bus satisfies it.
type Bus interface {
	Subscribe(realm string, f stream.Filter, buffer int) (<-chan stream.Event, func())
}

// WatchRequest is the watch payload DTO, exactly as the client sends it.
//
// GroupName is a top-level field, measured against upstream (channels.json,
// 2026-08-22): a group_name nested inside simple_trigger is refused by
// upstream's changeset, and one at the top level is what its authorization
// path is built from. Astrate used to read it only from simple_trigger, which
// meant an upstream-shaped group watch silently degraded into a device-shaped
// path check.
type WatchRequest struct {
	Name          string          `json:"name"`
	DeviceID      string          `json:"device_id"`
	GroupName     string          `json:"group_name"`
	SimpleTrigger json.RawMessage `json:"simple_trigger"`
}

type watchEntry struct {
	name     string
	deviceID string
	tg       *triggers.Trigger
}

// Registry is a topic→room map.
type Registry struct {
	bus   Bus
	mu    sync.Mutex
	rooms map[string]*Room
}

// NewRegistry creates a new Registry.
func NewRegistry(b Bus) *Registry {
	return &Registry{bus: b, rooms: make(map[string]*Room)}
}

// Join returns the room for topic, creating it on first join.
func (r *Registry) Join(realm, topic string) *Room {
	r.mu.Lock()
	defer r.mu.Unlock()
	if rm, ok := r.rooms[topic]; ok {
		return rm
	}
	ch, cancel := r.bus.Subscribe(realm, stream.Filter{}, 0)
	rm := &Room{
		reg:      r,
		topic:    topic,
		members:  make(map[*Member]struct{}),
		events:   ch,
		cancelFn: cancel,
	}
	r.rooms[topic] = rm
	go rm.dispatch()
	return rm
}

// Rooms reports the live room count.
func (r *Registry) Rooms() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.rooms)
}

// Room is one Phoenix topic.
type Room struct {
	reg      *Registry
	topic    string
	watches  map[string]watchEntry
	watchMu  sync.Mutex
	members  map[*Member]struct{}
	memberMu sync.Mutex
	events   <-chan stream.Event
	cancelFn func()
	retired  sync.Once
}

// retire releases the bus subscription and unregisters the topic, exactly
// once however the room ends: the last member leaving, or the bus closing
// underneath it. The caller holds memberMu.
func (rm *Room) retire() {
	rm.retired.Do(func() {
		rm.cancelFn()
		rm.reg.mu.Lock()
		delete(rm.reg.rooms, rm.topic)
		rm.reg.mu.Unlock()
	})
}

// AddMember registers a member with a mailbox of the given capacity.
// A non-positive buffer selects DefaultMailbox.
func (rm *Room) AddMember(buffer int) *Member {
	if buffer <= 0 {
		buffer = DefaultMailbox
	}
	m := &Member{
		room:    rm,
		mailbox: make(chan triggers.SimpleEvent, buffer),
	}
	rm.memberMu.Lock()
	rm.members[m] = struct{}{}
	rm.memberMu.Unlock()
	return m
}

// Watch compiles and stores a watch.
func (rm *Room) Watch(req WatchRequest) error {
	rm.watchMu.Lock()
	defer rm.watchMu.Unlock()
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.SimpleTrigger) == 0 {
		return fmt.Errorf("simple_trigger is required")
	}
	tg, err := triggers.CompileCondition(req.Name, req.SimpleTrigger)
	if err != nil {
		return fmt.Errorf("CompileCondition: %w", err)
	}
	if rm.watches == nil {
		rm.watches = make(map[string]watchEntry)
	}
	rm.watches[req.Name] = watchEntry{
		name:     req.Name,
		deviceID: req.DeviceID,
		tg:       tg,
	}
	return nil
}

// Unwatch removes a watch, reporting whether one was there.
func (rm *Room) Unwatch(name string) bool {
	rm.watchMu.Lock()
	defer rm.watchMu.Unlock()
	if rm.watches == nil {
		return false
	}
	_, ok := rm.watches[name]
	delete(rm.watches, name)
	return ok
}

// Watches reports the number of active watches.
func (rm *Room) Watches() int {
	rm.watchMu.Lock()
	defer rm.watchMu.Unlock()
	return len(rm.watches)
}

// dispatch processes bus events and delivers to members.
func (rm *Room) dispatch() {
	// The bus channel closing (shutdown, or the subscription being cancelled)
	// retires the room: every mailbox closes, the subscription is released and
	// the topic is dropped from the registry, so a later Join builds a live
	// room rather than handing out this one with a dead event channel.
	defer func() {
		rm.memberMu.Lock()
		for m := range rm.members {
			close(m.mailbox)
			delete(rm.members, m)
		}
		rm.retire()
		rm.memberMu.Unlock()
	}()

	for ev := range rm.events {
		rm.watchMu.Lock()
		watches := make([]watchEntry, 0, len(rm.watches))
		for _, w := range rm.watches {
			watches = append(watches, w)
		}
		rm.watchMu.Unlock()

		rm.memberMu.Lock()
		members := make([]*Member, 0, len(rm.members))
		for m := range rm.members {
			members = append(members, m)
		}

		for _, w := range watches {
			if w.deviceID != "" && w.deviceID != ev.DeviceID {
				continue
			}

			var matched bool
			switch ev.Kind {
			case stream.KindIncomingData:
				matched = w.tg.MatchesData(triggers.DataEvent{
					DeviceID:  ev.DeviceID,
					Interface: ev.Interface,
					Major:     ev.InterfaceMajor,
					Path:      ev.Path,
					Value:     ev.Value,
				})
			case stream.KindDeviceConnected, stream.KindDeviceDisconnected, stream.KindDeviceError:
				matched = w.tg.MatchesDevice(triggers.DeviceEvent{
					DeviceID: ev.DeviceID,
					On:       ev.Kind,
				})
			default:
				continue
			}
			if !matched {
				continue
			}

			var event any
			switch ev.Kind {
			case stream.KindIncomingData:
				event = triggers.NewIncomingDataEvent(ev.Interface, ev.Path, ev.Value)
			case stream.KindDeviceConnected:
				event = triggers.NewDeviceConnectedEvent(ev.IP)
			case stream.KindDeviceDisconnected:
				event = triggers.NewDeviceDisconnectedEvent()
			case stream.KindDeviceError:
				event = triggers.NewDeviceErrorEvent(ev.ErrorName, ev.ErrorMetadata)
			default:
				continue
			}

			envelope := triggers.SimpleEvent{
				Timestamp:   ev.Timestamp,
				DeviceID:    ev.DeviceID,
				TriggerName: w.name,
				Event:       event,
			}

			for _, m := range members {
				select {
				case m.mailbox <- envelope:
				default:
					atomic.AddUint64(&m.dropped, 1)
				}
			}
		}
		rm.memberMu.Unlock()
	}
}

// Member is a joined socket with its own mailbox.
type Member struct {
	room    *Room
	mailbox chan triggers.SimpleEvent
	dropped uint64
}

// Events returns the mailbox channel.
func (m *Member) Events() <-chan triggers.SimpleEvent {
	return m.mailbox
}

// Dropped reports how many envelopes were discarded because the mailbox was full.
func (m *Member) Dropped() uint64 {
	return atomic.LoadUint64(&m.dropped)
}

// Leave deregisters the member and closes its mailbox exactly once.
func (m *Member) Leave() {
	rm := m.room
	rm.memberMu.Lock()
	defer rm.memberMu.Unlock()
	if _, ok := rm.members[m]; !ok {
		return
	}
	delete(rm.members, m)
	if len(rm.members) == 0 {
		rm.retire()
	}
	close(m.mailbox)
}
