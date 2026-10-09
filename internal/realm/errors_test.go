package realm

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// rmWriter names one of the two RM error-mapping entry points so the
// containers-free tables below can drive both through one runner. Neither
// method dereferences a.svc, so a zero-value *API is enough.
type rmWriter func(*API, http.ResponseWriter, error)

func callWriter(fn rmWriter, err error) (int, string) {
	rec := httptest.NewRecorder()
	fn(&API{}, rec, err)
	return rec.Code, rec.Body.String()
}

// retentionBody is the named 422 envelope shared by both mapping layers.
const retentionBody = `{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`

// TestWriteInterfaceErrorMapping pins every arm of writeInterfaceError
// (http.go:352-399) to the probe-frozen upstream body, including the switch
// precedence and the 403-vs-409/404 split that only the integration suite
// exercises today.
func TestWriteInterfaceErrorMapping(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{"retention ceiling", ErrMaximumDatabaseRetentionExceeded,
			http.StatusUnprocessableEntity, retentionBody},
		{"violations envelope", &interfaceschema.ViolationsError{
			Violations:   []interfaceschema.Violation{{Field: "description", Messages: []string{"too long"}, MappingIndex: -1}},
			MappingCount: 1,
		}, http.StatusUnprocessableEntity, `{"errors":{"description":["too long"]}}`},
		{"major not found", errMajorNotFound, http.StatusNotFound,
			`{"errors":{"detail":"Interface major not found"}}`},
		{"name collision", ErrNameCollision, http.StatusConflict,
			`{"errors":{"detail":"Interface name collision detected. Make sure that the difference between two interface names is not limited to the casing or the presence of hyphens."}}`},
		{"name mismatch", ErrNameMismatch, http.StatusConflict,
			`{"errors":{"detail":"Interface name doesn't match the one in the interface json"}}`},
		{"major mismatch", ErrMajorMismatch, http.StatusConflict,
			`{"errors":{"detail":"Interface major version doesn't match the one in the interface json"}}`},
		{"minor not increased", interfaceschema.ErrMinorNotIncreased, http.StatusConflict,
			`{"errors":{"detail":"Interface minor version was not increased"}}`},
		{"downgrade not allowed", interfaceschema.ErrDowngradeNotAllowed, http.StatusConflict,
			`{"errors":{"detail":"Interface downgrade not allowed"}}`},
		{"missing endpoints", interfaceschema.ErrMissingEndpoints, http.StatusConflict,
			`{"errors":{"detail":"Interface update has missing endpoints"}}`},
		{"incompatible endpoint change", interfaceschema.ErrIncompatibleEndpointChange, http.StatusConflict,
			`{"errors":{"detail":"Interface update contains incompatible endpoint changes"}}`},
		{"already exists", store.ErrAlreadyExists, http.StatusConflict,
			`{"errors":{"detail":"Interface already exists"}}`},
		{"major not zero is forbidden", store.ErrInterfaceMajorNotZero, http.StatusForbidden,
			`{"errors":{"detail":"Interface can't be deleted"}}`},
		{"interface in use is forbidden", store.ErrInterfaceInUse, http.StatusForbidden,
			`{"errors":{"detail":"Interface can't be deleted since it's currently used"}}`},
		{"validation detail stripped", fmt.Errorf("%w: endpoint is missing", ErrValidation), http.StatusUnprocessableEntity,
			`{"errors":{"detail":"endpoint is missing"}}`},
		{"not found", store.ErrNotFound, http.StatusNotFound,
			`{"errors":{"detail":"Interface not found"}}`},
		{"unknown falls to 500", errors.New("boom"), http.StatusInternalServerError,
			`{"errors":{"detail":"Internal server error"}}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, body := callWriter((*API).writeInterfaceError, tc.err)
			if code != tc.wantStatus || body != tc.wantBody {
				t.Errorf("got %d %s, want %d %s", code, body, tc.wantStatus, tc.wantBody)
			}
		})
	}
}

// TestWriteErrorMapping pins the generic writeError arms (http.go:402-425),
// notably that it maps the two interface delete sentinels to 422 rather than
// writeInterfaceError's 403, and that it falls back to the generic 404 body.
func TestWriteErrorMapping(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{"retention ceiling", ErrMaximumDatabaseRetentionExceeded,
			http.StatusUnprocessableEntity, retentionBody},
		{"violations envelope", &interfaceschema.ViolationsError{
			Violations:   []interfaceschema.Violation{{Field: "description", Messages: []string{"too long"}, MappingIndex: -1}},
			MappingCount: 1,
		}, http.StatusUnprocessableEntity, `{"errors":{"description":["too long"]}}`},
		{"validation detail stripped", fmt.Errorf("%w: trigger requires a name", ErrValidation), http.StatusUnprocessableEntity,
			`{"errors":{"detail":"trigger requires a name"}}`},
		{"already exists", store.ErrAlreadyExists, http.StatusConflict,
			`{"errors":{"detail":"Already exists"}}`},
		{"major not zero is 422", store.ErrInterfaceMajorNotZero, http.StatusUnprocessableEntity,
			`{"errors":{"detail":"Interface major version is not 0, can't be deleted"}}`},
		{"interface in use is 422", store.ErrInterfaceInUse, http.StatusUnprocessableEntity,
			`{"errors":{"detail":"Cannot delete an interface that is used by a device introspection"}}`},
		{"not found is generic", store.ErrNotFound, http.StatusNotFound,
			`{"errors":{"detail":"Not Found"}}`},
		{"unknown falls to 500", errors.New("boom"), http.StatusInternalServerError,
			`{"errors":{"detail":"Internal server error"}}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			code, body := callWriter((*API).writeError, tc.err)
			if code != tc.wantStatus || body != tc.wantBody {
				t.Errorf("got %d %s, want %d %s", code, body, tc.wantStatus, tc.wantBody)
			}
		})
	}
}

