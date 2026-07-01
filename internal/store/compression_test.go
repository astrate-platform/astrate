//go:build integration

package store

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

const compressObjDef = `{
	"interface_name": "com.astrate.test.CompressObj",
	"version_major": 1,
	"version_minor": 0,
	"type": "datastream",
	"ownership": "device",
	"aggregation": "object",
	"mappings": [
		{"endpoint": "/%{id}/temp", "type": "double"},
		{"endpoint": "/%{id}/hum", "type": "double"}
	]
}`

// testCompressionLive is the ROADMAP §10 T5 gate bullet: force
// compress_chunk on aged data in both hypertables, record the before/after
// size ratio (informational), and prove queries read identical rows back
// from compressed chunks.
func testCompressionLive(t *testing.T, s *Store) {
	ctx := context.Background()
	realm := mustCreateRealm(t, s)
	device := mustRegisterDevice(t, s, realm.ID)
	si := mustInstallInterface(t, s, realm.ID, allTypesDef)
	so := mustInstallInterface(t, s, realm.ID, compressObjDef)

	// Aged rows: 30 days old, in chunks well past the 7-day compression
	// horizon, so compress_chunk targets sealed chunks.
	base := time.Now().UTC().Add(-30 * 24 * time.Hour).Truncate(time.Hour)
	const n = 500
	var batch DatastreamBatch
	values := make([]float64, n)
	for i := range n {
		v := float64(i) / 3
		values[i] = v
		ts := base.Add(time.Duration(i) * time.Second)
		batch.Individual = append(batch.Individual, IndividualRow{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID,
			EndpointID: si.Endpoints["/d"], Path: "/d",
			TS: ts, ReceptionTS: ts,
			ValueDouble: &v,
		})
		batch.Objects = append(batch.Objects, ObjectRow{
			RealmID: realm.ID, DeviceID: device, InterfaceID: so.ID, Path: "/12",
			TS: ts, ReceptionTS: ts,
			Value: fmt.Appendf(nil, `{"temp": %g, "hum": %d}`, v, 40+i%10),
		})
	}
	if err := s.AppendDatastreams(ctx, batch); err != nil {
		t.Fatalf("AppendDatastreams: %v", err)
	}

	readBack := func(t *testing.T) ([]IndividualRow, []ObjectRow) {
		t.Helper()
		ind, err := s.Series(ctx, SeriesQuery{
			RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID, Path: "/d",
		})
		if err != nil {
			t.Fatalf("Series: %v", err)
		}
		obj, err := s.ObjectSeries(ctx, SeriesQuery{
			RealmID: realm.ID, DeviceID: device, InterfaceID: so.ID, Path: "/12",
		})
		if err != nil {
			t.Fatalf("ObjectSeries: %v", err)
		}
		return ind, obj
	}
	beforeInd, beforeObj := readBack(t)
	if len(beforeInd) != n || len(beforeObj) != n {
		t.Fatalf("pre-compression read: %d individual, %d object rows, want %d each",
			len(beforeInd), len(beforeObj), n)
	}

	// Force-compress every sealed chunk of both hypertables. On a reused
	// database earlier runs may have left chunks compressed or partial;
	// if_not_compressed tolerates both.
	for _, table := range []string{"individual_datastreams", "object_datastreams"} {
		var chunks int
		if err := s.pool.QueryRow(ctx, `
			SELECT count(*) FROM (
				SELECT compress_chunk(c, true)
				FROM show_chunks($1::regclass, older_than => now() - INTERVAL '14 days') AS c
			) AS compressed`, table).Scan(&chunks); err != nil {
			t.Fatalf("compress_chunk sweep on %s: %v", table, err)
		}
		if chunks == 0 {
			t.Fatalf("no aged chunks found on %s — the aged batch did not land where expected", table)
		}

		var beforeBytes, afterBytes int64
		if err := s.pool.QueryRow(ctx, `
			SELECT COALESCE(sum(before_compression_total_bytes), 0),
			       COALESCE(sum(after_compression_total_bytes), 0)
			FROM chunk_compression_stats($1::regclass)`, table).Scan(&beforeBytes, &afterBytes); err != nil {
			t.Fatalf("chunk_compression_stats on %s: %v", table, err)
		}
		// Informational per the gate: the ratio depends on what else shares
		// the chunks, so it is recorded, not asserted.
		if afterBytes > 0 {
			t.Logf("%s: %d chunk(s) compressed, %d → %d bytes (%.1fx)",
				table, chunks, beforeBytes, afterBytes, float64(beforeBytes)/float64(afterBytes))
		}
	}

	// Queries must read the exact same rows back from compressed chunks.
	afterInd, afterObj := readBack(t)
	if len(afterInd) != n || len(afterObj) != n {
		t.Fatalf("post-compression read: %d individual, %d object rows, want %d each",
			len(afterInd), len(afterObj), n)
	}
	for i := range n {
		b, a := beforeInd[i], afterInd[i]
		if !a.TS.Equal(b.TS) || a.ValueDouble == nil || *a.ValueDouble != *b.ValueDouble {
			t.Fatalf("individual row %d changed after compression: before ts=%v v=%v, after ts=%v v=%v",
				i, b.TS, *b.ValueDouble, a.TS, a.ValueDouble)
		}
		ob, oa := beforeObj[i], afterObj[i]
		if !oa.TS.Equal(ob.TS) || !strings.EqualFold(string(oa.Value), string(ob.Value)) {
			t.Fatalf("object row %d changed after compression: before ts=%v %s, after ts=%v %s",
				i, ob.TS, ob.Value, oa.TS, oa.Value)
		}
	}
}
