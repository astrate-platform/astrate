package engine

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	promtest "github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// compiledFixture parses and compiles one M1 fixture with synthetic IDs.
func compiledFixture(t testing.TB, ifaceID int64, name string) *interfaceschema.CompiledInterface {
	t.Helper()
	si := fixtureStored(t, 1, ifaceID, name)
	iface, err := interfaceschema.ParseInterface(si.Definition)
	if err != nil {
		t.Fatalf("parsing %s: %v", name, err)
	}
	ci, err := interfaceschema.Compile(iface, si)
	if err != nil {
		t.Fatalf("compiling %s: %v", name, err)
	}
	return ci
}

// TestEncodeValueJSON pins the canonical jsonb rendering (docs/DESIGN.md
// §2.3): base64 blobs, RFC 3339 milli datetimes, decimal-string
// longintegers beyond 2^53.
func TestEncodeValueJSON(t *testing.T) {
	ts := time.Date(2026, 6, 1, 12, 0, 0, 500_000_000, time.UTC)
	cases := []struct {
		name string
		v    payload.Value
		want string
	}{
		{"double", 1.5, `1.5`},
		{"integer", int32(7), `7`},
		{"longinteger small", int64(5), `5`},
		{"longinteger negative", int64(-9007199254740991), `-9007199254740991`},
		{"longinteger beyond 2^53", int64(1) << 60, `"1152921504606846976"`},
		{"boolean", true, `true`},
		{"string", "héllo", `"héllo"`},
		{"binaryblob", []byte{0xde, 0xad}, `"3q0="`},
		{"datetime", ts, `"2026-06-01T12:00:00.500Z"`},
		{"doublearray", []float64{1.5, -2.5}, `[1.5,-2.5]`},
		{"integerarray", []int32{1, 2}, `[1,2]`},
		{"longintegerarray big", []int64{1 << 60}, `["1152921504606846976"]`},
		{"booleanarray", []bool{true, false}, `[true,false]`},
		{"stringarray", []string{"a", "b"}, `["a","b"]`},
		{"binaryblobarray", [][]byte{{0xde, 0xad}}, `["3q0="]`},
		{"datetimearray", []time.Time{ts}, `["2026-06-01T12:00:00.500Z"]`},
		{"object document", map[string]payload.Value{"b": "x", "a": 1.5}, `{"a":1.5,"b":"x"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := encodeValueJSON(tc.v)
			if err != nil {
				t.Fatalf("encodeValueJSON: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("encodeValueJSON(%v) = %s, want %s", tc.v, got, tc.want)
			}
		})
	}

	if _, err := encodeValueJSON(struct{}{}); err == nil {
		t.Error("unsupported Go type encoded without error")
	}
	if _, err := encodeValueJSON(map[string]payload.Value{"k": struct{}{}}); err == nil {
		t.Error("unsupported nested value encoded without error")
	}
}

// TestIndividualRowColumns: each scalar type lands in exactly its typed
// column; arrays land pre-encoded in value_array (docs/DESIGN.md §2.4).
func TestIndividualRowColumns(t *testing.T) {
	scalars := compiledFixture(t, 10, "com.astrate.test.AllScalarTypes.json")
	arrays := compiledFixture(t, 16, "com.astrate.test.AllArrayTypes.json")
	ts := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)

	op := func(ci *interfaceschema.CompiledInterface, path string, v payload.Value) *PersistOp {
		t.Helper()
		m, ok := ci.Trie.Match(path)
		if !ok {
			t.Fatalf("fixture path %s does not match", path)
		}
		return &PersistOp{
			Kind: OpIndividual, RealmID: 1, DeviceID: devAlpha, Interface: ci,
			Mapping: m, Path: path, Value: v, TS: ts, ReceptionTS: ts,
		}
	}

	cases := []struct {
		path  string
		ci    *interfaceschema.CompiledInterface
		v     payload.Value
		check func(t *testing.T, r *store.IndividualRow)
	}{
		{"/double", scalars, 1.5, func(t *testing.T, r *store.IndividualRow) {
			if r.ValueDouble == nil || *r.ValueDouble != 1.5 {
				t.Errorf("value_double = %v", r.ValueDouble)
			}
		}},
		{"/integer", scalars, int32(7), func(t *testing.T, r *store.IndividualRow) {
			if r.ValueInteger == nil || *r.ValueInteger != 7 {
				t.Errorf("value_integer = %v", r.ValueInteger)
			}
		}},
		{"/longinteger", scalars, int64(1) << 60, func(t *testing.T, r *store.IndividualRow) {
			if r.ValueLonginteger == nil || *r.ValueLonginteger != 1<<60 {
				t.Errorf("value_longinteger = %v", r.ValueLonginteger)
			}
		}},
		{"/boolean", scalars, true, func(t *testing.T, r *store.IndividualRow) {
			if r.ValueBoolean == nil || !*r.ValueBoolean {
				t.Errorf("value_boolean = %v", r.ValueBoolean)
			}
		}},
		{"/string", scalars, "x", func(t *testing.T, r *store.IndividualRow) {
			if r.ValueString == nil || *r.ValueString != "x" {
				t.Errorf("value_string = %v", r.ValueString)
			}
		}},
		{"/binaryblob", scalars, []byte{0xde}, func(t *testing.T, r *store.IndividualRow) {
			if len(r.ValueBinaryblob) != 1 || r.ValueBinaryblob[0] != 0xde {
				t.Errorf("value_binaryblob = %v", r.ValueBinaryblob)
			}
		}},
		{"/datetime", scalars, ts, func(t *testing.T, r *store.IndividualRow) {
			if r.ValueDatetime == nil || !r.ValueDatetime.Equal(ts) {
				t.Errorf("value_datetime = %v", r.ValueDatetime)
			}
		}},
		{"/g/doublearray", arrays, []float64{1.5, 2.5}, func(t *testing.T, r *store.IndividualRow) {
			if string(r.ValueArray) != `[1.5,2.5]` {
				t.Errorf("value_array = %s", r.ValueArray)
			}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			o := op(tc.ci, tc.path, tc.v)
			row, err := individualRow(o)
			if err != nil {
				t.Fatalf("individualRow: %v", err)
			}
			set := 0
			for _, p := range []bool{
				row.ValueDouble != nil, row.ValueInteger != nil, row.ValueLonginteger != nil,
				row.ValueBoolean != nil, row.ValueString != nil, row.ValueBinaryblob != nil,
				row.ValueDatetime != nil, row.ValueArray != nil,
			} {
				if p {
					set++
				}
			}
			if set != 1 {
				t.Errorf("%d value columns set, want exactly 1", set)
			}
			if row.EndpointID != o.Mapping.EndpointID || row.InterfaceID != o.Interface.ID || row.Path != tc.path {
				t.Errorf("row addressing: %+v", row)
			}
			tc.check(t, row)
		})
	}

	// A value/type mismatch (impossible for pipeline output) is an error.
	bad := op(scalars, "/double", "not a double")
	if _, err := individualRow(bad); err == nil {
		t.Error("mismatched value converted without error")
	}
}

// TestObjectAndPropertyRows covers the remaining conversions.
func TestObjectAndPropertyRows(t *testing.T) {
	flat := compiledFixture(t, 11, "com.astrate.test.ObjectFlat.json")
	props := compiledFixture(t, 14, "com.astrate.test.PropertyArrays.json")
	ts := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)

	obj := &PersistOp{
		Kind: OpObject, RealmID: 1, DeviceID: devAlpha, Interface: flat, Path: "",
		Value: map[string]payload.Value{"latitude": 45.07, "longitude": 7.69},
		TS:    ts, ReceptionTS: ts,
	}
	row, err := objectRow(obj)
	if err != nil {
		t.Fatalf("objectRow: %v", err)
	}
	if string(row.Value) != `{"latitude":45.07,"longitude":7.69}` {
		t.Errorf("object value = %s", row.Value)
	}
	if row.InterfaceID != flat.ID || row.Path != "" {
		t.Errorf("object row addressing: %+v", row)
	}
	obj.Value = 1.5
	if _, err := objectRow(obj); err == nil {
		t.Error("non-document object value converted without error")
	}

	m, _ := props.Trie.Match("/config/thresholds")
	set := &PersistOp{
		Kind: OpPropertySet, RealmID: 1, DeviceID: devAlpha, Interface: props,
		Mapping: m, Path: "/config/thresholds", Value: []float64{1.5},
		TS: ts, ReceptionTS: ts,
	}
	p, err := propertyRow(set)
	if err != nil {
		t.Fatalf("propertyRow: %v", err)
	}
	if string(p.Value) != `[1.5]` || p.ValueType != interfaceschema.DoubleArray ||
		p.EndpointID != m.EndpointID || p.InterfaceID != props.ID {
		t.Errorf("property row: %+v", p)
	}
}

// TestFlushAckAfterCommit: nothing is acknowledged before the store calls
// succeed; everything (including in-order property ops) is afterwards.
func TestFlushAckAfterCommit(t *testing.T) {
	ctx := context.Background()
	rig, fs := newPipelineRig(t, Config{BatchMaxRows: 1000})

	acks := make([]*ackCounter, 4)
	for i := range acks {
		acks[i] = &ackCounter{}
	}
	rig.handle(deviceMsg("com.astrate.test.Minimal", "/value", 1, enc(t, 1.0, nil, payload.FormatBSON), acks[0]))
	rig.handle(deviceMsg("com.astrate.test.PropertyArrays", "/config/thresholds", 2, []byte(`{"v":[1.0]}`), acks[1]))
	rig.handle(deviceMsg("com.astrate.test.PropertyArrays", "/config/thresholds", 2, []byte(`{"v":[2.0]}`), acks[2]))
	rig.handle(deviceMsg("com.astrate.test.PropertyArrays", "/config/thresholds", 2, nil, acks[3]))

	for i, a := range acks {
		if a.acked() {
			t.Fatalf("message %d acknowledged before flush", i)
		}
	}
	rig.sh.batch.flush(ctx)
	for i, a := range acks {
		if !a.acked() {
			t.Errorf("message %d not acknowledged after flush", i)
		}
	}

	if len(fs.upserts) != 2 || string(fs.upserts[0].Value) != `[1]` || string(fs.upserts[1].Value) != `[2]` {
		t.Errorf("property upserts out of order: %+v", fs.upserts)
	}
	if len(fs.unsets) != 1 || fs.unsets[0].path != "/config/thresholds" || fs.unsets[0].interfaceID != ifacePropArrays {
		t.Errorf("property unsets: %+v", fs.unsets)
	}
	if rows := fs.individualRows(); len(rows) != 1 || *rows[0].ValueDouble != 1.0 {
		t.Errorf("datastream rows: %+v", rows)
	}
	for kind, want := range map[OpKind]float64{
		OpIndividual: 1, OpPropertySet: 2, OpPropertyUnset: 1,
	} {
		if got := promtest.ToFloat64(rig.e.met.persistOps.WithLabelValues(kind.String())); got != want {
			t.Errorf("persist ops[%s] = %v, want %v", kind, got, want)
		}
	}
}

// TestFlushRetryTransient: transient store failures park and retry; acks
// fire exactly once, after the eventually-successful commit.
func TestFlushRetryTransient(t *testing.T) {
	ctx := context.Background()
	rig, fs := newPipelineRig(t, Config{BatchMaxRows: 1000})
	fs.mu.Lock()
	fs.appendErrs = []error{errTransient, errTransient}
	fs.mu.Unlock()

	ack := &ackCounter{}
	rig.handle(deviceMsg("com.astrate.test.Minimal", "/value", 1, enc(t, 1.0, nil, payload.FormatBSON), ack))
	rig.sh.batch.flush(ctx)

	if ack.n.Load() != 1 {
		t.Errorf("ack called %d times, want exactly 1", ack.n.Load())
	}
	if len(fs.individualRows()) != 1 {
		t.Errorf("%d rows after retries, want 1", len(fs.individualRows()))
	}
	if got := promtest.ToFloat64(rig.e.met.flushRetries); got != 2 {
		t.Errorf("flush retry counter = %v, want 2", got)
	}
}

// TestFlushPoisonedBatch: an integrity-class database rejection splits the
// batch into per-op commits instead of wedging the shard forever.
func TestFlushPoisonedBatch(t *testing.T) {
	ctx := context.Background()
	pgErr := &pgconn.PgError{Code: "23503"}

	t.Run("SplitRecovers", func(t *testing.T) {
		rig, fs := newPipelineRig(t, Config{BatchMaxRows: 1000})
		fs.mu.Lock()
		fs.appendErrs = []error{pgErr}
		fs.mu.Unlock()

		a1, a2 := &ackCounter{}, &ackCounter{}
		rig.handle(deviceMsg("com.astrate.test.Minimal", "/value", 1, enc(t, 1.0, nil, payload.FormatBSON), a1))
		rig.handle(deviceMsg("com.astrate.test.Minimal", "/value", 1, enc(t, 2.0, nil, payload.FormatBSON), a2))
		rig.sh.batch.flush(ctx)

		if !a1.acked() || !a2.acked() {
			t.Error("split-committed ops not acknowledged")
		}
		if got := len(fs.individualRows()); got != 2 {
			t.Errorf("%d rows after split, want 2", got)
		}
	})

	t.Run("SinglePoisonedOpDropped", func(t *testing.T) {
		rig, fs := newPipelineRig(t, Config{BatchMaxRows: 1000})
		fs.mu.Lock()
		fs.appendErrs = []error{pgErr, pgErr, pgErr}
		fs.mu.Unlock()

		a1, a2 := &ackCounter{}, &ackCounter{}
		rig.handle(deviceMsg("com.astrate.test.Minimal", "/value", 1, enc(t, 1.0, nil, payload.FormatBSON), a1))
		rig.handle(deviceMsg("com.astrate.test.Minimal", "/value", 1, enc(t, 2.0, nil, payload.FormatBSON), a2))
		rig.sh.batch.flush(ctx)

		if !a1.acked() || !a2.acked() {
			t.Error("dropped poisoned ops must still be acknowledged")
		}
		if got := len(fs.individualRows()); got != 0 {
			t.Errorf("%d rows persisted from a fully poisoned batch, want 0", got)
		}
		if got := promtest.ToFloat64(rig.e.met.internalErrors); got != 2 {
			t.Errorf("internal error counter = %v, want 2", got)
		}
	})
}

// TestFlushAbandonOnQuit: once drain has begun, a failing flush abandons its
// batch unacknowledged instead of parking forever (docs/DESIGN.md §5.3).
func TestFlushAbandonOnQuit(t *testing.T) {
	ctx := context.Background()
	rig, fs := newPipelineRig(t, Config{BatchMaxRows: 1000})
	fs.mu.Lock()
	fs.appendErrs = []error{errTransient}
	fs.mu.Unlock()

	ack := &ackCounter{}
	rig.handle(deviceMsg("com.astrate.test.Minimal", "/value", 1, enc(t, 1.0, nil, payload.FormatBSON), ack))
	rig.e.quitOnce.Do(func() { close(rig.e.quit) })

	flushDone := make(chan struct{})
	go func() {
		rig.sh.batch.flush(ctx)
		close(flushDone)
	}()
	select {
	case <-flushDone:
	case <-time.After(5 * time.Second):
		t.Fatal("flush kept parking after quit")
	}
	if ack.acked() {
		t.Error("abandoned message was acknowledged")
	}
	if rig.sh.batch.size() != 0 {
		t.Errorf("abandoned batch still holds %d ops", rig.sh.batch.size())
	}
	if len(fs.individualRows()) != 0 {
		t.Error("abandoned batch was persisted")
	}
}

// TestBatchTriggers covers the two flush triggers through the running
// pipeline: the row cap and the wait timer (docs/DESIGN.md §1.4).
func TestBatchTriggers(t *testing.T) {
	t.Run("RowCap", func(t *testing.T) {
		fs := newFakeStore()
		seedAlpha(t, fs)
		e := startTestEngine(t, fs, Config{Shards: 1, BatchMaxRows: 4, BatchMaxWait: time.Hour})
		ack := &ackCounter{}
		for i := range 8 {
			e.Submit(deviceMsg("com.astrate.test.Minimal", "/value", 1,
				enc(t, float64(i), nil, payload.FormatBSON), ack))
		}
		waitFor(t, 5*time.Second, "both row-cap batches", func() bool {
			return len(fs.individualRows()) == 8
		})
		if got := fs.batchCount(); got != 2 {
			t.Errorf("%d batches, want 2 (row cap of 4)", got)
		}
	})

	t.Run("WaitTimer", func(t *testing.T) {
		fs := newFakeStore()
		seedAlpha(t, fs)
		e := startTestEngine(t, fs, Config{Shards: 1, BatchMaxRows: 1000, BatchMaxWait: 30 * time.Millisecond})
		ack := &ackCounter{}
		for i := range 3 {
			e.Submit(deviceMsg("com.astrate.test.Minimal", "/value", 1,
				enc(t, float64(i), nil, payload.FormatBSON), ack))
		}
		waitFor(t, 5*time.Second, "timer-driven flush", func() bool {
			return len(fs.individualRows()) == 3
		})
		if got := fs.batchCount(); got != 1 {
			t.Errorf("%d batches, want 1 (single timer flush)", got)
		}
	})
}

// TestAfterCommitHook: the M6b observer sees committed ops, in order,
// excluding broken ones.
func TestAfterCommitHook(t *testing.T) {
	ctx := context.Background()
	rig, fs := newPipelineRig(t, Config{BatchMaxRows: 1000})

	var observed []PersistOp
	rig.e.afterCommit = func(ops []PersistOp) { observed = append(observed, ops...) }

	ack := &ackCounter{}
	rig.handle(deviceMsg("com.astrate.test.Minimal", "/value", 1, enc(t, 1.0, nil, payload.FormatBSON), ack))
	rig.handle(deviceMsg("com.astrate.test.PropertyArrays", "/config/thresholds", 2, []byte(`{"v":[1.0]}`), ack))
	// A hand-broken op: value/type mismatch that validation rules out.
	scalars := compiledFixture(t, 10, "com.astrate.test.AllScalarTypes.json")
	m, _ := scalars.Trie.Match("/double")
	broken := &ackCounter{}
	rig.sh.batch.add(PersistOp{
		Kind: OpIndividual, Realm: realmAlpha, RealmID: realmAlphaID, DeviceID: devAlpha,
		Interface: scalars, Mapping: m, Path: "/double", Value: "not a double",
		TS: time.Now(), ReceptionTS: time.Now(), ack: broken.fn(),
	})

	rig.sh.batch.flush(ctx)

	if len(observed) != 2 {
		t.Fatalf("afterCommit saw %d ops, want 2 (broken op excluded)", len(observed))
	}
	if observed[0].Kind != OpIndividual || observed[1].Kind != OpPropertySet {
		t.Errorf("afterCommit order: %s, %s", observed[0].Kind, observed[1].Kind)
	}
	if !broken.acked() {
		t.Error("broken op not consumed")
	}
	if len(fs.individualRows()) != 1 {
		t.Errorf("%d datastream rows, want 1 (broken op skipped)", len(fs.individualRows()))
	}
}
