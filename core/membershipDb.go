package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (pg *postgres) insertMembership(ctx *gin.Context, createGroupMembershipDto GroupMembershipDto) error {
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

func (pg *postgres) updateMembership(ctx *gin.Context, groupMembershipId string, groupMembershipDto GroupMembershipDto) {
	query := "update group_membership set enabled_at = @enabled_at, disabled_at @disabledAt, deleted_at = @deletedAt where id = @id"
	args := pgx.NamedArgs{
		"enabledAt":  groupMembershipDto.EnabledAt,
		"disabledAt": groupMembershipDto.DisabledAt,
		"deletedAt":  groupMembershipDto.DeletedAt,
		"id":         groupMembershipId,
	}

	_, err := pg.db.Exec(ctx, query, args)

	if err != nil {
		log.Printf("Error while updating membership %s, %s", groupMembershipId, err)
	}
}