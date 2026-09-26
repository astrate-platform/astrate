slug: forward-status-error-body
verdict: done
at:  915c260
ran: 2026-09-26T18:24:49Z on DietPi in 209s


> build · big-pickle

I'll start by reading the two files.
→ Read .mule/task.md
→ Read .mule/MULE.md
→ Read internal/engine/forward/http.go
→ Read internal/engine/forward/http_test.go
→ Read internal/engine/triggers/actions.go [offset=600, limit=50]

$ rg -n "forward: status" --glob '!*.md' . 2>/dev/null; echo "---"; rg -n "status %d" -g '*.go' . | head -30
./internal/engine/forward/http.go:152:		return fmt.Errorf("forward: status %d", resp.StatusCode)
---
./test/conformance/upstream/recordchannels/main.go:943:		return "", fmt.Errorf("pairing returned no certificate (status %d)", resp.StatusCode)
./internal/flow/blocks/httpblocks.go:192:		return nil, fmt.Errorf("http_source: GET %s: status %d", url, resp.StatusCode)
./internal/flow/blocks/httpblocks.go:341:		return fmt.Errorf("http_sink: %s %s: status %d", cfg.method, cfg.url, resp.StatusCode)
./internal/flow/blocks/container/httpbridge.go:81:			last = fmt.Errorf("status %d", resp.StatusCode)
./internal/engine/forward/http.go:152:		return fmt.Errorf("forward: status %d", resp.StatusCode)
./internal/engine/forward/http_test.go:250:				t.Errorf("expected non-nil error for status %d", tt.code)
./internal/engine/forward/http_test.go:255:					t.Errorf("error %q does not mention status %d", err.Error(), tt.code)
./internal/httpx/cors_test.go:55:		t.Fatalf("actual request did not reach the inner handler (status %d)", rec.Code)
./internal/httpx/cors_test.go:70:		t.Fatalf("disallowed origin blocked the request (status %d) — the browser enforces, not us", rec.Code)
./internal/httpx/cors_test.go:104:		t.Fatalf("non-CORS request blocked (status %d)", rec.Code)
./internal/httpx/cors_test.go:120:		t.Fatalf("plain OPTIONS swallowed by CORS middleware (status %d)", rec.Code)
./internal/appengine/dashboard_compat_test.go:89:				t.Fatalf("page %d: status %d: %s", page, rec.Code, rec.Body.String())
./internal/engine/triggers/policy.go:78:		return fmt.Sprintf("status %d", h.codes[0])
./internal/engine/triggers/policy.go:296:	return Decision{Strategy: StrategyDiscard, Reason: fmt.Sprintf("no handler claims status %d: discard", status)}
./internal/engine/triggers/policy_test.go:387:				t.Errorf("empty reason for status %d on %s", s, def)
./internal/flowapi/http_test.go:34:			t.Errorf("%v: status %d, want %d body=%s", tc.err, rec.Code, tc.code, rec.Body.String())
./internal/realm/http_test.go:465:				t.Errorf("got status %d body %s, want %d %s", rec.Code, rec.Body, step.wantStatus, step.wantBody)

