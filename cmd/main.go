package main

import (
	"github.com/Azpect3120/AuthenticationServer/core/applications"
	"github.com/Azpect3120/AuthenticationServer/core/database"
	"github.com/Azpect3120/AuthenticationServer/core/model"
	s "github.com/Azpect3120/AuthenticationServer/core/server"
	"github.com/gin-gonic/gin"
)

const DB_CONN_STRING string = "postgres://lnnhgkzj:kR6RcBmeiyhkkkEnWKmfnCJw3oovszRQ@bubble.db.elephantsql.com/lnnhgkzj"

func main() {
	server := s.NewServer(3000, "")

    db := database.NewDatabase(DB_CONN_STRING)

    
    // `POST` v2/applications -> Create an application
    s.AddRoute(server, "post", "/v2/applications", func(ctx *gin.Context) {
        var req model.CreateApplicationRequest
        if err := ctx.ShouldBindJSON(&req); err != nil {
            ctx.JSON(404, gin.H{ "error": err.Error() })
        }

        message := applications.MatchColumns(&req.Columns)
        app := applications.New(req.Name, req.Columns)

        if err := applications.Insert(db, app); err != nil {
            ctx.JSON(500, gin.H{ "status": 500, "error": err.Error() })
            return
        }

        ctx.JSON(200, gin.H{ "status": 200, "message": message, "application": app })
    })
    

    s.Listen(server)
}


