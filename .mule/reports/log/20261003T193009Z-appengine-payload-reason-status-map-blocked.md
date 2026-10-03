slug: appengine-payload-reason-status-map
verdict: blocked
at:  f5aa263
ran: 2026-10-03T19:30:09Z on DietPi in 1200s


> build · big-pickle

I'll start by reading the two files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ rg -n "RejectReason|ReasonOf" --glob '!*_test.go' -l
pkg/payload/value.go
pkg/payload/payload.go
internal/engine/data.go
internal/engine/router.go
internal/appengine/http.go

$ ssh -o ConnectTimeout=5 legion 'echo ok' 2>&1 | head -5
ok

Legion Go answers. Let me read the code.
$ rg -n "writeError|RejectReason|ReasonOf" internal/appengine/http.go
92:		a.writeError(w, err)
134:		a.writeError(w, err)
146:		a.writeError(w, err)
155:		a.writeError(w, err)
173:		a.writeError(w, err)
182:		a.writeError(w, err)
192:		a.writeError(w, err)
221:		a.writeError(w, err)
242:		a.writeError(w, err)
257:		a.writeError(w, err)
262:		a.writeError(w, err)
306:		a.writeError(w, err)
316:		a.writeError(w, err)
335:		a.writeError(w, err)
354:		a.writeError(w, err)
364:		a.writeError(w, err)
374:		a.writeError(w, err)
385:		a.writeError(w, err)
420:		a.writeError(w, err)
431:		a.writeError(w, err)
444:		a.writeError(w, err)
494:		a.writeError(w, err)
502:		a.writeError(w, err)
513:		a.writeError(w, err)
534:		a.writeError(w, err)
553:// back as FieldErrors so writeError renders upstream's changeset-shaped 422
643:// writeError maps service/store errors onto upstream-shaped responses.
644:func (a *API) writeError(w http.ResponseWriter, err error) {
684:	case payload.ReasonOf(err) == payload.ReasonBadObject:

$ rg -n "RejectReason|ReasonOf|bad_object|Reason" pkg/payload/payload.go | head -60
17:// RejectReason keys the engine's per-reason metrics and device_error
40:		return FormatInvalid, rejectf(ReasonTooLarge, "%d byte payload exceeds %d byte cap", len(p), d.limit())
44:		return f, rejectf(ReasonUnknownFormat, "payload is neither empty, BSON, nor JSON")
71:			return DecodedPayload{}, rejectf(ReasonUnsetNotAllowed, "empty payload on a mapping without allow_unset")
117:		return DecodedPayload{}, rejectf(ReasonUnsetNotAllowed, "empty payload on an object-aggregated datastream")
148:			return DecodedPayload{}, rejectf(ReasonMissingRequired,
160:		return DecodedPayload{}, rejectf(ReasonBadTimestamp, "mapping declares explicit_timestamp but payload carries no t")

→ Read pkg/payload/payload.go [offset=1, limit=60]
→ Read internal/appengine/http.go [offset=640, limit=80]

$ rg -n "Reason[A-Z][A-Za-z]* +RejectReason|RejectReason =|ReasonOf|type RejectReason" pkg/payload/*.go
pkg/payload/value.go:122:type RejectReason uint8
pkg/payload/value.go:127:	ReasonNone RejectReason = iota
pkg/payload/value.go:227:// ReasonOf extracts the RejectReason from err, or ReasonNone if err is nil
pkg/payload/value.go:229:func ReasonOf(err error) RejectReason {
pkg/payload/value_test.go:38:			} else if ReasonOf(err) != ReasonTypeMismatch {
pkg/payload/value_test.go:39:				t.Errorf("doubleFromInt64(%d) reason = %v; want type_mismatch", tc.in, ReasonOf(err))
pkg/payload/value_test.go:101:	if _, err := int32FromInt64(math.MaxInt32 + 1); ReasonOf(err) != ReasonTypeMismatch {
pkg/payload/value_test.go:102:		t.Errorf("int32FromInt64(MaxInt32+1) reason = %v; want type_mismatch", ReasonOf(err))
pkg/payload/value_test.go:112:		if _, err := checkDouble(bad); ReasonOf(err) != ReasonTypeMismatch {
pkg/payload/value_test.go:113:			t.Errorf("checkDouble(%v) reason = %v; want type_mismatch", bad, ReasonOf(err))
pkg/payload/value_test.go:131:		if _, err := dateTimeFromMillis(bad); ReasonOf(err) != ReasonBadTimestamp {
pkg/payload/value_test.go:132:			t.Errorf("dateTimeFromMillis(%d) reason = %v; want bad_timestamp", bad, ReasonOf(err))
pkg/payload/value_test.go:136:	if _, err := checkDateTime(MaxDateTime.Add(time.Millisecond)); ReasonOf(err) != ReasonBadTimestamp {
pkg/payload/value_test.go:137:		t.Errorf("checkDateTime(max+1ms) reason = %v; want bad_timestamp", ReasonOf(err))
pkg/payload/value_test.go:139:	if _, err := checkDateTime(MinDateTime.Add(-time.Nanosecond)); ReasonOf(err) != ReasonBadTimestamp {
pkg/payload/value_test.go:140:		t.Errorf("checkDateTime(min-1ns) reason = %v; want bad_timestamp", ReasonOf(err))
pkg/payload/value_test.go:153:	if _, err := checkString(string([]byte{0xff, 0xfe})); ReasonOf(err) != ReasonTypeMismatch {
pkg/payload/value_test.go:154:		t.Errorf("invalid UTF-8 reason = %v; want type_mismatch", ReasonOf(err))
pkg/payload/value_test.go:160:	if _, err := checkString(atCap + "x"); ReasonOf(err) != ReasonValueTooLarge {
pkg/payload/value_test.go:161:		t.Errorf("string over cap reason = %v; want value_too_large", ReasonOf(err))
pkg/payload/value_test.go:167:	if _, err := collectArray[int32](MaxArrayLen+1, nil); ReasonOf(err) != ReasonValueTooLarge {
pkg/payload/value_test.go:168:		t.Errorf("oversize array reason = %v; want value_too_large", ReasonOf(err))
pkg/payload/value_test.go:213:	if ReasonOf(nil) != ReasonNone || ReasonOf(errors.New("plain")) != ReasonNone {
pkg/payload/value_test.go:214:		t.Error("ReasonOf must return ReasonNone for nil and untyped errors")
pkg/payload/required_test.go:38:	if _, err := DecodeObject([]byte(`{"v":{"lon":9.0}}`), leaves); ReasonOf(err) != ReasonMissingRequired {
pkg/payload/required_test.go:39:		t.Fatalf("JSON rejection reason = %v (err %v); want %v", ReasonOf(err), err, ReasonMissingRequired)
pkg/payload/required_test.go:47:	if _, err := DecodeObject(bsonBad, leaves); ReasonOf(err) != ReasonMissingRequired {
pkg/payload/required_test.go:48:		t.Fatalf("BSON rejection reason = %v (err %v); want %v", ReasonOf(err), err, ReasonMissingRequired)
pkg/payload/fuzz_test.go:84:				if ReasonOf(err) == ReasonNone {
pkg/payload/fuzz_test.go:92:			if ReasonOf(err) == ReasonNone {
pkg/payload/payload_test.go:375:				if ReasonOf(err) != tc.wantReason {
pkg/payload/payload_test.go:376:					t.Fatalf("Decode(%s) err = %v (reason %v); want reason %v", tc.in, err, ReasonOf(err), tc.wantReason)
pkg/payload/payload_test.go:421:		if _, err := DecodeObject([]byte(in), leaves); ReasonOf(err) != wantReason {
pkg/payload/payload_test.go:422:			t.Errorf("DecodeObject(%s) reason = %v; want %v", in, ReasonOf(err), wantReason)
pkg/payload/payload_test.go:427:	if _, err := DecodeObject([]byte(`{"v":{"lat":1.0}}`), objectLeaves(true)); ReasonOf(err) != ReasonBadTimestamp {
pkg/payload/payload_test.go:428:		t.Errorf("explicit object without t: reason = %v; want bad_timestamp", ReasonOf(err))
pkg/payload/payload_test.go:535:				if ReasonOf(err) != tc.wantReason {
pkg/payload/payload_test.go:536:					t.Fatalf("reason = %v (err %v); want %v", ReasonOf(err), err, tc.wantReason)
pkg/payload/payload_test.go:554:	if _, err := DecodeObject(dup, objectLeaves(false)); ReasonOf(err) != ReasonBadObject {
pkg/payload/payload_test.go:555:		t.Errorf("duplicate object key reason = %v; want bad_object", ReasonOf(err))
pkg/payload/payload_test.go:597:			if ReasonOf(err) != ReasonBadObject {
pkg/payload/payload_test.go:598:				t.Fatalf("reason = %v (err %v); want bad_object", ReasonOf(err), err)
pkg/payload/payload_test.go:630:	if _, err := Decode(bigBlob, mapping(interfaceschema.BinaryBlob, false)); ReasonOf(err) != ReasonTooLarge {
pkg/payload/payload_test.go:631:		t.Errorf("BSON over default cap: reason = %v; want too_large", ReasonOf(err))
pkg/payload/payload_test.go:635:	if _, err := Decode(bigJSON, mapping(interfaceschema.String, false)); ReasonOf(err) != ReasonTooLarge {
pkg/payload/payload_test.go:636:		t.Errorf("JSON over default cap: reason = %v; want too_large", ReasonOf(err))
pkg/payload/payload_test.go:640:	if _, err := small.Individual([]byte(`{"v":true}`), mapping(interfaceschema.Boolean, false)); ReasonOf(err) != ReasonTooLarge {
pkg/payload/payload_test.go:644:	if _, err := small.Individual(okDoc, mapping(interfaceschema.Boolean, false)); ReasonOf(err) != ReasonTooLarge {
pkg/payload/payload_test.go:650:	if _, err := wide.Individual(overString, mapping(interfaceschema.String, false)); ReasonOf(err) != ReasonValueTooLarge {
pkg/payload/payload_test.go:651:		t.Errorf("over-cap string: reason = %v; want value_too_large", ReasonOf(err))
pkg/payload/payload_test.go:656:	if _, err := Decode(overJSON, mapping(interfaceschema.IntegerArray, false)); ReasonOf(err) != ReasonValueTooLarge {
pkg/payload/payload_test.go:657:		t.Errorf("JSON 1025-element array: reason = %v; want value_too_large", ReasonOf(err))
pkg/payload/payload_test.go:674:	if _, err := Decode(overBSON, mapping(interfaceschema.IntegerArray, false)); ReasonOf(err) != ReasonValueTooLarge {
pkg/payload/payload_test.go:675:		t.Errorf("BSON 1025-element array: reason = %v; want value_too_large", ReasonOf(err))
pkg/payload/payload_test.go:691:	if _, err := Decode([]byte{}, mapping(interfaceschema.Double, false)); ReasonOf(err) != ReasonUnsetNotAllowed {
pkg/payload/payload_test.go:692:		t.Errorf("unset without allow_unset: reason = %v; want unset_not_allowed", ReasonOf(err))
pkg/payload/payload_test.go:694:	if _, err := DecodeObject(nil, objectLeaves(false)); ReasonOf(err) != ReasonUnsetNotAllowed {
pkg/payload/payload_test.go:695:		t.Errorf("unset object: reason = %v; want unset_not_allowed", ReasonOf(err))
pkg/payload/payload_test.go:697:	if _, err := Decode([]byte{0xDE, 0xAD, 0xBE, 0xEF, 0x99}, mapping(interfaceschema.Double, false)); ReasonOf(err) != ReasonUnknownFormat {
pkg/payload/payload_test.go:702:	if _, err := Decode([]byte(`{"v":1}`), nil); err == nil || ReasonOf(err) != ReasonNone {
pkg/payload/payload_test.go:705:	if _, err := DecodeObject([]byte(`{"v":{}}`), nil); err == nil || ReasonOf(err) != ReasonNone {
pkg/payload/payload_test.go:827:		if _, err := Encode(math.NaN(), nil, f); ReasonOf(err) != ReasonTypeMismatch {
pkg/payload/payload_test.go:828:			t.Errorf("%v: NaN reason = %v; want type_mismatch", f, ReasonOf(err))
pkg/payload/payload_test.go:830:		if _, err := Encode(math.Inf(-1), nil, f); ReasonOf(err) != ReasonTypeMismatch {
pkg/payload/payload_test.go:831:			t.Errorf("%v: -Inf reason = %v; want type_mismatch", f, ReasonOf(err))
pkg/payload/payload_test.go:833:		if _, err := Encode(uint(5), nil, f); ReasonOf(err) != ReasonTypeMismatch {
pkg/payload/payload_test.go:834:			t.Errorf("%v: unsupported Go type reason = %v; want type_mismatch", f, ReasonOf(err))
pkg/payload/payload_test.go:836:		if _, err := Encode(strings.Repeat("a", MaxStringLen+1), nil, f); ReasonOf(err) != ReasonValueTooLarge {
pkg/payload/payload_test.go:837:			t.Errorf("%v: over-cap string reason = %v; want value_too_large", f, ReasonOf(err))
pkg/payload/payload_test.go:839:		if _, err := Encode(make([]bool, MaxArrayLen+1), nil, f); ReasonOf(err) != ReasonValueTooLarge {
pkg/payload/payload_test.go:840:			t.Errorf("%v: over-cap array reason = %v; want value_too_large", f, ReasonOf(err))
pkg/payload/payload_test.go:842:		if _, err := Encode(time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC), nil, f); ReasonOf(err) != ReasonBadTimestamp {
pkg/payload/payload_test.go:843:			t.Errorf("%v: out-of-window datetime reason = %v; want bad_timestamp", f, ReasonOf(err))
pkg/payload/payload_test.go:846:		if _, err := Encode(22.5, &badTS, f); ReasonOf(err) != ReasonBadTimestamp {
pkg/payload/payload_test.go:847:			t.Errorf("%v: out-of-window t reason = %v; want bad_timestamp", f, ReasonOf(err))
pkg/payload/payload_test.go:849:		if _, err := Encode(map[string]Value{}, nil, f); ReasonOf(err) != ReasonBadObject {
pkg/payload/payload_test.go:850:			t.Errorf("%v: empty object reason = %v; want bad_object", f, ReasonOf(err))
pkg/payload/payload_test.go:852:		if _, err := Encode(map[string]Value{"a": map[string]Value{"b": 1.0}}, nil, f); ReasonOf(err) != ReasonBadObject {
pkg/payload/payload_test.go:853:			t.Errorf("%v: nested object reason = %v; want bad_object", f, ReasonOf(err))

→ Read pkg/payload/value.go

# Todos
[•] Read existing writeerror_test.go + engine/serverdata.go:134 too_large caller
[ ] Probe upstream b6d46ad4 status+detail for each of 9 reasons on Legion Go
[ ] Extend writeError into a full switch with one status+detail per reason
[ ] Add table-driven cases in writeerror_test.go per reason
[ ] Check too_large REST reachability (Decoder.MaxSize fed from len(envelope))
[ ] Run gate: go vet ./... && go test ./... && gofmt -l .

→ Read internal/appengine/writeerror_test.go
→ Read internal/engine/serverdata.go [limit=60, offset=110]

$ rg -n "payload.Decoder\{|Decoder\{" --glob '!*_test.go' internal/ pkg/ cmd/ 2>/dev/null
internal/engine/control.go:235:	dp, err := payload.Decoder{MaxSize: len(env)}.Individual(env, &interfaceschema.CompiledMapping{
pkg/payload/payload.go:182:	return Decoder{}.Individual(p, m)
pkg/payload/payload.go:188:	return Decoder{}.Object(p, leaves)
internal/engine/router.go:175:		dec:     payload.Decoder{MaxSize: cfg.MaxPayloadBytes},
internal/engine/serverdata.go:134:	dec := payload.Decoder{MaxSize: len(envelope)}

→ Read .mule/recipes/legion-go.md

$ ssh legion 'ls ~/ ; echo ---; ls ~/astarte* 2>/dev/null | head -20'
drwxr-xr-x    - atsetilam 26 ago 23:04  .agents
drwxr-x---    - atsetilam  3 ott 17:50  .android
drwxr-xr-x    - atsetilam  5 ago 16:04  .bun
drwx------    - atsetilam  3 ott 18:46  .cache
drwxr-xr-x    - atsetilam  4 ago 03:59  .cargo
drwxr-xr-x    - atsetilam  3 ott 18:31  .claude
drwxr-xr-x    - atsetilam  3 ott 18:46  .config
drwxr-xr-x    - atsetilam 27 set 16:45  .docker
drwx------    - atsetilam 11 lug 16:55  .gnupg
drwxr-xr-x    - atsetilam  6 ago 19:19  .local
drwxr-xr-x    - atsetilam  5 ago 16:04  .npm
drwxr-xr-x    - atsetilam 11 lug 20:14  .opencode
drwx------    - atsetilam 18 mar 01:06  .pki
drwxr-xr-x    - atsetilam  3 ago 17:17  .rustup
drwx------    - atsetilam 24 ago 13:58  .sigstore
drwx------    - atsetilam  3 ott 18:47 󰢬 .ssh
drwxr-xr-x    - atsetilam  3 ott 16:45  .steam
drwxr-xr-x    - atsetilam 19 mar 20:20  .supermodel
drwxr-xr-x    - atsetilam  5 ago 22:42  .terminfo
drwxr-xr-x    - atsetilam  5 ago 16:06  .triton
drwxr-xr-x    - atsetilam 26 set 01:32  .unsloth
drwxr-xr-x    - atsetilam 18 mar 01:11  .var
drwxr-xr-x    - atsetilam 11 lug 20:32  .vcpkg
drwx------    - atsetilam  3 ott 00:22  .w3m
drwxr-xr-x    - atsetilam  3 ott 16:57  agent-skills
drwxr-xr-x    - atsetilam 25 lug 17:01  Applications
drwxr-xr-x    - atsetilam 22 ago 03:35  astrate
drwxr-xr-x    - atsetilam 11 lug 22:50 󱧼 build
drwxr-xr-x    - atsetilam 27 set 16:36  Desktop
drwxr-xr-x    - atsetilam  3 ott 17:50  dev
drwxr-xr-x    - atsetilam 18 mar 01:00  Documenti
drwxr-xr-x    - atsetilam  3 ott 16:44  dot-steam.bak.1791038671
drwxr-xr-x    - atsetilam 22 ago 03:56  edgehog-test
drwxr-xr-x    - atsetilam 19 mar 19:46  emudeck
drwxr-xr-x    - atsetilam 19 mar 20:26  Emulation
drwxr-xr-x    - atsetilam 19 mar 20:33  ES-DE
drwxr-xr-x    - atsetilam  3 ott 16:04  g-mind
drwxr-xr-x    - atsetilam 26 set 00:29  Games
drwxr-xr-x    - atsetilam 27 lug 06:05  go
drwxr-xr-x    - atsetilam 21 lug 11:24  Immagini
drwxr-xr-x    - atsetilam 31 ago 00:44  legion-gmind-backup-20260831-004404
drwxr-xr-x    - atsetilam  4 ago 03:59  listello
drwxr-xr-x    - atsetilam  3 ago 18:09  listello-wip
drwxr-xr-x    - atsetilam  6 ago 02:57  llama.cpp
drwxr-xr-x    - atsetilam 18 mar 01:00  Modelli
drwxr-xr-x    - atsetilam 12 lug 07:16  models
drwxr-xr-x    - atsetilam 18 mar 01:00  Musica
drwxr-xr-x    - atsetilam 27 ago 03:38  orange-cat
drwxr-xr-x    - atsetilam 11 lug 18:30  Progetti
drwxr-xr-x    - atsetilam 21 lug 11:19  prova-sito-phi
drwxr-xr-x    - atsetilam 18 mar 01:00  Pubblici
drwxr-xr-x    - atsetilam 27 set 17:31  Scaricati
drwxr-xr-x    - atsetilam 20 lug 18:10  steamvr
drwxr-xr-x    - atsetilam  6 ago 02:42  test
drwxr-xr-x    - atsetilam  3 ott 16:16  toshy-src
drwxr-xr-x    - atsetilam  5 ago 17:13  unsloth-studio-test
drwxr-xr-x    - atsetilam  5 ago 16:06  unsloth_compiled_cache
drwxr-xr-x    - atsetilam 18 mar 01:00  Video
drwxr-xr-x    - atsetilam 30 mar 20:27  whisper.cpp
.rw-r--r--   21 atsetilam 18 mar 00:48 󱆃 .bash_logout
.rw-r--r--   57 atsetilam 18 mar 00:48 󱆃 .bash_profile
.rw-r--r--  172 atsetilam 18 mar 00:48 󱆃 .bashrc
.rw-------  74k atsetilam  3 ott 20:54  .claude.json
.rw-r--r--  133 atsetilam 21 lug 02:35 󰊢 .gitconfig
.rw-r--r--  348 atsetilam  3 ott 18:46  .gtkrc-2.0
.rw-------   30 atsetilam  6 ago 19:19  .npmrc
.rw-------  256 atsetilam 18 mar 01:00  .pulse-cookie
lrwxrwxrwx    - atsetilam  3 ott 16:45  .steampath -> /home/atsetilam/.steam/sdk32/steam
lrwxrwxrwx    - atsetilam  3 ott 16:45  .steampid -> /home/atsetilam/.steam/steam.pid
.rw-------  885 atsetilam 19 mar 20:11  .viminfo
.rw-r--r--  180 atsetilam 19 mar 20:20  .wget-hsts
.rw-r--r--  57k atsetilam  3 ott 16:21  .zcompdump-atsetilam-legion-5.9.2
.r--r--r-- 122k atsetilam  3 ott 16:21  .zcompdump-atsetilam-legion-5.9.2.zwc
.rw-r--r--   56 atsetilam 18 mar 00:48 󱆃 .zshrc
.rw-r--r--   64 atsetilam 20 lug 18:12  crash_log.txt
.rw-r--r--  22k atsetilam 30 mar 17:46  install-docker.sh
.rw-r--r--  19k atsetilam 21 lug 01:49  scale-gfx1103-llamacpp-report.md
.rw-r--r-- 1,4G atsetilam 18 lug 10:43  scale-latest-amd64.tar.xz
---
fish: No matches for wildcard '~/astarte*'. See `help language#wildcards-globbing`.
ls ~/ ; echo ---; ls ~/astarte* 2>/dev/null | head -20
                     ^~~~~~~~~^

$ ssh legion 'docker ps --format "{{.Names}}\t{{.Image}}\t{{.Ports}}" 2>&1' | head -30
failed to connect to the docker API at unix:///home/atsetilam/.docker/desktop/docker.sock; check if the path is correct and if the daemon is running: dial unix /home/atsetilam/.docker/desktop/docker.sock: connect: no such file or directory

$ ssh legion 'systemctl --user is-active docker 2>&1; echo "---"; docker ps 2>&1 | head -5' | tail -15
inactive
---
failed to connect to the docker API at unix:///home/atsetilam/.docker/desktop/docker.sock; check if the path is correct and if the daemon is running: dial unix /home/atsetilam/.docker/desktop/docker.sock: connect: no such file or directory

$ ssh legion 'ls ~/astrate 2>&1 | head -15' | tail -10
drwxr-xr-x    - atsetilam 26 lug 23:38  cmd
drwxr-xr-x    - atsetilam 26 lug 23:38  deploy
drwxr-xr-x    - atsetilam 22 ago 03:35  docs
drwxr-xr-x    - atsetilam 22 ago 03:35  examples
drwxr-xr-x    - atsetilam 22 ago 03:35  internal
drwxr-xr-x    - atsetilam 22 ago 03:35  migrations
drwxr-xr-x    - atsetilam 26 lug 23:38  pkg
drwxr-xr-x    - atsetilam 26 lug 23:38  test
drwxr-xr-x    - atsetilam 22 ago 03:35  tools
.rw-r--r--  207 atsetilam 22 ago 03:35  .dockerignore

$ ssh legion 'cd ~/astrate && go version && git rev-parse --short HEAD && git log --oneline -1' | head -10
go version go1.26.5 linux/amd64
fb3f335
fb3f335 feat(flow): block-death detection + auto-restart with backoff (#45 phase 1)

$ ssh legion 'cd ~/astrate && go test -v ./pkg/payload/... ./internal/appengine/... 2>&1 | tail -20' | head -15
--- PASS: TestWatchAuthPath/device_trigger_naming_a_different_device_is_refused (0.00s)
    --- PASS: TestWatchAuthPath/device_trigger_with_the_wildcard_"*"_is_refused (0.00s)
    --- PASS: TestWatchAuthPath/group_data_trigger (0.00s)
    --- PASS: TestWatchAuthPath/group_device_trigger (0.00s)
=== RUN   TestWireSession_WatchPathIsMatchPathScoped
=== RUN   TestWireSession_WatchPathIsMatchPathScoped/claim_without_the_match_path_is_refused
=== RUN   TestWireSession_WatchPathIsMatchPathScoped/claim_with_the_match_path_is_accepted
--- PASS: TestWireSession_WatchPathIsMatchPathScoped (0.12s)
    --- PASS: TestWireSession_WatchPathIsMatchPathScoped/claim_without_the_match_path_is_refused (0.02s)
    --- PASS: TestWireSession_WatchPathIsMatchPathScoped/claim_with_the_match_path_is_accepted (0.09s)
PASS
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	1.129s
=== RUN   TestSocketSSE
=== RUN   TestSocketSSE/Unauthorized
=== RUN   TestSocketSSE/StreamsEventAsSSE

$ ssh legion 'find ~ -name "*astarte*" -not -path "*/cache/*" 2>&1 | head -20' | tail -15
/home/atsetilam/astrate/bench/.astarte/.github/workflows/astarte-end-to-end-test-workflow.yaml
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/appengine_astarte_overview.svg
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api.ex
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/device/astarte_value.ex
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web.ex
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/priv/static/astarte_appengine_api.yaml
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/test/astarte_appengine_api
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/test/astarte_appengine_api/device/astarte_value_test.exs
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/test/astarte_appengine_api_web
/home/atsetilam/astrate/bench/.astarte/apps/astarte_data_updater_plant
/home/atsetilam/astrate/bench/.astarte/apps/astarte_data_updater_plant/data_updater_astarte_overview.svg
/home/atsetilam/astrate/bench/.astarte/apps/astarte_data_updater_plant/lib/astarte_data_updater_plant

$ ssh legion 'find ~/astrate/bench/.astarte -name "*.ex" -path "*payload*controller*error*" 2>&1 | head -10' | tail -5
(no output)

$ ssh legion 'grep -rn "payload" ~/astrate/bench/.astarte/apps/astarte_appengine_api 2>&1 | head -40' | tail -30
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/data_transmitter.ex:83:    exchanged_bytes = byte_size(bson_payload) + byte_size(interface) + byte_size(path)
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/data_transmitter.ex:91:    VMQPlugin.publish(topic, bson_payload, @property_qos)
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/data_transmitter.ex:114:  defp make_payload_map(payload, nil, nil) do
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/data_transmitter.ex:115:    %{v: payload}
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/data_transmitter.ex:118:  defp make_payload_map(payload, timestamp, nil) do
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/data_transmitter.ex:119:    %{v: payload, t: timestamp}
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/data_transmitter.ex:122:  defp make_payload_map(payload, nil, metadata) do
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/data_transmitter.ex:123:    %{v: payload, m: metadata}
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/data_transmitter.ex:126:  defp make_payload_map(payload, timestamp, metadata) do
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/data_transmitter.ex:127:    %{v: payload, t: timestamp, m: metadata}
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/rooms/amqp_client.ex:82:  def handle_info({:basic_deliver, payload, meta}, chan) do
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/rooms/amqp_client.ex:83:    _ = Logger.debug("Got event, payload: #{inspect(payload)} meta: #{inspect(meta)}.")
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/rooms/amqp_client.ex:85:    EventsDispatcher.dispatch(payload)
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/rooms/room.ex:179:        payload = %{
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/rooms/room.ex:191:        Endpoint.broadcast("rooms:" <> room_name, "new_event", payload)
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/vmq_plugin.ex:43:  def publish(topic, payload, qos)
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/vmq_plugin.ex:44:      when is_binary(topic) and is_binary(payload) and is_integer(qos) and qos >= 0 and qos <= 2 do
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/vmq_plugin.ex:50:        payload: payload,
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/channels/rooms_channel.ex:32:  def join("rooms:" <> room_name, _payload, socket) do
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/channels/rooms_channel.ex:50:  def handle_in("watch", payload, socket) do
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/channels/rooms_channel.ex:51:    changeset = WatchRequest.changeset(%WatchRequest{}, payload)
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/channels/rooms_channel.ex:57:      payload = WatchRequestView.render("watch_request.json", %{watch_request: request})
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/channels/rooms_channel.ex:58:      broadcast(socket, "watch_added", payload)
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/channels/rooms_channel.ex:79:  def handle_in("unwatch", payload, socket) do
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/channels/rooms_channel.ex:80:    changeset = UnwatchRequest.changeset(%UnwatchRequest{}, payload)
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/channels/user_socket.ex:42:  def connect(%{"realm" => realm} = payload, socket) do
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/channels/user_socket.ex:45:    with token <- Map.get(payload, "token"),
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/priv/static/astarte_appengine_api.yaml:1118:          description: Total ammount of received payload bytes.
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/test/astarte_appengine_api/data_transmitter_test.exs:43:  @payload "importantdata"
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/test/astarte_appengine_api/data_transmitter_test.exs:57:      encoded_payload = Cyanide.encode!(%{v: @payload})

$ ssh legion 'grep -rn "RejectReason\|missing_required\|too_large\|malformed\|no_value\|type_mismatch\|value_too_large\|unset_not_allowed\|unknown_format\|bad_timestamp" ~/astrate/bench/.astarte/apps/astarte_appengine_api/lib 2>&1 | head -60' | tail -40
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/device/queries.ex:370:          tag: "unset_not_allowed"

$ ssh legion 'grep -rn "missing_required\|too_large.*400\|status.*too_large" ~/astrate/bench/.astarte 2>&1 | head -30' | tail -25
grep: /home/atsetilam/astrate/bench/.astarte/compose/cfssl-config/ca-key.pem.old: Permesso negato
grep: /home/atsetilam/astrate/bench/.astarte/compose/vernemq-certs/ca.old/private/ca.key.pem: Permesso negato

$ ssh legion 'grep -rn "too_large" ~/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/device 2>&1 | head -15' | tail -10
(no output)

$ ssh legion 'grep -rn "missing_required\|bad_timestamp\|type_mismatch\|value_too_large\|unset_not_allowed" ~/astrate/bench/.astarte/apps/astarte_appengine_api 2>&1 | head -40' | tail -20
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/device/queries.ex:370:          tag: "unset_not_allowed"

$ ssh legion 'grep -rn "payload.*reject\|data.*write.*error" ~/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers 2>&1 | head -30' | tail -25
(no output)

$ ssh legion 'find ~/astrate/bench/.astarte -name "*.ex" | grep controller | head -15' | tail -10
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/interface_values_by_device_alias_controller.ex
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/interface_values_by_group_controller.ex
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/interface_values_controller.ex
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/stats_controller.ex
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/version_controller.ex
/home/atsetilam/astrate/bench/.astarte/apps/astarte_housekeeping_api/lib/astarte_housekeeping_api_web/controllers/fallback_controller.ex
/home/atsetilam/astrate/bench/.astarte/apps/astarte_housekeeping_api/lib/astarte_housekeeping_api_web/controllers/realm_controller.ex
/home/atsetilam/astrate/bench/.astarte/apps/astarte_housekeeping_api/lib/astarte_housekeeping_api_web/controllers/version_controller.ex
/home/atsetilam/astrate/bench/.astarte/apps/astarte_pairing_api/lib/astarte_pairing_api_web/controllers/agent_controller.ex
/home/atsetilam/astrate/bench/.astarte/apps/astarte_pairing_api/lib/astarte_pairing_api_web/controllers/device_controller.ex

$ ssh legion 'grep -rn "bad_object\|Unexpected object key\|missing_required" ~/astrate/bench/.astarte/apps/astarte_appengine_api 2>&1 | head -20' | tail -15
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api/device/queries.ex:531:            "Unexpected object key #{inspect(obj_key)} with value #{inspect(obj_value)}."
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex:98:    %{errors: %{detail: "Unexpected object key"}}
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_values_by_device_alias_controller_test.exs:146:      assert json_response(conn, 400)["errors"] == %{"detail" => "Unexpected object key"}
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_values_by_group_controller_test.exs:235:      assert json_response(conn, 400)["errors"] == %{"detail" => "Unexpected object key"}
/home/atsetilam/astrate/bench/.astarte/apps/astarte_appengine_api/test/astarte_appengine_api_web/controllers/interface_values_controller_test.exs:179:      assert json_response(conn, 400)["errors"] == %{"detail" => "Unexpected object key"}

$ ssh legion 'cat ~/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex 2>&1 | head -120' | tail -80
def render("404.json", _assigns) do
    %{errors: %{detail: "Not found"}}
  end

  def render("404_device.json", _assigns) do
    %{errors: %{detail: "Device not found"}}
  end

  def render("404_endpoint_not_found.json", _assigns) do
    %{errors: %{detail: "Endpoint not found"}}
  end

  def render("404_interface_not_found.json", _assigns) do
    %{errors: %{detail: "Interface not found"}}
  end

  def render("404_interface_not_in_introspection.json", _assigns) do
    %{errors: %{detail: "Interface not found in device introspection"}}
  end

  def render("404_path.json", _assigns) do
    %{errors: %{detail: "Path not found"}}
  end

  def render("404_group.json", _assigns) do
    %{errors: %{detail: "Group not found"}}
  end

  def render("422_attribute_key_not_found.json", _assigns) do
    %{errors: %{detail: "Attribute key not found"}}
  end

  def render("409_group_already_exists.json", _assigns) do
    %{errors: %{detail: "Group already exists"}}
  end

  def render("409_device_already_in_group.json", _assigns) do
    %{errors: %{detail: "Device already in group"}}
  end

  def render("409_alias_already_in_use.json", _assigns) do
    %{errors: %{detail: "Alias already in use"}}
  end

  def render("422_alias_tag_not_found.json", _assigns) do
    %{errors: %{detail: "Alias tag not found"}}
  end

  def render("422_invalid_alias.json", _assigns) do
    %{errors: %{detail: "Invalid alias"}}
  end

  def render("422_invalid_attributes.json", _assigns) do
    %{errors: %{detail: "Invalid attributes"}}
  end

  def render("422_unexpected_object_key.json", _assigns) do
    %{errors: %{detail: "Unexpected object key"}}
  end

  def render("500.json", _assigns) do
    %{errors: %{detail: "Internal server error"}}
  end

  def render("401.json", _assigns) do
    %{errors: %{detail: "Unauthorized"}}
  end

  def render("403.json", _assigns) do
    %{errors: %{detail: "Forbidden"}}
  end

  def render("503_cannot_push_to_device.json", _assigns) do
    %{errors: %{detail: "Cannot push to device"}}
  end

  def render("503_service_unavailable.json", _assigns) do
    %{errors: %{detail: "Service unavailable"}}
  end

$ ssh legion 'grep -rn "missing_required\|too_large\|malformed\|no_value\|type_mismatch\|value_too_large\|unset_not_allowed\|bad_timestamp\|unknown_format" ~/astrate/bench/.astarte/apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex' | tail -40
(no output)

