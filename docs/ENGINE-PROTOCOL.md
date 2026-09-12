# Swarm Engine Sidecar Protocol

## Status

Development protocol, version 1. The Go service implementation exists in `internal/engine/sidecar.go`. No production BitTorrent sidecar is connected yet.

## Purpose

The sidecar protocol keeps the BitTorrent implementation behind a GoreeCloud-owned process and policy boundary. `swarmd` remains authoritative for application state, API behavior, privacy/security decisions, storage policy, and platform integration. The sidecar performs protocol execution only.

## Process launch

A sidecar is enabled only when `SWARM_ENGINE_PATH` contains an absolute executable path. `SWARM_DOWNLOAD_ROOT` must also identify an absolute existing directory.

Security properties of the current launcher:

- no shell is used;
- no PATH lookup is used;
- a relative executable path is rejected;
- the configured download root is required, resolved through symlinks once at startup, and must remain a directory;
- the child does not inherit the `swarmd` environment;
- the child receives only `SWARM_ENGINE_PROTOCOL_VERSION=1` through its environment;
- the resolved download root is sent explicitly in the versioned handshake;
- stderr may be forwarded to the service diagnostic sink;
- stdin/stdout are reserved for the protocol.

If an explicitly configured sidecar cannot start or complete its handshake, `swarmd` fails closed instead of silently falling back to another engine.

## Framing

Messages are UTF-8 JSON objects separated by a newline. Each message is limited to 24 MiB by the client implementation. This accommodates the 16 MiB torrent-metainfo ingestion ceiling after JSON/base64 expansion while retaining a fixed IPC resource bound.

Request:

```json
{"id":1,"method":"hello","params":{"protocol_version":1,"client":"swarmd","download_root":"/authorized/download/root"}}
```

Successful response:

```json
{"id":1,"result":{}}
```

Error response:

```json
{"id":1,"error":{"code":"unavailable","message":"engine is not ready"}}
```

Response IDs must exactly match request IDs.

## Handshake

The first request is `hello`.

The engine receives the Swarm-approved resolved download root and must return:

- the accepted protocol version;
- a truthful readiness state;
- engine name and version;
- only capabilities actually supported by the running engine build/configuration.

Example:

```json
{
  "id": 1,
  "result": {
    "protocol_version": 1,
    "capabilities": {
      "ready": true,
      "name": "libtorrent-rasterbar",
      "version": "<runtime version>",
      "supported": ["magnet", "torrent_v1", "torrent_v2", "torrent_hybrid", "dht", "pex", "ipv6"]
    }
  }
}
```

A protocol-version mismatch or `ready:false` fails startup.

## Version 1 methods

Current client methods:

- `hello`
- `add`
- `add_torrent`
- `pause`
- `resume`
- `remove`
- `shutdown`

`add_torrent` carries the already Swarm-validated raw metainfo bytes through JSON byte encoding. The protocol will be extended with live transfer state, files, peers, trackers, pieces, network policy, and persistence/recovery operations as those behaviors are implemented.

## Authority and state

A successful transport-level response does not automatically establish that a transfer is private, secure, persisted, protected, or complete. The responsible Swarm subsystem must independently accept and expose those states.

The download root establishes only the initial filesystem authority boundary supplied to the engine. Future Storage Manager operations must further constrain file allocation, relocation, symlink behavior, removable storage, and destructive deletion.

## Failure behavior

Malformed JSON, oversized messages, response-ID mismatches, protocol-version mismatches, invalid download-root configuration, unexpected stream closure, and engine-reported errors are treated as explicit failures. Context expiry terminates the engine process rather than allowing an unbounded request to remain active.
