// Package pairing implements Astarte's Pairing API surface
// (docs/DESIGN.md §4.4 flows A–C): device registration with show-once
// credentials secrets, CSR-based credential issuance through the embedded
// per-realm CA, broker discovery, and certificate verification — all
// wire-compatible with upstream so official device SDKs and astartectl run
// unmodified.
package pairing

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/astrate-platform/astrate/internal/pairing/ca"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// Store is the persistence surface the pairing service consumes
// (hexagonal-lite, docs/DESIGN.md §1.3). *store.Store satisfies it; tests
// use an in-memory fake.
type Store interface {
	GetRealmByName(ctx context.Context, name string) (*store.Realm, error)
	RegisterDevice(ctx context.Context, realmID int16, id deviceid.ID, secretHash string) error
	UnregisterDevice(ctx context.Context, realmID int16, id deviceid.ID) error
	GetDevice(ctx context.Context, realmID int16, id deviceid.ID) (*store.Device, error)
	SetDeviceCredentials(ctx context.Context, realmID int16, id deviceid.ID, certSerial, certAKI string, requestIP netip.Addr) error
	SetPayloadFormatHint(ctx context.Context, realmID int16, id deviceid.ID, hint string) error
	CountDevices(ctx context.Context, realmID int16) (int64, error)
}

// Service-level sentinel errors; the HTTP layer maps them onto upstream
// statuses and envelopes.
var (
	// ErrInvalidHWID reports a hw_id that is not a 22-character unpadded
	// base64url 128-bit device ID (422 upstream changeset shape).
	ErrInvalidHWID = errors.New("pairing: invalid hw_id")
	// ErrInvalidPayloadFormat reports an initial_payload_format outside
	// {bson, json} (Astrate extension, docs/DESIGN.md §3.5.4).
	ErrInvalidPayloadFormat = errors.New("pairing: invalid initial_payload_format")
	// ErrAlreadyRegistered reports re-registration of a device that has
	// already requested credentials (upstream error_name
	// "already_registered", 422).
	ErrAlreadyRegistered = errors.New("pairing: device already registered")
	// ErrRegistrationLimitReached reports the realm's
	// device_registration_limit being hit (upstream error_name
	// "device_registration_limit_reached", 422).
	ErrRegistrationLimitReached = errors.New("pairing: device registration limit reached")
	// ErrUnauthorized is the uniform device-authentication failure: unknown
	// device, unregistered device, or wrong credentials secret all produce
	// this same error (docs/DESIGN.md §4.4: no oracle distinguishing them).
	ErrUnauthorized = errors.New("pairing: unauthorized")
	// ErrInhibited reports a device blocked by credentials_inhibited (403).
	ErrInhibited = errors.New("pairing: credentials request inhibited")
	// ErrInvalidCSR reports an unusable certificate signing request (422).
	ErrInvalidCSR = errors.New("pairing: invalid CSR")
)

// Certificate verification causes on the wire (upstream
// CertificateValidationError enum subset, docs/DESIGN.md §4.4).
const (
	// CauseExpired marks a certificate outside its validity window.
	CauseExpired = "EXPIRED"
	// CauseInvalid marks a certificate that does not parse or does not
	// chain to the realm CA.
	CauseInvalid = "INVALID"
	// CauseRevoked marks a certificate superseded by a newer issuance
	// (latest-serial enforcement, docs/DESIGN.md §4.3).
	CauseRevoked = "REVOKED"
)

const (
	// secretBytes is the credentials secret entropy: 32 random bytes,
	// base64-encoded to the upstream-parity 44-character form.
	secretBytes = 32

	// DefaultVersion is reported by the info endpoint when Config.Version
	// is empty.
	DefaultVersion = "0.1.0-astrate"
)

// dummySecretHash is a bcrypt hash of an unguessable value, compared against
// when the device row (or its hash) is missing so authentication failures
// burn comparable time regardless of cause.
var dummySecretHash = func() string {
	h, err := bcrypt.GenerateFromPassword([]byte("astrate-dummy-secret-equalizer"), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Sprintf("pairing: generating dummy bcrypt hash: %v", err))
	}
	return string(h)
}()

