package pairing

import (
	"net/http"
	"net/netip"
	"testing"
)

func TestRemoteIP(t *testing.T) {
	tests := []struct {
		name   string
		remote string
		want   netip.Addr
	}{
		// valid addr:port — the common http.Server RemoteAddr form
		{"ipv4 addr:port", "192.0.2.1:5678", netip.MustParseAddr("192.0.2.1")},

		// IPv4-in-IPv6 (RFC 4291 §2.5.5) is unmapped back to IPv4
		{"ipv4-mapped addr:port", "[::ffff:192.0.2.1]:80", netip.MustParseAddr("192.0.2.1")},

		// bare address with no port
		{"bare ipv4", "192.0.2.1", netip.MustParseAddr("192.0.2.1")},
		{"bare ipv4-mapped", "::ffff:192.0.2.1", netip.MustParseAddr("192.0.2.1")},

		// unparseable RemoteAddr (unix socket / proxy) → IPv4 unspecified
		{"unix socket", "/var/run/pairing.sock", netip.IPv4Unspecified()},
		{"empty", "", netip.IPv4Unspecified()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.remote
			if got := remoteIP(req); got != tt.want {
				t.Errorf("remoteIP(%q) = %v, want %v", tt.remote, got, tt.want)
			}
		})
	}
}
