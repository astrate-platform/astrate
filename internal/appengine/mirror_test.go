//go:build integration

package appengine

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

const mirrorMT = "application/merge-patch+json"

// newGroup creates a group with the given devices through the API surface.
func newGroup(t *testing.T, r *rig, name string, devs ...deviceid.ID) {
	t.Helper()
	parts := make([]string, len(devs))
	for i := range devs {
		parts[i] = jsonStr(devs[i].String())
	}
	body := `{"group_name":` + jsonStr(name) + `,"devices":[` + strings.Join(parts, ",") + `]}`
	if rec := r.req(t, http.MethodPost, "/groups", body, r.token); rec.Code != http.StatusCreated {
		t.Fatalf("create group %s: got %d, want 201 (%s)", name, rec.Code, rec.Body)
	}
}

func TestMirrorInterfacesList(t *testing.T) {
	r := newRig(t)
	want := []string{aeConf, aeSensors} // sorted ascending

	var viaDevice []string
	decodeData(t, r.req(t, http.MethodGet, r.dpath("/interfaces"), "", r.token), &viaDevice)
	if len(viaDevice) != len(want) || viaDevice[0] != want[0] || viaDevice[1] != want[1] {
		t.Errorf("device interfaces = %v, want %v", viaDevice, want)
	}

	var viaAlias []string
	decodeData(t, r.req(t, http.MethodGet, "/devices-by-alias/sensor-1/interfaces", "", r.token), &viaAlias)
	if len(viaAlias) != len(want) || viaAlias[0] != want[0] || viaAlias[1] != want[1] {
		t.Errorf("alias interfaces = %v, want %v", viaAlias, want)
	}

	newGroup(t, r, "mg", r.dev)
	var viaGroup []string
	decodeData(t, r.req(t, http.MethodGet,
		"/groups/mg/devices/"+r.dev.String()+"/interfaces", "", r.token), &viaGroup)
	if len(viaGroup) != len(want) || viaGroup[0] != want[0] || viaGroup[1] != want[1] {
		t.Errorf("group interfaces = %v, want %v", viaGroup, want)
	}

	for _, tc := range []struct{ label, url string }{
		{"unknown alias", "/devices-by-alias/nope/interfaces"},
		{"unknown device", "/devices/" + unknownID(t) + "/interfaces"},
	} {
		if rec := r.req(t, http.MethodGet, tc.url, "", r.token); rec.Code != http.StatusNotFound {
			t.Errorf("%s: got %d, want 404 (%s)", tc.label, rec.Code, rec.Body)
		}
	}
}

func TestMirrorAliasPatch(t *testing.T) {
	r := newRig(t)

	var ds DeviceStatus
	decodeData(t, r.reqCT(t, http.MethodPatch, "/devices-by-alias/sensor-1",
		`{"aliases":{"label":"renamed"}}`, r.token, mirrorMT), &ds)
	if ds.Aliases["label"] != "renamed" {
		t.Errorf("aliases after patch = %v, want label=renamed", ds.Aliases)
	}

	rec := r.reqCT(t, http.MethodPatch, "/devices-by-alias/sensor-1",
		`{"credentials_inhibited":true}`, r.token, "application/json")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("wrong content-type: got %d (%s), want exactly 400", rec.Code, rec.Body)
	}
	if !strings.Contains(rec.Body.String(), "Bad request") {
		t.Errorf(`wrong content-type body = %s, want it to contain "Bad request"`, rec.Body)
	}

	rec = r.reqCT(t, http.MethodPatch, "/devices-by-alias/nope",
		`{"credentials_inhibited":true}`, r.token, mirrorMT)
	if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "Device not found") {
		t.Errorf("unknown alias: got %d (%s), want 404 Device not found", rec.Code, rec.Body)
	}
}

