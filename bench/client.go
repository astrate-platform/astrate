package main

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Endpoints are the per-service API bases, each including the service prefix
// but not the version segment (e.g. http://host/pairing). The default
// derivation from one base URL matches both Astarte's standard gateway layout
// and Astrate's single mux, which is also what astartectl assumes.
type Endpoints struct {
	Housekeeping string `json:"housekeeping"`
	RealmMgmt    string `json:"realm_management"`
	Pairing      string `json:"pairing"`
	AppEngine    string `json:"appengine"`
}

func deriveEndpoints(base string) Endpoints {
	b := strings.TrimRight(base, "/")
	return Endpoints{
		Housekeeping: b + "/housekeeping",
		RealmMgmt:    b + "/realmmanagement",
		Pairing:      b + "/pairing",
		AppEngine:    b + "/appengine",
	}
}

// client is a thin Astarte REST client: JSON bodies in the {"data": ...}
// envelope, bearer auth, and backoff on 429 (Astrate rate-limits the pairing
// API; provisioning bursts hit it).
type client struct {
	http *http.Client
}

func newClient(timeout time.Duration) *client {
	return &client{http: &http.Client{Timeout: timeout}}
}

type apiError struct {
	Status int
	Body   string
}

func (e *apiError) Error() string {
	body := e.Body
	if len(body) > 300 {
		body = body[:300] + "…"
	}
	return fmt.Sprintf("HTTP %d: %s", e.Status, body)
}

// doJSON performs one API call. A nil in body sends no payload; a non-nil out
// receives the decoded "data" member of the response envelope.
func (c *client) doJSON(method, rawurl, bearer string, in, out any) error {
	var lastErr error
	for attempt := 0; attempt < 6; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(1<<attempt) * 100 * time.Millisecond)
		}

		var body io.Reader
		if in != nil {
			raw, err := json.Marshal(map[string]any{"data": in})
			if err != nil {
				return err
			}
			body = bytes.NewReader(raw)
		}
		req, err := http.NewRequest(method, rawurl, body)
		if err != nil {
			return err
		}
		if in != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		if bearer != "" {
			req.Header.Set("Authorization", "Bearer "+bearer)
		}

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = err
			continue // transient transport error: retry
		}
		raw, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = &apiError{Status: resp.StatusCode, Body: string(raw)}
			continue // rate-limited: back off and retry
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return &apiError{Status: resp.StatusCode, Body: string(raw)}
		}
		if out == nil {
			return nil
		}
		var envelope struct {
			Data json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(raw, &envelope); err != nil {
			return fmt.Errorf("decoding response envelope: %w (body %q)", err, raw)
		}
		return json.Unmarshal(envelope.Data, out)
	}
	return fmt.Errorf("giving up after retries: %w", lastErr)
}

// --- JWT minting -----------------------------------------------------------

// parsePrivateKeyPEM accepts RSA (PKCS1/PKCS8) and EC (SEC1/PKCS8) private
// keys — Astarte tooling has generated both across versions.
func parsePrivateKeyPEM(pemBytes []byte) (crypto.Signer, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found")
	}
	if k, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return k, nil
	}
	if k, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		return k, nil
	}
	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %w", err)
	}
	signer, ok := k.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("unsupported private key type %T", k)
	}
	return signer, nil
}

// mintJWT signs an Astarte API token carrying the given authorization claims,
// each granted the match-all "<verb>::<path>" string.
func mintJWT(key crypto.Signer, claims ...string) (string, error) {
	now := time.Now()
	mc := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(2 * time.Hour).Unix(),
	}
	for _, c := range claims {
		mc[c] = []string{".*::.*"}
	}
	var method jwt.SigningMethod
	switch key.(type) {
	case *rsa.PrivateKey:
		method = jwt.SigningMethodRS256
	case *ecdsa.PrivateKey:
		method = jwt.SigningMethodES256
	default:
		return "", fmt.Errorf("unsupported signing key type %T", key)
	}
	return jwt.NewWithClaims(method, mc).SignedString(key)
}

// --- API operations --------------------------------------------------------

func (c *client) createRealm(ep Endpoints, hkKey crypto.Signer, realm, realmPubPEM string) error {
	token, err := mintJWT(hkKey, "a_ha")
	if err != nil {
		return err
	}
	err = c.doJSON(http.MethodPost, ep.Housekeeping+"/v1/realms", token, map[string]any{
		"realm_name":         realm,
		"jwt_public_key_pem": realmPubPEM,
	}, nil)
	if isStatus(err, http.StatusConflict) {
		return nil // realm already exists: provision is re-runnable
	}
	return err
}

