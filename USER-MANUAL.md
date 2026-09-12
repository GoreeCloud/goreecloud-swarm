# GoreeCloud Swarm — Development User Manual

## Current scope

Swarm is in development/pre-alpha. The default build is a service and metainfo-inspection foundation. An opt-in tagged build also contains the first development BitTorrent engine adapter. Neither mode is production-ready.

## Requirements

- Go 1.23 or newer.
- A local development environment capable of binding a loopback TCP port.
- For the experimental engine build, dependency resolution for `github.com/anacrolix/torrent` v1.61.0.

## Start the default service

```bash
go run ./cmd/swarmd
```

The default listener is `127.0.0.1:8080`. A different loopback address may be supplied:

```bash
SWARM_LISTEN_ADDR=127.0.0.1:9090 go run ./cmd/swarmd
```

The service intentionally refuses non-loopback addresses until authenticated remote management is implemented.

## Check health and capabilities

```bash
curl http://127.0.0.1:8080/api/v1/health
curl http://127.0.0.1:8080/api/v1/capabilities
```

For an ordinary build, the expected state is `degraded` with `ready: false` because the transfer engine is unavailable.

## Inspect a `.torrent` file

The development API can inspect bounded v1, v2, or hybrid metainfo without starting a transfer:

```bash
curl --data-binary @example.torrent \
  http://127.0.0.1:8080/api/v1/metainfo/inspect
```

The response may include the torrent name, format, private flag, piece length, file count, total size, accepted tracker URLs, and v1/v2 info hashes. The current inspection endpoint accepts at most 16 MiB of metainfo.

## Default transfer endpoint behavior

A normal development build can exercise magnet validation with:

```bash
curl -i \
  -H 'Content-Type: application/json' \
  -d '{"magnet_uri":"magnet:?xt=urn:btih:0123456789abcdef0123456789abcdef01234567"}' \
  http://127.0.0.1:8080/api/v1/transfers
```

Because the ordinary build contains no active transfer engine, a valid request returns HTTP 503 with `engine_unavailable`.

## Experimental Anacrolix engine mode

The first real engine adapter is deliberately excluded from ordinary builds. Build it explicitly:

```bash
go build -mod=mod -tags anacrolix_engine -o swarmd ./cmd/swarmd
```

Start it only with an explicit download root:

```bash
SWARM_ENGINE=anacrolix \
SWARM_DOWNLOAD_DIR=/absolute/path/to/swarm-downloads \
./swarmd
```

In this mode a valid magnet request may start real peer-to-peer network activity and write transfer data beneath the configured download directory.

The tagged adapter is development-only and has not yet completed hosted exact-head validation, controlled transfer acceptance, Network Lock, proxy enforcement, process isolation, persistence, or GoreeCloud Platform System runtime acceptance.

Do not expose the local management API to non-loopback addresses. Do not treat the tagged engine as anonymous, VPN-bound, Wardveil-protected, Privacy-Shield-accepted, or recovery-ready.

## Validate a development checkout

Default build:

```bash
gofmt -w .
go vet ./...
go test ./...
go build ./cmd/swarmd
```

Tagged adapter:

```bash
go test -mod=mod -tags anacrolix_engine ./internal/engine ./cmd/swarmd
go build -mod=mod -tags anacrolix_engine ./cmd/swarmd
```

## Unsupported or unaccepted in this development state

Starting transfers from inspected `.torrent` files, durable persistence/resume, Network Lock, proxy routing, engine process isolation, remote administration, Web UI, desktop UI, mobile clients, safe engine-managed data deletion, and production deployment remain unsupported or unaccepted.
