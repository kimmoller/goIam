package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func getIdentities(ctx *gin.Context) {
	response, err := http.Get(os.Getenv("CORE_URL") + "identity")
	if err != nil {
		log.Printf("Error while requesting identites, %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error while reading response, %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
	}
	log.Printf("Got response %s", body)

	ctx.Data(http.StatusOK, "application/json", body)
}

func getIdentity(ctx *gin.Context) {
	identityId := ctx.Param("id")
	log.Printf("Fetching identity %s", identityId)
	response, err := http.Get(os.Getenv("CORE_URL") + "extendedIdentity/" + identityId)
	if err != nil {
		log.Printf("Error while requesting identity %s, %s", identityId, err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error while reading response, %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
	}
	log.Printf("Got response %s", body)

	ctx.Data(http.StatusOK, "application/json", body)
}

func createIdentity(ctx *gin.Context) {
	response, err := http.Post(os.Getenv("CORE_URL")+"identity", "application/json", ctx.Request.Body)

	if err != nil {
		log.Printf("Error while creating identity %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error while reading response, %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
	}
	ctx.Data(http.StatusOK, "application/json", body)
}
