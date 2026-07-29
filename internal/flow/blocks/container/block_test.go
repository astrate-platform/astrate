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

	in := &flow.FlowMessage{
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

	outs, err := b.Process(&flow.FlowMessage{
		Key: "k", Type: flow.TypeString, Data: "x",
		Metadata: map[string]string{"echo_drop": "1"},
	})
	if err != nil || len(outs) != 0 {
		t.Fatalf("drop: outs=%v err=%v", outs, err)
	}

	_, err = b.Process(&flow.FlowMessage{
		Key: "k", Type: flow.TypeString, Data: "x",
		Metadata: map[string]string{"echo_fail": "1"},
	})
	if err == nil || !strings.Contains(err.Error(), "500") {
		t.Fatalf("fail err = %v", err)
	}

	outs, err = b.Process(&flow.FlowMessage{
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
	_, err = b.Process(&flow.FlowMessage{Key: "k", Type: flow.TypeString, Data: "x"})
	if err == nil {
		t.Fatal("want error after stop")
	}
}

func TestBridge_WaitReadyTimeout(t *testing.T) {
	// Server that never becomes ready on /healthz.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	outs, err := b.Process(&flow.FlowMessage{Key: "k", Type: flow.TypeInteger, Data: int64(7)})
	if err != nil || len(outs) != 1 {
		t.Fatalf("Process: outs=%v err=%v", outs, err)
	}
	if outs[0].Data != int64(7) {
		t.Errorf("data = %v (%T)", outs[0].Data, outs[0].Data)
	}
}

func TestNew_CleansUpWhenNotReady(t *testing.T) {
	// Port open but healthz fails forever — WaitReady times out, Stop must run.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
