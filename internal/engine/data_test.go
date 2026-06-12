package engine

import (
	"context"
	"strings"
	"testing"
	"time"

	promtest "github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// pipelineRig is a synchronous harness: an unstarted single-shard engine
// whose handle method the tests drive directly.
type pipelineRig struct {
	e  *Engine
	sh *shard
}

func newPipelineRig(t testing.TB, cfg Config) (*pipelineRig, *fakeStore) {
	t.Helper()
	fs := newFakeStore()
	seedAlpha(t, fs)
	cfg.Shards = 1
	e := newTestEngine(t, fs, cfg)
	return &pipelineRig{e: e, sh: e.shards[0]}, fs
}

// handle drives one message through the pipeline synchronously.
func (r *pipelineRig) handle(m broker.InboundMessage) {
	r.e.handle(context.Background(), r.sh, m)
}

// hugeArrayJSON builds {"v":[0,0,...]} with n elements.
func hugeArrayJSON(n int) []byte {
	return []byte(`{"v":[` + strings.TrimSuffix(strings.Repeat("0,", n), ",") + `]}`)
}

// TestRejectReasons is the docs/ROADMAP.md §7.3 reject-reason table: one row
// per §2.6 failure mode, driven with the M1 fixtures. Every rejected message
// is acknowledged (consumed), counted under exactly its reason, and never
// reaches the batch.
func TestRejectReasons(t *testing.T) {
	ts := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	ghost := deviceid.ID{0xff, 0xee}

	cases := []struct {
		name   string
		cfg    Config
		msg    func(t *testing.T, ack *ackCounter) broker.InboundMessage
		reason string
	}{
		{
			name: "realm unknown",
			msg: func(t *testing.T, ack *ackCounter) broker.InboundMessage {
				m := deviceMsg("com.astrate.test.Minimal", "/value", 1, enc(t, 1.0, nil, payload.FormatBSON), ack)
				m.Realm = "ghost"
				m.Topic = "ghost/" + devAlpha.String() + "/com.astrate.test.Minimal/value"
				return m
			},
			reason: reasonRealmUnknown,
		},
		{
			name: "malformed topic",
			msg: func(t *testing.T, ack *ackCounter) broker.InboundMessage {
				m := deviceMsg("com.astrate.test.Minimal", "/value", 1, enc(t, 1.0, nil, payload.FormatBSON), ack)
				m.Topic = "alpha/garbage"
				return m
			},
			reason: reasonMalformedTopic,
		},
		{
			name: "device unknown",
			msg: func(t *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsgFor(ghost, "com.astrate.test.Minimal", "/value", 1, enc(t, 1.0, nil, payload.FormatBSON), ack)
			},
			reason: reasonDeviceUnknown,
		},
		{
			name: "interface not in introspection",
			msg: func(t *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.undeclared.Iface", "/x", 1, enc(t, 1.0, nil, payload.FormatBSON), ack)
			},
			reason: reasonInterfaceNotDeclared,
		},
		{
			name: "interface declared but not installed",
			msg: func(t *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.astrate.test.DeclaredButMissing", "/x", 1, enc(t, 1.0, nil, payload.FormatBSON), ack)
			},
			reason: reasonInterfaceNotInstalled,
		},
		{
			name: "ownership violation",
			msg: func(_ *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.astrate.test.ServerProperties", "/limits/maxConnections", 2,
					[]byte(`{"v":5}`), ack)
			},
			reason: reasonOwnershipViolation,
		},
		{
			name: "unexpected path individual",
			msg: func(t *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.astrate.test.AllScalarTypes", "/nope", 1, enc(t, 1.0, nil, payload.FormatBSON), ack)
			},
			reason: reasonUnexpectedPath,
		},
		{
			name: "unexpected path object depth",
			msg: func(_ *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("org.astarte-platform.genericsensors.Geolocation", "/a/b", 2,
					[]byte(`{"v":{"latitude":1.0},"t":"2026-06-01T12:00:00.000Z"}`), ack)
			},
			reason: reasonUnexpectedPath,
		},
		{
			name: "unexpected path object missing prefix",
			msg: func(_ *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("org.astarte-platform.genericsensors.Geolocation", "", 2,
					[]byte(`{"v":{"latitude":1.0},"t":"2026-06-01T12:00:00.000Z"}`), ack)
			},
			reason: reasonUnexpectedPath,
		},
		{
			name: "payload too large",
			cfg:  Config{MaxPayloadBytes: 16},
			msg: func(_ *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.astrate.test.AllScalarTypes", "/string", 1,
					[]byte(`{"v":"`+strings.Repeat("x", 64)+`"}`), ack)
			},
			reason: payload.ReasonTooLarge.String(),
		},
		{
			name: "unknown format",
			msg: func(_ *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.astrate.test.AllScalarTypes", "/double", 1, []byte("garbage"), ack)
			},
			reason: payload.ReasonUnknownFormat.String(),
		},
		{
			name: "malformed document",
			msg: func(_ *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.astrate.test.AllScalarTypes", "/double", 1, []byte("{nope"), ack)
			},
			reason: payload.ReasonMalformed.String(),
		},
		{
			name: "no value member",
			msg: func(_ *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.astrate.test.AllScalarTypes", "/double", 1, []byte("{}"), ack)
			},
			reason: payload.ReasonNoValue.String(),
		},
		{
			name: "missing explicit timestamp",
			msg: func(t *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.astrate.test.AllScalarTypes", "/datetime", 1,
					enc(t, ts, nil, payload.FormatBSON), ack)
			},
			reason: payload.ReasonBadTimestamp.String(),
		},
		{
			name: "type mismatch",
			msg: func(_ *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.astrate.test.AllScalarTypes", "/double", 1, []byte(`{"v":"nope"}`), ack)
			},
			reason: payload.ReasonTypeMismatch.String(),
		},
		{
			name: "value too large",
			msg: func(_ *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.astrate.test.PropertyArrays", "/config/thresholds", 2,
					hugeArrayJSON(payload.MaxArrayLen+1), ack)
			},
			reason: payload.ReasonValueTooLarge.String(),
		},
		{
			name: "bad object key",
			msg: func(_ *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.astrate.test.ObjectFlat", "", 1,
					[]byte(`{"v":{"latitude":1.0,"bogus":2.0},"t":"2026-06-01T12:00:00.000Z"}`), ack)
			},
			reason: payload.ReasonBadObject.String(),
		},
		{
			name: "unset on datastream",
			msg: func(_ *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.astrate.test.AllScalarTypes", "/double", 1, nil, ack)
			},
			reason: payload.ReasonUnsetNotAllowed.String(),
		},
		{
			name: "unset without allow_unset",
			msg: func(_ *testing.T, ack *ackCounter) broker.InboundMessage {
				return deviceMsg("com.astrate.test.PropertyArrays", "/config/labels", 2, nil, ack)
			},
			reason: payload.ReasonUnsetNotAllowed.String(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rig, _ := newPipelineRig(t, tc.cfg)
			ack := &ackCounter{}
			rig.handle(tc.msg(t, ack))

			if got := promtest.ToFloat64(rig.e.met.rejects.WithLabelValues(tc.reason)); got != 1 {
				t.Errorf("reject counter[%s] = %v, want 1", tc.reason, got)
			}
			if !ack.acked() {
				t.Error("rejected message was not acknowledged (device would stall)")
			}
			if rig.sh.batch.size() != 0 {
				t.Errorf("rejected message reached the batch: %d ops", rig.sh.batch.size())
			}
		})
	}
}