// TestRenderViolationsBody pins the byte-exact wire rendering (http.go:443-533):
// interface-level fields merged per field in first-seen order, then the
// full-length index-aligned mappings array with {} for clean entries.
func TestRenderViolationsBody(t *testing.T) {
	t.Run("mixed top-level and aligned mappings", func(t *testing.T) {
		ve := &interfaceschema.ViolationsError{
			Violations: []interfaceschema.Violation{
				{Field: "description", Messages: []string{"first"}, MappingIndex: -1},
				{Field: "description", Messages: []string{"second"}, MappingIndex: -1},
				{Field: "interface_name", Messages: []string{"bad name"}, MappingIndex: -1},
				{Field: "type", Messages: []string{"bad type", "again"}, MappingIndex: 1},
				{Field: "endpoint", Messages: []string{"bad endpoint"}, MappingIndex: 1},
			},
			MappingCount: 3,
		}
		want := `{"errors":{"description":["first","second"],"interface_name":["bad name"],"mappings":[{},{"type":["bad type","again"],"endpoint":["bad endpoint"]},{}]}}`
		if got := string(renderViolationsBody(ve)); got != want {
			t.Errorf("got  %s\nwant %s", got, want)
		}
	})

	t.Run("top-level only omits mappings array", func(t *testing.T) {
		ve := &interfaceschema.ViolationsError{
			Violations: []interfaceschema.Violation{
				{Field: "type", Messages: []string{"required"}, MappingIndex: -1},
			},
		}
		want := `{"errors":{"type":["required"]}}`
		if got := string(renderViolationsBody(ve)); got != want {
			t.Errorf("got  %s\nwant %s", got, want)
		}
	})

	t.Run("mappings only is full length with clean holes", func(t *testing.T) {
		ve := &interfaceschema.ViolationsError{
			Violations: []interfaceschema.Violation{
				{Field: "endpoint", Messages: []string{"bad"}, MappingIndex: 2},
			},
			MappingCount: 4,
		}
		want := `{"errors":{"mappings":[{},{},{"endpoint":["bad"]},{}]}}`
		if got := string(renderViolationsBody(ve)); got != want {
			t.Errorf("got  %s\nwant %s", got, want)
		}
	})
}

// TestValidationDetail pins the ErrValidation prefix-strip rule (http.go:560-569),
// including the strict len-guard so a message exactly equal to the prefix is
// left intact.
func TestValidationDetail(t *testing.T) {
	const prefix = "realm: validation failed: "
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"prefix stripped", fmt.Errorf("%w: endpoint is missing", ErrValidation), "endpoint is missing"},
		{"non-prefix passthrough", errors.New("plain failure"), "plain failure"},
		{"exact prefix not stripped", errors.New(prefix), prefix},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := validationDetail(tc.err); got != tc.want {
				t.Errorf("validationDetail(%v) = %q, want %q", tc.err, got, tc.want)
			}
		})
	}
}
