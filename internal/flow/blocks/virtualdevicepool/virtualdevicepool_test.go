package virtualdevicepool_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool"
)

// ingestCall records one invocation of the fake engine ingest.
type ingestCall struct {
	realm    string
	deviceID string
	iface    string
	path     string
	payload  string
	ts       *time.Time
}

// fakeIngest is a capturing flow.Deps.Ingest: calls are appended to a
// mutex-guarded slice so tests may assert after processing.
type fakeIngest struct {
	mu    sync.Mutex
	calls []ingestCall
	err   error
}

func (f *fakeIngest) ingest(_ context.Context, realm, deviceID, ifaceName, path string, payload json.RawMessage, ts *time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return f.err
	}
	f.calls = append(f.calls, ingestCall{
		realm:    realm,
		deviceID: deviceID,
		iface:    ifaceName,
		path:     path,
		payload:  string(payload),
		ts:       ts,
	})
	return nil
}

func (f *fakeIngest) recorded() []ingestCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]ingestCall(nil), f.calls...)
}

func newBlock(t *testing.T, name string, config map[string]any, deps flow.Deps) flow.Block {
	t.Helper()
	b, err := virtualdevicepool.Constructor(name, config, deps)
	if err != nil {
		t.Fatalf("Constructor: %v", err)
	}
	return b
}

func poolConfig(devices ...string) map[string]any {
	return map[string]any{"devices": devices}
}

// registerCall records one invocation of the fake pairing register.
type registerCall struct {
	realm    string
	deviceID string
}

// fakeRegister is a capturing flow.Deps.Register: calls are appended to a
// mutex-guarded slice so tests may assert after processing.
type fakeRegister struct {
	mu    sync.Mutex
	calls []registerCall
	err   error
}

func (f *fakeRegister) register(_ context.Context, realm, deviceID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, registerCall{realm: realm, deviceID: deviceID})
	return f.err
}

func (f *fakeRegister) recorded() []registerCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]registerCall(nil), f.calls...)
}

func TestVirtualDevicePool_AcceptsAndAddressesTarget(t *testing.T) {
	fi := &fakeIngest{}
	deps := flow.Deps{Realm: "testrealm", Ingest: fi.ingest}
	b := newBlock(t, "vdp", poolConfig("dev1", "dev2"), deps)

	out, err := b.Process(&flow.Message{
		Key:       "dev1/com.example.Temp/value",
		Type:      flow.TypeString,
		Data:      "21.5",
		Timestamp: 1700000000000000,
	})
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if out != nil {
		t.Errorf("Process outputs = %v, want nil", out)
	}
	calls := fi.recorded()
	if len(calls) != 1 {
		t.Fatalf("ingest calls = %d, want exactly 1 (%v)", len(calls), calls)
	}
	c := calls[0]
	if c.realm != "testrealm" {
		t.Errorf("realm = %q, want testrealm", c.realm)
	}
	if c.deviceID != "dev1" {
		t.Errorf("deviceID = %q, want dev1", c.deviceID)
	}
	if c.iface != "com.example.Temp" {
		t.Errorf("iface = %q, want com.example.Temp", c.iface)
	}
	if c.path != "/value" {
		t.Errorf("path = %q, want /value", c.path)
	}
	if c.payload != `"21.5"` {
		t.Errorf("payload = %s, want %q", c.payload, `"21.5"`)
	}
	wantTS := time.UnixMicro(1700000000000000).UTC()
	if c.ts == nil || !c.ts.Equal(wantTS) {
		t.Errorf("ts = %v, want %v", c.ts, wantTS)
	}
}

func TestVirtualDevicePool_DeepPath(t *testing.T) {
	fi := &fakeIngest{}
	deps := flow.Deps{Realm: "r", Ingest: fi.ingest}
	b := newBlock(t, "vdp", poolConfig("dev1"), deps)

	if _, err := b.Process(&flow.Message{
		Key:       "dev1/com.example.Gps/coords/lat",
		Type:      flow.TypeReal,
		Data:      45.5,
		Timestamp: 1700000000000000,
	}); err != nil {
		t.Fatalf("Process: %v", err)
	}
	calls := fi.recorded()
	if len(calls) != 1 {
		t.Fatalf("ingest calls = %d, want 1", len(calls))
	}
	if calls[0].path != "/coords/lat" {
		t.Errorf("path = %q, want /coords/lat", calls[0].path)
	}
}

