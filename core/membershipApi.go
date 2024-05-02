package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v3"
)

func createGroupMembership(ctx *gin.Context) {
	var createGroupMembershipDto GroupMembershipDto

	if err := ctx.BindJSON(&createGroupMembershipDto); err != nil {
		log.Print(err)
		return
	}

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

	// TODO: Add check for existing accounts
	for index := range systemIds {
		var deleteTime null.Time
		if createGroupMembershipDto.DisabledAt.Valid {
			deleteTime.Valid = true
			deleteTime.Time = createGroupMembershipDto.DisabledAt.Time.AddDate(0, 0, 7)
		}
		account := CreateAccount{
			identityId: identityId,
			username:   strings.ToLower(identity.FirstName + identity.LastName),
			systemId:   systemIds[index],
			enabledAt:  createGroupMembershipDto.EnabledAt,
			disabledAt: createGroupMembershipDto.DisabledAt,
			deletedAt:  deleteTime,
		}
		pgInstance.insertAccount(ctx, account)
	}

	err = pgInstance.insertMembership(ctx, createGroupMembershipDto)

	if err != nil {
		log.Printf("Error while fetching identity %s: %s", identityId, err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func updateGroupMembership(ctx *gin.Context) {
	membershipId := ctx.Param("id")
	var groupMembershipDto GroupMembershipDto

	if err := ctx.BindJSON(&groupMembershipDto); err != nil {
		log.Print(err)
		return
	}

	systemIds, err := pgInstance.getGroupPermissions(ctx, groupMembershipDto.GroupId)
	if err != nil {
		log.Printf("Error while fetching permissions: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
	}

	identityId := groupMembershipDto.IdentityId

	accounts, err := pgInstance.getIdentityAccountsForSystemIdFromDb(ctx, identityId, systemIds)

	if err != nil {
		log.Printf("Error while fetching identity %s accounts, %s", identityId, err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
	}

	for i := range accounts {
		account := accounts[i]
		updateAccount := UpdateAccount{
			identityId: identityId,
			systemId:   account.SystemId.String,
			enabledAt:  groupMembershipDto.EnabledAt,
			disabledAt: groupMembershipDto.DisabledAt,
		}
		// TODO: Add validation here to only re-enable account when necessary
		pgInstance.updateAccount(ctx, updateAccount, true)
	}

	pgInstance.updateMembership(ctx, membershipId, groupMembershipDto)
}
