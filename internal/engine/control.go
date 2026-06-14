package engine

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// Control-channel subpaths (docs/DESIGN.md §3.3–3.4).
const (
	controlEmptyCache         = "emptyCache"
	controlProducerProperties = "producer/properties"
	controlConsumerProperties = "consumer/properties"
)

// maxControlInflated is the absolute ceiling on the inflated size of a zlib
// control payload — the zip-bomb guard of docs/DESIGN.md §4.5 ("hard-capped
// at the declared size and an absolute ceiling"). 1 MiB comfortably holds
// tens of thousands of property entries.
const maxControlInflated = 1 << 20

// controlFrameHeader is the 4-byte big-endian uncompressed-size prefix.
const controlFrameHeader = 4

// handleControl dispatches a device control publish
// ("<realm>/<device_id>/control/<subpath>", docs/ROADMAP.md §7.2 file 6.8)
// on the device's shard goroutine. Unknown subpaths are rejected and
// consumed — the broker ACL only admits the two device-publishable control
// topics, so anything else here is a defensive path.
func (e *Engine) handleControl(ctx context.Context, m broker.InboundMessage, realm *realmSchema, subpath string) {
	switch subpath {
	case controlEmptyCache:
		e.handleEmptyCache(ctx, m, realm)
	case controlProducerProperties:
		e.handleProducerProperties(ctx, m, realm)
	default:
		e.reject(m, reasonControlUnknown, "unknown control subpath "+subpath)
	}
}

// handleEmptyCache implements `control/emptyCache` (docs/DESIGN.md §3.3):
// the device lost its local cache, so Astrate re-sends every server-owned
// property on its data topic (QoS 2, format per the device hint) and then
// publishes the `consumer/properties` purge message. It also arms the
// §3.5.4 hint reset: if the device's next data payload is BSON, the sticky
// json hint flips back.
func (e *Engine) handleEmptyCache(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
	dev, ok := e.deviceState(ctx, m, realm)
	if !ok {
		return
	}
	dev.armHintReset()

	format := formatForHint(dev.hint())
	if !e.retryStore(ctx, m, "emptyCache resync", func() error {
		if err := e.resendServerProperties(ctx, realm, m.DeviceID, format); err != nil {
			return err
		}
		return e.sendConsumerProperties(ctx, realm, m.DeviceID)
	}) {
		return
	}
	m.Ack()
}

// handleProducerProperties implements `control/producer/properties`
// (docs/DESIGN.md §3.3): the payload is the exhaustive zlib-compressed list
// of device-owned properties the device still holds, and every device-owned
// row not in it is purged — the device is the source of truth for its own
// properties. Malformed payloads (bad frame, zip bomb) are rejected and
// consumed.
func (e *Engine) handleProducerProperties(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
	entries, err := inflateProperties(m.Payload)
	if err != nil {
		e.reject(m, reasonControlInvalid, err.Error())
		return
	}
	dev, ok := e.deviceState(ctx, m, realm)
	if !ok {
		return
	}

	keep := e.resolvePropertyRefs(realm, dev, entries)
	if !e.retryStore(ctx, m, "producer/properties purge", func() error {
		purged, err := e.st.PurgeDeviceOwnedExcept(ctx, realm.id, m.DeviceID, keep)
		if err == nil && purged > 0 {
			e.log.Info("purged device-owned properties",
				"realm", m.Realm, "device", m.DeviceID.String(), "purged", purged, "kept", len(keep))
		}
		return err
	}) {
		return
	}
	m.Ack()
}

// resolvePropertyRefs maps "interface_name/path" entries to storage
// references. The interface major is taken from the device's introspection,
// falling back to the highest installed major (an entry sent during an
// introspection update race must still protect its row). Unresolvable
// entries are skipped with a log: they cannot reference an existing row.
func (e *Engine) resolvePropertyRefs(realm *realmSchema, dev *deviceState, entries []string) []store.PropertyRef {
	keep := make([]store.PropertyRef, 0, len(entries))
	for _, entry := range entries {
		name, rest, ok := strings.Cut(entry, "/")
		if !ok || name == "" || rest == "" {
			e.log.Warn("skipping malformed producer/properties entry",
				"realm", realm.name, "entry", entry)
			continue
		}
		var ci *interfaceschema.CompiledInterface
		if v, declared := dev.declares(name); declared {
			ci = realm.iface(name, v.Major)
		}
		if ci == nil {
			ci = realm.latestIface(name)
		}
		if ci == nil {
			e.log.Warn("skipping producer/properties entry for unknown interface",
				"realm", realm.name, "interface", name)
			continue
		}
		keep = append(keep, store.PropertyRef{InterfaceID: ci.ID, Path: "/" + rest})
	}
	return keep
}

// resendServerProperties publishes every server-owned property of the
// device on its data topic (QoS 2, the emptyCache resync of docs/DESIGN.md
// §3.3). Rows that no longer resolve against the compiled schema — an
// interface deleted mid-flight — are skipped loudly; store errors propagate
// for the caller's retry loop.
func (e *Engine) resendServerProperties(ctx context.Context, realm *realmSchema, id deviceid.ID, format payload.Format) error {
	props, err := e.st.ListServerOwnedProperties(ctx, realm.id, id)
	if err != nil {
		return err
	}
	for i := range props {
		p := &props[i]
		ci := realm.ifaceByID(p.InterfaceID)
		if ci == nil {
			e.log.Warn("server-owned property references an uninstalled interface; skipping resend",
				"realm", realm.name, "device", id.String(), "interface_id", p.InterfaceID, "path", p.Path)
			continue
		}
		mapping, ok := ci.Trie.Match(p.Path)
		if !ok {
			e.log.Warn("server-owned property path matches no endpoint; skipping resend",
				"realm", realm.name, "device", id.String(), "interface", ci.Name, "path", p.Path)
			continue
		}
		val, err := e.rehydrateValue(p.Value, mapping)
		if err != nil {
			e.log.Error("stored property value does not rehydrate; skipping resend",
				"realm", realm.name, "device", id.String(), "interface", ci.Name, "path", p.Path, "err", err)
			continue
		}
		wire, err := payload.Encode(val, nil, format)
		if err != nil {
			e.log.Error("property value does not encode; skipping resend",
				"realm", realm.name, "device", id.String(), "interface", ci.Name, "path", p.Path, "err", err)
			continue
		}
		topic := deviceTopic(realm.name, id, ci.Name+p.Path)
		if err := e.broker.Publish(topic, wire, 2, false, 0); err != nil {
			return fmt.Errorf("engine: republishing %s: %w", topic, err)
		}
	}
	return nil
}

