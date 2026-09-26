slug: forward-static-header-validation
verdict: done
at:  1cbebd5
ran: 2026-09-26T18:02:57Z on DietPi in 874s


> build · big-pickle

I'll start by reading the two files.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/engine/forward/http.go

$ wc -l internal/engine/forward/*.go && rg -n "validMethod|TestNewRejectsBadConfig" internal/engine/forward/
29 internal/engine/forward/envelope.go
  125 internal/engine/forward/http.go
  310 internal/engine/forward/http_test.go
   53 internal/engine/forward/nats.go
  196 internal/engine/forward/nats_test.go
  713 total
internal/engine/forward/http.go:57:		if !validMethod(m) {
internal/engine/forward/http.go:73:// validMethod returns true when s is a token per RFC 7230 §3.2.6 that
internal/engine/forward/http.go:77:func validMethod(s string) bool {
internal/engine/forward/http_test.go:271:func TestNewRejectsBadConfig(t *testing.T) {

→ Read internal/engine/forward/http_test.go [offset=250, limit=60]
→ Read internal/engine/forward/http_test.go [offset=306]

$ rg -n "StaticHeaders|static_headers" --glob '!*.md' . | head -40
./cmd/astrate/main.go:265:			StaticHeaders: f.StaticHeaders,
./docs/api/astarte_realm_management_api.yaml:1083:        http_static_headers:
./internal/config/config.go:169:	StaticHeaders map[string]string `toml:"static_headers"`
./internal/config/config.go:345:	// static_headers — a stale url left behind after disabling forwarding must
./internal/config/config.example.toml:60:# static_headers = { Authorization = "Bearer ..." }
./internal/config/config_test.go:266:static_headers = { Authorization = "Bearer tok", X-Forward = "yes" }
./internal/config/config_test.go:282:	if f.StaticHeaders["Authorization"] != "Bearer tok" || f.StaticHeaders["X-Forward"] != "yes" {
./internal/config/config_test.go:283:		t.Errorf("StaticHeaders = %v", f.StaticHeaders)
./internal/engine/forward/http.go:22:	// StaticHeaders are extra request headers, applied after the fixed ones.
./internal/engine/forward/http.go:23:	StaticHeaders map[string]string
./internal/engine/forward/http.go:68:		static: cfg.StaticHeaders,
./internal/engine/forward/http_test.go:86:func TestStaticHeaders(t *testing.T) {
./internal/engine/forward/http_test.go:110:		StaticHeaders: map[string]string{"X-Custom": "custom-val", "X-Other": "other-val"},
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

Let me check the config path and probe real net/http behaviour for header values before writing the check.
→ Read internal/config/config.go [offset=330, limit=40]

← Write /tmp/opencode/hdrprobe/main.go
Wrote file successfully.

$ go version && go run /tmp/opencode/hdrprobe/main.go
go version go1.26.7 linux/arm64
space        "a b"        accepted (status 200)
tab          "a\tb"       accepted (status 200)
lf           "a\nb"       REJECTED: Post "http://127.0.0.1:36483": net/http: invalid header field value for "X-Probe"
cr           "a\rb"       REJECTED: Post "http://127.0.0.1:36483": net/http: invalid header field value for "X-Probe"
crlf         "a\r\nb"     REJECTED: Post "http://127.0.0.1:36483": net/http: invalid header field value for "X-Probe"
nul          "a\x00b"     REJECTED: Post "http://127.0.0.1:36483": net/http: invalid header field value for "X-Probe"
del          "a\x7fb"     REJECTED: Post "http://127.0.0.1:36483": net/http: invalid header field value for "X-Probe"
lead-space   " a"         accepted (status 200)
trail-space  "a "         accepted (status 200)
obs-text     "a\x80b"     accepted (status 200)
empty        ""           accepted (status 200)

→ Read internal/engine/forward/http_test.go [limit=40]
Now the implementation:

← Edit internal/engine/forward/http.go
Index: /root/astrate-mule/internal/engine/forward/http.go
===================================================================
--- /root/astrate-mule/internal/engine/forward/http.go
+++ /root/astrate-mule/internal/engine/forward/http.go
@@ -33,9 +33,11 @@
 	client *http.Client
 }
 
 // New builds an HTTP bus forwarder. It returns an error when URL is empty or
-// does not parse, or when Method is not a valid HTTP method token.
+// does not parse, when Method is not a valid HTTP method token, or when a
+// StaticHeaders name is not a header field name or a value is not a legal
+// field value.
 func New(cfg Config) (*HTTP, error) {
 	if cfg.URL == "" {
 		return nil, fmt.Errorf("forward: URL must not be empty")
 	}

← Edit internal/engine/forward/http.go
Index: /root/astrate-mule/internal/engine/forward/http.go
===================================================================
--- /root/astrate-mule/internal/engine/forward/http.go
+++ /root/astrate-mule/internal/engine/forward/http.go
@@ -55,12 +55,24 @@
 m := cfg.Method
 if m == "" {
 	m = http.MethodPost
 } else {
-	if !validMethod(m) {
+	if !validToken(m) {
 		return nil, fmt.Errorf("forward: invalid method %q", m)
 	}
 }
+// Static headers are applied with Header.Set per delivery, and net/http
+// only checks field names and values when it writes the request — so a
+// typo in triggers.forward.static_headers would boot clean and then fail
+// every custom action. Same rule as the URL above: reject it here.
+for k, v := range cfg.StaticHeaders {
+	if !validToken(k) {
+		return nil, fmt.Errorf("forward: invalid static header name %q", k)
+	}
+	if !validHeaderValue(v) {
+		return nil, fmt.Errorf("forward: invalid value for static header %q: %q", k, v)
+	}
+}
 c := cfg.Client
 if c == nil {
 	c = http.DefaultClient
 }

← Edit internal/engine/forward/http.go
Index: /root/astrate-mule/internal/engine/forward/http.go
===================================================================
--- /root/astrate-mule/internal/engine/forward/http.go
+++ /root/astrate-mule/internal/engine/forward/http.go
@@ -83,13 +83,14 @@
 		client: c,
 	}, nil
 }
 
-// validMethod returns true when s is a token per RFC 7230 §3.2.6 that
-// is also a recognised HTTP method. We accept any single-token string
-// so that user-configured methods like PUT or PATCH are allowed without
-// a whitelist.
-func validMethod(s string) bool {
+// validToken returns true when s is a non-empty token per RFC 7230 §3.2.6.
+// That is the grammar of the method (RFC 7230 §4.1) and, identically, of the
+// header field name (RFC 7230 §3.2), so both are checked with this one
+// predicate. We accept any single-token string so that user-configured
+// methods like PUT or PATCH are allowed without a whitelist.
+func validToken(s string) bool {
 	if s == "" {
 		return false
 	}
 	for i := 0; i < len(s); i++ {
@@ -107,8 +108,23 @@
 	}
 	return true
 }
 
+// validHeaderValue reports whether v can be written as a header field value.
+// It mirrors net/http's own write-time check (httpguts.ValidHeaderFieldValue),
+// measured rather than assumed: SP, HTAB, obs-text and the empty string are
+// accepted by a real request, while CR, LF and the other control bytes are
+// refused with `net/http: invalid header field value`.
+func validHeaderValue(v string) bool {
+	for i := 0; i < len(v); i++ {
+		c := v[i]
+		if (c < 0x20 && c != '\t') || c == 0x7f {
+			return false
+		}
+	}
+	return true
+}
+
 // Forward implements the triggers.Forwarder contract.
 func (h *HTTP) Forward(ctx context.Context, realm, trigger string, action json.RawMessage, event []byte) error {
 	raw, err := marshalEnvelope(realm, trigger, action, event)
 	if err != nil {

Now the test:
← Edit internal/engine/forward/http_test.go
Index: /root/astrate-mule/internal/engine/forward/http_test.go
===================================================================
--- /root/astrate-mule/internal/engine/forward/http_test.go
+++ /root/astrate-mule/internal/engine/forward/http_test.go
@@ -288,8 +288,16 @@
 		{"relative URL", Config{URL: "nope"}},
 		{"no host", Config{URL: "http://"}},
 		{"non-http scheme", Config{URL: "ftp://bus.example/x"}},
 		{"bad method", Config{URL: srv.URL, Method: "GET NOPE"}},
+		// A static header is applied per delivery, and net/http validates
+		// names and values only at write time, so these two are exactly the
+		// typos that would boot clean and then fail every delivery. A space in
+		// a value is legal, so it must not be caught by the value case.
+		{"bad static header name", Config{URL: srv.URL, StaticHeaders: map[string]string{"X Bad Name": "v"}}},
+		{"static header name with colon", Config{URL: srv.URL, StaticHeaders: map[string]string{"X:Foo": "v"}}},
+		{"static header value with newline", Config{URL: srv.URL, StaticHeaders: map[string]string{"X-Foo": "a\nb"}}},
+		{"static header value with carriage return", Config{URL: srv.URL, StaticHeaders: map[string]string{"X-Foo": "a\rb"}}},
 	}
 	for _, b := range bad {
 		t.Run(b.name, func(t *testing.T) {
 			_, err := New(b.cfg)
@@ -297,14 +305,20 @@
 				t.Error("expected non-nil error")
 			}
 		})
 	}
+	if n := errCount.Load(); n != 0 {
+		t.Errorf("%d request(s) reached the endpoint for rejected configs, want 0", n)
+	}
 
 	// A valid config with the same URL must be accepted.
-	f, err := New(Config{URL: srv.URL})
+	f, err := New(Config{URL: srv.URL, StaticHeaders: map[string]string{"Authorization": "Bearer tok"}})
 	if err != nil {
 		t.Fatalf("New with valid config: %v", err)
 	}
 	if err := f.Forward(context.Background(), "r", "t", json.RawMessage(`{}`), []byte(`{}`)); err != nil {
 		t.Fatalf("Forward: %v", err)
 	}
+	if n := errCount.Load(); n != 1 {
+		t.Errorf("%d request(s) reached the endpoint after one Forward, want 1", n)
+	}
 }

$ go vet ./internal/engine/forward/ && go test ./internal/engine/forward/ 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/engine/forward	0.100s

Now confirm the test actually breaks without the implementation:
$ cp internal/engine/forward/http.go /tmp/opencode/http.go.bak && python3 - <<'EOF'
import re
p='internal/engine/forward/http.go'
s=open(p).read()
block=s[s.index('\t// Static headers are applied'):s.index('\tc := cfg.Client')]
open(p,'w').write(s.replace(block,''))
EOF
go test ./internal/engine/forward/ -run TestNewRejectsBadConfig 2>&1 | tail -20; cp /tmp/opencode/http.go.bak internal/engine/forward/http.go && rm /tmp/opencode/http.go.bak
--- FAIL: TestNewRejectsBadConfig (0.01s)
    --- FAIL: TestNewRejectsBadConfig/bad_static_header_name (0.00s)
        http_test.go:305: expected non-nil error
    --- FAIL: TestNewRejectsBadConfig/static_header_name_with_colon (0.00s)
        http_test.go:305: expected non-nil error
    --- FAIL: TestNewRejectsBadConfig/static_header_value_with_newline (0.00s)
        http_test.go:305: expected non-nil error
    --- FAIL: TestNewRejectsBadConfig/static_header_value_with_carriage_return (0.00s)
        http_test.go:305: expected non-nil error
FAIL
FAIL	github.com/astrate-platform/astrate/internal/engine/forward	0.025s
FAIL

$ git diff --stat && gofmt -l . && go vet ./... && go test -race ./... 2>&1 | tail -25
internal/engine/forward/http.go      | 44 ++++++++++++++++++++++++++++++------
 internal/engine/forward/http_test.go | 16 ++++++++++++-
 2 files changed, 52 insertions(+), 8 deletions(-)
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/realm	0.057s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/store	0.054s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.061s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/testutil	0.089s
?   	github.com/astrate-platform/astrate/migrations	[no test files]
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/astarteapi	0.064s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/deviceid	0.078s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/interfaceschema	0.071s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/payload	0.078s
FAIL

Race detector can't run on this 39-bit-VMA kernel (the documented Pi limitation), so running the standard gate:
$ go test ./... 2>&1 | rg -v "^ok|no test files" | head -20; echo "---exit chain done---"; gofmt -l . ; echo "gofmt clean"
---exit chain done---
gofmt clean

$ git diff
diff --git a/internal/engine/forward/http.go b/internal/engine/forward/http.go
index d4f4357..637ab95 100644
--- a/internal/engine/forward/http.go
+++ b/internal/engine/forward/http.go
@@ -34,7 +34,9 @@ type HTTP struct {
 }
 
 // New builds an HTTP bus forwarder. It returns an error when URL is empty or
-// does not parse, or when Method is not a valid HTTP method token.
+// does not parse, when Method is not a valid HTTP method token, or when a
+// StaticHeaders name is not a header field name or a value is not a legal
+// field value.
 func New(cfg Config) (*HTTP, error) {
 	if cfg.URL == "" {
 		return nil, fmt.Errorf("forward: URL must not be empty")
@@ -54,10 +56,22 @@ func New(cfg Config) (*HTTP, error) {
 	if m == "" {
 		m = http.MethodPost
 	} else {
-		if !validMethod(m) {
+		if !validToken(m) {
 			return nil, fmt.Errorf("forward: invalid method %q", m)
 		}
 	}
+	// Static headers are applied with Header.Set per delivery, and net/http
+	// only checks field names and values when it writes the request — so a
+	// typo in triggers.forward.static_headers would boot clean and then fail
+	// every custom action. Same rule as the URL above: reject it here.
+	for k, v := range cfg.StaticHeaders {
+		if !validToken(k) {
+			return nil, fmt.Errorf("forward: invalid static header name %q", k)
+		}
+		if !validHeaderValue(v) {
+			return nil, fmt.Errorf("forward: invalid value for static header %q: %q", k, v)
+		}
+	}
 	c := cfg.Client
 	if c == nil {
 		c = http.DefaultClient
@@ -70,11 +84,12 @@ func New(cfg Config) (*HTTP, error) {
 	}, nil
 }
 
-// validMethod returns true when s is a token per RFC 7230 §3.2.6 that
-// is also a recognised HTTP method. We accept any single-token string
-// so that user-configured methods like PUT or PATCH are allowed without
-// a whitelist.
-func validMethod(s string) bool {
+// validToken returns true when s is a non-empty token per RFC 7230 §3.2.6.
+// That is the grammar of the method (RFC 7230 §4.1) and, identically, of the
+// header field name (RFC 7230 §3.2), so both are checked with this one
+// predicate. We accept any single-token string so that user-configured
+// methods like PUT or PATCH are allowed without a whitelist.
+func validToken(s string) bool {
 	if s == "" {
 		return false
 	}
@@ -94,6 +109,21 @@ func validMethod(s string) bool {
 	return true
 }
 
+// validHeaderValue reports whether v can be written as a header field value.
+// It mirrors net/http's own write-time check (httpguts.ValidHeaderFieldValue),
+// measured rather than assumed: SP, HTAB, obs-text and the empty string are
+// accepted by a real request, while CR, LF and the other control bytes are
+// refused with `net/http: invalid header field value`.
+func validHeaderValue(v string) bool {
+	for i := 0; i < len(v); i++ {
+		c := v[i]
+		if (c < 0x20 && c != '\t') || c == 0x7f {
+			return false
+		}
+	}
+	return true
+}
+
 // Forward implements the triggers.Forwarder contract.
 func (h *HTTP) Forward(ctx context.Context, realm, trigger string, action json.RawMessage, event []byte) error {
 	raw, err := marshalEnvelope(realm, trigger, action, event)
diff --git a/internal/engine/forward/http_test.go b/internal/engine/forward/http_test.go
index c80745b..08309e4 100644
--- a/internal/engine/forward/http_test.go
+++ b/internal/engine/forward/http_test.go
@@ -289,6 +289,14 @@ func TestNewRejectsBadConfig(t *testing.T) {
 		{"no host", Config{URL: "http://"}},
 		{"non-http scheme", Config{URL: "ftp://bus.example/x"}},
 		{"bad method", Config{URL: srv.URL, Method: "GET NOPE"}},
+		// A static header is applied per delivery, and net/http validates
+		// names and values only at write time, so these two are exactly the
+		// typos that would boot clean and then fail every delivery. A space in
+		// a value is legal, so it must not be caught by the value case.
+		{"bad static header name", Config{URL: srv.URL, StaticHeaders: map[string]string{"X Bad Name": "v"}}},
+		{"static header name with colon", Config{URL: srv.URL, StaticHeaders: map[string]string{"X:Foo": "v"}}},
+		{"static header value with newline", Config{URL: srv.URL, StaticHeaders: map[string]string{"X-Foo": "a\nb"}}},
+		{"static header value with carriage return", Config{URL: srv.URL, StaticHeaders: map[string]string{"X-Foo": "a\rb"}}},
 	}
 	for _, b := range bad {
 		t.Run(b.name, func(t *testing.T) {
@@ -298,13 +306,19 @@ func TestNewRejectsBadConfig(t *testing.T) {
 			}
 		})
 	}
+	if n := errCount.Load(); n != 0 {
+		t.Errorf("%d request(s) reached the endpoint for rejected configs, want 0", n)
+	}
 
 	// A valid config with the same URL must be accepted.
-	f, err := New(Config{URL: srv.URL})
+	f, err := New(Config{URL: srv.URL, StaticHeaders: map[string]string{"Authorization": "Bearer tok"}})
 	if err != nil {
 		t.Fatalf("New with valid config: %v", err)
 	}
 	if err := f.Forward(context.Background(), "r", "t", json.RawMessage(`{}`), []byte(`{}`)); err != nil {
 		t.Fatalf("Forward: %v", err)
 	}
+	if n := errCount.Load(); n != 1 {
+		t.Errorf("%d request(s) reached the endpoint after one Forward, want 1", n)
+	}
 }

← Edit internal/engine/forward/http_test.go
Index: /root/astrate-mule/internal/engine/forward/http_test.go
===================================================================
--- /root/astrate-mule/internal/engine/forward/http_test.go
+++ /root/astrate-mule/internal/engine/forward/http_test.go
@@ -289,11 +289,11 @@
 {"no host", Config{URL: "http://"}},
 {"non-http scheme", Config{URL: "ftp://bus.example/x"}},
 {"bad method", Config{URL: srv.URL, Method: "GET NOPE"}},
 // A static header is applied per delivery, and net/http validates
-// names and values only at write time, so these two are exactly the
-// typos that would boot clean and then fail every delivery. A space in
-// a value is legal, so it must not be caught by the value case.
+// names and values only at write time, so these are exactly the typos
+// that would boot clean and then fail every delivery. A space in a
+// value is legal, so the value cases must be about control bytes.
 {"bad static header name", Config{URL: srv.URL, StaticHeaders: map[string]string{"X Bad Name": "v"}}},
 {"static header name with colon", Config{URL: srv.URL, StaticHeaders: map[string]string{"X:Foo": "v"}}},
 {"static header value with newline", Config{URL: srv.URL, StaticHeaders: map[string]string{"X-Foo": "a\nb"}}},
 {"static header value with carriage return", Config{URL: srv.URL, StaticHeaders: map[string]string{"X-Foo": "a\rb"}}},

$ gofmt -l . ; go vet ./... && go test ./... 2>&1 | rg -v "^ok|no test files"; echo "gate exit: $?"; git status --short
gate exit: 1
 M internal/engine/forward/http.go
 M internal/engine/forward/http_test.go

Gate is green (the `exit: 1` there is `rg` reporting "no lines matched", i.e. no failures). Working tree has only the two intended files.

```
Done: forward.New now rejects static header names that are not RFC 7230 tokens (reusing the method predicate, renamed validToken) and values net/http would refuse to write (CR/LF and other control bytes), so a static_headers typo fails at boot instead of once per delivery.
Files: internal/engine/forward/http.go (New, validMethod→validToken, new validHeaderValue), internal/engine/forward/http_test.go (TestNewRejectsBadConfig)
Verified: go vet ./... && go test ./... -> pass (all packages ok); gofmt -l . -> clean; go test -race ./... attempted -> impossible here, "ThreadSanitizer: unsupported VMA range / Found 39 - Supported 48" (the documented 39-bit-VMA kernel, so no race coverage for this change); net/http value rules measured with a throwaway probe in /tmp: space/tab/obs-text/empty accepted, LF/CR/NUL/DEL refused
Unsure: the task asked for a CR/LF rejection per value; I rejected every control byte except HTAB (plus DEL) instead, i.e. exactly what net/http's ValidHeaderFieldValue refuses, so nothing can still fail at write time. If a config somewhere carries a stray NUL or DEL in a static header value it now fails at boot where it previously failed per delivery — intended, but a wider blast radius than the literal ask. Also: with several bad headers the reported one is whichever the map iteration yields first (nondeterministic); the message is always correct, and I did not sort to keep the diff small.
Follow-ups: none
```
