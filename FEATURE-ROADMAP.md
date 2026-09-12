# GoreeCloud Swarm — Feature Roadmap

The canonical product direction is maintained in `GoreeCloud/Projects/Project Specification — Swarm`. This repository roadmap tracks implementation sequencing and must remain subordinate to that specification.

## Foundation — in progress

- Native repository and service process.
- Engine abstraction.
- Local API and truthful health state.
- Magnet validation.
- Bounded bencode parsing and v1/v2/hybrid metainfo inspection.
- Storage path safety.
- Transfer domain model.
- CI and repository controls.

## Engine adapter — development / validation pending

- Initial dependency selected: `github.com/anacrolix/torrent` v1.61.0.
- Provenance and architectural boundary documented.
- Adapter is opt-in behind `anacrolix_engine`.
- Source exists for magnet add, pause, resume, non-destructive remove, and shutdown.
- Default builds remain on the unavailable engine.
- Hosted tagged compilation/runtime validation is still required before the adapter is accepted.
- Libtorrent-rasterbar remains an explicitly retained alternative for later out-of-process engine evaluation.

## V1 — next implementation priorities

1. Complete exact-head tagged engine compilation and controlled magnet-transfer validation.
2. Add dependency checksum/lock evidence, exact license-file verification, and SBOM-ready provenance.
3. Add fuzz targets and broader compatibility fixtures for bencode and metainfo parsing.
4. Add `.torrent` ingestion through the engine adapter while preserving native preflight validation.
5. Add durable state/resume persistence with migrations and crash recovery.
6. Implement the Storage Manager and safe file allocation/deletion/relocation behavior.
7. Implement interface binding, Network Lock, proxy policy, bandwidth limits, and network diagnostics.
8. Add files, peers, trackers, and pieces read models through `/api/v1`.
9. Move protocol execution behind a process-isolated engine boundary and add sandboxing where practical.
10. Publish and contract-test OpenAPI for supported API resources.
11. Begin the native desktop Glaze UI client after service behavior is sufficiently stable.

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

A dedicated central `FEATURE-ROADMAP.docx` for Swarm has not yet been verified during this development pass. Until it is created or located and reconciled, this repository remains development/nonconformant rather than claiming full documentation conformance.
