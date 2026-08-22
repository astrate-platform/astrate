//go:build integration

package realm

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// Interface fixtures (one name, evolving versions, plus a deletable draft).
const (
	rmIface  = "com.ex.M7a.Sensors"
	rmDraft  = "com.ex.M7a.Draft"
	ifaceV1  = `{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`
	ifaceV1b = `{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":1,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"},{"endpoint":"/count","type":"integer"}]}`
	ifaceV1x = `{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":2,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"integer"}]}`
	ifaceV2  = `{"interface_name":"com.ex.M7a.Sensors","version_major":2,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"string"}]}`
	draftV0  = `{"interface_name":"com.ex.M7a.Draft","version_major":0,"version_minor":1,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/x","type":"double"}]}`
)

const triggerJSON = `{"name":"on_value","action":{"http_url":"https://example.com/hook","http_method":"post"},` +
	`"simple_triggers":[{"type":"data_trigger","on":"incoming_data","interface_name":"com.ex.M7a.Sensors","interface_major":1,"match_path":"/value","value_match_operator":"*"}]}`

type rig struct {
	st       *store.Store
	svc      *Service
	mux      *http.ServeMux
	realm    string
	realmID  int16
	rmaToken string // valid a_rma
	wrongTok string // valid token, no a_rma
	jwtKey   *rsa.PrivateKey
	otherPub string
}

func newRig(t *testing.T) *rig {
	t.Helper()
	ctx := context.Background()
	pool := testutil.StartTimescale(t)
	st, err := store.New(ctx, pool.Config().ConnString())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(st.Close)

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	other, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	var suffix [4]byte
	_, _ = rand.Read(suffix[:])
	realmName := "rm" + hex.EncodeToString(suffix[:])
	realm, err := st.CreateRealm(ctx, store.NewRealm{
		Name:               realmName,
		JWTPublicKeysPEM:   []string{pubPEM(t, &key.PublicKey)},
		CACertificatePEM:   "test-ca",
		CAPrivateKeySealed: []byte("sealed"),
	})
	if err != nil {
		t.Fatalf("CreateRealm: %v", err)
	}

	svc := NewService(st, nil, discardLogger())
	mux := http.NewServeMux()
	NewAPI(svc, auth.NewMiddleware(st)).Mount(mux)

	return &rig{
		st: st, svc: svc, mux: mux, realm: realmName, realmID: realm.ID, jwtKey: key,
		otherPub: pubPEM(t, &other.PublicKey),
		rmaToken: mintToken(t, key, jwt.MapClaims{"a_rma": []string{".*::.*"}}),
		wrongTok: mintToken(t, key, jwt.MapClaims{"a_aea": []string{".*::.*"}}),
	}
}

// req drives one authenticated request; rawBody (if non-empty) is wrapped in
// the {"data": ...} envelope.
func (r *rig) req(t *testing.T, method, path, rawBody, token string) *httptest.ResponseRecorder {
	t.Helper()
	var body io.Reader
	if rawBody != "" {
		body = strings.NewReader(`{"data":` + rawBody + `}`)
	}
	httpReq := httptest.NewRequest(method, "/realmmanagement/v1/"+r.realm+path, body)
	if token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	r.mux.ServeHTTP(rec, httpReq)
	return rec
}

