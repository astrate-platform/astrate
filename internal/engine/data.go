package engine

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// Engine-level reject reasons, complementing the payload.RejectReason set
// (docs/DESIGN.md §2.6: "failures are never silent"). The strings are the
// metric labels and the M6b device_error trigger error names.
const (
	// reasonRealmUnknown: the message realm is not in the schema snapshot.
	reasonRealmUnknown = "realm_unknown"
	// reasonDeviceUnknown: the broker delivered a message for a device the
	// store does not know (cannot normally happen).
	reasonDeviceUnknown = "device_unknown"
	// reasonMalformedTopic: the topic does not carry the expected
	// "<realm>/<device_id>" root (cannot normally happen: the ACL pins it).
	reasonMalformedTopic = "malformed_topic"
	// reasonInterfaceNotDeclared: §2.6 step 1 — no introspected interface
	// name prefixes the topic.
	reasonInterfaceNotDeclared = "interface_not_in_introspection"
	// reasonInterfaceNotInstalled: the introspection declares a name:major
	// the realm has no compiled interface for.
	reasonInterfaceNotInstalled = "interface_not_installed"
	// reasonOwnershipViolation: §2.6 step 3 — a device published on a
	// server-owned interface.
	reasonOwnershipViolation = "ownership_violation"
	// reasonUnexpectedPath: §2.6 step 2 — the path resolves no endpoint
	// mapping (or no valid object aggregation prefix).
	reasonUnexpectedPath = "unexpected_path"
	// reasonIntrospectionInvalid: the introspection payload is oversized or
	// malformed (docs/ROADMAP.md §7.2 file 6.7).
	reasonIntrospectionInvalid = "introspection_invalid"
	// reasonControlUnknown: a control publish on an unknown subpath.
	reasonControlUnknown = "control_unknown"
	// reasonControlInvalid: a control payload with a bad zlib frame, a
	// declared size above the ceiling, or a lying size header (§4.5).
	reasonControlInvalid = "control_payload_invalid"
)

// engineRejectReasons pre-registers the labels above (newMetrics).
var engineRejectReasons = []string{
	reasonRealmUnknown, reasonDeviceUnknown, reasonMalformedTopic,
	reasonInterfaceNotDeclared, reasonInterfaceNotInstalled,
	reasonOwnershipViolation, reasonUnexpectedPath,
	reasonIntrospectionInvalid, reasonControlUnknown, reasonControlInvalid,
}

// Shard-parking backoff bounds for transient store failures
// (docs/DESIGN.md §5.3 "DB outage: shards park with exponential backoff").
const (
	parkBackoffStart = 100 * time.Millisecond
	parkBackoffCap   = 5 * time.Second
)

// OpKind discriminates PersistOp values.
type OpKind uint8

// OpKind values.
const (
	// OpIndividual is one individual-datastream insert.
	OpIndividual OpKind = iota + 1
	// OpObject is one object-datastream insert.
	OpObject
	// OpPropertySet is a property upsert.
	OpPropertySet
	// OpPropertyUnset is a property delete (empty payload, allow_unset).
	OpPropertyUnset
)

// opKinds pre-registers the kind labels (newMetrics).
var opKinds = []OpKind{OpIndividual, OpObject, OpPropertySet, OpPropertyUnset}

// String returns the stable snake_case metrics label.
func (k OpKind) String() string {
	switch k {
	case OpIndividual:
		return "individual"
	case OpObject:
		return "object"
	case OpPropertySet:
		return "property_set"
	case OpPropertyUnset:
		return "property_unset"
	default:
		return fmt.Sprintf("OpKind(%d)", uint8(k))
	}
}

// PersistOp is one validated operation emitted by the §2.6 pipeline
// (docs/ROADMAP.md §7.1 file 6.4), carried through the shard micro-batch and
// handed to the afterCommit observers (triggers and live fan-out, M6b).
type PersistOp struct {
	// Kind discriminates the persistence action.
	Kind OpKind
	// Realm and RealmID identify the tenant.
	Realm   string
	RealmID int16
	// DeviceID is the publishing device.
	DeviceID deviceid.ID
	// Interface is the compiled interface the message validated against.
	Interface *interfaceschema.CompiledInterface
	// Mapping is the matched endpoint mapping; nil for object aggregation.
	Mapping *interfaceschema.CompiledMapping
	// Path is the concrete data path ("" only for flat object aggregation).
	Path string
	// Value is the decoded payload value (payload.Value closed set;
	// map[string]payload.Value for objects; nil for property unset).
	Value payload.Value
	// Format is the wire format the payload arrived in.
	Format payload.Format
	// TS is the effective sample timestamp: the explicit `t` when the
	// mapping declares explicit_timestamp, broker reception time otherwise.
	TS time.Time
	// ReceptionTS is the broker reception timestamp.
	ReceptionTS time.Time

	// ack releases the broker acknowledgment; the batcher calls it after
	// commit (docs/DESIGN.md §5.3), or immediately if the op is consumed
	// without persistence (broken-op path).
	ack func()
	// broken marks an op skipped by the batcher (conversion or poisoned-row
	// failure); it is consumed (acked) but never persisted nor observed.
	broken bool
}

