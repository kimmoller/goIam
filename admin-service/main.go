package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-type"},
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE"},
	}))

	router.GET("/identity", getIdentities)
	router.GET("/identity/:id", getIdentity)
	router.POST("/identity", createIdentity)

	router.POST("/membership", createMembership)
	router.DELETE("/membership/:id", deleteMembership)

	router.Run(":8083")
}
