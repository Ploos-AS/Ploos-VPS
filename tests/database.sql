BEGIN;
DO $$
DECLARE
    t1 uuid;
    t2 uuid;
    s uuid;
    c uuid;
    n uuid;
    p uuid;
    i uuid;
BEGIN
    INSERT INTO tenants(name) VALUES ('test one') RETURNING id INTO t1;
    INSERT INTO tenants(name) VALUES ('test two') RETURNING id INTO t2;
    INSERT INTO sites(code,country_code) VALUES ('TEST','NO') RETURNING id INTO s;
    INSERT INTO clusters(site_id,name,provider) VALUES (s,'lab','fake') RETURNING id INTO c;
    INSERT INTO nodes(cluster_id,provider_ref) VALUES (c,'node') RETURNING id INTO n;
    INSERT INTO plans(code,vcpu,ram_mib,disk_gib) VALUES ('test',1,512,10) RETURNING id INTO p;
    INSERT INTO instances(tenant_id,node_id,plan_id,name) VALUES (t1,n,p,'vm') RETURNING id INTO i;
    INSERT INTO operations(tenant_id,instance_id,idempotency_key,request_hash,kind) VALUES (t1,i,'key','hash','create');
    BEGIN
        INSERT INTO operations(tenant_id,instance_id,idempotency_key,request_hash,kind) VALUES (t2,i,'key','hash','create');
        RAISE EXCEPTION 'cross-tenant operation accepted';
    EXCEPTION WHEN foreign_key_violation THEN NULL;
    END;
    BEGIN
        INSERT INTO operations(tenant_id,instance_id,idempotency_key,request_hash,kind) VALUES (t1,i,'key','hash','create');
        RAISE EXCEPTION 'duplicate idempotency key accepted';
    EXCEPTION WHEN unique_violation THEN NULL;
    END;
    BEGIN
        INSERT INTO plans(code,vcpu,ram_mib,disk_gib) VALUES ('invalid',0,512,10);
        RAISE EXCEPTION 'invalid plan accepted';
    EXCEPTION WHEN check_violation THEN NULL;
    END;
END $$;
ROLLBACK;
