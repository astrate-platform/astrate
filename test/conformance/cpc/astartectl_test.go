package cpc

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// astartectlVersion is the pinned conformance release
// (docs/ROADMAP.md §0.3: pinned, upgraded deliberately). Kept in sync with
// the cpa harness; the binary is shared through the user cache.
const astartectlVersion = "26.5.0"

// ensureAstartectl locates the pinned astartectl binary: the ASTARTECTL_BIN
// override first, then the user cache, downloading the release artifact on a
// cache miss.
func ensureAstartectl(t *testing.T) string {
	t.Helper()
	if bin := os.Getenv("ASTARTECTL_BIN"); bin != "" {
		return bin
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Fatalf("resolving user cache dir: %v", err)
	}
	dir := filepath.Join(cacheDir, "astrate-conformance")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(dir, "astartectl-v"+astartectlVersion)
	if _, err := os.Stat(bin); err == nil {
		return bin
	}

	url, err := astartectlURL()
	if err != nil {
		t.Skipf("astartectl: %v (set ASTARTECTL_BIN to run this checkpoint)", err)
	}
	t.Logf("downloading pinned astartectl %s from %s", astartectlVersion, url)
	if err := downloadAstartectl(url, bin); err != nil {
		t.Fatalf("downloading astartectl: %v", err)
	}
	return bin
}

// astartectlURL maps GOOS/GOARCH onto the release artifact name.
func astartectlURL() (string, error) {
	osName, ok := map[string]string{"darwin": "macOS", "linux": "Linux"}[runtime.GOOS]
	if !ok {
		return "", fmt.Errorf("no pinned astartectl artifact for GOOS=%s", runtime.GOOS)
	}
	archName, ok := map[string]string{"amd64": "x86_64", "arm64": "arm64"}[runtime.GOARCH]
	if !ok {
		return "", fmt.Errorf("no pinned astartectl artifact for GOARCH=%s", runtime.GOARCH)
	}
	return fmt.Sprintf(
		"https://github.com/astarte-platform/astartectl/releases/download/v%s/astartectl_%s_%s_%s.tar.gz",
		astartectlVersion, astartectlVersion, osName, archName), nil
}

// downloadAstartectl fetches the tar.gz artifact and extracts the
// `astartectl` member to dest.
func downloadAstartectl(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %s", url, resp.Status)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("artifact has no astartectl member")
		}
		if err != nil {
			return err
		}
		if filepath.Base(hdr.Name) != "astartectl" || hdr.Typeflag != tar.TypeReg {
			continue
		}
		// Write to a temp file and rename into place: cpa and cpc race this
		// download in parallel `go test` processes, and exec-ing a binary
		// another process still holds open for writing is ETXTBSY on Linux.
		// Rename is atomic, so an exec only ever sees a complete inode.
		tmp, err := os.CreateTemp(filepath.Dir(dest), filepath.Base(dest)+".tmp-*")
		if err != nil {
			return err
		}
		if _, err := io.Copy(tmp, tr); err != nil { //nolint:gosec // pinned release artifact
			_ = tmp.Close()
			_ = os.Remove(tmp.Name())
			return err
		}
		if err := tmp.Chmod(0o755); err != nil {
			_ = tmp.Close()
			_ = os.Remove(tmp.Name())
			return err
		}
		if err := tmp.Close(); err != nil {
			_ = os.Remove(tmp.Name())
			return err
		}
		return os.Rename(tmp.Name(), dest)
	}
}

// runAstartectl executes the binary and returns trimmed stdout, failing the
// test on a non-zero exit (with stderr in the failure message). An isolated
// --config-dir keeps the user's real astartectl contexts out of the run.
func runAstartectl(t *testing.T, bin string, args ...string) string {
	t.Helper()
	out, err := tryAstartectl(bin, args...)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			t.Fatalf("astartectl %s failed: %v\nstderr: %s", strings.Join(args, " "), err, exitErr.Stderr)
		}
		t.Fatalf("astartectl %s failed: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(out)
}

// tryAstartectl runs the binary and returns trimmed stdout and any error,
// letting callers assert on the failure path.
func tryAstartectl(bin string, args ...string) (string, error) {
	cmd := exec.Command(bin, args...)
	cmd.Args = append(cmd.Args, "--config-dir", os.TempDir())
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}
