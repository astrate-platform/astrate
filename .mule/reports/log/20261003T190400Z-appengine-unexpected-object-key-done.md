slug: appengine-unexpected-object-key
verdict: done
at:  f7d60af
ran: 2026-10-03T19:04:00Z on DietPi in 956s


> build · big-pickle

I'll start by reading the two files.
→ Read .mule/MULE.md
→ Read .mule/task.md
✱ Grep "ReasonBadObject|ReasonMissingRequired|Reason[A-Za-z]+" 100 matches

$ wc -l internal/appengine/http.go pkg/payload/*.go internal/engine/serverdata.go 2>/dev/null; ls internal/appengine/ pkg/payload/
696 internal/appengine/http.go
   343 pkg/payload/bson.go
   122 pkg/payload/fuzz_test.go
   408 pkg/payload/json.go
   217 pkg/payload/payload.go
   792 pkg/payload/payload_test.go
    71 pkg/payload/required_test.go
    41 pkg/payload/sniff.go
   120 pkg/payload/sniff_test.go
   353 pkg/payload/value.go
   223 pkg/payload/value_test.go
   349 internal/engine/serverdata.go
  3735 total
internal/appengine/:
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

pkg/payload/:
bson.go
fuzz_test.go
json.go
payload.go
payload_test.go
required_test.go
sniff.go
sniff_test.go
testdata
value.go
value_test.go

→ Read internal/appengine/writeerror_test.go
→ Read internal/appengine/http.go [offset=600, limit=96]

$ ls internal/astarteapi/ 2>/dev/null || rg -l "astarteapi" --type go | head -20
pkg/astarteapi/envelope.go
pkg/astarteapi/envelope_test.go
internal/housekeeping/http.go
internal/pairing/http.go
internal/flowapi/http.go
internal/auth/middleware.go
internal/httpx/notfound.go
internal/realm/http.go
internal/swagger/swagger_test.go
internal/appengine/http.go
internal/appengine/channels/ws.go

✱ Grep "func Write|func .*Detail" in internal/astarteapi · 18 matches
→ Read pkg/astarteapi/envelope.go
→ Read pkg/payload/value.go [limit=80, offset=100]
→ Read pkg/payload/bson.go [limit=50, offset=130]
→ Read pkg/payload/value.go [limit=30, offset=179]
→ Read pkg/payload/json.go [limit=60, offset=190]
→ Read internal/appengine/http.go [limit=80, offset=285]

$ rg -n "ReasonOf|RejectError" --type go -g '!*_test.go' | grep -v "^pkg/payload/"
internal/engine/data.go:289:		reason := payload.ReasonOf(err)
internal/engine/serverdata.go:18:// statuses; payload validation failures surface as *payload.RejectError.

$ rg -n "publishAsOwner|dec.Object|DecodeObject|annotate" internal/engine/serverdata.go internal/engine/data.go pkg/payload/payload.go | head -40
internal/engine/serverdata.go:51:	out, err := e.publishAsOwner(ctx, realm, id, ifaceName, path, value, ts, interfaceschema.OwnershipServer)
internal/engine/serverdata.go:73:	_, err := e.publishAsOwner(ctx, realm, id, ifaceName, path, value, ts, interfaceschema.OwnershipDevice)
internal/engine/serverdata.go:88:// publishAsOwner is the shared validate-and-persist core of the publish
internal/engine/serverdata.go:98:func (e *Engine) publishAsOwner(ctx context.Context, realm string, id deviceid.ID,
internal/engine/serverdata.go:139:		dp, err = dec.Object(envelope, ci.ObjectLeaves)
internal/engine/data.go:278:		dp, err = e.dec.Object(m.Payload, ci.ObjectLeaves)
pkg/payload/payload.go:185:// DecodeObject decodes an object-aggregation payload with the default
pkg/payload/payload.go:187:func DecodeObject(p []byte, leaves map[string]*interfaceschema.CompiledMapping) (DecodedPayload, error) {

→ Read internal/engine/data.go [limit=40, offset=275]
→ Read internal/engine/serverdata.go [limit=60, offset=125]

$ rg -n "putData|/data|newTestAPI|func newAPI|httptest.NewServer" internal/appengine/patch_test.go internal/appengine/e2e_test.go | head -40
internal/appengine/e2e_test.go:237:	srv := httptest.NewServer(mux)

→ Read internal/appengine/e2e_test.go [limit=120, offset=150]

$ head -12 internal/appengine/e2e_test.go internal/appengine/patch_test.go internal/appengine/http_test.go; rg -n "api\)|api\.putData|\.putData\(|Mount\(" internal/appengine/*_test.go | head -30
==> internal/appengine/e2e_test.go <==
//go:build integration && e2e

package appengine

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"

==> internal/appengine/patch_test.go <==
//go:build integration

package appengine

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

==> internal/appengine/http_test.go <==
//go:build integration

package appengine

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
internal/appengine/e2e_test.go:236:	channels.NewAPI(env.engine.Bus(), env.st, env.st).Mount(mux)
internal/appengine/http_test.go:149:	NewAPI(NewService(st, sd, discardLogger()), auth.NewMiddleware(st)).Mount(mux)

$ rg -n "type API struct" -A 20 internal/appengine/service.go | head -40; rg -ln "go:build" internal/appengine/*_test.go internal/engine/*_test.go | head; rg -c "go:build integration" internal/engine/*_test.go | head -30
internal/appengine/patch_test.go
internal/appengine/http_test.go
internal/appengine/groups_parity_test.go
internal/appengine/dashboard_compat_test.go
internal/engine/triggers_change_integration_test.go
internal/appengine/e2e_test.go
internal/engine/engine_e2e_test.go
internal/appengine/mirror_test.go
internal/appengine/formats_parity_test.go
internal/engine/engine_integration_test.go
internal/engine/engine_integration_test.go:1
internal/engine/triggers_change_integration_test.go:1

$ rg -n "type API" -A 18 internal/appengine/*.go | head -50
internal/appengine/http.go:28:type API struct {
internal/appengine/http.go-29-	svc     *Service
internal/appengine/http.go-30-	require func(http.Handler) http.Handler
internal/appengine/http.go-31-}
internal/appengine/http.go-32-
internal/appengine/http.go-33-// NewAPI wires the AppEngine service to its HTTP surface.
internal/appengine/http.go-34-func NewAPI(svc *Service, mw *auth.Middleware) *API {
internal/appengine/http.go-35-	return &API{svc: svc, require: mw.RequireRealm(auth.ClaimAppEngine)}
internal/appengine/http.go-36-}
internal/appengine/http.go-37-
internal/appengine/http.go-38-// Mount registers the routes on mux (paths wire-frozen to upstream
internal/appengine/http.go-39-// astarte_appengine_api).
internal/appengine/http.go-40-func (a *API) Mount(mux *http.ServeMux) {
internal/appengine/http.go-41-	h := func(f http.HandlerFunc) http.Handler { return a.require(f) }
internal/appengine/http.go-42-	const base = "/appengine/v1/{realm}"
internal/appengine/http.go-43-	mux.Handle("GET "+base+"/devices", h(a.listDevices))
internal/appengine/http.go-44-	mux.Handle("GET "+base+"/stats/devices", h(a.devicesStats))
internal/appengine/http.go-45-	mux.Handle("GET "+base+"/devices/{device}", h(a.getDevice))
internal/appengine/http.go-46-	mux.Handle("PATCH "+base+"/devices/{device}", h(a.patchDevice))

$ rg -ln "go:build" internal/engine/*_test.go | head -20; ls internal/engine/*_test.go
internal/engine/engine_integration_test.go
internal/engine/engine_e2e_test.go
internal/engine/triggers_change_integration_test.go
internal/engine/batch_test.go
internal/engine/bus_events_test.go
internal/engine/cache_test.go
internal/engine/capabilities_test.go
internal/engine/control_test.go
internal/engine/data_test.go
internal/engine/devicedata_test.go
internal/engine/engine_e2e_test.go
internal/engine/engine_integration_test.go
internal/engine/fakestore_test.go
internal/engine/introspection_test.go
internal/engine/mqttv1_parity_test.go
internal/engine/router_test.go
internal/engine/serverdata_test.go
internal/engine/topics_test.go
internal/engine/triggers_change_integration_test.go
internal/engine/triggers_change_test.go

$ rg -n "func NewService" -A 12 internal/appengine/service.go | head -25; rg -rn "func NewMem|func NewMemory|memory store|type MemStore" internal/store/*.go | head -20
86:func NewService(st *store.Store, sd ServerData, log *slog.Logger) *Service {
87-	if log == nil {
88-		log = slog.Default()
89-	}
90-	return &Service{st: st, sd: sd, log: log}
91-}
92-
93-func (s *Service) realmID(ctx context.Context, realm string) (int16, error) {
94-	r, err := s.st.GetRealmByName(ctx, realm)
95-	if err != nil {
96-		return 0, err
97-	}
98-	return r.ID, nil

$ rg -n "^func Test|^func new|^type |Object" internal/engine/serverdata_test.go | head -50
41:func TestPublishServerProperty(t *testing.T) {
81:func TestPublishServerPropertyJSONHint(t *testing.T) {
109:func TestPublishServerDatastream(t *testing.T) {
154:func TestPublishServerValueErrors(t *testing.T) {
213:func TestUnsetServerProperty(t *testing.T) {
268:// TestServerObjectAggregate: object-aggregated server interfaces accept one
270:func TestServerObjectAggregate(t *testing.T) {
273:		"interface_name": "com.astrate.test.ServerObject",
287:		"com.astrate.test.ServerObject", "/setpoints",
309:		"com.astrate.test.ServerObject", "/nope", json.RawMessage(`{"heating": 1.0}`), nil)
336:func TestServerPublishFiresBusAndTriggers(t *testing.T) {

→ Read internal/engine/serverdata_test.go [limit=45, offset=265]
→ Read internal/engine/serverdata_test.go [limit=70, offset=150]

$ rg -n "matches no declared object leaf|RejectError\{" --type go | head -20
pkg/payload/value.go:194:	return &RejectError{Reason: reason, Detail: fmt.Sprintf(format, args...)}
pkg/payload/value.go:317:		return &RejectError{Reason: re.Reason, Detail: fmt.Sprintf(format, args...) + ": " + re.Detail}
pkg/payload/json.go:76:			return time.Time{}, &RejectError{Reason: ReasonBadTimestamp, Detail: `"t": ` + re.Detail}
pkg/payload/json.go:216:			return nil, rejectf(ReasonBadObject, "key %q matches no declared object leaf", key)
pkg/payload/bson.go:157:			return nil, rejectf(ReasonBadObject, "key %q matches no declared object leaf", key)

→ Read pkg/payload/value.go [limit=45, offset=308]
→ Read pkg/payload/value.go [limit=30]
Now the payload-side change: carry the offending keys, sorted.

← Edit pkg/payload/value.go
Index: /root/astrate-mule/pkg/payload/value.go
===================================================================
--- /root/astrate-mule/pkg/payload/value.go
+++ /root/astrate-mule/pkg/payload/value.go
@@ -181,8 +181,14 @@
 	// Reason is the rejection class (metrics label).
 	Reason RejectReason
 	// Detail is a human-readable explanation for logs and trigger events.
 	Detail string
+	// UnexpectedKeys carries, sorted, the object-aggregation keys that
+	// resolve to no declared leaf (ReasonBadObject only; nil otherwise).
+	// Upstream reports the whole set as `unexpected_keys` beside the error
+	// detail (master b6d46ad4, #2237), so the HTTP surfaces name the keys
+	// instead of collapsing them into one log line.
+	UnexpectedKeys []string
 }
 
 // Error implements the error interface.
 func (e *RejectError) Error() string {

← Edit pkg/payload/value.go
Index: /root/astrate-mule/pkg/payload/value.go
===================================================================
--- /root/astrate-mule/pkg/payload/value.go
+++ /root/astrate-mule/pkg/payload/value.go
@@ -199,8 +199,30 @@
 func rejectf(reason RejectReason, format string, args ...any) error {
 	return &RejectError{Reason: reason, Detail: fmt.Sprintf(format, args...)}
 }
 
+// rejectUnexpectedKeys builds the object-shape rejection for the given
+// undeclared keys, reporting them sorted so the answer does not depend on map
+// or document order. Every offending key is named, not just the first one
+// found (master b6d46ad4, #2237).
+func rejectUnexpectedKeys(keys []string) error {
+	sorted := append([]string(nil), keys...)
+	sort.Strings(sorted)
+	return &RejectError{
+		Reason:         ReasonBadObject,
+		Detail:         fmt.Sprintf("object keys %s match no declared object leaf", clipKeys(sorted)),
+		UnexpectedKeys: sorted,
+	}
+}
+
+// clipKeys renders a key list for the Detail line, bounded like clip.
+func clipKeys(keys []string) string {
+	if len(keys) <= 8 {
+		return "[" + strings.Join(keys, " ") + "]"
+	}
+	return fmt.Sprintf("[%s ... +%d more]", strings.Join(keys[:8], " "), len(keys)-8)
+}
+
 // ReasonOf extracts the RejectReason from err, or ReasonNone if err is nil
 // or not a *RejectError.
 func ReasonOf(err error) RejectReason {
 	var re *RejectError

← Edit pkg/payload/value.go
Index: /root/astrate-mule/pkg/payload/value.go
===================================================================
--- /root/astrate-mule/pkg/payload/value.go
+++ /root/astrate-mule/pkg/payload/value.go
@@ -15,8 +15,10 @@
 import (
 	"errors"
 	"fmt"
 	"math"
+	"sort"
+	"strings"
 	"time"
 	"unicode/utf8"
 
 	"github.com/astrate-platform/astrate/pkg/interfaceschema"

$ rg -n "func clip" -A 12 pkg/payload/*.go
pkg/payload/value.go:220:func clipKeys(keys []string) string {
pkg/payload/value.go-221-	if len(keys) <= 8 {
pkg/payload/value.go-222-		return "[" + strings.Join(keys, " ") + "]"
pkg/payload/value.go-223-	}
pkg/payload/value.go-224-	return fmt.Sprintf("[%s ... +%d more]", strings.Join(keys[:8], " "), len(keys)-8)
pkg/payload/value.go-225-}
pkg/payload/value.go-226-
pkg/payload/value.go-227-// ReasonOf extracts the RejectReason from err, or ReasonNone if err is nil
pkg/payload/value.go-228-// or not a *RejectError.
pkg/payload/value.go-229-func ReasonOf(err error) RejectReason {
pkg/payload/value.go-230-	var re *RejectError
pkg/payload/value.go-231-	if errors.As(err, &re) {
pkg/payload/value.go-232-		return re.Reason
--
pkg/payload/json.go:286:func clip(raw []byte) string {
pkg/payload/json.go-287-	const maxDetail = 64
pkg/payload/json.go-288-	if len(raw) > maxDetail {
pkg/payload/json.go-289-		return string(raw[:maxDetail]) + "…"
pkg/payload/json.go-290-	}
pkg/payload/json.go-291-	return string(raw)
pkg/payload/json.go-292-}
pkg/payload/json.go-293-
pkg/payload/json.go-294-// ---------------------------------------------------------------------------
pkg/payload/json.go-295-// Outbound encoder
pkg/payload/json.go-296-// ---------------------------------------------------------------------------
pkg/payload/json.go-297-
pkg/payload/json.go-298-// encodeJSON builds the outbound JSON-profile `{v, t}` document. Member

← Edit pkg/payload/json.go
Index: /root/astrate-mule/pkg/payload/json.go
===================================================================
--- /root/astrate-mule/pkg/payload/json.go
+++ /root/astrate-mule/pkg/payload/json.go
@@ -208,15 +208,20 @@
 }
 if len(obj) == 0 {
 	return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
 }
+var unexpected []string
+for key := range obj {
+	if leaf, ok := leaves[key]; !ok || leaf == nil {
+		unexpected = append(unexpected, key)
+	}
+}
+if len(unexpected) > 0 {
+	return nil, rejectUnexpectedKeys(unexpected)
+}
 out := make(map[string]Value, len(obj))
 for key, eraw := range obj {
-	leaf, ok := leaves[key]
-	if !ok || leaf == nil {
-		return nil, rejectf(ReasonBadObject, "key %q matches no declared object leaf", key)
-	}
-	val, err := decodeJSONValue(eraw, leaf.ValueType)
+	val, err := decodeJSONValue(eraw, leaves[key].ValueType)
 	if err != nil {
 		return nil, annotate(err, "key %q", key)
 	}
 	out[key] = val

← Edit pkg/payload/bson.go
Index: /root/astrate-mule/pkg/payload/bson.go
===================================================================
--- /root/astrate-mule/pkg/payload/bson.go
+++ /root/astrate-mule/pkg/payload/bson.go
@@ -145,18 +145,27 @@
 }
 if len(elems) == 0 {
 	return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
 }
-out := make(map[string]Value, len(elems))
-for _, el := range elems {
+var unexpected []string
+keys := make([]string, len(elems))
+for i, el := range elems {
 	key, err := el.KeyErr()
 	if err != nil {
 		return nil, rejectf(ReasonMalformed, "invalid BSON object element: %v", err)
 	}
-	leaf, ok := leaves[key]
-	if !ok || leaf == nil {
-		return nil, rejectf(ReasonBadObject, "key %q matches no declared object leaf", key)
+	keys[i] = key
+	if leaf, ok := leaves[key]; !ok || leaf == nil {
+		unexpected = append(unexpected, key)
 	}
+}
+if len(unexpected) > 0 {
+	return nil, rejectUnexpectedKeys(unexpected)
+}
+out := make(map[string]Value, len(elems))
+for i, el := range elems {
+	key := keys[i]
+	leaf := leaves[key]
 	if _, dup := out[key]; dup {
 		return nil, rejectf(ReasonBadObject, "duplicate object key %q", key)
 	}
 	ev, err := el.ValueErr()

→ Read internal/appengine/http.go [limit=28]

← Edit internal/appengine/http.go
Index: /root/astrate-mule/internal/appengine/http.go
===================================================================
--- /root/astrate-mule/internal/appengine/http.go
+++ /root/astrate-mule/internal/appengine/http.go
@@ -12,8 +12,9 @@
 	"github.com/astrate-platform/astrate/internal/auth"
 	"github.com/astrate-platform/astrate/internal/engine"
 	"github.com/astrate-platform/astrate/internal/store"
 	"github.com/astrate-platform/astrate/pkg/astarteapi"
+	"github.com/astrate-platform/astrate/pkg/payload"
 )
 
 // maxBodyBytes caps AppEngine request bodies.
 const maxBodyBytes int64 = 1 << 20

Now the mapping itself:
← Edit internal/appengine/http.go
Index: /root/astrate-mule/internal/appengine/http.go
===================================================================
--- /root/astrate-mule/internal/appengine/http.go
+++ /root/astrate-mule/internal/appengine/http.go
@@ -680,13 +680,44 @@
 	case errors.Is(err, engine.ErrPathNotFound):
 		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Endpoint not found")
 	case errors.Is(err, store.ErrNotFound):
 		_ = astarteapi.WriteDeviceNotFound(w)
+	case payload.ReasonOf(err) == payload.ReasonBadObject:
+		// Object-shape rejection. Upstream answers 400 naming the offending
+		// keys (master b6d46ad4, #2237); before that commit the shape was
+		// still a 400, so this is a detail fix, not a status change. Placed
+		// after the sentinels because a payload rejection carries none of
+		// them, and last because every other reason stays unmapped (report
+		// only — see .mule/todo.md).
+		a.writeBadObjectError(w, err)
 	default:
 		_ = astarteapi.WriteInternalServerError(w)
 	}
 }
 
+// DetailUnexpectedObjectKey is upstream's 400 detail for an object-aggregated
+// write carrying undeclared keys (master b6d46ad4, #2237). Frozen: astartectl
+// and the SDKs match the string.
+const DetailUnexpectedObjectKey = "Unexpected object key"
+
+// writeBadObjectError renders the object-shape rejection in the measured
+// upstream envelope: {"errors": {"detail": "Unexpected object key",
+// "unexpected_keys": [...]}}. The other object-shape failures — not a
+// document, an empty document, a duplicate key — carry no key list and answer
+// the canonical 400, which is where upstream's fallback controller still
+// routes them (fallback_controller.ex:201-206).
+func (a *API) writeBadObjectError(w http.ResponseWriter, err error) {
+	var re *payload.RejectError
+	if !errors.As(err, &re) || len(re.UnexpectedKeys) == 0 {
+		_ = astarteapi.WriteBadRequest(w)
+		return
+	}
+	_ = astarteapi.WriteRawErrors(w, http.StatusBadRequest, struct {
+		Detail string   `json:"detail"`
+		Keys   []string `json:"unexpected_keys"`
+	}{Detail: DetailUnexpectedObjectKey, Keys: re.UnexpectedKeys})
+}
+
 // validationDetail strips the ErrValidation prefix for the response detail.
 func validationDetail(err error) string {
 	const prefix = "appengine: validation failed: "
 	msg := err.Error()

← Edit internal/appengine/writeerror_test.go
Index: /root/astrate-mule/internal/appengine/writeerror_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/writeerror_test.go
+++ /root/astrate-mule/internal/appengine/writeerror_test.go
@@ -8,8 +8,9 @@
 	"net/http/httptest"
 	"testing"
 
 	"github.com/astrate-platform/astrate/internal/engine"
+	"github.com/astrate-platform/astrate/pkg/payload"
 )
 
 func TestWriteErrorTaxonomy(t *testing.T) {
 	a := &API{}
@@ -17,21 +18,35 @@
 		name   string
 		err    error
 		status int
 		detail string
+		keys   []string
 	}{
 		{"device-owned write", fmt.Errorf("%w: context", engine.ErrNotServerOwned),
-			http.StatusMethodNotAllowed, "Cannot write to device owned resource"},
+			http.StatusMethodNotAllowed, "Cannot write to device owned resource", nil},
 		{"unset on server-owned datastream", fmt.Errorf("%w: context", engine.ErrNotAProperty),
-			http.StatusMethodNotAllowed, "Cannot write to read-only resource"},
+			http.StatusMethodNotAllowed, "Cannot write to read-only resource", nil},
 		{"unknown interface", fmt.Errorf("%w: context", engine.ErrInterfaceNotFound),
-			http.StatusNotFound, "Interface not found in device introspection"},
+			http.StatusNotFound, "Interface not found in device introspection", nil},
 		{"unknown endpoint", fmt.Errorf("%w: context", engine.ErrPathNotFound),
-			http.StatusBadRequest, "Endpoint not found"},
+			http.StatusBadRequest, "Endpoint not found", nil},
 		{"read-path path-not-found collision guard", ErrPathNotFound,
-			http.StatusNotFound, "Path not found"},
+			http.StatusNotFound, "Path not found", nil},
+		{"object write with an undeclared key answers 400 naming it",
+			&payload.RejectError{Reason: payload.ReasonBadObject,
+				Detail: `object keys [ghost] match no declared object leaf`,
+				UnexpectedKeys: []string{"ghost"}},
+			http.StatusBadRequest, "Unexpected object key", []string{"ghost"}},
+		{"object write with several undeclared keys lists them sorted",
+			fmt.Errorf("engine: %w", &payload.RejectError{Reason: payload.ReasonBadObject,
+				Detail:         "object keys [a b] match no declared object leaf",
+				UnexpectedKeys: []string{"a", "b"}}),
+			http.StatusBadRequest, "Unexpected object key", []string{"a", "b"}},
+		{"object write with no key list stays the canonical 400",
+			&payload.RejectError{Reason: payload.ReasonBadObject, Detail: "object-aggregation document is empty"},
+			http.StatusBadRequest, "Bad request", nil},
 		{"unknown cause stays 500", errors.New("boom"),
-			http.StatusInternalServerError, "Internal server error"},
+			http.StatusInternalServerError, "Internal server error", nil},
 	}
 	for _, tc := range cases {
 		t.Run(tc.name, func(t *testing.T) {
 			rec := httptest.NewRecorder()
@@ -40,16 +55,45 @@
 				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tc.status, rec.Body.String())
 			}
 			var body struct {
 				Errors struct {
-					Detail string `json:"detail"`
+					Detail string   `json:"detail"`
+					Keys   []string `json:"unexpected_keys"`
 				} `json:"errors"`
 			}
 			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
 				t.Fatalf("body %q: %v", rec.Body.String(), err)
 			}
 			if body.Errors.Detail != tc.detail {
 				t.Errorf("detail = %q, want %q", body.Errors.Detail, tc.detail)
 			}
+			if len(body.Errors.Keys) != len(tc.keys) {
+				t.Fatalf("unexpected_keys = %v, want %v", body.Errors.Keys, tc.keys)
+			}
+			for i, k := range tc.keys {
+				if body.Errors.Keys[i] != k {
+					t.Errorf("unexpected_keys[%d] = %q, want %q", i, body.Errors.Keys[i], k)
+				}
+			}
 		})
 	}
 }
+
+// TestWriteErrorBadObjectEnvelope freezes the exact bytes upstream emits for
+// an object write carrying an undeclared key (master b6d46ad4, #2237):
+// {"errors":{"detail":"Unexpected object key","unexpected_keys":["ghost"]}}.
+func TestWriteErrorBadObjectEnvelope(t *testing.T) {
+	a := &API{}
+	rec := httptest.NewRecorder()
+	a.writeError(rec, &payload.RejectError{
+		Reason:         payload.ReasonBadObject,
+		Detail:         "object keys [ghost] match no declared object leaf",
+		UnexpectedKeys: []string{"ghost"},
+	})
+	want := `{"errors":{"detail":"Unexpected object key","unexpected_keys":["ghost"]}}`
+	if got := rec.Body.String(); got != want {
+		t.Errorf("body = %s, want %s", got, want)
+	}
+	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
+		t.Errorf("Content-Type = %q", ct)
+	}
+}

$ go build ./... && go vet ./internal/appengine/ ./pkg/payload/ && go test ./internal/appengine/ -run 'TestWriteError' -count=1 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/appengine	0.032s

$ rg -n "ReasonBadObject|ObjectLeaves|leaves :?=|func .*[Ll]eaves" pkg/payload/*_test.go | head -20
pkg/payload/value_test.go:196:		ReasonBadObject:       "bad_object",
pkg/payload/fuzz_test.go:73:	leaves := map[string]*interfaceschema.CompiledMapping{
pkg/payload/payload_test.go:95:func objectLeaves(explicit bool) map[string]*interfaceschema.CompiledMapping {
pkg/payload/payload_test.go:395:	leaves := objectLeaves(false)
pkg/payload/payload_test.go:413:		`{"v":{"lat":1.0,"nope":2.0}}`: ReasonBadObject,    // undeclared key
pkg/payload/payload_test.go:414:		`{"v":{}}`:                     ReasonBadObject,    // empty document
pkg/payload/payload_test.go:415:		`{"v":5}`:                      ReasonBadObject,    // not a document
pkg/payload/payload_test.go:416:		`{"v":null}`:                   ReasonBadObject,    // null
pkg/payload/payload_test.go:417:		`{"v":[1]}`:                    ReasonBadObject,    // array
pkg/payload/payload_test.go:553:	if _, err := DecodeObject(dup, objectLeaves(false)); ReasonOf(err) != ReasonBadObject {
pkg/payload/payload_test.go:785:		if _, err := Encode(map[string]Value{}, nil, f); ReasonOf(err) != ReasonBadObject {
pkg/payload/payload_test.go:788:		if _, err := Encode(map[string]Value{"a": map[string]Value{"b": 1.0}}, nil, f); ReasonOf(err) != ReasonBadObject {
pkg/payload/required_test.go:18:	leaves := map[string]*interfaceschema.CompiledMapping{

→ Read pkg/payload/payload_test.go [limit=30, offset=535]
→ Read pkg/payload/payload_test.go [limit=40, offset=390]
→ Read pkg/payload/payload_test.go [limit=30, offset=90]
Now the payload test, appended to the JSON object test's neighbourhood as its own test:

← Edit pkg/payload/payload_test.go
Index: /root/astrate-mule/pkg/payload/payload_test.go
===================================================================
--- /root/astrate-mule/pkg/payload/payload_test.go
+++ /root/astrate-mule/pkg/payload/payload_test.go
@@ -554,8 +554,71 @@
 		t.Errorf("duplicate object key reason = %v; want bad_object", ReasonOf(err))
 	}
 }
 
+// TestUnexpectedObjectKeysSorted pins the object-shape rejection's key list:
+// every undeclared key is named, sorted, on both wire formats (upstream
+// `unexpected_keys`, master b6d46ad4, #2237). JSON map iteration is random, so
+// an unsorted list would make the reported body nondeterministic.
+func TestUnexpectedObjectKeysSorted(t *testing.T) {
+	leaves := objectLeaves(false)
+	cases := []struct {
+		name  string
+		in    []byte
+		want  []string
+		plain bool
+	}{
+		{name: "json one declared one undeclared",
+			in:   []byte(`{"v":{"lat":1.0,"nope":2.0}}`),
+			want: []string{"nope"}},
+		{name: "json several undeclared keys sorted",
+			in:   []byte(`{"v":{"zeta":1.0,"lat":2.0,"alpha":3.0,"mu":4.0}}`),
+			want: []string{"alpha", "mu", "zeta"}},
+		{name: "bson one declared one undeclared",
+			in: rawBSONDoc(t, rawBSONElem(byte(bson.TypeEmbeddedDocument), "v", rawBSONDoc(t,
+				rawBSONElem(byte(bson.TypeDouble), "lat", []byte{0, 0, 0, 0, 0, 0, 240, 63}),
+				rawBSONElem(byte(bson.TypeDouble), "nope", []byte{0, 0, 0, 0, 0, 0, 0, 0}),
+			))),
+			want: []string{"nope"}},
+		{name: "bson several undeclared keys sorted",
+			in: rawBSONDoc(t, rawBSONElem(byte(bson.TypeEmbeddedDocument), "v", rawBSONDoc(t,
+				rawBSONElem(byte(bson.TypeDouble), "zeta", []byte{0, 0, 0, 0, 0, 0, 0, 0}),
+				rawBSONElem(byte(bson.TypeDouble), "lat", []byte{0, 0, 0, 0, 0, 0, 240, 63}),
+				rawBSONElem(byte(bson.TypeDouble), "alpha", []byte{0, 0, 0, 0, 0, 0, 0, 0}),
+			))),
+			want: []string{"alpha", "zeta"}},
+		{name: "json empty document carries no key list",
+			in:   []byte(`{"v":{}}`),
+			plain: true},
+	}
+	for _, tc := range cases {
+		t.Run(tc.name, func(t *testing.T) {
+			_, err := DecodeObject(tc.in, leaves)
+			if ReasonOf(err) != ReasonBadObject {
+				t.Fatalf("reason = %v (err %v); want bad_object", ReasonOf(err), err)
+			}
+			var re *RejectError
+			if !errors.As(err, &re) {
+				t.Fatalf("err = %v; want a *RejectError", err)
+			}
+			if tc.plain {
+				if len(re.UnexpectedKeys) != 0 {
+					t.Fatalf("UnexpectedKeys = %v; want none", re.UnexpectedKeys)
+				}
+				return
+			}
+			if len(re.UnexpectedKeys) != len(tc.want) {
+				t.Fatalf("UnexpectedKeys = %v; want %v", re.UnexpectedKeys, tc.want)
+			}
+			for i, k := range tc.want {
+				if re.UnexpectedKeys[i] != k {
+					t.Errorf("UnexpectedKeys[%d] = %q; want %q", i, re.UnexpectedKeys[i], k)
+				}
+			}
+		})
+	}
+}
+
 // TestSizeAndCardinalityCaps pins the 64 KiB payload cap and the 1024-element
 // array cap for both formats (ROADMAP §2.4 gate).
 func TestSizeAndCardinalityCaps(t *testing.T) {
 	// Default 64 KiB cap, BSON: a >64 KiB binary payload.

$ go test ./pkg/payload/ -run 'TestUnexpectedObjectKeysSorted|TestJSONObjectAggregation|TestBSON' -count=1 2>&1 | tail -20
# github.com/astrate-platform/astrate/pkg/payload [github.com/astrate-platform/astrate/pkg/payload.test]
pkg/payload/payload_test.go:600:8: undefined: errors
FAIL	github.com/astrate-platform/astrate/pkg/payload [build failed]
FAIL

→ Read pkg/payload/payload_test.go [limit=22]

← Edit pkg/payload/payload_test.go
Index: /root/astrate-mule/pkg/payload/payload_test.go
===================================================================
--- /root/astrate-mule/pkg/payload/payload_test.go
+++ /root/astrate-mule/pkg/payload/payload_test.go
@@ -3,8 +3,9 @@
 import (
 	"bytes"
 	"encoding/binary"
 	"encoding/hex"
+	"errors"
 	"math"
 	"os"
 	"path/filepath"
 	"reflect"

$ go test ./pkg/payload/ -count=1 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/pkg/payload	0.049s

→ Read internal/engine/serverdata_test.go [limit=30, offset=305]

← Edit internal/engine/serverdata_test.go
Index: /root/astrate-mule/internal/engine/serverdata_test.go
===================================================================
--- /root/astrate-mule/internal/engine/serverdata_test.go
+++ /root/astrate-mule/internal/engine/serverdata_test.go
@@ -311,8 +311,52 @@
 		t.Errorf("bad prefix err = %v, want ErrPathNotFound", err)
 	}
 }
 
+// TestServerObjectUnexpectedKey: an object-aggregated write carrying a
+// declared and an undeclared key is rejected naming the undeclared one, which
+// is what the AppEngine surface renders as the upstream 400 envelope
+// (master b6d46ad4, #2237). Nothing is persisted or published.
+func TestServerObjectUnexpectedKey(t *testing.T) {
+	rig, fs, port := newWiredRig(t, Config{})
+	const iface = "com.astrate.test.ServerObject"
+	const def = `{
+		"interface_name": "com.astrate.test.ServerObject",
+		"version_major": 1, "version_minor": 0,
+		"type": "datastream", "ownership": "server", "aggregation": "object",
+		"mappings": [
+			{"endpoint": "/setpoints/heating", "type": "double"},
+			{"endpoint": "/setpoints/cooling", "type": "double"}
+		]
+	}`
+	fs.addInterface(realmAlphaID, storedInterface(t, realmAlphaID, 19, []byte(def)))
+	if err := rig.e.RefreshInterfaces(context.Background(), realmAlphaID); err != nil {
+		t.Fatalf("RefreshInterfaces: %v", err)
+	}
+
+	err := rig.e.PublishServerValue(context.Background(), realmAlpha, devAlpha, iface,
+		"/setpoints", json.RawMessage(`{"heating": 21.5, "ghost": 1.0}`), nil)
+	if payload.ReasonOf(err) != payload.ReasonBadObject {
+		t.Fatalf("reason = %v (err %v); want bad_object", payload.ReasonOf(err), err)
+	}
+	var re *payload.RejectError
+	if !errors.As(err, &re) {
+		t.Fatalf("err = %v; want a *payload.RejectError", err)
+	}
+	if len(re.UnexpectedKeys) != 1 || re.UnexpectedKeys[0] != "ghost" {
+		t.Errorf("UnexpectedKeys = %v; want [ghost]", re.UnexpectedKeys)
+	}
+
+	fs.mu.Lock()
+	rows := len(fs.objectRows())
+	upserts := len(fs.upserts)
+	fs.mu.Unlock()
+	if rows != 0 || upserts != 0 || len(port.published()) != 0 {
+		t.Errorf("rejected write left traces: %d object rows, %d upserts, %d publishes",
+			rows, upserts, len(port.published()))
+	}
+}
+
 // nextBusEvent returns the next live bus event, failing the test if none
 // arrives within the timeout.
 func nextBusEvent(t *testing.T, events <-chan stream.Event) stream.Event {
 	t.Helper()

$ go test ./internal/engine/ -run 'TestServerObject' -count=1 2>&1 | tail -20
(no output)

<shell_metadata>
shell tool terminated command after exceeding timeout 120000 ms. If this command is expected to take longer and is not waiting for interactive input, retry with a larger timeout value in milliseconds.
</shell_metadata>

$ go test ./internal/engine/ -run 'TestServerObject' -count=1 -timeout 90s 2>&1 | tail -25
goroutine 33 [chan receive]:
github.com/astrate-platform/astrate/internal/engine/triggers.(*Executor).worker(0x7f392301c0)
	/root/astrate-mule/internal/engine/triggers/actions.go:393 +0x108
created by github.com/astrate-platform/astrate/internal/engine/triggers.NewExecutor in goroutine 30
	/root/astrate-mule/internal/engine/triggers/actions.go:314 +0x5c4

goroutine 32 [chan receive]:
github.com/astrate-platform/astrate/internal/engine/triggers.(*Executor).worker(0x7f392301c0)
	/root/astrate-mule/internal/engine/triggers/actions.go:393 +0x108
created by github.com/astrate-platform/astrate/internal/engine/triggers.NewExecutor in goroutine 30
	/root/astrate-mule/internal/engine/triggers/actions.go:314 +0x5c4

goroutine 50 [chan receive]:
github.com/astrate-platform/astrate/internal/engine/triggers.(*Executor).worker(0x7f392301c0)
	/root/astrate-mule/internal/engine/triggers/actions.go:393 +0x108
created by github.com/astrate-platform/astrate/internal/engine/triggers.NewExecutor in goroutine 30
	/root/astrate-mule/internal/engine/triggers/actions.go:314 +0x5c4

goroutine 31 [chan receive]:
github.com/astrate-platform/astrate/internal/engine/triggers.(*Executor).worker(0x7f392301c0)
	/root/astrate-mule/internal/engine/triggers/actions.go:393 +0x108
created by github.com/astrate-platform/astrate/internal/engine/triggers.NewExecutor in goroutine 30
	/root/astrate-mule/internal/engine/triggers/actions.go:314 +0x5c4
FAIL	github.com/astrate-platform/astrate/internal/engine	90.103s
FAIL

$ go test ./internal/engine/ -run 'TestServerObject' -count=1 -timeout 60s 2>&1 | head -40
panic: test timed out after 1m0s
	running tests:
		TestServerObjectUnexpectedKey (1m0s)

goroutine 50 [running]:
testing.(*M).startAlarm.func1()
	/usr/lib/go-1.26/src/testing/testing.go:2802 +0x2cc
created by time.goFunc
	/usr/lib/go-1.26/src/time/sleep.go:215 +0x38

goroutine 1 [chan receive]:
testing.(*T).Run(0x7f3ceaa6c8, {0x898f68?, 0x19cb157b00153a20?}, 0x8c8868)
	/usr/lib/go-1.26/src/testing/testing.go:2109 +0x3bc
testing.runTests.func1(0x7f3ceaa6c8)
	/usr/lib/go-1.26/src/testing/testing.go:2585 +0x40
testing.tRunner(0x7f3ceaa6c8, 0x7f3d03fc48)
	/usr/lib/go-1.26/src/testing/testing.go:2036 +0xc4
testing.runTests({0x8a0302, 0x23}, {0x8aeb2c, 0x33}, 0x7f3cca2bd0, {0xfbea00, 0x46, 0x46}, {0x8080808080801c74?, 0x889e1d?, ...})
	/usr/lib/go-1.26/src/testing/testing.go:2583 +0x3f0
testing.(*M).Run(0x7f3ced6820)
	/usr/lib/go-1.26/src/testing/testing.go:2443 +0x578
main.main()
	_testmain.go:184 +0x80

goroutine 34 [sync.Mutex.Lock]:
internal/sync.runtime_SemacquireMutex(0x7616c0?, 0xc0?, 0x7f3d148960?)
	/usr/lib/go-1.26/src/runtime/sema.go:95 +0x28
internal/sync.(*Mutex).lockSlow(0x7f3cea5a20)
	/usr/lib/go-1.26/src/internal/sync/mutex.go:149 +0x170
internal/sync.(*Mutex).Lock(...)
	/usr/lib/go-1.26/src/internal/sync/mutex.go:70
sync.(*Mutex).Lock(...)
	/usr/lib/go-1.26/src/sync/mutex.go:46
github.com/astrate-platform/astrate/internal/engine.(*fakeStore).objectRows(0x7f3cea5a20)
	/root/astrate-mule/internal/engine/fakestore_test.go:537 +0x80
github.com/astrate-platform/astrate/internal/engine.TestServerObjectUnexpectedKey(0x7f3ceaad88)
	/root/astrate-mule/internal/engine/serverdata_test.go:350 +0x410
testing.tRunner(0x7f3ceaad88, 0x8c8868)
	/usr/lib/go-1.26/src/testing/testing.go:2036 +0xc4
created by testing.(*T).Run in goroutine 1

← Edit internal/engine/serverdata_test.go
Index: /root/astrate-mule/internal/engine/serverdata_test.go
===================================================================
--- /root/astrate-mule/internal/engine/serverdata_test.go
+++ /root/astrate-mule/internal/engine/serverdata_test.go
@@ -345,10 +345,10 @@
 if len(re.UnexpectedKeys) != 1 || re.UnexpectedKeys[0] != "ghost" {
 	t.Errorf("UnexpectedKeys = %v; want [ghost]", re.UnexpectedKeys)
 }
 
+rows := fs.objectRows()
 fs.mu.Lock()
-rows := len(fs.objectRows())
 upserts := len(fs.upserts)
 fs.mu.Unlock()
 if rows != 0 || upserts != 0 || len(port.published()) != 0 {
 	t.Errorf("rejected write left traces: %d object rows, %d upserts, %d publishes",

$ go test ./internal/engine/ -run 'TestServerObject' -count=1 -timeout 60s 2>&1 | head -20
# github.com/astrate-platform/astrate/internal/engine [github.com/astrate-platform/astrate/internal/engine.test]
internal/engine/serverdata_test.go:353:13: invalid operation: rows != 0 (mismatched types []store.ObjectRow and untyped int)
FAIL	github.com/astrate-platform/astrate/internal/engine [build failed]
FAIL

← Edit internal/engine/serverdata_test.go
Index: /root/astrate-mule/internal/engine/serverdata_test.go
===================================================================
--- /root/astrate-mule/internal/engine/serverdata_test.go
+++ /root/astrate-mule/internal/engine/serverdata_test.go
@@ -349,9 +349,9 @@
 	rows := fs.objectRows()
 	fs.mu.Lock()
 	upserts := len(fs.upserts)
 	fs.mu.Unlock()
-	if rows != 0 || upserts != 0 || len(port.published()) != 0 {
+	if len(rows) != 0 || upserts != 0 || len(port.published()) != 0 {
 		t.Errorf("rejected write left traces: %d object rows, %d upserts, %d publishes",
 			rows, upserts, len(port.published()))
 	}
 }

$ go test ./internal/engine/ -run 'TestServerObject' -count=1 -timeout 60s 2>&1 | head -20
# github.com/astrate-platform/astrate/internal/engine
# [github.com/astrate-platform/astrate/internal/engine]
internal/engine/serverdata_test.go:354:41: (*testing.common).Errorf format %d has arg rows of wrong type []github.com/astrate-platform/astrate/internal/store.ObjectRow
FAIL	github.com/astrate-platform/astrate/internal/engine [build failed]
FAIL

← Edit internal/engine/serverdata_test.go
Index: /root/astrate-mule/internal/engine/serverdata_test.go
===================================================================
--- /root/astrate-mule/internal/engine/serverdata_test.go
+++ /root/astrate-mule/internal/engine/serverdata_test.go
@@ -351,9 +351,9 @@
 	upserts := len(fs.upserts)
 	fs.mu.Unlock()
 	if len(rows) != 0 || upserts != 0 || len(port.published()) != 0 {
 		t.Errorf("rejected write left traces: %d object rows, %d upserts, %d publishes",
-			rows, upserts, len(port.published()))
+			len(rows), upserts, len(port.published()))
 	}
 }
 
 // nextBusEvent returns the next live bus event, failing the test if none

$ go test ./internal/engine/ -run 'TestServerObject' -count=1 -timeout 60s 2>&1 | head -20
ok  	github.com/astrate-platform/astrate/internal/engine	0.030s

$ gofmt -l . 2>/dev/null; echo "--- vet ---"; go vet ./... 2>&1 | tail -5; echo "--- vet integration ---"; go vet -tags integration ./internal/appengine/ ./internal/engine/ 2>&1 | tail -5
internal/appengine/writeerror_test.go
pkg/payload/payload_test.go
--- vet ---
--- vet integration ---

$ gofmt -d internal/appengine/writeerror_test.go pkg/payload/payload_test.go | head -60
diff internal/appengine/writeerror_test.go.orig internal/appengine/writeerror_test.go
--- internal/appengine/writeerror_test.go.orig
+++ internal/appengine/writeerror_test.go
@@ -33,7 +33,7 @@
 			http.StatusNotFound, "Path not found", nil},
 		{"object write with an undeclared key answers 400 naming it",
 			&payload.RejectError{Reason: payload.ReasonBadObject,
-				Detail: `object keys [ghost] match no declared object leaf`,
+				Detail:         `object keys [ghost] match no declared object leaf`,
 				UnexpectedKeys: []string{"ghost"}},
 			http.StatusBadRequest, "Unexpected object key", []string{"ghost"}},
 		{"object write with several undeclared keys lists them sorted",
diff pkg/payload/payload_test.go.orig pkg/payload/payload_test.go
--- pkg/payload/payload_test.go.orig
+++ pkg/payload/payload_test.go
@@ -588,7 +588,7 @@
 			))),
 			want: []string{"alpha", "zeta"}},
 		{name: "json empty document carries no key list",
-			in:   []byte(`{"v":{}}`),
+			in:    []byte(`{"v":{}}`),
 			plain: true},
 	}
 	for _, tc := range cases {

$ gofmt -w internal/appengine/writeerror_test.go pkg/payload/payload_test.go && gofmt -l . ; echo "fmt clean"
fmt clean

Now the full container-free gate:
$ go test ./... 2>&1 | grep -v "^ok\|no test files" | head -30; echo "EXIT=$?"
EXIT=0

$ go test -race ./pkg/payload/ ./internal/engine/ ./internal/appengine/ -count=1 -timeout 300s 2>&1 | tail -5
FAIL	github.com/astrate-platform/astrate/internal/engine	0.031s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine	0.031s
FAIL

$ cp internal/appengine/http.go /tmp/opencode/http.go.bak && python3 - <<'EOF'
import re
p='internal/appengine/http.go'
s=open(p).read()
s=s.replace("""	case payload.ReasonOf(err) == payload.ReasonBadObject:""","""	case false && payload.ReasonOf(err) == payload.ReasonBadObject:""")
open(p,'w').write(s)
EOF
go test ./internal/appengine/ -run 'TestWriteError' -count=1 2>&1 | tail -12; cp /tmp/opencode/http.go.bak internal/appengine/http.go; go build ./... && echo RESTORED
--- FAIL: TestWriteErrorTaxonomy (0.00s)
    --- FAIL: TestWriteErrorTaxonomy/object_write_with_an_undeclared_key_answers_400_naming_it (0.00s)
        writeerror_test.go:55: status = 500, want 400 (body: {"errors":{"detail":"Internal server error"}})
    --- FAIL: TestWriteErrorTaxonomy/object_write_with_several_undeclared_keys_lists_them_sorted (0.00s)
        writeerror_test.go:55: status = 500, want 400 (body: {"errors":{"detail":"Internal server error"}})
    --- FAIL: TestWriteErrorTaxonomy/object_write_with_no_key_list_stays_the_canonical_400 (0.00s)
        writeerror_test.go:55: status = 500, want 400 (body: {"errors":{"detail":"Internal server error"}})
--- FAIL: TestWriteErrorBadObjectEnvelope (0.00s)
    writeerror_test.go:94: body = {"errors":{"detail":"Internal server error"}}, want {"errors":{"detail":"Unexpected object key","unexpected_keys":["ghost"]}}
FAIL
FAIL	github.com/astrate-platform/astrate/internal/appengine	0.029s
FAIL
RESTORED

$ cat > internal/appengine/zz_probe_test.go <<'EOF'
package appengine

import (
	"net/http/httptest"
	"testing"

	"github.com/astrate-platform/astrate/pkg/payload"
)

func TestZZProbeMissingRequired(t *testing.T) {
	for _, r := range []payload.RejectReason{
		payload.ReasonMissingRequired, payload.ReasonTooLarge, payload.ReasonMalformed,
		payload.ReasonNoValue, payload.ReasonTypeMismatch, payload.ReasonValueTooLarge,
		payload.ReasonUnsetNotAllowed, payload.ReasonBadTimestamp, payload.ReasonUnknownFormat,
	} {
		rec := httptest.NewRecorder()
		(&API{}).writeError(rec, &payload.RejectError{Reason: r, Detail: "probe"})
		t.Logf("%-18s -> %d %s", r, rec.Code, rec.Body.String())
	}
}
EOF
go test ./internal/appengine/ -run TestZZProbe -count=1 -v 2>&1 | grep -E "zz_probe|PASS|FAIL" | head -20; rm internal/appengine/zz_probe_test.go
zz_probe_test.go:18: missing_required   -> 500 {"errors":{"detail":"Internal server error"}}
    zz_probe_test.go:18: too_large          -> 500 {"errors":{"detail":"Internal server error"}}
    zz_probe_test.go:18: malformed          -> 500 {"errors":{"detail":"Internal server error"}}
    zz_probe_test.go:18: no_value           -> 500 {"errors":{"detail":"Internal server error"}}
    zz_probe_test.go:18: type_mismatch      -> 500 {"errors":{"detail":"Internal server error"}}
    zz_probe_test.go:18: value_too_large    -> 500 {"errors":{"detail":"Internal server error"}}
    zz_probe_test.go:18: unset_not_allowed  -> 500 {"errors":{"detail":"Internal server error"}}
    zz_probe_test.go:18: bad_timestamp      -> 500 {"errors":{"detail":"Internal server error"}}
    zz_probe_test.go:18: unknown_format     -> 500 {"errors":{"detail":"Internal server error"}}
--- PASS: TestZZProbeMissingRequired (0.00s)
PASS

$ git status --short; echo ---; tail -12 .mule/todo.md; echo ---; grep -c "^- \[ \]" .mule/todo.md
M internal/appengine/http.go
 M internal/appengine/writeerror_test.go
 M internal/engine/serverdata_test.go
 M pkg/payload/bson.go
 M pkg/payload/json.go
 M pkg/payload/payload_test.go
 M pkg/payload/value.go
---
- [x] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, housekeeping and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
- [!] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy alias fields to `InterfaceDefinition` and `InterfaceMapping` in docs/api/astarte_realm_management_api.yaml — both declare `additionalProperties: false` (yaml:935, 972) yet omit interface-level `quality` (alias of `ownership`) and `aggregate` (alias of `aggregation`) and mapping-level `path` (alias of `endpoint`), all three of which the parser accepts and canonicalises: decoded into dedicated fields (pkg/interfaceschema/parse.go:116-122, 136-138), normalised before any rule runs (parse.go:212-236), and re-encoded canonically by `ParseInterfaceCanonical` so the store never sees them (parse.go:165-170, 307-311). A spec-conformant client — or any generated client with this schema baked in — therefore rejects a valid upstream-shaped install/update with a spurious 422 while the service would have accepted and stored it as `ownership`/`aggregation`/`endpoint`. Add the three as documented deprecated aliases (enums mirroring ownership/aggregation, `path` a plain string) or relax `additionalProperties`, and record the three rules the code enforces: `ownership` + `quality` is a violation (parse.go:218) and `aggregation` + `aggregate` is a violation (parse.go:227), while `endpoint` + `path` silently keeps `endpoint` (parse.go:307-311). None of the three may appear in a `required` list. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstrained string described only as "The device hardware ID." — `Service.DeleteDevice` runs the path value through `deviceid.Parse` (internal/realm/service.go:118-121), which accepts exactly 22 characters of unpadded base64url in strict mode and rejects padding, the standard alphabet's `+` and `/`, and non-canonical trailing bits (pkg/deviceid/deviceid.go:24-26, 33-37, 42-60); every rejection is folded into `store.ErrNotFound` (service.go:119-121) and answered `404 {"errors":{"detail":"Device not found"}}` (internal/realm/http.go:107-109), so a client that sends a well-formed standard-base64 or 23-character padded id is told the device does not exist rather than that its id is malformed. State the encoding and the 22-character length and add `pattern`/`minLength`/`maxLength`, exactly as docs-sync-pairing-deviceid-base64url did for the pairing twin. Do NOT encode the strict trailing-bits rule in the pattern: whether it is a deviation worth recording is still open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_realm_management_api.yaml — `data` is declared as a bare `type: object` with no properties (yaml:250-261) although the handler parses it with the same `interfaceschema.ParseInterfaceCanonical` as `installInterface` (internal/realm/service.go:196 vs 161) and then checks the body's `interface_name`/`version_major` against the URL's `{name}`/`{major}` (service.go:203-208), so the update body — the one carrying every interface rule — is the only interface body in the spec with nothing to validate against. `$ref` `#/components/schemas/InterfaceDefinition` as `installInterface` does (yaml:115-117), and say in the operation description that the body must name the same interface and major as the URL (disagreement is 409, already documented) and that `version_minor` must increase (409). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.
- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` shows `detail: Bad Request` (yaml:1732) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1792), while the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58 via `WriteBadRequest`/`WriteInternalServerError`, golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json); every 400/500 on that surface goes through the same astarteapi writers. The other detail examples on the file are already correct — change only those two — and add the docs-half assertion `TestAppEngineErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go, which pins the pairing and realm-management halves. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDetail` schema-level example (yaml:361-362) and the `BadRequest` response (yaml:385) show `detail: Bad Request`, and `InternalServerError` (yaml:477) shows `detail: Internal Server Error`, but the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58, via WriteBadRequest/WriteInternalServerError, pinned by pkg/astarteapi/testdata/error_bad_request.json and error_internal.json). While you are in `ErrorDetail`: its `example` sits at the schema level and shows a bare `{detail: ...}` where the envelope the wire sends is `{"errors":{"detail":...}}` — move it under `errors` (or drop it, as the realm-management and pairing `ErrorDetail` schemas carry no example) so a generated client is not shown the wrong shape. Add `TestHousekeepingErrorDetailExamples` beside the pairing and realm-management halves in internal/swagger/swagger_test.go. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `"Internal server error"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implements `device_empty_cache_received` end to end — `SimpleEvent` field 23, an empty `DeviceEmptyCacheReceivedEvent` whose encoder emits exactly `{"type":"device_empty_cache_received"}` (astarte_core encoder.ex:54-62), fired by DUP at `control_handler.ex:161` *inside* the `with` that runs `send_control_consumer_properties` → `resend_all_properties` → `set_pending_empty_cache`, so it fires only when all three succeed and never on failure, and the commit deletes the `# TODO: implement on_empty_cache_received` from `astarte_events/lib/astarte_events/triggers/core.ex`. Astrate already accepts the trigger but never evaluates it: `OnDeviceEmptyCacheReceived` sits in `deviceOns` with value `false` (internal/engine/triggers/match.go:123, the `evaluated` flag read at match.go:406) and the comment's stated reason ("upstream defines no SimpleEvent variant for it") is now false. Add the event body constructor to internal/engine/triggers/events.go (the payload is `{"type":"device_empty_cache_received"}` and nothing else — no device id, no metadata), flip the `deviceOns` value to true, and fire it from `handleEmptyCache` (internal/engine/control.go:58) with the same idiom `handleLifecycle` uses for `device_connected` (`e.fireDevice` + `stream.Publish`, internal/engine/engine.go:346-353), at the point where `m.Ack()` runs — i.e. after `resendServerProperties` and `sendConsumerProperties` have both succeeded inside `retryStore`. Add a T1 test that a trigger on it fires exactly once per accepted emptyCache and does not fire when the resend fails. — BLOCKED: wrote nothing
- [ ] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
---
15

$ tail -c 400 .mule/todo.md | tail -2 | cut -c1-80
t key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}`

$ cat >> .mule/todo.md <<'EOF'
- [ ] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress.
- [ ] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
EOF
grep -c "^- \[ \]" .mule/todo.md
17

$ git --no-pager diff -- pkg/payload internal/appengine/http.go internal/engine | head -160
diff --git a/internal/appengine/http.go b/internal/appengine/http.go
index 6837941..2df566e 100644
--- a/internal/appengine/http.go
+++ b/internal/appengine/http.go
@@ -13,6 +13,7 @@ import (
 	"github.com/astrate-platform/astrate/internal/engine"
 	"github.com/astrate-platform/astrate/internal/store"
 	"github.com/astrate-platform/astrate/pkg/astarteapi"
+	"github.com/astrate-platform/astrate/pkg/payload"
 )
 
 // maxBodyBytes caps AppEngine request bodies.
@@ -680,11 +681,42 @@ func (a *API) writeError(w http.ResponseWriter, err error) {
 		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Endpoint not found")
 	case errors.Is(err, store.ErrNotFound):
 		_ = astarteapi.WriteDeviceNotFound(w)
+	case payload.ReasonOf(err) == payload.ReasonBadObject:
+		// Object-shape rejection. Upstream answers 400 naming the offending
+		// keys (master b6d46ad4, #2237); before that commit the shape was
+		// still a 400, so this is a detail fix, not a status change. Placed
+		// after the sentinels because a payload rejection carries none of
+		// them, and last because every other reason stays unmapped (report
+		// only — see .mule/todo.md).
+		a.writeBadObjectError(w, err)
 	default:
 		_ = astarteapi.WriteInternalServerError(w)
 	}
 }
 
+// DetailUnexpectedObjectKey is upstream's 400 detail for an object-aggregated
+// write carrying undeclared keys (master b6d46ad4, #2237). Frozen: astartectl
+// and the SDKs match the string.
+const DetailUnexpectedObjectKey = "Unexpected object key"
+
+// writeBadObjectError renders the object-shape rejection in the measured
+// upstream envelope: {"errors": {"detail": "Unexpected object key",
+// "unexpected_keys": [...]}}. The other object-shape failures — not a
+// document, an empty document, a duplicate key — carry no key list and answer
+// the canonical 400, which is where upstream's fallback controller still
+// routes them (fallback_controller.ex:201-206).
+func (a *API) writeBadObjectError(w http.ResponseWriter, err error) {
+	var re *payload.RejectError
+	if !errors.As(err, &re) || len(re.UnexpectedKeys) == 0 {
+		_ = astarteapi.WriteBadRequest(w)
+		return
+	}
+	_ = astarteapi.WriteRawErrors(w, http.StatusBadRequest, struct {
+		Detail string   `json:"detail"`
+		Keys   []string `json:"unexpected_keys"`
+	}{Detail: DetailUnexpectedObjectKey, Keys: re.UnexpectedKeys})
+}
+
 // validationDetail strips the ErrValidation prefix for the response detail.
 func validationDetail(err error) string {
 	const prefix = "appengine: validation failed: "
diff --git a/internal/engine/serverdata_test.go b/internal/engine/serverdata_test.go
index ebafa09..e9b4cf9 100644
--- a/internal/engine/serverdata_test.go
+++ b/internal/engine/serverdata_test.go
@@ -312,6 +312,50 @@ func TestServerObjectAggregate(t *testing.T) {
 	}
 }
 
+// TestServerObjectUnexpectedKey: an object-aggregated write carrying a
+// declared and an undeclared key is rejected naming the undeclared one, which
+// is what the AppEngine surface renders as the upstream 400 envelope
+// (master b6d46ad4, #2237). Nothing is persisted or published.
+func TestServerObjectUnexpectedKey(t *testing.T) {
+	rig, fs, port := newWiredRig(t, Config{})
+	const iface = "com.astrate.test.ServerObject"
+	const def = `{
+		"interface_name": "com.astrate.test.ServerObject",
+		"version_major": 1, "version_minor": 0,
+		"type": "datastream", "ownership": "server", "aggregation": "object",
+		"mappings": [
+			{"endpoint": "/setpoints/heating", "type": "double"},
+			{"endpoint": "/setpoints/cooling", "type": "double"}
+		]
+	}`
+	fs.addInterface(realmAlphaID, storedInterface(t, realmAlphaID, 19, []byte(def)))
+	if err := rig.e.RefreshInterfaces(context.Background(), realmAlphaID); err != nil {
+		t.Fatalf("RefreshInterfaces: %v", err)
+	}
+
+	err := rig.e.PublishServerValue(context.Background(), realmAlpha, devAlpha, iface,
+		"/setpoints", json.RawMessage(`{"heating": 21.5, "ghost": 1.0}`), nil)
+	if payload.ReasonOf(err) != payload.ReasonBadObject {
+		t.Fatalf("reason = %v (err %v); want bad_object", payload.ReasonOf(err), err)
+	}
+	var re *payload.RejectError
+	if !errors.As(err, &re) {
+		t.Fatalf("err = %v; want a *payload.RejectError", err)
+	}
+	if len(re.UnexpectedKeys) != 1 || re.UnexpectedKeys[0] != "ghost" {
+		t.Errorf("UnexpectedKeys = %v; want [ghost]", re.UnexpectedKeys)
+	}
+
+	rows := fs.objectRows()
+	fs.mu.Lock()
+	upserts := len(fs.upserts)
+	fs.mu.Unlock()
+	if len(rows) != 0 || upserts != 0 || len(port.published()) != 0 {
+		t.Errorf("rejected write left traces: %d object rows, %d upserts, %d publishes",
+			len(rows), upserts, len(port.published()))
+	}
+}
+
 // nextBusEvent returns the next live bus event, failing the test if none
 // arrives within the timeout.
 func nextBusEvent(t *testing.T, events <-chan stream.Event) stream.Event {
diff --git a/pkg/payload/bson.go b/pkg/payload/bson.go
index 5217c60..372ee12 100644
--- a/pkg/payload/bson.go
+++ b/pkg/payload/bson.go
@@ -146,16 +146,25 @@ func decodeBSONObject(rv bson.RawValue, leaves map[string]*interfaceschema.Compi
 	if len(elems) == 0 {
 		return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
 	}
-	out := make(map[string]Value, len(elems))
-	for _, el := range elems {
+	var unexpected []string
+	keys := make([]string, len(elems))
+	for i, el := range elems {
 		key, err := el.KeyErr()
 		if err != nil {
 			return nil, rejectf(ReasonMalformed, "invalid BSON object element: %v", err)
 		}
-		leaf, ok := leaves[key]
-		if !ok || leaf == nil {
-			return nil, rejectf(ReasonBadObject, "key %q matches no declared object leaf", key)
+		keys[i] = key
+		if leaf, ok := leaves[key]; !ok || leaf == nil {
+			unexpected = append(unexpected, key)
 		}
+	}
+	if len(unexpected) > 0 {
+		return nil, rejectUnexpectedKeys(unexpected)
+	}
+	out := make(map[string]Value, len(elems))
+	for i, el := range elems {
+		key := keys[i]
+		leaf := leaves[key]
 		if _, dup := out[key]; dup {
 			return nil, rejectf(ReasonBadObject, "duplicate object key %q", key)
 		}
diff --git a/pkg/payload/json.go b/pkg/payload/json.go
index a18a896..52555b5 100644
--- a/pkg/payload/json.go
+++ b/pkg/payload/json.go
@@ -209,13 +209,18 @@ func decodeJSONObject(raw json.RawMessage, leaves map[string]*interfaceschema.Co
 	if len(obj) == 0 {
 		return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
 	}
+	var unexpected []string
+	for key := range obj {
+		if leaf, ok := leaves[key]; !ok || leaf == nil {
+			unexpected = append(unexpected, key)
+		}
+	}
+	if len(unexpected) > 0 {

$ python3 - <<'EOF'
p='internal/appengine/http.go'
s=open(p).read()
s=s.replace("""// DetailUnexpectedObjectKey is upstream's 400 detail for an object-aggregated
// write carrying undeclared keys (master b6d46ad4, #2237). Frozen: astartectl
// and the SDKs match the string.
const DetailUnexpectedObjectKey = "Unexpected object key\"""","""// detailUnexpectedObjectKey is upstream's 400 detail for an object-aggregated
// write carrying undeclared keys (master b6d46ad4, #2237). Frozen string:
// astartectl and the SDKs match on it.
const detailUnexpectedObjectKey = "Unexpected object key\"""")
s=s.replace("}{Detail: DetailUnexpectedObjectKey, Keys: re.UnexpectedKeys})","}{Detail: detailUnexpectedObjectKey, Keys: re.UnexpectedKeys})")
open(p,'w').write(s)
EOF
gofmt -l . ; go vet ./internal/appengine/ && go test ./internal/appengine/ ./pkg/payload/ ./internal/engine/ -count=1 2>&1 | tail -5
ok  	github.com/astrate-platform/astrate/internal/appengine	0.026s
ok  	github.com/astrate-platform/astrate/pkg/payload	0.051s
ok  	github.com/astrate-platform/astrate/internal/engine	1.390s

$ go vet ./... && go test ./... 2>&1 | grep -v "^ok\|no test files"; echo "GATE_EXIT=$?"; git status --short
GATE_EXIT=1
 M .mule/todo.md
 M internal/appengine/http.go
 M internal/appengine/writeerror_test.go
 M internal/engine/serverdata_test.go
 M pkg/payload/bson.go
 M pkg/payload/json.go
 M pkg/payload/payload_test.go
 M pkg/payload/value.go

$ go test ./... > /tmp/opencode/gate.log 2>&1; echo "go test exit=$?"; grep -c FAIL /tmp/opencode/gate.log; go vet ./... > /tmp/opencode/vet.log 2>&1; echo "go vet exit=$?"
go test exit=0
0
go vet exit=0

Done: an object-aggregated write carrying an undeclared key now answers 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":[...]}}` on all six AppEngine data write ops instead of 500, with every offending key named (sorted) by the payload decoder; the `ReasonMissingRequired` half is reported only, not guessed.

Files: `pkg/payload/value.go` (new `RejectError.UnexpectedKeys` + `rejectUnexpectedKeys`/`clipKeys`), `pkg/payload/json.go` + `pkg/payload/bson.go` (collect **all** undeclared keys, sorted, before decoding values — was first-key-only, nondeterministic for JSON map order), `internal/appengine/http.go` (`writeError` maps `payload.ReasonBadObject`, new `writeBadObjectError`, `detailUnexpectedObjectKey`; other object-shape failures with no key list answer the canonical 400 per fallback_controller.ex:201-206), tests in `pkg/payload/payload_test.go`, `internal/engine/serverdata_test.go`, `internal/appengine/writeerror_test.go`, `.mule/todo.md` (2 follow-ups).

Verified: `go vet ./...` -> pass (exit 0); `go test ./...` -> pass, 0 FAIL, 0 failures anywhere (run **without** `-race`: this host has the same 39-bit VMA kernel as the Pi, `go test -race` dies with `ThreadSanitizer: unsupported VMA range`); `gofmt -l .` -> clean; `go vet -tags integration ./internal/appengine/ ./internal/engine/` -> clean (builds, not run); mutation check — disabling only the new `writeError` case makes `TestWriteErrorTaxonomy` and `TestWriteErrorBadObjectEnvelope` fail with `status = 500, want 400`, confirming the old behaviour and that the tests are not vacuous; the probe measured the report claim below.

Unsure: two things I did not measure on upstream (no container-free way to reach it, and the Legion Go is off per AGENTS.md, so I did not probe it):
1. That pre-b6d46ad4 upstream answered 400 for *every* object-shape failure — taken from the task line's `fallback_controller.ex:201-206` note. Hence the keyless variants (not a document / empty document / duplicate key) answer the canonical `"Bad request"`, which is a guess about the detail string only.
2. `unexpected_keys` list length is unbounded (a payload can carry up to the 64 KiB cap of junk keys); upstream's cap is unmeasured. The `Detail` log line *is* bounded (8 keys + a count).
Also: a document carrying both an undeclared key and a type error now reports the unexpected key rather than the type error (the key check moved ahead of value coercion, mirroring upstream's `unexpected_keys` computed first). This changes the metric label from `type_mismatch` to `bad_object` in that mixed case.

Report half (measured, not changed): the hole is not specific to `bad_object`. A throwaway probe over `writeError` returned **500 "Internal server error"** for every one of the nine remaining reasons — `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp`, `unknown_format`. So `ReasonMissingRequired` has the same hole, and the upstream 422 `missing_required_mapping` status/detail must be probed on the Legion Go before mapping — left as a follow-up line, not guessed.

Follow-ups: appended to `.mule/todo.md` — `appengine-payload-reason-status-map` (map the remaining nine reasons, each probed on the Legion Go; notes that `too_large` may be unreachable from REST because `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner`) and `appengine-missing-required-422` (probe the status/detail, then map `ReasonMissingRequired`).
