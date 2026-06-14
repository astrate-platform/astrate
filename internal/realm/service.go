// Package realm is the Realm Management API (docs/DESIGN.md §3.7, ROADMAP
// §8.1 files 7.1–7.2): the operator-facing surface for installing and
// versioning interfaces, managing triggers, and rotating a realm's JWT
// auth key. It is wire-shaped to upstream astarte_realm_management so
// astartectl and the dashboard work unmodified, and every mutation both
// emits the store NOTIFY and calls the in-process engine invalidation
// callback so changes take effect immediately.
package realm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sort"

	"github.com/astrate-platform/astrate/internal/engine/triggers"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// ErrValidation wraps a request that is well-formed JSON but violates an
// interface/trigger schema rule (maps to a 422). Use errors.Is against it
// and fmt.Errorf("%w: ...", ErrValidation, ...) to attach the detail.
var ErrValidation = errors.New("realm: validation failed")

// Invalidator is the in-process cache-invalidation callback the engine
// satisfies (*engine.Engine's RefreshInterfaces / RefreshTriggers). After a
// realm mutation the service refreshes the engine's compiled snapshot so the
// change takes effect without waiting for the LISTEN/NOTIFY round-trip. A nil
// Invalidator disables the in-process path; the store NOTIFY still fires.
type Invalidator interface {
	RefreshInterfaces(ctx context.Context, realmID int16) error
	RefreshTriggers(ctx context.Context, realmID int16) error
}

// Service implements the Realm Management business logic over the store.
type Service struct {
	st  *store.Store
	inv Invalidator
	log *slog.Logger
}

// NewService builds the service. inv may be nil (e.g. management-only
// deployments without a local engine); log defaults to slog.Default().
func NewService(st *store.Store, inv Invalidator, log *slog.Logger) *Service {
	if log == nil {
		log = slog.Default()
	}
	return &Service{st: st, inv: inv, log: log}
}

// realmID resolves a realm name to its id; an unknown realm surfaces
// store.ErrNotFound (the HTTP layer maps it to 404, though the auth
// middleware normally rejects unknown realms first).
func (s *Service) realmID(ctx context.Context, realm string) (int16, error) {
	r, err := s.st.GetRealmByName(ctx, realm)
	if err != nil {
		return 0, err
	}
	return r.ID, nil
}

// --- interfaces -------------------------------------------------------------

// InstallInterface validates and installs a new interface major
// (docs/ROADMAP.md §8.1). A duplicate (name, major) yields
// store.ErrAlreadyExists; a schema violation yields ErrValidation.
func (s *Service) InstallInterface(ctx context.Context, realm string, def []byte) (*store.StoredInterface, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	if _, err := interfaceschema.ParseInterface(def); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	si, err := s.st.InstallInterface(ctx, rid, def)
	if err != nil {
		return nil, err
	}
	s.interfacesChanged(ctx, rid, realm)
	return si, nil
}

// UpdateInterface applies a minor upgrade, enforcing the additive-only
// upstream parity rules via interfaceschema.CheckMinorUpgrade (no mutated
// mapping attributes, same type/ownership/aggregation, strictly higher
// minor). The interface major must already exist.
func (s *Service) UpdateInterface(ctx context.Context, realm string, def []byte) (*store.StoredInterface, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	next, err := interfaceschema.ParseInterface(def)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	stored, err := s.st.GetInterface(ctx, rid, next.Name, next.Major)
	if err != nil {
		return nil, err
	}
	prev, err := interfaceschema.ParseInterface(stored.Definition)
	if err != nil {
		return nil, fmt.Errorf("realm: stored interface %s v%d does not parse: %w", next.Name, next.Major, err)
	}
	if err := interfaceschema.CheckMinorUpgrade(prev, next); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	si, err := s.st.UpdateInterface(ctx, rid, def)
	if err != nil {
		return nil, err
	}
	s.interfacesChanged(ctx, rid, realm)
	return si, nil
}

// DeleteInterface removes an interface major. The store enforces the upstream
// draining rules (store.ErrInterfaceMajorNotZero, store.ErrInterfaceInUse).
func (s *Service) DeleteInterface(ctx context.Context, realm, name string, major int) error {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return err
	}
	if err := s.st.DeleteInterface(ctx, rid, name, major); err != nil {
		return err
	}
	s.interfacesChanged(ctx, rid, realm)
	return nil
}

// ListInterfaces returns the distinct interface names installed in the realm
// (upstream GET /interfaces), sorted for stable output.
func (s *Service) ListInterfaces(ctx context.Context, realm string) ([]string, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	ifaces, err := s.st.LoadRealmInterfaces(ctx, rid)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{}, len(ifaces))
	names := make([]string, 0, len(ifaces))
	for _, si := range ifaces {
		if _, ok := seen[si.Name]; ok {
			continue
		}
		seen[si.Name] = struct{}{}
		names = append(names, si.Name)
	}
	sort.Strings(names)
	return names, nil
}

// ListInterfaceMajors returns the installed major versions of one interface
// name (upstream GET /interfaces/{name}), ascending. An unknown name yields
// store.ErrNotFound.
func (s *Service) ListInterfaceMajors(ctx context.Context, realm, name string) ([]int, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	ifaces, err := s.st.LoadRealmInterfaces(ctx, rid)
	if err != nil {
		return nil, err
	}
	var majors []int
	for _, si := range ifaces {
		if si.Name == name {
			majors = append(majors, si.Major)
		}
	}
	if len(majors) == 0 {
		return nil, fmt.Errorf("%w: interface %s", store.ErrNotFound, name)
	}
	sort.Ints(majors)
	return majors, nil
}

