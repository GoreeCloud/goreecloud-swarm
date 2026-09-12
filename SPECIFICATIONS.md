# GoreeCloud Swarm — Repository Specifications

## Authority

The canonical product specification is maintained in Google Drive as **Project Specification — Swarm** under `GoreeCloud/Projects`. This repository copy summarizes implementation-coupled requirements and must not silently redefine the canonical specification.

## Current implementation status

**Lifecycle:** Development / pre-alpha  
**Development model:** Native GoreeCloud repository  
**Production ready:** No

The repository implements the first Swarm service/metainfo foundation and now contains an opt-in development adapter for `anacrolix/torrent` v1.61.0. The adapter is not accepted as production-capable peer transfer until exact-head hosted compilation and controlled runtime/interoperability validation pass.

Implemented foundation:

- independent `swarmd` service process;
- versioned local HTTP API under `/api/v1`;
- truthful health and engine-capability reporting;
- replaceable transfer-engine contract;
- magnet URI parsing and validation;
- bounded native bencode decoding;
- v1/v2/hybrid `.torrent` metainfo inspection and raw-info hashing;
- authorized storage-path guards;
- application-owned transfer lifecycle model;
- loopback-only management listener until remote authentication exists;
- fail-closed engine selection;
- opt-in tagged Anacrolix adapter source for magnet add, pause, resume, non-destructive remove, and shutdown;
- unit/API validation and separate default/tagged CI gates.

## Architectural contract

```text
Clients
  -> Swarm Application API
  -> Swarm service and policy layer
  -> Swarm-owned engine interface
  -> bounded transfer-engine adapter
  -> BitTorrent dependency
  -> network and filesystem
```

The initial development dependency is Anacrolix/torrent v1.61.0. It remains subordinate to the Swarm-owned engine contract. Swarm clients and product-facing packages must not depend directly on third-party engine APIs.

Libtorrent-rasterbar remains an evaluated alternative for later interoperability/performance work and for the planned out-of-process engine boundary. The current dependency choice is not permanent architecture.

## Required behavior

- Requested operations must not be represented as accepted or completed until the responsible subsystem confirms them.
- Untrusted torrent, magnet, tracker, peer, path, API, and imported configuration data must be validated and bounded.
- Metainfo parsing must impose explicit input, nesting, string, and container limits before engine ingestion.
- Local operation must not require a GoreeCloud account.
- Remote administration must remain disabled until authentication, authorization, session, transport, rate-limit, and audit requirements are implemented and verified.
- Downloaded content must never be automatically executed.
- File paths supplied by torrent metadata must remain inside an authorized storage boundary.
- External engine capabilities must be represented conservatively and must not be confused with runtime privacy/security acceptance.
- Engine activation must remain explicit during development; default builds must not silently begin peer-to-peer networking.
- Deleting downloaded data through an external engine must remain blocked until Swarm's native storage policy can verify the deletion boundary.

## GoreeCloud Platform Systems

Swarm must evaluate and eventually integrate all currently required GoreeCloud Platform Systems where applicable:

- GoreeCloud Manager
- Privacy Shield
- Wardveil Security
- Everkeep
- Glaze UI
- GoreeCloud Mesh
- GoreeCloud Identity
- GoreeCloud Sync

None of these runtime integrations is complete in the current development state. Repository metadata and UI must not imply otherwise.

## V1 implementation target

The production-oriented V1 target remains broader than the current branch and includes a validated BitTorrent engine adapter, `.torrent` and magnet transfer operation, durable resume state, storage management, interface binding and Network Lock, proxy support, transfer inspection, process isolation, secure updates, local API, basic Web UI, accessibility, and interoperability/security/privacy testing.
