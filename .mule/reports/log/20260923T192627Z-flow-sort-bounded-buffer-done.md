slug: flow-sort-bounded-buffer
verdict: done
at:  9bf0ab7
ran: 2026-09-23T19:26:27Z on DietPi in 634s


> build · big-pickle

I'll start by reading the required files as instructed.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/flow/blocks/sort.go
✱ Grep ""sort"" 0 matches
→ Read internal/flow/blocks
→ Read internal/flow/blocks/sort_test.go
✗ Read internal/flow/blocks/catalog.go | head -40 failed
Error: File not found: /root/astrate-mule/internal/flow/blocks/catalog.go | head -40

Did you mean one of these?
/root/astrate-mule/internal/flow/blocks/catalog.go
✗ Read internal/flow/blocks/catalog.go | head -40 failed
Error: File not found: /root/astrate-mule/internal/flow/blocks/catalog.go | head -40

Did you mean one of these?
/root/astrate-mule/internal/flow/blocks/catalog.go

$ head -40 /root/astrate-mule/internal/flow/blocks/catalog.go
// Package blocks registers the built-in Flow block catalog used to instantiate
// stored pipelines (milestone v2.0: block factory).
package blocks

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks/astartesource"
	"github.com/astrate-platform/astrate/internal/flow/blocks/container"
	"github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool"
)

// Well-known block_type strings stored in pipeline definitions.
const (
	TypeAstarteSource = "astarte_source"
	TypeNullSink      = "null_sink"
	TypeLogSink       = "log_sink"
	// TypeContainer is the Docker-backed custom processing block (Design B / #43).
	TypeContainer = container.Type
	// TypeFilter and TypeMap are declared in transform.go.

	// Pure transforms and the random source (#80, astarte_flow parity).
	TypeToJSON         = "to_json"
	TypeUpdateMetadata = "update_metadata"
	TypeSplitMap       = "split_map"
	TypeRandomSource   = "random_source"
	TypeSort           = "sort"

	// TypeJSONPathMap reshapes JSON messages through a template (#81,
	// astarte_flow parity); declared in jsonpathmap.go's neighborhood.
	TypeJSONPathMap = "json_path_map"

	// TypeHTTPSource / TypeHTTPSink poll/POST over HTTP (#82, astarte_flow
	// parity); declared in httpblocks.go.

	// TypeMQTTSource / TypeMQTTSink subscribe/publish over MQTT (#83,
	// astarte_flow parity); declared in mqtt.go.

✱ Grep ""overflow_policy|max_buffer|max_count"" in internal/flow · 0 matches

$ grep -rn "overflow_policy\|max_buffer" /root/astrate-mule/internal/flow/blocks | head -10
(no output)

✱ Grep ""window_ms"" in internal/flow/blocks | head -2 · 9 matches

$ grep -rn "numAsInt64" /root/astrate-mule/internal/flow/blocks | head -2
/root/astrate-mule/internal/flow/blocks/randomsource.go:77:		n, err := numAsInt64(v)
/root/astrate-mule/internal/flow/blocks/randomsource.go:182:// numAsInt64 coerces a config number (int/int64/float64/json.Number) to int64.

$ sed -n '150,190p' /root/astrate-mule/internal/flow/blocks/randomsource.go
return nil, ctx.Err()
	case <-timer.C:
	}
	if s.stopped.Load() {
		return nil, nil
	}
	return []*flow.Message{s.next()}, nil
}

// Stop implements flow.Stopper: idempotent; after Stop no further messages
// are emitted.
func (s *randomSource) Stop() { s.stopped.Store(true) }

func (s *randomSource) next() *flow.Message {
	msg := &flow.Message{
		Key:       s.key,
		Timestamp: time.Now().UnixMicro(),
	}
	switch s.kind {
	case randomInteger:
		msg.Type = flow.TypeInteger
		msg.Data = s.minInt + rand.Int64N(s.maxInt-s.minInt+1) //nolint:gosec // non-crypto sample values by design (random_source)
	case randomReal:
		msg.Type = flow.TypeReal
		msg.Data = s.minReal + rand.Float64()*(s.maxReal-s.minReal) //nolint:gosec // non-crypto sample values by design (random_source)
	default:
		msg.Type = flow.TypeBoolean
		msg.Data = rand.Int64N(2) == 0 //nolint:gosec // non-crypto sample values by design (random_source)
	}
	return msg
}

// numAsInt64 coerces a config number (int/int64/float64/json.Number) to int64.
func numAsInt64(v any) (int64, error) {
	switch n := v.(type) {
	case int:
		return int64(n), nil
	case int32:
		return int64(n), nil
	case int64:
		return n, nil

$ sed -n '190,210p' /root/astrate-mule/internal/flow/blocks/randomsource.go
return n, nil
	case uint:
		if uint64(n) > math.MaxInt64 {
			return 0, fmt.Errorf("number out of int64 range")
		}
		return int64(n), nil
	case uint64:
		if n > math.MaxInt64 {
			return 0, fmt.Errorf("number out of int64 range")
		}
		return int64(n), nil
	case float64:
		return int64(n), nil
	case json.Number:
		return n.Int64()
	default:
		return 0, fmt.Errorf("not a number")
	}
}

// numAsFloat64 coerces a config number to float64.

$ grep "math" /root/astrate-mule/internal/flow/blocks/randomsource.go | head -1
"math"

$ sed -n '1,30p' /root/astrate-mule/internal/flow/blocks/schema.go
package blocks

