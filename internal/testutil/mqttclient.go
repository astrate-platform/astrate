package testutil

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// mqttConnectTimeout bounds CONNECT round trips in tests.
const mqttConnectTimeout = 10 * time.Second

// ServerTLSCert generates a self-signed server certificate for the broker's
// TLS listener (valid for localhost/127.0.0.1/::1) and returns it together
// with a root pool that trusts it, for client-side verification.
func ServerTLSCert(t testing.TB) (tls.Certificate, *x509.CertPool) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("testutil: generating server key: %v", err)
	}
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "astrate-test-broker"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true, // self-signed: acts as its own root
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("testutil: self-signing server certificate: %v", err)
	}
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("testutil: re-parsing server certificate: %v", err)
	}
	pool := x509.NewCertPool()
	pool.AddCert(leaf)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key, Leaf: leaf}, pool
}

// DeviceCSR generates a fresh P-256 device key and a PEM CSR for it (the
// CSR's subject is irrelevant: the pairing CA overrides everything).
func DeviceCSR(t testing.TB) (*ecdsa.PrivateKey, string) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("testutil: generating device key: %v", err)
	}
	der, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: "ignored"},
	}, key)
	if err != nil {
		t.Fatalf("testutil: creating device CSR: %v", err)
	}
	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: der})
	return key, string(csrPEM)
}

// DeviceTLSConfig assembles the client-side mTLS configuration from an
// issued certificate PEM, its private key, and the pool trusting the
// broker's server certificate.
func DeviceTLSConfig(t testing.TB, clientCertPEM string, key *ecdsa.PrivateKey, roots *x509.CertPool) *tls.Config {
	t.Helper()

	block, _ := pem.Decode([]byte(clientCertPEM))
	if block == nil || block.Type != "CERTIFICATE" {
		t.Fatalf("testutil: client certificate is not PEM")
	}
	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		RootCAs:      roots,
		Certificates: []tls.Certificate{{Certificate: [][]byte{block.Bytes}, PrivateKey: key}},
	}
}

// MQTTTryConnect dials the broker and waits for the CONNACK. It returns the
// connected client and the CONNACK's session-present flag, or an error when
// the broker refuses (or drops) the connection. brokerURL uses paho schemes:
// "ssl://host:port" (TLS) or "tcp://host:port". Extra option tweaks (e.g.
// SetDefaultPublishHandler) apply before connecting. On success, disconnect
// is registered on t.Cleanup.
func MQTTTryConnect(t testing.TB, brokerURL, clientID string, cleanSession bool, tlsCfg *tls.Config, tweaks ...func(*paho.ClientOptions)) (paho.Client, bool, error) {
	t.Helper()

	opts := paho.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		SetProtocolVersion(4). // MQTT 3.1.1, what the Astarte SDKs speak
		SetCleanSession(cleanSession).
		SetAutoReconnect(false).
		SetConnectRetry(false).
		SetResumeSubs(false).
		SetKeepAlive(30 * time.Second).
		SetConnectTimeout(mqttConnectTimeout)
	if tlsCfg != nil {
		opts.SetTLSConfig(tlsCfg)
	}
	for _, tweak := range tweaks {
		tweak(opts)
	}

	client := paho.NewClient(opts)
	token := client.Connect()
	if !token.WaitTimeout(mqttConnectTimeout) {
		client.Disconnect(0)
		return nil, false, context.DeadlineExceeded
	}
	if err := token.Error(); err != nil {
		return nil, false, err
	}
	ct, ok := token.(*paho.ConnectToken)
	if !ok {
		t.Fatalf("testutil: unexpected token type %T", token)
	}
	t.Cleanup(func() {
		if client.IsConnected() {
			client.Disconnect(100)
		}
	})
	return client, ct.SessionPresent(), nil
}

// MQTTConnect is MQTTTryConnect that fails the test on connection errors.
func MQTTConnect(t testing.TB, brokerURL, clientID string, cleanSession bool, tlsCfg *tls.Config, tweaks ...func(*paho.ClientOptions)) (paho.Client, bool) {
	t.Helper()

	client, sessionPresent, err := MQTTTryConnect(t, brokerURL, clientID, cleanSession, tlsCfg, tweaks...)
	if err != nil {
		t.Fatalf("testutil: connecting to %s as %q: %v", brokerURL, clientID, err)
	}
	return client, sessionPresent
}

// WaitToken waits for a paho token with a timeout and fails the test on
// token errors.
func WaitToken(t testing.TB, token paho.Token, timeout time.Duration) {
	t.Helper()

	if !token.WaitTimeout(timeout) {
		t.Fatalf("testutil: MQTT operation timed out after %s", timeout)
	}
	if err := token.Error(); err != nil {
		t.Fatalf("testutil: MQTT operation failed: %v", err)
	}
}
