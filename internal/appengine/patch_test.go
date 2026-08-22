//go:build integration

package appengine

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// reqCT issues an authenticated request with an explicit Content-Type, which
// req() cannot express: the device PATCH behavior depends on it byte-for-byte.
func (r *rig) reqCT(t *testing.T, method, path, rawBody, token, contentType string) *httptest.ResponseRecorder {
	t.Helper()
	var body io.Reader
	if rawBody != "" {
		body = strings.NewReader(`{"data":` + rawBody + `}`)
	}
	hr := httptest.NewRequest(method, "/appengine/v1/"+r.realm+path, body)
	hr.Header.Set("Content-Type", contentType)
	if token != "" {
		hr.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	r.mux.ServeHTTP(rec, hr)
	return rec
}

func TestPatchDeviceContentType(t *testing.T) {
	r := newRig(t)
	body := `{"credentials_inhibited":true}`

	// Upstream compares the Content-Type header to exactly one value; both a
	// plain JSON PATCH and a charset-parametered merge-patch miss it.
	for _, ct := range []string{"application/json", "application/merge-patch+json; charset=utf-8", ""} {
		rec := r.reqCT(t, http.MethodPatch, r.dpath(""), body, r.token, ct)
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Content-Type %q: got %d, want 500 (upstream unmapped fallback)", ct, rec.Code)
		}
	}

	var ds DeviceStatus
	decodeData(t, r.reqCT(t, http.MethodPatch, r.dpath(""), body, r.token, "application/merge-patch+json"), &ds)
	if !ds.CredentialsInhibited {
		t.Errorf("patch not applied: %+v", ds)
	}
}

func TestPatchDeviceAliasAttributesTaxonomy(t *testing.T) {
	r := newRig(t)
	const mt = "application/merge-patch+json"

	// A second device holding an alias, so ownership conflicts are reachable.
	ctx := context.Background()
	other, _ := deviceid.Random()
	if err := r.st.RegisterDevice(ctx, r.realmID, other, "h"); err != nil {
		t.Fatal(err)
	}
	if err := r.st.PatchDeviceAliases(ctx, r.realmID, other,
		map[string]*string{"site": ptr("plant-2")}); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name   string
		body   string
		want   int
		detail string
	}{
		{"empty alias tag", `{"aliases":{"":"x"}}`, http.StatusBadRequest, "Invalid alias"},
		{"empty alias value", `{"aliases":{"k":""}}`, http.StatusBadRequest, "Invalid alias"},
		{"alias owned elsewhere", `{"aliases":{"zone":"plant-2"}}`, http.StatusConflict, "Alias already in use"},
		{"delete unknown tag", `{"aliases":{"nope":null}}`, http.StatusBadRequest, "Alias tag not found"},
		{"empty attribute key", `{"attributes":{"":"v"}}`, http.StatusBadRequest, "Invalid attributes"},
		{"delete unknown attribute", `{"attributes":{"nope":null}}`, http.StatusUnprocessableEntity, "Attribute key not found"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := r.reqCT(t, http.MethodPatch, r.dpath(""), tc.body, r.token, mt)
			if rec.Code != tc.want {
				t.Fatalf("got %d, want %d (%s)", rec.Code, tc.want, rec.Body.String())
			}
			var got struct {
				Errors struct {
					Detail string `json:"detail"`
				} `json:"errors"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode error body: %v (%s)", err, rec.Body)
			}
			if got.Errors.Detail != tc.detail {
				t.Errorf("detail = %q, want %q", got.Errors.Detail, tc.detail)
			}
		})
	}

	// The happy paths: setting a fresh alias and deleting an existing tag must
	// both go through.
	var ds DeviceStatus
	decodeData(t, r.reqCT(t, http.MethodPatch, r.dpath(""), `{"aliases":{"label":null,"room":"lab"}}`, r.token, mt), &ds)
	if _, ok := ds.Aliases["label"]; ok {
		t.Errorf("label not deleted: %+v", ds.Aliases)
	}
	if ds.Aliases["room"] != "lab" {
		t.Errorf("room = %q, want lab", ds.Aliases["room"])
	}
}