func TestMirrorAliasDataAccess(t *testing.T) {
	r := newRig(t)

	var direct, viaAlias []Sample
	decodeData(t, r.req(t, http.MethodGet,
		r.dpath("/interfaces/"+aeSensors+"/value"), "", r.token), &direct)
	decodeData(t, r.req(t, http.MethodGet,
		"/devices-by-alias/sensor-1/interfaces/"+aeSensors+"/value", "", r.token), &viaAlias)
	if len(viaAlias) != len(direct) {
		t.Fatalf("series lengths differ: alias %d vs device %d", len(viaAlias), len(direct))
	}
	for i := range direct {
		if !viaAlias[i].Timestamp.Equal(direct[i].Timestamp) {
			t.Errorf("sample %d timestamp %v != device-scope %v", i, viaAlias[i].Timestamp, direct[i].Timestamp)
		}
	}

	if rec := r.req(t, http.MethodPut,
		"/devices-by-alias/sensor-1/interfaces/"+aeSensors+"/value", "7.5", r.token); rec.Code != http.StatusOK {
		t.Errorf("put via alias: got %d, want 200 (%s)", rec.Code, rec.Body)
	}
	r.sd.mu.Lock()
	if len(r.sd.pubs) != 1 {
		t.Errorf("captured pubs = %v, want exactly one", r.sd.pubs)
	} else if pub := r.sd.pubs[0]; !strings.Contains(pub, r.realm+"|"+r.dev.String()+"|"+aeSensors+"/value|7.5") {
		t.Errorf("captured pub = %q, want real device id and path", pub)
	}
	r.sd.mu.Unlock()

	if rec := r.req(t, http.MethodDelete,
		"/devices-by-alias/sensor-1/interfaces/"+aeSensors+"/value", "", r.token); rec.Code != http.StatusNoContent {
		t.Errorf("delete via alias: got %d, want 204", rec.Code)
	}
	r.sd.mu.Lock()
	if len(r.sd.unsets) != 1 || !strings.Contains(r.sd.unsets[0], r.dev.String()) {
		t.Errorf("captured unsets = %v, want one carrying the real device id", r.sd.unsets)
	}
	r.sd.mu.Unlock()

	if rec := r.req(t, http.MethodGet,
		"/devices-by-alias/nope/interfaces/"+aeSensors+"/value", "", r.token); rec.Code != http.StatusNotFound {
		t.Errorf("unknown alias: got %d, want 404", rec.Code)
	}
}

