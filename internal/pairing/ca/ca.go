// Package ca implements Astrate's embedded per-realm certificate authority
// (docs/DESIGN.md §4.3), replacing upstream Astarte's CFSSL sidecar. Each
// realm owns one ECDSA P-256 CA; the CA issues short-lived client
// certificates against device CSRs, treating the CSR purely as proof of key
// possession: every requested attribute, including the subject, is ignored
// and overridden, exactly as upstream does.
package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"
)

const (
	// DefaultCALifetime is the self-signed realm CA validity
	// (docs/DESIGN.md §4.3: default 10 years).
	DefaultCALifetime = 10 * 365 * 24 * time.Hour

	// DefaultCertTTL is the client certificate validity
	// (docs/DESIGN.md §4.3: default 30 days).
	DefaultCertTTL = 30 * 24 * time.Hour

	// clockSkewBackdate is subtracted from NotBefore on every issued
	// certificate so devices with slightly-behind clocks accept fresh
	// certs immediately.
	clockSkewBackdate = 5 * time.Minute

	// serialBytes sizes certificate serials: 128-bit random
	// (docs/DESIGN.md §4.3).
	serialBytes = 16
)

// Sentinel errors. The pairing service maps them onto the wire causes
// (EXPIRED/INVALID) and HTTP statuses.
var (
	// ErrInvalidCSR reports a CSR that does not parse or whose
	// proof-of-possession signature does not verify.
	ErrInvalidCSR = errors.New("ca: invalid certificate signing request")
	// ErrCAExpired reports an issuance attempt outside the CA certificate's
	// own validity window.
	ErrCAExpired = errors.New("ca: realm CA certificate is not currently valid")
	// ErrCertificateExpired reports a client certificate outside its
	// validity window (wire cause EXPIRED).
	ErrCertificateExpired = errors.New("ca: certificate expired")
	// ErrCertificateInvalid reports a client certificate that does not
	// parse or does not chain to the realm CA (wire cause INVALID).
	ErrCertificateInvalid = errors.New("ca: certificate invalid")
)

// CA is a realm certificate authority: the CA certificate plus its private
// key. Immutable and safe for concurrent use.
type CA struct {
	cert    *x509.Certificate
	key     *ecdsa.PrivateKey
	certPEM string
	keyDER  []byte
}

// Generate creates a fresh self-signed ECDSA P-256 realm CA. A
// zero lifetime selects DefaultCALifetime; negative lifetimes are allowed
// (they produce an already-expired CA, used by tests).
func Generate(realm string, lifetime time.Duration) (*CA, error) {
	if lifetime == 0 {
		lifetime = DefaultCALifetime
	}
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("ca: generating CA key: %w", err)
	}
	serial, err := randomSerial()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   "Astrate Realm " + realm + " CA",
			Organization: []string{"Astrate"},
		},
		NotBefore:             now.Add(-clockSkewBackdate),
		NotAfter:              now.Add(lifetime),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true, // issues leaves only
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, fmt.Errorf("ca: self-signing CA certificate: %w", err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, fmt.Errorf("ca: re-parsing CA certificate: %w", err)
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("ca: marshalling CA key: %w", err)
	}
	return &CA{
		cert:    cert,
		key:     key,
		certPEM: encodeCertPEM(der),
		keyDER:  keyDER,
	}, nil
}

// Load reconstructs a CA from its stored material: the PEM certificate and
// the PKCS#8 DER private key (the plaintext that store.KeySealer sealed).
// Operator-provided CA imports go through the same path.
func Load(certPEM string, keyDER []byte) (*CA, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("ca: CA material is not a PEM certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ca: parsing CA certificate: %w", err)
	}
	parsed, err := x509.ParsePKCS8PrivateKey(keyDER)
	if err != nil {
		return nil, fmt.Errorf("ca: parsing CA private key: %w", err)
	}
	key, ok := parsed.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("ca: CA private key is %T, want *ecdsa.PrivateKey", parsed)
	}
	if !key.PublicKey.Equal(cert.PublicKey) {
		return nil, errors.New("ca: CA private key does not match CA certificate")
	}
	return &CA{cert: cert, key: key, certPEM: certPEM, keyDER: keyDER}, nil
}

// CertificatePEM returns the CA certificate in PEM form (the `ca_crt`
// delivered to devices).
func (c *CA) CertificatePEM() string { return c.certPEM }

// PrivateKeyDER returns the PKCS#8 DER encoding of the CA private key — the
// plaintext handed to store.KeySealer for at-rest encryption.
func (c *CA) PrivateKeyDER() []byte { return c.keyDER }

