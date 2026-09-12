# Swarm Architecture — Initial Native Foundation

## Status

Development / pre-alpha. This document describes code that exists in this repository and clearly separates it from planned or unaccepted work.

## Current process boundary

Default build:

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

Opt-in development build (`-tags anacrolix_engine`):

```text
HTTP client
    |
    v
Swarm API
    |
    v
Transfer Manager
    |
    v
Swarm-owned Engine interface
    |
    v
Anacrolix adapter
    |
    v
anacrolix/torrent v1.61.0
    |
    v
network + configured download root
```

The default build intentionally refuses transfer operations. The tagged adapter is development code and is not yet an accepted production engine boundary.

## Package ownership

- `cmd/swarmd`: service entry point, engine selection, and process lifecycle.
- `internal/api`: versioned HTTP surface and response minimization.
- `internal/core`: application-owned transfer lifecycle and authoritative transfer state.
- `internal/engine`: stable internal engine contract, capability reporting, engine selection, and bounded dependency adapters.
- `internal/protocol`: protocol-input parsing that treats remote/user supplied values as untrusted.
- `internal/storage`: path-boundary and storage-safety primitives.

## Authority boundaries

The transfer manager owns application lifecycle state. A BitTorrent dependency owns protocol execution state, but it does not own authorization, API policy, storage policy, privacy decisions, presentation semantics, or GoreeCloud Platform System acceptance.

The HTTP API exposes only accepted Swarm application state. It must not convert requested behavior into a successful, private, protected, or recoverable state without evidence from the responsible subsystem.

## Local API posture

The service binds to loopback by default. This is not a remote-management implementation. Authentication and authorization are mandatory before any externally reachable management profile can be considered supported.

## Initial engine dependency

The first development adapter targets `github.com/anacrolix/torrent` v1.61.0 because it is a Go library and supports major BitTorrent capabilities without adding a C++ ABI to the first service implementation.

The selection is bounded and reversible. See `docs/ENGINE-DEPENDENCY.md`.

The adapter is opt-in rather than compiled into ordinary Swarm builds. This prevents development protocol code from silently enabling peer-to-peer networking or filesystem writes.

## Process isolation

The current tagged adapter is in-process. That is a transitional development state, not the final security architecture.

The preferred V1 architecture remains:

```text
Swarm UI / clients
        |
        v
swarmd service + policy boundary
        |
        v
local authenticated engine IPC
        |
        v
isolated Swarm engine host
        |
        v
bounded BitTorrent dependency
```

Process isolation should eventually allow separate filesystem, network, process-spawn, device, configuration, and secret access restrictions where the target operating system supports them.

## Adapter acceptance requirements

Before an engine adapter is accepted, validation should cover at least:

- exact-head tagged compilation;
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
- storage-boundary behavior;
- error translation;
- recovery after service restart;
- interoperability with independent clients;
- license/provenance and dependency security review.

## Integral Platform Systems

The repository architecture reserves explicit boundaries for the applicable GoreeCloud Platform Systems, but none are represented as implemented yet:

- GoreeCloud Manager — planned.
- Privacy Shield — planned.
- Wardveil Security — planned.
- Everkeep — planned.
- Glaze UI — planned.
- GoreeCloud Mesh — planned.
- GoreeCloud Identity — planned.
- GoreeCloud Sync — planned.

Runtime integration must be added only with the authoritative current contract for each system and must be backed by tests/evidence.
