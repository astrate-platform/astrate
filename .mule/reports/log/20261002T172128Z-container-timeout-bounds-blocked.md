slug: container-timeout-bounds
verdict: blocked
at:  927aa1b
ran: 2026-10-02T17:21:28Z on DietPi in 1201s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/flow/blocks/container/block.go [offset=160, limit=120]
→ Read internal/flow/blocks/container/httpbridge.go
→ Read internal/flow/blocks/container/block.go [limit=160]
→ Read internal/flow/blocks/schema.go [limit=60]

$ ls internal/flow/blocks/container/ && rg -n "parseConfig|timeout_ms|ready_timeout_ms" --glob '!*.md' -l
block.go
block_test.go
docker.go
httpbridge.go
internal/flow/blocks/info.go
internal/flow/blocks/schema.go
internal/flow/blocks/httpblocks.go
internal/flow/blocks/container/block.go
internal/flow/blocks/container/block_test.go
internal/flow/blocks/httpblocks_test.go

→ Read internal/flow/blocks/container/block_test.go

$ rg -n "timeout_ms|timeoutMs|Timeout" internal/flow/blocks/httpblocks.go internal/flow/blocks/mqtt*.go 2>/dev/null | head -50; echo ---; rg -n "must be|maximum|too large|exceed" internal/flow/blocks/*.go | head -40
internal/flow/blocks/httpblocks.go:56://   - timeout_ms (int, default 5000): per-request timeout; must be positive
internal/flow/blocks/httpblocks.go:97:	if cfg.timeout, err = positiveMillis(config, "timeout_ms", cfg.timeout); err != nil {
internal/flow/blocks/httpblocks.go:176:	reqCtx, cancel := context.WithTimeout(ctx, s.timeout)
internal/flow/blocks/httpblocks.go:208://   - timeout_ms (int, default 5000): per-request timeout; must be positive
internal/flow/blocks/httpblocks.go:250:	if cfg.timeout, err = positiveMillis(config, "timeout_ms", cfg.timeout); err != nil {
internal/flow/blocks/httpblocks.go:321:	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
internal/flow/blocks/mqtt.go:28:	mqttConnectTimeout    = 10 * time.Second
internal/flow/blocks/mqtt.go:29:	mqttSubscribeTimeout  = 10 * time.Second
internal/flow/blocks/mqtt.go:30:	mqttPublishTimeout    = 5 * time.Second
internal/flow/blocks/mqtt.go:62:	if !token.WaitTimeout(mqttConnectTimeout) {
internal/flow/blocks/mqtt.go:202:		if !token.WaitTimeout(mqttSubscribeTimeout) {
internal/flow/blocks/mqtt.go:349:	if !token.WaitTimeout(mqttPublishTimeout) {
internal/flow/blocks/mqtt_test.go:59:	if !token.WaitTimeout(5*time.Second) || token.Error() != nil {
internal/flow/blocks/mqtt_test.go:70:	if !token.WaitTimeout(5*time.Second) || token.Error() != nil {
internal/flow/blocks/mqtt_test.go:81:		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
internal/flow/blocks/mqtt_test.go:115:const mqttTestTimeout = 10 * time.Second
internal/flow/blocks/mqtt_test.go:139:	msg := pollEmit(t, src, mqttTestTimeout)[0]
internal/flow/blocks/mqtt_test.go:171:	deadline := time.Now().Add(mqttTestTimeout)
internal/flow/blocks/mqtt_test.go:204:	deadline := time.Now().Add(mqttTestTimeout)
internal/flow/blocks/mqtt_test.go:206:		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
internal/flow/blocks/mqtt_test.go:246:	m := waitForMessage(t, received, mqttTestTimeout)
internal/flow/blocks/mqtt_test.go:261:	m = waitForMessage(t, received, mqttTestTimeout)
internal/flow/blocks/mqtt_test.go:295:	m := waitForMessage(t, received, mqttTestTimeout)
---
internal/flow/blocks/splitmap_test.go:147:	// Mutate output[0]'s Metadata after Process returns; output[1] must be
internal/flow/blocks/sort.go:18://     buffered timestamp before it is released; must be >= 0
internal/flow/blocks/sort.go:21://   - max_buffered (int, default 0): maximum number of messages to buffer
internal/flow/blocks/sort.go:22://     (0 means unlimited). When exceeded, apply overflow_policy
internal/flow/blocks/sort.go:24://     exceeds max_buffered. Values: "drop_oldest", "error", "force_flush"
internal/flow/blocks/sort.go:33:// cause the buffer to exceed max_buffered, the policy controls behavior:
internal/flow/blocks/sort.go:44:		return nil, fmt.Errorf("sort: window_ms too large")
internal/flow/blocks/sort.go:83:		// Now handle overflow if max_buffered is set and buffer still exceeds limit
internal/flow/blocks/sort.go:136:			return cfg, fmt.Errorf("window_ms must be non-negative")
internal/flow/blocks/sort.go:146:			return cfg, fmt.Errorf("max_buffered must be non-negative")
internal/flow/blocks/sort.go:154:			return cfg, fmt.Errorf("max_buffered_count must be non-negative")
internal/flow/blocks/sort.go:161:			return cfg, fmt.Errorf("overflow_policy must be a string")
internal/flow/blocks/sort.go:167:			return cfg, fmt.Errorf("overflow_policy must be one of: drop_oldest, error, force_flush")
internal/flow/blocks/mqtt_test.go:315:		{"qos too big", map[string]any{"url": url, "topics": []string{"t/x"}, "qos": 3}, "mqtt_source: qos must be 0, 1 or 2"},
internal/flow/blocks/mqtt_test.go:316:		{"qos negative", map[string]any{"url": url, "topics": []string{"t/x"}, "qos": -1}, "mqtt_source: qos must be 0, 1 or 2"},
internal/flow/blocks/mqtt_test.go:317:		{"qos fractional", map[string]any{"url": url, "topics": []string{"t/x"}, "qos": 1.5}, "mqtt_source: qos must be 0, 1 or 2"},
internal/flow/blocks/mqtt_test.go:351:		{"qos too big", map[string]any{"url": url, "topic": "t/x", "qos": 3}, "mqtt_sink: qos must be 0, 1 or 2"},
internal/flow/blocks/mqtt_test.go:352:		{"qos fractional", map[string]any{"url": url, "topic": "t/x", "qos": 0.5}, "mqtt_sink: qos must be 0, 1 or 2"},
internal/flow/blocks/randomsource.go:19://   - interval_ms (int, default 1000): must be positive; every Emit (including
internal/flow/blocks/randomsource.go:79:			return cfg, fmt.Errorf("interval_ms must be positive")
internal/flow/blocks/randomsource.go:89:			return cfg, fmt.Errorf("min must be a number")
internal/flow/blocks/randomsource.go:99:			return cfg, fmt.Errorf("max must be a number")
internal/flow/blocks/randomsource.go:107:		return cfg, fmt.Errorf("min must be <= max")
internal/flow/blocks/info.go:80:		Config:  "window_ms (default 1000, must be ≥ 0), dedup (bool, default false), max_buffered (default 0), overflow_policy (default drop_oldest)",
internal/flow/blocks/schema.go:13:	TypeSort:              `{"type":"object","properties":{"window_ms":{"type":"integer","description":"Buffering window in milliseconds (default 1000, must be ≥ 0)"},"dedup":{"type":"boolean","description":"Release each key once per window (default false)"},"max_buffered":{"type":"integer","description":"Maximum number of buffered messages (default 0, unlimited)"},"overflow_policy":{"type":"string","description":"Action when buffer exceeds max_buffered: drop_oldest, error, or force_flush (default drop_oldest)","enum":["drop_oldest","error","force_flush"]}},"required":[]}`,
internal/flow/blocks/randomsource_test.go:36:			want: `random_source: interval_ms must be positive`,
internal/flow/blocks/randomsource_test.go:42:			want: `random_source: interval_ms must be positive`,
internal/flow/blocks/randomsource_test.go:48:			want: `random_source: min must be <= max`,
internal/flow/blocks/sort_test.go:111:	const wantErr = "sort: window_ms must be non-negative"
internal/flow/blocks/mqtt.go:114:		return 0, fmt.Errorf("qos must be 0, 1 or 2")
internal/flow/blocks/mqtt.go:118:		return 0, fmt.Errorf("qos must be 0, 1 or 2")
internal/flow/blocks/mqtt.go:164://   - qos (int, default 0): subscription QoS, must be 0, 1 or 2
internal/flow/blocks/httpblocks.go:34:		return 0, fmt.Errorf("%s must be positive", key)
internal/flow/blocks/httpblocks.go:55://     first); must be positive
internal/flow/blocks/httpblocks.go:56://   - timeout_ms (int, default 5000): per-request timeout; must be positive
internal/flow/blocks/httpblocks.go:124:		return nil, fmt.Errorf("urls must be a list of strings")
internal/flow/blocks/httpblocks.go:208://   - timeout_ms (int, default 5000): per-request timeout; must be positive
internal/flow/blocks/httpblocks.go:279:				return nil, fmt.Errorf("headers values must be strings")
internal/flow/blocks/httpblocks.go:285:		return nil, fmt.Errorf("headers must be an object of string→string")
internal/flow/blocks/httpblocks_test.go:168:		t.Fatal("Emit succeeded despite server exceeding timeout_ms")

→ Read internal/flow/blocks/sort.go [offset=28, limit=60]
→ Read internal/flow/blocks/httpblocks.go [offset=22, limit=45]
Task understood: bound `timeout_ms`/`ready_timeout_ms` in `parseConfig` so overflow and absurd values are rejected. First, verifying the claims with a throwaway probe in /tmp.

$ mkdir -p /tmp/opencode/probe && cat > /tmp/opencode/probe/main.go <<'EOF'
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

func main() {
	// 1. wrapping of time.Duration(n) * time.Millisecond
	for _, n := range []int64{18446744073709, 9223372036854, 9223372036854775807, 1 << 62} {
		d := time.Duration(n) * time.Millisecond
		fmt.Printf("n=%d -> %v (positive=%v)\n", n, d, d > 0)
	}
	fmt.Println("MaxInt64/ms =", int64(^uint64(0)>>1)/int64(time.Millisecond))

	// 2. net/http: non-positive Client.Timeout = no deadline
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()
	for _, to := range []time.Duration{-551 * time.Microsecond, 0, 50 * time.Millisecond} {
		start := time.Now()
		resp, err := (&http.Client{Timeout: to}).Get(srv.URL)
		if resp != nil {
			resp.Body.Close()
		}
		fmt.Printf("Timeout=%v -> err=%v after %v\n", to, err, time.Since(start).Round(10*time.Millisecond))
	}
}
EOF
cd /tmp/opencode/probe && go mod init probe >/dev/null 2>&1; go run main.go
n=18446744073709 -> -551.616µs (positive=false)
n=9223372036854 -> 2562047h47m16.854s (positive=true)
n=9223372036854775807 -> -1ms (positive=false)
n=4611686018427387904 -> 0s (positive=false)
MaxInt64/ms = 9223372036854
Timeout=-551µs -> err=<nil> after 3.01s
Timeout=0s -> err=<nil> after 3s
Timeout=50ms -> err=Get "http://127.0.0.1:39781": context deadline exceeded (Client.Timeout exceeded while awaiting headers) after 50ms

$ rg -n "ready_timeout_ms|timeout_ms" docs/ internal/flow/blocks/info.go 2>/dev/null | head -20; echo ---; rg -n "TypeContainer" -A3 internal/flow/blocks/info.go
internal/flow/blocks/info.go:92:		Config:  "urls (required string array, at least one), interval_ms (default 1000), timeout_ms (default 5000)",
internal/flow/blocks/info.go:98:		Config:  `url (required), method (default "POST"), timeout_ms (default 5000), headers (object string→string)`,
internal/flow/blocks/info.go:116:		Config:  "image (required), config (object→ASTRATE_FLOW_CONFIG), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",
---
112:	TypeContainer: {
113:		Type:    TypeContainer,
114-		Role:    RoleTransform,
115-		Summary: "Run a Docker image as a transform (HTTP POST /v1/message); local Docker only (PoC)",
116-		Config:  "image (required), config (object→ASTRATE_FLOW_CONFIG), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",

$ rg -n "builtinSchemas|Schema\(" --glob '*.go' internal/ | head -20; echo ---; rg -rn "minimum|maximum" internal/flow/blocks/*.go | head; echo ---; rg -n "parseConfig" -g '*.go' internal/
internal/flowapi/userblocks.go:134:	if err := s.validateConfigSchema(configSchema); err != nil {
internal/flowapi/userblocks.go:146:func (s *Service) validateConfigSchema(schema json.RawMessage) error {
internal/flowapi/userblocks.go:150:	if _, err := compileJSONSchema(schema); err != nil {
internal/flow/blocks/info.go:143:		if s := builtinSchemas[blockType]; s != "" {
internal/flow/blocks/schema.go:3:// builtinSchemas maps block_type → Draft-07-style object-schema JSON for the
internal/flow/blocks/schema.go:5:var builtinSchemas = map[string]string{
internal/flowapi/service.go:825:		sch, err := compileJSONSchema(ub.ConfigSchema)
internal/flowapi/service.go:836:func compileJSONSchema(schema json.RawMessage) (*jsonschema.Schema, error) {
internal/flow/blocks/schema_test.go:11:		if _, ok := builtinSchemas[typ]; !ok {
internal/flow/blocks/schema_test.go:12:			t.Errorf("builtinSchemas missing entry for %q (builtinInfo has docs)", typ)
internal/flow/blocks/schema_test.go:27:	keys := make([]string, 0, len(builtinSchemas))
internal/flow/blocks/schema_test.go:28:	for k := range builtinSchemas {
internal/flow/blocks/schema_test.go:33:		if err := json.Unmarshal([]byte(builtinSchemas[typ]), &schema); err != nil {
internal/flow/blocks/schema_test.go:60:	if err := json.Unmarshal([]byte(builtinSchemas[TypeContainer]), &container); err != nil {
internal/flow/blocks/schema_test.go:72:	if err := json.Unmarshal([]byte(builtinSchemas[TypeRandomSource]), &randomSource); err != nil {
internal/flow/blocks/schema_test.go:84:	if err := json.Unmarshal([]byte(builtinSchemas[TypeHTTPSource]), &httpSource); err != nil {
internal/flow/blocks/schema_test.go:94:	if err := json.Unmarshal([]byte(builtinSchemas[TypeNullSink]), &nullSink); err != nil {
internal/flow/blocks/schema_test.go:102:func TestLookupInfo_ConfigSchema(t *testing.T) {
internal/flow/blocks/mqtt_test.go:373:func TestMQTTRegistration_InfoAndSchema(t *testing.T) {
internal/flow/message_test.go:218:func TestFlowMessage_WireFormatSchema(t *testing.T) {
---
internal/flow/blocks/catalog.go:// DefaultRegistry returns a registry with the n useful built-in set:
internal/flow/blocks/sort.go://   - max_buffered (int, default 0): n number of messages to buffer
---
internal/flow/blocks/container/block.go:74:	cfg, err := parseConfig(config)
internal/flow/blocks/container/block.go:172:func parseConfig(config map[string]any) (Config, error) {

→ Read internal/flow/blocks/schema_test.go [limit=115]
All three claims confirmed by probe. Now the change.

← Edit internal/flow/blocks/container/block.go
Index: /root/astrate-mule/internal/flow/blocks/container/block.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/container/block.go
+++ /root/astrate-mule/internal/flow/blocks/container/block.go
@@ -23,8 +23,16 @@
 const Type = "container"
 
 const defaultContainerPort = 8080
 
+// maxMillis bounds every operator-supplied *_ms duration in the block
+// config. Past math.MaxInt64/int64(time.Millisecond) (~292 years) the
+// multiply into time.Duration wraps negative, and a non-positive Timeout
+// is read two ways: Bridge.timeout() falls back to defaultTimeout while
+// net/http treats it as "no deadline at all". One hour is far above any
+// useful per-message or readiness wait and keeps the value representable.
+const maxMillis = int64(time.Hour / time.Millisecond)
+
 // defaultRunner is used by the catalog constructor. Tests inject via New.
 var defaultRunner Runner = &CLIRunner{}
 
 // SetDefaultRunner replaces the Docker runner used by the catalog constructor.

← Edit internal/flow/blocks/container/block.go
Index: /root/astrate-mule/internal/flow/blocks/container/block.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/container/block.go
+++ /root/astrate-mule/internal/flow/blocks/container/block.go
@@ -228,8 +228,27 @@
 
 	return cfg, nil
 }
 
+// boundedMillis reads key from config as a millisecond duration, applying def
+// when the key is absent or null. Shared by both duration keys so their
+// validation messages cannot drift. A value above maxMillis is rejected by
+// name and by value rather than silently wrapping into a nonsense Duration.
+func boundedMillis(config map[string]any, key string, def time.Duration) (time.Duration, error) {
+	v, ok := config[key]
+	if !ok || v == nil {
+		return def, nil
+	}
+	n, err := asInt(v)
+	if err != nil || n <= 0 {
+		return 0, fmt.Errorf("container: %s must be a positive integer", key)
+	}
+	if int64(n) > maxMillis {
+		return 0, fmt.Errorf("container: %s must be at most %d ms, got %d", key, maxMillis, n)
+	}
+	return time.Duration(n) * time.Millisecond, nil
+}
+
 func asInt(v any) (int, error) {
 	switch n := v.(type) {
 	case int:
 		return n, nil

Now wire `parseConfig` to the helper:
← Edit internal/flow/blocks/container/block.go
Index: /root/astrate-mule/internal/flow/blocks/container/block.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/container/block.go
+++ /root/astrate-mule/internal/flow/blocks/container/block.go
@@ -209,23 +209,19 @@
 		}
 		cfg.ContainerPort = n
 	}
 
-	if v, ok := config["timeout_ms"]; ok && v != nil {
-		n, err := asInt(v)
-		if err != nil || n <= 0 {
-			return cfg, fmt.Errorf("container: timeout_ms must be a positive integer")
-		}
-		cfg.Timeout = time.Duration(n) * time.Millisecond
+	timeout, err := boundedMillis(config, "timeout_ms", defaultTimeout)
+	if err != nil {
+		return cfg, err
 	}
+	cfg.Timeout = timeout
 
-	if v, ok := config["ready_timeout_ms"]; ok && v != nil {
-		n, err := asInt(v)
-		if err != nil || n <= 0 {
-			return cfg, fmt.Errorf("container: ready_timeout_ms must be a positive integer")
-		}
-		cfg.ReadyTimeout = time.Duration(n) * time.Millisecond
+	readyTimeout, err := boundedMillis(config, "ready_timeout_ms", defaultReadyWait)
+	if err != nil {
+		return cfg, err
 	}
+	cfg.ReadyTimeout = readyTimeout
 
 	return cfg, nil
 }
 

→ Read internal/flow/blocks/container/block.go [limit=65, offset=180]

← Edit internal/flow/blocks/schema.go
Index: /root/astrate-mule/internal/flow/blocks/schema.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/schema.go
+++ /root/astrate-mule/internal/flow/blocks/schema.go
@@ -15,9 +15,9 @@
 	TypeHTTPSource:        `{"type":"object","properties":{"urls":{"type":"array","description":"URLs to poll GET round-robin (at least one)"},"interval_ms":{"type":"integer","description":"Milliseconds between polls (default 1000)"},"timeout_ms":{"type":"integer","description":"Per-request timeout in milliseconds (default 5000)"}},"required":["urls"]}`,
 	TypeHTTPSink:          `{"type":"object","properties":{"url":{"type":"string","description":"Base URL to POST payloads to"},"method":{"type":"string","description":"HTTP method (default \"POST\")"},"timeout_ms":{"type":"integer","description":"Per-request timeout in milliseconds (default 5000)"},"headers":{"type":"object","description":"Extra request headers string→string"}},"required":["url"]}`,
 	TypeMQTTSource:        `{"type":"object","properties":{"url":{"type":"string","description":"Broker URL (e.g. tcp://127.0.0.1:1883)"},"topics":{"type":"array","description":"Topics to subscribe to (at least one, non-empty strings)"},"qos":{"type":"integer","description":"Subscription QoS level 0|1|2 (default 0)"},"client_id":{"type":"string","description":"Optional MQTT client ID (random when unset)"},"username":{"type":"string","description":"Optional username"},"password":{"type":"string","description":"Optional password"}},"required":["url","topics"]}`,
 	TypeMQTTSink:          `{"type":"object","properties":{"url":{"type":"string","description":"Broker URL (e.g. tcp://127.0.0.1:1883)"},"topic":{"type":"string","description":"Topic to publish to"},"qos":{"type":"integer","description":"Publish QoS level 0|1|2 (default 0)"},"retained":{"type":"boolean","description":"Set the retain flag on published messages (default false)"},"client_id":{"type":"string","description":"Optional MQTT client ID (random when unset)"},"username":{"type":"string","description":"Optional username"},"password":{"type":"string","description":"Optional password"}},"required":["url","topic"]}`,
-	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as ASTRATE_FLOW_CONFIG"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
+	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as ASTRATE_FLOW_CONFIG"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","minimum":1,"maximum":3600000,"description":"Per-message POST timeout in milliseconds (default 5000, maximum 3600000)"},"ready_timeout_ms":{"type":"integer","minimum":1,"maximum":3600000,"description":"Health-check wait in milliseconds (default 15000, maximum 3600000)"}},"required":["image"]}`,
 	TypeVirtualDevicePool: `{"type":"object","properties":{"devices":{"type":"array","items":{"type":"string"},"minItems":1,"description":"Registered device_ids this pool may publish as"}},"required":["devices"]}`,
 	TypeNullSink:          `{"type":"object","properties":{}}`,
 	TypeLogSink:           `{"type":"object","properties":{}}`,
 }

Now the test — table-driven, asserting `New` refuses the values it accepts today.
← Edit internal/flow/blocks/container/block_test.go
Index: /root/astrate-mule/internal/flow/blocks/container/block_test.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/container/block_test.go
+++ /root/astrate-mule/internal/flow/blocks/container/block_test.go
@@ -467,8 +467,85 @@
 		t.Fatalf("stops after failed ready = %d, want 1", inst.stops.Load())
 	}
 }
 
+// TestNew_RejectsOutOfBoundsDurations: both *_ms keys are multiplied into a
+// time.Duration unchecked, so a value past math.MaxInt64/1e6 wraps negative
+// (measured: 18446744073709 -> -551.616us, which Bridge.timeout() reads as
+// "use the 5s default" while net/http reads it as no deadline at all) and
+// 9223372036854 stays positive at 2562047h47m16s, so New never returns.
+// Both must be refused by name and value instead.
+func TestNew_RejectsOutOfBoundsDurations(t *testing.T) {
+	cases := []struct {
+		name   string
+		key    string
+		value  int64
+		wantIn string
+	}{
+		{"timeout_ms wraps negative", "timeout_ms", 18446744073709,
+			"container: timeout_ms must be at most 3600000 ms, got 18446744073709"},
+		{"timeout_ms absurd but positive", "timeout_ms", 9223372036854,
+			"container: timeout_ms must be at most 3600000 ms, got 9223372036854"},
+		{"timeout_ms max int64", "timeout_ms", math.MaxInt64,
+			"container: timeout_ms must be at most 3600000 ms, got 9223372036854775807"},
+		{"ready_timeout_ms wraps negative", "ready_timeout_ms", 18446744073709,
+			"container: ready_timeout_ms must be at most 3600000 ms, got 18446744073709"},
+		{"ready_timeout_ms absurd but positive", "ready_timeout_ms", 9223372036854,
+			"container: ready_timeout_ms must be at most 3600000 ms, got 9223372036854"},
+		{"ready_timeout_ms max int64", "ready_timeout_ms", math.MaxInt64,
+			"container: ready_timeout_ms must be at most 3600000 ms, got 9223372036854775807"},
+	}
+	for _, tc := range cases {
+		t.Run(tc.name, func(t *testing.T) {
+			srv := echoServer(t)
+			defer srv.Close()
+
+			r := &fakeRunner{inst: &fakeInstance{base: srv.URL, id: "bounds"}}
+			cfg := map[string]any{"image": "img", tc.key: tc.value}
+			b, err := container.New("c", cfg, flow.Deps{}, r)
+			if err == nil {
+				_ = b.(flow.Stopper).Stop()
+				t.Fatalf("New accepted %s=%d, want refusal", tc.key, tc.value)
+			}
+			if err.Error() != tc.wantIn {
+				t.Errorf("err = %q, want %q", err, tc.wantIn)
+			}
+			if r.started.Load() != 0 {
+				t.Errorf("runner started %d containers, want 0 (config must be rejected first)", r.started.Load())
+			}
+		})
+	}
+}
+
+// TestNew_AcceptsBoundAndDefaults guards the other side: the documented
+// maximum is accepted, and absent keys keep their defaults.
+func TestNew_AcceptsBoundAndDefaults(t *testing.T) {
+	cases := []struct {
+		name string
+		cfg  map[string]any
+	}{
+		{"defaults", map[string]any{}},
+		{"at maximum", map[string]any{"timeout_ms": 3600000, "ready_timeout_ms": 3600000}},
+		{"one millisecond", map[string]any{"timeout_ms": 1, "ready_timeout_ms": 1}},
+	}
+	for _, tc := range cases {
+		t.Run(tc.name, func(t *testing.T) {
+			srv := echoServer(t)
+			defer srv.Close()
+
+			cfg := map[string]any{"image": "img"}
+			for k, v := range tc.cfg {
+				cfg[k] = v
+			}
+			b, err := container.New("c", cfg, flow.Deps{}, &fakeRunner{inst: &fakeInstance{base: srv.URL, id: "ok"}})
+			if err != nil {
+				t.Fatalf("New(%v): %v", tc.cfg, err)
+			}
+			_ = b.(flow.Stopper).Stop()
+		})
+	}
+}
+
 // waitingInstance adds Waiter so the death watcher engages.
 type waitingInstance struct {
 	fakeInstance
 	waitFn func(ctx context.Context) (int, error)

← Edit internal/flow/blocks/container/block_test.go
Index: /root/astrate-mule/internal/flow/blocks/container/block_test.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/container/block_test.go
+++ /root/astrate-mule/internal/flow/blocks/container/block_test.go
@@ -4,8 +4,9 @@
 "context"
 "encoding/json"
 "fmt"
 "io"
+"math"
 "net/http"
 "net/http/httptest"
 "strings"
 "sync"

$ gofmt -l . && go vet ./internal/flow/blocks/... && go test ./internal/flow/blocks/... 2>&1 | tail -20
# github.com/astrate-platform/astrate/internal/flow/blocks/container_test
# [github.com/astrate-platform/astrate/internal/flow/blocks/container_test]
vet: internal/flow/blocks/container/block_test.go:507:9: b.(flow.Stopper).Stop() (no value) used as value

$ rg -n "func.*Stop\(\)" internal/flow/flow.go internal/flow/*.go | head
internal/flow/factory_test.go:149:func (b *stoppableBlock) Stop() { *b.stopped++ }

$ python3 - <<'EOF'
import re
p='internal/flow/blocks/container/block_test.go'
s=open(p).read()
s=s.replace("""			if err == nil {
				_ = b.(flow.Stopper).Stop()
				t.Fatalf("New accepted %s=%d, want refusal", tc.key, tc.value)
			}""","""			if err == nil {
				b.(flow.Stopper).Stop()
				t.Fatalf("New accepted %s=%d, want refusal", tc.key, tc.value)
			}""")
s=s.replace("""			_ = b.(flow.Stopper).Stop()""","""			b.(flow.Stopper).Stop()""")
open(p,'w').write(s)
EOF
gofmt -l . ; go vet ./internal/flow/blocks/... && go test ./internal/flow/blocks/... 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	1.295s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
2026/10/02 19:17:52 ERROR container process error block=c err="container: http 500: boom"
--- FAIL: TestNew_AcceptsBoundAndDefaults (0.01s)
    --- FAIL: TestNew_AcceptsBoundAndDefaults/one_millisecond (0.00s)
        block_test.go:542: New(map[ready_timeout_ms:1 timeout_ms:1]): container "c": container: not ready: context deadline exceeded (last: Get "http://127.0.0.1:34201/healthz": context deadline exceeded)
FAIL
FAIL	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.545s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
FAIL

← Edit internal/flow/blocks/container/block_test.go
Index: /root/astrate-mule/internal/flow/blocks/container/block_test.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/container/block_test.go
+++ /root/astrate-mule/internal/flow/blocks/container/block_test.go
@@ -525,9 +525,9 @@
 	cfg  map[string]any
 }{
 	{"defaults", map[string]any{}},
 	{"at maximum", map[string]any{"timeout_ms": 3600000, "ready_timeout_ms": 3600000}},
-	{"one millisecond", map[string]any{"timeout_ms": 1, "ready_timeout_ms": 1}},
+	{"small positive", map[string]any{"timeout_ms": 500, "ready_timeout_ms": 2000}},
 }
 for _, tc := range cases {
 	t.Run(tc.name, func(t *testing.T) {
 		srv := echoServer(t)

$ gofmt -l . ; go vet ./... && go test ./internal/flow/... 2>&1 | tail -15
ok  	github.com/astrate-platform/astrate/internal/flow	0.557s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.544s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)

$ cp internal/flow/blocks/container/block.go /tmp/opencode/block.go.new && python3 - <<'EOF'
p='internal/flow/blocks/container/block.go'
s=open(p).read()
new = """	timeout, err := boundedMillis(config, "timeout_ms", defaultTimeout)
	if err != nil {
		return cfg, err
	}
	cfg.Timeout = timeout

	readyTimeout, err := boundedMillis(config, "ready_timeout_ms", defaultReadyWait)
	if err != nil {
		return cfg, err
	}
	cfg.ReadyTimeout = readyTimeout
"""
old = """	if v, ok := config["timeout_ms"]; ok && v != nil {
		n, err := asInt(v)
		if err != nil || n <= 0 {
			return cfg, fmt.Errorf("container: timeout_ms must be a positive integer")
		}
		cfg.Timeout = time.Duration(n) * time.Millisecond
	}

	if v, ok := config["ready_timeout_ms"]; ok && v != nil {
		n, err := asInt(v)
		if err != nil || n <= 0 {
			return cfg, fmt.Errorf("container: ready_timeout_ms must be a positive integer")
		}
		cfg.ReadyTimeout = time.Duration(n) * time.Millisecond
	}
"""
assert new in s
s = s.replace(new, old).replace('\t"sync"\n', '\t"sync"\n\t"golang.org/x/sync/errgroup"\n')  # keep helper used
open(p,'w').write(s)
EOF
go test ./internal/flow/blocks/container/ -run 'TestNew_RejectsOutOfBoundsDurations|TestNew_AcceptsBoundAndDefaults' -timeout 30s 2>&1 | tail -25; cp /tmp/opencode/block.go.new internal/flow/blocks/container/block.go
# github.com/astrate-platform/astrate/internal/flow/blocks/container
internal/flow/blocks/container/block.go:17:2: "golang.org/x/sync/errgroup" imported and not used
FAIL	github.com/astrate-platform/astrate/internal/flow/blocks/container [build failed]
FAIL

$ cp internal/flow/blocks/container/block.go /tmp/opencode/block.go.new && python3 - <<'EOF'
p='internal/flow/blocks/container/block.go'
s=open(p).read()
new = """	timeout, err := boundedMillis(config, "timeout_ms", defaultTimeout)
	if err != nil {
		return cfg, err
	}
	cfg.Timeout = timeout

	readyTimeout, err := boundedMillis(config, "ready_timeout_ms", defaultReadyWait)
	if err != nil {
		return cfg, err
	}
	cfg.ReadyTimeout = readyTimeout
"""
old = """	if v, ok := config["timeout_ms"]; ok && v != nil {
		n, err := asInt(v)
		if err != nil || n <= 0 {
			return cfg, fmt.Errorf("container: timeout_ms must be a positive integer")
		}
		cfg.Timeout = time.Duration(n) * time.Millisecond
	}

	if v, ok := config["ready_timeout_ms"]; ok && v != nil {
		n, err := asInt(v)
		if err != nil || n <= 0 {
			return cfg, fmt.Errorf("container: ready_timeout_ms must be a positive integer")
		}
		cfg.ReadyTimeout = time.Duration(n) * time.Millisecond
	}
"""
assert new in s
s = s.replace(new, old).replace("func dockerName(", "var _ = boundedMillis\n\nfunc dockerName(")
open(p,'w').write(s)
EOF
go test ./internal/flow/blocks/container/ -run 'TestNew_RejectsOutOfBoundsDurations|TestNew_AcceptsBoundAndDefaults' -timeout 30s 2>&1 | tail -25; cp /tmp/opencode/block.go.new internal/flow/blocks/container/block.go
--- FAIL: TestNew_RejectsOutOfBoundsDurations (0.02s)
    --- FAIL: TestNew_RejectsOutOfBoundsDurations/timeout_ms_wraps_negative (0.00s)
        block_test.go:508: New accepted timeout_ms=18446744073709, want refusal
    --- FAIL: TestNew_RejectsOutOfBoundsDurations/timeout_ms_absurd_but_positive (0.00s)
        block_test.go:508: New accepted timeout_ms=9223372036854, want refusal
    --- FAIL: TestNew_RejectsOutOfBoundsDurations/timeout_ms_max_int64 (0.00s)
        block_test.go:508: New accepted timeout_ms=9223372036854775807, want refusal
    --- FAIL: TestNew_RejectsOutOfBoundsDurations/ready_timeout_ms_wraps_negative (0.00s)
        block_test.go:511: err = "container \"c\": container: not ready after wait", want "container: ready_timeout_ms must be at most 3600000 ms, got 18446744073709"
        block_test.go:514: runner started 1 containers, want 0 (config must be rejected first)
    --- FAIL: TestNew_RejectsOutOfBoundsDurations/ready_timeout_ms_absurd_but_positive (0.00s)
        block_test.go:508: New accepted ready_timeout_ms=9223372036854, want refusal
    --- FAIL: TestNew_RejectsOutOfBoundsDurations/ready_timeout_ms_max_int64 (0.00s)
        block_test.go:511: err = "container \"c\": container: not ready after wait", want "container: ready_timeout_ms must be at most 3600000 ms, got 9223372036854775807"
        block_test.go:514: runner started 1 containers, want 0 (config must be rejected first)
FAIL
FAIL	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.052s
FAIL

