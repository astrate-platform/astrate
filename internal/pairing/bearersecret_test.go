package pairing

import (
	"net/http"
	"testing"
)

func TestBearerSecret(t *testing.T) {
	tests := []struct {
		name       string
		auth       string
		wantSecret string
		wantOK     bool
	}{
		// happy paths — all must return the same secret
		{"standard", "Bearer secretvalue", "secretvalue", true},
		{"lowercase", "bearer secretvalue", "secretvalue", true},
		{"uppercase", "BEARER secretvalue", "secretvalue", true},
		{"mixed case", "BeArEr secretvalue", "secretvalue", true},

		// colon form — upstream ~r/bearer\:?\s+(.*)$/i
		{"colon", "Bearer: secretvalue", "secretvalue", true},
		{"colon lowercase", "bearer: secretvalue", "secretvalue", true},
		{"colon uppercase", "BEARER: secretvalue", "secretvalue", true},

		// extra whitespace around the value is trimmed
		{"extra spaces", "Bearer   secretvalue  ", "secretvalue", true},
		{"colon extra spaces", "Bearer:   secretvalue  ", "secretvalue", true},

		// reject paths
		{"empty header", "", "", false},
		{"no space", "Bearersecret", "", false},
		{"non-bearer scheme", "Basic dXNlcjpwYXNz", "", false},
		{"empty value", "Bearer ", "", false},
		{"empty value colon", "Bearer: ", "", false},
		{"whitespace only value", "Bearer   ", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.auth != "" {
				req.Header.Set("Authorization", tt.auth)
			}
			got, ok := bearerSecret(req)
			if ok != tt.wantOK {
				t.Errorf("ok = %v, want %v", ok, tt.wantOK)
			}
			if got != tt.wantSecret {
				t.Errorf("secret = %q, want %q", got, tt.wantSecret)
			}
		})
	}
}
