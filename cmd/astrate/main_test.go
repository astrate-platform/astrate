//go:build integration

package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"

	"github.com/astrate-platform/astrate/internal/config"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

const composeDSN = "postgres://astrate:astrate@127.0.0.1:5432/astrate?sslmode=disable"

// dialDSN finds a reachable database: ASTRATE_TEST_DSN first, then compose.
func dialDSN(t *testing.T) string {
	t.Helper()
	for _, dsn := range []string{envOr("ASTRATE_TEST_DSN", ""), composeDSN} {
		if dsn == "" {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		conn, err := pgx.Connect(ctx, dsn)
		cancel()
		if err == nil {
			_ = conn.Close(context.Background())
			return dsn
		}
	}
	t.Skip("the boot suite needs a database: run `make up` or set ASTRATE_TEST_DSN")
	return ""
}

func envOr(k, def string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return def
}

// freeAddr returns a currently-free loopback address (small TOCTOU window,
// acceptable in a test).
func freeAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	_ = l.Close()
	return addr
}

// bootConfig builds a dev-mode config booting against dsn, auto-provisioning a
// realm whose JWT public key is pub.
func bootConfig(t *testing.T, dsn, httpAddr, sessPath, realmName, pub string) config.Config {
	t.Helper()
	cfg := config.Default()
	cfg.Database.DSN = dsn
	cfg.HTTP.Addr = httpAddr
	cfg.MQTT.Addr = "127.0.0.1:0"
	cfg.MQTT.DevAddr = "127.0.0.1:0"
	cfg.MQTT.InsecureDevMode = true
	cfg.MQTT.SessionStorePath = sessPath
	cfg.Engine.Shards = 2
	cfg.Realm.Name = realmName
	cfg.Realm.JWTPublicKey = pub
	cfg.Log.Level = "error"
	return cfg
}

func TestBoot(t *testing.T) {
	dsn := dialDSN(t)

	masterKey := make([]byte, 32)
	_, _ = rand.Read(masterKey)
	t.Setenv("ASTRATE_MASTER_KEY", hex.EncodeToString(masterKey))

	realmKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	pubDER, _ := x509.MarshalPKIXPublicKey(&realmKey.PublicKey)
	realmPub := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))

	var suffix [4]byte
	_, _ = rand.Read(suffix[:])
	realmName := "boot" + hex.EncodeToString(suffix[:])
	sessPath := t.TempDir() + "/sessions.db"
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	httpAddr := freeAddr(t)
	cfg := bootConfig(t, dsn, httpAddr, sessPath, realmName, realmPub)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- run(ctx, cfg, log) }()

	base := "http://" + httpAddr
	waitReady(t, base)

	t.Run("HealthAndMetrics", func(t *testing.T) {
		if code, _ := get(t, base+"/astrate/v1/health"); code != http.StatusOK {
			t.Errorf("health = %d, want 200", code)
		}
		code, body := get(t, base+"/astrate/v1/metrics")
		if code != http.StatusOK {
			t.Fatalf("metrics = %d, want 200", code)
		}
		for _, want := range []string{"astrate_broker_sessions", "astrate_db_pool_max_conns"} {
			if !strings.Contains(body, want) {
				t.Errorf("metrics missing %q", want)
			}
		}
	})

	t.Run("AutoProvisionedRealmServesPairing", func(t *testing.T) {
		id, _ := deviceid.Random()
		token := mintAgentToken(t, realmKey)
		secret := registerDevice(t, base, realmName, id.String(), token)
		if len(secret) != 44 {
			t.Errorf("credentials secret length = %d, want 44 (%q)", len(secret), secret)
		}
	})

	// Graceful shutdown: run returns cleanly within the drain budget.
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("run returned error on shutdown: %v", err)
		}
	case <-time.After(shutdownTimeout + 10*time.Second):
		t.Fatal("run did not shut down within the drain budget")
	}

	// Restart against the same DB + session store: the DB survives (the realm
	// auto-provision is a no-op) and the bbolt session store reopens.
	t.Run("RestartReusesState", func(t *testing.T) {
		cfg2 := cfg
		cfg2.HTTP.Addr = freeAddr(t)
		ctx2, cancel2 := context.WithCancel(context.Background())
		done2 := make(chan error, 1)
		go func() { done2 <- run(ctx2, cfg2, log) }()
		waitReady(t, "http://"+cfg2.HTTP.Addr)
		cancel2()
		select {
		case err := <-done2:
			if err != nil {
				t.Fatalf("restart run returned error: %v", err)
			}
		case <-time.After(shutdownTimeout + 10*time.Second):
			t.Fatal("restart did not shut down in time")
		}
	})
}

// waitReady polls /readiness until it returns 200.
func waitReady(t *testing.T, base string) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		if code, _ := get(t, base+"/astrate/v1/readiness"); code == http.StatusOK {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatal("service never became ready")
}

func get(t *testing.T, url string) (int, string) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}

// mintAgentToken signs a realm JWT carrying the pairing-agent claim.
func mintAgentToken(t *testing.T, key *rsa.PrivateKey) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"a_pa": []string{".*::.*"},
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	s, err := tok.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// registerDevice runs Flow A against the live pairing API and returns the
// credentials secret.
func registerDevice(t *testing.T, base, realm, deviceID, token string) string {
	t.Helper()
	body := fmt.Sprintf(`{"data":{"hw_id":%q}}`, deviceID)
	req, _ := http.NewRequest(http.MethodPost,
		base+"/pairing/v1/"+realm+"/agent/devices", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("register request: %v", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("register status = %d, want 201 (%s)", resp.StatusCode, raw)
	}
	var env struct {
		Data struct {
			CredentialsSecret string `json:"credentials_secret"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("decode register response: %v (%s)", err, raw)
	}
	return env.Data.CredentialsSecret
}
