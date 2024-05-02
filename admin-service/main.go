package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin"},
	}))

	router.GET("/identity", getIdentities)
	router.GET("/identity/:id", getIdentity)

	router.Run(":8083")
}
