slug: channels-group-watch-membership-scope
verdict: done
at:  aeaeb47
ran: 2026-09-18T20:23:46Z on DietPi in 677s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/appengine/channels/room.go
→ Read internal/appengine/channels/ws.go
→ Read internal/engine/triggers/match.go [offset=340, limit=80]

$ ls /root/astrate-mule/internal/appengine/channels/
protocol.go
protocol_test.go
room.go
room_test.go
ws.go
ws_test.go

→ Read internal/appengine/channels/room_test.go
✗ Read internal/engine/stream/stream.go failed
Error: File not found: /root/astrate-mule/internal/engine/stream/stream.go
✱ Glob "internal/engine/stream/*.go" 2 matches

$ grep -rn "NewRegistry\|config.Registry\|Registry{" /root/astrate-mule/internal/appengine/ --include="*.go" | grep -v "_test.go" | head -30
/root/astrate-mule/internal/appengine/channels/ws.go:40:		reg:   NewRegistry(bus),
/root/astrate-mule/internal/appengine/channels/room.go:49:// NewRegistry creates a new Registry.
/root/astrate-mule/internal/appengine/channels/room.go:50:func NewRegistry(b Bus) *Registry {
/root/astrate-mule/internal/appengine/channels/room.go:51:	return &Registry{bus: b, rooms: make(map[string]*Room)}

→ Read internal/engine/stream/bus.go

$ grep -rn "NewRegistry\|NewAPI\|channels\.New" /root/astrate-mule --include="*.go" | grep -v "func " | head -30
/root/astrate-mule/internal/observability/metrics.go:28:	reg := prometheus.NewRegistry()
/root/astrate-mule/internal/realm/http_test.go:90:	NewAPI(svc, auth.NewMiddleware(st)).Mount(mux)
/root/astrate-mule/internal/realm/http.go:28:// NewAPI wires the realm service to its HTTP surface. mw provides the
/root/astrate-mule/internal/flowapi/http.go:25:// NewAPI wires the service to its HTTP surface.
/root/astrate-mule/internal/pairing/http_test.go:100:	api := NewAPI(svc, auth.NewMiddleware(st), apiCfg)
/root/astrate-mule/internal/pairing/http.go:51:// NewAPI wires the pairing service to its HTTP surface. mw provides the
/root/astrate-mule/internal/housekeeping/http_test.go:88:	NewAPI(NewService(st, sealer, nil, discardLogger()).
/root/astrate-mule/internal/housekeeping/http.go:26:// NewAPI wires the housekeeping service to its HTTP surface. instanceKeysPEM
/root/astrate-mule/internal/appengine/http_test.go:149:	NewAPI(NewService(st, sd, discardLogger()), auth.NewMiddleware(st)).Mount(mux)
/root/astrate-mule/internal/appengine/e2e_test.go:236:	channels.NewAPI(env.engine.Bus(), env.st).Mount(mux)
/root/astrate-mule/internal/appengine/channels/ws.go:37:// NewAPI creates an API backed by the given bus and key source.
/root/astrate-mule/internal/appengine/channels/ws.go:40:		reg:   NewRegistry(bus),
/root/astrate-mule/internal/appengine/channels/ws_test.go:54:	api := NewAPI(bus, fakeRealmKeys{realm})
/root/astrate-mule/internal/appengine/channels/ws_test.go:375:	api := NewAPI(bus, fakeRealmKeys{realm})
/root/astrate-mule/internal/appengine/channels/room.go:49:// NewRegistry creates a new Registry.
/root/astrate-mule/internal/appengine/channels/room_test.go:34:	reg := NewRegistry(bus)
/root/astrate-mule/internal/appengine/channels/room_test.go:141:	reg := NewRegistry(bus)
/root/astrate-mule/internal/appengine/channels/room_test.go:168:	reg := NewRegistry(bus)
/root/astrate-mule/internal/appengine/channels/room_test.go:226:	reg := NewRegistry(bus)
/root/astrate-mule/internal/appengine/channels/room_test.go:251:	reg := NewRegistry(bus)
/root/astrate-mule/internal/appengine/channels/room_test.go:261:	reg := NewRegistry(bus)
/root/astrate-mule/internal/appengine/channels/room_test.go:305:	reg := NewRegistry(bus)
/root/astrate-mule/internal/appengine/http.go:33:// NewAPI wires the AppEngine service to its HTTP surface.
/root/astrate-mule/internal/appengine/stream/ws.go:39:// NewAPI wires the socket to the engine bus and the a_ch middleware.
/root/astrate-mule/internal/appengine/stream/ws_test.go:60:	NewAPI(bus, auth.NewMiddleware(fakeKeys{realm})).Mount(mux)
/root/astrate-mule/internal/engine/triggers/actions_test.go:35:		cfg.Registerer = prometheus.NewRegistry()
/root/astrate-mule/internal/engine/stream/bus_test.go:111:	reg := prometheus.NewRegistry()
/root/astrate-mule/internal/flow/factory_test.go:25:	reg := flow.NewRegistry()
/root/astrate-mule/internal/flow/factory_test.go:103:	reg := flow.NewRegistry()
/root/astrate-mule/internal/flow/factory_test.go:120:	reg := flow.NewRegistry()

→ Read internal/appengine/channels/ws_test.go

$ grep -rln "GroupName\|group_name\|Groups()" /root/astrate-mule/internal --include="*.go" | head -30
/root/astrate-mule/internal/appengine/http_test.go
/root/astrate-mule/internal/appengine/mirror_test.go
/root/astrate-mule/internal/appengine/channels/ws.go
/root/astrate-mule/internal/appengine/channels/ws_test.go
/root/astrate-mule/internal/appengine/channels/room.go
/root/astrate-mule/internal/appengine/dashboard_compat_test.go
/root/astrate-mule/internal/appengine/http.go
/root/astrate-mule/internal/appengine/groups_parity_test.go
/root/astrate-mule/internal/engine/triggers/change_test.go
/root/astrate-mule/internal/engine/triggers/match.go
/root/astrate-mule/internal/engine/triggers/match_test.go

$ ls /root/astrate-mule/internal/appengine/ /root/astrate-mule/internal/appengine/appengine.go 2>/dev/null; ls /root/astrate-mule/internal/appengine/stream/ 2>/dev/null
/root/astrate-mule/internal/appengine/:
channels
dashboard_compat_test.go
data.go
downsample.go
downsample_test.go
e2e_test.go
formats_parity_test.go
groups_parity_test.go
http.go
http_test.go
links_test.go
mirror_test.go
patch_test.go
query_opts_test.go
service.go
stream
writeerror_test.go
ws.go
ws_test.go

$ grep -n "roup" /root/astrate-mule/internal/appengine/http.go | head -40
57:	// ones; Go 1.22 precedence resolves the overlap with the group routes
67:	mux.Handle("GET "+base+"/groups", h(a.listGroups))
68:	mux.Handle("POST "+base+"/groups", h(a.createGroup))
69:	mux.Handle("GET "+base+"/groups/{group}", h(a.getGroup))
70:	mux.Handle("GET "+base+"/groups/{group}/devices", h(a.listGroupDevices))
71:	mux.Handle("POST "+base+"/groups/{group}/devices", h(a.addGroupDevice))
72:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}", h(a.getGroupDevice))
73:	mux.Handle("PATCH "+base+"/groups/{group}/devices/{device}", h(a.patchGroupDevice))
74:	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}", h(a.removeGroupDevice))
75:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces", h(a.listInterfacesInGroup))
76:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}", h(a.getDataInGroup))
77:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.getDataInGroup))
78:	mux.Handle("PUT "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.putDataInGroup))
79:	mux.Handle("POST "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.putDataInGroup))
80:	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteDataInGroup))
187:func (a *API) listInterfacesInGroup(w http.ResponseWriter, r *http.Request) {
188:	names, err := a.svc.ListInterfacesInGroup(r.Context(), r.PathValue("realm"),
189:		r.PathValue("group"), r.PathValue("device"))
285:func (a *API) getDataInGroup(w http.ResponseWriter, r *http.Request) {
287:		return a.svc.GetDataInGroup(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device"),
340:func (a *API) putDataInGroup(w http.ResponseWriter, r *http.Request) {
350:	err := a.svc.PublishDataInGroup(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device"),
369:func (a *API) deleteDataInGroup(w http.ResponseWriter, r *http.Request) {
370:	err := a.svc.UnsetPropertyInGroup(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device"),
379:// --- groups -----------------------------------------------------------------
381:func (a *API) listGroups(w http.ResponseWriter, r *http.Request) {
382:	names, err := a.svc.ListGroups(r.Context(), r.PathValue("realm"))
390:// groupBody is the POST /groups wire shape. Devices is a pointer so a MISSING
393:type groupBody struct {
394:	GroupName string    `json:"group_name"`
398:func (a *API) createGroup(w http.ResponseWriter, r *http.Request) {
399:	var body groupBody
405:	case body.GroupName == "":
407:			FieldErrors{"group_name": {"can't be blank"}})
418:	if err := a.svc.CreateGroup(r.Context(), r.PathValue("realm"), body.GroupName, *body.Devices); err != nil {
426:// getGroup serves GET /groups/{group}: the group identity body, distinct from
428:func (a *API) getGroup(w http.ResponseWriter, r *http.Request) {
429:	if err := a.svc.GetGroup(r.Context(), r.PathValue("realm"), r.PathValue("group")); err != nil {
433:	_ = astarteapi.WriteData(w, http.StatusOK, map[string]string{"group_name": r.PathValue("group")})
436:func (a *API) listGroupDevices(w http.ResponseWriter, r *http.Request) {

$ grep -rn "Group\|groups" /root/astrate-mule/internal/appengine/service.go | head -40
3:// server-owned publishing, groups, and the live event socket. It is wire-shaped
41:// ErrGroupNotFound marks a missing group (maps to 404 "Group not found").
42:var ErrGroupNotFound = errors.New("Group not found") //nolint:staticcheck // ST1005: upstream wire text
44:// ErrGroupAlreadyExists marks a duplicate group name (maps to 409 "Group
47:var ErrGroupAlreadyExists = fmt.Errorf("Group already exists: %w", store.ErrAlreadyExists) //nolint:staticcheck // ST1005: upstream wire text
49:// ErrDeviceAlreadyInGroup marks re-adding an existing member (maps to 409
51:var ErrDeviceAlreadyInGroup = fmt.Errorf("Device already in group: %w", store.ErrAlreadyExists) //nolint:staticcheck // ST1005: upstream wire text
64:// missingGroup wraps a failed group resolution so it satisfies BOTH
65:// ErrGroupNotFound (→ 404 "Group not found") and store.ErrNotFound.
66:func missingGroup(name string) error {
67:	return fmt.Errorf("%w: %w: group %q", ErrGroupNotFound, store.ErrNotFound, name)
111:	Groups               []string                      `json:"groups"`
198:	groups, err := s.st.ListDeviceGroupsBatch(ctx, rid, ids)
204:		statuses[i] = deviceStatusView(&devs[i], groups[devs[i].ID])
243:	groups, err := s.st.ListDeviceGroups(ctx, rid, d.ID)
247:	return deviceStatusView(d, groups), nil
251:// pre-resolved group names) into the upstream status body. groups renders as
254:func deviceStatusView(d *store.Device, groups []string) *DeviceStatus {
255:	if groups == nil {
256:		groups = []string{}
264:		Groups:               groups,
381:// PatchGroupDevice applies a device patch inside a group (upstream
382:// PATCH /groups/{group}/devices/{device}). Unknown group → ErrGroupNotFound;
384:func (s *Service) PatchGroupDevice(ctx context.Context, realm, groupName, deviceID string, p DevicePatch) (*DeviceStatus, error) {
392:// GetGroupDevice returns one member device's status (upstream
393:// GET /groups/{group}/devices/{device}).
394:func (s *Service) GetGroupDevice(ctx context.Context, realm, groupName, deviceID string) (*DeviceStatus, error) {
403:// ordering: unknown group first (ErrGroupNotFound), then membership — a
411:	if _, err := s.st.GetGroupByName(ctx, rid, groupName); err != nil {
412:		return 0, deviceid.ID{}, nil, missingGroup(groupName)
418:	groups, err := s.st.ListDeviceGroups(ctx, rid, id)
422:	if !slices.Contains(groups, groupName) {
455:// --- groups -----------------------------------------------------------------
457:// CreateGroup creates a group with its initial device membership (upstream
458:// POST /groups requires a non-empty device list). Blank-name and empty-device
462:func (s *Service) CreateGroup(ctx context.Context, realm, name string, devices []string) error {
492:	g, err := s.st.CreateGroup(ctx, rid, name)
495:			return fmt.Errorf("%w: %w", ErrGroupAlreadyExists, err)
500:		if err := s.st.AddGroupDevice(ctx, g.ID, rid, id); err != nil {
507:// GetGroup resolves one group by name (upstream GET /groups/{group}).

$ grep -rn "AddGroupDevice\|ListGroupDevices\|GroupMembers\|GetGroupByName" /root/astrate-mule/internal/store --include="*.go" | head -30
/root/astrate-mule/internal/store/devices_test.go:354:		if err := s.AddGroupDevice(ctx, g.ID, realm.ID, d1); err != nil {
/root/astrate-mule/internal/store/devices_test.go:368:		members, err := s.ListGroupDevices(ctx, g.ID)
/root/astrate-mule/internal/store/groups.go:40:// GetGroupByName fetches one group.
/root/astrate-mule/internal/store/groups.go:41:func (s *Store) GetGroupByName(ctx context.Context, realmID int16, name string) (*Group, error) {
/root/astrate-mule/internal/store/groups.go:91:// AddGroupDevice adds a device to a group. The INSERT...SELECT guard ensures
/root/astrate-mule/internal/store/groups.go:95:func (s *Store) AddGroupDevice(ctx context.Context, groupID int64, realmID int16, deviceID deviceid.ID) error {
/root/astrate-mule/internal/store/groups.go:129:// ListGroupDevices returns the IDs of every device in a group, ordered.
/root/astrate-mule/internal/store/groups.go:130:func (s *Store) ListGroupDevices(ctx context.Context, groupID int64) ([]deviceid.ID, error) {
/root/astrate-mule/internal/store/groups.go:152:// ListGroupDevicesPage returns one page of member IDs ordered by device_id,
/root/astrate-mule/internal/store/groups.go:154:func (s *Store) ListGroupDevicesPage(ctx context.Context, groupID int64, offset, limit int) ([]deviceid.ID, error) {
/root/astrate-mule/internal/store/groups_test.go:30:	got, err := s.GetGroupByName(ctx, realm.ID, "fleet1")
/root/astrate-mule/internal/store/groups_test.go:32:		t.Fatalf("GetGroupByName: %v", err)
/root/astrate-mule/internal/store/groups_test.go:49:	if err := s.AddGroupDevice(ctx, group.ID, realm.ID, d1); err != nil {
/root/astrate-mule/internal/store/groups_test.go:50:		t.Fatalf("AddGroupDevice d1: %v", err)
/root/astrate-mule/internal/store/groups_test.go:52:	if err := s.AddGroupDevice(ctx, group.ID, realm.ID, d2); err != nil {
/root/astrate-mule/internal/store/groups_test.go:53:		t.Fatalf("AddGroupDevice d2: %v", err)
/root/astrate-mule/internal/store/groups_test.go:55:	if err := s.AddGroupDevice(ctx, group.ID, realm.ID, d1); !errors.Is(err, ErrAlreadyExists) {
/root/astrate-mule/internal/store/groups_test.go:60:	if err := s.AddGroupDevice(ctx, group.ID, realm.ID, ghost); !errors.Is(err, ErrNotFound) {
/root/astrate-mule/internal/store/groups_test.go:63:	if err := s.AddGroupDevice(ctx, group.ID+9999, realm.ID, d1); !errors.Is(err, ErrNotFound) {
/root/astrate-mule/internal/store/groups_test.go:68:	if err := s.AddGroupDevice(ctx, group.ID, otherRealm.ID, d1); !errors.Is(err, ErrNotFound) {
/root/astrate-mule/internal/store/groups_test.go:72:	members, err := s.ListGroupDevices(ctx, group.ID)
/root/astrate-mule/internal/store/groups_test.go:74:		t.Fatalf("ListGroupDevices: %v", err)
/root/astrate-mule/internal/store/groups_test.go:93:	members, err = s.ListGroupDevices(ctx, group.ID)
/root/astrate-mule/internal/store/groups_test.go:111:	if _, err := s.GetGroupByName(ctx, realm.ID, "fleet1"); !errors.Is(err, ErrNotFound) {
/root/astrate-mule/internal/store/groups_test.go:138:		if err := s.AddGroupDevice(ctx, m.g.ID, realm.ID, m.d); err != nil {
/root/astrate-mule/internal/store/realms_test.go:314:		if err := s.AddGroupDevice(ctx, group.ID, realm.ID, device); err != nil {
/root/astrate-mule/internal/store/realms_test.go:315:			t.Fatalf("AddGroupDevice: %v", err)

→ Read internal/appengine/e2e_test.go [offset=200, limit=70]
→ Read internal/appengine/service.go [offset=455, limit=120]

$ grep -rn "channels\." /root/astrate-mule/cmd /root/astrate-mule/internal --include="*.go" | grep -v "_test.go" | grep -v "channels.NewRegistry\|channels/room.go\|channels/ws.go\|channels/ws_test\|channels/room_test" | head -20
/root/astrate-mule/cmd/astrate/main.go:377:	channels.NewAPI(e.Bus(), st).Mount(mux)
/root/astrate-mule/internal/auth/claims.go:195:// see test/conformance/upstream/channels.json, where upstream refuses a join
/root/astrate-mule/internal/appengine/channels/protocol.go:50:		return fmt.Errorf("channels.Frame: invalid JSON: %w", err)
/root/astrate-mule/internal/appengine/channels/protocol.go:53:		return fmt.Errorf("channels.Frame: expected a 5-element JSON array, got %d elements", len(raw))
/root/astrate-mule/internal/appengine/channels/protocol.go:58:		return fmt.Errorf("channels.Frame: join_ref: %w", err)
/root/astrate-mule/internal/appengine/channels/protocol.go:63:		return fmt.Errorf("channels.Frame: ref: %w", err)
/root/astrate-mule/internal/appengine/channels/protocol.go:67:		return fmt.Errorf("channels.Frame: topic: %w", err)
/root/astrate-mule/internal/appengine/channels/protocol.go:70:		return fmt.Errorf("channels.Frame: event: %w", err)
/root/astrate-mule/internal/appengine/channels/protocol.go:98:			return Frame{}, fmt.Errorf("channels.Reply: marshal response: %w", err)
/root/astrate-mule/internal/appengine/channels/protocol.go:108:		return Frame{}, fmt.Errorf("channels.Reply: marshal payload: %w", err)

→ Read internal/appengine/e2e_test.go [limit=60, offset=120]
→ Read cmd/astrate/main.go [limit=60, offset=340]

$ grep -n "GetRealmByName\|func (s \*Store) GetRealm" /root/astrate-mule/internal/store/*.go | head
/root/astrate-mule/internal/store/realms.go:78:func (s *Store) GetRealm(ctx context.Context, id int16) (*Realm, error) {
/root/astrate-mule/internal/store/realms.go:90:// GetRealmByName fetches a realm by name.
/root/astrate-mule/internal/store/realms.go:91:func (s *Store) GetRealmByName(ctx context.Context, name string) (*Realm, error) {
/root/astrate-mule/internal/store/realms_test.go:50:		got, err := s.GetRealmByName(ctx, created.Name)
/root/astrate-mule/internal/store/realms_test.go:52:			t.Fatalf("GetRealmByName: %v", err)
/root/astrate-mule/internal/store/realms_test.go:133:		got, err := s.GetRealmByName(ctx, realm.Name)
/root/astrate-mule/internal/store/realms_test.go:135:			t.Fatalf("GetRealmByName: %v", err)
/root/astrate-mule/internal/store/realms_test.go:173:		got, err := s.GetRealmByName(ctx, created.Name)
/root/astrate-mule/internal/store/realms_test.go:175:			t.Fatalf("GetRealmByName: %v", err)
/root/astrate-mule/internal/store/realms_test.go:182:		bare, err := s.GetRealmByName(ctx, mustCreateRealm(t, s).Name)

→ Read internal/store/groups.go

$ grep -rn "integration\|testcontainers\|build integration" /root/astrate-mule/internal/store/*_test.go | head -15
/root/astrate-mule/internal/store/compression_test.go:1://go:build integration
/root/astrate-mule/internal/store/datastreams_test.go:1://go:build integration
/root/astrate-mule/internal/store/devices_test.go:1://go:build integration
/root/astrate-mule/internal/store/flows_test.go:1://go:build integration
/root/astrate-mule/internal/store/groups_test.go:1://go:build integration
/root/astrate-mule/internal/store/interfaces_test.go:1://go:build integration
/root/astrate-mule/internal/store/notify_test.go:1://go:build integration
/root/astrate-mule/internal/store/pipelines_test.go:1://go:build integration
/root/astrate-mule/internal/store/policies_test.go:1://go:build integration
/root/astrate-mule/internal/store/properties_test.go:1://go:build integration
/root/astrate-mule/internal/store/realms_test.go:1://go:build integration
/root/astrate-mule/internal/store/retention_test.go:1://go:build integration
/root/astrate-mule/internal/store/store_test.go:1://go:build integration
/root/astrate-mule/internal/store/triggers_test.go:1://go:build integration
/root/astrate-mule/internal/store/userblocks_test.go:1://go:build integration

$ grep -rn "DBPool\|pool " /root/astrate-mule/internal/store/store.go | head; ls /root/astrate-mule/internal/store/ | head -40
26:// defaultPoolMaxConns bounds the pool when the DSN does not specify
56:// Store is the shared database handle: a bounded pgx pool over an
60:	pool *pgxpool.Pool
68:// optional capabilities, and returns a ready Store. The pool is capped at
104:	// does not close the pool itself (pgx stdlib.OpenDBFromPool contract).
161:// Stat returns a snapshot of the connection pool statistics
162:// (docs/DESIGN.md §5.2: the DB-pool observability gauges read it).
compression_test.go
crypto.go
crypto_test.go
datastreams.go
datastreams_test.go
devices.go
devices_test.go
flows.go
flows_test.go
groups.go
groups_test.go
interfaces.go
interfaces_test.go
notify.go
notify_test.go
pipelines.go
pipelines_test.go
pipelines_validate_test.go
policies.go
policies_test.go
properties.go
properties_test.go
realms.go
realms_test.go
retention.go
retention_test.go
store.go
store_test.go
testdata
triggers.go
triggers_test.go
userblocks.go
userblocks_test.go

$ head -3 /root/astrate-mule/internal/appengine/e2e_test.go /root/astrate-mule/internal/appengine/channels/ws_test.go /root/astrate-mule/internal/appengine/channels/room_test.go; grep -rn "//go:build" /root/astrate-mule/internal/appengine/*_test.go | head
==> /root/astrate-mule/internal/appengine/e2e_test.go <==
//go:build integration && e2e

package appengine

==> /root/astrate-mule/internal/appengine/channels/ws_test.go <==
package channels

import (

==> /root/astrate-mule/internal/appengine/channels/room_test.go <==
package channels

import (
/root/astrate-mule/internal/appengine/dashboard_compat_test.go:1://go:build integration
/root/astrate-mule/internal/appengine/e2e_test.go:1://go:build integration && e2e
/root/astrate-mule/internal/appengine/formats_parity_test.go:1://go:build integration
/root/astrate-mule/internal/appengine/groups_parity_test.go:1://go:build integration
/root/astrate-mule/internal/appengine/http_test.go:1://go:build integration
/root/astrate-mule/internal/appengine/mirror_test.go:1://go:build integration
/root/astrate-mule/internal/appengine/patch_test.go:1://go:build integration

$ ls /root/astrate-mule/test/conformance/upstream/ 2>/dev/null | head; wc -l /root/astrate-mule/test/conformance/upstream/channels.json 2>/dev/null
README.md
channels.json
channels.transcript.txt
record
recordchannels
rest-errors.json
rest-errors.transcript.txt
verify-server-writes.json
verify-server-writes.transcript.txt
verify-versions.json
492 /root/astrate-mule/test/conformance/upstream/channels.json

✱ Grep "group_name|groups/|Group" in test/conformance/upstream/channels.json · 58 matches

$ grep -rn "func (id ID) String\|func (idID\|func (id ID)" /root/astrate-mule/pkg/deviceid/*.go | head
/root/astrate-mule/pkg/deviceid/deviceid.go:63:func (id ID) String() string {
/root/astrate-mule/pkg/deviceid/deviceid.go:87:func (id ID) UUID() string {

$ ls /root/astrate-mule/pkg/deviceid/
deviceid.go
deviceid_test.go

→ Read internal/store/realms.go [limit=30, offset=85]

$ sed -n '1,25p' /root/astrate-mule/internal/store/store.go
// Package store is Astrate's single PostgreSQL/TimescaleDB access layer
// (docs/DESIGN.md §1.3): every domain package reads and writes the database
// through it, and it imports none of them. It owns connection pooling, the
// embedded schema migrations (docs/DESIGN.md §2.2–2.5), CA-key sealing, and
// one repository file per aggregate (realms, interfaces, devices, properties,
// datastreams, groups, triggers) plus LISTEN/NOTIFY plumbing.
package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/astrate-platform/astrate/migrations"
)

Now I understand the structure. Let me look at the remaining pieces before editing.
✱ Grep "st \*store.Store|st:\s+st|cdEnv struct" in internal/appengine/e2e_test.go · 5 matches
→ Read internal/appengine/e2e_test.go [limit=40, offset=80]

← Edit internal/appengine/channels/room.go
Index: /root/astrate-mule/internal/appengine/channels/room.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room.go
+++ /root/astrate-mule/internal/appengine/channels/room.go
@@ -1,7 +1,8 @@
 package channels
 
 import (
+	"context"
 	"encoding/json"
 	"fmt"
 	"sync"
 	"sync/atomic"
@@ -17,8 +18,16 @@
 type Bus interface {
 	Subscribe(realm string, f stream.Filter, buffer int) (<-chan stream.Event, func())
 }
 
+// GroupMembers lists the devices that belong to a group, so a group watch
+// delivers only its members' events (a group-scoped token authorizes
+// WATCH::groups/<name>/…, not every device in the realm). A *store.Store
+// satisfies it; tests supply a fake.
+type GroupMembers interface {
+	ListGroupMembers(ctx context.Context, realm, group string) ([]string, error)
+}
+
 // WatchRequest is the watch payload DTO, exactly as the client sends it.
 //
 // GroupName is a top-level field, measured against upstream (channels.json,
 // 2026-08-22): a group_name nested inside simple_trigger is refused by

Now I have the full picture. Implementing the changes. Starting with `room.go`:
← Edit internal/appengine/channels/room.go
Index: /root/astrate-mule/internal/appengine/channels/room.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room.go
+++ /root/astrate-mule/internal/appengine/channels/room.go
@@ -42,23 +42,29 @@
 	SimpleTrigger json.RawMessage `json:"simple_trigger"`
 }
 
 type watchEntry struct {
-	name     string
-	deviceID string
-	tg       *triggers.Trigger
+	name      string
+	deviceID  string
+	groupName string
+	// members holds the group's device IDs as encoded strings, resolved at
+	// Watch time; nil for a device-scoped watch (no membership filter).
+	members map[string]struct{}
+	tg      *triggers.Trigger
 }
 
 // Registry is a topic→room map.
 type Registry struct {
-	bus   Bus
-	mu    sync.Mutex
-	rooms map[string]*Room
+	bus    Bus
+	groups GroupMembers
+	mu     sync.Mutex
+	rooms  map[string]*Room
 }
 
-// NewRegistry creates a new Registry.
-func NewRegistry(b Bus) *Registry {
-	return &Registry{bus: b, rooms: make(map[string]*Room)}
+// NewRegistry creates a new Registry. groups resolves group watches'
+// membership; nil makes a group watch fail at Watch time.
+func NewRegistry(bus Bus, groups GroupMembers) *Registry {
+	return &Registry{bus: bus, groups: groups, rooms: make(map[string]*Room)}
 }
 
 // Join returns the room for topic, creating it on first join.
 func (r *Registry) Join(realm, topic string) *Room {
@@ -69,8 +75,9 @@
 	}
 	ch, cancel := r.bus.Subscribe(realm, stream.Filter{}, 0)
 	rm := &Room{
 		reg:      r,
+		realm:    realm,
 		topic:    topic,
 		members:  make(map[*Member]struct{}),
 		events:   ch,
 		cancelFn: cancel,

← Edit internal/appengine/channels/room.go
Index: /root/astrate-mule/internal/appengine/channels/room.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room.go
+++ /root/astrate-mule/internal/appengine/channels/room.go
@@ -96,8 +96,9 @@
 
 // Room is one Phoenix topic.
 type Room struct {
 	reg      *Registry
+	realm    string
 	topic    string
 	watches  map[string]watchEntry
 	watchMu  sync.Mutex
 	members  map[*Member]struct{}

Now the `Watch` method and membership resolution:
← Edit internal/appengine/channels/room.go
Index: /root/astrate-mule/internal/appengine/channels/room.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room.go
+++ /root/astrate-mule/internal/appengine/channels/room.go
@@ -135,10 +135,15 @@
 	rm.memberMu.Unlock()
 	return m
 }
 
-// Watch compiles and stores a watch.
-func (rm *Room) Watch(req WatchRequest) error {
+// Watch compiles and stores a watch. A group watch resolves its membership
+// now and keeps the member set in the entry, so dispatch filters deliveries by
+// it: without that, a token scoped to WATCH::groups/<name>/… would receive
+// every realm event matching the trigger's interface/path/on. The group must
+// resolve here — a missing group, or no resolver in the registry, refuses the
+// watch.
+func (rm *Room) Watch(ctx context.Context, req WatchRequest) error {
 	rm.watchMu.Lock()
 	defer rm.watchMu.Unlock()
 	if req.Name == "" {
 		return fmt.Errorf("name is required")
@@ -149,19 +154,44 @@
 	tg, err := triggers.CompileCondition(req.Name, req.SimpleTrigger)
 	if err != nil {
 		return fmt.Errorf("CompileCondition: %w", err)
 	}
+	var members map[string]struct{}
+	if req.GroupName != "" {
+		members, err = rm.resolveMembers(ctx, req.GroupName)
+		if err != nil {
+			return err
+		}
+	}
 	if rm.watches == nil {
 		rm.watches = make(map[string]watchEntry)
 	}
 	rm.watches[req.Name] = watchEntry{
-		name:     req.Name,
-		deviceID: req.DeviceID,
-		tg:       tg,
+		name:      req.Name,
+		deviceID:  req.DeviceID,
+		groupName: req.GroupName,
+		members:   members,
+		tg:        tg,
 	}
 	return nil
 }
 
+// resolveMembers resolves one group's devices to a lookup set of encoded IDs.
+func (rm *Room) resolveMembers(ctx context.Context, group string) (map[string]struct{}, error) {
+	if rm.reg.groups == nil {
+		return nil, fmt.Errorf("group watch of %q: no group-membership resolver", group)
+	}
+	ids, err := rm.reg.groups.ListGroupMembers(ctx, rm.realm, group)
+	if err != nil {
+		return nil, fmt.Errorf("resolving group %q: %w", group, err)
+	}
+	set := make(map[string]struct{}, len(ids))
+	for _, id := range ids {
+		set[id] = struct{}{}
+	}
+	return set, nil
+}
+
 // Unwatch removes a watch, reporting whether one was there.
 func (rm *Room) Unwatch(name string) bool {
 	rm.watchMu.Lock()
 	defer rm.watchMu.Unlock()

Now the dispatch filter:
← Edit internal/appengine/channels/room.go
Index: /root/astrate-mule/internal/appengine/channels/room.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room.go
+++ /root/astrate-mule/internal/appengine/channels/room.go
@@ -243,8 +243,16 @@
 for _, w := range watches {
 	if w.deviceID != "" && w.deviceID != ev.DeviceID {
 		continue
 	}
+	// A group watch is scoped to its members, resolved at Watch time;
+	// an event from outside the group is the leak a group-scoped
+	// token must not see.
+	if w.members != nil {
+		if _, ok := w.members[ev.DeviceID]; !ok {
+			continue
+		}
+	}
 
 	var matched bool
 	switch ev.Kind {
 	case stream.KindIncomingData:

← Edit internal/appengine/channels/ws.go
Index: /root/astrate-mule/internal/appengine/channels/ws.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/ws.go
+++ /root/astrate-mule/internal/appengine/channels/ws.go
@@ -33,12 +33,12 @@
 	keys  auth.KeySource
 	cache *auth.Cache
 }
 
-// NewAPI creates an API backed by the given bus and key source.
-func NewAPI(bus Bus, keys auth.KeySource) *API {
+// NewAPI creates an API backed by the given bus, key source and group resolver.
+func NewAPI(bus Bus, keys auth.KeySource, groups GroupMembers) *API {
 	return &API{
-		reg:   NewRegistry(bus),
+		reg:   NewRegistry(bus, groups),
 		keys:  keys,
 		cache: auth.NewCache(auth.DefaultCacheSize),
 	}
 }

Now `ws.go` — the `API` plumbing and the `Watch` call:
← Edit internal/appengine/channels/ws.go
Index: /root/astrate-mule/internal/appengine/channels/ws.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/ws.go
+++ /root/astrate-mule/internal/appengine/channels/ws.go
@@ -355,9 +355,9 @@
 	s.sendErr(f, "unauthorized")
 	return
 }
 
-if err := j.room.Watch(req); err != nil {
+if err := j.room.Watch(s.ctx, req); err != nil {
 	s.sendErr(f, "invalid trigger")
 	return
 }
 rep, _ := OK(f, nil)

Now the store method so `*store.Store` satisfies the interface:
← Edit internal/store/groups.go
Index: /root/astrate-mule/internal/store/groups.go
===================================================================
--- /root/astrate-mule/internal/store/groups.go
+++ /root/astrate-mule/internal/store/groups.go
@@ -151,30 +151,9 @@
 
 // ListGroupDevicesPage returns one page of member IDs ordered by device_id,
 // skipping the first offset rows; at most limit rows are returned.
 func (s *Store) ListGroupDevicesPage(ctx context.Context, groupID int64, offset, limit int) ([]deviceid.ID, error) {
-	rows, err := s.pool.Query(ctx,
-		`SELECT device_id FROM group_devices WHERE group_id = $1 ORDER BY device_id OFFSET $2 LIMIT $3`,
-		groupID, offset, limit)
-	if err != nil {
-		return nil, fmt.Errorf("store: listing devices of group %d: %w", groupID, err)
-	}
-	defer rows.Close()
 
-	var out []deviceid.ID
-	for rows.Next() {
-		var u pgtype.UUID
-		if err := rows.Scan(&u); err != nil {
-			return nil, fmt.Errorf("store: scanning group device: %w", err)
-		}
-		out = append(out, deviceid.ID(u.Bytes))
-	}
-	if err := rows.Err(); err != nil {
-		return nil, fmt.Errorf("store: listing devices of group %d: %w", groupID, err)
-	}
-	return out, nil
-}
-
 // ListDeviceGroups returns the names of every group the device belongs to.
 func (s *Store) ListDeviceGroups(ctx context.Context, realmID int16, deviceID deviceid.ID) ([]string, error) {
 	rows, err := s.pool.Query(ctx, `
 		SELECT g.name FROM groups g

→ Read internal/store/groups.go [limit=55, offset=128]
I broke the file. Restoring `ListGroupDevicesPage` and adding `ListGroupMembers`:

← Edit internal/store/groups.go
Index: /root/astrate-mule/internal/store/groups.go
===================================================================
--- /root/astrate-mule/internal/store/groups.go
+++ /root/astrate-mule/internal/store/groups.go
@@ -151,9 +151,54 @@
 
 // ListGroupDevicesPage returns one page of member IDs ordered by device_id,
 // skipping the first offset rows; at most limit rows are returned.
 func (s *Store) ListGroupDevicesPage(ctx context.Context, groupID int64, offset, limit int) ([]deviceid.ID, error) {
+	rows, err := s.pool.Query(ctx,
+		`SELECT device_id FROM group_devices WHERE group_id = $1 ORDER BY device_id OFFSET $2 LIMIT $3`,
+		groupID, offset, limit)
+	if err != nil {
+		return nil, fmt.Errorf("store: listing devices of group %d: %w", groupID, err)
+	}
+	defer rows.Close()
 
+	var out []deviceid.ID
+	for rows.Next() {
+		var u pgtype.UUID
+		if err := rows.Scan(&u); err != nil {
+			return nil, fmt.Errorf("store: scanning group device: %w", err)
+		}
+		out = append(out, deviceid.ID(u.Bytes))
+	}
+	if err := rows.Err(); err != nil {
+		return nil, fmt.Errorf("store: listing devices of group %d: %w", groupID, err)
+	}
+	return out, nil
+}
+
+// ListGroupMembers resolves the encoded IDs of every device in a named group
+// of a named realm — the shape the Channels socket's group-watch delivery
+// filter needs (channels.Room.Watch), which must not itself know a realm's id
+// or a group's id.
+func (s *Store) ListGroupMembers(ctx context.Context, realm, group string) ([]string, error) {
+	r, err := s.GetRealmByName(ctx, realm)
+	if err != nil {
+		return nil, err
+	}
+	g, err := s.GetGroupByName(ctx, r.ID, group)
+	if err != nil {
+		return nil, err
+	}
+	ids, err := s.ListGroupDevices(ctx, g.ID)
+	if err != nil {
+		return nil, err
+	}
+	out := make([]string, len(ids))
+	for i, id := range ids {
+		out[i] = id.String()
+	}
+	return out, nil
+}
+
 // ListDeviceGroups returns the names of every group the device belongs to.
 func (s *Store) ListDeviceGroups(ctx context.Context, realmID int16, deviceID deviceid.ID) ([]string, error) {
 	rows, err := s.pool.Query(ctx, `
 		SELECT g.name FROM groups g

← Edit cmd/astrate/main.go
Index: /root/astrate-mule/cmd/astrate/main.go
===================================================================
--- /root/astrate-mule/cmd/astrate/main.go
+++ /root/astrate-mule/cmd/astrate/main.go
@@ -373,9 +373,9 @@
 appengine.NewAPI(appengine.NewService(st, e, log), mw).Mount(mux)
 apstream.NewAPI(e.Bus(), mw).Mount(mux)
 // Phoenix Channels socket (phoenix.js V2), alongside the Astrate-native
 // socket above: two protocols, one bus.
-channels.NewAPI(e.Bus(), st).Mount(mux)
+channels.NewAPI(e.Bus(), st, st).Mount(mux)
 // Flow operator API: pipelines CRUD + start/stop/status (v2.0).
 flowapi.NewAPI(flowSvc, mw).Mount(mux)
 
 // Upstream-parity per-service health endpoints (the dashboard's API

Now update the wiring at all call sites:
← Edit internal/appengine/e2e_test.go
Index: /root/astrate-mule/internal/appengine/e2e_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/e2e_test.go
+++ /root/astrate-mule/internal/appengine/e2e_test.go
@@ -232,9 +232,9 @@
 env, newDevice := newCDEnv(t)
 dev := newDevice(t)
 
 mux := http.NewServeMux()
-channels.NewAPI(env.engine.Bus(), env.st).Mount(mux)
+channels.NewAPI(env.engine.Bus(), env.st, env.st).Mount(mux)
 srv := httptest.NewServer(mux)
 defer srv.Close()
 
 claims := jwt.MapClaims{"a_ch": []string{"JOIN::.*", "WATCH::.*"}, "exp": time.Now().Add(time.Hour).Unix()}

← Edit internal/appengine/channels/ws_test.go
Index: /root/astrate-mule/internal/appengine/channels/ws_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/ws_test.go
+++ /root/astrate-mule/internal/appengine/channels/ws_test.go
@@ -50,9 +50,9 @@
 	der, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
 	pub := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
 	realm := &store.Realm{Name: testRealmName, JWTPublicKeysPEM: []string{pub}}
 	bus := &fakeBus{ch: make(chan stream.Event)}
-	api := NewAPI(bus, fakeRealmKeys{realm})
+	api := NewAPI(bus, fakeRealmKeys{realm}, nil)
 	return api, key
 }
 
 func TestPreUpgrade_MissingToken(t *testing.T) {

← Edit internal/appengine/channels/ws_test.go
Index: /root/astrate-mule/internal/appengine/channels/ws_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/ws_test.go
+++ /root/astrate-mule/internal/appengine/channels/ws_test.go
@@ -371,9 +371,9 @@
 	der, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
 	pub := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
 	realm := &store.Realm{Name: testRealmName, JWTPublicKeysPEM: []string{pub}}
 	bus := &fakeBus{ch: make(chan stream.Event, 16)}
-	api := NewAPI(bus, fakeRealmKeys{realm})
+	api := NewAPI(bus, fakeRealmKeys{realm}, nil)
 	return api, key, bus
 }
 
 func joinTopic(t *testing.T, conn *websocket.Conn, topic, ref string) {

Now update the six `NewRegistry(bus)` calls in room_test.go and add the new test with a fake resolver:
$ cd /root/astrate-mule && sed -i 's/NewRegistry(bus)/NewRegistry(bus, nil)/g' internal/appengine/channels/room_test.go && grep -n "NewRegistry" internal/appengine/channels/room_test.go
34:	reg := NewRegistry(bus, nil)
141:	reg := NewRegistry(bus, nil)
168:	reg := NewRegistry(bus, nil)
226:	reg := NewRegistry(bus, nil)
251:	reg := NewRegistry(bus, nil)
261:	reg := NewRegistry(bus, nil)
305:	reg := NewRegistry(bus, nil)

$ cd /root/astrate-mule && grep -rn "\.Watch(" internal/appengine/channels/*_test.go
internal/appengine/channels/room_test.go:59:		if err := rm.Watch(req); err != nil {
internal/appengine/channels/room_test.go:149:	if err := rm.Watch(dataReq); err != nil {
internal/appengine/channels/room_test.go:175:	if err := rm.Watch(dataReq); err != nil {
internal/appengine/channels/room_test.go:289:			err := rm.Watch(tt.req)

← Edit internal/appengine/channels/room_test.go
Index: /root/astrate-mule/internal/appengine/channels/room_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room_test.go
+++ /root/astrate-mule/internal/appengine/channels/room_test.go
@@ -1,7 +1,8 @@
 package channels
 
 import (
+	"context"
 	"encoding/json"
 	"strings"
 	"sync/atomic"
 	"testing"

← Edit internal/appengine/channels/room_test.go
Index: /root/astrate-mule/internal/appengine/channels/room_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room_test.go
+++ /root/astrate-mule/internal/appengine/channels/room_test.go
@@ -56,9 +56,9 @@
 	SimpleTrigger: json.RawMessage(`{"type":"data_trigger","on":"incoming_data","interface_name":"*","value_match_operator":"*","match_path":"/*"}`),
 }
 
 for _, req := range []WatchRequest{connReq, discReq, errReq, dataReq} {
-	if err := rm.Watch(req); err != nil {
+	if err := rm.Watch(context.Background(), req); err != nil {
 		t.Fatalf("Watch(%s): %v", req.Name, err)
 	}
 }
 

← Edit internal/appengine/channels/room_test.go
Index: /root/astrate-mule/internal/appengine/channels/room_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room_test.go
+++ /root/astrate-mule/internal/appengine/channels/room_test.go
@@ -146,9 +146,9 @@
 	Name:          "datatrigger-D",
 	DeviceID:      validDeviceIDA,
 	SimpleTrigger: json.RawMessage(`{"type":"data_trigger","on":"incoming_data","interface_name":"*","value_match_operator":"*","match_path":"/*"}`),
 }
-if err := rm.Watch(dataReq); err != nil {
+if err := rm.Watch(context.Background(), dataReq); err != nil {
 	t.Fatalf("Watch: %v", err)
 }
 
 m := rm.AddMember(0)

← Edit internal/appengine/channels/room_test.go
Index: /root/astrate-mule/internal/appengine/channels/room_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room_test.go
+++ /root/astrate-mule/internal/appengine/channels/room_test.go
@@ -172,9 +172,9 @@
 dataReq := WatchRequest{
 	Name:          "datatrigger",
 	SimpleTrigger: json.RawMessage(`{"type":"data_trigger","on":"incoming_data","interface_name":"*","value_match_operator":"*","match_path":"/*"}`),
 }
-if err := rm.Watch(dataReq); err != nil {
+if err := rm.Watch(context.Background(), dataReq); err != nil {
 	t.Fatalf("Watch: %v", err)
 }
 
 slow := rm.AddMember(1)

← Edit internal/appengine/channels/room_test.go
Index: /root/astrate-mule/internal/appengine/channels/room_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room_test.go
+++ /root/astrate-mule/internal/appengine/channels/room_test.go
@@ -286,9 +286,9 @@
 
 for _, tt := range tests {
 	t.Run(tt.name, func(t *testing.T) {
 		before := rm.Watches()
-		err := rm.Watch(tt.req)
+		err := rm.Watch(context.Background(), tt.req)
 		if err == nil {
 			t.Fatal("expected error, got nil")
 		}
 		if tt.wantErr != "" && !strings.Contains(err.Error(), tt.wantErr) {

Now add the fake resolver and the group-membership test:
← Edit internal/appengine/channels/room_test.go
Index: /root/astrate-mule/internal/appengine/channels/room_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room_test.go
+++ /root/astrate-mule/internal/appengine/channels/room_test.go
@@ -24,8 +24,18 @@
 func (f *fakeBus) cancelCalled() bool {
 	return f.cancelCount.Load() > 0
 }
 
+// fakeGroups resolves group membership from a static map: group name → encoded
+// device IDs. It is the container-free stand-in for *store.Store.
+type fakeGroups struct {
+	byGroup map[string][]string
+}
+
+func (f *fakeGroups) ListGroupMembers(_ context.Context, _ string, group string) ([]string, error) {
+	return f.byGroup[group], nil
+}
+
 const (
 	validDeviceIDA = "f0VMRgIBAQAAAAAAAAAAAA"
 	validDeviceIDB = "f0VMRgIBAQAAAAAAAAAAAQ"
 )

← Edit internal/appengine/channels/room_test.go
Index: /root/astrate-mule/internal/appengine/channels/room_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room_test.go
+++ /root/astrate-mule/internal/appengine/channels/room_test.go
@@ -173,8 +173,61 @@
 	case <-time.After(200 * time.Millisecond):
 	}
 }
 
+func TestGroupWatchDeliversMembersOnly(t *testing.T) {
+	bus := &fakeBus{ch: make(chan stream.Event, 16)}
+	resolver := &fakeGroups{byGroup: map[string][]string{
+		"probe": {validDeviceIDA},
+	}}
+	reg := NewRegistry(bus, resolver)
+	rm := reg.Join("test", "rooms:test:devices")
+
+	// A group watch: WATCH::groups/probe/... authorizes the watch, but
+	// delivery must be confined to the group's members — the boundary a
+	// group-scoped token is granted. The trigger itself carries no device or
+	// group scope (upstream keeps group_name off the trigger), so without the
+	// resolved membership this watch would receive every realm event.
+	req := WatchRequest{
+		Name:      "group-data",
+		GroupName: "probe",
+		SimpleTrigger: json.RawMessage(
+			`{"type":"data_trigger","on":"incoming_data","interface_name":"*","value_match_operator":"*","match_path":"/*"}`),
+	}
+	if err := rm.Watch(context.Background(), req); err != nil {
+		t.Fatalf("Watch: %v", err)
+	}
+
+	m := rm.AddMember(0)
+	defer m.Leave()
+
+	now := time.Now().UTC().Truncate(time.Millisecond)
+
+	// The non-member publishes first. A leak would deliver it; the filter must
+	// swallow it.
+	bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: validDeviceIDB, Interface: "org.example.V1", Path: "/sensor/temp", Value: 1.0, InterfaceMajor: 1, Timestamp: now}
+	select {
+	case ev := <-m.Events():
+		t.Fatalf("non-member event leaked through group watch: %+v", ev)
+	case <-time.After(300 * time.Millisecond):
+	}
+
+	bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: validDeviceIDA, Interface: "org.example.V1", Path: "/sensor/temp", Value: 2.0, InterfaceMajor: 1, Timestamp: now}
+	select {
+	case ev := <-m.Events():
+		switch body := ev.Event.(type) {
+		case triggers.IncomingDataEvent:
+			if body.DeviceID != "" && body.Path != "/sensor/temp" {
+				t.Errorf("Path = %q, want /sensor/temp", body.Path)
+			}
+		default:
+			t.Fatalf("unexpected event type %T", ev.Event)
+		}
+	case <-time.After(time.Second):
+		t.Fatal("timeout waiting for member's event")
+	}
+}
+
 func TestFullMailboxDropsNeverBlocks(t *testing.T) {
 	bus := &fakeBus{ch: make(chan stream.Event, 16)}
 	reg := NewRegistry(bus, nil)
 	rm := reg.Join("test", "rooms:test:devices")

$ cd /root/astrate-mule && grep -n "type IncomingDataEvent\|type SimpleEvent" internal/engine/triggers/triggers.go
grep: internal/engine/triggers/triggers.go: No such file or directory

$ cd /root/astrate-mule && grep -rn "type IncomingDataEvent\|type SimpleEvent" internal/engine/triggers/*.go
internal/engine/triggers/events.go:33:type SimpleEvent struct {
internal/engine/triggers/events.go:63:type IncomingDataEvent struct {

→ Read internal/engine/triggers/events.go [limit=50, offset=33]
The envelope's DeviceID is correct. Let me fix my test to assert on `ev.DeviceID`:

← Edit internal/appengine/channels/room_test.go
Index: /root/astrate-mule/internal/appengine/channels/room_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room_test.go
+++ /root/astrate-mule/internal/appengine/channels/room_test.go
@@ -213,16 +213,14 @@
 
 	bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: validDeviceIDA, Interface: "org.example.V1", Path: "/sensor/temp", Value: 2.0, InterfaceMajor: 1, Timestamp: now}
 	select {
 	case ev := <-m.Events():
-		switch body := ev.Event.(type) {
-		case triggers.IncomingDataEvent:
-			if body.DeviceID != "" && body.Path != "/sensor/temp" {
-				t.Errorf("Path = %q, want /sensor/temp", body.Path)
-			}
-		default:
-			t.Fatalf("unexpected event type %T", ev.Event)
+		if ev.DeviceID != validDeviceIDA {
+			t.Errorf("event device %q, want %q", ev.DeviceID, validDeviceIDA)
 		}
+		if body, ok := ev.Event.(triggers.IncomingDataEvent); !ok || body.Path != "/sensor/temp" {
+			t.Fatalf("unexpected event %#v", ev.Event)
+		}
 	case <-time.After(time.Second):
 		t.Fatal("timeout waiting for member's event")
 	}
 }

← Edit internal/appengine/channels/room_test.go
Index: /root/astrate-mule/internal/appengine/channels/room_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/room_test.go
+++ /root/astrate-mule/internal/appengine/channels/room_test.go
@@ -342,8 +342,16 @@
 		name:    "invalid on value",
 		req:     WatchRequest{Name: "t", DeviceID: validDeviceIDA, SimpleTrigger: json.RawMessage(`{"type":"device_trigger","on":"not_a_real_event","device_id":"` + validDeviceIDA + `"}`)},
 		wantErr: "CompileCondition",
 	},
+	{
+		// A group watch with no resolver must refuse rather than accept
+		// and silently leak every realm event (the registry is wired nil
+		// in this reject-table).
+		name:    "group watch without a resolver",
+		req:     WatchRequest{Name: "t", GroupName: "probe", SimpleTrigger: json.RawMessage(`{"type":"data_trigger","on":"incoming_data","interface_name":"*","value_match_operator":"*","match_path":"/*"}`)},
+		wantErr: "resolver",
+	},
 }
 
 for _, tt := range tests {
 	t.Run(tt.name, func(t *testing.T) {

Let me also pin that a nil resolver refuses a group watch (guards against re-introducing the leak by dropping the wiring):
Now build and run the gate:
$ cd /root/astrate-mule && gofmt -l . && go vet ./internal/appengine/channels/ ./internal/store/ ./cmd/astrate/ ./internal/appengine/ 2>&1 | head -30
(no output)

$ cd /root/astrate-mule && go test ./internal/appengine/channels/ ./internal/store/ ./cmd/ 2>&1 | tail -20
# ./cmd
no Go files in /root/astrate-mule/cmd
FAIL	./cmd [setup failed]
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	23.448s
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
FAIL

$ cd /root/astrate-mule && go test ./... 2>&1 | tail -30
ok  	github.com/astrate-platform/astrate/internal/appengine	0.076s
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	0.122s
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

Full suite passes. Now the race detector on the package I touched (disk/concurrency changed in `dispatch` and `Watch`):
$ cd /root/astrate-mule && go test -race ./internal/appengine/channels/ 2>&1 | tail -5
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine/channels	0.048s
FAIL

Race detector is unavailable here (39-bit VMA kernel — the Pi case). Running the standard Pi gate `go vet ./... && go test ./...` to finish:
$ cd /root/astrate-mule && go vet ./... 2>&1 | tail -5 && echo "VET-DONE" && gofmt -l . && echo "FMT-DONE"
VET-DONE
FMT-DONE

Let me review the complete diff before reporting:
$ cd /root/astrate-mule && git diff --stat && echo "---" && git diff internal/appengine/channels/room.go internal/appengine/channels/ws.go internal/store/groups.go cmd/astrate/main.go
cmd/astrate/main.go                      |  2 +-
 internal/appengine/channels/room.go      | 83 +++++++++++++++++++++++-----
 internal/appengine/channels/room_test.go | 92 ++++++++++++++++++++++++++++----
 internal/appengine/channels/ws.go        |  8 +--
 internal/appengine/channels/ws_test.go   |  4 +-
 internal/appengine/e2e_test.go           |  2 +-
 internal/store/groups.go                 | 24 +++++++++
 7 files changed, 182 insertions(+), 33 deletions(-)
---
diff --git a/cmd/astrate/main.go b/cmd/astrate/main.go
index 2adfb86..a600f44 100644
--- a/cmd/astrate/main.go
+++ b/cmd/astrate/main.go
@@ -374,7 +374,7 @@ func mountAPIs(cfg config.Config, st *store.Store, e *engine.Engine, b *broker.B
 	apstream.NewAPI(e.Bus(), mw).Mount(mux)
 	// Phoenix Channels socket (phoenix.js V2), alongside the Astrate-native
 	// socket above: two protocols, one bus.
-	channels.NewAPI(e.Bus(), st).Mount(mux)
+	channels.NewAPI(e.Bus(), st, st).Mount(mux)
 	// Flow operator API: pipelines CRUD + start/stop/status (v2.0).
 	flowapi.NewAPI(flowSvc, mw).Mount(mux)
 
diff --git a/internal/appengine/channels/room.go b/internal/appengine/channels/room.go
index 2babf25..4c19734 100644
--- a/internal/appengine/channels/room.go
+++ b/internal/appengine/channels/room.go
@@ -1,6 +1,7 @@
 package channels
 
 import (
+	"context"
 	"encoding/json"
 	"fmt"
 	"sync"
@@ -18,6 +19,14 @@ type Bus interface {
 	Subscribe(realm string, f stream.Filter, buffer int) (<-chan stream.Event, func())
 }
 
+// GroupMembers lists the devices that belong to a group, so a group watch
+// delivers only its members' events (a group-scoped token authorizes
+// WATCH::groups/<name>/…, not every device in the realm). A *store.Store
+// satisfies it; tests supply a fake.
+type GroupMembers interface {
+	ListGroupMembers(ctx context.Context, realm, group string) ([]string, error)
+}
+
 // WatchRequest is the watch payload DTO, exactly as the client sends it.
 //
 // GroupName is a top-level field, measured against upstream (channels.json,
@@ -34,21 +43,27 @@ type WatchRequest struct {
 }
 
 type watchEntry struct {
-	name     string
-	deviceID string
-	tg       *triggers.Trigger
+	name      string
+	deviceID  string
+	groupName string
+	// members holds the group's device IDs as encoded strings, resolved at
+	// Watch time; nil for a device-scoped watch (no membership filter).
+	members map[string]struct{}
+	tg      *triggers.Trigger
 }
 
 // Registry is a topic→room map.
 type Registry struct {
-	bus   Bus
-	mu    sync.Mutex
-	rooms map[string]*Room
+	bus    Bus
+	groups GroupMembers
+	mu     sync.Mutex
+	rooms  map[string]*Room
 }
 
-// NewRegistry creates a new Registry.
-func NewRegistry(b Bus) *Registry {
-	return &Registry{bus: b, rooms: make(map[string]*Room)}
+// NewRegistry creates a new Registry. groups resolves group watches'
+// membership; nil makes a group watch fail at Watch time.
+func NewRegistry(bus Bus, groups GroupMembers) *Registry {
+	return &Registry{bus: bus, groups: groups, rooms: make(map[string]*Room)}
 }
 
 // Join returns the room for topic, creating it on first join.
@@ -61,6 +76,7 @@ func (r *Registry) Join(realm, topic string) *Room {
 	ch, cancel := r.bus.Subscribe(realm, stream.Filter{}, 0)
 	rm := &Room{
 		reg:      r,
+		realm:    realm,
 		topic:    topic,
 		members:  make(map[*Member]struct{}),
 		events:   ch,
@@ -81,6 +97,7 @@ func (r *Registry) Rooms() int {
 // Room is one Phoenix topic.
 type Room struct {
 	reg      *Registry
+	realm    string
 	topic    string
 	watches  map[string]watchEntry
 	watchMu  sync.Mutex
@@ -119,8 +136,13 @@ func (rm *Room) AddMember(buffer int) *Member {
 	return m
 }
 
-// Watch compiles and stores a watch.
-func (rm *Room) Watch(req WatchRequest) error {
+// Watch compiles and stores a watch. A group watch resolves its membership
+// now and keeps the member set in the entry, so dispatch filters deliveries by
+// it: without that, a token scoped to WATCH::groups/<name>/… would receive
+// every realm event matching the trigger's interface/path/on. The group must
+// resolve here — a missing group, or no resolver in the registry, refuses the
+// watch.
+func (rm *Room) Watch(ctx context.Context, req WatchRequest) error {
 	rm.watchMu.Lock()
 	defer rm.watchMu.Unlock()
 	if req.Name == "" {
@@ -133,17 +155,42 @@ func (rm *Room) Watch(req WatchRequest) error {
 	if err != nil {
 		return fmt.Errorf("CompileCondition: %w", err)
 	}
+	var members map[string]struct{}
+	if req.GroupName != "" {
+		members, err = rm.resolveMembers(ctx, req.GroupName)
+		if err != nil {
+			return err
+		}
+	}
 	if rm.watches == nil {
 		rm.watches = make(map[string]watchEntry)
 	}
 	rm.watches[req.Name] = watchEntry{
-		name:     req.Name,
-		deviceID: req.DeviceID,
-		tg:       tg,
+		name:      req.Name,
+		deviceID:  req.DeviceID,
+		groupName: req.GroupName,
+		members:   members,
+		tg:        tg,
 	}
 	return nil
 }
 
+// resolveMembers resolves one group's devices to a lookup set of encoded IDs.
+func (rm *Room) resolveMembers(ctx context.Context, group string) (map[string]struct{}, error) {
+	if rm.reg.groups == nil {
+		return nil, fmt.Errorf("group watch of %q: no group-membership resolver", group)
+	}
+	ids, err := rm.reg.groups.ListGroupMembers(ctx, rm.realm, group)
+	if err != nil {
+		return nil, fmt.Errorf("resolving group %q: %w", group, err)
+	}
+	set := make(map[string]struct{}, len(ids))
+	for _, id := range ids {
+		set[id] = struct{}{}
+	}
+	return set, nil
+}
+
 // Unwatch removes a watch, reporting whether one was there.
 func (rm *Room) Unwatch(name string) bool {
 	rm.watchMu.Lock()
@@ -197,6 +244,14 @@ func (rm *Room) dispatch() {
 			if w.deviceID != "" && w.deviceID != ev.DeviceID {
 				continue
 			}
+			// A group watch is scoped to its members, resolved at Watch time;
+			// an event from outside the group is the leak a group-scoped
+			// token must not see.
+			if w.members != nil {
+				if _, ok := w.members[ev.DeviceID]; !ok {
+					continue
+				}
+			}
 
 			var matched bool
 			switch ev.Kind {
diff --git a/internal/appengine/channels/ws.go b/internal/appengine/channels/ws.go
index 91de541..9a3a608 100644
--- a/internal/appengine/channels/ws.go
+++ b/internal/appengine/channels/ws.go
@@ -34,10 +34,10 @@ type API struct {
 	cache *auth.Cache
 }
 
-// NewAPI creates an API backed by the given bus and key source.
-func NewAPI(bus Bus, keys auth.KeySource) *API {
+// NewAPI creates an API backed by the given bus, key source and group resolver.
+func NewAPI(bus Bus, keys auth.KeySource, groups GroupMembers) *API {
 	return &API{
-		reg:   NewRegistry(bus),
+		reg:   NewRegistry(bus, groups),
 		keys:  keys,
 		cache: auth.NewCache(auth.DefaultCacheSize),
 	}
@@ -356,7 +356,7 @@ func (s *session) handleWatch(f Frame) {
 		return
 	}
 
-	if err := j.room.Watch(req); err != nil {
+	if err := j.room.Watch(s.ctx, req); err != nil {
 		s.sendErr(f, "invalid trigger")
 		return
 	}
diff --git a/internal/store/groups.go b/internal/store/groups.go
index 0618e93..6e5628a 100644
--- a/internal/store/groups.go
+++ b/internal/store/groups.go
@@ -174,6 +174,30 @@ func (s *Store) ListGroupDevicesPage(ctx context.Context, groupID int64, offset,
 	return out, nil
 }
 
+// ListGroupMembers resolves the encoded IDs of every device in a named group
+// of a named realm — the shape the Channels socket's group-watch delivery
+// filter needs (channels.Room.Watch), which must not itself know a realm's id
+// or a group's id.
+func (s *Store) ListGroupMembers(ctx context.Context, realm, group string) ([]string, error) {
+	r, err := s.GetRealmByName(ctx, realm)
+	if err != nil {
+		return nil, err
+	}
+	g, err := s.GetGroupByName(ctx, r.ID, group)
+	if err != nil {
+		return nil, err
+	}
+	ids, err := s.ListGroupDevices(ctx, g.ID)
+	if err != nil {
+		return nil, err
+	}
+	out := make([]string, len(ids))
+	for i, id := range ids {
+		out[i] = id.String()
+	}
+	return out, nil
+}
+
 // ListDeviceGroups returns the names of every group the device belongs to.
 func (s *Store) ListDeviceGroups(ctx context.Context, realmID int16, deviceID deviceid.ID) ([]string, error) {
 	rows, err := s.pool.Query(ctx, `

$ cd /root/astrate-mule && go test ./internal/appengine/channels/ -run 'TestGroupWatch|TestRejectedWatches' -v 2>&1 | tail -25
=== RUN   TestGroupWatchDeliversMembersOnly
--- PASS: TestGroupWatchDeliversMembersOnly (0.30s)
=== RUN   TestRejectedWatches
=== RUN   TestRejectedWatches/empty_name
=== RUN   TestRejectedWatches/missing_simple_trigger
=== RUN   TestRejectedWatches/invalid_on_value
=== RUN   TestRejectedWatches/group_watch_without_a_resolver
--- PASS: TestRejectedWatches (0.00s)
    --- PASS: TestRejectedWatches/empty_name (0.00s)
    --- PASS: TestRejectedWatches/missing_simple_trigger (0.00s)
    --- PASS: TestRejectedWatches/invalid_on_value (0.00s)
    --- PASS: TestRejectedWatches/group_watch_without_a_resolver (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	0.339s

Both new tests pass and the full diff is coherent. The `go test -race` binary can't run on this kernel (39-bit VMA — the Pi limitation the rules warn about), so I ran the plain `go test ./...` + `go vet ./...` + `gofmt` gate.

Done: group watches now resolve membership at Watch time via a new GroupMembers resolver and dispatch filters deliveries to members; wired *store.Store in production, nil resolver refuses group watches.
Files: internal/appengine/channels/room.go, internal/appengine/channels/ws.go, internal/appengine/channels/room_test.go, internal/appengine/channels/ws_test.go, internal/store/groups.go, cmd/astrate/main.go, internal/appengine/e2e_test.go
Verified: gofmt -l . -> pass; go vet ./... -> pass; go test ./... -> pass. (go test -race ./... aborted here: unsupported VMA range — ran plain tests only.)
Unsure: The member-filter rule itself is asserted by the task and stands regardless; upstream behavior for a *nonexistent* group at Watch time is inferred (refused as "invalid trigger" when resolution errors) — not measured on the Legion, which is off and only needed "if disputed". The membership resolution runs under watchMu like trigger compilation, so a slow DB stalls a room's dispatch briefly — accepted, watch-time cost only.
Follow-ups: none
