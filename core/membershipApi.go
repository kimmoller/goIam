package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v3"
)

func containsSystemId(systemIds []string, systemId string) bool {
	for s := range systemIds {
		if systemIds[s] == systemId {
			return true
		}
	}
	return false
}

func getAccountWithSystemId(accounts []Account, systemId string) *Account {
	for a := range accounts {
		if accounts[a].SystemId == systemId {
			return &accounts[a]
		}
	}
	return nil
}

func getAccountDeleteTime(disabledAt null.Time) null.Time {
	var deleteTime null.Time
	if disabledAt.Valid {
		deleteTime.Valid = true
		deleteTime.Time = disabledAt.Time.AddDate(0, 0, 7)
	}
	return deleteTime
}

func createGroupMembership(ctx *gin.Context) {
	var createGroupMembershipDto GroupMembershipDto

	if err := ctx.BindJSON(&createGroupMembershipDto); err != nil {
		log.Printf("Error while reading response to json, %s", err)
		return
	}

	log.Printf("Create group membership for identity %s to group %s", createGroupMembershipDto.IdentityId, createGroupMembershipDto.GroupId)

	systemIds, err := pgInstance.getGroupPermissions(ctx, createGroupMembershipDto.GroupId)
	if err != nil {
		log.Printf("Error while fetching permissions: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	identityId := createGroupMembershipDto.IdentityId

	identity, err := pgInstance.getExtendedIdentityFromDb(ctx, identityId)

	if err != nil {
		log.Printf("Error while fetching identity %s: %s", identityId, err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	for i := range systemIds {
		systemId := systemIds[i]
		existingAccount := getAccountWithSystemId(identity.Accounts, systemId)
		if existingAccount != nil {
			log.Printf("Update existing account %s for identity %s", existingAccount.SystemId, identityId)
			err = updateExistingAccount(ctx, createGroupMembershipDto, identity, existingAccount)
			if err != nil {
				log.Printf("error while updating identity %s account %s: %s", identityId, systemId, err)
				ctx.IndentedJSON(http.StatusInternalServerError, err)
				return
			}
		} else {
			log.Printf("Create new account %s for identity %s", systemId, identityId)
			err = createAccount(ctx, createGroupMembershipDto, identity, systemId)
			if err != nil {
				log.Printf("error while creating identity %s account %s: %s", identityId, systemId, err)
				ctx.IndentedJSON(http.StatusInternalServerError, err)
				return
			}
		}
	}

	err = pgInstance.insertMembership(ctx, createGroupMembershipDto)

	if err != nil {
		log.Printf("Error while fetching identity %s: %s", identityId, err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, "Created identity membership")
}

func getIntervalForAccount(ctx *gin.Context, membershipDto GroupMembershipDto, existingMemberships []GroupMembershipWithGroup,
	account *Account) (*AccountInterval, error) {
	memberships := []GroupMembershipDto{membershipDto}
	for i := range existingMemberships {
		membership := existingMemberships[i]
		systemIds, err := pgInstance.getGroupPermissions(ctx, membership.Group.ID)
		if err != nil {
			return nil, fmt.Errorf("rror while fetching permissions: %s", err)
		}
		if containsSystemId(systemIds, account.SystemId) {
			dto := GroupMembershipDto{
				IdentityId: membership.IdentityId,
				GroupId:    membership.Group.ID,
				EnabledAt:  membership.EnabledAt,
				DisabledAt: membership.DisabledAt,
			}
			memberships = append(memberships, dto)
		}
	}
	reEnable := account.DisabledAt.Valid && account.DisabledAt.Time.Before(time.Now())
	interval := AccountInterval{
		enabledAt:  account.EnabledAt,
		disabledAt: account.DisabledAt,
		deletedAt:  account.DeletedAt,
		reEnable:   reEnable,
	}
	for m := range memberships {
		membership := memberships[m]
		if (!interval.enabledAt.Before(time.Now()) && membership.EnabledAt.Before(interval.enabledAt)) ||
			(reEnable && interval.enabledAt.Before(time.Now())) ||
			(reEnable && membership.EnabledAt.Before(interval.enabledAt)) {
			interval.enabledAt = membership.EnabledAt
		}
		if (interval.disabledAt.Valid && membership.DisabledAt.Valid && membership.DisabledAt.Time.After(interval.disabledAt.Time)) ||
			(!interval.disabledAt.Valid && membership.DisabledAt.Valid) ||
			(interval.disabledAt.Valid && !membership.DisabledAt.Valid) {
			interval.disabledAt = membership.DisabledAt
		}
	}
	interval.deletedAt = getAccountDeleteTime(interval.disabledAt)
	log.Printf("Calculated interval for identity %s account %s with enable time %s, disable time %s, delete time %s, reEnable %s",
		account.IdentityId, account.SystemId, interval.enabledAt, interval.disabledAt.Time, interval.deletedAt.Time, fmt.Sprint(reEnable))
	return &interval, nil
}

func updateExistingAccount(ctx *gin.Context, createGroupMembershipDto GroupMembershipDto,
	identity *ExtendedIdentity, account *Account) error {
	interval, err := getIntervalForAccount(ctx, createGroupMembershipDto, identity.Memberships, account)
	if err != nil {
		return fmt.Errorf("error while getting interval for identity %s account %s", account.IdentityId, err)
	}
	updateAccount := UpdateAccount{
		identityId: identity.ID,
		systemId:   account.SystemId,
		enabledAt:  interval.enabledAt,
		disabledAt: interval.disabledAt,
		deletedAt:  interval.deletedAt,
	}
	log.Printf("Update identity %s account %s with enable time %s, disable time %s, delete time %s",
		identity.ID, account.SystemId, interval.enabledAt, interval.disabledAt.Time, interval.deletedAt.Time)
	return pgInstance.updateAccount(ctx, updateAccount, interval.reEnable)
}

func createAccount(ctx *gin.Context, createGroupMembershipDto GroupMembershipDto, identity *ExtendedIdentity, systemId string) error {
	deleteTime := getAccountDeleteTime(createGroupMembershipDto.DisabledAt)
	account := CreateAccount{
		identityId: identity.ID,
		username:   strings.ToLower(identity.FirstName + identity.LastName),
		systemId:   systemId,
		enabledAt:  createGroupMembershipDto.EnabledAt,
		disabledAt: createGroupMembershipDto.DisabledAt,
		deletedAt:  deleteTime,
	}
	log.Printf("Create identity %s account %s with enable time %s, disable time %s, delete time %s",
		identity.ID, systemId, account.enabledAt, account.disabledAt.Time, account.deletedAt.Time)
	return pgInstance.insertAccount(ctx, account)
}

func updateGroupMembership(ctx *gin.Context) {
	membershipId := ctx.Param("id")
	var groupMembershipDto GroupMembershipDto

	if err := ctx.BindJSON(&groupMembershipDto); err != nil {
		log.Print(err)
		return
	}

	identityId := groupMembershipDto.IdentityId

	identity, err := pgInstance.getExtendedIdentityFromDb(ctx, identityId)

	if err != nil {
		log.Printf("Error while fetching identity %s, %s", identityId, err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	for i := range identity.Accounts {
		account := identity.Accounts[i]
		err := updateExistingAccount(ctx, groupMembershipDto, identity, &account)
		if err != nil {
			log.Printf("Error while updating identity %s account %s, %s", identityId, account.SystemId, err)
			ctx.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
	}

	err = pgInstance.updateMembership(ctx, membershipId, groupMembershipDto)
	if err != nil {
		log.Println(err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, "Updated membership "+membershipId)
}
