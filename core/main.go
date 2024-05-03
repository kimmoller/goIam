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

	router.GET("/extendedIdentity", getExtendedIdentities)
	router.GET("/extendedIdentity/:id", getExtendedIdentity)

	router.POST("/membership", createGroupMembership)
	router.PATCH("/membership/:id", updateGroupMembership)

	go syncAccounts()
	go commitAccounts()
	router.Run(":8081")
}
