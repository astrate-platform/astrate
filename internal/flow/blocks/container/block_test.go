package container_test

import (
	"context"
	"encoding/json"
	"errors"
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
					// Return two copies as array.
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte("[" + string(body) + "," + string(body) + "]"))
					return
				}
				if v, _ := md["echo_fail"].(string); v == "1" {
					http.Error(w, "boom", http.StatusInternalServerError)
					return
				}
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(body)
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestConstructor_RequiresImage(t *testing.T) {
	_, err := container.New("c", map[string]any{}, flow.Deps{}, &fakeRunner{})
	if err == nil || !strings.Contains(err.Error(), "image") {
		t.Fatalf("err = %v, want image required", err)
	}
}

// errRunnerReached is what a runner returns when parseConfig let a config
// through. A test expecting a config rejection that sees this error has lost
// the validation rule it pins — and, worse, a real deployment would have handed
// the bad value to `docker run`.
var errRunnerReached = errors.New("runner reached: config was not rejected")

// newConfigRulesTests wires the validation cases below to a ready httptest
// echo, so a case that should be *accepted* is accepted for a real reason
// (parseConfig passed, container became ready) and not by accident.
func newConfigRulesServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := echoServer(t)
	t.Cleanup(srv.Close)
	return srv
}

func TestConstructor_ConfigMustBeJSONObject(t *testing.T) {
	// Anything that is not an object is a config authoring error: it would be
	// silently dropped from ASTRATE_FLOW_CONFIG, so it must be refused.
	for _, tc := range []struct {
		name string
		cfg  any
	}{
		{"string", "threshold=1"},
		{"number", 42},
		{"float", 0.5},
		{"bool", true},
		{"array", []any{1, 2}},
		{"empty string", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := &fakeRunner{err: errRunnerReached}
			_, err := container.New("c", map[string]any{
				"image":  "img",
				"config": tc.cfg,
			}, flow.Deps{}, r)
			if err == nil || !strings.Contains(err.Error(), "config must be a JSON object") {
				t.Fatalf("config %v (%T): err = %v, want config must be a JSON object", tc.cfg, tc.cfg, err)
			}
			if r.started.Load() != 0 {
				t.Errorf("runner started %d times, want 0", r.started.Load())
			}
		})
	}
}

func TestConstructor_ConfigObjectAccepted(t *testing.T) {
	srv := newConfigRulesServer(t)

	t.Run("object", func(t *testing.T) {
		r := &fakeRunner{inst: &fakeInstance{base: srv.URL, id: "cfg-obj"}}
		b, err := container.New("c", map[string]any{
			"image":  "img",
			"config": map[string]any{"threshold": 0.5},
		}, flow.Deps{}, r)
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		t.Cleanup(func() { b.(flow.Stopper).Stop() })
		r.mu.Lock()
		got := r.last.FlowConfigJSON
		r.mu.Unlock()
		if !strings.Contains(got, `"threshold":0.5`) {
			t.Errorf("FlowConfigJSON = %q, want threshold", got)
		}
	})

	t.Run("null is absent", func(t *testing.T) {
		// An explicit JSON null carries no type to check, so it must not be
		// rejected the way a string or number is.
		r := &fakeRunner{inst: &fakeInstance{base: srv.URL, id: "cfg-null"}}
		b, err := container.New("c", map[string]any{
			"image":  "img",
			"config": nil,
		}, flow.Deps{}, r)
		if err != nil {
			t.Fatalf("New: %v", err)
		}
		t.Cleanup(func() { b.(flow.Stopper).Stop() })
		r.mu.Lock()
		got := r.last.FlowConfigJSON
		r.mu.Unlock()
		if got != "{}" {
			t.Errorf("FlowConfigJSON = %q, want {}", got)
		}
	})
}

