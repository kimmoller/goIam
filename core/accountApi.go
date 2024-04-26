package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getIdentityAccounts(c *gin.Context) {
	identityId := c.Param("identityId")
	accounts, err := pgInstance.getIdentityAccountsFromDb(context.Background(), identityId)

	if err != nil {
		log.Printf("Accounts not found %s", err)
		c.IndentedJSON(http.StatusNotFound, err)
	} else {
		c.IndentedJSON(http.StatusOK, accounts)
	}
}

func createAccount(c *gin.Context) {
	var account Account

	if err := c.BindJSON(&account); err != nil {
		log.Print(err)
		return
	}

	pgInstance.insertAccount(context.Background(), account)

	c.IndentedJSON(http.StatusCreated, account)
}
