package interfaceschema

import (
	"fmt"
	"strings"
)

// Violation is one upstream-shaped field rejection. MappingIndex < 0 marks
// an interface-level field; >= 0 marks a field of mappings[MappingIndex].
type Violation struct {
	Field        string
	Messages     []string
	MappingIndex int
}

// ViolationsError carries field-shaped rejections for the rules whose
// upstream wire envelopes were probe-verified against Astarte 1.2.0
// (issue #61). ParseInterface returns it wrapped together with ErrInvalid,
// so errors.Is(err, ErrInvalid) and errors.As(err, &ve) both hit on the
// same value.
type ViolationsError struct {
	Violations []Violation
	// MappingCount is the number of declared mappings, needed to render
	// upstream's full-length aligned "mappings" error array.
	MappingCount int
}

// violationError couples the structured violations with the ErrInvalid
// sentinel in a single error value: Error() renders exactly the violations'
// own message (once), while Unwrap exposes both ErrInvalid and the
// *ViolationsError for errors.Is / errors.As.
type violationError struct {
	ve *ViolationsError
}

func (e *violationError) Error() string { return e.ve.Error() }

func (e *violationError) Unwrap() []error { return []error{ErrInvalid, e.ve} }

// add appends an interface-level violation.
func (e *ViolationsError) add(field string, msgs ...string) {
	e.Violations = append(e.Violations, Violation{Field: field, Messages: msgs, MappingIndex: -1})
}

// addMapping appends a violation against mappings[index].
func (e *ViolationsError) addMapping(index int, field string, msgs ...string) {
	e.Violations = append(e.Violations, Violation{Field: field, Messages: msgs, MappingIndex: index})
}

// empty reports whether no violation was collected yet.
func (e *ViolationsError) empty() bool { return len(e.Violations) == 0 }

// errOrNil returns the combined error value when any violation was collected
// and nil otherwise, so callers can abort a parse with a single check.
func (e *ViolationsError) errOrNil() error {
	if e.empty() {
		return nil
	}
	return &violationError{ve: e}
}

// Error renders every collected rejection on one line with the package's
// standard rejection prefix, e.g.
//
//	invalid interface: description: should be at most 1000 character(s); mappings[1].doc: ...
func (e *ViolationsError) Error() string {
	var b strings.Builder
	b.WriteString("invalid interface: ")
	for i := range e.Violations {
		if i > 0 {
			b.WriteString("; ")
		}
		v := &e.Violations[i]
		if v.MappingIndex >= 0 {
			fmt.Fprintf(&b, "mappings[%d].", v.MappingIndex)
		}
		b.WriteString(v.Field)
		b.WriteString(": ")
		b.WriteString(strings.Join(v.Messages, ", "))
	}
	return b.String()
}
