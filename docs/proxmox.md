# Proxmox provider

The first adapter lists nodes through GET /api2/json/nodes using a scoped API token. Configure a trusted HTTPS certificate (or install the private CA in the service trust store). TLS verification remains enabled. Use read-only permissions appropriate to inventory; do not use root credentials. No live endpoint or credentials are configured. Tests use a local TLS test server. VM lifecycle, task polling, durable reconciliation, fake provider and database integration are future work. Do not expose this adapter directly to customers.
