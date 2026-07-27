// Command recordchannels captures upstream Astarte's Channels (Phoenix
// websocket) frames and the error_name values it puts on device_error events,
// and writes them to test/conformance/upstream as a fixture plus a verbatim
// transcript (M12 plan, phase 04b).
//
// It is the second half of the oracle. Phase 04a recorded upstream's REST
// error envelopes, which needed nothing but HTTP and a housekeeping key; this
// half needs a real device driven over MQTT, because the only way to learn
// what upstream calls a rejection is to make upstream reject something.
//
// The prize is M11 phase 08's device_error -> error_name mapping table, which
// was an invention: an architect read an enum out of the Dashboard's JavaScript
// bundle and chose, by judgement, which of Astrate's own reject reasons mapped
// to which upstream name. Every value in it is accepted by the Dashboard, and
// not one of them had been checked against what upstream actually emits for the
// same rejection. This command checks.
//
// Two instruments, deliberately kept apart, because they see different things:
//
//   - the Channels room, which is what a *client* sees. This is the surface
//     Astrate has to match, so it is the one the fixture is built from.
//   - upstream's own reject behaviour, inferred from whether an event arrives
//     at all. Several provocations are detected by upstream and never reach the
//     room; that absence is recorded as an observation rather than hidden,
//     because "upstream stays silent here" is exactly the kind of claim this
//     directory exists to stop people guessing about.
//
// Each provocation gets a fresh websocket *and* a fresh device session. That is
// not caution for its own sake: several of these errors make upstream force a
// clean session and disconnect the device, so a run that reuses one session
// records the first result honestly and every later one as silence.
//
// Usage (needs a reachable upstream and a provisioned realm):
//
//	ASTARTE_UPSTREAM_URL=http://api.astarte.localhost:8080 \
//	ASTARTE_UPSTREAM_STATE=bench/upstream.json \
//	    go run ./recordchannels
//
// The state file is what `bench provision` writes: it carries the realm signing
// key (needed to mint the a_ch token) and the device credentials secrets. It is
// gitignored and nothing from it is written to the fixture.
//
// Nothing here imports Astrate: the recorder must not be able to describe
// upstream in Astrate's own terms. BSON is hand-encoded below rather than
// pulled in as a dependency, because half these provocations are *malformed*
// BSON and a library exists precisely to stop you emitting that.
package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coder/websocket"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang-jwt/jwt/v5"
)

// ---------------------------------------------------------------- state file

// state is the subset of `bench provision`'s output this recorder needs.
type state struct {
	BaseURL   string `json:"base_url"`
	Endpoints struct {
		Pairing string `json:"pairing"`
	} `json:"endpoints"`
	Realm     string `json:"realm"`
	RealmKey  string `json:"realm_key_pem"`
	BrokerURL string `json:"broker_url"`
	Devices   []struct {
		ID     string `json:"id"`
		Secret string `json:"secret"`
	} `json:"devices"`
}

// The interfaces `bench provision` installs. Both are device-owned datastreams,
// which bounds what this recorder can provoke: with no server-owned interface
// in the realm, write_on_server_owned_interface cannot be reached, and that is
// recorded as a gap rather than guessed at.
const (
	ifaceIndividual = "org.astrate.bench.Individual" // /value, double
	ifaceObject     = "org.astrate.bench.Object"     // /data/{temp,hum,status}
	introspection   = ifaceIndividual + ":1:0;" + ifaceObject + ":1:0"

	// otherDeviceID is a well-formed device id that is deliberately NOT the
	// device this recorder drives. It exists to pair a narrow WATCH claim's
	// acceptance with a refusal that differs only in which device is named.
	otherDeviceID = "aaaaaaaaaaaaaaaaaaaaaa"
)

// --------------------------------------------------------------- the fixture

