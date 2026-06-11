package auth

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// testKeys carries one RSA and one EC signing pair for the whole suite
// (generated once; RSA keygen is the slow part).
type testKeys struct {
	rsaKey *rsa.PrivateKey
	ecKey  *ecdsa.PrivateKey
}

var sharedKeys *testKeys

func keys(t *testing.T) *testKeys {
	t.Helper()
	if sharedKeys != nil {
		return sharedKeys
	}
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generating RSA key: %v", err)
	}
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating EC key: %v", err)
	}
	sharedKeys = &testKeys{rsaKey: rsaKey, ecKey: ecKey}
	return sharedKeys
}

// publicPEM encodes a public key as a PKIX PEM string.
func publicPEM(t *testing.T, pub crypto.PublicKey) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatalf("marshalling public key: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
}

// signToken builds and signs a token with the given claims map.
func signToken(t *testing.T, method jwt.SigningMethod, key any, claims jwt.MapClaims) string {
	t.Helper()
	s, err := jwt.NewWithClaims(method, claims).SignedString(key)
	if err != nil {
		t.Fatalf("signing token: %v", err)
	}
	return s
}

func TestVerifyValidTokens(t *testing.T) {
	tk := keys(t)
	now := time.Now

	cases := []struct {
		name   string
		method jwt.SigningMethod
		key    any
		pub    crypto.PublicKey
	}{
		{"RS256", jwt.SigningMethodRS256, tk.rsaKey, &tk.rsaKey.PublicKey},
		{"RS512", jwt.SigningMethodRS512, tk.rsaKey, &tk.rsaKey.PublicKey},
		{"ES256", jwt.SigningMethodES256, tk.ecKey, &tk.ecKey.PublicKey},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := signToken(t, tc.method, tc.key, jwt.MapClaims{
				"a_pa": []string{".*::.*"},
				"exp":  time.Now().Add(time.Hour).Unix(),
			})
			tok, err := Verify(s, []crypto.PublicKey{tc.pub}, now)
			if err != nil {
				t.Fatalf("Verify: %v", err)
			}
			if !tok.Authorizes(ClaimPairing, "POST", "agent/devices") {
				t.Error("catch-all a_pa claim should authorize POST agent/devices")
			}
			if exp, ok := tok.ExpiresAt(); !ok || !exp.After(time.Now()) {
				t.Errorf("ExpiresAt: got (%v, %v), want future expiry", exp, ok)
			}
		})
	}
}

func TestVerifyNoExpiry(t *testing.T) {
	tk := keys(t)
	// `exp` is optional and `iat` not required (docs/DESIGN.md §4.2 parity):
	// a claims-only token must verify.
	s := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{
		"a_aea": []string{".*::.*"},
	})
	tok, err := Verify(s, []crypto.PublicKey{&tk.rsaKey.PublicKey}, time.Now)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if _, ok := tok.ExpiresAt(); ok {
		t.Error("ExpiresAt: want none for exp-less token")
	}
}

