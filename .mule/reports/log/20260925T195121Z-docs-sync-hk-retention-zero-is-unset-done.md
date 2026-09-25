slug: docs-sync-hk-retention-zero-is-unset
verdict: done
at:  c2ad442
ran: 2026-09-25T19:51:21Z on DietPi in 325s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ ls /root/astrate-mule && ls /root/astrate-mule/docs /root/astrate-mule/docs/api 2>/dev/null
AGENTS.md
Dockerfile
LICENSE
Makefile
README.md
astrate
bench
cmd
deploy
docker-compose.yml
docs
examples
flow-container-echo
go.mod
go.sum
internal
migrations
opencode.json
pkg
site-dist
test
tools
/root/astrate-mule/docs:
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
UPSTREAM-EXPERIMENTAL.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui

/root/astrate-mule/docs/api:
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

→ Read docs/api/astarte_housekeeping_api.yaml [offset=100, limit=240]

$ sed -n 195,230p /root/astrate-mule/internal/housekeeping/http.go; echo ---; sed -n 160,180p /root/astrate-mule/internal/store/realms.go; echo ---; sed -n 140,170p /root/astrate-mule/internal/housekeeping/service.go
return
	}
	if pb.DatastreamMaximumStorageRetention.present && !pb.DatastreamMaximumStorageRetention.null && pb.DatastreamMaximumStorageRetention.val < 0 {
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"datastream_maximum_storage_retention": {"is invalid"}})
		return
	}

	u := RealmUpdate{}
	if pb.JWTPublicKeyPEM.present {
		u.PatchJWTPublicKeyPEM = true
		u.SetJWTPublicKeyPEM = pb.JWTPublicKeyPEM.val
	}
	if pb.DeviceRegistrationLimit.present {
		if pb.DeviceRegistrationLimit.null {
			u.ClearRegistrationLimit = true
		} else {
			u.PatchRegistrationLimit = true
			u.SetRegistrationLimit = pb.DeviceRegistrationLimit.val
		}
	}
	if pb.DatastreamMaximumStorageRetention.present {
		if pb.DatastreamMaximumStorageRetention.null || pb.DatastreamMaximumStorageRetention.val == 0 {
			u.ClearRetention = true
		} else {
			u.PatchRetention = true
			u.SetRetention = pb.DatastreamMaximumStorageRetention.val
		}
	}

	if err := a.svc.UpdateRealm(r.Context(), r.PathValue("realm"), u); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
---
	add := func(col string, v any) {
		sets = append(sets, fmt.Sprintf("%s = $%d", col, len(args)+1))
		args = append(args, v)
	}
	if p.PatchJWTPublicKeyPEM {
		add("jwt_public_keys", []string{p.SetJWTPublicKeyPEM})
	}
	if p.ClearRegistrationLimit {
		add("device_registration_limit", nil)
	} else if p.PatchRegistrationLimit {
		add("device_registration_limit", p.SetRegistrationLimit)
	}
	if p.ClearRetention || (p.PatchRetention && p.SetRetention <= 0) {
		add("datastream_maximum_storage_retention", nil)
	} else if p.PatchRetention {
		add("datastream_maximum_storage_retention", p.SetRetention)
	}
	if len(sets) == 0 {
		return nil
	}