// observation is one provocation and everything upstream said about it.
type observation struct {
	Name string `json:"name"`
	Why  string `json:"why"`

	// What the device did.
	Topic   string `json:"topic"`   // relative to <realm>/<device_id>
	Payload string `json:"payload"` // base64 of the exact bytes published

	// What the room saw for *this* publish. An event counts only when its
	// metadata echoes back the exact bytes we published: the device session
	// itself emits device_session_not_found events (upstream answers the
	// emptyCache control that way until Data Updater Plant has registered the
	// session), and taking the first event to arrive attributes that noise to
	// whatever was being provoked. That bug made two consecutive runs of this
	// recorder disagree, which is how it was found.
	// Delivery of these events is not reliable on this stack, so each row is
	// provoked Attempts times. ErrorName is the name upstream gave whenever it
	// gave one; Delivered says how often an event actually arrived. A row with
	// Delivered == 0 means "never observed in Attempts tries" — which is a
	// weaker claim than "upstream emits nothing", and is deliberately not
	// written as the stronger one. See the README on why delivery is flaky.
	ErrorName string   `json:"error_name"`
	Metadata  []string `json:"metadata_keys,omitempty"`
	Attempts  int      `json:"attempts"`
	Delivered int      `json:"delivered"`
	NoEvent   bool     `json:"never_reached_the_room,omitempty"`

	// Events that arrived in the window but belong to the session rather than
	// the provocation. Recorded, not discarded, so the correlation above can be
	// audited rather than trusted.
	UncorrelatedFrames []string `json:"uncorrelated_frames,omitempty"`

	Frames []string `json:"frames"` // verbatim, in arrival order
}

// authObservation records one a_ch claim tried against one operation. These
// settle M11's decision 5, which no test could reach before: a blanket
// `.*::.*` grant authorizes every reading of the WATCH path equally, so only a
// real upstream can say whether upstream reads it that way too.
type authObservation struct {
	Name      string   `json:"name"`
	Why       string   `json:"why"`
	Claims    []string `json:"a_ch_claims"`
	Operation string   `json:"operation"`
	Payload   string   `json:"payload"`
	Reply     string   `json:"reply"` // the phx_reply frame, verbatim
	Accepted  bool     `json:"accepted"`
}

