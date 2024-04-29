package main

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type SimpleIdentity struct {
	id        string
	firstName string
	lastName  string
	email     string
}

type CreateAccount struct {
	identityId string
	username   string
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

type CreateGroupMembership struct {
	IdentityId string    `json:"identityId"`
	GroupId    string    `json:"groupId"`
	EnabledAt  time.Time `json:"enabledAt"`
	DisabledAt null.Time `json:"disabledAt"`
	DeletedAt  null.Time `json:"deletedAt"`
}
