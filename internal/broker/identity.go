package broker

import (
	"errors"
	"fmt"
	"strings"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// maxRealmNameLen caps realm names (upstream parity: realm names become
// Cassandra keyspace names there, capped at 48 characters).
const maxRealmNameLen = 48

// Sentinel identity/topic errors.
var (
	// ErrBadCN reports a certificate CN (or claimed client ID) that is not
	// a well-formed "<realm>/<device_id>".
	ErrBadCN = errors.New("broker: malformed identity CN")
	// ErrBadTopic reports a topic that does not start with the
	// "<realm>/<device_id>" prefix scheme (docs/DESIGN.md §3.3).
	ErrBadTopic = errors.New("broker: malformed device topic")
)

// Identity is an authenticated device: the parsed form of the certificate
// CN "<realm>/<device_id>" (docs/DESIGN.md §3.1, §4.3).
type Identity struct {
	// Realm is the device's realm name.
	Realm string
	// DeviceID is the 128-bit Astarte device ID.
	DeviceID deviceid.ID
}

// CN renders the identity back to its canonical "<realm>/<device_id>" form —
// the certificate CN, the required MQTT client ID, and the device's base
// topic, which are all the same string by design.
func (i Identity) CN() string {
	return i.Realm + "/" + i.DeviceID.String()
}

// BaseTopic returns the device's base MQTT topic B = "<realm>/<device_id>"
// (docs/DESIGN.md §3.2). Identical to CN; named separately for call-site
// clarity.
func (i Identity) BaseTopic() string {
	return i.CN()
}

// ParseCN parses a certificate CN (equivalently: a claimed MQTT client ID)
// of the form "<realm>/<device_id>". The realm must be a valid Astarte realm
// name and the device part a 22-character base64url 128-bit device ID;
// anything else fails with ErrBadCN.
func ParseCN(cn string) (Identity, error) {
	realm, dev, ok := strings.Cut(cn, "/")
	if !ok {
		return Identity{}, fmt.Errorf("%w: %q has no realm/device separator", ErrBadCN, cn)
	}
	if !validRealmName(realm) {
		return Identity{}, fmt.Errorf("%w: invalid realm name in %q", ErrBadCN, cn)
	}
	id, err := deviceid.Parse(dev)
	if err != nil {
		return Identity{}, fmt.Errorf("%w: invalid device ID in %q: %v", ErrBadCN, cn, err)
	}
	return Identity{Realm: realm, DeviceID: id}, nil
}

// SplitTopic splits a wire topic into its realm, device, and rest parts
// (docs/DESIGN.md §3.3 parsing note): "<realm>/<device_id>" yields an empty
// rest (the introspection topic); "<realm>/<device_id>/<rest>" yields the
// remainder verbatim (control suffix or "<interface_name><path>"). The
// realm and device segments are shape-checked; the rest is not interpreted.
func SplitTopic(topic string) (realm, device, rest string, err error) {
	realm, tail, ok := strings.Cut(topic, "/")
	if !ok || !validRealmName(realm) {
		return "", "", "", fmt.Errorf("%w: %q lacks a valid realm segment", ErrBadTopic, topic)
	}
	device, rest, _ = strings.Cut(tail, "/")
	if len(device) != deviceid.EncodedLen {
		return "", "", "", fmt.Errorf("%w: %q lacks a device ID segment", ErrBadTopic, topic)
	}
	return realm, device, rest, nil
}

// validRealmName checks the Astarte realm-name shape: a lowercase letter
// followed by lowercase letters or digits (upstream housekeeping pattern
// ^[a-z][a-z0-9]*$), at most maxRealmNameLen characters.
func validRealmName(s string) bool {
	if len(s) == 0 || len(s) > maxRealmNameLen {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z':
		case c >= '0' && c <= '9':
			if i == 0 {
				return false
			}
		default:
			return false
		}
	}
	return true
}
