package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// randomDeviceID returns a random Astarte device ID: 16 bytes,
// base64url-encoded without padding (22 characters).
func randomDeviceID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

// mqttDevice is one simulated device: an mTLS identity plus a live paho
// session that has completed the Astarte connect handshake.
type mqttDevice struct {
	realm  string
	id     string
	base   string // "<realm>/<device_id>" — base topic and certificate CN
	tlsCfg *tls.Config
	client paho.Client
}

// newIdentity generates a key + CSR, obtains the client certificate through
// the pairing API, and returns a ready-to-connect device. ECDSA P-256 by
// default (fast enough to generate thousands); rsa flips to RSA-2048 for
// parity with SDKs that only do RSA.
func newIdentity(c *client, ep Endpoints, realm string, dev Device, useRSA, skipVerify bool) (*mqttDevice, error) {
	var (
		key any
		err error
	)
	if useRSA {
		key, err = rsa.GenerateKey(rand.Reader, 2048)
	} else {
		key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}
	if err != nil {
		return nil, err
	}

	cn := realm + "/" + dev.ID
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: cn},
	}, key)
	if err != nil {
		return nil, err
	}
	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})

	certPEM, err := c.obtainCertificate(ep, realm, dev, string(csrPEM))
	if err != nil {
		return nil, fmt.Errorf("device %s: obtaining certificate: %w", dev.ID, err)
	}

	keyDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
	pair, err := tls.X509KeyPair([]byte(certPEM), keyPEM)
	if err != nil {
		return nil, fmt.Errorf("device %s: assembling TLS pair: %w", dev.ID, err)
	}

	return &mqttDevice{
		realm: realm,
		id:    dev.ID,
		base:  cn,
		tlsCfg: &tls.Config{
			Certificates:       []tls.Certificate{pair},
			InsecureSkipVerify: skipVerify, //nolint:gosec // both dev stacks use self-signed broker certs
			MinVersion:         tls.VersionTLS12,
		},
	}, nil
}

// connect dials the broker and runs the Astarte session handshake:
// introspection, then emptyCache. It returns the time the CONNACK took.
func (d *mqttDevice) connect(brokerURL string, timeout time.Duration, interfaces map[string]string) (time.Duration, error) {
	// The pairing API advertises mqtts://; paho speaks ssl://.
	dial := strings.Replace(brokerURL, "mqtts://", "ssl://", 1)
	dial = strings.Replace(dial, "mqtt://", "tcp://", 1)

	opts := paho.NewClientOptions().
		AddBroker(dial).
		SetClientID(d.base).
		SetProtocolVersion(4). // MQTT 3.1.1, what the SDKs speak
		SetCleanSession(true).
		SetAutoReconnect(false).
		SetKeepAlive(60 * time.Second).
		SetConnectTimeout(timeout).
		SetTLSConfig(d.tlsCfg)

	d.client = paho.NewClient(opts)
	start := time.Now()
	tok := d.client.Connect()
	if !tok.WaitTimeout(timeout) {
		return 0, fmt.Errorf("device %s: CONNACK timeout after %s", d.id, timeout)
	}
	if tok.Error() != nil {
		return 0, fmt.Errorf("device %s: connect: %w", d.id, tok.Error())
	}
	connLatency := time.Since(start)

	if len(interfaces) > 0 {
		var parts []string
		for name, ver := range interfaces {
			parts = append(parts, name+":"+ver)
		}
		intro := d.client.Publish(d.base, 2, false, strings.Join(parts, ";"))
		if !intro.WaitTimeout(timeout) || intro.Error() != nil {
			return 0, fmt.Errorf("device %s: introspection publish failed: %v", d.id, intro.Error())
		}
		ec := d.client.Publish(d.base+"/control/emptyCache", 2, false, "1")
		if !ec.WaitTimeout(timeout) || ec.Error() != nil {
			return 0, fmt.Errorf("device %s: emptyCache publish failed: %v", d.id, ec.Error())
		}
	}
	return connLatency, nil
}

func (d *mqttDevice) disconnect() {
	if d.client != nil && d.client.IsConnected() {
		d.client.Disconnect(250)
	}
}

// publishDouble sends one individual-datastream double with an explicit
// timestamp and returns the paho token (PUBACK latency is measured on it).
func (d *mqttDevice) publishDouble(iface, path string, value float64, ts time.Time) (paho.Token, error) {
	payload, err := bson.Marshal(bson.M{"v": value, "t": ts.UTC()})
	if err != nil {
		return nil, err
	}
	return d.client.Publish(d.base+"/"+iface+path, 1, false, payload), nil
}

// publishObject sends one object-aggregated document.
func (d *mqttDevice) publishObject(iface, path string, doc bson.M, ts time.Time) (paho.Token, error) {
	payload, err := bson.Marshal(bson.M{"v": doc, "t": ts.UTC()})
	if err != nil {
		return nil, err
	}
	return d.client.Publish(d.base+"/"+iface+path, 1, false, payload), nil
}
