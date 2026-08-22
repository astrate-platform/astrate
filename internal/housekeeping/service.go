// Package housekeeping is the instance-admin Housekeeping API
// (docs/DESIGN.md §3.7, ROADMAP §8.1 files 7.3–7.4): realm provisioning and
// teardown, guarded by instance-level JWT keys carrying a_ha. Creating a
// realm mints its embedded CA (docs/DESIGN.md §4.3) and seals the CA private
// key before it ever touches the database.
package housekeeping

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"

	"github.com/astrate-platform/astrate/internal/pairing/ca"
	"github.com/astrate-platform/astrate/internal/store"
)

// ErrValidation wraps a well-formed request that violates a realm-creation
// rule (maps to 422).
var ErrValidation = errors.New("housekeeping: validation failed")

// RealmView is the API projection of a realm: never the CA private key.
type RealmView struct {
	Name                              string
	JWTPublicKeyPEM                   string
	DeviceRegistrationLimit           *int32
	DatastreamMaximumStorageRetention *int64
}

// Reloader is the broker port a realm mutation notifies so a freshly-created
// (or torn-down) realm's CA pool is trusted for new TLS handshakes without a
// restart (docs/DESIGN.md §3.1). *broker.Broker satisfies it via ReloadRealms;
// it may be nil (read-only or broker-less deployments).
type Reloader interface {
	ReloadRealms(ctx context.Context) error
}

// Service implements the Housekeeping business logic over the store, holding
// the key sealer used to protect freshly-minted CA private keys.
type Service struct {
	st       *store.Store
	sealer   *store.KeySealer
	reloader Reloader
	log      *slog.Logger

	// defaultRetention is the realm-level datastream retention ceiling
	// injected at creation when the caller omits the field (issue #73).
	// nil disables injection: an explicit value always wins.
	defaultRetention *int64
}

// NewService builds the service. reloader (the broker) may be nil; log
// defaults to slog.Default().
func NewService(st *store.Store, sealer *store.KeySealer, reloader Reloader, log *slog.Logger) *Service {
	if log == nil {
		log = slog.Default()
	}
	return &Service{st: st, sealer: sealer, reloader: reloader, log: log}
}

// WithDefaultDatastreamMaximumStorageRetention sets the realm-level datastream
// retention ceiling (seconds) injected at creation when the caller omits the
// field (#73). nil (the zero behavior) disables injection: an explicit value
// always wins, and deployments that never configure it behave exactly as
// before.
func (s *Service) WithDefaultDatastreamMaximumStorageRetention(defaultRetention *int64) *Service {
	s.defaultRetention = defaultRetention
	return s
}

// notifyBrokerReload asks the broker to rebuild its per-realm CA pools. A
// failure is logged but never fails the mutation: the realm row is already
// committed, and a later reload (or the broker's own self-heal) recovers.
func (s *Service) notifyBrokerReload(ctx context.Context, realm string) {
	if s.reloader == nil {
		return
	}
	if err := s.reloader.ReloadRealms(ctx); err != nil {
		s.log.Warn("broker realm reload failed after housekeeping mutation", "realm", realm, "error", err)
	}
}

// CreateRealm provisions a realm: mint a self-signed realm CA (ECDSA P-256,
// default 10-year lifetime), seal its private key, and persist the realm row
// plus CA material in one store transaction (docs/ROADMAP.md §8.1 file 7.3).
// A blank/invalid name or missing JWT key yields ErrValidation; a duplicate
// realm yields store.ErrAlreadyExists. retention is the realm-level
// datastream storage ceiling in seconds; when the caller omits it and a
// default was configured on the service, the default is injected (#73) — an
// explicit value always wins.
func (s *Service) CreateRealm(ctx context.Context, name, jwtPublicKeyPEM string, regLimit *int32, retention *int64) (*RealmView, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: realm_name can't be blank", ErrValidation)
	}
	if jwtPublicKeyPEM == "" {
		return nil, fmt.Errorf("%w: jwt_public_key_pem can't be blank", ErrValidation)
	}
	if regLimit != nil && *regLimit < 0 {
		return nil, fmt.Errorf("%w: device_registration_limit must be non-negative", ErrValidation)
	}
	if retention != nil && *retention < 0 {
		return nil, fmt.Errorf("%w: datastream_maximum_storage_retention must be non-negative", ErrValidation)
	}
	if retention == nil {
		retention = s.defaultRetention
	}

	realmCA, err := ca.Generate(name, 0)
	if err != nil {
		return nil, fmt.Errorf("housekeeping: minting realm CA: %w", err)
	}
	sealed, err := s.sealer.Seal(realmCA.PrivateKeyDER())
	if err != nil {
		return nil, fmt.Errorf("housekeeping: sealing realm CA key: %w", err)
	}

	r, err := s.st.CreateRealm(ctx, store.NewRealm{
		Name:                              name,
		JWTPublicKeysPEM:                  []string{jwtPublicKeyPEM},
		CACertificatePEM:                  realmCA.CertificatePEM(),
		CAPrivateKeySealed:                sealed,
		DeviceRegistrationLimit:           regLimit,
		DatastreamMaximumStorageRetention: retention,
	})
	if err != nil {
		if errors.Is(err, store.ErrInvalidRealmName) {
			return nil, fmt.Errorf("%w: realm_name is invalid", ErrValidation)
		}
		return nil, err
	}
	s.notifyBrokerReload(ctx, name)
	return view(r), nil
}

