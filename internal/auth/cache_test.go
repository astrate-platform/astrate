package auth

import (
	"crypto"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestCacheHitAvoidsReverification(t *testing.T) {
	tk := keys(t)
	pubPEM := publicPEM(t, &tk.rsaKey.PublicKey)
	token := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{
		"a_pa": []string{".*::.*"},
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	c := NewCache(8)
	calls := 0
	c.verify = func(s string, keys []crypto.PublicKey, now func() time.Time) (*Token, error) {
		calls++
		return Verify(s, keys, now)
	}

	first, err := c.Verify(token, []string{pubPEM})
	if err != nil {
		t.Fatalf("first Verify: %v", err)
	}
	second, err := c.Verify(token, []string{pubPEM})
	if err != nil {
		t.Fatalf("second Verify: %v", err)
	}
	if calls != 1 {
		t.Errorf("verification calls: got %d, want 1 (second lookup must hit the cache)", calls)
	}
	if first != second {
		t.Error("cache hit should return the same *Token instance")
	}
}

func TestCacheDistinctTokensDoNotCollide(t *testing.T) {
	tk := keys(t)
	pubPEM := publicPEM(t, &tk.rsaKey.PublicKey)

	tokenA := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{
		"a_pa": []string{".*::.*"},
	})
	tokenB := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{
		"a_aea": []string{".*::.*"},
	})

	c := NewCache(8)
	tokA, err := c.Verify(tokenA, []string{pubPEM})
	if err != nil {
		t.Fatal(err)
	}
	tokB, err := c.Verify(tokenB, []string{pubPEM})
	if err != nil {
		t.Fatal(err)
	}

	if !tokA.Authorizes(ClaimPairing, "POST", "x") || tokA.Authorizes(ClaimAppEngine, "GET", "x") {
		t.Error("token A grants mixed up")
	}
	if !tokB.Authorizes(ClaimAppEngine, "GET", "x") || tokB.Authorizes(ClaimPairing, "POST", "x") {
		t.Error("token B grants mixed up")
	}
}

func TestCacheKeyRotationInvalidates(t *testing.T) {
	tk := keys(t)
	pemA := publicPEM(t, &tk.rsaKey.PublicKey)
	pemB := publicPEM(t, &tk.ecKey.PublicKey)
	token := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{
		"a_pa": []string{".*::.*"},
	})

	c := NewCache(8)
	calls := 0
	c.verify = func(s string, keys []crypto.PublicKey, now func() time.Time) (*Token, error) {
		calls++
		return Verify(s, keys, now)
	}

	if _, err := c.Verify(token, []string{pemA}); err != nil {
		t.Fatal(err)
	}
	// Same token, different key set: must re-verify (and still succeed,
	// because the signing key is in the rotated set).
	if _, err := c.Verify(token, []string{pemB, pemA}); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Errorf("verification calls: got %d, want 2 (key set change must miss the cache)", calls)
	}

	// A rotated-away key set rejects even though the token was cached
	// under the old set.
	if _, err := c.Verify(token, []string{pemB}); !errors.Is(err, ErrNoKeyMatched) {
		t.Errorf("verify under rotated-away keys: got %v, want ErrNoKeyMatched", err)
	}
}

func TestCacheHitHonoursExpiry(t *testing.T) {
	tk := keys(t)
	pubPEM := publicPEM(t, &tk.rsaKey.PublicKey)
	token := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{
		"a_pa": []string{".*::.*"},
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	c := NewCache(8)
	clock := time.Now()
	c.now = func() time.Time { return clock }

	if _, err := c.Verify(token, []string{pubPEM}); err != nil {
		t.Fatalf("initial Verify: %v", err)
	}

	// Two hours later the cached entry must be evicted and rejected.
	clock = clock.Add(2 * time.Hour)
	if _, err := c.Verify(token, []string{pubPEM}); !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("expired-while-cached: got %v, want jwt.ErrTokenExpired", err)
	}
	if c.lru.Len() != 0 {
		t.Errorf("expired entry not evicted: %d entries left", c.lru.Len())
	}
}

func TestCacheFailedVerificationNotCached(t *testing.T) {
	tk := keys(t)
	pubPEM := publicPEM(t, &tk.ecKey.PublicKey) // wrong key on purpose
	token := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{})

	c := NewCache(8)
	if _, err := c.Verify(token, []string{pubPEM}); err == nil {
		t.Fatal("verification against the wrong key must fail")
	}
	if c.lru.Len() != 0 {
		t.Errorf("failed verification must not be cached, got %d entries", c.lru.Len())
	}
}