func TestRealmManagement(t *testing.T) {
	r := newRig(t)

	t.Run("Auth401WithoutToken", func(t *testing.T) {
		if rec := r.req(t, http.MethodGet, "/interfaces", "", ""); rec.Code != http.StatusUnauthorized {
			t.Errorf("no token: got %d, want 401", rec.Code)
		}
	})
	t.Run("Auth403WrongClaim", func(t *testing.T) {
		if rec := r.req(t, http.MethodGet, "/interfaces", "", r.wrongTok); rec.Code != http.StatusForbidden {
			t.Errorf("wrong claim: got %d, want 403", rec.Code)
		}
	})

	t.Run("InstallInterface", func(t *testing.T) {
		if rec := r.req(t, http.MethodPost, "/interfaces", ifaceV1, r.rmaToken); rec.Code != http.StatusCreated {
			t.Fatalf("install: got %d, want 201 (%s)", rec.Code, rec.Body)
		}
		// Duplicate → 409.
		if rec := r.req(t, http.MethodPost, "/interfaces", ifaceV1, r.rmaToken); rec.Code != http.StatusConflict {
			t.Errorf("duplicate install: got %d, want 409", rec.Code)
		}
	})

	t.Run("ListAndGet", func(t *testing.T) {
		var got []string
		decodeData(t, r.req(t, http.MethodGet, "/interfaces", "", r.rmaToken), &got)
		if !contains(got, rmIface) {
			t.Errorf("list interfaces = %v, want to contain %s", got, rmIface)
		}
		var majors []int
		decodeData(t, r.req(t, http.MethodGet, "/interfaces/"+rmIface, "", r.rmaToken), &majors)
		if len(majors) != 1 || majors[0] != 1 {
			t.Errorf("majors = %v, want [1]", majors)
		}
		if rec := r.req(t, http.MethodGet, "/interfaces/"+rmIface+"/1", "", r.rmaToken); rec.Code != http.StatusOK {
			t.Errorf("get interface: got %d, want 200", rec.Code)
		}
	})

	t.Run("MinorUpgradeAccepted", func(t *testing.T) {
		if rec := r.req(t, http.MethodPut, "/interfaces/"+rmIface+"/1", ifaceV1b, r.rmaToken); rec.Code != http.StatusNoContent {
			t.Fatalf("additive minor upgrade: got %d, want 204 (%s)", rec.Code, rec.Body)
		}
	})
	t.Run("MappingMutationRejected", func(t *testing.T) {
		// Changing /value's type is not an additive upgrade; upstream 1.2
		// answers with its measured 409 detail (#62).
		want := `{"errors":{"detail":"Interface update contains incompatible endpoint changes"}}`
		rec := r.req(t, http.MethodPut, "/interfaces/"+rmIface+"/1", ifaceV1x, r.rmaToken)
		if rec.Code != http.StatusConflict || rec.Body.String() != want {
			t.Errorf("mapping mutation: got %d %s, want 409 %s", rec.Code, rec.Body, want)
		}
	})

	t.Run("MajorCoexistence", func(t *testing.T) {
		if rec := r.req(t, http.MethodPost, "/interfaces", ifaceV2, r.rmaToken); rec.Code != http.StatusCreated {
			t.Fatalf("install major 2: got %d, want 201 (%s)", rec.Code, rec.Body)
		}
		var majors []int
		decodeData(t, r.req(t, http.MethodGet, "/interfaces/"+rmIface, "", r.rmaToken), &majors)
		if len(majors) != 2 || majors[0] != 1 || majors[1] != 2 {
			t.Errorf("majors after coexistence = %v, want [1 2]", majors)
		}
	})

	t.Run("DeleteRules", func(t *testing.T) {
		ctx := context.Background()
		// Major != 0 can't be deleted — upstream 403s this (#62).
		want := `{"errors":{"detail":"Interface can't be deleted"}}`
		rec := r.req(t, http.MethodDelete, "/interfaces/"+rmIface+"/1", "", r.rmaToken)
		if rec.Code != http.StatusForbidden || rec.Body.String() != want {
			t.Errorf("delete major 1: got %d %s, want 403 %s", rec.Code, rec.Body, want)
		}
		// Install a draft (major 0), reference it in an introspection → can't delete.
		if rec := r.req(t, http.MethodPost, "/interfaces", draftV0, r.rmaToken); rec.Code != http.StatusCreated {
			t.Fatalf("install draft: got %d (%s)", rec.Code, rec.Body)
		}
		dev, _ := deviceid.Random()
		if err := r.st.RegisterDevice(ctx, r.realmID, dev, "h"); err != nil {
			t.Fatal(err)
		}
		if _, err := r.st.UpdateIntrospection(ctx, r.realmID, dev,
			map[string]store.InterfaceVersion{rmDraft: {Major: 0, Minor: 1}}); err != nil {
			t.Fatal(err)
		}
		wantInUse := `{"errors":{"detail":"Interface can't be deleted since it's currently used"}}`
		rec = r.req(t, http.MethodDelete, "/interfaces/"+rmDraft+"/0", "", r.rmaToken)
		if rec.Code != http.StatusForbidden || rec.Body.String() != wantInUse {
			t.Errorf("delete introspected draft: got %d %s, want 403 %s", rec.Code, rec.Body, wantInUse)
		}
		// Drop it from the introspection → now deletable.
		if _, err := r.st.UpdateIntrospection(ctx, r.realmID, dev, map[string]store.InterfaceVersion{}); err != nil {
			t.Fatal(err)
		}
		if rec := r.req(t, http.MethodDelete, "/interfaces/"+rmDraft+"/0", "", r.rmaToken); rec.Code != http.StatusNoContent {
			t.Errorf("delete unused draft: got %d, want 204", rec.Code)
		}
	})

	t.Run("Triggers", func(t *testing.T) {
		if rec := r.req(t, http.MethodPost, "/triggers", triggerJSON, r.rmaToken); rec.Code != http.StatusCreated {
			t.Fatalf("create trigger: got %d, want 201 (%s)", rec.Code, rec.Body)
		}
		if rec := r.req(t, http.MethodPost, "/triggers", triggerJSON, r.rmaToken); rec.Code != http.StatusConflict {
			t.Errorf("duplicate trigger: got %d, want 409", rec.Code)
		}
		if rec := r.req(t, http.MethodPost, "/triggers", `{"name":"bad","action":{},"simple_triggers":[]}`, r.rmaToken); rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("invalid trigger: got %d, want 422", rec.Code)
		}
		var names []string
		decodeData(t, r.req(t, http.MethodGet, "/triggers", "", r.rmaToken), &names)
		if !contains(names, "on_value") {
			t.Errorf("trigger list = %v, want to contain on_value", names)
		}
		if rec := r.req(t, http.MethodGet, "/triggers/on_value", "", r.rmaToken); rec.Code != http.StatusOK {
			t.Errorf("get trigger: got %d, want 200", rec.Code)
		}
		if rec := r.req(t, http.MethodDelete, "/triggers/on_value", "", r.rmaToken); rec.Code != http.StatusNoContent {
			t.Errorf("delete trigger: got %d, want 204", rec.Code)
		}
		if rec := r.req(t, http.MethodGet, "/triggers/on_value", "", r.rmaToken); rec.Code != http.StatusNotFound {
			t.Errorf("get deleted trigger: got %d, want 404", rec.Code)
		}
	})

	t.Run("ConfigRetention", func(t *testing.T) {
		// Upstream's default for a realm that never set a retention ceiling is
		// 0, and Astrate has no way to set one yet — so 0 is the only value.
		var retention int64
		decodeData(t, r.req(t, http.MethodGet, "/config/datastream_maximum_storage_retention", "", r.rmaToken), &retention)
		if retention != 0 {
			t.Errorf("datastream_maximum_storage_retention = %d, want 0", retention)
		}
	})

	t.Run("ConfigAuth", func(t *testing.T) {
		// Rotate to a 2-key set that still includes the original, so the test
		// token keeps verifying.
		rotated := pubPEM(t, &r.jwtKey.PublicKey) + "\n" + r.otherPub
		if rec := r.req(t, http.MethodPut, "/config/auth", `{"jwt_public_key_pem":`+jsonStr(rotated)+`}`, r.rmaToken); rec.Code != http.StatusNoContent {
			t.Fatalf("put config/auth: got %d, want 204 (%s)", rec.Code, rec.Body)
		}
		var cfg authConfig
		decodeData(t, r.req(t, http.MethodGet, "/config/auth", "", r.rmaToken), &cfg)
		if cfg.JWTPublicKeyPEM != rotated {
			t.Errorf("config/auth key not rotated")
		}
	})
}

