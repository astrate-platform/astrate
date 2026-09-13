slug: config-fail-loud-dev-mode
verdict: done
at:  27d8165
ran: 2026-09-13T18:03:48Z on DietPi in 227s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/config/config.go [offset=220, limit=80]
→ Read internal/config/config.go [offset=299]
✱ Glob "internal/config/*.go" 2 matches
→ Read internal/config/config_test.go
→ Read internal/config/config.go [limit=30]
Now I understand the task: make `ASTRATE_MQTT_INSECURE_DEV_MODE` fail loud on a bad boolean, like the sibling env vars do. Let me make the change and add a test.

← Edit internal/config/config.go
Index: /root/astrate-mule/internal/config/config.go
===================================================================
--- /root/astrate-mule/internal/config/config.go
+++ /root/astrate-mule/internal/config/config.go
@@ -252,9 +252,13 @@
 str("ASTRATE_LOG_LEVEL", &cfg.Log.Level)
 str("ASTRATE_LOG_FORMAT", &cfg.Log.Format)
 
 if v, ok := os.LookupEnv("ASTRATE_MQTT_INSECURE_DEV_MODE"); ok {
-	cfg.MQTT.InsecureDevMode, _ = strconv.ParseBool(v)
+	b, err := strconv.ParseBool(v)
+	if err != nil {
+		return fmt.Errorf("config: mqtt.insecure_dev_mode %q must be a boolean", v)
+	}
+	cfg.MQTT.InsecureDevMode = b
 }
 // Like the realm-deletion gate above: the shard count is a fail-loud
 // value — a malformed ASTRATE_ENGINE_SHARDS says operator error rather
 // than silently keeping the TOML/default count.

← Edit internal/config/config_test.go
Index: /root/astrate-mule/internal/config/config_test.go
===================================================================
--- /root/astrate-mule/internal/config/config_test.go
+++ /root/astrate-mule/internal/config/config_test.go
@@ -382,8 +382,49 @@
 		}
 	})
 }
 
+func TestMQTTInsecureDevModeEnv(t *testing.T) {
+	const env = "ASTRATE_MQTT_INSECURE_DEV_MODE"
+	body := `
+[database]
+dsn = "x"
+[mqtt]
+insecure_dev_mode = true
+`
+	load := func(t *testing.T) (Config, error) {
+		t.Helper()
+		return Load(writeTOML(t, body))
+	}
+
+	t.Run("env overrides toml", func(t *testing.T) {
+		t.Setenv(env, "false")
+		cfg, err := load(t)
+		if err != nil {
+			t.Fatalf("Load: %v", err)
+		}
+		if cfg.MQTT.InsecureDevMode {
+			t.Errorf("env %s=false → insecure_dev_mode true, want false", env)
+		}
+	})
+	t.Run("true parses", func(t *testing.T) {
+		t.Setenv(env, "true")
+		cfg, err := load(t)
+		if err != nil {
+			t.Fatalf("Load: %v", err)
+		}
+		if !cfg.MQTT.InsecureDevMode {
+			t.Errorf("env %s=true → insecure_dev_mode false, want true", env)
+		}
+	})
+	t.Run("yes fails loud", func(t *testing.T) {
+		t.Setenv(env, "yes")
+		if _, err := load(t); err == nil {
+			t.Error("malformed insecure_dev_mode env: got nil error")
+		}
+	})
+}
+
 func TestHousekeepingRealmDeletionDisabledEnv(t *testing.T) {
 	const env = "ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"
 	body := `
 [database]

Now verify the fix passes the gate, then confirm the new test actually fails without the fix (what the runner will check).
$ gofmt -l internal/config/ && go test ./internal/config/ && go vet ./internal/config/
--- FAIL: TestMQTTInsecureDevModeEnv (0.00s)
    --- FAIL: TestMQTTInsecureDevModeEnv/env_overrides_toml (0.00s)
        config_test.go:403: Load: config: mqtt.tls_cert_file and mqtt.tls_key_file are required unless mqtt.insecure_dev_mode is set
FAIL
FAIL	github.com/astrate-platform/astrate/internal/config	0.071s
FAIL

The TOML body needs TLS files so dev-mode `false` still validates. Let me fix the test body.
✗ Edit internal/config/config_test.go failed
Error: Found multiple matches for oldString. Provide more surrounding context to make the match unique.

← Edit internal/config/config_test.go
Index: /root/astrate-mule/internal/config/config_test.go
===================================================================
--- /root/astrate-mule/internal/config/config_test.go
+++ /root/astrate-mule/internal/config/config_test.go
@@ -389,8 +389,10 @@
 [database]
 dsn = "x"
 [mqtt]
 insecure_dev_mode = true
+tls_cert_file = "cert.pem"
+tls_key_file = "key.pem"
 `
 	load := func(t *testing.T) (Config, error) {
 		t.Helper()
 		return Load(writeTOML(t, body))

$ gofmt -l internal/config/ ; go test ./internal/config/ && go vet ./internal/config/
ok  	github.com/astrate-platform/astrate/internal/config	0.053s

Gate passes. Now verify the test breaks when the implementation is reverted (the runner's criterion):
$ cp internal/config/config.go /tmp/opencode/config.go.bak && perl -0pi -e 's/if v, ok := os\.LookupEnv\("ASTRATE_MQTT_INSECURE_DEV_MODE"\); ok \{\n\t\tb, err := strconv\.ParseBool\(v\)\n\t\tif err != nil \{\n\t\t\treturn fmt\.Errorf\("config: mqtt\.insecure_dev_mode %q must be a boolean", v\)\n\t\t\}\n\t\tcfg\.MQTT\.InsecureDevMode = b\n\t\}/if v, ok := os.LookupEnv("ASTRATE_MQTT_INSECURE_DEV_MODE"); ok {\n\t\tcfg.MQTT.InsecureDevMode, _ = strconv.ParseBool(v)\n\t}/' internal/config/config.go && go test ./internal/config/ -run TestMQTTInsecureDevModeEnv; cp /tmp/opencode/config.go.bak internal/config/config.go && echo RESTORED
--- FAIL: TestMQTTInsecureDevModeEnv (0.00s)
    --- FAIL: TestMQTTInsecureDevModeEnv/yes_fails_loud (0.00s)
        config_test.go:424: malformed insecure_dev_mode env: got nil error
FAIL
FAIL	github.com/astrate-platform/astrate/internal/config	0.016s
FAIL
RESTORED

Confirmed the test fails without the fix. Re-verify the restored implementation passes the full gate:
$ gofmt -l . ; go vet ./... && go test ./... 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

Done: `applyEnv` now returns the parse error for a malformed `ASTRATE_MQTT_INSECURE_DEV_MODE` instead of silently disabling dev mode; a new test proves `ASTRATE_MQTT_INSECURE_DEV_MODE=yes` makes Load fail (verified it breaks when the fix is reverted).
Files: internal/config/config.go, internal/config/config_test.go
Verified: `gofmt -l .`, `go vet ./...`, `go test ./...` -> pass (all); reverting the fix makes `TestMQTTInsecureDevModeEnv/yes_fails_loud` fail as required
Unsure: nothing
Follow-ups: none
