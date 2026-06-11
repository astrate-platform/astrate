// Package auth implements the JWT validation and Astarte authorization-claim
// layer shared by every Astrate REST surface (docs/DESIGN.md §4.2). It
// reproduces upstream Astarte's token semantics exactly — asymmetric keys
// only (RSA/ECDSA, `none` and HMAC hard-rejected), per-realm multi-key
// rotation, and `"<verb-regex>:<opts>:<path-regex>"` authorization strings
// matched with implicit anchoring against the request method and the path
// relative to the realm base — so tokens minted by astartectl and existing
// operator tooling work unmodified.
//
// The package is pure (no database): key material is injected through the
// KeySource interface, which *store.Store already satisfies.
package auth

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// allowedSigningMethods is the docs/DESIGN.md §4.2 algorithm allowlist.
// Everything else — most notably `none` and the HS* HMAC family (the classic
// public-key-as-HMAC-secret confusion attack) — is rejected at parse time.
var allowedSigningMethods = []string{"RS256", "RS384", "RS512", "ES256", "ES384", "ES512"}

// Sentinel errors returned by token verification. All of them map to 401 at
// the HTTP layer; they are distinct so logs and tests can tell causes apart.
var (
	// ErrNoRealmKeys reports that the realm has no JWT public keys
	// configured, so no token can possibly verify.
	ErrNoRealmKeys = errors.New("auth: realm has no JWT public keys")
	// ErrNoKeyMatched reports that the token signature verified against
	// none of the realm's keys (wrong key, tampered token, or disallowed
	// algorithm).
	ErrNoKeyMatched = errors.New("auth: token matches none of the realm keys")
	// ErrUnsupportedKey reports PEM key material that is neither an RSA nor
	// an ECDSA public key.
	ErrUnsupportedKey = errors.New("auth: unsupported public key type")
)

// ParsePublicKeysPEM parses a realm's JWT public key set. Each entry may
// carry one or more PEM blocks; supported block types are PKIX "PUBLIC KEY"
// and PKCS#1 "RSA PUBLIC KEY", and the decoded keys must be RSA or ECDSA
// (matching the signing-method allowlist). An entry with no usable key is an
// error: silently dropping keys would turn a key-set typo into a hard 401
// for every token holder.
func ParsePublicKeysPEM(pems []string) ([]crypto.PublicKey, error) {
	var keys []crypto.PublicKey
	for i, entry := range pems {
		found := false
		rest := []byte(entry)
		for {
			var block *pem.Block
			block, rest = pem.Decode(rest)
			if block == nil {
				break
			}
			key, err := parsePublicKeyBlock(block)
			if err != nil {
				return nil, fmt.Errorf("auth: key entry %d: %w", i, err)
			}
			keys = append(keys, key)
			found = true
		}
		if !found {
			return nil, fmt.Errorf("auth: key entry %d: no PEM block found", i)
		}
	}
	return keys, nil
}

// parsePublicKeyBlock decodes a single PEM block into an RSA or ECDSA public key.
func parsePublicKeyBlock(block *pem.Block) (crypto.PublicKey, error) {
	switch block.Type {
	case "PUBLIC KEY":
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing PKIX public key: %w", err)
		}
		switch key.(type) {
		case *rsa.PublicKey, *ecdsa.PublicKey:
			return key, nil
		default:
			return nil, fmt.Errorf("%w: %T", ErrUnsupportedKey, key)
		}
	case "RSA PUBLIC KEY":
		key, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing PKCS#1 public key: %w", err)
		}
		return key, nil
	default:
		return nil, fmt.Errorf("%w: PEM block type %q", ErrUnsupportedKey, block.Type)
	}
}

// Verify parses tokenString, verifies its signature against the key set, and
// validates its registered claims (`exp` and `nbf` are honoured when
// present; `iat` is not required — upstream parity, docs/DESIGN.md §4.2).
// The token verifies if *any* key in the set matches, which is what makes
// zero-downtime key rotation work. now supplies the validation clock.
func Verify(tokenString string, keys []crypto.PublicKey, now func() time.Time) (*Token, error) {
	if len(keys) == 0 {
		return nil, ErrNoRealmKeys
	}
	if now == nil {
		now = time.Now
	}
	parser := jwt.NewParser(
		jwt.WithValidMethods(allowedSigningMethods),
		jwt.WithTimeFunc(now),
	)

	var lastErr error
	for _, key := range keys {
		claims := &astarteClaims{}
		_, err := parser.ParseWithClaims(tokenString, claims, func(*jwt.Token) (any, error) {
			return key, nil
		})
		if err == nil {
			return newToken(claims), nil
		}
		// A signature mismatch only means "not this key": keep rotating
		// through the set. Any other failure (malformed token, disallowed
		// alg surfaced as unverifiable, expired, not yet valid) is final.
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			lastErr = err
			continue
		}
		return nil, fmt.Errorf("auth: invalid token: %w", err)
	}
	return nil, fmt.Errorf("%w: %w", ErrNoKeyMatched, lastErr)
}
