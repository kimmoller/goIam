package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (pg *postgres) insertMembership(ctx *gin.Context, createGroupMembershipDto GroupMembershipDto) error {
	query := "insert into group_membership(identity_id, group_id, enabled_at, disabled_at)" +
		" values (@identityId, @groupId, @enabledAt, @disabledAt)"
	args := pgx.NamedArgs{
		"identityId": createGroupMembershipDto.IdentityId,
		"groupId":    createGroupMembershipDto.GroupId,
		"enabledAt":  createGroupMembershipDto.EnabledAt,
		"disabledAt": createGroupMembershipDto.DisabledAt,
	}

	_, err := pg.db.Exec(ctx, query, args)

	if err != nil {
		return fmt.Errorf("error while inserting membership with group %s for identity %s, %w",
			createGroupMembershipDto.GroupId, createGroupMembershipDto.IdentityId, err)
	}

	return nil
}

func (pg *postgres) updateMembership(ctx *gin.Context, groupMembershipId string, groupMembershipDto GroupMembershipDto) error {
	query := "update group_membership set enabled_at = @enabledAt, disabled_at = @disabledAt where id = @id"
	args := pgx.NamedArgs{
		"enabledAt":  groupMembershipDto.EnabledAt,
		"disabledAt": groupMembershipDto.DisabledAt,
		"id":         groupMembershipId,
	}

	_, err := pg.db.Exec(ctx, query, args)

	if err != nil {
		return fmt.Errorf("error while updating membership %s, %s", groupMembershipId, err)
	}
	return nil
}

func (pg *postgres) findMembership(ctx *gin.Context, membershipId string) (*GroupMembership, error) {
	query := "select * from group_membership where id = @id"
	args := pgx.NamedArgs{
		"id": membershipId,
	}

	row := pg.db.QueryRow(ctx, query, args)
	var membership GroupMembership
	err := row.Scan(&membership.ID, &membership.GroupId, &membership.IdentityId,
		&membership.EnabledAt, &membership.DisabledAt, &membership.DisabledAt)

	if err != nil {
		return nil, fmt.Errorf("error while getting membership %s, %w", membershipId, err)
	}
	return &membership, nil
}
