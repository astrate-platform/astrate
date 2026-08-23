package flowapi

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/flow/blocks"
)

// testUserBlockService is a validation-layer-only service: no store, so every
// test here must stay inside the rules that run before persistence.
func newUserBlockTestService() *Service {
	return &Service{reg: blocks.DefaultRegistry()}
}

// validUserBlockSource expands cleanly against built-ins only (random_source,
// log_sink), so validating it never needs the store.
const validUserBlockSource = `{"blocks":[` +
	`{"name":"picked","block_type":"random_source"},` +
	`{"name":"sink","block_type":"log_sink","config":{"level":"info"}}` +
	`],"connections":[{"from":"picked","to":"sink"}]}`

func assertValidation(t *testing.T, err error, wantIn string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected ErrValidation, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
	if wantIn != "" && !strings.Contains(err.Error(), wantIn) {
		t.Fatalf("error %q does not mention %q", err.Error(), wantIn)
	}
}

func TestCreateUserBlock_RejectsEmptyName(t *testing.T) {
	s := newUserBlockTestService()
	_, err := s.CreateUserBlock(t.Context(), "acme", "", "producer", json.RawMessage(validUserBlockSource), nil)
	assertValidation(t, err, "user block name is required")
}

func TestCreateUserBlock_RejectsBuiltInCollision(t *testing.T) {
	s := newUserBlockTestService()
	_, err := s.CreateUserBlock(t.Context(), "acme", "filter", "producer", json.RawMessage(validUserBlockSource), nil)
	assertValidation(t, err, `block type "filter" collides with a built-in block`)
}

func TestUpdateUserBlock_RefusesBuiltIn(t *testing.T) {
	s := newUserBlockTestService()
	src := json.RawMessage(validUserBlockSource)
	_, err := s.UpdateUserBlock(t.Context(), "acme", "filter", "producer", src, nil)
	assertValidation(t, err, "built-in blocks cannot be modified")
	err = s.DeleteUserBlock(t.Context(), "acme", "filter")
	assertValidation(t, err, "built-in blocks cannot be modified")
}

func TestUserBlockBody_BlockTypeRules(t *testing.T) {
	s := newUserBlockTestService()
	src := json.RawMessage(validUserBlockSource)
	err := s.validateUserBlockBody(t.Context(), "acme", "nope", src, nil)
	assertValidation(t, err, "block_type must be producer, consumer or producer_consumer")
	for _, good := range []string{"producer", "consumer", "producer_consumer"} {
		if err := s.validateUserBlockBody(t.Context(), "acme", good, src, nil); err != nil {
			t.Errorf("block_type %q must pass validation: %v", good, err)
		}
	}
}

func TestUserBlockBody_ConfigSchemaRules(t *testing.T) {
	s := newUserBlockTestService()
	src := json.RawMessage(validUserBlockSource)
	for name, schema := range map[string]string{
		"not JSON":   `not json`,
		"bad type":   `{"type":"nope"}`,
		"non-object": `[]`,
	} {
		err := s.validateUserBlockBody(t.Context(), "acme", "producer", src, json.RawMessage(schema))
		assertValidation(t, err, "config_schema")
		t.Logf("%s rejected as: %v", name, err)
	}
	for _, schema := range []string{`{}`, `{"type":"object","properties":{"x":{"type":"number"}}}`} {
		if err := s.validateUserBlockBody(t.Context(), "acme", "producer", src, json.RawMessage(schema)); err != nil {
			t.Errorf("config_schema %s must pass: %v", schema, err)
		}
	}
	if err := s.validateUserBlockBody(t.Context(), "acme", "producer", src, nil); err != nil {
		t.Errorf("nil config_schema must pass: %v", err)
	}
}

func TestCreateUserBlock_RejectsUnparsableSource(t *testing.T) {
	s := newUserBlockTestService()
	_, err := s.CreateUserBlock(t.Context(), "acme", "combo", "producer", json.RawMessage(`{"blocks": [`), nil)
	assertValidation(t, err, "does not parse")
}
