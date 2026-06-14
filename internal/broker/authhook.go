package broker

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"log/slog"
	"net/netip"
	"sync"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	"github.com/astrate-platform/astrate/internal/pairing/ca"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// Store is the persistence surface the broker consumes (hexagonal-lite,
// docs/DESIGN.md §1.3). *store.Store satisfies it; tests substitute fakes.
type Store interface {
	// ListRealms feeds the per-realm client-CA pools.
	ListRealms(ctx context.Context) ([]store.Realm, error)
	// GetDevice authenticates connections and loads introspection for ACLs.
	GetDevice(ctx context.Context, realmID int16, id deviceid.ID) (*store.Device, error)
	// GetInterface resolves the ownership of introspected interfaces.
	GetInterface(ctx context.Context, realmID int16, name string, major int) (*store.StoredInterface, error)
	// SetDeviceConnected records a connection on the device row.
	SetDeviceConnected(ctx context.Context, realmID int16, id deviceid.ID, at time.Time, ip netip.Addr) error
	// SetDeviceDisconnected records a disconnection on the device row.
	SetDeviceDisconnected(ctx context.Context, realmID int16, id deviceid.ID, at time.Time) error
}

const (
	// hookDBTimeout bounds the database work performed inside broker hooks
	// (which receive no caller context).
	hookDBTimeout = 5 * time.Second

	// introspectionReloadDebounce rate-limits per-session introspection
	// reloads triggered by ACL misses on unknown interface names, so an
	// adversarial topic flood cannot hammer the database.
	introspectionReloadDebounce = time.Second
)

// realmCA is one realm's client-certificate trust anchor.
type realmCA struct {
	id   int16
	name string
	cert *x509.Certificate
	pool *x509.CertPool
}

// realmPools holds the per-realm CA pools and the TLS listener configuration
// snapshot (docs/DESIGN.md §3.1). Reload rebuilds both from the store; M7's
// realm CRUD calls Broker.ReloadRealms to hot-reload after creating or
// re-keying a realm.
type realmPools struct {
	st         Store
	serverCert tls.Certificate
	log        *slog.Logger

	mu     sync.RWMutex
	byName map[string]*realmCA
	tlsCfg *tls.Config
}

func newRealmPools(st Store, serverCert tls.Certificate, log *slog.Logger) *realmPools {
	p := &realmPools{st: st, serverCert: serverCert, log: log, byName: map[string]*realmCA{}}
	p.rebuildTLSLocked()
	return p
}

// Reload re-reads every realm's CA certificate and rebuilds the union client
// pool used for the TLS handshake. A realm whose stored CA fails to parse is
// skipped with a warning rather than failing the whole broker: one corrupt
// realm must not deny service to every other realm (docs/DESIGN.md §6
// single-process blast-radius). Devices in a skipped realm simply fail the
// handshake until the realm is re-keyed.
func (p *realmPools) Reload(ctx context.Context) error {
	realms, err := p.st.ListRealms(ctx)
	if err != nil {
		return err
	}
	byName := make(map[string]*realmCA, len(realms))
	for i := range realms {
		r := &realms[i]
		cert, err := ca.ParseCertificatePEM(r.CACertificatePEM)
		if err != nil {
			p.log.Warn("skipping realm with unparseable CA certificate", "realm", r.Name, "error", err)
			continue
		}
		pool := x509.NewCertPool()
		pool.AddCert(cert)
		byName[r.Name] = &realmCA{id: r.ID, name: r.Name, cert: cert, pool: pool}
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.byName = byName
	p.rebuildTLSLocked()
	return nil
}

// rebuildTLSLocked recomputes the handshake snapshot. Callers hold p.mu (or
// own p exclusively, during construction).
func (p *realmPools) rebuildTLSLocked() {
	union := x509.NewCertPool()
	for _, rc := range p.byName {
		union.AddCert(rc.cert)
	}
	p.tlsCfg = &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{p.serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    union,
	}
}

// Lookup returns the realm's trust anchor, if known.
func (p *realmPools) Lookup(name string) (*realmCA, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	rc, ok := p.byName[name]
	return rc, ok
}

// handshakeConfig returns the listener TLS configuration. The outer config
// defers to GetConfigForClient so Reload takes effect for new handshakes
// without rebinding the listener.
func (p *realmPools) handshakeConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		GetConfigForClient: func(*tls.ClientHelloInfo) (*tls.Config, error) {
			p.mu.RLock()
			defer p.mu.RUnlock()
			return p.tlsCfg, nil
		},
	}
}

