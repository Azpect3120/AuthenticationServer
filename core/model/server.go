package model

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
	Config cors.Config
	Port   string
}

