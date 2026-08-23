package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// Sentinel errors of the server-owned publish path (docs/ROADMAP.md §7.2
// file 6.9). The AppEngine layer (M7) maps them onto upstream-shaped HTTP
// statuses; payload validation failures surface as *payload.RejectError.
var (
	// ErrRealmUnknown: the realm is not in the schema snapshot.
	ErrRealmUnknown = errors.New("engine: realm unknown")
	// ErrInterfaceNotFound: the realm has no installed interface by that name.
	ErrInterfaceNotFound = errors.New("engine: interface not installed")
	// ErrNotServerOwned: the interface is device-owned (docs/DESIGN.md §2.6
	// step 3: AppEngine publishes only on ownership: server).
	ErrNotServerOwned = errors.New("engine: interface is not server-owned")
	// ErrNotDeviceOwned: the interface is server-owned (the mirrored gate of
	// PublishDeviceValue, issue #84: device-owned ingest refuses them).
	ErrNotDeviceOwned = errors.New("engine: interface is not device-owned")
	// ErrPathNotFound: the path resolves no endpoint mapping.
	ErrPathNotFound = errors.New("engine: path matches no endpoint")
	// ErrNotAProperty: a property operation on a datastream interface.
	ErrNotAProperty = errors.New("engine: interface is not a properties interface")
	// ErrUnsetNotAllowed: unset on a mapping without allow_unset.
	ErrUnsetNotAllowed = errors.New("engine: mapping does not allow unset")
)

// PublishServerValue validates, persists, and delivers one server-owned
// value (docs/ROADMAP.md §7.2 file 6.9): AppEngine PUT/POST bodies land
// here. value is the raw JSON value (the unwrapped "data" field); it is
// validated against the mapping exactly like an inbound JSON-profile
// payload (docs/DESIGN.md §2.6 steps 4–5). ts is the optional explicit
// timestamp; reception time applies otherwise.
//
// Persistence happens before delivery: a property upsert or a datastream
// insert, then a broker publish on the device's data topic with the
// mapping's QoS, retain for properties, per-message expiry for datastreams,
// and the wire format chosen by the device's payload_format_hint
// (docs/DESIGN.md §3.4, §3.5.4).
func (e *Engine) PublishServerValue(ctx context.Context, realm string, id deviceid.ID, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
	out, err := e.publishAsOwner(ctx, realm, id, ifaceName, path, value, ts, interfaceschema.OwnershipServer)
	if err != nil {
		return err
	}
	wire, err := payload.Encode(out.dp.Value, out.wireTS, formatForHint(out.hint))
	if err != nil {
		return fmt.Errorf("engine: encoding server value: %w", err)
	}
	topic := deviceTopic(realm, id, ifaceName+path)
	if err := e.broker.Publish(topic, wire, out.qos, out.retain, out.expiry); err != nil {
		return fmt.Errorf("engine: publishing %s: %w", topic, err)
	}
	return nil
}

// PublishDeviceValue validates and persists one device-owned value without
// delivering it (issue #84): virtual-device ingest lands storage rows
// exactly like a real device's data would, but produces no MQTT traffic —
// the device speaks to the engine directly, so there is no upstream hop to
// relay onto the broker. Validation and persistence mirror
// PublishServerValue with the ownership gate flipped to device.
func (e *Engine) PublishDeviceValue(ctx context.Context, realm string, id deviceid.ID, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
	_, err := e.publishAsOwner(ctx, realm, id, ifaceName, path, value, ts, interfaceschema.OwnershipDevice)
	return err
}

// ownerDelivery carries one persisted value plus everything wire delivery
// needs to know about it.
type ownerDelivery struct {
	dp     payload.DecodedPayload
	hint   string
	wireTS *time.Time
	qos    byte
	retain bool
	expiry time.Duration
}