// TestRealmManagementRetentionCeiling pins the #72 enforcement surface:
// installs and minor updates whose mapping TTL exceeds the realm's
// datastream_maximum_storage_retention answer upstream's named 422 envelope
// and leave nothing behind, a sub-ceiling install passes, and the config GET
// reports the stored ceiling (0 when unset — the #60 wire contract).
func TestRealmManagementRetentionCeiling(t *testing.T) {
	r := newRig(t)

	const retIface = "com.ex.M7a.Retention"
	retDef := func(minor, ttl int) string {
		return fmt.Sprintf(`{"interface_name":"%s","version_major":0,"version_minor":%d,"type":"datastream","ownership":"device",`+
			`"mappings":[{"endpoint":"/value","type":"double","database_retention_policy":"use_ttl","database_retention_ttl":%d}]}`,
			retIface, minor, ttl)
	}
	const overEnvelope = `{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`

	var retention int64
	decodeData(t, r.req(t, http.MethodGet, "/config/datastream_maximum_storage_retention", "", r.rmaToken), &retention)
	if retention != 0 {
		t.Errorf("retention before patch = %d, want 0", retention)
	}

	// Set the ceiling through the store: this suite has only RM; housekeeping
	// owns realm PATCH.
	if err := r.st.UpdateRealm(context.Background(), r.realm,
		store.RealmPatch{PatchRetention: true, SetRetention: 3600}); err != nil {
		t.Fatalf("UpdateRealm(retention=3600): %v", err)
	}
	decodeData(t, r.req(t, http.MethodGet, "/config/datastream_maximum_storage_retention", "", r.rmaToken), &retention)
	if retention != 3600 {
		t.Errorf("retention after patch = %d, want 3600", retention)
	}

	// Install above the ceiling → upstream's named 422, nothing installed.
	rec := r.req(t, http.MethodPost, "/interfaces", retDef(1, 7200), r.rmaToken)
	if rec.Code != http.StatusUnprocessableEntity || rec.Body.String() != overEnvelope {
		t.Errorf("install ttl 7200 under a 3600 ceiling: got %d %s, want 422 %s",
			rec.Code, rec.Body, overEnvelope)
	}
	if rec := r.req(t, http.MethodGet, "/interfaces/"+retIface+"/0", "", r.rmaToken); rec.Code != http.StatusNotFound {
		t.Errorf("get rejected install: got %d, want 404", rec.Code)
	}

	// Same shape within the ceiling installs fine.
	if rec := r.req(t, http.MethodPost, "/interfaces", retDef(1, 1800), r.rmaToken); rec.Code != http.StatusCreated {
		t.Fatalf("install ttl 1800: got %d, want 201 (%s)", rec.Code, rec.Body)
	}

	// A minor update raising the TTL past the ceiling is refused too, and
	// leaves the installed interface untouched at v0.1.
	rec = r.req(t, http.MethodPut, "/interfaces/"+retIface+"/0", retDef(2, 7200), r.rmaToken)
	if rec.Code != http.StatusUnprocessableEntity || rec.Body.String() != overEnvelope {
		t.Errorf("minor update to ttl 7200: got %d %s, want 422 %s",
			rec.Code, rec.Body, overEnvelope)
	}
	var majors []int
	decodeData(t, r.req(t, http.MethodGet, "/interfaces/"+retIface, "", r.rmaToken), &majors)
	if len(majors) != 1 || majors[0] != 0 {
		t.Errorf("majors after refused update = %v, want [0]", majors)
	}
}

