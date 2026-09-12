# Swarm User Manual

## Current supported use

Swarm is currently a development service foundation, not a usable BitTorrent client. The present build is intended for developer validation only.

## Start the development service

Requires Go 1.23 or newer compatible with the current module declaration.

```bash
go run ./cmd/swarmd
```

The service binds to `127.0.0.1:8080` by default. A different **loopback** address may be selected:

```bash
SWARM_LISTEN_ADDR=127.0.0.1:9090 go run ./cmd/swarmd
```

Non-loopback addresses are intentionally rejected because remote authentication and authorization are not implemented.

## Check service state

```bash
curl http://127.0.0.1:8080/api/v1/health
```

The expected current result is `status: degraded` and `ready: false` because no production BitTorrent engine adapter is connected.

## Engine-development configuration

A compatible sidecar may be configured only with an absolute executable path and an explicit existing absolute download root:

```bash
SWARM_ENGINE_PATH=/absolute/path/to/swarm-engine \
SWARM_DOWNLOAD_ROOT=/absolute/path/to/downloads \
go run ./cmd/swarmd
```

Swarm resolves the download root before starting the sidecar and fails closed if the engine path, root, or protocol handshake is invalid. No production libtorrent sidecar is included yet.

## Current transfer behavior

Submitting a valid magnet link to `POST /api/v1/transfers` with the default unavailable engine returns HTTP `503` with `engine_unavailable`. This is intentional and prevents the development API from pretending a transfer began when no engine exists.

Validated raw `.torrent` metainfo may be submitted to:

```text
POST /api/v1/transfers/torrent
```

Swarm parses the metainfo, determines whether it is v1, v2, or hybrid, checks the engine's advertised capability, and only then dispatches the already validated bytes to the Engine contract. With the default unavailable engine, this does not start a real transfer.

## Pause/resume/remove lifecycle

The development API defines lifecycle routes for accepted transfers:

- `POST /api/v1/transfers/{id}/pause`
- `POST /api/v1/transfers/{id}/resume`
- `DELETE /api/v1/transfers/{id}?delete_data=false`

These routes are operational only when a configured engine has accepted the transfer. With the default unavailable engine, no transfer can be created, so these controls are not yet usable for real BitTorrent work. `delete_data=true` is an explicit destructive intent flag and remains subject to future storage-policy authorization before production use.

## Troubleshooting

- If the configured port is already in use, select another loopback port with `SWARM_LISTEN_ADDR`.
- If the process refuses to start on `0.0.0.0` or another external address, this is expected security behavior in the current development stage.
- If a configured engine fails startup, verify that both `SWARM_ENGINE_PATH` and `SWARM_DOWNLOAD_ROOT` are absolute and that the download root already exists.

Normal torrent downloading, desktop UI, remote management, and mobile operation are not yet available.
