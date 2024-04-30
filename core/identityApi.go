package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getIdentities(c *gin.Context) {
	identities, err := pgInstance.getIdentitiesFromDb(context.Background())

	if err != nil {
		log.Printf("Identities not found %s", err)
		c.IndentedJSON(http.StatusNotFound, err)
	} else {
		c.IndentedJSON(http.StatusOK, identities)
	}
}

func getExtendedIdentities(c *gin.Context) {
	identities, err := pgInstance.getExtendedIdentitiesFromDb(c)
	if err != nil {
		log.Printf("Identities not found %s", err)
		c.IndentedJSON(http.StatusNotFound, err)
	} else {
		c.IndentedJSON(http.StatusOK, identities)
	}
}

func createIdentity(c *gin.Context) {
	var identity Identity

	if err := c.BindJSON(&identity); err != nil {
		log.Print(err)
		return
	}

	pgInstance.insertIdentity(context.Background(), identity)

	c.IndentedJSON(http.StatusCreated, identity)
}
