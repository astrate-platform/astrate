slug: forward-static-headers-override-test
verdict: done
at:  3a0f153
ran: 2026-09-26T18:19:16Z on DietPi in 263s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/engine/forward/http.go
→ Read internal/engine/forward/http_test.go [limit=140]

$ rg -n "static_headers|StaticHeaders" --glob '!internal/engine/forward/*' -S . | head -40
./cmd/astrate/main.go:265:			StaticHeaders: f.StaticHeaders,
./internal/config/config.go:169:	StaticHeaders map[string]string `toml:"static_headers"`
./internal/config/config.go:345:	// static_headers — a stale url left behind after disabling forwarding must
./internal/config/config.example.toml:60:# static_headers = { Authorization = "Bearer ..." }
./internal/config/config_test.go:266:static_headers = { Authorization = "Bearer tok", X-Forward = "yes" }
./internal/config/config_test.go:282:	if f.StaticHeaders["Authorization"] != "Bearer tok" || f.StaticHeaders["X-Forward"] != "yes" {
./internal/config/config_test.go:283:		t.Errorf("StaticHeaders = %v", f.StaticHeaders)
./internal/engine/triggers/match_test.go:561:			"http_static_headers": {"X-Custom": "v"}, "ignore_ssl_errors": true},
./internal/engine/triggers/match_test.go:564:	if headers.Action.Method != "DELETE" || headers.Action.StaticHeaders["X-Custom"] != "v" || !headers.Action.IgnoreSSLErrors {
./internal/engine/triggers/actions.go:30:	// StaticHeaders are the action's extra request headers.
./internal/engine/triggers/actions.go:31:	StaticHeaders map[string]string
./internal/engine/triggers/actions.go:52:	HTTPStaticHeaders map[string]string `json:"http_static_headers"`
./internal/engine/triggers/actions.go:94:		if len(h.HTTPStaticHeaders) > 0 {
./internal/engine/triggers/actions.go:95:			fe.add("http_static_headers", "must be blank")
./internal/engine/triggers/actions.go:104:			StaticHeaders: h.HTTPStaticHeaders, IgnoreSSLErrors: h.IgnoreSSLErrors,
./internal/engine/triggers/actions.go:121:		validateStaticHeaders(fe, h.HTTPStaticHeaders)
./internal/engine/triggers/actions.go:130:			StaticHeaders: h.HTTPStaticHeaders, IgnoreSSLErrors: h.IgnoreSSLErrors,
./internal/engine/triggers/actions.go:613:	for k, v := range a.StaticHeaders {
./internal/engine/triggers/validation.go:138:// validateStaticHeaders applies the modern-branch header checks: the
./internal/engine/triggers/validation.go:140:func validateStaticHeaders(e *FieldErrors, headers map[string]string) {
./internal/engine/triggers/validation.go:144:			e.add("http_static_headers", "must contain only allowed http headers")
./internal/engine/triggers/validation.go:153:		e.add("http_static_headers", fmt.Sprintf("headers total size must be lower than %d", maxHeadersBytes))
./internal/engine/triggers/actions_test.go:110:			`"http_static_headers":{"X-Foo":"bar"},"ignore_ssl_errors":true}`))
./internal/engine/triggers/actions_test.go:114:		if a.Method != http.MethodPut || a.StaticHeaders["X-Foo"] != "bar" || !a.IgnoreSSLErrors {
./internal/engine/triggers/actions_test.go:277:			raw:  `{"http_post_url":"https://y/legacy","http_static_headers":{"X-A":"b"}}`,
./internal/engine/triggers/actions_test.go:278:			want: map[string][]string{"http_static_headers": {"must be blank"}},
./internal/engine/triggers/actions_test.go:290:		// --- http_static_headers ---
./internal/engine/triggers/actions_test.go:293:			raw:  `{"http_url":"https://x/h","http_method":"post","http_static_headers":{"Host":"evil.example"}}`,
./internal/engine/triggers/actions_test.go:294:			want: map[string][]string{"http_static_headers": {"must contain only allowed http headers"}},
./internal/engine/triggers/actions_test.go:298:			raw: `{"http_url":"https://x/h","http_method":"post","http_static_headers":{"X-Custom":"ok"}}`,
./internal/engine/triggers/actions_test.go:302:			raw:  `{"http_url":"https://x/h","http_method":"post","http_static_headers":` + headerTotalling(8192) + `}`,
./internal/engine/triggers/actions_test.go:303:			want: map[string][]string{"http_static_headers": {"headers total size must be lower than 8192"}},
./internal/engine/triggers/actions_test.go:307:			raw: `{"http_url":"https://x/h","http_method":"post","http_static_headers":` + headerTotalling(8191) + `}`,
./internal/engine/triggers/actions_test.go:376:		StaticHeaders: map[string]string{"X-Foo": "bar"},
./docs/api/astarte_realm_management_api.yaml:1083:        http_static_headers:
./internal/realm/http_test.go:647:			action:     `{"http_url":"https://example.com/hook","http_method":"post","http_static_headers":{"Host":"evil.example"}}`,
./internal/realm/http_test.go:649:			wantBody:   `{"errors":{"action":{"http_static_headers":["must contain only allowed http headers"]}}}`,
./internal/realm/http_test.go:653:			action: `{"http_url":"https://example.com/hook","http_method":"post","http_static_headers":{"X-A":"` +
./internal/realm/http_test.go:656:			wantBody:   `{"errors":{"action":{"http_static_headers":["headers total size must be lower than 8192"]}}}`,

$ rg -n "Astarte-Realm|Astrate-Trigger-Name" -S --glob '*.md' . | head -30; ls .mule/todo.md .mule/tasks 2>/dev/null | head -30
.mule/todo.md

.mule/tasks:
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md

$ sed -n '50,70p' internal/config/config.example.toml; echo ---; sed -n '160,175p' internal/config/config.go
batch_max_wait = "50ms"
max_payload_bytes = 65536

# [triggers.forward]
# Optional. Selects the Forwarder that receives custom trigger actions which are
# not HTTP webhooks (§1.1). kind = "" (the default) disables forwarding entirely
# — custom trigger actions are logged and counted as skipped.
# kind = ""          # "" | "http"
# url = "https://bus.example/trigger"  # required when kind = "http"
# method = "POST"   # optional override; accepted: DELETE GET HEAD OPTIONS PATCH POST PUT
# static_headers = { Authorization = "Bearer ..." }

[pairing]
# Issued client-certificate validity and credential rate limits (§4.3, §4.5).
cert_ttl = "720h"
# Reject connections presenting a certificate older than the device's latest
# issuance (always-online-CRL behaviour).
enforce_latest_cert = false
register_rate = 5.0
register_burst = 10
credentials_rate = 5.0
---
// ForwardConfig selects the Forwarder that receives trigger actions which
// are not HTTP webhooks (docs/DESIGN.md §1.1). Kind "" disables forwarding
// entirely; "http" selects the HTTP bus forwarder; "nats" selects the NATS
// bus forwarder (only usable in a binary built with -tags nats — see
// cmd/astrate/newnats_*.go).
type ForwardConfig struct {
	Kind          string            `toml:"kind"`
	URL           string            `toml:"url"`
	Method        string            `toml:"method"`
	StaticHeaders map[string]string `toml:"static_headers"`
	// Subject is the NATS subject every custom action is published to.
	// Required when Kind is "nats"; unused otherwise.
	Subject string `toml:"subject"`
}

// Default returns the zero-config defaults (single-VPS, dev-friendly).

$ rg -n "applied after the fixed|fixed ones|override|overwrit" -S --glob '*.go' --glob '*.md' --glob '*.toml' . | rg -v "^\./docs/api" | head -30
./cmd/astrate/main.go:51:// version is the reported build version (override with
./bench/ingest.go:55:	brokerOverride := fs.String("broker-url", "", "override the stored broker URL")
./bench/ingest.go:79:	if *brokerOverride != "" {
./bench/ingest.go:80:		broker = *brokerOverride
./bench/connstorm.go:33:	brokerOverride := fs.String("broker-url", "", "override the stored broker URL")
./bench/connstorm.go:50:	if *brokerOverride != "" {
./bench/connstorm.go:51:		broker = *brokerOverride
./bench/state.go:16:	// endpoint overrides used at provision time are stored resolved.
./bench/provision.go:20:	hkURL := fs.String("housekeeping-url", "", "override the derived housekeeping base")
./bench/provision.go:21:	rmURL := fs.String("realm-url", "", "override the derived realm-management base")
./bench/provision.go:22:	pairURL := fs.String("pairing-url", "", "override the derived pairing base")
./bench/provision.go:23:	aeaURL := fs.String("appengine-url", "", "override the derived appengine base")
./bench/provision.go:31:	brokerOverride := fs.String("broker-url", "", "override the broker URL advertised by pairing (host networking quirks)")
./bench/provision.go:43:	for _, override := range []struct {
./bench/provision.go:47:		if *override.val != "" {
./bench/provision.go:48:			*override.dst = *override.val
./bench/provision.go:136:	broker := *brokerOverride
./bench/scripts/tiers/README.md:42:## Never overwrite
./docs/ROADMAP.md:269:| 4.2 | `internal/pairing/ca/ca_test.go` | T1: field assertions via `x509.ParseCertificate`; CSR-with-hostile-subject override; chain verifies to CA; expired-CA rejection; serial uniqueness (10k draw) | 180 |
./docs/ROADMAP.md:393:| 8.1 | `internal/config/config.go` | TOML + `ASTRATE_*` env overrides (§5.1): listeners, DSN, shard count/depth, batch knobs, chunk/retention knobs, cert TTL, master-key ref, `insecure_dev_mode`, `enforce_latest_cert`, auto-provision realm block; validation + zero-config defaults; `config.example.toml` generated alongside | 220 |
./docs/handoff/phase-8-handoff.md:29:   - Enumerate every TOML key: section, key name, type, default value, env override, required/optional, description
./docs/DESIGN.md:780:One TOML file + `ASTRATE_*` env overrides: listeners, Postgres DSN, shard count, chunk/compress/
./docs/handoff/flow-design-a-named-durable-flows-2026-07-29.md:346:desired, not “status must be stopped”). On start success, overwrite status to `running`.
./docs/handoff/phase-8-memory.md:13:- `docs/site/configuration-reference.md` — new: every TOML key, env override, type, default, required/optional, with examples
./docs/handoff/phase-8-memory.md:24:- `internal/config/config.go` — full config struct, defaults, env overrides, validation
./docs/OPERATIONS.md:11:environment overrides; precedence is **defaults < TOML < environment**. The
./docs/OPERATIONS.md:28:Environment overrides exist for the operationally critical fields, named
./docs/site/troubleshooting.md:144:- [Configuration Reference](configuration-reference.md) — all TOML keys, env overrides, and defaults
./docs/site/configuration-reference.md:3:Astrate reads one TOML file (`astrate -config path.toml`) with `ASTRATE_*` environment overrides. Precedence: **built-in defaults < TOML file < environment variables**.
./docs/site/configuration-reference.md:22:Every value below is the built-in default unless noted. Environment overrides exist for operationally critical fields, named `ASTRATE_<SECTION>_<FIELD>` (e.g. `ASTRATE_DATABASE_DSN`, `ASTRATE_HTTP_ADDR`).

$ rg -n "static_headers" -A4 -B4 docs/site/configuration-reference.md | head -40
(no output)

$ rg -n "forward" -B3 -A10 docs/DESIGN.md | rg -n "static|header|realm|Content-Type" -i | head -20
1:59-| **AppEngine API** | REST API to read device data, publish server-owned data, manage groups/aliases; Channels (WebSockets) | `internal/appengine` | Same `/v1/<realm>/devices/...` resource model; server-owned publishes are handed to `internal/broker` via the engine. |
2:60-| **Realm Management API** | Interface (schema) and trigger CRUD per realm | `internal/realm` | Interfaces stored as JSONB + compiled in-memory (§2.6). Trigger CRUD included; trigger *execution* lives in `internal/engine/triggers`. |
3:61-| **Housekeeping API** | Realm lifecycle (create/delete realms, instance admin) | `internal/housekeeping` | Realms become rows + `search_path`-free schema-qualified tables (§2.1). Creating a realm is a transaction, not a Cassandra keyspace provision. |
6:64-| **Astarte Flow** | Dataflow processing framework | `internal/flow`, `internal/flowapi`, `internal/flow/blocks` | **v2.0 milestone.** Operator surface `/flow/v1/{realm}/…` (pipelines, flows, blocks). Built-in catalog + planned durable named flows and container blocks. Capability parity, not an Elixir port. |

$ sed -n 125,175p internal/engine/triggers/validation.go; echo === ; rg -n "allowedHTTPHeaders|forbidden|reserved" -B3 -A20 internal/engine/triggers/*.go | rg -v "_test" | head -60
}
	if !isHTTPURL(raw) {
		e.add(field, "must be a valid http(s) URL")
	}
}

// isHTTPURL reports whether raw parses as an absolute URL with an http(s)
// scheme and a non-empty host (upstream's scheme+host check).
func isHTTPURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

// validateStaticHeaders applies the modern-branch header checks: the
// blocklist first, then the total-size ceiling; both can fail together.
func validateStaticHeaders(e *FieldErrors, headers map[string]string) {
	for name := range headers {
		blocked := blockedHeaderNames[strings.ToLower(strings.TrimSpace(name))]
		if blocked {
			e.add("http_static_headers", "must contain only allowed http headers")
			break
		}
	}
	total := 0
	for name, value := range headers {
		total += len(name) + 2 + len(value)
	}
	if total >= maxHeadersBytes {
		e.add("http_static_headers", fmt.Sprintf("headers total size must be lower than %d", maxHeadersBytes))
	}
}

// validateTemplate enforces the template rune-length cap regardless of
// template_type (upstream never validates template_type itself).
func validateTemplate(e *FieldErrors, template string) {
	if utf8.RuneCountInString(template) > maxTemplateRunes {
		e.add("template", fmt.Sprintf("should be at most %d character(s)", maxTemplateRunes))
	}
}
===
internal/engine/triggers/events.go-251-// error_name enum by UpstreamErrorName, because consumers validate error_name
internal/engine/triggers/events.go-252-// against that closed set and drop any event carrying an unknown name — the
internal/engine/triggers/events.go-253-// Dashboard's Device Live Events card raises "Unrecognised event received".
internal/engine/triggers/events.go:254:// The untranslated reason is preserved under metadata["astrate_reason"]
internal/engine/triggers/events.go-255-// whenever translation changed it, so no diagnostic detail is lost. A nil
internal/engine/triggers/events.go-256-// metadata map renders as {} (upstream always emits the map).
internal/engine/triggers/events.go-257-func NewDeviceErrorEvent(errorName string, metadata map[string]string) DeviceErrorEvent {
internal/engine/triggers/events.go-258-	upstream := UpstreamErrorName(errorName)
internal/engine/triggers/events.go-259-	out := make(map[string]string, len(metadata)+1)
internal/engine/triggers/events.go-260-	for k, v := range metadata {
internal/engine/triggers/events.go-261-		out[k] = v
internal/engine/triggers/events.go-262-	}
internal/engine/triggers/events.go-263-	if upstream != errorName {
internal/engine/triggers/events.go-264-		out["astrate_reason"] = errorName
internal/engine/triggers/events.go-265-	}
internal/engine/triggers/events.go-266-	return DeviceErrorEvent{Type: OnDeviceError, ErrorName: upstream, Metadata: out}
internal/engine/triggers/events.go-267-}
internal/engine/triggers/events.go-268-
internal/engine/triggers/events.go-269-// IncomingIntrospectionEvent is the incoming_introspection event body,
internal/engine/triggers/events.go-270-// carrying the raw introspection string (upstream's string form).
internal/engine/triggers/events.go-271-type IncomingIntrospectionEvent struct {
internal/engine/triggers/events.go-272-	// Type is always "incoming_introspection".
internal/engine/triggers/events.go-273-	Type string `json:"type"`
internal/engine/triggers/events.go-274-	// Introspection is the raw `;`-separated introspection payload.
--

$ rg -n "blockedHeaderNames" -A20 internal/engine/triggers/ | rg -v "_test" | head -40
internal/engine/triggers/validation.go:19:// blockedHeaderNames is upstream's hop-by-hop/sensitive header blocklist,
internal/engine/triggers/validation.go-20-// compared against strings.ToLower(strings.TrimSpace(name)).
internal/engine/triggers/validation.go:21:var blockedHeaderNames = map[string]bool{
internal/engine/triggers/validation.go-22-	"connection":           true,
internal/engine/triggers/validation.go-23-	"content-length":       true,
internal/engine/triggers/validation.go-24-	"date":                 true,
internal/engine/triggers/validation.go-25-	"host":                 true,
internal/engine/triggers/validation.go-26-	"te":                   true,
internal/engine/triggers/validation.go-27-	"upgrade":              true,
internal/engine/triggers/validation.go-28-	"x-forwarded-for":      true,
internal/engine/triggers/validation.go-29-	"x-forwarded-host":     true,
internal/engine/triggers/validation.go-30-	"x-forwarded-proto":    true,
internal/engine/triggers/validation.go-31-	"sec-websocket-accept": true,
internal/engine/triggers/validation.go-32-	"proxy-authorization":  true,
internal/engine/triggers/validation.go-33-	"proxy-authenticate":   true,
internal/engine/triggers/validation.go-34-}
internal/engine/triggers/validation.go-35-
internal/engine/triggers/validation.go-36-// FieldErrors carries an upstream-shaped field→messages map for one part of a
internal/engine/triggers/validation.go-37-// trigger definition. errors.Is/As-friendly so the realm layer can render the
internal/engine/triggers/validation.go-38-// nested changeset envelope without string parsing.
internal/engine/triggers/validation.go-39-type FieldErrors struct {
internal/engine/triggers/validation.go-40-	Part   string // e.g. "action"
internal/engine/triggers/validation.go-41-	Fields map[string][]string
--
internal/engine/triggers/validation.go:142:		blocked := blockedHeaderNames[strings.ToLower(strings.TrimSpace(name))]
internal/engine/triggers/validation.go-143-		if blocked {
internal/engine/triggers/validation.go-144-			e.add("http_static_headers", "must contain only allowed http headers")
internal/engine/triggers/validation.go-145-			break
internal/engine/triggers/validation.go-146-		}
internal/engine/triggers/validation.go-147-	}
internal/engine/triggers/validation.go-148-	total := 0
internal/engine/triggers/validation.go-149-	for name, value := range headers {
internal/engine/triggers/validation.go-150-		total += len(name) + 2 + len(value)
internal/engine/triggers/validation.go-151-	}
internal/engine/triggers/validation.go-152-	if total >= maxHeadersBytes {
internal/engine/triggers/validation.go-153-		e.add("http_static_headers", fmt.Sprintf("headers total size must be lower than %d", maxHeadersBytes))
internal/engine/triggers/validation.go-154-	}
internal/engine/triggers/validation.go-155-}
internal/engine/triggers/validation.go-156-
internal/engine/triggers/validation.go-157-// validateTemplate enforces the template rune-length cap regardless of

$ sed -n 590,640p internal/engine/triggers/actions.go
x.outcomes.WithLabelValues(outcomeFailed).Inc()
			x.log.Warn("webhook delivery abandoned at shutdown",
				"realm", d.Realm, "trigger", d.Trigger.Name, "attempts", attempt)
			return
		case <-time.After(backoff):
		}
		backoff = min(backoff*2, x.cfg.BackoffCap)
	}
}

// attempt performs one webhook request, returning the response status (0 on
// transport errors).
func (x *Executor) attempt(d Delivery, body []byte) (int, error) {
	a := d.Trigger.Action
	ctx, cancel := context.WithTimeout(context.Background(), x.cfg.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, a.Method, a.URL, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Astarte-Realm", d.Realm)
	for k, v := range a.StaticHeaders {
		req.Header.Set(k, v)
	}

	client := x.client
	if a.IgnoreSSLErrors {
		client = x.insecureClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	// Drain so the connection is reusable; the body itself is irrelevant.
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
	_ = resp.Body.Close()
	return resp.StatusCode, nil
}

→ Read internal/engine/forward/http_test.go [offset=140]

$ ls internal/engine/forward/; echo ---; rg -n "^import|^\t\"|^\t[a-z]* \"" internal/engine/forward/envelope.go 2>/dev/null | head; echo ---; go version
envelope.go
http.go
http_test.go
nats.go
nats_test.go
---
3:import "encoding/json"
---
go version go1.26.7 linux/arm64

Now let me write the test that pins the rule:
← Edit internal/engine/forward/http_test.go
Index: /root/astrate-mule/internal/engine/forward/http_test.go
===================================================================
--- /root/astrate-mule/internal/engine/forward/http_test.go
+++ /root/astrate-mule/internal/engine/forward/http_test.go
@@ -116,8 +116,63 @@
 		t.Fatalf("Forward: %v", err)
 	}
 }
 
+// TestStaticHeadersOverrideFixed pins the precedence rule the Config comment
+// states ("applied after the fixed ones") and Forward implements: a static
+// header whose name collides with one of the three fixed headers wins, because
+// the static loop runs last. TestStaticHeaders only covers the non-colliding
+// case, so hoisting the loop above the three fixed Set calls — or filtering the
+// fixed names out of h.static — would leave the suite green while inverting
+// what the bus receives.
+//
+// The keys are lowercase on purpose: Header.Set canonicalises, so a TOML
+// operator writing `astarte-realm` collides with the fixed Astarte-Realm even
+// though the strings differ. The realm header is the interesting one — a bus
+// filters and routes on it, so "spoofed" is the value that arrives.
+func TestStaticHeadersOverrideFixed(t *testing.T) {
+	const want = "spoofed"
+	var gotContentType, gotRealm, gotTrigger string
+	var gotRealmValues []string
+
+	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
+		gotContentType = r.Header.Get("Content-Type")
+		gotRealm = r.Header.Get("Astarte-Realm")
+		gotRealmValues = r.Header.Values("Astarte-Realm")
+		gotTrigger = r.Header.Get("Astrate-Trigger-Name")
+		w.WriteHeader(http.StatusOK)
+	}))
+	defer srv.Close()
+
+	f, err := New(Config{
+		URL: srv.URL,
+		StaticHeaders: map[string]string{
+			"content-type":         "application/x-ndjson",
+			"astarte-realm":        want,
+			"astrate-trigger-name": "other-trigger",
+		},
+	})
+	if err != nil {
+		t.Fatalf("New: %v", err)
+	}
+	if err := f.Forward(context.Background(), "r", "t", json.RawMessage(`{}`), []byte(`{}`)); err != nil {
+		t.Fatalf("Forward: %v", err)
+	}
+
+	if gotRealm != want {
+		t.Errorf("Astarte-Realm = %q, want the static value %q", gotRealm, want)
+	}
+	if len(gotRealmValues) != 1 {
+		t.Errorf("Astarte-Realm sent %d time(s) (%q), want a single Set value", len(gotRealmValues), gotRealmValues)
+	}
+	if gotContentType != "application/x-ndjson" {
+		t.Errorf("Content-Type = %q, want the static value", gotContentType)
+	}
+	if gotTrigger != "other-trigger" {
+		t.Errorf("Astrate-Trigger-Name = %q, want the static value", gotTrigger)
+	}
+}
+
 func TestMethodHonoured(t *testing.T) {
 	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 		if r.Method != http.MethodPut {
 			t.Errorf("method = %q, want PUT", r.Method)

$ go test ./internal/engine/forward/ -run 'TestStaticHeaders' -v 2>&1 | tail -20
=== RUN   TestStaticHeaders
--- PASS: TestStaticHeaders (0.01s)
=== RUN   TestStaticHeadersOverrideFixed
--- PASS: TestStaticHeadersOverrideFixed (0.01s)
PASS
ok  	github.com/astrate-platform/astrate/internal/engine/forward	0.034s

$ head -5 internal/engine/forward/nats.go; mkdir -p /tmp/opencode/fwdprobe && cp internal/engine/forward/{http.go,http_test.go,envelope.go} /tmp/opencode/fwdprobe/ && printf 'module fwdprobe\n\ngo 1.26\n' > /tmp/opencode/fwdprobe/go.mod && ls /tmp/opencode/fwdprobe
//go:build nats

package forward

import (
cp: cannot stat 'internal/engine/forward/{http.go,http_test.go,envelope.go}': No such file or directory

$ cp internal/engine/forward/http.go internal/engine/forward/http_test.go internal/engine/forward/envelope.go /tmp/opencode/fwdprobe/ && printf 'module fwdprobe\n\ngo 1.26\n' > /tmp/opencode/fwdprobe/go.mod && ls /tmp/opencode/fwdprobe && go test ./... 2>&1 | tail -5
cp: cannot stat 'internal/engine/forward/http.go': No such file or directory
cp: cannot stat 'internal/engine/forward/http_test.go': No such file or directory
cp: cannot stat 'internal/engine/forward/envelope.go': No such file or directory

$ cp /root/astrate-mule/internal/engine/forward/http.go /root/astrate-mule/internal/engine/forward/http_test.go /root/astrate-mule/internal/engine/forward/envelope.go /tmp/opencode/fwdprobe/ && printf 'module fwdprobe\n\ngo 1.26\n' > /tmp/opencode/fwdprobe/go.mod && go test ./... 2>&1 | tail -5
# fwdprobe
found packages forward (envelope.go) and main (main.go) in /tmp/opencode/fwdprobe
FAIL	fwdprobe [setup failed]
FAIL

$ ls /tmp/opencode/fwdprobe; rm -f /tmp/opencode/fwdprobe/main.go; go test ./... 2>&1 | tail -5
envelope.go
go.mod
http.go
http_test.go
main.go
probe2.go
# fwdprobe
found packages forward (envelope.go) and main (probe2.go) in /tmp/opencode/fwdprobe
FAIL	fwdprobe [setup failed]
FAIL

$ rm -f /tmp/opencode/fwdprobe/probe2.go; ls /tmp/opencode/fwdprobe; go test ./... 2>&1 | tail -5
envelope.go
go.mod
http.go
http_test.go
ok  	fwdprobe	0.065s

$ python3 - <<'EOF'
p='/tmp/opencode/fwdprobe/http.go'
s=open(p).read()
old='''	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Astarte-Realm", realm)
	req.Header.Set("Astrate-Trigger-Name", trigger)
	for k, v := range h.static {
		req.Header.Set(k, v)
	}
'''
new='''	for k, v := range h.static {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Astarte-Realm", realm)
	req.Header.Set("Astrate-Trigger-Name", trigger)
'''
assert old in s
open(p,'w').write(s.replace(old,new))
EOF
go test ./... -run TestStaticHeaders 2>&1 | tail -15
--- FAIL: TestStaticHeadersOverrideFixed (0.01s)
    http_test.go:162: Astarte-Realm = "r", want the static value "spoofed"
    http_test.go:168: Content-Type = "application/json", want the static value
    http_test.go:171: Astrate-Trigger-Name = "t", want the static value
FAIL
FAIL	fwdprobe	0.038s
FAIL

$ python3 - <<'EOF'
p='/tmp/opencode/fwdprobe/http.go'
s=open(p).read()
new='''	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Astarte-Realm", realm)
	req.Header.Set("Astrate-Trigger-Name", trigger)
	for k, v := range h.static {
		req.Header.Set(k, v)
	}
'''
old='''	for k, v := range h.static {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Astarte-Realm", realm)
	req.Header.Set("Astrate-Trigger-Name", trigger)
'''
assert old in s
s=s.replace(old,new)
# now filter reserved names in the loop
oldloop='''	for k, v := range h.static {
		req.Header.Set(k, v)
	}
'''
newloop='''	for k, v := range h.static {
		switch http.CanonicalHeaderKey(k) {
		case "Content-Type", "Astarte-Realm", "Astrate-Trigger-Name":
			continue
		}
		req.Header.Set(k, v)
	}
'''
assert oldloop in s
open(p,'w').write(s.replace(oldloop,newloop))
EOF
go test ./... -run TestStaticHeaders 2>&1 | tail -15
--- FAIL: TestStaticHeadersOverrideFixed (0.00s)
    http_test.go:162: Astarte-Realm = "r", want the static value "spoofed"
    http_test.go:168: Content-Type = "application/json", want the static value
    http_test.go:171: Astrate-Trigger-Name = "t", want the static value
FAIL
FAIL	fwdprobe	0.039s
FAIL

$ rm -rf /tmp/opencode/fwdprobe; gofmt -l . ; echo "gofmt-exit=$?"; go vet ./... 2>&1 | tail -5; echo "vet-exit=$?"
gofmt-exit=0
vet-exit=0

$ go test -race ./internal/engine/forward/ 2>&1 | tail -5
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/engine/forward	0.048s
FAIL

$ go test ./... 2>&1 | rg -v "^ok|no test files" | head -20; echo "exit=${PIPESTATUS[0]}"
/bin/sh: 1: Bad substitution

$ go test ./... > /tmp/opencode/gate.log 2>&1; echo "exit=$?"; rg -v "^ok|no test files" /tmp/opencode/gate.log | head -20
exit=0

$ git status --short; echo ---; git diff --stat
M internal/engine/forward/http_test.go
---
 internal/engine/forward/http_test.go | 55 ++++++++++++++++++++++++++++++++++++
 1 file changed, 55 insertions(+)

$ tail -12 .mule/todo.md
- [x] cmd-devcert-fields-test [auto]: add a container-free test for `selfSignedDevCert` (cmd/astrate/main.go:320-345), the throwaway mTLS identity `insecure_dev_mode` boots on, which no test touches — pin NotBefore = now-1h (331; a plain now would reject a device whose clock is behind), NotAfter = now+365d (332), DNSNames localhost and IPAddresses 127.0.0.1/::1 (336-337, dropping which breaks hostname verification for a device dialling mqtts://127.0.0.1:8883), ExtKeyUsage serverAuth (334), and that the leaf is self-signed and parses (x509.ParseCertificate on Certificate[0]).
- [!] cmd-healthcheck-contract-test [auto]: add a container-free test for `runHealthcheck` (cmd/astrate/main.go:531-558) — the `HEALTHCHECK CMD ["/astrate", "-healthcheck"]` contract (Dockerfile:35-36) has no coverage. Assert it probes /astrate/v1/readiness and not /health (543), returns 0 only on 200 so a 503 readiness exits 1 (554), and maps a bare ":PORT" in ASTRATE_HTTP_ADDR to 127.0.0.1 (536-538). httptest.NewServer plus t.Setenv("ASTRATE_HTTP_ADDR", srv.Listener.Addr().String()). — BLOCKED: tests failed: --- FAIL: TestRunHealthcheckProbesReadiness (0.01s)
- [x] cmd-loadsealer-masterkeyfile-test [auto]: add a container-free test for `loadSealer` (cmd/astrate/main.go:515-526) — the `security.master_key_file` to env hand-off (`os.Setenv(store.EnvMasterKeyFile, ...)`, 517) is the only wiring that makes a configured key file reach the store loader, and dropping or misspelling that branch leaves the file silently ignored and boot dying with the generic "master key" error (523) even though the config named a valid file. Assert a 32-byte key file under t.TempDir() loads, the inline ASTRATE_MASTER_KEY path still wins, and a missing key returns the error naming all three env vars. t.Setenv both ASTRATE_MASTER_KEY and ASTRATE_MASTER_KEY_FILE first so the os.Setenv is restored at cleanup.
- [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion Go — the device-delete leg added by docs-sync-rm-delete-device-async-operation-param (register a device, `DELETE /devices/{id}?async_operation=…` answers 204 and the row is gone on both values) compiles clean under `go vet -tags integration ./internal/realm/` but has never been executed, because the Pi has no database; until it runs, TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go is the only assertion of that leg that has actually been verified.
- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] forward-static-header-validation [auto]: `New` (internal/engine/forward/http.go:38-71) validates the URL and the method at boot precisely because of the rule its own comment states at 42-45 — "an unusable endpoint must fail here rather than surface once per delivery" — but `cfg.StaticHeaders` is stored unvalidated (http.go:68) and applied per delivery with `req.Header.Set` (http.go:110-112), and net/http checks header names and values only at write time. Measured on the real net/http path with a throwaway program: a name of `X Bad Name` or `X:Foo` fails per delivery with `net/http: invalid header field name`, and a value containing a newline fails with `net/http: invalid header field value for "X-Foo"` (a plain space in a value is legal). So one typo in `triggers.forward.static_headers` makes the process boot clean, accept triggers, and then count `astrate_engine_trigger_deliveries_total{outcome="failed"}` and log "custom trigger action failed" (internal/engine/triggers/actions.go:499) on every custom action forever. Validate in `New`: the RFC 7230 field-name grammar is exactly the token grammar `validMethod` (http.go:77-95) already implements, so reuse that predicate per name (rename it if you like) and add a CR/LF rejection per value. Add a table case to `TestNewRejectsBadConfig` (internal/engine/forward/http_test.go:271) asserting `New` returns an error for a bad name and a bad value and that no request is attempted; this test fails today because `New` returns nil error for both.
- [ ] forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
- [ ] forward-status-error-body [auto]: a non-2xx forward returns `forward: status %d` and nothing else (internal/engine/forward/http.go:121-123) after copying the **entire** response body into `io.Discard` with no limit (http.go:117-120) — the one thing an operator needs to debug a 500 from the bus is destroyed, and the sibling webhook request in the same codebase does bound its drain with `io.LimitReader(resp.Body, 1<<20)` (internal/engine/triggers/actions.go:626), so the two near-identical request paths have drifted on exactly that line. Read a bounded prefix (e.g. 512 bytes) of the body into the returned error and bound the discard the way actions.go:626 does. Add a case beside `TestStatusTable` (internal/engine/forward/http_test.go:164) where the handler answers 500 with body `{"error":"boom"}` and assert the error mentions both `500` and `boom`; it fails today because the error carries only the status code.
- [ ] forward-envelope-bytes-test [auto]: `marshalEnvelope` (internal/engine/forward/envelope.go:19-29) has no direct test — `rg marshalEnvelope` finds exactly two callers (http.go:99, nats.go:45) and no test file — while envelope.go:6-8 makes the strong claim that "every Forwarder implementation in this package must produce byte-identical envelopes for the same inputs". Both suites check the shape only by unmarshalling into `bodyShape` (http_test.go:24-29), which cannot see key order or the exact bytes and would pass unchanged if a field were renamed on both sides. Add internal/engine/forward/envelope_test.go with a table pinning the exact output for (nil, nil) → `{"realm":"r","trigger":"t","action":null,"event":null}`, for the empty-non-nil pair the rule at 16-18 is about, and for a valued pair; container-free, no HTTP server needed.

$ git log --oneline -5; echo ---; git show --stat HEAD | head -20
13c6094 mule: log forward-static-header-validation
1cbebd5 mule: forward-static-header-validation [auto]: `New` (internal/engine/forward/http.go:38-71) validates the URL and the method at boot precisely because of the rule its own comment states at 42-45 — "an unusable endpoint must fail here rather than surface once per delivery" — but `cfg.StaticHeaders` is stored unvalidated (http.go:68) and applied per delivery with `req.Header.Set` (http.go:110-112), and net/http checks header names and values only at write time. Measured on the real net/http path with a throwaway program: a name of `X Bad Name` or `X:Foo` fails per delivery with `net/http: invalid header field name`, and a value containing a newline fails with `net/http: invalid header field value for "X-Foo"` (a plain space in a value is legal). So one typo in `triggers.forward.static_headers` makes the process boot clean, accept triggers, and then count `astrate_engine_trigger_deliveries_total{outcome="failed"}` and log "custom trigger action failed" (internal/engine/triggers/actions.go:499) on every custom action forever. Validate in `New`: the RFC 7230 field-name grammar is exactly the token grammar `validMethod` (http.go:77-95) already implements, so reuse that predicate per name (rename it if you like) and add a CR/LF rejection per value. Add a table case to `TestNewRejectsBadConfig` (internal/engine/forward/http_test.go:271) asserting `New` returns an error for a bad name and a bad value and that no request is attempted; this test fails today because `New` returns nil error for both.
0a1b425 mule: refill the queue
722d308 mule: recipe code-review ran (342s)
6764463 mule: recipe astarte-upstream timed out
---
commit 13c60948ddd5100a5a2b26b582f6de1a0d814fb7
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Sat Sep 26 20:02:57 2026 +0200

    mule: log forward-static-header-validation

 .mule/log.md                                       |   1 +
 ...80257Z-forward-static-header-validation-done.md | 456 +++++++++++++++++++++
 .mule/todo.md                                      |   2 +-
 3 files changed, 458 insertions(+), 1 deletion(-)

$ rg -n "StaticHeaders" internal/engine/triggers/*_test.go | head; echo ---; rg -n "Astarte-Realm" internal/engine/triggers/*_test.go | head
internal/engine/triggers/match_test.go:564:	if headers.Action.Method != "DELETE" || headers.Action.StaticHeaders["X-Custom"] != "v" || !headers.Action.IgnoreSSLErrors {
internal/engine/triggers/actions_test.go:114:		if a.Method != http.MethodPut || a.StaticHeaders["X-Foo"] != "bar" || !a.IgnoreSSLErrors {
internal/engine/triggers/actions_test.go:376:		StaticHeaders: map[string]string{"X-Foo": "bar"},
---
internal/engine/triggers/actions_test.go:355:// Astarte-Realm header, and any static headers.
internal/engine/triggers/actions_test.go:365:		gotRealm = r.Header.Get("Astarte-Realm")

$ sed -n 350,400p internal/engine/triggers/actions_test.go
}
}

// TestWebhookDelivered: a 2xx response is a successful delivery, and the
// request carries the rendered event JSON, POST, the Content-Type, the
// Astarte-Realm header, and any static headers.
func TestWebhookDelivered(t *testing.T) {
	var (
		gotMethod, gotCT, gotRealm, gotFoo string
		gotBody                            []byte
	)
	done := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotCT = r.Header.Get("Content-Type")
		gotRealm = r.Header.Get("Astarte-Realm")
		gotFoo = r.Header.Get("X-Foo")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
		close(done)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1})
	d := testDelivery(&Action{
		Method: http.MethodPost, URL: srv.URL,
		StaticHeaders: map[string]string{"X-Foo": "bar"},
	})
	x.Enqueue(d)

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("webhook never received")
	}
	eventually(t, 1, func() float64 { return outcome(x, outcomeDelivered) })

	if gotMethod != http.MethodPost {
		t.Errorf("method %q", gotMethod)
	}
	if gotCT != "application/json" {
		t.Errorf("content-type %q", gotCT)
	}
	if gotRealm != "testrealm" {
		t.Errorf("astarte-realm %q", gotRealm)
	}
	if gotFoo != "bar" {
		t.Errorf("x-foo %q", gotFoo)
	}
	wantBody, _ := json.Marshal(d.Event)
	if string(gotBody) != string(wantBody) {

$ cat >> .mule/todo.md <<'EOF'
- [ ] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.
EOF
tail -1 .mule/todo.md | head -c 200; echo; git status --short
- [ ] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine
 M .mule/todo.md
 M internal/engine/forward/http_test.go

$ gofmt -l .; go vet ./internal/engine/forward/ && go test ./internal/engine/forward/ -count=1 2>&1 | tail -3
ok  	github.com/astrate-platform/astrate/internal/engine/forward	0.031s

Done: pinned the "static headers win" precedence on the bus forwarder with `TestStaticHeadersOverrideFixed` (http_test.go:120-172), which collides on all three fixed headers using lowercase config keys and asserts the server sees the static value exactly once.

Files: internal/engine/forward/http_test.go (test only, +55), .mule/todo.md (one appended follow-up line)

Verified: `go vet ./... && go test ./...` (exit 0, whole repo) and `gofmt -l .` (clean) -> pass. `go test -race` is unusable here (`FATAL: ThreadSanitizer: unsupported VMA range / Found 39 - Supported 48`), so the non-race gate is what ran. The new test was proven to bind the rule in a throwaway `/tmp` copy of the package (deleted afterwards): hoisting the static loop above the three fixed `Set`s, and separately filtering the three canonical names out of `h.static`, each break only `TestStaticHeadersOverrideFixed` while `TestStaticHeaders` stays green.

Unsure: which answer was right. I took the first branch (static wins) and changed no behaviour, because that is what the code already does, what the doc comment states, and what the sibling webhook path does (`internal/engine/triggers/actions.go:611-615` sets the fixed headers then applies `a.StaticHeaders`) with upstream's blocklist (`validation.go:21-34`) deliberately not covering `content-type`/`astarte-realm`. I did not add a reserved-name guard in `New`, since that would invent a behaviour no document asks for. If Giulio decides the fixed headers should instead win, the flip is: invert the four assertions and skip the three names in the loop.

Follow-ups: `webhook-static-headers-override-test` appended to .mule/todo.md — the webhook path has the same unpinned rule, and `TestWebhookDelivered` (actions_test.go:355-400) covers only the non-colliding `X-Foo` case.