---
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
	} else if *retention == 0 {
		// 0 ≡ unset, upstream parity (measured Astarte v1.2.0): the engine
		// folds retention 0 to nil at the DB layer ("ScyllaDB considers
		// TTL=0 as unset"), PATCH maps 0 to ClearRetention, and the store's
		// retention sweep only treats values > 0 as ceilings. An explicit 0
		// is still explicit: it must not pick up the configured default
		// (#73), so fold before the nil-injection above.
		retention = nil
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

$ cat /root/astrate-mule/docs/Makefile; echo ---; ls /root/astrate-mule/docs/handoff 2>/dev/null | head -30
.PHONY: serve build clean install sync

install:
	pip install -r requirements.txt

sync:
	cp DESIGN.md site/DESIGN.md
	cp ROADMAP.md site/ROADMAP.md
	cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
	rm -rf site/api && cp -r api site/api
	rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui

serve: sync
	mkdocs serve -f mkdocs.yml

build: sync
	mkdocs build -f mkdocs.yml

clean:
	rm -rf ../site-dist
---
README.md
flow-design-a-named-durable-flows-2026-07-29.md
flow-design-b-container-block-2026-07-29.md
flow-lifecycle-followups-design-2026-08-22.md
flow-parity-audit-2026-07-29.md
flow-v2-decisions-2026-07-29.md
phase-1-handoff.md
phase-1-memory.md
phase-10-handoff.md
phase-10-memory.md
phase-11-handoff.md
phase-2-handoff.md
phase-2-memory.md
phase-3-handoff.md
phase-3-memory.md
phase-4-handoff.md
phase-4-memory.md
phase-5-handoff.md
phase-5-memory.md
phase-6-handoff.md
phase-6-memory.md
phase-7-memory.md
phase-8-handoff.md
phase-8-memory.md
phase-9-handoff.md
phase-9-memory.md
pokemon-agent-handoff.md
pokemon-agent-memory.md
session-2026-07-29-flow-factory-handoff.md
session-2026-07-29-handoff.md

$ cd /root/astrate-mule && rg -l "docs/api|astarte_housekeeping_api|site-dist" --glob '!site-dist/**' | head -30; echo ---; which mkdocs python3; echo ---; rg -n "retention" docs/COMPATIBILITY.md | head -20
docs/swagger-ui/index.html
docs/handoff/phase-1-handoff.md
docs/handoff/phase-9-memory.md
docs/handoff/phase-5-handoff.md
docs/handoff/phase-7-memory.md
docs/handoff/phase-2-handoff.md
docs/site/quickstart.md
docs/handoff/phase-4-memory.md
docs/handoff/phase-6-memory.md
docs/handoff/phase-10-memory.md
docs/handoff/phase-3-handoff.md
docs/handoff/phase-4-handoff.md
docs/mkdocs.yml
docs/handoff/phase-3-memory.md
docs/Makefile
docs/handoff/phase-2-memory.md
docs/handoff/phase-5-memory.md
docs/handoff/phase-6-handoff.md
docs/handoff/phase-1-memory.md
docs/handoff/phase-8-memory.md
internal/swagger/swagger_test.go
---
/usr/bin/python3
---

→ Read internal/swagger/swagger_test.go

$ cd /root/astrate-mule && git status --porcelain | head; echo ---; git check-ignore -v docs/site/api docs/site/swagger-ui site-dist 2>&1; echo ---; cat .gitignore 2>/dev/null | head -30
---
.gitignore:34:/docs/site/api/	docs/site/api
.gitignore:35:/docs/site/swagger-ui/	docs/site/swagger-ui
.gitignore:29:/site-dist/	site-dist
---
# Build artifacts
dist/
/astrate
*.exe
*.test
*.out

# Coverage
coverage.*
*.cover

# bbolt session stores (broker persistence, §3.1)
*.db

# Local env / editor
.env
.envrc
.idea/
.vscode/
*.swp
.DS_Store

# Local AI documentation workflow. The handoff files under docs/handoff/ ARE tracked:
# they are the project's phase history, and the agents that need it most are the cheap ones
# that cannot re-derive it.
/docs/AI-DOCUMENTATION-WORKFLOW.md

# MkDocs build output
/site-dist/

$ cd /root/astrate-mule && rg -n "retention|deviation|upstream" docs/api/astarte_housekeeping_api.yaml; echo === COMPAT; rg -n "retention|housekeeping" docs/COMPATIBILITY.md | head -30
9:    This API surface mirrors the upstream Astarte Housekeeping API paths and
148:        `device_registration_limit` or `datastream_maximum_storage_retention`
172:                datastream_maximum_storage_retention: 3600
243:        response is identical either way (deviation 17 in COMPATIBILITY.md).
283:        datastream_maximum_storage_retention:
287:          description: Maximum datastream storage retention in seconds. Null means unlimited.
304:        datastream_maximum_storage_retention:
308:          description: Maximum datastream storage retention in seconds. Null means unlimited.
322:        datastream_maximum_storage_retention:
326:          description: Maximum datastream storage retention in seconds. Null clears the retention.
=== COMPAT

$ cd /root/astrate-mule && rg -n "Retention|retention" internal/housekeeping/*_test.go | head -40
internal/housekeeping/http_test.go:41:	defaultRetention *int64
internal/housekeeping/http_test.go:50:func newHKRigWithDefault(t *testing.T, defaultRetention *int64) *hkRig {
internal/housekeeping/http_test.go:52:	return newHKRigOpts(t, hkRigOpts{defaultRetention: defaultRetention})
internal/housekeeping/http_test.go:89:		WithDefaultDatastreamMaximumStorageRetention(opts.defaultRetention).
internal/housekeeping/http_test.go:232:	t.Run("PatchRetention", func(t *testing.T) {
internal/housekeeping/http_test.go:233:		if rec := patch(t, `{"datastream_maximum_storage_retention":86400}`); rec.Code != http.StatusOK || rec.Body.Len() != 0 {
internal/housekeeping/http_test.go:234:			t.Fatalf("patch retention: got %d %q, want 200 empty", rec.Code, rec.Body)
internal/housekeeping/http_test.go:236:		if body := getBody(t); body.DatastreamMaximumStorageRetention == nil || *body.DatastreamMaximumStorageRetention != 86400 {
internal/housekeeping/http_test.go:237:			t.Errorf("retention after patch = %v, want 86400", body.DatastreamMaximumStorageRetention)
internal/housekeeping/http_test.go:263:		if rec := patch(t, `{"device_registration_limit":null,"datastream_maximum_storage_retention":null}`); rec.Code != http.StatusOK || rec.Body.Len() != 0 {
internal/housekeeping/http_test.go:267:		if body.DeviceRegistrationLimit != nil || body.DatastreamMaximumStorageRetention != nil {
internal/housekeeping/http_test.go:268:			t.Errorf("fields after null patch = %v/%v, want nil/nil", body.DeviceRegistrationLimit, body.DatastreamMaximumStorageRetention)
internal/housekeeping/http_test.go:272:	t.Run("ZeroRetentionUnsets", func(t *testing.T) {
internal/housekeeping/http_test.go:273:		if rec := patch(t, `{"datastream_maximum_storage_retention":60}`); rec.Code != http.StatusOK {
internal/housekeeping/http_test.go:274:			t.Fatalf("set retention: got %d (%s)", rec.Code, rec.Body)
internal/housekeeping/http_test.go:276:		if rec := patch(t, `{"datastream_maximum_storage_retention":0}`); rec.Code != http.StatusOK || rec.Body.Len() != 0 {
internal/housekeeping/http_test.go:277:			t.Fatalf("zero-set retention: got %d %q, want 200 empty", rec.Code, rec.Body)
internal/housekeeping/http_test.go:279:		if body := getBody(t); body.DatastreamMaximumStorageRetention != nil {
internal/housekeeping/http_test.go:280:			t.Errorf("retention after zero patch = %v, want nil", *body.DatastreamMaximumStorageRetention)
internal/housekeeping/http_test.go:285:		wantFieldErrors(t, `{"datastream_maximum_storage_retention":-5}`,
internal/housekeeping/http_test.go:286:			`{"errors":{"datastream_maximum_storage_retention":["is invalid"]}}`)
internal/housekeeping/http_test.go:307:// TestHousekeepingRetentionDefault covers #73: HOUSEKEEPING_DEFAULT_...
internal/housekeeping/http_test.go:309:func TestHousekeepingRetentionDefault(t *testing.T) {
internal/housekeeping/http_test.go:313:	create := func(t *testing.T, name string, retentionJSON string) {
internal/housekeeping/http_test.go:315:		body := `{"realm_name":` + jsonStr(name) + `,"jwt_public_key_pem":` + jsonStr(r.realmKey) + retentionJSON + `}`
internal/housekeeping/http_test.go:320:	retentionOf := func(t *testing.T, name string) *int64 {
internal/housekeeping/http_test.go:324:		return body.DatastreamMaximumStorageRetention
internal/housekeeping/http_test.go:329:		create(t, name, `,"datastream_maximum_storage_retention":60`)
internal/housekeeping/http_test.go:330:		if got := retentionOf(t, name); got == nil || *got != 60 {
internal/housekeeping/http_test.go:331:			t.Errorf("explicit retention = %v, want 60", got)
internal/housekeeping/http_test.go:337:		if got := retentionOf(t, name); got == nil || *got != def {
internal/housekeeping/http_test.go:338:			t.Errorf("injected retention = %v, want %d", got, def)
internal/housekeeping/http_test.go:350:		if rb.DatastreamMaximumStorageRetention != nil {
internal/housekeeping/http_test.go:351:			t.Errorf("retention without default = %v, want nil", *rb.DatastreamMaximumStorageRetention)
internal/housekeeping/service_test.go:40:		DatastreamMaximumStorageRetention: nr.DatastreamMaximumStorageRetention,
internal/housekeeping/service_test.go:93:// retention all fail before the store or sealer is touched.
internal/housekeeping/service_test.go:97:	negRetention := int64(-1)
internal/housekeeping/service_test.go:102:		retention               *int64
internal/housekeeping/service_test.go:107:		{"negative retention", "plant", "key", nil, &negRetention},
internal/housekeeping/service_test.go:114:			if _, err := svc.CreateRealm(ctx, tc.realmName, tc.jwtKey, tc.regLimit, tc.retention); !errors.Is(err, ErrValidation) {

→ Read internal/housekeeping/http_test.go [offset=300, limit=60]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=100]

$ cd /root/astrate-mule && sed -n 1,40p internal/housekeeping/http_test.go && sed -n 85,100p internal/housekeeping/http_test.go
//go:build integration

package housekeeping

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/pairing"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

type hkRig struct {
	st       *store.Store
	sealer   *store.KeySealer
	mux      *http.ServeMux
	haToken  string // valid a_ha
	wrongTok string // valid token, no a_ha
	realmKey string // a JWT public key PEM realms are created with
}

// hkRigOpts selects the optional service knobs a rig wires.
type hkRigOpts struct {
	instancePub := pubPEM(t, &instanceKey.PublicKey)

	mux := http.NewServeMux()
	NewAPI(NewService(st, sealer, nil, discardLogger()).
		WithDefaultDatastreamMaximumStorageRetention(opts.defaultRetention).
		WithRealmDeletionDisabled(opts.deletionDisabled),
		auth.NewMiddleware(st), []string{instancePub}).Mount(mux)

	return &hkRig{
		st: st, sealer: sealer, mux: mux,
		realmKey: pubPEM(t, &realmKey.PublicKey),
		haToken:  mintToken(t, instanceKey, jwt.MapClaims{"a_ha": []string{".*::.*"}}),
		wrongTok: mintToken(t, instanceKey, jwt.MapClaims{"a_aea": []string{".*::.*"}}),
	}
}

$ cd /root/astrate-mule && rg -n "HOUSEKEEPING_DEFAULT" --glob '!site-dist/**' -g '!docs/handoff/**' | head; echo === LIMIT; rg -rn "" -e "" /dev/null 2>/dev/null; rg -n "DeviceRegistrationLimit|device_registration_limit" internal --glob '!*_test.go' | head -20
internal/config/config.go:279:		"ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
internal/config/config.go:280:		"HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
internal/config/config_test.go:392:		prefixed = "ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION"
internal/config/config_test.go:393:		bare     = "HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION"
internal/housekeeping/http_test.go:307:// TestHousekeepingRetentionDefault covers #73: HOUSEKEEPING_DEFAULT_...
=== LIMIT
internal/config/config.go:146:	DeviceRegistrationLimit *int32 `toml:"device_registration_limit"`
internal/config/config.example.toml:94:# device_registration_limit = 0
internal/store/devices.go:186:// device_registration_limit in the pairing layer).
internal/housekeeping/http.go:49:	DeviceRegistrationLimit           *int32 `json:"device_registration_limit"`
internal/housekeeping/http.go:68:	rv, err := a.svc.CreateRealm(r.Context(), req.RealmName, req.JWTPublicKeyPEM, req.DeviceRegistrationLimit, req.DatastreamMaximumStorageRetention)
internal/housekeeping/http.go:98:		DeviceRegistrationLimit:           rv.DeviceRegistrationLimit,
internal/housekeeping/http.go:151:	DeviceRegistrationLimit           optInt32  `json:"device_registration_limit"`
internal/housekeeping/http.go:159:	"device_registration_limit":            {},
internal/housekeeping/http.go:192:	if pb.DeviceRegistrationLimit.present && !pb.DeviceRegistrationLimit.null && pb.DeviceRegistrationLimit.val < 0 {
internal/housekeeping/http.go:194:			map[string][]string{"device_registration_limit": {"is invalid"}})
internal/housekeeping/http.go:208:	if pb.DeviceRegistrationLimit.present {
internal/housekeeping/http.go:209:		if pb.DeviceRegistrationLimit.null {
internal/housekeeping/http.go:213:			u.SetRegistrationLimit = pb.DeviceRegistrationLimit.val
internal/housekeeping/service.go:39:	DeviceRegistrationLimit           *int32
internal/housekeeping/service.go:144:		return nil, fmt.Errorf("%w: device_registration_limit must be non-negative", ErrValidation)
internal/housekeeping/service.go:175:		DeviceRegistrationLimit:           regLimit,
internal/housekeeping/service.go:209:		return fmt.Errorf("%w: device_registration_limit must be non-negative", ErrValidation)
internal/housekeeping/service.go:285:		DeviceRegistrationLimit:           r.DeviceRegistrationLimit,
internal/store/realms.go:22:	DeviceRegistrationLimit           *int32
internal/store/realms.go:34:	DeviceRegistrationLimit           *int32

$ cd /root/astrate-mule && sed -n 175,200p internal/store/devices.go
return nil, fmt.Errorf("store: scanning device: %w", err)
		}
		out = append(out, *d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: listing devices: %w", err)
	}
	return out, nil
}

// CountDevices returns the number of devices in a realm (used to enforce
// device_registration_limit in the pairing layer).
func (s *Store) CountDevices(ctx context.Context, realmID int16) (int64, error) {
	var n int64
	if err := s.pool.QueryRow(ctx,
		`SELECT count(*) FROM devices WHERE realm_id = $1`, realmID).Scan(&n); err != nil {
		return 0, fmt.Errorf("store: counting devices: %w", err)
	}
	return n, nil
}

// UpdateIntrospection replaces the device's introspection. Every
// (name, major) pair present in the current introspection but missing from
// the new one is merged into old_introspection (docs/ROADMAP.md §3.1 file
// 2.10) and returned, so the engine can react (cache eviction, property
// cleanup).

$ cd /root/astrate-mule && rg -n "CountDevices" internal --glob '!*_test.go' -B4 -A8 | head -40
internal/store/devices.go-181-	}
internal/store/devices.go-182-	return out, nil
internal/store/devices.go-183-}
internal/store/devices.go-184-
internal/store/devices.go:185:// CountDevices returns the number of devices in a realm (used to enforce
internal/store/devices.go-186-// device_registration_limit in the pairing layer).
internal/store/devices.go:187:func (s *Store) CountDevices(ctx context.Context, realmID int16) (int64, error) {
internal/store/devices.go-188-	var n int64
internal/store/devices.go-189-	if err := s.pool.QueryRow(ctx,
internal/store/devices.go-190-		`SELECT count(*) FROM devices WHERE realm_id = $1`, realmID).Scan(&n); err != nil {
internal/store/devices.go-191-		return 0, fmt.Errorf("store: counting devices: %w", err)
internal/store/devices.go-192-	}
internal/store/devices.go-193-	return n, nil
internal/store/devices.go-194-}
internal/store/devices.go-195-
--
internal/pairing/service.go-31-	UnregisterDevice(ctx context.Context, realmID int16, id deviceid.ID) error
internal/pairing/service.go-32-	GetDevice(ctx context.Context, realmID int16, id deviceid.ID) (*store.Device, error)
internal/pairing/service.go-33-	SetDeviceCredentials(ctx context.Context, realmID int16, id deviceid.ID, certSerial, certAKI string, requestIP netip.Addr) error
internal/pairing/service.go-34-	SetPayloadFormatHint(ctx context.Context, realmID int16, id deviceid.ID, hint string) error
internal/pairing/service.go:35:	CountDevices(ctx context.Context, realmID int16) (int64, error)
internal/pairing/service.go-36-	// Health verifies database liveness behind the realm-scoped health
internal/pairing/service.go-37-	// probe (#71).
internal/pairing/service.go-38-	Health(ctx context.Context) error
internal/pairing/service.go-39-}
internal/pairing/service.go-40-
internal/pairing/service.go-41-// Service-level sentinel errors; the HTTP layer maps them onto upstream
internal/pairing/service.go-42-// statuses and envelopes.
internal/pairing/service.go-43-var (
--
internal/pairing/service.go-187-
internal/pairing/service.go-188-	// Upstream parity: the limit gates every registration attempt, secret
internal/pairing/service.go-189-	// rotations for existing devices included.
internal/pairing/service.go-190-	if realm.DeviceRegistrationLimit != nil {
internal/pairing/service.go:191:		n, err := s.st.CountDevices(ctx, realm.ID)
internal/pairing/service.go-192-		if err != nil {
internal/pairing/service.go-193-			return "", err
internal/pairing/service.go-194-		}
internal/pairing/service.go-195-		if n >= int64(*realm.DeviceRegistrationLimit) {
internal/pairing/service.go-196-			return "", ErrRegistrationLimitReached

Now the yaml edits.
← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -145,9 +145,14 @@
 description: |
   Updates one or more fields of an existing realm. Anything not named here
   is rejected with `invalid_update_parameters`. A null
   `device_registration_limit` or `datastream_maximum_storage_retention`
-  clears that field. Returns 200 with an empty body on success.
+  clears that field. For `datastream_maximum_storage_retention`, an
+  explicit `0` clears it too: 0 is folded to unset, so it cannot ask for
+  a literal zero-second retention, and the realm is left unlimited
+  (upstream parity, measured on Astarte v1.2.0). `device_registration_limit`
+  has no such fold — `0` there is stored as the literal limit `0`, and
+  only `null` clears it. Returns 200 with an empty body on success.
 tags: [Housekeeping]
 security:
   - a_ha: []
 parameters:

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -304,14 +304,22 @@
     device_registration_limit:
       type: integer
       format: int32
       nullable: true
-      description: Maximum number of devices allowed in this realm. Null means unlimited.
+      description: |
+        Maximum number of devices allowed in this realm. Null means unlimited.
+        An explicit `0` is kept literally and is not the same as null: the
+        limit is then reached before the first device, so nothing can
+        register. Only `null` clears the limit.
     datastream_maximum_storage_retention:
       type: integer
       format: int64
       nullable: true
-      description: Maximum datastream storage retention in seconds. Null means unlimited.
+      description: |
+        Maximum datastream storage retention in seconds. Null means unlimited.
+        An explicit `0` is folded to unset and means unlimited, not a literal
+        zero-second retention; unlike omitting the field it does not pick up
+        the deployment's configured default retention.
 
 RealmPatch:
   type: object
   properties:

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -330,14 +330,21 @@
     device_registration_limit:
       type: integer
       format: int32
       nullable: true
-      description: Maximum number of devices allowed in this realm. Null clears the limit (unlimited).
+      description: |
+        Maximum number of devices allowed in this realm. Null means unlimited.
+        An explicit `0` is stored literally and admits no device, so it is
+        not the same as null: only `null` means unlimited.
     datastream_maximum_storage_retention:
       type: integer
       format: int64
       nullable: true
-      description: Maximum datastream storage retention in seconds. Null clears the retention.
+      description: |
+        Maximum datastream storage retention in seconds. Null clears the
+        retention. `0` clears it as well: 0 is folded to unset, so it cannot
+        ask for a literal zero-second retention and the realm is left
+        unlimited (upstream parity, measured on Astarte v1.2.0).
 
 # ── Error envelopes ──
 
 ErrorDetail:

Now a doc test pinning the new wording, following the existing pattern in that file.
← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -208,8 +208,96 @@
 		t.Error("components parameter AsyncOperation does not say the value is accepted and ignored")
 	}
 }
 
+// propertyBlock returns the lines of the schema property with the given name
+// inside a components.schemas entry, up to the next property or the end of the
+// entry.
+func propertyBlock(t *testing.T, lines []string, name string) []string {
+	t.Helper()
+	key := "        " + name + ":"
+	start := -1
+	for i, l := range lines {
+		if l == key {
+			start = i
+			break
+		}
+	}
+	if start < 0 {
+		t.Fatalf("schema declares no property %q", name)
+	}
+	for i := start + 1; i < len(lines); i++ {
+		if strings.TrimSpace(lines[i]) == "" {
+			continue
+		}
+		if !strings.HasPrefix(lines[i], "        ") {
+			return lines[start+1 : i]
+		}
+	}
+	return lines[start+1:]
+}
+
+// TestHousekeepingRetentionZeroFoldDocumented guards that the housekeeping spec
+// tells clients that `datastream_maximum_storage_retention: 0` means unset, not
+// a literal zero-second retention. The wire folds it on both paths: PATCH maps
+// `null || val == 0` to ClearRetention (internal/housekeeping/http.go) and
+// create folds 0 to nil before injecting the configured default
+// (internal/housekeeping/service.go) — upstream parity measured on v1.2.0. A
+// client that follows the spec without this note sends 0 and silently gets
+// unlimited instead. `device_registration_limit` has no such fold, so the same
+// 0 is stored literally there; the asymmetry is documented on both fields so
+// the two do not read as interchangeable. This is the documentation half of the
+// behaviour pinned by ZeroRetentionUnsets and TestHousekeepingRetentionDefault
+// in internal/housekeeping.
+func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	patch := strings.Join(operationBlock(t, lines, "patchRealm"), "\n")
+	for _, want := range []string{
+		"explicit `0` clears it too",
+		"0 is folded to unset",
+		"has no such fold",
+		"stored as the literal limit `0`",
+	} {
+		if !strings.Contains(patch, want) {
+			t.Errorf("patchRealm description does not say %q", want)
+		}
+	}
+
+	for _, schema := range []string{"    RealmCreate:", "    RealmPatch:"} {
+		block := componentBlock(t, lines, schema)
+
+		retention := strings.Join(propertyBlock(t, block, "datastream_maximum_storage_retention"), "\n")
+		if !strings.Contains(retention, "0") || !strings.Contains(retention, "unlimited") {
+			t.Errorf("%s retention does not say what 0 means", schema)
+		}
+		if !strings.Contains(retention, "unset") {
+			t.Errorf("%s retention does not say 0 is folded to unset", schema)
+		}
+
+		limit := strings.Join(propertyBlock(t, block, "device_registration_limit"), "\n")
+		if !strings.Contains(limit, "`0`") {
+			t.Errorf("%s device_registration_limit does not record what 0 means", schema)
+		}
+		if strings.Contains(limit, "folded") {
+			t.Errorf("%s device_registration_limit claims a 0 fold that the wire does not have", schema)
+		}
+		if !strings.Contains(limit, "null") {
+			t.Errorf("%s device_registration_limit does not say null clears the limit", schema)
+		}
+	}
+
+	create := strings.Join(propertyBlock(t,
+		componentBlock(t, lines, "    RealmCreate:"), "datastream_maximum_storage_retention"), "\n")
+	if !strings.Contains(create, "configured default retention") {
+		t.Error("RealmCreate retention does not say 0 skips the configured default retention")
+	}
+}
+
 // operationBlock returns the lines of the operation with the given operationId,
 // up to the next operation, path, or the components section.
 func operationBlock(t *testing.T, lines []string, operationID string) []string {
 	t.Helper()

$ cd /root/astrate-mule && gofmt -l . && go vet ./internal/swagger/ && go test ./internal/swagger/ -run 'TestHousekeeping|TestSpecs|TestMount' -v 2>&1 | tail -30
=== RUN   TestMount
=== RUN   TestMount/GET_/swagger_redirects_to_/swagger/index.html
=== RUN   TestMount/GET_/swagger/_serves_the_embedded_UI
=== RUN   TestMount/GET_/api/_serves_every_OpenAPI_YAML_spec
--- PASS: TestMount (0.05s)
    --- PASS: TestMount/GET_/swagger_redirects_to_/swagger/index.html (0.01s)
    --- PASS: TestMount/GET_/swagger/_serves_the_embedded_UI (0.02s)
    --- PASS: TestMount/GET_/api/_serves_every_OpenAPI_YAML_spec (0.01s)
=== RUN   TestMountSubPanicsOnBrokenFS
--- PASS: TestMountSubPanicsOnBrokenFS (0.00s)
=== RUN   TestSpecs
--- PASS: TestSpecs (0.00s)
=== RUN   TestHousekeepingAsyncOperationParamDocumented
--- PASS: TestHousekeepingAsyncOperationParamDocumented (0.01s)
=== RUN   TestHousekeepingRetentionZeroFoldDocumented
    swagger_test.go:286:     RealmCreate: device_registration_limit claims a 0 fold that the wire does not have
    swagger_test.go:286:     RealmPatch: device_registration_limit claims a 0 fold that the wire does not have
--- FAIL: TestHousekeepingRetentionZeroFoldDocumented (0.00s)
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.092s
FAIL

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -228,9 +228,9 @@
 for i := start + 1; i < len(lines); i++ {
 	if strings.TrimSpace(lines[i]) == "" {
 		continue
 	}
-	if !strings.HasPrefix(lines[i], "        ") {
+	if !strings.HasPrefix(lines[i], "          ") {
 		return lines[start+1 : i]
 	}
 }
 return lines[start+1:]

$ cd /root/astrate-mule && gofmt -l . ; go test ./internal/swagger/ -run 'TestHousekeepingRetentionZeroFoldDocumented' -v 2>&1 | tail -20
=== RUN   TestHousekeepingRetentionZeroFoldDocumented
--- PASS: TestHousekeepingRetentionZeroFoldDocumented (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.026s

$ cd /root/astrate-mule && cp docs/api/astarte_housekeeping_api.yaml /tmp/opencode/hk.yaml.new && git show HEAD:docs/api/astarte_housekeeping_api.yaml > docs/api/astarte_housekeeping_api.yaml && go test ./internal/swagger/ -run 'TestHousekeepingRetentionZeroFoldDocumented' 2>&1 | tail -12; cp /tmp/opencode/hk.yaml.new docs/api/astarte_housekeeping_api.yaml && git diff --stat
swagger_test.go:266: patchRealm description does not say "has no such fold"
    swagger_test.go:266: patchRealm description does not say "stored as the literal limit `0`"
    swagger_test.go:275:     RealmCreate: retention does not say what 0 means
    swagger_test.go:278:     RealmCreate: retention does not say 0 is folded to unset
    swagger_test.go:283:     RealmCreate: device_registration_limit does not record what 0 means
    swagger_test.go:275:     RealmPatch: retention does not say what 0 means
    swagger_test.go:278:     RealmPatch: retention does not say 0 is folded to unset
    swagger_test.go:283:     RealmPatch: device_registration_limit does not record what 0 means
    swagger_test.go:296: RealmCreate retention does not say 0 skips the configured default retention
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.027s
FAIL
 docs/api/astarte_housekeeping_api.yaml | 30 ++++++++++--
 internal/swagger/swagger_test.go       | 88 ++++++++++++++++++++++++++++++++++
 2 files changed, 113 insertions(+), 5 deletions(-)

$ cd /root/astrate-mule && python3 -c "import yaml,sys; d=yaml.safe_load(open('docs/api/astarte_housekeeping_api.yaml')); print('yaml ok, paths:', list(d['paths'])); print('schemas:', list(d['components']['schemas']))" ; uname -m; python3 -m pip --version 2>&1 | head -2
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml,sys; d=yaml.safe_load(open('docs/api/astarte_housekeeping_api.yaml')); print('yaml ok, paths:', list(d['paths'])); print('schemas:', list(d['components']['schemas']))
    ^^^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'
aarch64
/usr/bin/python3: No module named pip

$ cd /root/astrate-mule && go vet ./... 2>&1 | tail -5; echo "VET=$?"; go test ./... 2>&1 | rg -v "^ok|no test files" | head -20; echo "TEST_DONE"
VET=0
TEST_DONE

$ cd /root/astrate-mule && rg -n "yaml" go.mod; ls /root/go/pkg/mod/gopkg.in 2>/dev/null; ls /root/go/pkg/mod/github.com 2>/dev/null | head -5; echo "=== make"; timeout 60 make -C docs build 2>&1 | tail -5
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect
check.v1@v0.0.0-20161208181325-20d25e280405
check.v1@v1.0.0-20201130134442-10cb98267c6c
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1
!abirdcfly
!ada!logics
!admin!benni
!alwx!sin
!antonboom
=== make
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'

$ cd /root/astrate-mule && git status --porcelain; echo "=== net"; timeout 6 getent hosts pypi.org || echo "no dns"; python3 -c "import ensurepip; print('ensurepip ok')" 2>&1 | tail -1
M docs/api/astarte_housekeeping_api.yaml
 M internal/swagger/swagger_test.go
=== net
2a04:4e42::223  pypi.org
2a04:4e42:200::223 pypi.org
2a04:4e42:400::223 pypi.org
2a04:4e42:600::223 pypi.org
ensurepip ok

$ cd /root/astrate-mule && cat docs/requirements.txt; python3 -m venv /tmp/opencode/mkv >/dev/null 2>&1 && echo venv-ok; timeout 240 /tmp/opencode/mkv/bin/pip install -q -r docs/requirements.txt 2>&1 | tail -5; echo "PIP_EXIT=$?"
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10
venv-ok
PIP_EXIT=0

$ cd /root/astrate-mule && /tmp/opencode/mkv/bin/mkdocs --version && PATH="/tmp/opencode/mkv/bin:$PATH" timeout 300 make -C docs build 2>&1 | tail -25
mkdocs, version 1.6.1 from /tmp/opencode/mkv/lib/python3.13/site-packages/mkdocs (Python 3.13)
mkdocs build -f mkdocs.yml

 │  ⚠  Warning from the Material for MkDocs team
 │
 │  MkDocs 2.0, the underlying framework of Material for MkDocs,
 │  will introduce backward-incompatible changes, including:
 │
 │  × All plugins will stop working – the plugin system has been removed
 │  × All theme overrides will break – the theming system has been rewritten
 │  × No migration path exists – existing projects cannot be upgraded
 │  × Closed contribution model – community members can't report bugs
 │  × Currently unlicensed – unsuitable for production use
 │
 │  Our full analysis:
 │
 │  https://squidfunk.github.io/mkdocs-material/blog/2026/02/18/mkdocs-2.0/

INFO    -  Cleaning site directory
INFO    -  Building documentation to directory: /root/astrate-mule/site-dist
INFO    -  The following pages exist in the docs directory, but are not included in the "nav" configuration:
  - DESIGN.md
  - ROADMAP.md
INFO    -  Doc file 'index.md' contains an unrecognized relative link 'api/', it was left as is.
INFO    -  Documentation built in 9.30 seconds
make: Leaving directory '/root/astrate-mule/docs'

$ cd /root/astrate-mule && rg -c "folded to unset" site-dist/api/astarte_housekeeping_api.yaml; ls site-dist/api/; echo "=== race"; timeout 300 go test -race ./internal/swagger/ 2>&1 | tail -5
3
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
=== race
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.065s
FAIL

$ cd /root/astrate-mule && gofmt -l . ; echo "GOFMT_CLEAN"; go vet ./... && echo VET_OK; go test ./... 2>&1 | rg -v "no test files" | tail -25; git status --porcelain
GOFMT_CLEAN
VET_OK
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
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)
 M docs/api/astarte_housekeeping_api.yaml
 M internal/swagger/swagger_test.go

$ cd /root/astrate-mule && git diff docs/api/astarte_housekeeping_api.yaml
diff --git a/docs/api/astarte_housekeeping_api.yaml b/docs/api/astarte_housekeeping_api.yaml
index bea5e74..73fc0c5 100644
--- a/docs/api/astarte_housekeeping_api.yaml
+++ b/docs/api/astarte_housekeeping_api.yaml
@@ -146,7 +146,12 @@ paths:
         Updates one or more fields of an existing realm. Anything not named here
         is rejected with `invalid_update_parameters`. A null
         `device_registration_limit` or `datastream_maximum_storage_retention`
-        clears that field. Returns 200 with an empty body on success.
+        clears that field. For `datastream_maximum_storage_retention`, an
+        explicit `0` clears it too: 0 is folded to unset, so it cannot ask for
+        a literal zero-second retention, and the realm is left unlimited
+        (upstream parity, measured on Astarte v1.2.0). `device_registration_limit`
+        has no such fold — `0` there is stored as the literal limit `0`, and
+        only `null` clears it. Returns 200 with an empty body on success.
       tags: [Housekeeping]
       security:
         - a_ha: []
@@ -300,12 +305,20 @@ components:
           type: integer
           format: int32
           nullable: true
-          description: Maximum number of devices allowed in this realm. Null means unlimited.
+          description: |
+            Maximum number of devices allowed in this realm. Null means unlimited.
+            An explicit `0` is kept literally and is not the same as null: the
+            limit is then reached before the first device, so nothing can
+            register. Only `null` clears the limit.
         datastream_maximum_storage_retention:
           type: integer
           format: int64
           nullable: true
-          description: Maximum datastream storage retention in seconds. Null means unlimited.
+          description: |
+            Maximum datastream storage retention in seconds. Null means unlimited.
+            An explicit `0` is folded to unset and means unlimited, not a literal
+            zero-second retention; unlike omitting the field it does not pick up
+            the deployment's configured default retention.
 
     RealmPatch:
       type: object
@@ -318,12 +331,19 @@ components:
           type: integer
           format: int32
           nullable: true
-          description: Maximum number of devices allowed in this realm. Null clears the limit (unlimited).
+          description: |
+            Maximum number of devices allowed in this realm. Null means unlimited.
+            An explicit `0` is stored literally and admits no device, so it is
+            not the same as null: only `null` means unlimited.
         datastream_maximum_storage_retention:
           type: integer
           format: int64
           nullable: true
-          description: Maximum datastream storage retention in seconds. Null clears the retention.
+          description: |
+            Maximum datastream storage retention in seconds. Null clears the
+            retention. `0` clears it as well: 0 is folded to unset, so it cannot
+            ask for a literal zero-second retention and the realm is left
+            unlimited (upstream parity, measured on Astarte v1.2.0).
 
     # ── Error envelopes ──

Done: documented the `0 ≡ unset` retention fold (and the deliberate `device_registration_limit` asymmetry) in the housekeeping spec — `patchRealm` description, `RealmCreate.datastream_maximum_storage_retention`, `RealmCreate.device_registration_limit`, `RealmPatch.datastream_maximum_storage_retention`, `RealmPatch.device_registration_limit` — plus a new doc test pinning that wording.
Files: docs/api/astarte_housekeeping_api.yaml, internal/swagger/swagger_test.go
Verified: `go vet ./...` -> pass; `go test ./...` -> pass (all packages ok); `gofmt -l .` -> clean; `make -C docs build` (mkdocs 1.6.1 + material installed into a /tmp venv, since mkdocs was missing) -> pass, site-dist/api copy carries the new text; new test confirmed to fail with the yaml reverted to HEAD (10 assertion failures) and pass with it. `go test -race` unavailable: `ThreadSanitizer: unsupported VMA range` on this 39-bit-VMA box.
Unsure: the create-side claim (explicit `0` stored as nil and *not* replaced by the configured default) is read from `internal/housekeeping/service.go:149-159`, not executed — the pinning tests there are `//go:build integration` and need a DB, so the Legion Go would be needed to confirm on the wire. I did not run anything on the Legion. Wording choices ("admits no device" for a literal `device_registration_limit: 0`) are derived from `internal/pairing/service.go:190-196`, verified by reading, not by running.
Follow-ups: none
