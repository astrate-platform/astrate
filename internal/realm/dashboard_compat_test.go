//go:build integration

package realm

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// fakeDisconnecter records device-kick calls.
type fakeDisconnecter struct {
	kicked []string
}

func (f *fakeDisconnecter) DisconnectDevice(realm string, id deviceid.ID) {
	f.kicked = append(f.kicked, realm+"/"+id.String())
}

// TestDashboardCompat covers the M10 realm-management surface the Astarte
// Dashboard v1.2.2 requires: version, device_registration_limit, delivery
// policies CRUD, and synchronous device deletion.
func TestDashboardCompat(t *testing.T) {
	r := newRig(t)
	ctx := context.Background()

	t.Run("Version", func(t *testing.T) {
		var v string
		decodeData(t, r.req(t, http.MethodGet, "/version", "", r.rmaToken), &v)
		if v != APICompatVersion {
			t.Errorf("version = %q, want %q", v, APICompatVersion)
		}
	})

	t.Run("DeviceRegistrationLimit", func(t *testing.T) {
		rec := r.req(t, http.MethodGet, "/config/device_registration_limit", "", r.rmaToken)
		if rec.Code != http.StatusOK {
			t.Fatalf("limit: got %d (%s)", rec.Code, rec.Body)
		}
		if body := rec.Body.String(); body != `{"data":null}` {
			t.Errorf("unlimited realm body = %s, want {\"data\":null}", body)
		}
	})

	t.Run("Policies", func(t *testing.T) {
		def := `{"name":"retry5xx","error_handlers":[{"on":"server_error","strategy":"retry"}],"maximum_capacity":100,"retry_times":3}`
		if rec := r.req(t, http.MethodPost, "/policies", def, r.rmaToken); rec.Code != http.StatusCreated {
			t.Fatalf("create policy: %d (%s)", rec.Code, rec.Body)
		}
		if rec := r.req(t, http.MethodPost, "/policies", def, r.rmaToken); rec.Code != http.StatusConflict {
			t.Errorf("duplicate policy: %d, want 409", rec.Code)
		}
		if rec := r.req(t, http.MethodPost, "/policies",
			`{"name":"bad","error_handlers":[],"maximum_capacity":1}`, r.rmaToken); rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("invalid policy: %d, want 422", rec.Code)
		}
		var names []string
		decodeData(t, r.req(t, http.MethodGet, "/policies", "", r.rmaToken), &names)
		if len(names) != 1 || names[0] != "retry5xx" {
			t.Errorf("policy names = %v", names)
		}
		var got struct {
			Name          string `json:"name"`
			RetryTimes    int    `json:"retry_times"`
			ErrorHandlers []any  `json:"error_handlers"`
		}
		decodeData(t, r.req(t, http.MethodGet, "/policies/retry5xx", "", r.rmaToken), &got)
		if got.Name != "retry5xx" || got.RetryTimes != 3 || len(got.ErrorHandlers) != 1 {
			t.Errorf("policy = %+v", got)
		}
		if rec := r.req(t, http.MethodDelete, "/policies/retry5xx", "", r.rmaToken); rec.Code != http.StatusNoContent {
			t.Errorf("delete policy: %d, want 204", rec.Code)
		}
		if rec := r.req(t, http.MethodGet, "/policies/retry5xx", "", r.rmaToken); rec.Code != http.StatusNotFound {
			t.Errorf("deleted policy get: %d, want 404", rec.Code)
		}
	})

	t.Run("DeviceDeletion", func(t *testing.T) {
		// A device with data, wired through a service that has the kick seam.
		disc := &fakeDisconnecter{}
		r.svc.WithDisconnecter(disc)

		dev, err := deviceid.Random()
		if err != nil {
			t.Fatal(err)
		}
		if err := r.st.RegisterDevice(ctx, r.realmID, dev, "h"); err != nil {
			t.Fatal(err)
		}
		si, err := r.st.InstallInterface(ctx, r.realmID, []byte(`{
			"interface_name": "org.astrate.rm.DelData", "version_major": 1, "version_minor": 0,
			"type": "datastream", "ownership": "device",
			"mappings": [{"endpoint": "/v", "type": "double"}]}`))
		if err != nil {
			t.Fatal(err)
		}
		v := 4.2
		ts := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
		if err := r.st.AppendDatastreams(ctx, store.DatastreamBatch{Individual: []store.IndividualRow{{
			RealmID: r.realmID, DeviceID: dev, InterfaceID: si.ID, EndpointID: si.Endpoints["/v"],
			Path: "/v", TS: ts, ReceptionTS: ts, ValueDouble: &v,
		}}}); err != nil {
			t.Fatal(err)
		}

		if rec := r.req(t, http.MethodDelete, "/devices/"+dev.String(), "", r.rmaToken); rec.Code != http.StatusNoContent {
			t.Fatalf("delete device: %d (%s)", rec.Code, rec.Body)
		}
		if len(disc.kicked) != 1 || disc.kicked[0] != r.realm+"/"+dev.String() {
			t.Errorf("kick calls = %v", disc.kicked)
		}
		if _, err := r.st.GetDevice(ctx, r.realmID, dev); err == nil {
			t.Error("device row survived the delete")
		}
		rows, err := r.st.Series(ctx, store.SeriesQuery{RealmID: r.realmID, DeviceID: dev, InterfaceID: si.ID, Path: "/v"})
		if err != nil || len(rows) != 0 {
			t.Errorf("data survived the delete: %d rows, err %v", len(rows), err)
		}
		// Unknown and malformed IDs → the upstream device 404 envelope.
		for _, path := range []string{"/devices/" + dev.String(), "/devices/not-an-id"} {
			rec := r.req(t, http.MethodDelete, path, "", r.rmaToken)
			if rec.Code != http.StatusNotFound {
				t.Errorf("DELETE %s: %d, want 404", path, rec.Code)
			}
			if body := rec.Body.String(); body != `{"errors":{"detail":"Device not found"}}` {
				t.Errorf("DELETE %s body = %s", path, body)
			}
		}
	})
}

