package main

import (
	"crypto/x509"
	"net"
	"slices"
	"testing"
	"time"
)

// TestSelfSignedDevCertFields pins the shape of the throwaway mTLS identity
// insecure_dev_mode boots the broker's TLS listener on: the validity window, the
// loopback names a device verifies, the serverAuth usage, and that the leaf is
// self-signed. Each of these is load-bearing — a plain NotBefore of now rejects
// a device whose clock runs behind ours, and dropping the loopback IPs breaks
// hostname verification for a device dialling mqtts://127.0.0.1:8883 — so the
// fields are asserted exactly rather than only through the handshake that would
// otherwise cover them.
func TestSelfSignedDevCertFields(t *testing.T) {
	before := time.Now()
	cert, err := selfSignedDevCert()
	after := time.Now()
	if err != nil {
		t.Fatalf("selfSignedDevCert: %v", err)
	}
	if len(cert.Certificate) != 1 {
		t.Fatalf("chain length = %d, want a single self-signed leaf", len(cert.Certificate))
	}
	if cert.PrivateKey == nil {
		t.Error("PrivateKey is nil: the broker listener has nothing to present")
	}
	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatalf("parsing the generated certificate: %v", err)
	}

	t.Run("ValidityWindow", func(t *testing.T) {
		// Backdated by an hour: a device whose clock is behind the broker's
		// would otherwise see a not-yet-valid certificate. Truncation to whole
		// seconds is why the window is a minute wide rather than exact.
		if age := before.Sub(leaf.NotBefore); age < 30*time.Minute || age > 2*time.Hour {
			t.Errorf("NotBefore is %s back, want roughly 1h (a plain now() would reject a device clock behind ours)", age.Round(time.Second))
		}
		if left := leaf.NotAfter.Sub(after); left < 365*24*time.Hour-time.Minute || left > 365*24*time.Hour+time.Minute {
			t.Errorf("NotAfter is %s out, want 365d", left.Round(time.Second))
		}
	})

	t.Run("LoopbackNames", func(t *testing.T) {
		if !slices.Contains(leaf.DNSNames, "localhost") {
			t.Errorf("DNSNames = %v, want localhost", leaf.DNSNames)
		}
		for _, want := range []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback} {
			if !slices.ContainsFunc(leaf.IPAddresses, want.Equal) {
				t.Errorf("IPAddresses = %v, want %s", leaf.IPAddresses, want)
			}
		}
		// What a device actually does when it dials the dev broker.
		for _, host := range []string{"localhost", "127.0.0.1", "::1"} {
			if err := leaf.VerifyHostname(host); err != nil {
				t.Errorf("VerifyHostname(%q): %v", host, err)
			}
		}
	})

	t.Run("ServerAuthOnly", func(t *testing.T) {
		if !slices.Contains(leaf.ExtKeyUsage, x509.ExtKeyUsageServerAuth) {
			t.Errorf("ExtKeyUsage = %v, want serverAuth", leaf.ExtKeyUsage)
		}
		if len(leaf.ExtKeyUsage) != 1 {
			t.Errorf("ExtKeyUsage = %v, want only serverAuth", leaf.ExtKeyUsage)
		}
		if leaf.KeyUsage&x509.KeyUsageDigitalSignature == 0 {
			t.Errorf("KeyUsage = %b, want digitalSignature", leaf.KeyUsage)
		}
	})

	t.Run("SelfSignedChainVerifies", func(t *testing.T) {
		// Self-signed: the leaf signed itself, so it verifies against a pool
		// holding nothing but itself.
		if err := leaf.CheckSignature(leaf.SignatureAlgorithm, leaf.RawTBSCertificate, leaf.Signature); err != nil {
			t.Errorf("leaf signature does not verify under its own key: %v", err)
		}
		if leaf.Issuer.String() != leaf.Subject.String() {
			t.Errorf("issuer = %q, want the subject %q (self-signed)", leaf.Issuer, leaf.Subject)
		}
		roots := x509.NewCertPool()
		roots.AddCert(leaf)
		for _, host := range []string{"localhost", "127.0.0.1", "::1"} {
			_, err := leaf.Verify(x509.VerifyOptions{
				Roots:     roots,
				DNSName:   host,
				KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			})
			if err != nil {
				t.Errorf("Verify(%q) as a self-signed server cert: %v", host, err)
			}
		}
	})
}
