package main

import (
	"context"

	"github.com/gin-gonic/gin"
)

func main() {
	dbUrl := "postgres://iamcore:iamcore@localhost:5432/iamcore"

	migratedb(dbUrl)
	NewPG(context.Background(), dbUrl)

	router := gin.Default()

	router.GET("/identity", getIdentities)
	router.POST("/identity", createIdentity)

	router.GET("/account/:identityId", getIdentityAccounts)
	router.POST("/account", createAccount)

	go syncAccounts()
	go commitAccounts()
	router.Run(":8081")
}