// deviceSession is the broker-side state of one authenticated device
// connection: its identity plus the introspection-derived interface
// ownership map the ACL hook consults (docs/DESIGN.md §3.2).
type deviceSession struct {
	identity Identity
	realmID  int16
	client   *mqtt.Client // pointer identity guards against takeover races
	remote   netip.Addr

	mu            sync.Mutex
	ownership     map[string]interfaceschema.Ownership
	lastIntroLoad time.Time
}

// ownershipOf reports the ownership of an introspected interface name.
func (s *deviceSession) ownershipOf(iface string) (interfaceschema.Ownership, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.ownership[iface]
	return o, ok
}

// refresh reloads the device's introspection and ownership map from the
// store. The engine calls it (via Broker.RefreshIntrospection) after
// persisting a new introspection; the ACL hook calls it, debounced, when a
// topic names an interface the cached map does not know.
func (s *deviceSession) refresh(ctx context.Context, st Store, log *slog.Logger) error {
	dev, err := st.GetDevice(ctx, s.realmID, s.identity.DeviceID)
	if err != nil {
		return err
	}
	ownership := loadOwnership(ctx, st, s.realmID, dev.Introspection, log)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ownership = ownership
	s.lastIntroLoad = time.Now()
	return nil
}

// refreshIfStale runs refresh at most once per introspectionReloadDebounce.
func (s *deviceSession) refreshIfStale(ctx context.Context, st Store, log *slog.Logger) {
	s.mu.Lock()
	stale := time.Since(s.lastIntroLoad) >= introspectionReloadDebounce
	if stale {
		s.lastIntroLoad = time.Now() // claim the slot before the slow reload
	}
	s.mu.Unlock()
	if !stale {
		return
	}
	if err := s.refresh(ctx, st, log); err != nil {
		log.Warn("introspection refresh failed", "client", s.identity.CN(), "error", err)
	}
}

// loadOwnership resolves each introspected interface's ownership. Interfaces
// missing from the realm (not installed, or version mismatch) are skipped:
// the ACL then denies their topics, which is the correct posture.
func loadOwnership(ctx context.Context, st Store, realmID int16, intro map[string]store.InterfaceVersion, log *slog.Logger) map[string]interfaceschema.Ownership {
	ownership := make(map[string]interfaceschema.Ownership, len(intro))
	for name, ver := range intro {
		si, err := st.GetInterface(ctx, realmID, name, ver.Major)
		if err != nil {
			log.Debug("introspected interface not resolvable", "interface", name, "major", ver.Major, "error", err)
			continue
		}
		ownership[name] = si.Ownership
	}
	return ownership
}

// sessionRegistry maps live client IDs (= CNs) to their device sessions. It
// is shared by the auth, ACL, intake, and lifecycle hooks.
type sessionRegistry struct {
	mu       sync.RWMutex
	sessions map[string]*deviceSession
}

func newSessionRegistry() *sessionRegistry {
	return &sessionRegistry{sessions: map[string]*deviceSession{}}
}

func (r *sessionRegistry) get(clientID string) *deviceSession {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.sessions[clientID]
}

func (r *sessionRegistry) put(s *deviceSession) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[s.identity.CN()] = s
}

// removeIfOwner deletes and returns the session for clientID only when it
// belongs to cl: after a session takeover the registry entry points at the
// new connection, and the old connection's disconnect must not evict it.
func (r *sessionRegistry) removeIfOwner(clientID string, cl *mqtt.Client) *deviceSession {
	r.mu.Lock()
	defer r.mu.Unlock()
	s := r.sessions[clientID]
	if s == nil || s.client != cl {
		return nil
	}
	delete(r.sessions, clientID)
	return s
}

// authHook implements OnConnectAuthenticate (docs/ROADMAP.md §6 file 5.3):
// CN parse, chain verification against that realm's CA, device existence and
// inhibition checks, client-ID == CN, and latest-serial enforcement behind
// the pairing.enforce_latest_cert flag (docs/DESIGN.md §3.1, §4.3).
type authHook struct {
	mqtt.HookBase
	st                Store
	pools             *realmPools
	registry          *sessionRegistry
	enforceLatestCert bool
	devListenerID     string // empty when insecure_dev_mode is off
	log               *slog.Logger
}

// ID implements mqtt.Hook.
func (h *authHook) ID() string { return "astrate-auth" }

// Provides implements mqtt.Hook.
func (h *authHook) Provides(b byte) bool {
	return b == mqtt.OnConnectAuthenticate || b == mqtt.OnConnect
}

