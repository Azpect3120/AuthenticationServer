package model

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	port string
}

func CreateServer(port string, database Database) *Server {
	server := &Server{
		router: gin.Default(),
		port: port,
	}
	server.LoadRoutes(database)
	return server
}

func (server *Server) Listen () {
	server.router.Run(":" + server.port)
}
