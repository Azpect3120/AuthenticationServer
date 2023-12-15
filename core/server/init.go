package server

import (
	"fmt"
	"strings"

	"github.com/Azpect3120/AuthenticationServer/core/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(port int, connstring string) *model.Server {
	s := &model.Server{
		Router: gin.Default(),
		Config: cors.DefaultConfig(),
		Port:   port,
	}

	// db := database.NewDatabase(connstring)

	s.Config.AllowOrigins = []string{"*"}
	s.Router.Use(cors.New(s.Config))

	return s
}

func AddRoute (s *model.Server, method, endpoint string, handler func(*gin.Context)) {
	method = strings.ToLower(method)

	// Add route to server router
	switch method {
		case "get":
			s.Router.GET(endpoint, handler)
		case "post":
			s.Router.POST(endpoint, handler)
		case "put":
			s.Router.PUT(endpoint, handler)
		case "patch":
			s.Router.PATCH(endpoint, handler)
		case "delete":
			s.Router.DELETE(endpoint, handler)
		default:
			fmt.Printf("Invalid method: %s\n", method)
	}
}

func Listen (s *model.Server) error {
    return s.Router.Run(fmt.Sprintf(":%d", s.Port))
}
