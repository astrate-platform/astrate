package flowapi

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/store"
)

func TestWriteErrorMapping(t *testing.T) {
	a := &API{}
	cases := []struct {
		err  error
		code int
		sub  string
	}{
		{fmt.Errorf("%w: bad", ErrValidation), http.StatusUnprocessableEntity, "bad"},
		{flow.ErrFlowExists, http.StatusConflict, "already"},
		{flow.ErrFlowNotFound, http.StatusNotFound, "Not Found"},
		{store.ErrAlreadyExists, http.StatusConflict, "Already"},
		{store.ErrPipelineCyclic, http.StatusUnprocessableEntity, "cycle"},
		{store.ErrNotFound, http.StatusNotFound, "Not Found"},
		{errors.New("other"), http.StatusInternalServerError, "Internal"},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		a.writeError(rec, tc.err)
		if rec.Code != tc.code {
			t.Errorf("%v: status %d, want %d body=%s", tc.err, rec.Code, tc.code, rec.Body.String())
		}
		if !strings.Contains(strings.ToLower(rec.Body.String()), strings.ToLower(tc.sub)) {
			t.Errorf("%v: body %q should contain %q", tc.err, rec.Body.String(), tc.sub)
		}
	}
}

func TestValidationDetail(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "strips prefix",
			err:  fmt.Errorf("flowapi: validation failed: bad block"),
			want: "bad block",
		},
		{
			name: "non-prefix passes through",
			err:  errors.New("something else"),
			want: "something else",
		},
		{
			name: "prefix only not stripped",
			err:  errors.New("flowapi: validation failed: "),
			want: "flowapi: validation failed: ",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := validationDetail(tc.err)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
