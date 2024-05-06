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

func deleteMembership(ctx *gin.Context) {
	id := ctx.Param("id")
	log.Printf("Got request to delete membership %s", id)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, os.Getenv("CORE_URL")+"membership/"+id, nil)

	if err != nil {
		log.Printf("Error while creating request to delete membership %s, %s", id, err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error while deleting membership %s, %s", id, err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	defer resp.Body.Close()

	if err != nil {
		log.Printf("Error while reading response, %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
