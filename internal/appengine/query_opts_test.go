package appengine

import (
	"errors"
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
		{"minimum valid", "downsample_to=3", 3, false},
		{"large value", "downsample_to=1000", 1000, false},
		// Upstream's changeset rejects 2 with "must be greater than 2".
		{"two", "downsample_to=2", 0, true},
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
			if err != nil && tt.wantErr {
				var fe FieldErrors
				if !errors.As(err, &fe) || len(fe["downsample_to"]) == 0 {
					t.Errorf("err = %v, want downsample_to field error", err)
				}
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

func TestParseQueryOptsFieldErrors(t *testing.T) {
	tests := []struct {
		name   string
		raw    string
		fields FieldErrors
	}{
		{"bad since", "since=nope", FieldErrors{"since": {"is invalid"}}},
		{"bad to", "to=yesterday", FieldErrors{"to": {"is invalid"}}},
		{
			"since_after conflicts with since",
			"since=2026-01-01T00:00:00Z&since_after=2026-01-02T00:00:00Z",
			FieldErrors{"since_after": {"conflicts already set parameter"}},
		},
		{"negative limit", "limit=-1", FieldErrors{"limit": {"must be greater than or equal to 0"}}},
		{"non numeric limit", "limit=abc", FieldErrors{"limit": {"is invalid"}}},
		{"two rejected", "downsample_to=2&limit=-1", FieldErrors{
			"downsample_to": {"must be greater than 2"},
			"limit":         {"must be greater than or equal to 0"},
		}},
		{"bogus format", "format=bogus", FieldErrors{"format": {"is invalid"}}},
		{"non bool allow_bigintegers", "allow_bigintegers=yes", FieldErrors{"allow_bigintegers": {"is invalid"}}},
		{"non bool allow_safe_bigintegers", "allow_safe_bigintegers=1", FieldErrors{"allow_safe_bigintegers": {"is invalid"}}},
		{"non bool retrieve_metadata", "retrieve_metadata=maybe", FieldErrors{"retrieve_metadata": {"is invalid"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/x?"+tt.raw, nil)
			opts, err := parseQueryOpts(r)
			var fe FieldErrors
			if !errors.As(err, &fe) {
				t.Fatalf("err = %v (%T), want FieldErrors", err, err)
			}
			for field, msgs := range tt.fields {
				got, ok := fe[field]
				if !ok {
					t.Fatalf("missing field %q in %v", field, fe)
				}
				if len(got) != len(msgs) || got[0] != msgs[0] {
					t.Errorf("%s = %v, want %v", field, got, msgs)
				}
				delete(fe, field)
			}
			if len(fe) > 0 {
				t.Errorf("unexpected extra fields %v", fe)
			}
			if opts.Format == "" {
				t.Error("Format empty even on rejection")
			}
		})
	}
}

func TestParseQueryOptsNewParams(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/x?format=table&retrieve_metadata=true"+
		"&allow_bigintegers=false&allow_safe_bigintegers=true&downsample_key=abc&limit=5", nil)
	opts, err := parseQueryOpts(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Format != "table" {
		t.Errorf("Format=%q, want table", opts.Format)
	}
	if !opts.RetrieveMetadata {
		t.Error("RetrieveMetadata=false, want true")
	}
	if opts.AllowBigIntegers == nil || *opts.AllowBigIntegers {
		t.Errorf("AllowBigIntegers=%v, want false", opts.AllowBigIntegers)
	}
	if opts.AllowSafeBigIntegers == nil || !*opts.AllowSafeBigIntegers {
		t.Errorf("AllowSafeBigIntegers=%v, want true", opts.AllowSafeBigIntegers)
	}
	if opts.DownsampleKey != "abc" {
		t.Errorf("DownsampleKey=%q, want abc", opts.DownsampleKey)
	}

	// Absent options take their defaults.
	r = httptest.NewRequest(http.MethodGet, "/x", nil)
	opts, err = parseQueryOpts(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Format != "structured" {
		t.Errorf("default Format=%q, want structured", opts.Format)
	}
	if opts.AllowBigIntegers != nil || opts.AllowSafeBigIntegers != nil {
		t.Error("absent bigint flags must stay nil")
	}
	if opts.RetrieveMetadata {
		t.Error("RetrieveMetadata default true, want false")
	}

	// An empty format value is accepted and means structured.
	r = httptest.NewRequest(http.MethodGet, "/x?format=", nil)
	opts, err = parseQueryOpts(r)
	if err != nil {
		t.Fatalf("empty format rejected: %v", err)
	}
	if opts.Format != "structured" {
		t.Errorf("empty Format=%q, want structured", opts.Format)
	}

	// since_after alone is legal.
	r = httptest.NewRequest(http.MethodGet, "/x?since_after=2026-01-02T00:00:00Z", nil)
	if _, err = parseQueryOpts(r); err != nil {
		t.Errorf("since_after alone rejected: %v", err)
	}
}
