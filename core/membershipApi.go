package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createGroupMembership(ctx *gin.Context) {
	var createGroupMembershipDto CreateGroupMembership

	if err := ctx.BindJSON(&createGroupMembershipDto); err != nil {
		log.Print(err)
		return
	}

	log.Printf("Create membership for identity %s with group %s, enabledAt %s, disabledAt %s and deleteAt %s",
		createGroupMembershipDto.IdentityId,
		createGroupMembershipDto.GroupId,
		createGroupMembershipDto.EnabledAt,
		createGroupMembershipDto.DisabledAt.Time,
		createGroupMembershipDto.DeletedAt.Time)

	systemIds, err := pgInstance.getGroupPermissions(ctx, createGroupMembershipDto.GroupId)
	if err != nil {
		log.Printf("Error while fetching permissions: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
	}

	identityId := createGroupMembershipDto.IdentityId

	identity, err := pgInstance.getIdentityFromDb(ctx, identityId)

	if err != nil {
		log.Printf("Error while fetching identity %s: %s", identityId, err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
	}

	for index := range systemIds {
		account := CreateAccount{
			identityId: identityId,
			username:   identity.firstName + identity.lastName,
			systemId:   systemIds[index],
			enabledAt:  createGroupMembershipDto.EnabledAt,
			disabledAt: createGroupMembershipDto.DisabledAt,
			deletedAt:  createGroupMembershipDto.DeletedAt,
		}
		pgInstance.insertAccount(ctx, account)
	}

	err = pgInstance.insertMembership(ctx, createGroupMembershipDto)

	if err != nil {
		log.Printf("Error while fetching identity %s: %s", identityId, err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
	}
}