// TestRealmManagementErrorCodes pins the install/update/delete error taxonomy
// measured against upstream 1.2 through the tunnel (#62): every rejection row
// compares the exact {"errors":{"detail":...}} body, interleaved with success
// twins proving nothing got stuck. Steps share one rig and run in order.
func TestRealmManagementErrorCodes(t *testing.T) {
	r := newRig(t)

	const (
		bodyDup           = `{"errors":{"detail":"Interface already exists"}}`
		bodyCollision     = `{"errors":{"detail":"Interface name collision detected. Make sure that the difference between two interface names is not limited to the casing or the presence of hyphens."}}`
		bodyNameMismatch  = `{"errors":{"detail":"Interface name doesn't match the one in the interface json"}}`
		bodyMajorMismatch = `{"errors":{"detail":"Interface major version doesn't match the one in the interface json"}}`
		bodyMinorSame     = `{"errors":{"detail":"Interface minor version was not increased"}}`
		bodyDowngrade     = `{"errors":{"detail":"Interface downgrade not allowed"}}`
		bodyMutated       = `{"errors":{"detail":"Interface update contains incompatible endpoint changes"}}`
		bodyDropped       = `{"errors":{"detail":"Interface update has missing endpoints"}}`
		bodyMajorNotFound = `{"errors":{"detail":"Interface major not found"}}`
		bodyNotFound      = `{"errors":{"detail":"Interface not found"}}`
		bodyCantDelete    = `{"errors":{"detail":"Interface can't be deleted"}}`
		bodyInUse         = `{"errors":{"detail":"Interface can't be deleted since it's currently used"}}`
		bodyTriggerDup    = `{"errors":{"detail":"Already exists"}}`
	)

	sensors := func(minor int, mappings string) string {
		return fmt.Sprintf(`{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":%d,`+
			`"type":"datastream","ownership":"device","mappings":[%s]}`, minor, mappings)
	}
	named := func(name string) string {
		return fmt.Sprintf(`{"interface_name":"%s","version_major":1,"version_minor":0,`+
			`"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`, name)
	}

	var draftDevice deviceid.ID
	setupDraftInUse := func(t *testing.T) {
		t.Helper()
		if rec := r.req(t, http.MethodPost, "/interfaces", draftV0, r.rmaToken); rec.Code != http.StatusCreated {
			t.Fatalf("setup install draft: got %d (%s)", rec.Code, rec.Body)
		}
		dev, _ := deviceid.Random()
		if err := r.st.RegisterDevice(context.Background(), r.realmID, dev, "h"); err != nil {
			t.Fatal(err)
		}
		draftDevice = dev
		if _, err := r.st.UpdateIntrospection(context.Background(), r.realmID, dev,
			map[string]store.InterfaceVersion{rmDraft: {Major: 0, Minor: 1}}); err != nil {
			t.Fatal(err)
		}
	}
	clearDraftIntrospection := func(t *testing.T) {
		t.Helper()
		if _, err := r.st.UpdateIntrospection(context.Background(), r.realmID, draftDevice,
			map[string]store.InterfaceVersion{}); err != nil {
			t.Fatal(err)
		}
	}

	steps := []struct {
		name       string
		method     string
		path       string
		body       string
		setup      func(t *testing.T)
		wantStatus int
		wantBody   string // "" = don't compare (success twins)
	}{
		{"get_unknown_name", http.MethodGet, "/interfaces/com.ex.M7a.Absent", "",
			nil, http.StatusNotFound, bodyNotFound},
		{"get_unknown_major", http.MethodGet, "/interfaces/com.ex.M7a.Absent/0", "",
			nil, http.StatusNotFound, bodyNotFound},

		{"fresh_install_twin", http.MethodPost, "/interfaces", ifaceV1,
			nil, http.StatusCreated, ""},
		{"duplicate_install", http.MethodPost, "/interfaces", ifaceV1,
			nil, http.StatusConflict, bodyDup},

		{"install_case_collision", http.MethodPost, "/interfaces", named("com.ex.M7A.sensors"),
			nil, http.StatusConflict, bodyCollision},
		// Upstream's collision key lowercases and strips '-' only (queries.ex
		// normalize_interface_name): dropping the dot yields a genuinely
		// distinct interface, while a hyphenated middle label collides.
		{"install_dotless_distinct_twin", http.MethodPost, "/interfaces", named("com.ex.M7aSensors"),
			nil, http.StatusCreated, ""},
		{"install_hyphen_collision", http.MethodPost, "/interfaces", named("com.ex.M-7a.Sensors"),
			nil, http.StatusConflict, bodyCollision},
		{"distinct_name_twin", http.MethodPost, "/interfaces", named("com.ex.M7a.Distinct"),
			nil, http.StatusCreated, ""},

		{"put_name_mismatch", http.MethodPut, "/interfaces/" + rmIface + "/1",
			named("com.ex.Mismatched"), nil, http.StatusConflict, bodyNameMismatch},
		{"put_major_mismatch", http.MethodPut, "/interfaces/" + rmIface + "/1",
			fmt.Sprintf(`{"interface_name":"%s","version_major":2,"version_minor":0,"type":"datastream",`+
				`"ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`, rmIface),
			nil, http.StatusConflict, bodyMajorMismatch},

		{"put_same_minor", http.MethodPut, "/interfaces/" + rmIface + "/1", sensors(0, `{"endpoint":"/value","type":"double"}`),
			nil, http.StatusConflict, bodyMinorSame},
		{"additive_twin_after_same_minor", http.MethodPut, "/interfaces/" + rmIface + "/1", ifaceV1b,
			nil, http.StatusNoContent, ""},
		{"put_downgrade", http.MethodPut, "/interfaces/" + rmIface + "/1", sensors(0, `{"endpoint":"/value","type":"double"}`),
			nil, http.StatusConflict, bodyDowngrade},
		{"put_mutated_mapping_type", http.MethodPut, "/interfaces/" + rmIface + "/1",
			sensors(2, `{"endpoint":"/value","type":"integer"},{"endpoint":"/count","type":"integer"}`),
			nil, http.StatusConflict, bodyMutated},
		{"additive_twin_after_mutation", http.MethodPut, "/interfaces/" + rmIface + "/1",
			sensors(2, `{"endpoint":"/value","type":"double"},{"endpoint":"/count","type":"integer"},{"endpoint":"/extra","type":"boolean"}`),
			nil, http.StatusNoContent, ""},
		{"put_dropped_endpoint", http.MethodPut, "/interfaces/" + rmIface + "/1", sensors(3, `{"endpoint":"/value","type":"double"}`),
			nil, http.StatusConflict, bodyDropped},

		{"put_uninstalled_major", http.MethodPut, "/interfaces/" + rmIface + "/7",
			fmt.Sprintf(`{"interface_name":"%s","version_major":7,"version_minor":0,"type":"datastream",`+
				`"ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`, rmIface),
			nil, http.StatusNotFound, bodyMajorNotFound},
		{"delete_unknown_interface", http.MethodDelete, "/interfaces/com.ex.M7a.Absent/0", "",
			nil, http.StatusNotFound, bodyMajorNotFound},
		{"delete_major_not_zero_unused", http.MethodDelete, "/interfaces/com.ex.M7a.Distinct/1", "",
			nil, http.StatusForbidden, bodyCantDelete},
		{"delete_introspected_draft", http.MethodDelete, "/interfaces/" + rmDraft + "/0", "",
			setupDraftInUse, http.StatusForbidden, bodyInUse},
		{"delete_draft_after_introspection_cleared", http.MethodDelete, "/interfaces/" + rmDraft + "/0", "",
			clearDraftIntrospection, http.StatusNoContent, ""},

		{"trigger_create_twin", http.MethodPost, "/triggers", triggerJSON,
			nil, http.StatusCreated, ""},
		{"duplicate_trigger_keeps_generic_detail", http.MethodPost, "/triggers", triggerJSON,
			nil, http.StatusConflict, bodyTriggerDup},
	}

	for _, step := range steps {
		t.Run(step.name, func(t *testing.T) {
			if step.setup != nil {
				step.setup(t)
			}
			rec := r.req(t, step.method, step.path, step.body, r.rmaToken)
			if rec.Code != step.wantStatus {
				t.Errorf("got status %d body %s, want %d %s", rec.Code, rec.Body, step.wantStatus, step.wantBody)
			}
			if step.wantBody != "" && rec.Body.String() != step.wantBody {
				t.Errorf("got body %s, want %s", rec.Body, step.wantBody)
			}
		})
	}
}

