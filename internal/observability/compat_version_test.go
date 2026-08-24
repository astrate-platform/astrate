package observability

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMountVersionCompat(t *testing.T) {
	mux := http.NewServeMux()
	for _, svc := range []string{"housekeeping", "appengine", "realmmanagement", "pairing"} {
		MountVersionCompat(mux, svc, "0.0.0-test")
	}
	srv := httptest.NewServer(mux)
	defer srv.Close()

	for _, svc := range []string{"housekeeping", "appengine", "realmmanagement", "pairing"} {
		resp, err := http.Get(srv.URL + "/" + svc + "/version")
		if err != nil {
			t.Fatalf("%s/version: %v", svc, err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("%s/version = %d, want 200", svc, resp.StatusCode)
		}
		if string(body) != `{"data":"0.0.0-test"}` {
			t.Errorf("%s/version body = %s", svc, body)
		}
		if ct := resp.Header.Get("Content-Type"); ct == "" {
			t.Errorf("%s/version missing Content-Type", svc)
		}
	}
}

func TestVersionHandlerRealmScoped(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/{realm}/version", VersionHandler("x"))
	srv := httptest.NewServer(mux)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/v1/somerealm/version")
	if err != nil {
		t.Fatalf("/v1/somerealm/version: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("/v1/somerealm/version = %d, want 200", resp.StatusCode)
	}
	if string(body) != `{"data":"x"}` {
		t.Errorf("/v1/somerealm/version body = %s", body)
	}
}
