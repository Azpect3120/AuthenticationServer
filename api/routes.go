package api

import (
	"net/http"
//	"fmt"

	"github.com/gin-gonic/gin"
)

// Load routes in to the server
func (server *Server) LoadRoutes (database Database) {
	server.router.POST("/createApplication", func (ctx *gin.Context) { createApplication(ctx, database) })
	server.router.POST("/createUser", func (ctx *gin.Context) { createUser(ctx, database) })
}

// Request map: createApplication
type CreateApplicationRequest struct {
	Name	string		`json:"name"`
}

// Request map: createUser
type CreateUserRequest struct {
	ApplicationID	uuid.UUID	`json:"applicationID"`
	Username		string		`json:"username"`
	Password		string		`json:"password"`
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
func createApplication (ctx *gin.Context, database Database) {
	var appReq CreateApplicationRequest

	if err := ctx.ShouldBindJSON(&appReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
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
func createUser (ctx *gin.Context, database Database) {
	var userReq CreateUserRequest

	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	user := *database.CreateUser(userReq.applicationID, userReq.username, userReq.password)

	ctx.JSON(http.StatusCreated, gin.H{ "status": 201, "user": &user })
}
