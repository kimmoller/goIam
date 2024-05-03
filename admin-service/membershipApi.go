package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func createMembership(ctx *gin.Context) {
	response, err := http.Post(os.Getenv("CORE_URL")+"membership", "application/json", ctx.Request.Body)

	if err != nil {
		log.Printf("Error while creating identity membership, %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Printf("Error while reading response, %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.Data(http.StatusOK, "application/json", body)
}
