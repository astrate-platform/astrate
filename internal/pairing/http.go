package pairing

import (
	"errors"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

// maxBodyBytes caps pairing request bodies (docs/DESIGN.md §4.5 input
// bounds): CSRs and certificates fit comfortably in 64 KiB.
const maxBodyBytes int64 = 64 << 10

// Default rate-limit parameters (docs/DESIGN.md §4.5). Registration is
// operator-driven (JWT-protected, fleet provisioning bursts); credentials
// requests are device-driven and rare (a renewal per device per cert TTL).
const (
	DefaultRegisterRate     = 5.0
	DefaultRegisterBurst    = 20
	DefaultCredentialsRate  = 1.0
	DefaultCredentialsBurst = 5
)

// detailTooManyRequests is the 429 envelope detail.
const detailTooManyRequests = "Too Many Requests"

// APIConfig tunes the HTTP layer's rate limits; zero values select the
// defaults above.
type APIConfig struct {
	RegisterRate     float64
	RegisterBurst    int
	CredentialsRate  float64
	CredentialsBurst int
}

// API is the /pairing/v1 HTTP surface (docs/DESIGN.md §4.4, §3.7). Agent
// endpoints are guarded by realm JWTs carrying a_pa; device endpoints
// authenticate with the device's credentials secret as a bearer token.
type API struct {
	svc          *Service
	requireAgent func(http.Handler) http.Handler
	regLimiter   *Limiter
	credLimiter  *Limiter
}

// NewAPI wires the pairing service to its HTTP surface. mw provides the
// realm-JWT middleware (M3).
func NewAPI(svc *Service, mw *auth.Middleware, cfg APIConfig) *API {
	if cfg.RegisterRate == 0 {
		cfg.RegisterRate = DefaultRegisterRate
	}
	if cfg.RegisterBurst == 0 {
		cfg.RegisterBurst = DefaultRegisterBurst
	}
	if cfg.CredentialsRate == 0 {
		cfg.CredentialsRate = DefaultCredentialsRate
	}
	if cfg.CredentialsBurst == 0 {
		cfg.CredentialsBurst = DefaultCredentialsBurst
	}
	return &API{
		svc:          svc,
		requireAgent: mw.RequireRealm(auth.ClaimPairing),
		regLimiter:   NewLimiter(cfg.RegisterRate, cfg.RegisterBurst),
		credLimiter:  NewLimiter(cfg.CredentialsRate, cfg.CredentialsBurst),
	}
}

// Mount registers the pairing routes on mux. Paths are wire-frozen
// (docs/DESIGN.md §4.4): they are exactly what the official SDKs and
// astartectl call.
func (a *API) Mount(mux *http.ServeMux) {
	mux.Handle("POST /pairing/v1/{realm}/agent/devices",
		a.requireAgent(http.HandlerFunc(a.handleRegister)))
	mux.Handle("DELETE /pairing/v1/{realm}/agent/devices/{deviceID}",
		a.requireAgent(http.HandlerFunc(a.handleUnregister)))
	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
		a.handleCredentials)
	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
		a.handleInfo)
	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
		a.handleVerify)
}

// --- flow A: agent endpoints ------------------------------------------------

// registerRequest is the flow A body: hw_id plus the Astrate
// initial_payload_format extension (docs/DESIGN.md §3.5.4; upstream-shaped
// requests simply omit it).
type registerRequest struct {
	HwID                 string `json:"hw_id"`
	InitialPayloadFormat string `json:"initial_payload_format"`
}

// registerResponse is the show-once secret envelope payload.
type registerResponse struct {
	CredentialsSecret string `json:"credentials_secret"`
}

