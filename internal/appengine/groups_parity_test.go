//go:build integration

package appengine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// linkedPage is a paged-collection envelope with its links object.
type linkedPage struct {
	Data  json.RawMessage `json:"data"`
	Links struct {
		Self string `json:"self"`
		Next string `json:"next"`
	} `json:"links"`
}

// decodeLinkedData decodes a 200 envelope into dst and returns its links.
func decodeLinkedData(t *testing.T, rec *httptest.ResponseRecorder) (linkedPage, func(dst any)) {
	t.Helper()
	if rec.Code != http.StatusOK {
		t.Fatalf("non-200 response %d: %s", rec.Code, rec.Body)
	}
	var env linkedPage
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode envelope: %v (%s)", err, rec.Body)
	}
	return env, func(dst any) {
		if err := json.Unmarshal(env.Data, dst); err != nil {
			t.Fatalf("decode data: %v (%s)", err, env.Data)
		}
	}
}

// followGET requests a returned relative link verbatim through the rig.
func followGET(t *testing.T, r *rig, next string) *httptest.ResponseRecorder {
	t.Helper()
	return r.req(t, http.MethodGet, strings.TrimPrefix(next, "/v1/"+r.realm), "", r.token)
}

func TestShowGroup(t *testing.T) {
	r := newRig(t)
	newGroup(t, r, "sg", r.dev)

	var got map[string]string
	decodeData(t, r.req(t, http.MethodGet, "/groups/sg", "", r.token), &got)
	if got["group_name"] != "sg" {
		t.Errorf("data.group_name = %q, want %q", got["group_name"], "sg")
	}

	rec := r.req(t, http.MethodGet, "/groups/nope", "", r.token)
	if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "Group not found") {
		t.Errorf("unknown group: got %d (%s), want 404 Group not found", rec.Code, rec.Body)
	}
}

func TestCreateGroupValidation(t *testing.T) {
	r := newRig(t)
	newGroup(t, r, "taken", r.dev)

	ghost := unknownID(t)
	for _, tc := range []struct {
		label    string
		body     string
		status   int
		fragment string
	}{
		{
			label:    "missing group_name",
			body:     `{"devices":[` + jsonStr(r.dev.String()) + `]}`,
			status:   http.StatusUnprocessableEntity,
			fragment: `"group_name":["can't be blank"]`,
		},
		{
			label:    "blank group_name",
			body:     `{"group_name":"","devices":[` + jsonStr(r.dev.String()) + `]}`,
			status:   http.StatusUnprocessableEntity,
			fragment: `"group_name":["can't be blank"]`,
		},
		{
			label:    "missing devices key",
			body:     `{"group_name":"ghost-missing-devices"}`,
			status:   http.StatusUnprocessableEntity,
			fragment: `"devices":["can't be blank"]`,
		},
		{
			label:    "empty devices array",
			body:     `{"group_name":"ghost-empty-devices","devices":[]}`,
			status:   http.StatusUnprocessableEntity,
			fragment: `"devices":["should have at least 1 item(s)"]`,
		},
		{
			label:    "malformed device id",
			body:     `{"group_name":"ghost-malformed","devices":["not-a-device-id"]}`,
			status:   http.StatusUnprocessableEntity,
			fragment: `"must exist (not-a-device-id not found)"`,
		},
		{
			label:    "unknown well-formed device id",
			body:     `{"group_name":"ghost-unknown","devices":[` + jsonStr(ghost) + `]}`,
			status:   http.StatusUnprocessableEntity,
			fragment: `"must exist (` + ghost + ` not found)"`,
		},
		{
			label:    "duplicate group name",
			body:     `{"group_name":"taken","devices":[` + jsonStr(r.dev.String()) + `]}`,
			status:   http.StatusConflict,
			fragment: `"Group already exists"`,
		},
	} {
		if rec := r.req(t, http.MethodPost, "/groups", tc.body, r.token); rec.Code != tc.status ||
			!strings.Contains(rec.Body.String(), tc.fragment) {
			t.Errorf("%s: got %d (%s), want %d containing %s",
				tc.label, rec.Code, rec.Body, tc.status, tc.fragment)
		}
	}

	// A failed create must leave no group behind.
	var names []string
	decodeData(t, r.req(t, http.MethodGet, "/groups", "", r.token), &names)
	for _, ghost := range []string{"ghost-missing-devices", "ghost-empty-devices", "ghost-malformed", "ghost-unknown"} {
		if contains(names, ghost) {
			t.Errorf("failed create left group %q behind: %v", ghost, names)
		}
	}

	// A valid body with one id duplicated is accepted; repeating it hits 409.
	dupBody := `{"group_name":"dup","devices":[` + jsonStr(r.dev.String()) + "," + jsonStr(r.dev.String()) + `]}`
	if rec := r.req(t, http.MethodPost, "/groups", dupBody, r.token); rec.Code != http.StatusCreated {
		t.Errorf("create with duplicated id: got %d, want 201 (%s)", rec.Code, rec.Body)
	}
	if rec := r.req(t, http.MethodPost, "/groups", dupBody, r.token); rec.Code != http.StatusConflict ||
		!strings.Contains(rec.Body.String(), `"Group already exists"`) {
		t.Errorf("repeat create: got %d (%s), want 409 Group already exists", rec.Code, rec.Body)
	}

	names = nil
	decodeData(t, r.req(t, http.MethodGet, "/groups", "", r.token), &names)
	if got := strings.Count(strings.Join(names, "\x00"), "dup"); got != 1 {
		t.Errorf("dup listed %d times, want exactly once: %v", got, names)
	}
	var members []string
	decodeData(t, r.req(t, http.MethodGet, "/groups/dup/devices", "", r.token), &members)
	if len(members) != 1 || members[0] != r.dev.String() {
		t.Errorf("duplicated id inserted %d times, want once: %v", len(members), members)
	}
}

