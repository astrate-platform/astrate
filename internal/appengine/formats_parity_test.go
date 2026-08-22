//go:build integration

package appengine

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

const aeObjFmt = "com.ex.M7b.Formats"

var objFmtDef = `{"interface_name":"` + aeObjFmt + `","version_major":1,"version_minor":0,` +
	`"type":"datastream","ownership":"device","aggregation":"object",` +
	`"mappings":[{"endpoint":"/data/temp","type":"double"},` +
	`{"endpoint":"/data/hum","type":"double"},` +
	`{"endpoint":"/data/big","type":"longinteger"}]}`

// Bigint rendering bands for the /data/big column: one value inside the safe
// range and two outside it. All three are float64-exact, so decoded cells
// compare exactly while their JSON type (number vs string) still distinguishes
// the flag behavior.
var objFmtBigs = []int64{4503599627370495, 9007199254740992, int64(1) << 60} // 2^52-1, 2^53, 2^60

func TestQueryParamValidation(t *testing.T) {
	r := newRig(t)
	path := r.dpath("/interfaces/"+aeSensors+"/value") + "?"
	tests := []struct {
		name     string
		query    string
		wantCode int
		fragment string
	}{
		{
			"conflicting anchors",
			"since=" + iso(r.t2) + "&since_after=" + iso(r.t2), http.StatusUnprocessableEntity,
			`{"errors":{"since_after":["conflicts already set parameter"]}}`,
		},
		{"invalid since", "since=nope", http.StatusUnprocessableEntity, `{"errors":{"since":["is invalid"]}}`},
		{
			"negative limit",
			"limit=-1", http.StatusUnprocessableEntity,
			`{"errors":{"limit":["must be greater than or equal to 0"]}}`,
		},
		{
			"downsample_to=2",
			"downsample_to=2", http.StatusUnprocessableEntity,
			`{"errors":{"downsample_to":["must be greater than 2"]}}`,
		},
		{"bogus format", "format=bogus", http.StatusUnprocessableEntity, `{"errors":{"format":["is invalid"]}}`},
		{"since_after alone accepted", "since_after=" + iso(r.t2), http.StatusOK, ""},
		{"positive limit accepted", "limit=5", http.StatusOK, ""},
		{"table format accepted", "format=table", http.StatusOK, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := r.req(t, http.MethodGet, path+tt.query, "", r.token)
			if rec.Code != tt.wantCode {
				t.Fatalf("got %d want %d (body %s)", rec.Code, tt.wantCode, rec.Body)
			}
			if tt.fragment != "" && !strings.Contains(rec.Body.String(), tt.fragment) {
				t.Errorf("body %s missing %s", rec.Body, tt.fragment)
			}
		})
	}
}

