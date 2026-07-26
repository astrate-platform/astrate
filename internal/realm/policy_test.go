package realm

import (
	"errors"
	"testing"
)

func TestValidatePolicy(t *testing.T) {
	cases := map[string]struct {
		def      string
		wantName string
		wantErr  bool
	}{
		"dashboard retry form": {
			def:      `{"name":"retry5xx","error_handlers":[{"on":"server_error","strategy":"retry"}],"maximum_capacity":100,"retry_times":3}`,
			wantName: "retry5xx",
		},
		"dashboard discard form": {
			def:      `{"name":"drop.all","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1}`,
			wantName: "drop.all",
		},
		"status code array": {
			def:      `{"name":"codes","error_handlers":[{"on":[500,503],"strategy":"retry"}],"maximum_capacity":10,"retry_times":10}`,
			wantName: "codes",
		},
		"custom_status_codes object form": {
			def:      `{"name":"custom","error_handlers":[{"on":{"custom_status_codes":[418]},"strategy":"discard"}],"maximum_capacity":5}`,
			wantName: "custom",
		},
		"event_ttl ok": {
			def:      `{"name":"ttl","error_handlers":[{"on":"client_error","strategy":"discard"}],"maximum_capacity":5,"event_ttl":60}`,
			wantName: "ttl",
		},

		"not json":                  {def: `nope`, wantErr: true},
		"empty name":                {def: `{"name":"","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1}`, wantErr: true},
		"reserved @ name":           {def: `{"name":"@default","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1}`, wantErr: true},
		"name bad chars":            {def: `{"name":"a b","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1}`, wantErr: true},
		"no handlers":               {def: `{"name":"x","error_handlers":[],"maximum_capacity":1}`, wantErr: true},
		"bad strategy":              {def: `{"name":"x","error_handlers":[{"on":"any_error","strategy":"panic"}],"maximum_capacity":1}`, wantErr: true},
		"bad keyword":               {def: `{"name":"x","error_handlers":[{"on":"sometimes","strategy":"discard"}],"maximum_capacity":1}`, wantErr: true},
		"code out of range":         {def: `{"name":"x","error_handlers":[{"on":[200],"strategy":"discard"}],"maximum_capacity":1}`, wantErr: true},
		"empty code list":           {def: `{"name":"x","error_handlers":[{"on":[],"strategy":"discard"}],"maximum_capacity":1}`, wantErr: true},
		"retry without retry_times": {def: `{"name":"x","error_handlers":[{"on":"any_error","strategy":"retry"}],"maximum_capacity":1}`, wantErr: true},
		"retry_times without retry": {def: `{"name":"x","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1,"retry_times":3}`, wantErr: true},
		"retry_times over 100":      {def: `{"name":"x","error_handlers":[{"on":"any_error","strategy":"retry"}],"maximum_capacity":1,"retry_times":101}`, wantErr: true},
		"zero capacity":             {def: `{"name":"x","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":0}`, wantErr: true},
		"negative event_ttl":        {def: `{"name":"x","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":1,"event_ttl":-1}`, wantErr: true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := validatePolicy([]byte(tc.def))
			if tc.wantErr {
				if err == nil {
					t.Fatalf("accepted invalid policy, name %q", got)
				}
				if !errors.Is(err, ErrValidation) {
					t.Errorf("error %v does not wrap ErrValidation", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("rejected valid policy: %v", err)
			}
			if got != tc.wantName {
				t.Errorf("name = %q, want %q", got, tc.wantName)
			}
		})
	}
}
