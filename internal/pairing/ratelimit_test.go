package pairing

import (
	"fmt"
	"testing"
	"time"
)

func TestLimiterBurstThenDeny(t *testing.T) {
	l := NewLimiter(1, 3)
	clock := time.Now()
	l.now = func() time.Time { return clock }

	for i := 0; i < 3; i++ {
		if !l.Allow("ip|192.0.2.1") {
			t.Fatalf("request %d within burst must be allowed", i)
		}
	}
	if l.Allow("ip|192.0.2.1") {
		t.Fatal("request beyond burst must be denied")
	}

	// Refill: after 1 second at 1 token/s exactly one more request passes.
	clock = clock.Add(time.Second)
	if !l.Allow("ip|192.0.2.1") {
		t.Fatal("request after refill must be allowed")
	}
	if l.Allow("ip|192.0.2.1") {
		t.Fatal("second request after a 1-token refill must be denied")
	}
}

func TestLimiterKeysAreIndependent(t *testing.T) {
	l := NewLimiter(1, 1)
	clock := time.Now()
	l.now = func() time.Time { return clock }

	if !l.Allow("ip|192.0.2.1") {
		t.Fatal("first key must be allowed")
	}
	if !l.Allow("ip|192.0.2.2") {
		t.Fatal("a fresh key must have its own bucket")
	}
	if !l.Allow("dev|test/x") {
		t.Fatal("device keys are independent of IP keys")
	}
	if l.Allow("ip|192.0.2.1") {
		t.Fatal("exhausted key must stay denied")
	}
}

func TestLimiterCapsTokensAtBurst(t *testing.T) {
	l := NewLimiter(100, 2)
	clock := time.Now()
	l.now = func() time.Time { return clock }

	if !l.Allow("k") {
		t.Fatal("initial request")
	}
	// A long idle period must not accumulate more than burst tokens.
	clock = clock.Add(time.Hour)
	allowed := 0
	for i := 0; i < 10; i++ {
		if l.Allow("k") {
			allowed++
		}
	}
	if allowed != 2 {
		t.Fatalf("after long idle: got %d allowed, want burst=2", allowed)
	}
}

func TestLimiterEvictionBoundsMemory(t *testing.T) {
	l := NewLimiter(1, 1)
	for i := 0; i < limiterMaxKeys+100; i++ {
		l.Allow(fmt.Sprintf("ip|10.0.%d.%d", i/256%256, i%256) + fmt.Sprint(i))
	}
	if got := l.buckets.Len(); got > limiterMaxKeys {
		t.Fatalf("bucket count %d exceeds the %d cap", got, limiterMaxKeys)
	}
}