// TestRejectSeam: rejections feed the M6b device_error seam.
func TestRejectSeam(t *testing.T) {
	rig, _ := newPipelineRig(t, Config{})
	var gotReason, gotDetail string
	rig.e.onDeviceError = func(_ broker.InboundMessage, reason, detail string) {
		gotReason, gotDetail = reason, detail
	}
	ack := &ackCounter{}
	rig.handle(deviceMsg("com.undeclared.Iface", "/x", 1, enc(t, 1.0, nil, payload.FormatBSON), ack))
	if gotReason != reasonInterfaceNotDeclared || gotDetail == "" {
		t.Errorf("device_error seam saw (%q, %q)", gotReason, gotDetail)
	}
}

// TestAcceptedOps covers the happy paths of the §2.6 pipeline: each message
// class emits the right PersistOp.
func TestAcceptedOps(t *testing.T) {
	rig, _ := newPipelineRig(t, Config{BatchMaxRows: 1000})
	explicit := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	ack := &ackCounter{}

	// 1. Individual BSON double, server-side timestamp.
	m := deviceMsg("com.astrate.test.AllScalarTypes", "/double", 0, enc(t, 22.5, nil, payload.FormatBSON), ack)
	rig.handle(m)
	// 2. Individual BSON datetime with explicit timestamp.
	rig.handle(deviceMsg("com.astrate.test.AllScalarTypes", "/datetime", 1,
		enc(t, explicit, &explicit, payload.FormatBSON), ack))
	// 3. Individual JSON integer.
	rig.handle(deviceMsg("com.astrate.test.AllScalarTypes", "/integer", 1, []byte(`{"v":7}`), ack))
	// 4. Flat object aggregation (empty path prefix).
	objVal := map[string]payload.Value{"latitude": 45.07, "longitude": 7.69}
	rig.handle(deviceMsg("com.astrate.test.ObjectFlat", "", 1, enc(t, objVal, &explicit, payload.FormatBSON), ack))
	// 5. Parametric object aggregation on "/gps".
	rig.handle(deviceMsg("org.astarte-platform.genericsensors.Geolocation", "/gps", 2,
		enc(t, objVal, &explicit, payload.FormatBSON), ack))
	// 6. Property set (device-owned).
	rig.handle(deviceMsg("com.astrate.test.PropertyArrays", "/config/thresholds", 2,
		[]byte(`{"v":[1.5,2.5]}`), ack))
	// 7. Property unset (allow_unset).
	rig.handle(deviceMsg("com.astrate.test.PropertyArrays", "/config/thresholds", 2, nil, ack))

	ops := rig.sh.batch.ops
	if len(ops) != 7 {
		t.Fatalf("%d ops emitted, want 7", len(ops))
	}
	check := func(i int, kind OpKind, path string, format payload.Format) *PersistOp {
		t.Helper()
		op := &ops[i]
		if op.Kind != kind || op.Path != path || op.Format != format {
			t.Errorf("op %d = kind %s path %q format %s, want %s %q %s",
				i, op.Kind, op.Path, op.Format, kind, path, format)
		}
		if op.Realm != realmAlpha || op.RealmID != realmAlphaID || op.DeviceID != devAlpha {
			t.Errorf("op %d addressing: %+v", i, op)
		}
		return op
	}

	op := check(0, OpIndividual, "/double", payload.FormatBSON)
	if v, ok := op.Value.(float64); !ok || v != 22.5 {
		t.Errorf("op 0 value %v", op.Value)
	}
	if !op.TS.Equal(m.ReceivedAt) {
		t.Errorf("op 0 ts %s, want reception %s", op.TS, m.ReceivedAt)
	}
	if op.Interface.ID != ifaceAllScalars || op.Mapping == nil || op.Mapping.EndpointID != ifaceAllScalars*1000 {
		t.Errorf("op 0 schema refs: iface %d endpoint %+v", op.Interface.ID, op.Mapping)
	}

	op = check(1, OpIndividual, "/datetime", payload.FormatBSON)
	if !op.TS.Equal(explicit) {
		t.Errorf("op 1 ts %s, want explicit %s", op.TS, explicit)
	}

	op = check(2, OpIndividual, "/integer", payload.FormatJSON)
	if v, ok := op.Value.(int32); !ok || v != 7 {
		t.Errorf("op 2 value %v (%T)", op.Value, op.Value)
	}

	op = check(3, OpObject, "", payload.FormatBSON)
	if v, ok := op.Value.(map[string]payload.Value); !ok || v["latitude"] != 45.07 {
		t.Errorf("op 3 value %v", op.Value)
	}
	if op.Mapping != nil {
		t.Error("op 3 (object) carries an endpoint mapping")
	}
	if !op.TS.Equal(explicit) {
		t.Errorf("op 3 ts %s, want explicit %s", op.TS, explicit)
	}

	check(4, OpObject, "/gps", payload.FormatBSON)

	op = check(5, OpPropertySet, "/config/thresholds", payload.FormatJSON)
	if v, ok := op.Value.([]float64); !ok || len(v) != 2 || v[0] != 1.5 {
		t.Errorf("op 5 value %v (%T)", op.Value, op.Value)
	}

	op = check(6, OpPropertyUnset, "/config/thresholds", payload.FormatEmpty)
	if op.Value != nil {
		t.Errorf("op 6 (unset) value %v", op.Value)
	}
}

