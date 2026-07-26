package appengine

import "time"

// bucketFor derives the time_bucket width that yields about points samples over a
// series spanning span.
//
// We round up so that span/bucket ≤ points: a truncated bucket would produce more
// points than the client asked for, which is the one outcome a caller requesting N
// points does not want.
func bucketFor(span time.Duration, points int) time.Duration {
	if points <= 0 {
		return 0
	}
	if span <= 0 {
		return 0
	}

	bucket := span / time.Duration(points)
	if bucket*time.Duration(points) < span {
		bucket++
	}

	const micro = time.Microsecond
	if bucket < micro {
		bucket = micro
	}
	return bucket
}