// GetInterface returns the stored definition JSON of one interface major
// (upstream GET /interfaces/{name}/{major}).
func (s *Service) GetInterface(ctx context.Context, realm, name string, major int) (json.RawMessage, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	si, err := s.st.GetInterface(ctx, rid, name, major)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(si.Definition), nil
}

// --- triggers ---------------------------------------------------------------

// triggerName is the minimal projection used to read a trigger's name out of
// its definition body.
type triggerName struct {
	Name string `json:"name"`
}

// CreateTrigger validates a trigger definition (name + action + simple
// triggers, via triggers.Compile — the same validation the engine applies)
// and installs it. A duplicate name yields store.ErrAlreadyExists.
func (s *Service) CreateTrigger(ctx context.Context, realm string, def []byte) (*store.Trigger, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	var tn triggerName
	if err := json.Unmarshal(def, &tn); err != nil {
		return nil, fmt.Errorf("%w: trigger does not parse: %v", ErrValidation, err)
	}
	if tn.Name == "" {
		return nil, fmt.Errorf("%w: trigger requires a name", ErrValidation)
	}
	if _, err := triggers.Compile(tn.Name, def); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	tr, err := s.st.CreateTrigger(ctx, rid, tn.Name, def)
	if err != nil {
		return nil, err
	}
	s.triggersChanged(ctx, rid, realm)
	return tr, nil
}

// GetTrigger returns one trigger's definition JSON.
func (s *Service) GetTrigger(ctx context.Context, realm, name string) (json.RawMessage, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	tr, err := s.st.GetTrigger(ctx, rid, name)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(tr.Definition), nil
}

// DeleteTrigger removes one trigger by name.
func (s *Service) DeleteTrigger(ctx context.Context, realm, name string) error {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return err
	}
	if err := s.st.DeleteTrigger(ctx, rid, name); err != nil {
		return err
	}
	s.triggersChanged(ctx, rid, realm)
	return nil
}

// ListTriggers returns the realm's trigger names, sorted.
func (s *Service) ListTriggers(ctx context.Context, realm string) ([]string, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	trs, err := s.st.ListTriggers(ctx, rid)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(trs))
	for i := range trs {
		names = append(names, trs[i].Name)
	}
	sort.Strings(names)
	return names, nil
}

// --- config/auth ------------------------------------------------------------

// GetAuthKey returns the realm's JWT public key PEM (upstream GET
// /config/auth → {"jwt_public_key_pem": "..."}). Astrate stores a list for
// rotation; the wire field carries them concatenated, which the verifier
// already splits into individual keys.
func (s *Service) GetAuthKey(ctx context.Context, realm string) (string, error) {
	r, err := s.st.GetRealmByName(ctx, realm)
	if err != nil {
		return "", err
	}
	return joinPEM(r.JWTPublicKeysPEM), nil
}

// SetAuthKey rotates the realm's JWT public key (upstream PUT /config/auth).
// The supplied PEM may concatenate multiple keys for a rotation window.
func (s *Service) SetAuthKey(ctx context.Context, realm, keyPEM string) error {
	if keyPEM == "" {
		return fmt.Errorf("%w: jwt_public_key_pem can't be blank", ErrValidation)
	}
	if err := s.st.SetRealmJWTPublicKeys(ctx, realm, []string{keyPEM}); err != nil {
		return err
	}
	return nil
}

// joinPEM concatenates PEM blocks with a blank line, the form the verifier
// parses back into individual keys.
func joinPEM(keys []string) string {
	switch len(keys) {
	case 0:
		return ""
	case 1:
		return keys[0]
	default:
		out := keys[0]
		for _, k := range keys[1:] {
			out += "\n" + k
		}
		return out
	}
}

// --- invalidation -----------------------------------------------------------

// interfacesChanged emits the store NOTIFY and refreshes the engine snapshot
// after an interface mutation. Failures are logged, never fatal: the change
// is already committed and the engine self-heals on its next reload.
func (s *Service) interfacesChanged(ctx context.Context, rid int16, realm string) {
	if err := s.st.NotifyInterfacesChanged(ctx, rid); err != nil {
		s.log.Warn("realm: NOTIFY after interface change failed", "realm", realm, "err", err)
	}
	if s.inv != nil {
		if err := s.inv.RefreshInterfaces(ctx, rid); err != nil {
			s.log.Warn("realm: engine interface refresh failed", "realm", realm, "err", err)
		}
	}
}

// triggersChanged refreshes the engine snapshot after a trigger mutation
// (triggers ride in the same realm snapshot as interfaces).
func (s *Service) triggersChanged(ctx context.Context, rid int16, realm string) {
	if err := s.st.NotifyInterfacesChanged(ctx, rid); err != nil {
		s.log.Warn("realm: NOTIFY after trigger change failed", "realm", realm, "err", err)
	}
	if s.inv != nil {
		if err := s.inv.RefreshTriggers(ctx, rid); err != nil {
			s.log.Warn("realm: engine trigger refresh failed", "realm", realm, "err", err)
		}
	}
}
