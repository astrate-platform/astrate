# Minimal echo container for Astrate Flow container-block PoC (#43).
# Build: docker build -t astrate/flow-container-echo:poc .
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY main.go .
RUN go mod init echo && CGO_ENABLED=0 go build -o /echo .

FROM alpine:3.21
RUN adduser -D -H -u 10001 echo
USER echo
COPY --from=build /echo /echo
EXPOSE 8080
ENTRYPOINT ["/echo"]
