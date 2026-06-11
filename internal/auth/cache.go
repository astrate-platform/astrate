package auth

import (
	"crypto"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	lru "github.com/hashicorp/golang-lru/v2"
)

// DefaultCacheSize is the verified-token LRU capacity (docs/DESIGN.md §4.2).
const DefaultCacheSize = 1024

// Cache memoizes verified tokens so the signature check, claim validation,
// and regex compilation run once per (token, key set) instead of once per
// request. Entries are keyed by SHA-256 of the token *and* of the key set,
// so rotating a realm's keys naturally invalidates its cached tokens.
type Cache struct {
	lru *lru.Cache[[sha256.Size]byte, *Token]
	now func() time.Time
	// verify is the verification function; tests swap in a spy.
	verify func(string, []crypto.PublicKey, func() time.Time) (*Token, error)
}

// NewCache builds a Cache with the given capacity (values < 1 fall back to
// DefaultCacheSize).
func NewCache(size int) *Cache {
	if size < 1 {
		size = DefaultCacheSize
	}
	// lru.New only fails for size < 1, which is excluded above.
	cache, err := lru.New[[sha256.Size]byte, *Token](size)
	if err != nil {
		panic(fmt.Sprintf("auth: building %d-entry LRU: %v", size, err))
	}
	return &Cache{lru: cache, now: time.Now, verify: Verify}
}

// cacheKey fingerprints the token together with the PEM key set. Domain
// separation between token and keys (and between keys) uses a 0x00 byte,
// which cannot appear in base64url tokens or PEM text.
func cacheKey(tokenString string, keysPEM []string) [sha256.Size]byte {
	h := sha256.New()
	h.Write([]byte(tokenString))
	for _, k := range keysPEM {
		h.Write([]byte{0})
		h.Write([]byte(k))
	}
	var key [sha256.Size]byte
	h.Sum(key[:0])
	return key
}

// Verify returns the verified Token for tokenString against the realm's PEM
// key set, from cache when possible. Cache hits still honour `exp`: a token
// that expired while cached is evicted and rejected. Failed verifications
// are not cached.
func (c *Cache) Verify(tokenString string, keysPEM []string) (*Token, error) {
	key := cacheKey(tokenString, keysPEM)

	if tok, ok := c.lru.Get(key); ok {
		if exp, has := tok.ExpiresAt(); has && !c.now().Before(exp) {
			c.lru.Remove(key)
			return nil, fmt.Errorf("auth: cached token: %w", jwt.ErrTokenExpired)
		}
		return tok, nil
	}

	keys, err := ParsePublicKeysPEM(keysPEM)
	if err != nil {
		return nil, err
	}
	tok, err := c.verify(tokenString, keys, c.now)
	if err != nil {
		return nil, err
	}
	c.lru.Add(key, tok)
	return tok, nil
}
