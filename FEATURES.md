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
| Bounded bencode decoder | Implemented foundation | `internal/protocol/bencode` |
| v1/v2/hybrid `.torrent` metainfo inspection | Implemented foundation | `internal/protocol/torrent.go` |
| Metainfo inspection API | Implemented foundation | `POST /api/v1/metainfo/inspect` |
| Storage path traversal protection | Implemented foundation | `internal/storage` |
| Transfer lifecycle domain model | Implemented foundation | `internal/core` |
| Opt-in Anacrolix engine adapter source | Development / validation pending | `internal/engine/factory_anacrolix.go` |
| Fail-closed engine selection | Implemented development | `internal/engine/factory*.go` |
| Separate default/tagged engine CI gates | Implemented configuration | `.github/workflows/ci.yml` |

## Engine adapter development state

The opt-in `anacrolix_engine` build currently contains source for:

- initializing `anacrolix/torrent` v1.61.0 with an explicit download root;
- adding magnet transfers;
- starting download-all behavior;
- pausing upload and download data flow;
- resuming upload and download data flow;
- removing a transfer from the engine without deleting downloaded data;
- clean engine shutdown;
- conservative capability reporting.

This code is **not yet accepted implementation evidence** for peer transfer. Hosted exact-head compilation, controlled runtime transfer tests, interoperability tests, and dependency/license review remain outstanding.

## Not yet implemented or accepted

The following remain planned/development and must not be represented as production-ready:

- accepted/validated peer-to-peer transfer support in the default build;
- using `.torrent` metainfo to start transfers through the engine;
- persistent state and resume recovery;
- storage profiles, safe engine-owned deletion, and relocation;
- Network Lock, proxy routing, NAT traversal, and bandwidth profiles;
- file, peer, tracker, and piece inspection backed by accepted live engine state;
- remote authentication, authorization, sessions, and external Web UI;
- automation, RSS/Atom, watch folders, notifications, and CLI management;
- desktop/mobile clients and Glaze UI implementation;
- runtime GoreeCloud Manager, Privacy Shield, Wardveil, Everkeep, Mesh, Identity, and Sync integration;
- production release, deployment, recovery, and Stable qualification.

Feature status must continue to track implementation evidence rather than roadmap intent.