func TestAddGroupDeviceErrors(t *testing.T) {
	r := newRig(t)
	newGroup(t, r, "ag", r.dev)

	rec := r.req(t, http.MethodPost, "/groups/ag/devices",
		`{"device_id":"not-a-device-id"}`, r.token)
	if rec.Code != http.StatusUnprocessableEntity || !strings.Contains(rec.Body.String(), "is not a valid device id") {
		t.Errorf("malformed device id: got %d (%s), want 422 is not a valid device id", rec.Code, rec.Body)
	}

	rec = r.req(t, http.MethodPost, "/groups/ag/devices",
		`{"device_id":`+jsonStr(r.dev.String())+`}`, r.token)
	if rec.Code != http.StatusConflict || !strings.Contains(rec.Body.String(), "Device already in group") {
		t.Errorf("re-add member: got %d (%s), want 409 Device already in group", rec.Code, rec.Body)
	}
}

func TestGroupListingPagination(t *testing.T) {
	r := newRig(t)
	ctx := context.Background()
	const pg = "pg"
	ids := make([]deviceid.ID, 0, 5)
	for range 5 {
		id, err := deviceid.Random()
		if err != nil {
			t.Fatal(err)
		}
		if err := r.st.RegisterDevice(ctx, r.realmID, id, "h"); err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}
	newGroup(t, r, pg, ids...)

	// Without params: the whole set, bare self link, no next.
	env, decode := decodeLinkedData(t, r.req(t, http.MethodGet, "/groups/"+pg+"/devices", "", r.token))
	var all []string
	decode(&all)
	if env.Links.Self != "/v1/"+r.realm+"/groups/"+pg+"/devices" {
		t.Errorf("links.self = %q, want no query string", env.Links.Self)
	}
	if env.Links.Next != "" {
		t.Errorf("links.next = %q on an unpaged listing, want empty", env.Links.Next)
	}
	if !sameStringSet(all, idsToStrings(ids)) {
		t.Errorf("unpaged listing = %v, want exactly the seeded set %v", all, ids)
	}

	// limit=2 walks three pages whose union is the seeded set.
	seen := map[string]bool{}
	next := "/groups/" + pg + "/devices?limit=2"
	pages := 0
	for next != "" {
		pages++
		if pages > 3 {
			t.Fatal("pagination did not terminate after three pages")
		}
		env, decode = decodeLinkedData(t, followGET(t, r, next))
		var pageIDs []string
		decode(&pageIDs)
		if pages == 1 && env.Links.Self != "/v1/"+r.realm+"/groups/"+pg+"/devices?limit=2" {
			t.Errorf("links.self = %q, want the request query echoed", env.Links.Self)
		}
		wantLen := map[int]int{1: 2, 2: 2, 3: 1}[pages]
		if len(pageIDs) != wantLen {
			t.Errorf("page %d returned %v, want %d ids", pages, pageIDs, wantLen)
		}
		for _, id := range pageIDs {
			seen[id] = true
		}
		nextURL := env.Links.Next
		if pages < 3 && nextURL == "" {
			t.Fatalf("page %d has no links.next: %s", pages, env.Links.Next)
		}
		if nextURL != "" {
			u, err := url.Parse(nextURL)
			if err != nil {
				t.Fatalf("parse links.next %q: %v", nextURL, err)
			}
			if u.Path != "/v1/"+r.realm+"/groups/"+pg+"/devices" {
				t.Errorf("links.next path = %q", u.Path)
			}
			if got := u.Query().Get("from_token"); got == "" {
				t.Errorf("links.next %q carries no from_token", nextURL)
			} else if _, ok := parseGroupToken(got); !ok {
				t.Errorf("links.next from_token %q is not a UUID-v1-format token", got)
			}
			if got := u.Query().Get("limit"); got != "2" {
				t.Errorf("links.next limit = %q, want 2", got)
			}
			if strings.Contains(u.RawQuery, "details") {
				t.Errorf("links.next %q carries details although the request had none", nextURL)
			}
		}
		next = nextURL
	}
	if !sameStringSet(keys(seen), idsToStrings(ids)) {
		t.Errorf("pages covered %v, want exactly the seeded set %v", keys(seen), ids)
	}

	// from_token must be a v1-format uuid: junk and v4 are both rejected.
	rec := r.req(t, http.MethodGet, "/groups/"+pg+"/devices?from_token=zzz", "", r.token)
	if rec.Code != http.StatusUnprocessableEntity || !strings.Contains(rec.Body.String(), "is invalid") {
		t.Errorf("from_token=zzz: got %d (%s), want 422 is invalid", rec.Code, rec.Body)
	}
	v4, err := deviceid.Random()
	if err != nil {
		t.Fatal(err)
	}
	rec = r.req(t, http.MethodGet, "/groups/"+pg+"/devices?from_token="+v4.UUID(), "", r.token)
	if rec.Code != http.StatusUnprocessableEntity || !strings.Contains(rec.Body.String(), "is invalid") {
		t.Errorf("from_token=v4 uuid: got %d (%s), want 422 is invalid", rec.Code, rec.Body)
	}
	rec = r.req(t, http.MethodGet, "/groups/"+pg+"/devices?from_token="+groupTokenFor(0), "", r.token)
	if rec.Code != http.StatusOK {
		t.Errorf("from_token=groupTokenFor(0): got %d (%s), want 200", rec.Code, rec.Body)
	}

	// A negative limit is refused with the upstream message.
	rec = r.req(t, http.MethodGet, "/groups/"+pg+"/devices?limit=-1", "", r.token)
	if rec.Code != http.StatusUnprocessableEntity ||
		!strings.Contains(rec.Body.String(), "must be greater than or equal to 0") {
		t.Errorf("limit=-1: got %d (%s), want 422 must be greater than or equal to 0", rec.Code, rec.Body)
	}
}

func idsToStrings(ids []deviceid.ID) []string {
	out := make([]string, len(ids))
	for i := range ids {
		out[i] = ids[i].String()
	}
	return out
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func sameStringSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	x := slices.Clone(a)
	y := slices.Clone(b)
	slices.Sort(x)
	slices.Sort(y)
	return slices.Equal(x, y)
}