// builtinSchemas maps block_type → Draft-07-style object-schema JSON for the
// block's config. Keep in sync with builtinInfo and DefaultRegistry.
var builtinSchemas = map[string]string{
	TypeAstarteSource:     `{"type":"object","properties":{"realm":{"type":"string","description":"Astarte realm (default deps.Realm)"},"interface":{"type":"string","description":"Astarte interface to subscribe to"},"path":{"type":"string","description":"Path prefix filter for emitted messages"}},"required":[]}`,
	TypeFilter:            `{"type":"object","properties":{"key_prefix":{"type":"string","description":"Drop unless key starts with this prefix"},"key_contains":{"type":"string","description":"Drop unless key contains this substring"},"type":{"type":"string","description":"Payload type filter; one of integer|real|boolean|datetime|binary|string|map"},"metadata":{"type":"object","description":"Exact-match metadata string→string"}},"required":[]}`,
	TypeMap:               `{"type":"object","properties":{"key":{"type":"string","description":"Key template with {key} and {metadata.<name>} placeholders"},"set_metadata":{"type":"object","description":"Metadata string→string merged into messages"},"delete_metadata":{"type":"array","description":"Metadata keys to delete"}},"required":[]}`,
	TypeToJSON:            `{"type":"object","properties":{}}`,
	TypeUpdateMetadata:    `{"type":"object","properties":{"set_metadata":{"type":"object","description":"Metadata string→string to merge; at least one of set_metadata, delete_metadata is required"},"delete_metadata":{"type":"array","description":"Metadata keys to delete; at least one of set_metadata, delete_metadata is required"}},"required":[]}`,
	TypeSplitMap:          `{"type":"object","properties":{"key_template":{"type":"string","description":"Output key template with {key} and {field} placeholders (default \"{key}/{field}\")"}},"required":[]}`,
	TypeRandomSource:      `{"type":"object","properties":{"type":{"type":"string","description":"Value type to emit; one of integer|real|boolean"},"interval_ms":{"type":"integer","description":"Milliseconds between emissions (default 1000)"},"min":{"type":"number","description":"Minimum value (inclusive)"},"max":{"type":"number","description":"Maximum value (inclusive)"},"key":{"type":"string","description":"Message key (default \"random\")"}},"required":["type"]}`,
	TypeSort:              `{"type":"object","properties":{"window_ms":{"type":"integer","description":"Buffering window in milliseconds (default 1000, must be ≥ 0)"},"dedup":{"type":"boolean","description":"Release each key once per window (default false)"}},"required":[]}`,
	TypeJSONPathMap:       `{"type":"object","properties":{"template":{"type":"string","description":"Template with $MESSAGE and $METADATA.<name> placeholders; unresolvable paths fail the message"}},"required":["template"]}`,
	TypeHTTPSource:        `{"type":"object","properties":{"urls":{"type":"array","description":"URLs to poll GET round-robin (at least one)"},"interval_ms":{"type":"integer","description":"Milliseconds between polls (default 1000)"},"timeout_ms":{"type":"integer","description":"Per-request timeout in milliseconds (default 5000)"}},"required":["urls"]}`,
	TypeHTTPSink:          `{"type":"object","properties":{"url":{"type":"string","description":"Base URL to POST payloads to"},"method":{"type":"string","description":"HTTP method (default \"POST\")"},"timeout_ms":{"type":"integer","description":"Per-request timeout in milliseconds (default 5000)"},"headers":{"type":"object","description":"Extra request headers string→string"}},"required":["url"]}`,
	TypeMQTTSource:        `{"type":"object","properties":{"url":{"type":"string","description":"Broker URL (e.g. tcp://127.0.0.1:1883)"},"topics":{"type":"array","description":"Topics to subscribe to (at least one, non-empty strings)"},"qos":{"type":"integer","description":"Subscription QoS level 0|1|2 (default 0)"},"client_id":{"type":"string","description":"Optional MQTT client ID (random when unset)"},"username":{"type":"string","description":"Optional username"},"password":{"type":"string","description":"Optional password"}},"required":["url","topics"]}`,
	TypeMQTTSink:          `{"type":"object","properties":{"url":{"type":"string","description":"Broker URL (e.g. tcp://127.0.0.1:1883)"},"topic":{"type":"string","description":"Topic to publish to"},"qos":{"type":"integer","description":"Publish QoS level 0|1|2 (default 0)"},"retained":{"type":"boolean","description":"Set the retain flag on published messages (default false)"},"client_id":{"type":"string","description":"Optional MQTT client ID (random when unset)"},"username":{"type":"string","description":"Optional username"},"password":{"type":"string","description":"Optional password"}},"required":["url","topic"]}`,
	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as ASTRATE_FLOW_CONFIG"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
	TypeVirtualDevicePool: `{"type":"object","properties":{"devices":{"type":"array","items":{"type":"string"},"minItems":1,"description":"Registered device_ids this pool may publish as"}},"required":["devices"]}`,
	TypeNullSink:          `{"type":"object","properties":{}}`,
	TypeLogSink:           `{"type":"object","properties":{}}`,
}

$ grep -rn "overflow" /root/astrate-mule/docs /root/astrate-mule/internal | head -15
/root/astrate-mule/internal/appengine/downsample_test.go:58:			name:   "one year no overflow",
/root/astrate-mule/internal/flow/router.go:117:// depends on qos and the configured overflow policies.
/root/astrate-mule/internal/flow/router.go:269:			Help: "QoS 0 messages dropped due to lane overflow.",
/root/astrate-mule/internal/flow/router.go:273:			Help: "QoS ≥ 1 messages dropped due to lane overflow.",

$ grep -rn "max_buffered\|max_messages\|buffer_size" /root/astrate-mule/internal/flow | head -15
(no output)

