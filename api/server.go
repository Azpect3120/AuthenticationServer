package api

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	port string
}

func CreateServer(port string) *Server {
	server := &Server{
		router: gin.Default(),
		port: port,
	}
	server.LoadRoutes()
	return server
}

func (server *Server) Listen () {
	server.router.Run(":" + server.port)
}
