package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// recordingTB captures Errorf/Fatalf calls so the helpers under test can be
// observed failing without failing the real test.
type recordingTB struct {
	testing.TB
	errors []string
	fatals []string
}

func (r *recordingTB) Helper() {}

func (r *recordingTB) Logf(string, ...any) {}

func (r *recordingTB) Errorf(format string, args ...any) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}

func (r *recordingTB) Fatalf(format string, args ...any) {
	r.fatals = append(r.fatals, fmt.Sprintf(format, args...))
}

// setUpdateFlag overrides the -update flag for one test and restores it after.
func setUpdateFlag(t *testing.T, v bool) {
	t.Helper()
	old := *updateGolden
	*updateGolden = v
	t.Cleanup(func() { *updateGolden = old })
}

func TestGoldenUpdateWritesThenMatches(t *testing.T) {
	t.Chdir(t.TempDir())

	content := []byte("frozen wire bytes\n")

	setUpdateFlag(t, true)
	Golden(t, "nested/sample.golden", content)

	written, err := os.ReadFile(filepath.Join("testdata", "nested", "sample.golden"))
	if err != nil {
		t.Fatalf("golden file not written: %v", err)
	}
	if string(written) != string(content) {
		t.Fatalf("golden file content = %q, want %q", written, content)
	}

	setUpdateFlag(t, false)
	rec := &recordingTB{TB: t}
	Golden(rec, "nested/sample.golden", content)
	if len(rec.errors) != 0 || len(rec.fatals) != 0 {
		t.Fatalf("expected clean comparison, got errors=%v fatals=%v", rec.errors, rec.fatals)
	}
}

func TestGoldenMismatchReportsOffset(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := os.MkdirAll("testdata", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join("testdata", "x.golden"), []byte("abcdef"), 0o644); err != nil {
		t.Fatal(err)
	}

	rec := &recordingTB{TB: t}
	Golden(rec, "x.golden", []byte("abcXef"))

	if len(rec.errors) != 1 {
		t.Fatalf("expected exactly one error, got %v", rec.errors)
	}
	if !strings.Contains(rec.errors[0], "mismatch at byte 3") {
		t.Errorf("error should pinpoint byte 3: %s", rec.errors[0])
	}
}

func TestGoldenMissingFileHintsUpdate(t *testing.T) {
	t.Chdir(t.TempDir())

	rec := &recordingTB{TB: t}
	Golden(rec, "absent.golden", []byte("anything"))

	if len(rec.fatals) != 1 {
		t.Fatalf("expected exactly one fatal, got %v", rec.fatals)
	}
	if !strings.Contains(rec.fatals[0], "-update") {
		t.Errorf("missing-file error should mention -update: %s", rec.fatals[0])
	}
}

func TestFirstDiff(t *testing.T) {
	cases := []struct {
		name string
		a, b string
		want int
	}{
		{"identical", "same", "same", 4},
		{"first byte", "x", "y", 0},
		{"middle", "abcd", "abXd", 2},
		{"prefix shorter", "ab", "abcd", 2},
		{"both empty", "", "", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := firstDiff([]byte(tc.a), []byte(tc.b)); got != tc.want {
				t.Errorf("firstDiff(%q, %q) = %d, want %d", tc.a, tc.b, got, tc.want)
			}
		})
	}
}