func TestConstructor_PortAccepted(t *testing.T) {
	// Every numeric shape a JSON pipeline config can arrive as: int (hand-built
	// map), int64, float64 and json.Number (decoded with UseNumber).
	for _, tc := range []struct {
		name string
		port any // nil means "key absent"
		omit bool
		want int
	}{
		{name: "absent", omit: true, want: 8080},
		{name: "null", port: nil, want: 8080},
		{name: "int lowest", port: 1, want: 1},
		{name: "int highest", port: 65535, want: 65535},
		{name: "int64", port: int64(9090), want: 9090},
		{name: "float64", port: float64(9091), want: 9091},
		{name: "json.Number", port: json.Number("9092"), want: 9092},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := newConfigRulesServer(t)
			cfg := map[string]any{"image": "img"}
			if !tc.omit {
				cfg["port"] = tc.port
			}
			r := &fakeRunner{inst: &fakeInstance{base: srv.URL, id: "port-ok"}}
			b, err := container.New("c", cfg, flow.Deps{}, r)
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			t.Cleanup(func() { b.(flow.Stopper).Stop() })
			r.mu.Lock()
			got := r.last.ContainerPort
			r.mu.Unlock()
			if got != tc.want {
				t.Errorf("Spec.ContainerPort = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestConstructor_PortRejected(t *testing.T) {
	// The port lands in `docker run -p 127.0.0.1::<port>` (docker.go), so an
	// out-of-range or non-numeric value must die here, before any runner call.
	for _, tc := range []struct {
		name string
		port any
	}{
		{"zero", 0},
		{"negative", -1},
		{"above max", 65536},
		{"99999", 99999},
		{"string", "8080"},
		{"bool", true},
		{"array", []any{8080}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := &fakeRunner{err: errRunnerReached}
			_, err := container.New("c", map[string]any{
				"image": "img",
				"port":  tc.port,
			}, flow.Deps{}, r)
			if err == nil || !strings.Contains(err.Error(), "port must be an integer 1–65535") {
				t.Fatalf("port %v (%T): err = %v, want port must be an integer 1–65535", tc.port, tc.port, err)
			}
			if r.started.Load() != 0 {
				t.Errorf("runner started %d times: the port reached the runner", r.started.Load())
			}
		})
	}
}

func TestConstructor_TimeoutsMustBePositive(t *testing.T) {
	// Same rule for both keys; the key is part of the case so a regression in
	// one is not masked by the other.
	for _, key := range []string{"timeout_ms", "ready_timeout_ms"} {
		for _, tc := range []struct {
			name string
			val  any
		}{
			{"zero", 0},
			{"negative", -1},
			{"string", "500"},
			{"bool", false},
		} {
			t.Run(key+"/"+tc.name, func(t *testing.T) {
				r := &fakeRunner{err: errRunnerReached}
				_, err := container.New("c", map[string]any{
					"image": "img",
					key:     tc.val,
				}, flow.Deps{}, r)
				want := key + " must be a positive integer"
				if err == nil || !strings.Contains(err.Error(), want) {
					t.Fatalf("%s=%v (%T): err = %v, want %q", key, tc.val, tc.val, err, want)
				}
				if r.started.Load() != 0 {
					t.Errorf("runner started %d times, want 0", r.started.Load())
				}
			})
		}
	}
}

// TestConstructor_TimeoutMsBoundsRequest: an accepted timeout_ms must actually
// reach the bridge client, otherwise a 20s container would still hang the flow
// for the 5s default.
func TestConstructor_TimeoutMsBoundsRequest(t *testing.T) {
	const handlerDelay = 2 * time.Second
	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		time.Sleep(handlerDelay)
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(slow.Close)

	inst := &fakeInstance{base: slow.URL, id: "slow"}
	b, err := container.New("c", map[string]any{
		"image":      "img",
		"timeout_ms": 100,
	}, flow.Deps{}, &fakeRunner{inst: inst})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { b.(flow.Stopper).Stop() })

	start := time.Now()
	_, err = b.Process(&flow.Message{Key: "k", Type: flow.TypeString, Data: "x"})
	elapsed := time.Since(start)
	if err == nil {
		t.Fatalf("Process succeeded, want timeout after 100ms (handler sleeps %v)", handlerDelay)
	}
	if elapsed >= handlerDelay {
		t.Errorf("Process took %v, want the 100ms timeout to cut it off", elapsed)
	}
}

func TestConstructor_StartFailure(t *testing.T) {
	r := &fakeRunner{err: fmt.Errorf("docker: image not found")}
	_, err := container.New("c", map[string]any{"image": "missing:latest"}, flow.Deps{}, r)
	if err == nil || !strings.Contains(err.Error(), "image not found") {
		t.Fatalf("err = %v", err)
	}
}

func TestBlock_RoundTripEcho(t *testing.T) {
	srv := echoServer(t)
	t.Cleanup(srv.Close)

	inst := &fakeInstance{base: srv.URL, id: "fake-1"}
	r := &fakeRunner{inst: inst}

	b, err := container.New("enrich", map[string]any{
		"image":  "astrate/flow-container-echo:poc",
		"config": map[string]any{"threshold": 0.5},
	}, flow.Deps{Realm: "acme", FlowName: "demo"}, r)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() {
		if s, ok := b.(flow.Stopper); ok {
			s.Stop()
		}
	})

	if b.Name() != "enrich" {
		t.Errorf("Name = %q", b.Name())
	}
	if r.started.Load() != 1 {
		t.Fatalf("started = %d", r.started.Load())
	}
	r.mu.Lock()
	if r.last.Image != "astrate/flow-container-echo:poc" {
		t.Errorf("image = %q", r.last.Image)
	}
	if r.last.Labels["astrate.realm"] != "acme" {
		t.Errorf("labels realm = %v", r.last.Labels)
	}
	if r.last.Labels["astrate.flow_name"] != "demo" {
		t.Errorf("labels flow_name = %v", r.last.Labels)
	}
	if !strings.Contains(r.last.FlowConfigJSON, "threshold") {
		t.Errorf("FlowConfigJSON = %q", r.last.FlowConfigJSON)
	}
	r.mu.Unlock()

	in := &flow.Message{
		Key:       "device/path",
		Type:      flow.TypeString,
		Data:      "hello",
		Timestamp: 42,
	}
	outs, err := b.Process(in)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(outs) != 1 {
		t.Fatalf("outs len = %d", len(outs))
	}
	if outs[0].Key != "device/path" || outs[0].Data != "hello" {
		t.Errorf("out = %+v", outs[0])
	}
}

