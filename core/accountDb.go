package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"gopkg.in/guregu/null.v3"
)

func (pg *postgres) insertAccount(ctx context.Context, account CreateAccount) error {
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
		return fmt.Errorf("error while creation identity %s account %s", account.identityId, account.systemId)
	}
	return nil
}

func (pg *postgres) updateAccount(ctx *gin.Context, updateAccount UpdateAccount, reEnable bool) error {
	query := "update account set enabled_at = @enabledAt, disabled_at = @disabledAt, deleted_at = @deletedAt" +
		" where identity_id = @identityId and system_id = @systemId"

	if reEnable {
		query = "update account set enabled_at = @enabledAt, disabled_at = @disabledAt, deleted_at = @deletedAt" +
			", enable_provisioned_at = null, enable_committed_at = null, disable_provisioned_at = null," +
			" disable_committed_at = null where identity_id = @identityId and system_id = @systemId"
	}

	args := pgx.NamedArgs{
		"enabledAt":  updateAccount.enabledAt,
		"disabledAt": updateAccount.disabledAt,
		"deletedAt":  updateAccount.deletedAt,
		"identityId": updateAccount.identityId,
		"systemId":   updateAccount.systemId,
	}

	_, err := pg.db.Exec(ctx, query, args)

	if err != nil {
		return fmt.Errorf("error while updating account %s for identity %s, %s", updateAccount.systemId, updateAccount.identityId, err)
	}

	return nil
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

func (pg *postgres) getEnableAccountsForProvisioning(ctx context.Context) ([]AccountProvision, error) {
	query := "select identity.id as identity_id, first_name, last_name, email, account.id as account_id, username, system_id" +
		" from account join identity on identity.id = identity_id where provisioned_at is not null and committed_at is not null" +
		" and enabled_at < now() and enable_provisioned_at is null and enable_committed_at is null"

	rows, err := pg.db.Query(ctx, query)

	if err != nil {
		log.Printf("Error while getting accounts for enable provisioning, %s", err)
	}

	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[AccountProvision])

	if err != nil {
		return nil, fmt.Errorf("error while collecting rows to accounts: %w", err)
	}

	return accounts, nil
}

func (pg *postgres) getEnableAccountsForRetryProvisioning(ctx context.Context) ([]AccountProvision, error) {
	query := "select identity.id as identity_id, first_name, last_name, email, account.id as account_id, username, system_id" +
		" from account join identity on identity.id = identity_id where provisioned_at is not null and committed_at is not null" +
		" and enabled_at < now() and enable_provisioned_at < now() - interval '1 minutes' and enable_committed_at is null"

	rows, err := pg.db.Query(ctx, query)

	if err != nil {
		log.Printf("Error while getting accounts for enable retry provisioning, %s", err)
	}

	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[AccountProvision])

	if err != nil {
		return nil, fmt.Errorf("unable to scan row: %w", err)
	}

	return accounts, nil
}

func (pg *postgres) getDisableAccountsForProvisioning(ctx context.Context) ([]AccountProvision, error) {
	query := "select identity.id as identity_id, first_name, last_name, email, account.id as account_id, username, system_id" +
		" from account join identity on identity.id = identity_id where enable_provisioned_at is not null and enable_committed_at is not null" +
		" and disabled_at < now() and disable_provisioned_at is null and disable_committed_at is null"

	rows, err := pg.db.Query(ctx, query)

	if err != nil {
		log.Printf("Error while getting accounts for disable provisioning, %s", err)
	}

	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[AccountProvision])

	if err != nil {
		return nil, fmt.Errorf("error while collecting rows to accounts: %w", err)
	}

	return accounts, nil
}

func (pg *postgres) getDisableAccountsForRetryProvisioning(ctx context.Context) ([]AccountProvision, error) {
	query := "select identity.id as identity_id, first_name, last_name, email, account.id as account_id, username, system_id" +
		" from account join identity on identity.id = identity_id where enable_provisioned_at is not null and enable_committed_at is not null" +
		" and disabled_at < now() and disable_provisioned_at < now() - interval '1 minutes' and disable_committed_at is null"

	rows, err := pg.db.Query(ctx, query)

	if err != nil {
		log.Printf("Error while getting accounts for disable retry provisioning, %s", err)
	}

	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[AccountProvision])

	if err != nil {
		return nil, fmt.Errorf("unable to scan row: %w", err)
	}

	return accounts, nil
}

