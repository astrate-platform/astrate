// Package appengine is the AppEngine API (docs/DESIGN.md §3.7, ROADMAP §8.2):
// the operator/application surface for device status, interface data queries,
// server-owned publishing, groups, and the live event socket. It is wire-shaped
// to upstream astarte_appengine_api so astartectl and applications work
// unmodified, reading through the store and writing server-owned values through
// the engine.
package appengine

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// DefaultDeviceLimit is the device-list page size when the caller gives none.
const DefaultDeviceLimit = 100

// ErrValidation wraps a well-formed request that violates a rule (maps to 422).
var ErrValidation = errors.New("appengine: validation failed")

// Device-PATCH error taxonomy, frozen to upstream astarte_appengine_api 1.2.2
// (FallbackController + ErrorView): the detail strings are what astartectl and
// the dashboard match on.
var (
	ErrInvalidAlias         = errors.New("Invalid alias")           //nolint:staticcheck // ST1005: detail strings are upstream astarte_appengine_api 1.2.2 wire text
	ErrAliasAlreadyInUse    = errors.New("Alias already in use")    //nolint:staticcheck // ST1005: upstream wire text
	ErrAliasTagNotFound     = errors.New("Alias tag not found")     //nolint:staticcheck // ST1005: upstream wire text
	ErrInvalidAttributes    = errors.New("Invalid attributes")      //nolint:staticcheck // ST1005: upstream wire text
	ErrAttributeKeyNotFound = errors.New("Attribute key not found") //nolint:staticcheck // ST1005: upstream wire text
)

// ErrGroupNotFound marks a missing group (maps to 404 "Group not found").
var ErrGroupNotFound = errors.New("Group not found") //nolint:staticcheck // ST1005: upstream wire text

// ErrGroupAlreadyExists marks a duplicate group name (maps to 409 "Group
// already exists"); it also satisfies store.ErrAlreadyExists so generic
// duplicate handling keeps working.
var ErrGroupAlreadyExists = fmt.Errorf("Group already exists: %w", store.ErrAlreadyExists) //nolint:staticcheck // ST1005: upstream wire text

// ErrDeviceAlreadyInGroup marks re-adding an existing member (maps to 409
// "Device already in group"); it too satisfies store.ErrAlreadyExists.
var ErrDeviceAlreadyInGroup = fmt.Errorf("Device already in group: %w", store.ErrAlreadyExists) //nolint:staticcheck // ST1005: upstream wire text

// FieldErrors carries per-field validation failures rendered as the
// Phoenix-changeset envelope {"errors":{"<field>":["<message>",...]}} (422).
type FieldErrors map[string][]string

func (fe FieldErrors) Error() string { return "appengine: invalid request payload" }

// addf appends one formatted message under field.
func (fe FieldErrors) addf(field, format string, args ...any) {
	fe[field] = append(fe[field], fmt.Sprintf(format, args...))
}

// missingGroup wraps a failed group resolution so it satisfies BOTH
// ErrGroupNotFound (→ 404 "Group not found") and store.ErrNotFound.
func missingGroup(name string) error {
	return fmt.Errorf("%w: %w: group %q", ErrGroupNotFound, store.ErrNotFound, name)
}

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

// DevicePage is one device-list page: IDs, or full statuses when details was
// requested. Next is the next-page cursor ("" at the end).
type DevicePage struct {
	IDs      []string
	Statuses []*DeviceStatus
	Next     string
}

// ListDevices returns one page of devices (upstream GET /devices). after is
// the cursor (exclusive); limit <= 0 selects DefaultDeviceLimit; details
// switches the page from bare IDs to full status objects (the dashboard's
// device list).
func (s *Service) ListDevices(ctx context.Context, realm string, after string, limit int, details bool) (*DevicePage, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = DefaultDeviceLimit
	}
	var cursor *deviceid.ID
	if after != "" {
		id, err := deviceid.Parse(after)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid cursor", ErrValidation)
		}
		cursor = &id
	}
	devs, err := s.st.ListDevices(ctx, rid, cursor, limit+1)
	if err != nil {
		return nil, err
	}
	page := &DevicePage{}
	if len(devs) > limit {
		page.Next = devs[limit-1].ID.String()
		devs = devs[:limit]
	}
	if !details {
		page.IDs = make([]string, len(devs))
		for i := range devs {
			page.IDs[i] = devs[i].ID.String()
		}
		return page, nil
	}
	statuses, err := s.deviceStatusBatch(ctx, rid, devs)
	if err != nil {
		return nil, err
	}
	page.Statuses = statuses
	return page, nil
}

