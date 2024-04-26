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
	ID         null.String `json:"id"`
	SystemId   null.String `json:"systemId"`
	IdentityId null.String `json:"identityId"`
}