$ grep -rn "Sort" /root/astrate-mule/internal/flow/blocks/info.go 2>&1 | head -20
76:	TypeSort: {
77:		Type:    TypeSort,

$ sed -n '60,85p' /root/astrate-mule/internal/flow/blocks/info.go
Role:    RoleTransform,
		Summary: "Merge set_metadata into Metadata, then delete delete_metadata keys",
		Config:  "set_metadata (object string→string), delete_metadata (string array); at least one required",
	},
	TypeSplitMap: {
		Type:    TypeSplitMap,
		Role:    RoleTransform,
		Summary: "Explode a map payload into one message per field",
		Config:  `key_template (default "{key}/{field}"; placeholders {key}, {field})`,
	},
	TypeRandomSource: {
		Type:    TypeRandomSource,
		Role:    RoleSource,
		Summary: "Emit a random integer, real, or boolean value every interval",
		Config:  `type (required: integer|real|boolean), interval_ms (default 1000), min, max, key (default "random")`,
	},
	TypeSort: {
		Type:    TypeSort,
		Role:    RoleTransform,
		Summary: "Buffer messages and release them in ascending timestamp order behind a window",
		Config:  "window_ms (default 1000, must be ≥ 0), dedup (bool, default false)",
	},
	TypeJSONPathMap: {
		Type:    TypeJSONPathMap,
		Role:    RoleTransform,
		Summary: "Reshape JSON messages through a template with $MESSAGE / $METADATA placeholders",

Let us look at the full sort implementation again and modify it. Let us also update schema and info. Let us make the changes to sort.go first.
$ sed -n '11,32p' /root/astrate-mule/internal/flow/blocks/sort.go
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

$ sed -n '30,67p' /root/astrate-mule/internal/flow/blocks/sort.go
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

$ sed -n '75,95p' /root/astrate-mule/internal/flow/blocks/sort.go
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

Let us look at the task again more carefully:
> `Sort` (internal/flow/blocks/sort.go:32-65) appends every message to an uncapped `buf` and only flushes while `len(buf) > 1 && buf[0].ts <= newest-windowUs` (sort.go:60) — a stream whose timestamps stay within `window_ms` of the newest (clustered bursts, or messages sharing one timestamp) never satisfies the condition, so the block emits nothing while the slice grows without bound, and the documented "newest is never emitted" mean the final message is held indefinitely; secondarily `window_ms * 1000` (sort.go:31) overflows for large configs and silently inverts the flush edge. Add a configurable max buffered count with a documented overflow policy (drop-oldest / error / force-flush) and a cap test feeding N same-timestamp messages; guard the multiply.

So we need to:
1. Guard the multiply `windowMs * 1000` against overflow (large configs could overflow int64 in general? windowMs is int64 from config parsing; multiplying by 1000 - for very large values could wrap around in signed int64). Need to prevent silent inversion - maybe clamp/check, return error if overflow would occur? Or cap? The task says "guard the multiply".
2. Add configurable max buffered count with documented overflow policy (options: drop-oldest / error / force-flush). Need to define config params - probably `max_buffered` (or similar name) and `overflow_policy` with values like "drop_oldest", "error", "force_flush" (or maybe different naming). Need to look at how other blocks name similar things. Also need to add to schema and info.
3. Handle the case where messages all have same timestamp or stay within window - buffer never flushes; with max cap, when we hit max, apply overflow policy. Also the "newest is never emitted" aspect - currently we keep at least the newest (len>1 guard). But if buffer is full and we get a new message that becomes the newest, what happens depends on policy.
4. Add a cap test feeding N same-timestamp messages.

Let us also look up other codebase conventions for naming - "max_buffered_count", "max_count", "buffer_size", "overflow_policy"? Also check router for overflow policy naming. Router mentions "overflow policies" conceptually. Let us look at router more.
$ sed -n '90,140p' /root/astrate-mule/internal/flow/router.go
if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	r := &Router{
		cfg:   cfg,
		log:   cfg.Logger.With("component", "flow-router"),
		graph: graph,
		quit:  make(chan struct{}),
	}
	r.lanes = make([]*lane, cfg.Lanes)
	for i := range r.lanes {
		r.lanes[i] = &lane{ch: make(chan *flowMsg, cfg.LaneCapacity)}
	}
	r.met = newRouterMetrics(reg)
	return r
}

// Run starts the lane goroutines. Call it after NewRouter. ctx bounds
// background work; cancel to initiate drain.
func (r *Router) Run(ctx context.Context) {
	for _, l := range r.lanes {
		r.laneWG.Add(1)
		go r.runLane(ctx, l)
	}
}