// TestFormatHintFlip covers the docs/DESIGN.md §3.5.4 sticky hint:
// bson→json on the first JSON payload, json→bson only after an armed
// emptyCache reset followed by a BSON payload.
func TestFormatHintFlip(t *testing.T) {
	ctx := context.Background()
	rig, fs := newPipelineRig(t, Config{BatchMaxRows: 1000})
	ack := &ackCounter{}
	jsonMsg := func() broker.InboundMessage {
		return deviceMsg("com.astrate.test.AllScalarTypes", "/integer", 0, []byte(`{"v":1}`), ack)
	}
	bsonMsg := func() broker.InboundMessage {
		return deviceMsg("com.astrate.test.AllScalarTypes", "/double", 0, enc(t, 1.0, nil, payload.FormatBSON), ack)
	}

	// BSON traffic on a bson-hinted device: no flip, nothing persisted.
	rig.handle(bsonMsg())
	if h := fs.hintFor(realmAlpha, devAlpha); h != "" {
		t.Errorf("hint persisted (%q) without a flip", h)
	}

	// First JSON payload flips to json.
	rig.handle(jsonMsg())
	if h := fs.hintFor(realmAlpha, devAlpha); h != hintJSON {
		t.Errorf("hint after first JSON payload = %q, want json", h)
	}
	dev, err := rig.e.devices.get(ctx, realmAlpha, realmAlphaID, devAlpha)
	if err != nil {
		t.Fatalf("device state: %v", err)
	}
	if dev.hint() != hintJSON {
		t.Errorf("cached hint = %q, want json", dev.hint())
	}

	// BSON without an armed reset: sticky json.
	rig.handle(bsonMsg())
	if dev.hint() != hintJSON {
		t.Error("hint flipped back to bson without an emptyCache reset")
	}

	// Armed reset followed by JSON: stays json, disarms.
	dev.mu.Lock()
	dev.resetHintOnBSON = true
	dev.mu.Unlock()
	rig.handle(jsonMsg())
	dev.mu.Lock()
	armed := dev.resetHintOnBSON
	dev.mu.Unlock()
	if armed || dev.hint() != hintJSON {
		t.Errorf("armed reset + JSON payload: hint %q armed %v, want json/disarmed", dev.hint(), armed)
	}

	// Armed reset followed by BSON: flips back to bson.
	dev.mu.Lock()
	dev.resetHintOnBSON = true
	dev.mu.Unlock()
	rig.handle(bsonMsg())
	if dev.hint() != hintBSON || fs.hintFor(realmAlpha, devAlpha) != hintBSON {
		t.Errorf("armed reset + BSON payload: cached %q persisted %q, want bson",
			dev.hint(), fs.hintFor(realmAlpha, devAlpha))
	}
}

