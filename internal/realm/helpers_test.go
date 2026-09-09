package realm

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/astrate-platform/astrate/internal/engine/triggers"
)

func TestNormaliseIfaceName(t *testing.T) {
	cases := map[string]struct {
		input string
		want  string
	}{
		"lowercase":          {input: "com.example.sensors", want: "com.example.sensors"},
		"uppercase":          {input: "Com.Example.Sensors", want: "com.example.sensors"},
		"hyphen strip":       {input: "com-example-sensors", want: "comexamplesensors"},
		"mixed case+hyphens": {input: "Com-Example-Sensors", want: "comexamplesensors"},
		"empty":              {input: "", want: ""},
		"already normalised": {input: "com.example.v1", want: "com.example.v1"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := normaliseIfaceName(tc.input)
			if got != tc.want {
				t.Errorf("normaliseIfaceName(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestJoinPEM(t *testing.T) {
	cases := map[string]struct {
		keys []string
		want string
	}{
		"empty":    {keys: nil, want: ""},
		"single":   {keys: []string{"-----BEGIN PUBLIC KEY-----\nabc\n-----END PUBLIC KEY-----"}, want: "-----BEGIN PUBLIC KEY-----\nabc\n-----END PUBLIC KEY-----"},
		"two keys": {keys: []string{"key1", "key2"}, want: "key1\nkey2"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := joinPEM(tc.keys)
			if got != tc.want {
				t.Errorf("joinPEM(...) = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTriggerErrorBody(t *testing.T) {
	cases := map[string]struct {
		input *triggers.TriggerErrors
		want  map[string]any
	}{
		"both empty": {
			input: &triggers.TriggerErrors{},
			want:  map[string]any{},
		},
		"action only": {
			input: &triggers.TriggerErrors{
				Action: map[string][]string{"http_url": {"invalid url"}},
			},
			want: map[string]any{
				"action": map[string][]string{"http_url": {"invalid url"}},
			},
		},
		"simple_triggers only with index-aligned empty": {
			input: &triggers.TriggerErrors{
				SimpleTriggers: []map[string][]string{
					{"type": {"bad type"}},
					nil,
				},
			},
			want: map[string]any{
				"simple_triggers": []map[string][]string{
					{"type": {"bad type"}},
					{},
				},
			},
		},
		"both present": {
			input: &triggers.TriggerErrors{
				Action:         map[string][]string{"http_method": {"put"}},
				SimpleTriggers: []map[string][]string{{"on": {"bad val"}}},
			},
			want: map[string]any{
				"action": map[string][]string{"http_method": {"put"}},
				"simple_triggers": []map[string][]string{
					{"on": {"bad val"}},
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := triggerErrorBody(tc.input)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("triggerErrorBody(...) = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestMajorParam(t *testing.T) {
	cases := map[string]struct {
		major    string
		wantVal  int
		wantOK   bool
		wantCode int
	}{
		"valid positive": {major: "3", wantVal: 3, wantOK: true},
		"zero":           {major: "0", wantVal: 0, wantOK: true},
		"negative":       {major: "-1", wantCode: http.StatusNotFound},
		"non-numeric":    {major: "abc", wantCode: http.StatusNotFound},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			r.SetPathValue("major", tc.major)
			val, ok := majorParam(w, r)
			if ok != tc.wantOK {
				t.Fatalf("majorParam ok = %v, want %v", ok, tc.wantOK)
			}
			if ok && val != tc.wantVal {
				t.Errorf("majorParam val = %d, want %d", val, tc.wantVal)
			}
			if tc.wantCode != 0 && w.Code != tc.wantCode {
				t.Errorf("response status = %d, want %d", w.Code, tc.wantCode)
			}
		})
	}
}