func (pg *postgres) getDeleteAccountsForProvisioning(ctx context.Context) ([]AccountProvision, error) {
	query := "select identity.id as identity_id, first_name, last_name, email, account.id as account_id, username, system_id" +
		" from account join identity on identity.id = identity_id where disable_provisioned_at is not null and disable_committed_at is not null" +
		" and deleted_at < now() and delete_provisioned_at is null and delete_committed_at is null"

	rows, err := pg.db.Query(ctx, query)

	if err != nil {
		log.Printf("Error while getting accounts delete for provisioning, %s", err)
	}

	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[AccountProvision])

	if err != nil {
		return nil, fmt.Errorf("error while collecting rows to accounts: %w", err)
	}

	return accounts, nil
}

func (pg *postgres) getDeleteAccountsForRetryProvisioning(ctx context.Context) ([]AccountProvision, error) {
	query := "select identity.id as identity_id, first_name, last_name, email, account.id as account_id, username, system_id" +
		" from account join identity on identity.id = identity_id where disable_provisioned_at is not null and disable_committed_at is not null" +
		" and deleted_at < now() and delete_provisioned_at < now() - interval '1 minutes' and delete_committed_at is null"

	rows, err := pg.db.Query(ctx, query)

	if err != nil {
		log.Printf("Error while getting accounts for delete retry provisioning, %s", err)
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
		log.Printf("Failed to mark account as committed: %s", err)
	}
}

func (pg *postgres) markAccountEnableAsProvisioned(id null.String) {
	log.Printf("Mark account enable %s as provisioned", fmt.Sprint(id))
	query := "update account set enable_provisioned_at = now() where id=@id"
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := pg.db.Exec(context.Background(), query, args)

	if err != nil {
		log.Printf("Failed to mark account as provisioned: %s", err)
	}
}

func (pg *postgres) markAccountEnableAsCommitted(id null.String) {
	log.Printf("Mark account enable %s as committed", fmt.Sprint(id))
	query := "update account set enable_committed_at = now() where id=@id"
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := pg.db.Exec(context.Background(), query, args)

	if err != nil {
		log.Printf("Failed to mark account as committed: %s", err)
	}
}

func (pg *postgres) markAccountDisableAsProvisioned(id null.String) {
	log.Printf("Mark account disable %s as provisioned", fmt.Sprint(id))
	query := "update account set disable_provisioned_at = now() where id=@id"
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := pg.db.Exec(context.Background(), query, args)

	if err != nil {
		log.Printf("Failed to mark account as provisioned: %s", err)
	}
}

func (pg *postgres) markAccountDisableAsCommitted(id null.String) {
	log.Printf("Mark account disable %s as committed", fmt.Sprint(id))
	query := "update account set disable_committed_at = now() where id=@id"
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := pg.db.Exec(context.Background(), query, args)

	if err != nil {
		log.Printf("Failed to mark account as committed: %s", err)
	}
}

func (pg *postgres) markAccountDeleteAsProvisioned(id null.String) {
	log.Printf("Mark account delete %s as provisioned", fmt.Sprint(id))
	query := "update account set delete_provisioned_at = now() where id=@id"
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := pg.db.Exec(context.Background(), query, args)

	if err != nil {
		log.Printf("Failed to mark account as provisioned: %s", err)
	}
}

func (pg *postgres) markAccountDeleteAsCommitted(id null.String) {
	log.Printf("Mark account delete %s as committed", fmt.Sprint(id))
	query := "update account set delete_committed_at = now() where id=@id"
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := pg.db.Exec(context.Background(), query, args)

	if err != nil {
		log.Printf("Failed to mark account as committed: %s", err)
	}
}
