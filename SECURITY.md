# GoreeCloud Swarm — Security

## Current status

Swarm is development/pre-alpha software and is not production-ready.

## Security controls implemented in the foundation

- The management service binds to loopback by default and refuses non-loopback addresses until authenticated remote management is implemented.
- HTTP request bodies for transfer creation are bounded.
- Unknown JSON request fields are rejected.
- Responses include baseline `nosniff`, no-referrer, and no-store headers.
- Magnet URIs and tracker schemes are validated before they reach an engine.
- Storage path helpers reject absolute paths and parent traversal outside the authorized download root.
- A failed engine operation is not recorded as an accepted transfer.
- The current engine placeholder performs no network transfer and truthfully reports itself unavailable.

## Security boundaries not yet implemented

The current foundation does not provide:

- remote authentication or authorization;
- TLS termination or externally reachable Web UI;
- session/device management;
- rate limiting or brute-force protection for remote access;
- BitTorrent peer/tracker protocol processing;
- torrent metainfo/bencode parsing;
- engine process sandboxing;
- persistent-state encryption or database hardening;
- signed release/update verification;
- Wardveil runtime acceptance or security evidence.

These gaps block production and Stable qualification.

## Secrets

Do not commit passwords, API tokens, private keys, recovery codes, authentication cookies, or other reusable secrets. Local secret-bearing `.env` files are ignored; `.env.example` is intentionally sanitized.

## Dependency security

A future BitTorrent engine or other material dependency must have its provenance, version, license, security implications, update source, and architectural boundary documented before acceptance. Third-party code must remain subordinate to Swarm's native application and policy layers.

## Vulnerability reporting

Do not publish active exploit details, credentials, private user data, or other sensitive security information in a public issue. Use an authorized private GoreeCloud security-reporting channel. A dedicated public Swarm security-reporting address has not yet been established in this repository.

## Wardveil

Wardveil integration is planned but not implemented. Repository labels, documentation, or UI must not claim Wardveil protection unless authoritative runtime evidence has actually been accepted.
