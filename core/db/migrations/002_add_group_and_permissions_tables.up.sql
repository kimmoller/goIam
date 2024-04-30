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
    enabled_at timestamptz NOT NULL,
    disabled_at timestamptz,
    deleted_at timestamptz
);
