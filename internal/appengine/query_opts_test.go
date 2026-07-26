package appengine

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestParseQueryOptsDownsampleTo(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantPts int
		wantErr bool
	}{
		{"empty query", "", 0, false},
		{"present but empty", "downsample_to=", 0, false},
		{"minimum valid", "downsample_to=2", 2, false},
		{"large value", "downsample_to=1000", 1000, false},
		{"too small 1", "downsample_to=1", 0, true},
		{"zero", "downsample_to=0", 0, true},
		{"negative", "downsample_to=-5", 0, true},
		{"not a number", "downsample_to=abc", 0, true},
		{"float", "downsample_to=1.5", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/x?"+tt.raw, nil)
			opts, err := parseQueryOpts(r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err=%v, wantErr=%v", err, tt.wantErr)
			}
			if opts.DownsamplePoints != tt.wantPts {
				t.Errorf("DownsamplePoints=%d, want %d", opts.DownsamplePoints, tt.wantPts)
			}
		})
	}
}

func TestParseQueryOptsDownsampleAlongside(t *testing.T) {
	raw := "since=2026-01-01T00:00:00Z&limit=10&sort=ascending&downsample_to=50"
	r := httptest.NewRequest(http.MethodGet, "/x?"+raw, nil)
	opts, err := parseQueryOpts(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.DownsamplePoints != 50 {
		t.Errorf("DownsamplePoints=%d, want 50", opts.DownsamplePoints)
	}
	if opts.Limit != 10 {
		t.Errorf("Limit=%d, want 10", opts.Limit)
	}
	if opts.Descending {
		t.Error("Descending=true, want false")
	}
	if opts.Since == nil || !opts.Since.Equal(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("Since=%v, want 2026-01-01T00:00:00Z", opts.Since)
	}
}
