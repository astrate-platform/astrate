package broker

import (
	"errors"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

func TestParseCN(t *testing.T) {
	devID, err := deviceid.Parse("h4-Dx_RYTU-RbpDOTabhRg")
	if err != nil {
		t.Fatalf("parsing fixture device ID: %v", err)
	}

	valid := []struct {
		cn    string
		realm string
	}{
		{"test/h4-Dx_RYTU-RbpDOTabhRg", "test"},
		{"a/h4-Dx_RYTU-RbpDOTabhRg", "a"},
		{"realm42/h4-Dx_RYTU-RbpDOTabhRg", "realm42"},
		{strings.Repeat("r", 48) + "/h4-Dx_RYTU-RbpDOTabhRg", strings.Repeat("r", 48)},
	}
	for _, tc := range valid {
		id, err := ParseCN(tc.cn)
		if err != nil {
			t.Errorf("ParseCN(%q): unexpected error: %v", tc.cn, err)
			continue
		}
		if id.Realm != tc.realm || id.DeviceID != devID {
			t.Errorf("ParseCN(%q) = %+v, want realm %q device %s", tc.cn, id, tc.realm, devID)
		}
		if got := id.CN(); got != tc.cn {
			t.Errorf("CN() round trip: got %q, want %q", got, tc.cn)
		}
		if got := id.BaseTopic(); got != tc.cn {
			t.Errorf("BaseTopic(): got %q, want %q", got, tc.cn)
		}
	}

	invalid := []string{
		"",                                // empty
		"test",                            // no separator
		"/h4-Dx_RYTU-RbpDOTabhRg",         // empty realm
		"test/",                           // empty device
		"Test/h4-Dx_RYTU-RbpDOTabhRg",     // uppercase realm
		"9test/h4-Dx_RYTU-RbpDOTabhRg",    // digit-leading realm
		"te-st/h4-Dx_RYTU-RbpDOTabhRg",    // hyphen in realm
		"te st/h4-Dx_RYTU-RbpDOTabhRg",    // space in realm
		"te\x00st/h4-Dx_RYTU-RbpDOTabhRg", // NUL in realm
		"t#st/h4-Dx_RYTU-RbpDOTabhRg",     // wildcard char in realm
		strings.Repeat("r", 49) + "/h4-Dx_RYTU-RbpDOTabhRg", // realm too long
		"test/h4-Dx_RYTU-RbpDOTabhR",                        // 21-char device
		"test/h4-Dx_RYTU-RbpDOTabhRgg",                      // 23-char device
		"test/h4+Dx/RYTU-RbpDOTabhRg",                       // std-base64 chars + extra slash
		"test/h4-Dx_RYTU-RbpDOTabhRg==",                     // padded
		"test/h4-Dx_RYTU-RbpDOTabhRg/extra",                 // trailing segment
		"test/aaaaaaaaaaaaaaaaaaaaa!",                       // invalid base64url byte
	}
	for _, cn := range invalid {
		if _, err := ParseCN(cn); !errors.Is(err, ErrBadCN) {
			t.Errorf("ParseCN(%q): got %v, want ErrBadCN", cn, err)
		}
	}
}

func TestSplitTopic(t *testing.T) {
	valid := []struct {
		topic               string
		realm, device, rest string
	}{
		{"test/h4-Dx_RYTU-RbpDOTabhRg", "test", "h4-Dx_RYTU-RbpDOTabhRg", ""},
		{"test/h4-Dx_RYTU-RbpDOTabhRg/control/emptyCache", "test", "h4-Dx_RYTU-RbpDOTabhRg", "control/emptyCache"},
		{"test/h4-Dx_RYTU-RbpDOTabhRg/com.ex.Sensors/0/value", "test", "h4-Dx_RYTU-RbpDOTabhRg", "com.ex.Sensors/0/value"},
	}
	for _, tc := range valid {
		realm, device, rest, err := SplitTopic(tc.topic)
		if err != nil {
			t.Errorf("SplitTopic(%q): unexpected error: %v", tc.topic, err)
			continue
		}
		if realm != tc.realm || device != tc.device || rest != tc.rest {
			t.Errorf("SplitTopic(%q) = (%q, %q, %q), want (%q, %q, %q)",
				tc.topic, realm, device, rest, tc.realm, tc.device, tc.rest)
		}
	}

	invalid := []string{
		"",
		"test",
		"test/",
		"test/short",
		"Test/h4-Dx_RYTU-RbpDOTabhRg",
		"/h4-Dx_RYTU-RbpDOTabhRg",
	}
	for _, topic := range invalid {
		if _, _, _, err := SplitTopic(topic); !errors.Is(err, ErrBadTopic) {
			t.Errorf("SplitTopic(%q): got %v, want ErrBadTopic", topic, err)
		}
	}
}
