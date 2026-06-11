package pairing

import (
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

// limiterMaxKeys bounds the number of tracked buckets; least-recently-used
// keys are evicted, so a brute-forcer cannot grow memory without bound
// (docs/DESIGN.md §4.5).
const limiterMaxKeys = 65536

// Limiter is a keyed token-bucket rate limiter for the pairing endpoints
// (docs/DESIGN.md §4.4–4.5): one bucket per key (client IP, device address),
// refilled at rate tokens/second up to burst.
type Limiter struct {
	mu      sync.Mutex
	buckets *lru.Cache[string, *bucket]
	rate    float64
	burst   float64
	now     func() time.Time
}

// bucket is one token bucket; guarded by the owning Limiter's mutex.
type bucket struct {
	tokens float64
	last   time.Time
}

// NewLimiter builds a Limiter refilling rate tokens per second with the
// given burst capacity. Non-positive parameters are clamped to minimal
// sane values (1 token/minute, burst 1).
func NewLimiter(rate float64, burst int) *Limiter {
	if rate <= 0 {
		rate = 1.0 / 60
	}
	if burst < 1 {
		burst = 1
	}
	// lru.New only fails for sizes < 1; limiterMaxKeys is a positive const.
	cache, _ := lru.New[string, *bucket](limiterMaxKeys)
	return &Limiter{
		buckets: cache,
		rate:    rate,
		burst:   float64(burst),
		now:     time.Now,
	}
}

// Allow reports whether one request under key may proceed, consuming a
// token when it does.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	b, ok := l.buckets.Get(key)
	if !ok {
		b = &bucket{tokens: l.burst, last: now}
		l.buckets.Add(key, b)
	}

	elapsed := now.Sub(b.last).Seconds()
	if elapsed > 0 {
		b.tokens += elapsed * l.rate
		if b.tokens > l.burst {
			b.tokens = l.burst
		}
		b.last = now
	}

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}
