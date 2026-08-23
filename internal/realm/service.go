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
	"strings"
	"time"

	"github.com/astrate-platform/astrate/internal/engine/triggers"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// ErrValidation wraps a request that is well-formed JSON but violates an
// interface/trigger schema rule (maps to a 422). Use errors.Is against it
// and fmt.Errorf("%w: ...", ErrValidation, ...) to attach the detail.
var ErrValidation = errors.New("realm: validation failed")

// ErrMaximumDatabaseRetentionExceeded rejects an interface install/update
// whose mapping TTL exceeds the realm's datastream_maximum_storage_retention
// ceiling (upstream error_name maximum_database_retention_exceeded, #72).
var ErrMaximumDatabaseRetentionExceeded = errors.New("realm: maximum_database_retention_exceeded")

var (
	// ErrNameMismatch is a PUT whose body interface_name disagrees with the
	// URL (upstream name_not_matching, #62).
	ErrNameMismatch = errors.New("realm: interface name does not match")
	// ErrMajorMismatch is a PUT whose body version_major disagrees with the
	// URL (upstream major_version_not_matching, #62).
	ErrMajorMismatch = errors.New("realm: interface major does not match")
	// ErrNameCollision rejects an install whose name equals an installed one
	// modulo case and hyphens (upstream interface_name_collision, #62).
	ErrNameCollision = errors.New("realm: interface name collision")
)

// errMajorNotFound marks a PUT/DELETE lookup miss on the (name, major) pair:
// upstream answers those 404 "Interface major not found" rather than
// "Interface not found" (#62). The HTTP layer matches it first; a bare
// store.ErrNotFound still renders "Interface not found".
var errMajorNotFound = errors.New("realm: interface major not found")

// Invalidator is the in-process cache-invalidation callback the engine
// satisfies (*engine.Engine's RefreshInterfaces / RefreshTriggers). After a
// realm mutation the service refreshes the engine's compiled snapshot so the
// change takes effect without waiting for the LISTEN/NOTIFY round-trip. A nil
// Invalidator disables the in-process path; the store NOTIFY still fires.
type Invalidator interface {
	RefreshInterfaces(ctx context.Context, realmID int16) error
	RefreshTriggers(ctx context.Context, realmID int16) error
}

// Disconnecter force-closes a device's live MQTT session (*broker.Broker
// satisfies it). Device deletion kicks the session before wiping the rows,
// mirroring upstream's disconnect-first deletion order. A nil Disconnecter
// skips the kick; the deleted device is then refused at its next reconnect
// or credential request instead.
type Disconnecter interface {
	DisconnectDevice(realm string, id deviceid.ID)
}

// OnDeletionFunc is called around a synchronous device delete so the engine
// can emit device_deletion_started / device_deletion_finished (issue #21).
// The callback receives the realm name, encoded device ID, and the instant.
type OnDeletionFunc func(realmName, deviceID string, at time.Time)

// Service implements the Realm Management business logic over the store.
type Service struct {
	st   *store.Store
	inv  Invalidator
	disc Disconnecter
	// OnDeletionStart / OnDeletionFinish bookend DeleteDevice (nil-safe).
	OnDeletionStart  OnDeletionFunc
	OnDeletionFinish OnDeletionFunc
	log              *slog.Logger
}

// NewService builds the service. inv may be nil (e.g. management-only
// deployments without a local engine); log defaults to slog.Default().
func NewService(st *store.Store, inv Invalidator, log *slog.Logger) *Service {
	if log == nil {
		log = slog.Default()
	}
	return &Service{st: st, inv: inv, log: log}
}

// WithDisconnecter attaches the broker seam used to kick a live session at
// device deletion (nil-safe, mirrors housekeeping's Reloader wiring).
func (s *Service) WithDisconnecter(d Disconnecter) *Service {
	s.disc = d
	return s
}

