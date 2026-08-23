package blocks

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/astrate-platform/astrate/internal/flow"
)

// Sort buffers messages and releases them in ascending Timestamp order once
// a newer arrival makes them older than the window behind the newest buffered
// timestamp.
//
// Config keys:
//   - window_ms (int, default 1000): how long a message may trail the newest
//     buffered timestamp before it is released; must be >= 0
//   - dedup (bool, default false): drop a message whose full wire encoding
//     equals one currently buffered
//
// Window semantics: a message is held until either a newer arrival pushes it
// behind the window edge of the newest timestamp, or it is superseded in the
// buffer. The len>1 flush guard keeps at least the newest message buffered —
// the tail is never emitted by this block itself; it is only released when a
// later arrival overtakes it (buffered state is lost on teardown).
func Sort(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
	cfg, err := parseSortConfig(config)
	if err != nil {
		return nil, fmt.Errorf("sort: %w", err)
	}
	windowUs := cfg.windowMs * 1000
	var buf []*sortEntry
	return flow.NewTransformBlock(name, func(msg *flow.Message) ([]*flow.Message, error) {
		if msg == nil {
			return nil, nil
		}
		wire, err := msg.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("sort: %w", err)
		}
		if cfg.dedup {
			for _, e := range buf {
				if bytes.Equal(e.wire, wire) {
					return nil, nil
				}
			}
		}
		// Insert keeping ascending Timestamp order; equal timestamps keep
		// arrival order (the new message goes AFTER existing equals).
		i := sort.Search(len(buf), func(i int) bool { return buf[i].ts > msg.Timestamp })
		buf = append(buf, nil)
		copy(buf[i+1:], buf[i:])
		buf[i] = &sortEntry{msg: msg, ts: msg.Timestamp, wire: wire}

		newest := msg.Timestamp
		if last := buf[len(buf)-1]; last.ts > newest {
			newest = last.ts
		}
		var out []*flow.Message
		for len(buf) > 1 && buf[0].ts <= newest-windowUs {
			out = append(out, buf[0].msg)
			buf = buf[1:]
		}
		return out, nil
	}), nil
}

// sortEntry is one buffered message with its precomputed wire encoding.
type sortEntry struct {
	msg  *flow.Message
	ts   int64
	wire []byte
}

type sortConfig struct {
	windowMs int64
	dedup    bool
}

func parseSortConfig(config map[string]any) (sortConfig, error) {
	cfg := sortConfig{windowMs: 1000}
	if v, ok := config["window_ms"]; ok && v != nil {
		n, err := numAsInt64(v)
		if err != nil || n < 0 {
			return cfg, fmt.Errorf("window_ms must be non-negative")
		}
		cfg.windowMs = n
	}
	if b, ok := config["dedup"].(bool); ok {
		cfg.dedup = b
	}
	return cfg, nil
}
