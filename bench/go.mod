// Standalone benchmark harness (see README.md). Deliberately a separate
// module that imports NOTHING from astrate: it speaks only the wire protocols
// (Astarte MQTT v1 + the REST APIs), so the same binary drives an upstream
// Astarte deployment and an Astrate deployment interchangeably.
module github.com/astrate-platform/astrate/bench

go 1.26.1

require (
	github.com/eclipse/paho.mqtt.golang v1.5.1 // same pin as the main module
	github.com/golang-jwt/jwt/v5 v5.3.1
	go.mongodb.org/mongo-driver/v2 v2.6.0
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
)
