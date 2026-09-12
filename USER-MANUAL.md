# GoreeCloud Swarm — Development User Manual

## Current scope

Swarm is in development/pre-alpha. The current build is a service foundation for developers and testers; it is not yet a usable BitTorrent client.

## Requirements

- Go 1.23 or newer.
- A local development environment capable of binding a loopback TCP port.

## Start the service

```bash
go run ./cmd/swarmd
```

The default listener is:

```text
127.0.0.1:8080
```

A different loopback address may be supplied:

```bash
SWARM_LISTEN_ADDR=127.0.0.1:9090 go run ./cmd/swarmd
```

The service intentionally refuses non-loopback addresses until authenticated remote management is implemented.

## Check health

```bash
curl http://127.0.0.1:8080/api/v1/health
```

In the current foundation, the expected state is `degraded` with `ready: false` because no transfer engine adapter is connected.

## Inspect engine capability state

```bash
curl http://127.0.0.1:8080/api/v1/capabilities
```

## Transfer endpoint

A development caller may exercise request validation with:

```bash
curl -i \
  -H 'Content-Type: application/json' \
  -d '{"magnet_uri":"magnet:?xt=urn:btih:0123456789abcdef0123456789abcdef01234567"}' \
  http://127.0.0.1:8080/api/v1/transfers
```

Until the first real engine adapter exists, a valid request returns HTTP 503 with `engine_unavailable`. This is expected and prevents the application from claiming that a transfer was accepted when it was not.

## Validate a development checkout

```bash
gofmt -w .
go vet ./...
go test ./...
go build ./cmd/swarmd
```

## Unsupported in this foundation

Downloading, seeding, torrent-file import, persistence, remote administration, Web UI, desktop UI, mobile clients, and production deployment are not yet supported.