// Config carries the service's operational knobs.
type Config struct {
	// BrokerURL is handed to devices by the info endpoint
	// (e.g. "mqtts://host:8883").
	BrokerURL string
	// CertTTL is the client certificate validity; zero selects
	// ca.DefaultCertTTL (30 days).
	CertTTL time.Duration
	// EnforceLatestCert enables the always-online-CRL behaviour
	// (docs/DESIGN.md §4.3): verify reports REVOKED for certificates whose
	// serial differs from the device's latest issuance.
	EnforceLatestCert bool
	// Version is reported by the info endpoint; empty selects
	// DefaultVersion.
	Version string
	// BcryptCost hashes credentials secrets; zero selects
	// bcrypt.DefaultCost (10, docs/DESIGN.md §4.1).
	BcryptCost int
}

// Service implements pairing flows A–C over a Store and the per-realm CA.
type Service struct {
	st     Store
	sealer *store.KeySealer
	cfg    Config
	now    func() time.Time
}

// New builds a pairing Service. The sealer opens realms' AES-GCM-sealed CA
// private keys (docs/DESIGN.md §4.3); it is owned by the caller and shared
// with housekeeping (which seals new realm CAs).
func New(st Store, sealer *store.KeySealer, cfg Config) *Service {
	if cfg.Version == "" {
		cfg.Version = DefaultVersion
	}
	if cfg.BcryptCost == 0 {
		cfg.BcryptCost = bcrypt.DefaultCost
	}
	if cfg.CertTTL <= 0 {
		cfg.CertTTL = ca.DefaultCertTTL
	}
	return &Service{st: st, sealer: sealer, cfg: cfg, now: time.Now}
}

// ProvisionCA generates a realm CA and seals its private key, producing the
// (ca_certificate, ca_private_key) pair stored on the realm row. Called by
// housekeeping at realm creation; lifetime zero selects the 10-year default.
func ProvisionCA(realmName string, lifetime time.Duration, sealer *store.KeySealer) (certPEM string, sealedKey []byte, err error) {
	realmCA, err := ca.Generate(realmName, lifetime)
	if err != nil {
		return "", nil, err
	}
	sealed, err := sealer.Seal(realmCA.PrivateKeyDER())
	if err != nil {
		return "", nil, fmt.Errorf("pairing: sealing CA key: %w", err)
	}
	return realmCA.CertificatePEM(), sealed, nil
}

// Register implements flow A: it validates hw_id, enforces the realm's
// device registration limit, generates a 44-character credentials secret,
// stores its bcrypt hash, and returns the secret (shown exactly once).
// Re-registering a device that has not yet requested credentials rotates the
// secret; afterwards it fails with ErrAlreadyRegistered. initialFormat is
// the Astrate initial_payload_format extension ("", "bson" or "json").
func (s *Service) Register(ctx context.Context, realmName, hwID, initialFormat string) (string, error) {
	id, err := deviceid.Parse(hwID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidHWID, err)
	}
	if initialFormat != "" && initialFormat != "bson" && initialFormat != "json" {
		return "", fmt.Errorf("%w: %q", ErrInvalidPayloadFormat, initialFormat)
	}

	realm, err := s.st.GetRealmByName(ctx, realmName)
	if err != nil {
		return "", err
	}

	// Upstream parity: the limit gates every registration attempt, secret
	// rotations for existing devices included.
	if realm.DeviceRegistrationLimit != nil {
		n, err := s.st.CountDevices(ctx, realm.ID)
		if err != nil {
			return "", err
		}
		if n >= int64(*realm.DeviceRegistrationLimit) {
			return "", ErrRegistrationLimitReached
		}
	}

	secret, err := generateSecret()
	if err != nil {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), s.cfg.BcryptCost)
	if err != nil {
		return "", fmt.Errorf("pairing: hashing credentials secret: %w", err)
	}

	if err := s.st.RegisterDevice(ctx, realm.ID, id, string(hash)); err != nil {
		if errors.Is(err, store.ErrDeviceAlreadyConfirmed) {
			return "", fmt.Errorf("%w: %s", ErrAlreadyRegistered, hwID)
		}
		return "", err
	}

	if initialFormat != "" {
		if err := s.st.SetPayloadFormatHint(ctx, realm.ID, id, initialFormat); err != nil {
			return "", err
		}
	}
	return secret, nil
}

