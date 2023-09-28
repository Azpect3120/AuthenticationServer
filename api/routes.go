package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
//	"github.com/google/uuid"
)

// Load routes in to the server
func (server *Server) LoadRoutes (database Database) {
	server.router.POST("/createApplication", func (c *gin.Context) { createApplication(c, database) })
	server.router.POST("/createUser", func (c *gin.Context) { createUser(c, database) })
	server.router.POST("/verifyUser", func(c *gin.Context) { verifyUser(c, database) })
}

// Creates a new application in the database
func createApplication (ctx *gin.Context, database Database) {
	var appReq CreateApplicationRequest

	if err := ctx.ShouldBindJSON(&appReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	application := *database.CreateApplication(appReq.Name)

	ctx.JSON(http.StatusCreated, gin.H{ "status": 201, "application": &application })
}

// Create a new user in an application in the database
// User the ApplicationID as a key
func createUser (ctx *gin.Context, database Database) {
	var userReq CreateUserRequest

	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
		return
	}

	hashedPassword, err := HashString(userReq.Password)

	if err != nil {
		panic(err)
	}

	user := *database.CreateUser(userReq.ApplicationID, userReq.Username, hashedPassword)

	ctx.JSON(http.StatusCreated, gin.H{ "status": 201, "user": &user })
}

// Verify a user in the database
// Require ApplicationID as a key
/*
	body: {
		ApplicationID: string(uuid),
		Username: string,
		Password: string
	}

	return: {
		ApplicationID: string(uuid) || nil
		ID: string(uuid) || nil
	}
*/
func verifyUser (ctx *gin.Context, database Database) {
	var verifyReq VerifyUserRequest

	if err := ctx.ShouldBindJSON(&verifyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
	}

	user, err := database.VerifyUser(verifyReq.ApplicationID, verifyReq.Username, verifyReq.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{ "status": 400, "error": err.Error() })
	} else {
		ctx.JSON(http.StatusOK, gin.H{ "status": 200, "user": &user })
	}
}


