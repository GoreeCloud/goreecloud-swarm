# GoreeCloud Swarm — Feature Roadmap

The canonical product direction is maintained in `GoreeCloud/Projects/Project Specification — Swarm`. This repository roadmap tracks implementation sequencing and must remain subordinate to that specification.

## Foundation — in progress

- Native repository and service process.
- Engine abstraction.
- Local API and truthful health state.
- Magnet validation.
- Storage path safety.
- Transfer domain model.
- CI and repository controls.

## V1 — next implementation priorities

1. Bencode and `.torrent` metainfo parser with strict limits and fuzz-test targets.
2. Select and document the initial mature BitTorrent engine dependency and licensing/provenance.
3. Implement the Swarm Engine Adapter without exposing third-party APIs to clients.
4. Add real magnet and `.torrent` ingestion, pause/resume/remove, and capability negotiation.
5. Add durable state/resume persistence with migrations and crash recovery.
6. Implement the Storage Manager and safe file allocation/relocation behavior.
7. Implement interface binding, Network Lock, proxy policy, bandwidth limits, and network diagnostics.
8. Add files, peers, trackers, and pieces read models through `/api/v1`.
9. Publish and contract-test OpenAPI for supported API resources.
10. Begin the native desktop Glaze UI client after service behavior is sufficiently stable.

## V2 — planned

- automation rules;
- RSS/Atom workflows;
- watch folders;
- storage and bandwidth profiles;
- authenticated remote API and session management;
- mature Web UI;
- mobile remote client;
- notifications and advanced diagnostics.

## V3 — planned

- approved Privacy Shield runtime integration;
- approved Wardveil evidence integration;
- approved Everkeep preservation/recovery integration;
- GoreeCloud Manager workflows;
- GoreeCloud Mesh event/capability integration;
- GoreeCloud Identity authentication integration;
- GoreeCloud Sync integration where approved and applicable;
- browser/media handoff and additional engine/protocol evaluation.

## Documentation synchronization

A dedicated central `FEATURE-ROADMAP.docx` for Swarm has not yet been verified during this foundation pass. Until it is created or located and reconciled, this repository remains development/nonconformant rather than claiming full documentation conformance.