// TestRealmManagementViolationEnvelopes pins the probe-frozen 422 bodies of
// the #61 rules: exact bytes, deterministic field order, and the full-length
// aligned mappings array.
func TestRealmManagementViolationEnvelopes(t *testing.T) {
	r := newRig(t)

	t.Run("DescriptionTooLong", func(t *testing.T) {
		def := fmt.Sprintf(`{"interface_name":"com.ex.M7a.Long","version_major":0,"version_minor":1,`+
			`"type":"datastream","ownership":"device","description":"%s",`+
			`"mappings":[{"endpoint":"/value","type":"double"}]}`, strings.Repeat("x", 1001))
		want := `{"errors":{"description":["should be at most 1000 character(s)"]}}`
		rec := r.req(t, http.MethodPost, "/interfaces", def, r.rmaToken)
		if rec.Code != http.StatusUnprocessableEntity || rec.Body.String() != want {
			t.Errorf("got %d %s, want 422 %s", rec.Code, rec.Body, want)
		}
	})

	t.Run("SecondMappingViolatesAlignedArray", func(t *testing.T) {
		def := fmt.Sprintf(`{"interface_name":"com.ex.M7a.Two","version_major":0,"version_minor":1,`+
			`"type":"datastream","ownership":"device","mappings":[`+
			`{"endpoint":"/a","type":"double"},`+
			`{"endpoint":"/b","type":"double","description":"%s"}]}`, strings.Repeat("x", 1001))
		want := `{"errors":{"mappings":[{},{"description":["should be at most 1000 character(s)"]}]}}`
		rec := r.req(t, http.MethodPost, "/interfaces", def, r.rmaToken)
		if rec.Code != http.StatusUnprocessableEntity || rec.Body.String() != want {
			t.Errorf("got %d %s, want 422 %s", rec.Code, rec.Body, want)
		}
	})
}

