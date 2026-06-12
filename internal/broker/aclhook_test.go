package broker

import (
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// TestCheckACLMatrix is the exhaustive §3.2 table test (docs/ROADMAP.md §6):
// every allowed row, plus adversarial topics that must all be denied. The
// device has introspected one device-owned and one server-owned interface.
func TestCheckACLMatrix(t *testing.T) {
	const (
		base  = "test/h4-Dx_RYTU-RbpDOTabhRg"
		other = "test/AAAAAAAAAAAAAAAAAAAAAA" // another device, same realm
	)
	ownership := func(iface string) (interfaceschema.Ownership, bool) {
		switch iface {
		case "com.ex.DeviceData":
			return interfaceschema.OwnershipDevice, true
		case "com.ex.ServerData":
			return interfaceschema.OwnershipServer, true
		}
		return 0, false
	}

	cases := []struct {
		name  string
		topic string
		write bool
		want  bool
	}{
		// --- PUBLISH: allowed rows ---
		{"pub introspection base", base, true, true},
		{"pub control emptyCache", base + "/control/emptyCache", true, true},
		{"pub control producer properties", base + "/control/producer/properties", true, true},
		{"pub device-owned path", base + "/com.ex.DeviceData/value", true, true},
		{"pub device-owned nested path", base + "/com.ex.DeviceData/0/sample/value", true, true},

		// --- PUBLISH: denied ---
		{"pub other device base", other, true, false},
		{"pub other device data", other + "/com.ex.DeviceData/value", true, false},
		{"pub other realm", "evil/h4-Dx_RYTU-RbpDOTabhRg/com.ex.DeviceData/value", true, false},
		{"pub realm-prefix trick", base + "x/com.ex.DeviceData/value", true, false},
		{"pub server-owned interface", base + "/com.ex.ServerData/value", true, false},
		{"pub uninstalled interface", base + "/com.ex.Unknown/value", true, false},
		{"pub interface without path", base + "/com.ex.DeviceData", true, false},
		{"pub interface with empty path", base + "/com.ex.DeviceData/", true, false},
		{"pub control consumer properties", base + "/control/consumer/properties", true, false},
		{"pub control unknown", base + "/control/selfDestruct", true, false},
		{"pub control prefix only", base + "/control", true, false},
		{"pub bare realm", "test", true, false},
		{"pub empty topic", "", true, false},
		{"pub sys topic", "$SYS/broker/uptime", true, false},
		{"pub wildcard hash", base + "/#", true, false},
		{"pub wildcard plus interface", base + "/+/value", true, false},
		{"pub oversized topic", base + "/com.ex.DeviceData/" + strings.Repeat("x", maxTopicBytes), true, false},

		// --- SUBSCRIBE / delivery: allowed rows ---
		{"sub control consumer properties", base + "/control/consumer/properties", false, true},
		{"sub superset hash", base + "/#", false, true},
		{"sub server-owned wildcard", base + "/com.ex.ServerData/#", false, true},
		{"sub server-owned inner wildcard", base + "/com.ex.ServerData/+/value", false, true},
		{"deliver server-owned concrete", base + "/com.ex.ServerData/value", false, true},
		{"deliver server-owned nested", base + "/com.ex.ServerData/0/sample", false, true},

		// --- SUBSCRIBE / delivery: denied ---
		{"sub base topic", base, false, false},
		{"sub other device hash", other + "/#", false, false},
		{"sub other device consumer properties", other + "/control/consumer/properties", false, false},
		{"sub other realm hash", "evil/h4-Dx_RYTU-RbpDOTabhRg/#", false, false},
		{"sub global hash", "#", false, false},
		{"sub global plus pair", "+/+", false, false},
		{"sub global plus hash", "+/#", false, false},
		{"sub realm hash", "test/#", false, false},
		{"sub device-owned interface", base + "/com.ex.DeviceData/#", false, false},
		{"deliver device-owned concrete", base + "/com.ex.DeviceData/value", false, false},
		{"sub uninstalled interface", base + "/com.ex.Unknown/#", false, false},
		{"sub control producer properties", base + "/control/producer/properties", false, false},
		{"sub control hash", base + "/control/#", false, false},
		{"sub wildcard interface segment", base + "/+/value", false, false},
		{"sub sys topics", "$SYS/#", false, false},
		{"sub empty topic", "", false, false},
		{"sub oversized filter", base + "/" + strings.Repeat("x", maxTopicBytes) + "/#", false, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := checkACL(base, tc.topic, tc.write, ownership); got != tc.want {
				t.Errorf("checkACL(%q, %q, write=%v) = %v, want %v",
					base, tc.topic, tc.write, got, tc.want)
			}
		})
	}
}

// TestCheckACLNoIntrospection covers a freshly-registered device that has
// not introspected yet: only the base topic and control publishes plus the
// control/superset subscriptions are allowed.
func TestCheckACLNoIntrospection(t *testing.T) {
	const base = "test/h4-Dx_RYTU-RbpDOTabhRg"
	none := func(string) (interfaceschema.Ownership, bool) { return 0, false }

	allowed := []struct {
		topic string
		write bool
	}{
		{base, true},
		{base + "/control/emptyCache", true},
		{base + "/control/producer/properties", true},
		{base + "/control/consumer/properties", false},
		{base + "/#", false},
	}
	for _, tc := range allowed {
		if !checkACL(base, tc.topic, tc.write, none) {
			t.Errorf("checkACL(%q, write=%v) = false, want true", tc.topic, tc.write)
		}
	}
	denied := []struct {
		topic string
		write bool
	}{
		{base + "/com.ex.DeviceData/value", true},
		{base + "/com.ex.ServerData/#", false},
	}
	for _, tc := range denied {
		if checkACL(base, tc.topic, tc.write, none) {
			t.Errorf("checkACL(%q, write=%v) = true, want false", tc.topic, tc.write)
		}
	}
}