// SignCSR issues a client certificate for a device against csrPEM
// (docs/DESIGN.md §4.3):
//
//   - Subject CN is forced to "<realm>/<deviceID>"; everything the CSR
//     requested (subject, extensions, attributes) is ignored;
//   - serial is 128-bit random; KeyUsage is digitalSignature only;
//     ExtKeyUsage is clientAuth;
//   - validity is now-clockSkewBackdate .. now+ttl, clamped to the CA's own
//     NotAfter; a non-positive ttl selects DefaultCertTTL;
//   - issuance is refused outside the CA certificate's validity window.
//
// It returns the certificate PEM, its serial (decimal string) and its
// authority key identifier (lowercase hex) for the device row's
// latest-certificate trail.
func (c *CA) SignCSR(csrPEM, realm, deviceID string, ttl time.Duration) (certPEM, serial, aki string, err error) {
	csr, err := parseCSR(csrPEM)
	if err != nil {
		return "", "", "", err
	}
	if err := csr.CheckSignature(); err != nil {
		return "", "", "", fmt.Errorf("%w: proof-of-possession signature: %v", ErrInvalidCSR, err)
	}

	now := time.Now()
	if now.Before(c.cert.NotBefore) || now.After(c.cert.NotAfter) {
		return "", "", "", ErrCAExpired
	}
	if ttl <= 0 {
		ttl = DefaultCertTTL
	}
	notAfter := now.Add(ttl)
	if notAfter.After(c.cert.NotAfter) {
		notAfter = c.cert.NotAfter
	}

	serialNum, err := randomSerial()
	if err != nil {
		return "", "", "", err
	}
	template := &x509.Certificate{
		SerialNumber:          serialNum,
		Subject:               pkix.Name{CommonName: realm + "/" + deviceID},
		NotBefore:             now.Add(-clockSkewBackdate),
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, c.cert, csr.PublicKey, c.key)
	if err != nil {
		return "", "", "", fmt.Errorf("ca: signing client certificate: %w", err)
	}
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		return "", "", "", fmt.Errorf("ca: re-parsing client certificate: %w", err)
	}
	return encodeCertPEM(der), leaf.SerialNumber.String(), hex.EncodeToString(leaf.AuthorityKeyId), nil
}

// Verify checks a client certificate against this realm CA at the given
// instant (the `credentials/verify` endpoint backend). On success it returns
// the certificate's NotAfter. Failures wrap ErrCertificateExpired (outside
// the validity window) or ErrCertificateInvalid (parse failure, foreign CA,
// wrong usage) — the precedence is expiry first, so an expired certificate
// reports EXPIRED even when other problems coexist.
func (c *CA) Verify(certPEM string, at time.Time) (until time.Time, err error) {
	cert, err := ParseCertificatePEM(certPEM)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %v", ErrCertificateInvalid, err)
	}
	if at.After(cert.NotAfter) || at.Before(cert.NotBefore) {
		return cert.NotAfter, ErrCertificateExpired
	}

	roots := x509.NewCertPool()
	roots.AddCert(c.cert)
	_, err = cert.Verify(x509.VerifyOptions{
		Roots:       roots,
		CurrentTime: at,
		KeyUsages:   []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})
	if err != nil {
		var invalid x509.CertificateInvalidError
		if errors.As(err, &invalid) && invalid.Reason == x509.Expired {
			return cert.NotAfter, ErrCertificateExpired
		}
		return cert.NotAfter, fmt.Errorf("%w: %v", ErrCertificateInvalid, err)
	}
	return cert.NotAfter, nil
}

// ParseCertificatePEM decodes a single PEM-encoded X.509 certificate.
func ParseCertificatePEM(certPEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("ca: not a PEM certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ca: parsing certificate: %w", err)
	}
	return cert, nil
}

// parseCSR decodes a PEM certificate signing request.
func parseCSR(csrPEM string) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode([]byte(csrPEM))
	if block == nil || (block.Type != "CERTIFICATE REQUEST" && block.Type != "NEW CERTIFICATE REQUEST") {
		return nil, fmt.Errorf("%w: not a PEM certificate request", ErrInvalidCSR)
	}
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCSR, err)
	}
	return csr, nil
}

// randomSerial draws a uniform non-zero 128-bit serial.
func randomSerial() (*big.Int, error) {
	for {
		buf := make([]byte, serialBytes)
		if _, err := rand.Read(buf); err != nil {
			return nil, fmt.Errorf("ca: gathering serial randomness: %w", err)
		}
		serial := new(big.Int).SetBytes(buf)
		if serial.Sign() > 0 {
			return serial, nil
		}
	}
}

// encodeCertPEM wraps DER certificate bytes in PEM.
func encodeCertPEM(der []byte) string {
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
}
