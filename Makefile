# Astrate build & verification entry points (docs/ROADMAP.md §1.2, verification tiers §0.2).
#
#   make build               static binary / compile check (CGO disabled)
#   make lint                golangci-lint (config: .golangci.yml)
#   make test                T1 unit tests (no Docker, no network)
#   make test-integration    T2 tests (testcontainers -> timescale/timescaledb:latest-pg16)
#   make test-e2e            T3 component/E2E tests (full wired binary; populated from M5/M6)
#   make test-conformance    T4 official-SDK conformance harness (lands in M9 under test/)
#   make up / make down      local TimescaleDB via docker compose

GO                     ?= go
GOLANGCI_LINT_VERSION  ?= v2.12.2
GOFLAGS                := -trimpath
LDFLAGS                := -s -w
DIST                   := dist

.DEFAULT_GOAL := build

.PHONY: build lint test test-integration test-e2e test-conformance tools up down clean

## build: compile every package statically; emits $(DIST)/astrate once cmd/astrate lands (M8).
build:
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' ./...
	@if [ -d cmd/astrate ]; then \
		mkdir -p $(DIST); \
		CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(DIST)/astrate ./cmd/astrate; \
		echo "built $(DIST)/astrate"; \
	fi

## lint: run the pinned golangci-lint (install with `make tools`).
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not found; run 'make tools' (installs $(GOLANGCI_LINT_VERSION))"; exit 1; }
	golangci-lint run ./...

## test: T1 — pure unit tests, race detector on.
test:
	$(GO) test -race ./...

## test-integration: T2 — requires a Docker daemon (or ASTRATE_TEST_DSN to reuse a database).
test-integration:
	$(GO) test -race -count=1 -tags integration ./...

## test-e2e: T3 — component tests behind the e2e build tag (suites land from M5 onward).
test-e2e:
	$(GO) test -race -count=1 -tags "integration e2e" ./...

## test-conformance: T4 — official-SDK harness; lives in test/ with its own go.mod (M9).
test-conformance:
	@if [ -f test/conformance/go.mod ]; then \
		cd test/conformance && $(GO) test -count=1 ./...; \
	else \
		echo "conformance harness lands in M9 (docs/ROADMAP.md §10); nothing to run yet"; \
	fi

## tools: install the pinned developer toolchain.
tools:
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

## up/down: local TimescaleDB for development and DSN-reuse test runs.
up:
	docker compose up -d

down:
	docker compose down

clean:
	rm -rf $(DIST)