func TestVerifyRejections(t *testing.T) {
	tk := keys(t)
	rsaPub := []crypto.PublicKey{&tk.rsaKey.PublicKey}

	t.Run("Expired", func(t *testing.T) {
		s := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{
			"exp": time.Now().Add(-time.Hour).Unix(),
		})
		_, err := Verify(s, rsaPub, time.Now)
		if !errors.Is(err, jwt.ErrTokenExpired) {
			t.Errorf("expired token: got %v, want jwt.ErrTokenExpired", err)
		}
	})

	t.Run("NotYetValid", func(t *testing.T) {
		s := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{
			"nbf": time.Now().Add(time.Hour).Unix(),
			"exp": time.Now().Add(2 * time.Hour).Unix(),
		})
		_, err := Verify(s, rsaPub, time.Now)
		if !errors.Is(err, jwt.ErrTokenNotValidYet) {
			t.Errorf("nbf-future token: got %v, want jwt.ErrTokenNotValidYet", err)
		}
	})

	t.Run("Garbage", func(t *testing.T) {
		for _, s := range []string{"", "garbage", "a.b", "a.b.c", "Bearer x"} {
			if _, err := Verify(s, rsaPub, time.Now); err == nil {
				t.Errorf("Verify(%q): want error", s)
			}
		}
	})

	t.Run("HMACConfusionAttack", func(t *testing.T) {
		// Classic alg-confusion: sign an HS256 token using the realm's
		// *public* key bytes as the HMAC secret. A verifier that feeds key
		// material into HMAC verification would accept it.
		pubBytes := []byte(publicPEM(t, &tk.rsaKey.PublicKey))
		s := signToken(t, jwt.SigningMethodHS256, pubBytes, jwt.MapClaims{
			"a_pa": []string{".*::.*"},
			"exp":  time.Now().Add(time.Hour).Unix(),
		})
		if _, err := Verify(s, rsaPub, time.Now); err == nil {
			t.Fatal("HS256 token signed with public key bytes must be rejected")
		}
	})

	t.Run("AlgNone", func(t *testing.T) {
		s := signToken(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, jwt.MapClaims{
			"a_pa": []string{".*::.*"},
		})
		if _, err := Verify(s, rsaPub, time.Now); err == nil {
			t.Fatal("alg=none token must be rejected")
		}
	})

	t.Run("NoKeys", func(t *testing.T) {
		s := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{})
		if _, err := Verify(s, nil, time.Now); !errors.Is(err, ErrNoRealmKeys) {
			t.Errorf("empty key set: got %v, want ErrNoRealmKeys", err)
		}
	})
}

func TestVerifyKeyRotation(t *testing.T) {
	tk := keys(t)
	otherRSA, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	s := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{
		"a_rma": []string{".*::.*"},
	})

	// The signing key sits last in a 3-key set (EC + stranger RSA first):
	// rotation means any key in the set may verify the token.
	set := []crypto.PublicKey{&tk.ecKey.PublicKey, &otherRSA.PublicKey, &tk.rsaKey.PublicKey}
	if _, err := Verify(s, set, time.Now); err != nil {
		t.Fatalf("Verify with rotated key set: %v", err)
	}

	// Without the signing key the token must fail with ErrNoKeyMatched.
	set = []crypto.PublicKey{&tk.ecKey.PublicKey, &otherRSA.PublicKey}
	if _, err := Verify(s, set, time.Now); !errors.Is(err, ErrNoKeyMatched) {
		t.Errorf("Verify without signing key: got %v, want ErrNoKeyMatched", err)
	}
}

func TestParsePublicKeysPEM(t *testing.T) {
	tk := keys(t)

	t.Run("PKIXAndPKCS1", func(t *testing.T) {
		pkcs1 := string(pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&tk.rsaKey.PublicKey),
		}))
		got, err := ParsePublicKeysPEM([]string{
			publicPEM(t, &tk.rsaKey.PublicKey),
			publicPEM(t, &tk.ecKey.PublicKey),
			pkcs1,
		})
		if err != nil {
			t.Fatalf("ParsePublicKeysPEM: %v", err)
		}
		if len(got) != 3 {
			t.Fatalf("got %d keys, want 3", len(got))
		}
	})

	t.Run("MultipleBlocksInOneEntry", func(t *testing.T) {
		combined := publicPEM(t, &tk.rsaKey.PublicKey) + publicPEM(t, &tk.ecKey.PublicKey)
		got, err := ParsePublicKeysPEM([]string{combined})
		if err != nil {
			t.Fatalf("ParsePublicKeysPEM: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("got %d keys, want 2", len(got))
		}
	})

	t.Run("Rejections", func(t *testing.T) {
		bad := [][]string{
			{"not a pem"},
			{""},
			{"-----BEGIN PUBLIC KEY-----\nZ2FyYmFnZQ==\n-----END PUBLIC KEY-----\n"},
			{string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte{1, 2, 3}}))},
		}
		for _, pems := range bad {
			if _, err := ParsePublicKeysPEM(pems); err == nil {
				t.Errorf("ParsePublicKeysPEM(%q): want error", pems)
			}
		}
	})
}
