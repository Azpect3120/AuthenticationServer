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
	server.loadCORS()
	return server
}

func (server *Server) Listen () {
	server.router.Run(":" + server.port)
}

func (server *Server) loadCORS () {
	server.router.Use(func (c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
}
