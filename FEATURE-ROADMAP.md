# Swarm Feature Roadmap

## Status

This is the repository-coupled implementation roadmap. The canonical Swarm project specification in `GoreeCloud/Projects` controls product scope. A dedicated central `FEATURE-ROADMAP.docx` counterpart has not yet been verified for Swarm, so central-roadmap synchronization remains an open documentation task rather than being falsely claimed complete.

## Foundation — in progress

- Native repository and service boundaries — implemented in development.
- Local-only API and truthful health/capability reporting — implemented in development.
- Magnet validation and storage-boundary primitives — implemented in development.
- CI validation — implemented.
- Repository governance controls — implemented in the initial foundation.
- Native bounded bencode/metainfo inspection — implemented in development.
- Versioned isolated engine-sidecar client/handshake — implemented in development.

## V1 — next implementation priorities

1. Build the first libtorrent-rasterbar sidecar against the versioned Swarm engine protocol and pin verified provenance/licensing.
2. Add real magnet and `.torrent` ingestion, pause/resume/remove, and runtime capability negotiation.
3. Add engine state read models and lifecycle/error translation without exposing libtorrent APIs.
4. Add durable state/resume persistence with migrations and crash recovery.
5. Implement the Storage Manager and safe file allocation/relocation behavior.
6. Implement interface binding, Network Lock, proxy policy, bandwidth limits, and network diagnostics.
7. Add files, peers, trackers, and pieces read models through `/api/v1`.
8. Publish and contract-test OpenAPI for supported API resources.
9. Add fuzz/property tests for bencode, metainfo, magnet, path, and sidecar-protocol boundaries.
10. Begin the native desktop Glaze UI client after service behavior is sufficiently stable.

## V2 — planned

- Advanced automation rules.
- Watch folders.
- RSS/Atom workflows.
- Storage and bandwidth profiles.
- Authenticated remote API and Web UI.
- Remote session/device management.
- Mobile remote client.
- Notifications and advanced diagnostics.

## V3 — planned

- Approved runtime integrations with GoreeCloud Manager, Privacy Shield, Wardveil Security, Everkeep, GoreeCloud Mesh, GoreeCloud Identity, GoreeCloud Sync, and Glaze UI contracts.
- Browser and media handoff integrations.
- Additional peer-to-peer engines/protocols through the engine abstraction when justified.

## Release gate

No milestone is considered complete until implementation and its required tests/evidence exist. Planned entries must not be converted to implemented status based on documentation alone.
