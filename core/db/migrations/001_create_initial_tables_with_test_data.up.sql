CREATE TABLE IF NOT EXISTS identity(
    id serial PRIMARY KEY,
    first_name varchar(255) not null,
    last_name varchar(255) not null,
    email varchar(255) not null
);

INSERT INTO identity(first_name, last_name, email) values ('Test', 'Person', 'test.person@example.org');
INSERT INTO identity(first_name, last_name, email) values ('Simple', 'Identity', 'simple.identity@example.org');
INSERT INTO identity(first_name, last_name, email) values ('New', 'Name', 'new.name@example.org');

CREATE TABLE IF NOT EXISTS account(
    id serial PRIMARY KEY,
    system_id varchar(255) not null,
    identity_id int REFERENCES identity(id)
);

INSERT INTO account(system_id, identity_id) values ('keycloak', (select id from identity where email = 'test.person@example.org'));
INSERT INTO account(system_id, identity_id) values ('keycloak', (select id from identity where email = 'simple.identity@example.org'));
INSERT INTO account(system_id, identity_id) values ('keycloak', (select id from identity where email = 'new.name@example.org'));