type fixture struct {
	RecordedAt      string `json:"recorded_at"`
	AstarteVersion  string `json:"astarte_version"`
	Realm           string `json:"realm"`
	RecorderCommand string `json:"recorder_command"`

	SocketPath    string            `json:"socket_path"`
	Authorization []authObservation `json:"authorization"`
	DeviceErrors  []observation     `json:"device_errors"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "recordchannels:", err)
		os.Exit(1)
	}
}

var (
	st      state
	baseURL string
	tr      = &strings.Builder{} // the transcript, written as we go
)

func logf(format string, a ...any) {
	fmt.Fprintf(tr, format+"\n", a...)
	fmt.Printf(format+"\n", a...)
}

func run() error {
	baseURL = env("ASTARTE_UPSTREAM_URL", "http://api.astarte.localhost:8080")
	statePath := env("ASTARTE_UPSTREAM_STATE", "bench/upstream.json")

	b, err := os.ReadFile(statePath)
	if err != nil {
		return fmt.Errorf("reading state (%s): %w — run `go run . provision` in bench/ first", statePath, err)
	}
	if err := json.Unmarshal(b, &st); err != nil {
		return fmt.Errorf("parsing state: %w", err)
	}
	if len(st.Devices) == 0 {
		return fmt.Errorf("state has no devices")
	}
	dev := st.Devices[0]

	logf("# Channels transcript, recorded %s", time.Now().UTC().Format(time.RFC3339))
	logf("# upstream %s realm %q device %q", baseURL, st.Realm, dev.ID)
	logf("# socket: /appengine/v1/socket/websocket?vsn=2.0.0&realm=<realm>&token=<a_ch JWT>")
	logf("# The token is never printed; it is a realm-signed JWT carrying only the a_ch claims named per row.")

	fx := fixture{
		RecordedAt:      time.Now().UTC().Format(time.RFC3339),
		AstarteVersion:  "v1.2.0",
		Realm:           st.Realm,
		RecorderCommand: "go run ./recordchannels",
		SocketPath:      "/appengine/v1/socket/websocket?vsn=2.0.0&realm=<realm>&token=<a_ch JWT>",
	}

	// ASTARTE_UPSTREAM_SECTION re-records one section and carries the other over
	// from the committed fixture. Delivery of device_error to the room is
	// unreliable on this stack (see the README), so a full re-run to settle an
	// authorization question would churn attempts/delivered counts that were
	// measured carefully and are not what is being asked about.
	section := env("ASTARTE_UPSTREAM_SECTION", "all")
	if section != "all" && section != "auth" && section != "errors" {
		return fmt.Errorf("ASTARTE_UPSTREAM_SECTION must be all, auth or errors (got %q)", section)
	}
	if section != "errors" {
		fx.Authorization = recordAuthorization(dev.ID)
	}
	if section != "auth" {
		fx.DeviceErrors = recordDeviceErrors(dev)
	}

	// Resolve the output directory rather than guessing from the working
	// directory: `go run ./upstream/recordchannels` and `go run .` from inside
	// the package are both normal, and a wrong guess silently writes the
	// fixture somewhere nobody commits it.
	dir := os.Getenv("ASTARTE_UPSTREAM_OUT")
	if dir == "" {
		for _, cand := range []string{"upstream", "test/conformance/upstream", "."} {
			if fi, err := os.Stat(filepath.Join(cand, "README.md")); err == nil && !fi.IsDir() {
				dir = cand
				break
			}
		}
	}
	if dir == "" {
		return fmt.Errorf("cannot find test/conformance/upstream; set ASTARTE_UPSTREAM_OUT")
	}
	jsonPath := filepath.Join(dir, "channels.json")
	txtPath := filepath.Join(dir, "channels.transcript.txt")

	// Carry the un-recorded section over from the committed fixture, and refuse
	// to write a fixture with a section missing — a half-empty channels.json
	// would look like "upstream does nothing here" to every later reader.
	if section != "all" {
		var prev fixture
		b, err := os.ReadFile(jsonPath)
		if err != nil {
			return fmt.Errorf("reading %s to carry over the other section: %w", jsonPath, err)
		}
		if err := json.Unmarshal(b, &prev); err != nil {
			return fmt.Errorf("parsing %s: %w", jsonPath, err)
		}
		if section == "auth" {
			fx.DeviceErrors = prev.DeviceErrors
		} else {
			fx.Authorization = prev.Authorization
		}
		logf("\n# section=%s — the other section was carried over unchanged from %s", section, jsonPath)
	}
	if len(fx.Authorization) == 0 || len(fx.DeviceErrors) == 0 {
		return fmt.Errorf("refusing to write a fixture with an empty section "+
			"(authorization=%d, device_errors=%d)", len(fx.Authorization), len(fx.DeviceErrors))
	}

	out, err := json.MarshalIndent(fx, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(jsonPath, append(out, '\n'), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(txtPath, []byte(tr.String()), 0o644); err != nil {
		return err
	}
	fmt.Println("\nwrote", jsonPath, "and", txtPath)
	return nil
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// ------------------------------------------------------- authorization rows

// recordAuthorization establishes the a_ch grammar upstream actually enforces.
//
// Every rejection row is paired with an acceptance that differs only in the
// rule under test, because a refusal on its own proves nothing: a token that
// authorized nothing, or a server that refused every join, would satisfy the
// rejections and tell us we had learned something.
func recordAuthorization(deviceID string) []authObservation {
	logf("\n## Authorization: which a_ch claim authorizes which operation")

	join := func(name, why string, claims []string) authObservation {
		reply, ok := attempt(claims, "phx_join", nil, deviceID)
		return authObservation{Name: name, Why: why, Claims: claims,
			Operation: "phx_join rooms:<realm>:probe", Reply: reply, Accepted: ok}
	}
	watch := func(name, why string, claims []string, trigger map[string]any) authObservation {
		pl, _ := json.Marshal(trigger)
		reply, ok := attempt(claims, "watch", trigger, deviceID)
		return authObservation{Name: name, Why: why, Claims: claims,
			Operation: "watch", Payload: string(pl), Reply: reply, Accepted: ok}
	}

	devTrigger := func(on string, id string) map[string]any {
		return map[string]any{"name": "probe", "device_id": deviceID,
			"simple_trigger": map[string]any{"type": "device_trigger", "on": on, "device_id": id}}
	}
	dataTrigger := map[string]any{"name": "probe", "device_id": deviceID,
		"simple_trigger": map[string]any{"type": "data_trigger", "on": "incoming_data",
			"interface_name": "*", "match_path": "/*", "value_match_operator": "*"}}
	// pathTrigger names a concrete interface and match path, because the narrow
	// claims below have to be written as regexes: the wildcards in dataTrigger
	// would make the claim match by accident and prove nothing.
	//
	// interface_major is required as soon as interface_name is not "*", and
	// omitting it is not an authorization failure but a *validation* one —
	// `{"errors":{"simple_trigger":{"interface_major":["can't be blank"]}}}`.
	// That reply is a refusal like any other to a recorder that only checks for
	// "status":"ok", so both narrow rows would have been recorded as refused and
	// read as evidence about the authorization path. They are not: a validation
	// error is returned before authorization is consulted at all.
	pathTrigger := map[string]any{"name": "probe", "device_id": deviceID,
		"simple_trigger": map[string]any{"type": "data_trigger", "on": "incoming_data",
			"interface_name": ifaceIndividual, "interface_major": 1,
			"match_path": "/value", "value_match_operator": "*"}}

	rows := []authObservation{
		join("join with blanket .*::.*",
			"M11 decision 5 assumed a blanket grant authorizes everything. Upstream refuses it: the claim is matched against JOIN::<room>, and `.*::.*` requires a `::` the room name does not have.",
			[]string{".*::.*"}),
		join("join with blanket .*",
			"The other shape of a blanket grant, refused for the same reason. Paired with the acceptances below so the refusals cannot be explained by a server that refuses every join.",
			[]string{".*"}),
		join("join with JOIN::.*",
			"The acceptance that makes the two refusals above mean something: an explicit JOIN verb is what upstream requires.",
			[]string{"JOIN::.*"}),
		join("join with JOIN::<this room>",
			"An exact-room grant is honoured, so the verb is matched per room and not merely present.",
			[]string{"JOIN::probe"}),
		join("join with JOIN::<other room>",
			"The paired refusal for the row above: same verb, wrong room. Together they show the room name is really being matched.",
			[]string{"JOIN::other"}),

		watch("watch data_trigger with WATCH::.*",
			"A data trigger is authorized by the WATCH verb, as expected.",
			[]string{"JOIN::.*", "WATCH::.*"}, dataTrigger),
		watch("watch device_trigger, device_id inside simple_trigger",
			"A device trigger is authorized only when simple_trigger carries device_id itself.",
			[]string{"JOIN::.*", "WATCH::.*"}, devTrigger("device_error", deviceID)),
		watch("watch device_trigger, device_id only at payload top level",
			"The paired refusal: identical except that device_id sits where the AppEngine REST API puts it. Upstream calls this `unauthorized`, which is a misleading reason for a payload-shape problem — a client author would look at their token.",
			[]string{"JOIN::.*", "WATCH::.*"}, map[string]any{"name": "probe", "device_id": deviceID,
				"simple_trigger": map[string]any{"type": "device_trigger", "on": "device_error"}}),
		watch("watch device_trigger for device_id \"*\"",
			"A wildcard device is refused even under WATCH::.*, so device triggers cannot be watched realm-wide from a room.",
			[]string{"JOIN::.*", "WATCH::.*"}, devTrigger("device_error", "*")),

		// The rows above all used WATCH::.*, which accepts whatever string
		// upstream builds and therefore cannot say what that string IS. These
		// use narrow claims, so an acceptance identifies the authorization path
		// and a refusal excludes a candidate. Astrate builds
		// "<device>/<interface>"; upstream's source appends the trigger's match
		// path. Only one of the next two rows can be accepted.
		watch("watch data_trigger, claim WATCH::<device>/<interface> (Astrate's shape)",
			"Astrate authorizes a data trigger against `<device_id>/<interface>`, with no match path. If upstream built the same string this is accepted; if upstream appends the match path it is refused, because the anchored claim regex cannot match the longer string.",
			[]string{"JOIN::.*", "WATCH::" + deviceID + "/" + ifaceIndividual},
			pathTrigger),
		watch("watch data_trigger, claim WATCH::<device>/<interface><match_path>",
			"The paired alternative, and upstream's shape according to its source (`\"#{device_id}/#{interface_name}#{path}\"`, with no separator because the path already starts with a slash). Exactly one of this row and the one above can be accepted, so together they identify the authorization path rather than merely confirming a guess.",
			[]string{"JOIN::.*", "WATCH::" + deviceID + "/" + ifaceIndividual + "/value"},
			pathTrigger),
		watch("watch device_trigger, claim WATCH::<this device>",
			"A device trigger's authorization path is claimed to be the bare device id. A narrow, exact claim accepts it, which is what identifies the path.",
			[]string{"JOIN::.*", "WATCH::" + deviceID},
			devTrigger("device_error", deviceID)),
		watch("watch device_trigger, claim WATCH::<another device>",
			"The paired refusal for the row above: a well-formed claim naming a different device. Together they show the device id is really the string being matched, rather than the acceptance coming from any claim at all.",
			[]string{"JOIN::.*", "WATCH::" + otherDeviceID},
			devTrigger("device_error", deviceID)),
	}

	for _, r := range rows {
		logf("\n### %s\n  a_ch: %v\n  op: %s", r.Name, r.Claims, r.Operation)
		if r.Payload != "" {
			logf("  payload: %s", r.Payload)
		}
		logf("  <- %s", r.Reply)
		logf("  accepted: %v", r.Accepted)
		if validationRefusal(r.Reply) {
			// Loud on purpose. Upstream validates the payload before it consults
			// the token, so this row was refused for a reason that has nothing to
			// do with its claim — but it is a refusal like any other to a reader
			// scanning for "accepted: false", and reading it as evidence about
			// authorization is exactly how a recorder invents a finding.
			logf("  !! VALIDATION ERROR, NOT AN AUTHORIZATION REFUSAL — this row says")
			logf("  !! nothing about the a_ch grammar. Fix the payload and re-record.")
		}
	}
	return rows
}

