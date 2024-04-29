CREATE TABLE IF NOT EXISTS permission(
    id serial PRIMARY KEY,
    system_id varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS permission_group(
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS permission_to_group(
    id serial PRIMARY KEY,
    permission_id int REFERENCES permission(id),
    group_id int REFERENCES permission_group(id)
);

CREATE TABLE IF NOT EXISTS group_membership(
    id serial PRIMARY KEY,
    group_id int REFERENCES permission_group(id),
    identity_id int REFERENCES identity(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    provisioned_at timestamptz,
    committed_at timestamptz,
    enabled_at timestamptz NOT NULL,
    enable_provisioned_at timestamptz,
    enable_committed_at timestamptz,
    disabled_at timestamptz,
    disable_provisioned_at timestamptz,
    disable_committed_at timestamptz,
    deleted_at timestamptz,
    delete_provisioned_at timestamptz,
    delete_committed_at timestamptz
);