func TestVirtualDevicePool_ZeroTimestampMeansNilTS(t *testing.T) {
	fi := &fakeIngest{}
	deps := flow.Deps{Realm: "r", Ingest: fi.ingest}
	b := newBlock(t, "vdp", poolConfig("dev1"), deps)

	if _, err := b.Process(&flow.Message{
		Key:  "dev1/com.example.Temp/value",
		Type: flow.TypeString,
		Data: "x",
	}); err != nil {
		t.Fatalf("Process: %v", err)
	}
	calls := fi.recorded()
	if len(calls) != 1 {
		t.Fatalf("ingest calls = %d, want 1", len(calls))
	}
	if calls[0].ts != nil {
		t.Errorf("ts = %v, want nil", calls[0].ts)
	}
}

func TestVirtualDevicePool_DropsBadMessages(t *testing.T) {
	rows := []struct {
		name string
		msg  *flow.Message
	}{
		{"unknown device", &flow.Message{Key: "stranger/com.example.X/y", Type: flow.TypeString, Data: "x"}},
		{"device only", &flow.Message{Key: "dev1", Type: flow.TypeString, Data: "x"}},
		{"device and slash", &flow.Message{Key: "dev1/", Type: flow.TypeString, Data: "x"}},
		{"empty device segment", &flow.Message{Key: "/com.example.X/y", Type: flow.TypeString, Data: "x"}},
		{"empty key", &flow.Message{Key: "", Type: flow.TypeString, Data: "x"}},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			fi := &fakeIngest{}
			deps := flow.Deps{Realm: "r", Ingest: fi.ingest}
			b := newBlock(t, "vdp", poolConfig("dev1"), deps)

			out, err := b.Process(row.msg)
			if err != nil {
				t.Fatalf("Process err = %v, want nil (drop)", err)
			}
			if out != nil {
				t.Errorf("Process outputs = %v, want nil", out)
			}
			if calls := fi.recorded(); len(calls) != 0 {
				t.Errorf("ingest calls = %d, want 0", len(calls))
			}
		})
	}
}

func TestVirtualDevicePool_NilMessage(t *testing.T) {
	fi := &fakeIngest{}
	deps := flow.Deps{Realm: "r", Ingest: fi.ingest}
	b := newBlock(t, "vdp", poolConfig("dev1"), deps)

	out, err := b.Process(nil)
	if err != nil {
		t.Fatalf("Process(nil) err = %v, want nil", err)
	}
	if out != nil {
		t.Errorf("outputs = %v, want nil", out)
	}
	if calls := fi.recorded(); len(calls) != 0 {
		t.Errorf("ingest calls = %d, want 0", len(calls))
	}
}

func TestVirtualDevicePool_IngestErrorIsWrapped(t *testing.T) {
	sentinelErr := errors.New("engine ingest unavailable")
	fi := &fakeIngest{err: sentinelErr}
	deps := flow.Deps{Realm: "r", Ingest: fi.ingest}
	b := newBlock(t, "vdp", poolConfig("dev1"), deps)

	_, err := b.Process(&flow.Message{
		Key:       "dev1/com.example.Temp/value",
		Type:      flow.TypeString,
		Data:      "21.5",
		Timestamp: 1700000000000000,
	})
	if err == nil {
		t.Fatal("Process err = nil, want wrapped sentinel")
	}
	if !errors.Is(err, sentinelErr) {
		t.Fatalf("err = %v, want it to wrap %v", err, sentinelErr)
	}
	wantPrefix := "virtual_device_pool vdp:"
	if got := err.Error(); len(got) < len(wantPrefix) || got[:len(wantPrefix)] != wantPrefix {
		t.Errorf("err = %q, want prefix %q", got, wantPrefix)
	}
}

