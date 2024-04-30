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
		log.Printf("Error while fetching accounts: %s", err)
		c.IndentedJSON(http.StatusNotFound, err)
	} else {
		c.IndentedJSON(http.StatusOK, accounts)
	}
}
