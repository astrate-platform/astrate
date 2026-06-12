package engine

import (
	"strings"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// Topic classification (docs/DESIGN.md §3.3): an inbound device topic is
// "<realm>/<device_id>[/<rest>]". rest "" is the introspection publish,
// "control/..." is the control channel, and anything else is data — split
// into interface name and path by longest-prefix match against the device's
// introspected interface names.

// topicKind discriminates the §3.3 message taxonomy.
type topicKind uint8

// topicKind values.
const (
	// kindIntrospection is the bare "<realm>/<device_id>" topic.
	kindIntrospection topicKind = iota + 1
	// kindControl is "<realm>/<device_id>/control/...".
	kindControl
	// kindData is a data publish on "<realm>/<device_id>/<interface><path>".
	kindData
)

// controlPrefix introduces control-channel topics under the device root.
const controlPrefix = "control/"

// deviceIDStringLen is the canonical encoded device-ID length (22 characters
// of unpadded base64url).
const deviceIDStringLen = 22

// splitDeviceTopic strips the "<realm>/<device_id>" root off topic and
// returns the remainder ("" for the introspection topic). ok is false when
// the topic does not carry that exact root — which the broker ACL prevents,
// so a mismatch is a defensive rejection, not an expected path.
func splitDeviceTopic(topic, realm string, id deviceid.ID) (rest string, ok bool) {
	root := len(realm) + 1 + deviceIDStringLen
	if len(topic) < root ||
		topic[:len(realm)] != realm ||
		topic[len(realm)] != '/' ||
		topic[len(realm)+1:root] != id.String() {
		return "", false
	}
	if len(topic) == root {
		return "", true
	}
	if topic[root] != '/' {
		return "", false
	}
	return topic[root+1:], true
}

// classify buckets rest into the §3.3 taxonomy. For control topics, subpath
// is the remainder after "control/" (e.g. "emptyCache",
// "producer/properties"); a bare "control" yields an empty subpath, which
// the control handler rejects.
func classify(rest string) (kind topicKind, subpath string) {
	switch {
	case rest == "":
		return kindIntrospection, ""
	case rest == "control":
		return kindControl, ""
	case strings.HasPrefix(rest, controlPrefix):
		return kindControl, rest[len(controlPrefix):]
	default:
		return kindData, ""
	}
}

// matchInterface resolves a data-topic rest into (interface name, path) by
// longest-prefix match over '/' boundaries against the device's introspected
// interface names (docs/DESIGN.md §3.3 parsing note). Interface names never
// contain '/', so in practice the first segment wins, but the longest-prefix
// contract is implemented as specified. path is "" or starts with '/'.
func matchInterface(rest string, declares func(string) (store.InterfaceVersion, bool)) (name, path string, ver store.InterfaceVersion, ok bool) {
	end := len(rest)
	for {
		if v, found := declares(rest[:end]); found {
			return rest[:end], rest[end:], v, true
		}
		i := strings.LastIndexByte(rest[:end], '/')
		if i < 0 {
			return "", "", store.InterfaceVersion{}, false
		}
		end = i
	}
}
