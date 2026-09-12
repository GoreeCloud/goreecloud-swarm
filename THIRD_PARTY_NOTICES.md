# Third-Party Notices

This file records material third-party software selected or used by GoreeCloud Swarm. It is a development provenance record and does not replace the governing upstream license texts.

## anacrolix/torrent

- **Purpose:** Initial bounded BitTorrent protocol engine dependency behind Swarm's `internal/engine` interface.
- **Module:** `github.com/anacrolix/torrent`
- **Pinned development version:** `v1.61.0`
- **Upstream:** https://github.com/anacrolix/torrent
- **Current integration state:** Development / opt-in build tag `anacrolix_engine`.
- **Product-defining role:** No. It supplies BitTorrent protocol execution behind a GoreeCloud-owned adapter; Swarm owns application behavior, API, policy, storage governance, privacy/security boundaries, lifecycle state, and user experience.
- **License status:** Public Go package metadata currently reports MPL-2.0 and MIT. Exact upstream license files and file-level applicability must be verified and preserved before distribution or Stable qualification.

Swarm must not copy upstream branding or present the dependency as the GoreeCloud product. Dependency upgrades require provenance, compatibility, security, license, and exact-head validation.
