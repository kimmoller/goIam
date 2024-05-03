package main

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type ExtendedIdentity struct {
	ID                  string                     `json:"id"`
	FirstName           string                     `json:"firstName"`
	LastName            string                     `json:"lastName"`
	Email               string                     `json:"email"`
	Accounts            []Account                  `json:"accounts"`
	Memberships         []GroupMembershipWithGroup `json:"memberships"`
	DisabledMemberships []GroupMembershipWithGroup `json:"disabledMemberships"`
}

type CreateAccount struct {
	identityId string
	username   string
	systemId   string
	enabledAt  time.Time
	disabledAt null.Time
	deletedAt  null.Time
}

type UpdateAccount struct {
	identityId string
	systemId   string
	enabledAt  time.Time
	disabledAt null.Time
	deletedAt  null.Time
}

type AccountProvision struct {
	IdentityID string      `json:"identityId"`
	FirstName  string      `json:"firstName"`
	LastName   string      `json:"lastName"`
	Email      string      `json:"email"`
	AccountID  null.String `json:"accountId"`
	Username   null.String `json:"username"`
	SystemId   null.String `json:"systemId"`
}

type GroupMembershipDto struct {
	IdentityId string    `json:"identityId"`
	GroupId    string    `json:"groupId"`
	EnabledAt  time.Time `json:"enabledAt"`
	DisabledAt null.Time `json:"disabledAt"`
}

type GroupMembershipWithGroup struct {
	ID         string          `json:"id"`
	IdentityId string          `json:"identityId"`
	Group      PermissionGroup `json:"group"`
	EnabledAt  time.Time       `json:"enabledAt"`
	DisabledAt null.Time       `json:"disabledAt"`
}

type AccountInterval struct {
	enabledAt  time.Time
	disabledAt null.Time
	deletedAt  null.Time
	reEnable   bool
}
