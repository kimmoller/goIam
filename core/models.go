package main

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type Identity struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type Account struct {
	ID                   string    `json:"id"`
	Username             string    `json:"username"`
	SystemId             string    `json:"systemId"`
	IdentityId           string    `json:"identityId"`
	CreatedAt            time.Time `json:"createdAt"`
	ProvisionedAt        null.Time `json:"provisionedAt"`
	CommittedAt          null.Time `json:"committedAt"`
	EnabledAt            time.Time `json:"enabledAt"`
	EnableProvisionedAt  null.Time `json:"enableProvisionedAt"`
	EnableCommittedAt    null.Time `json:"enableCommittedAt"`
	DisabledAt           null.Time `json:"disabledAt"`
	DisableProvisionedAt null.Time `json:"disableProvisionedAt"`
	DisableCommittedAt   null.Time `json:"disableCommittedAt"`
	DeletedAt            null.Time `json:"deletedAt"`
	DeleteProvisionedAt  null.Time `json:"deleteProvisionedAt"`
	DeleteCommittedAt    null.Time `json:"deleteCommittedAt"`
}

type PermissionGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PermissionToGroup struct {
	ID           string
	PermissionId string
	GroupId      string
}

type GroupMembership struct {
	ID         string    `json:"id"`
	GroupId    string    `json:"groupId"`
	IdentityId string    `json:"identityId"`
	CreatedAt  time.Time `json:"createdAt"`
	EnabledAt  time.Time `json:"enabledAt"`
	DisabledAt null.Time `json:"disabledAt"`
}
