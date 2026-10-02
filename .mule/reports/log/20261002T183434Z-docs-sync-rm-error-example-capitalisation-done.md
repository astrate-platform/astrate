slug: docs-sync-rm-error-example-capitalisation
verdict: done
at:  1d1e6f4
ran: 2026-10-02T18:34:34Z on DietPi in 751s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ rg -n "detail: (Bad|Internal)" /root/astrate-mule/docs/api/astrate_realm_management_api.yaml
rg: /root/astrate-mule/docs/api/astrate_realm_management_api.yaml: IO error for operation on /root/astrate-mule/docs/api/astrate_realm_management_api.yaml: No such file or directory (os error 2)

$ sed -n '25,65p' /root/astrate-mule/pkg/astarteapi/envelope.go
// DefaultMaxBodyBytes is the request body size cap callers are expected to
// pass to DecodeData unless an endpoint has a documented reason to differ
// (interface uploads, for example, may need more than a pairing request).
const DefaultMaxBodyBytes int64 = 1 << 20 // 1 MiB

// Canonical upstream error detail strings. These are frozen: SDK and
// astartectl error paths match on them.
const (
	// DetailBadRequest is the canonical 400 detail. Measured upstream 1.2.0
	// (2026-08-24, verify batch #57): Phoenix renders "Bad request" — the
	// capital-R form this constant carried before was reconstructed, not
	// observed, and every probed 400 row answers with the lowercase r.
	DetailBadRequest = "Bad request"
	// DetailUnauthorized is the canonical 401 detail.
	DetailUnauthorized = "Unauthorized"
	// DetailForbidden is the canonical 403 detail.
	DetailForbidden = "Forbidden"
	// DetailNotFound is the canonical generic 404 detail.
	DetailNotFound = "Not Found"
	// DetailDeviceNotFound is the canonical 404 detail for unknown devices
	// (upstream AppEngine/Pairing shape).
	DetailDeviceNotFound = "Device not found"
	// DetailRouteNotFound is the 404 detail upstream AppEngine and
	// RealmManagement emit for a path matching no route at all. Distinct from
	// DetailNotFound, which a handler emits for a route that matched and found
	// no resource; the capitalisation difference is upstream's, not a typo.
	DetailRouteNotFound = "Not found"
	// DetailPageNotFound is the 404 detail upstream Pairing emits for a path
	// matching no route at all.
	DetailPageNotFound = "Page not found"
	// DetailInternalServerError is the canonical 500 detail. Measured
	// upstream 1.2.0 (2026-08-24, verify batch #57): Phoenix renders
	// "Internal server error" — the capital-S form was reconstructed.
	DetailInternalServerError = "Internal server error"
)

// ErrMissingData is wrapped by DecodeData when the request body has no
// "data" key (or it is JSON null) — upstream rejects such bodies uniformly.
var ErrMissingData = errors.New(`missing "data" key in request body`)

// ErrBodyTooLarge is wrapped by DecodeData when the request body exceeds the

$ ls /root/astrate-mule/docs/ && ls /root/astrate-mule/docs/api 2>/dev/null || rg --files /root/astrate-mule/docs | head -50
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
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

$ ls -la /root/astrate-mule/docs/api/
total 172
drwxr-xr-x 2 root root  4096 Sep  4 18:55 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 59526 Sep 27 13:31 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 15547 Sep 25 21:48 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 21724 Sep 28 19:59 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 45387 Sep 25 22:07 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 14977 Sep 20 21:25 astrate_native_api.yaml

$ pwd; rg -n "detail:" docs/api/astrate_realm_management_api.yaml | head -30
/root/astrate-mule
rg: docs/api/astrate_realm_management_api.yaml: IO error for operation on docs/api/astrate_realm_management_api.yaml: No such file or directory (os error 2)

✗ Read docs/api/astrate_realm_management_api.yaml failed [offset=1245, limit=120]
Error: File not found: /root/astrate-mule/docs/api/astrate_realm_management_api.yaml

