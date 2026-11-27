# Ploos VPS — Product strategy and architecture

Status: proposed MVP architecture. No production deployment or commercial commitments.

## Product
Build an internationally usable, self-hostable VPS control plane with English and Norwegian localization, initially operated by Ploos AS. Open-source distribution is intended; license selection requires explicit approval. Support owned and leased infrastructure, with NO-SIR1 as the first planned owned location. Avoid hard-coded site names or provider dependencies.

## Initial product hypotheses
Offer KVM Linux VPS plans with IPv6 by default, optional shared NAT IPv4 or dedicated IPv4 where available. Clearly disclose inbound NAT port limitations. Start with Debian and Ubuntu cloud images. Illustrative monthly NOK prices excluding VAT, not approved offers:

| Plan | vCPU | RAM | SSD | NOK/month |
|---|---:|---:|---:|---:|
| Nano | 1 | 512 MiB | 10 GB | 19 |
| Micro | 1 | 1 GiB | 20 GB | 39 |
| Small | 2 | 2 GiB | 40 GB | 69 |
| Medium | 2 | 4 GiB | 80 GB | 119 |
| Large | 4 | 8 GiB | 160 GB | 229 |
| XL | 8 | 16 GiB | 320 GB | 449 |

Validate margins, storage capacity, CPU contention, traffic, address costs, support, payment fees and VAT before publication. No RAM overcommit initially. Consumer prices must include applicable taxes.

## Architecture
Use a Go modular monolith initially, PostgreSQL for durable state and jobs, a versioned REST API with OpenAPI, and a separate worker process. Evaluate SvelteKit for customer and operator portals. Use a provider interface with a fake implementation for tests and Proxmox VE as the first real adapter. Keep hypervisor APIs private; use scoped credentials and secure TLS verification. Redis is optional, not a mandatory initial dependency.

Core entities: tenant, user, membership, region, site, cluster, node, plan, image, instance, volume, network, IP allocation, operation, audit event and usage record. Separate customer identity from infrastructure identity. Enforce tenant scope server-side. Use transactional capacity reservations and unique IP allocation constraints.

## Provisioning
Accept authenticated requests with idempotency keys; validate entitlement, quotas and capacity, reserve resources, enqueue durable work, create VM from approved image, apply cloud-init and SSH keys, validate outcome, reconcile provider state and record audit events. Model pending, provisioning, running, stopped, deleting, deleted and error states. Retry safely and handle ambiguous provider responses without duplicate VMs. Release reservations only after verified cleanup. Never log secrets or cloud-init credentials.

## Security and operations
MFA for operators, least privilege, tenant isolation, rate limiting, audit trails, backup and restore testing, abuse reporting, explicit support hours and documented incident response. Keep billing, identity and DNS integrations behind interfaces. No direct customer access to hypervisor management. Do not promise HA, backup or data residency beyond verified infrastructure and contracts.

## Roadmap
1. Go repository scaffold, configuration, health endpoints, CI and local development.
2. PostgreSQL schema and migrations, fake provider and API contracts.
3. Proxmox adapter and durable provisioning worker.
4. Capacity reservations, IPAM and tenant authentication.
5. Customer and operator portals, plans and billing integration.
6. Backup, abuse management, observability and deployment.
7. Isolated lab acceptance tests before commercial pilot.

Integrate with Ploos Hosting, Domains and DNS through documented APIs later. Keep development independent of NO-SIR1 construction. No hardware purchases or live infrastructure changes authorized.
