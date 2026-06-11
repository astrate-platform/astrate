package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
)

// CA private keys are encrypted at rest (docs/DESIGN.md §2.2, §4.3): the
// realms.ca_private_key bytea column holds an AES-256-GCM sealed box produced
// by KeySealer under a 32-byte master key supplied via environment reference.
// Losing the master key means re-issuing realm CAs; devices re-pair
// automatically at their next credential rotation.

const (
	// EnvMasterKey names the env var holding the master key itself,
	// encoded as 64 hex characters or base64 (std or raw) of 32 bytes.
	EnvMasterKey = "ASTRATE_MASTER_KEY"
	// EnvMasterKeyFile names the env var holding a path to a file that
	// contains the master key (raw 32 bytes, or its hex/base64 text form).
	EnvMasterKeyFile = "ASTRATE_MASTER_KEY_FILE"
	// MasterKeySize is the required decoded master key length (AES-256).
	MasterKeySize = 32
)

// ErrNoMasterKey reports that neither master-key env reference is set.
var ErrNoMasterKey = fmt.Errorf("store: neither %s nor %s is set", EnvMasterKey, EnvMasterKeyFile)

// KeySealer seals and opens small secrets (realm CA private keys) with
// AES-256-GCM. The sealed form is nonce || ciphertext+tag with a fresh
// random 96-bit nonce per Seal.
type KeySealer struct {
	aead cipher.AEAD
}

// NewKeySealer builds a KeySealer from a raw 32-byte master key.
func NewKeySealer(masterKey []byte) (*KeySealer, error) {
	if len(masterKey) != MasterKeySize {
		return nil, fmt.Errorf("store: master key must be %d bytes, got %d", MasterKeySize, len(masterKey))
	}
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, fmt.Errorf("store: initializing AES: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("store: initializing GCM: %w", err)
	}
	return &KeySealer{aead: aead}, nil
}

// NewKeySealerFromEnv builds a KeySealer from the environment references
// (EnvMasterKey first, then EnvMasterKeyFile).
func NewKeySealerFromEnv() (*KeySealer, error) {
	key, err := LoadMasterKey()
	if err != nil {
		return nil, err
	}
	return NewKeySealer(key)
}

// LoadMasterKey resolves the master key from the environment: EnvMasterKey
// (hex or base64 text) wins over EnvMasterKeyFile (raw 32 bytes, or hex or
// base64 text, surrounding whitespace ignored).
func LoadMasterKey() ([]byte, error) {
	if v := os.Getenv(EnvMasterKey); v != "" {
		key, err := decodeMasterKey(v)
		if err != nil {
			return nil, fmt.Errorf("store: %s: %w", EnvMasterKey, err)
		}
		return key, nil
	}
	if path := os.Getenv(EnvMasterKeyFile); path != "" {
		raw, err := os.ReadFile(path) // #nosec G304 G703 -- operator-supplied key file reference is the feature
		if err != nil {
			return nil, fmt.Errorf("store: reading %s: %w", EnvMasterKeyFile, err)
		}
		if len(raw) == MasterKeySize {
			return raw, nil
		}
		key, err := decodeMasterKey(strings.TrimSpace(string(raw)))
		if err != nil {
			return nil, fmt.Errorf("store: %s (%s): %w", EnvMasterKeyFile, path, err)
		}
		return key, nil
	}
	return nil, ErrNoMasterKey
}

// decodeMasterKey decodes a textual master key: 64 hex characters, or the
// standard/raw base64 of 32 bytes.
func decodeMasterKey(s string) ([]byte, error) {
	if len(s) == 2*MasterKeySize {
		if key, err := hex.DecodeString(s); err == nil {
			return key, nil
		}
	}
	for _, enc := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding} {
		key, err := enc.DecodeString(s)
		if err == nil {
			if len(key) != MasterKeySize {
				return nil, fmt.Errorf("decoded master key is %d bytes, want %d", len(key), MasterKeySize)
			}
			return key, nil
		}
	}
	return nil, errors.New("master key is neither 64 hex characters nor base64 of 32 bytes")
}

// Seal encrypts plaintext, returning nonce || ciphertext+tag.
func (ks *KeySealer) Seal(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, ks.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("store: generating nonce: %w", err)
	}
	return ks.aead.Seal(nonce, nonce, plaintext, nil), nil
}

// Open decrypts a Seal-produced box, authenticating it in the process.
func (ks *KeySealer) Open(sealed []byte) ([]byte, error) {
	ns := ks.aead.NonceSize()
	if len(sealed) < ns+ks.aead.Overhead() {
		return nil, errors.New("store: sealed box too short")
	}
	plaintext, err := ks.aead.Open(nil, sealed[:ns], sealed[ns:], nil)
	if err != nil {
		return nil, fmt.Errorf("store: opening sealed box: %w", err)
	}
	return plaintext, nil
}