func TestVirtualDevicePool_ConfigValidation(t *testing.T) {
	rows := []struct {
		name    string
		config  map[string]any
		deps    flow.Deps
		wantErr string
	}{
		{
			name:    "nil ingest",
			config:  poolConfig("ok"),
			deps:    flow.Deps{},
			wantErr: "virtual_device_pool requires the engine ingest dependency",
		},
		{
			name:    "missing devices",
			config:  map[string]any{},
			deps:    flow.Deps{Ingest: func(context.Context, string, string, string, string, json.RawMessage, *time.Time) error { return nil }},
			wantErr: "virtual_device_pool: devices is required",
		},
		{
			name:    "empty devices",
			config:  map[string]any{"devices": []string{}},
			deps:    flow.Deps{Ingest: func(context.Context, string, string, string, string, json.RawMessage, *time.Time) error { return nil }},
			wantErr: "virtual_device_pool: devices must not be empty",
		},
		{
			name:    "empty entry",
			config:  map[string]any{"devices": []any{"ok", ""}},
			deps:    flow.Deps{Ingest: func(context.Context, string, string, string, string, json.RawMessage, *time.Time) error { return nil }},
			wantErr: "virtual_device_pool: devices[1] must be a non-empty string",
		},
		{
			name:    "acceptance twin json array",
			config:  map[string]any{"devices": []any{"ok"}},
			deps:    flow.Deps{Ingest: func(context.Context, string, string, string, string, json.RawMessage, *time.Time) error { return nil }},
			wantErr: "",
		},
		{
			name:    "auto_register not boolean",
			config:  map[string]any{"auto_register": "yes"},
			deps:    flow.Deps{Ingest: func(context.Context, string, string, string, string, json.RawMessage, *time.Time) error { return nil }},
			wantErr: "virtual_device_pool: auto_register must be a boolean",
		},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			b, err := virtualdevicepool.Constructor("vdp", row.config, row.deps)
			if row.wantErr != "" {
				if err == nil || err.Error() != row.wantErr {
					t.Fatalf("err = %v, want %q", err, row.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Constructor: %v", err)
			}
			if b.Name() != "vdp" {
				t.Errorf("Name = %q, want vdp", b.Name())
			}
			if _, err := b.Process(&flow.Message{
				Key:  "ok/com.example.X/y",
				Type: flow.TypeString,
				Data: "x",
			}); err != nil {
				t.Errorf("working construction Process: %v", err)
			}
		})
	}
}

func TestVirtualDevicePool_StaticInvariantNoRegisterWithoutAutoRegister(t *testing.T) {
	fi := &fakeIngest{}
	fr := &fakeRegister{}
	deps := flow.Deps{Realm: "r", Ingest: fi.ingest, Register: fr.register}
	b := newBlock(t, "vdp", poolConfig("dev1"), deps)

	out, err := b.Process(&flow.Message{Key: "stranger/com.example.X/y", Type: flow.TypeString, Data: "x"})
	if err != nil {
		t.Fatalf("Process err = %v, want nil (drop)", err)
	}
	if out != nil {
		t.Errorf("outputs = %v, want nil", out)
	}
	if calls := fi.recorded(); len(calls) != 0 {
		t.Errorf("ingest calls = %d, want 0", len(calls))
	}
	if calls := fr.recorded(); len(calls) != 0 {
		t.Errorf("register calls = %d, want 0 (auto_register off)", len(calls))
	}
}

func TestVirtualDevicePool_RegisterDepRequiredWithAutoRegister(t *testing.T) {
	deps := flow.Deps{Realm: "r", Ingest: func(context.Context, string, string, string, string, json.RawMessage, *time.Time) error { return nil }}
	_, err := virtualdevicepool.Constructor("vdp", map[string]any{"auto_register": true}, deps)
	if err == nil {
		t.Fatal("Constructor err = nil, want register-dependency error")
	}
	if !strings.Contains(err.Error(), "register") {
		t.Fatalf("err = %v, want it to mention register", err)
	}
}

func TestVirtualDevicePool_FirstSeenRegisters(t *testing.T) {
	fi := &fakeIngest{}
	fr := &fakeRegister{}
	deps := flow.Deps{Realm: "testrealm", Ingest: fi.ingest, Register: fr.register}
	b := newBlock(t, "vdp", map[string]any{"auto_register": true}, deps)

	if _, err := b.Process(&flow.Message{
		Key:       "newdev/com.example.Temp/value",
		Type:      flow.TypeString,
		Data:      "21.5",
		Timestamp: 1700000000000000,
	}); err != nil {
		t.Fatalf("first Process: %v", err)
	}
	regCalls := fr.recorded()
	if len(regCalls) != 1 {
		t.Fatalf("register calls = %d, want exactly 1 (%v)", len(regCalls), regCalls)
	}
	if regCalls[0].realm != "testrealm" || regCalls[0].deviceID != "newdev" {
		t.Errorf("register call = %+v, want realm testrealm device newdev", regCalls[0])
	}
	ingests := fi.recorded()
	if len(ingests) != 1 {
		t.Fatalf("ingest calls after first message = %d, want 1", len(ingests))
	}
	if ingests[0].deviceID != "newdev" || ingests[0].payload != `"21.5"` {
		t.Errorf("ingest call = %+v, want newdev with %q", ingests[0], `"21.5"`)
	}

	if _, err := b.Process(&flow.Message{
		Key:  "newdev/com.example.Temp/value",
		Type: flow.TypeString,
		Data: "22.5",
	}); err != nil {
		t.Fatalf("second Process: %v", err)
	}
	if calls := fr.recorded(); len(calls) != 1 {
		t.Errorf("register calls after second message = %d, want still 1 (pool hit)", len(calls))
	}
	if calls := fi.recorded(); len(calls) != 2 {
		t.Errorf("ingest calls after second message = %d, want 2", len(calls))
	}
}

