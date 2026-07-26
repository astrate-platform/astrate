// Standalone capture harness (see main.go); kept out of the main Astrate
// module so the official SDK's mongo-driver v1 dependency tree never leaks
// into the production go.mod. Versions are pinned to exactly what
// astarte-device-sdk-go v0.90.2 declares.
module github.com/astrate-platform/astrate/tools/bsoncapture

go 1.24

require (
	github.com/astarte-platform/astarte-go v0.90.4
	go.mongodb.org/mongo-driver v1.8.0
)

require github.com/go-stack/stack v1.8.0 // indirect
