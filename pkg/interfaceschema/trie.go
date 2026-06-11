package interfaceschema

import (
	"fmt"
	"strings"
	"time"
)

// MaxPlaceholderValueLen is the maximum byte length of a concrete path
// segment matched by a %{placeholder} (docs/DESIGN.md §2.6 hygiene rule).
const MaxPlaceholderValueLen = 256

// CompiledMapping is the hot-path form of one endpoint mapping
// (docs/DESIGN.md §2.6). It is what the engine validates and persists
// against after a trie match.
type CompiledMapping struct {
	// EndpointID is the storage identifier of this endpoint.
	EndpointID int64
	// ValueType drives both validation and BSON/JSON decoding.
	ValueType ValueType
	// Reliability is the MQTT QoS byte (0, 1, or 2).
	Reliability byte
	// Retention states what happens to undeliverable values.
	Retention Retention
	// Expiry is the retention expiry; 0 means never.
	Expiry time.Duration
	// ExplicitTimestamp states whether publishes carry their own timestamp.
	ExplicitTimestamp bool
	// AllowUnset permits property unset via empty payload.
	AllowUnset bool
	// DBRetentionTTL is the database TTL; 0 means no_ttl.
	DBRetentionTTL time.Duration
}

// EndpointTrie matches concrete inbound paths (for example "/4/value")
// against declared endpoint patterns (for example "/%{sensor_id}/value"),
// segment by segment. Exact-literal children take priority over the (single)
// parametric child; placeholder values are charset- and length-checked but
// not semantically interpreted. Match is O(depth) and allocation-free.
type EndpointTrie struct {
	root trieNode
}

// trieNode is one level of the trie. A node may simultaneously be a leaf
// (mapping set) and an interior node (children/param set): "/a" and "/a/b"
// coexist but match paths of different depths.
type trieNode struct {
	children  map[string]*trieNode
	param     *trieNode
	paramName string
	mapping   *CompiledMapping
}

// NewEndpointTrie returns an empty trie.
func NewEndpointTrie() *EndpointTrie {
	return &EndpointTrie{}
}

// Add inserts one endpoint pattern with its compiled mapping. The endpoint
// must be syntactically valid (same rules as ParseInterface); duplicate
// endpoints and sibling placeholders with different names are rejected.
// Cross-pattern conflict checking (ambiguous literal/parametric overlap) is
// ParseInterface's responsibility — the trie itself resolves such overlaps
// deterministically, exact match first.
func (t *EndpointTrie) Add(endpoint string, m *CompiledMapping) error {
	if m == nil {
		return fmt.Errorf("interfaceschema: Add(%q): nil mapping", endpoint)
	}
	segs, err := splitEndpoint(endpoint)
	if err != nil {
		return fmt.Errorf("interfaceschema: %w", err)
	}
	node := &t.root
	for _, seg := range segs {
		if seg.param {
			if node.param == nil {
				node.param = &trieNode{}
				node.paramName = seg.literal
			} else if node.paramName != seg.literal {
				return fmt.Errorf(
					"interfaceschema: Add(%q): conflicting placeholder %%{%s} already declared as %%{%s}",
					endpoint, seg.literal, node.paramName)
			}
			node = node.param
			continue
		}
		if node.children == nil {
			node.children = make(map[string]*trieNode)
		}
		child, ok := node.children[seg.literal]
		if !ok {
			child = &trieNode{}
			node.children[seg.literal] = child
		}
		node = child
	}
	if node.mapping != nil {
		return fmt.Errorf("interfaceschema: Add(%q): duplicate endpoint", endpoint)
	}
	node.mapping = m
	return nil
}

// Match resolves a concrete path to its compiled mapping. The path must be
// '/'-rooted with non-empty segments; parametric segments must additionally
// be ≤ 256 bytes and contain no '+' or '#'. Match never allocates.
func (t *EndpointTrie) Match(path string) (*CompiledMapping, bool) {
	if t == nil || len(path) < 2 || path[0] != '/' {
		return nil, false
	}
	m := matchNode(&t.root, path[1:])
	return m, m != nil
}

// matchNode matches rest (the remaining path, leading '/' stripped) from n,
// preferring exact-literal children and backtracking into the parametric
// child when the literal branch dead-ends.
func matchNode(n *trieNode, rest string) *CompiledMapping {
	var seg, tail string
	last := false
	if i := strings.IndexByte(rest, '/'); i < 0 {
		seg, last = rest, true
	} else {
		seg, tail = rest[:i], rest[i+1:]
	}
	if seg == "" {
		return nil // empty segment: "//", trailing '/', or bare "/"
	}
	if child, ok := n.children[seg]; ok {
		if last {
			if child.mapping != nil {
				return child.mapping
			}
		} else if m := matchNode(child, tail); m != nil {
			return m
		}
	}
	if n.param != nil && placeholderValueOK(seg) {
		if last {
			return n.param.mapping
		}
		return matchNode(n.param, tail)
	}
	return nil
}

// placeholderValueOK applies the §2.6 hygiene rules to a concrete segment
// consumed by a placeholder: non-empty (checked by the caller), at most 256
// bytes, and free of MQTT wildcard bytes. It must not allocate.
func placeholderValueOK(seg string) bool {
	if len(seg) > MaxPlaceholderValueLen {
		return false
	}
	return strings.IndexByte(seg, '+') < 0 && strings.IndexByte(seg, '#') < 0
}
