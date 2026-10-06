slug: docs-native-socket-handshake-errors
verdict: blocked
at:  e3f5780
ran: 2026-10-06T18:57:40Z on DietPi in 91s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astrate_native_api.yaml [offset=300, limit=130]
→ Read docs/api/astrate_native_api.yaml [limit=80, offset=430]
✱ Grep "content:|text/plain|application/json" in docs/api/astrate_native_api.yaml · 100 matches

$ grep -rn "coder/websocket" go.mod; ls docs/
9:	github.com/coder/websocket v1.8.14 // live stream socket (M7)
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

→ Read internal/appengine/stream/ws.go [limit=110]
→ Read internal/appengine/channels/ws.go [limit=100]

$ cat go.mod | head -40; echo ---; cat docs/Makefile; echo ---; ls docs/api
module github.com/astrate-platform/astrate

go 1.26.1

// Pinned dependency decisions — docs/ROADMAP.md §1.1 (frozen in docs/DESIGN.md).
// HTTP routing is stdlib net/http (Go >= 1.22 ServeMux patterns): deliberately no framework.
require (
	github.com/BurntSushi/toml v1.6.0 // TOML config (§5.1)
	github.com/coder/websocket v1.8.14 // live stream socket (M7)
	github.com/eclipse/paho.mqtt.golang v1.5.1 // test-only MQTT client (same family as official Go SDK)
	github.com/golang-jwt/jwt/v5 v5.3.1 // JWT verification, alg allowlist (§4.2)
	github.com/golang-migrate/migrate/v4 v4.19.1 // migrations via source/iofs + go:embed (§1.3)
	github.com/hashicorp/golang-lru/v2 v2.0.7 // verified-token LRU cache (§4.2)
	github.com/jackc/pgx/v5 v5.10.0 // Postgres/TimescaleDB driver + pgxpool (§1.3)
	github.com/mochi-mqtt/server/v2 v2.7.9 // embedded MQTT broker (§1.1)
	github.com/prometheus/client_golang v1.23.2 // metrics (§5.2)
	github.com/testcontainers/testcontainers-go v0.43.0 // T2 integration harness (§5.4)
	github.com/testcontainers/testcontainers-go/modules/postgres v0.42.0 // Timescale container module
	go.etcd.io/bbolt v1.4.3 // broker session store (§3.1; pinned over pebble)
	go.mongodb.org/mongo-driver/v2 v2.6.0 // BSON raw-document API only (§3.5.5)
	golang.org/x/crypto v0.53.0 // bcrypt cost 10 for credentials secrets (§4.1)
)

require (
	github.com/cbroglie/mustache v1.4.0
	github.com/nats-io/nats.go v1.52.0
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.3
	github.com/testcontainers/testcontainers-go/modules/nats v0.43.0
)

require (
	dario.cat/mergo v1.0.2 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
---
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
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

! permission requested: external_directory (/tmp/*, /tmp/wsprobe/*); auto-rejecting
✗ grep -n "^replace\|^=> " go.mod; mkdir -p /tmp/wsprobe && sed '1s|.*|module wsprobe|' go.mod > /tmp/wsprobe/go.mod && cp go.sum /tmp/wsprobe/go.sum && grep -c websocket /tmp/wsprobe/go.sum failed
Error: The user rejected permission to use this specific tool call.