// Unregister implements the flow A DELETE: the device becomes registrable
// again, its data is kept (store.UnregisterDevice clears only the
// credential trail). store.ErrNotFound is returned for unknown devices.
func (s *Service) Unregister(ctx context.Context, realmName, deviceIDStr string) error {
	id, err := deviceid.Parse(deviceIDStr)
	if err != nil {
		return fmt.Errorf("%w: device %q", store.ErrNotFound, deviceIDStr)
	}
	realm, err := s.st.GetRealmByName(ctx, realmName)
	if err != nil {
		return err
	}
	return s.st.UnregisterDevice(ctx, realm.ID, id)
}

// Credentials implements flow B, the SDK hot path: it authenticates the
// device by credentials secret (uniform ErrUnauthorized on any mismatch),
// rejects inhibited devices, signs the CSR with the realm CA, and records
// the new certificate's serial/AKI plus the requesting IP (which also flips
// the device status registered → confirmed).
func (s *Service) Credentials(ctx context.Context, realmName, deviceIDStr, secret, csrPEM string, ip netip.Addr) (string, error) {
	realm, dev, err := s.authenticateDevice(ctx, realmName, deviceIDStr, secret)
	if err != nil {
		return "", err
	}
	if dev.Status == store.DeviceStatusInhibited {
		return "", ErrInhibited
	}

	realmCA, err := s.loadCA(realm)
	if err != nil {
		return "", err
	}
	certPEM, serial, aki, err := realmCA.SignCSR(csrPEM, realmName, dev.ID.String(), s.cfg.CertTTL)
	if err != nil {
		if errors.Is(err, ca.ErrInvalidCSR) {
			return "", fmt.Errorf("%w: %v", ErrInvalidCSR, err)
		}
		return "", err
	}
	if err := s.st.SetDeviceCredentials(ctx, realm.ID, dev.ID, serial, aki, ip); err != nil {
		return "", err
	}
	return certPEM, nil
}

// Info is the flow C device-info document.
type Info struct {
	Version   string
	Status    string
	BrokerURL string
	CACertPEM string
}

// Info implements the flow C GET: broker discovery plus device status,
// authenticated by credentials secret. Status strings are upstream parity:
// "pending" until the first credentials request, then "confirmed";
// "inhibited" when blocked (inhibited devices may still read info — only new
// credentials and connections are blocked).
func (s *Service) Info(ctx context.Context, realmName, deviceIDStr, secret string) (*Info, error) {
	realm, dev, err := s.authenticateDevice(ctx, realmName, deviceIDStr, secret)
	if err != nil {
		return nil, err
	}

	var status string
	switch dev.Status {
	case store.DeviceStatusInhibited:
		status = "inhibited"
	case store.DeviceStatusConfirmed:
		status = "confirmed"
	default:
		status = "pending"
	}
	return &Info{
		Version:   s.cfg.Version,
		Status:    status,
		BrokerURL: s.cfg.BrokerURL,
		CACertPEM: realm.CACertificatePEM,
	}, nil
}

// VerifyResult is the flow C credentials/verify outcome. With Valid true,
// Until carries the certificate expiry; otherwise Cause carries one of the
// Cause* constants. Timestamp is the verification instant.
type VerifyResult struct {
	Valid     bool
	Timestamp time.Time
	Until     time.Time
	Cause     string
}

