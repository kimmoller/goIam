CREATE TABLE IF NOT EXISTS identity(
    id serial PRIMARY KEY,
    first_name varchar(255) NOT NULL,
    last_name varchar(255) NOT NULL,
    email varchar(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS account(
    id serial PRIMARY KEY,
    username varchar(255) NOT NULL,
    system_id varchar(255) NOT NULL,
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