// DeleteDevice synchronously removes a device and all its data (upstream
// starts an async deletion with a transient deletion_in_progress state;
// Astrate is single-process and deletes in one transaction —
// docs/COMPATIBILITY.md). Emits device_deletion_started immediately before
// the store delete and device_deletion_finished immediately after (issue
// #21: back-to-back around the sync path so imported trigger configs still
// fire). finished is emitted even when the store delete fails, so a started
// lifecycle always closes.
func (s *Service) DeleteDevice(ctx context.Context, realm, deviceID string) error {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return err
	}
	id, err := deviceid.Parse(deviceID)
	if err != nil {
		return fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
	}
	if s.disc != nil {
		s.disc.DisconnectDevice(realm, id)
	}
	at := time.Now().UTC()
	if s.OnDeletionStart != nil {
		s.OnDeletionStart(realm, deviceID, at)
	}
	err = s.st.DeleteDevice(ctx, rid, id)
	if s.OnDeletionFinish != nil {
		s.OnDeletionFinish(realm, deviceID, time.Now().UTC())
	}
	return err
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
// store.ErrAlreadyExists; a schema violation yields ErrValidation — carrying
// a *interfaceschema.ViolationsError when the rejection has a probe-verified
// upstream wire shape (#61). Documents posted with legacy alias fields
// (quality/aggregate/path) are stored canonically so GET renders them the
// way upstream does.
func (s *Service) InstallInterface(ctx context.Context, realm string, def []byte) (*store.StoredInterface, error) {
	r, err := s.st.GetRealmByName(ctx, realm)
	if err != nil {
		return nil, err
	}
	iface, canon, err := interfaceschema.ParseInterfaceCanonical(def)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
	}
	if canon != nil {
		def = canon
	}
	// Upstream rejects names that differ from an installed one only by case
	// or hyphens (#62); identical raw names are skipped here so a duplicate
	// install keeps its store.ErrAlreadyExists shape.
	if err := s.checkNameCollision(ctx, r.ID, iface.Name); err != nil {
		return nil, err
	}
	if err := s.checkRetentionCeiling(r, iface); err != nil {
		return nil, err
	}
	si, err := s.st.InstallInterface(ctx, r.ID, def)
	if err != nil {
		return nil, err
	}
	s.interfacesChanged(ctx, r.ID, realm)
	return si, nil
}

// UpdateInterface applies a minor upgrade, enforcing the additive-only
// upstream parity rules via interfaceschema.CheckMinorUpgrade (no mutated
// mapping attributes, same type/ownership/aggregation, strictly higher
// minor). urlName/urlMajor are the interface identity from the URL path: the
// parsed body must agree with both (upstream 409s the mismatches, #62), and
// a lookup miss on that identity yields the "major not found" marker.
func (s *Service) UpdateInterface(ctx context.Context, realm, urlName string, urlMajor int, def []byte) (*store.StoredInterface, error) {
	r, err := s.st.GetRealmByName(ctx, realm)
	if err != nil {
		return nil, err
	}
	next, canon, err := interfaceschema.ParseInterfaceCanonical(def)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
	}
	if canon != nil {
		def = canon
	}
	if next.Name != urlName {
		return nil, fmt.Errorf("%w: body declares %q, URL says %q", ErrNameMismatch, next.Name, urlName)
	}
	if next.Major != urlMajor {
		return nil, fmt.Errorf("%w: body declares %d, URL says %d", ErrMajorMismatch, next.Major, urlMajor)
	}
	if err := s.checkRetentionCeiling(r, next); err != nil {
		return nil, err
	}
	stored, err := s.st.GetInterface(ctx, r.ID, next.Name, next.Major)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, fmt.Errorf("%w: %w", errMajorNotFound, err)
		}
		return nil, err
	}
	prev, err := interfaceschema.ParseInterface(stored.Definition)
	if err != nil {
		return nil, fmt.Errorf("realm: stored interface %s v%d does not parse: %w", next.Name, next.Major, err)
	}
	if err := interfaceschema.CheckMinorUpgrade(prev, next); err != nil {
		// Wrapped so the classification sentinels survive to the HTTP layer.
		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
	}
	si, err := s.st.UpdateInterface(ctx, r.ID, def)
	if err != nil {
		return nil, err
	}
	s.interfacesChanged(ctx, r.ID, realm)
	return si, nil
}