// TestRealmManagementLegacyAliases pins issue #61 end-to-end: an interface
// posted entirely in upstream's legacy alias form installs, and GET renders
// the canonical fields because the service stores the canonical re-encoding.
func TestRealmManagementLegacyAliases(t *testing.T) {
	r := newRig(t)

	const legacy = `{"interface_name":"com.ex.M7a.Legacy","version_major":0,"version_minor":1,` +
		`"type":"properties","quality":"device","aggregate":false,` +
		`"mappings":[{"path":"/value","type":"string"}]}`
	if rec := r.req(t, http.MethodPost, "/interfaces", legacy, r.rmaToken); rec.Code != http.StatusCreated {
		t.Fatalf("install legacy-alias interface: got %d, want 201 (%s)", rec.Code, rec.Body)
	}

	var stored struct {
		Quality   any    `json:"quality"`
		Ownership string `json:"ownership"`
		Mappings  []struct {
			Path     any    `json:"path"`
			Endpoint string `json:"endpoint"`
		} `json:"mappings"`
	}
	decodeData(t, r.req(t, http.MethodGet, "/interfaces/com.ex.M7a.Legacy/0", "", r.rmaToken), &stored)
	if stored.Ownership != "device" || stored.Quality != nil {
		t.Errorf("stored ownership = %q quality = %v, want device / absent", stored.Ownership, stored.Quality)
	}
	if len(stored.Mappings) != 1 || stored.Mappings[0].Endpoint != "/value" || stored.Mappings[0].Path != nil {
		t.Errorf("stored mapping = %+v, want canonical endpoint /value without path", stored.Mappings)
	}
}