// validationRefusal reports whether a refusal came from upstream's changeset
// validation rather than from authorization. An authorization refusal carries
// {"reason":"unauthorized"}; a validation refusal carries {"errors":{…}}.
func validationRefusal(reply string) bool {
	return strings.Contains(reply, `"status":"error"`) && strings.Contains(reply, `"errors":`)
}

// attempt opens a socket with the given claims, joins, optionally watches, and
// returns the reply frame for whichever operation was under test.
func attempt(claims []string, op string, trigger map[string]any, deviceID string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c, err := dial(ctx, claims)
	if err != nil {
		return "dial error: " + err.Error(), false
	}
	defer c.Close(websocket.StatusNormalClosure, "")

	room := "rooms:" + st.Realm + ":probe"
	send := func(ref, event string, payload any) {
		if payload == nil {
			payload = map[string]any{}
		}
		p, _ := json.Marshal([]any{"1", ref, room, event, payload})
		_ = c.Write(ctx, websocket.MessageText, p)
	}
	send("1", "phx_join", nil)
	wantRef := `"1"`
	if op == "watch" {
		send("2", "watch", trigger)
		wantRef = `"2"`
	}
	for {
		_, data, err := c.Read(ctx)
		if err != nil {
			return "read error: " + err.Error(), false
		}
		var f []json.RawMessage
		if json.Unmarshal(data, &f) != nil || len(f) < 5 {
			continue
		}
		if string(f[1]) != wantRef {
			continue
		}
		return string(data), strings.Contains(string(data), `"status":"ok"`)
	}
}

