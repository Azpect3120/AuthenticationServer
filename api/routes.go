package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Load routes in to the server
func (server *Server) LoadRoutes () {
	server.router.POST("/createApplication", createApplication)
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
func createApplication (ctx *gin.Context) {

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
