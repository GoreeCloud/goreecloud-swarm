# Swarm Architecture — Initial Native Foundation

## Status

Development / pre-alpha. This document describes code that exists in this repository and clearly separates it from planned work.

## Current process boundary

```text
HTTP client
    |
    v
Swarm API (internal/api)
    |
    v
Transfer Manager (internal/core)
    |
    v
Engine interface (internal/engine)
    |
    v
Unavailable engine placeholder
```

The current engine placeholder intentionally refuses transfer operations. This is a truthful development boundary: the service and API exist, but peer-to-peer transfer does not yet exist.

## Package ownership

- `cmd/swarmd`: service entry point and process lifecycle.
- `internal/api`: versioned HTTP surface and response minimization.
- `internal/core`: application-owned transfer lifecycle and authoritative transfer state.
- `internal/engine`: stable internal engine contract and engine capability reporting.
- `internal/protocol`: protocol-input parsing that must treat remote/user supplied values as untrusted.
- `internal/storage`: path-boundary and storage-safety primitives.

## Authority boundaries

The transfer manager owns application lifecycle state. A future BitTorrent engine owns protocol execution state, but it does not own authorization, API policy, storage policy, privacy decisions, or presentation semantics.

The HTTP API exposes only accepted Swarm application state. It must not convert requested behavior into a successful or protected state without evidence from the responsible subsystem.

## Local API posture

The service binds to loopback by default. This is not a remote-management implementation. Authentication and authorization are mandatory before any externally reachable management profile can be considered supported.

## Transfer engine integration

The first real BitTorrent engine must implement `internal/engine.Engine`. A mature BitTorrent implementation may be used as a bounded dependency. The adapter must translate external engine state into Swarm-owned capability and transfer concepts rather than leaking third-party APIs through the rest of the application.

Before an engine adapter is accepted, it should have tests for at least:

- startup and shutdown;
- magnet ingestion;
- `.torrent` ingestion;
- pause/resume;
- file selection and priority;
- piece verification;
- rate limits;
- tracker state;
- peer state;
- DHT capability reporting;
- error translation;
- recovery after service restart.

## Integral Platform Systems

The repository architecture reserves explicit boundaries for the applicable GoreeCloud platform systems, but none are represented as implemented yet:

- GoreeCloud Manager — planned.
- Privacy Shield — planned.
- Wardveil Security — planned.
- Everkeep — planned.
- Glaze UI — planned.
- GoreeCloud Mesh — planned.
- GoreeCloud Identity — planned.

Runtime integration must be added only with the authoritative current contract for each system and must be backed by tests/evidence.
