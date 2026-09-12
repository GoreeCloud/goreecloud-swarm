# GoreeCloud Swarm

GoreeCloud Swarm is the planned native GoreeCloud peer-to-peer transfer platform and BitTorrent client.

## Development status

**Status: Development / pre-alpha.**

The repository currently contains the first native service foundation. It does **not** yet perform BitTorrent transfers. The service reports the transfer engine as unavailable until a real, validated engine adapter is connected.

Implemented in this foundation:

- Native GoreeCloud repository and service layout.
- Authoritative service/core separation.
- Replaceable transfer-engine interface.
- Local-only HTTP API default (`127.0.0.1`).
- `/api/v1/health` and `/api/v1/capabilities` endpoints.
- Truthful degraded health while no transfer engine is attached.
- Magnet URI parsing and validation.
- Download-path sanitization primitives.
- Transfer lifecycle domain model.
- Unit and API tests.
- CI build, format, vet, and test gates.

Not yet implemented:

- Peer connections or piece transfer.
- Trackers, DHT, PEX, LPD, or web seeds.
- `.torrent` import and bencode parsing.
- Persistence/resume database.
- Network Lock and proxy routing.
- Authentication/authorization for remote management.
- Web, desktop, mobile, or CLI clients beyond the initial `swarmd` service.
- GoreeCloud Manager, Privacy Shield, Wardveil Security, Everkeep, Glaze UI, GoreeCloud Mesh, GoreeCloud Identity, or GoreeCloud Sync runtime integrations.

Those capabilities remain planned until implementation evidence exists.

## Run locally

```bash
go run ./cmd/swarmd
```

The development service binds to `127.0.0.1:8080` by default. Override the address with `SWARM_LISTEN_ADDR`.

```bash
SWARM_LISTEN_ADDR=127.0.0.1:9090 go run ./cmd/swarmd
```

Check health:

```bash
curl http://127.0.0.1:8080/api/v1/health
```

## Validate

```bash
gofmt -w .
go vet ./...
go test ./...
go build ./cmd/swarmd
```

## Repository controls

Repository-coupled specifications, feature state, roadmap, privacy, security, branding, and platform-contract records are maintained at the repository root. The canonical full product specification remains under `GoreeCloud/Projects` in Google Drive.

## Architecture

See [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md).

## Repository origin and dependency policy

This repository is a GoreeCloud-created native product repository. It is not a direct fork of another BitTorrent client. Mature third-party protocol engines may be integrated later only behind GoreeCloud-owned adapter and policy boundaries, with provenance and licensing documented at the time they are introduced.