func TestBlock_DropAndError(t *testing.T) {
	srv := echoServer(t)
	t.Cleanup(srv.Close)
	inst := &fakeInstance{base: srv.URL, id: "fake-2"}
	b, err := container.New("c", map[string]any{"image": "img"}, flow.Deps{}, &fakeRunner{inst: inst})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { b.(flow.Stopper).Stop() })

	outs, err := b.Process(&flow.Message{
		Key: "k", Type: flow.TypeString, Data: "x",
		Metadata: map[string]string{"echo_drop": "1"},
	})
	if err != nil || len(outs) != 0 {
		t.Fatalf("drop: outs=%v err=%v", outs, err)
	}

	_, err = b.Process(&flow.Message{
		Key: "k", Type: flow.TypeString, Data: "x",
		Metadata: map[string]string{"echo_fail": "1"},
	})
	if err == nil || !strings.Contains(err.Error(), "500") {
		t.Fatalf("fail err = %v", err)
	}

	outs, err = b.Process(&flow.Message{
		Key: "k", Type: flow.TypeString, Data: "x",
		Metadata: map[string]string{"echo_array": "1"},
	})
	if err != nil || len(outs) != 2 {
		t.Fatalf("array: len=%d err=%v", len(outs), err)
	}
}

func TestBlock_StopRemovesContainer(t *testing.T) {
	srv := echoServer(t)
	t.Cleanup(srv.Close)
	inst := &fakeInstance{base: srv.URL, id: "fake-3"}
	b, err := container.New("c", map[string]any{"image": "img"}, flow.Deps{}, &fakeRunner{inst: inst})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	s := b.(flow.Stopper)
	s.Stop()
	s.Stop() // idempotent
	if inst.stops.Load() != 1 {
		t.Fatalf("stops = %d, want 1", inst.stops.Load())
	}
	_, err = b.Process(&flow.Message{Key: "k", Type: flow.TypeString, Data: "x"})
	if err == nil {
		t.Fatal("want error after stop")
	}
}

func TestBridge_WaitReadyTimeout(t *testing.T) {
	// Server that never becomes ready on /healthz.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "no", http.StatusServiceUnavailable)
	}))
	t.Cleanup(srv.Close)

	bridge := &container.Bridge{BaseURL: srv.URL, Timeout: 50 * time.Millisecond}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	err := bridge.WaitReady(ctx)
	if err == nil {
		t.Fatal("want not ready error")
	}
}

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

