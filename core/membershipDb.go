package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (pg *postgres) insertMembership(ctx *gin.Context, createGroupMembershipDto CreateGroupMembership) error {
	query := "insert into group_membership(identity_id, group_id, enabled_at, disabled_at, deleted_at)" +
		" values (@identityId, @groupId, @enabledAt, @disabledAt, @deletedAt)"
	args := pgx.NamedArgs{
		"identityId": createGroupMembershipDto.IdentityId,
		"groupId":    createGroupMembershipDto.GroupId,
		"enabledAt":  createGroupMembershipDto.EnabledAt,
		"disabledAt": createGroupMembershipDto.DisabledAt,
		"deletedAt":  createGroupMembershipDto.DeletedAt,
	}

	_, err := pg.db.Exec(ctx, query, args)

	if err != nil {
		return fmt.Errorf("error while inserting membership with group %s for identity %s, %w",
			createGroupMembershipDto.GroupId, createGroupMembershipDto.IdentityId, err)
	}

	return nil
}