func (a *API) handleRegister(w http.ResponseWriter, r *http.Request) {
	if !a.regLimiter.Allow("ip|" + remoteIP(r).String()) {
		_ = astarteapi.WriteError(w, http.StatusTooManyRequests, detailTooManyRequests)
		return
	}

	var req registerRequest
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &req); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if req.HwID == "" {
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"hw_id": {"can't be blank"}})
		return
	}

	secret, err := a.svc.Register(r.Context(), r.PathValue("realm"), req.HwID, req.InitialPayloadFormat)
	if err != nil {
		a.writeServiceError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusCreated, registerResponse{CredentialsSecret: secret})
}

func (a *API) handleUnregister(w http.ResponseWriter, r *http.Request) {
	err := a.svc.Unregister(r.Context(), r.PathValue("realm"), r.PathValue("deviceID"))
	if err != nil {
		a.writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- flow B/C: device endpoints ---------------------------------------------

// credentialsRequest is the flow B body.
type credentialsRequest struct {
	CSR string `json:"csr"`
}

// credentialsResponse carries the issued client certificate.
type credentialsResponse struct {
	ClientCrt string `json:"client_crt"`
}

func (a *API) handleCredentials(w http.ResponseWriter, r *http.Request) {
	realm, deviceID := r.PathValue("realm"), r.PathValue("deviceID")
	ip := remoteIP(r)
	if !a.credLimiter.Allow("ip|"+ip.String()) || !a.credLimiter.Allow("dev|"+realm+"/"+deviceID) {
		_ = astarteapi.WriteError(w, http.StatusTooManyRequests, detailTooManyRequests)
		return
	}

	secret, ok := bearerSecret(r)
	if !ok {
		_ = astarteapi.WriteUnauthorized(w)
		return
	}
	var req credentialsRequest
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &req); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if req.CSR == "" {
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"csr": {"can't be blank"}})
		return
	}

	clientCrt, err := a.svc.Credentials(r.Context(), realm, deviceID, secret, req.CSR, ip)
	if err != nil {
		a.writeServiceError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusCreated, credentialsResponse{ClientCrt: clientCrt})
}

// infoResponse is the flow C device-info payload. Field order matches the
// upstream-rendered JSON; ca_crt inside the protocol entry is an Astrate
// extension (docs/DESIGN.md §4.4) — a pure superset upstream clients ignore.
type infoResponse struct {
	Protocols struct {
		AstarteMQTTV1 struct {
			BrokerURL string `json:"broker_url"`
			CACrt     string `json:"ca_crt"`
		} `json:"astarte_mqtt_v1"`
	} `json:"protocols"`
	Status  string `json:"status"`
	Version string `json:"version"`
}

func (a *API) handleInfo(w http.ResponseWriter, r *http.Request) {
	secret, ok := bearerSecret(r)
	if !ok {
		_ = astarteapi.WriteUnauthorized(w)
		return
	}
	info, err := a.svc.Info(r.Context(), r.PathValue("realm"), r.PathValue("deviceID"), secret)
	if err != nil {
		a.writeServiceError(w, err)
		return
	}
	var resp infoResponse
	resp.Protocols.AstarteMQTTV1.BrokerURL = info.BrokerURL
	resp.Protocols.AstarteMQTTV1.CACrt = info.CACertPEM
	resp.Status = info.Status
	resp.Version = info.Version
	_ = astarteapi.WriteData(w, http.StatusOK, resp)
}

// verifyRequest is the flow C verify body.
type verifyRequest struct {
	ClientCrt string `json:"client_crt"`
}

// verifyValidResponse / verifyInvalidResponse mirror upstream's
// CredentialsStatusView shapes: timestamps are rendered as Elixir
// DateTime.to_string of a millisecond-precision UTC datetime, and the
// invalid shape carries an always-null details field.
type verifyValidResponse struct {
	Timestamp string `json:"timestamp"`
	Until     string `json:"until"`
	Valid     bool   `json:"valid"`
}

type verifyInvalidResponse struct {
	Cause     string  `json:"cause"`
	Details   *string `json:"details"`
	Timestamp string  `json:"timestamp"`
	Valid     bool    `json:"valid"`
}

