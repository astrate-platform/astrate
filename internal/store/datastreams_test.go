//go:build integration

package store

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

const allTypesDef = `{
	"interface_name": "com.astrate.test.AllTypes",
	"version_major": 1,
	"version_minor": 0,
	"type": "datastream",
	"ownership": "device",
	"mappings": [
		{"endpoint": "/d", "type": "double"},
		{"endpoint": "/i", "type": "integer"},
		{"endpoint": "/l", "type": "longinteger"},
		{"endpoint": "/b", "type": "boolean"},
		{"endpoint": "/s", "type": "string"},
		{"endpoint": "/bb", "type": "binaryblob"},
		{"endpoint": "/dt", "type": "datetime"},
		{"endpoint": "/da", "type": "doublearray"}
	]
}`

func testDatastreams(t *testing.T, s *Store) {
	ctx := context.Background()

	t.Run("ExactlyOneValueColumn", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		device := mustRegisterDevice(t, s, realm.ID)
		si := mustInstallInterface(t, s, realm.ID, allTypesDef)

		now := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
		base := IndividualRow{RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, TS: now, ReceptionTS: now}

		fl := 1.5
		i32 := int32(42)
		i64 := int64(1 << 40)
		bo := true
		st := "hello"
		dt := time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)
		rows := map[string]IndividualRow{}
		for path, set := range map[string]func(*IndividualRow){
			"/d":  func(r *IndividualRow) { r.ValueDouble = &fl },
			"/i":  func(r *IndividualRow) { r.ValueInteger = &i32 },
			"/l":  func(r *IndividualRow) { r.ValueLonginteger = &i64 },
			"/b":  func(r *IndividualRow) { r.ValueBoolean = &bo },
			"/s":  func(r *IndividualRow) { r.ValueString = &st },
			"/bb": func(r *IndividualRow) { r.ValueBinaryblob = []byte{0xde, 0xad} },
			"/dt": func(r *IndividualRow) { r.ValueDatetime = &dt },
			"/da": func(r *IndividualRow) { r.ValueArray = []byte(`[1.5,2.5]`) },
		} {
			r := base
			r.EndpointID = si.Endpoints[path]
			r.Path = path
			set(&r)
			rows[path] = r
		}

		var batch DatastreamBatch
		for _, r := range rows {
			batch.Individual = append(batch.Individual, r)
		}
		if err := s.AppendDatastreams(ctx, batch); err != nil {
			t.Fatalf("AppendDatastreams: %v", err)
		}

		// Each row read back must have exactly its declared column set.
		for path := range rows {
			got, err := s.Series(ctx, SeriesQuery{
				RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: path,
			})
			if err != nil {
				t.Fatalf("Series %s: %v", path, err)
			}
			if len(got) != 1 {
				t.Fatalf("Series %s: %d rows", path, len(got))
			}
			r := got[0]
			nonNil := map[string]bool{
				"/d":  r.ValueDouble != nil,
				"/i":  r.ValueInteger != nil,
				"/l":  r.ValueLonginteger != nil,
				"/b":  r.ValueBoolean != nil,
				"/s":  r.ValueString != nil,
				"/bb": r.ValueBinaryblob != nil,
				"/dt": r.ValueDatetime != nil,
				"/da": r.ValueArray != nil,
			}
			for col, set := range nonNil {
				if set != (col == path) {
					t.Errorf("row %s: column of %s set=%v", path, col, set)
				}
			}
		}

		// Zero or two value columns must be rejected before reaching COPY.
		bad := base
		bad.EndpointID = si.Endpoints["/d"]
		bad.Path = "/d"
		if err := s.AppendDatastreams(ctx, DatastreamBatch{Individual: []IndividualRow{bad}}); err == nil {
			t.Error("row with zero value columns accepted")
		}
		bad.ValueDouble = &fl
		bad.ValueInteger = &i32
		if err := s.AppendDatastreams(ctx, DatastreamBatch{Individual: []IndividualRow{bad}}); err == nil {
			t.Error("row with two value columns accepted")
		}
	})

	t.Run("SeriesWindows", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		device := mustRegisterDevice(t, s, realm.ID)
		si := mustInstallInterface(t, s, realm.ID, allTypesDef)

		base := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
		var batch DatastreamBatch
		for i := range 10 {
			v := float64(i)
			batch.Individual = append(batch.Individual, IndividualRow{
				RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID,
				EndpointID: si.Endpoints["/d"], Path: "/d",
				TS: base.Add(time.Duration(i) * time.Second), ReceptionTS: base,
				ValueDouble: &v,
			})
		}
		if err := s.AppendDatastreams(ctx, batch); err != nil {
			t.Fatalf("AppendDatastreams: %v", err)
		}

		at := func(i int) *time.Time { ts := base.Add(time.Duration(i) * time.Second); return &ts }
		q := SeriesQuery{RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/d"}

		var sb strings.Builder
		render := func(name string, q SeriesQuery) {
			rows, err := s.Series(ctx, q)
			if err != nil {
				t.Fatalf("Series %s: %v", name, err)
			}
			fmt.Fprintf(&sb, "[%s]\n", name)
			for _, r := range rows {
				fmt.Fprintf(&sb, "%s v=%g\n", r.TS.UTC().Format(time.RFC3339), *r.ValueDouble)
			}
		}

		full := q
		render("full", full)

		since := q
		since.Since = at(2) // inclusive: starts at v=2
		render("since=+2s", since)

		sinceAfter := q
		sinceAfter.SinceAfter = at(2) // exclusive: starts at v=3
		render("since_after=+2s", sinceAfter)

		to := q
		to.To = at(8) // inclusive: ends at v=8
		render("to=+8s", to)

		window := q
		window.Since = at(3)
		window.To = at(6)
		render("since=+3s,to=+6s", window)

		latest := q
		latest.Descending = true
		latest.Limit = 3
		render("desc,limit=3", latest)

		testutil.Golden(t, "series_windows.golden", []byte(sb.String()))
	})

	t.Run("DuplicateSeriesTimestamp", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		device := mustRegisterDevice(t, s, realm.ID)
		si := mustInstallInterface(t, s, realm.ID, allTypesDef)

		ts := time.Date(2026, 1, 20, 8, 0, 0, 0, time.UTC)
		v := 9.9
		row := IndividualRow{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID,
			EndpointID: si.Endpoints["/d"], Path: "/d",
			TS: ts, ReceptionTS: ts, ValueDouble: &v,
		}
		// At-least-once redelivery: the same (series, ts) twice must not error
		// (docs/DESIGN.md §5.3).
		if err := s.AppendDatastreams(ctx, DatastreamBatch{Individual: []IndividualRow{row}}); err != nil {
			t.Fatalf("first insert: %v", err)
		}
		if err := s.AppendDatastreams(ctx, DatastreamBatch{Individual: []IndividualRow{row}}); err != nil {
			t.Fatalf("duplicate insert: %v", err)
		}
		got, err := s.Series(ctx, SeriesQuery{RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/d"})
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 2 {
			t.Errorf("duplicate rows stored: %d, want 2", len(got))
		}
	})

	t.Run("ObjectSeries", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		device := mustRegisterDevice(t, s, realm.ID)
		si := mustInstallInterface(t, s, realm.ID, `{
			"interface_name": "com.astrate.test.ObjectAgg",
			"version_major": 1,
			"version_minor": 0,
			"type": "datastream",
			"ownership": "device",
			"aggregation": "object",
			"mappings": [
				{"endpoint": "/%{id}/temp", "type": "double"},
				{"endpoint": "/%{id}/hum", "type": "double"}
			]
		}`)

		base := time.Date(2026, 1, 16, 9, 0, 0, 0, time.UTC)
		var batch DatastreamBatch
		for i := range 3 {
			batch.Objects = append(batch.Objects, ObjectRow{
				RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/12",
				TS: base.Add(time.Duration(i) * time.Minute), ReceptionTS: base,
				Value: fmt.Appendf(nil, `{"temp": %d.5, "hum": %d}`, 20+i, 40+i),
			})
		}
		if err := s.AppendDatastreams(ctx, batch); err != nil {
			t.Fatalf("AppendDatastreams: %v", err)
		}

		rows, err := s.ObjectSeries(ctx, SeriesQuery{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/12",
			SinceAfter: &base,
		})
		if err != nil {
			t.Fatalf("ObjectSeries: %v", err)
		}
		if len(rows) != 2 {
			t.Fatalf("ObjectSeries since_after: %d rows, want 2", len(rows))
		}
		if !strings.Contains(string(rows[0].Value), `"temp": 21.5`) {
			t.Errorf("first windowed object row: %s", rows[0].Value)
		}
	})

	t.Run("Downsample", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		device := mustRegisterDevice(t, s, realm.ID)
		si := mustInstallInterface(t, s, realm.ID, allTypesDef)

		// One double per minute for 10 minutes; 5-minute buckets aligned to
		// the epoch → avg(0..4)=2 and avg(5..9)=7.
		base := time.Date(2026, 1, 18, 12, 0, 0, 0, time.UTC)
		var batch DatastreamBatch
		for i := range 10 {
			v := float64(i)
			batch.Individual = append(batch.Individual, IndividualRow{
				RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID,
				EndpointID: si.Endpoints["/d"], Path: "/d",
				TS: base.Add(time.Duration(i) * time.Minute), ReceptionTS: base,
				ValueDouble: &v,
			})
		}
		// Integer series on its own path: Downsample must coalesce
		// value_integer into the average as well.
		for i := range 4 {
			v := int32(10 * i)
			batch.Individual = append(batch.Individual, IndividualRow{
				RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID,
				EndpointID: si.Endpoints["/i"], Path: "/i",
				TS: base.Add(time.Duration(i) * time.Minute), ReceptionTS: base,
				ValueInteger: &v,
			})
		}
		if err := s.AppendDatastreams(ctx, batch); err != nil {
			t.Fatalf("AppendDatastreams: %v", err)
		}

		got, err := s.Downsample(ctx, SeriesQuery{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/d",
		}, 5*time.Minute)
		if err != nil {
			t.Fatalf("Downsample: %v", err)
		}
		want := []DownsamplePoint{
			{Bucket: base, Value: 2},
			{Bucket: base.Add(5 * time.Minute), Value: 7},
		}
		if len(got) != len(want) {
			t.Fatalf("Downsample buckets: got %v, want %v", got, want)
		}
		for i := range want {
			if !got[i].Bucket.Equal(want[i].Bucket) || got[i].Value != want[i].Value {
				t.Errorf("bucket %d: got (%v, %g), want (%v, %g)",
					i, got[i].Bucket, got[i].Value, want[i].Bucket, want[i].Value)
			}
		}

		intBuckets, err := s.Downsample(ctx, SeriesQuery{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/i",
		}, 5*time.Minute)
		if err != nil {
			t.Fatalf("Downsample integers: %v", err)
		}
		if len(intBuckets) != 1 || intBuckets[0].Value != 15 { // avg(0,10,20,30)
			t.Errorf("integer downsample: %v", intBuckets)
		}

		if _, err := s.Downsample(ctx, SeriesQuery{}, 0); err == nil {
			t.Error("zero bucket accepted")
		}
	})

	t.Run("BatchSmoke10k", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		si := mustInstallInterface(t, s, realm.ID, allTypesDef)

		devices := make([]deviceid.ID, 10)
		for i := range devices {
			devices[i] = mustRegisterDevice(t, s, realm.ID)
		}

		const total = 10_000
		base := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
		batch := DatastreamBatch{Individual: make([]IndividualRow, total)}
		values := make([]float64, total)
		for i := range total {
			values[i] = float64(i) / 10
			batch.Individual[i] = IndividualRow{
				RealmID: realm.ID, DeviceID: devices[i%len(devices)], InterfaceID: si.ID,
				EndpointID: si.Endpoints["/d"], Path: "/d",
				TS: base.Add(time.Duration(i) * time.Millisecond), ReceptionTS: base,
				ValueDouble: &values[i],
			}
		}

		start := time.Now()
		if err := s.AppendDatastreams(ctx, batch); err != nil {
			t.Fatalf("AppendDatastreams 10k: %v", err)
		}
		elapsed := time.Since(start)
		// Non-binding smoke budget (docs/ROADMAP.md §3.2): log, don't fail.
		t.Logf("10k-row batch landed in %s", elapsed)

		var n int
		if err := s.pool.QueryRow(ctx,
			`SELECT count(*) FROM individual_datastreams WHERE realm_id = $1 AND interface_id = $2`,
			realm.ID, si.ID).Scan(&n); err != nil {
			t.Fatal(err)
		}
		if n != total {
			t.Errorf("row count after batch: %d, want %d", n, total)
		}
	})
}

