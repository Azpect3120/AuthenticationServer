package api

import (
	"net/http"
//	"fmt"

	"github.com/gin-gonic/gin"
)

// Load routes in to the server
func (server *Server) LoadRoutes (database Database) {
	server.router.POST("/createApplication", func (ctx *gin.Context) { createApplication(ctx, database)})
	server.router.POST("/createUser", createUser)
}

/*
	Creates a new application in the database

	body: {
		name: string
	}
	 
	return: {
		ID: uuid,
		Name: string,
		Key: uuid
	}
	
*/
type CreateApplicationRequest struct {
	Name string `json:"name"`
}

func createApplication (ctx *gin.Context, database Database) {
	var appReq CreateApplicationRequest

	if err := ctx.ShouldBindJSON(&appReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
		return
	}

	application := *database.CreateApplication(appReq.Name)

	ctx.JSON(http.StatusCreated, gin.H{ "status": 201, "application": &application })

}

/*
	Create a new user in an application in the database

	body: {
		ApplicationID: uuid,
		Username: string,
		Password: string,
	}

	return: {
		ID: uuid,
		ApplicationID: uuid,
		Username: string,
		Password: string,
	}

*/
func createUser (ctx *gin.Context) {

}
