package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

func (pg *postgres) getIdentitiesFromDb(ctx context.Context) ([]Identity, error) {
	query := "select * from identity left join account on identity.id=identity_id"

	rows, err := pg.db.Query(ctx, query)

	if err != nil {
		log.Printf("Unable to get identities, %s", err)
		return nil, err
	}

	identities := []Identity{}
	for rows.Next() {
		identity := Identity{}
		err := rows.Scan(&identity.ID, &identity.FirstName, &identity.LastName, &identity.Email,
			&identity.Account.ID, &identity.Account.Username, &identity.Account.SystemId, &identity.Account.IdentityId,
			&identity.Account.CreatedAt, &identity.Account.ProvisionedAt, &identity.Account.CommittedAt)

		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		identities = append(identities, identity)
	}

	return identities, nil
}

func (pg *postgres) insertIdentity(ctx context.Context, indetity Identity) {
	query := "insert into identity (first_name, last_name, email) values (@firstName, @lastName, @email)"
	args := pgx.NamedArgs{
		"firstName": indetity.FirstName,
		"lastName":  indetity.LastName,
		"email":     indetity.Email,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		log.Printf("Unable to insert row: %s", err)
	}
}
