// Command bsoncapture captures golden BSON payload vectors from the official
// Astarte Go device SDK encoder for pkg/payload/testdata/bson/.
//
// This is a standalone module (deliberately outside the main Astrate module)
// so the SDK's dependency tree — mongo-driver v1, astarte-go — never leaks
// into Astrate's go.mod. Provenance, the pinned versions, and the regenerate
// command are documented in pkg/payload/testdata/bson/README.md.
//
// Usage:
//
//	go run . ../../pkg/payload/testdata/bson
package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/astarte-platform/astarte-go/interfaces"
	"go.mongodb.org/mongo-driver/bson"
)

// sdkEncode reproduces, verbatim, the payload construction performed by the
// official astarte-device-sdk-go v0.90.2 in device/protocol_mqtt_v1.go,
// enqueueMqttV1Message (lines 417-428):
//
//	payload := map[string]interface{}{"v": interfaces.NormalizePayload(values, false)}
//	if !timestamp.IsZero() {
//		payload["t"] = timestamp.UTC()
//	}
//	doc, err := bson.Marshal(payload)
//
// That function is unexported and reachable only through a connected Device,
// so this harness calls the exact same expression against the exact same
// dependencies the SDK pins (astarte-go v0.90.4, mongo-driver v1.8.0).
func sdkEncode(values interface{}, timestamp time.Time) ([]byte, error) {
	payload := map[string]interface{}{"v": interfaces.NormalizePayload(values, false)}
	if !timestamp.IsZero() {
		payload["t"] = timestamp.UTC()
	}
	return bson.Marshal(payload)
}

// captureTime is the fixed explicit timestamp used by every *_t.hex vector:
// 2026-06-10T12:34:56.789Z.
var captureTime = time.Date(2026, 6, 10, 12, 34, 56, 789000000, time.UTC)

// vectors lists one representative value per Astarte mapping type, plus an
// object-aggregation document and an empty array. Go types match canonical
// SDK usage for each Astarte type.
var vectors = []struct {
	name  string
	value interface{}
}{
	{"double", float64(22.5)},
	{"integer", int32(42)},
	{"boolean", true},
	{"longinteger", int64(9007199254740993)}, // 2^53 + 1: not float64-representable
	{"string", "héllo, Astarte ✓"},
	{"binaryblob", []byte{0x00, 0x01, 0x02, 0xFE, 0xFF}},
	{"datetime", time.Date(2025, 12, 31, 23, 59, 59, 999000000, time.UTC)},
	{"doublearray", []float64{1.5, -2.25, 0}},
	{"integerarray", []int32{-1, 0, 2147483647}},
	{"booleanarray", []bool{true, false, true}},
	{"longintegerarray", []int64{1, -9007199254740993, 9223372036854775807}},
	{"stringarray", []string{"a", "β", ""}},
	{"binaryblobarray", [][]byte{{0xDE, 0xAD}, {0xBE, 0xEF}, {}}},
	{"datetimearray", []time.Time{
		time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 6, 10, 12, 0, 0, 500000000, time.UTC),
	}},
	{"doublearray_empty", []float64{}},
	{"object", map[string]interface{}{
		"lat":     45.4642,
		"lon":     9.19,
		"samples": int32(3),
		"ok":      true,
	}},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <output-dir>\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}
	outDir := os.Args[1]
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "bsoncapture: %v\n", err)
		os.Exit(1)
	}
	for _, vec := range vectors {
		for _, withT := range []bool{false, true} {
			ts := time.Time{}
			suffix := ""
			if withT {
				ts = captureTime
				suffix = "_t"
			}
			doc, err := sdkEncode(vec.value, ts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "bsoncapture: encoding %s%s: %v\n", vec.name, suffix, err)
				os.Exit(1)
			}
			path := filepath.Join(outDir, vec.name+suffix+".hex")
			if err := os.WriteFile(path, []byte(hex.EncodeToString(doc)+"\n"), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "bsoncapture: writing %s: %v\n", path, err)
				os.Exit(1)
			}
			fmt.Printf("wrote %s (%d bytes)\n", path, len(doc))
		}
	}
}