func dial(ctx context.Context, claims []string) (*websocket.Conn, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(st.RealmKey))
	if err != nil {
		return nil, fmt.Errorf("parsing realm key: %w", err)
	}
	tok, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"a_ch": claims,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour).Unix(),
	}).SignedString(key)
	if err != nil {
		return nil, err
	}
	ws := strings.Replace(strings.Replace(baseURL, "https://", "wss://", 1), "http://", "ws://", 1)
	c, _, err := websocket.Dial(ctx,
		ws+"/appengine/v1/socket/websocket?vsn=2.0.0&realm="+st.Realm+"&token="+url.QueryEscape(tok), nil)
	return c, err
}

// ------------------------------------------------------- device_error rows

func recordDeviceErrors(dev struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}) []observation {
	logf("\n## device_error: what upstream calls each rejection")
	logf("# One fresh websocket and one fresh device session per row. Several of these")
	logf("# provocations make upstream force a clean session and disconnect the device,")
	logf("# so a shared session would record the first row honestly and the rest as silence.")

	cases := []struct {
		name, why, topic string
		payload          []byte
	}{
		{"valid publish", "The acceptance row. Without it, a stack that emitted device_error for everything would satisfy every row below.",
			ifaceIndividual + "/value", bdoc(bDouble("v", 42))},
		{"interface not in introspection", "Astrate maps its own interface_not_in_introspection to invalid_interface. This row is what upstream actually says.",
			"org.astrate.nonexistent.Iface/value", bdoc(bDouble("v", 1))},
		{"unknown path on a known interface", "Astrate maps unexpected_path to mapping_not_found.",
			ifaceIndividual + "/nosuchpath", bdoc(bDouble("v", 1))},
		{"payload that is not BSON at all", "Astrate maps unknown_format/malformed to undecodable_bson_payload.",
			ifaceIndividual + "/value", []byte("definitely not bson")},
		{"string where the mapping says double", "Astrate maps type_mismatch to unexpected_value_type.",
			ifaceIndividual + "/value", bdoc(bString("v", "hello"))},
		{"BSON document with no v key", "Astrate maps no_value to undecodable_bson_payload.",
			ifaceIndividual + "/value", bdoc(bDouble("x", 1))},
		{"object with an unexpected key", "Astrate maps bad_object to unexpected_object_key.",
			ifaceObject + "/data", bdoc(bSub("v", bdoc(bDouble("temp", 1), bDouble("bogus_key", 2))))},
		{"unknown control message", "Astrate maps control_unknown to unexpected_control_message.",
			"control/bogusControl", []byte("1")},
		{"invalid introspection", "Astrate maps introspection_invalid to invalid_introspection.",
			"", []byte("this is not an introspection")},
	}

	const attempts = 3

	var out []observation
	for _, tc := range cases {
		agg := observation{Name: tc.name, Why: tc.why, Topic: tc.topic,
			Payload: b64(tc.payload), Attempts: attempts}
		for i := 0; i < attempts; i++ {
			o := provoke(dev.ID, dev.Secret, tc.name, tc.why, tc.topic, tc.payload)
			if o.ErrorName != "" {
				agg.Delivered++
				if agg.ErrorName == "" {
					agg.ErrorName, agg.Metadata = o.ErrorName, o.Metadata
				} else if agg.ErrorName != o.ErrorName {
					// Two different names for one provocation would mean the
					// row is not measuring what it claims to. Record it loudly
					// rather than letting first-wins hide it.
					logf("  !! attempt %d gave %q, earlier attempts gave %q",
						i+1, o.ErrorName, agg.ErrorName)
					agg.ErrorName += " | " + o.ErrorName
				}
			}
			agg.Frames = append(agg.Frames, o.Frames...)
			agg.UncorrelatedFrames = append(agg.UncorrelatedFrames, o.UncorrelatedFrames...)
		}
		agg.NoEvent = agg.Delivered == 0
		logf("  => %s: %s, delivered %d/%d", tc.name,
			orNone(agg.ErrorName), agg.Delivered, attempts)
		out = append(out, agg)
	}
	return out
}