$ ls -la /root/ | head -20; echo ---; readlink -f /root/astrate-mule; echo ---; stat -c '%n %s %F' /root/astrate-mule/docs/api/*.yaml; echo ---; head -c 100 /root/astrate-mule/docs/api/astrate_realm_management_api.yaml; echo; echo ---; wc -l /root/astrate-mule/docs/api/astrate_realm_management_api.yaml
total 384
drwx------ 29 root root    4096 Sep 22 12:57 .
drwxr-xr-x 18 root root    4096 Aug  9 00:18 ..
-rw-------  1 root root      51 Jul 21 00:10 .Xauthority
-rw-------  1 root root      58 Jul 17 17:24 .bash_history
-rw-r--r--  1 root root     629 Jul 17 21:11 .bashrc
drwx------ 21 root root    4096 Sep  2 21:44 .cache
drwxr-xr-x  4 root root    4096 Aug 20 16:08 .cargo
drwxr-xr-x 16 root root    4096 Sep  4 23:23 .claude
drwxr-xr-x  4 root root    4096 Jul 27 19:17 .claude-skills
-rw-------  1 root root   44525 Jul 23 18:13 .claude.json
drwx------ 20 root root    4096 Aug  3 14:55 .config
-rw-r--r--  1 root root      35 Jul 17 18:13 .dmrc
-rw-r--r--  1 root root     298 Aug  4 03:04 .gitconfig
drwx------  4 root root    4096 Aug 31 01:48 .gnupg
-rw-r--r--  1 root root     139 Jul 17 18:13 .gtk-bookmarks
drwx------ 23 root root    4096 Oct  2 20:22 .hermes
drwxr-xr-x  2 root root    4096 Jul 17 17:25 .kbd
drwx------  5 root root    4096 Jul 17 21:01 .local
drwxr-xr-x  6 root root    4096 Jul 27 02:54 .npm
---
/root/astrate-mule
---
/root/astrate-mule/docs/api/astarte_appengine_api.yaml 59526 regular file
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml 15547 regular file
/root/astrate-mule/docs/api/astarte_pairing_api.yaml 21724 regular file
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml 45387 regular file
/root/astrate-mule/docs/api/astrate_native_api.yaml 14977 regular file
---

---
head: cannot open '/root/astrate-mule/docs/api/astrate_realm_management_api.yaml' for reading: No such file or directory
wc: /root/astrate-mule/docs/api/astrate_realm_management_api.yaml: No such file or directory

$ mount | grep -i astrate; echo "--- lsof/xattr"; getfattr -d -m - /root/astrate-mule/docs/api/astarte_realm_management_api.yaml 2>&1 | head; echo "--- python"; python3 -c "
p='/root/astrate-mule/docs/api/astarte_realm_management_api.yaml'
try:
    d=open(p,'rb').read(120); print(d)
except Exception as e: print('ERR',e)
"; echo "--- pair"; python3 -c "
p='/root/astrate-mule/docs/api/astarte_pairing_api.yaml'
try:
    d=open(p,'rb').read(60); print(d)
except Exception as e: print('ERR',e)
"
--- lsof/xattr
/bin/sh: 1: getfattr: not found
--- python
b'openapi: "3.0.3"\ninfo:\n  title: Astarte Realm Management API\n  description: |\n    Realm-level administration API for Ast'
--- pair
b'openapi: "3.0.3"\ninfo:\n  title: Astarte Pairing API\n  descri'

$ python3 - <<'EOF'
p='/root/astrate-mule/docs/api/astarte_realm_management_api.yaml'
lines=open(p).read().split('\n')
for i,l in enumerate(lines,1):
    if 'detail' in l:
        print(i, repr(l))
EOF
33 '        `detailed` query parameter is `true`, serves the additive 1.4-style'
34 '        detailed listing instead (issue #66): one fully materialised interface'
43 '        - name: detailed'
47 '            When `true`, return the detailed interface listing (one materialised'
54 '          description: A list of interface names, or of detailed interface documents when `detailed=true`.'
68 '                detailed:'
69 '                  summary: Detailed interface listing (`?detailed=true`)'
148 '                flat_detail:'
152 '                      detail: \'missing "interface_name"\''
284 '                flat_detail:'
288 '                      detail: \'missing "interface_name"\''
341 '                      detail: "Interface can\'t be deleted"'
346 '                      detail: "Interface can\'t be deleted since it\'s currently used"'
430 '                flat_detail:'
434 '                      detail: "trigger requires a name"'
1191 '          required: [detail]'
1193 '            detail:'
1265 '              detail: Bad Request'
1275 '              detail: Unauthorized'
1285 '              detail: Forbidden'
1295 '              detail: Not Found'
1305 '              detail: Device not found'
1318 '                  detail: Already exists'
1323 "                  detail: Interface name doesn't match the one in the interface json"
1328 "                  detail: Interface major version doesn't match the one in the interface json"
1338 '              detail: "invalid interface: missing \\"interface_name\\""'
1348 '              detail: Internal Server Error'

$ cd /root/astrate-mule && rg -n "detail" pkg/astarteapi/testdata/error_bad_request.json pkg/astarteapi/testdata/error_internal.json; echo "=== docs sync tests ==="; rg -ln "docs/api" --glob '*.go' | head -20; echo "=== pairing spec detail lines ==="; python3 -c "
lines=open('docs/api/astarte_pairing_api.yaml').read().split('\n')
print('\n'.join(f'{i}: {l}' for i,l in enumerate(lines,1) if 'detail:' in l and 'flat_detail' not in l and 'flat_detail' not in l))"
pkg/astarteapi/testdata/error_bad_request.json:1:{"errors":{"detail":"Bad request"}}
pkg/astarteapi/testdata/error_internal.json:1:{"errors":{"detail":"Internal server error"}}
=== docs sync tests ===
=== pairing spec detail lines ===
59:                   detail: Not Found
68:                   detail: unhealthy
563:             detail:
586:               detail: Bad request
596:               detail: Unauthorized
606:               detail: Forbidden
616:               detail: Device not found
676:               detail: Too Many Requests
686:               detail: Internal server error

→ Read docs/embed.go
→ Read docs/Makefile

$ cd /root/astrate-mule && rg -n "APIYAML" --glob '*.go' | head -30; echo "=== yaml/openapi tests ==="; rg -ln "openapi|yaml.Unmarshal" --glob '*_test.go' | head -20
docs/embed.go:12:// APIYAML holds the OpenAPI 3.0 YAML specifications for all five API surfaces.
docs/embed.go:15:var APIYAML embed.FS
internal/swagger/swagger.go:19:	apiRoot := mustSub(docs.APIYAML, "api")
internal/swagger/swagger.go:59:	_ = fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
internal/swagger/swagger_test.go:61:			want, err := docs.APIYAML.ReadFile("api/" + name)
internal/swagger/swagger_test.go:108:		if _, err := docs.APIYAML.ReadFile("api/" + name); err != nil {
internal/swagger/swagger_test.go:109:			t.Errorf("Specs() entry %q cannot be read from docs.APIYAML: %v", name, err)
internal/swagger/swagger_test.go:147:	b, err := docs.APIYAML.ReadFile("api/" + filename)
internal/swagger/swagger_test.go:203:	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
internal/swagger/swagger_test.go:246:	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
internal/swagger/swagger_test.go:321:	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
internal/swagger/swagger_test.go:381:	b, err := docs.APIYAML.ReadFile("api/astarte_appengine_api.yaml")
internal/swagger/swagger_test.go:425:	b, err := docs.APIYAML.ReadFile("api/astarte_appengine_api.yaml")
internal/swagger/swagger_test.go:470:	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
internal/swagger/swagger_test.go:519:	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
internal/swagger/swagger_test.go:594:	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
internal/swagger/swagger_test.go:707:	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
=== yaml/openapi tests ===

→ Read internal/swagger/swagger_test.go [offset=190, limit=140]

$ cd /root/astrate-mule && rg -n "Bad request|Internal server error|capitalis|Capitalis" internal/swagger/swagger_test.go
463:// request" and "Internal server error" — the capital-R/capital-S spellings the

→ Read internal/swagger/swagger_test.go [limit=130, offset=440]

$ cd /root/astrate-mule && python3 -c "
lines=open('docs/api/astarte_realm_management_api.yaml').read().split('\n')
print('\n'.join(f'{i}: {l}' for i,l in enumerate(lines[1170:1360],1171)))"
1171:             type: integer
1172: 
1173:     DataEnvelopeInterfaceDetailedArray:
1174:       type: object
1175:       required: [data]
1176:       properties:
1177:         data:
1178:           type: array
1179:           description: One fully materialised interface document per installed interface major.
1180:           items:
1181:             $ref: "#/components/schemas/InterfaceDefinition"
1182: 
1183:     # ── Error envelopes ──
1184: 
1185:     ErrorDetail:
1186:       type: object
1187:       required: [errors]
1188:       properties:
1189:         errors:
1190:           type: object
1191:           required: [detail]
1192:           properties:
1193:             detail:
1194:               type: string
1195: 
1196:     TriggerErrorsDetail:
1197:       type: object
1198:       required: [errors]
1199:       properties:
1200:         errors:
1201:           type: object
1202:           description: Upstream nested changeset envelope for a failed trigger compile.
1203:           properties:
1204:             action:
1205:               type: object
1206:               description: Per-field error messages on the trigger action (changeset shape).
1207:               additionalProperties:
1208:                 type: array
1209:                 items:
1210:                   type: string
1211:             simple_triggers:
1212:               type: array
1213:               description: Index-aligned per-condition errors; empty objects mark entries without errors.
1214:               items:
1215:                 type: object
1216:                 additionalProperties:
1217:                   type: array
1218:                   items:
1219:                     type: string
1220: 
1221:     InterfaceViolationsDetail:
1222:       type: object
1223:       required: [errors]
1224:       properties:
1225:         errors:
1226:           type: object
1227:           description: Interface schema violations in upstream's changeset shape.
1228:           properties:
1229:             mappings:
1230:               type: array
1231:               description: Index-aligned per-mapping errors; empty objects mark entries without errors.
1232:               items:
1233:                 type: object
1234:                 additionalProperties:
1235:                   type: array
1236:                   items:
1237:                     type: string
1238:           additionalProperties:
1239:             type: array
1240:             items:
1241:               type: string
1242: 
1243:     InterfaceFieldErrorsDetail:
1244:       type: object
1245:       required: [errors]
1246:       properties:
1247:         errors:
1248:           type: object
1249:           required: [error_name]
1250:           properties:
1251:             error_name:
1252:               type: array
1253:               items:
1254:                 type: string
1255: 
1256:   responses:
1257:     BadRequest:
1258:       description: Bad Request — malformed or missing data envelope.
1259:       content:
1260:         application/json; charset=utf-8:
1261:           schema:
1262:             $ref: "#/components/schemas/ErrorDetail"
1263:           example:
1264:             errors:
1265:               detail: Bad Request
1266: 
1267:     Unauthorized:
1268:       description: Unauthorized — missing or invalid JWT.
1269:       content:
1270:         application/json; charset=utf-8:
1271:           schema:
1272:             $ref: "#/components/schemas/ErrorDetail"
1273:           example:
1274:             errors:
1275:               detail: Unauthorized
1276: 
1277:     Forbidden:
1278:       description: Forbidden — valid JWT but insufficient permissions.
1279:       content:
1280:         application/json; charset=utf-8:
1281:           schema:
1282:             $ref: "#/components/schemas/ErrorDetail"
1283:           example:
1284:             errors:
1285:               detail: Forbidden
1286: 
1287:     NotFound:
1288:       description: Not Found — resource does not exist.
1289:       content:
1290:         application/json; charset=utf-8:
1291:           schema:
1292:             $ref: "#/components/schemas/ErrorDetail"
1293:           example:
1294:             errors:
1295:               detail: Not Found
1296: 
1297:     DeviceNotFound:
1298:       description: Device not found — unknown device ID.
1299:       content:
1300:         application/json; charset=utf-8:
1301:           schema:
1302:             $ref: "#/components/schemas/ErrorDetail"
1303:           example:
1304:             errors:
1305:               detail: Device not found
1306: 
1307:     Conflict:
1308:       description: Conflict — resource already exists or the submitted data disagrees with the resource.
1309:       content:
1310:         application/json; charset=utf-8:
1311:           schema:
1312:             $ref: "#/components/schemas/ErrorDetail"
1313:           examples:
1314:             alreadyExists:
1315:               summary: The resource already exists.
1316:               value:
1317:                 errors:
1318:                   detail: Already exists
1319:             nameMismatch:
1320:               summary: The interface name in the body does not match the one in the URL.
1321:               value:
1322:                 errors:
1323:                   detail: Interface name doesn't match the one in the interface json
1324:             majorMismatch:
1325:               summary: The interface major version in the body does not match the one in the URL.
1326:               value:
1327:                 errors:
1328:                   detail: Interface major version doesn't match the one in the interface json
1329: 
1330:     ValidationError:
1331:       description: Unprocessable Entity — request body failed validation.
1332:       content:
1333:         application/json; charset=utf-8:
1334:           schema:
1335:             $ref: "#/components/schemas/ErrorDetail"
1336:           example:
1337:             errors:
1338:               detail: "invalid interface: missing \"interface_name\""
1339: 
1340:     InternalServerError:
1341:       description: Internal Server Error — unexpected failure.
1342:       content:
1343:         application/json; charset=utf-8:
1344:           schema:
1345:             $ref: "#/components/schemas/ErrorDetail"
1346:           example:
1347:             errors:
1348:               detail: Internal Server Error
1349:

$ cd /root/astrate-mule && rg -n "DetailBadRequest|DetailInternalServerError" internal/realm internal/astarte* --glob '!*_test.go' | head -20; echo "=== test refs ==="; rg -n "DetailBadRequest|DetailInternalServerError" internal/realm --glob '*_test.go' | head -20
rg: internal/astarte*: No such file or directory (os error 2)
=== test refs ===

$ cd /root/astrate-mule && ls internal/ && echo "=== realm tree ===" && ls internal/realm 2>/dev/null; echo "=== who uses DetailBadRequest ==="; rg -ln "DetailBadRequest" --glob '*.go'
appengine
auth
broker
config
engine
flow
flowapi
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil
=== realm tree ===
async_operation_test.go
dashboard_compat_test.go
detailed.go
detailed_required_encrypted_test.go
detailed_test.go
helpers_test.go
http.go
http_test.go
policy.go
policy_test.go
service.go
=== who uses DetailBadRequest ===
pkg/astarteapi/envelope.go
internal/swagger/swagger_test.go

$ cd /root/astrate-mule && rg -n "Bad request|Internal server error" internal/realm/http.go internal/httpx/*.go pkg/astarteapi/*.go | grep -v _test | head -30
pkg/astarteapi/envelope.go:34:	// (2026-08-24, verify batch #57): Phoenix renders "Bad request" — the
pkg/astarteapi/envelope.go:37:	DetailBadRequest = "Bad request"
pkg/astarteapi/envelope.go:57:	// "Internal server error" — the capital-S form was reconstructed.
pkg/astarteapi/envelope.go:58:	DetailInternalServerError = "Internal server error"

$ cd /root/astrate-mule && rg -n "Detail[A-Za-z]*\b" internal/realm/http.go | head -20; echo "=== badRequest helpers ==="; rg -n "func .*BadRequest|func .*InternalError|func WriteError" internal/httpx/*.go pkg/astarteapi/*.go | grep -v _test | head
151:		docs, err := a.svc.ListInterfacesDetailed(r.Context(), r.PathValue("realm"))
393:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
411:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
560:// validationDetail strips the ErrValidation sentinel prefix so the response
562:func validationDetail(err error) string {
=== badRequest helpers ===
pkg/astarteapi/envelope.go:163:func WriteError(w http.ResponseWriter, status int, detail string) error {
pkg/astarteapi/envelope.go:191:func WriteBadRequest(w http.ResponseWriter) error {

→ Read pkg/astarteapi/envelope.go [offset=155, limit=55]

$ cd /root/astrate-mule && sed -n '209,249p' pkg/astarteapi/envelope.go; echo "=== realm uses ==="; rg -c "astarteapi.WriteBadRequest|astarteapi.WriteInternalServerError" internal/realm/*.go
// WriteDeviceNotFound writes the canonical 404 "Device not found" envelope.
func WriteDeviceNotFound(w http.ResponseWriter) error {
	return WriteError(w, http.StatusNotFound, DetailDeviceNotFound)
}

// WriteInternalServerError writes the canonical 500 envelope.
func WriteInternalServerError(w http.ResponseWriter) error {
	return WriteError(w, http.StatusInternalServerError, DetailInternalServerError)
}

// DecodeData reads at most maxBytes bytes from r, unwraps the mandatory
// {"data": ...} request envelope, and unmarshals the "data" value into dst.
//
// It fails with an error wrapping ErrBodyTooLarge when the body exceeds
// maxBytes, with one wrapping ErrMissingData when the "data" key is absent or
// null, and with the underlying JSON error for malformed bodies (including
// trailing garbage after the top-level value). Sibling keys next to "data"
// are ignored, matching upstream parameter handling.
func DecodeData(r io.Reader, maxBytes int64, dst any) error {
	body, err := io.ReadAll(io.LimitReader(r, maxBytes+1))
	if err != nil {
		return fmt.Errorf("astarteapi: reading request body: %w", err)
	}
	if int64(len(body)) > maxBytes {
		return fmt.Errorf("astarteapi: %w: exceeds %d bytes", ErrBodyTooLarge, maxBytes)
	}
	var env struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("astarteapi: decoding request envelope: %w", err)
	}
	if len(env.Data) == 0 || bytes.Equal(env.Data, []byte("null")) {
		return fmt.Errorf("astarteapi: %w", ErrMissingData)
	}
	if err := json.Unmarshal(env.Data, dst); err != nil {
		return fmt.Errorf("astarteapi: decoding request data: %w", err)
	}
	return nil
}
=== realm uses ===
internal/realm/http.go:7

$ cd /root/astrate-mule && python3 - <<'EOF'
p='docs/api/astarte_realm_management_api.yaml'
src=open(p).read()
pairs=[
 ('          example:\n            errors:\n              detail: Bad Request\n',
  '          example:\n            errors:\n              detail: Bad request\n'),
 ('          example:\n            errors:\n              detail: Internal Server Error\n',
  '          example:\n            errors:\n              detail: Internal server error\n'),
]
for old,new in pairs:
    n=src.count(old)
    assert n==1, (n, old)
    src=src.replace(old,new)
open(p,'w').write(src)
lines=src.split('\n')
print('\n'.join(f'{i}: {l}' for i,l in enumerate(lines,1) if l.strip().startswith('detail:') and 'flat_detail' not in l))
EOF
152:                       detail: 'missing "interface_name"'
288:                       detail: 'missing "interface_name"'
341:                       detail: "Interface can't be deleted"
346:                       detail: "Interface can't be deleted since it's currently used"
434:                       detail: "trigger requires a name"
1193:             detail:
1265:               detail: Bad request
1275:               detail: Unauthorized
1285:               detail: Forbidden
1295:               detail: Not Found
1305:               detail: Device not found
1318:                   detail: Already exists
1323:                   detail: Interface name doesn't match the one in the interface json
1328:                   detail: Interface major version doesn't match the one in the interface json
1338:               detail: "invalid interface: missing \"interface_name\""
1348:               detail: Internal server error

Now the test that proves it, following the pairing precedent:
→ Read internal/swagger/swagger_test.go [limit=30, offset=1]

$ cd /root/astrate-mule && rg -n "func componentBlock" -A 22 internal/swagger/swagger_test.go
673:func componentBlock(t *testing.T, lines []string, key string) []string {
674-	t.Helper()
675-	start := -1
676-	for i, l := range lines {
677-		if l == key {
678-			start = i
679-			break
680-		}
681-	}
682-	if start < 0 {
683-		t.Fatalf("spec declares no components entry %q", key)
684-	}
685-	for i := start + 1; i < len(lines); i++ {
686-		if strings.TrimSpace(lines[i]) == "" {
687-			continue
688-		}
689-		if !strings.HasPrefix(lines[i], "      ") {
690-			return lines[start+1 : i]
691-		}
692-	}
693-	return lines[start+1:]
694-}
695-

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -504,8 +504,58 @@
 		}
 	}
 }
 
+// TestRealmManagementErrorDetailExamples pins the error-detail examples in the
+// realm management spec's components.responses to the frozen canonical strings
+// the wire emits (astarteapi/envelope.go), the same guard the pairing spec gets
+// from TestPairingErrorDetailExamples. The BadRequest and InternalServerError
+// examples carried the capital-R/capital-S spellings, which upstream Phoenix
+// never renders and no handler emits: every 400 on this surface goes through
+// astarteapi.WriteBadRequest and every 500 through WriteInternalServerError
+// (internal/realm/http.go), and the SDK/astartectl error paths match on the
+// canonical strings. The Conflict and ValidationError examples are
+// endpoint-specific details with no constant behind them, so they stay free-form.
+func TestRealmManagementErrorDetailExamples(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	const detailPrefix = "              detail: "
+
+	for _, tc := range []struct {
+		response string
+		want     string
+	}{
+		{"    BadRequest:", astarteapi.DetailBadRequest},
+		{"    Unauthorized:", astarteapi.DetailUnauthorized},
+		{"    Forbidden:", astarteapi.DetailForbidden},
+		{"    NotFound:", astarteapi.DetailNotFound},
+		{"    DeviceNotFound:", astarteapi.DetailDeviceNotFound},
+		{"    InternalServerError:", astarteapi.DetailInternalServerError},
+	} {
+		block := componentBlock(t, lines, tc.response)
+
+		got := ""
+		for _, l := range block {
+			if strings.HasPrefix(l, detailPrefix) {
+				got = strings.TrimPrefix(l, detailPrefix)
+				break
+			}
+		}
+		if got == "" {
+			t.Errorf("response %q carries no %q example", tc.response, "detail")
+			continue
+		}
+		if got != tc.want {
+			t.Errorf("response %q example detail = %q, want the canonical %q",
+				strings.TrimSpace(tc.response), got, tc.want)
+		}
+	}
+}
+
 // TestPairingDeviceIDEncodingDocumented guards that the pairing spec states the
 // exact wire form of a device identifier — the 22-character unpadded base64url
 // string deviceid.Parse accepts — instead of the "base64-encoded 128-bit"
 // wording, which read as standard base64 and let a generated client send a

$ cd /root/astrate-mule && gofmt -l . && go vet ./internal/swagger/... && go test ./internal/swagger/... 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/swagger	0.082s

$ cd /root/astrate-mule && cp docs/api/astrate_realm_management_api.yaml /tmp/opencode/rm.yaml.new && python3 - <<'EOF'
p='docs/api/astrate_realm_management_api.yaml'
s=open(p).read()
s=s.replace('detail: Bad request','detail: Bad Request').replace('detail: Internal server error','detail: Internal Server Error')
open(p,'w').write(s)
EOF
go test ./internal/swagger/ -run TestRealmManagementErrorDetailExamples 2>&1 | tail -10; cp /tmp/opencode/rm.yaml.new docs/api/astrate_realm_management_api.yaml; echo "--- restored ---"; go test ./internal/swagger/ -run TestRealmManagementErrorDetailExamples 2>&1 | tail -3
cp: cannot stat 'docs/api/astrate_realm_management_api.yaml': No such file or directory
ok  	github.com/astrate-platform/astrate/internal/swagger	0.021s
cp: cannot stat '/tmp/opencode/rm.yaml.new': No such file or directory
--- restored ---
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)

$ cd /root/astrate-mule && python3 - <<'EOF'
import subprocess
p='docs/api/astrate_realm_management_api.yaml'
orig=open(p).read()
assert 'detail: Bad request\n' in orig and 'detail: Internal server error\n' in orig
assert 'detail: Bad Request\n' not in orig and 'detail: Internal Server Error\n' not in orig

mut=orig.replace('detail: Bad request\n','detail: Bad Request\n').replace('detail: Internal server error\n','detail: Internal Server Error\n')
try:
    open(p,'w').write(mut)
    r=subprocess.run(['go','test','./internal/swagger/','-run','TestRealmManagementErrorDetailExamples','-count=1'],capture_output=True,text=True)
    print('WITHOUT FIX: exit',r.returncode)
    print('\n'.join(l for l in (r.stdout+r.stderr).split('\n') if 'detail' in l or 'FAIL' in l or 'ok ' in l))
finally:
    open(p,'w').write(orig)
print('restored ok:', open(p).read()==orig)
EOF
Traceback (most recent call last):
  File "<stdin>", line 3, in <module>
FileNotFoundError: [Errno 2] No such file or directory: 'docs/api/astrate_realm_management_api.yaml'

$ pwd; ls docs/api/ ; echo "---"; ls -d docs; python3 -c "import os; print(os.getcwd()); print(sorted(os.listdir('docs/api')))"
/root/astrate-mule
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
---
docs
/root/astrate-mule
['astarte_appengine_api.yaml', 'astarte_housekeeping_api.yaml', 'astarte_pairing_api.yaml', 'astarte_realm_management_api.yaml', 'astrate_native_api.yaml']

$ cd /root/astrate-mule && python3 - <<'EOF'
import subprocess, time
p='/root/astrate-mule/docs/api/astarte_realm_management_api.yaml'
def rd():
    for _ in range(10):
        try: return open(p).read()
        except FileNotFoundError: time.sleep(0.3)
    raise SystemExit('cannot read spec')
orig=rd()
assert 'detail: Bad request\n' in orig and 'detail: Internal server error\n' in orig
assert 'detail: Bad Request\n' not in orig and 'detail: Internal Server Error\n' not in orig
mut=orig.replace('detail: Bad request\n','detail: Bad Request\n').replace('detail: Internal server error\n','detail: Internal Server Error\n')
try:
    open(p,'w').write(mut)
    r=subprocess.run(['go','test','./internal/swagger/','-run','TestRealmManagementErrorDetailExamples','-count=1'],capture_output=True,text=True)
    print('WITHOUT FIX: exit',r.returncode)
    print('\n'.join(l for l in (r.stdout+r.stderr).split('\n') if 'detail =' in l or 'FAIL' in l or l.startswith('ok ')))
finally:
    open(p,'w').write(orig)
print('restored unchanged:', rd()==orig)
EOF
WITHOUT FIX: exit 1
--- FAIL: TestRealmManagementErrorDetailExamples (0.00s)
    swagger_test.go:552: response "BadRequest:" example detail = "Bad Request", want the canonical "Bad request"
    swagger_test.go:552: response "InternalServerError:" example detail = "Internal Server Error", want the canonical "Internal server error"
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.030s
FAIL
restored unchanged: True

$ cd /root/astrate-mule && python3 -c "
for f in ['.gitignore','docs/.gitignore']:
    try:
        print('==',f); print(open(f).read())
    except FileNotFoundError: print('missing',f)
"; echo "=== python yaml? ==="; python3 -c "import yaml; print(yaml.__version__)" 2>&1 | tail -1; echo "=== mkdocs? ==="; command -v mkdocs || echo "no mkdocs"; ls docs/site | head
== .gitignore
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

# Docs site copies (refreshed by docs/make sync)
/docs/site/DESIGN.md
/docs/site/ROADMAP.md
/docs/site/api/
/docs/site/swagger-ui/

# Local-only compose overrides
docker-compose.override.yml
__pycache__/

# bench harness artifacts (bench/README.md). Note bench/results/ is deliberately NOT here:
# a benchmark number without its host and its history is worthless, so results are evidence
# and evidence is versioned.
/bench/.astarte/
/bench/keys/
/bench/astrate/conf/
/bench/*.json
/bench/bench

# The mule host's tick transcript: local to whichever machine ran it, grows without
# bound, and got committed by accident on 2026-09-04.
.mule/cron.log

== docs/.gitignore
missing docs/.gitignore
=== python yaml? ===
ModuleNotFoundError: No module named 'yaml'
=== mkdocs? ===
no mkdocs
DESIGN.md
ROADMAP.md
api
appengine-api.md
architecture.md
compatibility.md
configuration-reference.md
contributing.md
data-modeling.md
deployment.md

$ cd /root/astrate-mule && sed -n '750,800p' internal/swagger/swagger_test.go; echo "=== swagger.go parse ==="; sed -n '1,70p' internal/swagger/swagger.go; echo "=== yaml dep ==="; python3 -c "print(open('go.mod').read())" | grep -i yaml; cat docs/requirements.txt
}
	}
	return false
}

func embeddedYAMLFilenames() ([]string, error) {
	var names []string
	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".yaml") {
			names = append(names, strings.TrimPrefix(path, "api/"))
		}
		return nil
	})
	sort.Strings(names)
	return names, err
}

func get(t *testing.T, url string) string {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s: status = %d, want %d", url, resp.StatusCode, http.StatusOK)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading %s: %v", url, err)
	}
	return string(b)
}
=== swagger.go parse ===
// Package swagger serves the embedded Swagger UI and OpenAPI YAML specs at
// /swagger/ and /api/ respectively.
package swagger

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	docs "github.com/astrate-platform/astrate/docs"
)

// Mount registers the /swagger and /api routes on the given mux.
// /swagger redirects to /swagger/index.html; /swagger/ serves the static UI;
// /api/ serves the OpenAPI YAML specs.
func Mount(mux *http.ServeMux) {
	uiRoot := mustSub(docs.SwaggerUI, "swagger-ui")
	apiRoot := mustSub(docs.APIYAML, "api")

	MountWithFS(mux, uiRoot, apiRoot)
}

// MountWithFS registers the /swagger and /api routes on the given mux from
// already-subtreed file systems. Mount computes the sub-trees from the
// embedded docs and delegates here; injecting fs.FS values makes the served
// roots testable and lets a broken embed fail fast instead of silently
// serving an empty tree.
func MountWithFS(mux *http.ServeMux, uiRoot, apiRoot fs.FS) {
	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
	})
	mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(uiRoot))))

	// Serve YAML files at /api/ so the relative ../api/*.yaml in index.html
	// resolves correctly when the page is loaded from /swagger/index.html.
	mux.Handle("GET /api/", http.StripPrefix("/api/", http.FileServer(http.FS(apiRoot))))
}

// mustSub returns fs.Sub(fsys, name), panicking if the sub-tree is absent so
// a broken docs embed fails fast instead of silently serving an empty tree.
// fs.Sub itself only errors on an invalid path (embed.FS has no Sub method),
// so the sub-root is also opened to catch an embed whose layout dropped the
// tree entirely.
func mustSub(fsys fs.FS, name string) fs.FS {
	sub, err := fs.Sub(fsys, name)
	if err != nil {
		panic(fmt.Sprintf("swagger: fs.Sub(%q): %v", name, err))
	}
	if _, err := sub.Open("."); err != nil {
		panic(fmt.Sprintf("swagger: %q missing from embedded fs: %v", name, err))
	}
	return sub
}

// Specs returns the list of available YAML spec filenames (without path prefix).
func Specs() []string {
	var names []string
	_ = fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".yaml") {
			names = append(names, strings.TrimPrefix(path, "api/"))
		}
		return nil
	})
	return names
}
=== yaml dep ===
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10

$ cd /root/astrate-mule && timeout 90 python3 -m pip install --quiet pyyaml 2>&1 | tail -5; python3 -c "import yaml;print('pyyaml', yaml.__version__)" 2>&1 | tail -1
/usr/bin/python3: No module named pip
ModuleNotFoundError: No module named 'yaml'

$ mkdir -p /tmp/opencode/yamlprobe && cd /tmp/opencode/yamlprobe && python3 -c "
open('go.mod','w').write('module yamlprobe\n\ngo 1.24\n\nrequire gopkg.in/yaml.v3 v3.0.1\n')
open('main.go','w').write('''package main

import (
	\"fmt\"
	\"os\"

	\"gopkg.in/yaml.v3\"
)

func main() {
	for _, p := range os.Args[1:] {
		b, err := os.ReadFile(p)
		if err != nil {
			fmt.Println(\"READ-FAIL\", p, err)
			os.Exit(1)
		}
		var doc map[string]any
		if err := yaml.Unmarshal(b, &doc); err != nil {
			fmt.Println(\"PARSE-FAIL\", p, err)
			os.Exit(1)
		}
		var node yaml.Node
		if err := yaml.Unmarshal(b, &node); err != nil {
			fmt.Println(\"NODE-FAIL\", p, err)
			os.Exit(1)
		}
		resp, _ := doc[\"components\"].(map[string]any)[\"responses\"].(map[string]any)
		for _, name := range []string{\"BadRequest\", \"InternalServerError\"} {
			ex, _ := resp[name].(map[string]any)[\"content\"].(map[string]any)[\"application/json; charset=utf-8\"].(map[string]any)[\"example\"].(map[string]any)
			fmt.Printf(\"%s %s example detail = %#v\\n\", p, name, ex[\"errors\"].(map[string]any)[\"detail\"])
		}
		fmt.Println(\"OK\", p, \"openapi=\", doc[\"openapi\"], \"nodes=\", len(node.Content))
	}
}
''')
" && GOFLAGS=-mod=mod GOPROXY=off go run . /root/astrate-mule/docs/api/astrate_realm_management_api.yaml
# yamlprobe
./main.go:30:42: newline in string
./main.go:30:42: syntax error: unexpected newline in argument list; possibly missing comma or )
./main.go:31:53: newline in string

$ cd /tmp/opencode/yamlprobe && python3 - <<'PYEOF'
src = '''package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	for _, p := range os.Args[1:] {
		b, err := os.ReadFile(p)
		if err != nil {
			fmt.Println("READ-FAIL", p, err)
			os.Exit(1)
		}
		var doc map[string]any
		if err := yaml.Unmarshal(b, &doc); err != nil {
			fmt.Println("PARSE-FAIL", p, err)
			os.Exit(1)
		}
		var node yaml.Node
		if err := yaml.Unmarshal(b, &node); err != nil {
			fmt.Println("NODE-FAIL", p, err)
			os.Exit(1)
		}
		comps, _ := doc["components"].(map[string]any)
		resp, _ := comps["responses"].(map[string]any)
		for _, name := range []string{"BadRequest", "InternalServerError"} {
			ex, _ := resp[name].(map[string]any)["content"].(map[string]any)["application/json; charset=utf-8"].(map[string]any)["example"].(map[string]any)
			errs, _ := ex["errors"].(map[string]any)
			fmt.Printf("%s %s example detail = %q\\n", name, name, errs["detail"])
		}
		fmt.Println("OK", p, "openapi =", doc["openapi"], "top-level nodes =", len(node.Content))
	}
}
'''
open('main.go','w').write(src)
PYEOF
GOFLAGS=-mod=mod GOPROXY=off go run . /root/astrate-mule/docs/api/astrate_realm_management_api.yaml
READ-FAIL /root/astrate-mule/docs/api/astrate_realm_management_api.yaml open /root/astrate-mule/docs/api/astrate_realm_management_api.yaml: no such file or directory
exit status 1

$ cd /tmp/opencode/yamlprobe && for i in 1 2 3; do GOFLAGS=-mod=mod GOPROXY=off go run . /root/astrate-mule/docs/api/astarte_realm_management_api.yaml 2>&1 | tail -2; sleep 1; done
InternalServerError InternalServerError example detail = "Internal server error"
OK /root/astrate-mule/docs/api/astarte_realm_management_api.yaml openapi = 3.0.3 top-level nodes = 1
InternalServerError InternalServerError example detail = "Internal server error"
OK /root/astrate-mule/docs/api/astarte_realm_management_api.yaml openapi = 3.0.3 top-level nodes = 1
InternalServerError InternalServerError example detail = "Internal server error"
OK /root/astrate-mule/docs/api/astarte_realm_management_api.yaml openapi = 3.0.3 top-level nodes = 1

$ cd /root/astrate-mule && make -C docs build 2>&1 | tail -5; echo "=== exit above ==="; ls docs/site/api/
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'
=== exit above ===
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

$ cd /root/astrate-mule && python3 -c "
for l in open('docs/site/api/astarte_realm_management_api.yaml').read().split('\n'):
    if l.strip() in ('detail: Bad request','detail: Internal server error'): print('site copy:', l)
"; echo "=== gate ==="; gofmt -l . ; go vet ./... && echo "vet ok"
site copy:               detail: Bad request
site copy:               detail: Internal server error
=== gate ===
vet ok

$ cd /root/astrate-mule && timeout 900 go test -race ./... 2>&1 | grep -v "^ok\|no test files" | head -20; echo "EXIT=${PIPESTATUS:-?}"
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/cmd/astrate	0.072s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/examples/flow-container-echo	0.060s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine	0.083s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine/channels	0.077s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine/stream	0.104s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/auth	0.064s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
EXIT=?

$ cd /root/astrate-mule && timeout 1200 go test ./... 2>&1 | grep -v "^ok\|no test files" | head -20; echo "done"
done

$ cd /root/astrate-mule && python3 -c "
import re
s=open('.mule/todo.md').read()
print('lines:', len(s.split('\n')))
for i,l in enumerate(s.split('\n'),1):
    if re.search('capital|detail|Bad request|Internal', l, re.I): print(i, l)
print('--- tail ---')
print('\n'.join(s.split('\n')[-8:]))
"
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_0fde2d65b001O0jnzMXkSmrAf5

207 - [x] broker-offlineacl-entry-eviction [auto]: `offlineACL.entries` (internal/broker/aclhook.go:100-162) is append-only — `ownershipOf` upserts one entry per CN ever ACL-checked offline and nothing ever deletes, so a long-lived broker grows the map unboundedly across device churn even after the 10s `offlineACLCacheTTL` passes (the entry is kept and re-stamped, not reaped). Evict stale entries lazily on access (drop an entry whose `loadedAt` predates some multiple of the TTL before inserting another, or cap+LRU), add an injectable `now func() time.Time` in the same style as `lifecycleHook.now`, and cover in a container-free unit test: seed N entries, advance the clock, access one, assert the map stays bounded.
208 - [!] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine_api.yaml — `retrieve_metadata` and `downsample_key` are parsed for every data read by `parseQueryOpts` (internal/appengine/http.go:624 `retrieve_metadata`, 633 `downsample_key`) yet appear on none of the read endpoints: device-scoped GET `{interface}` (yaml:583-606) and `{interface}/{path}` (yaml:623-649), by-alias (yaml:324-368, 369-417), and in-group (yaml:1108-1132, 1149-1176) — the parameter block on each lists only since/since_after/to/limit/downsample_to/sort/format/allow_bigintegers/allow_safe_bigintegers (components DataSince..DataAllowSafeBigIntegers, yaml:1360-1439). Add the two `$ref`s (new DataRetrieveMetadata and DataDownsampleKey components, or inline — mirror the existing style) to all six, and update the DataSort-independent `sort` default note only if verified. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
210 - [!] flow-mqtt-source-reconnect-recovery [auto]: `mqttSource` latches `lost` on paho's ConnectionLost handler (internal/flow/blocks/mqtt.go:185-187, gates in Emit/Process at 234-236 and 252-254) but nothing ever clears it, even though `newMQTTClient` enables `SetAutoReconnect(true)` (mqtt.go:54) and the docstring promises background auto-reconnect — so any transient broker blip permanently kills the source for the flow's lifetime ("connection lost" forever) while the flow stays `running`; `TestMQTTSource_ConnectionLost` (mqtt_test.go:190) only asserts the error surfaces, never recovery. Clear the latch on reconnect (a paho OnConnect handler, or gate on `client.IsConnected()` instead of a sticky flag), and add a container-free test that stops and restarts the embedded mochi broker and asserts Emit resumes. — BLOCKED: wrote nothing
212 - [x] flow-sort-bounded-buffer [auto]: `Sort` (internal/flow/blocks/sort.go:32-65) appends every message to an uncapped `buf` and only flushes while `len(buf) > 1 && buf[0].ts <= newest-windowUs` (sort.go:60) — a stream whose timestamps stay within `window_ms` of the newest (clustered bursts, or messages sharing one timestamp) never satisfies the condition, so the block emits nothing while the slice grows without bound, and the documented "newest is never emitted" mean the final message is held indefinitely; secondarily `window_ms * 1000` (sort.go:31) overflows for large configs and silently inverts the flush edge. Add a configurable max buffered count with a documented overflow policy (drop-oldest / error / force-flush) and a cap test feeding N same-timestamp messages; guard the multiply. Container-free.
213 - [!] flow-randomsource-span-overflow [auto]: `randomSource.next()` computes `rand.Int64N(s.maxInt-s.minInt+1)` (internal/flow/blocks/randomsource.go:171) whose argument overflows int64 to a non-positive value when min/max span more than the int64 range — config like `min:-9e18, max:9e18` passes the single `min <= max` check (randomsource.go:106-108) but makes `rand.Int64N` panic (it panics on n <= 0) inside the source pump goroutine (flow.go:205-239, no recover), crashing the whole process. Reject a span wider than `math.MaxInt64` at construction alongside the existing check, and add a construct/Emit test with the wide span. Container-free. — BLOCKED: lint failed: internal/flow/blocks/randomsource.go:113:40: G115: integer overflow conversion int64 -> uint64 (gosec)
214 - [x] flow-filter-key-contains-test [auto]: `Filter`'s `key_contains` rule (`strings.Contains`, internal/flow/blocks/transform.go:69) has zero test coverage — `transform_test.go` only exercises `key_prefix`, `type`, unknown-type and the at-least-one construct rule, so a regression silently dropping the Contains check would pass the whole suite. Add substring-match, substring-absent, and prefix+contains combined rows to the filter tests. Container-free.
217 - [x] docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` (yaml:231-232) in docs/api/astarte_realm_management_api.yaml as a oneOf — both handlers go through `writeInterfaceError`, which answers three distinct 422 bodies: the flat ErrorDetail (`ErrValidation` via `validationDetail`, http.go:392-393), the nested violations changeset envelope (`writeViolations`, http.go:358-359 + 434-533, e.g. `{"errors":{"description":["should be at most 1000 character(s)"]}}` and the aligned full-length `mappings` array, http_test.go:480-501), and the named FieldErrors envelope for `ErrMaximumDatabaseRetentionExceeded` (`{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, http.go:355-357, http_test.go:282-320); today only the flat ValidationError is referenced. Mirror the createTrigger oneOf pattern (yaml:344-368). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
218 - [x] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "realm: validation failed: interface definition is invalid"`, but `validationDetail` strips the `realm: validation failed: ` prefix on the wire (http.go:562-568) and the parser's real messages are `invalid interface: ...` (pkg/interfaceschema/parse.go:38, wrapped at service.go:163); the example should show the stripped, real message. Same class as the already-fixed docs-sync-hk-validation-example (housekeeping http.go:251-258). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
219 - [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
220 - [x] swagger-sub-failfast [auto]: make `Mount` fail fast (panic) when `fs.Sub` over the `docs.SwaggerUI`/`docs.APIYAML` embeds errors — `uiRoot, _ := fs.Sub(...)` and `apiRoot, _ := fs.Sub(...)` (internal/swagger/swagger.go:17-18) silently serve an empty tree → 404s at /swagger/ and /api/ with no log the day the docs embed layout changes; add an fs-injecting internal variant (`MountWithFS(mux, uiRoot, apiRoot fs.FS)` called by `Mount` after the subs) so a broken-FS unit test in swagger_test.go can prove the fail-fast, keeping the public `Mount(mux)` signature and its single main.go:414 caller untouched. Container-free.
221 - [ ] flow-boot-resumes-stopped-durable-flows [legion] [auto]: RehydrateAutoRestart (cmd/astrate/main.go:216) restarts every flow with auto_restart=true whatever its last status — `ListAutoRestartFlows` filters on `f.auto_restart = true` alone (internal/store/flows.go:125) and nothing ever clears the flag, which is written once at create (flowapi/http.go:141-146 -> service.go:318) and never on stop — so a durable flow stopped through the API is running again after the next boot, while the shutdown mark that says otherwise (MarkRunningFlowsStopped, main.go:497 -> internal/flowapi/service.go:885-901, which writes status="stopped") is never read by that query. The two halves disagree: service.go:325 documents "every durable flow", the shutdown call assumes "only what was running". Say which is the intent, then make them agree — either add `AND f.status = 'running'` to flows.go:125 (making the shutdown mark load-bearing) or delete the mark as dead code — and add a case in internal/store/flows_test.go pinning the chosen rule (a stopped auto_restart row must, or must not, come back). Needs the DB.
222 - [x] drain-per-stage-budget [auto]: `shutdown` spends a single 30s sctx (cmd/astrate/main.go:479) across three sequential stages — srv.Shutdown (484), b.Close (489), flowSvc.Manager().Shutdown (494) and MarkRunningFlowsStopped (497) — while `drainEngine` alone opens a second one (504), so the constant's own comment "bounds the whole graceful drain" (main.go:55-56) is wrong in both directions, and one slow in-flight HTTP request (the server sets no per-request deadline, 221) burns the whole first budget and leaves the flow stage on an expired context: StopFlow then returns at `f.router.Drain(ctx)` before step 3 releases block resources (internal/flow/flow.go:265-273) and MarkRunningFlowsStopped skips every flow (flowapi/service.go:895-897), both surfacing as one log.Warn. Give each stage its own budget through an extracted helper (shutdownTimeout has to become an injectable var for the test) and pin it container-free: a second stage still receives a live context after a first stage exhausts its own.
226 - [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
227 - [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
228 - [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
229 - [x] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion Go — the device-delete leg added by docs-sync-rm-delete-device-async-operation-param (register a device, `DELETE /devices/{id}?async_operation=…` answers 204 and the row is gone on both values) compiles clean under `go vet -tags integration ./internal/realm/` but has never been executed, because the Pi has no database; until it runs, TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go is the only assertion of that leg that has actually been verified.
230 - [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
231 - [x] forward-static-header-validation [auto]: `New` (internal/engine/forward/http.go:38-71) validates the URL and the method at boot precisely because of the rule its own comment states at 42-45 — "an unusable endpoint must fail here rather than surface once per delivery" — but `cfg.StaticHeaders` is stored unvalidated (http.go:68) and applied per delivery with `req.Header.Set` (http.go:110-112), and net/http checks header names and values only at write time. Measured on the real net/http path with a throwaway program: a name of `X Bad Name` or `X:Foo` fails per delivery with `net/http: invalid header field name`, and a value containing a newline fails with `net/http: invalid header field value for "X-Foo"` (a plain space in a value is legal). So one typo in `triggers.forward.static_headers` makes the process boot clean, accept triggers, and then count `astrate_engine_trigger_deliveries_total{outcome="failed"}` and log "custom trigger action failed" (internal/engine/triggers/actions.go:499) on every custom action forever. Validate in `New`: the RFC 7230 field-name grammar is exactly the token grammar `validMethod` (http.go:77-95) already implements, so reuse that predicate per name (rename it if you like) and add a CR/LF rejection per value. Add a table case to `TestNewRejectsBadConfig` (internal/engine/forward/http_test.go:271) asserting `New` returns an error for a bad name and a bad value and that no request is attempted; this test fails today because `New` returns nil error for both.
232 - [x] forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
233 - [x] forward-status-error-body [auto]: a non-2xx forward returns `forward: status %d` and nothing else (internal/engine/forward/http.go:121-123) after copying the **entire** response body into `io.Discard` with no limit (http.go:117-120) — the one thing an operator needs to debug a 500 from the bus is destroyed, and the sibling webhook request in the same codebase does bound its drain with `io.LimitReader(resp.Body, 1<<20)` (internal/engine/triggers/actions.go:626), so the two near-identical request paths have drifted on exactly that line. Read a bounded prefix (e.g. 512 bytes) of the body into the returned error and bound the discard the way actions.go:626 does. Add a case beside `TestStatusTable` (internal/engine/forward/http_test.go:164) where the handler answers 500 with body `{"error":"boom"}` and assert the error mentions both `500` and `boom`; it fails today because the error carries only the status code.
234 - [x] forward-envelope-bytes-test [auto]: `marshalEnvelope` (internal/engine/forward/envelope.go:19-29) has no direct test — `rg marshalEnvelope` finds exactly two callers (http.go:99, nats.go:45) and no test file — while envelope.go:6-8 makes the strong claim that "every Forwarder implementation in this package must produce byte-identical envelopes for the same inputs". Both suites check the shape only by unmarshalling into `bodyShape` (http_test.go:24-29), which cannot see key order or the exact bytes and would pass unchanged if a field were renamed on both sides. Add internal/engine/forward/envelope_test.go with a table pinning the exact output for (nil, nil) → `{"realm":"r","trigger":"t","action":null,"event":null}`, for the empty-non-nil pair the rule at 16-18 is about, and for a valued pair; container-free, no HTTP server needed.
235 - [x] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.
236 - [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml — the string `"403"` appears nowhere in the file and `components` defines no `Forbidden` response, yet every route is wrapped by `a.require = mw.RequireRealm(auth.ClaimAppEngine)` (internal/appengine/http.go:35), which answers `403 Forbidden` for a verified a_aea JWT whose grants do not authorize the method+path (`astarteapi.WriteForbidden`, internal/auth/middleware.go:90; `DetailForbidden`, pkg/astarteapi/envelope.go:41; pinned today by internal/appengine/http_test.go:184-185). Add a `Forbidden` response component copied from the realm-management one (docs/api/astarte_realm_management_api.yaml:1277-1285 — `ErrorDetail` schema, example `errors.detail: Forbidden`) and `$ref` it from all 32 operations, which is exactly the parity docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go:123-170) pins; add the same docs-half test for the AppEngine spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
237 - [x] docs-sync-ae-patch-by-alias-409 [auto]: add the missing `409` response to `patchDeviceByAlias` in docs/api/astarte_appengine_api.yaml (yaml:241-291, currently 200/400/401/404/422/500) — the handler resolves the alias then calls the same `applyPatch` as the device- and group-scoped PATCH, and a rename to an alias already in use returns `ErrAliasAlreadyInUse` (internal/appengine/service.go:337) which `writeError` maps to `409 ConflictAliasInUse` (internal/appengine/http.go:650-651); the other two PATCH operations already document that response thanks to docs-sync-appengine-query-params-status and docs-sync-appengine-group-patch-status, so this is the last PATCH of the three with a hole. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
238 - [x] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml (yaml:939-980, currently 201/400/401/404/409/500) — a well-formed body whose `device_id` does not parse is rejected as a field error (`FieldErrors{"device_id": {"is not a valid device id"}}`, internal/appengine/service.go:631-634) and answered 422 with the `FieldErrorsDetail` body (`astarteapi.WriteFieldErrors`, internal/appengine/http.go:647), so `$ref` the existing `ValidationErrors` response (yaml:1707-1716) rather than `BadRequest` — the documented 400 only covers a body `DecodeData` cannot parse (internal/appengine/http.go:488-491). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
239 - [x] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data DELETE operations in docs/api/astarte_appengine_api.yaml (`deleteDataByAlias` yaml:519-547, `deleteData` yaml:756-784, `deleteDataInGroup` yaml:1285-1313 — each currently 204/401/404/405/500) — unsetting a path that matches no endpoint mapping returns `engine.ErrPathNotFound` (`UnsetServerProperty`, internal/engine/serverdata.go:241-244) which `writeError` answers as `400 "Endpoint not found"` (internal/appengine/http.go:679-680), the same 400 the PUT/POST twins already document via `BadRequest`; internal/appengine/writeerror_test.go pins the mapping. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
240 - [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 200/401/500) — a `from_token` that does not parse as a device id is refused with `ErrValidation: invalid cursor` (internal/appengine/service.go:150-155) which `writeError` maps to 422 with the flat `{"errors":{"detail":"invalid cursor"}}` body, i.e. the `ErrorDetail` shape and not `FieldErrorsDetail`, so `$ref` the same response the group twin uses for the identical parameter (`ValidationErrors`, yaml:934-935) only after fixing its example, or add a dedicated one. Extend the `from_token` parameter description (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
241 - [x] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:556) and `InternalServerError` shows `detail: Internal Server Error` (yaml:638), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37, golden fixture pkg/astarteapi/testdata/error_bad_request.json) and `"Internal server error"` (envelope.go:58, testdata/error_internal.json), both measured on upstream 1.2.0 (envelope.go:33-36, 55-57). The other five detail examples on this file (`Not Found`, `Device not found`, `Unauthorized`, `Forbidden`, `Too Many Requests`) are already correct — change only those two. The identical wrong pair is in the other four specs (appengine 1732/1792, housekeeping 362/385/477, realm_management 1265/1348, native 466); leave those for a run of their surface. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
242 - [x] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pairing_api.yaml, which shows `version: "0.1.0-test"` (yaml:252) — a value no build can produce. The field is filled from `pairing.Config.Version` (internal/pairing/service.go:302), wired from `main.version` (cmd/astrate/main.go:151), whose default is `"0.1.0-dev"` (main.go:53) and which is otherwise set only via `-ldflags "-X main.version=..."`. Use `"0.1.0-dev"`, matching the four compat-version examples docs-sync-native-version-value already corrected to that (docs/api/astrate_native_api.yaml:203, 223, 243, 263). Do NOT document the service's own fallback `DefaultVersion = "0.1.0-astrate"` (service.go:86-88): it is unreachable from `main`, which always passes a non-empty value. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
243 - [x] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` path parameter (docs/api/astarte_pairing_api.yaml:360) and `RegisterRequest.hw_id` (yaml:373) — but `deviceid.Parse` accepts only the 22-character **unpadded base64url** form, in strict mode, and rejects the standard alphabet's `+` and `/` and any `=` padding with an error wrapping `ErrInvalid` (pkg/deviceid/deviceid.go:26, 37, 39-60; upstream parity cited at 36). Measured: with the same `encoding` the parser uses, `dT6hS2W9TT6LEnP25ks+lg` and `.../lg` fail at byte 19 and the padded 23-character form fails at byte 21, so a client generated from the current spec can send a well-formed standard-base64 ID and get 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}` (internal/pairing/http.go:307-309). State the exact encoding and the 22-character length in both descriptions and add a `pattern`/`minLength`/`maxLength` to `hw_id`. The existing `hw_id: dT6hS2W9TT6LEnP25ks_lg` example (yaml:123) is valid and stays. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
244 - [x] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in docs/api/astarte_pairing_api.yaml (yaml:374-376), an unconstrained string described only as "Optional Astrate extension specifying the device's initial payload format" — so a generated client can send any string, while `Service.Register` accepts only `"bson"` or `"json"` and rejects everything else with 422 `{"errors":{"initial_payload_format":["is invalid"]}}` (internal/pairing/service.go:179-181 -> internal/pairing/http.go:310-312; the allowed set is stated at service.go:173). Add `enum: [bson, json]` plus a note that absent/empty means unset, and add the three missing `ValidationError` examples while you are in the component: `initial_payload_format: ["is invalid"]` (http.go:310-312), `csr: ["can't be blank"]` (http.go:173-177) and `client_crt: ["can't be blank"]` (http.go:253-257) — the component currently shows only the four hw_id/error_name rows (yaml:594-618) and the `csr`/`client_crt` blank-field cases are what a device actually hits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
245 - [x] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` description in docs/api/astarte_pairing_api.yaml (yaml:153), "Removes a device from the realm. Returns 204 on success.", which reads as a device-and-data deletion. The row and every byte of its stored data survive: `store.UnregisterDevice` only clears the credential trail (`credentials_secret_hash = ''`, `first_credentials_request = NULL`, `cert_serial = NULL`, `cert_aki = NULL`) and flips `status` back to `'registered'` (internal/store/devices.go:93-110), and the service docstring says so outright ("its data is kept", internal/pairing/service.go:227-229). State the real semantics: the device becomes registrable again, its interfaces, datastream data and group memberships are retained, and a second DELETE is a 204 no-op (the row is still there, so RowsAffected is 1, devices.go:106-108) — not the documented 404. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
248 - [x] compat-note-v1.4-rc.6 [auto]: record in `.mule/for-giulio.md` that v1.4.0-rc.6 is wire-inert for Astrate, as a delta on the open v1.3.2/v1.3.4 wording proposal rather than a re-derivation. Measured, the whole non-FDO delta is inert: DUP RPC-server reliability is OTP supervisor/GC tuning; the Dashboard's new `unknown_status` is computed client-side from `last_connection && !last_disconnection && !isConnected` (models/Device/index.ts:213-215) and both timestamps are already on Astrate's device shape (internal/appengine/service.go:118), so a device whose disconnection was never recorded reads as "Unknown status" rather than "Disconnected"; the `on_device_registered` dispatch moved from `Agent.register_device` to `Engine.register_device` (pairing/engine.ex:157-165, 186-193) so it also covers the FDO path, but the Agent path is byte-for-byte the same rule and Astrate already fires it on every successful registration (internal/pairing/service.go:221-223); the `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"); and `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma. So the proposal's "v1.4.0 is still rc.5-only" reference becomes rc.6, and its "FDO authentication (pairing, disabled by default)" line understates the 1.4 line, which ships a full ownership-voucher surface Astrate does not emulate (#78, parked). Do not edit docs/COMPATIBILITY.md and do not restate the standing proposal.
250 - [x] container-stop-deadline [auto]: make `cliInstance.Stop` (internal/flow/blocks/container/docker.go:167-175) derive its 15s bound whenever the passed context carries **no deadline**, not only when it is nil — both cleanup callers pass `context.Background()` (docker.go:122, the "Best-effort cleanup so we do not leave orphans on mapping failure" path, and block.go:121, the not-ready cleanup in `New`), so a wedged docker daemon leaves an unbounded `docker rm -f` inside `exec.CommandContext` holding flow instantiation, while `Block.Stop` (block.go:310-311) does the same job correctly with an explicit 15s timeout. Measured with the package's own `Run` hook: the `rm` invocation arrives with `ctx.Deadline()` unset on both paths. Add a case asserting `inst.Stop(context.Background())` reaches the injected `Run` with a deadline set (it does not today) and that the hostPort-failure path still issues `rm -f <id>` (docker.go:120-124, never executed by any test today). Container-free.
251 - [!] container-timeout-bounds [auto]: bound `timeout_ms` and `ready_timeout_ms` in `parseConfig` (internal/flow/blocks/container/block.go:205-219) — both are only checked `> 0` and then multiplied unchecked (`time.Duration(n) * time.Millisecond`), which wraps: measured `ready_timeout_ms: 18446744073709` gives `-551.616µs`, so `New` fails with "container: not ready after wait" without ever issuing a probe (the wrapped negative deadline fails the loop guard at httpbridge.go:64), while `ready_timeout_ms: 9223372036854` gives 2562047h47m16s and `New` never returns (still polling `/healthz` after 3s in the probe). The two consumers of a non-positive duration also disagree — `Bridge.timeout()` (httpbridge.go:39-44) treats it as "use the 5s default" while net/http treats a non-positive `Client.Timeout` as no deadline at all (measured: a client with `Timeout: -551µs` returned err=nil after a 3s handler). Reject a value that would overflow or exceed a documented maximum, naming the key and the value in the error, and add table cases asserting `New` refuses both (it accepts them today). Nothing bounds these on the wire either: internal/flow/blocks/schema.go:19 advertises both as bare `"type":"integer"`. Container-free. — BLOCKED: TIMEOUT after 1201s — task too big, split it
252 - [x] container-parseconfig-rules-test [auto]: pin the four `parseConfig` rules that no test touches (internal/flow/blocks/container/block.go:188-219 — all `count=0` in `go test -coverprofile`, which is 71.2% for the package): `config` must be a JSON object (192-193), `port` an integer in 1–65535 including the 0 / 65536 / non-integer rejections (197-202), and `timeout_ms` and `ready_timeout_ms` positive (205-210, 215-217) — beside the existing `TestConstructor_RequiresImage` (block_test.go:99), which is the only config case the suite has. The port rule is the one with teeth: it is what keeps `127.0.0.1::99999` out of the `docker run` argv (docker.go:93). Container-free, pure.
253 - [!] container-response-cap-test [auto]: pin `Bridge.RoundTrip`'s 1 MiB response cap (internal/flow/blocks/container/httpbridge.go:133-140 — `io.LimitReader(resp.Body, maxBody+1)` plus the explicit rejection; lines 136, 138.16 and 138.40 are all uncovered), the only bound between a tenant-supplied container and process memory on the message path: assert a handler answering exactly `maxBody+1` bytes is rejected with the "exceeds" error and that exactly `maxBody` bytes still round-trip, driving a `Bridge{MaxBodyBytes: N}` so the fixture stays small. Also worth one line in the same change: the cap is absent from the operator-facing catalog entry (internal/flow/blocks/info.go:116) and from the schema, so an author whose container emits 2 MiB learns about it only from a runtime error. Container-free. — BLOCKED: lint failed: internal/flow/blocks/container/block_test.go:500:2: redefines-builtin-id: redefinition of the built-in function cap (rev
254 - [ ] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, housekeeping and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
256 - [ ] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
257 - [ ] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstrained string described only as "The device hardware ID." — `Service.DeleteDevice` runs the path value through `deviceid.Parse` (internal/realm/service.go:118-121), which accepts exactly 22 characters of unpadded base64url in strict mode and rejects padding, the standard alphabet's `+` and `/`, and non-canonical trailing bits (pkg/deviceid/deviceid.go:24-26, 33-37, 42-60); every rejection is folded into `store.ErrNotFound` (service.go:119-121) and answered `404 {"errors":{"detail":"Device not found"}}` (internal/realm/http.go:107-109), so a client that sends a well-formed standard-base64 or 23-character padded id is told the device does not exist rather than that its id is malformed. State the encoding and the 22-character length and add `pattern`/`minLength`/`maxLength`, exactly as docs-sync-pairing-deviceid-base64url did for the pairing twin. Do NOT encode the strict trailing-bits rule in the pattern: whether it is a deviation worth recording is still open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
258 - [ ] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_realm_management_api.yaml — `data` is declared as a bare `type: object` with no properties (yaml:250-261) although the handler parses it with the same `interfaceschema.ParseInterfaceCanonical` as `installInterface` (internal/realm/service.go:196 vs 161) and then checks the body's `interface_name`/`version_major` against the URL's `{name}`/`{major}` (service.go:203-208), so the update body — the one carrying every interface rule — is the only interface body in the spec with nothing to validate against. `$ref` `#/components/schemas/InterfaceDefinition` as `installInterface` does (yaml:115-117), and say in the operation description that the body must name the same interface and major as the URL (disagreement is 409, already documented) and that `version_minor` must increase (409). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
259 - [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.
--- tail ---
- [!] container-response-cap-test [auto]: pin `Bridge.RoundTrip`'s 1 MiB response cap (internal/flow/blocks/container/httpbridge.go:133-140 — `io.LimitReader(resp.Body, maxBody+1)` plus the explicit rejection; lines 136, 138.16 and 138.40 are all uncovered), the only bound between a tenant-supplied container and process memory on the message path: assert a handler answering exactly `maxBody+1` bytes is rejected with the "exceeds" error and that exactly `maxBody` bytes still round-trip, driving a `Bridge{MaxBodyBytes: N}` so the fixture stays small. Also worth one line in the same change: the cap is absent from the operator-facing catalog entry (internal/flow/blocks/info.go:116) and from the schema, so an author whose container emits 2 MiB learns about it only from a runtime error. Container-free. — BLOCKED: lint failed: internal/flow/blocks/container/block_test.go:500:2: redefines-builtin-id: redefinition of the built-in function cap (rev
- [ ] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, housekeeping and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
- [ ] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy alias fields to `InterfaceDefinition` and `InterfaceMapping` in docs/api/astarte_realm_management_api.yaml — both declare `additionalProperties: false` (yaml:935, 972) yet omit interface-level `quality` (alias of `ownership`) and `aggregate` (alias of `aggregation`) and mapping-level `path` (alias of `endpoint`), all three of which the parser accepts and canonicalises: decoded into dedicated fields (pkg/interfaceschema/parse.go:116-122, 136-138), normalised before any rule runs (parse.go:212-236), and re-encoded canonically by `ParseInterfaceCanonical` so the store never sees them (parse.go:165-170, 307-311). A spec-conformant client — or any generated client with this schema baked in — therefore rejects a valid upstream-shaped install/update with a spurious 422 while the service would have accepted and stored it as `ownership`/`aggregation`/`endpoint`. Add the three as documented deprecated aliases (enums mirroring ownership/aggregation, `path` a plain string) or relax `additionalProperties`, and record the three rules the code enforces: `ownership` + `quality` is a violation (parse.go:218) and `aggregation` + `aggregate` is a violation (parse.go:227), while `endpoint` + `path` silently keeps `endpoint` (parse.go:307-311). None of the three may appear in a `required` list. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstrained string described only as "The device hardware ID." — `Service.DeleteDevice` runs the path value through `deviceid.Parse` (internal/realm/service.go:118-121), which accepts exactly 22 characters of unpadded base64url in strict mode and rejects padding, the standard alphabet's `+` and `/`, and non-canonical trailing bits (pkg/deviceid/deviceid.go:24-26, 33-37, 42-60); every rejection is folded into `store.ErrNotFound` (service.go:119-121) and answered `404 {"errors":{"detail":"Device not found"}}` (internal/realm/http.go:107-109), so a client that sends a well-formed standard-base64 or 23-character padded id is told the device does not exist rather than that its id is malformed. State the encoding and the 22-character length and add `pattern`/`minLength`/`maxLength`, exactly as docs-sync-pairing-deviceid-base64url did for the pairing twin. Do NOT encode the strict trailing-bits rule in the pattern: whether it is a deviation worth recording is still open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_realm_management_api.yaml — `data` is declared as a bare `type: object` with no properties (yaml:250-261) although the handler parses it with the same `interfaceschema.ParseInterfaceCanonical` as `installInterface` (internal/realm/service.go:196 vs 161) and then checks the body's `interface_name`/`version_major` against the URL's `{name}`/`{major}` (service.go:203-208), so the update body — the one carrying every interface rule — is the only interface body in the spec with nothing to validate against. `$ref` `#/components/schemas/InterfaceDefinition` as `installInterface` does (yaml:115-117), and say in the operation description that the body must name the same interface and major as the URL (disagreement is 409, already documented) and that `version_minor` must increase (409). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.

$ cd /root/astrate-mule && python3 - <<'EOF'
for f in ['astarte_appengine_api.yaml','astarte_housekeeping_api.yaml','astrate_native_api.yaml']:
    p='docs/api/'+f
    lines=open(p).read().split('\n')
    print('==',f)
    for i,l in enumerate(lines,1):
        s=l.strip()
        if s in ('detail: Bad Request','detail: Internal Server Error'):
            # find enclosing component
            comp=None
            for j in range(i-1,0,-1):
                if lines[j].startswith('    ') and lines[j].rstrip().endswith(':') and not lines[j].startswith('     '):
                    comp=lines[j].strip(); break
            print(f'  {i}: {s}   (component {comp})')
EOF
== astarte_appengine_api.yaml
  1732: detail: Bad Request   (component BadRequest:)
  1792: detail: Internal Server Error   (component InternalServerError:)
== astarte_housekeeping_api.yaml
  362: detail: Bad Request   (component ErrorDetail:)
  385: detail: Bad Request   (component BadRequest:)
  477: detail: Internal Server Error   (component InternalServerError:)
== astrate_native_api.yaml
  466: detail: Internal Server Error   (component InternalServerError:)

$ cd /root/astrate-mule && python3 - <<'EOF'
lines=open('docs/api/astrate_housekeeping_api.yaml').read().split('\n')
print('\n'.join(f'{i}: {l}' for i,l in enumerate(lines[340:395],341)))
print('--- native details ---')
n=open('docs/api/astrate_native_api.yaml').read().split('\n')
print('\n'.join(f'{i}: {l}' for i,l in enumerate(n,1) if 'detail' in l))
EOF
Traceback (most recent call last):
  File "<stdin>", line 1, in <module>
FileNotFoundError: [Errno 2] No such file or directory: 'docs/api/astrate_housekeeping_api.yaml'

$ cd /root/astrate-mule && python3 - <<'EOF'
import time
def rd(p):
    for _ in range(15):
        try: return open(p).read().split('\n')
        except FileNotFoundError: time.sleep(0.2)
    raise SystemExit('unreadable '+p)
lines=rd('/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml')
print('\n'.join(f'{i}: {l}' for i,l in enumerate(lines[340:392],341)))
print('--- native detail lines ---')
n=rd('/root/astrate-mule/docs/api/astrate_native_api.yaml')
print('\n'.join(f'{i}: {l}' for i,l in enumerate(n,1) if 'detail' in l))
EOF
341:           nullable: true
342:           description: |
343:             Maximum datastream storage retention in seconds. Null clears the
344:             retention. `0` clears it as well: 0 is folded to unset, so it cannot
345:             ask for a literal zero-second retention and the realm is left
346:             unlimited (upstream parity, measured on Astarte v1.2.0).
347: 
348:     # ── Error envelopes ──
349: 
350:     ErrorDetail:
351:       type: object
352:       required: [errors]
353:       properties:
354:         errors:
355:           type: object
356:           required: [detail]
357:           properties:
358:             detail:
359:               type: string
360:               description: Human-readable error detail string.
361:           example:
362:             detail: Bad Request
363: 
364:     FieldErrorsDetail:
365:       type: object
366:       required: [errors]
367:       properties:
368:         errors:
369:           type: object
370:           additionalProperties:
371:             type: array
372:             items:
373:               type: string
374:           description: Per-field error messages (changeset shape).
375: 
376:   responses:
377:     BadRequest:
378:       description: Bad Request — malformed or missing data envelope.
379:       content:
380:         application/json; charset=utf-8:
381:           schema:
382:             $ref: "#/components/schemas/ErrorDetail"
383:           example:
384:             errors:
385:               detail: Bad Request
386: 
387:     Unauthorized:
388:       description: Unauthorized — missing or invalid JWT.
389:       content:
390:         application/json; charset=utf-8:
391:           schema:
392:             $ref: "#/components/schemas/ErrorDetail"
--- native detail lines ---
432:           required: [detail]
434:             detail:
446:               detail: Unauthorized
456:               detail: Forbidden
466:               detail: Internal Server Error

$ cd /root/astrate-mule && python3 - <<'EOF'
import time
p='/root/astrate-mule/.mule/todo.md'
for _ in range(15):
    try:
        s=open(p).read(); break
    except FileNotFoundError: time.sleep(0.2)
if not s.endswith('\n'): s+='\n'
s += "- [ ] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` shows `detail: Bad Request` (yaml:1732) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1792), while the wire emits the frozen lower-r/lower-s constants `\"Bad request\"` and `\"Internal server error\"` (pkg/astarteapi/envelope.go:37, 58 via `WriteBadRequest`/`WriteInternalServerError`, golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json); every 400/500 on that surface goes through the same astarteapi writers. The other detail examples on the file are already correct — change only those two — and add the docs-half assertion `TestAppEngineErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go, which pins the pairing and realm-management halves. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.\n"
s += "- [ ] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDetail` schema-level example (yaml:361-362) and the `BadRequest` response (yaml:385) show `detail: Bad Request`, and `InternalServerError` (yaml:477) shows `detail: Internal Server Error`, but the wire emits the frozen lower-r/lower-s constants `\"Bad request\"` and `\"Internal server error\"` (pkg/astarteapi/envelope.go:37, 58, via WriteBadRequest/WriteInternalServerError, pinned by pkg/astarteapi/testdata/error_bad_request.json and error_internal.json). While you are in `ErrorDetail`: its `example` sits at the schema level and shows a bare `{detail: ...}` where the envelope the wire sends is `{\"errors\":{\"detail\":...}}` — move it under `errors` (or drop it, as the realm-management and pairing `ErrorDetail` schemas carry no example) so a generated client is not shown the wrong shape. Add `TestHousekeepingErrorDetailExamples` beside the pairing and realm-management halves in internal/swagger/swagger_test.go. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.\n"
s += "- [ ] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `\"Internal server error\"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.\n"
open(p,'w').write(s)
lines=open(p).read().split('\n')
print('\n'.join(l[:120] for l in lines[-4:]))
print('total lines', len(lines))
EOF
- [ ] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appe
- [ ] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_ho
- [ ] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_n

total lines 263

$ cd /root/astrate-mule && gofmt -l . ; go vet ./internal/swagger/... && go test -count=1 ./internal/swagger/... ./pkg/astarteapi/... 2>&1 | tail -5
ok  	github.com/astrate-platform/astrate/internal/swagger	0.108s
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	0.084s

Done: fixed the two wrong error-detail examples in the realm-management spec (`BadRequest` → `Bad request` yaml:1265, `InternalServerError` → `Internal server error` yaml:1348) and pinned them with a docs-half test mirroring the pairing one.
Files: docs/api/astrate_realm_management_api.yaml (2 lines), internal/swagger/swagger_test.go (`TestRealmManagementErrorDetailExamples`, 6 constant-backed responses), .mule/todo.md (3 follow-up lines)
Verified: `go vet ./...` + `go test ./...` -> pass (`gofmt -l .` clean); `go test -race ./...` unusable here (ThreadSanitizer: unsupported VMA range 39 — expected on this host); negative check — reverting only the two YAML lines makes the new test fail with `example detail = "Bad Request", want "Bad request"` and `"Internal Server Error"`/`"Internal server error"`, restored afterwards; YAML still loads (throwaway `gopkg.in/yaml.v3` probe in /tmp: parses, `openapi = 3.0.3`, parsed examples equal the two constants); `make -C docs build` -> **fails at `mkdocs build`**: mkdocs absent and there is no `pip` module on this host, so only the `sync` half ran (the gitignored `docs/site/api/` copy carries the fix) — the Swagger-referenced YAML was verified by the embed/HTTP path (`go test ./internal/swagger/` fetches every spec over HTTP) instead
Unsure: nothing about the code; the only unverified leg is the real mkdocs render, which this host cannot install (PyYAML/mkdocs both unavailable without pip)
Follow-ups: 3 lines appended to .mule/todo.md — `docs-sync-ae-error-example-capitalisation` (appengine 1732/1792), `docs-sync-hk-error-example-capitalisation` (housekeeping 385/477 plus the schema-level `ErrorDetail` example at 361-362, which also shows a bare `{detail: …}` instead of `{"errors":{"detail":…}}`), `docs-sync-native-error-example-capitalisation` (native 466 only — that spec has no BadRequest component)
