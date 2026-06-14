package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTOML(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "astrate.toml")
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestDefaultsAndPrecedence(t *testing.T) {
	// Defaults fill everything but the required DSN; the TOML overrides a
	// subset; the environment overrides the TOML.
	path := writeTOML(t, `
[database]
dsn = "postgres://toml/db"

[http]
addr = ":9090"

[engine]
shards = 8
batch_max_wait = "10ms"

[mqtt]
insecure_dev_mode = true
`)
	t.Setenv("ASTRATE_DATABASE_DSN", "postgres://env/db")
	t.Setenv("ASTRATE_ENGINE_SHARDS", "32")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	// env beats TOML.
	if cfg.Database.DSN != "postgres://env/db" {
		t.Errorf("dsn = %q, want env override", cfg.Database.DSN)
	}
	if cfg.Engine.Shards != 32 {
		t.Errorf("shards = %d, want 32 (env)", cfg.Engine.Shards)
	}
	// TOML beats default.
	if cfg.HTTP.Addr != ":9090" {
		t.Errorf("http.addr = %q, want :9090 (toml)", cfg.HTTP.Addr)
	}
	if cfg.Engine.BatchMaxWait.Std() != 10*time.Millisecond {
		t.Errorf("batch_max_wait = %v, want 10ms (toml)", cfg.Engine.BatchMaxWait.Std())
	}
	// default untouched.
	if cfg.MQTT.Addr != ":8883" {
		t.Errorf("mqtt.addr = %q, want default :8883", cfg.MQTT.Addr)
	}
	if cfg.Engine.ShardQueue != 4096 || cfg.Pairing.BcryptCost != 10 {
		t.Errorf("defaults not applied: %+v", cfg.Engine)
	}
}

func TestValidation(t *testing.T) {
	cases := map[string]struct {
		body    string
		wantErr bool
	}{
		"missing dsn":              {`[http]` + "\naddr = \":1\"\n[mqtt]\ninsecure_dev_mode=true", true},
		"tls required outside dev": {`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=false", true},
		"dev mode ok no tls":       {`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true", false},
		"http tls half set":        {`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[http]\ntls_cert_file=\"c\"", true},
		"realm without key":        {`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[realm]\nname=\"r\"", true},
		"bad log level":            {`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[log]\nlevel=\"trace\"", true},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Load(writeTOML(t, tc.body))
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestExampleParses(t *testing.T) {
	// The shipped example must load (the env reset keeps a developer's own
	// ASTRATE_* vars from leaking in).
	for _, e := range os.Environ() {
		if k, _, ok := cut(e, '='); ok && len(k) > 8 && k[:8] == "ASTRATE_" {
			t.Setenv(k, "")
			_ = os.Unsetenv(k)
		}
	}
	cfg, err := Load("config.example.toml")
	if err != nil {
		t.Fatalf("config.example.toml does not load: %v", err)
	}
	if cfg.Database.DSN == "" {
		t.Error("example should set a database dsn")
	}
}

func cut(s string, sep byte) (string, string, bool) {
	for i := range len(s) {
		if s[i] == sep {
			return s[:i], s[i+1:], true
		}
	}
	return s, "", false
}
