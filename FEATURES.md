# Swarm Features

## Implemented — development foundation

| Capability | State | Evidence |
| --- | --- | --- |
| Native service process | Implemented, development | `cmd/swarmd` |
| Local-only management listener | Implemented, development | `cmd/swarmd/main.go`, tests |
| Versioned HTTP API shell | Implemented, development | `internal/api` |
| Health and capability reporting | Implemented, development | `/api/v1/health`, `/api/v1/capabilities` |
| Engine abstraction | Implemented, development | `internal/engine` |
| Magnet URI parser | Implemented, development | `internal/protocol`, tests |
| Storage path guard | Implemented, development | `internal/storage`, tests |
| Transfer state model | Implemented, development | `internal/core` |
| Add/list/pause/resume/remove lifecycle API | Implemented foundation | `internal/core`, `internal/api`, tests |
| CI format/vet/test/build gate | Implemented | `.github/workflows/ci.yml` |

## Explicitly not implemented yet

- BitTorrent peer connections or piece transfer.
- `.torrent` metainfo ingestion.
- Tracker, DHT, PEX, LPD, web-seed, NAT traversal, or peer-encryption execution.
- Persistence, resume data, or state database.
- Network Lock or proxy routing.
- Storage profiles or relocation workflows.
- Authentication, authorization, remote sessions, or externally reachable Web UI.
- Automation, RSS/Atom, watch folders, or post-completion hooks.
- Native desktop/mobile clients or Glaze UI surfaces.
- Runtime GoreeCloud Manager, Privacy Shield, Wardveil, Everkeep, Mesh, Identity, or Sync integrations.

The canonical Drive specification remains the authority for planned capability scope.