// RealmUpdate carries a PATCH payload; each Patch* flag gates its value
// field, Clear* flags unset it (Clear beats Set).
type RealmUpdate struct {
	PatchJWTPublicKeyPEM   bool
	SetJWTPublicKeyPEM     string
	PatchRegistrationLimit bool
	SetRegistrationLimit   int32
	ClearRegistrationLimit bool
	PatchRetention         bool
	SetRetention           int64
	ClearRetention         bool
}

// UpdateRealm applies an update to one realm after validating the touched
// fields (upstream PATCH /housekeeping/v1/realms/{realm}); an unknown realm
// yields store.ErrNotFound.
func (s *Service) UpdateRealm(ctx context.Context, name string, p RealmUpdate) error {
	if p.PatchJWTPublicKeyPEM && p.SetJWTPublicKeyPEM == "" {
		return fmt.Errorf("%w: jwt_public_key_pem can't be blank", ErrValidation)
	}
	if p.PatchRegistrationLimit && p.SetRegistrationLimit < 0 {
		return fmt.Errorf("%w: device_registration_limit must be non-negative", ErrValidation)
	}
	if p.PatchRetention && p.SetRetention < 0 {
		return fmt.Errorf("%w: datastream_maximum_storage_retention must be non-negative", ErrValidation)
	}
	return s.st.UpdateRealm(ctx, name, store.RealmPatch{
		PatchJWTPublicKeyPEM:   p.PatchJWTPublicKeyPEM,
		SetJWTPublicKeyPEM:     p.SetJWTPublicKeyPEM,
		PatchRegistrationLimit: p.PatchRegistrationLimit,
		SetRegistrationLimit:   p.SetRegistrationLimit,
		ClearRegistrationLimit: p.ClearRegistrationLimit,
		PatchRetention:         p.PatchRetention,
		SetRetention:           p.SetRetention,
		ClearRetention:         p.ClearRetention,
	})
}

// GetRealm returns one realm's public view (upstream GET
// /housekeeping/v1/realms/{realm}).
func (s *Service) GetRealm(ctx context.Context, name string) (*RealmView, error) {
	r, err := s.st.GetRealmByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return view(r), nil
}

// ListRealms returns the realm names, sorted (upstream GET
// /housekeeping/v1/realms).
func (s *Service) ListRealms(ctx context.Context) ([]string, error) {
	realms, err := s.st.ListRealms(ctx)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(realms))
	for i := range realms {
		names = append(names, realms[i].Name)
	}
	sort.Strings(names)
	return names, nil
}

// DeleteRealm tears a realm down, cascading its interfaces, devices,
// properties, and datastream rows (store.DeleteRealm, docs/DESIGN.md §2.1).
func (s *Service) DeleteRealm(ctx context.Context, name string) error {
	if err := s.st.DeleteRealm(ctx, name); err != nil {
		return err
	}
	s.notifyBrokerReload(ctx, name)
	return nil
}

// view projects a stored realm into its API shape, dropping CA material.
func view(r *store.Realm) *RealmView {
	key := ""
	if len(r.JWTPublicKeysPEM) > 0 {
		key = r.JWTPublicKeysPEM[0]
	}
	return &RealmView{
		Name:                              r.Name,
		JWTPublicKeyPEM:                   key,
		DeviceRegistrationLimit:           r.DeviceRegistrationLimit,
		DatastreamMaximumStorageRetention: r.DatastreamMaximumStorageRetention,
	}
}
