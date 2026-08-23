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
//   - devices (required string array): registered device_ids this pool may
//     publish as; at least one entry, every entry a non-empty string.
//
// Messages with an unusable key, a device outside the pool, or an
// unmarshalable payload are logged and dropped — never an error. Ingest
// failures are infrastructure problems and are returned wrapped.
func Constructor(name string, config map[string]any, deps flow.Deps) (flow.Block, error) {
	devices, err := parseDevices(config)
	if err != nil {
		return nil, err
	}
	if deps.Ingest == nil {
		return nil, errors.New("virtual_device_pool requires the engine ingest dependency")
	}
	pool := make(map[string]struct{}, len(devices))
	for _, d := range devices {
		pool[d] = struct{}{}
	}
	return flow.NewSinkBlock(name, func(msg *flow.Message) error {
		return publish(name, pool, deps.Ingest, deps.Realm, msg)
	}), nil
}

// parseDevices coerces the devices config entry into []string, accepting a
// []string or a JSON-decoded []any of strings, rejecting empty entries by
// index.
func parseDevices(config map[string]any) ([]string, error) {
	const prefix = "virtual_device_pool: "
	v, ok := config["devices"]
	if !ok || v == nil {
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

// publish maps one message to a virtual-device ingest call. Bad messages are
// warned about and dropped; only ingest failures surface as errors.
func publish(
	name string,
	pool map[string]struct{},
	ingest func(ctx context.Context, realm, deviceID, ifaceName, path string, payload json.RawMessage, ts *time.Time) error,
	realm string,
	msg *flow.Message,
) error {
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
	if _, ok := pool[device]; !ok {
		logger.Warn("virtual_device_pool dropping message: device not in pool",
			"block", name, "key", msg.Key, "device", device)
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
	if err := ingest(context.Background(), realm, device, iface, path, payload, ts); err != nil {
		return fmt.Errorf("virtual_device_pool %s: %w", name, err)
	}
	return nil
}
