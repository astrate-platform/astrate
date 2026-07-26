// Command record captures upstream Astarte's REST error envelopes and writes
// them to test/conformance/upstream as a fixture plus a verbatim transcript
// (M12 plan, phase 04).
//
// It exists because Astrate's error envelopes were reconstructed rather than
// observed: phase 01 shipped a 404 envelope whose *shape* and status were
// inferred from the project's own conventions, with only the two detail
// strings actually seen upstream. A conformance test that hits the network to
// decide whether Astrate is correct fails whenever the upstream is down; a
// test that compares Astrate against a recorded transcript fails only when
// Astrate changes. So recording is a deliberate, separate act — this command —
// and comparison is an ordinary offline test.
//
// The transcript is the point. Every fixture entry carries the request that
// produced it, so a later reader can tell an observation from an assumption.
//
// Usage (needs a reachable upstream and the housekeeping key):
//
//	ASTARTE_UPSTREAM_URL=http://api.astarte.localhost:8080 \
//	ASTARTE_UPSTREAM_REALM=bench \
//	ASTARTE_UPSTREAM_HOUSEKEEPING_KEY=/path/to/housekeeping_private.pem \
//	go run ./record
//
// Nothing here imports Astrate: the recorder must not be able to describe
// upstream in Astrate's own terms.
package main

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// observation is one recorded request/response pair. The fields are ordered
// so the JSON reads as "we asked this, and got that".
type observation struct {
	Name string `json:"name"`
	Why  string `json:"why"`

	Method string `json:"method"`
	Path   string `json:"path"`
	Auth   string `json:"auth"` // how the Authorization header was formed

	Status int    `json:"status"`
	Body   string `json:"body"`
}