func TestCLIRunner_StopBoundsDeadlineLessContext(t *testing.T) {
	type rmCall struct {
		args        []string
		hasDeadline bool
		within      time.Duration
	}
	var rm []rmCall
	r := &container.CLIRunner{
		Run: func(ctx context.Context, _ string, args ...string) (string, string, error) {
			if len(args) > 0 && args[0] == "run" {
				return "abc123deadbeef\n", "", nil
			}
			if len(args) > 0 && args[0] == "port" {
				return "127.0.0.1:34567\n", "", nil
			}
			if len(args) > 0 && args[0] == "rm" {
				d, ok := ctx.Deadline()
				c := rmCall{args: args, hasDeadline: ok}
				if ok {
					c.within = time.Until(d)
				}
				rm = append(rm, c)
				return "", "", nil
			}
			return "", "unexpected", fmt.Errorf("bad")
		},
	}
	inst, err := r.Start(context.Background(), container.Spec{Image: "img:tag"})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	// Both cleanup callers (docker.go Start, block.go New) pass
	// context.Background(), so Stop must bound the `docker rm -f` itself.
	if err := inst.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(rm) != 1 {
		t.Fatalf("rm calls = %#v, want 1", rm)
	}
	if !rm[0].hasDeadline {
		t.Fatalf("Stop(context.Background()) reached Run with no deadline: %#v", rm[0])
	}
	if rm[0].within <= 0 || rm[0].within > 15*time.Second {
		t.Errorf("Stop deadline within %v, want (0, 15s]", rm[0].within)
	}
	if got := strings.Join(rm[0].args, " "); got != "rm -f abc123deadbeef" {
		t.Errorf("rm args = %q", got)
	}
}

func TestCLIRunner_StopKeepsShorterCallerDeadline(t *testing.T) {
	var within time.Duration
	var seen bool
	r := &container.CLIRunner{
		Run: func(ctx context.Context, _ string, args ...string) (string, string, error) {
			if len(args) > 0 && args[0] == "run" {
				return "abc123deadbeef\n", "", nil
			}
			if len(args) > 0 && args[0] == "port" {
				return "127.0.0.1:34567\n", "", nil
			}
			if len(args) > 0 && args[0] == "rm" {
				d, ok := ctx.Deadline()
				seen = true
				if ok {
					within = time.Until(d)
				}
				return "", "", nil
			}
			return "", "unexpected", fmt.Errorf("bad")
		},
	}
	inst, err := r.Start(context.Background(), container.Spec{Image: "img:tag"})
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	// A deadline the caller set must win: Stop must not widen it to 15s.
	caller, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := inst.Stop(caller); err != nil {
		t.Fatal(err)
	}
	if !seen {
		t.Fatal("Stop did not reach Run")
	}
	if within <= 0 || within > 100*time.Millisecond {
		t.Errorf("Stop deadline within %v, want (0, 100ms]", within)
	}
}

func TestCLIRunner_StartCleansUpOnPortFailure(t *testing.T) {
	var calls [][]string
	r := &container.CLIRunner{
		Run: func(_ context.Context, name string, args ...string) (string, string, error) {
			calls = append(calls, append([]string{name}, args...))
			if len(args) > 0 && args[0] == "run" {
				return "abc123deadbeef\n", "", nil
			}
			if len(args) > 0 && args[0] == "port" {
				return "", "Error: No public port '8080/tcp' published", fmt.Errorf("exit 1")
			}
			if len(args) > 0 && args[0] == "rm" {
				return "", "", nil
			}
			return "", "unexpected", fmt.Errorf("bad")
		},
	}
	if _, err := r.Start(context.Background(), container.Spec{Image: "img:tag"}); err == nil {
		t.Fatal("Start succeeded despite docker port failure")
	}
	// Best-effort cleanup so we do not leave orphans on mapping failure.
	if got := strings.Join(calls[len(calls)-1], " "); got != "docker rm -f abc123deadbeef" {
		t.Fatalf("last call = %q, want docker rm -f abc123deadbeef (all %#v)", got, calls)
	}
}

