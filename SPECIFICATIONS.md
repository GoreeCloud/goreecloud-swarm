# GoreeCloud Swarm — Repository Specifications

## Authority

The canonical product specification is maintained in Google Drive as **Project Specification — Swarm** under `GoreeCloud/Projects`. This repository copy summarizes implementation-coupled requirements and must not silently redefine the canonical specification.

## Current implementation status

**Lifecycle:** Development / pre-alpha  
**Development model:** Native GoreeCloud repository  
**Production ready:** No

The current repository implements the first Swarm service and metainfo-inspection foundation. It does not yet perform BitTorrent transfers.

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
- unit/API validation and CI build gates.

## Architectural contract

```text
Clients
  -> Swarm Application API
  -> Swarm service and policy layer
  -> Swarm-owned engine interface
  -> bounded transfer-engine adapter
  -> network and filesystem
```

A future BitTorrent implementation may be a mature third-party dependency, but it must remain behind the Swarm-owned engine and policy boundaries. Swarm must not become a direct fork of another torrent client.

## Required behavior

- Requested operations must not be represented as accepted or completed until the responsible subsystem confirms them.
- Untrusted torrent, magnet, tracker, peer, path, API, and imported configuration data must be validated and bounded.
- Metainfo parsing must impose explicit input, nesting, string, and container limits before engine ingestion.
- Local operation must not require a GoreeCloud account.
- Remote administration must remain disabled until authentication, authorization, session, transport, rate-limit, and audit requirements are implemented and verified.
- Downloaded content must never be automatically executed.
- File paths supplied by torrent metadata must remain inside an authorized storage boundary.
- External engine capabilities must be discovered and represented truthfully rather than assumed.

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

None of these runtime integrations is complete in this foundation. Repository metadata and UI must not imply otherwise.

## V1 implementation target

The production-oriented V1 target remains broader than this foundation and includes a real BitTorrent engine adapter, `.torrent` and magnet transfer operation, storage management, interface binding and Network Lock, proxy support, transfer inspection, persistence/recovery, secure updates, local API, basic Web UI, accessibility, and interoperability/security/privacy testing.
