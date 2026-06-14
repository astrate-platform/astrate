// Package config loads Astrate's single TOML configuration file with
// ASTRATE_* environment overrides (docs/DESIGN.md §5.1, ROADMAP §9 file 8.1).
// Precedence is default < TOML < environment. Zero-config defaults target the
// single-VPS case: only a database DSN and (outside dev mode) the broker's TLS
// identity must be supplied. The master encryption key is referenced here but
// read by internal/store (env or file), never inlined.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// Duration is a time.Duration that (un)marshals as a Go duration string
// ("50ms", "720h") in TOML — time.Duration itself is not a TextUnmarshaler.
type Duration time.Duration

// UnmarshalText parses a Go duration string.
func (d *Duration) UnmarshalText(b []byte) error {
	v, err := time.ParseDuration(string(b))
	if err != nil {
		return err
	}
	*d = Duration(v)
	return nil
}

// MarshalText renders the duration as a Go duration string.
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

// Std returns the underlying time.Duration.
func (d Duration) Std() time.Duration { return time.Duration(d) }

// Config is the whole Astrate configuration.
type Config struct {
	HTTP         HTTPConfig         `toml:"http"`
	MQTT         MQTTConfig         `toml:"mqtt"`
	Database     DatabaseConfig     `toml:"database"`
	Engine       EngineConfig       `toml:"engine"`
	Pairing      PairingConfig      `toml:"pairing"`
	Housekeeping HousekeepingConfig `toml:"housekeeping"`
	Storage      StorageConfig      `toml:"storage"`
	Security     SecurityConfig     `toml:"security"`
	Realm        RealmConfig        `toml:"realm"`
	Log          LogConfig          `toml:"log"`
}

// HTTPConfig is the single REST listener (§3.7). TLS is optional: leave the
// files empty to serve plaintext behind a TLS-terminating reverse proxy.
type HTTPConfig struct {
	Addr        string `toml:"addr"`
	TLSCertFile string `toml:"tls_cert_file"`
	TLSKeyFile  string `toml:"tls_key_file"`
}

// MQTTConfig is the embedded broker (§3.1).
type MQTTConfig struct {
	Addr             string `toml:"addr"`
	TLSCertFile      string `toml:"tls_cert_file"`
	TLSKeyFile       string `toml:"tls_key_file"`
	InsecureDevMode  bool   `toml:"insecure_dev_mode"`
	DevAddr          string `toml:"dev_addr"`
	SessionStorePath string `toml:"session_store_path"`
	// AdvertisedURL is the broker URL handed to devices by the pairing info
	// endpoint. Empty derives "mqtts://<Addr>" (fine for localhost; set it
	// explicitly when devices reach the broker by a different host).
	AdvertisedURL  string `toml:"advertised_url"`
	MaxPacketBytes uint32 `toml:"max_packet_bytes"`
}

// DatabaseConfig is the PostgreSQL/TimescaleDB connection.
type DatabaseConfig struct {
	DSN string `toml:"dsn"`
}

// EngineConfig tunes the ingestion pipeline (§1.4).
type EngineConfig struct {
	Shards          int      `toml:"shards"`
	ShardQueue      int      `toml:"shard_queue"`
	BatchMaxRows    int      `toml:"batch_max_rows"`
	BatchMaxWait    Duration `toml:"batch_max_wait"`
	MaxPayloadBytes int      `toml:"max_payload_bytes"`
}

// PairingConfig tunes credential issuance and its rate limits (§4.3, §4.5).
type PairingConfig struct {
	CertTTL           Duration `toml:"cert_ttl"`
	EnforceLatestCert bool     `toml:"enforce_latest_cert"`
	RegisterRate      float64  `toml:"register_rate"`
	RegisterBurst     int      `toml:"register_burst"`
	CredentialsRate   float64  `toml:"credentials_rate"`
	CredentialsBurst  int      `toml:"credentials_burst"`
	BcryptCost        int      `toml:"bcrypt_cost"`
}

