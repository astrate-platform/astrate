package config

import (
	"os"
	"path/filepath"
	"strings"
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
	t.Setenv("ASTRATE_HTTP_CORS_ALLOWED_ORIGINS", "http://a.example, http://b.example:4040")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(cfg.HTTP.CORSAllowedOrigins) != 2 ||
		cfg.HTTP.CORSAllowedOrigins[0] != "http://a.example" ||
		cfg.HTTP.CORSAllowedOrigins[1] != "http://b.example:4040" {
		t.Errorf("cors origins = %v, want the two trimmed env entries", cfg.HTTP.CORSAllowedOrigins)
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
		"cors wildcard ok":         {`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[http]\naddr=\":1\"\ncors_allowed_origins=[\"*\"]", false},
		"cors origin ok":           {`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[http]\naddr=\":1\"\ncors_allowed_origins=[\"http://localhost:4040\"]", false},
		"cors origin with path":    {`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[http]\naddr=\":1\"\ncors_allowed_origins=[\"http://h/app\"]", true},
		"cors origin bare host":    {`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[http]\naddr=\":1\"\ncors_allowed_origins=[\"localhost:4040\"]", true},
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

func TestTriggersZeroValue(t *testing.T) {
	body := `
[database]
dsn = "x"
[mqtt]
insecure_dev_mode = true
`
	cfg, err := Load(writeTOML(t, body))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Triggers.Forward.Kind != "" {
		t.Errorf("Triggers.Forward.Kind = %q, want empty", cfg.Triggers.Forward.Kind)
	}
}

func TestTriggersForwardDisabled(t *testing.T) {
	body := `
[database]
dsn = "x"
[mqtt]
insecure_dev_mode = true
[triggers.forward]
url = "nope"
`
	cfg, err := Load(writeTOML(t, body))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Triggers.Forward.Kind != "" {
		t.Errorf("Kind = %q, want empty", cfg.Triggers.Forward.Kind)
	}
}

func TestTriggersForwardHTTPRoundTrip(t *testing.T) {
	body := `
[database]
dsn = "x"
[mqtt]
insecure_dev_mode = true
[triggers.forward]
kind = "http"
url = "https://bus.example/trigger"
method = "POST"
static_headers = { Authorization = "Bearer tok", X-Forward = "yes" }
`
	cfg, err := Load(writeTOML(t, body))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	f := cfg.Triggers.Forward
	if f.Kind != "http" {
		t.Errorf("Kind = %q, want http", f.Kind)
	}
	if f.URL != "https://bus.example/trigger" {
		t.Errorf("URL = %q", f.URL)
	}
	if f.Method != "POST" {
		t.Errorf("Method = %q, want POST", f.Method)
	}
	if f.StaticHeaders["Authorization"] != "Bearer tok" || f.StaticHeaders["X-Forward"] != "yes" {
		t.Errorf("StaticHeaders = %v", f.StaticHeaders)
	}
}

func TestTriggersForwardNATSRoundTrip(t *testing.T) {
	body := `
[database]
dsn = "x"
[mqtt]
insecure_dev_mode = true
[triggers.forward]
kind = "nats"
url = "nats://bus.example:4222"
subject = "astrate.triggers"
`
	cfg, err := Load(writeTOML(t, body))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	f := cfg.Triggers.Forward
	if f.Kind != "nats" {
		t.Errorf("Kind = %q, want nats", f.Kind)
	}
	if f.URL != "nats://bus.example:4222" {
		t.Errorf("URL = %q", f.URL)
	}
	if f.Subject != "astrate.triggers" {
		t.Errorf("Subject = %q, want astrate.triggers", f.Subject)
	}
}

func TestTriggersForwardValidation(t *testing.T) {
	cases := map[string]struct {
		body    string
		wantErr string
	}{
		"kind unknown": {
			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"amqp\"",
			"kind",
		},
		"nats missing url": {
			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"nats\"",
			"url is required",
		},
		"nats missing subject": {
			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"nats\"\nurl=\"nats://bus.example:4222\"",
			"subject is required",
		},
		// An absent url is refused by the absolute-URL rule as well, so
		// asserting only that the message mentions "url" leaves the dedicated
		// required-url check unbound — deleting it keeps this row green. The
		// wording is the rule: "is required" is the message an operator who
		// set no url at all should get, not "is not an absolute http or
		// https URL".
		"http missing url": {
			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"http\"",
			"url is required",
		},
		"http relative url": {
			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"http\"\nurl=\"nope\"",
			"url",
		},
		"http ftp scheme": {
			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"http\"\nurl=\"ftp://bus.example/x\"",
			"url",
		},
		"http no host": {
			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"http\"\nurl=\"http://\"",
			"url",
		},
		"http bad method": {
			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"http\"\nurl=\"https://bus.example/x\"\nmethod=\"FETCH\"",
			"method",
		},
		"http method post lowercase": {
			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"http\"\nurl=\"https://bus.example/x\"\nmethod=\"post\"",
			"",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Load(writeTOML(t, tc.body))
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error %q does not contain %q", err.Error(), tc.wantErr)
			}
		})
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

func TestHousekeepingDefaultRetentionEnv(t *testing.T) {
	const (
		prefixed = "ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION"
		bare     = "HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION"
	)
	body := `
[database]
dsn = "x"
[mqtt]
insecure_dev_mode = true
`
	load := func(t *testing.T) (Config, error) {
		t.Helper()
		return Load(writeTOML(t, body))
	}

	t.Run("absent leaves nil", func(t *testing.T) {
		t.Setenv(prefixed, "")
		t.Setenv(bare, "")
		cfg, err := load(t)
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if cfg.Housekeeping.DefaultDatastreamMaximumStorageRetention != nil {
			t.Errorf("default retention without env = %d, want nil",
				*cfg.Housekeeping.DefaultDatastreamMaximumStorageRetention)
		}
	})
	t.Run("prefixed name parses", func(t *testing.T) {
		t.Setenv(bare, "")
		t.Setenv(prefixed, "86400")
		cfg, err := load(t)
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if cfg.Housekeeping.DefaultDatastreamMaximumStorageRetention == nil ||
			*cfg.Housekeeping.DefaultDatastreamMaximumStorageRetention != 86400 {
			t.Errorf("default retention = %v, want 86400", cfg.Housekeeping.DefaultDatastreamMaximumStorageRetention)
		}
	})
	t.Run("bare upstream name wins", func(t *testing.T) {
		t.Setenv(prefixed, "111")
		t.Setenv(bare, "222")
		cfg, err := load(t)
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if cfg.Housekeeping.DefaultDatastreamMaximumStorageRetention == nil ||
			*cfg.Housekeeping.DefaultDatastreamMaximumStorageRetention != 222 {
			t.Errorf("default retention = %v, want 222 (bare name)", cfg.Housekeeping.DefaultDatastreamMaximumStorageRetention)
		}
	})
	t.Run("malformed fails loud", func(t *testing.T) {
		t.Setenv(prefixed, "")
		t.Setenv(bare, "soon")
		if _, err := load(t); err == nil {
			t.Error("malformed retention env: got nil error")
		}
	})
	t.Run("negative fails loud", func(t *testing.T) {
		t.Setenv(prefixed, "")
		t.Setenv(bare, "-5")
		if _, err := load(t); err == nil {
			t.Error("negative retention env: got nil error")
		}
	})
}

func TestHousekeepingRealmDeletionDisabledEnv(t *testing.T) {
	const env = "ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"
	body := `
[database]
dsn = "x"
[mqtt]
insecure_dev_mode = true
`
	load := func(t *testing.T) (Config, error) {
		t.Helper()
		return Load(writeTOML(t, body))
	}

	t.Run("absent is false", func(t *testing.T) {
		cfg, err := load(t)
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if cfg.Housekeeping.RealmDeletionDisabled {
			t.Error("realm deletion disabled without env = true, want false")
		}
	})
	for _, truthy := range []string{"true", "TRUE", "True", "1"} {
		t.Run("true "+truthy, func(t *testing.T) {
			t.Setenv(env, truthy)
			cfg, err := load(t)
			if err != nil {
				t.Fatalf("Load: %v", err)
			}
			if !cfg.Housekeeping.RealmDeletionDisabled {
				t.Errorf("env %s=%s → false, want true", env, truthy)
			}
		})
	}
	for _, falsy := range []string{"0", "false"} {
		t.Run("false "+falsy, func(t *testing.T) {
			t.Setenv(env, falsy)
			cfg, err := load(t)
			if err != nil {
				t.Fatalf("Load: %v", err)
			}
			if cfg.Housekeeping.RealmDeletionDisabled {
				t.Errorf("env %s=%s → true, want false", env, falsy)
			}
		})
	}
	t.Run("banana fails loud", func(t *testing.T) {
		t.Setenv(env, "banana")
		if _, err := load(t); err == nil {
			t.Error("malformed boolean env: got nil error")
		}
	})
}