// DevicesStats counts the realm's devices (upstream GET /stats/devices).
func (s *Service) DevicesStats(ctx context.Context, realm string) (total, connected int64, err error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return 0, 0, err
	}
	return s.st.DeviceStats(ctx, rid)
}

// deviceStatusBatch projects a page of stored devices into status bodies with
// one batched group-membership query (details=true listings would otherwise
// go N+1).
func (s *Service) deviceStatusBatch(ctx context.Context, rid int16, devs []store.Device) ([]*DeviceStatus, error) {
	ids := make([]deviceid.ID, len(devs))
	for i := range devs {
		ids[i] = devs[i].ID
	}
	groups, err := s.st.ListDeviceGroupsBatch(ctx, rid, ids)
	if err != nil {
		return nil, err
	}
	statuses := make([]*DeviceStatus, len(devs))
	for i := range devs {
		statuses[i] = deviceStatusView(&devs[i], groups[devs[i].ID])
	}
	return statuses, nil
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
	return deviceStatusView(d, groups), nil
}

// deviceStatusView is the pure projection of a stored device (plus its
// pre-resolved group names) into the upstream status body. groups renders as
// [] rather than null for group-less devices (upstream shape; the dashboard's
// DTO validation expects an array).
func deviceStatusView(d *store.Device, groups []string) *DeviceStatus {
	if groups == nil {
		groups = []string{}
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
	return ds
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
	d, err := s.st.GetDevice(ctx, rid, id)
	if err != nil {
		return nil, err
	}
	return s.applyPatch(ctx, rid, id, d, p)
}

// PatchDeviceByAlias applies a device patch addressed by alias (upstream
// PATCH /devices-by-alias/{alias}). Unknown alias behaves like GetDeviceByAlias.
func (s *Service) PatchDeviceByAlias(ctx context.Context, realm, alias string, p DevicePatch) (*DeviceStatus, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	d, err := s.st.GetDeviceByAlias(ctx, rid, alias)
	if err != nil {
		return nil, err
	}
	return s.applyPatch(ctx, rid, d.ID, d, p)
}

// applyPatch validates and applies p to the already-fetched device row,
// preserving upstream merge_device_status's validation order: alias format,
// alias ownership realm-wide, attribute format, then the changeset's
// missing-tag/key rejections at apply time.
func (s *Service) applyPatch(ctx context.Context, rid int16, id deviceid.ID, d *store.Device, p DevicePatch) (*DeviceStatus, error) {
	for tag, value := range p.Aliases {
		if tag == "" || (value != nil && *value == "") {
			return nil, ErrInvalidAlias
		}
	}
	if len(p.Aliases) > 0 {
		taken, err := s.aliasValuesTaken(ctx, rid, id, d, p.Aliases)
		if err != nil {
			return nil, err
		}
		if taken {
			return nil, ErrAliasAlreadyInUse
		}
	}
	for key := range p.Attributes {
		if key == "" {
			return nil, ErrInvalidAttributes
		}
	}
	for tag, value := range p.Aliases {
		if value == nil {
			if _, ok := d.Aliases[tag]; !ok {
				return nil, ErrAliasTagNotFound
			}
		}
	}
	for key, value := range p.Attributes {
		if value == nil {
			if _, ok := d.Attributes[key]; !ok {
				return nil, ErrAttributeKeyNotFound
			}
		}
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
	d2, err := s.st.GetDevice(ctx, rid, id)
	if err != nil {
		return nil, err
	}
	return s.deviceStatus(ctx, rid, d2)
}

// PatchGroupDevice applies a device patch inside a group (upstream
// PATCH /groups/{group}/devices/{device}). Unknown group → ErrGroupNotFound;
// non-member device behaves like an unknown device.
func (s *Service) PatchGroupDevice(ctx context.Context, realm, groupName, deviceID string, p DevicePatch) (*DeviceStatus, error) {
	rid, id, d, err := s.groupMember(ctx, realm, groupName, deviceID)
	if err != nil {
		return nil, err
	}
	return s.applyPatch(ctx, rid, id, d, p)
}

// GetGroupDevice returns one member device's status (upstream
// GET /groups/{group}/devices/{device}).
func (s *Service) GetGroupDevice(ctx context.Context, realm, groupName, deviceID string) (*DeviceStatus, error) {
	rid, _, d, err := s.groupMember(ctx, realm, groupName, deviceID)
	if err != nil {
		return nil, err
	}
	return s.deviceStatus(ctx, rid, d)
}

// groupMember resolves (realm, group name, device id) with upstream's
// ordering: unknown group first (ErrGroupNotFound), then membership — a
// registered-but-not-member (or unparseable/unknown) device is store.ErrNotFound,
// indistinguishable from a missing device.
func (s *Service) groupMember(ctx context.Context, realm, groupName, deviceID string) (int16, deviceid.ID, *store.Device, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return 0, deviceid.ID{}, nil, err
	}
	if _, err := s.st.GetGroupByName(ctx, rid, groupName); err != nil {
		return 0, deviceid.ID{}, nil, missingGroup(groupName)
	}
	id, err := deviceid.Parse(deviceID)
	if err != nil {
		return 0, deviceid.ID{}, nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
	}
	groups, err := s.st.ListDeviceGroups(ctx, rid, id)
	if err != nil {
		return 0, deviceid.ID{}, nil, err
	}
	if !slices.Contains(groups, groupName) {
		return 0, deviceid.ID{}, nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
	}
	d, err := s.st.GetDevice(ctx, rid, id)
	if err != nil {
		return 0, deviceid.ID{}, nil, err
	}
	return rid, id, d, nil
}

// aliasValuesTaken reports whether any alias value the patch touches is owned
// by another device in the realm (upstream find_all_aliases ownership check).
// Deleted tags only contribute their value when the device actually holds
// them, matching upstream's Map.take over the current aliases.
func (s *Service) aliasValuesTaken(ctx context.Context, rid int16, id deviceid.ID,
	d *store.Device, patch map[string]*string,
) (bool, error) {
	values := make([]string, 0, len(patch))
	for tag, value := range patch {
		if value != nil {
			values = append(values, *value)
			continue
		}
		if v, ok := d.Aliases[tag]; ok {
			values = append(values, v)
		}
	}
	if len(values) == 0 {
		return false, nil
	}
	return s.st.AliasValuesTaken(ctx, rid, id, values)
}

// --- groups -----------------------------------------------------------------

// CreateGroup creates a group with its initial device membership (upstream
// POST /groups requires a non-empty device list). Blank-name and empty-device
// rejections live at the handler level as field errors; here every listed id
// must parse AND name an existing device row — a single failure rejects the
// whole body with the upstream changeset message, before any row is created.
func (s *Service) CreateGroup(ctx context.Context, realm, name string, devices []string) error {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return err
	}
	fe := FieldErrors{}
	var ids []deviceid.ID
	seen := make(map[deviceid.ID]struct{}, len(devices))
	for _, raw := range devices {
		id, err := deviceid.Parse(raw)
		if err != nil {
			fe.addf("devices", "must exist (%s not found)", raw)
			continue
		}
		if _, err := s.st.GetDevice(ctx, rid, id); err != nil {
			if errors.Is(err, store.ErrNotFound) {
				fe.addf("devices", "must exist (%s not found)", raw)
				continue
			}
			return err
		}
		// Duplicate ids inside one body are accepted and inserted once.
		if _, dup := seen[id]; !dup {
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	if len(fe) > 0 {
		return fe
	}
	g, err := s.st.CreateGroup(ctx, rid, name)
	if err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			return fmt.Errorf("%w: %w", ErrGroupAlreadyExists, err)
		}
		return err
	}
	for _, id := range ids {
		if err := s.st.AddGroupDevice(ctx, g.ID, rid, id); err != nil {
			return err
		}
	}
	return nil
}

