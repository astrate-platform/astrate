package appengine

import (
	"net/url"
	"strings"
	"testing"
)

// TestDeviceListLinks pins the pagination-links contract the dashboard
// depends on: it parses links.next with JavaScript's
// `new URLSearchParams(next)` over the RAW path+query string, which corrupts
// the FIRST key/value pair ("/v1/…?details" becomes the first key). from_token
// must therefore never be the first query parameter.
func TestDeviceListLinks(t *testing.T) {
	q := url.Values{"details": {"true"}, "limit": {"10"}}
	links := deviceListLinks("test", q, "CURSOR")

	if links.Self != "/v1/test/devices?details=true&limit=10" {
		t.Errorf("self = %q", links.Self)
	}
	if links.Next == "" {
		t.Fatal("next missing with a cursor present")
	}
	if strings.Contains(links.Next, "?from_token=") {
		t.Fatalf("from_token is the first query parameter in %q — the dashboard's URLSearchParams parse would corrupt it", links.Next)
	}

	// Recover the cursor exactly the way the dashboard does: URLSearchParams
	// over the raw string splits on & and =, mangling only the first pair.
	got := ""
	for i, pair := range strings.Split(links.Next, "&") {
		k, v, _ := strings.Cut(pair, "=")
		if i == 0 {
			continue // the "…?details" mangled pair
		}
		if k == "from_token" {
			got, _ = url.QueryUnescape(v)
		}
	}
	if got != "CURSOR" {
		t.Errorf("dashboard-style parse recovered from_token %q, want CURSOR", got)
	}
}

func TestDeviceListLinksLastPageAndBareQuery(t *testing.T) {
	links := deviceListLinks("test", url.Values{}, "")
	if links.Next != "" {
		t.Errorf("next = %q on the last page, want empty", links.Next)
	}
	if links.Self != "/v1/test/devices" {
		t.Errorf("self = %q", links.Self)
	}

	// With sparse request params the effective defaults are filled in, so
	// "details" always precedes "from_token" even when the client sent
	// neither.
	links = deviceListLinks("test", url.Values{"limit": {"5"}}, "N")
	if links.Next != "/v1/test/devices?details=false&from_token=N&limit=5" {
		t.Errorf("next = %q", links.Next)
	}
	links = deviceListLinks("test", url.Values{}, "N")
	if strings.Contains(links.Next, "?from_token=") {
		t.Errorf("from_token first in %q", links.Next)
	}
}
