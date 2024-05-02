package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/identity", getIdentities)
	router.GET("/identity/:id", getIdentity)

	router.Run(":8083")
}