func orNone(s string) string {
	if s == "" {
		return "(never reached the room)"
	}
	return s
}

func provoke(deviceID, secret, name, why, topic string, payload []byte) observation {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	o := observation{Name: name, Why: why, Topic: topic, Payload: b64(payload)}

	c, err := dial(ctx, []string{"JOIN::.*", "WATCH::.*"})
	if err != nil {
		o.Frames = []string{"dial error: " + err.Error()}
		return o
	}
	defer c.Close(websocket.StatusNormalClosure, "")

	room := "rooms:" + st.Realm + ":probe"
	write := func(ref, event string, payload any) {
		p, _ := json.Marshal([]any{"1", ref, room, event, payload})
		_ = c.Write(ctx, websocket.MessageText, p)
	}
	write("1", "phx_join", map[string]any{})
	write("2", "watch", map[string]any{"name": "err", "device_id": deviceID,
		"simple_trigger": map[string]any{"type": "device_trigger", "on": "device_error", "device_id": deviceID}})

	frames := make(chan string, 256)
	closed := make(chan struct{})
	go func() {
		defer close(closed)
		for {
			_, data, err := c.Read(ctx)
			if err != nil {
				return
			}
			frames <- string(data)
		}
	}()
	// Phoenix closes an idle socket. Without this heartbeat the socket dies
	// part-way through a run and every later row records silence that looks
	// exactly like "upstream emitted nothing".
	go func() {
		t := time.NewTicker(15 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-closed:
				return
			case <-t.C:
				p, _ := json.Marshal([]any{nil, "hb", "phoenix", "heartbeat", map[string]any{}})
				if c.Write(ctx, websocket.MessageText, p) != nil {
					return
				}
			}
		}
	}()

	mc, err := session(deviceID, secret)
	if err != nil {
		o.Frames = []string{"session error: " + err.Error()}
		return o
	}
	defer mc.Disconnect(200)
	drain(frames, 2*time.Second) // discard session-setup noise

	base := st.Realm + "/" + deviceID
	full := base
	if topic != "" {
		full = base + "/" + topic
	}
	logf("\n### %s\n  topic: %s\n  payload (base64): %s", name, full, b64(payload))

	tk := mc.Publish(full, 2, false, payload)
	if !tk.WaitTimeout(15*time.Second) || tk.Error() != nil {
		logf("  publish problem: %v", tk.Error())
	}

	for _, f := range drain(frames, 12*time.Second) {
		if strings.Contains(f, `"phoenix"`) {
			continue // heartbeat replies
		}
		n, keys, echoed := parseDeviceError(f)
		if n == "" {
			o.Frames = append(o.Frames, f)
			logf("  <- %s", f)
			continue
		}
		if echoed != o.Payload {
			o.UncorrelatedFrames = append(o.UncorrelatedFrames, f)
			logf("  (uncorrelated, belongs to the session) <- %s", f)
			continue
		}
		o.Frames = append(o.Frames, f)
		logf("  <- %s", f)
		if o.ErrorName == "" {
			o.ErrorName = n
			o.Metadata = keys
		}
	}
	if o.ErrorName == "" {
		logf("  (no device_error reached the room on this attempt)")
	}
	return o
}

