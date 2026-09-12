# GoreeCloud Swarm

GoreeCloud Swarm is the planned native GoreeCloud peer-to-peer transfer platform and BitTorrent client.

## Development status

**Status: Development / pre-alpha.**

The repository contains the native service/metainfo foundation plus an **opt-in development engine adapter** for `anacrolix/torrent` v1.61.0. Ordinary builds still use the unavailable engine. The tagged adapter is not considered accepted until hosted exact-head validation and controlled interoperability/runtime tests pass.

Implemented foundation:

- Native GoreeCloud repository and service layout.
- Authoritative service/core separation.
- Replaceable transfer-engine interface.
- Local-only HTTP API default (`127.0.0.1`).
- `/api/v1/health` and `/api/v1/capabilities` endpoints.
- Bounded native bencode decoding.
- v1/v2/hybrid `.torrent` metainfo inspection at `/api/v1/metainfo/inspect`.
- Raw `info` dictionary hashing for v1 and v2 identifiers.
- Magnet URI parsing and validation.
- Download-path sanitization primitives.
- Transfer lifecycle domain model.
- Opt-in Anacrolix engine adapter source for magnet add, pause, resume, non-destructive remove, and shutdown.
- Fail-closed engine selection and explicit download-root requirement.
- Unit and API tests.
- Separate CI gates for default builds and the tagged engine build.

Not yet accepted or implemented:

- Accepted/validated peer-transfer support in the default build.
- Starting transfers from inspected `.torrent` files.
- Persistent resume/session state.
- Network Lock and proxy routing.
- Storage deletion/relocation through the engine adapter.
- Engine process isolation/sandboxing.
- Authentication/authorization for remote management.
- Web, desktop, mobile, or CLI clients beyond the initial `swarmd` service.
- GoreeCloud Manager, Privacy Shield, Wardveil Security, Everkeep, Glaze UI, GoreeCloud Mesh, GoreeCloud Identity, or GoreeCloud Sync runtime integrations.

Those capabilities remain development/planned until implementation evidence exists.

## Run the default development service

```bash
go run ./cmd/swarmd
```

The default build binds to `127.0.0.1:8080` and keeps the transfer engine unavailable.

```bash
SWARM_LISTEN_ADDR=127.0.0.1:9090 go run ./cmd/swarmd
```

Check health:

```bash
curl http://127.0.0.1:8080/api/v1/health
```

Inspect a torrent without starting it:

```bash
curl --data-binary @example.torrent http://127.0.0.1:8080/api/v1/metainfo/inspect
```

## Experimental Anacrolix engine build

The first engine adapter is deliberately opt-in:

```bash
go build -mod=mod -tags anacrolix_engine -o swarmd ./cmd/swarmd
SWARM_ENGINE=anacrolix \
SWARM_DOWNLOAD_DIR=/absolute/path/to/swarm-downloads \
./swarmd
```

This development mode can create peer-to-peer network traffic and write transfer data under the configured download root. It is not production-ready and must not be treated as providing Network Lock, proxy enforcement, anonymity, Wardveil protection, Privacy Shield acceptance, or recovery guarantees.

See [`docs/ENGINE-DEPENDENCY.md`](docs/ENGINE-DEPENDENCY.md) and [`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md).

## Validate

Default build:

```bash
gofmt -w .
go vet ./...
go test ./...
go build ./cmd/swarmd
```

Tagged engine build (requires dependency resolution):

```bash
go test -mod=mod -tags anacrolix_engine ./internal/engine ./cmd/swarmd
go build -mod=mod -tags anacrolix_engine ./cmd/swarmd
```

## Repository controls

Repository-coupled specifications, feature state, roadmap, privacy, security, branding, and platform-contract records are maintained at the repository root. The canonical full product specification remains under `GoreeCloud/Projects` in Google Drive.

## Architecture

See [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md).

## Repository origin and dependency policy

This repository is a GoreeCloud-created native product repository. It is not a direct fork of another BitTorrent client. Third-party protocol code may be used only as a bounded dependency behind GoreeCloud-owned adapter and policy boundaries, with provenance, licensing, security implications, and replacement boundaries documented.
