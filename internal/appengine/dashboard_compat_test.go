//go:build integration

package appengine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// pagedBody is the upstream paged-collection envelope the dashboard consumes.
type pagedBody struct {
	Data  json.RawMessage `json:"data"`
	Links struct {
		Self string `json:"self"`
		Next string `json:"next"`
	} `json:"links"`
}

// dashboardCursor extracts from_token the way the dashboard does:
// URLSearchParams over the raw links.next string, which mangles the first
// key/value pair.
func dashboardCursor(t *testing.T, next string) string {
	t.Helper()
	for i, pair := range strings.Split(next, "&") {
		if i == 0 {
			continue // "…?<firstkey>" — corrupted by URLSearchParams
		}
		if k, v, _ := strings.Cut(pair, "="); k == "from_token" {
			cursor, err := url.QueryUnescape(v)
			if err != nil {
				t.Fatalf("unescaping cursor: %v", err)
			}
			return cursor
		}
	}
	return ""
}

// TestDashboardCompat covers the M10 surface the Astarte Dashboard v1.2.2
// requires: /stats/devices, details=true listings, and body-links pagination.
func TestDashboardCompat(t *testing.T) {
	r := newRig(t)
	ctx := context.Background()

	// Four more devices → five total (one connected from the rig).
	for range 4 {
		id, err := deviceid.Random()
		if err != nil {
			t.Fatal(err)
		}
		if err := r.st.RegisterDevice(ctx, r.realmID, id, "h"); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("Stats", func(t *testing.T) {
		var stats struct {
			Connected int64 `json:"connected_devices"`
			Total     int64 `json:"total_devices"`
		}
		decodeData(t, r.req(t, http.MethodGet, "/stats/devices", "", r.token), &stats)
		if stats.Total != 5 || stats.Connected != 1 {
			t.Errorf("stats = %+v, want total 5 connected 1", stats)
		}
	})

	t.Run("PaginationWalk", func(t *testing.T) {
		// Walk ?details=true&limit=2 exactly like the dashboard: 3 pages,
		// cursor recovered from links.next, 5 distinct devices, no next on
		// the last page.
		seen := map[string]bool{}
		cursor := ""
		for page := 0; ; page++ {
			if page > 3 {
				t.Fatal("pagination did not terminate")
			}
			path := "/devices?details=true&limit=2"
			if cursor != "" {
				path += "&from_token=" + url.QueryEscape(cursor)
			}
			rec := r.req(t, http.MethodGet, path, "", r.token)
			if rec.Code != http.StatusOK {
				t.Fatalf("page %d: status %d: %s", page, rec.Code, rec.Body.String())
			}
			var body pagedBody
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			var statuses []DeviceStatus
			if err := json.Unmarshal(body.Data, &statuses); err != nil {
				t.Fatalf("page %d data is not a status array: %v", page, err)
			}
			for _, ds := range statuses {
				if seen[ds.ID] {
					t.Errorf("device %s appeared twice", ds.ID)
				}
				seen[ds.ID] = true
			}
			if body.Links.Self == "" {
				t.Error("links.self missing")
			}
			if body.Links.Next == "" {
				if len(statuses) == 2 && len(seen) < 5 {
					t.Fatalf("next missing before the end (%d seen)", len(seen))
				}
				break
			}
			if strings.Contains(body.Links.Next, "?from_token=") {
				t.Fatalf("from_token first in %q", body.Links.Next)
			}
			cursor = dashboardCursor(t, body.Links.Next)
			if cursor == "" {
				t.Fatalf("dashboard-style parse found no cursor in %q", body.Links.Next)
			}
		}
		if len(seen) != 5 {
			t.Errorf("walk found %d devices, want 5", len(seen))
		}
	})

	t.Run("PlainListStillIDs", func(t *testing.T) {
		rec := r.req(t, http.MethodGet, "/devices?limit=10", "", r.token)
		var body pagedBody
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		var ids []string
		if err := json.Unmarshal(body.Data, &ids); err != nil {
			t.Fatalf("data is not an ID array: %v (%s)", err, body.Data)
		}
		if len(ids) != 5 {
			t.Errorf("ids = %d, want 5", len(ids))
		}
	})

	t.Run("GroupDevicesDetails", func(t *testing.T) {
		if rec := r.req(t, http.MethodPost, "/groups",
			`{"group_name":"dash","devices":["`+r.dev.String()+`"]}`, r.token); rec.Code != http.StatusCreated {
			t.Fatalf("creating group: %d %s", rec.Code, rec.Body.String())
		}
		var statuses []DeviceStatus
		decodeData(t, r.req(t, http.MethodGet, "/groups/dash/devices?details=true", "", r.token), &statuses)
		if len(statuses) != 1 || statuses[0].ID != r.dev.String() {
			t.Fatalf("group details = %+v", statuses)
		}
		found := false
		for _, g := range statuses[0].Groups {
			found = found || g == "dash"
		}
		if !found {
			t.Errorf("device groups %v missing \"dash\"", statuses[0].Groups)
		}
	})
}
