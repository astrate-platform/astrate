slug: config-fail-loud-engine-shards
verdict: done
at:  83482e4
ran: 2026-09-13T17:57:23Z on DietPi in 150s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
✱ Glob "internal/config/*_test.go" 1 match
→ Read internal/config/config.go [offset=230, limit=90]
✱ Grep "ASTRATE_ENGINE_SHARDS|ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED|t.Setenv|func Test" in internal/config/config_test.go · 31 matches
→ Read internal/config/config_test.go [limit=75, offset=355]
→ Read internal/config/config.go [limit=45, offset=195]
→ Read internal/config/config_test.go [limit=75, offset=1]

← Edit internal/config/config.go
Index: /root/astrate-mule/internal/config/config.go
===================================================================
--- /root/astrate-mule/internal/config/config.go
+++ /root/astrate-mule/internal/config/config.go
@@ -254,12 +254,17 @@
 
 if v, ok := os.LookupEnv("ASTRATE_MQTT_INSECURE_DEV_MODE"); ok {
 	cfg.MQTT.InsecureDevMode, _ = strconv.ParseBool(v)
 }
+// Like the realm-deletion gate above: the shard count is a fail-loud
+// value — a malformed ASTRATE_ENGINE_SHARDS says operator error rather
+// than silently keeping the TOML/default count.
 if v, ok := os.LookupEnv("ASTRATE_ENGINE_SHARDS"); ok {
-	if n, err := strconv.Atoi(v); err == nil {
-		cfg.Engine.Shards = n
+	n, err := strconv.Atoi(v)
+	if err != nil {
+		return fmt.Errorf("config: engine.shards %q must be an integer", v)
 	}
+	cfg.Engine.Shards = n
 }
 
 // The realm default-retention override (#73) accepts both the ASTRATE_-
 // prefixed name and the upstream bare one; the bare name wins when both

Task understood: make `ASTRATE_ENGINE_SHARDS` fail loud on a malformed value, matching the realm-deletion env, and prove it with a test that fails without the change.
← Edit internal/config/config_test.go
Index: /root/astrate-mule/internal/config/config_test.go
===================================================================
--- /root/astrate-mule/internal/config/config_test.go
+++ /root/astrate-mule/internal/config/config_test.go
@@ -351,8 +351,39 @@
 		}
 	})
 }
 
+func TestEngineShardsEnv(t *testing.T) {
+	const env = "ASTRATE_ENGINE_SHARDS"
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
+	t.Run("valid applies", func(t *testing.T) {
+		t.Setenv(env, "4")
+		cfg, err := load(t)
+		if err != nil {
+			t.Fatalf("Load: %v", err)
+		}
+		if cfg.Engine.Shards != 4 {
+			t.Errorf("env %s=4 → shards %d, want 4", env, cfg.Engine.Shards)
+		}
+	})
+	t.Run("banana fails loud", func(t *testing.T) {
+		t.Setenv(env, "banana")
+		if _, err := load(t); err == nil {
+			t.Error("malformed shards env: got nil error")
+		}
+	})
+}
+
 func TestHousekeepingRealmDeletionDisabledEnv(t *testing.T) {
 	const env = "ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"
 	body := `
 [database]

$ go test ./internal/config/ && go vet ./internal/config/ && gofmt -l internal/config/
ok  	github.com/astrate-platform/astrate/internal/config	0.067s

$ go vet ./... && go test ./...
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
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

Full gate green.

```
Done: ASTRATE_ENGINE_SHARDS now fails loud on a malformed value (returns error from applyEnv); added TestEngineShardsEnv proving banana errors.
Files: internal/config/config.go, internal/config/config_test.go
Verified: go vet ./... && go test ./... (non-race) -> pass; gofmt -l internal/config/ -> clean
Unsure: nothing
Follow-ups: none
```
