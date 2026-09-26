slug: recipe-code-review
verdict: proposed
at:  6764463
ran: 2026-09-26T17:42:59Z on DietPi in 342s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md

$ ls internal/ pkg/ 2>/dev/null
internal/:
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

pkg/:
astarteapi
deviceid
interfaceschema
payload

$ ls .mule/ && echo "--- reviews ---" && ls .mule/reviews/ 2>/dev/null | tail -30 && echo "--- for-giulio ---" && ls .mule/for-giulio.md 2>/dev/null && echo "--- date ---" && date +%F
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
recipes
reports
reviews
tasks
todo.md
--- reviews ---
appengine-2026-09-07.md
astarteapi-2026-09-15.md
auth-2026-09-14.md
broker-2026-09-04.md
broker-2026-09-21.md
channels-2026-09-18.md
cmd-astrate-2026-09-25.md
config-2026-09-13.md
deviceid-2026-09-18.md
engine-2026-09-06.md
engine-triggers-2026-09-20.md
flow-2026-09-05.md
flow-2026-09-23.md
flowapi-2026-09-10.md
housekeeping-2026-09-15.md
httpx-2026-09-17.md
interfaceschema-2026-09-11.md
observability-2026-09-16.md
pairing-2026-09-09.md
payload-2026-09-12.md
realm-2026-09-09.md
store-2026-09-08.md
store-2026-09-21.md
swagger-2026-09-24.md
--- for-giulio ---
.mule/for-giulio.md
--- date ---
2026-09-26

$ tail -40 .mule/log.md
| 2026-09-16 | observability-dbpool-gauge-coverage | done | 62s | 19cd4f1 |
| 2026-09-16 | observability-health-content-type-test | done | 250s | f6c27f5 |
| 2026-09-17 | docs-sync-pairing-version-value | done | 1059s | c69a164 |
| 2026-09-17 | docs-sync-native-version-value | done | 503s | 08048b3 |
| 2026-09-17 | httpx-cors-vary-origin-passthrough | done | 143s | 986c100 |
| 2026-09-17 | docs-sync-rm-policies-delete-422 | done | 117s | c0521bc |
| 2026-09-17 | docs-sync-rm-triggers-422-nested-envelope | done | 161s | 8bfedd6 |
| 2026-09-18 | compat-note-v134 | done | 154s | bcd4383 |
| 2026-09-18 | channels-group-watch-membership-scope | done | 677s | aeaeb47 |
| 2026-09-18 | channels-rejoin-joinref-tagging-test | done | 208s | d48e17f |
| 2026-09-19 | channels-rejoin-authz-mismatch | blocked | 179s | wrote nothing |
| 2026-09-20 | errorname-missing-required-test | done | 135s | 739df67 |
| 2026-09-20 | triggers-custom-action-policy-nodecide | done | 459s | 3c8b89a |
| 2026-09-20 | compat-note-custom-action-policy-boundary | done | 97s | 8cd17a8 |
| 2026-09-20 | docs-sync-native-metrics-example-fake-series | done | 93s | 43c4a8b |
| 2026-09-20 | docs-sync-native-socket-missing-403-500 | done | 116s | f9d8b27 |
| 2026-09-21 | store-pipelines-empty-name-zero-blocks-test | done | 91s | 82fc0a4 |
| 2026-09-21 | broker-acl-coldstart-fallback-flood | done | 218s | 0e5078f |
| 2026-09-21 | broker-offlineacl-entry-eviction | done | 406s | 72b2686 |
| 2026-09-22 | docs-sync-ae-read-query-params | blocked | 159s | wrote nothing |
| 2026-09-23 | examples-echo-container-contract-test | done | 251s | f35af21 |
| 2026-09-23 | flow-mqtt-source-reconnect-recovery | blocked | 84s | wrote nothing |
| 2026-09-23 | flow-msg-json-integer-precision | done | 612s | 4df3f6b |
| 2026-09-23 | flow-sort-bounded-buffer | done | 634s | 9bf0ab7 |
| 2026-09-23 | flow-randomsource-span-overflow | blocked | 288s | lint failed: internal/flow/blocks/randomsource.go:113:40: G115: integer overflow conversion int64 -> uint64 (gosec) |
| 2026-09-23 | flow-filter-key-contains-test | done | 107s | 3875ba6 |
| 2026-09-23 | docs-sync-rm-put-interface-409 | done | 356s | 7256af9 |
| 2026-09-24 | docs-sync-rm-interface-422-shapes | done | 490s | 93e784e |
| 2026-09-24 | docs-sync-rm-validation-example-prefix | done | 238s | 02e269b |
| 2026-09-24 | docs-sync-rm-auth-403 | done | 668s | 9c05370 |
| 2026-09-24 | swagger-sub-failfast | done | 371s | 1d321ea |
| 2026-09-25 | drain-per-stage-budget | done | 265s | 84a050e |
| 2026-09-25 | cmd-devcert-fields-test | done | 357s | f1434eb |
| 2026-09-25 | cmd-healthcheck-contract-test | blocked | 199s | tests failed: --- FAIL: TestRunHealthcheckProbesReadiness (0.01s) |
| 2026-09-25 | cmd-loadsealer-masterkeyfile-test | done | 221s | a6711f9 |
| 2026-09-25 | docs-sync-hk-async-operation-param | done | 248s | ed1cd72 |
| 2026-09-25 | docs-sync-hk-retention-zero-is-unset | done | 325s | c2ad442 |
| 2026-09-25 | docs-sync-rm-async-operation-param | done | 425s | 3f0eca8 |
| 2026-09-25 | docs-sync-rm-delete-device-async-operation-param | done | 903s | f965f39 |
| 2026-09-25 | docs-sync-rm-delete-device-async-operation-param | blocked | 177s | wrote nothing |

$ find internal pkg -type d | sort
internal
internal/appengine
internal/appengine/channels
internal/appengine/stream
internal/auth
internal/auth/testdata
internal/broker
internal/config
internal/engine
internal/engine/forward
internal/engine/stream
internal/engine/testdata
internal/engine/triggers
internal/engine/triggers/testdata
internal/flow
internal/flow/blocks
internal/flow/blocks/astartesource
internal/flow/blocks/container
internal/flow/blocks/virtualdevicepool
internal/flowapi
internal/housekeeping
internal/httpx
internal/observability
internal/pairing
internal/pairing/ca
internal/pairing/testdata
internal/pairing/testdata/http
internal/realm
internal/store
internal/store/testdata
internal/swagger
internal/testutil
pkg
pkg/astarteapi
pkg/astarteapi/testdata
pkg/deviceid
pkg/interfaceschema
pkg/interfaceschema/testdata
pkg/interfaceschema/testdata/invalid
pkg/interfaceschema/testdata/valid
pkg/payload
pkg/payload/testdata
pkg/payload/testdata/bson

