# GoreeCloud Swarm — Features

## Current implemented foundation

| Capability | State | Evidence |
| --- | --- | --- |
| Native GoreeCloud service repository | Implemented | Repository history and README |
| Independent `swarmd` service process | Implemented | `cmd/swarmd` |
| Local-only management listener | Implemented | Loopback enforcement in `cmd/swarmd` |
| Versioned API shell | Implemented | `internal/api` |
| Health and engine capability reporting | Implemented | `/api/v1/health`, `/api/v1/capabilities` |
| Replaceable transfer-engine boundary | Implemented | `internal/engine` |
| Magnet URI parser | Implemented foundation | `internal/protocol` |
| Storage path traversal protection | Implemented foundation | `internal/storage` |
| Transfer lifecycle domain model | Implemented foundation | `internal/core` |
| Automated format/vet/test/build validation | Implemented | `.github/workflows/ci.yml` |

## Not yet implemented

The following remain planned and must not be represented as available:

- actual peer-to-peer transfer;
- `.torrent` metainfo/bencode parsing;
- tracker communication;
- DHT, PEX, Local Peer Discovery, and web seeds;
- v1/v2/hybrid torrent operation through a real engine;
- persistent state and resume recovery;
- storage profiles and relocation;
- Network Lock, proxy routing, NAT traversal, and bandwidth profiles;
- file, peer, tracker, and piece inspection backed by live engine state;
- remote authentication, authorization, sessions, and external Web UI;
- automation, RSS/Atom, watch folders, notifications, and CLI management;
- desktop/mobile clients and Glaze UI implementation;
- runtime GoreeCloud Manager, Privacy Shield, Wardveil, Everkeep, Mesh, Identity, and Sync integration;
- production release, deployment, recovery, and Stable qualification.

Feature status must continue to track implementation evidence rather than roadmap intent.
