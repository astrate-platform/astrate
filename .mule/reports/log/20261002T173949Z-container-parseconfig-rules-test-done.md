slug: container-parseconfig-rules-test
verdict: done
at:  32c85af
ran: 2026-10-02T17:39:49Z on DietPi in 274s


> build · big-pickle

I'll read the two files first.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/flow/blocks/container/block.go
→ Read internal/flow/blocks/container/block_test.go
→ Read internal/flow/blocks/container/docker.go

$ ls -la internal/flow/blocks/container/ && wc -l internal/flow/blocks/container/*.go
total 52
drwxr-xr-x 2 root root  4096 Oct  2 19:21 .
drwxr-xr-x 5 root root  4096 Oct  2 19:21 ..
-rw-r--r-- 1 root root  8316 Oct  2 19:21 block.go
-rw-r--r-- 1 root root 15539 Oct  2 19:21 block_test.go
-rw-r--r-- 1 root root  6629 Oct  2 18:55 docker.go
-rw-r--r-- 1 root root  4472 Sep  4 20:46 httpbridge.go
  320 internal/flow/blocks/container/block.go
  540 internal/flow/blocks/container/block_test.go
  223 internal/flow/blocks/container/docker.go
  178 internal/flow/blocks/container/httpbridge.go
 1261 total

→ Read internal/flow/blocks/container/httpbridge.go

$ uname -m; nproc; go version
aarch64
4
go version go1.26.7 linux/arm64

$ go test -coverprofile=/tmp/cov-before.out ./internal/flow/blocks/container/ 2>&1 | tail -3 && go tool cover -func=/tmp/cov-before.out | tail -5
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.538s	coverage: 73.7% of statements
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:39:	timeout			66.7%
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:46:	maxBody			66.7%
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:54:	WaitReady		56.7%
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:108:	RoundTrip		78.3%
total:											(statements)		73.7%

$ grep 'block.go:1[89][0-9]\|block.go:2[0-2][0-9]' /tmp/cov-before.out
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:181.2,183.15 3 1
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:183.15,185.3 1 1
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:186.2,188.51 2 1
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:188.51,189.26 1 1
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:190.23,191.18 1 1
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:192.11,193.69 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:197.2,197.45 1 1
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:197.45,199.40 2 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:199.40,201.4 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:202.3,202.24 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:205.2,205.51 1 1
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:205.51,207.27 2 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:207.27,209.4 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:210.3,210.52 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:213.2,213.57 1 1
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:213.57,215.27 2 1
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:215.27,217.4 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:218.3,218.57 1 1
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:221.2,221.17 1 1
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:224.32,225.23 1 1
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:226.11,227.16 1 1
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:228.13,229.21 1 0

Now I'll add the tests. Inserting after `TestConstructor_RequiresImage`:
← Edit internal/flow/blocks/container/block_test.go
Index: /root/astrate-mule/internal/flow/blocks/container/block_test.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/container/block_test.go
+++ /root/astrate-mule/internal/flow/blocks/container/block_test.go
@@ -102,8 +102,233 @@
 		t.Fatalf("err = %v, want image required", err)
 	}
 }
 
+// errRunnerReached is what a runner returns when parseConfig let a config
+// through. A test expecting a config rejection that sees this error has lost
+// the validation rule it pins — and, worse, a real deployment would have handed
+// the bad value to `docker run`.
+var errRunnerReached = errors.New("runner reached: config was not rejected")
+
+// newConfigRulesTests wires the validation cases below to a ready httptest
+// echo, so a case that should be *accepted* is accepted for a real reason
+// (parseConfig passed, container became ready) and not by accident.
+func newConfigRulesServer(t *testing.T) *httptest.Server {
+	t.Helper()
+	srv := echoServer(t)
+	t.Cleanup(srv.Close)
+	return srv
+}
+
+func TestConstructor_ConfigMustBeJSONObject(t *testing.T) {
+	// Anything that is not an object is a config authoring error: it would be
+	// silently dropped from ASTRATE_FLOW_CONFIG, so it must be refused.
+	for _, tc := range []struct {
+		name string
+		cfg  any
+	}{
+		{"string", "threshold=1"},
+		{"number", 42},
+		{"float", 0.5},
+		{"bool", true},
+		{"array", []any{1, 2}},
+		{"empty string", ""},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			r := &fakeRunner{err: errRunnerReached}
+			_, err := container.New("c", map[string]any{
+				"image":  "img",
+				"config": tc.cfg,
+			}, flow.Deps{}, r)
+			if err == nil || !strings.Contains(err.Error(), "config must be a JSON object") {
+				t.Fatalf("config %v (%T): err = %v, want config must be a JSON object", tc.cfg, tc.cfg, err)
+			}
+			if r.started.Load() != 0 {
+				t.Errorf("runner started %d times, want 0", r.started.Load())
+			}
+		})
+	}
+}
+
+func TestConstructor_ConfigObjectAccepted(t *testing.T) {
+	srv := newConfigRulesServer(t)
+
+	t.Run("object", func(t *testing.T) {
+		r := &fakeRunner{inst: &fakeInstance{base: srv.URL, id: "cfg-obj"}}
+		b, err := container.New("c", map[string]any{
+			"image":  "img",
+			"config": map[string]any{"threshold": 0.5},
+		}, flow.Deps{}, r)
+		if err != nil {
+			t.Fatalf("New: %v", err)
+		}
+		t.Cleanup(func() { b.(flow.Stopper).Stop() })
+		r.mu.Lock()
+		got := r.last.FlowConfigJSON
+		r.mu.Unlock()
+		if !strings.Contains(got, `"threshold":0.5`) {
+			t.Errorf("FlowConfigJSON = %q, want threshold", got)
+		}
+	})
+
+	t.Run("null is absent", func(t *testing.T) {
+		// An explicit JSON null carries no type to check, so it must not be
+		// rejected the way a string or number is.
+		r := &fakeRunner{inst: &fakeInstance{base: srv.URL, id: "cfg-null"}}
+		b, err := container.New("c", map[string]any{
+			"image":  "img",
+			"config": nil,
+		}, flow.Deps{}, r)
+		if err != nil {
+			t.Fatalf("New: %v", err)
+		}
+		t.Cleanup(func() { b.(flow.Stopper).Stop() })
+		r.mu.Lock()
+		got := r.last.FlowConfigJSON
+		r.mu.Unlock()
+		if got != "{}" {
+			t.Errorf("FlowConfigJSON = %q, want {}", got)
+		}
+	})
+}
+
+func TestConstructor_PortAccepted(t *testing.T) {
+	// Every numeric shape a JSON pipeline config can arrive as: int (hand-built
+	// map), int64, float64 and json.Number (decoded with UseNumber).
+	for _, tc := range []struct {
+		name string
+		port any // nil means "key absent"
+		omit bool
+		want int
+	}{
+		{name: "absent", omit: true, want: 8080},
+		{name: "null", port: nil, want: 8080},
+		{name: "int lowest", port: 1, want: 1},
+		{name: "int highest", port: 65535, want: 65535},
+		{name: "int64", port: int64(9090), want: 9090},
+		{name: "float64", port: float64(9091), want: 9091},
+		{name: "json.Number", port: json.Number("9092"), want: 9092},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			srv := newConfigRulesServer(t)
+			cfg := map[string]any{"image": "img"}
+			if !tc.omit {
+				cfg["port"] = tc.port
+			}
+			r := &fakeRunner{inst: &fakeInstance{base: srv.URL, id: "port-ok"}}
+			b, err := container.New("c", cfg, flow.Deps{}, r)
+			if err != nil {
+				t.Fatalf("New: %v", err)
+			}
+			t.Cleanup(func() { b.(flow.Stopper).Stop() })
+			r.mu.Lock()
+			got := r.last.ContainerPort
+			r.mu.Unlock()
+			if got != tc.want {
+				t.Errorf("Spec.ContainerPort = %d, want %d", got, tc.want)
+			}
+		})
+	}
+}
+
+func TestConstructor_PortRejected(t *testing.T) {
+	// The port lands in `docker run -p 127.0.0.1::<port>` (docker.go), so an
+	// out-of-range or non-numeric value must die here, before any runner call.
+	for _, tc := range []struct {
+		name string
+		port any
+	}{
+		{"zero", 0},
+		{"negative", -1},
+		{"above max", 65536},
+		{"99999", 99999},
+		{"string", "8080"},
+		{"bool", true},
+		{"array", []any{8080}},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			r := &fakeRunner{err: errRunnerReached}
+			_, err := container.New("c", map[string]any{
+				"image": "img",
+				"port":  tc.port,
+			}, flow.Deps{}, r)
+			if err == nil || !strings.Contains(err.Error(), "port must be an integer 1–65535") {
+				t.Fatalf("port %v (%T): err = %v, want port must be an integer 1–65535", tc.port, tc.port, err)
+			}
+			if r.started.Load() != 0 {
+				t.Errorf("runner started %d times: the port reached the runner", r.started.Load())
+			}
+		})
+	}
+}
+
+func TestConstructor_TimeoutsMustBePositive(t *testing.T) {
+	// Same rule for both keys; the key is part of the case so a regression in
+	// one is not masked by the other.
+	for _, key := range []string{"timeout_ms", "ready_timeout_ms"} {
+		for _, tc := range []struct {
+			name string
+			val  any
+		}{
+			{"zero", 0},
+			{"negative", -1},
+			{"string", "500"},
+			{"bool", false},
+		} {
+			t.Run(key+"/"+tc.name, func(t *testing.T) {
+				r := &fakeRunner{err: errRunnerReached}
+				_, err := container.New("c", map[string]any{
+					"image": "img",
+					key:    tc.val,
+				}, flow.Deps{}, r)
+				want := key + " must be a positive integer"
+				if err == nil || !strings.Contains(err.Error(), want) {
+					t.Fatalf("%s=%v (%T): err = %v, want %q", key, tc.val, tc.val, err, want)
+				}
+				if r.started.Load() != 0 {
+					t.Errorf("runner started %d times, want 0", r.started.Load())
+				}
+			})
+		}
+	}
+}
+
+// TestConstructor_TimeoutMsBoundsRequest: an accepted timeout_ms must actually
+// reach the bridge client, otherwise a 20s container would still hang the flow
+// for the 5s default.
+func TestConstructor_TimeoutMsBoundsRequest(t *testing.T) {
+	const handlerDelay = 2 * time.Second
+	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
+		if r.URL.Path == "/healthz" {
+			w.WriteHeader(http.StatusNoContent)
+			return
+		}
+		time.Sleep(handlerDelay)
+		w.WriteHeader(http.StatusNoContent)
+	}))
+	t.Cleanup(slow.Close)
+
+	inst := &fakeInstance{base: slow.URL, id: "slow"}
+	b, err := container.New("c", map[string]any{
+		"image":      "img",
+		"timeout_ms": 100,
+	}, flow.Deps{}, &fakeRunner{inst: inst})
+	if err != nil {
+		t.Fatalf("New: %v", err)
+	}
+	t.Cleanup(func() { b.(flow.Stopper).Stop() })
+
+	start := time.Now()
+	_, err = b.Process(&flow.Message{Key: "k", Type: flow.TypeString, Data: "x"})
+	elapsed := time.Since(start)
+	if err == nil {
+		t.Fatalf("Process succeeded, want timeout after 100ms (handler sleeps %v)", handlerDelay)
+	}
+	if elapsed >= handlerDelay {
+		t.Errorf("Process took %v, want the 100ms timeout to cut it off", elapsed)
+	}
+}
+
 func TestConstructor_StartFailure(t *testing.T) {
 	r := &fakeRunner{err: fmt.Errorf("docker: image not found")}
 	_, err := container.New("c", map[string]any{"image": "missing:latest"}, flow.Deps{}, r)
 	if err == nil || !strings.Contains(err.Error(), "image not found") {

← Edit internal/flow/blocks/container/block_test.go
Index: /root/astrate-mule/internal/flow/blocks/container/block_test.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/container/block_test.go
+++ /root/astrate-mule/internal/flow/blocks/container/block_test.go
@@ -2,8 +2,9 @@
 
 import (
 	"context"
 	"encoding/json"
+	"errors"
 	"fmt"
 	"io"
 	"net/http"
 	"net/http/httptest"

$ gofmt -l . && go vet ./internal/flow/blocks/container/ && go test -count=1 ./internal/flow/blocks/container/ 2>&1 | tail -20
internal/flow/blocks/container/block_test.go
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	2.565s

$ gofmt -d internal/flow/blocks/container/block_test.go | head -40
diff internal/flow/blocks/container/block_test.go.orig internal/flow/blocks/container/block_test.go
--- internal/flow/blocks/container/block_test.go.orig
+++ internal/flow/blocks/container/block_test.go
@@ -279,7 +279,7 @@
 				r := &fakeRunner{err: errRunnerReached}
 				_, err := container.New("c", map[string]any{
 					"image": "img",
-					key:    tc.val,
+					key:     tc.val,
 				}, flow.Deps{}, r)
 				want := key + " must be a positive integer"
 				if err == nil || !strings.Contains(err.Error(), want) {

$ gofmt -w internal/flow/blocks/container/block_test.go && gofmt -l . ; echo "gofmt clean: $?"
gofmt clean: 0

Now the important check — do these tests actually fail if the rules are removed? Backing up `block.go`, then mutating it:
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cp internal/flow/blocks/container/block.go /tmp/block.go.orig && md5sum /tmp/block.go.orig internal/flow/blocks/container/block.go failed
Error: The user rejected permission to use this specific tool call.
