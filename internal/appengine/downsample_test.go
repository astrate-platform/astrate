package appengine

import (
	"testing"
	"time"
)

func TestBucketFor(t *testing.T) {
	tests := []struct {
		name   string
		span   time.Duration
		points int
		want   time.Duration // 0 means "check property only" when span>0 && points>0
	}{
		{
			name:   "exact division",
			span:   60 * time.Second,
			points: 6,
			want:   10 * time.Second,
		},
		{
			name:   "rounds up",
			span:   10 * time.Second,
			points: 3,
			want:   0, // property: span/bucket must not exceed points
		},
		{
			name:   "points zero",
			span:   60 * time.Second,
			points: 0,
			want:   0,
		},
		{
			name:   "points negative",
			span:   60 * time.Second,
			points: -5,
			want:   0,
		},
		{
			name:   "span zero",
			span:   0,
			points: 10,
			want:   0,
		},
		{
			name:   "span negative",
			span:   -time.Hour,
			points: 10,
			want:   0,
		},
		{
			name:   "tiny span clamps to microsecond",
			span:   1 * time.Nanosecond,
			points: 1000,
			want:   0, // property: bucket >= 1µs, span/bucket ≤ points
		},
		{
			name:   "one year no overflow",
			span:   365 * 24 * time.Hour,
			points: 2,
			want:   0, // property: span/bucket ≤ points
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bucketFor(tt.span, tt.points)

			if tt.want != 0 {
				if got != tt.want {
					t.Errorf("bucketFor(%v, %d) = %v, want %v", tt.span, tt.points, got, tt.want)
				}
				return
			}

			if tt.points <= 0 || tt.span <= 0 {
				if got != 0 {
					t.Errorf("bucketFor(%v, %d) = %v, want 0", tt.span, tt.points, got)
				}
				return
			}

			if got < time.Microsecond {
				t.Errorf("bucketFor(%v, %d) = %v, want ≥ 1µs", tt.span, tt.points, got)
			}

			buckets := (tt.span + got - 1) / got
			if buckets > time.Duration(tt.points) {
				t.Errorf("ceil(span/bucket) = %d > %d points", buckets, tt.points)
			}
		})
	}
}