// TestPolicyReferentialIntegrity covers the two checks that tie triggers to
// delivery policies. Neither is reachable from T1: internal/realm's whole
// behavioural surface needs a live database, so without this the checks would
// ship proven only to compile.
func TestPolicyReferentialIntegrity(t *testing.T) {
	r := newRig(t)

	const (
		policy = `{"name":"refpolicy","error_handlers":[{"on":"server_error","strategy":"retry"}],"maximum_capacity":10,"retry_times":2}`
		hook   = `"action":{"http_url":"https://example.com/hook","http_method":"post"},` +
			`"simple_triggers":[{"type":"device_trigger","on":"device_connected"}]`
	)

	// A trigger naming a policy that was never installed is refused, and
	// nothing is stored for it.
	unknown := `{"name":"wants_missing","policy":"nosuchpolicy",` + hook + `}`
	if rec := r.req(t, http.MethodPost, "/triggers", unknown, r.rmaToken); rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("trigger naming an unknown policy: %d, want 422 (%s)", rec.Code, rec.Body)
	}
	if rec := r.req(t, http.MethodGet, "/triggers/wants_missing", "", r.rmaToken); rec.Code != http.StatusNotFound {
		t.Errorf("refused trigger was still stored: %d, want 404", rec.Code)
	}

	// With the policy installed, the same trigger goes in.
	if rec := r.req(t, http.MethodPost, "/policies", policy, r.rmaToken); rec.Code != http.StatusCreated {
		t.Fatalf("create policy: %d (%s)", rec.Code, rec.Body)
	}
	named := `{"name":"wants_refpolicy","policy":"refpolicy",` + hook + `}`
	if rec := r.req(t, http.MethodPost, "/triggers", named, r.rmaToken); rec.Code != http.StatusCreated {
		t.Fatalf("trigger naming an installed policy: %d, want 201 (%s)", rec.Code, rec.Body)
	}

	// The policy cannot be deleted out from under it.
	rec := r.req(t, http.MethodDelete, "/policies/refpolicy", "", r.rmaToken)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("delete referenced policy: %d, want 422 (%s)", rec.Code, rec.Body)
	}
	if body := rec.Body.String(); !strings.Contains(body, "wants_refpolicy") {
		t.Errorf("error should name the referencing trigger, got %s", body)
	}
	if rec := r.req(t, http.MethodGet, "/policies/refpolicy", "", r.rmaToken); rec.Code != http.StatusOK {
		t.Errorf("refused delete removed the policy anyway: %d", rec.Code)
	}

	// Once the last reference is gone, the delete succeeds.
	if rec := r.req(t, http.MethodDelete, "/triggers/wants_refpolicy", "", r.rmaToken); rec.Code != http.StatusNoContent {
		t.Fatalf("delete trigger: %d (%s)", rec.Code, rec.Body)
	}
	if rec := r.req(t, http.MethodDelete, "/policies/refpolicy", "", r.rmaToken); rec.Code != http.StatusNoContent {
		t.Errorf("delete unreferenced policy: %d, want 204 (%s)", rec.Code, rec.Body)
	}
}
