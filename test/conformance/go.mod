// Conformance harness module (docs/ROADMAP.md §10): official Astarte SDK and
// tooling dependencies live here, isolated from the main astrate module.
// Conformance clients are pinned deliberately (docs/ROADMAP.md §0.3) and
// upgraded on purpose, never by drift:
//
//   - astarte-device-sdk-go v0.90.2 (latest tagged release at M4 start)
//   - astarte-go v0.90.4 (the client library the SDK delegates pairing to)
//   - astartectl v26.5.0 (release binary, fetched by the harness — see
//     cpa/astartectl_test.go)
module github.com/astrate-platform/astrate/test/conformance

go 1.26.1

require (
	github.com/astarte-platform/astarte-device-sdk-go v0.90.2
	github.com/astarte-platform/astarte-go v0.90.4
	github.com/astrate-platform/astrate v0.0.0-00010101000000-000000000000
	github.com/jackc/pgx/v5 v5.10.0
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cristalhq/jwt/v3 v3.1.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/golang-migrate/migrate/v4 v4.19.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/iancoleman/orderedmap v0.2.0 // indirect
	github.com/ispirata/paho.mqtt.golang v1.3.91-0.20220121155423-6e36d43b2ec9 // indirect
	github.com/jackc/pgerrcode v0.0.0-20220416144525-469b46aa5efa // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.3 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	go.mongodb.org/mongo-driver v1.8.0 // indirect
	golang.org/x/crypto v0.53.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	gorm.io/driver/sqlite v1.2.6 // indirect
	gorm.io/gorm v1.22.3 // indirect
)

replace github.com/astrate-platform/astrate => ../..