// sendConsumerProperties publishes the `control/consumer/properties` purge
// message: the zlib-framed exhaustive list of currently-set server-owned
// property paths (docs/DESIGN.md §3.4). It is sent after emptyCache, after
// a device (re)connects, and after a server-owned property unset; the list
// is always the full current truth, so the message is idempotent and safe
// to send at any time. Control payloads keep the zlib+size framing for JSON
// devices too (docs/DESIGN.md §3.5.4).
func (e *Engine) sendConsumerProperties(ctx context.Context, realm *realmSchema, id deviceid.ID) error {
	props, err := e.st.ListServerOwnedProperties(ctx, realm.id, id)
	if err != nil {
		return err
	}
	entries := make([]string, 0, len(props))
	for i := range props {
		ci := realm.ifaceByID(props[i].InterfaceID)
		if ci == nil {
			continue
		}
		entries = append(entries, ci.Name+props[i].Path)
	}
	frame, err := deflateProperties(entries)
	if err != nil {
		return fmt.Errorf("engine: building consumer/properties payload: %w", err)
	}
	topic := deviceTopic(realm.name, id, controlPrefix+controlConsumerProperties)
	if err := e.broker.Publish(topic, frame, 2, false, 0); err != nil {
		return fmt.Errorf("engine: publishing %s: %w", topic, err)
	}
	return nil
}

// rehydrateValue turns a stored jsonb property value back into a typed
// payload.Value by wrapping it in a JSON-profile envelope and running the
// standard decoder: the stored rendering (docs/DESIGN.md §2.3 — base64
// blobs, RFC 3339 datetimes, big longintegers as decimal strings) is exactly
// the §3.5.3 JSON profile.
func (e *Engine) rehydrateValue(stored []byte, m *interfaceschema.CompiledMapping) (payload.Value, error) {
	env := make([]byte, 0, len(stored)+8)
	env = append(env, `{"v":`...)
	env = append(env, stored...)
	env = append(env, '}')
	dp, err := payload.Decoder{MaxSize: len(env)}.Individual(env, &interfaceschema.CompiledMapping{
		ValueType:  m.ValueType,
		AllowUnset: m.AllowUnset,
		// ExplicitTimestamp deliberately false: stored properties carry no t.
	})
	if err != nil {
		return nil, err
	}
	return dp.Value, nil
}

// deviceTopic builds "<realm>/<device_id>/<rest>".
func deviceTopic(realm string, id deviceid.ID, rest string) string {
	return realm + "/" + id.String() + "/" + rest
}

// formatForHint maps a devices.payload_format_hint value to the outbound
// wire format (docs/DESIGN.md §3.5.4).
func formatForHint(hint string) payload.Format {
	if hint == hintJSON {
		return payload.FormatJSON
	}
	return payload.FormatBSON
}

// deflateProperties builds a zlib control payload: 4-byte big-endian
// uncompressed size, then the zlib-deflated `;`-joined entry list
// (docs/DESIGN.md §3.3–3.4).
func deflateProperties(entries []string) ([]byte, error) {
	plain := strings.Join(entries, ";")
	var buf bytes.Buffer
	buf.Write(binary.BigEndian.AppendUint32(nil, uint32(len(plain)))) // #nosec G115 -- property lists are far below 4 GiB
	zw := zlib.NewWriter(&buf)
	if _, err := io.WriteString(zw, plain); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// inflateProperties parses a zlib control payload with the docs/DESIGN.md
// §4.5 bounds: the declared size is capped by maxControlInflated, and the
// stream may not inflate beyond what it declared (a lying header is a
// zip-bomb attempt, not a tolerable client quirk). It returns the entry
// list; an empty payload yields no entries (the SDK deliberately sends an
// empty list when the device holds no properties).
func inflateProperties(p []byte) ([]string, error) {
	if len(p) < controlFrameHeader {
		return nil, fmt.Errorf("control payload is %d bytes, below the %d byte size prefix", len(p), controlFrameHeader)
	}
	declared := binary.BigEndian.Uint32(p[:controlFrameHeader])
	if declared > maxControlInflated {
		return nil, fmt.Errorf("control payload declares %d inflated bytes, above the %d byte ceiling", declared, maxControlInflated)
	}
	zr, err := zlib.NewReader(bytes.NewReader(p[controlFrameHeader:]))
	if err != nil {
		return nil, fmt.Errorf("control payload is not a zlib stream: %w", err)
	}
	defer func() { _ = zr.Close() }()

	plain, err := io.ReadAll(io.LimitReader(zr, int64(declared)+1))
	if err != nil {
		return nil, fmt.Errorf("inflating control payload: %w", err)
	}
	if len(plain) > int(declared) {
		return nil, fmt.Errorf("control payload inflates beyond its declared %d bytes", declared)
	}
	if len(plain) == 0 {
		return nil, nil
	}
	return strings.Split(string(plain), ";"), nil
}