// HousekeepingConfig carries the instance-admin JWT public keys (a_ha): either
// inline PEM blocks or file references; both are concatenated.
type HousekeepingConfig struct {
	JWTPublicKeys     []string `toml:"jwt_public_keys"`
	JWTPublicKeyFiles []string `toml:"jwt_public_key_files"`
}

// StorageConfig holds runtime storage policy. Retention applies a global
// TimescaleDB drop-chunks policy; zero (the default) disables it, leaving
// per-endpoint TTL (§2.5) as the only expiry.
type StorageConfig struct {
	Retention Duration `toml:"retention"`
}

// SecurityConfig references the master encryption key that seals realm CA
// private keys (§4.3). When MasterKeyFile is empty, the store falls back to
// ASTRATE_MASTER_KEY / ASTRATE_MASTER_KEY_FILE.
type SecurityConfig struct {
	MasterKeyFile string `toml:"master_key_file"`
}

// RealmConfig optionally auto-provisions a realm on boot (§5.1). Name empty
// disables it; when set, JWTPublicKey (PEM) is required.
type RealmConfig struct {
	Name                    string `toml:"name"`
	JWTPublicKey            string `toml:"jwt_public_key"`
	JWTPublicKeyFile        string `toml:"jwt_public_key_file"`
	DeviceRegistrationLimit *int32 `toml:"device_registration_limit"`
}

// LogConfig configures the slog handler (§5.2).
type LogConfig struct {
	Level  string `toml:"level"`
	Format string `toml:"format"`
}

// Default returns the zero-config defaults (single-VPS, dev-friendly).
func Default() Config {
	return Config{
		HTTP: HTTPConfig{Addr: ":8080"},
		MQTT: MQTTConfig{
			Addr:             ":8883",
			DevAddr:          ":1883",
			SessionStorePath: "sessions.db",
			MaxPacketBytes:   272 * 1024,
		},
		Engine: EngineConfig{
			Shards:          16,
			ShardQueue:      4096,
			BatchMaxRows:    64,
			BatchMaxWait:    Duration(50 * time.Millisecond),
			MaxPayloadBytes: 64 * 1024,
		},
		Pairing: PairingConfig{
			CertTTL:      Duration(30 * 24 * time.Hour),
			RegisterRate: 5, RegisterBurst: 10,
			CredentialsRate: 5, CredentialsBurst: 10,
			BcryptCost: 10,
		},
		Log: LogConfig{Level: "info", Format: "json"},
	}
}

// Load reads the config: defaults, then the TOML file at path (skipped when
// path is empty — env-only operation), then ASTRATE_* env overrides, then
// validation. The returned Config is ready to wire.
func Load(path string) (Config, error) {
	cfg := Default()
	if path != "" {
		if _, err := toml.DecodeFile(path, &cfg); err != nil {
			return Config{}, fmt.Errorf("config: reading %s: %w", path, err)
		}
	}
	applyEnv(&cfg)
	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// applyEnv overlays the operationally critical fields from ASTRATE_* variables
// (the documented overrides — full reflection is intentionally avoided).
func applyEnv(cfg *Config) {
	str := func(env string, dst *string) {
		if v, ok := os.LookupEnv(env); ok {
			*dst = v
		}
	}
	str("ASTRATE_HTTP_ADDR", &cfg.HTTP.Addr)
	str("ASTRATE_HTTP_TLS_CERT_FILE", &cfg.HTTP.TLSCertFile)
	str("ASTRATE_HTTP_TLS_KEY_FILE", &cfg.HTTP.TLSKeyFile)
	str("ASTRATE_MQTT_ADDR", &cfg.MQTT.Addr)
	str("ASTRATE_MQTT_TLS_CERT_FILE", &cfg.MQTT.TLSCertFile)
	str("ASTRATE_MQTT_TLS_KEY_FILE", &cfg.MQTT.TLSKeyFile)
	str("ASTRATE_MQTT_SESSION_STORE_PATH", &cfg.MQTT.SessionStorePath)
	str("ASTRATE_MQTT_ADVERTISED_URL", &cfg.MQTT.AdvertisedURL)
	str("ASTRATE_DATABASE_DSN", &cfg.Database.DSN)
	str("ASTRATE_SECURITY_MASTER_KEY_FILE", &cfg.Security.MasterKeyFile)
	str("ASTRATE_REALM_NAME", &cfg.Realm.Name)
	str("ASTRATE_REALM_JWT_PUBLIC_KEY", &cfg.Realm.JWTPublicKey)
	str("ASTRATE_REALM_JWT_PUBLIC_KEY_FILE", &cfg.Realm.JWTPublicKeyFile)
	str("ASTRATE_LOG_LEVEL", &cfg.Log.Level)
	str("ASTRATE_LOG_FORMAT", &cfg.Log.Format)

	if v, ok := os.LookupEnv("ASTRATE_MQTT_INSECURE_DEV_MODE"); ok {
		cfg.MQTT.InsecureDevMode, _ = strconv.ParseBool(v)
	}
	if v, ok := os.LookupEnv("ASTRATE_ENGINE_SHARDS"); ok {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Engine.Shards = n
		}
	}
}

