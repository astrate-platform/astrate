package blocks_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
)

// staticServer serves a fixed status/content-type/body on every request and
// closes itself when the test finishes.
func staticServer(t *testing.T, status int, contentType, body string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(status)
		_, _ = io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func mustHTTPSource(t *testing.T, config map[string]any) flow.Source {
	t.Helper()
	b, err := blocks.HTTPSource("src", config, flow.Deps{})
	if err != nil {
		t.Fatalf("HTTPSource: %v", err)
	}
	src, ok := b.(flow.Source)
	if !ok {
		t.Fatal("block is not a Source")
	}
	return src
}

func emitOnce(t *testing.T, src flow.Source) *flow.Message {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := src.Emit(ctx)
	if err != nil {
		t.Fatalf("Emit: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("len(out) = %d, want 1", len(out))
	}
	return out[0]
}

func TestHTTPSource_HappyPath(t *testing.T) {
	const body = `{"value":7}`
	srv := staticServer(t, http.StatusOK, "application/json; charset=utf-8", body)
	src := mustHTTPSource(t, map[string]any{"urls": []string{srv.URL}, "interval_ms": 1})

	msg := emitOnce(t, src)
	if string(msg.Data.([]byte)) != body {
		t.Errorf("Data = %q, want %q", msg.Data, body)
	}
	if msg.Type != flow.TypeBinary {
		t.Errorf("Type = %v, want binary", msg.Type)
	}
	if msg.Subtype != "application/json" {
		t.Errorf("Subtype = %q, want parameters stripped to application/json", msg.Subtype)
	}
	if msg.Key != srv.URL {
		t.Errorf("Key = %q, want %q", msg.Key, srv.URL)
	}
	if msg.Metadata != nil {
		t.Errorf("Metadata = %v, want nil", msg.Metadata)
	}
	if msg.Timestamp == 0 {
		t.Error("timestamp not set")
	}
}

func TestHTTPSource_RoundRobin(t *testing.T) {
	srv1 := staticServer(t, http.StatusOK, "text/plain", "one")
	srv2 := staticServer(t, http.StatusOK, "", "two")
	src := mustHTTPSource(t, map[string]any{
		"urls":        []string{srv1.URL, srv2.URL},
		"interval_ms": 1,
	})

	first := emitOnce(t, src)
	if first.Key != srv1.URL || string(first.Data.([]byte)) != "one" {
		t.Errorf("first poll = key %q data %q, want %q/one", first.Key, first.Data, srv1.URL)
	}
	second := emitOnce(t, src)
	if second.Key != srv2.URL || string(second.Data.([]byte)) != "two" {
		t.Errorf("second poll = key %q data %q, want %q/two", second.Key, second.Data, srv2.URL)
	}
	third := emitOnce(t, src)
	if third.Key != srv1.URL {
		t.Errorf("third poll = key %q, want wrap-around to %q", third.Key, srv1.URL)
	}
}

func TestHTTPSource_IntervalRespectedIncludingFirstCall(t *testing.T) {
	srv := staticServer(t, http.StatusOK, "text/plain", "x")
	src := mustHTTPSource(t, map[string]any{
		"urls":        []string{srv.URL},
		"interval_ms": 30,
	})

	start := time.Now()
	emitOnce(t, src)
	elapsed := time.Since(start)
	if elapsed < 25*time.Millisecond {
		t.Errorf("first Emit returned after %v, want >= 25ms (interval wait including first call)", elapsed)
	}
}

func TestHTTPSource_StatusError(t *testing.T) {
	srv := staticServer(t, http.StatusInternalServerError, "text/plain", "boom")
	src := mustHTTPSource(t, map[string]any{"urls": []string{srv.URL}, "interval_ms": 1})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := src.Emit(ctx)
	if err == nil || err.Error() != "http_source: GET "+srv.URL+": status 500" {
		t.Fatalf("err = %v, want %q", err, "http_source: GET "+srv.URL+": status 500")
	}
}

func TestHTTPSource_TransportErrorWrapped(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	url := srv.URL
	srv.Close()

	src := mustHTTPSource(t, map[string]any{"urls": []string{url}, "interval_ms": 1})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := src.Emit(ctx)
	if err == nil || !strings.HasPrefix(err.Error(), "http_source: GET "+url+": ") {
		t.Fatalf("err = %v, want wrapped \"http_source: GET %s: ...\" transport error", err, url)
	}
}

func TestHTTPSource_RequestTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(300 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	src := mustHTTPSource(t, map[string]any{
		"urls":        []string{srv.URL},
		"interval_ms": 1,
		"timeout_ms":  50,
	})

	start := time.Now()
	_, err := src.Emit(context.Background())
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("Emit succeeded despite server exceeding timeout_ms")
	}
	if elapsed > 250*time.Millisecond {
		t.Errorf("Emit took %v, want timeout_ms to bound the request", elapsed)
	}
}

func TestHTTPSource_EmitCancelledWhileWaiting(t *testing.T) {
	srv := staticServer(t, http.StatusOK, "", "x")
	src := mustHTTPSource(t, map[string]any{
		"urls":        []string{srv.URL},
		"interval_ms": 60000,
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	_, err := src.Emit(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
}

func TestHTTPSource_ConstructRejections(t *testing.T) {
	tests := []struct {
		name string
		bad  map[string]any
		good map[string]any
		want string
	}{
		{
			name: "missing urls",
			bad:  map[string]any{},
			good: map[string]any{"urls": []string{"http://localhost/poll"}},
			want: `http_source: at least one url is required`,
		},
		{
			name: "empty url list",
			bad:  map[string]any{"urls": []any{}},
			good: map[string]any{"urls": []any{"http://localhost/poll"}},
			want: `http_source: at least one url is required`,
		},
		{
			name: "blank url entry",
			bad:  map[string]any{"urls": []any{"http://localhost/poll", ""}},
			good: map[string]any{"urls": []any{"http://localhost/poll"}},
			want: `http_source: at least one url is required`,
		},
		{
			name: "zero interval",
			bad:  map[string]any{"urls": []string{"http://localhost/poll"}, "interval_ms": 0},
			good: map[string]any{"urls": []string{"http://localhost/poll"}, "interval_ms": 1},
			want: `http_source: interval_ms must be positive`,
		},
		{
			name: "negative interval",
			bad:  map[string]any{"urls": []string{"http://localhost/poll"}, "interval_ms": -5},
			good: map[string]any{"urls": []string{"http://localhost/poll"}, "interval_ms": 5},
			want: `http_source: interval_ms must be positive`,
		},
		{
			name: "zero timeout",
			bad:  map[string]any{"urls": []string{"http://localhost/poll"}, "timeout_ms": 0},
			good: map[string]any{"urls": []string{"http://localhost/poll"}, "timeout_ms": 1},
			want: `http_source: timeout_ms must be positive`,
		},
		{
			name: "negative timeout",
			bad:  map[string]any{"urls": []string{"http://localhost/poll"}, "timeout_ms": -5},
			good: map[string]any{"urls": []string{"http://localhost/poll"}, "timeout_ms": 5},
			want: `http_source: timeout_ms must be positive`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+"/reject", func(t *testing.T) {
			_, err := blocks.HTTPSource("src", tt.bad, flow.Deps{})
			if err == nil || err.Error() != tt.want {
				t.Fatalf("err = %v, want %q", err, tt.want)
			}
		})
		t.Run(tt.name+"/accept", func(t *testing.T) {
			if _, err := blocks.HTTPSource("src", tt.good, flow.Deps{}); err != nil {
				t.Fatalf("twin rejected: %v", err)
			}
		})
	}
}

// --- http_sink ---

type capturedRequest struct {
	method      string
	path        string
	body        []byte
	contentType string
	header      http.Header
}

type captureState struct {
	mu   sync.Mutex
	reqs []capturedRequest
}

func (c *captureState) record(r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.reqs = append(c.reqs, capturedRequest{
		method:      r.Method,
		path:        r.URL.Path,
		body:        body,
		contentType: r.Header.Get("Content-Type"),
		header:      r.Header.Clone(),
	})
}

func (c *captureState) snapshot() []capturedRequest {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]capturedRequest(nil), c.reqs...)
}

// newSinkServer records every request; /fail answers 500 and /nocontent
// answers 204, everything else 200.
func newSinkServer(t *testing.T) (*httptest.Server, *captureState) {
	t.Helper()
	state := &captureState{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state.record(r)
		switch r.URL.Path {
		case "/fail":
			w.WriteHeader(http.StatusInternalServerError)
		case "/nocontent":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, state
}

func mustHTTPSink(t *testing.T, config map[string]any) flow.Block {
	t.Helper()
	b, err := blocks.HTTPSink("snk", config, flow.Deps{})
	if err != nil {
		t.Fatalf("HTTPSink: %v", err)
	}
	return b
}

func TestHTTPSink_PayloadKinds(t *testing.T) {
	srv, state := newSinkServer(t)
	sink := mustHTTPSink(t, map[string]any{"url": srv.URL + "/ingest"})

	payloads := []*flow.Message{
		{
			Key:  "bin-bare",
			Type: flow.TypeBinary,
			Data: []byte{0x00, 0xff},
		},
		{
			Key:     "bin-csv",
			Type:    flow.TypeBinary,
			Subtype: "text/csv",
			Data:    []byte("a,b\n1,2\n"),
		},
		{
			Key:  "str",
			Type: flow.TypeString,
			Data: "héllo",
		},
		{
			Key:  "map",
			Type: flow.TypeMap,
			FieldTypes: map[string]flow.DataType{
				"a": flow.TypeInteger,
				"b": flow.TypeString,
			},
			Data: map[string]any{"a": int64(42), "b": "hello"},
		},
	}
	wantBodies := [][]byte{
		{0x00, 0xff},
		[]byte("a,b\n1,2\n"),
		[]byte("héllo"),
		nil, // asserted as decoded JSON below
	}
	wantContentTypes := []string{
		"application/octet-stream",
		"text/csv",
		"text/plain; charset=utf-8",
		"application/json",
	}

	for i, msg := range payloads {
		if _, err := sink.Process(msg); err != nil {
			t.Fatalf("payload %d (%s): Process: %v", i, msg.Key, err)
		}
	}

	got := state.snapshot()
	if len(got) != len(payloads) {
		t.Fatalf("server saw %d requests, want %d", len(got), len(payloads))
	}
	for i, req := range got {
		if req.method != http.MethodPost {
			t.Errorf("request %d method = %q, want POST", i, req.method)
		}
		if req.path != "/ingest" {
			t.Errorf("request %d path = %q, want /ingest", i, req.path)
		}
		if req.contentType != wantContentTypes[i] {
			t.Errorf("request %d (%s) Content-Type = %q, want %q",
				i, payloads[i].Key, req.contentType, wantContentTypes[i])
		}
		if wantBodies[i] != nil && string(req.body) != string(wantBodies[i]) {
			t.Errorf("request %d (%s) body = %q, want %q",
				i, payloads[i].Key, req.body, wantBodies[i])
		}
	}

	var decoded map[string]any
	if err := json.Unmarshal(got[3].body, &decoded); err != nil {
		t.Fatalf("map payload body %q is not JSON: %v", got[3].body, err)
	}
	if decoded["a"] != float64(42) || decoded["b"] != "hello" {
		t.Errorf("map payload JSON = %v, want a=42 b=hello", decoded)
	}
}

func TestHTTPSink_OperatorContentTypeWins(t *testing.T) {
	srv, state := newSinkServer(t)
	sink := mustHTTPSink(t, map[string]any{
		"url":     srv.URL + "/ingest",
		"headers": map[string]any{"Content-Type": "application/custom"},
	})

	msg := &flow.Message{Key: "k", Type: flow.TypeBinary, Data: []byte("raw")}
	if _, err := sink.Process(msg); err != nil {
		t.Fatalf("Process: %v", err)
	}

	got := state.snapshot()
	if len(got) != 1 {
		t.Fatalf("server saw %d requests, want 1", len(got))
	}
	if ct := got[0].contentType; ct != "application/custom" {
		t.Errorf("Content-Type = %q, want operator header to win (application/custom)", ct)
	}
}

func TestHTTPSink_CustomHeaderArrives(t *testing.T) {
	srv, state := newSinkServer(t)
	sink := mustHTTPSink(t, map[string]any{
		"url":     srv.URL + "/ingest",
		"headers": map[string]any{"X-Test": "yes"},
	})

	msg := &flow.Message{Key: "k", Type: flow.TypeString, Data: "x"}
	if _, err := sink.Process(msg); err != nil {
		t.Fatalf("Process: %v", err)
	}

	got := state.snapshot()
	if len(got) != 1 {
		t.Fatalf("server saw %d requests, want 1", len(got))
	}
	if v := got[0].header.Get("X-Test"); v != "yes" {
		t.Errorf("X-Test = %q, want yes", v)
	}
}

func TestHTTPSink_StatusHandling(t *testing.T) {
	srv, _ := newSinkServer(t)
	sink := mustHTTPSink(t, map[string]any{"url": srv.URL + "/fail"})

	msg := &flow.Message{Key: "k", Type: flow.TypeString, Data: "x"}
	_, err := sink.Process(msg)
	if err == nil || err.Error() != "http_sink: POST "+srv.URL+"/fail: status 500" {
		t.Fatalf("err = %v, want %q", err, "http_sink: POST "+srv.URL+"/fail: status 500")
	}

	okSink := mustHTTPSink(t, map[string]any{"url": srv.URL + "/nocontent"})
	if _, err := okSink.Process(msg); err != nil {
		t.Fatalf("204 No Content should succeed, got %v", err)
	}
}

func TestHTTPSink_NilMessageIsNoop(t *testing.T) {
	srv, state := newSinkServer(t)
	sink := mustHTTPSink(t, map[string]any{"url": srv.URL + "/ingest"})
	if _, err := sink.Process(nil); err != nil {
		t.Fatalf("Process(nil): %v", err)
	}
	if got := state.snapshot(); len(got) != 0 {
		t.Errorf("nil message produced %d requests, want 0", len(got))
	}
}

func TestHTTPSink_TransportErrorWrapped(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	url := srv.URL + "/ingest"
	srv.Close()

	sink := mustHTTPSink(t, map[string]any{"url": url})
	msg := &flow.Message{Key: "k", Type: flow.TypeString, Data: "x"}
	_, err := sink.Process(msg)
	if err == nil || !strings.HasPrefix(err.Error(), "http_sink: POST "+url+": ") {
		t.Fatalf("err = %v, want wrapped \"http_sink: POST %s: ...\" transport error", err, url)
	}
}

func TestHTTPSink_RequestTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(300 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	sink := mustHTTPSink(t, map[string]any{
		"url":        srv.URL + "/ingest",
		"timeout_ms": 50,
	})

	start := time.Now()
	_, err := sink.Process(&flow.Message{Key: "k", Type: flow.TypeString, Data: "x"})
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("Process succeeded despite server exceeding timeout_ms")
	}
	if elapsed > 250*time.Millisecond {
		t.Errorf("Process took %v, want timeout_ms to bound the request", elapsed)
	}
}

func TestHTTPSink_ConstructRejections(t *testing.T) {
	tests := []struct {
		name string
		bad  map[string]any
		good map[string]any
		want string
	}{
		{
			name: "missing url",
			bad:  map[string]any{},
			good: map[string]any{"url": "http://localhost/ingest"},
			want: `http_sink: url is required`,
		},
		{
			name: "empty url",
			bad:  map[string]any{"url": ""},
			good: map[string]any{"url": "http://localhost/ingest"},
			want: `http_sink: url is required`,
		},
		{
			name: "zero timeout",
			bad:  map[string]any{"url": "http://localhost/ingest", "timeout_ms": 0},
			good: map[string]any{"url": "http://localhost/ingest", "timeout_ms": 1},
			want: `http_sink: timeout_ms must be positive`,
		},
		{
			name: "negative timeout",
			bad:  map[string]any{"url": "http://localhost/ingest", "timeout_ms": -5},
			good: map[string]any{"url": "http://localhost/ingest", "timeout_ms": 5},
			want: `http_sink: timeout_ms must be positive`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name+"/reject", func(t *testing.T) {
			_, err := blocks.HTTPSink("snk", tt.bad, flow.Deps{})
			if err == nil || err.Error() != tt.want {
				t.Fatalf("err = %v, want %q", err, tt.want)
			}
		})
		t.Run(tt.name+"/accept", func(t *testing.T) {
			if _, err := blocks.HTTPSink("snk", tt.good, flow.Deps{}); err != nil {
				t.Fatalf("twin rejected: %v", err)
			}
		})
	}
}

func TestHTTPBlocks_RegistrationAndInfo(t *testing.T) {
	reg := blocks.DefaultRegistry()
	if !reg.Has(blocks.TypeHTTPSource) {
		t.Errorf("registry missing %q", blocks.TypeHTTPSource)
	}
	if !reg.Has(blocks.TypeHTTPSink) {
		t.Errorf("registry missing %q", blocks.TypeHTTPSink)
	}

	info, ok := blocks.LookupInfo(blocks.TypeHTTPSource)
	if !ok {
		t.Fatalf("LookupInfo missing docs for %q", blocks.TypeHTTPSource)
	}
	if info.Role != blocks.RoleSource {
		t.Errorf("%q role = %q, want source", blocks.TypeHTTPSource, info.Role)
	}

	info, ok = blocks.LookupInfo(blocks.TypeHTTPSink)
	if !ok {
		t.Fatalf("LookupInfo missing docs for %q", blocks.TypeHTTPSink)
	}
	if info.Role != blocks.RoleSink {
		t.Errorf("%q role = %q, want sink", blocks.TypeHTTPSink, info.Role)
	}
}