func TestDefaultRegistry_HasContainer(t *testing.T) {
	reg := blocks.DefaultRegistry()
	if !reg.Has(blocks.TypeContainer) {
		t.Fatal("missing container type")
	}
	info, ok := blocks.LookupInfo(blocks.TypeContainer)
	if !ok || info.Role != blocks.RoleTransform {
		t.Fatalf("info = %+v ok=%v", info, ok)
	}
}

func TestCatalogConstructor_UsesInjectedRunner(t *testing.T) {
	srv := echoServer(t)
	t.Cleanup(srv.Close)
	inst := &fakeInstance{base: srv.URL, id: "catalog"}
	restore := container.SetDefaultRunner(&fakeRunner{inst: inst})
	t.Cleanup(restore)

	b, err := container.Constructor("c", map[string]any{"image": "img"}, flow.Deps{Realm: "r"})
	if err != nil {
		t.Fatalf("Constructor: %v", err)
	}
	t.Cleanup(func() { b.(flow.Stopper).Stop() })
	outs, err := b.Process(&flow.Message{Key: "k", Type: flow.TypeInteger, Data: int64(7)})
	if err != nil || len(outs) != 1 {
		t.Fatalf("Process: outs=%v err=%v", outs, err)
	}
	if outs[0].Data != int64(7) {
		t.Errorf("data = %v (%T)", outs[0].Data, outs[0].Data)
	}
}

func TestNew_CleansUpWhenNotReady(t *testing.T) {
	// Port open but healthz fails forever — WaitReady times out, Stop must run.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "nope", http.StatusServiceUnavailable)
	}))
	t.Cleanup(srv.Close)

	inst := &fakeInstance{base: srv.URL, id: "cleanup"}
	r := &fakeRunner{inst: inst}
	_, err := container.New("c", map[string]any{
		"image":            "img",
		"ready_timeout_ms": 150,
	}, flow.Deps{}, r)
	if err == nil {
		t.Fatal("want ready error")
	}
	if inst.stops.Load() != 1 {
		t.Fatalf("stops after failed ready = %d, want 1", inst.stops.Load())
	}
}

// waitingInstance adds Waiter so the death watcher engages.
type waitingInstance struct {
	fakeInstance
	waitFn func(ctx context.Context) (int, error)
}

func (w *waitingInstance) Wait(ctx context.Context) (int, error) { return w.waitFn(ctx) }

// TestWatchExit_UnexpectedDeath: when the container exits while the flow is
// running, the block reports fatal with its name.
func TestWatchExit_UnexpectedDeath(t *testing.T) {
	srv := echoServer(t)
	defer srv.Close()

	notifyCh := make(chan string, 1)
	inst := &waitingInstance{
		fakeInstance: fakeInstance{base: srv.URL, id: "deadbeef"},
		waitFn: func(context.Context) (int, error) {
			return 137, nil // container died mid-run
		},
	}
	runner := &fakeRunner{inst: inst}
	b, err := container.New("n1", map[string]any{"image": "x:1"}, flow.Deps{
		NotifyFatal: func(block string, _ error) { notifyCh <- block },
	}, runner)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case got := <-notifyCh:
		if got != "n1" {
			t.Fatalf("fatal block = %q, want n1", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("expected NotifyFatal on unexpected container exit")
	}
	b.(flow.Stopper).Stop()
}

// TestWatchExit_CleanStopIsSilent: a clean Stop must not report fatal even
// though tearing the container down ends the wait.
func TestWatchExit_CleanStopIsSilent(t *testing.T) {
	srv := echoServer(t)
	defer srv.Close()

	var fired atomic.Bool
	blockCtx, cancel := context.WithCancel(context.Background())
	inst := &waitingInstance{
		fakeInstance: fakeInstance{base: srv.URL, id: "deadbeef"},
		waitFn: func(ctx context.Context) (int, error) {
			<-ctx.Done() // wait returns because Stop cancelled it
			return -1, ctx.Err()
		},
	}
	_ = cancel
	_ = blockCtx
	runner := &fakeRunner{inst: inst}
	b, err := container.New("n1", map[string]any{"image": "x:1"}, flow.Deps{
		NotifyFatal: func(string, error) { fired.Store(true) },
	}, runner)
	if err != nil {
		t.Fatal(err)
	}
	b.(flow.Stopper).Stop()
	time.Sleep(100 * time.Millisecond)
	if fired.Load() {
		t.Fatal("clean stop must not trigger NotifyFatal")
	}
}