func (a *API) handleVerify(w http.ResponseWriter, r *http.Request) {
	secret, ok := bearerSecret(r)
	if !ok {
		_ = astarteapi.WriteUnauthorized(w)
		return
	}
	var req verifyRequest
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &req); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if req.ClientCrt == "" {
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"client_crt": {"can't be blank"}})
		return
	}

	res, err := a.svc.VerifyCredentials(r.Context(), r.PathValue("realm"), r.PathValue("deviceID"), secret, req.ClientCrt)
	if err != nil {
		a.writeServiceError(w, err)
		return
	}
	if res.Valid {
		_ = astarteapi.WriteData(w, http.StatusOK, verifyValidResponse{
			Timestamp: formatDateTime(res.Timestamp),
			Until:     formatDateTime(res.Until),
			Valid:     true,
		})
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, verifyInvalidResponse{
		Cause:     res.Cause,
		Timestamp: formatDateTime(res.Timestamp),
		Valid:     false,
	})
}

// --- shared plumbing --------------------------------------------------------

// writeServiceError maps service errors onto upstream statuses and bodies.
//
// Two shapes are upstream-verbatim quirks worth naming: 422 validation
// failures use the Phoenix changeset envelope ({"errors": {"<field>":
// ["<msg>"]}}), and RPC-originated conflicts surface as the literal
// "error_name" field. One mapping deviates deliberately: upstream answers
// wrong-secret/unknown-device with 403; Astrate uses 401 (frozen in
// docs/ROADMAP.md §5 — semantically authentication, and uniform across
// causes). Inhibited devices stay 403.
func (a *API) writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidHWID):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"hw_id": {"is not a valid base64 encoded 128 bits id"}})
	case errors.Is(err, ErrInvalidPayloadFormat):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"initial_payload_format": {"is invalid"}})
	case errors.Is(err, ErrAlreadyRegistered):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"error_name": {"already_registered"}})
	case errors.Is(err, ErrRegistrationLimitReached):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"error_name": {"device_registration_limit_reached"}})
	case errors.Is(err, ErrInvalidCSR):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"csr": {"is invalid"}})
	case errors.Is(err, ErrUnauthorized):
		_ = astarteapi.WriteUnauthorized(w)
	case errors.Is(err, ErrInhibited):
		_ = astarteapi.WriteForbidden(w)
	case errors.Is(err, store.ErrNotFound):
		_ = astarteapi.WriteDeviceNotFound(w)
	default:
		_ = astarteapi.WriteInternalServerError(w)
	}
}

// formatDateTime renders t exactly as upstream's views do
// (Elixir DateTime.to_string on a millisecond-precision UTC datetime):
// "2024-05-30 13:49:57.045Z".
func formatDateTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05.000Z")
}

// bearerSecret extracts the device credentials secret from the
// Authorization header with upstream's permissive matching
// (~r/bearer\:?\s+(.*)$/i).
func bearerSecret(r *http.Request) (string, bool) {
	header := r.Header.Get("Authorization")
	scheme, rest, found := strings.Cut(header, " ")
	if !found {
		return "", false
	}
	if !strings.EqualFold(strings.TrimSuffix(scheme, ":"), "bearer") {
		return "", false
	}
	secret := strings.TrimSpace(rest)
	return secret, secret != ""
}

// remoteIP extracts the peer address from the request, falling back to the
// IPv4 unspecified address when RemoteAddr is unparseable (e.g. unix
// sockets behind a proxy).
func remoteIP(r *http.Request) netip.Addr {
	if ap, err := netip.ParseAddrPort(r.RemoteAddr); err == nil {
		return ap.Addr().Unmap()
	}
	if a, err := netip.ParseAddr(r.RemoteAddr); err == nil {
		return a.Unmap()
	}
	return netip.IPv4Unspecified()
}
