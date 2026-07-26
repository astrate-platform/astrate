package realm

import (
	"fmt"

	"github.com/astrate-platform/astrate/internal/engine/triggers"
)

// validatePolicy checks a delivery-policy document by delegating to
// triggers.CompilePolicy. Violations wrap ErrValidation (mapped to 422).
func validatePolicy(def []byte) (string, error) {
	p, err := triggers.CompilePolicy(def)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return p.Name, nil
}