// VerifyCredentials implements the flow C verify endpoint, authenticated by
// credentials secret. Classification precedence: certificates outside their
// validity window report EXPIRED; certificates that fail to parse or to
// chain to the realm CA report INVALID; valid-but-superseded certificates
// report REVOKED when EnforceLatestCert is on.
func (s *Service) VerifyCredentials(ctx context.Context, realmName, deviceIDStr, secret, clientCrtPEM string) (*VerifyResult, error) {
	realm, dev, err := s.authenticateDevice(ctx, realmName, deviceIDStr, secret)
	if err != nil {
		return nil, err
	}
	if dev.Status == store.DeviceStatusInhibited {
		return nil, ErrInhibited
	}

	res := &VerifyResult{Timestamp: s.now()}

	realmCA, err := s.loadCA(realm)
	if err != nil {
		return nil, err
	}
	until, err := realmCA.Verify(clientCrtPEM, res.Timestamp)
	switch {
	case errors.Is(err, ca.ErrCertificateExpired):
		res.Cause = CauseExpired
		return res, nil
	case err != nil:
		res.Cause = CauseInvalid
		return res, nil
	}

	if s.cfg.EnforceLatestCert && dev.CertSerial != nil {
		cert, perr := ca.ParseCertificatePEM(clientCrtPEM)
		if perr != nil {
			res.Cause = CauseInvalid
			return res, nil
		}
		if cert.SerialNumber.String() != *dev.CertSerial {
			res.Cause = CauseRevoked
			return res, nil
		}
	}

	res.Valid = true
	res.Until = until
	return res, nil
}

// authenticateDevice resolves the realm and device and bcrypt-compares the
// presented secret. All failure modes (unknown realm or device, unregistered
// device, wrong secret) return the same ErrUnauthorized after a comparable
// amount of bcrypt work, so response timing does not leak device existence.
func (s *Service) authenticateDevice(ctx context.Context, realmName, deviceIDStr, secret string) (*store.Realm, *store.Device, error) {
	id, err := deviceid.Parse(deviceIDStr)
	if err != nil {
		burnBcrypt(secret)
		return nil, nil, ErrUnauthorized
	}
	realm, err := s.st.GetRealmByName(ctx, realmName)
	if errors.Is(err, store.ErrNotFound) {
		burnBcrypt(secret)
		return nil, nil, ErrUnauthorized
	}
	if err != nil {
		return nil, nil, err
	}
	dev, err := s.st.GetDevice(ctx, realm.ID, id)
	if errors.Is(err, store.ErrNotFound) {
		burnBcrypt(secret)
		return nil, nil, ErrUnauthorized
	}
	if err != nil {
		return nil, nil, err
	}
	if dev.CredentialsSecretHash == "" {
		// Unregistered (or never-registered) device: no valid hash exists.
		burnBcrypt(secret)
		return nil, nil, ErrUnauthorized
	}
	if bcrypt.CompareHashAndPassword([]byte(dev.CredentialsSecretHash), []byte(secret)) != nil {
		return nil, nil, ErrUnauthorized
	}
	return realm, dev, nil
}

// loadCA opens the realm's sealed CA key and reconstructs its CA.
func (s *Service) loadCA(realm *store.Realm) (*ca.CA, error) {
	keyDER, err := s.sealer.Open(realm.CAPrivateKeySealed)
	if err != nil {
		return nil, fmt.Errorf("pairing: opening sealed CA key of realm %q: %w", realm.Name, err)
	}
	realmCA, err := ca.Load(realm.CACertificatePEM, keyDER)
	if err != nil {
		return nil, fmt.Errorf("pairing: loading CA of realm %q: %w", realm.Name, err)
	}
	return realmCA, nil
}

// burnBcrypt spends one bcrypt comparison against the dummy hash, equalizing
// the work done on authentication paths that lack a real hash.
func burnBcrypt(secret string) {
	_ = bcrypt.CompareHashAndPassword([]byte(dummySecretHash), []byte(secret))
}

// generateSecret draws the 32-byte credentials secret and encodes it as the
// upstream-parity 44-character standard base64 string.
func generateSecret() (string, error) {
	buf := make([]byte, secretBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("pairing: gathering secret randomness: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}