// GetGroup resolves one group by name (upstream GET /groups/{group}).
func (s *Service) GetGroup(ctx context.Context, realm, name string) error {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return err
	}
	if _, err := s.st.GetGroupByName(ctx, rid, name); err != nil {
		return missingGroup(name)
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

// ListGroupDevices returns one page of the devices in a group — bare IDs, or
// full status objects with details (the dashboard's group page). The bare
// listing paginates with an OFFSET cursor carried inside a UUID-v1-format
// from_token (upstream's own tokens are insertion-time v1 uuids; any
// well-formed v1 string is accepted); details stays unpaginated.
func (s *Service) ListGroupDevices(ctx context.Context, realm, name string, details bool, fromToken string, limit int) (*DevicePage, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	offset := 0
	if fromToken != "" {
		var ok bool
		if offset, ok = parseGroupToken(fromToken); !ok {
			return nil, FieldErrors{"from_token": {"is invalid"}}
		}
	}
	if limit < 0 {
		return nil, FieldErrors{"limit": {"must be greater than or equal to 0"}}
	}
	if limit == 0 {
		limit = DefaultDeviceLimit
	}
	g, err := s.st.GetGroupByName(ctx, rid, name)
	if err != nil {
		return nil, missingGroup(name)
	}
	page := &DevicePage{}
	if details {
		ids, err := s.st.ListGroupDevices(ctx, g.ID)
		if err != nil {
			return nil, err
		}
		devs := make([]store.Device, 0, len(ids))
		for _, id := range ids {
			d, err := s.st.GetDevice(ctx, rid, id)
			if err != nil {
				return nil, err
			}
			devs = append(devs, *d)
		}
		statuses, err := s.deviceStatusBatch(ctx, rid, devs)
		if err != nil {
			return nil, err
		}
		page.Statuses = statuses
		return page, nil
	}
	rows, err := s.st.ListGroupDevicesPage(ctx, g.ID, offset, limit+1)
	if err != nil {
		return nil, err
	}
	if len(rows) > limit {
		rows = rows[:limit]
		page.Next = groupTokenFor(offset + limit)
	}
	page.IDs = make([]string, len(rows))
	for i := range rows {
		page.IDs[i] = rows[i].String()
	}
	return page, nil
}

// groupTokenFor renders offset as a UUID-v1-format cursor string:
// time_low = offset (big-endian), time_mid = 0, version nibble = 1,
// clock-seq variant = 0b10, node = 0.
func groupTokenFor(offset int) string {
	var u deviceid.ID
	// The cursor format is 32-bit by design (UUID v1 time_low): offsets above
	// math.MaxUint32 cannot occur — limit is capped and offset comes from a
	// token parsed back through parseGroupToken (uint32).
	binary.BigEndian.PutUint32(u[0:4], uint32(offset)) //nolint:gosec // G115: see comment
	u[6] = 0x10                                        // version 1
	u[8] = 0x80                                        // RFC 4122 variant
	return u.UUID()
}

// parseGroupToken accepts only strings that parse as a canonical UUID whose
// version NIBBLE is 1 (upstream :uuid.is_v1()); returns (offset, true) using
// time_low, or (0, false) otherwise. The zero fields are not demanded —
// upstream accepts any well-formed v1 uuid.
func parseGroupToken(s string) (int, bool) {
	u, err := deviceid.FromUUID(s)
	if err != nil || u[6]>>4 != 1 {
		return 0, false
	}
	return int(binary.BigEndian.Uint32(u[0:4])), true
}

// AddGroupDevice adds a device to a group.
func (s *Service) AddGroupDevice(ctx context.Context, realm, name, deviceID string) error {
	rid, g, err := s.group(ctx, realm, name)
	if err != nil {
		return err
	}
	id, err := deviceid.Parse(deviceID)
	if err != nil {
		return FieldErrors{"device_id": {"is not a valid device id"}}
	}
	if err := s.st.AddGroupDevice(ctx, g.ID, rid, id); err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			return fmt.Errorf("%w: %w", ErrDeviceAlreadyInGroup, err)
		}
		return err
	}
	return nil
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