func TestIndividualFormats(t *testing.T) {
	r := newRig(t)
	path := r.dpath("/interfaces/" + aeSensors + "/value")
	wantV := []float64{3, 2, 1}
	wantTs := []time.Time{r.t3, r.t2, r.t2.Add(-time.Minute)}

	t.Run("StructuredUnchanged", func(t *testing.T) {
		var got []Sample
		decodeData(t, r.req(t, http.MethodGet, path, "", r.token), &got)
		if len(got) != 3 {
			t.Fatalf("samples = %d, want 3", len(got))
		}
		for i := range got {
			v, _ := got[i].Value.(float64)
			if v != wantV[i] || !got[i].Timestamp.Equal(wantTs[i]) {
				t.Errorf("sample[%d] = %v@%v, want %v@%v", i, got[i].Value, got[i].Timestamp, wantV[i], wantTs[i])
			}
		}
	})

	t.Run("Table", func(t *testing.T) {
		rec := r.req(t, http.MethodGet, path+"?format=table", "", r.token)
		if rec.Code != http.StatusOK {
			t.Fatalf("got %d (body %s)", rec.Code, rec.Body)
		}
		var env struct {
			Data     [][]any `json:"data"`
			Metadata struct {
				Columns     map[string]int `json:"columns"`
				TableHeader []string       `json:"table_header"`
			} `json:"metadata"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("decode: %v (%s)", err, rec.Body)
		}
		if len(env.Data) != 3 {
			t.Fatalf("rows = %d, want 3 (%s)", len(env.Data), rec.Body)
		}
		for i, row := range env.Data {
			if len(row) != 2 {
				t.Fatalf("row %d = %v, want [timestamp,value]", i, row)
			}
			// Timestamp FIRST in individual tables.
			s, ok := row[0].(string)
			if !ok {
				t.Fatalf("row %d first element %v (%T), want timestamp", i, row[0], row[0])
			}
			got, err := time.Parse(time.RFC3339, s)
			if err != nil || !got.Equal(wantTs[i]) {
				t.Errorf("row %d timestamp %q, want %v", i, s, wantTs[i])
			}
			if v, ok := row[1].(float64); !ok || v != wantV[i] {
				t.Errorf("row %d value %v (%T), want %v", i, row[1], row[1], wantV[i])
			}
		}
		if env.Metadata.Columns["timestamp"] != 0 || env.Metadata.Columns["value"] != 1 {
			t.Errorf("columns = %v", env.Metadata.Columns)
		}
		if fmt.Sprint(env.Metadata.TableHeader) != "[timestamp value]" {
			t.Errorf("table_header = %v", env.Metadata.TableHeader)
		}
	})

	t.Run("NoMetadataWithoutFormat", func(t *testing.T) {
		rec := r.req(t, http.MethodGet, path, "", r.token)
		var body map[string]json.RawMessage
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if _, ok := body["metadata"]; ok {
			t.Error("structured response carries a metadata key")
		}
	})

	t.Run("DisjointTablesValueFirst", func(t *testing.T) {
		rec := r.req(t, http.MethodGet, path+"?format=disjoint_tables", "", r.token)
		if rec.Code != http.StatusOK {
			t.Fatalf("got %d (body %s)", rec.Code, rec.Body)
		}
		var env struct {
			Data map[string][][]any `json:"data"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("decode: %v (%s)", err, rec.Body)
		}
		col, ok := env.Data["value"]
		if !ok || len(col) != 3 {
			t.Fatalf("data = %v (%s)", env.Data, rec.Body)
		}
		for i, pair := range col {
			if len(pair) != 2 {
				t.Fatalf("pair %d = %v", i, pair)
			}
			// Value FIRST, then timestamp — upstream's pair order.
			v, vOK := pair[0].(float64)
			if !vOK || v != wantV[i] {
				t.Errorf("pair %d first element %v (%T), want number %v", i, pair[0], pair[0], wantV[i])
			}
			s, sOK := pair[1].(string)
			got, err := time.Parse(time.RFC3339, s)
			if !sOK || err != nil || !got.Equal(wantTs[i]) {
				t.Errorf("pair %d second element %v (%T), want timestamp %v", i, pair[1], pair[1], wantTs[i])
			}
		}
	})

	t.Run("LongIntegerTableStaysString", func(t *testing.T) {
		rec := r.req(t, http.MethodGet, r.dpath("/interfaces/"+aeSensors+"/big")+"?format=table", "", r.token)
		if rec.Code != http.StatusOK {
			t.Fatalf("got %d (body %s)", rec.Code, rec.Body)
		}
		var env struct {
			Data [][]json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("decode: %v (%s)", err, rec.Body)
		}
		if len(env.Data) != 1 || len(env.Data[0]) != 2 {
			t.Fatalf("rows = %v (%s)", env.Data, rec.Body)
		}
		if got := string(env.Data[0][1]); got != `"1152921504606846976"` {
			t.Errorf("big cell = %s, want quoted decimal string", got)
		}
	})
}

// objFmtSeed installs aeObjFmt on a fresh device and seeds three documents on
// path "/12", oldest first: temp/hum vary by row, big spans the bigint bands.
func objFmtSeed(t *testing.T, r *rig) deviceid.ID {
	t.Helper()
	ctx := context.Background()
	dev, _ := deviceid.Random()
	if err := r.st.RegisterDevice(ctx, r.realmID, dev, "h"); err != nil {
		t.Fatal(err)
	}
	si, err := r.st.InstallInterface(ctx, r.realmID, []byte(objFmtDef))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.st.UpdateIntrospection(ctx, r.realmID, dev, map[string]store.InterfaceVersion{
		aeObjFmt: {Major: 1, Minor: 0},
	}); err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 6, 3, 9, 0, 0, 0, time.UTC)
	var batch store.DatastreamBatch
	for i := range 3 {
		ts := base.Add(time.Duration(i) * time.Minute)
		batch.Objects = append(batch.Objects, store.ObjectRow{
			RealmID: r.realmID, DeviceID: dev, InterfaceID: si.ID, Path: "/12",
			TS: ts, ReceptionTS: ts,
			Value: fmt.Appendf(nil, `{"temp":%v,"hum":%v,"big":%d}`, 21.5+float64(i), 40+i, objFmtBigs[i]),
		})
	}
	if err := r.st.AppendDatastreams(ctx, batch); err != nil {
		t.Fatal(err)
	}
	return dev
}

