package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"testing"
	"time"
)

// newCSR builds a CSR with a deliberately hostile subject: every field must
// be ignored by SignCSR.
func newCSR(t *testing.T, key any) string {
	t.Helper()
	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:         "evilrealm/evil-device",
			Organization:       []string{"Evil Corp"},
			OrganizationalUnit: []string{"Impersonation Dept"},
			Country:            []string{"XX"},
		},
		DNSNames:       []string{"broker.evil.example"},
		EmailAddresses: []string{"root@evil.example"},
	}
	der, err := x509.CreateCertificateRequest(rand.Reader, template, key)
	if err != nil {
		t.Fatalf("creating CSR: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: der}))
}

func newECKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

func mustGenerate(t *testing.T, realm string, lifetime time.Duration) *CA {
	t.Helper()
	c, err := Generate(realm, lifetime)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	return c
}

func TestGenerateCACertificate(t *testing.T) {
	c := mustGenerate(t, "test", 0)

	cert, err := ParseCertificatePEM(c.CertificatePEM())
	if err != nil {
		t.Fatalf("parsing CA PEM: %v", err)
	}
	if !cert.IsCA || !cert.BasicConstraintsValid {
		t.Error("CA certificate must carry CA basic constraints")
	}
	if cert.KeyUsage&x509.KeyUsageCertSign == 0 {
		t.Error("CA certificate must have certSign key usage")
	}
	if len(cert.SubjectKeyId) == 0 {
		t.Error("CA certificate must have a subject key identifier")
	}
	if got := cert.Subject.CommonName; got != "Astrate Realm test CA" {
		t.Errorf("CA CN: got %q", got)
	}
	wantAfter := time.Now().Add(DefaultCALifetime - 24*time.Hour)
	if cert.NotAfter.Before(wantAfter) {
		t.Errorf("CA NotAfter %v too early for the 10-year default", cert.NotAfter)
	}
	if _, ok := cert.PublicKey.(*ecdsa.PublicKey); !ok {
		t.Errorf("CA key: got %T, want ECDSA P-256", cert.PublicKey)
	}
}

func TestLoadRoundTrip(t *testing.T) {
	c := mustGenerate(t, "test", 0)

	loaded, err := Load(c.CertificatePEM(), c.PrivateKeyDER())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	// The reloaded CA must be able to issue and verify.
	certPEM, _, _, err := loaded.SignCSR(newCSR(t, newECKey(t)), "test", "device", time.Hour)
	if err != nil {
		t.Fatalf("SignCSR after Load: %v", err)
	}
	if _, err := loaded.Verify(certPEM, time.Now()); err != nil {
		t.Fatalf("Verify after Load: %v", err)
	}

	// Mismatched key/cert pairs are rejected.
	other := mustGenerate(t, "other", 0)
	if _, err := Load(c.CertificatePEM(), other.PrivateKeyDER()); err == nil {
		t.Error("Load with mismatched key must fail")
	}
	if _, err := Load("garbage", c.PrivateKeyDER()); err == nil {
		t.Error("Load with garbage certificate must fail")
	}
}

func TestSignCSRFieldAssertions(t *testing.T) {
	c := mustGenerate(t, "test", 0)
	deviceID := "h4-Dx_RYTU-RbpDOTabhRg"

	certPEM, serial, aki, err := c.SignCSR(newCSR(t, newECKey(t)), "test", deviceID, time.Hour)
	if err != nil {
		t.Fatalf("SignCSR: %v", err)
	}
	cert, err := ParseCertificatePEM(certPEM)
	if err != nil {
		t.Fatalf("parsing issued certificate: %v", err)
	}

	// CN forced to <realm>/<device_id>; the hostile CSR subject is gone.
	if got := cert.Subject.CommonName; got != "test/"+deviceID {
		t.Errorf("CN: got %q, want %q", got, "test/"+deviceID)
	}
	if len(cert.Subject.Organization) != 0 || len(cert.Subject.Country) != 0 ||
		len(cert.Subject.OrganizationalUnit) != 0 {
		t.Errorf("CSR subject fields leaked into the certificate: %v", cert.Subject)
	}
	if len(cert.DNSNames) != 0 || len(cert.EmailAddresses) != 0 {
		t.Errorf("CSR SANs leaked into the certificate: %v / %v", cert.DNSNames, cert.EmailAddresses)
	}

	// 128-bit random serial, reported as its decimal string.
	if cert.SerialNumber.Sign() <= 0 || cert.SerialNumber.BitLen() > 128 {
		t.Errorf("serial out of 128-bit range: %v", cert.SerialNumber)
	}
	if serial != cert.SerialNumber.String() {
		t.Errorf("returned serial %q != certificate serial %q", serial, cert.SerialNumber.String())
	}

	// KeyUsage digitalSignature only; EKU clientAuth only.
	if cert.KeyUsage != x509.KeyUsageDigitalSignature {
		t.Errorf("KeyUsage: got %v, want digitalSignature only", cert.KeyUsage)
	}
	if len(cert.ExtKeyUsage) != 1 || cert.ExtKeyUsage[0] != x509.ExtKeyUsageClientAuth {
		t.Errorf("ExtKeyUsage: got %v, want [clientAuth]", cert.ExtKeyUsage)
	}
	if cert.IsCA {
		t.Error("issued certificate must not be a CA")
	}

	// AKI matches the CA's subject key identifier.
	if aki != hex.EncodeToString(c.cert.SubjectKeyId) {
		t.Errorf("AKI: got %q, want CA SKI %q", aki, hex.EncodeToString(c.cert.SubjectKeyId))
	}

	// Chain verifies to the CA.
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM([]byte(c.CertificatePEM())) {
		t.Fatal("appending CA PEM to pool")
	}
	if _, err := cert.Verify(x509.VerifyOptions{
		Roots:     roots,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}); err != nil {
		t.Errorf("issued certificate does not chain to the CA: %v", err)
	}

	// TTL honoured (1 hour, within a minute of tolerance).
	if d := time.Until(cert.NotAfter); d > time.Hour+time.Minute || d < time.Hour-time.Minute {
		t.Errorf("NotAfter %v not ~1h away", cert.NotAfter)
	}
}

func TestSignCSRRSAKey(t *testing.T) {
	c := mustGenerate(t, "test", 0)
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	certPEM, _, _, err := c.SignCSR(newCSR(t, rsaKey), "test", "dev", time.Hour)
	if err != nil {
		t.Fatalf("SignCSR with RSA CSR: %v", err)
	}
	cert, err := ParseCertificatePEM(certPEM)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := cert.PublicKey.(*rsa.PublicKey); !ok {
		t.Errorf("issued key: got %T, want RSA from the CSR", cert.PublicKey)
	}
}

func TestSignCSRTTLClampedToCA(t *testing.T) {
	// A CA with 1 hour left cannot issue a 30-day certificate.
	c := mustGenerate(t, "test", time.Hour)
	certPEM, _, _, err := c.SignCSR(newCSR(t, newECKey(t)), "test", "dev", DefaultCertTTL)
	if err != nil {
		t.Fatalf("SignCSR: %v", err)
	}
	cert, err := ParseCertificatePEM(certPEM)
	if err != nil {
		t.Fatal(err)
	}
	if cert.NotAfter.After(c.cert.NotAfter) {
		t.Errorf("leaf NotAfter %v exceeds CA NotAfter %v", cert.NotAfter, c.cert.NotAfter)
	}
}

func TestSignCSRRejections(t *testing.T) {
	c := mustGenerate(t, "test", 0)

	t.Run("ExpiredCA", func(t *testing.T) {
		expired := mustGenerate(t, "test", -time.Hour)
		_, _, _, err := expired.SignCSR(newCSR(t, newECKey(t)), "test", "dev", time.Hour)
		if !errors.Is(err, ErrCAExpired) {
			t.Errorf("expired CA: got %v, want ErrCAExpired", err)
		}
	})

	t.Run("GarbageCSR", func(t *testing.T) {
		for _, csr := range []string{"", "garbage", "-----BEGIN CERTIFICATE REQUEST-----\nZ2FyYmFnZQ==\n-----END CERTIFICATE REQUEST-----\n"} {
			if _, _, _, err := c.SignCSR(csr, "test", "dev", time.Hour); !errors.Is(err, ErrInvalidCSR) {
				t.Errorf("SignCSR(%.20q): got %v, want ErrInvalidCSR", csr, err)
			}
		}
	})

	t.Run("TamperedCSRSignature", func(t *testing.T) {
		// A CSR whose embedded public key was swapped after signing fails
		// proof-of-possession.
		template := &x509.CertificateRequest{Subject: pkix.Name{CommonName: "x"}}
		der, err := x509.CreateCertificateRequest(rand.Reader, template, newECKey(t))
		if err != nil {
			t.Fatal(err)
		}
		der[len(der)-5] ^= 0xff // corrupt the signature bytes
		csrPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: der}))
		if _, _, _, err := c.SignCSR(csrPEM, "test", "dev", time.Hour); !errors.Is(err, ErrInvalidCSR) {
			t.Errorf("tampered CSR: got %v, want ErrInvalidCSR", err)
		}
	})
}

