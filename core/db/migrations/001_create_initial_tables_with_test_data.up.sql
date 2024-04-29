CREATE TABLE IF NOT EXISTS identity(
    id serial PRIMARY KEY,
    first_name varchar(255) NOT NULL,
    last_name varchar(255) NOT NULL,
    email varchar(255) NOT NULL
);

INSERT INTO identity(first_name, last_name, email) values ('Test', 'Person', 'test.person@example.org');
INSERT INTO identity(first_name, last_name, email) values ('Simple', 'Identity', 'simple.identity@example.org');
INSERT INTO identity(first_name, last_name, email) values ('New', 'Name', 'new.name@example.org');

CREATE TABLE IF NOT EXISTS account(
    id serial PRIMARY KEY,
    username varchar(255) NOT NULL,
    system_id varchar(255) NOT NULL,
    identity_id int REFERENCES identity(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    provisioned_at timestamptz,
    committed_at timestamptz
);

INSERT INTO account(system_id, username, identity_id) values (
    'keycloak_private',
    'testperson',
    (select id from identity where email = 'test.person@example.org'));
INSERT INTO account(system_id, username, identity_id) values (
    'keycloak_private',
    'simpleidentity',
    (select id from identity where email = 'simple.identity@example.org'));
INSERT INTO account(system_id, username, identity_id) values (
    'keycloak_private',
    'newname',
    (select id from identity where email = 'new.name@example.org'));
