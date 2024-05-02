package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

func (pg *postgres) getIdentityFromDb(ctx context.Context, identityId string) (*Identity, error) {
	query := "select * from identity where id=@identityId"
	args := pgx.NamedArgs{
		"identityId": identityId,
	}
	row := pg.db.QueryRow(ctx, query, args)
	var identity Identity
	err := row.Scan(&identity.ID, &identity.FirstName, &identity.LastName, &identity.Email)

	if err != nil {
		return nil, fmt.Errorf("error while scanning identity %s, %w", identityId, err)
	}

	return &identity, nil
}

func (pg *postgres) getIdentitiesFromDb(ctx context.Context) ([]Identity, error) {
	query := "select * from identity"
	rows, err := pg.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error while fetching identities %w", err)
	}
	identities, err := pgx.CollectRows(rows, pgx.RowToStructByName[Identity])
	if err != nil {
		return nil, fmt.Errorf("error while scanning identities %w", err)
	}
	return identities, nil
}

func (pg *postgres) getExtendedIdentityFromDb(ctx context.Context, identityId string) (*ExtendedIdentity, error) {
	query := "select * from identity where id = @identityId"
	args := pgx.NamedArgs{
		"identityId": identityId,
	}
	identities, err := pg.getExtendedIdentitiesFromDbWithQuery(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("error while fetching identity %s, %w", identityId, err)
	}
	return &identities[0], nil
}

func (pg *postgres) getExtendedIdentitiesFromDb(ctx context.Context) ([]ExtendedIdentity, error) {
	query := "select * from identity"
	return pg.getExtendedIdentitiesFromDbWithQuery(ctx, query, nil)
}

func (pg *postgres) getExtendedIdentitiesFromDbWithQuery(ctx context.Context, query string, args pgx.NamedArgs) ([]ExtendedIdentity, error) {
	rows, err := pg.db.Query(ctx, query, args)

	if err != nil {
		log.Printf("Unable to get identities, %s", err)
		return nil, err
	}

	identities := []ExtendedIdentity{}
	var identity Identity
	_, err = pgx.ForEachRow(rows, []any{&identity.ID, &identity.FirstName, &identity.LastName, &identity.Email}, func() error {
		accountQuery := "select * from account where identity_id = @identityId"
		args := pgx.NamedArgs{
			"identityId": identity.ID,
		}
		accountRows, err := pg.db.Query(ctx, accountQuery, args)
		if err != nil {
			return fmt.Errorf("error while fetching identity %s accounts, %w", identity.ID, err)
		}
		accounts, err := pgx.CollectRows(accountRows, pgx.RowToStructByName[Account])
		if err != nil {
			return fmt.Errorf("error while scanning identity %s accounts, %w", identity.ID, err)
		}

		membershipQuery := "select gm.id, identity_id, pg.id, name, enabled_at, disabled_at, deleted_at from group_membership gm" +
			" left join permission_group pg on gm.group_id=pg.id where identity_id = @identityId"
		membershipRows, err := pg.db.Query(ctx, membershipQuery, args)
		if err != nil {
			return fmt.Errorf("error while fetching identity %s memberships, %w", identity.ID, err)
		}

		var memberships []GroupMembershipWithGroup
		for membershipRows.Next() {
			var membership GroupMembershipWithGroup
			err := membershipRows.Scan(&membership.ID, &membership.IdentityId, &membership.Group.ID, &membership.Group.Name,
				&membership.EnabledAt, &membership.DisabledAt, &membership.DeletedAt)
			if err != nil {
				return fmt.Errorf("error while scanning identity %s memberships, %w", identity.ID, err)
			}
			memberships = append(memberships, membership)
		}

		extendedIdentity := ExtendedIdentity{
			ID:          identity.ID,
			FirstName:   identity.FirstName,
			LastName:    identity.LastName,
			Email:       identity.Email,
			Accounts:    accounts,
			Memberships: memberships,
		}

		identities = append(identities, extendedIdentity)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error while fetching extended identities %w", err)
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
