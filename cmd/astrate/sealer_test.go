package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/config"
	"github.com/astrate-platform/astrate/internal/store"
)

// randomMasterKey returns a distinct 32-byte master key, so two keys from the
// same test are never equal and a sealer built from one is guaranteed not to
// open the other's boxes.
func randomMasterKey(t *testing.T) []byte {
	t.Helper()
	key := make([]byte, store.MasterKeySize)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	return key
}

// writeKeyFile drops raw 32-byte key material in a fresh t.TempDir() and
// returns its path.
func writeKeyFile(t *testing.T, key []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "master.key")
	if err := os.WriteFile(path, key, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// assertSealerHoldsKey pins *which* key a sealer was built from without
// reaching inside it: a box sealed by ks opens under want and must not open
// under other. Two candidate keys, one sealer, no false pass from a sealer
// that merely round-trips with itself.
func assertSealerHoldsKey(t *testing.T, ks *store.KeySealer, want, other []byte) {
	t.Helper()
	sealed, err := ks.Seal([]byte("realm-ca-private-key"))
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	for _, c := range []struct {
		role string
		key  []byte
		open bool
	}{
		{"configured key", want, true},
		{"other key", other, false},
	} {
		sealer, err := store.NewKeySealer(c.key)
		if err != nil {
			t.Fatalf("NewKeySealer(%s): %v", c.role, err)
		}
		_, err = sealer.Open(sealed)
		if c.open && err != nil {
			t.Errorf("box sealed by the loaded sealer did not open under the %s: %v", c.role, err)
		}
		if !c.open && err == nil {
			t.Errorf("box sealed by the loaded sealer opened under the %s", c.role)
		}
	}
}

// clearMasterKeyEnv unsets both master-key references for the rest of the
// test. t.Setenv is also what makes loadSealer's own os.Setenv safe here: the
// registered cleanup restores the previous value, so the key-file path the
// function exports does not leak into the process environment.
func clearMasterKeyEnv(t *testing.T) {
	t.Helper()
	t.Setenv(store.EnvMasterKey, "")
	t.Setenv(store.EnvMasterKeyFile, "")
}

// TestLoadSealerReadsConfiguredKeyFile pins the security.master_key_file ->
// store.EnvMasterKeyFile hand-off inside loadSealer. With both env references
// empty, that os.Setenv is the only thing that lets the configured file reach
// store.LoadMasterKey: drop it (or misspell the var) and the file is silently
// ignored and boot dies on the generic "master key" error even though the
// config named a valid file.
func TestLoadSealerReadsConfiguredKeyFile(t *testing.T) {
	fileKey := randomMasterKey(t)
	otherKey := randomMasterKey(t)
	keyPath := writeKeyFile(t, fileKey)
	clearMasterKeyEnv(t)

	cfg := config.Default()
	cfg.Security.MasterKeyFile = keyPath

	ks, err := loadSealer(cfg)
	if err != nil {
		t.Fatalf("loadSealer with security.master_key_file=%s: %v", keyPath, err)
	}
	if ks == nil {
		t.Fatal("loadSealer returned a nil sealer and a nil error")
	}
	assertSealerHoldsKey(t, ks, fileKey, otherKey)

	// The hand-off is the documented contract: the exported var is what the
	// store loader reads, so it must hold the configured path while the test
	// runs (t.Setenv restores it at cleanup).
	if got := os.Getenv(store.EnvMasterKeyFile); got != keyPath {
		t.Errorf("%s = %q, want %q", store.EnvMasterKeyFile, got, keyPath)
	}
}

// TestLoadSealerPrefersInlineEnvKey pins the store's precedence: a non-empty
// ASTRATE_MASTER_KEY wins over the file named in the config, so an operator's
// inline key is not silently overridden by a stale key_file.
func TestLoadSealerPrefersInlineEnvKey(t *testing.T) {
	inlineKey := randomMasterKey(t)
	fileKey := randomMasterKey(t)
	keyPath := writeKeyFile(t, fileKey)
	clearMasterKeyEnv(t)
	t.Setenv(store.EnvMasterKey, hex.EncodeToString(inlineKey))

	cfg := config.Default()
	cfg.Security.MasterKeyFile = keyPath

	ks, err := loadSealer(cfg)
	if err != nil {
		t.Fatalf("loadSealer with ASTRATE_MASTER_KEY set: %v", err)
	}
	assertSealerHoldsKey(t, ks, inlineKey, fileKey)
}

// TestLoadSealerMissingKeyNamesEveryReference pins the operator-facing error:
// when no key can be resolved, the message must name all three places one can
// be set, so an operator with a valid key_file does not conclude the key is
// broken.
func TestLoadSealerMissingKeyNamesEveryReference(t *testing.T) {
	for name, file := range map[string]string{
		"no reference at all":       "",
		"configured file is absent": filepath.Join(t.TempDir(), "never-written.key"),
	} {
		t.Run(name, func(t *testing.T) {
			clearMasterKeyEnv(t)
			cfg := config.Default()
			cfg.Security.MasterKeyFile = file

			ks, err := loadSealer(cfg)
			if err == nil {
				t.Fatalf("loadSealer with no usable key returned sealer %v and no error", ks)
			}
			if ks != nil {
				t.Errorf("loadSealer returned a sealer alongside error %v", err)
			}
			if file == "" && !errors.Is(err, store.ErrNoMasterKey) {
				t.Errorf("no-reference case: got %v, want it to wrap store.ErrNoMasterKey", err)
			}
			for _, want := range []string{
				store.EnvMasterKey,
				store.EnvMasterKeyFile,
				"security.master_key_file",
			} {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q does not name %s", err, want)
				}
			}
		})
	}
}