// --- mirror: by-alias / by-group interface access ---------------------------

// resolveOnDevice finishes a mirror resolution with the same introspection +
// GetInterface steps resolve() ends with.
func (s *Service) resolveOnDevice(ctx context.Context, rid int16, d *store.Device, ifaceName string) (*resolved, error) {
	v, ok := d.Introspection[ifaceName]
	if !ok {
		return nil, fmt.Errorf("%w: interface %s not in device introspection", store.ErrNotFound, ifaceName)
	}
	si, err := s.st.GetInterface(ctx, rid, ifaceName, v.Major)
	if err != nil {
		return nil, err
	}
	return &resolved{rid: rid, id: d.ID, iface: si}, nil
}

// resolveByAlias maps (realm, alias, interface name) to a resolved interface.
func (s *Service) resolveByAlias(ctx context.Context, realm, alias, ifaceName string) (*resolved, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	d, err := s.st.GetDeviceByAlias(ctx, rid, alias)
	if err != nil {
		return nil, err
	}
	return s.resolveOnDevice(ctx, rid, d, ifaceName)
}

// resolveInGroup maps (realm, group, device id, interface name) to a resolved
// interface. Membership gates before the interface lookup: upstream answers
// "Device not found" for a non-member even when the interface would also be
// missing.
func (s *Service) resolveInGroup(ctx context.Context, realm, groupName, deviceID, ifaceName string) (*resolved, error) {
	rid, _, d, err := s.groupMember(ctx, realm, groupName, deviceID)
	if err != nil {
		return nil, err
	}
	return s.resolveOnDevice(ctx, rid, d, ifaceName)
}

