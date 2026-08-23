package blocks

import (
	"fmt"

	"github.com/astrate-platform/astrate/internal/flow"
)

// UpdateMetadata merges set_metadata into a message's Metadata and then
// deletes delete_metadata keys (delete wins for the same key). Payload, key
// and timestamp are untouched.
//
// Config keys (all optional; at least one required):
//   - set_metadata (object string→string): merged in, overwriting existing keys
//   - delete_metadata ([]string): removed after set applies
//
// The message is shallow-copied so concurrent lanes never share Metadata maps.
func UpdateMetadata(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
	cfg, err := parseUpdateMetadataConfig(config)
	if err != nil {
		return nil, fmt.Errorf("update_metadata: %w", err)
	}
	return flow.NewTransformBlock(name, func(msg *flow.Message) ([]*flow.Message, error) {
		if msg == nil {
			return nil, nil
		}
		out := cloneMessage(msg)
		if len(cfg.setMetadata) > 0 {
			if out.Metadata == nil {
				out.Metadata = make(map[string]string, len(cfg.setMetadata))
			}
			for k, v := range cfg.setMetadata {
				out.Metadata[k] = v
			}
		}
		for _, k := range cfg.deleteMetadata {
			delete(out.Metadata, k)
		}
		return []*flow.Message{out}, nil
	}), nil
}

type updateMetadataConfig struct {
	setMetadata    map[string]string
	deleteMetadata []string
}

func parseUpdateMetadataConfig(config map[string]any) (updateMetadataConfig, error) {
	cfg := updateMetadataConfig{
		setMetadata:    stringMapConfig(config, "set_metadata"),
		deleteMetadata: stringSliceConfig(config, "delete_metadata"),
	}
	if len(cfg.setMetadata) == 0 && len(cfg.deleteMetadata) == 0 {
		return cfg, fmt.Errorf("at least one of set_metadata, delete_metadata is required")
	}
	return cfg, nil
}