// TestRealmManagementTTLBand pins the probe-verified [60, 630720000) band on
// the wire: 59 is rejected with the exact upstream body, 60 installs.
func TestRealmManagementTTLBand(t *testing.T) {
	r := newRig(t)

	ttlDef := func(name string, ttl int) string {
		return fmt.Sprintf(`{"interface_name":"%s","version_major":0,"version_minor":1,`+
			`"type":"datastream","ownership":"device","mappings":[`+
			`{"endpoint":"/value","type":"double","database_retention_policy":"use_ttl",`+
			`"database_retention_ttl":%d}]}`, name, ttl)
	}

	low := ttlDef("com.ex.M7a.TTLLow", 59)
	want := `{"errors":{"mappings":[{"database_retention_ttl":["must be greater than or equal to 60"]}]}}`
	if rec := r.req(t, http.MethodPost, "/interfaces", low, r.rmaToken); rec.Code != http.StatusUnprocessableEntity || rec.Body.String() != want {
		t.Errorf("ttl 59: got %d %s, want 422 %s", rec.Code, rec.Body, want)
	}
	if rec := r.req(t, http.MethodPost, "/interfaces", ttlDef("com.ex.M7a.TTLMin", 60), r.rmaToken); rec.Code != http.StatusCreated {
		t.Errorf("ttl 60: got %d, want 201 (%s)", rec.Code, rec.Body)
	}
}

// TestRealmManagementAmqpTriggerAndPolicyGuards pins #64/#65 on the wire:
// an amqp-exchange trigger action is rejected at creation with 422 (upstream
// would accept it; Astrate deliberately refuses — validate-and-reject beats
// silent misbehavior), the plain HTTP twin still creates, and policy
// overlap/prefetch violations answer 422 with the probed upstream wording
// while a disjoint policy installs.
func TestRealmManagementAmqpTriggerAndPolicyGuards(t *testing.T) {
	r := newRig(t)

	t.Run("amqp action rejected", func(t *testing.T) {
		def := `{"name":"on_amqp","action":{"amqp_exchange":"astarte_events_` + r.realm + `",` +
			`"amqp_routing_key":"rk","amqp_queue":"q"},"simple_triggers":[{"type":"data_trigger",` +
			`"on":"incoming_data","interface_name":"com.ex.M7a.Sensors","interface_major":1,` +
			`"match_path":"/value","value_match_operator":"*"}]}`
		rec := r.req(t, http.MethodPost, "/triggers", def, r.rmaToken)
		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("amqp trigger: got %d (%s), want 422", rec.Code, rec.Body)
		}
		if body := rec.Body.String(); !strings.Contains(body, "amqp trigger actions are not supported") {
			t.Errorf("amqp trigger detail %q does not mention the rejection", body)
		}
	})

	t.Run("http trigger twin still creates", func(t *testing.T) {
		rec := r.req(t, http.MethodPost, "/triggers", triggerJSON, r.rmaToken)
		if rec.Code != http.StatusCreated {
			t.Errorf("create trigger: got %d (%s), want 201", rec.Code, rec.Body)
		}
	})

	t.Run("overlapping handlers rejected", func(t *testing.T) {
		def := `{"name":"overlap","error_handlers":[{"on":[400,401],"strategy":"discard"},` +
			`{"on":[401,402],"strategy":"discard"}],"maximum_capacity":1}`
		rec := r.req(t, http.MethodPost, "/policies", def, r.rmaToken)
		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("overlapping policy: got %d (%s), want 422", rec.Code, rec.Body)
		}
		if body := rec.Body.String(); !strings.Contains(body, "must all handle distinct errors") {
			t.Errorf("overlap detail %q does not mention disjointness", body)
		}
	})

	t.Run("prefetch_count out of band rejected", func(t *testing.T) {
		def := `{"name":"badprefetch","error_handlers":[{"on":"any_error","strategy":"discard"}],` +
			`"maximum_capacity":1,"prefetch_count":301}`
		rec := r.req(t, http.MethodPost, "/policies", def, r.rmaToken)
		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("prefetch_count 301: got %d (%s), want 422", rec.Code, rec.Body)
		}
		if body := rec.Body.String(); !strings.Contains(body, "prefetch_count") {
			t.Errorf("prefetch detail %q does not mention prefetch_count", body)
		}
	})

	t.Run("disjoint policy creates", func(t *testing.T) {
		def := `{"name":"disjoint-ok","error_handlers":[{"on":"server_error","strategy":"retry"},` +
			`{"on":[401],"strategy":"discard"}],"maximum_capacity":1,"retry_times":1}`
		if rec := r.req(t, http.MethodPost, "/policies", def, r.rmaToken); rec.Code != http.StatusCreated {
			t.Errorf("disjoint policy: got %d (%s), want 201", rec.Code, rec.Body)
		}
	})
}

