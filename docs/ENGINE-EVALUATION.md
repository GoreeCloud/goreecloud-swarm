# Initial BitTorrent Engine Evaluation

## Status

Development evaluation. This document records the current engine direction; it is not evidence that a production engine is integrated.

## Requirements carried from the canonical Swarm specification

The initial engine must be capable of supporting the V1 interoperability target, including at minimum:

- magnet links;
- BitTorrent v1;
- BitTorrent v2;
- hybrid v1/v2 torrents;
- trackers;
- DHT;
- Peer Exchange;
- IPv4/IPv6;
- selective files and priorities;
- pause/resume/remove;
- piece verification;
- rate/network controls;
- inspectable peer/tracker/transfer state.

The engine must remain subordinate to the GoreeCloud-owned service/API/policy architecture.

## libtorrent-rasterbar — leading candidate

As verified during the September 12, 2026 implementation pass, the official libtorrent project describes itself as a feature-complete C++ BitTorrent implementation focused on efficiency and scalability and documents BitTorrent v2 support in the 2.x line. The project is distributed under a BSD license. The latest upstream non-prerelease release observed during this pass is **v2.1.1**, published **August 10, 2026**.

Why it currently leads:

- broad mature BitTorrent feature coverage;
- explicit BitTorrent v2 support;
- established use by major BitTorrent clients;
- suitable library architecture rather than a requirement to fork a complete client;
- permissive licensing compatible with use as a bounded dependency;
- fits Swarm's specification, which already named libtorrent-rasterbar as a candidate.

Swarm will not directly expose libtorrent APIs. The planned integration is a dedicated engine sidecar implementing `docs/ENGINE-PROTOCOL.md`.

## anacrolix/torrent — evaluated, not selected for the first engine

`anacrolix/torrent` is attractive because it is a Go library and would simplify build integration. However, its public BEP 52 / BitTorrent v2 support issue remains open as of the current verification pass. Because Swarm's V1 explicitly targets v1, v2, and hybrid torrents, adopting it as the sole initial engine would require weakening an existing interoperability requirement or adding substantial missing protocol work.

It may remain useful for future evaluation, testing, or a secondary engine if its capability set changes.

## Decision for the next implementation slice

Proceed with the isolated sidecar architecture and use libtorrent-rasterbar **v2.1.1** as the initial evaluation pin. This is not yet an accepted production dependency. Do not mark it integrated/accepted until the following are verified in an actual Swarm sidecar build:

1. exact version and source provenance;
2. license/notice requirements;
3. v1/v2/hybrid behavior;
4. DHT/PEX/tracker/IPv6 behavior;
5. pause/resume/remove semantics;
6. file selection and priority APIs;
7. piece/hash verification;
8. proxy/interface binding and Network Lock feasibility;
9. runtime capability discovery;
10. crash/restart and persistence behavior;
11. security and sandboxing implications;
12. reproducible build and CI support on target platforms.

Until those checks pass, `goreecloud.platform.yaml` must continue to report the transfer engine as incomplete and Swarm as nonconformant/pre-alpha.