// OnConnect strips any Will message from the connection. Wills are not part
// of the Astarte MQTT v1 protocol, and mochi publishes them without a
// publish-side ACL check — accepting them would let a device plant a
// retained message on an arbitrary topic at disconnect time.
func (h *authHook) OnConnect(cl *mqtt.Client, _ packets.Packet) error {
	cl.Properties.Will = mqtt.Will{}
	return nil
}

// OnConnectAuthenticate implements the §3.1 connection contract.
func (h *authHook) OnConnectAuthenticate(cl *mqtt.Client, _ packets.Packet) bool {
	ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
	defer cancel()

	deny := func(reason string, err error) bool {
		h.log.Warn("mqtt connection rejected",
			"client", cl.ID, "remote", cl.Net.Remote, "listener", cl.Net.Listener,
			"reason", reason, "error", err)
		return false
	}

	identity, err := ParseCN(cl.ID)
	if err != nil {
		return deny("client ID is not a <realm>/<device_id> CN", err)
	}
	rc, ok := h.pools.Lookup(identity.Realm)
	if !ok {
		return deny("unknown realm", nil)
	}

	if tc, isTLS := cl.Net.Conn.(*tls.Conn); isTLS {
		peers := tc.ConnectionState().PeerCertificates
		if len(peers) == 0 {
			return deny("no client certificate presented", nil)
		}
		leaf := peers[0]
		if leaf.Subject.CommonName != cl.ID {
			return deny("client ID does not match certificate CN", nil)
		}
		if _, err := leaf.Verify(x509.VerifyOptions{
			Roots:         rc.pool,
			Intermediates: intermediatePool(peers[1:]),
			KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		}); err != nil {
			return deny("certificate does not chain to the realm CA", err)
		}
		dev, err := h.st.GetDevice(ctx, rc.id, identity.DeviceID)
		if err != nil {
			return deny("device lookup failed", err)
		}
		if dev.Status == store.DeviceStatusInhibited {
			return deny("device is inhibited", nil)
		}
		if h.enforceLatestCert && !isLatestCert(leaf, dev) {
			return deny("certificate superseded by a newer issuance", nil)
		}
		h.admit(ctx, cl, identity, rc.id, dev)
		return true
	}

	// Plaintext connection: only the insecure_dev_mode listener accepts
	// them, authenticating by claimed client ID alone (docs/DESIGN.md §3.1).
	if h.devListenerID == "" || cl.Net.Listener != h.devListenerID {
		return deny("plaintext connection outside insecure_dev_mode", nil)
	}
	dev, err := h.st.GetDevice(ctx, rc.id, identity.DeviceID)
	if err != nil {
		return deny("device lookup failed", err)
	}
	if dev.Status == store.DeviceStatusInhibited {
		return deny("device is inhibited", nil)
	}
	h.admit(ctx, cl, identity, rc.id, dev)
	return true
}

// admit registers the authenticated session for the ACL/intake/lifecycle
// hooks.
func (h *authHook) admit(ctx context.Context, cl *mqtt.Client, identity Identity, realmID int16, dev *store.Device) {
	sess := &deviceSession{
		identity:      identity,
		realmID:       realmID,
		client:        cl,
		remote:        remoteAddr(cl.Net.Remote),
		ownership:     loadOwnership(ctx, h.st, realmID, dev.Introspection, h.log),
		lastIntroLoad: time.Now(),
	}
	h.registry.put(sess)
}

// isLatestCert compares the presented leaf against the device row's recorded
// latest issuance: decimal serial plus lowercase-hex authority key ID
// (docs/DESIGN.md §4.3 always-online-CRL equivalent).
func isLatestCert(leaf *x509.Certificate, dev *store.Device) bool {
	if dev.CertSerial == nil || dev.CertAKI == nil {
		return false
	}
	return leaf.SerialNumber.String() == *dev.CertSerial &&
		hex.EncodeToString(leaf.AuthorityKeyId) == *dev.CertAKI
}

// intermediatePool builds a pool from any extra presented certificates
// (devices normally present a bare leaf).
func intermediatePool(certs []*x509.Certificate) *x509.CertPool {
	if len(certs) == 0 {
		return nil
	}
	pool := x509.NewCertPool()
	for _, c := range certs {
		pool.AddCert(c)
	}
	return pool
}

// remoteAddr parses mochi's "ip:port" remote string; unparseable remotes
// (e.g. in-memory test connections) degrade to the unspecified IPv4 address.
func remoteAddr(remote string) netip.Addr {
	if ap, err := netip.ParseAddrPort(remote); err == nil {
		return ap.Addr().Unmap()
	}
	return netip.IPv4Unspecified()
}
