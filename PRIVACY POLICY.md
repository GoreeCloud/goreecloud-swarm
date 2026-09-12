# GoreeCloud Swarm — Privacy Policy

## Current development state

GoreeCloud Swarm is in development/pre-alpha. This policy describes the behavior implemented by the current repository foundation; planned BitTorrent and GoreeCloud platform capabilities are identified separately.

## Current data behavior

The current service:

- binds its management API to loopback only by default and refuses non-loopback binds;
- contains no advertising framework;
- contains no behavioral analytics or remote telemetry integration;
- does not require a GoreeCloud account;
- does not establish BitTorrent peer, tracker, DHT, PEX, or LPD connections;
- does not persist a transfer library or downloaded content;
- does not provide remote-management sessions;
- emits local process logs for service lifecycle and errors.

A magnet URI submitted to the current transfer endpoint is validated, but because no transfer engine is connected the operation is rejected and is not added to the authoritative transfer library.

## Sensitive transfer information

Future Swarm implementations should minimize exposure and retention of torrent names, file names, magnet URIs, tracker URLs, peer addresses, local paths, credentials, tokens, and other sensitive operational details. Ordinary logs and generic GoreeCloud evidence must not become a second permanent record of transfer content.

## BitTorrent privacy reality

BitTorrent is a peer-to-peer protocol. When peer transfer is implemented, participating peers ordinarily require network addressing information sufficient to establish connections. Peer-connection encryption alone will not make participation anonymous. Proxy or VPN protection will apply only to traffic that is actually routed through the configured path.

## Privacy Shield

Privacy Shield runtime integration is planned but not implemented. Swarm must not display or publish a Privacy Shield-approved state until authoritative Privacy Shield evidence has actually been accepted.

## Future changes

This policy must be updated before introducing material telemetry, persistent transfer history, authenticated remote access, account/identity processing, external integrations, or other data-handling behavior not described above.
