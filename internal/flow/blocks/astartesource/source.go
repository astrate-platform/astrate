// Package astartesource implements the AstarteSource Flow block (issue #27,
// astarte_flow parity): a Source that subscribes to Astrate's existing live
// event bus (internal/engine/stream) and converts device events into
// FlowMessages, connecting device ingestion to operator-defined pipelines.
//
// Wiring (issue #37): the flow manager's source pump calls Emit; StopFlow
// calls Stop to drop the bus subscription.
package astartesource

import (
	"context"
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
//
// Source implements flow.Block, flow.Source, and flow.Stopper.
type Source struct {
	cfg    Config
	ch     <-chan stream.Event
	cancel func()
}

// New subscribes to bus for cfg.Realm/cfg.Interface and returns a Source
// ready to be used as a flow.Block. Call Stop when the flow using it stops,
// to release the subscription (Manager.StopFlow does this automatically when
// the block is in the flow graph).
func New(bus *stream.Bus, cfg Config) *Source {
	ch, cancel := bus.Subscribe(cfg.Realm, stream.Filter{Interface: cfg.Interface}, 0)
	return &Source{cfg: cfg, ch: ch, cancel: cancel}
}

// Name implements flow.Block.
func (s *Source) Name() string { return "astarte_source" }

// Process implements flow.Block as a non-blocking drain of currently buffered
// events. Prefer Emit for the live pump path (it waits for the next event).
func (s *Source) Process(_ *flow.Message) ([]*flow.Message, error) {
	return s.drain(context.TODO(), false)
}

// Emit implements flow.Source: it blocks until at least one accepted event is
// available or ctx is cancelled, then drains any further buffered events.
func (s *Source) Emit(ctx context.Context) ([]*flow.Message, error) {
	return s.drain(ctx, true)
}

// drain converts bus events into FlowMessages. When block is true it waits
// for the first accepted event (or ctx cancellation); otherwise it only
// takes events already buffered.
func (s *Source) drain(ctx context.Context, block bool) ([]*flow.Message, error) {
	var out []*flow.Message

	if block {
		for len(out) == 0 {
			select {
			case <-ctx.Done():
				return out, ctx.Err()
			case ev, ok := <-s.ch:
				if !ok {
					return out, nil
				}
				if s.cfg.Path != "" && !strings.HasPrefix(ev.Path, s.cfg.Path) {
					continue
				}
				out = append(out, toFlowMessage(&ev))
			}
		}
	}

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

// Stop implements flow.Stopper: unsubscribes from the bus. Safe to call once;
// the underlying channel is closed by the bus, after which Process/Emit
// return cleanly with no events.
func (s *Source) Stop() {
	s.cancel()
}

// toFlowMessage converts one stream.Event to the upstream Message shape
// (astarte_flow/message/v0.1): key "<realm>/<device_id>", data/subtype
// events map from Event.Value; lifecycle events (no Interface/Path) carry
// their kind and metadata instead.
func toFlowMessage(ev *stream.Event) *flow.Message {
	msg := &flow.Message{
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
func setValue(msg *flow.Message, v any) {
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