func TestVirtualDevicePool_AlreadyRegisteredDropsWithoutCaching(t *testing.T) {
	fi := &fakeIngest{}
	fr := &fakeRegister{err: fmt.Errorf("%w: x", flow.ErrVirtualDeviceRegistered)}
	deps := flow.Deps{Realm: "r", Ingest: fi.ingest, Register: fr.register}
	b := newBlock(t, "vdp", map[string]any{"auto_register": true}, deps)

	out, err := b.Process(&flow.Message{Key: "newdev/com.example.X/y", Type: flow.TypeString, Data: "x"})
	if err != nil {
		t.Fatalf("Process err = %v, want nil (drop)", err)
	}
	if out != nil {
		t.Errorf("outputs = %v, want nil", out)
	}
	if calls := fi.recorded(); len(calls) != 0 {
		t.Errorf("ingest calls = %d, want 0", len(calls))
	}

	if _, err := b.Process(&flow.Message{Key: "newdev/com.example.X/y", Type: flow.TypeString, Data: "x"}); err != nil {
		t.Fatalf("second Process err = %v, want nil (drop)", err)
	}
	if calls := fr.recorded(); len(calls) != 2 {
		t.Errorf("register calls = %d, want 2 (no negative caching)", len(calls))
	}
	if calls := fi.recorded(); len(calls) != 0 {
		t.Errorf("ingest calls after second message = %d, want 0", len(calls))
	}
}

func TestVirtualDevicePool_InfraRegisterErrorSurfaces(t *testing.T) {
	fi := &fakeIngest{}
	dbDown := errors.New("db down")
	fr := &fakeRegister{err: dbDown}
	deps := flow.Deps{Realm: "r", Ingest: fi.ingest, Register: fr.register}
	b := newBlock(t, "vdp", map[string]any{"auto_register": true}, deps)

	_, err := b.Process(&flow.Message{Key: "newdev/com.example.X/y", Type: flow.TypeString, Data: "x"})
	if err == nil {
		t.Fatal("Process err = nil, want wrapped register error")
	}
	if !errors.Is(err, dbDown) {
		t.Fatalf("err = %v, want it to wrap %v (db down)", err, dbDown)
	}
	if !strings.Contains(err.Error(), "db down") {
		t.Errorf("err = %q, want it to contain \"db down\"", err.Error())
	}
	if calls := fi.recorded(); len(calls) != 0 {
		t.Errorf("ingest calls = %d, want 0", len(calls))
	}
}

func TestVirtualDevicePool_ConcurrentFirstSeenSingleRegister(t *testing.T) {
	fi := &fakeIngest{}
	fr := &fakeRegister{}
	deps := flow.Deps{Realm: "r", Ingest: fi.ingest, Register: fr.register}
	b := newBlock(t, "vdp", map[string]any{"auto_register": true}, deps)

	const n = 8
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if _, err := b.Process(&flow.Message{
				Key:  fmt.Sprintf("hotdev/com.example.N/value-%d", i),
				Type: flow.TypeString,
				Data: i,
			}); err != nil {
				t.Errorf("Process(%d): %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	if calls := fr.recorded(); len(calls) != 1 {
		t.Errorf("register calls = %d, want exactly 1 (%v)", len(calls), calls)
	}
	calls := fi.recorded()
	if len(calls) != n {
		t.Fatalf("ingest calls = %d, want %d", len(calls), n)
	}
	got := make(map[string]bool, n)
	for _, c := range calls {
		got[c.payload] = true
	}
	for i := 0; i < n; i++ {
		want := strconv.Itoa(i)
		if !got[want] {
			t.Errorf("payload %s missing from ingested set %v", want, got)
		}
	}
}
