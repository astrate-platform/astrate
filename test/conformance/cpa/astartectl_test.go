package cpa

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
// (docs/ROADMAP.md §0.3: pinned, upgraded deliberately).
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
		out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
		if err != nil {
			return err
		}
		if _, err := io.Copy(out, tr); err != nil {
			_ = out.Close()
			return err
		}
		return out.Close()
	}
}

// runAstartectl executes the binary and returns trimmed stdout, failing the
// test on a non-zero exit (with stderr in the failure message).
func runAstartectl(t *testing.T, bin string, args ...string) string {
	t.Helper()
	cmd := exec.Command(bin, args...)
	// An isolated config dir silences "config not found" warnings and keeps
	// the user's real astartectl contexts out of the run.
	cmd.Args = append(cmd.Args, "--config-dir", t.TempDir())

	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			t.Fatalf("astartectl %s failed: %v\nstderr: %s", strings.Join(args, " "), err, exitErr.Stderr)
		}
		t.Fatalf("astartectl %s failed: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out))
}
