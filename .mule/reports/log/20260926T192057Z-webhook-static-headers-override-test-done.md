slug: webhook-static-headers-override-test
verdict: done
at:  f9db0a6
ran: 2026-09-26T19:20:57Z on DietPi in 178s


> build · big-pickle

I'll read the two files first.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/engine/triggers/actions.go [offset=570, limit=80]
→ Read internal/engine/triggers/validation.go [offset=1, limit=50]

$ rg -n "StaticHeaders" --glob '*.go' -C 3 . | head -80
./cmd/astrate/main.go-262-		h, err := forward.New(forward.Config{
./cmd/astrate/main.go-263-			URL:           f.URL,
./cmd/astrate/main.go-264-			Method:        f.Method,
./cmd/astrate/main.go:265:			StaticHeaders: f.StaticHeaders,
./cmd/astrate/main.go-266-		})
./cmd/astrate/main.go-267-		if err != nil {
./cmd/astrate/main.go-268-			return nil, err
--
./internal/config/config.go-166-	Kind          string            `toml:"kind"`
./internal/config/config.go-167-	URL           string            `toml:"url"`
./internal/config/config.go-168-	Method        string            `toml:"method"`
./internal/config/config.go:169:	StaticHeaders map[string]string `toml:"static_headers"`
./internal/config/config.go-170-	// Subject is the NATS subject every custom action is published to.
./internal/config/config.go-171-	// Required when Kind is "nats"; unused otherwise.
./internal/config/config.go-172-	Subject string `toml:"subject"`
--
./internal/config/config_test.go-279-	if f.Method != "POST" {
./internal/config/config_test.go-280-		t.Errorf("Method = %q, want POST", f.Method)
./internal/config/config_test.go-281-	}
./internal/config/config_test.go:282:	if f.StaticHeaders["Authorization"] != "Bearer tok" || f.StaticHeaders["X-Forward"] != "yes" {
./internal/config/config_test.go:283:		t.Errorf("StaticHeaders = %v", f.StaticHeaders)
./internal/config/config_test.go-284-	}
./internal/config/config_test.go-285-}
./internal/config/config_test.go-286-
--
./internal/engine/triggers/match_test.go-561-			"http_static_headers": {"X-Custom": "v"}, "ignore_ssl_errors": true},
./internal/engine/triggers/match_test.go-562-		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
./internal/engine/triggers/match_test.go-563-	}`)
./internal/engine/triggers/match_test.go:564:	if headers.Action.Method != "DELETE" || headers.Action.StaticHeaders["X-Custom"] != "v" || !headers.Action.IgnoreSSLErrors {
./internal/engine/triggers/match_test.go-565-		t.Errorf("action attributes: %+v", headers.Action)
./internal/engine/triggers/match_test.go-566-	}
./internal/engine/triggers/match_test.go-567-}
--
./internal/engine/triggers/actions.go-27-	Method string
./internal/engine/triggers/actions.go-28-	// URL is the webhook endpoint.
./internal/engine/triggers/actions.go-29-	URL string
./internal/engine/triggers/actions.go:30:	// StaticHeaders are the action's extra request headers.
./internal/engine/triggers/actions.go:31:	StaticHeaders map[string]string
./internal/engine/triggers/actions.go-32-	// IgnoreSSLErrors disables server-certificate verification for this
./internal/engine/triggers/actions.go-33-	// action's requests (upstream "ignore_ssl_errors").
./internal/engine/triggers/actions.go-34-	IgnoreSSLErrors bool
--
./internal/engine/triggers/actions.go-49-	HTTPURL           string            `json:"http_url"`
./internal/engine/triggers/actions.go-50-	HTTPMethod        string            `json:"http_method"`
./internal/engine/triggers/actions.go-51-	HTTPPostURL       string            `json:"http_post_url"` // pre-1.1 legacy: implies POST
./internal/engine/triggers/actions.go:52:	HTTPStaticHeaders map[string]string `json:"http_static_headers"`
./internal/engine/triggers/actions.go-53-	IgnoreSSLErrors   bool              `json:"ignore_ssl_errors"`
./internal/engine/triggers/actions.go-54-	Template          string            `json:"template"`
./internal/engine/triggers/actions.go-55-	TemplateType      string            `json:"template_type"`
--
./internal/engine/triggers/actions.go-91-		if h.HTTPMethod != "" {
./internal/engine/triggers/actions.go-92-			fe.add("http_method", "must be blank")
./internal/engine/triggers/actions.go-93-		}
./internal/engine/triggers/actions.go:94:		if len(h.HTTPStaticHeaders) > 0 {
./internal/engine/triggers/actions.go-95-			fe.add("http_static_headers", "must be blank")
./internal/engine/triggers/actions.go-96-		}
./internal/engine/triggers/actions.go-97-		validateURLField(fe, "http_post_url", h.HTTPPostURL)
--
./internal/engine/triggers/actions.go-101-		}
./internal/engine/triggers/actions.go-102-		return &Action{
./internal/engine/triggers/actions.go-103-			Method: http.MethodPost, URL: h.HTTPPostURL,
./internal/engine/triggers/actions.go:104:			StaticHeaders: h.HTTPStaticHeaders, IgnoreSSLErrors: h.IgnoreSSLErrors,
./internal/engine/triggers/actions.go-105-			Template: h.Template, TemplateType: h.TemplateType,
./internal/engine/triggers/actions.go-106-		}, unsupported, nil
./internal/engine/triggers/actions.go-107-	case h.HTTPURL != "":
--
./internal/engine/triggers/actions.go-118-			fe.add("http_post_url", "must be blank")
./internal/engine/triggers/actions.go-119-		}
./internal/engine/triggers/actions.go-120-		validateURLField(fe, "http_url", h.HTTPURL)
./internal/engine/triggers/actions.go:121:		validateStaticHeaders(fe, h.HTTPStaticHeaders)
./internal/engine/triggers/actions.go-122-		validateTemplate(fe, h.Template)
./internal/engine/triggers/actions.go-123-		if len(fe.Fields) > 0 {
./internal/engine/triggers/actions.go-124-			return nil, nil, fe
--
./internal/engine/triggers/actions.go-127-			// Methods arrive lowercase on the wire (upstream convention) but
./internal/engine/triggers/actions.go-128-			// HTTP methods are case-sensitive uppercase tokens.
./internal/engine/triggers/actions.go-129-			Method: strings.ToUpper(h.HTTPMethod), URL: h.HTTPURL,
./internal/engine/triggers/actions.go:130:			StaticHeaders: h.HTTPStaticHeaders, IgnoreSSLErrors: h.IgnoreSSLErrors,
./internal/engine/triggers/actions.go-131-			Template: h.Template, TemplateType: h.TemplateType,
./internal/engine/triggers/actions.go-132-		}, unsupported, nil

$ rg -n "forward-static-headers-override|static-headers" .mule/ | head -20
.mule/log.md:135:| 2026-09-26 | forward-static-headers-override-test | done | 263s | 3a0f153 |
.mule/reviews/forward-2026-09-26.md:49:**2. The static-headers-win rule is documented and completely unasserted.** http.go:22 says
.mule/reviews/forward-2026-09-26.md:90:- `forward-static-headers-override-test` — pin the documented apply-after-the-fixed-headers rule
.mule/todo.md:232:- [x] forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
.mule/todo.md:235:- [ ] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.
.mule/reports/log/20260926T174259Z-recipe-code-review-proposed.md:919:- [ ] forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
.mule/reports/log/20260926T174259Z-recipe-code-review-proposed.md:965:Follow-ups: forward-static-header-validation, forward-static-headers-override-test, forward-status-error-body, forward-envelope-bytes-test
.mule/reports/log/20260926T181916Z-forward-static-headers-override-test-done.md:1:slug: forward-static-headers-override-test
.mule/reports/log/20260926T181916Z-forward-static-headers-override-test-done.md:535:- [ ] forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
.mule/reports/log/20260926T181916Z-forward-static-headers-override-test-done.md:619:- [ ] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.
.mule/reports/log/20260926T181916Z-forward-static-headers-override-test-done.md:622:- [ ] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine
.mule/reports/log/20260926T181916Z-forward-static-headers-override-test-done.md:637:Follow-ups: `webhook-static-headers-override-test` appended to .mule/todo.md — the webhook path has the same unpinned rule, and `TestWebhookDelivered` (actions_test.go:355-400) covers only the non-colliding `X-Foo` case.

$ rg -n "static wins|static headers win|precedence" -i /root/astrate-mule/.mule/reports/log/20260926T181916Z-forward-static-headers-override-test-done.md | head -30
134:./docs/OPERATIONS.md:11:environment overrides; precedence is **defaults < TOML < environment**. The
137:./docs/site/configuration-reference.md:3:Astrate reads one TOML file (`astrate -config path.toml`) with `ASTRATE_*` environment overrides. Precedence: **built-in defaults < TOML file < environment variables**.
324:+// TestStaticHeadersOverrideFixed pins the precedence rule the Config comment
619:- [ ] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.
622:- [ ] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine
629:Done: pinned the "static headers win" precedence on the bus forwarder with `TestStaticHeadersOverrideFixed` (http_test.go:120-172), which collides on all three fixed headers using lowercase config keys and asserts the server sees the static value exactly once.
635:Unsure: which answer was right. I took the first branch (static wins) and changed no behaviour, because that is what the code already does, what the doc comment states, and what the sibling webhook path does (`internal/engine/triggers/actions.go:611-615` sets the fixed headers then applies `a.StaticHeaders`) with upstream's blocklist (`validation.go:21-34`) deliberately not covering `content-type`/`astarte-realm`. I did not add a reserved-name guard in `New`, since that would invent a behaviour no document asks for. If Giulio decides the fixed headers should instead win, the flip is: invert the four assertions and skip the three names in the loop.

→ Read internal/engine/triggers/actions_test.go [offset=300, limit=130]
→ Read internal/engine/forward/http_test.go [limit=80, offset=110]
→ Read internal/engine/forward/http.go [limit=40, offset=15]

$ rg -n "func newTestExecutor|func testDelivery|func eventually|func outcome" -A 30 /root/astrate-mule/internal/engine/triggers/actions_test.go | head -100
26:func newTestExecutor(t *testing.T, cfg ExecutorConfig) *Executor {
27-	t.Helper()
28-	cfg.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
29-	if cfg.BackoffStart == 0 {
30-		cfg.BackoffStart = time.Millisecond
31-	}
32-	if cfg.BackoffCap == 0 {
33-		cfg.BackoffCap = 5 * time.Millisecond
34-	}
35-	if cfg.Registerer == nil {
36-		cfg.Registerer = prometheus.NewRegistry()
37-	}
38-	x := NewExecutor(cfg)
39-	t.Cleanup(func() {
40-		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
41-		defer cancel()
42-		if err := x.Close(ctx); err != nil {
43-			t.Errorf("executor close: %v", err)
44-		}
45-	})
46-	return x
47-}
48-
49-// testDelivery builds a delivery whose action is the given action.
50:func testDelivery(action *Action) Delivery {
51-	return Delivery{
52-		Realm:   "testrealm",
53-		Trigger: &Trigger{Name: "hook", Action: action},
54-		Event: SimpleEvent{
55-			Timestamp:   time.Date(2026, 6, 12, 10, 0, 0, 0, time.UTC),
56-			DeviceID:    "f0VMRgIBAQAAAAAAAAAAAA",
57-			TriggerName: "hook",
58-			Event:       NewIncomingDataEvent("com.ex.Sensors", "/v", 1.0),
59-		},
60-	}
61-}
62-
63-// outcome reads one outcome counter.
64:func outcome(x *Executor, label string) float64 {
65-	return promtest.ToFloat64(x.outcomes.WithLabelValues(label))
66-}
67-
68-// eventually polls get until it equals want or the deadline passes.
69:func eventually(t *testing.T, want float64, get func() float64) {
70-	t.Helper()
71-	deadline := time.Now().Add(3 * time.Second)
72-	for time.Now().Before(deadline) {
73-		if get() == want {
74-			return
75-		}
76-		time.Sleep(time.Millisecond)
77-	}
78-	t.Fatalf("value did not reach %v (last %v)", want, get())
79-}
80-
81-// TestParseAction validates the action-object parser (docs/ROADMAP.md §7.2
82-// file 6.12): http_url+method, the legacy http_post_url, custom (non-HTTP)
83-// actions routed to the Forwarder, and the rejection / unsupported paths.
84-func TestParseAction(t *testing.T) {
85-	t.Run("http_url+method lowercases on the wire, uppercases for net/http", func(t *testing.T) {
86-		a, unsupported, err := parseAction([]byte(`{"http_url":"https://x/h","http_method":"post"}`))
87-		if err != nil {
88-			t.Fatal(err)
89-		}
90-		if a.Method != http.MethodPost || a.URL != "https://x/h" {
91-			t.Errorf("got %q %q", a.Method, a.URL)
92-		}
93-		if a.Custom != nil || len(unsupported) != 0 {
94-			t.Errorf("custom=%s unsupported=%v", a.Custom, unsupported)
95-		}
96-	})
97-
98-	t.Run("legacy http_post_url implies POST", func(t *testing.T) {
99-		a, _, err := parseAction([]byte(`{"http_post_url":"https://y/legacy"}`))

→ Read internal/engine/triggers/validation.go [offset=95, limit=50]
→ Read internal/engine/triggers/actions_test.go [limit=26]

← Edit internal/engine/triggers/actions_test.go
Index: /root/astrate-mule/internal/engine/triggers/actions_test.go
===================================================================
--- /root/astrate-mule/internal/engine/triggers/actions_test.go
+++ /root/astrate-mule/internal/engine/triggers/actions_test.go
@@ -401,8 +401,68 @@
 		t.Errorf("body = %s, want %s", gotBody, wantBody)
 	}
 }
 
+// TestWebhookStaticHeadersOverrideFixed pins the precedence rule attempt
+// implements on the webhook path: the fixed Content-Type and Astarte-Realm are
+// set first and a.StaticHeaders second, so an action's http_static_headers wins
+// on a colliding name. TestWebhookDelivered only covers the non-colliding X-Foo
+// case, so hoisting the static loop above the two fixed Set calls — or filtering
+// the two names out of the map — would leave the suite green while inverting
+// what the endpoint receives.
+//
+// "Static wins" is the deliberate parity call, for the reasons the bus
+// forwarder already settled on (TestStaticHeadersOverrideFixed): the fixed
+// Content-Type is only a default for the JSON envelope, and upstream's own
+// blocklist (blockedHeaderNames) rejects content-length, host and the
+// hop-by-hop set but deliberately not content-type or astarte-realm, so naming
+// them in a trigger definition is allowed and upstream would deliver the
+// value. The action is built here by parseAction, so the test also pins that
+// the collision survives validation instead of relying on a hand-built Action.
+//
+// The header names are lowercase on purpose: that is what a trigger
+// definition carries, and Header.Set canonicalises them into a collision with
+// the fixed names.
+func TestWebhookStaticHeadersOverrideFixed(t *testing.T) {
+	const wantRealm = "spoofed"
+	var gotCT, gotRealm string
+	var gotRealmValues []string
+
+	done := make(chan struct{})
+	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
+		gotCT = r.Header.Get("Content-Type")
+		gotRealm = r.Header.Get("Astarte-Realm")
+		gotRealmValues = r.Header.Values("Astarte-Realm")
+		w.WriteHeader(http.StatusNoContent)
+		close(done)
+	}))
+	defer srv.Close()
+
+	a, _, err := parseAction([]byte(`{"http_url":"` + srv.URL +
+		`","http_method":"post","http_static_headers":{"content-type":"application/x-ndjson","astarte-realm":"spoofed"}}`))
+	if err != nil {
+		t.Fatalf("parseAction rejected a colliding header: %v", err)
+	}
+	x := newTestExecutor(t, ExecutorConfig{Workers: 1})
+	x.Enqueue(testDelivery(a))
+
+	select {
+	case <-done:
+	case <-time.After(5 * time.Second):
+		t.Fatal("webhook never received")
+	}
+
+	if gotRealm != wantRealm {
+		t.Errorf("Astarte-Realm = %q, want the static value %q", gotRealm, wantRealm)
+	}
+	if len(gotRealmValues) != 1 {
+		t.Errorf("Astarte-Realm sent %d time(s) (%q), want a single Set value", len(gotRealmValues), gotRealmValues)
+	}
+	if gotCT != "application/x-ndjson" {
+		t.Errorf("Content-Type = %q, want the static value", gotCT)
+	}
+}
+
 // TestWebhookMustacheTemplateRendered: a mustache action renders the body
 // from realm/device/trigger/event fields instead of the default envelope.
 func TestWebhookMustacheTemplateRendered(t *testing.T) {
 	var gotBody []byte

$ go test -race -run 'TestWebhook' ./internal/engine/triggers/ 2>&1 | tail -20; gofmt -l . ; go vet ./... 2>&1 | tail -5
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/engine/triggers	0.022s
FAIL

$ go test -run 'TestWebhook' -v ./internal/engine/triggers/ 2>&1 | tail -25
=== RUN   TestWebhookDelivered
--- PASS: TestWebhookDelivered (0.01s)
=== RUN   TestWebhookStaticHeadersOverrideFixed
--- PASS: TestWebhookStaticHeadersOverrideFixed (0.00s)
=== RUN   TestWebhookMustacheTemplateRendered
--- PASS: TestWebhookMustacheTemplateRendered (0.00s)
=== RUN   TestWebhookMustacheTemplateMalformedFallsBack
--- PASS: TestWebhookMustacheTemplateMalformedFallsBack (0.00s)
=== RUN   TestWebhookRetriesThenSucceeds
--- PASS: TestWebhookRetriesThenSucceeds (0.01s)
=== RUN   TestWebhook4xxIsPermanent
--- PASS: TestWebhook4xxIsPermanent (0.02s)
=== RUN   TestWebhookFailsAfterMaxAttempts
--- PASS: TestWebhookFailsAfterMaxAttempts (0.01s)
PASS
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	0.245s

Now the inversion check — temporarily hoist the static loop above the fixed `Set`s and confirm the new test fails.
← Edit internal/engine/triggers/actions.go
Index: /root/astrate-mule/internal/engine/triggers/actions.go
===================================================================
--- /root/astrate-mule/internal/engine/triggers/actions.go
+++ /root/astrate-mule/internal/engine/triggers/actions.go
@@ -607,13 +607,13 @@
 req, err := http.NewRequestWithContext(ctx, a.Method, a.URL, bytes.NewReader(body))
 if err != nil {
 	return 0, err
 }
-req.Header.Set("Content-Type", "application/json")
-req.Header.Set("Astarte-Realm", d.Realm)
 for k, v := range a.StaticHeaders {
 	req.Header.Set(k, v)
 }
+req.Header.Set("Content-Type", "application/json")
+req.Header.Set("Astarte-Realm", d.Realm)
 
 client := x.client
 if a.IgnoreSSLErrors {
 	client = x.insecureClient

$ go test -run 'TestWebhookStaticHeadersOverrideFixed' ./internal/engine/triggers/ 2>&1 | tail -12
--- FAIL: TestWebhookStaticHeadersOverrideFixed (0.01s)
    actions_test.go:455: Astarte-Realm = "testrealm", want the static value "spoofed"
    actions_test.go:461: Content-Type = "application/json", want the static value
FAIL
FAIL	github.com/astrate-platform/astrate/internal/engine/triggers	0.149s
FAIL

← Edit internal/engine/triggers/actions.go
Index: /root/astrate-mule/internal/engine/triggers/actions.go
===================================================================
--- /root/astrate-mule/internal/engine/triggers/actions.go
+++ /root/astrate-mule/internal/engine/triggers/actions.go
@@ -607,13 +607,13 @@
 req, err := http.NewRequestWithContext(ctx, a.Method, a.URL, bytes.NewReader(body))
 if err != nil {
 	return 0, err
 }
+req.Header.Set("Content-Type", "application/json")
+req.Header.Set("Astarte-Realm", d.Realm)
 for k, v := range a.StaticHeaders {
 	req.Header.Set(k, v)
 }
-req.Header.Set("Content-Type", "application/json")
-req.Header.Set("Astarte-Realm", d.Realm)
 
 client := x.client
 if a.IgnoreSSLErrors {
 	client = x.insecureClient

$ git diff --stat; go vet ./... && go test ./... 2>&1 | tail -25; gofmt -l .
internal/engine/triggers/actions_test.go | 60 ++++++++++++++++++++++++++++++++
 1 file changed, 60 insertions(+)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	1.948s
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

```
Done: pinned the "static headers win" precedence on the webhook path with `TestWebhookStaticHeadersOverrideFixed` (actions_test.go:405-465), which builds its Action through `parseAction` with lowercase `content-type`/`astarte-realm` in `http_static_headers` — so the test also proves the collision survives validation — and asserts the server sees the static values, Astarte-Realm exactly once; no behaviour changed.
Files: internal/engine/triggers/actions_test.go
Verified: `go vet ./... && go test ./...` -> pass; `gofmt -l .` -> clean; `go test -race` unusable here (ThreadSanitizer: unsupported VMA range, 39-bit — same kernel limitation as the Pi, so the new test is unverified for races). Inversion check: hoisting the static loop above the two fixed `Set`s in `attempt` makes the new test fail on both headers (restored afterwards; `git diff --stat` shows only the test file).
Unsure: the parity call itself — I took "static wins" (the code's current behaviour, the bus forwarder's already-pinned rule, and upstream's blocklist not covering these two names) rather than "fixed wins", and stated that reasoning in the test comment. Flip is: invert three assertions and skip the two names in the loop.
Follow-ups: none
```