func TestVerify(t *testing.T) {
	c := mustGenerate(t, "test", 0)
	certPEM, _, _, err := c.SignCSR(newCSR(t, newECKey(t)), "test", "dev", time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Valid", func(t *testing.T) {
		until, err := c.Verify(certPEM, time.Now())
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		cert, _ := ParseCertificatePEM(certPEM)
		if !until.Equal(cert.NotAfter) {
			t.Errorf("until: got %v, want NotAfter %v", until, cert.NotAfter)
		}
	})

	t.Run("Expired", func(t *testing.T) {
		// Time-travel past NotAfter instead of sleeping.
		_, err := c.Verify(certPEM, time.Now().Add(2*time.Hour))
		if !errors.Is(err, ErrCertificateExpired) {
			t.Errorf("expired: got %v, want ErrCertificateExpired", err)
		}
	})

	t.Run("ForeignCA", func(t *testing.T) {
		foreign := mustGenerate(t, "test", 0) // same realm name, different key
		foreignCert, _, _, err := foreign.SignCSR(newCSR(t, newECKey(t)), "test", "dev", time.Hour)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := c.Verify(foreignCert, time.Now()); !errors.Is(err, ErrCertificateInvalid) {
			t.Errorf("foreign CA: got %v, want ErrCertificateInvalid", err)
		}
	})

	t.Run("Garbage", func(t *testing.T) {
		if _, err := c.Verify("not a certificate", time.Now()); !errors.Is(err, ErrCertificateInvalid) {
			t.Errorf("garbage: got %v, want ErrCertificateInvalid", err)
		}
	})
}

func TestSerialUniqueness10k(t *testing.T) {
	if testing.Short() {
		t.Skip("10k issuance draw skipped in -short mode")
	}
	c := mustGenerate(t, "test", 0)
	csr := newCSR(t, newECKey(t)) // key reuse is fine; serials must still differ

	seen := make(map[string]struct{}, 10_000)
	for i := 0; i < 10_000; i++ {
		_, serial, _, err := c.SignCSR(csr, "test", "dev", time.Hour)
		if err != nil {
			t.Fatalf("SignCSR #%d: %v", i, err)
		}
		if _, dup := seen[serial]; dup {
			t.Fatalf("duplicate serial after %d issuances: %s", i, serial)
		}
		seen[serial] = struct{}{}
	}
}
