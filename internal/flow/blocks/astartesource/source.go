// Package astartesource implements the AstarteSource Flow block (issue #27,
// astarte_flow parity): a Source that subscribes to Astrate's existing live
// event bus (internal/engine/stream) and converts device events into
// FlowMessages, connecting device ingestion to operator-defined pipelines.
package astartesource

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/flow"
)

// Config selects which events a Source consumes.
type Config struct {
	// Realm is the tenant to subscribe to (required).
	Realm string
	// Interface, when set, keeps only that interface's data events;
	// lifecycle events (no interface) are dropped when this is set.
	Interface string
	// Path, when set, keeps only data events whose path has this prefix.
	Path string
}

// Source is a Flow Source block backed by a stream.Bus subscription. The
// zero value is not usable; construct with New.
type Source struct {
	cfg    Config
	ch     <-chan stream.Event
	cancel func()
}

// New subscribes to bus for cfg.Realm/cfg.Interface and returns a Source
// ready to be used as a flow.Block. Call Stop when the flow using it stops,
// to release the subscription.
func New(bus *stream.Bus, cfg Config) *Source {
	ch, cancel := bus.Subscribe(cfg.Realm, stream.Filter{Interface: cfg.Interface}, 0)
	return &Source{cfg: cfg, ch: ch, cancel: cancel}
}

// Name implements flow.Block.
func (s *Source) Name() string { return "astarte_source" }

// Process implements flow.Block as a Source: it drains every event
// currently buffered on the subscription (non-blocking) and converts each
// to a FlowMessage. Called repeatedly by the router's polling loop.
func (s *Source) Process(_ *flow.FlowMessage) ([]*flow.FlowMessage, error) {
	var out []*flow.FlowMessage
	for {
		select {
		case ev, ok := <-s.ch:
			if !ok {
				return out, nil
			}
			if s.cfg.Path != "" && !strings.HasPrefix(ev.Path, s.cfg.Path) {
				continue
			}
			out = append(out, toFlowMessage(&ev))
		default:
			return out, nil
		}
	}
}

// Stop unsubscribes from the bus. Safe to call once; the underlying channel
// is closed by the bus, after which Process returns cleanly with no events.
func (s *Source) Stop() {
	s.cancel()
}

// toFlowMessage converts one stream.Event to the upstream FlowMessage shape
// (astarte_flow/message/v0.1): key "<realm>/<device_id>", data/subtype
// events map from Event.Value; lifecycle events (no Interface/Path) carry
// their kind and metadata instead.
func toFlowMessage(ev *stream.Event) *flow.FlowMessage {
	msg := &flow.FlowMessage{
		Key:       ev.Realm + "/" + ev.DeviceID,
		Timestamp: ev.Timestamp.UnixMicro(),
		Metadata: map[string]string{
			"kind":      ev.Kind,
			"interface": ev.Interface,
			"path":      ev.Path,
		},
	}
	setValue(msg, ev.Value)
	return msg
}

// setValue maps a stream.Event's JSON-friendly Value into the message's
// typed Data/Type, defaulting to a string rendering for anything that
// isn't one of the wire-representable scalar kinds.
func setValue(msg *flow.FlowMessage, v any) {
	switch val := v.(type) {
	case nil:
		msg.Type = flow.TypeString
		msg.Data = ""
	case bool:
		msg.Type = flow.TypeBoolean
		msg.Data = val
	case float64:
		msg.Type = flow.TypeReal
		msg.Data = val
	case int64:
		msg.Type = flow.TypeInteger
		msg.Data = val
	case string:
		msg.Type = flow.TypeString
		msg.Data = val
	case map[string]any:
		msg.Type = flow.TypeMap
		msg.FieldTypes = make(map[string]flow.DataType, len(val))
		data := make(map[string]any, len(val))
		for k, fv := range val {
			switch fval := fv.(type) {
			case bool:
				msg.FieldTypes[k] = flow.TypeBoolean
				data[k] = fval
			case float64:
				msg.FieldTypes[k] = flow.TypeReal
				data[k] = fval
			default:
				msg.FieldTypes[k] = flow.TypeString
				data[k] = fmt.Sprint(fval)
			}
		}
		msg.Data = data
	default:
		msg.Type = flow.TypeString
		msg.Data = strconv.Quote(fmt.Sprint(val))
	}
}
