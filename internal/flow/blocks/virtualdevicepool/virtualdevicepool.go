// Package virtualdevicepool implements the virtual_device_pool block
// (issue #84): it publishes pipeline messages as registered virtual
// devices through the engine ingest path — storage rows without MQTT.
package virtualdevicepool

import (
	"errors"

	"github.com/astrate-platform/astrate/internal/flow"
)

// Type is the block_type string stored in pipeline definitions.
const Type = "virtual_device_pool"

// Constructor builds a virtual_device_pool sink. STUB pending phase 84b.
func Constructor(name string, config map[string]any, deps flow.Deps) (flow.Block, error) {
	return nil, errors.New("virtual_device_pool: not implemented")
}