// handle processes one inbound message on its shard goroutine: resolve the
// realm, classify the topic (§3.3), and dispatch. Introspection and control
// messages go to the M6b handler seams; until those are wired they are
// acknowledged, counted, and dropped.
func (e *Engine) handle(ctx context.Context, sh *shard, m broker.InboundMessage) {
	realm := e.schemas.realmOrReload(ctx, m.Realm)
	if realm == nil {
		e.reject(m, reasonRealmUnknown, "realm not in schema snapshot")
		return
	}
	rest, ok := splitDeviceTopic(m.Topic, m.Realm, m.DeviceID)
	if !ok {
		e.reject(m, reasonMalformedTopic, "topic does not carry the device root")
		return
	}

	kind, subpath := classify(rest)
	switch kind {
	case kindIntrospection:
		if e.onIntrospection != nil {
			// Flush first: introspection (and control) handlers act on
			// persisted state, so the device's pending data ops must commit
			// before them to preserve in-shard ordering semantics.
			e.flushShard(ctx, sh)
			e.onIntrospection(ctx, m, realm)
			return
		}
		e.met.unhandled.WithLabelValues("introspection").Inc()
		e.log.Debug("introspection handler not wired; message consumed",
			"realm", m.Realm, "device", m.DeviceID.String())
		m.Ack()
		return
	case kindControl:
		if e.onControl != nil {
			e.flushShard(ctx, sh)
			e.onControl(ctx, m, realm, subpath)
			return
		}
		e.met.unhandled.WithLabelValues("control").Inc()
		e.log.Debug("control handler not wired; message consumed",
			"realm", m.Realm, "device", m.DeviceID.String(), "subpath", subpath)
		m.Ack()
		return
	case kindData:
	}

	dev, ok := e.deviceState(ctx, m, realm)
	if !ok {
		return // rejected or abandoned; deviceState handled the message
	}
	e.processData(ctx, sh, m, realm, dev, rest)
}

// deviceState loads the per-device cache entry, parking the shard with
// exponential backoff on transient store failures (the §1.4 backpressure
// then propagates through the filling shard channel). On shutdown or context
// cancellation the message is abandoned unacknowledged (re-sent for
// QoS >= 1, dropped for QoS 0 — §5.3 degradation order).
func (e *Engine) deviceState(ctx context.Context, m broker.InboundMessage, realm *realmSchema) (*deviceState, bool) {
	backoff := parkBackoffStart
	for {
		dev, err := e.devices.get(ctx, m.Realm, realm.id, m.DeviceID)
		if err == nil {
			return dev, true
		}
		if errors.Is(err, store.ErrNotFound) {
			e.reject(m, reasonDeviceUnknown, err.Error())
			return nil, false
		}
		e.met.internalErrors.Inc()
		e.log.Warn("device state load failed; parking shard",
			"realm", m.Realm, "device", m.DeviceID.String(), "backoff", backoff, "err", err)
		select {
		case <-e.quit:
			return nil, false
		case <-ctx.Done():
			return nil, false
		case <-time.After(backoff):
		}
		backoff = min(backoff*2, parkBackoffCap)
	}
}