// readResolved dispatches an already-resolved interface read on its type —
// the shared tail of GetData and its mirror variants.
func (s *Service) readResolved(ctx context.Context, r *resolved, path string, opts QueryOpts) (any, error) {
	if r.iface.Type == interfaceschema.Properties {
		// See GetData: reducing a snapshot is refused, not silently ignored.
		if opts.DownsamplePoints > 0 {
			return nil, fmt.Errorf("%w: downsample_to is not supported on properties interfaces", ErrValidation)
		}
		return s.propertiesData(ctx, r, path)
	}
	return s.datastreamData(ctx, r, path, opts)
}

// GetDataByAlias reads an interface endpoint on an alias-addressed device
// (upstream GET /devices-by-alias/{alias}/interfaces/{iface}[/{path}]).
func (s *Service) GetDataByAlias(ctx context.Context, realm, alias, ifaceName, path string, opts QueryOpts) (any, error) {
	r, err := s.resolveByAlias(ctx, realm, alias, ifaceName)
	if err != nil {
		return nil, err
	}
	return s.readResolved(ctx, r, path, opts)
}

// GetDataInGroup reads an interface endpoint on a group member (upstream
// GET /groups/{group}/devices/{device}/interfaces/{iface}[/{path}]).
func (s *Service) GetDataInGroup(ctx context.Context, realm, groupName, deviceID, ifaceName, path string, opts QueryOpts) (any, error) {
	r, err := s.resolveInGroup(ctx, realm, groupName, deviceID, ifaceName)
	if err != nil {
		return nil, err
	}
	return s.readResolved(ctx, r, path, opts)
}

// PublishDataByAlias publishes through an alias address; the resolution gates
// access and carries the real device ID to the engine write.
func (s *Service) PublishDataByAlias(ctx context.Context, realm, alias, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
	r, err := s.resolveByAlias(ctx, realm, alias, ifaceName)
	if err != nil {
		return err
	}
	return s.PublishData(ctx, realm, r.id.String(), ifaceName, path, value, ts)
}

// PublishDataInGroup publishes through a group address.
func (s *Service) PublishDataInGroup(ctx context.Context, realm, groupName, deviceID, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
	r, err := s.resolveInGroup(ctx, realm, groupName, deviceID, ifaceName)
	if err != nil {
		return err
	}
	return s.PublishData(ctx, realm, r.id.String(), ifaceName, path, value, ts)
}

// UnsetPropertyByAlias unsets through an alias address.
func (s *Service) UnsetPropertyByAlias(ctx context.Context, realm, alias, ifaceName, path string) error {
	r, err := s.resolveByAlias(ctx, realm, alias, ifaceName)
	if err != nil {
		return err
	}
	return s.UnsetProperty(ctx, realm, r.id.String(), ifaceName, path)
}

// UnsetPropertyInGroup unsets through a group address.
func (s *Service) UnsetPropertyInGroup(ctx context.Context, realm, groupName, deviceID, ifaceName, path string) error {
	r, err := s.resolveInGroup(ctx, realm, groupName, deviceID, ifaceName)
	if err != nil {
		return err
	}
	return s.UnsetProperty(ctx, realm, r.id.String(), ifaceName, path)
}

// listInterfacesNames renders a device's introspection names sorted ascending.
func listInterfacesNames(d *store.Device) []string {
	names := make([]string, 0, len(d.Introspection))
	for name := range d.Introspection {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// ListInterfaces returns a device's introspection interface names, sorted
// ascending (upstream GET .../devices/{d}/interfaces).
func (s *Service) ListInterfaces(ctx context.Context, realm, deviceID string) ([]string, error) {
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
	return listInterfacesNames(d), nil
}

// ListInterfacesByAlias is ListInterfaces addressed by alias.
func (s *Service) ListInterfacesByAlias(ctx context.Context, realm, alias string) ([]string, error) {
	rid, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	d, err := s.st.GetDeviceByAlias(ctx, rid, alias)
	if err != nil {
		return nil, err
	}
	return listInterfacesNames(d), nil
}

// ListInterfacesInGroup is ListInterfaces inside a group (membership gate first).
func (s *Service) ListInterfacesInGroup(ctx context.Context, realm, groupName, deviceID string) ([]string, error) {
	_, _, d, err := s.groupMember(ctx, realm, groupName, deviceID)
	if err != nil {
		return nil, err
	}
	return listInterfacesNames(d), nil
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