// testTTLJob inserts aged rows for use_ttl and no_ttl endpoints, runs the
// TTL action, and verifies only aged+TTL'd rows are gone (M2 gate).
func testTTLJob(t *testing.T, s *Store) {
	ctx := context.Background()
	realm := mustCreateRealm(t, s)
	device := mustRegisterDevice(t, s, realm.ID)

	individual := mustInstallInterface(t, s, realm.ID, `{
		"interface_name": "com.astrate.test.TTLValues",
		"version_major": 1,
		"version_minor": 0,
		"type": "datastream",
		"ownership": "device",
		"mappings": [
			{"endpoint": "/ttl", "type": "double", "database_retention_policy": "use_ttl", "database_retention_ttl": 3600},
			{"endpoint": "/keep", "type": "double"}
		]
	}`)
	object := mustInstallInterface(t, s, realm.ID, `{
		"interface_name": "com.astrate.test.TTLObject",
		"version_major": 1,
		"version_minor": 0,
		"type": "datastream",
		"ownership": "device",
		"aggregation": "object",
		"mappings": [
			{"endpoint": "/%{id}/a", "type": "double", "database_retention_policy": "use_ttl", "database_retention_ttl": 3600},
			{"endpoint": "/%{id}/b", "type": "double", "database_retention_policy": "use_ttl", "database_retention_ttl": 3600}
		]
	}`)

	now := time.Now().UTC()
	aged := now.Add(-2 * time.Hour) // beyond the 3600 s TTL
	v := 1.0
	mk := func(path string, endpointID int64, ts time.Time) IndividualRow {
		return IndividualRow{
			RealmID: realm.ID, DeviceID: device, InterfaceID: individual.ID,
			EndpointID: endpointID, Path: path, TS: ts, ReceptionTS: now, ValueDouble: &v,
		}
	}
	batch := DatastreamBatch{
		Individual: []IndividualRow{
			mk("/ttl", individual.Endpoints["/ttl"], aged),
			mk("/ttl", individual.Endpoints["/ttl"], now),
			mk("/keep", individual.Endpoints["/keep"], aged),
		},
		Objects: []ObjectRow{
			{RealmID: realm.ID, DeviceID: device, InterfaceID: object.ID, Path: "/1",
				TS: aged, ReceptionTS: now, Value: []byte(`{"a":1}`)},
			{RealmID: realm.ID, DeviceID: device, InterfaceID: object.ID, Path: "/1",
				TS: now, ReceptionTS: now, Value: []byte(`{"a":2}`)},
		},
	}
	if err := s.AppendDatastreams(ctx, batch); err != nil {
		t.Fatalf("AppendDatastreams: %v", err)
	}

	// Zero-argument Exec goes through the simple protocol, so the procedure
	// may COMMIT between chunks.
	if _, err := s.pool.Exec(ctx, "CALL astrate_apply_endpoint_ttl()"); err != nil {
		t.Fatalf("CALL astrate_apply_endpoint_ttl: %v", err)
	}

	count := func(table, path string, ifaceID int64) int {
		var n int
		if err := s.pool.QueryRow(ctx,
			`SELECT count(*) FROM `+table+` WHERE realm_id = $1 AND interface_id = $2 AND path = $3`,
			realm.ID, ifaceID, path).Scan(&n); err != nil {
			t.Fatalf("counting %s %s: %v", table, path, err)
		}
		return n
	}

	if n := count("individual_datastreams", "/ttl", individual.ID); n != 1 {
		t.Errorf("/ttl rows after TTL run: %d, want 1 (only the fresh row)", n)
	}
	if n := count("individual_datastreams", "/keep", individual.ID); n != 1 {
		t.Errorf("/keep rows after TTL run: %d, want 1 (no_ttl must keep aged rows)", n)
	}
	if n := count("object_datastreams", "/1", object.ID); n != 1 {
		t.Errorf("object rows after TTL run: %d, want 1 (only the fresh row)", n)
	}
}