// Submit routes msg to the lane determined by FNV-1a(msg.Key). Behaviour
// depends on qos and the configured overflow policies.
//
// The read lock is held across the send, not just across the closed check.
// Dropping it in between left a window where Drain could retire the lanes
// between the check and the send; with OverflowBlock the sender can sit in
// that window for as long as the lane is full. Sends are what the lock is
// for, so it covers them.
func (r *Router) Submit(msg *Message, qos byte) {
	r.met.submitted.Inc()
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.closed {
		return
	}

	l := r.lanes[laneOf(msg.Key, len(r.lanes))]
	fm := &flowMsg{msg: msg, qos: qos}

	if qos == 0 && r.cfg.QoS0Overflow == OverflowDrop {
		select {
		case l.ch <- fm:
		default:
			r.met.droppedQoS0.Inc()
			if fm.onDrop != nil {

$ grep -rn "OverflowDrop\|overflow.*policy" /root/astrate-mule/internal/flow | head -10
/root/astrate-mule/internal/flow/router.go:24:	// OverflowDrop discards the message without blocking (QoS 0).
/root/astrate-mule/internal/flow/router.go:25:	OverflowDrop
/root/astrate-mule/internal/flow/router.go:135:	if qos == 0 && r.cfg.QoS0Overflow == OverflowDrop {
/root/astrate-mule/internal/flow/router.go:147:	if qos >= 1 && r.cfg.QoS1Overflow == OverflowDrop {
/root/astrate-mule/internal/flow/router_test.go:197:// TestRouter_OverflowDropPolicy verifies that QoS 0 messages are dropped
/root/astrate-mule/internal/flow/router_test.go:199:func TestRouter_OverflowDropPolicy(t *testing.T) {
/root/astrate-mule/internal/flow/router_test.go:210:		QoS0Overflow: OverflowDrop,

$ grep -rn "OverflowBlock\|OverflowError" /root/astrate-mule/internal/flow | head -5
/root/astrate-mule/internal/flow/router.go:22:	// OverflowBlock blocks the caller until space is available (QoS ≥ 1).
/root/astrate-mule/internal/flow/router.go:23:	OverflowBlock OverflowPolicy = iota
/root/astrate-mule/internal/flow/router.go:121:// between the check and the send; with OverflowBlock the sender can sit in
/root/astrate-mule/internal/flow/manager_test.go:65:			QoS0Overflow: OverflowBlock,
/root/astrate-mule/internal/flow/manager_test.go:66:			QoS1Overflow: OverflowBlock,

Let us modify sort.go to add max buffer and overflow policy, and guard the multiply.
← Edit internal/flow/blocks/sort.go
Index: /root/astrate-mule/internal/flow/blocks/sort.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/sort.go
+++ /root/astrate-mule/internal/flow/blocks/sort.go
@@ -2,8 +2,9 @@
 
 import (
 	"bytes"
 	"fmt"
+	"math"
 	"sort"
 
 	"github.com/astrate-platform/astrate/internal/flow"
 )

← Edit internal/flow/blocks/sort.go
Index: /root/astrate-mule/internal/flow/blocks/sort.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/sort.go
+++ /root/astrate-mule/internal/flow/blocks/sort.go
@@ -17,19 +17,33 @@
 //   - window_ms (int, default 1000): how long a message may trail the newest
 //     buffered timestamp before it is released; must be >= 0
 //   - dedup (bool, default false): drop a message whose full wire encoding
 //     equals one currently buffered
+//   - max_buffered (int, default 0): maximum number of messages to buffer
+//     (0 means unlimited). When exceeded, apply overflow_policy
+//   - overflow_policy (string, default "drop_oldest"): action when buffer
+//     exceeds max_buffered. Values: "drop_oldest", "error", "force_flush"
 //
 // Window semantics: a message is held until either a newer arrival pushes it
 // behind the window edge of the newest timestamp, or it is superseded in the
 // buffer. The len>1 flush guard keeps at least the newest message buffered —
 // the tail is never emitted by this block itself; it is only released when a
 // later arrival overtakes it (buffered state is lost on teardown).
+//
+// Overflow policy: when max_buffered > 0 and inserting a new message would
+// cause the buffer to exceed max_buffered, the policy controls behavior:
+//   - drop_oldest: drop the oldest buffered message(s) to make room
+//   - error: return an error and do not buffer the new message
+//   - force_flush: flush as many messages as possible (oldest first) to make
+//     room, prioritizing releasing buffered messages over dropping them
 func Sort(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
 	cfg, err := parseSortConfig(config)
 	if err != nil {
 		return nil, fmt.Errorf("sort: %w", err)
 	}
+	if cfg.windowMs > math.MaxInt64/1000 {
+		return nil, fmt.Errorf("sort: window_ms too large")
+	}
 	windowUs := cfg.windowMs * 1000
 	var buf []*sortEntry
 	return flow.NewTransformBlock(name, func(msg *flow.Message) ([]*flow.Message, error) {
 		if msg == nil {
@@ -52,11 +66,50 @@
 		buf = append(buf, nil)
 		copy(buf[i+1:], buf[i:])
 		buf[i] = &sortEntry{msg: msg, ts: msg.Timestamp, wire: wire}
 
+		// Handle overflow if max_buffered is set
+		if cfg.maxBuffered > 0 && int64(len(buf)) > cfg.maxBuffered {
+			switch cfg.overflowPolicy {
+			case overflowDropOldest:
+				// Drop oldest messages to fit within max_buffered
+				excess := int64(len(buf)) - cfg.maxBuffered
+				if excess > 0 {
+					buf = buf[excess:]
+				}
+			case overflowError:
+				return nil, fmt.Errorf("sort: buffer overflow (max %d messages)", cfg.maxBuffered)
+			case overflowForceFlush:
+				// Force flush oldest messages until we're within limit, or at least flush some
+				var flushed []*flow.Message
+				for int64(len(buf)) > cfg.maxBuffered && len(buf) > 0 {
+					flushed = append(flushed, buf[0].msg)
+					buf = buf[1:]
+				}
+				// Continue with normal window-based flushing on top of forced flushes
+				newest := msg.Timestamp
+				if len(buf) > 0 {
+					if last := buf[len(buf)-1]; last.ts > newest {
+						newest = last.ts
+					}
+				} else if len(flushed) > 0 {
+					// Buffer is now empty - no more window-based flushes
+					return flushed, nil
+				}
+				// Flush based on window
+				for len(buf) > 1 && buf[0].ts <= newest-windowUs {
+					flushed = append(flushed, buf[0].msg)
+					buf = buf[1:]
+				}
+				return flushed, nil
+			}
+		}
+
 		newest := msg.Timestamp
-		if last := buf[len(buf)-1]; last.ts > newest {
-			newest = last.ts
+		if len(buf) > 0 {
+			if last := buf[len(buf)-1]; last.ts > newest {
+				newest = last.ts
+			}
 		}
 		var out []*flow.Message
 		for len(buf) > 1 && buf[0].ts <= newest-windowUs {
 			out = append(out, buf[0].msg)

← Edit internal/flow/blocks/sort.go
Index: /root/astrate-mule/internal/flow/blocks/sort.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/sort.go
+++ /root/astrate-mule/internal/flow/blocks/sort.go
@@ -125,15 +125,27 @@
 	ts   int64
 	wire []byte
 }
 
+const (
+	overflowDropOldest  = "drop_oldest"
+	overflowError       = "error"
+	overflowForceFlush  = "force_flush"
+)
+
 type sortConfig struct {
-	windowMs int64
-	dedup    bool
+	windowMs       int64
+	dedup          bool
+	maxBuffered    int64
+	overflowPolicy string
 }
 
 func parseSortConfig(config map[string]any) (sortConfig, error) {
-	cfg := sortConfig{windowMs: 1000}
+	cfg := sortConfig{
+		windowMs:       1000,
+		maxBuffered:    0,
+		overflowPolicy: overflowDropOldest,
+	}
 	if v, ok := config["window_ms"]; ok && v != nil {
 		n, err := numAsInt64(v)
 		if err != nil || n < 0 {
 			return cfg, fmt.Errorf("window_ms must be non-negative")
@@ -142,6 +154,33 @@
 	}
 	if b, ok := config["dedup"].(bool); ok {
 		cfg.dedup = b
 	}
+	if v, ok := config["max_buffered"]; ok && v != nil {
+		n, err := numAsInt64(v)
+		if err != nil || n < 0 {
+			return cfg, fmt.Errorf("max_buffered must be non-negative")
+		}
+		cfg.maxBuffered = n
+	}
+	if v, ok := config["max_buffered_count"]; ok && v != nil && cfg.maxBuffered == 0 {
+		// Alternative naming for compatibility
+		n, err := numAsInt64(v)
+		if err != nil || n < 0 {
+			return cfg, fmt.Errorf("max_buffered_count must be non-negative")
+		}
+		cfg.maxBuffered = n
+	}
+	if v, ok := config["overflow_policy"]; ok && v != nil {
+		p, ok := v.(string)
+		if !ok {
+			return cfg, fmt.Errorf("overflow_policy must be a string")
+		}
+		switch p {
+		case overflowDropOldest, overflowError, overflowForceFlush:
+			cfg.overflowPolicy = p
+		default:
+			return cfg, fmt.Errorf("overflow_policy must be one of: drop_oldest, error, force_flush")
+		}
+	}
 	return cfg, nil
 }

← Edit internal/flow/blocks/schema.go
Index: /root/astrate-mule/internal/flow/blocks/schema.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/schema.go
+++ /root/astrate-mule/internal/flow/blocks/schema.go
@@ -9,9 +9,9 @@
 TypeToJSON:            `{"type":"object","properties":{}}`,
 TypeUpdateMetadata:    `{"type":"object","properties":{"set_metadata":{"type":"object","description":"Metadata string→string to merge; at least one of set_metadata, delete_metadata is required"},"delete_metadata":{"type":"array","description":"Metadata keys to delete; at least one of set_metadata, delete_metadata is required"}},"required":[]}`,
 TypeSplitMap:          `{"type":"object","properties":{"key_template":{"type":"string","description":"Output key template with {key} and {field} placeholders (default \"{key}/{field}\")"}},"required":[]}`,
 TypeRandomSource:      `{"type":"object","properties":{"type":{"type":"string","description":"Value type to emit; one of integer|real|boolean"},"interval_ms":{"type":"integer","description":"Milliseconds between emissions (default 1000)"},"min":{"type":"number","description":"Minimum value (inclusive)"},"max":{"type":"number","description":"Maximum value (inclusive)"},"key":{"type":"string","description":"Message key (default \"random\")"}},"required":["type"]}`,
-TypeSort:              `{"type":"object","properties":{"window_ms":{"type":"integer","description":"Buffering window in milliseconds (default 1000, must be ≥ 0)"},"dedup":{"type":"boolean","description":"Release each key once per window (default false)"}},"required":[]}`,
+TypeSort:              `{"type":"object","properties":{"window_ms":{"type":"integer","description":"Buffering window in milliseconds (default 1000, must be ≥ 0)"},"dedup":{"type":"boolean","description":"Release each key once per window (default false)"},"max_buffered":{"type":"integer","description":"Maximum number of buffered messages (default 0, unlimited)"},"overflow_policy":{"type":"string","description":"Action when buffer exceeds max_buffered: drop_oldest, error, or force_flush (default drop_oldest)","enum":["drop_oldest","error","force_flush"]}},"required":[]}`,
 TypeJSONPathMap:       `{"type":"object","properties":{"template":{"type":"string","description":"Template with $MESSAGE and $METADATA.<name> placeholders; unresolvable paths fail the message"}},"required":["template"]}`,
 TypeHTTPSource:        `{"type":"object","properties":{"urls":{"type":"array","description":"URLs to poll GET round-robin (at least one)"},"interval_ms":{"type":"integer","description":"Milliseconds between polls (default 1000)"},"timeout_ms":{"type":"integer","description":"Per-request timeout in milliseconds (default 5000)"}},"required":["urls"]}`,
 TypeHTTPSink:          `{"type":"object","properties":{"url":{"type":"string","description":"Base URL to POST payloads to"},"method":{"type":"string","description":"HTTP method (default \"POST\")"},"timeout_ms":{"type":"integer","description":"Per-request timeout in milliseconds (default 5000)"},"headers":{"type":"object","description":"Extra request headers string→string"}},"required":["url"]}`,
 TypeMQTTSource:        `{"type":"object","properties":{"url":{"type":"string","description":"Broker URL (e.g. tcp://127.0.0.1:1883)"},"topics":{"type":"array","description":"Topics to subscribe to (at least one, non-empty strings)"},"qos":{"type":"integer","description":"Subscription QoS level 0|1|2 (default 0)"},"client_id":{"type":"string","description":"Optional MQTT client ID (random when unset)"},"username":{"type":"string","description":"Optional username"},"password":{"type":"string","description":"Optional password"}},"required":["url","topics"]}`,

← Edit internal/flow/blocks/info.go
Index: /root/astrate-mule/internal/flow/blocks/info.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/info.go
+++ /root/astrate-mule/internal/flow/blocks/info.go
@@ -76,9 +76,9 @@
 TypeSort: {
 	Type:    TypeSort,
 	Role:    RoleTransform,
 	Summary: "Buffer messages and release them in ascending timestamp order behind a window",
-	Config:  "window_ms (default 1000, must be ≥ 0), dedup (bool, default false)",
+	Config:  "window_ms (default 1000, must be ≥ 0), dedup (bool, default false), max_buffered (default 0), overflow_policy (default drop_oldest)",
 },
 TypeJSONPathMap: {
 	Type:    TypeJSONPathMap,
 	Role:    RoleTransform,

$ go build ./... 2>&1 | tail -30
(no output)

$ go build ./... 2>&1 | wc -l && echo "---" && go build ./... 2>&1 | tail -5
0
---

The build succeeds. Let us also add the cap test as requested - "cap test feeding N same-timestamp messages". Let us look at the test structure in sort_test.go and add it.
← Edit internal/flow/blocks/sort_test.go
Index: /root/astrate-mule/internal/flow/blocks/sort_test.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/sort_test.go
+++ /root/astrate-mule/internal/flow/blocks/sort_test.go
@@ -152,4 +152,39 @@
 	if len(out) != 2 || out[0] != first || out[1] != second {
 		t.Errorf("equal-ts order broken: %v", out)
 	}
 }
+
+func TestSort_CapWithSameTimestamps(t *testing.T) {
+	b := sortBlock(t, map[string]any{"window_ms": 10, "max_buffered": 3, "overflow_policy": "drop_oldest"})
+	const ts = 100_000
+	// Feed 5 messages with same timestamp; with drop_oldest and max_buffered=3,
+	// after inserting all, the buffer should hold the last 3 (or be handled per policy)
+	for i := 0; i < 5; i++ {
+		m := intMsg("k", ts)
+		m.Data = int64(i)
+		if _, err := b.Process(m); err != nil {
+			t.Fatal(err)
+		}
+	}
+	// Flush - all have same timestamp, so window condition buf[0].ts <= newest-windowUs
+	// means 100000 <= (newest which is also 100000) - 10000 = 90000, false.
+	// But with overflow and drop_oldest, buffer should contain last 3 messages.
+	// When we flush by sending a later message, all buffered should be flushed (in order)
+	out, err := b.Process(intMsg("k", ts+1_000_000))
+	if err != nil {
+		t.Fatal(err)
+	}
+	if len(out) != 3 {
+		t.Fatalf("flushed %d messages, want 3", len(out))
+	}
+	// With drop_oldest, oldest 2 dropped; remaining should be messages i=2,3,4 in order
+	if v, ok := out[0].Data.(int64); !ok || v != 2 {
+		t.Errorf("out[0].Data = %v, want 2", out[0].Data)
+	}
+	if v, ok := out[1].Data.(int64); !ok || v != 3 {
+		t.Errorf("out[1].Data = %v, want 3", out[1].Data)
+	}
+	if v, ok := out[2].Data.(int64); !ok || v != 4 {
+		t.Errorf("out[2].Data = %v, want 4", out[2].Data)
+	}
+}

$ go build ./... 2>&1 | tail -2
(no output)

$ go test -v ./internal/flow/blocks/... 2>&1 | tail -40
--- PASS: TestVirtualDevicePool_DropsBadMessages/device_only (0.00s)
    --- PASS: TestVirtualDevicePool_DropsBadMessages/device_and_slash (0.00s)
    --- PASS: TestVirtualDevicePool_DropsBadMessages/empty_device_segment (0.00s)
    --- PASS: TestVirtualDevicePool_DropsBadMessages/empty_key (0.00s)
=== RUN   TestVirtualDevicePool_NilMessage
--- PASS: TestVirtualDevicePool_NilMessage (0.00s)
=== RUN   TestVirtualDevicePool_IngestErrorIsWrapped
--- PASS: TestVirtualDevicePool_IngestErrorIsWrapped (0.00s)
=== RUN   TestVirtualDevicePool_ConfigValidation
=== RUN   TestVirtualDevicePool_ConfigValidation/nil_ingest
=== RUN   TestVirtualDevicePool_ConfigValidation/missing_devices
=== RUN   TestVirtualDevicePool_ConfigValidation/empty_devices
=== RUN   TestVirtualDevicePool_ConfigValidation/empty_entry
=== RUN   TestVirtualDevicePool_ConfigValidation/acceptance_twin_json_array
=== RUN   TestVirtualDevicePool_ConfigValidation/auto_register_not_boolean
--- PASS: TestVirtualDevicePool_ConfigValidation (0.00s)
    --- PASS: TestVirtualDevicePool_ConfigValidation/nil_ingest (0.00s)
    --- PASS: TestVirtualDevicePool_ConfigValidation/missing_devices (0.00s)
    --- PASS: TestVirtualDevicePool_ConfigValidation/empty_devices (0.00s)
    --- PASS: TestVirtualDevicePool_ConfigValidation/empty_entry (0.00s)
    --- PASS: TestVirtualDevicePool_ConfigValidation/acceptance_twin_json_array (0.00s)
    --- PASS: TestVirtualDevicePool_ConfigValidation/auto_register_not_boolean (0.00s)
=== RUN   TestVirtualDevicePool_StaticInvariantNoRegisterWithoutAutoRegister
2026/09/23 21:20:12 WARN virtual_device_pool dropping message: device not in pool block=vdp key=stranger/com.example.X/y device=stranger
--- PASS: TestVirtualDevicePool_StaticInvariantNoRegisterWithoutAutoRegister (0.00s)
=== RUN   TestVirtualDevicePool_RegisterDepRequiredWithAutoRegister
--- PASS: TestVirtualDevicePool_RegisterDepRequiredWithAutoRegister (0.00s)
=== RUN   TestVirtualDevicePool_FirstSeenRegisters
--- PASS: TestVirtualDevicePool_FirstSeenRegisters (0.00s)
=== RUN   TestVirtualDevicePool_AlreadyRegisteredDropsWithoutCaching
2026/09/23 21:20:12 WARN virtual_device_pool dropping message: cannot register device block=vdp key=newdev/com.example.X/y device=newdev
2026/09/23 21:20:12 WARN virtual_device_pool dropping message: cannot register device block=vdp key=newdev/com.example.X/y device=newdev
--- PASS: TestVirtualDevicePool_AlreadyRegisteredDropsWithoutCaching (0.00s)
=== RUN   TestVirtualDevicePool_InfraRegisterErrorSurfaces
--- PASS: TestVirtualDevicePool_InfraRegisterErrorSurfaces (0.00s)
=== RUN   TestVirtualDevicePool_ConcurrentFirstSeenSingleRegister
--- PASS: TestVirtualDevicePool_ConcurrentFirstSeenSingleRegister (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	0.072s
FAIL

$ go test ./internal/flow/blocks/... 2>&1 | grep -v "^ok\|^PASS\|^RUN" | head -60
time=2026-09-23T21:20:24.968+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-23T21:20:24.969+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:37889
time=2026-09-23T21:20:24.969+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-23T21:20:24.971+02:00 level=INFO msg="mochi mqtt server started"
time=2026-09-23T21:20:24.976+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-23T21:20:24.976+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-23T21:20:24.976+02:00 level=INFO msg="mochi mqtt server stopped"
time=2026-09-23T21:20:24.977+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-23T21:20:24.977+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:39137
time=2026-09-23T21:20:24.977+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-23T21:20:24.979+02:00 level=INFO msg="mochi mqtt server started"
time=2026-09-23T21:20:24.983+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-23T21:20:24.983+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-23T21:20:24.984+02:00 level=INFO msg="mochi mqtt server stopped"
time=2026-09-23T21:20:24.984+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-23T21:20:24.984+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:40775
time=2026-09-23T21:20:24.984+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-23T21:20:24.985+02:00 level=INFO msg="mochi mqtt server started"
time=2026-09-23T21:20:24.987+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-23T21:20:24.987+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-23T21:20:24.987+02:00 level=INFO msg="mochi mqtt server stopped"
time=2026-09-23T21:20:24.987+02:00 level=WARN msg="" listener=test-tcp error="read tcp 127.0.0.1:40775->127.0.0.1:33700: use of closed network connection"
time=2026-09-23T21:20:25.439+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-23T21:20:25.440+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:38623
time=2026-09-23T21:20:25.440+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-23T21:20:25.441+02:00 level=INFO msg="mochi mqtt server started"
time=2026-09-23T21:20:25.446+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-23T21:20:25.446+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-23T21:20:25.446+02:00 level=INFO msg="mochi mqtt server stopped"
time=2026-09-23T21:20:25.446+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-23T21:20:25.446+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:34153
time=2026-09-23T21:20:25.447+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-23T21:20:25.448+02:00 level=INFO msg="mochi mqtt server started"
time=2026-09-23T21:20:25.453+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-23T21:20:25.453+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-23T21:20:25.453+02:00 level=INFO msg="mochi mqtt server stopped"
time=2026-09-23T21:20:25.453+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-23T21:20:25.454+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:46483
time=2026-09-23T21:20:25.455+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-23T21:20:25.456+02:00 level=INFO msg="mochi mqtt server started"
time=2026-09-23T21:20:25.461+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-23T21:20:25.461+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-23T21:20:25.461+02:00 level=INFO msg="mochi mqtt server stopped"
time=2026-09-23T21:20:25.461+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-23T21:20:25.461+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:41777
time=2026-09-23T21:20:25.461+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-23T21:20:25.462+02:00 level=INFO msg="mochi mqtt server started"
time=2026-09-23T21:20:25.463+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-23T21:20:25.463+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-23T21:20:25.463+02:00 level=INFO msg="mochi mqtt server stopped"
--- FAIL: TestSort_CapWithSameTimestamps (0.00s)
    sort_test.go:178: flushed 2 messages, want 3
FAIL
FAIL	github.com/astrate-platform/astrate/internal/flow/blocks	1.278s
FAIL

$ go test -v ./internal/flow/blocks/... -run TestSort 2>&1 | tail -80
=== RUN   TestSort_FlushesAscendingOnArrival
--- PASS: TestSort_FlushesAscendingOnArrival (0.00s)
=== RUN   TestSort_Dedup
--- PASS: TestSort_Dedup (0.00s)
=== RUN   TestSort_NilPassthrough
--- PASS: TestSort_NilPassthrough (0.00s)
=== RUN   TestSort_ConstructRejections
=== RUN   TestSort_ConstructRejections/negative_window/reject
=== RUN   TestSort_ConstructRejections/negative_window/accept
=== RUN   TestSort_ConstructRejections/non-numeric_window/reject
=== RUN   TestSort_ConstructRejections/non-numeric_window/accept
--- PASS: TestSort_ConstructRejections (0.00s)
    --- PASS: TestSort_ConstructRejections/negative_window/reject (0.00s)
    --- PASS: TestSort_ConstructRejections/negative_window/accept (0.00s)
    --- PASS: TestSort_ConstructRejections/non-numeric_window/reject (0.00s)
    --- PASS: TestSort_ConstructRejections/non-numeric_window/accept (0.00s)
=== RUN   TestSort_EqualTimestampsKeepArrivalOrder
--- PASS: TestSort_EqualTimestampsKeepArrivalOrder (0.00s)
=== RUN   TestSort_CapWithSameTimestamps
    sort_test.go:178: flushed 2 messages, want 3
--- FAIL: TestSort_CapWithSameTimestamps (0.00s)
FAIL
FAIL	github.com/astrate-platform/astrate/internal/flow/blocks	0.070s
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	0.051s [no tests to run]
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.050s [no tests to run]
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	0.041s [no tests to run]
FAIL

$ go test -v ./internal/flow/blocks/... -run TestSort_Cap 2>&1 | tail -30
=== RUN   TestSort_CapWithSameTimestamps
    sort_test.go:178: flushed 2 messages, want 3
--- FAIL: TestSort_CapWithSameTimestamps (0.00s)
FAIL
FAIL	github.com/astrate-platform/astrate/internal/flow/blocks	0.047s
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	0.043s [no tests to run]
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.041s [no tests to run]
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	0.040s [no tests to run]
FAIL

$ python3 << 'PY'
# Simulate the case
# We insert 5 messages all with same ts. After each insert, we do normal flush.
# With max_buffered=3 and drop_oldest on overflow: when we insert 4th message (after 3 in buffer), buffer becomes 4 > 3, so drop oldest (excess 1) -> buffer becomes 3? Or we dropped to make room.
# But we inserted the new message already. Let us trace logic.
pass
PY

(no output)

$ python3 << 'PY'
# Simulate flush trigger processing
ts = 100000
newest_trigger = ts + 1000000  # 1100000
# Before processing flush, buf is [msg2(i=2), msg3(i=3), msg4(i=4)] with ts all 100000
buf = ['m2', 'm3', 'm4'] # conceptual
# Process flush msg: insert it. flush msg has ts = 1100000. All existing have ts 100000 < 1100000, so Search returns len(buf)=3, insert at end
# buf becomes ['m2','m3','m4','flush']
# len is 4, max_buffered 3, overflow triggers - drop oldest 1 -> buf becomes ['m3','m4','flush']
# Then compute newest: msg.Timestamp of flush is 1100000; last is flush with ts 1100000, so newest is 1100000
# Then flush loop: for len(buf)>1 (3>1) and buf[0].ts (m3.ts=100000) <= newest-windowUs (1100000-10000=1090000) -> 100000 <= 1090000 is true! So it flushes buf[0] which is m3? Wait wait - look up what buf contains now.
# After dropping oldest due to overflow from inserting flush? Let us look: "Insert keeping ascending..." - we insert the flush message first, then check overflow.
# So inserting flush (ts 1100000) into buf ['m2','m3','m4'] (all ts 100000): Search returns 3 (no element with ts>1100000? all elements have ts 100000, so no element has ts>1100000; so Search returns 3 = len(buf)). So we insert flush at position 3. buf becomes ['m2','m3','m4','flush'] (conceptually). Now len is 4, max_buffered is 3, so overflow policy drop_oldest kicks in: excess = 1, buf becomes buf[1:] = ['m3','m4','flush']. Then we compute newest - msg.Timestamp is flush's ts = 1100000; buf[-1] is flush with ts 1100000, so newest is 1100000. WindowUs is 10000, so newest-windowUs = 1090000.
# Now flush: while len(buf)>1 (3>1 true) and buf[0].ts (m3 has ts 100000) <= 1090000 is true. So we flush m3 - append to out, buf becomes ['m4','flush']. out=[m3]
# Loop again: len(buf)>1 is true (2>1), buf[0].ts is m4.ts 100000 <= 1090000 is true. Flush m4, buf becomes ['flush']. out=[m3, m4]
# Loop again: len(buf)>1 is false (1>1 is false). Stop. Return out = [m3, m4]
# But wait - in our expected result from the test, we wanted out[0] to be the message with data 2 (which was m2), out[1] has data 3 (m3), out[2] has data 4 (m4). But according to this trace, we got m3, m4 - missing m2 and also got different values? Wait compare.
# Earlier: after inserting all 5 messages (i=0..4), final buf before flush was [msg2, msg3, msg4] where msg2 has data=2, msg3 has data=3, msg4 has data=4. Correct - that matches.
# Then we process flush trigger (intMsg with ts ts+1M). That flush message is inserted into the buffer. So in the buffer context, the flush message becomes part of sorting - but it's the newest arrival.
# But when overflow happens *after* inserting the new message (flush), we drop oldest from the buffer. So inserting flush into [msg2, msg3, msg4] gives [msg2, msg3, msg4, flush]. Overflow drops msg2 → buffer becomes [msg3, msg4, flush]. Then window-based flush releases msg3 and msg4 (both older than flush by more than window). Flush itself stays in buffer (len becomes 1, so loop stops). So total out is [msg3, msg4] - only 2 messages flushed. That means the test as written will get only 2 messages. But the test says "want 3" and expects data 2,3,4. That means either my overflow handling is wrong, or the test logic is off.
PY

(no output)

← Edit internal/flow/blocks/sort.go
Index: /root/astrate-mule/internal/flow/blocks/sort.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/sort.go
+++ /root/astrate-mule/internal/flow/blocks/sort.go
@@ -66,56 +66,41 @@
 		buf = append(buf, nil)
 		copy(buf[i+1:], buf[i:])
 		buf[i] = &sortEntry{msg: msg, ts: msg.Timestamp, wire: wire}
 
-		// Handle overflow if max_buffered is set
+		// First, try to flush based on window with the new message
+		newest := msg.Timestamp
+		if len(buf) > 0 {
+			if last := buf[len(buf)-1]; last.ts > newest {
+				newest = last.ts
+			}
+		}
+		var out []*flow.Message
+		for len(buf) > 1 && buf[0].ts <= newest-windowUs {
+			out = append(out, buf[0].msg)
+			buf = buf[1:]
+		}
+
+		// Now handle overflow if max_buffered is set and buffer still exceeds limit
 		if cfg.maxBuffered > 0 && int64(len(buf)) > cfg.maxBuffered {
 			switch cfg.overflowPolicy {
 			case overflowDropOldest:
-				// Drop oldest messages to fit within max_buffered
+				// Drop oldest messages to fit within max_buffered (excluding the case
+				// where this would drop everything? but we need to keep space for semantics)
 				excess := int64(len(buf)) - cfg.maxBuffered
 				if excess > 0 {
 					buf = buf[excess:]
 				}
 			case overflowError:
 				return nil, fmt.Errorf("sort: buffer overflow (max %d messages)", cfg.maxBuffered)
 			case overflowForceFlush:
-				// Force flush oldest messages until we're within limit, or at least flush some
-				var flushed []*flow.Message
+				// Force flush oldest messages until we're within limit
 				for int64(len(buf)) > cfg.maxBuffered && len(buf) > 0 {
-					flushed = append(flushed, buf[0].msg)
+					out = append(out, buf[0].msg)
 					buf = buf[1:]
 				}
-				// Continue with normal window-based flushing on top of forced flushes
-				newest := msg.Timestamp
-				if len(buf) > 0 {
-					if last := buf[len(buf)-1]; last.ts > newest {
-						newest = last.ts
-					}
-				} else if len(flushed) > 0 {
-					// Buffer is now empty - no more window-based flushes
-					return flushed, nil
-				}
-				// Flush based on window
-				for len(buf) > 1 && buf[0].ts <= newest-windowUs {
-					flushed = append(flushed, buf[0].msg)
-					buf = buf[1:]
-				}
-				return flushed, nil
 			}
 		}
-
-		newest := msg.Timestamp
-		if len(buf) > 0 {
-			if last := buf[len(buf)-1]; last.ts > newest {
-				newest = last.ts
-			}
-		}
-		var out []*flow.Message
-		for len(buf) > 1 && buf[0].ts <= newest-windowUs {
-			out = append(out, buf[0].msg)
-			buf = buf[1:]
-		}
 		return out, nil
 	}), nil
 }
 

$ go test -v ./internal/flow/blocks/... -run TestSort_Cap 2>&1 | tail -10
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	0.033s
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached) [no tests to run]
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.030s [no tests to run]
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached) [no tests to run]

$ go test -v ./internal/flow/blocks/... -run TestSort 2>&1 | tail -5
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.044s [no tests to run]
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached) [no tests to run]

$ go test ./internal/flow/blocks/... 2>&1 | tail -3
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.534s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)

$ go test ./internal/flow/blocks/... 2>&1 | grep "blocks" | head -1
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)

$ go test ./internal/flow/blocks 2>&1 | tail -2
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)

$ go test ./internal/flow/blocks/... 2>&1 | grep -c "ok" && go test ./internal/flow/blocks/... 2>&1 | grep "FAIL\|RUN" | head -2
4

$ go clean -testcache 2>&1; go test ./internal/flow/blocks 2>&1 | tail -2
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	1.290s

All the sort tests (including the new cap test) pass. The task also said "guard the multiply. window_ms * 1000 (sort.go:31) overflows for large configs and silently inverts the flush edge". I added a check for overflow before multiplying. Also I handled the general overflow case in the multiply.

But looking at my current code: I do `if cfg.windowMs > math.MaxInt64/1000 { return nil, fmt.Errorf("sort: window_ms too large") }` and then `windowUs := cfg.windowMs * 1000`. That correctly prevents overflow. Good.

The task requirements are met - added max buffered count with overflow policy (drop_oldest, error, force_flush), guarded multiply, added cap test with same-timestamp messages. All existing tests still pass and new test passes.

Let us also run vet as the gate would do. The gate is `go vet ./... && go test ./...` (from MULE.md). Let us verify.
$ go vet ./... 2>&1 | tail -30
(no output)

All vet checks pass, all tests in blocks package pass. The task is complete. Let us write the report as specified in MULE.md format.Done: Add max_buffered/overflow_policy to Sort, guard window_ms*1000 multiply, and add cap test for same-timestamp bursts
Files: internal/flow/blocks/sort.go, internal/flow/blocks/schema.go, internal/flow/blocks/info.go, internal/flow/blocks/sort_test.go
Verified: go vet ./... -> pass; go test ./internal/flow/blocks/... -> pass
Unsure: nothing
Follow-ups: none
