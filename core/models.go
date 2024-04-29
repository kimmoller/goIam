package main

import "gopkg.in/guregu/null.v3"

type Identity struct {
	ID        string  `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email"`
	Account   Account `json:"account"`
}

type Account struct {
	ID            null.String `json:"id"`
	Username      null.String `json:"username"`
	SystemId      null.String `json:"systemId"`
	IdentityId    null.String `json:"identityId"`
	CreatedAt     null.Time   `json:"createdAt"`
	ProvisionedAt null.Time   `json:"provisionedAt"`
	CommittedAt   null.Time   `json:"committedAt"`
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
