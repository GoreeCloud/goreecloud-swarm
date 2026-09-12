# GoreeCloud Swarm — Security

## Current status

Swarm is development/pre-alpha software and is not production-ready.

## Security controls implemented in the foundation

- The management service binds to loopback by default and refuses non-loopback addresses until authenticated remote management is implemented.
- HTTP request bodies are bounded, including a 16 MiB development metainfo-inspection ceiling.
- Unknown JSON transfer-request fields are rejected.
- Responses include baseline `nosniff`, no-referrer, and no-store headers.
- Magnet URIs and tracker schemes are validated before they reach an engine.
- The native bencode decoder bounds total input size, nesting depth, byte-string length, and container entries; duplicate dictionary keys and non-canonical integers are rejected.
- `.torrent` inspection validates essential v1/v2 structure before exposing metadata and computes hashes from the original raw `info` bytes.
- Storage path helpers reject absolute paths and parent traversal outside the authorized download root.
- A failed engine operation is not recorded as an accepted transfer.
- Default builds contain no active BitTorrent engine and truthfully report the engine unavailable.
- The development Anacrolix adapter is opt-in behind a build tag and runtime engine selector.
- The tagged adapter requires an explicit download root and refuses engine-managed data deletion until a native safe-deletion policy exists.

## Development engine security boundary

The current `anacrolix_engine` adapter is a transitional in-process engine boundary. It is useful for establishing the Swarm-owned engine contract, but it does not yet satisfy the preferred long-term process-isolation design.

Enabling the tagged engine may create outbound/inbound BitTorrent traffic and write data beneath the configured download root. It must not be represented as providing:

- Network Lock;
- VPN binding or leak prevention;
- proxy enforcement;
- anonymity;
- Wardveil runtime acceptance;
- Privacy Shield runtime acceptance;
- Everkeep recovery coverage;
- production sandboxing.

## Security boundaries not yet implemented or accepted

The current development state does not provide:

- remote authentication or authorization;
- TLS termination or externally reachable Web UI;
- session/device management;
- rate limiting or brute-force protection for remote access;
- accepted/sandboxed BitTorrent peer and tracker protocol processing;
- metainfo fuzzing and broad real-world compatibility corpus validation;
- engine process sandboxing;
- Network Lock or explicit interface-binding enforcement;
- proxy-routing enforcement;
- persistent-state encryption or database hardening;
- signed release/update verification;
- Wardveil runtime acceptance or security evidence.

These gaps block production and Stable qualification.

## Secrets

Do not commit passwords, API tokens, private keys, recovery codes, authentication cookies, or other reusable secrets. Local secret-bearing `.env` files are ignored; `.env.example` is intentionally sanitized.

## Dependency security

The initial development engine dependency is pinned to `github.com/anacrolix/torrent` v1.61.0 and documented in `docs/ENGINE-DEPENDENCY.md` and `THIRD_PARTY_NOTICES.md`.

Before acceptance or distribution, GoreeCloud must complete:

- exact upstream license-file verification;
- dependency checksum/lock evidence;
- vulnerability and maintenance review;
- transitive dependency review appropriate to release risk;
- controlled tagged-build tests;
- interoperability tests;
- runtime storage/network boundary validation;
- SBOM/release provenance as required.

Third-party code must remain subordinate to Swarm's native application and policy layers.

## Vulnerability reporting

Do not publish active exploit details, credentials, private user data, or other sensitive security information in a public issue. Use an authorized private GoreeCloud security-reporting channel. A dedicated public Swarm security-reporting address has not yet been established in this repository.

## Wardveil

Wardveil integration is planned but not implemented. Repository labels, documentation, or UI must not claim Wardveil protection unless authoritative runtime evidence has actually been accepted.
