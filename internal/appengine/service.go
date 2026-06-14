// Package appengine is the AppEngine API (docs/DESIGN.md §3.7, ROADMAP §8.2):
// the operator/application surface for device status, interface data queries,
// server-owned publishing, groups, and the live event socket. It is wire-shaped
// to upstream astarte_appengine_api so astartectl and applications work
// unmodified, reading through the store and writing server-owned values through
// the engine.
package appengine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// DefaultDeviceLimit is the device-list page size when the caller gives none.
const DefaultDeviceLimit = 100

// ErrValidation wraps a well-formed request that violates a rule (maps to 422).
var ErrValidation = errors.New("appengine: validation failed")

// ServerData is the engine port for server-owned writes (docs/ROADMAP.md §8.2
// file 7.7). *engine.Engine satisfies it; tests substitute a fake.
type ServerData interface {
	PublishServerValue(ctx context.Context, realm string, id deviceid.ID, iface, path string, value json.RawMessage, ts *time.Time) error
	UnsetServerProperty(ctx context.Context, realm string, id deviceid.ID, iface, path string) error
}

// Service implements the AppEngine business logic over the store and engine.
type Service struct {
	st  *store.Store
	sd  ServerData
	log *slog.Logger
}

// NewService builds the service. sd may be nil (read-only deployments); log
// defaults to slog.Default().
func NewService(st *store.Store, sd ServerData, log *slog.Logger) *Service {
	if log == nil {
		log = slog.Default()
	}
	return &Service{st: st, sd: sd, log: log}
}

func (s *Service) realmID(ctx context.Context, realm string) (int16, error) {
	r, err := s.st.GetRealmByName(ctx, realm)
	if err != nil {
		return 0, err
	}
	return r.ID, nil
}

// --- devices ----------------------------------------------------------------

// DeviceStatus is the AppEngine device status body (upstream
// DeviceStatusView): a JSON-friendly projection of a devices row.
type DeviceStatus struct {
	ID                   string                        `json:"id"`
	Connected            bool                          `json:"connected"`
	Introspection        map[string]introspectionEntry `json:"introspection"`
	Aliases              map[string]string             `json:"aliases"`
	Attributes           map[string]string             `json:"attributes"`
	Groups               []string                      `json:"groups"`
	CredentialsInhibited bool                          `json:"credentials_inhibited"`
	TotalReceivedMsgs    int64                         `json:"total_received_msgs"`
	TotalReceivedBytes   int64                         `json:"total_received_bytes"`
	FirstRegistration    *time.Time                    `json:"first_registration"`
	FirstCredentialsReq  *time.Time                    `json:"first_credentials_request"`
	LastConnection       *time.Time                    `json:"last_connection"`
	LastDisconnection    *time.Time                    `json:"last_disconnection"`
	LastSeenIP           string                        `json:"last_seen_ip,omitempty"`
	PreviousInterfaces   map[string]introspectionEntry `json:"previous_interfaces,omitempty"`
}

// introspectionEntry renders one introspection pair (upstream shape).
type introspectionEntry struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

// ListDevices returns one page of device IDs (upstream GET /devices). after is
// the cursor (exclusive); limit <= 0 selects DefaultDeviceLimit. The second
// result is the next-page cursor, or "" at the end.
func (s *Service) ListDevices(ctx context.Context, realm string, after string, limit int) ([]string, string, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, "", err
	}
	if limit <= 0 {
		limit = DefaultDeviceLimit
	}
	var cursor *deviceid.ID
	if after != "" {
		id, err := deviceid.Parse(after)
		if err != nil {
			return nil, "", fmt.Errorf("%w: invalid cursor", ErrValidation)
		}
		cursor = &id
	}
	devs, err := s.st.ListDevices(ctx, rid, cursor, limit+1)
	if err != nil {
		return nil, "", err
	}
	next := ""
	if len(devs) > limit {
		next = devs[limit-1].ID.String()
		devs = devs[:limit]
	}
	ids := make([]string, len(devs))
	for i := range devs {
		ids[i] = devs[i].ID.String()
	}
	return ids, next, nil
}

// GetDevice returns one device's status (upstream GET /devices/{id}).
func (s *Service) GetDevice(ctx context.Context, realm, deviceID string) (*DeviceStatus, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	id, err := deviceid.Parse(deviceID)
	if err != nil {
		return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
	}
	d, err := s.st.GetDevice(ctx, rid, id)
	if err != nil {
		return nil, err
	}
	return s.deviceStatus(ctx, rid, d)
}

// GetDeviceByAlias returns the status of the device owning an alias (upstream
// GET /devices-by-alias/{alias}).
func (s *Service) GetDeviceByAlias(ctx context.Context, realm, alias string) (*DeviceStatus, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	d, err := s.st.GetDeviceByAlias(ctx, rid, alias)
	if err != nil {
		return nil, err
	}
	return s.deviceStatus(ctx, rid, d)
}

