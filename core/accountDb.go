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
		return nil, fmt.Errorf("unable to get identity %s accounts, %w", identityId, err)
	}

	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Account])

	if err != nil {
		return nil, fmt.Errorf("unable to scan row: %w", err)
	}

	return accounts, nil
}

func (pg *postgres) insertAccount(ctx context.Context, account CreateAccount) {
	query := "insert into account (system_id, username, identity_id, enabled_at, disabled_at, deleted_at)" +
		" values (@systemId, @username, @identityId, @enabledAt, @disabledAt, @deletedAt)"
	args := pgx.NamedArgs{
		"systemId":   account.systemId,
		"username":   account.username,
		"identityId": account.identityId,
		"enabledAt":  account.enabledAt,
		"disabledAt": account.disabledAt,
		"deletedAt":  account.deletedAt,
	}

	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		log.Printf("Unable to insert row: %s", err)
	}
}

func (pg *postgres) getAccountsForProvisioning(ctx context.Context) ([]AccountProvision, error) {
	query := "select identity.id as identity_id, first_name, last_name, email, account.id as account_id, username, system_id" +
		" from account join identity on identity.id = identity_id where provisioned_at is null and committed_at is null"

	rows, err := pg.db.Query(ctx, query)

	if err != nil {
		log.Printf("Error while getting accounts for provisioning, %s", err)
	}

	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[AccountProvision])

	if err != nil {
		return nil, fmt.Errorf("error while collecting rows to accounts: %w", err)
	}

	return accounts, nil
}

func (pg *postgres) getAccountsForRetryProvisioning(ctx context.Context) ([]AccountProvision, error) {
	query := "select identity.id as identity_id, first_name, last_name, email, account.id as account_id, username, system_id" +
		" from account join identity on identity.id = identity_id where provisioned_at < now() - interval '1 minutes' and committed_at is null"

	rows, err := pg.db.Query(ctx, query)

	if err != nil {
		log.Printf("Error while getting accounts for retry provisioning, %s", err)
	}

	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[AccountProvision])

	if err != nil {
		return nil, fmt.Errorf("unable to scan row: %w", err)
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
