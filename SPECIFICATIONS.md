# Swarm Repository Specifications

## Authority

The canonical product specification is **Project Specification — Swarm.docx** in `GoreeCloud/Projects` on the authoritative GoreeCloud Google Drive. This repository document records only implementation-coupled requirements needed to understand and validate the source tree; it does not replace the canonical specification.

## Current lifecycle

- Product: GoreeCloud Swarm
- Repository: `GoreeCloud/goreecloud-swarm`
- Development model: native GoreeCloud repository; not a direct upstream fork
- Lifecycle: Development / pre-alpha
- Production ready: no
- Stable eligible: no

## Current implemented foundation

The current repository implements the first Swarm service foundation, native torrent metainfo inspection, validated magnet/`.torrent` ingestion to the Engine contract, lifecycle control APIs, and a versioned isolated engine-sidecar client. It does not yet perform BitTorrent peer transfer.

Implemented foundation:

- `swarmd`, a local service process.
- A versioned `/api/v1/` HTTP surface.
- Loopback-only binding until remote authentication exists.
- Transfer-engine abstraction with explicit capability reporting.
- Isolated versioned engine-sidecar launcher and handshake.
- Explicit absolute/resolved engine download-root authority.
- Truthful degraded health when no engine is connected.
- Magnet URI validation.
- v1/v2/hybrid `.torrent` validation and Engine-contract dispatch.
- Storage path-boundary validation.
- Application-owned transfer lifecycle state including pause/resume/remove.
- CI gates for formatting, vetting, tests, and build.

## Required architecture

The source must preserve this authority chain:

```text
Clients
  -> Swarm Application API
  -> Swarm service/policy layer
  -> replaceable Engine interface
  -> BitTorrent / peer-to-peer engine adapter
  -> network and filesystem
```

A third-party BitTorrent implementation may be introduced only as a bounded engine dependency behind the GoreeCloud-owned interface. Its APIs must not become the application contract.

## Security and privacy requirements

- Treat torrent metadata, magnet links, trackers, peer data, files, API requests, and imported configuration as untrusted input.
- Do not expose the administration API beyond loopback until authentication and authorization are implemented and tested.
- Do not claim anonymity, Privacy Shield protection, Wardveil protection, Everkeep recovery, or other platform acceptance without authoritative evidence.
- Do not automatically execute downloaded content.
- Keep telemetry off unless separately implemented and governed.
- Prevent torrent-controlled paths from escaping authorized storage roots.
- Require an explicit resolved storage root before starting a real transfer engine.

## Platform systems

All eight GoreeCloud Platform Systems must be evaluated: GoreeCloud Manager, Privacy Shield, Wardveil Security, Everkeep, Glaze UI, GoreeCloud Mesh, GoreeCloud Identity, and GoreeCloud Sync. None is currently represented as runtime-complete in this repository.

## Validation

A material capability is not complete merely because an API or UI exists. Implementation claims require code plus applicable tests, security/privacy review, interoperability validation, accessibility validation, runtime evidence, and documentation alignment.