func objFmtDeviceOnly(t *testing.T, r *rig) deviceid.ID {
	t.Helper()
	ctx := context.Background()
	dev, _ := deviceid.Random()
	if err := r.st.RegisterDevice(ctx, r.realmID, dev, "h"); err != nil {
		t.Fatal(err)
	}
	if _, err := r.st.UpdateIntrospection(ctx, r.realmID, dev, map[string]store.InterfaceVersion{
		aeObjFmt: {Major: 1, Minor: 0},
	}); err != nil {
		t.Fatal(err)
	}
	return dev
}

func TestObjectFormats(t *testing.T) {
	r := newRig(t)
	dev := objFmtSeed(t, r)
	path := "/devices/" + dev.String() + "/interfaces/" + aeObjFmt + "/12"
	// Rows come back descending; doc i was seeded at base.Add(i*minute).
	wantTemp := []float64{23.5, 22.5, 21.5}
	wantHum := []float64{42, 41, 40}
	wantBig := []float64{1152921504606846976, 9007199254740992, 4503599627370495}

	t.Run("StructuredUnchanged", func(t *testing.T) {
		var got []map[string]any
		decodeData(t, r.req(t, http.MethodGet, path, "", r.token), &got)
		if len(got) != 3 {
			t.Fatalf("docs = %d, want 3", len(got))
		}
		for i := range got {
			if got[i]["temp"] != wantTemp[i] || got[i]["hum"] != wantHum[i] {
				t.Errorf("doc[%d] = %v", i, got[i])
			}
			if ts, ok := got[i]["timestamp"].(string); !ok || ts == "" {
				t.Errorf("doc[%d] missing timestamp: %v", i, got[i])
			}
			// Default longinteger rendering is untouched: numbers.
			if big, ok := got[i]["big"].(float64); !ok || big != wantBig[i] {
				t.Errorf("doc[%d] big = %v (%T)", i, got[i]["big"], got[i]["big"])
			}
		}
	})

	t.Run("Table", func(t *testing.T) {
		rec := r.req(t, http.MethodGet, path+"?format=table", "", r.token)
		if rec.Code != http.StatusOK {
			t.Fatalf("got %d (body %s)", rec.Code, rec.Body)
		}
		var env struct {
			Data     [][]any `json:"data"`
			Metadata struct {
				Columns     map[string]int `json:"columns"`
				TableHeader []string       `json:"table_header"`
			} `json:"metadata"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("decode: %v (%s)", err, rec.Body)
		}
		if fmt.Sprint(env.Metadata.TableHeader) != "[big hum temp timestamp]" {
			t.Fatalf("table_header = %v (%s)", env.Metadata.TableHeader, rec.Body)
		}
		for name, idx := range map[string]int{"big": 0, "hum": 1, "temp": 2, "timestamp": 3} {
			if env.Metadata.Columns[name] != idx {
				t.Errorf("columns[%s] = %d, want %d", name, env.Metadata.Columns[name], idx)
			}
		}
		if len(env.Data) != 3 {
			t.Fatalf("rows = %d, want 3 (%s)", len(env.Data), rec.Body)
		}
		for i, row := range env.Data {
			if len(row) != 4 {
				t.Fatalf("row %d = %v", i, row)
			}
			big, ok := row[0].(float64)
			if !ok || big != wantBig[i] {
				t.Errorf("row %d big = %v (%T)", i, row[0], row[0])
			}
			if row[1] != wantHum[i] || row[2] != wantTemp[i] {
				t.Errorf("row %d = %v", i, row)
			}
			ts, ok := row[3].(string)
			if !ok || ts == "" {
				t.Errorf("row %d timestamp = %v (%T)", i, row[3], row[3])
			}
		}
	})

	t.Run("DisjointTablesValueFirst", func(t *testing.T) {
		rec := r.req(t, http.MethodGet, path+"?format=disjoint_tables", "", r.token)
		if rec.Code != http.StatusOK {
			t.Fatalf("got %d (body %s)", rec.Code, rec.Body)
		}
		var env struct {
			Data map[string][][]any `json:"data"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("decode: %v (%s)", err, rec.Body)
		}
		for _, name := range []string{"big", "hum", "temp"} {
			col, ok := env.Data[name]
			if !ok || len(col) != 3 {
				t.Fatalf("column %q = %v (%s)", name, env.Data[name], rec.Body)
			}
			for i, pair := range col {
				var want any
				switch name {
				case "big":
					want = wantBig[i]
				case "hum":
					want = wantHum[i]
				default:
					want = wantTemp[i]
				}
				if pair[0] != want {
					t.Errorf("%s[%d] value = %v, want %v", name, i, pair[0], want)
				}
				ts, ok := pair[1].(string)
				if !ok || ts == "" {
					t.Errorf("%s[%d] second element = %v (%T), want timestamp after value", name, i, pair[1], pair[1])
				}
			}
		}
		if _, ok := env.Data["timestamp"]; ok {
			t.Error("disjoint_tables carries its own timestamp column")
		}
	})

	t.Run("AllowBigIntegersFalseStringsEverywhere", func(t *testing.T) {
		wantStr := []string{"1152921504606846976", "9007199254740992", "4503599627370495"}
		for _, format := range []string{"structured", "table", "disjoint_tables"} {
			t.Run(format, func(t *testing.T) {
				query := "?allow_bigintegers=false&format=" + format
				rec := r.req(t, http.MethodGet, path+query, "", r.token)
				if rec.Code != http.StatusOK {
					t.Fatalf("got %d (body %s)", rec.Code, rec.Body)
				}
				switch format {
				case "structured":
					var docs []map[string]any
					if err := json.Unmarshal(envData(rec), &docs); err != nil {
						t.Fatal(err)
					}
					for i := range docs {
						s, ok := docs[i]["big"].(string)
						if !ok || s != wantStr[i] {
							t.Errorf("doc[%d] big = %v (%T)", i, docs[i]["big"], docs[i]["big"])
						}
					}
				case "table":
					var env struct {
						Data [][]any `json:"data"`
					}
					if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
						t.Fatal(err)
					}
					for i, row := range env.Data {
						s, ok := row[0].(string)
						if !ok || s != wantStr[i] {
							t.Errorf("row %d big = %v (%T)", i, row[0], row[0])
						}
					}
				case "disjoint_tables":
					var env struct {
						Data map[string][][]any `json:"data"`
					}
					if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
						t.Fatal(err)
					}
					for i, pair := range env.Data["big"] {
						s, ok := pair[0].(string)
						if !ok || s != wantStr[i] {
							t.Errorf("big[%d] = %v (%T)", i, pair[0], pair[0])
						}
					}
				}
			})
		}
	})

	t.Run("AllowSafeBigIntegersBand", func(t *testing.T) {
		rec := r.req(t, http.MethodGet, path+"?allow_safe_bigintegers=true", "", r.token)
		var docs []map[string]any
		decodeData(t, rec, &docs)
		// Above the bound → decimal strings...
		for i := range 2 {
			s, ok := docs[i]["big"].(string)
			if !ok || s != fmt.Sprint(int64(wantBig[i])) {
				t.Errorf("doc[%d] big = %v (%T), want unsafe band as string", i, docs[i]["big"], docs[i]["big"])
			}
		}
		// ...the single in-band value stays a number.
		if n, ok := docs[2]["big"].(float64); !ok || n != wantBig[2] {
			t.Errorf("doc[2] big = %v (%T), want in-band number", docs[2]["big"], docs[2]["big"])
		}
	})

	t.Run("EmptySeriesRootRendersEmptyObject", func(t *testing.T) {
		dev2 := objFmtDeviceOnly(t, r)
		root := "/devices/" + dev2.String() + "/interfaces/" + aeObjFmt
		for _, query := range []string{"", "?format=table", "?format=disjoint_tables"} {
			rec := r.req(t, http.MethodGet, root+query, "", r.token)
			if rec.Code != http.StatusOK {
				t.Fatalf("%s: got %d (body %s)", query, rec.Code, rec.Body)
			}
			if !strings.Contains(rec.Body.String(), `"data":{}`) {
				t.Errorf("%s: body %s missing \"data\":{}", query, rec.Body)
			}
		}
	})

	t.Run("EmptySeriesConcretePathNotFound", func(t *testing.T) {
		rec := r.req(t, http.MethodGet, path+"/none", "", r.token)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("got %d (body %s)", rec.Code, rec.Body)
		}
		if !strings.Contains(rec.Body.String(), "Path not found") {
			t.Errorf("body %s missing \"Path not found\"", rec.Body)
		}
	})
}

// envData extracts the raw "data" member of an envelope response.
func envData(rec *httptest.ResponseRecorder) json.RawMessage {
	var env struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		return nil
	}
	return env.Data
}
