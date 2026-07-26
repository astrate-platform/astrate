package triggers

import (
	"encoding/json"
	"fmt"
)

// CompileCondition compiles a single simple_triggers condition into a
// matcher-only Trigger. A transient trigger is compiled per subscription
// for live-event watching: it is never stored, never delivered, and its
// Action is deliberately nil.
func CompileCondition(name string, raw json.RawMessage) (*Trigger, error) {
	var c simpleTriggerConfig
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, fmt.Errorf("triggers: %q condition does not parse: %w", name, err)
	}
	t := &Trigger{Name: name}
	if err := t.compileSimple(&c); err != nil {
		return nil, fmt.Errorf("triggers: %q condition: %w", name, err)
	}
	return t, nil
}
