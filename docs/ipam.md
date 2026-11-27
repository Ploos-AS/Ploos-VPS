# IPAM foundation

Initial address selection uses Go net/netip and supports IPv4 and IPv6 prefixes. It reserves the first address and excludes IPv4 broadcast; /31 and /32 pools are not supported by this policy. Explicitly mark gateways and other reserved addresses occupied. The helper is intended for small or sparsely allocated pools, not large dense pools. It does not allocate routed IPv6 prefixes.

## Product modes

- IPv6-only: public IPv6 address or routed prefix, subject to upstream routing and firewall policy.
- NAT IPv4: private IPv4 plus shared public IPv4 with explicit inbound port allocations and documented limitations.
- Dedicated IPv4: exclusive public address, subject to availability and cost.

## Production requirements

Persist pools by site and routing domain, enforce non-overlapping pools within each domain and unique address reservations, use database transactions for concurrent allocation, and reconcile with provider state before releasing addresses. NAT port mappings must be unique by public address, protocol and port; reserve infrastructure ports and enforce tenant isolation. Track allocation ownership, audit events, quarantine and abuse handling. Validate upstream routing, reverse DNS, anti-spoofing and firewall policies. Never assign addresses without authorization from the network operator.

This PR implements only a pure address-selection helper and tests. No persistent IPAM, NAT rules, prefix delegation, routing or live network changes are implemented. Documentation prefixes are used in tests.
