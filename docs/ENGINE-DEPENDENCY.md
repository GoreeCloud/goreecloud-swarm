# Swarm Initial BitTorrent Engine Dependency

## Status

**Decision state:** Development selection for the first bounded engine adapter.  
**Selected dependency:** `github.com/anacrolix/torrent`  
**Pinned development version:** `v1.61.0`  
**Activation:** opt-in build tag `anacrolix_engine`; runtime selection `SWARM_ENGINE=anacrolix`.

This decision does not make Anacrolix the permanent Swarm engine. Swarm owns the application API, policy, transfer state, storage policy, privacy/security boundaries, and engine interface so the protocol implementation can be replaced or supplemented later.

## Why this dependency is selected first

The current Swarm service is written in Go. Anacrolix/torrent is designed for use as a Go library and publicly documents support for major BitTorrent capabilities including protocol encryption, DHT, PEX, uTP, and v1/v2/hybrid metainfo workflows.

Using a Go library for the first adapter avoids introducing a C++ ABI and cross-language binding boundary into the initial Linux service while Swarm's own engine contract is still being established.

The dependency remains subordinate to `internal/engine`. Swarm callers must not import or depend on Anacrolix APIs directly.

## Alternatives evaluated

### libtorrent-rasterbar

Libtorrent remains a serious alternative. Its official project documentation describes it as a feature-complete C++ BitTorrent implementation focused on efficiency and scalability, and its 2.x API supports BitTorrent v2 concepts.

Reasons it is not the first adapter:

- the current Swarm service is Go-native;
- direct use would introduce C++ toolchain, ABI, packaging, and binding complexity early in the project;
- Swarm has not yet established the process-isolated engine boundary that would make a language-neutral native engine host preferable;
- selecting it now would increase build and packaging surface before core lifecycle, persistence, Network Lock, and storage policy are mature.

Libtorrent should be re-evaluated when Swarm's out-of-process engine protocol is implemented and during interoperability/performance testing.

## Current adapter scope

The tagged adapter currently implements:

- engine initialization with an explicit download root;
- magnet addition;
- start/download-all behavior;
- pause of upload and download data flow;
- resume of upload and download data flow;
- removal from the engine without deleting downloaded data;
- engine shutdown;
- conservative capability reporting.

It deliberately does **not** implement:

- `.torrent` file ingestion into the engine;
- persistent resume/session state;
- deletion of downloaded files;
- Network Lock;
- proxy routing policy;
- storage relocation;
- file priorities;
- queue policy;
- tracker/peer/piece read models;
- process isolation;
- GoreeCloud Platform System runtime acceptance.

## Fail-closed activation

Normal builds do not include the dependency adapter. `SWARM_ENGINE=anacrolix` is rejected unless the binary was built with:

```bash
go build -tags anacrolix_engine ./cmd/swarmd
```

When the tagged adapter is included, `SWARM_DOWNLOAD_DIR` is mandatory. The adapter creates only that configured root before handing it to the dependency.

This opt-in model prevents a development dependency from silently enabling network traffic or filesystem writes in ordinary Swarm builds.

## Licensing and provenance

Public Go package metadata for the selected module currently reports MPL-2.0 and MIT licensing. Before any packaged or Stable release includes this dependency, GoreeCloud must verify the exact upstream license files and file-level licensing that apply to the pinned revision, preserve required notices, and record the result in release/SBOM evidence.

No dependency license is reclassified or replaced by the GoreeCloud project license.

## Validation requirements

Before this adapter can be treated as accepted implementation rather than proposed development code, the exact PR head must pass:

- default format/vet/test/build validation;
- tagged `anacrolix_engine` compilation and tests;
- a controlled magnet download test using authorized test content;
- pause/resume verification;
- clean shutdown verification;
- storage-boundary verification;
- dependency/license review;
- later interoperability testing against independent BitTorrent clients.

The adapter must remain Development until that evidence exists.
