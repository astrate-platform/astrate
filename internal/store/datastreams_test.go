//go:build integration

package store

import (
	"context"
	"fmt"
	"strconv"
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

	t.Run("IndividualSnapshot", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		device := mustRegisterDevice(t, s, realm.ID)
		si := mustInstallInterface(t, s, realm.ID, allTypesDef)

		base := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
		var batch DatastreamBatch
		// Two distinct paths, three samples each: the snapshot must return
		// only the newest sample per path.
		for _, path := range []string{"/d", "/s"} {
			for i := range 3 {
				r := IndividualRow{
					RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID,
					EndpointID: si.Endpoints[path], Path: path,
					TS: base.Add(time.Duration(i) * time.Minute), ReceptionTS: base,
				}
				switch path {
				case "/d":
					v := float64(i)
					r.ValueDouble = &v
				case "/s":
					v := "v" + strconv.Itoa(i)
					r.ValueString = &v
				}
				batch.Individual = append(batch.Individual, r)
			}
		}
		if err := s.AppendDatastreams(ctx, batch); err != nil {
			t.Fatalf("AppendDatastreams: %v", err)
		}

		rows, err := s.IndividualSnapshot(ctx, realm.ID, device, si.ID)
		if err != nil {
			t.Fatalf("IndividualSnapshot: %v", err)
		}
		if len(rows) != 2 {
			t.Fatalf("snapshot rows = %d, want 2 (one per path)", len(rows))
		}
		latest := map[string]IndividualRow{}
		for _, r := range rows {
			latest[r.Path] = r
		}
		if r, ok := latest["/d"]; !ok || r.ValueDouble == nil || *r.ValueDouble != 2 {
			t.Errorf("/d snapshot = %+v, want value 2", latest["/d"])
		}
		if r, ok := latest["/s"]; !ok || r.ValueString == nil || *r.ValueString != "v2" {
			t.Errorf("/s snapshot = %+v, want value v2", latest["/s"])
		}
		if !latest["/d"].TS.Equal(base.Add(2 * time.Minute)) {
			t.Errorf("/d snapshot ts = %v, want newest", latest["/d"].TS)
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

	t.Run("DownsampleLTTB", func(t *testing.T) {
		if !s.HasToolkitLTTB() {
			t.Skip("timescaledb_toolkit not installed; lttb path unavailable")
		}
		realm := mustCreateRealm(t, s)
		device := mustRegisterDevice(t, s, realm.ID)
		si := mustInstallInterface(t, s, realm.ID, allTypesDef)

		// A flat series with one sharp spike. This is the whole reason lttb
		// exists: time_bucket+avg would dilute the spike into its bucket's
		// mean and the shape would be lost, whereas lttb selects the real
		// sample. 60 points, spike of 1000 at index 30.
		base := time.Date(2026, 1, 18, 12, 0, 0, 0, time.UTC)
		const spikeIdx, spikeVal = 30, 1000.0
		var batch DatastreamBatch
		for i := range 60 {
			v := 1.0
			if i == spikeIdx {
				v = spikeVal
			}
			batch.Individual = append(batch.Individual, IndividualRow{
				RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID,
				EndpointID: si.Endpoints["/d"], Path: "/d",
				TS: base.Add(time.Duration(i) * time.Minute), ReceptionTS: base,
				ValueDouble: &v,
			})
		}
		if err := s.AppendDatastreams(ctx, batch); err != nil {
			t.Fatalf("AppendDatastreams: %v", err)
		}
		q := SeriesQuery{RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/d"}

		got, err := s.DownsampleLTTB(ctx, q, 8)
		if err != nil {
			t.Fatalf("DownsampleLTTB: %v", err)
		}
		if len(got) != 8 {
			t.Fatalf("resolution 8: got %d points, want 8", len(got))
		}
		// Ascending by time, first and last samples preserved.
		for i := 1; i < len(got); i++ {
			if !got[i].Bucket.After(got[i-1].Bucket) {
				t.Errorf("point %d not strictly after its predecessor: %v", i, got)
			}
		}
		if !got[0].Bucket.Equal(base) {
			t.Errorf("first point: got %v, want %v", got[0].Bucket, base)
		}
		if want := base.Add(59 * time.Minute); !got[len(got)-1].Bucket.Equal(want) {
			t.Errorf("last point: got %v, want %v", got[len(got)-1].Bucket, want)
		}
		// The spike survives, at its own real timestamp and exact value. An
		// averaging downsample cannot pass this: with 60 points reduced to 8
		// the spike shares a bucket with flat neighbours and its mean is far
		// below 1000.
		spikeTS := base.Add(spikeIdx * time.Minute)
		var found bool
		for _, p := range got {
			if p.Bucket.Equal(spikeTS) {
				found = true
				if p.Value != spikeVal {
					t.Errorf("spike value: got %g, want %g", p.Value, spikeVal)
				}
			}
		}
		if !found {
			t.Errorf("lttb dropped the spike at %v: %v", spikeTS, got)
		}
		// Every returned value is a real sample, never an average.
		for _, p := range got {
			if p.Value != 1 && p.Value != spikeVal {
				t.Errorf("point %v is not a real input sample", p)
			}
		}

		// Descending reverses the output without changing which samples lttb
		// picked: the aggregate sorts internally, so the order clause is
		// applied outside it.
		desc, err := s.DownsampleLTTB(ctx, SeriesQuery{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/d", Descending: true,
		}, 8)
		if err != nil {
			t.Fatalf("DownsampleLTTB descending: %v", err)
		}
		if len(desc) != len(got) {
			t.Fatalf("descending: got %d points, want %d", len(desc), len(got))
		}
		for i := range desc {
			if !desc[i].Bucket.Equal(got[len(got)-1-i].Bucket) {
				t.Errorf("descending point %d: got %v, want %v", i, desc[i].Bucket, got[len(got)-1-i].Bucket)
			}
		}

		// Limit trims the output but must not change the resolution lttb ran
		// at — the first three points are the same as the unlimited run's.
		lim, err := s.DownsampleLTTB(ctx, SeriesQuery{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/d", Limit: 3,
		}, 8)
		if err != nil {
			t.Fatalf("DownsampleLTTB limit: %v", err)
		}
		if len(lim) != 3 {
			t.Fatalf("limit 3: got %d points", len(lim))
		}
		for i := range lim {
			if !lim[i].Bucket.Equal(got[i].Bucket) || lim[i].Value != got[i].Value {
				t.Errorf("limited point %d: got %v, want %v", i, lim[i], got[i])
			}
		}

		// Integers coalesce into the value expression like the bucket path.
		var intBatch DatastreamBatch
		for i := range 10 {
			v := int32(i)
			intBatch.Individual = append(intBatch.Individual, IndividualRow{
				RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID,
				EndpointID: si.Endpoints["/i"], Path: "/i",
				TS: base.Add(time.Duration(i) * time.Minute), ReceptionTS: base,
				ValueInteger: &v,
			})
		}
		if err := s.AppendDatastreams(ctx, intBatch); err != nil {
			t.Fatalf("AppendDatastreams integers: %v", err)
		}
		ints, err := s.DownsampleLTTB(ctx, SeriesQuery{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/i",
		}, 4)
		if err != nil {
			t.Fatalf("DownsampleLTTB integers: %v", err)
		}
		if len(ints) != 4 {
			t.Errorf("integer lttb: got %d points, want 4", len(ints))
		}

		// Fewer samples than the resolution returns them all, not an error.
		few, err := s.DownsampleLTTB(ctx, SeriesQuery{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/i",
		}, 50)
		if err != nil {
			t.Fatalf("DownsampleLTTB oversized resolution: %v", err)
		}
		if len(few) != 10 {
			t.Errorf("oversized resolution: got %d points, want all 10", len(few))
		}

		// A resolution below 3 is refused by us, before Postgres raises
		// "resolution must be greater than 2".
		for _, n := range []int{-1, 0, 1, 2} {
			if _, err := s.DownsampleLTTB(ctx, q, n); err == nil {
				t.Errorf("resolution %d accepted", n)
			}
		}
	})

	t.Run("SeriesSpan", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		device := mustRegisterDevice(t, s, realm.ID)
		si := mustInstallInterface(t, s, realm.ID, allTypesDef)

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
		if err := s.AppendDatastreams(ctx, batch); err != nil {
			t.Fatalf("AppendDatastreams: %v", err)
		}

		// Happy path: span covers all 10 minutes.
		first, last, ok, err := s.SeriesSpan(ctx, SeriesQuery{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/d",
		})
		if err != nil {
			t.Fatalf("SeriesSpan: %v", err)
		}
		if !ok {
			t.Fatal("SeriesSpan: ok = false, want true")
		}
		if !first.Equal(base) {
			t.Errorf("first = %v, want %v", first, base)
		}
		if !last.Equal(base.Add(9 * time.Minute)) {
			t.Errorf("last = %v, want %v", last, base.Add(9*time.Minute))
		}

		// Invariant: Limit and Descending must not change the span.
		first2, last2, ok2, err2 := s.SeriesSpan(ctx, SeriesQuery{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/d",
			Limit:      1,
			Descending: true,
		})
		if err2 != nil {
			t.Fatalf("SeriesSpan with Limit/Descending: %v", err2)
		}
		if !ok2 {
			t.Fatal("SeriesSpan with Limit/Descending: ok = false")
		}
		if !first.Equal(first2) || !last.Equal(last2) {
			t.Errorf("span changed with Limit/Descending: got (%v, %v), want (%v, %v)",
				first2, last2, first, last)
		}

		// Filters are honoured: Since and To narrow the span.
		since := base.Add(3 * time.Minute)
		to := base.Add(6 * time.Minute)
		first3, last3, ok3, err3 := s.SeriesSpan(ctx, SeriesQuery{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/d",
			Since: &since, To: &to,
		})
		if err3 != nil {
			t.Fatalf("SeriesSpan with Since/To: %v", err3)
		}
		if !ok3 {
			t.Fatal("SeriesSpan with Since/To: ok = false")
		}
		if !first3.Equal(since) {
			t.Errorf("filtered first = %v, want %v", first3, since)
		}
		if !last3.Equal(to) {
			t.Errorf("filtered last = %v, want %v", last3, to)
		}

		// Empty series: path that received nothing.
		_, _, okEmpty, errEmpty := s.SeriesSpan(ctx, SeriesQuery{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/i",
		})
		if errEmpty != nil {
			t.Fatalf("SeriesSpan empty: %v", errEmpty)
		}
		if okEmpty {
			t.Error("SeriesSpan empty: ok = true, want false")
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

// testRetention exercises the config-driven global drop-chunks policy
// (docs/DESIGN.md §2.5): applying it registers a policy_retention job on both
// hypertables, and clearing it (d <= 0) removes them.
func testRetention(t *testing.T, s *Store) {
	ctx := context.Background()

	policies := func() int {
		var n int
		if err := s.pool.QueryRow(ctx, `
			SELECT count(*) FROM timescaledb_information.jobs
			WHERE proc_name = 'policy_retention'
			  AND hypertable_name IN ('individual_datastreams', 'object_datastreams')`).Scan(&n); err != nil {
			t.Fatalf("counting retention jobs: %v", err)
		}
		return n
	}

	if err := s.ApplyGlobalRetention(ctx, 90*24*time.Hour); err != nil {
		t.Fatalf("ApplyGlobalRetention(set): %v", err)
	}
	if n := policies(); n != 2 {
		t.Fatalf("retention jobs after set = %d, want 2 (one per hypertable)", n)
	}

	// Idempotent re-apply with a new duration must not error or duplicate.
	if err := s.ApplyGlobalRetention(ctx, 30*24*time.Hour); err != nil {
		t.Fatalf("ApplyGlobalRetention(re-set): %v", err)
	}
	if n := policies(); n != 2 {
		t.Fatalf("retention jobs after re-set = %d, want 2", n)
	}

	if err := s.ApplyGlobalRetention(ctx, 0); err != nil {
		t.Fatalf("ApplyGlobalRetention(clear): %v", err)
	}
	if n := policies(); n != 0 {
		t.Fatalf("retention jobs after clear = %d, want 0", n)
	}
}
