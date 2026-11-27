BEGIN;
CREATE TABLE tenants (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL CHECK (length(trim(name)) > 0),
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE sites (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code text NOT NULL UNIQUE,
    country_code text NOT NULL CHECK (country_code ~ '^[A-Z]{2}$'),
    enabled boolean NOT NULL DEFAULT false
);
CREATE TABLE clusters (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    site_id uuid NOT NULL REFERENCES sites(id),
    name text NOT NULL,
    provider text NOT NULL,
    UNIQUE (site_id, name)
);
CREATE TABLE nodes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id uuid NOT NULL REFERENCES clusters(id),
    provider_ref text NOT NULL,
    enabled boolean NOT NULL DEFAULT false,
    UNIQUE (cluster_id, provider_ref)
);
CREATE TABLE plans (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code text NOT NULL UNIQUE,
    vcpu integer NOT NULL CHECK (vcpu > 0),
    ram_mib integer NOT NULL CHECK (ram_mib > 0),
    disk_gib integer NOT NULL CHECK (disk_gib > 0),
    enabled boolean NOT NULL DEFAULT false
);
CREATE TABLE instances (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL REFERENCES tenants(id),
    node_id uuid NOT NULL REFERENCES nodes(id),
    plan_id uuid NOT NULL REFERENCES plans(id),
    name text NOT NULL,
    state text NOT NULL DEFAULT 'pending' CHECK (state IN ('pending','provisioning','running','stopped','deleting','deleted','error')),
    provider_ref text,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, id),
    UNIQUE (node_id, provider_ref)
);
CREATE INDEX instances_tenant_idx ON instances(tenant_id);
CREATE TABLE operations (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid NOT NULL REFERENCES tenants(id),
    instance_id uuid NOT NULL,
    idempotency_key text NOT NULL CHECK (length(idempotency_key) BETWEEN 1 AND 255),
    request_hash text NOT NULL,
    kind text NOT NULL CHECK (kind IN ('create','start','stop','delete')),
    state text NOT NULL DEFAULT 'queued' CHECK (state IN ('queued','running','succeeded','failed')),
    attempts integer NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (tenant_id, instance_id) REFERENCES instances(tenant_id, id),
    UNIQUE (tenant_id, idempotency_key)
);
CREATE INDEX operations_queue_idx ON operations(created_at, id) WHERE state = 'queued';
COMMIT;