// TestUnhandledDispatch: introspection and control messages are consumed
// (acked + counted) while the M6b handlers are not wired, and routed to the
// seams once they are.
func TestUnhandledDispatch(t *testing.T) {
	rig, _ := newPipelineRig(t, Config{})

	ack := &ackCounter{}
	rig.handle(deviceMsg("", "", 2, []byte("com.astrate.test.Minimal:0:1"), ack))
	if !ack.acked() {
		t.Error("unhandled introspection not acknowledged")
	}
	if got := promtest.ToFloat64(rig.e.met.unhandled.WithLabelValues("introspection")); got != 1 {
		t.Errorf("unhandled[introspection] = %v, want 1", got)
	}

	ack = &ackCounter{}
	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
	if !ack.acked() {
		t.Error("unhandled control not acknowledged")
	}
	if got := promtest.ToFloat64(rig.e.met.unhandled.WithLabelValues("control")); got != 1 {
		t.Errorf("unhandled[control] = %v, want 1", got)
	}

	// Wired seams take over (and own the acknowledgment).
	var introBody string
	rig.e.onIntrospection = func(_ context.Context, m broker.InboundMessage) { introBody = string(m.Payload) }
	var controlSub string
	rig.e.onControl = func(_ context.Context, _ broker.InboundMessage, subpath string) { controlSub = subpath }

	ack = &ackCounter{}
	rig.handle(deviceMsg("", "", 2, []byte("a:1:0;b:0:3"), ack))
	if introBody != "a:1:0;b:0:3" {
		t.Errorf("introspection seam saw %q", introBody)
	}
	if ack.acked() {
		t.Error("dispatcher acked a message owned by the introspection handler")
	}
	rig.handle(deviceMsg("control", "/producer/properties", 2, []byte{0, 0, 0, 0}, ack))
	if controlSub != "producer/properties" {
		t.Errorf("control seam saw subpath %q", controlSub)
	}
}

// TestDeviceLoadParking: transient store failures while loading device
// state park the shard (with backoff) instead of rejecting, then proceed.
func TestDeviceLoadParking(t *testing.T) {
	rig, fs := newPipelineRig(t, Config{BatchMaxRows: 1000})
	fs.mu.Lock()
	fs.getDeviceErrs = []error{errTransient}
	fs.mu.Unlock()

	ack := &ackCounter{}
	startAt := time.Now()
	rig.handle(deviceMsg("com.astrate.test.Minimal", "/value", 1, enc(t, 4.5, nil, payload.FormatBSON), ack))
	if elapsed := time.Since(startAt); elapsed < parkBackoffStart {
		t.Errorf("handle returned after %s, want >= %s of parking", elapsed, parkBackoffStart)
	}
	if rig.sh.batch.size() != 1 {
		t.Fatalf("message not enqueued after transient failure: %d ops", rig.sh.batch.size())
	}
	if got := promtest.ToFloat64(rig.e.met.internalErrors); got != 1 {
		t.Errorf("internal error counter = %v, want 1", got)
	}
}
