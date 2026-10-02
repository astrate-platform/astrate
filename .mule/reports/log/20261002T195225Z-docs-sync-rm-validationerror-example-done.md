slug: docs-sync-rm-validationerror-example
verdict: done
at:  8a09a1a
ran: 2026-10-02T19:52:25Z on DietPi in 732s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/task.md
→ Read .mule/MULE.md

$ ls /root/astrate-mule/docs/ && ls /root/astrate-mule/docs/api/
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

$ wc -l /root/astrate-mule/docs/api/*.yaml && ls /root/astrate-mule/docs/*.mk /root/astrate-mule/docs/Makefile 2>/dev/null
1853 /root/astrate-mule/docs/api/astarte_appengine_api.yaml
   477 /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
   686 /root/astrate-mule/docs/api/astarte_pairing_api.yaml
  1348 /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
   466 /root/astrate-mule/docs/api/astrate_native_api.yaml
  4830 total
/root/astrate-mule/docs/Makefile

→ Read docs/api/astarte_realm_management_api.yaml [offset=120, limit=200]
→ Read docs/api/astarte_realm_management_api.yaml [offset=1300, limit=49]
→ Read docs/api/astarte_realm_management_api.yaml [limit=130, offset=520]
→ Read docs/api/astarte_realm_management_api.yaml [limit=140, offset=750]

$ cd /root/astrate-mule && rg -n "can't be blank|is still used by trigger|must have at least one error handler|policy name must be" --glob '!*_test.go' | head -30
docs/api/astarte_pairing_api.yaml:630:                    - "can't be blank"
docs/api/astarte_pairing_api.yaml:660:                    - "can't be blank"
docs/api/astarte_pairing_api.yaml:666:                    - "can't be blank"
docs/api/astarte_housekeeping_api.yaml:467:              detail: "realm_name can't be blank"
docs/api/astarte_appengine_api.yaml:1803:                - can't be blank
internal/housekeeping/http.go:189:			map[string][]string{"jwt_public_key_pem": {"can't be blank"}})
internal/housekeeping/service.go:138:		return nil, fmt.Errorf("%w: realm_name can't be blank", ErrValidation)
internal/housekeeping/service.go:141:		return nil, fmt.Errorf("%w: jwt_public_key_pem can't be blank", ErrValidation)
internal/housekeeping/service.go:206:		return fmt.Errorf("%w: jwt_public_key_pem can't be blank", ErrValidation)
internal/engine/triggers/actions.go:111:			fe.add("http_method", "can't be blank")
internal/engine/triggers/policy.go:113:		return nil, fmt.Errorf("triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]")
internal/engine/triggers/policy.go:116:		return nil, fmt.Errorf("triggers: policy must have at least one error handler")
internal/engine/triggers/match.go:321:		fe.add("on", "can't be blank")
internal/engine/triggers/match.go:330:		fe.add("interface_name", "can't be blank")
internal/engine/triggers/match.go:333:		fe.add("match_path", "can't be blank")
internal/engine/triggers/match.go:338:		fe.add("value_match_operator", "can't be blank")
internal/engine/triggers/match.go:356:			fe.add("interface_major", "can't be blank")
internal/engine/triggers/match.go:368:		fe.add("known_value", "can't be blank")
internal/engine/triggers/match.go:408:		fe.add("on", "can't be blank")
internal/engine/triggers/match.go:430:				fe.add("interface_major", "can't be blank")
internal/engine/triggers/match.go:441:				fe.add("interface_major", "can't be blank")
internal/appengine/http.go:407:			FieldErrors{"group_name": {"can't be blank"}})
internal/appengine/http.go:411:			FieldErrors{"devices": {"can't be blank"}})
test/conformance/upstream/recordchannels/main.go:376:	// `{"errors":{"simple_trigger":{"interface_major":["can't be blank"]}}}`.
test/conformance/upstream/README.md:247:`{"errors":{"simple_trigger":{"interface_major":["can't be blank"]}}}` — a *validation* error
internal/pairing/http.go:122:			map[string][]string{"hw_id": {"can't be blank"}})
internal/pairing/http.go:175:			map[string][]string{"csr": {"can't be blank"}})
internal/pairing/http.go:255:			map[string][]string{"client_crt": {"can't be blank"}})
internal/pairing/testdata/http/register_422_blank_hwid.json:1:{"errors":{"hw_id":["can't be blank"]}}
internal/realm/service.go:574:			return fmt.Errorf("%w: policy %q is still used by trigger %q", ErrValidation, name, trs[i].Name)

$ cd /root/astrate-mule && rg -n "validationDetail" --glob '!*_test.go' | head -20
internal/appengine/http.go:659:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
internal/appengine/http.go:688:// validationDetail strips the ErrValidation prefix for the response detail.
internal/appengine/http.go:689:func validationDetail(err error) string {
internal/housekeeping/http.go:241:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
internal/housekeeping/http.go:251:// validationDetail strips the ErrValidation prefix for the response detail.
internal/housekeeping/http.go:252:func validationDetail(err error) string {
internal/flowapi/http.go:278:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
internal/flowapi/http.go:303:func validationDetail(err error) string {
internal/realm/http.go:393:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
internal/realm/http.go:411:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
internal/realm/http.go:560:// validationDetail strips the ErrValidation sentinel prefix so the response
internal/realm/http.go:562:func validationDetail(err error) string {

→ Read internal/realm/http.go [limit=60, offset=370]
→ Read internal/realm/http.go [limit=30, offset=552]

$ cd /root/astrate-mule && rg -n "func \(a \*API\) (putAuthConfig|createPolicy|deletePolicy)" internal/realm/http.go
71:func (a *API) createPolicy(w http.ResponseWriter, r *http.Request) {
94:func (a *API) deletePolicy(w http.ResponseWriter, r *http.Request) {

→ Read internal/realm/policy.go [limit=40]

$ cd /root/astrate-mule && rg -n "jwt_public_key_pem" internal/realm/service.go
619:// /config/auth → {"jwt_public_key_pem": "..."}). Astrate stores a list for
634:		return fmt.Errorf("%w: jwt_public_key_pem can't be blank", ErrValidation)

→ Read internal/realm/http.go [limit=70, offset=60]
→ Read internal/realm/service.go [limit=30, offset=615]
→ Read internal/realm/service.go [limit=25, offset=560]

$ cd /root/astrate-mule && rg -n "ErrValidation" internal/realm/*.go | rg -v _test | head -20 && echo "---" && rg -n "config/auth|setAuthConfig|putAuthConfig|SetAuthKey" internal/realm/*.go internal/api/*.go 2>/dev/null | rg -v _test | head -20
internal/realm/policy.go:10:// triggers.CompilePolicy. Violations wrap ErrValidation (mapped to 422).
internal/realm/policy.go:14:		return "", fmt.Errorf("%w: %v", ErrValidation, err)
internal/realm/service.go:26:// ErrValidation wraps a request that is well-formed JSON but violates an
internal/realm/service.go:28:// and fmt.Errorf("%w: ...", ErrValidation, ...) to attach the detail.
internal/realm/service.go:29:var ErrValidation = errors.New("realm: validation failed")
internal/realm/service.go:151:// store.ErrAlreadyExists; a schema violation yields ErrValidation — carrying
internal/realm/service.go:163:		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
internal/realm/service.go:198:		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
internal/realm/service.go:225:		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
internal/realm/service.go:426:		return nil, fmt.Errorf("%w: trigger does not parse: %v", ErrValidation, err)
internal/realm/service.go:429:		return nil, fmt.Errorf("%w: trigger requires a name", ErrValidation)
internal/realm/service.go:436:		// generic ErrValidation.
internal/realm/service.go:441:		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
internal/realm/service.go:446:				return nil, fmt.Errorf("%w: policy %q does not exist in realm %q", ErrValidation, ct.PolicyName, realm)
internal/realm/service.go:574:			return fmt.Errorf("%w: policy %q is still used by trigger %q", ErrValidation, name, trs[i].Name)
internal/realm/service.go:634:		return fmt.Errorf("%w: jwt_public_key_pem can't be blank", ErrValidation)
internal/realm/http.go:392:	case errors.Is(err, ErrValidation):
internal/realm/http.go:410:	case errors.Is(err, ErrValidation):
internal/realm/http.go:560:// validationDetail strips the ErrValidation sentinel prefix so the response
---
internal/realm/http.go:48:	mux.Handle("GET /realmmanagement/v1/{realm}/config/auth", h(a.getAuth))
internal/realm/http.go:49:	mux.Handle("PUT /realmmanagement/v1/{realm}/config/auth", h(a.putAuth))
internal/realm/http.go:304:// --- config/auth ------------------------------------------------------------
internal/realm/http.go:306:// authConfig is the GET/PUT /config/auth body shape (upstream
internal/realm/http.go:327:	if err := a.svc.SetAuthKey(r.Context(), r.PathValue("realm"), cfg.JWTPublicKeyPEM); err != nil {
internal/realm/service.go:616:// --- config/auth ------------------------------------------------------------
internal/realm/service.go:619:// /config/auth → {"jwt_public_key_pem": "..."}). Astrate stores a list for
internal/realm/service.go:630:// SetAuthKey rotates the realm's JWT public key (upstream PUT /config/auth).
internal/realm/service.go:632:func (s *Service) SetAuthKey(ctx context.Context, realm, keyPEM string) error {

→ Read internal/realm/http.go [limit=40, offset=304]
→ Read internal/engine/triggers/policy.go [limit=30, offset=100]

$ cd /root/astrate-mule && rg -n "func \(s \*Service\) CreatePolicy" -A 30 internal/realm/service.go
515:func (s *Service) CreatePolicy(ctx context.Context, realm string, def []byte) (*store.TriggerPolicy, error) {
516-	rid, err := s.realmID(ctx, realm)
517-	if err != nil {
518-		return nil, err
519-	}
520-	name, err := validatePolicy(def)
521-	if err != nil {
522-		return nil, err
523-	}
524-	return s.st.CreateTriggerPolicy(ctx, rid, name, def)
525-}
526-
527-// GetPolicy returns one policy's definition.
528-func (s *Service) GetPolicy(ctx context.Context, realm, name string) (json.RawMessage, error) {
529-	rid, err := s.realmID(ctx, realm)
530-	if err != nil {
531-		return nil, err
532-	}
533-	p, err := s.st.GetTriggerPolicy(ctx, rid, name)
534-	if err != nil {
535-		return nil, err
536-	}
537-	return p.Definition, nil
538-}
539-
540-// ListPolicies returns the realm's policy names.
541-func (s *Service) ListPolicies(ctx context.Context, realm string) ([]string, error) {
542-	rid, err := s.realmID(ctx, realm)
543-	if err != nil {
544-		return nil, err
545-	}

$ cd /root/astrate-mule && rg -n 'fmt.Errorf' internal/engine/triggers/policy.go | head -30
110:		return nil, fmt.Errorf("triggers: policy does not parse as JSON: %w", err)
113:		return nil, fmt.Errorf("triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]")
116:		return nil, fmt.Errorf("triggers: policy must have at least one error handler")
130:			return nil, fmt.Errorf("triggers: error handler strategy must be discard or retry")
134:			return nil, fmt.Errorf("triggers: error handler %d: %w", i+1, err)
144:		return nil, fmt.Errorf("triggers: %w", err)
149:			return nil, fmt.Errorf("triggers: retry_times must be 1-100 when any handler retries")
152:		return nil, fmt.Errorf("triggers: retry_times requires a retry handler")
155:		return nil, fmt.Errorf("triggers: maximum_capacity must be a positive integer")
158:		return nil, fmt.Errorf("triggers: prefetch_count must be between 1 and 300")
161:		return nil, fmt.Errorf("triggers: event_ttl must be non-negative")
197:			return compiledHandler{}, fmt.Errorf("on-keyword must be any_error, client_error, or server_error")
207:			return compiledHandler{}, fmt.Errorf("on must be a keyword or a status-code array")
212:		return compiledHandler{}, fmt.Errorf("status-code list must be non-empty")
216:			return compiledHandler{}, fmt.Errorf("status codes must be in 400..599")
247:				return fmt.Errorf("error handlers must all handle distinct errors")

$ cd /root/astrate-mule && ls internal/realm/ && rg -n "func newTestAPI|func testAPI|httptest.NewServer" internal/realm/*_test.go | head -20
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

$ cd /root/astrate-mule && rg -n "func " internal/realm/helpers_test.go | head -20
12:func TestNormaliseIfaceName(t *testing.T) {
34:func TestJoinPEM(t *testing.T) {
53:func TestTriggerErrorBody(t *testing.T) {
107:func TestMajorParam(t *testing.T) {

$ cd /root/astrate-mule && rg -ln "realm.NewAPI|realm\.New|newRig" internal/realm/*_test.go test/ 2>/dev/null | head -10
internal/realm/async_operation_test.go
internal/realm/http_test.go
internal/realm/detailed_test.go
internal/realm/dashboard_compat_test.go
test/conformance/cpc/cpc_test.go
test/security/security_test.go

$ cd /root/astrate-mule && rg -n "realm.NewAPI|func newRig|func newAPI" internal/realm/http_test.go | head -10
56:func newRig(t *testing.T) *rig {

→ Read internal/realm/http_test.go [limit=110]

$ cd /root/astrate-mule && rg -n "func WriteError" -A 25 internal/astarteapi/*.go | head -45
rg: internal/astarteapi/*.go: IO error for operation on internal/astarteapi/*.go: No such file or directory (os error 2)

$ cd /root/astrate-mule && rg -rn --files-with-matches "func WriteError" . | head
./pkg/astarteapi/envelope.go

$ cd /root/astrate-mule && rg -n "func WriteError" -A 22 pkg/astarteapi/envelope.go
163:func WriteError(w http.ResponseWriter, status int, detail string) error {
164-	return write(w, status, detailEnvelope{Errors: detailBody{Detail: detail}})
165-}
166-
167-// WriteFieldErrors writes the Phoenix-changeset-shaped error envelope
168-// {"errors": {"<field>": ["<message>", ...]}} used by upstream 422 validation
169-// failures (for example {"errors": {"hw_id": ["is invalid"]}}). Keys are
170-// emitted in sorted order (Go map marshalling), which matches the
171-// deterministic bodies upstream produces for single-field failures.
172-func WriteFieldErrors(w http.ResponseWriter, status int, fields map[string][]string) error {
173-	return write(w, status, fieldsEnvelope{Errors: fields})
174-}
175-
176-// rawEnvelope is the caller-shaped error wrapper: {"errors": <anything>}.
177-type rawEnvelope struct {
178-	Errors any `json:"errors"`
179-}
180-
181-// WriteRawErrors writes {"errors": <anything>} — caller-shaped error bodies
182-// for the nested changeset envelopes upstream uses inside sub-objects (for
183-// example a trigger's {"action": {"http_url": [...]}}). Map keys are emitted
184-// in sorted order, without HTML escaping, and without a trailing newline,
185-// like every envelope here.

← Write /tmp/opencode/probe422.go
Wrote file successfully.

$ cd /root/astrate-mule && go run /tmp/opencode/probe422.go
package command-line-arguments
	/tmp/opencode/probe422.go:9:2: use of internal package github.com/astrate-platform/astrate/internal/engine/triggers not allowed

$ cd /root/astrate-mule && rg -n "policy must have at least one error handler|policy name must be 1-128|is still used by trigger|jwt_public_key_pem can't be blank" --glob '*_test.go' | head -20
(no output)

$ cd /root/astrate-mule && rg -rn "at least one error handler|still used by trigger" --glob '*_test.go' -l | head; echo "--- json fixtures:"; ls internal/realm/testdata 2>/dev/null; rg -ln "at least one error handler" . | head
internal/engine/triggers/policy_test.go
--- json fixtures:
./internal/engine/triggers/policy.go
./internal/engine/triggers/policy_test.go

$ cd /root/astrate-mule && rg -n "at least one error handler|still used by trigger" internal/engine/triggers/policy_test.go | head; echo ---; rg -n "CompilePolicy" -B3 -A12 internal/engine/triggers/policy_test.go | sed -n '1,60p'
92:			wantErr: "at least one error handler",
---
4-	"testing"
5-)
6-
7:func TestCompilePolicy(t *testing.T) {
8-	tests := []struct {
9-		name              string
10-		def               string
11-		wantErr           string
12-		wantPrefetchCount int // 0 means don't check
13-	}{
14-		{
15-			name:    "valid discard policy",
16-			def:     `{"name":"test","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":10}`,
17-			wantErr: "",
18-		},
19-		{
--
204-	}
205-	for _, tt := range tests {
206-		t.Run(tt.name, func(t *testing.T) {
207:			p, err := CompilePolicy([]byte(tt.def))
208-			if tt.wantErr == "" {
209-				if err != nil {
210-					t.Fatalf("unexpected error: %v", err)
211-				}
212-				if p == nil {
213-					t.Fatal("expected non-nil policy")
214-				}
215-				if tt.wantPrefetchCount != 0 && p.PrefetchCount != tt.wantPrefetchCount {
216-					t.Errorf("PrefetchCount = %d, want %d", p.PrefetchCount, tt.wantPrefetchCount)
217-				}
218-				return
219-			}
--
318-	}
319-	for _, tt := range tests {
320-		t.Run(tt.name, func(t *testing.T) {
321:			p, err := CompilePolicy([]byte(tt.def))
322-			if err != nil {
323-				t.Fatalf("compile: %v", err)
324-			}
325-			d := p.Decide(tt.status)
326-			if d.Strategy != tt.wantStrat {
327-				t.Errorf("strategy: got %v, want %v", d.Strategy, tt.wantStrat)
328-			}
329-			if d.Reason == "" {
330-				t.Error("reason must be non-empty")
331-			}
332-			if tt.wantReason != "" && !contains(d.Reason, tt.wantReason) {
333-				t.Errorf("reason %q does not contain %q", d.Reason, tt.wantReason)
--
352-
353-func TestDecideTransportReasonDistinct(t *testing.T) {
354-	def := `{"name":"k","error_handlers":[{"on":"server_error","strategy":"retry"}],"maximum_capacity":1,"retry_times":3}`
355:	p, err := CompilePolicy([]byte(def))
356-	if err != nil {
357-		t.Fatalf("compile: %v", err)
358-	}
359-	dTransport := p.Decide(StatusTransport)
360-	dCode := p.Decide(500)

$ cd /root/astrate-mule && go test ./internal/engine/triggers/ -run TestCompilePolicy -v 2>&1 | tail -30
--- PASS: TestCompilePolicy/retry_times_forbidden_when_no_retry_handler (0.00s)
    --- PASS: TestCompilePolicy/retry_times_out_of_range_low (0.00s)
    --- PASS: TestCompilePolicy/retry_times_out_of_range_high (0.00s)
    --- PASS: TestCompilePolicy/maximum_capacity_too_low (0.00s)
    --- PASS: TestCompilePolicy/negative_event_ttl (0.00s)
    --- PASS: TestCompilePolicy/name_starts_with_@ (0.00s)
    --- PASS: TestCompilePolicy/no_handlers (0.00s)
    --- PASS: TestCompilePolicy/invalid_JSON (0.00s)
    --- PASS: TestCompilePolicy/empty_on_keyword (0.00s)
    --- PASS: TestCompilePolicy/on_empty_array (0.00s)
    --- PASS: TestCompilePolicy/on_code_out_of_range (0.00s)
    --- PASS: TestCompilePolicy/overlapping_explicit_lists_rejected (0.00s)
    --- PASS: TestCompilePolicy/any_error_twice_rejected (0.00s)
    --- PASS: TestCompilePolicy/any_error_plus_explicit_code_rejected (0.00s)
    --- PASS: TestCompilePolicy/server_error_plus_any_error_rejected (0.00s)
    --- PASS: TestCompilePolicy/client_error_plus_explicit_404_rejected_cleanly (0.00s)
    --- PASS: TestCompilePolicy/same_code_discarded_twice_rejected_regardless_of_strategy (0.00s)
    --- PASS: TestCompilePolicy/chained_three-way_overlap_caught_by_pairwise_checks (0.00s)
    --- PASS: TestCompilePolicy/disjoint_explicit_codes_accepted (0.00s)
    --- PASS: TestCompilePolicy/server_error_plus_client-range_code_accepted (0.00s)
    --- PASS: TestCompilePolicy/prefetch_count_round-trips (0.00s)
    --- PASS: TestCompilePolicy/prefetch_count_defaults_to_1_when_omitted (0.00s)
    --- PASS: TestCompilePolicy/prefetch_count_1_accepted_(lower_bound) (0.00s)
    --- PASS: TestCompilePolicy/prefetch_count_300_accepted_(upper_bound) (0.00s)
    --- PASS: TestCompilePolicy/prefetch_count_0_rejected (0.00s)
    --- PASS: TestCompilePolicy/prefetch_count_301_rejected (0.00s)
    --- PASS: TestCompilePolicy/retry_times_1_accepted_(lower_bound) (0.00s)
    --- PASS: TestCompilePolicy/retry_times_100_accepted_(upper_bound) (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	0.171s

→ Read docs/Makefile

$ cd /root/astrate-mule && rg -n "ErrorDetail:" -A 14 docs/api/astarte_realm_management_api.yaml
1185:    ErrorDetail:
1186-      type: object
1187-      required: [errors]
1188-      properties:
1189-        errors:
1190-          type: object
1191-          required: [detail]
1192-          properties:
1193-            detail:
1194-              type: string
1195-
1196-    TriggerErrorsDetail:
1197-      type: object
1198-      required: [errors]
1199-      properties:

$ cd /root/astrate-mule && rg -rln "astarte_realm_management_api" --glob '*.go' --glob '*.html' --glob '*.yml' --glob '*.yaml' --glob '*.mk' --glob 'Makefile' . | head -20
./docs/swagger-ui/index.html:      <option value="../api/ln.yaml">Realm Management API — Interfaces, triggers, policies, a_rma auth</option>
./internal/swagger/swagger_test.go:	specDocuments403(t, "ln.yaml")
./internal/swagger/swagger_test.go:	b, err := docs.APIYAML.ReadFile("api/ln.yaml")
./internal/swagger/swagger_test.go:		t.Fatalf("reading ln.yaml: %v", err)
./internal/swagger/swagger_test.go:	b, err := docs.APIYAML.ReadFile("api/ln.yaml")
./internal/swagger/swagger_test.go:		t.Fatalf("reading ln.yaml: %v", err)

$ cd /root/astrate-mule && cat .gitignore 2>/dev/null; echo "=== swagger-ui:"; ls docs/swagger-ui; echo "=== which mkdocs:"; which mkdocs || echo "no mkdocs"; python3 -c "import yaml, sys; print('pyyaml ok')" 2>&1 | tail -1
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
=== swagger-ui:
index.html
=== which mkdocs:
no mkdocs
ModuleNotFoundError: No module named 'yaml'

$ cd /root/astrate-mule && ls internal/swagger/ && rg -rn "ln.yaml" . --glob '!docs/site/**' | head -20
swagger.go
swagger_test.go

$ cd /root/astrate-mule && rg -n "ln\.yaml" -g '!docs/site/**' . | head -20; echo "exit=$?"
exit=0

$ cd /root/astrate-mule && rg -n "yaml" internal/swagger/swagger.go | head -20; echo "=== embed.go:"; cat docs/embed.go
35:	// Serve YAML files at /api/ so the relative ../api/*.yaml in index.html
63:		if strings.HasSuffix(path, ".yaml") {
=== embed.go:
// Package docs embeds the Swagger UI static files and OpenAPI YAML specs so
// they are served directly from the astrate binary without external files.
package docs

import "embed"

// SwaggerUI holds the static Swagger UI files (index.html, CSS, JS refs).
//
//go:embed swagger-ui/*
var SwaggerUI embed.FS

// APIYAML holds the OpenAPI 3.0 YAML specifications for all five API surfaces.
//
//go:embed api/*.yaml
var APIYAML embed.FS

$ cd /root/astrate-mule && grep -rn "ln\.yaml" --include="*.go" --include="*.html" . | head -20
(no output)

→ Read internal/swagger/swagger.go
→ Read internal/swagger/swagger_test.go

$ cd /root/astrate-mule && grep -rn "yaml" go.mod; ls ../examples/pokemon-agent/docs/DESIGN.md 2>&1; ls ../examples 2>&1 | head
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect
ls: cannot access '../examples/pokemon-agent/docs/DESIGN.md': No such file or directory
ls: cannot access '../examples': No such file or directory

$ cd /root/astrate-mule && grep -rn "responses/ValidationError\b\|ValidationError:" --include="*.go" --include="*.yaml" --include="*.md" --include="*.html" . | grep -v "docs/site/" | head -20; echo "=== responses/ValidationError refs:"; grep -rn "responses/ValidationError" --include="*.yaml" . | grep -v "docs/site/" | head
./.mule/reports/log/20261002T181712Z-recipe-docs-sync-proposed.md:454:1330:    ValidationError:
./.mule/reports/log/20261002T181712Z-recipe-docs-sync-proposed.md:719:$ rg -n 'responses/ValidationError' astarte_realm_management_api.yaml; echo "--- Conflict refs"; rg -n 'responses/Conflict' astarte_realm_management_api.yaml; echo "--- NotFound refs"; rg -n 'responses/NotFound' astarte_realm_management_api.yaml
./.mule/reports/log/20261002T181712Z-recipe-docs-sync-proposed.md:720:574:          $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20261002T181712Z-recipe-docs-sync-proposed.md:721:782:          $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20261002T181712Z-recipe-docs-sync-proposed.md:722:839:          $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20260928T180120Z-docs-sync-pairing-unregister-description-done.md:27:          $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20260924T202134Z-recipe-docs-sync-proposed.md:159:218-          $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20260905T172234Z-docs-sync-pairing-register-404-done.md:51:   $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20260917T202137Z-docs-sync-rm-triggers-422-nested-envelope-done.md:39:-        $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20260911T170716Z-docs-sync-rm-put-auth-422-done.md:35:+        $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20260911T170716Z-docs-sync-rm-put-auth-422-done.md:152:PUT /config/auth response 422 -> map[$ref:#/components/responses/ValidationError]
./.mule/reports/log/20260916T165431Z-docs-sync-pairing-deviceendpoints-dead-404-403-done.md:38:   $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20260916T165431Z-docs-sync-pairing-deviceendpoints-dead-404-403-done.md:49:   $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20260916T165431Z-docs-sync-pairing-deviceendpoints-dead-404-403-done.md:60:   $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20260916T165431Z-docs-sync-pairing-deviceendpoints-dead-404-403-done.md:105:   $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20260928T172627Z-docs-sync-pairing-initial-payload-format-enum-done.md:99:    ValidationError:
./.mule/reports/log/20260923T201746Z-docs-sync-rm-put-interface-409-done.md:67:   $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20260923T201746Z-docs-sync-rm-put-interface-409-done.md:108: ValidationError:
./.mule/reports/log/20260923T201746Z-docs-sync-rm-put-interface-409-done.md:293:           $ref: "#/components/responses/ValidationError"
./.mule/reports/log/20260923T201746Z-docs-sync-rm-put-interface-409-done.md:325:     ValidationError:
=== responses/ValidationError refs:
./docs/api/astarte_realm_management_api.yaml:574:          $ref: "#/components/responses/ValidationError"
./docs/api/astarte_realm_management_api.yaml:782:          $ref: "#/components/responses/ValidationError"
./docs/api/astarte_realm_management_api.yaml:839:          $ref: "#/components/responses/ValidationError"
./docs/api/astarte_appengine_api.yaml:383:          $ref: "#/components/responses/ValidationErrors"
./docs/api/astarte_appengine_api.yaml:435:          $ref: "#/components/responses/ValidationErrors"
./docs/api/astarte_appengine_api.yaml:652:          $ref: "#/components/responses/ValidationErrors"
./docs/api/astarte_appengine_api.yaml:697:          $ref: "#/components/responses/ValidationErrors"
./docs/api/astarte_appengine_api.yaml:907:          $ref: "#/components/responses/ValidationErrors"
./docs/api/astarte_appengine_api.yaml:987:          $ref: "#/components/responses/ValidationErrors"
./docs/api/astarte_appengine_api.yaml:1033:          $ref: "#/components/responses/ValidationErrors"

$ cd /root/astrate-mule && ls .mule/recipes/ && rg -n "^#|^##" .mule/recipes/docs-sync.md | head -40
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md
1:# Recipe — docs sync
11:## Find drift, cheaply
24:## What to propose
36:## Verify the site still builds

→ Read .mule/recipes/docs-sync.md
Now the three inline 422 blocks. First `putAuthConfig`:

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -570,9 +570,20 @@
         $ref: "#/components/responses/Unauthorized"
       "403":
         $ref: "#/components/responses/Forbidden"
       "422":
-        $ref: "#/components/responses/ValidationError"
+        description: |
+          Unprocessable Entity — the auth configuration failed validation.
+          A body with an empty `jwt_public_key_pem` is answered 422, not 400:
+          it parses, and SetAuthKey rejects a blank key (Service.SetAuthKey,
+          internal/realm/service.go) as an ErrValidation detail.
+        content:
+          application/json; charset=utf-8:
+            schema:
+              $ref: "#/components/schemas/ErrorDetail"
+            example:
+              errors:
+                detail: "jwt_public_key_pem can't be blank"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
 /realmmanagement/v1/{realm}/config/device_registration_limit:

→ Read docs/api/astarte_realm_management_api.yaml [limit=30, offset=564]

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -789,9 +789,29 @@
         $ref: "#/components/responses/Forbidden"
       "409":
         $ref: "#/components/responses/Conflict"
       "422":
-        $ref: "#/components/responses/ValidationError"
+        description: |
+          Unprocessable Entity — the policy definition failed validation.
+          The detail is the trigger compiler's own message with the
+          `realm: validation failed: ` prefix stripped, so it names the
+          policy constraint that was broken (triggers.CompilePolicy via
+          realm.validatePolicy).
+        content:
+          application/json; charset=utf-8:
+            schema:
+              $ref: "#/components/schemas/ErrorDetail"
+            examples:
+              no_error_handlers:
+                summary: The policy declares no error handler.
+                value:
+                  errors:
+                    detail: "triggers: policy must have at least one error handler"
+              bad_name:
+                summary: The policy name is empty or holds characters outside [a-zA-Z0-9_.~-].
+                value:
+                  errors:
+                    detail: "triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
 /realmmanagement/v1/{realm}/policies/{name}:

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -866,9 +866,20 @@
           $ref: "#/components/responses/Forbidden"
         "404":
           $ref: "#/components/responses/NotFound"
         "422":
-          $ref: "#/components/responses/ValidationError"
+          description: |
+            Unprocessable Entity — the policy is still referenced by a trigger.
+            Deletion is refused while any trigger in the realm names this
+            policy; the detail names both (Service.DeletePolicy,
+            internal/realm/service.go).
+          content:
+            application/json; charset=utf-8:
+              schema:
+                $ref: "#/components/schemas/ErrorDetail"
+              example:
+                errors:
+                  detail: 'policy "audit" is still used by trigger "on_audit_failure"'
         "500":
           $ref: "#/components/responses/InternalServerError"
 
 components:

Now remove the now-unreferenced shared component:
← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -1368,18 +1368,8 @@
           value:
             errors:
               detail: Interface major version doesn't match the one in the interface json
 
-ValidationError:
-  description: Unprocessable Entity — request body failed validation.
-  content:
-    application/json; charset=utf-8:
-      schema:
-        $ref: "#/components/schemas/ErrorDetail"
-      example:
-        errors:
-          detail: "invalid interface: missing \"interface_name\""
-
 InternalServerError:
   description: Internal Server Error — unexpected failure.
   content:
     application/json; charset=utf-8:

$ cd /root/astrate-mule && grep -n "ValidationError" docs/api/astarte_realm_management_api.yaml; echo "exit=$?"; sed -n '790,830p;860,880p' docs/api/astarte_realm_management_api.yaml
exit=1
        "409":
          $ref: "#/components/responses/Conflict"
        "422":
          description: |
            Unprocessable Entity — the policy definition failed validation.
            The detail is the trigger compiler's own message with the
            `realm: validation failed: ` prefix stripped, so it names the
            policy constraint that was broken (triggers.CompilePolicy via
            realm.validatePolicy).
          content:
            application/json; charset=utf-8:
              schema:
                $ref: "#/components/schemas/ErrorDetail"
              examples:
                no_error_handlers:
                  summary: The policy declares no error handler.
                  value:
                    errors:
                      detail: "triggers: policy must have at least one error handler"
                bad_name:
                  summary: The policy name is empty or holds characters outside [a-zA-Z0-9_.~-].
                  value:
                    errors:
                      detail: "triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]"
        "500":
          $ref: "#/components/responses/InternalServerError"

  /realmmanagement/v1/{realm}/policies/{name}:
    get:
      operationId: getPolicy
      summary: Get a trigger delivery policy
      description: Returns the full policy definition.
      tags: [Realm Management - Policies]
      security:
        - a_rma: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - $ref: "#/components/parameters/PolicyName"
      responses:
        "200":
          description: Policy definition.
      responses:
        "204":
          description: Policy deleted. No body returned.
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "404":
          $ref: "#/components/responses/NotFound"
        "422":
          description: |
            Unprocessable Entity — the policy is still referenced by a trigger.
            Deletion is refused while any trigger in the realm names this
            policy; the detail names both (Service.DeletePolicy,
            internal/realm/service.go).
          content:
            application/json; charset=utf-8:
              schema:
                $ref: "#/components/schemas/ErrorDetail"
              example:
                errors:

$ cd /root/astrate-mule && tail -20 docs/api/astarte_realm_management_api.yaml
nameMismatch:
              summary: The interface name in the body does not match the one in the URL.
              value:
                errors:
                  detail: Interface name doesn't match the one in the interface json
            majorMismatch:
              summary: The interface major version in the body does not match the one in the URL.
              value:
                errors:
                  detail: Interface major version doesn't match the one in the interface json

    InternalServerError:
      description: Internal Server Error — unexpected failure.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Internal server error

✗ Edit internal/swagger/swagger_test.go failed
Error: Could not find oldString in the file. It must match exactly, including whitespace, indentation, and line endings.
✗ Edit internal/swagger/swagger_test.go failed
Error: Could not find oldString in the file. It must match exactly, including whitespace, indentation, and line endings.
→ Read internal/swagger/swagger_test.go [offset=508, limit=14]

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -512,10 +512,12 @@
 // examples carried the capital-R/capital-S spellings, which upstream Phoenix
 // never renders and no handler emits: every 400 on this surface goes through
 // astarteapi.WriteBadRequest and every 500 through WriteInternalServerError
 // (internal/realm/http.go), and the SDK/astartectl error paths match on the
-// canonical strings. The Conflict and ValidationError examples are
-// endpoint-specific details with no constant behind them, so they stay free-form.
+// canonical strings. The Conflict example is an endpoint-specific detail with no
+// constant behind it, so it stays free-form; the three 422s that used to share
+// one ValidationError component are pinned by
+// TestRealmManagement422ValidationDetails.
 func TestRealmManagementErrorDetailExamples(t *testing.T) {
 	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
 	if err != nil {
 		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)

→ Read internal/swagger/swagger_test.go [limit=14, offset=546]
Now the new test, placed after `TestRealmManagementErrorDetailExamples`:
✗ Edit internal/swagger/swagger_test.go failed
Error: Found multiple matches for oldString. Provide more surrounding context to make the match unique.

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -556,8 +556,92 @@
 		}
 	}
 }
 
+// TestRealmManagement422ValidationDetails pins the 422 examples the three
+// non-interface realm-management operations carry to the details the wire
+// actually emits, so a generated client that copies one does not learn a
+// message no handler can produce.
+//
+// The three operations used to share components.responses.ValidationError, whose
+// single example read `invalid interface: missing "interface_name"`. That string
+// is wrong for all three: it names an interface, and none of them installs or
+// updates one. Each now carries its own inline 422, the pattern
+// installInterface and updateInterface already use for their own details.
+//
+// The createPolicy examples are derived from the compiler here instead of being
+// written out, so a change to the message fails here instead of drifting the
+// spec — the coupling TestPairingDeviceIDEncodingDocumented gets from
+// deviceid. The putAuthConfig and deletePolicy details have no exported
+// constant behind them (Service.SetAuthKey and Service.DeletePolicy,
+// internal/realm/service.go), so they are spelled out and cited here.
+func TestRealmManagement422ValidationDetails(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	if strings.Contains(string(b), "#/components/responses/ValidationError") {
+		t.Error("spec still $refs a shared ValidationError response; one example cannot be right for every surface that refs it")
+	}
+	if containsLine(lines, "    ValidationError:") {
+		t.Error("spec still defines components.responses.ValidationError")
+	}
+
+	_, noHandlers := triggers.CompilePolicy([]byte(`{"name":"p","error_handlers":[]}`))
+	if noHandlers == nil {
+		t.Fatal("triggers.CompilePolicy accepted a policy with no error handler; this test's premise no longer holds")
+	}
+	_, badName := triggers.CompilePolicy(
+		[]byte(`{"name":"","error_handlers":[{"on":"any_error","strategy":"discard"}]}`))
+	if badName == nil {
+		t.Fatal("triggers.CompilePolicy accepted an empty policy name; this test's premise no longer holds")
+	}
+
+	for _, tc := range []struct {
+		op   string
+		want []string
+	}{
+		{"putAuthConfig", []string{`jwt_public_key_pem can't be blank`}},
+		{"createPolicy", []string{noHandlers.Error(), badName.Error()}},
+		{"deletePolicy", []string{`policy "audit" is still used by trigger "on_audit_failure"`}},
+	} {
+		block := strings.Join(responseBlock(t, operationBlock(t, lines, tc.op), "422"), "\n")
+		for _, want := range tc.want {
+			if !strings.Contains(block, want) {
+				t.Errorf("%s 422 carries no example with detail %q", tc.op, want)
+			}
+		}
+		if strings.Contains(block, "$ref") {
+			t.Errorf("%s 422 still $refs a component instead of carrying its own example", tc.op)
+		}
+	}
+}
+
+// responseBlock returns the lines of the operation's response with the given
+// status code, up to the next status code or the end of the operation block.
+func responseBlock(t *testing.T, block []string, status string) []string {
+	t.Helper()
+	key := `        "` + status + `":`
+	start := -1
+	for i, l := range block {
+		if l == key {
+			start = i
+			break
+		}
+	}
+	if start < 0 {
+		t.Fatalf("operation block declares no %s response", status)
+	}
+	for i := start + 1; i < len(block); i++ {
+		if strings.HasPrefix(block[i], `        "`) {
+			return block[start+1 : i]
+		}
+	}
+	return block[start+1:]
+}
+
 // TestPairingDeviceIDEncodingDocumented guards that the pairing spec states the
 // exact wire form of a device identifier — the 22-character unpadded base64url
 // string deviceid.Parse accepts — instead of the "base64-encoded 128-bit"
 // wording, which read as standard base64 and let a generated client send a

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -13,8 +13,9 @@
 	"strconv"
 	"strings"
 	"testing"
 
+	"github.com/astrate-platform/astrate/internal/engine/triggers"
 	docs "github.com/astrate-platform/astrate/docs"
 	"github.com/astrate-platform/astrate/pkg/astarteapi"
 	"github.com/astrate-platform/astrate/pkg/deviceid"
 )

$ cd /root/astrate-mule && gofmt -l . && go vet ./internal/swagger/ && go test ./internal/swagger/ -run 'TestRealmManagement' -v 2>&1 | tail -25
internal/swagger/swagger_test.go
=== RUN   TestRealmManagement403
--- PASS: TestRealmManagement403 (0.00s)
=== RUN   TestRealmManagementAsyncOperationParamDocumented
--- PASS: TestRealmManagementAsyncOperationParamDocumented (0.01s)
=== RUN   TestRealmManagementErrorDetailExamples
--- PASS: TestRealmManagementErrorDetailExamples (0.00s)
=== RUN   TestRealmManagement422ValidationDetails
    swagger_test.go:617: putAuthConfig 422 still $refs a component instead of carrying its own example
    swagger_test.go:617: createPolicy 422 still $refs a component instead of carrying its own example
    swagger_test.go:617: deletePolicy 422 still $refs a component instead of carrying its own example
--- FAIL: TestRealmManagement422ValidationDetails (0.01s)
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.046s
FAIL

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -612,10 +612,10 @@
 			if !strings.Contains(block, want) {
 				t.Errorf("%s 422 carries no example with detail %q", tc.op, want)
 			}
 		}
-		if strings.Contains(block, "$ref") {
-			t.Errorf("%s 422 still $refs a component instead of carrying its own example", tc.op)
+		if strings.Contains(block, "#/components/responses/") {
+			t.Errorf("%s 422 still $refs a response component instead of carrying its own example", tc.op)
 		}
 	}
 }
 

$ cd /root/astrate-mule && gofmt -w internal/swagger/swagger_test.go && gofmt -l . ; sed -n '1,25p' internal/swagger/swagger_test.go
package swagger

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	docs "github.com/astrate-platform/astrate/docs"
	"github.com/astrate-platform/astrate/internal/engine/triggers"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

func TestMount(t *testing.T) {
	mux := http.NewServeMux()
	Mount(mux)

$ cd /root/astrate-mule && gofmt -l . ; go vet ./... && go test ./internal/swagger/ -v -run 'TestRealmManagement|TestMount|TestSpecs' 2>&1 | tail -20
=== RUN   TestMount/GET_/swagger/_serves_the_embedded_UI
=== RUN   TestMount/GET_/api/_serves_every_OpenAPI_YAML_spec
--- PASS: TestMount (0.04s)
    --- PASS: TestMount/GET_/swagger_redirects_to_/swagger/index.html (0.00s)
    --- PASS: TestMount/GET_/swagger/_serves_the_embedded_UI (0.03s)
    --- PASS: TestMount/GET_/api/_serves_every_OpenAPI_YAML_spec (0.01s)
=== RUN   TestMountSubPanicsOnBrokenFS
--- PASS: TestMountSubPanicsOnBrokenFS (0.00s)
=== RUN   TestSpecs
--- PASS: TestSpecs (0.00s)
=== RUN   TestRealmManagement403
--- PASS: TestRealmManagement403 (0.00s)
=== RUN   TestRealmManagementAsyncOperationParamDocumented
--- PASS: TestRealmManagementAsyncOperationParamDocumented (0.01s)
=== RUN   TestRealmManagementErrorDetailExamples
--- PASS: TestRealmManagementErrorDetailExamples (0.00s)
=== RUN   TestRealmManagement422ValidationDetails
--- PASS: TestRealmManagement422ValidationDetails (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.086s

Now verify the test genuinely fails without the YAML change (reverse-patch, then restore from a saved copy):
$ cd /root/astrate-mule && cp docs/api/astarte_realm_management_api.yaml /tmp/opencode/new.yaml && sha256sum /tmp/opencode/new.yaml
a4cd9c3b58927187c4224c5c72f80a10dc0556e55e373528848b703ef2041942  /tmp/opencode/new.yaml

← Write /tmp/opencode/reverse.py
Wrote file successfully.

$ cd /root/astrate-mule && python3 /tmp/opencode/reverse.py && go test ./internal/swagger/ -run 'TestRealmManagement422ValidationDetails' 2>&1 | tail -12
reverse-patched: 44683 bytes, 3 ValidationError refs
    swagger_test.go:586: spec still $refs a shared ValidationError response; one example cannot be right for every surface that refs it
    swagger_test.go:589: spec still defines components.responses.ValidationError
    swagger_test.go:613: putAuthConfig 422 carries no example with detail "jwt_public_key_pem can't be blank"
    swagger_test.go:617: putAuthConfig 422 still $refs a response component instead of carrying its own example
    swagger_test.go:613: createPolicy 422 carries no example with detail "triggers: policy must have at least one error handler"
    swagger_test.go:613: createPolicy 422 carries no example with detail "triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]"
    swagger_test.go:617: createPolicy 422 still $refs a response component instead of carrying its own example
    swagger_test.go:613: deletePolicy 422 carries no example with detail "policy \"audit\" is still used by trigger \"on_audit_failure\""
    swagger_test.go:617: deletePolicy 422 still $refs a response component instead of carrying its own example
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.033s
FAIL

$ cd /root/astrate-mule && cp /tmp/opencode/new.yaml docs/api/astarte_realm_management_api.yaml && sha256sum docs/api/astarte_realm_management_api.yaml && go test ./internal/swagger/ -run 'TestRealmManagement422ValidationDetails' 2>&1 | tail -4
a4cd9c3b58927187c4224c5c72f80a10dc0556e55e373528848b703ef2041942  docs/api/astarte_realm_management_api.yaml
ok  	github.com/astrate-platform/astrate/internal/swagger	0.037s

$ cd /root/astrate-mule && (make -C docs build 2>&1 | tail -5); echo "--- mkdocs?"; python3 -m mkdocs --version 2>&1 | tail -2; ls /root/examples 2>&1 | head -2
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'
--- mkdocs?
/usr/bin/python3: No module named mkdocs
ls: cannot access '/root/examples': No such file or directory

$ cd /root/astrate-mule && make -C docs sync; echo "sync exit=$?"; echo "=== git status:"; git status --porcelain
make: Entering directory '/root/astrate-mule/docs'
cp DESIGN.md site/DESIGN.md
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
make: Leaving directory '/root/astrate-mule/docs'
sync exit=0
=== git status:
 M docs/api/astarte_realm_management_api.yaml
 M internal/swagger/swagger_test.go

! permission requested: external_directory (/root/examples/pokemon-agent/docs/*); auto-rejecting
✗ cd /root/astrate-mule/docs && cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md; echo "cp exit=$?"; ls -la site/pokemon-agent.md; ls /root/astrate-mule/../ | head failed
Error: The user rejected permission to use this specific tool call.
