package main

import (
	"context"

	"github.com/gin-gonic/gin"
)

func main() {
	dbUrl := "postgres://iamcore:iamcore@0.0.0.0:9080/iamcore"

	migratedb(dbUrl)
	NewPG(context.Background(), dbUrl)

	router := gin.Default()

	router.GET("/identity", getIdentities)
	router.POST("/identity", createIdentity)

	router.GET("/account/:identityId", getIdentityAccounts)

	router.POST("/membership", createGroupMembership)
	router.PATCH("/membership/:id", updateGroupMembership)

	go syncAccounts()
	go commitAccounts()
	router.Run(":8081")
}