// parseDeviceError pulls the error_name, the metadata keys and the echoed
// base64 payload out of a new_event frame. It returns "" for any frame that is
// not a device_error. The echoed payload is what lets a row be attributed to
// the publish that caused it.
func parseDeviceError(frame string) (string, []string, string) {
	var f []json.RawMessage
	if json.Unmarshal([]byte(frame), &f) != nil || len(f) < 5 {
		return "", nil, ""
	}
	var pl struct {
		Event struct {
			Type      string            `json:"type"`
			ErrorName string            `json:"error_name"`
			Metadata  map[string]string `json:"metadata"`
		} `json:"event"`
	}
	if json.Unmarshal(f[4], &pl) != nil || pl.Event.Type != "device_error" {
		return "", nil, ""
	}
	keys := make([]string, 0, len(pl.Event.Metadata))
	for k := range pl.Event.Metadata {
		keys = append(keys, k)
	}
	return pl.Event.ErrorName, keys, pl.Event.Metadata["base64_payload"]
}

func drain(ch chan string, d time.Duration) []string {
	var out []string
	deadline := time.After(d)
	for {
		select {
		case f := <-ch:
			out = append(out, f)
		case <-deadline:
			return out
		}
	}
}

// ------------------------------------------------------------ device session

// session brings up a fresh mTLS device session: new credentials through the
// pairing API, connect, introspection, emptyCache, then a settle window.
// Without the settle the first publishes race Data Updater Plant's registration
// and come back as device_session_not_found, which would be recorded as the
// answer to whatever was being provoked.
func session(deviceID, secret string) (paho.Client, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	cn := st.Realm + "/" + deviceID
	csrDER, err := x509.CreateCertificateRequest(rand.Reader,
		&x509.CertificateRequest{Subject: pkix.Name{CommonName: cn}}, priv)
	if err != nil {
		return nil, err
	}
	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})
	certPEM, err := obtainCert(deviceID, secret, string(csrPEM))
	if err != nil {
		return nil, err
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}
	pair, err := tls.X509KeyPair([]byte(certPEM),
		pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER}))
	if err != nil {
		return nil, err
	}
	dialURL := strings.Replace(strings.TrimSuffix(st.BrokerURL, "/"), "mqtts://", "ssl://", 1)
	mc := paho.NewClient(paho.NewClientOptions().AddBroker(dialURL).SetClientID(cn).
		SetTLSConfig(&tls.Config{ //nolint:gosec // the dev stack uses a self-signed broker cert
			Certificates:       []tls.Certificate{pair},
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}).
		SetConnectTimeout(20 * time.Second).SetCleanSession(true))
	t := mc.Connect()
	if !t.WaitTimeout(25*time.Second) || t.Error() != nil {
		return nil, fmt.Errorf("connecting: %v", t.Error())
	}
	mc.Publish(cn, 2, false, introspection).WaitTimeout(10 * time.Second)
	mc.Publish(cn+"/control/emptyCache", 2, false, "1").WaitTimeout(10 * time.Second)
	time.Sleep(8 * time.Second)
	return mc, nil
}

