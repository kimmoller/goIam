package main

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type CreateIdentity struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type CreateMembership struct {
	IdentityId string    `json:"identityId"`
	GroupId    string    `json:"groupId"`
	EnabledAt  time.Time `json:"enabledAt"`
	DisabledAt null.Time `json:"disabledAt"`
}