func TestMirrorGroupDevice(t *testing.T) {
	r := newRig(t)
	ctx := context.Background()
	const gm = "gm"
	newGroup(t, r, gm, r.dev)

	// A registered device that is NOT a member, with introspection installed so
	// a missing-membership gate would still let its data access resolve.
	stranger, _ := deviceid.Random()
	if err := r.st.RegisterDevice(ctx, r.realmID, stranger, "h"); err != nil {
		t.Fatal(err)
	}
	if _, err := r.st.UpdateIntrospection(ctx, r.realmID, stranger, map[string]store.InterfaceVersion{
		aeSensors: {Major: 1, Minor: 0},
	}); err != nil {
		t.Fatal(err)
	}

	var ds DeviceStatus
	decodeData(t, r.req(t, http.MethodGet,
		"/groups/"+gm+"/devices/"+r.dev.String(), "", r.token), &ds)
	if ds.ID != r.dev.String() {
		t.Errorf("group device status id = %q, want %q", ds.ID, r.dev.String())
	}

	decodeData(t, r.reqCT(t, http.MethodPatch, "/groups/"+gm+"/devices/"+r.dev.String(),
		`{"credentials_inhibited":true}`, r.token, mirrorMT), &ds)
	if !ds.CredentialsInhibited {
		t.Error("credentials_inhibited not applied through the group patch")
	}

	rec := r.reqCT(t, http.MethodPatch, "/groups/"+gm+"/devices/"+r.dev.String(),
		`{"credentials_inhibited":false}`, r.token, "application/json")
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("wrong content-type: got %d, want 500 (upstream unmapped fallback)", rec.Code)
	}

	nonmember := []struct {
		name   string
		method string
		path   string
	}{
		{"status", http.MethodGet, "/groups/" + gm + "/devices/" + stranger.String()},
		{"data", http.MethodGet, "/groups/" + gm + "/devices/" + stranger.String() + "/interfaces/" + aeSensors + "/value"},
	}
	for _, x := range nonmember {
		t.Run("nonmember-"+x.name, func(t *testing.T) {
			got := r.req(t, x.method, x.path, "", r.token)
			if got.Code != http.StatusNotFound || !strings.Contains(got.Body.String(), "Device not found") {
				t.Errorf("%s nonmember: got %d (%s), want 404 Device not found", x.name, got.Code, got.Body)
			}
		})
	}
	rec = r.reqCT(t, http.MethodPatch, "/groups/"+gm+"/devices/"+stranger.String(),
		`{"credentials_inhibited":true}`, r.token, mirrorMT)
	if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "Device not found") {
		t.Errorf("patch nonmember: got %d (%s), want 404 Device not found", rec.Code, rec.Body)
	}

	// Data access under the group behaves like the device scope.
	var inGroup []Sample
	decodeData(t, r.req(t, http.MethodGet,
		"/groups/"+gm+"/devices/"+r.dev.String()+"/interfaces/"+aeSensors+"/value", "", r.token), &inGroup)
	if len(inGroup) != 3 {
		t.Errorf("group series = %d samples, want 3", len(inGroup))
	}
	if rec := r.req(t, http.MethodPut,
		"/groups/"+gm+"/devices/"+r.dev.String()+"/interfaces/"+aeSensors+"/value", "9.5", r.token); rec.Code != http.StatusOK {
		t.Errorf("put in group: got %d, want 200 (%s)", rec.Code, rec.Body)
	}
	r.sd.mu.Lock()
	if len(r.sd.pubs) != 1 || !strings.Contains(r.sd.pubs[0], r.realm+"|"+r.dev.String()+"|"+aeSensors+"/value|9.5") {
		t.Errorf("captured pubs = %v, want one carrying the real device id", r.sd.pubs)
	}
	r.sd.mu.Unlock()
	if rec := r.req(t, http.MethodDelete,
		"/groups/"+gm+"/devices/"+r.dev.String()+"/interfaces/"+aeSensors+"/value", "", r.token); rec.Code != http.StatusNoContent {
		t.Errorf("delete in group: got %d, want 204", rec.Code)
	}
	r.sd.mu.Lock()
	if len(r.sd.unsets) != 1 || !strings.Contains(r.sd.unsets[0], r.dev.String()) {
		t.Errorf("captured unsets = %v, want one carrying the real device id", r.sd.unsets)
	}
	r.sd.mu.Unlock()

	// Unknown group name maps to 404 "Group not found" on every mirrored surface.
	unknownGroup := []struct {
		name   string
		method string
		path   string
	}{
		{"status", http.MethodGet, "/groups/nope/devices/" + r.dev.String()},
		{"interfaces", http.MethodGet, "/groups/nope/devices/" + r.dev.String() + "/interfaces"},
		{"data", http.MethodGet, "/groups/nope/devices/" + r.dev.String() + "/interfaces/" + aeSensors + "/value"},
	}
	for _, x := range unknownGroup {
		t.Run("unknown-group-"+x.name, func(t *testing.T) {
			got := r.req(t, x.method, x.path, "", r.token)
			if got.Code != http.StatusNotFound || !strings.Contains(got.Body.String(), "Group not found") {
				t.Errorf("%s unknown group: got %d (%s), want 404 Group not found", x.name, got.Code, got.Body)
			}
		})
	}
	rec = r.reqCT(t, http.MethodPatch, "/groups/nope/devices/"+r.dev.String(),
		`{"credentials_inhibited":true}`, r.token, mirrorMT)
	if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "Group not found") {
		t.Errorf("patch unknown group: got %d (%s), want 404 Group not found", rec.Code, rec.Body)
	}
}
