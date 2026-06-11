package interfaceschema

import (
	"fmt"
	"time"
)

// EndpointIDResolver supplies the storage identifiers stamped into a
// CompiledInterface: the interface row ID and one ID per declared endpoint.
// The store layer implements it against the interfaces/endpoints tables; a
// nil resolver compiles with all IDs zero (pure in-memory use, tests).
type EndpointIDResolver interface {
	// ResolveInterface returns the storage ID of the interface itself.
	ResolveInterface(name string, major int) (int64, error)
	// ResolveEndpoint returns the storage ID of one declared endpoint
	// pattern (for example "/%{sensor_id}/value").
	ResolveEndpoint(endpoint string) (int64, error)
}

// CompiledInterface is the hot-path form of a validated interface
// (docs/DESIGN.md §2.6): an endpoint trie for individual matching plus, for
// object aggregation, the last-level key → mapping table.
type CompiledInterface struct {
	// ID is the storage identifier of the interface.
	ID int64
	// Name is the reverse-domain interface name.
	Name string
	// Major and Minor are the interface version.
	Major, Minor int
	// Type discriminates datastream from properties.
	Type InterfaceType
	// Ownership states which side publishes on this interface.
	Ownership Ownership
	// Aggregation is individual or object.
	Aggregation Aggregation
	// Trie matches concrete full endpoint paths, segment-wise.
	Trie *EndpointTrie
	// ObjectLeaves maps, for object aggregation, each last-level name to its
	// mapping; nil for individual aggregation.
	ObjectLeaves map[string]*CompiledMapping
}

// Compile turns a validated Interface into its hot-path form. The input must
// come from ParseInterface (or satisfy the same invariants); endpoint or
// aggregation violations surface as errors, not panics. With a nil resolver
// every ID is zero.
func Compile(iface *Interface, ids EndpointIDResolver) (*CompiledInterface, error) {
	if iface == nil {
		return nil, fmt.Errorf("interfaceschema: Compile(nil interface)")
	}
	ci := &CompiledInterface{
		Name:        iface.Name,
		Major:       iface.Major,
		Minor:       iface.Minor,
		Type:        iface.Type,
		Ownership:   iface.Ownership,
		Aggregation: iface.Aggregation,
		Trie:        NewEndpointTrie(),
	}
	if ids != nil {
		id, err := ids.ResolveInterface(iface.Name, iface.Major)
		if err != nil {
			return nil, fmt.Errorf("interfaceschema: resolving interface %s v%d: %w", iface.Name, iface.Major, err)
		}
		ci.ID = id
	}
	if iface.Aggregation == AggregationObject {
		ci.ObjectLeaves = make(map[string]*CompiledMapping, len(iface.Mappings))
	}

	for i := range iface.Mappings {
		m := &iface.Mappings[i]
		cm := &CompiledMapping{
			ValueType:         m.Type,
			Reliability:       m.Reliability.QoS(),
			Retention:         m.Retention,
			Expiry:            time.Duration(m.Expiry) * time.Second,
			ExplicitTimestamp: m.ExplicitTimestamp,
			AllowUnset:        m.AllowUnset,
			DBRetentionTTL:    time.Duration(m.DatabaseRetentionTTL) * time.Second,
		}
		if ids != nil {
			id, err := ids.ResolveEndpoint(m.Endpoint)
			if err != nil {
				return nil, fmt.Errorf("interfaceschema: resolving endpoint %q: %w", m.Endpoint, err)
			}
			cm.EndpointID = id
		}
		if err := ci.Trie.Add(m.Endpoint, cm); err != nil {
			return nil, err
		}
		if ci.ObjectLeaves != nil {
			leaf, err := lastLiteralSegment(m.Endpoint)
			if err != nil {
				return nil, err
			}
			if _, dup := ci.ObjectLeaves[leaf]; dup {
				return nil, fmt.Errorf("interfaceschema: duplicate object leaf %q in %s", leaf, iface.Name)
			}
			ci.ObjectLeaves[leaf] = cm
		}
	}
	return ci, nil
}

// lastLiteralSegment returns the final endpoint level, which for
// object-aggregated interfaces must be a literal (it is the object document
// key).
func lastLiteralSegment(endpoint string) (string, error) {
	segs, err := splitEndpoint(endpoint)
	if err != nil {
		return "", fmt.Errorf("interfaceschema: %w", err)
	}
	last := segs[len(segs)-1]
	if last.param {
		return "", fmt.Errorf(
			"interfaceschema: object-aggregated mapping %q must end in a literal segment", endpoint)
	}
	return last.literal, nil
}