// validate rejects an unusable configuration.
func (c *Config) validate() error {
	if c.Database.DSN == "" {
		return fmt.Errorf("config: database.dsn is required (or set ASTRATE_DATABASE_DSN)")
	}
	if c.HTTP.Addr == "" {
		return fmt.Errorf("config: http.addr is required")
	}
	if c.MQTT.SessionStorePath == "" {
		return fmt.Errorf("config: mqtt.session_store_path is required")
	}
	if !c.MQTT.InsecureDevMode && (c.MQTT.TLSCertFile == "" || c.MQTT.TLSKeyFile == "") {
		return fmt.Errorf("config: mqtt.tls_cert_file and mqtt.tls_key_file are required unless mqtt.insecure_dev_mode is set")
	}
	if (c.HTTP.TLSCertFile == "") != (c.HTTP.TLSKeyFile == "") {
		return fmt.Errorf("config: http.tls_cert_file and http.tls_key_file must be set together")
	}
	if c.Engine.Shards <= 0 {
		return fmt.Errorf("config: engine.shards must be positive")
	}
	if c.Realm.Name != "" && c.Realm.JWTPublicKey == "" && c.Realm.JWTPublicKeyFile == "" {
		return fmt.Errorf("config: realm.jwt_public_key (or _file) is required when realm.name is set")
	}
	switch c.Log.Level {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("config: log.level %q is not one of debug|info|warn|error", c.Log.Level)
	}
	switch c.Log.Format {
	case "json", "text":
	default:
		return fmt.Errorf("config: log.format %q is not json|text", c.Log.Format)
	}
	return nil
}

// HousekeepingKeys resolves the instance-admin JWT public keys, reading any
// referenced files and appending them to the inline blocks.
func (c *Config) HousekeepingKeys() ([]string, error) {
	keys := append([]string(nil), c.Housekeeping.JWTPublicKeys...)
	for _, f := range c.Housekeeping.JWTPublicKeyFiles {
		b, err := os.ReadFile(f) //nolint:gosec // operator-controlled config path
		if err != nil {
			return nil, fmt.Errorf("config: reading housekeeping key %s: %w", f, err)
		}
		keys = append(keys, strings.TrimSpace(string(b)))
	}
	return keys, nil
}

// RealmJWTPublicKey resolves the auto-provision realm's JWT public key,
// preferring the inline block over the file reference.
func (c *Config) RealmJWTPublicKey() (string, error) {
	if c.Realm.JWTPublicKey != "" {
		return c.Realm.JWTPublicKey, nil
	}
	if c.Realm.JWTPublicKeyFile == "" {
		return "", nil
	}
	b, err := os.ReadFile(c.Realm.JWTPublicKeyFile) //nolint:gosec // operator-controlled config path
	if err != nil {
		return "", fmt.Errorf("config: reading realm key %s: %w", c.Realm.JWTPublicKeyFile, err)
	}
	return strings.TrimSpace(string(b)), nil
}