$ for d in internal/appengine/stream internal/engine/forward internal/engine/stream internal/flow/blocks/container internal/flow/blocks/virtualdevicepool internal/flow/blocks/astartesource internal/pairing/ca internal/testutil internal/housekeeping; do n=$(rg -c '' --glob "$d/*.go" 2>/dev/null | awk -F: '{s+=$2} END {print s+0}'); f=$(ls $d/*.go 2>/dev/null | wc -l); echo "$n lines  $f files  $d"; done | sort -rn
1372 lines  5 files  internal/housekeeping
1133 lines  4 files  internal/flow/blocks/container
713 lines  5 files  internal/engine/forward
707 lines  2 files  internal/flow/blocks/virtualdevicepool
697 lines  6 files  internal/testutil
596 lines  2 files  internal/pairing/ca
443 lines  2 files  internal/engine/stream
324 lines  2 files  internal/flow/blocks/astartesource
272 lines  2 files  internal/appengine/stream

$ rg -c '' internal/engine/forward/*.go internal/engine/stream/*.go internal/appengine/stream/*.go internal/pairing/ca/*.go 2>/dev/null | sort -t: -k2 -rn
internal/engine/forward/http_test.go:310
internal/pairing/ca/ca_test.go:309
internal/pairing/ca/ca.go:287
internal/engine/stream/bus_test.go:234
internal/engine/stream/bus.go:209
internal/engine/forward/nats_test.go:196
internal/appengine/stream/ws.go:162
internal/engine/forward/http.go:125
internal/appengine/stream/ws_test.go:110
internal/engine/forward/nats.go:53
internal/engine/forward/envelope.go:29

→ Read internal/engine/forward/nats.go
→ Read internal/engine/forward/http.go
→ Read internal/engine/forward/envelope.go
→ Read internal/engine/forward/http_test.go
→ Read internal/engine/forward/nats_test.go

$ rg -n "forward\.|Forwarder|Forward\(" --glob '*.go' --glob '!internal/engine/forward/*' -l | head -20
cmd/astrate/main.go
cmd/astrate/forward_test.go
cmd/astrate/forward_nats_test.go
cmd/astrate/newnats_default.go
cmd/astrate/newnats_nats.go
internal/config/config.go
internal/engine/triggers_change_test.go
internal/engine/router.go
internal/engine/engine.go
internal/engine/serverdata_test.go
internal/engine/triggers/match_test.go
internal/engine/triggers/actions.go
internal/engine/triggers/actions_test.go
internal/engine/triggers_change_integration_test.go
internal/engine/devicedata_test.go

$ rg -n "Forwarder|forward\." --glob '*.go' --glob '!internal/engine/forward/*' | head -40
cmd/astrate/main.go:116:	fwd, err := newForwarder(cfg, log)
cmd/astrate/main.go:122:		Forwarder:       fwd,
cmd/astrate/main.go:248:// newForwarder builds the Forwarder that receives trigger actions which are
cmd/astrate/main.go:254:// the executor tests its Forwarder against nil, and a typed nil would satisfy
cmd/astrate/main.go:256:func newForwarder(cfg config.Config, log *slog.Logger) (triggers.Forwarder, error) {
cmd/astrate/main.go:262:		h, err := forward.New(forward.Config{
cmd/astrate/main.go:273:		h, err := newNATSForwarder(f)
cmd/astrate/forward_test.go:12:// quietLogger keeps newForwarder's info line out of the test output.
cmd/astrate/forward_test.go:17:// TestForwarderDisabledIsUntypedNil is the one rule in newForwarder that a
cmd/astrate/forward_test.go:19:// decides whether to forward by comparing its Forwarder against nil, so
cmd/astrate/forward_test.go:20:// returning a typed nil *forward.HTTP would pass that comparison and then
cmd/astrate/forward_test.go:24:func TestForwarderDisabledIsUntypedNil(t *testing.T) {
cmd/astrate/forward_test.go:25:	f, err := newForwarder(config.Config{}, quietLogger())
cmd/astrate/forward_test.go:27:		t.Fatalf("newForwarder: %v", err)
cmd/astrate/forward_test.go:34:// TestForwarderHTTPBuilt checks the configured kind produces a usable
cmd/astrate/forward_test.go:37:func TestForwarderHTTPBuilt(t *testing.T) {
cmd/astrate/forward_test.go:43:	f, err := newForwarder(cfg, quietLogger())
cmd/astrate/forward_test.go:45:		t.Fatalf("newForwarder: %v", err)
cmd/astrate/forward_test.go:52:// TestForwarderRejectsBadURL: config.validate normally catches this, but
cmd/astrate/forward_test.go:53:// newForwarder must not silently ignore a construction failure if a config
cmd/astrate/forward_test.go:55:func TestForwarderRejectsBadURL(t *testing.T) {
cmd/astrate/forward_test.go:58:	if _, err := newForwarder(cfg, quietLogger()); err == nil {
cmd/astrate/forward_test.go:63:// TestForwarderNATSWithoutBuildTag: this file (and newnats_default.go) build
cmd/astrate/forward_test.go:66:func TestForwarderNATSWithoutBuildTag(t *testing.T) {
cmd/astrate/forward_test.go:72:	_, err := newForwarder(cfg, quietLogger())
cmd/astrate/forward_nats_test.go:15:// TestForwarderNATSBuilt proves the config→newForwarder→forward.NewNATS wire
cmd/astrate/forward_nats_test.go:17:// TestForwarderNATSWithoutBuildTag's untagged-build error path.
cmd/astrate/forward_nats_test.go:18:func TestForwarderNATSBuilt(t *testing.T) {
cmd/astrate/forward_nats_test.go:32:	f, err := newForwarder(cfg, quietLogger())
cmd/astrate/forward_nats_test.go:34:		t.Fatalf("newForwarder: %v", err)
cmd/astrate/newnats_default.go:16:// newNATSForwarder reports that this binary was not built with -tags nats, so
cmd/astrate/newnats_default.go:17:// triggers.forward.kind = "nats" fails at boot rather than silently
cmd/astrate/newnats_default.go:19:func newNATSForwarder(_ config.ForwardConfig) (triggers.Forwarder, error) {
cmd/astrate/newnats_default.go:20:	return nil, fmt.Errorf("triggers.forward.kind \"nats\" requires building astrate with -tags nats")
cmd/astrate/newnats_nats.go:15:// newNATSForwarder builds the NATS bus forwarder. Only compiled into a binary
cmd/astrate/newnats_nats.go:17:func newNATSForwarder(f config.ForwardConfig) (triggers.Forwarder, error) {
cmd/astrate/newnats_nats.go:18:	return forward.NewNATS(forward.NATSConfig{URL: f.URL, Subject: f.Subject})
internal/config/config.go:160:// ForwardConfig selects the Forwarder that receives trigger actions which
internal/config/config.go:342:		return fmt.Errorf("config: triggers.forward.kind %q is not one of \"\"|\"http\"|\"nats\"", c.Triggers.Forward.Kind)
internal/config/config.go:349:			return fmt.Errorf("config: triggers.forward.url is required when kind is \"http\"")

→ Read cmd/astrate/main.go [offset=95, limit=80]
→ Read cmd/astrate/main.go [offset=240, limit=50]
→ Read internal/engine/triggers/actions.go

$ rg -n "Forward|forward" internal/config/config.go | head -40
157:	Forward ForwardConfig `toml:"forward"`
160:// ForwardConfig selects the Forwarder that receives trigger actions which
161:// are not HTTP webhooks (docs/DESIGN.md §1.1). Kind "" disables forwarding
162:// entirely; "http" selects the HTTP bus forwarder; "nats" selects the NATS
163:// bus forwarder (only usable in a binary built with -tags nats — see
165:type ForwardConfig struct {
339:	switch c.Triggers.Forward.Kind {
342:		return fmt.Errorf("config: triggers.forward.kind %q is not one of \"\"|\"http\"|\"nats\"", c.Triggers.Forward.Kind)
344:	// When kind is "" forwarding is disabled; do not validate url, method, or
345:	// static_headers — a stale url left behind after disabling forwarding must
347:	if c.Triggers.Forward.Kind == "http" {
348:		if c.Triggers.Forward.URL == "" {
349:			return fmt.Errorf("config: triggers.forward.url is required when kind is \"http\"")
351:		u, err := url.Parse(c.Triggers.Forward.URL)
353:			return fmt.Errorf("config: triggers.forward.url %q is not an absolute http or https URL", c.Triggers.Forward.URL)
355:		if c.Triggers.Forward.Method != "" {
356:			switch strings.ToUpper(c.Triggers.Forward.Method) {
359:				return fmt.Errorf("config: triggers.forward.method %q is not a recognised HTTP method", c.Triggers.Forward.Method)
364:	// server) happens at construction time in forward.NewNATS, the same way
368:	if c.Triggers.Forward.Kind == "nats" {
369:		if c.Triggers.Forward.URL == "" {
370:			return fmt.Errorf("config: triggers.forward.url is required when kind is \"nats\"")
372:		if c.Triggers.Forward.Subject == "" {
373:			return fmt.Errorf("config: triggers.forward.subject is required when kind is \"nats\"")

$ rg -n "func shutdown" -A 40 cmd/astrate/main.go
482:func shutdown(srv *http.Server, b *broker.Broker, flowSvc *flowapi.Service, e *engine.Engine, log *slog.Logger) {
483-	if srv != nil {
484-		log.Info("draining http listener")
485-		stage("http shutdown", log, srv.Shutdown)
486-	}
487-	log.Info("stopping broker")
488-	if err := b.Close(); err != nil {
489-		log.Warn("broker close", "error", err)
490-	}
491-	if flowSvc != nil {
492-		log.Info("stopping flows")
493-		mgr := flowSvc.Manager()
494-		stage("flow shutdown", log, mgr.Shutdown)
495-		stage("flows marked stopped", log, func(ctx context.Context) error {
496-			flowSvc.MarkRunningFlowsStopped(ctx)
497-			return nil
498-		})
499-	}
500-	drainEngine(e, log)
501-	log.Info("shutdown complete")
502-}
503-
504-// stage runs one drain stage under a fresh shutdownTimeout context and logs
505-// whatever it returns, so an exhausted stage never leaves the next one with a
506-// dead context.
507-func stage(label string, log *slog.Logger, fn func(context.Context) error) {
508-	sctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
509-	defer cancel()
510-	if err := fn(sctx); err != nil {
511-		log.Warn(label, "error", err)
512-	}
513-}
514-
515-func drainEngine(e *engine.Engine, log *slog.Logger) {
516-	log.Info("draining engine")
517-	stage("engine drain", log, e.Drain)
518-}
519-
520-// loadSealer builds the CA-key sealer: the configured master-key file feeds the
521-// store's env-based loader (ASTRATE_MASTER_KEY[_FILE]); decoding rules are
522-// shared with internal/store.

$ rg -n "func validateStaticHeaders" -A 45 internal/engine/triggers/*.go
internal/engine/triggers/validation.go:140:func validateStaticHeaders(e *FieldErrors, headers map[string]string) {
internal/engine/triggers/validation.go-141-	for name := range headers {
internal/engine/triggers/validation.go-142-		blocked := blockedHeaderNames[strings.ToLower(strings.TrimSpace(name))]
internal/engine/triggers/validation.go-143-		if blocked {
internal/engine/triggers/validation.go-144-			e.add("http_static_headers", "must contain only allowed http headers")
internal/engine/triggers/validation.go-145-			break
internal/engine/triggers/validation.go-146-		}
internal/engine/triggers/validation.go-147-	}
internal/engine/triggers/validation.go-148-	total := 0
internal/engine/triggers/validation.go-149-	for name, value := range headers {
internal/engine/triggers/validation.go-150-		total += len(name) + 2 + len(value)
internal/engine/triggers/validation.go-151-	}
internal/engine/triggers/validation.go-152-	if total >= maxHeadersBytes {
internal/engine/triggers/validation.go-153-		e.add("http_static_headers", fmt.Sprintf("headers total size must be lower than %d", maxHeadersBytes))
internal/engine/triggers/validation.go-154-	}
internal/engine/triggers/validation.go-155-}
internal/engine/triggers/validation.go-156-
internal/engine/triggers/validation.go-157-// validateTemplate enforces the template rune-length cap regardless of
internal/engine/triggers/validation.go-158-// template_type (upstream never validates template_type itself).
internal/engine/triggers/validation.go-159-func validateTemplate(e *FieldErrors, template string) {
internal/engine/triggers/validation.go-160-	if utf8.RuneCountInString(template) > maxTemplateRunes {
internal/engine/triggers/validation.go-161-		e.add("template", fmt.Sprintf("should be at most %d character(s)", maxTemplateRunes))
internal/engine/triggers/validation.go-162-	}
internal/engine/triggers/validation.go-163-}

→ Read internal/config/config.go [offset=150, limit=40]

$ rg -n "blockedHeaderNames" -B 12 -A 20 internal/engine/triggers/validation.go | head -60
7-	"strings"
8-	"unicode/utf8"
9-)
10-
11-// Upstream http_action.ex v1.2.0 limits, probe-frozen (issue #63).
12-const (
13-	minURLRunes      = 8
14-	maxURLRunes      = 8192
15-	maxHeadersBytes  = 8192 // Σ(byte_size(name)+2+byte_size(value)) must be < this
16-	maxTemplateRunes = 1048576
17-)
18-
19:// blockedHeaderNames is upstream's hop-by-hop/sensitive header blocklist,
20-// compared against strings.ToLower(strings.TrimSpace(name)).
21:var blockedHeaderNames = map[string]bool{
22-	"connection":           true,
23-	"content-length":       true,
24-	"date":                 true,
25-	"host":                 true,
26-	"te":                   true,
27-	"upgrade":              true,
28-	"x-forwarded-for":      true,
29-	"x-forwarded-host":     true,
30-	"x-forwarded-proto":    true,
31-	"sec-websocket-accept": true,
32-	"proxy-authorization":  true,
33-	"proxy-authenticate":   true,
34-}
35-
36-// FieldErrors carries an upstream-shaped field→messages map for one part of a
37-// trigger definition. errors.Is/As-friendly so the realm layer can render the
38-// nested changeset envelope without string parsing.
39-type FieldErrors struct {
40-	Part   string // e.g. "action"
41-	Fields map[string][]string
--
130-
131-// isHTTPURL reports whether raw parses as an absolute URL with an http(s)
132-// scheme and a non-empty host (upstream's scheme+host check).
133-func isHTTPURL(raw string) bool {
134-	u, err := url.Parse(raw)
135-	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
136-}
137-
138-// validateStaticHeaders applies the modern-branch header checks: the
139-// blocklist first, then the total-size ceiling; both can fail together.
140-func validateStaticHeaders(e *FieldErrors, headers map[string]string) {
141-	for name := range headers {
142:		blocked := blockedHeaderNames[strings.ToLower(strings.TrimSpace(name))]
143-		if blocked {
144-			e.add("http_static_headers", "must contain only allowed http headers")
145-			break
146-		}
147-	}
148-	total := 0
149-	for name, value := range headers {
150-		total += len(name) + 2 + len(value)
151-	}
152-	if total >= maxHeadersBytes {
153-		e.add("http_static_headers", fmt.Sprintf("headers total size must be lower than %d", maxHeadersBytes))

$ rg -n "forward" internal/config/config_test.go | head -30
244:[triggers.forward]
262:[triggers.forward]
293:[triggers.forward]
320:			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"amqp\"",
324:			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"nats\"",
328:			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"nats\"\nurl=\"nats://bus.example:4222\"",
338:			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"http\"",
342:			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"http\"\nurl=\"nope\"",
346:			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"http\"\nurl=\"ftp://bus.example/x\"",
350:			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"http\"\nurl=\"http://\"",
354:			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"http\"\nurl=\"https://bus.example/x\"\nmethod=\"FETCH\"",
358:			`[database]` + "\ndsn=\"x\"\n[mqtt]\ninsecure_dev_mode=true\n[triggers.forward]\nkind=\"http\"\nurl=\"https://bus.example/x\"\nmethod=\"post\"",

$ rg -n "nats" go.mod; echo "--- module cache ---"; ls -d $(go env GOMODCACHE)/github.com/nats-io 2>/dev/null
26:	github.com/nats-io/nats.go v1.52.0
28:	github.com/testcontainers/testcontainers-go/modules/nats v0.43.0
72:	github.com/nats-io/nkeys v0.4.15 // indirect
73:	github.com/nats-io/nuid v1.0.1 // indirect
--- module cache ---
/root/go/pkg/mod/github.com/nats-io

$ rg -n "forward|Forward|bus" .mule/todo.md | head -40; echo "=== total lines ==="; wc -l .mule/todo.md
2:- [x] triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
3:- [x] compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).
86:- [ ] broker-external-bus-intake [legion]: implement the second Intake implementation named in the TODO at internal/broker/intake.go:56-62 — an external-bus-backed (e.g. NATS JetStream) Intake for multi-instance deployment / restart survival, reproducing the embedded broker's per-device ordering, deferred-ack backpressure and QoS 0 drop semantics; an implementation of this interface needs containerised integration tests, so it runs on the Legion Go.
131:- [x] server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.
201:- [x] docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate_devices_total{realm="test"} 42`, but no such series is ever registered — the registry carries `astrate_broker_sessions` (internal/observability/metrics.go:54-62), `astrate_db_pool_{acquired,idle,total,max}_conns` (metrics.go:65-79), plus the engine/flow/trigger families (`astrate_engine_*` engine/router.go:411-444, `astrate_engine_stream_dropped_total` engine/stream/bus.go:119, `astrate_engine_trigger_*` engine/triggers/actions.go:296-300, `astrate_flow_router_*` flow/router.go:260-276) and the standard `go_*`/`process_*` runtime collectors — so the example advertises a scrape result the wire can never return. Swap it for a real series (e.g. `astrate_broker_sessions 0`), keeping the HELP/TYPE comment style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
219:- [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
=== total lines ===
230 .mule/todo.md

$ rg -n "^func \(nc \*Conn\) Publish\b" -A 30 /root/go/pkg/mod/github.com/nats-io/nats.go@v1.52.0/nats.go | head -50
4229:func (nc *Conn) Publish(subj string, data []byte) error {
4230-	return nc.publish(subj, _EMPTY_, false, nil, data)
4231-}
4232-
4233-// Header represents the optional Header for a NATS message,
4234-// based on the implementation of http.Header.
4235-type Header map[string][]string
4236-
4237-// Add adds the key, value pair to the header. It is case-sensitive
4238-// and appends to any existing values associated with key.
4239-func (h Header) Add(key, value string) {
4240-	h[key] = append(h[key], value)
4241-}
4242-
4243-// Set sets the header entries associated with key to the single
4244-// element value. It is case-sensitive and replaces any existing
4245-// values associated with key.
4246-func (h Header) Set(key, value string) {
4247-	h[key] = []string{value}
4248-}
4249-
4250-// Get gets the first value associated with the given key.
4251-// It is case-sensitive.
4252-func (h Header) Get(key string) string {
4253-	if h == nil {
4254-		return _EMPTY_
4255-	}
4256-	if v := h[key]; v != nil {
4257-		return v[0]
4258-	}
4259-	return _EMPTY_

$ rg -n "func \(nc \*Conn\) publish" -A 60 /root/go/pkg/mod/github.com/nats-io/nats.go@v1.52.0/nats.go | head -80
4424:func (nc *Conn) publish(subj, reply string, validateReply bool, hdr, data []byte) error {
4425-	if nc == nil {
4426-		return ErrInvalidConnection
4427-	}
4428-	if !nc.Opts.SkipSubjectValidation {
4429-		if err := validateSubject(subj); err != nil {
4430-			return err
4431-		}
4432-		if validateReply {
4433-			if err := validateSubject(reply); err != nil {
4434-				return ErrBadSubject
4435-			}
4436-		}
4437-	} else if subj == _EMPTY_ {
4438-		return ErrBadSubject
4439-	}
4440-	nc.mu.Lock()
4441-
4442-	// Check if headers attempted to be sent to server that does not support them.
4443-	if len(hdr) > 0 && !nc.info.Headers {
4444-		nc.mu.Unlock()
4445-		return ErrHeadersNotSupported
4446-	}
4447-
4448-	if nc.isClosed() {
4449-		nc.mu.Unlock()
4450-		return ErrConnectionClosed
4451-	}
4452-
4453-	if nc.isDrainingPubs() {
4454-		nc.mu.Unlock()
4455-		return ErrConnectionDraining
4456-	}
4457-
4458-	// Proactively reject payloads over the threshold set by server.
4459-	msgSize := int64(len(data) + len(hdr))
4460-	// Skip this check if we are not yet connected (RetryOnFailedConnect)
4461-	if !nc.initc && msgSize > nc.info.MaxPayload {
4462-		nc.mu.Unlock()
4463-		return ErrMaxPayload
4464-	}
4465-
4466-	// Check if we are reconnecting, and if so check if
4467-	// we have exceeded our reconnect outbound buffer limits.
4468-	if nc.bw.atLimitIfUsingPending() {
4469-		nc.mu.Unlock()
4470-		return ErrReconnectBufExceeded
4471-	}
4472-
4473-	var mh []byte
4474-	if hdr != nil {
4475-		mh = nc.scratch[:len(_HPUB_P_)]
4476-	} else {
4477-		mh = nc.scratch[1:len(_HPUB_P_)]
4478-	}
4479-	mh = append(mh, subj...)
4480-	mh = append(mh, ' ')
4481-	if reply != "" {
4482-		mh = append(mh, reply...)
4483-		mh = append(mh, ' ')
4484-	}

$ rg -n "func \(nc \*Conn\) processErr" -A 40 /root/go/pkg/mod/github.com/nats-io/nats.go@v1.52.0/nats.go | head -55
4181:func (nc *Conn) processErr(ie string) {
4182-	// Trim, remove quotes
4183-	ne := normalizeErr(ie)
4184-	// convert to lower case.
4185-	e := strings.ToLower(ne)
4186-
4187-	var close bool
4188-
4189-	// FIXME(dlc) - process Slow Consumer signals special.
4190-	if e == STALE_CONNECTION {
4191-		close = nc.processOpErr(ErrStaleConnection, false)
4192-	} else if e == MAX_CONNECTIONS_ERR {
4193-		close = nc.processOpErr(ErrMaxConnectionsExceeded, false)
4194-	} else if e == MAX_ACCOUNT_CONNECTIONS_ERR {
4195-		close = nc.processOpErr(ErrMaxAccountConnectionsExceeded, false)
4196-	} else if strings.HasPrefix(e, PERMISSIONS_ERR) {
4197-		nc.processTransientError(fmt.Errorf("%w: %s", ErrPermissionViolation, ne))
4198-	} else if strings.HasPrefix(e, MAX_SUBSCRIPTIONS_ERR) {
4199-		nc.processTransientError(ErrMaxSubscriptionsExceeded)
4200-	} else if authErr := checkAuthError(e); authErr != nil {
4201-		nc.mu.Lock()
4202-		close = nc.processAuthError(authErr)
4203-		nc.mu.Unlock()
4204-	} else {
4205-		close = true
4206-		nc.mu.Lock()
4207-		nc.err = errors.New("nats: " + ne)
4208-		nc.mu.Unlock()
4209-	}
4210-	if close {
4211-		nc.close(CLOSED, true, nil)
4212-	}
4213-}
4214-
4215-// kickFlusher will send a bool on a channel to kick the
4216-// flush Go routine to flush data to the server.
4217-func (nc *Conn) kickFlusher() {
4218-	if nc.bw != nil {
4219-		select {
4220-		case nc.fch <- struct{}{}:
4221-		default:

$ rg -n "ErrorHandler|AsyncErrorCB" /root/go/pkg/mod/github.com/nats-io/nats.go@v1.52.0/nats.go | rg -n "if |!= nil|Default" | head -20
12:1854:	if nc.Opts.AsyncErrorCB == nil {
14:3024:			if asyncErrorCB := nc.Opts.AsyncErrorCB; asyncErrorCB != nil {
15:3043:			if asyncErrorCB := nc.Opts.AsyncErrorCB; asyncErrorCB != nil {
16:3056:				if asyncErrorCB := nc.Opts.AsyncErrorCB; asyncErrorCB != nil {
17:3719:			if asyncErrorCB := nc.Opts.AsyncErrorCB; asyncErrorCB != nil {
18:3897:		if asyncErrorCB := nc.Opts.AsyncErrorCB; asyncErrorCB != nil {
19:3941:	if asyncErrorCB := nc.Opts.AsyncErrorCB; asyncErrorCB != nil {
20:3955:		if asyncErrorCB := nc.Opts.AsyncErrorCB; asyncErrorCB != nil {
21:4002:				if asyncErrorCB := nc.Opts.AsyncErrorCB; asyncErrorCB != nil {
22:5250:					if errCB := nc.Opts.AsyncErrorCB; errCB != nil {

$ sed -n '1840,1870p' /root/go/pkg/mod/github.com/nats-io/nats.go@v1.52.0/nats.go
// the Secure boolean is true.
	if nc.Opts.TLSHandshakeFirst {
		nc.Opts.Secure = true
	}

	if err := nc.setupServerPool(); err != nil {
		return nil, err
	}

	// Create the async callback handler.
	nc.ach = &asyncCallbacksHandler{}
	nc.ach.cond = sync.NewCond(&nc.ach.mu)

	// Set a default error handler that will print to stderr.
	if nc.Opts.AsyncErrorCB == nil {
		nc.Opts.AsyncErrorCB = defaultErrHandler
	}

	// Create reader/writer
	nc.newReaderWriter()

	// Spin up the async cb dispatcher before connect so it's ready
	// to handle callbacks, especially when RetryOnFailedConnect is used
	// and initial connection fails.
	go nc.ach.asyncCBDispatcher()

	connectionEstablished, err := nc.connect()
	if err != nil {
		nc.ach.close()
		return nil, err
	}

$ sed -n '1,40p' .mule/reviews/swagger-2026-09-24.md; echo "=== for-giulio tail ==="; tail -12 .mule/for-giulio.md; echo "=== todo tail ==="; tail -6 .mule/todo.md
# Review — internal/swagger (2026-09-24)

First review of this area, so no prior file to extend.

## What I read

- `internal/swagger/swagger.go` (all 43 lines)
- `internal/swagger/swagger_test.go` (all 192 lines)
- `cmd/astrate/main.go:414` (the single production `swagger.Mount` call site)
- `docs/*.go` (to confirm `SwaggerUI` embeds `swagger-ui/*` and `APIYAML` embeds `api/*.yaml`)
- Grep for `swagger.Mount` and `Specs()` callers — Mount is called only in main; `Specs()` has no production callers (test support only).

## What I found

The package is small and its behaviours are all asserted: the `/swagger` redirect,
the served embedded UI, every `/api/*.yaml` spec served, `Specs()` completeness
against a fresh embed walk, and the realm-management 403 contract. Test coverage is
genuinely good for a 43-line file.

One real (if minor) defect:

**`internal/swagger/swagger.go:17-18` — `fs.Sub` errors silently ignored.** Both
`uiRoot, _ := fs.Sub(...)` and `apiRoot, _ := fs.Sub(...)` discard their error. The
embeds are fixed at compile time today, but the moment anyone renames the directory
inside `docs/` (say `swagger-ui` to `static`), `fs.Sub` fails, `Mount` serves an
empty/invalid tree, and `/swagger/` plus `/api/` degrade to silent 404s in
production with no log. That is an unguarded path that fails quietly. Proposal below
makes it verifiable by injecting an FS.

## Suggested task (appended to todo.md)

- swagger-sub-failfast: make Mount fail fast when fs.Sub over docs.SwaggerUI/APIYAML
  fails (inject fs via an internal MountFS/mount-override so a broken-FS test can
  prove the panic/error), internal/swagger/swagger.go

## Decided NOT to propose, and why

- `Specs()` ignoring its own `fs.WalkDir` error (swagger.go:33): the function has no
  production caller — it exists to feed tests. Not worth a change.
- `TestRealmManagement403`'s line-regex YAML parsing being fragile in principle: it
=== for-giulio tail ===
`config.validate` accepts `triggers.forward.kind = "nats"` (internal/config/config.go:339-343)
and the untagged binary then fails at boot with the build-tag message
(cmd/astrate/newnats_default.go). I measured it today: `go build -tags nats ./cmd/...`
and `go vet -tags nats ./cmd/...` both exit 0, so there is nothing broken — just nothing
keeping it unbroken. Worth a `make build-nats` (and the NATS test tier, which needs a
container, so it would be a `[legion]` line) in the gate; not a queue line because the
only test I can imagine for it shells out to `go build`, and the mule's own check strips
the implementation and the guard together, so it would never fail.
Not queued for the same reason: `IdleTimeout` is missing on the HTTP server
(cmd/astrate/main.go:221) and `-healthcheck` ignores the config file
(cmd/astrate/main.go:531) — both real, both yours to call, neither has a test worth
writing. See `.mule/reviews/cmd-astrate-2026-09-25.md`.
=== todo tail ===
- [x] cmd-loadsealer-masterkeyfile-test [auto]: add a container-free test for `loadSealer` (cmd/astrate/main.go:515-526) — the `security.master_key_file` to env hand-off (`os.Setenv(store.EnvMasterKeyFile, ...)`, 517) is the only wiring that makes a configured key file reach the store loader, and dropping or misspelling that branch leaves the file silently ignored and boot dying with the generic "master key" error (523) even though the config named a valid file. Assert a 32-byte key file under t.TempDir() loads, the inline ASTRATE_MASTER_KEY path still wins, and a missing key returns the error naming all three env vars. t.Setenv both ASTRATE_MASTER_KEY and ASTRATE_MASTER_KEY_FILE first so the os.Setenv is restored at cleanup.
- [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion Go — the device-delete leg added by docs-sync-rm-delete-device-async-operation-param (register a device, `DELETE /devices/{id}?async_operation=…` answers 204 and the row is gone on both values) compiles clean under `go vet -tags integration ./internal/realm/` but has never been executed, because the Pi has no database; until it runs, TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go is the only assertion of that leg that has actually been verified.
- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing

$ mkdir -p /tmp/opencode/fwdprobe && cat > /tmp/opencode/fwdprobe/main.go <<'EOF'
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/http/httputil"
	"strings"
)

type Fwd struct{ url, method string; static map[string]string }

func (f *Fwd) Forward(ctx context.Context, realm string, action json.RawMessage) error {
	req, _ := http.NewRequestWithContext(ctx, f.method, f.url, strings.NewReader(`{"a":1}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Astarte-Realm", realm)
	for k, v := range f.static {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("forward: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func main() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer srv.Close()

	// (1) invalid static header name: does it fail per delivery?
	f := &Fwd{url: srv.URL, method: "POST", static: map[string]string{"X Bad Name": "v"}}
	err := f.Forward(context.Background(), "r")
	fmt.Printf("invalid header name  -> err=%v\n", err)
	fmt.Printf("  canonicalised key  -> %q\n", textproto.CanonicalMIMEHeaderKey("X Bad Name"))

	// (2) static header overriding a fixed one (documented at http.go:110)
	var seen string
	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.Header.Get("Astarte-Realm")
	}))
	defer srv2.Close()
	f2 := &Fwd{url: srv2.URL, method: "POST", static: map[string]string{"astarte-realm": "spoofed"}}
	if err := f2.Forward(context.Background(), "real-realm"); err != nil {
		fmt.Println("override err:", err)
	}
	fmt.Printf("static override wins -> server saw Astarte-Realm=%q\n", seen)

	// (3) does DefaultClient follow a cross-host redirect, re-sending body+headers?
	var dst *http.Request
	dstSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := httputil.DumpRequest(r, false)
		dst = r.Clone(context.Background())
		_ = b
		w.WriteHeader(200)
	}))
	defer dstSrv.Close()
	redir := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, dstSrv.URL+"/sink", http.StatusTemporaryRedirect)
	}))
	defer redir.Close()
	f3 := &Fwd{url: redir.URL, method: "POST"}
	err = f3.Forward(context.Background(), "leaky-realm")
	fmt.Printf("redirect -> err=%v  dst saw realm=%q method=%q\n", err, dst.Header.Get("Astarte-Realm"), dst.Method)
}
EOF
cd /tmp/opencode/fwdprobe && go mod init probe >/dev/null 2>&1; go run . 2>&1 | tail -20
# probe
./main.go:37:41: not enough arguments in call to f.Forward
	have (context.Context, string)
	want (context.Context, string, json.RawMessage)
./main.go:48:45: not enough arguments in call to f2.Forward
	have (context.Context, string)
	want (context.Context, string, json.RawMessage)
./main.go:67:41: not enough arguments in call to f3.Forward
	have (context.Context, string)
	want (context.Context, string, json.RawMessage)

$ cd /tmp/opencode/fwdprobe && sed -i 's/func (f \*Fwd) Forward(ctx context.Context, realm string, action json.RawMessage) error {/func (f *Fwd) Forward(ctx context.Context, realm string, _ json.RawMessage) error {/' main.go && go run . 2>&1 | tail -20
# probe
./main.go:37:41: not enough arguments in call to f.Forward
	have (context.Context, string)
	want (context.Context, string, json.RawMessage)
./main.go:48:45: not enough arguments in call to f2.Forward
	have (context.Context, string)
	want (context.Context, string, json.RawMessage)
./main.go:67:41: not enough arguments in call to f3.Forward
	have (context.Context, string)
	want (context.Context, string, json.RawMessage)

$ cd /tmp/opencode/fwdprobe && sed -i 's/f\.Forward(context.Background(), "r")/f.Forward(context.Background(), "r", nil)/; s/f2\.Forward(context.Background(), "real-realm")/f2.Forward(context.Background(), "real-realm", nil)/; s/f3\.Forward(context.Background(), "leaky-realm")/f3.Forward(context.Background(), "leaky-realm", nil)/' main.go && go run . 2>&1 | tail -20
invalid header name  -> err=forward: Post "http://127.0.0.1:42021": net/http: invalid header field name "X Bad Name"
  canonicalised key  -> "X Bad Name"
static override wins -> server saw Astarte-Realm="spoofed"
redirect -> err=<nil>  dst saw realm="leaky-realm" method="POST"

$ cd /tmp/opencode/fwdprobe && cat > probe2.go <<'EOF'
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

type Fwd2 struct {
	url    string
	method string
	static map[string]string
}

func (f *Fwd2) Forward(ctx context.Context, realm string, _ json.RawMessage) error {
	req, _ := http.NewRequestWithContext(ctx, f.method, f.url, strings.NewReader(`{"a":1}`))
	req.Header.Set("Astarte-Realm", realm)
	for k, v := range f.static {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("forward: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

func probe2() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	defer srv.Close()
	for name, hv := range map[string]map[string]string{
		"value with newline": {"X-Foo": "a\nb"},
		"value with space":   {"X-Foo": "a b"},
		"name with colon":    {"X:Foo": "v"},
		"name lowercase ok":  {"x-foo": "v"},
	} {
		f := &Fwd2{url: srv.URL, method: "POST", static: hv}
		err := f.Forward(context.Background(), "r", nil)
		fmt.Printf("%-20s -> err=%v\n", name, err)
	}
}
EOF
sed -i 's/^func main() {/func main() {\n\tprobe2()/' main.go && go run . 2>&1 | tail -12
value with space     -> err=<nil>
name with colon      -> err=forward: Post "http://127.0.0.1:39663": net/http: invalid header field name "X:Foo"
name lowercase ok    -> err=<nil>
value with newline   -> err=forward: Post "http://127.0.0.1:39663": net/http: invalid header field value for "X-Foo"
invalid header name  -> err=forward: Post "http://127.0.0.1:45609": net/http: invalid header field name "X Bad Name"
  canonicalised key  -> "X Bad Name"
static override wins -> server saw Astarte-Realm="spoofed"
redirect -> err=<nil>  dst saw realm="leaky-realm" method="POST"

$ head -30 .mule/for-giulio.md; echo "=== wc ==="; wc -l .mule/for-giulio.md; echo "=== git log forward ==="; git log --oneline -8 -- internal/engine/forward/ cmd/astrate/newnats_nats.go; echo "=== marshalEnvelope callers/tests ==="; rg -n "marshalEnvelope" --glob '*.go' | head
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **github-issues triage run, 2026-09-26: nothing proposable (as every run since 2026-09-05), and the 19-issue alarm pile is a false alarm — it says the mule has landed nothing in 19 days and it has landed something every day.** 23 open issues, the same set plus **#112** (today, 09-26 11:04:46Z, "nothing has landed in 14h"). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction.
  **The "idle streak" that 19 consecutive entries here have repeated is wrong, and this run retracts it.** `.mule/log.md` — the script's own row-per-attempt log, starting exactly at 2026-09-04, the day the streak is claimed to begin — records attempts on 20 of the 22 days to 09-25 and `done` lands on 17 of them (3–8/day: 6, 10, 8, 6, 8, 7, 7, 6, 3, 5, 6, 5, 3, 0, 5, 3, 0, 5, 4, 7 for 09-04→09-25). The last three landed shas — `84a050e`, `ed1cd72`, `f965f39` — are all reachable from `origin/mule/queue` (`git branch -r --contains`), so the work landed *and* pushed, which is the exact failure the alarm exists to catch.
  **What the alarm really measures, and why it fires daily by construction.** `check_pulse` (mule.sh:739-774) compares `.heartbeat` against `MULE_IDLE_ALARM_HOURS` (default 8, mule.sh:740), and `beat` is called only on a `done` land or a `checked` run (mule.sh:578, 540) — a refill does not beat. Last beat: **2026-09-25T20:19:13Z**; the pulse that files the issue runs ~10:54–11:12Z. 20:19Z → 11:04Z is ~14.75h, comfortably over 8h, so the arithmetic is right and the conclusion is wrong: with lands bunched into an evening window, the gap to the next morning's tick is >8h *every day*. It is a mis-set threshold, not a dead mule. The two days in this window that landed nothing (09-19, 09-22) are the only ones where the alarm had anything real to say.
  **Two candidate fixes, both yours — how many hours of silence is a real problem is a judgement, so I did not queue either:** (a) raise `MULE_IDLE_ALARM_HOURS` to ~24 in `mule.service`; or (b) make the pulse read `.mule/log.md` instead of the heartbeat and alarm only when a *whole day* has no `done` row — more honest, since that file is the record of what actually landed, and it needs no threshold tuning as the cadence changes.
  **Proposal: close #94–#112, all nineteen.** Each is a one-day false alarm superseded by the next; closing is inert either way (nothing reconciles them, as established on 2026-09-25), but 19 open "the mule is dead" issues is exactly the noise that hides a real one. Expect #113 tomorrow unless the threshold changes.
  **Correction on a fact 19 entries have repeated as settled: there is no `.mule/waiting-on.md`.** #92's own comment (2026-09-04) states the issue "has a row in `.mule/waiting-on.md`" and that the weekly `mule-upstream-watch` job reads that file as its first step. The file does not exist here, and `rg 'waiting-on|upstream-watch'` over the whole tree — including gitignored `.mule/` and `tools/` — matches only prose in this file and copies of it under `.mule/reports/log/`. So the escalation that was supposed to catch a stable upstream `v1.4.0` is not in the repo: either it lives in a Pi-side unit outside it, or it was never wired. One command settles it: `systemctl list-timers 'mule-upstream*' --all; grep -rl 'upstream-watch\|waiting-on' /etc/systemd/system /root 2>/dev/null`. If it was never wired, #92 is parked on a condition nothing is watching, and the 2026-09-04 claim that "the mechanism was added because it was missing here" is false. The gate is being caught by hand in the meantime — today's milestone run re-derived it from `gh api .../releases`.

---

- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Status: not started", which is no longer true.** Release sweep (`gh api repos/astarte-platform/astarte/releases`, run this session) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.5` (2026-08-20, prerelease) — byte-identical to what the 2026-09-24 run saw, so nothing shipped upstream in the two days since. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Open non-alarm issues remain exactly **#93** (aclhook comment, `mule-review`, its own recipe path), **#92** (keyAgreement, parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** (FDO, `milestone-4.0`), **#1** (untouched per standing instruction). No new gaps, **no issues filed, no `.mule/todo.md` task lines added** — there is nothing machine-checkable left, and filing anything here would duplicate #92/#93/#78.
  **New, and the reason this run is not a copy of the last thirteen: `.mule/milestones.md:127-128` still reads "Status: not started. First recipe job: triage #47–#89 into an ordered plan (which are audits vs features vs decisions), file sub-issues where work splits, escalate the 'decide' set in one batch."** All of that happened: of #47–#89 every issue is closed except the deliberately-parked #78. The section describes as future work a triage that is finished, and it is the exact text every future run of this recipe reads to pick its target — so the recipe has been re-deriving "the backlog is done, the gate is upstream" from scratch on a timer for three weeks. Suggest rewriting those two lines to state the real position: *1.3 surface delivered (#47–#89 closed); milestone blocked on upstream v1.4.0 going stable; final phase = reconcile the two `docs/UPSTREAM-EXPERIMENTAL.md` rows, answer #92, bump `APICompatVersion`.* That also makes the blocking condition greppable instead of inferred.
  **The recipe's standing item is satisfied, and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still has exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental* with reconcile trigger "upstream 1.4 final". v1.4 is not final, so both rows stay as written: nothing to promote, nothing to deprecate, and the register is accurate. `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision; the milestone-bump issue the recipe asks for **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": the milestone's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is done, and the 1.4 half is unresolved (#92 open, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut. Also still true and already escalated on 2026-09-06, so not repeated: the section's source document `.mule/research/upstream-parity-2026-08-22.md` does not exist in the repo.
  **Standing ask, unchanged, for whenever a stable v1.4.0 ships:** answer #92, reconcile the two experimental rows, run the final-phase `APICompatVersion` bump (and update `docs/COMPATIBILITY.md` deviation #10 to match), then cut the v3.0 tag. One process note while this stays blocked: this is the fifteenth milestone-run entry in this file (fourteen prior, 2026-09-06 → 2026-09-24) and the header asks for a queue rather than a log — worth deleting the superseded duplicates whenever you next prune, or the signal-to-noise here keeps falling.

---

## 2026-09-25 — docs-sync run, surface: **housekeeping** (the least-covered of the five)

=== wc ===
464 .mule/for-giulio.md
=== git log forward ===
3e23d8a M12-08c/08d: NATS bus Forwarder, wired end to end behind -tags nats (#11)
9c6b22b Add an HTTP bus Forwarder that posts a custom trigger action verbatim
=== marshalEnvelope callers/tests ===
internal/engine/forward/envelope.go:16:// marshalEnvelope builds the envelope's JSON bytes. A nil json.RawMessage
internal/engine/forward/envelope.go:19:func marshalEnvelope(realm, trigger string, action, event []byte) ([]byte, error) {
internal/engine/forward/nats.go:45:	payload, err := marshalEnvelope(realm, trigger, action, event)
internal/engine/forward/http.go:99:	raw, err := marshalEnvelope(realm, trigger, action, event)

$ go test ./internal/engine/forward/ 2>&1 | tail -3; echo "=== todo head ==="; sed -n '1,1p' .mule/todo.md | cut -c1-400; echo; echo "=== a recent open line ==="; rg -n "^- \[ \] " .mule/todo.md | tail -3 | cut -c1-500
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
=== todo head ===
- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUps

=== a recent open line ===
203:- [ ] store-devices-inhibit-re-register [legion] [auto]: `RegisterDevice` (internal/store/devices.go:75-91) silently clears the inhibit flag — its `ON CONFLICT ... DO UPDATE SET status = 'registered'` fires for any device with `first_credentials_request IS NULL`, including one the admin inhibited via `SetDeviceInhibited` (devices.go:251-268), which sets `status='inhibited'` on unconfirmed devices too; §5.3 says an inhibited device blocks new credentials and connections, so the re-registra
205:- [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (internal/store/devices.go:130-131, "if several devices share an alias the lowest device ID wins", `ORDER BY id LIMIT 1`) — no test drives it today; add a case in internal/store/devices_test.go that gives two devices the same alias and asserts the lookup resolves the lower ID. Needs the integration DB.
221:- [ ] flow-boot-resumes-stopped-durable-flows [legion] [auto]: RehydrateAutoRestart (cmd/astrate/main.go:216) restarts every flow with auto_restart=true whatever its last status — `ListAutoRestartFlows` filters on `f.auto_restart = true` alone (internal/store/flows.go:125) and nothing ever clears the flag, which is written once at create (flowapi/http.go:141-146 -> service.go:318) and never on stop — so a durable flow stopped through the API is running again after the next boot, while th

$ go vet ./internal/engine/forward/ && go test -count=1 ./internal/engine/forward/ 2>&1 | tail -3
ok  	github.com/astrate-platform/astrate/internal/engine/forward	0.083s

Baseline green. Now writing the review.
← Write .mule/reviews/forward-2026-09-26.md
Wrote file successfully.

Now the four task lines in the existing house style.
$ cat >> .mule/todo.md <<'EOF'
- [ ] forward-static-header-validation [auto]: `New` (internal/engine/forward/http.go:38-71) validates the URL and the method at boot precisely because of the rule its own comment states at 42-45 — "an unusable endpoint must fail here rather than surface once per delivery" — but `cfg.StaticHeaders` is stored unvalidated (http.go:68) and applied per delivery with `req.Header.Set` (http.go:110-112), and net/http checks header names and values only at write time. Measured on the real net/http path with a throwaway program: a name of `X Bad Name` or `X:Foo` fails per delivery with `net/http: invalid header field name`, and a value containing a newline fails with `net/http: invalid header field value for "X-Foo"` (a plain space in a value is legal). So one typo in `triggers.forward.static_headers` makes the process boot clean, accept triggers, and then count `astrate_engine_trigger_deliveries_total{outcome="failed"}` and log "custom trigger action failed" (internal/engine/triggers/actions.go:499) on every custom action forever. Validate in `New`: the RFC 7230 field-name grammar is exactly the token grammar `validMethod` (http.go:77-95) already implements, so reuse that predicate per name (rename it if you like) and add a CR/LF rejection per value. Add a table case to `TestNewRejectsBadConfig` (internal/engine/forward/http_test.go:271) asserting `New` returns an error for a bad name and a bad value and that no request is attempted; this test fails today because `New` returns nil error for both.
- [ ] forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
- [ ] forward-status-error-body [auto]: a non-2xx forward returns `forward: status %d` and nothing else (internal/engine/forward/http.go:121-123) after copying the **entire** response body into `io.Discard` with no limit (http.go:117-120) — the one thing an operator needs to debug a 500 from the bus is destroyed, and the sibling webhook request in the same codebase does bound its drain with `io.LimitReader(resp.Body, 1<<20)` (internal/engine/triggers/actions.go:626), so the two near-identical request paths have drifted on exactly that line. Read a bounded prefix (e.g. 512 bytes) of the body into the returned error and bound the discard the way actions.go:626 does. Add a case beside `TestStatusTable` (internal/engine/forward/http_test.go:164) where the handler answers 500 with body `{"error":"boom"}` and assert the error mentions both `500` and `boom`; it fails today because the error carries only the status code.
- [ ] forward-envelope-bytes-test [auto]: `marshalEnvelope` (internal/engine/forward/envelope.go:19-29) has no direct test — `rg marshalEnvelope` finds exactly two callers (http.go:99, nats.go:45) and no test file — while envelope.go:6-8 makes the strong claim that "every Forwarder implementation in this package must produce byte-identical envelopes for the same inputs". Both suites check the shape only by unmarshalling into `bodyShape` (http_test.go:24-29), which cannot see key order or the exact bytes and would pass unchanged if a field were renamed on both sides. Add internal/engine/forward/envelope_test.go with a table pinning the exact output for (nil, nil) → `{"realm":"r","trigger":"t","action":null,"event":null}`, for the empty-non-nil pair the rule at 16-18 is about, and for a valued pair; container-free, no HTTP server needed.
EOF
tail -1 .mule/todo.md | cut -c1-120; echo "..."; rg -c "" .mule/todo.md
- [ ] forward-envelope-bytes-test [auto]: `marshalEnvelope` (internal/engine/forward/envelope.go:19-29) has no direct te
...
234

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,16 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Code review of `internal/engine/forward` (2026-09-26): four decisions on the bus-forwarding seam, none queued — all four are either a policy call or untestable-in-one-task.** Detail in `.mule/reviews/forward-2026-09-26.md`. The package itself is in good shape; the seam is thin and its tests are real.
+  **NATS async publish failures never reach Astrate's logger.** `forward.NATS.Forward` (internal/engine/forward/nats.go:49) calls `nats.Conn.Publish`, which appends to the client's buffered writer and returns (nats.go@v1.52.0 nats.go:4229 → 4424-4471) — the socket write happens later on the flusher goroutine, and a failure there goes to `Opts.AsyncErrorCB`. `NewNATS` (nats.go:36) passes no options, so that is nats.go's own `defaultErrHandler` printing to stderr (nats.go:1853-1856): no `*slog.Logger`, no realm, no trigger name — while `Executor.forward` has already counted the delivery `forwarded` (internal/engine/triggers/actions.go:503). `Forward` also takes `_ context.Context` (nats.go:44), so the 10s `RequestTimeout` the executor sets at actions.go:487 is not honoured on this path. **Not queued because** the fix is a `Logger` on `NATSConfig` plumbed through `cmd/astrate/newnats_nats.go` and `main.go:273` (three files), and the only test I can imagine kills the container mid-publish, which is racy — a task whose test cannot be made to fail reliably is a wasted run. Your call: plumb the logger, or switch to `PublishSync` and accept a round trip per delivery.
+  **The HTTP forwarder follows redirects, re-sending the realm header and the event body to a host the bus names.** Measured: a 307 to a second httptest server arrives there as `POST` with `Astarte-Realm="leaky-realm"` and the body intact (`http.DefaultClient`, http.go:63 — `main.go:262` passes no client). The webhook path behaves identically (actions.go:621), so a forwarder-only `CheckRedirect` would just make the two paths inconsistent. **Yours:** should Astrate pin `CheckRedirect` on both (refuse cross-host, or refuse all redirects), or is a redirecting bus a legitimate configuration?
+  **`triggers.forward.static_headers` gets no validation at all, while trigger *action* headers get a blocklist and a size cap.** `config.validate` checks kind, url, method and subject (internal/config/config.go:339-375) and nothing else; the action path runs `validateStaticHeaders` (internal/engine/triggers/validation.go:140-155) against upstream's hop-by-hop/sensitive list (validation.go:21-34) plus an 8 KiB total. Measured that the practical harm is smaller than it looks — net/http ignores a static `Host` and derives `Content-Length` from the body, so a blocked name is silently dropped rather than corrupting the request. **Yours:** should bus-forwarder headers carry the same policy as action headers, and if so is a boot-time rejection of a currently-"working" config acceptable?
+  **The NATS connection is never closed or flushed.** `main.go:116` builds the forwarder and `shutdown` (main.go:482-502) tears down the http listener, broker, flows and engine, with no forwarder step — so whatever is in the client's pending buffer at SIGTERM is lost, and the connection outlives every drain. The fix is a lifecycle method on `triggers.Forwarder` (internal/engine/triggers/actions.go:161), i.e. three files and an interface change: a design decision, so it is yours rather than a queue line.
+
+---
+
 - **github-issues triage run, 2026-09-26: nothing proposable (as every run since 2026-09-05), and the 19-issue alarm pile is a false alarm — it says the mule has landed nothing in 19 days and it has landed something every day.** 23 open issues, the same set plus **#112** (today, 09-26 11:04:46Z, "nothing has landed in 14h"). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction.
   **The "idle streak" that 19 consecutive entries here have repeated is wrong, and this run retracts it.** `.mule/log.md` — the script's own row-per-attempt log, starting exactly at 2026-09-04, the day the streak is claimed to begin — records attempts on 20 of the 22 days to 09-25 and `done` lands on 17 of them (3–8/day: 6, 10, 8, 6, 8, 7, 7, 6, 3, 5, 6, 5, 3, 0, 5, 3, 0, 5, 4, 7 for 09-04→09-25). The last three landed shas — `84a050e`, `ed1cd72`, `f965f39` — are all reachable from `origin/mule/queue` (`git branch -r --contains`), so the work landed *and* pushed, which is the exact failure the alarm exists to catch.
   **What the alarm really measures, and why it fires daily by construction.** `check_pulse` (mule.sh:739-774) compares `.heartbeat` against `MULE_IDLE_ALARM_HOURS` (default 8, mule.sh:740), and `beat` is called only on a `done` land or a `checked` run (mule.sh:578, 540) — a refill does not beat. Last beat: **2026-09-25T20:19:13Z**; the pulse that files the issue runs ~10:54–11:12Z. 20:19Z → 11:04Z is ~14.75h, comfortably over 8h, so the arithmetic is right and the conclusion is wrong: with lands bunched into an evening window, the gap to the next morning's tick is >8h *every day*. It is a mis-set threshold, not a dead mule. The two days in this window that landed nothing (09-19, 09-22) are the only ones where the alarm had anything real to say.
   **Two candidate fixes, both yours — how many hours of silence is a real problem is a judgement, so I did not queue either:** (a) raise `MULE_IDLE_ALARM_HOURS` to ~24 in `mule.service`; or (b) make the pulse read `.mule/log.md` instead of the heartbeat and alarm only when a *whole day* has no `done` row — more honest, since that file is the record of what actually landed, and it needs no threshold tuning as the cadence changes.

$ git status --short; echo "=== gate ==="; go vet ./internal/engine/forward/ && go test -count=1 ./internal/engine/forward/ 2>&1 | tail -2
M .mule/for-giulio.md
 M .mule/todo.md
?? .mule/reviews/forward-2026-09-26.md
=== gate ===
ok  	github.com/astrate-platform/astrate/internal/engine/forward	0.060s

Proposal run only — no code changed, no git touched.

```
Done: reviewed internal/engine/forward (207 non-test lines + both test files + its caller and wiring) and proposed 4 tasks, escalating 4 decisions
Files: .mule/reviews/forward-2026-09-26.md (new), .mule/todo.md (+4 lines), .mule/for-giulio.md (+4 decisions)
Verified: go vet ./internal/engine/forward/ && go test -count=1 ./internal/engine/forward/ -> pass (0.060s, container-free half; the nats build tag is off); three claims about net/http and nats.go Publish measured with throwaway programs in /tmp, not reasoned about
Unsure: whether static headers *should* win over the fixed Astarte-Realm header — the task asserts the documented behaviour and says to flip it if you disagree
Follow-ups: forward-static-header-validation, forward-static-headers-override-test, forward-status-error-body, forward-envelope-bytes-test
```