func obtainCert(id, secret, csr string) (string, error) {
	body, _ := json.Marshal(map[string]any{"data": map[string]string{"csr": csr}})
	u := st.Endpoints.Pairing + "/v1/" + st.Realm + "/devices/" + id +
		"/protocols/astarte_mqtt_v1/credentials"
	req, err := http.NewRequest("POST", u, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+secret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var out struct {
		Data struct {
			ClientCrt string `json:"client_crt"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if out.Data.ClientCrt == "" {
		return "", fmt.Errorf("pairing returned no certificate (status %d)", resp.StatusCode)
	}
	return out.Data.ClientCrt, nil
}

// ------------------------------------------- minimal BSON, malformed on purpose

func bdoc(fields ...[]byte) []byte {
	var b bytes.Buffer
	for _, f := range fields {
		b.Write(f)
	}
	b.WriteByte(0)
	out := make([]byte, 4+b.Len())
	binary.LittleEndian.PutUint32(out, uint32(4+b.Len()))
	copy(out[4:], b.Bytes())
	return out
}

func bDouble(k string, v float64) []byte {
	var b bytes.Buffer
	b.WriteByte(0x01)
	b.WriteString(k)
	b.WriteByte(0)
	_ = binary.Write(&b, binary.LittleEndian, math.Float64bits(v))
	return b.Bytes()
}

func bString(k, v string) []byte {
	var b bytes.Buffer
	b.WriteByte(0x02)
	b.WriteString(k)
	b.WriteByte(0)
	_ = binary.Write(&b, binary.LittleEndian, uint32(len(v)+1))
	b.WriteString(v)
	b.WriteByte(0)
	return b.Bytes()
}

func bSub(k string, d []byte) []byte {
	var b bytes.Buffer
	b.WriteByte(0x03)
	b.WriteString(k)
	b.WriteByte(0)
	b.Write(d)
	return b.Bytes()
}

func b64(b []byte) string { return base64.StdEncoding.EncodeToString(b) }
