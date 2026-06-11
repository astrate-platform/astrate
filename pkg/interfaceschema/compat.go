package interfaceschema

import (
	"errors"
	"fmt"
)

// ErrIncompatibleUpgrade is wrapped by every CheckMinorUpgrade rejection, so
// callers can classify upgrade failures with errors.Is.
var ErrIncompatibleUpgrade = errors.New("incompatible interface upgrade")

// CheckMinorUpgrade enforces the Astarte minor-version compatibility rule
// exactly as Realm Management does on interface update (docs/DESIGN.md §2.6
// versioning parity): next must keep the same name, major version, type,
// ownership, and aggregation; strictly increase the minor version; keep
// every existing mapping with identical attributes (description and doc may
// change, and placeholders may be renamed — they are not semantic); and only
// add mappings, never remove them.
//
// Both arguments must be validated interfaces (from ParseInterface).
func CheckMinorUpgrade(old, next *Interface) error {
	if old == nil || next == nil {
		return fmt.Errorf("%w: nil interface", ErrIncompatibleUpgrade)
	}
	if next.Name != old.Name {
		return fmt.Errorf("%w: name changed from %q to %q", ErrIncompatibleUpgrade, old.Name, next.Name)
	}
	if next.Major != old.Major {
		return fmt.Errorf("%w: major version changed from %d to %d (new majors are new interfaces)",
			ErrIncompatibleUpgrade, old.Major, next.Major)
	}
	if next.Minor <= old.Minor {
		return fmt.Errorf("%w: minor version must increase (%d.%d -> %d.%d)",
			ErrIncompatibleUpgrade, old.Major, old.Minor, next.Major, next.Minor)
	}
	if next.Type != old.Type {
		return fmt.Errorf("%w: type changed from %s to %s", ErrIncompatibleUpgrade, old.Type, next.Type)
	}
	if next.Ownership != old.Ownership {
		return fmt.Errorf("%w: ownership changed from %s to %s", ErrIncompatibleUpgrade, old.Ownership, next.Ownership)
	}
	if next.Aggregation != old.Aggregation {
		return fmt.Errorf("%w: aggregation changed from %s to %s",
			ErrIncompatibleUpgrade, old.Aggregation, next.Aggregation)
	}

	nextByEndpoint := make(map[string]*Mapping, len(next.Mappings))
	for i := range next.Mappings {
		key, err := normalizedKey(next.Mappings[i].Endpoint)
		if err != nil {
			return err
		}
		nextByEndpoint[key] = &next.Mappings[i]
	}
	for i := range old.Mappings {
		om := &old.Mappings[i]
		key, err := normalizedKey(om.Endpoint)
		if err != nil {
			return err
		}
		nm, ok := nextByEndpoint[key]
		if !ok {
			return fmt.Errorf("%w: mapping %q removed (minor upgrades are additive only)",
				ErrIncompatibleUpgrade, om.Endpoint)
		}
		if err := sameMappingAttributes(om, nm); err != nil {
			return err
		}
	}
	return nil
}

// normalizedKey is the placeholder-erased endpoint form used to pair
// mappings across versions.
func normalizedKey(endpoint string) (string, error) {
	segs, err := splitEndpoint(endpoint)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrIncompatibleUpgrade, err)
	}
	return normalizeEndpoint(segs), nil
}

// sameMappingAttributes rejects any attribute mutation between two versions
// of the same mapping. Description and doc are non-semantic and may change.
func sameMappingAttributes(om, nm *Mapping) error {
	attr := ""
	switch {
	case nm.Type != om.Type:
		attr = "type"
	case nm.Reliability != om.Reliability:
		attr = "reliability"
	case nm.Retention != om.Retention:
		attr = "retention"
	case nm.Expiry != om.Expiry:
		attr = "expiry"
	case nm.DatabaseRetentionPolicy != om.DatabaseRetentionPolicy:
		attr = "database_retention_policy"
	case nm.DatabaseRetentionTTL != om.DatabaseRetentionTTL:
		attr = "database_retention_ttl"
	case nm.AllowUnset != om.AllowUnset:
		attr = "allow_unset"
	case nm.ExplicitTimestamp != om.ExplicitTimestamp:
		attr = "explicit_timestamp"
	}
	if attr != "" {
		return fmt.Errorf("%w: mapping %q changed %s (existing mappings are immutable)",
			ErrIncompatibleUpgrade, om.Endpoint, attr)
	}
	return nil
}