// checkRetentionCeiling rejects iface when any mapping's database TTL exceeds
// the realm's datastream_maximum_storage_retention ceiling (#72). Realms
// without a ceiling pass untouched.
func (s *Service) checkRetentionCeiling(r *store.Realm, iface *interfaceschema.Interface) error {
	if r.DatastreamMaximumStorageRetention == nil {
		return nil
	}
	for i := range iface.Mappings {
		m := &iface.Mappings[i]
		if m.DatabaseRetentionPolicy == interfaceschema.UseTTL &&
			m.DatabaseRetentionTTL > *r.DatastreamMaximumStorageRetention {
			return ErrMaximumDatabaseRetentionExceeded
		}
	}
	return nil
}

// DeleteInterface removes an interface major. The store enforces the upstream
// draining rules (store.ErrInterfaceMajorNotZero, store.ErrInterfaceInUse);
// a lookup miss is marked errMajorNotFound for the HTTP layer (#62).
func (s *Service) DeleteInterface(ctx context.Context, realm, name string, major int) error {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return err
	}
	if err := s.st.DeleteInterface(ctx, rid, name, major); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("%w: %w", errMajorNotFound, err)
		}
		return err
	}
	s.interfacesChanged(ctx, rid, realm)
	return nil
}

// checkNameCollision rejects a new interface whose name equals an installed
// one modulo case and hyphens (upstream interface_name_collision, #62).
// Identical raw names are skipped so the duplicate install keeps its
// store.ErrAlreadyExists path.
func (s *Service) checkNameCollision(ctx context.Context, rid int16, name string) error {
	ifaces, err := s.st.LoadRealmInterfaces(ctx, rid)
	if err != nil {
		return err
	}
	norm := normaliseIfaceName(name)
	for _, si := range ifaces {
		if si.Name == name {
			continue
		}
		if normaliseIfaceName(si.Name) == norm {
			return fmt.Errorf("%w: installed interface %q differs from %q only by case or hyphens",
				ErrNameCollision, si.Name, name)
		}
	}
	return nil
}

// normaliseIfaceName lowercases name and drops '-', upstream's collision key
// for interface names.
func normaliseIfaceName(name string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r == '-':
			return -1
		case 'A' <= r && r <= 'Z':
			return r + ('a' - 'A')
		default:
			return r
		}
	}, name)
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

// ListInterfacesDetailed renders the additive 1.4-style detailed listing
// (issue #66): one fully materialised document per installed interface
// major, sorted by (name, major). The names-only listing upstream 1.2
// serves stays on ListInterfaces.
func (s *Service) ListInterfacesDetailed(ctx context.Context, realm string) ([]json.RawMessage, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	ifaces, err := s.st.LoadRealmInterfaces(ctx, rid)
	if err != nil {
		return nil, err
	}
	sort.Slice(ifaces, func(i, j int) bool {
		if ifaces[i].Name != ifaces[j].Name {
			return ifaces[i].Name < ifaces[j].Name
		}
		return ifaces[i].Major < ifaces[j].Major
	})
	docs := make([]json.RawMessage, 0, len(ifaces))
	for _, si := range ifaces {
		doc, err := renderDetailedInterface(si.Definition)
		if err != nil {
			return nil, fmt.Errorf("realm: stored interface %s v%d does not parse: %w", si.Name, si.Major, err)
		}
		docs = append(docs, doc)
	}
	return docs, nil
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

// triggerPolicyRef is the minimal projection used to read the policy a trigger
// references out of its stored definition.
type triggerPolicyRef struct {
	Policy string `json:"policy"`
}

// CreateTrigger validates a trigger definition (name + action + simple
// triggers, via triggers.Compile — the same validation the engine applies)
// and installs it. A duplicate name yields store.ErrAlreadyExists.
// If the trigger references a policy, that policy must already exist in the
// realm; otherwise the request is rejected as a validation error.
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
	ct, err := triggers.Compile(tn.Name, def)
	if err != nil {
		// Trigger validation failures carry the upstream-shaped field errors
		// across action and simple_triggers (issues #63/#70); the HTTP layer
		// renders them as the nested 422 envelope. Everything else stays a
		// generic ErrValidation.
		var te *triggers.TriggerErrors
		if errors.As(err, &te) {
			return nil, te
		}
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	if ct.PolicyName != "" {
		if _, err := s.st.GetTriggerPolicy(ctx, rid, ct.PolicyName); err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return nil, fmt.Errorf("%w: policy %q does not exist in realm %q", ErrValidation, ct.PolicyName, realm)
			}
			return nil, err
		}
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

// --- trigger delivery policies ----------------------------------------------
//
// The engine compiles policies into its realm snapshot and attaches them to
// triggers, yet these mutations deliberately do not invalidate it: the
// referential rules below mean a policy is unreferenced both when it is
// created and when it is deleted, and every trigger mutation — the only way a
// policy enters or leaves a snapshot — already refreshes it. Relax either
// rule, or add an update operation, and this stops holding: call
// triggersChanged here at that point.

// CreatePolicy validates and stores a delivery policy (upstream POST
// /policies). Duplicates yield store.ErrAlreadyExists.
func (s *Service) CreatePolicy(ctx context.Context, realm string, def []byte) (*store.TriggerPolicy, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	name, err := validatePolicy(def)
	if err != nil {
		return nil, err
	}
	return s.st.CreateTriggerPolicy(ctx, rid, name, def)
}

// GetPolicy returns one policy's definition.
func (s *Service) GetPolicy(ctx context.Context, realm, name string) (json.RawMessage, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	p, err := s.st.GetTriggerPolicy(ctx, rid, name)
	if err != nil {
		return nil, err
	}
	return p.Definition, nil
}

// ListPolicies returns the realm's policy names.
func (s *Service) ListPolicies(ctx context.Context, realm string) ([]string, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	ps, err := s.st.ListTriggerPolicies(ctx, rid)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(ps))
	for i := range ps {
		names[i] = ps[i].Name
	}
	return names, nil
}