// TestRealmManagementTriggerActionLimits pins the #63 upstream action
// validation limits on the wire: POST /triggers answers upstream's nested
// changeset 422 envelope with exact bytes, and an acceptance twin creates
// alongside the rejections to prove nothing else moved.
func TestRealmManagementTriggerActionLimits(t *testing.T) {
	r := newRig(t)

	trigger := func(name, action string) string {
		return `{"name":"` + name + `","action":` + action +
			`,"simple_triggers":[{"type":"data_trigger","on":"incoming_data",` +
			`"interface_name":"com.ex.M7a.Sensors","interface_major":1,` +
			`"match_path":"/value","value_match_operator":"*"}]}`
	}

	rows := []struct {
		name       string
		action     string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "short url and wrong scheme",
			action:     `{"http_url":"http://","http_method":"post"}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantBody:   `{"errors":{"action":{"http_url":["should be at least 8 character(s)","must be a valid http(s) URL"]}}}`,
		},
		{
			name:       "blocked header",
			action:     `{"http_url":"https://example.com/hook","http_method":"post","http_static_headers":{"Host":"evil.example"}}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantBody:   `{"errors":{"action":{"http_static_headers":["must contain only allowed http headers"]}}}`,
		},
		{
			name: "oversized header value",
			action: `{"http_url":"https://example.com/hook","http_method":"post","http_static_headers":{"X-A":"` +
				strings.Repeat("v", 8187) + `"}}`, // name(3) + ": "(2) + value == 8192 bytes
			wantStatus: http.StatusUnprocessableEntity,
			wantBody:   `{"errors":{"action":{"http_static_headers":["headers total size must be lower than 8192"]}}}`,
		},
		{
			name:       "bad method",
			action:     `{"http_url":"https://example.com/hook","http_method":"fetch"}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantBody:   `{"errors":{"action":{"http_method":["is invalid"]}}}`,
		},
	}
	for _, tc := range rows {
		t.Run(tc.name, func(t *testing.T) {
			rec := r.req(t, http.MethodPost, "/triggers", trigger("limit_"+strings.ReplaceAll(tc.name, " ", "_"), tc.action), r.rmaToken)
			if rec.Code != tc.wantStatus || rec.Body.String() != tc.wantBody {
				t.Errorf("got %d %s, want %d %s", rec.Code, rec.Body, tc.wantStatus, tc.wantBody)
			}
		})
	}

	t.Run("valid trigger twin creates", func(t *testing.T) {
		if rec := r.req(t, http.MethodPost, "/triggers", triggerJSON, r.rmaToken); rec.Code != http.StatusCreated {
			t.Errorf("create trigger: got %d (%s), want 201", rec.Code, rec.Body)
		}
	})
}

// --- helpers ----------------------------------------------------------------

func discardLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func mintToken(t *testing.T, key *rsa.PrivateKey, claims jwt.MapClaims) string {
	t.Helper()
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	s, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

func pubPEM(t *testing.T, pub *rsa.PublicKey) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
}

func decodeData(t *testing.T, rec *httptest.ResponseRecorder, dst any) {
	t.Helper()
	if rec.Code/100 != 2 {
		t.Fatalf("non-2xx response %d: %s", rec.Code, rec.Body)
	}
	var env struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode envelope: %v (%s)", err, rec.Body)
	}
	if err := json.Unmarshal(env.Data, dst); err != nil {
		t.Fatalf("decode data: %v (%s)", err, env.Data)
	}
}

func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

func jsonStr(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