// publishAsOwner is the shared validate-and-persist core of the publish
// paths (docs/ROADMAP.md §7.2 file 6.9): realm, device, and interface
// resolution with the ownership gate parameterised, §2.6 payload
// validation, and the persist switch. ownership selects which side may
// write and which sentinel an ownership mismatch carries; delivery stays
// with the caller.
func (e *Engine) publishAsOwner(ctx context.Context, realm string, id deviceid.ID,
	ifaceName, path string, value json.RawMessage, ts *time.Time,
	ownership interfaceschema.Ownership) (*ownerDelivery, error) {
	notOwned := ErrNotDeviceOwned
	if ownership == interfaceschema.OwnershipServer {
		notOwned = ErrNotServerOwned
	}
	rs := e.schemas.realmOrReload(ctx, realm)
	if rs == nil {
		return nil, fmt.Errorf("%w: %s", ErrRealmUnknown, realm)
	}
	view, err := e.deviceView(ctx, rs, id)
	if err != nil {
		return nil, err
	}
	ci := resolveServerInterface(rs, view, ifaceName)
	if ci == nil {
		return nil, fmt.Errorf("%w: %s", ErrInterfaceNotFound, ifaceName)
	}
	if ci.Ownership != ownership {
		return nil, fmt.Errorf("%w: %s", notOwned, ifaceName)
	}

	effTS := time.Now().UTC()
	if ts != nil {
		effTS = ts.UTC()
	}
	envelope, err := serverEnvelope(value, effTS)
	if err != nil {
		return nil, err
	}

	var (
		dp      payload.DecodedPayload
		mapping *interfaceschema.CompiledMapping
	)
	dec := payload.Decoder{MaxSize: len(envelope)}
	if ci.Aggregation == interfaceschema.AggregationObject {
		if !objectPathOK(ci, path) {
			return nil, fmt.Errorf("%w: %q is not an aggregation prefix of %s", ErrPathNotFound, path, ifaceName)
		}
		dp, err = dec.Object(envelope, ci.ObjectLeaves)
		mapping = anyObjectLeaf(ci)
	} else {
		var ok bool
		mapping, ok = ci.Trie.Match(path)
		if !ok {
			return nil, fmt.Errorf("%w: %q matches no endpoint of %s", ErrPathNotFound, path, ifaceName)
		}
		dp, err = dec.Individual(envelope, mapping)
	}
	if err != nil {
		return nil, err
	}
	if dp.IsUnset() {
		return nil, fmt.Errorf("%w: null value (use the property DELETE path for unset)", ErrUnsetNotAllowed)
	}

	op := PersistOp{
		Realm:       realm,
		RealmID:     rs.id,
		DeviceID:    id,
		Interface:   ci,
		Mapping:     mapping,
		Path:        path,
		Value:       dp.Value,
		Format:      dp.Format,
		TS:          effTS,
		ReceptionTS: effTS,
	}
	out := &ownerDelivery{dp: dp, hint: view.hint, qos: 2}
	switch {
	case ci.Type == interfaceschema.Properties:
		op.Kind = OpPropertySet
		out.retain = true // properties stay retained (§3.4)
		row, err := propertyRow(&op)
		if err != nil {
			return nil, err
		}
		if err := e.st.UpsertProperty(ctx, *row); err != nil {
			return nil, err
		}
	case ci.Aggregation == interfaceschema.AggregationObject:
		op.Kind = OpObject
		out.wireTS = &effTS
		out.qos = mapping.Reliability
		out.expiry = mapping.Expiry
		row, err := objectRow(&op)
		if err != nil {
			return nil, err
		}
		if err := e.st.AppendDatastreams(ctx, store.DatastreamBatch{Objects: []store.ObjectRow{*row}}); err != nil {
			return nil, err
		}
	default:
		op.Kind = OpIndividual
		out.wireTS = &effTS
		out.qos = mapping.Reliability
		out.expiry = mapping.Expiry
		row, err := individualRow(&op)
		if err != nil {
			return nil, err
		}
		if err := e.st.AppendDatastreams(ctx, store.DatastreamBatch{Individual: []store.IndividualRow{*row}}); err != nil {
			return nil, err
		}
	}
	e.met.persistOps.WithLabelValues(op.Kind.String()).Inc()
	return out, nil
}