type fixture struct {
	// Provenance, so a stale recording announces itself rather than being
	// mistaken for a fresh one.
	RecordedAt      string `json:"recorded_at"`
	AstarteVersion  string `json:"astarte_version"`
	Realm           string `json:"realm"`
	RecorderCommand string `json:"recorder_command"`

	Observations []observation `json:"observations"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "record:", err)
		os.Exit(1)
	}
}

func env(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}

func run() error {
	base := strings.TrimSuffix(env("ASTARTE_UPSTREAM_URL", ""), "/")
	realm := env("ASTARTE_UPSTREAM_REALM", "bench")
	keyPath := env("ASTARTE_UPSTREAM_HOUSEKEEPING_KEY", "")
	version := env("ASTARTE_UPSTREAM_VERSION", "v1.2.0")
	if base == "" || keyPath == "" {
		return fmt.Errorf("ASTARTE_UPSTREAM_URL and ASTARTE_UPSTREAM_HOUSEKEEPING_KEY are required")
	}

	hkKey, err := loadKey(keyPath)
	if err != nil {
		return err
	}
	// A key upstream has never seen: its rejection is what distinguishes
	// "this token is not for you" from "this token is unreadable".
	strangerKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	hkToken, err := mint(hkKey, "a_ha")
	if err != nil {
		return err
	}
	strangerToken, err := mint(strangerKey, "a_ha")
	if err != nil {
		return err
	}

	// Each row pairs a rejection with the reason it is interesting. The
	// services are probed separately because upstream does not answer them
	// uniformly — that non-uniformity is precisely what Astrate has to match.
	type probe struct {
		name, why, method, path, auth, token string
		body                                 any
	}
	probes := []probe{
		// --- unmatched routes: what phase 01 assumed rather than observed ---
		{
			name: "appengine unmatched route", why: "phase 01 inferred the envelope shape and the 404 status; only the detail string was observed",
			method: "GET", path: "/appengine/v1/" + realm + "/no-such-route", auth: "none",
		},
		{
			name: "realmmanagement unmatched route", why: "same inference as appengine; recorded to confirm the two really do agree",
			method: "GET", path: "/realmmanagement/v1/" + realm + "/no-such-route", auth: "none",
		},
		{
			name: "pairing unmatched route", why: "pairing was observed to differ in wording (Page not found); confirm it still does",
			method: "GET", path: "/pairing/v1/" + realm + "/no-such-route", auth: "none",
		},
		{
			name: "housekeeping unmatched route", why: "Astrate deliberately leaves this as Go's plain-text 404; record what upstream actually does",
			method: "GET", path: "/housekeeping/v1/no-such-route", auth: "none",
		},
		{
			name: "unknown service prefix", why: "bounds the envelope: a prefix belonging to no service should not get a service envelope",
			method: "GET", path: "/no-such-service/v1/thing", auth: "none",
		},

		// --- token-shape rejections: the evidence behind deviation 5 ---
		{
			name: "authenticated route, no token", why: "deviation 5: the baseline every other token row is compared against",
			method: "GET", path: "/housekeeping/v1/realms", auth: "none",
		},
		{
			name: "authenticated route, malformed token", why: "deviation 5: upstream distinguishes an unreadable token from an unauthorized one; Astrate does not",
			method: "GET", path: "/housekeeping/v1/realms", auth: "Bearer not-a-jwt", token: "not-a-jwt",
		},
		{
			name: "authenticated route, well-formed token signed by an unknown key", why: "deviation 5: a readable token upstream cannot verify",
			method: "GET", path: "/housekeeping/v1/realms", auth: "Bearer <a_ha token signed by a key upstream has never seen>", token: strangerToken,
		},
		{
			name: "authenticated route, valid token", why: "the acceptance that makes the rejections above mean something: a blanket refusal would pass every row without it",
			method: "GET", path: "/housekeeping/v1/realms", auth: "Bearer <valid a_ha token>", token: hkToken,
		},

		// --- realm creation conflict: found while fixing bench provision ---
		{
			name: "create an existing realm", why: "upstream reports 422 existing_realm where Astrate reports 409; found by bench provision failing on a re-run",
			method: "POST", path: "/housekeeping/v1/realms", auth: "Bearer <valid a_ha token>", token: hkToken,
			body: map[string]any{"realm_name": realm, "jwt_public_key_pem": "-----BEGIN PUBLIC KEY-----\nrecorded\n-----END PUBLIC KEY-----\n"},
		},
	}

	out := fixture{
		RecordedAt:     time.Now().UTC().Format(time.RFC3339),
		AstarteVersion: version,
		Realm:          realm,
		RecorderCommand: "ASTARTE_UPSTREAM_URL=… ASTARTE_UPSTREAM_REALM=" + realm +
			" ASTARTE_UPSTREAM_HOUSEKEEPING_KEY=… go run ./record",
	}
	var transcript bytes.Buffer
	fmt.Fprintf(&transcript, "Recorded %s against upstream Astarte %s, realm %q.\n",
		out.RecordedAt, version, realm)
	transcript.WriteString("Every entry in rest-errors.json was produced by exactly one exchange below.\n" +
		"Tokens are described, never printed: a transcript is committed, a credential is not.\n")

	client := &http.Client{Timeout: 15 * time.Second}
	for _, p := range probes {
		var body io.Reader
		if p.body != nil {
			raw, err := json.Marshal(map[string]any{"data": p.body})
			if err != nil {
				return err
			}
			body = bytes.NewReader(raw)
		}
		req, err := http.NewRequest(p.method, base+p.path, body)
		if err != nil {
			return err
		}
		if p.body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		if p.token != "" {
			req.Header.Set("Authorization", "Bearer "+p.token)
		}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("%s: %w", p.name, err)
		}
		raw, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return err
		}

		got := strings.TrimSpace(string(raw))
		out.Observations = append(out.Observations, observation{
			Name: p.name, Why: p.why,
			Method: p.method, Path: p.path, Auth: p.auth,
			Status: resp.StatusCode, Body: got,
		})

		fmt.Fprintf(&transcript, "\n--- %s\n%s %s\nAuthorization: %s\n",
			p.name, p.method, p.path, p.auth)
		if p.body != nil {
			fmt.Fprintf(&transcript, "(request carries a JSON body)\n")
		}
		fmt.Fprintf(&transcript, "→ %d\n%s\n", resp.StatusCode, redact(got))
	}

	sort.SliceStable(out.Observations, func(i, j int) bool {
		return out.Observations[i].Name < out.Observations[j].Name
	})

	dir := env("ASTARTE_UPSTREAM_OUT", ".")

	enc, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "rest-errors.json"), append(enc, '\n'), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "rest-errors.transcript.txt"), transcript.Bytes(), 0o644); err != nil {
		return err
	}
	fmt.Printf("recorded %d observations into %s\n", len(out.Observations), dir)
	return nil
}

// redact keeps a realm's public key out of the committed transcript. It is not
// secret, but a PEM blob in a transcript invites someone to reuse it.
func redact(body string) string {
	if i := strings.Index(body, "-----BEGIN PUBLIC KEY-----"); i >= 0 {
		return body[:i] + "…(public key elided)…"
	}
	return body
}

func loadKey(path string) (crypto.Signer, error) {
	raw, err := os.ReadFile(path) //nolint:gosec // operator-supplied path to a local dev key
	if err != nil {
		return nil, err
	}
	blk, _ := pem.Decode(raw)
	if blk == nil {
		return nil, fmt.Errorf("%s: no PEM block", path)
	}
	if k, err := x509.ParseECPrivateKey(blk.Bytes); err == nil {
		return k, nil
	}
	if k, err := x509.ParsePKCS1PrivateKey(blk.Bytes); err == nil {
		return k, nil
	}
	k, err := x509.ParsePKCS8PrivateKey(blk.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%s: unsupported private key", path)
	}
	signer, ok := k.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("%s: key cannot sign", path)
	}
	return signer, nil
}

func mint(key crypto.Signer, claim string) (string, error) {
	var method jwt.SigningMethod = jwt.SigningMethodRS256
	if _, ok := key.(*ecdsa.PrivateKey); ok {
		method = jwt.SigningMethodES256
	}
	now := time.Now()
	return jwt.NewWithClaims(method, jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(time.Hour).Unix(),
		claim: []string{".*::.*"},
	}).SignedString(key)
}