// DeletePolicy removes a policy. The request is rejected if any trigger in
// the realm still references this policy.
func (s *Service) DeletePolicy(ctx context.Context, realm, name string) error {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return err
	}
	trs, err := s.st.ListTriggers(ctx, rid)
	if err != nil {
		return err
	}
	for i := range trs {
		var ref triggerPolicyRef
		if err := json.Unmarshal(trs[i].Definition, &ref); err != nil {
			continue
		}
		if ref.Policy == name {
			return fmt.Errorf("%w: policy %q is still used by trigger %q", ErrValidation, name, trs[i].Name)
		}
	}
	return s.st.DeleteTriggerPolicy(ctx, rid, name)
}

// APICompatVersion is the upstream Realm Management API level Astrate
// emulates, reported by GET /v1/{realm}/version. The Astarte Dashboard
// feature-gates UI sections on it (the Trigger Delivery Policies page
// requires >= 1.1.1), so this is a compatibility declaration, not Astrate's
// own release version (docs/COMPATIBILITY.md). Bump rule: only in the same
// change that completes the full surface of the new level, after reconciling
// every row in docs/UPSTREAM-EXPERIMENTAL.md tagged with it — never via
// configuration, never speculatively.
const APICompatVersion = "1.2.2"

// GetDeviceRegistrationLimit returns the realm's device registration limit
// (nil = unlimited), upstream GET /config/device_registration_limit.
func (s *Service) GetDeviceRegistrationLimit(ctx context.Context, realm string) (*int32, error) {
	r, err := s.st.GetRealmByName(ctx, realm)
	if err != nil {
		return nil, err
	}
	return r.DeviceRegistrationLimit, nil
}

// GetDatastreamMaximumStorageRetention returns the realm's datastream maximum
// storage retention in seconds, upstream GET
// /config/datastream_maximum_storage_retention (served since 1.2.0). An unset
// ceiling renders as 0 — the #60 wire contract (upstream answers null here,
// a recorded deviation), pinned by tests.
func (s *Service) GetDatastreamMaximumStorageRetention(ctx context.Context, realm string) (int64, error) {
	r, err := s.st.GetRealmByName(ctx, realm)
	if err != nil {
		return 0, err
	}
	if r.DatastreamMaximumStorageRetention == nil {
		return 0, nil
	}
	return *r.DatastreamMaximumStorageRetention, nil
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
