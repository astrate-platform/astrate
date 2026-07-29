package flow

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// configPlaceholder matches ${config.KEY} with KEY = [A-Za-z_][A-Za-z0-9_]*.
var configPlaceholder = regexp.MustCompile(`\$\{config\.([A-Za-z_][A-Za-z0-9_]*)\}`)

// SubstituteConfig walks a pipeline definition JSON and replaces
// ${config.key} placeholders inside string values only. Missing keys or
// non-stringable config values fail loudly. Non-string JSON leaves are
// left unchanged.
func SubstituteConfig(definition []byte, config map[string]any) ([]byte, error) {
	if len(definition) == 0 {
		return nil, fmt.Errorf("flow: empty pipeline definition")
	}
	if config == nil {
		config = map[string]any{}
	}
	var root any
	if err := json.Unmarshal(definition, &root); err != nil {
		return nil, fmt.Errorf("flow: pipeline definition does not parse: %w", err)
	}
	if err := walkSubstitute(root, config); err != nil {
		return nil, err
	}
	out, err := json.Marshal(root)
	if err != nil {
		return nil, fmt.Errorf("flow: re-encoding substituted definition: %w", err)
	}
	return out, nil
}

func walkSubstitute(v any, config map[string]any) error {
	switch n := v.(type) {
	case map[string]any:
		for k, child := range n {
			if s, ok := child.(string); ok {
				replaced, err := replacePlaceholders(s, config)
				if err != nil {
					return err
				}
				n[k] = replaced
				continue
			}
			if err := walkSubstitute(child, config); err != nil {
				return err
			}
		}
	case []any:
		for i, child := range n {
			if s, ok := child.(string); ok {
				replaced, err := replacePlaceholders(s, config)
				if err != nil {
					return err
				}
				n[i] = replaced
				continue
			}
			if err := walkSubstitute(child, config); err != nil {
				return err
			}
		}
	}
	return nil
}

func replacePlaceholders(s string, config map[string]any) (string, error) {
	var firstErr error
	out := configPlaceholder.ReplaceAllStringFunc(s, func(match string) string {
		if firstErr != nil {
			return match
		}
		sub := configPlaceholder.FindStringSubmatch(match)
		if len(sub) != 2 {
			firstErr = fmt.Errorf("flow: invalid config placeholder %q", match)
			return match
		}
		key := sub[1]
		val, ok := config[key]
		if !ok {
			firstErr = fmt.Errorf("flow: missing config key %q for placeholder ${config.%s}", key, key)
			return match
		}
		str, err := configValueString(val)
		if err != nil {
			firstErr = fmt.Errorf("flow: config key %q: %w", key, err)
			return match
		}
		return str
	})
	if firstErr != nil {
		return "", firstErr
	}
	// Any remaining ${config...} that did not match the key regex is also an error.
	if strings.Contains(out, "${config.") {
		return "", fmt.Errorf("flow: unresolved or invalid config placeholder in %q", s)
	}
	return out, nil
}

func configValueString(v any) (string, error) {
	switch t := v.(type) {
	case string:
		return t, nil
	case float64:
		// JSON numbers decode as float64.
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10), nil
		}
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(t), nil
	case json.Number:
		return t.String(), nil
	case nil:
		return "", fmt.Errorf("value is null (stringable string/number/bool required)")
	default:
		return "", fmt.Errorf("value has type %T (stringable string/number/bool required)", v)
	}
}
