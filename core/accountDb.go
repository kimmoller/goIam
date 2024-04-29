package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"gopkg.in/guregu/null.v3"
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

func (pg *postgres) getAccountsForProvisioning(ctx context.Context) ([]Account, error) {
	query := "select * from account where provisioned_at is null and committed_at is null"

	rows, err := pg.db.Query(ctx, query)

	if err != nil {
		log.Printf("Error while getting accounts for provisioning, %s", err)
	}

	accounts := []Account{}
	for rows.Next() {
		account := Account{}
		err := rows.Scan(&account.ID, &account.SystemId, &account.IdentityId, &account.CreatedAt, &account.ProvisionedAt, &account.CommittedAt)

		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (pg *postgres) getAccountsForRetryProvisioning(ctx context.Context) ([]Account, error) {
	query := "select * from account where provisioned_at < now() - interval '1 minutes' and committed_at is null"

	rows, err := pg.db.Query(ctx, query)

	if err != nil {
		log.Printf("Error while getting accounts for retry provisioning, %s", err)
	}

	accounts := []Account{}
	for rows.Next() {
		account := Account{}
		err := rows.Scan(&account.ID, &account.SystemId, &account.IdentityId, &account.CreatedAt, &account.ProvisionedAt, &account.CommittedAt)

		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (pg *postgres) markAccountAsProvisioned(id null.String) {
	log.Printf("Mark account %s as provisioned", fmt.Sprint(id))
	query := "update account set provisioned_at = now() where id=@id"
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := pg.db.Exec(context.Background(), query, args)

	if err != nil {
		log.Printf("Failed to mark account as provisioned: %s", err)
	}
}

func (pg *postgres) markAccountAsCommitted(id null.String) {
	log.Printf("Mark account %s as committed", fmt.Sprint(id))
	query := "update account set committed_at = now() where id=@id"
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := pg.db.Exec(context.Background(), query, args)

	if err != nil {
		log.Printf("Failed to mark account as provisioned: %s", err)
	}
}
