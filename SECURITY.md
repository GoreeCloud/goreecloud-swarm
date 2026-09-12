# Swarm Security

## Status

Development / pre-alpha. Swarm is not production-ready and no Stable security claim is made.

## Current enforced boundaries

- `swarmd` refuses non-loopback HTTP binds while remote authentication is absent.
- HTTP request bodies for transfer creation are bounded.
- Unknown JSON fields are rejected for the current magnet transfer-create request.
- Responses include baseline `nosniff`, no-referrer, and no-store headers.
- Magnet URIs and tracker schemes are validated before they reach an engine.
- `.torrent` metainfo is parsed and classified as v1, v2, or hybrid before Engine-contract dispatch.
- Storage path helpers reject absolute paths and parent traversal outside the authorized download root.
- A failed engine operation is not recorded as an accepted transfer.
- The current default engine placeholder performs no network transfer and truthfully reports itself unavailable.
- A configured engine sidecar must use an absolute executable path, is launched without a shell or PATH lookup, and does not inherit the `swarmd` environment.
- A real sidecar additionally requires an explicit existing absolute download root; Swarm resolves that root before launch and sends it through the versioned handshake.

## Untrusted input model

Future peer messages, tracker responses, downloaded file names/paths, imported configuration, API calls, and integration events must all be treated as untrusted input. Parser and resource-exhaustion testing are release requirements, not optional hardening.

## Remote management

Remote administration is not implemented. Do not expose the current development API through a reverse proxy, public listener, port forward, or similar external path and represent it as a supported remote-management deployment.

## Secrets

Do not commit passwords, API tokens, private keys, recovery codes, authentication cookies, or other reusable secrets. `.env` and `.env.*` are ignored except sanitized example files.

## Vulnerability reporting

Do not publish active credentials or sensitive exploit details in a public issue. Use an authorized private GoreeCloud security reporting channel. A dedicated public Swarm security-contact address has not yet been established in this repository.

## Platform security state

Wardveil Security integration is planned but not currently accepted at runtime. Ordinary Swarm checks or a green CI run must not be represented as a Wardveil verdict.