// processData runs the §2.6 validation pipeline on a data publish and, on
// success, appends the resulting PersistOp to the shard batch.
func (e *Engine) processData(ctx context.Context, sh *shard, m broker.InboundMessage, realm *realmSchema, dev *deviceState, rest string) {
	// Step 1 — introspection gate (doubles as interface-name resolution).
	name, path, ver, ok := matchInterface(rest, dev.declares)
	if !ok {
		e.reject(m, reasonInterfaceNotDeclared, "no introspected interface matches "+rest)
		return
	}
	ci := realm.iface(name, ver.Major)
	if ci == nil {
		e.reject(m, reasonInterfaceNotInstalled,
			fmt.Sprintf("introspection declares %s v%d but the realm has no such interface", name, ver.Major))
		return
	}
	// Step 3 — ownership (checked before decoding: cheaper, and the payload
	// of a misdirected publish is meaningless anyway).
	if ci.Ownership != interfaceschema.OwnershipDevice {
		e.reject(m, reasonOwnershipViolation, "device publish on server-owned interface "+name)
		return
	}

	// Steps 2 and 4–6 — endpoint match, decode, type check, aggregation
	// shape (pkg/payload owns 4–6).
	var (
		dp      payload.DecodedPayload
		mapping *interfaceschema.CompiledMapping
		err     error
	)
	if ci.Aggregation == interfaceschema.AggregationObject {
		if !objectPathOK(ci, path) {
			e.reject(m, reasonUnexpectedPath,
				fmt.Sprintf("%q is not an aggregation prefix of %s", path, name))
			return
		}
		dp, err = e.dec.Object(m.Payload, ci.ObjectLeaves)
	} else {
		mapping, ok = ci.Trie.Match(path)
		if !ok {
			e.reject(m, reasonUnexpectedPath,
				fmt.Sprintf("%q matches no endpoint of %s", path, name))
			return
		}
		dp, err = e.dec.Individual(m.Payload, mapping)
	}
	if err != nil {
		reason := payload.ReasonOf(err)
		if reason == payload.ReasonNone {
			// Not a payload rejection: a caller-bug error (nil mapping et
			// al.), which the pipeline construction rules out. Be loud.
			e.met.internalErrors.Inc()
			e.log.Error("payload decode failed without a reject reason",
				"realm", m.Realm, "device", m.DeviceID.String(), "topic", m.Topic, "err", err)
			m.Ack()
			return
		}
		e.reject(m, reason.String(), err.Error())
		return
	}

	e.noteFormat(ctx, m.DeviceID, realm, dev, dp.Format)

	ts := m.ReceivedAt
	if dp.Timestamp != nil {
		ts = *dp.Timestamp
	}
	op := PersistOp{
		Realm:       m.Realm,
		RealmID:     realm.id,
		DeviceID:    m.DeviceID,
		Interface:   ci,
		Mapping:     mapping,
		Path:        path,
		Value:       dp.Value,
		Format:      dp.Format,
		TS:          ts,
		ReceptionTS: m.ReceivedAt,
		ack:         m.Ack,
	}
	switch {
	case ci.Type == interfaceschema.Properties && dp.IsUnset():
		op.Kind = OpPropertyUnset
	case ci.Type == interfaceschema.Properties:
		op.Kind = OpPropertySet
	case ci.Aggregation == interfaceschema.AggregationObject:
		op.Kind = OpObject
	default:
		op.Kind = OpIndividual
	}
	sh.batch.add(op)
}

// objectPathOK validates an object-aggregation publish path: appending any
// declared last-level key must resolve in the endpoint trie, which checks
// prefix depth and placeholder hygiene in one go (docs/DESIGN.md §2.6
// step 6).
func objectPathOK(ci *interfaceschema.CompiledInterface, prefix string) bool {
	for leaf := range ci.ObjectLeaves {
		_, ok := ci.Trie.Match(prefix + "/" + leaf)
		return ok
	}
	return false
}

// noteFormat maintains the device's sticky payload-format hint
// (docs/DESIGN.md §3.5.4): the first JSON data payload flips it to json;
// after an emptyCache armed the reset, the next BSON data payload flips it
// back. Hint persistence failures are logged, not fatal — the in-memory
// state self-heals on the next flip.
func (e *Engine) noteFormat(ctx context.Context, id deviceid.ID, realm *realmSchema, dev *deviceState, f payload.Format) {
	var flip string
	dev.mu.Lock()
	switch f {
	case payload.FormatJSON:
		dev.resetHintOnBSON = false
		if dev.formatHint != hintJSON {
			dev.formatHint = hintJSON
			flip = hintJSON
		}
	case payload.FormatBSON:
		if dev.resetHintOnBSON {
			dev.resetHintOnBSON = false
			if dev.formatHint != hintBSON {
				dev.formatHint = hintBSON
				flip = hintBSON
			}
		}
	case payload.FormatEmpty, payload.FormatInvalid:
		// Property unsets carry no format information.
	}
	dev.mu.Unlock()
	if flip == "" {
		return
	}
	if err := e.st.SetPayloadFormatHint(ctx, realm.id, id, flip); err != nil {
		e.log.Warn("payload format hint not persisted",
			"realm", realm.name, "device", id.String(), "hint", flip, "err", err)
	}
}

// reject consumes an invalid message: per-reason metric, structured log,
// device_error seam (M6b), and acknowledgment — devices must not stall on
// data the platform will never accept (docs/DESIGN.md §2.6).
func (e *Engine) reject(m broker.InboundMessage, reason, detail string) {
	e.met.rejects.WithLabelValues(reason).Inc()
	e.log.Warn("message rejected",
		"reason", reason, "realm", m.Realm, "device", m.DeviceID.String(),
		"topic", m.Topic, "detail", detail)
	if e.onDeviceError != nil {
		e.onDeviceError(m, reason, detail)
	}
	m.Ack()
}
