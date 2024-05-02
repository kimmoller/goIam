package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	dbUrl := os.Getenv("DB_URL")

	migratedb(dbUrl)
	NewPG(context.Background(), dbUrl)

	router := gin.Default()

	router.GET("/identity", getIdentities)
	router.POST("/identity", createIdentity)

	router.GET("/extendedIdentities", getExtendedIdentities)

	router.GET("/account/:identityId", getIdentityAccounts)

	router.POST("/membership", createGroupMembership)
	router.PATCH("/membership/:id", updateGroupMembership)

	go syncAccounts()
	go commitAccounts()
	router.Run(":8081")
}
