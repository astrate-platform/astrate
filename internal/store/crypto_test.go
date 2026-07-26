package store

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func testMasterKey(t *testing.T) []byte {
	t.Helper()
	key := make([]byte, MasterKeySize)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("generating master key: %v", err)
	}
	return key
}

func TestKeySealerRoundTrip(t *testing.T) {
	ks, err := NewKeySealer(testMasterKey(t))
	if err != nil {
		t.Fatalf("NewKeySealer: %v", err)
	}

	plaintext := []byte("-----BEGIN EC PRIVATE KEY-----\nfake\n-----END EC PRIVATE KEY-----\n")
	sealed, err := ks.Seal(plaintext)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	if bytes.Contains(sealed, plaintext) {
		t.Fatal("sealed box contains the plaintext")
	}

	got, err := ks.Open(sealed)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Fatalf("Open round-trip mismatch: got %q want %q", got, plaintext)
	}
}

func TestKeySealerNonceUniqueness(t *testing.T) {
	ks, err := NewKeySealer(testMasterKey(t))
	if err != nil {
		t.Fatalf("NewKeySealer: %v", err)
	}
	plaintext := []byte("same plaintext")

	a, err := ks.Seal(plaintext)
	if err != nil {
		t.Fatalf("Seal a: %v", err)
	}
	b, err := ks.Seal(plaintext)
	if err != nil {
		t.Fatalf("Seal b: %v", err)
	}
	if bytes.Equal(a, b) {
		t.Fatal("two seals of the same plaintext produced identical boxes (nonce reuse)")
	}
	if bytes.Equal(a[:12], b[:12]) {
		t.Fatal("two seals share the same nonce")
	}
}

func TestKeySealerWrongKeyFails(t *testing.T) {
	ks1, err := NewKeySealer(testMasterKey(t))
	if err != nil {
		t.Fatalf("NewKeySealer 1: %v", err)
	}
	ks2, err := NewKeySealer(testMasterKey(t))
	if err != nil {
		t.Fatalf("NewKeySealer 2: %v", err)
	}

	sealed, err := ks1.Seal([]byte("secret"))
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	if _, err := ks2.Open(sealed); err == nil {
		t.Fatal("Open with a different master key succeeded")
	}
}

func TestKeySealerTamperDetected(t *testing.T) {
	ks, err := NewKeySealer(testMasterKey(t))
	if err != nil {
		t.Fatalf("NewKeySealer: %v", err)
	}
	sealed, err := ks.Seal([]byte("secret"))
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	sealed[len(sealed)-1] ^= 0x01
	if _, err := ks.Open(sealed); err == nil {
		t.Fatal("Open of a tampered box succeeded")
	}
}

func TestKeySealerShortInput(t *testing.T) {
	ks, err := NewKeySealer(testMasterKey(t))
	if err != nil {
		t.Fatalf("NewKeySealer: %v", err)
	}
	for _, n := range []int{0, 1, 12, 27} {
		if _, err := ks.Open(make([]byte, n)); err == nil {
			t.Errorf("Open of %d-byte box succeeded", n)
		}
	}
}

func TestNewKeySealerRejectsBadLength(t *testing.T) {
	for _, n := range []int{0, 16, 31, 33, 64} {
		if _, err := NewKeySealer(make([]byte, n)); err == nil {
			t.Errorf("NewKeySealer accepted a %d-byte key", n)
		}
	}
}

func TestLoadMasterKeyFromEnvHex(t *testing.T) {
	key := testMasterKey(t)
	t.Setenv(EnvMasterKey, hex.EncodeToString(key))
	got, err := LoadMasterKey()
	if err != nil {
		t.Fatalf("LoadMasterKey: %v", err)
	}
	if !bytes.Equal(got, key) {
		t.Fatal("hex-loaded key mismatch")
	}
}

func TestLoadMasterKeyFromEnvBase64(t *testing.T) {
	key := testMasterKey(t)
	for name, enc := range map[string]*base64.Encoding{
		"std": base64.StdEncoding,
		"raw": base64.RawStdEncoding,
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv(EnvMasterKey, enc.EncodeToString(key))
			got, err := LoadMasterKey()
			if err != nil {
				t.Fatalf("LoadMasterKey: %v", err)
			}
			if !bytes.Equal(got, key) {
				t.Fatal("base64-loaded key mismatch")
			}
		})
	}
}

func TestLoadMasterKeyFromFile(t *testing.T) {
	key := testMasterKey(t)

	t.Run("raw", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "master.key")
		if err := os.WriteFile(path, key, 0o600); err != nil {
			t.Fatal(err)
		}
		t.Setenv(EnvMasterKey, "")
		t.Setenv(EnvMasterKeyFile, path)
		got, err := LoadMasterKey()
		if err != nil {
			t.Fatalf("LoadMasterKey: %v", err)
		}
		if !bytes.Equal(got, key) {
			t.Fatal("file-loaded raw key mismatch")
		}
	})

	t.Run("hex-with-newline", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "master.key")
		if err := os.WriteFile(path, []byte(hex.EncodeToString(key)+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		t.Setenv(EnvMasterKey, "")
		t.Setenv(EnvMasterKeyFile, path)
		got, err := LoadMasterKey()
		if err != nil {
			t.Fatalf("LoadMasterKey: %v", err)
		}
		if !bytes.Equal(got, key) {
			t.Fatal("file-loaded hex key mismatch")
		}
	})
}

func TestLoadMasterKeyErrors(t *testing.T) {
	t.Setenv(EnvMasterKey, "")
	t.Setenv(EnvMasterKeyFile, "")
	if _, err := LoadMasterKey(); !errors.Is(err, ErrNoMasterKey) {
		t.Fatalf("unset env: got %v, want ErrNoMasterKey", err)
	}

	t.Setenv(EnvMasterKey, "not-a-key")
	if _, err := LoadMasterKey(); err == nil {
		t.Fatal("garbage env value accepted")
	}

	t.Setenv(EnvMasterKey, base64.StdEncoding.EncodeToString(make([]byte, 16)))
	if _, err := LoadMasterKey(); err == nil {
		t.Fatal("16-byte base64 key accepted")
	}
}
