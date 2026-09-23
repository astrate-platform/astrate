package blocks

import (
	"bytes"
	"fmt"
	"math"
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
//   - max_buffered (int, default 0): maximum number of messages to buffer
//     (0 means unlimited). When exceeded, apply overflow_policy
//   - overflow_policy (string, default "drop_oldest"): action when buffer
//     exceeds max_buffered. Values: "drop_oldest", "error", "force_flush"
//
// Window semantics: a message is held until either a newer arrival pushes it
// behind the window edge of the newest timestamp, or it is superseded in the
// buffer. The len>1 flush guard keeps at least the newest message buffered —
// the tail is never emitted by this block itself; it is only released when a
// later arrival overtakes it (buffered state is lost on teardown).
//
// Overflow policy: when max_buffered > 0 and inserting a new message would
// cause the buffer to exceed max_buffered, the policy controls behavior:
//   - drop_oldest: drop the oldest buffered message(s) to make room
//   - error: return an error and do not buffer the new message
//   - force_flush: flush as many messages as possible (oldest first) to make
//     room, prioritizing releasing buffered messages over dropping them
func Sort(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
	cfg, err := parseSortConfig(config)
	if err != nil {
		return nil, fmt.Errorf("sort: %w", err)
	}
	if cfg.windowMs > math.MaxInt64/1000 {
		return nil, fmt.Errorf("sort: window_ms too large")
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

		// First, try to flush based on window with the new message
		newest := msg.Timestamp
		if len(buf) > 0 {
			if last := buf[len(buf)-1]; last.ts > newest {
				newest = last.ts
			}
		}
		var out []*flow.Message
		for len(buf) > 1 && buf[0].ts <= newest-windowUs {
			out = append(out, buf[0].msg)
			buf = buf[1:]
		}

		// Now handle overflow if max_buffered is set and buffer still exceeds limit
		if cfg.maxBuffered > 0 && int64(len(buf)) > cfg.maxBuffered {
			switch cfg.overflowPolicy {
			case overflowDropOldest:
				// Drop oldest messages to fit within max_buffered (excluding the case
				// where this would drop everything? but we need to keep space for semantics)
				excess := int64(len(buf)) - cfg.maxBuffered
				if excess > 0 {
					buf = buf[excess:]
				}
			case overflowError:
				return nil, fmt.Errorf("sort: buffer overflow (max %d messages)", cfg.maxBuffered)
			case overflowForceFlush:
				// Force flush oldest messages until we're within limit
				for int64(len(buf)) > cfg.maxBuffered && len(buf) > 0 {
					out = append(out, buf[0].msg)
					buf = buf[1:]
				}
			}
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

const (
	overflowDropOldest = "drop_oldest"
	overflowError      = "error"
	overflowForceFlush = "force_flush"
)

type sortConfig struct {
	windowMs       int64
	dedup          bool
	maxBuffered    int64
	overflowPolicy string
}

func parseSortConfig(config map[string]any) (sortConfig, error) {
	cfg := sortConfig{
		windowMs:       1000,
		maxBuffered:    0,
		overflowPolicy: overflowDropOldest,
	}
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
	if v, ok := config["max_buffered"]; ok && v != nil {
		n, err := numAsInt64(v)
		if err != nil || n < 0 {
			return cfg, fmt.Errorf("max_buffered must be non-negative")
		}
		cfg.maxBuffered = n
	}
	if v, ok := config["max_buffered_count"]; ok && v != nil && cfg.maxBuffered == 0 {
		// Alternative naming for compatibility
		n, err := numAsInt64(v)
		if err != nil || n < 0 {
			return cfg, fmt.Errorf("max_buffered_count must be non-negative")
		}
		cfg.maxBuffered = n
	}
	if v, ok := config["overflow_policy"]; ok && v != nil {
		p, ok := v.(string)
		if !ok {
			return cfg, fmt.Errorf("overflow_policy must be a string")
		}
		switch p {
		case overflowDropOldest, overflowError, overflowForceFlush:
			cfg.overflowPolicy = p
		default:
			return cfg, fmt.Errorf("overflow_policy must be one of: drop_oldest, error, force_flush")
		}
	}
	return cfg, nil
}
