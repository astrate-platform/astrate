package flowapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v6"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/store"
)

// UserBlockView is the operator JSON shape for one stored user-defined
// composite block (#85). ConfigSchema is omitted when the block has none.
type UserBlockView struct {
	Name         string          `json:"name"`
	BlockType    string          `json:"block_type"`
	Source       json.RawMessage `json:"source"`
	ConfigSchema json.RawMessage `json:"config_schema,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// CreateUserBlock validates and stores a user-defined composite block. The
// source must expand cleanly right now against built-ins plus the realm's
// stored blocks, so self-cycles and references to not-yet-existing composites
// are refused at creation time (accepted v1 behaviour; upstream defers the
// same failure to build time).
func (s *Service) CreateUserBlock(ctx context.Context, realm, name, blockType string, source, configSchema json.RawMessage) (*UserBlockView, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: user block name is required", ErrValidation)
	}
	if s.reg.Has(name) {
		return nil, fmt.Errorf("%w: block type %q collides with a built-in block", ErrValidation, name)
	}
	if err := s.validateUserBlockBody(ctx, realm, blockType, source, configSchema); err != nil {
		return nil, err
	}
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	row, err := s.st.CreateUserBlock(ctx, id, &store.UserBlock{
		Name:         name,
		BlockType:    blockType,
		Source:       source,
		ConfigSchema: nilIfEmptyJSON(configSchema),
	})
	if err != nil {
		return nil, err
	}
	return toUserBlockView(row), nil
}

// GetUserBlock returns one stored user block by name.
func (s *Service) GetUserBlock(ctx context.Context, realm, name string) (*UserBlockView, error) {
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	row, err := s.st.GetUserBlock(ctx, id, name)
	if err != nil {
		return nil, err
	}
	return toUserBlockView(row), nil
}

// UpdateUserBlock replaces a user block's type, source and schema. Built-in
// (default) blocks cannot be modified through the user-block surface.
func (s *Service) UpdateUserBlock(ctx context.Context, realm, name, blockType string, source, configSchema json.RawMessage) (*UserBlockView, error) {
	if s.reg.Has(name) {
		return nil, fmt.Errorf("%w: built-in blocks cannot be modified", ErrValidation)
	}
	if err := s.validateUserBlockBody(ctx, realm, blockType, source, configSchema); err != nil {
		return nil, err
	}
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	row, err := s.st.UpdateUserBlock(ctx, id, name, blockType, source, nilIfEmptyJSON(configSchema))
	if err != nil {
		return nil, err
	}
	return toUserBlockView(row), nil
}

// DeleteUserBlock removes a stored user block. Built-in (default) blocks
// cannot be deleted through the user-block surface.
func (s *Service) DeleteUserBlock(ctx context.Context, realm, name string) error {
	if s.reg.Has(name) {
		return fmt.Errorf("%w: built-in blocks cannot be modified", ErrValidation)
	}
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return err
	}
	return s.st.DeleteUserBlock(ctx, id, name)
}

// ListUserBlocks returns every stored user block of a realm, ordered by name.
func (s *Service) ListUserBlocks(ctx context.Context, realm string) ([]UserBlockView, error) {
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	rows, err := s.st.ListUserBlocks(ctx, id)
	if err != nil {
		return nil, err
	}
	out := make([]UserBlockView, 0, len(rows))
	for i := range rows {
		out = append(out, *toUserBlockView(&rows[i]))
	}
	return out, nil
}

// HasBuiltInBlock reports whether blockType names a registered built-in;
// GET /blocks/{name} dispatches on it.
func (s *Service) HasBuiltInBlock(blockType string) bool {
	return s.reg.Has(blockType)
}

// validateUserBlockBody enforces the rules shared by create and update:
// block_type membership, optional config_schema legality, and a source that
// expands cleanly against the current stored state. Every failure wraps
// ErrValidation.
func (s *Service) validateUserBlockBody(ctx context.Context, realm, blockType string, source, configSchema json.RawMessage) error {
	switch blockType {
	case "producer", "consumer", "producer_consumer":
	default:
		return fmt.Errorf("%w: block_type must be producer, consumer or producer_consumer", ErrValidation)
	}
	if err := s.validateConfigSchema(configSchema); err != nil {
		return err
	}
	if _, err := flow.ExpandComposites(source, s.userBlockResolver(ctx, realm)); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return nil
}

// validateConfigSchema compiles a non-empty config_schema with the verified
// jsonschema/v6 pattern: decode, register under any unused name, compile
// (which rejects non-object schemas via the metaschema).
func (s *Service) validateConfigSchema(schema json.RawMessage) error {
	if len(schema) == 0 {
		return nil
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(schema))
	if err != nil {
		return fmt.Errorf("%w: config_schema is not valid JSON: %v", ErrValidation, err)
	}
	c := jsonschema.NewCompiler()
	if err := c.AddResource("config_schema.json", doc); err != nil {
		return fmt.Errorf("%w: config_schema is not a valid JSON Schema: %v", ErrValidation, err)
	}
	if _, err := c.Compile("config_schema.json"); err != nil {
		return fmt.Errorf("%w: config_schema is not a valid JSON Schema: %v", ErrValidation, err)
	}
	return nil
}

// userBlockResolver adapts the realm's stored blocks to the
// flow.ExpandComposites contract: built-ins pass through untouched, stored
// composites resolve by name, and anything else is a dangling reference.
// Realm resolution happens lazily inside the closure so sources that fail
// structurally never touch the store.
func (s *Service) userBlockResolver(ctx context.Context, realm string) func(string) (*flow.UserBlock, error) {
	return func(name string) (*flow.UserBlock, error) {
		if s.reg.Has(name) {
			return nil, nil
		}
		id, err := s.realmID(ctx, realm)
		if err != nil {
			return nil, err
		}
		rows, err := s.st.ListUserBlocks(ctx, id)
		if err != nil {
			return nil, err
		}
		for i := range rows {
			if rows[i].Name == name {
				return &flow.UserBlock{
					Name:         rows[i].Name,
					BlockType:    rows[i].BlockType,
					Source:       rows[i].Source,
					ConfigSchema: json.RawMessage(rows[i].ConfigSchema),
				}, nil
			}
		}
		return nil, fmt.Errorf("references composite %q which does not exist yet", name)
	}
}

func toUserBlockView(b *store.UserBlock) *UserBlockView {
	v := &UserBlockView{
		Name:      b.Name,
		BlockType: b.BlockType,
		Source:    json.RawMessage(b.Source),
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
	if len(b.ConfigSchema) > 0 {
		v.ConfigSchema = json.RawMessage(b.ConfigSchema)
	}
	return v
}

// nilIfEmptyJSON normalizes an absent value so nullable bytea columns stay
// NULL instead of zero-length.
func nilIfEmptyJSON(raw json.RawMessage) []byte {
	if len(raw) == 0 {
		return nil
	}
	return raw
}
