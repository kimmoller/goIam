package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func (pg *postgres) getIdentityAccountsFromDb(ctx context.Context, identityId string) ([]Account, error) {
	query := "select * from account where identity_id = @identityId"

	args := pgx.NamedArgs{
		"identityId": identityId,
	}

	rows, err := pg.db.Query(ctx, query, args)

	if err != nil {
		log.Printf("Unable to get identity %s accounts, %s", identityId, err)
		return nil, err
	}

	accounts := []Account{}
	for rows.Next() {
		account := Account{}
		err := rows.Scan(&account.ID, &account.SystemId, &account.IdentityId)

		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (pg *postgres) insertAccount(ctx context.Context, account Account) {
	query := "insert into account (system_id, identity_id) values (@systemId, @identityId)"
	args := pgx.NamedArgs{
		"systemId":   account.SystemId,
		"identityId": account.IdentityId,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		log.Printf("Unable to insert row: %s", err)
	}
}