// UnsetServerProperty deletes a server-owned property (AppEngine DELETE,
// docs/ROADMAP.md §8.2 file 7.7): the row is removed, an empty retained
// payload clears the topic and signals the unset to a connected device, and
// the `consumer/properties` purge message follows so offline-window state
// converges (docs/DESIGN.md §3.4). Unsetting an absent property is a no-op,
// not an error.
func (e *Engine) UnsetServerProperty(ctx context.Context, realm string, id deviceid.ID, ifaceName, path string) error {
	rs := e.schemas.realmOrReload(ctx, realm)
	if rs == nil {
		return fmt.Errorf("%w: %s", ErrRealmUnknown, realm)
	}
	view, err := e.deviceView(ctx, rs, id)
	if err != nil {
		return err
	}
	ci := resolveServerInterface(rs, view, ifaceName)
	if ci == nil {
		return fmt.Errorf("%w: %s", ErrInterfaceNotFound, ifaceName)
	}
	if ci.Ownership != interfaceschema.OwnershipServer {
		return fmt.Errorf("%w: %s", ErrNotServerOwned, ifaceName)
	}
	if ci.Type != interfaceschema.Properties {
		return fmt.Errorf("%w: %s", ErrNotAProperty, ifaceName)
	}
	mapping, ok := ci.Trie.Match(path)
	if !ok {
		return fmt.Errorf("%w: %q matches no endpoint of %s", ErrPathNotFound, path, ifaceName)
	}
	if !mapping.AllowUnset {
		return fmt.Errorf("%w: %s%s", ErrUnsetNotAllowed, ifaceName, path)
	}

	existed, err := e.st.UnsetProperty(ctx, rs.id, id, ci.ID, path)
	if err != nil {
		return err
	}
	if existed {
		e.met.persistOps.WithLabelValues(OpPropertyUnset.String()).Inc()
	}

	topic := deviceTopic(realm, id, ifaceName+path)
	if err := e.broker.Publish(topic, []byte{}, 2, true, 0); err != nil {
		return fmt.Errorf("engine: publishing %s: %w", topic, err)
	}
	return e.sendConsumerProperties(ctx, rs, id)
}

// deviceView is the introspection + format-hint view of a device used by
// the server-data path: the live cache entry when the device is cached,
// otherwise a one-shot store read (offline devices must not populate a
// cache that is only evicted on disconnect).
type deviceView struct {
	declares func(name string) (store.InterfaceVersion, bool)
	hint     string
}

// deviceView resolves the view; store.ErrNotFound propagates (M7 maps it to
// the 404 device envelope).
func (e *Engine) deviceView(ctx context.Context, rs *realmSchema, id deviceid.ID) (deviceView, error) {
	if dev := e.devices.peek(rs.name, id); dev != nil {
		return deviceView{declares: dev.declares, hint: dev.hint()}, nil
	}
	d, err := e.st.GetDevice(ctx, rs.id, id)
	if err != nil {
		return deviceView{}, err
	}
	hint := d.PayloadFormatHint
	if hint == "" {
		hint = hintBSON
	}
	return deviceView{
		declares: func(name string) (store.InterfaceVersion, bool) {
			v, ok := d.Introspection[name]
			return v, ok
		},
		hint: hint,
	}, nil
}

// resolveServerInterface picks the interface major to publish against: the
// device's introspected major when declared, the highest installed major
// otherwise (so server-owned properties can be staged before a device's
// first connection — a documented lenient superset of upstream).
func resolveServerInterface(rs *realmSchema, view deviceView, name string) *interfaceschema.CompiledInterface {
	if v, declared := view.declares(name); declared {
		if ci := rs.iface(name, v.Major); ci != nil {
			return ci
		}
	}
	return rs.latestIface(name)
}

// anyObjectLeaf returns one leaf mapping of an object-aggregated interface
// (reliability/expiry are uniform across an object's mappings by upstream
// validation, so any leaf is representative).
func anyObjectLeaf(ci *interfaceschema.CompiledInterface) *interfaceschema.CompiledMapping {
	for _, m := range ci.ObjectLeaves {
		return m
	}
	return nil
}

// serverEnvelope wraps a raw JSON value into the §3.5.3 `{v, t}` envelope so
// the standard decoder applies the §2.6 type rules. The timestamp is always
// present: explicit-timestamp mappings require it, and the decoder ignores
// it elsewhere (upstream leniency).
func serverEnvelope(value json.RawMessage, ts time.Time) ([]byte, error) {
	if len(value) == 0 {
		return nil, fmt.Errorf("%w: empty value", ErrUnsetNotAllowed)
	}
	env := struct {
		V json.RawMessage `json:"v"`
		T string          `json:"t"`
	}{V: value, T: ts.UTC().Format(jsonTimeLayout)}
	return json.Marshal(env)
}
