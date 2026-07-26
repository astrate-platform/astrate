package testutil

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

// updateGolden rewrites golden files instead of comparing against them:
//
//	go test ./... -update
var updateGolden = flag.Bool("update", false, "rewrite golden files with the received content")

// Golden compares got against the golden file testdata/<name> (relative to the
// calling package). With -update the file is (re)written instead and the
// comparison always passes. Golden fixtures are wire-frozen bytes (envelopes,
// payload vectors), so the comparison is exact — byte for byte.
func Golden(t testing.TB, name string, got []byte) {
	t.Helper()

	path := filepath.Join("testdata", filepath.FromSlash(name))
	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("golden: creating %s: %v", filepath.Dir(path), err)
			return
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatalf("golden: writing %s: %v", path, err)
			return
		}
		t.Logf("golden: wrote %s (%d bytes)", path, len(got))
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("golden: reading %s: %v (create it with `go test -update`)", path, err)
		return
	}
	if bytes.Equal(got, want) {
		return
	}

	t.Errorf("golden: %s mismatch at byte %d\n--- want (%d bytes)\n%s\n--- got (%d bytes)\n%s",
		path, firstDiff(got, want), len(want), want, len(got), got)
}

// firstDiff returns the offset of the first differing byte between a and b.
func firstDiff(a, b []byte) int {
	n := min(len(a), len(b))
	for i := range n {
		if a[i] != b[i] {
			return i
		}
	}
	return n
}
