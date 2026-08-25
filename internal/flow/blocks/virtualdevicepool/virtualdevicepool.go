// Package virtualdevicepool implements the virtual_device_pool block
// (issue #84): it publishes pipeline messages as registered virtual
// devices through the engine ingest path — storage rows without MQTT.
package virtualdevicepool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
)

// Type is the block_type string stored in pipeline definitions.
const Type = "virtual_device_pool"

// Constructor builds a virtual_device_pool sink (issue #84). Each message's
// key addresses the target as <device_id>/<interface></path…> — upstream's
// key carries a realm segment first; astrate drops it because flows are
// per-realm. The payload is JSON-encoded and landed through deps.Ingest,
// the engine's device-owned ingest.
//
// Config keys:
//   - devices (string array): registered device_ids this pool may publish
//     as. Required unless auto_register is true, in which case it may be
//     omitted or empty and any entries present pre-seed the pool.
//   - auto_register (optional bool): when true a message addressed to an
//     unknown device id registers that id through deps.Register on first
//     sight (requires the pairing register dependency) and is published as
//     usual. The minted secret is discarded — astrate keeps credentials
//     server-side. An id that turns out to be already registered and
//     confirmed drops the message and is never cached, so every later
//     message re-attempts registration (upstream parity).
//
// Messages with an unusable key, a device outside the pool (with
// auto_register off), or an unmarshalable payload are logged and dropped —
// never an error. Ingest failures — and with auto_register on, register
// failures other than already-registered — are infrastructure problems and
// are returned wrapped.
func Constructor(name string, config map[string]any, deps flow.Deps) (flow.Block, error) {
	autoRegister, err := parseAutoRegister(config)
	if err != nil {
		return nil, err
	}
	devices, err := parseDevices(config, autoRegister)
	if err != nil {
		return nil, err
	}
	if deps.Ingest == nil {
		return nil, errors.New("virtual_device_pool requires the engine ingest dependency")
	}
	if autoRegister && deps.Register == nil {
		return nil, errors.New("virtual_device_pool requires the pairing register dependency")
	}
	p := &pool{
		devices:      make(map[string]struct{}, len(devices)),
		ingest:       deps.Ingest,
		register:     deps.Register,
		autoRegister: autoRegister,
		realm:        deps.Realm,
	}
	for _, d := range devices {
		p.devices[d] = struct{}{}
	}
	return flow.NewSinkBlock(name, func(msg *flow.Message) error {
		return p.publish(name, msg)
	}), nil
}

// parseAutoRegister coerces the auto_register config entry into a bool,
// accepting a Go bool or a JSON-decoded bool (the same Go type). Absent or
// nil means false.
func parseAutoRegister(config map[string]any) (bool, error) {
	v, ok := config["auto_register"]
	if !ok || v == nil {
		return false, nil
	}
	b, isBool := v.(bool)
	if !isBool {
		return false, errors.New("virtual_device_pool: auto_register must be a boolean")
	}
	return b, nil
}

// parseDevices coerces the devices config entry into []string, accepting a
// []string or a JSON-decoded []any of strings, rejecting empty entries by
// index. With auto_register the entry is optional: absent, nil or empty
// yields no pre-seeded devices; entries still validate the same way.
func parseDevices(config map[string]any, autoRegister bool) ([]string, error) {
	const prefix = "virtual_device_pool: "
	v, ok := config["devices"]
	if !ok || v == nil {
		if autoRegister {
			return nil, nil
		}
		return nil, errors.New(prefix + "devices is required")
	}
	var raw []any
	switch items := v.(type) {
	case []string:
		raw = make([]any, len(items))
		for i, s := range items {
			raw[i] = s
		}
	case []any:
		raw = items
	default:
		return nil, errors.New(prefix + "devices is required")
	}
	if len(raw) == 0 {
		if autoRegister {
			return nil, nil
		}
		return nil, errors.New(prefix + "devices must not be empty")
	}
	out := make([]string, len(raw))
	for i, item := range raw {
		s, isStr := item.(string)
		if !isStr || s == "" {
			return nil, fmt.Errorf("%sdevices[%d] must be a non-empty string", prefix, i)
		}
		out[i] = s
	}
	return out, nil
}

// pool is one block instance's shared state across messages: the known
// device set plus the ingest door, and — with auto_register — the pairing
// register door for first-seen ids. mu guards devices across the whole
// check-register-add sequence so concurrent first-seen messages cannot
// double-register a device.
type pool struct {
	mu           sync.Mutex
	devices      map[string]struct{}
	ingest       flow.IngestFunc
	register     flow.RegisterFunc
	autoRegister bool
	realm        string
}

// publish maps one message to a virtual-device ingest call. Bad messages are
// warned about and dropped; only infrastructure failures (ingest errors,
// and with auto_register unexpected register errors) surface as errors.
func (p *pool) publish(name string, msg *flow.Message) error {
	logger := slog.Default()
	if msg == nil {
		return nil
	}
	segments := strings.Split(msg.Key, "/")
	if len(segments) < 3 || segments[0] == "" || segments[1] == "" {
		logger.Warn("virtual_device_pool dropping message: key must be <device>/<interface></path>",
			"block", name, "key", msg.Key)
		return nil
	}
	device, iface := segments[0], segments[1]
	dropWarning, err := p.resolveDevice(context.Background(), device)
	if err != nil {
		return fmt.Errorf("virtual_device_pool %s: %w", name, err)
	}
	if dropWarning != "" {
		logger.Warn(dropWarning, "block", name, "key", msg.Key, "device", device)
		return nil
	}
	path := "/" + strings.Join(segments[2:], "/")

	payload, err := json.Marshal(msg.Data)
	if err != nil {
		logger.Warn("virtual_device_pool dropping message: payload does not encode",
			"block", name, "key", msg.Key, "error", err)
		return nil
	}
	var ts *time.Time
	if msg.Timestamp != 0 {
		t := time.UnixMicro(msg.Timestamp)
		ts = &t
	}
	if err := p.ingest(context.Background(), p.realm, device, iface, path, payload, ts); err != nil {
		return fmt.Errorf("virtual_device_pool %s: %w", name, err)
	}
	return nil
}

// resolveDevice makes sure device may be published as, registering a
// first-seen id through the pairing door when auto_register is on. It
// returns "" when the message may proceed, a non-empty drop warning to log
// when the message must be dropped, or an unwrapped infrastructure error
// for the caller to wrap.
func (p *pool) resolveDevice(ctx context.Context, device string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.devices[device]; ok {
		return "", nil
	}
	if !p.autoRegister {
		return "virtual_device_pool dropping message: device not in pool", nil
	}
	switch err := p.register(ctx, p.realm, device); {
	case err == nil:
		p.devices[device] = struct{}{}
		return "", nil
	case errors.Is(err, flow.ErrVirtualDeviceRegistered):
		// Already registered and confirmed elsewhere: drop, and do NOT add
		// the id to the pool — every later message re-attempts
		// registration, mirroring upstream's dynamic pool.
		return "virtual_device_pool dropping message: cannot register device", nil
	default:
		return "", err
	}
}
