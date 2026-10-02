slug: container-stop-deadline
verdict: done
at:  f9af73b
ran: 2026-10-02T16:55:25Z on DietPi in 410s


> build · big-pickle

I'll start by reading the two files.
→ Read .mule/task.md
→ Read .mule/MULE.md
→ Read internal/flow/blocks/container/docker.go

$ ls internal/flow/blocks/container/
block.go
block_test.go
docker.go
httpbridge.go

$ sed -n '110,130p;300,320p' internal/flow/blocks/block.go
sed: can't read internal/flow/blocks/block.go: No such file or directory

$ sed -n '100,135p;295,325p' internal/flow/blocks/container/block.go
inst, err := runner.Start(ctx, Spec{
		Image:          cfg.Image,
		ContainerPort:  cfg.ContainerPort,
		FlowConfigJSON: flowCfgJSON,
		Labels:         labels,
		Name:           cname,
	})
	if err != nil {
		return nil, err
	}

	bridge := &Bridge{
		BaseURL: inst.BaseURL(),
		Timeout: cfg.Timeout,
		Client:  &http.Client{Timeout: cfg.Timeout},
	}

	readyCtx, readyCancel := context.WithTimeout(context.Background(), cfg.ReadyTimeout)
	defer readyCancel()
	if err := bridge.WaitReady(readyCtx); err != nil {
		_ = inst.Stop(context.Background())
		return nil, fmt.Errorf("container %q: %w", name, err)
	}

	b := &Block{
		name:   name,
		cfg:    cfg,
		log:    slog.Default(),
		bridge: bridge,
		inst:   inst,
	}

	// Death watch (#45): when the flow service asks to be told about block
	// death and the runner can report container exit, watch in the background.
	// A clean Stop() marks the block stopped before tearing the container
	// Mark stopped before touching the container so the exit watcher sees a
	// clean shutdown rather than an unexpected death.
	b.stopped = true
	inst := b.inst
	b.inst = nil
	cancelWatch := b.cancelWatch
	b.cancelWatch = nil
	b.mu.Unlock()
	if cancelWatch != nil {
		cancelWatch()
	}

	if inst == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := inst.Stop(ctx); err != nil && b.log != nil {
		b.log.Error("container stop failed", "block", b.name, "err", err)
	}
}

var (
	_ flow.Block   = (*Block)(nil)
	_ flow.Stopper = (*Block)(nil)
)

$ wc -l internal/flow/blocks/container/block_test.go && sed -n '1,80p' internal/flow/blocks/container/block_test.go
424 internal/flow/blocks/container/block_test.go
package container_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
	"github.com/astrate-platform/astrate/internal/flow/blocks/container"
)

// fakeInstance is a test Docker handle that points at an httptest.
type fakeInstance struct {
	base   string
	id     string
	stops  atomic.Int32
	stopFn func()
}

func (f *fakeInstance) BaseURL() string { return f.base }
func (f *fakeInstance) ID() string      { return f.id }
func (f *fakeInstance) Stop(context.Context) error {
	f.stops.Add(1)
	if f.stopFn != nil {
		f.stopFn()
	}
	return nil
}

type fakeRunner struct {
	mu      sync.Mutex
	last    container.Spec
	inst    container.Instance
	err     error
	started atomic.Int32
}

func (r *fakeRunner) Start(_ context.Context, spec container.Spec) (container.Instance, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.last = spec
	r.started.Add(1)
	if r.err != nil {
		return nil, r.err
	}
	return r.inst, nil
}

func echoServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/healthz":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPost && r.URL.Path == "/v1/message":
			body, _ := io.ReadAll(r.Body)
			if len(body) == 0 {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			var probe map[string]any
			if err := json.Unmarshal(body, &probe); err != nil {
				http.Error(w, "bad json", http.StatusBadRequest)
				return
			}
			if md, ok := probe["metadata"].(map[string]any); ok {
				if v, _ := md["echo_drop"].(string); v == "1" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				if v, _ := md["echo_array"].(string); v == "1" {

$ rg -n "func Test" internal/flow/blocks/container/block_test.go
99:func TestConstructor_RequiresImage(t *testing.T) {
106:func TestConstructor_StartFailure(t *testing.T) {
114:func TestBlock_RoundTripEcho(t *testing.T) {
173:func TestBlock_DropAndError(t *testing.T) {
208:func TestBlock_StopRemovesContainer(t *testing.T) {
228:func TestBridge_WaitReadyTimeout(t *testing.T) {
244:func TestCLIRunner_ParsesDockerPort(t *testing.T) {
290:func TestCLIRunner_RunFailure(t *testing.T) {
302:func TestDefaultRegistry_HasContainer(t *testing.T) {
313:func TestCatalogConstructor_UsesInjectedRunner(t *testing.T) {
334:func TestNew_CleansUpWhenNotReady(t *testing.T) {
365:func TestWatchExit_UnexpectedDeath(t *testing.T) {
397:func TestWatchExit_CleanStopIsSilent(t *testing.T) {

$ sed -n '244,305p' internal/flow/blocks/container/block_test.go
func TestCLIRunner_ParsesDockerPort(t *testing.T) {
	var calls [][]string
	r := &container.CLIRunner{
		Run: func(_ context.Context, name string, args ...string) (string, string, error) {
			calls = append(calls, append([]string{name}, args...))
			if len(args) > 0 && args[0] == "run" {
				return "abc123deadbeef\n", "", nil
			}
			if len(args) > 0 && args[0] == "port" {
				return "127.0.0.1:34567\n", "", nil
			}
			if len(args) > 0 && args[0] == "rm" {
				return "", "", nil
			}
			return "", "unexpected", fmt.Errorf("bad")
		},
	}
	inst, err := r.Start(context.Background(), container.Spec{
		Image:  "img:tag",
		Labels: map[string]string{"astrate.flow": "1"},
	})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	if inst.BaseURL() != "http://127.0.0.1:34567" {
		t.Errorf("BaseURL = %q", inst.BaseURL())
	}
	if inst.ID() != "abc123deadbeef" {
		t.Errorf("ID = %q", inst.ID())
	}
	if err := inst.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	// docker run should publish 127.0.0.1::8080
	foundPublish := false
	for _, c := range calls {
		joined := strings.Join(c, " ")
		if strings.Contains(joined, "127.0.0.1::8080") {
			foundPublish = true
		}
	}
	if !foundPublish {
		t.Fatalf("expected publish flag in calls %#v", calls)
	}
}

func TestCLIRunner_RunFailure(t *testing.T) {
	r := &container.CLIRunner{
		Run: func(context.Context, string, ...string) (string, string, error) {
			return "", "Unable to find image 'nope:latest' locally", fmt.Errorf("exit 125")
		},
	}
	_, err := r.Start(context.Background(), container.Spec{Image: "nope:latest"})
	if err == nil || !strings.Contains(err.Error(), "Unable to find image") {
		t.Fatalf("err = %v", err)
	}
}

func TestDefaultRegistry_HasContainer(t *testing.T) {
	reg := blocks.DefaultRegistry()
	if !reg.Has(blocks.TypeContainer) {
		t.Fatal("missing container type")

$ rg -n "defaultContainerPort|stopTimeout|15 \* time.Second" internal/flow/blocks/container/
internal/flow/blocks/container/block.go:25:const defaultContainerPort = 8080
internal/flow/blocks/container/block.go:174:		ContainerPort: defaultContainerPort,
internal/flow/blocks/container/docker.go:87:		port = defaultContainerPort
internal/flow/blocks/container/httpbridge.go:21:	defaultReadyWait    = 15 * time.Second

Now the fix in `docker.go`:
← Edit internal/flow/blocks/container/docker.go
Index: /root/astrate-mule/internal/flow/blocks/container/docker.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/container/docker.go
+++ /root/astrate-mule/internal/flow/blocks/container/docker.go
@@ -29,9 +29,11 @@
 	// BaseURL is the host URL to POST messages to (e.g. http://127.0.0.1:49152).
 	BaseURL() string
 	// ID is the docker container id (or a test fake id).
 	ID() string
-	// Stop stops and removes the container (best-effort).
+	// Stop stops and removes the container (best-effort). The ctx is bounded
+	// by stopTimeout unless the caller already gave it a deadline, so a wedged
+	// docker daemon cannot hold the caller forever.
 	Stop(ctx context.Context) error
 }
 
 // Runner starts containers. PoC default is CLIRunner (docker CLI).

← Edit internal/flow/blocks/container/docker.go
Index: /root/astrate-mule/internal/flow/blocks/container/docker.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/container/docker.go
+++ /root/astrate-mule/internal/flow/blocks/container/docker.go
@@ -40,8 +40,12 @@
 type Runner interface {
 	Start(ctx context.Context, spec Spec) (Instance, error)
 }
 
+// stopTimeout bounds the `docker rm -f` issued by cliInstance.Stop when the
+// caller's context carries no deadline of its own.
+const stopTimeout = 15 * time.Second
+
 // Waiter is implemented by Instances that can block until their container
 // exits. The block's exit watcher uses it to detect unexpected death (#45).
 type Waiter interface {
 	// Wait blocks until the container exits (or ctx is cancelled) and returns

← Edit internal/flow/blocks/container/docker.go
Index: /root/astrate-mule/internal/flow/blocks/container/docker.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/container/docker.go
+++ /root/astrate-mule/internal/flow/blocks/container/docker.go
@@ -174,10 +174,16 @@
 if c == nil || c.id == "" {
 	return nil
 }
 if ctx == nil {
+	ctx = context.Background()
+}
+// Both cleanup callers pass context.Background(), so a deadline-less ctx
+// is the common case: bound `docker rm -f` ourselves instead of letting a
+// wedged daemon hold flow instantiation (or Block.Stop) open indefinitely.
+if _, ok := ctx.Deadline(); !ok {
 	var cancel context.CancelFunc
-	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
+	ctx, cancel = context.WithTimeout(ctx, stopTimeout)
 	defer cancel()
 }
 _, stderr, err := c.runner.run(ctx, "rm", "-f", c.id)
 if err != nil {

Now the tests. Two cases: `Stop(context.Background())` must reach `Run` with a deadline, and the hostPort-failure path must still `rm -f <id>`.
→ Read internal/flow/blocks/container/block_test.go [offset=286, limit=18]

← Edit internal/flow/blocks/container/block_test.go
Index: /root/astrate-mule/internal/flow/blocks/container/block_test.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/container/block_test.go
+++ /root/astrate-mule/internal/flow/blocks/container/block_test.go
@@ -298,8 +298,124 @@
 		t.Fatalf("err = %v", err)
 	}
 }
 
+func TestCLIRunner_StopBoundsDeadlineLessContext(t *testing.T) {
+	type rmCall struct {
+		args       []string
+		hasDeadlne bool
+		within     time.Duration
+	}
+	var rm []rmCall
+	r := &container.CLIRunner{
+		Run: func(ctx context.Context, _ string, args ...string) (string, string, error) {
+			if len(args) > 0 && args[0] == "run" {
+				return "abc123deadbeef\n", "", nil
+			}
+			if len(args) > 0 && args[0] == "port" {
+				return "127.0.0.1:34567\n", "", nil
+			}
+			if len(args) > 0 && args[0] == "rm" {
+				d, ok := ctx.Deadline()
+				c := rmCall{args: args, hasDeadlne: ok}
+				if ok {
+					c.within = time.Until(d)
+				}
+				rm = append(rm, c)
+				return "", "", nil
+			}
+			return "", "unexpected", fmt.Errorf("bad")
+		},
+	}
+	inst, err := r.Start(context.Background(), container.Spec{Image: "img:tag"})
+	if err != nil {
+		t.Fatalf("Start: %v", err)
+	}
+	// Both cleanup callers (docker.go Start, block.go New) pass
+	// context.Background(), so Stop must bound the `docker rm -f` itself.
+	if err := inst.Stop(context.Background()); err != nil {
+		t.Fatal(err)
+	}
+	if len(rm) != 1 {
+		t.Fatalf("rm calls = %#v, want 1", rm)
+	}
+	if !rm[0].hasDeadlne {
+		t.Fatalf("Stop(context.Background()) reached Run with no deadline: %#v", rm[0])
+	}
+	if rm[0].within <= 0 || rm[0].within > 15*time.Second {
+		t.Errorf("Stop deadline within %v, want (0, 15s]", rm[0].within)
+	}
+	if got := strings.Join(rm[0].args, " "); got != "rm -f abc123deadbeef" {
+		t.Errorf("rm args = %q", got)
+	}
+}
+
+func TestCLIRunner_StopKeepsShorterCallerDeadline(t *testing.T) {
+	var within time.Duration
+	var seen bool
+	r := &container.CLIRunner{
+		Run: func(ctx context.Context, _ string, args ...string) (string, string, error) {
+			if len(args) > 0 && args[0] == "run" {
+				return "abc123deadbeef\n", "", nil
+			}
+			if len(args) > 0 && args[0] == "port" {
+				return "127.0.0.1:34567\n", "", nil
+			}
+			if len(args) > 0 && args[0] == "rm" {
+				d, ok := ctx.Deadline()
+				seen = true
+				if ok {
+					within = time.Until(d)
+				}
+				return "", "", nil
+			}
+			return "", "unexpected", fmt.Errorf("bad")
+		},
+	}
+	inst, err := r.Start(context.Background(), container.Spec{Image: "img:tag"})
+	if err != nil {
+		t.Fatalf("Start: %v", err)
+	}
+	// A deadline the caller set must win: Stop must not widen it to 15s.
+	caller, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
+	defer cancel()
+	if err := inst.Stop(caller); err != nil {
+		t.Fatal(err)
+	}
+	if !seen {
+		t.Fatal("Stop did not reach Run")
+	}
+	if within <= 0 || within > 100*time.Millisecond {
+		t.Errorf("Stop deadline within %v, want (0, 100ms]", within)
+	}
+}
+
+func TestCLIRunner_StartCleansUpOnPortFailure(t *testing.T) {
+	var calls [][]string
+	r := &container.CLIRunner{
+		Run: func(_ context.Context, name string, args ...string) (string, string, error) {
+			calls = append(calls, append([]string{name}, args...))
+			if len(args) > 0 && args[0] == "run" {
+				return "abc123deadbeef\n", "", nil
+			}
+			if len(args) > 0 && args[0] == "port" {
+				return "", "Error: No public port '8080/tcp' published", fmt.Errorf("exit 1")
+			}
+			if len(args) > 0 && args[0] == "rm" {
+				return "", "", nil
+			}
+			return "", "unexpected", fmt.Errorf("bad")
+		},
+	}
+	if _, err := r.Start(context.Background(), container.Spec{Image: "img:tag"}); err == nil {
+		t.Fatal("Start succeeded despite docker port failure")
+	}
+	// Best-effort cleanup so we do not leave orphans on mapping failure.
+	if got := strings.Join(calls[len(calls)-1], " "); got != "docker rm -f abc123deadbeef" {
+		t.Fatalf("last call = %q, want docker rm -f abc123deadbeef (all %#v)", got, calls)
+	}
+}
+
 func TestDefaultRegistry_HasContainer(t *testing.T) {
 	reg := blocks.DefaultRegistry()
 	if !reg.Has(blocks.TypeContainer) {
 		t.Fatal("missing container type")

$ sed -i 's/hasDeadlne/hasDeadline/g' internal/flow/blocks/container/block_test.go && gofmt -l . && go vet ./internal/flow/blocks/container/ && go test ./internal/flow/blocks/container/ 2>&1 | tail -20
internal/flow/blocks/container/block_test.go
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.540s

$ gofmt -d internal/flow/blocks/container/block_test.go
diff internal/flow/blocks/container/block_test.go.orig internal/flow/blocks/container/block_test.go
--- internal/flow/blocks/container/block_test.go.orig
+++ internal/flow/blocks/container/block_test.go
@@ -301,9 +301,9 @@
 
 func TestCLIRunner_StopBoundsDeadlineLessContext(t *testing.T) {
 	type rmCall struct {
-		args       []string
+		args        []string
 		hasDeadline bool
-		within     time.Duration
+		within      time.Duration
 	}
 	var rm []rmCall
 	r := &container.CLIRunner{

$ gofmt -w internal/flow/blocks/container/block_test.go && gofmt -l . ; echo "gofmt-clean=$?"
gofmt-clean=0

Now verifying the tests actually break without the fix:
$ cp internal/flow/blocks/container/docker.go /tmp/opencode/docker.go.new && git stash list >/dev/null 2>&1; git diff --stat internal/flow/blocks/container/docker.go
internal/flow/blocks/container/docker.go | 16 ++++++++++++++--
 1 file changed, 14 insertions(+), 2 deletions(-)

$ python3 - <<'EOF'
import re
p='internal/flow/blocks/container/docker.go'
s=open(p).read()
new='''	if ctx == nil {
		ctx = context.Background()
	}
	// Both cleanup callers pass context.Background(), so a deadline-less ctx
	// is the common case: bound `docker rm -f` ourselves instead of letting a
	// wedged daemon hold flow instantiation (or Block.Stop) open indefinitely.
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, stopTimeout)
		defer cancel()
	}
'''
old='''	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
	}
'''
assert new in s
open(p,'w').write(s.replace(new,old,1))
EOF
go test -run 'TestCLIRunner_Stop' ./internal/flow/blocks/container/ 2>&1 | tail -20
--- FAIL: TestCLIRunner_StopBoundsDeadlineLessContext (0.00s)
    block_test.go:342: Stop(context.Background()) reached Run with no deadline: container_test.rmCall{args:[]string{"rm", "-f", "abc123deadbeef"}, hasDeadline:false, within:0}
FAIL
FAIL	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.030s
FAIL

$ cp /tmp/opencode/docker.go.new internal/flow/blocks/container/docker.go && rm /tmp/opencode/docker.go.new && gofmt -l . ; go vet ./... && go test ./... 2>&1 | grep -v "^ok\|no test files" ; echo "gate-exit=$?"
gate-exit=1

$ go test ./... 2>&1 | grep -v "^ok\|no test files"; echo "test-exit=${?}"
test-exit=1

$ go test -count=1 ./internal/flow/blocks/container/ && go test -race -count=1 ./internal/flow/blocks/container/ 2>&1 | tail -5
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.523s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.050s
FAIL

$ git diff --stat && git status --short
internal/flow/blocks/container/block_test.go | 116 +++++++++++++++++++++++++++
 internal/flow/blocks/container/docker.go     |  16 +++-
 2 files changed, 130 insertions(+), 2 deletions(-)
 M internal/flow/blocks/container/block_test.go
 M internal/flow/blocks/container/docker.go

Done: `cliInstance.Stop` now derives its 15s `stopTimeout` whenever ctx carries no deadline (nil or deadline-less), so both `context.Background()` cleanup callers bound the `docker rm -f`; added 3 tests, one of which fails without the fix.
Files: internal/flow/blocks/container/docker.go, internal/flow/blocks/container/block_test.go
Verified: `go vet ./... && go test ./...` -> pass; `gofmt -l .` -> clean; `go test -race` -> cannot run here (39-bit VMA, ThreadSanitizer unsupported), so no race verificationUnsure: nothing on the change itself; deriving the timeout from the caller's ctx (not `context.Background()`) preserves values/cancellation — that part is my choice, the task only specified the deadline
Follow-ups: none
