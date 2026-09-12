# GoreeCloud Swarm — Repository Notes

## Current foundation decisions

- Swarm is being built in a GoreeCloud-created native repository rather than as a fork of another torrent client.
- The initial service implementation is Go and currently uses only the standard library.
- No production BitTorrent engine is connected yet.
- libtorrent-rasterbar v2.1.1 is the current evaluation pin, not an accepted production dependency.
- The engine boundary is intentionally stable so a mature protocol implementation can be introduced as a bounded dependency without defining Swarm's application architecture.
- Engine integration uses an isolated versioned sidecar protocol rather than exposing third-party APIs to Swarm clients.
- A configured sidecar requires an absolute executable path and an explicit existing absolute download root; the root is resolved before launch and sent through the handshake.
- The management API is intentionally loopback-only until authenticated remote administration exists.
- Health remains degraded while no production transfer engine is connected.
- Transfer requests rejected by the engine are not persisted as accepted transfers.
- Magnet and validated `.torrent` ingestion, plus pause/resume/remove lifecycle controls, now reach the Engine contract; peer transfer itself remains unimplemented.

## Next engineering decision

Build and validate the first libtorrent-rasterbar sidecar against the existing protocol and evaluate its real runtime support for v1/v2/hybrid torrents, trackers, DHT, PEX, IPv6, lifecycle controls, storage behavior, network binding/proxy controls, restart/recovery behavior, and reproducible builds.

## Documentation follow-up

A dedicated central Swarm `FEATURE-ROADMAP.docx` under the GoreeCloud Feature Roadmap location has not yet been verified during this implementation pass. Repository conformance must remain truthful until that synchronization requirement is satisfied.
