slug: config-key-resolvers-test
verdict: done
at:  9e68a01
ran: 2026-09-13T18:09:35Z on DietPi in 143s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/config/config.go [offset=350, limit=90]
✱ Glob "internal/config/*_test.go" 1 match
→ Read internal/config/config.go [limit=180]
→ Read internal/config/config_test.go

$ head -20 go.mod && go version
module github.com/astrate-platform/astrate

go 1.26.1

// Pinned dependency decisions — docs/ROADMAP.md §1.1 (frozen in docs/DESIGN.md).
// HTTP routing is stdlib net/http (Go >= 1.22 ServeMux patterns): deliberately no framework.
require (
	github.com/BurntSushi/toml v1.6.0 // TOML config (§5.1)
	github.com/coder/websocket v1.8.14 // live stream socket (M7)
	github.com/eclipse/paho.mqtt.golang v1.5.1 // test-only MQTT client (same family as official Go SDK)
	github.com/golang-jwt/jwt/v5 v5.3.1 // JWT verification, alg allowlist (§4.2)
	github.com/golang-migrate/migrate/v4 v4.19.1 // migrations via source/iofs + go:embed (§1.3)
	github.com/hashicorp/golang-lru/v2 v2.0.7 // verified-token LRU cache (§4.2)
	github.com/jackc/pgx/v5 v5.10.0 // Postgres/TimescaleDB driver + pgxpool (§1.3)
	github.com/mochi-mqtt/server/v2 v2.7.9 // embedded MQTT broker (§1.1)
	github.com/prometheus/client_golang v1.23.2 // metrics (§5.2)
	github.com/testcontainers/testcontainers-go v0.43.0 // T2 integration harness (§5.4)
	github.com/testcontainers/testcontainers-go/modules/postgres v0.42.0 // Timescale container module
	go.etcd.io/bbolt v1.4.3 // broker session store (§3.1; pinned over pebble)
	go.mongodb.org/mongo-driver/v2 v2.6.0 // BSON raw-document API only (§3.5.5)
go version go1.26.7 linux/arm64

The task: add container-free unit tests for `HousekeepingKeys()` (config.go:394) and `RealmJWTPublicKey()` (config.go:408). Both are untested file-reading helpers. I'll build `Config` structs directly and use `t.TempDir()` + `os.WriteFile` fixtures.
← Edit internal/config/config_test.go
Index: /root/astrate-mule/internal/config/config_test.go
===================================================================
--- /root/astrate-mule/internal/config/config_test.go
+++ /root/astrate-mule/internal/config/config_test.go
@@ -2,8 +2,9 @@
 
 import (
 	"os"
 	"path/filepath"
+	"slices"
 	"strings"
 	"testing"
 	"time"
 )

← Edit internal/config/config_test.go
Index: /root/astrate-mule/internal/config/config_test.go
===================================================================
--- /root/astrate-mule/internal/config/config_test.go
+++ /root/astrate-mule/internal/config/config_test.go
@@ -17,8 +17,109 @@
 	}
 	return p
 }
 