// deviceStatus projects a stored device into its status body, resolving its
// group memberships.
func (s *Service) deviceStatus(ctx context.Context, rid int16, d *store.Device) (*DeviceStatus, error) {
	groups, err := s.st.ListDeviceGroups(ctx, rid, d.ID)
	if err != nil {
		return nil, err
	}
	ds := &DeviceStatus{
		ID:                   d.ID.String(),
		Connected:            d.Connected,
		Introspection:        introspectionView(d.Introspection),
		Aliases:              orEmptyStr(d.Aliases),
		Attributes:           orEmptyStr(d.Attributes),
		Groups:               groups,
		CredentialsInhibited: d.Status == store.DeviceStatusInhibited,
		TotalReceivedMsgs:    d.TotalReceivedMsgs,
		TotalReceivedBytes:   d.TotalReceivedBytes,
		LastConnection:       d.LastConnection,
		LastDisconnection:    d.LastDisconnection,
		FirstCredentialsReq:  d.FirstCredentialsRequest,
		PreviousInterfaces:   introspectionView(d.OldIntrospection),
	}
	reg := d.FirstRegistration
	ds.FirstRegistration = &reg
	if d.LastSeenIP != nil {
		ds.LastSeenIP = d.LastSeenIP.String()
	}
	return ds, nil
}

// DevicePatch carries the JSON-merge-style PATCH fields (upstream
// DevicePatch): a nil pointer leaves a field unchanged, a present map patches
// aliases/attributes (a nil map value removes that key).
type DevicePatch struct {
	Aliases              map[string]*string
	Attributes           map[string]*string
	CredentialsInhibited *bool
}

// PatchDevice applies a device patch (upstream PATCH /devices/{id}).
func (s *Service) PatchDevice(ctx context.Context, realm, deviceID string, p DevicePatch) (*DeviceStatus, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	id, err := deviceid.Parse(deviceID)
	if err != nil {
		return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
	}
	if _, err := s.st.GetDevice(ctx, rid, id); err != nil {
		return nil, err
	}
	if len(p.Aliases) > 0 {
		if err := s.st.PatchDeviceAliases(ctx, rid, id, p.Aliases); err != nil {
			return nil, err
		}
	}
	if len(p.Attributes) > 0 {
		if err := s.st.PatchDeviceAttributes(ctx, rid, id, p.Attributes); err != nil {
			return nil, err
		}
	}
	if p.CredentialsInhibited != nil {
		if err := s.st.SetDeviceInhibited(ctx, rid, id, *p.CredentialsInhibited); err != nil {
			return nil, err
		}
	}
	d, err := s.st.GetDevice(ctx, rid, id)
	if err != nil {
		return nil, err
	}
	return s.deviceStatus(ctx, rid, d)
}

// --- groups -----------------------------------------------------------------

// CreateGroup creates a group with its initial device membership (upstream
// POST /groups requires a non-empty device list).
func (s *Service) CreateGroup(ctx context.Context, realm, name string, devices []string) error {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return err
	}
	if name == "" {
		return fmt.Errorf("%w: group_name can't be blank", ErrValidation)
	}
	if len(devices) == 0 {
		return fmt.Errorf("%w: a group must contain at least one device", ErrValidation)
	}
	ids, err := parseDeviceIDs(devices)
	if err != nil {
		return err
	}
	g, err := s.st.CreateGroup(ctx, rid, name)
	if err != nil {
		return err
	}
	for _, id := range ids {
		if err := s.st.AddGroupDevice(ctx, g.ID, rid, id); err != nil {
			return err
		}
	}
	return nil
}

// ListGroups returns the realm's group names.
func (s *Service) ListGroups(ctx context.Context, realm string) ([]string, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	groups, err := s.st.ListGroups(ctx, rid)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(groups))
	for i := range groups {
		names[i] = groups[i].Name
	}
	return names, nil
}

// ListGroupDevices returns the device IDs in a group.
func (s *Service) ListGroupDevices(ctx context.Context, realm, name string) ([]string, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	g, err := s.st.GetGroupByName(ctx, rid, name)
	if err != nil {
		return nil, err
	}
	ids, err := s.st.ListGroupDevices(ctx, g.ID)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(ids))
	for i := range ids {
		out[i] = ids[i].String()
	}
	return out, nil
}

// AddGroupDevice adds a device to a group.
func (s *Service) AddGroupDevice(ctx context.Context, realm, name, deviceID string) error {
	rid, g, err := s.group(ctx, realm, name)
	if err != nil {
		return err
	}
	id, err := deviceid.Parse(deviceID)
	if err != nil {
		return fmt.Errorf("%w: invalid device id", ErrValidation)
	}
	return s.st.AddGroupDevice(ctx, g.ID, rid, id)
}

// RemoveGroupDevice removes a device from a group.
func (s *Service) RemoveGroupDevice(ctx context.Context, realm, name, deviceID string) error {
	rid, g, err := s.group(ctx, realm, name)
	if err != nil {
		return err
	}
	id, err := deviceid.Parse(deviceID)
	if err != nil {
		return fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
	}
	return s.st.RemoveGroupDevice(ctx, g.ID, rid, id)
}

// group resolves a realm + group name to (realmID, group).
func (s *Service) group(ctx context.Context, realm, name string) (int16, *store.Group, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return 0, nil, err
	}
	g, err := s.st.GetGroupByName(ctx, rid, name)
	if err != nil {
		return 0, nil, err
	}
	return rid, g, nil
}

// --- helpers ----------------------------------------------------------------

func introspectionView(in map[string]store.InterfaceVersion) map[string]introspectionEntry {
	out := make(map[string]introspectionEntry, len(in))
	for name, v := range in {
		out[name] = introspectionEntry{Major: v.Major, Minor: v.Minor}
	}
	return out
}

func orEmptyStr(m map[string]string) map[string]string {
	if m == nil {
		return map[string]string{}
	}
	return m
}

func parseDeviceIDs(ids []string) ([]deviceid.ID, error) {
	out := make([]deviceid.ID, len(ids))
	for i, s := range ids {
		id, err := deviceid.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid device id %q", ErrValidation, s)
		}
		out[i] = id
	}
	return out, nil
}