$ rg -rn "forward\." -g '*.go' --glob '!internal/engine/forward/*' | grep -iv "fmt\.\|// " | head -20
cmd/astrate/main.go:		h, err := nNew(nConfig{
cmd/astrate/newnats_nats.go:	return nNewNATS(nNATSConfig{URL: f.URL, Subject: f.Subject})

$ rg -n "engine/forward" -g '*.go' . | head; echo ---; rg -n "Forward\(ctx" -g '*.go' internal/engine | head -20
./cmd/astrate/main.go:36:	"github.com/astrate-platform/astrate/internal/engine/forward"
./cmd/astrate/newnats_nats.go:7:	"github.com/astrate-platform/astrate/internal/engine/forward"
./internal/engine/triggers/actions.go:156:// not an HTTP webhook. internal/engine/forward provides the HTTP and NATS
---
internal/engine/forward/http.go:128:func (h *HTTP) Forward(ctx context.Context, realm, trigger string, action json.RawMessage, event []byte) error {
internal/engine/forward/http_test.go:18:	Forward(ctx context.Context, realm, trigger string, action json.RawMessage, event []byte) error
internal/engine/forward/http_test.go:320:	err = f.Forward(ctx, "r", "t", json.RawMessage(`{}`), []byte(`{}`))
internal/engine/triggers/actions.go:166:	Forward(ctx context.Context, realm, trigger string, action json.RawMessage, event []byte) error
internal/engine/triggers/actions.go:489:	if err := x.cfg.Forwarder.Forward(ctx, d.Realm, d.Trigger.Name, d.Trigger.Action.Custom, body); err != nil {
internal/engine/forward/nats_test.go:61:	if err := fwd.Forward(ctx, "realm1", "my.trigger", action, event); err != nil {
internal/engine/forward/nats_test.go:110:	if err := fwd.Forward(ctx, "r", "t", nil, nil); err != nil {
internal/engine/forward/nats_test.go:153:	if err := fwd.Forward(ctx, "r", "t", json.RawMessage{}, []byte{}); err != nil {

→ Read internal/engine/triggers/actions.go [offset=480, limit=20]

$ rg -n "LimitReader" -g '*.go' . | head -20; echo "=== body-in-error precedent ==="; rg -n "respBody|bodyPrefix|truncat|quoted" -g '*.go' internal/ | head -20
./examples/flow-container-echo/main.go:50:	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
./pkg/astarteapi/envelope.go:229:	body, err := io.ReadAll(io.LimitReader(r, maxBytes+1))
./internal/housekeeping/http.go:164:	raw, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes+1))
./internal/flow/blocks/container/httpbridge.go:133:	limited := io.LimitReader(resp.Body, b.maxBody()+1)
./internal/engine/control.go:305:	plain, err := io.ReadAll(io.LimitReader(zr, int64(declared)+1))
./internal/engine/triggers/actions.go:626:	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
=== body-in-error precedent ===
internal/flow/blocks/container/httpbridge.go:134:	respBody, err := io.ReadAll(limited)
internal/flow/blocks/container/httpbridge.go:138:	if int64(len(respBody)) > b.maxBody() {
internal/flow/blocks/container/httpbridge.go:146:		snippet := strings.TrimSpace(string(respBody))
internal/flow/blocks/container/httpbridge.go:153:	respBody = bytes.TrimSpace(respBody)
internal/flow/blocks/container/httpbridge.go:154:	if len(respBody) == 0 {
internal/flow/blocks/container/httpbridge.go:159:	if respBody[0] == '[' {
internal/flow/blocks/container/httpbridge.go:161:		if err := json.Unmarshal(respBody, &list); err != nil {
internal/flow/blocks/container/httpbridge.go:174:	if err := json.Unmarshal(respBody, &one); err != nil {
internal/flow/blocks/jsonpathmap.go:29:// Whole-message and path tokens insert raw (unquoted) JSON; time.Time becomes
internal/flow/blocks/jsonpathmap.go:226:// marshalRawJSON renders v as unquoted JSON, converting time.Time to its
internal/flow/message.go:376:// silently truncate or wrap (3.7 → 3, 1e300 → garbage). Wire decoding
internal/engine/control_test.go:86:// TestInflatePropertiesRejects covers the §4.5 bounds: truncated frames,
internal/engine/control_test.go:110:		{name: "truncated header", in: []byte{0, 0, 1}},
internal/appengine/formats_parity_test.go:197:			t.Errorf("big cell = %s, want quoted decimal string", got)
internal/appengine/downsample.go:8:// We round up so that span/bucket ≤ points: a truncated bucket would produce more

→ Read internal/flow/blocks/container/httpbridge.go [offset=120, limit=45]
Now the change. Precedent for body-in-error is at httpbridge.go:145-151, so I'll match that shape.

← Edit internal/engine/forward/http.go
Index: /root/astrate-mule/internal/engine/forward/http.go
===================================================================
--- /root/astrate-mule/internal/engine/forward/http.go
+++ /root/astrate-mule/internal/engine/forward/http.go
@@ -144,12 +144,38 @@
 	if err != nil {
 		return fmt.Errorf("forward: %w", err)
 	}
 	defer func() { _ = resp.Body.Close() }()
-	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
+
+	ok := resp.StatusCode >= 200 && resp.StatusCode <= 299
+	// A non-2xx needs its body. "forward: status 500" tells an operator
+	// nothing they can act on, and the bus's own explanation ("{"error":
+	// "unknown realm"}") is the whole difference between a fixable failure and
+	// a mystery — so read a bounded prefix of it for the error, in the same
+	// shape as the container http bridge (blocks/container/httpbridge.go).
+	// Bounded because the endpoint is not us: an unbounded read would pull an
+	// arbitrarily large body into memory and into a log line.
+	var snippet string
+	if !ok {
+		var b []byte
+		b, _ = io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
+		snippet = strings.TrimSpace(string(b))
+		if len(b) == maxErrorBody {
+			snippet += "…"
+		}
+	}
+	// Drain so the connection is reusable. Bounded for the same reason, and to
+	// match the sibling request path in this codebase, which caps its drain at
+	// the same 1 MiB (triggers/actions.go): two near-identical request paths
+	// should not drift on the line that decides how much of a peer's body we
+	// are willing to hold.
+	if _, err := io.Copy(io.Discard, io.LimitReader(resp.Body, maxDrain)); err != nil {
 		return fmt.Errorf("forward: drain: %w", err)
 	}
-	if resp.StatusCode < 200 || resp.StatusCode > 299 {
-		return fmt.Errorf("forward: status %d", resp.StatusCode)
+	if !ok {
+		if snippet == "" {
+			return fmt.Errorf("forward: status %d", resp.StatusCode)
+		}
+		return fmt.Errorf("forward: status %d: %s", resp.StatusCode, snippet)
 	}
 	return nil
 }

← Edit internal/engine/forward/http.go
Index: /root/astrate-mule/internal/engine/forward/http.go
===================================================================
--- /root/astrate-mule/internal/engine/forward/http.go
+++ /root/astrate-mule/internal/engine/forward/http.go
@@ -10,10 +10,21 @@
 	"fmt"
 	"io"
 	"net/http"
 	"net/url"
+	"strings"
 )
 
+const (
+	// maxErrorBody is how much of a non-2xx response body is quoted back in
+	// the error. Enough for a JSON error document from a bus, small enough
+	// that a hostile or broken endpoint cannot flood the logs.
+	maxErrorBody = 512
+	// maxDrain caps the body read purely to make the connection reusable,
+	// matching the webhook path in internal/engine/triggers/actions.go.
+	maxDrain = 1 << 20
+)
+
 // Config is the HTTP bus forwarder's configuration.
 type Config struct {
 	// URL is the bus endpoint every custom action is posted to. Required.
 	URL string

Now the test, beside `TestStatusTable`:
← Edit internal/engine/forward/http_test.go
Index: /root/astrate-mule/internal/engine/forward/http_test.go
===================================================================
--- /root/astrate-mule/internal/engine/forward/http_test.go
+++ /root/astrate-mule/internal/engine/forward/http_test.go
@@ -263,8 +263,67 @@
 // json.RawMessage does not need: a non-nil, zero-length action or event
 // marshals to nothing and would corrupt the whole envelope. Removing the
 // len() guards leaves TestNilActionAndEvent green, so this is the row that
 // actually binds the rule.
+// TestStatusErrorCarriesBody pins the half of a failed delivery an operator
+// needs. A bare "forward: status 500" survives the bus's own explanation
+// nowhere — the drain at the end of Forward throws the body away — so a 500
+// from the bus is unactionable. This is the row that breaks if the body prefix
+// is dropped from the error: TestStatusTable only pins that the status code is
+// mentioned, and it writes no body.
+func TestStatusErrorCarriesBody(t *testing.T) {
+	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
+		w.WriteHeader(http.StatusInternalServerError)
+		_, _ = w.Write([]byte(`{"error":"boom"}`))
+	}))
+	defer srv.Close()
+
+	f, err := New(Config{URL: srv.URL})
+	if err != nil {
+		t.Fatalf("New: %v", err)
+	}
+	err = f.Forward(context.Background(), "r", "t", json.RawMessage(`{}`), []byte(`{}`))
+	if err == nil {
+		t.Fatal("expected non-nil error for status 500")
+	}
+	if !strings.Contains(err.Error(), "500") {
+		t.Errorf("error %q does not mention the status", err.Error())
+	}
+	if !strings.Contains(err.Error(), "boom") {
+		t.Errorf("error %q does not mention the response body", err.Error())
+	}
+}
+
+// TestStatusErrorBodyBounded covers the other half: the body quoted into the
+// error is capped, so a peer that answers 500 with a megabyte of text cannot
+// turn one log line into a megabyte. Without the cap this test still passes on
+// content, so it asserts the bound itself.
+func TestStatusErrorBodyBounded(t *testing.T) {
+	const filler = "0123456789"
+	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
+		w.WriteHeader(http.StatusInternalServerError)
+		_, _ = w.Write([]byte(strings.Repeat(filler, 4096)))
+	}))
+	defer srv.Close()
+
+	f, err := New(Config{URL: srv.URL})
+	if err != nil {
+		t.Fatalf("New: %v", err)
+	}
+	err = f.Forward(context.Background(), "r", "t", json.RawMessage(`{}`), []byte(`{}`))
+	if err == nil {
+		t.Fatal("expected non-nil error for status 500")
+	}
+	// One "forward: status 500: " prefix plus the 512-byte cap, the trailing
+	// ellipsis marking the cut, and slack for the numbers in the prefix.
+	if n := len(err.Error()); n > 600 {
+		t.Errorf("error is %d bytes, want the body prefix bounded to ~512", n)
+	}
+	if !strings.Contains(err.Error(), "…") {
+		t.Errorf("error %q does not mark the body as truncated", err.Error())
+	}
+}
+
 func TestEmptyNonNilActionAndEvent(t *testing.T) {
 	var body []byte
 	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
 		body, _ = io.ReadAll(r.Body)

$ gofmt -l . && go vet ./internal/engine/... && go test -race ./internal/engine/forward/... 2>&1 | tail -20
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/engine/forward	0.045s
FAIL

Race build is unavailable on this kernel (39-bit VMA), as MULE.md documents. Running the normal gate:
$ go test ./internal/engine/... 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	0.085s
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)

Now verify the new tests actually break without the implementation change (the runner does this too):
$ cp internal/engine/forward/http.go /tmp/opencode/http.go.new && git stash list >/dev/null 2>&1; git show HEAD:internal/engine/forward/http.go > /tmp/opencode/http.go.orig && python3 - <<'EOF'
import re
src = open('/tmp/opencode/http.go.orig').read()
# reapply only the test-visible part of the change: the new error text, keeping old unbounded drain
new_body = src.replace('''	defer func() { _ = resp.Body.Close() }()
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return fmt.Errorf("forward: drain: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("forward: status %d", resp.StatusCode)
	}
	return nil''','''	defer func() { _ = resp.Body.Close() }()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
	if _, err := io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20)); err != nil {
		return fmt.Errorf("forward: drain: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if s := strings.TrimSpace(string(b)); s != "" {
			return fmt.Errorf("forward: status %d: %s", resp.StatusCode, s)
		}
		return fmt.Errorf("forward: status %d", resp.StatusCode)
	}
	return nil''')
assert new_body != src
new_body = new_body.replace('"net/url"\n)', '"net/url"\n\t"strings"\n)')
open('/root/astrate-mule/internal/engine/forward/http.go','w').write(new_body)
EOF
go test ./internal/engine/forward/... 2>&1 | tail -25
--- FAIL: TestStatusErrorBodyBounded (0.00s)
    http_test.go:322: error "forward: status 500: 01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901" does not mark the body as truncated
FAIL
FAIL	github.com/astrate-platform/astrate/internal/engine/forward	0.074s
FAIL

The marker assertion binds. Now the key check — the original status-only error must break the main new test:
$ git show HEAD:internal/engine/forward/http.go > internal/engine/forward/http.go && go test ./internal/engine/forward/... -run 'TestStatusError' 2>&1 | tail -12; echo "=== restoring ==="; cp /tmp/opencode/http.go.new internal/engine/forward/http.go && git diff --stat
--- FAIL: TestStatusErrorCarriesBody (0.01s)
    http_test.go:292: error "forward: status 500" does not mention the response body
--- FAIL: TestStatusErrorBodyBounded (0.01s)
    http_test.go:322: error "forward: status 500" does not mark the body as truncated
FAIL
FAIL	github.com/astrate-platform/astrate/internal/engine/forward	0.037s
FAIL
=== restoring ===
 internal/engine/forward/http.go      | 43 ++++++++++++++++++++++++--
 internal/engine/forward/http_test.go | 59 ++++++++++++++++++++++++++++++++++++
 2 files changed, 99 insertions(+), 3 deletions(-)

Both new tests fail against the original implementation and pass against the change. Final full gate:
$ gofmt -l . ; go vet ./... && go test ./... 2>&1 | grep -v "^ok\|no test files" | head -20; echo "GATE-EXIT=$?"
GATE-EXIT=0

→ Read internal/engine/forward/http.go [offset=138, limit=45]
Two small polish items — the `var b []byte` dance and the all-whitespace-body edge:

← Edit internal/engine/forward/http.go
Index: /root/astrate-mule/internal/engine/forward/http.go
===================================================================
--- /root/astrate-mule/internal/engine/forward/http.go
+++ /root/astrate-mule/internal/engine/forward/http.go
@@ -166,12 +166,14 @@
 // Bounded because the endpoint is not us: an unbounded read would pull an
 // arbitrarily large body into memory and into a log line.
 var snippet string
 if !ok {
-	var b []byte
-	b, _ = io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
+	// The read error is deliberately dropped: a truncated or unreadable
+	// body still has to surface as a status error, not as a read error
+	// that hides the status.
+	b, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
 	snippet = strings.TrimSpace(string(b))
-	if len(b) == maxErrorBody {
+	if snippet != "" && len(b) == maxErrorBody {
 		snippet += "…"
 	}
 }
 // Drain so the connection is reusable. Bounded for the same reason, and to

$ gofmt -l . ; go vet ./... && go test ./internal/engine/... 2>&1 | tail -6
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	0.108s
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)

$ git diff --stat; echo ===; git status --short
internal/engine/forward/http.go      | 45 +++++++++++++++++++++++++--
 internal/engine/forward/http_test.go | 59 ++++++++++++++++++++++++++++++++++++
 2 files changed, 101 insertions(+), 3 deletions(-)
===
 M internal/engine/forward/http.go
 M internal/engine/forward/http_test.go

Work left in the working tree, no git touched.

```
Done: forward.Forward now quotes a bounded 512-byte prefix of a non-2xx body into the error ("forward: status 500: {"error":"boom"}") and bounds the connection-reuse drain at 1 MiB via io.LimitReader, matching the sibling webhook path (triggers/actions.go:626).
Files: internal/engine/forward/http.go (Forward + new maxErrorBody/maxDrain consts), internal/engine/forward/http_test.go (TestStatusErrorCarriesBody, TestStatusErrorBodyBounded, both beside TestStatusTable)
Verified: gofmt -l . -> clean; go vet ./... -> pass; go test ./... -> pass; go test -race attempted and FAILED TO BUILD on this box ("ThreadSanitizer: unsupported VMA range / Found 39 - Supported 48"), the known 39-bit VMA kernel limitation, so I ran the non-race gate. Both new tests were confirmed to FAIL against the original HEAD implementation ("error \"forward: status 500\" does not mention the response body") and pass after the change; the truncation-marker assertion also fails if only the "…" is dropped.
Unsure: two judgement calls the task left open. (1) I kept a trailing "…" to mark a capped body, following the existing httpbridge.go:146-149 snippet precedent; an operator reading a truncated body can tell it was cut. (2) I did not sanitise control bytes/newlines in the body before it goes into the log line, so a multi-line or binary body still spans lines — httpbridge.go has the same exposure, and inventing a sanitiser here was more invention than the task asked for. Error text is not matched anywhere else in the tree (checked), and the sole caller logs it via x.log.Warn in triggers/actions.go:489.
Follow-ups: none
```