+func writePEMFile(t *testing.T, body string) string {
+	t.Helper()
+	p := filepath.Join(t.TempDir(), "key.pem")
+	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
+		t.Fatal(err)
+	}
+	return p
+}
+
+func TestHousekeepingKeys(t *testing.T) {
+	t.Run("inline plus file PEM concatenated and trimmed", func(t *testing.T) {
+		filePEM := "-----BEGIN PUBLIC KEY-----\nfile-key\n-----END PUBLIC KEY-----"
+		cfg := Config{
+			Housekeeping: HousekeepingConfig{
+				JWTPublicKeys:     []string{"inline-a", "inline-b"},
+				JWTPublicKeyFiles: []string{writePEMFile(t, "\n"+filePEM+"\n\n")},
+			},
+		}
+		keys, err := cfg.HousekeepingKeys()
+		if err != nil {
+			t.Fatalf("HousekeepingKeys: %v", err)
+		}
+		want := []string{"inline-a", "inline-b", filePEM}
+		if !slices.Equal(keys, want) {
+			t.Errorf("keys = %q, want %q", keys, want)
+		}
+	})
+	t.Run("inline keys only", func(t *testing.T) {
+		cfg := Config{
+			Housekeeping: HousekeepingConfig{JWTPublicKeys: []string{"a"}},
+		}
+		keys, err := cfg.HousekeepingKeys()
+		if err != nil {
+			t.Fatalf("HousekeepingKeys: %v", err)
+		}
+		if !slices.Equal(keys, []string{"a"}) {
+			t.Errorf("keys = %q, want [a]", keys)
+		}
+	})
+	t.Run("missing file error names the path", func(t *testing.T) {
+		missing := filepath.Join(t.TempDir(), "no-such.pem")
+		cfg := Config{
+			Housekeeping: HousekeepingConfig{JWTPublicKeyFiles: []string{missing}},
+		}
+		if _, err := cfg.HousekeepingKeys(); err == nil {
+			t.Fatal("expected error, got nil")
+		} else if !strings.Contains(err.Error(), missing) {
+			t.Errorf("error %q does not mention path %q", err.Error(), missing)
+		}
+	})
+}
+
+func TestRealmJWTPublicKey(t *testing.T) {
+	t.Run("file PEM read and trimmed", func(t *testing.T) {
+		filePEM := "-----BEGIN PUBLIC KEY-----\nrealm-key\n-----END PUBLIC KEY-----"
+		cfg := Config{
+			Realm: RealmConfig{JWTPublicKeyFile: writePEMFile(t, "\n"+filePEM+"\n\n")},
+		}
+		got, err := cfg.RealmJWTPublicKey()
+		if err != nil {
+			t.Fatalf("RealmJWTPublicKey: %v", err)
+		}
+		if got != filePEM {
+			t.Errorf("key = %q, want %q", got, filePEM)
+		}
+	})
+	t.Run("inline wins over file, not trimmed", func(t *testing.T) {
+		cfg := Config{
+			Realm: RealmConfig{
+				JWTPublicKey:     "\n inline-pem \n",
+				JWTPublicKeyFile: filepath.Join(t.TempDir(), "no-such.pem"),
+			},
+		}
+		got, err := cfg.RealmJWTPublicKey()
+		if err != nil {
+			t.Fatalf("RealmJWTPublicKey: %v", err)
+		}
+		if got != "\n inline-pem \n" {
+			t.Errorf("key = %q, want inline block verbatim", got)
+		}
+	})
+	t.Run("neither set returns empty", func(t *testing.T) {
+		got, err := (Config{}).RealmJWTPublicKey()
+		if err != nil {
+			t.Fatalf("RealmJWTPublicKey: %v", err)
+		}
+		if got != "" {
+			t.Errorf("key = %q, want empty", got)
+		}
+	})
+	t.Run("missing file error names the path", func(t *testing.T) {
+		missing := filepath.Join(t.TempDir(), "no-such.pem")
+		cfg := Config{Realm: RealmConfig{JWTPublicKeyFile: missing}}
+		if _, err := cfg.RealmJWTPublicKey(); err == nil {
+			t.Fatal("expected error, got nil")
+		} else if !strings.Contains(err.Error(), missing) {
+			t.Errorf("error %q does not mention path %q", err.Error(), missing)
+		}
+	})
+}
+
 func TestDefaultsAndPrecedence(t *testing.T) {
 	// Defaults fill everything but the required DSN; the TOML overrides a
 	// subset; the environment overrides the TOML.
 	path := writeTOML(t, `

Now run the gate:
$ go test ./internal/config/ && go vet ./internal/config/ && gofmt -l internal/config/
# github.com/astrate-platform/astrate/internal/config [github.com/astrate-platform/astrate/internal/config.test]
internal/config/config_test.go:103:26: cannot call pointer method RealmJWTPublicKey on Config
FAIL	github.com/astrate-platform/astrate/internal/config [build failed]
FAIL

Pointer receiver — need an addressable Config:
← Edit internal/config/config_test.go
Index: /root/astrate-mule/internal/config/config_test.go
===================================================================
--- /root/astrate-mule/internal/config/config_test.go
+++ /root/astrate-mule/internal/config/config_test.go
@@ -99,9 +99,10 @@
 		t.Errorf("key = %q, want inline block verbatim", got)
 	}
 })
 t.Run("neither set returns empty", func(t *testing.T) {
-	got, err := (Config{}).RealmJWTPublicKey()
+	var cfg Config
+	got, err := cfg.RealmJWTPublicKey()
 	if err != nil {
 		t.Fatalf("RealmJWTPublicKey: %v", err)
 	}
 	if got != "" {

$ go test ./internal/config/ && go vet ./internal/config/ && gofmt -l internal/config/
ok  	github.com/astrate-platform/astrate/internal/config	0.079s

$ go test -race ./internal/config/
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/config	0.027s
FAIL

Gate is green; `-race` is unavailable on this box's 39-bit VMA kernel (same documented constraint as the Pi) — irrelevant here since the tests touch no concurrency.

```
Done: added container-free unit tests for HousekeepingKeys() and RealmJWTPublicKey() covering happy path, missing-file error naming the path, inline-over-file priority, and TrimSpace
Files: internal/config/config_test.go
Verified: go test ./internal/config/ -> pass; go vet ./internal/config/ -> pass; gofmt -l internal/config/ -> clean; go test -race ./internal/config/ -> fail (TSan unsupported VMA range, platform, not code)
Unsure: nothing
Follow-ups: none
```
