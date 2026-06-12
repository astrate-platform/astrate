package engine

import (
	"testing"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// testDeviceID is a fixed device for topic tests; its String() form is the
// 22-character base64url encoding.
var testDeviceID = deviceid.ID{0xde, 0xad, 0xbe, 0xef, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

func TestSplitDeviceTopic(t *testing.T) {
	dev := testDeviceID.String()
	other := deviceid.ID{1}.String()

	cases := []struct {
		name     string
		topic    string
		wantRest string
		wantOK   bool
	}{
		{"introspection root", "alpha/" + dev, "", true},
		{"control", "alpha/" + dev + "/control/emptyCache", "control/emptyCache", true},
		{"data", "alpha/" + dev + "/com.ex.S/4/value", "com.ex.S/4/value", true},
		{"trailing slash", "alpha/" + dev + "/", "", true},
		{"wrong realm", "beta/" + dev, "", false},
		{"wrong device", "alpha/" + other, "", false},
		{"device id prefix glued", "alpha/" + dev + "x/y", "", false},
		{"too short", "alpha/abc", "", false},
		{"no separator after realm", "alpha" + dev, "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rest, ok := splitDeviceTopic(tc.topic, "alpha", testDeviceID)
			if ok != tc.wantOK || rest != tc.wantRest {
				t.Errorf("splitDeviceTopic(%q) = (%q, %v), want (%q, %v)",
					tc.topic, rest, ok, tc.wantRest, tc.wantOK)
			}
		})
	}
}

func TestClassify(t *testing.T) {
	cases := []struct {
		rest        string
		wantKind    topicKind
		wantSubpath string
	}{
		{"", kindIntrospection, ""},
		{"control", kindControl, ""},
		{"control/emptyCache", kindControl, "emptyCache"},
		{"control/producer/properties", kindControl, "producer/properties"},
		{"control/", kindControl, ""},
		{"controlfreak/x", kindData, ""},
		{"com.ex.Sensors/4/value", kindData, ""},
		{"com.ex.Sensors", kindData, ""},
	}
	for _, tc := range cases {
		kind, subpath := classify(tc.rest)
		if kind != tc.wantKind || subpath != tc.wantSubpath {
			t.Errorf("classify(%q) = (%d, %q), want (%d, %q)",
				tc.rest, kind, subpath, tc.wantKind, tc.wantSubpath)
		}
	}
}

func TestMatchInterface(t *testing.T) {
	declared := map[string]store.InterfaceVersion{
		"com.ex.Sensors": {Major: 1, Minor: 2},
		"com.ex":         {Major: 0, Minor: 1},
		"single":         {Major: 3},
	}
	declares := func(name string) (store.InterfaceVersion, bool) {
		v, ok := declared[name]
		return v, ok
	}

	cases := []struct {
		rest      string
		wantName  string
		wantPath  string
		wantMajor int
		wantOK    bool
	}{
		{"com.ex.Sensors/4/value", "com.ex.Sensors", "/4/value", 1, true},
		{"com.ex.Sensors", "com.ex.Sensors", "", 1, true},
		{"com.ex/x", "com.ex", "/x", 0, true},
		// Longest declared prefix wins even when a shorter one also exists.
		{"com.ex.Sensors/value", "com.ex.Sensors", "/value", 1, true},
		{"single/a/b/c", "single", "/a/b/c", 3, true},
		{"com.other.Iface/x", "", "", 0, false},
		{"/leading/slash", "", "", 0, false},
		{"", "", "", 0, false},
	}
	for _, tc := range cases {
		name, path, ver, ok := matchInterface(tc.rest, declares)
		if ok != tc.wantOK || name != tc.wantName || path != tc.wantPath {
			t.Errorf("matchInterface(%q) = (%q, %q, ok=%v), want (%q, %q, ok=%v)",
				tc.rest, name, path, ok, tc.wantName, tc.wantPath, tc.wantOK)
			continue
		}
		if ok && ver.Major != tc.wantMajor {
			t.Errorf("matchInterface(%q) major = %d, want %d", tc.rest, ver.Major, tc.wantMajor)
		}
	}
}