func (c *client) installInterface(ep Endpoints, realmKey crypto.Signer, realm string, ifaceJSON []byte) error {
	token, err := mintJWT(realmKey, "a_rma")
	if err != nil {
		return err
	}
	var body any
	if err := json.Unmarshal(ifaceJSON, &body); err != nil {
		return err
	}
	err = c.doJSON(http.MethodPost, ep.RealmMgmt+"/v1/"+realm+"/interfaces", token, body, nil)
	if isStatus(err, http.StatusConflict) {
		return nil // already installed
	}
	return err
}

func (c *client) registerDevice(ep Endpoints, realmKey crypto.Signer, realm, hwID string) (string, error) {
	token, err := mintJWT(realmKey, "a_pa")
	if err != nil {
		return "", err
	}
	var out struct {
		CredentialsSecret string `json:"credentials_secret"`
	}
	err = c.doJSON(http.MethodPost, ep.Pairing+"/v1/"+realm+"/agent/devices", token,
		map[string]any{"hw_id": hwID}, &out)
	if err != nil {
		return "", err
	}
	return out.CredentialsSecret, nil
}

func (c *client) obtainCertificate(ep Endpoints, realm string, dev Device, csrPEM string) (string, error) {
	var out struct {
		ClientCrt string `json:"client_crt"`
	}
	err := c.doJSON(http.MethodPost,
		ep.Pairing+"/v1/"+realm+"/devices/"+dev.ID+"/protocols/astarte_mqtt_v1/credentials",
		dev.Secret, map[string]any{"csr": csrPEM}, &out)
	if err != nil {
		return "", err
	}
	return out.ClientCrt, nil
}

func (c *client) brokerURL(ep Endpoints, realm string, dev Device) (string, error) {
	var out struct {
		Protocols struct {
			AstarteMQTTV1 struct {
				BrokerURL string `json:"broker_url"`
			} `json:"astarte_mqtt_v1"`
		} `json:"protocols"`
	}
	err := c.doJSON(http.MethodGet, ep.Pairing+"/v1/"+realm+"/devices/"+dev.ID, dev.Secret, nil, &out)
	if err != nil {
		return "", err
	}
	if out.Protocols.AstarteMQTTV1.BrokerURL == "" {
		return "", fmt.Errorf("pairing info returned no astarte_mqtt_v1 broker_url")
	}
	return out.Protocols.AstarteMQTTV1.BrokerURL, nil
}

// sample is one datastream row as the AppEngine API returns it. Value stays
// raw because object-aggregate rows carry a document there, not a scalar.
type sample struct {
	Value     json.RawMessage `json:"value"`
	Timestamp string          `json:"timestamp"`
}

// float64Value parses a scalar double row (the individual-datastream case).
func (s sample) float64Value() (float64, error) {
	var v float64
	err := json.Unmarshal(s.Value, &v)
	return v, err
}

// getSamples fetches rows of a concrete interface path. Extra query
// parameters (since, to, since_after, limit, sort) come via q. Both platforms
// return either a bare array or a single object for the data member; the
// caller only ever needs arrays here, so single objects are wrapped.
func (c *client) getSamples(ep Endpoints, token, realm, device, iface, path string, q url.Values) ([]sample, error) {
	u := ep.AppEngine + "/v1/" + realm + "/devices/" + device + "/interfaces/" + iface + path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}
	var raw json.RawMessage
	if err := c.doJSON(http.MethodGet, u, token, nil, &raw); err != nil {
		return nil, err
	}
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) > 0 && trimmed[0] == '{' {
		var one sample
		if err := json.Unmarshal(trimmed, &one); err != nil {
			return nil, err
		}
		return []sample{one}, nil
	}
	var many []sample
	if err := json.Unmarshal(trimmed, &many); err != nil {
		return nil, err
	}
	return many, nil
}

// getRaw performs a GET and discards the payload; used by the query mix where
// only latency matters.
func (c *client) getRaw(u, token string) error {
	var raw json.RawMessage
	return c.doJSON(http.MethodGet, u, token, nil, &raw)
}

// isStatus reports whether err carries the given HTTP status.
func isStatus(err error, status int) bool {
	var ae *apiError
	return errors.As(err, &ae) && ae.Status == status
}
