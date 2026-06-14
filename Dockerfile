# Astrate — single static binary on distroless (docs/ROADMAP.md §9 file 8.6,
# docs/DESIGN.md §4.5/§5.4). Two stages: a Go builder and a non-root,
# read-only distroless runtime carrying only the binary.

# --- build -------------------------------------------------------------------
FROM golang:1.26 AS build
WORKDIR /src

# Cache modules before copying the full tree.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=0.1.0-dev
RUN CGO_ENABLED=0 go build -trimpath \
        -ldflags "-s -w -X main.version=${VERSION}" \
        -o /astrate ./cmd/astrate

# Stage the session-store directory so the runtime's named volume inherits
# non-root ownership on first mount (distroless has no shell to mkdir/chown).
RUN mkdir -p /data

# --- runtime -----------------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /astrate /astrate
COPY --from=build --chown=65532:65532 /data /var/lib/astrate

USER 65532:65532
# HTTP API, mTLS MQTT, and (insecure_dev_mode only) plaintext MQTT.
EXPOSE 8080 8883 1883
# Only writable path; the rest of the FS can be mounted read-only.
VOLUME ["/var/lib/astrate"]

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=5 \
    CMD ["/astrate", "-healthcheck"]

ENTRYPOINT ["/astrate"]
