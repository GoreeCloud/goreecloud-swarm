# GoreeCloud Swarm — Repository Notes

## Current foundation decisions

- Swarm is being built in a GoreeCloud-created native repository rather than as a fork of another torrent client.
- The initial service implementation is Go and currently uses only the standard library.
- No BitTorrent engine has been selected or connected yet.
- The engine boundary is intentionally stable so a mature protocol implementation can be introduced as a bounded dependency without defining Swarm's application architecture.
- The management API is intentionally loopback-only until authenticated remote administration exists.
- Health remains degraded while no transfer engine is connected.
- Transfer requests rejected by the engine are not persisted as accepted transfers.

## Next engineering decision

Evaluate candidate BitTorrent engines against protocol coverage, license compatibility, security history, platform support, API stability, testability, resource controls, and ability to remain behind the Swarm Engine Adapter.

## Documentation follow-up

A dedicated central Swarm `FEATURE-ROADMAP.docx` under the GoreeCloud Feature Roadmap location has not yet been verified during this implementation pass. Repository conformance must remain truthful until that synchronization requirement is satisfied.
